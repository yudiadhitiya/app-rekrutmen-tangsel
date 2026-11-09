package sheets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/store"
)

// testPrivateKeyPEM generates an RSA key for the tests to sign fake JWTs with.
func testPrivateKeyPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
}

// fakeSheets is an in-memory stand-in for the Google Sheets API. It stores
// each tab as a grid of cells so tests can exercise the real code paths.
type fakeSheets struct {
	mu   sync.Mutex
	grid map[string][][]string
	// tokenCalls counts OAuth token requests.
	tokenCalls int
}

func newFakeSheets() *fakeSheets {
	return &fakeSheets{grid: map[string][][]string{}}
}

func (f *fakeSheets) tab(name string) [][]string {
	if _, ok := f.grid[name]; !ok {
		f.grid[name] = nil
	}
	return f.grid[name]
}

func (f *fakeSheets) serve(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/token":
			f.mu.Lock()
			f.tokenCalls++
			f.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"access_token":"fake-token","expires_in":3600}`)
			return
		}

		// Sheet endpoints look like /x/values/<Tab>... ; the spreadsheet id is
		// the first path segment.
		rest, ok := strings.CutPrefix(path, "/x/values/")
		if !ok {
			if path == "/x" && r.Method == http.MethodPost {
				// Add a sheet.
				fmt.Fprint(w, `{}`)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error":"not found: `+path+`"}`)
			return
		}

		switch {
		case strings.HasSuffix(rest, ":append"):
			name := strings.TrimSuffix(rest, ":append")
			var body struct {
				Values [][]string `json:"values"`
			}
			json.NewDecoder(r.Body).Decode(&body)
			f.mu.Lock()
			for _, row := range body.Values {
				f.grid[name] = append(f.grid[name], row)
			}
			f.mu.Unlock()
			fmt.Fprint(w, `{}`)
			return

		case strings.HasSuffix(rest, ":clear"):
			// Form: <Tab>!A5:Z5:clear
			spec := strings.TrimSuffix(rest, ":clear")
			name := spec[:strings.Index(spec, "!")]
			rowNum := rowFromSpec(spec)
			f.mu.Lock()
			grid := f.grid[name]
			if rowNum-1 < len(grid) {
				f.grid[name][rowNum-1] = make([]string, len(grid[rowNum-1]))
			}
			f.mu.Unlock()
			fmt.Fprint(w, `{}`)
			return

		case r.Method == http.MethodPut:
			// Form: <Tab>!A5:Z5
			name := rest[:strings.Index(rest, "!")]
			rowNum := rowFromSpec(rest)
			var body struct {
				Values [][]string `json:"values"`
			}
			json.NewDecoder(r.Body).Decode(&body)
			f.mu.Lock()
			grid := f.grid[name]
			for len(grid) < rowNum {
				grid = append(grid, nil)
			}
			if len(body.Values) > 0 {
				grid[rowNum-1] = body.Values[0]
			}
			f.grid[name] = grid
			f.mu.Unlock()
			fmt.Fprint(w, `{}`)
			return

		case r.Method == http.MethodGet:
			name := rest[:strings.Index(rest, "!")]
			f.mu.Lock()
			values := f.grid[name]
			f.mu.Unlock()
			if values == nil {
				// Mirrors the API error for a missing tab.
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `{"error":{"message":"Unable to parse range: `+name+`"}}`)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"values": values})
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":"not found: `+path+`"}`)
	}))
}

// rowFromSpec extracts the 1-based row number from a range like "Rekrutmen!A5:Z5".
func rowFromSpec(spec string) int {
	start := strings.Index(spec, "!A") + 2
	end := strings.Index(spec[start:], ":")
	num, _ := strconv.Atoi(spec[start : start+end])
	return num
}

// newTestStore builds a Store pointed at the fake API.
func newTestStore(t *testing.T) (*Store, *fakeSheets) {
	t.Helper()
	fake := newFakeSheets()
	srv := fake.serve(t)
	t.Cleanup(srv.Close)

	keyFile := filepath.Join(t.TempDir(), "key.json")
	key := map[string]string{
		"client_email": "svc@example.iam.gserviceaccount.com",
		"private_key":  testPrivateKeyPEM(t),
		"token_uri":    srv.URL + "/token",
	}
	b, _ := json.Marshal(key)
	if err := os.WriteFile(keyFile, b, 0o600); err != nil {
		t.Fatal(err)
	}

	s, err := New("x", keyFile)
	if err != nil {
		t.Fatal(err)
	}
	// Point the API base URLs at the fake server.
	s.baseURL = srv.URL
	return s, fake
}

func TestSheetsEnsureTabsCreatesHeaders(t *testing.T) {
	s, fake := newTestStore(t)
	if err := s.Ping(); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	fake.mu.Lock()
	defer fake.mu.Unlock()
	reks := fake.grid["Rekrutmen"]
	if len(reks) == 0 {
		t.Fatal("expected the Rekrutmen header row to be created")
	}
	if reks[0][0] != "id" || reks[0][1] != "nama" {
		t.Errorf("unexpected header: %v", reks[0])
	}
	if _, ok := fake.grid["Aktivitas"]; !ok {
		t.Error("expected the Aktivitas tab to be created")
	}
}

func TestSheetsRekrutmenCRUD(t *testing.T) {
	s, _ := newTestStore(t)
	if err := s.Ping(); err != nil {
		t.Fatal(err)
	}

	r := &model.Rekrutmen{
		ID:                "r1",
		Nama:              "Ahmad Fauzi",
		NoWA:              "628123456789",
		Alamat:            "Jl Melati 5",
		RT:                "003",
		RW:                "005",
		Kecamatan:         "Serpong",
		Kelurahan:         "Buaran",
		Aktivitas:         "Wiraswasta",
		TingkatPenerimaan: "Terbuka",
		Metode:            "Halaqoh",
		Status:            "Berhasil",
		UPL:               "UPLS",
	}
	if err := s.CreateRekrutmen(r); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := s.GetRekrutmen("r1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Nama != "Ahmad Fauzi" || got.NoWA != "628123456789" || got.UPL != "UPLS" {
		t.Errorf("round-trip mismatch: %+v", got)
	}
	if got.Kecamatan != "Serpong" || got.Kelurahan != "Buaran" {
		t.Errorf("wilayah mismatch: %+v", got)
	}

	// Duplicate number is rejected.
	dup := *r
	dup.ID = "r2"
	if err := s.CreateRekrutmen(&dup); err != store.ErrDuplicateWA {
		t.Errorf("expected ErrDuplicateWA, got %v", err)
	}

	// Update.
	r.Nama = "Ahmad Update"
	r.Status = "Proses"
	if err := s.UpdateRekrutmen(r); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, err = s.GetRekrutmen("r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Nama != "Ahmad Update" || got.Status != "Proses" {
		t.Errorf("update not applied: %+v", got)
	}

	// Delete.
	if err := s.DeleteRekrutmen("r1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.GetRekrutmen("r1"); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// Deleting one row must not corrupt the rows after it.
func TestSheetsDeleteKeepsOtherRowsIntact(t *testing.T) {
	s, _ := newTestStore(t)
	if err := s.Ping(); err != nil {
		t.Fatal(err)
	}
	for i, id := range []string{"r1", "r2", "r3"} {
		r := &model.Rekrutmen{
			ID: id, Nama: "Nama " + id, NoWA: "6281111111" + strconv.Itoa(i),
			Kecamatan: "Serpong", Aktivitas: "Wiraswasta",
			TingkatPenerimaan: "Terbuka", Metode: "Halaqoh", Status: "Berhasil",
		}
		if err := s.CreateRekrutmen(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.DeleteRekrutmen("r1"); err != nil {
		t.Fatal(err)
	}
	// r2 and r3 must still be findable.
	if _, err := s.GetRekrutmen("r2"); err != nil {
		t.Errorf("r2 lost after deleting r1: %v", err)
	}
	if _, err := s.GetRekrutmen("r3"); err != nil {
		t.Errorf("r3 lost after deleting r1: %v", err)
	}
	got, err := s.GetRekrutmen("r3")
	if err != nil {
		t.Fatal(err)
	}
	if got.Nama != "Nama r3" {
		t.Errorf("r3 data corrupted: %+v", got)
	}
}

func TestSheetsListAndFilter(t *testing.T) {
	s, _ := newTestStore(t)
	if err := s.Ping(); err != nil {
		t.Fatal(err)
	}
	a := &model.Rekrutmen{
		ID: "r1", Nama: "Budi Santoso", NoWA: "628111111111",
		Kecamatan: "Ciputat", Kelurahan: "Cipayung", Aktivitas: "Pedagang",
		TingkatPenerimaan: "Terbuka", Metode: "Program Antara", Status: "Pending",
	}
	b := &model.Rekrutmen{
		ID: "r2", Nama: "Siti Aminah", NoWA: "628222222222",
		Kecamatan: "Serpong", Kelurahan: "Buaran", Aktivitas: "Pedagang",
		TingkatPenerimaan: "Terbuka", Metode: "Halaqoh", Status: "Berhasil",
	}
	s.CreateRekrutmen(a)
	s.CreateRekrutmen(b)

	cases := []struct {
		name string
		f    model.Filter
		want int
	}{
		{"all", model.Filter{}, 2},
		{"name", model.Filter{Q: "siti"}, 1},
		{"wa", model.Filter{Q: "628111"}, 1},
		{"kecamatan", model.Filter{Kecamatan: "Ciputat"}, 1},
		{"kelurahan", model.Filter{Kelurahan: "Buaran"}, 1},
		{"metode", model.Filter{Metode: "Program Antara"}, 1},
		{"status", model.Filter{Status: "Pending"}, 1},
	}
	for _, c := range cases {
		data, total, err := s.ListRekrutmen(c.f)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if total != c.want || len(data) != c.want {
			t.Errorf("%s: expected %d, got total=%d len=%d", c.name, c.want, total, len(data))
		}
	}
}

func TestSheetsFindByWhatsApp(t *testing.T) {
	s, _ := newTestStore(t)
	s.Ping()
	r := &model.Rekrutmen{
		ID: "r1", Nama: "Ahmad", NoWA: "628123456789",
		Kecamatan: "Serpong", Aktivitas: "Wiraswasta",
		TingkatPenerimaan: "Terbuka", Metode: "Halaqoh", Status: "Berhasil",
	}
	s.CreateRekrutmen(r)
	got, err := s.FindRekrutmenByWA("628123456789")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "r1" {
		t.Errorf("expected r1, got %q", got.ID)
	}
	if _, err := s.FindRekrutmenByWA("628000000000"); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSheetsAktivitasMaster(t *testing.T) {
	s, _ := newTestStore(t)
	s.Ping()
	a := &model.Aktivitas{ID: "a1", Nama: "Wiraswasta", Aktif: true}
	if err := s.CreateAktivitas(a); err != nil {
		t.Fatal(err)
	}
	list, err := s.ListAktivitas()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Nama != "Wiraswasta" {
		t.Fatalf("unexpected list: %+v", list)
	}

	// Duplicate name is rejected.
	if err := s.CreateAktivitas(&model.Aktivitas{ID: "a2", Nama: "Wiraswasta"}); err != store.ErrDuplicateName {
		t.Errorf("expected ErrDuplicateName, got %v", err)
	}

	// Rename cascades to the rekrutmen rows.
	r := &model.Rekrutmen{
		ID: "r1", Nama: "Ahmad", NoWA: "628123456789",
		Kecamatan: "Serpong", Aktivitas: "Wiraswasta",
		TingkatPenerimaan: "Terbuka", Metode: "Halaqoh", Status: "Berhasil",
	}
	s.CreateRekrutmen(r)
	a.Nama = "Wiraswasta Tunggal"
	if err := s.UpdateAktivitas(a); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetRekrutmen("r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Aktivitas != "Wiraswasta Tunggal" {
		t.Errorf("cascade rename failed: %q", got.Aktivitas)
	}

	// An activity in use cannot be deleted.
	if err := s.DeleteAktivitas("a1"); err != store.ErrAktivitasUsed {
		t.Errorf("expected ErrAktivitasUsed, got %v", err)
	}
}

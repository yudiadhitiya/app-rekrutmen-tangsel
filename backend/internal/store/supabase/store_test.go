package supabase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/store"
)

// fakePostgREST is an in-memory stand-in for the PostgREST API.
type fakePostgREST struct {
	rekrutmen []map[string]any
	aktivitas []map[string]any
	// nextID generates unique ids for inserted rows.
	nextID int
	// failWith, when set, is returned as a 409 unique violation.
	failWith string
}

func newFakePostgREST() *fakePostgREST {
	return &fakePostgREST{nextID: 1}
}

func (f *fakePostgREST) serve(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		table := strings.TrimPrefix(r.URL.Path, "/rest/v1/")
		q := r.URL.Query()

		if f.failWith != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, `{"code":"23505","message":%q}`, f.failWith)
			return
		}

		switch r.Method {
		case http.MethodGet:
			rows := f.table(table)
			// Apply the eq. filters the store sends.
			for key, vals := range q {
				if !strings.HasPrefix(vals[0], "eq.") {
					continue
				}
				want := strings.TrimPrefix(vals[0], "eq.")
				var kept []map[string]any
				for _, row := range rows {
					if fmt.Sprint(row[key]) == want {
						kept = append(kept, row)
					}
				}
				rows = kept
			}
			// Count header for the total.
			total := len(rows)
			if offset := q.Get("offset"); offset != "" {
				if off, err := strconv.Atoi(offset); err == nil && off < len(rows) {
					rows = rows[off:]
				} else {
					rows = nil
				}
			}
			if limit := q.Get("limit"); limit != "" {
				if lim, err := strconv.Atoi(limit); err == nil && lim < len(rows) {
					rows = rows[:lim]
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Range", fmt.Sprintf("0-%d/%d", len(rows)-1, total))
			json.NewEncoder(w).Encode(rows)

		case http.MethodPost:
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if _, ok := body["id"]; !ok {
				body["id"] = fmt.Sprintf("id-%d", f.nextID)
				f.nextID++
			}
			rows := append(f.table(table), body)
			f.setTable(table, rows)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode([]map[string]any{body})

		case http.MethodPatch:
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			rows := f.table(table)
			var updated []map[string]any
			for _, row := range rows {
				if matchRow(row, q) {
					for k, v := range body {
						row[k] = v
					}
					updated = append(updated, row)
				}
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updated)

		case http.MethodDelete:
			rows := f.table(table)
			var kept []map[string]any
			var deleted []map[string]any
			for _, row := range rows {
				if matchRow(row, q) {
					deleted = append(deleted, row)
				} else {
					kept = append(kept, row)
				}
			}
			f.setTable(table, kept)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(deleted)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func matchRow(row map[string]any, q map[string][]string) bool {
	for key, vals := range q {
		if !strings.HasPrefix(vals[0], "eq.") {
			continue
		}
		want := strings.TrimPrefix(vals[0], "eq.")
		if fmt.Sprint(row[key]) != want {
			return false
		}
	}
	return true
}

func (f *fakePostgREST) table(name string) []map[string]any {
	if name == "aktivitas" {
		return f.aktivitas
	}
	return f.rekrutmen
}

func (f *fakePostgREST) setTable(name string, rows []map[string]any) {
	if name == "aktivitas" {
		f.aktivitas = rows
		return
	}
	f.rekrutmen = rows
}

func newTestStore(t *testing.T) (*Store, *fakePostgREST) {
	t.Helper()
	fake := newFakePostgREST()
	srv := fake.serve(t)
	t.Cleanup(srv.Close)
	s := New(srv.URL, "test-key", "")
	s.baseURL = srv.URL + "/rest/v1"
	return s, fake
}

func seedRow(id, nama, noWA string) *model.Rekrutmen {
	return &model.Rekrutmen{
		ID: id, Nama: nama, NoWA: noWA,
		Alamat: "Jl Melati 5", RT: "003", RW: "005",
		Kecamatan: "Serpong", Kelurahan: "Buaran",
		Aktivitas:         "Wiraswasta",
		TingkatPenerimaan: "Terbuka",
		Metode:            "Halaqoh",
		Status:            "Berhasil",
		UPL:               "UPLS",
	}
}

func TestSupabaseCreateAndGet(t *testing.T) {
	s, _ := newTestStore(t)
	r := seedRow("r1", "Ahmad Fauzi", "628123456789")
	if err := s.CreateRekrutmen(r); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := s.GetRekrutmen("r1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Nama != "Ahmad Fauzi" || got.NoWA != "628123456789" {
		t.Errorf("unexpected record: %+v", got)
	}
	if got.Kecamatan != "Serpong" || got.Kelurahan != "Buaran" {
		t.Errorf("wilayah mismatch: %+v", got)
	}
}

func TestSupabaseUniqueViolationMapsToConflict(t *testing.T) {
	s, fake := newTestStore(t)
	fake.failWith = "duplicate key value violates unique constraint"
	if err := s.CreateRekrutmen(seedRow("r1", "Ahmad", "628123456789")); err != store.ErrDuplicateWA {
		t.Fatalf("expected ErrDuplicateWA, got %v", err)
	}
}

func TestSupabaseUpdateAndDelete(t *testing.T) {
	s, _ := newTestStore(t)
	if err := s.CreateRekrutmen(seedRow("r1", "Ahmad", "628123456789")); err != nil {
		t.Fatal(err)
	}
	r := seedRow("r1", "Ahmad Update", "628123456789")
	r.Status = "Proses"
	if err := s.UpdateRekrutmen(r); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, err := s.GetRekrutmen("r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Nama != "Ahmad Update" || got.Status != "Proses" {
		t.Errorf("update not applied: %+v", got)
	}
	if err := s.DeleteRekrutmen("r1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.GetRekrutmen("r1"); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSupabaseFindByWhatsApp(t *testing.T) {
	s, _ := newTestStore(t)
	s.CreateRekrutmen(seedRow("r1", "Ahmad", "628123456789"))
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

func TestSupabaseAktivitasMaster(t *testing.T) {
	s, _ := newTestStore(t)
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

	// Renaming cascades to the rekrutmen rows.
	r := seedRow("r1", "Ahmad", "628123456789")
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

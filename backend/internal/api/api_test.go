package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/settings"
)

type testEnv struct {
	t      *testing.T
	server *Server
	client *http.Client
	base   string
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	dir := t.TempDir()
	users := auth.DefaultUsers()
	authMgr := auth.NewManager(users, []byte("test-secret"))
	setStore, err := settings.NewStore(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)

	jar, err := cookieJar()
	if err != nil {
		t.Fatal(err)
	}
	return &testEnv{t: t, server: srv, client: &http.Client{Jar: jar}, base: ts.URL}
}

// request performs an API call and decodes the JSON response.
func (e *testEnv) request(method, path string, body any, out any) int {
	e.t.Helper()
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			e.t.Fatal(err)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, e.base+path, rdr)
	if err != nil {
		e.t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := e.client.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer resp.Body.Close()
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			e.t.Fatalf("%s %s: decode: %v", method, path, err)
		}
	}
	return resp.StatusCode
}

func (e *testEnv) login(username, password string) int {
	return e.request(http.MethodPost, "/api/login", map[string]string{
		"username": username, "password": password,
	}, nil)
}

func validBody() map[string]any {
	return map[string]any{
		"nama":               "Ahmad Fauzi",
		"no_wa":              "628123456789",
		"alamat":             "Jl. Melati No. 5",
		"rt":                 "003",
		"rw":                 "005",
		"kecamatan":          "Serpong",
		"kelurahan":          "Buaran",
		"aktivitas":          "Wiraswasta",
		"tingkat_penerimaan": "Terbuka",
		"metode":             "Halaqoh",
		"status":             "Berhasil",
		"keterangan":         "test",
		"upl":                "UPLS",
	}
}

func TestLoginSetsSessionAndMeReturnsUser(t *testing.T) {
	e := newTestEnv(t)
	if code := e.login("admin", "admin123"); code != http.StatusOK {
		t.Fatalf("login: expected 200, got %d", code)
	}
	var out struct {
		User model.User `json:"user"`
	}
	if code := e.request(http.MethodGet, "/api/me", nil, &out); code != http.StatusOK {
		t.Fatalf("me: expected 200, got %d", code)
	}
	if out.User.Username != "admin" || out.User.Role != model.RoleAdmin {
		t.Errorf("unexpected user: %+v", out.User)
	}
}

func TestLoginRejectsBadPassword(t *testing.T) {
	e := newTestEnv(t)
	if code := e.login("admin", "wrong"); code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", code)
	}
}

func TestProtectedRoutesRequireAuth(t *testing.T) {
	e := newTestEnv(t)
	for _, path := range []string{"/api/rekrutmen", "/api/wilayah", "/api/aktivitas", "/api/settings"} {
		if code := e.request(http.MethodGet, path, nil, nil); code != http.StatusUnauthorized {
			t.Errorf("%s: expected 401, got %d", path, code)
		}
	}
}

func TestRekrutmenCRUDLifecycle(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")

	var created model.Rekrutmen
	if code := e.request(http.MethodPost, "/api/rekrutmen", validBody(), &created); code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", code)
	}
	if created.ID == "" {
		t.Fatal("expected an id on the created record")
	}
	if created.CreatedBy != "admin" {
		t.Errorf("expected created_by=admin, got %q", created.CreatedBy)
	}
	if created.NoWA != "628123456789" {
		t.Errorf("no_wa not stored as string: %q", created.NoWA)
	}

	// List.
	var list struct {
		Total int               `json:"total"`
		Data  []model.Rekrutmen `json:"data"`
	}
	if code := e.request(http.MethodGet, "/api/rekrutmen", nil, &list); code != http.StatusOK {
		t.Fatal("list failed")
	}
	if list.Total != 1 || len(list.Data) != 1 {
		t.Fatalf("expected 1 record, got total=%d len=%d", list.Total, len(list.Data))
	}

	// Get.
	var got model.Rekrutmen
	if code := e.request(http.MethodGet, "/api/rekrutmen/"+created.ID, nil, &got); code != http.StatusOK {
		t.Fatal("get failed")
	}
	if got.Nama != "Ahmad Fauzi" {
		t.Errorf("unexpected name %q", got.Nama)
	}

	// Update.
	body := validBody()
	body["nama"] = "Ahmad Fauzi Update"
	body["status"] = "Proses"
	if code := e.request(http.MethodPut, "/api/rekrutmen/"+created.ID, body, &got); code != http.StatusOK {
		t.Fatal("update failed")
	}
	if got.Nama != "Ahmad Fauzi Update" || got.Status != "Proses" {
		t.Errorf("update not applied: %+v", got)
	}

	// Delete.
	if code := e.request(http.MethodDelete, "/api/rekrutmen/"+created.ID, nil, nil); code != http.StatusOK {
		t.Fatal("delete failed")
	}
	if code := e.request(http.MethodGet, "/api/rekrutmen/"+created.ID, nil, nil); code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", code)
	}
}

// PRD section 9: the record must carry created_at/updated_at and the user
// who made the change. The request body carries only the editable fields, so
// an update must preserve the stored creation audit instead of zeroing it.
func TestUpdatePreservesCreationAuditFields(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")

	var created model.Rekrutmen
	if code := e.request(http.MethodPost, "/api/rekrutmen", validBody(), &created); code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", code)
	}
	if created.CreatedAt.IsZero() {
		t.Fatal("expected a non-zero created_at on the created record")
	}

	body := validBody()
	body["nama"] = "Ahmad Fauzi Update"
	var updated model.Rekrutmen
	if code := e.request(http.MethodPut, "/api/rekrutmen/"+created.ID, body, &updated); code != http.StatusOK {
		t.Fatal("update failed")
	}

	if updated.CreatedAt.IsZero() {
		t.Error("created_at was wiped by the update")
	}
	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("created_at changed: got %v, want %v", updated.CreatedAt, created.CreatedAt)
	}
	if updated.CreatedBy != "admin" {
		t.Errorf("created_by was wiped by the update: got %q", updated.CreatedBy)
	}
	if updated.UpdatedBy != "admin" {
		t.Errorf("expected updated_by=admin, got %q", updated.UpdatedBy)
	}
	if updated.UpdatedAt.Before(updated.CreatedAt) {
		t.Errorf("updated_at %v is earlier than created_at %v", updated.UpdatedAt, updated.CreatedAt)
	}

	// The persisted record must keep the audit trail too, not just the response.
	var got model.Rekrutmen
	if code := e.request(http.MethodGet, "/api/rekrutmen/"+created.ID, nil, &got); code != http.StatusOK {
		t.Fatal("get failed")
	}
	if got.CreatedBy != "admin" || got.CreatedAt.IsZero() {
		t.Errorf("stored record lost its creation audit: %+v", got)
	}
}

// PRD acceptance: "Nomor Whatsapp bersifat unik untuk setiap orang".
func TestDuplicateWhatsAppReturnsConflict(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")
	if code := e.request(http.MethodPost, "/api/rekrutmen", validBody(), nil); code != http.StatusCreated {
		t.Fatal("first create failed")
	}
	if code := e.request(http.MethodPost, "/api/rekrutmen", validBody(), nil); code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate, got %d", code)
	}
}

// PRD: "tampilkan peringatan sebelum menyimpan nomor yang sudah ada".
func TestCheckWhatsAppReportsExistence(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")
	e.request(http.MethodPost, "/api/rekrutmen", validBody(), nil)

	var check struct {
		Exists bool   `json:"exists"`
		ID     string `json:"id"`
	}
	if code := e.request(http.MethodGet, "/api/rekrutmen/check-wa?no=628123456789", nil, &check); code != http.StatusOK {
		t.Fatal("check-wa failed")
	}
	if !check.Exists || check.ID == "" {
		t.Errorf("expected exists=true with an id, got %+v", check)
	}

	var missing struct {
		Exists bool `json:"exists"`
	}
	e.request(http.MethodGet, "/api/rekrutmen/check-wa?no=628999999999", nil, &missing)
	if missing.Exists {
		t.Error("expected exists=false for an unknown number")
	}
}

// PRD: "Nomor Whatsapp ... menolak karakter selain digit".
func TestValidationErrorsAreReportedPerField(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")

	body := validBody()
	body["no_wa"] = "628-abc"
	body["rt"] = "00A"
	body["kecamatan"] = ""

	var out struct {
		Error  string            `json:"error"`
		Fields map[string]string `json:"fields"`
	}
	if code := e.request(http.MethodPost, "/api/rekrutmen", body, &out); code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", code)
	}
	for _, field := range []string{"no_wa", "rt", "kecamatan"} {
		if _, ok := out.Fields[field]; !ok {
			t.Errorf("expected an error on %q, got %v", field, out.Fields)
		}
	}
}

// PRD: "Kecamatan wajib dan hanya menampilkan 7 kecamatan Tangsel".
func TestWilayahEndpointReturnsSevenKecamatan(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")

	var out struct {
		Kecamatan []struct {
			ID        string `json:"id"`
			Kelurahan []struct {
				ID string `json:"id"`
			} `json:"kelurahan"`
		} `json:"kecamatan"`
	}
	if code := e.request(http.MethodGet, "/api/wilayah", nil, &out); code != http.StatusOK {
		t.Fatal("wilayah failed")
	}
	if len(out.Kecamatan) != 7 {
		t.Fatalf("expected 7 kecamatan, got %d", len(out.Kecamatan))
	}
	total := 0
	for _, k := range out.Kecamatan {
		total += len(k.Kelurahan)
	}
	if total != 54 {
		t.Errorf("expected 54 kelurahan, got %d", total)
	}
}

// PRD: master aktivitas dapat dikelola Admin.
func TestAktivitasMasterIsAdminOnly(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")

	var list []model.Aktivitas
	if code := e.request(http.MethodGet, "/api/aktivitas", nil, &list); code != http.StatusOK {
		t.Fatal("list failed")
	}
	if len(list) == 0 {
		t.Fatal("expected the seeded activity master to be non-empty")
	}

	var created model.Aktivitas
	if code := e.request(http.MethodPost, "/api/aktivitas", map[string]string{"nama": "Nelayan"}, &created); code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", code)
	}
	created.Nama = "Nelayanbigril"
	if code := e.request(http.MethodPut, "/api/aktivitas/"+created.ID, created, &created); code != http.StatusOK {
		t.Fatal("update failed")
	}
	if code := e.request(http.MethodDelete, "/api/aktivitas/"+created.ID, nil, nil); code != http.StatusOK {
		t.Fatal("delete failed")
	}
}

func TestPetugasCannotChangeMasterAktivitas(t *testing.T) {
	e := newTestEnv(t)
	e.login("petugas", "petugas123")
	if code := e.request(http.MethodPost, "/api/aktivitas", map[string]string{"nama": "X"}, nil); code != http.StatusForbidden {
		t.Errorf("expected 403 for petugas, got %d", code)
	}
}

// Petugas may still add recruitment data.
func TestPetugasCanCreateRekrutmen(t *testing.T) {
	e := newTestEnv(t)
	e.login("petugas", "petugas123")
	if code := e.request(http.MethodPost, "/api/rekrutmen", validBody(), nil); code != http.StatusCreated {
		t.Errorf("expected 201, got %d", code)
	}
}

func TestSettingsRequireAdmin(t *testing.T) {
	e := newTestEnv(t)
	e.login("petugas", "petugas123")
	if code := e.request(http.MethodGet, "/api/settings", nil, nil); code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", code)
	}
}

func TestSettingsSwitchProvider(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")

	var cur struct {
		Provider string `json:"provider"`
	}
	if code := e.request(http.MethodGet, "/api/settings", nil, &cur); code != http.StatusOK {
		t.Fatal("settings failed")
	}
	if cur.Provider != settings.ProviderLocal {
		t.Errorf("expected the local provider by default, got %q", cur.Provider)
	}

	// Switching to an unconfigured provider must fail without changing state.
	if code := e.request(http.MethodPut, "/api/settings", map[string]string{
		"provider": settings.ProviderSupa,
	}, nil); code != http.StatusBadRequest {
		t.Errorf("expected 400 for an unconfigured provider, got %d", code)
	}
	e.request(http.MethodGet, "/api/settings", nil, &cur)
	if cur.Provider != settings.ProviderLocal {
		t.Errorf("provider changed despite a failed switch: %q", cur.Provider)
	}
}

func TestSearchAndFilter(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")

	a := validBody()
	a["nama"] = "Budi Santoso"
	e.request(http.MethodPost, "/api/rekrutmen", a, nil)

	b := validBody()
	b["nama"] = "Siti Aminah"
	b["no_wa"] = "628999999999"
	b["kecamatan"] = "Ciputat"
	b["kelurahan"] = "Cipayung"
	b["status"] = "Pending"
	e.request(http.MethodPost, "/api/rekrutmen", b, nil)

	var out struct {
		Total int `json:"total"`
	}
	e.request(http.MethodGet, "/api/rekrutmen?q=siti", nil, &out)
	if out.Total != 1 {
		t.Errorf("search by name: expected 1, got %d", out.Total)
	}
	e.request(http.MethodGet, "/api/rekrutmen?q=628999", nil, &out)
	if out.Total != 1 {
		t.Errorf("search by wa: expected 1, got %d", out.Total)
	}
	e.request(http.MethodGet, "/api/rekrutmen?kecamatan=Ciputat", nil, &out)
	if out.Total != 1 {
		t.Errorf("filter by kecamatan: expected 1, got %d", out.Total)
	}
	e.request(http.MethodGet, "/api/rekrutmen?status=Pending", nil, &out)
	if out.Total != 1 {
		t.Errorf("filter by status: expected 1, got %d", out.Total)
	}
	e.request(http.MethodGet, "/api/rekrutmen", nil, &out)
	if out.Total != 2 {
		t.Errorf("expected 2 total, got %d", out.Total)
	}
}

// The Metode list has three options; a record submitted with a retired
// spelling must be stored as "Program Antara" and be findable by the
// current filter value.
func TestMetodeProgramAntaraIsStoredAndFilterable(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")

	body := validBody()
	body["metode"] = "Program"
	var created model.Rekrutmen
	if code := e.request(http.MethodPost, "/api/rekrutmen", body, &created); code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", code)
	}
	if created.Metode != "Program Antara" {
		t.Errorf("expected metode normalized to %q, got %q", "Program Antara", created.Metode)
	}

	// The retired spelling must also be accepted by the filter.
	for _, filter := range []string{"Program Antara", "Program"} {
		var out struct {
			Total int `json:"total"`
		}
		e.request(http.MethodGet, "/api/rekrutmen?metode="+url.QueryEscape(filter), nil, &out)
		if out.Total != 1 {
			t.Errorf("filter metode=%q: expected 1, got %d", filter, out.Total)
		}
	}
}

// An activity referenced by a record must not be deletable.
func TestAktivitasInUseCannotBeDeleted(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")

	var list []model.Aktivitas
	e.request(http.MethodGet, "/api/aktivitas", nil, &list)
	if len(list) == 0 {
		t.Fatal("no seeded activities")
	}
	used := list[0].Nama

	body := validBody()
	body["aktivitas"] = used
	e.request(http.MethodPost, "/api/rekrutmen", body, nil)

	for _, a := range list {
		if a.Nama == used {
			if code := e.request(http.MethodDelete, "/api/aktivitas/"+a.ID, nil, nil); code != http.StatusConflict {
				t.Errorf("expected 409 for an activity in use, got %d", code)
			}
		}
	}
}

func TestLogoutClearsSession(t *testing.T) {
	e := newTestEnv(t)
	e.login("admin", "admin123")
	if code := e.request(http.MethodPost, "/api/logout", nil, nil); code != http.StatusOK {
		t.Fatal("logout failed")
	}
	if code := e.request(http.MethodGet, "/api/me", nil, nil); code != http.StatusUnauthorized {
		t.Errorf("expected 401 after logout, got %d", code)
	}
}

// legacyAuthenticator stands in for the accounts of the previous
// application.
type legacyAuthenticator struct {
	user  model.User
	ok    bool
	calls int
}

func (l *legacyAuthenticator) VerifyLegacy(username, password string) (model.User, bool, error) {
	l.calls++
	return l.user, l.ok, nil
}

func TestLoginWithLegacyAccount(t *testing.T) {
	dir := t.TempDir()
	users := auth.DefaultUsers()
	authMgr := auth.NewManager(users, []byte("test-secret"))
	authMgr.SetUsersPath(filepath.Join(dir, "users.json"))
	legacy := &legacyAuthenticator{
		ok:   true,
		user: model.User{ID: "u1", Username: "anggota01", Nama: "Anggota Satu", Role: model.RolePetugas},
	}
	authMgr.SetLegacyAuthenticator(legacy)

	setStore, err := settings.NewStore(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)

	jar, err := cookieJar()
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Jar: jar}

	body, _ := json.Marshal(map[string]string{"username": "anggota01", "password": "rahasia"})
	resp, err := client.Post(ts.URL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login warisan harus berhasil, dapat %d", resp.StatusCode)
	}
	var out struct {
		User     model.User `json:"user"`
		Migrated bool       `json:"migrated"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.User.Username != "anggota01" {
		t.Errorf("user: %+v", out.User)
	}
	if !out.Migrated {
		t.Error("login warisan yang berhasil harus ditandai sudah dimigrasi")
	}

	// The account is stored locally now, so a second login does not need the
	// sheet any more.
	callsBefore := legacy.calls
	body, _ = json.Marshal(map[string]string{"username": "anggota01", "password": "rahasia"})
	resp2, err := client.Post(ts.URL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("login kedua harus berhasil, dapat %d", resp2.StatusCode)
	}
	if legacy.calls != callsBefore {
		t.Error("akun yang sudah dimigrasi tidak boleh menyentuh sheet lagi")
	}
}

func TestLoginWithWrongLegacyPassword(t *testing.T) {
	dir := t.TempDir()
	users := auth.DefaultUsers()
	authMgr := auth.NewManager(users, []byte("test-secret"))
	authMgr.SetUsersPath(filepath.Join(dir, "users.json"))
	authMgr.SetLegacyAuthenticator(&legacyAuthenticator{ok: false})

	setStore, err := settings.NewStore(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)

	body, _ := json.Marshal(map[string]string{"username": "anggota01", "password": "salah"})
	resp, err := http.Post(ts.URL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("password salah harus ditolak, dapat %d", resp.StatusCode)
	}
}

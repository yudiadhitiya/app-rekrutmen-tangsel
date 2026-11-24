package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/settings"
)

// districtUser builds an account with a given level and district.
func districtUser(t *testing.T, username, password string, level int, kecamatan string) auth.User {
	t.Helper()
	salt, hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	return auth.User{
		User: model.User{
			ID:        username,
			Username:  username,
			Role:      model.RoleAdmin,
			Level:     level,
			Kecamatan: kecamatan,
		},
		PasswordHash: hash,
		Salt:         salt,
	}
}

// newScopeEnv builds a server with one district administrator and one member.
func newScopeEnv(t *testing.T, allowRegistration bool) (*Server, string, string) {
	t.Helper()
	dir := t.TempDir()
	users := []auth.User{
		districtUser(t, "pwkec", "kecamatan-secret-2024", model.LevelAdminKecamatan, "Serpong"),
		districtUser(t, "petugas", "petugas-secret-2024", model.LevelAnggota, ""),
	}
	authMgr := auth.NewManager(users, []byte("test-secret"))
	authMgr.SetRegistrar(&fakeRegistrar{})
	setStore, err := settings.NewStore(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir, AllowSelfRegistration: allowRegistration}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}
	return srv, loginCookie(t, srv, "pwkec", "kecamatan-secret-2024"), loginCookie(t, srv, "petugas", "petugas-secret-2024")
}

// loginCookie signs in and returns the session cookie.
func loginCookie(t *testing.T, srv *Server, username, password string) string {
	t.Helper()
	rec := httptest.NewRecorder()
	body := `{"username":"` + username + `","password":"` + password + `"}`
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("login %s gagal: %d %s", username, rec.Code, rec.Body.String())
	}
	return rec.Header().Get("Set-Cookie")
}

// TestDistrictAdminCannotMoveRecordToAnotherDistrict checks that a district
// administrator cannot move a record out of their own district.
func TestDistrictAdminCannotMoveRecordToAnotherDistrict(t *testing.T) {
	srv, kecCookie, _ := newScopeEnv(t, false)

	if err := srv.Store().CreateRekrutmen(&model.Rekrutmen{
		ID: "r1", Nama: "Data Asli", NoWA: "628999999999", Kecamatan: "Serpong",
		Aktivitas: "Wiraswasta", TingkatPenerimaan: "Terbuka", Metode: "Halaqoh",
		Status: "Berhasil", CreatedBy: "admin",
	}); err != nil {
		t.Fatal(err)
	}

	body := `{"nama":"Data Asli","no_wa":"628999999999","kecamatan":"Ciputat","aktivitas":"Wiraswasta","tingkat_penerimaan":"Terbuka","metode":"Halaqoh","status":"Berhasil"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/rekrutmen/r1", strings.NewReader(body))
	req.Header.Set("Cookie", kecCookie)
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("admin kec harus ditolak, kode %d: %s", rec.Code, rec.Body.String())
	}

	stored, err := srv.Store().GetRekrutmen("r1")
	if err != nil {
		t.Fatal(err)
	}
	if stored.Kecamatan != "Serpong" {
		t.Errorf("kecamatan record berubah jadi %q", stored.Kecamatan)
	}
}

// TestDistrictAdminCannotCreateRecordInAnotherDistrict checks that a
// district administrator cannot put a new record into another district.
func TestDistrictAdminCannotCreateRecordInAnotherDistrict(t *testing.T) {
	srv, kecCookie, _ := newScopeEnv(t, false)

	body := `{"nama":"Data Asli","no_wa":"628888888888","kecamatan":"Ciputat","aktivitas":"Wiraswasta","tingkat_penerimaan":"Terbuka","metode":"Halaqoh","status":"Berhasil"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/rekrutmen", strings.NewReader(body))
	req.Header.Set("Cookie", kecCookie)
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("admin kec harus ditolak, kode %d: %s", rec.Code, rec.Body.String())
	}
}

// TestMemberCannotMoveOwnRecordToAnotherDistrict checks that a member
// cannot move its own record into the queue of another district.
func TestMemberCannotMoveOwnRecordToAnotherDistrict(t *testing.T) {
	srv, _, petugasCookie := newScopeEnv(t, false)

	if err := srv.Store().CreateRekrutmen(&model.Rekrutmen{
		ID: "r2", Nama: "Data Sendiri", NoWA: "628777777777", Kecamatan: "Serpong",
		Aktivitas: "Wiraswasta", TingkatPenerimaan: "Terbuka", Metode: "Halaqoh",
		Status: "Berhasil", CreatedBy: "petugas",
	}); err != nil {
		t.Fatal(err)
	}

	body := `{"nama":"Data Sendiri","no_wa":"628777777777","kecamatan":"Ciputat","aktivitas":"Wiraswasta","tingkat_penerimaan":"Terbuka","metode":"Halaqoh","status":"Berhasil"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/rekrutmen/r2", strings.NewReader(body))
	req.Header.Set("Cookie", petugasCookie)
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("anggota harus ditolak, kode %d: %s", rec.Code, rec.Body.String())
	}
}

// TestSelfRegistrationIsOffByDefault checks that the public registration
// endpoint refuses strangers unless the operator turns it on.
func TestSelfRegistrationIsOffByDefault(t *testing.T) {
	srv, _, _ := newScopeEnv(t, false)

	body := `{"username":"asing2024","password":"KuatBanget-9924x","nama":"Orang Asing","kecamatan":"Serpong","question1":"` +
		model.SecurityQuestions[0] + `","answer1":"budi","question2":"` + model.SecurityQuestions[1] + `","answer2":"jakarta"}`
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(body)))
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("pendaftaran harus mati secara bawaan, kode %d: %s", rec.Code, rec.Body.String())
	}

	// The options endpoint tells the form that registration is closed.
	rec2 := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec2, httptest.NewRequest(http.MethodGet, "/api/register/options", nil))
	if !strings.Contains(rec2.Body.String(), `"available":false`) {
		t.Errorf("opsi pendaftaran harus melaporkan available=false: %s", rec2.Body.String())
	}
}

// TestSelfRegistrationWorksWhenEnabled checks that the operator can still
// open the endpoint on purpose.
func TestSelfRegistrationWorksWhenEnabled(t *testing.T) {
	srv, _, _ := newScopeEnv(t, true)

	body := `{"username":"asing2024","password":"KuatBanget-9924x","nama":"Orang Asing","kecamatan":"Serpong","question1":"` +
		model.SecurityQuestions[0] + `","answer1":"budi","question2":"` + model.SecurityQuestions[1] + `","answer2":"jakarta"}`
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Errorf("pendaftaran harus jalan saat diizinkan, kode %d: %s", rec.Code, rec.Body.String())
	}
}

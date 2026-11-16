package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/settings"
)

// fakeLegacyStore stands in for the accounts of the previous application.
//
// It accepts one password and can store a new one, so the effect of a
// password change on that store can be checked.
type fakeLegacyStore struct {
	password string
	found    bool
	changed  bool
}

func (f *fakeLegacyStore) VerifyLegacy(username, password string) (model.User, bool, error) {
	if !f.found || password != f.password {
		return model.User{}, false, nil
	}
	return model.User{
		ID:       "legacy-" + username,
		Username: username,
		Nama:     username,
		Role:     model.RolePetugas,
		Level:    model.LevelAnggota,
	}, true, nil
}

func (f *fakeLegacyStore) SetLegacyPassword(username, password string) (bool, error) {
	if !f.found {
		return false, nil
	}
	f.password = password
	f.changed = true
	return true, nil
}

// newTestEnvWithLegacy starts an environment whose accounts of the previous
// application are known to the given store.
func newTestEnvWithLegacy(t *testing.T, legacy *fakeLegacyStore) *testEnv {
	t.Helper()
	dir := t.TempDir()
	authMgr := auth.NewManager(auth.DefaultUsers(), []byte("test-secret"))
	authMgr.SetLegacyAuthenticator(legacy)
	authMgr.SetUsersPath(filepath.Join(dir, "users.json"))
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

// changePassword calls the endpoint with the given passwords.
func (e *testEnv) changePassword(current, next string) int {
	return e.request(http.MethodPut, "/api/me/password", map[string]string{
		"current_password": current,
		"new_password":     next,
	}, nil)
}

// TestChangePasswordNeedsSession checks that the endpoint is protected.
func TestChangePasswordNeedsSession(t *testing.T) {
	env := newTestEnv(t)

	if code := env.changePassword("petugas123", "petugas456"); code != http.StatusUnauthorized {
		t.Fatalf("tanpa sesi harus ditolak: %d", code)
	}
}

// TestChangePasswordOfMemberAccount checks that an account without any access
// level may change its own password.
func TestChangePasswordOfMemberAccount(t *testing.T) {
	env := newTestEnv(t)
	if code := env.login("petugas", "petugas123"); code != http.StatusOK {
		t.Fatalf("login: %d", code)
	}

	if code := env.changePassword("petugas123", "petugas456"); code != http.StatusOK {
		t.Fatalf("anggota harus boleh mengganti password: %d", code)
	}
}

// TestChangePasswordOfEveryLevel checks that the access level does not decide
// whether an account may change its own password.
func TestChangePasswordOfEveryLevel(t *testing.T) {
	env := newTestEnvWithUsers(t, []auth.User{
		newTestUser(t, "kec", "kec12345", model.LevelAdminKecamatan, "Serpong"),
	})

	cases := []struct {
		username string
		password string
	}{
		{"petugas", "petugas123"}, // level 1 - anggota
		{"kec", "kec12345"},       // level 2 - admin kecamatan
		{"admin", "admin123"},     // level 3 - admin penuh
	}
	for _, c := range cases {
		if code := env.login(c.username, c.password); code != http.StatusOK {
			t.Fatalf("login %s: %d", c.username, code)
		}
		if code := env.changePassword(c.password, "baru12345"); code != http.StatusOK {
			t.Fatalf("level %s harus boleh mengganti password: %d", c.username, code)
		}
	}
}

// TestChangePasswordRejectsWrongCurrentPassword checks that the current
// password is needed, so a borrowed session cannot take the account over.
func TestChangePasswordRejectsWrongCurrentPassword(t *testing.T) {
	env := newTestEnv(t)
	if code := env.login("petugas", "petugas123"); code != http.StatusOK {
		t.Fatalf("login: %d", code)
	}

	var out map[string]any
	code := env.request(http.MethodPut, "/api/me/password", map[string]string{
		"current_password": "salah",
		"new_password":     "petugas456",
	}, &out)
	if code != http.StatusBadRequest {
		t.Fatalf("password lama salah harus ditolak: %d", code)
	}
	if out["error"] != "password lama salah" {
		t.Errorf("pesan: %v", out["error"])
	}
	if code := env.login("petugas", "petugas123"); code != http.StatusOK {
		t.Errorf("password lama harus tetap berlaku: %d", code)
	}
}

// TestChangePasswordRejectsShortPassword checks the length rule.
func TestChangePasswordRejectsShortPassword(t *testing.T) {
	env := newTestEnv(t)
	if code := env.login("petugas", "petugas123"); code != http.StatusOK {
		t.Fatalf("login: %d", code)
	}

	var out map[string]any
	code := env.request(http.MethodPut, "/api/me/password", map[string]string{
		"current_password": "petugas123",
		"new_password":     "abc",
	}, &out)
	if code != http.StatusBadRequest {
		t.Fatalf("password pendek harus ditolak: %d", code)
	}
	fields, _ := out["fields"].(map[string]any)
	if fields == nil || fields["password"] == nil {
		t.Errorf("pesan validasi per field tidak ada: %v", out)
	}
}

// TestChangePasswordRejectsSamePassword checks that the new password differs
// from the current one.
func TestChangePasswordRejectsSamePassword(t *testing.T) {
	env := newTestEnv(t)
	if code := env.login("petugas", "petugas123"); code != http.StatusOK {
		t.Fatalf("login: %d", code)
	}

	if code := env.changePassword("petugas123", "petugas123"); code != http.StatusBadRequest {
		t.Fatalf("password yang sama harus ditolak: %d", code)
	}
}

// TestChangePasswordNeedsCurrentPassword checks that the current password is
// mandatory.
func TestChangePasswordNeedsCurrentPassword(t *testing.T) {
	env := newTestEnv(t)
	if code := env.login("petugas", "petugas123"); code != http.StatusOK {
		t.Fatalf("login: %d", code)
	}

	if code := env.changePassword("", "petugas456"); code != http.StatusBadRequest {
		t.Fatalf("password lama kosong harus ditolak: %d", code)
	}
}

// TestChangePasswordRejectsBrokenBody checks the body of the request.
func TestChangePasswordRejectsBrokenBody(t *testing.T) {
	env := newTestEnv(t)
	if code := env.login("petugas", "petugas123"); code != http.StatusOK {
		t.Fatalf("login: %d", code)
	}

	code := env.request(http.MethodPut, "/api/me/password", "bukan-json", nil)
	if code != http.StatusBadRequest {
		t.Fatalf("body rusak harus ditolak: %d", code)
	}
}

// TestChangePasswordKeepsTheSession checks that the session of the account
// stays usable after the password was changed.
func TestChangePasswordKeepsTheSession(t *testing.T) {
	env := newTestEnv(t)
	if code := env.login("petugas", "petugas123"); code != http.StatusOK {
		t.Fatalf("login: %d", code)
	}
	if code := env.changePassword("petugas123", "petugas456"); code != http.StatusOK {
		t.Fatalf("ubah password: %d", code)
	}

	var out map[string]any
	if code := env.request(http.MethodGet, "/api/me", nil, &out); code != http.StatusOK {
		t.Fatalf("sesi harus tetap berlaku: %d", code)
	}
}

// TestChangePasswordIsUsedByTheNextLogin checks that the new password works
// and the old one does not.
func TestChangePasswordIsUsedByTheNextLogin(t *testing.T) {
	env := newTestEnv(t)
	if code := env.login("petugas", "petugas123"); code != http.StatusOK {
		t.Fatalf("login: %d", code)
	}
	if code := env.changePassword("petugas123", "petugas456"); code != http.StatusOK {
		t.Fatalf("ubah password: %d", code)
	}

	if code := env.login("petugas", "petugas456"); code != http.StatusOK {
		t.Errorf("login dengan password baru: %d", code)
	}
	if code := env.login("petugas", "petugas123"); code != http.StatusUnauthorized {
		t.Errorf("password lama masih diterima: %d", code)
	}
}

// TestChangePasswordUpdatesTheAccountOfThePreviousApplication checks that an
// account of the previous application keeps the same password.
func TestChangePasswordUpdatesTheAccountOfThePreviousApplication(t *testing.T) {
	legacy := &fakeLegacyStore{password: "lama12345", found: true}
	env := newTestEnvWithLegacy(t, legacy)

	if code := env.login("anggota01", "lama12345"); code != http.StatusOK {
		t.Fatalf("login: %d", code)
	}
	if code := env.changePassword("lama12345", "baru12345"); code != http.StatusOK {
		t.Fatalf("ubah password: %d", code)
	}
	if !legacy.changed || legacy.password != "baru12345" {
		t.Errorf("password warisan: %q (ditulis: %v)", legacy.password, legacy.changed)
	}
}

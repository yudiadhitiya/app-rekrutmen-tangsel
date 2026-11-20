package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

// TestAPIResponsesAreNotStoredInACache checks that the answers of the API are
// not kept in a cache of the browser or of a proxy.
//
// The answers hold the data of the people who registered, so they must not
// stay anywhere after the page is closed.
func TestAPIResponsesAreNotStoredInACache(t *testing.T) {
	h := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/rekrutmen", nil))
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, ingin no-store", got)
	}

	// The static files of the frontend may be cached, so the header is not
	// set for them.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/assets/app.js", nil))
	if got := rec.Header().Get("Cache-Control"); got != "" {
		t.Errorf("aset statis tidak perlu no-store, dapat %q", got)
	}
}

// TestLoginLimiterSurvivesARestart checks that the failed sign ins are still
// counted after the server is started again.
//
// Otherwise a restart would hand a fresh set of attempts to whoever tries to
// guess a password.
func TestLoginLimiterSurvivesARestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "limiter_login.json")

	first := newLoginLimiter(3, newLimiterStore(path))
	for i := 0; i < 3; i++ {
		first.fail("admin", "192.0.2.1")
	}
	if first.allow("admin", "192.0.2.1") {
		t.Fatal("limiter harus menolak setelah tiga kegagalan")
	}

	// The server starts again and reads the counters from the disk.
	second := newLoginLimiter(3, newLimiterStore(path))
	if second.allow("admin", "192.0.2.1") {
		t.Error("penghitung harus bertahan setelah server dimulai ulang")
	}

	// A sign in that succeeds clears the counter of the account.
	second.reset("admin", "192.0.2.1")
	third := newLoginLimiter(3, newLimiterStore(path))
	if !third.allow("admin", "192.0.2.1") {
		t.Error("penghitung akun harus bersih setelah masuk berhasil")
	}
}

// TestRegisterRefusesAWeakPassword checks that a registration with a password
// that everybody knows is refused.
func TestRegisterRefusesAWeakPassword(t *testing.T) {
	registrar := &fakeRegistrar{}
	env := newRegisterEnv(t, registrar)

	body := validRegistration()
	body["password"] = "password12345"
	var out map[string]any
	if code := env.request(http.MethodPost, "/api/register", body, &out); code != http.StatusBadRequest {
		t.Errorf("password umum harus ditolak, kode %d", code)
	}
	if len(registrar.accounts) != 0 {
		t.Error("akun tidak boleh dibuat dari password yang ditolak")
	}
}

// TestRegisterRefusesTheUsernameAsPassword checks that the name of the
// account cannot be its own password.
func TestRegisterRefusesTheUsernameAsPassword(t *testing.T) {
	registrar := &fakeRegistrar{}
	env := newRegisterEnv(t, registrar)

	body := validRegistration()
	body["password"] = body["username"]
	var out map[string]any
	if code := env.request(http.MethodPost, "/api/register", body, &out); code != http.StatusBadRequest {
		t.Errorf("password yang sama dengan username harus ditolak, kode %d", code)
	}
}

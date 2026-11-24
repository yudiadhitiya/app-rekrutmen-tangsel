package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/authtest"
	"rekrutmen-tangsel/backend/internal/settings"
)

// newLoginTestEnv starts an environment with a given number of allowed failed
// sign ins.
func newLoginTestEnv(t *testing.T, maxAttempts int) *testEnv {
	t.Helper()
	dir := t.TempDir()
	authMgr := auth.NewManager(authtest.Users(), []byte("test-secret"))
	setStore, err := settings.NewStore(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir, LoginMaxAttempts: maxAttempts}, authMgr, setStore)
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

// TestLockedAccountRejectsWrongPasswordsButLetsRightPasswordIn checks
// that an account under brute-force attack still accepts the correct
// password, so an attacker cannot lock the account out.
func TestLockedAccountRejectsWrongPasswordsButLetsRightPasswordIn(t *testing.T) {
	e := newLoginTestEnv(t, 3)

	// The account limit is 2 * max = 6 failures in 3 hours.
	// The attacker sends 6 wrong passwords from one address.
	for i := 0; i < 6; i++ {
		code := e.loginFrom("admin", "salah", "198.51.100.1:5555")
		if code != http.StatusUnauthorized {
			t.Fatalf("percobaan %d: status %d, ingin 401", i+1, code)
		}
	}
	// The 7th wrong attempt trips the account counter and is
	// refused with 429.
	code := e.loginFrom("admin", "salah", "198.51.100.1:5555")
	if code != http.StatusTooManyRequests {
		t.Fatalf("percobaan ke-7 (terkunci): status %d, ingin 429", code)
	}

	// The correct password from a different address (the admin's
	// own office) still works even while the account is locked.
	// This is the core defense against lockout abuse.
	code = e.loginFrom("admin", "admin123", "203.0.113.5:5555")
	if code != http.StatusOK {
		t.Errorf("password yang benar harus diterima dari alamat baru, kode %d", code)
	}
}

// TestLockedAccountRejectsWrongPasswordsFromDifferentAddresses checks
// that a distributed brute force against an account is also stopped
// by the account counter.
func TestLockedAccountRejectsWrongPasswordsFromDifferentAddresses(t *testing.T) {
	e := newLoginTestEnv(t, 3)

	// Fail from 6 different addresses. The account limit is 6,
	// so the 7th failure trips the counter.
	for i := 0; i < 6; i++ {
		code := e.loginFrom("admin", "salah", "10.0.0."+string(rune('1'+i))+":5555")
		if code != http.StatusUnauthorized {
			t.Fatalf("percobaan %d dari 10.0.0.%d: status %d, ingin 401", i+1, i+1, code)
		}
	}
	// The 7th failure from yet another address trips the counter.
	code := e.loginFrom("admin", "salah", "10.0.0.99:5555")
	if code != http.StatusTooManyRequests {
		t.Fatalf("percobaan ke-7 (terkunci): status %d, ingin 429", code)
	}

	// The correct password from yet another address still works.
	code = e.loginFrom("admin", "admin123", "203.0.113.1:5555")
	if code != http.StatusOK {
		t.Errorf("password yang benar harus diterima dari alamat baru, kode %d", code)
	}
}

// TestLoginUnlockedAfterSuccess checks that a successful sign in clears the
// counter, so a user who mistyped the password is not locked out.
func TestLoginUnlockedAfterSuccess(t *testing.T) {
	e := newLoginTestEnv(t, 3)

	for i := 0; i < 2; i++ {
		e.login("admin", "salah")
	}
	if code := e.login("admin", "admin123"); code != http.StatusOK {
		t.Fatalf("login yang benar ditolak: %d", code)
	}
	for i := 0; i < 2; i++ {
		if code := e.login("admin", "salah"); code != http.StatusUnauthorized {
			t.Fatalf("percobaan %d setelah login berhasil: %d", i+1, code)
		}
	}
	if code := e.login("admin", "admin123"); code != http.StatusOK {
		t.Errorf("login yang benar ditolak setelah percobaan gagal: %d", code)
	}
}

// TestLoginLockIsPerAccount checks that the lock of one account does not lock
// another account.
func TestLoginLockIsPerAccount(t *testing.T) {
	e := newLoginTestEnv(t, 3)

	for i := 0; i < 3; i++ {
		e.login("admin", "salah")
	}
	if code := e.login("petugas", "petugas123"); code != http.StatusOK {
		t.Errorf("akun lain ikut terkunci: %d", code)
	}
}

// TestLoginLockPerAddress checks that one address cannot attack many
// accounts, so a single address cannot guess the password of anybody.
//
// The limit of an address is a multiple of the limit of an account, because
// an address is shared by many people.
func TestLoginLockPerAddress(t *testing.T) {
	e := newLoginTestEnv(t, 1)

	names := []string{"admin", "petugas", "petugas2", "petugas3", "petugas4", "petugas5"}
	locked := false
	for _, name := range names {
		code := e.login(name, "salah")
		if code == http.StatusTooManyRequests {
			locked = true
			break
		}
	}
	if !locked {
		t.Error("percobaan dari satu alamat tidak pernah dikunci")
	}
}

// mustParse reads the address of a test server.
func mustParse(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

// TestPasswordChangeInvalidatesOldSession checks that a session that was made
// before a password change stops working.
func TestPasswordChangeInvalidatesOldSession(t *testing.T) {
	e := newTestEnv(t)

	if code := e.login("admin", "admin123"); code != http.StatusOK {
		t.Fatalf("login gagal: %d", code)
	}
	// The cookie of the session before the change is kept aside.
	var oldCookie *http.Cookie
	for _, c := range e.client.Jar.Cookies(mustParse(t, e.base)) {
		if c.Name == "rekr_session" {
			oldCookie = c
		}
	}
	if oldCookie == nil {
		t.Fatal("cookie sesi tidak ditemukan")
	}

	if code := e.request(http.MethodPut, "/api/me/password", map[string]string{
		"current_password": "admin123",
		"new_password":     "admin-sekarang-2024",
	}, nil); code != http.StatusOK {
		t.Fatalf("ganti password gagal: %d", code)
	}

	// The session of the browser that changed the password keeps working.
	if code := e.request(http.MethodGet, "/api/me", nil, nil); code != http.StatusOK {
		t.Errorf("sesi setelah ganti password tidak berlaku: %d", code)
	}

	// The old session is refused.
	req, err := http.NewRequest(http.MethodGet, e.base+"/api/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(oldCookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("sesi lama masih berlaku: %d, ingin 401", resp.StatusCode)
	}
}

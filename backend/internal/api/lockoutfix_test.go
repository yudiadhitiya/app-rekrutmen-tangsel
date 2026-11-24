package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// loginAs sends a sign in to the handler from the given address and returns
// the status code.
//
// The request is handed to the handler directly, so the address of the
// caller can be chosen. A real client would always come from the address of
// the test machine.
func loginAs(h http.Handler, username, password, address string) int {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login",
		strings.NewReader(`{"username":"`+username+`","password":"`+password+`"}`))
	req.RemoteAddr = address
	h.ServeHTTP(rec, req)
	return rec.Code
}

// TestAttackerCannotLockOutTheAdmin checks the finding that an attacker
// could lock the account of the admin for 15 minutes, over and over again.
//
// The counter of the account refuses the wrong password, but a password that
// is right always gets in. So the attacker cannot keep the admin out of the
// application.
func TestAttackerCannotLockOutTheAdmin(t *testing.T) {
	e := newLoginTestEnv(t, 0)
	h := e.server.Handler()

	// The attacker guesses the password of the admin. The limit of the
	// account is 2 * 10 = 20 failures in 3 hours.
	for i := 0; i < 30; i++ {
		loginAs(h, "admin", "salah", "198.51.100.9:5555")
	}

	// The admin signs in from his own address with the right password.
	// Before the fix this answer was 429 and the application was
	// unusable until the window had passed.
	code := loginAs(h, "admin", "admin123", "203.0.113.9:5555")
	if code != http.StatusOK {
		t.Errorf("admin harus tetap bisa masuk, kode %d", code)
	}
}

// TestLockedAccountStillRefusesWrongPasswords checks that the brute force is
// still stopped while the account is under attack.
func TestLockedAccountStillRefusesWrongPasswords(t *testing.T) {
	e := newLoginTestEnv(t, 0)
	h := e.server.Handler()

	// 20 failures trip the counter of the account.
	for i := 0; i < 20; i++ {
		loginAs(h, "admin", "salah", "198.51.100.9:5555")
	}

	// The wrong password is refused from an address that never failed.
	code := loginAs(h, "admin", "salah", "203.0.113.9:5555")
	if code != http.StatusTooManyRequests {
		t.Errorf("password salah harus ditolak saat akun dikunci, kode %d", code)
	}
}

// TestLockedAccountSlowsDownTheAttempts checks that a flood against one
// account cannot use up the server.
//
// The password of a locked account is only checked every few seconds, so the
// work of the attacker does not grow with the number of attempts.
func TestLockedAccountSlowsDownTheAttempts(t *testing.T) {
	l := newLoginLimiter(3, nil)

	for i := 0; i < 6; i++ {
		l.fail("admin", "198.51.100.9")
	}
	if !l.locked("admin") {
		t.Fatal("akun harus terkunci setelah kegagalan cukup banyak")
	}

	// The address that already failed gets one chance, and the next
	// attempt is slowed down.
	if !l.mayTry("admin", "198.51.100.9") {
		t.Error("percobaan pertama dari alamat itu harus boleh")
	}
	if l.mayTry("admin", "198.51.100.9") {
		t.Error("percobaan berikutnya harus dibatasi")
	}
}

// TestNewAddressGetsItsChance checks that a user who signs in from a new
// connection is not kept waiting by an attack on somebody else's account.
func TestNewAddressGetsItsChance(t *testing.T) {
	l := newLoginLimiter(3, nil)

	for i := 0; i < 6; i++ {
		l.fail("admin", "198.51.100.9")
	}
	if !l.locked("admin") {
		t.Fatal("akun harus terkunci setelah kegagalan cukup banyak")
	}

	if !l.mayTry("admin", "203.0.113.9") {
		t.Error("alamat yang belum pernah gagal harus dapat kesempatan")
	}
}

// TestKnownAddressKeepsWorkingUnderAttack checks that the user who signed in
// before keeps working from his own address, even while the account is being
// guessed at.
func TestKnownAddressKeepsWorkingUnderAttack(t *testing.T) {
	e := newLoginTestEnv(t, 0)
	h := e.server.Handler()

	// The admin signs in from his own address once.
	if code := loginAs(h, "admin", "admin123", "203.0.113.9:5555"); code != http.StatusOK {
		t.Fatalf("login admin: %d", code)
	}

	// The attacker floods the account from another address.
	for i := 0; i < 30; i++ {
		loginAs(h, "admin", "salah", "198.51.100.9:5555")
	}

	// The admin signs in again from his own address.
	code := loginAs(h, "admin", "admin123", "203.0.113.9:5555")
	if code != http.StatusOK {
		t.Errorf("admin dari alamatnya sendiri harus tetap bisa masuk, kode %d", code)
	}
}

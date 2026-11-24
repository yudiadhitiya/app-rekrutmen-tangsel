package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/authtest"
	"rekrutmen-tangsel/backend/internal/settings"
)

// TestKnownGoodAccountKeepsWorkingWhenAddressIsBlocked checks that one
// attacker cannot lock out everybody behind a shared address.
//
// The address of a proxy is shared by every user of the application, so the
// limit of an address must not refuse an account that signed in successfully.
func TestKnownGoodAccountKeepsWorkingWhenAddressIsBlocked(t *testing.T) {
	dir := t.TempDir()
	authMgr := auth.NewManager(authtest.Users(), []byte("test-secret"))
	setStore, err := settings.NewStore(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}
	h := srv.Handler()

	// The admin signs in once from the address of the office, so the
	// account is known there.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login",
		strings.NewReader(`{"username":"admin","password":"admin123"}`))
	req.RemoteAddr = "10.0.0.7:5555"
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login admin: %d", rec.Code)
	}

	// An attacker from the same address fails 60 times with other names.
	for i := 0; i < 60; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/login",
			strings.NewReader(`{"username":"orang`+strings.Repeat("x", i%3)+`","password":"salah"}`))
		req.RemoteAddr = "10.0.0.7:5555"
		h.ServeHTTP(rec, req)
	}

	// The admin signs in again from the same address.
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/api/login",
		strings.NewReader(`{"username":"admin","password":"admin123"}`))
	req2.RemoteAddr = "10.0.0.7:5555"
	h.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Errorf("admin yang sudah masuk harus tetap bisa masuk, kode %d: %s", rec2.Code, rec2.Body.String())
	}
}

// TestUnknownAccountIsStillBlockedByAddressLimit checks that the limit of an
// address still stops a brute force against an account nobody signs in from.
func TestUnknownAccountIsStillBlockedByAddressLimit(t *testing.T) {
	dir := t.TempDir()
	authMgr := auth.NewManager(authtest.Users(), []byte("test-secret"))
	setStore, err := settings.NewStore(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}
	h := srv.Handler()

	for i := 0; i < 60; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/login",
			strings.NewReader(`{"username":"orang`+strings.Repeat("x", i%3)+`","password":"salah"}`))
		req.RemoteAddr = "10.0.0.9:5555"
		h.ServeHTTP(rec, req)
	}

	// A name nobody signed in with is refused from that address.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login",
		strings.NewReader(`{"username":"orangxx","password":"salahlagi"}`))
	req.RemoteAddr = "10.0.0.9:5555"
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("batas alamat harus tetap berlaku, kode %d: %s", rec.Code, rec.Body.String())
	}
}

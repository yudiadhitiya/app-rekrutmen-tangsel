package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestClientAddressIgnoresForwardedHeaderByDefault checks that a header a
// caller sets by itself is not used, so the limits of the login and of the
// registration cannot be sidestepped.
func TestClientAddressIgnoresForwardedHeaderByDefault(t *testing.T) {
	srv := &Server{cfg: Config{TrustProxy: false}}
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "10.0.0.5:12345"
	req.Header.Set("X-Forwarded-For", "1.2.3.4")

	if got := srv.clientAddress(req); got != "10.0.0.5" {
		t.Errorf("alamat harus dari koneksi, dapat %q", got)
	}
}

// TestClientAddressUsesLastForwardedEntry checks that the address of the
// proxy, the one it appended itself, is used when the proxy is trusted.
func TestClientAddressUsesLastForwardedEntry(t *testing.T) {
	srv := &Server{cfg: Config{TrustProxy: true}}
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "10.0.0.5:12345"
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 203.0.113.9")

	if got := srv.clientAddress(req); got != "203.0.113.9" {
		t.Errorf("alamat harus entri terakhir, dapat %q", got)
	}
}

// TestClientAddressFallsBackToConnection checks that the address of the
// connection is used when the trusted proxy sent no header.
func TestClientAddressFallsBackToConnection(t *testing.T) {
	srv := &Server{cfg: Config{TrustProxy: true}}
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "10.0.0.5:12345"

	if got := srv.clientAddress(req); got != "10.0.0.5" {
		t.Errorf("alamat harus dari koneksi, dapat %q", got)
	}
}

// TestClientAddressWithoutPort checks that an address without a port is
// still usable.
func TestClientAddressWithoutPort(t *testing.T) {
	srv := &Server{cfg: Config{TrustProxy: false}}
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "10.0.0.5"

	if got := srv.clientAddress(req); got != "10.0.0.5" {
		t.Errorf("alamat harus tetap utuh, dapat %q", got)
	}
}

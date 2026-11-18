package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/authtest"
	"rekrutmen-tangsel/backend/internal/settings"
)

// TestSecurityHeaders checks that every response carries the headers that
// protect the application in production.
func TestSecurityHeaders(t *testing.T) {
	h := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/health", nil))

	want := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "no-referrer",
	}
	for key, value := range want {
		if got := rec.Header().Get(key); got != value {
			t.Errorf("header %s = %q, ingin %q", key, got, value)
		}
	}
	csp := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "default-src 'self'") {
		t.Errorf("CSP tidak membatasi sumber: %q", csp)
	}
	if !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Errorf("CSP tidak melarang pembingkaian: %q", csp)
	}
	if !strings.Contains(csp, "https://cdn.sheetjs.com") {
		t.Errorf("CSP tidak mengizinkan pustaka Excel: %q", csp)
	}
	if got := rec.Header().Get("Strict-Transport-Security"); got != "" {
		t.Errorf("HSTS tidak boleh dikirim pada permintaan biasa, dapat %q", got)
	}
}

// TestSecurityHeadersHSTSOnSecureRequest checks that the HSTS header is sent
// when the request arrives over TLS or through a proxy that terminated TLS.
func TestSecurityHeadersHSTSOnSecureRequest(t *testing.T) {
	h := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Strict-Transport-Security"); got == "" {
		t.Error("HSTS tidak dikirim pada permintaan https")
	}
}

// TestSessionCookieSecureBehindProxy checks that the session cookie carries
// the Secure flag when the proxy reports a secure request.
func TestSessionCookieSecureBehindProxy(t *testing.T) {
	dir := t.TempDir()
	authMgr := auth.NewManager(authtest.Users(), []byte("test-secret"))
	setStore, err := settings.NewStore(dir + "/settings.json")
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}

	login := func(proto string) *http.Cookie {
		body := strings.NewReader(`{"username":"admin","password":"admin123"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/login", body)
		req.Header.Set("Content-Type", "application/json")
		if proto != "" {
			req.Header.Set("X-Forwarded-Proto", proto)
		}
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("login %q gagal: %d", proto, rec.Code)
		}
		for _, c := range rec.Result().Cookies() {
			if c.Name == "rekr_session" {
				return c
			}
		}
		t.Fatalf("cookie sesi tidak ditemukan pada login %q", proto)
		return nil
	}

	if c := login("https"); !c.Secure {
		t.Error("cookie sesi di belakang proxy https harus punya flag Secure")
	}
	if c := login("http"); c.Secure {
		t.Error("cookie sesi pada permintaan http biasa tidak harus punya flag Secure")
	}
}

// TestSessionCookieSecureForced checks that the flag can be forced, which is
// the setting of a server that always sits behind TLS.
func TestSessionCookieSecureForced(t *testing.T) {
	dir := t.TempDir()
	authMgr := auth.NewManager(authtest.Users(), []byte("test-secret"))
	authMgr.SetSecureCookies(true)
	setStore, err := settings.NewStore(dir + "/settings.json")
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/login",
		strings.NewReader(`{"username":"admin","password":"admin123"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login gagal: %d", rec.Code)
	}

	var found bool
	for _, c := range rec.Result().Cookies() {
		if c.Name == "rekr_session" {
			found = true
			if !c.Secure {
				t.Error("cookie sesi wajib punya flag Secure saat dipaksa")
			}
		}
	}
	if !found {
		t.Error("cookie sesi tidak ditemukan")
	}
}

// TestCORSOnlyAllowedOrigins checks that only the listed origins may call the
// API with the session of the user.
func TestCORSOnlyAllowedOrigins(t *testing.T) {
	dir := t.TempDir()
	authMgr := auth.NewManager(authtest.Users(), []byte("test-secret"))
	setStore, err := settings.NewStore(dir + "/settings.json")
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir, AllowedOrigins: []string{"https://contoh.go.id"}}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}

	call := func(origin string) string {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		req.Header.Set("Origin", origin)
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		return rec.Header().Get("Access-Control-Allow-Origin")
	}

	if got := call("https://contoh.go.id"); got != "https://contoh.go.id" {
		t.Errorf("origin yang diizinkan diberi jawaban %q", got)
	}
	if got := call("https://penyerang.go.id"); got != "" {
		t.Errorf("origin asing diberi izin %q", got)
	}
}

// TestLoginBodyIsNotLogged checks that the response of a login does not carry
// the password of the account.
func TestLoginBodyIsNotLogged(t *testing.T) {
	dir := t.TempDir()
	authMgr := auth.NewManager(authtest.Users(), []byte("test-secret"))
	setStore, err := settings.NewStore(dir + "/settings.json")
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/login",
		strings.NewReader(`{"username":"admin","password":"admin123"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("respons login bukan JSON: %v", err)
	}
	if _, ok := body["password"]; ok {
		t.Error("respons login memuat password")
	}
}

// TestRequestBodyLimit checks that a request with a body that is too large
// is refused, so one request cannot take the memory of the server.
func TestRequestBodyLimit(t *testing.T) {
	dir := t.TempDir()
	authMgr := auth.NewManager(authtest.Users(), []byte("test-secret"))
	setStore, err := settings.NewStore(dir + "/settings.json")
	if err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{DataDir: dir}, authMgr, setStore)
	if err != nil {
		t.Fatal(err)
	}

	big := strings.Repeat("a", MaxJSONBody+1024)
	body := strings.NewReader(`{"username":"admin","password":"` + big + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status body yang terlalu besar: %d, ingin 413", rec.Code)
	}
}

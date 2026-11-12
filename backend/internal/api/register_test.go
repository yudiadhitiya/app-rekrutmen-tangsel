package api

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/settings"
	"rekrutmen-tangsel/backend/internal/store/sheets"
)

// fakeRegistrar keeps the registered accounts in memory.
type fakeRegistrar struct {
	mu       sync.Mutex
	accounts []model.Registration
	err      error
}

func (f *fakeRegistrar) RegisterUser(reg model.Registration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return f.err
	}
	f.accounts = append(f.accounts, reg)
	return nil
}

// newRegisterEnv builds a server whose accounts are registered through the
// given registrar. A nil registrar means no account store is configured.
func newRegisterEnv(t *testing.T, registrar auth.Registrar) *testEnv {
	t.Helper()
	dir := t.TempDir()
	authMgr := auth.NewManager(auth.DefaultUsers(), []byte("test-secret"))
	if registrar != nil {
		authMgr.SetRegistrar(registrar)
	}
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

func validRegistration() map[string]string {
	return map[string]string{
		"username":  "anggota_baru",
		"password":  "rahasia123",
		"nama":      "Anggota Baru",
		"kecamatan": "Ciputat",
		"question1": model.SecurityQuestions[0],
		"answer1":   "Boneka",
		"question2": model.SecurityQuestions[1],
		"answer2":   "Bandung",
	}
}

func TestRegisterCreatesAnAccount(t *testing.T) {
	registrar := &fakeRegistrar{}
	env := newRegisterEnv(t, registrar)

	var out map[string]any
	code := env.request(http.MethodPost, "/api/register", validRegistration(), &out)
	if code != http.StatusCreated {
		t.Fatalf("kode: %d, ingin 201", code)
	}
	if out["status"] != "ACTIVE" {
		t.Errorf("status: %v", out["status"])
	}
	if len(registrar.accounts) != 1 {
		t.Fatalf("akun tersimpan: %d, ingin 1", len(registrar.accounts))
	}
	if registrar.accounts[0].Username != "anggota_baru" {
		t.Errorf("username: %q", registrar.accounts[0].Username)
	}
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(map[string]string)
	}{
		{"username pendek", func(m map[string]string) { m["username"] = "ab" }},
		{"password pendek", func(m map[string]string) { m["password"] = "12345" }},
		{"nama kosong", func(m map[string]string) { m["nama"] = "" }},
		{"kecamatan tak dikenal", func(m map[string]string) { m["kecamatan"] = "Bogor" }},
		{"pertanyaan tak dikenal", func(m map[string]string) { m["question1"] = "Pertanyaan lain?" }},
		{"pertanyaan sama", func(m map[string]string) { m["question2"] = m["question1"] }},
		{"jawaban kosong", func(m map[string]string) { m["answer1"] = "" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			registrar := &fakeRegistrar{}
			env := newRegisterEnv(t, registrar)
			body := validRegistration()
			tc.mutate(body)

			code := env.request(http.MethodPost, "/api/register", body, nil)
			if code != http.StatusBadRequest {
				t.Fatalf("kode: %d, ingin 400", code)
			}
			if len(registrar.accounts) != 0 {
				t.Errorf("akun tersimpan: %d, ingin 0", len(registrar.accounts))
			}
		})
	}
}

func TestRegisterReportsTakenUsername(t *testing.T) {
	env := newRegisterEnv(t, &fakeRegistrar{})

	// The sheet answers with this error when the username is already used.
	srv := env.server
	srv.auth.SetRegistrar(errorRegistrar{err: sheets.ErrUsernameTaken})

	code := env.request(http.MethodPost, "/api/register", validRegistration(), nil)
	if code != http.StatusConflict {
		t.Fatalf("kode: %d, ingin 409", code)
	}
}

// errorRegistrar always fails with the given error.
type errorRegistrar struct{ err error }

func (e errorRegistrar) RegisterUser(model.Registration) error { return e.err }

func TestRegisterWithoutAccountStoreIsUnavailable(t *testing.T) {
	env := newRegisterEnv(t, nil)

	code := env.request(http.MethodPost, "/api/register", validRegistration(), nil)
	if code != http.StatusServiceUnavailable {
		t.Fatalf("kode: %d, ingin 503", code)
	}
}

func TestRegisterOptionsArePublic(t *testing.T) {
	env := newRegisterEnv(t, &fakeRegistrar{})

	var out struct {
		Questions []string `json:"questions"`
		Kecamatan []string `json:"kecamatan"`
		Available bool     `json:"available"`
	}
	code := env.request(http.MethodGet, "/api/register/options", nil, &out)
	if code != http.StatusOK {
		t.Fatalf("kode: %d, ingin 200", code)
	}
	if len(out.Questions) != len(model.SecurityQuestions) {
		t.Errorf("pertanyaan: %d, ingin %d", len(out.Questions), len(model.SecurityQuestions))
	}
	if len(out.Kecamatan) == 0 {
		t.Error("kecamatan kosong")
	}
	if !out.Available {
		t.Error("ketersediaan harus true saat account store ada")
	}
}

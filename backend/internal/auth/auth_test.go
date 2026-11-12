package auth

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
)

// fakeLegacy stands in for the accounts of the previous application.
type fakeLegacy struct {
	user  model.User
	ok    bool
	err   error
	calls int
}

func (f *fakeLegacy) VerifyLegacy(username, password string) (model.User, bool, error) {
	f.calls++
	return f.user, f.ok, f.err
}

func testManager(t *testing.T) *Manager {
	t.Helper()
	users := DefaultUsers()
	m := NewManager(users, []byte("test-secret"))
	m.SetUsersPath(filepath.Join(t.TempDir(), "users.json"))
	return m
}

func TestVerifyLoginUsesLocalAccount(t *testing.T) {
	m := testManager(t)
	legacy := &fakeLegacy{ok: true, user: model.User{Username: "admin"}}
	m.SetLegacyAuthenticator(legacy)

	user, fromLegacy, err := m.VerifyLogin("admin", "admin123")
	if err != nil {
		t.Fatalf("VerifyLogin: %v", err)
	}
	if fromLegacy {
		t.Error("akun lokal tidak boleh dilaporkan sebagai warisan")
	}
	if user.Username != "admin" {
		t.Errorf("user: %+v", user)
	}
	if legacy.calls != 0 {
		t.Errorf("akun lokal tidak boleh menyentuh sumber warisan, dipanggil %d kali", legacy.calls)
	}
}

func TestVerifyLoginFallsBackToLegacy(t *testing.T) {
	m := testManager(t)
	legacy := &fakeLegacy{
		ok:   true,
		user: model.User{ID: "u1", Username: "anggota01", Nama: "Anggota Satu", Role: model.RolePetugas},
	}
	m.SetLegacyAuthenticator(legacy)

	user, fromLegacy, err := m.VerifyLogin("anggota01", "rahasia")
	if err != nil {
		t.Fatalf("VerifyLogin: %v", err)
	}
	if !fromLegacy {
		t.Error("login warisan harus dilaporkan sebagai warisan")
	}
	if user.Username != "anggota01" || user.Role != model.RolePetugas {
		t.Errorf("user: %+v", user)
	}
	if legacy.calls != 1 {
		t.Errorf("sumber warisan harus dipanggil tepat satu kali, dipanggil %d", legacy.calls)
	}
}

func TestVerifyLoginRejectsWrongLegacyPassword(t *testing.T) {
	m := testManager(t)
	m.SetLegacyAuthenticator(&fakeLegacy{ok: false})

	if _, _, err := m.VerifyLogin("anggota01", "salah"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("harus ErrInvalidCredentials, dapat %v", err)
	}
}

func TestVerifyLoginReportsStoreFailure(t *testing.T) {
	m := testManager(t)
	m.SetLegacyAuthenticator(&fakeLegacy{err: errors.New("sheet tidak terbaca")})

	_, _, err := m.VerifyLogin("anggota01", "rahasia")
	if err == nil {
		t.Fatal("kegagalan sumber warisan harus dilaporkan")
	}
	if errors.Is(err, ErrInvalidCredentials) {
		t.Error("kegagalan penyimpanan bukan salah password")
	}
}

func TestVerifyLoginWithoutLegacySource(t *testing.T) {
	m := testManager(t)
	if _, _, err := m.VerifyLogin("tidak-ada", "apa-saja"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("harus ErrInvalidCredentials, dapat %v", err)
	}
}

func TestMigrateLegacyUserStoresAccountLocally(t *testing.T) {
	m := testManager(t)
	path := filepath.Join(t.TempDir(), "users.json")
	m.SetUsersPath(path)
	legacy := &fakeLegacy{
		ok:   true,
		user: model.User{ID: "u1", Username: "anggota01", Nama: "Anggota Satu", Role: model.RolePetugas},
	}
	m.SetLegacyAuthenticator(legacy)

	if _, _, err := m.VerifyLogin("anggota01", "rahasia"); err != nil {
		t.Fatalf("VerifyLogin: %v", err)
	}
	if err := m.MigrateLegacyUser(model.User{ID: "u1", Username: "anggota01", Nama: "Anggota Satu", Role: model.RolePetugas}, "rahasia"); err != nil {
		t.Fatalf("MigrateLegacyUser: %v", err)
	}

	// The account is now local, so the sheet is not consulted any more.
	callsBefore := legacy.calls
	user, fromLegacy, err := m.VerifyLogin("anggota01", "rahasia")
	if err != nil {
		t.Fatalf("VerifyLogin setelah migrasi: %v", err)
	}
	if fromLegacy {
		t.Error("akun yang sudah dimigrasi bukan lagi login warisan")
	}
	if user.Username != "anggota01" {
		t.Errorf("user: %+v", user)
	}
	if legacy.calls != callsBefore {
		t.Error("akun lokal tidak boleh menyentuh sheet")
	}

	// The account was written to the file.
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("baca file akun: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("file akun kosong")
	}
	if _, _, err := m.VerifyLogin("anggota01", "salah"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatal("password salah harus ditolak setelah migrasi")
	}
}

// fakeProfiler stands in for the profile of an account of the previous
// application.
type fakeProfiler struct {
	user  model.User
	found bool
}

func (f *fakeProfiler) LegacyProfile(username string) (model.User, bool, error) {
	return f.user, f.found, nil
}

func TestVerifyLoginCompletesLevelOfMigratedAccount(t *testing.T) {
	m := testManager(t)
	path := filepath.Join(t.TempDir(), "users.json")
	m.SetUsersPath(path)

	// An account that was migrated before the levels were introduced: it
	// carries a role but no level.
	salt, hash, err := HashPassword("rahasia")
	if err != nil {
		t.Fatal(err)
	}
	m.users["anggota01"] = User{
		User:         model.User{ID: "a1b2c3d4-0000-0000-0000-000000000001", Username: "anggota01", Role: model.RoleAdmin},
		PasswordHash: hash,
		Salt:         salt,
	}

	profiler := &fakeProfiler{
		found: true,
		user: model.User{
			Username:  "anggota01",
			Role:      model.RoleAdmin,
			Level:     model.LevelAdminKecamatan,
			Kecamatan: "Serpong",
		},
	}
	m.SetLegacyAuthenticator(&fakeLegacy{})
	m.profiler = profiler

	user, _, err := m.VerifyLogin("anggota01", "rahasia")
	if err != nil {
		t.Fatalf("VerifyLogin: %v", err)
	}
	if user.Level != model.LevelAdminKecamatan {
		t.Errorf("level harus %d, dapat %d", model.LevelAdminKecamatan, user.Level)
	}
	if user.Kecamatan != "Serpong" {
		t.Errorf("kecamatan harus Serpong, dapat %q", user.Kecamatan)
	}

	// The completed profile is stored, so it survives the next sign in.
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("baca file akun: %v", err)
	}
	if !strings.Contains(string(b), `"level": 2`) {
		t.Errorf("level harus tersimpan di file akun: %s", string(b))
	}
}

func TestMigrateLegacyUserKeepsExistingAccount(t *testing.T) {
	m := testManager(t)
	m.SetUsersPath(filepath.Join(t.TempDir(), "users.json"))

	// The local account already exists, so nothing is overwritten.
	if err := m.MigrateLegacyUser(model.User{Username: "admin", Role: model.RolePetugas}, "rahasia"); err != nil {
		t.Fatalf("MigrateLegacyUser: %v", err)
	}
	if _, err := m.VerifyPassword("admin", "admin123"); err != nil {
		t.Fatalf("akun lokal harus tetap berlaku: %v", err)
	}
}

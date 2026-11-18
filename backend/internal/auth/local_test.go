package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"path/filepath"
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
)

// fakeStatusLegacy stands in for the account store of the previous application
// that also reports the access and the active flag of an account.
type fakeStatusLegacy struct {
	user   model.User
	active bool
	found  bool
	err    error

	password string
	calls    int
}

func (f *fakeStatusLegacy) VerifyLegacy(username, password string) (model.User, bool, error) {
	f.calls++
	if f.err != nil {
		return model.User{}, false, f.err
	}
	if !f.found || password != f.password {
		return model.User{}, false, nil
	}
	return f.user, true, nil
}

func (f *fakeStatusLegacy) LegacyAccountStatus(username string) (model.User, bool, error) {
	if f.err != nil {
		return model.User{}, false, f.err
	}
	if !f.found {
		return model.User{}, false, nil
	}
	return f.user, f.active, nil
}

// fakeRegistrar stands in for the account store that receives the new member
// accounts.
type fakeRegistrar struct {
	regs []model.Registration
	err  error
}

func (f *fakeRegistrar) RegisterUser(reg model.Registration) error {
	if f.err != nil {
		return f.err
	}
	f.regs = append(f.regs, reg)
	return nil
}

func registration() model.Registration {
	return model.Registration{
		Username:  "anggota01",
		Password:  "rahasia-kader",
		Nama:      "Anggota Satu",
		Kecamatan: "Ciputat",
	}
}

// TestRegisterUserStoresSaltedHashLocally checks that a new member account
// is stored locally with a salted hash, so its password is never checked
// against the unsalted hash of the previous application.
func TestRegisterUserStoresSaltedHashLocally(t *testing.T) {
	m := testManager(t)
	registrar := &fakeRegistrar{}
	m.SetRegistrar(registrar)

	reg := registration()
	if err := m.RegisterUser(reg); err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}
	if len(registrar.regs) != 1 {
		t.Fatalf("penyimpanan warisan harus dipanggil sekali, dipanggil %d kali", len(registrar.regs))
	}

	m.mu.RLock()
	u, ok := m.users["anggota01"]
	m.mu.RUnlock()
	if !ok {
		t.Fatal("akun baru harus disimpan secara lokal")
	}
	if u.Salt == "" || u.PasswordHash == "" {
		t.Fatal("akun baru harus punya salt dan hash")
	}
	// The unsalted SHA-256 of the password must not be the stored hash.
	if u.PasswordHash == sha256HexOf(reg.Password) {
		t.Error("hash unsalted tidak boleh dipakai")
	}
	if u.Level != 0 {
		t.Errorf("level harus diambil dari aplikasi lama, dapat %d", u.Level)
	}
}

// TestVerifyLoginUsesLocalHashOfRegisteredAccount checks that the sign in of
// a new member is decided by the salted local hash.
func TestVerifyLoginUsesLocalHashOfRegisteredAccount(t *testing.T) {
	m := testManager(t)
	m.SetRegistrar(&fakeRegistrar{})
	legacy := &fakeStatusLegacy{found: true, active: true, user: model.User{
		ID: model.LegacyUserID("anggota01"), Username: "anggota01", Role: model.RolePetugas, Level: model.LevelAnggota,
	}}
	m.SetLegacyAuthenticator(legacy)

	reg := registration()
	if err := m.RegisterUser(reg); err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}

	user, fromLegacy, err := m.VerifyLogin("anggota01", reg.Password)
	if err != nil {
		t.Fatalf("VerifyLogin: %v", err)
	}
	if fromLegacy {
		t.Error("akun baru harus memakai hash lokal")
	}
	if user.Username != "anggota01" {
		t.Errorf("user: %+v", user)
	}
}

// TestVerifyLoginTakesOverPasswordChangedInLegacyStore checks that a password
// that was replaced in the previous application still lets the member in, and
// that the local copy is hashed again.
func TestVerifyLoginTakesOverPasswordChangedInLegacyStore(t *testing.T) {
	m := testManager(t)
	path := filepath.Join(t.TempDir(), "users.json")
	m.SetUsersPath(path)

	salt, hash, err := HashPassword("password-lama")
	if err != nil {
		t.Fatal(err)
	}
	m.users["anggota01"] = User{
		User:         model.User{ID: model.LegacyUserID("anggota01"), Username: "anggota01", Role: model.RolePetugas},
		PasswordHash: hash,
		Salt:         salt,
	}
	legacy := &fakeStatusLegacy{
		found:    true,
		active:   true,
		password: "password-baru",
		user:     model.User{ID: model.LegacyUserID("anggota01"), Username: "anggota01", Role: model.RolePetugas},
	}
	m.SetLegacyAuthenticator(legacy)

	if _, _, err := m.VerifyLogin("anggota01", "password-baru"); err != nil {
		t.Fatalf("password baru harus diterima: %v", err)
	}

	// The local copy must now hold the salted hash of the new password.
	if _, err := m.VerifyPassword("anggota01", "password-baru"); err != nil {
		t.Fatalf("hash lokal harus diperbarui: %v", err)
	}
	if _, err := m.VerifyPassword("anggota01", "password-lama"); err == nil {
		t.Error("password lama harus ditolak")
	}
}

// TestVerifyLoginRefreshesAccessFromLegacyStore checks that a change of the
// role in the previous application takes effect at the next sign in.
func TestVerifyLoginRefreshesAccessFromLegacyStore(t *testing.T) {
	m := testManager(t)
	legacy := &fakeStatusLegacy{
		found:  true,
		active: true,
		user: model.User{
			ID: model.LegacyUserID("anggota01"), Username: "anggota01",
			Role: model.RoleAdmin, Level: model.LevelAdminKecamatan, Kecamatan: "Ciputat",
		},
	}
	m.SetLegacyAuthenticator(legacy)

	reg := registration()
	m.SetRegistrar(&fakeRegistrar{})
	if err := m.RegisterUser(reg); err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}

	user, _, err := m.VerifyLogin("anggota01", reg.Password)
	if err != nil {
		t.Fatalf("VerifyLogin: %v", err)
	}
	if user.Level != model.LevelAdminKecamatan {
		t.Errorf("level harus mengikuti aplikasi lama, dapat %d", user.Level)
	}
	if user.Kecamatan != "Ciputat" {
		t.Errorf("kecamatan harus mengikuti aplikasi lama, dapat %q", user.Kecamatan)
	}
}

// TestVerifyLoginRefusesDeactivatedAccount checks that an account that was
// deactivated in the previous application may not sign in.
func TestVerifyLoginRefusesDeactivatedAccount(t *testing.T) {
	m := testManager(t)
	legacy := &fakeStatusLegacy{
		found:  true,
		active: false,
		user:   model.User{ID: model.LegacyUserID("anggota01"), Username: "anggota01"},
	}
	m.SetLegacyAuthenticator(legacy)

	reg := registration()
	m.SetRegistrar(&fakeRegistrar{})
	if err := m.RegisterUser(reg); err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}

	if _, _, err := m.VerifyLogin("anggota01", reg.Password); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("akun nonaktif harus ditolak, dapat %v", err)
	}
}

// TestVerifyLoginKeepsLocalAccountWhenLegacyStoreFails checks that a store
// that cannot be reached does not stop a member from signing in.
func TestVerifyLoginKeepsLocalAccountWhenLegacyStoreFails(t *testing.T) {
	m := testManager(t)
	legacy := &fakeStatusLegacy{found: true, active: true, err: errors.New("jaringan gagal")}
	m.SetLegacyAuthenticator(legacy)

	reg := registration()
	m.SetRegistrar(&fakeRegistrar{})
	if err := m.RegisterUser(reg); err != nil {
		t.Fatalf("RegisterUser: %v", err)
	}

	if _, _, err := m.VerifyLogin("anggota01", reg.Password); err != nil {
		t.Fatalf("gangguan store lama tidak boleh menolak login: %v", err)
	}
}

// sha256HexOf returns the unsalted SHA-256 of a password, the way the
// previous application stores it.
func sha256HexOf(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

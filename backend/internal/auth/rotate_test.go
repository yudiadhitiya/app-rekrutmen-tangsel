package auth

import (
	"os"
	"path/filepath"
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
)

// writeUserFile stores the given accounts in a file.
func writeUserFile(t *testing.T, path string, users []User) {
	t.Helper()
	if err := saveUsers(path, users); err != nil {
		t.Fatalf("simpan akun: %v", err)
	}
}

// userWithPassword makes an account with the given password.
func userWithPassword(t *testing.T, username, password string) User {
	t.Helper()
	salt, hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	return User{
		User: model.User{
			ID:       username,
			Username: username,
			Nama:     username,
			Role:     model.RoleAdmin,
			Level:    model.LevelAdmin,
		},
		PasswordHash: hash,
		Salt:         salt,
	}
}

// verifyUser reports whether the password belongs to the account.
func verifyUser(u User, password string) bool {
	m := NewManager([]User{u}, []byte("secret-untuk-uji"))
	_, err := m.VerifyPassword(u.Username, password)
	return err == nil
}

// TestRotateDefaultPasswordsReplacesKnownDefaults checks that an account with
// a password that was written down in the source code gets a new password.
func TestRotateDefaultPasswordsReplacesKnownDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "users.json")
	writeUserFile(t, path, []User{userWithPassword(t, "admin", "admin123")})

	rotated, err := RotateDefaultPasswords(path)
	if err != nil {
		t.Fatalf("rotasi: %v", err)
	}
	if len(rotated) != 1 {
		t.Fatalf("harus ada satu akun yang dirotasi, dapat %d", len(rotated))
	}
	if rotated[0].Username != "admin" {
		t.Errorf("username: %q", rotated[0].Username)
	}
	if rotated[0].Password == "admin123" {
		t.Error("password baru tidak boleh sama dengan password bawaan")
	}

	users, err := LoadUsers(path)
	if err != nil {
		t.Fatalf("baca akun: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("harus ada satu akun, dapat %d", len(users))
	}
	if !verifyUser(users[0], rotated[0].Password) {
		t.Error("password baru harus cocok dengan hash yang tersimpan")
	}
	if verifyUser(users[0], "admin123") {
		t.Error("password lama tidak boleh berlaku lagi")
	}
}

// TestRotateDefaultPasswordsIsIdempotent checks that a second run changes
// nothing, because the new password is not a known default password.
func TestRotateDefaultPasswordsIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "users.json")
	writeUserFile(t, path, []User{
		userWithPassword(t, "admin", "admin123"),
		userWithPassword(t, "petugas", "petugas123"),
	})

	rotated, err := RotateDefaultPasswords(path)
	if err != nil {
		t.Fatalf("rotasi: %v", err)
	}
	if len(rotated) != 2 {
		t.Fatalf("harus ada dua akun yang dirotasi, dapat %d", len(rotated))
	}

	again, err := RotateDefaultPasswords(path)
	if err != nil {
		t.Fatalf("rotasi kedua: %v", err)
	}
	if len(again) != 0 {
		t.Fatalf("rotasi kedua tidak boleh mengubah apa pun, dapat %d", len(again))
	}
}

// TestRotateDefaultPasswordsKeepsOtherPasswords checks that an account with
// a password of its own is left alone.
func TestRotateDefaultPasswordsKeepsOtherPasswords(t *testing.T) {
	path := filepath.Join(t.TempDir(), "users.json")
	writeUserFile(t, path, []User{userWithPassword(t, "pengawas", "password-pribadi")})

	rotated, err := RotateDefaultPasswords(path)
	if err != nil {
		t.Fatalf("rotasi: %v", err)
	}
	if len(rotated) != 0 {
		t.Fatalf("tidak boleh ada akun yang dirotasi, dapat %d", len(rotated))
	}

	users, err := LoadUsers(path)
	if err != nil {
		t.Fatalf("baca akun: %v", err)
	}
	if !verifyUser(users[0], "password-pribadi") {
		t.Error("password sendiri harus tetap berlaku")
	}
}

// TestRotateDefaultPasswordsWithoutFile checks that a missing file is not an
// error, because there is nothing to rotate.
func TestRotateDefaultPasswordsWithoutFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "users.json")
	rotated, err := RotateDefaultPasswords(path)
	if err != nil {
		t.Fatalf("rotasi: %v", err)
	}
	if len(rotated) != 0 {
		t.Fatalf("tidak boleh ada akun yang dirotasi, dapat %d", len(rotated))
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("rotasi tidak boleh membuat berkas akun")
	}
}

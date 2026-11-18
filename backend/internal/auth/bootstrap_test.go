package auth

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseUsersCreatesOnlyTheBootstrapAdmin checks that a fresh
// installation gets exactly one account: the administrator of the bootstrap.
//
// The old behaviour created a pair of accounts with passwords that were
// written down in the source code, so anybody could sign in as an
// administrator of a fresh installation.
func TestParseUsersCreatesOnlyTheBootstrapAdmin(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "users.json")

	users, err := ParseUsers(path, Bootstrap{
		Username: "admin",
		Password: "password-yang-panjang",
		Nama:     "Administrator",
	})
	if err != nil {
		t.Fatalf("ParseUsers: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("harus ada tepat satu akun, dapat %d", len(users))
	}
	if users[0].Username != "admin" {
		t.Errorf("username: %q", users[0].Username)
	}
	if users[0].Role != "admin" || users[0].Level != 3 {
		t.Errorf("akun awal harus jadi admin penuh, dapat %+v", users[0].User)
	}
	if users[0].Salt == "" {
		t.Error("akun harus disimpan dengan salt")
	}

	// The password of the old default accounts must not work any more.
	m := NewManager(users, []byte("test-secret"))
	for _, old := range []string{"admin123", "petugas123"} {
		if err := m.verifyStored(users[0], old); err == nil {
			t.Errorf("password lama %q masih diterima", old)
		}
	}

	// The file holds the account, so a restart keeps it.
	again, err := ParseUsers(path, Bootstrap{Username: "lain", Password: "password-yang-panjang"})
	if err != nil {
		t.Fatalf("ParseUsers kedua: %v", err)
	}
	if len(again) != 1 || again[0].Username != "admin" {
		t.Errorf("akun yang tersimpan harus dipakai kembali, dapat %+v", again)
	}
}

// TestParseUsersKeepsExistingAccounts checks that an installation that
// already has accounts keeps them.
func TestParseUsersKeepsExistingAccounts(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "users.json")

	if _, err := ParseUsers(path, Bootstrap{Username: "admin", Password: "password-yang-panjang"}); err != nil {
		t.Fatalf("ParseUsers: %v", err)
	}
	users, err := ParseUsers(path, Bootstrap{Username: "admin", Password: "password-yang-panjang"})
	if err != nil {
		t.Fatalf("ParseUsers: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("harus ada satu akun, dapat %d", len(users))
	}
	m := NewManager(users, []byte("test-secret"))
	if err := m.verifyStored(users[0], "password-yang-panjang"); err != nil {
		t.Errorf("password bootstrap tidak cocok: %v", err)
	}
}

// TestParseUsersRejectsShortBootstrapPassword checks that a password that is
// too short for the root account is refused.
func TestParseUsersRejectsShortBootstrapPassword(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "users.json")

	_, err := ParseUsers(path, Bootstrap{Username: "admin", Password: "pendek"})
	if err == nil {
		t.Fatal("password yang terlalu pendek harus ditolak")
	}
	if _, statErr := os.Stat(path); statErr == nil {
		t.Error("berkas akun tidak boleh dibuat bila bootstrap ditolak")
	}
}

// TestParseUsersRejectsPasswordEqualToUsername checks that the name of the
// account cannot be its own password.
func TestParseUsersRejectsPasswordEqualToUsername(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "users.json")

	_, err := ParseUsers(path, Bootstrap{Username: "administrator", Password: "administrator"})
	if err == nil {
		t.Fatal("password yang sama dengan nama akun harus ditolak")
	}
}

// TestRandomPasswordIsRandomAndLongEnough checks that a generated password
// is long enough and not the same twice.
func TestRandomPasswordIsRandomAndLongEnough(t *testing.T) {
	first, err := RandomPassword(16)
	if err != nil {
		t.Fatalf("RandomPassword: %v", err)
	}
	if len(first) < MinBootstrapPasswordLength {
		t.Errorf("password terlalu pendek: %d", len(first))
	}
	second, err := RandomPassword(16)
	if err != nil {
		t.Fatalf("RandomPassword: %v", err)
	}
	if first == second {
		t.Error("dua password acak harus berbeda")
	}
	if strings.ContainsAny(first, " \t\n") {
		t.Error("password acak tidak boleh berisi spasi")
	}
}

// TestLoadUsersWithoutFileYieldsNoAccount checks that a missing file gives no
// account, so no account with a known password can come back.
func TestLoadUsersWithoutFileYieldsNoAccount(t *testing.T) {
	users, err := loadUsers(filepath.Join(t.TempDir(), "tidak-ada.json"))
	if err != nil {
		t.Fatalf("loadUsers: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("file yang tidak ada tidak boleh menghasilkan akun, dapat %d", len(users))
	}
}

// TestParseUsersRejectsUsernameWithADash checks that the name of the first
// administrator cannot look like the name of an account of the previous
// application.
func TestParseUsersRejectsUsernameWithADash(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "users.json")

	_, err := ParseUsers(path, Bootstrap{Username: "admin-1", Password: "password-yang-panjang"})
	if err == nil {
		t.Fatal("nama akun dengan tanda hubung harus ditolak")
	}
}

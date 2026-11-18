package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPrepareUsersCreatesOneAdminWithoutAKnownPassword checks that a fresh
// installation gets exactly one administrator and that its password is not
// one of the passwords that were written down in the source code before.
func TestPrepareUsersCreatesOneAdminWithoutAKnownPassword(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ADMIN_PASSWORD", "")
	t.Setenv("ADMIN_USERNAME", "admin")

	users, err := prepareUsers(dir)
	if err != nil {
		t.Fatalf("prepareUsers: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("harus ada tepat satu akun, dapat %d", len(users))
	}
	if users[0].Username != "admin" {
		t.Errorf("username: %q", users[0].Username)
	}
	if users[0].Level != 3 {
		t.Errorf("akun awal harus punya level penuh, dapat %d", users[0].Level)
	}

	// The password is written to a file that only the owner may read.
	passFile := filepath.Join(dir, "initial_admin_password.txt")
	info, err := os.Stat(passFile)
	if err != nil {
		t.Fatalf("berkas password awal harus ada: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("izin berkas password awal: %o, harus 600", perm)
	}
	content, err := os.ReadFile(passFile)
	if err != nil {
		t.Fatalf("baca berkas password awal: %v", err)
	}
	password := ""
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "password: ") {
			password = strings.TrimPrefix(line, "password: ")
		}
	}
	if len(password) < 12 {
		t.Errorf("password awal terlalu pendek: %d", len(password))
	}
	for _, old := range []string{"admin123", "petugas123"} {
		if password == old {
			t.Errorf("password awal tidak boleh %q", old)
		}
	}
}

// TestPrepareUsersUsesThePasswordFromTheEnvironment checks that the password
// of the first administrator may be given by the environment.
func TestPrepareUsersUsesThePasswordFromTheEnvironment(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ADMIN_PASSWORD", "password-dari-environment")
	t.Setenv("ADMIN_USERNAME", "pengawas")

	users, err := prepareUsers(dir)
	if err != nil {
		t.Fatalf("prepareUsers: %v", err)
	}
	if len(users) != 1 || users[0].Username != "pengawas" {
		t.Fatalf("akun: %+v", users)
	}
	// No file with a password is written, because the operator knows it.
	if _, err := os.Stat(filepath.Join(dir, "initial_admin_password.txt")); err == nil {
		t.Error("berkas password awal tidak boleh dibuat bila ADMIN_PASSWORD diisi")
	}
}

// TestPrepareUsersKeepsExistingAccounts checks that an installation that
// already has accounts keeps them and creates nothing.
func TestPrepareUsersKeepsExistingAccounts(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ADMIN_PASSWORD", "password-dari-environment")

	if _, err := prepareUsers(dir); err != nil {
		t.Fatalf("prepareUsers: %v", err)
	}
	users, err := prepareUsers(dir)
	if err != nil {
		t.Fatalf("prepareUsers kedua: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("harus ada satu akun, dapat %d", len(users))
	}
}

// TestPrepareUsersRefusesAShortPassword checks that a password that is too
// short for the root account stops the server.
func TestPrepareUsersRefusesAShortPassword(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ADMIN_PASSWORD", "pendek")

	if _, err := prepareUsers(dir); err == nil {
		t.Fatal("password yang terlalu pendek harus ditolak")
	}
}

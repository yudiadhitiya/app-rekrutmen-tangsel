package sheets

import (
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
)

// seedUsersTab fills the USERS tab of the fake spreadsheet.
func seedUsersTab(t *testing.T, fake *fakeSheets, rows [][]string) {
	t.Helper()
	fake.mu.Lock()
	defer fake.mu.Unlock()
	fake.tabs["USERS"] = 10
	fake.grid["USERS"] = rows
}

func legacyUsersRows(hash string) [][]string {
	return [][]string{
		{"USERNAME", "PASSWORD_HASH", "ROLE", "NAMA ANGGOTA", "NAMA KECAMATAN", "AKTIF",
			"SECURITY_Q1", "SECURITY_A1_HASH", "SECURITY_Q2", "SECURITY_A2_HASH"},
		{"anggota01", hash, "ANGGOTA", "Anggota Satu", "Pamulang", "Y", "", "", "", ""},
		{"admin02", hash, "SYSTEM ADMIN", "Admin Dua", "Ciputat", "Y", "", "", "", ""},
		{"adminkec", hash, "ADMIN KECAMATAN", "Admin Kecamatan", "Serpong", "Y", "", "", "", ""},
		{"nonaktif", hash, "ANGGOTA", "Nonaktif", "Ciputat", "N", "", "", "", ""},
		{"tanpahash", "", "ANGGOTA", "Tanpa Hash", "Ciputat", "Y", "", "", "", ""},
	}
}

func TestSheetsVerifyLegacyUser(t *testing.T) {
	s, fake := newTestStore(t)
	seedUsersTab(t, fake, legacyUsersRows(sha256Hex("rahasia")))

	user, ok, err := s.VerifyLegacy("anggota01", "rahasia")
	if err != nil {
		t.Fatalf("VerifyLegacy: %v", err)
	}
	if !ok {
		t.Fatal("kredensial yang benar harus diterima")
	}
	if user.Username != "anggota01" {
		t.Errorf("username: %q", user.Username)
	}
	if user.Nama != "Anggota Satu" {
		t.Errorf("nama: %q", user.Nama)
	}
	if user.Role != model.RolePetugas {
		t.Errorf("role ANGGOTA harus jadi petugas, dapat %q", user.Role)
	}
	if user.ID == "" {
		t.Error("id harus terisi")
	}
}

func TestSheetsVerifyLegacyUserIsCaseInsensitive(t *testing.T) {
	s, fake := newTestStore(t)
	seedUsersTab(t, fake, legacyUsersRows(sha256Hex("rahasia")))

	if _, ok, err := s.VerifyLegacy("  ANGGOTA01 ", "rahasia"); err != nil || !ok {
		t.Fatalf("username harus tidak membedakan huruf besar-kecil: ok=%v err=%v", ok, err)
	}
}

func TestSheetsVerifyLegacyUserRejectsBadCredentials(t *testing.T) {
	s, fake := newTestStore(t)
	seedUsersTab(t, fake, legacyUsersRows(sha256Hex("rahasia")))

	cases := []struct {
		username string
		password string
	}{
		{"anggota01", "salah"},
		{"nonaktif", "rahasia"},
		{"tanpahash", "rahasia"},
		{"tidak-ada", "rahasia"},
	}
	for _, c := range cases {
		if _, ok, err := s.VerifyLegacy(c.username, c.password); ok || err != nil {
			t.Errorf("%s/%s harus ditolak (ok=%v err=%v)", c.username, c.password, ok, err)
		}
	}
}

func TestSheetsVerifyLegacyUserMapsRoles(t *testing.T) {
	s, fake := newTestStore(t)
	seedUsersTab(t, fake, legacyUsersRows(sha256Hex("rahasia")))

	for _, username := range []string{"admin02", "adminkec"} {
		user, ok, err := s.VerifyLegacy(username, "rahasia")
		if err != nil || !ok {
			t.Fatalf("%s harus diterima: ok=%v err=%v", username, ok, err)
		}
		if user.Role != model.RoleAdmin {
			t.Errorf("%s harus jadi admin, dapat %q", username, user.Role)
		}
	}
}

func TestSheetsVerifyLegacyUserReportsMissingTab(t *testing.T) {
	s, _ := newTestStore(t)
	// The USERS tab does not exist in the fake spreadsheet.
	if _, _, err := s.VerifyLegacy("anggota01", "rahasia"); err == nil {
		t.Fatal("kegagalan membaca tab harus dilaporkan")
	}
}

func TestLegacyRoleMapping(t *testing.T) {
	cases := map[string]string{
		"SYSTEM ADMIN":    model.RoleAdmin,
		"ADMIN PTD":       model.RoleAdmin,
		"ADMIN DPD":       model.RoleAdmin,
		"ADMIN KECAMATAN": model.RoleAdmin,
		"ANGGOTA":         model.RolePetugas,
		"":                model.RolePetugas,
		"PERAN KUSTOM":    model.RolePetugas,
	}
	for role, want := range cases {
		if got := LegacyRole(role); got != want {
			t.Errorf("role %q harus jadi %q, dapat %q", role, want, got)
		}
	}
}

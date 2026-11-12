package sheets

import (
	"context"
	"os"
	"strings"
	"testing"
)

// TestLegacyLoginIntegration signs in against the real spreadsheet.
//
// It is skipped unless the spreadsheet is configured and a test account is
// given, so it never runs in an ordinary test run.
//
//	GOOGLE_SHEET_ID=... GOOGLE_KEY_FILE=... \
//	  GOOGLE_LEGACY_TEST_USER=... GOOGLE_LEGACY_TEST_PASSWORD=... \
//	  go test ./internal/store/sheets/ -run TestLegacyLoginIntegration -v
func TestLegacyLoginIntegration(t *testing.T) {
	sid := os.Getenv("GOOGLE_SHEET_ID")
	key := os.Getenv("GOOGLE_KEY_FILE")
	if sid == "" || key == "" {
		t.Skip("GOOGLE_SHEET_ID atau GOOGLE_KEY_FILE belum diisi")
	}
	username := os.Getenv("GOOGLE_LEGACY_TEST_USER")
	password := os.Getenv("GOOGLE_LEGACY_TEST_PASSWORD")
	if username == "" || password == "" {
		t.Skip("GOOGLE_LEGACY_TEST_USER atau GOOGLE_LEGACY_TEST_PASSWORD belum diisi")
	}

	s, err := New(Config{SpreadsheetID: sid, KeyFile: key})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	user, ok, err := s.VerifyLegacy(username, password)
	if err != nil {
		t.Fatalf("VerifyLegacy: %v", err)
	}
	if !ok {
		t.Fatalf("akun %q tidak diterima", username)
	}
	t.Logf("berhasil: username=%s role=%s nama=%s", user.Username, user.Role, user.Nama)
}

// TestLegacyHashFormat checks that the password hashes of the real USERS
// sheet have the shape this application expects.
//
// The previous application stored SHA-256 digests as lowercase hex, so every
// stored hash is 64 hexadecimal characters. The check reports only counts, so
// no hash is printed.
func TestLegacyHashFormat(t *testing.T) {
	sid := os.Getenv("GOOGLE_SHEET_ID")
	key := os.Getenv("GOOGLE_KEY_FILE")
	if sid == "" || key == "" {
		t.Skip("GOOGLE_SHEET_ID atau GOOGLE_KEY_FILE belum diisi")
	}
	s, err := New(Config{SpreadsheetID: sid, KeyFile: key})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx := context.Background()
	headers, rows, err := s.ReadTab(ctx, s.tabPengguna)
	if err != nil {
		t.Fatalf("ReadTab: %v", err)
	}
	idx := indexOfHeader(headers, colPasswordHash)
	if idx < 0 {
		t.Fatalf("tab %q tidak punya kolom %q", s.tabPengguna, colPasswordHash)
	}

	var total, filled, matching int
	for _, row := range rows {
		if idx >= len(row) {
			continue
		}
		total++
		v := strings.TrimSpace(row[idx])
		if v == "" {
			continue
		}
		filled++
		if len(v) == 64 && isHex(v) {
			matching++
		}
	}
	t.Logf("baris=%d hash terisi=%d hash berbentuk sha256 hex=%d", total, filled, matching)
	if filled > 0 && matching != filled {
		t.Errorf("hanya %d dari %d hash yang berbentuk sha256 hex", matching, filled)
	}
}

// isHex reports whether every character of s is a hexadecimal digit.
func isHex(s string) bool {
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9', r >= 'a' && r <= 'f':
		default:
			return false
		}
	}
	return true
}

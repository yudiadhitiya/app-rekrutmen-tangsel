package sheets

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
)

// DefaultTabPengguna is the tab of the previous application that holds its
// accounts.
const DefaultTabPengguna = "USERS"

// Legacy user column names, as written by the previous application.
const (
	colUsername     = "username"
	colPasswordHash = "password_hash"
	colRole         = "role"
	colNama         = "nama anggota"
	colAktif        = "aktif"
)

// LegacyUser is one account of the previous application.
type LegacyUser struct {
	// Username is the account name.
	Username string
	// PasswordHash is the SHA-256 hash of the password, as stored by the
	// previous application.
	PasswordHash string
	// Role is the role of the account.
	Role string
	// Nama is the display name of the account.
	Nama string
	// Kecamatan is the district of the account. It is only filled for a
	// district administrator.
	Kecamatan string
	// Aktif is the active flag: only "Y" accounts may sign in.
	Aktif string
}

// LegacyRole maps a role of the previous application to a role of this
// application.
//
// The previous application knew four roles; this application only distinguishes
// an administrator from an ordinary officer. The finer distinction between
// the levels is carried by LegacyLevel.
func LegacyRole(role string) string {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case "SYSTEM ADMIN", "ADMIN PTD", "ADMIN DPD", "ADMIN KECAMATAN", "TARBIYAH KECAMATAN":
		return model.RoleAdmin
	default:
		return model.RolePetugas
	}
}

// LegacyLevel maps a role of the previous application to an access level of
// this application.
//
// The previous application scoped the data per role: a district administrator
// only saw its own district, a member only its own records, and the other
// administrators saw everything. The level carries that distinction.
func LegacyLevel(role string) int {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case "SYSTEM ADMIN", "ADMIN PTD", "ADMIN DPD":
		return model.LevelAdmin
	case "ADMIN KECAMATAN", "TARBIYAH KECAMATAN":
		return model.LevelAdminKecamatan
	default:
		return model.LevelAnggota
	}
}

// VerifyLegacy checks a username and password against the accounts of the
// previous application.
//
// It implements the auth.LegacyAuthenticator contract, so the auth package
// can consult it for usernames that are unknown locally. The account is
// returned as a user of this application, with the legacy role mapped to a
// role of this application.
//
// The returned bool reports whether the credentials were accepted.
func (s *Store) VerifyLegacy(username, password string) (model.User, bool, error) {
	user, found, err := s.findLegacyUser(username)
	if err != nil {
		return model.User{}, false, err
	}
	if !found {
		return model.User{}, false, nil
	}
	// Only active accounts may sign in.
	if !strings.EqualFold(strings.TrimSpace(user.Aktif), "Y") {
		return model.User{}, false, nil
	}
	if user.PasswordHash == "" {
		return model.User{}, false, nil
	}
	if sha256Hex(password) != user.PasswordHash {
		return model.User{}, false, nil
	}
	return model.User{
		ID:        LegacyUserID(user.Username),
		Username:  user.Username,
		Nama:      user.Nama,
		Role:      LegacyRole(user.Role),
		Level:     LegacyLevel(user.Role),
		Kecamatan: user.Kecamatan,
	}, true, nil
}

// LegacyProfile returns the access level and the district of an account of
// the previous application, without checking a password.
//
// It is used to complete an account that was migrated before the levels were
// introduced.
func (s *Store) LegacyProfile(username string) (model.User, bool, error) {
	user, found, err := s.findLegacyUser(username)
	if err != nil {
		return model.User{}, false, err
	}
	if !found {
		return model.User{}, false, nil
	}
	return model.User{
		ID:        LegacyUserID(user.Username),
		Username:  user.Username,
		Nama:      user.Nama,
		Role:      LegacyRole(user.Role),
		Level:     LegacyLevel(user.Role),
		Kecamatan: user.Kecamatan,
	}, true, nil
}

// findLegacyUser looks up an account of the previous application by name.
func (s *Store) findLegacyUser(username string) (LegacyUser, bool, error) {
	want := strings.ToLower(strings.TrimSpace(username))
	if want == "" {
		return LegacyUser{}, false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	headers, rows, err := s.ReadTab(ctx, s.tabPengguna)
	if err != nil {
		return LegacyUser{}, false, fmt.Errorf("google sheets: baca tab pengguna: %w", err)
	}
	idxUsername := indexOfHeader(headers, colUsername)
	idxHash := indexOfHeader(headers, colPasswordHash)
	idxRole := indexOfHeader(headers, colRole)
	idxNama := indexOfHeader(headers, colNama)
	idxAktif := indexOfHeader(headers, colAktif)
	idxKecamatan := indexOfHeader(headers, colNamaKecamatan)
	if idxUsername < 0 || idxHash < 0 {
		return LegacyUser{}, false, fmt.Errorf("google sheets: tab %q tidak punya kolom %q dan %q",
			s.tabPengguna, colUsername, colPasswordHash)
	}

	for _, row := range rows {
		if idxUsername >= len(row) {
			continue
		}
		if strings.ToLower(strings.TrimSpace(row[idxUsername])) != want {
			continue
		}
		u := LegacyUser{Username: strings.TrimSpace(row[idxUsername])}
		if idxHash < len(row) {
			u.PasswordHash = strings.ToLower(strings.TrimSpace(row[idxHash]))
		}
		if idxRole < len(row) {
			u.Role = strings.TrimSpace(row[idxRole])
		}
		if idxNama < len(row) {
			u.Nama = strings.TrimSpace(row[idxNama])
		}
		if idxAktif < len(row) {
			u.Aktif = strings.TrimSpace(row[idxAktif])
		}
		if idxKecamatan >= 0 && idxKecamatan < len(row) {
			u.Kecamatan = strings.TrimSpace(row[idxKecamatan])
		}
		return u, true, nil
	}
	return LegacyUser{}, false, nil
}

// sha256Hex returns the lowercase hex SHA-256 digest of a password.
//
// The previous application stored passwords this way, without a salt.
func sha256Hex(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

// LegacyUserID builds a stable id for an account of the previous
// application, which does not carry one.
func LegacyUserID(username string) string {
	sum := sha256.Sum256([]byte("user:" + strings.ToLower(strings.TrimSpace(username))))
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// indexOfHeader returns the position of a column, matched without regard to
// case and punctuation, or -1 when the column is missing.
func indexOfHeader(headers []string, name string) int {
	want := normalizeHeader(name)
	for i, h := range headers {
		if normalizeHeader(h) == want {
			return i
		}
	}
	return -1
}

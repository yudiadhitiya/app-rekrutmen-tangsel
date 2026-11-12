package sheets

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
)

// Legacy security question columns, as written by the previous application.
const (
	colSecurityQ1     = "security_q1"
	colSecurityA1Hash = "security_a1_hash"
	colSecurityQ2     = "security_q2"
	colSecurityA2Hash = "security_a2_hash"
	colNamaKecamatan  = "nama kecamatan"
	legacyAktifValue  = "Y"
	registerTimeout   = 30 * time.Second
)

// ErrUsernameTaken is returned when the username is already used by another
// account of the previous application.
var ErrUsernameTaken = errors.New("username sudah digunakan")

// RegisterUser adds a member account to the USERS sheet.
//
// It follows the layout of the previous application: the password and the two
// security answers are stored as unsalted SHA-256 hashes, the role is ANGGOTA
// and the account is active right away, so no approval is needed.
func (s *Store) RegisterUser(reg model.Registration) error {
	ctx, cancel := context.WithTimeout(context.Background(), registerTimeout)
	defer cancel()

	headers, rows, err := s.ReadTab(ctx, s.tabPengguna)
	if err != nil {
		return fmt.Errorf("google sheets: baca tab pengguna: %w", err)
	}
	if len(headers) == 0 {
		return fmt.Errorf("google sheets: tab %q tidak punya baris header", s.tabPengguna)
	}

	idxUsername := indexOfHeader(headers, colUsername)
	if idxUsername < 0 {
		return fmt.Errorf("google sheets: tab %q tidak punya kolom %q", s.tabPengguna, colUsername)
	}

	want := strings.ToLower(strings.TrimSpace(reg.Username))
	for _, row := range rows {
		if idxUsername >= len(row) {
			continue
		}
		if strings.ToLower(strings.TrimSpace(row[idxUsername])) == want {
			return ErrUsernameTaken
		}
	}

	values := make([]string, len(headers))
	set := func(column, value string) {
		if i := indexOfHeader(headers, column); i >= 0 {
			values[i] = value
		}
	}

	set(colUsername, strings.TrimSpace(reg.Username))
	set(colPasswordHash, sha256Hex(reg.Password))
	set(colRole, model.RoleAnggotaLegacy)
	set(colNama, strings.TrimSpace(reg.Nama))
	set(colNamaKecamatan, strings.TrimSpace(reg.Kecamatan))
	set(colAktif, legacyAktifValue)
	set(colSecurityQ1, strings.TrimSpace(reg.Question1))
	set(colSecurityA1Hash, sha256Hex(normalizeAnswer(reg.Answer1)))
	set(colSecurityQ2, strings.TrimSpace(reg.Question2))
	set(colSecurityA2Hash, sha256Hex(normalizeAnswer(reg.Answer2)))

	if err := s.appendRow(ctx, s.tabPengguna, values); err != nil {
		return err
	}
	return nil
}

// normalizeAnswer prepares a security answer for hashing.
//
// The previous application lowercased the answer and reduced every run of
// spaces to one, so the same answer is hashed the same way here.
func normalizeAnswer(v string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(v))), " ")
}

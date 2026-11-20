// Package validate holds the input rules shared by the API layer.
package validate

import (
	"fmt"
	"regexp"
	"strings"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/seed"
)

var (
	digitsOnly = regexp.MustCompile(`^[0-9]+$`)
)

// Fields maps a form field name to an error message.
type Fields map[string]string

// Rekrutmen validates an incoming record. aktivitasNames is the set of master
// activity names currently allowed.
//
// It returns a Fields map; an empty map means the record is valid.
func Rekrutmen(r *model.Rekrutmen, aktivitasNames []string) Fields {
	f := Fields{}

	r.Nama = strings.TrimSpace(r.Nama)
	r.NoWA = strings.TrimSpace(r.NoWA)
	r.Alamat = strings.TrimSpace(r.Alamat)
	r.RT = strings.TrimSpace(r.RT)
	r.RW = strings.TrimSpace(r.RW)
	r.Kecamatan = strings.TrimSpace(r.Kecamatan)
	r.Kelurahan = strings.TrimSpace(r.Kelurahan)
	r.Aktivitas = strings.TrimSpace(r.Aktivitas)
	r.Keterangan = strings.TrimSpace(r.Keterangan)
	r.UPL = strings.TrimSpace(r.UPL)

	// Accept the retired "Program"/"Antara" spellings and fold them into the
	// current "Program Antara" option.
	r.Metode = model.NormalizeMetode(strings.TrimSpace(r.Metode))

	if r.Nama == "" {
		f["nama"] = "nama wajib diisi"
	}
	if r.NoWA == "" {
		f["no_wa"] = "nomor whatsapp wajib diisi"
	} else if !digitsOnly.MatchString(r.NoWA) {
		f["no_wa"] = "hanya boleh angka 0-9"
	}
	if r.RT != "" && !digitsOnly.MatchString(r.RT) {
		f["rt"] = "hanya boleh angka 0-9"
	}
	if r.RW != "" && !digitsOnly.MatchString(r.RW) {
		f["rw"] = "hanya boleh angka 0-9"
	}
	if r.Kecamatan == "" {
		f["kecamatan"] = "kecamatan wajib dipilih"
	} else if !seed.IsKecamatan(r.Kecamatan) {
		f["kecamatan"] = "kecamatan tidak valid"
	}
	if r.Kelurahan != "" {
		if r.Kecamatan == "" {
			f["kelurahan"] = "pilih kecamatan terlebih dahulu"
		} else if !seed.IsKelurahan(r.Kecamatan, r.Kelurahan) {
			f["kelurahan"] = "kelurahan tidak sesuai dengan kecamatan"
		}
	}
	if r.Aktivitas == "" {
		f["aktivitas"] = "aktivitas wajib dipilih"
	} else if !contains(aktivitasNames, r.Aktivitas) {
		f["aktivitas"] = "aktivitas tidak dikenal"
	}
	if !contains(model.TingkatPenerimaan, r.TingkatPenerimaan) {
		f["tingkat_penerimaan"] = "tingkat penerimaan wajib dipilih"
	}
	if !contains(model.Metode, r.Metode) {
		f["metode"] = "metode wajib dipilih"
	}
	if !contains(model.Status, r.Status) {
		f["status"] = "status wajib dipilih"
	}
	if !contains(model.UPL, r.UPL) {
		f["upl"] = "pilihan UPL tidak valid"
	}
	return f
}

// Aktivitas validates a master activity entry.
func Aktivitas(nama string) Fields {
	f := Fields{}
	nama = strings.TrimSpace(nama)
	if nama == "" {
		f["nama"] = "nama wajib diisi"
	}
	return f
}

// MinPasswordLength is the shortest password that is accepted.
//
// The password of an account is the only secret that protects the data of the
// people who registered, so it has to be long enough to survive a guess.
const MinPasswordLength = 12

// MaxPasswordLength is the longest password that is accepted.
//
// A password beyond this length brings nothing but a slow hash, so it is
// refused.
const MaxPasswordLength = 200

// weakPasswords are the passwords that are tried first in any guess.
//
// They are refused even when they are long enough, because a long password
// that everybody knows is no password at all.
var weakPasswords = map[string]bool{
	"password":     true,
	"password123":  true,
	"password1234": true,
	"admin":        true,
	"admin123":     true,
	"admin1234":    true,
	"petugas":      true,
	"petugas123":   true,
	"petugas1234":  true,
	"rahasia":      true,
	"rahasia123":   true,
	"rahasia1234":  true,
	"12345678":     true,
	"123456789":    true,
	"1234567890":   true,
	"123456789012": true,
	"qwerty":       true,
	"qwerty123":    true,
	"111111111111": true,
	"000000000000": true,
}

// Password validates a new password.
//
// The rules are the same everywhere a password is set: a password that is
// changed and a password of a new registration are both checked with it.
func Password(password string) Fields {
	f := Fields{}
	p := strings.TrimSpace(password)
	if len(p) < MinPasswordLength {
		f["password"] = fmt.Sprintf("password minimal %d karakter", MinPasswordLength)
		return f
	}
	if len(p) > MaxPasswordLength {
		f["password"] = fmt.Sprintf("password maksimal %d karakter", MaxPasswordLength)
		return f
	}
	if isWeakPassword(p) {
		f["password"] = "password ini terlalu umum, pilih yang lain"
	}
	return f
}

// weakFragments are the parts of a password that are tried first.
//
// A password that holds one of them is refused, whatever else is added to
// it, because the guess starts with that word anyway.
var weakFragments = []string{
	"password",
	"admin123",
	"petugas123",
	"rahasia123",
	"qwerty",
	"123456",
	"letmein",
	"welcome",
}

// isWeakPassword reports whether the password is one that is guessed first.
//
// The word "password" is refused in any spelling, because it is the first
// thing anybody tries, whatever is added to it.
func isWeakPassword(password string) bool {
	lower := strings.ToLower(password)
	if weakPasswords[lower] {
		return true
	}
	for _, fragment := range weakFragments {
		if strings.Contains(lower, fragment) {
			return true
		}
	}
	runes := []rune(lower)
	same := true
	for _, r := range runes {
		if r != runes[0] {
			same = false
			break
		}
	}
	return same
}

// PasswordDiffersFromUsername reports whether the password is the name of the
// account itself.
//
// Such a password is known to everybody who knows the account name, so it is
// refused.
func PasswordDiffersFromUsername(password, username string) bool {
	return strings.EqualFold(strings.TrimSpace(password), strings.TrimSpace(username))
}

func contains(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

// Error renders a field map as a single message, for logging.
func (f Fields) Error() string {
	if len(f) == 0 {
		return ""
	}
	parts := make([]string, 0, len(f))
	for k, v := range f {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(parts, "; ")
}

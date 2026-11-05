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

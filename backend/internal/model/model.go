// Package model holds the domain types shared by the API and the storage providers.
package model

import "time"

// Picklist enums. These are fixed by the PRD.
var (
	TingkatPenerimaan = []string{"Sangat Terbuka", "Terbuka", "Netral", "Kurang Terbuka", "Tertutup"}
	Metode            = []string{"Halaqoh", "Program Antara", "Dakwah Fardiyah"}
	Status            = []string{"Berhasil", "Proses", "Pending", "Gagal"}
	UPL               = []string{"", "UPLS", "UPLM"}
)

// MetodeLama maps the retired Metode values to their current replacement.
//
// The MVP initially shipped "Program" and "Antara" as two separate options.
// The list was later corrected to a single "Program Antara" option, so
// records stored with the old values are migrated on read.
var MetodeLama = map[string]string{
	"Program": "Program Antara",
	"Antara":  "Program Antara",
}

// NormalizeMetode maps a retired Metode value to its current replacement.
// Values that are already current are returned unchanged.
func NormalizeMetode(v string) string {
	if replacement, ok := MetodeLama[v]; ok {
		return replacement
	}
	return v
}

// MetodeMatchValues returns every stored Metode value that should match the
// given filter value.
//
// Filtering by "Program Antara" must also return records still stored with
// the retired "Program" or "Antara" spelling, so those are included too.
func MetodeMatchValues(filter string) []string {
	if filter == "" {
		return nil
	}
	if filter == "Program Antara" {
		return []string{"Program Antara", "Program", "Antara"}
	}
	return []string{filter}
}

// MetodeMatches reports whether a stored Metode value satisfies the filter.
func MetodeMatches(stored, filter string) bool {
	for _, v := range MetodeMatchValues(filter) {
		if stored == v {
			return true
		}
	}
	return false
}

// Rekrutmen is one recruitment record.
type Rekrutmen struct {
	ID                string    `json:"id"`
	Nama              string    `json:"nama"`
	NoWA              string    `json:"no_wa"`
	Alamat            string    `json:"alamat"`
	RT                string    `json:"rt"`
	RW                string    `json:"rw"`
	Kecamatan         string    `json:"kecamatan"`
	Kelurahan         string    `json:"kelurahan"`
	Aktivitas         string    `json:"aktivitas"`
	TingkatPenerimaan string    `json:"tingkat_penerimaan"`
	Metode            string    `json:"metode"`
	Status            string    `json:"status"`
	Keterangan        string    `json:"keterangan"`
	UPL               string    `json:"upl"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedBy         string    `json:"updated_by"`
}

// Aktivitas is a master entry for job/activity categories, manageable by an admin.
type Aktivitas struct {
	ID    string `json:"id"`
	Nama  string `json:"nama"`
	Aktif bool   `json:"aktif"`
}

// Filter narrows a rekrutmen listing.
type Filter struct {
	Q         string
	Kecamatan string
	Kelurahan string
	Metode    string
	Status    string
	Page      int
	Limit     int
}

// User is an authenticated account.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Nama     string `json:"nama"`
	Role     string `json:"role"`
}

// Roles.
const (
	RoleAdmin   = "admin"
	RolePetugas = "petugas"
)

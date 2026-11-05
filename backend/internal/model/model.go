// Package model holds the domain types shared by the API and the storage providers.
package model

import "time"

// Picklist enums. These are fixed by the PRD.
var (
	TingkatPenerimaan = []string{"Sangat Terbuka", "Terbuka", "Netral", "Kurang Terbuka", "Tertutup"}
	Metode            = []string{"Halaqoh", "Program", "Antara", "Dakwah Fardiyah"}
	Status            = []string{"Berhasil", "Proses", "Pending", "Gagal"}
	UPL               = []string{"", "UPLS", "UPLM"}
)

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

// Package seed holds the master reference data shipped with the app.
package seed

import (
	_ "embed"
	"encoding/json"
)

//go:embed wilayah.json
var wilayahJSON []byte

// Kecamatan is one district with its villages.
type Kecamatan struct {
	ID        string      `json:"id"`
	Nama      string      `json:"nama"`
	Kelurahan []Kelurahan `json:"kelurahan"`
}

// Kelurahan is one village.
type Kelurahan struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

// All returns the master list of districts and villages.
func All() ([]Kecamatan, error) {
	var out []Kecamatan
	if err := json.Unmarshal(wilayahJSON, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// KelurahanByKecamatan returns the villages belonging to a district.
func KelurahanByKecamatan(kecamatan string) []Kelurahan {
	all, err := All()
	if err != nil {
		return nil
	}
	for _, k := range all {
		if k.ID == kecamatan {
			return k.Kelurahan
		}
	}
	return nil
}

// IsKecamatan reports whether a district name is part of the master list.
func IsKecamatan(name string) bool {
	all, err := All()
	if err != nil {
		return false
	}
	for _, k := range all {
		if k.ID == name {
			return true
		}
	}
	return false
}

// IsKelurahan reports whether a village name is valid for the given district.
func IsKelurahan(kecamatan, kelurahan string) bool {
	for _, k := range KelurahanByKecamatan(kecamatan) {
		if k.ID == kelurahan {
			return true
		}
	}
	return false
}

// Aktivitas is the initial master list of job/activity categories, seeded on
// first run so the picklist is never empty. Admins may rename or add to it.
var Aktivitas = []string{
	"Agen Asuransi",
	"Buruh Pabrik",
	"Driver Ojek Online",
	"Driver Taksi Online",
	"Guru",
	"Ibu Rumah Tangga",
	"Karyawan Swasta",
	"Pedagang",
	"Pedagang Online",
	"Pelajar / Mahasiswa",
	"Pensiunan",
	"Petani",
	"Perawat",
	"Pedagang Keeping",
	"Wiraswasta",
	"Ustadz / Daiyah",
	"Tukang Bangun",
	"Pengusaha",
}

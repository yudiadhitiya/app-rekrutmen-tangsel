// Command e2eseed makes the accounts used by the end-to-end checks.
//
// It writes a users.json with the test accounts of every access level.
package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
)

func main() {
	dir := os.Args[1]
	specs := []struct {
		username, password, nama, role string
		level                          int
		kecamatan                      string
	}{
		{"petugas", "petugas123", "Petugas", model.RolePetugas, model.LevelAnggota, ""},
		{"admin", "admin123", "Administrator", model.RoleAdmin, model.LevelAdmin, ""},
		{"pwpetugas", "pwpetugas123", "Petugas Uji", model.RolePetugas, model.LevelAnggota, ""},
		{"pwkec", "pwkec12345", "Kecamatan Uji", model.RoleAdmin, model.LevelAdminKecamatan, "Serpong"},
		{"pwadmin", "pwadmin123", "Admin Uji", model.RoleAdmin, model.LevelAdmin, ""},
	}
	users := []auth.User{}
	for _, s := range specs {
		salt, hash, err := auth.HashPassword(s.password)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, auth.User{
			User: model.User{
				ID:        s.username,
				Username:  s.username,
				Nama:      s.nama,
				Role:      s.role,
				Level:     s.level,
				Kecamatan: s.kecamatan,
			},
			PasswordHash: hash,
			Salt:         salt,
		})
	}
	b, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "users.json"), b, 0o600); err != nil {
		log.Fatal(err)
	}
	log.Printf("akun uji ditulis ke %s", dir)
}

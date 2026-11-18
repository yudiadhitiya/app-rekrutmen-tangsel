// Package authtest holds the accounts used by the tests of the other
// packages.
//
// The accounts carry passwords that are written down here, so they must only
// ever be used by a test. The application itself never creates an account
// with a known password.
package authtest

import (
	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
)

// Users returns the accounts the tests sign in with: a full administrator
// and an ordinary officer.
func Users() []auth.User {
	out := []auth.User{}
	for _, spec := range []struct {
		username, password, nama, role string
		level                          int
	}{
		{"admin", "admin123", "Administrator", model.RoleAdmin, model.LevelAdmin},
		{"petugas", "petugas123", "Petugas", model.RolePetugas, model.LevelAnggota},
	} {
		salt, hash, err := auth.HashPassword(spec.password)
		if err != nil {
			continue
		}
		out = append(out, auth.User{
			User: model.User{
				ID:       spec.username,
				Username: spec.username,
				Nama:     spec.nama,
				Role:     spec.role,
				Level:    spec.level,
			},
			PasswordHash: hash,
			Salt:         salt,
		})
	}
	return out
}

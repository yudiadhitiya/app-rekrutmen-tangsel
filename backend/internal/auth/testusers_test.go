package auth

import (
	"rekrutmen-tangsel/backend/internal/model"
)

// testUsers returns the accounts the tests of this package sign in with: a
// full administrator and an ordinary officer.
//
// The passwords are written down here, so these accounts must only be used
// by a test. The application itself never creates an account with a known
// password.
func testUsers() []User {
	out := []User{}
	for _, spec := range []struct {
		username, password, nama, role string
		level                          int
	}{
		{"admin", "admin123", "Administrator", model.RoleAdmin, model.LevelAdmin},
		{"petugas", "petugas123", "Petugas", model.RolePetugas, model.LevelAnggota},
	} {
		salt, hash, err := HashPassword(spec.password)
		if err != nil {
			continue
		}
		out = append(out, User{
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

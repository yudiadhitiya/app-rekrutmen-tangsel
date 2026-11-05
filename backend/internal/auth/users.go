package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"rekrutmen-tangsel/backend/internal/model"
)

// DefaultUsers are created on first run so the app is usable immediately.
// Passwords are documented in the README and should be changed in production.
func DefaultUsers() []User {
	out := []User{}
	for _, spec := range []struct {
		username, password, nama, role string
	}{
		{"admin", "admin123", "Administrator", model.RoleAdmin},
		{"petugas", "petugas123", "Petugas", model.RolePetugas},
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
			},
			PasswordHash: hash,
			Salt:         salt,
		})
	}
	return out
}

func loadUsers(path string) ([]User, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		users := DefaultUsers()
		if err := saveUsers(path, users); err != nil {
			return nil, err
		}
		return users, nil
	}
	if err != nil {
		return nil, err
	}
	var users []User
	if err := json.Unmarshal(b, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func saveUsers(path string, users []User) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}

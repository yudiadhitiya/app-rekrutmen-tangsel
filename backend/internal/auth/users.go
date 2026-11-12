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
	// An account that was stored before the access levels were introduced
	// carries none. A local account keeps the access its role gives it. An
	// account of the previous application keeps it empty, so its level can be
	// taken from the previous application instead of being guessed.
	for i := range users {
		if users[i].Level == 0 && !isLegacyID(users[i].ID) {
			users[i].Level = users[i].LevelOrDefault()
		}
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

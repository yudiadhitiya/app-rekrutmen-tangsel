package auth

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"rekrutmen-tangsel/backend/internal/model"
)

// MinBootstrapPasswordLength is the shortest password that may be given to
// the first administrator of a fresh installation.
//
// The length is larger than the one of an ordinary registration, because
// this password is the root of the application: it may reach every record
// and every setting.
const MinBootstrapPasswordLength = 12

// Bootstrap describes the first administrator of a fresh installation.
//
// It is used when the account file does not exist yet. The account is
// created with a salted hash, so the password is never written to the disk
// in plain text.
type Bootstrap struct {
	// Username is the name of the first administrator.
	Username string
	// Password is the password of the first administrator.
	Password string
	// Nama is the display name of the first administrator.
	Nama string
}

// NeedsBootstrap reports whether the account file still has to be created.
//
// A file that does not exist, or one that holds no account at all, needs the
// first administrator of the installation.
func NeedsBootstrap(path string) (bool, error) {
	users, err := loadUsers(path)
	if err != nil {
		return false, err
	}
	return len(users) == 0, nil
}

// loadUsers reads the accounts from the given file.
//
// A file that does not exist yields no accounts, so the caller decides how
// the first administrator is created. No account with a known password is
// ever made here.
func loadUsers(path string) ([]User, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
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

// bootstrapAdmin creates the first administrator of a fresh installation.
//
// The account is stored with a salted PBKDF2 hash, so the password is never
// written to the disk in plain text.
func bootstrapAdmin(boot Bootstrap) (User, error) {
	username := strings.TrimSpace(boot.Username)
	if username == "" {
		username = "admin"
	}
	// The id of a local account is its own name, and an id with a dash marks
	// an account of the previous application. A name with a dash would
	// therefore make this account look like one of those, so it is refused.
	if strings.Contains(username, "-") {
		return User{}, errors.New("nama akun administrator tidak boleh berisi tanda hubung")
	}
	nama := strings.TrimSpace(boot.Nama)
	if nama == "" {
		nama = "Administrator"
	}
	if err := ValidateBootstrapPassword(boot.Password, username); err != nil {
		return User{}, err
	}

	salt, hash, err := HashPassword(boot.Password)
	if err != nil {
		return User{}, err
	}
	return User{
		User: model.User{
			ID:       username,
			Username: username,
			Nama:     nama,
			Role:     model.RoleAdmin,
			Level:    model.LevelAdmin,
		},
		PasswordHash: hash,
		Salt:         salt,
	}, nil
}

// ValidateBootstrapPassword checks the password of the first administrator.
//
// The password is the root of the application, so it must be long enough and
// must not be the name of the account itself.
func ValidateBootstrapPassword(password, username string) error {
	if len(password) < MinBootstrapPasswordLength {
		return fmt.Errorf("password administrator minimal %d karakter", MinBootstrapPasswordLength)
	}
	if strings.EqualFold(strings.TrimSpace(username), strings.TrimSpace(password)) {
		return errors.New("password administrator tidak boleh sama dengan nama akun")
	}
	return nil
}

// passwordAlphabet is the set of characters a generated password is built
// from. The characters that look alike are left out, so a password that is
// written down is not misread.
const passwordAlphabet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// RandomPassword returns a random password of the given length.
//
// It is used when no password was given for the first administrator, so the
// application is usable without a password that is written down somewhere.
func RandomPassword(length int) (string, error) {
	if length < MinBootstrapPasswordLength {
		length = MinBootstrapPasswordLength
	}
	out := make([]byte, 0, length)
	max := big.NewInt(int64(len(passwordAlphabet)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("auth: %w", err)
		}
		out = append(out, passwordAlphabet[n.Int64()])
	}
	return string(out), nil
}

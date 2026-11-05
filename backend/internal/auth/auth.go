// Package auth provides username/password login with signed session cookies.
//
// Passwords are stored as PBKDF2-HMAC-SHA256 hashes and sessions are carried
// in an HMAC-signed cookie, so no server-side session store is needed.
package auth

import (
	"crypto/hmac"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
)

const (
	cookieName   = "rekr_session"
	sessionTTL   = 12 * time.Hour
	pbkdf2Iters  = 120000
	pbkdf2KeyLen = 32
)

// ErrInvalidCredentials is returned for a bad username or password.
var ErrInvalidCredentials = errors.New("username atau password salah")

// User is a stored account with a password hash.
type User struct {
	model.User
	PasswordHash string `json:"password_hash"`
	Salt         string `json:"salt"`
}

// Manager authenticates users and issues session cookies.
type Manager struct {
	users  map[string]User // keyed by username
	secret []byte
}

// NewManager builds a Manager from a list of users and a signing secret.
func NewManager(users []User, secret []byte) *Manager {
	m := &Manager{users: make(map[string]User, len(users)), secret: secret}
	for _, u := range users {
		m.users[strings.ToLower(u.Username)] = u
	}
	return m
}

// HashPassword returns a salt and PBKDF2 hash for a plaintext password.
func HashPassword(password string) (salt, hash string, err error) {
	s := make([]byte, 16)
	if _, err := rand.Read(s); err != nil {
		return "", "", err
	}
	key, err := pbkdf2.Key(sha256.New, password, s, pbkdf2Iters, pbkdf2KeyLen)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(s), hex.EncodeToString(key), nil
}

// VerifyPassword checks a plaintext password against a stored hash.
func (m *Manager) VerifyPassword(username, password string) (model.User, error) {
	u, ok := m.users[strings.ToLower(username)]
	if !ok {
		// Spend comparable time so a missing user is not distinguishable by timing.
		_, _ = pbkdf2.Key(sha256.New, password, []byte("dummy-salt"), pbkdf2Iters, pbkdf2KeyLen)
		return model.User{}, ErrInvalidCredentials
	}
	salt, err := hex.DecodeString(u.Salt)
	if err != nil {
		return model.User{}, ErrInvalidCredentials
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, pbkdf2Iters, pbkdf2KeyLen)
	if err != nil {
		return model.User{}, err
	}
	want, err := hex.DecodeString(u.PasswordHash)
	if err != nil {
		return model.User{}, ErrInvalidCredentials
	}
	if subtle.ConstantTimeCompare(key, want) != 1 {
		return model.User{}, ErrInvalidCredentials
	}
	return u.User, nil
}

// SetSession writes a signed session cookie for the user.
func (m *Manager) SetSession(w http.ResponseWriter, user model.User) error {
	payload := struct {
		Username string `json:"u"`
		Role     string `json:"r"`
		Expires  int64  `json:"e"`
	}{
		Username: user.Username,
		Role:     user.Role,
		Expires:  time.Now().Add(sessionTTL).Unix(),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	enc := base64.RawURLEncoding.EncodeToString(b)
	sig := m.sign(enc)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    enc + "." + sig,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(sessionTTL),
	})
	return nil
}

// ClearSession expires the session cookie.
func (m *Manager) ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// UserFromRequest returns the authenticated user, or false when there is no
// valid session.
func (m *Manager) UserFromRequest(r *http.Request) (model.User, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return model.User{}, false
	}
	parts := strings.SplitN(c.Value, ".", 2)
	if len(parts) != 2 {
		return model.User{}, false
	}
	if !hmac.Equal([]byte(parts[1]), []byte(m.sign(parts[0]))) {
		return model.User{}, false
	}
	b, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return model.User{}, false
	}
	var payload struct {
		Username string `json:"u"`
		Role     string `json:"r"`
		Expires  int64  `json:"e"`
	}
	if err := json.Unmarshal(b, &payload); err != nil {
		return model.User{}, false
	}
	if time.Now().Unix() > payload.Expires {
		return model.User{}, false
	}
	u, ok := m.users[strings.ToLower(payload.Username)]
	if !ok {
		return model.User{}, false
	}
	return u.User, true
}

func (m *Manager) sign(value string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// ParseUsers reads the users file, creating a default admin/petugas pair when
// the file does not exist yet.
func ParseUsers(path string) ([]User, error) {
	users, err := loadUsers(path)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// RandomSecret returns a random secret used to sign sessions.
func RandomSecret() ([]byte, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("auth: %w", err)
	}
	return b, nil
}

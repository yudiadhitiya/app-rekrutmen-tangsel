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
	"sort"
	"strings"
	"sync"
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

// ErrNoRegistrar is returned when no account store for new members is
// configured, so a self-registration cannot be written anywhere.
var ErrNoRegistrar = errors.New("pendaftaran anggota belum tersedia")

// LegacyAuthenticator verifies credentials against an external account store,
// such as the USERS sheet of the previous application.
//
// It is consulted only for usernames that are not known locally, so an
// account that was already migrated always uses the local hash.
type LegacyAuthenticator interface {
	// VerifyLegacy reports whether the credentials are accepted, and returns
	// the account as a user of this application.
	VerifyLegacy(username, password string) (model.User, bool, error)
}

// Registrar creates a new account in an external account store, such as the
// USERS sheet of the previous application.
type Registrar interface {
	// RegisterUser stores a new member account.
	RegisterUser(reg model.Registration) error
}

// User is a stored account with a password hash.
type User struct {
	model.User
	PasswordHash string `json:"password_hash"`
	Salt         string `json:"salt"`
}

// Manager authenticates users and issues session cookies.
type Manager struct {
	mu     sync.RWMutex
	users  map[string]User // keyed by lowercase username
	secret []byte

	// legacy verifies accounts of the previous application.
	legacy LegacyAuthenticator

	// registrar stores new member accounts.
	registrar Registrar

	// path is the file the accounts are stored in. It is empty when the
	// manager must not write to disk.
	path string
}

// NewManager builds a Manager from a list of users and a signing secret.
func NewManager(users []User, secret []byte) *Manager {
	m := &Manager{users: make(map[string]User, len(users)), secret: secret}
	for _, u := range users {
		m.users[strings.ToLower(u.Username)] = u
	}
	return m
}

// SetLegacyAuthenticator sets the source that verifies accounts of the
// previous application.
func (m *Manager) SetLegacyAuthenticator(a LegacyAuthenticator) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.legacy = a
}

// SetRegistrar sets the store that keeps new member accounts.
func (m *Manager) SetRegistrar(r Registrar) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registrar = r
}

// CanRegister reports whether new member accounts can be stored.
func (m *Manager) CanRegister() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.registrar != nil
}

// RegisterUser stores a new member account.
// The account is written to the account store of the previous application,
// so the member can sign in with it in both applications. It stays there
// until the first sign in, which migrates it to the local store.
func (m *Manager) RegisterUser(reg model.Registration) error {
	m.mu.RLock()
	registrar := m.registrar
	m.mu.RUnlock()

	if registrar == nil {
		return ErrNoRegistrar
	}
	return registrar.RegisterUser(reg)
}

// SetUsersPath sets the file the accounts are stored in, so a migrated
// account can be written to it.
func (m *Manager) SetUsersPath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.path = path
}

// VerifyLogin checks a password against the stored accounts.
//
// Local accounts are checked first. When the username is unknown locally,
// the accounts of the previous application are consulted, so a user can sign
// in with the credentials they already had.
//
// The returned bool reports whether the login came from the previous
// application, which means the account can be migrated.
func (m *Manager) VerifyLogin(username, password string) (model.User, bool, error) {
	key := strings.ToLower(strings.TrimSpace(username))

	m.mu.RLock()
	u, ok := m.users[key]
	legacy := m.legacy
	m.mu.RUnlock()

	if ok {
		if err := m.verifyStored(u, password); err != nil {
			return model.User{}, false, err
		}
		return u.User, false, nil
	}

	if legacy == nil {
		// Spend comparable time so a missing user is not distinguishable by
		// timing.
		_, _ = pbkdf2.Key(sha256.New, password, []byte("dummy-salt"), pbkdf2Iters, pbkdf2KeyLen)
		return model.User{}, false, ErrInvalidCredentials
	}

	user, valid, err := legacy.VerifyLegacy(username, password)
	if err != nil {
		return model.User{}, false, err
	}
	if !valid {
		return model.User{}, false, ErrInvalidCredentials
	}
	return user, true, nil
}

// MigrateLegacyUser stores an account of the previous application locally.
//
// The account keeps the password the user just typed, but it is stored with a
// salted PBKDF2 hash, so the next sign in no longer needs the sheet.
func (m *Manager) MigrateLegacyUser(user model.User, password string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := strings.ToLower(strings.TrimSpace(user.Username))
	if _, ok := m.users[key]; ok {
		return nil
	}
	salt, hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	u := User{
		User:         user,
		PasswordHash: hash,
		Salt:         salt,
	}
	m.users[key] = u

	if m.path == "" {
		return nil
	}
	all := make([]User, 0, len(m.users))
	for _, item := range m.users {
		all = append(all, item)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Username < all[j].Username })
	return saveUsers(m.path, all)
}

// verifyStored checks a password against a stored account.
func (m *Manager) verifyStored(u User, password string) error {
	salt, err := hex.DecodeString(u.Salt)
	if err != nil {
		return ErrInvalidCredentials
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, pbkdf2Iters, pbkdf2KeyLen)
	if err != nil {
		return err
	}
	want, err := hex.DecodeString(u.PasswordHash)
	if err != nil {
		return ErrInvalidCredentials
	}
	if subtle.ConstantTimeCompare(key, want) != 1 {
		return ErrInvalidCredentials
	}
	return nil
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
	m.mu.RLock()
	u, ok := m.users[strings.ToLower(username)]
	m.mu.RUnlock()
	if !ok {
		// Spend comparable time so a missing user is not distinguishable by timing.
		_, _ = pbkdf2.Key(sha256.New, password, []byte("dummy-salt"), pbkdf2Iters, pbkdf2KeyLen)
		return model.User{}, ErrInvalidCredentials
	}
	if err := m.verifyStored(u, password); err != nil {
		return model.User{}, err
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
	m.mu.RLock()
	u, ok := m.users[strings.ToLower(payload.Username)]
	m.mu.RUnlock()
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

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
	"os"
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

// ErrAccountNotFound is returned when the account is not known to any store.
var ErrAccountNotFound = errors.New("akun tidak ditemukan")

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

// LegacyLookup reports whether an account of the previous application exists,
// without checking its password.
//
// It lets the manager tell an account that the previous application does not
// know from an account whose password is simply wrong.
type LegacyLookup interface {
	// LookupLegacy returns the account as it is known to the previous
	// application.
	LookupLegacy(username string) (model.User, bool, error)
}

// LegacyPasswordSetter changes the password of an account of the previous
// application.
//
// It is implemented by the account store of the previous application, so a
// user keeps one password for both applications.
type LegacyPasswordSetter interface {
	// SetLegacyPassword stores a new password for the account.
	//
	// It reports whether the account is known to the store.
	SetLegacyPassword(username, password string) (bool, error)
}

// Registrar creates a new account in an external account store, such as the
// USERS sheet of the previous application.
type Registrar interface {
	// RegisterUser stores a new member account.
	RegisterUser(reg model.Registration) error
}

// LegacyProfiler reports the profile of an account of the previous
// application, without checking its password.
//
// It is consulted for an account that was migrated before the access levels
// were introduced, so that the level and the district are filled in.
type LegacyProfiler interface {
	// LegacyProfile returns the account as it is known to the previous
	// application.
	LegacyProfile(username string) (model.User, bool, error)
}

// LegacyStatus reports the profile and the active flag of an account of the
// previous application.
//
// It is optional. An account store that implements it lets this application
// keep the salted local hash of a password, while the access of the account
// and its active flag still come from the previous application.
type LegacyStatus interface {
	// LegacyAccountStatus returns the account as it is known to the previous
	// application, and whether that account may sign in.
	LegacyAccountStatus(username string) (model.User, bool, error)
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

	// profiler completes the profile of an account that was migrated before
	// the access levels were introduced.
	profiler LegacyProfiler

	// registrar stores new member accounts.
	registrar Registrar

	// sheetOnly means every sign in is verified against the account store of
	// the previous application, so a local copy of an account is never
	// authoritative.
	sheetOnly bool

	// path is the file the accounts are stored in. It is empty when the
	// manager must not write to disk.
	path string

	// secureCookies forces the Secure flag on the session cookie, whatever the
	// request looks like. It is used in production, where the server always
	// sits behind TLS.
	secureCookies bool

	// versions holds the session token of every account whose password was
	// changed. A session that carries an older token is refused, so a stolen
	// session stops working as soon as the password is replaced.
	versions map[string]string

	// versionsPath is the file the session tokens are stored in. It is empty
	// when they must not be written to disk.
	versionsPath string
}

// NewManager builds a Manager from a list of users and a signing secret.
func NewManager(users []User, secret []byte) *Manager {
	m := &Manager{
		users:    make(map[string]User, len(users)),
		secret:   secret,
		versions: map[string]string{},
	}
	for _, u := range users {
		m.users[strings.ToLower(u.Username)] = u
	}
	return m
}

// SetSessionVersionsPath sets the file the session tokens are stored in, so a
// password change keeps its effect over a restart.
//
// The tokens that are already in the file are loaded, so the sessions that
// were made invalid stay invalid.
func (m *Manager) SetSessionVersionsPath(path string) {
	m.mu.Lock()
	m.versionsPath = path
	m.mu.Unlock()

	if path == "" {
		return
	}
	stored := map[string]string{}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &stored)
	}
	m.mu.Lock()
	for key, value := range stored {
		m.versions[strings.ToLower(key)] = value
	}
	m.mu.Unlock()
}

// sessionVersion returns the token of an account.
//
// It is empty for an account whose password was never changed, so the
// sessions of that account keep working.
func (m *Manager) sessionVersion(username string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.versions[strings.ToLower(strings.TrimSpace(username))]
}

// bumpVersion replaces the token of an account, so every session that was
// made before the change stops working.
//
// The token is random, so it cannot be guessed, and it is stored on disk, so
// the sessions stay invalid after a restart.
func (m *Manager) bumpVersion(username string) error {
	key := strings.ToLower(strings.TrimSpace(username))
	if key == "" {
		return ErrAccountNotFound
	}
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return err
	}

	m.mu.Lock()
	m.versions[key] = hex.EncodeToString(token)
	path := m.versionsPath
	all := make(map[string]string, len(m.versions))
	for name, value := range m.versions {
		all[name] = value
	}
	m.mu.Unlock()

	if path == "" {
		return nil
	}
	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}

// SetLegacyAuthenticator sets the source that verifies accounts of the
// previous application.
func (m *Manager) SetLegacyAuthenticator(a LegacyAuthenticator) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.legacy = a
	if p, ok := a.(LegacyProfiler); ok {
		m.profiler = p
	}
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
//
// The account is written to the account store of the previous application,
// so the member can sign in with it in both applications. It is also stored
// locally with a salted hash, so this application never has to check such a
// password against the unsalted hash of the previous application.
func (m *Manager) RegisterUser(reg model.Registration) error {
	m.mu.RLock()
	registrar := m.registrar
	m.mu.RUnlock()

	if registrar == nil {
		return ErrNoRegistrar
	}
	if err := registrar.RegisterUser(reg); err != nil {
		return err
	}
	return m.storeRegisteredUser(reg)
}

// SetSheetOnly makes the account store of the previous application the only
// authority for a sign in.
//
// A local copy of an account is then no longer used to verify a password, so a
// change of the password or of the role in that store takes effect at once.
func (m *Manager) SetSheetOnly(v bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sheetOnly = v
}

// SetUsersPath sets the file the accounts are stored in, so a migrated
// account can be written to it.
func (m *Manager) SetUsersPath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.path = path
}

// SetSecureCookies forces the Secure flag on the session cookie.
//
// Behind a reverse proxy that terminates TLS, the request itself is plain
// HTTP, so the flag cannot always be derived from the request. Setting this
// makes every session cookie secure anyway.
func (m *Manager) SetSecureCookies(v bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.secureCookies = v
}

// isSecure reports whether the session cookie must carry the Secure flag.
//
// A request that arrives over TLS is secure. A request that arrives over
// plain HTTP but carries the forwarding header of a proxy that terminated
// TLS is secure as well, so the cookie is still protected.
func (m *Manager) isSecure(r *http.Request) bool {
	m.mu.RLock()
	forced := m.secureCookies
	m.mu.RUnlock()
	if forced {
		return true
	}
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	proto := r.Header.Get("X-Forwarded-Proto")
	if i := strings.IndexByte(proto, ','); i >= 0 {
		proto = proto[:i]
	}
	return strings.EqualFold(strings.TrimSpace(proto), "https")
}

// VerifyLogin checks a password against the stored accounts.
//
// Local accounts are checked first. When the username is unknown locally,
// the accounts of the previous application are consulted, so a user can sign
// in with the credentials they already had.
//
// When the manager is set to sheet only, the accounts of the previous
// application are the authority: their password decides the sign in. A local
// account that the previous application does not know still signs in locally.
//
// The returned bool reports whether the login came from the previous
// application, which means the account can be migrated.
func (m *Manager) VerifyLogin(username, password string) (model.User, bool, error) {
	key := strings.ToLower(strings.TrimSpace(username))

	m.mu.RLock()
	u, ok := m.users[key]
	legacy := m.legacy
	sheetOnly := m.sheetOnly
	m.mu.RUnlock()

	if sheetOnly && legacy != nil {
		user, valid, err := legacy.VerifyLegacy(username, password)
		if err != nil {
			return model.User{}, false, err
		}
		if valid {
			return user, true, nil
		}
		// The previous application knows the account, so its password is
		// wrong or the account is not active. A local copy must not accept
		// it.
		if known, determined := legacyKnown(legacy, username); determined && known {
			return model.User{}, false, ErrInvalidCredentials
		}
	}

	if ok {
		user, err := m.verifyLocal(u, password)
		if err != nil {
			return model.User{}, false, err
		}
		return user, false, nil
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

// legacyKnown reports whether the account store of the previous application
// knows the account.
//
// The second bool reports whether the store could be asked. A store that
// cannot tell is treated as not knowing the account.
func legacyKnown(legacy LegacyAuthenticator, username string) (known, determined bool) {
	l, ok := legacy.(LegacyLookup)
	if !ok {
		return false, false
	}
	_, found, err := l.LookupLegacy(username)
	if err != nil {
		return false, false
	}
	return found, true
}

// isLegacyID reports whether the id belongs to an account of the previous
// application.
//
// Such an account carries a generated identifier, while a local account
// carries its own name as the id.
func isLegacyID(id string) bool {
	return strings.Contains(id, "-")
}

// completeProfile fills in the level and the district of an account that was
// migrated before the access levels were introduced.
//
// The account is stored again, so the levels survive the next sign in.
func (m *Manager) completeProfile(user model.User) (model.User, bool) {
	m.mu.RLock()
	profiler := m.profiler
	m.mu.RUnlock()
	if profiler == nil {
		return user, false
	}
	profile, found, err := profiler.LegacyProfile(user.Username)
	if err != nil || !found {
		return user, false
	}
	user.Level = profile.Level
	if profile.Kecamatan != "" {
		user.Kecamatan = profile.Kecamatan
	}
	m.storeUser(user)
	return user, true
}

// storeUser writes the account to the local store.
func (m *Manager) storeUser(user model.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := strings.ToLower(strings.TrimSpace(user.Username))
	u, ok := m.users[key]
	if !ok {
		return
	}
	u.User = user
	m.users[key] = u
	if m.path == "" {
		return
	}
	all := make([]User, 0, len(m.users))
	for _, item := range m.users {
		all = append(all, item)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Username < all[j].Username })
	_ = saveUsers(m.path, all)
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

// ChangePassword replaces the password of an account.
//
// The current password must be given, so a stolen session alone cannot take
// the account over. The new password is written to the local account store
// and, when the account also exists in the account store of the previous
// application, there too. One password then works in both applications.
//
// Every account may change its own password; the access level does not matter.
func (m *Manager) ChangePassword(username, current, next string) error {
	key := strings.ToLower(strings.TrimSpace(username))
	if key == "" {
		return ErrAccountNotFound
	}

	m.mu.RLock()
	u, localExists := m.users[key]
	legacy := m.legacy
	m.mu.RUnlock()

	if !localExists && legacy == nil {
		return ErrAccountNotFound
	}

	// The current password is checked against every store that knows the
	// account, because the store that decides a sign in depends on the mode.
	verified := false
	if localExists {
		if err := m.verifyStored(u, current); err == nil {
			verified = true
		}
	}
	if !verified && legacy != nil {
		_, valid, err := legacy.VerifyLegacy(username, current)
		if err != nil {
			return err
		}
		verified = valid
	}
	if !verified {
		return ErrInvalidCredentials
	}

	if localExists {
		salt, hash, err := HashPassword(next)
		if err != nil {
			return err
		}
		m.mu.Lock()
		u.PasswordHash = hash
		u.Salt = salt
		m.users[key] = u
		m.mu.Unlock()
		if err := m.saveAll(); err != nil {
			return err
		}
	}

	// The account of the previous application keeps the same password, so the
	// user does not have to remember two of them.
	if setter, ok := legacy.(LegacyPasswordSetter); ok {
		found, err := setter.SetLegacyPassword(username, next)
		if err != nil {
			return err
		}
		if !found && !localExists {
			return ErrAccountNotFound
		}
	}

	// The sessions that were made before the change stop working, so a stolen
	// session cannot be used once the password was replaced.
	if err := m.bumpVersion(username); err != nil {
		return err
	}
	return nil
}

// saveAll writes the accounts to the local store.
func (m *Manager) saveAll() error {
	m.mu.RLock()
	path := m.path
	all := make([]User, 0, len(m.users))
	for _, item := range m.users {
		all = append(all, item)
	}
	m.mu.RUnlock()

	if path == "" {
		return nil
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Username < all[j].Username })
	return saveUsers(path, all)
}

// SetSession writes a signed session cookie for the user.
//
// The cookie carries the Secure flag whenever the request is secure, so it is
// never sent over a plain connection in production.
func (m *Manager) SetSession(w http.ResponseWriter, r *http.Request, user model.User) error {
	payload := struct {
		Username  string `json:"u"`
		Role      string `json:"r"`
		Level     int    `json:"l"`
		Kecamatan string `json:"k"`
		Expires   int64  `json:"e"`
		Version   string `json:"v,omitempty"`
	}{
		Username:  user.Username,
		Role:      user.Role,
		Level:     user.LevelOrDefault(),
		Kecamatan: user.Kecamatan,
		Expires:   time.Now().Add(sessionTTL).Unix(),
		Version:   m.sessionVersion(user.Username),
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
		Secure:   m.isSecure(r),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(sessionTTL),
	})
	return nil
}

// ClearSession expires the session cookie.
//
// The flags match the ones of the session cookie, so the browser removes it
// even when the connection is secure.
func (m *Manager) ClearSession(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   m.isSecure(r),
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
		Username  string `json:"u"`
		Role      string `json:"r"`
		Level     int    `json:"l"`
		Kecamatan string `json:"k"`
		Expires   int64  `json:"e"`
		Version   string `json:"v"`
	}
	if err := json.Unmarshal(b, &payload); err != nil {
		return model.User{}, false
	}
	if time.Now().Unix() > payload.Expires {
		return model.User{}, false
	}
	// A session that was made before a password change carries an older
	// token, so it is refused. The token of an account whose password was
	// never changed is empty, so those sessions keep working.
	if payload.Version != m.sessionVersion(payload.Username) {
		return model.User{}, false
	}
	m.mu.RLock()
	u, ok := m.users[strings.ToLower(payload.Username)]
	sheetOnly := m.sheetOnly
	m.mu.RUnlock()
	if !ok {
		// The account is not stored locally, which is the case when every
		// sign in is verified against the sheet of the previous application.
		// The session is signed, so its payload can be trusted.
		m.mu.RLock()
		legacy := m.legacy
		m.mu.RUnlock()
		if legacy == nil {
			return model.User{}, false
		}
		return model.User{
			ID:        payload.Username,
			Username:  payload.Username,
			Role:      payload.Role,
			Level:     payload.Level,
			Kecamatan: payload.Kecamatan,
		}, true
	}
	// The stored account is authoritative, but an account that was migrated
	// before the levels were introduced carries none. The signed session then
	// supplies them.
	if u.User.Level == 0 {
		u.User.Level = payload.Level
	}
	if u.User.Kecamatan == "" {
		u.User.Kecamatan = payload.Kecamatan
	}
	if sheetOnly {
		// The account store of the previous application is the authority, so
		// the values it gave at the sign in win over the local copy.
		u.User.Role = payload.Role
		u.User.Level = payload.Level
		u.User.Kecamatan = payload.Kecamatan
	}
	// An account of the previous application that was migrated before the
	// levels were introduced still needs them. They are taken from the
	// previous application and stored, so this happens only once.
	if u.User.Level == 0 && isLegacyID(u.User.ID) {
		if user, done := m.completeProfile(u.User); done {
			u.User = user
		}
	}
	return u.User, true
}

func (m *Manager) sign(value string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// ParseUsers reads the accounts of the given file.
//
// When the file does not exist yet, the first administrator of the
// installation is created from the given bootstrap. No account with a
// password that is written down in the source code is ever made, so an
// installation that is started without a password of its own gets a random
// one instead.
func ParseUsers(path string, boot Bootstrap) ([]User, error) {
	users, err := loadUsers(path)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return users, nil
	}

	admin, err := bootstrapAdmin(boot)
	if err != nil {
		return nil, err
	}
	if err := saveUsers(path, []User{admin}); err != nil {
		return nil, err
	}
	return []User{admin}, nil
}

// RandomSecret returns a random secret used to sign sessions.
func RandomSecret() ([]byte, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("auth: %w", err)
	}
	return b, nil
}

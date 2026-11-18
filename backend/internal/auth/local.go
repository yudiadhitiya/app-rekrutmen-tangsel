package auth

import (
	"log"
	"sort"
	"strings"

	"rekrutmen-tangsel/backend/internal/model"
)

// verifyLocal checks the password of an account that is stored locally.
//
// The password is checked against the salted hash of the local store, so the
// unsalted hash of the previous application is never used for an account that
// this application knows. An account of the previous application keeps its
// access and its active flag up to date from that application, and a password
// that was replaced there is taken over again, so the salted hash stays the
// one that decides a sign in.
func (m *Manager) verifyLocal(u User, password string) (model.User, error) {
	if err := m.verifyStored(u, password); err != nil {
		// The password may have been replaced in the previous application.
		user, valid, lerr := m.verifyLegacyFallback(u, password)
		if lerr != nil {
			log.Printf("gagal memeriksa akun %s di sumber warisan: %v", u.Username, lerr)
			return model.User{}, err
		}
		if !valid {
			return model.User{}, err
		}
		// The previous application accepted the password, so it is stored
		// again, this time with a salted hash.
		if rerr := m.rehashUser(u.Username, password); rerr != nil {
			log.Printf("gagal menyimpan hash akun %s: %v", u.Username, rerr)
		}
		return user, nil
	}

	user := u.User
	// An account that lives only in this application is decided by the local
	// store alone.
	if !isLegacyID(user.ID) {
		return user, nil
	}

	refreshed, active, known := m.refreshLegacyStatus(user)
	if known {
		if !active {
			// The account was deactivated in the previous application.
			return model.User{}, ErrInvalidCredentials
		}
		if refreshed != user {
			m.storeUser(refreshed)
		}
		return refreshed, nil
	}

	// An account that was migrated before the access levels were introduced
	// carries none. They are taken from the previous application, so an
	// administrator of a district does not get the access of a full
	// administrator.
	if user.Level == 0 {
		if completed, done := m.completeProfile(user); done {
			return completed, nil
		}
	}
	return user, nil
}

// verifyLegacyFallback checks the password against the previous application.
//
// It is used when the local hash does not match, because the password may
// have been replaced there. It is only used for an account of the previous
// application, so an account that lives only in this application never
// reaches the other store.
func (m *Manager) verifyLegacyFallback(u User, password string) (model.User, bool, error) {
	m.mu.RLock()
	legacy := m.legacy
	m.mu.RUnlock()
	if legacy == nil || !isLegacyID(u.ID) {
		return model.User{}, false, nil
	}
	user, valid, err := legacy.VerifyLegacy(u.Username, password)
	if err != nil {
		return model.User{}, false, err
	}
	return user, valid, nil
}

// refreshLegacyStatus updates the access of an account from the store of the
// previous application.
//
// The last bool reports whether that store knows the account, the second one
// whether the account may sign in. A store that cannot be reached is not a
// reason to refuse a sign in, so the local access is kept in that case.
func (m *Manager) refreshLegacyStatus(user model.User) (model.User, bool, bool) {
	m.mu.RLock()
	legacy := m.legacy
	m.mu.RUnlock()
	if legacy == nil {
		return user, true, false
	}
	status, ok := legacy.(LegacyStatus)
	if !ok {
		return user, true, false
	}
	profile, active, err := status.LegacyAccountStatus(user.Username)
	if err != nil {
		log.Printf("gagal membaca status akun %s dari sumber warisan: %v", user.Username, err)
		return user, true, false
	}
	if profile.Username == "" {
		// The store does not know the account, so the local copy decides.
		return user, true, false
	}
	if !active {
		return profile, false, true
	}
	refreshed := user
	if profile.Level > 0 {
		refreshed.Level = profile.Level
	}
	if profile.Kecamatan != "" {
		refreshed.Kecamatan = profile.Kecamatan
	}
	if profile.Nama != "" {
		refreshed.Nama = profile.Nama
	}
	if profile.Role != "" {
		refreshed.Role = profile.Role
	}
	return refreshed, true, true
}

// rehashUser replaces the salted hash of an account.
//
// It is used when the password was accepted by the store of the previous
// application, so the local copy holds a salted hash of it too.
func (m *Manager) rehashUser(username, password string) error {
	salt, hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	key := strings.ToLower(strings.TrimSpace(username))

	m.mu.Lock()
	u, ok := m.users[key]
	if !ok {
		m.mu.Unlock()
		return nil
	}
	u.PasswordHash = hash
	u.Salt = salt
	m.users[key] = u
	all := make([]User, 0, len(m.users))
	for _, item := range m.users {
		all = append(all, item)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Username < all[j].Username })
	m.mu.Unlock()

	if m.path == "" {
		return nil
	}
	return saveUsers(m.path, all)
}

// storeRegisteredUser stores a new member account in the local store.
//
// The account is stored with a salted hash, so its password is never checked
// against the unsalted hash of the previous application. Its access is left
// empty on purpose: it is taken from the previous application at the first
// sign in, so a change of the role there takes effect at once.
func (m *Manager) storeRegisteredUser(reg model.Registration) error {
	salt, hash, err := HashPassword(reg.Password)
	if err != nil {
		return err
	}
	username := strings.TrimSpace(reg.Username)
	user := model.User{
		ID:        model.LegacyUserID(username),
		Username:  username,
		Nama:      strings.TrimSpace(reg.Nama),
		Role:      model.RolePetugas,
		Kecamatan: strings.TrimSpace(reg.Kecamatan),
	}

	m.mu.Lock()
	key := strings.ToLower(username)
	if _, exists := m.users[key]; exists {
		m.mu.Unlock()
		return nil
	}
	m.users[key] = User{User: user, PasswordHash: hash, Salt: salt}
	all := make([]User, 0, len(m.users))
	for _, item := range m.users {
		all = append(all, item)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Username < all[j].Username })
	m.mu.Unlock()

	if m.path == "" {
		return nil
	}
	return saveUsers(m.path, all)
}

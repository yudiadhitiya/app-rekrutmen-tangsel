package auth

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

// KnownDefaultPasswords are the passwords that the earlier versions of this
// application gave to the accounts they created by themselves.
//
// They are written down in the source code and in the documentation, so an
// installation that still has one of them is open to everybody who read that
// code.
var KnownDefaultPasswords = []string{
	"admin123",
	"petugas123",
}

// RotatedAccount is an account whose default password was replaced by a
// random one.
type RotatedAccount struct {
	// Username is the name of the account.
	Username string
	// Password is the new password of the account.
	Password string
}

// RotateDefaultPasswords replaces the password of every local account that
// still has one of the known default passwords.
//
// The new password is random, so the account is no longer known to whoever
// read the source code. The caller writes the new passwords to a file that
// only the owner of the service may read, so the operator can hand them out.
//
// Nothing is written when no account has a default password, so the function
// may be run on every start.
//
// The accounts of the previous application are left alone: their password is
// not kept in this file, so there is nothing to rotate here.
func RotateDefaultPasswords(path string) ([]RotatedAccount, error) {
	users, err := loadUsers(path)
	if err != nil {
		return nil, err
	}

	var rotated []RotatedAccount
	for i := range users {
		if isLegacyID(users[i].ID) {
			continue
		}
		if !hasKnownDefaultPassword(users[i]) {
			continue
		}
		password, err := RandomPassword(16)
		if err != nil {
			return nil, err
		}
		salt, hash, err := HashPassword(password)
		if err != nil {
			return nil, err
		}
		users[i].Salt = salt
		users[i].PasswordHash = hash
		rotated = append(rotated, RotatedAccount{
			Username: users[i].Username,
			Password: password,
		})
	}
	if len(rotated) == 0 {
		return nil, nil
	}
	if err := saveUsers(path, users); err != nil {
		return nil, err
	}
	return rotated, nil
}

// hasKnownDefaultPassword reports whether the account still has one of the
// known default passwords.
//
// The stored hash is compared with the hash of every known default password,
// so the check works for the accounts of every version that used them.
func hasKnownDefaultPassword(u User) bool {
	if u.Salt == "" || u.PasswordHash == "" {
		return false
	}
	salt, err := hex.DecodeString(u.Salt)
	if err != nil {
		return false
	}
	want, err := hex.DecodeString(u.PasswordHash)
	if err != nil {
		return false
	}
	for _, candidate := range KnownDefaultPasswords {
		key, err := pbkdf2.Key(sha256.New, candidate, salt, pbkdf2Iters, pbkdf2KeyLen)
		if err != nil {
			continue
		}
		if subtle.ConstantTimeCompare(key, want) == 1 {
			return true
		}
	}
	return false
}

package api

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// Limits of a sign in attempt.
const (
	// DefaultLoginMaxAttempts is the number of failed sign ins one account
	// or one address may make within loginWindow.
	DefaultLoginMaxAttempts = 10

	// loginWindow is the period the failed attempts are counted in.
	loginWindow = 15 * time.Minute

	// loginAddressFactor is how many times the limit of an account applies to
	// an address. An address is shared by many people, for example in an
	// office, so the limit of an address must be higher. Without this, the
	// failures of one account would lock out everybody behind the same
	// address.
	loginAddressFactor = 5
)

// loginLimiter counts the failed sign ins, so the login cannot be used to
// guess a password by brute force.
//
// Two counters are kept for every attempt: one for the account, so one
// account cannot be attacked from many addresses, and one for the address, so
// one address cannot attack many accounts. The limit of the address is a
// multiple of the limit of the account, because an address is shared by many
// users. A successful sign in clears the counter of the account.
type loginLimiter struct {
	mu      sync.Mutex
	entries map[string][]time.Time
	max     int

	// store keeps the counters on the disk, so a restart of the server does
	// not hand a new set of attempts to whoever tries to guess a password.
	store *limiterStore
}

// newLoginLimiter builds a limiter that allows a given number of failed
// attempts. A number of zero means the default.
//
// The store may be empty, which keeps the counters in the memory only.
func newLoginLimiter(maxAttempts int, store *limiterStore) *loginLimiter {
	if maxAttempts <= 0 {
		maxAttempts = DefaultLoginMaxAttempts
	}
	l := &loginLimiter{entries: map[string][]time.Time{}, max: maxAttempts, store: store}
	if store != nil {
		for key, times := range store.entries() {
			l.entries[key] = times
		}
	}
	return l
}

// keys returns the counters that belong to an attempt.
func (l *loginLimiter) keys(username, address string) []string {
	keys := make([]string, 0, 2)
	if username != "" {
		keys = append(keys, "u:"+strings.ToLower(username))
	}
	if address != "" {
		keys = append(keys, "ip:"+address)
	}
	return keys
}

// limitFor returns the number of failures a key tolerates.
//
// The counter of an account uses the limit of the limiter. The counter of an
// address uses a multiple of it, because one address serves many accounts.
func (l *loginLimiter) limitFor(key string) int {
	if strings.HasPrefix(key, "ip:") {
		return l.max * loginAddressFactor
	}
	return l.max
}

// allow reports whether another attempt is permitted.
//
// A rejected attempt is not counted, so a locked account stays locked until
// the window has passed.
func (l *loginLimiter) allow(username, address string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for _, key := range l.keys(username, address) {
		if l.count(key, now) >= l.limitFor(key) {
			return false
		}
	}
	return true
}

// fail records a failed attempt.
func (l *loginLimiter) fail(username, address string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for _, key := range l.keys(username, address) {
		kept := l.prune(key, now)
		l.entries[key] = append(kept, now)
		if l.store != nil {
			l.store.set(key, l.entries[key])
		}
	}
}

// reset clears the counter of an account, so a successful sign in starts
// from zero.
//
// The counter of the address is kept, because the address is shared by many
// accounts and a single successful sign in must not wipe the protection
// against a brute force that runs from that address.
func (l *loginLimiter) reset(username, address string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, key := range l.keys(username, address) {
		if strings.HasPrefix(key, "u:") {
			delete(l.entries, key)
			if l.store != nil {
				l.store.remove(key)
			}
		}
	}
}

// count returns the number of attempts that are still inside the window.
func (l *loginLimiter) count(key string, now time.Time) int {
	return len(l.prune(key, now))
}

// prune drops the attempts that fell out of the window.
func (l *loginLimiter) prune(key string, now time.Time) []time.Time {
	kept := l.entries[key][:0]
	for _, t := range l.entries[key] {
		if now.Sub(t) < loginWindow {
			kept = append(kept, t)
		}
	}
	l.entries[key] = kept
	return kept
}

// loginLocked reports whether a sign in is refused because too many attempts
// failed.
//
// It writes the answer and returns true when the attempt is refused.
func (s *Server) loginLocked(w http.ResponseWriter, r *http.Request, username string) bool {
	address := s.clientAddress(r)
	if s.loginLimiter.allow(username, address) {
		return false
	}
	w.Header().Set("Retry-After", "900")
	writeError(w, http.StatusTooManyRequests,
		"terlalu banyak percobaan login. Coba lagi dalam 15 menit.")
	return true
}

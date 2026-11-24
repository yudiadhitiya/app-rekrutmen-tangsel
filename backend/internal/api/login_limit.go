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

	// loginWindow is the period the failed attempts of an address are
	// counted in.
	loginWindow = 15 * time.Minute

	// loginAddressFactor is how many times the limit of an account applies to
	// an address. An address is shared by many people, for example in an
	// office, so the limit of an address must be higher. Without this, the
	// failures of one account would lock out everybody behind the same
	// address.
	loginAddressFactor = 5

	// loginAccountFactor is how many times the limit of an account applies to
	// the counter of the account itself.
	//
	// The counter of an account is weaker than the counter of an address, so
	// it may be reached a few times more before it refuses.
	loginAccountFactor = 2

	// loginAccountWindow is the period the failed attempts of an account are
	// counted in.
	//
	// It is longer than the window of an address, so one unlucky moment does
	// not follow the account for long, while a patient attacker cannot wait
	// the window out in a few minutes.
	loginAccountWindow = 3 * time.Hour

	// knownGoodWindow is how long an account that signed in successfully
	// from an address keeps its exemption.
	//
	// The limit of an address is a limit of an address, not of the people
	// behind it. An office behind one address would be locked out by one
	// attacker. So the limit of an address only refuses the accounts that
	// never signed in there, and an account that signed in successfully from
	// there keeps working.
	knownGoodWindow = 30 * 24 * time.Hour

	// goodKeyPrefix marks the entries that hold the last successful sign in
	// of an account on an address.
	goodKeyPrefix = "g:"

	// keySeparator joins the parts of a key, so two parts cannot be mixed up.
	keySeparator = "\x00"
)

// lockedTryEvery is how often the password of an account under attack may be
// checked.
//
// A locked account is being guessed at, so the server must not spend its time
// on every guess that comes in. One try every few seconds is enough for a
// user who knows the password and useless for an attacker.
//
// It is a variable, so a test does not have to wait for it.
var lockedTryEvery = 10 * time.Second

// loginLimiter counts the failed sign ins, so the login cannot be used to
// guess a password by brute force.
//
// Two counters are kept for every attempt: one for the account, so one
// account cannot be attacked from many addresses, and one for the address, so
// one address cannot attack many accounts.
//
// The limit of the address refuses before the password is checked, so a flood
// from one address costs the server nothing. The limit of an account refuses
// only the attempts with a wrong password, so an attacker cannot lock an
// account out: a password that is right always gets in.
type loginLimiter struct {
	mu      sync.Mutex
	entries map[string][]time.Time
	max     int

	// good holds the last successful sign in of every account on every
	// address.
	//
	// An address is shared by many people, for example in an office behind
	// one connection. The limit of an address must therefore not lock out
	// the people who live there: an account that signed in successfully from
	// there keeps working, while the accounts that never signed in are still
	// refused.
	good map[string]time.Time

	// tries holds the moment the password of an account under attack was
	// checked the last time, so the attempts can be slowed down.
	tries map[string]time.Time

	// bad holds the addresses that already failed for an account.
	//
	// An address that never failed gets its chance, so a user who signs in
	// from a new connection is never kept waiting by an attacker.
	bad map[string]time.Time

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
	l := &loginLimiter{
		entries: map[string][]time.Time{},
		good:    map[string]time.Time{},
		tries:   map[string]time.Time{},
		bad:     map[string]time.Time{},
		max:     maxAttempts,
		store:   store,
	}
	if store != nil {
		for key, times := range store.entries() {
			if strings.HasPrefix(key, goodKeyPrefix) {
				if len(times) > 0 {
					l.good[strings.TrimPrefix(key, goodKeyPrefix)] = times[len(times)-1]
				}
				continue
			}
			l.entries[key] = times
		}
	}
	return l
}

// accountKey returns the key of the counter of an account.
func accountKey(username string) string {
	name := strings.ToLower(strings.TrimSpace(username))
	if name == "" {
		return ""
	}
	return "u:" + name
}

// addressKey returns the key of the counter of an address.
func addressKey(address string) string {
	if address == "" {
		return ""
	}
	return "ip:" + address
}

// goodKey returns the key of the last successful sign in of an account on an
// address.
func goodKey(username, address string) string {
	name := strings.ToLower(strings.TrimSpace(username))
	if name == "" {
		return ""
	}
	return name + keySeparator + address
}

// badKey returns the key of the addresses that already failed for an
// account.
func badKey(username, address string) string {
	name := strings.ToLower(strings.TrimSpace(username))
	if name == "" {
		return ""
	}
	return name + keySeparator + address
}

// keys returns the counters that belong to an attempt.
func (l *loginLimiter) keys(username, address string) []string {
	keys := make([]string, 0, 2)
	if key := accountKey(username); key != "" {
		keys = append(keys, key)
	}
	if key := addressKey(address); key != "" {
		keys = append(keys, key)
	}
	return keys
}

// windowFor returns the period the counter of a key is counted in.
//
// The failures of an account are counted longer than the failures of an
// address, because the counter of an account is the weaker of the two.
func windowFor(key string) time.Duration {
	if strings.HasPrefix(key, "u:") {
		return loginAccountWindow
	}
	return loginWindow
}

// limitFor returns the number of failures a key tolerates.
//
// The counter of an account uses the limit of the limiter times a small
// factor, because it is weaker than the counter of an address. The counter of
// an address uses a larger multiple of it, because one address serves many
// accounts.
func (l *loginLimiter) limitFor(key string) int {
	switch {
	case strings.HasPrefix(key, "u:"):
		return l.max * loginAccountFactor
	case strings.HasPrefix(key, "ip:"):
		return l.max * loginAddressFactor
	}
	return l.max
}

// allowAddress reports whether an attempt from this address may be made.
//
// It is the only limit that refuses before the password is checked, so a
// flood from one address costs the server nothing. An account that signed in
// successfully from this address is not refused, because one address serves
// many people and one attacker must not be able to lock them all out.
func (l *loginLimiter) allowAddress(username, address string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := addressKey(address)
	if key == "" {
		return true
	}
	now := time.Now()
	if l.count(key, now) < l.limitFor(key) {
		return true
	}
	return l.isKnownGood(username, address, now)
}

// locked reports whether the account is being guessed at, so its wrong
// passwords are refused.
//
// A locked account is not cut off: a password that is right still gets in.
func (l *loginLimiter) locked(username string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.lockedAt(username, time.Now())
}

// lockedAt reports whether the counter of the account is over its limit.
//
// It expects the lock of the limiter to be held.
func (l *loginLimiter) lockedAt(username string, now time.Time) bool {
	key := accountKey(username)
	if key == "" {
		return false
	}
	return l.count(key, now) >= l.limitFor(key)
}

// mayTry reports whether the password of an account may be checked now.
//
// A locked account is being guessed at, so the server must not spend its
// time on every guess that comes in. One try every few seconds is enough for
// a user who knows the password and useless for an attacker.
//
// An address that the account used successfully is not slowed down, so the
// real user is never kept waiting. An address that never failed for the
// account gets its chance as well, so a user who signs in from a new
// connection is not punished for the attack of somebody else.
func (l *loginLimiter) mayTry(username, address string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	name := strings.ToLower(strings.TrimSpace(username))
	if name == "" {
		return true
	}
	now := time.Now()
	if !l.lockedAt(name, now) {
		return true
	}
	if l.isKnownGood(name, address, now) {
		return true
	}

	// The address never failed for this account, so it gets its chance.
	bad := badKey(name, address)
	if last, ok := l.bad[bad]; !ok || now.Sub(last) > loginAccountWindow {
		l.bad[bad] = now
		return true
	}

	// The address already failed here, so it shares the one try every few
	// seconds that a locked account gets.
	if last, ok := l.tries[name]; ok && now.Sub(last) < lockedTryEvery {
		return false
	}
	l.tries[name] = now
	return true
}

// isKnownGood reports whether the account signed in successfully from this
// address recently.
//
// It expects the lock of the limiter to be held.
func (l *loginLimiter) isKnownGood(username, address string, now time.Time) bool {
	key := goodKey(username, address)
	if key == "" {
		return false
	}
	last, ok := l.good[key]
	if !ok {
		return false
	}
	return now.Sub(last) < knownGoodWindow
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

	// The address failed for this account, so it does not get a free try
	// while the account is being guessed at.
	if key := badKey(username, address); key != "" {
		l.bad[key] = now
		l.pruneBad(now)
	}
}

// pruneBad drops the addresses that failed long ago, so the memory of the
// limiter stays small even when it is attacked from many addresses.
func (l *loginLimiter) pruneBad(now time.Time) {
	if len(l.bad) < 4096 {
		return
	}
	for key, last := range l.bad {
		if now.Sub(last) > loginAccountWindow {
			delete(l.bad, key)
		}
	}
}

// reset clears the counter of an account, so a successful sign in starts
// from zero, and remembers that the account signed in successfully.
//
// The counter of the address is kept, because the address is shared by many
// accounts and a single successful sign in must not wipe the protection
// against a brute force that runs from that address.
func (l *loginLimiter) reset(username, address string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for _, key := range l.keys(username, address) {
		if strings.HasPrefix(key, "u:") {
			delete(l.entries, key)
			delete(l.tries, strings.ToLower(strings.TrimSpace(username)))
			if l.store != nil {
				l.store.remove(key)
			}
		}
	}

	// The account is not under attack any more, so the addresses that failed
	// for it do not have to be remembered.
	prefix := strings.ToLower(strings.TrimSpace(username)) + keySeparator
	for key := range l.bad {
		if strings.HasPrefix(key, prefix) {
			delete(l.bad, key)
		}
	}

	// The account signed in successfully on this address, so it keeps
	// working there even when the address runs into its limit.
	key := goodKey(username, address)
	if key == "" {
		return
	}
	l.good[key] = now
	if l.store != nil {
		l.store.set(goodKeyPrefix+key, []time.Time{now})
		l.pruneGood(now)
	}
}

// pruneGood drops the successful sign ins that are too old to matter, so
// the memory of the limiter and the file of the counters stay small.
func (l *loginLimiter) pruneGood(now time.Time) {
	if len(l.good) < 1024 {
		return
	}
	for key, last := range l.good {
		if now.Sub(last) > knownGoodWindow {
			delete(l.good, key)
			if l.store != nil {
				l.store.remove(goodKeyPrefix + key)
			}
		}
	}
}

// lockedFor returns how long the account stays locked.
//
// The counter of an account is counted over a long window, so the lock lasts
// until the oldest attempts fall out of it.
//
// It expects the lock of the limiter to be held.
func (l *loginLimiter) lockedFor(username string) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := accountKey(username)
	if key == "" {
		return 0
	}
	now := time.Now()
	times := l.prune(key, now)
	limit := l.limitFor(key)
	if len(times) < limit {
		return 0
	}
	// The attempts are kept in the order they happened, so the one that
	// must fall out first is the oldest one that is not needed to keep the
	// lock.
	oldest := times[len(times)-limit]
	wait := oldest.Add(windowFor(key)).Sub(now)
	if wait < 0 {
		return 0
	}
	return wait
}

// count returns the number of attempts that are still inside the window.
func (l *loginLimiter) count(key string, now time.Time) int {
	return len(l.prune(key, now))
}

// prune drops the attempts that fell out of the window.
func (l *loginLimiter) prune(key string, now time.Time) []time.Time {
	kept := l.entries[key][:0]
	for _, t := range l.entries[key] {
		if now.Sub(t) < windowFor(key) {
			kept = append(kept, t)
		}
	}
	l.entries[key] = kept
	return kept
}

// loginAddressLocked reports whether a sign in is refused because the address
// made too many failed attempts.
//
// It writes the answer and returns true when the attempt is refused.
func (s *Server) loginAddressLocked(w http.ResponseWriter, r *http.Request, username string) bool {
	address := s.clientAddress(r)
	if s.loginLimiter.allowAddress(username, address) {
		return false
	}
	w.Header().Set("Retry-After", "900")
	writeError(w, http.StatusTooManyRequests,
		"terlalu banyak percobaan login. Coba lagi dalam 15 menit.")
	return true
}

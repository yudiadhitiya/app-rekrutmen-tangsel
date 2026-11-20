package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// limiterStore keeps the counters of the limiters on the disk.
//
// A restart of the server must not hand an attacker a fresh set of attempts,
// so the counters of the failed sign ins, of the registrations and of the
// uploads survive it.
//
// The file holds nothing but the time of the attempts, so it says nothing
// about the accounts themselves.
type limiterStore struct {
	mu   sync.Mutex
	path string
	data map[string][]int64
}

// newLimiterStore builds a store that keeps its counters in the given file.
//
// An empty path keeps the counters in the memory only, which is what the
// tests use.
func newLimiterStore(path string) *limiterStore {
	s := &limiterStore{path: path, data: map[string][]int64{}}
	if path == "" {
		return s
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	_ = json.Unmarshal(b, &s.data)
	return s
}

// entries returns every counter of the store.
//
// It is used when the server starts, so the attempts that were made before the
// restart are counted as well.
func (s *limiterStore) entries() map[string][]time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string][]time.Time, len(s.data))
	for key, stamps := range s.data {
		times := make([]time.Time, 0, len(stamps))
		for _, v := range stamps {
			times = append(times, time.Unix(v, 0))
		}
		out[key] = times
	}
	return out
}

// list returns the times of the attempts of a key.
func (s *limiterStore) list(key string) []time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	stamps := s.data[key]
	out := make([]time.Time, 0, len(stamps))
	for _, v := range stamps {
		out = append(out, time.Unix(v, 0))
	}
	return out
}

// set stores the times of the attempts of a key.
func (s *limiterStore) set(key string, times []time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	stamps := make([]int64, 0, len(times))
	for _, t := range times {
		stamps = append(stamps, t.Unix())
	}
	s.data[key] = stamps
	s.save()
}

// remove drops the counter of a key.
func (s *limiterStore) remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; !ok {
		return
	}
	delete(s.data, key)
	s.save()
}

// save writes the counters to the disk.
//
// The file is written beside the real one and then moved, so a crash never
// leaves a half written file behind.
func (s *limiterStore) save() {
	if s.path == "" {
		return
	}
	b, err := json.Marshal(s.data)
	if err != nil {
		return
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return
	}
	_ = os.Rename(tmp, s.path)
}

// limiterStorePath returns the file the counters are kept in.
//
// The file lives beside the accounts, in the directory that only the owner of
// the service may read. Every limiter has its own file, so the counters of
// one never mix with the counters of another.
func limiterStorePath(dataDir, name string) string {
	if dataDir == "" {
		return ""
	}
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return ""
	}
	return filepath.Join(dataDir, name)
}

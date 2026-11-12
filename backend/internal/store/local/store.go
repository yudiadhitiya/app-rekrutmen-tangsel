// Package local implements a file-backed Store. It is the default provider so
// the app runs out of the box, and serves as the reference implementation of
// the Store interface.
package local

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/store"
)

type fileState struct {
	Rekrutmen []model.Rekrutmen `json:"rekrutmen"`
	Aktivitas []model.Aktivitas `json:"aktivitas"`
}

// Store persists records in a single JSON file, guarded by a mutex.
type Store struct {
	mu    sync.RWMutex
	path  string
	state fileState
}

// New opens (or creates) the JSON file at path.
func New(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		s.state = fileState{}
		return s.save()
	}
	if err != nil {
		return err
	}
	if len(b) == 0 {
		s.state = fileState{}
		return nil
	}
	if err := json.Unmarshal(b, &s.state); err != nil {
		return err
	}
	// Migrate records still holding the retired Metode spellings.
	if s.migrateMetode() {
		return s.save()
	}
	return nil
}

// migrateMetode folds the retired Metode values into "Program Antara" and
// reports whether anything changed.
func (s *Store) migrateMetode() bool {
	changed := false
	for i := range s.state.Rekrutmen {
		if replacement, ok := model.MetodeLama[s.state.Rekrutmen[i].Metode]; ok {
			s.state.Rekrutmen[i].Metode = replacement
			changed = true
		}
	}
	return changed
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s.state, "", "  ")
	if err != nil {
		return err
	}
	// Write to a temp file then rename so a crash cannot truncate the store.
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Name implements store.Store.
func (s *Store) Name() string { return "local" }

// Ping implements store.Store.
func (s *Store) Ping() error { return nil }

func (s *Store) ListRekrutmen(f model.Filter) ([]model.Rekrutmen, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(f.Q))
	var out []model.Rekrutmen
	for _, r := range s.state.Rekrutmen {
		if q != "" && !strings.Contains(strings.ToLower(r.Nama), q) &&
			!strings.Contains(r.NoWA, q) {
			continue
		}
		if !model.MatchScope(f.Scope, r) {
			continue
		}
		if f.Kecamatan != "" && r.Kecamatan != f.Kecamatan {
			continue
		}
		if f.Kelurahan != "" && r.Kelurahan != f.Kelurahan {
			continue
		}
		if f.Metode != "" && !model.MetodeMatches(r.Metode, f.Metode) {
			continue
		}
		if f.Status != "" && r.Status != f.Status {
			continue
		}
		out = append(out, r)
	}

	// Newest first.
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})

	total := len(out)
	limit := f.Limit
	if limit <= 0 {
		limit = 25
	}
	page := f.Page
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	if start >= total {
		return []model.Rekrutmen{}, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}
	return out[start:end], total, nil
}

func (s *Store) GetRekrutmen(id string) (*model.Rekrutmen, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.state.Rekrutmen {
		if s.state.Rekrutmen[i].ID == id {
			r := s.state.Rekrutmen[i]
			return &r, nil
		}
	}
	return nil, store.ErrNotFound
}

func (s *Store) FindRekrutmenByWA(noWA string) (*model.Rekrutmen, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.state.Rekrutmen {
		if s.state.Rekrutmen[i].NoWA == noWA {
			r := s.state.Rekrutmen[i]
			return &r, nil
		}
	}
	return nil, store.ErrNotFound
}

func (s *Store) CreateRekrutmen(r *model.Rekrutmen) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.state.Rekrutmen {
		if existing.NoWA == r.NoWA {
			return store.ErrDuplicateWA
		}
	}
	r.Metode = model.NormalizeMetode(r.Metode)
	s.state.Rekrutmen = append(s.state.Rekrutmen, *r)
	return s.save()
}

func (s *Store) UpdateRekrutmen(r *model.Rekrutmen) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existing := range s.state.Rekrutmen {
		if existing.NoWA == r.NoWA && existing.ID != r.ID {
			return store.ErrDuplicateWA
		}
	}
	r.Metode = model.NormalizeMetode(r.Metode)
	for i := range s.state.Rekrutmen {
		if s.state.Rekrutmen[i].ID == r.ID {
			s.state.Rekrutmen[i] = *r
			return s.save()
		}
	}
	return store.ErrNotFound
}

func (s *Store) DeleteRekrutmen(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.state.Rekrutmen {
		if s.state.Rekrutmen[i].ID == id {
			s.state.Rekrutmen = append(s.state.Rekrutmen[:i], s.state.Rekrutmen[i+1:]...)
			return s.save()
		}
	}
	return store.ErrNotFound
}

func (s *Store) ListAktivitas() ([]model.Aktivitas, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]model.Aktivitas, len(s.state.Aktivitas))
	copy(out, s.state.Aktivitas)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Nama < out[j].Nama
	})
	return out, nil
}

func (s *Store) CreateAktivitas(a *model.Aktivitas) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.state.Aktivitas {
		if e.Nama == a.Nama {
			return store.ErrDuplicateName
		}
	}
	s.state.Aktivitas = append(s.state.Aktivitas, *a)
	return s.save()
}

// UpdateAktivitas renames an activity. Records referencing the old name are
// rewritten so a rename never orphans existing rekrutmen data.
func (s *Store) UpdateAktivitas(a *model.Aktivitas) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	oldName := ""
	for i := range s.state.Aktivitas {
		if s.state.Aktivitas[i].ID == a.ID {
			oldName = s.state.Aktivitas[i].Nama
			s.state.Aktivitas[i] = *a
			break
		}
	}
	if oldName == "" {
		return store.ErrNotFound
	}
	if oldName != a.Nama {
		for i := range s.state.Rekrutmen {
			if s.state.Rekrutmen[i].Aktivitas == oldName {
				s.state.Rekrutmen[i].Aktivitas = a.Nama
			}
		}
	}
	return s.save()
}

func (s *Store) DeleteAktivitas(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	name := ""
	for _, a := range s.state.Aktivitas {
		if a.ID == id {
			name = a.Nama
			break
		}
	}
	if name == "" {
		return store.ErrNotFound
	}
	for _, r := range s.state.Rekrutmen {
		if r.Aktivitas == name {
			return store.ErrAktivitasUsed
		}
	}
	for i := range s.state.Aktivitas {
		if s.state.Aktivitas[i].ID == id {
			s.state.Aktivitas = append(s.state.Aktivitas[:i], s.state.Aktivitas[i+1:]...)
			return s.save()
		}
	}
	return store.ErrNotFound
}

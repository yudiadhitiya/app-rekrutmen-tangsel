// Package settings holds the runtime application configuration, including the
// active storage provider chosen by an admin.
package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Provider identifiers.
const (
	ProviderLocal  = "local"
	ProviderSupa   = "supabase"
	ProviderSheets = "google_sheets"
)

// Providers lists the storage providers an admin may choose from.
var Providers = []string{ProviderLocal, ProviderSupa, ProviderSheets}

// Settings is the persisted, admin-editable configuration.
type Settings struct {
	Provider string `json:"provider"`
}

// Store persists settings to disk.
type Store struct {
	mu   sync.RWMutex
	path string
	cur  Settings
}

// NewStore loads settings from path, defaulting to the local provider.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, cur: Settings{Provider: ProviderLocal}}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, s.save()
	}
	if err != nil {
		return nil, err
	}
	if len(b) > 0 {
		if err := json.Unmarshal(b, &s.cur); err != nil {
			return nil, err
		}
	}
	if s.cur.Provider == "" {
		s.cur.Provider = ProviderLocal
	}
	return s, nil
}

// Get returns the current settings.
func (s *Store) Get() Settings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cur
}

// Set updates the settings and persists them.
func (s *Store) Set(next Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cur = next
	return s.save()
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s.cur, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o600)
}

// IsValidProvider reports whether a provider name is known.
func IsValidProvider(name string) bool {
	for _, p := range Providers {
		if p == name {
			return true
		}
	}
	return false
}

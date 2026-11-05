// Package store defines the storage abstraction that lets the app switch
// between Supabase and Google Sheets without touching the API layer.
package store

import (
	"errors"

	"rekrutmen-tangsel/backend/internal/model"
)

// Sentinel errors shared by every provider.
var (
	ErrNotFound      = errors.New("tidak ditemukan")
	ErrDuplicateWA   = errors.New("nomor whatsapp sudah terdaftar")
	ErrDuplicateName = errors.New("nama sudah digunakan")
	ErrAktivitasUsed = errors.New("aktivitas masih dipakai data rekrutmen")
)

// Store is the repository interface implemented by each storage provider.
//
// A provider owns two collections: the rekrutmen records and the master
// aktivitas list. Implementations are expected to be safe for concurrent use.
type Store interface {
	// Name returns the provider identifier used in settings.
	Name() string

	// Rekrutmen CRUD.
	ListRekrutmen(f model.Filter) ([]model.Rekrutmen, int, error)
	GetRekrutmen(id string) (*model.Rekrutmen, error)
	CreateRekrutmen(r *model.Rekrutmen) error
	UpdateRekrutmen(r *model.Rekrutmen) error
	DeleteRekrutmen(id string) error

	// FindRekrutmenByWA returns the record with the given WhatsApp number.
	// It returns ErrNotFound when no record matches.
	FindRekrutmenByWA(noWA string) (*model.Rekrutmen, error)

	// Aktivitas master CRUD.
	ListAktivitas() ([]model.Aktivitas, error)
	CreateAktivitas(a *model.Aktivitas) error
	UpdateAktivitas(a *model.Aktivitas) error
	DeleteAktivitas(id string) error

	// Ping verifies the provider is reachable and configured.
	Ping() error
}

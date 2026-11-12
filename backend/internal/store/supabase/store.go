// Package supabase implements Store against the Supabase PostgREST API.
//
// It talks to PostgREST directly over HTTPS rather than pulling in the
// Supabase SDK, so the backend stays on the Go standard library.
package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/store"
)

const (
	tableRekrutmen = "rekrutmen"
	tableAktivitas = "aktivitas"
)

// Store is a PostgREST-backed Store.
type Store struct {
	baseURL string
	apiKey  string
	token   string // optional user JWT; empty means the service/anon key is used
	http    *http.Client
}

// baseURL is a field so tests can point the store at a fake PostgREST server.

// New builds a Store from the project URL and API key.
func New(projectURL, apiKey, token string) *Store {
	return &Store{
		baseURL: strings.TrimRight(projectURL, "/") + "/rest/v1",
		apiKey:  apiKey,
		token:   token,
		http:    &http.Client{Timeout: 20 * time.Second},
	}
}

// Name implements store.Store.
func (s *Store) Name() string { return "supabase" }

// Ping verifies the project is reachable.
func (s *Store) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := s.newRequest(ctx, http.MethodGet, "/"+tableRekrutmen+"?limit=1", nil)
	if err != nil {
		return err
	}
	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("supabase: %s", s.errBody(resp))
	}
	return nil
}

func (s *Store) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+path, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Authorization", "Bearer "+s.authToken())
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (s *Store) authToken() string {
	if s.token != "" {
		return s.token
	}
	return s.apiKey
}

func (s *Store) errBody(resp *http.Response) string {
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	return strings.TrimSpace(string(b))
}

// do executes a request and decodes the JSON body into out (when out != nil).
func (s *Store) do(req *http.Request, out any) (*http.Response, error) {
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("supabase: %d %s", resp.StatusCode, s.errBody(resp))
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *Store) ListRekrutmen(f model.Filter) ([]model.Rekrutmen, int, error) {
	q := url.Values{}
	if s := strings.TrimSpace(f.Q); s != "" {
		// Case-insensitive match on name or WhatsApp number.
		q.Set("or", fmt.Sprintf("(nama.ilike.*%s*,no_wa.ilike.*%s*)", s, s))
	}
	if f.Scope.CreatedBy != "" {
		q.Set("created_by", "eq."+f.Scope.CreatedBy)
	}
	// The scope of the account and the district filter of the request both
	// narrow the district. A request for another district than the account
	// may reach yields no record at all.
	switch {
	case f.Scope.Kecamatan != "" && f.Kecamatan != "" && !strings.EqualFold(f.Scope.Kecamatan, f.Kecamatan):
		return []model.Rekrutmen{}, 0, nil
	case f.Scope.Kecamatan != "":
		q.Set("kecamatan", "eq."+f.Scope.Kecamatan)
	case f.Kecamatan != "":
		q.Set("kecamatan", "eq."+f.Kecamatan)
	}
	if f.Kelurahan != "" {
		q.Set("kelurahan", "eq."+f.Kelurahan)
	}
	if f.Metode != "" {
		// Use an "in" filter so the retired "Program"/"Antara" spellings still
		// match the current "Program Antara" option.
		values := model.MetodeMatchValues(f.Metode)
		if len(values) == 1 {
			q.Set("metode", "eq."+f.Metode)
		} else {
			quoted := make([]string, 0, len(values))
			for _, v := range values {
				quoted = append(quoted, `"`+v+`"`)
			}
			q.Set("metode", "in.("+strings.Join(quoted, ",")+")")
		}
	}
	if f.Status != "" {
		q.Set("status", "eq."+f.Status)
	}
	limit := f.Limit
	if limit <= 0 {
		limit = 25
	}
	page := f.Page
	if page <= 0 {
		page = 1
	}
	q.Set("limit", fmt.Sprint(limit))
	q.Set("offset", fmt.Sprint((page-1)*limit))
	q.Set("order", "created_at.desc")
	q.Set("select", "*")

	req, err := s.newRequest(context.Background(), http.MethodGet, "/"+tableRekrutmen+"?"+q.Encode(), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Prefer", "count=exact")

	var rows []model.Rekrutmen
	resp, err := s.do(req, &rows)
	if err != nil {
		return nil, 0, err
	}
	normalizeRows(rows)
	total := len(rows)
	if cr := resp.Header.Get("Content-Range"); cr != "" {
		// Format: "0-24/120"
		if idx := strings.LastIndex(cr, "/"); idx >= 0 {
			fmt.Sscanf(cr[idx+1:], "%d", &total)
		}
	}
	return rows, total, nil
}

func (s *Store) GetRekrutmen(id string) (*model.Rekrutmen, error) {
	req, err := s.newRequest(context.Background(), http.MethodGet,
		"/"+tableRekrutmen+"?id=eq."+url.QueryEscape(id)+"&select=*", nil)
	if err != nil {
		return nil, err
	}
	var rows []model.Rekrutmen
	if _, err := s.do(req, &rows); err != nil {
		return nil, err
	}
	normalizeRows(rows)
	if len(rows) == 0 {
		return nil, store.ErrNotFound
	}
	return &rows[0], nil
}

func (s *Store) FindRekrutmenByWA(noWA string) (*model.Rekrutmen, error) {
	req, err := s.newRequest(context.Background(), http.MethodGet,
		"/"+tableRekrutmen+"?no_wa=eq."+url.QueryEscape(noWA)+"&select=*", nil)
	if err != nil {
		return nil, err
	}
	var rows []model.Rekrutmen
	if _, err := s.do(req, &rows); err != nil {
		return nil, err
	}
	normalizeRows(rows)
	if len(rows) == 0 {
		return nil, store.ErrNotFound
	}
	return &rows[0], nil
}

func (s *Store) CreateRekrutmen(r *model.Rekrutmen) error {
	r.Metode = model.NormalizeMetode(r.Metode)
	req, err := s.newRequest(context.Background(), http.MethodPost, "/"+tableRekrutmen, r)
	if err != nil {
		return err
	}
	req.Header.Set("Prefer", "return=representation")
	var rows []model.Rekrutmen
	if _, err := s.do(req, &rows); err != nil {
		if isUniqueViolation(err) {
			return store.ErrDuplicateWA
		}
		return err
	}
	return nil
}

func (s *Store) UpdateRekrutmen(r *model.Rekrutmen) error {
	r.Metode = model.NormalizeMetode(r.Metode)
	req, err := s.newRequest(context.Background(), http.MethodPatch,
		"/"+tableRekrutmen+"?id=eq."+url.QueryEscape(r.ID), r)
	if err != nil {
		return err
	}
	req.Header.Set("Prefer", "return=representation")
	var rows []model.Rekrutmen
	if _, err := s.do(req, &rows); err != nil {
		if isUniqueViolation(err) {
			return store.ErrDuplicateWA
		}
		return err
	}
	if len(rows) == 0 {
		return store.ErrNotFound
	}
	return nil
}

func (s *Store) DeleteRekrutmen(id string) error {
	req, err := s.newRequest(context.Background(), http.MethodDelete,
		"/"+tableRekrutmen+"?id=eq."+url.QueryEscape(id), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Prefer", "return=representation")
	var rows []model.Rekrutmen
	if _, err := s.do(req, &rows); err != nil {
		return err
	}
	if len(rows) == 0 {
		return store.ErrNotFound
	}
	return nil
}

func (s *Store) ListAktivitas() ([]model.Aktivitas, error) {
	req, err := s.newRequest(context.Background(), http.MethodGet,
		"/"+tableAktivitas+"?select=*&order=nama.asc", nil)
	if err != nil {
		return nil, err
	}
	var rows []model.Aktivitas
	if _, err := s.do(req, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *Store) CreateAktivitas(a *model.Aktivitas) error {
	req, err := s.newRequest(context.Background(), http.MethodPost, "/"+tableAktivitas, a)
	if err != nil {
		return err
	}
	req.Header.Set("Prefer", "return=representation")
	var rows []model.Aktivitas
	if _, err := s.do(req, &rows); err != nil {
		if isUniqueViolation(err) {
			return store.ErrDuplicateName
		}
		return err
	}
	return nil
}

// UpdateAktivitas renames an activity and rewrites records that reference the
// old name, mirroring the behaviour of the other providers.
func (s *Store) UpdateAktivitas(a *model.Aktivitas) error {
	current, err := s.getAktivitas(a.ID)
	if err != nil {
		return err
	}
	req, err := s.newRequest(context.Background(), http.MethodPatch,
		"/"+tableAktivitas+"?id=eq."+url.QueryEscape(a.ID), a)
	if err != nil {
		return err
	}
	req.Header.Set("Prefer", "return=representation")
	var rows []model.Aktivitas
	if _, err := s.do(req, &rows); err != nil {
		if isUniqueViolation(err) {
			return store.ErrDuplicateName
		}
		return err
	}
	if current.Nama != a.Nama {
		if err := s.renameAktivitas(current.Nama, a.Nama); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) getAktivitas(id string) (*model.Aktivitas, error) {
	req, err := s.newRequest(context.Background(), http.MethodGet,
		"/"+tableAktivitas+"?id=eq."+url.QueryEscape(id)+"&select=*", nil)
	if err != nil {
		return nil, err
	}
	var rows []model.Aktivitas
	if _, err := s.do(req, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, store.ErrNotFound
	}
	return &rows[0], nil
}

// renameAktivitas rewrites the aktivitas column for every record using oldName.
func (s *Store) renameAktivitas(oldName, newName string) error {
	req, err := s.newRequest(context.Background(), http.MethodPatch,
		"/"+tableRekrutmen+"?aktivitas=eq."+url.QueryEscape(oldName),
		map[string]string{"aktivitas": newName})
	if err != nil {
		return err
	}
	_, err = s.do(req, nil)
	return err
}

func (s *Store) DeleteAktivitas(id string) error {
	a, err := s.getAktivitas(id)
	if err != nil {
		return err
	}
	// Refuse to delete while records still reference this activity.
	req, err := s.newRequest(context.Background(), http.MethodGet,
		"/"+tableRekrutmen+"?aktivitas=eq."+url.QueryEscape(a.Nama)+"&select=id&limit=1", nil)
	if err != nil {
		return err
	}
	var rows []model.Rekrutmen
	if _, err := s.do(req, &rows); err != nil {
		return err
	}
	if len(rows) > 0 {
		return store.ErrAktivitasUsed
	}

	req, err = s.newRequest(context.Background(), http.MethodDelete,
		"/"+tableAktivitas+"?id=eq."+url.QueryEscape(id), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Prefer", "return=representation")
	var del []model.Aktivitas
	if _, err := s.do(req, &del); err != nil {
		return err
	}
	if len(del) == 0 {
		return store.ErrNotFound
	}
	return nil
}

// isUniqueViolation detects a Postgres unique-constraint error reported by
// PostgREST, so it can be mapped to a friendly conflict.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "23505") || strings.Contains(s, "duplicate key")
}

// normalizeRows folds the retired Metode spellings into "Program Antara" so
// rows already stored in the database read back with current values.
func normalizeRows(rows []model.Rekrutmen) {
	for i := range rows {
		rows[i].Metode = model.NormalizeMetode(rows[i].Metode)
	}
}

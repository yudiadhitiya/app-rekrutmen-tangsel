// Package sheets implements Store against the Google Sheets API.
//
// Credentials are a service-account JSON key held in the backend environment;
// they never reach the browser. Authentication uses the standard OAuth2
// service-account JWT flow, signed with the Go standard library.
package sheets

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/store"
)

const (
	// SheetRekrutmen and SheetAktivitas are the tab names used in the spreadsheet.
	SheetRekrutmen = "Rekrutmen"
	SheetAktivitas = "Aktivitas"

	tokenURL = "https://oauth2.googleapis.com/token"
	scope    = "https://www.googleapis.com/auth/spreadsheets"

	// sheetsAPI is the Google Sheets API root.
	sheetsAPI = "https://sheets.googleapis.com/v4/spreadsheets"
)

// Store is a Google Sheets backed Store.
type Store struct {
	spreadsheetID string
	key           *serviceKey
	http          *http.Client

	// baseURL is the Sheets API root. It is a field so tests can point the
	// store at a fake server.
	baseURL string

	mu        sync.Mutex
	accessTok string
	tokExpiry time.Time
}

// serviceKey holds the fields of a Google service-account JSON key.
type serviceKey struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
}

// New builds a Store from a spreadsheet ID and a service-account key file.
func New(spreadsheetID, keyFile string) (*Store, error) {
	raw, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("google sheets: baca kunci: %w", err)
	}
	var k serviceKey
	if err := json.Unmarshal(raw, &k); err != nil {
		return nil, fmt.Errorf("google sheets: parse kunci: %w", err)
	}
	if k.ClientEmail == "" || k.PrivateKey == "" {
		return nil, fmt.Errorf("google sheets: kunci tidak lengkap")
	}
	if k.TokenURI == "" {
		k.TokenURI = tokenURL
	}
	return &Store{
		spreadsheetID: spreadsheetID,
		key:           &k,
		http:          &http.Client{Timeout: 30 * time.Second},
		baseURL:       sheetsAPI,
	}, nil
}

// Name implements store.Store.
func (s *Store) Name() string { return "google_sheets" }

// Ping verifies the spreadsheet is readable.
func (s *Store) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// A fresh spreadsheet has no tabs yet; create them with their headers so
	// the first run works without manual setup.
	return s.ensureTabs(ctx)
}

// ensureTabs creates the two tabs with their header row when they are missing.
func (s *Store) ensureTabs(ctx context.Context) error {
	for _, spec := range []struct {
		sheet   string
		headers []string
	}{
		{SheetRekrutmen, rekrutmenHeaders},
		{SheetAktivitas, aktivitasHeaders},
	} {
		rows, err := s.readSheet(ctx, spec.sheet)
		if err != nil {
			// The tab does not exist yet; create it below.
			if !s.isMissingTab(err) {
				return err
			}
			rows = nil
		}
		if len(rows) > 0 {
			continue
		}
		if err := s.createSheet(ctx, spec.sheet); err != nil {
			return err
		}
		if err := s.appendRow(ctx, spec.sheet, spec.headers); err != nil {
			return err
		}
	}
	return nil
}

// isMissingTab reports whether the error means the tab does not exist.
func (s *Store) isMissingTab(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Unable to parse range")
}

// createSheet adds a new tab to the spreadsheet.
func (s *Store) createSheet(ctx context.Context, sheet string) error {
	tok, err := s.token(ctx)
	if err != nil {
		return err
	}
	u := fmt.Sprintf("%s/%s?includeGridData=false",
		s.baseURL, url.PathEscape(s.spreadsheetID))
	b, err := json.Marshal(map[string]any{
		"requests": []map[string]any{{
			"addSheet": map[string]any{"properties": map[string]any{"title": sheet}},
		}},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("google sheets: create tab %s: %d: %s", sheet, resp.StatusCode, string(body))
	}
	return nil
}

// ---------------------------------------------------------------------------
// OAuth2 service-account authentication
// ---------------------------------------------------------------------------

// token returns a cached access token, refreshing it when it is close to expiry.
func (s *Store) token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.accessTok != "" && time.Now().Before(s.tokExpiry.Add(-time.Minute)) {
		return s.accessTok, nil
	}
	assertion, err := s.assertion()
	if err != nil {
		return "", err
	}
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Set("assertion", assertion)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.key.TokenURI, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google sheets: oauth %d: %s", resp.StatusCode, string(body))
	}
	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	s.accessTok = out.AccessToken
	s.tokExpiry = time.Now().Add(time.Duration(out.ExpiresIn) * time.Second)
	return s.accessTok, nil
}

// assertion builds the signed JWT used as the OAuth2 client assertion.
func (s *Store) assertion() (string, error) {
	now := time.Now()
	header := map[string]string{"alg": "RS256", "typ": "JWT"}
	claims := map[string]any{
		"iss":   s.key.ClientEmail,
		"scope": scope,
		"aud":   s.key.TokenURI,
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}
	hb, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	cb, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	enc := base64.RawURLEncoding
	signing := enc.EncodeToString(hb) + "." + enc.EncodeToString(cb)

	block, _ := pem.Decode([]byte(s.key.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("google sheets: private key bukan PEM valid")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("google sheets: parse private key: %w", err)
	}
	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("google sheets: private key bukan RSA")
	}
	sum := sha256.Sum256([]byte(signing))
	sig, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, sum[:])
	if err != nil {
		return "", err
	}
	return signing + "." + enc.EncodeToString(sig), nil
}

// ---------------------------------------------------------------------------
// Sheet I/O
// ---------------------------------------------------------------------------

// sheetRow is one spreadsheet row, keyed by header name.
type sheetRow map[string]string

// sheetEntry is a data row together with its 1-based row number in the sheet.
// Tracking the real row number keeps updates and deletes correct even after
// other rows have been cleared (a cleared row still occupies its position).
type sheetEntry struct {
	num int
	row sheetRow
}

// readSheet returns all data rows of a tab, keyed by the header row.
func (s *Store) readSheet(ctx context.Context, sheet string) ([]sheetEntry, error) {
	tok, err := s.token(ctx)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s/%s/values/%s!A1:Z",
		s.baseURL, url.PathEscape(s.spreadsheetID), url.PathEscape(sheet))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google sheets: read %s: %d: %s", sheet, resp.StatusCode, string(body))
	}
	var out struct {
		Values [][]string `json:"values"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if len(out.Values) == 0 {
		return nil, nil
	}
	header := out.Values[0]
	entries := make([]sheetEntry, 0, len(out.Values)-1)
	for i, vals := range out.Values[1:] {
		if len(vals) == 0 || (len(vals) == 1 && vals[0] == "") {
			continue
		}
		row := sheetRow{}
		for i, h := range header {
			if i < len(vals) {
				row[h] = vals[i]
			}
		}
		// Row 1 is the header, so the first data row is row 2.
		entries = append(entries, sheetEntry{num: i + 2, row: row})
	}
	return entries, nil
}

// appendRow appends a row to a tab.
func (s *Store) appendRow(ctx context.Context, sheet string, values []string) error {
	tok, err := s.token(ctx)
	if err != nil {
		return err
	}
	u := fmt.Sprintf("%s/%s/values/%s:append?valueInputOption=RAW&insertDataOption=INSERT_ROWS",
		s.baseURL, url.PathEscape(s.spreadsheetID), url.PathEscape(sheet))
	b, err := json.Marshal(map[string]any{"values": [][]string{values}})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("google sheets: append %s: %d: %s", sheet, resp.StatusCode, string(body))
	}
	return nil
}

// updateRow writes values into a specific row (1-based, including the header).
func (s *Store) updateRow(ctx context.Context, sheet string, row int, values []string) error {
	tok, err := s.token(ctx)
	if err != nil {
		return err
	}
	end := len(values)
	u := fmt.Sprintf("%s/%s/values/%s!A%d:Z%d",
		s.baseURL, url.PathEscape(s.spreadsheetID), url.PathEscape(sheet), row, row)
	b, err := json.Marshal(map[string]any{"values": [][]string{values}})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("google sheets: update %s: %d: %s", sheet, resp.StatusCode, string(body))
	}
	_ = end
	return nil
}

// deleteRow clears the contents of a row.
func (s *Store) deleteRow(ctx context.Context, sheet string, row int) error {
	tok, err := s.token(ctx)
	if err != nil {
		return err
	}
	u := fmt.Sprintf("%s/%s/values/%s!A%d:Z%d:clear",
		s.baseURL, url.PathEscape(s.spreadsheetID), url.PathEscape(sheet), row, row)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader("{}"))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("google sheets: delete %s: %d: %s", sheet, resp.StatusCode, string(body))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Column definitions
// ---------------------------------------------------------------------------

var rekrutmenHeaders = []string{
	"id", "nama", "no_wa", "alamat", "rt", "rw", "kecamatan", "kelurahan",
	"aktivitas", "tingkat_penerimaan", "metode", "status", "keterangan", "upl",
	"created_at", "updated_at", "created_by", "updated_by",
}

var aktivitasHeaders = []string{"id", "nama", "aktif"}

func rekrutmenToValues(r model.Rekrutmen) []string {
	return []string{
		r.ID, r.Nama, r.NoWA, r.Alamat, r.RT, r.RW, r.Kecamatan, r.Kelurahan,
		r.Aktivitas, r.TingkatPenerimaan, r.Metode, r.Status, r.Keterangan, r.UPL,
		r.CreatedAt.UTC().Format(time.RFC3339), r.UpdatedAt.UTC().Format(time.RFC3339),
		r.CreatedBy, r.UpdatedBy,
	}
}

func valuesToRekrutmen(row sheetRow) model.Rekrutmen {
	parseTime := func(s string) time.Time {
		t, _ := time.Parse(time.RFC3339, s)
		return t
	}
	return model.Rekrutmen{
		ID:                row["id"],
		Nama:              row["nama"],
		NoWA:              row["no_wa"],
		Alamat:            row["alamat"],
		RT:                row["rt"],
		RW:                row["rw"],
		Kecamatan:         row["kecamatan"],
		Kelurahan:         row["kelurahan"],
		Aktivitas:         row["aktivitas"],
		TingkatPenerimaan: row["tingkat_penerimaan"],
		Metode:            model.NormalizeMetode(row["metode"]),
		Status:            row["status"],
		Keterangan:        row["keterangan"],
		UPL:               row["upl"],
		CreatedAt:         parseTime(row["created_at"]),
		UpdatedAt:         parseTime(row["updated_at"]),
		CreatedBy:         row["created_by"],
		UpdatedBy:         row["updated_by"],
	}
}

func aktivitasToValues(a model.Aktivitas) []string {
	return []string{a.ID, a.Nama, fmt.Sprint(a.Aktif)}
}

func valuesToAktivitas(row sheetRow) model.Aktivitas {
	return model.Aktivitas{
		ID:    row["id"],
		Nama:  row["nama"],
		Aktif: strings.EqualFold(row["aktif"], "true"),
	}
}

// ---------------------------------------------------------------------------
// Store implementation
// ---------------------------------------------------------------------------

func (s *Store) ListRekrutmen(f model.Filter) ([]model.Rekrutmen, int, error) {
	ctx := context.Background()
	rows, err := s.readSheet(ctx, SheetRekrutmen)
	if err != nil {
		return nil, 0, err
	}
	q := strings.ToLower(strings.TrimSpace(f.Q))
	var out []model.Rekrutmen
	for _, e := range rows {
		r := valuesToRekrutmen(e.row)
		if q != "" && !strings.Contains(strings.ToLower(r.Nama), q) && !strings.Contains(r.NoWA, q) {
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
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
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
	rows, err := s.readSheet(context.Background(), SheetRekrutmen)
	if err != nil {
		return nil, err
	}
	for _, e := range rows {
		if e.row["id"] == id {
			r := valuesToRekrutmen(e.row)
			return &r, nil
		}
	}
	return nil, store.ErrNotFound
}

func (s *Store) FindRekrutmenByWA(noWA string) (*model.Rekrutmen, error) {
	rows, err := s.readSheet(context.Background(), SheetRekrutmen)
	if err != nil {
		return nil, err
	}
	for _, e := range rows {
		if e.row["no_wa"] == noWA {
			r := valuesToRekrutmen(e.row)
			return &r, nil
		}
	}
	return nil, store.ErrNotFound
}

func (s *Store) CreateRekrutmen(r *model.Rekrutmen) error {
	ctx := context.Background()
	if _, err := s.FindRekrutmenByWA(r.NoWA); err == nil {
		return store.ErrDuplicateWA
	} else if err != store.ErrNotFound {
		return err
	}
	r.Metode = model.NormalizeMetode(r.Metode)
	return s.appendRow(ctx, SheetRekrutmen, rekrutmenToValues(*r))
}

func (s *Store) UpdateRekrutmen(r *model.Rekrutmen) error {
	ctx := context.Background()
	rows, err := s.readSheet(ctx, SheetRekrutmen)
	if err != nil {
		return err
	}
	for _, e := range rows {
		if e.row["id"] != r.ID && e.row["no_wa"] == r.NoWA {
			return store.ErrDuplicateWA
		}
	}
	for _, e := range rows {
		if e.row["id"] == r.ID {
			r.Metode = model.NormalizeMetode(r.Metode)
			return s.updateRow(ctx, SheetRekrutmen, e.num, rekrutmenToValues(*r))
		}
	}
	return store.ErrNotFound
}

func (s *Store) DeleteRekrutmen(id string) error {
	ctx := context.Background()
	rows, err := s.readSheet(ctx, SheetRekrutmen)
	if err != nil {
		return err
	}
	for _, e := range rows {
		if e.row["id"] == id {
			return s.deleteRow(ctx, SheetRekrutmen, e.num)
		}
	}
	return store.ErrNotFound
}

func (s *Store) ListAktivitas() ([]model.Aktivitas, error) {
	rows, err := s.readSheet(context.Background(), SheetAktivitas)
	if err != nil {
		return nil, err
	}
	out := make([]model.Aktivitas, 0, len(rows))
	for _, e := range rows {
		out = append(out, valuesToAktivitas(e.row))
	}
	return out, nil
}

func (s *Store) CreateAktivitas(a *model.Aktivitas) error {
	ctx := context.Background()
	existing, err := s.ListAktivitas()
	if err != nil {
		return err
	}
	for _, e := range existing {
		if strings.EqualFold(e.Nama, a.Nama) {
			return store.ErrDuplicateName
		}
	}
	return s.appendRow(ctx, SheetAktivitas, aktivitasToValues(*a))
}

func (s *Store) UpdateAktivitas(a *model.Aktivitas) error {
	ctx := context.Background()
	rows, err := s.readSheet(ctx, SheetAktivitas)
	if err != nil {
		return err
	}
	for _, e := range rows {
		if e.row["id"] != a.ID && strings.EqualFold(e.row["nama"], a.Nama) {
			return store.ErrDuplicateName
		}
	}
	oldName := ""
	for _, e := range rows {
		if e.row["id"] == a.ID {
			oldName = e.row["nama"]
			if err := s.updateRow(ctx, SheetAktivitas, e.num, aktivitasToValues(*a)); err != nil {
				return err
			}
			if oldName != a.Nama {
				return s.renameAktivitas(ctx, oldName, a.Nama)
			}
			return nil
		}
	}
	return store.ErrNotFound
}

// renameAktivitas rewrites the aktivitas column for every record using oldName.
func (s *Store) renameAktivitas(ctx context.Context, oldName, newName string) error {
	rows, err := s.readSheet(ctx, SheetRekrutmen)
	if err != nil {
		return err
	}
	for _, e := range rows {
		if e.row["aktivitas"] != oldName {
			continue
		}
		r := valuesToRekrutmen(e.row)
		r.Aktivitas = newName
		if err := s.updateRow(ctx, SheetRekrutmen, e.num, rekrutmenToValues(r)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) DeleteAktivitas(id string) error {
	ctx := context.Background()
	rows, err := s.readSheet(ctx, SheetAktivitas)
	if err != nil {
		return err
	}
	name := ""
	for _, e := range rows {
		if e.row["id"] == id {
			name = e.row["nama"]
			break
		}
	}
	if name == "" {
		return store.ErrNotFound
	}
	reks, err := s.readSheet(ctx, SheetRekrutmen)
	if err != nil {
		return err
	}
	for _, e := range reks {
		if e.row["aktivitas"] == name {
			return store.ErrAktivitasUsed
		}
	}
	for _, e := range rows {
		if e.row["id"] == id {
			return s.deleteRow(ctx, SheetAktivitas, e.num)
		}
	}
	return store.ErrNotFound
}

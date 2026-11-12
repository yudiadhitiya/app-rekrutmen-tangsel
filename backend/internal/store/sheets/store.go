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
	// Default tab names used when the configuration does not name them.
	DefaultTabRekrutmen = "Rekrutmen"
	DefaultTabAktivitas = "Aktivitas"

	// SheetRekrutmen and SheetAktivitas are the default tab names used in the
	// spreadsheet.
	SheetRekrutmen = DefaultTabRekrutmen
	SheetAktivitas = DefaultTabAktivitas

	tokenURL = "https://oauth2.googleapis.com/token"
	scope    = "https://www.googleapis.com/auth/spreadsheets"

	// sheetsAPI is the Google Sheets API root.
	sheetsAPI = "https://sheets.googleapis.com/v4/spreadsheets"

	// defaultColumns is the column count assumed for a tab whose size is not
	// known yet.
	defaultColumns = 26

	// tabsCacheTTL is how long the spreadsheet tab list is cached.
	tabsCacheTTL = 60 * time.Second
)

// Modes for the Sheets provider.
const (
	// ModeApp means the application owns its tabs inside the spreadsheet: it
	// creates them when missing and stores its records there.
	ModeApp = "app"

	// ModeForm means the application reads an existing Google Form response
	// tab. In this mode the provider never creates a tab; it only verifies
	// that the configured tab exists.
	ModeForm = "form"
)

// Config configures a Google Sheets backed Store.
type Config struct {
	// SpreadsheetID is the id of the target spreadsheet.
	SpreadsheetID string

	// KeyFile is the path to the service-account JSON key.
	KeyFile string

	// TabRekrutmen is the tab holding the recruitment records.
	TabRekrutmen string

	// TabAktivitas is the tab holding the master activity list.
	TabAktivitas string

	// TabSumber is the tab holding the source rows to import. It is only used
	// in ModeForm.
	TabSumber string

	// TabPengguna is the tab holding the accounts of the previous
	// application. It is used to sign in with the old credentials.
	TabPengguna string

	// Mode is ModeApp or ModeForm. An empty value means ModeApp.
	Mode string
}

// Store is a Google Sheets backed Store.
type Store struct {
	spreadsheetID string
	key           *serviceKey
	http          *http.Client

	// baseURL is the Sheets API root. It is a field so tests can point the
	// store at a fake server.
	baseURL string

	// tabRekrutmen, tabAktivitas, tabSumber and tabPengguna are the
	// configured tab names.
	tabRekrutmen string
	tabAktivitas string
	tabSumber    string
	tabPengguna  string

	// mode is ModeApp or ModeForm.
	mode string

	mu        sync.Mutex
	accessTok string
	tokExpiry time.Time

	// tabsMu guards the cached tab list.
	tabsMu  sync.Mutex
	tabs    []sheetInfo
	tabsAt  time.Time
	tabsErr error
}

// sheetInfo describes one tab of the spreadsheet.
type sheetInfo struct {
	// Title is the tab name.
	Title string

	// Columns is the number of columns of the tab's grid.
	Columns int
}

// serviceKey holds the fields of a Google service-account JSON key.
type serviceKey struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
}

// New builds a Store from a configuration.
//
// The tab names and the mode are taken from cfg; empty values fall back to
// the defaults, so a caller that only knows the spreadsheet id and the key
// file still gets a working store.
func New(cfg Config) (*Store, error) {
	spreadsheetID, err := normalizeSpreadsheetID(cfg.SpreadsheetID)
	if err != nil {
		return nil, err
	}
	cfg.SpreadsheetID = spreadsheetID
	if err := checkKeyFilePerm(cfg.KeyFile); err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(cfg.KeyFile)
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
	if cfg.TabRekrutmen == "" {
		cfg.TabRekrutmen = DefaultTabRekrutmen
	}
	if cfg.TabAktivitas == "" {
		cfg.TabAktivitas = DefaultTabAktivitas
	}
	if cfg.TabPengguna == "" {
		cfg.TabPengguna = DefaultTabPengguna
	}
	if cfg.Mode == "" {
		cfg.Mode = ModeApp
	}
	if cfg.Mode != ModeApp && cfg.Mode != ModeForm {
		return nil, fmt.Errorf("google sheets: mode tidak dikenal: %s", cfg.Mode)
	}
	return &Store{
		spreadsheetID: cfg.SpreadsheetID,
		key:           &k,
		http:          &http.Client{Timeout: 30 * time.Second},
		baseURL:       sheetsAPI,
		tabRekrutmen:  cfg.TabRekrutmen,
		tabAktivitas:  cfg.TabAktivitas,
		tabSumber:     cfg.TabSumber,
		tabPengguna:   cfg.TabPengguna,
		mode:          cfg.Mode,
	}, nil
}

// Name implements store.Store.
func (s *Store) Name() string { return "google_sheets" }

// Ping verifies the spreadsheet is readable.
//
// In ModeApp it also prepares the tabs the application owns: they are created
// with their header row when missing, so the first run works without manual
// setup. In ModeForm it only checks that the configured source tab exists and
// never writes to the spreadsheet.
func (s *Store) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if s.mode == ModeForm {
		return s.checkSourceTab(ctx)
	}
	// A fresh spreadsheet has no tabs yet; create them with their headers so
	// the first run works without manual setup.
	return s.ensureTabs(ctx)
}

// ensureTabs creates the two tabs with their header row when they are missing.
func (s *Store) ensureTabs(ctx context.Context) error {
	tabs, err := s.listTabs(ctx)
	if err != nil {
		return err
	}
	for _, spec := range []struct {
		sheet   string
		headers []string
	}{
		{s.tabRekrutmen, rekrutmenHeaders},
		{s.tabAktivitas, aktivitasHeaders},
	} {
		if !hasTab(tabs, spec.sheet) {
			if err := s.createSheet(ctx, spec.sheet); err != nil {
				return err
			}
			if err := s.appendRow(ctx, spec.sheet, spec.headers); err != nil {
				return err
			}
			continue
		}
		rows, err := s.readSheet(ctx, spec.sheet)
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			if err := s.appendRow(ctx, spec.sheet, spec.headers); err != nil {
				return err
			}
		}
	}
	return nil
}

// checkSourceTab verifies that the tab holding the source rows exists.
//
// In ModeForm the provider must never create a tab, because the spreadsheet
// belongs to a Google Form and its tab layout is fixed by the form.
func (s *Store) checkSourceTab(ctx context.Context) error {
	if s.tabSumber == "" {
		return fmt.Errorf("google sheets: tab sumber belum dikonfigurasi")
	}
	tabs, err := s.listTabs(ctx)
	if err != nil {
		return err
	}
	if !hasTab(tabs, s.tabSumber) {
		return fmt.Errorf("google sheets: tab %q tidak ditemukan; tab yang ada: %s",
			s.tabSumber, strings.Join(tabTitles(tabs), ", "))
	}
	return nil
}

// hasTab reports whether the list contains a tab with the given title.
func hasTab(tabs []sheetInfo, title string) bool {
	for _, t := range tabs {
		if t.Title == title {
			return true
		}
	}
	return false
}

// tabTitles returns the titles of the given tabs.
func tabTitles(tabs []sheetInfo) []string {
	out := make([]string, 0, len(tabs))
	for _, t := range tabs {
		out = append(out, t.Title)
	}
	return out
}

// listTabs returns the tabs of the spreadsheet.
//
// The result is cached briefly, because every read needs to know the tab
// layout and the metadata endpoint is one request on its own.
func (s *Store) listTabs(ctx context.Context) ([]sheetInfo, error) {
	s.tabsMu.Lock()
	defer s.tabsMu.Unlock()
	if s.tabs != nil && time.Since(s.tabsAt) < tabsCacheTTL {
		return s.tabs, nil
	}
	tabs, err := s.fetchTabs(ctx)
	if err != nil {
		// A stale list is better than failing the request outright, but only
		// when one was ever fetched successfully.
		if s.tabs != nil {
			return s.tabs, nil
		}
		return nil, err
	}
	s.tabs = tabs
	s.tabsAt = time.Now()
	return tabs, nil
}

// invalidateTabs drops the cached tab list after the spreadsheet changed.
func (s *Store) invalidateTabs() {
	s.tabsMu.Lock()
	defer s.tabsMu.Unlock()
	s.tabs = nil
}

// fetchTabs reads the spreadsheet metadata and returns its tabs.
func (s *Store) fetchTabs(ctx context.Context) ([]sheetInfo, error) {
	tok, err := s.token(ctx)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s/%s", s.baseURL, url.PathEscape(s.spreadsheetID))
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
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google sheets: metadata spreadsheet: %d: %s", resp.StatusCode, string(body))
	}
	var out struct {
		Sheets []struct {
			Properties struct {
				Title          string `json:"title"`
				GridProperties struct {
					ColumnCount int `json:"columnCount"`
				} `json:"gridProperties"`
			} `json:"properties"`
		} `json:"sheets"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("google sheets: parse metadata: %w", err)
	}
	tabs := make([]sheetInfo, 0, len(out.Sheets))
	for _, sh := range out.Sheets {
		cols := sh.Properties.GridProperties.ColumnCount
		if cols <= 0 {
			cols = defaultColumns
		}
		tabs = append(tabs, sheetInfo{Title: sh.Properties.Title, Columns: cols})
	}
	return tabs, nil
}

// tabColumns returns the column count of a tab, falling back to the standard
// 26 columns when the tab is unknown.
func (s *Store) tabColumns(ctx context.Context, tab string) int {
	tabs, err := s.listTabs(ctx)
	if err != nil {
		return defaultColumns
	}
	for _, t := range tabs {
		if t.Title == tab {
			return t.Columns
		}
	}
	return defaultColumns
}

// colLetter converts a 1-based column number to its spreadsheet letter, so
// column 27 becomes "AA".
func colLetter(n int) string {
	if n <= 0 {
		n = defaultColumns
	}
	letters := ""
	for n > 0 {
		n--
		letters = string(rune('A'+(n%26))) + letters
		n /= 26
	}
	return letters
}

// createSheet adds a new tab to the spreadsheet.
func (s *Store) createSheet(ctx context.Context, sheet string) error {
	tok, err := s.token(ctx)
	if err != nil {
		return err
	}
	// A tab is added through the batchUpdate endpoint of the spreadsheet.
	u := fmt.Sprintf("%s/%s:batchUpdate", s.baseURL, url.PathEscape(s.spreadsheetID))
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
	s.invalidateTabs()
	return nil
}

// normalizeSpreadsheetID accepts either a bare spreadsheet id or the full
// sheet URL, because the id is usually copied out of the browser address bar.
//
// The returned id is the part between "/d/" and the next path separator.
func normalizeSpreadsheetID(v string) (string, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", fmt.Errorf("google sheets: spreadsheet id kosong")
	}
	if strings.Contains(v, "/") {
		i := strings.Index(v, "/d/")
		if i < 0 {
			return "", fmt.Errorf("google sheets: url spreadsheet tidak dikenali: %q", v)
		}
		v = v[i+3:]
		if j := strings.IndexAny(v, "/?#&"); j >= 0 {
			v = v[:j]
		}
	}
	if strings.ContainsAny(v, " \t") {
		return "", fmt.Errorf("google sheets: id spreadsheet tidak valid: %q", v)
	}
	if len(v) < 20 {
		return "", fmt.Errorf("google sheets: id spreadsheet terlalu pendek: %q", v)
	}
	return v, nil
}

// checkKeyFilePerm rejects a credentials file that anyone but the owner can
// read.
//
// The service-account private key grants full read/write access to every
// spreadsheet shared with the account, so the file must stay owner-only. A
// group- or world-readable key is rejected outright rather than warned about,
// because the deployment environment is the only place the file exists.
func checkKeyFilePerm(keyFile string) error {
	info, err := os.Stat(keyFile)
	if err != nil {
		return fmt.Errorf("google sheets: stat kunci: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("google sheets: kunci bukan file biasa: %s", keyFile)
	}
	if perm := info.Mode().Perm(); perm&0o077 != 0 {
		return fmt.Errorf("google sheets: izin kunci %s terlalu terbuka (%#o); jalankan: chmod 600 %s",
			keyFile, perm, keyFile)
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
	// The range covers the whole grid of the tab, so a tab that is wider or
	// narrower than the default 26 columns is read completely.
	last := colLetter(s.tabColumns(ctx, sheet))
	u := fmt.Sprintf("%s/%s/values/%s!A1:%s",
		s.baseURL, url.PathEscape(s.spreadsheetID), url.PathEscape(sheet), last)
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
	last := colLetter(s.tabColumns(ctx, sheet))
	u := fmt.Sprintf("%s/%s/values/%s!A%d:%s%d",
		s.baseURL, url.PathEscape(s.spreadsheetID), url.PathEscape(sheet), row, last, row)
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
	return nil
}

// deleteRow clears the contents of a row.
func (s *Store) deleteRow(ctx context.Context, sheet string, row int) error {
	tok, err := s.token(ctx)
	if err != nil {
		return err
	}
	last := colLetter(s.tabColumns(ctx, sheet))
	u := fmt.Sprintf("%s/%s/values/%s!A%d:%s%d:clear",
		s.baseURL, url.PathEscape(s.spreadsheetID), url.PathEscape(sheet), row, last, row)
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
	rows, err := s.readSheet(ctx, s.tabRekrutmen)
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
	rows, err := s.readSheet(context.Background(), s.tabRekrutmen)
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
	rows, err := s.readSheet(context.Background(), s.tabRekrutmen)
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
	return s.appendRow(ctx, s.tabRekrutmen, rekrutmenToValues(*r))
}

func (s *Store) UpdateRekrutmen(r *model.Rekrutmen) error {
	ctx := context.Background()
	rows, err := s.readSheet(ctx, s.tabRekrutmen)
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
			return s.updateRow(ctx, s.tabRekrutmen, e.num, rekrutmenToValues(*r))
		}
	}
	return store.ErrNotFound
}

func (s *Store) DeleteRekrutmen(id string) error {
	ctx := context.Background()
	rows, err := s.readSheet(ctx, s.tabRekrutmen)
	if err != nil {
		return err
	}
	for _, e := range rows {
		if e.row["id"] == id {
			return s.deleteRow(ctx, s.tabRekrutmen, e.num)
		}
	}
	return store.ErrNotFound
}

func (s *Store) ListAktivitas() ([]model.Aktivitas, error) {
	rows, err := s.readSheet(context.Background(), s.tabAktivitas)
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
	return s.appendRow(ctx, s.tabAktivitas, aktivitasToValues(*a))
}

func (s *Store) UpdateAktivitas(a *model.Aktivitas) error {
	ctx := context.Background()
	rows, err := s.readSheet(ctx, s.tabAktivitas)
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
			if err := s.updateRow(ctx, s.tabAktivitas, e.num, aktivitasToValues(*a)); err != nil {
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
	rows, err := s.readSheet(ctx, s.tabRekrutmen)
	if err != nil {
		return err
	}
	for _, e := range rows {
		if e.row["aktivitas"] != oldName {
			continue
		}
		r := valuesToRekrutmen(e.row)
		r.Aktivitas = newName
		if err := s.updateRow(ctx, s.tabRekrutmen, e.num, rekrutmenToValues(r)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) DeleteAktivitas(id string) error {
	ctx := context.Background()
	rows, err := s.readSheet(ctx, s.tabAktivitas)
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
	reks, err := s.readSheet(ctx, s.tabRekrutmen)
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
			return s.deleteRow(ctx, s.tabAktivitas, e.num)
		}
	}
	return store.ErrNotFound
}

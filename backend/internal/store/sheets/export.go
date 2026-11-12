package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
)

// Tabs returns the titles of the tabs in the spreadsheet.
func (s *Store) Tabs(ctx context.Context) ([]string, error) {
	tabs, err := s.listTabs(ctx)
	if err != nil {
		return nil, err
	}
	return tabTitles(tabs), nil
}

// ReadTab returns the header row and the data rows of a tab.
//
// The rows are returned as they are in the sheet, so a caller that knows the
// column layout can read them positionally. Empty rows are skipped.
func (s *Store) ReadTab(ctx context.Context, tab string) ([]string, [][]string, error) {
	tok, err := s.token(ctx)
	if err != nil {
		return nil, nil, err
	}
	last := colLetter(s.tabColumns(ctx, tab))
	u := fmt.Sprintf("%s/%s/values/%s!A1:%s",
		s.baseURL, url.PathEscape(s.spreadsheetID), url.PathEscape(tab), last)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("google sheets: read %s: %d: %s", tab, resp.StatusCode, string(body))
	}
	var out struct {
		Values [][]string `json:"values"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, nil, fmt.Errorf("google sheets: parse %s: %w", tab, err)
	}
	if len(out.Values) == 0 {
		return nil, nil, nil
	}
	headers := out.Values[0]
	rows := make([][]string, 0, len(out.Values)-1)
	for _, vals := range out.Values[1:] {
		if len(vals) == 0 || (len(vals) == 1 && vals[0] == "") {
			continue
		}
		rows = append(rows, vals)
	}
	return headers, rows, nil
}

// EnsureTab creates the tab with a header row when it is missing.
//
// It is used by the import tool, which owns its own target tab.
func (s *Store) EnsureTab(ctx context.Context, tab string, headers []string) error {
	tabs, err := s.listTabs(ctx)
	if err != nil {
		return err
	}
	if !hasTab(tabs, tab) {
		if err := s.createSheet(ctx, tab); err != nil {
			return err
		}
		return s.appendRow(ctx, tab, headers)
	}
	_, rows, err := s.ReadTab(ctx, tab)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return s.appendRow(ctx, tab, headers)
	}
	return nil
}

// AppendRow appends one row to a tab.
func (s *Store) AppendRow(ctx context.Context, tab string, values []string) error {
	return s.appendRow(ctx, tab, values)
}

// ClearRow empties the contents of one row of a tab.
//
// The row number is the 1-based number of the row in the sheet. The row
// itself stays in place; only its cells are emptied.
func (s *Store) ClearRow(tab string, row int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.deleteRow(ctx, tab, row)
}

// RekrutmenHeaders returns the column layout of the rekrutmen tab.
func RekrutmenHeaders() []string {
	out := make([]string, len(rekrutmenHeaders))
	copy(out, rekrutmenHeaders)
	return out
}

// AktivitasHeaders returns the column layout of the aktivitas tab.
func AktivitasHeaders() []string {
	out := make([]string, len(aktivitasHeaders))
	copy(out, aktivitasHeaders)
	return out
}

// RekrutmenToValues converts a record into the values of a sheet row.
func RekrutmenToValues(r model.Rekrutmen) []string {
	return rekrutmenToValues(r)
}

// AktivitasToValues converts a master entry into the values of a sheet row.
func AktivitasToValues(a model.Aktivitas) []string {
	return aktivitasToValues(a)
}

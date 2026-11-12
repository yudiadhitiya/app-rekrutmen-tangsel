// Command impor-responses imports the rows of a Google Form response sheet
// into the tab the application owns.
//
// The import is add-only: a row whose WhatsApp number is already present in
// the target tab is skipped, so the command can be run more than once without
// creating duplicates. Without --apply nothing is written to the spreadsheet;
// the command only prints what it would do.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/store/sheets"
	"rekrutmen-tangsel/backend/internal/validate"
)

// report is the outcome of one import run.
type report struct {
	SpreadsheetID string    `json:"spreadsheet_id"`
	TabSumber     string    `json:"tab_sumber"`
	TabTujuan     string    `json:"tab_tujuan"`
	Applied       bool      `json:"applied"`
	StartedAt     time.Time `json:"started_at"`
	FinishedAt    time.Time `json:"finished_at"`

	// Read is the number of data rows in the source tab.
	Read int `json:"read"`
	// Imported is the number of rows written to the target tab.
	Imported int `json:"imported"`
	// SkippedDuplicate counts rows whose number was already present.
	SkippedDuplicate int `json:"skipped_duplicate"`
	// SkippedInvalid counts rows that could not be imported at all.
	SkippedInvalid int `json:"skipped_invalid"`
	// NeedsReview counts imported rows that do not satisfy the application
	// rules and were marked for verification.
	NeedsReview int `json:"needs_review"`
	// NewAktivitas lists the activity names that were not in the master list.
	NewAktivitas []string `json:"new_aktivitas,omitempty"`
	// Problems describes every row that was skipped or needs review.
	Problems []problem `json:"problems,omitempty"`
}

// problem is one row that needs attention.
type problem struct {
	// Row is the 1-based row number in the source sheet.
	Row int `json:"row"`
	// NoWA and Nama identify the row.
	NoWA string `json:"no_wa"`
	Nama string `json:"nama"`
	// Reasons explains why the row was skipped or marked for review.
	Reasons []string `json:"reasons"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gagal:", err)
		os.Exit(1)
	}
}

func run() error {
	spreadsheetID := flag.String("spreadsheet-id", os.Getenv("GOOGLE_SHEET_ID"),
		"id spreadsheet, boleh juga URL lengkap spreadsheet")
	keyFile := flag.String("key-file", os.Getenv("GOOGLE_KEY_FILE"),
		"path file JSON service account")
	tabSumber := flag.String("tab-sumber", envOr("GOOGLE_SHEETS_TAB_SUMBER", "Form Responses 1"),
		"tab yang berisi respons formulir")
	tabTujuan := flag.String("tab-tujuan", envOr("GOOGLE_SHEETS_TAB_REKRUTMEN", sheets.DefaultTabRekrutmen),
		"tab milik aplikasi yang menerima hasil impor")
	tabAktivitas := flag.String("tab-aktivitas", envOr("GOOGLE_SHEETS_TAB_AKTIVITAS", sheets.DefaultTabAktivitas),
		"tab master aktivitas")
	apply := flag.Bool("apply", false,
		"tulis hasil ke spreadsheet; tanpa flag ini hanya menampilkan rencana")
	dryRun := flag.Bool("dry-run", false,
		"tidak menulis apa pun ke spreadsheet (perilaku bawaan)")
	reportPath := flag.String("report", "", "path file laporan JSON")
	inspect := flag.Bool("inspect", false,
		"tampilkan daftar tab beserta header dan jumlah baris, tanpa isi data")
	inspectKolom := flag.String("kolom", "",
		"untuk --inspect: daftar kolom (dipisah koma) yang nilai uniknya ditampilkan")
	inspectFill := flag.Bool("fill", false,
		"untuk --inspect: tampilkan jumlah baris terisi per kolom")
	flag.Parse()

	if *spreadsheetID == "" || *keyFile == "" {
		return fmt.Errorf("GOOGLE_SHEET_ID dan GOOGLE_KEY_FILE harus diisi")
	}
	if *dryRun && *apply {
		return fmt.Errorf("--dry-run dan --apply tidak boleh dipakai bersamaan")
	}
	if *dryRun {
		*apply = false
	}

	store, err := sheets.New(sheets.Config{
		SpreadsheetID: *spreadsheetID,
		KeyFile:       *keyFile,
		TabRekrutmen:  *tabTujuan,
		TabAktivitas:  *tabAktivitas,
		Mode:          sheets.ModeApp,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	rep := report{
		SpreadsheetID: *spreadsheetID,
		TabSumber:     *tabSumber,
		TabTujuan:     *tabTujuan,
		StartedAt:     time.Now().UTC(),
	}

	// The inspect mode reports the layout of the spreadsheet without reading
	// any row content.
	if *inspect {
		return inspectSpreadsheet(ctx, store, *tabSumber, splitList(*inspectKolom), *inspectFill)
	}

	// The source tab must exist: the tool never creates the tab of a form.
	tabs, err := store.Tabs(ctx)
	if err != nil {
		return err
	}
	if !contains(tabs, *tabSumber) {
		return fmt.Errorf("tab %q tidak ditemukan; tab yang ada: %s",
			*tabSumber, strings.Join(tabs, ", "))
	}

	headers, rows, err := store.ReadTab(ctx, *tabSumber)
	if err != nil {
		return err
	}
	rep.Read = len(rows)

	mapping, err := sheets.ResolveColumnMap(headers, sheets.FormColumnSpec())
	if err != nil {
		return err
	}

	// The numbers already in the target tab, so a second run skips them.
	existing := map[string]bool{}
	if targetHeaders, targetRows, err := store.ReadTab(ctx, *tabTujuan); err == nil {
		idx := indexOf(targetHeaders, "no_wa")
		for _, row := range targetRows {
			if idx < 0 || idx >= len(row) {
				continue
			}
			if no := sheets.NormalizeNoWA(row[idx]); no != "" {
				existing[no] = true
			}
		}
	}

	// The activity names allowed by the master list, plus the names that only
	// appear in the form responses.
	master := map[string]bool{}
	if masterHeaders, masterRows, err := store.ReadTab(ctx, *tabAktivitas); err == nil {
		idx := indexOf(masterHeaders, "nama")
		for _, row := range masterRows {
			if idx < 0 || idx >= len(row) {
				continue
			}
			if name := strings.TrimSpace(row[idx]); name != "" {
				master[name] = true
			}
		}
	}
	// The activity names used in the responses are allowed even when the
	// master list does not know them yet: they are added to the master list
	// by the import itself.
	known := make([]string, 0, len(master))
	for name := range master {
		known = append(known, name)
	}
	sort.Strings(known)

	allowed := append([]string{}, known...)
	for _, row := range rows {
		name := mapping.Get(row, sheets.FieldAktivitas)
		if name != "" && !contains(allowed, name) {
			allowed = append(allowed, name)
		}
	}
	sort.Strings(allowed)

	// Plan the import first, so nothing is written before every row is known.
	res := buildPlan(mapping, rows, allowed, known, existing, time.Now().UTC())
	rep.SkippedDuplicate = res.SkippedDuplicate
	rep.SkippedInvalid = res.SkippedInvalid
	rep.NeedsReview = res.NeedsReview
	rep.Problems = res.Problems
	rep.NewAktivitas = res.NewAktivitas
	plan := res.Plan

	// Write only when asked to.
	if *apply {
		if err := store.EnsureTab(ctx, *tabTujuan, sheets.RekrutmenHeaders()); err != nil {
			return err
		}
		if err := store.EnsureTab(ctx, *tabAktivitas, sheets.AktivitasHeaders()); err != nil {
			return err
		}
		for _, name := range res.NewAktivitas {
			if err := store.AppendRow(ctx, *tabAktivitas, sheets.AktivitasToValues(model.Aktivitas{
				ID:    sheets.DerivedID("aktivitas:"+name, time.Time{}),
				Nama:  name,
				Aktif: true,
			})); err != nil {
				return err
			}
		}
		for _, p := range plan {
			if err := store.AppendRow(ctx, *tabTujuan, sheets.RekrutmenToValues(p.Record)); err != nil {
				return fmt.Errorf("baris %d gagal ditulis: %w", p.Row, err)
			}
		}
		rep.Imported = len(plan)
		rep.Applied = true
	}

	rep.FinishedAt = time.Now().UTC()

	printSummary(rep)

	if *reportPath != "" {
		if err := writeReport(*reportPath, rep); err != nil {
			return err
		}
		fmt.Printf("laporan ditulis ke %s\n", *reportPath)
	}
	return nil
}

// inspectSpreadsheet prints the tabs of the spreadsheet with their header row
// and the number of data rows.
//
// Only column names and counts are printed: the content of the rows stays in
// the spreadsheet. When a list of columns is given, the distinct values of
// those columns are shown, which is what an operator needs to check the
// enumeration of a sheet without exposing personal data.
func inspectSpreadsheet(ctx context.Context, store *sheets.Store, tab string, kolom []string, fill bool) error {
	tabs, err := store.Tabs(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("spreadsheet punya %d tab\n", len(tabs))
	for _, name := range tabs {
		if tab != "" && !strings.EqualFold(name, tab) {
			continue
		}
		headers, rows, err := store.ReadTab(ctx, name)
		if err != nil {
			fmt.Printf("  %-22s GAGAL dibaca: %v\n", name, err)
			continue
		}
		if len(headers) == 0 {
			fmt.Printf("  %-22s (kosong)\n", name)
			continue
		}
		fmt.Printf("  %-22s %d baris, %d kolom: %s\n",
			name, len(rows), len(headers), strings.Join(headers, " | "))

		if fill {
			for i, h := range headers {
				n := 0
				for _, row := range rows {
					if i < len(row) && strings.TrimSpace(row[i]) != "" {
						n++
					}
				}
				fmt.Printf("      %-22s terisi %d/%d\n", h, n, len(rows))
			}
		}

		for _, want := range kolom {
			idx := indexOf(headers, want)
			if idx < 0 {
				fmt.Printf("      kolom %q tidak ada\n", want)
				continue
			}
			counts := map[string]int{}
			for _, row := range rows {
				if idx < len(row) {
					if v := strings.TrimSpace(row[idx]); v != "" {
						counts[v]++
					}
				}
			}
			uniq := make([]string, 0, len(counts))
			for v := range counts {
				uniq = append(uniq, v)
			}
			sort.Strings(uniq)
			shown := uniq
			if len(shown) > 5 {
				shown = shown[:5]
			}
			fmt.Printf("      %s (%d nilai unik): %s\n",
				want, len(uniq), strings.Join(shown, " | "))
		}
	}
	return nil
}

// splitList splits a comma separated list into its parts.
func splitList(v string) []string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// plannedRow is one record that is ready to be written to the target tab.
type plannedRow struct {
	// Row is the 1-based row number in the source sheet.
	Row int
	// Record is the record to write.
	Record model.Rekrutmen
}

// planResult is the outcome of planning an import.
type planResult struct {
	// Plan holds the records to write.
	Plan []plannedRow
	// Problems describes the rows that were skipped or need review.
	Problems []problem
	// SkippedDuplicate counts rows whose number was already present.
	SkippedDuplicate int
	// SkippedInvalid counts rows that could not be imported at all.
	SkippedInvalid int
	// NeedsReview counts imported rows that need verification.
	NeedsReview int
	// NewAktivitas lists the activity names missing from the master list.
	NewAktivitas []string
}

// buildPlan turns the source rows into the records to write.
//
// Rows whose WhatsApp number is already known are skipped, so the import can
// be run twice without creating duplicates. Rows that do not satisfy the
// application rules are still imported, but they are marked for verification
// and reported.
func buildPlan(mapping sheets.ColumnMap, rows [][]string, allowed, known []string, existing map[string]bool, now time.Time) planResult {
	res := planResult{}
	seen := map[string]bool{}
	newAktivitas := map[string]bool{}

	for i, row := range rows {
		rowNum := i + 2 // row 1 is the header
		rec := mapping.Record(row)
		ts, hasTS := mapping.Timestamp(row)

		if rec.NoWA == "" {
			res.SkippedInvalid++
			res.Problems = append(res.Problems, problem{
				Row: rowNum, Nama: rec.Nama,
				Reasons: []string{"nomor whatsapp kosong"},
			})
			continue
		}
		if existing[rec.NoWA] || seen[rec.NoWA] {
			res.SkippedDuplicate++
			continue
		}
		seen[rec.NoWA] = true

		// A status the application does not know is not trusted blindly.
		if rec.Status == "" || !contains(model.Status, rec.Status) {
			rec.Status = "Pending"
		}
		if rec.Aktivitas != "" && !contains(known, rec.Aktivitas) {
			newAktivitas[rec.Aktivitas] = true
		}

		reasons := fieldReasons(validate.Rekrutmen(&rec, allowed))

		rec.ID = sheets.DerivedID(rec.NoWA, ts)
		rec.CreatedAt = ts
		if !hasTS {
			rec.CreatedAt = now
		}
		rec.UpdatedAt = rec.CreatedAt
		rec.CreatedBy = "form"
		rec.UpdatedBy = "form"

		if len(reasons) > 0 {
			rec.Status = "Pending"
			rec.Keterangan = strings.TrimSpace(rec.Keterangan + " hasil impor form, perlu verifikasi: " +
				strings.Join(reasons, "; "))
			res.NeedsReview++
			res.Problems = append(res.Problems, problem{
				Row: rowNum, NoWA: rec.NoWA, Nama: rec.Nama, Reasons: reasons,
			})
		}

		res.Plan = append(res.Plan, plannedRow{Row: rowNum, Record: rec})
	}

	res.NewAktivitas = sortedKeys(newAktivitas)
	return res
}

// fieldReasons turns a validation result into readable reasons.
func fieldReasons(fields validate.Fields) []string {
	if len(fields) == 0 {
		return nil
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k+": "+fields[k])
	}
	return out
}

// printSummary writes the outcome of the run to stdout.
//
// The summary carries row numbers only. Names and phone numbers stay in the
// JSON report, so the console output can be shared or pasted into a ticket
// without exposing the personal data of the respondents.
func printSummary(rep report) {
	mode := "rencana (tidak ada yang ditulis)"
	if rep.Applied {
		mode = "diterapkan"
	}
	fmt.Printf("impor %s -> %s [%s]\n", rep.TabSumber, rep.TabTujuan, mode)
	fmt.Printf("  terbaca      : %d\n", rep.Read)
	fmt.Printf("  diimpor      : %d\n", rep.Imported)
	fmt.Printf("  duplikat     : %d\n", rep.SkippedDuplicate)
	fmt.Printf("  tidak valid  : %d\n", rep.SkippedInvalid)
	fmt.Printf("  perlu review : %d\n", rep.NeedsReview)
	if len(rep.NewAktivitas) > 0 {
		fmt.Printf("  aktivitas baru: %d (lihat laporan)\n", len(rep.NewAktivitas))
	}
	for _, p := range rep.Problems {
		fmt.Printf("  baris %d: %s\n", p.Row, strings.Join(p.Reasons, "; "))
	}
	if len(rep.Problems) > 0 {
		fmt.Println("  detail tiap baris ada di file laporan")
	}
}

// writeReport stores the report as JSON.
func writeReport(path string, rep report) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}

// envOr returns the environment value, or the fallback when it is empty.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// contains reports whether the list holds the value.
func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// indexOf returns the position of the header, or -1.
func indexOf(headers []string, name string) int {
	for i, h := range headers {
		if strings.EqualFold(strings.TrimSpace(h), name) {
			return i
		}
	}
	return -1
}

// sortedKeys returns the map keys in alphabetical order.
func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

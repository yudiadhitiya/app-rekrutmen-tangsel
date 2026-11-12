package main

import (
	"strings"
	"testing"
	"time"

	"rekrutmen-tangsel/backend/internal/store/sheets"
)

// formHeaders is a plausible header row of a Google Form response sheet.
var formHeaders = []string{
	"Timestamp", "Nama", "No. WA", "Alamat", "RT", "RW",
	"Kecamatan", "Kelurahan", "Aktivitas", "Tingkat Penerimaan",
	"Metode", "Status", "Keterangan", "UPL",
}

func testMapping(t *testing.T) sheets.ColumnMap {
	t.Helper()
	m, err := sheets.ResolveColumnMap(formHeaders, sheets.FormColumnSpec())
	if err != nil {
		t.Fatalf("ResolveColumnMap: %v", err)
	}
	return m
}

func TestBuildPlanImportsValidRow(t *testing.T) {
	m := testMapping(t)
	rows := [][]string{{
		"06/10/2025 13:45:00", "Ahmad Fauzi", "081234567890", "Jl. Melati No. 5",
		"01", "02", "Pamulang", "Pamulang Barat", "Guru", "Terbuka",
		"Halaqoh", "Proses", "", "",
	}}
	res := buildPlan(m, rows, []string{"Guru"}, []string{"Guru"}, map[string]bool{}, time.Now())

	if len(res.Plan) != 1 {
		t.Fatalf("harus ada 1 baris di rencana, dapat %d", len(res.Plan))
	}
	rec := res.Plan[0].Record
	if rec.Nama != "Ahmad Fauzi" {
		t.Errorf("nama: %q", rec.Nama)
	}
	if rec.NoWA != "081234567890" {
		t.Errorf("no_wa: %q", rec.NoWA)
	}
	if rec.ID == "" {
		t.Error("id harus terisi")
	}
	if rec.CreatedBy != "form" || rec.UpdatedBy != "form" {
		t.Errorf("jejak audit harus menandai sumber form, dapat %q/%q", rec.CreatedBy, rec.UpdatedBy)
	}
	if rec.CreatedAt.Year() != 2025 {
		t.Errorf("created_at harus dari kolom timestamp, dapat %v", rec.CreatedAt)
	}
	if res.NeedsReview != 0 {
		t.Errorf("baris yang valid tidak perlu review, dapat %d", res.NeedsReview)
	}
}

func TestBuildPlanSkipsDuplicateNumbers(t *testing.T) {
	m := testMapping(t)
	row := []string{
		"06/10/2025 13:45:00", "Ahmad Fauzi", "081234567890", "Jl. Melati No. 5",
		"01", "02", "Pamulang", "Pamulang Barat", "Guru", "Terbuka",
		"Halaqoh", "Proses", "", "",
	}
	other := []string{
		"06/10/2025 14:00:00", "Budi", "081234567890", "Jl. Kenanga",
		"01", "02", "Pamulang", "Pamulang Barat", "Guru", "Terbuka",
		"Halaqoh", "Proses", "", "",
	}
	// The number is already in the target tab.
	res := buildPlan(m, [][]string{row}, []string{"Guru"}, []string{"Guru"},
		map[string]bool{"081234567890": true}, time.Now())
	if len(res.Plan) != 0 {
		t.Fatalf("baris yang sudah ada tidak boleh diimpor, dapat %d", len(res.Plan))
	}
	if res.SkippedDuplicate != 1 {
		t.Errorf("harus mencatat 1 duplikat, dapat %d", res.SkippedDuplicate)
	}

	// The same number twice inside the source sheet.
	res = buildPlan(m, [][]string{row, other}, []string{"Guru"}, []string{"Guru"},
		map[string]bool{}, time.Now())
	if len(res.Plan) != 1 {
		t.Fatalf("hanya baris pertama yang boleh diimpor, dapat %d", len(res.Plan))
	}
	if res.SkippedDuplicate != 1 {
		t.Errorf("harus mencatat 1 duplikat, dapat %d", res.SkippedDuplicate)
	}
}

func TestBuildPlanMarksInvalidRowForReview(t *testing.T) {
	m := testMapping(t)
	rows := [][]string{{
		"06/10/2025 13:45:00", "Siti", "081234567891", "Alamat",
		"01", "02", "Kota Luar", "Kelurahan Luar", "Guru", "Terbuka",
		"Halaqoh", "Proses", "", "",
	}}
	res := buildPlan(m, rows, []string{"Guru"}, []string{"Guru"}, map[string]bool{}, time.Now())

	if len(res.Plan) != 1 {
		t.Fatalf("baris dengan masalah harus tetap diimpor, dapat %d baris", len(res.Plan))
	}
	rec := res.Plan[0].Record
	if rec.Status != "Pending" {
		t.Errorf("status harus Pending, dapat %q", rec.Status)
	}
	if !strings.Contains(rec.Keterangan, "perlu verifikasi") {
		t.Errorf("keterangan harus menjelaskan perlunya verifikasi, dapat %q", rec.Keterangan)
	}
	if res.NeedsReview != 1 {
		t.Errorf("harus mencatat 1 baris perlu review, dapat %d", res.NeedsReview)
	}
	if len(res.Problems) != 1 {
		t.Fatalf("harus ada 1 masalah tercatat, dapat %d", len(res.Problems))
	}
	if !strings.Contains(strings.Join(res.Problems[0].Reasons, " "), "kecamatan") {
		t.Errorf("alasan harus menyebut kecamatan, dapat %v", res.Problems[0].Reasons)
	}
}

func TestBuildPlanSkipsRowWithoutNumber(t *testing.T) {
	m := testMapping(t)
	rows := [][]string{{
		"06/10/2025 13:45:00", "Tanpa Nomor", "", "Alamat",
		"01", "02", "Pamulang", "Pamulang Barat", "Guru", "Terbuka",
		"Halaqoh", "Proses", "", "",
	}}
	res := buildPlan(m, rows, []string{"Guru"}, []string{"Guru"}, map[string]bool{}, time.Now())
	if len(res.Plan) != 0 {
		t.Fatalf("baris tanpa nomor tidak boleh diimpor, dapat %d", len(res.Plan))
	}
	if res.SkippedInvalid != 1 {
		t.Errorf("harus mencatat 1 baris tidak valid, dapat %d", res.SkippedInvalid)
	}
}

func TestBuildPlanReportsNewActivity(t *testing.T) {
	m := testMapping(t)
	rows := [][]string{{
		"06/10/2025 13:45:00", "Ahmad", "081234567890", "Alamat",
		"01", "02", "Pamulang", "Pamulang Barat", "Pedagang Pasar", "Terbuka",
		"Halaqoh", "Proses", "", "",
	}}
	res := buildPlan(m, rows, []string{"Guru", "Pedagang Pasar"}, []string{"Guru"},
		map[string]bool{}, time.Now())
	if len(res.NewAktivitas) != 1 || res.NewAktivitas[0] != "Pedagang Pasar" {
		t.Fatalf("aktivitas baru harus dilaporkan, dapat %v", res.NewAktivitas)
	}
}

func TestBuildPlanIsIdempotent(t *testing.T) {
	m := testMapping(t)
	rows := [][]string{{
		"06/10/2025 13:45:00", "Ahmad Fauzi", "081234567890", "Alamat",
		"01", "02", "Pamulang", "Pamulang Barat", "Guru", "Terbuka",
		"Halaqoh", "Proses", "", "",
	}}
	first := buildPlan(m, rows, []string{"Guru"}, []string{"Guru"}, map[string]bool{}, time.Now())

	// The second run sees the number that the first run produced.
	existing := map[string]bool{first.Plan[0].Record.NoWA: true}
	second := buildPlan(m, rows, []string{"Guru"}, []string{"Guru"}, existing, time.Now())

	if len(second.Plan) != 0 {
		t.Fatalf("import kedua tidak boleh menambah baris, dapat %d", len(second.Plan))
	}
	if second.SkippedDuplicate != 1 {
		t.Errorf("import kedua harus melewati baris sebagai duplikat, dapat %d", second.SkippedDuplicate)
	}
}

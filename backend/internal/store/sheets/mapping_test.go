package sheets

import (
	"strings"
	"testing"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
)

func TestResolveColumnMapMatchesAliases(t *testing.T) {
	headers := []string{
		"Timestamp", "Nama", "No. WA", "Alamat", "RT", "RW",
		"Kecamatan", "Kelurahan", "Aktivitas", "Tingkat Penerimaan",
		"Metode", "Status", "Keterangan", "UPL",
	}
	m, err := ResolveColumnMap(headers, FormColumnSpec())
	if err != nil {
		t.Fatalf("ResolveColumnMap: %v", err)
	}
	cases := map[string]int{
		FieldTimestamp:         0,
		FieldNama:              1,
		FieldNoWA:              2,
		FieldRT:                4,
		FieldKecamatan:         6,
		FieldKelurahan:         7,
		FieldAktivitas:         8,
		FieldTingkatPenerimaan: 9,
		FieldMetode:            10,
		FieldStatus:            11,
		FieldKeterangan:        12,
		FieldUPL:               13,
	}
	for field, want := range cases {
		got, ok := m.Index[field]
		if !ok {
			t.Errorf("field %q tidak terpetakan", field)
			continue
		}
		if got != want {
			t.Errorf("field %q terpetakan ke kolom %d, ingin %d", field, got, want)
		}
	}
}

func TestResolveColumnMapIgnoresPunctuation(t *testing.T) {
	headers := []string{"nama lengkap", "KECAMATAN", "no_wa"}
	m, err := ResolveColumnMap(headers, FormColumnSpec())
	if err != nil {
		t.Fatalf("ResolveColumnMap: %v", err)
	}
	if got := m.Index[FieldNoWA]; got != 2 {
		t.Errorf("no_wa harus terpetakan ke kolom 2, dapat %d", got)
	}
	if got := m.Index[FieldNama]; got != 0 {
		t.Errorf("nama harus terpetakan ke kolom 0, dapat %d", got)
	}
	if got := m.Index[FieldKecamatan]; got != 1 {
		t.Errorf("kecamatan harus terpetakan ke kolom 1, dapat %d", got)
	}
}

func TestResolveColumnMapUsesAliasWhenNeeded(t *testing.T) {
	// The sheet spells the phone column "No Whatsapp" and the name column
	// "Nama Pendaftar".
	headers := []string{"Nama Pendaftar", "No Whatsapp", "Kecamatan"}
	m, err := ResolveColumnMap(headers, FormColumnSpec())
	if err != nil {
		t.Fatalf("ResolveColumnMap: %v", err)
	}
	if got := m.Index[FieldNoWA]; got != 1 {
		t.Errorf("no_wa harus terpetakan ke kolom 1, dapat %d", got)
	}
	if got := m.Index[FieldNama]; got != 0 {
		t.Errorf("nama harus terpetakan ke kolom 0, dapat %d", got)
	}
}

func TestResolveColumnMapPrefersFirstMatch(t *testing.T) {
	// Two columns match the same field; the leftmost one wins.
	headers := []string{"Nama", "no_wa", "No Whatsapp", "Kecamatan"}
	m, err := ResolveColumnMap(headers, FormColumnSpec())
	if err != nil {
		t.Fatalf("ResolveColumnMap: %v", err)
	}
	if got := m.Index[FieldNoWA]; got != 1 {
		t.Errorf("kolom paling kiri harus dipakai, dapat %d", got)
	}
}

func TestResolveColumnMapReportsMissingRequiredColumn(t *testing.T) {
	headers := []string{"Timestamp", "Nama"}
	_, err := ResolveColumnMap(headers, FormColumnSpec())
	if err == nil {
		t.Fatal("expected an error for the missing required columns")
	}
	if !strings.Contains(err.Error(), "no_wa") || !strings.Contains(err.Error(), "kecamatan") {
		t.Fatalf("error harus menyebut kolom yang hilang, dapat: %v", err)
	}
	if !strings.Contains(err.Error(), "Timestamp") {
		t.Fatalf("error harus menyebut header yang tersedia, dapat: %v", err)
	}
}

func TestColumnMapRecord(t *testing.T) {
	headers := []string{"Timestamp", "Nama", "No. WA", "Kecamatan", "Metode", "Status"}
	m, err := ResolveColumnMap(headers, FormColumnSpec())
	if err != nil {
		t.Fatalf("ResolveColumnMap: %v", err)
	}
	row := []string{"06/10/2025 13:45:00", "Ahmad Fauzi", "+62 812-3456-7890", "Pamulang", "Program", "Berhasil"}
	rec := m.Record(row)

	if rec.Nama != "Ahmad Fauzi" {
		t.Errorf("nama: %q", rec.Nama)
	}
	if rec.NoWA != "6281234567890" {
		t.Errorf("no_wa harus dinormalisasi menjadi digit saja, dapat %q", rec.NoWA)
	}
	if rec.Kecamatan != "Pamulang" {
		t.Errorf("kecamatan: %q", rec.Kecamatan)
	}
	if rec.Metode != "Program Antara" {
		t.Errorf("metode lama harus dilipat ke Program Antara, dapat %q", rec.Metode)
	}
	ts, ok := m.Timestamp(row)
	if !ok {
		t.Fatal("timestamp harus terbaca")
	}
	if ts.Year() != 2025 || ts.Month() != time.October || ts.Day() != 6 {
		t.Errorf("timestamp: %v", ts)
	}
}

func TestParseTimestampFormats(t *testing.T) {
	cases := []string{
		"06/10/2025 13:45:00",
		"06/10/2025 13:45",
		"2025-10-06 13:45:00",
		"2025-10-06T13:45:00Z",
		"2025-10-06",
	}
	for _, v := range cases {
		if _, ok := ParseTimestamp(v); !ok {
			t.Errorf("timestamp %q harus terbaca", v)
		}
	}
	if _, ok := ParseTimestamp("bukan tanggal"); ok {
		t.Error("nilai bukan tanggal harus ditolak")
	}
	if _, ok := ParseTimestamp(""); ok {
		t.Error("nilai kosong harus ditolak")
	}
}

func TestDerivedIDIsStable(t *testing.T) {
	ts := time.Date(2025, 10, 6, 13, 45, 0, 0, time.UTC)
	a := DerivedID("6281234567890", ts)
	b := DerivedID("6281234567890", ts)
	if a != b {
		t.Fatalf("id harus stabil: %q != %q", a, b)
	}
	if len(a) != 36 {
		t.Fatalf("id harus berbentuk uuid, dapat %q", a)
	}
	c := DerivedID("6281234567890", ts.Add(time.Minute))
	if a == c {
		t.Fatal("id yang berbeda harus berbeda")
	}
}

// legacyHeaders is the header row of the tab of the previous application.
var legacyHeaders = []string{
	"ID", "NAMA ANGGOTA", "NAMA TARGET", "NO WHATSAPP", "NAMA KECAMATAN",
	"ALAMAT", "RT", "RW", "KECAMATAN", "KELURAHAN", "UPLS", "UPLM",
	"AKTIFITAS", "TINGKAT PENERIMAAN", "METODE", "STATUS", "KETERANGAN",
	"CREATED_AT", "UPDATED_AT", "CREATED_BY",
}

func TestAppColumnSpecResolvesLegacyHeaders(t *testing.T) {
	cm, err := ResolveColumnMap(legacyHeaders, AppColumnSpec())
	if err != nil {
		t.Fatalf("ResolveColumnMap: %v", err)
	}
	cases := []struct {
		field string
		col   string
	}{
		{FieldID, "ID"},
		{FieldNama, "NAMA TARGET"},
		{FieldNoWA, "NO WHATSAPP"},
		{FieldKecamatan, "KECAMATAN"},
		{FieldAktivitas, "AKTIFITAS"},
		{FieldTingkatPenerimaan, "TINGKAT PENERIMAAN"},
		{FieldUPL, "UPLS"},
		{FieldCreatedAt, "CREATED_AT"},
		{FieldCreatedBy, "CREATED_BY"},
		{FieldUploader, "NAMA ANGGOTA"},
	}
	for _, c := range cases {
		idx, ok := cm.Index[c.field]
		if !ok {
			t.Errorf("field %q tidak terpetakan", c.field)
			continue
		}
		if legacyHeaders[idx] != c.col {
			t.Errorf("field %q harus ke kolom %q, dapat %q", c.field, c.col, legacyHeaders[idx])
		}
	}
}

func TestRekrutmenRowValuesFollowsTheTabLayout(t *testing.T) {
	cm, err := ResolveColumnMap(legacyHeaders, AppColumnSpec())
	if err != nil {
		t.Fatalf("ResolveColumnMap: %v", err)
	}
	rec := model.Rekrutmen{
		ID:        "abc",
		Nama:      "Uju Juhiwa",
		NoWA:      "85198971242",
		Kecamatan: "Pondok Aren",
		Aktivitas: "Pedagang",
		CreatedBy: "umar",
	}
	vals := rekrutmenRowValues(cm, rec, nil)
	if len(vals) != len(legacyHeaders) {
		t.Fatalf("jumlah kolom: %d, harus %d", len(vals), len(legacyHeaders))
	}
	for field, want := range map[string]string{
		FieldID:        rec.ID,
		FieldNama:      rec.Nama,
		FieldNoWA:      rec.NoWA,
		FieldKecamatan: rec.Kecamatan,
		FieldAktivitas: rec.Aktivitas,
		FieldCreatedBy: rec.CreatedBy,
		FieldUploader:  rec.CreatedBy,
	} {
		idx := cm.Index[field]
		if vals[idx] != want {
			t.Errorf("kolom %q: %q, harus %q", legacyHeaders[idx], vals[idx], want)
		}
	}
}

func TestRekrutmenRowValuesKeepsUnknownColumns(t *testing.T) {
	cm, err := ResolveColumnMap(legacyHeaders, AppColumnSpec())
	if err != nil {
		t.Fatalf("ResolveColumnMap: %v", err)
	}
	existing := sheetRow{
		"NAMA KECAMATAN": "Serpong",
		"UPLM":          "UPLM",
		"NAMA ANGGOTA":  "Kader Serpong",
	}
	rec := model.Rekrutmen{
		ID:        "abc",
		Nama:      "Uju Juhiwa",
		CreatedBy: "umar",
	}
	vals := rekrutmenRowValues(cm, rec, existing)
	for _, col := range []string{"NAMA KECAMATAN", "UPLM", "NAMA ANGGOTA"} {
		idx := cm.Index[""]
		for i, h := range legacyHeaders {
			if h == col {
				idx = i
			}
		}
		if vals[idx] != existing[col] {
			t.Errorf("kolom %q: %q, harus tetap %q", col, vals[idx], existing[col])
		}
	}
}

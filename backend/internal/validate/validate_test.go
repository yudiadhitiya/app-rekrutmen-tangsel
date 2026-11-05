package validate

import (
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
)

func valid() *model.Rekrutmen {
	return &model.Rekrutmen{
		Nama:              "Ahmad Fauzi",
		NoWA:              "628123456789",
		RT:                "003",
		RW:                "005",
		Kecamatan:         "Serpong",
		Kelurahan:         "Buaran",
		Aktivitas:         "Wiraswasta",
		TingkatPenerimaan: "Terbuka",
		Metode:            "Halaqoh",
		Status:            "Berhasil",
		UPL:               "UPLS",
	}
}

var aktivitasNames = []string{"Wiraswasta", "Pedagang"}

func TestValidRecordHasNoErrors(t *testing.T) {
	f := Rekrutmen(valid(), aktivitasNames)
	if len(f) != 0 {
		t.Fatalf("expected no errors, got %v", f)
	}
}

func TestRequiredFields(t *testing.T) {
	cases := []struct {
		field  string
		mutate func(*model.Rekrutmen)
	}{
		{"nama", func(r *model.Rekrutmen) { r.Nama = "" }},
		{"no_wa", func(r *model.Rekrutmen) { r.NoWA = "" }},
		{"kecamatan", func(r *model.Rekrutmen) { r.Kecamatan = "" }},
		{"aktivitas", func(r *model.Rekrutmen) { r.Aktivitas = "" }},
		{"tingkat_penerimaan", func(r *model.Rekrutmen) { r.TingkatPenerimaan = "" }},
		{"metode", func(r *model.Rekrutmen) { r.Metode = "" }},
		{"status", func(r *model.Rekrutmen) { r.Status = "" }},
	}
	for _, c := range cases {
		r := valid()
		c.mutate(r)
		f := Rekrutmen(r, aktivitasNames)
		if _, ok := f[c.field]; !ok {
			t.Errorf("expected error on %q, got %v", c.field, f)
		}
	}
}

// PRD: "Nomor Whatsapp ... menolak karakter selain digit" and
// "Nomor Whatsapp ... simpan sebagai string agar angka 0 di awal tidak hilang".
func TestWhatsAppMustBeDigitsOnly(t *testing.T) {
	for _, bad := range []string{"6281abc", "+628123", "62 81", "62-81"} {
		r := valid()
		r.NoWA = bad
		f := Rekrutmen(r, aktivitasNames)
		if _, ok := f["no_wa"]; !ok {
			t.Errorf("expected error for no_wa=%q", bad)
		}
	}
}

func TestWhatsAppKeepsLeadingZero(t *testing.T) {
	r := valid()
	r.NoWA = "08123456789"
	Rekrutmen(r, aktivitasNames)
	if r.NoWA != "08123456789" {
		t.Errorf("leading zero lost: got %q", r.NoWA)
	}
}

// PRD: "RT/RW menolak karakter selain digit".
func TestRTAndRWMustBeDigitsOnly(t *testing.T) {
	for _, field := range []string{"rt", "rw"} {
		r := valid()
		if field == "rt" {
			r.RT = "00A"
		} else {
			r.RW = "1x"
		}
		f := Rekrutmen(r, aktivitasNames)
		if _, ok := f[field]; !ok {
			t.Errorf("expected error on %q", field)
		}
	}
}

func TestRTAndRWMayBeEmpty(t *testing.T) {
	r := valid()
	r.RT = ""
	r.RW = ""
	if f := Rekrutmen(r, aktivitasNames); len(f) != 0 {
		t.Fatalf("expected no errors, got %v", f)
	}
}

// PRD: "Kecamatan wajib dan hanya menampilkan 7 kecamatan Tangsel".
func TestKecamatanMustBeOneOfSeven(t *testing.T) {
	r := valid()
	r.Kecamatan = "Jakarta Selatan"
	f := Rekrutmen(r, aktivitasNames)
	if _, ok := f["kecamatan"]; !ok {
		t.Errorf("expected error for invalid kecamatan, got %v", f)
	}
}

// PRD: "Kelurahan berubah mengikuti kecamatan yang dipilih".
func TestKelurahanMustBelongToKecamatan(t *testing.T) {
	r := valid()
	r.Kecamatan = "Ciputat"
	f := Rekrutmen(r, aktivitasNames)
	if _, ok := f["kelurahan"]; !ok {
		t.Errorf("expected error for kelurahan not in kecamatan, got %v", f)
	}
}

func TestKelurahanMayBeEmpty(t *testing.T) {
	r := valid()
	r.Kelurahan = ""
	if f := Rekrutmen(r, aktivitasNames); len(f) != 0 {
		t.Fatalf("expected no errors, got %v", f)
	}
}

func TestAktivitasMustExistInMaster(t *testing.T) {
	r := valid()
	r.Aktivitas = "Programmer"
	f := Rekrutmen(r, aktivitasNames)
	if _, ok := f["aktivitas"]; !ok {
		t.Errorf("expected error for unknown aktivitas, got %v", f)
	}
}

func TestEnumValuesMustBeValid(t *testing.T) {
	r := valid()
	r.TingkatPenerimaan = "Sangat Buka"
	r.Metode = "Manual"
	r.Status = "Selesai"
	f := Rekrutmen(r, aktivitasNames)
	for _, field := range []string{"tingkat_penerimaan", "metode", "status"} {
		if _, ok := f[field]; !ok {
			t.Errorf("expected error on %q, got %v", field, f)
		}
	}
}

// The UPL field is a single optional radio choice.
func TestUPLAcceptsEmptyUPLSorUPLM(t *testing.T) {
	for _, v := range []string{"", "UPLS", "UPLM"} {
		r := valid()
		r.UPL = v
		if f := Rekrutmen(r, aktivitasNames); len(f) != 0 {
			t.Errorf("upl=%q should be valid, got %v", v, f)
		}
	}
	r := valid()
	r.UPL = "UPLX"
	if f := Rekrutmen(r, aktivitasNames); len(f) == 0 {
		t.Errorf("expected error for invalid upl")
	}
}

func TestFieldsAreTrimmed(t *testing.T) {
	r := valid()
	r.Nama = "  Ahmad Fauzi  "
	r.NoWA = " 628123456789 "
	Rekrutmen(r, aktivitasNames)
	if r.Nama != "Ahmad Fauzi" {
		t.Errorf("nama not trimmed: %q", r.Nama)
	}
	if r.NoWA != "628123456789" {
		t.Errorf("no_wa not trimmed: %q", r.NoWA)
	}
}

package api

import (
	"net/http"
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
)

// uploadRow builds one row of an upload with the headers of the template of
// the previous application.
func uploadRow(nama, noWA string) map[string]string {
	return map[string]string{
		"NAMA TARGET":        nama,
		"NO WHATSAPP":        noWA,
		"ALAMAT":             "Jl. Melati No. 5",
		"RT":                 "003",
		"RW":                 "005",
		"KECAMATAN":          "Serpong",
		"KELURAHAN":          "Buaran",
		"AKTIFITAS":          "Wiraswasta",
		"TINGKAT PENERIMAAN": "Terbuka",
		"METODE":             "Halaqoh",
		"STATUS":             "Berhasil",
		"KETERANGAN":         "upload",
	}
}

// uploadBatch posts rows to the batch endpoint and returns the status code
// and the decoded answer.
func (e *testEnv) uploadBatch(rows []map[string]string) (int, map[string]any) {
	e.t.Helper()
	var out map[string]any
	code := e.request(http.MethodPost, "/api/rekrutmen/batch", rows, &out)
	return code, out
}

func TestMemberCanUploadBatch(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("petugas", "petugas123")
	code, out := e.uploadBatch([]map[string]string{
		uploadRow("Ahmad Fauzi", "628123456781"),
		uploadRow("Siti Aminah", "628123456782"),
	})
	if code != http.StatusCreated {
		t.Fatalf("upload: expected 201, got %d (%v)", code, out)
	}
	if created, _ := out["created"].(float64); created != 2 {
		t.Fatalf("created: expected 2, got %v", out["created"])
	}

	// The records belong to the member that uploaded them.
	rows, total := e.listRekrutmen()
	if total != 2 || len(rows) != 2 {
		t.Fatalf("list: expected 2 records, got %d (total %d)", len(rows), total)
	}
	for _, rec := range rows {
		if rec.CreatedBy != "petugas" {
			t.Errorf("created_by: expected petugas, got %q", rec.CreatedBy)
		}
		if rec.Nama != "Ahmad Fauzi" && rec.Nama != "Siti Aminah" {
			t.Errorf("nama: got %q", rec.Nama)
		}
	}
}

func TestBatchIgnoresTheOwnerColumnsOfTheFile(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("petugas", "petugas123")
	row := uploadRow("Ahmad Fauzi", "628123456783")
	row["NAMA ANGGOTA"] = "Orang Lain"
	row["NAMA KECAMATAN"] = "Ciputat"
	row["CREATED_BY"] = "admin"

	code, out := e.uploadBatch([]map[string]string{row})
	if code != http.StatusCreated {
		t.Fatalf("upload: expected 201, got %d (%v)", code, out)
	}

	rows, _ := e.listRekrutmen()
	if len(rows) != 1 {
		t.Fatalf("list: expected 1 record, got %d", len(rows))
	}
	if rows[0].CreatedBy != "petugas" {
		t.Errorf("created_by: expected petugas, got %q", rows[0].CreatedBy)
	}
}

func TestBatchSavesNothingWhenOneRowIsWrong(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("petugas", "petugas123")
	bad := uploadRow("Salah", "bukan-angka")
	code, out := e.uploadBatch([]map[string]string{
		uploadRow("Ahmad Fauzi", "628123456784"),
		bad,
	})
	if code != http.StatusBadRequest {
		t.Fatalf("upload: expected 400, got %d (%v)", code, out)
	}
	if out["error"] == nil {
		t.Error(" jawaban harus menyebut kegagalannya")
	}
	if _, ok := out["errors"].([]any); !ok {
		t.Errorf("jawaban harus memuat daftar error baris: %v", out)
	}

	// Nothing is saved when one row is wrong.
	rows, total := e.listRekrutmen()
	if total != 0 || len(rows) != 0 {
		t.Fatalf("list: expected 0 records, got %d (total %d)", len(rows), total)
	}
}

func TestBatchRejectsTwoRowsWithTheSameNumber(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("petugas", "petugas123")
	code, _ := e.uploadBatch([]map[string]string{
		uploadRow("Ahmad Fauzi", "628123456785"),
		uploadRow("Siti Aminah", "628123456785"),
	})
	if code != http.StatusBadRequest {
		t.Fatalf("upload: expected 400, got %d", code)
	}

	rows, total := e.listRekrutmen()
	if total != 0 || len(rows) != 0 {
		t.Fatalf("list: expected 0 records, got %d (total %d)", len(rows), total)
	}
}

func TestBatchSkipsANumberThatIsAlreadyStored(t *testing.T) {
	e := newTestEnv(t)

	// The record is stored by the full administrator.
	e.loginAs("admin", "admin123")
	var rec model.Rekrutmen
	if code := e.request(http.MethodPost, "/api/rekrutmen", validBody(), &rec); code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", code)
	}

	// The member uploads the same number again.
	e.loginAs("petugas", "petugas123")
	code, out := e.uploadBatch([]map[string]string{
		uploadRow("Ahmad Fauzi", rec.NoWA),
		uploadRow("Siti Aminah", "628123456786"),
	})
	if code != http.StatusCreated {
		t.Fatalf("upload: expected 201, got %d (%v)", code, out)
	}
	if created, _ := out["created"].(float64); created != 1 {
		t.Errorf("created: expected 1, got %v", out["created"])
	}
	skipped, _ := out["skipped"].([]any)
	if len(skipped) != 1 {
		t.Fatalf("skipped: expected 1 row, got %v", out["skipped"])
	}
	first, _ := skipped[0].(map[string]any)
	if reason, _ := first["reason"].(string); reason == "" {
		t.Errorf("baris yang dilewati harus menjelaskan sebabnya: %v", first)
	}
}

func TestBatchIsOnlyForMembers(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("admin", "admin123")
	code, _ := e.uploadBatch([]map[string]string{uploadRow("Ahmad Fauzi", "628123456787")})
	if code != http.StatusForbidden {
		t.Fatalf("upload as admin: expected 403, got %d", code)
	}
}

func TestBatchRejectsMoreThanTheAllowedRows(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("petugas", "petugas123")
	rows := make([]map[string]string, 0, MaxBatchRows+1)
	for i := 0; i <= MaxBatchRows; i++ {
		rows = append(rows, uploadRow("Ahmad Fauzi", "62812345678"+string(rune('a'+i%26))+string(rune('a'+i/26))))
	}
	code, _ := e.uploadBatch(rows)
	if code != http.StatusBadRequest {
		t.Fatalf("upload: expected 400, got %d", code)
	}
}

func TestBatchAcceptsTheSpellingsOfThePreviousApplication(t *testing.T) {
	e := newTestEnv(t)

	// The previous application spelled two Metode values differently
	// (old-app/index.html:306), so a file that was filled there still
	// uploads here.
	e.loginAs("petugas", "petugas123")
	row := uploadRow("Ahmad Fauzi", "628123456789")
	row["METODE"] = "Dakwah Fardliyah"
	code, out := e.uploadBatch([]map[string]string{row})
	if code != http.StatusCreated {
		t.Fatalf("upload: expected 201, got %d (%v)", code, out)
	}

	rows, _ := e.listRekrutmen()
	if len(rows) != 1 {
		t.Fatalf("list: expected 1 record, got %d", len(rows))
	}
	if rows[0].Metode != "Dakwah Fardiyah" {
		t.Errorf("metode: expected Dakwah Fardiyah, got %q", rows[0].Metode)
	}
}

func TestBatchAcceptsTheFieldNamesOfThisApplication(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("petugas", "petugas123")
	var out map[string]any
	code := e.request(http.MethodPost, "/api/rekrutmen/batch", []map[string]any{
		{
			"nama":               "Ahmad Fauzi",
			"no_wa":              "628123456788",
			"kecamatan":          "Serpong",
			"kelurahan":          "Buaran",
			"aktivitas":          "Wiraswasta",
			"tingkat_penerimaan": "Terbuka",
			"metode":             "Halaqoh",
			"status":             "Berhasil",
			"upl":                "UPLS",
		},
	}, &out)
	if code != http.StatusCreated {
		t.Fatalf("upload: expected 201, got %d (%v)", code, out)
	}
	rows, _ := e.listRekrutmen()
	if len(rows) != 1 {
		t.Fatalf("list: expected 1 record, got %d", len(rows))
	}
	if rows[0].UPL != "UPLS" {
		t.Errorf("upl: expected UPLS, got %q", rows[0].UPL)
	}
}

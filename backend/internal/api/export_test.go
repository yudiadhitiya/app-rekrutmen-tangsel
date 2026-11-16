package api

import (
	"fmt"
	"net/http"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
)

// exportResponse is the answer of the export endpoint.
type exportResponse struct {
	Headers  []string   `json:"headers"`
	Rows     [][]string `json:"rows"`
	Total    int        `json:"total"`
	Filename string     `json:"filename"`
}

// exportRekrutmen downloads the records and returns the answer.
func (e *testEnv) exportRekrutmen(query string) exportResponse {
	e.t.Helper()
	var out exportResponse
	path := "/api/rekrutmen/export"
	if query != "" {
		path += "?" + query
	}
	if code := e.request(http.MethodGet, path, nil, &out); code != http.StatusOK {
		e.t.Fatalf("export: expected 200, got %d", code)
	}
	return out
}

// createRecord stores a record with the given district and status.
func (e *testEnv) createRecord(noWA, kecamatan, status string) {
	e.t.Helper()
	body := validBody()
	body["no_wa"] = noWA
	body["kecamatan"] = kecamatan
	body["kelurahan"] = "Buaran"
	if kecamatan == "Ciputat" {
		body["kelurahan"] = "Cipayung"
	}
	body["status"] = status
	if code := e.request(http.MethodPost, "/api/rekrutmen", body, nil); code != http.StatusCreated {
		e.t.Fatalf("create %s: expected 201, got %d", noWA, code)
	}
}

func TestExportRequiresAuth(t *testing.T) {
	e := newTestEnv(t)
	if code := e.request(http.MethodGet, "/api/rekrutmen/export", nil, nil); code != http.StatusUnauthorized {
		t.Fatalf("export without session: expected 401, got %d", code)
	}
}

func TestExportMemberHoldsOnlyOwnRecords(t *testing.T) {
	e := newTestEnv(t)

	// The full administrator creates a record.
	e.loginAs("admin", "admin123")
	e.createRecord("628111100001", "Serpong", "Berhasil")

	// The member of level 1 creates a record of its own.
	e.loginAs("petugas", "petugas123")
	e.createRecord("628111100002", "Serpong", "Berhasil")

	// The member downloads only its own record.
	out := e.exportRekrutmen("")
	if len(out.Rows) != 1 {
		t.Fatalf("member export: expected 1 row, got %d", len(out.Rows))
	}
	if out.Rows[0][0] != "petugas" {
		t.Errorf("member export: expected owner petugas, got %q", out.Rows[0][0])
	}

	// The full administrator downloads both records.
	e.loginAs("admin", "admin123")
	out = e.exportRekrutmen("")
	if len(out.Rows) != 2 {
		t.Fatalf("admin export: expected 2 rows, got %d", len(out.Rows))
	}
	if out.Total != 2 {
		t.Errorf("admin export: expected total 2, got %d", out.Total)
	}
}

func TestExportAdminKecamatanHoldsOnlyOwnDistrict(t *testing.T) {
	e := newTestEnvWithUsers(t, []auth.User{
		newTestUser(t, "kec", "kec12345", model.LevelAdminKecamatan, "Serpong"),
	})

	e.loginAs("admin", "admin123")
	e.createRecord("628222100001", "Serpong", "Berhasil")
	e.createRecord("628222100002", "Ciputat", "Berhasil")

	e.loginAs("kec", "kec12345")
	out := e.exportRekrutmen("")
	if len(out.Rows) != 1 {
		t.Fatalf("admin kecamatan export: expected 1 row, got %d", len(out.Rows))
	}
	if out.Rows[0][6] != "Serpong" {
		t.Errorf("admin kecamatan export: expected district Serpong, got %q", out.Rows[0][6])
	}
}

func TestExportFollowsActiveFilter(t *testing.T) {
	e := newTestEnv(t)
	e.loginAs("admin", "admin123")
	e.createRecord("628333100001", "Serpong", "Berhasil")
	e.createRecord("628333100002", "Ciputat", "Proses")

	out := e.exportRekrutmen("kecamatan=Serpong")
	if len(out.Rows) != 1 {
		t.Fatalf("export with district filter: expected 1 row, got %d", len(out.Rows))
	}
	if out.Rows[0][6] != "Serpong" {
		t.Errorf("export with district filter: expected Serpong, got %q", out.Rows[0][6])
	}

	out = e.exportRekrutmen("status=Proses")
	if len(out.Rows) != 1 {
		t.Fatalf("export with status filter: expected 1 row, got %d", len(out.Rows))
	}
	if out.Rows[0][13] != "Proses" {
		t.Errorf("export with status filter: expected status Proses, got %q", out.Rows[0][13])
	}
}

func TestExportIsNotPaginated(t *testing.T) {
	e := newTestEnv(t)
	e.loginAs("admin", "admin123")
	for i := 0; i < 30; i++ {
		e.createRecord(fmt.Sprintf("6284441000%02d", i), "Serpong", "Berhasil")
	}

	out := e.exportRekrutmen("")
	if len(out.Rows) != 30 {
		t.Fatalf("export: expected 30 rows, got %d", len(out.Rows))
	}
	if out.Total != 30 {
		t.Errorf("export: expected total 30, got %d", out.Total)
	}
	if len(out.Headers) != len(out.Rows[0]) {
		t.Errorf("export: expected %d headers, got %d", len(out.Rows[0]), len(out.Headers))
	}
}

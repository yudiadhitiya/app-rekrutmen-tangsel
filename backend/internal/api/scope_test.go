package api

import (
	"net/http"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
)

// loginAs signs in with a fresh session, so a test can act as more than one
// account within the same environment.
func (e *testEnv) loginAs(username, password string) {
	e.t.Helper()
	jar, err := cookieJar()
	if err != nil {
		e.t.Fatal(err)
	}
	e.client = &http.Client{Jar: jar}
	if code := e.login(username, password); code != http.StatusOK {
		e.t.Fatalf("login %s: expected 200, got %d", username, code)
	}
}

// listRekrutmen returns the records and the total the list endpoint reports.
func (e *testEnv) listRekrutmen() ([]model.Rekrutmen, int) {
	e.t.Helper()
	var out struct {
		Total int               `json:"total"`
		Data  []model.Rekrutmen `json:"data"`
	}
	if code := e.request(http.MethodGet, "/api/rekrutmen?limit=100", nil, &out); code != http.StatusOK {
		e.t.Fatalf("list: expected 200, got %d", code)
	}
	return out.Data, out.Total
}

func TestAdminKecamatanOnlySeesOwnDistrict(t *testing.T) {
	e := newTestEnvWithUsers(t, []auth.User{
		newTestUser(t, "kec", "kec12345", model.LevelAdminKecamatan, "Serpong"),
	})

	// The full administrator creates a record of every district.
	e.loginAs("admin", "admin123")
	for i, kec := range []string{"Serpong", "Ciputat"} {
		body := validBody()
		body["kecamatan"] = kec
		body["kelurahan"] = "Buaran"
		if kec == "Ciputat" {
			body["kelurahan"] = "Cipayung"
		}
		body["no_wa"] = "62812345678" + string(rune('0'+i))
		if code := e.request(http.MethodPost, "/api/rekrutmen", body, nil); code != http.StatusCreated {
			t.Fatalf("create %s: expected 201, got %d", kec, code)
		}
	}

	// The district administrator only reaches its own district.
	e.loginAs("kec", "kec12345")
	rows, total := e.listRekrutmen()
	if total != 1 || len(rows) != 1 {
		t.Fatalf("admin kecamatan list: expected 1 record, got %d (total %d)", len(rows), total)
	}
	if rows[0].Kecamatan != "Serpong" {
		t.Errorf("admin kecamatan list: expected district Serpong, got %q", rows[0].Kecamatan)
	}

	// The full administrator still reaches every district.
	e.loginAs("admin", "admin123")
	rows, total = e.listRekrutmen()
	if total != 2 || len(rows) != 2 {
		t.Fatalf("admin list: expected 2 records, got %d (total %d)", len(rows), total)
	}
}

func TestAdminKecamatanCannotOpenAnotherDistrictRecord(t *testing.T) {
	e := newTestEnvWithUsers(t, []auth.User{
		newTestUser(t, "kec", "kec12345", model.LevelAdminKecamatan, "Serpong"),
	})

	e.loginAs("admin", "admin123")
	body := validBody()
	body["kecamatan"] = "Ciputat"
	body["kelurahan"] = "Cipayung"
	var rec model.Rekrutmen
	if code := e.request(http.MethodPost, "/api/rekrutmen", body, &rec); code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", code)
	}

	e.loginAs("kec", "kec12345")
	if code := e.request(http.MethodGet, "/api/rekrutmen/"+rec.ID, nil, nil); code != http.StatusNotFound {
		t.Fatalf("get other district record: expected 404, got %d", code)
	}
}

func TestPetugasOnlySeesOwnRecords(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("admin", "admin123")
	var adminRec model.Rekrutmen
	if code := e.request(http.MethodPost, "/api/rekrutmen", validBody(), &adminRec); code != http.StatusCreated {
		t.Fatalf("create as admin: expected 201, got %d", code)
	}

	e.loginAs("petugas", "petugas123")
	body := validBody()
	body["no_wa"] = "628999999999"
	body["nama"] = "Budi Santoso"
	var petugasRec model.Rekrutmen
	if code := e.request(http.MethodPost, "/api/rekrutmen", body, &petugasRec); code != http.StatusCreated {
		t.Fatalf("create as petugas: expected 201, got %d", code)
	}

	// The officer sees the record they created, and only that one.
	rows, total := e.listRekrutmen()
	if total != 1 || len(rows) != 1 {
		t.Fatalf("petugas list: expected 1 record, got %d (total %d)", len(rows), total)
	}
	if rows[0].CreatedBy != "petugas" {
		t.Errorf("petugas list: expected own record, got created_by %q", rows[0].CreatedBy)
	}

	// The admin sees both records.
	e.loginAs("admin", "admin123")
	rows, total = e.listRekrutmen()
	if total != 2 || len(rows) != 2 {
		t.Fatalf("admin list: expected 2 records, got %d (total %d)", len(rows), total)
	}
}

func TestPetugasCannotOpenAnotherAccountRecord(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("admin", "admin123")
	var adminRec model.Rekrutmen
	if code := e.request(http.MethodPost, "/api/rekrutmen", validBody(), &adminRec); code != http.StatusCreated {
		t.Fatalf("create as admin: expected 201, got %d", code)
	}

	e.loginAs("petugas", "petugas123")
	if code := e.request(http.MethodGet, "/api/rekrutmen/"+adminRec.ID, nil, nil); code != http.StatusNotFound {
		t.Fatalf("get other record: expected 404, got %d", code)
	}

	body := validBody()
	body["nama"] = "Diubah Petugas"
	if code := e.request(http.MethodPut, "/api/rekrutmen/"+adminRec.ID, body, nil); code != http.StatusNotFound {
		t.Fatalf("update other record: expected 404, got %d", code)
	}
}

func TestOnlyAdminCanDeleteRecord(t *testing.T) {
	e := newTestEnv(t)

	e.loginAs("admin", "admin123")
	var rec model.Rekrutmen
	if code := e.request(http.MethodPost, "/api/rekrutmen", validBody(), &rec); code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", code)
	}

	e.loginAs("petugas", "petugas123")
	if code := e.request(http.MethodDelete, "/api/rekrutmen/"+rec.ID, nil, nil); code != http.StatusForbidden {
		t.Fatalf("delete as petugas: expected 403, got %d", code)
	}

	e.loginAs("admin", "admin123")
	if code := e.request(http.MethodDelete, "/api/rekrutmen/"+rec.ID, nil, nil); code != http.StatusOK {
		t.Fatalf("delete as admin: expected 200, got %d", code)
	}
}

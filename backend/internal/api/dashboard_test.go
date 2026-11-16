package api

import (
	"net/http"
	"testing"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/seed"
)

// dashboardResponse is the shape of the dashboard endpoint.
type dashboardResponse struct {
	Scope   string `json:"scope"`
	Filter  string `json:"filter"`
	Total   int    `json:"total"`
	Methods []struct {
		Name  string `json:"name"`
		Total int    `json:"total"`
	} `json:"methods"`
	MethodSummary []struct {
		Name         string `json:"name"`
		Total        int    `json:"total"`
		Berhasil     int    `json:"berhasil"`
		Proses       int    `json:"proses"`
		GagalPending int    `json:"gagalPending"`
	} `json:"methodSummary"`
	Acceptance []struct {
		Name  string `json:"name"`
		Total int    `json:"total"`
	} `json:"acceptance"`
	Status []struct {
		Name  string `json:"name"`
		Total int    `json:"total"`
	} `json:"status"`
	KecamatanSummary []struct {
		Kecamatan string `json:"kecamatan"`
		Target    int    `json:"target"`
		Berhasil  int    `json:"berhasil"`
	} `json:"kecamatanSummary"`
	Top10 []struct {
		Nama  string `json:"nama"`
		Total int    `json:"total"`
	} `json:"top10"`
}

// dashboard calls the dashboard endpoint and returns the answer.
func (e *testEnv) dashboard(query string) dashboardResponse {
	e.t.Helper()
	var out dashboardResponse
	if code := e.request(http.MethodGet, "/api/dashboard"+query, nil, &out); code != http.StatusOK {
		e.t.Fatalf("dashboard: expected 200, got %d", code)
	}
	return out
}

// createRecord stores a record with the given district and method.
func (e *testEnv) createDashboardRecord(noWA, kecamatan, metode, status string) {
	e.t.Helper()
	body := validBody()
	body["no_wa"] = noWA
	body["nama"] = "Pengguna " + noWA
	body["kecamatan"] = kecamatan
	body["kelurahan"] = ""
	if villages := seed.KelurahanByKecamatan(kecamatan); len(villages) > 0 {
		body["kelurahan"] = villages[0].Nama
	}
	body["metode"] = metode
	body["status"] = status
	if code := e.request(http.MethodPost, "/api/rekrutmen", body, nil); code != http.StatusCreated {
		e.t.Fatalf("create %s: expected 201, got %d", noWA, code)
	}
}

// dashboardTarget counts the total of a breakdown list.
func dashboardTotal(items []struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}, name string) int {
	for _, item := range items {
		if item.Name == name {
			return item.Total
		}
	}
	return 0
}

func TestDashboardCountsEveryRecordForFullAdmin(t *testing.T) {
	e := newTestEnv(t)
	e.loginAs("admin", "admin123")

	e.createDashboardRecord("62811111111", "Serpong", "Halaqoh", "Berhasil")
	e.createDashboardRecord("62822222222", "Ciputat", "Program Antara", "Proses")

	data := e.dashboard("")
	if data.Total != 2 {
		t.Fatalf("dashboard total: expected 2, got %d", data.Total)
	}
	if got := dashboardTotal(data.Methods, "Halaqoh"); got != 1 {
		t.Errorf("dashboard methods Halaqoh: expected 1, got %d", got)
	}
	if got := dashboardTotal(data.Methods, "Program Antara"); got != 1 {
		t.Errorf("dashboard methods Program Antara: expected 1, got %d", got)
	}
	if got := dashboardTotal(data.Status, "Berhasil"); got != 1 {
		t.Errorf("dashboard status Berhasil: expected 1, got %d", got)
	}
	// Every district of Tangsel is listed, so the panel keeps its layout.
	if len(data.KecamatanSummary) != 7 {
		t.Errorf("dashboard kecamatan summary: expected 7 districts, got %d", len(data.KecamatanSummary))
	}
	for _, item := range data.KecamatanSummary {
		if item.Kecamatan == "Serpong" && item.Target != 1 {
			t.Errorf("dashboard Serpong target: expected 1, got %d", item.Target)
		}
		if item.Kecamatan == "Ciputat" && item.Target != 1 {
			t.Errorf("dashboard Ciputat target: expected 1, got %d", item.Target)
		}
	}
	if len(data.Top10) != 2 {
		t.Errorf("dashboard top10: expected 2 entries, got %d", len(data.Top10))
	}
}

func TestDashboardDistrictAdminOnlyCountsOwnDistrict(t *testing.T) {
	e := newTestEnvWithUsers(t, []auth.User{
		newTestUser(t, "kec", "kec12345", model.LevelAdminKecamatan, "Serpong"),
	})

	e.loginAs("admin", "admin123")
	e.createDashboardRecord("62811111111", "Serpong", "Halaqoh", "Berhasil")
	e.createDashboardRecord("62822222222", "Ciputat", "Halaqoh", "Berhasil")

	e.loginAs("kec", "kec12345")
	data := e.dashboard("")
	if data.Total != 1 {
		t.Fatalf("dashboard total: expected 1, got %d", data.Total)
	}
	for _, item := range data.KecamatanSummary {
		if item.Kecamatan == "Ciputat" && item.Target != 0 {
			t.Errorf("dashboard Ciputat target: expected 0, got %d", item.Target)
		}
	}
}

func TestDashboardMemberFollowsTheChosenFilter(t *testing.T) {
	e := newTestEnvWithUsers(t, []auth.User{
		newTestUser(t, "anggota", "anggota123", model.LevelAnggota, "Serpong"),
	})

	// The administrator adds one record of another district.
	e.loginAs("admin", "admin123")
	e.createDashboardRecord("62811111111", "Ciputat", "Halaqoh", "Berhasil")

	// The member adds one record of its own district.
	e.loginAs("anggota", "anggota123")
	e.createDashboardRecord("62822222222", "Serpong", "Halaqoh", "Berhasil")

	// Every district, like in the previous application.
	if data := e.dashboard("?filter=ALL"); data.Total != 2 {
		t.Errorf("dashboard filter ALL: expected 2, got %d", data.Total)
	}
	// The district of the account.
	if data := e.dashboard("?filter=KECAMATAN"); data.Total != 1 {
		t.Errorf("dashboard filter KECAMATAN: expected 1, got %d", data.Total)
	}
	// The records the account created itself.
	if data := e.dashboard("?filter=ANGGOTA"); data.Total != 1 {
		t.Errorf("dashboard filter ANGGOTA: expected 1, got %d", data.Total)
	}
}

func TestDashboardMethodSummaryCountsTheStatuses(t *testing.T) {
	e := newTestEnv(t)
	e.loginAs("admin", "admin123")

	e.createDashboardRecord("62811111111", "Serpong", "Halaqoh", "Berhasil")
	e.createDashboardRecord("62822222222", "Serpong", "Halaqoh", "Proses")
	e.createDashboardRecord("62833333333", "Serpong", "Halaqoh", "Gagal")

	data := e.dashboard("")
	if len(data.MethodSummary) == 0 {
		t.Fatal("dashboard method summary: expected rows, got none")
	}
	for _, item := range data.MethodSummary {
		if item.Name != "Halaqoh" {
			continue
		}
		if item.Total != 3 || item.Berhasil != 1 || item.Proses != 1 || item.GagalPending != 1 {
			t.Errorf("dashboard method summary Halaqoh: unexpected counts %+v", item)
		}
		return
	}
	t.Error("dashboard method summary: the row Halaqoh is missing")
}

func TestDashboardRequiresAuthentication(t *testing.T) {
	e := newTestEnv(t)
	if code := e.request(http.MethodGet, "/api/dashboard", nil, nil); code != http.StatusUnauthorized {
		t.Fatalf("dashboard without session: expected 401, got %d", code)
	}
}

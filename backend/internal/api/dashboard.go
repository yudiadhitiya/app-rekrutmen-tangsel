package api

import (
	"math"
	"net/http"
	"sort"
	"strings"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/seed"
)

// The dashboard is the high level view of the recruitment data. It was part
// of the previous application and is restored here, so the managers keep the
// same numbers they used to see.
//
// The endpoint only returns aggregates. It never returns a single record, so
// it cannot be used to reach the data of another account.

// The dashboard filters of the previous application.
const (
	// dashboardFilterAll aggregates every record of the sheet.
	dashboardFilterAll = "ALL"
	// dashboardFilterKecamatan aggregates the records of the district of the
	// account.
	dashboardFilterKecamatan = "KECAMATAN"
	// dashboardFilterAnggota aggregates the records the account created.
	dashboardFilterAnggota = "ANGGOTA"
)

// dashboardPageSize is the number of records read per call while the
// dashboard collects every record it aggregates.
const dashboardPageSize = 500

// dashboardItem is one entry of a breakdown list.
type dashboardItem struct {
	Name    string  `json:"name"`
	Total   int     `json:"total"`
	Percent float64 `json:"percent"`
}

// dashboardMethodSummary is one row of the method x status summary.
type dashboardMethodSummary struct {
	Name         string  `json:"name"`
	Total        int     `json:"total"`
	Berhasil     int     `json:"berhasil"`
	Proses       int     `json:"proses"`
	GagalPending int     `json:"gagalPending"`
	Percent      float64 `json:"percent"`
}

// dashboardKecamatanSummary is one row of the summary per district.
type dashboardKecamatanSummary struct {
	Kecamatan      string  `json:"kecamatan"`
	Target         int     `json:"target"`
	Berhasil       int     `json:"berhasil"`
	Proses         int     `json:"proses"`
	PersenBerhasil float64 `json:"persenBerhasil"`
	PersenProses   float64 `json:"persenProses"`
}

// dashboardTopItem is one entry of the top ten list.
type dashboardTopItem struct {
	Nama  string `json:"nama"`
	Total int    `json:"total"`
}

// handleDashboard returns the aggregation of the recruitment data.
//
// The scope of the account decides which records are counted:
//
//   - a full administrator counts every record;
//   - a district administrator counts the records of its own district;
//   - a member counts the records of the chosen filter, which defaults to
//     every district. Like in the previous application, a member only sees
//     aggregates, never the records themselves.
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request, user model.User) {
	filter := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("filter")))
	switch filter {
	case dashboardFilterAll, dashboardFilterKecamatan, dashboardFilterAnggota:
	default:
		filter = dashboardFilterAll
	}

	records, err := s.allRekrutmen()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat data dashboard: "+err.Error())
		return
	}

	scoped := make([]model.Rekrutmen, 0, len(records))
	for _, record := range records {
		if dashboardRowAllowed(record, user, filter) {
			scoped = append(scoped, record)
		}
	}

	writeJSON(w, http.StatusOK, buildDashboard(user, scoped, filter))
}

// allRekrutmen reads every record of the store.
//
// The dashboard needs all of them, so the records are read page by page
// until the store reports no more.
func (s *Server) allRekrutmen() ([]model.Rekrutmen, error) {
	var out []model.Rekrutmen
	for page := 1; ; page++ {
		rows, total, err := s.Store().ListRekrutmen(model.Filter{Page: page, Limit: dashboardPageSize})
		if err != nil {
			return nil, err
		}
		out = append(out, rows...)
		if len(rows) == 0 || len(out) >= total {
			break
		}
	}
	return out, nil
}

// dashboardRowAllowed reports whether a record is counted for the account.
func dashboardRowAllowed(record model.Rekrutmen, user model.User, filter string) bool {
	switch user.LevelOrDefault() {
	case model.LevelAdmin:
		return true
	case model.LevelAdminKecamatan:
		// The scope of the account covers the district of the account. A
		// district administrator without a district only reaches the
		// records it created itself.
		return model.MatchScope(user.ScopeFor(), record)
	default:
		switch filter {
		case dashboardFilterAnggota:
			return strings.EqualFold(strings.TrimSpace(record.CreatedBy), strings.TrimSpace(user.Username))
		case dashboardFilterKecamatan:
			if user.Kecamatan == "" {
				return false
			}
			return strings.EqualFold(strings.TrimSpace(record.Kecamatan), strings.TrimSpace(user.Kecamatan))
		default:
			// Every district, but only as an aggregate.
			return true
		}
	}
}

// buildDashboard aggregates the records into the panels of the dashboard.
func buildDashboard(user model.User, records []model.Rekrutmen, filter string) map[string]any {
	total := len(records)

	methodCounts := map[string]int{}
	methodStatus := map[string]map[string]int{}
	otherMethodStatus := map[string]int{}
	acceptanceCounts := map[string]int{}
	statusCounts := map[string]int{}
	kecamatanCounts := map[string]*dashboardKecamatanSummary{}
	nameCounts := map[string]int{}

	// Every district of Tangsel is listed, so the panel keeps the same layout
	// no matter how much data there is.
	districts, err := seed.All()
	if err == nil {
		for _, district := range districts {
			kecamatanCounts[strings.ToLower(strings.TrimSpace(district.Nama))] = &dashboardKecamatanSummary{
				Kecamatan: district.Nama,
			}
		}
	}

	for _, record := range records {
		method := model.NormalizeMetode(strings.TrimSpace(record.Metode))
		methodCounts[method]++
		status := strings.TrimSpace(record.Status)
		if isKnownMethod(method) {
			if methodStatus[method] == nil {
				methodStatus[method] = map[string]int{}
			}
			methodStatus[method][status]++
		} else {
			otherMethodStatus[status]++
		}

		acceptanceCounts[strings.TrimSpace(record.TingkatPenerimaan)]++
		statusCounts[status]++

		name := strings.TrimSpace(record.Nama)
		if name == "" {
			name = "(Tanpa Nama)"
		}
		nameCounts[name]++

		key := strings.ToLower(strings.TrimSpace(record.Kecamatan))
		summary, ok := kecamatanCounts[key]
		if !ok {
			// A district outside the master list is not part of the panel.
			continue
		}
		summary.Target++
		switch status {
		case "Berhasil":
			summary.Berhasil++
		case "Proses", "Pending":
			summary.Proses++
		}
	}

	methods := dashboardBreakdown(model.Metode, methodCounts, total, "Lainnya / Belum Diisi")
	acceptance := dashboardBreakdown(model.TingkatPenerimaan, acceptanceCounts, total, "Belum Diisi / Lainnya")

	methodSummary := make([]dashboardMethodSummary, 0, len(methods))
	for _, method := range methods {
		counts := methodStatus[method.Name]
		if method.Name == "Lainnya / Belum Diisi" {
			counts = otherMethodStatus
		}
		methodSummary = append(methodSummary, dashboardMethodSummary{
			Name:         method.Name,
			Total:        method.Total,
			Berhasil:     counts["Berhasil"],
			Proses:       counts["Proses"],
			GagalPending: counts["Pending"] + counts["Gagal"],
			Percent:      method.Percent,
		})
	}

	status := make([]dashboardItem, 0, len(model.Status))
	for _, name := range model.Status {
		status = append(status, dashboardItem{
			Name:    name,
			Total:   statusCounts[name],
			Percent: dashboardPercent(statusCounts[name], total),
		})
	}

	statusDetail := make([]dashboardItem, 0, len(statusCounts))
	for name, count := range statusCounts {
		statusDetail = append(statusDetail, dashboardItem{
			Name:    name,
			Total:   count,
			Percent: dashboardPercent(count, total),
		})
	}
	sort.Slice(statusDetail, func(i, j int) bool {
		if statusDetail[i].Total != statusDetail[j].Total {
			return statusDetail[i].Total > statusDetail[j].Total
		}
		return statusDetail[i].Name < statusDetail[j].Name
	})

	kecamatanSummary := make([]dashboardKecamatanSummary, 0, len(kecamatanCounts))
	for _, summary := range kecamatanCounts {
		summary.PersenBerhasil = dashboardPercent(summary.Berhasil, summary.Target)
		summary.PersenProses = dashboardPercent(summary.Proses, summary.Target)
		kecamatanSummary = append(kecamatanSummary, *summary)
	}
	sort.Slice(kecamatanSummary, func(i, j int) bool {
		return kecamatanSummary[i].Kecamatan < kecamatanSummary[j].Kecamatan
	})

	top10 := make([]dashboardTopItem, 0, len(nameCounts))
	for name, count := range nameCounts {
		top10 = append(top10, dashboardTopItem{Nama: name, Total: count})
	}
	sort.Slice(top10, func(i, j int) bool {
		if top10[i].Total != top10[j].Total {
			return top10[i].Total > top10[j].Total
		}
		return top10[i].Nama < top10[j].Nama
	})
	if len(top10) > 10 {
		top10 = top10[:10]
	}

	return map[string]any{
		"role":             user.LevelOrDefault(),
		"scope":            dashboardScope(user, filter),
		"filter":           filter,
		"total":            total,
		"methods":          methods,
		"methodSummary":    methodSummary,
		"acceptance":       acceptance,
		"status":           status,
		"statusDetail":     statusDetail,
		"kecamatanSummary": kecamatanSummary,
		"top10":            top10,
	}
}

// dashboardBreakdown builds a breakdown list in the given order.
//
// A value outside the list is counted as the other label, so the total of the
// panel always matches the total of the dashboard.
func dashboardBreakdown(order []string, counts map[string]int, total int, otherLabel string) []dashboardItem {
	items := make([]dashboardItem, 0, len(order)+1)
	for _, name := range order {
		count := counts[name]
		items = append(items, dashboardItem{
			Name:    name,
			Total:   count,
			Percent: dashboardPercent(count, total),
		})
	}
	other := 0
	for name, count := range counts {
		if !containsName(order, name) {
			other += count
		}
	}
	if other > 0 {
		items = append(items, dashboardItem{
			Name:    otherLabel,
			Total:   other,
			Percent: dashboardPercent(other, total),
		})
	}
	return items
}

// isKnownMethod reports whether the method is one of the current options.
func isKnownMethod(name string) bool {
	for _, item := range model.Metode {
		if item == name {
			return true
		}
	}
	return false
}

// containsName reports whether the list holds the name.
func containsName(list []string, name string) bool {
	for _, item := range list {
		if item == name {
			return true
		}
	}
	return false
}

// dashboardPercent returns the share of the value in the total, with one
// decimal, the way the previous application calculated it.
func dashboardPercent(value, total int) float64 {
	if total <= 0 {
		return 0
	}
	return math.Round(float64(value)*1000/float64(total)) / 10
}

// dashboardScope describes the data the dashboard holds.
func dashboardScope(user model.User, filter string) string {
	switch user.LevelOrDefault() {
	case model.LevelAdmin:
		return "Seluruh data pendaftaran"
	case model.LevelAdminKecamatan:
		return "Data KECAMATAN " + orDash(user.Kecamatan)
	default:
		switch filter {
		case dashboardFilterAnggota:
			return "Dashboard data anggota Anda"
		case dashboardFilterKecamatan:
			return "Dashboard KECAMATAN " + orDash(user.Kecamatan)
		default:
			return "Dashboard seluruh kecamatan"
		}
	}
}

// orDash returns the value, or a dash when it is empty.
func orDash(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}

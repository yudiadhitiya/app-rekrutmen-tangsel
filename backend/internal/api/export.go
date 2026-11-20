package api

import (
	"net/http"
	"strconv"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
)

// maxExportRows caps the number of rows a single download may carry.
//
// The export follows the filter of the screen, so it stays small in normal
// use. The cap protects the server and the browser from a request that would
// ask for every record at once.
const maxExportRows = 10000

// exportHeaders are the columns of the downloaded file.
//
// The order follows the sheet of the previous application, so a file that is
// downloaded today keeps the same shape as a file from the old one.
var exportHeaders = []string{
	"NAMA ANGGOTA",
	"NAMA TARGET",
	"NO WHATSAPP",
	"ALAMAT",
	"RT",
	"RW",
	"KECAMATAN",
	"KELURAHAN",
	"UPLS",
	"UPLM",
	"AKTIFITAS",
	"TINGKAT PENERIMAAN",
	"METODE",
	"STATUS",
	"KETERANGAN",
	"CREATED_AT",
	"UPDATED_AT",
}

// exportRow turns a record into the cells of the downloaded file.
func exportRow(r model.Rekrutmen) []string {
	upls, uplm := "", ""
	if r.UPL == "UPLS" {
		upls = "UPLS"
	} else if r.UPL == "UPLM" {
		uplm = "UPLM"
	}
	return []string{
		r.CreatedBy,
		r.Nama,
		r.NoWA,
		r.Alamat,
		r.RT,
		r.RW,
		r.Kecamatan,
		r.Kelurahan,
		upls,
		uplm,
		r.Aktivitas,
		r.TingkatPenerimaan,
		r.Metode,
		r.Status,
		r.Keterangan,
		formatExportTime(r.CreatedAt),
		formatExportTime(r.UpdatedAt),
	}
}

// formatExportTime writes a timestamp the way a spreadsheet shows it.
func formatExportTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// rekrutmenFilterFromQuery reads the filter of the request.
//
// The scope of the account is part of the filter, so it is applied before the
// rows are selected. The listing and the export share this helper, so a
// download always holds the records the screen shows.
func rekrutmenFilterFromQuery(r *http.Request, user model.User) model.Filter {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	return model.Filter{
		Q:         q.Get("q"),
		Kecamatan: q.Get("kecamatan"),
		Kelurahan: q.Get("kelurahan"),
		Metode:    model.NormalizeMetode(q.Get("metode")),
		Status:    q.Get("status"),
		Page:      page,
		Limit:     limit,
		Scope:     user.ScopeFor(),
	}
}

// handleExportRekrutmen sends the records of the caller for download.
//
// Every account may use it, including a member of level 1: the scope of the
// account limits the export to the records the account may reach. The rows
// are not paginated, so the file holds the whole result of the filter.
func (s *Server) handleExportRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	f := rekrutmenFilterFromQuery(r, user)
	f.Page = 1
	f.Limit = maxExportRows

	data, total, err := s.Store().ListRekrutmen(f)
	if err != nil {
		writeStoreError(w, "gagal menyiapkan data unduhan", err)
		return
	}

	rows := make([][]string, 0, len(data))
	for _, rec := range data {
		rows = append(rows, exportRow(rec))
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"headers":  exportHeaders,
		"rows":     rows,
		"total":    total,
		"filename": exportFilename(),
	})
}

// exportFilename names the downloaded file with the current time.
func exportFilename() string {
	return "Data_Rekrutmen_Tangsel_" + time.Now().Format("20060102_150405") + ".xlsx"
}

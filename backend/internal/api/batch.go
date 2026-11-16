package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/validate"
)

// Limits of the batch upload.
//
// The previous application allowed 1.000 rows per upload
// (old-app/code.gs:739) and the file itself stays small, so the body limit
// only guards against a runaway request.
const (
	// MaxBatchRows is the number of records one upload may carry.
	MaxBatchRows = 1000
	// MaxBatchBody is the largest accepted request body.
	MaxBatchBody = 5 << 20
	// maxBatchErrors is the number of row errors reported at once.
	maxBatchErrors = 30
	// DefaultUploadMaxAttempts is the number of uploads one account may
	// start in ten minutes.
	DefaultUploadMaxAttempts = 5
)

// batchColumn maps a column header of the template to the field of the
// record.
//
// The headers of the template of the previous application are accepted, as
// well as the field names of this application. The headers that carry the
// owner of the record ("NAMA ANGGOTA", "NAMA KECAMATAN") are not mapped: the
// server fills them from the session, so a file cannot pretend to be
// another account.
var batchColumn = map[string]string{
	"nama":               "nama",
	"nama target":        "nama",
	"no_wa":              "no_wa",
	"no wa":              "no_wa",
	"no whatsapp":        "no_wa",
	"nomor whatsapp":     "no_wa",
	"alamat":             "alamat",
	"rt":                 "rt",
	"rw":                 "rw",
	"kecamatan":          "kecamatan",
	"kelurahan":          "kelurahan",
	"upl":                "upl",
	"upls":               "upls",
	"uplm":               "uplm",
	"aktivitas":          "aktivitas",
	"aktifitas":          "aktivitas",
	"tingkat penerimaan": "tingkat_penerimaan",
	"tingkat_penerimaan": "tingkat_penerimaan",
	"metode":             "metode",
	"status":             "status",
	"keterangan":         "keterangan",
}

// batchRowError is one problem in one row of the upload.
type batchRowError struct {
	// Row is the row number in the file, counting the header as row 1.
	Row int `json:"row"`
	// Field is the field the problem belongs to.
	Field string `json:"field,omitempty"`
	// Reason explains the problem.
	Reason string `json:"reason"`
}

// batchSkipped is a row that was valid but not saved.
type batchSkipped struct {
	// Row is the row number in the file, counting the header as row 1.
	Row int `json:"row"`
	// Nama is the name in the row.
	Nama string `json:"nama"`
	// NoWA is the WhatsApp number in the row.
	NoWA string `json:"no_wa"`
	// Reason explains why the row was skipped.
	Reason string `json:"reason"`
}

// batchRow is one row of the upload with the row number it had in the file.
type batchRow struct {
	number int
	record model.Rekrutmen
}

// requireUpload wraps a handler so it only runs for a member.
//
// The previous application let only the role ANGGOTA upload data
// (old-app/code.gs:735), and the button was hidden for the other roles
// (old-app/index.html:530-536).
func (s *Server) requireUpload(next func(http.ResponseWriter, *http.Request, model.User)) http.HandlerFunc {
	return s.requireAuth(func(w http.ResponseWriter, r *http.Request, user model.User) {
		if user.LevelOrDefault() != model.LevelAnggota {
			writeError(w, http.StatusForbidden, "hanya anggota yang dapat melakukan upload data")
			return
		}
		next(w, r, user)
	})
}

// handleBatchCreateRekrutmen stores many records in one call.
//
// It follows the previous application: the rows are validated first and
// nothing is saved when one of them is wrong, the owner of the records comes
// from the session, and a record whose WhatsApp number is already known is
// skipped and reported.
func (s *Server) handleBatchCreateRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	if !s.uploadLimiter.allow(user.Username) {
		writeError(w, http.StatusTooManyRequests, "terlalu banyak percobaan upload, coba lagi nanti")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxBatchBody)
	var raw []map[string]any
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if len(raw) == 0 {
		writeError(w, http.StatusBadRequest, "tidak ada data yang dapat di-upload")
		return
	}
	if len(raw) > MaxBatchRows {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("maksimal %d baris per upload", MaxBatchRows))
		return
	}

	now := time.Now().UTC()
	rows := make([]batchRow, 0, len(raw))
	var errs []batchRowError
	seen := map[string]int{}

	for i, cells := range raw {
		number := i + 2
		fields := batchFields(cells)
		rec := model.Rekrutmen{
			Nama:              fields["nama"],
			NoWA:              fields["no_wa"],
			Alamat:            fields["alamat"],
			RT:                fields["rt"],
			RW:                fields["rw"],
			Kecamatan:         fields["kecamatan"],
			Kelurahan:         fields["kelurahan"],
			Aktivitas:         fields["aktivitas"],
			TingkatPenerimaan: fields["tingkat_penerimaan"],
			Metode:            fields["metode"],
			Status:            fields["status"],
			Keterangan:        fields["keterangan"],
		}

		// The template has two UPL columns and only one of them may be
		// filled (old-app/code.gs:1003).
		upls := fields["upls"]
		uplm := fields["uplm"]
		upl := fields["upl"]
		filled := 0
		for _, v := range []string{upls, uplm, upl} {
			if v != "" {
				filled++
			}
		}
		switch {
		case filled > 1:
			errs = append(errs, batchRowError{Row: number, Field: "upl", Reason: "UPLS dan UPLM hanya boleh memilih salah satu"})
			continue
		case upls != "":
			rec.UPL = upls
		case uplm != "":
			rec.UPL = uplm
		case upl != "":
			rec.UPL = upl
		}

		if f := validate.Rekrutmen(&rec, s.aktivitasNames()); len(f) > 0 {
			for field, reason := range f {
				errs = append(errs, batchRowError{Row: number, Field: field, Reason: reason})
			}
			continue
		}

		// The same number twice in one file is a mistake in the file, so it
		// stops the upload.
		if first, dup := seen[rec.NoWA]; dup {
			errs = append(errs, batchRowError{
				Row:    number,
				Field:  "no_wa",
				Reason: fmt.Sprintf("nomor whatsapp sama dengan baris %d", first),
			})
			continue
		}
		seen[rec.NoWA] = number

		// The owner and the timestamps come from the session, so a file
		// cannot fill them in (old-app/code.gs:795-809).
		rec.ID = newID()
		rec.CreatedAt = now
		rec.UpdatedAt = now
		rec.CreatedBy = user.Username
		rec.UpdatedBy = user.Username

		rows = append(rows, batchRow{number: number, record: rec})
	}

	if len(errs) > 0 {
		sort.Slice(errs, func(i, j int) bool {
			if errs[i].Row != errs[j].Row {
				return errs[i].Row < errs[j].Row
			}
			return errs[i].Field < errs[j].Field
		})
		shown := errs
		if len(shown) > maxBatchErrors {
			shown = shown[:maxBatchErrors]
		}
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error":       "validasi upload gagal. tidak ada data yang disimpan.",
			"errors":      shown,
			"error_count": len(errs),
		})
		return
	}

	// A number that is already known is skipped and reported, so one old
	// record does not stop the whole upload.
	known, err := s.knownWAs()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memeriksa nomor whatsapp: "+err.Error())
		return
	}
	toCreate := make([]model.Rekrutmen, 0, len(rows))
	var skipped []batchSkipped
	for _, row := range rows {
		if known[row.record.NoWA] {
			skipped = append(skipped, batchSkipped{
				Row:    row.number,
				Nama:   row.record.Nama,
				NoWA:   row.record.NoWA,
				Reason: "nomor whatsapp sudah terdaftar",
			})
			continue
		}
		toCreate = append(toCreate, row.record)
	}

	if len(toCreate) > 0 {
		if err := s.Store().CreateRekrutmenBatch(toCreate); err != nil {
			writeError(w, http.StatusInternalServerError, "gagal menyimpan data: "+err.Error())
			return
		}
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"created": len(toCreate),
		"skipped": skipped,
	})
}

// batchFields reads the values of one row of the upload.
//
// The header of a column decides which field it fills. A header that is not
// known is ignored, so a file with extra columns still works.
func batchFields(cells map[string]any) map[string]string {
	out := make(map[string]string, len(batchColumn))
	for header, value := range cells {
		field, known := batchColumn[normalizeBatchHeader(header)]
		if !known {
			continue
		}
		out[field] = batchValue(value)
	}
	return out
}

// normalizeBatchHeader folds a header into the shape the mapping uses.
func normalizeBatchHeader(header string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(header)), " "))
}

// batchValue turns a value of the file into text.
//
// A number that came from a spreadsheet keeps its digits; a fractional number
// keeps its decimals.
func batchValue(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%v", v)
	case json.Number:
		return v.String()
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

// knownWAs returns the WhatsApp numbers that are already stored.
//
// The listing is not scoped, because a number of another account is a number
// that may not be used again.
func (s *Server) knownWAs() (map[string]bool, error) {
	list, _, err := s.Store().ListRekrutmen(model.Filter{Limit: 100000})
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(list))
	for _, rec := range list {
		if rec.NoWA != "" {
			out[rec.NoWA] = true
		}
	}
	return out, nil
}

// uploadLimiter counts the upload attempts of an account.
//
// The previous application had no limit, but an upload writes many rows at
// once, so the number of attempts per account is capped.
type uploadLimiter struct {
	mu          sync.Mutex
	hits        map[string][]time.Time
	maxAttempts int
}

// newUploadLimiter builds a limiter that allows a few uploads per account.
//
// The limit can be raised through the configuration, so the end-to-end
// checks can upload more files in a row.
func newUploadLimiter(maxAttempts int) *uploadLimiter {
	if maxAttempts <= 0 {
		maxAttempts = DefaultUploadMaxAttempts
	}
	return &uploadLimiter{hits: map[string][]time.Time{}, maxAttempts: maxAttempts}
}

// allow reports whether another upload may start.
//
// The limit is five uploads in ten minutes, unless the configuration raises
// the number of attempts.
func (l *uploadLimiter) allow(key string) bool {
	const window = 10 * time.Minute
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	kept := l.hits[key][:0]
	for _, t := range l.hits[key] {
		if now.Sub(t) < window {
			kept = append(kept, t)
		}
	}
	if len(kept) >= l.maxAttempts {
		l.hits[key] = kept
		return false
	}
	l.hits[key] = append(kept, now)
	return true
}

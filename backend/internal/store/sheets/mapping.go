package sheets

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"rekrutmen-tangsel/backend/internal/model"
)

// Field names used as mapping keys. They match the json tags of
// model.Rekrutmen, plus the source-only timestamp field.
const (
	FieldID                = "id"
	FieldNama              = "nama"
	FieldNoWA              = "no_wa"
	FieldAlamat            = "alamat"
	FieldRT                = "rt"
	FieldRW                = "rw"
	FieldKecamatan         = "kecamatan"
	FieldKelurahan         = "kelurahan"
	FieldAktivitas         = "aktivitas"
	FieldTingkatPenerimaan = "tingkat_penerimaan"
	FieldMetode            = "metode"
	FieldStatus            = "status"
	FieldKeterangan        = "keterangan"
	FieldUPL               = "upl"

	// FieldTimestamp is not a model field: it is the column a Google Form
	// response sheet uses for the submission time.
	FieldTimestamp = "timestamp"

	// Fields the model keeps beside the recruitment itself.
	FieldCreatedAt = "created_at"
	FieldUpdatedAt = "updated_at"
	FieldCreatedBy = "created_by"
	FieldUpdatedBy = "updated_by"

	// FieldUploader is the member who recorded the recruitment. It is a
	// column of the tab of the previous application.
	FieldUploader = "uploader"
)

// ColumnSpec says which sheet headers may carry one field.
type ColumnSpec struct {
	// Field is the mapping key.
	Field string

	// Aliases are the accepted header spellings. The first match wins.
	Aliases []string

	// Required marks a field the sheet must provide. A missing required
	// field fails the resolution instead of silently producing empty values.
	Required bool
}

// ColumnMap maps model fields to the columns of a sheet.
type ColumnMap struct {
	// Index maps a field name to the 0-based column index of the sheet.
	Index map[string]int

	// Headers is the header row the map was resolved from.
	Headers []string
}

// FormColumnSpec is the expected column layout of a Google Form response
// sheet.
//
// The aliases cover the spellings a form administrator is likely to use; the
// resolver also compares a punctuation-free form of every header, so
// "No. WhatsApp", "no whatsapp" and "no_wa" all match the same field.
func FormColumnSpec() []ColumnSpec {
	return []ColumnSpec{
		{Field: FieldTimestamp, Aliases: []string{"timestamp", "waktu", "waktu pengisian", "tanggal"}},
		{Field: FieldNama, Aliases: []string{"nama", "nama lengkap", "nama pendaftar"}, Required: true},
		{Field: FieldNoWA, Aliases: []string{"no wa", "no whatsapp", "nomor whatsapp", "whatsapp", "no hp", "nomor hp"}, Required: true},
		{Field: FieldAlamat, Aliases: []string{"alamat", "alamat lengkap"}},
		{Field: FieldRT, Aliases: []string{"rt"}},
		{Field: FieldRW, Aliases: []string{"rw"}},
		{Field: FieldKecamatan, Aliases: []string{"kecamatan"}, Required: true},
		{Field: FieldKelurahan, Aliases: []string{"kelurahan"}},
		{Field: FieldAktivitas, Aliases: []string{"aktivitas", "pekerjaan", "aktivitas pekerjaan"}},
		{Field: FieldTingkatPenerimaan, Aliases: []string{"tingkat penerimaan", "tingkat"}},
		{Field: FieldMetode, Aliases: []string{"metode"}},
		{Field: FieldStatus, Aliases: []string{"status"}},
		{Field: FieldKeterangan, Aliases: []string{"keterangan", "catatan"}},
		{Field: FieldUPL, Aliases: []string{"upl"}},
	}
}

// AppColumnSpec is the layout of the tab the application owns.
//
// Besides the headers of this application, it accepts the headers of the tab
// of the previous application, so one provider reads and writes both layouts.
// In the older layout "NAMA TARGET" carries the name of the recruited person
// and "NAMA ANGGOTA" the name of the member who recorded it.
func AppColumnSpec() []ColumnSpec {
	return []ColumnSpec{
		{Field: FieldID, Aliases: []string{"id"}},
		{Field: FieldNama, Aliases: []string{"nama", "nama target", "nama anggota"}},
		{Field: FieldNoWA, Aliases: []string{"no_wa", "no wa", "no whatsapp", "nomor whatsapp", "whatsapp"}},
		{Field: FieldAlamat, Aliases: []string{"alamat"}},
		{Field: FieldRT, Aliases: []string{"rt"}},
		{Field: FieldRW, Aliases: []string{"rw"}},
		{Field: FieldKecamatan, Aliases: []string{"kecamatan", "nama kecamatan"}},
		{Field: FieldKelurahan, Aliases: []string{"kelurahan"}},
		{Field: FieldAktivitas, Aliases: []string{"aktivitas", "aktifitas"}},
		{Field: FieldTingkatPenerimaan, Aliases: []string{"tingkat_penerimaan", "tingkat penerimaan"}},
		{Field: FieldMetode, Aliases: []string{"metode"}},
		{Field: FieldStatus, Aliases: []string{"status"}},
		{Field: FieldKeterangan, Aliases: []string{"keterangan", "catatan"}},
		{Field: FieldUPL, Aliases: []string{"upl", "upls", "uplm"}},
		{Field: FieldCreatedAt, Aliases: []string{"created_at"}},
		{Field: FieldUpdatedAt, Aliases: []string{"updated_at"}},
		{Field: FieldCreatedBy, Aliases: []string{"created_by"}},
		{Field: FieldUpdatedBy, Aliases: []string{"updated_by"}},
		{Field: FieldUploader, Aliases: []string{"nama anggota"}},
	}
}

// ResolveColumnMap matches the header row of a sheet against a spec.
//
// A header matches a field when it equals one of the field's aliases, or when
// both reduce to the same string once spaces, punctuation and case are
// removed. The first matching column wins, so a sheet with two equally named
// columns is read from the left.
func ResolveColumnMap(headers []string, spec []ColumnSpec) (ColumnMap, error) {
	normalized := make([]string, len(headers))
	for i, h := range headers {
		normalized[i] = normalizeHeader(h)
	}

	m := ColumnMap{Index: map[string]int{}, Headers: headers}
	used := make([]bool, len(headers))
	missing := []string{}

	for _, s := range spec {
		idx := -1
		for _, alias := range s.Aliases {
			for i, h := range normalized {
				if used[i] {
					continue
				}
				if h == normalizeHeader(alias) {
					idx = i
					break
				}
			}
			if idx >= 0 {
				break
			}
		}
		if idx < 0 {
			if s.Required {
				missing = append(missing, s.Field)
			}
			continue
		}
		used[idx] = true
		m.Index[s.Field] = idx
	}

	if len(missing) > 0 {
		return m, fmt.Errorf("kolom wajib tidak ditemukan: %s; header yang tersedia: %s",
			strings.Join(missing, ", "), strings.Join(headers, ", "))
	}
	return m, nil
}

// normalizeHeader reduces a header to a comparable form: lower case, without
// spaces, punctuation or separators.
func normalizeHeader(h string) string {
	h = strings.ToLower(strings.TrimSpace(h))
	var b strings.Builder
	for _, r := range h {
		switch r {
		case ' ', '\t', '.', '_', '-', '/', '(', ')', ':', ',':
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Has reports whether the map knows the given field.
func (m ColumnMap) Has(field string) bool {
	_, ok := m.Index[field]
	return ok
}

// Get returns the value of a field in a row, or an empty string when the
// sheet does not have that column.
func (m ColumnMap) Get(row []string, field string) string {
	idx, ok := m.Index[field]
	if !ok || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

// Record converts one source row into a recruitment record.
//
// The conversion is intentionally forgiving: values that do not fit the
// model are passed through as they are, because the import reports the
// problems instead of dropping the row.
func (m ColumnMap) Record(row []string) model.Rekrutmen {
	return model.Rekrutmen{
		Nama:              m.Get(row, FieldNama),
		NoWA:              NormalizeNoWA(m.Get(row, FieldNoWA)),
		Alamat:            m.Get(row, FieldAlamat),
		RT:                m.Get(row, FieldRT),
		RW:                m.Get(row, FieldRW),
		Kecamatan:         m.Get(row, FieldKecamatan),
		Kelurahan:         m.Get(row, FieldKelurahan),
		Aktivitas:         m.Get(row, FieldAktivitas),
		TingkatPenerimaan: m.Get(row, FieldTingkatPenerimaan),
		Metode:            model.NormalizeMetode(m.Get(row, FieldMetode)),
		Status:            m.Get(row, FieldStatus),
		Keterangan:        m.Get(row, FieldKeterangan),
		UPL:               m.Get(row, FieldUPL),
	}
}

// Timestamp returns the submission time of a row, if the sheet has one.
func (m ColumnMap) Timestamp(row []string) (time.Time, bool) {
	raw := m.Get(row, FieldTimestamp)
	if raw == "" {
		return time.Time{}, false
	}
	t, ok := ParseTimestamp(raw)
	if !ok {
		return time.Time{}, false
	}
	return t, true
}

// NormalizeNoWA keeps the digits of a phone number and drops everything else.
//
// Form responses often carry a number with a country code, spaces or
// punctuation; the application compares numbers as plain digit strings.
func NormalizeNoWA(v string) string {
	var b strings.Builder
	for _, r := range v {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// timestampLayouts are the formats a Google Form response sheet may use.
var timestampLayouts = []string{
	"02/01/2006 15:04:05",
	"02/01/2006 15:04",
	"02-01-2006 15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	time.RFC3339,
	"02/01/2006",
	"2006-01-02",
}

// ParseTimestamp reads a timestamp written in one of the common sheet
// formats. It reports whether the value could be read.
func ParseTimestamp(v string) (time.Time, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return time.Time{}, false
	}
	for _, layout := range timestampLayouts {
		if t, err := time.Parse(layout, v); err == nil {
			return t.UTC(), true
		}
	}
	// Some sheets export the time in a locale format the layouts above miss.
	for _, layout := range []string{"02/01/2006 03:04:05 PM", "01/02/2006 15:04:05"} {
		if t, err := time.Parse(layout, v); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}

// DerivedID builds a stable id for a record that has none of its own.
//
// The id is derived from the phone number and the submission time, so running
// the import twice produces the same id and the import stays idempotent.
func DerivedID(noWA string, submittedAt time.Time) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(noWA) + "|" + submittedAt.UTC().Format(time.RFC3339)))
	// Format the first 16 bytes as a version 4 uuid.
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

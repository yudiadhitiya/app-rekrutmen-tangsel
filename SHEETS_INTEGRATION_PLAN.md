# Rencana Integrasi Google Sheets — Spreadsheet "Lapor Rekrutmen DPD PKS Tangsel (Responses)"

Dokumen ini adalah rencana implementasi (bukan perubahan kode). Kode yang sudah
ada **tidak diubah** dalam dokumen ini.

Status item memakai penanda yang sama dengan dokumen plan lain:
⬜ belum · 🔧 sedang dikerjakan · ✅ selesai & terverifikasi.

---

## 1. Tujuan & batasan

**Tujuan:** aplikasi memakai spreadsheet Google Sheets eksisting
`Lapor Rekrutmen DPD PKS Tangsel (Responses)` sebagai sumber penyimpanan aktif,
tanpa kehilangan data formulir yang sudah terkumpul.

**Yang sudah ada (tidak diulang):**

| Komponen | Status | Lokasi |
| --- | --- | --- |
| Provider `google_sheets` (Store) | ✅ | `backend/internal/store/sheets/store.go` |
| OAuth2 service-account JWT (stdlib) | ✅ | `sheets/store.go:176-256` |
| CRUD rekrutmen + master aktivitas | ✅ | `sheets/store.go:474-701` |
| Auto-create tab & header | ✅ | `sheets/store.go:103-130` |
| Test provider (fake Sheets API) | ✅ 6 test | `sheets/store_test.go` |
| Ganti provider dari UI | ✅ | `api.go:590`, `SettingsView.vue` |

**Batasan yang harus diputuskan lebih dulu (Bagian 2).**

---

## 2. Keputusan arsitektur (blocking — perlu konfirmasi)

Sasaran spreadsheet adalah hasil Google Form, sehingga tab respons-nya
(`Form Responses 1`) dimiliki oleh Form: jumlah, urutan, dan nama kolomnya
dikunci, dan barisnya ditulis oleh Form.

Terdapat dua pendekatan:

### Opsi A — Tab terpisah di spreadsheet yang sama (REKOMENDASI)

```
Spreadsheet "Lapor Rekrutmen DPD PKS Tangsel (Responses)"
├── Form Responses 1   ← read-only, milik Form. Dipakai sebagai sumber impor.
├── Rekrutmen          ← milik aplikasi. CRUD aplikasi terjadi di sini.
└── Aktivitas          ← milik aplikasi. Master aktivitas.
```

- Data formulir diimpor **sekali** ke tab `Rekrutmen` (Fase 4).
- Aplikasi tetap menulis ke tab miliknya sendiri → tidak ada baris sintetis
  bercampur dengan jawaban responden.
- Form boleh tetap aktif dipakai; perubahan di app tidak merusak data Form.
- Kolom yang tidak ada di Form (`id`, `created_at`, `created_by`, `updated_at`,
  `updated_by`) di-*derive* saat impor dan disimpan sebagai kolom tambahan di
  tab `Rekrutmen`.

**Kelemahan:** sampai impor dijalankan, app tidak melihat data Form.

### Opsi B — Aplikasi menulis langsung ke tab `Form Responses 1`

- `Rekrutmen` dan `Aktivitas` tidak dipakai; app memakai tab Form apa adanya.
- **Tidak direkomendasikan.** alasan:
  1. Kolom `id`/audit tidak bisa ditambahkan — Form mengunci skema tab.
  2. Baris yang ditulis app indistinguishable dari jawaban asli; Form berikutnya
     akan menyisipkan baris baru dan menggeser urutan.
  3. `DeleteRekrutmen` pada provider Sheets berarti *clear* baris — untuk tab
     Form ini berarti menghapus jawaban responden, termasuk jawaban yang sudah
     dipakai sebagai rujukan anderer.
  4. Halaman rekap/analisis Form bisa rusak.

> **Rekomendasi: Opsi A.** Seluruh rencana berikut memakai Opsi A. Fase 3
> (mapping) tetap berguna untuk Opsi B dengan perubahan kecil pada target tab.

---

## 3. Status per fase

| Fase | Status | Item selesai | Catatan |
| --- | --- | --- | --- |
| 0 — Discovery & keputusan | 🔧 | — | spreadsheet ID, daftar tab, dan header kolom masih perlu diambil dari spreadsheet sebenarnya (D-105…D-110) |
| 1 — Kredensial & akses | 🔧 | D-111, D-112, D-114 | service account sudah dibuat dan spreadsheet sudah dibagikan; verifikasi akses (D-113) belum dilakukan |
| 2 — Konfigurasi provider | ✅ | D-116…D-124 | tab & mode configurable, metadata spreadsheet menggantikan deteksi string, `batchUpdate` diperbaiki, id spreadsheet divalidasi |
| 3 — Lapisan mapping kolom | ✅ | D-125…D-131 | `internal/store/sheets/mapping.go` + `mapping_test.go` |
| 4 — Impor data Responses | ✅ | D-132…D-139 | `cmd/impor-responses` + `main_test.go` |
| 5 — UI & status | ⬜ | — | panel status Sheets di Pengaturan |
| 6 — Keamanan & keandalan | ⬜ | — | optimistic concurrency (D-145) belum ada |
| 7 — Pengujian | 🔧 | D-151…D-156 | test unit selesai; test integrasi spreadsheet sungguhan (D-158) dan smoke test (D-159) menyusul |
| 8 — Dokumentasi & rollout | 🔧 | D-160, D-161, D-162, D-163 | README, DEPLOYMENT_PLAN, TESTING_PLAN, dan dokumen ini sudah diperbarui |

Status di atas merujuk nomor item di bagian 3 dan 4. Item yang belum disebut
di kolom "Item selesai" masih ⬜.

## 3. Fase kerja

Estimasi total: **8–10 hari kerja** (1 developer),excluding waktu tunggu
persetujuan admin Google Workspace.

### Fase 0 — Discovery & keputusan (0.5–1 hari) 🔧

| ID | Item |
| --- | --- |
| D-105 | Ambil ID spreadsheet dari URL "Lapor Rekrutmen DPD PKS Tangsel (Responses)". ID adalah segmen `/d/` pada URL. |
| D-106 | Inventarisasi tab: panggil `GET https://sheets.googleapis.com/v4/spreadsheets/{id}` dan catat `sheets[].properties.title`. Target Hopefully: `Form Responses 1`. |
| D-107 | Dump header baris 1 tab respons untuk menyusun pemetaan kolom (Bagian 4). |
| D-108 | Klarifikasi ke pemilik organisasi: apakah Form masih dipakai aktif? Data Form perlu tampil di app? Siapa yang boleh jadi editor spreadsheet? |
| D-109 | Backup: "File → Make a copy" spreadsheet, dan unduh ekspor CSV tab respons. Simpan di luar repo. |

**Keluar Fase 0:** spreadsheet ID, daftar tab, header kolom, konfirmasi Opsi A,
backup tersedia.

### Fase 1 — Kredensial & akses (0.5 hari) 🔧

| ID | Item |
| --- | --- |
| D-110 | Buat project Google Cloud, aktifkan **Google Sheets API**. |
| D-111 | Buat service account, unduh key JSON, simpan di path di luar repo dengan izin `0600`. Jangan pernah di-commit. |
| D-112 | Share spreadsheet ke email service account dengan role **Editor**. Bila spreadsheet milik domain yang adminnya tidak tersedia, buat **salinan** spreadsheet yang dikelola service account, lalu opsional: salinan itu yang jadi target, dan lakukan impor dua arah manual. |
| D-113 | Uji akses manual sebelum kode: ambil access token via JWT assertion, lalu `GET /v4/spreadsheets/{id}` dan `GET .../values/Form Responses 1!A1:Z1`. Catat hasilnya sebagai baseline. |
| D-114 | Tambah `backend/.env.example` (atau `docs/env.example`) berisi semua variabel yang dibutuhkan Fase 2. |

**Variabel lingkungan (target akhir):**

| Variabel | Keterangan | Default |
| --- | --- | --- |
| `GOOGLE_SHEET_ID` | ID spreadsheet | — (wajib) |
| `GOOGLE_KEY_FILE` | path key JSON service account | — (wajib) |
| `GOOGLE_SHEETS_MODE` | `app` (tab sendiri) atau `form` (tab Form) | `app` |
| `GOOGLE_SHEETS_TAB_REKRUTMEN` | nama tab data | `Rekrutmen` |
| `GOOGLE_SHEETS_TAB_AKTIVITAS` | nama tab master aktivitas | `Aktivitas` |
| `GOOGLE_SHEETS_TAB_SUMBER` | nama tab sumber untuk impor | `Form Responses 1` |

### Fase 2 — Konfigurasi provider di backend (1–1.5 hari) ✅

| ID | Item | File |
| --- | --- | --- |
| D-116 | Ubah `SheetRekrutmen`/`SheetAktivitas` dari `const` menjadi field pada `Store`, dengan default nama lama (backward compatible). | `backend/internal/store/sheets/store.go` |
| D-117 | Tambahkan mode `app`/`form`. Pada mode `app`: `ensureTabs` tetap berjalan. Pada mode `form`: `ensureTabs` **dilarang** membuat tab, dan hanya memvalidasi bahwa tab sumber ada (jika tidak → error informatif, bukan membuat tab). | `sheets/store.go` |
| D-118 | Ganti deteksi tab hilang berbasis pencocokan string `"Unable to parse range"` dengan daftar tab dari metadata spreadsheet (di-cache, TTL mis. 5 menit). Buang `isMissingTab`. | `sheets/store.go:102-135` |
| D-119 | Buat range baca dinamis: `A1:AZ` atau mengikuti `gridProperties` dari metadata, bukan `A1:Z` hardcoded. | `sheets/store.go:279` |
| D-120 | `Ping()` untuk mode `app` tetap boleh menulis (auto-create tab). Tambahkan `PingReadOnly(ctx)` yang hanya membaca, dipakai health check. | `sheets/store.go:93-100` |
| D-121 | Tambahkan `ListTabs()` dan `CountRows()` pada Store Sheets untuk keperluan status/verifikasi. | `sheets/store.go` |
| D-122 | `Config` API gaining field baru; `storeFactory` meneruskannya. | `backend/internal/api/api.go:29-37, 59-77` |
| D-123 | Baca env baru di entrypoint. | `backend/cmd/server/main.go:57-63` |
| D-124 | `handleGetSettings` menambah info Sheets: mode, tab terpakai, jumlah baris, status koneksi, dan daftar tab yang ditemukan. | `api.go:579-588` |

### Fase 3 — Lapisan mapping kolom (1.5 hari) ✅

Tabel pemetaan kolom (hasil D-107) disimpan sebagai konfigurasi, bukan
hardcode. Buatkan `backend/internal/store/sheets/mapping.go`:

| ID | Item |
| --- | --- |
| D-125 | Tipe `ColumnMap`: pemetaan dua arah `field model ↔ nama header di sheet`. |
| D-126 | `ResolveColumnMap(headers []string, spec ColumnSpec) (ColumnMap, error)`: pencocokan exact → case-insensitive → alias (mis. `No. WA`, `No WA`, `Nomor Whatsapp`, `WhatsApp`, `HP`). |
| D-127 | Kolom wajib yang hilang → error yang menyebutkan daftar kolom yang tidak ditemukan, plus daftar header yang tersedia. |
| D-128 | `rekrutmenToValues` / `valuesToRekrutmen` menjadi mapping-aware (berbasis indeks kolom hasil resolve, bukan urutan tetap). |
| D-129 | Kolom turunan untuk baris hasil impor: `id` (UUID deterministik dari hash `no_wa` + timestamp Form → idempoten), `created_at` (dari kolom Timestamp Form, format `dd/MM/yyyy HH:mm:ss` dan variasinya), `updated_at`, `created_by`=`form`, `updated_by`=`form`. |
| D-130 | Normalisasi tambahan saat baca: `model.NormalizeMetode` (sudah ada, tangani "Program"/"Antara" lama), trim spasi, `no_wa` dinormalisasi (buang spasi/tanda baca non-digit) agar pengecekan duplikat konsisten. |
| D-131 | aktivitas yang tidak ada di master → tandai baris untuk review, jangan gagalkan seluruh impor. |

### Fase 4 — Impor data Responses (1.5–2 hari) ✅

Buat binary terpisah `backend/cmd/impor-responses/main.go` (stdlib only, tidak
ikut ter-deploy sebagai server).

| ID | Item |
| --- | --- |
| D-132 | Flag `--dry-run` (default), `--apply`, `--spreadsheet-id`, `--tab-sumber`, `--tab-tujuan`, `--report <path>`. |
| D-133 | Baca semua baris tab sumber, resolve mapping, konversi ke `model.Rekrutmen`, validasi dengan aturan `internal/validate` yang sudah ada. |
| D-134 | Validasi yang gagal **tidak** menggagalkan impor: baris ditandai `status` = `Pending` + `keterangan` = "hasil impor form, perlu verifikasi", dan alasan dicatat di laporan. |
| D-135 | Aktivitas yang belum ada di master → ditambahkan ke tab `Aktivitas` (aktif) supaya nilainya valid saat diedit nanti. |
| D-136 | Pemeriksaan uniqueness `no_wa`: baris duplikat terhadap data yang sudah ada di tab tujuan **dilewati** (idempoten, tidak menimpa). |
| D-137 | Laporan impor `data/impor-responses-report.json`: jumlah terbaca, berhasil diimpor, dilewati, duplikat, gagal validasi + alasan per baris. |
| D-138 | `--apply` bersifat add-only (tidak pernah update/menimpa baris lama). |
| D-139 | Uji idempoten: jalankan dua kali → jumlah baris tidak bertambah. |

**Contoh perintah:**

```bash
cd backend
go run ./cmd/impor-responses \
  --spreadsheet-id "$GOOGLE_SHEET_ID" \
  --tab-sumber "Form Responses 1" \
  --tab-tujuan "Rekrutmen" \
  --dry-run
# periksa laporan, baru:
go run ./cmd/impor-responses ... --apply
```

### Fase 5 — UI & status (1 hari) ⬜

| ID | Item | File |
| --- | --- | --- |
| D-140 | Panel status Google Sheets di Pengaturan: mode, tab, jumlah baris, waktu terakhir cek, hasil ping. | `frontend/src/views/SettingsView.vue` |
| D-141 | Peringatan eksplisit: mengganti provider tidak memindahkan data lama; jalankan impor dulu. | `SettingsView.vue` |
| D-142 | Banner ringan di Daftar & Form saat provider aktif `google_sheets`: "data disimpan di Google Sheets". | `frontend/src/views/ListView.vue`, `FormView.vue` |
| D-143 | Endpoint status read-only: `GET /api/settings/sheets-status` (admin). Tidak mencederai pengguna non-admin. | `api.go` |
| D-144 | Konfirmasi hapus pada provider Sheets memperjelas: baris **dikosongkan**, bukan dihapus, dan data bisa dipulihkan dari backup. | `ConfirmDialog` |

### Fase 6 — Keamanan & keandalan (1 hari) ⬜

| ID | Item |
| --- | --- |
| D-145 | Optimistic concurrency: sebelum update/delete, baca baris dan bandingkan `updated_at`; berbeda → 409 "baris sudah diubah di sheet". |
| D-146 | Cache daftar tab + jumlah baris (TTL 60s) agar tidak memukul API Sheets tiap request. Token OAuth sudah di-cache (sudah ada). |
| D-147 | Rate limit endpoint Sheets status (admin, mis. 1 req/10 detik). |
| D-148 | Log audit provider switch + hasil impor ke `data/audit.log` (append-only). |
| D-149 | Rotasi key service account: catat tanggalKedaluwarsa/rotasi di DEPLOYMENT_PLAN. |
| D-150 | Ensure `GOOGLE_KEY_FILE` tidak pernah ikut ter-deploy di image container; pakai secret mount. |

### Fase 7 — Pengujian (1–1.5 hari) 🔧

| ID | Test | ID Test |
| --- | --- | --- |
| D-151 | Fake Sheets API diperluas dengan multi-tab + endpoint metadata spreadsheet. | `sheets/store_test.go` |
| D-152 | Test resolusi mapping: exact, alias, case-insensitive, kolom wajib hilang. | T-81 |
| D-153 | Test konversi nilai: timestamp Form, normalisasi metode lama, normalisasi no_wa. | T-82 |
| D-154 | Test UUID turunan deterministik (idempoten). | T-83 |
| D-155 | Test impor pada fixture CSV respons form: dry-run tidak menulis, apply menulis, dua kali apply idempoten. | T-84 |
| D-156 | Test mode `form` tidak pernah membuat tab. | T-85 |
| D-157 | Test optimistic concurrency → 409. | T-86 |
| D-158 | Smoke test manual dengan spreadsheet sungguhan (skip bila `GOOGLE_SHEET_ID` kosong) — item ini sudah ada di TESTING_PLAN sebagai **T-70**. | T-70 |
| D-159 | Checklist smoke test_end-to-end: login → daftar → tambah → ubah → hapus → ganti provider → impor → verifikasi jumlah baris. | T-87 |

### Fase 8 — Dokumentasi & rollout (0.5 hari) 🔧

| ID | Item |
| --- | --- |
| D-160 | README: bagian "Menghubungkan ke Spreadsheet Responses", contoh env, prosedur impor. |
| D-161 | DEPLOYMENT_PLAN: variabel env baru, prosedur share service account, rollback. |
| D-162 | DEVELOPMENT_PLAN: tambahkan item D-105…D-162 dengan status. |
| D-163 | TESTING_PLAN: tambahkan T-81…T-87. |
| D-164 | Prosedur rollback terdokumentasi: kembali ke provider `local`/`supabase`, impor ulang dari tab `Rekrutmen` (yang tetap utuh), atau restore dari backup. |

---

## 4. Pemetaan kolom (draft — wajib dikonfirmasi di D-107)

| Kolom di tab respons Form | Field model | Aturan konversi |
| --- | --- | --- |
| Timestamp | `created_at` | parse `dd/MM/yyyy HH:mm:ss` → RFC3339 UTC |
| Nama | `nama` | trim |
| No. WhatsApp | `no_wa` | non-digit dibuang, harus ≥ 9 digit |
| Alamat | `alamat` | trim |
| RT | `rt` | numeric string, leading zero dipertahankan |
| RW | `rw` | numeric string |
| Kecamatan | `kecamatan` | harus ada di `seed.wilayah.json` |
| Kelurahan | `kelurahan` | harus sesuai kecamatan |
| Aktivitas / Pekerjaan | `aktivitas` | harus ada di master aktivitas (atau ditambahkan saat impor) |
| Tingkat Penerimaan | `tingkat_penerimaan` | enum: Sangat Terbuka / Terbuka / Netral / Kurang Terbuka / Tertutup |
| Metode | `metode` | enum: Halaqoh / Program Antara / Dakwah Fardiyah; nilai lama "Program"/"Antara" → `Program Antara` |
| Status | `status` | enum: Berhasil / Proses / Pending / Gagal; nilai tak dikenal → `Pending` |
| Keterangan | `keterangan` | trim |
| UPL | `upl` | kosong / UPLS / UPLM |
| — | `id` | UUID deterministik dari `sha256(no_wa + "|" + timestamp)` |
| — | `created_by` / `updated_by` | `form` |
| — | `updated_at` | = `created_at` |

Kolom yang tidak ada di Form (`id`, audit) **tidak** dapat ditambahkan ke tab
`Form Responses 1` — Form mengunci skema. Karena itu Opsi A memakai tab
`Rekrutmen` yang skemanya dikendalikan aplikasi.

---

## 5. Alur rollout

1. Backup spreadsheet + ekspor CSV. *(Fase 0)*
2. Setup service account, share spreadsheet, verifikasi akses via `curl`. *(Fase 1)*
3. Merge kode Fase 2–3 → `go test ./...` hijau.
4. Terapkan kode Fase 4–5.
5. Jalankan impor `--dry-run`,-review laporan bersama pemilik data.
6. Jalankan impor `--apply`; verifikasi jumlah baris = jumlah baris Form (minus duplikat/invalid).
7. Admin → Pengaturan → provider **Google Sheets**.
8. Smoke test end-to-end (T-87) + pantau `data/audit.log`.
9. Update README/DEPLOYMENT_PLAN.

**Rollback:** kembalikan provider ke `local` atau `supabase`; data tetap utuh di
tab `Rekrutmen`; bila perlu, restore spreadsheet dari backup.

---

## 6. Risiko

| Risiko | Dampak | Mitigasi |
| --- | --- | --- |
| Spreadsheet milik domain, service account tak bisa jadi Editor | Tinggi | D-108 lebih awal; siapkan jalur salinan spreadsheet |
| Kolom Form tidak cocok dengan master wilayah/aktivitas | Sedang | Impor toleran (D-134, D-135) + laporan |
| Menulis ke spreadsheet destroys data Form | Tinggi | Opsi A (tab terpisah) |
| Concurrent write tanpa transaksi | Sedang | Optimistic concurrency (D-145) + operator tunggal per sheet |
| Kuota API Sheets | Rendah–Sedang | Cache (D-146), operasi batch, mengukur jumlah baris dulu |
| Key service account bocor | Tinggi | Jangan commit, izin 0600, secret mount (D-150) |
| Baris yang di-*clear* dianggap hilang | Sedang | Dokumentasi di UI (D-144) + backup berkala |

---

## 7. Definition of Done

- [ ] Impor berjalan dari spreadsheet nyata: jumlah baris terp-import cocok dengan sumber.
- [ ] Impor idempoten (dijalankan dua kali → jumlah baris sama).
- [ ] Provider `google_sheets` aktif di production dan CRUD app bekerja penuh.
- [ ] Data formulir lama dapat dilihat di UI.
- [ ] Semua test baru (T-81…T-87) hijau; `go test ./...` dan `npm run build` lolos.
- [ ] Tidak ada kredensial di repo; `git grep -i "private_key"` bersih.
- [ ] README, DEVELOPMENT_PLAN, TESTING_PLAN, DEPLOYMENT_PLAN diperbarui.

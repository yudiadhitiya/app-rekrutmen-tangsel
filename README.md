# Aplikasi Rekrutmen Tangsel

Aplikasi web untuk mencatat dan mengelola data rekrutmen di wilayah Kota
Tangerang Selatan. Dibangun sesuai PRD MVP
([`PRD_Rekrutmen_Tangsel_MVP.md`](PRD_Rekrutmen_Tangsel_MVP.md)).

## Dokumen Proyek

| Dokumen | Isi |
| --- | --- |
| [`PRD_Rekrutmen_Tangsel_MVP.md`](PRD_Rekrutmen_Tangsel_MVP.md) | Kebutuhan produk & kriteria penerimaan |
| [`API_CONTRACT.md`](API_CONTRACT.md) | Kontrak API (bentuk request/response, enum, aturan validasi) |
| [`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md) | Rencana pengembangan + status per item + traceability ke PRD |
| [`TESTING_PLAN.md`](TESTING_PLAN.md) | Rencana pengujian + daftar test + checklist smoke test |
| [`DEPLOYMENT_PLAN.md`](DEPLOYMENT_PLAN.md) | Rencana deployment + prosedur rilis/rollback + checklist pra-rilis |

Ketiga dokumen plan (`DEVELOPMENT_`, `TESTING_`, `DEPLOYMENT_PLAN.md`) berisi
status **per item** dengan penanda ✅ (selesai & terverifikasi), 🔧 (sedang
dikerjakan), dan ⬜ (belum). Gunakan ketiganya sebagai catatan checkpoint:
kalau proses terhenti, item ✅ terakhir sebelum titik berhenti adalah posisi
kemajuan yang sudah aman.

## Fitur

- Form tambah / ubah data rekrutmen dengan validasi di frontend dan backend.
- Daftar data dengan pencarian (nama / nomor WhatsApp) dan filter kecamatan,
  kelurahan, metode, serta status.
- Hapus data dengan konfirmasi.
- Master wilayah (7 kecamatan, 54 kelurahan) dan master aktivitas/pekerjaan
  yang dapat dikelola admin.
- Pilihan sumber penyimpanan: **Supabase**, **Google Sheets**, atau penyimpanan
  lokal, dapat diganti admin lewat menu Pengaturan.
- Autentikasi berbasis peran (admin / petugas).
- Tampilan tabel pada layar lebar dan kartu pada layar sempit.

## Stack

| Bagian    | Teknologi                                    |
| --------- | -------------------------------------------- |
| Frontend  | Vue 3 + Vite + vue-router (JavaScript)       |
| Backend   | Go 1.24, REST API, **hanya standard library** |
| Penyimpanan | Supabase (PostgREST) / Google Sheets API / file lokal |

Backend sengaja tidak memakai dependensi eksternal: Supabase diakses langsung
lewat API PostgREST, dan Google Sheets lewat REST API + OAuth2 service account.
`go.mod` tidak punya `require` sama sekali.

## Menjalankan

### Backend

```bash
cd backend
go run ./cmd/server
```

Default: `http://localhost:8080`.

Environment (opsional):

| Variabel            | Kegunaan                                                |
| ------------------- | ------------------------------------------------------- |
| `ADDR`              | alamat server, default `:8080`                          |
| `DATA_DIR`          | lokasi data runtime, default `./data`                   |
| `FRONTEND_DIR`      | folder hasil build frontend, default `frontend/dist`     |
| `SESSION_SECRET`    | kunci penandatangan cookie sesi                         |
| `SUPABASE_URL`      | URL project Supabase, contoh `https://xxx.supabase.co`  |
| `SUPABASE_API_KEY`  | service-role key Supabase                               |
| `GOOGLE_SHEET_ID`   | ID spreadsheet Google Sheets                            |
| `GOOGLE_KEY_FILE`   | path file JSON service account Google                   |

### Frontend

```bash
cd frontend
npm install
npm run dev      # http://localhost:5173, proxy /api ke localhost:8080
npm run build    # menghasilkan frontend/dist, dilayani backend
```

## Akun awal

Dibuat otomatis pada jalankan pertama, tersimpan di `data/users.json`:

| Username  | Password    | Peran   |
| --------- | ----------- | ------- |
| `admin`   | `admin123`  | admin   |
| `petugas` | `petugas123`| petugas |

**Ganti password ini sebelum dipakai di produksi.**

## Menyiapkan Supabase

1. Jalankan [`backend/supabase/schema.sql`](backend/supabase/schema.sql) di
   SQL Editor Supabase.
2. Isi `SUPABASE_URL` dan `SUPABASE_API_KEY` (service-role key).
3. Admin → Pengaturan → pilih provider **Supabase**.

Skema sudah mengaktifkan Row Level Security: aplikasi hanya mengakses lewat
backend, akses langsung dari browser ditolak.

## Menyiapkan Google Sheets

1. Buat spreadsheet kosong, catat ID-nya.
2. Buat service account di Google Cloud, unduh key JSON.
3. Bagikan spreadsheet ke email service account dengan akses editor.
4. Isi `GOOGLE_SHEET_ID` dan `GOOGLE_KEY_FILE`.
5. Admin → Pengaturan → pilih provider **Google Sheets**.

Tab `Rekrutmen` dan `Aktivitas` beserta header-nya dibuat otomatis saat
provider pertama kali dipakai.

## Test

```bash
cd backend && go test ./...
cd frontend && npm run build
```

Cakupan test: aturan validasi, keunikan nomor WhatsApp, filter/paginasi,
cascade rename master aktivitas, serta hak akses per peran, serta ketiga
provider penyimpanan (local, Supabase, Google Sheets).

## Struktur

```
backend/
  cmd/server/           entrypoint
  internal/api/         handler HTTP, routing, RBAC
  internal/auth/        login, hash PBKDF2, cookie sesi
  internal/model/       tipe domain
  internal/seed/        master wilayah (embed) + seed aktivitas
  internal/settings/    provider aktif
  internal/store/       antarmuka + 3 implementasi provider
  internal/validate/    aturan validasi input
frontend/
  src/api.js            wrapper fetch
  src/views/            Login, List, Form, Aktivitas, Settings
  src/components/       Toast, ConfirmDialog
```

Kontrak API terdokumentasi di [`API_CONTRACT.md`](API_CONTRACT.md).

## Catatan

- Mengganti provider **tidak** memindahkan data lama; provider yang diganti
  akan terisi kosong (kecuali aktivitas, yang di-seed ulang).
- Nomor WhatsApp disimpan sebagai string agar angka `0` di awal tidak hilang,
  dan wajib unik.
- Data master wilayah mengikuti Permendagri (7 kecamatan, 54 kelurahan).

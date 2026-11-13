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
| `GOOGLE_SHEETS_MODE` | `app` (tab milik aplikasi) atau `form` (tab respons formulir), default `app` |
| `GOOGLE_SHEETS_TAB_REKRUTMEN` | nama tab data rekrutmen, default `Rekrutmen` |
| `GOOGLE_SHEETS_TAB_AKTIVITAS` | nama tab master aktivitas, default `Aktivitas` |
| `GOOGLE_SHEETS_TAB_SUMBER`    | nama tab sumber untuk impor, default `Form Responses 1` |
| `GOOGLE_SHEETS_TAB_PENGGUNA`  | nama tab akun app lama, default `USERS` |
| `LEGACY_AUTH_MODE`            | `upgrade` (default), `sheet`, atau `off` — cara memakai akun app lama |

Backend membaca konfigurasi langsung dari environment; tidak ada file `.env`
yang dimuat otomatis. Untuk pengembangan, ekspor variabelnya di shell:

```bash
export GOOGLE_SHEET_ID=1AbCdEf...
export GOOGLE_KEY_FILE="$HOME/.config/rekrutmen-tangsel/google-key.json"
```

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
provider pertama kali dipakai. Tab `REKRUTMEN` milik aplikasi lama juga bisa
langsung dipakai: header lamanya dikenali, jadi data lama terbaca dan data
baru ditulis dengan kolom tab tersebut.

Provider aktif tersimpan di `DATA_DIR/settings.json` (`{"provider":
"google_sheets"}`), jadi pilihan itu tetap berlaku di setiap penjalanan
backend. Mengganti provider lewat menu **Pengaturan** juga menulis file itu.

### Lokasi file kredensial

File JSON service account berisi private key yang memberi akses baca/tulis ke
seluruh spreadsheet yang dibagikan ke akun tersebut. Backend menolak file yang
izinnya lebih longgar dari `0600`:

```bash
mkdir -p ~/.config/rekrutmen-tangsel
cp service-account.json ~/.config/rekrutmen-tangsel/google-key.json
chmod 600 ~/.config/rekrutmen-tangsel/google-key.json
```

Di server produksi, letakkan file di `/etc/rekrutmen-tangsel/google-key.json`
dengan izin `0600`, lalu isi `GOOGLE_KEY_FILE` di
`/etc/rekrutmen-tangsel.env`. Nama file kredensial umum sudah masuk
`.gitignore`, tetapi file tetap sebaiknya diletakkan di luar repo.

### Masuk dengan akun dari aplikasi lama

Aplikasi ini bisa menerima akun dari aplikasi lama (Google Apps Script) yang
datanya berada di tab `USERS` spreadsheet yang sama. Akun lama disimpan dengan
hash SHA-256 tanpa salt, sehingga tidak bisa dipakai ulang oleh aplikasi ini
sebagai hash.

Mekanismenya: saat pertama kali login, kredensial dicek ke tab `USERS`. Kalau cocok,
akun itu langsung disimpan di `data/users.json` dengan hash PBKDF2 bersalt —
jadi login berikutnya tidak perlu menyentuh spreadsheet lagi.

| `LEGACY_AUTH_MODE` | Perilaku |
| --- | --- |
| `upgrade` (default) | login pertama memakai akun lama, lalu akun disimpan lokal |
| `sheet` | setiap login dicek ke tab `USERS`; salinan lokal di `data/users.json` tidak dipakai sebagai penentu password |
| `off` | akun lama tidak dipakai sama sekali |

Peran dipetakan: `SYSTEM ADMIN`, `ADMIN PTD`, `ADMIN KECAMATAN` menjadi
`admin`; `ANGGOTA` menjadi `petugas`. Akun dengan `AKTIF` selain `Y` dan akun
tanpa `PASSWORD_HASH` tidak dapat masuk.

Setiap akun juga membawa tingkat akses (`level`), seperti di aplikasi lama:

| Peran lama | `level` | Cakupan data |
| --- | --- | --- |
| `SYSTEM ADMIN`, `ADMIN PTD`, `ADMIN DPD` | 3 | semua kecamatan |
| `ADMIN KECAMATAN`, `TARBIYAH KECAMATAN` | 2 | kecamatan sendiri |
| `ANGGOTA` | 1 | data yang ia buat sendiri |

Admin kecamatannya diambil dari kolom `NAMA KECAMATAN` di tab `USERS`. Akun
lokal yang sudah ada sebelumnya mendapat tingkatnya saat login pertama kali
setelah perubahan ini.

Hak managerial mengikuti `level`: hanya `level` 3 (admin penuh) yang boleh
menghapus data, mengelola master aktivitas, dan membuka Pengaturan. Admin
kecamatan (`level` 2) hanya mengelola data kecamatannya sendiri.

### Daftar anggota mandiri (self-registration)

Halaman login punya tombol **Belum punya akun? Daftar anggota**. Formulirnya
minta username, kata sandi, nama, kecamatan, dan dua pertanyaan keamanan.

Akun baru ditulis ke tab `USERS` dengan layout yang sama seperti aplikasi lama:
peran `ANGGOTA`, `AKTIF` = `Y`, serta password dan kedua jawaban keamanan
disimpan sebagai hash SHA-256 tanpa salt. Jadi akun itu langsung aktif dan bisa
dipakai di kedua aplikasi. Pada login pertama, akun dimindahkan ke penyimpanan
lokal dengan hash PBKDF2 bersalt.

Syaratnya: `GOOGLE_SHEET_ID` dan `GOOGLE_KEY_FILE` terisi (agar tab `USERS`
bisa ditulis) dan `LEGACY_AUTH_MODE` bukan `off`. Tanpa itu, endpoint
pendaftaran menjawab `503` dan tombol pendaftaran disembunyikan.

Endpoint publiknya:

| Endpoint | Fungsi |
| --- | --- |
| `GET /api/register/options` | daftar pertanyaan keamanan dan kecamatan |
| `POST /api/register` | membuat akun anggota baru |

Pendaftaran dibatasi 10 percobaan per alamat dalam 10 menit, dan username
yang sudah dipakai ditolak dengan `409`.

### Mengimpor data dari Google Form

Untuk spreadsheet yang berisi respons Google Form, gunakan perkakas impor.
Tagian dipetakan otomatis berdasarkan nama kolom, dan impor bersifat
*tambah-saja*: baris yang nomornya sudah ada dilewati, sehingga aman
dijalankan berkali-kali.

```bash
cd backend
go run ./cmd/impor-responses --dry-run              # hanya rencana
go run ./cmd/impor-responses --apply --report ../data/impor-responses-report.json
```

Tanpa `--apply` tidak ada yang ditulis ke spreadsheet. Rincian rencana dan
baris yang perlu diperiksa ada di file laporan.

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

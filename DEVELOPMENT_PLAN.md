# Development Plan — Rekrutmen Tangsel

**Status dokumen:** aktif, diperbarui per checkpoint
**Terakhir diperbarui:** 2026-11-14
**Commit baseline:** `20ddf68` (update aplikasi)
**Rujukan:** [`PRD_Rekrutmen_Tangsel_MVP.md`](PRD_Rekrutmen_Tangsel_MVP.md) v0.1,
[`API_CONTRACT.md`](API_CONTRACT.md)

---

## 1. Cara memakai dokumen ini

Dokumen ini adalah **catatan checkpoint**. Setiap item punya status yang bisa
diperiksa ulang kapan saja:

| Penanda | Arti |
|---|---|
| ✅ | **Selesai** — terimplementasi, sudah diverifikasi (build/test/runtime), sudah ter-commit |
| 🔧 | **Sedang dikerjakan** — ada perubahan di working tree, belum ter-commit |
| ⬜ | **Belum** — belum dikerjakan |
| ⏸️ | **Ditunda** — sengaja tidak dikerjakan di fase ini, dengan alasan |

Status **hanya** berubah menjadi ✅ setelah item terverifikasi lewat perintah
yang tertulis di bagian **Verifikasi**, lalu ter-commit. Kalau proses terputus
(token habis / sesi mati), checkpoint terakhir yang ✅ adalah titik kembali —
periksa `git log --oneline` dan `git status` untuk memastikan tidak ada
perubahan yang menggantung.

> **Konvensi:** nomor item (`D-xx`, `T-xx`, `P-xx`) dipakai silang antar tiga
> dokumen. `D` = development, `T` = testing, `P` = deployment. Satu item yang
> disebut di dokumen mana pun punya nomor yang sama.

---

## 2. Ringkasan

| Fase | Selesai | Sedang | Belum | Total |
|---|---|---|---|---|
| F1 — Fondasi & model data | 8 | 0 | 0 | 8 |
| F2 — API backend | 23 | 0 | 4 | 27 |
| F3 — Penyimpanan | 10 | 0 | 2 | 12 |
| F4 — Frontend | 27 | 0 | 4 | 31 |
| F5 — Keamanan | 4 | 0 | 7 | 11 |
| F6 — Kualitas & CI | 2 | 1 | 5 | 8 |
| F7 — Fitur pelengkap PRD | 9 | 2 | 8 | 19 |
| F8 — Granularity role (paritas app lama) | 4 | 4 | 0 | 8 |
| **Total** | **87** | **7** | **30** | **124** |

> Item D-176 dan sub-item D-176…D-186 di §9.1 dihitung **satu kali** di F7
> (D-176 sebagai item induk; D-177…D-186 sebagai rinciannya).
> Item D-187…D-194 di §9.2 dihitung di F8.
> Status codebase: **build hijau, test hijau, siap didemokan** — belum siap
> diproduksi karena F5 (keamanan) dan F6 (CI) masih kosong.
> Item **🔧** = sudah ada kodenya di working tree dan teruji, tetapi belum
> ter-commit; statusnya menjadi ✅ setelah commit (lihat §11).

---

## 3. F1 — Fondasi & model data ✅

| # | Item | Status | Verifikasi | Catatan |
|---|---|---|---|---|
| D-01 | Struktur repo (backend Go, frontend Vue) | ✅ | `ls` | `backend/`, `frontend/` |
| D-02 | Model domain rekrutmen | ✅ | `backend/internal/model/model.go` | 14 field PRD + 4 field audit |
| D-03 | Kontrak API tertulis | ✅ | `API_CONTRACT.md` | 17 endpoint terdokumentasi |
| D-04 | Master wilayah 7 kecamatan / 54 kelurahan | ✅ | `GET /api/wilayah` → 7/54 | `internal/seed/wilayah.json` (embed) |
| D-05 | Master aktivitas seed awal | ✅ | `POST /api/aktivitas` → 201 | 18 kategori bawaan |
| D-06 | Aturan validasi terpusat | ✅ | `internal/validate/validate.go` | dipakai backend **dan** ditiru frontend |
| D-07 | `.gitignore` aman | ✅ | `git check-ignore -v` | `data/`, `*.env`, `*.secret` diabaikan |
| D-08 | Tidak ada secret bocor di file terlacak | ✅ | `grep -rn "ghp_\|password\|api_key"` | hanya default demo di `auth/users.go` |

### Decision log

| Tanggal | Keputusan | Alasan |
|---|---|---|
| 2026-11-06 | Frontend **tidak** memakai library validasi; aturan diduplikasi dari `constants.js` | API sudah jadi sumber kebenaran; menambah dependensi hanya memperbesar bundle |
| 2026-11-06 | `no_wa` disimpan sebagai `string`, bukan angka | PRD §4 — agar angka nol di awal tidak hilang (`validate_test.go:70` mengunci ini) |
| 2026-11-06 | UPLS/UPLM jadi **satu** field `upl` dengan radio group | PRD §10 butir 176 — "user secara opsional bisa memilih salah satu" |
| 2026-11-09 | Enum `metode` jadi **3 opsi**: `Halaqoh`, `Program Antara`, `Dakwah Fardiyah` | Koreksi pemilik produk; "Program" + "Antara" sebelumnya terpisah. Nilai lama tetap diterima dan dinormalisasi, tidak ada data hilang |
| 2026-11-09 | Normalisasi enum di **4 lapisan** (validasi, ketiga store, frontend) | Data lama bisa datang dari storage yang sudah terisi, bukan hanya dari form baru |

---

## 4. F2 — API backend ✅

Handler: `backend/internal/api/api.go`. Semua 17 route terdaftar di
`api.go:197-224`.

| # | Item | Status | Verifikasi |
|---|---|---|---|
| D-09 | `POST /api/login` | ✅ | login admin & petugas → 200 |
| D-10 | `POST /api/logout` | ✅ | `TestLogoutClearsSession` |
| D-11 | `GET /api/me` | ✅ | `TestLoginSetsSessionAndMeReturnsUser` |
| D-12 | `GET /api/wilayah` | ✅ | 7 kecamatan / 54 kelurahan |
| D-13 | `GET/POST/PUT/DELETE /api/aktivitas` | ✅ | `TestAktivitasMasterIsAdminOnly` |
| D-14 | `GET /api/rekrutmen` (list + filter + paginasi) | ✅ | `TestSearchAndFilter` |
| D-15 | `GET /api/rekrutmen/:id` | ✅ | `TestRekrutmenCRUDLifecycle` |
| D-16 | `POST /api/rekrutmen` | ✅ | 201 + isi `created_at`/`created_by` |
| D-17 | `PUT /api/rekrutmen/:id` | ✅ | `TestUpdatePreservesCreationAuditFields` |
| D-18 | `DELETE /api/rekrutmen/:id` | ✅ | 200 lalu 404 |
| D-19 | `GET /api/rekrutmen/check-wa` | ✅ | `TestCheckWhatsAppReportsExistence` |
| D-20 | `GET/PUT /api/settings` | ✅ | `TestSettingsSwitchProvider` |
| D-21 | `GET /api/health` | ✅ | `{"ok":true,"provider":"local"}` |
| D-22 | RBAC: `requireAuth` / `requireAdmin` | ✅ | petugas → 403, tanpa auth → 401 |
| D-23 | Error per-field (`{error, fields}`) | ✅ | `TestValidationErrorsAreReportedPerField` |
| D-24 | Konflik 409 untuk no_wa duplikat | ✅ | `TestDuplicateWhatsAppReturnsConflict` |
| D-25 | 409 saat aktivitas masih dipakai | ✅ | `TestAktivitasInUseCannotBeDeleted` |
| D-26 | 404 untuk rekrutmen/aktivitas tak ada | ✅ | diuji di lifecycle test |
| D-27 | CORS untuk dev server Vite | ✅ | `withCORS` `api.go:229` |
| D-28 | Pesan validasi berbahasa Indonesia | ✅ | konsisten di seluruh respons |
| D-29 | **Fix: `PUT` menghapus jejak audit** | ✅ | `api.go:531-545`; regresi test `api_test.go:196` |
| D-30 | Sinkronisasi frontend↔backend saat field berubah | ✅ | `API_CONTRACT.md` diperbarui bersama kode |
| D-104 | **Fix: enum `metode` — "Program" + "Antara" → "Program Antara"** | ✅ | `model.go:9` (3 opsi); normalisasi legacy di `validate.go`, ketiga store, dan frontend; regresi test `TestMetodeProgramAntaraIsStoredAndFilterable` |
| D-31 | Endpoint ganti password | ✅ | `PUT /api/me/password` (`api.go`), `auth.Manager.ChangePassword` (`auth.go`); semua role boleh, password lama wajib benar |
| D-32 | Manajemen pengguna (tambah/hapus petugas) | ⬜ | user hanya dari `auth/users.go` |
| D-33 | Rate limiting pada `POST /api/login` | ⬜ | lihat D-83 |
| D-34 | Request ID / logging terstruktur | ⬜ | server pakai `log.Printf` polos |

---

## 5. F3 — Penyimpanan ✅

Antarmuka: `backend/internal/store/store.go`. Tiga provider mengimplementasikan
`Store` yang sama.

| # | Item | Status | Verifikasi |
|---|---|---|---|
| D-35 | Antarmuka `Store` + sentinel error | ✅ | `store.go:23` |
| D-36 | Provider `local` (JSON di disk) | ✅ | coverage 86.5% |
| D-37 | Provider `supabase` | ✅ | coverage 51.8% |
| D-38 | Provider `google_sheets` | ✅ | coverage 78.7% |
| D-39 | Tab + header Sheets dibuat otomatis | ✅ | `TestSheetsEnsureTabsCreatesHeaders` |
| D-40 | Unikensi no_wa di ketiga provider | ✅ | `ErrDuplicateWA` → 409 |
| D-41 | Cascade rename master aktivitas | ✅ | `TestAktivitasRenameCascadesToRekrutmen` |
| D-42 | Aktivasi aktivitas `aktif=false` | ✅ | `PUT /api/aktivitas/:id` |
| D-43 | Skema Supabase + RLS + trigger `updated_at` | ✅ | `backend/supabase/schema.sql` |
| D-44 | Halaman `settings.json` untuk provider aktif | ✅ | `internal/settings/settings.go` |
| D-45 | **Ekspor/import data** | ⬜ | PRD §7 minta prosedur ekspor/impor saat ganti provider |
| D-46 | **Verifikasi perpindahan provider** | ⬜ | PRD §7 minta verifikasi sebelum provider dipakai |

> **Catatan D-45/D-46:** README hanya memberi peringatan "meng ganti provider
> tidak memindahkan data". PRD §7 eksplisit meminta prosedur ekspor/impor +
> verifikasi. Ini gap yang nyata, lihat F7.

---

## 6. F4 — Frontend ✅

Vue 3 + Vue Router + Vite. Sumber: `frontend/src/`.

| # | Item | Status | Verifikasi |
|---|---|---|---|
| D-47 | Bootstrap aplikasi + router | ✅ | `main.js`, `router.js` |
| D-48 | Route guard (login + admin-only) | ✅ | `router.js:26-41` |
| D-49 | Judul tab dinamis per route | ✅ | `router.js:43-47` |
| D-50 | View Login | ✅ | `LoginView.vue` |
| D-51 | View List (tabel + kartu) | ✅ | `ListView.vue` |
| D-52 | View Form (tambah + ubah) | ✅ | `FormView.vue` |
| D-53 | View Master Aktivitas | ✅ | `AktivitasView.vue` |
| D-54 | View Pengaturan (ganti provider) | ✅ | `SettingsView.vue` |
| D-55 | Komponen Toast | ✅ | `components/Toast.vue` |
| D-56 | Komponen ConfirmDialog | ✅ | `components/ConfirmDialog.vue` |
| D-57 | Wrapper API + class error | ✅ | `api.js` |
| D-58 | Semua 14 field PRD ada di form | ✅ | `FormView.vue:268-508` |
| D-59 | Validasi klien (nama, no_wa, RT/RW, enum wajib) | ✅ | `FormView.vue:83-104` |
| D-60 | Kelompok form: Identitas / Wilayah / Info Rekrutmen / Keterangan | ✅ | `FormView.vue:261,320,395` |
| D-61 | Kelurahan mengikuti kecamatan + reset | ✅ | `FormView.vue` computed |
| D-62 | UPLS/UPLM radio single-choice | ✅ | `FormView.vue:507-508` |
| D-63 | Peringatan duplikat WA sebelum simpan | ✅ | `checkDuplicateWa()` → `check-wa` |
| D-64 | Notifikasi sukses/gagal | ✅ | `toast.js` |
| D-65 | Loading state saat menyimpan | ✅ | `FormView.vue` `saving` |
| D-66 | Hapus dengan konfirmasi | ✅ | `ConfirmDialog` + `ListView.vue:137` |
| D-67 | Pencarian nama/nomor WA (debounce) | ✅ | `ListView.vue:30-40` |
| D-68 | Filter kecamatan/kelurahan/metode/status | ✅ | `ListView.vue:24-27` |
| D-69 | Paginasi | ✅ | `ListView.vue:52-55,125` |
| D-70 | Responsif: tabel → kartu (<768px) | ✅ | `style.css:1286` |
| D-71 | Ganti provider di menu Pengaturan (admin) | ✅ | `SettingsView.vue` |
| D-72 | Master aktivitas dikelola admin | ✅ | `AktivitasView.vue` |
| D-73 | Error handling + pesan fallback | ✅ | `ApiError` di `api.js` |
| D-74 | **Test runner frontend** | ⬜ | `package.json` tidak punya vitest/jest |
| D-75 | **Manajemen pengguna di UI** | ⬜ | lihat D-32 |
| D-76 | **Ekspor data dari UI** | ⬜ | lihat D-45 |
| D-77 | **Skor aksesibilitas (a11y) audit** | ⬜ | label & `for` sudah ada, audit belum |

---

## 7. F5 — Keamanan 🔐

| # | Item | Status | Verifikasi |
|---|---|---|---|
| D-78 | Hash password PBKDF2 + salt acak | ✅ | `auth.go:58` |
| D-79 | Timing-equalization saat username tak ada | ✅ | `auth.go:73-75` |
| D-80 | Cookie `HttpOnly` + `SameSite=Lax` | ✅ | `auth.go:117-118` |
| D-81 | RLS Supabase menolak akses `anon` | ✅ | `schema.sql:66-73` |
| D-82 | **Flag `Secure` pada cookie sesi** | ⬜ | `auth.go:113` & `auth.go:126` — wajib saat HTTPS |
| D-83 | **Rate limit login (anti brute-force)** | ⬜ | `api.go:251` — endpoint login terbuka |
| D-84 | **Endpoint ganti password** | ✅ | lihat D-31 |
| D-85 | **Hapus password default / paksa ganti saat login pertama** | ⬜ | `auth/users.go:20-22` |
| D-86 | **CORS allow-list** | ⬜ | `api.go:234` mencerminkan `Origin` apa pun |
| D-87 | **Header keamanan (CSP, X-Frame-Options, HSTS)** | ⬜ | belum ada middleware |
| D-88 | **Rate limit / pagination batas atas `limit`** | ⬜ | `api.go:405` — `limit` tanpa cap |

> D-82 dan D-87 adalah **blocker produksi**: cookie tanpa `Secure` dan tanpa
> HSTS berarti sesi bisa dicuri lewat jaringan tidak tepercaya.

---

## 8. F6 — Kualitas & CI

| # | Item | Status | Verifikasi |
|---|---|---|---|
| D-89 | `go vet` bersih | ✅ | `go vet ./...` |
| D-90 | `gofmt` konsisten | ✅ | `gofmt -l backend/` |
| D-91 | **GitHub Actions CI** | ⬜ | `.github/` belum ada |
| D-92 | **Lint frontend** | ⬜ | tidak ada ESLint |
| D-93 | **Test frontend (unit)** | ⬜ | lihat D-74 |
| D-94 | **Test E2E (Playwright)** | 🔧 | `playwright.config.ts` + `tests/e2e/browser-check.spec.ts` (hanya cek browser ke `example.com`); skenario aplikasi belum |
| D-95 | **Pinning versi toolchain** | ⬜ | `go.mod` butuh `go 1.24`; Node belum di-pin |
| D-96 | **Pre-commit hook** | ⬜ | `.claude/settings.json` hanya berisi hook graphify |

---

## 9. F7 — Fitur pelengkap PRD

| # | Item | Status | Keterkaitan PRD |
|---|---|---|---|
| D-97 | **Ekspor/import + verifikasi ganti provider** | ⬜ | §7, §10 |
| D-98 | **Manajemen pengguna (CRUD petugas)** | ⬜ | §2, §9 |
| D-99 | **Dashboard & laporan analitik** | ✅ | `internal/api/dashboard.go`, `views/DashboardView.vue` — lihat rencana §14 (D-207…D-215) |
| D-100 | **Impor massal Excel/CSV** | ⬜ | §11 — *dinaikkan prioritas*, lihat rencana §9.1 (D-176) |
| D-101 | **Integrasi WhatsApp kirim pesan** | ⬜ | §11 — *di luar scope MVP* |
| D-102 | **Workflow persetujuan berjenjang** | ⬜ | §11 — *di luar scope MVP* |
| D-103 | **Sinkronisasi dua arah Supabase↔Sheets** | ⬜ | §11 — *di luar scope MVP* |
| D-165 | **Login memakai akun aplikasi lama (tab `USERS`)** | ✅ | `internal/store/sheets/users.go`, `internal/auth` |
| D-166 | **Pemetaan peran aplikasi lama → peran baru** | ✅ | `SYSTEM ADMIN`/`ADMIN PTD`/`ADMIN KECAMATAN` → `admin`, `ANGGOTA` → `petugas` |
| D-167 | **Pemindahan akun lama ke hash PBKDF2 saat login pertama** | ✅ | `Manager.MigrateLegacyUser`, mode `upgrade` |
| D-168 | **Mode login warisan (`upgrade`/`sheet`/`off`)** | ✅ | env `LEGACY_AUTH_MODE` |
| D-169 | **Penolakan akun tidak aktif / tanpa hash** | ✅ | `AKTIF` harus `Y`, `PASSWORD_HASH` wajib terisi |
| D-170 | **Batas laju login warisan** | ⬜ | D-83 masih terbuka; hash SHA-256 tanpa salt mudah dipecah |
| D-171 | **Migrasi 161 akun aplikasi lama ke direktori pengguna** | 🔧 | sebagian: akun aktif otomatis saat login pertama; 8 akun tanpa hash belum bisa masuk |
| D-172 | **Daftar anggota mandiri (self-registration)** | ✅ | `POST /api/register`, `internal/store/sheets/register.go`, `LoginView.vue` |
| D-173 | **Akun baru memakai hash SHA-256 yang sama dengan aplikasi lama** | ✅ | password & jawaban keamanan di-hash SHA-256 tanpa salt,role `ANGGOTA`, `AKTIF`=`Y` |
| D-174 | **Batas laju pendaftaran mandiri** | ✅ | 10 percobaan per alamat per 10 menit (`internal/api/register.go`) |
| D-175 | **Pilihan pertanyaan keamanan & kecamatan untuk form daftar** | ✅ | `GET /api/register/options` |
| D-176 | **Upload batch Excel/CSV untuk anggota (paritas app lama)** | 🔧 | rencana di §9.1 — backend selesai, UI & dokumentasi belum |

> D-99 s.d. D-103 secara eksplisit **di luar scope MVP**. Item-item ini dicatat
> agar tidak hilang, bukan untuk dikerjakan sekarang.

> **Catatan 2026-11-14:** D-100 (impor massal) **tidak lagi "di luar scope MVP"**.
> Aplikasi lama sudah punya fitur ini dan dipakai harian oleh anggota
> (`old-app/code.gs:732` `uploadRecruitments`, `old-app/index.html:322` tombol
> Upload Excel) — di aplikasi baru fitur ini hilang, sehingga anggota hanya
> input satu per satu. Rencana rinci ada di §9.1.

---

## 9.1 Rencana — Upload batch Excel untuk anggota (paritas app lama)

**Status:** 🔧 sebagian — backend selesai & teruji, UI & dokumentasi belum ·
**Tanggal:** 2026-11-14 · **Sumber:** analisis `old-app/code.gs` +
`old-app/index.html`

### 9.1.1 Perilaku aplikasi lama yang harus dipertahankan

| Aspek | Perilaku app lama (sumber) |
|---|---|
| Siapa boleh | Hanya role `ANGGOTA` — `requireRole_(me, [ANGGOTA])` (`code.gs:735`), tombol disembunyikan untuk role lain (`index.html:530-536`) |
| Format | Excel `.xlsx/.xls/.csv` via SheetJS (`index.html:325`), template 14 kolom (`index.html:992`) |
| Template | `NAMA TARGET, NO WHATSAPP, ALAMAT, RT, RW, KECAMATAN, KELURAHAN, UPLS, UPLM, AKTIFITAS, TINGKAT PENERIMAAN, METODE, STATUS, KETERANGAN` + sheet `PETUNJUK` (`index.html:1004-1029`) |
| Batas | Maksimal **1.000 baris** per upload (`code.gs:739`, `index.html:1080`) |
| Validasi kolom | 8 kolom wajib; `UPLS` XOR `UPLM`; `KECAMATAN` harus ada di `TANGSEL_WILAYAH`; `KELURAHAN` harus milik kecamatan itu; `AKTIFITAS` harus dari daftar opsi (`code.gs:997-1009`) |
| Sifat kegagalan | **All-or-nothing**: ada 1 baris salah → tidak ada baris tersimpan, pesan maks 30 error (`code.gs:766-780`) |
| Field protektif | Diisi server dari session, tak bisa dipalsukan dari Excel: `ID` (UUID), `NAMA ANGGOTA` = `me.anggota`, `NAMA KECAMATAN` = `me.kecamatan`, `CREATED_AT`/`UPDATED_AT` = now, `CREATED_BY` = `me.username` (`code.gs:795-809`) |
| Penulisan | Satu operasi tulis di bawah `LockService` (tunggu 30 detik) (`code.gs:785-813`) |
| Notifikasi | Alert jumlah data tersimpan, lalu refresh dashboard + tabel (`index.html:1094-1097`) |

> Catatan: app lama **tidak** mengecek duplikat `no_wa` saat upload. Server
> baru sudah punya `store.ErrDuplicateWA`, jadi aturan baru harus diputuskan
> (lihat 9.1.3).

### 9.1.2 Prinsip implementasi

1. **Parsing di client, validasi di server** — sama seperti app lama: browser
   membaca file menjadi array objek, server memvalidasi ulang dan menyimpan.
   Tidak perlu parser Excel di Go.
2. **Tidak ada parser Excel di backend** → **keputusan berubah 2026-11-14:**
   pustaka `xlsx` **tidak** dipasang ke `frontend/package.json`. Berkas
   `frontend/src/excel.js` memuat SheetJS **dari CDN** saat file `.xlsx`
   benar-benar dipilih (`cdn.sheetjs.com`, sama seperti app lama), dan
   berkas `.csv` dibaca tanpa pustaka apa pun. Bundle tetap kecil, tidak ada
   dependensi npm baru — tetapi butuh internet saat memakai berkas Excel.
3. **Field protektif diisi server** — `created_by`, `updated_by`, `created_at`,
   `updated_at`, dan kolom `NAMA ANGGOTA`/uploader diambil dari session, bukan
   dari payload.
4. **Semua-or-nothing**, dengan laporan error per baris (maks 30 ditampilkan).
5. **Scope konsisten** — baris hasil upload otomatis berada di scope
   penginput (anggota = `created_by` = username), sehingga anggota langsung
   melihat datanya sendiri.

### 9.1.3 Keputusan (sudah diambil 2026-11-14, sudah diterapkan di kode)

| # | Pertanyaan | Keputusan | Penerapan |
|---|---|---|---|
| K-1 | Siapa boleh upload? | **Anggota (level 1) saja** — paritas app lama | `requireUpload` → 403 untuk selain anggota (`api/batch.go`) |
| K-2 | `no_wa` duplikat terhadap data yang sudah ada | **Lewati baris duplikat + laporkan di ringkasan** | `skipped[]` pada respons 201 |
| K-3 | Duplikat di dalam file yang sama | **Gagal** (hanya baris itu yang salah) | pesan "nomor whatsapp sama dengan baris N" |
| K-4 | Batas ukuran file | 5 MB | `MaxBatchBody` (5 << 20) |
| K-5 | Nama endpoint | `POST /api/rekrutmen/batch` | terdaftar di `api.go` |
| K-6 | Header `X-UUpload-Mode` | **Tidak dipakai** — tidak ada kebutuhan mode berbeda | body JSON polos berisi array baris |

### 9.1.4 Item pekerjaan

| # | Item | Status | Verifikasi |
|---|---|---|---|
| D-176 | Kontrak endpoint `POST /api/rekrutmen/batch` (body: array baris) | 🔧 | route ada di `api.go`; `API_CONTRACT.md` **belum** diperbarui |
| D-177 | Hak akses endpoint (level 1) + penolakan 403 | 🔧 | `requireUpload`; `TestBatchIsOnlyForMembers` lulus |
| D-178 | Validasi per baris all-or-nothing + laporan error bernomor baris | 🔧 | `TestBatchSavesNothingWhenOneRowIsWrong` lulus |
| D-179 | Pemetaan kolom Excel (alias) → `model.Rekrutmen` | 🔧 | peta lokal di `api/batch.go` (**bukan** `AppColumnSpec()`); alias app lama & nama field app ini diterima |
| D-180 | Field protektif diisi server (created_by/updated_by/timestamp/uploader) | 🔧 | `TestBatchIgnoresTheOwnerColumnsOfTheFile` lulus |
| D-181 | Parser Excel di frontend + unduh template | 🔧 | `frontend/src/excel.js` (CSV tanpa pustaka, `.xlsx` via CDN SheetJS) dipakai `components/BatchUpload.vue` |
| D-182 | UI: tombol Upload + Template (khusus anggota), progress, ringkasan hasil | 🔧 | `frontend/src/components/BatchUpload.vue` + `canUpload()` di `auth.js`; hanya level 1 yang melihat tombolnya |
| D-183 | Batas: maks 1.000 baris, 5 MB, rate limit sederhana | 🔧 | `MaxBatchRows`, `MaxBatchBody`, `uploadLimiter` (5×/10 mnt); **timeout** belum ada |
| D-184 | Penanganan `no_wa` duplikat (K-2/K-3) | 🔧 | `TestBatchSkipsANumberThatIsAlreadyStored`, `TestBatchRejectsTwoRowsWithTheSameNumber` lulus |
| D-185 | Test otomatis: validasi, scope, field protektif, all-or-nothing | 🔧 | 8 test di `api/batch_test.go` lulus; test `CreateRekrutmenBatch` per provider belum |
| D-186 | Dokumentasi: README, API_CONTRACT, ANALISIS_AKSES_DATA (paritas app lama) | ⬜ | ketiga dokumen belum disentuh |
| D-195 | Ejaan enum `Metode` dari app lama (`Halaqah`, `Dakwah Fardliyah`) diterima saat upload dan dinormalkan | 🔧 | `model.MetodeLama` + `constants.js` `METODE_LAMA`; `TestBatchAcceptsTheSpellingsOfThePreviousApplication` |

### 9.1.5 Urutan pengerjaan —Kemajuan

| # | Langkah | Status |
|---|---|---|
| 1 | D-179 + D-180 (model & pemetaan kolom) | 🔧 selesai (peta lokal, bukan `AppColumnSpec()`) |
| 2 | D-176 + D-177 + D-178 (endpoint + hak akses + validasi) + D-184 (duplikat) | 🔧 selesai, `API_CONTRACT.md` menyusul |
| 3 | D-185 test backend | 🔧 sebagian — 8 test lulus |
| 4 | D-181 + D-182 (frontend: parser, template, tombol, ringkasan) | 🔧 selesai — `BatchUpload.vue` terpasang di `ListView.vue` |
| 5 | D-183 (batas & rate limit), D-186 (dokumentasi) | 🔧 sebagian / ⬜ |

### 9.1.6 Peringatan / risiko

- **Bukan menambah dependency npm** — keputusan berubah: SheetJS dimuat dari
  CDN saat dipakai, jadi bundle tidak naik. Risikonya: berkas `.xlsx` butuh
  internet; berkas `.csv` tetap jalan tanpa internet.
- **Penulisan massal ke Google Sheets** bisa melewati kuota & batas 1.000 baris
  per update → tulis bertahap (mis. 200 baris per panggilan) dengan jeda.
- `store.ErrDuplicateWA` berbeda dari app lama → lihat K-2.
- Endpoint baru membuka jalur input baru → wajib rate limit (D-183) dan
  validasi ulang di server, bukan hanya di client.
- **Kode frontend sudah tersambung** — `frontend/src/excel.js` dipakai
  `frontend/src/components/BatchUpload.vue`, yang terpasang di `ListView.vue`
  dan hanya tampil untuk anggota (level 1). Alur upload sudah diuji lewat
  server lokal: unggah valid → `created`; baris salah → 400 + daftar error
  bernomor baris; nomor ganda → 400; admin penuh → 403.

### 9.2 Granularity role — pemulangan hak akses aplikasi lama

**Status:** 🔧 sebagian · **Tanggal:** 2026-11-14 · **Sumber:** analisis
`old-app/code.gs` + `old-app/index.html`, dan `ANALISIS_AKSES_DATA.md` §5.1–5.3

#### 9.2.1 Perilaku aplikasi lama yang dipulangkan

| Aspek | Perilaku app lama (sumber) |
|---|---|
| Tingkat akses | 3 tingkat: `ANGGOTA` (1), `ADMIN KECAMATAN` (2), `ADMIN PTD`/`SYSTEM ADMIN` (3) |
| Scope baca | anggota → data sendiri; admin kecamatan → data kecamatannya; admin penuh → semua |
| Hak hapus | `requireRole_(me, [ANGGOTA])` **+** `CREATED_BY` harus sama (`old-app/code.gs:873-882`); tombol Hapus hanya untuk `ANGGOTA` (`old-app/index.html:896`) |
| Master & pengaturan | khusus admin penuh |

> Catatan: aplikasi lama **tidak** memberi hak hapus ke admin. Versi ini
> menambah hak hapus untuk admin penuh (level 3) dan tetap menutupnya untuk
> admin kecamatan — keputusan yang lebih longgar bagi admin penuh, dicatat di
> `ANALISIS_AKSES_DATA.md` §5.3.

#### 9.2.2 Item pekerjaan

| # | Item | Status | Verifikasi |
|---|---|---|---|
| D-187 | Model level akses (1/2/3) + `ScopeFor` per level | ✅ | `model/model.go` (commit `9ea813d`) |
| D-188 | Pemetaan peran app lama → level (`ANGGOTA`→1, `ADMIN KECAMATAN`/`TARBIYAH KECAMATAN`→2, `ADMIN PTD`/`SYSTEM ADMIN`/`ADMIN DPD`→3) | ✅ | `store/sheets/users.go:53-71` |
| D-189 | Fallback akun lama tanpa level → level diturunkan dari `role` | ✅ | `model.User.LevelOrDefault` |
| D-190 | Rekaman di luar scope dijawab **404** (bukan 403) — cegah kebocoran keberadaan data | ✅ | `canAccess` di `api/api.go` (commit `9ea813d`) |
| D-191 | `requireDelete`: anggota boleh hapus data sendiri, admin penuh boleh hapus semua, admin kecamatan tidak | 🔧 | `api/api.go` (`requireDelete`) |
| D-192 | Frontend `canDelete()` + tombol Hapus mengikuti level | 🔧 | `frontend/src/auth.js`, `frontend/src/views/ListView.vue` |
| D-193 | Test otomatis hak hapus per level | 🔧 | `api/scope_test.go` `TestMemberCanDeleteOnlyItsOwnRecord` (menggantikan `TestOnlyAdminCanDeleteRecord`) |
| D-194 | Dokumentasi hak hapus | 🔧 | `README.md`, `API_CONTRACT.md`, `ANALISIS_AKSES_DATA.md` §5.3 |

> Item D-191…D-194 sudah ada di working tree dan **belum ter-commit**; setelah
> commit statusnya menjadi ✅. Test lama `TestOnlyAdminCanDeleteRecord`
> dihapus karena bertentangan dengan paritas app lama.

---

## 10. Traceability ke PRD

Setiap kriteria penerimaan PRD §10 dipetakan ke implementasi dan test.

| # | Kriteria PRD | D- | Uji otomatis |
|---|---|---|---|
| 1 | Tambah, lihat, ubah, hapus data | D-14…D-18 | `TestRekrutmenCRUDLifecycle` |
| 2 | No WA wajib, tolak non-digit | D-16, D-36…D-38 | `TestWhatsAppMustBeDigitsOnly` |
| 3 | RT/RW tolak non-digit | D-36…D-38 | `TestRTAndRWMustBeDigitsOnly` |
| 4 | Kecamatan wajib, hanya 7 Tangsel | D-06, D-04 | `TestKecamatanMustBeOneOfSeven` |
| 5 | Kelurahan mengikuti kecamatan | D-61 | `TestKelurahanMustBelongToKecamatan` |
| 6 | Aktivitas/tingkat/metode/status wajib dari opsi | D-06 | `TestEnumValuesMustBeValid` |
| 7 | UPLS & UPLM opsional, satu pilihan | D-62 | `TestUPLAcceptsEmptyUPLSorUPLM` |
| 8 | Tersimpan di provider aktif | D-35…D-38 | test ketiga provider |
| 9 | Dipakai di desktop & mobile | D-70 | perlu uji manual / E2E (D-94) |
| 10 | Nama orang wajib | D-06, D-59 | `TestRequiredFields` |
| 11 | UPLS/UPLM salah satu | D-62 | `TestUPLAcceptsEmptyUPLSorUPLM` |
| 12 | No WA unik | D-24, D-40 | `TestDuplicateWhatsAppReturnsConflict` |
| 13 | Admin ganti provider | D-20, D-71 | `TestSettingsSwitchProvider` |
| 14 | Master aktivitas dikelola admin | D-13, D-72 | `TestAktivitasMasterIsAdminOnly` |

---

## 11. Checkpoint terakhir

| Item | Nilai |
|---|---|
| Commit terakhir | `20ddf68` (update aplikasi) |
| Remote `main` | belum dicek pada pembaruan ini |
| Working tree | **belum bersih** — tiga kelompok perubahan:<br>1. **Upload batch Excel (§9.1)** — `backend/internal/api/batch.go`, `backend/internal/api/batch_test.go`, `frontend/src/excel.js`, `frontend/src/components/BatchUpload.vue` (baru), `CreateRekrutmenBatch` di ketiga store, route + `api.uploadBatch()` di frontend, `canUpload()` di `auth.js`.<br>2. **Granularity role (§9.2)** — `requireDelete` di backend (anggota boleh hapus miliknya, admin penuh boleh hapus semua, admin kecamatan tidak) + `canDelete()` di `frontend/src/auth.js` + `ListView.vue`.<br>3. **Ejaan enum warisan (D-195)** — `model.MetodeLama` + `constants.js`. |
| `go build ./...` | ✅ |
| `go vet ./...` | ✅ bersih |
| `go test ./...` | ✅ semua paket lulus (termasuk 8 test batch) |
| `npm run build` | ✅ 5.09 s, 146 kB (SheetJS tidak masuk bundle) |
| Test regression untuk D-29 | ✅ lulus; terbukti menangkap bug bila fix dibatalkan |

> **Catatan 2026-11-14:** item bertanda 🔧 (§9.1 D-176…D-185 & D-195, §9.2
> D-191…D-194, D-171) sudah terimplementasi dan teruji, tetapi belum
> ter-commit. Setelah commit, ubah statusnya menjadi ✅. Item bertanda ⬜
> tetap belum dikerjakan.

### Titik kembali kalau terputus

1. `git log --oneline` — cari commit `✅` terakhir.
2. `git status` — pastikan tidak ada file `🔧` yang menggantung; kalau ada,
   `git diff` untuk melihat apakah layak di-commit atau di-buang.
3. Jalankan `go test ./...` dan `npm run build` untuk memastikan tree sehat.
4. Lanjutkan dari item pertama yang statusnya bukan ✅ di bagian 2.

---

## 12. Urutan pengerjaan berikutnya

Prioritas dihitung dari (risiko produksi) × (regulatori PRD).

| Urutan | Item | Alasan |
|---|---|---|
| 1 | D-82 flag `Secure` + D-87 header keamanan | blocker produksi |
| 2 | D-84 + D-85 ganti password / hapus default | kredensial default = risiko |
| 3 | D-176 upload batch Excel untuk anggota (§9.1) | backend selesai; UI (D-182) & dokumentasi (D-186) menyusul |
| 4 | D-97 ekspor/import + verifikasi provider | gap PRD §7 yang eksplisit |
| 5 | D-91 CI + D-93 test frontend | mencegah regresi |
| 6 | D-31 + D-32 manajemen pengguna | operasional harian |
| 7 | D-83 + D-86 rate limit & CORS allow-list | hardening |
| 8 | D-94 E2E | mengunci kriteria PRD §9 |
| 9 | D-99 dashboard (post-MVP) | hanya setelah MVP stabil |

> **Housekeeping:** item bertanda 🔧 (D-171, D-176…D-185, D-191…D-194) sudah
> teruji dan sebaiknya **dicommit lebih dulu** — setelah itu statusnya menjadi
> ✅ dan §11 tidak lagi menyebut working tree yang belum bersih.

> Keputusan K-1…K-6 di §9.1.3 **sudah diambil** dan diterapkan di kode.

---

## 13. Rencana Task — Unduh Data Rekrutmen (export XLSX)

**Konteks:** aplikasi lama (`old-app/`) sudah punya tombol **Export Excel**
(`old-app/index.html:326`, server `getRecruitmentsForExport` di
`old-app/code.gs:884`) yang bisa dipakai **user level 1 (ANGGOTA)** — datanya
dibatasi ke record miliknya sendiri. Aplikasi baru belum punya fitur ini.

**Keputusan pemilik produk (2026-11-14):**

| # | Keputusan | Alasan |
|---|---|---|
| K-7 | Format file **XLSX** | Setara dengan app lama; dibuka langsung di Excel |
| K-8 | Export **mengikuti filter aktif** (pencarian, kecamatan, kelurahan, metode, status) | File yang diunduh sama dengan yang terlihat di layar |
| K-9 | Backend mengirim **JSON** (headers + rows), frontend menyusun XLSX dengan SheetJS | Bundle tetap kecil (SheetJS dari CDN, seperti `frontend/src/excel.js`), backend tanpa dependensi baru |

### Daftar task

| # | Task | Status | Verifikasi |
|---|---|---|---|
| D-197 | Rencana task tercatat di dokumen ini | ✅ | §13 |
| D-198 | Endpoint `GET /api/rekrutmen/export` (`requireAuth`, **bukan** `requireAdmin`) | ✅ | `backend/internal/api/export.go` |
| D-199 | Scope per level: level 1 = data sendiri, level 2 = kecamatannya, level 3 = semua; tanpa paginasi | ✅ | `user.ScopeFor()` + `Limit` besar; `TestExportMemberHoldsOnlyOwnRecords` |
| D-200 | Format respons JSON: `{headers, rows, total, filename}` | ✅ | `export.go` |
| D-201 | Klien API `api.rekrutmen.export(filters)` | ✅ | `frontend/src/api.js` |
| D-202 | Penulis XLSX di frontend + fallback CSV bila SheetJS gagal dimuat | ✅ | `frontend/src/excel.js` `downloadWorkbook()` |
| D-203 | Tombol "⬇ Unduh Data" di `ListView.vue`, memakai filter aktif | ✅ | `frontend/src/views/ListView.vue` |
| D-204 | Test backend: scope export per level + apply filter + tanpa paginasi | ✅ | `backend/internal/api/export_test.go` (5 test) |
| D-205 | Dokumentasi endpoint di `API_CONTRACT.md` | ✅ | § Rekrutmen + § "Unduh data" |
| D-206 | Build, test, dan `graphify update .` | ✅ | `go test ./...` ✅, `npm run build` ✅ |

> Status item berubah menjadi ✅ setelah diverifikasi lewat kolom **Verifikasi**.

---

## 14. Rencana Task — Dashboard Rekapitulasi (D-99)

**Konteks:** aplikasi lama (`old-app/`) sudah punya halaman dashboard
rekapitulasi yang dipakai pengelola untuk memantau target pendaftaran. Item
D-99 tadinya *di luar scope MVP*; sekarang dikerjakan sebagai fase lanjutan.

**Keputusan desain:**

| # | Keputusan | Alasan |
|---|---|---|
| K-8 | Endpoint mengembalikan **agregat saja**, tidak pernah record tunggal | Dashboard tidak boleh jadi jalan masuk ke data akun lain |
| K-9 | Cakupan mengikuti **level** akun: level 3 = semua, level 2 = kecamatannya, level 1 = filter yang dipilih | Konsisten dengan list dan export |
| K-10 | Anggota (level 1) boleh memilih cakupan: `ALL` / `KECAMATAN` / `ANGGOTA` | Sama seperti aplikasi lama |
| K-11 | `kecamatanSummary` selalu memuat 7 kecamatan, walau kosong | Panel tetap rapi saat data sedikit |

### Daftar task

| # | Task | Status | Verifikasi |
|---|---|---|---|
| D-207 | Endpoint `GET /api/dashboard` (`requireAuth`) | ✅ | `backend/internal/api/dashboard.go` |
| D-208 | Agregat: total, metode, metode×status, tingkat penerimaan, status, per kecamatan, top 10 nama | ✅ | `buildDashboard()` |
| D-209 | Scope per level + filter anggota | ✅ | `dashboardRowAllowed()`, `TestDashboard…` |
| D-210 | Halaman `DashboardView.vue` + menu navigasi + route `/dashboard` | ✅ | `frontend/src/views/DashboardView.vue` |
| D-211 | Klien API `api.dashboard(query)` | ✅ | `frontend/src/api.js` |
| D-212 | Gaya kartu/panel dashboard | ✅ | `frontend/src/style.css` (`.dash-*`) |
| D-213 | Test backend: scope per level, ringkasan metode, tanpa sesi | ✅ | `backend/internal/api/dashboard_test.go` (5 test) |
| D-214 | Test E2E: panel, cakupan admin penuh, filter anggota, cakupan admin kecamatan | ✅ | `tests/e2e/dashboard.spec.ts` (4 test) |
| D-215 | Dokumentasi endpoint di `API_CONTRACT.md` + fitur di `README.md` | ✅ | § Dashboard |
| D-216 | Perbaikan upload massal: `api.rekrutmen.uploadBatch` (pemanggilan salah) | ✅ | `frontend/src/components/BatchUpload.vue` |
| D-217 | Batas upload dapat dikonfigurasi lewat `UPLOAD_MAX_ATTEMPTS` | ✅ | `backend/internal/api/batch.go`, `cmd/server/main.go` |
| D-218 | Akun uji lengkap untuk E2E (`petugas`, `admin`, `pw*`) | ✅ | `backend/cmd/e2eseed/main.go` |
| D-219 | Test backend batas upload (429 setelah batas) | ✅ | `TestBatchStopsAfterTooManyUploads` |

> Status item berubah menjadi ✅ setelah diverifikasi lewat kolom **Verifikasi**.

---

**Dokumen terkait:** [`TESTING_PLAN.md`](TESTING_PLAN.md) ·
[`DEPLOYMENT_PLAN.md`](DEPLOYMENT_PLAN.md)

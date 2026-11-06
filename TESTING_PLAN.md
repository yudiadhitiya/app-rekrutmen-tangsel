# Testing Plan — Rekrutmen Tangsel

**Status dokumen:** aktif, diperbarui per checkpoint
**Terakhir diperbarui:** 2026-11-06
**Commit baseline:** `3d7d05a`
**Rujukan:** [`PRD_Rekrutmen_Tangsel_MVP.md`](PRD_Rekrutmen_Tangsel_MVP.md) v0.1,
[`API_CONTRACT.md`](API_CONTRACT.md),
[`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md)

---

## 1. Cara memakai dokumen ini

Plan ini memakai penanda status yang sama dengan
[`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md):

| Penanda | Arti |
|---|---|
| ✅ | **Selesai** — test ada, lulus, dan sudah ter-commit |
| 🔧 | **Sedang dikerjakan** — test ditulis tapi belum lulus / belum ter-commit |
| ⬜ | **Belum** — test belum ada |

Nomor item `T-xx` dipakai di dokumen ini; `D-xx` merujuk ke item di
[`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md). Perubahan kode di `D-xx` yang
menuntut test baru akan dicatat di bagian [5. Matriks Coverage](#5-matriks-coverage).

---

## 2. Ringkasan

| Lapisan | Sudah ada & lulus | Belum ada | Total | Catatan |
|---|---|---|---|---|
| Unit — validasi | 13 | 0 | 13 | coverage 81.0% |
| Integrasi — API | 17 | 0 | 17 | coverage 63.1% |
| Integrasi — store local | 12 | 0 | 12 | coverage 86.5% |
| Integrasi — store sheets | 6 | 0 | 6 | coverage 78.7% |
| Integrasi — store supabase | 5 | 0 | 5 | coverage 51.8% |
| Unit — auth / settings / seed / model | 0 | 4 | 4 | **0% coverage, tanpa test file** |
| Unit — frontend | 0 | 2 | 2 | **tidak ada test runner** |
| E2E + responsif + a11y | 0 | 4 | 4 | belum ada |
| Integrasi provider sungguhan | 0 | 2 | 2 | mock `httptest` saja |
| Keamanan | 0 | 3 | 3 | cookie `Secure`, rate limit, batas `limit` |
| Lainnya (kontrak, ekspor/impor, migrasi, performa, beban, browser, restore) | 0 | 7 | 7 | lihat §6 |
| **Total** | **53** | **22** | **75** | |

Kondisi sekarang: **semua test yang ada lulus** (`go test ./...` → 5 paket ok,
`npm run build` → sukses). Yang belum ada: test frontend, E2E, dan unit test
untuk tiga paket backend yang benar-benar 0%.

---

## 3. Perintah verifikasi

Jalankan dari root repo.

```bash
# Backend — build, vet, test, coverage
cd backend
export PATH=$PATH:/usr/local/go/bin
go build ./...
go vet ./...
go test -count=1 ./...
go test -count=1 -cover ./...

# Test satu paket
go test -count=1 -run TestRekrutmenCRUDLifecycle -v ./internal/api/

# Frontend — saat ini hanya build (belum ada test runner)
cd ../frontend
npm ci
npm run build
```

> `go test` memakai cache. Selalu pakai `-count=1` saat memverifikasi,
> supaya hasil benar-benar mencerminkan kode terkini.

### Baseline terukur (2026-11-06)

```
ok  rekrutmen-tangsel/backend/internal/api              4.014s  coverage: 63.1%
ok  rekrutmen-tangsel/backend/internal/store/local      0.341s  coverage: 86.5%
ok  rekrutmen-tangsel/backend/internal/store/sheets     3.378s  coverage: 78.7%
ok  rekrutmen-tangsel/backend/internal/store/supabase   0.145s  coverage: 51.8%
ok  rekrutmen-tangsel/backend/internal/validate         0.050s  coverage: 81.0%
```

Paket tanpa test file: `cmd/server`, `internal/auth`, `internal/model`,
`internal/seed`, `internal/settings`, `internal/store`.

---

## 4. Test yang sudah ada

### 4.1 Validasi — `backend/internal/validate/validate_test.go` (13 test, 81.0%)

| # | Test | Mengunci aturan |
|---|---|---|
| T-1 | `TestValidRecordHasNoErrors` | record valid → tanpa error |
| T-2 | `TestRequiredFields` | semua field wajib (nama, no_wa, kecamatan, aktivitas, tingkat, metode, status) |
| T-3 | `TestWhatsAppMustBeDigitsOnly` | PRD §10 butir 2 |
| T-4 | `TestWhatsAppKeepsLeadingZero` | PRD §4 — nol di awal tidak hilang |
| T-5 | `TestRTAndRWMustBeDigitsOnly` | PRD §10 butir 3 |
| T-6 | `TestRTAndRWMayBeEmpty` | RT/RW opsional |
| T-7 | `TestKecamatanMustBeOneOfSeven` | PRD §10 butir 4 |
| T-8 | `TestKelurahanMustBelongToKecamatan` | PRD §10 butir 5 |
| T-9 | `TestKelurahanMayBeEmpty` | kelurahan opsional |
| T-10 | `TestAktivitasMustExistInMaster` | aktivitas harus di master |
| T-11 | `TestEnumValuesMustBeValid` | PRD §10 butir 6 |
| T-12 | `TestUPLAcceptsEmptyUPLSorUPLM` | PRD §10 butir 7 & 11 |
| T-13 | `TestFieldsAreTrimmed` | input di-trim |

### 4.2 API — `backend/internal/api/api_test.go` (17 test, 63.1%)

| # | Test | Mengunci |
|---|---|---|
| T-14 | `TestLoginSetsSessionAndMeReturnsUser` | alur login + `/me` |
| T-15 | `TestLoginRejectsBadPassword` | password salah → ditolak |
| T-16 | `TestProtectedRoutesRequireAuth` | tanpa sesi → 401 |
| T-17 | `TestRekrutmenCRUDLifecycle` | PRD §10 butir 1 (lengkap) |
| T-18 | `TestUpdatePreservesCreationAuditFields` | **regresi D-29** — jejak audit bertahan |
| T-19 | `TestDuplicateWhatsAppReturnsConflict` | PRD §10 butir 12 |
| T-20 | `TestCheckWhatsAppReportsExistence` | PRD §9 — peringatan sebelum simpan |
| T-21 | `TestValidationErrorsAreReportedPerField` | bentuk respons `{error, fields}` |
| T-22 | `TestWilayahEndpointReturnsSevenKecamatan` | 7 kecamatan |
| T-23 | `TestAktivitasMasterIsAdminOnly` | PRD §10 butir 14 — hak admin |
| T-24 | `TestPetugasCannotChangeMasterAktivitas` | petugas ditolak |
| T-25 | `TestPetugasCanCreateRekrutmen` | petugas boleh CRUD rekrutmen |
| T-26 | `TestSettingsRequireAdmin` | pengaturan khusus admin |
| T-27 | `TestSettingsSwitchProvider` | PRD §10 butir 13 |
| T-28 | `TestSearchAndFilter` | pencarian + filter (PRD §6) |
| T-29 | `TestAktivitasInUseCannotBeDeleted` | 409 — aktivitas terpakai |
| T-30 | `TestLogoutClearsSession` | logout membersihkan sesi |

### 4.3 Store local — `backend/internal/store/local/store_test.go` (12 test, 86.5%)

| # | Test | Mengunci |
|---|---|---|
| T-31 | `TestCreateRejectsDuplicateWhatsApp` | unikensi saat create |
| T-32 | `TestUpdateRejectsDuplicateWhatsApp` | unikensi saat update |
| T-33 | `TestUpdateWithSameWhatsAppSucceeds` | update dengan WA sendiri tetap boleh |
| T-34 | `TestListFilters` | filter Works di provider |
| T-35 | `TestListPagination` | paginasi benar |
| T-36 | `TestFindByWhatsApp` | lookup by WA |
| T-37 | `TestDeleteRekrutmen` | hapus |
| T-38 | `TestAktivitasMasterCRUD` | CRUD master aktivitas |
| T-39 | `TestAktivitasRenameCascadesToRekrutmen` | cascade rename |
| T-40 | `TestDeleteAktivitasInUseFails` | proteksi hapus |
| T-41 | `TestCreateAktivitasRejectsDuplicateName` | nama aktivitas unik |
| T-42 | `TestDataPersistsAcrossReopen` | data bertahan setelah store ditutup & dibuka ulang |

### 4.4 Store Google Sheets — `backend/internal/store/sheets/store_test.go` (6 test, 78.7%)

| # | Test | Mengunci |
|---|---|---|
| T-43 | `TestSheetsEnsureTabsCreatesHeaders` | tab & header dibuat otomatis |
| T-44 | `TestSheetsRekrutmenCRUD` | CRUD di Sheets |
| T-45 | `TestSheetsDeleteKeepsOtherRowsIntact` | hapus tidak memotong baris lain |
| T-46 | `TestSheetsListAndFilter` | list + filter di Sheets |
| T-47 | `TestSheetsFindByWhatsApp` | lookup by WA |
| T-48 | `TestSheetsAktivitasMaster` | master aktivitas di Sheets |

> Test Sheets memakai `httptest` sebagai **fake Google API** (bukan API asli).
> Lihat T-75 untuk uji integrasi sungguhan.

### 4.5 Store Supabase — `backend/internal/store/supabase/store_test.go` (5 test, 51.8%)

| # | Test | Mengunci |
|---|---|---|
| T-49 | `TestSupabaseCreateAndGet` | create + get |
| T-50 | `TestSupabaseUniqueViolationMapsToConflict` | error unique → `ErrDuplicateWA` |
| T-51 | `TestSupabaseUpdateAndDelete` | update + delete |
| T-52 | `TestSupabaseFindByWhatsApp` | lookup by WA |
| T-53 | `TestSupabaseAktivitasMaster` | master aktivitas |

---

## 5. Matriks coverage

Menjadi ✅ bila ada test yang gagal saat aturan dilanggar, dan lulus saat
dipulihkan.

| Area | Kasus diuji | Status |
|---|---|---|
| Validasi input | 13 aturan PRD | ✅ |
| Autentikasi | login benar/salah, sesi, logout, tanpa auth | ✅ |
| Otorisasi peran | admin vs petugas, endpoint per peran | ✅ |
| CRUD rekrutmen | create/read/update/delete | ✅ |
| Jejak audit | created_at/by bertahan saat update | ✅ |
| Unikensi no_wa | create, update, lookup | ✅ |
| Master wilayah | 7 kecamatan / 54 kelurahan | ✅ |
| Master aktivitas | CRUD, cascade, proteksi hapus | ✅ |
| Filter & paginasi | q, status, filter, page/limit | ✅ |
| Bentuk respons error | 400 per-field, 401, 403, 404, 409 | ✅ |
| Penyimpanan | ketiga provider, CRUD + filter | ✅ |
| Ganti provider | validasi provider, provider tak terpasang | ✅ |
| Hash password | PBKDF2, salt, timing-equalization | ⬜ |
| Session cookie | TTL, tanda tangan, kedaluwarsa | ⬜ |
| Settings store | persist, load, default, file rusak | ⬜ |
| Seed wilayah | integritas 7/54, waliRT | ⬜ |
| **Frontend — form** | validasi klien, submit, error field | ⬜ |
| **Frontend — daftar** | render, filter, paginasi, delete | ⬜ |
| **Frontend — auth** | route guard, state sesi | ⬜ |
| **E2E alur utama** | PRD §5 tujuh langkah | ⬜ |
| **E2E peran** | petugas vs admin | ⬜ |
| **Responsif** | tampilan mobile | ⬜ |
| **Aksesibilitas** | keyboard, label, kontras | ⬜ |

---

## 6. Yang belum ada

| # | Item | Status | Detail |
|---|---|---|---|
| T-54 | **Unit test `internal/auth`** | ⬜ | `HashPassword`, `VerifyPassword` (user ada/tidak, password benar/salah), `RandomSecret`, parsing `users.json` |
| T-55 | **Unit test `internal/settings`** | ⬜ | default `local`, persist, load, file JSON rusak |
| T-56 | **Unit test `internal/seed`** | ⬜ | integritas 7 kecamatan/54 kelurahan, tidak ada kelurahan tanpa kecamatan |
| T-57 | **Unit test `internal/model`** | ⬜ | serialisasi JSON, field audit |
| T-58 | **Test runner frontend** | ⬜ | pasang Vitest + `@vue/test-utils` |
| T-59 | **Test unit komponen frontend** | ⬜ | `FormView` validasi, `ListView` filter, guard router |
| T-60 | **Test E2E (Playwright)** | ⬜ | alur PRD §5 end-to-end |
| T-61 | **Test E2E peran** | ⬜ | petugas tidak bisa masuk menu admin |
| T-62 | **Test responsif** | ⬜ | viewport mobile, tabel→kartu |
| T-63 | **Audit a11y** | ⬜ | axe-core di E2E |
| T-64 | **Integration test Supabase sungguhan** | ⬜ | skip kalau `SUPABASE_URL` kosong |
| T-65 | **Integration test Sheets sungguhan** | ⬜ | skip kalau `GOOGLE_SHEET_ID` kosong |
| T-66 | **Test keamanan cookie `Secure`** | ⬜ | terkait D-82 |
| T-67 | **Test rate limit login** | ⬜ | terkait D-83 |
| T-68 | **Test batas atas `limit`** | ⬜ | terkait D-88 |
| T-69 | **Test ekspor/impor** | ⬜ | terkait D-97 |
| T-70 | **Test kontrak API (schema)** | ⬜ | pastikan respons sesuai `API_CONTRACT.md` |
| T-71 | **Test naik-turunkan versi (migration)** | ⬜ | bila skema berubah |
| T-72 | **Benchmark performa list** | ⬜ | data besar (10k+ baris) |
| T-73 | **Uji beban (load test)** | ⬜ | pra-pemcctlon produksi |
| T-74 | **Uji kompatibilitas browser** | ⬜ | Chrome/Firefox/Safari/Edge |
| T-75 | **Uji restore backup** | ⬜ | terkait P-35 |

---

## 7. Smoke test manual

Checklist yang dijalankan setelah deploy atau sebelum rilis.aty.

### 7.1 Autentikasi

| # | Langkah | Harapan |
|---|---|---|
| S-01 | Buka `/pengaturan` tanpa login | Diarahkan ke `/login` |
| S-02 | Login dengan password salah | Pesan "username atau password salah" |
| S-03 | Login sebagai `admin` | Masuk, menu admin tampil |
| S-04 | Login sebagai `petugas` | Masuk, menu admin tersembunyi |
| S-05 | Klik logout | Kembali ke `/login`, akses `/` ditolak |

### 7.2 CRUD & validasi

| # | Langkah | Harapan |
|---|---|---|
| S-06 | Tambah data lengkap | Toast sukses, data muncul di daftar |
| S-07 | Kosongkan nama | Pesan error di bawah field, tidak terkirim |
| S-08 | Isi no_wa dengan huruf | Pesan "hanya boleh berisi angka 0-9" |
| S-09 | Pilih kecamatan, lalu kelurahan | Hanya kelurahan sesuai kecamatan |
| S-10 | Ganti kecamatan | Pilihan kelurahan ter-reset |
| S-11 | Simpan no_wa yang sudah ada | Peringatan duplikat sebelum simpan |
| S-12 | Ubah data | Perubahan tersimpan, `updated_at` berubah, `created_at` tetap |
| S-13 | Hapus data | Dialog konfirmasi, data hilang |
| S-14 | Pilih UPLS lalu UPLM | Hanya satu aktif (radio) |

### 7.3 Master & pengaturan

| # | Langkah | Harapan |
|---|---|---|
| S-15 | Tambah aktivitas baru sebagai admin | Muncul di dropdown form |
| S-16 | Hapus aktivitas yang terpakai | Ditolak dengan pesan jelas |
| S-17 | Ganti provider di Pengaturan | Toast sukses, badge provider berubah |
| S-18 | Buka Pengaturan sebagai petugas | Diarahkan ke `/` |

### 7.4 Responsif & izin

| # | Langkah | Harapan |
|---|---|---|
| S-19 | Lebar layar desktop | Tabel tampil |
| S-20 | Lebar layar mobile | Kartu tampil, tabel disembunyikan |
| S-21 | DevTools → Application → Cookies | `HttpOnly` = true, `SameSite` = Lax |
| S-22 | DevTools → Network → /api/* | Semua 2xx, tidak ada 500 |

---

## 8. Strategi test per lapisan

| Lapisan | Tool | Cakupan | Kapan dijalankan |
|---|---|---|---|
| Unit (backend) | `go test` bawaan | fungsi tunggal, tanpa I/O | tiap commit |
| Integrasi (backend) | `httptest` + provider fake | handler + routing + RBAC | tiap commit |
| Kontrak API | `httptest` + golden file | respons vs `API_CONTRACT.md` | tiap commit |
| Unit (frontend) | Vitest + `@vue/test-utils` | komponen & composable | tiap commit |
| E2E | Playwright | alur pengguna nyata | nightly + pre-release |
| A11y | axe-core via Playwright | WCAG AA | pre-release |
| Manual | checklist §7 | sebab-akibat + smoke | tiap deploy |
| Load | k6 / hey | list + create | pra-release |

> Golden-file test kontrak (T-70) sangat berguna di sini: karena
> `API_CONTRACT.md` adalah sumber kebenaran yang dipakai frontend **dan**
> backend, test golden memastikan tidak ada diam-diamnya perubahan bentuk
> respons.

---

## 9. Standar penulisan test

1. **Satu test, satu aturan.** Kalau nama test berisi "dan", pecah jadi dua.
2. **Format nama:** `Test<Subjek><Perilaku>`, cermin nama fungsi yang diuji.
3. **Gunakan `newTestEnv(t)`** (`api_test.go:24`) — sudah menangani direktori
   sementara, cookie jar, dan server uji.
4. **Gunakan `validBody()`** (`api_test.go:84`) sebagai basis; ubah hanya field
   yang diuji.
5. **Assertion harus spesifik.** Hindari `if got != want` tanpa pesan; sertakan
   nilai yang diharapkan agar kegagalan langsung terbaca.
6. **Pakai `t.Fatalf` untuk kondisi yang memblokir** test lanjutan, `t.Errorf`
   untuk yang independen — supaya semua kegagalan terlihat sekaligus.
7. **Test regresi wajib** untuk setiap bugfix, dan harus **gagal tanpa fix**.
   Uji ini pada bug D-29: cabut fix → test gagal → pasang fix → test lulus.
8. **Tidak ada sleeps.** Pakai `httptest` server, bukan penundaan waktu.

---

## 10. Checkpoint terakhir

| Item | Nilai |
|---|---|
| Commit terakhir | `3d7d05a` |
| `go test ./...` | ✅ 5 paket ok |
| Coverage | api 63.1% · local 86.5% · sheets 78.7% · supabase 51.8% · validate 81.0% |
| `npm run build` | ✅ sukses |
| Total test | 53 lulus |
| Test belum ada | 22 item |
| Test runner frontend | belum ada |

### Titik kembali kalau terputus

1. `cd backend && go test -count=1 ./...` — pastikan hijau.
2. `git status` — cari test yang ditulis tapi belum di-commit (status `🔧`).
3. Lanjutkan dari `T-54` (unit test auth) — paket dengan celah coverage terbesar
   dan paling cepat ditulis.

---

## 11. Target cakupan

| Target | Cum-bu | Saat |
|---|---|---|
| Backend total | ≥ 75% | sebelum rilis |
| `internal/auth` | ≥ 80% | sebelum rilis (T-54) |
| `internal/settings` | ≥ 80% | sebelum rilis (T-55) |
| `internal/store/supabase` | ≥ 70% | sebelum rilis |
| Frontend | ≥ 60% | setelah Vitest masuk (T-58) |
| Kriteria PRD §10 | 100% punya test | sebelum rilis — **sudah 13/14 terpenuhi**, sisa butir 9 (responsif) menunggu T-62 |

---

**Dokumen terkait:** [`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md) ·
[`DEPLOYMENT_PLAN.md`](DEPLOYMENT_PLAN.md)

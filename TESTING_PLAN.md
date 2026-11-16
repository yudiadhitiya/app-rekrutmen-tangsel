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
| Unit — validasi | 15 | 0 | 15 | coverage 76.2% |
| Integrasi — API | 69 | 0 | 69 | coverage 73.4% |
| Integrasi — store local | 14 | 0 | 14 | coverage 78.8% |
| Integrasi — store sheets | 35 | 0 | 35 | coverage 69.5% |
| Integrasi — store supabase | 5 | 0 | 5 | coverage 47.4% |
| Unit — auth / model | 20 | 0 | 20 | auth 61.6% · model 53.3% |
| Unit — settings / seed | 0 | 2 | 2 | **0% coverage, tanpa test file** |
| Unit — frontend | 0 | 2 | 2 | **tidak ada test runner** |
| E2E Playwright | 20 | 2 | 22 | upload, dashboard, ubah password, browser |
| E2E responsif + a11y | 0 | 2 | 2 | belum ada |
| Integrasi provider sungguhan | 0 | 2 | 2 | mock `httptest` saja |
| Keamanan | 0 | 3 | 3 | cookie `Secure`, rate limit, batas `limit` |
| Lainnya (kontrak, ekspor/impor, migrasi, performa, beban, restore) | 0 | 6 | 6 | lihat §6 |
| **Total** | **178** | **19** | **197** | |

Kondisi sekarang: **semua test yang ada lulus** (`go test ./...` → 8 paket ok,
164 test Go; `npm run build` → sukses; `npx playwright test` → 20 test E2E
lulus). Yang belum ada: test frontend (Vitest), responsif, a11y, dan unit
test untuk dua paket backend yang benar-benar 0%.

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

### E2E (Playwright)

E2E butuh server yang sudah berjalan di `http://localhost:8100`, memakai data
uji yang terpisah dari `data/` repo.

```bash
# 1. Seed akun uji (petugas, admin, pwpetugas, pwkec, pwadmin)
cd backend
go run ./cmd/e2eseed /tmp/e2e-data
printf '{"rekrutmen":[],"aktivitas":[]}' > /tmp/e2e-data/rekrutmen.json

# 2. Jalankan server uji
#    UPLOAD_MAX_ATTEMPTS dinaikkan, sebab batas upload bawaan (5 per 10 menit)
#    menahan rangkaian test upload.
APP_ROOT="$(pwd)/.." DATA_DIR=/tmp/e2e-data \
FRONTEND_DIR="$(pwd)/../frontend/dist" ADDR=:8100 \
LEGACY_AUTH_MODE=off UPLOAD_MAX_ATTEMPTS=50 go run ./cmd/server

# 3. Jalankan test E2E dari root repo
npx playwright test
```

> `go test` memakai cache. Selalu pakai `-count=1` saat memverifikasi,
> supaya hasil benar-benar mencerminkan kode terkini.

### Baseline terukur (2026-11-06)

```
ok  rekrutmen-tangsel/backend/internal/api              5.028s  coverage: 63.1%
ok  rekrutmen-tangsel/backend/internal/store/local      0.555s  coverage: 86.9%
ok  rekrutmen-tangsel/backend/internal/store/sheets     4.620s  coverage: 78.8%
ok  rekrutmen-tangsel/backend/internal/store/supabase   0.176s  coverage: 51.4%
ok  rekrutmen-tangsel/backend/internal/validate         0.098s  coverage: 81.4%
```

Paket tanpa test file: `cmd/server`, `internal/auth`, `internal/model`,
`internal/seed`, `internal/settings`, `internal/store`.

---

## 4. Test yang sudah ada

### 4.1 Validasi — `backend/internal/validate/validate_test.go` (15 test, 81.4%)

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
| T-14 | `TestMetodeOptionsAreTheThreeCurrentValues` | enum `metode` = 3 opsi (**regresi D-104**) |
| T-15 | `TestRetiredMetodeSpellingsAreNormalized` | "Program"/"Antara" diterima & dinormalisasi |

### 4.2 API — `backend/internal/api/api_test.go` (18 test, 63.1%)

| # | Test | Mengunci |
|---|---|---|
| T-16 | `TestLoginSetsSessionAndMeReturnsUser` | alur login + `/me` |
| T-17 | `TestLoginRejectsBadPassword` | password salah → ditolak |
| T-18 | `TestProtectedRoutesRequireAuth` | tanpa sesi → 401 |
| T-19 | `TestRekrutmenCRUDLifecycle` | PRD §10 butir 1 (lengkap) |
| T-20 | `TestUpdatePreservesCreationAuditFields` | **regresi D-29** — jejak audit bertahan |
| T-21 | `TestMetodeProgramAntaraIsStoredAndFilterable` | **regresi D-104** — metode tersimpan & terfilter |
| T-22 | `TestDuplicateWhatsAppReturnsConflict` | PRD §10 butir 12 |
| T-23 | `TestCheckWhatsAppReportsExistence` | PRD §9 — peringatan sebelum simpan |
| T-24 | `TestValidationErrorsAreReportedPerField` | bentuk respons `{error, fields}` |
| T-25 | `TestWilayahEndpointReturnsSevenKecamatan` | 7 kecamatan |
| T-26 | `TestAktivitasMasterIsAdminOnly` | PRD §10 butir 14 — hak admin |
| T-27 | `TestPetugasCannotChangeMasterAktivitas` | petugas ditolak |
| T-28 | `TestPetugasCanCreateRekrutmen` | petugas boleh CRUD rekrutmen |
| T-29 | `TestSettingsRequireAdmin` | pengaturan khusus admin |
| T-30 | `TestSettingsSwitchProvider` | PRD §10 butir 13 |
| T-31 | `TestSearchAndFilter` | pencarian + filter (PRD §6) |
| T-32 | `TestAktivitasInUseCannotBeDeleted` | 409 — aktivitas terpakai |
| T-33 | `TestLogoutClearsSession` | logout membersihkan sesi |

### 4.3 Store local — `backend/internal/store/local/store_test.go` (14 test, 86.9%)

| # | Test | Mengunci |
|---|---|---|
| T-34 | `TestCreateRejectsDuplicateWhatsApp` | unikensi saat create |
| T-35 | `TestUpdateRejectsDuplicateWhatsApp` | unikensi saat update |
| T-36 | `TestUpdateWithSameWhatsAppSucceeds` | update dengan WA sendiri tetap boleh |
| T-37 | `TestListFilters` | filter Works di provider |
| T-38 | `TestListPagination` | paginasi benar |
| T-39 | `TestFindByWhatsApp` | lookup by WA |
| T-40 | `TestDeleteRekrutmen` | hapus |
| T-41 | `TestAktivitasMasterCRUD` | CRUD master aktivitas |
| T-42 | `TestAktivitasRenameCascadesToRekrutmen` | cascade rename |
| T-43 | `TestDeleteAktivitasInUseFails` | proteksi hapus |
| T-44 | `TestCreateAktivitasRejectsDuplicateName` | nama aktivitas unik |
| T-45 | `TestRetiredMetodeIsMigratedOnLoad` | **regresi D-104** — data lama dimigrasi saat store dibuka |
| T-46 | `TestFilterByProgramAntaraMatchesRetiredSpellings` | **regresi D-104** — filter cocok dengan nilai lama |
| T-47 | `TestDataPersistsAcrossReopen` | data bertahan setelah store ditutup & dibuka ulang |

### 4.4 Store Google Sheets — `backend/internal/store/sheets/store_test.go` (6 test, 78.8%)

| # | Test | Mengunci |
|---|---|---|
| T-48 | `TestSheetsEnsureTabsCreatesHeaders` | tab & header dibuat otomatis |
| T-49 | `TestSheetsRekrutmenCRUD` | CRUD di Sheets |
| T-50 | `TestSheetsDeleteKeepsOtherRowsIntact` | hapus tidak memotong baris lain |
| T-51 | `TestSheetsListAndFilter` | list + filter di Sheets |
| T-52 | `TestSheetsFindByWhatsApp` | lookup by WA |
| T-53 | `TestSheetsAktivitasMaster` | master aktivitas di Sheets |

> Test Sheets memakai `httptest` sebagai **fake Google API** (bukan API asli).
> Lihat T-80 untuk uji integrasi sungguhan.

### 4.5 Store Supabase — `backend/internal/store/supabase/store_test.go` (5 test, 51.4%)

| # | Test | Mengunci |
|---|---|---|
| T-54 | `TestSupabaseCreateAndGet` | create + get |
| T-55 | `TestSupabaseUniqueViolationMapsToConflict` | error unique → `ErrDuplicateWA` |
| T-56 | `TestSupabaseUpdateAndDelete` | update + delete |
| T-57 | `TestSupabaseFindByWhatsApp` | lookup by WA |
| T-58 | `TestSupabaseAktivitasMaster` | master aktivitas |

### 4.6 Auth & ubah password — `backend/internal/auth/auth_test.go`, `backend/internal/api/password_test.go` (18 test)

| # | Test | Mengunci |
|---|---|---|
| T-98 | `TestVerifyLoginUsesLocalAccount` | akun lokal tidak menyentuh sumber warisan |
| T-99 | `TestChangePasswordUpdatesTheLocalAccount` | hash baru tersimpan di `users.json` |
| T-100 | `TestChangePasswordRejectsWrongCurrentPassword` | password lama salah → `ErrInvalidCredentials` |
| T-101 | `TestChangePasswordInSheetOnlyMode` | akun hanya di sheet → password ditulis ke sheet |
| T-102 | `TestChangePasswordUpdatesBothStores` | akun ada di lokal & sheet → keduanya diperbarui |
| T-103 | `TestChangePasswordOfUnknownAccount` | akun tak dikenal → `ErrAccountNotFound` |
| T-104 | `TestChangePasswordNeedsSession` | tanpa sesi → 401 |
| T-105 | `TestChangePasswordOfEveryLevel` | level 1, 2, 3 boleh mengganti password sendiri |
| T-106 | `TestChangePasswordRejectsShortPassword` | password < 6 karakter → 400 + `fields.password` |
| T-107 | `TestChangePasswordRejectsSamePassword` | password baru harus berbeda |
| T-108 | `TestChangePasswordRejectsBrokenBody` | body rusak → 400 `body tidak valid` |
| T-109 | `TestChangePasswordKeepsTheSession` | sesi tidak diputus setelah ganti password |
| T-110 | `TestChangePasswordIsUsedByTheNextLogin` | password baru berlaku, yang lama ditolak |
| T-111 | `TestChangePasswordUpdatesTheAccountOfThePreviousApplication` | hash baru ikut ditulis ke store warisan |

> Rincian matriks dan hasil eksekusi: [`testing_plan.md`](testing_plan.md).

### 4.7 E2E — `tests/e2e/ubah-password.spec.ts` (7 test)

| # | Test | Mengunci |
|---|---|---|
| T-112 | `Tombol Password tersedia untuk semua level akun` | tombol header untuk level 1, 2, 3 |
| T-113 | `Dialog mengubah password menampilkan tiga field` | form lama/baru/ulang, `type=password` |
| T-114 | `Konfirmasi password yang tidak sama ditolak di sisi klien` | validasi konfirmasi |
| T-115 | `Password lama yang salah ditolak oleh server` | pesan `password lama salah` |
| T-116 | `Password baru terlalu pendek ditolak di sisi klien` | minimal 6 karakter |
| T-117 | `Dialog bisa ditutup dan isiannya dikosongkan` | Batal → form kosong |
| T-118 | `Password berhasil diganti dan dipakai untuk masuk lagi` | toast sukses, sesi bertahan, password lama ditolak |

### 4.8 Dashboard — `backend/internal/api/dashboard_test.go` (5 test)

| # | Test | Mengunci |
|---|---|---|
| T-119 | `TestDashboardCountsEveryRecordForFullAdmin` | admin penuh menghitung semua record; 7 kecamatan selalu tampil |
| T-120 | `TestDashboardDistrictAdminOnlyCountsOwnDistrict` | admin kecamatan hanya menghitung kecamatannya |
| T-121 | `TestDashboardMemberFollowsTheChosenFilter` | anggota mengikuti filter `ALL` / `KECAMATAN` / `ANGGOTA` |
| T-122 | `TestDashboardMethodSummaryCountsTheStatuses` | rekap metode menghitung berhasil, proses, gagal/pending |
| T-123 | `TestDashboardRequiresAuthentication` | tanpa sesi → 401 |

### 4.9 E2E Dashboard — `tests/e2e/dashboard.spec.ts` (4 test)

| # | Test | Mengunci |
|---|---|---|
| T-124 | `Dashboard terbuka dari navigasi dan menampilkan semua panel` | 4 kartu + 5 panel (penerimaan, status, metode, kecamatan, top 10) |
| T-125 | `Admin penuh melihat seluruh data pada dashboard` | scope "Seluruh data pendaftaran", data baru masuk ke rekap |
| T-126 | `Anggota memilih cakupan dashboard sendiri` | filter 3 pilihan; scope berubah sesuai filter |
| T-127 | `Admin kecamatan hanya melihat kecamatannya` | scope "KECAMATAN Serpong", tanpa pilihan filter |

---

### 4.10 Upload massal — `backend/internal/api/batch_test.go` (10 test)

| # | Test | Mengunci |
|---|---|---|
| T-128 | `TestMemberCanUploadBatch` | anggota boleh mengunggah; record milik pengunggah |
| T-129 | `TestBatchIgnoresTheOwnerColumnsOfTheFile` | kolom pemilik dari file lama diabaikan |
| T-130 | `TestBatchSavesNothingWhenOneRowIsWrong` | satu baris salah → tidak ada baris tersimpan |
| T-131 | `TestBatchRejectsTwoRowsWithTheSameNumber` | nomor ganda di satu file ditolak |
| T-132 | `TestBatchSkipsANumberThatIsAlreadyStored` | nomor terdaftar dilewati dan alasannya dilaporkan |
| T-133 | `TestBatchIsOnlyForMembers` | admin penuh tidak boleh mengunggah |
| T-134 | `TestBatchRejectsMoreThanTheAllowedRows` | lebih dari 1.000 baris ditolak |
| T-135 | `TestBatchStopsAfterTooManyUploads` | batas upload per akun (429 setelah batas) |
| T-136 | `TestBatchAcceptsTheSpellingsOfThePreviousApplication` | ejaan metode app lama diterima |
| T-137 | `TestBatchAcceptsTheFieldNamesOfThisApplication` | nama field aplikasi ini diterima |

### 4.11 E2E Upload massal — `tests/e2e/batch-upload.spec.ts` (8 test)

| # | Test | Mengunci |
|---|---|---|
| T-138 | `Anggota melihat kartu upload dan bisa mengunduh template` | kartu upload + unduh template |
| T-139 | `Anggota penuh (level 3) tidak melihat kartu upload` | upload hanya untuk anggota |
| T-140 | `Upload CSV yang valid menyimpan data dan menampilkan ringkasan` | ringkasan "2 data tersimpan." |
| T-141 | `Baris yang salah tidak disimpan dan nomor barisnya dilaporkan` | nomor baris dilapor |
| T-142 | `Nomor yang sama di dalam satu file menggagalkan upload` | duplikat di dalam file |
| T-143 | `Nomor yang sudah terdaftar dilewati dan dilaporkan` | baris dilewati + sebab |
| T-144 | `Berkas dengan lebih dari 1.000 baris ditolak di sisi klien` | batas baris di klien |
| T-145 | `Ejaan metode dari aplikasi lama tetap diterima` | ejaan lama → "Dakwah Fardiyah" |

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
| Upload massal | anggota saja, baris salah, duplikat, batas baris, batas percobaan | ✅ |
| Penyimpanan | ketiga provider, CRUD + filter | ✅ |
| Ganti provider | validasi provider, provider tak terpasang | ✅ |
| Hash password | PBKDF2, salt, timing-equalization | ✅ |
| Ubah password | endpoint, 3 level akun, warisan, E2E dialog | ✅ |
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
| T-59 | **Unit test `internal/auth`** | ✅ (sebagian) | `HashPassword`, `VerifyPassword`, `RandomSecret`, `users.json`, `ChangePassword` (T-98…T-103) |
| T-60 | **Unit test `internal/settings`** | ⬜ | default `local`, persist, load, file JSON rusak |
| T-61 | **Unit test `internal/seed`** | ⬜ | integritas 7 kecamatan/54 kelurahan, tidak ada kelurahan tanpa kecamatan |
| T-62 | **Unit test `internal/model`** | ⬜ | serialisasi JSON, field audit |
| T-63 | **Test runner frontend** | ⬜ | pasang Vitest + `@vue/test-utils` |
| T-64 | **Test unit komponen frontend** | ⬜ | `FormView` validasi, `ListView` filter, guard router |
| T-65 | **Test E2E (Playwright)** | ✅ (sebagian) | alur ubah password (T-112…T-118) dan dashboard (T-124…T-127); alur PRD §5 menyusul |
| T-66 | **Test E2E peran** | ⬜ | petugas tidak bisa masuk menu admin |
| T-67 | **Test responsif** | ⬜ | viewport mobile, tabel→kartu |
| T-68 | **Audit a11y** | ⬜ | axe-core di E2E |
| T-69 | **Integration test Supabase sungguhan** | ⬜ | skip kalau `SUPABASE_URL` kosong |
| T-70 | **Integration test Sheets sungguhan** | ⬜ | skip kalau `GOOGLE_SHEET_ID` kosong |
| T-81 | **Resolusi mapping kolom Sheets** | ✅ | `TestResolveColumnMapMatchesAliases`, `...IgnoresPunctuation`, `...UsesAliasWhenNeeded`, `...PrefersFirstMatch`, `...ReportsMissingRequiredColumn` |
| T-82 | **Konversi nilai baris Sheets** | ✅ | `TestColumnMapRecord`, `TestParseTimestampFormats`, `TestNormalizeNoWA` |
| T-83 | **ID turunan stabil (idempoten)** | ✅ | `TestDerivedIDIsStable` |
| T-84 | **Perencanaan impor respons form** | ✅ | `TestBuildPlanImportsValidRow`, `...SkipsDuplicateNumbers`, `...MarksInvalidRowForReview`, `...SkipsRowWithoutNumber`, `...ReportsNewActivity`, `...IsIdempotent` |
| T-85 | **Mode `form` tidak pernah membuat tab** | ✅ | `TestSheetsFormModeNeverCreatesTabs` |
| T-86 | **Konfigurasi tab dan id spreadsheet** | ✅ | `TestSheetsUsesConfiguredTabNames`, `TestSheetsAcceptsSpreadsheetURL`, `TestSheetsRejectsInvalidSpreadsheetID` |
| T-87 | **File kredensial dengan izin longgar ditolak** | ✅ | `TestSheetsRejectsReadableKeyFile` |
| T-88 | **Login memakai akun aplikasi lama** | ✅ | `TestSheetsVerifyLegacyUser`, `...IsCaseInsensitive`, `...RejectsBadCredentials`, `...MapsRoles`, `...ReportsMissingTab` |
| T-89 | **Pemetaan peran aplikasi lama** | ✅ | `TestLegacyRoleMapping` |
| T-90 | **Login lokal didahulukan atas akun lama** | ✅ | `TestVerifyLoginUsesLocalAccount` |
| T-91 | **Login lama dipindahkan ke akun lokal** | ✅ | `TestMigrateLegacyUserStoresAccountLocally`, `...KeepsExistingAccount` |
| T-92 | **Endpoint login menerima akun lama** | ✅ | `TestLoginWithLegacyAccount`, `TestLoginWithWrongLegacyPassword` |
| T-93 | **Integrasi login ke spreadsheet sungguhan** | ⬜ | `TestLegacyLoginIntegration`, skip tanpa `GOOGLE_LEGACY_TEST_USER` |
| T-94 | **Format hash akun lama di spreadsheet sungguhan** | ✅ | `TestLegacyHashFormat` — 153/153 hash berupa sha256 hex |
| T-95 | **Pendaftaran mandiri menulis baris sesuai layout aplikasi lama** | ✅ | `TestRegisterUserAppendsRowWithLegacyLayout`, `...HashesPasswordWithSHA256`, `...NormalizesSecurityAnswers`, `...RejectsTakenUsername`, `TestRegisteredUserCanSignIn` |
| T-96 | **Validasi & endpoint pendaftaran mandiri** | ✅ | `TestRegisterCreatesAnAccount`, `TestRegisterRejectsInvalidInput`, `TestRegisterReportsTakenUsername`, `TestRegisterWithoutAccountStoreIsUnavailable`, `TestRegisterOptionsArePublic` |
| T-97 | **Uji langsung pendaftaran ke spreadsheet sungguhan** | ✅ | manual: daftar → baris muncul di tab `USERS` → login berhasil → baris uji dihapus |
| T-71 | **Test keamanan cookie `Secure`** | ⬜ | terkait D-82 |
| T-72 | **Test rate limit login** | ⬜ | terkait D-83 |
| T-73 | **Test batas atas `limit`** | ⬜ | terkait D-88 |
| T-74 | **Test ekspor/impor** | ⬜ | terkait D-97 |
| T-75 | **Test kontrak API (schema)** | ⬜ | pastikan respons sesuai `API_CONTRACT.md` |
| T-76 | **Test naik-turunkan versi (migration)** | ⬜ | bila skema berubah |
| T-77 | **Benchmark performa list** | ⬜ | data besar (10k+ baris) |
| T-78 | **Uji beban (load test)** | ⬜ | pra-pemcctlon produksi |
| T-79 | **Uji kompatibilitas browser** | ⬜ | Chrome/Firefox/Safari/Edge |
| T-80 | **Uji restore backup** | ⬜ | terkait P-35 |

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

> Golden-file test kontrak (T-75) sangat berguna di sini: karena
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
| Commit terakhir | `5acf1c1` (pekerjaan dashboard & upload belum di-commit) |
| `go test ./...` | ✅ 8 paket ok (164 test Go) |
| Coverage | api 73.4% · auth 61.6% · local 78.8% · sheets 69.5% · supabase 47.4% · model 53.3% · validate 76.2% |
| `npm run build` | ✅ sukses |
| `npx playwright test` | ✅ 20 lulus (upload 8, dashboard 4, ubah password 7, browser 1) |
| Total test | 164 lulus (Go) + 20 E2E Playwright |
| Test belum ada | 19 item |
| Test runner frontend | belum ada |
| Uji fitur ubah password | ✅ selesai 2026-11-16 — lihat [`testing_plan.md`](testing_plan.md) |
| Uji fitur dashboard | ✅ selesai 2026-11-16 — T-119…T-123 (Go) & T-124…T-127 (E2E) |
| Uji fitur upload massal | ✅ selesai 2026-11-16 — T-128…T-137 (Go) & T-138…T-145 (E2E) |

### Titik kembali kalau terputus

1. `cd backend && go test -count=1 ./...` — pastikan hijau.
2. `git status` — cari test yang ditulis tapi belum di-commit (status `🔧`).
3. Lanjutkan dari `T-59` (unit test auth) — paket dengan celah coverage terbesar
   dan paling cepat ditulis.

---

## 11. Target cakupan

| Target | Cum-bu | Saat |
|---|---|---|
| Backend total | ≥ 75% | sebelum rilis |
| `internal/auth` | ≥ 80% | sebelum rilis (T-59) |
| `internal/settings` | ≥ 80% | sebelum rilis (T-60) |
| `internal/store/supabase` | ≥ 70% | sebelum rilis |
| Frontend | ≥ 60% | setelah Vitest masuk (T-63) |
| Kriteria PRD §10 | 100% punya test | sebelum rilis — **sudah 13/14 terpenuhi**, sisa butir 9 (responsif) menunggu T-67 |

---

**Dokumen terkait:** [`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md) ·
[`DEPLOYMENT_PLAN.md`](DEPLOYMENT_PLAN.md)

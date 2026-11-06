# Development Plan — Rekrutmen Tangsel

**Status dokumen:** aktif, diperbarui per checkpoint
**Terakhir diperbarui:** 2026-11-06
**Commit baseline:** `3d7d05a` (fix: pertahankan jejak audit saat update rekrutmen)
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
| F2 — API backend | 22 | 0 | 4 | 26 |
| F3 — Penyimpanan | 10 | 0 | 2 | 12 |
| F4 — Frontend | 27 | 0 | 4 | 31 |
| F5 — Keamanan | 4 | 0 | 7 | 11 |
| F6 — Kualitas & CI | 2 | 0 | 6 | 8 |
| F7 — Fitur pelengkap PRD | 0 | 0 | 7 | 7 |
| **Total** | **73** | **0** | **30** | **103** |

Status codebase: **build hijau, test hijau, siap didemokan** — belum siap
diproduksi karena F5 (keamanan) dan F6 (CI) masih kosong.

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
| D-31 | Endpoint ganti password | ⬜ | tidak ada; `auth.go` hanya punya verify/hash |
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
| D-84 | **Endpoint ganti password** | ⬜ | lihat D-31 |
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
| D-94 | **Test E2E (Playwright)** | ⬜ | belum ada |
| D-95 | **Pinning versi toolchain** | ⬜ | `go.mod` butuh `go 1.24`; Node belum di-pin |
| D-96 | **Pre-commit hook** | ⬜ | `.claude/settings.json` hanya berisi hook graphify |

---

## 9. F7 — Fitur pelengkap PRD

| # | Item | Status | Keterkaitan PRD |
|---|---|---|---|
| D-97 | **Ekspor/import + verifikasi ganti provider** | ⬜ | §7, §10 |
| D-98 | **Manajemen pengguna (CRUD petugas)** | ⬜ | §2, §9 |
| D-99 | **Dashboard & laporan analitik** | ⬜ | §11 — *di luar scope MVP*, fase berikutnya |
| D-100 | **Impor massal Excel/CSV** | ⬜ | §11 — *di luar scope MVP* |
| D-101 | **Integrasi WhatsApp kirim pesan** | ⬜ | §11 — *di luar scope MVP* |
| D-102 | **Workflow persetujuan berjenjang** | ⬜ | §11 — *di luar scope MVP* |
| D-103 | **Sinkronisasi dua arah Supabase↔Sheets** | ⬜ | §11 — *di luar scope MVP* |

> D-99 s.d. D-103 secara eksplisit **di luar scope MVP**. Item-item ini dicatat
> agar tidak hilang, bukan untuk dikerjakan sekarang.
> tidak hilang, bukan untuk dikerjakan sekarang.

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
| Commit terakhir | `3d7d05a` |
| Remote `main` | `3d7d05a` (selaras dengan lokal) |
| Working tree | bersih |
| `go build ./...` | ✅ |
| `go vet ./...` | ✅ |
| `go test ./...` | ✅ 5 paket lulus |
| `npm run build` | ✅ 3.8 s, 140 kB |
| Test regression untuk D-29 | ✅ lulus; terbukti menangkap bug bila fix dibatalkan |

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
| 3 | D-97 ekspor/import + verifikasi provider | gap PRD §7 yang eksplisit |
| 4 | D-91 CI + D-93 test frontend | mencegah regresi |
| 5 | D-31 + D-32 manajemen pengguna | operasional harian |
| 6 | D-83 + D-86 rate limit & CORS allow-list | hardening |
| 7 | D-94 E2E | mengunci kriteria PRD §9 |
| 8 | D-99 dashboard (post-MVP) | hanya setelah MVP stabil |

---

**Dokumen terkait:** [`TESTING_PLAN.md`](TESTING_PLAN.md) ·
[`DEPLOYMENT_PLAN.md`](DEPLOYMENT_PLAN.md)

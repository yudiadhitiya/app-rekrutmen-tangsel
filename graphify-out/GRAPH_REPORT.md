# Graph Report - app-rekrutmen-tangsel  (2026-11-12)

## Corpus Check
- 63 files · ~54,233 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 1, .example 1, .css 1)

## Summary
- 781 nodes · 2073 edges · 36 communities (30 shown, 6 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 299 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a17b2195`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- api.go
- impor-responses/main.go
- seed.go
- testing.T
- Server
- ListView.vue
- Deployment Plan — Rekrutmen Tangsel
- App.vue
- FormView.vue
- Store
- Rekrutmen
- local/store_test.go
- 3. Fase kerja
- Aktivitas
- sheets/store_test.go
- Store
- AktivitasView.vue
- Store
- package.json
- Aplikasi Rekrutmen Tangsel
- Development Plan — Rekrutmen Tangsel
- api.js
- SettingsView.vue
- PRD --- Aplikasi Rekrutmen Tangsel
- API Contract — Rekrutmen Tangsel
- model.go
- 6. Yang belum ada
- LoginView.vue
- Testing Plan — Rekrutmen Tangsel
- Menyiapkan Google Sheets
- CLAUDE.md
- rekrutmen-tangsel/backend
- users_test.go
- toast.js
- run-dev.sh

## God Nodes (most connected - your core abstractions)
1. `Rekrutmen` - 35 edges
2. `Server` - 34 edges
3. `Store` - 34 edges
4. `User` - 30 edges
5. `6. Yang belum ada` - 28 edges
6. `newTestEnv()` - 26 edges
7. `Rekrutmen()` - 24 edges
8. `writeJSON()` - 23 edges
9. `Manager` - 22 edges
10. `Store` - 21 edges

## Surprising Connections (you probably didn't know these)
- `9. Standar penulisan test` --references--> `validBody()`  [INFERRED]
  TESTING_PLAN.md → backend/internal/api/api_test.go
- `Fase 4 — Impor data Responses (1.5–2 hari) ✅` --references--> `Rekrutmen`  [INFERRED]
  SHEETS_INTEGRATION_PLAN.md → backend/internal/model/model.go
- `4. F2 — API backend ✅` --references--> `withCORS()`  [INFERRED]
  DEVELOPMENT_PLAN.md → backend/internal/api/api.go
- `6. Yang belum ada` --references--> `TestRegisterUserAppendsRowWithLegacyLayout()`  [INFERRED]
  TESTING_PLAN.md → backend/internal/store/sheets/register_test.go
- `6. Yang belum ada` --references--> `TestRegisteredUserCanSignIn()`  [INFERRED]
  TESTING_PLAN.md → backend/internal/store/sheets/register_test.go

## Import Cycles
- None detected.

## Communities (36 total, 6 thin omitted)

### Community 0 - "api.go"
Cohesion: 0.08
Nodes (20): errorRegistrar, fakeRegistrar, registerLimiter, isSecurityQuestion(), newRegisterLimiter(), validateRegistration(), Registration, normalizeHeader() (+12 more)

### Community 1 - "impor-responses/main.go"
Cohesion: 0.08
Nodes (44): buildPlan(), contains(), envOr(), fieldReasons(), indexOf(), inspectSpreadsheet(), main(), printSummary() (+36 more)

### Community 2 - "seed.go"
Cohesion: 0.46
Nodes (6): All(), IsKecamatan(), IsKelurahan(), KelurahanByKecamatan(), Kecamatan, Kelurahan

### Community 3 - "testing.T"
Cohesion: 0.13
Nodes (50): testEnv, newTestEnv(), TestAktivitasInUseCannotBeDeleted(), TestAktivitasMasterIsAdminOnly(), TestCheckWhatsAppReportsExistence(), TestDuplicateWhatsAppReturnsConflict(), TestLoginRejectsBadPassword(), TestLoginSetsSessionAndMeReturnsUser() (+42 more)

### Community 4 - "Server"
Cohesion: 0.12
Nodes (20): legacyAuthenticator, Server, fakeLegacy, LegacyAuthenticator, User, newID(), seedAktivitasIfEmpty(), ServeFrontend() (+12 more)

### Community 5 - "ListView.vue"
Cohesion: 0.07
Nodes (29): METODE, METODE_LAMA, normalizeMetode(), STATUS, statusTone(), TINGKAT_PENERIMAAN, tingkatTone(), UPL_OPTIONS (+21 more)

### Community 6 - "Deployment Plan — Rekrutmen Tangsel"
Cohesion: 0.07
Nodes (29): 10. P7 — Monitoring & operasional, 11. P8 — Kontainer (opsional, non-blocking), 12.1 Rilis pertama, 12.2 Rilis berikutnya, 12.3 Rollback, 12. Prosedur deploy, 13. Checklist pra-rilis, 14. Checklist pasca-rilis (+21 more)

### Community 7 - "App.vue"
Cohesion: 0.12
Nodes (21): dependencies, vue@^3.5.13, vue-router@^4.5.0, api, onUnauthorized(), admin, handleLogout(), isPublicPage (+13 more)

### Community 8 - "FormView.vue"
Cohesion: 0.12
Nodes (19): aktivitasMaster, aktivitasOptions, checkDuplicateWa(), clearErrors(), errors, form, formError, handleSubmit() (+11 more)

### Community 9 - "Store"
Cohesion: 0.16
Nodes (7): colLetter(), Store, hasTab(), tabTitles(), serviceKey, sheetInfo, Store

### Community 10 - "Rekrutmen"
Cohesion: 0.19
Nodes (7): Rekrutmen, NormalizeMetode(), RekrutmenToValues(), rekrutmenToValues(), valuesToRekrutmen(), Store, normalizeRows()

### Community 11 - "local/store_test.go"
Cohesion: 0.30
Nodes (18): New(), newStore(), rec(), TestAktivitasMasterCRUD(), TestAktivitasRenameCascadesToRekrutmen(), TestCreateAktivitasRejectsDuplicateName(), TestCreateRejectsDuplicateWhatsApp(), TestDataPersistsAcrossReopen() (+10 more)

### Community 12 - "3. Fase kerja"
Cohesion: 0.10
Nodes (19): 1. Tujuan & batasan, 2. Keputusan arsitektur (blocking — perlu konfirmasi), 3. Fase kerja, 3. Status per fase, 4. Pemetaan kolom (draft — wajib dikonfirmasi di D-107), 5. Alur rollout, 6. Risiko, 7. Definition of Done (+11 more)

### Community 13 - "Aktivitas"
Cohesion: 0.19
Nodes (7): Aktivitas, AktivitasToValues(), aktivitasToValues(), valuesToAktivitas(), isUniqueViolation(), sheetEntry, sheetRow

### Community 14 - "sheets/store_test.go"
Cohesion: 0.09
Nodes (31): newFakeSheets(), newTestStore(), newTestStoreWith(), rowFromSpec(), splitPath(), testPrivateKeyPEM(), TestSheetsAcceptsSpreadsheetURL(), TestSheetsAktivitasMaster() (+23 more)

### Community 16 - "AktivitasView.vue"
Cohesion: 0.13
Nodes (16): confirmDialog(), confirmState, resolveConfirm(), adding, addItem(), admin, editForm, editingId (+8 more)

### Community 18 - "package.json"
Cohesion: 0.16
Nodes (13): devDependencies, vite@^6.0.7, @vitejs/plugin-vue@^5.2.1, name, private, scripts, build, dev (+5 more)

### Community 19 - "Aplikasi Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): Akun awal, Aplikasi Rekrutmen Tangsel, Backend, Catatan, Dokumen Proyek, Fitur, Frontend, Menjalankan (+4 more)

### Community 20 - "Development Plan — Rekrutmen Tangsel"
Cohesion: 0.18
Nodes (11): 11. Checkpoint terakhir, 12. Urutan pengerjaan berikutnya, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. F1 — Fondasi & model data ✅, 6. F4 — Frontend ✅, 7. F5 — Keamanan 🔐, 8. F6 — Kualitas & CI (+3 more)

### Community 21 - "api.js"
Cohesion: 0.43
Nodes (5): ApiError, buildUrl(), defaultMessage(), request(), unauthorizedHandlers

### Community 22 - "SettingsView.vue"
Cohesion: 0.20
Nodes (11): providerLabel(), PROVIDERS, currentProvider, errorMessage, loading, loadSettings(), provider, PROVIDER_DESCRIPTIONS (+3 more)

### Community 23 - "PRD --- Aplikasi Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): 10. Kriteria Penerimaan MVP, 11. Di Luar Scope MVP, 1. Tujuan, 2. Pengguna, 3. Fitur MVP, 4. Struktur Data dan Aturan Input, 5. Alur Utama, 6. UI/UX (+4 more)

### Community 24 - "API Contract — Rekrutmen Tangsel"
Cohesion: 0.22
Nodes (9): API Contract — Rekrutmen Tangsel, Auth, Error shape, Frontend routes (vue-router), Master aktivitas (job/activity categories), Master wilayah (read-only, all roles), Rekrutmen, Settings (admin only) (+1 more)

### Community 25 - "model.go"
Cohesion: 0.52
Nodes (3): Filter, MetodeMatches(), MetodeMatchValues()

### Community 26 - "6. Yang belum ada"
Cohesion: 0.08
Nodes (33): Config, env(), main(), main(), New(), TestLoginWithLegacyAccount(), TestLoginWithWrongLegacyPassword(), newRegisterEnv() (+25 more)

### Community 27 - "LoginView.vue"
Cohesion: 0.10
Nodes (22): login(), answer1, answer2, errorMessage, handleRegister(), handleSubmit(), kecamatan, kecamatanOptions (+14 more)

### Community 29 - "Testing Plan — Rekrutmen Tangsel"
Cohesion: 0.12
Nodes (16): 10. Checkpoint terakhir, 11. Target cakupan, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. Perintah verifikasi, 5. Matriks coverage, 7.1 Autentikasi, 7.2 CRUD & validasi (+8 more)

### Community 30 - "Menyiapkan Google Sheets"
Cohesion: 0.40
Nodes (5): Daftar anggota mandiri (self-registration), Lokasi file kredensial, Masuk dengan akun dari aplikasi lama, Mengimpor data dari Google Form, Menyiapkan Google Sheets

### Community 33 - "users_test.go"
Cohesion: 0.31
Nodes (11): LegacyRole(), LegacyUserID(), sha256Hex(), legacyUsersRows(), seedUsersTab(), TestLegacyRoleMapping(), TestSheetsVerifyLegacyUser(), TestSheetsVerifyLegacyUserIsCaseInsensitive() (+3 more)

### Community 34 - "toast.js"
Cohesion: 0.43
Nodes (5): TITLES, dismissToast(), showToast(), toast, toasts

## Knowledge Gaps
- **2 isolated node(s):** `PATH`, `rekrutmen-tangsel/backend`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 210 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Development Plan — Rekrutmen Tangsel` connect `Development Plan — Rekrutmen Tangsel` to `6. Yang belum ada`, `testing.T`, `DEPLOYMENT_PLAN.md`, `sheets/store_test.go`?**
  _High betweenness centrality (0.405) - this node is a cross-community bridge._
- **What connects `PATH`, `rekrutmen-tangsel/backend` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `api.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08315789473684211 - nodes in this community are weakly interconnected._
- **Why does `6. F4 — Frontend ✅` connect `Development Plan — Rekrutmen Tangsel` to `FormView.vue`, `api.js`?**
  _High betweenness centrality (0.335) - this node is a cross-community bridge._
- **Should `impor-responses/main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07597402597402597 - nodes in this community are weakly interconnected._
- **Why does `ApiError` connect `api.js` to `FormView.vue`, `Development Plan — Rekrutmen Tangsel`?**
  _High betweenness centrality (0.208) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.13283208020050125 - nodes in this community are weakly interconnected._
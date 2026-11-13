# Graph Report - app-rekrutmen-tangsel  (2026-11-13)

## Corpus Check
- 101 files · ~108,309 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 6 file(s) not represented in the graph (top: .csv 2, (none) 1, .example 1)

## Summary
- 1453 nodes · 3475 edges · 73 communities (57 shown, 16 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 413 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9ea813d4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- User
- canAccess
- api_test.go
- Server
- ListView.vue
- Deployment Plan — Rekrutmen Tangsel
- App.vue
- FormView.vue
- Store
- buildPlan
- Manager
- 3. Fase kerja
- sheets/store_test.go
- supabase/store_test.go
- Rekrutmen
- AktivitasView.vue
- Store
- package.json
- Aplikasi Rekrutmen Tangsel
- validate_data.py
- api.js
- SettingsView.vue
- PRD --- Aplikasi Rekrutmen Tangsel
- API Contract — Rekrutmen Tangsel
- Rekrutmen
- testing.T
- LoginView.vue
- NormalizeMetode
- local/store_test.go
- CLAUDE.md
- rekrutmen-tangsel/backend
- Registration
- Development Plan — Rekrutmen Tangsel
- run-dev.sh
- search
- 6. Yang belum ada
- TestMetricMath
- Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)
- design_system.py
- BM25
- mapping.go
- Menyiapkan Google Sheets
- design_system.py
- sheets/register_test.go
- UI/UX Pro Max - Design Intelligence
- CatalogRefreshTest
- search_stack
- Pre-Delivery Checklist (canonical — the only one)
- Quick Reference
- impor-responses/main.go
- DesignSystemGenerator
- Testing Plan — Rekrutmen Tangsel
- Store
- RtlCliTest
- TestThresholdGate
- TestFixtureValidation
- toast.js
- .handleGetSettings
- users_test.go
- json
- seed.go
- TestStackFlagWithDesignSystem
- Store
- 4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi
- _filter_anti_patterns_for_mode
- 12. Prosedur deploy
- rekrutmen-tangsel/backend/internal/model.Aktivitas
- Aktivitas
- core.py
- test_design_system_mode.py
- test_catalog_summary_line_endings.py

## God Nodes (most connected - your core abstractions)
1. `search()` - 54 edges
2. `Store` - 48 edges
3. `User` - 45 edges
4. `DesignSystemGenerator` - 38 edges
5. `search_stack()` - 36 edges
6. `Server` - 34 edges
7. `Rekrutmen` - 28 edges
8. `6. Yang belum ada` - 28 edges
9. `newTestEnv()` - 25 edges
10. `Manager` - 24 edges

## Surprising Connections (you probably didn't know these)
- `1. Ringkasan Eksekutif` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `(b) Header lama tidak cocok dengan kunci yang dibaca — **data kosong buta**` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `9. Standar penulisan test` --references--> `validBody()`  [INFERRED]
  TESTING_PLAN.md → backend/internal/api/api_test.go
- `5.5 Alias importer diperluas` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go
- `(d) Importer tidak mengenal header lama` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go

## Import Cycles
- None detected.

## Communities (73 total, 16 thin omitted)

### Community 1 - "User"
Cohesion: 0.21
Nodes (6): legacyAuthenticator, fakeLegacy, fakeLookupLegacy, fakeProfiler, LegacyProfiler, User

### Community 2 - "canAccess"
Cohesion: 0.10
Nodes (23): 1. Ringkasan Eksekutif, 2.1 Model role lama — `old-app/code.gs:6-11`, 2.2 Aturan akses data — `old-app/code.gs:373-392` (dari `getRecruitmentsPaged`), 2.3 Sumber data lama, 2.4 Perilaku hak akses lain di UI lama (`old-app/index.html`), 2. Logika yang ada di `old-app` (yang seharusnya dipertahankan), 3. Kondisi aplikasi baru (verified), 5.1 — Granularity role (+15 more)

### Community 3 - "api_test.go"
Cohesion: 0.21
Nodes (26): testEnv, newTestEnv(), TestAktivitasInUseCannotBeDeleted(), TestAktivitasMasterIsAdminOnly(), TestCheckWhatsAppReportsExistence(), TestDuplicateWhatsAppReturnsConflict(), TestLoginRejectsBadPassword(), TestLoginSetsSessionAndMeReturnsUser() (+18 more)

### Community 4 - "Server"
Cohesion: 0.28
Nodes (8): Server, ServeFrontend(), withCORS(), writeError(), writeFieldError(), writeJSON(), registerAddress(), IsValidProvider()

### Community 5 - "ListView.vue"
Cohesion: 0.06
Nodes (33): METODE, METODE_LAMA, normalizeMetode(), PROVIDERS, STATUS, statusTone(), TINGKAT_PENERIMAAN, tingkatTone() (+25 more)

### Community 6 - "Deployment Plan — Rekrutmen Tangsel"
Cohesion: 0.08
Nodes (25): 10. P7 — Monitoring & operasional, 11. P8 — Kontainer (opsional, non-blocking), 13. Checklist pra-rilis, 14. Checklist pasca-rilis, 15. Environment per stadium, 16. Yang harus dikerjakan berikutnya, 17. Checkpoint terakhir, 1. Cara memakai dokumen ini (+17 more)

### Community 7 - "App.vue"
Cohesion: 0.15
Nodes (17): 5.3 Jaga hak hapus seperti lama ✅, admin, handleLogout(), isPublicPage, navLinks, route, router, userRoleLabel (+9 more)

### Community 8 - "FormView.vue"
Cohesion: 0.12
Nodes (19): aktivitasMaster, aktivitasOptions, checkDuplicateWa(), clearErrors(), errors, form, formError, handleSubmit() (+11 more)

### Community 9 - "Store"
Cohesion: 0.14
Nodes (7): colLetter(), hasTab(), isMissingTab(), tabTitles(), serviceKey, sheetInfo, Store

### Community 10 - "buildPlan"
Cohesion: 0.18
Nodes (15): buildPlan(), sortedKeys(), TestBuildPlanImportsValidRow(), TestBuildPlanIsIdempotent(), TestBuildPlanMarksInvalidRowForReview(), TestBuildPlanReportsNewActivity(), TestBuildPlanSkipsDuplicateNumbers(), TestBuildPlanSkipsRowWithoutNumber() (+7 more)

### Community 11 - "Manager"
Cohesion: 0.13
Nodes (8): LegacyAuthenticator, Manager, Registrar, isLegacyID(), legacyKnown(), loadUsers(), saveUsers(), 9. F7 — Fitur pelengkap PRD

### Community 12 - "3. Fase kerja"
Cohesion: 0.10
Nodes (19): 1. Tujuan & batasan, 2. Keputusan arsitektur (blocking — perlu konfirmasi), 3. Fase kerja, 3. Status per fase, 4. Pemetaan kolom (draft — wajib dikonfirmasi di D-107), 5. Alur rollout, 6. Risiko, 7. Definition of Done (+11 more)

### Community 13 - "sheets/store_test.go"
Cohesion: 0.13
Nodes (19): newFakeSheets(), newTestStore(), newTestStoreWith(), rowFromSpec(), splitPath(), testPrivateKeyPEM(), TestSheetsAcceptsSpreadsheetURL(), TestSheetsAktivitasMaster() (+11 more)

### Community 14 - "supabase/store_test.go"
Cohesion: 0.24
Nodes (11): matchRow(), newFakePostgREST(), newTestStore(), seedRow(), TestSupabaseAktivitasMaster(), TestSupabaseCreateAndGet(), TestSupabaseFindByWhatsApp(), TestSupabaseUniqueViolationMapsToConflict() (+3 more)

### Community 15 - "Rekrutmen"
Cohesion: 0.14
Nodes (9): testEnv, Rekrutmen, MetodeMatches(), MetodeMatchValues(), RekrutmenToValues(), rekrutmenToValues(), valuesToRekrutmen(), normalizeRows() (+1 more)

### Community 16 - "AktivitasView.vue"
Cohesion: 0.13
Nodes (16): confirmDialog(), confirmState, resolveConfirm(), adding, addItem(), admin, editForm, editingId (+8 more)

### Community 18 - "package.json"
Cohesion: 0.13
Nodes (18): dependencies, vue@^3.5.13, vue-router@^4.5.0, devDependencies, vite@^6.0.7, @vitejs/plugin-vue@^5.2.1, name, private (+10 more)

### Community 19 - "Aplikasi Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): Akun awal, Aplikasi Rekrutmen Tangsel, Backend, Catatan, Dokumen Proyek, Fitur, Frontend, Menjalankan (+4 more)

### Community 20 - "validate_data.py"
Cohesion: 0.06
Nodes (46): read_rows(), TestAccessibilityGuidance, TestChartsTypographyAndIcons, TestCurrentReactGuidance, TestSemanticColors, TestGeneratedCatalogContract, resolve(), shipped_invocations() (+38 more)

### Community 21 - "api.js"
Cohesion: 0.23
Nodes (8): ApiError, buildUrl(), defaultMessage(), onUnauthorized(), request(), unauthorizedHandlers, app, router

### Community 22 - "SettingsView.vue"
Cohesion: 0.20
Nodes (11): api, providerLabel(), currentProvider, errorMessage, loading, loadSettings(), provider, PROVIDER_DESCRIPTIONS (+3 more)

### Community 23 - "PRD --- Aplikasi Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): 10. Kriteria Penerimaan MVP, 11. Di Luar Scope MVP, 1. Tujuan, 2. Pengguna, 3. Fitur MVP, 4. Struktur Data dan Aturan Input, 5. Alur Utama, 6. UI/UX (+4 more)

### Community 24 - "API Contract — Rekrutmen Tangsel"
Cohesion: 0.22
Nodes (9): API Contract — Rekrutmen Tangsel, Auth, Error shape, Frontend routes (vue-router), Master aktivitas (job/activity categories), Master wilayah (read-only, all roles), Rekrutmen, Settings (admin only) (+1 more)

### Community 25 - "Rekrutmen"
Cohesion: 0.35
Nodes (19): Rekrutmen(), TestAktivitasMustExistInMaster(), TestEnumValuesMustBeValid(), TestFieldsAreTrimmed(), TestKecamatanMustBeOneOfSeven(), TestKelurahanMayBeEmpty(), TestKelurahanMustBelongToKecamatan(), TestMetodeOptionsAreTheThreeCurrentValues() (+11 more)

### Community 26 - "testing.T"
Cohesion: 0.34
Nodes (13): testManager(), TestMigrateLegacyUserKeepsExistingAccount(), TestMigrateLegacyUserStoresAccountLocally(), TestSheetOnlyFallsBackToAccountUnknownToTheStore(), TestSheetOnlyRejectsWrongPasswordOfKnownAccount(), TestSheetOnlyUsesTheAccountStore(), TestVerifyLoginCompletesLevelOfMigratedAccount(), TestVerifyLoginFallsBackToLegacy() (+5 more)

### Community 27 - "LoginView.vue"
Cohesion: 0.10
Nodes (22): login(), answer1, answer2, errorMessage, handleRegister(), handleSubmit(), kecamatan, kecamatanOptions (+14 more)

### Community 30 - "local/store_test.go"
Cohesion: 0.30
Nodes (18): New(), newStore(), rec(), TestAktivitasMasterCRUD(), TestAktivitasRenameCascadesToRekrutmen(), TestCreateAktivitasRejectsDuplicateName(), TestCreateRejectsDuplicateWhatsApp(), TestDataPersistsAcrossReopen() (+10 more)

### Community 33 - "Registration"
Cohesion: 0.14
Nodes (10): errorRegistrar, fakeRegistrar, registerLimiter, isSecurityQuestion(), newRegisterLimiter(), validateRegistration(), Registration, normalizeAnswer() (+2 more)

### Community 34 - "Development Plan — Rekrutmen Tangsel"
Cohesion: 0.18
Nodes (11): 11. Checkpoint terakhir, 12. Urutan pengerjaan berikutnya, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. F1 — Fondasi & model data ✅, 6. F4 — Frontend ✅, 7. F5 — Keamanan 🔐, 8. F6 — Kualitas & CI (+3 more)

### Community 36 - "search"
Cohesion: 0.07
Nodes (7): search(), TestSearchDomains, TestImeGuidance, read_rows(), RtlDataContractTest, RtlFilterTest, TestStyleTaxonomy

### Community 37 - "6. Yang belum ada"
Cohesion: 0.09
Nodes (37): Config, User, env(), main(), main(), New(), newTestEnvWithUsers(), newTestUser() (+29 more)

### Community 39 - "Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)"
Cohesion: 0.06
Nodes (32): 1. Ringkasan Eksekutif, 2. Skor Audit Teknis, 3. Skor Heuristik Nielsen, 4.1 App Shell & Navigasi (`App.vue`), 4.2 Daftar Rekrutmen (`ListView.vue`), 4.3 Form Tambah/Ubah (`FormView.vue`), 4.4 Login & Registrasi (`LoginView.vue`), 4.5 Master Aktivitas (`AktivitasView.vue`) (+24 more)

### Community 40 - "design_system.py"
Cohesion: 0.12
Nodes (3): read_rows(), read_rows(), TestTextLayoutDataContracts

### Community 41 - "BM25"
Cohesion: 0.06
Nodes (13): BM25, detect_domain(), _domain_keywords(), _file_signature(), _get_bm25(), _load_csv(), _load_csv_snapshot(), _load_product_keywords() (+5 more)

### Community 42 - "mapping.go"
Cohesion: 0.20
Nodes (18): AppColumnSpec(), FormColumnSpec(), normalizeHeader(), ResolveColumnMap(), TestAppColumnSpecResolvesLegacyHeaders(), TestColumnMapRecord(), TestRekrutmenRowValuesFollowsTheTabLayout(), TestRekrutmenRowValuesKeepsUnknownColumns() (+10 more)

### Community 43 - "Menyiapkan Google Sheets"
Cohesion: 0.40
Nodes (5): Daftar anggota mandiri (self-registration), Lokasi file kredensial, Masuk dengan akun dari aplikasi lama, Mengimpor data dari Google Form, Menyiapkan Google Sheets

### Community 44 - "design_system.py"
Cohesion: 0.06
Nodes (18): ansi_ljust(), _detect_page_type(), format_ascii_box(), add_wrapped(), wrap_text(), format_markdown(), format_page_override_md(), generate_design_system() (+10 more)

### Community 45 - "sheets/register_test.go"
Cohesion: 0.50
Nodes (7): registration(), TestRegisteredUserCanSignIn(), TestRegisterUserAppendsRowWithLegacyLayout(), TestRegisterUserHashesPasswordWithSHA256(), TestRegisterUserNormalizesSecurityAnswers(), TestRegisterUserRejectsTakenUsername(), usersGrid()

### Community 46 - "UI/UX Pro Max - Design Intelligence"
Cohesion: 0.11
Nodes (17): Before Delivering App UI, Example Workflow, If a search returns 0 results, Output Formats, Query Contract, Rule Categories by Priority, Running the search tool, Step 1: Analyze User Requirements (+9 more)

### Community 48 - "search_stack"
Cohesion: 0.08
Nodes (7): search_stack(), TestDiagnosticsContracts, _rows(), TestNativeDesktopStackFreshness, TestTextLayoutRetrieval, _rows(), TestWebStackFreshness

### Community 49 - "Pre-Delivery Checklist (canonical — the only one)"
Cohesion: 0.15
Nodes (12): Accessibility, Common Rules for Professional UI + Pre-Delivery Checklist, Icons & Visual Elements, Interaction, Interaction (App), Layout, Layout & Spacing, Light/Dark Mode (+4 more)

### Community 50 - "Quick Reference"
Cohesion: 0.15
Nodes (12): 10. Charts & Data (LOW), 1. Accessibility (CRITICAL), 2. Touch & Interaction (CRITICAL), 3. Performance (HIGH), 4. Style Selection (HIGH), 5. Layout & Responsive (HIGH), 6. Typography & Color (MEDIUM), 7. Animation (MEDIUM) (+4 more)

### Community 51 - "impor-responses/main.go"
Cohesion: 0.22
Nodes (15): contains(), envOr(), indexOf(), inspectSpreadsheet(), main(), printSummary(), run(), splitList() (+7 more)

### Community 52 - "DesignSystemGenerator"
Cohesion: 0.06
Nodes (10): DesignSystemGenerator, _resolve_dial(), TestReasoningMatch, read_rows(), split_values(), style_identities(), TestLandingAndStackContract, TestReasoningContract (+2 more)

### Community 53 - "Testing Plan — Rekrutmen Tangsel"
Cohesion: 0.12
Nodes (16): 10. Checkpoint terakhir, 11. Target cakupan, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. Perintah verifikasi, 5. Matriks coverage, 7.1 Autentikasi, 7.2 CRUD & validasi (+8 more)

### Community 58 - "toast.js"
Cohesion: 0.43
Nodes (5): TITLES, dismissToast(), showToast(), toast, toasts

### Community 59 - ".handleGetSettings"
Cohesion: 0.28
Nodes (4): sheetsMode(), sheetsTab(), Store, Settings

### Community 60 - "users_test.go"
Cohesion: 0.21
Nodes (15): LegacyLookup, LegacyLevel(), LegacyRole(), LegacyUserID(), sha256Hex(), legacyUsersRows(), seedUsersTab(), TestLegacyLevelMapping() (+7 more)

### Community 61 - "json"
Cohesion: 0.17
Nodes (5): apply_decision_rules(), _object_without_duplicates(), parse_decision_rules(), _validate_action(), _check_reasoning_contract()

### Community 62 - "seed.go"
Cohesion: 0.46
Nodes (6): All(), IsKecamatan(), IsKelurahan(), KelurahanByKecamatan(), Kecamatan, Kelurahan

### Community 64 - "Store"
Cohesion: 0.67
Nodes (3): newID(), seedAktivitasIfEmpty(), Store

### Community 65 - "4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi"
Cohesion: 0.29
Nodes (7): 4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi, (a) Tab yang dibaca ≠ tab yang berisi data historis — **penyebab utama**, (b) Header lama tidak cocok dengan kunci yang dibaca — **data kosong buta**, (c) Kolom `UPL` terpecah di data lama, Cara verifikasi cepat (2 menit, tanpa kode), (d) Importer tidak mengenal header lama, (e) Mode `form` tidak menyuap data sumber ke daftar

### Community 67 - "12. Prosedur deploy"
Cohesion: 0.50
Nodes (4): 12.1 Rilis pertama, 12.2 Rilis berikutnya, 12.3 Rollback, 12. Prosedur deploy

### Community 68 - "rekrutmen-tangsel/backend/internal/model.Aktivitas"
Cohesion: 0.21
Nodes (3): AktivitasToValues(), aktivitasToValues(), valuesToAktivitas()

### Community 70 - "Aktivitas"
Cohesion: 0.33
Nodes (4): fieldReasons(), Aktivitas(), contains(), Fields

### Community 76 - "core.py"
Cohesion: 0.07
Nodes (21): _contains_phrase(), _exact_match_diagnostic(), _exact_row_identity(), _exact_stack_identifier(), _legacy_successor_guidance(), _load_rows_or_empty(), _normalize(), _normalize_rtl_level() (+13 more)

### Community 77 - "test_design_system_mode.py"
Cohesion: 0.08
Nodes (13): _button_outline_text_color(), _contrast_ratio(), _derive_dark_palette(), format_master_md(), _palette_is_dark(), _query_wants_dark(), _relative_luminance(), _resolve_color_mode() (+5 more)

## Knowledge Gaps
- **2 isolated node(s):** `rekrutmen-tangsel/backend`, `PATH`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 422 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `5.3 — Hak hapus & master data berbasis level` connect `canAccess` to `6. Yang belum ada`, `App.vue`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `User` (e.g. with `5.1 — Granularity role` and `5.1 Kembalikan granularity role ✅`) actually correct?**
  _`User` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `rekrutmen-tangsel/backend`, `PATH` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `canAccess` be split into smaller, more focused modules?**
  _Cohesion score 0.09852216748768473 - nodes in this community are weakly interconnected._
- **Why does `isFullAdmin()` connect `App.vue` to `AktivitasView.vue`, `canAccess`, `ListView.vue`?**
  _High betweenness centrality (0.060) - this node is a cross-community bridge._
- **Are the 4 inferred relationships involving `DesignSystemGenerator` (e.g. with `TestReasoningMatch` and `TestReasoningContract`) actually correct?**
  _`DesignSystemGenerator` has 4 INFERRED edges - model-reasoned connections that need verification._
- **Should `ListView.vue` be split into smaller, more focused modules?**
  _Cohesion score 0.06072874493927125 - nodes in this community are weakly interconnected._
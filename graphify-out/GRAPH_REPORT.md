# Graph Report - app-rekrutmen-tangsel  (2026-11-17)

## Corpus Check
- 121 files · ~131,491 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 6 file(s) not represented in the graph (top: .csv 2, (none) 1, .example 1)

## Summary
- 1717 nodes · 4176 edges · 88 communities (68 shown, 20 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 578 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `545af112`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- model.go
- BatchUpload.vue
- 5.2 — Scope di API
- newTestEnv
- User
- ListView.vue
- Deployment Plan — Rekrutmen Tangsel
- App.vue
- FormView.vue
- Store
- New
- batch-upload.spec.ts
- 3. Fase kerja
- testing.T
- api/register.go
- Rekrutmen
- AktivitasView.vue
- Store
- package.json
- Aplikasi Rekrutmen Tangsel
- validate_data.py
- DashboardView.vue
- SettingsView.vue
- PRD --- Aplikasi Rekrutmen Tangsel
- API Contract — Rekrutmen Tangsel
- mapping_test.go
- read_rows
- LoginView.vue
- DesignSystemGenerator
- dashboard.go
- CLAUDE.md
- rekrutmen-tangsel/backend
- sheets/store_test.go
- local/store_test.go
- run-dev.sh
- search
- password_test.go
- model.go
- Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)
- design_system.py
- BM25
- impor-responses/main.go
- Menyiapkan Google Sheets
- format_ascii_box
- sheets/store.go
- UI/UX Pro Max - Design Intelligence
- CatalogRefreshTest
- search_stack
- Pre-Delivery Checklist (canonical — the only one)
- Quick Reference
- exportResponse
- _select_palette_for_mode
- validate.go
- Store
- TestStackFlagWithDesignSystem
- test_relevance_evaluator.py
- PasswordDialog.vue
- batch_test.go
- test_core_data_quality.py
- batch.go
- _style_is_dark_primary
- Rekrutmen
- 4. Test yang sudah ada
- core.py
- users_test.go
- _palette_is_dark
- Testing Plan — Rekrutmen Tangsel
- rekrutmen-tangsel/backend/internal/model.Aktivitas
- testEnv
- 6. Yang belum ada
- parse_decision_rules
- _filter_anti_patterns_for_mode
- buildPlan
- split_values
- RtlCliTest
- test_catalog_summary_line_endings.py
- TestGeneratedCatalogContract
- api.js
- 7. Smoke test manual
- Testing Plan — Fitur Ubah Password (Change Password)
- test_skill_script_paths.py
- test_text_layout_resilience.py

## God Nodes (most connected - your core abstractions)
1. `User` - 56 edges
2. `search()` - 54 edges
3. `Store` - 51 edges
4. `newTestEnv()` - 50 edges
5. `Rekrutmen` - 44 edges
6. `Server` - 43 edges
7. `DesignSystemGenerator` - 38 edges
8. `search_stack()` - 36 edges
9. `Manager` - 28 edges
10. `6. Yang belum ada` - 28 edges

## Surprising Connections (you probably didn't know these)
- `1. Ringkasan Eksekutif` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `(b) Header lama tidak cocok dengan kunci yang dibaca — **data kosong buta**` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `9.1.5 Urutan pengerjaan —Kemajuan` --references--> `AppColumnSpec()`  [INFERRED]
  DEVELOPMENT_PLAN.md → backend/internal/store/sheets/mapping.go
- `5.5 Alias importer diperluas` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go
- `(d) Importer tidak mengenal header lama` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go

## Import Cycles
- None detected.

## Communities (88 total, 20 thin omitted)

### Community 0 - "model.go"
Cohesion: 0.16
Nodes (6): env(), envInt(), main(), sheetsMode(), RandomSecret(), IsValidProvider()

### Community 1 - "BatchUpload.vue"
Cohesion: 0.07
Nodes (38): 13. Rencana Task — Unduh Data Rekrutmen (export XLSX), Daftar task, emit, errorMessage, errorTotal, fileInput, fileName, getTemplate() (+30 more)

### Community 2 - "5.2 — Scope di API"
Cohesion: 0.52
Nodes (3): Filter, MetodeMatches(), MetodeMatchValues()

### Community 3 - "newTestEnv"
Cohesion: 0.26
Nodes (25): newTestEnv(), TestAktivitasInUseCannotBeDeleted(), TestAktivitasMasterIsAdminOnly(), TestCheckWhatsAppReportsExistence(), TestDuplicateWhatsAppReturnsConflict(), TestLoginRejectsBadPassword(), TestLoginSetsSessionAndMeReturnsUser(), TestLogoutClearsSession() (+17 more)

### Community 4 - "User"
Cohesion: 0.06
Nodes (45): 5.1 — Granularity role, 5.1 Kembalikan granularity role ✅, 5.2 — Scope di API, 5.2 Terapkan scope di API, bukan di store ✅, 5.3 — Hak hapus & master data berbasis level, 5.3 Jaga hak hapus seperti lama ✅, 5.4 Migrasi data historis sekali jalan, 5.5 Alias importer diperluas (+37 more)

### Community 5 - "ListView.vue"
Cohesion: 0.07
Nodes (27): canUpload(), activeFilters, canDeleteRecord, deleteRecord(), deletingId, emptyMessage, emptyTitle, errorMessage (+19 more)

### Community 6 - "Deployment Plan — Rekrutmen Tangsel"
Cohesion: 0.07
Nodes (29): 10. P7 — Monitoring & operasional, 11. P8 — Kontainer (opsional, non-blocking), 12.1 Rilis pertama, 12.2 Rilis berikutnya, 12.3 Rollback, 12. Prosedur deploy, 13. Checklist pra-rilis, 14. Checklist pasca-rilis (+21 more)

### Community 7 - "App.vue"
Cohesion: 0.13
Nodes (17): onUnauthorized(), admin, handleLogout(), isPublicPage, navLinks, passwordOpen, route, router (+9 more)

### Community 8 - "FormView.vue"
Cohesion: 0.08
Nodes (32): 6. F4 — Frontend ✅, METODE, METODE_LAMA, normalizeMetode(), STATUS, statusTone(), TINGKAT_PENERIMAAN, tingkatTone() (+24 more)

### Community 9 - "Store"
Cohesion: 0.12
Nodes (7): colLetter(), hasTab(), isMissingTab(), tabTitles(), serviceKey, sheetInfo, Store

### Community 10 - "New"
Cohesion: 0.12
Nodes (22): main(), testEnv, newTestEnvWithUploadLimit(), newTestEnvWithUsers(), newTestUser(), TestLoginWithLegacyAccount(), TestLoginWithWrongLegacyPassword(), TestDashboardDistrictAdminOnlyCountsOwnDistrict() (+14 more)

### Community 11 - "batch-upload.spec.ts"
Cohesion: 0.09
Nodes (6): devDependencies, @playwright/test@^1.65.1, @playwright/test, RUN, RUN, ACCOUNTS

### Community 12 - "3. Fase kerja"
Cohesion: 0.10
Nodes (19): 1. Tujuan & batasan, 2. Keputusan arsitektur (blocking — perlu konfirmasi), 3. Fase kerja, 3. Status per fase, 4. Pemetaan kolom (draft — wajib dikonfirmasi di D-107), 5. Alur rollout, 6. Risiko, 7. Definition of Done (+11 more)

### Community 13 - "testing.T"
Cohesion: 0.19
Nodes (24): TestExportFollowsActiveFilter(), TestExportIsNotPaginated(), TestExportMemberHoldsOnlyOwnRecords(), TestExportRequiresAuth(), TestChangePasswordInSheetOnlyMode(), TestChangePasswordOfUnknownAccount(), TestChangePasswordRejectsWrongCurrentPassword(), TestChangePasswordUpdatesBothStores() (+16 more)

### Community 14 - "api/register.go"
Cohesion: 0.24
Nodes (6): errorRegistrar, fakeRegistrar, Registration, normalizeAnswer(), indexOfHeader(), LegacyUser

### Community 15 - "Rekrutmen"
Cohesion: 0.15
Nodes (9): batchRow, Rekrutmen, NormalizeMetode(), RekrutmenToValues(), rekrutmenToValues(), valuesToRekrutmen(), isUniqueViolation(), normalizeRows() (+1 more)

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
Cohesion: 0.10
Nodes (33): _catalog_date(), _check_app_interface_contract(), _check_catalog_contract(), _check_catalog_summary(), _check_color_contract(), _check_core_data_contract(), _check_file(), _check_font_catalog() (+25 more)

### Community 21 - "DashboardView.vue"
Cohesion: 0.13
Nodes (19): isMemberOnly(), cards, data, errorMessage, filter, filters, formatNumber(), kecamatanTotals (+11 more)

### Community 22 - "SettingsView.vue"
Cohesion: 0.20
Nodes (11): providerLabel(), PROVIDERS, currentProvider, errorMessage, loading, loadSettings(), provider, PROVIDER_DESCRIPTIONS (+3 more)

### Community 23 - "PRD --- Aplikasi Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): 10. Kriteria Penerimaan MVP, 11. Di Luar Scope MVP, 1. Tujuan, 2. Pengguna, 3. Fitur MVP, 4. Struktur Data dan Aturan Input, 5. Alur Utama, 6. UI/UX (+4 more)

### Community 24 - "API Contract — Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): API Contract — Rekrutmen Tangsel, Auth, Dashboard, Error shape, Frontend routes (vue-router), Master aktivitas (job/activity categories), Master wilayah (read-only, all roles), Rekrutmen (+4 more)

### Community 25 - "mapping_test.go"
Cohesion: 0.30
Nodes (14): AppColumnSpec(), FormColumnSpec(), normalizeHeader(), ResolveColumnMap(), TestAppColumnSpecResolvesLegacyHeaders(), TestColumnMapRecord(), TestRekrutmenRowValuesFollowsTheTabLayout(), TestRekrutmenRowValuesKeepsUnknownColumns() (+6 more)

### Community 26 - "read_rows"
Cohesion: 0.10
Nodes (4): _resolve_dial(), read_rows(), TestLandingAndStackContract, TestReasoningContract

### Community 27 - "LoginView.vue"
Cohesion: 0.08
Nodes (28): login(), answer1, answer2, errorMessage, errors, focusField(), handleRegister(), handleSubmit() (+20 more)

### Community 30 - "dashboard.go"
Cohesion: 0.12
Nodes (21): dashboardItem, dashboardKecamatanSummary, dashboardMethodSummary, dashboardResponse, dashboardTopItem, buildDashboard(), containsName(), dashboardBreakdown() (+13 more)

### Community 33 - "sheets/store_test.go"
Cohesion: 0.10
Nodes (28): newFakeSheets(), newTestStore(), newTestStoreWith(), rowFromSpec(), splitPath(), testPrivateKeyPEM(), TestSheetsAcceptsSpreadsheetURL(), TestSheetsAktivitasMaster() (+20 more)

### Community 34 - "local/store_test.go"
Cohesion: 0.09
Nodes (39): testEnv, New(), newStore(), rec(), TestAktivitasMasterCRUD(), TestAktivitasRenameCascadesToRekrutmen(), TestCreateAktivitasRejectsDuplicateName(), TestCreateRejectsDuplicateWhatsApp() (+31 more)

### Community 36 - "search"
Cohesion: 0.08
Nodes (6): search(), TestDiagnosticsContracts, TestSearchDomains, TestImeGuidance, RtlFilterTest, TestStyleTaxonomy

### Community 37 - "password_test.go"
Cohesion: 0.20
Nodes (16): fakeLegacyStore, TestChangePasswordIsUsedByTheNextLogin(), TestChangePasswordKeepsTheSession(), TestChangePasswordNeedsCurrentPassword(), TestChangePasswordNeedsSession(), TestChangePasswordOfEveryLevel(), TestChangePasswordOfMemberAccount(), TestChangePasswordRejectsBrokenBody() (+8 more)

### Community 39 - "Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)"
Cohesion: 0.06
Nodes (32): 1. Ringkasan Eksekutif, 2. Skor Audit Teknis, 3. Skor Heuristik Nielsen, 4.1 App Shell & Navigasi (`App.vue`), 4.2 Daftar Rekrutmen (`ListView.vue`), 4.3 Form Tambah/Ubah (`FormView.vue`), 4.4 Login & Registrasi (`LoginView.vue`), 4.5 Master Aktivitas (`AktivitasView.vue`) (+24 more)

### Community 40 - "design_system.py"
Cohesion: 0.06
Nodes (14): _button_outline_text_color(), _detect_page_type(), format_master_md(), format_page_override_md(), generate_design_system(), _generate_intelligent_overrides(), persist_design_system(), safe_slug() (+6 more)

### Community 41 - "BM25"
Cohesion: 0.06
Nodes (17): BM25, _contains_phrase(), detect_domain(), _domain_keywords(), _file_signature(), _get_bm25(), _load_csv(), _load_csv_snapshot() (+9 more)

### Community 42 - "impor-responses/main.go"
Cohesion: 0.19
Nodes (16): contains(), envOr(), indexOf(), inspectSpreadsheet(), main(), printSummary(), run(), splitList() (+8 more)

### Community 43 - "Menyiapkan Google Sheets"
Cohesion: 0.40
Nodes (5): Daftar anggota mandiri (self-registration), Lokasi file kredensial, Masuk dengan akun dari aplikasi lama, Mengimpor data dari Google Form, Menyiapkan Google Sheets

### Community 44 - "format_ascii_box"
Cohesion: 0.13
Nodes (9): ansi_ljust(), format_ascii_box(), add_wrapped(), wrap_text(), format_markdown(), hex_to_ansi(), _mode_support_labels(), section_header() (+1 more)

### Community 45 - "sheets/store.go"
Cohesion: 0.16
Nodes (14): main(), altUPLColumn(), checkKeyFilePerm(), New(), normalizeSpreadsheetID(), rekrutmenRowValues(), TestSheetsRejectsInvalidSpreadsheetID(), TestSheetsRejectsReadableKeyFile() (+6 more)

### Community 46 - "UI/UX Pro Max - Design Intelligence"
Cohesion: 0.11
Nodes (17): Before Delivering App UI, Example Workflow, If a search returns 0 results, Output Formats, Query Contract, Rule Categories by Priority, Running the search tool, Step 1: Analyze User Requirements (+9 more)

### Community 48 - "search_stack"
Cohesion: 0.10
Nodes (5): search_stack(), _rows(), TestNativeDesktopStackFreshness, _rows(), TestWebStackFreshness

### Community 49 - "Pre-Delivery Checklist (canonical — the only one)"
Cohesion: 0.15
Nodes (12): Accessibility, Common Rules for Professional UI + Pre-Delivery Checklist, Icons & Visual Elements, Interaction, Interaction (App), Layout, Layout & Spacing, Light/Dark Mode (+4 more)

### Community 50 - "Quick Reference"
Cohesion: 0.15
Nodes (12): 10. Charts & Data (LOW), 1. Accessibility (CRITICAL), 2. Touch & Interaction (CRITICAL), 3. Performance (HIGH), 4. Style Selection (HIGH), 5. Layout & Responsive (HIGH), 6. Typography & Color (MEDIUM), 7. Animation (MEDIUM) (+4 more)

### Community 52 - "_select_palette_for_mode"
Cohesion: 0.22
Nodes (4): _contrast_ratio(), _derive_dark_palette(), _select_palette_for_mode(), TestPaletteSelection

### Community 53 - "validate.go"
Cohesion: 0.40
Nodes (3): fieldReasons(), Fields, Password()

### Community 54 - "Store"
Cohesion: 0.10
Nodes (3): Store, fileState, Store

### Community 56 - "test_relevance_evaluator.py"
Cohesion: 0.12
Nodes (3): TestFixtureValidation, TestMetricMath, TestThresholdGate

### Community 57 - "PasswordDialog.vue"
Cohesion: 0.15
Nodes (14): api, confirmPassword, currentPassword, emit, errorMessage, handleSubmit(), newPassword, props (+6 more)

### Community 58 - "batch_test.go"
Cohesion: 0.40
Nodes (13): TestBatchAcceptsTheFieldNamesOfThisApplication(), TestBatchAcceptsTheSpellingsOfThePreviousApplication(), TestBatchIgnoresTheOwnerColumnsOfTheFile(), TestBatchIsOnlyForMembers(), TestBatchRejectsMoreThanTheAllowedRows(), TestBatchRejectsTwoRowsWithTheSameNumber(), TestBatchSavesNothingWhenOneRowIsWrong(), TestBatchSkipsANumberThatIsAlreadyStored() (+5 more)

### Community 59 - "test_core_data_quality.py"
Cohesion: 0.16
Nodes (10): read_rows(), TestAccessibilityGuidance, TestChartsTypographyAndIcons, TestCurrentReactGuidance, TestSemanticColors, _check_chart_contract(), _check_typography_contract(), _configured_font_names() (+2 more)

### Community 60 - "batch.go"
Cohesion: 0.17
Nodes (10): batchRowError, batchSkipped, registerLimiter, uploadLimiter, batchFields(), batchValue(), newUploadLimiter(), normalizeBatchHeader() (+2 more)

### Community 61 - "_style_is_dark_primary"
Cohesion: 0.24
Nodes (4): _query_wants_dark(), _resolve_color_mode(), _style_is_dark_primary(), TestModeResolution

### Community 62 - "Rekrutmen"
Cohesion: 0.32
Nodes (20): contains(), Rekrutmen(), TestAktivitasMustExistInMaster(), TestEnumValuesMustBeValid(), TestFieldsAreTrimmed(), TestKecamatanMustBeOneOfSeven(), TestKelurahanMayBeEmpty(), TestKelurahanMustBelongToKecamatan() (+12 more)

### Community 63 - "4. Test yang sudah ada"
Cohesion: 0.22
Nodes (9): dashboardTotal(), TestDashboardCountsEveryRecordForFullAdmin(), TestDashboardMethodSummaryCountsTheStatuses(), TestDashboardRequiresAuthentication(), 4.11 E2E Upload massal — `tests/e2e/batch-upload.spec.ts` (8 test), 4.7 E2E — `tests/e2e/ubah-password.spec.ts` (7 test), 4.8 Dashboard — `backend/internal/api/dashboard_test.go` (5 test), 4.9 E2E Dashboard — `tests/e2e/dashboard.spec.ts` (4 test) (+1 more)

### Community 64 - "core.py"
Cohesion: 0.08
Nodes (17): _exact_match_diagnostic(), _exact_row_identity(), _exact_stack_identifier(), _legacy_successor_guidance(), _normalize(), _normalize_rtl_level(), _passes_threshold(), _project_row() (+9 more)

### Community 65 - "users_test.go"
Cohesion: 0.07
Nodes (34): 1. Ringkasan Eksekutif, 2.1 Model role lama — `old-app/code.gs:6-11`, 2.2 Aturan akses data — `old-app/code.gs:373-392` (dari `getRecruitmentsPaged`), 2.3 Sumber data lama, 2.4 Perilaku hak akses lain di UI lama (`old-app/index.html`), 2. Logika yang ada di `old-app` (yang seharusnya dipertahankan), 3. Kondisi aplikasi baru (verified), 4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi (+26 more)

### Community 66 - "_palette_is_dark"
Cohesion: 0.18
Nodes (4): _palette_is_dark(), _relative_luminance(), TestEndToEndCoherence, TestLuminance

### Community 67 - "Testing Plan — Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): 10. Checkpoint terakhir, 11. Target cakupan, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. Perintah verifikasi, 5. Matriks coverage, 8. Strategi test per lapisan, 9. Standar penulisan test (+4 more)

### Community 68 - "rekrutmen-tangsel/backend/internal/model.Aktivitas"
Cohesion: 0.25
Nodes (3): Aktivitas, aktivitasToValues(), valuesToAktivitas()

### Community 70 - "6. Yang belum ada"
Cohesion: 0.25
Nodes (15): newRegisterEnv(), TestRegisterCreatesAnAccount(), TestRegisterOptionsArePublic(), TestRegisterRejectsInvalidInput(), TestRegisterReportsTakenUsername(), TestRegisterWithoutAccountStoreIsUnavailable(), validRegistration(), registration() (+7 more)

### Community 71 - "parse_decision_rules"
Cohesion: 0.31
Nodes (4): apply_decision_rules(), _object_without_duplicates(), parse_decision_rules(), _validate_action()

### Community 76 - "buildPlan"
Cohesion: 0.13
Nodes (18): buildPlan(), sortedKeys(), TestBuildPlanImportsValidRow(), TestBuildPlanIsIdempotent(), TestBuildPlanMarksInvalidRowForReview(), TestBuildPlanReportsNewActivity(), TestBuildPlanSkipsDuplicateNumbers(), TestBuildPlanSkipsRowWithoutNumber() (+10 more)

### Community 77 - "split_values"
Cohesion: 0.47
Nodes (3): split_values(), style_identities(), TestStyleIdentityContract

### Community 81 - "api.js"
Cohesion: 0.43
Nodes (5): ApiError, buildUrl(), defaultMessage(), request(), unauthorizedHandlers

### Community 82 - "7. Smoke test manual"
Cohesion: 0.40
Nodes (5): 7.1 Autentikasi, 7.2 CRUD & validasi, 7.3 Master & pengaturan, 7.4 Responsif & izin, 7. Smoke test manual

### Community 83 - "Testing Plan — Fitur Ubah Password (Change Password)"
Cohesion: 0.25
Nodes (8): 1. Ruang lingkup, 2. Kriteria lulus (exit criteria), 3. Prasyarat, 5. Hasil uji manual yang sudah dijalankan, 6. Temuan, 7. Perintah verifikasi, 8. Urutan pengerjaan test yang disarankan, Testing Plan — Fitur Ubah Password (Change Password)

### Community 84 - "test_skill_script_paths.py"
Cohesion: 0.38
Nodes (3): resolve(), shipped_invocations(), SkillScriptPathsTest

### Community 89 - "test_text_layout_resilience.py"
Cohesion: 0.20
Nodes (3): read_rows(), TestTextLayoutDataContracts, TestTextLayoutRetrieval

## Knowledge Gaps
- **3 isolated node(s):** `testEnv`, `rekrutmen-tangsel/backend`, `PATH`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 517 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **20 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `9.1.4 Item pekerjaan` connect `batch_test.go` to `local/store_test.go`, `ListView.vue`, `Rekrutmen`, `mapping_test.go`, `batch.go`?**
  _High betweenness centrality (0.045) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `User` (e.g. with `5.1 — Granularity role` and `5.1 Kembalikan granularity role ✅`) actually correct?**
  _`User` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `testEnv`, `rekrutmen-tangsel/backend`, `PATH` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `BatchUpload.vue` be split into smaller, more focused modules?**
  _Cohesion score 0.06829268292682927 - nodes in this community are weakly interconnected._
- **Why does `5.3 — Hak hapus & master data berbasis level` connect `User` to `New`, `newTestEnv`, `DashboardView.vue`, `App.vue`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **Are the 28 inferred relationships involving `newTestEnv()` (e.g. with `TestBatchAcceptsTheFieldNamesOfThisApplication()` and `TestBatchAcceptsTheSpellingsOfThePreviousApplication()`) actually correct?**
  _`newTestEnv()` has 28 INFERRED edges - model-reasoned connections that need verification._
- **Should `User` be split into smaller, more focused modules?**
  _Cohesion score 0.06084070796460177 - nodes in this community are weakly interconnected._
# Graph Report - app-rekrutmen-tangsel  (2026-11-18)

## Corpus Check
- 131 files · ~138,721 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 3, .csv 2, .service 2)

## Summary
- 1872 nodes · 4525 edges · 86 communities (67 shown, 19 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 622 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `334b8e63`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- model.go
- BatchUpload.vue
- 5.2 — Scope di API
- newTestEnv
- Server
- ListView.vue
- Deployment Plan — Rekrutmen Tangsel
- App.vue
- FormView.vue
- Store
- New
- batch-upload.spec.ts
- 3. Fase kerja
- testing.T
- steps
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
- environment
- LoginView.vue
- DesignSystemGenerator
- dashboard.go
- CLAUDE.md
- rekrutmen-tangsel/backend
- sheets/store.go
- local/store_test.go
- run-dev.sh
- search
- password_test.go
- sheets/register_test.go
- Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)
- design_system.py
- test_core.py
- impor-responses/main.go
- Menyiapkan Google Sheets
- format_ascii_box
- User
- UI/UX Pro Max - Design Intelligence
- CatalogRefreshTest
- search_stack
- Pre-Delivery Checklist (canonical — the only one)
- Quick Reference
- exportResponse
- Rekrutmen
- Manager
- TestStackFlagWithDesignSystem
- csv
- PasswordDialog.vue
- batch_test.go
- backup/main.go
- batch.go
- test_design_system_mode.py
- time.Time
- newTestEnvWithUsers
- core.py
- users_test.go
- detect_domain
- Testing Plan — Rekrutmen Tangsel
- rekrutmen-tangsel/backend/internal/model.Aktivitas
- .handleChangePassword
- 6. Yang belum ada
- .clientAddress
- supabase/store_test.go
- CatalogSummaryLineEndingsTest
- RtlCliTest
- TestFixtureValidation
- api.js
- 7. Smoke test manual
- Testing Plan — Fitur Ubah Password (Change Password)
- test_text_layout_resilience.py

## God Nodes (most connected - your core abstractions)
1. `User` - 56 edges
2. `search()` - 54 edges
3. `Store` - 51 edges
4. `newTestEnv()` - 51 edges
5. `Server` - 46 edges
6. `Rekrutmen` - 44 edges
7. `DesignSystemGenerator` - 38 edges
8. `search_stack()` - 36 edges
9. `Manager` - 33 edges
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

## Communities (86 total, 19 thin omitted)

### Community 0 - "model.go"
Cohesion: 0.13
Nodes (8): env(), envBool(), envInt(), main(), splitList(), isSecurityQuestion(), validateRegistration(), RandomSecret()

### Community 1 - "BatchUpload.vue"
Cohesion: 0.07
Nodes (37): emit, errorMessage, errorTotal, fileInput, fileName, getTemplate(), onFileSelected(), rowCount (+29 more)

### Community 2 - "5.2 — Scope di API"
Cohesion: 0.27
Nodes (8): 5.2 — Scope di API, Filter, MatchScope(), MetodeMatches(), MetodeMatchValues(), TestMatchScope(), TestScopeFor(), Scope

### Community 3 - "newTestEnv"
Cohesion: 0.26
Nodes (25): newTestEnv(), TestAktivitasInUseCannotBeDeleted(), TestAktivitasMasterIsAdminOnly(), TestCheckWhatsAppReportsExistence(), TestDuplicateWhatsAppReturnsConflict(), TestLoginRejectsBadPassword(), TestLoginSetsSessionAndMeReturnsUser(), TestLogoutClearsSession() (+17 more)

### Community 4 - "Server"
Cohesion: 0.20
Nodes (9): Server, newID(), ServeFrontend(), sheetsMode(), sheetsTab(), writeError(), writeJSON(), exportFilename() (+1 more)

### Community 5 - "ListView.vue"
Cohesion: 0.07
Nodes (27): canUpload(), activeFilters, canDeleteRecord, deleteRecord(), deletingId, emptyMessage, emptyTitle, errorMessage (+19 more)

### Community 6 - "Deployment Plan — Rekrutmen Tangsel"
Cohesion: 0.06
Nodes (37): IsSecureRequest(), normalizeOrigins(), originAllowed(), SecurityHeaders(), TestSecurityHeaders(), TestSecurityHeadersHSTSOnSecureRequest(), withCORS(), 10. P7 — Monitoring & operasional (+29 more)

### Community 7 - "App.vue"
Cohesion: 0.13
Nodes (17): onUnauthorized(), admin, handleLogout(), isPublicPage, navLinks, passwordOpen, route, router (+9 more)

### Community 8 - "FormView.vue"
Cohesion: 0.08
Nodes (32): 6. F4 — Frontend ✅, METODE, METODE_LAMA, normalizeMetode(), STATUS, statusTone(), TINGKAT_PENERIMAAN, tingkatTone() (+24 more)

### Community 9 - "Store"
Cohesion: 0.11
Nodes (9): colLetter(), hasTab(), isMissingTab(), tabTitles(), indexOfHeader(), LegacyUser, serviceKey, sheetInfo (+1 more)

### Community 10 - "New"
Cohesion: 0.12
Nodes (23): Config, testEnv, New(), newTestEnvWithUploadLimit(), TestLoginWithLegacyAccount(), TestLoginWithWrongLegacyPassword(), newLoginTestEnv(), TestLoginLockedAfterTooManyFailures() (+15 more)

### Community 11 - "batch-upload.spec.ts"
Cohesion: 0.09
Nodes (6): devDependencies, @playwright/test@^1.65.1, @playwright/test, RUN, RUN, ACCOUNTS

### Community 12 - "3. Fase kerja"
Cohesion: 0.10
Nodes (19): 1. Tujuan & batasan, 2. Keputusan arsitektur (blocking — perlu konfirmasi), 3. Fase kerja, 3. Status per fase, 4. Pemetaan kolom (draft — wajib dikonfirmasi di D-107), 5. Alur rollout, 6. Risiko, 7. Definition of Done (+11 more)

### Community 13 - "testing.T"
Cohesion: 0.16
Nodes (26): TestExportFollowsActiveFilter(), TestExportIsNotPaginated(), TestExportMemberHoldsOnlyOwnRecords(), TestExportRequiresAuth(), mustParse(), TestPasswordChangeInvalidatesOldSession(), TestChangePasswordInSheetOnlyMode(), TestChangePasswordOfUnknownAccount() (+18 more)

### Community 14 - "steps"
Cohesion: 0.05
Nodes (45): name, runs-on, steps, concurrency, cancel-in-progress, group, name, needs (+37 more)

### Community 15 - "Rekrutmen"
Cohesion: 0.14
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
Cohesion: 0.06
Nodes (43): read_rows(), TestAccessibilityGuidance, TestChartsTypographyAndIcons, TestCurrentReactGuidance, TestSemanticColors, TestGeneratedCatalogContract, _catalog_date(), _check_app_interface_contract() (+35 more)

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
Cohesion: 0.14
Nodes (22): AppColumnSpec(), FormColumnSpec(), normalizeHeader(), NormalizeNoWA(), ParseTimestamp(), ResolveColumnMap(), TestAppColumnSpecResolvesLegacyHeaders(), TestColumnMapRecord() (+14 more)

### Community 26 - "environment"
Cohesion: 0.05
Nodes (36): build, environment, healthcheck, image, ports, restart, volumes, command (+28 more)

### Community 27 - "LoginView.vue"
Cohesion: 0.08
Nodes (28): login(), answer1, answer2, errorMessage, errors, focusField(), handleRegister(), handleSubmit() (+20 more)

### Community 29 - "DesignSystemGenerator"
Cohesion: 0.05
Nodes (13): DesignSystemGenerator, _filter_anti_patterns_for_mode(), _resolve_dial(), TestReasoningMatch, read_rows(), split_values(), style_identities(), TestLandingAndStackContract (+5 more)

### Community 30 - "dashboard.go"
Cohesion: 0.13
Nodes (19): dashboardItem, dashboardKecamatanSummary, dashboardMethodSummary, dashboardResponse, dashboardTopItem, buildDashboard(), containsName(), dashboardBreakdown() (+11 more)

### Community 33 - "sheets/store.go"
Cohesion: 0.11
Nodes (16): main(), checkKeyFilePerm(), New(), normalizeSpreadsheetID(), newFakeSheets(), newTestStoreWith(), rowFromSpec(), splitPath() (+8 more)

### Community 34 - "local/store_test.go"
Cohesion: 0.07
Nodes (51): testEnv, New(), newStore(), rec(), TestAktivitasMasterCRUD(), TestAktivitasRenameCascadesToRekrutmen(), TestCreateAktivitasRejectsDuplicateName(), TestCreateRejectsDuplicateWhatsApp() (+43 more)

### Community 36 - "search"
Cohesion: 0.08
Nodes (7): _normalize_rtl_level(), _rtl_row_filter(), search(), _style_search_destination(), TestSearchDomains, RtlFilterTest, TestStyleTaxonomy

### Community 37 - "password_test.go"
Cohesion: 0.16
Nodes (20): TestChangePasswordIsUsedByTheNextLogin(), TestChangePasswordKeepsTheSession(), TestChangePasswordNeedsCurrentPassword(), TestChangePasswordNeedsSession(), TestChangePasswordOfEveryLevel(), TestChangePasswordOfMemberAccount(), TestChangePasswordRejectsBrokenBody(), TestChangePasswordRejectsSamePassword() (+12 more)

### Community 38 - "sheets/register_test.go"
Cohesion: 0.50
Nodes (7): registration(), TestRegisteredUserCanSignIn(), TestRegisterUserAppendsRowWithLegacyLayout(), TestRegisterUserHashesPasswordWithSHA256(), TestRegisterUserNormalizesSecurityAnswers(), TestRegisterUserRejectsTakenUsername(), usersGrid()

### Community 39 - "Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)"
Cohesion: 0.06
Nodes (32): 1. Ringkasan Eksekutif, 2. Skor Audit Teknis, 3. Skor Heuristik Nielsen, 4.1 App Shell & Navigasi (`App.vue`), 4.2 Daftar Rekrutmen (`ListView.vue`), 4.3 Form Tambah/Ubah (`FormView.vue`), 4.4 Login & Registrasi (`LoginView.vue`), 4.5 Master Aktivitas (`AktivitasView.vue`) (+24 more)

### Community 40 - "design_system.py"
Cohesion: 0.08
Nodes (13): _detect_page_type(), format_page_override_md(), generate_design_system(), _generate_intelligent_overrides(), persist_design_system(), safe_slug(), _write_persisted_file(), apply_decision_rules() (+5 more)

### Community 41 - "test_core.py"
Cohesion: 0.12
Nodes (5): BM25, _suggest_identities(), TestBm25CoreBehavior, TestDiagnosticsContracts, TestTokenizer

### Community 42 - "impor-responses/main.go"
Cohesion: 0.15
Nodes (27): buildPlan(), contains(), envOr(), fieldReasons(), indexOf(), inspectSpreadsheet(), main(), printSummary() (+19 more)

### Community 43 - "Menyiapkan Google Sheets"
Cohesion: 0.40
Nodes (5): Daftar anggota mandiri (self-registration), Lokasi file kredensial, Masuk dengan akun dari aplikasi lama, Mengimpor data dari Google Form, Menyiapkan Google Sheets

### Community 44 - "format_ascii_box"
Cohesion: 0.13
Nodes (9): ansi_ljust(), format_ascii_box(), add_wrapped(), wrap_text(), format_markdown(), hex_to_ansi(), _mode_support_labels(), section_header() (+1 more)

### Community 45 - "User"
Cohesion: 0.09
Nodes (23): 5.1 — Granularity role, 5.1 Kembalikan granularity role ✅, 5.2 Terapkan scope di API, bukan di store ✅, 5.3 — Hak hapus & master data berbasis level, 5.3 Jaga hak hapus seperti lama ✅, 5.4 Migrasi data historis sekali jalan, 5.5 Alias importer diperluas, 5. Rekomendasi pemulihan logika lama (+15 more)

### Community 46 - "UI/UX Pro Max - Design Intelligence"
Cohesion: 0.11
Nodes (17): Before Delivering App UI, Example Workflow, If a search returns 0 results, Output Formats, Query Contract, Rule Categories by Priority, Running the search tool, Step 1: Analyze User Requirements (+9 more)

### Community 48 - "search_stack"
Cohesion: 0.09
Nodes (8): _exact_stack_identifier(), _project_row(), search_stack(), _valid_max_results(), _rows(), TestNativeDesktopStackFreshness, _rows(), TestWebStackFreshness

### Community 49 - "Pre-Delivery Checklist (canonical — the only one)"
Cohesion: 0.15
Nodes (12): Accessibility, Common Rules for Professional UI + Pre-Delivery Checklist, Icons & Visual Elements, Interaction, Interaction (App), Layout, Layout & Spacing, Light/Dark Mode (+4 more)

### Community 50 - "Quick Reference"
Cohesion: 0.15
Nodes (12): 10. Charts & Data (LOW), 1. Accessibility (CRITICAL), 2. Touch & Interaction (CRITICAL), 3. Performance (HIGH), 4. Style Selection (HIGH), 5. Layout & Responsive (HIGH), 6. Typography & Color (MEDIUM), 7. Animation (MEDIUM) (+4 more)

### Community 52 - "Rekrutmen"
Cohesion: 0.32
Nodes (20): contains(), Rekrutmen(), TestAktivitasMustExistInMaster(), TestEnumValuesMustBeValid(), TestFieldsAreTrimmed(), TestKecamatanMustBeOneOfSeven(), TestKelurahanMayBeEmpty(), TestKelurahanMustBelongToKecamatan() (+12 more)

### Community 53 - "Manager"
Cohesion: 0.08
Nodes (13): errorRegistrar, fakeRegistrar, LegacyAuthenticator, Registrar, Manager, isLegacyID(), legacyKnown(), ParseUsers() (+5 more)

### Community 56 - "csv"
Cohesion: 0.06
Nodes (9): TestImeGuidance, TestMetricMath, TestThresholdGate, read_rows(), RtlDataContractTest, resolve(), shipped_invocations(), SkillScriptPathsTest (+1 more)

### Community 57 - "PasswordDialog.vue"
Cohesion: 0.15
Nodes (14): api, confirmPassword, currentPassword, emit, errorMessage, handleSubmit(), newPassword, props (+6 more)

### Community 58 - "batch_test.go"
Cohesion: 0.40
Nodes (13): TestBatchAcceptsTheFieldNamesOfThisApplication(), TestBatchAcceptsTheSpellingsOfThePreviousApplication(), TestBatchIgnoresTheOwnerColumnsOfTheFile(), TestBatchIsOnlyForMembers(), TestBatchRejectsMoreThanTheAllowedRows(), TestBatchRejectsTwoRowsWithTheSameNumber(), TestBatchSavesNothingWhenOneRowIsWrong(), TestBatchSkipsANumberThatIsAlreadyStored() (+5 more)

### Community 59 - "backup/main.go"
Cohesion: 0.20
Nodes (15): fileSHA256(), listFiles(), main(), prune(), readArchive(), runBackup(), runList(), runRestore() (+7 more)

### Community 60 - "batch.go"
Cohesion: 0.27
Nodes (7): batchRowError, batchSkipped, uploadLimiter, batchFields(), batchValue(), newUploadLimiter(), normalizeBatchHeader()

### Community 61 - "test_design_system_mode.py"
Cohesion: 0.07
Nodes (13): _button_outline_text_color(), _contrast_ratio(), _derive_dark_palette(), format_master_md(), _palette_is_dark(), _query_wants_dark(), _relative_luminance(), _resolve_color_mode() (+5 more)

### Community 62 - "time.Time"
Cohesion: 0.18
Nodes (6): loginLimiter, registerLimiter, exportRow(), formatExportTime(), newLoginLimiter(), newRegisterLimiter()

### Community 63 - "newTestEnvWithUsers"
Cohesion: 0.19
Nodes (15): User, testEnv, newTestEnvWithUsers(), newTestUser(), dashboardTotal(), TestDashboardCountsEveryRecordForFullAdmin(), TestDashboardDistrictAdminOnlyCountsOwnDistrict(), TestDashboardMemberFollowsTheChosenFilter() (+7 more)

### Community 64 - "core.py"
Cohesion: 0.08
Nodes (22): _contains_phrase(), _domain_keywords(), _exact_match_diagnostic(), _exact_row_identity(), _file_signature(), _get_bm25(), _legacy_successor_guidance(), _load_csv() (+14 more)

### Community 65 - "users_test.go"
Cohesion: 0.08
Nodes (31): 1. Ringkasan Eksekutif, 2.1 Model role lama — `old-app/code.gs:6-11`, 2.2 Aturan akses data — `old-app/code.gs:373-392` (dari `getRecruitmentsPaged`), 2.3 Sumber data lama, 2.4 Perilaku hak akses lain di UI lama (`old-app/index.html`), 2. Logika yang ada di `old-app` (yang seharusnya dipertahankan), 3. Kondisi aplikasi baru (verified), 4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi (+23 more)

### Community 67 - "Testing Plan — Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): 10. Checkpoint terakhir, 11. Target cakupan, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. Perintah verifikasi, 5. Matriks coverage, 8. Strategi test per lapisan, 9. Standar penulisan test (+4 more)

### Community 68 - "rekrutmen-tangsel/backend/internal/model.Aktivitas"
Cohesion: 0.13
Nodes (8): seedAktivitasIfEmpty(), Aktivitas, AktivitasToValues(), aktivitasToValues(), valuesToAktivitas(), Store, fileState, Store

### Community 69 - ".handleChangePassword"
Cohesion: 0.36
Nodes (6): bodyError(), decodeJSON(), writeFieldError(), Aktivitas(), Fields, Password()

### Community 70 - "6. Yang belum ada"
Cohesion: 0.29
Nodes (13): main(), newRegisterEnv(), TestRegisterCreatesAnAccount(), TestRegisterOptionsArePublic(), TestRegisterRejectsInvalidInput(), TestRegisterReportsTakenUsername(), TestRegisterWithoutAccountStoreIsUnavailable(), validRegistration() (+5 more)

### Community 75 - ".clientAddress"
Cohesion: 0.36
Nodes (6): lastForwardedFor(), remoteAddress(), TestClientAddressFallsBackToConnection(), TestClientAddressIgnoresForwardedHeaderByDefault(), TestClientAddressUsesLastForwardedEntry(), TestClientAddressWithoutPort()

### Community 76 - "supabase/store_test.go"
Cohesion: 0.26
Nodes (12): New(), matchRow(), newFakePostgREST(), newTestStore(), seedRow(), TestSupabaseAktivitasMaster(), TestSupabaseCreateAndGet(), TestSupabaseFindByWhatsApp() (+4 more)

### Community 81 - "api.js"
Cohesion: 0.43
Nodes (5): ApiError, buildUrl(), defaultMessage(), request(), unauthorizedHandlers

### Community 82 - "7. Smoke test manual"
Cohesion: 0.40
Nodes (5): 7.1 Autentikasi, 7.2 CRUD & validasi, 7.3 Master & pengaturan, 7.4 Responsif & izin, 7. Smoke test manual

### Community 83 - "Testing Plan — Fitur Ubah Password (Change Password)"
Cohesion: 0.29
Nodes (7): 2. Kriteria lulus (exit criteria), 3. Prasyarat, 5. Hasil uji manual yang sudah dijalankan, 6. Temuan, 7. Perintah verifikasi, 8. Urutan pengerjaan test yang disarankan, Testing Plan — Fitur Ubah Password (Change Password)

### Community 89 - "test_text_layout_resilience.py"
Cohesion: 0.20
Nodes (3): read_rows(), TestTextLayoutDataContracts, TestTextLayoutRetrieval

## Knowledge Gaps
- **3 isolated node(s):** `testEnv`, `rekrutmen-tangsel/backend`, `PATH`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 577 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Rekrutmen` connect `Rekrutmen` to `model.go`, `5.2 — Scope di API`, `Server`, `rekrutmen-tangsel/backend/internal/model.Aktivitas`, `New`, `User`, `Rekrutmen`, `dashboard.go`, `mapping_test.go`, `batch_test.go`, `time.Time`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `User` (e.g. with `5.1 — Granularity role` and `5.1 Kembalikan granularity role ✅`) actually correct?**
  _`User` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `testEnv`, `rekrutmen-tangsel/backend`, `PATH` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `model.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12585034013605442 - nodes in this community are weakly interconnected._
- **Why does `9.1.4 Item pekerjaan` connect `batch_test.go` to `local/store_test.go`, `ListView.vue`, `Rekrutmen`, `mapping_test.go`, `batch.go`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **Are the 29 inferred relationships involving `newTestEnv()` (e.g. with `TestBatchAcceptsTheFieldNamesOfThisApplication()` and `TestBatchAcceptsTheSpellingsOfThePreviousApplication()`) actually correct?**
  _`newTestEnv()` has 29 INFERRED edges - model-reasoned connections that need verification._
- **Should `BatchUpload.vue` be split into smaller, more focused modules?**
  _Cohesion score 0.07179487179487179 - nodes in this community are weakly interconnected._
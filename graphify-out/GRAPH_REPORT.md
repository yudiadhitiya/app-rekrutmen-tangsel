# Graph Report - app-rekrutmen-tangsel  (2026-11-24)

## Corpus Check
- 150 files · ~152,745 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 3, .csv 2, .service 2)

## Summary
- 2021 nodes · 5171 edges · 94 communities (67 shown, 27 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 729 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `62686e19`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- DesignSystemGenerator
- BatchUpload.vue
- Filter
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
- test_design_system_mode.py
- environment
- LoginView.vue
- testing.T
- test_core_data_quality.py
- CLAUDE.md
- rekrutmen-tangsel/backend
- 4. Test yang sudah ada
- Rekrutmen
- run-dev.sh
- search
- 4.6 Auth & ubah password — `backend/internal/auth/auth_test.go`, `backend/internal/api/password_test.go` (18 test)
- server/main.go
- Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)
- core.py
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
- TestStyleTaxonomy
- Manager
- Store
- test_skill_script_paths.py
- TestThresholdGate
- PasswordDialog.vue
- batch_test.go
- _row_identities
- local/store_test.go
- read_rows
- TestStackFlagWithDesignSystem
- newTestEnvWithUsers
- _normalize
- users_test.go
- detect_domain
- Testing Plan — Rekrutmen Tangsel
- rekrutmen-tangsel/backend/internal/model.Aktivitas
- TestTextLayoutDataContracts
- TestImeGuidance
- TestTextLayoutRetrieval
- .clientAddress
- supabase/store_test.go
- CatalogSummaryLineEndingsTest
- RtlCliTest
- TestFixtureValidation
- backup/main.go
- loginLimiter
- sheets/register_test.go
- test_data_contracts.py
- .handleGetSettings
- TestMetricMath
- 12. Prosedur deploy
- 6. Yang belum ada
- 6. P3 — Server & reverse proxy

## God Nodes (most connected - your core abstractions)
1. `User` - 65 edges
2. `search()` - 54 edges
3. `newTestEnv()` - 52 edges
4. `Store` - 52 edges
5. `Server` - 49 edges
6. `Rekrutmen` - 47 edges
7. `DesignSystemGenerator` - 38 edges
8. `search_stack()` - 36 edges
9. `Manager` - 33 edges
10. `New()` - 31 edges

## Surprising Connections (you probably didn't know these)
- `9. Standar penulisan test` --references--> `validBody()`  [INFERRED]
  TESTING_PLAN.md → backend/internal/api/api_test.go
- `1. Ringkasan Eksekutif` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `(b) Header lama tidak cocok dengan kunci yang dibaca — **data kosong buta**` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `9.1.5 Urutan pengerjaan —Kemajuan` --references--> `AppColumnSpec()`  [INFERRED]
  DEVELOPMENT_PLAN.md → backend/internal/store/sheets/mapping.go
- `5.5 Alias importer diperluas` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go

## Import Cycles
- None detected.

## Communities (94 total, 27 thin omitted)

### Community 0 - "DesignSystemGenerator"
Cohesion: 0.09
Nodes (4): DesignSystemGenerator, _resolve_dial(), TestReasoningMatch, TestEndToEndCoherence

### Community 1 - "BatchUpload.vue"
Cohesion: 0.07
Nodes (40): 7. P4 — HTTPS & keamanan ⛔, 13. Rencana Task — Unduh Data Rekrutmen (export XLSX), Daftar task, emit, errorMessage, errorTotal, fileInput, fileName (+32 more)

### Community 2 - "Filter"
Cohesion: 0.27
Nodes (8): 5.2 — Scope di API, Filter, MatchScope(), MetodeMatches(), MetodeMatchValues(), TestMatchScope(), TestScopeFor(), Scope

### Community 3 - "newTestEnv"
Cohesion: 0.26
Nodes (25): newTestEnv(), TestAktivitasInUseCannotBeDeleted(), TestAktivitasMasterIsAdminOnly(), TestCheckWhatsAppReportsExistence(), TestDuplicateWhatsAppReturnsConflict(), TestLoginRejectsBadPassword(), TestLoginSetsSessionAndMeReturnsUser(), TestLogoutClearsSession() (+17 more)

### Community 4 - "Server"
Cohesion: 0.21
Nodes (13): Server, bodyError(), decodeJSON(), newID(), writeError(), writeFieldError(), writeJSON(), writeStoreError() (+5 more)

### Community 5 - "ListView.vue"
Cohesion: 0.06
Nodes (37): canUpload(), METODE, METODE_LAMA, MIN_PASSWORD_LENGTH, normalizeMetode(), STATUS, statusTone(), TINGKAT_PENERIMAAN (+29 more)

### Community 6 - "Deployment Plan — Rekrutmen Tangsel"
Cohesion: 0.10
Nodes (21): 10. P7 — Monitoring & operasional, 11. P8 — Kontainer (opsional, non-blocking), 13. Checklist pra-rilis, 14. Checklist pasca-rilis, 15. Environment per stadium, 16. Yang harus dikerjakan berikutnya, 17. Checkpoint terakhir, 1. Cara memakai dokumen ini (+13 more)

### Community 7 - "App.vue"
Cohesion: 0.13
Nodes (18): onUnauthorized(), admin, handleLogout(), isPublicPage, navLinks, passwordOpen, route, router (+10 more)

### Community 8 - "FormView.vue"
Cohesion: 0.09
Nodes (25): 6. F4 — Frontend ✅, aktivitasMaster, aktivitasOptions, checkDuplicateWa(), clearErrors(), errorCount, errors, FIELD_IDS (+17 more)

### Community 9 - "Store"
Cohesion: 0.14
Nodes (7): colLetter(), hasTab(), isMissingTab(), tabTitles(), serviceKey, sheetInfo, Store

### Community 10 - "New"
Cohesion: 0.14
Nodes (20): Config, testEnv, New(), newTestEnvWithUploadLimit(), TestLoginWithLegacyAccount(), TestLoginWithWrongLegacyPassword(), limiterStorePath(), TestKnownGoodAccountKeepsWorkingWhenAddressIsBlocked() (+12 more)

### Community 11 - "batch-upload.spec.ts"
Cohesion: 0.09
Nodes (6): devDependencies, @playwright/test@^1.65.1, @playwright/test, RUN, RUN, ACCOUNTS

### Community 12 - "3. Fase kerja"
Cohesion: 0.10
Nodes (19): 1. Tujuan & batasan, 2. Keputusan arsitektur (blocking — perlu konfirmasi), 3. Fase kerja, 3. Status per fase, 4. Pemetaan kolom (draft — wajib dikonfirmasi di D-107), 5. Alur rollout, 6. Risiko, 7. Definition of Done (+11 more)

### Community 13 - "testing.T"
Cohesion: 0.11
Nodes (33): main(), TestPrepareUsersRotatesAnOldDefaultPassword(), districtUser(), HashPassword(), TestChangePasswordInSheetOnlyMode(), TestChangePasswordOfUnknownAccount(), TestChangePasswordRejectsWrongCurrentPassword(), TestChangePasswordUpdatesBothStores() (+25 more)

### Community 14 - "steps"
Cohesion: 0.05
Nodes (45): name, runs-on, steps, concurrency, cancel-in-progress, group, name, needs (+37 more)

### Community 15 - "Rekrutmen"
Cohesion: 0.13
Nodes (11): batchRow, checkMemberDistrict(), Rekrutmen, NormalizeMetode(), RekrutmenToValues(), rekrutmenRowValues(), rekrutmenToValues(), valuesToRekrutmen() (+3 more)

### Community 16 - "AktivitasView.vue"
Cohesion: 0.13
Nodes (16): confirmDialog(), confirmState, resolveConfirm(), adding, addItem(), admin, editForm, editingId (+8 more)

### Community 18 - "package.json"
Cohesion: 0.13
Nodes (17): dependencies, vue@^3.5.13, vue-router@^4.5.0, devDependencies, vite@^6.0.7, @vitejs/plugin-vue@^5.2.1, name, private (+9 more)

### Community 19 - "Aplikasi Rekrutmen Tangsel"
Cohesion: 0.14
Nodes (14): Akun awal, Aplikasi Rekrutmen Tangsel, Aturan password, Backend, Catatan, Dokumen Proyek, Fitur, Frontend (+6 more)

### Community 20 - "validate_data.py"
Cohesion: 0.11
Nodes (32): _catalog_date(), _check_app_interface_contract(), _check_catalog_contract(), _check_catalog_summary(), _check_color_contract(), _check_core_data_contract(), _check_file(), _check_font_catalog() (+24 more)

### Community 21 - "DashboardView.vue"
Cohesion: 0.13
Nodes (19): isMemberOnly(), cards, data, errorMessage, filter, filters, formatNumber(), kecamatanTotals (+11 more)

### Community 22 - "SettingsView.vue"
Cohesion: 0.14
Nodes (16): TITLES, providerLabel(), PROVIDERS, dismissToast(), showToast(), toast, toasts, currentProvider (+8 more)

### Community 23 - "PRD --- Aplikasi Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): 10. Kriteria Penerimaan MVP, 11. Di Luar Scope MVP, 1. Tujuan, 2. Pengguna, 3. Fitur MVP, 4. Struktur Data dan Aturan Input, 5. Alur Utama, 6. UI/UX (+4 more)

### Community 24 - "API Contract — Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): API Contract — Rekrutmen Tangsel, Auth, Dashboard, Error shape, Frontend routes (vue-router), Master aktivitas (job/activity categories), Master wilayah (read-only, all roles), Rekrutmen (+4 more)

### Community 25 - "test_design_system_mode.py"
Cohesion: 0.07
Nodes (15): _button_outline_text_color(), _contrast_ratio(), _derive_dark_palette(), _filter_anti_patterns_for_mode(), format_master_md(), _palette_is_dark(), _query_wants_dark(), _relative_luminance() (+7 more)

### Community 26 - "environment"
Cohesion: 0.05
Nodes (40): build, environment, healthcheck, image, ports, restart, volumes, command (+32 more)

### Community 27 - "LoginView.vue"
Cohesion: 0.08
Nodes (28): login(), answer1, answer2, errorMessage, errors, focusField(), handleRegister(), handleSubmit() (+20 more)

### Community 29 - "testing.T"
Cohesion: 0.12
Nodes (25): TestExportFollowsActiveFilter(), TestExportIsNotPaginated(), TestExportMemberHoldsOnlyOwnRecords(), TestExportRequiresAuth(), loginAs(), TestAttackerCannotLockOutTheAdmin(), TestKnownAddressKeepsWorkingUnderAttack(), TestLockedAccountStillRefusesWrongPasswords() (+17 more)

### Community 30 - "test_core_data_quality.py"
Cohesion: 0.16
Nodes (10): read_rows(), TestAccessibilityGuidance, TestChartsTypographyAndIcons, TestCurrentReactGuidance, TestSemanticColors, _check_chart_contract(), _check_typography_contract(), _configured_font_names() (+2 more)

### Community 33 - "4. Test yang sudah ada"
Cohesion: 0.15
Nodes (16): newFakeSheets(), newTestStore(), newTestStoreWith(), testPrivateKeyPEM(), TestSheetsAcceptsSpreadsheetURL(), TestSheetsAktivitasMaster(), TestSheetsDeleteKeepsOtherRowsIntact(), TestSheetsEnsureTabsCreatesHeaders() (+8 more)

### Community 34 - "Rekrutmen"
Cohesion: 0.06
Nodes (61): dashboardItem, dashboardKecamatanSummary, dashboardMethodSummary, dashboardResponse, dashboardTopItem, testEnv, buildDashboard(), containsName() (+53 more)

### Community 36 - "search"
Cohesion: 0.10
Nodes (7): _normalize_rtl_level(), _rtl_row_filter(), search(), _style_search_destination(), TestDiagnosticsContracts, TestSearchDomains, RtlFilterTest

### Community 37 - "4.6 Auth & ubah password — `backend/internal/auth/auth_test.go`, `backend/internal/api/password_test.go` (18 test)"
Cohesion: 0.15
Nodes (17): TestChangePasswordIsUsedByTheNextLogin(), TestChangePasswordKeepsTheSession(), TestChangePasswordNeedsCurrentPassword(), TestChangePasswordNeedsSession(), TestChangePasswordOfMemberAccount(), TestChangePasswordRejectsBrokenBody(), TestChangePasswordRejectsSamePassword(), TestChangePasswordRejectsShortPassword() (+9 more)

### Community 38 - "server/main.go"
Cohesion: 0.06
Nodes (65): batchRowError, batchSkipped, Bootstrap, env(), envBool(), envInt(), main(), prepareUsers() (+57 more)

### Community 39 - "Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)"
Cohesion: 0.06
Nodes (32): 1. Ringkasan Eksekutif, 2. Skor Audit Teknis, 3. Skor Heuristik Nielsen, 4.1 App Shell & Navigasi (`App.vue`), 4.2 Daftar Rekrutmen (`ListView.vue`), 4.3 Form Tambah/Ubah (`FormView.vue`), 4.4 Login & Registrasi (`LoginView.vue`), 4.5 Master Aktivitas (`AktivitasView.vue`) (+24 more)

### Community 40 - "core.py"
Cohesion: 0.07
Nodes (12): _passes_threshold(), _query_coverage(), _suggest_terms(), _detect_page_type(), format_page_override_md(), _generate_intelligent_overrides(), persist_design_system(), safe_slug() (+4 more)

### Community 41 - "test_core.py"
Cohesion: 0.16
Nodes (5): BM25, _get_bm25(), _suggest_identities(), TestBm25CoreBehavior, TestTokenizer

### Community 42 - "impor-responses/main.go"
Cohesion: 0.07
Nodes (47): buildPlan(), contains(), envOr(), fieldReasons(), indexOf(), inspectSpreadsheet(), main(), printSummary() (+39 more)

### Community 43 - "Menyiapkan Google Sheets"
Cohesion: 0.40
Nodes (5): Daftar anggota mandiri (self-registration), Lokasi file kredensial, Masuk dengan akun dari aplikasi lama, Mengimpor data dari Google Form, Menyiapkan Google Sheets

### Community 44 - "format_ascii_box"
Cohesion: 0.10
Nodes (11): ansi_ljust(), format_ascii_box(), add_wrapped(), wrap_text(), format_markdown(), generate_design_system(), hex_to_ansi(), _mode_support_labels() (+3 more)

### Community 45 - "User"
Cohesion: 0.08
Nodes (25): 5.1 — Granularity role, 5.1 Kembalikan granularity role ✅, 5.2 Terapkan scope di API, bukan di store ✅, 5.3 — Hak hapus & master data berbasis level, 5.3 Jaga hak hapus seperti lama ✅, 5.4 Migrasi data historis sekali jalan, 5.5 Alias importer diperluas, 5. Rekomendasi pemulihan logika lama (+17 more)

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

### Community 53 - "Manager"
Cohesion: 0.10
Nodes (7): LegacyAuthenticator, Manager, Registrar, User, isLegacyID(), legacyKnown(), saveUsers()

### Community 55 - "test_skill_script_paths.py"
Cohesion: 0.38
Nodes (3): resolve(), shipped_invocations(), SkillScriptPathsTest

### Community 57 - "PasswordDialog.vue"
Cohesion: 0.15
Nodes (14): api, ApiError, buildUrl(), defaultMessage(), request(), unauthorizedHandlers, confirmPassword, currentPassword (+6 more)

### Community 58 - "batch_test.go"
Cohesion: 0.40
Nodes (13): TestBatchAcceptsTheFieldNamesOfThisApplication(), TestBatchAcceptsTheSpellingsOfThePreviousApplication(), TestBatchIgnoresTheOwnerColumnsOfTheFile(), TestBatchIsOnlyForMembers(), TestBatchRejectsMoreThanTheAllowedRows(), TestBatchRejectsTwoRowsWithTheSameNumber(), TestBatchSavesNothingWhenOneRowIsWrong(), TestBatchSkipsANumberThatIsAlreadyStored() (+5 more)

### Community 59 - "_row_identities"
Cohesion: 0.33
Nodes (3): _exact_row_identity(), _row_identities(), _style_identity()

### Community 60 - "local/store_test.go"
Cohesion: 0.30
Nodes (18): New(), newStore(), rec(), TestAktivitasMasterCRUD(), TestAktivitasRenameCascadesToRekrutmen(), TestCreateAktivitasRejectsDuplicateName(), TestCreateRejectsDuplicateWhatsApp(), TestDataPersistsAcrossReopen() (+10 more)

### Community 63 - "newTestEnvWithUsers"
Cohesion: 0.24
Nodes (14): newTestEnvWithUsers(), newTestUser(), dashboardTotal(), TestDashboardCountsEveryRecordForFullAdmin(), TestDashboardDistrictAdminOnlyCountsOwnDistrict(), TestDashboardMemberFollowsTheChosenFilter(), TestDashboardMethodSummaryCountsTheStatuses(), TestDashboardRequiresAuthentication() (+6 more)

### Community 64 - "_normalize"
Cohesion: 0.11
Nodes (9): _contains_phrase(), _exact_match_diagnostic(), _legacy_successor_guidance(), _normalize(), _rewrite_query_for_domain(), _search_csv(), _search_csv_detailed(), _stack_query_requests_legacy() (+1 more)

### Community 65 - "users_test.go"
Cohesion: 0.09
Nodes (34): 1. Ringkasan Eksekutif, 2.1 Model role lama — `old-app/code.gs:6-11`, 2.2 Aturan akses data — `old-app/code.gs:373-392` (dari `getRecruitmentsPaged`), 2.3 Sumber data lama, 2.4 Perilaku hak akses lain di UI lama (`old-app/index.html`), 2. Logika yang ada di `old-app` (yang seharusnya dipertahankan), 3. Kondisi aplikasi baru (verified), 4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi (+26 more)

### Community 66 - "detect_domain"
Cohesion: 0.12
Nodes (8): detect_domain(), _domain_keywords(), _file_signature(), _load_csv(), _load_csv_snapshot(), _load_product_keywords(), _load_rows_or_empty(), TestDomainDetection

### Community 67 - "Testing Plan — Rekrutmen Tangsel"
Cohesion: 0.12
Nodes (17): 10. Checkpoint terakhir, 11. Target cakupan, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. Perintah verifikasi, 5. Matriks coverage, 7.1 Autentikasi, 7.2 CRUD & validasi (+9 more)

### Community 68 - "rekrutmen-tangsel/backend/internal/model.Aktivitas"
Cohesion: 0.23
Nodes (6): seedAktivitasIfEmpty(), Aktivitas, aktivitasToValues(), valuesToAktivitas(), Store, sheetRow

### Community 75 - ".clientAddress"
Cohesion: 0.24
Nodes (7): ServeFrontend(), lastForwardedFor(), remoteAddress(), TestClientAddressFallsBackToConnection(), TestClientAddressIgnoresForwardedHeaderByDefault(), TestClientAddressUsesLastForwardedEntry(), TestClientAddressWithoutPort()

### Community 76 - "supabase/store_test.go"
Cohesion: 0.14
Nodes (15): rowFromSpec(), splitPath(), TestListRekrutmenSanitizesTheFilter(), New(), matchRow(), newFakePostgREST(), newTestStore(), seedRow() (+7 more)

### Community 81 - "backup/main.go"
Cohesion: 0.20
Nodes (15): fileSHA256(), listFiles(), main(), prune(), readArchive(), runBackup(), runList(), runRestore() (+7 more)

### Community 82 - "loginLimiter"
Cohesion: 0.11
Nodes (17): limiterStore, loginLimiter, registerLimiter, uploadLimiter, newUploadLimiter(), TestLoginLimiterSurvivesARestart(), newLimiterStore(), TestLockedAccountSlowsDownTheAttempts() (+9 more)

### Community 83 - "sheets/register_test.go"
Cohesion: 0.07
Nodes (32): errorRegistrar, fakeRegistrar, fakeRegistrar, LegacyPasswordSetter, isSecurityQuestion(), validateRegistration(), Registration, normalizeAnswer() (+24 more)

### Community 84 - "test_data_contracts.py"
Cohesion: 0.09
Nodes (12): apply_decision_rules(), _object_without_duplicates(), parse_decision_rules(), _validate_action(), read_rows(), split_values(), style_identities(), TestGeneratedCatalogContract (+4 more)

### Community 87 - ".handleGetSettings"
Cohesion: 0.28
Nodes (4): sheetsMode(), sheetsTab(), Store, Settings

### Community 93 - "12. Prosedur deploy"
Cohesion: 0.50
Nodes (4): 12.1 Rilis pertama, 12.2 Rilis berikutnya, 12.3 Rollback, 12. Prosedur deploy

### Community 95 - "6. Yang belum ada"
Cohesion: 0.30
Nodes (11): TestRegisterRefusesAWeakPassword(), TestRegisterRefusesTheUsernameAsPassword(), newRegisterEnv(), TestRegisterCreatesAnAccount(), TestRegisterOptionsArePublic(), TestRegisterRejectsInvalidInput(), TestRegisterReportsTakenUsername(), TestRegisterWithoutAccountStoreIsUnavailable() (+3 more)

### Community 97 - "6. P3 — Server & reverse proxy"
Cohesion: 0.50
Nodes (4): 6. P3 — Server & reverse proxy, Contoh environment systemd, Contoh konfigurasi nginx, Contoh unit systemd

## Knowledge Gaps
- **3 isolated node(s):** `testEnv`, `rekrutmen-tangsel/backend`, `PATH`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 592 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **27 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `5.3 — Hak hapus & master data berbasis level` connect `User` to `newTestEnv`, `App.vue`, `DashboardView.vue`, `newTestEnvWithUsers`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `User` (e.g. with `5.1 — Granularity role` and `5.1 Kembalikan granularity role ✅`) actually correct?**
  _`User` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `testEnv`, `rekrutmen-tangsel/backend`, `PATH` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `DesignSystemGenerator` be split into smaller, more focused modules?**
  _Cohesion score 0.08658536585365853 - nodes in this community are weakly interconnected._
- **Why does `9.1.4 Item pekerjaan` connect `batch_test.go` to `Rekrutmen`, `ListView.vue`, `impor-responses/main.go`, `Rekrutmen`, `loginLimiter`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Are the 30 inferred relationships involving `newTestEnv()` (e.g. with `TestBatchAcceptsTheFieldNamesOfThisApplication()` and `TestBatchAcceptsTheSpellingsOfThePreviousApplication()`) actually correct?**
  _`newTestEnv()` has 30 INFERRED edges - model-reasoned connections that need verification._
- **Should `BatchUpload.vue` be split into smaller, more focused modules?**
  _Cohesion score 0.06533776301218161 - nodes in this community are weakly interconnected._
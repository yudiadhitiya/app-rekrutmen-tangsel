# Graph Report - app-rekrutmen-tangsel  (2026-11-20)

## Corpus Check
- 147 files · ~148,271 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 3, .csv 2, .service 2)

## Summary
- 1983 nodes · 5013 edges · 90 communities (64 shown, 26 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 706 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9233a3f3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- DesignSystemGenerator
- BatchUpload.vue
- Rekrutmen
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
- design_system.py
- environment
- LoginView.vue
- DesignSystemGenerator
- BatchUpload.vue
- CLAUDE.md
- rekrutmen-tangsel/backend
- 4. Test yang sudah ada
- Rekrutmen
- run-dev.sh
- search
- 4. Matriks test
- server/main.go
- Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)
- csv
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
- canAccess
- Manager
- Store
- TestStackFlagWithDesignSystem
- TestThresholdGate
- PasswordDialog.vue
- batch_test.go
- api.js
- local/store_test.go
- test_design_system_mode.py
- newTestEnvWithUsers
- core.py
- users_test.go
- detect_domain
- Testing Plan — Rekrutmen Tangsel
- rekrutmen-tangsel/backend/internal/model.Aktivitas
- auth/users.go
- split_values
- TestLandingAndStackContract
- .clientAddress
- supabase/store_test.go
- CatalogSummaryLineEndingsTest
- RtlCliTest
- TestFixtureValidation
- 7. Smoke test manual
- Testing Plan — Fitur Ubah Password (Change Password)
- _style_is_dark_primary
- _filter_anti_patterns_for_mode
- TestMetricMath
- testEnv

## God Nodes (most connected - your core abstractions)
1. `User` - 64 edges
2. `search()` - 54 edges
3. `Store` - 52 edges
4. `newTestEnv()` - 52 edges
5. `Server` - 47 edges
6. `Rekrutmen` - 45 edges
7. `DesignSystemGenerator` - 38 edges
8. `search_stack()` - 36 edges
9. `Manager` - 33 edges
10. `New()` - 28 edges

## Surprising Connections (you probably didn't know these)
- `1. Ringkasan Eksekutif` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `(b) Header lama tidak cocok dengan kunci yang dibaca — **data kosong buta**` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `9. Standar penulisan test` --references--> `validBody()`  [INFERRED]
  TESTING_PLAN.md → backend/internal/api/api_test.go
- `9.1.5 Urutan pengerjaan —Kemajuan` --references--> `AppColumnSpec()`  [INFERRED]
  DEVELOPMENT_PLAN.md → backend/internal/store/sheets/mapping.go
- `5.5 Alias importer diperluas` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go

## Import Cycles
- None detected.

## Communities (90 total, 26 thin omitted)

### Community 0 - "DesignSystemGenerator"
Cohesion: 0.13
Nodes (3): DesignSystemGenerator, TestReasoningMatch, TestEndToEndCoherence

### Community 1 - "BatchUpload.vue"
Cohesion: 0.06
Nodes (39): limiterStore, loginLimiter, registerLimiter, uploadLimiter, newUploadLimiter(), TestLoginLimiterSurvivesARestart(), newLimiterStore(), newLoginLimiter() (+31 more)

### Community 2 - "Rekrutmen"
Cohesion: 0.18
Nodes (8): batchRow, Filter, Rekrutmen, MetodeMatches(), MetodeMatchValues(), RekrutmenToValues(), rekrutmenToValues(), normalizeRows()

### Community 3 - "newTestEnv"
Cohesion: 0.13
Nodes (41): newTestEnv(), TestAktivitasInUseCannotBeDeleted(), TestAktivitasMasterIsAdminOnly(), TestCheckWhatsAppReportsExistence(), TestDuplicateWhatsAppReturnsConflict(), TestLoginRejectsBadPassword(), TestLoginSetsSessionAndMeReturnsUser(), TestLogoutClearsSession() (+33 more)

### Community 4 - "Server"
Cohesion: 0.21
Nodes (14): Server, bodyError(), decodeJSON(), newID(), ServeFrontend(), sheetsMode(), sheetsTab(), writeError() (+6 more)

### Community 5 - "ListView.vue"
Cohesion: 0.07
Nodes (27): canUpload(), activeFilters, canDeleteRecord, deleteRecord(), deletingId, emptyMessage, emptyTitle, errorMessage (+19 more)

### Community 6 - "Deployment Plan — Rekrutmen Tangsel"
Cohesion: 0.07
Nodes (29): 10. P7 — Monitoring & operasional, 11. P8 — Kontainer (opsional, non-blocking), 12.1 Rilis pertama, 12.2 Rilis berikutnya, 12.3 Rollback, 12. Prosedur deploy, 13. Checklist pra-rilis, 14. Checklist pasca-rilis (+21 more)

### Community 7 - "App.vue"
Cohesion: 0.13
Nodes (17): api, admin, handleLogout(), isPublicPage, navLinks, passwordOpen, route, router (+9 more)

### Community 8 - "FormView.vue"
Cohesion: 0.08
Nodes (31): METODE, METODE_LAMA, normalizeMetode(), STATUS, statusTone(), TINGKAT_PENERIMAAN, tingkatTone(), UPL_OPTIONS (+23 more)

### Community 9 - "Store"
Cohesion: 0.14
Nodes (7): colLetter(), hasTab(), isMissingTab(), tabTitles(), serviceKey, sheetInfo, Store

### Community 10 - "New"
Cohesion: 0.09
Nodes (34): Config, New(), newTestEnvWithUploadLimit(), TestLoginWithLegacyAccount(), TestLoginWithWrongLegacyPassword(), TestRegisterRefusesAWeakPassword(), TestRegisterRefusesTheUsernameAsPassword(), limiterStorePath() (+26 more)

### Community 11 - "batch-upload.spec.ts"
Cohesion: 0.09
Nodes (6): devDependencies, @playwright/test@^1.65.1, @playwright/test, RUN, RUN, ACCOUNTS

### Community 12 - "3. Fase kerja"
Cohesion: 0.10
Nodes (19): 1. Tujuan & batasan, 2. Keputusan arsitektur (blocking — perlu konfirmasi), 3. Fase kerja, 3. Status per fase, 4. Pemetaan kolom (draft — wajib dikonfirmasi di D-107), 5. Alur rollout, 6. Risiko, 7. Definition of Done (+11 more)

### Community 13 - "testing.T"
Cohesion: 0.15
Nodes (31): mustParse(), TestPasswordChangeInvalidatesOldSession(), TestChangePasswordInSheetOnlyMode(), TestChangePasswordOfUnknownAccount(), TestChangePasswordRejectsWrongCurrentPassword(), TestChangePasswordUpdatesBothStores(), TestChangePasswordUpdatesTheAccountOfThePreviousApplication(), TestChangePasswordUpdatesTheLocalAccount() (+23 more)

### Community 14 - "steps"
Cohesion: 0.05
Nodes (45): name, runs-on, steps, concurrency, cancel-in-progress, group, name, needs (+37 more)

### Community 15 - "Rekrutmen"
Cohesion: 0.17
Nodes (6): NormalizeMetode(), rekrutmenRowValues(), valuesToRekrutmen(), isUniqueViolation(), Fase 3 — Lapisan mapping kolom (1.5 hari) ✅, sheetRow

### Community 16 - "AktivitasView.vue"
Cohesion: 0.13
Nodes (16): confirmDialog(), confirmState, resolveConfirm(), adding, addItem(), admin, editForm, editingId (+8 more)

### Community 18 - "package.json"
Cohesion: 0.13
Nodes (18): dependencies, vue@^3.5.13, vue-router@^4.5.0, devDependencies, vite@^6.0.7, @vitejs/plugin-vue@^5.2.1, name, private (+10 more)

### Community 19 - "Aplikasi Rekrutmen Tangsel"
Cohesion: 0.14
Nodes (14): Akun awal, Aplikasi Rekrutmen Tangsel, Aturan password, Backend, Catatan, Dokumen Proyek, Fitur, Frontend (+6 more)

### Community 20 - "validate_data.py"
Cohesion: 0.06
Nodes (47): read_rows(), TestAccessibilityGuidance, TestChartsTypographyAndIcons, TestCurrentReactGuidance, TestSemanticColors, TestGeneratedCatalogContract, resolve(), shipped_invocations() (+39 more)

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

### Community 25 - "design_system.py"
Cohesion: 0.07
Nodes (13): _detect_page_type(), format_page_override_md(), generate_design_system(), _generate_intelligent_overrides(), persist_design_system(), safe_slug(), _write_persisted_file(), apply_decision_rules() (+5 more)

### Community 26 - "environment"
Cohesion: 0.05
Nodes (39): build, environment, healthcheck, image, ports, restart, volumes, command (+31 more)

### Community 27 - "LoginView.vue"
Cohesion: 0.08
Nodes (28): login(), answer1, answer2, errorMessage, errors, focusField(), handleRegister(), handleSubmit() (+20 more)

### Community 29 - "DesignSystemGenerator"
Cohesion: 0.14
Nodes (3): _resolve_dial(), read_rows(), TestReasoningContract

### Community 30 - "BatchUpload.vue"
Cohesion: 0.13
Nodes (16): emit, errorMessage, errorTotal, fileInput, fileName, onFileSelected(), rowCount, rowErrors (+8 more)

### Community 33 - "4. Test yang sudah ada"
Cohesion: 0.12
Nodes (20): newFakeSheets(), newTestStore(), newTestStoreWith(), testPrivateKeyPEM(), TestSheetsAcceptsSpreadsheetURL(), TestSheetsAktivitasMaster(), TestSheetsDeleteKeepsOtherRowsIntact(), TestSheetsEnsureTabsCreatesHeaders() (+12 more)

### Community 34 - "Rekrutmen"
Cohesion: 0.06
Nodes (61): dashboardItem, dashboardKecamatanSummary, dashboardMethodSummary, dashboardResponse, dashboardTopItem, testEnv, buildDashboard(), containsName() (+53 more)

### Community 36 - "search"
Cohesion: 0.07
Nodes (8): _normalize_rtl_level(), _rtl_row_filter(), search(), _style_search_destination(), TestSearchDomains, TestImeGuidance, RtlFilterTest, TestStyleTaxonomy

### Community 37 - "4. Matriks test"
Cohesion: 0.67
Nodes (3): 4.3 E2E (Playwright, `tests/e2e/ubah-password.spec.ts`), 4.4 Manual / eksploratif, 4. Matriks test

### Community 38 - "server/main.go"
Cohesion: 0.06
Nodes (58): batchRowError, batchSkipped, fileSHA256(), listFiles(), main(), prune(), readArchive(), runBackup() (+50 more)

### Community 39 - "Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)"
Cohesion: 0.06
Nodes (32): 1. Ringkasan Eksekutif, 2. Skor Audit Teknis, 3. Skor Heuristik Nielsen, 4.1 App Shell & Navigasi (`App.vue`), 4.2 Daftar Rekrutmen (`ListView.vue`), 4.3 Form Tambah/Ubah (`FormView.vue`), 4.4 Login & Registrasi (`LoginView.vue`), 4.5 Master Aktivitas (`AktivitasView.vue`) (+24 more)

### Community 40 - "csv"
Cohesion: 0.10
Nodes (5): read_rows(), RtlDataContractTest, read_rows(), read_rows(), TestTextLayoutDataContracts

### Community 41 - "test_core.py"
Cohesion: 0.14
Nodes (5): BM25, _get_bm25(), _suggest_identities(), TestBm25CoreBehavior, TestTokenizer

### Community 42 - "impor-responses/main.go"
Cohesion: 0.06
Nodes (56): 4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi, (a) Tab yang dibaca ≠ tab yang berisi data historis — **penyebab utama**, (b) Header lama tidak cocok dengan kunci yang dibaca — **data kosong buta**, (c) Kolom `UPL` terpecah di data lama, Cara verifikasi cepat (2 menit, tanpa kode), (d) Importer tidak mengenal header lama, (e) Mode `form` tidak menyuap data sumber ke daftar, buildPlan() (+48 more)

### Community 43 - "Menyiapkan Google Sheets"
Cohesion: 0.40
Nodes (5): Daftar anggota mandiri (self-registration), Lokasi file kredensial, Masuk dengan akun dari aplikasi lama, Mengimpor data dari Google Form, Menyiapkan Google Sheets

### Community 44 - "format_ascii_box"
Cohesion: 0.13
Nodes (9): ansi_ljust(), format_ascii_box(), add_wrapped(), wrap_text(), format_markdown(), hex_to_ansi(), _mode_support_labels(), section_header() (+1 more)

### Community 45 - "User"
Cohesion: 0.13
Nodes (9): fakeLegacyStore, legacyAuthenticator, fakeLegacy, fakeLookupLegacy, fakeProfiler, fakeStatusLegacy, isLegacyID(), saveUsers() (+1 more)

### Community 46 - "UI/UX Pro Max - Design Intelligence"
Cohesion: 0.11
Nodes (17): Before Delivering App UI, Example Workflow, If a search returns 0 results, Output Formats, Query Contract, Rule Categories by Priority, Running the search tool, Step 1: Analyze User Requirements (+9 more)

### Community 48 - "search_stack"
Cohesion: 0.08
Nodes (8): _exact_stack_identifier(), search_stack(), TestDiagnosticsContracts, _rows(), TestNativeDesktopStackFreshness, TestTextLayoutRetrieval, _rows(), TestWebStackFreshness

### Community 49 - "Pre-Delivery Checklist (canonical — the only one)"
Cohesion: 0.15
Nodes (12): Accessibility, Common Rules for Professional UI + Pre-Delivery Checklist, Icons & Visual Elements, Interaction, Interaction (App), Layout, Layout & Spacing, Light/Dark Mode (+4 more)

### Community 50 - "Quick Reference"
Cohesion: 0.15
Nodes (12): 10. Charts & Data (LOW), 1. Accessibility (CRITICAL), 2. Touch & Interaction (CRITICAL), 3. Performance (HIGH), 4. Style Selection (HIGH), 5. Layout & Responsive (HIGH), 6. Typography & Color (MEDIUM), 7. Animation (MEDIUM) (+4 more)

### Community 52 - "canAccess"
Cohesion: 0.09
Nodes (28): 1. Ringkasan Eksekutif, 2.1 Model role lama — `old-app/code.gs:6-11`, 2.2 Aturan akses data — `old-app/code.gs:373-392` (dari `getRecruitmentsPaged`), 2.3 Sumber data lama, 2.4 Perilaku hak akses lain di UI lama (`old-app/index.html`), 2. Logika yang ada di `old-app` (yang seharusnya dipertahankan), 3. Kondisi aplikasi baru (verified), 5.1 — Granularity role (+20 more)

### Community 53 - "Manager"
Cohesion: 0.12
Nodes (5): LegacyAuthenticator, Manager, Registrar, legacyKnown(), 1. Ruang lingkup

### Community 57 - "PasswordDialog.vue"
Cohesion: 0.22
Nodes (9): confirmPassword, currentPassword, emit, errorMessage, handleSubmit(), newPassword, props, submitting (+1 more)

### Community 58 - "batch_test.go"
Cohesion: 0.40
Nodes (13): TestBatchAcceptsTheFieldNamesOfThisApplication(), TestBatchAcceptsTheSpellingsOfThePreviousApplication(), TestBatchIgnoresTheOwnerColumnsOfTheFile(), TestBatchIsOnlyForMembers(), TestBatchRejectsMoreThanTheAllowedRows(), TestBatchRejectsTwoRowsWithTheSameNumber(), TestBatchSavesNothingWhenOneRowIsWrong(), TestBatchSkipsANumberThatIsAlreadyStored() (+5 more)

### Community 59 - "api.js"
Cohesion: 0.31
Nodes (7): 6. F4 — Frontend ✅, ApiError, buildUrl(), defaultMessage(), onUnauthorized(), request(), unauthorizedHandlers

### Community 60 - "local/store_test.go"
Cohesion: 0.30
Nodes (18): New(), newStore(), rec(), TestAktivitasMasterCRUD(), TestAktivitasRenameCascadesToRekrutmen(), TestCreateAktivitasRejectsDuplicateName(), TestCreateRejectsDuplicateWhatsApp(), TestDataPersistsAcrossReopen() (+10 more)

### Community 61 - "test_design_system_mode.py"
Cohesion: 0.12
Nodes (9): _button_outline_text_color(), _contrast_ratio(), _derive_dark_palette(), format_master_md(), _palette_is_dark(), _relative_luminance(), _select_palette_for_mode(), TestLuminance (+1 more)

### Community 63 - "newTestEnvWithUsers"
Cohesion: 0.15
Nodes (16): testEnv, newTestEnvWithUsers(), newTestUser(), dashboardTotal(), TestDashboardCountsEveryRecordForFullAdmin(), TestDashboardDistrictAdminOnlyCountsOwnDistrict(), TestDashboardMemberFollowsTheChosenFilter(), TestDashboardMethodSummaryCountsTheStatuses() (+8 more)

### Community 64 - "core.py"
Cohesion: 0.07
Nodes (23): _contains_phrase(), _domain_keywords(), _exact_match_diagnostic(), _exact_row_identity(), _file_signature(), _legacy_successor_guidance(), _load_csv(), _load_csv_snapshot() (+15 more)

### Community 65 - "users_test.go"
Cohesion: 0.07
Nodes (36): errorRegistrar, fakeRegistrar, fakePasswordLegacy, fakeRegistrar, LegacyLookup, LegacyPasswordSetter, LegacyProfiler, LegacyStatus (+28 more)

### Community 67 - "Testing Plan — Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): 10. Checkpoint terakhir, 11. Target cakupan, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. Perintah verifikasi, 5. Matriks coverage, 8. Strategi test per lapisan, 9. Standar penulisan test (+4 more)

### Community 68 - "rekrutmen-tangsel/backend/internal/model.Aktivitas"
Cohesion: 0.27
Nodes (5): seedAktivitasIfEmpty(), Aktivitas, aktivitasToValues(), valuesToAktivitas(), Store

### Community 69 - "auth/users.go"
Cohesion: 0.08
Nodes (35): Bootstrap, main(), prepareUsers(), rotateDefaultPasswords(), TestPrepareUsersCreatesOneAdminWithoutAKnownPassword(), TestPrepareUsersKeepsExistingAccounts(), TestPrepareUsersRefusesAShortPassword(), TestPrepareUsersRotatesAnOldDefaultPassword() (+27 more)

### Community 70 - "split_values"
Cohesion: 0.47
Nodes (3): split_values(), style_identities(), TestStyleIdentityContract

### Community 75 - ".clientAddress"
Cohesion: 0.36
Nodes (6): lastForwardedFor(), remoteAddress(), TestClientAddressFallsBackToConnection(), TestClientAddressIgnoresForwardedHeaderByDefault(), TestClientAddressUsesLastForwardedEntry(), TestClientAddressWithoutPort()

### Community 76 - "supabase/store_test.go"
Cohesion: 0.14
Nodes (15): rowFromSpec(), splitPath(), TestListRekrutmenSanitizesTheFilter(), New(), matchRow(), newFakePostgREST(), newTestStore(), seedRow() (+7 more)

### Community 82 - "7. Smoke test manual"
Cohesion: 0.40
Nodes (5): 7.1 Autentikasi, 7.2 CRUD & validasi, 7.3 Master & pengaturan, 7.4 Responsif & izin, 7. Smoke test manual

### Community 83 - "Testing Plan — Fitur Ubah Password (Change Password)"
Cohesion: 0.29
Nodes (7): 2. Kriteria lulus (exit criteria), 3. Prasyarat, 5. Hasil uji manual yang sudah dijalankan, 6. Temuan, 7. Perintah verifikasi, 8. Urutan pengerjaan test yang disarankan, Testing Plan — Fitur Ubah Password (Change Password)

### Community 85 - "_style_is_dark_primary"
Cohesion: 0.21
Nodes (4): _query_wants_dark(), _resolve_color_mode(), _style_is_dark_primary(), TestModeResolution

## Knowledge Gaps
- **4 isolated node(s):** `testEnv`, `testEnv`, `rekrutmen-tangsel/backend`, `PATH`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 590 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **26 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `9.1.4 Item pekerjaan` connect `batch_test.go` to `BatchUpload.vue`, `Rekrutmen`, `Rekrutmen`, `ListView.vue`, `impor-responses/main.go`?**
  _High betweenness centrality (0.044) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `User` (e.g. with `5.1 — Granularity role` and `5.1 Kembalikan granularity role ✅`) actually correct?**
  _`User` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `testEnv`, `testEnv`, `rekrutmen-tangsel/backend` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `DesignSystemGenerator` be split into smaller, more focused modules?**
  _Cohesion score 0.13333333333333333 - nodes in this community are weakly interconnected._
- **Why does `canUpload()` connect `ListView.vue` to `Rekrutmen`, `batch_test.go`, `App.vue`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **Are the 30 inferred relationships involving `newTestEnv()` (e.g. with `TestBatchAcceptsTheFieldNamesOfThisApplication()` and `TestBatchAcceptsTheSpellingsOfThePreviousApplication()`) actually correct?**
  _`newTestEnv()` has 30 INFERRED edges - model-reasoned connections that need verification._
- **Should `BatchUpload.vue` be split into smaller, more focused modules?**
  _Cohesion score 0.06233766233766234 - nodes in this community are weakly interconnected._
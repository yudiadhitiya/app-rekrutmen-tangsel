# Graph Report - app-rekrutmen-tangsel  (2026-11-12)

## Corpus Check
- 101 files · ~104,220 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 6 file(s) not represented in the graph (top: .csv 2, (none) 1, .example 1)

## Summary
- 1404 nodes · 3284 edges · 86 communities (58 shown, 28 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 366 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6e084afc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- api.go
- impor-responses/main.go
- seed.go
- api_test.go
- Server
- ListView.vue
- Deployment Plan — Rekrutmen Tangsel
- App.vue
- FormView.vue
- Store
- NormalizeMetode
- local/store_test.go
- 3. Fase kerja
- testing.T
- supabase/store_test.go
- Store
- AktivitasView.vue
- Store
- package.json
- Aplikasi Rekrutmen Tangsel
- validate_data.py
- api.js
- SettingsView.vue
- PRD --- Aplikasi Rekrutmen Tangsel
- API Contract — Rekrutmen Tangsel
- model.go
- sheets/store.go
- LoginView.vue
- test_data_contracts.py
- test_core_data_quality.py
- CLAUDE.md
- rekrutmen-tangsel/backend
- TestWebStackFreshness
- toast.js
- run-dev.sh
- search
- DesignSystemGenerator
- _select_palette_for_mode
- Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)
- design_system.py
- BM25
- Rekrutmen
- impor-responses/main.go
- format_ascii_box
- test_design_system_mode.py
- UI/UX Pro Max - Design Intelligence
- CatalogRefreshTest
- search_stack
- Pre-Delivery Checklist (canonical — the only one)
- Quick Reference
- Registration
- DesignSystemGenerator
- Testing Plan — Rekrutmen Tangsel
- valuesToRekrutmen
- Development Plan — Rekrutmen Tangsel
- TestThresholdGate
- 6. Yang belum ada
- .ListAktivitas
- CatalogSummaryLineEndingsTest
- split_values
- test_data_contracts.py
- read_rows
- TestStackFlagWithDesignSystem
- TestMetricMath
- RtlCliTest
- TestTextLayoutDataContracts
- TestImeGuidance
- TestFixtureValidation
- 12. Prosedur deploy
- _filter_anti_patterns_for_mode
- sha256Hex
- model.go
- testEnv
- mapping_test.go
- _normalize
- _contrast_ratio
- TestStyleTaxonomy
- test_skill_script_paths.py
- _row_identities
- TestGeneratedCatalogContract
- TestDiagnosticsContracts
- TestTextLayoutRetrieval
- 6. P3 — Server & reverse proxy

## God Nodes (most connected - your core abstractions)
1. `search()` - 54 edges
2. `Store` - 41 edges
3. `User` - 39 edges
4. `DesignSystemGenerator` - 38 edges
5. `search_stack()` - 36 edges
6. `Server` - 33 edges
7. `newTestEnv()` - 29 edges
8. `6. Yang belum ada` - 28 edges
9. `Manager` - 23 edges
10. `Rekrutmen()` - 23 edges

## Surprising Connections (you probably didn't know these)
- `5.5 Alias importer diperluas` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go
- `(d) Importer tidak mengenal header lama` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go
- `9. Standar penulisan test` --references--> `validBody()`  [INFERRED]
  TESTING_PLAN.md → backend/internal/api/api_test.go
- `1. Ringkasan Eksekutif` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `6. Yang belum ada` --references--> `RandomSecret()`  [INFERRED]
  TESTING_PLAN.md → backend/internal/auth/auth.go

## Import Cycles
- None detected.

## Communities (86 total, 28 thin omitted)

### Community 0 - "api.go"
Cohesion: 0.19
Nodes (4): env(), main(), ParseUsers(), RandomSecret()

### Community 1 - "impor-responses/main.go"
Cohesion: 0.17
Nodes (16): buildPlan(), sortedKeys(), TestBuildPlanImportsValidRow(), TestBuildPlanIsIdempotent(), TestBuildPlanMarksInvalidRowForReview(), TestBuildPlanReportsNewActivity(), TestBuildPlanSkipsDuplicateNumbers(), TestBuildPlanSkipsRowWithoutNumber() (+8 more)

### Community 2 - "seed.go"
Cohesion: 0.10
Nodes (22): 1. Ringkasan Eksekutif, 2.1 Model role lama — `old-app/code.gs:6-11`, 2.2 Aturan akses data — `old-app/code.gs:373-392` (dari `getRecruitmentsPaged`), 2.3 Sumber data lama, 2.4 Perilaku hak akses lain di UI lama (`old-app/index.html`), 2. Logika yang ada di `old-app` (yang seharusnya dipertahankan), 3. Kondisi aplikasi baru (verified), 5.1 Kembalikan granularity role (+14 more)

### Community 3 - "api_test.go"
Cohesion: 0.26
Nodes (25): newTestEnv(), TestAktivitasInUseCannotBeDeleted(), TestAktivitasMasterIsAdminOnly(), TestCheckWhatsAppReportsExistence(), TestDuplicateWhatsAppReturnsConflict(), TestLoginRejectsBadPassword(), TestLoginSetsSessionAndMeReturnsUser(), TestLogoutClearsSession() (+17 more)

### Community 4 - "Server"
Cohesion: 0.06
Nodes (38): Config, legacyAuthenticator, registerLimiter, Server, fakeLegacy, fakeProfiler, LegacyAuthenticator, LegacyProfiler (+30 more)

### Community 5 - "ListView.vue"
Cohesion: 0.06
Nodes (32): METODE, METODE_LAMA, normalizeMetode(), PROVIDERS, STATUS, statusTone(), TINGKAT_PENERIMAAN, tingkatTone() (+24 more)

### Community 6 - "Deployment Plan — Rekrutmen Tangsel"
Cohesion: 0.09
Nodes (22): 10. P7 — Monitoring & operasional, 11. P8 — Kontainer (opsional, non-blocking), 13. Checklist pra-rilis, 14. Checklist pasca-rilis, 15. Environment per stadium, 16. Yang harus dikerjakan berikutnya, 17. Checkpoint terakhir, 1. Cara memakai dokumen ini (+14 more)

### Community 7 - "App.vue"
Cohesion: 0.19
Nodes (13): admin, handleLogout(), isPublicPage, navLinks, roleLabel, route, router, auth (+5 more)

### Community 8 - "FormView.vue"
Cohesion: 0.12
Nodes (19): aktivitasMaster, aktivitasOptions, checkDuplicateWa(), clearErrors(), errors, form, formError, handleSubmit() (+11 more)

### Community 9 - "Store"
Cohesion: 0.16
Nodes (7): colLetter(), hasTab(), tabTitles(), serviceKey, sheetEntry, sheetInfo, Store

### Community 11 - "local/store_test.go"
Cohesion: 0.30
Nodes (18): New(), newStore(), rec(), TestAktivitasMasterCRUD(), TestAktivitasRenameCascadesToRekrutmen(), TestCreateAktivitasRejectsDuplicateName(), TestCreateRejectsDuplicateWhatsApp(), TestDataPersistsAcrossReopen() (+10 more)

### Community 12 - "3. Fase kerja"
Cohesion: 0.10
Nodes (19): 1. Tujuan & batasan, 2. Keputusan arsitektur (blocking — perlu konfirmasi), 3. Fase kerja, 3. Status per fase, 4. Pemetaan kolom (draft — wajib dikonfirmasi di D-107), 5. Alur rollout, 6. Risiko, 7. Definition of Done (+11 more)

### Community 13 - "testing.T"
Cohesion: 0.11
Nodes (33): testManager(), TestMigrateLegacyUserKeepsExistingAccount(), TestMigrateLegacyUserStoresAccountLocally(), TestVerifyLoginCompletesLevelOfMigratedAccount(), TestVerifyLoginFallsBackToLegacy(), TestVerifyLoginRejectsWrongLegacyPassword(), TestVerifyLoginReportsStoreFailure(), TestVerifyLoginUsesLocalAccount() (+25 more)

### Community 14 - "supabase/store_test.go"
Cohesion: 0.21
Nodes (13): rowFromSpec(), splitPath(), matchRow(), newFakePostgREST(), newTestStore(), seedRow(), TestSupabaseAktivitasMaster(), TestSupabaseCreateAndGet() (+5 more)

### Community 15 - "Store"
Cohesion: 0.15
Nodes (3): New(), normalizeRows(), Store

### Community 16 - "AktivitasView.vue"
Cohesion: 0.12
Nodes (17): confirmDialog(), confirmState, resolveConfirm(), adding, addItem(), admin, editForm, editingId (+9 more)

### Community 18 - "package.json"
Cohesion: 0.13
Nodes (18): dependencies, vue@^3.5.13, vue-router@^4.5.0, devDependencies, vite@^6.0.7, @vitejs/plugin-vue@^5.2.1, name, private (+10 more)

### Community 19 - "Aplikasi Rekrutmen Tangsel"
Cohesion: 0.12
Nodes (17): Akun awal, Aplikasi Rekrutmen Tangsel, Backend, Catatan, Daftar anggota mandiri (self-registration), Dokumen Proyek, Fitur, Frontend (+9 more)

### Community 20 - "validate_data.py"
Cohesion: 0.12
Nodes (30): _catalog_date(), _check_app_interface_contract(), _check_catalog_contract(), _check_catalog_summary(), _check_color_contract(), _check_core_data_contract(), _check_file(), _check_font_catalog() (+22 more)

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

### Community 26 - "sheets/store.go"
Cohesion: 0.14
Nodes (7): main(), checkKeyFilePerm(), New(), normalizeSpreadsheetID(), TestSheetsRejectsInvalidSpreadsheetID(), TestSheetsRejectsReadableKeyFile(), Config

### Community 27 - "LoginView.vue"
Cohesion: 0.10
Nodes (22): login(), answer1, answer2, errorMessage, handleRegister(), handleSubmit(), kecamatan, kecamatanOptions (+14 more)

### Community 29 - "test_data_contracts.py"
Cohesion: 0.11
Nodes (5): _passes_threshold(), _query_coverage(), _suggest_terms(), read_rows(), read_rows()

### Community 30 - "test_core_data_quality.py"
Cohesion: 0.14
Nodes (12): read_rows(), TestAccessibilityGuidance, TestChartsTypographyAndIcons, TestCurrentReactGuidance, TestSemanticColors, _check_chart_contract(), _check_typography_contract(), _configured_font_names() (+4 more)

### Community 34 - "toast.js"
Cohesion: 0.43
Nodes (5): TITLES, dismissToast(), showToast(), toast, toasts

### Community 36 - "search"
Cohesion: 0.11
Nodes (6): _normalize_rtl_level(), _rtl_row_filter(), search(), _style_search_destination(), TestSearchDomains, RtlFilterTest

### Community 37 - "DesignSystemGenerator"
Cohesion: 0.10
Nodes (4): _resolve_dial(), read_rows(), TestLandingAndStackContract, TestReasoningContract

### Community 38 - "_select_palette_for_mode"
Cohesion: 0.13
Nodes (6): _palette_is_dark(), _relative_luminance(), _select_palette_for_mode(), TestEndToEndCoherence, TestLuminance, TestPaletteSelection

### Community 39 - "Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)"
Cohesion: 0.06
Nodes (32): 1. Ringkasan Eksekutif, 2. Skor Audit Teknis, 3. Skor Heuristik Nielsen, 4.1 App Shell & Navigasi (`App.vue`), 4.2 Daftar Rekrutmen (`ListView.vue`), 4.3 Form Tambah/Ubah (`FormView.vue`), 4.4 Login & Registrasi (`LoginView.vue`), 4.5 Master Aktivitas (`AktivitasView.vue`) (+24 more)

### Community 40 - "design_system.py"
Cohesion: 0.10
Nodes (9): _detect_page_type(), format_page_override_md(), generate_design_system(), _generate_intelligent_overrides(), persist_design_system(), safe_slug(), _write_persisted_file(), format_output() (+1 more)

### Community 41 - "BM25"
Cohesion: 0.06
Nodes (17): BM25, _contains_phrase(), detect_domain(), _domain_keywords(), _file_signature(), _get_bm25(), _load_csv(), _load_csv_snapshot() (+9 more)

### Community 42 - "Rekrutmen"
Cohesion: 0.32
Nodes (20): contains(), Rekrutmen(), TestAktivitasMustExistInMaster(), TestEnumValuesMustBeValid(), TestFieldsAreTrimmed(), TestKecamatanMustBeOneOfSeven(), TestKelurahanMayBeEmpty(), TestKelurahanMustBelongToKecamatan() (+12 more)

### Community 43 - "impor-responses/main.go"
Cohesion: 0.17
Nodes (17): contains(), envOr(), indexOf(), inspectSpreadsheet(), main(), printSummary(), run(), splitList() (+9 more)

### Community 44 - "format_ascii_box"
Cohesion: 0.13
Nodes (9): ansi_ljust(), format_ascii_box(), add_wrapped(), wrap_text(), format_markdown(), hex_to_ansi(), _mode_support_labels(), section_header() (+1 more)

### Community 45 - "test_design_system_mode.py"
Cohesion: 0.24
Nodes (4): _query_wants_dark(), _resolve_color_mode(), _style_is_dark_primary(), TestModeResolution

### Community 46 - "UI/UX Pro Max - Design Intelligence"
Cohesion: 0.11
Nodes (17): Before Delivering App UI, Example Workflow, If a search returns 0 results, Output Formats, Query Contract, Rule Categories by Priority, Running the search tool, Step 1: Analyze User Requirements (+9 more)

### Community 48 - "search_stack"
Cohesion: 0.16
Nodes (6): _exact_stack_identifier(), _project_row(), search_stack(), _valid_max_results(), _rows(), TestNativeDesktopStackFreshness

### Community 49 - "Pre-Delivery Checklist (canonical — the only one)"
Cohesion: 0.15
Nodes (12): Accessibility, Common Rules for Professional UI + Pre-Delivery Checklist, Icons & Visual Elements, Interaction, Interaction (App), Layout, Layout & Spacing, Light/Dark Mode (+4 more)

### Community 50 - "Quick Reference"
Cohesion: 0.15
Nodes (12): 10. Charts & Data (LOW), 1. Accessibility (CRITICAL), 2. Touch & Interaction (CRITICAL), 3. Performance (HIGH), 4. Style Selection (HIGH), 5. Layout & Responsive (HIGH), 6. Typography & Color (MEDIUM), 7. Animation (MEDIUM) (+4 more)

### Community 51 - "Registration"
Cohesion: 0.27
Nodes (6): errorRegistrar, fakeRegistrar, isSecurityQuestion(), validateRegistration(), Registration, normalizeAnswer()

### Community 53 - "Testing Plan — Rekrutmen Tangsel"
Cohesion: 0.12
Nodes (16): 10. Checkpoint terakhir, 11. Target cakupan, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. Perintah verifikasi, 5. Matriks coverage, 7.1 Autentikasi, 7.2 CRUD & validasi (+8 more)

### Community 54 - "valuesToRekrutmen"
Cohesion: 0.14
Nodes (11): 4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi, (a) Tab yang dibaca ≠ tab yang berisi data historis — **penyebab utama**, (b) Header lama tidak cocok dengan kunci yang dibaca — **data kosong buta**, (c) Kolom `UPL` terpecah di data lama, Cara verifikasi cepat (2 menit, tanpa kode), (d) Importer tidak mengenal header lama, (e) Mode `form` tidak menyuap data sumber ke daftar, RekrutmenToValues() (+3 more)

### Community 55 - "Development Plan — Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): 11. Checkpoint terakhir, 12. Urutan pengerjaan berikutnya, 1. Cara memakai dokumen ini, 2. Ringkasan, 3. F1 — Fondasi & model data ✅, 6. F4 — Frontend ✅, 7. F5 — Keamanan 🔐, 8. F6 — Kualitas & CI (+4 more)

### Community 57 - "6. Yang belum ada"
Cohesion: 0.24
Nodes (16): TestLoginWithLegacyAccount(), TestLoginWithWrongLegacyPassword(), newRegisterEnv(), TestRegisterCreatesAnAccount(), TestRegisterOptionsArePublic(), TestRegisterRejectsInvalidInput(), TestRegisterReportsTakenUsername(), TestRegisterWithoutAccountStoreIsUnavailable() (+8 more)

### Community 60 - "split_values"
Cohesion: 0.47
Nodes (3): split_values(), style_identities(), TestStyleIdentityContract

### Community 61 - "test_data_contracts.py"
Cohesion: 0.21
Nodes (5): apply_decision_rules(), _object_without_duplicates(), parse_decision_rules(), _validate_action(), _check_reasoning_contract()

### Community 69 - "12. Prosedur deploy"
Cohesion: 0.50
Nodes (4): 12.1 Rilis pertama, 12.2 Rilis berikutnya, 12.3 Rollback, 12. Prosedur deploy

### Community 71 - "sha256Hex"
Cohesion: 0.33
Nodes (10): newFakeSheets(), sha256Hex(), legacyUsersRows(), seedUsersTab(), TestSheetsVerifyLegacyUser(), TestSheetsVerifyLegacyUserIsCaseInsensitive(), TestSheetsVerifyLegacyUserMapsLevels(), TestSheetsVerifyLegacyUserMapsRoles() (+2 more)

### Community 72 - "model.go"
Cohesion: 0.19
Nodes (7): normalizeHeader(), indexOfHeader(), isHex(), TestLegacyHashFormat(), TestLegacyLoginIntegration(), Aktivitas, LegacyUser

### Community 75 - "mapping_test.go"
Cohesion: 0.44
Nodes (9): FormColumnSpec(), ResolveColumnMap(), TestColumnMapRecord(), TestResolveColumnMapIgnoresPunctuation(), TestResolveColumnMapMatchesAliases(), TestResolveColumnMapPrefersFirstMatch(), TestResolveColumnMapReportsMissingRequiredColumn(), TestResolveColumnMapUsesAliasWhenNeeded() (+1 more)

### Community 76 - "_normalize"
Cohesion: 0.22
Nodes (5): _exact_match_diagnostic(), _legacy_successor_guidance(), _normalize(), _stack_query_requests_legacy(), _stack_row_filter()

### Community 77 - "_contrast_ratio"
Cohesion: 0.22
Nodes (4): _button_outline_text_color(), _contrast_ratio(), _derive_dark_palette(), format_master_md()

### Community 79 - "test_skill_script_paths.py"
Cohesion: 0.38
Nodes (3): resolve(), shipped_invocations(), SkillScriptPathsTest

### Community 80 - "_row_identities"
Cohesion: 0.33
Nodes (3): _exact_row_identity(), _row_identities(), _style_identity()

### Community 85 - "6. P3 — Server & reverse proxy"
Cohesion: 0.67
Nodes (3): 6. P3 — Server & reverse proxy, Contoh konfigurasi nginx, Contoh unit systemd

## Knowledge Gaps
- **2 isolated node(s):** `rekrutmen-tangsel/backend`, `PATH`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 416 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **28 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Development Plan — Rekrutmen Tangsel` connect `Development Plan — Rekrutmen Tangsel` to `Rekrutmen`, `api_test.go`, `DEPLOYMENT_PLAN.md`, `testing.T`?**
  _High betweenness centrality (0.103) - this node is a cross-community bridge._
- **Are the 4 inferred relationships involving `DesignSystemGenerator` (e.g. with `TestReasoningMatch` and `TestReasoningContract`) actually correct?**
  _`DesignSystemGenerator` has 4 INFERRED edges - model-reasoned connections that need verification._
- **What connects `rekrutmen-tangsel/backend`, `PATH` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `seed.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10333333333333333 - nodes in this community are weakly interconnected._
- **Why does `6. F4 — Frontend ✅` connect `Development Plan — Rekrutmen Tangsel` to `FormView.vue`, `api.js`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Should `Server` be split into smaller, more focused modules?**
  _Cohesion score 0.06163677592249021 - nodes in this community are weakly interconnected._
- **Why does `ApiError` connect `api.js` to `FormView.vue`, `Development Plan — Rekrutmen Tangsel`?**
  _High betweenness centrality (0.050) - this node is a cross-community bridge._
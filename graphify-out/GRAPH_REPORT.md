# Graph Report - app-rekrutmen-tangsel  (2026-11-12)

## Corpus Check
- 101 files · ~105,384 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 6 file(s) not represented in the graph (top: .csv 2, (none) 1, .example 1)

## Summary
- 1418 nodes · 3382 edges · 72 communities (53 shown, 19 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 383 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4c942cba`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- api.go
- detect_domain
- Analisis: Admin tidak melihat semua data (Rekrutmen Tangsel)
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
- _domain_keywords
- testing.T
- LoginView.vue
- test_data_contracts.py
- test_core_data_quality.py
- CLAUDE.md
- rekrutmen-tangsel/backend
- generate_design_system
- toast.js
- run-dev.sh
- search
- DesignSystemGenerator
- _select_palette_for_mode
- Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)
- design_system.py
- test_core.py
- Rekrutmen
- Menyiapkan Google Sheets
- format_ascii_box
- test_design_system_mode.py
- UI/UX Pro Max - Design Intelligence
- CatalogRefreshTest
- search_stack
- Pre-Delivery Checklist (canonical — the only one)
- Quick Reference
- DesignSystemGenerator
- TestThresholdGate
- newTestEnvWithUsers
- Store
- split_values
- test_data_contracts.py
- read_rows
- TestStackFlagWithDesignSystem
- TestMetricMath
- RtlCliTest
- TestImeGuidance
- TestFixtureValidation
- 12. Prosedur deploy
- _filter_anti_patterns_for_mode
- _normalize
- _contrast_ratio
- pathlib
- _row_identities
- TestGeneratedCatalogContract
- TestTextLayoutRetrieval
- 6. P3 — Server & reverse proxy

## God Nodes (most connected - your core abstractions)
1. `search()` - 54 edges
2. `User` - 41 edges
3. `Store` - 41 edges
4. `DesignSystemGenerator` - 38 edges
5. `search_stack()` - 36 edges
6. `Server` - 34 edges
7. `6. Yang belum ada` - 28 edges
8. `newTestEnv()` - 25 edges
9. `Rekrutmen` - 25 edges
10. `writeJSON()` - 23 edges

## Surprising Connections (you probably didn't know these)
- `9. Standar penulisan test` --references--> `validBody()`  [INFERRED]
  TESTING_PLAN.md → backend/internal/api/api_test.go
- `1. Ringkasan Eksekutif` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `(b) Header lama tidak cocok dengan kunci yang dibaca — **data kosong buta**` --references--> `valuesToRekrutmen()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/store.go
- `(d) Importer tidak mengenal header lama` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go
- `5.5 Alias importer diperluas` --references--> `FormColumnSpec()`  [INFERRED]
  ANALISIS_AKSES_DATA.md → backend/internal/store/sheets/mapping.go

## Import Cycles
- None detected.

## Communities (72 total, 19 thin omitted)

### Community 0 - "api.go"
Cohesion: 0.05
Nodes (63): errorRegistrar, fakeRegistrar, buildPlan(), contains(), envOr(), indexOf(), inspectSpreadsheet(), main() (+55 more)

### Community 2 - "Analisis: Admin tidak melihat semua data (Rekrutmen Tangsel)"
Cohesion: 0.12
Nodes (15): 1. Ringkasan Eksekutif, 2.1 Model role lama — `old-app/code.gs:6-11`, 2.2 Aturan akses data — `old-app/code.gs:373-392` (dari `getRecruitmentsPaged`), 2.3 Sumber data lama, 2.4 Perilaku hak akses lain di UI lama (`old-app/index.html`), 2. Logika yang ada di `old-app` (yang seharusnya dipertahankan), 4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi, 6. Ringkasan akar masalah (+7 more)

### Community 3 - "api_test.go"
Cohesion: 0.26
Nodes (25): newTestEnv(), TestAktivitasInUseCannotBeDeleted(), TestAktivitasMasterIsAdminOnly(), TestCheckWhatsAppReportsExistence(), TestDuplicateWhatsAppReturnsConflict(), TestLoginRejectsBadPassword(), TestLoginSetsSessionAndMeReturnsUser(), TestLogoutClearsSession() (+17 more)

### Community 4 - "Server"
Cohesion: 0.05
Nodes (49): 5.1 — Granularity role, 5.1 Kembalikan granularity role ✅, 5.2 Terapkan scope di API, bukan di store ✅, 5.3 Jaga hak hapus seperti lama, 5.4 Migrasi data historis sekali jalan, 5.5 Alias importer diperluas, 5. Rekomendasi pemulihan logika lama, Config (+41 more)

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
Cohesion: 0.17
Nodes (6): colLetter(), hasTab(), tabTitles(), serviceKey, sheetInfo, Store

### Community 10 - "NormalizeMetode"
Cohesion: 0.17
Nodes (8): NormalizeMetode(), RekrutmenToValues(), rekrutmenToValues(), valuesToRekrutmen(), isUniqueViolation(), Fase 3 — Lapisan mapping kolom (1.5 hari) ✅, sheetEntry, sheetRow

### Community 11 - "local/store_test.go"
Cohesion: 0.14
Nodes (30): New(), newStore(), rec(), TestAktivitasMasterCRUD(), TestAktivitasRenameCascadesToRekrutmen(), TestCreateAktivitasRejectsDuplicateName(), TestCreateRejectsDuplicateWhatsApp(), TestDataPersistsAcrossReopen() (+22 more)

### Community 12 - "3. Fase kerja"
Cohesion: 0.10
Nodes (19): 1. Tujuan & batasan, 2. Keputusan arsitektur (blocking — perlu konfirmasi), 3. Fase kerja, 3. Status per fase, 4. Pemetaan kolom (draft — wajib dikonfirmasi di D-107), 5. Alur rollout, 6. Risiko, 7. Definition of Done (+11 more)

### Community 13 - "sheets/store_test.go"
Cohesion: 0.15
Nodes (17): newFakeSheets(), newTestStore(), newTestStoreWith(), rowFromSpec(), splitPath(), testPrivateKeyPEM(), TestSheetsAcceptsSpreadsheetURL(), TestSheetsAktivitasMaster() (+9 more)

### Community 14 - "supabase/store_test.go"
Cohesion: 0.09
Nodes (28): matchRow(), newFakePostgREST(), newTestStore(), seedRow(), TestSupabaseAktivitasMaster(), TestSupabaseCreateAndGet(), TestSupabaseFindByWhatsApp(), TestSupabaseUniqueViolationMapsToConflict() (+20 more)

### Community 15 - "Rekrutmen"
Cohesion: 0.12
Nodes (13): 3. Kondisi aplikasi baru (verified), 5.2 — Scope di API, 7. Implementasi 5.1 dan 5.2, testEnv, Filter, Rekrutmen, MatchScope(), MetodeMatches() (+5 more)

### Community 16 - "AktivitasView.vue"
Cohesion: 0.12
Nodes (17): confirmDialog(), confirmState, resolveConfirm(), adding, addItem(), admin, editForm, editingId (+9 more)

### Community 18 - "package.json"
Cohesion: 0.13
Nodes (18): dependencies, vue@^3.5.13, vue-router@^4.5.0, devDependencies, vite@^6.0.7, @vitejs/plugin-vue@^5.2.1, name, private (+10 more)

### Community 19 - "Aplikasi Rekrutmen Tangsel"
Cohesion: 0.17
Nodes (12): Akun awal, Aplikasi Rekrutmen Tangsel, Backend, Catatan, Dokumen Proyek, Fitur, Frontend, Menjalankan (+4 more)

### Community 20 - "validate_data.py"
Cohesion: 0.12
Nodes (31): _catalog_date(), _check_app_interface_contract(), _check_catalog_contract(), _check_catalog_summary(), _check_chart_contract(), _check_color_contract(), _check_core_data_contract(), _check_file() (+23 more)

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

### Community 25 - "_domain_keywords"
Cohesion: 0.22
Nodes (6): _domain_keywords(), _file_signature(), _load_csv(), _load_csv_snapshot(), _load_product_keywords(), _load_rows_or_empty()

### Community 26 - "testing.T"
Cohesion: 0.13
Nodes (29): main(), newRegisterEnv(), TestRegisterCreatesAnAccount(), TestRegisterOptionsArePublic(), TestRegisterRejectsInvalidInput(), TestRegisterReportsTakenUsername(), TestRegisterWithoutAccountStoreIsUnavailable(), validRegistration() (+21 more)

### Community 27 - "LoginView.vue"
Cohesion: 0.10
Nodes (22): login(), answer1, answer2, errorMessage, handleRegister(), handleSubmit(), kecamatan, kecamatanOptions (+14 more)

### Community 29 - "test_data_contracts.py"
Cohesion: 0.11
Nodes (6): _passes_threshold(), _query_coverage(), _suggest_terms(), read_rows(), read_rows(), TestTextLayoutDataContracts

### Community 30 - "test_core_data_quality.py"
Cohesion: 0.14
Nodes (11): read_rows(), TestAccessibilityGuidance, TestChartsTypographyAndIcons, TestCurrentReactGuidance, TestSemanticColors, _check_typography_contract(), _configured_font_names(), contrast_ratio() (+3 more)

### Community 34 - "toast.js"
Cohesion: 0.43
Nodes (5): TITLES, dismissToast(), showToast(), toast, toasts

### Community 36 - "search"
Cohesion: 0.08
Nodes (7): _normalize_rtl_level(), _rtl_row_filter(), search(), _style_search_destination(), TestSearchDomains, RtlFilterTest, TestStyleTaxonomy

### Community 37 - "DesignSystemGenerator"
Cohesion: 0.12
Nodes (4): _resolve_dial(), read_rows(), TestLandingAndStackContract, TestReasoningContract

### Community 38 - "_select_palette_for_mode"
Cohesion: 0.18
Nodes (4): _palette_is_dark(), _relative_luminance(), TestEndToEndCoherence, TestLuminance

### Community 39 - "Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)"
Cohesion: 0.06
Nodes (32): 1. Ringkasan Eksekutif, 2. Skor Audit Teknis, 3. Skor Heuristik Nielsen, 4.1 App Shell & Navigasi (`App.vue`), 4.2 Daftar Rekrutmen (`ListView.vue`), 4.3 Form Tambah/Ubah (`FormView.vue`), 4.4 Login & Registrasi (`LoginView.vue`), 4.5 Master Aktivitas (`AktivitasView.vue`) (+24 more)

### Community 40 - "design_system.py"
Cohesion: 0.11
Nodes (7): _detect_page_type(), format_page_override_md(), _generate_intelligent_overrides(), persist_design_system(), safe_slug(), _write_persisted_file(), format_output()

### Community 41 - "test_core.py"
Cohesion: 0.12
Nodes (6): BM25, _get_bm25(), _suggest_identities(), TestBm25CoreBehavior, TestDiagnosticsContracts, TestTokenizer

### Community 42 - "Rekrutmen"
Cohesion: 0.32
Nodes (20): contains(), Rekrutmen(), TestAktivitasMustExistInMaster(), TestEnumValuesMustBeValid(), TestFieldsAreTrimmed(), TestKecamatanMustBeOneOfSeven(), TestKelurahanMayBeEmpty(), TestKelurahanMustBelongToKecamatan() (+12 more)

### Community 43 - "Menyiapkan Google Sheets"
Cohesion: 0.40
Nodes (5): Daftar anggota mandiri (self-registration), Lokasi file kredensial, Masuk dengan akun dari aplikasi lama, Mengimpor data dari Google Form, Menyiapkan Google Sheets

### Community 44 - "format_ascii_box"
Cohesion: 0.13
Nodes (9): ansi_ljust(), format_ascii_box(), add_wrapped(), wrap_text(), format_markdown(), hex_to_ansi(), _mode_support_labels(), section_header() (+1 more)

### Community 45 - "test_design_system_mode.py"
Cohesion: 0.21
Nodes (4): _query_wants_dark(), _resolve_color_mode(), _style_is_dark_primary(), TestModeResolution

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

### Community 57 - "newTestEnvWithUsers"
Cohesion: 0.11
Nodes (19): env(), main(), testEnv, newTestEnvWithUsers(), newTestUser(), TestLoginWithLegacyAccount(), TestLoginWithWrongLegacyPassword(), TestAdminKecamatanCannotOpenAnotherDistrictRecord() (+11 more)

### Community 58 - "Store"
Cohesion: 0.15
Nodes (6): Aktivitas, AktivitasToValues(), aktivitasToValues(), valuesToAktivitas(), fileState, Store

### Community 60 - "split_values"
Cohesion: 0.47
Nodes (3): split_values(), style_identities(), TestStyleIdentityContract

### Community 61 - "test_data_contracts.py"
Cohesion: 0.19
Nodes (5): apply_decision_rules(), _object_without_duplicates(), parse_decision_rules(), _validate_action(), _check_reasoning_contract()

### Community 69 - "12. Prosedur deploy"
Cohesion: 0.50
Nodes (4): 12.1 Rilis pertama, 12.2 Rilis berikutnya, 12.3 Rollback, 12. Prosedur deploy

### Community 76 - "_normalize"
Cohesion: 0.11
Nodes (9): _contains_phrase(), _exact_match_diagnostic(), _legacy_successor_guidance(), _normalize(), _rewrite_query_for_domain(), _search_csv(), _search_csv_detailed(), _stack_query_requests_legacy() (+1 more)

### Community 77 - "_contrast_ratio"
Cohesion: 0.15
Nodes (6): _button_outline_text_color(), _contrast_ratio(), _derive_dark_palette(), format_master_md(), _select_palette_for_mode(), TestPaletteSelection

### Community 79 - "pathlib"
Cohesion: 0.11
Nodes (5): CatalogSummaryLineEndingsTest, _load_generator(), resolve(), shipped_invocations(), SkillScriptPathsTest

### Community 80 - "_row_identities"
Cohesion: 0.33
Nodes (3): _exact_row_identity(), _row_identities(), _style_identity()

### Community 85 - "6. P3 — Server & reverse proxy"
Cohesion: 0.67
Nodes (3): 6. P3 — Server & reverse proxy, Contoh konfigurasi nginx, Contoh unit systemd

## Knowledge Gaps
- **2 isolated node(s):** `rekrutmen-tangsel/backend`, `PATH`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 416 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Development Plan — Rekrutmen Tangsel` connect `local/store_test.go` to `newTestEnvWithUsers`, `Rekrutmen`, `api_test.go`, `DEPLOYMENT_PLAN.md`?**
  _High betweenness centrality (0.132) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `User` (e.g. with `5.1 — Granularity role` and `5.1 Kembalikan granularity role ✅`) actually correct?**
  _`User` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `rekrutmen-tangsel/backend`, `PATH` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `api.go` be split into smaller, more focused modules?**
  _Cohesion score 0.052832438878950505 - nodes in this community are weakly interconnected._
- **Why does `6. F4 — Frontend ✅` connect `local/store_test.go` to `FormView.vue`, `api.js`?**
  _High betweenness centrality (0.109) - this node is a cross-community bridge._
- **Are the 4 inferred relationships involving `DesignSystemGenerator` (e.g. with `TestReasoningMatch` and `TestReasoningContract`) actually correct?**
  _`DesignSystemGenerator` has 4 INFERRED edges - model-reasoned connections that need verification._
- **Should `Analisis: Admin tidak melihat semua data (Rekrutmen Tangsel)` be split into smaller, more focused modules?**
  _Cohesion score 0.125 - nodes in this community are weakly interconnected._
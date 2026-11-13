# Graph Report - backend  (2026-11-13)

## Corpus Check
- 39 files · ~34,476 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 1 file(s) not represented in the graph (top: .example 1)

## Summary
- 522 nodes · 1759 edges · 19 communities (13 shown, 6 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 217 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9ea813d4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- User
- testing.T
- api.go
- impor-responses/main.go
- newTestStore
- Rekrutmen
- New
- Store
- Rekrutmen
- local/store_test.go
- Store
- Aktivitas
- supabase/store_test.go
- Store
- Store
- run-dev.sh
- .listRekrutmen
- rekrutmen-tangsel/backend

## God Nodes (most connected - your core abstractions)
1. `User` - 43 edges
2. `Rekrutmen` - 37 edges
3. `Server` - 34 edges
4. `Store` - 34 edges
5. `Manager` - 26 edges
6. `newTestEnv()` - 25 edges
7. `Rekrutmen()` - 24 edges
8. `writeJSON()` - 23 edges
9. `newTestStore()` - 22 edges
10. `Store` - 21 edges

## Surprising Connections (you probably didn't know these)
- `run()` --calls--> `AktivitasHeaders()`  [EXTRACTED]
  cmd/impor-responses/main.go → internal/store/sheets/export.go
- `run()` --calls--> `AktivitasToValues()`  [EXTRACTED]
  cmd/impor-responses/main.go → internal/store/sheets/export.go
- `run()` --calls--> `RekrutmenHeaders()`  [EXTRACTED]
  cmd/impor-responses/main.go → internal/store/sheets/export.go
- `run()` --calls--> `RekrutmenToValues()`  [EXTRACTED]
  cmd/impor-responses/main.go → internal/store/sheets/export.go
- `run()` --calls--> `New()`  [EXTRACTED]
  cmd/impor-responses/main.go → internal/store/sheets/store.go

## Import Cycles
- None detected.

## Communities (19 total, 6 thin omitted)

### Community 0 - "User"
Cohesion: 0.09
Nodes (24): legacyAuthenticator, Server, fakeLegacy, fakeLookupLegacy, fakeProfiler, LegacyAuthenticator, LegacyLookup, LegacyProfiler (+16 more)

### Community 1 - "testing.T"
Cohesion: 0.07
Nodes (64): env(), main(), testEnv, newTestEnv(), newTestEnvWithUsers(), newTestUser(), TestAktivitasInUseCannotBeDeleted(), TestAktivitasMasterIsAdminOnly() (+56 more)

### Community 2 - "api.go"
Cohesion: 0.13
Nodes (13): main(), normalizeHeader(), checkKeyFilePerm(), New(), normalizeSpreadsheetID(), TestSheetsRejectsInvalidSpreadsheetID(), TestSheetsRejectsReadableKeyFile(), indexOfHeader() (+5 more)

### Community 3 - "impor-responses/main.go"
Cohesion: 0.09
Nodes (39): buildPlan(), contains(), envOr(), indexOf(), inspectSpreadsheet(), main(), printSummary(), run() (+31 more)

### Community 4 - "newTestStore"
Cohesion: 0.08
Nodes (37): registration(), TestRegisteredUserCanSignIn(), TestRegisterUserAppendsRowWithLegacyLayout(), TestRegisterUserHashesPasswordWithSHA256(), TestRegisterUserNormalizesSecurityAnswers(), TestRegisterUserRejectsTakenUsername(), usersGrid(), newFakeSheets() (+29 more)

### Community 5 - "Rekrutmen"
Cohesion: 0.13
Nodes (27): fieldReasons(), All(), IsKecamatan(), IsKelurahan(), KelurahanByKecamatan(), Aktivitas(), contains(), Fields (+19 more)

### Community 6 - "New"
Cohesion: 0.10
Nodes (14): Config, errorRegistrar, fakeRegistrar, registerLimiter, New(), isSecurityQuestion(), newRegisterLimiter(), validateRegistration() (+6 more)

### Community 7 - "Store"
Cohesion: 0.14
Nodes (8): colLetter(), Store, hasTab(), tabTitles(), serviceKey, sheetEntry, sheetInfo, sheetRow

### Community 8 - "Rekrutmen"
Cohesion: 0.19
Nodes (9): Filter, Rekrutmen, MatchScope(), MetodeMatches(), MetodeMatchValues(), TestMatchScope(), valuesToRekrutmen(), normalizeRows() (+1 more)

### Community 9 - "local/store_test.go"
Cohesion: 0.25
Nodes (17): New(), newStore(), rec(), TestAktivitasMasterCRUD(), TestAktivitasRenameCascadesToRekrutmen(), TestCreateAktivitasRejectsDuplicateName(), TestCreateRejectsDuplicateWhatsApp(), TestDataPersistsAcrossReopen() (+9 more)

### Community 11 - "Aktivitas"
Cohesion: 0.23
Nodes (4): Aktivitas, AktivitasToValues(), aktivitasToValues(), valuesToAktivitas()

### Community 12 - "supabase/store_test.go"
Cohesion: 0.29
Nodes (10): matchRow(), newFakePostgREST(), newTestStore(), seedRow(), TestSupabaseAktivitasMaster(), TestSupabaseCreateAndGet(), TestSupabaseFindByWhatsApp(), TestSupabaseUniqueViolationMapsToConflict() (+2 more)

### Community 13 - "Store"
Cohesion: 0.22
Nodes (5): NormalizeMetode(), RekrutmenToValues(), rekrutmenToValues(), Store, isUniqueViolation()

## Knowledge Gaps
- **2 isolated node(s):** `rekrutmen-tangsel/backend`, `PATH`
  These have ≤1 connection - possible missing edges. (Counts symbols only; 36 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Rekrutmen` connect `Rekrutmen` to `User`, `api.go`, `impor-responses/main.go`, `Rekrutmen`, `local/store_test.go`, `Store`, `supabase/store_test.go`, `Store`, `.listRekrutmen`?**
  _High betweenness centrality (0.070) - this node is a cross-community bridge._
- **What connects `rekrutmen-tangsel/backend`, `PATH` to the rest of the system?**
  _2 weakly-connected nodes found - possible missing edges._
- **Should `User` be split into smaller, more focused modules?**
  _Cohesion score 0.09087719298245614 - nodes in this community are weakly interconnected._
- **Why does `User` connect `User` to `testing.T`, `api.go`, `newTestStore`?**
  _High betweenness centrality (0.070) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.07063063063063063 - nodes in this community are weakly interconnected._
- **Why does `Store` connect `Store` to `api.go`, `impor-responses/main.go`, `New`, `Rekrutmen`, `Store`, `Aktivitas`, `Store`, `Store`?**
  _High betweenness centrality (0.067) - this node is a cross-community bridge._
- **Should `api.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12709181365897784 - nodes in this community are weakly interconnected._
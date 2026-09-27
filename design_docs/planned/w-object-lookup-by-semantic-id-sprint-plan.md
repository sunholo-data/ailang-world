# Sprint plan: w-object-lookup-by-semantic-id (iteration 201)

**Base:** detached worktree at dev `c30b967` plus `design_docs/planned/w-object-lookup-by-semantic-id.md`; no production edits in this plan.  
**Authority:** that design doc in full, especially B1–B8, §c–d, Prototype manifest, Mutations and the Quorum log. Round 1 and 2 dispositions override the banked r0 prototype at `/Users/voightkampff/.ailang/state/world-iter201/prototype/` (`prototype.diff` is its unified diff). Scope is the whole row: M1, M2, M3. F-2 remains a finding.

## Gates and measurements

Run every Go and AILANG command with `export PATH=/opt/homebrew/bin:$PATH AILANG_BIN=$HOME/.pinned-ailang/ailang GOCACHE=/Users/voightkampff/Library/Caches/go-build TMPDIR=/private/tmp`. The baseline below is the untouched worktree, before `.plan-scratch` existed. `go test ./... -run '^$'` is the compile fence; `go build` alone would miss test doubles and `_test.go` errors. Run `go test -race -count=1` on `./host/store ./host/daemon ./cmd/ailang-worldd` after **each** landing, plus the focused commands below. Run `./scripts/verify_ail.sh` after the M2 error ADT and M3 route-table edits, with rc 0 required.

| untouched-base gate | observed rc | interpretation |
|---|---:|---|
| `go vet ./...` | 0 | clean |
| `go test ./... -run '^$'` | 0 | all packages and test files compile |
| `./scripts/verify_ail.sh` | 0 | 16 identities, 40 named tests, 9 package steps, PUB011 |
| `go test -race -count=1 ./host/store ./host/daemon ./cmd/ailang-worldd` | 1 | store passed; daemon and CLI socket/`httptest` bind EPERM: **UNINFORMATIVE UNDER SANDBOX** |
| `go test -count=1 ./...` | 1 | socket/`httptest` bind EPERM: **UNINFORMATIVE UNDER SANDBOX**; broker `TestHandlerTimeoutKillsTheWholeProcessGroup` also lost its fork marker under load. Do not chase row-115 capsule/broker `TestOrdinaryOverflowCarriesNoESRCH` EPERM, row-117 pkgproj descendant-pid, or row-32 capsule F5 wall-clock flakes if encountered. |

The planner copied the 15 MB base into nested scratch modules, layered the prototype by dependency, and ran both compile commands from each nested module. The revised sketch was also applied to the last scratch slice and verified. The r1/r2 index verifier, read-only guard, 503 mapping, `next` removal, and AC-9 test have **no banked implementation**; their final compile and race results must be measured by the executor. The r0 compile results establish that the prototype's interfaces and test doubles compose; they are not acceptance of the superseded behavior. Scratch was removed after measurement.

| landing boundary, in order | scratch assembly | added production LOC at boundary (approx.) | `go vet ./...` | `go test ./... -run '^$'` |
|---|---|---:|---:|---:|
| M1a query | base + r0 `objects_by_semantic_id.go` and store tests, canonical base schema | 56 | 0 | 0 |
| M1b bounded index + guard | M1a + r0 schema hunk (compile surrogate); replace that hunk with B4 provision/verify/guard during execution | ~58 more, ~114 M1 cumulative (r2 estimate) | 0 on surrogate; revised code unmeasured | 0 on surrogate; revised code unmeasured |
| M2 read route + 503 | M1 surrogate + r0 daemon hunks and tests; B6/round-1 changes reasoned from prototype | ~68 more (r2 estimate) | 0 on surrogate; revised code unmeasured | 0 on surrogate; revised code unmeasured |
| M3 CLI + frozen table | M2 surrogate + r0 CLI hunks; r2 sketch edits actually assembled | ~42 more Go (r2 estimate), ~3 sketch lines | 0 on surrogate; revised code unmeasured | 0 on surrogate; revised code unmeasured |

The sketch-edited M3 scratch slice ran `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` with **rc 0**. Each landing is below 150 added production Go/SQL LOC; count additions after implementation and split M1b before landing if its measured count exceeds 150. No scratch test exercised sockets: any later socket/`httptest` gate outcome inside this sandbox is **UNINFORMATIVE UNDER SANDBOX**.

### Execution rule for mutations

Each mutation below is a one-at-a-time literal replacement (`old` → `new`) in the named file **after** the named landing exists. Require exactly one match, run the named focused test, then restore the original bytes and compare SHA-256 before the next mutation. Rows using `DELETE` remove the exact old line/block. The new-code literals in M1b/M2/M3 are implementation contracts: write them with the shown spelling so mutations remain runnable. A green comment-only control should remain green. Never mutate the production worktree permanently. AC-8 and AC-11 are measured properties, not refusal branches; AC-8 is explicitly unguarded SQL equivalence, while the served method must refuse a missing index.

## M1a — store keyset query (first half of design M1)

**Files:** `host/store/objects_by_semantic_id.go`, `host/store/objects_by_semantic_id_test.go`. Take the r0 constant `MaxSemanticIDPage=500`, metadata-only keyset SQL, hash parsing, `InvalidLimitError`, and seven-object reverse-hash fixture. Change the r0 SELECT/Scan to include the stored `semantic_id` column as B3 specifies, rather than assigning the request ID into every result. Retain the predicate/order `semantic_id = ? AND hash_ref > ? ORDER BY hash_ref LIMIT ?`, with no `payload`, `OFFSET`, or rowid. Its result is `[]store.Object` with `Payload=nil`. The lookup is not exposed through a route yet. In M1b add the guard before the query; M1a's temporary direct SQL path is solely an unserved compile boundary.

**Acceptance:** `go test -race -count=1 ./host/store -run '^TestObjectsBySemanticIDNonUniqueOrderedAndPaged$'`; then the global vet/compile fence and touched-package race gate. AC-1 non-unique exact match, AC-2 stable ascending hash order independent of insertion, AC-3 strict cursor/page partition and store bound, AC-4's store-side empty non-nil slice. Keep the fixture's prefix decoy and per-page `len<=limit` assertion; both were necessary to kill r0 mutations.

| mutation | exact file and `old` → `new` | must red |
|---|---|---|
| MU1 ordering | `host/store/objects_by_semantic_id.go`: `ORDER BY hash_ref LIMIT ?` → `ORDER BY hash_ref DESC LIMIT ?` | `TestObjectsBySemanticIDNonUniqueOrderedAndPaged` (AC-2) |
| MU2 cursor | same: `hash_ref > ?` → `hash_ref >= ?` | same (AC-3 duplicate cursor) |
| MU3 exact ID | same: `semantic_id = ?` → `semantic_id GLOB ? || '*'` | same (AC-1 prefix decoy) |
| MU4 SQL limit | same: `LIMIT ?` → `LIMIT 500 + 0*?` | same (AC-3 page length) |
| MU5 store bound | same: `if limit < 1 || limit > MaxSemanticIDPage {` → `if false {` | same (AC-3 invalid limits 0 and 501) |
| MU-empty-store | same: `objects := make([]Object, 0, limit)` → `var objects []Object` | same (AC-4 non-nil empty result) |

## M1b — bounded index provisioning and cached refusal (rest of design M1)

**Files:** `host/store/store.go`, `host/store/objects_by_semantic_id.go`, `host/store/schema.sql` (one-line comment only), `host/store/objects_by_semantic_id_test.go`; optionally a new `host/store/lookup_index.go` to keep the verifier legible. Change the r0 schema hunk: **do not** take its `CREATE INDEX` or its claim that `schema.sql` reapplication builds the index. Add `provisionLookupIndex` after successful `enforceSchemaVersion` in writable `Open`, including fresh and v3 reopen, using `ExecContext` with `context.WithTimeout(..., 30*time.Second)`. Surface deadline/DDL errors and close/release the handle on failure. `OpenReadOnly` verifies but performs no DDL. Cache one immutable verdict per handle. Verify `PRAGMA index_list('objects')` row name and `partial=0`, then `PRAGMA index_xinfo('objects_by_semantic_id')` exactly two `key=1` rows: `(seqno,name,coll,desc)=(0,semantic_id,BINARY,0),(1,hash_ref,BINARY,0)`; ignore the `key=0` auxiliary row. An incompatible same-named index is never dropped or rebuilt. Guard before `QueryContext`, return `*LookupIndexUnavailableError` with fixed remediation text, and give the test-only `objectsBySemanticIDBeforeQuery` hook zero calls on refusal. Keep schema version 3.

**Acceptance:** `go test -race -count=1 ./host/store -run '^(TestObjectsBySemanticIDUsesIndex|TestSemanticIDIndexBuiltOnWritableReopen|TestLookupIndexGuardReadOnlyThenProvision|TestLookupIndexIncompatibleIsRefused|TestLookupIndexProvisionFailureIsSurfaced)$'`; global vet/compile and touched-package race gate. Revise r0 `TestObjectsBySemanticIDUsesIndex`: after `DROP INDEX`, execute **unguarded SQL directly** for AC-8; guarded `ObjectsBySemanticID` must refuse if a handle's cached verdict was absent at open. AC-7 asserts both fresh and v3-reopen provisioning and `EXPLAIN` index use without TEMP B-TREE. AC-13 asserts absent read-only index → typed error and hook 0, then writable provision → reopened read-only success. AC-14 has four separate fixtures (wrong key count/order, partial, NOCASE, DESC), each with writable refusal, byte-unchanged `sqlite_master.sql`, read-only typed refusal and hook 0. AC-15 forces an expired provisioning deadline and requires `errors.Is(err, context.DeadlineExceeded)` with nil Store. Assert PRAGMA semantics through the Go SQLite driver, not only the CLI probe in V29.

| mutation | exact file and `old` → `new` | must red |
|---|---|---|
| MU6 provision absent | `host/store/lookup_index.go`: `CREATE INDEX IF NOT EXISTS objects_by_semantic_id ON objects(semantic_id, hash_ref)` → `SELECT 1` | `TestObjectsBySemanticIDUsesIndex`, `TestSemanticIDIndexBuiltOnWritableReopen` (AC-7) |
| MU7 one key | same: `ON objects(semantic_id, hash_ref)` → `ON objects(semantic_id)` | `TestObjectsBySemanticIDUsesIndex` (AC-7; no ordered index) |
| MU20 guard off | `host/store/objects_by_semantic_id.go`: `if !s.lookupIndexAvailable {` → `if false {` | `TestLookupIndexGuardReadOnlyThenProvision` (AC-13) |
| MU21 read-only check off | `host/store/store.go`: `available, err := verifyLookupIndex(db)` → `available, err := true, error(nil)` in `OpenReadOnly` only | same (AC-13) |
| MU22 wrong columns | `host/store/lookup_index.go`: `name == want.name` → `true` in key comparison | `TestLookupIndexIncompatibleIsRefused/wrong-column-order` (AC-14a) |
| MU26 partial | same: `partial == 0` → `true` | `TestLookupIndexIncompatibleIsRefused/partial` (AC-14b) |
| MU27 collation | same: `coll == "BINARY"` → `true` | `TestLookupIndexIncompatibleIsRefused/nocase` (AC-14c) |
| MU28 DESC | same: `desc == 0` → `true` | `TestLookupIndexIncompatibleIsRefused/descending` (AC-14d) |
| MU29 destructive repair | `host/store/lookup_index.go`: `if !available { return fmt.Errorf("store: incompatible semantic-id lookup index") }` → `if !available { _, _ = db.Exec("DROP INDEX objects_by_semantic_id"); return provisionLookupIndex(db) }` | every AC-14 fixture's unchanged SQL assertion |
| MU23 ignored failure | `host/store/store.go`: `if err := provisionLookupIndex(db); err != nil { return nil, err }` → `_ = provisionLookupIndex(db)` | `TestLookupIndexProvisionFailureIsSurfaced` (AC-15) |

For MU22, seed `objects(hash_ref,semantic_id)` as the same-named index to exercise column order, plus a separate one-column fixture if desired. The test must compare `sqlite_master.sql` before and after attempted writable open. For MU29, a wrong-definition index's SQL must remain byte-identical: a successful replacement is a failure even when subsequent lookup would work. The hook lives immediately before `QueryContext`; initialize it only in tests and restore it via `t.Cleanup`.

## M2 — read route, seam and 503 class

**Files:** `host/daemon/daemon.go`, `host/daemon/handlers.go`, `host/daemon/objects_by_semantic_id_test.go`, `host/daemon/read_deadline_test.go`, `design_docs/sketches/worlddapi.ail`. Take r0 mux line, sixth `readStore` method, handler parsing/clamp/read context, and read-deadline doubles. Change the handler's `objectPageResponse` to **only** `Items []objectResponse`; delete `Next`, its assignment, and all tests or comments that rely on it. Revise cap test to use the 500th item's hash for `after`. Add `errors.As` mapping for `*store.LookupIndexUnavailableError` after timeout handling and before sanitized 500; return class `LookupIndexUnavailable`, status 503, fixed host-detail-free remediation. Add the sketch `ApiError` variant, `httpStatus` arm and test vector in this landing. Update the read-seam comment from five to six getters. Extend r0 `seedReadRoutes` and all three fake-store methods; add fault-injected 503 and sketch-mirror tests.

**Acceptance:** `go test -race -count=1 ./host/daemon -run '^(TestObjectsBySemanticIDRoute|TestObjectsBySemanticIDRouteCapAndDefault|TestObjectsBySemanticIDRouteStoreErrors|TestLookupIndexUnavailableIs503|TestDaemonReadDeadline|TestInternalErrorsAreSanitized|TestGETRoutesRejectOtherMethods)$'`; `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` (rc 0); global vet/compile and touched-package race gate. AC-4 demands byte-exact `{"items":[]}` for unknown and 400 for empty; AC-5 demands default 100, cap 500, final 20, union 520, no `next` key; AC-6 has separate malformed cursor, noninteger limit and no-name arms. AC-10/12 demand timeout 503, sanitized 500, and readCtx wiring through the shared harness. AC-16 demands the 503 class/message and the sketch's 503 vector. Route `?payload=true` still omits payload. POST to the GET route is 405.

| mutation | exact file and `old` → `new` | must red |
|---|---|---|
| MU8 empty name | `host/daemon/handlers.go`: `if name == "" {` → `if false {` in new handler | `TestObjectsBySemanticIDRoute` (AC-4) |
| MU9 after parse | same: `ref, err := parseRef(text, "after")` → `ref, err := hashref.HashRef{}, error(nil)` | same (AC-6) |
| MU10 limit parse | same: `if err != nil {` → `if false {` immediately after `strconv.Atoi(text)` | same (AC-6) |
| MU11 default | same: `requested := 0` → `requested := 500` in new handler | `TestObjectsBySemanticIDRouteCapAndDefault` (AC-5) |
| MU12 cap | same: `limit := clampLimit(requested)` → `limit := requested` | same (AC-5) |
| MU13r next | same: `Items []objectResponse `json:"items"`` → `Items []objectResponse `json:"items"`; Next string `json:"next"`` in `objectPageResponse` | same, raw JSON key assertion (AC-5) |
| MU14 timeout class | same: `if timedOut(ctx, err) {` → `if false && timedOut(ctx, err) {` in new handler | `TestDaemonReadDeadline` (AC-10/12) |
| MU15 context | same: `ctx, cancel := d.readCtx(r)` → `ctx, cancel := context.WithCancel(r.Context())` in new handler | `TestDaemonReadDeadline` (AC-12) |
| MU16 leak | same: `d.writeInternalError(w, r, err)` → `writeAPIError(w, "Internal", err.Error(), http.StatusInternalServerError)` in new handler | `TestInternalErrorsAreSanitized` (AC-10/12) |
| MU17 payload | same: `SemanticID: object.SemanticID, Provenance: object.Provenance,` → `SemanticID: object.SemanticID, Provenance: object.Provenance, Payload: &[]byte{},` | `TestObjectsBySemanticIDRoute` (AC-4/B7) |
| MU18 mux | `host/daemon/daemon.go`: `mux.HandleFunc("GET /v1/objects/by-semantic-id/{name...}", d.handleObjectsBySemanticID)` → DELETE | `TestObjectsBySemanticIDRoute` (AC-4/6) |
| MU24 mapping | `host/daemon/handlers.go`: `if errors.As(err, &indexErr) { writeAPIError(w, "LookupIndexUnavailable", lookupIndexUnavailableMessage, http.StatusServiceUnavailable); return }` → DELETE | `TestLookupIndexUnavailableIs503` (AC-16) |
| MU25 sketch status | `design_docs/sketches/worlddapi.ail`: `LookupIndexUnavailable(_) => 503` → `LookupIndexUnavailable(_) => 500` | `TestLookupIndexUnavailableIs503` mirror (AC-16) and `verify_ail.sh` vector |
| MU-empty-route | `host/daemon/handlers.go`: `Items: make([]objectResponse, 0, len(objects))` → `Items: nil` | `TestObjectsBySemanticIDRoute` byte-exact empty JSON (AC-4) |

## M3 — CLI and frozen-table extension

**Files:** `cmd/ailang-worldd/cli.go`, `cmd/ailang-worldd/main.go`, `cmd/ailang-worldd/cli_test.go` or `main_test.go`, `host/daemon/daemon.go` count comments, `host/daemon/session_middleware_test.go` comment, `host/daemon/read_deadline_test.go` comment, `design_docs/sketches/worlddapi.ail`, and a new `host/daemon/route_table_test.go`. Take r0 `object find` dispatch, per-segment `PathEscape`, query encoding, usage text. Add focused CLI verification of slash-bearing semanticId and optional cursor/limit. Add the missing `/v1/log` row and the new `/v1/objects/by-semantic-id/{name...}` row to sketch `routes()`, and correct `/v1/registry/{name}` to `{name...}`. Add an AC-9 parity test that extracts **all** literal `/v1` `mux.HandleFunc` method patterns from the production `Handler` function and the sketch's `routes()` list, compares sorted exact method/path pairs, and requires count 9 (8 GET, 1 POST). Its positive control is the pre-existing registry wildcard, so neither row can be normalized away. Update comments to nine frozen routes/eight GET and six read getters; workbench/A2A remain outside the table.

**Acceptance:** `go test -race -count=1 ./cmd/ailang-worldd ./host/daemon -run '^(TestObjectFindCLI|TestFrozenV1RouteTableMatchesMux|TestGETRoutesRejectOtherMethods)$'`; `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` (rc 0); global vet/compile and touched-package race gate. AC-9 exact parity and CLI 1:1 reachability. Run the env-gated AC-11 walk `WORLD_WALK_N=10000 go test -count=1 ./host/daemon -run '^TestWalkRetimedAtScale$' -v` only where loopback sockets are permitted; require scan/route same ref and source hash checks, report p50 and call counts. Inside this sandbox its socket result is **UNINFORMATIVE UNDER SANDBOX**. The r0 measurements (route CLI p50 0.101 s versus scan 136.8 s) are prior evidence, not a fresh acceptance result.

| mutation | exact file and `old` → `new` | must red |
|---|---|---|
| MU-9-mux | `host/daemon/daemon.go`: `mux.HandleFunc("GET /v1/objects/by-semantic-id/{name...}", d.handleObjectsBySemanticID)` → DELETE | `TestFrozenV1RouteTableMatchesMux` (AC-9) |
| MU-9-sketch | `design_docs/sketches/worlddapi.ail`: `{ method: "GET", path: "/v1/objects/by-semantic-id/{name...}" },` → DELETE | same (AC-9) |
| MU-9-registry | same: `/v1/registry/{name...}` → `/v1/registry/{name}` | same (AC-9 sibling drift) |
| MU-9-log | same: `{ method: "GET", path: "/v1/log" },` → DELETE | same (AC-9 existing missing row) |
| MU-CLI-route | `cmd/ailang-worldd/cli.go`: `path := "/v1/objects/by-semantic-id/" + strings.Join(parts, "/")` → `path := "/v1/objects/" + strings.Join(parts, "/")` | `TestObjectFindCLI` |
| MU-CLI-cursor | same: `query.Set("after", *after)` → `query.Set("after", "not-a-ref")` | `TestObjectFindCLI` |
| MU-CLI-limit | same: `query.Set("limit", strconv.Itoa(*limit))` → `query.Set("limit", "1")` | `TestObjectFindCLI` |

## Closure

After all red mutations and restoration, run vet, compile fence, touched-package `-race`, full `go test -count=1 ./...`, and `verify_ail.sh`. Report socket gates as **UNINFORMATIVE UNDER SANDBOX** when the sandbox refuses bind; rerun them in a socket-capable environment before landing. Report each final measured production LOC count, the no-version-bump/reopen behavior, the four index-definition refusal outcomes, sketch/mux count 9, and the AC-11 walk numbers. Do not add a registry commit mutation or chronological order here; those remain F-2/G-2 and row 104 respectively.

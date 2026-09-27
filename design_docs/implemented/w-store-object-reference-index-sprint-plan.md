# Sprint plan: w-store-object-reference-index (iteration 202, row 103)

**Base:** detached worktree `39dec32`. **Authority:** `design_docs/implemented/w-store-object-reference-index.md` in full, including D1–D5, V16/V17, the appendix and its mutation table. The quorum outcome accepts the controller's V16 premise check and V17 scale measurement, with sonnet's D2 fix (a). This plan implements the whole entry/world reverse-reference row. Row 104 owns `committedBy`.

## Gates and measured base

Prefix **every** Go or AILANG command, including commands below and mutation runs, with:

```sh
export PATH=/opt/homebrew/bin:$PATH AILANG_BIN=$HOME/.pinned-ailang/ailang
```

The baseline is the untouched `39dec32` tree. The controller's socket-capable run measured the race gate; this planner re-ran vet, compile fence and AILANG verification here. A sandbox `bind: operation not permitted` result from socket/`httptest` tests is **UNINFORMATIVE UNDER SANDBOX**, never a pass or a failure. CI must run the race gate.

| Untouched-base gate | Observed rc | Measurement / interpretation |
|---|---:|---|
| `go vet ./...` | 0 | Controller and planner: clean. |
| `go test ./... -run '^$'` | 0 | Controller and planner: all packages including `_test.go` compile. |
| `./scripts/verify_ail.sh` | 0 | Controller and planner: 16 identities, 40 named tests, nine package steps, PUB011. |
| `go test -race -count=1 ./host/store ./host/daemon ./host/workbench` | 0 outside sandbox | Controller: store 22.0 s; daemon 24.9 s; workbench 1.9 s. Planner's local run: store rc 0 (15.0 s), workbench rc 0 (1.85 s), daemon rc 1 solely on forbidden loopback/`httptest` binds. The local aggregate is **UNINFORMATIVE UNDER SANDBOX**. |
| V10 file-store probe at N=10,000 | 0 | Three samples: indexed sparse p50 52–55 µs versus unindexed 11.1–11.7 ms; four-index writable Open 44–45 ms; closed DB grows 4,120,576 bytes. Prior design measurement, not final-code timing. |
| V17 provisioning DDL | 0 | Controller: 41.2 ms at 10k, 588.9 ms at 100k, 7.635 s at 1M entries plus worlds. This supports reuse of the existing 30 s ceiling, not a new policy duration. |

At **each landing**, in order after its focused acceptance command: `go vet ./...`, `go test ./... -run '^$'`, and `go test -race -count=1 ./host/store ./host/daemon ./host/workbench`. Require rc 0 in a socket-capable CI environment. At final closure run `./scripts/verify_ail.sh` with rc 0 and `go test -count=1 ./...`; report sandbox bind failures separately as uninformative. A landing is incomplete if its compile fence or applicable focused test is red. Test files introduced at landing *k* must pass at *k*: never snapshot a later milestone's tests into an earlier commit, as occurred in iteration 201. Mutation runs happen only after the owning landing passes, one at a time, with original bytes restored and compared after each run.

## Landing boundaries

| Landing | Production LOC estimate | Complete, compiling boundary | Tests first introduced here |
|---|---:|---|---|
| M1a | 115 | Four named index specs and strict context-aware verifier; no Open behavior changed. | Definition fixtures and verifier controls. |
| M1b | 95 | Both writable Open branches provision; read-only Open caches verdict; typed refusal exists. | Lifecycle, cleanup, version and reopen. |
| M2a | 70 | Public types, cursor/limit/ref validation and error types; no SQL call yet. | Validator unit tests only. |
| M2b | 145 | Four covering keyset branches and guarded read are usable; no daemon exposure. | Store pages, scope, guard, cancellation, plans. |
| M3a | 85 | Seventh read seam with all fakes updated; ReferenceView renders independently. | Renderer and seam compile tests. |
| M3b | 145 | Object page checks and displays first 100 references, 503 adapter and truthful truncation text. | First-page, source/error and preservation tests. |
| M4a | 105 | Strict versioned cursor codec, tested directly; handler still first-page only. | Cursor codec tests. |
| M4b | 125 | `refsAfter` grammar, continuation, safe links and recorder walk. | Paging, URL, walk and route-parity tests. |

Estimates count added production Go lines, not tests/docs; measure actual `git diff --numstat` for each landing. If any landing exceeds about 150 added production lines, split it **before** landing and give each new boundary its own passing tests and all three gates. No boundary contains tests for features unavailable at that boundary. Never split a seam signature from its fake implementations or a view field from its renderer.

## Implementation contracts used by mutations

These are literal source contracts for the future executor, not claims about base code. Keep each **old** literal below exactly once in its named file at the owning landing; use named SQL constants so whole SQL strings are unique. Do not replace a literal that appears more than once. A mutation harness first asserts one occurrence, performs one exact `old` → `new` replacement, runs the named test, requires red, restores the original bytes, then checks the positive test green. A comment-only edit is the green control. In the table, `DELETE` means replacement by the empty string. The named test must check the behavior the mutation changes; no test merely checks that an implementation string exists.

`reference_index.go` must own the four literal index names, in order: `log_entries_by_transition_ref`, `log_entries_by_transition_fn_ref`, `log_entries_by_interpreter_ref`, `worlds_by_state_root`. Its spec row for the last is `{"worlds_by_state_root", "worlds", "state_root", "world_ref"},`. The verifier uses one shared context for DDL plus `PRAGMA index_list` and `PRAGMA index_xinfo`, closes rows before another query on the single connection, and has exactly one `partial == 0`, `coll == "BINARY"`, `desc == 0`, and `return valid && keys == 2, nil`. It rejects table/name, expression, column/order, DESC, NOCASE, partial and extra-key mismatch without dropping/rebuilding an incompatible same-named index. `schema.sql` gets a pointer comment only; `user_version` stays unchanged. `store.go` calls a single `openWritableReferenceIndexes(db)` helper from both in-memory and file-backed writable branches; that helper contains `if err := provisionReferenceIndexes(db); err != nil {`. On failure both DB and writer lock are released. Read-only Open verifies without DDL and assigns `referenceIndexesAvailable: available`; PRAGMA errors fail Open, while absent/incompatible indexes leave other reads usable. Use `lookupIndexProvisionDeadline` for the whole group; no new duration. When row 23 supplies caller Open context, derive this group from it.

`object_references.go` defines `MaxObjectReferencePage = 500`, `ReferenceKind` in the fixed order TransitionRef, TransitionFn, Interpreter, StateRoot, `ObjectReferenceCursor{Kind, EntryIndex, WorldRef}`, `ObjectReference{Cursor}`, and `ObjectReferences(ctx, ref, after, limit)`. The validator uses `limit < 1 || limit > MaxObjectReferencePage`, rejects malformed refs/cursors with typed errors, requires inactive fields zero, and returns non-nil empty slices. The read starts by checking `if !s.referenceIndexesAvailable {` before SQL. Keep the four branch SQL literals exactly as in the M2b table; use `entry_index > ?` (internal first-page sentinel `int64(-1)`) or `world_ref > ?` (empty-text sentinel), with `ORDER BY` and `LIMIT ?`. The branch caller uses `remaining := limit - len(items)` and exactly one `rows, err := s.db.QueryContext(ctx, query, ref.String(), afterKey, remaining)` line. Close and check each result set before another; discard partial results on any error. Check caller cancellation between branches. Parse world keys with `hashref.Parse`; reject negative stored entry indices. An entry used in two roles produces two edges. No OFFSET, payload, per-request index probe or scan fallback.

`daemon.go` adds exactly one `ObjectReferences` method to `readStore`, without a `/v1` route. Update blocking, recording and failing `read_deadline_test.go` fakes in the same landing. `workbench.go` uses the existing `d.readCtx(r)` once; the reverse read and every source point check share `ctx`. `checkedReferenceEdge` has separate `if !entryOK {` and `if !worldOK {` branches for named unavailable rows with no links; a differing relation takes `if actual != ref {` to sanitized 500. Log links use `pageHref(min(index, math.MaxInt64-WorkbenchPageLimit), index)`; world links use the world ref. Request `WorkbenchPageLimit+1`, display 100, and define `Truncated` by `len(refs) > WorkbenchPageLimit`. `ObjectView` holds an optional `*ReferenceView`; the M3a renderer emits no `referencedBy` section while it is nil, and M3b supplies it on every successfully inspected object page. The `ReferenceView` renderer owns the exact scoped text **“Entries/worlds only: transitionRef, transitionFn, interpreter, stateRoot. Registry, journal and object-interface inbound references are not included.”** and empty text **“none recorded in these entry/world fields”**; continued empty pages say **“no further references recorded in these entry/world fields”**. It shows “Showing 100 references; more recorded” when truncated. The named `committedBy` stop and its rendered text must remain byte-identical; remove only `objectReferencedByMissing`. A missing inspected object stays 404, unavailable reference index becomes HTML 503 `ReferenceIndexUnavailable` with the design's fixed remediation, timeout stays 503 `Timeout`, other errors sanitized 500. Preserve nine frozen `/v1` registrations.

`workbench.go` cursor codec uses base64url JSON with exactly `v`, `kind`, `key`, cap 512 encoded bytes, and rejects unknown/duplicate fields, unsupported `v`, noncanonical decimal integer and invalid hash. Its version guard is exactly `v != 1`. `refsAfter` is accepted only with `object` and optional `payload`; reject duplicates and world/from/entry combinations. Build Next from `refs[WorkbenchPageLimit-1]`, never the hidden lookahead row, using `url.Values` and HTML escaping; preserve object and payload. Show restart notice and first-page link on continuation. A missing source still occupies its cursor position. This is keyset pagination without a cross-page snapshot.

## M1a — index specs and verifier

**Files:** new `host/store/reference_index.go`, new `host/store/reference_index_test.go`, `host/store/schema.sql` (comment only). **Acceptance:**

```sh
go test -race -count=1 ./host/store -run '^(TestReferenceIndexDefinitions|TestReferenceIndexVerifierContext)$'
```

| Mutation | `host/store/reference_index.go`: exact old → new | Must red |
|---|---|---|
| MU1 | `{"worlds_by_state_root", "worlds", "state_root", "world_ref"},` → DELETE | `TestReferenceIndexDefinitions/all-four` |
| MU2 | `partial == 0` → `true` | `TestReferenceIndexDefinitions/partial` |
| MU3 | `coll == "BINARY"` → `true` | `TestReferenceIndexDefinitions/nocase` |
| MU4 | `desc == 0` → `true` | `TestReferenceIndexDefinitions/descending` |
| MU5 | `return valid && keys == 2, nil` → `return true, nil` | `TestReferenceIndexDefinitions/extra-key` |

## M1b — Open lifecycle and cached verdict

**Files:** `host/store/store.go`, `host/store/reference_index.go`, `host/store/reference_index_test.go`. Keep the existing semantic lookup index and its tests green. **Acceptance:**

```sh
go test -race -count=1 ./host/store -run '^(TestReferenceIndexLifecycle|TestReferenceIndexFailureCleanup|TestLookupIndexGuardReadOnlyThenProvision)$'
```

| Mutation | `host/store/store.go` unless stated: exact old → new | Must red |
|---|---|---|
| MU23 | `referenceIndexesAvailable: available` → `referenceIndexesAvailable: true` | `TestReferenceIndexLifecycle/read-only-absent` |
| MU24 | `host/store/reference_index.go`: `if err := provisionReferenceIndexes(db); err != nil {` → `if err := provisionReferenceIndexes(db); err != nil && false {` | `TestReferenceIndexFailureCleanup` |

## M2a — API types and validation

**Files:** new `host/store/object_references.go`, new `host/store/object_references_test.go`. Validator tests call the validator directly; do not introduce query tests yet. **Acceptance:**

```sh
go test -race -count=1 ./host/store -run '^TestObjectReferencesValidation$'
```

| Mutation | `host/store/object_references.go`: exact old → new | Must red |
|---|---|---|
| MU11 | `limit < 1 || limit > MaxObjectReferencePage` → `false` | `TestObjectReferencesValidation/limits` |

## M2b — bounded reverse read

**Files:** `host/store/object_references.go`, `host/store/object_references_test.go`. Re-run the authority appendix's N=10,000 file-store oracle against the final read, recording ordered-page equality, p50/pmax for sparse/hot/world/absent, covering plans, and Open/disk cost separately; do not turn latency into an assertion threshold. **Acceptance:**

```sh
go test -race -count=1 ./host/store -run '^(TestObjectReferencesPage|TestObjectReferencesScope|TestObjectReferencesGuard|TestObjectReferencesCancellation|TestObjectReferencesPlans)$'
WORLD103_N=10000 go test -count=1 -timeout 180s ./host/store -run '^TestObjectReferencesMeasuredAtScale$' -v
```

The second command is a planned opt-in test introduced in this landing; it must pass here. Its fixture uses production Store writes, matches complete pages to a raw SQL oracle and reports timings. Use 0, 1, 100, 101 and 501 edges; numeric indices 2/10, out-of-order world hashes, same-source multi-role, a cursor row absent from storage, and cross-kind boundaries. Include present but out-of-scope registry/journal/interface objects and entry-chain positive controls. An unavailable-index fixture must assert no reverse SQL ran.

| Mutation | `host/store/object_references.go`: exact old → new | Must red |
|---|---|---|
| MU6 | `if !s.referenceIndexesAvailable {` → `if false {` | `TestObjectReferencesGuard/no-sql` |
| MU7a | `SELECT entry_index FROM log_entries WHERE transition_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?` → `SELECT entry_index FROM log_entries WHERE transition_ref = ? AND entry_index >= ? ORDER BY entry_index LIMIT ?` | `TestObjectReferencesPage/transition-cursor` |
| MU7b | `SELECT entry_index FROM log_entries WHERE transition_fn_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?` → `SELECT entry_index FROM log_entries WHERE transition_fn_ref = ? AND entry_index >= ? ORDER BY entry_index LIMIT ?` | `TestObjectReferencesPage/function-cursor` |
| MU7c | `SELECT entry_index FROM log_entries WHERE interpreter_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?` → `SELECT entry_index FROM log_entries WHERE interpreter_ref = ? AND entry_index >= ? ORDER BY entry_index LIMIT ?` | `TestObjectReferencesPage/interpreter-cursor` |
| MU8 | `SELECT world_ref FROM worlds WHERE state_root = ? AND world_ref > ? ORDER BY world_ref LIMIT ?` → `SELECT world_ref FROM worlds WHERE state_root = ? AND world_ref >= ? ORDER BY world_ref LIMIT ?` | `TestObjectReferencesPage/world-cursor` |
| MU9a | `SELECT entry_index FROM log_entries WHERE transition_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?` → `SELECT entry_index FROM log_entries WHERE transition_ref <> ? AND entry_index > ? ORDER BY entry_index LIMIT ?` | `TestObjectReferencesScope/transition` |
| MU9b | `SELECT entry_index FROM log_entries WHERE transition_fn_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?` → `SELECT entry_index FROM log_entries WHERE transition_fn_ref <> ? AND entry_index > ? ORDER BY entry_index LIMIT ?` | `TestObjectReferencesScope/function` |
| MU9c | `SELECT entry_index FROM log_entries WHERE interpreter_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?` → `SELECT entry_index FROM log_entries WHERE interpreter_ref <> ? AND entry_index > ? ORDER BY entry_index LIMIT ?` | `TestObjectReferencesScope/interpreter` |
| MU9d | `SELECT world_ref FROM worlds WHERE state_root = ? AND world_ref > ? ORDER BY world_ref LIMIT ?` → `SELECT world_ref FROM worlds WHERE state_root <> ? AND world_ref > ? ORDER BY world_ref LIMIT ?` | `TestObjectReferencesScope/world` |
| MU10 | `remaining := limit - len(items)` → `remaining := limit` | `TestObjectReferencesPage/mixed-kind-bound` |
| MU12 | `rows, err := s.db.QueryContext(ctx, query, ref.String(), afterKey, remaining)` → `rows, err := s.db.QueryContext(context.Background(), query, ref.String(), afterKey, remaining)` | `TestObjectReferencesCancellation/blocked-connection` |
| MU13 | `kind < after.Kind` → `kind <= after.Kind` | `TestObjectReferencesPage/within-kind` |
| MU14 | `SELECT entry_index FROM log_entries WHERE transition_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?` → `SELECT entry_index FROM log_entries WHERE transition_ref = ? AND entry_index > ? ORDER BY entry_index DESC LIMIT ?` | `TestObjectReferencesPage/numeric-order` |
| MU25 | `afterKey = int64(-1)` → `afterKey = int64(0)` | `TestObjectReferencesPage/index-zero` |

For MU7–MU9 and MU14, the `old` is the whole unique SQL literal, not the repeated fragment. The `new` is the whole old literal with exactly the indicated fragment changed; the harness asserts whole-literal one-match before replacement.

## M3a — read seam and renderer

**Files:** `host/daemon/daemon.go`, `host/daemon/read_deadline_test.go`, `host/workbench/render.go`, `host/workbench/render_test.go`. Add the seventh seam method and all direct fake implementations together. Introduce `ReferenceView` and render it from a direct render test; the daemon still supplies no references at this boundary. Existing object-page tests remain green. **Acceptance:**

```sh
go test -race -count=1 ./host/daemon ./host/workbench -run '^(TestReferenceView|TestWorkbenchObjectProvenanceWalk|TestFrozenV1RouteTableMatchesMux)$'
```

| Mutation | `host/workbench/render.go`: exact old → new | Must red |
|---|---|---|
| MU21 | `none recorded in these entry/world fields` → DELETE | `TestReferenceView/empty` |

## M3b — checked first object page

**Files:** `host/daemon/workbench.go`, `host/daemon/workbench_test.go`, new `host/daemon/workbench_references_test.go`, `host/workbench/render.go`. Remove only the `referencedBy` missing-stop literal. Keep the `committedBy` named stop byte-identical and assert it. Show truncation text at 101 but no Next link until M4b. A source missing between reverse read and point check is a named unavailable row; mismatch or malformed stored source key is 500. Test all source links by fetching their target through `Handler` + `httptest.NewRecorder`. **Acceptance:**

```sh
go test -race -count=1 ./host/daemon ./host/workbench -run '^(TestWorkbenchReferenceFirstPage|TestWorkbenchMissingReferenceSource|TestWorkbenchReferenceMismatch|TestWorkbenchReferenceUnavailable|TestWorkbenchReferenceCancellation|TestWorkbenchObjectProvenanceWalk|TestFrozenV1RouteTableMatchesMux|TestReferenceView)$'
```

| Mutation | Exact file and `old` → `new` | Must red |
|---|---|---|
| MU15a | `host/daemon/workbench.go`: `if !entryOK {` → `if false {` | `TestWorkbenchMissingReferenceSource/log` |
| MU15b | `host/daemon/workbench.go`: `if !worldOK {` → `if false {` | `TestWorkbenchMissingReferenceSource/world` |
| MU16 | `host/daemon/workbench.go`: `if actual != ref {` → `if false {` | `TestWorkbenchReferenceMismatch` |
| MU17 | `host/daemon/workbench.go`: `writeWorkbenchError(w, http.StatusServiceUnavailable, "ReferenceIndexUnavailable", referenceIndexUnavailableMessage)` → `writeWorkbenchError(w, http.StatusOK, "ReferenceIndexUnavailable", referenceIndexUnavailableMessage)` | `TestWorkbenchReferenceUnavailable` |
| MU18 | `host/daemon/workbench.go`: `WorkbenchPageLimit+1` → `WorkbenchPageLimit` | `TestWorkbenchReferenceFirstPage/101` |
| MU19 | `host/daemon/workbench.go`: `len(refs) > WorkbenchPageLimit` → `len(refs) >= WorkbenchPageLimit` | `TestWorkbenchReferenceFirstPage/100` |
| MU22 | `host/daemon/workbench.go`: `d.reads.ObjectReferences(ctx, ref, nil, WorkbenchPageLimit+1)` → `d.reads.ObjectReferences(context.Background(), ref, nil, WorkbenchPageLimit+1)` | `TestWorkbenchReferenceCancellation` |

MU18's old literal occurs exactly once in the file: the reverse call is the sole use of `WorkbenchPageLimit+1`. MU22 replaces the full call and therefore also has one match; run them separately.

## M4a — strict cursor codec

**Files:** new `host/daemon/workbench_reference_cursor.go`, new `host/daemon/workbench_reference_cursor_test.go`. Implement encode/decode independent of handler grammar. Decode exactly the three fields, reject duplicate/unknown fields and trailing JSON, verify canonical key format and 512-byte encoded cap. **Acceptance:**

```sh
go test -race -count=1 ./host/daemon -run '^TestWorkbenchReferenceCursor$'
```

| Mutation | `host/daemon/workbench_reference_cursor.go`: exact old → new | Must red |
|---|---|---|
| MU26 | `v != 1` → `false` | `TestWorkbenchReferenceCursor/version` |

## M4b — continuation and real-question walk

**Files:** `host/daemon/workbench.go`, `host/daemon/workbench_references_test.go`, `host/daemon/route_table_test.go` (assert unchanged route set), `host/workbench/render.go`, `host/workbench/render_test.go`. Query validation accepts only object/payload/refsAfter combinations. `TestWorkbenchReferenceWalk` builds a file-backed Store fixture with an object whose semantic ID is `incident/row103/transition`, known payload bytes, a committed entry that uses it as `transitionRef`, a shared interpreter in several entries and a state root in several worlds. Its literal question is “Which entry references incident/row103/transition, and what transition body does that entry name?” Call `ObjectsBySemanticID` for that ID; GET the located object page through `Daemon.Handler().ServeHTTP(httptest.NewRecorder(), request)`; extract and GET its labelled `transitionRef` `referencedBy` link through the same Handler; assert that the selected-entry panel renders the exact `transitionRef` hash. Follow that rendered transition-ref object link and assert the payload bytes match the fixture. Also walk the shared interpreter and state-root backlinks. Log the question, answer, object/entry/world refs, each GET status, end-to-end elapsed time and remaining stops with `t.Logf`; elapsed time is a measurement, **not** a threshold. No listening server. First-page and continuation links use exact object/payload state, and 100/101 partition checks prove no skipped lookahead. MaxInt64 selected-entry link uses bounded `from`. **Acceptance:**

```sh
go test -race -count=1 ./host/daemon ./host/workbench -run '^(TestWorkbenchReferencePaging|TestWorkbenchReferenceCursor|TestWorkbenchReferenceWalk|TestWorkbenchReferenceGrammar|TestWorkbenchObjectProvenanceWalk|TestFrozenV1RouteTableMatchesMux|TestReferenceView)$' -v
```

| Mutation | `host/daemon/workbench.go`: exact old → new | Must red |
|---|---|---|
| MU20 | `refs[WorkbenchPageLimit-1]` → `refs[WorkbenchPageLimit]` | `TestWorkbenchReferencePaging/no-skipped-lookahead` |

## Closure and scope

Run each mutation after its owning milestone, one at a time, requiring the specified test red and restored green. The table has **32** concrete mutation cases (design MU1–MU26, with three MU7, four MU9 and two MU15 branch cases). Record actual production LOC per landing, all gate rc values, the N=10,000 final-code timings, the recorder-walk time, and every mutation kill. Run final vet, compile, race, full `go test -count=1 ./...` and `./scripts/verify_ail.sh` rc 0. If sandbox sockets are denied, mark that run **UNINFORMATIVE UNDER SANDBOX** and use socket-capable CI for the gate; do not relabel it red or green.

No `/v1` route, CLI command, AILANG sketch change, schema-version bump, new duration, payload discovery, registry/journal/interface inbound relation, or commit membership is in scope. Row 104's `committedBy` line remains byte-identical. The reverse result is relationship evidence and is not a causality verdict. Provisioning beyond the existing deadline fails writable Open loudly; V17 measured 7.635 s at 1M, while larger stores need the pending row-23 B1 startup owner or offline provisioning. Read-only handles cache absence until reopened. New refs inserted before a cursor require restarting the walk.

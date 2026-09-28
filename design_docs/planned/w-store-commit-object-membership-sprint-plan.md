# Sprint plan: w-store-commit-object-membership (iteration 204, resumed 203, row 104)

**Base:** detached `07d90e0`. **Authority:** `w-store-commit-object-membership.md` §§3–7, 9–11, 13 in full. Quorum is closed: r1 blocked 2/2 (astra's integrity fix applied; gemini's V13 claim refuted); r2 blocked 2/2 and controller carve-out (gemini V17 verbatim, astra's bounded-context rebase contract applied, landing dependency declined). Do not reopen it. The designer's `~/.ailang/state/world-iter203/prototype_r2.diff` plus `prototype/object_commits.go` and `prototype/object_commits_test.go` are the M1/M2 reference, not changes already present in this worktree. M3 is not prototyped.

## Gates and measured base

Before **every** Go or AILANG command, including mutation and focused runs:

```sh
export PATH=/opt/homebrew/bin:$PATH AILANG_BIN=$HOME/.pinned-ailang/ailang
export GOCACHE=/Users/voightkampff/Library/Caches/go-build TMPDIR=/private/tmp
```

| Untouched `07d90e0` gate | Controller measurement outside sandbox | Interpretation |
|---|---|---|
| `go vet ./...` | rc 0 | Clean. |
| `go test ./... -run '^$'` | rc 0 | Production and `_test.go` compile fence. |
| `go test -race -count=1 ./host/store ./host/daemon ./host/workbench` | rc 0; store 16.466 s, daemon 16.738 s, workbench 1.324 s | Socket-capable result; CI repeats it. |
| `./scripts/verify_ail.sh` | rc 0 | AILANG gate. |

These are controller-provided measurements, not a planner rerun. A local `bind: operation not permitted` from socket/httptest tests is **UNINFORMATIVE UNDER SANDBOX**, never pass/fail. At each landing run the focused command below, then vet, compile fence, and the three-package race gate; all must be rc 0 in socket-capable CI. At final closure also require `go test -count=1 ./...` and `./scripts/verify_ail.sh` rc 0. Keep servers, if any, foreground and bounded; the recorder walk needs no listener. Do not use a PID-file kill, `kill -1/0`, or `pkill`.

## Landing boundaries

| Landing | Added production LOC estimate | Complete, compiling boundary | Tests first introduced here |
|---|---:|---|---|
| M1 | 45 | v4 table, pre-DDL integrity refusal, same-transaction Step 4b, authored ledger; daemon still shows old named stop. | Direct-SQL membership write/rollback/FK/replay/conflict, v3/v4 version and dropped-table tests only. |
| M2 | 64 | Bounded `ObjectCommits` read; M1 tests may now use its helper. | Read ordering, validation, read-only parity, missing-table and opt-in timing tests. |
| M3a | 35 | `readStore` seam and independently rendered optional `CommitView`; daemon supplies nil so the existing page and pinned tests stay green. | Renderer tests and seam/route compile tests only. |
| M3b | 125 | Checked page projection, cursor grammar, pagination, deliberate pinned-test updates, recorder walk. | Page, error, paging, order and three-question walk tests. |

Count added production Go/SQL lines, excluding tests/docs, with `git diff --numstat` per landing. If M3b approaches more than ~150 added production LOC, split it into first-page M3b and continuation/walk M3c; each must have its own passing tests and three gates. Never put later tests in an earlier cumulative snapshot. In particular the prototype's `object_commits_test.go` must be partitioned: M1 uses direct `SELECT` through `s.db` for its checks, and must not declare or call `ObjectCommits`/`MaxObjectCommitPage`; M2 introduces the prototype read helper and its tests, then may convert M1 assertions to it. M1's dropped-table test checks the pre-drop membership with direct SQL. This makes both boundaries compile and pass.

## Implementation contracts used by the mutations

Each mutation is a **one-at-a-time exact literal** `old` → `new` replacement in its named file. The executor must assert exactly one old-literal occurrence, run the named behavioral test to red, restore original bytes, verify green, then move on. Where a short fragment could occur twice, use the whole statement or a named constant. These source spellings are contracts for new code, not assertions about the base. A comment-only edit must stay green. MUT-5 is a guard-bypass proxy for the design's impossible pre-guard membership move: the FK prevents a pre-log membership row, so the original mutation cannot isolate AC5. This proxy directly tests the stale-head nonwrite outcome.

M1 adopts the prototype's `commit_objects` DDL exactly: `(object_ref, entry_index)` primary key, both FKs, `WITHOUT ROWID`; `currentSchemaVersion = 4`; `frozenFutureSchemaVersion = 5`; authored `const schemaV4SQL = schemaV3SQL +` and canonical DDL. In both open modes, current-version `enforceSchemaVersion` checks `sqlite_master` for exactly one `commit_objects` table **before** writable `db.Exec(schemaSQL)`; only fresh initialization creates it. It returns `*SchemaIntegrityError` and leaves a dropped table absent. v3 opens are `LegacySchemaVersionError`, unmodified. `Store.Commit` retains its signature and DR-1 guard; Step 4b is a plain `tx.Exec` after log insert and before head advance, inserting `(o.Hash.String(), h.EntryIndex)` for each `c.Objects`, with `OR IGNORE` and the transaction's deferred rollback. `h := c.Entry.Header` is the **new** entry's header; `[1 3]` pins it. Do not add a `/v1` route or edit `.ail`. When row 23/PR #153 lands later, rebase Step 4b to `tx.ExecContext(ctx, …)` in that bounded transaction; deadline exhaustion must leave no entry, membership, or head advancement. That later merge is not a prerequisite here.

M2 copies the prototype's `object_commits.go`: `MaxObjectCommitPage = 500`, `InvalidObjectCommitCursorError`, `objectCommitsSQL`, `ObjectCommits(ctx, ref, afterEntry, limit)`. Validate `limit < 1 || limit > MaxObjectCommitPage`, malformed ref, and `afterEntry < -1`; `-1` includes genesis entry 0. Query with the caller `ctx`, ascending strict-after keyset, `LIMIT ?`, reject negative stored indices and return errors instead of partial/empty success. A missing table on an already-open handle is a store error. The timing fixture remains `MEMBERSHIP_TIMING=1` opt-in: N=10,000 file-backed production commits, three new objects per commit plus hot object every 100th, 200 reads/sample; record read p50 and commit overhead against a base archive. Do not make latency a regular-suite assertion.

M3a adds `ObjectCommits(context.Context, hashref.HashRef, int64, int) ([]int64, error)` to `host/daemon/daemon.go`'s `readStore`; update `blockingStore`, `recordingStore`, and `failingStore` in `read_deadline_test.go` in the **same** landing, preserving their deadline behavior. Embedded-interface fakes need no method edits. Add `ObjectView.Commits *CommitView` with `Edges`, `Truncated`, `Continued`, `NextHref`, `FirstHref`. Render this optional section between `.Edges` (interface) and `.References`; nil preserves current pages. The exact section literals are `<h3>committedBy</h3>`, `Commits whose object set carried this object, oldest first. An object stored by PutObject or the journal before a commit carried it is attributed only to the commits that carried it.`, CB-NONE `no commit carried this object: it was stored outside any commit (PutObject or journal). Entries that only reference it are listed under referencedBy.`, CB-END `no further commits carried this object`, CB-MORE `Showing 100 commits; more recorded`, and link labels `next commits` / `first page`. The existing edge template produces available `<p>committedBy: <a …>entry I</a></p>` and CB-GONE `UNAVAILABLE: log entry I is not stored`. Its direct renderer tests cover available, empty, continued-empty, truncated and HTML-safe links. `objectCommittedByMissing` and the current stop remain until M3b so all base tests pass at M3a.

M3b removes `objectCommittedByMissing` and the stop from `objectEdges`, leaving interface as its single edge; updates its outdated comment. `committedByEdges(ctx, indexes)` calls `d.reads.GetLogEntry(ctx, index)` per displayed index. On error propagate to `writeWorkbenchStoreError`; on not-found return `{Relation: "committedBy", Target: fmt.Sprintf("entry %d", index), Missing: fmt.Sprintf("log entry %d is not stored", index)}`; on found link `pageHref(from,index)` with `from = min(index, math.MaxInt64-WorkbenchPageLimit)`. The target and label are `entry I`, and the page stays 200 for CB-GONE. Use the same `d.readCtx(r)` for lookup and checks; deadline maps to 503 `Timeout`, other store errors to sanitized 500. This check establishes target existence, not membership cross-verification.

Add `commitsAfter` to accepted query keys and parse a **non-negative** int64 only with `object`, optionally `payload` and `refsAfter`; reject duplicates and world/from/entry combinations. Malformed cursor gets 400 `malformedWorkbenchEntryMessage`. For an absent cursor use `afterEntry=-1`. Preserve both cursors and payload in page links via `url.Values`; `next commits` uses index `commits[WorkbenchPageLimit-1]`, not the hidden lookahead. Query `d.reads.ObjectCommits(ctx, ref, afterEntry, WorkbenchPageLimit+1)` and define truncation as `len(commits) > WorkbenchPageLimit`; display the first 100, with a first-page link on continuation. It may coexist with `refsAfter` so paging one section does not reset the other. New entries after a cursor do not invalidate keyset order; an earlier insertion requires restarting the walk. Preserve the nine frozen `/v1` routes.

Update `workbench_test.go`'s `named-stops` at ~857 to assert the **CB-NONE section** on its uncommitted plain/typed objects; its `edge-order` at ~865–880 must use the unique section heading (`<h3>committedBy</h3>`) between the single interface edge and `<h3>referencedBy</h3>`, not count `<p>committedBy:` once (there may be 100 edges). Update `workbench_references_test.go` ~364 to log the new committedBy answer rather than the deleted stop. Add page tests with a fake read seam for CB-GONE and read errors, 0/100/101 carriers, cursor grammar, payload/refsAfter coexistence, and link GETs via `Handler` + `httptest.NewRecorder`. The real-question walk must ask at least three explicit questions using row 103's fixture shape: which entry carried a transition body, which entries carried an identical state object, and which entry carried the interpreter object. Append the same state object to all three fixture commits so the second answer is `[0 1 2]`; append the interpreter to the first. For each question, locate by semantic ID or known fixture ref, GET object, parse/follow `committedBy`, require selected entry naming the index and its `transitionRef`, follow that object/payload, then log question, answer, refs, statuses, elapsed. Use no listening server and no timing assertion on this walk.

## M1 — schema v4 and atomic membership write

**Files:** `host/store/schema.sql`, `host/store/store.go`, `host/store/schema_version_test.go`, `host/store/journal_test.go`, new `host/store/object_commits_test.go` (M1 tests/helpers only). **Acceptance** (after the common environment exports):

```sh
go test -race -count=1 ./host/store -run '^(TestCommitMembershipRollsBackWithCommit|TestCommitMembershipForeignKeys|TestCommitMembershipDedupesWithinCommit|TestCommitReplayDoesNotDuplicateMembership|TestConflictWritesNoMembership|TestVersionThreeStoreIsRefusedUnmodified|TestDroppedMembershipTableIsRefusedNotRecreated|TestSchemaVersionLedgerIsIndependent|TestOpenReadOnlyEnforcesVersionFour|TestFreshWriterInitializesSchemaAndVersionFour)$'
go vet ./...
go test ./... -run '^$'
go test -race -count=1 ./host/store ./host/daemon ./host/workbench
```

Use the prototype's actual renamed version tests when transcribing; the focused regex must be checked against `go test -list` before landing so it does not silently omit a renamed test. M1 tests use direct SQL membership reads; AC6/7/9/10/11 and timing begin only at M2.

| Mutation | Named file: unique exact `old` → `new` | Must red |
|---|---|---|
| MUT-1 | `host/store/store.go`: `` `INSERT OR IGNORE INTO commit_objects (object_ref, entry_index) VALUES (?, ?);` `` → `` `SELECT ?, ?;` `` | `TestCommitMembershipDedupesWithinCommit` (direct SQL expects row). |
| MUT-2 | `host/store/store.go`: `defer func() { _ = tx.Rollback() }()` → `defer func() { _ = tx.Commit() }()` | `TestCommitMembershipRollsBackWithCommit/later-step-fails` (partial entry/membership must not survive). |
| MUT-3 | `host/store/store.go`: `` `INSERT OR IGNORE INTO commit_objects (object_ref, entry_index) VALUES (?, ?);` `` → `` `INSERT INTO commit_objects (object_ref, entry_index) VALUES (?, ?);` `` | `TestCommitMembershipDedupesWithinCommit`. |
| MUT-4 | `host/store/store.go`: `if resolved {` → `if resolved && false {` | `TestCommitReplayDoesNotDuplicateMembership` (the stale-head guard rejects replay). |
| MUT-5 | `host/store/store.go`: `if selected.String() != c.ObservedHead.String() {` → `if false {` | `TestConflictWritesNoMembership` (stale request must not write). |
| MUT-10 | `host/store/store.go`: `const currentSchemaVersion = 4` → `const currentSchemaVersion = 3` | `TestVersionThreeStoreIsRefusedUnmodified`. |
| MUT-11 | `host/store/schema.sql`: `CREATE TABLE IF NOT EXISTS commit_objects (` → `CREATE TABLE IF NOT EXISTS commit_objects_disabled (` | `TestFreshWriterInitializesSchemaAndVersionFour` and canonical DDL gate. |
| MUT-12 | `host/store/schema.sql`: `object_ref  TEXT    NOT NULL REFERENCES objects(hash_ref),\n    entry_index INTEGER NOT NULL REFERENCES log_entries(entry_index),` → same two lines without `REFERENCES …` | `TestCommitMembershipForeignKeys` (also ledger). |
| MUT-18 | `host/store/store.go`: `if membershipTables != 1 {` → `if false {` | `TestDroppedMembershipTableIsRefusedNotRecreated` (writable Open recreates). |
| MUT-19 | `host/store/store.go`: `if membershipTables != 1 {` → `if membershipTables != 1 && applySchema {` | `TestDroppedMembershipTableIsRefusedNotRecreated` (read-only half). |

MUT-1 and MUT-3 each use the same unique old literal, separately from a restored source. For MUT-12 the old is one two-line span, not either repeated `REFERENCES` fragment. Preserve exact whitespace from the prototype.

## M2 — bounded store read

**Files:** new `host/store/object_commits.go`, `host/store/object_commits_test.go` (add remaining prototype tests; optionally replace M1 direct-SQL helper). **Acceptance:**

```sh
go test -race -count=1 ./host/store -run '^(TestCommitRecordsMembershipPerCarryingCommit|TestObjectCommitsOrderAndContinuation|TestObjectCommitsValidation|TestObjectCommitsReadOnlyMatchesWriterAndEntriesExist|TestObjectCommitsMissingTableIsAnError|TestCommitMembershipRollsBackWithCommit|TestCommitMembershipForeignKeys|TestCommitMembershipDedupesWithinCommit|TestCommitReplayDoesNotDuplicateMembership|TestConflictWritesNoMembership|TestDroppedMembershipTableIsRefusedNotRecreated)$'
go vet ./...
go test ./... -run '^$'
go test -race -count=1 ./host/store ./host/daemon ./host/workbench
MEMBERSHIP_TIMING=1 go test -count=1 -timeout 180s ./host/store -run '^TestObjectCommitsTimingAt10k$' -v
```

The last command is opt-in measurement. Compare its write p50 with an interleaved untouched-base archive run on the same machine. Record AC13 p50 <1 ms and overhead <25% as a measurement gate; if sandbox prevents representative timings, name that limitation and rely on the designer's §7 runs until CI measures final code.

| Mutation | Named file: unique exact `old` → `new` | Must red |
|---|---|---|
| MUT-6 | `host/store/object_commits.go`: `ORDER BY entry_index LIMIT ?` → `ORDER BY entry_index DESC LIMIT ?` | `TestObjectCommitsOrderAndContinuation`. |
| MUT-7 | `host/store/object_commits.go`: `entry_index > ? ORDER BY entry_index` → `entry_index >= ? ORDER BY entry_index` | `TestObjectCommitsOrderAndContinuation` (bounded continuation loop). |
| MUT-8 | `host/store/object_commits.go`: `limit < 1 || limit > MaxObjectCommitPage` → `limit < 1` | `TestObjectCommitsValidation`. |
| MUT-9 | `host/store/object_commits.go`: `afterEntry < -1` → `false` | `TestObjectCommitsValidation`. |

## M3a — read seam and optional renderer

**Files:** `host/daemon/daemon.go`, `host/daemon/read_deadline_test.go`, `host/workbench/render.go`, `host/workbench/render_test.go`. **Acceptance:**

```sh
go test -race -count=1 ./host/daemon ./host/workbench -run '^(TestCommitView|TestWorkbenchObjectProvenanceWalk|TestFrozenV1RouteTableMatchesMux)$'
go vet ./...
go test ./... -run '^$'
go test -race -count=1 ./host/store ./host/daemon ./host/workbench
```

No §6 mutation belongs to this seam-only boundary. Direct `TestCommitView` must prove §3.7 text and link escaping; it passes while the daemon supplies `Commits=nil`.

## M3b — checked object projection and recorder walk

**Files:** `host/daemon/workbench.go`, `host/daemon/workbench_test.go`, `host/daemon/workbench_references_test.go`, new `host/daemon/workbench_commits_test.go`, `host/workbench/render.go`, `host/workbench/render_test.go`. **Acceptance:**

```sh
go test -race -count=1 ./host/daemon ./host/workbench -run '^(TestWorkbenchCommittedByEdges|TestWorkbenchCommitPaging|TestWorkbenchCommitGrammar|TestWorkbenchCommitWalk|TestWorkbenchObjectProvenanceWalk|TestWorkbenchReferenceWalk|TestCommitView|TestFrozenV1RouteTableMatchesMux)$' -v
go vet ./...
go test ./... -run '^$'
go test -race -count=1 ./host/store ./host/daemon ./host/workbench
```

The implementation must keep these M3b old literals unique in their named files: `if !entryOK {` within the new `committedByEdges` helper (if row 103 already uses it in the same file, use a unique whole-line `if !commitEntryOK {` instead and mutate that); `if err != nil {` is too common, so use a dedicated helper `if commitReadErr != nil {`; `len(commits) > WorkbenchPageLimit`; and `afterEntry := int64(-1)` followed by a uniquely named `commitsAfter` assignment. The table below uses the unique committedBy spellings.

| Mutation | Named file: unique exact `old` → `new` | Must red |
|---|---|---|
| MUT-13 | `host/daemon/workbench.go`: `if !commitEntryOK {` → `if false {` | `TestWorkbenchCommittedByEdges/missing-entry` (CB-GONE, no href). |
| MUT-14 | `host/workbench/render.go`: `no commit carried this object: it was stored outside any commit (PutObject or journal). Entries that only reference it are listed under referencedBy.` → `UNAVAILABLE: commit membership unknown` | `TestCommitView/none` and `TestWorkbenchCommittedByEdges/none`. |
| MUT-15 | `host/daemon/workbench.go`: `len(commits) > WorkbenchPageLimit` → `len(commits) >= WorkbenchPageLimit` | `TestWorkbenchCommitPaging/100` (no CB-MORE). |
| MUT-16 | `host/daemon/workbench.go`: `if commitReadErr != nil {` → `if false {` | `TestWorkbenchCommittedByEdges/read-error` (must be sanitized 500). |
| MUT-17 | `host/daemon/workbench.go`: `afterEntry = parsedCommitsAfter` → `afterEntry = -1` | `TestWorkbenchCommitPaging/continuation` (must not repeat first page). |

For MUT-16 the helper must declare `commitReadErr` and use it in the return branch so `if false` still compiles and the subsequent not-found path is exercised; the test checks status, not an implementation string. For MUT-17 `parsedCommitsAfter` is validated before assignment. If a compile error rather than behavioral red arises, refine the replacement/test contract before landing. MUT-13–17 are pre-registered, not measured kills.

## Closure, evidence and scope

The designer ran **13/13 kills** on the prototype's mutation forms: MUT-1–4, 6–12, 18–19. This plan's MUT-2 is a single-literal rollback proxy, so it needs a new kill; MUT-5 was deliberately not run in the prototype and now has the guard-bypass test above. M3's MUT-13–17 await implementation. Total planned: **19** one-at-a-time mutations. Record each actual kill, positive restoration, production LOC, gate rc and CI environment. The prototype's prior N=10,000 measurements (§7) were single read p50 12.0/12.6 µs, hot-100 p50 30.3/31.1 µs, none p50 11.7/12.1 µs; write overhead +16.2/+19.0%. These are historical prototype measurements, not final landed-code results. The controller's base gates above and designer's measurements are distinct evidence.

At final closure require all AC1–AC16, exact §3.7 text, 100/101 continuation, ≥3 recorder questions with statuses/refs/timings, frozen route-table test, full Go/race/vet/fence gates, and `verify_ail.sh` rc 0. No `/v1`, `.ail`, migration/backfill, or public CLI change. Preserve §11 residuals: real v3 history needs a separate upgrade row; this change does not repair other pre-existing silent table recreation; row 23 owns Commit deadline policy; row 111 owns lock-blocked read status. Recheck §10 conflict surface before landing: `schema.sql`, `store.go`, version tests, daemon read seam/workbench, and renderer. Do not make this landing depend on PR #153 or reopen `D-WORLD-40`.

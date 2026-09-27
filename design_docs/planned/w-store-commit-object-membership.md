# w-store-commit-object-membership — record which commits carried an object, and show it as checked committedBy edges

Status: PLANNED (r2 + controller carve-out after quorum r2) — designed iteration 203 (designer `claude-opus-5-5`), **revision r2**. Queue row
104. Quorum r1 was BLOCKED 2/2; the §13 Quorum log records each objection and how it was disposed
of. The prototype (store layer only: M1 + M2) is in this worktree; §8 lists every file it changed.
Item: row 104 `w-store-commit-object-membership` — record commit membership (a schema change) and
show it on the object page as a checked edge, in place of the `committedBy` named stop.
Clauses: clause 5 (a provenance walk answers real "why did X happen" questions in ≤5 minutes).
Estimate: ~1 day across three milestones (M1 ~60 production lines, M2 ~70, M3 ~110–150). If M3's
render/seam churn goes over 150 production lines, split the continuation link into M3b rather than
dropping mutations from the §6 table.
Verified against: dev `4431ed6` (V0), pinned `ailang` at `$HOME/.pinned-ailang/ailang`, Go tests with
`AILANG_BIN` exported. Every present-tense claim about the codebase cites a row in §9.

## 1. Problem

The object inspector's `committedBy` relation is a named stop. The store can't say which commit put
an object there (V1):

* `host/daemon/workbench.go:43` defines `objectCommittedByMissing = "the store records no
  commit-to-object relation, and the provenance field is a free-text label, not a reference"`, and
  `:198` renders `{Relation: "committedBy", Missing: objectCommittedByMissing}`. Two tests pin it:
  `workbench_test.go:857` (the exact stop HTML) and `workbench_references_test.go:364` (the row-103
  walk logs it as a remaining stop).
* `objects` has five columns: `hash_ref, interface_hash_ref, semantic_id, provenance, payload`. None
  of them points at a commit or entry (V2).
* Objects are written by `INSERT OR IGNORE` at three sites (V3):
  `Store.Commit` (`store.go:982`, inside the commit transaction, after the compare-and-append guard),
  `Store.PutObject` (`store.go:494`, outside any commit) and `insertJournalObjectTx`
  (`journal.go:401`, journal intents and outcomes). This means:
  1. content addressing merges duplicate object rows, but it doesn't merge the commits that carried
     them. One object can be in the `Objects` set of **several** commits.
  2. an object can exist with **no** commit at all: registry, broker, transitionreg and world-publish
     objects all go through `PutObject` (V4).
  3. an object can reach `Commit` after `PutObject` already stored it. The `INSERT OR IGNORE` is then
     a no-op, so "the commit that *inserted* the row" is not even a defined event for that object.

For clause 5, the walk "why did this output exist?" stops at the object page. The walker has to fall
back to grep-style archaeology over `log_entries` and payloads. That's exactly the pre-World method
the clause says a walk must replace.

## 2. Measured current state that decides the design

**Callers of `Store.Commit` and what they pass as `c.Objects`** (V4):

| Caller | `c.Objects` | Derivable later from stored rows? |
|---|---|---|
| `host/coordinator/coordinator.go:285` via `plan.go:99-104` | `{input, output, record}` | Partly. `record` = `log_entries.transition_ref`; `output` = the committed world's `state_root`; `input` only appears inside the record's JSON payload (`Input` field). This needs a payload parse keyed on a semantic ID. |
| `host/daemon/handlers.go:669` (`POST /v1/commit`, `decodeCommit` at `:741`) | whatever the HTTP caller sends | **No.** The object list isn't recorded anywhere except as the `INSERT OR IGNORE`d rows, and those can't be told apart from `PutObject` rows. |

So the membership of historical commits **can't be derived** in general. A backfill would be a guess
for every `/v1/commit` entry and a payload-format-coupled heuristic for coordinator entries. This
design doesn't backfill.

**Schema versioning precedent** (V5): both earlier table additions bumped `user_version` and refused
the older store loudly, with no in-place migration: `1856bfb` ("schema 1→2 … a fail-loud un-upgraded
store", `approval_claims`) and `a036062` (row 39, 2→3, `session_credentials`). The version test's
ledger gate, `TestSchemaVersionLedgerIsIndependent`, compares the **table** DDL of a fresh store
against an authored `schemaV3SQL` constant. So adding a table without a bump is exactly what that
gate exists to catch. Rows 96/103 added **indexes**, which are derived and maintained by SQLite on
every write, including writes by an older binary. A membership table has neither property.

**The rig's real store** (V6): `~/.ailang/world/world.db` is at `user_version = 2` with 0 log entries,
37 objects and 0 worlds. The current v3 binary already refuses it as legacy. No store in the rig
holds v3 commit history that a v4 bump would strand. The only v3 databases are regenerable scratch
stores under `~/.ailang/state/world-iter*/`.

## 3. Design

### 3.1 Storage: one table, schema version 3 → 4

```sql
-- Commit membership: one row per distinct object in a commit's Objects set.
-- Written only by Store.Commit, in the commit's own transaction. A (object,
-- entry) row means "entry_index's commit carried this object", not "this commit
-- first inserted the objects row" (PutObject may have stored it earlier).
CREATE TABLE IF NOT EXISTS commit_objects (
    object_ref  TEXT    NOT NULL REFERENCES objects(hash_ref),
    entry_index INTEGER NOT NULL REFERENCES log_entries(entry_index),
    PRIMARY KEY (object_ref, entry_index)
) WITHOUT ROWID;
```

* The key is exactly the read's access path: `WHERE object_ref = ? AND entry_index > ? ORDER BY
  entry_index LIMIT ?` is a range scan on the primary key. `WITHOUT ROWID` makes the table its own
  covering index, so there's no second index to provision or verify. This avoids row 103's
  verifier and 503 guard, because there is no optional structure that can be absent on an openable
  store.
* The foreign keys make "a membership row pointing at nothing" impossible to write: the store
  enables `PRAGMA foreign_keys = ON` in `openSQLite` (V17), and `Commit` inserts the object and entry rows
  before the membership rows. The page still existence-checks each edge (§3.6). The FKs are a write
  fence, not a read proof.
* **Version bump 3 → 4.** A v3 store opened by the new binary is refused **unmodified**, by both the
  writer and the reader, through the existing `LegacySchemaVersionError` path ("store: schema version
  legacy: … has user_version 3 with application schema; binary requires 4; refusing to modify"). A
  v4 store opened by an old binary is refused as future (`FutureSchemaVersionError`). An old binary
  therefore can't append entries that lack membership rows, which is the one way an honest-looking
  empty result could turn into a lie.
* **Integrity check before any DDL (r2, astra).** Writable Open on an existing current-version
  store re-runs `db.Exec(schemaSQL)` (base `store.go:385-387`). Without a guard,
  `CREATE TABLE IF NOT EXISTS commit_objects` would recreate a dropped table **empty**, and every
  historical carrier would then render CB-NONE. So `enforceSchemaVersion` checks, for
  `user_version == 4` in **both** open modes, that `sqlite_master` holds exactly one table named
  `commit_objects`. It does this **before** the `schemaSQL` re-apply (`store.go:394-405`, V14). If
  the table is absent, `Open` and `OpenReadOnly` both return `*SchemaIntegrityError` ("store: schema
  integrity: … is missing table commit_objects; refusing to recreate it empty (its rows cannot be
  re-derived)") and leave the database unmodified. **Only fresh initialization (`freshInitTx`, for a
  store at version 0 with no application objects) may create the table.** Schema re-application on
  writable Open does **not** guarantee that historical membership is complete. What does is this
  check combined with the same-transaction write (§3.2). AC16 pins it.
  The same silent-recreate property exists today for **every** other table on writable Open. It's
  measured for `worlds` on both base v3 and this v4 prototype (V15). It is pre-existing, not
  introduced by this row, and not fixed here (§11).
* Pins: `expectedCurrentSchemaVersion` 3 → 4, `frozenFutureSchemaVersion` 4 → 5. The ledger gets an
  authored `const schemaV4SQL = schemaV3SQL + …` and the ledger test moves to V4. The v3 "supported"
  test becomes a "v3 refused unmodified" test. This mirrors how row 39 handled v2.

**Rejected (b): a table provisioned on writable Open without a bump (the row 96/103 shape).** It
fails for three reasons:
1. It reds the ledger gate, and passing it would mean weakening that gate.
2. An older v3 binary could still open the store and commit without writing membership. Those
   entries would then read as "carried nothing", which is exactly the silent empty that must never
   happen.
3. It would need a "membership unknown before entry K" boundary for stores that have no history to
   protect (V6).

**Rejected (a′): bump with an in-place v3 → v4 migration and a boundary row.** This would be the
store's first migration. It would add a boundary table, a read-only "unavailable" branch and a third
rendered state, all to protect v3 history that doesn't exist in the rig (V6). If an operator later
has real v3 history to carry forward, that is a store-upgrade tool: **NEW ROW NEEDED** (§11).

### 3.2 Write path inside `Store.Commit`

This is a new Step 4b: after the log row is inserted and before the selected head advances, in the
same `tx`:

```go
for _, o := range c.Objects {
    INSERT OR IGNORE INTO commit_objects(object_ref, entry_index) VALUES (o.Hash, h.EntryIndex)
}
```

`h` here is `c.Entry.Header`, **the new entry's header** (`store.go:1021`). This is the same
`h.EntryIndex` that Step 4 has just inserted as `log_entries.entry_index` for this commit
(`:1023-1027`). It is not the parent head's index, so each membership row names the carrying entry
itself. AC6 (`[1 3]`) pins this attribution (V13).

* **Same transaction.** Any later failure (head advance, outcome append, `tx.Commit`) rolls the
  membership rows back together with the entry. AC1 proves this with an injected failure after
  Step 4b.
* **Duplicates within one commit.** `OR IGNORE` on the primary key collapses one hash listed twice
  into a single row. This is a real case: when a coordinator transition's output bytes equal its
  input bytes, `plan.go` builds `in` and `out` with the same `Hash` (V4). AC3 covers it.
* **Idempotent replay.** A replay whose `bindCommitIntentTx` resolves returns `nil` at Step 1, before
  Step 2. It writes nothing, membership included. The primary key backs this up: a second row for the
  same `(object, entry)` can't exist. AC4 replays a committed invocation and asserts the row count
  and the rows are unchanged.
* **Conflict.** `ConflictError` returns at Step 1, before any insert. AC5 asserts no membership rows
  are written.
* Signature, the unbounded `context.Background()` guard read, and DR-1 are all unchanged. The
  statement uses `tx.Exec`, like every other Commit step (row 23 / PR #153 owns deadline policy).
* **Bounded-context rebase contract (quorum r2, gpt6-astra's fix, applied verbatim where it does not
  override an open human decision).** "Step 4b must use that bounded transaction context, and
  deadline exhaustion must return an explicit error without a committed entry, membership rows, or
  head advancement." Concretely: when row 23's `CommitContext` lands (PR #153 converts every Commit
  step to `tx.ExecContext(ctx, …)`, V18), Step 4b becomes `tx.ExecContext(ctx, …)` in the same
  rebase, and the same-transaction rollback (AC1) guarantees the "no entry, no membership, no head"
  outcome on exhaustion. The **landing dependency** in that fix ("M1 must not land before row 23 /
  PR #153") is **not** applied: the unbounded Commit is pre-existing at base (DR-1, V18), Step 4b
  adds no new wait point (it writes inside a transaction that already holds SQLite's write lock from
  Step 2, V18; +0.2 ms, §7), and PR #153's merge is held on the open human decision `D-WORLD-40`.
  Making this row wait on it would re-decide that ledger row by the back door. Gate 2 routes a
  pre-existing defect to its owner, which here is row 23.

### 3.3 The three cases, defined

For an object `X`, `committedBy(X)` = the set of log entries whose commit had `X` in its `Objects`
set, in ascending `entry_index` order.

| Case | Store state | Page shows |
|---|---|---|
| Carried by exactly one commit | 1 row | one `committedBy` edge → entry i |
| Carried by several commits | k rows | k edges, ascending, capped at 100 per page, with continuation |
| Carried by none (`PutObject` / journal only) | 0 rows | the definite statement `CB-NONE` (§3.7), not UNAVAILABLE |

**Why every carrying commit rather than only the first:** in a "why did X happen" walk, X is usually
an output or state object. Identical bytes produced by two different invocations dedupe to one
object row. If only the first carrier were recorded, the walk would attribute the second
occurrence's cause to the first invocation. That would be a wrong answer presented as a verified
one. With every carrier recorded in ascending order, the first edge is still "the earliest commit
that carried it", so nothing is lost. The storage cost is one small row per (object, commit).
**Why "carried", not "inserted":** V3(iii). The row existence of an object that `PutObject` stored
earlier is not attributable to any commit, and the page must not claim otherwise.

### 3.4 Read query, ordering, bound, continuation

```go
// ObjectCommits returns entry indexes whose commit carried ref, ascending, strictly after afterEntry.
func (s *Store) ObjectCommits(ctx context.Context, ref hashref.HashRef, afterEntry int64, limit int) ([]int64, error)
// SQL: SELECT entry_index FROM commit_objects WHERE object_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?
```

* It validates `1 ≤ limit ≤ MaxObjectCommitPage (500)` (`InvalidLimitError`), `ref` parses
  (`InvalidRefError`) and `afterEntry ≥ -1`. The start sentinel is `-1`, because entry 0 is a valid
  genesis index.
* Ordering is the primary-key order `(object_ref, entry_index)`, so it's deterministic and stable
  under appends. Because entry indexes only grow, a cursor never skips a row that existed when the
  page was read.
* The page asks for `WorkbenchPageLimit + 1` (101) rows to detect truncation. Continuation is a new
  workbench query key `commitsAfter=<entry_index>` (M3). It mirrors row 103's `refsAfter`: accepted
  only together with `object`, and it may combine with `payload` and `refsAfter`. It's parsed as a
  non-negative int64, and a malformed value → 400 with `malformedWorkbenchEntryMessage`.
* A negative stored `entry_index` is refused as a store error, not rendered (same as row 103).

### 3.5 Error mapping (row 103 / 35+38 F6 pattern)

| Condition | Result |
|---|---|
| Read deadline exceeded (the handler's existing read `ctx`) | 503 `Timeout` via `writeWorkbenchStoreError` |
| Any other store error (for example, `commit_objects` dropped **while a handle is already open**) | sanitized 500 `Internal`, never an empty edge list |
| `commit_objects` absent when the store is opened | never reaches the page: both open modes refuse with `SchemaIntegrityError` (§3.1, AC16), so the daemon doesn't start |
| `GetLogEntry(i)` for an edge returns not-found | that edge is UNAVAILABLE with `CB-GONE`; the page stays 200 |
| `GetLogEntry(i)` returns an error | 5xx as above |
| v3 store | never reaches the page: `Open`/`OpenReadOnly` refuse it (§3.1), so the daemon doesn't start |

There's no 503 "unavailable" branch like row 103's. A v4 store that opens always has the table:
fresh init creates it, and every later open refuses a store without it rather than recreating it
(§3.1). So a missing table on an already-open handle can only mean store damage, which is a 5xx. Row 111 owns the 500-vs-503 question for lock-blocked reads; this row
uses the same adapter and inherits its answer.

### 3.6 Existence check

Each `committedBy` edge target is `GetLogEntry(i)` under the handler's read `ctx`. If the entry is
stored, the edge is a link `pageHref(from, i)` with `from = min(i, MaxInt64 - WorkbenchPageLimit)`.
That's the same link shape and overflow clamp as row 103's `transitionRef` edges. If it's absent, the
edge is `CB-GONE`. The entry row doesn't list its objects, so beyond existence the check can't
cross-verify membership (unlike row 103's `actual != ref` check). Membership integrity rests on the
same-transaction write and the FKs, which AC1/AC2 prove.

### 3.7 Rendered text (exact)

`objectCommittedByMissing` is deleted. `committedBy` becomes a section in the fixed walk order, like
`referencedBy`: interface → committedBy → referencedBy.

* Section heading: `<h3>committedBy</h3>`, followed by the scope line
  `Commits whose object set carried this object, oldest first. An object stored by PutObject or the
  journal before a commit carried it is attributed only to the commits that carried it.`
* Edge (available): `<p>committedBy: <a href="?from=F&amp;entry=I" …>entry I</a></p>`
* **CB-NONE** (0 rows, first page): `<p>no commit carried this object: it was stored outside any
  commit (PutObject or journal). Entries that only reference it are listed under referencedBy.</p>`
* **CB-END** (0 rows on a continued page): `<p>no further commits carried this object</p>`
* **CB-GONE** (edge target absent): `UNAVAILABLE: log entry I is not stored`
* **CB-MORE** (truncated): `<p>Showing 100 commits; more recorded</p>` plus a `next commits` link, and
  on continued pages a `first page` link, the same as referencedBy.

There's no "unknown" text on an openable store. Two things together guarantee this:
* the bump, which refuses v3 stores (so there is no pre-membership history), and
* the integrity check, which refuses a v4 store whose table was dropped rather than recreating it
  empty (§3.1).

So every store this binary opens has had membership recorded, in-transaction, since genesis. The
two "unknown" situations never reach a page. The operator sees the `LegacySchemaVersionError` or
`SchemaIntegrityError` text at startup instead.

### 3.8 Why no `.ail`, CLI or `/v1` change

This is a host storage fact, and effects stay at the host boundary. No AILANG kernel law changes: a
commit's semantics are the same, and the store only remembers more about it. The frozen `/v1` route
table (`daemon.go`) is untouched, because the surface is the workbench page, outside that table, as
in row 103. The clause-5 walk only needs a browser: locate → object → **committedBy → entry** →
transitionRef → record payload. A JSON/CLI projection isn't essential. If a scripted walk needs it,
that is **NEW ROW NEEDED** (§11).

## 4. Milestones (compile/dependency order; each ≤150 production lines)

* **M1: schema v4 + membership write** (~60 prod lines, prototyped). `schema.sql` table,
  `currentSchemaVersion = 4`, the pre-DDL `commit_objects` integrity check + `SchemaIntegrityError`,
  Commit Step 4b, and the version-test ledger/pin moves. ACs: AC1–AC5, AC8, AC9, AC16.
* **M2: store read** (~70 prod lines, prototyped). `host/store/object_commits.go` `ObjectCommits` +
  `MaxObjectCommitPage` + validation. ACs: AC6, AC7, AC10, AC11, AC13 (store timing).
* **M3: workbench projection** (~110–150 prod lines, not prototyped). `readStore.ObjectCommits`
  (`daemon.go:365`) and a `committedByEdges` existence check in `workbench.go`. It also covers
  `commitsAfter` parsing and continuation, the `ObjectView.Commits` section in
  `host/workbench/render.go`, deleting `objectCommittedByMissing`, and updating the pinned tests
  (`workbench_test.go:857` named-stops and `edge-order`, `workbench_references_test.go:364`). It
  closes with the real-question walk (AC14). ACs: AC12, AC14, AC15.

## 5. Acceptance criteria

* **AC1 same-tx, both directions.** An injected `RAISE(ABORT)` trigger on `commit_objects` (the
  membership write itself fails) makes Commit return an error and leaves **no log entry**. The same
  trigger on `store_heads` (a later step fails) leaves no entry and **no membership rows**.
* **AC2 FK fence.** A membership insert naming a non-existent entry fails inside Commit's
  transaction. (This is structural, via the FK. It's proven by a direct `tx` insert in the test.)
* **AC3 in-commit duplicate.** A commit listing the same hash twice records one row.
* **AC4 replay.** Replaying a committed `InvocationID` returns nil, and the membership rows are
  byte-identical before and after.
* **AC5 conflict.** A stale-head Commit writes no membership.
* **AC6 multi-commit.** An object carried by entries 1 and 3 (not 2) reads back as `[1 3]`.
* **AC7 no-commit.** A `PutObject`-only object reads back as `[]` with a nil error. A
  `PutObject`-then-Commit object reads back as exactly that commit.
* **AC8 v3 refused unmodified.** Writer and reader both return `LegacySchemaVersionError`, and the
  schema names and version are unchanged.
* **AC9 v4 read-only.** `OpenReadOnly` on a v4 store returns the same membership as the writer.
  There's no unavailable branch.
* **AC10 ordering + continuation.** 7 carriers read with `limit=3` give `[a b c] [d e f] [g]` with
  cursors, and their union equals the one-shot read.
* **AC11 bound.** `limit` 0 and 501 are refused (`InvalidLimitError`), `afterEntry = -2` is refused,
  and 500 is accepted.
* **AC12 existence-checked edge.** A stored entry renders a link that GETs 200 and names `entry I`.
  An absent entry (seam fake) renders `CB-GONE`, and a `GetLogEntry` error → 500.
* **AC13 timing.** At N = 10,000 entries, the `ObjectCommits` p50 and per-commit write overhead are
  recorded (§7). The gate is p50 < 1 ms, and write overhead < 25% of the base Commit p50.
* **AC14 real-question walk (M3).** On ≥3 real questions (row 103's walk fixture shape), the walk
  locate → object → committedBy → entry → transitionRef → payload completes. The question, answer,
  refs, statuses and elapsed time are logged.
* **AC15 render cases.** CB-NONE, CB-END, CB-MORE (101 carriers) and fixed walk order are each
  pinned.
* **AC16 dropped table refused, not recreated (r2).** The test commits a carrier (`ObjectCommits` →
  `[1]`), closes, and drops `commit_objects` through a raw connection. It then attempts `Open` and
  `OpenReadOnly` in turn. Each must return `*SchemaIntegrityError{Table: "commit_objects"}`, and
  after each attempt `sqlite_master` must still hold no `commit_objects`.
  (`TestDroppedMembershipTableIsRefusedNotRecreated`.)

## 6. Non-vacuity table (mutation → the AC that turns red)

The M1/M2 rows were **run** against the prototype (§8.2 records the observed red). The M3 rows are
pre-registered for the executor.

| # | Mutation | AC / test | Status |
|---|---|---|---|
| MUT-1 | Step 4b statement → `SELECT ?, ?` (no membership write) | AC6/AC7 `TestCommitRecordsMembershipPerCarryingCommit` | **killed** (§8.2) |
| MUT-2 | Step 4b removed from the tx and re-done by `s.db.Exec` **after** `tx.Commit()` | AC1 `TestCommitMembershipRollsBackWithCommit/membership-fails` | **killed** |
| MUT-3 | `INSERT OR IGNORE` → `INSERT` in Step 4b | AC3 `TestCommitMembershipDedupesWithinCommit` | **killed** |
| MUT-4 | replay path falls through (`return nil` in the resolved branch removed) | AC4 `TestCommitReplayDoesNotDuplicateMembership` (+ existing `TestIntentBindingMirrorsAllTenSketchRows`) | **killed**, but by the ConflictError the fall-through hits, not by the row-count assertion (§8.2 note) |
| MUT-5 | Step 4b moved before the compare-and-append guard | AC5 | **not run**: impossible by construction. Before Step 4 no entry row exists, so the FK refuses every commit. The whole suite would red, with no discrimination for AC5. AC5 stays as a regression pin. |
| MUT-6 | `ORDER BY entry_index` → `ORDER BY entry_index DESC` | AC10 `TestObjectCommitsOrderAndContinuation` (+ AC6, AC9) | **killed** |
| MUT-7 | `entry_index > ?` → `>=` | AC10 (continuation repeats a row) | **killed** |
| MUT-8 | limit upper check removed | AC11 `TestObjectCommitsValidation` | **killed** |
| MUT-9 | `afterEntry < -1` check removed | AC11 | **killed** |
| MUT-10 | `currentSchemaVersion` back to 3 | AC8 `TestVersionThreeStoreIsRefusedUnmodified` + 8 version/ledger tests | **killed** |
| MUT-11 | `commit_objects` renamed away in `schema.sql` | AC11 `TestObjectCommitsMissingTableIsAnError` path + every Commit test + canonical DDL + ledger | **killed** |
| MUT-12 | both FK clauses removed | AC2 `TestCommitMembershipForeignKeys` (+ canonical DDL, ledger) | **killed** |
| MUT-18 (r2) | integrity check disabled (`if membershipTables != 1` → `if false`) | AC16 `TestDroppedMembershipTableIsRefusedNotRecreated` | **killed**: writable Open recreates the table (§8.2) |
| MUT-19 (r2) | check limited to writable opens (`&& applySchema`) | AC16 (read-only half) | **killed** |
| MUT-13 (M3) | `GetLogEntry` not-found rendered as a link | AC12 | pre-registered |
| MUT-14 (M3) | 0 rows rendered as UNAVAILABLE instead of CB-NONE | AC15 | pre-registered |
| MUT-15 (M3) | truncation check `> 100` → `>= 100` | AC15 CB-MORE | pre-registered |
| MUT-16 (M3) | `GetLogEntry` error swallowed as CB-GONE | AC12 (500 expected) | pre-registered |
| MUT-17 (M3) | `commitsAfter` ignored | AC10-at-page continuation | pre-registered |

## 7. Measurements (N = 10,000)

The rig's real store has 0 entries (V6), so scale is synthetic, as in rows 96/103. The fixture is a
file-backed store (`Open(path)`, WAL/default DSN) built by 10,000 production `Commit` calls. Each
commit carries three fresh objects (input/output/record, the coordinator's shape), and every 100th
commit also carries one shared `hot` object. Each read is sampled 200 times. Test:
`TestObjectCommitsTimingAt10k`, opt-in with `MEMBERSHIP_TIMING=1`. The write baseline is the same
commit loop run in a `git archive 4431ed6` copy of the base (no membership). Base and new runs were
interleaved twice on the same machine.

| Measure | Round 1 | Round 2 |
|---|---:|---:|
| `ObjectCommits` p50, object in 1 commit | 12.0 µs | 12.6 µs |
| `ObjectCommits` p50, object in 100 commits (limit 101) | 30.3 µs | 31.1 µs |
| `ObjectCommits` p50, object in 0 commits | 11.7 µs | 12.1 µs |
| Commit p50, base `4431ed6` (no membership) | 1.084 ms | 1.080 ms |
| Commit p50, with membership (3 rows) | 1.260 ms | 1.285 ms |
| Write overhead | **+16.2 %** (+0.18 ms) | **+19.0 %** (+0.20 ms) |

Both AC13 gates hold: p50 < 1 ms, and overhead < 25%. p95 figures are in V11. For comparison, the
read is in the same range as row 103's reverse read (p50 45–86 µs), and far inside the clause-5
five-minute budget.

## 8. Prototype

The prototype covers M1 + M2 only: the store layer. **M3 (workbench projection) is not
prototyped.** The daemon still renders the `committedBy` named stop, and its pinned tests still pass
unchanged.

### 8.1 Manifest (files changed in this worktree)

| File | Kind | Change |
|---|---|---|
| `host/store/schema.sql` | prod | +11: `commit_objects` table (§3.1) |
| `host/store/store.go` | prod | +32/−2: `currentSchemaVersion = 4`; `SchemaIntegrityError` + pre-DDL `commit_objects` check in `enforceSchemaVersion` (r2); Commit Step 4b; doc-comment step list |
| `host/store/object_commits.go` | prod, new | 64: `ObjectCommits`, `MaxObjectCommitPage`, `InvalidObjectCommitCursorError` |
| `host/store/object_commits_test.go` | test, new | AC1–AC11 + AC8 v3 refusal + opt-in AC13 timing + AC16 dropped-table refusal (r2) |
| `host/store/schema_version_test.go` | test | pins 3→4 / 4→5; the authored `schemaV4SQL` ledger + ledger gate moved to v4; v3 now a read-only legacy fixture; supported/fresh tests renamed to Four. **r2:** the two writer-lock-release probes (in `…VersionOne…` and `…LegacyVersionZero…`) used to relabel a v1/v0 fixture as current and reopen it. They now apply the frozen `schemaV4SQL` before relabelling, because the integrity check (correctly) refuses a current-version store without `commit_objects`. The probes still prove only lock release. |
| `host/store/journal_test.go` | test | `canonicalTableDDL` gains `commit_objects` |
| `design_docs/planned/w-store-commit-object-membership.md` | doc | this file |

That's ~107 production lines. Scratch files (the mutation harness, the base copy, logs) live under
`~/.ailang/state/world-iter203/designer-scratch/`, outside the worktree.

### 8.2 Mutations run

Harness: `~/.ailang/state/world-iter203/designer-scratch/mutate.py <MUT>`. For each mutation it
backs up the file, applies one exact-string edit (asserting the edit matched exactly once), runs
`go test -count=1 -timeout 120s ./host/store/`, and restores the file. After all runs, `git status`
showed only the intended changes, and the unmutated package passed again.

| MUT | Result | Failing tests | Observed red (first line) |
|---|---|---|---|
| 1 | rc=1 | PerCarryingCommit, RollsBack, Dedupes, Replay, OrderAndContinuation, ReadOnlyMatchesWriter | `multi: ObjectCommits = [], want [1 3]` |
| 2 | rc=1 | `TestCommitMembershipRollsBackWithCommit` | `entry 1 after failed commit: ok=true err=<nil>` |
| 3 | rc=1 | `TestCommitMembershipDedupesWithinCommit` | `Commit entry 1: … UNIQUE constraint failed: commit_objects.object_ref, commit_objects.entry_index (1555)` |
| 4 | rc=1 | `TestIntentBindingMirrorsAllTenSketchRows`, `TestCommitReplayDoesNotDuplicateMembership` | `replay: store: stale observed world head …` |
| 6 | rc=1 | PerCarryingCommit, OrderAndContinuation, ReadOnlyMatchesWriter | `multi: ObjectCommits = [3 1], want [1 3]` |
| 7 | rc=1 | `TestObjectCommitsOrderAndContinuation` | `pages = [[1 2 4] [4 5 7] [7 8 10] [10] [10]], want [1 2 4 5 7 8 10] in pages of 3,3,1` |
| 8 | rc=1 | `TestObjectCommitsValidation` | `limit 501: want InvalidLimitError, got <nil>` |
| 9 | rc=1 | `TestObjectCommitsValidation` | `afterEntry -2: want InvalidObjectCommitCursorError, got <nil>` |
| 10 | rc=1 | 9 tests incl. `TestVersionThreeStoreIsRefusedUnmodified` | `Open: err = <nil>, want legacy Found=3 Current=4` |
| 11 | rc=1 | 19 tests | `store: commit membership …: SQL logic error: no such table: commit_objects` |
| 12 | rc=1 | `TestCommitMembershipForeignKeys`, canonical DDL, ledger | `membership naming absent entry 999 was accepted` |

| 18 (r2) | rc=1 | `TestDroppedMembershipTableIsRefusedNotRecreated` | `Open: err = <nil>, want SchemaIntegrityError for commit_objects` then `Open recreated commit_objects (count 1)` |
| 19 (r2) | rc=1 | `TestDroppedMembershipTableIsRefusedNotRecreated` | `OpenReadOnly: err = <nil>, want SchemaIntegrityError for commit_objects` |

13 run, 13 killed, 0 survived. MUT-5 was not run, for the reason given in §6. MUT-18's red
reproduces astra's premise directly: without the check, writable Open brought `commit_objects` back.

Notes, recorded plainly:
* **MUT-7, first attempt:** the continuation loop was unbounded, so MUT-7 went red only by a 120 s
  test timeout (a hang). I bounded the loop to 5 rounds, and it now reds in milliseconds with the
  line above.
* **MUT-4:** the fall-through is caught by the compare-and-append guard (ConflictError) before Step
  4b is reached. The replay test's membership-equality assertion is therefore a regression pin, not
  the discriminating check. Membership can't be duplicated on replay for two independent reasons:
  the early return (V10) and the `(object_ref, entry_index)` primary key, which MUT-3 shows is
  enforced.

### 8.3 Timings
See §7. The r2 change adds one `sqlite_master` point query at open time only. Commit and the read
path are unchanged, so the timings weren't re-run.

### 8.4 r2 gate results (run on the final r2 tree)

* `gofmt -l host/`: no output. `go vet ./...`: ok. `go test -run '^$' ./...`: no failures.
* `go test -count=1 ./host/store/`: `ok … 6.294s`. `go test -count=1 -race ./host/store/`:
  `ok … 16.955s`.
* `go test ./...`: 24 `ok`, 0 FAIL, `EXIT 0`. The r1 `host/pkgproj` flake did not recur.
* `./scripts/verify_ail.sh`: rc=0.

## 9. Verification Log

| V | Claim | Command | Observed |
|---|---|---|---|
| V0 | base | `git log --oneline -1` | `4431ed6 record(202): row 103 LANDED …` |
| V1 | the committedBy stop and its pins | `grep -n 'objectCommittedByMissing' host/daemon/*.go` | `workbench.go:43` const; `workbench.go:198` edge; `workbench_references_test.go:364`; `workbench_test.go:857` (4 hits) |
| V2 | `objects` has 5 columns, no commit link | `sed -n 13,19p host/store/schema.sql` | `hash_ref, interface_hash_ref, semantic_id, provenance, payload` |
| V3 | three object insert sites | `grep -n 'INSERT OR IGNORE INTO objects' host/store/*.go` (positive control in the same run: `grep -n 'INSERT OR IGNORE INTO worlds' host/store/*.go` → `store.go:551`, `store.go:994`) | `journal.go:401`, `store.go:494`, `store.go:982` |
| V4 | Commit callers + object sets; PutObject callers | `git grep -n '\.Commit(' -- '*.go' \| grep -v _test \| grep -v tx.Commit`; `sed -n 40,120p host/coordinator/plan.go`; `sed -n 640,760p host/daemon/handlers.go`; `git grep -n 'PutObject(' -- '*.go' \| grep -v _test` | `coordinator.go:285`, `handlers.go:669`; plan `objects := []store.Object{in, out, rec}`, `TransitionRef: rec.Hash`, `StateRoot: out.Hash`, `Hash: hashref.SumSHA256(payload)` (so `in`/`out` collide on equal bytes); decodeCommit copies `r.Objects` verbatim; PutObject callers in world-publish, broker (5), registry, transitionreg |
| V5 | version precedent + ledger gate | `git log --oneline -- host/store/schema.sql`; `grep -n 'SchemaVersion = ' host/store/*.go`; `sed -n 568,640p host/store/schema_version_test.go`; `grep -n -i 'schema version' design_docs/planned/w-session-authority.md` | `1856bfb … schema 1→2, and a fail-loud un-upgraded store`; `a036062 row 39`; `store.go:42 currentSchemaVersion = 3`; test `frozenFutureSchemaVersion = 4`, `expectedCurrentSchemaVersion = 3`; ledger compares `type='table'` DDL (`journal_test.go:969 tableDDL`) to `schemaV3SQL`; row-39 doc §74 "bumps 2 → 3 … an un-migrated v2 DB is refused loudly" |
| V6 | the rig's real store | `sqlite3 -readonly "file:$HOME/.ailang/world/world.db?mode=ro" 'PRAGMA user_version; SELECT count(*) FROM log_entries; SELECT count(*) FROM objects; SELECT count(*) FROM worlds; …'` | `2 0 37 0`; tables `approval_claims epoch_registry_heads journal log_entries objects store_heads verification_cache worlds` (no `session_credentials`, so v2) |
| V7 | row 103 surface shape | `sed -n 95,262p host/daemon/workbench.go`; `sed -n 370,440p …` | `pageHref(from, entry)`; `checkedReferenceEdge` clamps `from` to `MaxInt64-WorkbenchPageLimit`; `ObjectReferences(ctx, ref, refsAfter, WorkbenchPageLimit+1)`; `writeWorkbenchStoreError` → 503 Timeout / 503 ReferenceIndexUnavailable / 500 |
| V8 | readStore has 7 methods | `grep -n 'type readStore interface' -A9 host/daemon/daemon.go` | GetObject, GetWorld, GetLogEntry, GetRegistryHead, SelectedHead, ObjectsBySemanticID, ObjectReferences |
| V9 | the walk-order test pins one `committedBy` | `sed -n 865,880p host/daemon/workbench_test.go` | `edge-order`: `"<p>committedBy: "` must occur exactly once. M3 must change this deliberately. |
| V10 | replay returns before any insert | `sed -n 955,968p host/store/store.go` | `if resolved { // … idempotent no-op  return nil }` precedes Step 1 guard and Step 2 |
| V11 | N=10,000 timings | `bash ~/.ailang/state/world-iter203/designer-scratch/timing.sh` (runs `MEMBERSHIP_TIMING=1 go test -run TestObjectCommitsTimingAt10k -v ./host/store/` in the base copy and the worktree, twice, interleaved) | base r1 `commit … p50=1.084292ms p95=1.459291ms`; new r1 `read single p50=12µs`, `read hot-100 p50=30.292µs p95=33.167µs`, `read none p50=11.709µs`, `commit … p50=1.259666ms p95=1.642458ms`; base r2 `p50=1.080375ms`; new r2 `single 12.583µs`, `hot-100 31.083µs`, `none 12.083µs`, `commit p50=1.284833ms p95=1.682791ms` |
| V12 | full suite + gates on the prototype | `go vet ./...`; `go test -run '^$' ./...`; `gofmt -l host/`; `go test -race ./host/store/`; `go test ./...`; `./scripts/verify_ail.sh` | vet/compile fence/gofmt: no output; race `ok host/store 17.266s` (re-run on the final test file); full suite 23 `ok`, 1 FAIL `host/pkgproj` (`TestQueryInterfaceReturnsWhileADescendantHoldsStdout`, known flake row 117), re-run alone `ok host/pkgproj 11.596s`; `verify_ail.sh` rc=0 |

| V13 | Step 4b's `h` is the new entry's header (refutes gemini r1) | `grep -n 'h := c.Entry.Header\|INSERT INTO log_entries\|h.EntryIndex, c.Entry.EntryHash\|o.Hash.String(), h.EntryIndex' host/store/store.go`; `grep -n '{"multi", shared.Hash' host/store/object_commits_test.go` | `1021: h := c.Entry.Header`; `1023: INSERT INTO log_entries`; `1027: h.EntryIndex, c.Entry.EntryHash.String(), …` (Step 4); `1039: o.Hash.String(), h.EntryIndex,` (Step 4b); test `97: {"multi", shared.Hash, []int64{1, 3}}`. Same statement at base: `git show 4431ed6:host/store/store.go \| grep -n 'h := c.Entry.Header'` → `1002` (the controller's `:1003`/`:1009` were measured on the r1 prototype, which is base +1; the r2 check and error type add 18 lines above them) |
| V14 | open/initialization order after r2 | `grep -n 'func enforceSchemaVersion\|if version == 0 && applicationObjects == 0\|freshInitTx(db\|if version == currentSchemaVersion\|membershipTables\|db.Exec(schemaSQL)' host/store/store.go`; base: `git show 4431ed6:host/store/store.go \| grep -n 'db.Exec(schemaSQL)\|if version == currentSchemaVersion'`; `go test -count=1 -v -run DroppedMembership ./host/store/` | `361` enforceSchemaVersion → `369` fresh branch (`373 freshInitTx` is the only creator) → `394` current-version branch → `397-401` `commit_objects` count check, for both modes → `405 db.Exec(schemaSQL)` (writable only). Base: `385`/`387`, the re-apply with no check. Test: `--- PASS: TestDroppedMembershipTableIsRefusedNotRecreated (0.01s)`; MUT-18/19 red as in §8.2 |
| V15 | writable Open silently recreates any dropped table, pre-existing (scratch stores only) | a probe test `zz_drop_probe_test.go`, copied into `…/designer-scratch/base` (`git archive 4431ed6`) and `…/designer-scratch/newcopy` (an rsync of this worktree), neither inside the worktree: Open → PutWorld → Close → raw `DROP TABLE worlds` → Open writable → count; `go test -count=1 -run TestProbeWritableOpenRecreatesDroppedWorlds -v ./host/store/` | base: `user_version=3 worlds rows before drop=1; after writable reopen: open err=nil, worlds table present=1, rows=0`; newcopy: `user_version=4 … rows before drop=1; … open err=nil, worlds table present=1, rows=0` |
| V17 | `openSQLite` enables foreign keys (gemini r2) | `grep -n 'foreign_keys' host/store/*.go \| grep -v _test`; base: `git show 4431ed6:host/store/store.go \| grep -n 'PRAGMA foreign_keys'` | `store.go:350: if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {` (r2 tree); base `341` (same statement). MUT-12 (§8.2) is the behavioural control: removing the FK clauses lets an absent-entry membership insert through |
| V18 | Commit's unbounded path is pre-existing; Step 4b adds no wait point; PR #153's shape (astra r2) | base: `git show 4431ed6:host/store/store.go \| grep -n 'DR-1\|selectedHeadTx(context.Background\|Step 2: immutable'`; PR #153: `git diff 4431ed6...origin/sprint/w-store-bounded-durable-operations-m1 -- host/store/store.go \| grep -nE '^\+.*(CommitContext\|ExecContext)'` | base `964: // DR-1: Store.Commit keeps its signature and today's unbounded behaviour;`, `966: selectedHeadTx(context.Background(), tx)`, `979: // Step 2: immutable objects.` (the first write, so the deferred tx holds the write lock from here, before Step 4b). PR #153: `+func (s *Store) CommitContext(ctx context.Context, c Commit) (err error) {` and five `+ tx.ExecContext(ctx,` step statements. `D-WORLD-40` OPEN in the charter ledger (PR #153 held) |
| V16 | r2 gates | `gofmt -l host/`; `go vet ./...`; `go test -run '^$' ./...`; `go test -count=1 ./host/store/`; `go test -count=1 -race ./host/store/`; `go test ./...` | see §8.4 |

## 10. Conflict Surface

| File | This row | Open rows touching it |
|---|---|---|
| `host/store/schema.sql` | +table | 99/100 if they add store facts (both are clause-5 walk rows, D-WORLD-42 grooming); any schema change after this one must bump to 5 |
| `host/store/store.go` | version const, Commit Step 4b | **23 / PR #153** (Commit deadline policy, DR-1; Step 4b is a plain `tx.Exec` and rebases trivially onto a ctx-threaded Commit); 111 (lock-blocked read status, via the adapter, not this file) |
| `host/store/schema_version_test.go`, `journal_test.go` fixtures | version pins, ledger v4 | any row that changes schema |
| `host/store/object_commits.go` (new) | read | none |
| `host/daemon/daemon.go` (`readStore` only, not the route table) | +1 method | 106 / 108 if they extend `readStore`; 111 |
| `host/daemon/workbench.go`, `host/workbench/render.go` | M3 | 99, 100 (other object-page relations); 111 (`writeWorkbenchStoreError`) |

## 11. Residuals

* An upgrade tool for stores holding v2/v3 history: **NEW ROW NEEDED**. It isn't required by clause
  5, since no such store exists in the rig (V6).
* A JSON/CLI projection of committedBy for scripted walks: **NEW ROW NEEDED**. It's out of scope; the
  browser walk suffices.
* Attribution of `PutObject`/journal objects to the operation that stored them (registry publish,
  broker request): **NEW ROW NEEDED**. CB-NONE states the absence honestly, but it doesn't answer
  "who stored it".
* **NEW ROW NEEDED: writable Open silently recreates any dropped application table.**
  `enforceSchemaVersion` re-applies `schemaSQL` (`CREATE TABLE IF NOT EXISTS …`) on every writable
  open of a current-version store. A dropped `worlds`, `log_entries`, `objects`, etc. therefore comes
  back empty with no error (V15, measured for `worlds` on base v3 and v4). This row guards only
  `commit_objects`, whose emptiness would be read as a definite answer (CB-NONE). The general fix is
  a verify-before-DDL for every table, and it's out of scope here.
* Lock-blocked read status: row 111. Commit deadline policy: row 23 / PR #153 (D-WORLD-40).

## 12. Decision for Mark

None. The storage/upgrade choice follows the twice-applied fail-loud schema-bump precedent (V5), and
it strands no real history (V6).

## 13. Quorum log

**r1: BLOCKED 2/2.** `oc-glm` and `oc-kimi` were unreachable, so there were no verdicts from those
seats. The controller measured each objection's premise before r2.

| Reviewer | Verdict | Objection (one line) | Controller measurement | Disposition |
|---|---|---|---|---|
| gpt6-astra | REJECT | On an existing v4 store, writable Open's `db.Exec(schemaSQL)` would silently recreate a dropped `commit_objects` empty. Historical carriers would then render CB-NONE, so the §3.5/§3.7 "missing table → 500" and "no unknown" claims are false. | **Premise TRUE.** Base `store.go:385-387` re-applies `schemaSQL` when `user_version == current` on writable Open. | **Applied verbatim in substance.** Pre-DDL integrity check in both open modes, and only fresh init creates the table (§3.1). §3.5/§3.7 revised, and the "re-apply guarantees completeness" claim removed. AC16 + `TestDroppedMembershipTableIsRefusedNotRecreated` added. MUT-18/19 killed. V14 + V15 added. The pre-existing all-tables property is filed as NEW ROW NEEDED (§11). |
| gemini-3-1-pro | REJECT | Step 4b's `h.EntryIndex` is the parent head's index, so membership is attributed to the wrong entry. | **Premise FALSE.** `h := c.Entry.Header` is the new entry's header, and Step 4 inserts that same `h.EntryIndex` as this commit's `log_entries.entry_index`. AC6 `[1 3]` pins the attribution. | **Refuted by measurement.** No code change. V13 added, and §3.2 now states that `h` is the new entry's header. |
| oc-glm | — | unreachable | — | no verdict |
| oc-kimi | — | unreachable | — | no verdict |

**r2: BLOCKED 2/2** (same seats; `oc-glm`, `oc-kimi` unreachable again). Each objection was measured by
the controller, which then made the bounded narrow-refinement revision itself (no r3).

| Reviewer | Verdict | Objection (one line) | Controller measurement | Disposition |
|---|---|---|---|---|
| gemini-3-1-pro | REJECT | §3.1's `PRAGMA foreign_keys = ON` claim has no Verification Log row. | **Claim TRUE, citation missing.** `store.go:350` (base `341`). | **Applied verbatim.** V17 added; §3.1 cites it. |
| gpt6-astra | REJECT | Commit stays unbounded (`context.Background()` guard read, `tx.Exec`); make M1 wait for row 23 / PR #153's bounded Commit, or implement it here. | **Premise partly TRUE, and pre-existing.** The unbounded Commit is DR-1 at base (`store.go:964-966`). Step 4b adds no new wait point: it writes after Step 2 has taken the write lock (V18). PR #153 threads `ctx` through every step (`CommitContext`, `tx.ExecContext`) and is held on OPEN `D-WORLD-40`. | **The compatible half is applied verbatim** as the §3.2 rebase contract (Step 4b uses the bounded tx context; exhaustion leaves no entry, membership or head). **The landing dependency is not applied.** It would override an open human decision, and Gate 2 rule (c) routes a pre-existing defect to its owner (row 23), not into this doc. Recorded for the Gate-4 routing row and the report. |

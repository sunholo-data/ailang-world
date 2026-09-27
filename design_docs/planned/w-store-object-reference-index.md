# w-store-object-reference-index — checked reverse edges for object pages

Status: DESIGN, iteration 202, queue row 103. Base measured: `7fc05de` (V0).
Estimate: about one day, four independently testable milestones. This document is the
only deliverable; the prototype uses Go overlays under `/tmp`, not repository code edits.

## Problem

The object inspector currently existence-checks its interface edge, then renders named
`committedBy` and `referencedBy` stops (V3). Row 103 asks for the latter to become an
answer to “which entries or worlds reference this object?”; row 104 owns commit membership
(V8). An absent access path must never masquerade as proof that an object is unreferenced.

## Clause mapping

Clause 5 requires verified answers to real “why did X happen?” questions within five minutes.
This design supplies the reverse-hop step after locating an object. It does not assert that
an index microbenchmark completes the charter's real-question demonstration. The mission
row and row 96's implemented design establish that scope (V8, V9).

Acceptance walk: locate an incident object by semantic ID, open its object page, follow a
`transitionRef` backlink to the exact entry, and verify that entry's transition reference
and payload identity. Also exercise an interpreter object with several referring entries
and a state object with several worlds. Record the real question, answer, supporting refs,
elapsed end-to-end time and remaining stops. An interpreter backlink alone does not prove
causality; it is a recorded relationship for the walker to inspect.

## Measured current state

All present-tense repository claims below are supported by the Verification Log. Proposed
behavior in subsequent sections is a specification, not a claim of implementation.

* `schema.sql` declares nine tables and zero explicit indexes, with the semantic-ID index
  documented as provisioned outside the schema. This is not zero physical SQLite indexes:
  primary/unique constraints and the writable-open lookup index also exist (V1, V2, V10).
* Writable `Open` provisions `objects_by_semantic_id`; read-only open caches its verified
  availability. The verifier checks table membership, non-partialness, exact two key columns,
  BINARY collation and ascending order. Provisioning uses a 30-second deadline; the lookup
  refuses an unavailable index before querying, and the JSON adapter maps that error to 503
  (V2, V5). The new provisioning path should additionally use context for its PRAGMA reads.
* `ObjectsBySemanticID` is context-first, uses a strict keyset cursor, enforces 1–500 results,
  and selects metadata without payload (V2). `readStore` has six methods, including this
  lookup; the historical five-point-read description is stale (V4).
* The mux has nine `/v1` registrations plus the workbench. Workbench accepts object/payload
  queries, renders edge links from `EdgeView`, and has a 100-row page constant (V3, V4, V6).
  Its store error adapter handles timeout as 503 and other errors as sanitized 500 (V3).
* The pending row-23 document contains the D-WORLD-40 bound table; the mission records PR
  #153 as unmerged pending ratification. This is local-document evidence, not a fresh remote
  PR-status check (V7, V8).

### Reference census: actual file-backed Store, not a SQL-only mock (V10)

The probe uses existing `seedGenesis`, `obj`, and `testCommitIntent` helpers with production
`PutObject`, `PutWorld`, `SelectHead`, `Commit`, `SetRegistryHead`, and `AppendIntent` calls.
It commits 10,000 entries and deliberately stores function, interpreter, state and interface
objects so fixture omissions cannot produce false exclusions. It retains the helper's
original genesis world as a dangling-state control. Counts are source rows, not distinct refs.

| Column | Source rows | Join to `objects.hash_ref` | Decision |
|---|---:|---:|---|
| `log_entries.transition_ref` | 10,000 | 10,000 | Include |
| `log_entries.transition_fn_ref` | 10,000 | 10,000 | Include |
| `log_entries.interpreter_ref` | 10,000 | 10,000 | Include |
| `log_entries.prev_entry_hash_ref` | 10,000 | 0 | Exclude: entry identity, not object identity |
| `worlds.state_root` | 10,002 | 10,001 | Include |
| `worlds.log_head` | 10,002 | 0 | Exclude: entry identity, not object identity |
| `epoch_registry_heads.object_ref` | 1 | 1 | Genuine object ref; outside entries/worlds scope |
| `journal.object_ref` | 1 | 1 | Genuine object ref; outside entries/worlds scope |
| `objects.interface_hash_ref` | 10,005 | 1 | Genuine object ref; outside entries/worlds scope |

The same probe supplies positive controls: **9,999** previous-entry values resolve to
`log_entries.entry_hash_ref`, and **10,000** world log heads resolve there. All three
log object-ref columns resolve in the same census. Zero object joins are therefore observed
negatives with both object-namespace and entry-namespace controls. Schema comments likewise
define log heads/previous hashes as entry identities (V1).

“Never resolves” here means zero matches in this measured fixture, not a universal SQL
invariant. A HashRef is not a namespace tag, and coincident object bytes could have an entry's
hash. Do not classify columns by accidental cross-namespace equality. Registry, journal and
interface references must **not** be described as non-object references merely to narrow scope.

### Cost at N = 10,000 (V10)

Instrument: `TestRow103Measure`, reproduced in the appendix. Three independent file stores;
101 sequential warm-query samples per target and index state, p50 and maximum, no network.
The read executes the four ordered branches with a shared remaining limit of 101, stopping
when full. This measures query/scan cost, not HTML rendering or the later point checks.

| Target | Without indexes p50 range | Indexed p50 range | Worst without / indexed |
|---|---:|---:|---:|
| Sparse | 11.056–11.702 ms | 52.125–54.541 µs | 60.803 ms / 180.042 µs |
| Hot function | 1.737–1.937 ms | 52.208–54.250 µs | 5.021 ms / 66.500 µs |
| Hot world | 4.677–4.817 ms | 85.417–86.250 µs | 12.547 ms / 363.792 µs |
| Absent | 10.484–10.611 ms | 52.750–53.708 µs | 16.307 ms / 73.459 µs |

Sparse means the last transition body (one edge); hot means the shared function (101 returned
of 10,000); hot-world means the shared state (101 of 10,001); absent means no matches, paired
with the sparse and hot positive controls. Before indexing the plans search the entry integer
primary key and the world primary-key index, filtering the object-ref column. Afterwards
all four use the intended **covering indexes**, with both equality and cursor-range terms
(V10). No payloads are fetched by the measured reverse queries.

Writable Open without new provisioning: **0.517 / 0.518 / 0.521 ms**.
Writable Open with four-index provisioning plus strict verification:
**44.943 / 45.046 / 44.124 ms**. One measured Open per fresh seeded store, not a p99 claim.
After closing: **12,304,384 → 16,424,960 bytes** in all three runs,
a **4,120,576-byte (3.93 MiB, 33.5%)** increase. This measures the main database after Close,
not transient WAL high-water usage; seeding time is excluded.

The provisioning prototype adds the four indexes and strict PRAGMA inspection **inside
writable Open**, via a temporary overlay of `provisionLookupIndex`. It reuses the existing
provisioning deadline. This is a real Open timing, not DDL timing mislabeled as startup.
The final separate helper, read-only guard, cursor decoding, HTTP adapter and existence checks
are designed here but not prototyped. The sample is synthetic mission-scale data using the
real Store; no claim is made that it is a production incident archive or a cold-cache test.

## Decisions (chosen option and alternatives)

**D1 — Four entry/world reference relations, with explicit scope.** Include transition body,
transition function, interpreter and state root. Exclude entry-chain columns by their measured
namespace; defer registry, journal and object-interface inbound relations by scope, not by
claiming their values are not objects. Render the scope on every object page. Alternative:
index all seven genuine object-reference columns now. Rejected for this row because it adds
source kinds and registry/journal navigation semantics beyond “entries/worlds referencing X”.
Alternative: infer references from any equal hash or payload text. Rejected: equality alone
does not establish the intended typed relation. Evidence: census and controls (V10).

**D2 — Four composite indexes, writable-open provisioning, fail-closed availability.**
Follow row 96's access-path pattern without a schema version bump. Alternative: single-column
indexes; rejected because the world key needs explicit ordered hash pagination. Entry indexes
also declare their stable key explicitly, even though SQLite integer-primary-key storage may
already supply it. Alternative: unconditional DDL in schema or fallback scans; rejected in
favor of bounded provisioning and an explicit 503. Evidence: existing pattern (V2), covering
query plans and measured scan/index costs (V10).

**D3 — One context-first, globally bounded, relation-ordered keyset read.** Return reference
edges, not deduplicated source entities; an entry that references X as both function and
interpreter supplies two differently labelled facts. Alternative: OR across three log fields
plus global sorting/deduplication. Rejected because branch-local indexes and fixed relation
ordering give simple bounded work and an unambiguous cursor. Alternative: four seam methods;
rejected because the caller would own aggregate limits and pagination. Evidence: four branch
plans and hot-key measurements (V10), six-method seam (V4).

**D4 — Workbench only; 100 displayed edges, lookahead, explicit empty state, continuation.**
Keep all nine `/v1` routes unchanged. Add an object-only `refsAfter` HTML query parameter.
Alternative: a tenth `/v1` route and matching CLI/sketch work; unnecessary for the object-page
walk. Alternative: truncate without a continuation; rejected because a real explanation may
sit beyond the first page. Evidence: current mux and workbench grammar/page bound (V3, V4, V6).

**D5 — Share caller cancellation; inherit policy, do not ratify a new duration.**
Pass the workbench `readCtx` through the reverse read and every existence check. For new index
provisioning reuse row 96's existing 30-second ceiling as a compatibility precedent, not a new
row-23 policy value. Once row 23 changes Open to take ctx, derive the provisioning context
from that caller, so the earlier deadline wins; integrate under its B1 startup owner and B7
connection-wait rule. Alternative: Background per query or a freshly invented reverse-read
budget. Rejected: it would defeat the request deadline or bypass D-WORLD-40. Evidence: existing
`readCtx` and provisioning code (V2, V5), pending table (V7).

## Design

### Schema and index lifecycle

Provision these outside `schema.sql`; add only a pointer comment there:

```sql
CREATE INDEX IF NOT EXISTS log_entries_by_transition_ref
  ON log_entries(transition_ref, entry_index);
CREATE INDEX IF NOT EXISTS log_entries_by_transition_fn_ref
  ON log_entries(transition_fn_ref, entry_index);
CREATE INDEX IF NOT EXISTS log_entries_by_interpreter_ref
  ON log_entries(interpreter_ref, entry_index);
CREATE INDEX IF NOT EXISTS worlds_by_state_root
  ON worlds(state_root, world_ref);
```

Use one bounded context for the entire four-index provisioning/verification group. Execute
DDL and both PRAGMA forms with that ctx; close rows before the next query on the single
connection. Integrate both memory and file-backed writable Open branches. On any error,
close the DB and release an acquired writer lock before returning failure. Partial successful
DDL is harmless access-path residue: a later writable open retries idempotently.

For each index, `index_list(expectedTable)` must find the exact name with `partial == 0`;
`index_xinfo(name)` must report exactly the two declared key columns in sequence, BINARY,
ascending, ignoring only auxiliary `key == 0` rows. Reject expressions, wrong table, missing
columns, reversed order, DESC, NOCASE, partial and extra-key definitions. Do not drop/rebuild
an incompatible same-named index silently. Do not change `user_version`.

Cache `referenceIndexesAvailable` on the Store. Writable Open must provision then verify all
four or fail. Read-only Open only verifies under the reused provisioning ceiling (and, after row 23,
its caller context): absent/incompatible indexes leave other reads
usable, but any reverse read returns `*ReferenceIndexUnavailableError` before SQL. PRAGMA I/O
errors fail Open rather than caching a false verdict. A read-only handle stays unavailable
until reopened after provisioning; do not introduce per-request schema probes. This assumes
schema is not externally rewritten behind an open handle, as does the row-96 precedent (V2).

### Store API, ordering and cursor

```go
const MaxObjectReferencePage = 500
// Kind values, in order: TransitionRef, TransitionFn, Interpreter, StateRoot.
type ObjectReferenceCursor struct {
    Kind ReferenceKind
    EntryIndex int64       // log kinds only; nonnegative
    WorldRef hashref.HashRef // StateRoot only
}
type ObjectReference struct { Cursor ObjectReferenceCursor }
func (s *Store) ObjectReferences(ctx context.Context, ref hashref.HashRef,
    after *ObjectReferenceCursor, limit int) ([]ObjectReference, error)
```

Validate ref, cursor shape and `1 <= limit <= 500` at the store boundary. Use the existing
`InvalidLimitError` for limit rejection (V2); introduce a typed invalid-cursor error. Nil cursor
means before every relation, including entry index 0. Reject unknown kinds, negative log
indices and invalid world refs; inactive fields must be zero. A syntactically valid cursor
need not identify an existing row. No OFFSET, rowid cursor, wildcard ref matching or payload.

Order is `(Kind, sourceKey)`, not timestamp or hash of the referenced object: numeric ascending
entry index for the three log branches and BINARY ascending world_ref for the world branch.
Skip kinds before `after.Kind`; apply strict `>` in the matching kind; start later kinds from
the beginning. For a first log page use `-1` only as an internal SQL sentinel. For a first world
page use empty text. Each branch requests only `limit - len(items)`; stop at the total limit.
Empty results are a non-nil empty slice. Close each result set and check `rows.Err()` before
starting another; on any error discard accumulated results, never return partial success.

```sql
SELECT entry_index FROM log_entries
WHERE transition_ref = ? AND entry_index > ? ORDER BY entry_index LIMIT ?;
-- Same SQL shape for transition_fn_ref and interpreter_ref.
SELECT world_ref FROM worlds
WHERE state_root = ? AND world_ref > ? ORDER BY world_ref LIMIT ?;
```

All four queries are covering with the specified indexes, as measured in V10. GetLogEntry
supplies the source entry hash during the mandatory existence check; returning it from the
reverse read would add an unnecessary table lookup. Results contain no payload or source body.

Parse stored world refs with `hashref.Parse`; reject invalid or negative stored source keys.
Use `QueryContext(ctx, ...)` and check caller cancellation between branches. The read has no
snapshot across pages: immutable source keys give stable ordering, while new references that
sort before the cursor require restarting the walk. Do not claim historical completeness or
hide this limitation in the continuation UI.

### Daemon seam and workbench rendering

Add exactly the `ObjectReferences` method to `readStore`. Update blocking, recording and
failing seam implementations and their call-record assertions. They currently implement the
semantic-ID seam too (V4). Derive the context once via `d.readCtx(r)`; no timer per backlink.

For an existing object, request 101 edges from the new read, retain 100 and set `Truncated`
only if the 101st exists. For each retained log edge call `GetLogEntry(ctx, index)`, check
that the indicated relation still equals the inspected object's hash, then link through
`pageHref(index, index)` to `/workbench?from=<index>&entry=<index>`. For state-root edges call
`GetWorld(ctx, worldRef)`, compare StateRoot, then link to `/workbench?world=<ref>`. Use the
existing entry-page overflow constraint: if an index cannot be represented by that page's
`from` window, use a safe from value and the independent selected-entry panel; test MaxInt64 boundaries.
The independent selected-entry lookup is measured in V13. Specifically use
`from = min(index, math.MaxInt64-WorkbenchPageLimit)`; never emit an invalid
link. Display source kind, source key and reference role; these are factual links, not claims
that the entry committed the object.

A missing source from a seam/race yields a named UNAVAILABLE row with its key and no link;
a mismatched source relation or invalid stored key is a sanitized 500. A point-read error
aborts the page through the error adapter. Never silently drop missing rows or call them
“none”. At most 100 source point checks plus the four bounded reverse queries are added.

Keep interface and `committedBy` behavior unchanged. Remove only `objectReferencedByMissing`.
Add a `ReferenceView` to `ObjectView` with edges, `None`, `Truncated`, and `NextHref` instead of
misusing `EdgeView.Missing` for a successful empty result. Render a single `referencedBy`
section after committedBy with the fixed visible scope:

> Entries/worlds only: transitionRef, transitionFn, interpreter, stateRoot.
> Registry, journal and object-interface inbound references are not included.

* First page, zero rows: **“none recorded in these entry/world fields”** (HTTP 200).
* Continued page, zero rows: **“no further references recorded in these entry/world fields”**.
* Nonempty: labelled, existence-checked links (or named missing-source rows).
* More than 100: **“Showing 100 references; more recorded”** and a **Next references** link.
* Continuation: note **“New references may require restarting this walk.”** and a first-page link.

`refsAfter` encodes a versioned, base64url JSON cursor with exact fields `v`, `kind`, `key`
(integer decimal string for entries, canonical HashRef for worlds). Cap encoded input at
512 bytes, reject unknown/duplicate JSON fields, unsupported versions, noncanonical integers,
wrong kinds and invalid refs with 400. This is a cursor, not authorization or a snapshot token.
Accept it only with `object` and optional `payload`; preserve the current object's ref and
payload flag in next/first links, using `url.Values` and HTML escaping. Reject duplicate query
keys and combinations with world/from/entry. Build next from the **last displayed** edge,
never from the hidden lookahead row; preserve a missing source's position too.

Teach `writeWorkbenchStoreError` to recognize `ReferenceIndexUnavailableError` and return
HTTP **503**, HTML class `ReferenceIndexUnavailable`, with fixed remediation:
“object-reference indexes are absent or incompatible; open this store writable once to
provision them, then reopen this read-only handle”. Do this before writing success content.
Missing inspected objects still return 404. Context expiry remains 503 Timeout; other store
errors remain sanitized 500. No machine-route or AILANG API error variant is required.

## Failure modes

| Failure | Required outcome |
|---|---|
| Any one index missing/partial/incompatible on read-only handle | 503, no reverse SQL, never empty success |
| Writable DDL fails, times out or verifies incompatible | Open fails; cleanup releases writer ownership; retry is possible |
| Cursor bad or limit outside store bound | Typed rejection; HTTP cursor error 400; no query |
| Context cancelled before/during read or connection wait | Cancellation propagates; workbench 503 Timeout; no partial page |
| Source key malformed, relation differs, SQL/scan error | Sanitized 500, no false edges |
| Source no longer present | Named unavailable source with no link, not “none” |
| Exactly 100 / 101 matches | No truncation / explicit truncation and continuation respectively |
| One entry uses X in several fields | Several role-labelled edges; one cursor position per role/source pair |
| References added during pagination | Stable keyset semantics, restart notice; no snapshot promise |
| Read-only handle opened before another writer provisions | Remains 503 until reopened |

## Milestones

All names below describe tests to implement, not tests already passing. Each milestone must
leave a compiling, bisectable tree. Run the common fence after each:

```sh
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
go vet ./host/...
go test -run '^$' ./host/...
```

**M1 — Access paths and availability (store only).** Add four indexes, strict context-bound
verification, cached verdict, typed error and writable/read-only Open integration. Add schema
comment; no seam/renderer changes. Test fresh and pre-index files, read-only no-DDL, all
incompatible definitions, failure cleanup/reopen, unchanged version, existing semantic lookup.

```sh
go test -race -count=1 ./host/store -run '^(TestReferenceIndexLifecycle|TestReferenceIndexDefinitions|TestReferenceIndexFailureCleanup|TestLookupIndexGuardReadOnlyThenProvision)$'
```

**M2 — Bounded reverse read (store only).** Add typed cursors and four ordered queries. Test
each relation, simultaneous roles, empty result, strict cursor, equal-index cross-kind
boundaries, absent cursor row, invalid limits/refs, hot targets, cancellation and plans.
Keep it unexposed until M3. Re-run the appendix at N=10,000 against the final implementation.

```sh
go test -race -count=1 ./host/store -run '^(TestObjectReferencesPage|TestObjectReferencesScope|TestObjectReferencesValidation|TestObjectReferencesGuard|TestObjectReferencesCancellation|TestObjectReferencesPlans)$'
```

**M3 — Seam, checked first page and errors.** Add the seventh seam method and update every
fake in the same change. Add ReferenceView and render scoped empty/checked/truncated states,
503 unavailable mapping, point-read timeout/error paths. Until M4, show truncation text but
no broken Next link. Preserve committedBy test coverage and frozen-route parity.

```sh
go test -race -count=1 ./host/daemon ./host/workbench -run 'TestWorkbench.*Reference|TestWorkbenchObjectProvenanceWalk|TestFrozenV1RouteTableMatchesMux|TestReferenceView'
```

**M4 — Continuation and walk verification.** Extend object grammar, cursor codec, next/first
links, page partition tests and selected-entry boundary cases. Exercise Handler with
`httptest.NewRecorder` (no server). Run a real question walk as the acceptance evidence;
report latency separately from the microbenchmark. Execute mutations one at a time, restore
bytes, then final fences and touched-package race suite:

```sh
go test -race -count=1 ./host/daemon -run 'TestWorkbenchReferencePaging|TestWorkbenchReferenceCursor|TestWorkbenchReferenceWalk'
go test -race -count=1 ./host/store ./host/daemon ./host/workbench
go test -run '^$' ./host/...
```

Estimate M1/M2 around 100–150 production lines each and M3/M4 around 80–150 each. These are
estimates, not measured patch sizes. Split further if necessary; do not leave uncompilable
seam changes between milestones. If any socket-dependent gate cannot bind in the sandbox,
report **UNINFORMATIVE UNDER SANDBOX**, and rerun in a socket-capable environment.

## Test plan and MUTATION TABLE

Positive controls are mandatory alongside absence tests: a matching source per included
field, an absent target, a present but out-of-scope registry/journal/interface ref, and actual
entry-chain values that do not resolve to objects. Compare complete ordered pages with a
raw unguarded SQL oracle in tests; dropping indexes may be used only for that oracle and
cost measurement, never as a serving fallback. Assert each returned link's GET status and
identity, not merely its href text. Exercise 0, 1, 100, 101 and 501 edges, numeric indices 2
and 10, multiple worlds inserted out of hash order, same-source multi-role references, and
mixed kinds across a page boundary. Guard tests must prove no SQL ran after refusal.

Mutation literals below are implementation contracts for the proposed code. Index mutations
apply in `host/store/reference_index.go`, query mutations in `host/store/object_references.go`,
and adapter/cursor mutations in `host/daemon/workbench.go`; renderer text lives in
`host/workbench/render.go`. Name the source checker `checkedReferenceEdge`. Require exactly
one match (parameterize the named per-column SQL where stated), change one mutation at a
time, run the named test, require failure, restore and compare original bytes. A comment-only
mutation is the green control. These mutations have **not** been executed against future code.

| ID | Exact code mutation | Test that must kill it |
|---|---|---|
| MU1 | In reference-index specs, delete `{"worlds_by_state_root", "worlds", "state_root", "world_ref"}` | `TestReferenceIndexLifecycle` asserts all four names |
| MU2 | Reference verifier: `partial == 0` → `true` | `TestReferenceIndexDefinitions/partial` |
| MU3 | Reference verifier: `coll == "BINARY"` → `true` | `TestReferenceIndexDefinitions/nocase` |
| MU4 | Reference verifier: `desc == 0` → `true` | `TestReferenceIndexDefinitions/descending` |
| MU5 | Reference verifier: `return valid && keys == 2, nil` → `return true, nil` | `TestReferenceIndexDefinitions/extra-key` |
| MU6 | Reference query: `if !s.referenceIndexesAvailable {` → `if false {` | `TestObjectReferencesGuard` missing-index SQL-hook control |
| MU7 | Each log SQL constant, separately: `entry_index > ?` → `entry_index >= ?` | `TestObjectReferencesPage` full-page partition |
| MU8 | World SQL: `world_ref > ?` → `world_ref >= ?` | `TestObjectReferencesPage/world-cursor` |
| MU9 | Each column SQL, separately: `transition_ref = ?` (or its fn/interpreter/state counterpart) → that column followed by `<> ?` | `TestObjectReferencesScope` positive and decoy fixtures |
| MU10 | Branch limit argument: `limit - len(items)` → `limit` | `TestObjectReferencesPage/mixed-kind-bound` |
| MU11 | Store bound: `limit < 1 || limit > MaxObjectReferencePage` → `false` | `TestObjectReferencesValidation/limits` |
| MU12 | Query call: `QueryContext(ctx,` → `QueryContext(context.Background(),` | `TestObjectReferencesCancellation` cancelled ctx and blocked sole connection |
| MU13 | Cursor-kind skip: `kind < after.Kind` → `kind <= after.Kind` | `TestObjectReferencesPage/within-kind` |
| MU14 | Branch SQL: `ORDER BY entry_index` → `ORDER BY entry_index DESC` | `TestObjectReferencesPage/numeric-order` |
| MU15 | `checkedReferenceEdge`: `if !ok {` → `if false {` (separately for log/world branches) | `TestWorkbenchMissingReferenceSource` asserts no href |
| MU16 | Reference relation comparison: `if actual != ref {` → `if false {` | `TestWorkbenchReferenceMismatch` expects 500 |
| MU17 | Reference-index HTML error status: `http.StatusServiceUnavailable` → `http.StatusOK` | `TestWorkbenchReferenceUnavailable` expects 503, no “none” |
| MU18 | Lookahead: `WorkbenchPageLimit+1` → `WorkbenchPageLimit` in reverse call | `TestWorkbenchReferencePaging/101` |
| MU19 | Truncation: `len(refs) > WorkbenchPageLimit` → `len(refs) >= WorkbenchPageLimit` | `TestWorkbenchReferencePaging/100` |
| MU20 | Next cursor: `refs[WorkbenchPageLimit-1]` → `refs[WorkbenchPageLimit]` | `TestWorkbenchReferencePaging/no-skipped-lookahead` |
| MU21 | Empty-state text: `none recorded in these entry/world fields` → empty string | `TestReferenceView/empty` |
| MU22 | In object handler, reverse-read ctx argument `ctx` → `context.Background()` | `TestWorkbenchReferenceCancellation` checks parent cancellation and shared deadline |
| MU23 | Read-only Open: cached `referenceIndexesAvailable: available` → `referenceIndexesAvailable: true` | `TestReferenceIndexLifecycle/read-only-absent` |
| MU24 | Writable Open: `if err := provisionReferenceIndexes(db); err != nil {` → `if err := provisionReferenceIndexes(db); false {` | `TestReferenceIndexFailureCleanup` checks failure plus immediate reopen |
| MU25 | First log query cursor sentinel: `int64(-1)` → `int64(0)` | `TestObjectReferencesPage/index-zero` |
| MU26 | Cursor decoding version guard: `v != 1` → `false` | `TestWorkbenchReferenceCursor/version` |

## Non-goals

Row 104's `committedBy`, commit membership persistence, inference from provenance text,
payload reference discovery, registry/journal/interface backlinks, verification-cache reverse
reads, a tenth `/v1` route, a CLI command, global chronological order, exact reference counts,
cross-page snapshots, or ratification of D-WORLD-40. Do not edit the kernel or frozen API sketch.

## Conflict Surface

Planned implementation files (this design deliverable edits none of them):

* `host/store/schema.sql`: comment only.
* `host/store/store.go`: cached flag and both Open paths; overlap with row-23 Open/ctx lifecycle.
* New `host/store/reference_index.go`, `object_references.go` and corresponding `_test.go` files.
  Prefer a separate verifier helper over refactoring `lookup_index.go` in this row.
* `host/daemon/daemon.go`: seventh seam method only, no routes.
* `host/daemon/read_deadline_test.go`: blocking/recording/failing seam implementations.
* `host/daemon/workbench.go`, `workbench_test.go`, new `workbench_references_test.go`:
  reverse adapter, grammar, checked links, errors and tests. Coordinate with row 104 around
  `objectEdges`; do not remove or change its committedBy stop.
* `host/workbench/render.go`, `render_test.go`: ReferenceView and scoped renderer.

## VERIFICATION LOG

Commands run from the worktree root unless stated. Output is trimmed; absence claims include
known-positive controls within the same command or the same probe execution. V10 embeds the
full executable instrument so the measurements do not depend on retained `/tmp` files.

| ID | Exact command | Observed output / claim supported |
|---|---|---|
| V0 | `git rev-parse --short HEAD; git diff --name-only; $HOME/.pinned-ailang/ailang --version; go version; uname -m` | `7fc05de`; no tracked diff paths; positive git control is the HEAD hash; AILANG `v0.41.0`, Go `go1.26.6 darwin/arm64`, `arm64` |
| V1 | `cat host/store/schema.sql; rg -n 'CREATE TABLE\|CREATE INDEX' host/store/schema.sql` | Nine CREATE TABLE declarations, no CREATE INDEX (same search's table hits are positive control); fields and comments identify object, entry-chain and registry/journal refs |
| V2 | `cat host/store/lookup_index.go host/store/objects_by_semantic_id.go; sed -n '235,330p' host/store/store.go` | 30s `lookupIndexProvisionDeadline`; CREATE/strict PRAGMAs; writable provision and readonly verify/cache; `MaxSemanticIDPage = 500`, guard, strict `hash_ref > ?`, `QueryContext` |
| V3 | `sed -n '1,260p' host/daemon/workbench.go; sed -n '825,880p' host/daemon/workbench_test.go` | checked interface; committedBy/referencedBy named stops; allowed object/payload grammar; timeout 503/internal 500; named-stops and edge-order tests pin current rendering |
| V4 | `sed -n '365,373p' host/daemon/daemon.go; sed -n '630,650p' host/daemon/daemon.go; rg -n 'ObjectsBySemanticID' host/daemon/read_deadline_test.go` | Six seam signatures; nine /v1 patterns and workbench; blocking line 173, recording 235, failing 727 |
| V5 | `sed -n '256,287p' host/daemon/handlers.go; rg -n 'lookupIndex\|LookupIndex' host/daemon/handlers.go` | `WithTimeout(r.Context(), d.readDeadline)`; index error maps to `http.StatusServiceUnavailable` at line 461 |
| V6 | `sed -n '34,68p' host/workbench/render.go; rg -n 'EdgeView\|WorkbenchPageLimit\|UNAVAILABLE\|Relation' host/workbench` | EdgeView fields, ObjectView.Edges; limit 100; unavailable fallback and edge partial |
| V7 | `sed -n '20,39p' design_docs/planned/w-store-bounded-durable-operations.md` | Table “for Mark's ratification”; B1 startup, B7 no separate timer/parent connection budget |
| V8 | `rg -n '^103\.|^104\.' design_docs/world-mission.md; rg -n '^23\.' design_docs/world-mission.md` | Row 103 entry/world reverse relation; row 104 membership; row 23 says PR #153 NOT MERGED pending D-WORLD-40 |
| V9 | `cat design_docs/implemented/w-object-lookup-by-semantic-id.md design_docs/implemented/w-object-lookup-by-semantic-id-sprint-plan.md` | Read house problem/design/verification structure, M1a/M1b/M2/M3, one-at-a-time literal mutation rules, compile/race gates. Prior timing values are not reused as fresh measurements here |
| V10 | Appendix reconstruction, then `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -overlay /tmp/world103/overlay.json -count=3 -timeout 120s ./host/store -run '^TestRow103Measure$' -v` | Census, positive controls, plans, page timings, writable-Open timings and disk delta reported above; three PASS runs. Each empty join/query has nonempty controls in the same test |
| V11 | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -run '^$' ./host/...` | Exit 0; host packages compile, including test files (`[no tests to run]`) |
| V15 | `git diff --check; git diff --name-only; git ls-files --others --exclude-standard; wc -l design_docs/planned/w-store-object-reference-index.md` | No whitespace errors or tracked-file changes; positive untracked control lists only this new document; final line count 557 |
| V14 | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -overlay /tmp/world103/overlay.json -run '^$' ./host/store` after running the appendix reconstruction verbatim | Exit 0; reconstructed overlay compiles, including the test probe |
| V13 | `sed -n '310,415p' host/daemon/workbench.go` | Selected entry uses its own GetLogEntry; timeline uses a separate from/offset loop and bounded next arithmetic |
| V12 | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race -count=1 -timeout 90s ./host/store -run '^(TestObjectsBySemanticIDNonUniqueOrderedAndPaged\|TestLookupIndexGuardReadOnlyThenProvision\|TestLookupIndexIncompatibleIsRefused)$'` | Exit 0; `ok .../host/store 1.510s` |

Not measured: final HTML/point-check latency, final implementation cancellation and mutation
kills, and an end-to-end real incident walk. These need implementation and an identified real
incident, respectively. No server was started; no socket result is being treated as evidence.

## Reproducible measurement appendix

Run from the base worktree. This reconstructs two Go overlays under `/tmp`; it does not modify
repository code. The first copies an existing test file and adds the probe; the second injects
an opt-in provisioning call inside writable Open's existing lookup-index step. The test flag
is set only after seeding and the baseline Open timing. Each run owns a temporary database,
closes it before file-size sampling, and lets `testing` remove it. Test timeout bounds execution.

```sh
mkdir -p /tmp/world103
python3 - <<'PY'
from pathlib import Path
import json, re
root = Path.cwd()
p = Path('/tmp/world103')
doc = (root/'design_docs/planned/w-store-object-reference-index.md').read_text()
def block(name):
    return re.search(r'<!-- '+name+r' -->\n```go\n(.*?)\n```', doc, re.S).group(1)
s = (root/'host/store/store_test.go').read_text()
s = s.replace('"context"', '"context"\n "fmt"\n "os"\n "path/filepath"\n "sort"\n "time"')
(p/'store_test.go').write_text(s + '\n' + block('probe-test') + '\n')
s = (root/'host/store/lookup_index.go').read_text().replace('"fmt"', '"fmt"\n "os"')
s = s.replace('\treturn nil\n}', '\tif os.Getenv("WORLD103_PROVISION") == "1" { return row103Provision(db) }\n\treturn nil\n}', 1)
(p/'lookup_index.go').write_text(s + '\n' + block('probe-provision') + '\n')
(p/'overlay.json').write_text(json.dumps({'Replace': {
    str(root/'host/store/store_test.go'): str(p/'store_test.go'),
    str(root/'host/store/lookup_index.go'): str(p/'lookup_index.go')}}))
PY
AILANG_BIN=$HOME/.pinned-ailang/ailang go test -overlay /tmp/world103/overlay.json -count=3 -timeout 120s ./host/store -run '^TestRow103Measure$' -v
```

<!-- probe-test -->
```go
func TestRow103Measure(t *testing.T) {
 path := filepath.Join(t.TempDir(), "world.db")
 s, err := Open(path); if err != nil {t.Fatal(err)}
 must := func(err error) {t.Helper(); if err != nil {t.Fatal(err)}}
 iface := obj("reference-interface", "interface"); must(s.PutObject(iface))
 fn := obj("reference-fn", "fn"); fn.InterfaceHash=iface.Hash; must(s.PutObject(fn))
 interp := obj("reference-interpreter", "interpreter"); must(s.PutObject(interp))
 state := obj("reference-state", "state"); must(s.PutObject(state))
 current := seedGenesis(t,s); current.StateRoot=state.Hash
 // Replace fixture's genesis with a separate, stored-state genesis before commits.
 current.Ref=hashref.SumSHA256([]byte("reference-genesis")); must(s.PutWorld(current)); must(s.SelectHead(current.Ref))
 must(s.SetRegistryHead("reference-registry",fn.Hash))
 var last Commit
 for i:=int64(0); i<10000; i++ {
  body:=obj(fmt.Sprintf("body-%d",i),"transition/body")
  eh:=hashref.SumSHA256([]byte(fmt.Sprintf("entry-%d",i)))
  c:=Commit{ObservedHead:current.Ref,Objects:[]Object{body},NextWorld:World{Ref:hashref.SumSHA256([]byte(fmt.Sprintf("world-%d",i))),Revision:i+1,StateRoot:state.Hash,LogHead:eh},Entry:LogEntry{Header:LogHeader{EntryIndex:i,SemanticsEpoch:1,TransitionFn:fn.Hash,Interpreter:interp.Hash,PrevEntryHash:current.LogHead,WrittenBy:"row103"},EntryHash:eh,TransitionRef:body.Hash}}
  must(s.Commit(c)); current=c.NextWorld;last=c
 }
 _,_,err=s.AppendIntent("row103",testCommitIntent("row103",last));must(err)
 cols:=[][2]string{{"log_entries","transition_ref"},{"log_entries","transition_fn_ref"},{"log_entries","interpreter_ref"},{"log_entries","prev_entry_hash_ref"},{"worlds","state_root"},{"worlds","log_head"},{"epoch_registry_heads","object_ref"},{"journal","object_ref"},{"objects","interface_hash_ref"}}
 for _,c:=range cols {var total,resolved int;must(s.db.QueryRow("SELECT count(*),count(o.hash_ref) FROM "+c[0]+" r LEFT JOIN objects o ON r."+c[1]+"=o.hash_ref").Scan(&total,&resolved));t.Logf("census %s.%s total=%d resolved=%d",c[0],c[1],total,resolved)}
 var control int;must(s.db.QueryRow("SELECT count(*) FROM log_entries p JOIN log_entries e ON e.prev_entry_hash_ref=p.entry_hash_ref").Scan(&control));t.Logf("prev_entry positive entry-namespace control=%d",control)
 must(s.db.QueryRow("SELECT count(*) FROM worlds w JOIN log_entries e ON w.log_head=e.entry_hash_ref").Scan(&control));t.Logf("log_head positive entry-namespace control=%d",control)
 _=[]string{"CREATE INDEX log_entries_by_transition_ref ON log_entries(transition_ref,entry_index)","CREATE INDEX log_entries_by_transition_fn_ref ON log_entries(transition_fn_ref,entry_index)","CREATE INDEX log_entries_by_interpreter_ref ON log_entries(interpreter_ref,entry_index)","CREATE INDEX worlds_by_state_root ON worlds(state_root,world_ref)"}
 query:=func(ref string,limit int,plans bool) int {
  total:=0
  for _,c:=range cols[:5] {if c[1]=="prev_entry_hash_ref" {continue};key:="entry_index";after:=any(int64(-1));if c[0]=="worlds" {key="world_ref";after=""}
   q:="SELECT "+key+" FROM "+c[0]+" WHERE "+c[1]+"=? AND "+key+">? ORDER BY "+key+" LIMIT ?"
   if plans {r,e:=s.db.Query("EXPLAIN QUERY PLAN "+q,ref,after,limit);must(e);for r.Next(){var a,b,c int;var d string;must(r.Scan(&a,&b,&c,&d));t.Log(d)};must(r.Close())}
   r,e:=s.db.QueryContext(context.Background(),q,ref,after,limit-total);must(e);for r.Next(){var v any;must(r.Scan(&v));total++};must(r.Err());must(r.Close());if total>=limit {break}
  };return total
 }
 measure:=func(label string){for _,x:=range []struct{name,ref string}{{"sparse",last.Entry.TransitionRef.String()},{"hot",fn.Hash.String()},{"hot-world",state.Hash.String()},{"absent",hashref.SumSHA256([]byte("absent")).String()}} {samples:=[]int64{};count:=0;for i:=0;i<101;i++{start:=time.Now();count=query(x.ref,101,false);samples=append(samples,time.Since(start).Nanoseconds())};sort.Slice(samples,func(i,j int)bool{return samples[i]<samples[j]});t.Logf("%s %s rows=%d p50_us=%.3f pmax_us=%.3f",label,x.name,count,float64(samples[50])/1000,float64(samples[100])/1000)};query(last.Entry.TransitionRef.String(),101,true)}
 measure("unindexed")
 must(s.Close());info,err:=os.Stat(path);must(err);before:=info.Size()
 start:=time.Now();s,err=Open(path);must(err);base:=time.Since(start);must(s.Close());t.Logf("baseline_Open_us=%.3f",float64(base.Nanoseconds())/1000)
 t.Setenv("WORLD103_PROVISION","1")
 start=time.Now();s,err=Open(path);must(err)
 t.Logf("prototype_writable_Open_us=%.3f",float64(time.Since(start).Nanoseconds())/1000)
 measure("indexed");must(s.Close());info,err=os.Stat(path);must(err);t.Logf("db_before=%d db_after=%d delta=%d",before,info.Size(),info.Size()-before)
}
```

<!-- probe-provision -->
```go
func row103Provision(db *sql.DB) error {
 ctx,cancel:=context.WithTimeout(context.Background(),lookupIndexProvisionDeadline);defer cancel()
 specs:=[][4]string{{"log_entries_by_transition_ref","log_entries","transition_ref","entry_index"},{"log_entries_by_transition_fn_ref","log_entries","transition_fn_ref","entry_index"},{"log_entries_by_interpreter_ref","log_entries","interpreter_ref","entry_index"},{"worlds_by_state_root","worlds","state_root","world_ref"}}
 for _,s:=range specs {
  if _,err:=db.ExecContext(ctx,"CREATE INDEX IF NOT EXISTS "+s[0]+" ON "+s[1]+"("+s[2]+","+s[3]+")");err!=nil{return err}
  rows,err:=db.QueryContext(ctx,"PRAGMA index_list('"+s[1]+"')");if err!=nil{return err};found:=false
  for rows.Next(){var seq,u,p int;var name,origin string;if err:=rows.Scan(&seq,&name,&u,&origin,&p);err!=nil{rows.Close();return err};if name==s[0]&&p==0{found=true}}
  err=rows.Err();rows.Close();if err!=nil{return err};if !found{return fmt.Errorf("missing %s",s[0])}
  rows,err=db.QueryContext(ctx,"PRAGMA index_xinfo('"+s[0]+"')");if err!=nil{return err};keys:=0;valid:=true
  for rows.Next(){var seq,cid,desc,key int;var name sql.NullString;var coll string;if err:=rows.Scan(&seq,&cid,&name,&desc,&coll,&key);err!=nil{rows.Close();return err};if key==0{continue};if keys>=2||seq!=keys||name.String!=s[keys+2]||desc!=0||coll!="BINARY"{valid=false};keys++}
  err=rows.Err();rows.Close();if err!=nil{return err};if !valid||keys!=2{return fmt.Errorf("incompatible %s",s[0])}
 };return nil
}
```

# w-object-lookup-by-semantic-id — Find an object by semanticId without a log scan (row 96)

**Status**: Planned — design + prototype (iteration 201, designer `claude:claude-opus-5-5`). The prototype is in this worktree at its real paths (see "Prototype manifest"). It is a measured draft, not a landing.
**Item**: queue row 96 of `design_docs/world-mission.md` (`w-object-lookup-by-semantic-id`). The row closes finding F-1 of `design_docs/verification/w-1-0-value-demonstration.md` and rules on F-2.
**Clauses**: clause-5 (*"on ≥3 REAL 'why did X happen' questions, a provenance walk yields the verified answer in ≤5 minutes each, where the pre-World method was grep/log archaeology"*). The route is read-only and adds no effect, so no other clause is touched.
**Estimate**: ~1d as the row says. Three milestones of ≤150 production lines each (M1 store query + index, M2 route, M3 frozen-table extension: sketch, CLI verb, counts). They are sized from the prototype's measured lines (V20).
**Verified against**: dev `c30b967` (worktree `.wt-world-iter201-design`, branch `design/w-object-lookup-by-semantic-id`). Every claim of the form "the codebase does X" has a Verification Log row with its command and observed output. Pinned interpreter: `$HOME/.pinned-ailang/ailang` (v0.41.0, V0). Machine: Apple M4 Max, macOS 26.6.2 (V0).

---

## (a) Problem (measured at `c30b967`)

Row 92's value-demonstration walk found each incident **by linear scan**. The locate step of every
walk in `walk-transcript.md` (lines 237, 279, 352) runs:

```
for r in $(ailang-worldd log range --from 0 | jq transitionRef); do ailang-worldd object get $r; done
```

and then reads the `semanticId` of each result until it matches `world/mission/incident/<name>`.
Three facts force that shape:

1. `GET /v1/objects/{ref}` accepts a HashRef only. `parseRef` rejects any other text with 400 (V3).
2. The only name-keyed read is `GET /v1/registry/{name...}`. It reads `epoch_registry_heads`, and
   `POST /v1/commit` cannot write that table. Its only writers are `SetRegistryHead` /
   `CompareAndSetRegistryHead`, called by registry bootstrap, the broker and `transitionreg` (V4,
   F-2).
3. `objects.semantic_id` has no index. The schema declares **0** `CREATE INDEX` statements (V5).

The scan costs **one HTTP call per log entry**, plus one per 100-entry page. The CLI walk also
spawns **one process per call**. At row 92's 4 entries the scan was invisible. The measured cost at
mission scale is in the table below (V16–V19).

**What "realistic size" means here.** No World store on this machine has a meaningful log. The
largest real store, `~/.ailang/world/world.db`, holds **0** log entries and **37** objects, and it
is schema v2, which the current binary refuses (V1). The other four stores are empty v3 stores (V1).
So N cannot come from a real store. It is derived from the scale the bar implies:

- The mission has run **201 iterations**. Its repo holds **556** commits on `dev` (V2).
- If World recorded the mission, each iteration would commit its design, plan, milestones, judge
  verdict and record as transitions, plus their journal objects. That is conservatively ~50 log
  entries per iteration, or **~10,000 entries** for the mission so far.
- The walk is measured at **N = 1,000** (small), **N = 10,000** (mission scale, the headline) and
  **N = 50,000** (5× headroom). The incident sits at the **last** index, which is the worst case
  for a scan from 0.

**Measured walk time (the locate step plus the incident read, the same walk shape as row 92).**

| N (log entries) | method | transport | p50 | max | runs | calls per walk |
|---|---|---|---|---|---|---|
| **10,000** | **scan (today)** | **CLI, one process per call — the transcript's shape** | **136.8 s** | **165.0 s** | 3 | 10,104 |
| **10,000** | **route (this doc)** | **CLI** | **0.101 s** | **0.129 s** | 3 | 5 |
| 1,000 | scan | in-process HTTP client | 75.6 ms | 101.2 ms | 5 | 1,014 |
| 1,000 | route | in-process HTTP client | 0.34 ms | 0.65 ms | 5 | 5 |
| 10,000 | scan | in-process HTTP client | 1.10 s | 1.46 s | 5 | 10,104 |
| 10,000 | route | in-process HTTP client | 0.27 ms | 0.33 ms | 5 | 5 |
| 50,000 | scan | in-process HTTP client | 6.57 s | 9.86 s | 5 | 50,504 |
| 50,000 | route | in-process HTTP client | 1.11 ms | 1.59 ms | 5 | 5 |

**Reading it honestly.** The HTTP cost of the scan is small: 1.1 s at 10,000 entries. What
threatens the bar is **the walk as it is actually run**, one CLI process per `object get`. That
costs ~13.7 ms per entry (136.8 s / 10,000). A 10,000-entry log already spends **46%** of the
5-minute budget on locating the incident, before any reasoning, and the scan crosses 300 s at about
**22,000 entries** (300 / 0.0137). The route makes the locate step **one call**, independent of N:
the CLI walk drops from 136.8 s to 0.101 s (~1,350×). The HTTP-client route time grows only from
0.27 ms to 1.11 ms between N=10k and 50k. Both walks returned the same incident hash and verified
all three sources by re-hashing their payloads. The full commands are in §Measurements
(V16–V19).

## (b) Design

### B1. The route: a ninth frozen `/v1` route, not an unversioned side route

```
GET /v1/objects/by-semantic-id/{name...}?after=<hashref>&limit=<n>
```

**Decision: (i) — extend the frozen v1 machine table.** This doc *is* the extension that
`w-effect-broker-m3.md:675` requires ("a new surface must extend the frozen table through its own
doc + quorum"). The argument:

- The frozen table exists so machine clients have a versioned JSON contract. A walker is a machine
  client: clause 5's walk is driven by agents and scripts, and row 92's walk already ran as shell
  loops over JSON. A JSON lookup route outside the table would be exactly the unversioned machine
  contract the table forbids.
- The two precedents outside the table do not fit this route. `GET /workbench` is left out
  *because* its HTML "may evolve" and it is a human renderer (`daemon.go:630-631`, V6). The A2A
  projection routes are a different protocol (`daemon.go:647`, V6). This route is neither: it
  returns the same `objectResponse` JSON the frozen `GET /v1/objects/{ref}` returns.
- `w-worldd-m2.md` requires every frozen route to be reachable via a CLI verb, 1:1 (V7). So the
  route gets the verb `ailang-worldd object find <semanticId> [--after <hashref>] [--limit N]`.

**What the extension changes (M3):**

- The canonical frozen table is `routes()` in `design_docs/sketches/worlddapi.ail` (V8). M3 adds
  `{ method: "GET", path: "/v1/objects/by-semantic-id/{name}" }` to it. This is the one `.ail`
  change: a data row in an existing checked list. No new function, contract or type is added, and
  `scripts/verify_ail.sh` re-checks the sketch (V8).
- **Pre-existing drift, fixed in the same edit.** `routes()` lists **7** routes, but
  `daemon.go:630` says "eight … (seven GET, one POST)". The sketch has never listed
  `GET /v1/log` (range) (V8, V6). M3 adds both rows, so the sketch and the comment agree on
  **nine** (eight GET, one POST). The alternative is a doc that claims "nine" while the sketch
  still says seven. The fix is one line, and M3 edits that list anyway.
- `daemon.go:630` comment "eight … (seven GET, one POST)" becomes "nine … (eight GET, one POST)".
  The comment on `isProtected` ("The eight GET routes pass through unauthenticated", daemon.go
  :657) and `session_middleware_test.go:200` ("the eight GET routes") become "nine"/"eight GET".
  **No test asserts the route count as a number** (V9), so nothing pins "eight" mechanically. That
  is a gap; AC-9 closes it with a test that enumerates the mux's `/v1` patterns against `routes()`.
- The route is **read-only**. It is registered with a `GET` method pattern, so any other method is a
  405 from the mux (the existing `TestGETRoutesRejectOtherMethods` pattern, V10). It is not in
  `isProtected`, the same as the other GET routes: that is residual R1 of `w-session-authority`,
  owned there.

**Mux coexistence.** `GET /v1/objects/{ref}` matches exactly three segments. The new pattern needs
four or more, so neither shadows the other. `/v1/objects/by-semantic-id` (no trailing name) matches
`{ref}` and is a 400 from `parseRef`, not a lookup (AC-6). The multi-segment `{name...}` wildcard is
required for the same reason as the registry route: semantic IDs contain slashes, for example
`world/mission/incident/iter171-index-row` (V11). The CLI escapes each segment separately, the same
as `runRegistry` (V7).

### B2. F-2 is moot for the walk; row 96 does not solve it

F-2 says a walker cannot register a name for an incident via commit. The walk never needed a
*registry* name. It needed to find the object whose `semanticId` is the name, and every incident
already carries one (V11). A semanticId lookup answers the locate step directly. A commit-settable
registry would add a second, mutable naming authority (a head pointer, needing compare-and-set
policy) to answer a question the immutable envelope already answers. So row 96 **does not** make
`epoch_registry_heads` commit-settable. F-2 stays a finding. If a mutable "current incident X"
pointer is ever needed, that is a policy question and needs its own row (Residual G-2, NEW ROW
NEEDED only if a walk ever needs it — none does today).

### B3. Store query

A new store method sits on the read seam, mirroring `ScanUnreadableWorlds` (keyset on a TEXT
primary key, never OFFSET, never rowid; V12):

```go
// ObjectsBySemanticID returns up to limit objects whose semantic_id equals id
// and whose hash_ref sorts strictly after after, ascending by hash_ref.
func (s *Store) ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]Object, error)

SELECT hash_ref, interface_hash_ref, semantic_id, provenance
  FROM objects
 WHERE semantic_id = ? AND hash_ref > ?
 ORDER BY hash_ref
 LIMIT ?
```

- **Returns the existing `store.Object` with `Payload` nil** (no new type). The store enforces its
  own bound, `1 ≤ limit ≤ MaxSemanticIDPage = 500`, with the existing `InvalidLimitError`, the
  same kernel-owned-bound idiom as `ScanUnreadableLog` (V12).
- **No payload column.** The query never reads `payload`, so a page's allocation is bounded by
  metadata only (Decision 7 of `w-worldd-m2`). The walker fetches the payload of the one object it
  wants with the existing `object get --payload`.
- It runs under the caller's `ctx` (`QueryContext`), which is the existing read-deadline seam (B6).
- It parses every hash column with `hashref.Parse`. A malformed stored ref is an error (500), the
  same as `GetObject` (V12).
- `readStore` gains a sixth method. The seam comment's "five distinct getters" becomes six
  (`daemon.go:355-372`, V13).

### B4. Index

```sql
CREATE INDEX IF NOT EXISTS objects_by_semantic_id ON objects(semantic_id, hash_ref);
```

The composite key `(semantic_id, hash_ref)` makes the query a single index range scan in the
requested order. SQLite needs no sort step, and `EXPLAIN QUERY PLAN` shows
`SEARCH objects USING INDEX objects_by_semantic_id (semantic_id=? AND hash_ref>?)` (V14). The index is not covering, so the envelope columns come from the table, but only for the ≤500
rows returned. The one-column alternative `ON objects(semantic_id)` is measured worse: the plan
adds `USE TEMP B-TREE FOR ORDER BY`, a sort of every match. That matters for a hot, non-unique
id such as `world/journal-intent/v1`, which has one object per commit (MU7).

**What the schema change costs (measured, V15):**

- **No `user_version` bump.** `enforceSchemaVersion` refuses any store below the binary's version
  as legacy and has no migration path (`store.go:348,371`, V15). Bumping 3→4 would make every
  existing v3 store unopenable, which is too high a price for an access path. An index is not
  semantic state. The query returns identical rows with or without it (AC-8 proves this by
  dropping the index and re-querying). A v3 store with the index and a v3 store without it are
  therefore the same schema version in every observable way. Both binary directions work, because
  SQLite maintains an index whatever binary writes the table.
- **Existing v3 stores get the index on their next writable open.** At the current version a
  writable `Open` re-executes the whole embedded `schema.sql` (`store.go:353-356`, V15), so
  `CREATE INDEX IF NOT EXISTS` builds it once. The one-time build cost at 50,000 objects is in
  §Measurements (V19). A read-only handle does not apply the schema. On a store that has never
  been opened writable since the upgrade, the query falls back to a table scan and stays correct.
- **Pinned-text tests.** `TestSchemaDDLMatchesCanonicalManifest` compares `type='table'` rows only
  (`journal_test.go:969-972`, V15), so an index does not red it. The `schemaV{1,2,3}SQL` ledger
  pins *tables* per version and is left untouched. M1 adds a new test that pins the index DDL on
  fresh and reopened stores (AC-7).
- **The row-103 "0 indexes" measurement becomes 1.** Row 103's text is historical, so no edit
  is needed, but its designer must re-measure (Conflict Surface).

**Does not foreclose rows 103/104.** Both rows need *new tables or columns* (a reverse-reference
map; a commit-membership link). Those are real semantic schema changes, and they need the version
decision this doc avoids. This index adds no column and no table, reserves no name those rows need,
and uses the `objects` table only as it exists. The ordering contract (B5) names hash order
explicitly, so row 104 can add a chronological order later as a new, opt-in query parameter.

### B5. Ordering, limit, cap, pagination

- **Ordering key: `hash_ref` ascending (byte order, SQLite `BINARY` collation).** It is
  deterministic because `hash_ref` is the table's PRIMARY KEY, so it is unique and a total order.
  It is stable because objects are immutable and never deleted (`PutObject` is
  `INSERT OR IGNORE`, and no production `DELETE FROM objects` exists, V12).
- **Why not insertion order?** `objects` has no insertion-order column (5 columns, V5). The implicit
  `rowid` is not a stable key. `objects` has a `TEXT PRIMARY KEY`, not an `INTEGER PRIMARY KEY`,
  and SQLite documents that `VACUUM` may renumber the rowids of such tables. The repo already
  avoids rowid order for this reason (`scan.go:91` "never depends on OFFSET or SQLite rowids",
  V12). **Why not the committing entry's log index?** The store records no object→entry link.
  That link *is* row 104. The limitation is recorded as Residual G-1, owned by row 104.
- **Limit:** `limit` goes through the existing, Z3-mirrored `clampLimit`. Absent or ≤0 gives
  **100**; above 500 is clamped to **500**, which is the **hard cap** (`handlers.go:240-250`, V10).
  A non-integer `limit` is a 400, the same as `GET /v1/log`.
- **Pagination:** keyset. `after` is a canonical HashRef (parsed with `parseRef`; malformed gives
  400). The response carries `next` = the last item's hash **iff the page is full**
  (`len(items) == limit`). A client pages with `?after=<next>` until `next` is absent. A full last
  page yields one extra, empty request. That is the standard keyset trade and is deterministic.

### B6. Error mapping (consistent with the existing handlers)

| condition | status | class | precedent |
|---|---|---|---|
| empty name (`/v1/objects/by-semantic-id/`) | 400 | `BadRequest` "semantic id is empty" | `handleRegistry` empty name |
| malformed `after` | 400 | `BadRequest` (parseRef text) | `handleObject` |
| non-integer `limit` | 400 | `BadRequest` "limit must be an integer" | `handleLogRange` |
| **no object has this semanticId** | **200** | `{"items":[]}` (no `next`) | `GET /v1/log?from=<past end>` returns 200 `[]` |
| read deadline exceeded | 503 | `Timeout` via `writeReadTimeout` | every GET handler (`readCtx` + `timedOut`) |
| any other store error | 500 | `Internal` via `writeInternalError` (sanitized) | every GET handler |

**Empty vs unknown.** A list query has no "unknown" arm. Returning 404 for a name nobody has
used would make "no match" look like "no such route" to a generic client, and it would be
indistinguishable from a typo in the path prefix. The route answers "which objects carry this
name": the empty set is a valid answer, the same as a log range past the end. The walker's "not
found" is `items == []`. AC-4 pins both the empty-200 and the distinct 400.

The route uses `d.readCtx(r)` and the `timedOut → writeReadTimeout / else writeInternalError`
branch verbatim. It invents no deadline policy: row 23 / PR #153 own that, and row 111 owns
lock-blocked 500-vs-503. Whatever they change in `timedOut`/`readCtx` applies to this route with
no edit.

### B7. Response shape (reuses `objectResponse`)

```json
{"items":[{"hash":"sha256:…","interfaceHash":"sha256:…","semanticId":"world/mission/incident/iter171-index-row","provenance":"w-prove-1-0-phase-a"}],
 "next":"sha256:…"}
```

Each item is the existing `objectResponse` with `Payload` nil, so the `payload` key is absent,
exactly as in `GET /v1/objects/{ref}` without `?payload=true`. `?payload=true` is **not** honoured
here: 500 unbounded payloads in one response would break Decision 7. `next` is `omitempty`.

### B8. No `.ail` change beyond the route row

The route is host transport over an existing pure predicate (`clampLimit`, mirrored from the sketch)
and a store read. Coding standards put effects at the host boundary. There is no new pure decision
to lift into AILANG, since the ordering and limit rules are a SQL clause plus the existing
`clampLimit`. The only `.ail` edit is the `routes()` data rows (B1).

## (c) Milestones (dependency order; each ≤150 production lines)

| # | milestone | production files | est. prod LOC | depends on |
|---|---|---|---|---|
| M1 | store query + index | `host/store/schema.sql` (+1 stmt), `host/store/objects_by_semantic_id.go` (new) | ~60 | — |
| M2 | read route + seam | `host/daemon/handlers.go` (handler + response type), `host/daemon/daemon.go` (mux line, `readStore` method, comments) | ~60 | M1 |
| M3 | frozen-table extension | `design_docs/sketches/worlddapi.ail` (+2 route rows), `cmd/ailang-worldd/cli.go` + `main.go` usage (`object find`), comment counts in `daemon.go` | ~45 | M2 |

The prototype implements M1 and M2 in full and the M3 CLI verb. It does not implement the M3
sketch edit or AC-9 (see "What the prototype does not do").

## (d) Acceptance criteria + non-vacuity

Every AC names the mutation that turns it red. For refusal and limit branches there is one
mutation per branch.

| AC | statement | test | mutation → observed red |
|---|---|---|---|
| AC-1 | non-unique names: 3 objects sharing one semanticId are all returned; an object with a different id is not | `TestObjectsBySemanticIDNonUniqueOrderedAndPaged` | MU3 (prefix match) → store test red |
| AC-2 | ordering determinism: items are in strictly ascending `hash` order, identical across two calls and independent of insertion order | same test (inserts in reverse-hash order) | MU1 (DESC) → store + route tests red |
| AC-3 | pagination: `after` resumes strictly after the cursor; the union of pages = the full set, no dupes | same test | MU2 (`>=`) → store + route red; MU4 (LIMIT ignored) → store + route red |
| AC-4 | empty vs bad request: unknown name → 200 `items:[]`, no `next`; empty name → 400 | `TestObjectsBySemanticIDRoute` | MU8 (empty-name branch) → route test red |
| AC-5 | cap: `limit=100000` returns ≤500 and `next` set when more exist; default is 100 | `TestObjectsBySemanticIDRouteCapAndDefault` | MU11 (default 500) / MU12 (cap bypass) / MU13 (next always) → cap test red |
| AC-6 | refusals: malformed `after` → 400; non-integer `limit` → 400; `/v1/objects/by-semantic-id` without a name does not reach the lookup | `TestObjectsBySemanticIDRoute` | MU9 (after unparsed) / MU10 (limit error ignored) → route test red, one each |
| AC-7 | the index exists on a fresh store and the plan uses it | `TestObjectsBySemanticIDUsesIndex` | MU6 (index dropped) / MU7 (single-column index) → index test red |
| AC-8 | the index is not semantic: dropping it yields byte-identical results | `TestObjectsBySemanticIDUsesIndex` | the arm runs after DROP INDEX; a result that depended on the index (none today) would red. Recorded as a property, not killed by a mutation |
| AC-9 | (M3) the mux's `/v1` patterns equal `routes()` in the sketch, count 9 | NOT PROTOTYPED | M3 executor must show: delete the mux line → red; delete a sketch row → red |
| AC-12 | the new route is in the shared read-route harness, so the existing deadline (503) and sanitization (500) gates cover it | `seedReadRoutes` + `blockingStore`/`recordingStore`/`failingStore` overrides in `read_deadline_test.go` | MU14, MU15, MU16 |
| AC-10 | store error → 500 sanitized; deadline → 503 | `TestObjectsBySemanticIDRouteStoreErrors` (fault-injecting `readStore` wrapper, the V13 seam's purpose) | MU14 / MU16 → deadline / sanitize tests red |
| AC-11 | walk re-timed at N=10,000: route locate p50 ≤ 1 s and far below the scan | `TestWalkRetimedAtScale` (env-gated, `WORLD_WALK_N`) | reported as measurement, not a gate (see §Measurements) |

## Prototype manifest (files changed in this worktree)

From `git status --short` / `git diff --stat` against `c30b967` (V20). Production lines are
`+` lines, comments included.

| file | status | what | + lines |
|---|---|---|---|
| `host/store/schema.sql` | M | `CREATE INDEX IF NOT EXISTS objects_by_semantic_id ON objects(semantic_id, hash_ref)` + comment | 7 |
| `host/store/objects_by_semantic_id.go` | new | `MaxSemanticIDPage`, the keyset SQL, `(*Store).ObjectsBySemanticID` | 56 |
| `host/daemon/handlers.go` | M | `objectPageResponse`, `handleObjectsBySemanticID` | 62 |
| `host/daemon/daemon.go` | M | mux line; `readStore.ObjectsBySemanticID` | 2 |
| `cmd/ailang-worldd/cli.go` | M | `object find` → `runObjectFind` (per-segment escaping) | 34 |
| `cmd/ailang-worldd/main.go` | M | usage lines for `object find` | 2 |
| **production total** | | | **163** |
| `host/store/objects_by_semantic_id_test.go` | new | AC-1/2/3/7/8 + writable-reopen index build | test |
| `host/daemon/objects_by_semantic_id_test.go` | new | AC-4/5/6, and `TestWalkRetimedAtScale` (AC-11, env-gated) | test |
| `host/daemon/read_deadline_test.go` | M | new route in `seedReadRoutes`; `ObjectsBySemanticID` on `blockingStore`, `recordingStore`, `failingStore` (AC-12) | test |
| `design_docs/planned/w-object-lookup-by-semantic-id.md` | new | this doc | doc |

**What the prototype does not do (M3 remainder):** it does not edit
`design_docs/sketches/worlddapi.ail` `routes()`. It does not update the count comments ("eight" →
"nine") in `daemon.go` / `session_middleware_test.go`, or the `readStore` "five distinct getters"
comment and the matching `read_deadline_test.go` "SIX routes / FIVE getters" comments. It adds no
AC-9 route-table-vs-sketch test and no CLI test for `object find` (the verb was exercised live in
the CLI walk, V18, but has no unit test). The M3 executor owns these.

## Mutations run (all against the prototype's tests; each restored after its run)

Runner: `python3 ~/.ailang/state/world-iter201/designer-scratch/mutate.py` applies one textual
mutation, runs `go test -count=1 -run 'SemanticID' ./host/store/` and/or
`go test -count=1 -run 'SemanticID|TestDaemonReadDeadline|TestInternalErrorsAreSanitized' ./host/daemon/`,
then restores the file from a scratch backup. Output is in `mutations.log` next to it (V21).

| # | mutation | killed by | observed red (first failing line) |
|---|---|---|---|
| MU1 | `ORDER BY hash_ref` → `ORDER BY hash_ref DESC` | store NonUnique…Paged, UsesIndex; route Route, CapAndDefault | `full page = [sha256:f2fc…, sha256:b5c6…, …], want ascending …` |
| MU2 | `hash_ref > ?` → `hash_ref >= ?` | store NonUnique…Paged; route Route, CapAndDefault | `pages concatenated = […duplicates…]`; route `after=… page = {Items:[{Hash:<the cursor itself>…` |
| MU3 | `semantic_id = ?` → `semantic_id GLOB ? \|\| '*'` | store NonUnique…Paged, UsesIndex (daemon tests have no prefix decoy: survives there, killed in store) | `full page = [sha256:134a… (decoy) …]` |
| MU4 | `LIMIT ?` → `LIMIT 500 + 0*?` | store NonUnique…Paged (after tightening, see note); route CapAndDefault | `page after "" has 7 items, want at most the limit 3` |
| MU5 | store limit bound → `if false` | store NonUnique…Paged | `limit 0: err = <nil>, want InvalidLimitError` |
| MU6 | delete the `CREATE INDEX` | store UsesIndex, IndexBuiltOnWritableReopen | `query plan = "SEARCH objects USING INDEX sqlite_autoindex_objects_1 (hash_ref>?)"` |
| MU7 | index `(semantic_id, hash_ref)` → `(semantic_id)` | store UsesIndex | `query plan = "… (semantic_id=?) \| USE TEMP B-TREE FOR ORDER BY"` |
| MU8 | delete the empty-name 400 branch | route Route | `status = 200, want 400; body={"items":[]}` |
| MU9 | `after` used unparsed | route Route | `status = 200, want 400; body={"items":[…` |
| MU10 | `limit` Atoi error ignored | route Route | `status = 200, want 400` |
| MU11 | default page 500 instead of 100 | route CapAndDefault | `default page: 500 items …, want 100` |
| MU12 | cap bypassed for `limit > 500` | route CapAndDefault | `oversized limit: 0 items next="", want the 500 cap` (the store's own bound refuses → 500) |
| MU13 | `next` set on any non-empty page | route Route, CapAndDefault | `json: cannot unmarshal string into … []map` (a `next` key appeared on a non-full page) |
| MU14 | timeout branch disabled (`if false && timedOut`) | TestDaemonReadDeadline/real-store-expired-deadline, /blocking-store | `status = 500, want 503` |
| MU15 | lookup runs under `r.Context()`, not `readCtx` | TestDaemonReadDeadline (both arms) | `status = 200, want 503` |
| MU16 | 500 writes `err.Error()` instead of `writeInternalError` | TestInternalErrorsAreSanitized/read-routes | `the 500 body leaked the store's internal detail "kQ7v-store-detail-9f3c1d82"` |
| MU17 | items carry `Payload: &[]byte{}` | route Route | `item carries a payload key` |
| MU18 | mux line removed | route Route, CapAndDefault; TestDaemonReadDeadline/normal-deadline-answers-200 … | `match status = 404 body=404 page not found` |
| MU19 | **green control**: comment-only edit | — (both packages `ok`) | none, as required |

**18 killed, 0 survived, 1 green control stayed green.** MU3 survives the *daemon* package alone
but is killed in the store package, which is where the predicate lives. **Note on MU4:** on the
first battery MU4 **survived the store package**. The paging loop accepted a 7-item page for
limit 3, because the concatenation was still correct. The test now asserts `len(page) ≤ 3`, and
the re-run kills MU4 in both packages (V21, second block).

## Measurements

Machine: Apple M4 Max, macOS 26.6.2, Go per `go.mod`, SQLite via `modernc.org/sqlite` (in-process)
and `sqlite3` 3.51.0 (CLI, index timing only). The HTTP-client and CLI rows ran at different
times. The **N=50,000 HTTP rows overlapped the mutation battery on the same machine**, and their
max (9.86 s) is inflated by that contention. The N=1,000 and N=10,000 HTTP rows and all CLI rows
ran with nothing else of this designer's running (the CLI rows ran after the battery and the
HTTP walks had finished).

- **Fixture (`seedWalkLog`).** N real `store.Commit`s, each with 2 objects: the transition object
  (`world/mission/transition/<i>`, the entry's `transitionRef`) and a journal object under the
  shared, non-unique id `world/journal-intent/v1`. The **last** entry's transition is the incident
  `world/mission/incident/walk-target`. Its JSON payload lists 3 source objects committed with
  it. Total objects = 2N + 3 (+1 epoch-registry object).
- **The walk (both methods).** (1) locate the incident, (2) `object get <incident> --payload`,
  (3) `object get <src> --payload` for each of the 3 sources, re-hashing each payload against its
  ref. This is row 92's walk shape (`walk-transcript.md:237-246`). Step (1) is either the scan
  (`log range --from k` pages of the default 100, then `object get <transitionRef>` per entry until
  `semanticId` matches) or the route (`object find <id>`, one call).
- **HTTP-client rows:** `TestWalkRetimedAtScale` over a real loopback listener (V16, V17).
- **CLI rows:** `cliwalk.sh` against `ailang-worldd serve --db walk10k.db` (the store the N=10,000
  test run left behind), every call a fresh `ailang-worldd` process, exactly as the transcript's
  loops (V18).
- **One-time index build (B4 cost):** 19 ms on the 20,004-object N=10,000 store and 60 ms on a
  synthetic 100,000-object `objects` table (V19). So the writable-open build is negligible at
  mission scale.
- **Seeding cost** (context only): N=10,000 in 8.9 s and N=50,000 in 63.1 s through
  `store.Commit` (V17).

## (e) Verification Log

Commands ran from the worktree root at `c30b967` (the `git show c30b967:` form reads the
unmodified tree even after the prototype's edits). Scratch = `~/.ailang/state/world-iter201/designer-scratch`.

| # | claim | command | observed |
|---|---|---|---|
| V0 | pinned interpreter; machine | `$HOME/.pinned-ailang/ailang --version \| head -2`; `sysctl -n machdep.cpu.brand_string`; `sw_vers -productVersion` | `AILANG v0.41.0` / `Commit: 24ee108`; `Apple M4 Max`; `26.6.2` |
| V1 | no World store on this machine has a meaningful log | `find ~/.ailang -maxdepth 4 \( -name "*.db" -o -name "*.sqlite" \) -size +0`, then for each World-shaped db `sqlite3 -readonly "file:$f?mode=ro" "PRAGMA user_version; SELECT count(*) FROM log_entries; SELECT count(*) FROM objects; …"` | `~/.ailang/world/world.db`: user_version **2**, log **0**, objects **37** (12 distinct ids, max 6 per id). `world-iter195/exec-m7/world-demo.db`, `iter188-eval/r.db`, `iter188-planner/recon_{p,b}.db`: user_version 3, log 0, objects 0. **Positive control, same command:** the non-World dbs (`brain.db`, `coordinator.db`, …) were found by the same `find`, so the search is not empty-by-construction. |
| V2 | mission size | `git rev-list --count c30b967` | `556` |
| V3 | `GET /v1/objects/{ref}` is hashref-only | `git show c30b967:host/daemon/handlers.go \| grep -n 'parseRef(r.PathValue("ref"), "object ref")'` | `378:` match |
| V4 | registry heads are not commit-settable (F-2) | `git grep -n "SetRegistryHead(\|CompareAndSetRegistryHead(" c30b967 -- '*.go' \| grep -v _test`; `git show c30b967:host/store/store.go \| grep -n epoch_registry_heads` | writers: `broker/approve.go:211`, `registry/registry.go:168`, `transitionreg/transitionreg.go:266`; the only INSERTs into `epoch_registry_heads` are `store.go:623` (SetRegistryHead) and `:715` (CAS). `Commit`'s object insert is `store.go:950` (`INSERT OR IGNORE INTO objects`), with no registry write |
| V5 | 0 indexes, 9 tables, 5 object columns | `git show c30b967:host/store/schema.sql \| grep -c "CREATE INDEX"`; same `\| grep -c "CREATE TABLE"` (positive control) | `0`; `9`; the `objects` DDL has 5 columns (read in full) |
| V6 | the frozen-table comment and the side routes | `git show c30b967:host/daemon/daemon.go \| grep -n "complete frozen v1 machine table\|NOT part of the frozen table\|ADDITIVE: the frozen\|The eight GET routes"` | `630` "eight … (seven GET, one POST)", `632` workbench "NOT part of the frozen table", `647` A2A "ADDITIVE", `657` "The eight GET routes" |
| V7 | CLI 1:1 rule; the extension-by-doc rule | `grep -n "mapping 1:1 onto the frozen route table\|every route in the frozen table is reachable via a CLI verb" design_docs/implemented/w-worldd-m2.md`; `grep -n "extend the frozen table" design_docs/implemented/w-effect-broker-m3.md` | `286`, `459`; `675` |
| V8 | the sketch lists 7 routes and omits `GET /v1/log` | `git show c30b967:design_docs/sketches/worlddapi.ail \| grep -c '{ method: "'`; same `\| grep -n '/v1/log'` | `7`; only `77: { method: "GET", path: "/v1/log/{index}" }` (the positive control for the grep matching log paths) |
| V9 | no test asserts the route count | `git grep -nE 'len\(.*\) != (7\|8)\b\|"eight"\|eight (GET\|/v1)\|seven GET' c30b967 -- '*_test.go'` | one hit, `session_middleware_test.go:200`, a **comment** ("the eight GET routes"). It is the positive control that the regex matches the phrasing; no assertion |
| V10 | `clampLimit` and the 405 test exist | `git show c30b967:host/daemon/handlers.go \| grep -n "^func clampLimit"`; `git show c30b967:host/daemon/handlers_test.go \| grep -n "func TestGETRoutesRejectOtherMethods"` | `240:`; `267:` |
| V11 | incident semanticIds contain slashes; none has an empty segment | `grep -c '"semanticId":"world/mission/incident/' …/walk-transcript.md`; `git grep -nE 'semanticId":"[^"]*//' c30b967 -- design_docs/verification \| wc -l` | `13` (positive); `0` |
| V12 | objects are never deleted; keyset precedent | `git grep -n "DELETE FROM" c30b967 -- 'host/*.go' ':!*_test.go'`; `git show c30b967:host/store/scan.go \| grep -n "never depends on OFFSET"` | the only production DELETE is `store.go:1129 DELETE FROM session_credentials` (the positive control: the grep finds DELETEs), none on `objects`; `92:` |
| V13 | the read seam has five getters | `git show c30b967:host/daemon/daemon.go \| grep -n "five distinct getters\|^type readStore"` | `357`, `366` |
| V14 | the new query uses the index with no sort | `sqlite3 scratch/idx-copy.db "EXPLAIN QUERY PLAN SELECT hash_ref, interface_hash_ref, provenance FROM objects WHERE semantic_id = 'x' AND hash_ref > '' ORDER BY hash_ref LIMIT 100;"` | `SEARCH objects USING INDEX objects_by_semantic_id (semantic_id=? AND hash_ref>?)`; also asserted in-process by `TestObjectsBySemanticIDUsesIndex` |
| V15 | legacy refusal, re-apply at current version, table-only DDL pin | `git show c30b967:host/store/store.go \| grep -n "return &LegacySchemaVersionError\|if version == currentSchemaVersion\|db.Exec(schemaSQL)"`; `git show c30b967:host/store/journal_test.go \| grep -n "WHERE type='table'"` | `348`, `371` (legacy), `353`, `355` (re-apply); `972` |
| V16 | HTTP-client walks at N=1,000 / 50,000 | `WORLD_WALK_N=$n scratch/daemon.test -test.run '^TestWalkRetimedAtScale$' -test.v` (via `scratch/walks.sh`, log `walks.log`) | N=1000: scan p50 75.648667ms max 101.235125ms, 1014 calls; route p50 339.833µs max 654.5µs, 5 calls. N=50000: scan p50 6.5688685s max 9.857619125s; route p50 1.11075ms max 1.590917ms |
| V17 | HTTP-client walk at N=10,000 (store kept) | same, `WORLD_WALK_N=10000 WORLD_WALK_DB=scratch/walk10k.db` | seeded 10000 entries (20003 objects) in 8.928831916s; scan p50 1.103359167s max 1.461967167s, 10104 calls; route p50 273.458µs max 331.875µs, 5 calls; same incident `sha256:5696…d9d6` both ways |
| V18 | CLI walks at N=10,000 | `env -u AILANG_REGISTRY_API_KEY perl -e 'alarm shift; exec @ARGV' 1500 scratch/ailang-worldd serve --db scratch/walk10k.db --bind 127.0.0.1:7699`; then `scratch/cliwalk.sh` (log `cliwalk.log`) | route: 0.101 / 0.129 / 0.094 s; scan: 164.979 / 136.808 / 112.314 s; every run `incident=sha256:569623ab…d9d6` |
| V19 | one-time index build | `sqlite3 -readonly … ".backup scratch/idx-copy.db"`, `DROP INDEX`, then `.timer on` + `CREATE INDEX objects_by_semantic_id …`; same on a 100,000-row synthetic `objects` | `20004` objects: `Run Time: real 0.019`; `100000`: `real 0.060` |
| V20 | prototype size | `git diff --stat`; `wc -l host/store/objects_by_semantic_id.go`; `git diff <file> \| grep -c '^+[^+]'` | 6 tracked files, +122/−1 (incl. test file); new store file 56 lines; handlers +62, cli +34, daemon.go +2, main.go +2, schema +7 |
| V21 | mutation battery | `python3 scratch/mutate.py` then `python3 scratch/mutate.py MU4 MU19` | as §Mutations; `mutations.log` |
| V22 | gates on the prototype | `go vet ./...`; `go test -count=1 -run '^$' ./...`; `go test -race -count=1 ./host/store/ ./host/daemon/ ./cmd/ailang-worldd/`; `./scripts/verify_ail.sh`; `go test -count=1 ./...` (all with `AILANG_BIN` pinned, `GOCACHE`/`TMPDIR` outside the tree) | vet clean; compile fence clean; gofmt clean on touched files; `-race`: store `ok` 15.8s, daemon `ok` 17.9s, cmd/ailang-worldd `ok` 3.0s; `verify_ail.sh`: "verify gate PASSED: 16 required identities verified, 40 named tests pass". Full suite: **23 ok / 1 FAIL**. The FAIL is `host/broker` `TestOrdinaryOverflowCarriesNoESRCH` ("overflow kill: operation not permitted"), an EPERM on the kill path in a package this prototype does not touch (`git diff --stat c30b967 -- host/broker` is empty). **Re-run alone: `ok host/broker 66.481s`.** It is the same EPERM-under-load family as row 115 (capsule). I did not find a broker-specific row, so it is recorded here and not claimed as a known flake |

## (f) Conflict Surface

| file | this row's change | open rows touching the same file |
|---|---|---|
| `host/store/schema.sql` | +1 `CREATE INDEX` | **103** (adds indexes + reverse read), **104** (commit-membership schema change). Additive, with no shared name. Both will need a version decision this doc does not take (B4). |
| `host/store/objects_by_semantic_id.go` (new) | query | none |
| `host/daemon/handlers.go` | new handler + response type | **111** (`timedOut` / 500-vs-503), **23 / PR #153** (read deadline policy). This route calls their helpers unchanged, so their edits flow through (B6). Merge risk is textual only. |
| `host/daemon/daemon.go` | mux line, `readStore` method, comment counts | **23 / PR #153** (readCtx/deadline wiring), **103** (a reverse read also extends `readStore`), **106** (coordinator wiring in `New`). Adjacent lines in the seam interface, so whichever lands second rebases one line. |
| `cmd/ailang-worldd/cli.go`, `main.go` | `object find` verb | none known |
| `design_docs/sketches/worlddapi.ail` | +2 `routes()` rows | none known |

## (g) Residuals

- **G-1 — no chronological order.** Results are in hash order because no object→commit link
  exists. A walker wanting "the newest object named X" among many cannot ask for it. Owner:
  **row 104** (`w-store-commit-object-membership`), which can add `order=commit` as an opt-in
  parameter.
- **G-2 — F-2 (commit-settable registry names) stays a finding.** It is moot for the walk (B2).
  Owner: NEW ROW NEEDED only if a mutable name pointer is ever required. No current walk requires
  one.
- **G-3 — reverse lookups** ("which entries reference this incident?") still require a scan.
  Owner: **row 103**.
- **G-4 — unauthenticated GET.** Inherited residual R1 of `w-session-authority`; the route joins
  the other GETs. Owner: the existing R1 follow-up (not re-filed).
- **G-5 — semantic IDs with empty segments or `.`/`..` segments** (`a//b`) are cleaned or
  redirected by `ServeMux` before routing, exactly as for `GET /v1/registry/{name...}`. That
  limitation is inherited, not introduced. No such ID exists in any store or fixture (V11). Owner:
  NEW ROW NEEDED if one ever appears.
- **G-6 — the walk's CLI-per-call spawn cost** dominates the *remaining* walk steps (reading an
  incident's `sources`). It is linear in the number of sources, which is small (≤5 in row 92). No
  owner is needed; recorded for honesty.

## (h) Decision for Mark

None. Adding a frozen route is a technical extension governed by the existing
doc-plus-quorum rule, and the quorum on this doc is that ratification. Every other choice here
(hash order, empty-200, no version bump) is technical and argued above.

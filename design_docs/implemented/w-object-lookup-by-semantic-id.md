# w-object-lookup-by-semantic-id — Find an object by semanticId without a log scan (row 96)

**Status**: IMPLEMENTED 2026-09-27 (iteration 201) — PR #156 squash `e859501`, remote CI 2/2 green on the merge; judged PASS 95/100 zero blocking (sonnet). Design history: **REVISION 2 (final; narrow-refinement carve-out after quorum round 2, reviewer fixes applied verbatim, no r3 quorum)**. Round 1 and round 2 were each BLOCKED 2/2 present; see "Quorum log". Design + prototype (iteration 201, designer `claude:claude-opus-5-5`). The r0 prototype was measured in this worktree and then moved out by the controller to `~/.ailang/state/world-iter201/prototype/` (see "Prototype manifest"). The r1/r2 changes (the index guard and its strict compatibility rule, dropping `next`, the `{name...}` sketch spelling) are **designed, not prototyped**. It is a measured draft, not a landing.
**Item**: queue row 96 of `design_docs/world-mission.md` (`w-object-lookup-by-semantic-id`). The row closes finding F-1 of `design_docs/verification/w-1-0-value-demonstration.md` and rules on F-2.
**Clauses**: clause-5 (*"on ≥3 REAL 'why did X happen' questions, a provenance walk yields the verified answer in ≤5 minutes each, where the pre-World method was grep/log archaeology"*). The route is read-only and adds no effect, so no other clause is touched.
**Estimate**: ~1d as the row says. Three milestones of ≤150 production lines each (M1 store query + index provisioning + index guard, M2 route + 503 guard class, M3 frozen-table extension: sketch, CLI verb, counts). They are sized from the r0 prototype's measured lines (V20), plus a labelled r1 estimate delta (§c).
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
  `{ method: "GET", path: "/v1/objects/by-semantic-id/{name...}" }` to it (r2: the path is spelled
  exactly as the mux pattern, `{name...}`, so AC-9's parity test compares like with like). **The
  same fix applies to the sibling registry row.** `worlddapi.ail:78` lists
  `{ method: "GET", path: "/v1/registry/{name}" }`, while the mux pattern is
  `GET /v1/registry/{name...}` (`daemon.go:643`, V28). M3 changes line 78 to `{name...}` too.
  Otherwise AC-9 would red on a pre-existing spelling drift, not on a real route difference.
  **Exactly what row 96 adds**
  (corrected in r1; r0 overclaimed "no new function, contract or type"):
  - **one route** (the mux pattern, plus its `routes()` row, plus the `object find` CLI verb);
  - **one Go response struct**, `objectPageResponse{Items []objectResponse}`. It has the same
    `{"items":[…]}` envelope as `logRangeResponse` (`handlers.go:63-65`, V24) and no other field;
  - **one store method**, `ObjectsBySemanticID`, plus one `readStore` seam method, plus one store
    error type, `LookupIndexUnavailableError` (B4);
  - **one index**, `objects_by_semantic_id`, provisioned at writable open (B4);
  - **one sketch `ApiError` variant**, `LookupIndexUnavailable(string)`, with its `httpStatus` arm
    `=> 503` and one test vector (B6). It is the sketch-side authority for the guard's status,
    mirrored the same way `TestTimeoutStatusMirrorsSketch` mirrors `Timeout` today.

  No new contract is added. `scripts/verify_ail.sh` re-checks the sketch (V8).
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
  same kernel-owned-bound idiom as `ScanUnreadableLog` (V27).
- **No payload column.** The query never reads `payload`, so a page's allocation is bounded by
  metadata only (Decision 7 of `w-worldd-m2`). The walker fetches the payload of the one object it
  wants with the existing `object get --payload`.
- It runs under the caller's `ctx` (`QueryContext`), which is the existing read-deadline seam (B6).
- **Index guard, before any query.** If the handle's cached index verdict (B4) is "absent or
  incompatible", the method returns `*LookupIndexUnavailableError` **without executing the lookup
  query**. It never falls back to a table scan.
- It parses every hash column with `hashref.Parse`. A malformed stored ref is an error (500), the
  same as `GetObject` (V27).
- `readStore` gains a sixth method. The seam comment's "five distinct getters" becomes six
  (`daemon.go:355-372`, V13).

### B4. Index

```sql
CREATE INDEX IF NOT EXISTS objects_by_semantic_id ON objects(semantic_id, hash_ref);
```

**r1: where this DDL lives.** It is **not** in `schema.sql` (the r0 prototype put it there). A
writable open re-executes `schema.sql` with `db.Exec`, which has no deadline (`store.go:355`,
V15). Instead, M1 adds a provisioning step, `provisionLookupIndex`, that runs inside writable
`Open` after `enforceSchemaVersion` succeeds, for fresh and existing stores alike. It runs as
`ExecContext` under `context.WithTimeout(lookupIndexProvisionDeadline)`, a named constant of
**30 s**. That is ~500× the measured build time of 60 ms for 100,000 objects (V19). `schema.sql` gets
a one-line comment pointing at the step, and no DDL.

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
  semantic state. At store level the query returns identical rows with or without it (AC-8
  proves this by dropping the index and re-querying the *unguarded* SQL). **AC-8 is retained only
  as that store-level result-equivalence property. It is not permission for a serving
  fallback:** the served path is guarded (above). A v3 store with the index and a v3 store without it are
  therefore the same schema version in every observable way. Both binary directions work, because
  SQLite maintains an index whatever binary writes the table.
- **The index is a precondition of the lookup; there is no serving fallback (r1, quorum
  objection 1 applied verbatim).** The lookup requires the composite `objects_by_semantic_id`
  index. Before enabling the route, verify its definition. If the index is absent or incompatible
  on a read-only store, return an explicit 503 `LookupIndexUnavailable` response with remediation
  instructions; do not execute the lookup query. Provisioning occurs through a writable open under
  an explicit bounded deadline, and provisioning failures are surfaced rather than ignored.

  How that is realized:
  - **Where the check lives: once per handle, at `Open` / `OpenReadOnly`, cached on the `Store`.**
    The verdict is stored in an immutable field set before `Open` returns.
  - **Compatibility rule (r2, quorum objection applied verbatim):** Accept only a non-partial
    index on objects whose two key columns, inspected through PRAGMA index_list and PRAGMA
    index_xinfo, are semantic_id and hash_ref in that order with BINARY collation and the declared
    ascending ordering. All other definitions are incompatible; writable Open fails without
    modifying them, and read-only lookup refuses before executing SQL.

    Mechanically (probed, V29): `verifyLookupIndex` runs `PRAGMA index_list('objects')` and requires
    a row named `objects_by_semantic_id` with `partial = 0`. Listing *objects'* indexes is what
    establishes "on objects", so a same-named index on another table is absent here. It then runs
    `PRAGMA index_xinfo('objects_by_semantic_id')`, keeps the rows with `key = 1`, and requires
    exactly two: `(seqno 0, name semantic_id, coll BINARY, desc 0)` and
    `(seqno 1, name hash_ref, coll BINARY, desc 0)`. **The probe shows why both PRAGMAs are needed:**
    a partial index `… WHERE semantic_id <> ''` has an `index_xinfo` byte-identical to the accepted
    one. Only `index_list.partial = 1` tells them apart (V29). A `COLLATE NOCASE` key shows
    `coll = NOCASE`, and a `hash_ref DESC` key shows `desc = 1` (V29). The rule ignores the
    `key = 0` auxiliary rowid row that `index_xinfo` always appends. **Why not per call:** a per-call check would add one `sqlite_master`/PRAGMA
    round trip to every request, inside the read deadline, to detect a change that cannot happen
    under a writable handle. The store is single-writer (writer lock), and no code drops indexes
    (V12's DELETE census; there is no `DROP INDEX` in production code, V27). The one stale case is a
    read-only handle whose store another process provisions later. It keeps answering 503 until it
    is reopened, which is explicit, not silent. The cost is one PRAGMA per open and zero per call.
  - **Writable open:** provision (bounded, above), then verify. A provisioning error, a deadline
    expiry, or a post-provision verdict of "incompatible" makes `Open` **return an error** (the
    daemon then fails startup at its existing `store.Open` stage). It is never ignored, and an
    incompatible same-named index is never silently dropped and rebuilt. So a writable handle
    always carries a verified index, or does not exist.
  - **Read-only open:** verify only (a read-only handle must not attempt DDL, `store.go:266-268`).
    It caches "absent/incompatible" and opens anyway, so the other reads keep working, and only
    this lookup answers 503.
  - **The 503 branch is a guard, not an expected path.** Production has **0** `OpenReadOnly`
    callers. The daemon opens its store with `store.Open(cfg.DBPath)` (`daemon.go:470`), as do
    `cmd/ailang-worldd/session.go:145,211` and the four `cmd/world-publish` sites (V23). So every
    production handle is writable, and the index is provisioned and verified at daemon startup
    before the route can be served. The guard exists so that a future read-only consumer (for
    example a read-replica workbench) cannot silently degrade to an O(N) scan.
  - **Remediation text (fixed, host-detail-free, like the 500 sanitizer):** `"semantic-id lookup
    index is absent or incompatible on this read-only store; open the store writable once (for
    example start ailang-worldd serve on it) to provision objects_by_semantic_id"`.
- **Pinned-text tests.** `TestSchemaDDLMatchesCanonicalManifest` compares `type='table'` rows only
  (`journal_test.go:969-972`, V15), so an index does not red it. The `schemaV{1,2,3}SQL` ledger
  pins *tables* per version and is left untouched. Because r1 keeps the DDL out of `schema.sql`,
  the `schemaSQL`-applying fixtures in `schema_version_test.go` are unaffected too. M1 adds new
  tests that pin the index on fresh and reopened stores (AC-7) and the guard (AC-13..AC-15).
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
- **Pagination (r1: no `next` field; quorum objection 2's paging rule, verbatim):** keyset. `after`
  is a canonical HashRef (parsed with `parseRef`; malformed gives 400). The client must "read the
  `hash` of the last item to construct `?after=`, terminating when len(items) < limit or the
  returned page is empty". This needs no server-issued token, because the cursor *is* the ordering
  key of the last item. A full last page costs one extra, empty request, which is the standard
  keyset trade and is deterministic. The client must know the effective limit. It is `limit` if
  1..500 was sent, 500 if more was sent, and 100 if none. That is the same `clampLimit` rule
  `GET /v1/log` clients already rely on.

### B6. Error mapping (consistent with the existing handlers)

| condition | status | class | precedent |
|---|---|---|---|
| empty name (`/v1/objects/by-semantic-id/`) | 400 | `BadRequest` "semantic id is empty" | `handleRegistry` empty name |
| malformed `after` | 400 | `BadRequest` (parseRef text) | `handleObject` |
| non-integer `limit` | 400 | `BadRequest` "limit must be an integer" | `handleLogRange` |
| **no object has this semanticId** | **200** | `{"items":[]}` | `GET /v1/log?from=<past end>` returns 200 `{"items":[]}` |
| lookup index absent/incompatible (read-only handle only; unreachable in production, V23) | **503** | **`LookupIndexUnavailable`** with the fixed remediation text (B4); the lookup query is **not** executed | new sketch arm `LookupIndexUnavailable(_) => 503` (B1) |
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
{"items":[{"hash":"sha256:…","interfaceHash":"sha256:…","semanticId":"world/mission/incident/iter171-index-row","provenance":"w-prove-1-0-phase-a"}]}
```

The envelope is exactly `GET /v1/log`'s: `logRangeResponse` is `{"items":[…]}` with no other
field (`handlers.go:63-65`, `:492`, V24). r1 drops the r0 `next` field (quorum objection 2). Each
item is the existing `objectResponse` with `Payload` nil, so the `payload` key is absent, exactly
as in `GET /v1/objects/{ref}` without `?payload=true`. `?payload=true` is **not** honoured here:
500 unbounded payloads in one response would break Decision 7.

### B8. `.ail` changes: sketch data only

The route is host transport over an existing pure predicate (`clampLimit`, mirrored from the sketch)
and a store read. Coding standards put effects at the host boundary. There is no new pure decision
to lift into AILANG, since the ordering and limit rules are a SQL clause plus the existing
`clampLimit`. The `.ail` edits are all in the frozen sketch: the two `routes()` data rows (B1), and
the `ApiError` variant `LookupIndexUnavailable(string)` with its `httpStatus` arm and one test
vector (B1, B6). The index check itself is host I/O (a PRAGMA), so it belongs in Go.

## (c) Milestones (dependency order; each ≤150 production lines)

| # | milestone | production files | r0 measured prod LOC (V20) | **r1 delta (ESTIMATE, not prototyped)** | r1 est. total | depends on |
|---|---|---|---|---|---|---|
| M1 | store query + bounded index provisioning + cached index guard | `host/store/objects_by_semantic_id.go` (new), `host/store/store.go` (call `provisionLookupIndex`/`verifyLookupIndex` from `Open`/`OpenReadOnly`; one `Store` field), `host/store/schema.sql` (comment only) | 63 (56 store file + 7 schema) | **+~45**: `LookupIndexUnavailableError` ~10, `provisionLookupIndex` with deadline ~12, `verifyLookupIndex` ~18 (r1) **→ ~30 (r2: `index_list` partial check + `index_xinfo` key/coll/desc checks)**, guard + test hook ~5; **−6** schema DDL moved out | **~114** (r2; was ~102 in r1) | — |
| M2 | read route + seam + 503 guard class | `host/daemon/handlers.go` (handler, `objectPageResponse`, 503 mapping), `host/daemon/daemon.go` (mux line, `readStore` method), `design_docs/sketches/worlddapi.ail` (ADT variant + arm + vector) | 64 (62 handlers + 2 daemon.go) | **+~8** (`errors.As` → `writeAPIError(…"LookupIndexUnavailable"…, 503)`); **−4** (`next` field + assignment) | **~68** Go, + ~3 `.ail` lines | M1 |
| M3 | frozen-table extension | `design_docs/sketches/worlddapi.ail` (+2 route rows), `cmd/ailang-worldd/cli.go` + `main.go` usage (`object find`), comment counts in `daemon.go` | 36 (34 cli + 2 main) | **+~6** comment-count edits; 0 for `next` (the CLI never read it); r2: `worlddapi.ail:78` `{name}` → `{name...}` is one `.ail` line, not Go | **~42** Go, + ~3 `.ail` lines | M2 |

Every milestone stays ≤150. The r1 deltas are the designer's line estimates for code this round did
**not** write; the M1/M2 executor must report the measured numbers. Moving the sketch `ApiError`
arm into M2 (not M3) is deliberate. The handler's 503 class needs its sketch authority in the same
landing, the way `Timeout` is mirrored by `TestTimeoutStatusMirrorsSketch`.

The r0 prototype implements M1 and M2 as they stood in r0, and the M3 CLI verb. It does not
implement the M3 sketch edit, AC-9, or any r1 change (see "What the prototype does not do").

## (d) Acceptance criteria + non-vacuity

Every AC names the mutation that turns it red. For refusal and limit branches there is one
mutation per branch.

| AC | statement | test | mutation → observed red |
|---|---|---|---|
| AC-1 | non-unique names: 3 objects sharing one semanticId are all returned; an object with a different id is not | `TestObjectsBySemanticIDNonUniqueOrderedAndPaged` | MU3 (prefix match) → store test red |
| AC-2 | ordering determinism: items are in strictly ascending `hash` order, identical across two calls and independent of insertion order | same test (inserts in reverse-hash order) | MU1 (DESC) → store + route tests red |
| AC-3 | pagination: `after` resumes strictly after the cursor; the union of pages = the full set, no dupes | same test | MU2 (`>=`) → store + route red; MU4 (LIMIT ignored) → store + route red |
| AC-4 | empty vs bad request: unknown name → 200 exactly `{"items":[]}`; empty name → 400 | `TestObjectsBySemanticIDRoute` | MU8 (empty-name branch) → route test red |
| AC-5 | cap + paging rule: with 520 objects, the default page is 100; `limit=100000` returns exactly 500; `?after=<hash of item 500>&limit=500` returns the last 20 (< limit, so the client stops); the union has 520 distinct, ascending hashes; **no response carries a `next` key** (r1) | `TestObjectsBySemanticIDRouteCapAndDefault` (r1 revision: assert via the last item's hash, not `next`) | MU11 (default 500) / MU12 (cap bypass) → cap test red; **MU13r** (r1: re-add a `next` field) → the "no `next` key" assertion red. Designed, not run |
| AC-6 | refusals: malformed `after` → 400; non-integer `limit` → 400; `/v1/objects/by-semantic-id` without a name does not reach the lookup | `TestObjectsBySemanticIDRoute` | MU9 (after unparsed) / MU10 (limit error ignored) → route test red, one each |
| AC-7 | a writable open (fresh or existing-without-index) provisions the index, and the plan uses it | `TestObjectsBySemanticIDUsesIndex`, `TestSemanticIDIndexBuiltOnWritableReopen` | MU6 (provisioning statement deleted) / MU7 (single-column index) → index tests red. In r0 the DDL was in `schema.sql` and MU6/MU7 mutated it there; r1 mutates `provisionLookupIndex` |
| AC-8 | **store-level result equivalence only** (r1): the *unguarded* lookup SQL returns byte-identical rows with and without the index. This is not a serving fallback | `TestObjectsBySemanticIDUsesIndex` (r1: runs `objectsBySemanticIDSQL` directly after DROP INDEX, since the guarded method would now refuse) | recorded as a property, not killed by a mutation |
| AC-9 | (M3) the mux's `/v1` patterns equal `routes()` in the sketch, **spelled identically** (both `{name...}` rows: the new route and the corrected `/v1/registry/{name...}`, V28), count 9 | NOT PROTOTYPED | M3 executor must show: delete the mux line → red; delete a sketch row → red; revert `worlddapi.ail:78` to `{name}` → red |
| AC-12 | the new route is in the shared read-route harness, so the existing deadline (503) and sanitization (500) gates cover it | `seedReadRoutes` + `blockingStore`/`recordingStore`/`failingStore` overrides in `read_deadline_test.go` | MU14, MU15, MU16 |
| AC-10 | store error → 500 sanitized; deadline → 503 | `TestObjectsBySemanticIDRouteStoreErrors` (fault-injecting `readStore` wrapper, the V13 seam's purpose) | MU14 / MU16 → deadline / sanitize tests red |
| AC-11 | walk re-timed at N=10,000: route locate p50 ≤ 1 s and far below the scan | `TestWalkRetimedAtScale` (env-gated, `WORLD_WALK_N`) | reported as measurement, not a gate (see §Measurements) |
| AC-13 | **(r1, astra's named test)** an unindexed v3 store opened **read-only** answers `*LookupIndexUnavailableError`, and **no lookup query executes** (a test-only hook, `objectsBySemanticIDBeforeQuery`, is called 0 times; precedent `readObjectBetweenStatements`, V26). Then a **writable** open provisions the index, and a fresh read-only open's lookup succeeds and `EXPLAIN QUERY PLAN` shows `USING INDEX objects_by_semantic_id` | `TestLookupIndexGuardReadOnlyThenProvision` | **MU20** (guard removed: verdict ignored) → error expected, got rows, hook count 1 → red. **MU21** (`OpenReadOnly` skips `verifyLookupIndex`, zero-value verdict "present") → same red |
| AC-14 | (r1, **extended r2**) incompatible same-named definitions are refused, one fixture each: (a) `objects(semantic_id)` only; (b) **partial** `objects(semantic_id, hash_ref) WHERE semantic_id <> ''`; (c) **non-BINARY collation** `objects(semantic_id COLLATE NOCASE, hash_ref)`; (d) `objects(semantic_id, hash_ref DESC)`. **For each fixture:** writable `Open` returns an error and the index's `sqlite_master.sql` is byte-unchanged afterwards (not modified); a read-only open's lookup returns `*LookupIndexUnavailableError`; the lookup hook `objectsBySemanticIDBeforeQuery` is called **0** times. The successful provision-and-query-plan test for the accepted definition (AC-7 / AC-13's second half) is **retained** | `TestLookupIndexIncompatibleIsRefused` (table-driven over (a)–(d)) | **MU22** (column-name check removed) → (a) red. **MU26** (`index_list.partial` check removed) → (b) red: writable `Open` succeeds, read-only lookup runs, hook = 1. **MU27** (`coll = BINARY` check removed) → (c) red, same symptoms. **MU28** (`desc = 0` check removed) → (d) red, same symptoms. **MU29** (writable `Open` "repairs" by `DROP INDEX` + recreate) → the byte-unchanged `sqlite_master.sql` assertion reds on (a)–(d). All **designed, not run; executor runs them and records the observed failures** |
| AC-15 | (r1) provisioning failure is surfaced: with `lookupIndexProvisionDeadline` forced to an already-expired value (test seam), writable `Open` returns an error wrapping `context.DeadlineExceeded` and no `*Store` | `TestLookupIndexProvisionFailureIsSurfaced` | **MU23** (provision error discarded: `_ = provisionLookupIndex(…)`) → `Open` succeeds → red |
| AC-16 | (r1) route mapping: a `readStore` returning `*LookupIndexUnavailableError` → **503**, class `LookupIndexUnavailable`, the fixed remediation text, no host detail; the class/status equal the sketch's new `httpStatus` arm (parsed by `sketchHTTPStatusVectors`) | `TestLookupIndexUnavailableIs503` (fault-injecting `readStore` wrapper; production cannot reach this branch through `New`, V23) | **MU24** (mapping deleted → falls to `writeInternalError`) → 500 `Internal` → red; **MU25** (sketch arm `=> 503` edited to `=> 500`) → mirror assertion red |

## Prototype manifest (r0; files as measured in this worktree, now moved)

From `git status --short` / `git diff --stat` against `c30b967` (V20). Production lines are
`+` lines, comments included. **After r0 the controller moved these files out of the worktree to
`~/.ailang/state/world-iter201/prototype/` (same relative paths). The revision round edits the doc
only.** The table describes the **r0** prototype. It still has the index DDL in `schema.sql`, the
`next` field, and no index guard. Those are exactly the r1 changes, designed in B4/B5/B7 and
estimated in §c, and not prototyped.

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
| `design_docs/implemented/w-object-lookup-by-semantic-id.md` | new | this doc | doc |

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
| MU13 | `next` set on any non-empty page | route Route, CapAndDefault | `json: cannot unmarshal string into … []map` (a `next` key appeared on a non-full page). **Retired in r1**: the `next` field is dropped, and MU13r (§d AC-5) replaces it |
| MU14 | timeout branch disabled (`if false && timedOut`) | TestDaemonReadDeadline/real-store-expired-deadline, /blocking-store | `status = 500, want 503` |
| MU15 | lookup runs under `r.Context()`, not `readCtx` | TestDaemonReadDeadline (both arms) | `status = 200, want 503` |
| MU16 | 500 writes `err.Error()` instead of `writeInternalError` | TestInternalErrorsAreSanitized/read-routes | `the 500 body leaked the store's internal detail "kQ7v-store-detail-9f3c1d82"` |
| MU17 | items carry `Payload: &[]byte{}` | route Route | `item carries a payload key` |
| MU18 | mux line removed | route Route, CapAndDefault; TestDaemonReadDeadline/normal-deadline-answers-200 … | `match status = 404 body=404 page not found` |
| MU19 | **green control**: comment-only edit | — (both packages `ok`) | none, as required |

**r0 battery: 18 killed, 0 survived, 1 green control stayed green.** **r1 mutations MU13r and
MU20–MU25 (§d AC-5, AC-13..AC-16) and r2 mutations MU26–MU29 (§d AC-14) are DESIGNED, NOT RUN.** This round edits the doc only, and the
M1/M2 executor must run them and record the observed reds. MU3 survives the *daemon* package alone
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
| V23 | **`OpenReadOnly` has 0 production callers; every production open is writable** (census for B4's guard) | `/usr/bin/grep -rn "OpenReadOnly(\|store\.Open(" host cmd \| grep -v _test` (worktree) and `git grep -n "OpenReadOnly(\|store\.Open(" c30b967 -- host cmd \| grep -v _test` | `OpenReadOnly(` matches **only its definition** `host/store/store.go:273`. **Positive control, same command:** `store.Open(` matches **7** production call sites: `host/daemon/daemon.go:470`, `cmd/ailang-worldd/session.go:145`, `:211`, `cmd/world-publish/main.go:261`, `:295`, `:458`, `cmd/world-publish/transitions.go:69` |
| V24 | `GET /v1/log`'s list shape is an `{"items":[…]}` envelope with no `next` (refutes objection 2's raw-array premise) | `git show c30b967:host/daemon/handlers.go \| grep -n "type logRangeResponse\|Items \[\]logEntryResponse\|logRangeResponse{Items: items}"` | `63: type logRangeResponse struct {`, `64: Items []logEntryResponse` with struct tag `json:"items"`, `492: writeJSON(w, http.StatusOK, logRangeResponse{Items: items})` |
| V25 | the sketch's `ApiError` has `Timeout` → 503 (the arm pattern the new variant follows) | `git show c30b967:design_docs/sketches/worlddapi.ail \| grep -n "\| Timeout(string)\|Timeout(_) => 503"` | `96: \| Timeout(string)`, `113: Timeout(_) => 503` |
| V26 | a test-only hook precedent exists in the store | `git show c30b967:host/store/read_object.go \| grep -n "var readObjectBetweenStatements\|test-only scheduling input"` | `35: … is a test-only scheduling input`, `38: var readObjectBetweenStatements func()` |
| V27 | **(r2)** B3's precedents: `GetObject` parses with `hashref.Parse`; `ScanUnreadableLog` enforces its bound with `InvalidLimitError`; no production code drops an index | `git show c30b967:host/store/store.go \| awk 'NR>=475 && NR<=506' \| grep -n "func (s \*Store) GetObject\|hashref.Parse"`; `git show c30b967:host/store/scan.go \| grep -n "func (s \*Store) ScanUnreadableLog\|invalidScanLimit(op, limit)\|return &InvalidLimitError"`; `git grep -n "DROP INDEX\|CREATE TABLE IF NOT EXISTS objects" c30b967 -- host cmd ':!*_test.go'` | `GetObject` opens at `store.go:475` and calls `hashref.Parse(ifaceText)` at `store.go:493`. `scan.go:53 func (s *Store) ScanUnreadableLog`, which calls `invalidScanLimit(op, limit)` at `:55`, which returns `&InvalidLimitError{…}` at `:46`. So the bound is enforced through a helper, and `ScanUnreadableWorlds` shares it at `:95`. `DROP INDEX`: **no match**. **Positive control, same command:** `schema.sql:13 CREATE TABLE IF NOT EXISTS objects (` |
| V28 | **(r2)** the sketch spells the registry route `{name}`, and the mux spells it `{name...}` | `git show c30b967:design_docs/sketches/worlddapi.ail \| grep -n 'registry/{name'`; `git show c30b967:host/daemon/daemon.go \| grep -n 'GET /v1/registry/{name...}'` | `78: { method: "GET", path: "/v1/registry/{name}" },`; `643: mux.HandleFunc("GET /v1/registry/{name...}", d.handleRegistry)` |
| V29 | **(r2)** `PRAGMA index_list` reports `partial`, and `PRAGMA index_xinfo` reports `desc`, `coll`, `key` | scratch probe: `sqlite3 pragma-probe.db "CREATE TABLE objects (hash_ref TEXT PRIMARY KEY, semantic_id TEXT NOT NULL); CREATE INDEX good ON objects(semantic_id, hash_ref); CREATE INDEX part ON objects(semantic_id, hash_ref) WHERE semantic_id <> ''; CREATE INDEX nocase ON objects(semantic_id COLLATE NOCASE, hash_ref); CREATE INDEX descd ON objects(semantic_id, hash_ref DESC);"`, then `sqlite3 -header … "PRAGMA index_list('objects');"` and `"PRAGMA index_xinfo('<i>');"` for each (`sqlite3` 3.51.0; the store's driver is `modernc.org/sqlite v1.54.0`, go.mod:7. **The executor must re-assert these columns through the driver in AC-14**, since the CLI probe is not the production SQLite) | `index_list`: header `seq\|name\|unique\|origin\|partial`; `part … partial=1`, `good`/`nocase`/`descd` `partial=0`. `index_xinfo(good)`: `0\|1\|semantic_id\|0\|BINARY\|1`, `1\|0\|hash_ref\|0\|BINARY\|1`, `2\|-1\|\|0\|BINARY\|0` (header `seqno\|cid\|name\|desc\|coll\|key`). `index_xinfo(part)`: **identical to `good`**. `index_xinfo(nocase)`: `semantic_id … NOCASE`. `index_xinfo(descd)`: `hash_ref` `desc=1` |
| V22 | gates on the prototype | `go vet ./...`; `go test -count=1 -run '^$' ./...`; `go test -race -count=1 ./host/store/ ./host/daemon/ ./cmd/ailang-worldd/`; `./scripts/verify_ail.sh`; `go test -count=1 ./...` (all with `AILANG_BIN` pinned, `GOCACHE`/`TMPDIR` outside the tree) | vet clean; compile fence clean; gofmt clean on touched files; `-race`: store `ok` 15.8s, daemon `ok` 17.9s, cmd/ailang-worldd `ok` 3.0s; `verify_ail.sh`: "verify gate PASSED: 16 required identities verified, 40 named tests pass". Full suite: **23 ok / 1 FAIL**. The FAIL is `host/broker` `TestOrdinaryOverflowCarriesNoESRCH` ("overflow kill: operation not permitted"), an EPERM on the kill path in a package this prototype does not touch (`git diff --stat c30b967 -- host/broker` is empty). **Re-run alone: `ok host/broker 66.481s`.** It is the same EPERM-under-load family as row 115 (capsule). I did not find a broker-specific row, so it is recorded here and not claimed as a known flake |

## (f) Conflict Surface

| file | this row's change | open rows touching the same file |
|---|---|---|
| `host/store/schema.sql` | r1: **comment only** (the DDL moved to `provisionLookupIndex`) | **103** (adds indexes + reverse read), **104** (commit-membership schema change). No shared name. Both will need a version decision this doc does not take (B4). Row 103 may reuse the `provision…`/`verify…` pattern for its own indexes. |
| `host/store/store.go` | r1: `Open`/`OpenReadOnly` call provision/verify; one `Store` field | **23 / PR #153** (store deadline policy; `lookupIndexProvisionDeadline` is a new named bound its tranche should list), **111** (busy/lock status). Textual merge only. |
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

## Quorum log

**Round 1 — BLOCKED 2/2 present.** Seats: `gpt6-astra` **reject**, `gemini-3-1-pro` **reject**;
`oc-glm-5-3` and `oc-kimi-k3` **absent (unreachable)**. Revision by the same designer.

| seat | objection (summary; verbatim in the controller's round-1 record) | measurement | disposition | what changed |
|---|---|---|---|---|
| gpt6-astra | B4 let a read-only store without the index silently execute the lookup as a table scan, which violates no-silent-fallback. AC-8 shows result equivalence, not compliance with the operational contract | V23: `OpenReadOnly` has **0** production callers (positive control: 7 `store.Open(` sites), so production could not reach the fallback. **The principle stands anyway** | **Applied verbatim** | B4: the fallback paragraph is replaced by astra's text; the check lives at open and is cached (cost argued); provisioning moved out of `schema.sql` into a 30 s-bounded `provisionLookupIndex`, with failures making `Open` return an error; B3 guard before the query; B6 503 `LookupIndexUnavailable` row; sketch `ApiError` variant (B1/B8); AC-13 (astra's named test) + AC-14..AC-16; MU20–MU25; AC-8 narrowed to store-level equivalence |
| gemini-3-1-pro | the `{"items","next"}` envelope contradicts B1's "no new function, contract or type" and the list precedent; `GET /v1/log` returns a raw array; `next` is redundant with the last item's hash | V24: **premise false.** `GET /v1/log` returns `logRangeResponse{Items}` = `{"items":[…]}` (`handlers.go:63-65`, `:492`), not a raw array. **Two parts true:** B1 overclaimed, and `next` is not in the precedent and is redundant | **Partly applied** | envelope **kept** as `{"items":[…]}` to match `GET /v1/log` exactly; **`next` dropped**; the paging rule adopted verbatim (B5); B1 now lists exactly what is added; AC-4/AC-5 rewritten; MU13 retired → MU13r |

Not re-litigated: the ninth-frozen-route decision (B1), hash ordering (B5), empty-200 (B6), and no
version bump (B4). Neither seat objected to them.

**Round 2 — BLOCKED 2/2 present.** Seats: `gpt6-astra` **reject**, `gemini-3-1-pro` **reject**;
`oc-glm-5-3` and `oc-kimi-k3` **absent**. Neither objection disputes the design direction, and
both carry concrete reviewer-authored fixes. **Carve-out taken (narrow refinement, reviewer fixes
verbatim); no r3 quorum.** The doc goes to the planner.

| seat | objection (summary) | measurement | disposition | what changed |
|---|---|---|---|---|
| gpt6-astra | r1's compatibility check (columns + table only) accepts a partial index, a non-BINARY collation or a DESC key: an index that exists but does not serve the query's order or coverage | V29 probe: a partial index's `index_xinfo` is byte-identical to the accepted one, so only `index_list.partial` detects it; `coll`/`desc` are reported per key | **Applied verbatim** | B4 compatibility rule replaced with astra's text + the PRAGMA mechanics; AC-14 extended with (b) partial, (c) NOCASE and (d) DESC fixtures, each asserting writable-open refusal without modification, read-only `LookupIndexUnavailable` and 0 hook calls; accepted-definition provision + plan test retained; MU26–MU29 designed, not run; M1 estimate ~102 → ~114 |
| gemini-3-1-pro | B3 cited V12 for precedents V12 does not verify (`hashref.Parse` in `GetObject`, `InvalidLimitError` in `ScanUnreadableLog`); the B1 sketch path `{name}` would not match the mux's `{name...}` in AC-9's parity test | V27: both precedents hold (`store.go:493`; `scan.go:53→55→46`, via the `invalidScanLimit` helper). V28 (controller measurement, re-run): the **existing** registry row has the same `{name}` vs `{name...}` drift (`worlddapi.ail:78` vs `daemon.go:643`) | **Applied verbatim** | V27 added and B3 cites it instead of V12; B1's sketch path is now `{name...}`. **Sibling-row change:** M3 also changes `worlddapi.ail:78` to `/v1/registry/{name...}`. This is the same fix, applied under the reviewer's stated rationale (AC-9 compares mux patterns with `routes()` literally); AC-9 gains a matching mutation |

## (h) Decision for Mark

None. Adding a frozen route is a technical extension governed by the existing
doc-plus-quorum rule, and the quorum on this doc is that ratification. Every other choice here
(hash order, empty-200, no version bump) is technical and argued above.

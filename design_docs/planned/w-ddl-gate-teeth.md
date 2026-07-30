# w-ddl-gate-teeth — make schema drift observable

**Status:** **PARKED — `needs-human-review`.** The pick-time quorum ran its full two rounds and
**BLOCKED both times**. One blocking objection remains, and it **disputes the design DIRECTION**, so
the charter's narrow-refinement carve-out does not apply and the controller may not resolve it.
**This document is a decision packet for the human, not authorization to implement DG.A.** No
milestone may be routed to a planner or an executor until **OD-5** is answered. See the Quorum
verification log at the end.  
**Item:** `w-ddl-gate-teeth` (queue item 4d)  
**Clause:** clause-1  
**Date:** 2026-07-30  
**Iteration:** 41  
**Author:** `codex:gpt-5.6-sol (rotation designer)`; controller `claude-opus-5` supplied every
measured premise and added the Quorum verification log, the label key and OD-5  
**Estimate:** 0.25–0.5 day for DG.A, once unblocked

---

## Problem statement

The repository says that
`TestPreJournalMigrationPreservesExistingDDL` guards SQLite DDL, but the assertion credited
with doing that cannot distinguish any current production change.

The test takes its “old” schema from the same, already-edited `schema.sql` that `store.Open`
will execute. A source hash stops edits before a database exists. If a developer updates that
hash, the test creates both sides of its comparison from the new source. Existing-table edits
therefore compare equal even though SQLite's `CREATE TABLE IF NOT EXISTS` has not applied them.
The comparison can presently notice only a newly added non-journal table. It explicitly deletes
`journal`, so no repository test owns the journal table's exact DDL.

This is not a request to build a migration platform. The small repair is to make two facts
observable in tests:

1. a fresh store materializes the exact canonical DDL for every durable table, including
   `journal`; and
2. opening a store built from an independent historical schema produces the current expected
   DDL for every table that already existed.

The second fact is deliberately stronger than a source-change detector. If a developer makes a
real edit to an existing table and updates the canonical expectation, the fresh-store check
becomes green but the historical-store check stays red until an upgrade mechanism exists or an
explicit compatibility decision changes the expectation.

Detection and production policy are separate. Whether `store.Open` should use
`PRAGMA user_version`, fail loudly, or run migrations changes the durability kernel and is
ratification-class. This document does not settle or implement that behavior.

## Premise Verification Log

**Label key** (added by the controller; quorum round 2 correctly flagged `M4`/`M5` as undefined
references). `M1`–`M5` name the controller's five first-party measurements taken at `e5027df`
during iteration 41's Gate 2, before this document existed:
**M1** `MUT-JOURNAL-DDL-WIDEN` — widening journal's `kind` CHECK leaves all 10 packages green.
**M2** `MUT-DDL-DRIFT` — a real pre-journal DDL edit reds **only** the sha256 source pin, which
fires before any database is opened. **M3** `MUT-DDL-COMPARE-DEAD` — short-circuiting the
`sqlite_master` comparison leaves the package green, so it discriminates nothing today.
**M4** `MUT-DDL-DRIFT-REPINNED` — keeping M2's edit and updating the pinned literal, which is what
a developer legitimately does next, re-greens all 10 packages in one line. **M5** the
production-path probe — an existing store opened by a binary carrying the new schema gets
`store.Open` returning `nil` with the DDL **silently unapplied**.

All measured rows below are **VERIFIED BY THE CONTROLLER at `e5027df`, iteration 41**.
Every gate run used darwin/arm64, `GOTOOLCHAIN=go1.25.6`, a throwaway worktree whose
`go.mod` floor was temporarily lowered from 1.26.4 to 1.25.6,
`AILANG_BIN=/tmp/ailang-v0300/ailang`, z3 4.16.0, and sqlite3 CLI 3.51.0.
Each mutation took its backup in the same command that applied it, restored with `cp`, printed
before/after SHA-256, and ended with empty `git status`.

| Claim | How verified | Evidence | Verdict |
|---|---|---|---|
| Baseline is green on the measured rig | Full first-party suite | `go test ./... -count=1` returned rc=0; 10/10 packages green; `host/store` 2.832s | CONFIRMED |
| The cited gate hashes the pre-journal source before opening a database | First-party read of `host/store/journal_test.go:631-662` | Prefix through schema line 77; literal SHA-256 check at lines 634-637; database open begins only afterward | CONFIRMED |
| The cited comparison excludes journal | First-party code read | `delete(after, "journal")` precedes `reflect.DeepEqual(after, before)` | CONFIRMED |
| The journal DDL can change without any test failing | `MUT-JOURNAL-DDL-WIDEN` in `host/store/schema.sql` (**PRODUCTION**) added `MUTANT` to the `kind` CHECK | Full suite rc=0, 10/10 packages green | CONFIRMED |
| The journal mutation materially changes SQLite DDL | Mutated and backup schemas independently executed by sqlite3 CLI; `sqlite_master` read | Mutated materialized SQL contained `MUTANT`; backup materialized SQL did not | CONFIRMED; non-vacuity control |
| The named DDL-drift mutation reds by the source pin, not the database comparison | `MUT-DDL-DRIFT` in `host/store/schema.sql` (**PRODUCTION**) added `mutant_col` to `store_heads` | Exactly one test red at `journal_test.go:636`, reporting the new SHA-256; no database had yet opened | CONFIRMED |
| The `sqlite_master` comparison contributes no discrimination today | `MUT-DDL-COMPARE-DEAD` in `host/store/journal_test.go` (**TEST**) short-circuited the comparison while keeping both variables used | `go vet ./host/store/` rc=0; package and full suite rc=0, 10/10 packages green | CONFIRMED; test-side discrimination probe only |
| Re-pinning after a legitimate DDL edit re-greens everything | `MUT-DDL-DRIFT-REPINNED` edited `host/store/schema.sql` (**PRODUCTION**) and its SHA literal in `host/store/journal_test.go` (**TEST**) | Package and full suite rc=0, 10/10 packages green | CONFIRMED |
| `store.Open` silently leaves an existing table at old DDL | Old schema created by sqlite3; a Go program using mutated production `schema.sql` called real `store.Open` | `store.Open` returned nil; `store_heads` lacked `mutant_col` before and after | CONFIRMED |
| The comparison is self-referential for modified existing tables | First-party code read plus the preceding production-path experiment | `before` is created from the current pre-journal prefix; `store.Open` executes that same current source | CONFIRMED |
| Schema text runs on every open | First-party read of `host/store/store.go:262` | `db.Exec(schemaSQL)` executes the embedded schema verbatim | CONFIRMED |
| No host migration machinery exists | First-party repository inspection | No migration implementation under `host/`; no `user_version` use | CONFIRMED |
| The historical fixture has exact released provenance | First-party history plus SHA-256 | `d5774eb` added journal; `8133573:host/store/schema.sql` is the last pre-journal schema, has SHA-256 `35f09862e20ddc1c6b0467b69781b2d25fbc07d04c49f777a76a62793e14bbdd`, and contains exactly six `CREATE TABLE` statements | CONFIRMED |
| The pre-journal artifact materializes exactly six durable tables | sqlite3 CLI 3.51.0 executed `8133573:host/store/schema.sql` and queried `sqlite_master` | Exactly `objects`, `worlds`, `log_entries`, `epoch_registry_heads`, `store_heads`, and `verification_cache` | CONFIRMED |
| The current schema materializes exactly seven durable tables and no declared secondary objects | sqlite3 CLI 3.51.0 executed HEAD `schema.sql`; source and `sqlite_master` inspected | Exactly `epoch_registry_heads`, `journal`, `log_entries`, `objects`, `store_heads`, `verification_cache`, and `worlds`; source has seven `CREATE TABLE` and zero `CREATE INDEX`, `TRIGGER`, or `VIEW` statements | CONFIRMED |
| `tableDDL` returns only durable tables today, but its two filters do different work | First-party helper read plus materialized-store query | `type='table'` drops all seven `sqlite_autoindex_*` index rows; `name NOT LIKE 'sqlite_%'` drops zero rows today. The name filter remains a dormant guard for a future `sqlite_sequence` table from `AUTOINCREMENT`; current bare `INTEGER PRIMARY KEY` creates none | CONFIRMED; live and dormant limbs distinguished |
| The AC2 sentinel survives current `store.Open` by construction | First-party read of `host/store/store.go:250-268` | Open enables foreign keys and executes schema, never reads or validates `store_heads`; the schema declares no foreign keys | CONFIRMED |

| **The 6 tables in the `8133573` artifact match HEAD's DDL exactly after whitespace normalization.** Raised by quorum r2 as an unverified, load-bearing premise: if ANY drift existed, D2 would RED on baseline and this design would be broken as written | sqlite3 CLI 3.51.0 materialized `8133573:host/store/schema.sql` and HEAD's `schema.sql` into two independent databases; the `sqlite_master` text of the 6 shared tables was compared after collapsing whitespace runs | All **6 identical** — `objects`, `worlds`, `log_entries`, `epoch_registry_heads`, `store_heads`, `verification_cache`. **Zero drift.** HEAD materializes one further table (`journal`) and no others | **CONFIRMED — the premise HOLDS, so D2 is green on baseline** |
| gemini r2's contingency does not fire | Same measurement | The reviewer's fallback — *"if they do NOT match, the test must compare historical tables against a separate baseline manifest, not the current one"* — is **not needed**; one canonical manifest is correct | CONFIRMED |

No local gate result is added by this designer. Per the iteration directive, any such result
inside this sandbox would be **UNINFORMATIVE UNDER SANDBOX**.

## Conflict Surface

| Existing surface or precedent | Collision / reuse question | Resolution |
|---|---|---|
| `host/store/schema.sql` and `store.Open` | Both are the durability kernel; changing either would require human ratification | Milestone DG.A changes neither. It adds test-side observations only |
| `TestPreJournalMigrationPreservesExistingDDL` | Its historical intent is useful, but its source slice, SHA pin, and same-source comparison are the defect | Replace its mechanism in place; do not add a second overlapping “migration” test |
| `tableDDL` helper | Already reads exact table SQL from `sqlite_master`; which predicate excludes non-durable rows? | Reuse it. Today `type='table'` excludes seven autoindexes and the `sqlite_%` name filter is dormant; retain the latter for a future `sqlite_sequence` table |
| `scripts/verify_go.sh` | Already owns repository Go build/test execution and the loud `AILANG_BIN` guard | No new script leg. New tests run under the existing bounded CI job and ordinary package test command |
| `scripts/verify_ail.sh` required-check manifest | Hardcoded named identities are an anti-vacuity precedent | Follow the hardcoded-manifest principle, but in the Go test that owns SQLite DDL |
| M2.A benchmark-name manifest gate | A dropped benchmark otherwise leaves `go test` green | Follow the same exact-name-set precedent: missing and unexpected durable table names both fail |
| Existing runtime pre-journal SHA-256 assertion | It is a one-line change detector, not a correctness gate; M4 proves the sanctioned re-pin ritual is green | Remove the runtime assertion; retain the historical artifact's SHA only as fixture provenance |
| Existing journal behavior tests | They cover application behavior and uniqueness, but not the complete materialized DDL | Keep them. Exact DDL is complementary and intentionally stricter |
| SQLite formatting in `sqlite_master.sql` | Comparing raw source text would be brittle to comments and statement layout | Compare SQLite's materialized SQL, normalized only for whitespace as specified below |
| `design_docs/sketches/` and `.ail` CI | This item is host-boundary Go and contains no AILANG design | No `.ail` sketch; `scripts/verify_ail.sh` and its required-check manifest remain untouched |

### Why this is not a package

This change observes the Go/SQLite durability boundary and edits only its Go tests. An AILANG
package cannot inspect `sqlite_master`, create an old SQLite store, or constrain `store.Open`.
No new kernel API, world type, policy, or reusable package behavior is introduced.

## Design

### D1. One canonical materialized-DDL manifest

Add a test-local `map[string]string` in `host/store/journal_test.go` containing the expected
materialized `CREATE TABLE` SQL for all seven durable tables:

- `objects`
- `worlds`
- `log_entries`
- `epoch_registry_heads`
- `store_heads`
- `verification_cache`
- `journal`

The values are not copied source statements with comments. They are the `sqlite_master.sql`
forms SQLite materializes from the current schema. A small test-only normalizer collapses runs
of Unicode whitespace to one ASCII space and trims the ends. It does not lowercase, reorder
columns or constraints, remove quotes, or parse SQL. Column order, constraint text, and table
names therefore remain significant.

The manifest is hardcoded by design. It is review-visible policy, like the named AILANG checks
and benchmark-name manifest. It must not be derived from `schemaSQL`, a hash of `schemaSQL`, or
the database being checked. Deriving expected values from the subject would recreate the defect.

`TestSchemaDDLMatchesCanonicalManifest` will:

1. create a fresh temporary-file store with `Open`;
2. collect DDL with the existing `tableDDL`;
3. normalize actual and expected values;
4. sort actual and expected table names, then compare the exact name sets, failing loudly on
   zero tables, a missing table, or an unexpected table; and
5. walk the same sorted names and compare every table's full normalized DDL with a deterministic
   diff whose failure message names the offending table.

There is no loop without a ceiling: the loop is bounded by the seven-entry manifest plus the
finite result set returned by `sqlite_master`. The database is a temporary file and every handle
is closed by test cleanup.

This gives journal teeth. `MUT-JOURNAL-DDL-WIDEN` changes the production schema SQLite
materializes while the independent canonical entry remains unchanged, so the journal row fails.
The explicit seven-name assertion prevents “fixing” the test by filtering journal out.

This gate is intentionally a reviewed change detector. A developer can update the manifest
alongside an intended journal change. That is necessary—intentional DDL must be expressible—and
means this test alone does not prove compatibility, data preservation, or semantic correctness.

### D2. An independent historical-store fixture

Replace the current source-prefix construction with a test-local constant named
`preJournalSchemaV0`, copied verbatim from the released artifact
`8133573:host/store/schema.sql`, not hand-authored or derived from current `schemaSQL`. Its
fixture comment cites commit `8133573`, notes that `d5774eb` added journal, records SHA-256
`35f09862e20ddc1c6b0467b69781b2d25fbc07d04c49f777a76a62793e14bbdd`, and forbids
mechanically updating the fixture when current DDL changes.

`TestOpenAddsJournalAndDetectsStalePreJournalDDL` will:

1. execute `preJournalSchemaV0` into a temporary-file database;
2. insert one valid sentinel row into `store_heads`;
3. close that setup connection;
4. call the real production `Open` on the same path;
5. assert the sentinel row still exists unchanged;
6. assert `journal` exists; and
7. compare each of the six historical table names after open against that table's current
   canonical manifest entry.

The test name and an adjacent comment must say what the test proves: current `Open` creates the
absent journal table but does not upgrade existing tables. A mismatch means the historical
table is stale relative to the canonical manifest; a green result must not be described as an
upgrade performed by `Open`.

The comparison intentionally ranges over the six historical names rather than deleting journal
from a full before/after map. Journal creation has its direct assertion and D1 owns journal's
exact DDL. Keeping the scopes separate makes the named mutations discriminating:
`MUT-JOURNAL-DDL-WIDEN` reds D1, while a repinned edit to an existing table reds D2.

The test must also sort the six historical names, assert that the set is exactly the six
materialized names verified above, require each name in the canonical manifest, and perform
per-table comparisons in that order. Every mismatch is deterministic and names the offending
table. Those assertions make accidental fixture erosion red. The artifact SHA appears as
provenance in the fixture comment, not as a runtime pin: a runtime hash would only prove that a
developer did not update the constant.

The sentinel is a narrow data-preservation control, not a migration suite. One row is sufficient
because current `Open` only adds the absent journal table. Future migrations that rewrite data
must add migration-specific fixtures and assertions rather than treating this sentinel as proof
for every data shape.

### D3. Retire the dead comparison; retain its useful claim honestly

Remove:

- the `schemaSQL[:strings.Index(...)]` old-schema derivation;
- the SHA-256 pin and its imports;
- the before/after comparison built from the same source;
- `delete(after, "journal")`; and
- the old test name if it implies that unchanged DDL is itself a migration.

Reuse `tableDDL`, with any normalization kept in a separate test helper. The repaired claim,
placed beside the renamed test as well as here, is not “Open upgrades a pre-journal store.” It is:

> Opening the named historical store creates journal without losing the sentinel, and the test
> detects if any historical table remains stale relative to the current canonical DDL.

That claim can fail on a realistic production edit. If `store_heads` gains `mutant_col` and the
canonical manifest is legitimately updated, SQLite leaves the historical table unchanged and
D2 reports the exact mismatch. This directly closes M4 and M5's one-line re-green path.

### D4. Detection now; production response only after ratification

**Direction objection, stated fairly.** A test-only detector does not make runtime drift visible
to production callers: `store.Open` still returns success for a structurally stale store. The
reviewer therefore asks that DG.A be ratification-blocked and include a production version or
fingerprint check that returns distinct errors for legacy, stale, and future stores.

**Position: option (a), maintain the deferral.** The fail-open behavior is pre-existing at HEAD;
DG.A neither introduces nor widens it. Not yet repairing that kernel defect is not the same as
adding a silent fallback. DG.A strictly increases detection by making a proposed incompatible
DDL edit red before landing. The frozen-kernel ratification rule is itself the controlling
guardrail: a headless loop may add this test extension, but may not change which on-disk stores
`Open` accepts. Consequently DG.A is not conditional on ratification.

**Reconciliation with this document's status, so the two are not read as contradicting each other
(the paragraph above is this document's POSITION; the header records the QUORUM'S OUTCOME).** The
quorum did not accept that position: `gpt5-6-sol` rejected it in both rounds, so as a matter of
process DG.A is blocked on **OD-5** regardless of the argument's merits. What is being deferred to
the human is not whether the reasoning above is sound — it is who gets to decide between two
ratified guardrails when they point opposite ways. Until OD-5 is answered, **OD-5 governs and no
milestone here may be routed.**

Changing that requires a durable schema-version contract. `PRAGMA user_version` is the smallest
candidate, but setting or enforcing it changes what stores the kernel accepts. Automatically
migrating is larger still. OD-3 and OD-4 therefore remain human decisions.

The non-ratification milestone can land first because it changes no accepted store format and no
runtime behavior. Its red result becomes evidence required by the first future DDL-change design;
it is not permission for an executor to invent a migration.

## Acceptance Criteria

### AC1 — exact fresh-store DDL, including journal

Owned by **DG.A**. A fresh store's normalized `sqlite_master` map exactly equals a hardcoded
seven-table canonical manifest; zero, missing, and unexpected tables fail loudly.

Named RED mutation: **`MUT-JOURNAL-DDL-WIDEN`** edits
`host/store/schema.sql` (**PRODUCTION**) to add `MUTANT` to journal's `kind` CHECK.
`TestSchemaDDLMatchesCanonicalManifest` must red with a journal-specific mismatch. The mutation
is reachable as a realistic constraint change and materially alters SQLite DDL.

### AC2 — existing-store edits cannot be re-pinned green

Owned by **DG.A**. The independent pre-journal fixture, after real `store.Open`, matches current
canonical DDL for all six historical tables and preserves its sentinel row.

Named RED mutation: **`MUT-EXISTING-DDL-CHANGE-REMANIFESTED`** edits
`host/store/schema.sql` (**PRODUCTION**) to add
`mutant_col TEXT NOT NULL DEFAULT ''` to `store_heads`, and updates only the `store_heads`
entry in the test's canonical manifest in `host/store/journal_test.go` (**TEST**) as a developer
legitimately accepting that new DDL would. The fresh-store gate must remain green and
`TestOpenAddsJournalAndDetectsStalePreJournalDDL` alone must red because the old table was not
upgraded. This is the required discriminating M4/M5 control.

### AC3 — the historical fixture is independent and non-vacuous

Owned by **DG.A**. `preJournalSchemaV0` is a verbatim artifact-derived constant independent of
`schemaSQL`; the test asserts exactly six historical names, requires each in the canonical
manifest, and checks journal creation explicitly.

Named RED mutation: **`MUT-HISTORICAL-FIXTURE-DROP-STORE-HEADS`** edits
`host/store/journal_test.go` (**TEST**) to remove `store_heads` from the historical fixture.
The exact six-name fixture assertion must red before the DDL comparison. This is a test-side
drift check, not proof of a kernel property.

### AC4 — the dead mechanism is gone

Owned by **DG.A**. The repaired tests contain no pre-journal source slice, source SHA pin,
same-source before/after equality, or `delete(..., "journal")`; `tableDDL` remains the single
SQLite extraction helper.

Named RED mutation: **`MUT-CANONICAL-MANIFEST-DERIVED`** edits
`host/store/journal_test.go` (**TEST**) to replace journal's hardcoded expected DDL with the
actual journal DDL just read from the database. Re-running `MUT-JOURNAL-DDL-WIDEN` must then
leave AC1's required RED absent, so mutation review rejects the implementation. This is a
test-side discrimination probe: it proves independence of the gate, not journal semantics.

### AC5 — scoped verification and no unrelated surfaces

Owned by **DG.A**. The implementation changes only `host/store/journal_test.go`; formatting,
`go vet ./host/store/`, the focused two tests, `go test ./host/store/ -count=1`, and the normal
repository Go gate are run on a sound rig with their toolchain and `AILANG_BIN` recorded.
No `.ail` file, sketch, production Go file, schema file, gate script, or manifest changes.

Named RED mutation: **`MUT-UPGRADE-ASSERTION-DEAD`** edits
`host/store/journal_test.go` (**TEST**) to short-circuit only the six-table post-open DDL
comparison. Under `MUT-EXISTING-DDL-CHANGE-REMANIFESTED`, the expected D2 red disappears while
the fresh-store test stays green; mutation review rejects the implementation. This is a
test-side gate-discrimination probe.

## Milestones

### DG.A — replace the inert DDL gate (~0.25–0.5 day)

Modify only `host/store/journal_test.go`:

- add the seven-table canonical materialized-DDL manifest and whitespace normalizer;
- add the exact fresh-store manifest test;
- replace the self-derived pre-journal test with the independent `preJournalSchemaV0` fixture,
  sentinel, and six-table post-open comparison;
- remove the SHA/source-slice/dead-comparison mechanism and unused imports; and
- execute and record AC1's production mutation, AC2's combined realistic edit/re-manifest
  mutation, and AC3–AC5's test-side discrimination probes with byte-identical restoration.

DG.A owns **AC1**, **AC2**, **AC3**, **AC4**, and **AC5** individually. It owns
`host/store/journal_test.go`, which is capable of reddening every acceptance check. Mutation
commands must have a 120-second ceiling per focused/package invocation; the repository gate
retains its existing CI ceiling. No retry loop is allowed: one run per named mutation plus one
restored baseline run.

Anything that changes `host/store/schema.sql`, `store.Open`, accepted on-disk versions, or
migrates stored data is deferred behind OD-3/OD-4 and is not a second milestone hidden inside
DG.A.

**DG.A IS NOT ROUTABLE YET.** The pick-time quorum blocked this document in both rounds on a
directional objection, so DG.A is held on **OD-5** and must not be handed to a planner or an
executor until the human answers it. The milestone body above is what will be built *if* OD-5
resolves to alternative 1; it is a specification, not an authorization.

## Open Decisions for the human

> Numbering note (controller, iteration 41): these are **OD-3** and **OD-4**, not OD-1/OD-2.
> Iteration 40 already parked **OD-1** (lower the `go.mod` floor 1.26.4 → 1.25.6) and **OD-2**
> (file the go1.26 reproducer upstream at `golang/go`) for the same human, and the charter's
> parked-for-human list is a single global namespace — two live `OD-1`s meaning different
> things is the iter-31 ID-collision defect, where a collision reads as continuity.

### OD-5 — **THE CONTROLLING DECISION. Does the no-silent-fallback axiom oblige this item to change `store.Open` NOW, overriding the frozen-kernel deferral?**

*Added by the controller after quorum round 2. This is the decision that parks the item; OD-3 and
OD-4 are downstream of it. It is here because it is not a design question — it is a conflict
between two RATIFIED guardrails, and a headless loop must not pick a winner.*

**The conflict, stated with both sides at full strength.**

- **`gpt5-6-sol` rejected this doc twice, on the same ground, and sharpened it the second time:**
  *"The design knowingly preserves the production silent fallback it verified: `store.Open` returns
  success when an existing database's DDL is stale. A test-only detector does not satisfy the
  mission's no-silent-fallback axiom, and deferring the behavior to OD-3/OD-4 leaves the core defect
  in the proposed completed state."* And: *"That is axiom noncompliance, not merely an out-of-scope
  enhancement."* Its proposed fix is to mark DG.A **ratification-blocked** and add a ratified
  production milestone that classifies stores deterministically before accepting them — recognising
  fresh/current and the supported pre-journal shape, rejecting unknown-legacy, stale and
  future shapes with **distinct** errors, setting the version only after a successful bounded
  transaction, and shipping fixtures for all five cases.
- **The doc's answer (D4), also at full strength:** the fail-open is **pre-existing at HEAD** —
  the controller measured it as M5 — and DG.A neither introduces nor widens it. Not-yet-repairing an
  inherited kernel defect is not the same act as adding a silent fallback. And the charter's
  **frozen-kernel discipline** (*"kernel changes require explicit human ratification"*) is itself a
  ratified guardrail: changing which on-disk stores `Open` accepts is exactly a kernel behaviour
  change, so this loop **cannot** adopt the reviewer's fix, and Standing rule 2 forbids forcing
  past the objection in the other direction.

**Why the controller did not resolve it.** The charter's narrow-refinement carve-out permits a
bounded controller revision only when every remaining objection is non-directional and carries a
reviewer-authored fix. This objection disputes the direction outright, and resolving it would mean
either overriding frozen-kernel discipline headlessly or declaring a reviewer wrong on a judgment
call rather than on a measurement. Both are outside this loop's authority.

**Alternatives.**

1. **DG.A lands as designed (test-only), and OD-3 is approved in principle for the next item** —
   the cheap detector is not held hostage to the expensive kernel change, but the kernel change is
   committed to on a schedule rather than left open-ended.
2. **DG.A is ratification-blocked** — nothing lands until the `user_version`/fingerprint contract is
   designed, ratified and implemented together, as `gpt5-6-sol` asks. The repository keeps a gate
   with no teeth in the meantime.
3. **DG.A lands and the production response is dropped** — cheapest, and the option the reviewer's
   objection is squarely aimed at. Not recommended.

**Recommendation: alternative 1.** It is the only option that satisfies what the objection is
actually protecting — that a verified silent fallback does not get closed out as "done" — while
respecting the ratification gate. DG.A is strictly additive, changes no accepted store format and
no runtime behaviour, and its red is the evidence the `user_version` design will need; and pairing
it with an approved-in-principle OD-3 means the axiom is satisfied on a schedule instead of being
argued about. Note the honest cost either way: until OD-3 ships, **already-deployed binaries still
cannot diagnose schema drift**, which is the reviewer's point and it stands.

**Cost to defer this decision:** item 4d stays parked and the DDL gate stays inert — five measured
mutations show it cannot fail for any realistic production change, so every DDL edit until then
ships fail-open with no gate that would notice.

### OD-3 — add and enforce a `PRAGMA user_version` contract now?

**Question.** Should every new/current store receive a nonzero version and should
`store.Open` fail loudly when an existing store has an unsupported or un-upgraded version?
This is ratification-class because it changes durability-kernel acceptance behavior.

**Alternatives.**

1. **Land the pin and loud failure now.** Establish version 1, distinguish legacy version 0,
   and reject any version without an explicit supported path.
2. **Land a version marker but keep version 0 accepted.** Improves inventory but preserves the
   ambiguous fail-open path.
3. **Defer until the first item that needs an existing-table DDL change.** DG.A detects the
   missing upgrade path in review, but production continues to accept structurally stale stores.

**Recommendation:** alternative 1, in a separately ratified follow-up. Loud rejection is honest
and makes “not upgraded” distinguishable from “upgraded”; a marker that does not affect behavior
recreates the inert-gate shape. Before ratification, specify treatment of legacy version 0 and
prove fresh, supported, legacy, and future-version cases.

**Cost to defer:** no immediate implementation cost and DG.A still blocks future incompatible
DDL in the suite, but current production code retains M5's silent ambiguity. The first real DDL
change must stop for this decision and cannot claim an upgrade path merely because tests detect
its absence.

### OD-4 — when DDL changes, fail or migrate?

**Question.** Once a version mismatch is detected, should `store.Open` reject it, run bounded
in-process migrations, or require an explicit external migration command?

**Alternatives.**

1. **Fail loudly only.** Smallest kernel behavior; operators upgrade out of band.
2. **Run ordered in-process migrations.** Best opening experience, but adds transactional
   migration code, fixtures for every supported version, rollback policy, and support windows.
3. **Require an explicit migration command.** Keeps normal open simple and makes the mutation
   intentional, but adds an operator-facing surface and S7 documentation obligations.

**Recommendation:** alternative 1 until a concrete DDL change supplies real migration
requirements; then decide between 2 and 3 with that change's data and availability constraints.
Do not build generic migration machinery speculatively for this 0.25–0.5 day gate item.

**Cost to defer:** none for the current add-only journal transition, which DG.A verifies.
A future existing-table change cannot land until it either preserves compatibility or brings a
ratified, versioned upgrade design with production-path fixtures.

## Design Freeze

The executor may not quietly change these invariants:

- DG.A is test-only and modifies exactly `host/store/journal_test.go`.
- The canonical expected DDL is hardcoded and independent of `schemaSQL` and actual query output.
- The manifest names exactly seven current durable tables and includes exact journal DDL.
- Expected/actual normalization collapses whitespace only; it does not erase SQL structure.
- The historical fixture is copied verbatim from the pinned `8133573` artifact and is not
  hand-authored or updated to follow current DDL.
- The historical-store test uses real `store.Open` on a temporary-file store, not direct
  execution of the current schema as a substitute.
- The six historical tables and the newly created journal have separate assertions.
- The source SHA pin, same-source equality, and explicit journal deletion are retired.
- A manifest edit is allowed for an intentional DDL change, but it must not green an unapplied
  existing-table edit; AC2 is the controlling evidence.
- No production mismatch response, `user_version`, migration registry, migration CLI, or schema
  rewrite lands without human ratification.
- No `.ail` sketch is needed; `scripts/verify_ail.sh` remains untouched.
- Every named mutation states its edited file and production/test classification; test-side
  mutations are never presented as kernel proofs.
- All test processes have the ceilings stated in DG.A; no unbounded retry, polling, or migration
  loop is introduced.

## What these gates CANNOT fail

These limits are part of the design, not residual fine print:

- They cannot stop a developer from intentionally editing both production DDL and the canonical
  manifest. D1 is a review-visible change detector, not an authorization system.
- They cannot prove that a new constraint or column has correct application semantics.
- They cannot prove data preservation for a future rewrite beyond the one sentinel row; each
  migration needs fixtures for the data it transforms.
- They cannot detect schema drift in an already-deployed store at runtime. DG.A adds tests only.
- They cannot decide whether a mismatch should fail, migrate automatically, or use an external
  tool. OD-3 and OD-4 own that policy.
- They cannot verify indexes, triggers, views, or PRAGMA settings because `tableDDL` selects
  `type='table'` only. None are part of the current schema.
- They cannot prove compatibility with every historical database. The fixture covers the single
  pre-journal shape named by this item.
- They cannot detect semantically equivalent DDL whose materialized SQL differs only in a way the
  whitespace normalizer preserves; such a change reds and requires review even if SQLite behavior
  is equivalent.
- They cannot detect a semantic difference erased solely by whitespace normalization. Current
  schema constraints do not depend on whitespace token boundaries; any future use that does must
  replace the normalizer with a parser or exact comparison.
- They cannot make `CREATE TABLE IF NOT EXISTS` apply an edit to an existing table.
- They cannot validate the Go toolchain. Verification must run under the controller's known-good
  conditions; sandbox results remain uninformative.
- The test-side mutations in AC3–AC5 cannot prove a production kernel property. They only show
  that the proposed gate mechanism discriminates.

The essential promise is narrower and testable: an unreviewed journal DDL change goes red, and a
reviewed existing-table DDL change cannot become green merely by updating a literal while the
old-store path still silently fails to apply it.

---

## Quorum verification log (pick-time quorum, iteration 41)

Two rounds, the maximum the charter allows per revision cycle. **Both BLOCKED.** Both reviewers
were **present in both rounds** (`absent_reviewers: []`), so neither round was an N−1 degrade —
`--max-cost-usd` was raised to `0.35` up front precisely to prevent the iter-36 defect where
enriching a doc silently prices out a reviewer. Total metered cost **$0.1155**.

**Round 1** — `w-ddl-gate-teeth-2026-07-30T04-46-11Z.json` ($0.0547). BLOCKED 2/2.

- `gemini-3-1-pro` (reject): the Conflict Surface's claim that `tableDDL` *"excludes internal SQLite
  tables"* is load-bearing for D1's exact seven-table match but absent from the Premise Verification
  Log. Also: Go map iteration is randomised, so D1 needs deterministic diffing.
  → **ACCEPTED, but the reviewer's proposed WORDING was REFUTED by measurement.** It asked to log
  `name NOT LIKE 'sqlite_%'` as the filter. Measured on a materialized store, `sqlite_master` holds
  **14** rows — the 7 tables plus **7 `sqlite_autoindex_*` rows of `type='index'`** — so
  `count(type='table' AND name NOT LIKE 'sqlite_%')` = 7 and `count(type='table')` = 7 are
  **identical**: the name filter excludes **nothing today** and the live limb is `type='table'`.
  Logging the dormant limb as the mechanism would have reproduced this mission's own signature
  defect (iter-38: *a gate can have teeth from one mechanism while the mechanism it is documented as
  having is dead*). The row now distinguishes live from dormant and says why the dormant guard stays.
  Determinism accepted in full: sorted names in both legs, failures name the offending table.
- `gpt5-6-sol` (reject): D4 preserves a production silent fallback; also, `preJournalSchemaV0`'s
  provenance and the durable-table enumeration lacked verification-log evidence.
  → The two evidence asks were **ACCEPTED and MEASURED, not asserted**: `d5774eb` added the journal
  table, `8133573:host/store/schema.sql` is the last pre-journal schema (sha256 `35f09862e2…`,
  exactly 6 `CREATE TABLE`, materializing exactly the 6 named tables), and the fixture is now sourced
  **verbatim from that pinned artifact** instead of hand-authored. The direction half → **OD-5**.

**Round 2** — `w-ddl-gate-teeth-2026-07-30T04-54-21Z.json` ($0.0608). BLOCKED 2/2.

- `gemini-3-1-pro` (reject): a **new** unverified premise — that the 6 historical tables' DDL in
  `8133573` still matches HEAD's after normalization. If it did not, D2 would **RED on baseline** and
  the design would be broken as written.
  → **A genuinely load-bearing gap that neither the controller nor the designer had checked.
  MEASURED after round 2: all 6 identical, zero drift, so the premise HOLDS**, D2 is green on
  baseline, and the reviewer's own contingency (a separate historical manifest) is not needed. Logged
  as two Premise Verification Log rows. Its second ask — `M4`/`M5` are undefined references — was
  correct and is fixed by the label key above.
- `gpt5-6-sol` (reject): the same direction objection, sharpened to *"axiom noncompliance, not
  merely an out-of-scope enhancement"*, with a five-fixture production milestone as the proposed fix.
  → **NOT resolvable by this loop → OD-5, parked for the human.**

**Why no third round and no carve-out.** The charter caps revision cycles at two rounds and forbids
re-running them in search of a pass. The narrow-refinement carve-out would permit a bounded
controller-authored 2nd revision only if **every** remaining blocking objection were
non-directional with a reviewer-authored fix; `gpt5-6-sol`'s is directional. So the item parks
`needs-human-review` — which is the charter working as designed, not a failure of the round.

**What the two rounds bought, stated plainly.** Four objections, **three of which were factually
correct and improved the design** (fixture provenance, durable-table enumeration, the baseline-drift
premise, determinism), **one of which was factually wrong in its prescribed wording** and would have
logged a dead mechanism as live had it been adopted on authority. That ratio is the argument for
running the quorum and for measuring what it prescribes.


# w-ddl-gate-teeth — make schema drift observable

**Status:** **DG.A ROUTABLE — OD-5 RATIFIED by Mark 2026-08-03 (attended).** ~~PARKED —
`needs-human-review`~~. The pick-time quorum ran its full two rounds and **BLOCKED both times** on a
directional objection the controller could not resolve; the human has now resolved it.

> **RATIFICATION (attended, recorded in the charter STATUS of 2026-08-03; transcribed here by the
> controller at iteration 43 — this is a record of a human decision, not a design change).**
> Mark's words: *"4d RATIFIED — fix NOW: SQLite `user_version` pin failing LOUD on binary↔store
> schema mismatch + de-fang the sha256 self-bypass; the frozen-store touch is ratified (the
> guardrail collision resolves the S6 way: a gate that invites its own bypass is not a gate)."*
>
> **How that maps onto this document, stated precisely, because it does not match any one of OD-5's
> three listed alternatives.**
> - **OD-5 is answered AGAINST this document's D4 position and WITH `gpt5-6-sol`.** The frozen-kernel
>   deferral does not hold: the loop is now authorized to change which on-disk stores `store.Open`
>   accepts. The reviewer was right, and it is recorded here that the reviewer was right.
> - **OD-3 is ratified as its alternative 1** — establish a version, fail LOUD on an unsupported or
>   un-upgraded store. **But ratifying the DECISION is not the same as having a DESIGN.** OD-3's own
>   text requires treatment of legacy version 0 to be specified and the fresh / supported / legacy /
>   future cases to be proven with fixtures. None of that exists in this document: there are no ACs,
>   no named mutations, and no milestone for it. **It is therefore authorized-but-undesigned and is
>   carried as `DG.B`, which a designer must specify before any executor may touch `store.Open`.**
>   Building it headlessly from a three-alternative decision packet is exactly the fabrication this
>   mission's Gate-2 discipline forbids.
>
>   > **STALE AS OF ITERATION 44 — READ THIS BEFORE ACTING ON THE PARAGRAPH ABOVE.** The block you
>   > are reading is a **historical transcript of iteration 43's state**, retained per this
>   > document's convention of never rewriting its own record. **DG.B is now fully designed**:
>   > D5–D9, AC6–AC14, nine named mutations, the fixture set, and a sized milestone all exist below,
>   > and the design passed a two-round quorum plus a narrow-refinement carve-out. An executor may
>   > implement DG.B from this document. (Flagged by `gemini-3-1-pro` as its round-2 non-blocking
>   > catch: an executor reading top-down could otherwise abort, believing DG.B still out of scope.)
>   > What remains human-gated is **rollout**, not implementation — see `4d/OD-7`.
> - **"De-fang the sha256 self-bypass" is ALREADY what DG.A specifies** and needs no new design:
>   **AC4** deletes the source-SHA pin, the same-source before/after equality and the
>   `delete(..., "journal")` outright, and **AC2**'s `MUT-EXISTING-DDL-CHANGE-REMANIFESTED` is the
>   discriminating control proving the replacement manifest **cannot** be silenced the same way — a
>   developer who legitimately re-manifests a `store_heads` edit still gets a red from the
>   historical-store test. That is the S6 property Mark names: the gate no longer invites its own
>   bypass.
>
> **Consequently: DG.A is routable NOW, unchanged, exactly as specified below** (it is the de-fang,
> and it is the evidence DG.B's design will need). **DG.B is queued as the next item.** Splitting
> this way delivers the ratified de-fang immediately without inventing a durability-kernel
> acceptance contract that no reviewer has seen.

See the Quorum verification log at the end.  
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

| **DG.B / V-A:** `openSQLite` is the repository's single production SQLite driver-open path; `Open` passes `applySchema=true` and `OpenReadOnly` passes `false` | Controller command at HEAD `6246ee6`: `grep -rn 'sql\.Open\|modernc.org/sqlite' --include='*.go' host/ \| grep -v _test.go`, plus first-party read of `host/store/store.go:191-268` | One production `sql.Open("sqlite", dsn)` at `store.go:246`; calls at lines 193/208 pass true and line 236 passes false | CONFIRMED BY CONTROLLER |
| **DG.B / V-B:** nothing inspects a database before `openSQLite` applies schema | Controller first-party read of `host/store/store.go:243-268` | Order is `sql.Open`, `SetMaxOpenConns(1)`, `PRAGMA foreign_keys = ON`, then (when requested) `db.Exec(schemaSQL)`; no intervening inspection | CONFIRMED BY CONTROLLER |
| **DG.B / V-C:** the current schema leaves `user_version` at 0 | Controller command at HEAD `6246ee6`: `sqlite3 probe.db < host/store/schema.sql; sqlite3 probe.db "PRAGMA user_version;"`; same call then set `PRAGMA user_version=42` and re-read it | Current schema returned `0`; known-positive write/read control returned `42` | CONFIRMED BY CONTROLLER |
| **DG.B / V-D:** no repository host code reads or writes `user_version` | Controller same-call negative/control: `grep -rn "user_version" host/`; then `grep -rn "PRAGMA" host/` | First command empty; known-positive control returned `host/store/store.go:257` (`PRAGMA foreign_keys = ON;`) | CONFIRMED BY CONTROLLER |
| **DG.B / V-E:** production code performs no `sqlite_master` inspection | Controller same-call negative/control: `grep -rn "sqlite_master" --include='*.go' host/ \| grep -v '_test.go'`; then the same pattern restricted to `*_test.go` | Production command empty; known-positive test control returned `host/store/journal_test.go` | CONFIRMED BY CONTROLLER |
| **DG.B / V-F:** current `schema.sql` consists of seven table creates and materializes seven tables | Controller source counts and sqlite3 materialization at HEAD `6246ee6` | `grep -c 'CREATE TABLE'` = 7, `grep -c 'CREATE'` = 7; `sqlite_master` table count = 7 | CONFIRMED BY CONTROLLER |
| **DG.B / V-G:** the production blast radius outside package `store` is one caller, while tests call `Open` broadly | Controller command: `grep -rn "store\.Open(\|store\.OpenReadOnly(" --include='*.go' . \| grep -v _test.go \| grep -v '^./host/store/'`; known-positive control repeated including tests | Production result only `host/daemon/daemon.go:332`; including tests returned 9 files across broker, daemon, registry, and replay | CONFIRMED BY CONTROLLER |
| **DG.B / V-H:** there is no migration implementation to reuse | Controller repository inspection of `host/` | `journal.go` contains DML only; journal reaches old stores through add-only `CREATE TABLE IF NOT EXISTS`; no migration machinery found | CONFIRMED BY CONTROLLER |
| DG.B can wrap the current schema and version write in one SQLite transaction without a transaction-prohibited statement in `schema.sql` | Designer command at HEAD `6246ee6`: `rg -n '^(PRAGMA\|VACUUM\|ATTACH\|DETACH\|BEGIN\|COMMIT)\|CREATE ' host/store/schema.sql` | Exactly seven `CREATE TABLE IF NOT EXISTS` hits at lines 13, 23, 36, 50, 60, 69, 81; no listed transaction-control or transaction-prohibited statement | CONFIRMED |
| The proposed version-1 fixture has an exact landed artifact whose schema equals HEAD | Designer same-call provenance/hash check: `rev=ad619d8; git show "${rev}:host/store/schema.sql" \| shasum -a 256; shasum -a 256 host/store/schema.sql; git show --no-patch --format='%H %s' "$rev"` | Artifact and HEAD both SHA-256 `13893a296c394cc2c5b3997e8fff729467dc9ac83a03b458796634aa52fb5436`; commit is `ad619d81d8815f5db9ee2de16fcd19cacb1f3c6b` (`feat(4d): ... DG.A ... (#33)`) | CONFIRMED |
| DG.B necessarily collides with the landed pre-journal success test and has a bounded set of current direct test callers to audit | Designer commands at HEAD `6246ee6`: `rg -n 'preJournalSchemaV0\|TestOpenAddsJournalAndDetectsStalePreJournalDDL' host/store/journal_test.go`; `rg -n 'store\.Open\(|store\.OpenReadOnly\(' --glob='*.go' .` | Fixture at line 630 and success test at line 822; direct calls occur in store tests and 17 external-package test sites, plus the one production daemon call | CONFIRMED |
| `run_bounded` is the repository precedent that kills a whole process group and returns 124 on a wall-clock expiry | Designer read of `scripts/verify_ail.sh` Leg 0 | Python uses `start_new_session=True`, `wait(timeout=t)`, `os.killpg(..., SIGKILL)`, and `sys.exit(124)` | CONFIRMED |
| `4d/OD-7` is free in the mission-global planned-doc namespace | Designer same-call negative/control at HEAD `6246ee6`: `{ rg -n '^### OD-' design_docs/planned/; printf 'known-positive headings='; rg -l '^### OD-' design_docs/planned/ \| wc -l; }` | Enumerated only `4d/OD-3`, `4d/OD-4`, `4d/OD-5`, and `w-bench-load-confound/OD-6`; known-positive control reported 2 files containing OD headings; no `OD-7` | CONFIRMED |

**Controller measurements taken AFTER the designer returned, through the PINNED pure-Go driver
(iteration 44).** Rows V-A…V-H above were measured with the `sqlite3` C CLI 3.51.0. That is the
wrong instrument for a claim about this kernel: `host/store` runs on `modernc.org/sqlite v1.54.0`,
and iteration 43's planner already caught this exact instrument gap (its finding **F4**). The four
rows below re-measure the load-bearing premises through the driver the code actually uses, in a
scratch module pinned to the same version. They exist because D7's atomicity contract and D5's
freshness predicate were **asserted by the design and verified by nobody**.

| Claim | How verified | Result | Verdict |
|---|---|---|---|
| **DG.B / V-I (decides whether AC7 is achievable at all):** `PRAGMA user_version` is TRANSACTIONAL, so D7's induced version-write failure genuinely leaves a fresh store at version 0 | Controller Go probe on `modernc.org/sqlite v1.54.0`: `BEGIN`; `tx.Exec(schemaSQL)`; `tx.Exec("PRAGMA user_version = 1;")`; read in-tx; `tx.Rollback()`; re-read on the same handle. Known-positive control in the same run: set `user_version = 5` OUTSIDE any transaction and re-read | In-tx `user_version=1`; **after rollback `user_version=0` and application-objects=0**; control returned `5`, so the post-rollback `0` is a real rollback and not an instrument stuck at zero | **CONFIRMED — D7 and AC7 are buildable as written.** Had the pragma escaped the transaction, AC7's atomicity criterion would have been unachievable and the milestone would have shipped a vacuous gate |
| **DG.B / V-J:** the multi-statement `schemaSQL` can execute inside a single `database/sql` transaction on this driver, as D7 requires | Same probe: `tx.Exec(string(schemaSQL))` then count tables inside the transaction | ACCEPTED; 7 tables visible in-transaction | CONFIRMED |
| **DG.B / V-K (the iteration-41 `sqlite_%` scar, re-checked because D5 reuses that predicate):** is D5's `name NOT LIKE 'sqlite_%'` limb LIVE or a dead mechanism logged as live? | Controller probe on the materialized current schema: total `sqlite_master` rows vs rows surviving the limb, plus a `GROUP BY type` breakdown. Known-positive control in the same call: rows whose name **does** match `sqlite_%` | Fresh DB: 0 and 0. After `schemaSQL`: **14 total, 7 surviving — the limb excludes 7 rows**, all of `type=index` (`sqlite_autoindex_*` from the PRIMARY KEYs); the control matched exactly those 7 | **LIVE, and materially different from iteration 41's case** — there the same predicate excluded **zero** rows because `type='table'` had already dropped them. D5 applies no type filter, so the limb does real work. **Honest limit**: it is defensive, not discriminating — an `sqlite_autoindex_*` row cannot exist without its table, so no reachable store state makes the exclusion change the fresh/legacy verdict. **SUPERSEDED IN PART BY V-S, and the gap was in this row's own control design**: this row proved the limb *fires* (a known-positive matched) but never tested its **exclusion boundary** — whether a name that should SURVIVE the filter actually does. `gpt5-6-sol` caught exactly that, and it was hiding a real defect |
| **DG.B / V-S (the round-2 carve-out fix, measured before adopting it):** does the unescaped `NOT LIKE 'sqlite_%'` misclassify a reachable application object, and does the `ESCAPE` form fix it without losing the real exclusions? | Controller Go probe on the PINNED `modernc.org/sqlite v1.54.0`, exactly the evidence `gpt5-6-sol`'s `proposed_fix` demanded: attempt to create `sqlite_internal_probe` and `sqliteX_probe`; then count a store whose only application object is `sqliteX_probe` under both predicates; then a two-sided control on the pattern itself | `CREATE TABLE sqlite_internal_probe` → **REJECTED**, `object name reserved for internal use`; `CREATE TABLE sqliteX_probe` → **ACCEPTED**. That store's `sqlite_master` holds `sqliteX_probe` + `sqlite_autoindex_sqliteX_probe_1`. **Unescaped predicate counts 0** (→ classified FRESH); **`ESCAPE` predicate counts 1**, naming `sqliteX_probe` (→ classified LEGACY). Two-sided control: `'sqlite_autoindex_…' LIKE 'sqlite\_%' ESCAPE '\'` = **1** (real internals still excluded) and `'sqliteX_probe'` = **0** (no longer wrongly excluded) | **CONFIRMED — the objection was right and the defect was reachable.** Under the unescaped form DG.B would have initialized and version-stamped a store containing real application data. The `ESCAPE` form is adopted verbatim from the reviewer's `proposed_fix` |
| **DG.B / V-L:** two mechanical constraints the executor will hit | Same probe suite | `PRAGMA user_version = ?` is **REJECTED** (`near "?": syntax error`) — the statement cannot use a bind parameter and must be built as a literal, which is why D7 specifies the fixed pragma text; on a `mode=ro` handle the version is **readable but NOT writable** (`attempt to write a readonly database`), which is what makes D8's read-only enforcement implementable without a write path | CONFIRMED |
| **DG.B / V-M:** the C-CLI readings V-C/V-F are not contradicted by the real driver | Same probe: brand-new DB pre-schema, then post-schema, on `modernc.org/sqlite` | Pre-schema `user_version=0`, tables=0; post-schema `user_version=0`, tables=7; control set-to-7 read back 7, and the value **persisted across close/reopen** (so the pin lives in the file header and is durable) | CONFIRMED — both instruments agree, and the persistence fact is new |
| **DG.B / V-N (signed-domain control):** an ordinary positive `user_version` is readable through the pinned driver | Controller Go probe on `modernc.org/sqlite v1.54.0`: execute literal `PRAGMA user_version = 1`, then query `PRAGMA user_version` | Set `1`; read back `1` | **CONFIRMED — known-good reader control for V-O through V-Q** |
| **DG.B / V-O:** negative `user_version` values are reachable at both ends of the signed range | Controller same probe and driver: set/read literals `-1` and `-2147483648`, alongside V-N's known-positive control | `-1` read back `-1`; `-2147483648` read back `-2147483648` | **CONFIRMED — D5 must classify negative values** |
| **DG.B / V-P:** the maximum signed 32-bit `user_version` is reachable without truncation | Controller same probe and driver: set/read literal `2147483647`, alongside V-N's known-positive control | Set `2147483647`; read back `2147483647` | **CONFIRMED — the usable positive schema-version ceiling is INT32_MAX** |
| **DG.B / V-Q:** values above the signed 32-bit ceiling silently truncate to zero rather than error | Controller same probe and driver: set/read literals `2147483648` and `4294967296`, alongside V-N's known-positive control | Both statements succeeded; both values read back `0` | **CONFIRMED — an out-of-range version can erase its own evidence by becoming the legacy/fresh sentinel** |
| **DG.B / V-R:** a negative `user_version` is durable, not a transient connection value | Controller pinned-driver probe: set `-1`, close the database, reopen it, and query again; V-N is the reader control | After close/reopen, read back `-1` | **CONFIRMED — a store can arrive at `Open` with a persistent negative version** |

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
| **DG.B:** `host/store/store.go` — `Open`, `OpenReadOnly`, and `openSQLite` | This is frozen durability-kernel code and the only driver-open choke point | The frozen-store touch is ratified by Mark for `4d/OD-3`; keep classification and enforcement in this choke point and preserve writer-lock release on every error |
| **DG.B:** `host/store/schema.sql` and new production constants/error types | The schema needs a nonzero identity, but embedding enforcement only in the schema would apply it before legacy classification | Keep `schema.sql` DDL-only. Define `currentSchemaVersion = 1` in `store.go`; write it explicitly only inside the fresh-store transaction after the DDL succeeds |
| **DG.B:** `host/store/schema_version_test.go` (new) | Owns total-domain fixtures, typed-error assertions, atomicity, and the independent version-to-DDL ledger | Add one focused test file; freeze version-1 DDL as an artifact-derived test fixture independent of `schemaSQL` and `currentSchemaVersion` |
| **DG.B:** landed `host/store/journal_test.go` fixture/test from DG.A | Its real `Open` success assertion for a version-0 pre-journal store contradicts the newly ratified rejection | Preserve the fixture and DG.A history; change only the landed test's post-open expectation to the typed legacy rejection, while D1 continues to own current exact DDL. Do not stamp the fixture as version 1 |
| **DG.B:** current `Open` test corpus and daemon caller | Fresh `Open` calls should self-initialize; raw sqlite fixtures and reopen paths may expose version assumptions | Audit all direct calls named in V-G. Fresh/reopen tests require no bypass helper; any raw current-schema fixture must explicitly write version 1 as part of fixture construction, while an intentionally legacy fixture must remain 0 and assert rejection |
| **DG.B:** writer-lock discipline | A schema mismatch occurs after the file-backed writer lock is acquired | Return through `Open`'s existing error path so the lock is released; add a retry-after-mismatch test that proves no lock descriptor is stranded |
| **DG.B:** `scripts/verify_go.sh` and CI | Its commands are not the bounded mechanism specified for mutation execution | DG.B edits no script. Executor/controller invocations wrap focused tests, package tests, and the repository gate with a `run_bounded`-equivalent hard wall-clock deadline; no bare `$AILANG_BIN --version` remedy is allowed |

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

**DG.B supersession of D2/D3's success-path mechanism.** The paragraphs above remain DG.A's
landed design and historical rationale. Under DG.B the same fixture is rejected before schema, so
no post-open DDL comparison runs and those paragraphs must not be read as DG.B acceptance
behavior. AC2 records the supersession explicitly: AC6 owns the stronger production property that
an existing unupgraded store never yields a usable current `Store`, while AC10 owns the
edit/re-manifest DDL discriminator.

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
ratified guardrails when they point opposite ways. ~~Until OD-5 is answered, **OD-5 governs and no
milestone here may be routed.**~~

> **RESOLVED 2026-08-03 (Mark, attended) — AND D4's POSITION ABOVE DID NOT SURVIVE.** The human
> ratified the frozen-store touch and ordered the `user_version` pin *now*, which is `gpt5-6-sol`'s
> direction, not this section's. **D4 above is retained verbatim as the position that was
> overruled** — it is not edited to look prescient, because a design doc that quietly rewrites its
> own losing argument teaches the next reader nothing. What D4 got right: DG.A is strictly additive
> and did not need to wait. What D4 got wrong: it treated "pre-existing at HEAD" as sufficient
> reason to leave a *measured* production fail-open in the completed state, and the axiom does not
> have that exemption. The production response is now **DG.B** (authorized, undesigned).

Changing that requires a durable schema-version contract. `PRAGMA user_version` is the smallest
candidate, but setting or enforcing it changes what stores the kernel accepts. Automatically
migrating is larger still. OD-3 and OD-4 therefore remain human decisions.

The non-ratification milestone can land first because it changes no accepted store format and no
runtime behavior. Its red result becomes evidence required by the first future DDL-change design;
it is not permission for an executor to invent a migration.

### D5. DG.B — classify before schema application; version 1 is current

**DG.B addition (iteration 44).** Define the unexported production constant
`currentSchemaVersion = 1`. Version 1 names exactly the seven-table schema frozen by DG.A.
SQLite's observable `user_version` domain is the **signed 32-bit range
`-2147483648 .. 2147483647`**, as V-N through V-R establish; it is not an unconstrained
positive-integer field. Version 0 is the unversioned legacy namespace, not a supported schema
version. A database is
**genuinely fresh** only when both conditions hold before `schemaSQL` runs: `PRAGMA user_version`
is 0 and

```sql
SELECT count(*) FROM sqlite_master
WHERE name NOT LIKE 'sqlite\_%' ESCAPE '\';
```

returns 0. **The `ESCAPE` clause is load-bearing, not decoration** (carve-out revision, iteration 44;
`gpt5-6-sol` round 2). In SQL `LIKE`, `_` is a **single-character wildcard**, so the unescaped
`NOT LIKE 'sqlite_%'` excludes every name beginning `sqlite` **plus any character** — not the
literal reserved prefix `sqlite_`. SQLite reserves and refuses to create `sqlite_`-prefixed names,
but `sqliteX_probe` is **not** reserved and creates fine, so a version-0 store whose only
application object is `sqliteX_probe` counted **0** under the unescaped predicate and would have
been classified **fresh** — then initialized and stamped version 1 over real data. Measured, with
controls, in V-S. Either a literal-prefix test works; this design takes the `ESCAPE` form (the
reviewer also offered `NOT GLOB 'sqlite_*'`) because `LIKE` is ASCII-case-insensitive and therefore
matches SQLite's own case-insensitive treatment of the reserved prefix, and because it is the form
measured in V-S.

Inspect all application schema-object types, not tables alone: a version-0 database
containing only an application index, trigger, or view is not fresh. SQLite-internal objects are
excluded because they are engine bookkeeping rather than evidence of an AILANG World schema.
This is a structural definition, not a file-size or file-existence heuristic: a pre-created empty
file is fresh, while any version-0 application object makes the store legacy even if none of the
seven current tables exists.

Move this classification into `openSQLite` after the existing first connection-establishing
`PRAGMA foreign_keys = ON` and before the current `db.Exec(schemaSQL)`. It applies on both writer
and read-only paths. Nothing may execute `schemaSQL`, repair a table, or set a version before the
classification query. The extra cost is two scalar queries on the already-single connection per
open; no scan of table contents and no retry or wait is introduced.

The five writer cases below are total and mutually exclusive over every value the field can
actually hold:

| Pre-schema observation | `Open` result |
|---|---|
| version 0, zero application objects (**fresh**) | Atomically create schema and set version 1 as D7 specifies; return a usable store |
| version 0, one or more application objects (**legacy**) | Close and return `*LegacySchemaVersionError`; execute no DDL and write nothing |
| version 1 (**supported**) | Apply idempotent `schemaSQL`, leave version unchanged, and return a usable store |
| version 2 through 2147483647 (**future**) | Close and return `*FutureSchemaVersionError`; execute no DDL and write nothing |
| version -2147483648 through -1 (**invalid**) | Close and return `*InvalidSchemaVersionError`; execute no DDL and write nothing |

If a later binary supports more than one positive version, that is a new support-window/migration
design; DG.B supports exactly 1 and must not generalize speculatively.

### D6. Typed, distinct, operator-actionable failures

Add exported error types so callers and tests do not parse prose:

- `LegacySchemaVersionError { Path string; Found, Current int }`, returned only for non-empty
  version-0 stores. Its stable message shape is
  `store: schema version legacy: <quoted path> has user_version 0 with application schema; binary requires 1; refusing to modify`.
- `FutureSchemaVersionError { Path string; Found, Current int }`, returned for versions above
  current. Its stable message shape is
  `store: schema version future: <quoted path> has user_version <n>; binary supports 1; use a compatible binary`.
- `InvalidSchemaVersionError { Path string; Found, Current int }`, returned for negative versions.
  Its stable message shape is
  `store: schema version invalid: <quoted path> has negative user_version <n>; binary requires 1; refusing to modify`.
- `UninitializedReadOnlyStoreError { Path string }`, returned when a read-only open sees the
  structurally fresh version-0 case. Its stable message shape is
  `store: schema uninitialized: read-only store <quoted path> has no application schema; open writable once to initialize`.

Wrapping driver/query failures keeps the existing `store: open ...` context. The four policy
errors above are returned directly after closing the handle. Legacy, future, and invalid must
never share a generic “mismatch” type or identical text: legacy means no version contract was
established; future means the store was written by a newer contract; negative means the field is
invalid/corrupt and establishes neither. None of the messages tells an operator to paste a new
version into the database. In particular, legacy and invalid text says “refusing to modify,” not
“run `PRAGMA user_version=1`.”

### D7. Fresh initialization is one transaction; failure leaves version 0

For a writable fresh store, begin one transaction on the already-open single connection, re-read
both `user_version` and the application-object count inside that transaction, execute `schemaSQL`,
execute the literal `PRAGMA user_version = 1`, and commit. The in-transaction re-read closes the
gap between classification and mutation without adding a wait loop. Any begin, re-check, schema,
version-write, or commit error rolls back, closes the database, and returns a stage-specific
wrapped error. If the version write fails, schema creation must not commit; if commit outcome is an
error, `Open` fails rather than returning an unpinned store.

Factor only the transaction body into an unexported helper that accepts a version-write callback.
The production caller supplies a closure executing the fixed pragma; the atomicity test calls the
helper directly with a callback returning a sentinel error after DDL execution, then rolls through
the same rollback path. Do not use a mutable package-global test hook, environment switch, or DSN
escape hatch: those would add a shipping bypass or race between parallel tests.

For a supported version-1 writer, execute `schemaSQL` as today, then re-read `user_version` and
require it still equals 1 before returning. Do not rewrite it on every open. That final assertion
makes the marker part of acceptance behavior rather than inventory. DG.B adds no migration,
rollback policy, support window, or operator migration command; `4d/OD-4` remains open and its
recommended fail-loud alternative bounds this milestone.

### D8. `OpenReadOnly` enforces the same compatibility boundary

`OpenReadOnly` accepts exactly version 1. It performs the same read-only classification before
returning a `Store`, applies no schema, takes no lock, and writes no pragma. It returns
`*LegacySchemaVersionError` for non-empty version 0, `*FutureSchemaVersionError` for greater than
1, `*InvalidSchemaVersionError` for negative values, and `*UninitializedReadOnlyStoreError` for
the empty version-0 case.

Accepting incompatible stores “for observation” would be a silent capability downgrade: every
method on the returned `Store` is still compiled against the current tables and columns, and the
API has no limited-view type with which to express safety. Refusal also prevents a future-version
store from being misread through older query semantics. Operators who need forensic SQLite access
can use SQLite tooling; DG.B does not turn the normal typed store API into that surface.

### D9. Independent version ledger and fixture transition

Add `host/store/schema_version_test.go` with a frozen `schemaV1SQL` copied from
`ad619d8:host/store/schema.sql` and provenance recorded in its comment. It is independent of the
embedded production `schemaSQL`, the production `currentSchemaVersion`, and DG.A's editable
canonical manifest. A test-only literal `expectedCurrentSchemaVersion = 1` checks the production
constant. Materializing `schemaV1SQL` and current `schemaSQL` into separate databases and comparing
their exact normalized application DDL proves that version 1 still names the frozen schema. A DDL
change therefore requires two review-visible acts: bump the production version and add a new
independent frozen fixture/expectation; re-manifesting DG.A alone cannot silence this gate.

Fixtures cover fresh empty version 0, supported version 1, legacy non-empty version 0 (use the
landed `preJournalSchemaV0` artifact), future version 2, invalid versions -1 and INT32_MIN, and —
added by the iteration-44 carve-out revision — a version-0 store whose ONLY application object is
`sqliteX_probe`, which AC14 requires to classify as legacy through both open paths.
Each case is exercised through both `Open` and, for file-backed stores, `OpenReadOnly`; rejection
fixtures snapshot object names and `user_version` before/after to prove no mutation. A forced
version-write failure must prove the fresh transaction leaves zero application objects and version
0. The writer mismatch fixture is then corrected out-of-band only inside the test setup and
reopened to prove the writer lock was released; that setup action is a lock cleanup probe, not a
supported operator remedy.

Most existing tests call `Open` on an empty temporary file or `:memory:` and receive version 1
through the real fresh path; they must not call a helper that stamps every database. Tests that
construct a raw current store may set version 1 in that fixture only, adjacent to the independent
DDL construction. Intentionally legacy/future fixtures remain unstamped/mismatched. The landed
pre-journal test changes from expecting journal creation to expecting the typed legacy failure;
the fixture itself stays byte-for-byte historical. This transition keeps the test corpus passing
without a global “set version 1” bypass that would make legacy rejection untestable.

Existing deployed version-0 stores are deliberately rejected rather than silently blessed. DG.B
does not have authority to decide whether the rig's actual stores may be discarded, reconstructed,
or certified and stamped; that deployment transition is isolated as `4d/OD-7` below. DG.B may be
implemented and tested in 0.5–1.0 day, but rollout to a rig containing valuable version-0 data is
blocked until that decision is answered.

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

> **SUPERSEDED BY DG.B; original DG.A criterion and mutation retained verbatim above for the audit
> trail.** Once DG.B lands, the real `Open` in this fixture returns
> `*LegacySchemaVersionError` before DDL for both the baseline and
> `MUT-EXISTING-DDL-CHANGE-REMANIFESTED`; the mutation changes neither `user_version` nor the
> application-object count. The old AC2 mechanism is therefore vacuous under DG.B and is retired
> as an AC2 discriminator. The same edit/re-manifest action is re-pointed to AC10, where the
> independently frozen version-1 ledger reds against changed production DDL.
>
> **Active DG.B form of AC2.** “This existing store was never upgraded” is no longer reachable as
> a successful-open/post-DDL comparison: production refuses every non-empty version-0 store at the
> door, before executing schema. That rejection is stronger than observing stale DDL after a
> successful open because no current typed `Store` can escape for an unupgraded store. **AC6 carries
> this property as a PRODUCTION discriminator.** AC2 retains an independently named production RED
> as a cross-check: **`MUT-LEGACY-REJECTION-BYPASS`** changes the production legacy branch in
> `host/store/store.go` (**PRODUCTION**) to return a usable `Store`; the transitioned
> `TestOpenAddsJournalAndDetectsStalePreJournalDDL` must red because it requires
> `*LegacySchemaVersionError`. It must not compare post-open DDL under DG.B.

### AC3 — the historical fixture is independent and non-vacuous

Owned by **DG.A**, with its DG.B expectation transition. `preJournalSchemaV0` is a verbatim
artifact-derived constant independent of `schemaSQL`; fixture construction independently asserts
exactly six historical names, requires each in the canonical manifest, and asserts journal is
absent before the real `Open` returns the typed legacy rejection. DG.B supersedes the former
post-open journal-creation assertion; it does not supersede fixture non-vacuity.

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

Named RED mutation: **`MUT-UPGRADE-ASSERTION-DEAD`**, whose DG.A behavior remains recorded in the
implementation log, edits `host/store/journal_test.go` (**TEST**) under DG.B to short-circuit the
transitioned `*LegacySchemaVersionError` assertion. Under `MUT-LEGACY-REJECTION-BYPASS`, AC2/AC6's
required RED disappears; mutation review rejects the implementation. This remains a test-side
gate-discrimination probe, not production evidence.

### AC6 — fresh and legacy version 0 are separated before DDL

Owned by **DG.B**. A writable empty version-0 store becomes the exact current schema at version 1,
while the non-empty `preJournalSchemaV0` fixture returns `*LegacySchemaVersionError` with
`Found=0`, `Current=1`, unchanged object names, and unchanged pragma. Neither case is classified
after schema application.

Named RED mutation: **`MUT-LEGACY-AS-FRESH`** changes the production freshness predicate in
`host/store/store.go` (**PRODUCTION**) from `applicationObjects == 0` to
`applicationObjects >= 0`. The legacy fixture must red because `Open` no longer returns the typed
legacy error and/or mutates the fixture. This is the central discriminator for the ambiguity in
V-C/V-E.

### AC7 — fresh schema and version pin commit atomically

Owned by **DG.B**. Fresh initialization executes current DDL and `PRAGMA user_version = 1` in one
transaction. Success exposes seven tables and version 1. An induced version-write failure returns
an error and a separately reopened database has zero application objects and version 0.

Named RED mutation: **`MUT-FRESH-SKIP-PIN`** removes the production transaction's single
`PRAGMA user_version = 1` execution in `host/store/store.go` (**PRODUCTION**) while leaving schema
execution and commit intact. The successful-fresh fixture must red on version 0; it is not enough
that seven tables exist.

### AC8 — supported and future versions have different outcomes

Owned by **DG.B**. Version 1 opens without changing its pragma. The same current DDL marked version
2 returns `*FutureSchemaVersionError` with `Found=2`, `Current=1`, the future-specific stable
message, and no schema mutation. `errors.As` assertions own type identity; exact substrings own the
operator distinction.

Named RED mutation: **`MUT-FUTURE-AS-SUPPORTED`** changes the production supported predicate in
`host/store/store.go` (**PRODUCTION**) from `version == currentSchemaVersion` to
`version >= currentSchemaVersion`. The version-2 fixture must red because `Open` succeeds instead
of returning the future error.

### AC9 — read-only opens enforce without writing or locking

Owned by **DG.B**. `OpenReadOnly` accepts version 1, rejects non-empty version 0 and version 2 with
the same typed legacy/future errors, and rejects an empty version-0 file with
`*UninitializedReadOnlyStoreError`. Before/after snapshots prove all rejection cases retain their
objects and pragma; existing writer-lock coexistence coverage remains green.

Named RED mutation: **`MUT-READONLY-SKIP-VERSION-GATE`** changes the production call in
`host/store/store.go` (**PRODUCTION**) so version classification runs only when `applySchema` is
true. The read-only legacy/future fixtures must red by returning a `Store`, while writer fixtures
remain green.

### AC10 — schema version and schema text are independently pinned

Owned by **DG.B**. `expectedCurrentSchemaVersion = 1` is test-local, `schemaV1SQL` is frozen from
the named landed artifact, and its materialized DDL equals current materialized DDL. Neither
expected value is computed from production `currentSchemaVersion` or `schemaSQL`.

Primary named RED mutation: **`MUT-EXISTING-DDL-CHANGE-REMANIFESTED`**, re-pointed here from AC2,
edits `host/store/schema.sql` (**PRODUCTION**) to add
`mutant_col TEXT NOT NULL DEFAULT ''` to `store_heads` and updates DG.A's canonical manifest
(**TEST**) to the changed DDL. Fresh DDL is accepted and the legacy fixture is still rejected, but
the independently frozen `schemaV1SQL` comparison must red because changed production schema is
still called version 1. This makes AC10 a **PRODUCTION discriminator** after DG.B.

Secondary named RED mutation: **`MUT-V1-LEDGER-DERIVED`** edits
`host/store/schema_version_test.go` (**TEST-SIDE**) to replace the frozen `schemaV1SQL` input with
production `schemaSQL`. Under a one-line production DDL mutation, the required version-ledger RED
disappears; mutation review rejects that implementation. This proves probe independence, not a
runtime kernel property.

### AC11 — corpus transition has no universal stamping bypass

Owned by **DG.B**. Fresh/reopen tests initialize only through real `Open`; raw current fixtures set
version 1 locally; legacy/future fixtures retain their deliberate versions. The landed pre-journal
fixture remains version 0 and its real `Open` assertion expects the typed legacy failure. A writer
mismatch followed by a corrected test fixture reopens successfully, proving lock release.

Named RED mutation: **`MUT-LEGACY-FIXTURE-AUTOSTAMP`** edits the legacy fixture setup in
`host/store/schema_version_test.go` (**TEST-SIDE**) to execute `PRAGMA user_version = 1` after
creating it. The legacy rejection assertion must red. This is a test-wiring probe and is not
evidence that production rejects a genuine version-0 store.

### AC12 — negative versions are invalid and fail distinctly

Owned by **DG.B**. Both `Open` and `OpenReadOnly` reject persistent negative versions, including
-1 and INT32_MIN, with `*InvalidSchemaVersionError`, correct `Found`/`Current` fields, the
invalid-specific stable message, and unchanged objects and pragma. They do not fold the value into
fresh, legacy, or future handling.

Named RED mutation: **`MUT-NEGATIVE-AS-LEGACY`** changes the production classifier in
`host/store/store.go` (**PRODUCTION**) from the negative invalid branch to the version-0 legacy
branch. The -1 writer and read-only fixtures must red on error type and message while the genuine
version-0 legacy fixture remains green.

### AC13 — schema versions stay inside SQLite's non-truncating positive range

Owned by **DG.B**. Production `currentSchemaVersion` and every independently frozen future schema
version are integers in `1 .. 2147483647`. Tests assert both bounds before constructing a pragma
literal. This is mandatory because V-Q shows that 2147483648 and 4294967296 succeed but read back
as 0, which could silently turn a future store into legacy or, when structurally empty, fresh.

Named RED mutation: **`MUT-CURRENT-VERSION-OVERFLOW`** changes production
`currentSchemaVersion` in `host/store/store.go` (**PRODUCTION**) from 1 to 2147483648. The explicit
range assertion must red before any fixture can execute the truncating pragma.

### AC14 — the freshness predicate tests the LITERAL reserved prefix

Owned by **DG.B**. Added by the **narrow-refinement carve-out revision (iteration 44)**, applying
`gpt5-6-sol`'s round-2 `proposed_fix` verbatim: *"Add a writer and read-only version-0 fixture
containing only `sqliteX_probe`; both must classify it as legacy, return `LegacySchemaVersionError`,
and leave objects and `user_version` unchanged."*

A version-0 fixture whose only application object is `sqliteX_probe` (a name SQLite permits,
because only the literal `sqlite_` prefix is reserved — V-S) must classify as **legacy** through
both `Open` and `OpenReadOnly`, returning `*LegacySchemaVersionError`, with object names and
`user_version` snapshotted before and after to prove no mutation. This is the boundary case the
unescaped predicate silently admitted as fresh.

Named RED mutation: **`MUT-PREFIX-WILDCARD-REGRESSION`** reverts the production predicate in
`host/store/store.go` (**PRODUCTION**) from `name NOT LIKE 'sqlite\_%' ESCAPE '\'` to the
unescaped `name NOT LIKE 'sqlite_%'`. The `sqliteX_probe` fixture must red, because `Open` would
classify it as fresh, execute DDL and stamp version 1 instead of returning the legacy error. V-S
measured both readings (0 vs 1), so this mutation is known to discriminate before it is written.

**Honest revised DG.B tally:** **8 of 9** acceptance criteria are production discriminators
(AC6–AC10, AC12–AC14); **1 of 9** is a test-side discrimination probe (AC11). Combined with DG.A's
active DG.B form of AC2, the document has **10 of 14** production-discriminating criteria and
**4 of 14** test-side criteria (AC3–AC5, AC11). AC2's original DG.A mechanism is retained but not
double-counted; AC10 receives its re-pointed edit/re-manifest mutation.

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

**DG.A IS ROUTABLE (OD-5 ratified by Mark, 2026-08-03 — see Status).** ~~DG.A IS NOT ROUTABLE
YET.~~ The directional objection that held it has been resolved by the human **in the reviewer's
favour**, and the milestone body above is now an authorization as well as a specification. It is
unchanged by the ratification: DG.A was always the de-fang half, and the ratification's other half
(the `user_version` production contract) is **DG.B**, which is *not* specified here and must not be
smuggled into this milestone. The paragraph immediately above still binds — anything touching
`host/store/schema.sql`, `store.Open`, accepted on-disk versions, or stored data is **DG.B**, not
DG.A, and an executor that reaches for it is out of scope.

### DG.B — enforce the version-1 fail-loud store contract (~0.5–1.0 day)

**Added in iteration 44; owns AC6–AC13 except the DG.A-owned AC2 transition.** Modify `host/store/store.go`, add
`host/store/schema_version_test.go`, and make the narrow expectation transition in
`host/store/journal_test.go`:

- classify empty/fresh, supported version 1, non-empty legacy version 0, future positive versions,
  and invalid negative versions before any schema application;
- initialize a fresh writer's seven-table schema and version 1 atomically;
- return distinct typed legacy, future, invalid, and uninitialized-read-only errors and release
  any writer lock on every rejection;
- enforce version 1 on `OpenReadOnly` without schema, lock, or write;
- freeze the independent version-1 schema fixture, enforce the signed-32-bit version ceiling, test
  the literal reserved-prefix boundary (AC14), and execute all nine primary named mutations plus
  AC10's secondary test-side probe; and
- audit direct `Open` call sites without adding a universal version-stamping helper.

DG.B is fail-loud only. It does not edit `schema.sql`, add migrations, bless legacy stores, define
a support window, or ship an operator command. On a sound rig, every proposed command is run via
the `scripts/verify_ail.sh` `run_bounded` mechanism (or a byte-identical temporary wrapper): 30
seconds for formatting/static checks, 120 seconds for each focused mutation test, 180 seconds for
`go test ./host/store/ -count=1`, and 600 seconds for `scripts/verify_go.sh`. Each timeout kills the
child process group and reports exit 124; `go test -timeout` is additional defense only. The
`AILANG_BIN --version` performed inside the repository gate is therefore inside the same 600-second
outer wall. One run per mutation and one restored baseline; no retries.

The implementation fits the estimate. Deployment treatment of valuable existing version-0 stores
does not: it is intentionally excluded and routed to `4d/OD-7` rather than hidden as migration work.

## Open Decisions for the human

> Numbering note (controller, iteration 41): these are **OD-3** and **OD-4**, not OD-1/OD-2.
> Iteration 40 already parked **OD-1** (lower the `go.mod` floor 1.26.4 → 1.25.6) and **OD-2**
> (file the go1.26 reproducer upstream at `golang/go`) for the same human, and the charter's
> parked-for-human list is a single global namespace — two live `OD-1`s meaning different
> things is the iter-31 ID-collision defect, where a collision reads as continuity.

### OD-5 — ~~THE CONTROLLING DECISION~~ **ANSWERED 2026-08-03 (Mark, attended): YES — the axiom wins, the frozen-store touch is ratified.**

> **Answer: none of the three alternatives below as written — closest to alternative 2, but without
> holding DG.A hostage.** Mark ratified the kernel touch *and* ordered the `user_version` pin
> **now**, while DG.A (which changes no accepted store format) lands immediately rather than
> waiting for it. So: alternative 1's *sequencing* with alternative 2's *substance*. The three
> alternatives are kept below unedited as the options that were actually put to the human.

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

### OD-3 — ~~add and enforce a `PRAGMA user_version` contract now?~~ **RATIFIED 2026-08-03 (Mark, attended): YES — alternative 1, fail LOUD. → becomes milestone `DG.B`, still to be DESIGNED.**

> The decision is settled; the design is not. This section is a decision packet — three
> alternatives and a recommendation. It has **no acceptance criteria, no named mutations, and no
> fixtures**, and its own recommendation text demands *"specify treatment of legacy version 0 and
> prove fresh, supported, legacy, and future-version cases"* before ratification-class code lands.
> None of that has been written or reviewed. **DG.B is therefore the next item's deliverable: a
> designed, quorum-reviewed milestone for the `user_version` contract.** No executor may implement
> `store.Open` acceptance changes from this section alone.

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

### OD-7 — how are valuable deployed version-0 stores transitioned after DG.B starts refusing them?

**DG.B addition (iteration 44). Question.** Mark ratified changing what `store.Open` accepts, but
the repository cannot establish whether the rig's existing version-0 databases are disposable,
reconstructible, or authoritative. Should rollout recreate them, certify and stamp them through a
separately designed one-shot procedure, or pause until a real migration surface exists?

**Alternatives.**

1. **Back up, discard, and reconstruct version-0 stores from authoritative inputs.** This keeps
   DG.B fail-loud and adds no version-blessing path, but is valid only if all deployed state is
   demonstrably reconstructible.
2. **Design a separate, explicit certification-and-stamp procedure.** It would compare the whole
   application schema and required integrity/data invariants against an independent version-1
   specification before writing version 1. This preserves data but is an operator migration
   surface, outside DG.B, and needs its own ACs, mutations, rollback/backup instructions, and
   ratification.
3. **Delay production rollout of DG.B while keeping its code/gates ready.** Safest when store value
   and reconstructibility are unknown, but leaves deployed binaries fail-open until the decision.

**Recommendation: alternative 3 until the controller inventories the rig, then alternative 1 only
with first-party proof that every store is reconstructible; otherwise alternative 2 as a separate
milestone.** Blindly setting `user_version=1` proves only that a pragma was writable and would
recreate the inert-marker defect. DG.B therefore contains no command or helper that does it.

**Cost to defer:** DG.B can be implemented and mutation-tested, but it must not be rolled onto a
rig that may contain valuable version-0 stores. Already-deployed binaries retain the runtime
ambiguity during that delay.

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
- A manifest edit is allowed for an intentional DDL change, but it must not green a production DDL
  edit still called version 1; AC10 is the controlling discriminator. AC6's production rejection,
  not AC2's superseded post-open comparison, owns the fact that an unupgraded existing store is
  never accepted.
- No production mismatch response, `user_version`, migration registry, migration CLI, or schema
  rewrite lands without human ratification.
- No `.ail` sketch is needed; `scripts/verify_ail.sh` remains untouched.
- Every named mutation states its edited file and production/test classification; test-side
  mutations are never presented as kernel proofs.
- All test processes have the ceilings stated in DG.A; no unbounded retry, polling, or migration
  loop is introduced.
- **DG.B:** `currentSchemaVersion` is 1; version 1 denotes exactly the frozen seven-table schema.
- **DG.B:** `currentSchemaVersion` and every future schema version are constrained to
  `1 .. 2147483647`. V-Q measured that larger literals silently truncate to 0, so date-like,
  namespaced, hash-derived, or otherwise oversized version numbers are forbidden.
- **DG.B:** fresh means version 0 plus zero application schema objects, observed before DDL;
  file existence, file size, or table count alone may not substitute.
- **DG.B:** the application-object predicate MUST test the **literal** reserved prefix —
  `name NOT LIKE 'sqlite\_%' ESCAPE '\'` (or `NOT GLOB 'sqlite_*'`). The unescaped
  `NOT LIKE 'sqlite_%'` is FORBIDDEN: `_` is a single-character wildcard, and V-S measured that it
  classifies a store whose only application object is `sqliteX_probe` as fresh. An executor that
  "simplifies" this predicate by dropping the escape reintroduces a store-corrupting misclassification.
- **DG.B:** non-empty version 0 is legacy, version 2 through INT32_MAX is future, and negative is
  invalid/corrupt; all are rejected through distinct exported types and messages.
- **DG.B:** `OpenReadOnly` accepts only version 1 and never writes, applies schema, or takes a lock.
- **DG.B:** fresh DDL and the version write commit in one transaction; every failure returns no
  store, rolls back when possible, closes the handle, and releases a file-backed writer lock.
- **DG.B:** supported version 1 is re-read after idempotent schema application; the version marker
  is enforced acceptance state, not write-only inventory.
- **DG.B:** `schemaV1SQL` and `expectedCurrentSchemaVersion` are hardcoded test policy independent
  of production `schemaSQL`, production `currentSchemaVersion`, and DG.A's canonical manifest.
- **DG.B:** no universal test helper stamps version 1. Only raw fixtures that intentionally model
  current stores may set it; legacy/future fixtures preserve their versions.
- **DG.B:** no migration registry, support window, automatic migration, operator migration
  command, or manual stamping recommendation enters this milestone. `4d/OD-7` owns rollout.
- **DG.B:** every verification process, including remedies and the repository gate's nested binary
  version check, has the hard wall-clock/process-group ceilings stated in the milestone.

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
- **DG.B:** the runtime pin cannot prove that a version-1 store's contents satisfy application
  invariants or that an operator did not falsely stamp it outside this program. It proves the
  binary/store contract only when the marker was created through the designed fresh path or a
  separately ratified certification path.
- **DG.B:** rejecting a legacy store does not migrate, recover, or certify it. `4d/OD-7` owns the
  deployed-store transition, and `4d/OD-4` still owns any future migrate-versus-fail decision.
- **DG.B:** the empty-schema predicate cannot distinguish a genuinely new empty file from an old
  database from which every application object was removed. Both are structurally empty and are
  initialized as fresh; no surviving store state exists for a stronger classification.
- **DG.B:** version 1 does not detect out-of-band DDL tampering that leaves the pragma unchanged.
  DG.A's exact DDL tests constrain shipped schema changes, not arbitrary edits to a deployed file;
  runtime schema fingerprinting would be a separate kernel contract.
- **DG.B:** AC11 and AC10's secondary derivation mutation are test-side probes. AC10's primary
  re-pointed DDL mutation and AC6–AC9/AC12–AC14 are production discriminators; only the runtime
  criteria demonstrate runtime rejection.
- **DG.A carries the SAME wildcard blind spot in LANDED code, and DG.B does not close it**
  (found first-party by the controller while applying the iteration-44 carve-out; recorded rather
  than fixed, because DG.A is landed and this was a design iteration). `tableDDL` at
  `host/store/journal_test.go:877` filters `type='table' AND name NOT LIKE 'sqlite_%'` — the same
  unescaped predicate V-S proves is a single-character wildcard. A table named `sqliteX_probe`
  is therefore silently absent from the materialized DDL map, so **AC1's "unexpected tables fail
  loudly" cannot see it**. Impact is low — the canonical manifest is a hardcoded seven-table list
  and no such table exists (the search for the predicate returns exactly this one site, with a
  known-positive `sqlite_master` control firing in the same call) — but it is the same defect
  class in shipped code, and it means iteration 41's V-row calling the name filter a harmless
  "dormant guard" was understated. Carried as **CF-O-1**: the DG.B executor edits adjacent code
  and should apply the `ESCAPE` form the Design Freeze now mandates in the same pass.

**DG.B gate walk against the signature defect.** The developer edit this gate is meant to catch is
a shipping `schema.sql` change whose binary still calls that changed shape “version 1,” including
the previously sanctioned follow-up of updating DG.A's canonical manifest. That sequence remains
red against the independent frozen `schemaV1SQL`. Changing only `currentSchemaVersion` also reds
against the independent test literal and missing frozen fixture. The legitimate repair is not a
single pasted value: it is a reviewed version bump plus a new independently frozen schema fixture
and explicit acceptance policy for the now-old version. If that policy would require migration or
a support window, implementation stops for a new ratification-class design. The runtime legacy
error deliberately recommends no pragma value, so its own remedy cannot re-green by bypass. No
single action suggested by DG.B's failures both preserves the changed DDL and re-greens the suite
without establishing the new version contract; directly editing all independent expectations is
multiple review-visible policy changes, not an automatic repair supplied by one error message.

The essential promise is narrower and testable: an unreviewed journal DDL change goes red, and a
reviewed existing-table DDL change cannot become green merely by updating a literal while the
old-store path still silently fails to apply it.

---

## Quorum verification log (DG.B design quorum, iteration 44)

Two rounds plus the **narrow-refinement carve-out**. Both reviewers **present in both rounds**
(`absent_reviewers: []`), so neither round was an N−1 degrade; `--max-cost-usd` was set to `0.30`
up front. Total metered cost **$0.3192** against the `$5` ceiling.

**Round 1** — `w-ddl-gate-teeth-2026-08-03T15-21-14Z.json` ($0.1534). **BLOCKED, 3/3 reject.**
- `gemini-3-1-pro` and the **controller independently and separately** raised the SAME defect: DG.B
  makes DG.A's **AC2 vacuous** — `MUT-EXISTING-DDL-CHANGE-REMANIFESTED` alters neither
  `user_version` nor the application-object count, so after DG.B the pre-journal test returns the
  same typed legacy error with the mutation applied and unapplied. gemini: *"rendering AC2
  mathematically impossible to satisfy as written."* Two independent parties reaching one
  conclusion is why this is recorded as the round's finding rather than either party's.
- `gpt5-6-sol`: D5's case table (0 / 1 / >1) is **not total** — a negative `user_version` matches no
  branch. The controller **measured it rather than forwarding it**, and it came back **worse than
  filed**: negatives round-trip and persist, AND the field truncates silently (`2147483648` →
  reads back **0**), so an out-of-range version erases its own evidence and, on a structurally
  empty store, would be initialized as fresh. That truncation is a controller finding, not the
  reviewer's. Rows **V-N…V-R**.

**Round 2** — `w-ddl-gate-teeth-2026-08-03T15-30-28Z.json` ($0.1658). **BLOCKED 1/3** —
`gemini-3-1-pro` **PASS**, controller **PASS**, `gpt5-6-sol` reject.
- `gpt5-6-sol`: the freshness predicate is **unsound** — in `LIKE`, `_` is a single-character
  wildcard, so `NOT LIKE 'sqlite_%'` excludes `sqlite` + *any* character, not the literal reserved
  prefix. Its `catch` landed on the controller's own **V-K**, which had proved the limb *fires* but
  never tested its **exclusion boundary**.
- `gemini-3-1-pro` (non-blocking): the iteration-43 STATUS blockquote could read to an executor as
  "DG.B is still undesigned". Folded — that block now carries an explicit stale-transcript marker.

**Carve-out applied (2nd revision, controller).** The single remaining objection (a) carried a
concrete reviewer-authored `proposed_fix` and (b) did **not** dispute the design direction — a
one-predicate soundness bug plus fixtures. Per the charter's narrow-refinement carve-out, the
reviewer's fix was applied **verbatim** and **measured before adoption** (row **V-S**): SQLite
rejects `sqlite_internal_probe` (*"object name reserved for internal use"*) but accepts
`sqliteX_probe`; a store whose only application object is `sqliteX_probe` counts **0** under the
unescaped predicate (→ **fresh**, so DG.B would have stamped version 1 over real data) and **1**
under `NOT LIKE 'sqlite\_%' ESCAPE '\'` (→ **legacy**), with a two-sided control confirming real
`sqlite_` internals are still excluded and `sqliteX_probe` no longer is. Landed as the Design
Freeze prohibition, D5's predicate, D9's fixture set, and **AC14** with
`MUT-PREFIX-WILDCARD-REGRESSION`. Tally corrected twice across the cycle: **6/11 → 9/13 → 10/14**
production discriminators.

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


---

## DG.A implementation verification log (iteration 43, 2026-08-03)

**Conditions.** Worktree `.wt-iter43` branched from `origin/dev` @ `ef8e104`; toolchain `go1.26.4`
(read with `go -C <dir> env GOVERSION` — plain `go env GOVERSION` is directory-sensitive under
`GOTOOLCHAIN=auto`); `AILANG_BIN=/tmp/ailang-v0300/ailang` → `AILANG v0.30.0`; darwin/arm64.
Executor `codex:gpt-5.6-sol` (sandboxed, cannot commit); planner `opus`; controller `claude-opus-5`.

**Baseline before any edit — measured, not assumed.** `scripts/verify_go.sh` → rc=0, **10 packages
`ok`, 0 FAIL**. Counted with a TAB-safe pattern (`^ok[[:space:]]+github`) plus a known-positive
control, because `go test` prints `ok` + two spaces + a TAB and the obvious grep returns a plausible
`0` (the iteration-41/42 scar).

**Mutation ledger.** Executor-run, one run each, no retries, each restored byte-identically:

| Mutation | Class | Expected | Observed | Restored |
|---|---|---|---|---|
| `MUT-JOURNAL-DDL-WIDEN` | PRODUCTION | fresh RED on journal; historical GREEN | `canonical DDL mismatch for table "journal"`; historical PASS | sha256 ✓ |
| `MUT-EXISTING-DDL-CHANGE-REMANIFESTED` | PRODUCTION+TEST | fresh GREEN; historical RED on `store_heads` | fresh PASS; `stale historical DDL for table "store_heads"` | sha256 ✓ |
| `MUT-HISTORICAL-FIXTURE-DROP-STORE-HEADS` | TEST probe | early missing-table RED | `historical fixture missing table "store_heads"` | sha256 ✓ |
| `MUT-CANONICAL-MANIFEST-DERIVED` | TEST probe | GREEN — AC1's RED absent | both PASS; derivation kills the discriminator | sha256 ✓ |
| `MUT-UPGRADE-ASSERTION-DEAD` | TEST probe | GREEN — AC2's RED absent; vet rc=0 | `MUT-5_VET_RC=0`; both PASS | sha256 ✓ |

**Controller's independent reproduction of the load-bearing one.** An executor's gate verdict is a
claim, and AC2 is the claim the whole milestone rests on, so the controller re-ran
`MUT-EXISTING-DDL-CHANGE-REMANIFESTED` first-party rather than banking the executor's row:

- fresh-store gate `TestSchemaDDLMatchesCanonicalManifest` → **PASS**
- historical gate `TestOpenAddsJournalAndDetectsStalePreJournalDDL` → **FAIL** at `journal_test.go:870`,
  `stale historical DDL for table "store_heads"`, printing got (2 columns) vs want (3 columns)
- restored: sha256 of `schema.sql` and `journal_test.go` **identical to pre-mutation**, and
  `git status --porcelain host/store/schema.sql` empty; restored baseline `ok 0.355s`.

**Why that specific pair is the whole item.** It is the exact action the OLD gate invited — edit the
DDL, then update the expectation the failure message hands you — and under the old gate that action
turned **all 10 packages green** with the edit unapplied to every existing store (iteration 41's M4
and M5). Under the new gate the fresh-store check goes green and the historical check **reds**. That
asymmetry is the de-fang Mark ratified: the gate no longer invites its own bypass.

**Post-implementation gate, re-run OUTSIDE the executor's sandbox** (mandatory — a
`workspace-write` sandbox denies loopback binds, so `host/broker`, `host/daemon` and
`cmd/ailang-worldd` results from inside it are uninformative in both directions):
`scripts/verify_go.sh` → **rc=0, 10 packages `ok`, 0 FAIL**.

**Honest scope of what this milestone establishes.** **Two** ACs are genuine production
discriminators (AC1, AC2 — both reddened by a PRODUCTION mutation). **Three** are test-side
discrimination probes (AC3, AC4, AC5) that prove the gate's independence, not any kernel property.
Reporting "5/5 acceptance criteria with named REDs" would overstate the production evidence by 2.5×.

**What DG.A still does NOT do.** `store.Open` continues to accept a structurally stale store and
return `err=nil` — iteration 41's M5, unchanged and unfixed here. DG.A makes drift visible **in the
suite before a change lands**; it makes nothing visible to an already-deployed binary. That is
**DG.B**'s job, and DG.B is ratified but undesigned. Item 4d is **not** closed by this milestone.

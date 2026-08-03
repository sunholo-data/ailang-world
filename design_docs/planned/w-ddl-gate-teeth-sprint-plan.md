# Sprint plan — `w-ddl-gate-teeth` milestone **DG.A** (de-fang the inert DDL gate)

**Design doc:** [`design_docs/planned/w-ddl-gate-teeth.md`](w-ddl-gate-teeth.md)
**Milestone:** `DG.A` only. **`DG.B` (the `PRAGMA user_version` production contract) is OUT OF SCOPE
and must not be started** — it is ratified-but-undesigned (no ACs, no mutations, no fixtures).
**Item:** queue item 4d, clause-1
**Worktree:** `/Users/voightkampff/dev/sunholo-data/.wt-iter43`, branch `sprint/w-ddl-gate-teeth-dga`,
based on `origin/dev` @ `ef8e104`
**Planner:** `claude-opus-5` (sprint-planner), iteration 43
**Estimate:** **0.4 day** (~3.0–3.5 focused hours), risk **LOW-MEDIUM**
**Files modified by the whole sprint:** exactly one — `host/store/journal_test.go`
**LOC estimate:** ~+150 / −32 in `host/store/journal_test.go` (net ~+118), zero production LOC

---

## 0. Planner's own verification log (Gate-2: re-measured, not inherited)

Every row below was measured **by the planner, in this worktree, at `ef8e104`**, today. The design
doc's own measurements were taken at `e5027df`; the controller labelled that file unchanged since,
and rows 1–2 below independently re-confirm it. Rows marked **NEW** were *not* in the design doc and
are load-bearing for DG.A.

| Claim | How verified | Evidence | Verdict |
|---|---|---|---|
| The defect is present at HEAD as described | First-party read of `host/store/journal_test.go:631-662` | source pin at 634-636 `t.Fatalf`s before `sql.Open` at 638; `before := tableDDL(t, db)` (645) and `after := tableDDL(t, s.db)` (654) both derive from the same already-edited `schemaSQL`; `delete(after, "journal")` at 658 | CONFIRMED |
| The historical artifact has the pinned provenance | `git show '8133573:host/store/schema.sql'` piped to `shasum -a 256` | `35f09862e20ddc1c6b0467b69781b2d25fbc07d04c49f777a76a62793e14bbdd`, 76 lines, exactly 6 `CREATE TABLE`: `objects`, `worlds`, `log_entries`, `epoch_registry_heads`, `store_heads`, `verification_cache`. `d5774eb` = "durable journal substrate … (w-store-durability SD.B) (#18)" | CONFIRMED — matches the doc byte for byte |
| HEAD materializes exactly 7 durable tables; `tableDDL`'s name filter is dormant | sqlite3 CLI 3.51.0 executed HEAD `schema.sql`; `sqlite_master` enumerated | 14 rows = 7 `type='table'` + 7 `sqlite_autoindex_*` of `type='index'`. `count(type='table' AND name NOT LIKE 'sqlite_%')` = **7**, identical to `count(type='table')` | CONFIRMED — the doc's live/dormant distinction holds |
| **The 6 historical tables match HEAD after whitespace normalization** (quorum r2's load-bearing premise — if false, D2 reds on baseline) | Both schemas materialized into independent sqlite3 databases; per-table `sqlite_master.sql` compared after `tr -s ' \t\n' ' '` + trim; known-positive control in the same call (journal absent in old, present in new) | **All 6 MATCH. Zero drift.** | CONFIRMED — **D2 is green on baseline** |
| **NEW — the Go driver `modernc.org/sqlite` materializes byte-identical `sqlite_master.sql` to the sqlite3 C CLI** | Throwaway `TestZZTmpDDLProbe` added to `host/store/`, run through the real `Open`, `tableDDL` output logged with `%q`, then the file deleted and `git status --porcelain` re-checked | All 7 Go-driver strings are character-identical to the CLI's, including newlines and column alignment | **CONFIRMED — the doc never checked this and the whole manifest depends on it** |
| **NEW — SQLite strips `IF NOT EXISTS` from `sqlite_master.sql`** | Same probe | source `CREATE TABLE IF NOT EXISTS journal (` materializes as `CREATE TABLE journal (` | **CONFIRMED — the manifest must NOT contain `IF NOT EXISTS`. This is the single most likely way to burn the first mutation run.** |
| `MUT-EXISTING-DDL-CHANGE-REMANIFESTED` cannot confound other tests via `SELECT *` | `grep -rn 'SELECT \*' --include='*.go' host/` with a known-positive control (`SELECT name`) in the same call | zero `SELECT *`; all `store_heads` DML uses explicit column lists (`store.go:590,612,752`) | CONFIRMED — adding `mutant_col TEXT NOT NULL DEFAULT ''` reds nothing but the intended gate |
| Renaming the old test breaks no gate | `grep -rn 'Test[A-Z][A-Za-z]*' scripts/ .github/` with a known-positive control (`verify_go`) | **no Go test name is referenced by any script or CI manifest**; the old name survives only in `design_docs/implemented/w-effect-journal.md` and mission-log prose (historical records, not gates) | CONFIRMED — see §7 finding F1 |
| Baseline `host/store` runtime, and the ceiling is affordable | `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/store/ -count=1` | `ok … 0.356s`; `go vet ./host/store/` rc=0 | CONFIRMED — the 120 s ceiling is ~330× headroom |
| **NEW — `timeout(1)` and `gtimeout(1)` do NOT exist on this rig** | `which timeout gtimeout` | both "not found" | **CONFIRMED — the 120 s ceiling MUST be expressed as `go test -timeout=120s`, not `timeout 120 go test`, which would fail with `command not found` and look like a mutation result** |
| Toolchain and pinned binary | `go -C /Users/voightkampff/dev/sunholo-data/.wt-iter43 env GOVERSION`; `/tmp/ailang-v0300/ailang --version` | `go1.26.4`; `AILANG v0.30.0` (commit `e37b370`). `go.mod` floor is `go 1.26.4` | CONFIRMED |
| Import fates in `journal_test.go` | `grep -n 'sha256\.\|hex\.\|strings\.\|reflect\.'` | `crypto/sha256` and `encoding/hex` are used **only** at lines 634-635 → both must be deleted. `strings` is used **only** at 633 → but the normalizer will reintroduce a use (`strings.Fields`/`strings.Join`), so it stays. `reflect` is used at 149/158/272/626 → stays | CONFIRMED |

**Pre-existing dirt in the worktree:** `git status --porcelain` shows ` M design_docs/planned/w-ddl-gate-teeth.md` — the controller's uncommitted iteration-43 Status-block transcription (72+/14−). That is **not** the executor's; the executor must not revert it and must not count it as a scope violation. The executor's own diff must add exactly `host/store/journal_test.go` and this plan file to that set.

---

## 1. Scope boundary (the Design Freeze, restated as an executable checklist)

**IN SCOPE — DG.A modifies exactly one file: `host/store/journal_test.go`.**

**OUT OF SCOPE — an executor that touches any of these has failed the sprint:**

- `host/store/schema.sql` — *except* transiently, as `MUT-JOURNAL-DDL-WIDEN` /
  `MUT-EXISTING-DDL-CHANGE-REMANIFESTED`, each restored **byte-identically with sha256 printed on
  both sides** in the same session.
- `host/store/store.go`, `store.Open`, `openSQLite`, accepted on-disk versions, `PRAGMA user_version`.
- Any `.ail` file, `design_docs/sketches/`, `scripts/verify_ail.sh`, its required-check manifest,
  `scripts/verify_go.sh`, `.github/workflows/ci.yml`.
- Any migration registry, migration CLI, or schema rewrite.
- `design_docs/implemented/w-effect-journal.md` (see finding F1 — leave it alone).

**DG.B is not merely deprioritized, it is undesigned.** If the executor finds itself writing
`PRAGMA user_version`, it has left DG.A. Stop and report.

---

## 2. Day-by-day breakdown (~0.4 day / 3.0–3.5 h)

### Half-day 1 — implement (≈2.0 h)

| # | Task | Est | Output |
|---|---|---|---|
| **T1** | Copy the historical fixture in verbatim: `git show '8133573:host/store/schema.sql'` → a backtick `const preJournalSchemaV0` in `journal_test.go`. **Verify the constant round-trips**: the executor prints the sha256 of the file it copied from and eyeballs that no line was dropped. Fixture comment cites `8133573`, notes `d5774eb` added journal, records sha256 `35f09862e2…bbdd`, **and states in words that the fixture must never be updated to follow current DDL.** | 20 m | ~80 LOC |
| **T2** | Add `canonicalTableDDL` — a hardcoded `map[string]string` with all **seven** materialized statements (Appendix A gives them verbatim; **no `IF NOT EXISTS`**). Add a doc comment stating it must not be derived from `schemaSQL`, a hash of it, or the database under test. | 25 m | ~50 LOC |
| **T3** | Add `normalizeDDL(string) string` = `strings.Join(strings.Fields(s), " ")`. Whitespace only: no lowercasing, no reordering, no quote-stripping, no parsing. One-line body, comment stating what it deliberately does **not** do. | 5 m | ~5 LOC |
| **T4** | Add `TestSchemaDDLMatchesCanonicalManifest` (D1 / **AC1**): fresh temp-file `Open`; `tableDDL`; **sorted** name-set equality with explicit loud failures for zero / missing / unexpected tables; then per-table normalized comparison in sorted order, failure message naming the offending table. | 30 m | ~40 LOC |
| **T5** | Replace `TestPreJournalMigrationPreservesExistingDDL` with `TestOpenAddsJournalAndDetectsStalePreJournalDDL` (D2 / **AC2**, **AC3**): exec `preJournalSchemaV0`; insert sentinel `INSERT INTO store_heads (head_key, world_ref) VALUES ('selected_world_head', <ref>)`; close setup conn; real `Open(path)`; assert sentinel unchanged; assert `journal` present; assert the historical name set is **exactly** the six names; assert each of the six is present in `canonicalTableDDL`; compare each of the six, sorted, normalized. **Range over the six names — do NOT build a full map and `delete` journal.** | 35 m | ~55 LOC |
| **T6** | **AC4 deletions**: remove the `schemaSQL[:strings.Index(...)]` slice, the sha256 literal pin, `delete(after, "journal")`, the same-source `reflect.DeepEqual(after, before)`, and the imports `crypto/sha256` + `encoding/hex`. Keep `strings` (normalizer) and `reflect` (used by 4 other tests). Keep `tableDDL` as the single SQLite extraction helper. | 10 m | −32 LOC |
| **T7** | Honest claim in code: a comment beside the renamed test saying *"Open creates the absent journal table; it does NOT upgrade existing tables. A green result here is not an upgrade."* | 5 m | prose |
| **T8** | `gofmt -l host/store/journal_test.go` (expect empty) → `go vet ./host/store/` → `AILANG_BIN=… go test ./host/store/ -count=1 -timeout=120s`. Iterate here until green — **these are implementation runs, not mutation runs, and the no-retry-loop rule does not bind them** (see finding F3). | 10 m | GREEN |

**Exit criterion for half-day 1:** `host/store` green, `git diff --name-only` lists exactly
`host/store/journal_test.go` (plus the controller's pre-existing doc dirt and this plan).

### Half-day 2 — mutation testing + record (≈1.0–1.5 h)

Execute §3 in order. Budget ~10 min per mutation including restoration and sha256 proof.

| # | Task | Est |
|---|---|---|
| **T9** | Run MUT-1 … MUT-5 per §3, one run each, no retries, restoring byte-identically after each | 55 m |
| **T10** | Final restored-baseline run: `./scripts/verify_go.sh` with `AILANG_BIN=/tmp/ailang-v0300/ailang` — the **full repository gate**, expecting rc=0 and **10 packages `ok`**. Read the count with a TAB-safe pattern (see §5 rig note) | 10 m |
| **T11** | Record the Verification Log in the design doc's tail: toolchain `go1.26.4`, `AILANG_BIN=/tmp/ailang-v0300/ailang` (`AILANG v0.30.0`), sqlite3 CLI 3.51.0, darwin/arm64, one row per named mutation with its expected/observed RED and its restoration sha256 pair | 20 m |

---

## 3. Mutation sequence — the five named mutations

**Global rules (from the Design Freeze, binding):**

- **One run per named mutation. Plus one restored baseline run at the end. No retry loops.** If a
  mutation produces an unexpected result, that is a **finding to report**, not a prompt to re-run.
- **Ceiling: `-timeout=120s` on every focused/package invocation.** `timeout(1)` does not exist on
  this rig — use `go test -timeout=120s`, never `timeout 120 go test`.
- **Every mutation prints `shasum -a 256` of every file it edits, before the edit and after the
  restore, and the two must be equal.** Take the backup with `cp` **in the same command that applies
  the mutation**; restore with `cp`; end with `git status --porcelain` showing no unexpected entry.
- **Classification is stated for every mutation and never blurred.** A TEST-side mutation proves the
  gate mechanism discriminates; it never proves a kernel property.

### Sequence table

| # | Mutation | File(s) edited | Class | Expected result — the whole result, both limbs | Restoration |
|---|---|---|---|---|---|
| **MUT-1** | `MUT-JOURNAL-DDL-WIDEN` | `host/store/schema.sql` — add `MUTANT` to journal's `kind` CHECK: `CHECK (kind IN ('intent','outcome','MUTANT'))` | **PRODUCTION** | **RED**: `TestSchemaDDLMatchesCanonicalManifest` fails with a **journal-specific** mismatch naming `journal`. **GREEN**: `TestOpenAddsJournalAndDetectsStalePreJournalDDL` stays green (it ranges over the six historical names; journal is outside its comparison). Both limbs required — a red in *both* tests means D2's scope leaked into journal and the separation D2 §"the comparison intentionally ranges over the six historical names" was not implemented. | `cp` backup back; sha256 must equal the pre-edit value |
| **MUT-2** | `MUT-EXISTING-DDL-CHANGE-REMANIFESTED` | **two files**: `host/store/schema.sql` (add `mutant_col TEXT NOT NULL DEFAULT ''` to `store_heads`) **AND** `host/store/journal_test.go` (update **only** the `store_heads` entry of `canonicalTableDDL` to the new materialized form) | **PRODUCTION + TEST** — this is the developer's *legitimate* re-manifest ritual, reproduced exactly | **GREEN**: `TestSchemaDDLMatchesCanonicalManifest` — the fresh store really does have `mutant_col`, so the updated manifest matches. **RED**: `TestOpenAddsJournalAndDetectsStalePreJournalDDL` alone, with a `store_heads` mismatch, because `CREATE TABLE IF NOT EXISTS` did not apply the column to the pre-existing table. **This asymmetric pair IS the sprint's central claim.** If both go green, M4's one-line re-green path is still open and the sprint has recreated the bug it is fixing. | `cp` **both** backups back; **two** sha256 pairs printed |
| **MUT-3** | `MUT-HISTORICAL-FIXTURE-DROP-STORE-HEADS` | `host/store/journal_test.go` — delete the `store_heads` `CREATE TABLE` from `preJournalSchemaV0` | **TEST** (fixture-drift probe) | **RED**: the exact six-name fixture assertion fires **before** the DDL comparison — the failure message must name a *missing historical table*, not a DDL mismatch and not a sentinel-insert SQL error. If the sentinel `INSERT` panics first with `no such table: store_heads`, the six-name assertion is sequenced too late; that is a defect to fix (assert the name set before inserting the sentinel, or accept the earlier failure only if it is loud and named). | `cp` back; sha256 pair |
| **MUT-4** | `MUT-CANONICAL-MANIFEST-DERIVED` | **two files, compound**: `host/store/journal_test.go` (replace journal's hardcoded expected DDL with the actual journal DDL just read from the database) **AND** re-apply `MUT-JOURNAL-DDL-WIDEN` to `host/store/schema.sql` | **TEST** (discrimination probe) + PRODUCTION carrier | **GREEN — and the green IS the finding.** MUT-1's required RED must be **absent**. That absence is the proof that a derived manifest kills the gate, i.e. that the *hardcoded* manifest is load-bearing. Record it as "AC1's required RED absent under derivation ⇒ mutation review would reject a derived implementation." **Do not report this green as a passing test.** | `cp` **both** backups back; **two** sha256 pairs |
| **MUT-5** | `MUT-UPGRADE-ASSERTION-DEAD` | **two files, compound**: `host/store/journal_test.go` (short-circuit **only** the six-table post-open DDL comparison — keep all variables used so `go vet` stays rc=0) **AND** re-apply MUT-2's `schema.sql` + `store_heads` manifest edits | **TEST** (discrimination probe) + PRODUCTION carrier | **GREEN — again the green is the finding.** MUT-2's D2 RED must **disappear** while `TestSchemaDDLMatchesCanonicalManifest` stays green. Run `go vet ./host/store/` first and record rc=0, so the green cannot be dismissed as an unused-variable compile artifact. | `cp` **both** backups back; **two** sha256 pairs |
| **BASE** | *(restored baseline)* | none | — | **rc=0, 10 packages `ok`, 0 FAIL** from `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh`, and `git status --porcelain` carrying only the intended edits | — |

### Why this order

MUT-1 and MUT-2 are the **PRODUCTION** discriminators and run first, on a clean implementation, so
their results are attributable to the gate and nothing else. MUT-3 is a cheap single-file TEST probe.
MUT-4 and MUT-5 are **compound** probes that re-apply an earlier production mutation as a carrier —
they run last because they are the only ones that leave two files dirty at once, and because their
expected result is the *absence* of a previously observed RED, which is only meaningful once that RED
has actually been observed (MUT-1 for MUT-4; MUT-2 for MUT-5).

### Command shape (per mutation, one invocation)

```bash
cd /Users/voightkampff/dev/sunholo-data/.wt-iter43
F=host/store/schema.sql
cp "$F" /tmp/dga_bak_schema.sql && shasum -a 256 "$F"     # pre-edit sha256
# ... apply the single named edit ...
shasum -a 256 "$F"                                         # mutated sha256 (must DIFFER)
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/store/ -count=1 -timeout=120s -run 'TestSchemaDDLMatchesCanonicalManifest|TestOpenAddsJournalAndDetectsStalePreJournalDDL' -v
cp /tmp/dga_bak_schema.sql "$F" && shasum -a 256 "$F"      # post-restore sha256 (must EQUAL pre-edit)
git status --porcelain
```

**zsh hazards that have bitten this mission — the executor must obey all three:**

1. Quote glob-shaped flag values: `--include='*.go'`. Unquoted, zsh aborts before the command runs.
2. Brace any variable followed by a colon: `git show "${rev}:host/store/schema.sql"`. Unbraced, zsh
   applies `:h`/`:t` as history modifiers and **silently rewrites the path**, yielding a plausible
   zero rather than an error.
3. Any "no matches" / empty result is a **CLAIM, not a fact** — every negative check must carry a
   known-positive control **in the same call**.

---

## 4. Discriminating vs. probe — which ACs actually prove something

This distinction is explicit in the design doc and the evaluator needs it stated. Do not let a probe
be reported as a proof.

| AC | Mutation | Class | What it genuinely establishes |
|---|---|---|---|
| **AC1** | `MUT-JOURNAL-DDL-WIDEN` | **PRODUCTION — genuinely discriminating** | A real, realistic edit to shipped `schema.sql` reds the suite. This is the first time any repository test owns journal's exact DDL. Strongest single result in the sprint. |
| **AC2** | `MUT-EXISTING-DDL-CHANGE-REMANIFESTED` | **PRODUCTION — genuinely discriminating, and the controlling evidence** | The *sanctioned repair ritual itself* (edit DDL, update the literal) no longer buys a green. This closes M4/M5's one-line re-green path. **This is the only AC that proves the fix is not the old bug wearing a new manifest.** |
| **AC3** | `MUT-HISTORICAL-FIXTURE-DROP-STORE-HEADS` | **TEST — probe** | Only that the fixture cannot silently erode. It says nothing about `store.Open`. |
| **AC4** | `MUT-CANONICAL-MANIFEST-DERIVED` | **TEST — probe** | Only that the gate's independence from the subject is load-bearing. Proves nothing about journal semantics. |
| **AC5** | `MUT-UPGRADE-ASSERTION-DEAD` | **TEST — probe** | Only that the six-table comparison is the mechanism actually doing the work — i.e. it forecloses iteration 38's defect (*teeth from one mechanism while the documented mechanism is dead*). |

**Two discriminating ACs, three probes.** The sprint report must say exactly that. Reporting "5/5 ACs
with named RED mutations" without the split would overstate the result by a factor of 2.5.

---

## 5. The signature failure mode — planned against explicitly

> This mission's recurring defect is **a gate that cannot fail**, and twice the bug was reintroduced
> *inside the remedy for it*.

The specific trap here: **the replacement manifest is hardcoded, so a developer can update it.** The
design doc concedes this openly (D1: *"a reviewed change detector"*). AC2 exists to prove that
updating it **still** leaves the historical-store test red.

**Binding rule for this sprint — AC2 may only be satisfied by the full asymmetric pair:**

> Under `MUT-EXISTING-DDL-CHANGE-REMANIFESTED`, `TestSchemaDDLMatchesCanonicalManifest` is **GREEN**
> and `TestOpenAddsJournalAndDetectsStalePreJournalDDL` is **RED**, in the same run, with the red
> naming `store_heads`.

Anything weaker recreates the bug:

- ❌ "AC2 passes because a test failed under the mutation" — **insufficient**. If the fresh-store test
  is the one that reds, the manifest edit was incomplete and nothing about the historical path was
  proven.
- ❌ "AC2 passes because the historical test reds under a `schema.sql` edit alone" — **insufficient**.
  That is M2, the *old* behaviour; the whole point is that the developer **also updates the literal**.
  The re-manifest half is not optional.
- ❌ "AC2 passes on the strength of MUT-5's probe" — **insufficient and inverted**. MUT-5's expected
  result is a *green*; a probe whose success is an absence can never carry a discrimination claim.
- ❌ Deriving the `store_heads` manifest entry for MUT-2 by printing it from the mutated database and
  pasting it in — **this is `MUT-CANONICAL-MANIFEST-DERIVED` done by hand.** The entry must be written
  by adding `mutant_col TEXT NOT NULL DEFAULT ''` to the *existing literal text*, by hand.

**A related trap on the implementation side.** Generating Appendix A's manifest by running a probe and
pasting the output is legitimate **once, at authoring time, reviewed in the diff** — Appendix A below
was produced exactly that way and the probe file was deleted. It becomes the defect the moment any of
it survives as a **runtime** code path that reads DDL into the expected map. The final file must
contain **no** assignment flowing from `tableDDL` / `sqlite_master` / `schemaSQL` into an *expected*
value. MUT-4 is the check; the executor should also grep the final diff for it and report the grep
with a known-positive control.

---

## 6. Honesty limits carried into the plan (from "What these gates CANNOT fail")

These are part of the deliverable, not fine print. The sprint report and the code comments must both
carry them; a sprint that claims more than this list allows has failed even if all five mutations
behave.

DG.A **cannot**:

1. Stop a developer intentionally editing both production DDL and the manifest. D1 is a
   review-visible change detector, **not an authorization system**.
2. Prove a new constraint or column has correct application *semantics*.
3. Prove data preservation beyond the single sentinel row.
4. **Detect drift in an already-deployed store at runtime.** `store.Open` still returns success for a
   structurally stale store. **That is DG.B's job and it is still open** — the sprint must not claim
   item 4d is "done", only that DG.A's half is.
5. Decide whether a mismatch should fail, migrate, or use an external tool (OD-3 / OD-4).
6. Verify indexes, triggers, views, or PRAGMAs — `tableDDL` selects `type='table'` only.
7. Prove compatibility with every historical database; the fixture covers **one** pre-journal shape.
8. Distinguish DDL differing only in whitespace the normalizer collapses; nor avoid redding on a
   semantically equivalent change whose materialized text differs.
9. Make `CREATE TABLE IF NOT EXISTS` apply an edit to an existing table.
10. Validate the Go toolchain — results are only meaningful on the controller's known-good rig.
11. **AC3/AC4/AC5's test-side mutations prove no production kernel property whatsoever.**

The one promise the sprint may make, in these words: *an unreviewed journal DDL change goes red, and a
reviewed existing-table DDL change cannot become green merely by updating a literal while the old-store
path still silently fails to apply it.*

### Rig notes that have cost this mission time before

- `go env GOVERSION` is **directory-sensitive** under `GOTOOLCHAIN=auto`. Always `go -C <dir> env GOVERSION`.
- `go test` prints `ok` + **two spaces + TAB** + import path. The obvious `grep -cE 'ok  +github'`
  returns a **plausible 0**. Use `grep -c $'^ok\t'` or `grep -c 'ok.*github.com/sunholo-data'`, and
  print a known-positive control alongside.
- `AILANG_BIN` must be exported for `verify_go.sh`; unset, the host/replay tests `t.Skip()` silently.
  The `ailang` on PATH is a **v0.31.0-dirty dev build — never use it**.

---

## 7. Where the design doc is underspecified or wrong

Five findings. None blocks DG.A; F2 is a correctness trap that would have burned a mutation run.

- **F1 — dangling references to the renamed test, and no instruction about them.** AC4 requires
  retiring the old test name; `design_docs/implemented/w-effect-journal.md` names
  `TestPreJournalMigrationPreservesExistingDDL` at lines 92, 153, 532, 572, 579, 708, and the charter
  and mission log name it too. The doc never says what to do. **Verified: no script or CI manifest
  references any Go test name**, so nothing breaks. **Plan's ruling: leave every one of them
  untouched.** They are historical records of what was true when written, and editing them would
  violate the "modifies exactly `host/store/journal_test.go`" freeze. Flagged for the controller: a
  later bookkeeping commit may want a forward-pointer, but it is **not** DG.A's.

- **F2 — the doc never says the manifest must omit `IF NOT EXISTS`, and this is the most likely way to
  waste a mutation run.** D1 says the values are *"the `sqlite_master.sql` forms SQLite materializes"*,
  which is correct but abstract. Measured: SQLite **strips `IF NOT EXISTS`**, so copying `schema.sql`'s
  statements verbatim produces a manifest that reds on baseline for **all seven** tables. Appendix A
  below removes the ambiguity by giving the exact strings.

- **F3 — "one run per named mutation plus one restored baseline run" is ambiguous about the
  implementation phase.** Read literally it would forbid the ordinary write-test-run-fix loop of T1–T8.
  **Plan's reading, stated so the evaluator can reject it if wrong: the no-retry-loop rule binds
  MUTATION runs only.** Implementation runs are unbounded in count (each still capped at 120 s), and
  restoration is verified by **sha256 equality, not by a test run**, so exactly one restored-baseline
  run occurs, at the very end.

- **F4 — the doc never verified that the Go SQLite driver materializes the same text as the sqlite3
  CLI, though every one of its measurements used the CLI and every one of DG.A's assertions will use
  the driver.** This is a real gap of the same species as quorum r2's catch: if `modernc.org/sqlite`
  had differed by so much as one space, the whole manifest would have redded on baseline. **Measured
  by the planner: identical, all seven tables.** The premise holds; it just had not been checked.

- **F5 — AC3's assertion ordering is unspecified and may be unreachable as written.** AC3 requires
  *"The exact six-name fixture assertion must red **before** the DDL comparison"*, but D2's step list
  puts the sentinel `INSERT` at step 2 and the name-set assertion at step 7. Under
  `MUT-HISTORICAL-FIXTURE-DROP-STORE-HEADS` the `INSERT` will fail first with
  `no such table: store_heads`. **Plan's ruling: assert the historical name set immediately after
  executing the fixture and before inserting the sentinel.** This satisfies AC3's ordering literally
  and yields a failure message that names the missing table rather than an opaque SQL error. Recorded
  here rather than silently coded, because it is a deviation from D2's numbered sequence.

---

## Appendix A — the seven canonical materialized DDL statements

**Provenance:** produced once, at plan time, by executing HEAD `schema.sql` through the real
`store.Open` and logging `tableDDL` with `%q`; independently cross-checked against sqlite3 CLI 3.51.0
materializing the same file. The probe file was deleted and `git status --porcelain` re-verified.
These are **literals to be typed into the test**, never a runtime derivation.

Each is the corresponding `schema.sql` statement **with `IF NOT EXISTS ` removed and everything else
byte-identical** (leading comments are not part of `sqlite_master.sql`). Write them as backtick raw
strings so the normalizer is genuinely exercised on both sides.

```
objects              CREATE TABLE objects (\n    hash_ref           TEXT PRIMARY KEY,\n    interface_hash_ref TEXT NOT NULL,\n    semantic_id        TEXT NOT NULL,\n    provenance         TEXT NOT NULL,\n    payload            BLOB NOT NULL\n)

worlds               CREATE TABLE worlds (\n    world_ref  TEXT PRIMARY KEY,\n    revision   INTEGER NOT NULL,\n    state_root TEXT NOT NULL,\n    log_head   TEXT NOT NULL\n)

log_entries          CREATE TABLE log_entries (\n    entry_index         INTEGER PRIMARY KEY,\n    entry_hash_ref      TEXT NOT NULL UNIQUE,\n    semantics_epoch     INTEGER NOT NULL,\n    transition_fn_ref   TEXT NOT NULL,\n    interpreter_ref     TEXT NOT NULL,\n    prev_entry_hash_ref TEXT NOT NULL,\n    written_by          TEXT NOT NULL,\n    transition_ref      TEXT NOT NULL\n)

epoch_registry_heads CREATE TABLE epoch_registry_heads (\n    registry_name TEXT PRIMARY KEY,\n    object_ref    TEXT NOT NULL\n)

store_heads          CREATE TABLE store_heads (\n    head_key  TEXT PRIMARY KEY,\n    world_ref TEXT NOT NULL\n)

verification_cache   CREATE TABLE verification_cache (\n    transition_fn_ref TEXT NOT NULL,\n    interpreter_ref   TEXT NOT NULL,\n    semantics_epoch   INTEGER NOT NULL,\n    verified          INTEGER NOT NULL,\n    result_detail     TEXT NOT NULL,\n    PRIMARY KEY (transition_fn_ref, interpreter_ref)\n)

journal              CREATE TABLE journal (\n    seq           INTEGER PRIMARY KEY,\n    kind          TEXT NOT NULL CHECK (kind IN ('intent','outcome')),\n    invocation_id TEXT NOT NULL CHECK (invocation_id <> ''),\n    object_ref    TEXT NOT NULL CHECK (object_ref <> ''),\n    UNIQUE (invocation_id, kind)\n)
```

No statement contains a backtick, so raw string literals are safe. `journal` contains single quotes
(`'intent','outcome'`) — also safe in backticks.

**MUT-2's re-manifested `store_heads` entry**, for reference (typed by hand from the literal above,
**not** pasted from a mutated database):

```
CREATE TABLE store_heads (\n    head_key   TEXT PRIMARY KEY,\n    world_ref  TEXT NOT NULL,\n    mutant_col TEXT NOT NULL DEFAULT ''\n)
```

Exact interior whitespace is irrelevant under the normalizer; the **tokens** are what matter.

---

## Appendix B — success metrics

- [ ] `git diff --name-only` on the branch lists exactly `host/store/journal_test.go` and this plan
      (plus the controller's pre-existing `w-ddl-gate-teeth.md` dirt).
- [ ] `gofmt -l host/store/journal_test.go` → empty; `go vet ./host/store/` → rc=0.
- [ ] `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` → rc=0, **10 packages `ok`, 0 FAIL**,
      counted with a TAB-safe pattern plus a known-positive control.
- [ ] Five named mutations executed, **one run each**, each with a printed sha256 pair proving
      byte-identical restoration, and `git status --porcelain` clean of mutation residue.
- [ ] AC2's asymmetric pair observed **in one run**: fresh-store GREEN, historical-store RED naming
      `store_heads`.
- [ ] Final diff contains no runtime derivation of an expected DDL value (grep + known-positive control).
- [ ] Verification Log written into `design_docs/planned/w-ddl-gate-teeth.md`'s tail with toolchain,
      `AILANG_BIN`, sqlite3 version and platform recorded.
- [ ] Report states **two discriminating ACs, three probes** — and states that **DG.B remains open**,
      so item 4d is **not** closed.

**Explicit non-goals:** no `.ail` file, no sketch, no `user_version`, no migration machinery, no
`store.Open` change, no edit to `design_docs/implemented/w-effect-journal.md`.

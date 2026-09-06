# w-go-test-compile-fence — `go build ./...` is not a compile fence for a `_test.go` mutant

> **Controller disposition, iteration161:** Do not implement this draft. Both permitted
> quorum rounds are BLOCKED. Round2 has Astra/Gemini reject and GLM absent (invalid JSON).
> The frozen regex misses the ordinary Markdown sole-fence assertion it must detect;
> a plain-text positive control fires. Record boundaries and same-mutant companion
> association remain unspecified. R8's claimed12 matching files is refuted: the exact
> command returns13. File-level companion presence does not establish every drill's
> compliance; the all-corpus-clean conclusion below is UNVERIFIED. The draft's direct
> `timeout` commands are prescriptions, not commands the controller ran: the controller
> used Python subprocess deadlines of120s and restored captured bytes. No lint exists,
> no verification script changed, and no acceptance criterion is closed by this record.
> The rejected revision is retained below for review, not as an operative plan.
> [Evidence](../verification/world-iter161/README.md). D-WORLD-33: REDESIGN explicit
> mutation/fence evidence records through fresh gates (recommended), or DEFER this row;
> unattended default DEFER immediately. No narrow-refinement override was applied.

**Status**: PARKED — needs-human-review (D-WORLD-33)
**Date**: 2026-09-06 (revision 1 — the ONE protocol-mandated revision reconciling the blocked
3-reviewer quorum; see §Quorum Reconciliation)
**Queue item**: 65, `w-go-build-is-not-a-compile-fence-for-a-test-file` (clause-2, AC-instrument
class; planner-found iter-149, controller-reproduced first-party, declared out of scope by the row
it was discovered in, pre-registered as **Instance 1**)
**Priority**: P1
**Estimated**: ≤0.2d (audit done; canon + reproduction + an **executable** instrument-health guard +
optional clarity lines)
**Dependencies**: none
**Planner-Lane**: codex-ok (docs + **one narrowly-scoped non-fleet verification edit** — the guard's
executable and a single additive leg in `scripts/verify_ail.sh`; nothing in `world/ host/ kernel
tools/launchd`; no AILANG code)

---

## What this is

A design for **closing the AC-instrument gap** where a mutation drill asserts "the mutant builds"
using `go build ./...` while the mutant actually lives in a `_test.go` file. `go build ./...` does
not compile `_test.go`, so that assertion cannot tell "the mutant is sound" from "the mutant does
not compile" — the exact third-fact-wearing-the-same-exit-code class the mutation discipline exists
to close (coding-standards S6).

**Audit headline (the honest finding, measured — see Verification Log):** the *flawed sole-fence
pattern is effectively **not present** in the current not-yet-shipped planned corpus.* Every planned
sprint plan that lands mutants in a `_test.go` has already adopted a test-compiling companion fence
(`go vet ./host/…` and/or `go test -run '^$'`) as its load-bearing typecheck, because row 55
(iter-149) propagated the lesson. The pattern that *does* survive is (a) historical and append-only
(implemented docs, evidence/mutation logs, world-mission STATUS records — these **must not be
rewritten**, per this mission's standing rule), and (b) latent as an undocumented risk with **no
durable guard** pinning it. **Row 65's premise of a broad live rewrite is therefore not what the
audit found; the deliverable is a durable, *executable* guard + the canon, not a mass sweep and not a
grep written only in a markdown table.** This is the quorum-driven correction of the prior draft.

---

## Quorum Reconciliation (revision 1 — the one mandated revision)

The prior draft's guard was a **documented instrument-health grep living only in this doc**. All
three reviewers rejected it on the same axis: a triggerless grep that no verify leg or CI step
invokes is not a guard (it silently degrades to "no guard" the moment a future executor skips reading
this doc — the "no silent fallbacks" axiom violated by the guard's own absence), and R15/R16 were not
same-scope controls. This revision resolves it by **freezing the guard as a standalone executable
non-fleet script wired into the existing World verify leg**, and by re-labelling all guard-execution
evidence **PENDING** (it is designed, not yet executed). Reviewer-specific closures:

- **gemini-3-1-pro** — Unresolved Choice 1 resolved: a standalone executable lint
  `scripts/lint_planned_fences.sh` mechanically enforces the negative-existence check, wired into the
  existing verification surface (reuse of `scripts/verify_ail.sh` as host, §Conflict-surface reuse)
  instead of inert documentation.
- **oc-glm-5-2** — the guard is wired into `scripts/verify_ail.sh` as a bounded leg AND a **Guard
  Lifecycle** section names the exact trigger (primary: verify-leg; interim: mandated manual trigger
  with bounded deadline). R13/R14 now quote the actual `git grep -n` lines, not paraphrased line
  numbers. AC3 is explicitly instrument-health with a positive control, never load-bearing.
- **gpt6-astra** — candidate enumeration, semantic classification and enforcement are separated; M2
  is labelled **unexecuted/PENDING**; the false `-count=1 forces recompile` claim is corrected; R7 is
  narrowed to "no **local** init/TestMain" plus an explicit dependency-initialization assessment; every
  drill carries a bounded wall-clock deadline with timeout treated as failure; the guard is specified
  as failing closed (unclassified candidates fail review), not as natural-language semantic proof.
- **controller (in-session)** — axiom arithmetic corrected (net **+6**, the +1 rows total 6, not 7);
  row-55 iteration attribution corrected (**149**, not 160); pending evidence is labelled **pending,
  never confirmed**; controls stay in one enumerated corpus; the caching statement is corrected; the
  guard is frozen as an invoked durable check, not a docs-only marker.

---

## Problem Statement

**Current state.** The repo's own mutation discipline forbids vacuous gates (S6) and requires
load-bearing criteria to be discharged by a **mutation** that makes the test red. But the fence used
to prove a mutant "builds" has been, and could again be, `go build ./...` — which compiles only
non-test packages. Because **this mission's gates are tests, and all of `host/verifygate/` (8 files)
is test-only**, a `go build ./...` fence cannot see a single mutant this loop produces.

Measured first-party (controller-reproduced this iteration on HEAD `09726ae`, go1.26.6, darwin/arm64;
byte-identical restore; R5):

| Fence | Pristine | After appending a test-file type error |
|---|---|---|
| `go build ./...` | rc=0 | **rc=0** ← the broken fence is blind |
| `go test -count=1 -run '^$' ./host/verifygate/` | rc=0 | **rc=1** "cannot use ... as int" |
| `go vet ./host/verifygate/` | rc=0 | **rc=1** same error |

Two working fences exist; the one the flawed prose leans on is the broken one.

**Impact.** Any design doc or sprint plan that a future executor copies for a `_test.go` mutation
drill and that still says "assert the mutant builds with `go build ./...`" silently measures the
wrong thing: a non-compiling mutant reports a green compile-fence, so the drill's "sound" verdict is
meaningless for exactly the mutants that matter (test-side gates). The cost is not a CI hole (CI is
covered; R6) but an **AC-instrument** that cannot do what it claims.

**Why not just "use `go test` everywhere".** `go test -run '^$'` compiles *and links* the test binary
and, once built, executes `init()` and any `TestMain` — and packages it **imports** may run their own
`init()`. It is **not** a universally side-effect-free compile command (see §3.2). The corrected fence
must be chosen per-package. This repo's gates (`host/verifygate/`) have no local `init`/`TestMain`
(measured empty), but side-effect freedom also depends on what they import (§R7). A blanket
replacement that ignores this is a new, subtler bug. The canonical default is therefore **`go vet`**
(never executes code); `go test -run '^$'` is the linking/initializing alternative used only where the
package is assessed safe (R7).

---

## Goals

**Primary Goal:** make "the mutant compiles" an honest, package-scoped claim in any future mutation
drill, and pin a **durable, invoked** guard so the `go build` sole-fence cannot silently recur.

**Success Metrics**
- Corrected-fence canon documented (what each command proves, the init/TestMain **and dependency-init**
  caveat, and the caching truth). — §3.2
- A load-bearing reproduction (mutation) proving the corrected fence reds where `go build` stays
  green on a `_test.go` type-error, restore byte-identical. — §7 **M1 (confirmed, controller-reproduced)**
- An S6-labelled **instrument-health** guard over the operative (planned) corpus, shipped as an
  **executable** `scripts/lint_planned_fences.sh`, wired into the standard verify leg, with a
  same-scope positive control so an empty scan is non-vacuous. — §7 **M2 (designed; execution PENDING)**
- The audited corpus classification on record with real `git grep` lines quoted (no unverified line
  citations). — §5 / Verification Log
- verify legs green outside any drill. — §8 AC5 (pending the guard wiring)

---

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|---|---|---|---|---|
| Corrected fence = `go vet ./host/<pkg>/` (side-effect-free pure typecheck — **canonical default**) *or* `go test -count=1 -run '^$' ./host/<pkg>/` (compile+link; **may run init/TestMain and imported init**; chosen only where the package is assessed safe) | Determines what every future test-mutant "builds" proof means; wrong pick = side-effectful drills or a fence that still can't see test files | agent (with the §3.2 matrix + R7 dep-init assessment) | design | low (docs/guard) |
| `go build ./...` retained ONLY as a non-test-tree/whole-repo build fence | Keeps valid gates valid (verify_go.sh, CLAUDE.md, ci.yml, README) and prevents over-rotation | human/agent (mission-confirmed: ordinary gate is valid) | design | high if reversed |
| **Guard = a standalone executable non-fleet lint `scripts/lint_planned_fences.sh`, wired as one additive leg in `scripts/verify_ail.sh`** (reuse of the existing World verify checkpoint; not a new CI job, not a triggerless grep) | This is the quorum's central objection — a docs-only marker cannot be a durable guard if nothing invokes it | agent, **frozen by this revision** (resolves Unresolved Choice 1) | design | low (new script + one leg) |
| Implemented/evidence/mutation-log history is append-only and NOT rewritten | Preserves measured results (mission standing rule); prevents falsifying the historical record | human | design | high |

### Design Freeze
- [x] Corrected fence canon: `go vet ./host/<pkg>/` and `go test -count=1 -run '^$' ./host/<pkg>/` both
      compile test files; `go build ./...` never does.
- [x] `-count=1` **disables test-result caching only** — it does **not** force recompilation; the build
      cache is retained. (This corrects the prior draft.)
- [x] Guard mechanism **frozen**: standalone non-fleet executable `scripts/lint_planned_fences.sh`
      (a bounded candidate-enumerator + fail-closed classifier over the planned corpus), invoked as a
      single additive leg of `scripts/verify_ail.sh`. Proposed **non-fleet verification edit**; the
      actual edit is implemented in the sprint execution phase, **not by this designer session**.
- [x] The only gate edits are those two narrowly-scoped non-fleet additions. `verify_go.sh`, `ci.yml`,
      `README.md`, `CLAUDE.md`, `tools/launchd/`, coding-standards, shared skills, and every `.ail` are
      **untouched**.
- [x] Guard-execution evidence (**M2 / AC3 / AC5**) is **PENDING** until the lint exists and its
      fixtures run; it is not claimed confirmed. M1 (the fence reproduction) is **confirmed**
      (controller-reproduced).

### Deferred Decisions
- `go vet` vs `go test -run '^$'` as the *recommended default* — agent per package, using the §3.2
  matrix + R7 dependency-init assessment. **Lean: `go vet`** (never executes code) unless linking and
  initializing must also be proven and the package + its imports are assessed safe.
- Whether to apply the optional clarity-line edits (§5.2) to other in-flight planned plans now, or only
  if no lane is concurrently mutating them — agent defers if in-flight.
- Exact phrasing of the lint's planted positive-control scratch fixtures (the sole-fence vs paired-fence
  fixtures) — agent may choose, but they must be enumerated within the lint's own corpus (R16-M2).

---

## Solution Design

### 3.1 Overview

The audit (Verification Log) enumerates all `git grep -l 'go build'` sites (88 tracked files) and
classifies the planned corpus. The sprint (a) documents the corrected-fence canon, (b) adds one
**executable** instrument-health guard (`scripts/lint_planned_fences.sh`, wired into
`scripts/verify_ail.sh`) so the flawed sole-fence cannot silently recur in a not-yet-shipped plan, and
(c) optionally clarifies the 4 operative planned plans that *pair* a redundant `go build ./...` rc=0
reading with an already-correct companion fence, so no future reader copies a misleading "go build
proves the test mutant builds" sentence. It does **not** rewrite history, and it does **not** touch the
whole-repo gates (verify_go.sh, ci.yml, README, CLAUDE.md).

### 3.2 What each command actually proves

| Fence | Compiles `_test.go` | Executes package code | rc=0 means | Caveat |
|---|---|---|---|---|
| `go build ./...` | **no** | no | all **non-test** packages build | Blind to every `*_test.go`; skips a test-only package entirely (so `host/verifygate/` is invisible to it) |
| `go test -count=1 -run '^$' ./host/<pkg>/` | **yes** | **may** — runs local `init()`, `TestMain`, **and any imported package `init()`** before filtering (side-effect freedom not guaranteed) | the test package compiles, links, and initializes, with zero tests run | **not side-effect-free in general**; `-count=1` **disables test-result caching only** (build cache retained; no forced recompile); safe only where the package and its imports are assessed init-free (R7) |
| `go vet ./host/<pkg>/` | **yes** | **no** (pure static analysis — never executes code) | the test package type-checks; analyzers pass | the most side-effect-free; **canonical default** |

**Rule to encode:** for any `_test.go` mutant, the "builds/compiles" proof must be `go vet
./host/<pkg>/` (preferred) and/or `go test -count=1 -run '^$' ./host/<pkg>/` — never `go build ./...`
alone. Prefer `go vet` unless linking/initializing must also be proven **and** the package + its
import closure are assessed side-effect-free (R7). `go build ./...` is retained only as a non-test-tree
/ whole-repo build fence.

### 3.3 The guard: an executable instrument-health lint (frozen)

**`scripts/lint_planned_fences.sh`** — a small, non-fleet, standalone script. It is an
**instrument-health control, not a semantic proof**: it enumerates *candidate* sole `go build`
test-mutant fences and classifies each, failing closed on anything it cannot cleanly classify. It
never discharges a load-bearing AC (that is a mutation's job, S6).

**Detection scope (precisely, not generic natural-language grep).** For every `.md` in the operative
corpus (default `design_docs/planned/*sprint-plan*.md`, overridable via `--corpus=<glob>`), the lint:
1. **Enumerates candidates**: any line that contains a `go build` command (`go build ./...` or
   `go build PATH`) **and** an assertion marker for a compile/build verdict (`rc=0`, "compiles",
   "builds", "compiled with"). (`grep -nE 'go build[^`]*(rc=0|compil|build)'`)
2. **Classifies each candidate** by the same record (line, or the immediately adjacent fence-record):
   - **paired** → the record also names a test-compiling companion (`go vet` or
     `go test …-run '^$'`/`-count=1`) → **PASS** (safe).
   - **sole + test-subject** → the record claims a `_test.go`/`test` mutant "builds" with `go build`
     and **no companion** → **RED** (the defect; must be corrected).
   - **sole + unambiguously non-test subject** → the `go build` target is a production path (e.g.
     `./host/daemon/` for a `store.go` mutant) or the whole-repo gate → **PASS** (valid uses, R15).
   - **unambiguous/ambiguous → RED-for-review** (fail closed; every unclassified candidate must be
     reviewed, per gpt6-astra).
3. **Exits distinctly**: `0` = clean scan (zero unclassified/RED) **and** the run itself succeeded;
   `1` = ≥1 RED candidate; `≥2` = execution error (missing file, glob matched nothing / fixture
   absent, or `timeout` expiry), reported from the clean scan so a script fault can never impersonate
   a green.

**Residual (what the guard does *not* prove, stated):** the lint does not parse prose meaning and
cannot prove negative-existence beyond its enumerated patterns; a `go build` fence written outside the
pattern set (renamed/obfuscated) would not be caught unless the pattern set stays in lockstep with the
§3.2 canon. It may false-positive a `go build` in a sentence that is really a whole-repo/non-test
reference — such candidates fail **closed to review**, never silently pass. It cannot prove that the
corrected fence compiles the mutant — only **M1** (a real mutation) discharges that. It is
instrument-health, never a load-bearing kill.

**Non-vacuity:** the lint always includes a **self-checked positive-control fixture** under
`scripts/testdata/planned-fence/` in its corpus, so a zero-candidate clean scan is meaningful; if that
fixture is missing the lint exits with the execution-error code (S6 fail-loudly). See M2.

**Conflict-surface reuse (answering gemini-3-1-pro's reuse objection):** rather than adding a new CI
job or a new harness, the guard **reuses the existing World-owned verification checkpoint
`scripts/verify_ail.sh`** as its invocation host — the local verify leg every sprint already runs and
the `ailang-verify` CI job executes. The lint is added as a single **new additive leg** that does not
alter the three existing legs or the S8 constants `EXACT_TOTAL_VERIFIED`/`EXACT_TOTAL_TESTS`.

### 3.4 Implementation Plan (single milestone)

**Phase 1: Canon + reproduction (~0.05d)**
- [ ] Record the corrected-fence canon (§3.2) and the measured reproduction into this doc's
      Verification Log (M1), so it is rerunnable by any future role.
- [ ] State the `host/verifygate/` is-a-test-only-package fact, the package-scoped caveat, the caching
      correction, and the dependency-initialization assessment (R7).

**Phase 2: Executable instrument-health guard (~0.08d)**
- [ ] Create `scripts/lint_planned_fences.sh` (candidate-enumerator + fail-closed classifier, §3.3),
      with a self-included positive-control fixture set under `scripts/testdata/planned-fence/`.
- [ ] Add one additive leg to `scripts/verify_ail.sh` invoking the lint (proposed non-fleet
      verification edit; does not touch S8 constants or existing legs).
- [ ] Execute M2 fixtures (PENDING) and record them; see §7.

**Phase 3: Optional clarity edits, if no lane conflict (~0.05d)**
- [ ] At the sites in §5.2, rephrase so the "mutant compiles" claim names the test-compiling fence and
      `go build ./...` is labelled as non-test-tree build only (or dropped from the test-mutant
      sentence).

**Deliberately absent (never edited):** `verify_go.sh`, `CLAUDE.md`, `.github/workflows/ci.yml`,
`README.md`, any implemented/history file, `tools/launchd/`, coding-standards, shared skills, accounts,
or any `.ail`.

---

## Files to Modify/Create

**New files**
- `design_docs/planned/w-go-test-compile-fence.md` (~this doc) — the canon, audit, and guard spec.
- `scripts/lint_planned_fences.sh` — the executable instrument-health guard (§3.3). **Proposed**;
  created in the execution phase.
- `scripts/testdata/planned-fence/` — the guard's self-included positive-control fixtures (M2).

**Modified files (narrowly-scoped non-fleet verification edit)**
- `scripts/verify_ail.sh` — one **additive** leg invoking `scripts/lint_planned_fences.sh` as the
  guard's durable trigger (reuse of the existing World verify checkpoint; must not touch
  `EXACT_TOTAL_VERIFIED`/`EXACT_TOTAL_TESTS` or the existing three legs). **Proposed**; implemented in
  the execution phase.

**Modified files (operative, optional clarity — only if not concurrently in flight)**

Each of the following is already *load-bearingly correct* (it carries a `go vet`/`-run '^$'`
companion); an audit classifies it as "compliant but reads as if `go build ./...` alone suffices".
The sprint may add one clarifying sentence per site. If any plan is mid-mutation by another lane, skip
it (deferred, §9).

- `design_docs/planned/w-verify-binary-lockfile-vlb-sprint-plan.md` — L460 ("Test files are not built
  by `go build ./...`. Type-check them with **`go vet`**"), L486/L547 (`go build ./... && go vet
  ./host/verifygate/`); rephrase so `go vet` is named as the **build** proof for test mutants.
- `design_docs/planned/w-capsule-output-cap-load-flake-sprint-plan.md` — L370–371/L420 name `go build`
  **plus** `go test ./host/capsule -run '^$' -count=1`; make the `-run '^$'` half the named "builds"
  definition.
- `design_docs/planned/w-ci-recovery-lever-absent-sprint-plan.md` — L119 "`go build ./...` rc=0 — the
  mutant compiles" is immediately caveated at L120; tighten so the caveat is not required to read
  correctly.
- `design_docs/planned/w-inventory-test-blind-to-asymmetric-addition-sprint-plan.md` — AC11
  (L83/L344) `go build ./...` rc=0 as a tree-hygiene clause; tag it as non-test-tree (it is a hygiene
  guard, not a mutant fence).

**Explicitly NOT modified** (classification; see Verification Log): `verify_go.sh`,
`.github/workflows/ci.yml`, `README.md`, `CLAUDE.md`, coding-standards, shared skills, every `.ail`,
every `implemented/` doc, every `design_docs/verification/*/…` mutation log, `world-mission*.md`,
`host/…`, `tools/launchd/`.

**Size estimate:** ~0.2d all-in. A single milestone; a second is not warranted.

---

## Testing Strategy / Acceptance Criteria

Each load-bearing criterion is discharged by a **mutation** that makes the fence red (S6); the guard
grep appears only as an explicitly **instrument-health** control.

| # | Criterion | Load-bearing? | How discharged / status |
|---|---|---|---|
| **AC1** | The corrected fence (`go vet ./host/verifygate/` AND `go test -count=1 -run '^$' ./host/verifygate/`) reds rc=1 naming the type error when a type-invalid function is appended to a `_test.go`, while `go build ./...` stays rc=0; restore is byte-identical to the captured sha256 and the pristine control is green either side | yes | **Mutation M1** — **confirmed** (controller-reproduced this iteration at `09726ae`; inherited; see §7, R3–R5). No repeat in the designer tree (controller note). Bounded: each leg under `timeout 180s`, timeout = failure. |
| **AC2** | Canon documented: what each fence proves + init/TestMain **+ dependency-init** caveat + the caching correction | yes (documentation) | Verification-Log row R2 / §3.2 / R7 |
| **AC3** | Instrument-health guard: `scripts/lint_planned_fences.sh` is an executable, wired into `scripts/verify_ail.sh`; the negative-existence scan over the operative planned corpus reports no **sole** `go build` test-mutant fence, and the same-scope **positive control** fixtures DO match (M2); execution errors report separately from a clean scan | no — instrument-health (S6 label, explicit) | **Mutation/control M2** — **designed, NOT executed: PENDING.** Cannot be confirmed until the lint exists and the fixtures run. Cannot discharge AC1. |
| **AC4** | No whole-repo gate, coding-standard, shared-skill, history or account changed: `verify_go.sh`, `ci.yml`, `README.md`, `CLAUDE.md`, coding-standards (S8 ratified), implemented docs, evidence logs all untouched (git status shows only this doc + proposed guard files) | yes (negative) | `git status --porcelain` + `git diff --stat -- <listed paths>` = empty; Verification-Log R18 |
| **AC5** | verify legs green outside any drill: `./scripts/verify_ail.sh` rc=0 (incl. the new lint leg, once wired); `go build ./...` rc=0; `go vet ./...` rc=0; `go test ./... -count=1` rc=0, with `AILANG_BIN` pinned | no (regression) | **PENDING the guard wiring.** `verify_go.sh` is excluded (three pre-existing FLEET-owned drift reds, row 76). |

### Per-branch mutation plan

- **M1 (load-bearing AC1, confirmed)** — append `func zzBroken() { var x int = "not an int"; _ = x }` to
  `host/verifygate/dispatch_lever_gate_test.go`; assert `go build ./...` rc=0 (flawed fence stays
  green), `go test -count=1 -run '^$' ./host/verifygate/` rc=1, `go vet ./host/verifygate/` rc=1;
  restore byte-identical (sha256 `2f4d0efa…`, R5), pristine control green both sides. **This is the
  only branch that proves the correction.** (Controller-reproduced; no repeat in designer tree.)
- **M2 (instrument-health AC3, PENDING)** — once `scripts/lint_planned_fences.sh` exists, run it with
  `--corpus` pointed at a temp dir whose planted fixtures are **inside that enumerator corpus**:
  (a) a **sole-fence fixture** (a line asserting a `_test.go` mutant builds with `go build ./...`
  alone) → **must RED (rc=1)**; (b) a **paired-fence fixture** (`go build` + `go vet`/`-run '^$'`
  companion on the same record) → **must PASS (rc=0, no false red)**; (c) a **clean file** (no
  candidate) → **rc=0 clean**; (d) an **execution-error arm** (fixture/glob absent) → **distinct
  rc≥2, reported separately from a clean scan**. Then run the same lint over the real operand corpus
  (expected zero candidates → clean; non-vacuous because the self-included positive fixture is always
  present). Every arm bounded by `timeout 30s`; a timeout is an execution error, not a clean scan.
  Rows R16–R17 stay **PENDING** until this executes.

---

## Non-Goals

- **Not a CI job change and not a new harness.** `.github/workflows/ci.yml` is untouched; the guard
  rides the existing `ailang-verify`/local `scripts/verify_ail.sh` leg (reuse, not new surface).
  `verify_go.sh` is **not** edited (fleet-drift arm + whole-repo Go gate preserved; row 76).
- **Not a triggerless grep.** The guard is an executable wired into an invoked verify leg (primary) or,
  solely as an interim fallback before that leg lands, a **mandated manual trigger** with a bounded
  deadline (Guard Lifecycle, §3.3). A grep that nothing invokes is never presented as a guard.
- **No rewrite of historical measured records.** Implemented docs, evidence/mutation logs and
  world-mission STATUS records that used the flawed fence are append-only measurements; classified and
  untouched. — [Append-only standing rule]
- **No edit to frozen/core surfaces**: `tools/launchd/`, coding-standards (ratification-class, S8), the
  shared `mission-control` SKILL.md (upstream, proposed to V1/Mark), accounts, and every `.ail`.
- **No AILANG code or language change.** The fence is Go tooling only. — [Mission directive]
- **Not a repo-wide text ban on `go build`.** It remains a correct fence for non-test mutants and the
  whole-repo build (R15); a blanket replacement would be the over-rotation the audit rejects.

---

## Why this is not a package (S3)

The deliverable is a corrected **verification/convention** discipline on the host-layer gates and an
**executable verify-leg guard** — not a new domain behavior, transition law, policy, or tool that ships
as a registry-published AILANG package. The mutation drills and verify gates that consume it live in
the Go host / CI surface, and a doc/plan-integrity guard is a verify-gate instrument of exactly the
class S3 keeps in `scripts/`, not in a package-able domain module. Nothing here adds a package API or a
`.ail` surface. → **Not a package.**

---

## Axiom Compliance

| Axiom | Score | Justification |
|---|---|---|
| A1 Determinism | +1 | Makes "the mutant compiles" a deterministic, honest per-package check instead of a go-build exit code that never varies for test mutants |
| A2 Replayability | +1 | The pinned reproduction (M1) is rerunnable verbatim by any future role; the guard's fixtures are checked in |
| A3 Effect Legibility | +1 | Documents that `go test -run '^$'` may execute init/TestMain + imported init, making the fence's hidden effects legible |
| A4 Explicit Authority | 0 | No capability changes |
| A5 Bounded Verification | +1 | Enables local, cheap check for each package's test compile |
| A6 Safe Concurrency | 0 | No concurrency impact |
| A7 Machines First | +1 | Removes a known-exit-code ambiguity AI executors copy from docs; fewer mis-measures per mutant |
| A8 Minimal Syntax | 0 | No new syntax |
| A9 Cost Visibility | 0 | No resource changes |
| A10 Composability | 0 | No API change |
| A11 Structured Failure | +1 | Replaces a silent false-green fence with one that names the type error; the guard separates execution error from a clean scan |
| A12 System Boundary | 0 | No boundary change |

**Net Score: +6** ✅ Proceed. *(Corrected: the +1 rows total 6, not 7.)*

**Hard Violation Check:** A1 no; A3 no (effects made legible, not hidden); A4 no; A7 no.

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| A future plan regresses to the sole `go build` fence anyway | Med | The guard is an **executable** wired into `verify_ail.sh` (runs on every verify/CI leg), not a dormant doc note; its positive control makes the second instance recognized rather than rediscovered; canon is written down |
| The guard wiring into `verify_ail.sh` is judged a core-gate edit | Med | It is a single **additive** leg touching no S8 const and no existing leg; if the controller rules it off-limits, the frozen fallback is the mandated manual trigger (bounded deadline, §3.3). The executable itself is a new non-fleet script in either case |
| `go test -run '^$'` on a package whose init/import-init runs code | Med | Canon defaults to `go vet` (never executes) except where linking is required and R7's dependency-init assessment passes; `host/verifygate/` assessed (local init/TestMain empty; imports stdlib-only, dep-init PENDING) |
| Editing other planned plans conflicts with a concurrently-mutating lane | Low | §5.2 clarity edits are optional and deferred if a lane is in flight (§9) |
| The audit's "already compliant" finding is mistaken | Low | Every classification cites a real `git grep` line (R9–R15); M1 and the M2 positive control are on record |
| A reviewer treats AC3's grep as load-bearing | Low | AC3 explicitly labelled instrument-health; AC1 is the only load-bearing fence (mutation-discharged), per S6; M2 labelled PENDING until executed |

---

## References

- **Queue item**: `design_docs/world-mission.md` row 65 (this work).
- **Discovery**: `design_docs/planned/sprint_w-dispatch-lever-parser-false-reds-on-valid-yaml.json` —
  row 55's plan (landed iter-**149**), non-blocking finding 5, already uses the corrected `-run '^$'`
  shape (copy in-tree).
- **Related implemented**: `w-load-bearing-criteria-need-a-mutation-not-a-grep.md` (grep ≠ mutation,
  S6), `w-canary-fence-passes-a-gutted-canary.md` (instrument-class gap), `w-boundary-gate-tree-mutation.md`
  (build gate blind to a mutant), `w-miscompile-instrument-inert-in-ci.md` (instrument-inert class).
  Distinct: those kill vacuously-green *assertions*; this kills a *fence* that cannot see test files.
- **Standards**: `design_docs/coding-standards.md` S6 (load-bearing criteria require mutations).
- **Method**: related-doc search was **static** (`git grep`/`ls`), **not** neural/embedding (no
  local-model work this iteration).
- **Upstream proposal (out of World's hand)**: shared `mission-control` SKILL.md prescribes this fence;
  World has proposed it to V1/Mark, cannot edit it.

---

## Verification Log (current-code claims, command + observed output)

All reproduction claims (R2–R5) are **controller-reproduced first-party on HEAD `09726ae`**, go1.26.6
darwin/arm64, byte-identical restore — inherited; **not repeated in the designer tree** (controller
note: no need to repeat the full regression suite). All enumeration/citation/negative claims (R0, R1,
R6–R15) were **re-grepped first-party this designer session** at HEAD `09726ae` and quote the actual
output lines. Every guard-execution claim (R16–R17) is **PENDING** until the lint exists and runs.
Mutations were applied and **restored byte-identical** (captured/asserted sha256); `git status
--porcelain` returned clean (R18). All drills are bounded by `timeout`, with a timeout treated as
failure.

| # | Premise (claim) | Command (bounded) | Observed | Verdict |
|---|---|---|---|---|
| R0 | Corpus inventory: `go build` occurs in 88 tracked files | `git grep -l 'go build' \| wc -l` | **88** | confirmed |
| R1 | `host/verifygate/` compiles only as tests (test-only package) | `ls host/verifygate/*.go` | **8** `*_test.go`, **0** non-test `.go` (controller-confirmed same shape) | confirmed — `go build ./...` sees **none** of this package |
| R2 | Baseline is green on the pristine tree | controller-reproduced: `go build ./...`; `go test -count=1 -run '^$' ./host/verifygate/`; `go vet ./host/verifygate/` | **rc=0** / **rc=0** / **rc=0** | confirmed (inherited, controller-reproduced) |
| R3 | The broken fence is blind to a `_test.go` type-error **mutant** | append `func zzBroken() { var x int = "not an int"; _ = x }` to `dispatch_lever_gate_test.go`; `timeout 180s go build ./...` | **rc=0** (build stays green on a non-compiling test) | confirmed — the flaw |
| R4 | The corrected fences red on the SAME mutant | `timeout 180s go test -count=1 -run '^$' ./host/verifygate/`; `timeout 180s go vet ./host/verifygate/` | **rc=1** both: `cannot use "not an int" ... as int value in variable declaration` | confirmed — the fix |
| R5 | Restore is byte-identical to pristine | `cp`-free restore; `shasum -a 256` vs captured `2f4d0efa…`; `git status --porcelain` | sha `2f4d0efa13ad9dab015b85ce5eb7482d0888e3a338c8a7f4f24777dbd60bef2e` restored **YES**; porcelain empty | confirmed |
| R6 | CI has no hole — verify_go.sh compiles test files | `grep -n 'go test ./\.\.\. -count=1' scripts/verify_go.sh` | present at L472–474 (also L5) | confirmed — AC-instrument, not CI |
| R7 | `host/verifygate/` has no **local** `init`/`TestMain` | `grep -rn 'func TestMain\|func init(' host/verifygate/ \| wc -l` | **0** hits | confirmed (**local only**) |
| R7b | Dependency-initialization assessment for the `-run '^$'` caveat | read imports of the 8 `host/verifygate` test files | imports are **stdlib-only** (`go/parser`, `go/types`, `os/exec`, `regexp`, `runtime`, etc.); package-level `init()` of these is not locally measured | **PENDING assessment** — a local no-init grep cannot prove imported-init absence; mitigated because the **canonical default is `go vet`** (never executes), so the risk only applies when `-run '^$'` is deliberately chosen |
| R8 | **Enumeration + scope:** the operative corpus is `design_docs/planned/*sprint-plan*.md`; every `go build` occurrence there is classified (paired vs sole) | `git grep -nE 'go build' -- 'design_docs/planned/*sprint-plan*.md'` | all 12 matching sprint-plan files carry a test-compiling companion (`go vet`/`-run '^$'`); **no sole-fence test-mutant claim found** (R9–R15 classify each) | confirmed — corpus already compliant (§5) |
| R9 | row-55 / `w-dispatch-lever-parser` uses the corrected fence | read plan; grep | L203/L280–281: `go test -count=1 -run '^$' ./host/verifygate/` + `go build ./...` | confirmed (corrected copy in-tree) |
| R10 | `w-canary-fence-blind-to-a-skipped-canary` §0.1 names the gap and uses the corrected fence | grep L20/L22/L25/L99 | L20 "`go build` IS NOT A COMPILE FENCE FOR A `_test.go`"; L99 `go test -count=1 -run '^$' ./host/verifygate/` | confirmed compliant |
| R11 | `w-ail-gate-module-pin` typechecks test mutants with `go vet` | grep L401/L607/L750–751 | "`go vet ./host/verifygate/` … `go build ./...` **does not compile `_test.go` at all**" | confirmed compliant |
| R12 | `w-daemon-timeout-test-flake` builds its `_test.go` case separately and names the gap on a production mutant | grep L563–564 | L564 "`go build ./...` DOES NOT COMPILE _test.go AT ALL — a build on a test-only mutant is VACUOUS." | confirmed compliant |
| R13 | `w-ci-recovery-lever-absent` (planned) carries a caveated paired reading | grep L119–120 | L119 "`go build ./...` rc=0 — the mutant compiles"; L120 "…a *broken* test file would leave `go build ./...` at rc=0 while `go vet` reds.**" (immediately caveated); companion `go vet ./host/verifygate/` is M6's typecheck (L347) | compliant-but-ambiguous → §5.2 clarity candidate |
| R14 | `w-capsule`, `w-inventory AC11`, `w-verify-binary-lockfile` (planned) pair `go build` WITH a test-compiling fence, never alone | grep capsule L370–371/L420; inventory L83/L344; lockfile L460/L486/L547 | capsule L371 "`go test ./host/capsule -run '^$' -count=1` … `go build` does **not** compile"; lockfile L460 "Test files are not built by `go build ./...`. Type-check them with **`go vet ./host/verifygate/`**"; L486/L547 `go build ./... && go vet ./host/verifygate/` | compliant; §5.2 clarity candidates (line numbers corrected vs prior draft) |
| R15 | **In-corpus positive examples (same corpus as R8's negative):** where `go build` IS the right BUILD proof, the subject is non-test; where it is paired with a companion, the record is safe — the classifier (R8) distinguishes sole vs paired vs non-test | grep (this corpus, not external `verification/` docs) | `w-daemon` L460/L563 run `go build ./host/daemon/` for a **PRODUCTION** (`store.go`) mutant — valid sole `go build`; `lockfile`/`capsule` pair `go build` with `go vet`/`-run '^$'` — valid paired | confirmed — establishes the classification is real, in the enumerated corpus |
| R16 | **Known-negative (same scope):** a `_test.go` mutant fenced by `go build` alone would be **detectable by the lint** via a planted fixture | **PENDING** — requires `scripts/lint_planned_fences.sh` to exist and M2 fixture (a) run | not yet run | **PENDING, not confirmed** |
| R17 | The lint reports **execution errors separately from a clean scan** (and timeout = failure) | **PENDING** — requires M2 fixture (d) + timeout arm run | not yet run | **PENDING, not confirmed** |
| R18 | No whole-repo gate / coding-standard / shared-skill / history / account surface touched by this designer session | `git status --porcelain` (designer tree) | only `?? design_docs/planned/w-go-test-compile-fence.md` untracked; no listed-path diffs | confirmed (design-only; no gates edited) |

**Interpretation of R8–R15 (the "exact operative corpus"):** the sole-fence defect is **not live** in
the current planned corpus; every operative plan carries a test-compiling companion, and the valid
`go build` uses are for non-test subjects (R15). The historical documents that show the flawed fence
are append-only and are **not** changed. The durable value of this row is the **canon + an *executable,
invoked* guard** so the defect is never silently re-introduced — "Instance 1 pre-registered so the
second is recognisable rather than rediscovered" (row 65).

---

## Guard Lifecycle (explicit trigger — the guard is never triggerless)

- **Primary trigger (durable):** `scripts/lint_planned_fences.sh` runs as an additive leg of
  `scripts/verify_ail.sh`. Every standard local verify run and the CI `ailang-verify` job execute it;
  a future agent that runs the repo's standard verify gate gets the guard for free. **No executorship
  step can skip it by merely not reading this doc.**
- **Interim trigger (until the verify_ail.sh leg is wired, and as a fallback if a gate-core ruling
  blocks that wiring):** the sprint MUST run `timeout 30s ./scripts/lint_planned_fences.sh
  --corpus='design_docs/planned/*sprint-plan*.md'` at **every planned-doc merge review**, and record
  the result in the verification log; a red or an execution error blocks the merge. Bounded wall-clock
  deadline; timeout = execution error = block.
- **Non-vacuity invariant:** the lint always includes its self-checked positive-control fixture
  (`scripts/testdata/planned-fence/`), so a zero-candidate clean scan is meaningful and a missing
  fixture is a loud execution error (S6), never a silent green.

---

## Unresolved choices (for quorum / implementer)

1. **Guard host — RESOLVED (frozen):** standalone executable `scripts/lint_planned_fences.sh` wired as
   an additive leg of `scripts/verify_ail.sh`. Removes the previous "documented grep" option rejected
   by all three reviewers. (No longer open.)
2. **Default fence:** `go vet` (never executes code — **default**) vs `go test -run '^$'` (also
   links/initializes; only where R7's dependency-init assessment passes). — agent per package.
3. **§5.2 clarity edits:** apply now vs. defer if any target plan is concurrently in-flight. Agent
   defers on conflict.
4. **AC5's `verify_go.sh` leg:** excluded from AC5 because of its three inherited FLEET-owned drift
   reds (row 76). If quorum insists `verify_go.sh` must be in an AC, that is a row-76 dependency, not
   this row.

---

## Single-milestone statement

One milestone: (1) canon + reproduction (§3.2, M1 confirmed, R1–R7), (2) the **executable**
instrument-health guard + wiring (M2 executed → then AC3/AC5 **PENDING→confirmed**), (3) optional
clarity lines (§5.2) if no lane conflict. ~0.2d. No second milestone warranted. Guard-execution
evidence stays **pending** until the lint exists and its fixtures run; nothing is claimed confirmed
before execution.

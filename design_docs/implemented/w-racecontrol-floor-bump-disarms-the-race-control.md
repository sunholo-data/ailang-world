# w-racecontrol-floor-bump-disarms-the-race-control

**Status**: Planned
**Target**: current iteration (World iter-137 / P48)
**Priority**: P1
**Estimated**: ~0.1d (~1h)
**Dependencies**: None
**Planner-Lane**: codex-ok

## Problem Statement

`design_docs/verification/w-race-gate-blindspot/racecontrol/` is the second of the repo's
three nested `go.mod` modules, and the only one that hosts a **known-positive control** for
`verify_go.sh`'s race-detector gate. `scripts/verify_go.sh:229` runs
`cd .../racecontrol && go run -race .`, and `:232` requires `WARNING: DATA RACE` in the output
or it FATALs `"the race detector is not armed; every 0-races result in this gate is void"`.
The module deliberately contains a data race (`main.go`: 4 goroutines incrementing `shared`)
so the control is supposed to fire on every run.

That control is **disarmed by a blanket module-floor bump**, the same defect class queue row
42 just fixed on the *repro* module. Row 42's own Systemic-Issue Audit row **V17** examined
`racecontrol/` and declared it a *non-instance* — *"driven only by `verify_go.sh:229` `go run
-race .` under the **default** toolchain, so a floor bump there cannot disarm anything"* — and
round 2's quorum reviewer objected to that trailing clause as an unmeasured premise. Measured
first-party as **V24** in `w-canary-control-does-not-survive-a-floor-raise.md`, and reproduced
verbatim at HEAD for this item (V1, V2):

- Inside the module the ambient toolchain is **`go1.26.4`** (the Homebrew base install), *not*
  the root-selected `go1.26.6` — because `racecontrol/go.mod` declares `go 1.22`, which is
  *below* the ambient base, so `GOTOOLCHAIN=auto` has nothing to push up to and Go falls back
  to the base binary.
- Base `go run -race .` → rc=1 with **2** `WARNING: DATA RACE` lines (the known-positive fires).
- Bump the directive `go 1.22` → `go 1.26.6`: with `GOTOOLCHAIN=auto` the control *still* fires
  — **only because `go1.26.6` is already in the local toolchain cache** and `auto` silently
  switched toolchains to satisfy the new floor.
- Add the one variable that removes that rescue, `GOTOOLCHAIN=local`: rc=1, output is exactly
  `go: go.mod requires go >= 1.26.6 (running go 1.26.4; GOTOOLCHAIN=local)`, **0** `WARNING:
  DATA RACE` lines — so `verify_go.sh:232` FATALs and the race-detector gate is dead. The
  control fails *before* execution; there is no race warning to count.

**The disarming is real, and its only concealment is a warm toolchain cache** — the exact
condition a fresh runner does not have. Row 42 filed this as **queue row 48** precisely because
a non-instance that an audit *publishes as measured* is "an unbound claim wearing a
measurement's clothes."

**Quorum round 1 BLOCKED at full strength** (all three external reviewers present) on this
draft's original direction — anchoring the invariant to the **oldest `KNOWN_GOOD`** toolchain.
The block and its resolutions are recorded in [Quorum History](#quorum-history). Two of the
three objections force the redesign in this revision:

1. **DIRECTIONAL (gpt5-6-sol):** `run.sh`'s `KNOWN_GOOD` list does **not** constrain the
   ambient toolchain used at `verify_go.sh:229` (a bare `go run -race .` under nested
   auto-selection). Bounding the floor `<= go1.24.9` proves nothing about whether the control
   is executable under the toolchain that actually runs it. The fix must **bind execution**
   with machinery, not assert against an unrelated list.
2. **PREMISE (oc-glm-5-2):** whether a `LOAD-BEARING` comment containing `go 1.22` breaks
   `moduleGoFloor`. Measured and recorded in the Verification Log (V8); a post-comment
   acceptance arm is added to AC7.

The third (non-blocking) reviewer, **gemini**, corrected a syntax error: `go env` takes
space-separated keys (`go env GOTOOLCHAIN GOVERSION`), not a comma-separated shorthand. Fixed
throughout this revision.

**Quorum round 2 raised three blockers** (recorded in [Quorum History](#quorum-history)):
R2-a that the doc never proved the new test actually runs in CI; R2-b that a `go test -run`
selector written with a grep-style escaped alternation `\|` executed zero tests while `go test`
printed `ok … [no tests to run]`, rc=0, which the doc had banked as a green baseline; and R2-c
— the **design-DIRECTION dispute** that parked the item for the human — that the revised proof's
middle premise (`ACTIVE_GO >= root floor` is always true under `GOTOOLCHAIN=auto`) is **false in
general**: under `GOTOOLCHAIN=local` with a base below the root floor, `go env GOVERSION` can
return that lower base, so the static implication alone is refutable. Mark Edmondson ruled at an
attended session, **2026-09-01**, recorded in the mission charter's decision ledger as
**D-WORLD-28** (ANSWERED A): *fail closed unless the selected `ACTIVE_GO` is at-or-above the
root module floor, then bind the race-control invocation to that exact `ACTIVE_GO`, with the
static `racecontrol/go.mod` floor <= root floor requirement, including the runtime floor check
round 2 found necessary.* Option B (pin `GOTOOLCHAIN=local`) is rejected by that ruling. This
revision implements the ruling as one three-part proof (P1 runtime fail-closed, P2 execution
binding, P3 static invariant).

**Impact:**
- A blanket floor raise (`go mod tidy -go=…`, an IDE "update module", or a human matching the
  root `go 1.26.6`) silently destroys the *only* first-party positive control for the race
  detector, the instrument that stands behind every `-race` result in the go gate.
- The failure is **loud** (`verify_go.sh` FATALs) — which is why this is a row, not a blocker —
  but it is a **false** loud failure: it blames the race detector when the actual cause is a
  floor mismatch, and on a fresh runner it would fail every go gate run for the wrong reason.
- Row 42 explicitly did **not** claim systemic completeness; this item closes the gap it
  measured and deferred.

## Goals

**Primary Goal:** Bind the `racecontrol` module floor to the **root module's pinned toolchain**
— the toolchain the race control is now *made to run under* — so a floor raise that would
disarm the race-detector known-positive control reds the static gate before any
`go run -race .` runs, and so the runtime lane executes the control under a **known,
deny-list-vetted toolchain** rather than whatever nested auto-selection happens to pick. That
toolchain is now guaranteed at-or-above the root module floor **by enforcement** (D-WORLD-28):
`verify_go.sh` fails closed unless the selected `ACTIVE_GO` is at-or-above the floor read from
`./go.mod`, then binds the race leg to that exact `ACTIVE_GO` — the premise the round-2
proof relied on by assumption becomes a runtime check.

**Success Metrics:**
- `verify_go.sh`'s race leg invokes the control with `GOTOOLCHAIN="$ACTIVE_GO"` (the root
  module's `go env GOVERSION`), binding execution to the same toolchain the deny-list just
  vetted — no nested auto-selection can differ.
- `verify_go.sh` **fails closed** toward the root module floor: after `ACTIVE_GO` is captured
  and the miscompile deny-list exits, a runtime check FATALs + exit 1 unless `ACTIVE_GO` is
  at-or-above the floor read from `./go.mod` (D-WORLD-28). This makes `ACTIVE_GO >= root floor`
  true by enforcement, closing R2-c; the check is independent of the deny-list, so a base below
  the root floor but outside the deny-list (e.g. a `go1.25.x` base) is still caught.
- The `racecontrol` floor is bound by a static test in the gating lane (`host/verifygate`),
  asserting `racecontrol floor <= root module floor`, so CI reds on a disarming floor bump even
  though the runtime lane only garbles it.
- The static test also pins the execution-binding needle, so reverting the `verify_go.sh` edit
  (back to a bare `go run -race .`) is a named RED — the exact objection-1 hole cannot
  silently return.
- Every refusal/anti-vacuity branch of the new test has a named RED mutant; the disarming bump
  (floor above the root toolchain) is a named mutant.
- All three nested `go.mod` floors are bound after this item (root by existing tests, `repro`
  by row 42's Test A, `racecontrol` by this item) — the systemic audit closes out.
- The `racecontrol` floor directive is byte-unchanged (`go 1.22`); the only production edit to
  `racecontrol/go.mod` is a `LOAD-BEARING` comment.

## Quorum History

**Round 1 — BLOCKED at full strength** (all 3 external reviewers present). Decision: the round
is **not** a waiver, and the two blocking objections are accepted. The third is a non-blocking
syntax correction. This revision resolves all three; the redesign is grounded in first-party
measurements recorded in the [Verification Log](#verification-log).

| Reviewer | Objection | Verdict | Resolution (measured first-party) |
|----------|-----------|---------|-----------------------------------|
| gpt5-6-sol | **DIRECTIONAL** — anchoring the floor to the oldest `KNOWN_GOOD` (`go1.24.9`) is rejected. `run.sh`'s list does not constrain the ambient toolchain used at `verify_go.sh:229`; `floor <= go1.24.9` does not prove the control is executable under every possible host base. Re-design around machinery that actually binds execution. | **Blocking** | Re-anchored to the **root module floor** (`go1.26.6`) and added an **execution-binding production edit**: `verify_go.sh:229` now runs `GOTOOLCHAIN="$ACTIVE_GO" go run -race .`, where `ACTIVE_GO` is the root module's `go env GOVERSION` (V3). Under `GOTOOLCHAIN=auto`, `ACTIVE_GO >= root floor` always, so `racecontrol floor <= root floor` implies `racecontrol floor <= ACTIVE_GO` — the floor is satisfiable by the very toolchain that runs the control, under every possible host base (V2, V3). All `KNOWN_GOOD`-anchor claims are dropped; no claim that `KNOWN_GOOD` constrains this path remains. |
| oc-glm-5-2 | **PREMISE** — does a `LOAD-BEARING` comment containing `go 1.22` break `moduleGoFloor`'s exactly-one-`go `-line floor? | **Blocking** | Measured (V8): `moduleGoFloor` uses `strings.HasPrefix(line, "go ")`; a comment line begins with `//` at column 0, so it is ignored. A probe file containing the comment plus the real directive yields exactly one awk-equivalent column-0 `go ` line (`floors=[go1.22] count=1`), and the original file was restored byte-identically (sha256 `ab782f11…`). Recorded in the Verification Log (V8) and a **post-comment acceptance arm** is added to AC7. |
| gemini | **SYNTAX (non-blocking)** — write `go env GOTOOLCHAIN GOVERSION`, not a comma-separated shorthand. | Non-blocking, accepted | Corrected. All `go env` invocations in this revision use space-separated keys; confirmed valid (V3: `auto go1.26.6`). |

The round-1 reviewers' substantive direction (bind execution, don't assert against an unrelated
list) is now the design; the round-2 surface (row 44 CI inertness) is untouched by this
revision.

**Round 2 — BLOCKED at full strength** with three blockers, all dispositions (not waivers).
R2-a and R2-b are closed by first-party measurement in this revision (V10, V11); R2-c was a
design-DIRECTION dispute parked for the human and settled by the attended ruling **D-WORLD-28**
(2026-09-01). The ruling is binding and not re-litigated here.

| # | Block | Disposition (measured first-party at HEAD 8bb9214) |
|----|-------|------------------------------------------------------|
| R2-a | **CI linkage unverified** — the doc never proved the new test actually runs in CI | Closed by measurement (V10). `.github/workflows/ci.yml:166` runs `./scripts/verify_go.sh` (step name at `:163`, "go build + test gate (replay tests run against pinned AILANG_BIN)"); `verify_go.sh:258` runs `go test ./... -count=1`; `go list ./...` includes `github.com/sunholo-data/ailang-world/host/verifygate`. A test added to `host/verifygate/` is reached by CI through the `./...` leg. Added **AC9** asserting this linkage mechanically (three greps, each with its own control). |
| R2-b | **an escaped `\|` baseline regex executed zero tests** — `go test -run` with a grep-style escaped alternation `\|` is a *literal pipe* to Go's regexp, so it matches nothing and `go test` exits 0 printing `[no tests to run]`; the doc banked that rc=0 as a green baseline | Closed by audit (V11). Every `-run` selector in the doc was enumerated and re-run. The only malformed one was V9's escaped-`\|` alternation, which ran **0** tests (`ok … [no tests to run]`, rc=0) and covered the five sibling pin tests vacuously. Fixed to a single-pipe `\|` alternation: all five run (5× `=== RUN`, 5× `--- PASS`, rc=0). Every go-test AC now binds run-existence (counted `=== RUN`/`--- PASS`), never a bare rc (propagated from AC1); **AC10** adds the selector-vacuity rule. |
| R2-c | **GOTOOLCHAIN=local below the root floor refutes the revised proof** — the old proof's middle premise (`ACTIVE_GO >= root floor` is always true under `GOTOOLCHAIN=auto`) is false in general | Settled by D-WORLD-28 (2026-09-01). P1 makes the premise true **by enforcement** (fail closed below the root floor) instead of by assumption; P2 binds the `racecontrol` invocation to the exact `ACTIVE_GO`; P3 is the static `racecontrol floor <= root floor` test. Measured (V12): at root, `go env GOTOOLCHAIN GOVERSION` → `auto go1.26.6`; `GOTOOLCHAIN=local go env GOVERSION` → `go1.26.4` (below root floor `1.26.6`); `GOTOOLCHAIN=local go version` → `go1.26.4 darwin/arm64`; inside `racecontrol/` → `auto go1.26.4`. **Option B (pin `GOTOOLCHAIN=local`) is rejected** by the ruling; only the existing Deferred-Decisions note names it. |

The round-1 direction (bind execution, don't assert against an unrelated list) stands; the
ruling's P1 converts the static proof's middle premise from an assumption into an enforced
runtime floor check.

**Round 3 — BLOCKED at full strength, and every objection lands on ONE surface: P1's shell block
was specified in prose and measured nowhere.** `gpt5-6-sol` was absent on `budget` in the round
itself and was RESTORED by a single re-run at a raised cap (`--max-cost-usd 0.40`, $0.115775) per
the shared skill's absent-reviewer rule, so this is a 3-of-3 external reject, not a 2-of-3.
**No reviewer disputed the design DIRECTION** — it is ratified by the attended human ruling
`D-WORLD-28` (2026-09-01) — so all three are completeness/verification objections with concrete
reviewer-authored fixes, and the controller applied them under the narrow-refinement carve-out by
MEASURING what each asked for rather than asserting it.

| Reviewer | Objection | Verdict | Resolution (measured first-party by the controller, 2026-09-01, HEAD `8bb9214`) |
|----------|-----------|---------|-----------------------------------|
| gpt5-6-sol | The P1 fail-closed mechanism is asserted rather than verified: `V14` said "plan, not run here", the exact shell block was never given, and `V15` verified only an isolated comparator — not root-floor extraction, malformed input, placement before the race leg, exit behaviour, or deny-list independence. | **Blocking** | The exact block is now IN the design (above), and every branch it names is measured. `V14` is replaced by four runtime arms against a temp patched copy of `verify_go.sh`; `V16` is the branch battery (equal / above / below / not-deny-listed / ordering-trap / malformed / prerelease / devel / missing / duplicate / indented); `V18` is the extraction and the lexical fail-open control. Production `verify_go.sh` was never edited (sha256 `27eab122…` asserted before and after). |
| gemini-3-1-pro | The V15 comparator validates `^go1\.[0-9]+(\.[0-9]+)?$`, but the doc never measured the shell extraction that must produce such a token: `ACTIVE_GO` is `go1.26.6` while `go.mod` carries `go 1.26.6`, so a raw read yields `1.26.6` and the comparator would exit 2 and fail the gate closed for the wrong reason. | **Blocking** | **CONFIRMED and closed by measurement (V18).** The raw `go.mod` token `1.26.6` matches the V15 grammar **0** times; the `go`-prefixed form matches **1** — the normalisation is load-bearing, exactly as objected. The `ROOT_FLOOR="go$(awk '/^go /{print $2; exit}' go.mod)"` extraction is now written into the block above and measured in V16/V18, and the reviewer's own proposed extraction is the one adopted. |
| oc-glm-5-2 | `V14` admits M9/M10 are a plan, not a measurement, while AC11 and the Axiom-Compliance rows cite them as evidence; `V12`/`V13` do not exercise P1 at all. "P1 fails closed independently of the deny-list" was unmeasured. | **Blocking** | **CONFIRMED and closed by running exactly the prototype the reviewer specified.** A temp patched copy of `verify_go.sh` carries the P1 block; with the miscompile deny-list NEUTERED and `GOTOOLCHAIN=local` (`ACTIVE_GO=go1.26.4`), P1 FATALs with its own attributed message and the race leg is never reached (rc=1). **The discriminating control the objection demanded**: the byte-identical scenario with the P1 block REMOVED reaches the race leg and arms the control (rc=0, 2× `WARNING: DATA RACE`) — so P1 alone supplies the red. `V14`'s "plan, not run here" text is deleted, not softened. |

**Round 3, RE-REVIEWED AFTER the carve-out revision (iteration 145 controller, 2026-09-01).**
The absent-reviewer rule was re-applied to the *revised* doc rather than only to the round itself:
`ailang design-review … --reviewer gpt5-6-sol --max-cost-usd 0.30` → **`reject`**, $0.13903,
25,748 input tokens. Its objection is recorded here because it was **independently confirmed** and
because it is the reason this revision exists: *"The document explicitly measures that deleting the
entire mandatory P1 block is a silent green in both static and runtime lanes (M11), yet the
implementable design still specifies only the P2 execution-binding needle. No exact P1
semantic-presence assertion appears in the Components section, acceptance criteria, or the M1–M10
mutation table. The Conflict Surface merely promises an undefined P1 needle, while the planner notes
refer to nonexistent M11/M12 coverage."* That is **correct as stated**: iteration 144 recorded the
planner's nine refutations in this Quorum History as prose and never propagated them into the
normative sections, so the doc's own arms table still carried the two refuted arms (M6, M10) and no
P1 needle. The reviewer's `proposed_fix` — *"Add a concrete P1 semantic assertion … Add named RED
mutants for whole-block deletion and each partial semantic deletion, including run-existence checks,
and make them explicit ACs. Then replace the stale M6/M10/AC7 text with the measured M6′,
M10a′/M10b, and diff-based directive-preservation checks; renumber the mutation table consistently
and correct the verify_go.sh LOC estimate."* — has been applied in full and **verbatim in scope**:
**AC12** (P1a–P1f), mutants **M11–M18** with the **M17** green control, **M6′**, **M10a′/M10b**,
AC7's diff-shape clauses, the AC9 anchored grep, the AC2 citation repair, `mv` in M8, and the
`+35/−0 plus 1 altered line` correction. No controller-invented resolution was substituted for any
of it, and the design **DIRECTION is untouched** — it is ratified by the attended human ruling
`D-WORLD-28` and was disputed by no reviewer in any round. Every fix applied here was **already
measured first-party** by the `opus` planner against the LANDED artifacts before this revision was
written (rows V19–V23), so this is the narrow-refinement carve-out satisfying objections with
measurements, not with argument.

A harness fact worth recording because it produced four outcome-identical arms before it was
caught: `verify_go.sh` line 17 is `cd "$(dirname "$0")/.."`, so a patched copy executed from
`/tmp` cds to `/` and dies in the tracked-binary hygiene gate long before P1. The copies must live
in `scripts/` (they were deleted afterwards; porcelain 0). The only reason this did not become a
banked "all four arms red" was the explicit `reached_race_leg` marker counted in every arm — rule
3a's known-positive discipline aimed at the harness rather than at the subject.

**PLANNER REFUTATIONS (opus, lane `opus fail-closed:env-pin`, 2026-09-01) — nine, four of them
load-bearing, all measured against the LANDED artifacts rather than a replica. They are recorded
here because a doc that keeps asserting a refuted arm will have that arm re-derived by every later
reader.** The planner landed P1 verbatim and could NOT measure a defect in the block itself; every
refutation below is about the doc's *arms and citations*, not its design.

1. **M6 is VACUOUS.** Root `go banana` makes Go's module loader reject `go.mod` before the test
   binary builds: `rc=1` with `=== RUN` **0**, `--- PASS` **0**, `--- FAIL` **0**
   (`go.mod:3: invalid go version 'banana'`). Under this doc's own AC10 rule an arm with RUN=0 is
   not evidence, so the test's `!version.IsValid(rootFloor)` branch had **no killer at all**.
   Replaced by **M6′**: root floor `go 1.26.6 // pin` → RUN=1, FAIL=1 on that exact branch.
2. **M10 AS WRITTEN IS KILLED BY THE MISCOMPILE DENY-LIST, NOT BY P1 — the doc's own V14/A3
   attribution shape, failed by the doc.** Deleting the root `go` directive removes what pushes
   `GOTOOLCHAIN=auto` upward, so `go env GOVERSION` drops to `go1.26.4`, which `:218–224` already
   denies; the P1-REMOVED control produces a **byte-identical** red. Replaced by two arms that hold
   `GOVERSION=go1.26.6` and therefore isolate P1: **M10a′** TAB-indents the root directive (P1 reads
   `0` column-0 `go ` lines) and **M10b** duplicates it (`2` lines) — each with an rc=0 / marker=1 /
   2-races control confirming a sole killer.
3. **P1's `exit 2` (malformed-token) branch has NO live-script killer, and this doc claimed M9+M10
   "cover both of P1's branches with no overlap".** P1 has **three** refusal branches, not two, and
   the malformed one is unreachable through the script: every root-`go.mod` directive Go rejects also
   flips `GOVERSION` into the deny-list, and the `ACTIVE_GO` side cannot be forced on this rig. It is
   a **declared residual** whose only killer is V15/V16's standalone battery — stated, not hidden.
4. **AC7's sha256 clause was not computable** ("the post-comment hash is the base hash plus the
   comment block" — hashes do not compose). Replaced by a measured form: `git diff --unified=0 |
   grep '^-'` → **0** deleted lines, `grep -n '^go 1.22$'` → the directive's line, col-0 `go ` count
   → **1**.
5. **`+~8 LOC net` for `verify_go.sh` is ~4× low** — the verbatim block is **35** lines plus the one
   altered line at `:229`.
6. **AC9/V10 misreported a grep**: `grep -n 'go test ./\.\.\.' scripts/verify_go.sh` yields **3**
   matches (a comment at `:5`, an `echo` at `:256`, the executable line at `:258`), not one.
   Corrected to `grep -cE '^go test \./\.\.\. -count=1$'` → **1**.
7. **AC2 cited V13 for M9/M10**; V13 is the racecontrol ARM row and contains no P1 measurement. The
   correct citations are V14 and V16.
8. **M8 specified `git mv`**, a git write the sandboxed executor is forbidden to make. `mv` measured
   equivalent.
9. **The warm-cache framing is right and the cache is now named**: `$(go env GOMODCACHE)/golang.org`
   holds `toolchain@…go1.{24.9,25.6,26.0,26.2,26.3,26.5,26.6,27.0}`, which is why even a `go 1.27.0`
   floor still fires 2 races under `auto`. **Every M1 runtime reading must therefore be taken under
   `GOTOOLCHAIN="$ACTIVE_GO"`, never `auto`** — under `auto` the arm measures the cache, not the floor.

**AND THE ONE THAT CHANGED THE SPRINT'S SCOPE — M11: DELETING THE ENTIRE P1 BLOCK IS A SILENT GREEN
IN BOTH LANES.** Measured: remove the whole block and the static gate is `rc=0 RUN=1 PASS=1` and the
runtime lane is `rc=0 marker=1 races=2`. The block the attended ruling `D-WORLD-28` exists to mandate
had nothing protecting it, while P2's execution binding already had **M7** for exactly that reason
("the machinery that binds execution must not silently revert"). This doc's own Conflict Surface
already promised "a P1 presence needle that the failure-open deletion reds", and no AC, component
spec or mutant implemented it. **Controller decision D1: IMPLEMENT the needle** — the asymmetry
between a protected P2 and an unprotected P1 has no argument behind it, and withdrawing the Conflict
Surface's claim would be more work than honouring it. The needle is constrained away from this
mission's own row-49/row-59 defect (a `strings.Count` of a literal cannot see shape-gutting): it
binds the block's load-bearing SEMANTICS — the per-branch FATAL message texts, the `root_go_lines`
count guard, the comparator's three-way exit contract — never a comment or other decoration, and it
carries partial-gutting arms as well as **M12**'s whole-block deletion, because whole-block deletion
is the easy case and single-branch removal is what a careless future edit actually produces.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|----------|-----------------|-----------|----------|-------------|
| Anchor the invariant to the **root module floor** (`go1.26.6`, read from the root `go.mod` `go` line via `moduleGoFloor`), not the oldest `KNOWN_GOOD` and not a literal | `ACTIVE_GO` (the toolchain that runs the control after the execution-binding edit) is always `>= root floor` under `GOTOOLCHAIN=auto`; `KNOWN_GOOD`'s oldest `go1.24.9` never touches `verify_go.sh:229`'s ambient toolchain (round-1 R1). The root floor is the repo's own statement of its pinned toolchain — already the anchor of `TestGoToolchainPinsAgreeAndMatchJobList` and `PINNED` — so the bound derives from the checked value, not an authored list | agent (round-1-block-settled) | design | low |
| Bind as `<=` (floor at-or-below root toolchain), not `<` | Equality still arms the control (measured: floor `go1.26.6` under `GOTOOLCHAIN=go1.26.6` fires, V6); strict order over-constrains and adds no safety | agent (evidence-settled) | design | low |
| **Execution-binding production edit to `verify_go.sh`**: race leg runs `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` | This is the machinery that actually binds execution (round-1 R1). Without it, nested auto-selection inside `racecontrol/` can still pick a toolchain different from `ACTIVE_GO` (the `go1.26.4` base, V3), and the floor-to-root bound would be vacuous w.r.t. execution; with it, the control runs under the deny-list-vetted root toolchain and the static floor bound is meaningful | agent (round-1-block-settled) | design | low |
| **Not** option B (pin `GOTOOLCHAIN=local`) | `GOTOOLCHAIN=local` removes the cache rescue but does **not** prevent the disarm (floor > bound toolchain still fails loudly); it is a concealment fix, not a prevention | agent | design | low |
| New test lives in `host/verifygate/toolchain_pin_gate_test.go`, reusing `moduleGoFloor`/`shellAssignmentValues`/`repoRoot` | Same home and helpers as row 42's Test A; the file already owns run.sh/floor binding; `go/version` already imported | agent | design | low |
| **One-line production edit to `verify_go.sh`** is in scope (unlike the round-0 draft's "no production edit"); **expanded to P1+P2 by D-WORLD-28** — see the P1 row below and Design Freeze | The round-1 DIRECTIONAL block makes the execution-binding edit the crux of the fix; the deny-list (lines 217–224) and the FATAL/arming logic (lines 226–236) are untouched | agent (round-1-block-settled); scope expanded by human (D-WORLD-28) | design | low |
| **P1: runtime fail-closed in `verify_go.sh` — the selected `ACTIVE_GO` must be at-or-above the root module floor, else FATAL + exit 1** | Mandated by D-WORLD-28 (2026-09-01); it makes the static proof's middle premise true by enforcement rather than assumption, closing R2-c. The floor is **read from `./go.mod`**, never hardcoded (this repo's control-derived-from-the-checked-value rule; a literal would go stale on the next floor bump). It sits after the deny-list (`:218–224`) and before the race leg so every below-floor base — including a `go1.25.x` base the deny-list does not catch — reds on the floor check itself | human (attended ruling D-WORLD-28) | design | low |

## Design Freeze

- **One test file modified, no test file created**: the new test appends to
  `host/verifygate/toolchain_pin_gate_test.go`, reusing `moduleGoFloor` and `repoRoot`;
  `go/version` is already imported (V9).
- **Two-part production edit to `verify_go.sh`, no longer a single line** (D-WORLD-28): (P1) a
  runtime fail-closed block between the deny-list (`:218–224`, i.e. right after `esac`) and the
  race leg: after `ACTIVE_GO` is captured at `:217` and the deny-list exits, read the root
  module floor from `./go.mod` and FATAL + exit 1 unless `ACTIVE_GO` is at-or-above it; (P2)
  the race leg at `:229` changes `go run -race .` → `GOTOOLCHAIN="$ACTIVE_GO" go run -race .`.
  `ACTIVE_GO` is already captured at `:217`, so no new variable; the deny-list (`:218–224`) and
  the arming/FATAL logic (`:226–236`) are otherwise untouched. The floor literal is never
  hardcoded; it is read from the root `go.mod` `go` line each run.
- **Comment-only production edit**: `racecontrol/go.mod` gains a `LOAD-BEARING` comment block
  above the directive; **the `go 1.22` line itself is not moved**.
- **NOT touched**: `ci.yml`, `run.sh`, `repro/`, the root `go.mod`, any `host/store` assertion.
- **Anchor is derived, not authored**: the invariant reads the root `go.mod` floor via
  `moduleGoFloor` and the `racecontrol` floor via `moduleGoFloor`, and the test only compares
  them (row 43's control-derived-from-checked-value refutation, binding here as in row 42).
- The invariant does not call `requirePinned` or read `AILANG_BIN`: a static text scan must
  run in any lane (row 42's P9 split).
- The `racecontrol/go.mod` comment must not collide with any needle the new test counts. The
  test's only comment-adjacent needle — the execution-binding substring
  `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` — is scanned in `verify_go.sh`, a different file,
  and the fence comment spells the channel as `GOTOOLCHAIN=$ACTIVE_GO` (no quotes, no `go run
  -race .` suffix), so no token the test scans collides.

## Deferred Decisions

- **Option B (pin `GOTOOLCHAIN=local`)** — **rejected as the fix by D-WORLD-28 (2026-09-01); it
  stays a deferred note only** (this doc never brings it back as a resolution). Measured
  reasoning: it removes the *cache rescue* that conceals the disarm but leaves the *disarm*
  itself (floor > bound toolchain still reds `go run -race .`); it is a loudness improvement,
  not the prevention the ruling mandates (P1). If a future item explores it, it is a one-line
  edit to `verify_go.sh`'s race leg and belongs in row 44's lane, not here.
- **An explicit repo-wide "oldest supported toolchain" policy** — would give future items a
  single maintained constant instead of deriving the bound from the root floor. Today the root
  floor is the repository's own statement of the pinned toolchain; if the repo ever adds an
  explicit oldest-toolchain policy, re-anchoring is a three-line change.

## Solution Design

### Overview

The disarming is a *static* condition — a floor directive above the toolchain that actually
runs the control — and round 2 established the runtime lane can resolve that toolchain *below*
the root floor under `GOTOOLCHAIN=local` (R2-c). The fix therefore has **three** parts, which
D-WORLD-28 composes into one proof:

1. **P1 — runtime fail-closed (new, this revision)**: `verify_go.sh`, after `ACTIVE_GO` is
   captured at `:217` and the miscompile deny-list exits at `:224`, reads the **root module
   floor from `./go.mod`** and FATALs + exit 1 unless the selected `ACTIVE_GO` is at-or-above
   that floor. This replaces the round-1 proof's *assumption* `ACTIVE_GO >= root floor` with an
   *enforced* premise — it holds even under `GOTOOLCHAIN=local` with a base below the floor
   (V12). The floor is read, never hardcoded, so it cannot go stale when the root floor moves.
2. **P2 — execution binding**: `verify_go.sh:229` runs the control under
   `GOTOOLCHAIN="$ACTIVE_GO"`, where `ACTIVE_GO = go env GOVERSION` in the root module. This
   forces the race control to execute under the same toolchain the deny-list and P1 just
   vetted, removing the nested auto-selection that let an unrelated ambient base (`go1.26.4`)
   change which toolchain ran the control (V3).
3. **P3 — static gate (test)**: `TestRaceControlFloorStaysBelowRootToolchain` in
   `host/verifygate` asserts `racecontrol floor <= root module floor`, and pins the
   execution-binding needle so the round-1 hole cannot silently return.

The composed proof is the chain **floor(racecontrol) <= floor(root) <= ACTIVE_GO**: the first
`<=` is supplied by P3 (static); the second `<=` is supplied by P1 (enforced at runtime); and
P2 makes `ACTIVE_GO` the toolchain that actually runs the control, so the floor is satisfiable
by the very toolchain that executes it. This answers the round-1 DIRECTIONAL objection with
machinery and the round-2 direction dispute with an enforced premise.

### Architecture

The repository maintains exactly **three** `go.mod` modules (census, V1): the **root** (`go
1.26.6`, the protection), **`repro`** (`go 1.22`, the miscompile reproducer, bound by row 42's
Test A to oldest `KNOWN_BAD`), and **`racecontrol`** (`go 1.22`, the race-detector
known-positive, currently *unbound* — this item). The race control is driven only by
`verify_go.sh:229` (V4, the only reference anywhere). `verify_go.sh` captures `ACTIVE_GO =
$(go env GOVERSION)` in the root module (`:217`), refuses to proceed under the deny-listed
toolchains `{go1.26.0 … go1.26.5}` (`:218–224`), then — after this item — runs the control as
`GOTOOLCHAIN="$ACTIVE_GO" go run -race .` from the `racecontrol/` directory (`:229`).

The execution chain is now a single, statically knowable toolchain: **`ACTIVE_GO`**, which is
**enforced** to be `>= root floor` (`go1.26.6`). P1 reads the root module floor from `./go.mod`
and fails closed unless `ACTIVE_GO >= that floor`; P2 then binds the `racecontrol` execution to
that exact `ACTIVE_GO`; P3 statically asserts `racecontrol floor <= root floor`. Therefore
**`racecontrol floor <= root floor <= ACTIVE_GO`** — the floor is satisfiable by the very
toolchain that runs the control, for every possible host base, including under
`GOTOOLCHAIN=local` with a base below the root floor (V12). That chain is the evidence-settled
answer to round-1 R1 and round-2 R2-c: the invariant is anchored to the toolchain the control
is *made to run under*, and the run's lower bound is enforced, not assumed. Note that the
miscompile deny-list (`:218–224`, covering `go1.26.0..go1.26.5`) is a *miscompile* list, not a
floor; on this rig it happens to mask the refuting scenario because the local base is `go1.26.4`
(already denied), but a base below the root floor but outside the deny-list (e.g. `go1.25.x`)
passes it and still disarms the control — P1 is what actually carries the floor (V12, M9).

The test is byte-for-byte the shape of row 42's `TestReproModuleFloorStaysBelowKnownBad
Toolchains` (same `moduleGoFloor`/`version.Compare` machinery), with a different anchor because
the two nested modules have different jobs and because round 1 rejected the `KNOWN_GOOD` anchor.

**Components:**
1. **`verify_go.sh` P1 block (after `:224`) + `:229` one-line edit** — a runtime floor check
   that FATALs + exit 1 unless the selected `ACTIVE_GO` is at-or-above the root module floor
   read from `./go.mod`, then
   `race_control_output="$(cd
   design_docs/verification/w-race-gate-blindspot/racecontrol && GOTOOLCHAIN="$ACTIVE_GO" go
   run -race . 2>&1)"`. The control's execution is bound to the root toolchain.
2. **`TestRaceControlFloorStaysBelowRootToolchain`** (`host/verifygate`) — the static gate.
   Reads `racecontrol/go.mod`'s floor via `moduleGoFloor`; reads the root `go.mod`'s floor via
   `moduleGoFloor`; validates both with `version.IsValid`; asserts
   `version.Compare(racecontrolFloor, rootFloor) <= 0`, else `t.Fatalf` naming the exact
   consequence. Then reads `scripts/verify_go.sh` and asserts the execution-binding needle
   `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` occurs exactly once (the round-1 hole cannot
   return by reverting part 1).
3. **`racecontrol/go.mod` `LOAD-BEARING` comment** — the IDE-bump human's first sight, the
   same fence row 42 added to `repro/go.mod`, reworded for the root-toolchain anchor.

### Implementation Plan

**Phase 1: Static gate** (~0.5h)
- [ ] Append `TestRaceControlFloorStaysBelowRootToolchain` to
      `host/verifygate/toolchain_pin_gate_test.go` (reuses `moduleGoFloor` and `repoRoot`;
      `go/version` already imported).
- [ ] Run the AC1 run-existence form, the named mutations (M1–M10), and the pristine control.

**Phase 2: Runtime (P1 + P2) + fence** (~0.35h)
- [ ] **P1**: add a runtime fail-closed block to `scripts/verify_go.sh` between the deny-list
      `esac` (`:224`) and the race leg (`:229`): read the root module floor from `./go.mod`,
      and FATAL + exit 1 unless the selected `ACTIVE_GO` is at-or-above it. Confirm the
      deny-list and FATAL logic are untouched.
- [ ] **P2**: edit `scripts/verify_go.sh:229`: `go run -race .` → `GOTOOLCHAIN="$ACTIVE_GO" go
      run -race .`.
- [ ] Add the `LOAD-BEARING` comment block to `racecontrol/go.mod`; confirm the `go 1.22`
      directive is byte-unchanged (AC7).

**P1 comparison choice (designed here, stated for the implementer):** go version strings are not
lexically ordered in general (`go1.9` vs `go1.10`), so P1 compares numerically, not as text. The
doc's `go/version` `Compare`/`IsValid` machinery (V7) compares the numeric release components
(`1`, `26`, `6`) in order — a total order over well-formed `go1.X.Y` tokens. This is total over
every version this repo can actually see: the root floor, every `ACTIVE_GO` the toolchain
resolves, and every nested-module floor are release versions `go1.(major).(minor)` parsed by the
same `go/version` semantics V7 exercised. The claim is restricted to that set; a pre-release or
suffix-bearing token (`go1.26.6rc1`) is rejected by `version.IsValid` as an instrument floor, not
silently compared (V7: `IsValid` semantics).

**P1 runs in `bash`, so V7 specifies the ORDER and does not implement it — see V15**, a controller-run 13-case battery showing a portable `awk` component-compare reproduces every one of V7's verdicts (malformed tokens exiting on their own code, which is M10's branch), while the naive shell form `[[ "go1.9" < "go1.10" ]]` is measured FALSE and fails open. The implementer takes the order from V7 and the mechanism from V15; `sort -V` is ruled out because this rig ships BSD `sort` and CI runs GNU coreutils.

**The exact P1 block (round-3 `gpt5-6-sol`: "add the exact proposed P1 shell block to the
design"). Inserted verbatim after the deny-list `esac` at `:224`; measured as written in V14/V16.**

```bash

# --- P1 (queue row 48 / D-WORLD-28): the SELECTED toolchain must be at-or-above the
# root module floor. Below it, the nested race-control module's own floor can exceed
# the toolchain that runs it and the known-positive control at :229 goes silently
# unarmed. The floor is READ from ./go.mod, never hardcoded.
go_version_ge() {   # rc 0 = $1 >= $2 ; rc 1 = $1 < $2 ; rc 2 = a token is not a release version
  awk -v a="$1" -v b="$2" 'BEGIN{
    if (a !~ /^go1\.[0-9]+(\.[0-9]+)?$/ || b !~ /^go1\.[0-9]+(\.[0-9]+)?$/) exit 2
    sub(/^go/,"",a); sub(/^go/,"",b)
    na=split(a,A,"."); nb=split(b,B,".")
    n=(na>nb?na:nb)
    for(i=1;i<=n;i++){x=(i<=na?A[i]+0:0); y=(i<=nb?B[i]+0:0)
      if(x>y) exit 0; if(x<y) exit 1}
    exit 0}'
}
root_go_lines=$(awk '/^go /{n++} END{print n+0}' go.mod)
if [ "$root_go_lines" -ne 1 ]; then
  echo "verify_go.sh: FATAL: root go.mod has $root_go_lines column-0 'go ' lines, want exactly 1;" >&2
  echo "  the root module floor cannot be read, so the race-detector control cannot be bounded." >&2
  exit 1
fi
ROOT_FLOOR="go$(awk '/^go /{print $2; exit}' go.mod)"
set +e
go_version_ge "$ACTIVE_GO" "$ROOT_FLOOR"; floor_rc=$?
set -e
case "$floor_rc" in
  0) ;;
  1) echo "verify_go.sh: FATAL: active toolchain $ACTIVE_GO is BELOW the root module floor $ROOT_FLOOR;" >&2
     echo "  the race-detector known-positive control would be disarmed. Pin GOTOOLCHAIN to $ROOT_FLOOR or above." >&2
     exit 1 ;;
  *) echo "verify_go.sh: FATAL: cannot order toolchain tokens (ACTIVE_GO=$ACTIVE_GO ROOT_FLOOR=$ROOT_FLOOR);" >&2
     echo "  at least one is not a well-formed goX.Y[.Z] release version." >&2
     exit 1 ;;
esac
echo "   ✓ toolchain floor gate: $ACTIVE_GO >= root module floor $ROOT_FLOOR"
```

Three properties of that block are load-bearing and each has its own arm. **(i)** the floor is
READ from `./go.mod` and normalised with a `go` prefix, because the file's token is `1.26.6` and
the comparator's grammar is `go1.X[.Y]` (V18). **(ii)** the count check is a separate, earlier
refusal from the ordering check, so a missing/duplicate floor is attributed to the *instrument*
and a below-floor toolchain to the *input* (V16). **(iii)** the comparison is `awk`, not `[[ < ]]`
and not `sort -V` — the lexical form is measured to FAIL OPEN (V18) and `sort -V` would be two
implementations across BSD and GNU (V15).

**Phase 3: Hygiene + docs** (~0.1h)
- [ ] `go vet ./host/verifygate/`, `gofmt -l host/verifygate/`; confirm no test-name collision;
      note the row-48 linkage in `world-mission.md` if the iteration record needs it (doc-only).

## Files to Modify/Create

**Modified files:**
- `scripts/verify_go.sh` (**+35 / −0 plus 1 altered line** — the earlier `+~8 LOC net` estimate is
  ~4× low and is withdrawn: the verbatim P1 block is **35** lines, landing at `:225–:259`, plus **P2**
  1 line altered at the old `:229`, which becomes `:264`. The executor's diff review must expect
  exactly that shape) — P1 reads the root module floor
  from `./go.mod` and FATALs + exit 1 unless `ACTIVE_GO` is at-or-above it; P2 binds the race
  leg's invocation's toolchain. Deny-list (`:218–224`) and arming/FATAL logic (`:226–236`)
  unchanged.
- `host/verifygate/toolchain_pin_gate_test.go` (+~50 LOC for the bound/needle clauses, plus the
  **AC12 P1a–P1f semantic needles**) — the new `TestRaceControlFloorStaysBelowRootToolchain`; no new
  imports, no name collisions (V9).
- `design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod` (**+8 comment lines, −0**)
  — `LOAD-BEARING` fence; the `go 1.22` directive byte-unchanged, asserted by AC7's diff-shape form
  (0 deleted/modified lines; directive still at line **14**; column-0 `go ` count **1**) rather than
  by the withdrawn sha256-composition clause. Base sha256 `ab782f11…` (V1).

**NOT touched** (unchanged from the round-0/round-1 scope): `ci.yml` (the existing Design Freeze
"ci.yml is NOT touched" stays true), `run.sh`, `repro/`, the root `go.mod`, any `host/store`
assertion — the new `verify_go.sh` changes are confined to P1/P2 in that one script.

**New files:** none.

No other files. `ci.yml`, `run.sh`, `repro/`, the root `go.mod`, any `host/store` assertion —
untouched.

## Examples

### Example 1: the disarming a floor bump causes, and what the static gate makes it red instead

**The attack (measured, V2 — the disarm under the execution-bound runtime):**

```
$ ACTIVE_GO=$(go env GOVERSION)          # root module: go1.26.6 (V3)
$ (cd design_docs/verification/w-race-gate-blindspot/racecontrol \
     && GOTOOLCHAIN="$ACTIVE_GO" go run -race .)
   # rc=1, 2× WARNING: DATA RACE            <- bound to ACTIVE_GO, the known-positive fires
$ sed -i '' 's/^go 1\.22$/go 1.27.0/' go.mod
$ GOTOOLCHAIN=go1.26.6 go run -race .
   # rc=1, output: go: go.mod requires go >= 1.27.0 (running go 1.26.6; GOTOOLCHAIN=go1.26.6)
   # 0× WARNING: DATA RACE                  <- verify_go.sh:232 FATALs, control is dead
```

**After this item**, the same edit reds the *static* gate in CI (the mutation M1), before any
`go run -race .`:

```
$ GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1
   # racecontrol module floor "go1.27.0" is above the root module floor "go1.26.6": ...
   #   --- FAIL
```

And the runtime lane fails closed on **P1** when the selected `ACTIVE_GO` is below the root
floor, before the race leg runs (V12, V13):

```
$ GOTOOLCHAIN=local ./scripts/verify_go.sh   # base go1.26.4 < root floor go1.26.6
   verify_go.sh: FATAL: ACTIVE_GO go1.26.4 is below the root module floor go1.26.6 ...
   exit 1    <- the race leg never runs; no spurious "race detector not armed" blame
```

### Example 2: the new test, in the shape it is designed to guard

```
racecontrol/go.mod (comment-only edit):
// Deliberately a SEPARATE, NESTED module: `go build ./...` / `go test ./...` in the
// repository root module must NOT pick this race-detector control up. It is a
// diagnostic artifact, not host code, and it deliberately contains a data race.
//
// The `go 1.22` line below is LOAD-BEARING: it must stay at or below the repository
// root module floor (`go 1.26.6`, the `go` line of the root go.mod), because
// scripts/verify_go.sh runs this control under GOTOOLCHAIN=$ACTIVE_GO derived from
// the root module's GOVERSION; a floor above that toolchain makes `go run -race .`
// refuse before it can fire (verify_go.sh then FATALs "the race detector is not
// armed"). Enforced by TestRaceControlFloorStaysBelowRootToolchain (host/verifygate).
// Do not let `go mod tidy -go=…` or an IDE floor-bump touch it.
module ailang-world/verification/race-detector-control

go 1.22
```

## Success Criteria / Acceptance Criteria

Each AC carries its vacuity self-test and its **observed result on the unmodified tree at
`f1b5d1a`** (Verification Log rows cited).

- **AC1 — the test exists, RUNS, and passes on the post-sprint tree, in run-existence form.**
  `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run
  '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1 -v` → rc=0 with exactly 1 top-level
  `=== RUN` line and 1 `--- PASS`; a paired nonsense pattern (`-run
  TestNoSuchRaceControlFloorTestZZZ`) prints `[no tests to run]`, proving the instrument says
  so rather than passing vacuously. **Base @f1b5d1a: the verbatim named command → `no tests to
  run`, rc=0 (V9)** — the naive form greens at base measuring nothing; the `=== RUN` clauses
  are the repair and red at base (0 of 1).
- **AC2 — the pristine tree is green by construction and the named RED mutations red it.**
  Base verdict on the prototype logic = `PASS: racecontrol floor go1.22 <= root floor
  go1.26.6; execution needle count=1` (V5); in run-existence form the landed test records
  exactly 1 `=== RUN` and 1 `--- PASS` on the post-sprint tree; each of M1, M3–M8 reds on its
  named branch (V5); M2 is the equality GREEN control (V6); **M6′** (not M6) reds the
  root-floor-validity branch, and **M10a′/M10b** (not M10) red P1's instrument floor —
  see the Non-Vacuity table for why the doc's original M6 and M10 are not arms; M9 and
  M10a′/M10b red P1's refusal and instrument floors in the runtime lane (**V14** for M9,
  **V16** for the branch battery — the earlier V13 citation was wrong: V13 is the
  `racecontrol` ARM row and contains no P1 measurement); **M12–M18 red the six P1
  presence needles** with M17 as the green anti-brittleness control (AC12); every arm
  restored sha256-byte-identical, porcelain 0 (V1).
- **AC3 — M1 is the disarming edit and reds the static gate.** `racecontrol/go.mod` `go 1.22`
  → `go 1.27.0` (above the root toolchain `go1.26.6`) → `TestRaceControlFloorStaysBelowRoot
  Toolchain` FAILs on the bound clause (prototype verdict V5: `RACE-BOUND-FAIL`). The same edit
  *without* the gate is the disarm under the execution-bound runtime (V2: `0` DATA RACE,
  `verify_go.sh:232` would FATAL).
- **AC4 — M2 is the equality GREEN control with a runtime proof.** `racecontrol/go.mod` `go
  1.22` → `go 1.26.6` (exactly the root module floor): the gate stays GREEN (V5) **and** the
  control still fires at equality — `GOTOOLCHAIN=go1.26.6 go run -race .` from `racecontrol/`
  → rc=1, 2× `WARNING: DATA RACE` (V6). Had the runtime proof failed, `<=` would be the wrong
  bound and the strict form would return with that measurement as its justification; it did not
  fail.
- **AC5 — environment independence.** `env -u AILANG_BIN GOTOOLCHAIN=go1.26.6 go test
  ./host/verifygate/ -run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1` → rc=0
  with AC1's run-existence clauses; no network, no solver, no pinned binary (static text scan,
  V9's lane; the sibling static-lane run is env-independent at base, V9).
- **AC6 — hygiene.** `go vet ./host/verifygate/` rc=0 and `gofmt -l host/verifygate/` prints
  nothing. **Base: both green (V9)** — green-at-base by design; they measure the sprint's own
  edit.
- **AC7 — the fence comment landed, the directive did not move, and the post-comment floor
  still parses.** `grep -c "LOAD-BEARING"`
  `design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod` → 1 AND `grep -c '^go
  1.22$' …/racecontrol/go.mod` → 1, at line **14** AND
  `git diff --unified=0 -- …/racecontrol/go.mod | grep '^-' | grep -v '^---' | wc -l` → **0**
  AND `awk '/^go /{n++} END{print n+0}' …/racecontrol/go.mod` → **1**.
  **The sha256 clause this AC used to carry is WITHDRAWN as not computable:** it read *"the
  post-comment hash is the base hash plus the comment block"*, and cryptographic hashes do not
  compose — there is no operation turning `ab782f11…` "plus a comment block" into a checkable
  value, so the AC as written could not be executed. The three clauses above are executable and
  jointly assert the same thing, and strictly more precisely: *only `//` lines were added, and
  the `go 1.22` directive did not move or change*. Measured on the fenced file (V19). **Post-comment acceptance arm (round-1 R2):** `moduleGoFloor` on the commented file
  still returns exactly one valid floor (`floors=[go1.22] count=1`, V8) — the `//` comment
  line is ignored by `strings.HasPrefix(line, "go ")` and cannot break the exactly-one-floor
  refusal. **Base: `LOAD-BEARING` count = 0 (V8)** — the fence is AC7-red at base by design.
- **AC8 — systemic closure: all three `go.mod` floors are bound.** Census = exactly 3 modules
  (V1); root bound by `TestGoToolchainPinsAgreeAndMatchJobList`/`TestMiscompileInstrument
  ProbesPinnedToolchain`, `repro` bound by row 42's Test A, `racecontrol` bound by this item
  (V9: all existing pin tests green at base).
- **AC9 — the new test's CI linkage is proven mechanically (R2-a).** `.github/workflows/ci.yml`
  runs `./scripts/verify_go.sh`, which runs `go test ./... -count=1`, and `go list ./...`
  includes `host/verifygate`; a test added there therefore executes in CI. Three greps, each
  with its own control: `grep -n 'run: ./scripts/verify_go.sh' .github/workflows/ci.yml` → ≥1
  match at `:166` (control: step name at `:163` is "go build + test gate (replay tests run
  against pinned AILANG_BIN)"); `grep -cE '^go test \./\.\.\. -count=1$' scripts/verify_go.sh` → exactly **1** (at `:258`
  pre-sprint, `:293` post-sprint — **the AC binds the count, never the line number**, because P1
  adds 35 lines above it). The loose `grep -n 'go test ./\.\.\.'` form this AC used to carry is
  satisfiable by a comment alone and is used only as the read-control. **Correction, measured
  post-sprint at iteration 145: that loose form returns 4, not 3** — the planner's §4.5(b)
  enumeration listed `:5` (header comment), `:256` (`echo`) and `:258` (the executable line) and
  **missed `:262`, the `-race` leg's own `echo`**; the base tree returns 4 as well, so this is a
  miscount in the refutation and not a change the sprint introduced. Control: the loose form
  returns 4 at base and 4 post-sprint, proving the file was read; `go list ./... | grep verifygate` →
  `host/verifygate` (control: the package resolves). **Base @HEAD 8bb9214: all three readings
  recorded (V10).** RED: a committed edit that repoints the CI job elsewhere, drops the `./...`
  leg, or renames the package must make the corresponding grep red while its paired control
  stays green, so the failure is attributed to the linkage and not to a stale path.
- **AC10 — no `-run` selector in this doc is vacuous; the escaped-`\|` failure class is named
  (R2-b).** The rule: a go-test AC's binding evidence is a counted number of top-level `=== RUN`
  and `--- PASS` lines (run-existence), never a bare rc=0 — because a grep-style escaped
  alternation `\|` is a *literal pipe* to Go's regexp, so `go test -run 'A\|B'` matches nothing
  and exits 0 printing `ok … [no tests to run]`. Every `-run` selector in this doc is a valid
  unescaped pattern and matches its asserted test count at base (V11). A deliberately
  nonsense selector (`-run TestNoSuchRaceControlFloorTestZZZ`) prints `[no tests to run]`, and
  the harness must treat that as **RED** — vacuous green is exactly the class R2-b blocked.
- **AC11 — P1 fails closed: the selected `ACTIVE_GO` must be at-or-above the root module floor,
  else `verify_go.sh` exits non-zero before the race leg runs (D-WORLD-28).** **MEASURED, not
  asserted** — the round-3 objections were that this AC cited a plan. Evidence is now V14 (four
  runtime arms against a temp patched copy: GREEN at equality with the race leg reached and the
  control armed; M9 RED below the floor with the deny-list NEUTERED and the race leg unreached;
  the byte-identical no-P1 CONTROL reaching the race leg at rc=0, which is what makes M9's red
  attributable to P1 alone; and the deny-list-live arm showing the two guards' messages are
  distinguishable), V16 (one arm per refusal branch and per instrument floor, including the
  `go1.25.6` case the deny-list structurally cannot catch), V18 (the `go.mod` extraction and the
  lexical fail-open control) and V15 (the comparator battery). The sprint re-runs every one of
  these arms against the LANDED block rather than a prototype. Base is green by construction:
  with the real `go.mod` and an ambient toolchain, `ACTIVE_GO == root floor` and the race leg
  runs (V2, V14/A1).
- **AC12 — the P1 block cannot be silently deleted or gutted: six narrow semantic needles in
  `TestRaceControlFloorStaysBelowRootToolchain`, each with its own named RED.** This AC exists
  because the Conflict Surface promised a "P1 presence needle" that appeared in no component, no
  AC and no mutant, and the hole is **measured, not argued**: with P2 and P3 intact and the whole
  P1 block deleted from `scripts/verify_go.sh`, the static gate returns rc=0 / `=== RUN`=1 /
  `--- PASS`=1 **and** the runtime lane returns rc=0 with the race leg reached and 2 `WARNING:
  DATA RACE` — deleting the block `D-WORLD-28` mandates was **green in both lanes** (**M11**,
  V20). The needle set is deliberately **not** a `strings.Count(src, "<literal>") == 1`, which is
  row 49's defect and row 59's open item; the block is first split into its comparator half
  (`go_version_ge`'s body) and its gate half so an assertion about one cannot be satisfied by
  text in the other, and the six assertions are:
  **P1a** sentinel comment and `✓` success line occur exactly once each, sentinel first (presence
  + delimitation) · **P1b** the gate half contains, once each, the guard `[ "$root_go_lines" -ne
  1 ]`, the stem `is BELOW the root module floor` and the stem `cannot order toolchain tokens`,
  and exactly **3** `verify_go.sh: FATAL:` refusals (the three distinct failure modes exist and
  are separately attributed) · **P1c** the gate half contains exactly **3** `exit 1` and **zero**
  `exit 0` (every refusal actually refuses) · **P1d** the comparator half contains all of `exit
  0`/`exit 1`/`exit 2`, and the gate half calls `go_version_ge "$ACTIVE_GO" "$ROOT_FLOOR"`
  exactly once **in that operand order** · **P1e** the block contains **zero** concrete
  `go1.<d>.<d>` literals and does contain `awk '/^go /{print $2; exit}' go.mod` (the floor is
  READ, never hardcoded) · **P1f** byte offsets satisfy `deny-list case < P1 sentinel < P2
  needle` (P1 runs after the deny-list and before the race control).
  Two of the six (P1c, P1e) are **negative** assertions and one (P1f) is **positional** — the
  shapes a token count structurally lacks. Named REDs M12–M18 and the M17 green control are in
  the Non-Vacuity table. **The finding that makes this AC non-optional:** under the *ambient*
  toolchain, which is the only condition CI ever runs, all six gutted variants are `rc=0`,
  race leg reached, 2 races — byte-for-byte indistinguishable from the correct block, so the
  runtime lane is blind to every one of them in CI; and for **M14** (operands swapped) and
  **M16** (below-floor `exit 1` → `exit 0`) even the deliberately hostile M9 arm goes green
  (`rc=0` race leg reached, and `rc=0` after printing its own FATAL, respectively). For those
  two the static needle is the **sole killer in the entire sprint** (V21).
Explicitly rejected as an AC: "the full verify gate is green" alone — `go build ./... && go
test ./...` must of course pass on the final tree, but a package-wide `ok` can print while the
named test never ran (V9's `no tests to run` is rc=0; V11's escaped-`\|` selector is the same
class in the *positive* direction); AC1/AC10's run-existence form is the binding version.

## Conflict Surface

This change adds a **runtime floor-check block (P1)** and **one executable invocation line
(P2)** to `verify_go.sh`, which *strengthen* the runtime lane rather than altering the deny-list
or arming logic; it adds no executable line to the parsers/typecheckers/codegen. The conflict
surface is entirely in the shared files the static text gate reads and the shared helpers it
reuses.

- **`verify_go.sh`** — now *edited* (not merely observed) by this item in two places: the **P1
  fail-closed block** inserted between the deny-list `esac` (`:224`) and the race leg, and the
  line at **`:229`**. `ACTIVE_GO` (`:217`) and the deny-list (`:218–224`) are read, not changed;
  the arming/FATAL block (`:226–236`) is unchanged. The new test scans the P2
  execution-binding substring `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` (count=1) **and the six
  P1 semantic needles P1a–P1f specified in AC12** (sentinel/success delimitation; the three
  distinct refusal branches and their `FATAL:` count; the `exit 1`×3 / `exit 0`×0
  refusal-actually-refuses pair; the comparator three-way contract plus the
  `go_version_ge "$ACTIVE_GO" "$ROOT_FLOOR"` operand order; the zero-hardcoded-version /
  `awk`-extraction pair; and the deny-list < P1 < race-leg ordering). **This bullet previously
  promised "a P1 presence needle" that was specified nowhere** — not in the Solution Design, not
  in any AC, not in the mutation table, and with no mutant; that gap is closed here and its cost
  is measured as M11 in AC12; row 44's surface (CI
  `continue-on-error`) is untouched.
- **`racecontrol/go.mod`** — read by this test (`moduleGoFloor`) and executed by
  `verify_go.sh:229`; written by nobody. A *moved* file (M8) reds this test's read floor and
  breaks `verify_go.sh` at runtime; the test names the path, so a move must move the fence too.
- **Root `go.mod`** — read by this test (`moduleGoFloor`, as the anchor) and by
  `TestGoToolchainPinsAgreeAndMatchJobList` (floor == ci.yml pin) and
  `TestMiscompileInstrumentProbesPinnedToolchain` (PINNED). Shared fate on a *deleted* or
  *malformed* `go` line (M5, M6); the sibling reds on its own clauses. A root-floor *raise* above
  `racecontrol`'s floor keeps this test green (the bound is `<=`), which is correct: the control
  is still runnable under the higher root toolchain.
- **`run.sh`'s `KNOWN_GOOD=` line** — **no longer read by this test** (round-1 R1 dropped the
  `KNOWN_GOOD` anchor). It remains read by `TestMiscompileInstrumentProbesPinnedToolchain`
  (reads `KNOWN_GOOD`, `KNOWN_BAD`, `PINNED`). This test's execution bound no longer shares fate
  with `KNOWN_GOOD` edits.
- **Shared helpers** (`moduleGoFloor`, `repoRoot`) — a helper edit reds every consumer at once;
  already this file's accepted pattern (row 42).
- **`repro/` and the root module's other machinery** — not read by this test beyond the root
  `go.mod` `go` line; existing bindings are unaffected (V9 confirms all existing pin tests green
  at base).

### Programs that MUST still work

1. `host/verifygate/toolchain_pin_gate_test.go` → `TestReproModuleFloorStaysBelowKnownBad
   Toolchains` (row 42's Test A) — still green at base (V9).
2. … → `TestMiscompileInstrumentProbesPinnedToolchain` — still green at base (V9); unaffected
   by the `KNOWN_GOOD` anchor being dropped (that test owns `KNOWN_GOOD` on its own).
3. `scripts/verify_go.sh` race leg (lines 226–235) — after the `:229` edit, `go run -race .`
   still fires 2× `WARNING: DATA RACE` from `racecontrol/` under `GOTOOLCHAIN="$ACTIVE_GO"`
   (V2), and the FATAL/arming block is unchanged (V4).
4. `design_docs/verification/w-race-gate-blindspot/racecontrol/main.go` — `go run -race .`
   still fires 2× `WARNING: DATA RACE` under `go1.26.4`, `go1.26.6`, and the equality floor
   (V2, V6).
5. `design_docs/verification/w-race-gate-blindspot/repro/go.mod` — unchanged; row 42's binding
   intact (V9).

### What deliberately changes

At runtime, exactly one invocation: the race control is now run under an explicit
`GOTOOLCHAIN="$ACTIVE_GO"` instead of nested auto-selection. The deliberate effect is that a
`racecontrol` floor bump above the root toolchain is now a **red static gate** instead of a
silent (cache-concealed) or wrong-reason loud (fresh-runner) disarm, and that the runtime lane
executes the control under a known, deny-list-vetted toolchain. `racecontrol/go.mod`'s directive
is byte-unchanged.

## Testing Strategy

**Unit / gate tests:**
- `TestRaceControlFloorStaysBelowRootToolchain` — the static gate, with run-existence form
  (AC1) and environment independence (AC5).

**Selector vacuity (R2-b):**
- Every `-run` selector in this doc is a valid unescaped pattern; the binding evidence for a
go-test run is a counted number of top-level `=== RUN` / `--- PASS` lines, never a bare rc=0.
A nonsense selector must print `[no tests to run]` and the harness must treat that as RED; the
escaped-`\|` case is why the rule exists (AC10, V11).

**Mutation testing (non-vacuity):**
- M1–M10 as enumerated below; one named RED mutant for every refusal and anti-vacuity floor;
  M2 as the equality GREEN control with a runtime arming proof (AC4); M9/M10 red P1's refusal
  and instrument floors in the runtime lane (AC11, V13).

**Regression-surface:**
- All five "Programs that MUST still work" entries re-run green at base (V9) and re-confirmed
  by the sprint.

**Manual / runtime lane (P1/P2):**
- `verify_go.sh` passes P1 on the real tree: with the actual `go.mod`, `ACTIVE_GO` is >= the
  root floor and the race leg runs (V13).
- `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` from `racecontrol/` → control fires (V2); at the
  equality floor `GOTOOLCHAIN=go1.26.6 go run -race .` → control fires (V6);
  `GOTOOLCHAIN=local` with a below-floor base fails closed on P1 before the race leg (V12,
  V13).

## Non-Goals

- **Not** pinning `GOTOOLCHAIN=local` in `verify_go.sh` — option B; it removes concealment but
  not the disarm, and is rejected by D-WORLD-28. [Deferred Decisions only, never to be the fix]
- **Not** altering `verify_go.sh`'s deny-list (`:218–224`) or its arming/FATAL logic
  (`:226–236`). The P1 fail-closed block sits between them and the race leg; the deny-list and
  arming/FATAL blocks are unchanged. [Out of scope]
- **Not** binding the floor to `run.sh`'s `KNOWN_GOOD` list in any form — round-1 R1 rejected
  that anchor; the bound is the root module floor, derived from the toolchain that actually
  runs the control. [Quorum-round-1-settled]
- **Not** rewriting the race control's runtime behaviour; `main.go` unchanged. [Out of scope]
- **Not** a runtime pre-flight *that is the gating lane* — the static floor<=floor test (P3)
  is what gates CI. A *runtime* floor check (P1) is still added, because D-WORLD-28 mandates
  failing closed when the selected `ACTIVE_GO` is below the root floor; the two complement
  (P3 reds in CI, P1 protects the runtime lane under `GOTOOLCHAIN=local`; R2-c). [D-WORLD-28]
- **Not** re-binding the `repro` module (already bound by row 42's Test A). [Already shipped]

## Timeline / Milestones

**Milestones:**
- **M1 — gate + runtime (P1 fail-closed + P2 binding) land**: `verify_go.sh` gains the P1
  root-floor fail-closed block after `:224` and the `GOTOOLCHAIN="$ACTIVE_GO"` binding at
  `:229`; `TestRaceControlFloorStaysBelowRootToolchain` added, AC1/AC5 green, M1–M10 red on
  named branches (AC2/AC3/AC4/AC11), pristine control green.
- **M2 — fence lands**: `racecontrol/go.mod` `LOAD-BEARING` comment added, directive
  byte-unchanged (AC7), post-comment floor parses (V8).
- **M3 — hygiene + audit closure**: vet/gofmt clean (AC6), no collisions (V9), CI linkage
  proven (AC9), no vacuous selector (AC10), systemic closure confirmed (AC8).

**Effort:** Phase 1 ~0.5h, Phase 2 ~0.35h, Phase 3 ~0.1h → ~1h total (matches the ~0.1d row
estimate; a single sprint, no split).

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| The anchor (root module floor = `go1.26.6`) is the lower bound on the toolchain that runs the control; if the selected `ACTIVE_GO` ever resolved *below* the root floor the static bound would overstate the guarantee | Low | **P1 enforces it**: `verify_go.sh` fails closed unless the selected `ACTIVE_GO` is at-or-above the root floor read from `./go.mod`, so `ACTIVE_GO >= floor(root)` holds under every host base, including `GOTOOLCHAIN=local` with a base below the floor (V12). P2 makes that `ACTIVE_GO` the control's actual runtime toolchain, so the bound is a proof, not an assumption (V13) |
| A base toolchain *above* the root floor (future `go1.30`) raises `ACTIVE_GO` above `go1.26.6` | None | The static bound is `racecontrol floor <= root floor`; since `ACTIVE_GO >= root floor`, the control still runs under any higher `ACTIVE_GO` — the bound is monotone-safe in the host-base direction |
| The static test reads `verify_go.sh` text, so an execution binding spelled differently (computed, indirection) escapes it | Low | The exact needle `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` is pinned at count=1 (M7); any re-spelling reds the count — the honest bound, matching row 42's declared residual |
| M5/M6's root-`go.mod` mutants could be mistaken for live root-module changes | Low | They are declared instrument-failure floors on the anchor; the root module's own tests red on the same edits, so they measure this test's validity/comparison arms, not a real drift |
| The `LOAD-BEARING` comment could be misread as colliding with a scanned needle | Low | The only comment-adjacent needle lives in `verify_go.sh` (a different file) and the fence spells the channel `GOTOOLCHAIN=$ACTIVE_GO` (no quotes, no `go run -race .`), so no collision (V8) |
| The P1 floor read from `./go.mod` could be misparsed (two `go ` lines, a deleted directive, an unreadable file), which would make the fail-closed check silently vacuous | Low | M10 reds the unreadable/unparsable floor branch in the runtime lane (V13), and M4 does the same for the static test's exactly-one-`go `-line floor; a read failure exits before `:229` with a named FATAL, never a silent pass |
| The P1 version comparison could misorder versions if it compared lexically (`go1.9` vs `go1.10`) | Low | P1 uses the numeric release-component order of `go/version` (V7) — total over the well-formed `go1.X.Y` tokens this repo can actually see; a malformed token is rejected as an instrument floor, not compared (AC11, V13) |
| A below-root-floor base *outside* the miscompile deny-list (e.g. `go1.25.x`) could pass the deny-list and still disarm the control | Low | That is precisely the class P1 closes: it is an independent floor check, not a miscompile list; M9 proves the refusal fires with the deny-list neutered (V13, AC11) — the doc states this explicitly so no future reader mistakes the deny-list for the floor (Architecture, V12) |

### Declared residual

The static gate (P3) cannot know the *live* `ACTIVE_GO` on a given runner; it binds the floor
chain to the root module floor and lets P1 enforce the runtime lower bound. The chain
`floor(racecontrol) <= floor(root) <= ACTIVE_GO` holds by construction: P3 proves the first `<=`
statically and P1 enforces the second at runtime — including under `GOTOOLCHAIN=local` with a
base below the root floor, which round 2 showed `GOTOOLCHAIN=auto` cannot guarantee on its own
(V12). The only true residual is the pre-existing boundary where the root module itself cannot
run at all — a host whose base Go is below the root floor with no network to fetch it, where P1
fail-closes loudly before the race leg; that loud FATAL stays the backstop and now names the
floor mismatch rather than a spurious race-detector failure.

**Widened 2026-09-05 (queue row 61, iteration 157).** This residual was written as if the static
needle set's blind spot were reachable only by *rewriting* an existing branch (`case "$floor_rc"
in` → `case "0" in`, or inverting the comparator's awk verdicts). It was reachable by a one-line
**INSERTION**, which is strictly larger and much likelier as an accident. Measured first-party at
HEAD before the fix: `floor_rc=0` inserted immediately after `go_version_ge "$ACTIVE_GO"
"$ROOT_FLOOR"; floor_rc=$?` is `bash -n` clean, `go vet` clean, leaves every one of P1a–P1f green
(`TestRaceControlFloorStaysBelowRootToolchain` `ok`), and opens **all three** refusal branches at
once — not only the below-floor branch the row named, but the malformed-token instrument floor
too, which then printed `✓ toolchain floor gate: devel >= root module floor go1.26.6`. The
success line is an assertion, and a fail-open gate publishes it as a false one.

The disposition is **structural**, not another needle: the comparator's verdict is now consumed
by `if` directly, and every arm of the attributing `case` exits, so a dataflow break can change
*which* refusal is reported but cannot produce a success. Measured on the new shape: an inserted
line in the `else` branch keeps `rc=1` (attribution shifts from the below-floor message to the
cannot-order message); an inserted line in the `then` branch changes nothing. Note the option
this rules OUT — dropping the variable and reading `case "$?"` directly is fail-**open even
unmutated**, because the intervening `set -e` is itself a successful command that resets `$?`
(measured: below-floor input, no mutation, `rc=0`). Two new conjuncts back the shape: the call
must be the `if` condition, and the gate must contain no `<var>=$?`.

**The residual that remains, stated in its widened form:** a static scan of `verify_go.sh` cannot
follow dataflow, so it bounds the *shape* of the verdict path and not its *meaning*. What is now
closed is every break that leaves the shape intact — insertion included. What is not closed is a
rewrite of the comparator's own body (inverting the `awk` exit codes), which changes the meaning
while satisfying every needle; that class is covered by the comparator's `exit 0`/`exit 1`/
`exit 2` contract needle and by V15's 13-case battery, not by the consumption-shape binding.

## Axiom Compliance

**Scoring**

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +1 | The gate is a deterministic static text scan; the same tree gives the same verdict in any lane |
| A2: Replayability | 0 | No replay/session surface |
| A3: Effect Legibility | 0 | No effect/IO change to AILANG programs; one invocation-line strengthening in the shell gate |
| A4: Explicit Authority | 0 | No capability change |
| A5: Bounded Verification | +1 | Replaces a runtime-only, cache/auto-selection-dependent signal with a bounded static gate in CI (P3) plus a runtime fail-closed floor check (P1, D-WORLD-28) that binds execution to a known toolchain above the root floor |
| A6: Safe Concurrency | 0 | No concurrency change (the race control's own race is untouched and still fires) |
| A7: Machines First | +1 | A static gate is cheap, deterministic, and runnable by CI without a warm toolchain cache; the runtime lane now pins the toolchain it executes and fails closed below the root floor (P1) |
| A8: Minimal Syntax | 0 | No syntax change |
| A9: Cost Visibility | 0 | No cost surface |
| A10: Composability | +1 | Reuses row 42's helpers and shape; composes with the existing three-module binding |
| A11: Structured Failure | +1 | The disarm is now a named, attributed RED in both lanes: the static gate names the floor mismatch (P3), and the runtime lane fail-closes on the floor check (P1) naming the mismatch rather than a wrong-reason "race detector not armed" FATAL |
| A12: System Boundary | 0 | No boundary change |

**Net Score: +5** ✅ Proceed to implementation

**Hard Violation Check**
- [x] A1 (Determinism): no implicit nondeterminism — the gate reads files and compares versions
- [x] A3 (Effects): no hidden side effects — comment-only production edit plus a runtime floor
      check and an invocation binding in the gate's own shell script
- [x] A4 (Authority): no ambient access granted
- [x] A7 (Machines First): the static gate serves CI and fresh runners; the runtime lane
      executes under an explicit toolchain and fails closed below the root floor (P1)

## Non-Vacuity — named RED mutation for every refusal and anti-vacuity floor

Production side mutated (`racecontrol/go.mod` floor, the root `go.mod` floor, `verify_go.sh`'s
runtime check + execution binding, file placement) — never the test helpers, and never the
production edit that is the fix. Every assertion maps to a mutant.
The instrument-failure floors row 42's shape carries — *exactly one `go ` line* (M4) and *every
floor token valid* (M3, M6) — are enumerated here as branches with their own mutants per the
iter-108 floor-pinning rule: the anti-vacuity floors are a refusal about the *measurement*, not
about the input, and they get one arm each. M7 is the round-1 R1 regression arm: reverting the
execution binding re-opens the exact hole the quorum blocked, so it is a named RED.

**M9 and M10a′/M10b are the round-2 / D-WORLD-28 arms for P1's refusal and instrument-floor
branches**, each killing its branch **solely**: M9 reds the new runtime **refusal** (the selected
`ACTIVE_GO` is below the root module floor read from `./go.mod`), and is constructed so the
miscompile deny-list cannot supply the red; M10a′/M10b red the new runtime **instrument floor**
(the root floor could not be read as exactly one column-0 `go ` line). They live in the runtime
lane (`verify_go.sh`); M5/M6′ are the separate test-lane (P3) anchors. Each requires the race leg
to be unreached — that they fail *before* `:229` is the point — **but unreached-ness alone is not
attribution**, which is why every runtime arm carries a **P1-REMOVED control** that must be
`rc=0` with the race leg reached; the doc's original M10 was refuted precisely by failing that
control (V14, V23). **M11–M18 are the AC12 P1-presence arms** and live in the *static* lane;
they exist because M11 measures the whole block deletable with every other assertion green.

| # | Exact edit | Expected RED (single test name) | Shape / branch |
|---|---|---|---|
| M1 | `racecontrol/go.mod` `go 1.22` → `go 1.27.0` (above the root floor `go1.26.6`) | bound clause: floor `go1.27.0` above root floor `go1.26.6` | **threat-shaped: the disarming edit under the execution-bound runtime (V2: 0 DATA RACE)** — the runtime lane garbles it, the static gate must name it |
| M2 | `racecontrol/go.mod` `go 1.22` → `go 1.26.6` (exactly equal to the root floor) | **GREEN** (equality control) | boundary-pair, not an arm; runtime arming proof at equality (V6) |
| M3 | `racecontrol/go.mod` `go 1.22` → `go banana` | floor-validity floor: `racecontrol floor "gobanana" invalid` | malformed-floor-shaped: a floor `version.Compare` would misorder is itself a disarm |
| M4 | delete the `go 1.22` line from `racecontrol/go.mod` | exactly-one-go-line floor: `moduleGoFloor` fatal `found 0 go lines, want 1` | anti-vacuity floor: a control with no floor has no bound |
| M5 | root `go.mod` `go 1.26.6` → `go 1.20` (below the `racecontrol` floor `go1.22`) | bound clause: `racecontrol floor go1.22` above root floor `go1.20` | anchor-direction anti-vacuity: proves the root-anchor comparison is enforced, not a tautology |
| ~~M6~~ **REFUTED — vacuous, runs zero tests** | root `go.mod` `go 1.26.6` → `go banana` | *claimed* `root floor "gobanana" invalid` | **NOT AN ARM.** Measured: `rc=1  === RUN=0  --- PASS=0  --- FAIL=0`, `go: errors parsing go.mod: go.mod:3: invalid go version 'banana'`. Go's module loader rejects the mutated root `go.mod` **before the test binary is built**, so the test never runs — and under this doc's own AC10 rule an arm with `=== RUN`=0 is not evidence. As written it leaves `!version.IsValid(rootFloor)` with **no killer**. (M3, the *racecontrol*-side `go banana`, is unaffected: `racecontrol/go.mod` is a nested module `go test ./host/verifygate/` never loads — measured, M3 runs and reds on its named branch.) Replaced by M6′ (V22) |
| **M6′** *(the repair)* | root `go.mod` `go 1.26.6` → **`go 1.26.6 // pin`** | root-floor-validity floor: `instrument failure: root module floor "go1.26.6 // pin" is not a valid Go version` (`toolchain_pin_gate_test.go:574`) | anchor-validity anti-vacuity, **now runnable**: Go accepts a trailing comment so the module loads and the test runs; `moduleGoFloor`'s `canonicalizeVersionPin` yields the token `"go1.26.6 // pin"`, which `version.IsValid` rejects. Measured `rc=1  === RUN=1  --- PASS=0  --- FAIL=1` (V22) |
| M7 | `verify_go.sh:229` `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` → `go run -race .` (revert the binding) | execution-binding needle count=0, want 1 | **round-1-R1 regression: re-opens the nested-auto-selection hole the quorum blocked** — the machinery that binds execution must not silently revert |
| M8 | `mv racecontrol/go.mod racecontrol/go.mod.moved` (or move the dir) — **`mv`, never `git mv`: the executor is forbidden from git writes, and the plain form is measured to give the same red** | read floor: `moduleGoFloor` `os.ReadFile` fatal — the test names its target by path | placement-shaped: a moved control must move the gate in the same edit |
| M9 | **P1 refusal (runtime lane)**: neuter the miscompile deny-list (`:218–224` case body) so `go1.26.4` sails past it, and force the selected `ACTIVE_GO` below the root floor; run `verify_go.sh` | the P1 floor check FATALs + exit 1 (e.g. `verify_go.sh: FATAL: ACTIVE_GO go1.26.4 is below the root module floor …`) and the race leg at `:229` is never reached | **P1-only refusal arm**: proves the runtime floor check is live **independently of the deny-list** — on this rig the only below-floor base locally present is the deny-listed `go1.26.4`, so the arm **neuters the deny-list** (chosen over a `go1.25.x` base) so P1 alone supplies the red; this is the exact point that a below-root-floor base outside the deny-list (`go1.25.x`) would otherwise pass the deny-list and disarm the control |
| ~~M10~~ **REFUTED — the deny-list supplies this red, not P1** | delete the `go 1.26.6` line from the root `go.mod` | *claimed* P1 floor-read FATAL before `:229` | **NOT AN ARM.** The root `go.mod` floor is *what pushes `GOTOOLCHAIN=auto` up to `go1.26.6`*; delete it and `go env GOVERSION` falls back to the local base `go1.26.4`, which the **pre-existing miscompile deny-list** already refuses. Measured: P1-present → `rc=1 marker=0 FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go's`; **P1-REMOVED control → byte-identical red**. `go banana` in the root `go.mod` (M10c) fails identically for the same reason. M10's "the race leg is unreached" clause is satisfied by both arms — which is exactly why unreached-ness alone is not attribution. Replaced by M10a′/M10b (V23) |
| **M10a′** *(the repair)* | TAB-indent the root `go 1.26.6` directive | `verify_go.sh: FATAL: root go.mod has 0 column-0 'go ' lines, want exactly 1;` before `:229` | **P1 instrument-floor arm, sole-killer proven**: `go env GOVERSION` stays `go1.26.6` so the deny-list passes; P1 reds on its own attributed message (rc=1, marker=0, races=0) and the **P1-REMOVED control is rc=0 / marker=1 / races=2** (V23) |
| **M10b** *(the repair)* | duplicate the root `go 1.26.6` line | `verify_go.sh: FATAL: root go.mod has 2 column-0 'go ' lines, want exactly 1;` before `:229` | same branch from the other side; same sole-killer control (rc=0 / marker=1 / races=2 without P1) (V23) |
| — | **P1's third refusal branch — malformed token (`exit 2`)** — is **not reachable through the live script by any root-`go.mod` mutation on this rig**, because every `go` directive Go itself rejects also drops `go env GOVERSION` to the deny-listed base (above). It is reachable only from the `ACTIVE_GO` side (`devel`, `go1.26.6rc1`), which cannot be installed here. Its killer is therefore the **standalone comparator battery (V16/V18)**, recorded as such rather than banked as a runtime-lane arm. | — | scope correction to *"P1 has two branches"*: it has **three** |
| **M11** | **delete the whole P1 block** from `scripts/verify_go.sh`, leaving P2 and P3 intact | **GREEN in both lanes — this is the HOLE, not an arm**: static `rc=0 === RUN=1 --- PASS=1`; runtime `rc=0 marker=1 races=2` | the measurement that justifies AC12: without the needles, deleting the block `D-WORLD-28` mandates is invisible to every other assertion in this sprint (V20) |
| **M12** | delete the whole P1 block (with AC12's needles landed) | `P1 block sentinel "# --- P1 (queue row 48" count=0, want 1: the D-WORLD-28 fail-closed block is absent or duplicated` (**P1a**) | presence needle; rc=1 RUN=1 FAIL=1 (V21) |
| **M13** | delete **only** the `if [ "$root_go_lines" -ne 1 ]; then … fi` branch | `P1 floor-read count guard ("[ \"$root_go_lines\" -ne 1 ]") count=0, want 1: a refusal branch of the D-WORLD-28 floor gate was removed or reworded` (**P1b**) | per-branch needle — the runtime M9 arm **cannot** distinguish M13 from correct (both red, on the below-floor branch), so P1b is its only attributing killer (V21) |
| **M14** | swap the comparator operands | `P1 comparator call "go_version_ge \"$ACTIVE_GO\" \"$ROOT_FLOOR\"" count=0, want 1` (**P1d**) | **the static needle is the SOLE killer in the entire sprint**: swapping inverts the gate into a fail-open that is `rc=0 marker=1 races=2` under BOTH the ambient toolchain and the hostile M9 conditions (V21) |
| **M15** | `ROOT_FLOOR="go$(awk …)"` → `ROOT_FLOOR="go1.26.6"` | `P1 block contains hardcoded Go version literal(s) [go1.26.6]` (**P1e**) | negative assertion: a hardcoded floor is correct today and silently stale at the next root-floor bump; the runtime lane cannot see it until then (V21) |
| **M16** | below-floor branch `exit 1 ;;` → `exit 0 ;;` | `P1 block has 2 \`exit 1\` statements, want 3 (one per refusal)` (**P1c**) | **second sole-killer**: the gutted gate prints its own FATAL and then **exits success**, skipping the race control and everything after it — `rc=0 marker=0` even under M9 conditions (V21) |
| **M18** | move the whole P1 block **below** the race leg | `P1 floor gate is out of order (deny-list@…, P1@…, race leg@…)` (**P1f**) | positional needle: under M9 conditions the run still reds, but the race control was already invoked under the **unvetted** toolchain — the red is not attribution (V21) |
| **M17 — GREEN CONTROL** | reword a **comment** word inside the P1 block (`never hardcoded.` → `(never hardcoded, see D-WORLD-28).`) | **GREEN**: rc=0, `=== RUN`=1, `--- PASS`=1 | anti-brittleness: proves the needle set is *not* "any edit reds" and is not row 49's token count (V21) |

**Green control for all arms:** the unmutated post-sprint tree passes AC1/AC4/AC5, and every arm
ends restored sha256-byte-identical with `git status --porcelain` empty (V1/V2 recipe). The
sibling `TestMiscompileInstrumentProbesPinnedToolchain` and row 42's `TestReproModuleFloor
StaysBelowKnownBadToolchains` stay green at base (V9) and are re-confirmed by the sprint on the
post-sprint tree (AC3's re-confirm clause).

## Milestones (repeat for clarity)

- **M1 — gate + runtime (P1 fail-closed + P2 binding) land** (AC1, AC2, AC3, AC4, AC5, AC11):
  `verify_go.sh` P1 block + `:229` binding + new test in `host/verifygate`.
- **M2 — fence lands** (AC7): `LOAD-BEARING` comment, directive unchanged, post-comment floor
  parses.
- **M3 — hygiene + closure** (AC6, AC8, AC9, AC10): vet/gofmt clean, CI linkage proven, no
  vacuous selector, three-module systemic closure.

## Verification Log

All rows run first-party by the designer at `f1b5d1a` (clean tree, porcelain 0 before and after
every arm), shell `zsh`, `PATH=/opt/homebrew/bin:$PATH`, darwin/arm64, 2026-08-28. KP =
known-positive control carried in the same call. The mutation prototypes ran a standalone Go
replica of the proposed test's logic against the real `racecontrol/go.mod`, root `go.mod`, and
`verify_go.sh` (base) and against temp mutated copies (M1–M8), deleted after the run — the
sprint re-runs every arm against the landed test. **M9 and M10 are no longer deferred**: after the
round-3 block they were run first-party by the CONTROLLER against a temp patched COPY of
`verify_go.sh` placed in `scripts/` and deleted afterwards (V14, V16), with the production script's
sha256 asserted identical before and after; the sprint re-runs them against the LANDED block. Rows V1–V4, V6, V8, V9 are the round-1 quorum
re-measurements; every `go env` invocation uses the space-separated `go env GOTOOLCHAIN
GOVERSION` form (gemini's syntax correction). **Rows V10–V13 are the round-2 re-verification measurements, rows V14–V18 the round-3
ones, and rows V19–V23 the SPRINT-PLANNER's land-and-restore measurements against the actually
landed artifacts (iteration 144, `opus` planner lane) — all run first-party at HEAD `8bb9214` on
2026-09-01** (clean tree, porcelain 0 before
and after every ARM land-and-restore).

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Worktree is `f1b5d1a`, clean; module census is exactly 3; the `racecontrol` floor is `go 1.22` with the V24 hash; worktree `go.mod` untouched by the round | `git rev-parse HEAD`; `git status --porcelain \| wc -l`; `find . -name go.mod -not -path './.git/*'`; `shasum -a 256 design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod` | `f1b5d1a79d89252…`; `0` (planned doc + log are the only non-clean entries, untouched); exactly `./go.mod`, `…/repro/go.mod`, `…/racecontrol/go.mod`; sha256 `ab782f11db0f7f259f73dd55a58eaf5a30b871bb79bd98bacbe964d50efc025b` (identical to row 42's V24 restored hash) |
| V2 | **THE FINDING re-measured under the execution-bound runtime**: the disarming edit is a floor above the bound toolchain, and the bound toolchain is `ACTIVE_GO` | `ACTIVE_GO=$(go env GOVERSION)`; `race_control_output="$(cd racecontrol && GOTOOLCHAIN="$ACTIVE_GO" go run -race . 2>&1)"` (base floor `go 1.22`); then `sed go 1.22→go 1.27.0` and re-run the same command; `cp` restore + sha256 + porcelain | base: rc=1, **2** `WARNING: DATA RACE` (the known-positive fires under the binding); floor `go 1.27.0`: rc=1, output `go: go.mod requires go >= 1.27.0 (running go 1.26.6; GOTOOLCHAIN=go1.26.6)`, **0** DATA RACE → `verify_go.sh:232` would FATAL. Restore byte-identical `ab782f11…`, porcelain 0; post-restore control re-fires (rc=1, 2 DATA RACE). (Also confirmed: floor `go 1.26.6` == `ACTIVE_GO` still fires **2** DATA RACE — equality is armed, so the old V24 edit stops being a disarm once execution is bound) |
| V3 | **The execution-binding chain**: `ACTIVE_GO` (root `go env GOVERSION`) = root floor, and the nested module's ambient toolchain differs; `go env` takes space-separated keys (gemini's fix) | `go env GOTOOLCHAIN GOVERSION` at root; `go version` at root; `go env GOTOOLCHAIN GOVERSION` inside `racecontrol/` | root: `auto go1.26.6` (two-arg space form rc=0), root `go version go1.26.6 darwin/arm64`; inside `racecontrol/`: `go1.26.4 auto` — the nested auto-selection would pick `go1.26.4`, which the execution-binding edit removes; `ACTIVE_GO = go1.26.6 = root go.mod floor`, so `racecontrol floor <= root floor` ⇒ `floor <= ACTIVE_GO` |
| V4 | `racecontrol` is driven only by `verify_go.sh:229`; the race leg is lines 226–235 and the deny-list/capture are `:217–224` | `git grep -n "racecontrol" -- '*.go' '*.sh' '*.yml'` (excluding design docs); `grep -n "ACTIVE_GO=\|go run -race\|case \\"\\$ACTIVE_GO\\"" scripts/verify_go.sh` | exactly one hit: `scripts/verify_go.sh:229` `cd .../racecontrol && go run -race . 2>&1`; `:217` `ACTIVE_GO=$(go env GOVERSION)`, `:218` `case "$ACTIVE_GO"`, `:226` race-leg header, `:232` `grep -q 'WARNING: DATA RACE'` else FATAL |
| V5 | **Prototype: base GREEN and M1–M8 verdicts** (the test's logic, run against real + mutated files) | standalone replica of `moduleGoFloor`+`version.Compare` over `racecontrol/go.mod`, root `go.mod`, `verify_go.sh` (base) and temp mutated copies (M1–M8) | base `PASS: racecontrol floor go1.22 <= root floor go1.26.6; execution needle count=1`; M1 `RACE-BOUND-FAIL` (`go1.27.0` > `go1.26.6`); M2 `PASS` (`go1.26.6` == `go1.26.6`); M3 `racecontrol floor "gobanana" invalid`; M4 `found 0 go lines, want 1`; M5 `RACE-BOUND-FAIL` (`go1.22` > `go1.20`); M6 `root floor "gobanana" invalid`; M7 `execution needle count=0, want 1`; M8 `open …/go.mod: no such file or directory` |
| V6 | **M2 runtime proof: the control is armed at the equality floor** | `GOTOOLCHAIN=go1.26.6 go run -race .` from `racecontrol/` with floor `go 1.26.6`; `GOTOOLCHAIN=go1.26.6 go version` | rc=1, **2** `WARNING: DATA RACE` — the equality floor still builds and fires the control, so `<=` is the right bound |
| V7 | Version comparison semantics for every token the gate compares | standalone `go/version`: `Compare`/`IsValid` on the bound and mutant tokens | `Compare(go1.22,go1.26.6)=-1`; `Compare(go1.27.0,go1.26.6)=1`; `Compare(go1.26.6,go1.26.6)=0`; `Compare(go1.22,go1.20)=1`; `IsValid(gobanana)=false` — `go/version` is stdlib, no new dependency |
| V8 | **Round-1 R2 premise: a `LOAD-BEARING` comment containing `go 1.22` does not break `moduleGoFloor`**; and the base has no fence yet | `moduleGoFloor`-replica over the real `racecontrol/go.mod` and over a temp file containing the proposed comment block plus the real directive; `grep -c "LOAD-BEARING" racecontrol/go.mod` | real file: `floors=[go1.22] count=1`; comment+directive probe: `floors=[go1.22] count=1` (exactly one awk-equivalent column-0 `go ` line; the `//` comment is ignored by `strings.HasPrefix(line, "go ")`); base `LOAD-BEARING` count = **0** (AC7 base-RED); probe file discarded, original restored byte-identically |
| V9 | Base gate lane: all existing pin tests green; the new name is unclaimed and its scoped run reports `no tests to run` (the AC1 vacuity trap); no-AILANG_BIN static lane works; hygiene/collision baselines | `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run 'TestReproModuleFloorStaysBelowKnownBadToolchains\|TestCanaryDeclaresPositiveArmOnly\|TestMiscompileInstrumentProbesPinnedToolchain\|TestMiscompileInstrumentStepIsGatedInCI\|TestGoToolchainPinsAgreeAndMatchJobList' -count=1`; `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1`; `env -u AILANG_BIN GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestMiscompileInstrumentProbesPinnedToolchain$' -count=1`; `grep -rn "TestRaceControlFloorStaysBelowRootToolchain" host/` | existing: `ok github.com/sunholo-data/ailang-world/host/verifygate 0.276s`; new name: `ok ... [no tests to run]`, rc=0 (green-at-base measuring nothing); no-AILANG_BIN scoped rc=0 (static lane is env-independent); no existing test with that name |
| V10 | **R2-a: CI linkage proven mechanically** — a test added to `host/verifygate/` reaches CI through `verify_go.sh`'s `./...` leg | `grep -n 'run: ./scripts/verify_go.sh' .github/workflows/ci.yml`; `sed -n '163,167p' .github/workflows/ci.yml`; `grep -n 'go test ./\.\.\.' scripts/verify_go.sh`; `go list \| grep verifygate` | `ci.yml:166: run: ./scripts/verify_go.sh`; step at `:163` = "go build + test gate (replay tests run against pinned AILANG_BIN)"; `scripts/verify_go.sh:256` echo + `:258: go test ./... -count=1` (not commented); `go list ./...` includes `github.com/sunholo-data/ailang-world/host/verifygate`. So a new test in that package is exercised by CI (AC9) |
| V11 | **R2-b: selector audit — every `-run` selector in this doc, and what each matched when run.** The only malformed one was V9's escaped-`\|` alternation (grep-style `\|` = literal pipe to Go regexp → null match); corrected to single `\|`. Go-test ACs bind counted `=== RUN`/`--- PASS`, never a bare rc (AC10). | `-run '^TestRaceControlFloorStaysBelowRootToolchain$'`; `-run 'TestNoSuchRaceControlFloorTestZZZ'`; `-run '^TestMiscompileInstrumentProbesPinnedToolchain$'`; `-run '…\\|…'` (V9's escaped form); `-run '…\|…'` (corrected single-pipe form) | anchor (AC1/AC2, absent at base) → `ok … [no tests to run]` rc=0; nonsense (AC1 control) → `ok … [no tests to run]` rc=0 (RED under the rule); AC5 anchor → 1× `=== RUN` / 1× `--- PASS` rc=0; escaped `\\|` alternation → `ok … [no tests to run]` rc=0, **0** tests; corrected `\|` alternation → **5**× `=== RUN`, **5**× `--- PASS`, rc=0 (all five sibling pin tests) |
| V12 | **R2-c: under `GOTOOLCHAIN=local` the base sits below the root floor, refuting the old proof's assumption; and the deny-list only covers a sub-range** | at root `go env GOTOOLCHAIN GOVERSION`; `GOTOOLCHAIN=local go env GOVERSION`; `GOTOOLCHAIN=local go version`; inside `racecontrol/` `go env GOTOOLCHAIN GOVERSION`; `grep '^go ' go.mod` (root and racecontrol) | root ambient `auto go1.26.6`; `GOTOOLCHAIN=local` → `go1.26.4` (**below** root floor `1.26.6`); `go version go1.26.4 darwin/arm64`; inside `racecontrol/` → `auto go1.26.4`; root floor `go 1.26.6`, racecontrol floor `go 1.22`. The deny-list (`:218–224`) covers only `go1.26.0..go1.26.5`, so it happens to catch this rig's `go1.26.4`, but a base below the root floor **outside** the deny-list (e.g. `go1.25.x`) passes it — only P1 carries the floor (V13/M9, AC11) |
| V13 | **ARM re-verification at HEAD 8bb9214** — the racecontrol go.mod sha256 is unchanged at `ab782f11…`; the four ARM rows re-fire and the restore is byte-identical (2026-09-01) | ARM0 `go run -race .`; ARM0b `GOTOOLCHAIN=local go run -race .`; ARM1 floor `go 1.22→1.26.6` then `go run -race .`; ARM2 same floor then `GOTOOLCHAIN=local go run -race .`; restore + `shasum -a 256 go.mod`; post-restore `go run -race .`; `git status --porcelain` | ARM0: rc=1, **2**× `WARNING: DATA RACE` (pristine positive fires); ARM0b: rc=1, **2**× (floor 1.22 satisfiable by local `go1.26.4`); ARM1: rc=1, **2**× (warm cache rescues `GOTOOLCHAIN=auto`); ARM2: rc=1, **0**× DATA RACE, output exactly `go: go.mod requires go >= 1.26.6 (running go 1.26.4; GOTOOLCHAIN=local)` — **DISARMED**; restore sha256 byte-identical `ab782f11…`, porcelain 0; post-restore control re-fires (rc=1, **2**×) |
| V14 | **P1 RUNTIME ARMS — MEASURED, not planned (round-3 `oc-glm-5-2`/`gpt5-6-sol`; controller-run 2026-09-01, HEAD `8bb9214`). The `verify_go.sh` in the repo was NEVER edited: every arm ran against a temp patched COPY placed in `scripts/` (line 17 is `cd "$(dirname "$0")/.."`, so a copy run from `/tmp` resolves the root to `/`), deleted afterwards.** | `AILANG_BIN=/tmp/ailang-v0300/ailang`; A1 `./scripts/zz_vg_p1.sh`; A2 `GOTOOLCHAIN=local ./scripts/zz_vg_p1_nodeny.sh` (deny-list case body neutered); A3 `GOTOOLCHAIN=local ./scripts/zz_vg_nop1_nodeny.sh` (identical, P1 block ABSENT); A4 `GOTOOLCHAIN=local ./scripts/zz_vg_p1.sh` (deny-list live). Each arm counts a `HARNESS: reached end of race leg` marker and `WARNING: DATA RACE` lines. | **A1 GREEN**: rc=0, `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`, race leg reached (marker=1), **2** DATA RACE — P1 is a no-op on the normal path and the control still arms. **A2 = M9 RED**: rc=1, `FATAL: active toolchain go1.26.4 is BELOW the root module floor go1.26.6`, marker=**0**, races=**0** — P1 exits before `:229`, with the deny-list unable to supply the red. **A3 = the discriminating CONTROL**: rc=**0**, marker=**1**, **2** DATA RACE — the byte-identical scenario without the P1 block sails through, so **P1 alone** supplies A2's red. **A4**: rc=1, `FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go's` — with the deny-list live it fires FIRST and its message is distinguishable from P1's, so an arm's red is always attributable. Production `scripts/verify_go.sh` sha256 `27eab122f4b15ac1febe0fb3aed9886d900a03d6da65377878b08843f337cd2b` and root `go.mod` sha256 `7a298361…` identical before and after; `git status --porcelain` shows only the design doc. |
| V15 | **CONTROLLER-RUN (iter-144, 2026-09-01, HEAD `8bb9214`) — P1 lives in `bash`, not in Go, so V7's `go/version` machinery is a SPECIFICATION of the order and not an implementation of it. A portable `awk` comparator realises exactly that order; the naive shell form does not.** | `go_ge()` = `awk` over `^go1\.[0-9]+(\.[0-9]+)?$`-validated tokens, numeric component compare, `exit 0` = A>=B, `exit 1` = A<B, `exit 2` = malformed; 13-case battery; plus the lexical control `[[ "go1.9" < "go1.10" ]]` | **13 of 13** cases match V7's `go/version` verdicts: `go1.26.4 < go1.26.6` (1), `go1.25.6 < go1.26.6` (1), `go1.26.7 >= go1.26.6` (0), equality (0), `go1.26` vs `go1.26.0` (0 in both directions), and the ordering trap `go1.10 >= go1.9` (0) / `go1.9 < go1.10` (1). Malformed tokens (`gobanana`, `go1.26.6rc1`, `devel`) all exit **2** — M10's instrument-floor branch, distinct from M9's refusal branch. **Lexical control: `[[ "go1.9" < "go1.10" ]]` is FALSE**, i.e. a naive string compare orders `go1.9` ABOVE `go1.10` and would fail OPEN on exactly the shape V7 warns about. `sort -V` is not used: this rig ships BSD `sort 2.3-Apple (199)` and CI runs GNU coreutils, so a comparator built on it would be two implementations; `awk` is one. |
| V16 | **P1 BRANCH BATTERY — every refusal branch and every instrument floor, one arm each (controller-run 2026-09-01).** The comparison branch is isolated by forcing `ACTIVE_GO` against a fixed root floor `go1.26.6`; the floor-read branch is isolated by varying the `go.mod` shape. | P1 block run standalone in a `mktemp -d` whose `go.mod` is written per arm; `FORCE_ACTIVE_GO` overrides the `go env GOVERSION` read | **Comparison branch** (root floor `go1.26.6`): `go1.26.6` equal → rc=0 GREEN; `go1.27.0` above → rc=0 GREEN; `go1.26.4` below → rc=1 `is BELOW the root module floor`; **`go1.25.6` below and NOT deny-listed → rc=1** — the case the miscompile deny-list structurally cannot catch, which is why P1 is not redundant with it; `go1.9` (ordering trap) → rc=1. **Malformed branch** (distinct message): `devel +abc` → rc=1 `cannot order toolchain tokens`; `go1.26.6rc1` → rc=1 same branch. **Floor-read instrument branch**: 0 `go ` lines → rc=1 `has 0 column-0 'go ' lines, want exactly 1`; 2 `go ` lines → rc=1 `has 2 …`; `go banana` → rc=1 `ROOT_FLOOR=gobanana` on the ordering branch. **Negative control**: a TAB-indented `go 1.26.6` line counts **0**, matching `moduleGoFloor`'s column-0 semantics — so the shell reader and the Go reader agree on what a floor line is. Every refusal carries a message naming which branch fired; no two branches share one. |
| V18 | **R2/R3 `gemini-3-1-pro`: the `go.mod` → comparator extraction, and the fail-open the naive shell form would have shipped (controller-run 2026-09-01).** | `awk '/^go /{n++} END{print n+0}' go.mod`; `ROOT_FLOOR="go$(awk '/^go /{print $2; exit}' go.mod)"`; both the raw and the normalised token piped to `grep -cE '^go1\.[0-9]+(\.[0-9]+)?$'`; and the lexical control `[[ "go1.9" < "go1.26.6" ]]` / `[[ "go1.9" < "go1.10" ]]` in `bash` | column-0 `go ` lines = **1**; `ROOT_FLOOR=go1.26.6`; `ACTIVE_GO=go1.26.6`. **The normalisation is load-bearing exactly as objected**: the RAW `go.mod` token `1.26.6` matches the V15 grammar **0** times, the `go`-prefixed form matches **1** — without the prefix the comparator would exit 2 and the gate would fail closed for the wrong reason. **Lexical fail-open control**: `[[ "go1.9" < "go1.26.6" ]]` is **FALSE**, i.e. a naive string compare reports `go1.9 >= go1.26.6` and would let a toolchain seventeen minor versions below the floor through the gate; `[[ "go1.9" < "go1.10" ]]` is FALSE too. This is why the block uses `awk` and why V15's battery is the binding specification of the order. |
| V19 | **THE FENCE, LANDED AND MEASURED (AC7's replacement clauses) — planner lane, iteration 144.** The doc's Example-2 fence (8 comment lines) was appended above `module` and every AC7 clause run against the fenced file. | `grep -c 'LOAD-BEARING'`; `grep -n '^go 1.22$'`; `awk '/^go /{n++} END{print n+0}'`; `git diff --unified=0 -- <file> \| grep '^-' \| grep -v '^---' \| wc -l`; `grep -cF 'GOTOOLCHAIN="$ACTIVE_GO" go run -race .'` on the fenced file; then the P3 test and the race control | `LOAD-BEARING` **1** (base **0**); `^go 1.22$` **1**, at line **14**; column-0 `go ` count **1** (`moduleGoFloor` still finds exactly one floor); **deleted/modified lines = 0** — this is the executable form of "the directive byte-unchanged", and it is what replaces the non-computable sha256-composition clause; needle-collision check on the fenced file **0** (the fence cannot satisfy the P2 needle); post-fence `TestRaceControlFloorStaysBelowRootToolchain` rc=0 RUN=1 PASS=1 and the race control still rc=1 with **2** `WARNING: DATA RACE`. |
| V20 | **M11 — THE HOLE THAT JUSTIFIES AC12: deleting the whole P1 block is a silent GREEN in BOTH lanes.** P2 and P3 intact; only the D-WORLD-28 block removed. | static: `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1 -v`; runtime: the §V14 patched-copy harness with its `HARNESS: reached end of race leg` marker | static **rc=0, `=== RUN`=1, `--- PASS`=1**; runtime **rc=0, marker=1, races=2**. Deleting the block the attended ruling mandates was invisible to every other assertion in the sprint — while P2's binding already had M7 for exactly this reason. This row is the whole argument for AC12. |
| V21 | **THE AC12 NEEDLE SET, LANDED AND KILLED — 6 named RED mutants + 1 GREEN control, plus the runtime cross-check that shows the static lane is the sole killer for two of them.** Every mutant landed in `scripts/verify_go.sh`, asserted landed by a count that moved, then restored byte-identical; production sha256 `27eab122…` re-asserted after each. | per-arm: land the mutant → assert the landing count moved → `go vet` → the AC1 selector → restore → `cmp`. Runtime cross-check: the same six variants under (a) the **ambient** toolchain and (b) **M9 conditions** (deny-list neutered, `GOTOOLCHAIN=local`) | **M12** (whole block deleted) → RED on **P1a**, sentinel count 1→0. **M13** (one refusal branch deleted) → RED on **P1b**, guard 1→0 and in-block FATALs 3→2. **M14** (comparator operands swapped) → RED on **P1d**. **M15** (floor hardcoded) → RED on **P1e**. **M16** (below-floor `exit 1`→`exit 0`) → RED on **P1c**, in-block `exit 0` 2→3. **M18** (block moved below the race leg) → RED on **P1f**. **M17 GREEN CONTROL** (a *comment* word reworded inside the block) → **rc=0 RUN=1 PASS=1** — the set is not "any edit reds". **The cross-check is the finding:** under the **ambient** toolchain, the only condition CI ever runs, all six gutted variants are `rc=0, marker=1, races=2` — byte-for-byte indistinguishable from the correct block, so the runtime lane is blind to every one of them in CI. Under M9 conditions **M14 stays `rc=0 marker=1 races=2` (fail-open)** and **M16 is `rc=0 marker=0` — it prints its own FATAL and then exits success**, skipping the race control and everything after it. For M14 and M16 the static needle is the **sole killer in the entire sprint**. |
| V22 | **M6 REFUTED and M6′ measured** — the doc's root-`go banana` arm runs **zero** tests. | `sed` root `go 1.26.6` → `go banana`, then the AC1 selector; and root `go 1.26.6` → `go 1.26.6 // pin`, same selector | M6: **rc=1, `=== RUN`=0, `--- PASS`=0, `--- FAIL`=0**, `go: errors parsing go.mod: go.mod:3: invalid go version 'banana'` — Go's module loader rejects the file before the test binary is built, so under this doc's own AC10 rule it is not evidence and `!version.IsValid(rootFloor)` had no killer. **M6′**: **rc=1, `=== RUN`=1, `--- PASS`=0, `--- FAIL`=1**, `toolchain_pin_gate_test.go:574: instrument failure: root module floor "go1.26.6 // pin" is not a valid Go version` — Go accepts the trailing comment so the module loads and the test runs. Control: **M3** (the *racecontrol*-side `go banana`) is unaffected and reds on its own named branch, because `racecontrol/go.mod` is a nested module `go test ./host/verifygate/` never loads. |
| V23 | **M10 REFUTED by its own P1-REMOVED control; M10a′ and M10b measured as sole killers.** | M10: delete the root `go 1.26.6` line, run the patched-copy harness with P1 present and with P1 REMOVED. M10a′: TAB-indent the directive. M10b: duplicate it. Each arm reads `go env GOVERSION`, the marker count, the race count and the FATAL text. | **M10**: column-0 `go ` lines **0** and `go env GOVERSION` drops to **`go1.26.4`** — deleting the root floor removes what pushes `GOTOOLCHAIN=auto` upward. P1-present → `rc=1 marker=0 FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go's`; **P1-REMOVED control → byte-identical red**. The pre-existing miscompile deny-list supplies M10's red, not P1, so M10 fails the very attribution shape V14/A3 exists to enforce. `go banana` (M10c) fails identically. **M10a′**: `GOVERSION` stays `go1.26.6` (deny-list passes); `rc=1 marker=0 races=0`, `FATAL: root go.mod has 0 column-0 'go ' lines, want exactly 1;` — **P1-REMOVED control rc=0 / marker=1 / races=2** ⇒ sole killer. **M10b**: same, `has 2 column-0 'go ' lines`, same control. Root `go.mod` restored byte-identical after every arm. |

## Related Documents

- [`../implemented/w-canary-control-does-not-survive-a-floor-raise.md`](../implemented/w-canary-control-does-not-survive-a-floor-raise.md) —
  row 42: the item that measured **V24** (this item's source evidence), built the exact test
  shape (`TestReproModuleFloorStaysBelowKnownBadToolchains`, `moduleGoFloor`,
  `shellAssignmentValues`, `repoRoot`) and helpers this item reuses, and **explicitly deferred
  `racecontrol/` to queue row 48** because its correct binding ("buildable by whatever
  toolchain verify_go.sh happens to run") was not statically knowable there. This item decides
  that binding by *making* the toolchain statically knowable (`GOTOOLCHAIN="$ACTIVE_GO"`) and
  closes the deferral.
- [`../implemented/w-race-gate-blindspot.md`](../implemented/w-race-gate-blindspot.md) — where
  the canary and the nested-module pattern (including `racecontrol/` as the race-detector
  known-positive) were born.
- [`../implemented/w-setup-go-pin-unguarded.md`](../implemented/w-setup-go-pin-unguarded.md) —
  row 41: built the file this item extends and the root-floor binding this item's anchor (the
  root `go.mod` floor) sits alongside.
- [`../implemented/w-miscompile-instrument-inert-in-ci.md`](../implemented/w-miscompile-instrument-inert-in-ci.md) —
  row 44: owns `run.sh`'s CI inertness and `ci.yml:172`; the declared reason a *static* gate is
  this item's gating lane rather than the runtime lane.
- `design_docs/world-mission.md` — queue row 48 (this item); rows 42 (the sibling binding),
  43 (floor-raise coupling inventory), 44, 45, 46.
- `design_docs/verification/w-race-gate-blindspot/` — `racecontrol/` (this item's target),
  `repro/` (row 42's target), `run.sh` (whose `KNOWN_GOOD` anchor round 1 rejected; the
  `PINNED`/root-floor chain remains row 41/42's concern).

## Future Work

- **Option B (`GOTOOLCHAIN=local`)** — a loudness improvement that removes the cache rescue; if
  adopted, it belongs in `verify_go.sh`'s race leg (row 44's lane) as a one-line change
  complementing this static gate and execution binding.
- **An explicit repo-wide "oldest supported toolchain" policy** — would give future items a
  single maintained anchor to bind all nested modules to, instead of deriving the bound from the
  root module floor.

---
**DESIGN_DOC_PATH**: `design_docs/planned/w-racecontrol-floor-bump-disarms-the-race-control.md`

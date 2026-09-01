# w-wiring-test-step-scoping-imprecise-under-key-reorder

**Status**: DRAFT, REVISION 2 + ROUND-3 CARVE-OUT PATCH — direction **RATIFIED by attended ruling `D-WORLD-30`**
(Mark Edmondson, attended, 2026-09-01; ledger row in `design_docs/world-mission.md`), no longer
chosen by the loop. History: round 1 BLOCKED at full strength (`.synthesis.absent_reviewers` =
`[]`; `gpt5-6-sol` reject, `gemini-3-1-pro` reject, `oc-glm-5-2` pass, controller reject); the
one protocol-mandated revision (revision 1: relative rule → absolute step-column rule, V15-V17)
went back to quorum and round 2 BLOCKED with all three external reviewers rejecting, the
surviving objection disputing the DESIGN DIRECTION itself — so the loop parked the row and asked
the human. The ruling, verbatim: *"ANSWERED — A — LINE SCAN, hardened, deriving the anchor from
the SHALLOWEST enclosing `steps:` rather than the nearest one (measured to catch the round-2
attack that the doc's nearest-anchor derivation is blind to). This gate is a repo-internal
tripwire against our own accidental regressions, not an adversarial boundary: the threat model
that would justify B is an actor who can already edit `ci.yml`, and such an actor can simply
delete the test. Not worth adding the second direct dependency to a `go.mod` that has exactly
one. The residual is accepted and named: the scan stays text-level, and a future re-indent of
`ci.yml` makes the gate fatal until the constant is updated."*
**What changed since revision 1, all by measurement at HEAD `f5cfb1b` (controller-measured rows
V20-V26):** the anchor derivation in Mechanism step 1 goes from the NEAREST enclosing `steps:`
to the SHALLOWEST one — V25 measures the divergence on the round-2 `steps:`-decoy attack, where
the nearest anchor mis-derives `stepCol=12` and fails to adjudicate the step (an Invariant-A
FATAL on the controller's construction; GREEN/blind on the ledger's iter-142 layout — both
readings recorded in V25, neither overwritten) while only the shallowest anchor produces the
correct RED. ARM A and ARM B are re-derived first-party at this head, not inherited (V21, V23),
with two new controls (V22 attribution, V24 position); the V12 citation `oc-glm-5-2` caught as
falsified in round 2 is repaired by a WIDENED grep, not deletion (V12 addendum); the
`AILANG_BIN` requirement of the package-wide test command is baselined and propagated into the
ACs (V26); the mutation drill gains MUT-O and MUT-P for the two new arms, and MUT-O's OLD-verdict cell is no longer owed — **V28 measures it: the landed gate at HEAD CATCHES the round-2 `steps:` attack (rc=1 at `ci.yml:182`, the guarded step's own flag line), incidentally, by anchoring on the in-scalar decoy and widening. So the ratified anchor prevents a regression this doc would have INTRODUCED; it does not close a live hole. The live holes remain V21 (false positive) and V23 (fail-open).** This is the SECOND
protocol-mandated revision, taken under an attended ruling rather than under the quorum's
one-revision allowance. Option B (semantic YAML parse, a new `go.mod` dependency) is **REJECTED
BY THE HUMAN** and appears below only under Alternatives Considered as rejected-by-ruling.
**ROUND 3 (iteration 147, 2026-09-02) — BLOCKED at FULL STRENGTH, then closed under the
narrow-refinement carve-out; every objection MEASURED rather than forwarded.**
`.synthesis.absent_reviewers` = `[]`, cross-checked against
`[.reviewers[]|select(.present==false)]` with `has("synthesis")` as the sibling control;
`metered=$0.18372421` (`gpt5-6-sol` $0.10932 reject, `gemini-3-1-pro` $0.047032 reject,
`oc-glm-5-2` $0.02737221 **pass** — its first pass since round 1). Per-surface record, as the
round-3+ discipline requires: the two rejections land on **different, newly-introduced**
surfaces (ruling fidelity; AC baselining), not on one surface with the others clearing — so the
signal is an immature revision, not a doc that needs splitting.
- `gpt5-6-sol` — **ruling fidelity**: Residual 7 reclassified the ruling's *"fatal until the
  constant is updated"* clause as belonging to an alternative form, with no evidence the human
  drew that line. **Its first verbatim `proposed_fix` is APPLIED**: Mechanism step 1b adds
  `const expectedStepCol = 6` with a loud `t.Fatalf`, plus mutation arm **MUT-Q** proving it,
  and Residual 7 withdraws the reclassification. No amended ruling is needed, and the loop did
  not pick which half of a human's sentence to honour — it implements both.
- `gemini-3-1-pro` — **AC baselining**: `gofmt -l host/verifygate/` and `./scripts/verify_ail.sh`
  appeared in the ACs with no Verification Log row. Measured, not forwarded: **gofmt is clean at
  base (V29, firing control), and `verify_ail.sh` is genuinely rc=1 at base on this rig** — the
  PATH `ailang` is a `-dirty` dev build the gate's own guard refuses — so **AC7 was
  unsatisfiable as written** and now carries `AILANG_BIN=/tmp/ailang-v0300/ailang` (V30).
  The reviewer's fear was unfounded for one half and correct for the other.
- `oc-glm-5-2` — **pass**, with a traceability note: MUT-O/MUT-P pinned no sha256. Applied —
  all three mutant hashes re-derived and pinned (V32).
- **What the objections BOUGHT beyond their own text, and the reason this round was worth
  $0.18:** measuring `gpt5-6-sol`'s objection produced **Declared Residual 8** — the backward
  `steps:` scan is unbounded, so on a re-indented file the shallowest anchor jumps to the OTHER
  JOB's `steps:` key (V31). That is a defect the RATIFIED rule introduces. It is recorded as a
  residual rather than escalated, on this mission's own discriminator: it fails **loud** (an
  Invariant-A/B refusal, plus step 1b's pin), never silent, so it is not the false-green shape
  that re-parked row 50.
- **Carve-out justification, stated so it can be audited:** all three objections carry concrete
  reviewer-authored `proposed_fix` text; none disputes the design DIRECTION (which is in any case
  human-ratified and not re-litigable by a reviewer); every applied fix is the reviewer's own
  words, not a controller invention; and the one controller-authored addition (Residual 8) is a
  DECLARED MEASUREMENT, not a resolution. Routed to sprint-planner without a fourth quorum round.

**Target**: World iter-142 / queue row 52
**Priority**: P2
**Estimated**: ~0.1d (measured scope: ONE function in ONE test file; see V-A)
**Dependencies**: None (row 52 is gated on nothing)
**Filed by**: the `sonnet` evaluator at the row-44 landing (iter-133), measured non-exploitable
then; **measured exploitable in BOTH directions this iteration** (V2, V3)

## Problem Statement

**The row-44 wiring gate FAILS OPEN.** `TestMiscompileInstrumentStepIsGatedInCI`
(`host/verifygate/toolchain_pin_gate_test.go:492-533`) is the test that forbids
`continue-on-error` in the miscompile instrument's own CI step — the flag that killed this
instrument the first time by converting its loudest output into silence. I measured (V3, first
party, this worktree): put `continue-on-error: true` ON the miscompile step as a real, live step
key, and place a line whose *trimmed* text begins `- name:` inside that same step's `run: |`
block, between the flag and the reproducer path:

```yaml
      - name: Measure compiler reproducer (platform-conditional, gated)
        timeout-minutes: 15
        continue-on-error: true
        run: |
          cat <<'YAMLDOC'
          - name: documentation example
          YAMLDOC
          ./design_docs/verification/w-race-gate-blindspot/run.sh
```

Ruby's YAML parser confirms the mutant means what the attack needs it to mean: the miscompile
step carries `continue-on-error=true` as a step key, the nested line is script *content*, the
job still has exactly 10 steps. The test result: **rc=0, `--- PASS`.** The forbidden flag is
live on the very step this gate exists to guard, and the gate is green.

The mechanism: the test locates "the miscompile step's own block" by walking BACKWARD from the
line containing the reproducer path to the nearest line whose `strings.TrimSpace(...)` has
prefix `- name:` (`toolchain_pin_gate_test.go:511`), then FORWARD to the next such line (`:524`).
`TrimSpace` discards indentation, so **nesting depth is invisible to the scan**: the decoy
inside the `run:` block reads as a step boundary, `start` lands BELOW the flag, and the flag
falls outside the scanned range.

**The same defect also fires in the opposite direction — a false positive that violates the
test's own quorum-ratified contract** (V2, first party): add `continue-on-error: true` to the
PREVIOUS, unrelated step (`go build + test gate`), and reorder the miscompile step so `- name:`
is not its first key (`- timeout-minutes: 15` first). The backward walk skips the miscompile
step's dash line (its trimmed prefix is now `- timeout-`, not `- name:`) and lands on the
*previous* step's `- name:`. Result: **rc=1, `ci.yml:166 re-introduces "continue-on-error" in
the miscompile step`** — and ci.yml:166 is the `go build + test gate` step's flag line (V2's
mutant listing shows it verbatim). The test's own comment at `:502-505` cites quorum round-2
R1: *"A flag on an unrelated step is that step's business"* — the exact boundary this arm
breaks. The whole history of this test is a fight against over-broad banning (round 1 banned
the platform tokens repo-wide and redded against its own documentation on arrival); the scan
quietly re-creates the class round 2 removed.

**This CORRECTS queue row 52 where they disagree.** The row records the defect as "measured
non-exploitable … every reordering the judge tried widened the scanned range; none narrowed
it". That sentence was — exactly as the row itself predicted — *a property of the mutations
someone thought to try, not of the parser*. V3 is the narrowing counterexample, and V2 shows
the widening direction is not harmless either. Two further corrections: the row's "hand-rolled
line scans in three places" counts ci.yml *readers*, not step *scopers* — the step-scoping
defect has exactly ONE call site (V5, both scan lines in the same function; the third
`HasPrefix(strings.TrimSpace` hit matches a SQLite pragma and is unrelated); and the other two
ci.yml readers do whole-file substring counting and are unaffected (V6, V7).

**THE GENERALISATION (the row's, now discharged rather than hypothetical):** *a scoping bug
that currently fails safe is still a scoping bug, and "it happened to widen" is a property of
the mutations someone thought to try, not of the parser.*

## Evidence (first-party measurements, with commands)

Rows V1-V19 were made by me in this checkout at `dev` = `7077e455` (pristine,
`git status --porcelain` = 0 lines) on 2026-09-01, bracketed by pristine controls; ci.yml was
restored byte-identical (sha256 `aed8e186…` before, between, and after — V1, V4). The
controller's iteration-142 measurements are independently reproduced here; every number below
is my own. Full commands and observed output: §Verification Log.

**Revision-2 rows V20-V26 and V28 were measured by the CONTROLLER, first-party, at HEAD `f5cfb1b`** in
a scratch worktree that is a SIBLING of the repo (never `/tmp`), and handed to this revision
with commands and observed output. ARM A and ARM B — the two arms row 52 rests on — are
RE-DERIVED there against the landed gate, not inherited from V2/V3; every mutant was asserted
LANDED by sha256 before its result was read; intended effects were asserted against
`ruby -ryaml` (the system's own view of ci.yml), not against bytes; `go vet ./host/verifygate/`
rc=0 was read BEFORE any test result; and the pristine control was green before and after the
batch (porcelain 0). ci.yml's sha at `f5cfb1b` is the same `aed8e186…` as at `7077e455` (V20),
and the gate function under change is line-for-line unchanged between the two heads — the test
file's 94 intervening insertions all land after it (V27, designer-measured) — so the two
batches measured the same ci.yml bytes against the same locator code.

- **V1 pristine control**: targeted test rc=0 `--- PASS`; ci.yml sha256
  `aed8e186fb57036eb6b03509cbb668850d577d46e6cc68b30e7a4c042108ed85`.
- **V2 ARM A (false positive)**: coe on the unrelated `go build + test gate` step + miscompile
  key reorder. Mutant sha256 `28731ce9…`; ruby YAML view: step[5] keys
  `["name","timeout-minutes","continue-on-error","run"]` coe=true, step[6] keys
  `["timeout-minutes","name","run"]` coe=nil, 10 steps. Test **rc=1** blaming ci.yml:166 —
  the *unrelated* step's flag line.
- **V3 ARM B (false negative)**: coe live on the miscompile step + nested `- name:` decoy in
  its `run: |` block. Mutant sha256 `66bad3af…`; ruby YAML view: the miscompile step carries
  coe=true, 10 steps, decoy is script content. Test **rc=0 `--- PASS`. The gate fails open.**
- **V4 restore**: sha256 back to `aed8e186…`, test rc=0, porcelain 0.

## Solution Design

### Mechanism: ABSOLUTE step-column location + two loud invariants (revised after quorum R1; anchor hardened to the SHALLOWEST enclosing `steps:` per ruling `D-WORLD-30`)

Rewrite the locator inside `TestMiscompileInstrumentStepIsGatedInCI` (the ONE call site, V5) to
anchor on the **absolute column YAML gives block-sequence items under the enclosing `steps:`
key** — derived from the file, never hardcoded. Round 1 proposed a *relative* test
(`indent(dash) < indent(identifying line)`); that rule is superseded by measurement — see §Why
the round-1 relative rule was rejected. Revision 1 derived the anchor from the NEAREST
`steps:`-shaped line; that derivation is superseded by ruling `D-WORLD-30` and by measurement
V25 — see step 1. The revised rule:

```go
const miscompileStepName = "Measure compiler reproducer (platform-conditional, gated)"

// indentOf returns the count of leading spaces. ci.yml contains no tabs (V9),
// no anchors/aliases, no flow-style steps, no `if:` keys, and no block-scalar
// indentation indicators (V9, V13, V19).
func indentOf(s string) int { return len(s) - len(strings.TrimLeft(s, " ")) }
```

1. **stepCol (the anchor, derived — SHALLOWEST enclosing `steps:`, per ruling `D-WORLD-30`)**:
   from the (unique — the existing `count==1` fatal guards this) line `i` containing
   `miscompileReproducerPath`, scan EVERY line from `i` backward; among all lines whose
   *trimmed* text is exactly `steps:`, keep the one with the SMALLEST indent (`indentOf`),
   breaking a tie at that minimal indent toward the nearest such line (the enclosing one; with
   every step dash in today's file at column 6 — V9 — the tie-break cannot change `stepCol`,
   but the rule must still be total). `stepCol` is then the indent of the first line AFTER that
   anchor whose trimmed text has prefix `- `. EITHER lookup failing — no `steps:` line above
   `i`, or no following dash line — is a `t.Fatalf("instrument failure: …")`, never a skip and
   never a pass. **Why shallowest, not nearest**: a `steps:`-shaped decoy line inside the
   guarded step's own `run: |` block scalar is NEARER to `i` than the real anchor, so under
   revision 1's nearest rule the decoy captures the anchor — measured V25: on that attack the
   nearest anchor derives `stepCol=12`, `start=179`, `end=180` and fails to adjudicate the step
   (Invariant A violated → loud FATAL on the controller's construction; GREEN/blind on the
   ledger's iter-142 layout — V25's honesty note carries both readings) — while block-scalar
   content must sit strictly DEEPER than its step's own column (V19), so no in-scalar decoy can
   ever be the shallowest candidate. Under the shallowest anchor the same mutant yields
   `stepCol=6`, `start=174`, `end=184`, RED on the step's own flag line (V25). V9's measured
   "every dash in ci.yml sits at column 6" is a *corroboration* of this derivation, not the
   rule itself.
1b. **Column pin (`gpt5-6-sol` round-3 verbatim fix, applied under the narrow-refinement
   carve-out; ruling fidelity, not a direction change)**: immediately after deriving `stepCol`,
   cross-check it against a declared constant and refuse loudly on any difference —
   `const expectedStepCol = 6`, then
   `t.Fatalf("instrument failure: derived step column %d; update expectedStepCol after an intentional ci.yml re-indent", stepCol)`.
   **Why both, rather than one or the other.** `D-WORLD-30` names two forms of hardening in one
   breath — it CHOOSES the shallowest derivation and then accepts a residual phrased in the
   constant form (*"a future re-indent of `ci.yml` makes the gate fatal until the constant is
   updated"*). Revision 2 first read that second clause as belonging only to the unadopted
   alternative; round 3's `gpt5-6-sol` rejected that reading, correctly — the loop does not get
   to decide which half of a human's sentence was meant. Implementing BOTH satisfies the ruling
   literally, is strictly conservative (it can only add a loud refusal, never a pass), and needs
   no amended ruling. It is also non-vacuous by measurement rather than by argument: **V31**
   re-indents the whole `go-verify` `steps:` block two columns right and the derivation alone
   does NOT stay silent — but the reason is an accident (see Residual 8), which is exactly why
   the pin is worth its one line.

2. **start**: from `j := i` walking backward — **starting at `j == i`**, so a dash-line
   identifying line (the `- run:` shape, V16) matches itself — the nearest line whose trimmed
   text has prefix `- ` (ANY first key, not just `name:`) **and whose indent is exactly
   `stepCol`**. Run-block content cannot legally occupy that column (V19), so no decoy at ANY
   depth can be `start`; ARM A's `- timeout-minutes:` first line qualifies, so key order does
   not matter — a claim that is now TRUE of this rule, including the `- run:`-as-dash-key
   layout, and is backed by measurement rather than asserted: the legitimate `- run:`-first
   reorder locates `start == i` and stays GREEN, its flagged variant reds on the step's own
   flag line (V16).
3. **end**: from `start+1`, the first line that is either a dash line at indent exactly
   `stepCol` (the next sibling step) or any non-blank, non-comment line at indent strictly
   less than `stepCol` (end of the steps sequence); EOF otherwise. Blank lines and comment
   lines (trimmed prefix `#`) never terminate the block. **Precondition (declared, Residual
   3)**: a non-blank non-comment NON-dash line at exactly `stepCol` would be silently included;
   valid YAML cannot produce one under `steps:` (measured REJECTED by the oracle, V19).
4. **Invariant A (containment)**: `start <= i < end`, else
   `t.Fatalf("instrument failure: located step block [%d,%d) does not contain the identifying line %d", …)`.
   The row's candidate fix (b), kept as a tripwire — see honesty note below.
5. **Invariant B (identity)**: the located block must contain a line whose trimmed text is
   exactly `- name: `+miscompileStepName or `name: `+miscompileStepName, else
   `t.Fatalf("instrument failure: located block is not the miscompile step …")`. A
   misattributed range becomes a LOUD instrument failure, never a silently shifted scan.
6. The existing `continue-on-error` loop then runs over `[start,end)` unchanged, as does
   everything after it (the run.sh assertions are untouched).

Estimated diff: ~40 lines changed in one function plus one small helper and one const, all in
`host/verifygate/toolchain_pin_gate_test.go`. The names `indentOf`, `miscompileStepName`, and
`stepCol` are unallocated in the repo today (V12 **including its 2026-09-01 addendum**:
`oc-glm-5-2` caught revision 1's citation FALSIFIED in round 2 — V12's original greps never
searched `indentOf` — and the addendum repairs it by WIDENING the grep to all three names at HEAD
`f5cfb1b`, 0 hits each against a firing 5-hit control, rather than by deleting the sentence).

### Why the round-1 relative rule was rejected — measured, not asserted (the quorum's kill)

Round 1's start test was `trimmed prefix "- " AND indent < indentOf(lines[i])`. Two
counterexamples, both implemented and run through a harness carrying the round-1 rule and the
absolute rule verbatim (V15, V16):

- **`gemini-3-1-pro`'s attack (V15)**: indent the identifying line DEEPER (path at column 14
  inside the `run: |` scalar), place a decoy `- decoy item` at an intermediate column 12 above
  it, and duplicate the pinned step name below the decoy so Invariant B is satisfied. Under the
  relative rule: `start` = the decoy (12 < 14), both invariants PASS, the live
  `continue-on-error` at column 8 sits ABOVE `start` and is excluded — **GREEN with the
  forbidden flag live: ARM B's fail-open, fully recreated**. Round 1 closed the one instance,
  not the class. Under the absolute rule the decoy can never be `start` (12 ≠ 6) and the arm
  reds on the step's own flag line. (Honestly noted: the HEAD scan incidentally catches this
  particular mutant — its `- name:`-prefix walk skips a `- decoy item` line — so MUT-I is the
  kill of the round-1 *proposal*, not additional evidence against HEAD.)
- **The controller's `- run:` dash-key shape (V16)**: when `run:` is the dash line's own key,
  the identifying line IS the dash line at `stepCol`, and no dash exists at a smaller indent —
  the relative rule dies FATAL "start not found" even on the *legitimate*, flag-free layout.
  Loud, which is the acceptable direction, but round 1's "key order no longer matters" claim
  was false and the shape appeared in neither the Conflict Surface nor the Declared Residuals.
  The absolute rule's `j == i` start case handles it: `start == i`, GREEN when clean, RED on
  its flagged variant.

### Why (b) alone was rejected — measured, not asserted

The row's candidate fix (b) — "assert the located block CONTAINS the `run:` line that
identified it" — **kills NEITHER measured arm on its own**. In ARM B the misattributed block
*does* contain the identifying line (`start` = the decoy, which sits *above* the path; the
narrowing cut off the flag, not the path). In ARM A the misattributed range is a superset and
trivially contains it. So (b) alone leaves both measured failure modes intact. It survives here
only as Invariant A: two lines of defense-in-depth against locator shapes not yet imagined,
with its reachability honestly declared nil today (§Mutation Drill, MUT-H row).

Invariant B is what converts residual mis-scoping into loud failure: under the OLD scan it
would fire on both arms (ARM A's block starts at the wrong step's name; ARM B's block starts
at `- name: documentation example`). Under the new locator it is the tripwire for anything
that defeats the indentation rule.

### What this is NOT

- **Not a YAML parser — now by attended ruling, not by the loop's scoping judgment.** The
  honest question "does this fix need a real YAML parse?" was put to the human after round 2
  blocked on exactly this direction, and `D-WORLD-30` answers NO, for a stated threat-model
  reason: *"This gate is a repo-internal tripwire against our own accidental regressions, not
  an adversarial boundary: the threat model that would justify B is an actor who can already
  edit `ci.yml`, and such an actor can simply delete the test. Not worth adding the second
  direct dependency to a `go.mod` that has exactly one."* (The "exactly one" is measured: V11.)
  The repo has ZERO yaml dependencies (V11) and the sibling scans in the same file are
  deliberately line-based with their limits doc-commented
  (`TestGoToolchainPinsAgreeAndMatchJobList` at `:105-109` says so verbatim). Ruby's YAML
  parser is used in this doc and in the mutation drill **only as a measurement oracle on this
  rig**, never inside the gate.
- **Not an actionlint gate.** `actionlint` exists at `/opt/homebrew/bin/actionlint` on this
  rig but appears NOWHERE in any repo gate (V8 — grep of `.github/` and `scripts/` empty with
  a firing same-file control; there is no Makefile, V8). Row 41's V18 stands at HEAD as a
  statement about repo gates. Wiring a rig-local binary into CI acceptance is exactly the
  rig-dependence this mission was burned by at queue row 58; this design does not propose it,
  and says so rather than leaving it implicit. (actionlint also would not catch either arm:
  both mutants are VALID workflow YAML — V2, V3 ruby views.)

## Alternatives Considered (and why rejected)

1. **Fix (b) alone** — rejected by measurement: kills neither ARM A nor ARM B (§above).
2. **Real YAML parse (new go.mod dependency) — "Option B" of `D-WORLD-30`** — **REJECTED BY
   THE HUMAN** (attended ruling, 2026-09-01), no longer merely by this doc's scoping judgment.
   The ruling's reason, verbatim: *"This gate is a repo-internal tripwire against our own
   accidental regressions, not an adversarial boundary: the threat model that would justify B
   is an actor who can already edit `ci.yml`, and such an actor can simply delete the test.
   Not worth adding the second direct dependency to a `go.mod` that has exactly one."* A
   reviewer objection re-proposing this direction is answered by the ruling, not re-litigated
   by the loop; only Mark can reverse it.
3. **actionlint in CI** — rig-dependent + would not detect either arm (both mutants are valid
   YAML); rejected and named (§What this is NOT).
4. **Byte-pin the step's exact three lines** — would red on ANY legitimate edit of the step
   (timeout change, comment) and re-create the round-1 class (a gate that reds against its own
   documentation). Rejected.
5. **Leave as-is** ("it failed safe when found") — refuted this iteration: V3 is fail-OPEN on
   the exact flag this gate exists for, and V2 is a contract-violating false positive.

## Acceptance Criteria

All commands run from the repo root with `export PATH=/opt/homebrew/bin:$PATH`. The gate
command for the changed package is `go vet ./host/verifygate/` — **`go build
./host/verifygate/` is rc=1 on PRISTINE dev** ("no non-test Go files", V10; test-only
package), so a build-rc=0 AC would be broken at base. `./scripts/verify_go.sh` rc=0 is
**forbidden as an AC** on this rig (measured FLAKY at iter-141: 4 runs, 2 red / 2 green, 3
different failing sets, CI green throughout — recorded in
`w-inventory-test-blind-to-asymmetric-addition.md` V6 and the iter-141 log); the targeted
package commands below are the acceptance vehicle. **Every test AC carries the
`AILANG_BIN=/tmp/ailang-v0300/ailang` export, and any AC of the form "`go test
./host/verifygate/` passes" without it is BROKEN AS WRITTEN**: bare `go test
./host/verifygate/ -count=1` is rc=1 at base with 17 FAILs, ALL of them "AILANG_BIN is unset"
instrument-failure refusals — a deliberate anti-false-green guard, not a defect (baselined at
HEAD `f5cfb1b`, V26; with the export the same command is rc=0 in 67s).

- **AC1 (base green, scope-reach proven)**:
  `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1 -v`
  → rc=0 with `=== RUN TestMiscompileInstrumentStepIsGatedInCI` present in the output (the
  `-v` RUN line is the proof the pattern reaches the test — V1 shows both at base; a green
  from a run that never entered the function proves nothing).
- **AC2 (the fail-open family dies)**: with the V3 mutant applied verbatim to ci.yml, the AC1
  command → **rc=1**, and the failure message names a line that lies inside the miscompile
  step's own block *as the ruby YAML oracle sees it* (the step's `continue-on-error` line),
  not a line of any other step. Same for the V15 gemini-attack mutant (MUT-I, sha `2f795e56…`):
  **rc=1** on the step's own flag line. Restore ci.yml byte-identical (sha256 `aed8e186…`)
  after each.
- **AC3 (the false-positive/key-order family closes)**: with the V2 mutant applied verbatim,
  the AC1 command → **rc=0** (a flag on an unrelated step is that step's business — quorum
  round-2 R1); and with the legitimate `- run:`-first mutant (MUT-K, sha `58a2dd93…`) applied,
  the AC1 command → **rc=0** (key order, including the dash-key shape, is none of this test's
  business). Restore byte-identical after each.
- **AC4 (the original kill still fires, in every layout)**: with `continue-on-error: true`
  added to the miscompile step in its CANONICAL layout (MUT-C), the AC1 command → rc=1; and
  with the flagged `- run:`-first mutant (MUT-J, sha `96f51f64…`) applied, the AC1 command →
  rc=1 naming the step's own flag line.
- **AC5 (identity invariant live, proven by mutation, not grep)**: with the miscompile step
  renamed in ci.yml (MUT-G), the AC1 command → rc=1 with an `instrument failure` message from
  Invariant B; then, with Invariant B neutered in the test source via `if false && …`
  (MUT-N — neutered, not deleted, so a compile error cannot masquerade as the guard firing),
  the same MUT-G mutant → rc=0. Both edits reverted afterwards. (A `grep -c` on the invariant
  cannot prove it is live — iter-141's correction; only this flip does.)
- **AC6 (package hygiene + package-wide green)**: `go vet ./host/verifygate/` rc=0
  (rc-correct at base, V10 and V26) and `gofmt -l host/verifygate/` prints nothing
  (**baselined at `f5cfb1b`: rc=0, zero lines, with a firing known-positive control that DOES
  name the file when it is deliberately misformatted — V29**); and
  `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1` → rc=0 (baselined
  rc=0, 67s, at `f5cfb1b` — V26). The `AILANG_BIN` export is load-bearing, not decoration:
  without it the same command is rc=1 with 17 loud "AILANG_BIN is unset" refusals (V26), so a
  bare-command version of this AC would be red at base and unsatisfiable.
- **AC7 (no collateral)**: `.github/workflows/ci.yml` byte-identical at landing (sha256
  `aed8e186fb57036eb6b03509cbb668850d577d46e6cc68b30e7a4c042108ed85`) and
  `git status --porcelain` shows only the intended test-file (and doc) changes. The full
  verify gate **`AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` remains rc=0** (no
  `.ail` files are touched). **The export is load-bearing and this AC was BROKEN AS WRITTEN
  until round 3: the bare command is rc=1 on the pristine tree on this rig** — the PATH `ailang`
  is a development build and the gate's own anti-false-green guard refuses it
  (`✗ AILANG_BIN refused [DEV_BUILD]: version token 'v0.34.0-247-g4bd58bef6-dirty'`), while with
  the pin it is rc=0 at 11 required identities / 40 named tests / 9-of-9 package steps (**V30**).
  Filed by `gemini-3-1-pro` at round 3 as an unbaselined AC; measured rather than forwarded, and
  its fear was correct for this half and unfounded for the `gofmt` half.
- **AC8 (the column pin is live and loud, proven by mutation)**: with the `go-verify` `steps:`
  block re-indented two columns right (MUT-Q), the AC1 command → rc=1 carrying an
  `instrument failure` message, never rc=0; and with step 1b's `t.Fatalf` neutered in the test
  source via `if false && …` (neutered, not deleted, so a compile error cannot masquerade as
  the guard firing) the SAME mutant must still red — via Invariant A or B — so this AC also
  records WHICH refusal fired. Both edits reverted afterwards, byte-identity asserted by
  sha256.

Every AC command above has a Verification Log row proving it rc-correct and scope-reaching on
the pristine tree (V1 for the targeted test command incl. the RUN line; V10 for vet; V26 for
the package-wide command with and without `AILANG_BIN`; **V29 for `gofmt`, with its own firing
control; V30 for `verify_ail.sh`, which is the one that was red at base until it carried the
pin**) — none is red at base as now written, none can green vacuously. The enumeration was
COMPLETED at round 3 rather than asserted: `gemini-3-1-pro` observed that two AC commands had
no baseline row at all, which was true, and one of the two was genuinely unsatisfiable.

## Mutation Drill

One named mutant per acceptance-relevant branch. Three arms per row: "OLD" = the scan at HEAD
(measured at `7077e455` in revision 1 and re-measured at `f5cfb1b` in revision 2 — ci.yml is
byte-identical across the two heads, V20, and the test function under change is line-for-line
unchanged with the file's 94 intervening insertions all landing AFTER it, V27); "R1" =
round-1's relative rule (SUPERSEDED — kept because the drill must show what killed it); "NEW" = the absolute
step-column locator **as RATIFIED by `D-WORLD-30`: anchor = the SHALLOWEST enclosing `steps:`**.
Revision 1's NEW cells were measured against the nearest-`steps:`-anchor variant (V15-V17); by
construction the two anchor derivations differ only when a `steps:`-shaped line sits between
the real anchor and `i` at a deeper indent — no revision-1 mutant contains one, so every
revision-1 NEW cell carries over unchanged, and V21 corroborates (both anchors prototyped on
ARM A: identical `stepCol`/`start`/`end`). The lone divergence is MUT-O (V25). R1/NEW cells
marked *measured* were run through the harness carrying the rules verbatim (V15-V17) or the
controller's revision-2 locator prototype (V21-V25); OLD cells cite V2/V3, the harness, or the
revision-2 landed-gate runs (V21, V23, V24). MUT-A/B/C/D/I/J/K/E/O/P mutate ci.yml (executor:
restored byte-identical after each, sha-checked); MUT-N mutates the test source (reverted).

**Addition coverage, honestly labeled** (controller obj-1: round 1 had none): "a removal proves
the check FIRES, only an addition proves it LOOKS" demands additions whose VERDICT
discriminates between arms. Round-1's MUT-D (rc=0 under both arms) and MUT-E (rc=1 under both)
moved counts, not verdicts. The discriminating additions are now **MUT-I** (adds a sequence
item inside the run scalar; GREEN under R1, RED under NEW — the mutant that killed round 1) and
the reworked **MUT-D** (adds a flagged stranger step; rc=1 false positive under OLD, rc=0 under
NEW). MUT-E is kept but labeled as what it is: a NON-discriminating guard of the uniqueness
precondition, not addition coverage.

| ID | Shape | Mutation (applied to pristine) | OLD (HEAD) | R1 relative (superseded) | NEW absolute expected |
|---|---|---|---|---|---|
| MUT-A | reorder + unrelated flag (= ARM A) | coe on `go build + test gate` step; miscompile keys reordered (`- timeout-minutes:` first); sha `28731ce9…` | **rc=1 at ci.yml:166 — FALSE POSITIVE, measured V2** | rc=0 *(measured)* | rc=0 — the false positive DIES |
| MUT-B | narrow-by-decoy (= ARM B) | coe live on miscompile step + nested `- name:` decoy inside its `run: \|` block; sha `66bad3af…` | **rc=0 PASS — FALSE NEGATIVE, measured V3** | rc=1 coe@176 *(measured)* | rc=1: reds on the step's own flag line |
| MUT-C | canonical re-introduction | coe added to miscompile step, canonical layout | rc=1 (the row-44 kill) | rc=1 *(measured)* | rc=1: same red, correct attribution |
| MUT-D | **ADDITION, discriminating (reworked)** | stranger step `- name: injected` with coe inserted immediately BEFORE the miscompile step, AND miscompile keys reordered (`- timeout-minutes:` first); sha `5f61658a…` | **rc=1 blaming the STRANGER's flag (coe@175) — false positive, measured V17** | rc=0 *(measured)* | rc=0 — the LOOKS proof with a verdict flip: the locator bounds the block at its true dash, so the added flagged stranger is (correctly) not blamed. Round-1's non-discriminating MUT-D (identical verdicts both arms) is retired. |
| MUT-E | **ADDITION, non-discriminating (uniqueness guard)** | second line containing the reproducer path inside another step | rc=1 fatal `count(...)=2` | rc=1 same fatal | rc=1 same fatal — guards the locator's uniqueness precondition before any scan runs; NOT addition coverage, and labeled so |
| MUT-F | reorder alone | miscompile keys reordered, NO coe anywhere | rc=0 (superset range, no flag) | rc=0 | rc=0 — key order is officially none of this test's business |
| MUT-I | **ADDITION — the round-1 kill (`gemini-3-1-pro`)** | coe live on miscompile step; `run: \|` scalar holds a `- decoy item` at column 12 + duplicated pinned name + the path at column 14; sha `2f795e56…` | rc=1 coe@176 *(measured — incidental: the `- name:` prefix walk skips a `- decoy item` dash)* | **rc=0 GREEN — BLIND; start=178 (the decoy), InvA=PASS, InvB=PASS, flag above start excluded — measured V15** | rc=1 coe@176, start=174 (the true dash): the class closes, not the instance |
| MUT-J | `- run:` dash-key + flag (controller obj-2) | miscompile step rewritten `- run:` first, `name:`/`timeout-minutes:`/coe following; sha `96f51f64…` | rc=1 coe@177 *(measured — but block misattributed: start=163, the unrelated step, superset luck)* | **FATAL `start not found (i=174)` — loud, but "key order no longer matters" was FALSE — measured V16** | rc=1 coe@177 with `start == i == 174` (the `j == i` case) |
| MUT-K | `- run:` dash-key, legitimate (no flag) | same rewrite, NO coe anywhere; sha `58a2dd93…` | rc=0 *(measured — superset luck)* | **FATAL — a legitimate layout reds, measured V16** | rc=0 with `start == i` — the row that makes the key-order claim TRUE |
| MUT-O | **ADDITION — the round-2 attack the RULING names (`steps:`-shaped decoy in the scalar)** | coe live on miscompile step; its `run: \|` scalar holds `steps:` (col 12) + `- name: decoy` below it, the reproducer path BELOW the heredoc — the V25 construction verbatim; **sha256 pinned: `72acc72d239056d7772b769b345212d092dce715a229202dee939899415f9f31`** (base `aed8e186…`; re-derived and pinned at round 3 on `oc-glm-5-2`'s note — V32); assert LANDED by matching that sha AND by the ruby view: step[6] coe=true | **rc=1, RED at `ci.yml:182` — the guarded step's OWN flag line, i.e. CAUGHT (measured V28, controller, this revision).** The landed gate catches this mutant *incidentally*: its backward scan anchors on the in-scalar `- name: decoy` line, which WIDENS the range to include the flag rather than narrowing it (the MUT-P mechanism). **So this attack does not defeat the code at HEAD — it defeats revision 1's PROPOSED nearest-anchor locator (V25). The ruling's shallowest anchor therefore prevents a regression this doc would have INTRODUCED, and does not close a live hole; the doc must not claim otherwise.** | not measured (superseded rule) | rc=1, RED on the step's own flag line — *measured on the locator prototype*: shallowest anchor → stepCol=6, start=174, end=184 (V25). Revision 1's NEAREST anchor on the same mutant mis-derives stepCol=12, start=179, end=180 → Invariant-A FATAL (V25; its honesty note also carries the ledger's GREEN/blind reading on a differently-laid-out decoy — both wrong outcomes, neither adjudicates the step). Observable: the failure message's `ci.yml:<line>` naming the step's own coe line, produced BY the locator (the coe loop runs over the locator's own `[start,end)`), not alongside it. Vacuous if the mutant fails to parse or the flag is not genuinely live — assert the ruby view BEFORE reading any verdict. Package-wide red set NOT MEASURED — do not write a `-skip` rc=0 criterion for this mutant. |
| MUT-P | position control for MUT-B (a refinement row 52 does not state) | the V23 decoy construct with the `- name: decoy` line placed BEFORE the run.sh line instead of after it (V24 construction); **sha256 pinned: `1917c413b1942f230da15dd026a432ed674e447e30e407263d2307f190d88e56`** (V32); assert LANDED by matching that sha AND by the ruby view | **rc=1 at "ci.yml:181" — CAUGHT, measured V24: the backward scan anchors ON the decoy and WIDENS the range instead of narrowing it** | not measured (superseded rule) | rc=1 on the step's own flag line — EXPECTED, not prototyped: a decoy at scalar depth can never be `start` (its indent ≠ stepCol; scalar content sits strictly deeper, V19). Observable: the `ci.yml:<line>` in the failure message, produced by the locator's coe loop over its own range. Vacuous if the decoy drifts to AFTER the path line — that is MUT-B's shape and flips the OLD verdict (V23 vs V24) — so assert the decoy's position in the ruby scalar-content view before reading any verdict. This row exists so the doc CANNOT claim the decoy construct always fails open: the decoy's POSITION relative to the identifying line decides it (V24). Package-wide red set NOT MEASURED — no `-skip` rc=0 criterion. |
| MUT-Q | **the re-indent arm `gpt5-6-sol` asked for — proves step 1b's column pin is not decoration** | the whole `go-verify` `steps:` block (key + every line under it) shifted two columns right; still valid YAML; assert LANDED by sha-differs from `aed8e186…` and by the ruby view still parsing to 10 steps | rc=0 GREEN — the LANDED gate silently absorbs the re-indent (measured V31) | not measured (superseded rule) | rc=1 `instrument failure`, and the row must record WHICH refusal fired. Measured on the locator prototype (V31): the NEAREST anchor returns `stepCol=8` and GREEN; the SHALLOWEST anchor returns `stepCol=6` from the OTHER JOB's `steps:` key, locates `[96,99)` and violates Invariant A → loud FATAL. So the derivation alone already fails loud on this shape, but for an accidental reason (Residual 8) — step 1b's `expectedStepCol` pin is what makes the refusal deliberate and self-explaining. Observable: the `instrument failure` message text, produced BY the locator. Vacuous if the mutant is not valid YAML — assert the ruby view before reading any verdict. Package-wide red set NOT MEASURED — no `-skip` rc=0 criterion for this mutant. |
| MUT-G | identity break | rename the miscompile step's `name:` value | rc=0 (OLD pins no name) | rc=1 Invariant B | rc=1 fatal: Invariant B `instrument failure` — a block that is not the miscompile step is loud, never silently scanned |
| MUT-N | test-side NEUTER (liveness of Invariant B) | `if false && …` on Invariant B's condition in the test source, then re-apply MUT-G | n/a | n/a | MUT-G flips rc=1 → rc=0, proving Invariant B was the firing guard; revert both |
| MUT-H | (declared, not demonstrable) | Invariant A (containment) | n/a | n/a | **No valid-YAML mutant reaches it under the NEW locator** — re-derived for the absolute rule, and STRONGER than round 1's version: `start` is found walking back from `i` (so `start <= i`), and both `end` terminators (a dash at `stepCol`; a non-blank non-comment line left of `stepCol`) are ALSO block-scalar terminators in the YAML data model itself (V19: a `- ` at `stepCol` after a `\|`/`\|N` header parses as a NEW STEP; content below the indicated indent is REJECTED), so no valid document can place a terminator inside `(start, i]` while `i` remains this step's content. Kept as 2 lines of defense-in-depth; non-reachability DECLARED, not faked by a mutation row. |

Executor protocol per mutant: assert the mutant LANDED (sha256 differs, matching the sha pinned
in its row where one is pinned), assert the intended effect against the SYSTEM'S OWN VIEW
(`ruby -ryaml` step/keys/coe listing, as in V2/V3/V15 — a measurement oracle on this rig, not a
gate), run the AC1 command, record rc + message, restore, re-run the pristine control. Two
standing rules for the drill: (1) MUT-B's fail-open class has a SECOND measured layout — V23's
trailing-flag construction, re-derived at `f5cfb1b` — and the sprint may use either, asserting
the layout in the ruby view first; (2) every arm's verdict is read from the TARGETED AC1
command — where a ci.yml mutant's red set under a package-wide run may exceed this one test, it
is unmeasured, and no arm may carry a `-skip <test>` rc=0 criterion in place of the targeted
run.

## Conflict Surface

**What else reads `.github/workflows/ci.yml` (complete census, V7 — 3 Go files, 28 hits;
corrected in revision: round 1's prose said "four files, 25 hits" while enumerating three —
`gpt5-6-sol` caught the contradiction, and re-measurement shows the ENUMERATION was complete
and the cardinality word plus the hit count were the errors; the census was not widened because
there is no fourth file to find — the scope was right, the number was wrong):**

- `host/verifygate/toolchain_pin_gate_test.go` — this function (the only step scoper, V5) and
  `TestGoToolchainPinsAgreeAndMatchJobList`, which enumerates jobs by an indentation-anchored
  regex (`^  ([a-z0-9-]+):$` at `:124`) — precedent in the SAME file for indentation-aware
  line scanning; untouched.
- `host/verifygate/ail_binary_gate_test.go` (`TestZ3PinDeclaredOnceAndInstalledInBothJobs`,
  `:669`) — whole-file `strings.Contains`/`strings.Count` with its own known-positive controls
  (V6); no step scoping; untouched.
- `host/runbook/runbook_stageb_test.go` (`:339-368`) — whole-file line-contains sweep across
  ci.yml + scripts with two same-call known-positive controls (V6); no step scoping; untouched.
- No non-Go reader gates on step structure (the workflow is otherwise read only by GitHub
  Actions itself).

**What a stricter scope could newly reject that is legitimate today — the round-1 lesson
applied in advance:**

- **Renaming the miscompile step** now reds (Invariant B) until the test's
  `miscompileStepName` const is updated in the same edit. This is a new, intentional coupling —
  the same class as the existing `miscompileReproducerPath` count==1 pin, and strictly
  narrower than the OLD behavior's silent mis-scoping. It is the one legitimate edit this
  design makes louder, and it is stated here so nobody discovers it as a surprise red.
- **The ratified anchor change (nearest → SHALLOWEST `steps:`, `D-WORLD-30`) tightens nothing
  legitimate.** The two derivations pick a different anchor only when a `steps:`-shaped line
  occupies block-scalar content between the real anchor and the identifying line — on every
  other controller-prototyped layout they agree exactly (V21: both anchors, identical numbers;
  V25: the lone divergence, on an attack mutant, not a legitimate edit). And because there is
  NO hardcoded column constant in the adopted form, a whole-file re-indent of ci.yml is
  absorbed by re-derivation rather than turned into a fatal-until-edited red — the
  constant-pin alternative the ruling also sketches would behave differently, and this doc
  does not adopt it (Declared Residual 7).
- Nothing else tightens: pristine ci.yml contains ZERO `continue-on-error` anywhere (V14,
  with a firing same-file control), so no legitimate current use is at risk; and the fix
  NARROWS the flag ban relative to the OLD scan's misattributed superset — MUT-A/MUT-D show
  flags on unrelated steps becoming definitively that step's business, which is the
  quorum-ratified contract, not a loosening. ARM A's re-derivation at `f5cfb1b` (V21) shows
  the false positive still live in the landed gate today — blaming the unrelated step's own
  flag line — and the ratified locator returning the correct GREEN, with V22 (same reorder,
  flag ON the guarded step → RED) proving that GREEN is attribution, not permissiveness.

**What this change cannot see** (unchanged surface, inherited by name): everything in the
row-44 doc's Declared Residual 2 — a step-level `if:` that never evaluates true, YAML
anchors/aliases, flow-style step syntax — plus computed/obfuscated text (residual 3). This fix
touches only HOW the step block is located between those residuals; it neither discharges nor
widens them (measured: ci.yml today contains none of those constructs — V9 — so their absence
is a fact about the current file, not a capability of the scan).

## Declared Residuals

1. **Still a line scan, not a parse.** `if:`/anchors/flow-style (row-44 Declared Residual 2)
   and computed text (Residual 3) remain invisible, verbatim inherited. This item's fix is
   orthogonal to them and does not claim otherwise.
2. **Block-scalar indentation indicators — RE-REASONED under the absolute anchor, not copied
   forward.** `stepCol` is now the load-bearing column, so round 1's sentence ("could in
   principle hold content at or left of the step column") was re-measured rather than
   inherited (V19, psych oracle): content below a scalar's indicated indent is NOT content —
   it terminates the scalar, and a document that then places a non-item line at a non-key
   column is REJECTED outright; a `- ` line at exactly `stepCol` after a `|`/`|N` header
   parses as a genuine NEXT STEP (nsteps increments, scalar ends); and the indicator's minimum
   places content deeper than the step's key column (> `stepCol`). So on this oracle, valid
   YAML cannot use `|N` to put scalar *content* at or left of `stepCol` — the residual NARROWS
   from "could confuse the locator" to a parser-authority gap: psych is the rig's oracle while
   GitHub Actions' own YAML parser is the authority, and their agreement on these edge shapes
   is assumed, not proven. No `|N`/`>N` exists in ci.yml today (V13: zero, against 7 plain
   `: |` blocks as the control). Named, narrowed, not closed.
3. **The `end` rule's precondition (from `oc-glm-5-2`'s round-1 note), named explicitly rather
   than left inside the generic "still a line scan" residual:** a non-blank, non-comment,
   NON-dash line at exactly `stepCol` matches neither terminator and would be silently
   included in the block. Under `steps:` such a line is not valid YAML — a sibling of sequence
   items must be an item — and the oracle REJECTS the probe (V19). It is a PRECONDITION of the
   end rule (holds for every valid workflow), not a handled case; a file that violates it is
   already broken to Actions before this scan runs.
4. **Invariant A is currently unreachable** (MUT-H, re-derived for the absolute rule): its
   liveness cannot be shown by a valid-YAML input mutation today, for the stronger reason in
   the MUT-H row (the end terminators are scalar terminators in the data model itself, V19).
   It is carried as declared defense-in-depth; if a future shape fires it, that firing is an
   instrument failure demanding a locator fix, not a nuisance to delete.
5. **The name pin trusts uniqueness informally.** Invariant B checks the located block
   contains the pinned name; it does not assert the name appears in ci.yml exactly once
   (the reproducer-path count==1 fatal already anchors uniqueness of the step this test
   cares about, so a duplicated *name* elsewhere cannot move the located block — but it
   would make Invariant B satisfiable by a block this test did not intend if the path ever
   moved there too, at which point count==1 reds first).
6. **Tabs.** `indentOf` counts spaces; a tab-indented ci.yml would defeat it. YAML forbids
   tabs in indentation and ci.yml contains none (V9), so a tab is a parse error to Actions
   before it is a problem for this scan — stated, not assumed silently.
7. **The ruling's accepted residual — quoted, and now implemented in BOTH of its clauses
   (revised at round 3).** `D-WORLD-30` accepts, verbatim: *"The residual is accepted and named:
   the scan stays text-level, and a future re-indent of `ci.yml` makes the gate fatal until the
   constant is updated."* **Revision 2 scoped the second clause away**, reading it as belonging
   to the constant-pin alternative the ruling also names and which this doc did not adopt.
   `gpt5-6-sol` rejected that at round 3 — *"Residual 7 reclassifies the ruling's constant clause
   as belonging to an alternative without evidence that the human made that distinction"* — and
   it is right: the loop does not get to choose which half of a human's sentence was operative.
   **The reclassification is WITHDRAWN.** Mechanism step 1b now implements the constant clause
   directly (`expectedStepCol = 6` + a loud `t.Fatalf`), on top of the shallowest derivation the
   ruling chose, so both clauses hold as written and no amended ruling is needed. What remains
   accepted, and is the residual proper, is the FIRST clause: the scan stays text-level —
   `gpt5-6-sol`'s general point that block-scalar content is author-controlled remains true in
   principle, and Residuals 1-3 name the text-level edges (`if:`/anchors/flow-style, computed
   text, the psych-vs-Actions parser-authority gap).

8. **The backward `steps:` scan is UNBOUNDED, so the shallowest anchor can leave the JOB — found
   by measuring the round-3 objection rather than by arguing it (V31), and it is a defect the
   RATIFIED rule INTRODUCES.** `ci.yml` declares two jobs, each with its own `steps:` key, both
   at indent 4 today (V9), and the derivation scans backward to the top of the file. Today the
   minimal-indent tie-break resolves to the enclosing one and `stepCol` is 6 either way (V20,
   V31 base arm). **Re-indent the `go-verify` `steps:` block two columns right and it does not:**
   the *other* job's `steps:` is then strictly shallower, the anchor jumps jobs, and the located
   block is `[96,99)` — the ailang-verify job — while the nearest anchor quietly returns
   `stepCol=8` and a GREEN (V31). **This is stated as a residual and not as a decision for the
   human, on a discriminator this mission already uses: the introduced defect fails LOUD, never
   silent.** Invariant B is the backstop and it is structural, not incidental — a located block
   belonging to another job cannot contain the pinned miscompile step name, so a cross-job
   capture is a `t.Fatalf` instrument failure by construction, not a green; on the measured
   re-indent it is in fact Invariant A that fires first (containment), which is the same
   outcome. Step 1b's column pin is the second, independent loud refusal on the same shape.
   **Honest limit on this claim:** the enumeration of shapes is NOT exhaustive — what is
   measured is one re-indent mutant and one structural argument, and no false green was found
   in any arm run this iteration (V20-V32). A bounded-to-the-enclosing-job scan would close the
   class outright; that is a mechanism change nobody has ratified, so it is named here as the
   follow-up rather than taken.

## Milestones

- **MS1 (~0.05d)** — rewrite the locator in `TestMiscompileInstrumentStepIsGatedInCI`:
  `indentOf` helper, `stepCol` derivation from the **SHALLOWEST enclosing `steps:` key
  (`D-WORLD-30`; Mechanism step 1, including the total tie-break)**, backward start scan
  from `j == i` at exactly `stepCol`, forward end scan, Invariants A + B,
  `miscompileStepName` const; update the function's doc-comment (the `DECLARED RESIDUAL`
  block's scoping sentence and its reference to the ROW-44 doc's V19 — not this doc's V19 —
  now describe the new mechanism, and the comment should name the ruling). AC1, AC6.
- **MS2 (~0.04d)** — mutation drill MUT-A through MUT-P per the table (including the two
  revision-2 arms: MUT-O, the ruling-named `steps:` decoy, whose OLD cell is owed at sprint
  time; and MUT-P, the position control), each with landed-assertion, ruby-oracle effect
  assertion, rc + message capture, byte-identical restore, pristine re-control. AC2-AC5, AC7.
- **MS3 (~0.01d)** — evidence write-up in the sprint notes; doc moves to `implemented/` at
  landing per repo flow.

Total ~0.1d. The real-YAML-parse dependency decision is **no longer open for the quorum to
route**: `D-WORLD-30` (attended) has taken it and REJECTED the dependency — *"Not worth adding
the second direct dependency to a `go.mod` that has exactly one"* (the one is measured, V11).
A reviewer objection in that direction is answered by quoting the ruling, not by widening or
re-parking the item; only Mark can reverse it.

## Verification Log

Rows V1-V19 measured by me at `dev` = `7077e455ec22a4779392cd4db9b72d4e0ae2df11` (single base
commit; `git status --porcelain` = 0 lines before and after the batch) on 2026-09-01, zsh,
`PATH=/opt/homebrew/bin:$PATH`. Rows V20-V26, V28-V32 and the V12 addendum are the CONTROLLER's,
first-party, at HEAD `f5cfb1b` (their batch protocol is in the revision-2 header below; ci.yml
is byte-identical across the two heads — V20). V27 is the designer's own read-only
cross-head consistency check (2026-09-02). Empty/negative results carry a same-scope
known-positive control in the same row.

```zsh
# V1 — pristine control (also AC1's command + scope-reach proof)
cp .github/workflows/ci.yml /tmp/ci_yml_iter142_mine.bak
shasum -a 256 .github/workflows/ci.yml
# observed: aed8e186fb57036eb6b03509cbb668850d577d46e6cc68b30e7a4c042108ed85
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ \
  -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1 -v
# observed: rc=0; `=== RUN   TestMiscompileInstrumentStepIsGatedInCI`; `--- PASS`; `ok … host/verifygate`
git status --porcelain | wc -l
# observed: 0

# V2 — ARM A: coe on the unrelated previous step + miscompile key reorder
perl -0pi -e 's/(        timeout-minutes: 25\n)(        run: \.\/scripts\/verify_go\.sh\n)/$1        continue-on-error: true\n$2/' .github/workflows/ci.yml
perl -0pi -e 's/      - name: Measure compiler reproducer \(platform-conditional, gated\)\n        timeout-minutes: 15\n/      - timeout-minutes: 15\n        name: Measure compiler reproducer \(platform-conditional, gated\)\n/' .github/workflows/ci.yml
shasum -a 256 .github/workflows/ci.yml
# observed: 28731ce9cd0b1e582804a40b7a377ae24aea9131f5d42669f09c1ea8b5b01626 (mutant LANDED)
ruby -ryaml -e '…'   # (step census of the job containing the reproducer; full one-liner in iter-142 notes)
# observed (system's own view): job=go-verify nsteps=10;
#   step[5] "go build + test gate…"          keys=["name","timeout-minutes","continue-on-error","run"] coe=true
#   step[6] "Measure compiler reproducer…"   keys=["timeout-minutes","name","run"]                     coe=nil
go vet ./host/verifygate/                     # observed: rc=0 (mutant tree still vets clean)
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1
# observed: rc=1 — `toolchain_pin_gate_test.go:531: ci.yml:166 re-introduces "continue-on-error"
#   in the miscompile step…`; and mutant line 166 IS `        continue-on-error: true` inside the
#   `go build + test gate` step (sed -n '162,178p' listing captured). FALSE POSITIVE on the
#   unrelated step, violating the test's own :502-505 contract.
cp /tmp/ci_yml_iter142_mine.bak .github/workflows/ci.yml && shasum -a 256 .github/workflows/ci.yml
# observed: aed8e186… (restored byte-identical)

# V3 — ARM B: coe live on the miscompile step + nested `- name:` decoy in its run block
perl -0pi -e "s/      - name: Measure compiler reproducer \\(platform-conditional, gated\\)\n        timeout-minutes: 15\n        run: \\.\\/design_docs\\/verification\\/w-race-gate-blindspot\\/run\\.sh\n/      - name: Measure compiler reproducer (platform-conditional, gated)\n        timeout-minutes: 15\n        continue-on-error: true\n        run: |\n          cat <<'YAMLDOC'\n          - name: documentation example\n          YAMLDOC\n          .\\/design_docs\\/verification\\/w-race-gate-blindspot\\/run.sh\n/" .github/workflows/ci.yml
shasum -a 256 .github/workflows/ci.yml
# observed: 66bad3afc51faeb60bac1fbd844f1165a5a821e6dde450cb7d5145d6facbd215 (mutant LANDED)
ruby -ryaml -e '…'
# observed (system's own view): nsteps=10; the miscompile step has
#   keys=["name","timeout-minutes","continue-on-error","run"] coe=true; its run scalar CONTAINS
#   the `- name: documentation example` line as CONTENT (printed verbatim).
go vet ./host/verifygate/                     # observed: rc=0
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1 -v
# observed: rc=0, `--- PASS`. THE GATE FAILS OPEN with the forbidden flag LIVE on the guarded step.

# V4 — restore + post-batch pristine control
cp /tmp/ci_yml_iter142_mine.bak .github/workflows/ci.yml && shasum -a 256 .github/workflows/ci.yml
# observed: aed8e186… (byte-identical)
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1
# observed: rc=0, `ok … host/verifygate`;  git status --porcelain | wc -l → 0

# V5 — the step-scoping defect has ONE call site (corrects the row's "three places")
grep -rn --include='*.go' 'HasPrefix(strings.TrimSpace' .
# observed, exactly 3 lines:
#   host/verifygate/toolchain_pin_gate_test.go:511  (backward/start scan)   ← same function
#   host/verifygate/toolchain_pin_gate_test.go:524  (forward/end scan)      ← same function
#   host/store/writer_lock.go:194                    (SQLite `busy_timeout` pragma — unrelated, read in context)
# same-scope known-positive control (proves the pattern-space is visible):
grep -c 'strings.HasPrefix' host/verifygate/toolchain_pin_gate_test.go
# observed: 6

# V6 — the other ci.yml readers do NOT scope steps (read, not inferred)
sed -n '655,690p' host/verifygate/ail_binary_gate_test.go
# observed: TestZ3PinDeclaredOnceAndInstalledInBothJobs — whole-file strings.Contains/Count over
#   src, with its own known-positive controls ("ailang-verify:", "go-verify:", "./scripts/verify_go.sh")
sed -n '330,375p' host/runbook/runbook_stageb_test.go
# observed: whole-file line-contains sweep over ci.yml + scripts/*.sh, with TWO same-call
#   known-positive controls (verify_go.sh in ci.yml; verify_world_package.sh in verify_ail.sh)

# V7 — complete census of Go readers of ci.yml (the conflict surface)
# CORRECTED IN REVISION: round 1's row said "25 hits across exactly 4 files" over a 3-file
# enumeration. Re-measured from scratch (revision pass, same head):
grep -rln --include='*.go' 'ci\.yml' .
# observed, exactly 3 files: host/verifygate/ail_binary_gate_test.go,
#   host/verifygate/toolchain_pin_gate_test.go, host/runbook/runbook_stageb_test.go
grep -rn --include='*.go' 'ci\.yml' . | wc -l
# observed: 28
grep -rln --include='*.go' 'package verifygate' . | wc -l
# observed: 7 (same-scope known-positive control — the grep family sees this tree)
# The enumeration was COMPLETE; the prose cardinality and the hit count were the errors.

# V8 — no actionlint in any repo gate (row 41 V18 re-verified at THIS head); no Makefile exists
ls Makefile                                    # observed: rc=1, "No such file or directory"
grep -rn "actionlint" .github/ scripts/        # observed: rc=1, zero hits
grep -c "verify_go.sh" .github/workflows/ci.yml
# observed: 1  (same-file firing control — the grep sees this scope)
command -v actionlint
# observed: /opt/homebrew/bin/actionlint — a RIG fact, not a repo fact; not used by this design.

# V9 — ci.yml structural facts the indentation rule relies on
grep -nE '^( *)- ' .github/workflows/ci.yml | (indent census)
# observed: 17 dash lines, ALL at exactly 6 leading spaces; zero at any other column
grep -nE '&[A-Za-z]|\*[A-Za-z]|steps: *\[' .github/workflows/ci.yml   # observed: rc=1 (no anchors/aliases/flow steps)
grep -nE '^\s+if:' .github/workflows/ci.yml                            # observed: rc=1 (no step-level if:)
grep -nP '\t' .github/workflows/ci.yml                                 # observed: rc=1 (no tabs)
# firing control for this family of negatives: the dash census above returns 17 positives in the same file.
# miscompile step at PRISTINE lines 174-176; `go build + test gate` step name at line 163.

# V10 — the compile-gate trap, measured (why AC6 is vet, not build)
go build ./host/verifygate/
# observed: rc=1 — "no non-test Go files in …/host/verifygate" ON PRISTINE dev
go vet ./host/verifygate/
# observed: rc=0

# V11 — negative-existence: no YAML library dependency anywhere in go.mod
grep -in "yaml" go.mod                         # observed: rc=1, zero hits
grep -c "^require\|modernc.org/sqlite" go.mod  # observed: 2 (same-file firing control; one direct dep)

# V12 — negative-existence: proposed identifiers are unallocated
grep -rn --include='*.go' 'locateStepBlock\|miscompileStepName\|stepIndent' .
# observed: rc=1, zero hits
grep -c 'miscompileReproducerPath' host/verifygate/toolchain_pin_gate_test.go
# observed: 5 (same-file firing control)
# revision addendum — the absolute rule's new identifier:
grep -rn --include='*.go' 'stepCol\|stepDashCol' .   # observed: rc=1, zero hits (same control: 5)
# ── ADDENDUM 2026-09-01 (revision 2; controller-measured at HEAD f5cfb1b) ──
# oc-glm-5-2 (round 2) caught the citation FALSIFIED: the greps above never search `indentOf`,
# yet the Mechanism prose cited V12 for it. Repaired by WIDENING the grep and re-measuring per
# name at f5cfb1b — same scope, both directions controlled — not by deleting the sentence:
#   grep -rn --include='*.go' 'indentOf' .                    → observed: 0 hits
#   grep -rn --include='*.go' 'miscompileStepName' .          → observed: 0 hits
#   grep -rn --include='*.go' 'stepCol' .                     → observed: 0 hits
#   grep -rn --include='*.go' 'miscompileReproducerPath' .    → observed: 5 hits
#     (known-positive control, same scope — the grep family sees this tree)
#   negative control — a fresh, known-absent literal, same scope → observed: 0 hits
# The revision-1 sentence's CLAIM survives at this head; its CITATION did not. The original
# rows above stand as what was measured (and believed sufficient) at 7077e455.

# V13 — negative-existence: no block-scalar indentation indicators in ci.yml (Residual 2's shape)
grep -nE '\|[0-9]|>[0-9]' .github/workflows/ci.yml   # observed: rc=1, zero hits
grep -cE ': \|$' .github/workflows/ci.yml            # observed: 7 (plain `: |` blocks — firing control)
# also: zero dash-shaped content lines inside any run: | block today (awk sweep, zero output)

# V14 — negative-existence: pristine ci.yml has ZERO continue-on-error anywhere
grep -n "continue-on-error" .github/workflows/ci.yml # observed: rc=1, zero hits
grep -c "timeout-minutes" .github/workflows/ci.yml   # observed: 4 (same-file firing control)

# ═══ REVISION-PASS ROWS (2026-09-01, same head 7077e455; the repo's ci.yml was NEVER touched
# in this pass — every mutant is a /tmp copy, and git status --porcelain shows only this doc).
# Harness: /tmp/iter142rev/locator.rb — implements the HEAD scan ("head"), round-1's relative
# rule ("rel"), and this doc's absolute rule ("abs") verbatim from §Solution Design; prints
# rule, stepCol, start/end/i with indents, Invariants A/B, and the coe verdict, 1-based lines.
# Pristine harness control (all three rules agree on the untouched file):
#   head/rel/abs on pristine.yml → start=174(ind 6) end=178 i=176 InvA=PASS InvB=PASS GREEN
#   mutC.yml (canonical coe)     → RED coe@176 under all three rules

# V15 — gemini-3-1-pro's round-1 attack: coe live + deep-indent decoy dash + duplicated
# pinned name inside the run scalar (identifying path at column 14, decoy `- decoy item` at 12)
shasum -a 256 gemini.yml   # observed: 2f795e56… (mutant LANDED)
ruby -ryaml (step census)  # observed: keys=["name","timeout-minutes","continue-on-error","run"]
#   coe=true, nsteps=10; decoy dash line AND duplicated name line are run-scalar CONTENT
ruby locator.rb rel gemini.yml
# observed: rule=rel start=178(ind 12) end=182 i=180(ind 14) InvA=PASS InvB=PASS GREEN
#   ← BLIND. The live flag (column 8, line 176) sits ABOVE start and is excluded; both
#   invariants pass. ARM B's fail-open, recreated against the round-1 rule.
ruby locator.rb abs gemini.yml
# observed: rule=abs stepCol=6 start=174(ind 6) end=182 InvA=PASS InvB=PASS RED coe@176 ← caught
ruby locator.rb head gemini.yml
# observed: RED coe@176 — HEAD incidentally catches this mutant (its `- name:` prefix walk
#   skips a `- decoy item` dash); the attack is the counterexample to the round-1 PROPOSAL.

# V16 — controller obj-2: `- run:` as the dash-line key (flagged + legitimate variants)
shasum: dashrun.yml 96f51f64… (with coe), dashrun_legit.yml 58a2dd93… (no coe) — both LANDED
ruby -ryaml: keys=["run","name","timeout-minutes"(,"continue-on-error")], nsteps=10 both
ruby locator.rb rel dashrun.yml        # observed: FATAL start not found (i=174)
ruby locator.rb rel dashrun_legit.yml  # observed: FATAL start not found (i=174) — a
#   LEGITIMATE flag-free layout reds under the round-1 rule; its key-order claim was false
ruby locator.rb abs dashrun.yml
# observed: stepCol=6 start=174(ind 6) end=179 i=174 InvA=PASS InvB=PASS RED coe@177
#   ← the j==i start case: the identifying line IS the step dash and matches itself
ruby locator.rb abs dashrun_legit.yml  # observed: start=174 == i, GREEN — key order free
ruby locator.rb head dashrun.yml       # observed: RED coe@177 but start=163 — HEAD blames the
#   right line only via a block misattributed to begin at the UNRELATED previous step

# V17 — MUT-D reworked into a DISCRIMINATING addition: flagged stranger step injected BEFORE
# the miscompile step + miscompile keys reordered (`- timeout-minutes:` first)
shasum -a 256 mutDp.yml    # observed: 5f61658a… (mutant LANDED)
ruby -ryaml                # observed: 11 steps; step[6] "injected" coe=true; miscompile coe=nil
ruby locator.rb head mutDp.yml
# observed: start=174 (the INJECTED step) RED coe@175 — FALSE POSITIVE blaming the stranger
ruby locator.rb abs mutDp.yml
# observed: stepCol=6 start=178 (the miscompile's own reordered dash) GREEN — verdict flips

# V18 — census re-derived (gpt5-6-sol's round-1 objection). Commands + observed output live in
# the CORRECTED V7 row above: exactly 3 files, 28 hits, known-positive control 7. The round-1
# enumeration was complete; "four" and "25" were the errors. Not widened — no fourth file exists.

# V19 — block-scalar / end-rule preconditions probed (Residuals 2 & 3; MUT-H's reasoning).
# ruby -ryaml on a steps:-shaped document (steps at col 4→ items at col 6, keys at col 8):
#   run: |1 + content at col 9              → PARSED; content IS the scalar (firing control)
#   run: |1 + content at col 7              → REJECTED (below indicated indent: not content,
#                                             and not valid anywhere else either)
#   run: |1 + `- x` at col 6                → PARSED as a NEW STEP (nsteps=3, scalar empty)
#   run: |1 + col-9 content, then `- x` col 6 → scalar="ok\n", nsteps=3 — a dash at stepCol
#                                             TERMINATES the scalar in the data model itself
#   run: |  + plain `x` at col 6            → REJECTED (a non-dash line at stepCol is not
#                                             valid YAML under steps: — the end rule's
#                                             precondition holds for every valid workflow)
# Oracle is psych on this rig; GitHub Actions' parser is the authority — the remaining gap is
# named in Residual 2.

# ═══ REVISION-2 ROWS (2026-09-01, CONTROLLER-measured at HEAD f5cfb1b, handed to this revision
# with commands + observed output). Rig protocol for the whole batch: a scratch worktree that is
# a SIBLING of the repo (never /tmp); every mutant asserted LANDED by sha256 before its result
# was read; the intended effect asserted against ruby -ryaml — the system's own view of ci.yml —
# rather than against bytes; `go vet ./host/verifygate/` rc=0 read BEFORE any test result;
# ci.yml restored byte-identical from a `cp` backup; the pristine control green before AND after
# the batch (porcelain 0). "landed gate" = the AC1 command against the test at HEAD;
# "ratified locator" = the controller's prototype of §Mechanism as revised (V21 ran BOTH anchor
# derivations; V25 contrasts them). Mutant shas for V21-V25 live in the controller's iter-147
# notes and are not transcribed here; the sprint re-pins its own.

# V20 — revision-2 base control at f5cfb1b
shasum -a 256 .github/workflows/ci.yml
# observed: aed8e186fb57036eb6b03509cbb668850d577d46e6cc68b30e7a4c042108ed85 — IDENTICAL to
#   V1's base at 7077e455: both batches measured the same bytes.
# ruby -ryaml step census, observed: 10 steps; step[5]="go build + test gate…" coe=nil;
#   step[6]="Measure compiler reproducer…" coe=nil
# landed gate (AC1 command), observed: rc=0, RUN=1 PASS=1

# V21 — ARM A re-derived at f5cfb1b (row 52's false-positive claim, NOT inherited from V2):
# coe inserted on the PREVIOUS, unrelated `go build + test gate` step + miscompile keys
# reordered so `- name:` is not first (`- timeout-minutes: 15` / `  name: Measure compiler
# reproducer …`). An independent construction, not V2's 28731ce9 mutant byte-for-byte — its
# flag line lands at 164 where V2's landed at 166; the verdict class is identical.
# ruby view, observed: step[5] coe=true, step[6] coe=nil, step count unchanged at 10
# landed gate, observed: rc=1 RUN=1 PASS=0 FAIL=1 —
#   `ci.yml:164 re-introduces "continue-on-error" in the miscompile step`
#   — and mutant line 164 IS the unrelated step's own flag line. The gate blames the wrong
#   step: the FALSE POSITIVE is live in the landed gate TODAY.
# ratified locator (prototyped, BOTH anchors — nearest and shallowest agree on this layout):
#   stepCol=6, start=175, end=179, verdict GREEN — the CORRECT outcome: a flag on an unrelated
#   step is that step's business (row-44 quorum round-2 R1).

# V22 — ARM A2, the discriminating control for V21: the SAME key reorder, but the flag placed
# ON the miscompile step itself.
# ruby view, observed: step[6] coe=true
# ratified locator, observed: stepCol=6, start=174, end=179, RED at the flag line
# — so V21's new GREEN is step-attribution, not permissiveness.

# V23 — ARM B re-derived at f5cfb1b (row 52's fail-open claim, NOT inherited from V3):
# coe live on the guarded step AND its run: converted to a block scalar whose content contains
# a step-shaped line AFTER the reproducer path:
#     run: |
#       ./design_docs/verification/w-race-gate-blindspot/run.sh
#       cat <<'EOF' > /dev/null
#       - name: decoy
#       EOF
#     continue-on-error: true
# ruby view, observed: step[6] coe=true — the forbidden flag is genuinely live on the guarded
#   step (a different layout from V3's, whose decoy sat between the flag and the path; two
#   measured layouts of the same fail-open class).
# landed gate, observed: rc=0 RUN=1 PASS=1 — a clean green over a live forbidden flag. The
#   FAIL-OPEN is live in the landed gate TODAY.
# ratified locator, observed: stepCol=6, start=174, end=183, RED at the flag line.

# V24 — ARM B0, the POSITION control for V23 (a refinement queue row 52 does not state): the
# SAME decoy construct with the `- name: decoy` line placed BEFORE the run.sh line instead of
# after it.
# landed gate, observed: rc=1 RUN=1 PASS=0 FAIL=1, "ci.yml:181" — still CAUGHT. The reason:
#   the landed backward scan then anchors ON the decoy, which WIDENS the scanned range instead
#   of narrowing it. So "a `- name:`-shaped line inside the block scalar" is NOT sufficient to
#   fail the gate open — the decoy's POSITION relative to the identifying line decides it.
#   This doc therefore claims fail-open only for the measured layouts (V3, V23), never for the
#   construct as such.

# V25 — ARM C: the round-2 attack the RULING names — a `steps:`-shaped decoy inside the
# guarded step's own block scalar, flag live:
#     run: |
#       cat <<'EOF' > /dev/null
#       steps:
#         - name: decoy
#       EOF
#       ./design_docs/verification/w-race-gate-blindspot/run.sh
#     continue-on-error: true
# ruby view, observed: step[6] coe=true
# NEAREST anchor (revision 1's rule), observed: stepCol=12, start=179, end=180 — Invariant A
#   (containment start <= i < end) is VIOLATED, so on this construction it is a LOUD t.Fatalf
#   instrument failure, not a silent green.
# SHALLOWEST anchor (the ratified rule), observed: stepCol=6, start=174, end=184, RED at the
#   flag line — the only derivation that ADJUDICATES the step.
# HONESTY NOTE — two readings of the nearest-anchor outcome, neither overwritten: the
#   D-WORLD-30 ledger characterises this attack under the nearest anchor as "both invariants
#   PASS, flag excluded → GREEN/blind", measured at iteration 142 on a DIFFERENTLY-LAID-OUT
#   decoy; the controller's construction above instead produces the loud FATAL. Both are wrong
#   outcomes — a blind green and an unadjudicated fatal alike fail to rule on the step — and
#   only the shallowest anchor produces the correct RED. Each reading stands as measured on
#   its own layout.

# V26 — AC-command baselines at f5cfb1b, pristine worktree (why every test AC carries AILANG_BIN)
go test ./host/verifygate/ -count=1
# observed: rc=1 — 17 FAIL, ALL of them "AILANG_BIN is unset" instrument-failure refusals: a
#   deliberate anti-false-green guard, not a defect. Any AC of the form
#   "`go test ./host/verifygate/` passes" is broken as written.
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1
# observed: rc=0, ok (67s)
go vet ./host/verifygate/
# observed: rc=0 (read before any test result, per batch protocol)

# V27 — designer-measured (this revision, 2026-09-02, read-only): the V1-V19 head and the
# V20-V26 head measured the SAME code under change. ci.yml is byte-identical (V20's sha = V1's).
# The test file itself is NOT byte-identical — but every inserted line lands AFTER the function:
git diff --stat 7077e455..f5cfb1b -- host/verifygate/toolchain_pin_gate_test.go
# observed: 1 file changed, 94 insertions(+)
git diff 7077e455..f5cfb1b -- host/verifygate/toolchain_pin_gate_test.go | grep -n "^@@"
# observed: one hunk, `@@ -560,3 +560,97 @@` — appended after line 560, below the function
git show f5cfb1b:host/verifygate/toolchain_pin_gate_test.go | grep -n 'func TestMiscompileInstrumentStepIsGatedInCI\|HasPrefix(strings.TrimSpace'
# observed: 492 (func), 511 (backward/start scan), 524 (forward/end scan) — identical line
#   numbers to the 7077e455 view (positive control: the same grep against
#   `git show 7077e455:…` puts the func at 492). Every `:492-533`/`:511`/`:524`/`:502-505`
#   citation in this doc therefore holds at BOTH heads.
```

```bash
# V28 — controller-measured (this revision, 2026-09-02, HEAD f5cfb1b, scratch worktree that is a
# SIBLING of the repo, never /tmp): the OLD (landed, HEAD) verdict for MUT-O, which V25 left owed.
# Mutant = the V25 `steps:`-decoy construction; asserted LANDED by sha256 (aed8e186… → 72acc72d…),
# effect asserted against the system's own view before any verdict was read, `go vet` rc=0 read
# BEFORE the test result, restored byte-identical from a `cp` backup, pristine control green after.
ruby -ryaml -e '...steps[6]["continue-on-error"]...' .github/workflows/ci.yml
# observed: n=10, step[6].coe=true — the forbidden flag is genuinely live on the GUARDED step
go vet ./host/verifygate/                                   # observed: rc=0 (read before the test)
go test ./host/verifygate/ -run TestMiscompileInstrumentStepIsGatedInCI -count=1 -v
# observed: rc=1 RUN=1 PASS=0 FAIL=1 —
#   toolchain_pin_gate_test.go:531: ci.yml:182 re-introduces "continue-on-error" in the miscompile step
awk 'NR==182' .github/workflows/ci.yml                      # observed: "        continue-on-error: true"
#   i.e. the named line IS the guarded step's own flag line: the landed gate CATCHES this mutant.
# WHY it catches, and why that is not reassuring: the landed backward scan anchors on the in-scalar
#   `- name: decoy` line, so the located range WIDENS past the flag (the MUT-P mechanism, V24) —
#   the same accident that makes MUT-P caught and MUT-B (V23) fail open. It is not adjudication.
# CONSEQUENCE, stated so no downstream reader can over-read the ruling: the round-2 `steps:` attack
#   defeats revision 1's PROPOSED nearest-anchor locator (V25), NOT the code at HEAD. The ratified
#   shallowest anchor prevents a regression this doc would have introduced; it does not close a
#   live hole. The live holes are V21 (false positive) and V23 (fail-open), both still open at HEAD.
# restore control: sha256 back to aed8e186…; pristine gate rc=0; worktree porcelain clean of any
#   non-doc change.
```

```bash
# V29 — AC baseline `gemini-3-1-pro` asked for (round 3), gofmt half. HEAD f5cfb1b, pristine
# sibling worktree. The reviewer's fear was UNFOUNDED here — and the row exists anyway, because
# "it is probably clean" is not a baseline.
gofmt -l host/verifygate/                      # observed: rc=0, 0 bytes, 0 lines
# known-positive control, SAME scope, same call: append a deliberately misformatted func to
# host/verifygate/toolchain_pin_gate_test.go, re-run, then restore:
gofmt -l host/verifygate/                      # observed: rc=0, 1 line ->
#   host/verifygate/toolchain_pin_gate_test.go
git diff --name-only -- host/verifygate/toolchain_pin_gate_test.go   # observed: 0 (restored)
# So the empty result above is a measurement, not a broken instrument.

# V30 — AC baseline, verify_ail.sh half. The reviewer's fear was CORRECT here: AC7 was
# UNSATISFIABLE as written on this rig.
./scripts/verify_ail.sh                        # observed: rc=1
#   "✗ AILANG_BIN refused [DEV_BUILD]: version token 'v0.34.0-247-g4bd58bef6-dirty'
#    identifies a development build"
#   i.e. the PATH `ailang` on this rig is a dev build and the gate's own anti-false-green guard
#   refuses it. This is the gate working, not a defect — and it makes a bare-command AC red at
#   base, which is exactly the class rule "baseline every acceptance command" exists for.
AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh          # observed: rc=0
#   "✓ verify gate PASSED: 11 required identities verified, 40 named tests pass"
#   "✓ world package gate PASSED: 9/9 steps performed non-zero work"
#   "✓ compiler pinned by exact bytes: AILANG v0.30.0 on Darwin/arm64"
# AC7 now carries the export. Two arms, both codes captured without a pipe.

# V31 — the re-indent arm (`gpt5-6-sol` round-3 objection, MEASURED rather than forwarded),
# and the finding that produced Declared Residual 8. Mutant: the whole `go-verify` `steps:`
# block shifted two columns right. LANDED asserted by sha256 (aed8e186… -> 4fd378c4…);
# restored byte-identical from a `cp` backup; pristine control rc=0 after.
#   locator prototype, pristine file, BOTH anchors:
#     nearest    stepCol=6  start=174 end=178 -> GREEN
#     shallowest stepCol=6  start=174 end=178 -> GREEN      (they agree at base)
#   locator prototype, RE-INDENTED file:
#     nearest    stepCol=8  start=174 end=178 -> GREEN      <- silently absorbed
#     shallowest stepCol=6  start=96  end=99  -> FATAL (Invariant A: containment)
#   the landed gate at HEAD on the same mutant:  rc=0        <- silently absorbed
# READ IT CAREFULLY, because it cuts two ways. (a) It ANSWERS the objection: under the ratified
#   shallowest anchor a re-indent does NOT pass silently — it is a loud instrument failure, which
#   is the outcome the ruling's residual sentence demands. (b) It also shows WHY that is not
#   good enough on its own: `stepCol=6` there comes from the OTHER JOB's `steps:` key (both jobs
#   declare one; the backward scan is unbounded), so the loud refusal is an accident of this
#   file's shape rather than a designed behaviour. Hence step 1b's explicit column pin AND
#   Declared Residual 8. NOT MEASURED: whether any shape makes a cross-job capture pass GREEN —
#   Invariant B forbids it structurally (another job's block cannot contain the pinned step
#   name), but no mutant proving that was run this iteration.

# V32 — mutant sha256 pins (`oc-glm-5-2` round-3 note, applied). Each construction re-derived
# from its documented recipe in the pristine sibling worktree and hashed; base restored
# byte-identical afterwards and the pristine gate re-run rc=0.
# base ci.yml : aed8e186fb57036eb6b03509cbb668850d577d46e6cc68b30e7a4c042108ed85
# MUT-O (V25) : 72acc72d239056d7772b769b345212d092dce715a229202dee939899415f9f31
# MUT-P (V24) : 1917c413b1942f230da15dd026a432ed674e447e30e407263d2307f190d88e56
# MUT-B (V23) : c1903a869701e47626a4fb1cd537ffb42a935d02fa1ace26b2f39f330a6cecbf
# The sprint asserts LANDED by MATCHING these, not merely by "sha differs" — a regex or heredoc
# that lands a subtly different mutant is then a loud mismatch instead of a silent pass.
```

## Out of Scope

- The `run.sh` assertions in the same test (kernel-read pins, platform-token bans) — untouched.
- Row-44 Declared Residuals 1, 3, 4, 5 (darwin-blindness, computed text, saw_good, upstream-fix
  detection) — a different surface each; this item neither discharges nor worsens them.
- Any change to `.github/workflows/ci.yml` itself (AC7 pins it byte-identical).
- The rig-local `actionlint` binary and any repo-gate wiring for it (named refusal, §Solution
  Design; a future row may take it up as its own decision).
- The other two ci.yml readers (V6) — already carrying their own controls, no step scoping.

## Related Documents

- `design_docs/world-mission.md` — ledger row **`D-WORLD-30` (RESOLVED, attended, 2026-09-01)**:
  the ratified ruling this revision propagates — Option A (line scan, shallowest-anchor
  hardening) adopted, Option B (YAML parse dependency) rejected by the human; the ruling text
  is quoted verbatim in the Status header, and its iter-142 nearest-anchor characterisation is
  carried (not overwritten) in V25's honesty note.
- `design_docs/implemented/w-miscompile-instrument-inert-in-ci.md` — row 44: the test this item
  sharpens; its quorum round-2 R1 is the scoping contract ARM A violates; its Declared
  Residual 2 names the adjacent constructs this item explicitly does NOT absorb.
- `design_docs/implemented/w-setup-go-pin-unguarded.md` — row 41: V18 (no actionlint gate),
  re-verified here as V8.
- `design_docs/planned/w-inventory-test-blind-to-asymmetric-addition.md` — row 51 (iter-141):
  the `verify_go.sh`-is-flaky measurement (its V6) that forbids rc=0 ACs on it, and the
  "a grep cannot prove an assertion is live" correction that shapes AC5/MUT-N here.
- `design_docs/world-mission.md` queue row 52 — the item; corrected by V2/V3 (exploitable both
  directions) and V5 (one call site, not three).

**DESIGN_DOC_PATH**: `design_docs/planned/w-wiring-test-step-scoping-imprecise-under-key-reorder.md`

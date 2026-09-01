# w-wiring-test-step-scoping-imprecise-under-key-reorder

**Status**: DRAFT, REVISION 1 — quorum round 1 BLOCKED at FULL STRENGTH
(`.synthesis.absent_reviewers` = `[]`; verdicts: `gpt5-6-sol` reject, `gemini-3-1-pro` reject,
`oc-glm-5-2` pass, controller reject). What changed and why, all by measurement, not concession:
the round-1 RELATIVE indentation locator is replaced by an ABSOLUTE step-column rule —
`gemini-3-1-pro`'s deep-indent-decoy attack recreates ARM B's fail-open under the relative rule
(measured GREEN/blind, V15) and the controller's `- run:` dash-key shape makes its "key order no
longer matters" claim false (measured FATAL, V16); both are correctly handled under the absolute
rule (V15, V16). `gpt5-6-sol`'s census contradiction is corrected: the enumeration was COMPLETE,
the cardinality word was the error (V7 corrected, three files). The mutation drill now carries
genuinely discriminating addition coverage (controller obj-1; MUT-I, MUT-D). This is the one
protocol-mandated revision; resubmitted to quorum. (Unattended mission loop: the quorum ALWAYS
runs; none of the four attended triggers fires.)
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

All measurements were made by me in this checkout at `dev` = `7077e455` (pristine,
`git status --porcelain` = 0 lines) on 2026-09-01, bracketed by pristine controls; ci.yml was
restored byte-identical (sha256 `aed8e186…` before, between, and after — V1, V4). The
controller's iteration-142 measurements are independently reproduced here; every number below
is my own. Full commands and observed output: §Verification Log.

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

### Mechanism: ABSOLUTE step-column location + two loud invariants (revised after quorum R1)

Rewrite the locator inside `TestMiscompileInstrumentStepIsGatedInCI` (the ONE call site, V5) to
anchor on the **absolute column YAML gives block-sequence items under the enclosing `steps:`
key** — derived from the file, never hardcoded. Round 1 proposed a *relative* test
(`indent(dash) < indent(identifying line)`); that rule is superseded by measurement — see §Why
the round-1 relative rule was rejected. The revised rule:

```go
const miscompileStepName = "Measure compiler reproducer (platform-conditional, gated)"

// indentOf returns the count of leading spaces. ci.yml contains no tabs (V9),
// no anchors/aliases, no flow-style steps, no `if:` keys, and no block-scalar
// indentation indicators (V9, V13, V19).
func indentOf(s string) int { return len(s) - len(strings.TrimLeft(s, " ")) }
```

1. **stepCol (the anchor, derived)**: from the (unique — the existing `count==1` fatal guards
   this) line `i` containing `miscompileReproducerPath`, walk backward to the nearest line whose
   *trimmed* text is exactly `steps:`; `stepCol` is the indent of the first following line whose
   trimmed text has prefix `- `. Either lookup failing is a
   `t.Fatalf("instrument failure: …")`. V9's measured "every dash in ci.yml sits at column 6"
   is now a *corroboration* of this derivation, not the rule itself.
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
`stepCol` are unallocated in the repo today (V12).

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

- **Not a YAML parser.** The honest question "does this fix need a real YAML parse?" is
  answered: no, at this item's size. The repo has ZERO yaml dependencies (V11) and the sibling
  scans in the same file are deliberately line-based with their limits doc-commented
  (`TestGoToolchainPinsAgreeAndMatchJobList` at `:105-109` says so verbatim). Adding a YAML
  module dependency to sharpen one test's scoping is a dependency decision above a ~0.1d row;
  if a future row needs semantic YAML reads in more places, THAT row should carry the decision
  explicitly. Ruby's YAML parser is used in this doc and in the mutation drill **only as a
  measurement oracle on this rig**, never inside the gate.
- **Not an actionlint gate.** `actionlint` exists at `/opt/homebrew/bin/actionlint` on this
  rig but appears NOWHERE in any repo gate (V8 — grep of `.github/` and `scripts/` empty with
  a firing same-file control; there is no Makefile, V8). Row 41's V18 stands at HEAD as a
  statement about repo gates. Wiring a rig-local binary into CI acceptance is exactly the
  rig-dependence this mission was burned by at queue row 58; this design does not propose it,
  and says so rather than leaving it implicit. (actionlint also would not catch either arm:
  both mutants are VALID workflow YAML — V2, V3 ruby views.)

## Alternatives Considered (and why rejected)

1. **Fix (b) alone** — rejected by measurement: kills neither ARM A nor ARM B (§above).
2. **Real YAML parse (new dependency)** — rejected at this scope; scoped as an explicit future
   decision, not assumed (§What this is NOT).
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
package commands below are the acceptance vehicle.

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
- **AC6 (package hygiene)**: `go vet ./host/verifygate/` rc=0 (rc-correct at base, V10) and
  `gofmt -l host/verifygate/` prints nothing.
- **AC7 (no collateral)**: `.github/workflows/ci.yml` byte-identical at landing (sha256
  `aed8e186fb57036eb6b03509cbb668850d577d46e6cc68b30e7a4c042108ed85`) and
  `git status --porcelain` shows only the intended test-file (and doc) changes. The full
  verify gate `./scripts/verify_ail.sh` remains rc=0 (no `.ail` files are touched).

Every AC command above has a Verification Log row proving it rc-correct and scope-reaching on
the pristine tree (V1 for the test command incl. the RUN line; V10 for vet) — none is red at
base, none can green vacuously.

## Mutation Drill

One named mutant per acceptance-relevant branch. Three arms per row: "OLD" = the scan at HEAD
`7077e455`; "R1" = round-1's relative rule (SUPERSEDED — kept because the drill must show what
killed it); "NEW" = the absolute step-column locator. R1/NEW cells marked *measured* were run
through the harness carrying both rules verbatim (V15-V17); OLD cells cite V2/V3 or the
harness. MUT-A/B/C/D/I/J/K/E mutate ci.yml (executor: restored byte-identical after each,
sha-checked); MUT-N mutates the test source (reverted).

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
| MUT-G | identity break | rename the miscompile step's `name:` value | rc=0 (OLD pins no name) | rc=1 Invariant B | rc=1 fatal: Invariant B `instrument failure` — a block that is not the miscompile step is loud, never silently scanned |
| MUT-N | test-side NEUTER (liveness of Invariant B) | `if false && …` on Invariant B's condition in the test source, then re-apply MUT-G | n/a | n/a | MUT-G flips rc=1 → rc=0, proving Invariant B was the firing guard; revert both |
| MUT-H | (declared, not demonstrable) | Invariant A (containment) | n/a | n/a | **No valid-YAML mutant reaches it under the NEW locator** — re-derived for the absolute rule, and STRONGER than round 1's version: `start` is found walking back from `i` (so `start <= i`), and both `end` terminators (a dash at `stepCol`; a non-blank non-comment line left of `stepCol`) are ALSO block-scalar terminators in the YAML data model itself (V19: a `- ` at `stepCol` after a `\|`/`\|N` header parses as a NEW STEP; content below the indicated indent is REJECTED), so no valid document can place a terminator inside `(start, i]` while `i` remains this step's content. Kept as 2 lines of defense-in-depth; non-reachability DECLARED, not faked by a mutation row. |

Executor protocol per mutant: assert the mutant LANDED (sha256 differs, matching the sha pinned
in its row where one is pinned), assert the intended effect against the SYSTEM'S OWN VIEW
(`ruby -ryaml` step/keys/coe listing, as in V2/V3/V15 — a measurement oracle on this rig, not a
gate), run the AC1 command, record rc + message, restore, re-run the pristine control.

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
- Nothing else tightens: pristine ci.yml contains ZERO `continue-on-error` anywhere (V14,
  with a firing same-file control), so no legitimate current use is at risk; and the fix
  NARROWS the flag ban relative to the OLD scan's misattributed superset — MUT-A/MUT-D show
  flags on unrelated steps becoming definitively that step's business, which is the
  quorum-ratified contract, not a loosening.

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

## Milestones

- **MS1 (~0.05d)** — rewrite the locator in `TestMiscompileInstrumentStepIsGatedInCI`:
  `indentOf` helper, `stepCol` derivation from the enclosing `steps:` key, backward start scan
  from `j == i` at exactly `stepCol`, forward end scan, Invariants A + B,
  `miscompileStepName` const; update the function's doc-comment (the `DECLARED RESIDUAL`
  block's scoping sentence and its reference to the ROW-44 doc's V19 — not this doc's V19 —
  now describe the new mechanism). AC1, AC6.
- **MS2 (~0.04d)** — mutation drill MUT-A through MUT-N per the table, each with
  landed-assertion, ruby-oracle effect assertion, rc + message capture, byte-identical
  restore, pristine re-control. AC2-AC5, AC7.
- **MS3 (~0.01d)** — evidence write-up in the sprint notes; doc moves to `implemented/` at
  landing per repo flow.

Total ~0.1d. If quorum objects in a direction that requires a real YAML parse, that is a
**dependency decision** to route explicitly (new go.mod entry in a repo with one direct
dependency, V11) — widen the item consciously or park it; do not absorb it silently.

## Verification Log

All rows measured by me at `dev` = `7077e455ec22a4779392cd4db9b72d4e0ae2df11` (single base
commit; `git status --porcelain` = 0 lines before and after the batch) on 2026-09-01, zsh,
`PATH=/opt/homebrew/bin:$PATH`. Empty/negative results carry a same-scope known-positive
control in the same row.

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

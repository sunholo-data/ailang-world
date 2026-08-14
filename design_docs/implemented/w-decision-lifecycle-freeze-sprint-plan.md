# Sprint plan — `w-decision-lifecycle-freeze` (queue item 15, clause-5)

**Item**: queue item 15 — *the ratified three-constructor `TimeoutPolicy` set, the frozen v1
`DecisionPacket` world type, and five Z3-proven lifecycle laws in `world/types.ail`.*
**Status**: PLANNED · **NO SPLIT of the landing** · **4 milestones / 1 commit**
**Design doc**: [`design_docs/planned/w-decision-lifecycle-freeze.md`](w-decision-lifecycle-freeze.md)
(694 lines; round-2 revision, both external reviewers' objections measured and closed)
**Base**: `dev` @ `b9a7838`, clean worktree (`git status --porcelain` empty, V0)
**Worktree**: `/Users/voightkampff/dev/sunholo-data/.wt-iter85`, branch
`sprint/w-decision-lifecycle-freeze`. A **sibling of the repo, never `/tmp`** for the executor's
tree: `host/verifygate` and `host/boundary` derive `repoRoot` from `runtime.Caller` and copy live
trees, so a relocated checkout reds for its location rather than for the code.
**Planner**: mission-control iteration 85, opus sprint-planner, first-party measurement on this rig.

**Ratification carried into this plan** (attended, Mark, 2026-08-14, recorded in the charter at
`b9a7838`): HUMAN-SURFACE §7 point 3 = **OPTION A — FREEZE v1 NOW**; the semantic ID
`world/decision-packet/v1` is reserved; every later amendment is a NEW version, never an in-place
edit; **enforcement remains item 7's by explicit deferral.** §7 point 1 (the three-constructor set
and §2.1's resolution semantics) is taken as ratified **by inclusion** — `DecisionPacket` has a
`policy: TimeoutPolicy` field and all five laws are defined over that type, so Option A is not
executable without it. §7.4 of this plan states why that inference is safe here and what would
make it unsafe.

---

## 0. Headline: this sprint is pre-registered, not predicted

**Everything in M1 was composed, applied, verified, tested, mutated and gate-run end to end by the
planner this session, on the exact bytes M1 lands.** §9 is the full verification log. The
load-bearing outputs:

| Artifact | Pre-registered value | Evidence |
|---|---|---|
| `world/types.ail` after M1 | sha256 **`f8e08f68734724e4525cf2ccfff0dd98053d45ba54e5e6c47d48e4e1f3f1f454`** (132 → 279 lines) | V4 |
| `EXACT_TOTAL_VERIFIED` | **10** — measured per-module: contracts 1, logepoch 2, transitions 1, types **6** | V6 |
| `EXACT_TOTAL_TESTS` (`len(tests[])`) | **39** · `passed_tests` observed **43** — do NOT pin 43 | V5 |
| The 19 new named identities | exactly `outcomeCode_test_1..6`, `deferCode_test_1..3`, `scheduleCode_test_1..3`, `firedCode_test_1..3`, `escalationCode_test_1..4` | V5 |
| New ready-packet golden | the byte string in §3/M1.5, computed through `host/pkgproj` behind **two** known-positive controls | V9, V10, V11 |
| Full pinned AILANG gate with all six edits | **rc=0** · `✓ 10/10 required world/ identities verified across 11 module(s)` · `✓ all 39 required named tests pass` · `9/9 steps` | V12 |
| `go test ./host/verifygate/ -run TestModuleManifest` with the `10/10` marker | **ok**, 18.4 s | V13 |

**The executor's job is transcription plus a 16-arm mutation drill, not authoring.** If any
pre-registered value does not reproduce, that is a **STOP-and-report**, not a number to force.

---

## 1. Planner's first-party verification, and what it REFUTED

Every controller premise and every load-bearing doc claim was re-derived at `b9a7838` before
anything was planned. Controller premises: **all confirmed**. The design doc: **three defects
found, two of which are CI-red or stall-shaped.**

### 1.1 Controller premises — CONFIRMED

| # | Premise | Verdict |
|---|---|---|
| C1 | Worktree at `b9a7838`, branch `sprint/w-decision-lifecycle-freeze`, clean | **CONFIRMED** (V0) |
| C2 | `/tmp/ailang-v0300/ailang` = AILANG v0.30.0 (`e37b370`); PATH `ailang` is `-dirty` | **CONFIRMED** (V1). No claim in this plan uses the PATH binary. |
| C3 | Five gate pins at `verify_ail.sh:135 / :262-267 / :310 / :333 / :342` and `module_manifest_gate_test.go:128`; spaces around `=` on the python pins | **CONFIRMED, every line number exact** (V2). `EXACT_TOTAL_MODULES` does not exist: grep **0**, same-instrument control `EXACT_TOTAL_VERIFIED` **4 hits** — the instrument fires. |
| C4 | AC4's totals are subject to observed runner JSON; the gate reads `len(tests[])`, never `passed_tests` | **CONFIRMED AND MEASURED**: `len(tests[])=39`, `passed_tests=43` on the exact landing bytes (V5). `verify_ail.sh:367` is `n = len(tests)`; `passed_tests` appears nowhere in the parser. |
| C5 | Doc is fresh from `bc8f193`; 0 non-`design_docs` files moved | **CONFIRMED** (V3), and independently: every §11 line citation this plan relies on resolves to the exact construct it names. |
| C6 | `host/broker` carries a ~18% base flake; `GOTOOLCHAIN=go1.25.6` is a base condition | **CARRIED, not re-litigated.** One green `host/broker` run is not evidence the flake is gone. |

### 1.2 Three defects in the design doc — each MEASURED, each changing the plan

#### (i) **AC9 is UNSATISFIABLE as written. `interfaceHash` CANNOT move.** *(plan-killing)*

Doc §4.3: *"content/interface/tarball hashes all move."* Doc **AC9**: *"Fails on hand-authored JSON
or **unchanged hashes**."*

`host/pkgproj/pkgproj.go:86` — `InterfaceHash(manifest Manifest)` hashes **only** the package name,
edition, `ailang` constraint, sorted export **module** names, and sorted effects. **It never opens a
source file.** Adding exported *symbols* to an already-exported *module* is invisible to it. Item 15
adds no export module (`exports` stays the same four), so `interfaceHash` is invariant by
construction.

Measured (V11), with the known-positive control run **first** (V9/V10: the unmodified projection
reproduces the committed golden byte-for-byte, and does so path-independently):

| field | committed golden | after this item | verdict |
|---|---|---|---|
| `contentHash` | `sha256:489d5e5d47d5…` | `sha256:06acbb83ce88…` | **MOVES** |
| `interfaceHash` | `sha256:d16cc88270ff…` | `sha256:d16cc88270ff…` | **UNCHANGED** |
| `tarballSHA256` | `sha256:d0cdf42be80e…` | `sha256:5823edcfbb3f…` | **MOVES** |
| `tarballBytes` | `6236` | `7856` | **MOVES** |

An executor obeying AC9 literally has two exits and both are disasters: declare a correct
implementation failed, or hand-edit the golden to "make interfaceHash differ" — which AC9's own
next clause forbids and which step 9/9's `cmp -s` reds anyway. **AC9 is AMENDED in §2.**

**This is a repeated defect, and that is the finding.** Item 13 hit exactly this; its *sprint plan*
(`design_docs/implemented/w-evidence-grade-mapping-sprint-plan.md` §0.2 i) diagnosed and amended it.
But the correction was **never written back into the design doc**: `w-evidence-grade-mapping.md:291`
still says the change moves `interfaceHash`, and its `:387` still carries the unsatisfiable AC. Item
15's designer cites that document in §Related as its precedent and inherited the defect verbatim. A
correction recorded only in the artifact that consumed the doc does not reach the next reader of the
doc.

#### (ii) **A SIXTH pin the Conflict Surface omits: `verify_ail.sh:378`.** *(unguarded lie-in-waiting)*

```bash
378: echo "✓ verify gate PASSED: 5 required identities verified, 20 named tests pass"
```

Lines 315 and 372 interpolate `$EXACT_TOTAL_VERIFIED` / `EXACT_TOTAL_TESTS`; **378 is a hardcoded
literal.** Doc §8.1 enumerates its conflict surface exhaustively — *"fact E — all five listed,
including the unmoved one"* — and 378 is outside every entry.

Left alone the gate still exits **0**, so nothing catches it (V14): the only Go assertion on that
line is `host/verifygate/ail_binary_gate_test.go:118`, `passedMarker = "verify gate PASSED"` — a
**substring**, used via `strings.Contains`, which a stale `5 … 20` satisfies. Same-instrument
control: `grep -rn '5/5' --include='*.go'` finds `module_manifest_gate_test.go:128` (the guarded
pin), so the sweep fires; `grep -rn '5 required identities'` finds **only `verify_ail.sh:378`** —
nothing else pins its digits.

Consequence: the gate's own terminal summary — the line the controller's baseline quotes and every
future iteration reads — would announce `5 required identities verified, 20 named tests pass` while
verifying 10 and running 39. **This is pin 6.** It gets a grep-based criterion with a control
(**AC13**), not an invented mutation arm.

This is the *second* consecutive item to trip it: item 13's plan raised it as `:376` and added its
own AC13; the line moved to `:378` and the design doc's pin map lost it again.

#### (iii) **§4.2 claims to be "Exact AILANG code" and omits ~60 lines carrying all 19 test pins.**

§4.2's fenced block gives the three types and five laws. The five private adapters — `outcomeCode`,
`deferCode`, `scheduleCode`, `firedCode`, `escalationCode`, which *define every one of the 19 named
identities AC4 pins* — are described in prose only: *"whose exact bodies are in the V-P7/V-P13
probes."* Those probes live under `/tmp`. They happen to have survived (`/tmp/iso-item15-r2/`,
sha `fda7f30b…` reproduced, V7) — but a sprint whose acceptance criteria depend on a `/tmp`
directory surviving is one `tmpreaper` away from an executor re-deriving 19 identity names by guess.

**Consequence:** §3/M1.1 of this plan reproduces the complete 129-line addition **verbatim and
pre-registered by sha**, so the sprint no longer depends on any `/tmp` artifact.

### 1.3 Two comment amendments the ratification forces (planner-authored, flagged not smuggled)

The doc's §4.2 comment text is landed verbatim **except** for two lines, both direct transcriptions
of Mark's ratified decision as handed to this planner:

1. `-- PROPOSED pending ratification; the three members are the round-2 reviewers' candidates.`
   → `-- RATIFIED (attended, Mark, 2026-08-14): exactly these three members.`
   Landing "PROPOSED pending ratification" into the file being frozen as ratified v1 would commit a
   false comment.
2. The `DecisionPacket` comment gains `semantic ID world/decision-packet/v1` — that reservation is
   the operative content of Option A and belongs on the frozen type.

Both are inside the sha `f8e08f68…` that everything downstream is pre-registered against. If the
controller or quorum wants different comment text, the golden in §3/M1.5 must be recomputed — the
plan says so rather than letting a "harmless comment tweak" red step 9/9.

### 1.4 Three traps for the executor — not refutations

1. **`passed_tests` is 43, not 39.** The package leg (step 5/9) prints its own activity count
   (observed **42 activity lines** after the change) and has **no count pin**. An executor that
   "fixes" `EXACT_TOTAL_TESTS` to 43 or 42 reds the gate. The oracle is `len(tests[])`, and
   `verify_ail.sh:367` is the only place it is read.
2. **`ai-check` exits 0 on a Z3 encoding error.** Measured again this session on the exact landing
   bytes (V15/MU5): `check.passed=true`, **rc=0**, `verify.errors=1`, `timeoutOutcome status=error`.
   **No criterion in this sprint may read an exit code to prove totality.**
3. **`RecomputeReadyPacket` needs `Version`, and the repo's only manifest literal omits it.**
   `verify_world_package.sh:175` builds a `pkgproj.Manifest` with **no `Version`** — invisible
   there, because that helper only prints hashes. A regeneration helper copied from `:175` verbatim
   emits `"version":""` and reds step 9/9 with a *hash-shaped* diff for a *manifest-shaped* reason.
   §3/M1.5's helper sets `Version: "0.1.0"` and is proven by the reproduce-the-committed-golden
   control.

---

## 2. Acceptance criteria — doc AC1–AC12 adopted, with one amendment and one addition

Every command below was **baselined on the pristine tree at `b9a7838` first** (rule 3e). An AC
already red at base is a broken AC, not a defect in the change.

### Amendment

| AC | Amendment | Forced by |
|---|---|---|
| **AC9** | **`interfaceHash` MUST be byte-IDENTICAL**, not different. Require exactly **three** moved packet fields (`contentHash`, `tarballSHA256`, `tarballBytes`) and assert `interfaceHash` **unchanged** — that invariance is the positive evidence that no package export module was added, which is the property §4.1 actually cares about. AC9's *"unchanged hashes ⇒ fail"* clause is **struck**; *"hand-authored JSON ⇒ fail"* stands. Four modules / six tar entries unchanged: both re-asserted. | §1.2(i), measured V11 |

### Addition

| AC | Statement |
|---|---|
| **AC13 (new)** | `scripts/verify_ail.sh`'s terminal banner reads `✓ verify gate PASSED: 10 required identities verified, 39 named tests pass`. This pin is **NOT gate-enforced** (§1.2 ii); it is asserted by grep with a same-instrument control. No mutation arm is invented for it — a decorative arm would be worse than an honest note. |

### Deferrals that are RATIFIED and MUST NOT be "fixed"

The surviving round-2 quorum objection — *no acceptance criterion requires any of the five proven
laws to be INVOKED; they are provable and inert until a host path exists* — is **correct and is now
a declared, ratified deferral.** The emitter is item 7's; `context.Context` plumbing is item 18's.

**The executor MUST NOT** add a milestone, an AC, a host emitter, a call site, or a Go struct to
make the laws reachable. **AC11 fails the sprint** if the diff touches `host/store`, `host/daemon`,
`host/broker`, `schema.sql`, or adds any host emitter. Equally out of scope: flipping
HUMAN-SURFACE §7.1/§7.3 to CLOSED, moving the §8 premise rows, or reconciling the item-15/item-7
charter rows (doc §10) — those are **controller/human acts**, named as non-ACs by doc §7's own tail.

### Hold set — re-measured at every milestone exit and again at M3

| Invariant | Value at base (`b9a7838`) | After |
|---|---|---|
| `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | rc=0 · `✓ 5/5 … across 11 module(s)` · `✓ all 20 required named tests pass` · `9/9 steps` | rc=0 · `✓ 10/10 … across 11 module(s)` · `✓ all 39 …` · `9/9` |
| `LEG1_MODULES` | 11 entries, set-compared | **MUST NOT MOVE** |
| isolated-root `ai-check` line count | 11 | **MUST NOT MOVE** |
| package exports | 4 | **MUST NOT MOVE** |
| tar entries (step 8/9) | 6 | **MUST NOT MOVE** |
| `interfaceHash` | `sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083` | **MUST NOT MOVE** (§1.2 i) |
| `world/types.ail` lines 1–132 | sha of prefix identical to base file `2cf5b004…` | **byte-untouched** — AC2's "existing exported types are byte-untouched" |
| `GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang go build ./... && go test ./...` | rc=0, 17 packages ok | rc=0 |
| bare `go test ./...` without `AILANG_BIN` | rc=1, ~8 failures | **BASE CONDITION, not a regression** — never report it as one |

---

## 3. Milestones — 4 milestones, 1 commit

### Verdict up front: **NO SPLIT of M1.** Measured, not argued.

The five files of M1 are one atomic unit *by gate construction*. Every candidate seam reds the tree:

| Proposed seam | What reds | Evidence |
|---|---|---|
| `types.ail` first, script pins second | Leg 1: `✗ expected exactly 5 proven world/ contracts, got 10`; also Leg 3 step 3/9 | `verify_ail.sh:311-313`, arithmetic from V6 |
| pins first, `types.ail` second | Leg 1: `required identity (types, timeoutOutcome) MISSING from verify.results[]` | `verify_ail.sh:290` |
| code + pins, projection second | Leg 3 step 3/9 sha256 equality | `verify_world_package.sh` step 3/9 |
| projection rebuilt, golden second | Leg 3 step 9/9 `cmp -s` | that is mutation **MU12**, not a milestone |
| AILANG first, Go marker second | `go test ./host/verifygate/` — CI's go-verify job — **FAILS** | **MEASURED, V13b** |

M2–M4 leave the committed tree byte-identical and are therefore green at every boundary by
construction; each is independently verifiable (M2 by sha equality against M1's commit, M3 by the
gates, M4 by the issue/message URLs).

---

### M1 — the atomic five-file landing · **COMMIT 1**

**Closes: AC2, AC3, AC4, AC8, AC9 (amended), AC13.**
**Advances (does not close): AC1** (declaration lands; the MU6 totality guard closes it in M2);
**AC10, AC11** (both re-measured and closed in M3).

Order below is causal, not optional.

#### M1.0 — pre-flight (rule 3e; do this before touching anything)

```bash
export AILANG_BIN=/tmp/ailang-v0300/ailang
export GOTOOLCHAIN=go1.25.6
mkdir -p /tmp/w15_backup
cp world/types.ail scripts/verify_ail.sh \
   host/verifygate/module_manifest_gate_test.go \
   scripts/world_package_ready_packet.golden.json /tmp/w15_backup/
cp scripts/world_package_ready_packet.golden.json /tmp/w15_backup/golden.base.json   # AC9 needs the OLD bytes
shasum -a 256 /tmp/w15_backup/*                                                       # record every one
./scripts/verify_ail.sh; echo "BASELINE rc=$?"     # must be rc=0, 5/5, 20
./scripts/verify_go.sh; echo "BASELINE rc=$?"      # must be rc=0
```

Backups are `cp`. **Every restore in this sprint is `cp` from `/tmp/w15_backup/`, never
`git checkout -- <file>`** — the worktree files are uncommitted by construction during M2, and
`git checkout --` would delete the executor's work and then report the disaster rather than prevent
it.

#### M1.1 — `world/types.ail` (+147 lines; 132 → 279)

**Append at end of file**, leaving lines 1–132 byte-for-byte untouched. This exact text was applied,
`ai-check`ed, tested, mutated and gate-run by the planner (V4–V6, V12). Land it verbatim:

```ailang

-- The typed finite set of timeout policies (HUMAN-SURFACE §3.1 / §7 point 1).
-- RATIFIED (attended, Mark, 2026-08-14): exactly these three members. Nullary
-- by design: the escalation bound is packet state, not policy identity.
export type TimeoutPolicy
  = Cancel
  | EscalateBounded
  | ExecuteIfGranted

-- The typed outcome of the explicit Timeout transition at deadline. Silence
-- never synthesizes approval: no outcome constructor approves anything, and
-- ExecuteUnderPriorAuthority requires authority that already existed.
export type TimeoutOutcome
  = ResolveRejectedTimeout
  | EscalateWithNewDeadline
  | ExecuteUnderPriorAuthority

-- The frozen v1 decision packet (HUMAN-SURFACE §7 point 3), semantic ID
-- world/decision-packet/v1. References, not embeds: proposalHash reaches
-- evidence/caps; requestRef reaches the landed approval objects.
-- createdAt/deadlineAt are LOGICAL times (journal.go:26-27 idiom); no field
-- reads a wall clock. Amendments are /v2, never in-place.
export type DecisionPacket = {
  packetHash: HashRef,
  proposalHash: HashRef,
  requestRef: HashRef,
  createdAt: int,
  deadlineAt: int,
  escalationsRemaining: int,
  policy: TimeoutPolicy
}

-- The total timeout law: what the explicit Timeout transition must do, per
-- policy, given the packet's escalation budget and the host-derived fact of
-- whether independent authority is live at the recorded logical time.
export func timeoutOutcome(policy: TimeoutPolicy, escalationsRemaining: int, independentAuthority: bool) -> TimeoutOutcome ! {}
ensures { result == match policy {
  Cancel => ResolveRejectedTimeout,
  EscalateBounded => if escalationsRemaining > 0 then EscalateWithNewDeadline else ResolveRejectedTimeout,
  ExecuteIfGranted => if independentAuthority then ExecuteUnderPriorAuthority else ResolveRejectedTimeout
} }
{
  match policy {
    Cancel => ResolveRejectedTimeout,
    EscalateBounded => if escalationsRemaining > 0 then EscalateWithNewDeadline else ResolveRejectedTimeout,
    ExecuteIfGranted => if independentAuthority then ExecuteUnderPriorAuthority else ResolveRejectedTimeout
  }
}

-- DEFER rebound law: a defer is valid only with a strictly-future new deadline
-- and remaining escalation budget. DEFER MUST NOT park indefinitely.
export func validDefer(now: int, newDeadline: int, escalationsRemaining: int) -> bool ! {}
ensures { result == (newDeadline > now && escalationsRemaining > 0) }
{
  newDeadline > now && escalationsRemaining > 0
}

-- Packet schedule well-formedness: creation time is a valid logical instant,
-- the deadline is strictly after creation, and the budget is non-negative.
export func wellFormedSchedule(createdAt: int, deadlineAt: int, escalationsRemaining: int) -> bool ! {}
ensures { result == (createdAt >= 0 && deadlineAt > createdAt && escalationsRemaining >= 0) }
{
  createdAt >= 0 && deadlineAt > createdAt && escalationsRemaining >= 0
}

-- Timeout deadline law: an explicit Timeout transition is legal only at or
-- after the packet's deadline, judged against the RECORDED logical now.
-- This is what connects the deadline to the outcome: an early timeout is a
-- detectable lie under replay, not a compatible history.
export func timeoutFiredLegally(deadlineAt: int, recordedNow: int) -> bool ! {}
ensures { result == (recordedNow >= deadlineAt) }
{
  recordedNow >= deadlineAt
}

-- Escalation/DEFER rebind law: any recorded rebind must decrement the budget
-- by EXACTLY one and set a strictly-future new deadline. validDefer guards
-- "may a defer happen now"; this law validates the recorded before/after
-- pair, so a non-decrementing escalation is a detectable lie under replay.
export func validEscalation(oldEscalationsRemaining: int, newEscalationsRemaining: int, recordedNow: int, newDeadlineAt: int) -> bool ! {}
ensures { result == (oldEscalationsRemaining > 0 && newEscalationsRemaining == oldEscalationsRemaining - 1 && newDeadlineAt > recordedNow) }
{
  oldEscalationsRemaining > 0 && newEscalationsRemaining == oldEscalationsRemaining - 1 && newDeadlineAt > recordedNow
}

-- Private test adapters. Single-int case dispatchers because v0.30.0's inline
-- harness cannot take multiple arguments (parse fail) nor tuple-valued inputs
-- (collected, then "no pattern matched" at runtime) — both measured, both
-- routed upstream. Integer codes are test plumbing, not an exported ordering.
func outcomeCode(caseId: int) -> int
tests [
  (1, 1), (2, 2), (3, 1), (4, 3), (5, 1), (6, 1)
]
{
  let code = \o. match o {
    ResolveRejectedTimeout => 1,
    EscalateWithNewDeadline => 2,
    ExecuteUnderPriorAuthority => 3
  };
  if caseId == 1 then code(timeoutOutcome(Cancel, 5, true))
  else if caseId == 2 then code(timeoutOutcome(EscalateBounded, 2, false))
  else if caseId == 3 then code(timeoutOutcome(EscalateBounded, 0, true))
  else if caseId == 4 then code(timeoutOutcome(ExecuteIfGranted, 3, true))
  else if caseId == 5 then code(timeoutOutcome(ExecuteIfGranted, 3, false))
  else code(timeoutOutcome(Cancel, 0, false))
}

func deferCode(caseId: int) -> int
tests [
  (1, 1), (2, 0), (3, 0)
]
{
  if caseId == 1 then (if validDefer(10, 20, 2) then 1 else 0)
  else if caseId == 2 then (if validDefer(10, 10, 2) then 1 else 0)
  else (if validDefer(10, 20, 0) then 1 else 0)
}

func scheduleCode(caseId: int) -> int
tests [
  (1, 1), (2, 0), (3, 0)
]
{
  if caseId == 1 then (if wellFormedSchedule(0, 1, 0) then 1 else 0)
  else if caseId == 2 then (if wellFormedSchedule(5, 5, 0) then 1 else 0)
  else (if wellFormedSchedule(0-1, 1, 0) then 1 else 0)
}

func firedCode(caseId: int) -> int
tests [
  (1, 0), (2, 1), (3, 1)
]
{
  if caseId == 1 then (if timeoutFiredLegally(10, 9) then 1 else 0)
  else if caseId == 2 then (if timeoutFiredLegally(10, 10) then 1 else 0)
  else (if timeoutFiredLegally(10, 11) then 1 else 0)
}

func escalationCode(caseId: int) -> int
tests [
  (1, 1), (2, 0), (3, 0), (4, 0)
]
{
  if caseId == 1 then (if validEscalation(3, 2, 10, 20) then 1 else 0)
  else if caseId == 2 then (if validEscalation(0, 0-1, 10, 20) then 1 else 0)
  else if caseId == 3 then (if validEscalation(3, 3, 10, 20) then 1 else 0)
  else (if validEscalation(3, 2, 10, 10) then 1 else 0)
}
```

> Note the leading blank line: the block is appended after the existing final line
> `  | Denied(string)`, separated by exactly one blank line.

**LANDED-PROOF (mandatory).** `shasum -a 256 world/types.ail` must read
`2cf5b004f7f0573f…` **before** and **`f8e08f68734724e4…`** after, and
`head -132 world/types.ail | shasum -a 256` must equal the sha of the pre-edit file.
**If the post sha is anything else, STOP.** The golden pre-registered in M1.5 was computed from
`f8e08f68…` and will not match; and a green print from a file that was never written is
indistinguishable from a green print from a file that was.

#### M1.2 — `scripts/verify_ail.sh` — **six** edits, all in this commit

Line numbers are the *current* ones (base `b9a7838`) and **shift as you edit** (the file goes 378 →
386 lines). Anchor on the exact strings; assert an occurrence count of **1** before each replace.

| # | Line (base) | From | To | Mechanism |
|---|---|---|---|---|
| **1** | `135` | `LEG1_MODULES=(` … 11 entries | **UNCHANGED** | bash array, **set-compared** at `:232-236` with an empty-allowlist floor at `:227`. No module is added. **Do not "fix" it.** |
| **2** | `266` | `    "world/types.ail":       {"gradeOf"},` | `    "world/types.ail":       {"gradeOf", "timeoutOutcome", "timeoutFiredLegally",`<br>`                             "validEscalation", "validDefer", "wellFormedSchedule"},` | **python dict** inside a heredoc; Leg-1 **named** proof pin (`:287-294`) |
| **3** | `310` | `EXACT_TOTAL_VERIFIED=5` | `EXACT_TOTAL_VERIFIED=10` | **shell** var, **no spaces** around `=`; interpolated into `:315`'s marker |
| **4** | `333-341` | `REQUIRED_TESTS = { … }` | append the **19** names below | **python set**; Leg-2 **named** test pin (`:355-361`). Also update the header comment `# logepoch (8) + contracts (6) + types (6)` → `types (25)` |
| **5** | `342` | `EXACT_TOTAL_TESTS = 20` | `EXACT_TOTAL_TESTS = 39` | **python** var, **SPACES around `=`** — a shell-shaped `grep 'EXACT_TOTAL_TESTS='` MISSES this line |
| **6** | `378` | `echo "✓ verify gate PASSED: 5 required identities verified, 20 named tests pass"` | `echo "✓ verify gate PASSED: 10 required identities verified, 39 named tests pass"` | **literal, not interpolated; UNGUARDED** (§1.2 ii → **AC13**). Omitted from doc §8.1. |

The 19 names for edit 4, exactly:

```python
    "outcomeCode_test_1", "outcomeCode_test_2", "outcomeCode_test_3",
    "outcomeCode_test_4", "outcomeCode_test_5", "outcomeCode_test_6",
    "deferCode_test_1", "deferCode_test_2", "deferCode_test_3",
    "scheduleCode_test_1", "scheduleCode_test_2", "scheduleCode_test_3",
    "firedCode_test_1", "firedCode_test_2", "firedCode_test_3",
    "escalationCode_test_1", "escalationCode_test_2", "escalationCode_test_3",
    "escalationCode_test_4",
```

**Edits 4 and 5 are separate obligations.** Moving the total to 39 while leaving `REQUIRED_TESTS`
untouched still passes (`n==39` holds, all 20 old names present) while leaving all 19 new
identities unpinned **by name** — a weakening that prints green. **MU10** is the arm that proves it.
Symmetrically, edit 2 and edit 3 are separate: **MU10** covers the proof-name pin.

#### M1.3 — `packages/world-core/world/types.ail`

```bash
./scripts/build_world_package.sh                                          # wholesale re-projection
cmp -s world/types.ail packages/world-core/world/types.ail && echo "projection EQUAL"   # want 0
cmp -s world/types.ail world/contracts.ail; echo "control (must be 1): $?"              # instrument fires
shasum -a 256 packages/world-core/world/types.ail   # must be f8e08f68734724e4…
```

`build_world_package.sh` replaces `$DEST_DIR` wholesale from a fresh `mktemp -d` staging dir, so a
stale file cannot survive. Expect `allowlisted modules: iterated=4 wc-l=4` and
`projected 4 modules into packages/world-core/world`.

#### M1.4 — `host/verifygate/module_manifest_gate_test.go:128`

```diff
-	const marker = "✓ 5/5 required world/ identities verified across 11 module(s)"
+	const marker = "✓ 10/10 required world/ identities verified across 11 module(s)"
```

`11 module(s)` **stays 11.** This is a substitution of one digit-pair, not a widening. Five tests
call `requirePristineControl` (`:172, 215, …`). Introduce no rig-absolute path literal
(`TestNoRigAbsolutePaths` globs `host/verifygate/*.go`); none is needed.

#### M1.5 — `scripts/world_package_ready_packet.golden.json`

**Do NOT hand-edit.** Regenerate through `host/pkgproj`, and run the **known-positive control
first** (the planner did — it reproduced the committed golden byte-for-byte, V9, and did so from a
copied directory, V10, proving path-independence):

```go
// /tmp/w15_regen_golden.go — throwaway, NEVER committed
package main

import (
	"os"

	"github.com/sunholo-data/ailang-world/host/pkgproj"
)

func main() {
	// Version is REQUIRED here and is ABSENT from verify_world_package.sh:175 —
	// that helper only prints hashes, so an empty Version is invisible there.
	m := pkgproj.Manifest{
		Package: pkgproj.Package{Name: "world/core", Version: "0.1.0", Edition: "1", AILANG: ">=0.30.0"},
		Exports: pkgproj.Exports{Modules: []string{"world/types", "world/contracts", "world/transitions", "world/logepoch"}},
		Effects: pkgproj.Effects{Max: []string{}},
	}
	p, err := pkgproj.RecomputeReadyPacket(os.Args[1], m, "AILANG v0.30.0")
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(pkgproj.EncodeReadyPacket(p))
}
```

```bash
# CONTROL FIRST — run BEFORE M1.1, on the untouched tree. Must reproduce the COMMITTED golden.
GOTOOLCHAIN=go1.25.6 go run /tmp/w15_regen_golden.go packages/world-core > /tmp/w15_ctl.json
cmp -s /tmp/w15_ctl.json scripts/world_package_ready_packet.golden.json && echo "CONTROL PASS"

# then, AFTER M1.3:
GOTOOLCHAIN=go1.25.6 go run /tmp/w15_regen_golden.go packages/world-core \
  > scripts/world_package_ready_packet.golden.json
```

**Pre-registered expected bytes** (planner-measured, V11). A mismatch means `world/types.ail`'s
bytes differ from M1.1 and is a **STOP**, not a golden to accept:

```json
{"compilerVersion":"AILANG v0.30.0","contentHash":"sha256:06acbb83ce8824166a2852912776edfe749ba741052f3daa8dd16c99abfb526e","effects":[],"exports":["world/types","world/contracts","world/transitions","world/logepoch"],"interfaceHash":"sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083","package":"world/core","tarballBytes":7856,"tarballSHA256":"sha256:5823edcfbb3fa640080f405771578d84c71de72184b7cfda9033a50f76de608d","version":"0.1.0"}
```

`interfaceHash` is **identical** to the committed golden's (§1.2 i / AC9 amended). Two independent
implementations then cross-check these bytes: `verify_world_package.sh:230-244`'s python emitter
(`cmp -s`) and `host/pkgproj/readypacket_test.go`.

*(Fallback if the helper is unavailable: run `./scripts/verify_world_package.sh`, let step 9/9 fail,
and take the `+{…}` line from its own `diff -u` output — that line is `$tmp_ready`, produced by
`pkgproj`, so it is still derived and not hand-authored. Re-run to green.)*

#### M1.6 — commit

**COMMIT 1** = exactly these five files and nothing else:

```
world/types.ail
packages/world-core/world/types.ail
scripts/verify_ail.sh
host/verifygate/module_manifest_gate_test.go
scripts/world_package_ready_packet.golden.json
```

Verify before committing: `git status --porcelain` shows exactly those five, `M` only, no `??`.
A stray `.ail` under `world/` reds the `LEG1_MODULES` set compare; a stray file under
`packages/world-core` reds step 2/9.

---

### M2 — the mutation drill, 16 arms · **NO COMMIT** (tree restored byte-identical)

**Closes: AC1 (via MU6), AC5, AC6, AC7.**
**Re-proves non-vacuously: AC3 (MU10), AC4 (MU10), AC8 (MU13), AC9 (MU11, MU12).**

This is the sprint's spine and the bulk of its time. The milestone's deliverable is **evidence**,
not a diff: at exit, every touched file must be byte-identical to COMMIT 1 by sha256. That
satisfies "green at the boundary" trivially and is independently verifiable
(`git status --porcelain` → empty).

#### M2.1 — the protocol, every arm, no exceptions

1. **Backup exists** in `/tmp/w15_backup/` (M1.0). Record the pre sha.
2. **Apply** the mutation with a replacement that **asserts its exact occurrence count** (1 for
   body-only, 2 for "both"). A `sed` that silently matched nothing is the failure this step exists
   to catch.
3. **Assert it LANDED by sha256**: pre ≠ post. *A mutation that never applied and a mutation that
   failed to red share an exit code.*
4. **Assert the mutant still checks/builds** before reading any test result: `check.passed` must
   still be `true` for AILANG arms (`ai-check` JSON), `go vet ./host/verifygate/` must pass for Go
   arms. **"The mutant does not compile" is a third fact wearing the same exit code as "the mutant
   was killed."** For MU5 and MU6 this step is inverted deliberately: their *expected* observable
   IS `verify.errors == 1` with `check.passed` still `true` — record both fields.
5. **Read the SCOPED observable, never the exit code.** For test arms, parse the JSON and print
   **the list of non-`pass` identities** — the AC requires the named test to be the **SOLE**
   failure, which a suite-wide red/green cannot show. For Go arms use `-run '<TestName>'` and read
   the `t.Fatalf` message text.
6. **Restore by `cp` from `/tmp/w15_backup/`** — **never `git checkout -- <file>`.** Assert
   byte-identical by sha256 against the pre value.
7. **Re-run the positive control and require GREEN before the next arm.**

Reference implementation of steps 2–3 and 5 (adapt per arm):

```bash
pre=$(shasum -a 256 world/types.ail | cut -d' ' -f1)
python3 - <<'PY'
p='world/types.ail'; s=open(p,encoding='utf-8').read()
old='<EXACT OLD>'; new='<EXACT NEW>'
n=s.count(old); assert n==<EXPECTED_N>, "anchor count %d != <EXPECTED_N>" % n
open(p,'w',encoding='utf-8').write(s.replace(old,new,<EXPECTED_N>))
PY
post=$(shasum -a 256 world/types.ail | cut -d' ' -f1)
[ "$pre" != "$post" ] || { echo "MUTATION DID NOT LAND — STOP"; exit 1; }

$AILANG_BIN ai-check -timeout 15s world/types.ail > /tmp/w15_m.json      # rc is NOT the oracle
python3 -c "
import json;d=json.load(open('/tmp/w15_m.json'));v=d['verify']
print('check.passed=%s verified=%s cex=%s errors=%s'%(d['check']['passed'],v['verified'],v['counterexample'],v['errors']))
print([(r['function'],r['status']) for r in v['results']])"

$AILANG_BIN test --format json world/ > /tmp/w15_t.json 2>/dev/null || true
python3 -c "
import json;raw=open('/tmp/w15_t.json').read();d=json.loads(raw[raw.find('{'):])
print('n=%d failed=%s'%(len(d['tests']),d.get('failed_tests')))
print('NON-PASS:',[(x['name'],x['status']) for x in d['tests'] if x['status']!='pass'])"

cp /tmp/w15_backup/types.ail world/types.ail
[ "$(shasum -a 256 world/types.ail | cut -d' ' -f1)" = "$pre" ] || { echo "RESTORE FAILED — STOP"; exit 1; }
```

#### M2.2 — the 16 arms

Rows marked **PLANNER-MEASURED** were landed, read, restored and sha-verified by the planner **on
the exact bytes M1.1 lands** (§9). The rest carry the doc's V-P8/V-P9/V-P14 measurement on
semantically identical text; the executor runs all 16 live regardless.

| ID | Mutation (occurrence count) | AC | Required observable | Control |
|---|---|---|---|---|
| **MU1** | `Cancel => ResolveRejectedTimeout` → `EscalateWithNewDeadline`, **body only** (n=1) | 7 | `verify.counterexample >= 1`, `timeoutOutcome status=counterexample` | restore → `verified`, cex=0 |
| **MU2** | `ExecuteIfGranted` no-auth arm → `ExecuteUnderPriorAuthority`, **both** (n=2) — *silence synthesizes execution* | **5**, 7 | **Z3 stays GREEN** (`verified=6, cex=0, errors=0`) **and** `outcomeCode_test_5` is the **SOLE** failure | restore → 39/39 |
| **MU3** | `EscalateBounded` exhaustion arm → `EscalateWithNewDeadline`, **both** (n=2) — *unbounded park* | 7 | `outcomeCode_test_3` fails, alone | `outcomeCode_test_2` stays `pass` |
| **MU4** | `EscalateBounded` live arm → `ResolveRejectedTimeout`, **both** (n=2) | 7 | `outcomeCode_test_2` fails, alone | `outcomeCode_test_3` stays `pass` |
| **MU5** | delete the `EscalateBounded` arm from the `timeoutOutcome` **body**, contract kept (n=1) | **6** | `verify.errors == 1`, `timeoutOutcome status=error`, **`check.passed` still `true`, `ai-check` rc still 0** | restore → `verified`, errors=0 |
| **MU6** | add a 4th `TimeoutPolicy` constructor (e.g. `\| ParkForever`) with **no** arm | **1**, 6 | `verify.errors == 1` — the contract is the only totality guard | add an explicit arm → `verified` |
| **MU7** | `validDefer`: drop `newDeadline > now`, **both** (n=2) | 7 | `deferCode_test_2` fails, alone | `deferCode_test_3` stays `pass` |
| **MU8** | `validDefer`: drop `escalationsRemaining > 0`, **both** (n=2) | 7 | `deferCode_test_3` fails, alone | `deferCode_test_2` stays `pass` |
| **MU9** | `wellFormedSchedule`: `deadlineAt > createdAt` → `>=`, **both** (n=2) | 7 | `scheduleCode_test_2` fails, alone | `scheduleCode_test_1/3` stay `pass` |
| **MU10** | remove `"timeoutOutcome"` from `REQUIRED_VERIFIED["world/types.ail"]` | 3, 4 | the named identity leaves the Leg-1 manifest while `EXACT_TOTAL_VERIFIED=10` still passes — **the count pin does not imply the name pin** | restore literal → gate green |
| **MU11** | edit canonical `world/types.ail`, **skip** `build_world_package.sh` | 9 | Leg 3 step 3/9: `✗ projection mismatch` / sha256 inequality | rebuild projection → step 3/9 green |
| **MU12** | rebuild the projection, keep the **old** golden | 9 | Leg 3 step 9/9: `✗ ready packet differs byte-for-byte from golden` + a `diff -u` naming `contentHash`/`tarballSHA256`/`tarballBytes` — **and NOT `interfaceHash`** | recompute via `pkgproj` → step 9/9 green |
| **MU13** | revert `module_manifest_gate_test.go:128` to `5/5` | **8** | `go test ./host/verifygate/ -run TestModuleManifest` **FAILS**: `pristine isolated control missing "✓ 5/5 required world/ identities verified across 11 module(s)"` | restore → `ok` |
| **MU14** | `timeoutFiredLegally`: accept an early fire (`recordedNow >= deadlineAt` → `>= deadlineAt - 1`), **both** (n=2) | **5**, 7 | **Z3 stays GREEN** (`verified=6, cex=0, errors=0`); `firedCode_test_1` is the **SOLE** failure | restore → 39/39 |
| **MU15** | `validEscalation`: `newEscalationsRemaining == oldEscalationsRemaining - 1` → `<=`, **both** (n=2) | **5**, 7 | **Z3 stays GREEN**; `escalationCode_test_3` is the **SOLE** failure | `escalationCode_test_1/2/4` stay `pass` |
| **MU16** | `validEscalation`: `newDeadlineAt > recordedNow` → `>=`, **both** (n=2) | **5**, 7 | **Z3 stays GREEN**; `escalationCode_test_4` is the **SOLE** failure | `escalationCode_test_1/2/3` stay `pass` |
| — | **Fabricated authority** | — | **NO mutation exists that a pure gate can red.** `independentAuthority` is a host-derived bool INPUT; no pure law can verify a bool's provenance. Carried as the **declared item-7 host residual** (doc §5 coverage map), **never as a pretended kill.** **Do not invent an arm for it.** | n/a |

**MU15's `<=` subtlety.** `x <= y - 1` is equivalent to `x < y` over ints, so this arm admits any
*over*-decrement while still refusing a non-decrement. If it does not kill `escalationCode_test_3`
(inputs `(3, 3, 10, 20)` — a non-decrement), use the doc's stated intent instead: relax the clause
to `newEscalationsRemaining <= oldEscalationsRemaining` and re-read. **Record which form landed.**
An arm that cannot fire is decoration; say so rather than reporting a kill.

Also carry the **non-vacuity control in the other direction** (doc V-P15), which no AC names but
which is cheap and proves the new contracts are not tautological: mutate **only the body**
occurrence of `recordedNow >= deadlineAt` to `>` and require
`verify.counterexample == 1`, `timeoutFiredLegally status=counterexample`.

---

### M3 — full pinned gates, hold-set re-measurement, scope inspection · **NO COMMIT**

**Closes: AC10, AC11, AC13.**

```bash
export AILANG_BIN=/tmp/ailang-v0300/ailang
export GOTOOLCHAIN=go1.25.6

# AC10 — the pinned AILANG gate.
#   BASE: rc=0, "✓ 5/5 … across 11 module(s)", "✓ all 20 …", "9/9 steps".
#   AFTER: rc=0, "✓ 10/10 … across 11 module(s)", "✓ all 39 …", "9/9 steps".
./scripts/verify_ail.sh; echo "rc=$?"

# AC10 — the pinned Go gate. GOTOOLCHAIN is a BASE CONDITION of this rig, not a regression.
./scripts/verify_go.sh; echo "rc=$?"
go build ./...; echo "build rc=$?"
go test ./... -count=1; echo "test rc=$?"     # 17 packages; host/broker carries a ~18% BASE flake

# AC11 — scope. Exactly five files; nothing under host/store, host/daemon, host/broker, schema.sql.
git diff --name-only HEAD~1..HEAD
git status --porcelain                                   # want EMPTY (M2 restored everything)
git diff --name-only HEAD~1..HEAD | grep -E 'host/store|host/daemon|host/broker|schema\.sql'; echo "forbidden-hit rc=$? (want 1 = none)"
git diff --name-only HEAD~1..HEAD | grep -c 'world/types.ail'; echo "control (want >=1)"

# AC13 — the unguarded terminal banner, with a same-instrument control.
grep -c '✓ verify gate PASSED: 10 required identities verified, 39 named tests pass' scripts/verify_ail.sh   # want 1
grep -c '5 required identities verified\|20 named tests pass' scripts/verify_ail.sh                          # want 0
grep -c 'verify gate PASSED' scripts/verify_ail.sh                                                           # CONTROL, want 1
```

**Hold-set re-measurement** (§2's table) is run in full at M3 exit. Report each value beside its
base value; a moved must-not-move is a STOP.

**`host/broker` flake protocol**: if `go test ./...` reds only in `host/broker`, re-run
`go test ./host/broker/ -count=5` and compare against the recorded ~18% base rate. **Do not "fix"
`host/broker`** — AC11 forbids touching it, and a green single run is not evidence the flake is
gone either.

---

### M4 — upstream routing · **NO COMMIT** (no repo file changes)

**Closes: AC12.**

The mission's no-local-workarounds rule: language gaps route to `sunholo-data/ailang` as issues
plus an `ailang messages send mission-control` note. **Three** issues, one per §5 limitation, each
citing the design doc's evidence row:

| # | Title | Evidence | Body must contain |
|---|---|---|---|
| 1 | *v0.30.0: a record containing a bare ADT field is Z3-unencodable and fails silently at rc level* | **V-P2** (control **V-P3**) | the repro (record with an ADT-typed field, contract on the record param); observed `check.passed=true`, **rc 0**, `verify.errors=1`, `Z3 error … unknown sort '<T>'`; the flat-int-record control that verifies. Note this is **stricter** than the repo's previously recorded `list[ADT]` limitation — the list is not needed. |
| 2 | *v0.30.0: inline-test harness cannot execute tuple-valued inputs (collected, then "no pattern matched")* | **V-P5** | multi-argument `tests [...]` rows fail to **parse**; tuple-valued inputs are **collected** then fail at runtime with `no pattern matched in match expression`, in **both** nested-match and direct tuple-pattern forms. |
| 3 | *v0.30.0: ADT equality in a general boolean expression needs an `Eq` instance; the suggested `import std/prelude` workaround is unsupported* | **V-P12** | `No instance for Eq[<T>] in scope. Equality (==, !=) needs an Eq instance`; the suggested fix `import std/prelude` fails with `IMP012_UNSUPPORTED_NAMESPACE … namespace imports not yet supported`. Note the top-level `result == match …` postcondition form DOES verify — this is what forces a three-law split instead of one combined law. |

```bash
gh issue create --repo sunholo-data/ailang --title "…" --body "…"     # ×3, record the URLs
ailang messages send mission-control "<one note naming all three issue URLs and item 15>"
```

**AC12 fails if any of the three is silently absorbed as local lore.** Record the three issue URLs
and the message ID in the sprint record. `gh` is at `/opt/homebrew/bin/gh`. Note that
`ailang messages` runs on the PATH binary — that is a **messaging** tool, not a gate, so the
`-dirty` PATH build is acceptable there and **only** there.

---

## 4. Day-by-day

One working day, four blocks. Times are the planner's re-price (§7), not the doc's.

| Block | Elapsed | Work | Exit gate |
|---|---:|---|---|
| **1** | 0.00 → 0.30 d | M1.0 pre-flight (baselines, backups, golden control) + M1.1–M1.6 | `verify_ail.sh` rc=0 with `10/10` / `39` / `9/9`; `go test ./host/verifygate/` ok; COMMIT 1 landed with exactly 5 files |
| **2** | 0.30 → 0.62 d | M2, arms MU1–MU9 and MU14–MU16 (the AILANG arms; 12 arms) | every arm landed by sha, killed its named identity **alone**, restored byte-identical; `git status --porcelain` empty |
| **3** | 0.62 → 0.82 d | M2, arms MU10–MU13 (the gate-pin arms) + the body-only non-vacuity control + M3 gates and scope | all four arms red their named gate step; full pinned gates green; hold set unmoved |
| **4** | 0.82 → 0.90 d | M4 upstream routing; sprint record | 3 issue URLs + 1 message ID recorded |

Milestone totals: **M1 0.30 d** (0.10 addition + 0.05 pins + 0.08 projection/golden + 0.05 AC9
reconciliation + 0.02 pin 6) · **M2 0.42 d** · **M3 0.10 d** · **M4 0.08 d** = **0.90 d**, matching
§7.2 line for line.

---

## 5. Test plan

There is no new Go test and no new test file — **by design and by AC11.** The sprint's test surface
is entirely the pinned gates plus the mutation drill.

| Layer | Instrument | Oracle | Non-vacuity proof |
|---|---|---|---|
| Z3 contracts | `ailang ai-check --` JSON | `verify.results[].status` per named identity; `verify.errors`; `verify.counterexample`. **Never rc.** | MU1 (counterexample), MU5/MU6 (`errors==1`), plus the body-only V-P15 control |
| Named runtime identities | `ailang test --format json world/` | `len(tests[])` and per-identity `status`. **Never `passed_tests`** (43, not 39). | MU2/MU3/MU4/MU7/MU8/MU9/MU14/MU15/MU16 — each kills exactly one named identity while Z3 stays green |
| Gate name-pins | `verify_ail.sh` Legs 1–2 | the gate's own parsed refusal line | MU10 — count pins pass, name pin reds |
| Package projection | `verify_ail.sh` Leg 3, steps 3/9 and 9/9 | `sha256` equality; `cmp -s` against the golden | MU11, MU12 |
| Go marker | `go test ./host/verifygate/` | `requirePristineControl`'s `t.Fatalf` text | MU13 (**PLANNER-MEASURED**, V13b) |
| Terminal banner | `grep` with a same-instrument control | literal match count | **none exists** — declared unguarded (AC13, §1.2 ii) |

**The proof/test division, stated so nobody overclaims it.** Z3 proves each law is exactly its
stated conjunction, and proves `timeoutOutcome` **total** over the three-constructor set. Z3 is
**blind to a consistent lie**: MU2 and MU14 were measured this session to leave `verified=6, cex=0,
errors=0` while encoding, respectively, "silence synthesizes execution" and "an early timeout is
legal" — the two cardinal §3.1 violations. **The named runtime identities are the only thing that
kills them.** Any future weakening of `outcomeCode`/`firedCode`/`escalationCode` reopens a §3.1
violation, which is precisely why AC5 exists.

---

## 6. Risks

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Executor pursues AC9's *"unchanged hashes ⇒ fail"* and stalls, or hand-edits the golden | **HIGH without §1.2(i)** | sprint failure or a corrupt golden | AC9 **amended** in §2; `interfaceHash` invariance is in the hold set as a **must-not-move**; the expected bytes are pre-registered |
| Pin 6 (`verify_ail.sh:378`) omitted — the gate's own summary becomes a lie | **HIGH** (two items running) | silent, permanent misinformation | listed as edit 6 in M1.2; **AC13** with a control; `grep -c '5 required identities verified'` must read 0 |
| `grep 'EXACT_TOTAL_TESTS='` misses line 342 (spaces around `=`) | medium | pin not moved; gate reds for a confusing reason | M1.2 row 5 names the mechanism; AC4's `n==39` catches the omission |
| Executor reads `passed_tests=43` (or step 5/9's `42 activity lines`) and "fixes" the total | medium | red gate | §1.4(1); the oracle is `len(tests[])` at `verify_ail.sh:367`; step 5/9 has **no count pin** |
| Executor "fixes" `LEG1_MODULES` because the pin map mentions it | medium | red set-compare | M1.2 row 1 says **UNCHANGED** in bold; the 11-module / 11-`ai-check`-line census is in the hold set |
| Golden regenerated with an empty `"version"` | medium | step 9/9 reds with a hash-shaped diff for a manifest-shaped reason | §1.4(3); the reproduce-the-committed-golden **control runs first** |
| A mutation "passes" because it never landed | medium | a fake kill recorded as evidence | M2.1 step 3 — sha256 landed-proof on **every** arm, occurrence-count assert on every replace |
| A mutant that does not compile is scored as a kill | medium | vacuous arm | M2.1 step 4 — assert `check.passed` / `go vet` **before** reading any test result |
| A suite-wide red hides *which* identity died | medium | AC5's "SOLE failure" clause silently unverified | M2.1 step 5 — print the **list** of non-`pass` identities, never rc |
| `git checkout --` used to restore during M2 | low-medium | **destroys the executor's uncommitted work** | M1.0 and M2.1 both state `cp` from `/tmp/w15_backup/`, twice, in bold |
| Scope creep: an executor "completes" the item by wiring a host emitter | low-medium | violates a **ratified deferral**; AC11 red | §2's deferral block; AC11 names the five forbidden paths; the emitter is item 7's |
| Scope creep: executor flips HUMAN-SURFACE §7.1/§7.3 to CLOSED | low-medium | a human gate performed by an agent | doc §7 tail names these as controller/human acts; §2 restates it; AC11's five-file list excludes `design_docs/` |
| `host/broker` ~18% base flake reds `verify_go.sh` | low-medium | false regression report | known base condition; re-run `-count=5` and compare; **never edit `host/broker`** |
| MU15's `<=` form cannot fire | low | a decorative arm reported as a kill | M2.2 note gives the alternative form and requires recording which landed |
| §7.1 was ratified only *by inclusion*, not by a lettered A/B | low | the wrong constructor set is frozen | §7.4 below; the risk is bounded because Option A is not executable without the set, and MU6 makes any later member a loud red rather than a silent fall-through |

---

## 7. Estimate, velocity, and the pricing verdict

### 7.1 Velocity, calculated two ways

The only close precedent is **item 13** (`w-evidence-grade-mapping`), which is nearly a scale model
of this item: same file, same shape (ADT + Z3 contract + private int-code adapter), the same five
gate pins, the same projection and golden. It landed at `36f0c7a` with **127 insertions across 8
files**; the doc priced **0.65 d**, its planner re-priced **≈0.75 d**, and it landed inside one
iteration.

- **By LOC.** 127 insertions / 0.75 d ≈ **170 insertions/day**. Item 15 lands ≈ 147 `.ail` lines
  ×2 (canonical + projection) + ≈ 25 script lines + 2 = **≈ 320 insertions** → **≈ 1.9 d**.
- **By mutation arm.** Item 13 ran **11 arms in a 0.30 d band** → **≈ 0.027 d/arm**. Item 15 has
  **16 arms** → **≈ 0.44 d** for the drill.

**LOC velocity is the wrong instrument here and I am saying so rather than quoting it.** These
sprints are not authoring-bound: the `.ail` payload is written once and pasted, while the cost is
`ai-check` + `ailang test` + full-gate wall time multiplied by the arm count. The arm-count metric
reproduces the doc's own 0.40 d mutation band almost exactly, which is evidence it is the right
unit. I use arm-count velocity.

### 7.2 The re-price

| Work | Doc §9 | Plan | Delta driver |
|---|---:|---:|---|
| `world/types.ail` addition (3 types + 5 contracts + 5 adapters), pinned-binary green | 0.32 d | **0.10 d** | text composed, verified, tested and gate-run by the planner; the executor pastes and asserts a sha |
| Six gate-pin edits across two files | *(in 0.40)* | **0.05 d** | all six validated end-to-end in an isolated repo copy (V12) |
| MU1–MU16 with landed-proofs, scoped reads, controls and restores | *(in 0.40)* | **0.42 d** | 16 arms × 0.027 d; 4 arms pre-measured on the exact landing bytes |
| Package projection + golden | *(in 0.28)* | **0.08 d** | recipe prototyped with **two** known-positive controls; expected bytes pre-registered |
| Full pinned gates + hold set + scope inspection | *(in 0.28)* | **0.10 d** | unchanged in substance |
| **AC9 amendment** — reconciling AC9 against `InterfaceHash` | — | **0.05 d** | §1.2(i). *Undiagnosed this is an open-ended stall, not 0.05 d* |
| **Pin 6** — `verify_ail.sh:378`, absent from the doc's conflict surface | — | **0.02 d** | §1.2(ii) |
| **Three upstream issues + `mission-control` note** | *(in 0.28)* | **0.08 d** | 3 GitHub issues with repro bodies and evidence citations, plus one message |
| **Total** | **1.00 d** | **≈ 0.90 d** | |

### 7.3 Verdict: **I agree with the 1.00 d headline; I disagree with the band split.**

The doc's total is sound and stays inside the queued ~1 d band. Two band-level objections, stated
rather than silently re-priced:

1. **Band 3 is underpriced.** §9 charges *"package projection + golden, full pinned gates, three
   upstream filings"* to **0.28 d**. Item 13 — the doc's own cited precedent — spent **0.30 d** on
   projection + golden + gates with **zero** upstream filings. Three GitHub issues with repro
   bodies plus a `mission-control` note is real work folded into a band that was already 0.02 d
   short before they were added. My split gives that work **0.26 d** (0.08 + 0.10 + 0.08).
2. **Band 1 is overpriced *given this plan*.** 0.32 d for the `.ail` addition was right when the
   text lived only in a `/tmp` probe. With the text pre-registered and gate-validated it is
   **0.10 d**. This is a saving the *plan* creates, not a defect in the doc.

The two roughly cancel, which is why **≈0.90 d ≈ 1.00 d** and why I am **not** asking the
controller to re-price the queue row. **The honest counterfactual matters more than the number:**
without §1.2's three corrections and the pre-registered artifacts, this sprint prices at
**≈1.15 d** — AC9 alone is an unbounded stall, not a 0.05 d line — and the doc's 1.00 d would have
been an *under*-price at the top of a band it already claimed to be at the edge of.

### 7.4 On the Ask-1-by-inclusion inference — I find it SAFE, with one condition

The controller's inference (that §7 point 1's three-constructor set is ratified by inclusion in
Option A) is sound and I plan on it. The reasoning:

- `DecisionPacket` has a `policy: TimeoutPolicy` field, so freezing the packet without the set is
  not a smaller decision — it is an **incoherent** one.
- All five laws are defined over `TimeoutPolicy`; `timeoutOutcome`'s totality proof *is* a proof
  over exactly that constructor set.
- Doc §2.1 argues the set is **closed by §3.1's own text**: every candidate fourth member examined
  ("park indefinitely", "auto-approve at deadline") is explicitly forbidden by the binding
  principle. Ratifying Option A while rejecting the set would require Mark to have intended a
  *fourth* member that §3.1 forbids.

**The one condition, and it is already satisfied by the plan:** because the set was not separately
lettered, the sprint must make a later amendment **loud rather than silent**. It does — **MU6**
proves that adding a `TimeoutPolicy` constructor without an arm is a `verify.errors == 1` red, not
a silent fall-through (v0.30.0 *accepts* non-exhaustive ADT matches at check time with rc 0, so the
contract is the only totality guard). If Mark's intent was a different set, the cost is a `/v2`
version bump — the exact cost Option A's reasoning (3) already accepts — and the gate will refuse
to let it land quietly. I therefore see no reason to work around the inference.

---

## 8. Rollback

| Situation | Action |
|---|---|
| A mutation arm leaves the tree dirty | `cp /tmp/w15_backup/<file> <path>`; verify by sha256 against the M1 commit. **Never `git checkout --`** while M2 is in flight. |
| M1's post-edit sha ≠ `f8e08f68…` | **STOP.** Do not regenerate the golden from the wrong bytes. `cp` the backup back, re-apply M1.1 verbatim, re-sha. Report the diff. |
| Golden mismatch after M1.5 | **STOP.** It means `world/types.ail` differs from M1.1. Do not accept the observed golden — the pre-registered bytes and the observed bytes disagreeing is a *finding*, not a number to overwrite. |
| Gate red after COMMIT 1 | `git reset --soft HEAD~1` is the controller's call, not the executor's. The executor reports the exact refusal line and stops. |
| CI red on the PR head | almost certainly pin 5 (`10/10` marker, go-verify job) or pin 6. Check `grep -c '5/5' host/verifygate/` → want 0, and AC13's greps. |
| A hold-set invariant moved | **STOP and report.** Do not "fix" the invariant to match the observation. |

The commit is a single, atomic, five-file change with no schema, no host code and no new file, so
rollback is `git revert` of one commit. Nothing in M2–M4 changes the tree.

---

## 9. Verification Log

All commands run at `b9a7838`, 2026-08-14, instrument `/tmp/ailang-v0300/ailang`
(**AILANG v0.30.0**, `e37b370`). Semantic verdicts come from JSON or from a gate's parsed refusal
line, **never from process rc**. Probes ran under `/tmp/plan-item15/`; **no file in the worktree was
modified by the planner** except the two sprint artifacts.

| ID | Claim | Command | Observed |
|---|---|---|---|
| V0 | Worktree clean at `b9a7838`, correct branch | `git rev-parse --short HEAD; git branch --show-current; git status --porcelain` | `b9a7838`; `sprint/w-decision-lifecycle-freeze`; empty |
| V1 | Instrument pin | `/tmp/ailang-v0300/ailang --version` | `AILANG v0.30.0`, `Commit: e37b370` |
| V2 | All five doc pins resolve to the exact construct and line named | `grep -n` on `verify_ail.sh` and `module_manifest_gate_test.go` | `LEG1_MODULES=(` **:135** (11 entries); `REQUIRED_VERIFIED = {` **:262**, `"world/types.ail": {"gradeOf"}` **:266**; `EXACT_TOTAL_VERIFIED=5` **:310**; `REQUIRED_TESTS = {` **:333**; `EXACT_TOTAL_TESTS = 20` **:342**; marker `✓ 5/5 … across 11 module(s)` **:128**. `EXACT_TOTAL_MODULES` → **0 hits**, same-instrument control `EXACT_TOTAL_VERIFIED` → **4 hits** (fires) |
| V3 | Doc freshness (controller C5) | `git diff --name-only bc8f193..HEAD -- ':!design_docs'`; control without the pathspec | **0** files; control **6**, all under `design_docs/` |
| V4 | The exact M1.1 text applies to a copy of the real `world/` tree, leaving lines 1–132 byte-identical | `cp world/*.ail /tmp/plan-item15/world/`; append; `shasum`; `head -132 \| cmp` | `2cf5b004f7f0573f…` → **`f8e08f68734724e4…`**; 132 → **279** lines; prefix **IDENTICAL**; control (`cmp` vs `contracts.ail`) **differs** |
| V5 | Runtime identities and the true totals | `ailang test --format json world/` | **`len(tests[])=39`**, `failed_tests=0`, **`passed_tests=43`**; the 19 new identities exactly `outcomeCode_test_1..6`, `deferCode_test_1..3`, `scheduleCode_test_1..3`, `firedCode_test_1..3`, `escalationCode_test_1..4`; zero non-`pass` |
| V6 | `EXACT_TOTAL_VERIFIED` must become **10**, derived per module | `ai-check` each of the four `world/*.ail`, count `status=="verified"` | contracts **1**, logepoch **2**, transitions **1**, types **6** (`gradeOf`, `timeoutOutcome`, `timeoutFiredLegally`, `wellFormedSchedule`, `validEscalation`, `validDefer` — all `verified`, `errors=0`, `cex=0`) → **total 10** |
| V7 | The designer's V-P13 probe tree survives and reproduces its recorded sha | `shasum -a 256 /tmp/iso-item15-r2/types.r2.bak` | **`fda7f30bca720fb3…`** — matches doc V-P13 exactly. (M1.1's text differs from it only in comments; hence M1.1's own sha `f8e08f68…`.) |
| V8 | `InterfaceHash` cannot see a source-file change | read `host/pkgproj/pkgproj.go:86-104` | hashes **only** `Package.Name`, `Edition`, `AILANG`, sorted `Exports.Modules`, sorted `Effects.Max`. **Never opens a source file.** |
| V9 | Golden regeneration — **known-positive control first** | `go run /tmp/plan-item15/regen_golden.go packages/world-core \| cmp` vs the committed golden | **CONTROL PASS — byte-identical** |
| V10 | …and it is path-independent | same helper against an unmodified `cp -R` copy in `/tmp` | **CONTROL 2 PASS — byte-identical** |
| V11 | The new golden, and `interfaceHash` invariance | same helper against the copy carrying M1.1's `types.ail` (`f8e08f68…`) | `contentHash` `489d5e5d…`→**`06acbb83ce8824166a2852912776edfe749ba741052f3daa8dd16c99abfb526e`**; `tarballSHA256` `d0cdf42b…`→**`5823edcfbb3fa640080f405771578d84c71de72184b7cfda9033a50f76de608d`**; `tarballBytes` `6236`→**`7856`**; **`interfaceHash` UNCHANGED** `d16cc88270ff…`; `exports`, `effects`, `package`, `version`, `compilerVersion` all unchanged |
| V12 | **All six pin edits + projection + golden → the full pinned gate is GREEN** | isolated `tar`-copy of the worktree at `/tmp/plan-item15/repo`; apply all edits with occurrence-count asserts; `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | `✓ 10/10 required world/ identities verified across 11 module(s)`; `✓ all 39 required named tests pass (failed_tests=0)`; steps 1/9–9/9 all green (`4/4 projection hashes`, `4 exports`, `6 tar entries`, `canonical JSON equals committed golden byte-for-byte`); `✓ verify gate PASSED: 10 required identities verified, 39 named tests pass` |
| V13a | Pin 5 at `10/10` passes the Go gate | `go test ./host/verifygate/ -run TestModuleManifest -count=1` in that tree | **ok**, 18.399 s |
| V13b | **MU13** — a stale `5/5` marker REDS the go-verify job | revert the marker (sha `eb2239bf…`→`3616a8c8…`), re-run `-run TestModuleManifestRejectsStrayModule` | **FAIL**: `pristine isolated control missing "✓ 5/5 required world/ identities verified across 11 module(s)"`. Restored, sha back to `eb2239bf…`. *(The logged `rc=127` is an artifact of the isolated `/tmp` gate root's Leg 3; `requirePristineControl` does not assert rc, and the same test passes at `10/10`. Do not chase it.)* |
| V14 | Pin 6 is a literal and nothing pins its digits | `grep -rn 'verify gate PASSED' --include='*.go'`; `grep -rn '5 required identities'`; controls `5/5` and `10/10` | only Go assertion is `ail_binary_gate_test.go:118 passedMarker = "verify gate PASSED"` (a **substring**, via `strings.Contains`); `5 required identities` → **only `verify_ail.sh:378`**; control `5/5` → `module_manifest_gate_test.go:128` (fires); control `10/10` → **0** |
| V15 | **MU5** — the totality guard, and `ai-check`'s rc-blindness, on the exact landing bytes | delete the `EscalateBounded` arm from the body only (n=1 assert), sha-verified landed | `check.passed=**true**`, **ai-check rc=0**, `verify.errors=**1**`, `timeoutOutcome status=**error**`, the other five identities still `verified`. Restored byte-identical |
| V16 | **MU2** — the cardinal consistent lie is Z3-GREEN and killed by exactly one named test | replace the `ExecuteIfGranted`/no-auth arm in **both** contract and body (n=2 assert); sha `f8e08f68…`→`938ab2cd…` | Z3: `verified=6, cex=0, errors=0` — **Z3 verified the lie**. Tests: `n=39, failed=1`, **NON-PASS = `[('outcomeCode_test_5','fail')]` — sole failure.** Restored byte-identical |
| V17 | **MU14** — the round-1 reviewer's early timeout, same signature | `recordedNow >= deadlineAt` → `>= deadlineAt - 1`, **both** (n=2 assert), sha-verified | Z3: `verified=6, cex=0, errors=0`. Tests: `n=39, failed=1`, **NON-PASS = `[('firedCode_test_1','fail')]` — sole failure.** Restored byte-identical |
| V18 | The `AC9`-defect precedent was never written back | `grep -n 'interfaceHash' design_docs/implemented/w-evidence-grade-mapping.md` | `:291` still claims the change moves `interfaceHash`; `:387` still carries *"Fail on … unchanged interface hash"*. The correction exists **only** in that item's *sprint plan* §0.2(i) |

---

## 10. Handoff

**SPRINT_PLAN_PATH**: `design_docs/planned/w-decision-lifecycle-freeze-sprint-plan.md`
**SPRINT_JSON_PATH**: `sprint_w-decision-lifecycle-freeze.json`

Carried forward, **out of scope here** and each already owned:

- the host Timeout-transition emitter, escalation delivery, DEFER's evidence-ref check, and the
  honesty of `independentAuthority` → **item 7**;
- `context.Context` plumbing / bounded store reads → **item 18**;
- flipping HUMAN-SURFACE §7.1/§7.3 to CLOSED, moving the two §8 premise rows, and reconciling the
  item-15/item-7 charter rows (doc §10) → **controller/human acts after this sprint lands**;
- writing the `interfaceHash` correction back into
  `design_docs/implemented/w-evidence-grade-mapping.md` (§1.2 i / V18) → **a controller call**; it
  is `design_docs/`, so AC11 forbids the executor from doing it here.

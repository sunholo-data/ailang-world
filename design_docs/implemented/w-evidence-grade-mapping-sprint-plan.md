# Sprint plan — `w-evidence-grade-mapping` (queue item 13, clause-5)

**Item**: queue item 13 — *a Z3-proven total `gradeOf(Evidence) -> EvidenceGrade` in `world/types.ail`,
the repo's 5th proven identity*.
**Status**: PLANNED · **NO SPLIT** · one milestone, **5 tasks / 2 commits**
**Design doc**: [`design_docs/planned/w-evidence-grade-mapping.md`](w-evidence-grade-mapping.md)
(quorum round 2 clean, both external reviewers present; narrow-refinement carve-out applied —
`EvidenceGrade` is exactly `PROVEN | TESTED | ATTESTED | CLAIMED`)
**Base**: `dev` @ `20ed668`, clean tree (`git status --porcelain` empty, V0)
**Planner**: mission-control iteration 81, opus sprint-planner, first-party measurement on this rig
**Worktree**: a **sibling of the repo** — `/Users/voightkampff/dev/sunholo-data/ailang-world-w13`.
Never `/tmp`: `host/verifygate` and `host/boundary` derive `repoRoot` from `runtime.Caller` and copy
live trees, so a relocated checkout reds for the location rather than the code.

**Headline price: ≈0.75 day, i.e. 1.15× the doc's 0.65 day.** The doc's own arithmetic is close to
right. The overrun is not the AILANG — that was prototyped end to end this session and works exactly
as §3.2 writes it — it is that **the doc's Conflict Surface is missing two artifacts and one of its
acceptance criteria is unsatisfiable by construction** (§0.2). §6 prices it.

**Everything in T1 was PROTOTYPED, RUN and MUTATED end to end this session** (V1–V16). The AILANG
block in §3/T1 is reproduced verbatim from a file that verified and passed 20/20. The new ready-packet
golden in §3/T1.4 is a **pre-registered byte string** — the planner computed it with a known-positive
control that reproduced the *committed* golden byte-for-byte first.

---

## 0. Planner's first-party verification, and what it REFUTED

Every controller premise was re-measured at `20ed668` before anything was planned. **All six
controller premises survived** (§0.1). The design doc did not: **three of its claims are corrected by
measurement**, and two of them would have produced a *red CI* or a *stuck executor* (§0.2).

Full command/observed-output table in §8.

### 0.1 Controller premises — all CONFIRMED

| # | Premise | Verdict |
|---|---|---|
| C1 | Baseline gate GREEN at `20ed668`: rc=0, `✓ verify gate PASSED: 4 required identities verified, 14 named tests pass`, `✓ world package gate PASSED: 9/9 steps performed non-zero work` | **CONFIRMED** (V1) |
| C2 | `/tmp/ailang-v0300/ailang` → `AILANG v0.30.0` | **CONFIRMED** (V1, and step 9/9 pins it by exact bytes on Darwin/arm64) |
| C3 | `Evidence` has exactly five variants at `world/types.ail:23-28` | **CONFIRMED, exact lines** (V2) |
| C4 | The pin map: `EXACT_TOTAL_VERIFIED=4` shell `:310`; `EXACT_TOTAL_TESTS = 14` python **with spaces around `=`** `:340`; `REQUIRED_TESTS` python set `:333-339`; `LEG1_MODULES` a **set-compared** bash array `:135-147`; **no `EXACT_TOTAL_MODULES` exists**; adding the six names to `REQUIRED_TESTS` is a *separate* edit from moving the total | **CONFIRMED in every part**, and the separability warning is now **MEASURED**, not argued: mutation M12 (§4) leaves `len(tests[])==20` and reds *only* under the name pin (V11) |
| C5 | Doc is FRESH from base `6d12a79`: only `host/broker/{handlers.go,handlers_parallel_guard_test.go,handlers_stall_diag_test.go,handlers_test.go}` moved; control 10 files | **CONFIRMED, byte-exact** (V15). Stronger: every one of §8.1's five line citations resolves to the exact construct it names (V3) |
| C6 | Plan against the POST-carve-out doc; `EvidenceGrade` is exactly four constructors | **CONFIRMED and re-derived**: applying §3.2 verbatim to a copy of the live `world/` tree lands sha `91af8cea…` → **`2cf5b004…`**, the identical post-carve-out sha the doc's V25 records (V4) |

One nit, not a refutation: the directive cites `LEG1_MODULES` at `:135-146`; the array *opens* at 135
and *closes* at 147 (entries 136–146). The design doc's `:135-147` is the exact one. Nothing depends
on it — **`LEG1_MODULES` is UNCHANGED by this item**, because no module is added. Stated explicitly
here so nobody "fixes" an allowlist that is already correct: the isolated gate still enumerates
**11 modules** and still emits **11** `ai-check` lines after the change (V6).

### 0.2 Three corrections to the design doc — each MEASURED, each changing the plan

#### (i) **REFUTED, and plan-killing: AC9 is unsatisfiable. `interfaceHash` CANNOT change.**

Doc §3.4: *"The grade type and function change `contentHash`, `interfaceHash`, `tarballSHA256`, and
ordinarily `tarballBytes`."* Doc §3.1: *"its interface and package hashes change."* Doc **AC9**:
*"require content/interface/tar hashes to differ from the old golden … **Fail on** hand-authored JSON,
**unchanged interface hash**, or packet drift."*

`host/pkgproj/pkgproj.go:86` — `InterfaceHash` takes a **`Manifest`**, and reads **only** the package
name, edition, `ailang` constraint, sorted export *module names*, and sorted effects. **It never opens
a source file.** Adding an exported *symbol* to an already-exported *module* is invisible to it. Only
adding/removing/renaming a package export module — which §3.1 spends a page arguing this item must NOT
do — could move it.

Measured (V13), with a known-positive control run **first**: `pkgproj.RecomputeReadyPacket` against the
**unmodified** live projection reproduced the committed golden **byte-for-byte**; against the projection
carrying `gradeOf`:

| field | committed golden | after `gradeOf` | verdict |
|---|---|---|---|
| `contentHash` | `sha256:5ea15858fddc…` | `sha256:489d5e5d47d5…` | **changes** |
| `interfaceHash` | `sha256:d16cc88270ff…` | `sha256:d16cc88270ff…` | **UNCHANGED** |
| `tarballSHA256` | `sha256:a32806a069bb…` | `sha256:d0cdf42be80e…` | **changes** |
| `tarballBytes` | `5773` | `6236` | **changes** |

**Consequence.** An executor obeying AC9 literally has two exits, and both are disasters. Declare the
sprint failed on a correct implementation; or "make interfaceHash differ" by hand-editing the golden —
which AC9's *own next clause* forbids, and which step 9/9's `cmp -s` against the recomputed packet
would red anyway. **AC9 is AMENDED in §2**: three of four packet fields move, `interfaceHash` is
required to be **byte-IDENTICAL**, and that invariance is itself asserted (it is the evidence that no
package export was added — the property §3.1 actually cares about).

#### (ii) **REFUTED: the Conflict Surface omits a Go test that this change REDS. CI goes red without it.**

`host/verifygate/module_manifest_gate_test.go:128`:

```go
const marker = "✓ 4/4 required world/ identities verified across 11 module(s)"
```

`newIsolatedGateRoot` (`:47-88`) copies **the live `scripts/verify_ail.sh`, the live `world/*.ail` and
the live `design_docs/sketches/*.ail`** into a temp root and runs the gate there;
`requirePristineControl` then requires that exact substring. Five tests call it (`:172, 215, 241, 257,
286`).

Measured, not inferred (V6): a faithful reconstruction of `newIsolatedGateRoot` — **13 files / 11 `.ail`,
matching the Go test's own `files != 13 || ailFiles != 11` assertion** — carrying the §3.2 edit and the
script pins prints:

```
   ✓ 5/5 required world/ identities verified across 11 module(s)
   ✓ all 20 required named tests pass (failed_tests=0)
```

`strings.Contains(out, "…4/4…11 module(s)")` is then **false**, and `requirePristineControl` `t.Fatalf`s.
So `go test ./host/verifygate/` — i.e. CI's go-verify job — **reds** unless line 128 moves to `5/5` in
the same commit. The doc's §8 Conflict Surface never mentions `host/verifygate/`; the doc's **AC11**
says *"No production Go source change is expected"* (true, narrowly — this is a *test*) and **AC12**
says *"changes only the four metadata files"* (**false — there are five**).

**Consequence:** `host/verifygate/module_manifest_gate_test.go` is added to the atomic commit set, and
AC12's file list is amended from four to five in §2. Note the count is a **substitution, not a
widening**: the `11 module(s)` half of the same literal must stay `11`.

#### (iii) **REFUTED, and the one a green gate will hide: `verify_ail.sh:376` is a hardcoded lie-in-waiting.**

```bash
echo "✓ verify gate PASSED: 4 required identities verified, 14 named tests pass"
```

Line 376 is a **literal**. Lines 315 and 370 interpolate `EXACT_TOTAL_VERIFIED` / `EXACT_TOTAL_TESTS`;
376 does not. The doc's §8.1 conflict surface lists `:135-147`, `:262-267`, `:278-282`, `:310-314`,
`:317-340` — **376 is outside every one of them**, and the controller's pin map omits it too.

Left alone the gate still exits 0, so **nothing catches it**. Measured (V14): the only Go assertion on
this line, `host/verifygate/ail_binary_gate_test.go:118`, is
`passedMarker = "verify gate PASSED"` — a **substring**, used by `strings.Contains`. A stale `4 … 14`
satisfies it.

That is precisely the shape this mission keeps paying for: the gate's own terminal summary would
announce `4 required identities verified, 14 named tests pass` while verifying 5 and running 20 — and
it is the line the controller's baseline quotes and every future iteration reads. **This is an
UNGUARDED pin.** §4 records it honestly as such rather than inventing a mutation that "kills" it: it
gets a grep-based acceptance criterion (AC13) with a control, not a fake non-vacuity arm.

### 0.3 Not refutations — three traps for the arm author

1. **`RecomputeReadyPacket` needs `Version` and the repo's only manifest literal omits it.**
   `verify_world_package.sh:175` builds `pkgproj.Manifest{Package: {Name, Edition, AILANG}, …}` with
   **no `Version`** — invisible there, because that helper only prints hashes. `RecomputeReadyPacket`
   copies `manifest.Package.Version` into the packet's `"version"` field, so a regeneration helper
   that copies `:175` verbatim emits `"version":""` and reds step 9/9 with a *hash-shaped* diff for a
   *manifest-shaped* reason. The §3/T1.4 helper sets `Version: "0.1.0"` and is proven by the
   reproduce-the-committed-golden control (V13).
2. **`ailang test` reports two different totals for the same tree.** Repo-root Leg 2 parses
   `len(tests[])` = **20**; the package leg's human output says `26 tests: 20 passed, 0 failed,
   6 skipped` (V12). The 6 are contract-derived *properties*, which is exactly why `verify_ail.sh:320`
   gates on `len(tests[])` and not `passed_tests`. Step 5/9 has **no count pin** (`grep -Ec
   'PASS|pass|test' > 0`), so nothing to move there — but an arm author reading `26` and "fixing"
   `EXACT_TOTAL_TESTS` to 26 would red the gate.
3. **`ai-check` returns rc=0 on a totality failure.** Measured again this session (V5): mutation M1
   gives `verify.errors=1`, `gradeOf status=error`, `non-exhaustive pattern match` — with
   `check.passed` still `true` and **rc still `0`**. No acceptance criterion in this sprint may read an
   exit code to prove totality. Every one below parses JSON or reads the gate's parsed refusal line.

---

## 1. What this sprint is, and what it explicitly does NOT do

**Is**: add `EvidenceGrade` (four constructors) and a Z3-proven total `gradeOf(Evidence)` to
`world/types.ail`, plus a private `gradeCode` adapter carrying six inline cases; pin `gradeOf` and all
six `gradeCode_test_N` identities **by name** in `scripts/verify_ail.sh`; move the four count/summary
pins; rebuild the package projection and the ready-packet golden; move the one Go marker; prove the
whole thing non-vacuous with 11 landed mutations.

**Is NOT** — carried verbatim from doc §10, and each is an executor refusal:

- No `ProofReport`, `ReplayReport`, or any other evidence carrier. **No `PROVEN` arm.** `PROVEN` stays
  deliberately unreachable; mint authority belongs to queue item 17
  (`w-validated-proven-evidence-boundary`) and **must not be started here**.
- No `UNSUPPORTED` grade constructor — the round-2 carve-out removed it.
- No new module (`LEG1_MODULES` unchanged, 11), no new package export (`interfaceHash` unchanged), no
  new tar entry (6), no Go API, no renderer, no effect, no host codec, no producer wiring.
- No `HashRef` loading, validation, decoding or repair. No aggregation over `Proposal.evidence`. No
  grading of `Proposal.confidence`.
- No change to `Evidence` itself — it keeps exactly its five constructors.

---

## 2. Acceptance criteria

The doc's **AC1–AC12 are adopted**, with three amendments and one addition. Every command below was
**baselined at `20ed668` first**; an AC that is already red at base is a broken AC, not a defect.

### Amendments (each forced by a measurement in §0.2)

| AC | Amendment |
|---|---|
| **AC9** | **`interfaceHash` MUST be byte-IDENTICAL**, not different. Require exactly three moved fields (`contentHash`, `tarballSHA256`, `tarballBytes`) and assert `interfaceHash` unchanged as the positive evidence that no package export was added. The doc's "unchanged interface hash ⇒ fail" clause is **struck**. |
| **AC12** | File list is **five**, not four: `world/types.ail`, `packages/world-core/world/types.ail`, `scripts/verify_ail.sh`, `scripts/world_package_ready_packet.golden.json`, **`host/verifygate/module_manifest_gate_test.go`**. Design-doc prose repairs (T5) land in a **separate commit** and are excluded from this count. |
| **AC4** | Add: moving `EXACT_TOTAL_TESTS` 14 → 20 **without** adding the six names to `REQUIRED_TESTS` is a **weakening, not a partial fix** — measured green under M12. Both edits are required, and M12 is the arm that proves it. |

### New

- **AC13 (new).** `scripts/verify_ail.sh:376`'s terminal banner reads
  `✓ verify gate PASSED: 5 required identities verified, 20 named tests pass`.
  This pin is **NOT gate-enforced** (§0.2 iii); it is asserted by grep with a control.

### Hold set — re-measured at every task exit and at T4

| Invariant | Value at base (`20ed668`) |
|---|---|
| `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | rc=0 · `✓ 4/4 … across 11 module(s)` · `✓ all 14 required named tests pass` · `9/9 steps` → **must become `5/5 … 11 module(s)` / `20` / `9/9`** |
| `LEG1_MODULES` | 11 entries, set-compared — **must not move** |
| isolated-root file census (`module_manifest_gate_test.go:85`) | `13` files / `11` `.ail` — **must not move** |
| `\n   ai-check ` line count in the isolated control | `11` — **must not move** |
| package tar entries (step 8/9) | `6` — **must not move** |
| package exports (steps 4/9, 7/9, 9/9) | `4` — **must not move** |
| `interfaceHash` | `sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083` — **must not move** |
| `GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` | rc=0. **Base condition, not a regression:** the rig's `go` is `go1.26.4`; without `GOTOOLCHAIN=go1.25.6` this script FATALs on a `host/store/scan.go` miscompile. Set it before reading the result. |
| `host/boundary` `wantFileCount` | `1` — untouched; `ailArmMutate` appends a comment to `world/types.ail` and is content-agnostic (V16) |

---

## 3. Task breakdown — 5 tasks, 2 commits, one milestone

### T1 — the atomic five-file landing (**COMMIT 1**)

Nothing in T1 may be committed alone. Four of the five files are *derived from* the first; a partial
landing reds the gate for a reason unrelated to the change, and the executor then debugs the wrong
thing. Order below is causal, not optional.

#### T1.1 — `world/types.ail` (+47 lines, 85 → 132)

Insert **immediately after** `  | RecordedEffect(HashRef)` (the last `Evidence` arm, `:28`), leaving
`Evidence` byte-for-byte unchanged. Verbatim; this exact text verified and passed 20/20 (V4):

```ailang
-- The four ratified trust grades. gradeOf is total over decoded Evidence and
-- returns one of these four; absent/unreadable/malformed input is NOT a grade
-- and belongs to the deferred validated-boundary follow-on.
export type EvidenceGrade
  = PROVEN
  | TESTED
  | ATTESTED
  | CLAIMED

-- Canonical grading for the ratified five-constructor representation.
-- No current Evidence constructor has authority to mint PROVEN.
export func gradeOf(e: Evidence) -> EvidenceGrade ! {}
ensures { result == match e {
  CompilerOutput(_) => ATTESTED,
  TestReport(_, _) => TESTED,
  HumanApproval(_) => ATTESTED,
  AiReview(_, _) => CLAIMED,
  RecordedEffect(_) => ATTESTED
} }
{
  match e {
    CompilerOutput(_) => ATTESTED,
    TestReport(_, _) => TESTED,
    HumanApproval(_) => ATTESTED,
    AiReview(_, _) => CLAIMED,
    RecordedEffect(_) => ATTESTED
  }
}

-- Private adapter because v0.30.0 inline-test expected values cannot be ADT
-- identifiers. Integer codes are test plumbing, not an exported ordering.
func gradeCode(e: Evidence) -> int
tests [
  (CompilerOutput({ algo: "sha256", digest: "compiler" }), 2),
  (TestReport({ algo: "sha256", digest: "tests-pass" }, true), 3),
  (TestReport({ algo: "sha256", digest: "tests-fail" }, false), 3),
  (HumanApproval({ algo: "sha256", digest: "approval" }), 2),
  (AiReview({ algo: "sha256", digest: "review" }, 0.8), 1),
  (RecordedEffect({ algo: "sha256", digest: "effect" }), 2)
]
{
  match gradeOf(e) {
    PROVEN => 4, TESTED => 3, ATTESTED => 2,
    CLAIMED => 1
  }
}
```

**Landed-proof (mandatory, doc V26's lesson):** `shasum -a 256 world/types.ail` must read
`91af8cea316ac177…` **before** and **`2cf5b004f7f0573f…`** after. If the post sha is anything else,
**STOP** — the golden byte string pre-registered in T1.4 was computed from `2cf5b004…` and will not
match. A green print from a file that was never written is indistinguishable from a green print from a
file that was.

#### T1.2 — `scripts/verify_ail.sh` — **five** edits, all in this commit

| # | Line | From | To | Mechanism |
|---|---|---|---|---|
| 1 | `266` | `"world/types.ail":       set(),` | `"world/types.ail":       {"gradeOf"},` | python dict, Leg 1 named proof pin |
| 2 | `310` | `EXACT_TOTAL_VERIFIED=4` | `EXACT_TOTAL_VERIFIED=5` | **shell** var, no spaces |
| 3 | `333-339` | `REQUIRED_TESTS = { … }` | append the six `gradeCode_test_1`…`_6` names | python **set**, Leg 2 named test pin |
| 4 | `340` | `EXACT_TOTAL_TESTS = 14` | `EXACT_TOTAL_TESTS = 20` | **python** var, **spaces around `=`** — a shell-style `grep 'EXACT_TOTAL_TESTS='` misses it |
| 5 | `376` | `… 4 required identities verified, 14 named tests pass` | `… 5 required identities verified, 20 named tests pass` | **literal**, not interpolated; **unguarded** (AC13) |

`LEG1_MODULES` (`:135-147`) is **UNCHANGED**. Edits 3 and 4 are separate obligations: edit 4 alone
still passes (n==20 holds, all 14 old names present) while leaving all six new identities unpinned by
name — **measured green under M12** (V11).

#### T1.3 — `packages/world-core/world/types.ail`

```bash
./scripts/build_world_package.sh
```

Wholesale-replaces the projection from a fresh staging dir. Then assert equality directly rather than
waiting for step 3/9:

```bash
cmp -s world/types.ail packages/world-core/world/types.ail && echo "projection EQUAL"
cmp -s world/types.ail world/contracts.ail; echo "control (must be 1): $?"
```

#### T1.4 — `scripts/world_package_ready_packet.golden.json`

Do **not** hand-edit. Regenerate through `host/pkgproj`, and run the **known-positive control first**
(the planner did; it reproduced the committed golden byte-for-byte, V13):

```go
// /tmp/regen_golden.go — throwaway, never committed
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
# CONTROL FIRST, at HEAD before T1.1: must reproduce the COMMITTED golden byte-for-byte.
git stash && go run /tmp/regen_golden.go packages/world-core > /tmp/ctl.json && git stash pop
cmp -s /tmp/ctl.json scripts/world_package_ready_packet.golden.json && echo "CONTROL PASS"
# then, after T1.3:
go run /tmp/regen_golden.go packages/world-core > scripts/world_package_ready_packet.golden.json
```

**Pre-registered expected result** (planner-measured, V13). The executor must land exactly this; a
mismatch means the `world/types.ail` bytes differ from §3.2 and is a **STOP**, not a golden to accept:

```json
{"compilerVersion":"AILANG v0.30.0","contentHash":"sha256:489d5e5d47d5c3443ca698f69c6250db57831f463cad8abbec8510ae55ecf632","effects":[],"exports":["world/types","world/contracts","world/transitions","world/logepoch"],"interfaceHash":"sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083","package":"world/core","tarballBytes":6236,"tarballSHA256":"sha256:d0cdf42be80eee2de043645a6bcdb491cb439e3ba659b08b3ff67ee797e99001","version":"0.1.0"}
```

Note `interfaceHash` is **identical** to the committed golden's (§0.2 i). Two independent
implementations then cross-check these bytes: `verify_world_package.sh:236-244`'s python emitter
(`cmp -s`) and `host/pkgproj/readypacket_test.go`'s
`TestEncodeReadyPacketReproducesTheCommittedGoldenBytes`, which reads the golden and re-encodes it.

#### T1.5 — `host/verifygate/module_manifest_gate_test.go:128`

```diff
-	const marker = "✓ 4/4 required world/ identities verified across 11 module(s)"
+	const marker = "✓ 5/5 required world/ identities verified across 11 module(s)"
```

`11 module(s)` **stays 11**. This is a substitution of one digit-pair, not a widening. `TestNoRigAbsolutePaths`
globs `host/verifygate/*.go` — introduce no rig-absolute literal (none needed here).

**COMMIT 1** = exactly these five files, message referencing item 13 and the doc.

### T2 — the mutation sweep (11 mutations, **NO COMMIT**)

The bulk of the time. Protocol in §4.1. Every mutant is restored **byte-identical by sha256** from a
`cp` backup — never `git checkout -- <file>`, which would delete uncommitted work and then report the
disaster rather than prevent it.

### T3 — full pinned gates and scope inspection

### T4 — hold-set re-measurement and AC7/AC13 grep arms

### T5 — narrow prose staleness repair (**COMMIT 2**, `design_docs/` only)

Doc §8.5 requires it: HUMAN-SURFACE's *"no total mapping"* statement becomes stale on landing, while
its *"no `PROVEN` carrier"* statement stays true. Bounds, and they are tight:

- **Only** the mapping-gap sentences. Mark the mapping gap **closed**; leave every `PROVEN`/carrier
  statement **untouched and true**.
- Targets: `design_docs/HUMAN-SURFACE.md` and `design_docs/world-mission.md` (§8.5 / V18 cites
  `world-mission.md:2247-2253, 2781-2786`).
- **Do not** silently rewrite `design_docs/implemented/*` — historical docs stay historical.
- **Do not** restate or extend policy. If a repair would require asserting anything new about
  authority, grades, or the renderer, **STOP and escalate** — `HUMAN-SURFACE.md` is ratified surface,
  and widening it is a human gate, not an executor edit. (This is the iter-79 spine: a repo's own
  recorded limitation is a claim, and it was a *ratified* doc that silently widened one last time.)

---

## 4. Mutation discipline

### 4.1 Protocol — every arm, no exceptions

1. `cp` the target to `/tmp/w13_backup/` **before** the first mutation.
2. Apply the mutation. **Assert it LANDED by sha256** (before ≠ after). A mutation that never applied
   and a mutation that failed to red share an exit code — doc V26 paid for this exact lesson.
3. Run the observable. Record the **exact** refusal line, not just rc.
4. Restore by `cp` from the backup. **Assert byte-identical by sha256.**
5. Re-run the positive control and require GREEN before the next arm.

### 4.2 Rule 3j — the refusal branches this change installs

Five new refusal branches: (B1) Leg-1 `verify.errors` on a non-total match; (B2) Leg-1
`verify.counterexample` on a body/spec divergence; (B3) Leg-1 `REQUIRED_VERIFIED` **name** lookup;
(B4) Leg-2 `REQUIRED_TESTS` **name** lookup; (B5) Leg-3 steps 3/9 and 9/9. Each gets ≥1 mutation.
`EXACT_TOTAL_VERIFIED` and `EXACT_TOTAL_TESTS` are **secondary count guards**, not branches — M12 and
M13 are specifically the arms that pass the counts and red only the names.

### 4.3 The 11 mutations — and the observable each one reads (rule 3i)

Every observable below is **downstream of the mechanism**: each reads the gate's own parsed verdict on
freshly recomputed compiler/runner output, never a value written alongside the edit. Rows marked
**MEASURED** were run by the planner this session with a landed-proof sha and a restore control.

| ID | Mutation (lands by sha256) | Branch | Observable the arm reads | Measured result |
|---|---|---|---|---|
| **M1** | add `\| FutureEvidence(HashRef)` to `Evidence`, no `gradeOf` arm | B1 | Leg-1 python parse of `ai-check` JSON → `verify.errors` | **MEASURED** rc=1 · `✗ world/types.ail: verify.errors == 1 (Z3 encoding error, silent under exit codes V10)`. **`check.passed` stays `true`, `ai-check` rc stays `0`** — the exit code proves nothing |
| **M1-ctl** | M1 + an explicit `FutureEvidence(_) => ATTESTED` arm in both contract and body | B1 control | same | **MEASURED** (doc V25 shape) `gradeOf verified`, `errors=0` |
| **M2a** | **body-only** `CompilerOutput(_) => TESTED` (4-space indent; unique anchor) | B2 | Leg-1 `verify.counterexample` **and** Leg-2 named test | **MEASURED** `counterexample=1`, `gradeOf status=counterexample`; `gradeCode_test_1=fail` |
| **M2b** | **contract AND body** `CompilerOutput(_) => TESTED` (both occurrences) | B4 | Leg-2 named test, with **Leg 1 GREEN** | **MEASURED** rc=1 · Leg 1 `✓ 5/5 … 11 module(s)` · `✗ test leg: required named tests missing/failing: gradeCode_test_1=fail`. **This is the arm that proves Leg 2 is load-bearing independently of the proof.** The doc's M2 row is silent on which variant, and the body-only reading makes its AC6 ("require **only** its named case to fail") literally false |
| **M7** | **contract AND body** `CompilerOutput(_) => PROVEN` | B4 | Leg-2 named test | **MEASURED** Leg 1 **GREEN** (`verified`, `errors=0`, `cex=0`); `gradeCode_test_1=fail`. **Load-bearing finding: a `PROVEN` arm added consistently to spec *and* body is INVISIBLE to Z3.** The no-`PROVEN` property is enforced by the six integer expectations and AC7's grep — **not** by the proof. Say so; do not claim the contract forbids `PROVEN` |
| **M3–M6** | `HumanApproval`→`CLAIMED`; `TestReport`→`ATTESTED`; `AiReview`→`ATTESTED`; `RecordedEffect`→`CLAIMED` — each in **both** contract and body | B4 | Leg-2 named tests; each arm names the failing identity **and** the sibling identities that stay `pass` | to run (M2b/M7 establish the signature: Leg 1 green, exactly the affected identities fail; `TestReport` fails **both** `_test_2` and `_test_3`) |
| **M12** | rename `gradeCode` → `gradeKode` (all occurrences) | B4 | Leg-2 `REQUIRED_TESTS` **name** lookup, with `len(tests[]) == 20` **unchanged** | **MEASURED, and this is the controller's Fact-4 warning made real.** Under **count-only** pinning (`EXACT_TOTAL_TESTS=20`, `REQUIRED_TESTS` untouched): **GREEN**, n=20, `failed_tests=0`. Under **name** pinning: rc=1 · `✗ test leg: required named tests missing/failing: gradeCode_test_1=MISSING, …_2=MISSING, …_3=MISSING, …_4=MISSING, …_5=MISSING, …_6=MISSING` |
| **M13** | rename `gradeOf` → `gradeOfX` (all occurrences) | B3 | Leg-1 `REQUIRED_VERIFIED` **name** lookup, with `total_verified == 5` **unchanged** | **MEASURED** `verify.verified` stays **1** (so the 4→5 count guard passes) yet rc=1 · `✗ world/types.ail: required identity (types, gradeOf) MISSING from verify.results[] (vanished silently, V20)`. **This is the arm that proves the named proof pin is non-vacuous independently of the count.** |
| **M14** | delete the `ensures { … }` block, keep `gradeOf` and its body | B3 | Leg-1 name lookup **and** count | **MEASURED** `verify.verified=0`, `results[]` empty; rc=1 · same `required identity (types, gradeOf) MISSING` line. Weaker than M13 (the count fires too) — kept because it is the "contract silently dropped" case `verify_ail.sh:8` was written for |
| **M9** | edit canonical `world/types.ail`, **skip** `build_world_package.sh` | B5 | Leg-3 step 3/9 sha256 compare | **MEASURED as an arithmetic fact**: canonical `2cf5b004…` vs projected `91af8cea…` differ, and step 3/9 is a `sha256` equality → `✗ projection mismatch: world/types.ail=… packages/world-core/world/types.ail=…`. Executor must still run it live |
| **M10** | rebuild the projection, keep the **old** golden | B5 | Leg-3 step 9/9 `cmp -s` | **MEASURED as an arithmetic fact** (V13): the recomputed packet and the committed golden differ in 3 of 8 fields → `✗ ready packet differs byte-for-byte from golden` + a `diff -u` naming `contentHash`/`tarballSHA256`/`tarballBytes` |
| **M15 (new)** | revert `module_manifest_gate_test.go:128` to `4/4` | (ii) | `go test ./host/verifygate/` | **CONFIRMED by construction from a measured string** (V6): the isolated control emits `✓ 5/5 required world/ identities verified across 11 module(s)`; `strings.Contains` of the `4/4` literal is false → `pristine isolated control missing "✓ 4/4 …"`. Executor must run it live |

### 4.4 The unguarded pin — stated, not faked

`verify_ail.sh:376` has **no killing mutation and no gate**. Reverting it to `4 … 14` leaves the gate
rc=0 and satisfies `ail_binary_gate_test.go`'s substring `passedMarker` (V14). Writing a mutation row
claiming otherwise would be exactly the decorative arm rule 3i forbids. It gets AC13's grep instead:

```bash
grep -c '✓ verify gate PASSED: 5 required identities verified, 20 named tests pass' scripts/verify_ail.sh   # want 1
grep -c '4 required identities verified\|14 named tests pass' scripts/verify_ail.sh                          # want 0
grep -c 'verify gate PASSED' scripts/verify_ail.sh                                                           # control, want 1
```

---

## 5. Acceptance commands, as the executor must run them

Each is annotated with its **baselined-at-`20ed668`** result.

```bash
export AILANG_BIN=/tmp/ailang-v0300/ailang

# AC10 — the full pinned gate. BASE: rc=0, "✓ 4/4 … across 11 module(s)",
#        "✓ all 14 required named tests pass", "9/9 steps performed non-zero work".
#        AFTER: rc=0, "✓ 5/5 … across 11 module(s)", "✓ all 20 …", "9/9".
AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh; echo "rc=$?"

# AC3 — the named proof pin, read from JSON (NEVER from rc: ai-check exits 0 on a Z3 error).
#        BASE: gradeOf absent entirely.  AFTER: gradeOf=verified, errors=0, counterexample=0.
$AILANG_BIN ai-check -timeout 5s world/types.ail > /tmp/t.json
python3 -c "import json;d=json.load(open('/tmp/t.json'));v=d['verify'];print(d['check']['passed'],v['verified'],v['errors'],v['counterexample'],[(r['function'],r['status']) for r in v['results']])"

# AC4 — the six named runtime identities. BASE: n=14, zero gradeCode names.
#        AFTER: n=20, gradeCode_test_1..6 all "pass".
$AILANG_BIN test --format json world/ > /tmp/w.json
python3 -c "
import json;raw=open('/tmp/w.json').read();d=json.loads(raw[raw.find('{'):])
t=d['tests'];print(len(t),d['failed_tests'],sorted(x['name'] for x in t if x['name'].startswith('gradeCode_')))"

# AC8 — projection. BASE: cmp rc=0 (already equal), control rc=1.
./scripts/build_world_package.sh
cmp -s world/types.ail packages/world-core/world/types.ail; echo "projection=$?"   # want 0
cmp -s world/types.ail world/contracts.ail;                echo "control=$?"       # want 1

# AC9 (AMENDED) — 3 of 4 packet fields move; interfaceHash MUST NOT.
#        BASE: content 5ea15858…, interface d16cc882…, tar a32806a0…, bytes 5773.
#        AFTER: content 489d5e5d…, interface d16cc882… (SAME), tar d0cdf42b…, bytes 6236.
python3 -c "
import json
g=json.load(open('/tmp/golden.base.json')); n=json.load(open('scripts/world_package_ready_packet.golden.json'))
moved=[k for k in ('contentHash','tarballSHA256','tarballBytes') if g[k]!=n[k]]
assert len(moved)==3, moved
assert g['interfaceHash']==n['interfaceHash'], 'interfaceHash MOVED — a package export was added'
assert len(n['exports'])==4 and n['effects']==[]
print('AC9 OK: moved', moved, '| interfaceHash invariant')"

# AC11 — Go gate. BASE: rc=0 WITH GOTOOLCHAIN. Without it: FATAL (base condition, not a regression).
GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh; echo "rc=$?"

# AC2/AC7 — no top-grade mint, no withdrawn carriers. Counts below are PLANNER-MEASURED
# on the exact §3.2 text; BASE is 0 for every one of them (the words do not occur at all).
grep -cE '=> *PROVEN'  world/types.ail            # want 0 — NO ARM RETURNS PROVEN. This is the load-bearing one.
grep -cE '=> *ATTESTED' world/types.ail           # KNOWN-POSITIVE CONTROL for the same instrument, want 6 (3 spec + 3 body)
grep -c  'PROVEN'      world/types.ail            # want exactly 3: the `= PROVEN` declaration, the comment, and
                                                  #   `PROVEN => 4` in gradeCode — where PROVEN is a PATTERN, left of `=>`,
                                                  #   not a result. Do not "simplify" this to `want 0`.
grep -cE 'ProofReport|ReplayReport|UNSUPPORTED' world/types.ail   # want 0 — withdrawn carriers and the carved-out sentinel
grep -c  'ATTESTED'    world/types.ail            # want 8 (1 decl + 3 spec + 3 body + 1 gradeCode pattern)

# Instrument non-vacuity for the AC7 grep, run once against a THROWAWAY mutant (never the live file):
#   sed 's/CompilerOutput(_) => ATTESTED,/CompilerOutput(_) => PROVEN,/g' world/types.ail > /tmp/m7.ail
#   grep -cE '=> *PROVEN' /tmp/m7.ail    # MEASURED: 2 — the instrument fires. A `want 0` that cannot
#                                        # reach 1 is decoration, not a gate.

# AC1 — Evidence is byte-unchanged. MEASURED identical at base and after: sha256 prefix cdf389ae6302b183.
sed -n '23,28p' world/types.ail | shasum -a 256   # must be cdf389ae6302b183…

# AC12 — exactly five files, and nothing else.
git status --porcelain    # want exactly the 5 files of COMMIT 1

# AC13 (new) — the unguarded terminal banner (§4.4 greps).
```

---

## 6. Estimate — and why the doc's 0.65 d becomes ≈0.75 d

| Work | Doc | Plan | Delta driver |
|---|---:|---:|---|
| Grade ADT, contract/function, six cases | 0.15 d | **0.05 d** | prototyped verbatim this session; it is a paste + sha assert |
| Move named proof/test pins; RED/control mutations | 0.20 d | **0.30 d** | doc has 11 rows; plan has 11 **arms with landed-proof + restore controls**, and adds M12/M13/M15 — the three that prove the *name* pins non-vacuous |
| Rebuild projection and canonical golden | 0.15 d | **0.10 d** | recipe prototyped with a reproduce-the-committed-golden control; expected bytes pre-registered |
| Run pinned AILANG/Go gates and inspect | 0.15 d | **0.15 d** | unchanged |
| **The two artifacts the doc's Conflict Surface omits** (`module_manifest_gate_test.go:128`, `verify_ail.sh:376`) | — | **0.05 d** | §0.2 (ii), (iii) |
| **AC9 amendment** — reconciling the doc against `InterfaceHash` | — | **0.05 d** | §0.2 (i); undiagnosed this is an open-ended stall, not 0.05 d |
| Narrow prose staleness repair (T5) | — | **0.05 d** | doc §8.5 requires it; doc §9 does not price it |
| **Total** | **0.65 d** | **≈0.75 d** | |

Still inside the queued ~0.5–1 day band.

### 6.1 Verdict: **NO SPLIT**

The five files of COMMIT 1 are one atomic unit *by gate construction*, and every candidate seam fails:

- *"types.ail first, script pins second"* — reds Leg 1 immediately (`✗ expected exactly 4 proven
  world/ contracts, got 5`) and Leg 3 step 3/9. Not a landable commit.
- *"code first, golden second"* — reds step 9/9 by design; that is mutation **M10**, not a milestone.
- *"AILANG first, Go marker second"* — reds `go test ./host/verifygate/`, i.e. CI's go-verify job.
- *"pins first, mutations second"* — lands the pin with **no** committed non-vacuity proof. That is
  the guard-not-a-gate failure this repo has now paid for three iterations running.

T5 is a genuine seam (design-doc prose, zero gate coupling) and is therefore its own commit.

---

## 7. Execution protocol and risks

### Protocol

- Sibling worktree `/Users/voightkampff/dev/sunholo-data/ailang-world-w13`; **never `/tmp`**.
- `export AILANG_BIN=/tmp/ailang-v0300/ailang` in every shell; `GOTOOLCHAIN=go1.25.6` for any `go`
  command that touches `host/store`.
- Snapshot backups to `/tmp/w13_backup/` before T2; restores are **`cp`**, never `git checkout --`.
- Save the base golden to `/tmp/golden.base.json` **before** T1.4 — AC9 needs the old bytes.
- `.probe/` and `/tmp/regen_golden.go` are **measurements, not artifacts**: neither is committed. A
  stray `.ail` under `world/` would red the `LEG1_MODULES` set compare, and a stray file under
  `packages/world-core` would red step 2/9.

### Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Executor "fixes" `LEG1_MODULES` because the pin map mentions it | medium | §0.1 and §1 say **UNCHANGED, 11** three times; the isolated census 13/11 and the 11 `ai-check` lines are in the hold set |
| Executor tries to make `interfaceHash` differ per AC9 | **high without §0.2 (i)** | AC9 amended; `interfaceHash` invariance is in the hold set as a **must-not-move** |
| `grep 'EXACT_TOTAL_TESTS='` misses line 340 (spaces around `=`) | medium | T1.2 row 4 names the mechanism; AC4's `n==20` catches the omission |
| Executor reads `26 tests` from the package leg and "fixes" the total to 26 | medium | §0.3 (2); step 5/9 has no count pin |
| Golden regenerated with an empty `"version"` | medium | §0.3 (1); the reproduce-the-committed-golden control catches it before T1.1 |
| A mutation "passes" because it never landed | medium | §4.1 step 2 — sha256 landed-proof on every arm; doc V26 is the precedent |
| Scope creep into item 17 (`w-validated-proven-evidence-boundary`) | low | §1; AC7's grep arms; no carrier may land without the enforcement boundary |
| `host/broker` ~18% base flake reddens `verify_go.sh` | low-medium | known base condition; re-run and compare against the base measurement, do not "fix" broker |

---

## 8. Verification Log

All commands run from the repo root at `20ed668`, 2026-08-13, instrument
`/tmp/ailang-v0300/ailang` (`AILANG v0.30.0`). Semantic verdicts come from JSON or from the gate's
parsed refusal line, never from process rc. Probes ran in `mktemp -d` directories; **the repo was
never modified**.

| ID | Claim | Command | Observed |
|---|---|---|---|
| V0 | Base is clean at `20ed668` | `git rev-parse --short HEAD; git status --porcelain` | `20ed668`; empty |
| V1 | Baseline gate GREEN (controller C1/C2) | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | rc=0; `✓ 4/4 required world/ identities verified across 11 module(s)`; `✓ all 14 required named tests pass (failed_tests=0)`; `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED: 4 required identities verified, 14 named tests pass` |
| V2 | `Evidence` = five constructors at `:23-28` | `Read world/types.ail` | exactly the five; file is 85 lines |
| V3 | §8.1's five line citations are live and exact | `Read scripts/verify_ail.sh` | `LEG1_MODULES` 135–147 (11 entries); `REQUIRED_VERIFIED` 262–267 with `"world/types.ail": set()` at 266; errors/counterexample 278–282; `EXACT_TOTAL_VERIFIED=4` at 310; Leg 2 317–340 with `EXACT_TOTAL_TESTS = 14` at 340 |
| V4 | §3.2's literal code verifies and passes, on a copy of the real `world/` tree | `cp world/*.ail $P/world/`; baseline `ailang test --format json world/`; apply §3.2; `ai-check -timeout 5s world/types.ail`; `ailang test --format json world/`; sha before/after | **baseline control `n=14 passed=14 failed=0`**; landed `91af8cea…`→**`2cf5b004…`** (identical to doc V25's post-carve-out sha); `gradeOf verified`, `verify={verified:1, errors:0, counterexample:0}`; `n=20 passed=20 failed=0`; identities exactly `gradeCode_test_1`…`_6` |
| V5 | M1 reds Leg 1 with rc=0 from the compiler | apply M1; `ai-check` | `check.passed=true`, **rc=0**, `verify.errors=1`, `gradeOf status=error`, `non-exhaustive pattern match`; Leg 2 **stays n=20 failed=0** |
| V6 | Faithful `newIsolatedGateRoot` reconstruction + both pin edits ⇒ `5/5`, `20` | build 13-file/11-`.ail` iso root; apply 5 script edits with count-1 anchor asserts; run gate | `files=13 ail=11`; `✓ 5/5 required world/ identities verified across 11 module(s)`; `✓ all 20 required named tests pass (failed_tests=0)`; 11 `ai-check` lines. **Establishes (ii): the `4/4` literal no longer matches** |
| V7 | M2a — body-only | apply; ai-check + test | `counterexample=1`, `gradeOf status=counterexample`; `gradeCode_test_1=fail` |
| V8 | M2b — contract **and** body | apply; gate | Leg 1 `✓ 5/5 … 11 module(s)`; `✗ test leg: required named tests missing/failing: gradeCode_test_1=fail`; rc=1 |
| V9 | M7 — a consistent `PROVEN` arm is invisible to Z3 | apply; gate | `verified=1, errors=0, counterexample=0`, `gradeOf verified`; only `gradeCode_test_1=fail` |
| V10 | M13 — count passes, **name** catches | rename `gradeOf`→`gradeOfX`; gate | `verify.verified` stays `1`; `✗ world/types.ail: required identity (types, gradeOf) MISSING from verify.results[] (vanished silently, V20)`; rc=1 |
| V11 | M12 — the count-only pin is **measurably** a weakening | rename `gradeCode`→`gradeKode`; run Leg-2 parser under both pin sets | count-only (`EXACT_TOTAL_TESTS=20`, `REQUIRED_TESTS` untouched): **GREEN**, n=20, `failed_tests=0`. Name-pinned: `✗ test leg: required named tests missing/failing: gradeCode_test_1=MISSING … _6=MISSING`, rc=1. Emitted names were `gradeKode_test_1…_6` |
| V12 | M14, and the package-leg total trap | delete `ensures`; also `ailang test world/` inside the package copy | M14: `verify.verified=0`, `results[]` empty, same MISSING line. Package leg prints `26 tests: 20 passed, 0 failed, 6 skipped`; `check --package` → `✓ 5 files checked, all passed!`; smoke → `rev=1 state=sha256:bbbb log=sha256:3333 proposal=true` |
| V13 | Golden regeneration, with a known-positive control **first** | `go run /tmp/regen_golden.go packages/world-core` (live) → `cmp` vs committed golden; then against the `gradeOf`-bearing copy | **CONTROL PASS: byte-identical to the committed golden.** New: `contentHash sha256:489d5e5d47d5c3443ca698f69c6250db57831f463cad8abbec8510ae55ecf632`, `interfaceHash sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083` (**UNCHANGED**), `tarballSHA256 sha256:d0cdf42be80eee2de043645a6bcdb491cb439e3ba659b08b3ff67ee797e99001`, `tarballBytes 6236` (was 5773). `publish --dry-run` independently reported the same tarball/content/interface prefixes |
| V14 | `:376` is a literal, and nothing asserts its numbers | `grep -n 'verify gate PASSED\|EXACT_TOTAL' scripts/verify_ail.sh`; `grep -rn 'verify gate PASSED' --include='*.go'` | `:315` and `:370` interpolate; `:376` is a hardcoded literal. Only Go assertion is `ail_binary_gate_test.go:118 passedMarker = "verify gate PASSED"`, used via `strings.Contains` — a substring |
| V15 | Doc freshness (controller C5) | `git diff --name-only 6d12a79..HEAD` | exactly the four `host/broker` files outside `design_docs/`; control: 10 files total |
| V17 | AC7's grep instrument, with a known-positive control and a fired-on-a-mutant proof | `grep -cE '=> *PROVEN\|=> *ATTESTED' <post>`; then the same on an M7 mutant; `sed -n '23,28p' \| shasum -a 256` on base and post | post: `=> *PROVEN` **0**, control `=> *ATTESTED` **6**; bare `PROVEN` lines **3** (declaration, comment, and `PROVEN => 4` where PROVEN is a *pattern*), bare `ATTESTED` lines **8**; `ProofReport\|ReplayReport\|UNSUPPORTED` **0** (base: every count 0). On the M7 mutant `=> *PROVEN` → **2**, so the instrument fires. `Evidence` lines 23–28 hash `cdf389ae6302b183…` **identically before and after** |
| V16 | Blast-radius sweep: nothing else pins these totals | `grep -rn '4/4 required\|11 module(s)\|14 named tests\|4 required identities\|EXACT_TOTAL' --include='*.go' --include='*.sh' --include='*.yml' --include='*.json' .` | exactly `module_manifest_gate_test.go:128` and `verify_ail.sh:{310,311,312,315,340,366,368,370,376}` — **no other file**. `InterfaceHash` at `host/pkgproj/pkgproj.go:86` reads only manifest fields. `host/pkgproj/readypacket_test.go` **reads** the golden (no hardcoded hashes). `host/boundary`'s `ailArmMutate` appends a comment and is content-agnostic; `wantFileCount=1` counts `host/boundary` `.go` files and is untouched |

---

## 9. Handoff

**SPRINT_PLAN_PATH**: `design_docs/planned/w-evidence-grade-mapping-sprint-plan.md`
**SPRINT_JSON_PATH**: `.ailang/state/sprints/w-evidence-grade-mapping.plan.json`

Open obligations carried forward, unchanged and **out of scope here**: mint authority for `PROVEN`,
the validation boundary, the first production `Evidence` codec, and producer wiring — all owned by
queue item 17, `w-validated-proven-evidence-boundary`, to be priced only after a fresh inventory.

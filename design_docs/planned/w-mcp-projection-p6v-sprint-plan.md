# Sprint plan — `w-mcp-projection` milestone **P6.V** (VERIFIED commit-boundary law)

- **Design doc**: [`design_docs/planned/w-mcp-projection.md`](w-mcp-projection.md) (1361 lines);
  milestone section `### P6.V — VERIFIED commit-boundary law (pure-core world/*.ail) (~0.3d)` at
  line 530; `Decision 6` at line 374; `Files to Create/Modify` at line 576; `AC10`/`AC17` at
  lines 692/716; `Non-Vacuity` table at line 725.
- **Sprint id**: `w-mcp-projection-p6v` · **Plan JSON**: `.ailang/state/sprints/w-mcp-projection-p6v.plan.json`
- **Planner**: sprint-planner (opus sub-agent), mission-world iteration 127, 2026-08-26.
- **Base**: `592a221` on `dev` (`dev == origin/dev`, tree clean). Every measurement in this plan
  was taken **in the main checkout at `592a221`** on darwin/arm64 and the tree was left clean.
- **Risk**: LOW-MEDIUM. The language risk is *retired* (see §2); the residual risk is the
  five-file coupling the design doc does not name (§3).
- **Estimate**: ~0.3 World days (~2.5h), ~65 LOC net. ONE milestone ⇒ **ONE commit**.

---

## 1. Scope

**IN**
- `world/contracts.ail` — the commit-boundary law (pure core, `! {}`).
- `packages/world-core/world/contracts.ail` — the byte-identical reprojection (§3).
- `scripts/verify_ail.sh` — `REQUIRED_VERIFIED` **and** `EXACT_TOTAL_VERIFIED` (§3, doc defect D2).
- `scripts/world_package_ready_packet.golden.json` — regenerated (§3).
- `docs/SELF_MOD_PUBLISH.md` — the two moved digests (§3).

**OUT — do not re-plan, do not implement**
- `P6.T` (toolchain floor) — **LANDED** at `8b196c3` via PR #95. AC15 belongs to it.
- `P6.D` (dependency admission) — **DEFERRED out of this doc entirely** at the round-5 carve-out
  (2026-08-26). No `require github.com/sunholo-data/ailang v0.33.2`, no allowlist edit, no
  `host/daemon/protocol_use.go`. AC16 travelled with it.
- `P6.B-A2A` — **MOVED** to `w-a2a-session-projection.md`.
- **No Go *source* changes.** (But `go test ./...` is still in the blast radius — see §3.)
- No new `world/` module file. No new route, handler, effect, or dependency.

**ACs this milestone closes**: **AC10** (honest gates: the measured floor plus EXACTLY the P6.V
additions) and **AC17** (a named commit-boundary identity Z3-verified on the pinned v0.30.0
binary and pinned in `REQUIRED_VERIFIED`, with the JSON checks `verify.errors == 0`,
`counterexample == 0`, identity present and `verified`).

---

## 2. AILANG fluency brief — read this before writing a line of `.ail`

The charter's fluency protocol (S5) is binding. **Two of its three normal sources are
unavailable or stale, and this is stated here so the executor does not discover it mid-sprint:**

1. **The `ailang-docs` MCP declared in `.mcp.json` was NOT exposed in the planning session's tool
   surface.** Probe it once; if it is absent for the executor too, do not stall on it.
2. **Do NOT use `ailang prompt` as the reference.** The pinned v0.30.0 binary's prompt corpus tops
   out at **v0.12.1** — eighteen minor versions behind the language it compiles. It is a
   stale instrument; a syntax claim taken from it is not evidence.
3. **The authoritative fluency source for this sprint is the repo's own gate-verified corpus at
   exactly the pinned version**, plus `design_docs/coding-standards.md`. Read, in order:
   - `world/contracts.ail:107-120` — `isValidNextWorld`. **This is the template.** Z3-proven
     pattern: `-> bool ! {}`, an `ensures { result == (...) }` that restates the body *exactly*,
     and a body inlined to field/scalar comparisons (never calling an unencodable-bodied callee).
   - `world/types.ail:188-217` — `validDefer`, `wellFormedSchedule`, `timeoutFiredLegally`,
     `validEscalation`. **All-scalar laws — the exact shape P6.V needs.**
   - `world/transitions.ail:57-59` — `applyRevision`, the `requires { ... }` precondition form.
   - `world/contracts.ail:10-15` — the recorded reason Contracts 1–3 are TESTS-ONLY: a contract
     that reaches `Proposal.evidence` (an ADT) Z3-errors `unknown sort 'Proposal'` **while
     `ai-check` exits 0**. This is the trap `MUT-SORT-SILENT` guards.
   - `design_docs/coding-standards.md` §S1, §S2, §S3, §S6.
4. **Verify every syntax claim against the pinned binary before it enters the diff** (S5.3).
   `/tmp/ailang-v0300/ailang --version` → `AILANG v0.30.0`, commit `e37b370`. That binary is the
   ONLY one this milestone may be validated against.

### Syntax facts measured first-party by the planner on the pinned binary (2026-08-26)

| Claim | Evidence |
|---|---|
| Unary `!` does **not** exist | `ai-check` → `type error … unknown unary operator: !`, `check.passed=false` |
| `not(x)` is the negation form | same file with `not(...)` → `check.passed=true`, `verified: 1` |
| `requires` + `ensures` on an all-scalar law verifies | probe below → `verified: 1, counterexample: 0, errors: 0` |
| A parameter unused in the body is legal and does not block encoding | probe carries `clientDisconnected` unused → still `verified: 1` |
| `ai-check` rc on a counterexample = **1**; rc on a Z3 sort error = **0** (silent) | measured both arms |
| Adding a Z3-proven contract to `world/contracts.ail` moves `len(tests[])` by **0** (40 → 40) and `passed_tests` 44 → 45 | `ailang test --format json world/` on a scratch copy, before/after |

---

## 3. THE FINDING THIS PLAN EXISTS TO CARRY — P6.V touches **five** files, not two

The design doc's `Files to Create/Modify` table names exactly two rows for P6.V: `world/*.ail`
and `scripts/verify_ail.sh`. **That is incomplete, and shipping only those two REDs the gate.**
Measured first-party at `592a221`:

`scripts/verify_ail.sh` Leg 3 calls `scripts/verify_world_package.sh`, whose nine steps bind
`world/*.ail` to a projected, hashed, published package:

- **Step 3/9 — projection SHA-256 equality.** `world/<m>.ail` and `packages/world-core/world/<m>.ail`
  must be byte-identical. Confirmed SAME for all four modules at base.
- **Step 9/9 — canonical ready-packet golden.** The recomputed packet is `cmp -s`'d against
  `scripts/world_package_ready_packet.golden.json`. `contentHash`, `tarballSHA256` and
  `tarballBytes` are functions of the `.ail` bytes and **all three move** when the law lands.
- **`interfaceHash` does NOT move.** `host/pkgproj/pkgproj.go:87 InterfaceHash` hashes only the
  manifest (name, edition, ailang, exported *module names*, effects) and never opens a `.ail`
  file. Adding a function inside an existing module leaves it at
  `sha256:d16cc882…`. So exactly **two** of three digests move.
- **`docs/SELF_MOD_PUBLISH.md:85-87`** pins those three digests, and `host/runbook`'s AC28
  (`TestRunbookDigestsAppearVerbatimInTheCommittedGolden`) asserts document-against-artifact.
  **A stale digest table REDs `go test ./host/runbook/`** — i.e. `verify_go.sh` REDs even though
  P6.V changes no Go source. The directive's "no Go changes" is true of the *source* and false of
  the *gate*.

**Precedent (executable template): `cbd17de` "PE.A"** moved exactly this file set —
`world/types.ail`, `packages/world-core/world/types.ail`, `scripts/verify_ail.sh`,
`scripts/world_package_ready_packet.golden.json` (+ `docs/SELF_MOD_PUBLISH.md`) — in ONE commit,
and its own message records the two gates that caught the omission (its M18/M19 mutations).

**Corollary — do NOT create a new `world/` module.** A new file cascades into
`verify_ail.sh:LEG1_MODULES`, `verify_world_package.sh:MODULES`+`EXPORTS`, the *hardcoded frozen
manifest* in step 4, the hardcoded tar-entry allowlist in step 8, the hardcoded `exports` list in
step 9, and `packages/world-core/ailang.toml` — six additional hardcoded sites, several inside
the gate script itself. **The law lands in the existing `world/contracts.ail`.** It is also
where it belongs semantically: that module already holds `commitAllowed` and `isValidNextWorld`.
This is the S3 "why is this not a package?" answer: **it is not a new kernel surface — it is one
proof about commit behaviour the kernel already has, added to the module that already states
commit policy. Zero new modules, zero new exports, zero new effects, `interfaceHash` unchanged.**

---

## 4. Branch A vs Branch B — and which one the planner measured

The design doc is explicit that P6.V is **conditional, not a mandate**: if no encodable statement
of the law exists, the identity floor correctly stays at 10 and the milestone closes on a named
**test-only** law plus a limitation row (the `w-m1-ailang-hardening` pattern).

**The planner probed the empirical question and Branch A is REACHABLE.** On the pinned v0.30.0
binary, in `/tmp/p6vfinal` (a scratch module outside the repo), this law verified:

```ailang
export func commitBoundaryHolds(
  intentJournaled: bool,
  accepted: bool,
  clientDisconnected: bool,
  receiptCount: int,
  outcomeKnown: bool,
  reportedNotCommitted: bool
) -> bool ! {}
requires { receiptCount >= 0 }
ensures { result == (
     (receiptCount == 0 || intentJournaled)
  && (receiptCount <= 1)
  && (accepted == (receiptCount == 1))
  && (not(reportedNotCommitted) || outcomeKnown)
) }
{
     (receiptCount == 0 || intentJournaled)
  && (receiptCount <= 1)
  && (accepted == (receiptCount == 1))
  && (not(reportedNotCommitted) || outcomeKnown)
}
```

→ `check.passed: true`, `verify: {verified: 1, counterexample: 0, errors: 0}`,
`results: [("commitBoundaryHolds", "verified")]`.

**What that does and does not establish.** It establishes that *an* encodable, Z3-verified
statement of the reviewer's contract exists — Branch B is not forced by the language. It does
**not** establish that this exact text is the right law, that it survives the executor's own
re-reading of Decision 6, or that it verifies unchanged once pasted into `world/contracts.ail`
next to ADT-importing neighbours. The executor re-derives all of that first-party. If the
executor's chosen statement fails to encode, **Branch B is a legitimate close, not a failure** —
take it, and record the limitation.

### Clause-by-clause map to Decision 6 (`w-mcp-projection.md:387-392`, the reviewer's verbatim contract)

| Reviewer's clause | Law conjunct | Go anchor (re-derived first-party, all five CORRECT at `592a221`) |
|---|---|---|
| "Before the coordinator accepts a commit, cancellation guarantees no durable mutation" | `accepted == (receiptCount == 1)`, ⇐ direction | `bindCommitIntentTx` `host/store/store.go:1025` |
| "Commit uses a stable invocation/idempotency ID" | every parameter is scoped to ONE invocation ID; `receiptCount <= 1` is the idempotency claim | `JournalIntent.InvocationID` `host/store/journal.go:28-29` |
| "…and has a defined point of no return" | `accepted` **is** the point of no return | `bindCommitIntentTx` |
| "Once accepted, the critical section may complete despite client disconnect" | `accepted == (receiptCount == 1)` holds for **both** values of `clientDisconnected` (Z3 quantifies over it) | `recoverCommitPending` `host/broker/recover.go:126` |
| "its durable receipt is recorded and queryable/replayable" | `accepted == (receiptCount == 1)`, ⇒ direction; `receiptCount == 0 \|\| intentJournaled` | `Store.GetReceipt` `journal.go:813`, `GetEffectReceipt` `journal.go:852` |
| "never a definitive 'not committed' when the outcome is unknown" | `not(reportedNotCommitted) \|\| outcomeKnown` | — |

**P6.V changes no Go.** The Go surface is landed; what P6.V adds is the *proof*.

### Merge criteria — distinct per branch

| | **Branch A (Z3 path)** | **Branch B (encodable fallback fired)** |
|---|---|---|
| `REQUIRED_VERIFIED` | gains `"commitBoundaryHolds"` under `world/contracts.ail` | **unchanged** |
| `EXACT_TOTAL_VERIFIED` | `10 → 11` | **stays 10** |
| `REQUIRED_TESTS` / `EXACT_TOTAL_TESTS` | **unchanged (40)** — measured: a proven contract adds 0 to `len(tests[])` | gains the law's named test ids; `EXACT_TOTAL_TESTS` bumped to match |
| `verify_ail.sh` banner | `11 required identities verified, 40 named tests pass` | `10 required identities verified, N named tests pass` |
| AC17 | closed by the verified identity | closed by the named test-only law **+ an explicit limitation row written into `w-mcp-projection.md`** |
| Mutations | `MUT-LAW-BREAK`, `MUT-SORT-SILENT` both RED via the JSON checks | `MUT-LAW-BREAK` REDs via the failing named test; `MUT-SORT-SILENT` is **N/A and must be recorded as N/A, not as passed** |

**Do not write an unconditional "raise the floor to 11" into the implementation report.** Report
which branch fired and why, with the command output that decided it.

---

## 5. Milestone breakdown — one milestone, three ordered phases, one commit

### P6V_A — author the law and prove it (≈1.0h, ~30 LOC)
1. Read the fluency corpus in §2. Re-read `Decision 6` at `w-mcp-projection.md:374-405`.
2. Append the law to `world/contracts.ail` as **Contract 5**, with a header comment in the house
   style of Contracts 1–4: state (a) that it is Z3-PROVEN, (b) that the parameters are
   deliberately all-scalar because an ADT-bearing parameter Z3-errors *silently*, (c) the
   clause-by-clause map to Decision 6.
3. Prove it standalone before touching any gate:
   `AILANG_BIN=/tmp/ailang-v0300/ailang /tmp/ailang-v0300/ailang ai-check -timeout 5s world/contracts.ail`
   → require `check.passed: true`, `verify.errors == 0`, `verify.counterexample == 0`, and
   `commitBoundaryHolds` present with `status: "verified"`. **Record the JSON verbatim.**
4. If it does not encode: **stop, take Branch B**, and record the exact `reason` string from
   `verify.results[].reason` — that string is the evidence, not a paraphrase.

### P6V_B — couple the five files and green the gates (≈1.0h, ~35 LOC)
1. `cp world/contracts.ail packages/world-core/world/contracts.ail` (byte-identical reprojection).
2. `scripts/verify_ail.sh:274-279` — add `"commitBoundaryHolds"` to the `world/contracts.ail` set.
3. `scripts/verify_ail.sh:323` — `EXACT_TOTAL_VERIFIED=10` → `11`. **The design doc does not name
   this line; without it the gate REDs with `expected exactly 10 proven world/ contracts, got 11`.**
4. Regenerate `scripts/world_package_ready_packet.golden.json` with the recipe in §6 (Gate G5).
5. Update the two moved digests in `docs/SELF_MOD_PUBLISH.md:85-87` **by copying them out of the
   regenerated golden, never by retyping**. `interfaceHash` must be left alone — if it changed,
   something is wrong and you must stop.
6. Run the full gate set (§6). All green.

### P6V_C — mutation drill (≈0.5h, **zero net diff**)
Run §7. Every mutant is `cp`-restored and the restore is **byte-verified with `shasum -a 256`
before and after**. A `SURVIVED` result is a finding to record, never a result to hide.

---

## 6. Acceptance gates — exact runnable commands, with baselines

Every command is run from the repo root. **`AILANG_BIN=/tmp/ailang-v0300/ailang` is MANDATORY on
every line that invokes either verify script.** A bare `./scripts/verify_go.sh` with `AILANG_BIN`
unset is **rc=1 at base** (`✗ AILANG_BIN is unset — host/replay tests would t.Skip() silently and
this gate would be false-green`) — controller-baselined; a gate line without the prefix is
unsatisfiable, not strict.

| id | command | base rc | baselined by | expected after P6.V |
|---|---|---|---|---|
| **G1** | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | **0** | **controller**, this session at `592a221`. Final lines: `✓ world package gate PASSED: 9/9 steps performed non-zero work` and `✓ verify gate PASSED: 10 required identities verified, 40 named tests pass` | rc=0, banner reads **`11 required identities verified, 40 named tests pass`** (Branch A) |
| **G2** | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` | **0** | **controller**, this session at `592a221` (build clean, plain + `-race`, driver-drift gate green, 11 host packages ok) | rc=0, unchanged |
| **G3** | `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/runbook/... ./host/pkgproj/...` | **0** | **planner**, main checkout at `592a221`: `ok host/runbook 1.028s`, `ok host/pkgproj 0.695s` | rc=0 — this is the **narrow, fast** detector for the `docs/SELF_MOD_PUBLISH.md` digest coupling; run it before G2 |
| **G4** | `WORLD_PKG_AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_world_package.sh` | **0** | **planner**, main checkout at `592a221`: `✓ canonical JSON equals committed golden byte-for-byte`, `✓ world package gate PASSED: 9/9 steps performed non-zero work` | rc=0 — the **narrow** detector for the projection+golden coupling; run it before G1 |
| **G5** | golden regeneration (recipe below) | reproduces the committed golden **byte-for-byte** | **planner**, main checkout at `592a221` — `cmp` exit 0 against `scripts/world_package_ready_packet.golden.json` | emits the new golden |
| **G6** | `git diff --name-only` | — | — | **exactly five paths**: `world/contracts.ail`, `packages/world-core/world/contracts.ail`, `scripts/verify_ail.sh`, `scripts/world_package_ready_packet.golden.json`, `docs/SELF_MOD_PUBLISH.md` |
| **G7** | `shasum -a 256 world/contracts.ail packages/world-core/world/contracts.ail` | equal | planner (all four modules SAME at base) | the two digests must be **equal** |
| **G8** | CI on the PR: both `ailang-code verify gate` and `go host build + test gate` `success` | green at `592a221` (controller: SHA-addressed `check-runs`, `checks=2`, both `success`) | controller | both `success`, zero skips |

### G5 — golden regeneration recipe, **validated first-party by the planner**

The recipe was run at base and **reproduced the committed golden byte-for-byte** (`cmp` exit 0) —
that is its known-positive control. It is a verbatim extract of `verify_world_package.sh` steps
7 and 9, so it cannot drift from the gate it must satisfy.

```bash
mkdir -p /tmp/p6v-regen && cat > /tmp/p6v-regen/main.go <<'GO'
package main
import (
  "fmt"
  "os"
  "github.com/sunholo-data/ailang-world/host/pkgproj"
)
func main() {
  m := pkgproj.Manifest{Package: pkgproj.Package{Name:"world/core", Edition:"1", AILANG:">=0.30.0"}, Exports:pkgproj.Exports{Modules:[]string{"world/types","world/contracts","world/transitions","world/logepoch"}}, Effects:pkgproj.Effects{Max:[]string{}}}
  r, err := pkgproj.CrossCheck("packages/world-core", m, os.Args[1]); if err != nil { panic(err) }
  fmt.Printf("contentHash=%s\ninterfaceHash=%s\ntarballSHA256=%s\ntarballBytes=%d\n", r.Local.Content, r.Local.Interface, r.Local.Tarball, r.Local.TarballBytes)
}
GO
env -u AILANG_REGISTRY_API_KEY go run /tmp/p6v-regen/main.go /tmp/ailang-v0300/ailang > /tmp/p6v-regen/proj.txt
python3 - /tmp/p6v-regen/proj.txt scripts/world_package_ready_packet.golden.json \
  "$(/tmp/ailang-v0300/ailang --version | sed -n '1p')" <<'PY'
import json, sys
values = {}
for line in open(sys.argv[1], encoding="utf-8"):
    if "=" in line:
        k, v = line.rstrip("\n").split("=", 1); values[k] = v
packet = {"compilerVersion":sys.argv[3], "contentHash":values["contentHash"], "effects":[], "exports":["world/types","world/contracts","world/transitions","world/logepoch"], "interfaceHash":values["interfaceHash"], "package":"world/core", "tarballBytes":int(values["tarballBytes"]), "tarballSHA256":values["tarballSHA256"], "version":"0.1.0"}
with open(sys.argv[2], "wb") as f: f.write((json.dumps(packet, sort_keys=True, separators=(",", ":")) + "\n").encode())
PY
```

**Run it only AFTER step P6V_B.1** (the reprojection). Running it against a stale
`packages/world-core/` mints a golden for bytes that are not the source, and step 3/9 will still
catch it — but you will have written a wrong artifact.

Baseline values at `592a221` (what the recipe must reproduce before you change anything):
`contentHash=sha256:7616498e…4459d`, `interfaceHash=sha256:d16cc882…b4083`,
`tarballSHA256=sha256:2f18c5e8…b3a21`, `tarballBytes=7883`.

### Sandbox rule
The executor lane is `codex:gpt-5.6-sol` under `--sandbox workspace-write`, which denies loopback
socket binds. This surface is `.ail` + shell + `go test` on two pure packages, so socket denial is
unlikely to bite — but **any gate result the executor cannot obtain must be reported as
`UNINFORMATIVE UNDER SANDBOX`, never as pass and never as fail.** The controller re-runs G1–G4
outside the sandbox before committing. The executor performs **no git write operations**.

---

## 7. Mutation drill — the test plan

Per rule 3i, each row names its **observable** and states whether that observable is *downstream
of the mechanism* or merely adjacent to it. Restore protocol: `cp` backup before, `cp` restore
after, `shasum -a 256` equal before and after — recorded.

| id | status | mutation (exact edit) | required RED observation | observable — downstream? |
|---|---|---|---|---|
| **MUT-LAW-BREAK** (doc row, **MUST**) | **planner-validated** | In `world/contracts.ail`, delete the conjunct `&& (receiptCount == 0 \|\| intentJournaled)` **from the BODY ONLY**, leaving `ensures` intact | `ai-check` → `verify.counterexample == 1`, `commitBoundaryHolds` status `counterexample`; `verify_ail.sh` REDs. Planner-measured model: `$p_intentJournaled=false, $p_receiptCount=1` — **a receipt with no journal intent**, the exact semantics named in the doc. Raw `ai-check` rc=1 | **Downstream.** The counterexample is produced by Z3 from the mutated body against the unmutated postcondition. No other mechanism in the repo writes `verify.counterexample`. Mutating body *and* `ensures` together would re-verify — that is the trap; mutate the body only |
| **MUT-LAW-BREAK / arm 2** (doc row's second alternative, **MUST**) | **planner-validated** | Body only: `(receiptCount <= 1)` → `(receiptCount <= 2)` **and** `(accepted == (receiptCount == 1))` → `(accepted == (receiptCount >= 1))` — i.e. permit a second receipt under one invocation ID | `verify.counterexample == 1`. Planner-measured model: `$p_receiptCount=2, $p_accepted=true` — **two receipts under one ID** | Downstream, same mechanism |
| **MUT-SORT-SILENT** (doc row, **MUST**) | **planner-validated** | Retype one law parameter to an ADT-bearing `Proposal`-class sort (e.g. `p: Proposal`, already imported at `world/contracts.ail:4`), referencing it in body and `ensures` | `verify.errors == 1` in the JSON, `status: "error"`, `reason` containing `unknown sort`; **`ai-check`'s own exit code is 0** (planner-measured on a `Prop`-shaped ADT: `raw_rc=0`, `errors: 1`, `unknown sort 'Prop'`). `verify_ail.sh` REDs on `verify.errors == %s (Z3 encoding error, silent under exit codes V10)` | **Downstream, and this row is the whole point.** The exit code is 0, so *only* the gate's JSON `verify.errors` branch can see it. This mutation proves that branch is armed rather than assumed |
| **MUT-DISCONNECT-EXCUSE** (planner-added, SHOULD) | **planner-validated** | Body only: `(accepted == (receiptCount == 1))` → `((accepted && not(clientDisconnected)) == (receiptCount == 1))` | `verify.counterexample == 1`. Planner-measured model: `$p_accepted=true, $p_clientDisconnected=true, $p_receiptCount=1` | Downstream. This is the *only* row that discriminates the "critical section may complete despite client disconnect" clause. Without it, `clientDisconnected` is a decorative parameter |
| **MUT-STALE-PROJECTION** (planner-added, SHOULD — the §3 gap) | precedent `cbd17de` M18 | Revert `packages/world-core/world/contracts.ail` to its `592a221` bytes, leaving `world/contracts.ail` mutated | `verify_world_package.sh` step **3/9** REDs: `✗ projection mismatch: world/contracts.ail=<a> packages/world-core/world/contracts.ail=<b>` | Downstream. Step 3 hashes both files itself; nothing else writes that message |
| **MUT-STALE-GOLDEN** (planner-added, SHOULD — the §3 gap) | precedent `cbd17de` M19 | Revert `scripts/world_package_ready_packet.golden.json` to its `592a221` bytes with everything else correct | Step **9/9** REDs: `✗ ready packet differs byte-for-byte from golden` + a diff | Downstream. Step 9 recomputes the packet from disk and `cmp`s |
| **MUT-STALE-RUNBOOK** (planner-added, SHOULD — the §3 gap) | precedent `cbd17de` defect 2 | Revert `docs/SELF_MOD_PUBLISH.md:85-87` to its `592a221` digests with everything else correct | `go test ./host/runbook/` REDs in `TestRunbookDigestsAppearVerbatimInTheCommittedGolden` (AC28): *"names digest … which does not appear in …. The runbook is telling an attended operator to approve bytes that are not the reviewed artifact"* | Downstream. The test reads the digests **out of** the document and membership-tests them against the golden; it never recomputes them |
| **MUT-FLOOR-BLIND** (planner-added, SHOULD — AC10's own teeth) | — | Set `EXACT_TOTAL_VERIFIED` back to `10` with the law landed and `REQUIRED_VERIFIED` updated | `verify_ail.sh` REDs: `✗ expected exactly 10 proven world/ contracts, got 11` | Downstream, and it is the discriminating check for AC10's "nothing else moves": `REQUIRED_VERIFIED` alone would pass with an *extra* unnamed contract present. **This row is why §3's doc defect D2 matters** |

**Rows deliberately NOT written**: an assertion that `len(tests[]) == 40` is unchanged. It is
true (planner-measured), but `EXACT_TOTAL_TESTS` is written by the *test* leg, which this
milestone does not touch on Branch A — asserting it would be decorative, and it is already
enforced by G1. It is stated here as a *prediction the sprint will falsify or confirm*, not as a
mutation row.

---

## 8. What this plan CANNOT establish before the sprint runs

Stated explicitly, because fabricating confidence here is the failure mode this loop keeps
catching:

1. **Which branch fires is an empirical question the executor answers.** The planner showed
   Branch A is *reachable* with one specific law text in a scratch module. It did **not** prove
   the executor's final text encodes, nor that it encodes *inside* `world/contracts.ail` beside
   ADT-importing neighbours. Branch B remains a legitimate close.
2. **Whether the law as written is the *right* rendering of Decision 6.** That is a semantic
   judgment for the executor and then the evaluator. The clause map in §4 is the planner's
   reading, not a proof.
3. **Whether `contentHash`/`tarballSHA256`/`tarballBytes` land on any particular value.** They are
   functions of bytes that do not exist yet. Only the *recipe* is validated, not the outputs.
4. **CI behaviour.** Every green in this plan is darwin/arm64. The two CI legs (ubuntu-latest,
   linux/amd64) are unrun locally. Note that `verify_world_package.sh` pins the compiler by
   **exact bytes per platform** — a linux/amd64 SHA table entry exists, so this is expected to
   hold, but it is asserted, not measured here.
5. **Whether the `ailang-docs` MCP is reachable from the executor's surface.** It was not
   reachable from the planner's.
6. **Whether any *other* consumer of `world/contracts.ail`'s exports exists that the five-file
   set misses.** The planner swept the gates, not every Go test; G2 (full `verify_go.sh`) is the
   backstop and it must be run outside the sandbox before the commit.

---

## 9. Design-doc defects found by this plan

| id | severity | where | defect |
|---|---|---|---|
| **D1** | **MATERIAL** | `Files to Create/Modify`, `w-mcp-projection.md:583-584` | The P6.V rows name only `world/*.ail` and `scripts/verify_ail.sh`. Three more files are load-bearing (`packages/world-core/world/contracts.ail`, `scripts/world_package_ready_packet.golden.json`, `docs/SELF_MOD_PUBLISH.md`) and omitting any one REDs a gate. Precedent `cbd17de` moved exactly this set. Non-blocking: this plan carries the full set |
| **D2** | **MATERIAL** | `w-mcp-projection.md:534` and AC10/AC17 | Both name `REQUIRED_VERIFIED` (`verify_ail.sh:274-279`, anchor **correct**) as the only manifest change. On Branch A the **secondary exact total `EXACT_TOTAL_VERIFIED=10` at `verify_ail.sh:323`** must also move, or the gate REDs `expected exactly 10 proven world/ contracts, got 11`. Non-blocking: carried as P6V_B.3 and guarded by `MUT-FLOOR-BLIND` |
| **D3** | minor | `w-mcp-projection.md:536` | "Load the AILANG language reference via the `ailang-docs` MCP BEFORE writing any `.ail` (charter fluency protocol; binding)" — the MCP was not exposed in the planning session's tool surface, and `ailang prompt` on the pinned binary is stale at v0.12.1. The doc names a source that may not exist for the executor. Non-blocking: §2 supplies a repo-local substitute |
| **D4** | minor | AC17, `w-mcp-projection.md:716-721` | "the named test-only law is in the 40+ named tests" is written as if Branch B *adds* tests to a floor of 40 — true — but AC10's companion clause pins "40 named tests" as the floor that must not move. On Branch B they conflict textually. Non-blocking: §4's branch table resolves it (Branch B bumps `EXACT_TOTAL_TESTS`; Branch A does not) |
| **D5** | informational | `Axiom Compliance` table, `w-mcp-projection.md:747-760` | Rows `A4 +1` and `A10 +1` are justified entirely by "the allowlist admits ONE package path" / "upstream-owned wire contract admitted at a pinned version" — **that is P6.D, deferred out of this doc at the round-5 carve-out.** The net `+10` is stale by at least those two points. Non-blocking for P6.V; recorded for the doc's next revision |
| **D6** | informational | `Estimate honesty`, `w-mcp-projection.md:568-570` | Still reads "split #2 re-cuts it to **~0.55d in this doc** (P6.T ~0.1 + P6.D ~0.15 + P6.V ~0.3)". P6.D was deferred after that sentence was written; the remaining figure is ~0.4d, which the head text at line 12 already states correctly. Internal inconsistency, non-blocking |

**No doc↔plan divergence to reconcile**: there was no prior sprint plan for P6.V.
`.ailang/state/sprints/w-mcp-projection-p6t.plan.json` covers P6.T only.

---

## 10. Definition of done

1. G1–G4 all rc=0, **run outside the sandbox by the controller**, with G1's banner recorded verbatim.
2. `git diff --name-only` is exactly the five paths in G6 — no sixth file, no `tools/launchd/*`.
3. G7: source and projection digests equal.
4. `docs/SELF_MOD_PUBLISH.md`'s `interfaceHash` is **unchanged**; the other two match the
   regenerated golden character-for-character (copied, never retyped).
5. Every mutation row in §7 has a recorded result: RED with its literal message, or `SURVIVED`,
   or `N/A` with the reason. No predicted-red sets. Tree byte-identical after each restore.
6. The implementation report states **which branch fired**, the JSON that decided it, and — on
   Branch B — the limitation row text added to `w-mcp-projection.md`.
7. ONE commit on a `sprint/w-mcp-projection-p6v` branch, authored by the controller.
8. Both CI jobs `success` on the PR, zero skips.

---

## 11. Invariants (charter, non-negotiable)

- **Frozen core**: never edit `tools/launchd/*`; never copy skills into this repo. A
  `verify_go.sh: FATAL: DRIVER DRIFT` red means **the fleet must commit** — stop and report; it
  never licenses absorbing the driver into this change.
- **Pinned binary only**: `/tmp/ailang-v0300/ailang` (`AILANG v0.30.0`, commit `e37b370`). Never a
  `-dirty` dev build, never a syntax claim from memory.
- **Language gaps route upstream** as `sunholo-data/ailang` issues plus an `ailang messages send
  mission-control` note — never a local workaround, never a vendored fork.
- **Never touch** `~/.ailang/state/mission-v1*` or the V1 checkout.
- **S2**: `world/` is pure. The law carries `! {}` and describes state; it performs nothing.

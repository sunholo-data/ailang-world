# Sprint Plan: W-DAEMON-STARTUP-UNDER-LOAD

Authored 2026-11-07 by this `codex:gpt-6.1-sol` planner invocation (derive/resolve: both `codex:gpt-6.1-sol`, `declared:planner-lane-default-pin`). No delegated author, fallback or sub-agents. World iteration 362, row 256.

## Summary and authority

Apply the approved r3 frozen patch so registry bootstrap, template reporting and integrity scanning receive a fresh 9 s allowance from serve's original parent. Failed exec profiles remain excluded; doctor, health and logs expose recording-time observations with retained reasons and review advice.

**Duration:** one working day, approximately six hours including controller/evaluation work, with CI queue time variable. **Risk:** medium (startup lifetime and fail-closed boundary). **Dependencies:** no new package dependencies and no milestone dependencies. Outside-sandbox controller baselines and acceptance, independent evaluation and CI remain prerequisites for completion.

Source: [approved design](w-daemon-startup-under-load.md), especially §§3–8, V27–V29 and §11. Design commit: `caee00cae7dd6938c9017f1977209f4b61ef7c92`; implementation base: `61a8d1c47b45fffc0f36137d29a5285c32da0145`.

- Design SHA256: `94f8bf9d7d1f3ac51a848c89e4a1595f892c2991984dc11778ffa0f539d6430a`.
- Frozen patch SHA256: `90c6df45c6b4dcf09fd241d76eb87294f5b84d688eb9f2ff90d0eca1de09db16`.
- Frozen patch: **11 files, +881/-20, net +861, 1088 patch lines; four new test files**. All eleven post-apply hashes are pinned in the machine state and the existing prototype manifest.

D-WORLD-108's attended ruling authorizes row 256 after row 255. Row 233 stays parked on D-WORLD-105 (OPEN, default C); row 257 is separate. D-WORLD-110 is row 255's policy question and does not gate this patch. No retries, cache-seed changes, reverse/forward verdict changes, archive redesign or live binary replacement are authorized.

r1 quorum was **BLOCKED 3/3**; r2 was **BLOCKED 1/3**, all reviewers present. The controller applied concrete remaining fixes under the ratified narrow-refinement carve-out, design §11: V28 empirical evidence, fourth presence guard, parent-context statement and AC9/MUT-SERVE-ON-STARTUP-CLOCK. No third quorum; r2 did not pass. No unresolved human-decision gate exists for this row. Genuine new human judgment parks with terminal STOP; settled helper/evidence-layout details do not require approval.

## Active baseline STOP (receipt arrived during planning)

**Execution is BLOCKED before patch application.** The controller's pristine real-checkout `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` completed **rc 1** at **2026-11-07T05:39:00.911689 UTC**, with `TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter` (`cmd/world-publish`) and `TestCheckSourceUnderPinnedInterpreter` (`host/archive`) red. The retained output shows the pinned interpreter's trace-exporter warning ahead of a diagnostic truncated before the expected `LDR001`; this is the controller's recorded condition, not a new causal diagnosis by the planner. Both affected source/test surfaces are outside the frozen patch. Evidence: [`baseline-inherited-failures.json`](../verification/world-iter362/baseline-inherited-failures.json) and [`baseline-full-verify-go.txt`](../verification/world-iter362/baseline-full-verify-go.txt). Baseline load at end: 8.06/8.58/7.67; resolved Go 1.26.6 darwin/arm64. The baseline AILANG gate completed **rc 0** at 05:35:25.929788 UTC ([receipt](../verification/world-iter362/baseline-verify-ail-real.txt)); this does not clear the Go gate.

**STOP: preserve files, report `BLOCKED: pristine full Go baseline red outside row 256`, and end the implementation run.** Do not apply, launch candidate gates, change `host/archive`, suppress the warning, replace the pin, alter temp/environment policy or treat unrelated red as permission to widen this sprint. Controller owns reporting/parking and any work outside this sprint. The plan, machine state and file handoff are complete planning deliverables; all implementation/CI completion fields remain pending. A later resume requires a genuine rc-0 pristine baseline receipt without masking the red, recorded by the controller; the machine execution hold prevents automatic continuation. The one-day estimate below is conditional on that prerequisite being restored and has no promised start time.

## Current status and measured velocity

Read-only inspection at the design commit found all seven existing targets byte-identical to the implementation base and all four new files absent. The implementation is not applied. All milestone passes, starts and completions are pending. Documentation and evidence already committed with the design are not implemented functionality.

Bounded local history (`git log --max-count=80 --since=2026-10-31 --no-merges -- host cmd docs/QUICKSTART.md`) and per-commit numstat measured the two most recent comparable M1 feature commits:

| Commit / date | Work | Files | Insertions / deletions |
|---|---|---:|---:|
| `d29c4155` / 2026-11-06 | row 250 foreign-dial classification | 4 | +546/-57 |
| `fa5b9183` / 2026-11-06 | row 255 bounded failure reason | 22 | +1101/-37 |

The sample is **1647 inserted LOC on one calendar day** (1553 net), not a measured authoring-time rate or a stable throughput forecast. `docs/HOST_CHANGELOG.md` is a release-level record; `packages/world-core/CHANGELOG.md` last covers 2026-09-30, so neither supplies a recent daily rate. The sprint budgets 881 insertions, about 54% of this small daily sample. Executor estimate is 60 minutes plus a 25% buffer = **75 minutes**. Controller full gates, independent review and CI account for the rest of the day. Any preexisting red invalidates this schedule and stops execution; it does not buy scope expansion.

## Registry reuse audit

**M1: `none`; dependencies `[]`.** This adds no package-like capability. Design §3.6 answers S3: constructor context ownership and startup probe/doctor bookkeeping belong to the existing Go host boundary; an AILANG extension package cannot replace it. Inspected local `world/core` 0.1.1 manifest/CHANGELOG (pure modules, empty effects) and `packages/se-tools/transitions.json` (existing tool transitions). No package can own Go constructor clocks here; no dependency or version changes are planned.

Search concepts: startup, exec probe, doctor, registry phase. No live `ailang pkg search/info/docs` is claimed: this invocation is restricted to existing local evidence and there is no package-like deliverable. `record-registry.txt` was read and reused solely as an **isolated mission-registry/workdir routing record**, not as a package search or evidence about canonical registry contents. No canonical registry, shared state or main checkout is changed.

## Milestone and day-one tasks

### M1: Apply and verify the frozen startup/diagnostic patch (~881 LOC)

One M1 explicitly owns **AC1, AC2, AC3, AC4, AC5, AC6, AC7, AC8, AC9**, and every file below. The approved design §7 and user ruling override the skill's generic two-milestone minimum. Splitting the indivisible frozen patch would create intermediate incomplete construction/diagnostic/test states and break the measured-byte reuse contract. No synthetic M2 is added. The stock validator's minimum-two rule is an acknowledged overridden check; no shared skill/validator is edited.

**Estimate:** 129 implementation insertions + 747 test insertions + 5 documentation insertions = **881** (20 deletions). **Dependencies:** none between milestones; controller baseline receipts must exist before application. Full suites are controller/CI responsibilities, not sandboxed executor milestones.

| Day 1 task | Executor budget | Completion evidence |
|---|---:|---|
| Read r3, plan, file handoff; verify design/patch/current hashes and baseline receipts | 10 min | exact inputs; no preexisting red |
| Apply frozen patch once, four presence guards, eleven post-apply hashes | 15 min | exact application receipt; no manual patch edits |
| Focused build/vet/gofmt; inspect exact named-test/mutation ownership and serial seams | 20 min | focused receipts; controller supplies all socket verdicts outside sandbox |
| Record pending/full obligations and hand back artifacts for controller verification | 15 min | M1 remains pending until acceptance closes |
| Buffer | 15 min | STOP/report on changed premise rather than widening scope |

Controller follows sequentially: honor the active baseline STOP; only after a genuine rc-0 pristine baseline is established, baseline the remaining gates and validate application/focused gates, run candidate full plain Go then AILANG gates, obtain independent evaluation, and require PR plus merge Linux checks. The pristine baseline has now ended red; do not start another baseline merely to seek a green, and do not launch implementation.

## Exact files and size constraints

| File (all owned by M1) | Delta | Final lines | Kind |
|---|---:|---:|---|
| `cmd/ailang-worldd/doctor.go` | +23 / -1 | 538 | implementation |
| `cmd/ailang-worldd/doctor_exec_profiles_test.go` | +100 / -0 | 100 | test |
| `cmd/ailang-worldd/doctor_hosts.go` | +1 / -1 | 264 | implementation |
| `docs/QUICKSTART.md` | +5 / -1 | 1771 | documentation |
| `host/broker/exec_probe.go` | +31 / -4 | 1129 | implementation |
| `host/broker/exec_probe_out_of_time_test.go` | +108 / -0 | 108 | test |
| `host/broker/exec_startup.go` | +8 / -1 | 690 | implementation |
| `host/daemon/daemon.go` | +54 / -10 | 1410 | implementation |
| `host/daemon/startup_callers_test.go` | +225 / -0 | 225 | test |
| `host/daemon/startup_under_load_test.go` | +314 / -0 | 314 | test |
| `host/daemon/workspace.go` | +12 / -2 | 802 | implementation |

The 747 new test lines are 225 caller inventory + 314 daemon + 108 broker + 100 doctor. Existing large host files are intentionally preserved by the frozen patch: `daemon.go` 1410, `exec_probe.go` 1129, `workspace.go` 802. Companion tests are the explicitly named four new files; do not split/refactor production files to satisfy generic AILANG-repo size suggestions. World has no Makefile: generic `make test-changed`, `make test-core`, changelog-fragment and size targets are not executable World gates. Use the native gates below. No `HOST_CHANGELOG.md` edit is in scope.

## Acceptance ownership (all pending)

- [ ] **AC1 (M1) — Fresh bounded registry phase and successful Run:** The registry bootstrap, template report and integrity sweep use a fresh parent-derived 9 s context; startup remains 9 s. Run announces, answers health and shuts down in the slow-client fixture. No whole-startup wall-time bound is claimed. Tests: `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock`, `TestRegistryPhaseClockIsBounded`, `TestServeStartupClocksArePinned`. Mutations: `MUT-RUN-HEAD-SHAPE`, `MUT-SHARED-PARENT`, `MUT-REGISTRY-UNBOUNDED`, `MUT-SWEEP-ON-STARTUP-CLOCK`, `MUT-REGISTRY-BUDGET-30S`, `MUT-STARTUP-BUDGET-5S`.

- [ ] **AC2 (M1) — All failures remain excluded:** All failed profiles, including a genuine refusal recorded late, are excluded with reasons retained and no implicit default; a failed explicit --exec-profile refuses startup. The episode receives the existing named refusal. Tests: `TestExecProbeOutOfTimeNamesTheClockNotAnArm`, `TestExecProbeGenuineRefusalRecordedAfterDeadline`, `TestRegistryPhaseHasItsOwnClock`. Mutations: `MUT-LATE-STILL-SERVED`, `MUT-FLAG-NOT-CARRIED`, `MUT-SERVE-SPENT`, `MUT-SPENT-STILL-SERVED`.

- [ ] **AC3 (M1) — Observation-only operator surfaces:** Log, health comments, doctor, help and QUICKSTART report only recording times, retain genuine refusal reasons and advise review before conditional restart. Doctor names every exclusion as an exec failure (exit 13), including --quick. Tests: `TestExecProbeGenuineRefusalRecordedAfterDeadline`, `TestDoctorNamesLeftOutExecProfiles`, `TestDoctorLateRefusalObservationAndReasonCap`, `TestExecExclusionLogIsObservationOnly`. Mutations: `MUT-CAUSAL-OBSERVATION`, `MUT-LOG-CAUSE`, `MUT-DOCTOR-CAUSE`, `MUT-DOCTOR-NO-REVIEW`, `MUT-DOCTOR-NOT-CALLED`, `MUT-DOCTOR-WARN`, `MUT-DOCTOR-CLASS`, `MUT-DOCTOR-ONE-FIX`.

- [ ] **AC4 (M1) — Nonempty ALL-after timestamp rule:** OutOfTime is true only for a nonempty failed-arm set all recorded at/after the deadline. Mixed, no-failure and no-deadline rows remain false; late genuine refusal may be true. R3/R9 retain both ambiguous directions and establish no scheduling bound. Tests: `TestExecProbeOutOfTimeNamesTheClockNotAnArm`, `TestExecProbeGenuineRefusalRecordedAfterDeadline`. Mutations: `MUT-SPENT-ANY`, `MUT-LATE-NEVER-SPENT`, `MUT-FAIL-NEVER-SPENT`, `MUT-DEADLINE-UNSET`.

- [ ] **AC5 (M1) — Actual caller census and forward guard:** Inventory the actual Go sources, preserving imported/local function-value and direct-call detection plus alias/nested/shadow controls. Expect 60 test New references, zero production New references after apply, and four internal production construction edges. A new compatibility caller requires reviewed design, never an allowlist expansion to hide red. Tests: `TestDaemonConstructorCallerInventory`, `TestDaemonCallerInventoryAliasControl`. Mutations: `MUT-CALLER-IMPORTED-ALIAS`, `MUT-CALLER-DIRECT`, `MUT-CALLER-LOCAL-ALIAS`.

- [ ] **AC6 (M1) — Additive JSON and bounded reason:** False exclusion JSON retains old bytes via omitempty; true adds out_of_time. Doctor reason is one line and at most 240 bytes including the 3-byte ellipsis (237-byte prefix); project, file and fixed prose are outside this cap. Tests: `TestHealthExecExclusionIsAdditive`, `TestDoctorLateRefusalObservationAndReasonCap`. Mutations: `MUT-HEALTH-NO-OMITEMPTY`, `MUT-REASON-CAP-240-PREFIX`, `MUT-HEALTH-DROPS-FLAG`.

- [ ] **AC7 (M1) — Parent cancellation, writer release and serial race coverage:** Fresh-phase cancellation reaches registry-bootstrap with context.Canceled, settlement is bounded and writer authority reopens. Existing cancellation/settlement/context-root neighbours remain green. Mutating startup clocks and execStartupHook are serial; successful T1 lifecycle must be exercised outside sandbox under race. Tests: `TestRegistryPhasePreservesParentCancellation`, `TestServeStartsWhenTheProbeSpendsTheStartupClock`. Mutations: `MUT-PARENT-CANCEL-DETACHED`.

- [ ] **AC8 (M1) — Exact delivery and complete verification:** Apply the frozen 11-file patch byte-exact, assert all four new files and all thirteen new named tests, preserve docs/help and existing neighbours. Current-head focused gates, controller full gates, independent evaluation and PR/merge Linux plain/race CI are required; sandbox socket verdicts are uninformative. V27/V29 archive receipts are reusable only after all eleven hashes match. Tests: `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock`, `TestRegistryPhaseClockIsBounded`, `TestExecProbeOutOfTimeNamesTheClockNotAnArm`, `TestExecProbeGenuineRefusalRecordedAfterDeadline`, `TestDoctorNamesLeftOutExecProfiles`, `TestDoctorLateRefusalObservationAndReasonCap`, `TestServeStartupClocksArePinned`, `TestRegistryPhasePreservesParentCancellation`, `TestDaemonConstructorCallerInventory`, `TestDaemonCallerInventoryAliasControl`, `TestExecExclusionLogIsObservationOnly`, `TestHealthExecExclusionIsAdditive`. Mutations: all controls below plus exact-byte/presence/name assertions.

- [ ] **AC9 (M1) — Serve lifetime retains the original parent:** Listen, announce, serving, maintenance and shutdown remain owned by Run's original parent, never its spent startup clock. Require successful announced-service health plus the compiling serve-parent mutation killed at the health assertion. Tests: `TestServeStartsWhenTheProbeSpendsTheStartupClock`. Mutations: `MUT-SERVE-ON-STARTUP-CLOCK`.

M1 owns the production and test files capable of reddening every criterion; the machine `acceptance_checks` records each file explicitly. No criterion is assigned to another milestone or closed by a static T-label count.

## Exact named tests and examples

| Design label | Runnable top-level Go function |
|---|---|
| T1 | `TestServeStartsWhenTheProbeSpendsTheStartupClock` |
| T2 | `TestRegistryPhaseHasItsOwnClock` |
| T3 | `TestRegistryPhaseClockIsBounded` |
| T4 | `TestExecProbeOutOfTimeNamesTheClockNotAnArm` |
| T4b | `TestExecProbeGenuineRefusalRecordedAfterDeadline` |
| T5 | `TestDoctorNamesLeftOutExecProfiles` |
| T5b | `TestDoctorLateRefusalObservationAndReasonCap` |
| T6 | `TestServeStartupClocksArePinned` |
| T7 | `TestRegistryPhasePreservesParentCancellation` |
| T8 | `TestDaemonConstructorCallerInventory`, `TestDaemonCallerInventoryAliasControl` |
| T9 | `TestExecExclusionLogIsObservationOnly` |
| T10 | `TestHealthExecExclusionIsAdditive` |

There are **thirteen new top-level names**, including both caller inventory tests. The focused candidate command also requires these **ten existing neighbours**, each once:

- `TestStartupCallerCancellation`.
- `TestStartupBudgetQuarantinesPostCutoffBootstrap`.
- `TestStartupBudgetHeldConnectionSettles`.
- `TestHealthReportsBuildToolCredentialsAndProfiles`.
- `TestProfileForResolutionOrder`.
- `TestExecProfilesDirProbeFailureIsolatesOneProfile`.
- `TestDoctorDaemonVsForeignListener`.
- `TestDoctorHelpAndAPIKeyGuard`.
- `TestDoctorExecToolsWithAnM3Daemon`.
- `TestProductionContextRoots`.

The candidate runner asserts **23 exact top-level PASS identities**, no skipped/failed/duplicated identity and no race warning. Baseline requests the same regex and explicitly records all thirteen new names as absent expected feature gaps; its ten neighbour passes never count as new-test success. Missing names after application are a STOP.

No `.ail` is created or changed. Therefore loading a syntax teaching prompt, new Z3 contracts, effect rows and AILANG inline-test/demo modules are **not applicable**; no prompt-loaded version or invented `.ail` example is claimed. S1/S4/S5 remain binding on the untouched core, verified by controller `verify_ail.sh`. The Go truth tables and production mutants are the executable examples of these host-only rules. S7 is met by the exact QUICKSTART profiles-directory paragraph plus doctor help; code blocks and the existing runnable CLI happy path remain unchanged. No unmeasured numerical coverage percentage is promised, and no coverage command runs in this planner.

## Execution boundary and non-vacuity

Everything below is a **future controller execution contract**, not a record of this planning invocation running a gate. Planner performs only bounded reading and artifact validation. The controller runs socket and full gates **outside sandbox**. Any socket verdict inside workspace-write, green or red, is **UNINFORMATIVE UNDER SANDBOX** and cannot discharge acceptance.

Before any outside-sandbox execution the controller reviews the frozen diff, fixtures and runner for writes outside the worktree, including shared HOME logs/databases. Treat unreviewed writes as unverified and potentially destructive. Real artifact guards are **observation-only**: sha256 and line count before/after, with `absent` valid; never seed, truncate, append, rename or delete a real artifact to establish a guard. In particular, observe the caller's `~/.ailang/state/autopush.log` and `~/Library/Logs/ailang-world/worldd.log` without reading their contents into reports. Rehearse genuinely risky setup only under a synthetic root after review; never change HOME/temp/environment policy to conceal a red production gate. Normal test temp directories and compiler caches must be accounted for. Do not touch mission-v1 state or credentials.

Every command uses the dated [`run.py`](../verification/world-iter362/run.py), which logs UTC start/deadline, cwd, exact argv, uptime, resolved Go version and rc, preserves timeout failure as VOID, starts PATH with `/opt/homebrew/bin` and fixes `AILANG_BIN=$HOME/.pinned-ailang/ailang`. **No TMPDIR/GOTMPDIR overrides.** Pin version must be the clean released v0.41.0 accepted by repository gates; report a pin mismatch without swapping to a dev binary. Record pin version and hash with controller receipts. Full-suite logs additionally carry a bounded process CPU/load observation without environment dumps.

Baseline **every acceptance command** in a separate pristine real checkout, never in the tree roles edit. Baseline mutable/application checks use base-state assertions; candidate checks use post-apply assertions. Missing new tests/files are explicitly expected preimplementation gaps, never greens. Any other baseline red or timeout is reported with terminal STOP. The archive `git ls-tree HEAD` failure in V25 is neither a pass nor evidence of candidate regression. Do not seed Git history/index or invent a full-suite workaround. Plain native `go build ./...` is valid in **this World repo** (V24 baseline and candidate rc 0); false-greens' unrelated ailang-repo `cmd/wasm`/`gen/main` failure does not justify dropping it here.

Apply only after baseline acceptance. Exact steps (the JSON `apply` program also verifies inputs/head and all eleven output hashes):

```bash
git apply --check design_docs/planned/w-daemon-startup-under-load.patch
git apply design_docs/planned/w-daemon-startup-under-load.patch
test -f host/daemon/startup_callers_test.go
test -f host/daemon/startup_under_load_test.go
test -f host/broker/exec_probe_out_of_time_test.go
test -f cmd/ailang-worldd/doctor_exec_profiles_test.go
```

No fuzzy/reject/three-way/manual application. `apply-exact` recognizes an already exact candidate to avoid reapplication, otherwise requires all eleven base identities and the design commit as an ancestor (controller planning-only commits are permitted); partial/different content stops. The active baseline hold blocks application before either case. The controller may run it outside sandbox or accept the launched executor's once-only application and then run outside-sandbox identity/presence checks. Fresh-archive V27 plus fourth-presence V29 receipts may be reused on exact candidate identity; a new archive rehearsal is not required here. If one is separately authorized later, assert an independent Git root before application; enclosing-repo no-op applies returning rc 0 are not evidence.

## Concrete phased gates

Machine `verify_commands` contains the exact argv or reviewable Python `-c` program for each row. Programs are data, not automatically executed at planning time. All deadlines are ≤1800 s; mutable seams run serially with `-p 1 -parallel 1`; no mutation or broad-suite overlap, no background jobs.

| Gate ID | Phase | Deadline seconds | Acceptance |
|---|---|---:|---|
| `identity-baseline` | baseline | 120 | 11 base hashes match; four new files absent; frozen input hashes match |
| `build-baseline` | baseline | 300 | rc 0 on pristine World base |
| `vet-baseline` | baseline | 300 | rc 0 |
| `gofmt-baseline` | baseline | 120 | six existing Go files clean |
| `named-baseline` | baseline | 1200 | ten existing neighbours pass; thirteen absent new tests explicitly pending, not passed |
| `verify-go-baseline` | baseline | 1800 | pristine receipt observed rc 1: active STOP; rc-0 prerequisite not met |
| `verify-ail-baseline` | baseline | 1800 | baseline receipt observed rc 0; candidate gate still owed |
| `apply-exact` | implementation | 180 | git apply --check then once-only exact application; four presence guards and eleven candidate hashes |
| `identity-candidate` | focused | 120 | all eleven frozen hashes + four files; changed premise STOP |
| `build-candidate` | focused | 300 | rc 0 |
| `vet-candidate` | focused | 300 | rc 0 |
| `gofmt-candidate` | focused | 120 | all ten touched Go files clean |
| `named-candidate` | focused | 1200 | 23 exact top-level PASS identities once, including thirteen-new-test floor; no skip/fail/race |
| `reuse-v28-v29` | focused | 120 | eleven candidate hashes and frozen receipt hashes; only named empirical obligations reused |
| `verify-go-candidate` | controller_full | 1800 | real candidate checkout AFTER pristine rc-0 baseline; full plain suite rc 0 |
| `verify-ail-candidate` | controller_full | 1800 | full pinned-language gate rc 0; no .ail changes is not a waiver |

The executable core commands are `go build ./...`, `go vet ./host/broker/ ./host/daemon/ ./cmd/ailang-worldd/`, and read-only `gofmt -l` over the six existing baseline or ten candidate Go files. The named gate runs:

```bash
go test -race -count=1 -p 1 -parallel 1 -timeout 240s -v \
  -run '^(TestDoctorNamesLeftOutExecProfiles|TestDoctorLateRefusalObservationAndReasonCap|TestExecProbeOutOfTimeNamesTheClockNotAnArm|TestExecProbeGenuineRefusalRecordedAfterDeadline|TestDaemonConstructorCallerInventory|TestDaemonCallerInventoryAliasControl|TestServeStartsWhenTheProbeSpendsTheStartupClock|TestRegistryPhaseHasItsOwnClock|TestRegistryPhaseClockIsBounded|TestServeStartupClocksArePinned|TestRegistryPhasePreservesParentCancellation|TestExecExclusionLogIsObservationOnly|TestHealthExecExclusionIsAdditive|TestStartupCallerCancellation|TestStartupBudgetQuarantinesPostCutoffBootstrap|TestStartupBudgetHeldConnectionSettles|TestHealthReportsBuildToolCredentialsAndProfiles|TestProfileForResolutionOrder|TestExecProfilesDirProbeFailureIsolatesOneProfile|TestDoctorDaemonVsForeignListener|TestDoctorHelpAndAPIKeyGuard|TestDoctorExecToolsWithAnM3Daemon|TestProductionContextRoots)$' \
  ./host/daemon/ ./host/broker/ ./cmd/ailang-worldd/ ./host/store/
```

The JSON runner captures and validates each required name, rather than accepting `go test` rc alone. `TestProductionContextRoots` remains explicitly required even though this patch adds no root.

For exact invocation without shell interpolation of embedded programs, the controller uses this **reviewable driver snippet**. It resolves each named gate from the review copy and invokes the dated wrapper synchronously. Set `TARGET_CHECKOUT` to the controller's pristine sibling for baseline rows and the implementation checkout for candidate rows. `GATE_ID` selects one row above; controller must run every applicable row and may reuse its already completed matching baseline receipt instead of duplicating a running gate.

```python
from pathlib import Path
import json, subprocess
ARTIFACT_ROOT = Path('/Users/voightkampff/dev/sunholo-data/.wt-world-iter362-row256')
STATE = ARTIFACT_ROOT / 'design_docs/verification/world-iter362/planner-sprint-state.json'
# Supply TARGET_CHECKOUT and GATE_ID as explicit controller arguments, not secrets.
import sys
TARGET_CHECKOUT, GATE_ID = sys.argv[1:]
s = json.loads(STATE.read_text())
g = next(x for x in s['verify_commands'] if x['id'] == GATE_ID)
assert g['phase'] == 'baseline' or not s['execution_hold']['active'], s['execution_hold']['reason']
args = [x.replace('{state}', str(STATE)) for x in g['program_args']]
argv = ['python3', '-c', s['validation_programs'][g['program']], *args] if g['program'] else g['argv']
log = ARTIFACT_ROOT / 'design_docs/verification/world-iter362' / ('controller-m1-' + GATE_ID + '.txt')
assert not log.exists(), 'Do not overwrite an existing receipt; retain it and choose a new attempt name'
subprocess.run(['python3', str(ARTIFACT_ROOT / s['gate_execution']['runner']),
                str(g['deadline_seconds']), str(log), TARGET_CHECKOUT, *argv], check=True)
```

`verify-go-candidate` is exactly `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` in a **real candidate checkout after pristine baseline rc 0**. It is owed before independent evaluation, regardless of focused/reused evidence. `verify-ail-candidate` is exactly `./scripts/verify_ail.sh`; no `.ail` delta does not waive it. Both full gates remain controller responsibilities, never sandboxed executor milestone commands. The controller's pristine full verify_go is **rc 1** (active STOP above), and this plan claims no candidate full-suite success. Baseline verify_ail rc 0 is recorded separately; candidate verify_ail remains pending.

## Empirical evidence reuse and mutation ownership

Read V28/V29 current-byte design rows, controller summary, application receipt, final dated race, historical replay plus repaired no-log instrument, fourteen revision replay receipts, serve-parent source/control/build/mutant and the eleven-file prototype manifest. Evidence is Darwin arm64, Go 1.26.6, dated 2026-11-07; V28 final race is 05:29:08–05:29:18 UTC, load 6.49/9.00/7.25 at start and 6.10/8.83/7.21 at end. Earlier V23 sandbox results are not used as acceptance evidence.

Only if **all eleven post-apply file hashes** match the frozen manifest and the banked receipt hashes match, V28/V29 discharge:

1. Seventeen named outside-sandbox race/lifecycle passes: thirteen new plus startup caller cancellation, two settlement neighbours and context roots. The six remaining neighbours in the new 23-name command are still checked on the current implementation head; historical 15-test controls include the three doctor neighbours but do not replace that current-head command.
2. All nineteen historical production mutant names: **eighteen killed**, known `MUT-TEMPLATES-ON-STARTUP-CLOCK` survivor disclosed. The initial stale no-log anchor is **not** a kill; the repaired instrument's compile/red control is the evidence.
3. Fourteen revision production mutants compiled and killed their named assertions, with unmutated controls and exact restoration.
4. V29 `MUT-SERVE-ON-STARTUP-CLOCK`: T1 control rc 0; mutant `go build ./...` rc 0; T1 race rc 1 at announced-service health connection refusal; exact restoration.
5. V27 independent archive check/apply, eleven applied identities, with V29's fourth new-file guard. Current application and four current presence checks remain required.

This reuse **does not** discharge full candidate verify_go, verify_ail, independent evaluation, PR/merge Linux CI, real Linux srt or any new host-load/service-start claim. A changed implementation hash means **STOP/report changed premise**, not automatic replay or silently inheriting evidence. Do not launch new mutation runs while the controller baseline/full gates run. Any later separately authorized replay must assert the production edit changed bytes, compile it, name the exact firing assertion, restore from a byte-exact backup, and verify all eleven hashes. Never restore uncommitted work using `git checkout`.

| Production mutant (historical/revision duplicates retained as separate receipts) | Production file | Named assertion / disclosed gap |
|---|---|---|
| `MUT-RUN-HEAD-SHAPE` | `host/daemon/daemon.go` | `TestServeStartsWhenTheProbeSpendsTheStartupClock`: Run announces instead of registry-bootstrap expiry |
| `MUT-SHARED-PARENT` | `host/daemon/daemon.go` | `TestServeStartsWhenTheProbeSpendsTheStartupClock`: Run announces; T2 construction also succeeds |
| `MUT-REGISTRY-UNBOUNDED` | `host/daemon/daemon.go` | `TestRegistryPhaseClockIsBounded`: 1 ns clock returns registry-bootstrap deadline |
| `MUT-SWEEP-ON-STARTUP-CLOCK` | `host/daemon/daemon.go` | `TestRegistryPhaseHasItsOwnClock`: integrity sweep complete with no scan hole |
| `MUT-TEMPLATES-ON-STARTUP-CLOCK` | `host/daemon/daemon.go` | known survivor: no interpreter archive, reportTemplates returns at d.arch == nil |
| `MUT-SERVE-SPENT` | `host/broker/exec_startup.go` | `TestExecProbeOutOfTimeNamesTheClockNotAnArm`: only clean profile remains served |
| `MUT-SPENT-STILL-SERVED` | `host/broker/exec_startup.go` | `TestExecProbeOutOfTimeNamesTheClockNotAnArm`: only clean profile remains served |
| `MUT-FAIL-NEVER-SPENT` | `host/broker/exec_probe.go` | `TestExecProbeOutOfTimeNamesTheClockNotAnArm`: failure recorded after deadline is marked spent |
| `MUT-DEADLINE-UNSET` | `host/broker/exec_probe.go` | `TestServeStartsWhenTheProbeSpendsTheStartupClock`: health exclusion is flagged |
| `MUT-SPENT-ANY` | `host/broker/exec_probe.go` | `TestExecProbeOutOfTimeNamesTheClockNotAnArm`: mixed recording times false |
| `MUT-FLAG-NOT-CARRIED` | `host/broker/exec_startup.go` | `TestExecProbeOutOfTimeNamesTheClockNotAnArm`: excluded entry is flagged |
| `MUT-HEALTH-DROPS-FLAG` | `host/daemon/workspace.go` | `TestServeStartsWhenTheProbeSpendsTheStartupClock`: health carries true flag |
| `MUT-NO-LOG-LINE` | `host/daemon/workspace.go` | `TestServeStartsWhenTheProbeSpendsTheStartupClock`: startup log contains exact observation/review advice; moved anchor repaired |
| `MUT-DOCTOR-NOT-CALLED` | `cmd/ailang-worldd/doctor.go` | `TestDoctorNamesLeftOutExecProfiles`: two exclusion lines and exit 13 |
| `MUT-DOCTOR-WARN` | `cmd/ailang-worldd/doctor.go` | `TestDoctorNamesLeftOutExecProfiles`: exec profiles line is failure mark |
| `MUT-DOCTOR-CLASS` | `cmd/ailang-worldd/doctor_hosts.go` | `TestDoctorNamesLeftOutExecProfiles`: exec exit class is 13 |
| `MUT-DOCTOR-ONE-FIX` | `cmd/ailang-worldd/doctor.go` | `TestDoctorNamesLeftOutExecProfiles`: true and false flags have distinct conditional advice |
| `MUT-REGISTRY-BUDGET-30S` | `host/daemon/daemon.go` | `TestServeStartupClocksArePinned`: registry clock is exactly 9 s |
| `MUT-STARTUP-BUDGET-5S` | `host/daemon/daemon.go` | `TestServeStartupClocksArePinned`: startup clock is exactly 9 s |
| `MUT-LATE-NEVER-SPENT` | `host/broker/exec_probe.go` | `TestExecProbeGenuineRefusalRecordedAfterDeadline`: late genuine refusal permits OutOfTime=true |
| `MUT-SPENT-ANY` | `host/broker/exec_probe.go` | `TestExecProbeOutOfTimeNamesTheClockNotAnArm`: mixed recording times must yield false |
| `MUT-LATE-STILL-SERVED` | `host/broker/exec_startup.go` | `TestExecProbeGenuineRefusalRecordedAfterDeadline`: late refusal leaves served set empty and Only nil |
| `MUT-FLAG-NOT-CARRIED` | `host/broker/exec_startup.go` | `TestExecProbeGenuineRefusalRecordedAfterDeadline`: excluded late entry remains flagged |
| `MUT-CAUSAL-OBSERVATION` | `host/broker/exec_probe.go` | `TestExecProbeGenuineRefusalRecordedAfterDeadline`: shared diagnostic equals observation-only sentence |
| `MUT-LOG-CAUSE` | `host/daemon/workspace.go` | `TestExecExclusionLogIsObservationOnly`: log preserves reason and exact observation/review advice |
| `MUT-DOCTOR-CAUSE` | `cmd/ailang-worldd/doctor.go` | `TestDoctorLateRefusalObservationAndReasonCap`: late doctor detail contains observation-only text |
| `MUT-DOCTOR-NO-REVIEW` | `cmd/ailang-worldd/doctor.go` | `TestDoctorLateRefusalObservationAndReasonCap`: doctor advises review before restart |
| `MUT-REASON-CAP-240-PREFIX` | `cmd/ailang-worldd/doctor.go` | `TestDoctorLateRefusalObservationAndReasonCap`: one-line reason cap includes ellipsis and no 241-byte suffix |
| `MUT-HEALTH-NO-OMITEMPTY` | `host/daemon/daemon.go` | `TestHealthExecExclusionIsAdditive`: false form equals exact existing JSON bytes |
| `MUT-PARENT-CANCEL-DETACHED` | `host/daemon/daemon.go` | `TestRegistryPhasePreservesParentCancellation`: canceled phase parent returns registry-bootstrap cancellation |
| `MUT-CALLER-IMPORTED-ALIAS` | `cmd/ailang-worldd/startup_caller_mutant.go` | `TestDaemonConstructorCallerInventory`: unreviewed production startup caller is rejected |
| `MUT-CALLER-DIRECT` | `cmd/ailang-worldd/startup_caller_mutant.go` | `TestDaemonConstructorCallerInventory`: unreviewed production startup caller is rejected |
| `MUT-CALLER-LOCAL-ALIAS` | `host/daemon/startup_caller_mutant.go` | `TestDaemonConstructorCallerInventory`: unreviewed production startup caller is rejected |
| `MUT-SERVE-ON-STARTUP-CLOCK` | `host/daemon/daemon.go` | `TestServeStartsWhenTheProbeSpendsTheStartupClock`: announced daemon answers GET /v1/health; mutant connection refused |

The template-clock mutant survives because no rig archives an interpreter (`d.arch == nil`); no claim that it was killed is permitted. Bootstrap and sweep have dynamic clock tests; template context selection additionally has source inspection and an explicit residual, not a fabricated dynamic proof. Mutation coverage is behavioral, not a grep count.

## Pending independent evaluation and Linux landing obligations

Independent evaluation is controller-routed **after** candidate full Go and AILANG gates, with exact implementation head, all AC ownership, mutation evidence/reuse identity, known survivor and residuals included. No evaluator or executor is launched by this planner, and no message is sent. The evaluator's mutation work must be sequenced after full suites on a quiet rig; findings requiring new direction/policy stop rather than changing the frozen patch.

On **both PR head and merge head**, controller records CI completeness and Linux success for the existing `ailang-code verify gate`, `go host build + test gate` (plain) and `go race gate (go test -race, own job)`, including the existing verbose real-srt confinement/startup checks (`TestExecSrtStartupProbePasses` and their named pass guards). Require the thirteen new tests to execute rather than disappear/skip; collect named receipts if full-suite summaries do not expose them. CI uses its existing checksum-pinned Linux binary setup; a Darwin gate or sandbox result cannot substitute. No workflow/timeout changes are planned; existing Linux CI configuration remains authoritative. Local planned wrappers are ≤1800 s; existing CI race job's outer 35-minute ceiling is not being changed or claimed as a new local deadline. All CI obligations and milestone completion remain **pending**.

## Risks, STOP rules and handoff

Context-free archives/setup and scheduling remain unbounded; the two 9 s context allowances do not establish total startup below the service's 30 s wait. A per-run 8 s timeout recorded early can yield false, while a genuine refusal recorded late can yield true. The slow-client fixture reproduces spending a context allowance, **not CPU/I/O contention or incident cause**. No retry policy, scheduling bound, high-load acceptance or second/third historical-start reproduction is established. Linux/real-srt and launchd/service end-to-end remain outside the measured Darwin evidence; template reporting retains its disclosed mutation gap.

- STOP: preserve current files, report BLOCKED: <specific changed premise, preexisting red, or new judgment call>, and end the run. Do not commit or dispatch; controller owns commits and routing.

- Any design/patch hash mismatch, non-base/non-candidate implementation hash, partial apply, missing required file or changed empirical receipt stops reuse and execution; do not repair or regenerate the frozen patch.

- Every acceptance command is baselined in a separate pristine real checkout. Any preexisting red or timeout is reported/STOP; never widen scope, change test/environment policy, seed Git history, override TMPDIR/GOTMPDIR or weaken a gate to hide it.

- Baseline absence of thirteen new tests is an explicit expected feature gap, not a passed test, and is distinguished from any existing-neighbour red. Candidate missing/skip/duplicate tests stop.

- Socket results inside workspace-write, green or red, are UNINFORMATIVE UNDER SANDBOX. Controller outside-sandbox rerun is required for verdicts.

- Review writes outside the worktree before any outside-sandbox execution. Observe real artifacts only by sha256 and line count, absent valid, never seed/truncate. Unknown real-state writes stop; rehearse risky setup only under a synthetic root after review.

- Do not run mutations or broad suites concurrently with the controller pristine baseline or each other. No background jobs, stress spinners, CPU/I/O contention experiment or parallel seam mutations.

- No human-decision gates remain in row 256. D-WORLD-105 belongs to parked row 233, D-WORLD-110 to row 255. Genuine new policy/direction judgment parks with terminal STOP; helper names and evidence layout do not reopen approval.

Skill template STOP wording normally includes committing work. This authorized file-only task overrides it: **preserve files, report `BLOCKED: <reason>`, and end the run; no commits or messages**. The controller owns any later commit, PR, merge, reporting and routing. No existing ruled decision is reopened.

Deliverables are this plan, populated ignored state `.ailang/state/sprints/sprint_W-DAEMON-STARTUP-UNDER-LOAD.json`, its exact review copy `design_docs/verification/world-iter362/planner-sprint-state.json`, file-only executor handoff `design_docs/verification/world-iter362/planner-handoff.json`, and compact hash manifest `design_docs/verification/world-iter362/planner-artifacts.json`. A fresh executor worktree lacks ignored state: controller must install the review copy byte-exact at the ignored path and verify hashes before launch, or supply the complete review copy inline; do not assume ignored state transfers through Git. The executor must honor the active baseline STOP before any application or candidate gate. Do not invoke the skill's issue-import/message/commit scripts.

Success metrics are all ACs evidenced, exact frozen application, thirteen new named tests plus ten neighbours, behavioral mutation ownership with the one disclosed survivor, clean focused gates, controller full gates, independent evaluation and PR/merge Linux CI. They are requirements, **not passed results**. No background jobs were started by this planner.

# Row 267 — daemon race package timeout headroom

Iteration 443. Planner and executor lane `codex:gpt-6.2-sol`. One cumulative milestone, M1. Code base `6420426efcd0304c88adf347cf7af0f158d96863`. Design head `c50264d6fe6d1a23184a4289ec214d023528c909`.

## 0. Outcome and size

M1 gives the CI race leg a **900 s per-package bound**, a **995 s start allowance**, a guard at **1955 s** (900 + 60 + 995), a **2000 s leg bound**, and a **40-minute step ceiling**. The wrapper reports its process group before SIGKILL and returns 124. The diagnostic gets up to two seconds. The order in seconds is 900 < 960 < 1955 ≤ 2000 < 2300 ≤ 2400. Use the reviewed r2 §0 amendments wherever §1–§3 cite r1. Round 1 blocked 3/3. Round 2 returned glm PASS and kimi/astra REJECT with concrete fixes, which the controller applied under the narrow-refinement carve-out.

`$S` is `$HOME/.ailang/state/world-iter443/plan-scratch`. `$STATE` is its parent. `$W` is `/Users/voightkampff/dev/sunholo-data/.wt-world-iter443-row267`. `$B` is the full code base above.

Expected `git diff -M --numstat 6420426e HEAD -- scripts host .github docs` after M1 is **+294/−57 = 351 changed lines**, measured from the archive and stage. The design estimated about 262 lines. The watchdog, bounded cleanup and extra evidence control account for part of the difference. The cap is 400. At design HEAD the filtered command is empty. The unfiltered command at design HEAD reports `535	0	design_docs/planned/w-daemon-race-package-timeout-headroom.md`. With the design commits and M1, the expected unfiltered committed sum is **886 lines**. The two untracked planner artifacts are outside that HEAD-to-base sum.

````text
70	28	scripts/verify_go.sh
110	0	scripts/test_race_leg_bounds.sh
87	18	host/verifygate/race_job_gate_test.go
22	7	host/verifygate/mission_config_gate_test.go
1	1	host/verifygate/evidence_manifest_gate_test.go
4	3	.github/workflows/ci.yml
````

## 1. Measurements

### 1.1 Base, anchors and the two r2 pending items

| Input | Command and observed output |
|---|---|
| Worktree | `git status --short --branch; git rev-parse HEAD; git rev-parse 6420426e` → branch `sprint/world-iter443-row267-race-headroom`, initially clean, HEAD `c50264d6fe6d1a23184a4289ec214d023528c909`, base `6420426efcd0304c88adf347cf7af0f158d96863`. |
| Existing shell bounds and leg | `nl -ba scripts/verify_go.sh \| sed -n '217,226p;342,364p'` → legs dispatch 217–221, default `1500` at 222, lexical check 223–225, banner 226. Early anchors are mission helper 142, evidence dispatch 158, mission dispatch 167, unknown-option check 176, AILANG_BIN preflight 181 and tracked-binary fixture 239. Race block 349–364, old banner 350, Python command with `8m` at 355, process-group kill at 360, FATAL at 361 and exit 124 at 362. `grep -n -- '-timeout' scripts/verify_go.sh` matched 350 and 355 only. |
| Existing workflow step | `nl -ba .github/workflows/ci.yml \| sed -n '520,526p'` → race step 520, comment 521, `timeout-minutes: 35` at 522, `VERIFY_GO_LEGS: race` at 524, leg `"1500"` at 525, gate invocation 526. `grep -nF 'continue-on-error' .github/workflows/ci.yml` returned 1 with no matches. Positive controls were 19 `timeout-minutes:` lines and two Go gate invocations. |
| Existing Go test anchors | `grep -nE 'TestRaceLegRunsInItsOwnJob\|TestVerifyGoProductCheckSequence\|requireVerifyGoProductSequence\|TestEvidenceNamedManifestRejectsUnpinnedTest\|TestZ3PinDeclaredOnceAndInstalledInBothJobs' host/verifygate/{race_job,mission_config,evidence_manifest,ail_binary}_gate_test.go` → definitions at race 15, mission helper 443 and test 536, evidence 229 and Z3 699. `nl -ba` showed race env at race 37 and script needles 55–58, legacy mission fixture 412–441 (race literals 436–439), repaired helper literals 500–503, evidence copy/execution/rewrite at 98/106/169/187/212 and Z3 controls at 708. Toolchain controls at 123, jobs at 159–160, P1/P2 read and bind at 1189–1340; `nl -ba host/pins/scan_test.go \| sed -n '150,161p'` → `verify_go.sh` entry 158, `RoleGatePrimary, n: 4` at 159. |
| Daemon per-`k` | Copied `$STATE/design-scratch/classes.py` and `$STATE/racelogs/` to `$S/classes-input/`. `(cd "$S/classes-input" && python3 -I daemon-trend.py)` uses the original `classes.py` rows and calculates daemon seconds divided by `k`. Output `attempts (deduped, complete): 51`, 13 mid-class at 1.057–1.321, four daemon-capped and nine uncapped (1.0565–1.2449). Uncapped per-`k` by ascending `k` was **387.42, 385.95, 381.70, 365.26, 379.48, 407.15, 364.90, 361.63, 357.84** s. Pearson `r=-0.583`, least-squares slope `-137.1 seconds/k`. The observed mid-class trend falls with a local increase at 1.1509. Censoring leaves the slow-tail trend unmeasured. The r2 bounds still apply. Raw rows and IDs in `$S/daemon-trend.raw`. |
| Charter wording | `grep -n 'keep the package' design_docs/world-mission.md` → one match at **6372**, row 267, “Row 243's S5f2 adds process pins to `cmd/ailang-worldd`; its directive must keep the package's race time flat or move the bound.” It supports the S5f2 wording. The design's attribution to the S6 directive is a paraphrase without a matching charter entry. |

`go help testflag | grep -A3 -- '-timeout d'` reports the implicit 10-minute default. `grep -nF '696 s' design_docs/verification/world-iter442/readings.md` reports a completed run with an explicit `-timeout 45m`. The 696 s reading therefore cannot certify default plain-leg capacity. `grep -nF "testKillTimeout = testTimeout + 1*time.Minute" "$(go env GOROOT)/src/cmd/go/internal/test/test.go"` read Go's backup at line 841.

### 1.2 Reader inventory

`grep -RnE 'verify_go\|ci.yml\|VERIFY_GO_\|timeout-minutes\|go-race-verify' host/ scripts/ .github/` is recorded in `$S/readers-grep.raw` (230 lines). The literal positive controls `grep -cF './scripts/verify_go.sh' .github/workflows/ci.yml` and `grep -cF 'timeout-minutes:' .github/workflows/ci.yml` printed **2** and **19**. A second search for path components catches generic workflow enumeration. The reader result column gives the transcript test names exactly. Unmodified named baseline AC4 printed five `--- PASS:` lines, because the new race-ordering test is absent at base.

| Reader at 6420426e | Contract | M1 receipt |
|---|---|---|
| `race_job_gate_test.go:15,68` | Separate plain/race jobs, CI env and runner gate, pinned interpreter, known-positive control, loud bound, early unknown-leg refusal. | `--- PASS: TestRaceLegRunsInItsOwnJob`, `--- PASS: TestVerifyGoRefusesUnknownLegs`, and new `--- PASS: TestRaceBoundsAreOrdered`. |
| `mission_config_gate_test.go:191,206,253,264,272,311,350,383,443,536` | Real charter subcommand, bad-ledger refusals, blocked-child timeout classification, copied script bash syntax, preflight passes the tracked-binary gate before the ledger and stops on invalid ledger, exactly three active CI jobs free of old driver consumer, broad retired-token scan, retired option refuses early, historical and live ordered product anchors. | `--- PASS: TestMissionConfigRealCharter`, `--- PASS: TestMissionConfigRejectsBadLedger`, `--- PASS: TestMissionConfigBlockingTimeoutClassified`, `--- PASS: TestMissionConfigBashSyntax`, `--- PASS: TestActiveWiringHasNoDriverConsumers`, `--- PASS: TestRetirementCompleteness`, `--- PASS: TestRetiredModeIsRefused`, `--- PASS: TestVerifyGoProductCheckSequence`. `--- FAIL: TestMissionConfigDefaultFlowIntegration` timed out inside its pre-existing tracked-binary fixture, matching the controller base fail-set entry. See §1.3. |
| `evidence_manifest_gate_test.go:229` | Copies current gate, executes early evidence dispatch and rewrites its required-set marker. M1 pins the known-good output byte for byte, so a new top-level echo is observable. | `--- PASS: TestEvidenceNamedManifestRejectsUnpinnedTest`. |
| `ail_binary_gate_test.go:699` | Workflow controls and counts of solver declarations, checksum downloads and installs across the three jobs. | `--- PASS: TestZ3PinDeclaredOnceAndInstalledInBothJobs`. |
| `toolchain_pin_gate_test.go:115,1090,1322,1431` | Exact job set and per-job toolchain pins, miscompile-step placement/fail-closed wiring, P1 floor and P2 observed-toolchain binding. The three-FATAL count is scoped to the P1 region. | `--- PASS: TestGoToolchainPinsAgreeAndMatchJobList`, `--- PASS: TestMiscompileInstrumentStepIsGatedInCI`, `--- PASS: TestRaceControlFloorStaysBelowRootToolchain`, `--- PASS: TestP1NeedleSetIsInertUnderRenameAndCatchesReversal`. |
| `dispatch_lever_gate_test.go:274,325,420` | Real workflow triggers as parser and control fixtures. Generic workflow enumeration enforces `workflow_dispatch`. | `--- PASS: TestOnBlockTriggerParserShapes`, `--- PASS: TestOnBlockControlNeedle`, `--- PASS: TestEveryWorkflowDeclaresDispatchLever`. |
| `host/pins/manifest_test.go:214` and `scan_test.go:253,305,389` | Actual install asset/digest/order, per-file release counts (four in `verify_go.sh`, including a repeated token in its error), injected-token refusals and production file classification. | `--- PASS: TestCIInstallsMatchManifest`, `--- PASS: TestReleaseLiteralsExplainedByManifest`, `--- PASS: TestReleaseLiteralScanCatchesInjectedLiterals`, `--- PASS: TestIsProductionFile`. Receipt `$S/readers-stage-git.raw` used `GIT_DIR=$(git -C "$W" rev-parse --absolute-git-dir)` and `GIT_WORK_TREE="$S/st-M1" GIT_OPTIONAL_LOCKS=0` for read-only tracked-file enumeration from the archive. |
| `host/runbook/runbook_stageb_test.go:351` | Scans workflow, shell and tool files for the publish entrypoint. Positive controls counted six Go gate mentions and one AIL package-checker mention across 25 staged files (base 24). | `--- PASS: TestNoCIStepOrScriptReachesThePublishEntrypoint` in `$S/runbook-st-M1.raw`. |
| `scripts/test_world_contract_history.py:34` | Workflow ratchet event and base inputs and checker invocation. | `CI event/base wiring: PASS`, `world contract history controls PASSED` in `$S/readers-stage-git.raw` using `python3 -B scripts/test_world_contract_history.py` with the same read-only Git metadata. |

Other census matches are workflow invocations, fixture or historical data, and prose. The tracked-binary fixture repeats the live allowlist and tests its detector directly. `scripts/floor/canary.py:95` is a docstring and makes no assertion on the gate text. The direct shell/workflow invocation grep returned no match, but `grep -nF "import canary" scripts/floor/{run,test_canary}.py` found Python imports in the floor runner and its fixture. `host/archive/archive.go:66` already described a stale 600 s outer wait. The new race-ordering test is the first consumer of `timeout-minutes` in this inventory. `$S/st-readers.raw` contains the stage named runs. In the plain archive the release scan predictably refused at `git ls-files`. The read-only rerun above reached the actual predicate.


Additional named-reader replay from the stage root used these commands and the pinned environment from §6. The read-only index setup reached the same source bytes in the archive.

````bash
go test -count=1 -v ./host/verifygate/ -run '^(TestMissionConfigRealCharter|TestMissionConfigRejectsBadLedger|TestMissionConfigBlockingTimeoutClassified|TestMissionConfigBashSyntax|TestMissionConfigDefaultFlowIntegration|TestActiveWiringHasNoDriverConsumers|TestRetirementCompleteness|TestRetiredModeIsRefused|TestGoToolchainPinsAgreeAndMatchJobList|TestMiscompileInstrumentStepIsGatedInCI|TestRaceControlFloorStaysBelowRootToolchain|TestP1NeedleSetIsInertUnderRenameAndCatchesReversal|TestOnBlockTriggerParserShapes|TestOnBlockControlNeedle|TestEveryWorkflowDeclaresDispatchLever)$'
go test -count=1 -v ./host/runbook/ -run '^TestNoCIStepOrScriptReachesThePublishEntrypoint$'
GIT_DIR=$(git -C "$W" rev-parse --absolute-git-dir) GIT_WORK_TREE="$PWD" GIT_OPTIONAL_LOCKS=0 \
  go test -count=1 -v ./host/pins/ -run '^(TestCIInstallsMatchManifest|TestReleaseLiteralsExplainedByManifest|TestReleaseLiteralScanCatchesInjectedLiterals|TestIsProductionFile)$'
````

### 1.3 Matched base and stage

Before the M1 test runs, `git archive 6420426e | tar -x -C "$S/pristine"` created the pristine source export. The same archive seeded `$S/st-M1`. Local base AC1–AC3 reported ac1 rc=2, ac2-refuse rc=2, ac2-accept rc=2, ac3 rc=127, ac3-contended rc=127. The old script refused the new switches as unrecognized, and the new fixture did not exist. AC4 native vet rc 0, tests rc 0 (five existing names). Base build rc 0 and `verify_ail.sh` rc 0. AC5 is SKIPPED at base because its anchors are added by M1. AC7/AC8 require PR or merge artifacts.

Pristine AC6 used `go test ./... -count=1 -timeout 45m` from the completed archive. rc 1. **UNINFORMATIVE UNDER SANDBOX**. It hit the socket-bind denial above. The source export has no Git index, which also affects tracked-file and history readers. The controller's actual checkout baselines below are the acceptance comparator. The parser control is `FAILSET_CONTROL host/store TestParserPositive`. The raw logs are checked for `panic: test timed out`. The pristine all-package parser control fired; its text fail set is in `$S/base-archive-all-plain.fails`.

| Controller base leg | Result |
|---|---|
| `verifygate-plain` | rc 1, host/verifygate TestMissionConfigDefaultFlowIntegration, FAILSET_CONTROL fired, timeout panic absent. |
| `pins-plain` | rc 0, empty fail set, FAILSET_CONTROL fired, timeout panic absent. |
| `all-plain` | rc 1, 49 failing tests (`cmd/ailang-worldd` 20, `cmd/world-publish` 1, `host/broker` 7, `host/coordinator` 1, `host/daemon` 19, `host/verifygate` 1). Full fail set at `$STATE/base/all-plain.fails`, FAILSET_CONTROL fired, timeout panic absent. |

Baseline marker `ALL_DONE for 6420426efcd0304c88adf347cf7af0f158d96863`. The controller script uses `go test -count=1 -timeout 30m ./host/verifygate/`, `go test -count=1 -timeout 10m ./host/pins/`, and `go test -count=1 -timeout 45m ./...`.

The first aggregate gates rehearsal stopped after verifygate when its script was edited in flight. The final extracted script passes `bash -n`. Pins, direct all-package and both AC3 environments have separate raw receipts.

Stage AC1/AC2 logged the three expected values and exit codes in `$S/st-ac1-ac2.raw`. Strict AC4 named rc 0. Native and Linux-targeted vet rc 0 and 0. Build rc 0. AIL gate rc 0. Shell syntax and gofmt are in `$S/st-gates-progress.raw`. The exact `gofmt -l host/verifygate/` also printed no files.

| Stage gate | Result |
|---|---|
| `verifygate-plain` | rc 1, **INSTRUMENT FAILURE** (30-minute timeout panic). The attributable fail set was host/verifygate TestMissionConfigDefaultFlowIntegration, the same as controller base. Repeat the package outside under the §7 gate. |
| `pins-plain` | rc 0 with read-only index; empty fail set matches controller base. |
| `all-plain` | rc 1, **UNINFORMATIVE UNDER SANDBOX** (socket bind in `$S/st-archive-all-plain.raw`). Parser control fired, 37 attributable names, timeout panic absent. Repeat against the controller base outside. |

`$S/pins-stage-whole-git.raw` also records the entire `host/pins` package with the read-only index, rc 0 and an empty fail set. Its fail set matched `$STATE/base/pins-plain.fails` byte for byte.

A non-zero test exit with an empty parse is an instrument failure, even when in a leg with a separate sandbox limitation. Compare only attributable package/test names and inspect timeout lines separately.

## 2. M1 — apply the compiled patch

Appendix P carries every M1 edit. Apply that fence as a single patch in the executor's current worktree. `scripts/verify_go.sh` defines the limits and watchdog before dispatch, validates the main flow, and runs the shell fixture before the race call. `.github/workflows/ci.yml` sets the race step ceiling and repeats the leg bound. The host tests pin the numeric ordering, the exact call and repaired product order, and the evidence subcommand's known-good output. The new shell fixture is mode 0644 and the race gate invokes it through Bash. The measured integer components 554 + 93 + 347 sum to 994. R1 pins the reviewed allowance at 995.

The evidence assertion is one line in an existing sixth file. The original success control accepted an arbitrary leading echo, so M19 requires the stronger assertion. The embedded mission fixture still contains the original 8m command and ordered markers. The repaired arm locates the runtime guard as a whole line, so the indented diagnostic subcommand cannot satisfy M17.

## 3. Operator view

The race job prints the following lines. A timeout adds one line per process present in the leg group, with elapsed time and argv. A failed or stalled listing emits one `could not list race leg processes` diagnostic. The wrapper attempts SIGKILL on the group and exits 124.

```text
── legs: race (race bounds: per-package 900s, leg 2000s)
PASS 7/7 arms
── go test ./... -count=1 -race -timeout 900s (leg bound 2000s)
verify_go.sh: still running at the leg bound (elapsed 33:20): /…/b123/daemon.test -test.timeout=15m0s
verify_go.sh: FATAL: -race leg timed out after 2000s
```

A package alarm can report `panic: test timed out after 15m0s` with the test and package. An undersized override exits 2 before the legs, with the per-package, backup kill and start-allowance arithmetic in its diagnostic.

## 4. Executor handoff

From `$W`, extract the patch, apply it, run acceptance, then create the snapshot only on a successful acceptance. The plan fence extraction is also the controller's extraction protocol.

````bash
S="$HOME/.ailang/state/world-iter443/plan-scratch"
plan=design_docs/planned/w-daemon-race-package-timeout-headroom-sprint-plan.md
awk -v m=M1 '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$plan" > "$S/executor-M1.patch"
test -s "$S/executor-M1.patch" &&
  /usr/bin/patch -p1 -E < "$S/executor-M1.patch" &&
  bash "$S/accept.sh" &&
  bash "$S/snapshot.sh" M1
````

Use §7 fences for the executable scripts. The planner scratch copies also exist at those paths. The planner's sandbox denied real `ps`. If the executor has the same restriction, the controller runs full AC3 and AC5 on the applied tree outside the sandbox and establishes the `ACCEPT_OK` receipt before committing. The controller commit script refuses untracked `design_docs/*` files absent from the product snapshot. Preserve both plan artifacts outside that worktree or land them in the controller's documentation commit before calling `commit_ms.sh`.

## 5. Mutation ledger

Anchor lines refer to the frozen stage. The complete receipts live in `$S/mutants.txt` and the raw per-mutation directories. Go killers replay against the named test(s) and the whole `host/verifygate` package. Static text predicates and behavior checks are identified below. Each Go fail-set parser prints its positive control. Whole-run kills followed by a package timeout are marked explicitly. Repeat those legs outside the sandbox for a complete package result. Shell rows report the command's mutated/restored rc. For R1–R5 restores, the expected arm was observed green before sandbox-denied R6 stopped the same fixture.

| Mutation | Single edit on staged source | Named test verdict and whole-package replay |
|---|---|---|
| M1 | P=480 `scripts/verify_go.sh:161` | `R1` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole KILLED<br>`TestRaceBoundsAreOrdered` (static text) 1→0, **KILLED**, whole rc 1 KILLED (intended failure reached; package timed out afterward) |
| M2 | hardcode 8m in race call `scripts/verify_go.sh:405` | `TestRaceBoundsAreOrdered` (static text) 1→0, **KILLED**, whole rc 1 KILLED (intended failure reached; package timed out afterward)<br>`TestVerifyGoProductCheckSequence` (static text) 1→0, **KILLED**, whole rc 1 KILLED (intended failure reached; package timed out afterward) |
| M3 | CI L=1500 `.github/workflows/ci.yml:526` | `TestRaceBoundsAreOrdered` (static text) 1→0, **KILLED**, whole rc 1 KILLED (intended failure reached; package timed out afterward) |
| M4 | default L=1900 `scripts/verify_go.sh:163` | `R1` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole KILLED<br>`TestRaceBoundsAreOrdered` (static text) 1→0, **KILLED**, whole rc 1 KILLED (intended failure reached; package timed out afterward) |
| M5 | CI step=36 `.github/workflows/ci.yml:523` | `TestRaceBoundsAreOrdered` (static text) 1→0, **KILLED**, whole rc 1 KILLED (intended failure reached; package timed out afterward) |
| M6 | short-circuit guard `scripts/verify_go.sh:165` | `R3` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole SKIPPED |
| M7 | comparison -le `scripts/verify_go.sh:171` | `R2` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole SKIPPED |
| M8 | A=813 `scripts/verify_go.sh:162` | `R1` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole SKIPPED<br>`TestRaceBoundsAreOrdered` (static text) 1→0, **KILLED**, whole intended failure observed; terminal receipt PENDING; timeout panic in raw log |
| M9 | delete ps listing `scripts/verify_go.sh:190` | `R6` (behavior) 1→1, **UNINFORMATIVE UNDER SANDBOX**, whole SKIPPED (mutated: FAIL R6 hang: expected a process listing naming the probe (rc=124)) |
| M10 | kill child only `scripts/verify_go.sh:200` | `R6` (behavior) 1→1, **UNINFORMATIVE UNDER SANDBOX**, whole SKIPPED (mutated: FAIL R6 hang: grandchild 83726 still alive after reap poll (rc=124)) |
| M11 | swallow child rc `scripts/verify_go.sh:185` | `R5` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole SKIPPED |
| M12 | timeout rc 0 `scripts/verify_go.sh:204` | `R6` (behavior) 1→1, **UNINFORMATIVE UNDER SANDBOX**, whole SKIPPED (mutated: FAIL R6 hang: expected rc 124 and timeout diagnostic (rc=0)) |
| M13 | append -short `scripts/verify_go.sh:405` | `TestRaceBoundsAreOrdered` (static text) 1→0, **KILLED**, whole intended failure observed; terminal receipt PENDING; timeout panic in raw log |
| M14 | drop -race `scripts/verify_go.sh:405` | `TestRaceBoundsAreOrdered` (static text) 1→0, **KILLED**, whole intended failure observed; terminal receipt PENDING; timeout panic in raw log<br>`TestRaceLegRunsInItsOwnJob` (static text) 1→0, **KILLED**, whole intended failure observed; terminal receipt PENDING; timeout panic in raw log |
| M15 | continue-on-error true `.github/workflows/ci.yml:523` | `TestRaceBoundsAreOrdered` (static text) 1→0, **KILLED**, whole intended failure observed; terminal receipt PENDING; timeout panic in raw log |
| M16 | drop fixture call `scripts/verify_go.sh:402` | `TestVerifyGoProductCheckSequence` (static text) 1→0, **KILLED**, whole intended failure observed; terminal receipt PENDING; timeout panic in raw log |
| M17 | drop main-flow guard `scripts/verify_go.sh:282` | `TestVerifyGoProductCheckSequence` (static text) 1→0, **KILLED**, whole intended failure observed; terminal receipt PENDING; timeout panic in raw log |
| M18 | omit backup minute `scripts/verify_go.sh:170` | `R3` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole SKIPPED |
| M19 | top-level guard and echo `scripts/verify_go.sh:208` | `R1` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole SKIPPED<br>`TestEvidenceNamedManifestRejectsUnpinnedTest` (behavior) 1→0, **KILLED**, whole intended failure observed; terminal receipt PENDING; timeout panic in raw log |
| M20 | remove ps timeout `scripts/verify_go.sh:192` | `R7` (behavior) 1→0, **KILLED**, whole SKIPPED |
| M21 | move SIGKILL after listing outside finally `scripts/verify_go.sh:188` | `R7` (behavior) 1→0, **KILLED**, whole SKIPPED |
| X01 | feed constant into guard case `scripts/verify_go.sh:167` | `R4` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole SKIPPED |
| X02 | disable bounds dispatch `scripts/verify_go.sh:226` | `R1` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole SKIPPED |
| X03 | disable wrapper dispatch `scripts/verify_go.sh:232` | `R5` (behavior) 1→1, **KILLED (named arm restored; later R6 denied)**, whole SKIPPED |
| X04 | remove process-group filter `scripts/verify_go.sh:196` | `R6` (behavior) 1→1, **UNINFORMATIVE UNDER SANDBOX**, whole SKIPPED (mutated: FAIL R6 hang: expected a process listing naming the probe (rc=124)) |
| X05 | do not catch subprocess listing errors `scripts/verify_go.sh:201` | `R7` (behavior) 1→0, **KILLED**, whole SKIPPED |
| X06 | remove child wait timeout `scripts/verify_go.sh:185` | `R6` (behavior) 1→1, **UNINFORMATIVE UNDER SANDBOX**, whole SKIPPED (mutated: FAIL R6 hang: expected rc 124 and timeout diagnostic (rc=1)) |
| X07 | remove partial ps-row guard `scripts/verify_go.sh:196` | `malformed` (behavior) 1→0, **KILLED**, whole SKIPPED |
| X08 | allow a missing group to add a spurious listing diagnostic `scripts/verify_go.sh:199` | `missing` (behavior) 1→0, **KILLED**, whole SKIPPED |

Drop `-ww`: **UNOBSERVED** by design in planning and CI. Do not infer a kill from the unmutated Linux R6 run.

Controlled heredoc probes also rejected M9 and X04 by their missing-list and wrong-group diagnostics. The live R6 verdicts remain uninformative. R7 replay keeps the original timed body and skips denied R6 in a scratch copy. The X07 and X08 behavior probes load the exact staged Python heredoc with controlled subprocess results, exercise the malformed-output and disappeared-group paths, and assert the same 124, listing output and group kill. No extra PATH shims are used. Replay paths are in §6.


## 6. Acceptance

Run from the target tree with `export AILANG_BIN="$HOME/.pinned-ailang/ailang" GOTOOLCHAIN=go1.26.6`. The observed version lines were `AILANG v0.41.0` and `go version go1.26.6 darwin/arm64`.

**AC1** `bash scripts/verify_go.sh --race-bounds-check; echo rc=$?` prints `race bounds: per-package 900s, start allowance 995s, leg 2000s` and `rc=0`.

**AC2** `VERIFY_GO_RACE_TIMEOUT_S=1954 bash scripts/verify_go.sh --race-bounds-check; echo rc=$?` prints `leg bound 1954s must be at least per-package 900s + 60s backup kill + start allowance 995s (1955s)` and `rc=2`. `VERIFY_GO_RACE_TIMEOUT_S=1955 bash scripts/verify_go.sh --race-bounds-check; echo rc=$?` prints the leg at 1955s and `rc=0`.

**AC3** Plain and 1777-parent contention replays. Expect seven `ok R<n>` lines, `PASS 7/7 arms`, rc 0, R6 under 30 s, and R7 under 10 s. Run the two legs together using §7 `accept.sh`, or use these commands separately. The spinners share the bounded child group and the exit trap handles normal cleanup.

````bash
state="$HOME/.ailang/state/world-iter443"
python3 "$state/bounded.py" 90 bash scripts/test_race_leg_bounds.sh
python3 "$state/bounded.py" 90 bash -c '
  set -e
  d=$(mktemp -d /private/tmp/r267.XXXX)
  pids=()
  cleanup() {
    if [ "${#pids[@]}" -gt 0 ]; then
      kill "${pids[@]}" 2>/dev/null || :
      wait "${pids[@]}" 2>/dev/null || :
    fi
    rm -rf "$d"
  }
  trap cleanup EXIT
  trap "exit 2" HUP INT TERM
  ls -ld /private/tmp
  for ((i=0; i<8; i++)); do yes >/dev/null & pids+=("$!"); done
  TMPDIR="$d" bash scripts/test_race_leg_bounds.sh
'
````

Both complete fixture runs returned **UNINFORMATIVE UNDER SANDBOX** for R6 with raw rc 1. R1–R5 passed and the wrapper returned 124 with the 2 s FATAL, but launching real `ps` returned `[Errno 1] Operation not permitted`. The contention run recorded `drwxrwxrwt 777 /private/tmp` and eight spinners. R7 executed independently with the same fixture body and passed in both settings (wrapper 4s, arm 4s plain and wrapper 4s, arm 4s contended). R6 wall/diagnostic acceptance remains pending outside. The receipts are `$S/ac3-final/legs.txt` and `$S/ac3-r7-final/legs.txt`.

**AC4** The selector uses exact top-level names and the pipeline preserves the Go exit code.

````bash
set -o pipefail
go vet ./host/verifygate/ &&
  go test ./host/verifygate/ -count=1 -v -run '^(TestRaceBoundsAreOrdered|TestRaceLegRunsInItsOwnJob|TestVerifyGoProductCheckSequence|TestVerifyGoRefusesUnknownLegs|TestEvidenceNamedManifestRejectsUnpinnedTest|TestZ3PinDeclaredOnceAndInstalledInBothJobs)$' 2>&1 |
  grep -E '^(--- |ok|FAIL)'
````

Expect six `--- PASS:` lines and an `ok` package summary. Base and stage receipts are in §1.3.

**AC5** Run the following from the applied tree. The helper creates scratch copies and replays the named killers. It compares actual test names and diagnostic text, restores each edit, and runs the same test again. During planning the Go-killed mutations also run the whole package. Expect all 21 design mutations and eight added diff cases to be killed in their named replay. Outside, use the same helper without `--named-only` to repeat the Go mutations whose full-package replays were incomplete.

````bash
S="$HOME/.ailang/state/world-iter443/plan-scratch"
out=$(mktemp -d "$S/ac5-outside.XXXXXX")
python3 "$S/replay-mutants.py" --outside --named-only --tree "$PWD" --jobs 3 --out "$out/replay" \
  M1 M2 M3 M4 M5 M6 M7 M8 M9 M10 M11 M12 M13 M14 M15 M16 M17 M18 M19 M20 M21 \
  X01 X02 X03 X04 X05 X06 X07 X08
awk -F ' [|] ' '$6 != "KILLED" && $0 != "ALL_DONE" {bad=1} END {exit bad}' "$out/replay/summary.txt"
````

**AC6** `bash "$S/gates.sh" "$(mktemp -d "$S/gates-outside.XXXXXX")"` runs the gates below and compares text fail sets to the controller's pristine run. The stand-alone design commands are:

````bash
./scripts/verify_ail.sh
go build ./... && go test ./... -count=1 -timeout 45m
gofmt -l host/verifygate/
bash -n scripts/verify_go.sh && bash -n scripts/test_race_leg_bounds.sh
````

Treat a socket-denied test as **UNINFORMATIVE UNDER SANDBOX**. The decision for each whole leg uses the controller's baseline fail set after its final marker. A new failure, a timeout panic, or a non-zero exit without an attributable fail set requires diagnosis. The pinned commit command in §8 requires a green whole verifygate run even if an earlier comparison found a carried failure. See §1.3.

**AC7 — PR CI** Set `RUN_ID` and the race job's `RACE_JOB_ID`, then run:

````bash
gh run view "$RUN_ID"
log=$(mktemp)
gh run view "$RUN_ID" --job "$RACE_JOB_ID" --log > "$log"
grep -F '── legs: race (race bounds: per-package 900s, leg 2000s)' "$log"
grep -F 'PASS 7/7 arms' "$log"
grep -F 'ok R6 hang' "$log"
grep -F '── go test ./... -count=1 -race -timeout 900s (leg bound 2000s)' "$log"
````

Confirm the three jobs are green. The R6 receipt here exercises Linux procps. **UNMEASURED** until the PR run.

**AC8 — first ten post-merge race attempts** Set the corresponding `RUN_ID` and `ATTEMPT` for each of the ten attempts. Save its log with:

````bash
log="$S/ac8-$RUN_ID-a$ATTEMPT.raw"
gh run view "$RUN_ID" --attempt "$ATTEMPT" --log > "$log"
grep -nF 'panic: test timed out after 15m0s' "$log"
grep -nE '(ok|FAIL)[[:space:]]+github.com/sunholo-data/ailang-world/(cmd/ailang-worldd|host/(daemon|verifygate|store))[[:space:]]' "$log"
````

Expect the panic grep to return 1. Record `attempt | heavy >480 s | result | cause class`. Each >480 s instance must be green or have a non-timeout failure that reproduces on the same commit, using (for example) its plain gate or the controller's base record for comparison. If all ten `host/store` readings are ≤67.5 s, record AC8 as unobserved. **UNMEASURED** before merge.

## 7. Executor scripts

`accept.sh` returns success only with rc 0 and `ACCEPT_OK` as its last line. `snapshot.sh M1` defaults to `$STATE/snap/M1`. `gates.sh` writes each raw leg and parser control before classifying the result. Its last line is `ALL_DONE`.

````bash accept.sh
#!/bin/bash
# AC1–AC5. Run from the current tree root.
set -euo pipefail
export AILANG_BIN="$HOME/.pinned-ailang/ailang" GOTOOLCHAIN=go1.26.6
state=$HOME/.ailang/state/world-iter443
d=$(mktemp -d /private/tmp/r267.XXXX)
pids=()
cleanup() {
  if [ "${#pids[@]}" -gt 0 ]; then
    kill "${pids[@]}" 2>/dev/null || :
    wait "${pids[@]}" 2>/dev/null || :
  fi
  rm -rf "$d"
}
trap cleanup EXIT
trap 'exit 2' INT TERM HUP
export AILANG_STATE_DIR="$d/state" AILANG_CACHE_DIR="$d/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"

expected='race bounds: per-package 900s, start allowance 995s, leg 2000s'
env -u VERIFY_GO_RACE_TIMEOUT_S bash scripts/verify_go.sh --race-bounds-check >"$d/ac1"
cat "$d/ac1"
[ "$(cat "$d/ac1")" = "$expected" ]
rc=0
VERIFY_GO_RACE_TIMEOUT_S=1954 bash scripts/verify_go.sh --race-bounds-check >"$d/ac2" 2>&1 || rc=$?
cat "$d/ac2"
[ "$rc" -eq 2 ]
grep -qF 'leg bound 1954s must be at least per-package 900s + 60s backup kill + start allowance 995s (1955s)' "$d/ac2"
VERIFY_GO_RACE_TIMEOUT_S=1955 bash scripts/verify_go.sh --race-bounds-check
python3 "$state/bounded.py" 90 bash scripts/test_race_leg_bounds.sh | tee "$d/ac3-plain"
grep -qxF 'PASS 7/7 arms' "$d/ac3-plain"
ls -ld /private/tmp
for ((i=0; i<8; i++)); do
  yes >/dev/null &
  pids+=("$!")
done
printf 'spinners=%s\n' "${#pids[@]}"
TMPDIR="$d" python3 "$state/bounded.py" 90 bash scripts/test_race_leg_bounds.sh | tee "$d/ac3-contended"
kill "${pids[@]}" 2>/dev/null || :
wait "${pids[@]}" 2>/dev/null || :
pids=()
grep -qxF 'PASS 7/7 arms' "$d/ac3-contended"
python3 "$state/bounded.py" 600 go vet ./host/verifygate/
names='TestRaceBoundsAreOrdered|TestRaceLegRunsInItsOwnJob|TestVerifyGoProductCheckSequence|TestVerifyGoRefusesUnknownLegs|TestEvidenceNamedManifestRejectsUnpinnedTest|TestZ3PinDeclaredOnceAndInstalledInBothJobs'
python3 "$state/bounded.py" 240 go test ./host/verifygate/ -count=1 -v -run "^($names)$" >"$d/ac4" 2>&1 ||
  { cat "$d/ac4"; exit 1; }
grep -E '^(--- |ok|FAIL)' "$d/ac4"
[ "$(grep -c '^--- PASS:' "$d/ac4")" -eq 6 ]
python3 "$state/plan-scratch/replay-mutants.py" --outside --named-only --tree "$PWD" --jobs 3 --out "$d/ac5" \
  M1 M2 M3 M4 M5 M6 M7 M8 M9 M10 M11 M12 M13 M14 M15 M16 M17 M18 M19 M20 M21 X01 X02 X03 X04 X05 X06 X07 X08
awk -F ' [|] ' '$6 != "KILLED" && $0 != "ALL_DONE" {bad=1} END {exit bad}' "$d/ac5/summary.txt"
echo ACCEPT_OK
````

````bash snapshot.sh
#!/bin/bash
# Run from the current tree root. Overrides are for the archive rehearsal.
set -euo pipefail
[ "$#" -eq 1 ] && [ "$1" = M1 ] || { echo 'usage: snapshot.sh M1' >&2; exit 2; }
base=6420426efcd0304c88adf347cf7af0f158d96863
dest=${SNAPSHOT_ROOT:-$HOME/.ailang/state/world-iter443/snap}
mkdir -p "$dest"
work=$(mktemp -d "$dest/.M1.XXXXXX")
trap 'rm -rf "$work"' EXIT
[ ! -e "$dest/$1" ] || { echo "snapshot already exists: $dest/$1" >&2; exit 2; }
mkdir "$work/base" "$work/result"
git -C "${SNAPSHOT_BASE_REPO:-$PWD}" archive "$base" -- scripts host .github | tar -x -C "$work/base"
python3 - "$work/base" "$work/result" <<'PY'
import filecmp
import os
from pathlib import Path
import shutil
import stat
import sys

base, result = map(Path, sys.argv[1:])

def files(root):
    return {p.relative_to(root) for scope in ("scripts", "host", ".github")
            for p in (root / scope).rglob("*") if p.is_file() or p.is_symlink()}

tracked = files(base)
live = files(Path("."))
status = []
deleted = []
for name in sorted(tracked | live):
    current = Path(name)
    original = base / name
    if name not in live:
        status.append(f"D\t{name}\n")
        deleted.append(f"{name}\n")
        continue
    if name in tracked:
        if current.is_symlink() != original.is_symlink():
            same = False
        elif current.is_symlink():
            same = os.readlink(current) == os.readlink(original)
        else:
            same = (filecmp.cmp(current, original, shallow=False)
                    and stat.S_IMODE(current.stat().st_mode) == stat.S_IMODE(original.stat().st_mode))
        if same:
            continue
    if current.is_symlink():
        raise SystemExit(f"snapshot cannot represent changed symlink: {name}")
    status.append(f"{'M' if name in tracked else 'A'}\t{name}\n")
    target = result / name
    target.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(current, target)
(result / "STATUS").write_text("".join(status))
(result / "DELETED").write_text("".join(deleted))
PY
(
  cd "$work/result"
  find . -type f | sed 's|^\./||' | LC_ALL=C sort |
    while IFS= read -r name; do
      [ "$name" = MANIFEST.sha256 ] && continue
      shasum -a 256 "$name"
    done > MANIFEST.sha256
  shasum -a 256 -c MANIFEST.sha256
)
mv "$work/result" "$dest/$1"
echo "SNAPSHOT_OK $dest/$1"
````

````bash gates.sh
#!/bin/bash
# Run from the tree root. Optional argument: a fresh log directory.
set -uo pipefail
state=$HOME/.ailang/state/world-iter443
base=6420426efcd0304c88adf347cf7af0f158d96863
logs=${1:-$(mktemp -d "${TMPDIR:-/tmp}/r267-gates.XXXXXX")}
mkdir -p "$logs" || exit 2
logs=$(cd "$logs" && pwd) || exit 2
export AILANG_BIN="$HOME/.pinned-ailang/ailang" GOTOOLCHAIN=go1.26.6
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export AILANG_STATE_DIR="$logs/state" AILANG_CACHE_DIR="$logs/cache"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" || exit 2
echo "Gate logs $logs"
fail=0
uninformative=0
run() {
  name=$1; shift
  limit=$1; shift
  rc=0
  python3 "$state/bounded.py" "$limit" "$@" >"$logs/$name.raw" 2>&1 || rc=$?
  printf 'LEG %s rc=%s\n' "$name" "$rc"
  printf '%s\n' "$rc" >"$logs/$name.rc"
}
gate() {
  run "$@"
  if [ "$rc" -ne 0 ]; then
    cat "$logs/$name.raw"
    fail=1
  fi
}
gofmt_check() {
  bad=$(gofmt -l "$@") || return 1
  if [ -n "$bad" ]; then
    printf '%s\n' "$bad"
    return 1
  fi
}
export -f gofmt_check
gate gofmt 60 bash -c 'gofmt_check "$@"' _ \
  host/verifygate/race_job_gate_test.go host/verifygate/mission_config_gate_test.go \
  host/verifygate/evidence_manifest_gate_test.go
gate build 600 go build ./...
gate vet 600 go vet ./host/verifygate/
gate linux-vet 600 env GOOS=linux go vet ./host/verifygate/
gate verify-go-syntax 60 bash -n scripts/verify_go.sh
gate race-test-syntax 60 bash -n scripts/test_race_leg_bounds.sh
gate ail 900 ./scripts/verify_ail.sh

compare() {
  run "$@"
  grep -nF 'panic: test timed out' "$logs/$name.raw" || :
  prc=0
  python3 "$state/failset.py" "$logs/$name.raw" "$rc" >"$logs/$name.fails" 2>"$logs/$name.parser" || prc=$?
  cat "$logs/$name.parser"
  if ! grep -qF 'FAILSET_CONTROL host/store TestParserPositive' "$logs/$name.parser"; then
    echo "INSTRUMENT FAILURE $name: missing fail-set control"
    fail=1
  elif grep -Ei 'listen (tcp|unix).*([Oo]peration not permitted|[Pp]ermission denied)|socket.*bind.*([Oo]peration not permitted|[Pp]ermission denied)' "$logs/$name.raw"; then
    echo "UNINFORMATIVE UNDER SANDBOX $name: socket bind denied"
    if [ "$prc" -ne 0 ] || { [ "$rc" -ne 0 ] && [ ! -s "$logs/$name.fails" ]; }; then
      echo "INSTRUMENT FAILURE $name: no complete attributable fail set"
    fi
    uninformative=1
  elif [ "$prc" -ne 0 ] || grep -qF 'panic: test timed out' "$logs/$name.raw"; then
    echo "INSTRUMENT FAILURE $name: inspect $logs/$name.raw"
    fail=1
  elif ! tail -n 1 "$state/baseline.done" 2>/dev/null | grep -qxE "ALL_DONE( for $base)?"; then
    printf 'LEG %s-comparison SKIPPED (controller baseline incomplete)\n' "$name"
    uninformative=1
  elif ! grep -qF 'FAILSET_CONTROL host/store TestParserPositive' "$state/base/$name.parser" ||
       grep -Eq 'INSTRUMENT FAILURE|Traceback' "$state/base/$name.parser" ||
       grep -qF 'panic: test timed out' "$state/base/$name.raw" ||
       ! test -f "$state/base/$name.fails"; then
    echo "INSTRUMENT FAILURE $name: controller base receipt invalid"
    fail=1
  else
    LC_ALL=C comm -13 "$state/base/$name.fails" "$logs/$name.fails" >"$logs/$name.new" || fail=1
    if [ -s "$logs/$name.new" ]; then
      echo "NEW FAILURES $name"
      cat "$logs/$name.new"
      fail=1
    else
      echo "FAIL SET $name has no new failures"
    fi
  fi
}
compare verifygate-plain 1900 go test -count=1 -timeout 30m ./host/verifygate/
compare pins-plain 1900 go test -count=1 -timeout 30m ./host/pins/
compare all-plain 2800 go test -count=1 -timeout 45m ./...

ac3() {
  run "$@"
  cat "$logs/$name.raw"
  if [ "$rc" -eq 0 ] && grep -qxF 'PASS 7/7 arms' "$logs/$name.raw"; then
    return
  fi
  if grep -Eq "Operation not permitted: 'ps'|Permission denied: 'ps'" "$logs/$name.raw"; then
    echo "UNINFORMATIVE UNDER SANDBOX $name: process listing denied"
    uninformative=1
  else
    fail=1
  fi
}
ac3 ac3-plain 90 bash scripts/test_race_leg_bounds.sh
ac3_contended() (
  set -e
  d=$(mktemp -d /private/tmp/r267.XXXX)
  pids=()
  cleanup() {
    if [ "${#pids[@]}" -gt 0 ]; then
      kill "${pids[@]}" 2>/dev/null || :
      wait "${pids[@]}" 2>/dev/null || :
    fi
    rm -rf "$d"
  }
  trap cleanup EXIT
  trap 'exit 2' INT TERM HUP
  ls -ld /private/tmp
  for ((i=0; i<8; i++)); do
    yes >/dev/null &
    pids+=("$!")
  done
  printf 'spinners=%s\n' "${#pids[@]}"
  TMPDIR="$d" bash scripts/test_race_leg_bounds.sh
)
export -f ac3_contended
ac3 ac3-contended 90 bash -c ac3_contended
echo ALL_DONE
[ "$fail" -eq 0 ] || exit 1
[ "$uninformative" -eq 0 ] || exit 2
````

`SNAPSHOT_ROOT="$S/snap-rehearsal" SNAPSHOT_BASE_REPO="$W" bash "$S/snapshot.sh" M1` created six STATUS entries, an empty DELETED and eight checked SHA-256 entries including the two bookkeeping files. The controller's active snapshot path is untouched.

Extracted the Appendix P fence from this plan with the controller's awk expression. `/usr/bin/patch -p1 -E` applied cleanly to the fresh archive at `$S/rt-1`. `diff -r` returned rc 0 for `scripts`, `host` and `.github` with no output. See `$S/roundtrip.raw`.

## 8. Commit-time test

Pass the following command to `commit_ms.sh M1 MSGFILE` as the test command. The controller's commit script separately verifies the manifest, gofmt, build and native vet. The pinned interpreter, tool release, node and SSH controls match the baseline lane. Use temporary interpreter state and cache.

````bash
export AILANG_BIN="$HOME/.pinned-ailang/ailang" GOTOOLCHAIN=go1.26.6
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
d=$(mktemp -d "${TMPDIR:-/tmp}/r267-commit.XXXXXX") || exit 2
trap 'rm -rf "$d"' EXIT
export AILANG_STATE_DIR="$d/state" AILANG_CACHE_DIR="$d/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" || exit 2
go test -count=1 -v -run '^(TestRaceBoundsAreOrdered|TestRaceLegRunsInItsOwnJob|TestVerifyGoProductCheckSequence|TestVerifyGoRefusesUnknownLegs|TestEvidenceNamedManifestRejectsUnpinnedTest|TestZ3PinDeclaredOnceAndInstalledInBothJobs)$' ./host/verifygate/ &&
  bash scripts/test_race_leg_bounds.sh &&
  go test -count=1 ./host/verifygate/
````

## Appendix P — compiled patches

````diff M1
diff --git a/.github/workflows/ci.yml b/.github/workflows/ci.yml
index 64ed2dd6..8b200090 100644
--- a/.github/workflows/ci.yml
+++ b/.github/workflows/ci.yml
@@ -518,9 +518,10 @@ jobs:
 
 
       - name: go -race gate (replay tests run against pinned AILANG_BIN)
-        # The leg bound (VERIFY_GO_RACE_TIMEOUT_S) is loud on overrun; this ceiling sits above it.
-        timeout-minutes: 35
+        # Package (900s) < leg (2000s) < step (40m). See row 267 in
+        # design_docs/planned/w-daemon-race-package-timeout-headroom.md §0.
+        timeout-minutes: 40
         env:
           VERIFY_GO_LEGS: race
-          VERIFY_GO_RACE_TIMEOUT_S: "1500"
+          VERIFY_GO_RACE_TIMEOUT_S: "2000"
         run: ./scripts/verify_go.sh
diff --git a/host/verifygate/evidence_manifest_gate_test.go b/host/verifygate/evidence_manifest_gate_test.go
index c58ea6b0..028d85c7 100644
--- a/host/verifygate/evidence_manifest_gate_test.go
+++ b/host/verifygate/evidence_manifest_gate_test.go
@@ -202,7 +202,7 @@ func runSyntheticEvidenceGate(t *testing.T, root string, tests []string, exact i
 func requireSyntheticEvidenceControl(t *testing.T, root string) {
 	t.Helper()
 	rc, out := runSyntheticEvidenceGate(t, root, syntheticEvidenceTests, len(syntheticEvidenceTests))
-	if rc != 0 || !strings.Contains(out, "all 64 required top-level evidence tests passed exactly once") {
+	if rc != 0 || out != "   ✓ all 64 required top-level evidence tests passed exactly once\n" {
 		t.Fatalf("pristine synthetic control did not reach success: rc=%d\n%s", rc, out)
 	}
 }
diff --git a/host/verifygate/mission_config_gate_test.go b/host/verifygate/mission_config_gate_test.go
index 63a3b25b..85453428 100644
--- a/host/verifygate/mission_config_gate_test.go
+++ b/host/verifygate/mission_config_gate_test.go
@@ -476,11 +476,29 @@ func requireVerifyGoProductSequence(t *testing.T, label, src string, repaired bo
 			`bash scripts/tracked_binary_gate.sh --allow-png 'website/static/img/v1/' \`,
 		}
 	}
+	raceAnchors := []string{`echo "── go test ./... -count=1 -race -timeout 8m (leg bound ${VERIFY_GO_RACE_TIMEOUT_S}s)"`, `cmd = ["go", "test", "./...", "-count=1", "-race", "-timeout", "8m"]`, `    print(f"verify_go.sh: FATAL: -race leg timed out after {bound}s", file=sys.stderr)`, `    sys.exit(124)`}
+	if repaired {
+		raceAnchors = []string{
+			`bash scripts/test_race_leg_bounds.sh`,
+			`echo "── go test ./... -count=1 -race -timeout ${VERIFY_GO_RACE_PKG_TIMEOUT_S}s (leg bound ${VERIFY_GO_RACE_TIMEOUT_S}s)"`,
+			raceLegCall,
+		}
+		for _, needle := range []string{`    print(f"verify_go.sh: FATAL: -race leg timed out after {bound}s", file=sys.stderr)`, `    sys.exit(124)`} {
+			if n := strings.Count(src, needle); n != 1 {
+				t.Errorf("%s: race wrapper anchor %q has count %d, want 1", label, needle, n)
+			}
+		}
+	}
 	passedAnchor := `echo "✓ go gate PASSED ($VERIFY_GO_LEGS): build clean; legs run pass with pinned AILANG_BIN ($ver)"`
 	needles := []string{
 		`if [ -z "${AILANG_BIN:-}" ]; then`,
 		verTokAnchor,
 		`echo "── AILANG_BIN=$AILANG_BIN ($ver)"`,
+	}
+	if repaired {
+		needles = append(needles, `check_race_bounds || exit $?`)
+	}
+	needles = append(needles,
 		`echo "── tracked-binary hygiene gate"`,
 		binaryAnchors[0],
 		binaryAnchors[1],
@@ -497,15 +515,12 @@ func requireVerifyGoProductSequence(t *testing.T, label, src string, repaired bo
 		evidenceCall,
 		`echo "── go test ./... -count=1"`,
 		`go test ./... -count=1`,
-		`echo "── go test ./... -count=1 -race -timeout 8m (leg bound ${VERIFY_GO_RACE_TIMEOUT_S}s)"`,
-		`cmd = ["go", "test", "./...", "-count=1", "-race", "-timeout", "8m"]`,
-		`    print(f"verify_go.sh: FATAL: -race leg timed out after {bound}s", file=sys.stderr)`,
-		`    sys.exit(124)`,
-		passedAnchor,
-	}
+	)
+	needles = append(needles, raceAnchors...)
+	needles = append(needles, passedAnchor)
 	position := -1
 	for _, needle := range needles {
-		exactLine := needle == `go build ./...` || needle == `go test ./... -count=1` || (repaired && needle == `check_mission_config`)
+		exactLine := needle == `go build ./...` || needle == `go test ./... -count=1` || (repaired && (needle == `check_mission_config` || needle == `check_race_bounds || exit $?`))
 		count := strings.Count(src, needle)
 		next := strings.Index(src, needle)
 		if exactLine {
diff --git a/host/verifygate/race_job_gate_test.go b/host/verifygate/race_job_gate_test.go
index af0c5028..a550d985 100644
--- a/host/verifygate/race_job_gate_test.go
+++ b/host/verifygate/race_job_gate_test.go
@@ -5,6 +5,7 @@ import (
 	"os/exec"
 	"path/filepath"
 	"regexp"
+	"strconv"
 	"strings"
 	"testing"
 )
@@ -13,23 +14,7 @@ import (
 // split from silently dropping -race: the race job must exist and run the race legs, the main
 // job must run the plain legs, and verify_go.sh must keep the race leg behind a loud bound.
 func TestRaceLegRunsInItsOwnJob(t *testing.T) {
-	raw, err := os.ReadFile(filepath.Join(repoRoot, ".github", "workflows", "ci.yml"))
-	if err != nil {
-		t.Fatal(err)
-	}
-	src := string(raw)
-	job := func(name string) string {
-		i := strings.Index(src, "\n  "+name+":\n")
-		if i < 0 {
-			t.Fatalf("ci.yml has no job %q", name)
-		}
-		rest := src[i+1:]
-		next := regexp.MustCompile(`\n  [a-z][a-z0-9-]*:\n`).FindStringIndex(rest[1:])
-		if next == nil {
-			return rest
-		}
-		return rest[:next[0]+1]
-	}
+	job := ciRaceJobReader(t)
 	plain, race := job("go-verify"), job("go-race-verify")
 	for _, c := range []struct{ body, want, why string }{
 		{plain, "VERIFY_GO_LEGS: plain", "the main Go job must skip only the race leg"},
@@ -53,7 +38,7 @@ func TestRaceLegRunsInItsOwnJob(t *testing.T) {
 	g := string(gate)
 	for _, want := range []string{
 		`if [ "$VERIFY_GO_LEGS" != plain ]; then`,
-		`"go", "test", "./...", "-count=1", "-race"`,
+		raceLegCall,
 		`FATAL: -race leg timed out after {bound}s`,
 		`FATAL: the race detector is not armed`,
 	} {
@@ -63,6 +48,90 @@ func TestRaceLegRunsInItsOwnJob(t *testing.T) {
 	}
 }
 
+const raceLegCall = `run_race_leg "$VERIFY_GO_RACE_TIMEOUT_S" go test ./... -count=1 -race -timeout "${VERIFY_GO_RACE_PKG_TIMEOUT_S}s"`
+
+func ciRaceJobReader(t *testing.T) func(string) string {
+	t.Helper()
+	raw, err := os.ReadFile(filepath.Join(repoRoot, ".github", "workflows", "ci.yml"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	src := string(raw)
+	return func(name string) string {
+		t.Helper()
+		i := strings.Index(src, "\n  "+name+":\n")
+		if i < 0 {
+			t.Fatalf("ci.yml has no job %q", name)
+		}
+		rest := src[i+1:]
+		next := regexp.MustCompile(`\n  [a-z][a-z0-9-]*:\n`).FindStringIndex(rest[1:])
+		if next == nil {
+			return rest
+		}
+		return rest[:next[0]+1]
+	}
+}
+
+func TestRaceBoundsAreOrdered(t *testing.T) {
+	raw, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "verify_go.sh"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	gate := string(raw)
+	race := ciRaceJobReader(t)("go-race-verify")
+	one := func(src, pattern string) []string {
+		t.Helper()
+		matches := regexp.MustCompile(pattern).FindAllStringSubmatch(src, -1)
+		if len(matches) != 1 {
+			t.Fatalf("race bounds: %q matches %d times, want 1", pattern, len(matches))
+		}
+		return matches[0]
+	}
+	seconds := func(src, pattern string) int {
+		t.Helper()
+		n, err := strconv.Atoi(one(src, pattern)[1])
+		if err != nil {
+			t.Fatalf("race bounds: %q: %v", pattern, err)
+		}
+		return n
+	}
+	p := seconds(gate, `(?m)^VERIFY_GO_RACE_PKG_TIMEOUT_S=(\d+)$`)
+	a := seconds(gate, `(?m)^VERIFY_GO_RACE_START_ALLOWANCE_S=(\d+)$`)
+	l := seconds(gate, `(?m)^VERIFY_GO_RACE_TIMEOUT_S="\$\{VERIFY_GO_RACE_TIMEOUT_S:-(\d+)\}"$`)
+	one(race, `(?m)^    runs-on: ubuntu-24\.04$`)
+	if strings.Contains(race, "continue-on-error") {
+		t.Error("race job must have zero continue-on-error directives")
+	}
+	// Match the step containing the gate invocation. Stop at the next step.
+	step := one(race, `(?m)^      - [^\n]+\n(?:^ {8}[^\n]*\n)*?^        run: \./scripts/verify_go\.sh$`)[0]
+	j := seconds(step, `(?m)^        timeout-minutes: (\d+)$`)
+	ciLeg := seconds(step, `(?m)^          VERIFY_GO_RACE_TIMEOUT_S: "(\d+)"$`)
+	call := one(gate, `(?m)^run_race_leg "\$VERIFY_GO_RACE_TIMEOUT_S" .*$`)[0]
+	if call != raceLegCall {
+		t.Errorf("race call must include -race and use the package bound with no -short, -run, -skip or -failfast=false: %s", call)
+	}
+
+	// Row 267 design §2 and §0. The class-invariant daemon extension was 166s.
+	const slowNeedS, daemonInflatedNeedS = 646, 828
+	const latestHeavyStartS, worlddExtensionS, daemonExtensionInflatedS = 554, 93, 347
+	const backupKillS, preLegMaxS, stepMarginS = 60, 55, 300
+	if p < (slowNeedS*5+3)/4 || p <= daemonInflatedNeedS {
+		t.Errorf("package bound %ds must be at least 1.25 × %ds and above inflated need %ds", p, slowNeedS, daemonInflatedNeedS)
+	}
+	if a < latestHeavyStartS+worlddExtensionS+daemonExtensionInflatedS {
+		t.Errorf("start allowance %ds must be at least %ds + %ds + %ds", a, latestHeavyStartS, worlddExtensionS, daemonExtensionInflatedS)
+	}
+	if ciLeg != l {
+		t.Errorf("CI leg bound %ds must equal script default %ds", ciLeg, l)
+	}
+	if l < p+backupKillS+a || ciLeg < p+backupKillS+a {
+		t.Errorf("leg bounds (default %ds, CI %ds) must be at least package %ds + backup %ds + start allowance %ds", l, ciLeg, p, backupKillS, a)
+	}
+	if j*60 < ciLeg+stepMarginS || stepMarginS < 5*preLegMaxS {
+		t.Errorf("step %dmin must cover leg %ds + %ds, and the margin must cover 5 × %ds", j, ciLeg, stepMarginS, preLegMaxS)
+	}
+}
+
 // TestVerifyGoRefusesUnknownLegs runs the real script: a bad VERIFY_GO_LEGS refuses (rc 2)
 // instead of silently running fewer legs. It needs the pinned AILANG_BIN the gate checks first.
 func TestVerifyGoRefusesUnknownLegs(t *testing.T) {
diff --git a/scripts/test_race_leg_bounds.sh b/scripts/test_race_leg_bounds.sh
new file mode 100644
index 00000000..d91a446b
--- /dev/null
+++ b/scripts/test_race_leg_bounds.sh
@@ -0,0 +1,110 @@
+#!/usr/bin/env bash
+set -euo pipefail
+
+tmp=$(mktemp -d "${TMPDIR:-/tmp}/race-leg.XXXXXX")
+cleanup() {
+  local pid
+  if [ -s "$tmp/group.pid" ]; then
+    pid=$(cat "$tmp/group.pid")
+    kill -KILL -- "-$pid" 2>/dev/null || true
+  fi
+  rm -rf "$tmp"
+}
+trap cleanup EXIT
+
+gate=(bash scripts/verify_go.sh)
+expected='race bounds: per-package 900s, start allowance 995s, leg 2000s'
+printf '%s\n' "$expected" >"$tmp/expected"
+order='must be at least per-package 900s + 60s backup kill + start allowance 995s (1955s)'
+fail() { echo "FAIL $arm: $* (rc=$rc)" >&2; cat "$tmp/out" "$tmp/err" >&2; exit 1; }
+run() {
+  rc=0
+  "$@" >"$tmp/out" 2>"$tmp/err" || rc=$?
+}
+ok() { echo "ok $arm"; }
+
+arm='R1 default'
+run env -u VERIFY_GO_RACE_TIMEOUT_S "${gate[@]}" --race-bounds-check
+[ "$rc" -eq 0 ] && cmp -s "$tmp/out" "$tmp/expected" || fail "$expected"
+ok
+arm='R2 boundary accept'
+run env VERIFY_GO_RACE_TIMEOUT_S=1955 "${gate[@]}" --race-bounds-check
+[ "$rc" -eq 0 ] || fail 'expected rc 0 at 1955s'
+ok
+arm='R3 boundary refuse'
+run env VERIFY_GO_RACE_TIMEOUT_S=1954 "${gate[@]}" --race-bounds-check
+[ "$rc" -eq 2 ] && grep -qF "$order" "$tmp/err" || fail 'expected rc 2 and the bound calculation'
+ok
+arm='R4 malformed'
+run env VERIFY_GO_RACE_TIMEOUT_S=15m "${gate[@]}" --race-bounds-check
+[ "$rc" -eq 2 ] && grep -qF 'must be a positive integer' "$tmp/err" || fail 'expected rc 2 and the integer diagnostic'
+ok
+arm='R5 passthrough'
+run "${gate[@]}" --race-leg-wrapper-check 30 sh -c 'exit 3'
+[ "$rc" -eq 3 ] && ! grep -qF FATAL "$tmp/err" || fail 'expected child rc 3 without FATAL'
+ok
+
+export RACE_PROBE_DIR="$tmp"
+cat > "$tmp/hang_probe.test" <<'PROBE'
+#!/bin/sh
+echo "$$" > "$RACE_PROBE_DIR/group.pid"
+sleep 300 &
+echo "$!" > "$RACE_PROBE_DIR/gc.pid"
+sleep 300
+PROBE
+chmod +x "$tmp/hang_probe.test"
+watchdog() {
+  python3 - "$@" <<'PYTHON'
+import os, signal, subprocess, sys
+bound = int(sys.argv[1])
+p = subprocess.Popen(sys.argv[2:], start_new_session=True)
+try:
+    sys.exit(p.wait(timeout=bound))
+except subprocess.TimeoutExpired:
+    os.killpg(p.pid, signal.SIGKILL)
+    p.wait()
+    sys.exit(f"watchdog fired after {bound}s")
+PYTHON
+}
+check_reaped() {
+  [ -s "$tmp/gc.pid" ] || fail 'grandchild PID control missing'
+  local pid i
+  pid=$(cat "$tmp/gc.pid")
+  for ((i=0; i<50; i++)); do
+    if ! kill -0 "$pid" 2>/dev/null; then
+      rm -f "$tmp/group.pid" "$tmp/gc.pid"
+      return
+    fi
+    sleep 0.1
+  done
+  fail "grandchild $pid still alive after reap poll"
+}
+check_timeout() {
+  [ "$rc" -eq 124 ] && grep -qF 'FATAL: -race leg timed out after 2s' "$tmp/err" || fail 'expected rc 124 and timeout diagnostic'
+}
+
+arm='R6 hang'
+SECONDS=0
+run watchdog 24 "${gate[@]}" --race-leg-wrapper-check 2 "$tmp/hang_probe.test"
+elapsed=$SECONDS
+check_timeout
+check_reaped
+grep -F 'still running at the leg bound' "$tmp/err" | grep -qF 'hang_probe.test' || fail 'expected a process listing naming the probe'
+if grep -qF -- '--race-leg-wrapper-check' "$tmp/err"; then
+  fail 'listing included the wrapper outside the probe group'
+fi
+[ "$elapsed" -ge 2 ] && [ "$SECONDS" -lt 30 ] || fail "wall bound: wrapper ${elapsed}s, arm ${SECONDS}s"
+echo "ok $arm (wrapper ${elapsed}s, arm ${SECONDS}s)"
+
+arm='R7 stalled ps'
+printf '#!/bin/sh\nexec sleep 300\n' > "$tmp/ps"
+chmod +x "$tmp/ps"
+SECONDS=0
+run watchdog 10 env PATH="$tmp:$PATH" "${gate[@]}" --race-leg-wrapper-check 2 "$tmp/hang_probe.test"
+elapsed=$SECONDS
+check_timeout
+grep -qF 'could not list' "$tmp/err" || fail 'expected a failed listing diagnostic'
+check_reaped
+[ "$elapsed" -ge 2 ] && [ "$SECONDS" -lt 10 ] || fail "wall bound: wrapper ${elapsed}s, arm ${SECONDS}s"
+echo "ok $arm (wrapper ${elapsed}s, arm ${SECONDS}s)"
+echo 'PASS 7/7 arms'
diff --git a/scripts/verify_go.sh b/scripts/verify_go.sh
index 145db6e4..4789e8cb 100755
--- a/scripts/verify_go.sh
+++ b/scripts/verify_go.sh
@@ -155,6 +155,56 @@ check_mission_config() {
   echo "   World decision-ledger validated; centralized runtime currency is NOT certified here."
 }
 
+# Row 267: design_docs/planned/w-daemon-race-package-timeout-headroom.md §0.
+# P covers the 828s inflated daemon estimate. A stacks the 554s latest heavy
+# start and 93s + 347s uncapped extensions. Allow Go's 60s backup kill before L.
+VERIFY_GO_RACE_PKG_TIMEOUT_S=900
+VERIFY_GO_RACE_START_ALLOWANCE_S=995
+VERIFY_GO_RACE_TIMEOUT_S="${VERIFY_GO_RACE_TIMEOUT_S:-2000}"
+
+check_race_bounds() {
+  local min
+  case "$VERIFY_GO_RACE_TIMEOUT_S" in
+    ''|*[!0-9]*) echo "verify_go.sh: FATAL: VERIFY_GO_RACE_TIMEOUT_S must be a positive integer" >&2; return 2 ;;
+  esac
+  min=$((VERIFY_GO_RACE_PKG_TIMEOUT_S + 60 + VERIFY_GO_RACE_START_ALLOWANCE_S))
+  if [ "$VERIFY_GO_RACE_TIMEOUT_S" -lt "$min" ]; then
+    echo "verify_go.sh: FATAL: leg bound ${VERIFY_GO_RACE_TIMEOUT_S}s must be at least per-package ${VERIFY_GO_RACE_PKG_TIMEOUT_S}s + 60s backup kill + start allowance ${VERIFY_GO_RACE_START_ALLOWANCE_S}s (${min}s)" >&2
+    return 2
+  fi
+}
+
+run_race_leg() {
+  python3 - "$@" <<'PYTHON'
+import os, signal, subprocess, sys
+from contextlib import suppress
+
+bound = int(sys.argv[1])
+p = subprocess.Popen(sys.argv[2:], start_new_session=True)
+try:
+    sys.exit(p.wait(timeout=bound))
+except subprocess.TimeoutExpired:
+    try:
+        try:
+            # Diagnostics have at most two seconds of grace before SIGKILL.
+            listing = subprocess.run(
+                ["ps", "-ww", "-eo", "pgid=,etime=,args="],
+                capture_output=True, text=True, errors="replace", check=True, timeout=2,
+            )
+            for line in listing.stdout.splitlines():
+                fields = line.split(None, 2)
+                if len(fields) == 3 and fields[0] == str(p.pid):
+                    print(f"verify_go.sh: still running at the leg bound (elapsed {fields[1]}): {fields[2]}", file=sys.stderr)
+        finally:
+            with suppress(ProcessLookupError):
+                os.killpg(p.pid, signal.SIGKILL)
+    except (OSError, subprocess.SubprocessError) as exc:
+        print(f"verify_go.sh: could not list race leg processes: {exc}", file=sys.stderr)
+    print(f"verify_go.sh: FATAL: -race leg timed out after {bound}s", file=sys.stderr)
+    sys.exit(124)
+PYTHON
+}
+
 if [ "${1:-}" = "--evidence-manifest-check" ]; then
   if [ "$#" -ne 3 ]; then
     echo "usage: $0 --evidence-manifest-check JSON EXACT" >&2
@@ -173,6 +223,18 @@ if [ "${1:-}" = "--mission-config-check" ]; then
   exit $?
 fi
 
+if [ "${1:-}" = "--race-bounds-check" ]; then
+  check_race_bounds || exit $?
+  echo "race bounds: per-package ${VERIFY_GO_RACE_PKG_TIMEOUT_S}s, start allowance ${VERIFY_GO_RACE_START_ALLOWANCE_S}s, leg ${VERIFY_GO_RACE_TIMEOUT_S}s"
+  exit 0
+fi
+
+if [ "${1:-}" = "--race-leg-wrapper-check" ]; then
+  shift
+  run_race_leg "$@"
+  exit $?
+fi
+
 if [ "$#" -gt 0 ]; then
   echo "verify_go.sh: unrecognized first argument: $1" >&2
   exit 2
@@ -209,9 +271,7 @@ if [ "$ver_tok" != 'v0.41.0' ]; then
 fi
 echo "── AILANG_BIN=$AILANG_BIN ($ver)"
 
-# Row 181 (attended 2026-10-16, Mark: "split race into its own CI job"): the -race leg needs ~1,000 s
-# wall on the 2-core runner (measured run 38884339127: broker 343 s, verifygate 308 s, daemon 307 s
-# under -race, ~2,050 s of package time at -p 2), so CI runs it as its own parallel job.
+# Row 181 (attended 2026-10-16): CI runs -race as its own parallel job.
 # VERIFY_GO_LEGS: all (default, local) | plain (everything except the -race leg) | race (the
 # shared preconditions, the race-detector control and the -race leg only).
 VERIFY_GO_LEGS="${VERIFY_GO_LEGS:-all}"
@@ -219,11 +279,8 @@ case "$VERIFY_GO_LEGS" in
   all|plain|race) ;;
   *) echo "verify_go.sh: FATAL: VERIFY_GO_LEGS must be all, plain or race (got: $VERIFY_GO_LEGS)" >&2; exit 2 ;;
 esac
-VERIFY_GO_RACE_TIMEOUT_S="${VERIFY_GO_RACE_TIMEOUT_S:-1500}"
-case "$VERIFY_GO_RACE_TIMEOUT_S" in
-  ''|*[!0-9]*) echo "verify_go.sh: FATAL: VERIFY_GO_RACE_TIMEOUT_S must be a positive integer" >&2; exit 2 ;;
-esac
-echo "── legs: $VERIFY_GO_LEGS (race bound ${VERIFY_GO_RACE_TIMEOUT_S}s)"
+check_race_bounds || exit $?
+echo "── legs: $VERIFY_GO_LEGS (race bounds: per-package ${VERIFY_GO_RACE_PKG_TIMEOUT_S}s, leg ${VERIFY_GO_RACE_TIMEOUT_S}s)"
 
 echo "── tracked-binary hygiene gate"
 # SM.A once committed a 15.7 MB Mach-O that every reviewer and both CI jobs passed, because
@@ -339,28 +396,13 @@ go version
 go test ./... -count=1
 fi
 
-# Measured at 78 s wall on darwin/arm64 at 7550ee9 under load ~4.5;
-# host/broker was the 76.9 s critical path. The doc's ~179 s was not reproduced.
-# 2026-10-03 (row 135 revert 78db602): on the 2-core runner this leg overran 600 s; with -p 2
-# packages run two at a time and host/verifygate (alphabetically last, 412 s race pre-#190) is
-# the tail. Fixed by removing redundant gate reruns inside verifygate (one shared, digest-checked
-# pristine control; one shared run per shim version line), not by raising this budget or
-# dropping the package from -race. Local GOMAXPROCS=2 race leg: 465 s before, 312 s after.
+# Row 267: 51 complete race attempts set the three bounds in
+# design_docs/planned/w-daemon-race-package-timeout-headroom.md §0.
 if [ "$VERIFY_GO_LEGS" != plain ]; then
-echo "── go test ./... -count=1 -race -timeout 8m (leg bound ${VERIFY_GO_RACE_TIMEOUT_S}s)"
+bash scripts/test_race_leg_bounds.sh
+echo "── go test ./... -count=1 -race -timeout ${VERIFY_GO_RACE_PKG_TIMEOUT_S}s (leg bound ${VERIFY_GO_RACE_TIMEOUT_S}s)"
 go version
-VERIFY_GO_RACE_TIMEOUT_S="$VERIFY_GO_RACE_TIMEOUT_S" python3 - <<'PY'
-import os, signal, subprocess, sys
-bound = int(os.environ["VERIFY_GO_RACE_TIMEOUT_S"])
-cmd = ["go", "test", "./...", "-count=1", "-race", "-timeout", "8m"]
-p = subprocess.Popen(cmd, start_new_session=True)
-try:
-    sys.exit(p.wait(timeout=bound))
-except subprocess.TimeoutExpired:
-    os.killpg(os.getpgid(p.pid), signal.SIGKILL)
-    print(f"verify_go.sh: FATAL: -race leg timed out after {bound}s", file=sys.stderr)
-    sys.exit(124)
-PY
+run_race_leg "$VERIFY_GO_RACE_TIMEOUT_S" go test ./... -count=1 -race -timeout "${VERIFY_GO_RACE_PKG_TIMEOUT_S}s"
 fi
 
 echo "✓ go gate PASSED ($VERIFY_GO_LEGS): build clean; legs run pass with pinned AILANG_BIN ($ver)"
````

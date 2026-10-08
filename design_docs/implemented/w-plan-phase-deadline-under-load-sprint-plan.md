# Sprint plan — `w-plan-phase-deadline-under-load` (row 153, clause-4; closes row 161)

**Design doc:** [`w-plan-phase-deadline-under-load.md`](w-plan-phase-deadline-under-load.md) (r3 carve-out, `e37807e`)
**Sprint base:** `990a6da` (dev) + the design commits = branch `sprint/row153-plan-phase-deadline-under-load` @ `e37807e`. `git diff --stat 990a6da e37807e` = the design doc only.
**Planner:** mission iteration 241, claude-opus-5-5 (unattended). Planner probes and logs: `~/.ailang/state/mission-world-iter241/planner-baseline/` and `…/planner-probe/`.
**Estimate:** ~1.75 d over 5 milestones and 7 commits in ONE PR, run as **two executor runs** (§6). **Risk:** medium-high. M3 adds a new on-disk cache, a new hard failure in publication, and a daemon background task.
**Scope:** host Go and tests only. One QUICKSTART paragraph (S7). No `.ail`, no pin, no `tools/launchd/*`.

---

## 0. The short version

The plan phase cold-compiles the transition and 7 `std/*` modules on every call, under a 2 s cap nobody derived. The work, in order:
- **M1** makes the timeout typed and diagnosable.
- **M2** derives the cap (4 s).
- **M3** builds a compile-cache template in the publication check and copies it into every capsule run.
- **M4** (executor part) records the proxy evidence.

The planner re-ran every gate on the pristine base: **all green, no pre-existing reds** (§2). The planner also re-verified the design's citations (§3) and measured five things the design assumed or missed (§4). Three of them change how tests must be written:
1. A stale template still gives byte-identical output, so **MUT-STALE-KEY cannot be killed by output equality**.
2. A warm run on a complete template changes no cache content (only `manifest.json`'s mtime), so AC3.2's "a run that would rewrite it" needs a deliberately incomplete template.
3. The hard publication failure reds **19 existing tests in 4 packages** that use fake interpreters.

The interpreter does give a deterministic warmth witness: a `MOD010` warning that appears on stderr only when it compiles.

---

## 1. Environment (on every command)

```sh
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang            # AILANG v0.41.0, commit 24ee108 (gate interpreter)
export WORLD_EXEC_SRT_NODE_MODULES=$HOME/.ailang/state/mission-world-iter241/srt/node_modules   # srt 0.0.78, cli.js = pin
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
# local-only supplement (CI has no tool binary, so these e2e SKIP there; memory: row 160):
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.52.1/ailang
```

`go version` = `go1.26.6 darwin/arm64` (= CI's `GOTOOLCHAIN`). **Instrument rules:**
- Redirect to a file, then `echo rc=$?` on the next line. Never `| tail; echo $?`, and never `PIPESTATUS` (zsh).
- A `-run` that matches nothing exits 0, and `go test` without `-v` hides SKIPs. Every named test is proven by `grep -- "^--- PASS: <name> "` (or `^    --- PASS:` for subtests) in the same call.

---

## 2. Baseline — measured by the planner on the pristine tree (`e37807e` = `990a6da` code)

The logs are in `~/.ailang/state/mission-world-iter241/planner-baseline/` (`run.sh`, `run2.sh`, `summary.txt`, `summary2.txt`, one `.log` per step). The host was otherwise idle (load average 2.4–5.8 on 16 cores).

### 2.1 The controller's gate list — ALL GREEN

| Gate | Command | rc | Wall | Notes |
|---|---|---|---|---|
| `.ail` | `./scripts/verify_ail.sh` | 0 | 24 s | 16 identities, 40 named tests, 444 se-tools named tests, 9/9 package steps, PUB011 ok |
| vet | `go vet ./...` | 0 | — | the compile fence for `_test.go` |
| build | `go build ./...` | 0 | 1 s | |
| plain | `go test ./... -count=1` | 0 | 169 s | 27/27 packages ok; `host/daemon` 168 s is the tail, `verifygate` 155 s, `broker` 118 s |
| race | `go test -race ./... -count=1 -timeout 8m` | 0 | 186 s | `daemon` 184 s, `verifygate` 175 s, `broker` 161 s, `coordinator` 52 s, `projection` 14 s |
| real srt | CI's 20-test list, `-p 1 -v` (below) | 0 | 58 s | **20/20 `--- PASS`**, including `TestExecSrtMCPEndToEnd`, `TestExecSrtMCPGrantAndBudget` and `TestExecSrtReplayIsByteEqual` |
| fork-lock | CI's verbose verifygate step | 0 | 3 s | |
| cleanup | CI's row-24 verbose step | 0 | 24 s | |
| gofmt | `gofmt -l host/ cmd/` | 0 | — | empty |

**Pre-existing reds: none. Pre-existing flakes observed: none** (on an idle 16-core host the plan takes ~0.2 s, design V7). Row 161's flake is a CI-only, load-only signature (design V23). A red at any later commit boundary is therefore attributable to that commit, **except** a `TestExecSrtMCP*` red in CI's parallel legs carrying row 153's signature. That is expected through M1/M2 (design §6) and is the controller's to count.

### 2.2 CI steps the controller's list does not cover (derived from `.github/workflows/ci.yml`)

| CI step | Local result at base | Executor obligation |
|---|---|---|
| `python3 scripts/check_world_contract_quality.py` (ratchet) | rc 0 (`--event workflow_dispatch`) | none: no `.ail` |
| `AILANG_BIN=… python3 scripts/test_kernel_contracts.py` | rc 0 | none: no `.ail` |
| Row 140 srt confinement matrix (`m0_matrix.py`) | linux-only, not run | none: srt is untouched |
| `scripts/verify_go.sh` extras: binary-blob guard (`git diff --numstat` `-`), toolchain deny-list/floor, race known-positive control, the `host/evidence` JSON leg | covered by the plain and race legs, except the blob guard | **commit no binaries.** Run `git diff --numstat 990a6da HEAD \| awk '$1=="-"'` and require empty output (the rebuilt `daemon.race.test` must stay outside the repo) |
| `design_docs/verification/w-race-gate-blindspot/run.sh` | not run (platform-gated toolchain probe) | none |
| `bench_worldd.sh --smoke` / `--check-claims` / recorder-refusal | smoke rc 0, claims rc 0 | run smoke after M3 (`ailang-worldd` start path changes) |
| gate-1 / census / census-live / gate-0 suites | all rc 0 | none (the plan and the verification dir do not touch the charter) |
| `check_no_personal_email.sh` + self-test | rc 0, rc 0 | **run on the M4 commit:** the verification dir is new loop-written text |
| floor harness (`FLOOR_REQUIRE_TOOL=1 python3 -m unittest …`) | rc 0, 155 tests, 5 skipped | none |
| **local supplement**: `WORLD_TOOL_AILANG_BIN=… go test ./host/daemon/ ./cmd/ailang-worldd/ -count=1 -v` | rc 0, **236 PASS, 3 SKIP** (`TestWalkRetimedAtScale`, `TestMeasureStartup`, `TestSessionPtyChildProcess`: pre-existing), 16 `TestSeTools*`/`TestWorkspaceRead*` PASS | **run after M3a and M3b.** These rigs publish through `PublishSet` and run real plans, but SKIP in CI |

### 2.3 Gate list per commit boundary (what "green" means in this plan)

```sh
./scripts/verify_ail.sh                                              >g1.log 2>&1; echo rc=$?
go vet ./...                                                         >g2.log 2>&1; echo rc=$?
go build ./... && go test ./... -count=1                             >g3.log 2>&1; echo rc=$?
go test -race ./... -count=1 -timeout 8m                             >g4.log 2>&1; echo rc=$?
tests="TestExecSrtGoTest TestExecSrtExitStatusAndArgv TestExecSrtRefusalsSpawnNothing TestExecSrtConfinementMatrix TestExecSrtGitHeadIsNotWritable TestExecSrtTimeoutKillsTheGroup TestExecSrtPycacheStaysOutOfTheWorktree TestExecSrtProfileEnvNeverReachesTheHostNode TestExecSrtStartupProbePasses TestExecSrtProbeRefusesAPackageOnlySandbox TestExecProbePassThroughShimFailsEveryConfinementArm TestExecProbeNoopShimFailsArm1 TestExecProbeZeroShimFailsArm5 TestExecProbeFailingToolchainFailsArm7 TestExecProbeMissingClientFailsTheRefusalArms TestQuickstartSection10ProfilesLoad TestExecStartupRefusalTable TestExecSrtMCPEndToEnd TestExecSrtMCPGrantAndBudget TestExecSrtReplayIsByteEqual"
go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -timeout 600s -v -run "^($(echo $tests | tr ' ' '|'))\$" >g5.log 2>&1; echo rc=$?
for t in $tests; do grep -q -- "^--- PASS: $t " g5.log || echo "NOPASS $t"; done          # must print nothing
gofmt -l host/ cmd/                                                  # must print nothing
# + the milestone's NEW tests: go test <pkg> -count=1 -v -run '^(…)$' >gN.log; grep "^--- PASS: <each> "
```

Record per-package wall for `coordinator`, `projection` and `daemon` in g3/g4 at every boundary. **Row 154** (CI race leg ≈ 540–556 s of a 600 s cap) makes any growth in the `daemon` tail a finding to report, even when green.

---

## 3. Design citations re-verified against the tree (the plan is an instrument)

Every `file:line` the design cites was re-read at `e37807e`. **TRUE as cited:**
- `effectful.go`: :37–41 (the budget constants), :111–118 (`handlerBudget`; doc comment :111, func :112), :127–139 (`runPhase`; the wrap at :135), :206 (finish `runPhase`), :248 (`dispatchEffectful`), :260 (plan), :279 (`hctx`), :290, :299–300 (`unrecorded`), :302–306 (finish → `return unrecorded(err)`), :387 (replay plan).
- `coordinator.go`: :271 (`inFlight.Delete`), :288 (`GetReceipt`), :386 (`AppendIntent`).
- `projection.go`: :372 (`dispatchError` call), :416 (func), :437 (`unrecorded`, first case), :461 (the generic deadline `case`; its `return` is :462).
- `capsule.go`: :194, :198, :202, :221, :228, :359–371.
- `check.go`: :61–102.
- `verify.go`: :86–98.
- `publish.go`: :70.
- `daemon.go`: :99, :704–705.
- `mcp_test.go`: :165.
- `workspace_exec_test.go`: :127.
- `m7b_deadline_test.go`: :12 (30 s).
- `handlers_test.go`: :250–263.
- `workspace_test.go`: :143–148 (`newWSDaemon`, ErrorLog defaults to `io.Discard` at :146–147).
- `exec_e2e_srt_test.go`: :126, :152–154, :197.
- `setools_e2e_test.go`: :809 (`TestSeToolsMCPPostEffectFailureIsLogged`).
- `w-software-engineering-domain.md`: :181, :248.
- `packages/se-tools/transitions.json`: N = 9, all effectful, 9 distinct files.
- Upstream: `serveapi/protocol/envelope.go:21–32` (`CallbackMessage`) and `serveapi/protocol/mcphttp/handler.go:118,125,182` at pin `v0.47.2`.

**Drift / defects (the executor must not "fix" the design; these are recorded here):**

| # | Design says | Measured | Consequence |
|---|---|---|---|
| D1 | `mcp.go:126–131` | Dispatch + `logRefusal` + return are **:130–133** | cosmetic |
| D2 | V35: `CheckSource`'s "only caller is `EnsureSourceLoadable` → `verifyLoadable`" | `EnsureSourceLoadable` is **also** called by `cmd/world-publish/transitions.go:314` (the CLI's pre-`PutObject` check) | the CLI also builds a template, and also meets `TemplateBuildError`. Retry lives inside `EnsureSourceLoadable` (P8), so both callers inherit it |
| D3 | AC2.3: the plan "blocks until `ctx.Done()`" then the handler deadline is checked | a plan that blocks until `pctx` is done **times out**, so no handler ever runs | P2 |
| D4 | §7: `TestA2AFinishOverrunStaysEffectsUnrecorded` kills MUT-ORDER-FINISH | in `dispatchEffectful` the finish `pctx` inherits `fctx`'s deadline (`detached(ctx, FinishPhaseBudget)` :302, then `WithTimeout(fctx, FinishPhaseBudget)` :128). Both expire together, so a finish overrun classifies as a caller deadline, not `*PhaseTimeoutError{finish}`. Even when it is one, `Phase == "plan"` excludes it. **MUT-ORDER-FINISH survives that e2e** | P3 |
| D5 | §7: `TestCapsuleWarmEqualsColdForSeTools` kills MUT-STALE-KEY | measured §4 P-M3: a stale template gives **byte-identical** output, and the interpreter silently recompiles | P5 |
| D6 | AC3.2: hash the template "before and after a run that would rewrite it" | measured §4 P-M2: a warm run on a complete template changes **no file content**, only `compile/manifest.json`'s mtime | P6 |
| D7 | AC3.1's hard failure is "new; trigger surface bounded" | measured §4 P-M5: **19 existing tests** red, because their fake interpreters' `check` writes no cache | P9 |
| D8 | (not mentioned) | `host/archive/check_test.go:53,60` assert `staged-as:entry.ail` | amend with the staging path (M3a) |
| D9 | AC3.6: every daemon e2e rig fails on any cold line | `publishReadTool` (`workspace_test.go:477`) writes the registry head with `CompareAndSetRegistryHead`, **bypassing `PublishSet`**, so `TestWorkspaceReadThroughProductionWiring` (:618) would run cold | P10 |
| D10 | AC3.7(iii): rebuild "after the listener is ready" | the rigs serve `httptest.NewServer(d.Handler())` and never call `Run`/`Listen`. Only `Run` (`daemon.go:1005`; `Listen` :1012, `Serve` :1065) has a listener-ready point | P12 |
| D11 | AC3.6: the label says `(template missing\|rebuilding)` | the capsule cannot know the daemon is rebuilding | P13 |

---

## 4. Planner measurements (new; probes under `~/.ailang/state/mission-world-iter241/planner-probe/`, pinned v0.41.0)

| # | Probe | Reading | Used by |
|---|---|---|---|
| P-M1 | `envkey.py`: template built by `check` at `host/capsule/main.ail` under the **production check env** (`childenv.Scrubbed(os.Environ())` incl. HOME, plus `AILANG_RELAX_MODULES=1`, `AILANG_CACHE_DIR`), then `run` under the capsule's **minimal** env (`AILANG_FS_SANDBOX`, `AILANG_RELAX_MODULES`, `AILANG_CACHE_DIR`) in a different root, 3× each arm | cold 171–181 ms with 41/41 cache files written. Template-from-full-env 27–28 ms, template-from-minimal-env 27–31 ms: the run changes 1 of 41 files. Stdout `4b6c4d08…` in all 9 | the env difference between check and run does not break the key; check-built templates warm runs in another root |
| P-M2 | template built, then a warm run with `AILANG_CACHE_DIR` = the template itself; snapshot (content sha256, then content + mtime) | **content changed: 0 files, added: 0.** content + mtime: only `compile/manifest.json`. After deleting `compile/modules/host__capsule__main/`, the run (rc 0, plan emitted) **adds 5 files** and changes `manifest.json` content | AC3.2 test design (P6) |
| P-M3 | template from `read.ail`, run of `exec.ail` with it (the key mutated to "interpreter only") | rc 0, stdout **identical** to a cold `exec.ail` run (`4b6c4d08…`); **10 cache files rewritten** | MUT-STALE-KEY is invisible to output equality (P5) |
| P-M4 | `exec.ail` (header `module se_tools/exec`) cold, warm, and from a check-built template; stderr | cold stderr: `WARNING MOD010 (relaxed): module 'se_tools/exec' does not match canonical path 'host/capsule/main'…`. **Warm and from-template stderr: empty.** `check` prints the same MOD010 warning | a **deterministic warmth witness** for every se-tools source (all are `module se_tools/*` ≠ the staging path). Also a v0.41.0 behaviour: add it to the AC3.5 canary |
| P-M5 | throwaway worktree (removed): `CheckSource` staged at `host/capsule/main.ail`, `AILANG_CACHE_DIR` set, and an error if `Passed` but no `compile/manifest.json`; `go test ./host/archive/ ./host/transitionreg/ ./host/daemon/ ./host/broker/ ./cmd/... -v` with the tool binary set → `planner-baseline/sim_m3_hardfail.log` | **19 FAIL**, all `TEMPLATE-SIM` (fake interpreters): archive 1, transitionreg 12, daemon 3, cmd/world-publish 3 (enumerated in §5 C2). **Positive control:** every real-interpreter rig PASSED with the hard failure on (30 `--- PASS` across `TestSeTools*`, `TestExecSrt*` and `TestWorkspaceRead*`), so the pinned interpreter really writes the manifest at the capsule path in production env. `host/broker` and `cmd/ailang-worldd` stayed ok | P9 |

---

## 5. Cross-package couplings (found by the planner; each is a must-change or a must-not-break)

| # | Coupling | Constraint |
|---|---|---|
| C1 | `host/broker/registry_publish_test.go:1181–1206` maps **subprocess-site files** to drivers and fails if `len(drivers) != len(files)` | **No `exec.Command*` in any new file.** The template build stays inside `host/archive/check.go` (`CheckSource`). The capsule launch stays in `host/capsule/capsule.go`. The daemon rebuild calls `archive.CheckSource` and never execs. `capsulecache.go`/`cachetemplate.go` do file I/O only |
| C2 | Fake interpreters whose `check` exits 0 (P-M5). The fakes: `host/archive/check_test.go:18` `checkingInterpreter`, `host/transitionreg/verify_test.go:22` `fakeInterpreterScript` (also used by `publish_test.go` and `epoch_test.go`), `host/daemon/registry_publisher_test.go:24` `daemonFakeInterpreter`, `cmd/world-publish/transitions_test.go:83` `fakeInterpreter`. The tests: archive `TestCheckSourceReportsTheArchivedInterpretersVerdict`; transitionreg `TestPublishSetEpochTwinInterpretersKeepDistinctPins`, `…EpochCheckIsPerDescriptorInterpreter`, `TestGenesisCannotUseBuildNextWithZeroExpectedHead`, `TestPublishSetGenesisThenGrowth`, `…IdempotentRepublishIsNoOp`, `…CASRetryMergesWinner`, `…SecondConflictSurfacesTypedError`, `…SameIDConflictOnCASRetryRefuses`, `…CASRetrySameIDIdenticalBytesIsNoOp`, `…RefusesEpochNotNominatingTheInterpreter`, `…RefusesUnloadableTransitionSource`, `TestEnsureSourceLoadableStagesHermetically`; daemon `TestPublishedTransitionsAppearOnAgentCard`, `TestEpochRefusalKeepsCardUnchanged`, `TestEpochTwinInterpretersBothListedOnCard`; world-publish `TestTransitionsVerbHappyPathAndIdempotence`, `TestTransitionsVerbFences`, `TestTransitionsReadsTheConfirmationFromThePtyWhenStdinIsDevNull` | in the **same commit** as the hard failure, each fake's `check` arm also writes `"$AILANG_CACHE_DIR/compile/manifest.json"` (one `mkdir -p` + `printf '{}' >`). It must write **only** when `AILANG_CACHE_DIR` is set, so a fake still models "check passed, cache absent" where a test wants that (AC3.7 v). The executor re-measures the red set at that commit; **any red outside this list is a finding: stop and report it** |
| C3 | `host/archive/check_test.go:53,60` assert `staged-as:entry.ail` | change to the new staging path (`host/capsule/main.ail`) |
| C4 | `entryModulePath = "host/capsule/main.ail"` is a `capsule.go:32` const. `capsule` imports `archive`, never the reverse | export the path from `archive` (e.g. `archive.CapsuleEntryPath`) and make `capsule`'s const equal to it, so check and run stage the same path by construction |
| C5 | `host/broker/allowlist_test.go`: `go list -deps` of broker/capsule against a module allowlist | stdlib only (`io/fs`, `path/filepath`, `sync/atomic`, `syscall`) |
| C6 | `TestNoParallelWriteForkPackagesOutsideVerifygate` (fork-lock gate) | **no `t.Parallel()`** in any new test that writes a script and execs it. No such test exists in these packages today (grep: 0) |
| C7 | Daemon background work and `-race`: tests in a package run sequentially, but a goroutine left by an earlier test races later ones | the template-maintenance goroutine is **cancelled and joined** by `Shutdown`/`Close` (and by `Run`'s exit). A test `ErrorLog` writer must tolerate writes after the test ends: `t.Log` after completion panics, so guard it with a mutex + `done` flag set in `t.Cleanup` |
| C8 | `TestSeToolsMCPPostEffectFailureIsLogged` (`setools_e2e_test.go:809`) pins `"host callback failed"` for R14 | M1 must not wrap `EffectsUnrecordedError` causes in anything carrying `DeadlineExceeded` |
| C9 | `projection_test.go:831` (row R11, bare `DeadlineExceeded` → `"invocation exceeded its deadline"`) and `a2a_wire_test.go:267,295,427` | must stay green unchanged (AC1.2) |
| C10 | S7 (usage surfaces) | `docs/QUICKSTART.md:91–96` documents the operator-line formats. Add the new A2A message (M1) and the two operator lines `capsule: cold compile …` / `capsule template <id>: …` plus the `<db>.artifacts/capsule-cache/` directory (M3) to that paragraph. Run every `QUICKSTART`-binding test (`host/daemon/*quickstart*`, `cmd/ailang-worldd/quickstart_*_test.go`, `host/broker/exec_quickstart_test.go`) |

---

## 6. Landing shape

**One PR, 7 commits, each green on §2.3 at its boundary. Two executor runs:**
- **Run 1:** C0 → C1 → C2 → C3 → C4 (M0, M1, M2, M3a).
- **Run 2:** C5 → C6 (M3b, M4-executor).

The split is where the risk changes. M3a ends with templates built at publication and consumed by every run, gated by the rig tripwire. Run 2 adds a daemon background task (startup status, rebuild after listener-ready, GC), which is the only part that touches daemon lifecycle and `Run`. It also adds the proxy re-run, which needs ~30 min of an otherwise idle host. Splitting keeps each executor context under one subsystem, and lets the controller check M3a's CI legs before Run 2 builds on them.

| Commit | Milestone | Content |
|---|---|---|
| C0 | M0 | none (no commit). The executor re-runs §2.3 on its HEAD and diffs the result against §2.1. A mismatch → STOP, report |
| C1 | M1 | coordinator `PhaseTimeoutError` + `classifyPhaseErr` + budget fields; projection case + prefix; coordinator/projection tests |
| C2 | M1 | daemon: rig `ErrorLog` writer, `TestExecSrtMCP*` amendments, `TestMCPPlanTimeoutFrozenWire`; QUICKSTART line for the A2A message |
| C3 | M2 | `PlanPhaseBudget = 4s` + derivation; `TestPlanPhaseBudgetIsDerived`, `TestSlowPlanLeavesHandlerItsCap` |
| C4 | M3a | archive template build in `CheckSource` + `TemplateBuildError` + promotion/racer + retry; fakes (C2/C3); capsule copy-in, cold label, `ColdRuns`, lazy promotion; daemon wires the capsule log; rig tripwire; `publishReadTool` fix; AC3.2–3.5 tests; QUICKSTART |
| C5 | M3b | daemon per-descriptor template status, rebuild after listener-ready, GC; AC3.7 + GC tests |
| C6 | M4 | `design_docs/verification/world-row153/` (README, baseline, per-boundary gate results, mutation ledger, the AC3.8/AC4.1 proxy V-row, archived scripts) |

C4 may be split into C4a (archive + transitionreg + fakes) and C4b (capsule + daemon wiring + tripwire), if both boundaries are green. This is recommended, because C4a alone already has a 19-test blast radius to clear.

---

## 7. Milestones

Conventions: "MUT" = apply exactly this edit on top of the milestone's own commit, run the named test(s) `-count=1 -v`, require `--- FAIL: <name>` with the stated reason, then `git checkout -- <file>` and require `git diff --quiet`. Record every MUT (diff line, the failing assertion's output line, the revert check) in the C6 mutation ledger. A MUT that does not red is a **finding**: report it, and never weaken the mutation to make it red.

### M0 — re-baseline (~0.05 d)

Run §2.3 on HEAD before any edit. Compare it with §2.1, including that all 20 srt names PASS and `gofmt` is clean. Copy the planner's two summaries into the C6 README later.

### M1 — typed phase timeout + diagnosable rig (~0.45 d) — AC1.1–AC1.7

**Files:** `host/coordinator/errors.go`, `host/coordinator/effectful.go`, `host/coordinator/coordinator.go` (`New`: budget fields), `host/projection/projection.go`; new tests in `host/coordinator/phase_timeout_test.go`, `host/projection/projection_test.go` (table), `host/projection/a2a_wire_test.go`, `host/daemon/exec_e2e_srt_test.go`, a new `host/daemon/testlog_test.go`, and `host/daemon/mcp_plan_timeout_test.go`. Also `docs/QUICKSTART.md`.

**Steps:**
1. `errors.go`: add `type PhaseTimeoutError struct{ Phase string; Budget, Elapsed time.Duration }`. Then:
   - `Error()` = `"coordinator: <phase> phase exceeded its <Budget> budget after <Elapsed rounded to 10ms>"`. This is what the operator line shows: `mcp refusal … plan phase exceeded its 2s budget after 2.01s`.
   - `Unwrap() error { return context.DeadlineExceeded }`.
2. `effectful.go`: add the pure `classifyPhaseErr(phase string, budget, elapsed time.Duration, parentErr, phaseErr error) error`. It returns:
   - `nil` when `phaseErr == nil` (the caller then falls to `classifyExec`);
   - `*PhaseTimeoutError` iff `phaseErr != nil && parentErr == nil`;
   - otherwise `fmt.Errorf("coordinator: execute %s phase: %w", phase, phaseErr)`. That is today's text: the caller's deadline or cancel, unchanged.
   
   In `runPhase` (:127), take `start := time.Now()` before `RunContext`. Replace :134–136 with `if cerr := classifyPhaseErr(in.Phase, budget, time.Since(start), ctx.Err(), pctx.Err()); cerr != nil { return nil, cerr }`. If the parent expired first, both errors are set and the result is the caller's deadline, which is AC1.2.
3. **Budget hook (P1):** add unexported `planBudget, finishBudget time.Duration` on `Coordinator`. `New` sets them to `PlanPhaseBudget`/`FinishPhaseBudget`. Replace the constant at :260 (dispatch plan), :387 (replay plan), :302 (`detached(ctx, …)`) and :206 (the finish `runPhase` inside `composeEffectOutput`, shared by dispatch and replay), so internal tests can scale them. No exported knob, no `Config` field.
4. `projection.go`: add the exported `const PhaseTimeoutPrefix = "transition plan phase exceeded its"` and a `var phase *coordinator.PhaseTimeoutError`. Put a new case **immediately after `:437` (`unrecorded`) and before every other case**: `case errors.As(err, &phase) && phase.Phase == "plan": return codeInternal, fmt.Sprintf("%s %s budget; nothing ran or was committed; resend the same task id", PhaseTimeoutPrefix, phase.Budget)`. Placing it directly after the `unrecorded` case also puts it before `:461`. No other case can match a `PhaseTimeoutError` (it wraps only `DeadlineExceeded`).
5. `docs/QUICKSTART.md:91–96`: one sentence naming the A2A plan-timeout message and that it is safe to resend the same task id.
6. **C2 (daemon):**
   - (a) `testlog_test.go`: `func testErrorLog(t *testing.T) io.Writer` writes each line with `t.Log`, is mutex-guarded, and drops writes after `t.Cleanup` (C7).
   - (b) `newExecSrtRig` (`exec_e2e_srt_test.go:29`) passes `ErrorLog: testErrorLog(t)` into the `mustWSDaemon` config.
   - (c) Add `func (r *seRig) callOK(token, name string, args map[string]any) (mcpWire, time.Duration)`. It calls `r.call` and `t.Fatalf("%s: wire error %+v", name, *wire.Error)` when `wire.Error != nil` (deref: never the pointer). Use it for every **success-expecting** call: `:122`, `:152`, `:195`, `:204` (P14).
   - (d) The no-grant call at `:186` keeps `r.call`, and its refusal message prints `*wire.Error` when non-nil (`:197`'s `%+v` of the struct, not the pointer).

**Tests and the mutation each must kill:**

| Test (pkg) | Assertions | MUT to apply (one line) → must red |
|---|---|---|
| `TestClassifyPhaseErrMatrix` (coordinator) | 4 rows, no sleeps. (nil,nil) → nil. (nil,DE) → `*PhaseTimeoutError{plan, budget, elapsed}` and `errors.Is(err, context.DeadlineExceeded)`. (DE,DE) → **not** `*PhaseTimeoutError`, `errors.Is` DE, text contains `execute plan phase`. (Canceled,Canceled) → not typed, `errors.Is` Canceled | **MUT-RELABEL:** in `classifyPhaseErr` change `parentErr == nil` → `true`, so row (DE,DE) reds. **MUT-NOWRAP:** delete `Unwrap()`, so row (nil,DE)'s `errors.Is` reds |
| `TestA2APlanTimeoutIsTypedAndNamesThePhase` (projection) | (i) table rows added to `TestA2ADispatch` (`projection_test.go:808`), each `fmt.Errorf("wrapped: %w")` as the table does: `&PhaseTimeoutError{Phase:"plan", Budget: coordinator.PlanPhaseBudget}` → `codeInternal`, starts with `PhaseTimeoutPrefix`, contains `plan`, `PlanPhaseBudget.String()` and `nothing ran`; **`&PhaseTimeoutError{Phase:"finish"}` bare → `"invocation exceeded its deadline"`** (P4). (ii) wire arm in `a2a_wire_test.go`: an effectful variant of `newWireRig` (descriptor `DeclaredEffects` = the token's effect; binder `broker.OpenBinder(st, ep, grants, broker.Registry{effect: okHandler})`; `invokeWait` 20 s) with a `hookRunner` whose plan blocks on `<-ctx.Done()`. `tasks/send` answers `-32603` with the prefix. Real budget wait: 2 s at M1, 4 s after M2 | **MUT-ORDER:** move the new case below `:461`'s generic case → (i) and (ii) answer `invocation exceeded its deadline`. **MUT-GENERIC:** message → `"transition timed out"` → the prefix/phase assertions red. **MUT-ANYPHASE:** delete `&& phase.Phase == "plan"` → the finish row reds |
| `TestA2AFinishOverrunStaysEffectsUnrecorded` (projection) | (i) table row: `&EffectsUnrecordedError{Cause: &PhaseTimeoutError{Phase:"plan"}}` → starts with `EffectsUnrecordedPrefix`. (ii) wire arm: the plan returns one effect with `finish: true`, and the finish blocks on `ctx.Done()` (real `FinishPhaseBudget`, 2 s) → message starts with `EffectsUnrecordedPrefix`, does **not** start with `PhaseTimeoutPrefix`, and ≠ `"invocation exceeded its deadline"` | **MUT-ORDER-FINISH** (move the new case above `:437`) → (i) reds; (ii) survives by construction (D4). **MUT-UNWRAP-FINISH:** `effectful.go:306` `return unrecorded(err)` → `return Result{}, err` → (ii) reds (generic deadline) |
| `TestA2APlanTimeoutRetrySameTaskID` (coordinator, internal) | `fxRig` + `probe` handler. Set `c.planBudget = 100*time.Millisecond`. A runner whose **first** plan blocks on `ctx.Done()` and later plans return `readPlan(false)`. Dispatch #1 (Surface A2A, task `t1`) → `errors.As` `*PhaseTimeoutError{Phase:"plan"}`, probe count 0, no log entry. Dispatch #2, same task id → success, probe count 1, exactly 1 log entry, 1 effect record (V28) | **MUT-EARLY-INTENT (proxy):** insert `c.uncertainIntent.Store(id, struct{}{})` as the first line of `dispatchEffectful` → Dispatch #2 answers `UnconfirmedError` → red. If the executor can express "a real `AppendIntent` before the plan" in ≤ 3 lines, apply that too and record both |
| `TestMCPPlanTimeoutFrozenWire` (daemon) | binary-free. `mustWSDaemon` with `fakeToolBin` + workspace root. `publishReadTool(t, d, <fake interp hash>, …)` (as `workspace_test.go:534` does). Replace coordinator + projection (the `serveR14` pattern, `setools_e2e_test.go:728`) with a blocking-plan runner and `ErrorLog: &syncLog{}`, `InvokeWait: invokeDeadline`. MCP `tools/call` → `{"code":-32603,"message":"host callback timed out"}` exactly, and the log has one line containing `plan phase exceeded its`. Comment: row 159 / ailang#1602 flips this | **MUT-NOWRAP** (delete `Unwrap`) → wire becomes `host callback failed` → red |
| `TestExecSrtMCPEndToEnd`, `TestExecSrtMCPGrantAndBudget` (amended, AC1.6) | `callOK` everywhere success is expected; rig `ErrorLog` = `testErrorLog` | **MUT-SILENT-ERROR (proof that the assertion is reached and prints the message):** change the second call's command at `:152` from `"where"` to `"no-such-command"` → the test must fail **at `callOK`** with a printed `{Code:… Message:…}`, not at `where = "<nil>"` |

**AC1.7 note (D4):** the e2e arm's purpose is the end-to-end wrap; the synthetic row's purpose is the case order. Both stay.

### M2 — the derived plan cap (~0.15 d) — AC2.1–AC2.4

**Files:** `host/coordinator/effectful.go`; new `host/daemon/plan_budget_test.go`; new test in `host/coordinator/phase_timeout_test.go`.

**Steps:**
1. `effectful.go:37`: set `PlanPhaseBudget = 4 * time.Second`. Comment: the largest plan cap that leaves the handler its full cap is `invokeDeadline (20 s, daemon) − HandlerCap − HandlerHeadroom = 20 − 10 − (2 + 4)`. Pinned by `TestPlanPhaseBudgetIsDerived`.
2. `FinishPhaseBudget` stays 2 s, with a one-line comment citing AC2.4 / R2.

| Test (pkg) | Assertions | MUT → must red |
|---|---|---|
| `TestPlanPhaseBudgetIsDerived` (daemon) | `coordinator.PlanPhaseBudget == invokeDeadline − coordinator.HandlerCap − coordinator.HandlerHeadroom` | **MUT-GUESS:** `PlanPhaseBudget = 3 * time.Second`. Also apply `HandlerCap = 9 * time.Second` alone (this also reds `TestExecMaxTimeoutIsHandlerCapLessOneSecond`; record both) |
| `TestSlowPlanLeavesHandlerItsCap` (coordinator, internal, **default budgets, no hook**: this also pins `New`'s default) | **P2.** Caller ctx = `WithTimeout(20 s)` (literal `modelInvokeDeadline`; comment that `TestPlanPhaseBudgetIsDerived` and `mcp_test.go:165` pin it). The plan runner reads `pctx.Deadline()`, waits until `deadline − 500ms` on a timer (select with `ctx.Done()`), records `planReturned := time.Now()`, and returns `readPlan(false)`. The probe handler records `dl, _ := ctx.Deadline()`. Assert `dl.Sub(planReturned) ≥ HandlerCap − 100ms`. Wall ≈ 3.5 s | **MUT-STEAL:** `PlanPhaseBudget = 5 * time.Second` → `dl − planReturned` = 20 − 4.5 − 6 = 9.5 s < 9.9 s → red (`TestPlanPhaseBudgetIsDerived` reds too). Also **MUT-DEFAULT:** `New` sets `planBudget: 2 * PlanPhaseBudget` → red |

### M3a — templates at publication, copied into every run (~0.55 d) — AC3.1 (build, error, racer, retry), AC3.2–AC3.6

**Files:**
- `host/archive/check.go`; new `host/archive/capsulecache.go` (file I/O only, C1).
- `host/transitionreg/verify.go`.
- `host/capsule/capsule.go`; new `host/capsule/cachetemplate.go` (no exec, C1).
- `host/daemon/daemon.go` (:660, pass the log).
- The tests and fakes in C2/C3.
- `host/daemon/workspace_test.go` (`newWSDaemon` default writer, `publishReadTool`).
- `docs/QUICKSTART.md`.

**Steps (function level; shapes are the executor's, invariants are not):**
1. `capsulecache.go` (archive):
   - `CapsuleEntryPath` (C4).
   - `const capsuleCacheDir = "capsule-cache"`.
   - `const CacheDirEnv = "AILANG_CACHE_DIR"`.
   - `(a *Archive) CapsuleTemplateDir(interp hashref.HashRef, source []byte) string` = `<a.Root()>/capsule-cache/<interp.Digest()>/<hex sha256(source)>`.
   - `(a *Archive) CapsuleTemplateReady(interp, source) bool`: `compile/manifest.json` is a regular file.
   - `newTemplateTmp()`: `MkdirAll(capsule-cache, 0o700)` + `MkdirTemp(capsule-cache, ".tmp-*")`.
   - `PromoteTemplate(tmp, interp, source) error`: `MkdirAll(parent, 0o700)`, then `renameDir(tmp, final)`. A racer (`errors.Is(err, fs.ErrExist)` or `syscall.ENOTEMPTY`/`EEXIST`) → `RemoveAll(tmp)`, `nil`. Any other error → `RemoveAll(tmp)`, `&TemplateBuildError{…, Err: err}`.
   - `PromoteCopy(runCacheDir, interp, source) error`: deep-copy regular files and dirs only (refuse symlinks/devices) into a fresh `newTemplateTmp()` in `capsule-cache/`, then `PromoteTemplate`. It renames only within `capsule-cache/`, so it is EXDEV-safe.
   - `var renameDir = os.Rename` (seam).
   - `type TemplateBuildError struct{ Interpreter, Source hashref.HashRef; Reason string; Err error }` with `Error`/`Unwrap`.
2. `check.go` `CheckSource`:
   - Stage at `CapsuleEntryPath` (MkdirAll its dir). Keep `AILANG_RELAX_MODULES=1`, and update the :79–82 comment.
   - Append `CacheDirEnv + "=" + tmp` **last** in `cmd.Env` (last-wins over an ambient value).
   - On `!Passed`: `RemoveAll(tmp)`, unchanged verdict.
   - On `Passed` with no `compile/manifest.json`: `RemoveAll(tmp)`, return `CheckResult{Output, Passed: true}, &TemplateBuildError{Reason: "check passed but wrote no compile/manifest.json"}`.
   - On `Passed` with a manifest: `PromoteTemplate`; its error is returned the same way.
   - `err` stays "host-side infrastructure", which matches the doc comment at :51–60. Update that comment.
3. `verify.go` `EnsureSourceLoadable`: retry while `errors.As(err, &*archive.TemplateBuildError)`, up to `const templateBuildAttempts = 3` total. Return the last error wrapped (`%w`) after that. Placing it here covers both callers (D2, P8).
4. Fakes (C2) and `check_test.go` (C3), in the same commit as steps 2–3.
5. `capsule.go` / `cachetemplate.go`:
   - `Config.Log io.Writer` (nil → `io.Discard`).
   - `Runner` gets `log io.Writer` and `cold atomic.Int64`, plus `func (r *Runner) ColdRuns() int64`.
   - In `RunContext` after staging: `cacheDir := filepath.Join(root, ".capsule-cache")`. If `CapsuleTemplateReady` → `copyTree(templateDir, cacheDir)` (copy, never link or rename). Else write one line `capsule: cold compile <interp-digest>/<source-sha256> (template missing)\n` to `r.log` and `r.cold.Add(1)`.
   - Append `archive.CacheDirEnv + "=" + cacheDir` to `cmd.Env` (`:228`).
   - After `collectOutput`, **only** if the run was cold and `err == nil && runErr == nil`: call `r.archive.PromoteCopy(cacheDir, …)`. A promotion error is logged, never returned: the run already succeeded.
   - The single `exec.CommandContext` stays at `capsule.go:223` (C1).
6. `daemon.go:660`: `capsule.New(a, capsule.Config{Log: d.errLog})`. Keep the `*capsule.Runner` on `Daemon` (`d.capsule`) for M3b.
7. Rig tripwire (AC3.6):
   - Extend `testErrorLog` into `testOperatorLog(t)`, which also calls `t.Errorf` (not `Fatal`, since it may run off the test goroutine) on any line containing `capsule: cold compile`.
   - `newWSDaemon` (`workspace_test.go:143`) installs it **when `cfg.ErrorLog == nil`**, replacing the `io.Discard` default. Tests that pass their own log keep it.
   - **P10:** `publishReadTool` builds the template when `interp` resolves in `archive.New(db)`: call `archive.New(f.db).CheckSource(ctx, interp, source)` after `PutObject`. The `:534` caller pins a fake hash that never runs, so it is skipped.
8. QUICKSTART (C10).

| Test (pkg) | Assertions | MUT → must red |
|---|---|---|
| `TestPublishBuildsCapsuleTemplates` (transitionreg; real `AILANG_BIN`, `t.Fatal` if unset, never skip) | `PublishSet` of two real sources → for each, `CapsuleTemplateReady` and no `.tmp-*` left in `capsule-cache/` | **MUT-PUBLISH-NOTEMPLATE:** delete the `CacheDirEnv` append in `CheckSource` → `PublishSet` returns `TemplateBuildError` (the hard failure) → red. Also apply "delete the env append **and** the manifest check" → template absent → red |
| `TestTemplateRenameRacerIsSuccess` (archive) | pre-create the final dir with a manifest (the "winner"). `PromoteTemplate(tmp2)` → `nil`, `tmp2` removed, the winner's bytes unchanged. Same via `PromoteCopy` | **MUT-RACER-ERROR:** the racer predicate → `false` → returns `TemplateBuildError` → red |
| `TestPublishRetriesTemplateBuildError` (transitionreg) | a fake interpreter whose `check` writes the manifest **only from its 2nd invocation** (counter file in its dir) → `PublishSet` succeeds, the template is ready, and the fake ran `check` exactly 2 times. Control arm: the fake never writes it → `errors.As` `*archive.TemplateBuildError` after exactly 3 invocations | **MUT-WEDGE:** `templateBuildAttempts = 1` → the first arm reds |
| `TestCapsuleTemplateIsCopiedNotShared` (capsule; real binary) | **P6.** Build the template via `CheckSource`. Delete `compile/modules/host__capsule__main/` from it (P-M2: a run then adds 5 files). Digest the tree with (relpath, size, sha256, mtime_ns). Run → rc ok, plan emitted. The digest is unchanged | **MUT-SHARE:** set `AILANG_CACHE_DIR` to the template dir itself instead of `cacheDir` → 5 files added + `manifest.json` rewritten → red |
| `TestCapsulePromotionSurvivesCrossDevice` (**archive**, P7: the `renameDir` seam is package-internal) | `PromoteCopy` from a run dir under a separate `t.TempDir()`, with `renameDir` stubbed to return `syscall.EXDEV` whenever `filepath.Dir(old) != filepath.Dir(new)` → `nil`, template ready, no `.tmp-*` left | **MUT-CROSS-DEVICE-RENAME:** `PromoteCopy` calls `renameDir(runCacheDir, final)` directly → `EXDEV` → red |
| `TestCapsuleColdRunIsLabelled` (capsule; **planner addition**) | no template: one run → exactly one log line matching `^capsule: cold compile [0-9a-f]+/[0-9a-f]{64} \(template missing\)$`, `ColdRuns() == 1`. Then the template is ready (lazy promotion), and a second run logs nothing, `ColdRuns() == 1` | **MUT-SILENT-COLD (run side):** delete the `Fprintf` → red. **MUT-NO-PROMOTE:** delete the `PromoteCopy` call → the second-run assertions red |
| `TestCapsuleWarmEqualsColdForSeTools` (**coordinator**, P5: the fixtures `seToolsCases` :163 and the finish inputs of `TestSeToolsFinishPhases` :327 are package-internal) | `newSeToolsRig` keeps its archive. For each case's plan input (and each finish input): **cold arm** = delete `CapsuleTemplateDir` first, then run, assert stderr contains `MOD010` and `ColdRuns` +1. **Warm arm** = build via `CheckSource`, then run, assert stderr has **no** `MOD010` and `ColdRuns` unchanged. stdout warm == cold byte for byte. **Cross-source arm:** build templates for `ailang-read` and `workspace-exec`, run `workspace-exec`, assert no `MOD010` | **MUT-STALE-KEY:** `CapsuleTemplateDir` drops the source component (key = interpreter only) → the cross-source arm and the warm arms see `MOD010` (P-M3: silent recompile) → red. Output equality alone would **not** red (P-M3) |
| `TestCapsuleWarmRunsDoNotCompile` (capsule; real binary; `exec.ail` plan) | 5 cold + 5 warm, **interleaved C,W,C,W…** so load bursts hit both arms. Assert median warm < median cold / 2 **and** every warm stderr lacks `MOD010` (the deterministic tooth) | **MUT-NOTEMPLATE:** skip the `copyTree` call (the env still points at an empty dir) → `MOD010` on warm stderr + ratio ≈ 1 → red | r2: ratio assertion removed — machine-dependent (judge B1, CI runs 37696379845 / 37700015498)
| `TestPinnedInterpreterHonoursCacheDir` (capsule; canary for S8) | the pinned binary's `check` with `archive.CacheDirEnv` → `compile/manifest.json` exists. A cold `run` prints `MOD010` for a `module se_tools/*` header and a warm one does not (P-M4). The comment lists it for the S8 floor-raise inventory | **floor drift emulation:** `CacheDirEnv = "AILANG_CACHE_DIRX"` → the manifest is absent → red |
| rig tripwire (AC3.6; daemon) | `testOperatorLog` installed by `newWSDaemon` | **proof:** in `newExecSrtRig`, after `publishSeTools`, `os.RemoveAll` the `workspace-exec` template → `TestExecSrtMCPEndToEnd` must fail with the tripwire's `cold compile` error (revert) |

**After C4:** re-run P-M5's package set with the tool binary (§2.2 local supplement). Expect 0 FAIL and all 16 `TestSeTools*`/`TestWorkspaceRead*` PASS **with no tripwire line**. Also run `./scripts/bench_worldd.sh --smoke`.

### M3b — daemon template status, rebuild after listener-ready, GC (~0.4 d) — AC3.1 (GC), AC3.7

**Files:** `host/daemon/daemon.go` (+ optionally a new `host/daemon/templates.go`, with no exec, C1); tests `host/daemon/templates_test.go`; `docs/QUICKSTART.md`.

**Steps:**
1. `type TemplateStatus struct{ ID string; Interpreter, Source hashref.HashRef; State string }`, where `State` is `ready|missing|rebuilt|failed`. `Daemon.templates` sits behind a mutex. `func (d *Daemon) CapsuleTemplates() []TemplateStatus` returns a copy.
2. In `New`, after `bootstrapRegistry` and only when `d.coord != nil`: read the registry head (`transitionreg.NewReader(d.store).CurrentRevision`). For each **effectful** descriptor, stat `CapsuleTemplateReady` (source bytes from the store) and write `ailang-worldd: capsule template <id>: ready|missing` to `d.errLog`. This is a stat per descriptor and never blocks (P12).
3. `func (d *Daemon) startTemplateMaintenance(ctx context.Context)` starts one goroutine, tracked by a `sync.WaitGroup` and cancelled by a stored `cancel`:
   - (a) GC: prune `capsule-cache/<digest>` dirs whose digest is **neither** `d.interpreterRef`'s digest **nor** any head descriptor's interpreter digest (P11). Also prune `.tmp-*` entries older than 2 × the check bound.
   - (b) Rebuild every `missing` entry with `archive.CheckSource` (its own 10 s bound), at concurrency ≤ 2 (semaphore). Write `capsule template <id>: rebuilt in <d>` or `… rebuild failed: <err>`, and update the status.
   - An unexported `Config` test seam `templateRebuildGate <-chan struct{}` is awaited before each rebuild. Production leaves it nil.
4. `Run` (`daemon.go:1005`) calls `d.startTemplateMaintenance(ctx)` **after** `Listen` (:1012) and the announce write, before `go d.Serve()` (:1065). On exit, after drain, cancel and wait (bounded). `Close` also cancels and waits, which makes it idempotent.
5. QUICKSTART: the startup and rebuild lines.

| Test (pkg) | Assertions | MUT → must red |
|---|---|---|
| `TestDaemonStartReportsMissingTemplates` (daemon; needs only `AILANG_BIN` + `fakeToolBin`, so it **must run in CI**, not skip). Fixture: daemon #1 (New) + `seGenesis` + `publishSeTools` + mint a token (`Workspace.Read` on ep1) → Close. Delete `ailang-read`'s template. Daemon #2 via **`Run`** with `cfg.ErrorLog = &syncLog{}` and `templateRebuildGate` held closed | (i) exactly one `capsule template <id>:` line per effectful descriptor (9), with `ailang-read` `missing` and the rest `ready`, present **before** the announce line is consumed. (iv) While the gate is held: MCP `tools/call ailang-read` → the log gains exactly one `capsule: cold compile` line (the plan ran cold; the effect outcome is irrelevant). (iii) Then release the gate → a `capsule template ailang-read: rebuilt in` line follows, and the template is ready. The first request was served **before** the rebuild ended (by construction: the gate). (ii) A separate New-based arm: `CapsuleTemplates()` reports the same states as (i). (v) A stub interpreter whose `check` passes and writes no cache → `PublishSet` → `errors.As` `*archive.TemplateBuildError` | **MUT-SILENT-COLD (startup):** delete the status `Fprintf` → (i) red. **MUT-BLOCKING-PREWARM:** call the rebuild synchronously in `New` (before `Listen`) → `New` blocks on the held gate → `Run` never announces within the test's bound → red. The test releases the gate in `t.Cleanup`, so the mutant does not leak a blocked goroutine into later tests. **MUT-SOFT-PUBLISH:** `CheckSource` ignores a missing manifest → (v) red |
| `TestTemplateGCPrunesStaleInterpDigests` (daemon) | create `capsule-cache/<64×"0">/x/compile/manifest.json`, a stale `.tmp-old` (mtime set back) and a fresh `.tmp-new`. Start maintenance (direct call is fine: `Run` ordering is pinned above), then wait on its WaitGroup → the stale digest and `.tmp-old` are gone; the current digest's templates and `.tmp-new` remain | **MUT-NO-GC:** make the prune loop `continue` unconditionally → red |

**After C5:** re-run the local supplement and the §2.3 list. `TestDaemonStartReportsMissingTemplates` must show `--- PASS` in the plain and race legs.

### M4 (executor part) — evidence (~0.15 d + ~30 min idle-host proxy) — AC4.1 (+ AC3.8)

**Files:** `design_docs/verification/world-row153/` containing:
- `README.md`;
- `baseline.md` (this plan's §2 + the M0 re-run);
- `gates.md` (per commit: rc, wall, per-package wall for `coordinator`/`projection`/`daemon`);
- `mutations.md` (every MUT: applied line, failing assertion line, revert check);
- `proxy/` (scripts + trimmed logs).

**Proxy procedure (AC3.8 = AC4.1):**
- Copy `bg_load.sh` into `proxy/bg_load.sh`, **parameterised**: the test-binary path and the package dir as arguments. The original hardcodes the design worktree and a base-built binary.
- Re-apply the probe patch (`~/.ailang/state/mission-world-iter241/probe-phase-timing.patch`) to a **scratch copy** of the M3b head. It touches `capsule.go`/`effectful.go`, which M1–M3 changed, so expect manual re-application and keep the result **uncommitted**.
- `go test -race -c -o $S/daemon.race.m3.test ./host/daemon/`.
- Run `proxy/bg_load.sh 8 4` (on an otherwise idle host). Record, per test:
  - (a) startup probe PASS;
  - (b) per-descriptor template status after publish (all `ready`) + publish wall vs V34;
  - (c) the first plan's `child` ≤ 1 s vs V15's 1.96–2.98 s;
  - (d) time-to-ready vs base.
- Require **0 plan timeouts** and **0 `cold compile` lines**.
- Control: V13 on the base tree, 4/4 timeouts (re-run it with the base-built `daemon.race.test` only if it is still present; otherwise cite V13).
- 12-burner arm: failures only at startup probe/publication, never at a plan or a template.

Logs are trimmed and **contain no personal addresses**. Run `bash scripts/check_no_personal_email.sh` on C6.

**Not executor work:** AC4.2 (5/5 consecutive CI attempts, counting every attempt) and AC4.3 (charter rows 153/161, follow-up rows) belong to the controller.

---

## 8. Planner decisions (the executor applies them; it never resolves these differently without reporting)

| # | Ambiguity | Decision and reason |
|---|---|---|
| P1 | "test-only budget hook" with no mechanism; `projection.Handler.coord` is a concrete `*coordinator.Coordinator` (`projection.go:171`), so no fake coordinator exists | Unexported `planBudget`/`finishBudget` fields on `Coordinator`, set in `New` and used at every phase site. **Coordinator-internal tests** scale them. **Cross-package tests wait the real budget** (AC1.3: 2→4 s, AC1.5: 2→4 s, AC1.7: 2 s; ≈ 10 s total; they block on `ctx.Done`, with no sleeps). Rejected: an exported `Config` knob, which is production surface the design did not approve and S7 would then require documenting |
| P2 | AC2.3's "blocks until `ctx.Done()`" leaves no handler (D3) | The plan returns at `deadline − 500ms`. Assert `handlerDeadline − planReturned ≥ HandlerCap − 100ms` (jitter only *adds* to the unmutated value). Use the real default budget, so the test also pins `New` |
| P3 | MUT-ORDER-FINISH survives AC1.7's e2e (D4) | Add a synthetic `dispatchError` row (`EffectsUnrecorded` wrapping a plan `PhaseTimeoutError`). The e2e arm is kept and kills MUT-UNWRAP-FINISH instead |
| P4 | MUT-ANYPHASE has no killing row in the design | Add a bare `PhaseTimeoutError{Phase:"finish"}` row expecting the generic deadline message |
| P5 | MUT-STALE-KEY is invisible to output equality (P-M3); the fixtures live in `coordinator` | Move `TestCapsuleWarmEqualsColdForSeTools` to `host/coordinator`. Its tooth is the `MOD010` warmth witness plus a cross-source arm; byte equality stays as AC3.3's determinism gate |
| P6 | A complete template is not content-rewritten by a run (P-M2) | AC3.2's test uses an **incomplete** template and an mtime-inclusive digest |
| P7 | where template key/promotion logic lives | One implementation in `host/archive` (`capsulecache.go`), used by `CheckSource`, the capsule and the daemon. `capsule/cachetemplate.go` holds copy-in, the label and the counter. `TestCapsulePromotionSurvivesCrossDevice` lives in `archive`, where the seam is |
| P8 | "retryable at the publish caller" | The bounded retry (3 attempts) sits inside `EnsureSourceLoadable`, so both callers (`verifyLoadable` and `cmd/world-publish/transitions.go:314`) inherit it. The typed error surfaces only after exhaustion |
| P9 | the hard failure's blast radius (D7) | Fix the fakes in the same commit (C2 list). A red outside the list stops the run |
| P10 | `publishReadTool` bypasses `PublishSet` (D9) | The fixture builds the template through `archive.CheckSource`. The tripwire is installed by `newWSDaemon` only when no log was supplied. Never exempt a rig from the tripwire |
| P11 | GC "≠ the current pinned digest" would prune templates of a head still served under another archived interpreter, making a GC/rebuild ping-pong on every start | Keep the digests of (daemon pin ∪ head descriptors' interpreters). Also prune `.tmp-*` older than 2 × the check bound; never touch a fresh `.tmp-*` (a concurrent `world-publish` may own it) |
| P12 | status-line timing and the listener point (D10) | Status lines come from `New` (a stat only). Rebuild and GC run from `Run` after `Listen` + announce, in a cancelled-and-joined goroutine. The AC3.7 test drives `Run` with an unexported `Config` gate |
| P13 | `(template missing\|rebuilding)` (D11) | The capsule emits `(template missing)` only. The daemon's own line `capsule template <id>: …` carries the rebuild state. No capsule↔daemon state hook |
| P14 | AC1.6 "every `r.call` asserts `wire.Error == nil`", but `:186` expects a refusal | `callOK` on every success-expecting call. `:186` keeps `r.call` and prints `*wire.Error` |
| P15 | operator-line prefix | New daemon lines carry the house prefix `ailang-worldd: ` (as `logRefusal` does). The capsule line is `capsule: cold compile …` written to the daemon's `ErrorLog`. Tests match by `Contains`/anchored regex on the unprefixed part |
| P16 | M3 executor-run split | Yes: Run 1 = C0–C4, Run 2 = C5–C6 (§6) |

---

## 9. Risks the executor must watch (report; never paper over)

1. **CI race-leg headroom (row 154).** M1 adds ~6 s of real-budget waits to `projection` and ~4 s to `daemon`. M3 should cut every capsule run's cost. Report per-package wall at each boundary (§2.3).
2. **Linux behaviour (design U1/U5).** Neither the `MOD010` warmth witness (P-M4) nor `AILANG_CACHE_DIR` has been measured on linux. AC3.4/AC3.5 run in CI and are the measurement. If CI reds there, that is evidence to report, not a test to relax.
3. **`TestCapsuleWarmIsFaster` is a timing ratio under CI load.** The interleaving and the `MOD010` tooth are what make it robust. If the ratio arm flakes in CI while `MOD010` holds, report it with the observed ratios. Do not loosen 2×.
4. **M1/M2 CI reruns are expected** (design §6). The executor never reruns CI. It hands the PR to the controller.
5. **Template trust.** Copy only regular files and dirs into a run (refuse symlinks). Templates are mode 0700 and host-owned. Never point `AILANG_CACHE_DIR` at the template (MUT-SHARE).
6. **Daemon goroutine lifetime (C7).** A leaked maintenance goroutine shows up as `-race` reports or `database is closed` lines in later tests. The latter is row 161's red herring; do not reintroduce it.

---

## 10. Effort

| Milestone | Estimate | Executor run |
|---|---|---|
| M0 re-baseline | 0.05 d | 1 |
| M1 (C1 + C2) | 0.45 d | 1 |
| M2 (C3) | 0.15 d | 1 |
| M3a (C4 / C4a + C4b) | 0.55 d | 1 |
| M3b (C5) | 0.40 d | 2 |
| M4 executor part (C6 + proxy) | 0.15 d + ~30 min idle host | 2 |
| **Total** | **~1.75 d** | (design: 1.5 d; +0.25 d is the 19-test fake migration, the warmth-witness tests and the P1 real-budget wire tests) |

## 11. Not verified by the planner

- Linux behaviour of every M3 mechanism (U1/U5): CI is the measurement.
- Whether the design's probe patch re-applies cleanly on the M3 head. It will not apply verbatim (it touches `runPhase` and `RunContext`); the executor adapts it uncommitted.
- The AC3.7 fixture's exact token/grant shape for an `ailang-read` call with `fakeToolBin`. The plan only needs the plan phase to run, which happens before the effect.
- The exact wall-time growth of the `daemon` package under CI `-race` after M3: measured per boundary locally, and in CI by the controller.

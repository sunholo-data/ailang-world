# Row 153 executor run 1 log (C0-C4b)

Executor: mission iteration 241, claude-sonnet-5-5, unattended. Worktree branch
`sprint/row153-plan-phase-deadline-under-load`, base `8e15d28`. Raw logs: `~/.ailang/state/mission-world-iter241/executor-run1/`
(`c0 c1 c2 c3 c4a c4b`, one dir per commit boundary; `mut.log` is the verbatim mutation record below).

## Commits

| Plan | SHA | Subject |
|---|---|---|
| C1 | b5a4a25 | typed plan/finish phase timeout (PhaseTimeoutError, classifyPhaseErr), budget hooks, A2A plan-timeout message |
| C2 | ed6b8d7 | diagnosable rig (testErrorLog, callOK), frozen MCP plan-timeout wire test, QUICKSTART line |
| C3 | 2f9e377 | PlanPhaseBudget = 4s derived, pinned; slow-plan test |
| C4a | b3219b2 | template build in CheckSource, TemplateBuildError, racer/EXDEV-safe promotion, retry; fakes |
| C4b | a2d68d1 | capsule copy-in, cold label + ColdRuns, lazy promotion, daemon log wiring, rig tripwire, tests |

## C0 baseline (HEAD 8e15d28 = code of 990a6da)

g1 rc=0, g2 rc=0, g4 (race) rc=0, g5 rc=0 20/20 `--- PASS`, gofmt empty, no binary blobs.
**g3 (`go build && go test ./... -count=1`) rc=1: one red, `host/broker` `TestRunBoundedHeadTailTimeoutKeepsPartialOutput`
("partial stdout = \"\"", a 0.30 s budget).** The planner measured g3 green on an idle host. Here it fails deterministically
in g3 (every boundary, same line) and in one race run while the executor was loading the host. It passes alone 3/3 at C0
(`c0/flake{1,2,3}.log`) and in isolation under `-race` (`c3/broker-race-rerun.log`). The package is outside every commit's
footprint, so it is attributed to the C0 baseline, not to any commit. It is a load-sensitive 0.3 s test and a hygiene
candidate; report to the controller.

## Gate results per boundary (rc; wall in s)

| Boundary | g1 verify_ail | g2 vet | g3 build+test | g4 race | g5 srt (-p 1 -v) | gofmt | blobs |
|---|---|---|---|---|---|---|---|
| C0 | 0 (22) | 0 | 1* (252) | 0 (507) | 0, 20/20 | empty | empty |
| C1 b5a4a25 | 0 (22) | 0 | 1* (252) | 0 (508) | 0, 20/20 | empty | empty |
| C2 ed6b8d7 | 0 (23) | 0 | 1* (254) | 0 (510) | 0, 20/20 | empty | empty |
| C3 2f9e377 | 0 (22) | 0 | 1* (268) | 1** (527) | 0, 20/20 | empty | empty |
| C4a b3219b2 | 0 (23) | 0 | 1* (258) | 0 (517) | 0, 20/20 | empty | empty |
| C4b a2d68d1 | 0 (23) | 0 | 1* (254) | 0 (512) | 0, 20/20 | empty | empty |

`*` = only the C0 `host/broker` red above; every other package `ok`. `**` = the same broker test under race while the executor
ran scratch work in parallel; broker race passes alone (rc=0). C4a had a first run with one real red (below), fixed
and re-run in full (`c4a-first/` is the red run, `c4a/` the green one).

Per-package wall (g3 plain / g4 race), the row-154 headroom figure:

| Boundary | coordinator | projection | daemon |
|---|---|---|---|
| C0 | 50.5 / 51.8 | 10.7 / 13.2 | 229.0 / 253.5 |
| C1 | 48.7 / 52.2 | 14.3 / 17.0 | 228.0 / 254.2 |
| C2 | 48.7 / 51.6 | 14.2 / 17.6 | 230.3 / 254.2 |
| C3 | 51.7 / 55.0 | 16.9 / 19.7 | 244.0 / 257.1 |
| C4a | 51.3 / 52.5 | 16.8 / 19.5 | 233.7 / 257.1 |
| C4b | 52.0 / 54.7 | 16.8 / 19.6 | 229.3 / 255.5 |

The race `daemon` tail is +2 s over C0 (253.5 -> 255.5), `coordinator` +3 s, `projection` +6 s (the real-budget wire waits).
The whole race leg is 507 -> 512 s on this host.

Extra steps: the local supplement after C4b (`WORLD_TOOL_AILANG_BIN`, `go test ./host/daemon/ ./cmd/ailang-worldd/ -v`): rc=0,
238 PASS, 3 SKIP (the same three), 0 FAIL, **0 `cold compile` lines** (baseline 236 PASS = +2 new daemon tests).
`bench_worldd.sh --smoke`: rc=0 PASSED.

## New-test PASS evidence

`c1/new.log` (5 PASS: TestClassifyPhaseErrMatrix, TestA2APlanTimeoutRetrySameTaskID, TestA2APlanTimeoutIsTypedAndNamesThePhase,
TestA2AFinishOverrunStaysEffectsUnrecorded, TestA2ADispatch); `c4b-pre/subset.log` (96 PASS, 0 SKIP across daemon, ailang-worldd,
capsule, coordinator, broker incl. every new C2/C3/C4 test).

## Mutation ledger

Method: the one-line edit is applied to the milestone's own production file with a backup copy, the named test runs
`-count=1 -v`, the file is restored from the backup (a `git checkout` would have destroyed the uncommitted milestone), and the
build is re-proven by the next gate. `host/...` lines below are the verbatim `mut.log` (long lines cut at 400 bytes).

| Mutation | Test | Red? | Restored green? |
|---|---|---|---|
| MUT-RELABEL | TestClassifyPhaseErrMatrix | yes (:40 (DE,DE) typed) | yes (C1 gate) |
| MUT-NOWRAP | TestClassifyPhaseErrMatrix | yes (:32 does not unwrap) | yes |
| MUT-EARLY-INTENT | TestA2APlanTimeoutRetrySameTaskID | yes (:93 UnconfirmedError on resend) | yes |
| MUT-ORDER | TestA2APlanTimeoutIsTypedAndNamesThePhase + TestA2ADispatch/PlanPhaseTimeout | yes | yes |
| MUT-GENERIC | same | yes | yes |
| MUT-ANYPHASE | TestA2ADispatch/FinishPhaseTimeout_bare | yes | yes |
| MUT-ORDER-FINISH | TestA2ADispatch/EffectsUnrecorded_wrapping_plan_timeout | yes (table row; wire arm survives by construction, D4) | yes |
| MUT-UNWRAP-FINISH | TestA2AFinishOverrunStaysEffectsUnrecorded | yes (:618) | yes |
| MUT-NOWRAP (MCP) | TestMCPPlanTimeoutFrozenWire | yes (`host callback canceled`) | yes |
| MUT-SILENT-ERROR-b | TestExecSrtMCPEndToEnd | yes, **at callOK** with `{Code:-32602 Message:unknown tool ...}` | yes |
| MUT-GUESS (3 s) | TestPlanPhaseBudgetIsDerived | yes | yes |
| MUT-GUESS-HANDLERCAP | TestPlanPhaseBudgetIsDerived + TestExecMaxTimeoutIsHandlerCapLessOneSecond | yes, both | yes |
| MUT-STEAL (5 s) | TestSlowPlanLeavesHandlerItsCap (9.498 s < 9.9 s) and TestPlanPhaseBudgetIsDerived | yes, both | yes |
| MUT-DEFAULT | TestSlowPlanLeavesHandlerItsCap (6.498 s) | yes | yes |
| MUT-PUBLISH-NOTEMPLATE (a: env append deleted) | TestPublishBuildsCapsuleTemplates | yes (TemplateBuildError) | yes |
| MUT-PUBLISH-NOTEMPLATE-b (env append and manifest check deleted) | same | yes (template absent) | yes |
| MUT-RACER-ERROR | TestTemplateRenameRacerIsSuccess | yes (`file exists, want nil`) | yes |
| MUT-WEDGE | TestPublishRetriesTemplateBuildError/second_attempt_succeeds | yes | yes |
| MUT-CROSS-DEVICE-RENAME | TestCapsulePromotionSurvivesCrossDevice | yes (`cross-device link`) | yes |
| MUT-TRIPWIRE-PROOF | TestExecSrtMCPEndToEnd | yes (`row 153 tripwire: capsule: cold compile`) | yes |
| MUT-SHARE | TestCapsuleTemplateIsCopiedNotShared | yes (`a run changed the template`) | yes |
| MUT-SILENT-COLD (run side) | TestCapsuleColdRunIsLabelled | yes (log "") | yes |
| MUT-NO-PROMOTE | TestCapsuleColdRunIsLabelled | yes (`did not promote a template`) | yes |
| MUT-NOTEMPLATE | TestCapsuleWarmIsFaster | yes (warm run compiled, MOD010) | yes |
| MUT-STALE-KEY | TestCapsuleWarmEqualsColdForSeTools/cross-source | yes (MOD010) | yes |
| MUT-CACHEDIRX (floor drift) | TestPinnedInterpreterHonoursCacheDir | yes (CheckSource TemplateBuildError) | yes |

Not applicable until run 2 (M3b): MUT-SILENT-COLD (startup), MUT-BLOCKING-PREWARM, MUT-SOFT-PUBLISH, MUT-NO-GC.

### Verbatim mutation record

```
=== MUT MUT-RELABEL  host/coordinator/effectful.go:136  [if parentErr == nil {] -> [if true {]  test=^TestClassifyPhaseErrMatrix$
    phase_timeout_test.go:40: (DE,DE) = coordinator: plan phase exceeded its 2s budget after 1s, want the caller's deadline, not a phase timeout
--- FAIL: TestClassifyPhaseErrMatrix (0.00s)
FAIL	github.com/sunholo-data/ailang-world/host/coordinator	0.295s
=== MUT MUT-EARLY-INTENT  host/coordinator/effectful.go:276  [	now := c.cfg.Now() // one logical time] -> [	c.uncertainIntent.Store(id, struct{}{})
	now := c.cfg.Now() // one logical time]  test=^TestA2APlanTimeoutRetrySameTaskID$
    phase_timeout_test.go:93: dispatch #2 (same task id) = coordinator: invocation a2a:ep1:t1 outcome not confirmed: intent receipt absent after uncertain append, want success
--- FAIL: TestA2APlanTimeoutRetrySameTaskID (0.11s)
FAIL	github.com/sunholo-data/ailang-world/host/coordinator	0.409s
=== MUT MUT-NOWRAP  host/coordinator/errors.go:176  [Unwrap() error { return context.DeadlineExceeded }] -> [Unwrap() error { return context.Canceled }]  test=^TestClassifyPhaseErrMatrix$
    phase_timeout_test.go:32: (nil,DE) coordinator: plan phase exceeded its 2s budget after 2.01s does not unwrap to context.DeadlineExceeded
--- FAIL: TestClassifyPhaseErrMatrix (0.00s)
FAIL	github.com/sunholo-data/ailang-world/host/coordinator	0.504s
=== MUT MUT-ORDER  host/projection/projection.go  test=^(TestA2APlanTimeoutIsTypedAndNamesThePhase|TestA2ADispatch)$
< 	case errors.As(err, &phase) && phase.Phase == "plan":
< 		return codeInternal, fmt.Sprintf("%s %s budget; nothing ran or was committed; resend the same task id", PhaseTimeoutPrefix, phase.Budget)
> 	case errors.As(err, &phase) && phase.Phase == "plan":
> 		return codeInternal, fmt.Sprintf("%s %s budget; nothing ran or was committed; resend the same task id", PhaseTimeoutPrefix, phase.Budget)
    a2a_wire_test.go:594: error=-32603 "invocation exceeded its deadline", want -32603 "transition plan phase exceeded its 2s budget; nothing ran or was committed; resend the same task id"; body={"error":{"code":-32603,"message":"invocation exceeded its deadline"},"id":"wire-plan-timeout","jsonrpc":"2.0"}
--- FAIL: TestA2APlanTimeoutIsTypedAndNamesThePhase (2.01s)
    projection_test.go:857: got -32603 "invocation exceeded its deadline", want -32603 "transition plan phase exceeded its 2s budget; nothing ran or was committed; resend the same task id"
--- FAIL: TestA2ADispatch (0.02s)
    --- FAIL: TestA2ADispatch/PlanPhaseTimeout (0.00s)
FAIL	github.com/sunholo-data/ailang-world/host/projection	2.355s
=== MUT MUT-ORDER-FINISH  host/projection/projection.go  test=^(TestA2AFinishOverrunStaysEffectsUnrecorded|TestA2ADispatch)$
> 	case errors.As(err, &phase) && phase.Phase == "plan":
> 		return codeInternal, fmt.Sprintf("%s %s budget; nothing ran or was committed; resend the same task id", PhaseTimeoutPrefix, phase.Budget)
< 	case errors.As(err, &phase) && phase.Phase == "plan":
< 		return codeInternal, fmt.Sprintf("%s %s budget; nothing ran or was committed; resend the same task id", PhaseTimeoutPrefix, phase.Budget)
    projection_test.go:857: got -32603 "transition plan phase exceeded its 0s budget; nothing ran or was committed; resend the same task id", want -32603 "effects were requested but the invocation was not confirmed committed; effect records: none"
--- FAIL: TestA2ADispatch (0.02s)
    --- FAIL: TestA2ADispatch/EffectsUnrecorded_wrapping_plan_timeout (0.00s)
FAIL	github.com/sunholo-data/ailang-world/host/projection	2.357s
=== MUT MUT-GENERIC  host/projection/projection.go:444  ["%s %s budget; nothing ran or was committed; resend the same task id"] -> ["transition timed out%.0s%.0s"]  test=^(TestA2APlanTimeoutIsTypedAndNamesThePhase|TestA2ADispatch)$
    a2a_wire_test.go:594: error=-32603 "transition timed out", want -32603 "transition plan phase exceeded its 2s budget; nothing ran or was committed; resend the same task id"; body={"error":{"code":-32603,"message":"transition timed out"},"id":"wire-plan-timeout","jsonrpc":"2.0"}
--- FAIL: TestA2APlanTimeoutIsTypedAndNamesThePhase (2.01s)
    projection_test.go:857: got -32603 "transition timed out", want -32603 "transition plan phase exceeded its 2s budget; nothing ran or was committed; resend the same task id"
--- FAIL: TestA2ADispatch (0.01s)
    --- FAIL: TestA2ADispatch/PlanPhaseTimeout (0.00s)
FAIL	github.com/sunholo-data/ailang-world/host/projection	2.331s
=== MUT MUT-ANYPHASE  host/projection/projection.go:16  [ && phase.Phase == "plan":] -> [:]  test=^(TestA2APlanTimeoutIsTypedAndNamesThePhase|TestA2ADispatch)$
    projection_test.go:857: got -32603 "transition plan phase exceeded its 2s budget; nothing ran or was committed; resend the same task id", want -32603 "invocation exceeded its deadline"
--- FAIL: TestA2ADispatch (0.02s)
    --- FAIL: TestA2ADispatch/FinishPhaseTimeout_bare (0.00s)
FAIL	github.com/sunholo-data/ailang-world/host/projection	2.334s
=== MUT MUT-UNWRAP-FINISH  host/coordinator/effectful.go:91  [	if err != nil {
		return unrecorded(err)
	}
	p := planEffectInvocation(] -> [	if err != nil {
		return Result{}, err
	}
	p := planEffectInvocation(]  test=^TestA2AFinishOverrunStaysEffectsUnrecorded$
    a2a_wire_test.go:618: answer = -32603 "invocation exceeded its deadline", want -32603 with prefix "effects were requested but the invocation was not confirmed committed; effect records:"; body={"error":{"code":-32603,"message":"invocation exceeded its deadline"},"id":"wire-finish-overrun","jsonrpc":"2.0"}
--- FAIL: TestA2AFinishOverrunStaysEffectsUnrecorded (2.01s)
FAIL	github.com/sunholo-data/ailang-world/host/projection	2.335s
=== MUT MUT-NOWRAP-MCP  host/coordinator/errors.go:176  [Unwrap() error { return context.DeadlineExceeded }] -> [Unwrap() error { return context.Canceled }]  test=^TestMCPPlanTimeoutFrozenWire$
    mcp_plan_timeout_test.go:50: tools/call under a plan timeout error = {Code:-32603 Message:host callback canceled}, want the frozen -32603 "host callback timed out" (row 159 / ailang#1602 flips this)
--- FAIL: TestMCPPlanTimeoutFrozenWire (2.30s)
FAIL	github.com/sunholo-data/ailang-world/host/daemon	2.660s
=== MUT MUT-SILENT-ERROR-b  host/daemon/exec_e2e_srt_test.go:176  [r.callOK(token, "workspace-exec", map[string]any{"command": "where"})] -> [r.callOK(token, "workspace-exec-x", map[string]any{"command": "where"})]  test=^TestExecSrtMCPEndToEnd$
    exec_e2e_srt_test.go:148: workspace-exec test: map[argv:[printf [%s] -run=TestX a] duration_ms:126 exit_code:0 limit:<nil> profile:map[argv0_sha256:d7992bc59e048ee061716365839990227ab463d954d87e6c2e7f987c0054fd59 command:test digest:sha256:a1ce9fbbeef30ee600c242afe488ee9218235997fffb44c75a4bb00c9a13f96f project:proj] sandbox:map[cli_sha256:3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2
    exec_e2e_srt_test.go:176: workspace-exec-x: wire error {Code:-32602 Message:unknown tool "workspace-exec-x"}
--- FAIL: TestExecSrtMCPEndToEnd (4.27s)
FAIL	github.com/sunholo-data/ailang-world/host/daemon	4.630s
=== MUT MUT-GUESS  host/coordinator/effectful.go:45  [PlanPhaseBudget   = 4 * time.Second] -> [PlanPhaseBudget   = 3 * time.Second]  test=^TestPlanPhaseBudgetIsDerived$
    plan_budget_test.go:16: PlanPhaseBudget = 3s, want invokeDeadline − HandlerCap − HandlerHeadroom = 20s − 10s − 6s = 4s
--- FAIL: TestPlanPhaseBudgetIsDerived (0.00s)
FAIL	github.com/sunholo-data/ailang-world/host/daemon	0.357s
=== MUT MUT-GUESS-HANDLERCAP  host/coordinator/effectful.go:47  [HandlerCap        = 10 * time.Second] -> [HandlerCap        = 9 * time.Second]  test=^(TestPlanPhaseBudgetIsDerived|TestExecMaxTimeoutIsHandlerCapLessOneSecond)$
    plan_budget_test.go:16: PlanPhaseBudget = 4s, want invokeDeadline − HandlerCap − HandlerHeadroom = 20s − 9s − 6s = 5s
--- FAIL: TestPlanPhaseBudgetIsDerived (0.00s)
    workspace_exec_test.go:129: broker.ExecMaxTimeoutMS = 9000, want HandlerCap − 1 s = 8000
--- FAIL: TestExecMaxTimeoutIsHandlerCapLessOneSecond (0.00s)
FAIL	github.com/sunholo-data/ailang-world/host/daemon	0.367s
=== MUT MUT-STEAL  host/coordinator/effectful.go:45  [PlanPhaseBudget   = 4 * time.Second] -> [PlanPhaseBudget   = 5 * time.Second]  test=^TestSlowPlanLeavesHandlerItsCap$
    phase_timeout_test.go:161: the handler had 9.498192875s after the slow plan, want >= HandlerCap−100ms = 9.9s: the plan cap steals the handler's time
--- FAIL: TestSlowPlanLeavesHandlerItsCap (4.53s)
FAIL	github.com/sunholo-data/ailang-world/host/coordinator	4.853s
=== MUT MUT-STEAL-derived  host/coordinator/effectful.go:45  [PlanPhaseBudget   = 4 * time.Second] -> [PlanPhaseBudget   = 5 * time.Second]  test=^TestPlanPhaseBudgetIsDerived$
    plan_budget_test.go:16: PlanPhaseBudget = 5s, want invokeDeadline − HandlerCap − HandlerHeadroom = 20s − 10s − 6s = 4s
--- FAIL: TestPlanPhaseBudgetIsDerived (0.00s)
FAIL	github.com/sunholo-data/ailang-world/host/daemon	0.362s
=== MUT MUT-DEFAULT  host/coordinator/coordinator.go:83  [planBudget: PlanPhaseBudget,] -> [planBudget: 2 * PlanPhaseBudget,]  test=^TestSlowPlanLeavesHandlerItsCap$
    phase_timeout_test.go:161: the handler had 6.498212333s after the slow plan, want >= HandlerCap−100ms = 9.9s: the plan cap steals the handler's time
--- FAIL: TestSlowPlanLeavesHandlerItsCap (7.53s)
FAIL	github.com/sunholo-data/ailang-world/host/coordinator	7.847s
=== MUT MUT-PUBLISH-NOTEMPLATE  host/archive/check.go:16  [, CacheDirEnv+"="+tmp)] -> [)]  test=^TestPublishBuildsCapsuleTemplates$
    capsule_template_test.go:68: PublishSet: publish transition registry: capsule template for interpreter sha256:1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f source sha256:5ff3d1fc6c771a0d08ea087edbcb72b286fd2794a7a6bc1e7ff47c2242486085: check passed but wrote no compile/manifest.json
--- FAIL: TestPublishBuildsCapsuleTemplates (1.02s)
FAIL	github.com/sunholo-data/ailang-world/host/transitionreg	1.352s
=== MUT MUT-PUBLISH-NOTEMPLATE-b  host/archive/check.go  test=^TestPublishBuildsCapsuleTemplates$
< 	cmd.Env = append(childenv.Scrubbed(os.Environ()), "AILANG_RELAX_MODULES=1", CacheDirEnv+"="+tmp)
> 	cmd.Env = append(childenv.Scrubbed(os.Environ()), "AILANG_RELAX_MODULES=1")
< 	if info, err := os.Stat(filepath.Join(tmp, filepath.FromSlash(manifestRel))); err != nil || !info.Mode().IsRegular() {
< 		return passed, &TemplateBuildError{Interpreter: interpreter, Source: hashref.SumSHA256(source),
< 			Reason: "check passed but wrote no compile/manifest.json"}
< 	}
    capsule_template_test.go:72: tools.a: no ready capsule template at /var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/TestPublishBuildsCapsuleTemplates2502481231/001/world.db.artifacts/capsule-cache/1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f/5ff3d1fc6c771a0d08ea087edbcb72b286fd2794a7a6bc1e7ff47c2242486085 after publication
--- FAIL: TestPublishBuildsCapsuleTemplates (0.86s)
FAIL	github.com/sunholo-data/ailang-world/host/transitionreg	1.179s
=== MUT MUT-WEDGE  host/transitionreg/verify.go:80  [const templateBuildAttempts = 3] -> [const templateBuildAttempts = 1]  test=^TestPublishRetriesTemplateBuildError$
    capsule_template_test.go:149: PublishSet with a once-failing template build: publish transition registry: capsule template for interpreter sha256:5e6d352939f7147d9377731a1380bb38f41f771f78292834d541f8f22299e393 source sha256:c4c502d9c770f2a0fcd2da160b3e348297b2e55942c3f94c766c335150fe11d0: check passed but wrote no compile/manifest.json
--- FAIL: TestPublishRetriesTemplateBuildError (0.52s)
    --- FAIL: TestPublishRetriesTemplateBuildError/second_attempt_succeeds (0.26s)
FAIL	github.com/sunholo-data/ailang-world/host/transitionreg	0.830s
=== MUT MUT-CROSS-DEVICE-RENAME  host/archive/capsulecache.go:137  [if err := CopyCacheTree(runCacheDir, tmp); err != nil {] -> [if err := renameDir(runCacheDir, tmp); err != nil {]  test=^TestCapsulePromotionSurvivesCrossDevice$
    capsulecache_test.go:114: PromoteCopy across devices: capsule template for interpreter sha256:9666d9e8899447735cb9897b77dbb121754fd4d503609758755bd3fdae3a4b22 source sha256:7b79504c1c3ab1bf5e14d63a44f209cd217d9e1a47ba4a0e8dcb432c9647ea30: copy run cache: rename /var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/TestCapsulePromotionSurvivesCrossDevice1494409310/002 /var/folders/kr/h0wr2sj94vd6ljt
--- FAIL: TestCapsulePromotionSurvivesCrossDevice (0.01s)
FAIL	github.com/sunholo-data/ailang-world/host/archive	0.270s
=== MUT MUT-RACER-ERROR  host/archive/capsulecache.go:119  [return errors.Is(err, fs.ErrExist) || errors.Is(err, syscall.ENOTEMPTY) || errors.Is(err, syscall.EEXIST)] -> [return false && (errors.Is(err, fs.ErrExist) || errors.Is(err, syscall.ENOTEMPTY) || errors.Is(err, syscall.EEXIST))]  test=^TestTemplateRenameRacerIsSuccess$
    capsulecache_test.go:74: PromoteTemplate lost to a racer: capsule template for interpreter sha256:9666d9e8899447735cb9897b77dbb121754fd4d503609758755bd3fdae3a4b22 source sha256:7b79504c1c3ab1bf5e14d63a44f209cd217d9e1a47ba4a0e8dcb432c9647ea30: promote template: rename /var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/TestTemplateRenameRacerIsSuccess23562382/001/world.db.artifacts/capsule-cache/.
--- FAIL: TestTemplateRenameRacerIsSuccess (0.01s)
FAIL	github.com/sunholo-data/ailang-world/host/archive	0.272s
=== MUT MUT-TRIPWIRE-PROOF  host/daemon/exec_e2e_srt_test.go:72  [	publishSeTools(t, d, f.db)
	srv :=] -> [	publishSeTools(t, d, f.db)
	_ = os.RemoveAll(filepath.Join(f.db+".artifacts", "capsule-cache"))
	srv :=]  test=^TestExecSrtMCPEndToEnd$
    testlog_test.go:56: capsule: cold compile 1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f/51d93ce680eea3031eaf786530b3c2a4d61b3b547440a4e542786695b05c9a05 (template missing)
    testlog_test.go:58: a rig ran a transition with no capsule template (row 153 tripwire): capsule: cold compile 1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f/51d93ce680eea3031eaf786530b3c2a4d61b3b547440a4e542786695b05c9a05 (template missing)
    exec_e2e_srt_test.go:150: workspace-exec test: map[argv:[printf [%s] -run=TestX a] duration_ms:126 exit_code:0 limit:<nil> profile:map[argv0_sha256:d7992bc59e048ee061716365839990227ab463d954d87e6c2e7f987c0054fd59 command:test digest:sha256:a1ce9fbbeef30ee600c242afe488ee9218235997fffb44c75a4bb00c9a13f96f project:proj] sandbox:map[cli_sha256:3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2
--- FAIL: TestExecSrtMCPEndToEnd (4.56s)
FAIL	github.com/sunholo-data/ailang-world/host/daemon	4.920s
=== MUT MUT-SHARE  host/capsule/cachetemplate.go:41  [return cacheDir, false, nil] -> [return r.archive.CapsuleTemplateDir(interpreter, source), false, nil]  test=^TestCapsuleTemplateIsCopiedNotShared$
    cachetemplate_test.go:172: a run changed the template.
--- FAIL: TestCapsuleTemplateIsCopiedNotShared (1.07s)
FAIL	github.com/sunholo-data/ailang-world/host/capsule	1.319s
=== MUT MUT-SILENT-COLD  host/capsule/cachetemplate.go:50  [	fmt.Fprintf(r.cs.log, "capsule: cold compile %s/%s (template missing)\n", interpreter.Digest(), hashref.SumSHA256(source).Digest())] -> [	_ = hashref.SumSHA256]  test=^TestCapsuleColdRunIsLabelled$
    cachetemplate_test.go:190: after one cold run: log "", ColdRuns 1; want exactly one matching line and 1
--- FAIL: TestCapsuleColdRunIsLabelled (0.96s)
FAIL	github.com/sunholo-data/ailang-world/host/capsule	1.200s
=== MUT MUT-NO-PROMOTE  host/capsule/capsule.go:291  [		r.promoteCold(cacheDir, entry.Interpreter, entry.Source)] -> [		_ = cold]  test=^TestCapsuleColdRunIsLabelled$
    cachetemplate_test.go:193: the successful cold run did not promote a template
--- FAIL: TestCapsuleColdRunIsLabelled (1.00s)
FAIL	github.com/sunholo-data/ailang-world/host/capsule	1.238s
=== MUT MUT-STALE-KEY  host/archive/capsulecache.go:74  [interpreter.Digest(), hashref.SumSHA256(source).Digest())] -> [interpreter.Digest(), hashref.SumSHA256(nil).Digest())]  test=^TestCapsuleWarmEqualsColdForSeTools$
    capsule_template_test.go:95: workspace-exec compiled although its own template was built (stderr "WARNING MOD010 (relaxed): module 'se_tools/exec' does not match canonical path 'host/capsule/main'\n  Running under --relax-modules; mismatch ignored. For strict checking, omit --relax-modules flag.\n"): the template key ignores the source
--- FAIL: TestCapsuleWarmEqualsColdForSeTools (5.40s)
    --- FAIL: TestCapsuleWarmEqualsColdForSeTools/cross-source (0.43s)
FAIL	github.com/sunholo-data/ailang-world/host/coordinator	5.723s
=== MUT MUT-CACHEDIRX  host/archive/capsulecache.go:36  [const CacheDirEnv = "AILANG_CACHE_DIR"] -> [const CacheDirEnv = "AILANG_CACHE_DIRX"]  test=^TestPinnedInterpreterHonoursCacheDir$
    cachetemplate_test.go:245: CheckSource = ({Output:→ Type checking host/capsule/main.ail...
--- FAIL: TestPinnedInterpreterHonoursCacheDir (1.00s)
FAIL	github.com/sunholo-data/ailang-world/host/capsule	1.264s
=== MUT MUT-NOTEMPLATE  host/capsule/cachetemplate.go:40  [archive.CopyCacheTree(r.archive.CapsuleTemplateDir(interpreter, source), cacheDir)] -> [archive.CopyCacheTree(cacheDir, cacheDir)]  test=^TestCapsuleWarmIsFaster$
    cachetemplate_test.go:229: warm run 0 compiled (stderr "WARNING MOD010 (relaxed): module 'se_tools/exec' does not match canonical path 'host/capsule/main'\n  Running under --relax-modules; mismatch ignored. For strict checking, omit --relax-modules flag.\n"): the template was not copied in
--- FAIL: TestCapsuleWarmIsFaster (1.20s)
FAIL	github.com/sunholo-data/ailang-world/host/capsule	1.473s
```

## Deviations from the plan (each measured)

1. **MUT-NOWRAP form.** "Delete `Unwrap()`" does not compile (`context` becomes unused), so the mutation is
   `Unwrap() error { return context.Canceled }`. Same effect on `errors.Is(err, context.DeadlineExceeded)`. A first
   attempt as written is a build failure, not a red test, and was discarded.
2. **MUT-ORDER / MUT-ORDER-FINISH / MUT-RACER-ERROR / MUT-NOTEMPLATE** needed edit forms that keep the build green
   (move the case; `false && (...)`; `CopyCacheTree(cacheDir, cacheDir)`). Equivalent to the plan's lines.
3. **MUT-SILENT-ERROR** as planned (`command: "no-such-command"`) does NOT reach `callOK`: an unknown *command* is a
   successful tool call with an empty result, so the test fails later at `where = "<nil>"` (measured, first run
   in the first attempt, whose log entry was replaced by the proof below). The proof that `callOK` is reached and prints the wire error uses an unknown *tool name*
   (`workspace-exec-x`) -> `-32602 unknown tool` at `exec_e2e_srt_test.go:176`.
4. **MUT-STALE-KEY** reds only the *cross-source* arm. The per-phase warm arms build their own template immediately
   before running, so a source-blind key still serves them. The plan expected the warm arms to red as well; the
   cross-source arm is the only tooth, and it is sufficient.
5. **`TestCapsulePromotionSurvivesCrossDevice` stub.** The plan's "EXDEV whenever `Dir(old) != Dir(new)`" also rejects
   the legitimate `capsule-cache/.tmp-X` -> `capsule-cache/<digest>/<sha>` rename (different parent directories, one
   device): measured red on the unmutated code. The stub is "EXDEV for any source outside `capsule-cache/`", which
   models a run root on another device exactly. The test additionally asserts that a direct cross-device
   `PromoteTemplate` is a `*TemplateBuildError` wrapping `EXDEV`.
6. **MUT-CROSS-DEVICE-RENAME** applied as `renameDir(runCacheDir, tmp)` (the first rename PromoteCopy would make).
7. **Fifth fake interpreter.** The plan listed four fake scripts (C2). `TestEpochTwinInterpretersBothListedOnCard` has an
   additional inline fake in `host/daemon/registry_publisher_test.go` (`writeFake`, ~:251) whose `check` also exits 0
   without a cache. It was red in the first C4a gate (`c4a-first/g3.log`) with `check passed but wrote no
   compile/manifest.json`. The test was already on the plan's list; the red set was exactly the planner's 19 plus
   no unlisted test. Fixed in the same commit (C4a) and the full gates re-run. Pre-fix measured red set (subset run):
   archive 1, transitionreg 12, daemon 3, world-publish 3 = 19, matching P-M5.
8. **Rig tripwire for the exec rig.** `newExecSrtRig` passes its own ErrorLog, so the `newWSDaemon` default tripwire
   would not apply. It uses `testOperatorLog(t)` (tripwire on), which is what the plan's proof
   (MUT-TRIPWIRE-PROOF) requires.
9. **MCP tool name in `TestMCPPlanTimeoutFrozenWire`.** `ws.read` is addressed as `ws_dread` (MCP names escape `.`).
10. **P7 placement.** `copyTree` lives once, in `host/archive` as `CopyCacheTree` (exported), used by `PromoteCopy` and the
    capsule's copy-in, instead of a second copy in `capsule/cachetemplate.go`. `serveR14` was generalised into
    `serveWith(coordStore, runner, errLog)` so the frozen-wire test reuses it.
11. **QUICKSTART (C10).** Only the cold-compile line, the template directory and the hard publish failure are
    documented in C4b. The `capsule template <id>:` status lines belong to M3b (run 2).

## Findings for the controller

- **F1 (baseline red).** `host/broker` `TestRunBoundedHeadTailTimeoutKeepsPartialOutput` fails in `go test ./...` on this host at
  every boundary including C0, passes alone and under `-race` alone. Not caused by row 153. Row-114/hygiene candidate.
- **F2 (unlisted fake).** A fifth fake interpreter existed (deviation 7); a future fake that `check`s without writing a
  manifest will red in the hard failure path by design.
- **F3 (race leg).** race `daemon` 253.5 -> 255.5 s and the whole race leg 507 -> 512 s; no CI-headroom regression locally, but
  the CI runner is the measurement (row 154).
- **F4 (warm ratio).** `TestCapsuleWarmIsFaster` observed cold ~204 ms vs warm ~71 ms (2.9x) while a gate leg ran in parallel
  (planner idle-host: 171 vs 27 ms). The 2x bound has margin but the `MOD010` tooth is what holds under CI load.
- **F5 (CI pending).** M1/M2 CI reruns, AC4.2 5/5, and Linux behaviour of `AILANG_CACHE_DIR`/`MOD010` (U1/U5) are
  measured only by CI and the controller; nothing was pushed.

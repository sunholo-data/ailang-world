# Iteration 242 — row 152 evaluator r1 (independent judge)

Judge: claude-opus-5-5 (MISSION-ROLE evaluator; did not write the code).
Under review: `44f222b` (C1) · `3d3ee7b` (C2) · `d8e9c72` (C3) on plan `f1cbbb0`.
Worktree: `.wt-world-iter242-eval` at `d8e9c72`, left byte-clean (`git diff --quiet` → CLEAN).
Gate interpreter: `~/.pinned-ailang/ailang` = AILANG v0.41.0. Scratch: `~/.ailang/state/mission-world-iter242/eval/`.

## Verdict: PASS — 91/100, zero blocking findings

| Rubric area | Max | Score | Note |
|---|---|---|---|
| Acceptance criteria met (bar + AC1.1–1.9, AC2.1–2.5) | 35 | 34 | every code AC met; AC3.1 (records) is the controller's M3, not in this diff |
| Correctness / concurrency | 20 | 18 | lock discipline sound; stale-root overwrite in an unreachable state (N2) |
| Tests and non-vacuity | 20 | 17 | 10/12 executed mutants red; budget application unpinned (N1); MUT-SHARED-LOCK is equivalent at head (N3) |
| Code quality / standards | 10 | 10 | small, commented, no new surface; context-root census respected |
| Docs (QUICKSTART, comments) | 10 | 7 | mostly right; three inaccuracies (N4) |
| Gates | 5 | 5 | all green, measured by me |

## 1. The bar

Row bar: "a failure in one family must leave the other's typed behaviour intact", pinned by a test that times out the policy
summary and asserts `workspace-exec` still runs. **Met.** `TestWorkspaceExecSurvivesPolicySummaryTimeout/timeout`
(`host/daemon/workspace_family_test.go:195`) drives a summary that never answers (hold never released, 300 ms budget seam).
The test checks four things: (i) ws.exec commits through a real coordinator with `exit_code 0` and stdout = ep1 worktree;
(ii) ws.read gets R8 `[Workspace.Read]` with 0 plans; (iii) the operator log has the `policy summary: … timed out` line;
(iv) the error is `*broker.HandlerTimeoutError`. The `/healthy` control commits both calls. MUT-ALL-OR-NOTHING reds (i) (E1 below).

AC walk:
- AC1.1: `registry()` (`workspace.go:338`) binds exec whenever `enabled && episodeRoot ok`, and the eight names only on a built
  handler. Both `sandboxRoot` and `episodeHandler` failures go through `ailangFamily` and print the unchanged line. Met.
- AC1.2: `execMu` guards `execH` only (`workspace.go:382`). `execHandler` never takes `w.mu`. Met.
- AC1.3: `var workspaceHandlerBudget = 3 * time.Second` with a seam comment. Met.
- AC1.4: met (above). `publishTools` generalises `publishReadTool`, which now wraps it and keeps the row-153 template block intact.
- AC1.5: the four refusal tests assert exactly `[Workspace.Exec]` AND run one `where` call
  (`requireExecRunsIn`, exit 0 + stdout = worktree). Their summaries==0 and exact-log halves are unchanged. Met.
- AC1.6: `workspace_test.go` refusal tests are unedited (the diff only touches `publishReadTool`). Met.
- AC1.7: split across `TestWorkspaceExecNeverWaitsOnAilangBuild` (execHandler half) and `…DoesNotBlockOtherEpisodes`
  (registry half), plus `TestWorkspaceTwoEpisodeConcurrentRegistryExecRace` (2 eps × 8 goroutines × registry+execHandler,
  pointer-equal exec, asserted after join). Met.
- AC1.8: QUICKSTART and code comments updated. Met, with inaccuracies (N4).
- AC1.9: `TestWorkspaceExecSurvivesPolicyDirFault` plants a file at `<stateDir>/policies`, which is root-proof. It asserts exactly
  `[Workspace.Exec]`, 0 summaries, a committed exec Dispatch, and the ENOTDIR line. Met.
- AC2.1: the `ailangBuild` map is under `w.mu`, `w.mu` is released before the subprocess, and failures are uncached. Met.
- AC2.2: `…SingleFlight`: 1 spawn per episode, both callers get the result, 2 lines on failure, the same handler pointer on success. Met.
- AC2.3: met (DoesNotBlockOtherEpisodes + TwoEpisode re-run under -race at head).
- AC2.4: `git diff 3d3ee7b d8e9c72` touches no `workspace_test.go`/`workspace_layout_test.go`, so M2 edited none of the listed
  assertions, and all of them pass. Met.
- AC2.5: `…FailedBuildRetriesEachCall`: 3 spawns, 3 lines, exec-only each time, then the 4th call binds nine with no extra line. Met.
- AC3.1 (charter → LANDED, follow-up rows, M4 README note): controller/M3 scope, absent from this diff by design.

## 2. Correctness and concurrency (`host/daemon/workspace.go`)

- `episodeHandler` (`:433`):
  - Cache hit, join and new-build decisions all happen under `w.mu`, and the lock is dropped before the subprocess.
  - `b.tool`/`b.err` are written under `w.mu` before `close(b.done)`. Waiters read them only after `<-b.done`, so the
    channel close gives happens-before.
  - The in-flight entry is deleted only if it is still `b`, which is correct when a newer build replaced it.
  - No path returns with `w.mu` held.
  - My O2 (publish after close) is caught by `-race`.
- `execMu` and `w.mu` are never nested. `registry` takes them sequentially, never together. No lock-order cycle.
- Waiters block on `<-b.done` with no context. They are bounded only because the builder's ctx (3 s) bounds
  `NewAilangToolHandler`, and the broker's own `runBounded` adds a 10 s fallback. The mkdirs and the `O_NONBLOCK` lock check
  take no ctx, but they were unbounded before this change too. No new goroutine or process leak: the build runs on the
  caller's goroutine.
- Deviation (1), the clock starting before the mkdirs: verified as forced. `host/store/context_roots_test.go:46` pins
  `workspace.go|episodeHandler|Background` = 1, and `TestProductionContextRoots` passes. Effect: the budget now covers the mkdirs,
  the cache link and the policy write too, which only tightens it (they are µs–ms). Accepted.
- Deviations (2)–(4) accepted:
  - (2) `waitFor` is the existing bounded poller.
  - (3) `pendingRegistry` is a plain done+result pair, read only after done.
  - (4) the `exec_e2e_srt_test.go` switch to `execMu` is REQUIRED, since `execH` is now `execMu`-guarded. Under w.mu it would
    race. I re-ran the srt e2e trio under `-race -p 1` and all 3 passed.
- Episode-root changes: see N2. The persistent per-call 3 s stall in a broken episode, `workspace-exec` included, is kept as the
  design's stated cost (§5 c5, §11).

## 3. Mutations (all applied in my worktree via `eval/drill.py`; file sha256 `48ddbea4836d…` before and after every run → IDENTICAL)

| # | Mutation | Test | Result |
|---|---|---|---|
| E1 (executor's MUT-ALL-OR-NOTHING) | `return broker.Registry{}` after the AILANG operator line | `…SurvivesPolicySummaryTimeout` | RED `workspace_family_test.go:201: ws.exec Dispatch in ep1: … no registered handler [Workspace.Exec]` |
| E2 (MUT-SHARED-LOCK) | `execHandler` takes `w.mu` | `…ExecNeverWaitsOnAilangBuild` | **SURVIVED at head**: an equivalent mutant once C3 stopped holding `w.mu` across the subprocess. The executor's red was at C2 (`drill-C2-MUT-SHARED-LOCK.log`, valid there). See N3 |
| E3 race drill (MUT-EXECH-UNLOCKED) | no lock in `execHandler` | `…TwoEpisodeConcurrentRegistryExecRace -race` | RED `WARNING: DATA RACE` ×3 |
| E4 (MUT-NO-SINGLEFLIGHT) | never join an in-flight build | `…AilangBuildSingleFlight` | RED `:424 … ran 2 summaries, want 1` |
| E5 (MUT-NEGATIVE-CACHE) | keep the failed in-flight entry (`err == nil &&` on the delete) | `…FailedBuildRetriesEachCall` | RED `:495 after 2 failing calls: 1 summaries, want 2` |
| E6 (MUT-SANDBOX-GATES-EXEC + USES-SANDBOX) | bind exec only if `sandboxRoot` ok, with the sandbox as root | ModuleRootSymlink + CoversLock | RED ×12 `workspace_layout_test.go:197: registry = [], want exactly [Workspace.Exec]` |
| E7 (MUT-LOCK-ACROSS-SUBPROCESS) | hold `w.mu` across the build | `…DoesNotBlockOtherEpisodes` | RED `:479 ep1's AILANG build had already finished when ep2's registry returned` |
| O1 (mine) | waiter skips `<-b.done` | `…SingleFlight` | RED `:427 registry = [all nine], want exactly [Workspace.Exec]` |
| O2 (mine, race drill) | `close(b.done)` before publishing `b.tool, b.err` | `…SingleFlight -race` | RED `WARNING: DATA RACE` |
| O3 (mine) | drop the success-cache store | `^TestWorkspace` | RED `workspace_layout_test.go:129: summaries = 2, want 1` (+ RefusesEpisodesOutsideTheGrammarAndRoot) |
| O4 (mine) | budget not applied (`WithCancel` instead of `WithTimeout(budget)`) | `…SurvivesPolicySummaryTimeout` | **SURVIVED**: PASS in 30.45 s vs 1.43 s. The broker's internal 10 s cap ends each of the 3 builds instead. See N1 |
| O5 (mine) | join an in-flight build regardless of `b.root` | `^TestWorkspace` | **SURVIVED**: the root-mismatch state is unreached by any test, and in practice unreachable. See N2 |

Totals: executor's named mutants 6/7 red (E2 equivalent at head). Mine: 3/5 red, 2 survived (one gap, one unreachable state).

## 4. Load/timing fragility

- No elapsed-time or ratio assertions anywhere.
- Wait bounds:
  - `waitFor` uses 30 s on preconditions. A slow machine makes it red, never falsely green.
  - The 5 s budget seam (PD12) is the one real bound. For `NeverWaits`/`DoesNotBlock`/`SingleFlight`, correct code must finish
    ep2's work (one stub build + one exec) within 5 s of ep1's summary starting. Measured: 0.3–0.7 s under `-race`, a 7–15× margin.
  - Failure direction: a slow box gives a false RED ("ep1 had already finished"), never a false pass. Acceptable. If CI ever flakes
    here, raise it to `boundedTestContext`'s 30 s (the design wording "test's own bounded context").
- The 300 ms budget in AC1.4 only has to be exceeded, which is guaranteed (the hold is never released).
- The never-released hold leaves no orphan loop: the stub `sh` is killed by the ctx, and `t.TempDir` removes the marker.

## 5. Docs

`docs/QUICKSTART.md:984–988` correctly scopes the refusal to the AILANG tools and says `workspace-exec` keeps working. It has
three inaccuracies (N4), none dangerous.

## 6. Gates (measured by me)

- `go vet ./...` rc=0. `gofmt -l host cmd`: empty.
- `go test ./host/daemon/ -count=1 -race -v -run '^(TestWorkspace|TestExecMaxTimeout)'` rc=0: **32 top-level `--- PASS`**, 0 FAIL,
  2 SKIP (pre-existing real-tool tests needing `WORLD_TOOL_AILANG_BIN`). All 7 new tests and the 4 amended ones print `--- PASS`.
- Full `go test ./host/daemon/ -count=1 -race` rc=0 (165 s).
- Commit boundaries C1 `44f222b` and C2 `3d3ee7b`: `go vet` plus the daemon TestWorkspace/TestExecSrt subset under `-race`, rc=0 each.
- srt trio `-race -p 1` (WORLD_EXEC_SRT_NODE_MODULES = iter-241 v0.0.78 install): 3/3 PASS.
- `host/store` `TestProductionContextRoots` PASS.
- **Broker verdict: pre-existing, unrelated, not caused by this diff.**
  - `git diff fe7c2a4 d8e9c72` is empty over `host/broker` and every package it depends on (store, canon, hashref, childenv,
    pkgproj, procbound). `go list -deps ./host/broker` does not include `host/daemon`, so the broker test binary is source-identical
    on base and head.
  - The executor's two red runs fail DIFFERENT timing tests: C0 baseline `TestRunBoundedHeadTailTimeoutKeepsPartialOutput`, head
    `TestExecHandlerTimeoutKillsTheGroup`, both partial-head-under-timeout. That is the row-166 flake class, appearing only under
    full-parallel `./...` load.
  - The `WARNING: DATA RACE` lines in those logs come from verifygate's deliberate `racecontrol/main.go` known-positive control,
    not from broker.
  - My isolated `go test ./host/broker/ -count=1 -race` at head: rc=0 (144 s).

## Findings

BLOCKING: none.

- **N1 (tests)**: the production budget's application is unpinned.
  - Reproduction: O4 (`workspace.go:455` `WithTimeout(…, workspaceHandlerBudget)` → `WithCancel`) survives the whole suite. The
    broker's internal 10 s `handlerExecTimeout` masks it, and `HandlerTimeoutError.Timeout` names that 10 s either way (F6/R1).
  - Effect: a regression that drops the 3 s budget would triple-plus the stall in a broken episode, `workspace-exec` included,
    with nothing red.
  - Suggested pin, without timing: assert that an error built under a budget below the broker cap reports a deadline from the
    parent, e.g. `errors.Is(err, context.DeadlineExceeded)`. Or fold this into follow-up R1, the timeout-label fix.
- **N2 (correctness, low)**: two builds for one episode with different sandboxes can run concurrently and share the policy path and
  cache dir. The older one finishing last overwrites `w.handler` with a stale root (`workspace.go:461–465`).
  - It is self-healing: the next caller's `cached.root` check misses and rebuilds, and `writePolicy` is atomic.
  - Effectively unreachable: `sandbox` is a pure function of a static config plus the symlink-free `epRoot`. O5 survives for this
    reason.
  - Worth a one-line comment, or a guard that skips the store when `w.ailangBuild[ep] != b`.
- **N3 (tests)**: at head, MUT-SHARED-LOCK is an equivalent mutant. Exec isolation from AILANG builds now rests on M2's "w.mu
  never held across the subprocess", which E7 pins, rather than on `execMu`. `execMu` is still justified (AC1.2, race isolation),
  but the record should not cite `…NeverWaitsOnAilangBuild` as killing MUT-SHARED-LOCK at head.
- **N4 (docs, `docs/QUICKSTART.md:985–988`)**:
  - (a) "discover" is an effect (`Ailang.Discover`), not a tool. The tools are `builtins-search`/`examples-search` (§9, `:376–377`).
  - (b) "refused only when the worktree itself is missing or unsafe" omits the typed refusal with no exec profile (§9 `:379`) and
    the "workspace-exec unavailable" construction failure.
  - (c) It does not tell the operator that in a persistently broken episode every call, `workspace-exec` included, waits up to 3 s
    first (design §11).
- **N5 (fragility, low)**: the 5 s hold-test budget is the only load-sensitive bound (false-red direction, 7–15× margin). See §4.

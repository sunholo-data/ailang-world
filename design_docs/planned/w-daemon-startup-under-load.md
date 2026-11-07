# w-daemon-startup-under-load — the registry phase gets its own startup clock, and a left-out exec profile shows in `doctor`

- **Status**: Planned; iteration 362 **r3**, bounded reviewer-specified refinement after r2
  quorum blocked 1/3; narrow-refinement carve-out applied (§11). Priority P1; target next host change, no release decision.
- **Scope**: World row 256, under D-WORLD-108. Row 233 remains parked on D-WORLD-105
  (OPEN, default C); row 257 remains separate. Ratified charter, logs and dashboard are unchanged.
- **Author**: `codex:gpt-6.1-sol` (pinned designer, World iter-362); prior r1 by
  `claude:claude-opus-5-5`. All three objections accepted (§11).
- **Estimated**: ~1 day including controller validation. Reuses the existing measured patch,
  changing diagnostic wording and its guards, with a caller inventory and cancellation arm.
  Patch: **11 files, +881/−20**, 1088 lines; SHA-256
  `90c6df45c6b4dcf09fd241d76eb87294f5b84d688eb9f2ff90d0eca1de09db16` (V27 and the revision receipt).
- **Base**: `61a8d1c47b45fffc0f36137d29a5285c32da0145`, worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter362-row256`.
- **Evidence**: `$V` means `design_docs/verification/world-iter361` (historical evidence,
  read fully; measurements inherited and **not re-measured** this revision). `$R` means
  `design_docs/verification/world-iter362` (new commands, actual dates, deadlines, loads and rc).
  New prototypes use `git archive HEAD` under `.design-scratch/`, with their own `git init`;
  production source is untouched. Scratch is removed before handoff. No commits or other Git
  writes outside those disposable repositories; no quorum here (controller owns it).
- **Validation environment**: PATH starts `/opt/homebrew/bin`; AILANG_BIN is
  `$HOME/.pinned-ailang/ailang`; TMPDIR and GOTMPDIR unchanged. No load experiment, secret access,
  external artifact seeding or shared HOME writes. Every gate uses the dated, bounded `$R/run.py`
  wrapper (at most 1800 s) and observes `uptime`. Socket tests are **UNINFORMATIVE UNDER SANDBOX**,
  including green designer tests. Controller V28 reran current bytes outside sandbox; Linux remains unmeasured.

---

## 1. Row 256 and the inherited measurements

| # | Claim | Verdict at `61a8d1c4` | Evidence |
|---|---|---|---|
| 1 | One deadline is shared by the probes and the registry bootstrap | **CONFIRMED in r1; inherited.** `Run` calls `startBounded(ctx, cfg, 9*time.Second, New)` (`daemon.go:1284`), which wraps **all** of `New` in one `context.WithTimeout` (`:1277`). Inside `New`, `configureRunCaps(ctx)` (`:817`), `configureExec(ctx)` (`:820`), `bootstrapRegistry(ctx)` (`:862`), `reportTemplates(ctx)` (`:868`) and `scanIntegrity(ctx)` (`:870`) all take that one ctx. A deadline in the bootstrap goes to `abortBudget(StageRegistry, …)` (`:937–938`) | V1 |
| 2 | The probes spend the budget the bootstrap needs | **TRUE BUT INCOMPLETE.** The home log of 2026-11-03 shows that **arm 1 of all five profiles** already read `(the startup budget is spent; stderr "")`. A probe run returns that text only when **no** time is left when it starts (`exec_probe.go:398–401`). So the clock was spent **before the first probe run** of any profile. What spent it was the work in front of the first probe run: the tool binary's archive, `node --version`, the srt tree and its archive (0.93–0.97 s warm on this rig up to the probe's start, V3), then the probe's own scratch setup. The probes then spent nothing, and the bootstrap found the clock empty | V1 (P13), V3 |
| 3 | `/v1/health` names left-out profiles, `doctor` does not | **CONFIRMED in r1; inherited.** `exec_profiles_excluded` is in the health body (`daemon.go:651`, filled at `:1134`). `cmd/ailang-worldd/doctor*.go` has **0** matches for `excluded`, `exec_profiles` or `ExecProfilesExcluded` (case-insensitive); the control `exec_tool_status\|ExecToolStatus` matches **4** times in the same files | V1 (P7, P8) |
| 4 | An exec call for a left-out profile should refuse with its name | **ALREADY BUILT** (row 189 M2). The refusal is `ExecProfileNotLoadedRefusal` (`handlers_exec.go:74`). Its non-test callers are `profileFor` (`workspace.go:195`, `:199`), which `execHandler` calls (`:486`) and turns into `broker.ExecRefusalHandler{Why: why}` (`:488`). That is the handler the episode registry binds to `Workspace.Exec`. `TestProfileForResolutionOrder` pins it (`binding_test.go:76`). This row builds nothing for it, and its new test **T2** drives it end to end from an excluded profile | V1 (P9), §5 T2 |
| 5 | The bind-time `broker: handler subprocess timed out after 10s` | **A DIFFERENT BUDGET, OUT OF SCOPE (§9).** That line comes from building the episode's AILANG handler. The build runs under `workspaceHandlerBudget` (3 s, `workspace.go:58`, its own root at `:607`). The "10s" in the text is the handler's own `ailangToolExecTimeout` (`handlers_ailang.go:54`), not the clock that fired. `runBounded` labels any expiry of its `runCtx` with `bounds.execTimeout` (`handlers.go:155`, `:208`, `:229–231`). A throwaway test on HEAD code (parent ctx 150 ms, handler bound 10 s) returned `broker: handler subprocess timed out after 10s` after **150 ms** | V1 (P10), V5 |
| 6 | A second start at load ~100 failed the same way; a third at load 16 succeeded | **Not re-measured.** The home log holds one `daemon startup failed at registry-bootstrap` line (line 36), not two. This session did not load the host | V1 (P13) |

**Historical evidence only.** V3 records phase timing at ambient load 7–13, and V4
records a slow-client reproduction without host stress. Neither identifies the cause of each
failed arm in the original incident. The old ~9.3× extrapolation is not a causal or scheduling
bound and is withdrawn as a sizing guarantee. No whole-startup wall-time bound has been
established. This revision changes the diagnostic claim while retaining the two context clocks.

## 2. The mechanism being changed

At `61a8d1c4`, `ailang-worldd serve` builds the daemon in one call (V1, V3):

1. `Run` → `startBounded(ctx, cfg, 9 s, New)`. One startup clock starts here.
2. `New`: `MigrateCredentials(ctx)`, `store.Open`, and the tool binary's archive (`archiveToolBinary`).
   The archive copies and hashes about 100 MB and takes no ctx (`archive.go:225`, `:278`).
3. `configureRunCaps(ctx)`, which is a no-op on the home service (no `--run-*` flag).
4. `configureExec(ctx)` → `broker.StartExec(ctx, …)`:
   - `node --version`, which is ctx-bounded;
   - the srt tree, then `ArchiveExecSandbox`, which takes no ctx;
   - `ProbeExecProfilesEach`, which probes every profile at once. Each probe run's budget is
     `min(execProbeRunBudget = 8 s, time left on ctx)`. A run that starts with nothing left returns
     `the startup budget is spent` (`exec_probe.go:396–402`).
   - `applyProbe`: a failed `--exec-profile` refuses startup. A failed directory profile is left out
     with `errs[i].Error()` as its reason (`exec_startup.go:256`, `:264–265`).
5. The interpreter's archive, about 100 MB, with no ctx (`daemon.go:830`).
6. `bootstrapRegistry(ctx)` on the **same** ctx. If it is past its deadline,
   `abortBudget(StageRegistry)` follows, then `daemon startup failed at registry-bootstrap: …
   context deadline exceeded`, and `serve` exits.
7. `reportTemplates(ctx)` and `scanIntegrity(ctx)` also run on the same ctx.

`service start` then polls `/v1/health` for 30 s (`service.go:87`, `:240`). It stops the service
when nothing answers: `build did not answer /v1/health within 30 s; service stopped` (`:256`).
The plist has `KeepAlive` and `ThrottleInterval 30` (V1 P11), so launchd cannot relaunch inside
that window.

**Historical reproduction, no deliberate host load** (V4, inherited positive control). A probe client that never answers
stands in for the load. In a `git archive 61a8d1c4` copy, a throwaway test built HEAD's own shape,
`startBounded(ctx, cfg, 1.5 s, New)`, with one directory profile:
`daemon startup failed at registry-bootstrap: cannot bootstrap world/epoch-registry/v1 with release
"unpinned": context deadline exceeded` after 1.91 s, and the log shows `left out: exec startup probe
failed … (7 of 7 arms)`. HEAD's real `Run` (9 s) failed the same way after 9.21 s. The control, the
same rig with an answering client and the 9 s clock, started in 2.50 s, with the profile left out by
its arms.

## 3. Design

### 3.1 What this row builds, per candidate

| Candidate | Decision | Why (measured) |
|---|---|---|
| (1) The registry bootstrap gets its own deadline | **BUILT.** It covers the whole *registry phase* (bootstrap, template report, integrity sweep), §3.2 | The bootstrap itself takes 0–2 ms on this rig (V3). Only the shared clock made it fail. The template report and the sweep run on the same ctx after it, so a fresh clock for the bootstrap alone would leave the sweep on a spent clock. The sweep then reports a false `scan` hole, which `MUT-SWEEP-ON-STARTUP-CLOCK` shows (§5) |
| (2) Retry a profile whose probe ran out of time | **NARROWED.** The probe records when failures were observed (§3.3); `/v1/health`, the log and `doctor` expose that observation without attributing cause. **No retry loop.** The loop is proposed as a follow-on row (§9) | A post-startup re-probe would have to make the startup exec state mutable. Today 10 functions read it and **no lock guards it** (`mu` guards `handler`/`ailangBuild`, `execMu` guards `execH`). `ExecToolState` also captures the profile set once, when it is built (`exec_tools.go:148`, `:161`) (V6). That is a concurrency change to the fail-closed gate's state, with a different correctness bar from the clock split. Recent rounds blocked docs that bundled two such bars. No retry is justified solely by the timestamp flag, which can also mark a genuine refusal |
| (3) Name left-out profiles in `/v1/health` and in `doctor` | **BUILT for `doctor`.** `/v1/health` already had them; it gains `out_of_time` (§3.4) | V1 P7, P8 |
| (4) An exec call for a left-out profile refuses with its name | **ALREADY BUILT** (row 189). Nothing is built; T2 pins the path from an excluded profile | §1 claim 4 |

### 3.2 Two startup clocks, both bounded

`Run` now calls `startDaemon(ctx, cfg, startupBudget, registryPhaseBudget)` (`host/daemon/daemon.go`, revised patch):

- **The startup clock**, `startupBudget = 9 s`, is HEAD's constant and value, unchanged. It is still
  made by `startBounded`. Pre-registry work receives that context; context-aware operations
  observe its deadline, while the archives remain context-free.
- **The registry clock**, `registryPhaseBudget = 9 s`, is new. `newDaemon` starts it from serve's
  own ctx (`phaseParent`) at the moment the registry phase starts:
  `regCtx, cancelReg := context.WithTimeout(phaseParent, registryBudget)` (in `newDaemon`).
  `bootstrapRegistry`, `reportTemplates` and `scanIntegrity` take `regCtx` (in `newDaemon`). Nothing
  before it shares it, so it never inherits time the probes, or the work in front of them, spent.
- `New(ctx, cfg)` keeps its signature for its **60 test references** (V19), as
  `newDaemon(ctx, ctx, registryPhaseBudget, cfg)`. When there is no separate parent, the registry
  clock is the smaller of the caller's deadline and 9 s. At baseline the only production
  `New` reference is `Run` passing its function value to `startBounded`; patched `Run` uses
  `startDaemon`, leaving zero production `New` references. T8 guards future callers, including
  imported aliases and function-value bindings. A new compatibility caller requires review.
- No new context root. Both clocks derive from serve's ctx with `WithTimeout`, so the
  context-root census (`TestProductionContextRoots`) is unchanged and green (§5).

**Why 9 s for the registry clock.** This preserves the former 9 s allowance when earlier
startup work was fast. V3 historically measured bootstrap 0–2 ms, template report 8–19 ms and
sweep 2–7 ms; these measurements do not establish an upper bound. `commitBudget` is 3 s (V20),
and `integrityScanTimeBudget` is 2 s plus page work; the separate 9 s phase is a bounded allowance,
not a derived worst-case latency guarantee.

**Whole lifecycle and cancellation.** The two context-aware phases have 9 s + 9 s allowances.
The tool binary, srt and interpreter archives remain context-free, as do other existing
non-context-aware setup steps. Scheduling is also unbounded here. No proof establishes a total
startup wall time below `service start`'s 30 s wait; that wait may still expire, warm or cold.
`startDaemon` must pass the original serve parent to `newDaemon`; canceling that parent must
reach the fresh registry phase. Construction failure still releases writer authority (or reports
settlement); successful `Run` still owns listen, serving and shutdown. T7 isolates fresh-phase
parent cancellation and writer release; the existing `TestStartupCallerCancellation` covers
`New` and `Run`, and T1 covers successful announce/serve/cancel/shutdown outside sandbox.

Listen, announce, serving and the shutdown watch bind serve's parent ctx (the same phaseParent passed to newDaemon), never the startup clock; a spent startup clock cannot end a successfully built daemon.
Here “bind” means the lifetime controlled by Run's parent cancellation watch and its maintenance
workers; Listen and announce themselves have no context argument. V29 traces the actual code;
T1 and MUT-SERVE-ON-STARTUP-CLOCK exercise the successfully announced service lifetime (AC9).

### 3.3 Timestamp classification establishes an observation

- `execProber.deadline` comes from the startup context. `fail` marks `spent` exactly when it
  records a failure at or after that deadline (a zero deadline never marks it).
- `ExecProbeError.OutOfTime()` is true for a nonempty set of failed arms only if **all** were
  marked `spent`. Any failure recorded before the deadline makes it false.
- `applyProbe` copies the flag into the exclusion and deletes the failed profile from the served
  set regardless. A failed explicit `--exec-profile` still refuses startup.
- The flag proves neither that timeout caused a failure nor that host load caused a timeout.
  A genuine refusal may be recorded late; a per-run 8 s timeout may be recorded before the
  startup deadline. There is **no established causal or scheduling bound**, including no bound
  on the delay between a failure and recording it. Both ambiguous directions remain (R3, R9).

The shared `broker.ExecProbeDeadlineObservation` string is exactly:

> Every failed arm was recorded at or after the startup deadline; startup timeout may have contributed.

The operator is told to review the failure reasons and restart if appropriate. The `out_of_time`
name is retained for compatibility with the existing patch; its documented semantics are only
these timestamps. It is unrelated to row 255's exec-call class `time-budget-spent`.

### 3.4 Where an operator sees it

- **Log**: the existing `exec profile "<p>" (<file>) left out: <reason>` line keeps its bytes.
  `logExecProfileExclusion` adds the shared observation string only for `OutOfTime`, followed by
  `Review the failure reasons and restart serve if appropriate.` T9 drives the production helper
  with a genuine refusal and a true flag, without sockets.
- **`/v1/health`**: `exec_profiles_excluded[]` gains `"out_of_time": true` (`omitempty`). The false
  form preserves existing exclusion bytes; true is additive (T10, V21). The field's advice and
  comments make only the timestamp claim. Profile validity is never inferred from it.
- **`doctor`**: the existing health GET also decodes exclusions. Each becomes a ✗ `exec profiles`
  line (exec class, exit 13), including with `--quick`, naming project, file and reason. The reason
  is one line and at most **240 bytes including the ellipsis**; project, file and fixed text are
  outside that cap. `oneLine(x.Reason, 240-len("…"))` reserves the three UTF-8 ellipsis bytes (V22).
  A true flag adds the shared observation and the fix `Review the failure reasons and restart
  World if appropriate (ailang-worldd service stop, then service start; or restart serve)`.
  Other exclusions advise reviewing the named file and failure reasons, correcting the profile
  if needed and restarting if appropriate. `doctor --help` lists the check.
- **Exec refusal**: unchanged: `exec profile "<p>" of episode "<e>" was left out at serve startup:
  <reason>`. No retries or widening of the served set.
- **S7**: `docs/QUICKSTART.md` explains the flag as an observation and advises review before restart.

### 3.5 Fail-closed

No path serves a profile whose probe did not pass every arm. The served set is `st.Profiles` after
`applyProbe`. `applyProbe` deletes every profile with a non-nil probe result, whatever `OutOfTime`
says. The flag is written only into `st.Excluded` and never reaches `st.Profiles` (§4 Claim C).

The registry clock changes when the daemon gives up, never what it serves. A flagged
`--exec-profile` whose probe runs out of time still refuses startup, as any failed `--exec-profile`
does (T4's explicit-profile arm).

### 3.6 Pure core, and why this is not a package

This is Go host code: the daemon's startup and the probe's bookkeeping. **No `.ail` is touched**,
and S1's Z3 contracts bind `world/` and `.ail` code. `OutOfTime` and `fail`'s rule are pure
functions of their inputs and the clock reading, and T4 tests them with no process run. A package
cannot hold this, because it is the host boundary's own startup (S2, S3).

## 4. The by-construction argument

These are structural properties of the patch, paired with current-byte controller V28/V29
outside-sandbox lifecycle/race and mutation evidence. Historical V11 remains dated history;
V28 independently replayed its instruments against the unchanged r2 patch.

**Claim A (serve's registry phase never inherits the startup clock).** *Proof.*

- Baseline callers are 60 test references and the sole production path `Run` →
  `startBounded(..., New)` (V19); patched `Run` uses the path below. No production compatibility
  caller remains. T8's AST guard fails if a caller appears outside the reviewed edges.

- `Run` calls `startDaemon(ctx, …)` (`host/daemon/daemon.go`). `startDaemon` passes `ctx` as `phaseParent`
  and the startup clock as `ctx` to `newDaemon` (its constructor closure).
- `regCtx` is `WithTimeout(phaseParent, registryBudget)` (in `newDaemon`).
- The three registry-phase calls take `regCtx` and nothing else (`:886`, `:892`, `:894`).

∎ T1 (`Serve`, through `Run`) kills `MUT-RUN-HEAD-SHAPE`. T1 and T2 (`Phase`) kill
`MUT-SHARED-PARENT`. T2 kills `MUT-SWEEP-ON-STARTUP-CLOCK`.

**Claim B (every clock is bounded).** *Proof.*

- The startup clock is `startBounded`'s `WithTimeout`, unchanged.
- The registry clock is `WithTimeout(phaseParent, registryBudget)`.
- Both values are 9 s, pinned by T6 (`Pin`).

∎ T3 (`Bound`) runs the registry clock at 1 ns and gets `registry-bootstrap … context deadline
exceeded`, which kills `MUT-REGISTRY-UNBOUNDED`. T6 kills both budget mutants.

**Claim C (fail-closed: a profile whose probe did not pass is never served).** *Proof.*

- In `applyProbe` every `i` with `errs[i] != nil` reaches `delete(st.Profiles, p.Project)`
  (`exec_startup.go:272`) unconditionally.
- `OutOfTime` is computed above that line and written only into `st.Excluded`.
- `w.exec.byProject` is `st.Profiles` (`workspace.go:260`).
- `profileFor` refuses any project in `excludedWhy` before it reads `byProject` (`:194–197`).

∎ T4 (`Spent`), T1 and T2 kill `MUT-SERVE-SPENT` (an out-of-time profile is skipped, so it stays
served) and `MUT-SPENT-STILL-SERVED` (it is left out but not deleted).

**Timestamp guarantee (replaces the withdrawn universal Claim D).** `OutOfTime()` returns
false if there are no failures or any failure was recorded before the deadline. A true flag
certifies only that every failed arm was recorded at or after it. T4 includes a genuine refusal
recorded after a fixed past deadline: true is permitted and exclusion is preserved. T4b, T5b
and T9 assert observation-only text and retained failure reasons. This is no causal classifier;
there is no established scheduling bound. `MUT-SPENT-ANY`, `MUT-LATE-NEVER-SPENT`,
`MUT-CAUSAL-OBSERVATION`, `MUT-LOG-CAUSE` and `MUT-DOCTOR-CAUSE` pin these obligations.

**Claim E (doctor names every left-out profile as an exec failure).** *Proof.* `doctorExecProfiles`
calls `r.add(markFail, "exec profiles", …)` once for each element of `exec_profiles_excluded`, and
`doctorClass("exec profiles")` is `("exec", 13)`. ∎ T5 (`Doctor`) kills `MUT-DOCTOR-NOT-CALLED`,
`MUT-DOCTOR-WARN`, `MUT-DOCTOR-CLASS` and `MUT-DOCTOR-ONE-FIX`.

**Claim F (neighbours untouched).** *Proof.* The patch's file list (V27) contains none of row 255's
files (`failure.go`, `record.go`, `broker.go`, `handlers_exec.go`, `effectful.go`, `why.go`,
`mcpclient.go`, `projection.go`). In `exec_probe.go` it changes only `ExecProbeArmFailure`,
`ExecProbeError.OutOfTime`, `execProber.deadline`, `fail`, the observation constant and two lines in `probeExecProfile`. No
verdict (`judgeReverse`, `judgeLeg`, `net6Leg`, `reverseLeg`, `classify*`) reads `deadline` or
`spent`. ∎ V12 historically records the full broker run. Current V28 validates its named socket-dependent neighbours; the full suite and Linux CI
remain acceptance obligations. V27 verifies patch scope, not dynamic verdicts.

## 5. Test plan and mutation table

| Id | Test (file) | What it asserts |
|---|---|---|
| **T1** `Serve` | `TestServeStartsWhenTheProbeSpendsTheStartupClock` (`host/daemon/startup_under_load_test.go`) | Through **`Run`**, with the startup clock at 1.5 s (test seam) and a probe client that never answers, spending the startup context allowance. Run announces within 12 s. `GET /v1/health` lists no loaded profile and one exclusion: `proj`, `out_of_time: true`, with a reason holding the probe's own `the startup budget is spent`. The log has the out-of-time line, and Run stops cleanly on cancel |
| **T2** `Phase` | `TestRegistryPhaseHasItsOwnClock` (same file) | Through `startDaemon`, two rows. **Spent-clock row** (1.5 s): the integrity sweep is complete with no hole; health shows `proj` left out with `out_of_time: true`; the episode `ep1`, mapped to `proj`, gets `ExecRefusalHandler` with `exec profile "proj" of episode "ep1" was left out at serve startup: exec startup probe failed`, so candidate 4 is reached end to end and nothing runs. **Control row** (9 s, the client answers, and the pass-through srt shim lets real arms refuse at once): the profile is left out with `out_of_time: false` |
| **T3** `Bound` | `TestRegistryPhaseClockIsBounded` (same file) | A 1 ns registry clock refuses startup with `StartupError{Stage: registry-bootstrap}` and `context.DeadlineExceeded`. The control (9 s) starts |
| **T6** `Pin` | `TestServeStartupClocksArePinned` (same file) | The production clocks are 9 s and 9 s |
| **T4** `Spent` | `TestExecProbeOutOfTimeNamesTheClockNotAnArm` (`host/broker/exec_probe_out_of_time_test.go`) | `fail` against a past, a future and a zero deadline. `OutOfTime` for all spent, one refused in time, a refusal only, and none; **genuine refusal recorded after a fixed past deadline ⇒ true**. `applyProbe` over four profiles: the out-of-time one is left out, flagged and not served; the refused one and a non-probe error are left out and not flagged; only the clean one is served, and `Only` is nil. A flagged `--exec-profile` with all failures recorded late still refuses startup |
| **T5** `Doctor` | `TestDoctorNamesLeftOutExecProfiles` (`cmd/ailang-worldd/doctor_exec_profiles_test.go`) | `doctor --quick` against a fake health with one out-of-time and one refused exclusion: exit 13, the ✓ daemon line, two ✗ `exec profiles` lines with the observation/review advice for the true flag and conditional profile review for the other, and `2 check(s) failed`. The control (nothing left out) exits 0 with no `exec profiles` line |
| **T4b** `Late` | `TestExecProbeGenuineRefusalRecordedAfterDeadline` (broker file) | Fixed `time.Unix(1, 0)` deadline, injects a sibling-write-success refusal fixture through the actual `fail` method: true flag permitted, reason retained, profile excluded, no default; shared diagnostic equals the observation-only sentence |
| **T5b** `AdviceCap` | `TestDoctorLateRefusalObservationAndReasonCap` (doctor file) | Direct JSON report capture for the same genuine refusal with false/true flags; review advice, no causal diagnosis, reason cap including ellipsis, one line and exec exit class |
| **T7** `Cancel` | `TestRegistryPhasePreservesParentCancellation` (daemon file) | Live startup ctx plus canceled phase parent refuses at registry-bootstrap with `context.Canceled`; settlement bounded, writer reopens |
| **T8** `Callers` | `TestDaemonConstructorCallerInventory`, `TestDaemonCallerInventoryAliasControl` (`startup_callers_test.go`) | Parse actual Go sources across repo, split test/production, named paths and owners; detect alias imports, function-value bindings and direct calls; reject unreviewed production construction edges; shadow control survives |
| **T9** `Log` | `TestExecExclusionLogIsObservationOnly` (daemon file) | Production logging helper preserves the genuine refusal and adds exactly the shared observation plus review/restart advice |
| **T10** `HealthJSON` | `TestHealthExecExclusionIsAdditive` (daemon file) | False flag preserves exact existing exclusion JSON; true adds the field |

**T4 truth table (recording time, not cause):**

| Failed arms | OutOfTime | Result |
|---|---|---|
| All recorded at/after deadline | true | exclude; observation/review advice |
| One recorded before deadline, others after | false | exclude; review reason |
| Genuine refusal recorded before deadline | false | exclude; review reason |
| Genuine refusal recorded after fixed past deadline | true | exclude; same observation/review advice, no causal attribution |
| Per-run 8 s timeout recorded before startup deadline | false | exclude; timeout still possible in the reason (R3) |
| No failures | false | no failed-profile exclusion |
| No deadline with a failure | false | exclude; review reason |

**Mutation history.** The following table is inherited iteration-361 evidence (V11), not a
claim that the revised bytes passed those experiments. Current V28 separately replayed all
19 named mutants outside sandbox: 18 killed, one disclosed template-clock survivor. The
MUT-NO-LOG-LINE instrument needed its moved helper anchor repaired, then compiled and killed
Serve. The original stale-anchor run is retained and is not counted as a successful experiment. The revised serial mutation harness is
`$R/mutate.py`; it asserts each edit changed production bytes, names the firing assertion,
and restores each file byte-exact. The new-obligation table follows the historical table.

| Mutant | Edit | Expected red | Observed |
|---|---|---|---|
| CONTROL | none | — | 15/15 PASS, MATCH |
| **MUT-RUN-HEAD-SHAPE** (HEAD's shared clock) | `Run` calls `startBounded(ctx, cfg, startupBudget, New)` | Serve | Serve, MATCH |
| **MUT-SHARED-PARENT** | `startDaemon` passes the startup clock as `phaseParent` | Serve, Phase | Serve, Phase, MATCH |
| MUT-REGISTRY-UNBOUNDED | `regCtx` from `WithCancel` (no bound) | Bound | Bound, MATCH |
| MUT-SWEEP-ON-STARTUP-CLOCK | `scanIntegrity(ctx)` | Phase | Phase, MATCH (a `scan` hole) |
| MUT-TEMPLATES-ON-STARTUP-CLOCK | `reportTemplates(ctx)` | — (predicted survivor) | **SURVIVES**, as predicted. No rig has `--ailang-bin`, so `reportTemplates` returns at `d.arch == nil` (`templates.go:97`). Its only effect would be a `capsule template status unavailable` log line (§6 R6) |
| **MUT-SERVE-SPENT** (fail-closed) | `applyProbe` skips an out-of-time profile, so it stays served | Spent, Serve, Phase | Spent, Serve, Phase, MATCH |
| **MUT-SPENT-STILL-SERVED** (fail-closed) | an out-of-time profile is excluded but not deleted | Spent, Serve, Phase | Spent, Serve, Phase, MATCH |
| MUT-FAIL-NEVER-SPENT | `fail` ignores the deadline | Spent, Serve, Phase | Spent, Serve, Phase, MATCH |
| MUT-DEADLINE-UNSET | `probeExecProfile` never sets `pp.deadline` | Serve, Phase | Serve, Phase, MATCH |
| MUT-SPENT-ANY | `OutOfTime` true when any arm is spent | Spent | Spent, MATCH |
| MUT-FLAG-NOT-CARRIED | `applyProbe` writes `OutOfTime: false` | Spent, Serve, Phase | Spent, Serve, Phase, MATCH |
| MUT-HEALTH-DROPS-FLAG | `healthProfiles` drops the flag | Serve, Phase | Serve, Phase, MATCH |
| MUT-NO-LOG-LINE | no out-of-time log line | Serve | Serve, MATCH |
| MUT-DOCTOR-NOT-CALLED | `doctorDaemon` never calls `doctorExecProfiles` | Doctor | Doctor, MATCH |
| MUT-DOCTOR-WARN | ✗ becomes ! | Doctor | Doctor, MATCH |
| MUT-DOCTOR-CLASS | `exec profiles` falls out of the exec class (exit 1) | Doctor | Doctor, MATCH |
| MUT-DOCTOR-ONE-FIX | `doctor` ignores `out_of_time` | Doctor | Doctor, MATCH |
| MUT-REGISTRY-BUDGET-30S | `registryPhaseBudget = 30 s` | Pin | Pin, MATCH |
| MUT-STARTUP-BUDGET-5S | `startupBudget = 5 s` | Pin | Pin, MATCH |

**Revision mutations (production edits; results V26).** Each row owns the correspondingly
named assertion in the T-test above; no grep discharges a load-bearing obligation.

| Mutant | Production edit | Required red |
|---|---|---|
| MUT-LATE-NEVER-SPENT | broker `fail` never marks spent | T4b: late refusal permits true |
| MUT-SPENT-ANY | broker `OutOfTime` uses any | T4: mixed recorded times false |
| MUT-LATE-STILL-SERVED | `applyProbe` fails to delete flagged profile | T4b: served set empty |
| MUT-FLAG-NOT-CARRIED | `applyProbe` drops flag | T4b: excluded entry flagged |
| MUT-CAUSAL-OBSERVATION | shared broker diagnostic asserts cause/profile validity | T4b: exact observation text |
| MUT-LOG-CAUSE | workspace log replaces observation with cause | T9: exact reason and observation |
| MUT-DOCTOR-CAUSE | doctor asserts cause/profile validity | T5b: observation text |
| MUT-DOCTOR-NO-REVIEW | doctor drops review advice | T5b: review before restart |
| MUT-REASON-CAP-240-PREFIX | doctor passes 240 without reserving ellipsis | T5b: reason ≤240 bytes |
| MUT-HEALTH-NO-OMITEMPTY | daemon removes `omitempty` | T10: exact false-form JSON |
| MUT-PARENT-CANCEL-DETACHED | registry timeout uses `WithoutCancel(phaseParent)` | T7: parent cancellation at registry-bootstrap |
| MUT-CALLER-IMPORTED-ALIAS | new production `var f = renamed.New` | T8: unreviewed production startup caller |
| MUT-CALLER-DIRECT | new production `renamed.New(...)` | T8: same assertion |
| MUT-CALLER-LOCAL-ALIAS | new daemon package `var f = New` | T8: same assertion |
| MUT-SERVE-ON-STARTUP-CLOCK | rebind Run lifecycle ctx to the construction deadline after successful construction | T1: announced daemon must answer health; AC9, V29 |

## 6. Residuals and risks

- **R1**: No retries. A failed profile remains excluded until a restart. The true flag is
  insufficient grounds for automatic retry, because the reason may describe a genuine refusal.
- **R2**: Context-free archives and setup, and host scheduling, have no established total bound.
  Either a warm or cold start can still miss the service's 30 s health wait. Historical V3
  records timings, not a guarantee or a load-causation experiment.
- **R3 (first ambiguous direction)**: A per-run 8 s timeout can be recorded while startup time
  remains, giving false. That does not prove a bad profile or rule out a timeout contribution.
  Reasons remain visible and must be investigated; exclusion is unchanged.
- **R4**: A failed explicit `--exec-profile` still refuses startup, whatever the flag.
- **R5**: A fresh registry allowance can defer refusal by up to its 9 s context budget for
  context-aware work. This is not a bound on all startup work or a guarantee about the 30 s wait.
- **R6**: The historical template-clock mutant survived because no rig archived an interpreter;
  that test gap remains explicit. Bootstrap and sweep have dynamic clock tests.
- **R7**: Linux and real Linux srt remain unmeasured. V28 discharges the named Darwin outside-sandbox
  socket runs; Linux CI is still required. Every sandbox socket result, green or red, is uninformative.
- **R8**: The non-answering client spends the startup context clock; it establishes no host-load
  cause and reproduces no CPU or I/O contention. No load experiment was run this revision.
- **R9 (second ambiguous direction)**: A genuine refusal may be recorded after the deadline and
  yield true. No causal or scheduling bound has been established, and no millisecond window is
  claimed. Advice must remain observational. T4b/T5b/T9 pin this case; it is always excluded.
- **R10**: Mutable `startupBudget`/`registryPhaseBudget` and existing rig hooks are serial test
  seams. No test may mutate them in parallel; narrow `-race` validation is mandatory. V28 discharges the successful Run lifecycle and named
  cancellation neighbours under race outside sandbox; V23 denials remain historical limitations.

## 7. Milestones

**M1 (single, ~1 day)** owns **AC1, AC2, AC3, AC4, AC5, AC6, AC7, AC8, AC9** and all files
listed below. Baseline each gate before applying; report preexisting failures rather than
changing the gate or environment to hide them. Controller runs socket gates outside sandbox.

For a fresh archive rehearsal (the designer's only application target), invoke `$R/fresh_apply.py`
through `$R/run.py`. It archives HEAD, runs `git init`, and executes these exact steps in that
archive; an enclosing repo can otherwise skip every path with rc 0 (historical V17):

```bash
git apply --check "$PATCH"
git apply "$PATCH"
test -f host/daemon/startup_callers_test.go
test -f host/daemon/startup_under_load_test.go
test -f host/broker/exec_probe_out_of_time_test.go
test -f cmd/ailang-worldd/doctor_exec_profiles_test.go
```

Each command below must run through the dated `$R/run.py` wrapper with a deadline ≤1800 s,
PATH=/opt/homebrew/bin:$PATH and AILANG_BIN=$HOME/.pinned-ailang/ailang; neither TMPDIR nor
GOTMPDIR is overridden. The wrapper records `uptime` for every gate. Use `go test -timeout 240s`:

1. `go build ./...`; `go vet ./host/broker/ ./host/daemon/ ./cmd/ailang-worldd/`; assert touched
   Go files are gofmt-clean (V24).
2. Run all T1, T2, T3, T4, T4b, T5, T5b, T6, T7, T8, T9 and T10 **serially under `-race`**.
   Run existing `TestStartupCallerCancellation`, `TestStartupBudgetQuarantinesPostCutoffBootstrap`,
   `TestStartupBudgetHeldConnectionSettles`, `TestHealthReportsBuildToolCredentialsAndProfiles`,
   `TestProfileForResolutionOrder`, `TestExecProfilesDirProbeFailureIsolatesOneProfile`,
   `TestDoctorDaemonVsForeignListener`, `TestDoctorHelpAndAPIKeyGuard`,
   `TestDoctorExecToolsWithAnM3Daemon` and `TestProductionContextRoots` as neighbours.
   Do not treat a missing new test as a green baseline. V23 gives exact command lines.
3. `$R/mutate.py` runs the revised obligation mutants after its unmutated control. Outside sandbox,
   re-run the historical clock/lifecycle mutants from V11 with current anchors, including
   `MUT-RUN-HEAD-SHAPE`, `MUT-SHARED-PARENT`, `MUT-REGISTRY-UNBOUNDED`,
   `MUT-SWEEP-ON-STARTUP-CLOCK`, `MUT-SERVE-SPENT`, `MUT-SPENT-STILL-SERVED` and
   `MUT-DOCTOR-NOT-CALLED`. Record the known template-clock survivor. Also run MUT-SERVE-ON-STARTUP-CLOCK (V29).
   V28/V29 need not be repeated when all 11 implementation hashes match their frozen manifest;
   full verify_go and PR/merge CI are still owed. V28 replay includes explicit build controls.
4. `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` in the controller's real checkout after its
   pristine baseline. In a `git init` archive it fails at `git ls-tree HEAD` on both baseline
   and patch (V25); do not seed history or call that a passed gate. No `.ail` file changes.

| AC | Owner | Required evidence / named red |
|---|---|---|
| AC1 | M1 | Fresh registry clock, bounded 9 s values, successful Run lifecycle; T1/T2/T3/T6; historical clock mutants above, controller outside sandbox |
| AC2 | M1 | All failed profiles excluded, including late genuine refusal; T4/T4b, MUT-LATE-STILL-SERVED and MUT-FLAG-NOT-CARRIED |
| AC3 | M1 | Observation-only log/health advice/doctor/docs; T4b/T5b/T9, MUT-CAUSAL-OBSERVATION, MUT-LOG-CAUSE, MUT-DOCTOR-CAUSE, MUT-DOCTOR-NO-REVIEW; S7 paragraph and help updated |
| AC4 | M1 | Nonempty ALL-after timestamp rule and T4 late-refusal truth row; MUT-SPENT-ANY, MUT-LATE-NEVER-SPENT; R3/R9 document both ambiguous directions with no scheduling bound |
| AC5 | M1 | Actual caller census and forward guard including imported/local aliases and direct calls; T8, all three MUT-CALLER-* edits |
| AC6 | M1 | Additive health JSON and ≤240-byte one-line reason; T10/T5b, MUT-HEALTH-NO-OMITEMPTY and MUT-REASON-CAP-240-PREFIX |
| AC7 | M1 | Parent cancellation, writer release and race-safe serial seams; T7 plus existing cancellation/settlement neighbours, MUT-PARENT-CANCEL-DETACHED, narrow -race including T1; current socket portion discharged by V28 |
| AC8 | M1 | All required tests passing on implementation/CI head; fresh-archive check/apply plus four presence assertions and digest V27/V29; sandbox denials explicitly uninformative; historical and new mutation controls, documentation updated |
| AC9 | M1 | Successful serve lifetime follows the original parent, never its spent startup clock; §3.2, T1, source trace and compiling MUT-SERVE-ON-STARTUP-CLOCK killed by T1 health assertion (V29) |

**Acceptance gate**: controller outside-sandbox gates and CI on the PR head. Linux acceptance
requires plain `go build + test` in CI; no sandbox green substitutes for it. There is no human
freeze decision; these are verification obligations. Controller owns routing, commits and quorum.

## 8. Files

Exact per-file added/deleted counts and new file lengths are banked in V27.

- `host/daemon/daemon.go`: split constructor, two clocks, Run path and additive exclusion flag.
- `host/daemon/workspace.go`: carry flag and observation-only exclusion logging helper.
- `host/broker/exec_probe.go`: record failure times, ALL-after flag, shared observation text.
- `host/broker/exec_startup.go`: carry flag only into exclusions; preserve deletion and flagged refusal.
- `cmd/ailang-worldd/doctor.go`: decode exclusions, observation-only advice, capped reason and help.
- `cmd/ailang-worldd/doctor_hosts.go`: exec profiles use class 13.
- `docs/QUICKSTART.md`: explain the observation flag and review/restart advice (S7).
- `host/daemon/startup_under_load_test.go`: T1/T2/T3/T6/T7/T9/T10 and serial helpers.
- `host/broker/exec_probe_out_of_time_test.go`: T4/T4b including deterministic late genuine refusal.
- `cmd/ailang-worldd/doctor_exec_profiles_test.go`: T5/T5b, long reasons and late-refusal advice.
- `host/daemon/startup_callers_test.go`: T8 AST inventory, forward guard and alias/shadow controls.

## 9. What this item is NOT doing

- **Not the retry loop** (candidate 2's second half). **Proposed follow-on row**,
  `w-exec-reprobe-out-of-time-profiles`: any future re-probe policy must review failure reasons; the timestamp flag alone
  is insufficient to identify retryable failures. Serving still requires a full pass.
  - Its premise is V6: it must make `workspaceExec` safe for concurrent update across 10 reader
    functions that no lock guards today, and rebuild or extend `ExecToolState`.
  - It needs its own design and quorum, because it changes the fail-closed gate's state machine.
  - No retry trigger is established here; `out_of_time` alone establishes no cause.
- **Not the bind-time `handler subprocess timed out after 10s`.** That is the AILANG handler's build
  under `workspaceHandlerBudget` (3 s), not startup, and it does not stop the daemon.
  - A failed build is not cached; the next Binder open builds again (`workspace.go:136–137`), so it
    recovers when load drops.
  - Its message names the wrong bound (V5). **Proposed follow-on row**, `w-handler-timeout-names-the-bound`:
    `runBounded` should name the bound that fired, the caller's ctx or its own `execTimeout`.
- **Not row 233** (D-WORLD-105 OPEN, default C). No reverse-leg verdict, retry or reset rule changes.
- **Not row 250's** forward legs, and **not row 242's** listener hoist.
- **Not row 255.** No change to the failure line, its classes, the record or its readers. The probe
  is startup, not an exec call.
- **Not row 257.** The seed is still cloned at the first call; `prepareCache` is untouched.
- **Not the archives' cost** (§6 R2): they still take no ctx.
- No change to `tools/launchd/*`, any skill, any `.ail`, `~/.ailang/state/*`, or `HOST_CHANGELOG.md`.

There are no decisions for the human. The deferral of the retry loop is a routing call, and §3.1
gives its measured reason. The controller may file the two follow-on rows, or not.

## Conflict Surface

| File | Function or symbol | Change | Neighbours it must not regress |
|---|---|---|---|
| `host/daemon/daemon.go` | `New`, new `newDaemon`, `startBounded` (unchanged), new `startDaemon`, `Run`, `startupBudget`, `registryPhaseBudget`, `HealthExecExclusion` | the registry clock; `out_of_time` | Row 23's bounded startup: `TestStartupBudgetQuarantinesPostCutoffBootstrap`/`…HeldConnectionSettles` (current V28 outside-sandbox race and mutation controls). Row 153's template report and the integrity announce (`TestIntegrityWarningsFollowTheAnnouncementAndServeNormally`, in the full daemon run, V12). The context-root census (`Roots`, V23). Row 189 M3's health fields (`Health`, V23) |
| `host/daemon/workspace.go` | `healthProfiles`, `configureExec`, `logExecProfileExclusion` | carry the flag; one more log line | Row 189 M2's `exec profile … left out` line, unchanged byte for byte. `profileFor` and `execHandler` (candidate 4) are untouched |
| `host/broker/exec_probe.go` | `ExecProbeArmFailure`, `ExecProbeError.OutOfTime`, `execProber.deadline`, `fail`, `probeExecProfile` | bookkeeping only | **Row 233** (reverse leg, D-WORLD-105): `reverseLeg`, `judgeReverse`, `pollLoopback` untouched. **Row 250** (forward legs): `net6Leg`, `judgeLeg`, `classifyNet6Dial`, `readProbeLine` untouched. **Row 242** (net6 hoist): `ProbeExecProfilesEach` untouched. `ExecProbeError.Error()` is unchanged, so every reason string and every test that matches it is unchanged |
| `host/broker/exec_startup.go` | `ExecProfileExclusion`, `applyProbe` | carry the flag | Row 189 M2's isolation (`Isolate`, baseline V23) and the flagged-profile refusal (T4) |
| `cmd/ailang-worldd/doctor.go`, `doctor_hosts.go` | `doctorDaemon`, new `doctorExecProfiles`, `doctorClass`, `doctorHelp` | the ✗ line | Row 138 M4's doctor contract: read-only (no new I/O beyond the existing health GET), the class order, `DocDaemon`, `DocHelp` and `DocM3` (current V28 historical-replay unmutated control) |
| `docs/QUICKSTART.md` | the profiles-directory paragraph | observation-only advice | No code block changes, and §10's refusal table is untouched |
| Compatibility constructor | `New(ctx, cfg)` and all actual references (V19) | baseline sole production reference `Run → startBounded(..., New)` becomes `Run → startDaemon → startBounded → newDaemon`; New remains for 60 test references | T8 rejects new production callers outside reviewed Run/startBounded edges and internal constructor edges, including import aliases, nested expressions and function values. No unobserved production caller is assumed |
| **Row 255** | `failure.go`, `record.go`, `broker.go`, `handlers_exec.go`, `effectful.go`, `why.go`, `mcpclient.go`, `projection.go` | **none** | the failure-line contract (200-byte cap, classes, `time-budget-spent`) |
| **Row 257** | `handlers_exec.go` `prepareCache`, `NewExecHandler` | **none** | the seed clone at first call |

## 10. Verification Log

Historical V1–V18 commands ran from `/Users/voightkampff/dev/sunholo-data/.wt-world-iter361-row256` unless they say
otherwise. A bare line number refers to `61a8d1c4`; a "patched" line number refers to the
prototype. `$V` = `design_docs/verification/world-iter361`. These are inherited measurements,
not fresh validation of r2. New V19–V27 rows cite `$R` and supersede the specified historical claims.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Premise readings P1–P14 at the base, each negative with a positive control | `bash $V/designer-premise-greps.sh.txt` → `$V/designer-premise-greps-out.txt` (`git show 61a8d1c4:<file> \| grep -n …`; `plutil -p` of the home plist; a read-only `grep` of `~/Library/Logs/ailang-world/worldd.log`) | P1 `daemon.go:1276 func startBounded`, `:1277 startupCtx, cancel := context.WithTimeout(ctx, budget)`, `:1284 d, err := startBounded(ctx, cfg, 9*time.Second, New)`. P2 `:811` `archiveToolBinary`, `:817` `configureRunCaps(ctx`, `:820` `configureExec(ctx`, `:830` `a.Archive(cfg.AilangBin)`, `:862` `bootstrapRegistry(ctx, release)`, `:868` `reportTemplates(ctx)`, `:870` `scanIntegrity(ctx)`. P3 `:937–938` `abortBudget(StageRegistry`. P4 `exec_probe.go:108 const execProbeRunBudget = 8 * time.Second`, `:398 budget = min(budget, time.Until(dl))`, `:401 "the startup budget is spent"`. P5 `exec_startup.go:256`, `:264 Reason: errs[i].Error()`, `:265 delete(st.Profiles`. P6 `BudgetSpent\|spent bool` in host: 3 files, all row 255's call-time class `FailTimeBudgetSpent` (`failure.go:33,40`, `handlers_exec.go:320`), none in the probe (control `Reason` in `exec_startup.go`: 4). P7 `daemon.go:651`, `:1134`. P8 doctor `excluded\|exec_profiles\|ExecProfilesExcluded` **0**, control `exec_tool_status\|ExecToolStatus` **4**. P9 `handlers_exec.go:74` def; callers `workspace.go:195`, `:199`, `:486`, `:488 return broker.ExecRefusalHandler{Why: why}`, `:572`; test `binding_test.go:76`. P10 `workspace.go:58 workspaceHandlerBudget = 3 * time.Second`, `:607`, `:664 runCapsStartupBudget = 10 * time.Second`, `:691`. P11 `service.go:87 wait: 30 * time.Second`, `:240`, `:256`; plist `KeepAlive => true`, `ThrottleInterval => 30`, `--exec-profiles-dir`. P12 `archive.go:75 probeTimeout = 10 * time.Second`, `:225 func (a *Archive) Archive(execPath string)` (no ctx), `:278 io.Copy(tmp, io.TeeReader(src, hasher))`. P13 home log: 5 `left out` lines, every one `7-of-7 arm1: … (the startup budget is spent;`; line 36 `daemon startup failed at registry-bootstrap: cannot bootstrap world/epoch-registry/v1 with release "AILANG v0.41.0": context deadline exceeded`; line 29 `workspace tools unavailable … handler subprocess timed out after 10s`. P14 the census pins `daemon.go\|drain\|Background` only |
| V2 | Instruments | `go version`; `grep '^go ' go.mod`; `echo $TMPDIR $GOTMPDIR $GOTOOLCHAIN`; `$HOME/.pinned-ailang/ailang --version`; `sysctl -n hw.ncpu`; `du -sh` and `find … -type f \| wc -l` of the srt tree; `ls -la` of both pins | `go1.26.6 darwin/arm64` = `go 1.26.6`; TMPDIR `/var/folders/kr/…/T/`, GOTMPDIR and GOTOOLCHAIN empty; `AILANG v0.41.0`, `Commit: 24ee108`; 16 cores; srt `16M`, 849 files; tool pin 106 877 442 B, interpreter pin 100 594 162 B |
| V3 | Phase timing of HEAD's `New` on this rig, real flags, scratch store | HEAD plus `PhaseMark` lines only (`$V/designer-phase-instrumentation.diff.txt`), built to `worldd-instr`; `bash $V/designer-phase-timing.sh.txt <label>` for `cold`, `warm1–3` → `$V/designer-phase-timing-out.txt`; loads 11.3–13.4. Each mark is printed *before* its step, so a step's time is the difference between its mark's `t` and the next mark's `t` (±1 ms against `since-prev`) | **cold**: tool archive 1153 ms, `node --version` 163, srt tree 199, srt archive 552, **probe phase 4797**, interpreter archive 1059, bootstrap 2, templates 8, sweep 7; New returns at **7981 ms (1047 ms left of 9 s)**. **warm1–3**: tool archive 216/245/232, node 197/168/135, srt tree 235/228/216, srt archive 297/276/331, **probe phase 4191/4268/4476**, interpreter archive 132/151/171, bootstrap 1/1/0, templates 8/19/9, sweep 3/2/2; New returns at **5311/5376/5610 ms**. Pre-probe work (to the `exec:probe-profiles` mark) 974/933/950 ms warm, 2106 cold. All five home profiles passed in every run. Wall to the listen line 5.6–6.1 s warm, 8.7 s cold |
| V4 | HEAD reproduction with no host load (positive control) | throwaway `zz_head_shape_test.go` (banked as `$V/designer-head-shape-repro_test.go.txt`) in a `git archive 61a8d1c4` copy (`head/`): the slow-client rig with one directory profile; `startBounded(ctx, cfg, 1.5 s, New)`; HEAD's own `Run` (9 s); control with an answering client and 9 s → `$V/designer-head-shape-repro.txt` | 1.5 s: `took=1.911s stage="registry-bootstrap" deadline=true err=daemon startup failed at registry-bootstrap: cannot bootstrap world/epoch-registry/v1 with release "unpinned": context deadline exceeded`; log `left out: … (7 of 7 arms): arm1-write-inside: … (broker: handler subprocess timed out after 1.055337667s; …`. Run 9 s: `took=9.206s`, the same error; log `(6 of 7 arms): arm1-write-inside: … timed out after 8s`. Control: `took=2.502s … err=<nil>`, left out `(4 of 7 arms): arm2-write-outside: … token "WORLD-PROBE-WROTE"` (the pass-through shim does not confine) |
| V5 | The bind-time "10s" names the handler's bound, not the clock that fired | throwaway `zz_label_test.go` (banked as `$V/designer-bind-label-probe_test.go.txt`) in `head/host/broker`: `runBounded` with a parent ctx of 150 ms, `execTimeout` 10 s, `sh -c 'sleep 2'` → `$V/designer-bind-label-probe.txt` | `took=150ms err=broker: handler subprocess timed out after 10s` |
| V6 | What a re-probe would make mutable (candidate 2's measured cost) | `bash $V/designer-retry-surface.sh.txt` → `$V/designer-retry-surface-out.txt` | 30 lines in non-test `host/daemon` touch `w.exec`/`workspaceExec` fields (3 writes in `configureExec`, the rest reads); **10 distinct functions** (`Close`, `Shutdown`, `profileFor`, `configureExec`, `episodeHandler`, `execHandler`, `healthExecToolStatus`, `healthExecTools`, `healthProfiles`, `captureExecTools`); locks `workspace.go:116 mu guards handler and ailangBuild only`, `:123 execMu guards execH only`, so **none guards `exec`**; `exec_tools.go:122 func NewExecToolState`, `:148 for project, p := range cfg.Profiles`, `:161` |
| V7 | The prototype builds, vets and is gofmt-clean | in `proto/`: `go build ./...`; `go vet ./host/broker/ ./host/daemon/ ./cmd/ailang-worldd/`; `gofmt -l` over the 9 touched Go files → `$V/designer-final-gates.txt` | `BUILD-OK`, `VET-OK`, `gofmt -l (touched): []` |
| V8 | The named tests pass | `go test -count=1 -v -run '^(<6 new>)$' ./host/daemon/ ./host/broker/ ./cmd/ailang-worldd/` | `--- PASS` ×6 (Serve 1.95 s, Phase 4.36 s with subtests 1.92 s and 2.48 s, Bound 0.01 s, Pin, Spent, Doctor); 3 `ok` |
| V9 | (pre-rename) the full packages pass with the first prototype | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 -timeout 900s ./host/broker/ ./host/daemon/ ./cmd/ailang-worldd/` | `ok host/broker 276.820s`, `ok host/daemon 306.095s`, `ok cmd/ailang-worldd 153.608s`, `RC=0`. A later rename changed only identifiers (`BudgetSpent` → `OutOfTime`, `budget_spent` → `out_of_time`). V12 re-ran the full packages on the final bytes |
| V10 | **SUPERSEDED by V22/V23 (r2)**: old doctor output used unsupported causal/profile-validity advice | Historical T5 instrumentation, r1 only; original transcript retained in `$V` | Withdrawn. Current output uses only recorded-after-deadline observations, timeout may have contributed, review reasons and restart if appropriate. No r1 diagnostic is a current requirement |
| V11 | The mutation table | `python3 $V/designer-mutate.py.txt $WT/.design-scratch/mut $V/designer-mutants.txt` (a copy of the final prototype); then `diff -r proto mut` | `CONTROL: rc=0 pass=15 red=- … MATCH`; the 19 mutants as in §5; `== done: 20 MATCH, 0 MISMATCH`; `MUT-TREE-RESTORED-IDENTICAL` |
| V12 | The final prototype: stability and the full packages | `$V/designer-final-gates.txt`: `go test -count=3` of the 6 new tests; `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 -timeout 900s ./host/broker/ ./host/daemon/ ./cmd/ailang-worldd/ ./host/store/` (load 13.9 at start) | `-count=3`: `ok host/daemon 19.823s`, `ok host/broker 0.241s`, `ok cmd/ailang-worldd 0.482s`. Full: `ok host/broker 273.719s`, `ok host/daemon 306.966s`, `ok cmd/ailang-worldd 156.223s`, `ok host/store 12.030s`. The daemon package holds the integrity-announce, context-thread and QUICKSTART-bound tests |
| V16 | Historical r1 patch — **SUPERSEDED by V27 for r2** | rsync the prototype onto a fresh `git archive 61a8d1c4` repo (`base/`), `git add -A`, `git diff --cached` → the patch; `git apply --check` in the worktree | `10 files changed, 530 insertions(+), 18 deletions(-)`; 741 lines; sha256 `f9bcb6360a38172ba99dcc41d7887a79b33646ad1b928f1aeaa2f890665c9b3e`; `apply-check rc=0`; `git status --porcelain` → only `.design-scratch/` (deleted after), the patch and `$V/` |
| V17 | Historical r1 fresh apply — **SUPERSEDED by V27 for r2** | `bash $V/designer-fresh-apply.sh.txt` → `$V/designer-fresh-apply.txt`: `git archive 61a8d1c4` into `fresh/`, `git init`, `git apply --check`, `git apply`, `go build ./...`, `go vet` (3 packages), `gofmt -l host cmd`, the 15 named tests | `APPLY-CHECK rc=0`, `APPLY rc=0`, the new test file present, `newDaemon` ×4, `BUILD-VET-OK`, gofmt empty, `go test rc=0`, **15 `--- PASS`, 0 FAIL**, 4 `ok`. **Run 1 was vacuous and is discarded.** Without `git init`, `git apply` inside the worktree's repo skipped every path and still returned rc 0, so only the 9 neighbours ran. The new-file check now guards that |
| V18 | No `.ail` change | `grep -c '^diff --git.*\.ail' <patch>`, with the control `'^diff --git.*\.go'` | `.ail` 0; `.go` 9 |
| V19 | Actual caller census, aliases and sole production path | `$R/baseline-caller-census-complete.txt`: inventory test-only fixture copied into the pristine archive; exact `go test` command banked, plus `TestDaemonCallerInventoryAliasControl`. Current path also in `$R/fresh-focused-race.txt`; actual source in `$R/premises.txt` | Baseline **60 test New references (all direct calls)**, **1 production New reference**, `host/daemon/daemon.go:1284 Run → startBounded(..., New)` (function value). Patched **60 test / 0 production New**; 4 production internal construction edges. Baseline scanner sees 623 Go files including its one injected test fixture (622 pristine); patched 626. Renamed imports, function-value aliases, nested selector expressions detected; shadow/field controls survive. Initial census instrument omitted baseline Run→startBounded from its allowlist (rc 1); corrected instrument rc 0, counts unchanged. Forward guard has three compiling production mutants (V26), not a negative grep |
| V20 | commitBudget is 3 s, distinct from startup clocks | `$R/premises.py` → `$R/premises.txt`, actual baseline source with lines: `daemon.go:146–150` and existing constant pin `daemon_test.go:228` | `commitBudget = 3 * time.Second` (B6). Registry clock is a separate 9 s allowance, not a renamed commit budget |
| V21 | HealthResponse permits additive data; exclusion false-form bytes preserved | `$R/premises.txt` reads `daemon.go:607–651`; T10 in `$R/fresh-focused-race.txt`; MUT-HEALTH-NO-OMITEMPTY in `$R/mutations.json` | Existing comment at :610 names additive health data; false flag yields exact old three-field JSON, true adds `out_of_time`. T10 passes; removing `omitempty` reds it |
| V22 | Doctor reason cap and noncausal review advice | `$R/premises.txt` reads the actual `oneLine` body and new doctor helper; T5b/T9/T4b in `$R/fresh-focused-race.txt`; cap/cause/review mutants in `$R/mutations.json` | `oneLine(s, n)` may append a 3-byte ellipsis after n bytes. New calls reserve those bytes with 237, guaranteeing reason ≤240 bytes; entire detail/path/fix are not claimed capped. Long/newline genuine-refusal fixture passes for both flags; true reports observations, retained reason and conditional restart. V10 is explicitly superseded |
| V23 | Narrow race, serial seams, cancellation and limitations | `$R/baseline-neighbours-race.txt`; `$R/race-detector-control.txt`; `$R/fresh-new-race.txt`; `$R/fresh-focused-race.txt` (exact dated commands/load/rc in each) | Baseline rc 1: doctor `httptest` IPv6 loopback bind denied; baseline daemon/broker/store neighbours green. Known-positive race control rc 1 with **WARNING: DATA RACE**, arming detector. Final focused run **rc 0, 14 named top-level PASS**, no race warning: 10 non-socket new tests + 4 existing cancellation/settlement/root tests. All 13 new tests requested under race: run rc 1, Run bind 127.0.0.1 denied and doctor bind [::1] denied; doctor's panic prevents T5b there, discharged separately in focused run. T2 green still uses sockets and is **UNINFORMATIVE UNDER SANDBOX**. No parallel seam mutations; mutable seams in T1 execute under race, but successful Run lifecycle requires controller outside sandbox |
| V24 | Fresh archive build/vet/gofmt baseline and patch | `$R/baseline-build-vet.txt` then `$R/fresh-build-vet.txt` | Both rc 0; `go1.26.6 darwin/arm64`. Exact touched file hashes in V27. No claim of full Go-suite acceptance |
| V25 | Prescribed plain verify_go baseline failure retained | `$R/baseline-verify-go.txt`; `$R/fresh-verify-go.txt`, `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` via bounded runner | Both **rc 1** after 24 tracked-binary instrument controls pass: archive's empty Git repository has no HEAD for `git ls-tree HEAD`. No history/index seeding or gate workaround. Controller must baseline/rerun in its real checkout; this result is not a production regression or pass |
| V26 | New obligations are load-bearing | `$R/mutate.py`; `$R/final-mutation-run.txt`; `$R/mutations.json`; per-mutant logs. Caller-specific rerun after AST nested-selector hardening is also in `$R/caller-mutation-final-run.txt` | Unmutated control rc 0 (8 tests). **14/14** named production mutants compile and red their named assertion (each rc 1), changed hashes and byte-exact restoration asserted. Three caller mutants also independently re-killed after hardening. Historical socket-dependent clock/lifecycle mutations were subsequently replayed by the controller in V28 |
| V27 | Revised patch, independent archive apply and presence guards | `$R/make_patch.py` (Python diff, no Git index); `$R/fresh_apply.py` → `$R/fresh-apply.txt`; `$R/prototype-manifest.json` | Baseline `61a8d1c47b45fffc0f36137d29a5285c32da0145`; **11 files, +881/−20**, 1088 patch lines; SHA-256 `90c6df45c6b4dcf09fd241d76eb87294f5b84d688eb9f2ff90d0eca1de09db16`. Check rc 0, apply rc 0, all three required test-file presence assertions rc 0, all 11 applied hashes equal the measured prototype; independent `.git` asserted by construction. 747 new-test lines (225 callers +314 daemon +108 broker +100 doctor); 7 existing files changed. Scratch removed after evidence banking |
| V28 | Current r2 bytes outside sandbox: lifecycle, races and both mutation sets | `$R/controller-v28/final-dated-race.txt`, `summary.json`, `historical-mutation-results.txt`, repaired `historical-no-log-repaired-results.txt`, `revision-replay/mutations.json` and all per-case logs | All 13 new tests plus four cancellation/settlement/context-root neighbours: **17 PASS, 0 FAIL, race rc 0** on Darwin arm64. All 19 historical mutant names replayed: 18 killed, one predicted template survivor; stale moved-log anchor corrected and re-run. All 14 revision mutants compile and kill their named assertions. Unmutated controls pass, all 11 file hashes restored. Historical 15-test controls include doctor daemon/help/M3 neighbours. No sandbox result is used for this acceptance evidence |
| V29 | Serve lifetime uses its parent; four new-file presence guards | `$R/controller-v28/serve-parent-context-source.txt`, `serve-parent-control.txt`, `serve-parent-mutant-build.txt`, `serve-parent-mutant.txt`; `$R/controller-v28` application receipt | Actual Run retains original ctx for maintenance, surface follow and cancellation watch; startupCtx exists only inside construction. Added mutation rebinds Run to the construction deadline after successful construction: control T1 race rc 0, mutant build rc 0, mutant T1 race rc 1 at health connection refusal. Files restored byte-exact. All four new test files, including startup_callers_test.go, present |
| — | **Still unmeasured / not accepted here** | — | Real host at load ~110; Linux/real Linux srt; `service start`/launchd end to end; second/third historical starts; full repo verify_go in a real checkout; template-clock mutation gap. No causal or scheduling bound established. Earlier prototype race log is superseded for its cancellation fixture: the first hook canceled during exec setup (config-stage refusal); T7 now isolates fresh-phase cancellation, and the older log is retained |


## 11. Quorum log

r1 (iteration 361, `$V/quorum-r1.json`): **BLOCKED 3/3**, no absent reviewers.
All objections accepted. GLM's caller census, production-path guard, missing premises and race
coverage are addressed by V19–V26. Kimi's false universal, one-direction conservatism, unsupported
latency window and missing apply-presence checks are removed/addressed (§3.3, timestamp guarantee,
T4/T4b, R3/R9, M1). GPT's observation-only alternative is used in the shared diagnostic, logs,
health comments, doctor, help, QUICKSTART and tests. No causal/scheduling-bound measurement is
invented. The count in the proposed fix was re-measured: **60**, not 47.

r2 (iteration 362, `$R/quorum-r2.json`): **BLOCKED 1/3**, all three seats present;
Kimi rejects, GLM and Gemini pass. Cost $0.3415408. OpenAI's vendor was excluded because the
author is codex:gpt-6.1-sol; Gemini replaces the former OpenAI reviewer, GLM and Kimi retained.

r3: bounded second revision under the already-ratified Gate-2 narrow-refinement carve-out;
**no third quorum and no claim that r2 passed**. Every remaining blocker has a concrete
reviewer-authored fix and none disputes the two-clock design direction. All are applied:

- **Kimi:** add current-byte outside-sandbox V28 for T1/T2/T3/T5/T5b and cancellation neighbours,
  under race, plus clock/lifecycle mutation replay; correct pending/current-evidence claims and
  row-23 conflict evidence. V28 contains the actual dated race control, all historical names and
  revised mutation replays with compile controls. Its proposed revision-only runner omitted old
  socket mutants, so the original historical instrument was adapted too. One stale anchor was
  corrected rather than counted green. No unsupported proof or MATCH claim is inherited.
- **Gemini:** add the verbatim fourth `test -f host/daemon/startup_callers_test.go` guard and
  change AC8 to four presence assertions. V29 records the controller's four-file check.
- **GLM:** insert its parent-context sentence verbatim in §3.2, with the actual Run lifetime
  mechanism clarified; add source trace and MUT-SERVE-ON-STARTUP-CLOCK (V29). New **AC9**, owned
  by M1, requires T1 and that named production mutation; it compiles and fails T1 health.

The source patch is unchanged at SHA-256
`90c6df45c6b4dcf09fd241d76eb87294f5b84d688eb9f2ff90d0eca1de09db16`.
These are reviewer-specified completeness/verification refinements, not a controller resolution
of a contested direction. Plan/executor gates, independent evaluator and CI remain mandatory.

## Axiom compliance and related documents

This host-boundary change has no compiler pipeline or syntax conflict surface. It reuses
`HealthResponse`, the existing probe verdict and the existing doctor report/class machinery.
There are no human design-freeze choices; the independent design quorum was performed for this
unattended row; its remaining concrete fixes were applied under the ratified carve-out (§11). Implementation details such as helper
names and evidence layout are agent-resolvable within the ACs.

| Axiom | Score | Reason |
|---|---|---|
| A1 Determinism | 0 | Timestamp observation is explicit host bookkeeping; pure core unchanged |
| A2 Replayability | 0 | No effect records or replay schema changes |
| A3 Effect Legibility | +1 | Failed profiles and diagnostic uncertainty are visible |
| A4 Explicit Authority | 0 | Exclusion and confinement verdicts unchanged |
| A5 Bounded Verification | 0 | Existing proof gates unchanged; context allowances remain finite |
| A6 Safe Concurrency | 0 | Startup state remains immutable after construction; serial seams checked under race |
| A7 Machines First | +1 | Additive structured exclusion flag and exec failure class |
| A8 Minimal Syntax | 0 | No language change |
| A9 Cost Visibility | 0 | No cost/KPI semantics change |
| A10 Composability | 0 | Existing report and health mechanisms reused |
| A11 Structured Failure | +1 | Doctor surfaces each exclusion with retained reason and review advice |
| A12 System Boundary | 0 | Host startup stays at the existing boundary |

Net +3; no hard violation of A1/A3/A4/A7. S1–S5: no `.ail` change; host startup cannot be
an extension package (§3.6). S6: named production mutations and presence assertions. S7:
QUICKSTART and help change with the operator surface.

Related: [World charter guardrails](../world-mission.md#guardrails-mission-specific-the-skills-standing-rules-always-apply-on-top),
[coding standards](../coding-standards.md), [DESIGN §1/§14](../DESIGN.md),
[blocked r1 quorum](../verification/world-iter361/quorum-r1.json). Existing startup clock,
probe isolation and doctor tests are retained; row 233/D105 and row 257 remain separate.

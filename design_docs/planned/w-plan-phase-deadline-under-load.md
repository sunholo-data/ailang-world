# w-plan-phase-deadline-under-load — the plan phase recompiles from cold on every call, under a 2 s cap nobody derived (row 153; closes row 161)

**Status:** PLANNED. Designed in iteration 241 (designer lane claude-opus-5-5, unattended). No quorum yet.
**Queue:** World mission row 153 `w-plan-phase-deadline-under-load` (clause-4, +6). It is the loop's next pick by D-WORLD-68 = A (Mark, attended 2026-10-07). **This design also closes hygiene row 161** (`w-exec-srt-mcp-e2e-flake`). Both of 161's signatures are measured to be this defect (§4).
**Estimated:** ~1.5 d of host Go and tests (M1 ~0.4 d, M2 ~0.15 d, M3 ~0.7 d, M4 ~0.25 d). **Dependencies:** none. No `.ail` changes, no pin changes, no upstream gate. Row 159 (`#1602`) only affects how much of M1 reaches the MCP wire (§5a).
**Measured base:** worktree `sprint/row153-plan-phase-deadline-under-load` at `990a6da`. Gate interpreter `~/.pinned-ailang/ailang` = `AILANG v0.41.0` (commit `24ee108`, sha256 `1a67b014…`, 100,594,162 bytes). `go1.26.6 darwin/arm64` on an Apple M4 Max (12 performance + 4 efficiency cores). Real srt `@anthropic-ai/sandbox-runtime@0.0.78`, installed for this design under `~/.ailang/state/mission-world-iter241/srt/`, with `cli.js` sha256 `3c3092bd…` equal to `design_docs/verification/world-row140-m0/pin.json`. Probe artifacts, logs and the uncommitted timing patch are banked under `~/.ailang/state/mission-world-iter241/` (V-rows cite file names). **The probe patch was reverted; this commit touches only this document.**

## 1. Problem

An effectful transition runs plan → effect → (finish) → commit (`host/coordinator/effectful.go:248–316`). The plan phase is one capsule run of the transition's `.ail` in the pinned interpreter, under `PlanPhaseBudget = 2 * time.Second` (`effectful.go:37`). Under CPU contention the phase overruns that cap. The coordinator returns `fmt.Errorf("coordinator: execute plan phase: %w", context.DeadlineExceeded)` (`:135`). The MCP adapter passes it on unchanged (`host/projection/mcp.go:126–131`). The released `mcphttp` seam then turns any `DeadlineExceeded` into `-32603 host callback timed out` (V3). So the client is told the *callback* timed out when the call still had ~18 s of its 20 s `invokeDeadline` left (V5). The real cause was a 2 s sub-budget on a pure, effect-free step.

In CI this fails `TestExecSrtMCPEndToEnd` and `TestExecSrtMCPGrantAndBudget` in `verify_go.sh`'s parallel `go test` legs on most attempts since ~14:00Z 2026-10-06: 11 of the 13 job attempts fetched for this design were red (V23). The same tests pass in the dedicated `-p 1` srt step of the same runs. Every landing pays reruns.

Measured: **where the time goes.** Idle, the plan phase is ~0.2 s, and ~75% of it is the interpreter child (V7). **The child cold-compiles the transition and 7 `std/*` modules on every call.** The interpreter keeps a compile cache beside the source (`.ailang/cache/compile/`) or under `AILANG_CACHE_DIR`. The capsule stages a fresh temp root for every run and deletes it afterwards (`host/capsule/capsule.go:198–202`), so every phase compiles from cold (V17). Cold costs 5–7× warm, at every load measured: 171–181 ms vs 25–30 ms idle, and 2.36–2.52 s vs 0.41–0.56 s on the throttled 8-burner proxy (V18). The rest of the phase is a full-file SHA-256 re-verification of the 100 MB interpreter on every run, 18–30% of the phase (V6, V7, V15). Neither is the store or a callback queue: the store reads before the plan take 0.08–0.32 ms (V7), and `procbound.Admit` never blocks (V6).

## 2. Measured facts (summary)

| # | Claim | V |
|---|---|---|
| F1 | The plan cap is 2 s. Its only stated basis is the row-134 design's assertion "plan capsule run ≤ 2 s". No derivation or measured budget was recorded. | V1, V27 |
| F2 | A plan-cap expiry surfaces as MCP `-32603 host callback timed out` and A2A `invocation exceeded its deadline`. Both name the call deadline, not the phase. | V2–V5 |
| F3 | Idle split (16-core rig): plan 193–227 ms = verify 37–45 ms + child 152–179 ms + staging ≤ 2 ms. Store reads before the plan: ≤ 0.33 ms. | V7, V8 |
| F4 | The child is ~85% cold compile. In the same root, with the cache present, the same exec drops from 168 ms to 26 ms. | V17, V18, V26 |
| F5 | The failure reproduces locally under a throttled-CPU proxy (background QoS on 4 efficiency cores, plus N CPU burners). With 4 burners the plan takes 1.40–1.71 s and passes. With 5–8 burners it overruns 2 s, giving signature 1 and signature 2. | V11–V14 |
| F6 | Uncapped at 8 burners, the plan takes 2.42–3.72 s (n = 10), and every call that reached the plan succeeded. At 12 burners the rig fails *before* any plan: the startup probe and publication break first. | V15, V16 |
| F7 | Signature 2 (`where = "<nil>"`) is the same plan-phase timeout on the test's second call. The test never reads that call's `wire.Error`. `database is closed` comes from a different, passing test. | V14, V24, V25 |
| F8 | A plan timeout happens before any durable write, so resending is safe: A2A can reuse the same task id, MCP just makes a new call. | V28 |
| F9 | A warm template removes the cold cost. `ailang check` at the capsule's staging path fully pre-warms `run` (27–32 ms vs 171–174 ms). v0.41.0 honours `AILANG_CACHE_DIR`. Warm and cold stdout are byte-identical. | V19–V22 |

## 3. Verification Log

Shorthand: `B=~/.pinned-ailang/ailang` (v0.41.0); `S=~/.ailang/state/mission-world-iter241`; `W=` this worktree. Unless a row says otherwise, Go probes ran with `PATH=/opt/homebrew/bin:$PATH AILANG_BIN=$B WORLD_EXEC_SRT_NODE_MODULES=$S/srt/node_modules`. Timing rows used the uncommitted probe patch `$S/probe-phase-timing.patch`, enabled by `WORLD_PROBE_PHASE_TIMING=1`. It logs, per capsule run, `resolve / verify / stage / start / child / total`, and per phase `took` and `pctx.Err()`. Every timing row is descriptive of this rig, not an acceptance threshold. **Rule applied:** each null result is paired with a positive control on the same rig.

| V | Command (trimmed) | Reading | Control / scope |
|---|---|---|---|
| V1 | `sed -n 32,42p host/coordinator/effectful.go` | `PlanPhaseBudget = 2s`, `FinishPhaseBudget = 2s`, `HandlerCap = 10s`, `HandlerHeadroom = Finish + PostEffect`, `PostEffectBudget = 4s` | static, base `990a6da` |
| V2 | `sed -n 127,141p host/coordinator/effectful.go` | `runPhase`: `pctx := WithTimeout(ctx, budget)`; on error, `if pctx.Err() != nil` → `fmt.Errorf("coordinator: execute %s phase: %w", …, cerr)`. Plan call sites `:260` (dispatch) and `:387` (replay). | static |
| V3 | `grep 'sunholo-data/ailang ' go.mod`; `sed -n 21,32p $GOMODCACHE/…/ailang@v0.47.2/serveapi/protocol/envelope.go`; `grep -n CallbackMessage …/mcphttp/handler.go` | pin `v0.47.2`. `CallbackMessage`: `errors.Is(err, context.DeadlineExceeded)` → `"host callback timed out"`, written with code `-32603` by `WriteMCPEnvelope`. Called at `handler.go:118,125,182`. | Same file maps `Canceled` → `"host callback canceled"`, so the switch is live. Agrees with row 136 / ailang#1602: only 4 fixed strings reach the MCP wire. |
| V4 | `sed -n 416,470p host/projection/projection.go` | `dispatchError` (A2A only, called at `:372`): `case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded)` → `"invocation exceeded its deadline"` (`:461`) | static |
| V5 | `grep -n invokeDeadline host/daemon/daemon.go` | `invokeDeadline = 20 * time.Second` (`:99`), wired as `InvokeWait` and `CallbackTimeout` (`:704–705`) | pinned by `host/daemon/mcp_test.go:165` |
| V6 | `sed -n 189,256p host/capsule/capsule.go`; `sed -n 43,55p host/procbound/procbound.go`; `stat -f %z $B` | Order: `Resolve` → `verifyExecutable` (`os.ReadFile` + hash of the whole binary, `:194`, `:359–371`) → `MkdirTemp` fresh root (`:198`) + `defer RemoveAll` (`:202`) → `WithTimeout(parent, execTimeout)` (`:221`) → `Env = {AILANG_FS_SANDBOX, AILANG_RELAX_MODULES}` (`:228`). `procbound.Admit` is a CAS counter that fails fast and never blocks. Binary 100,594,162 bytes. | static |
| V7 | `go test ./host/daemon/ -count=3 -v -run '^(TestExecSrtMCPEndToEnd\|TestExecSrtMCPGrantAndBudget)$'` (probe on) → `$S/idle-probe.log` | 12 plans: `took` 193–227 ms; `verify` 37–45 ms; `child` 152–179 ms; `resolve` ≤ 13 µs; `start − stage` ≤ 2 ms; `loadSourceAndWorld` 0.08–0.32 ms. All 6 tests PASS. First call 399–467 ms. | Real srt and real interpreter: 0 skips (both tests print `--- PASS`) |
| V8 | same as V7 with `-race -count=2` → `$S/idle-race-probe.log` | plans 220–246 ms; verify 52–63 ms; child 159–179 ms | `-race` costs the daemon ~+15 ms of verify; the child is not race-built |
| V9 | `$S/suite_load.sh`: full `go test ./... -count=1 -race -timeout 12m` in the background (CI's race-leg command), with the two e2e tests looped `-race` until it ends → `$S/suiteload1.log{,.suite}` | suite `SUITE_RC=0` (all 27 packages ok). 15 loops: **30 PASS, 0 FAIL**. 60 plans 233–278 ms. Load average 15–20. | **Null on a 16-core host.** The positive control is V13/V14 on the same tree and binary. This rig cannot reach CI's per-core contention with real work alone. |
| V10 | `$S/stress_run.sh 64 5`: 64 `yes` burners (default QoS) + `-count=5` (no race) → `$S/stress64.log` | 10 PASS. Plans 237–448 ms. Load average reached 65. | macOS spreads default-QoS work across 12 P-cores, so the plan barely moves (≤ 2×) |
| V11 | `taskpolicy -b $S/daemon.race.test -test.count=2 -test.run '…'` (background QoS → the 4 efficiency cores, throttled) → `$S/bg-idle.log` | 8 plans **1.09–1.25 s**: verify 231–271 ms, child 812–959 ms. First call 2.37–2.44 s. 4/4 PASS. | Proxy for a slow machine: ~5× idle |
| V12 | `$S/bg_load.sh 4 3` (V11 + 4 burners under `taskpolicy -b`) → `$S/bg4.log` | 12 plans **1.40–1.71 s**. 6/6 PASS. First call 3.48–3.50 s. | Comparable to CI's signature-2 first calls of 2.37–3.36 s (V23) |
| V13 | `$S/bg_load.sh 8 3` → `$S/bg8.log` | **Signature 1 reproduced**: `took=2.03–2.96s err=context deadline exceeded`, client `{Code:-32603 Message:host callback timed out}` after 2209–2275 ms. 4 of 6 tests reached a plan and all 4 failed there; 2 failed earlier at publication (V16). | Positive control for V9/V10 |
| V14 | test temporarily logs the second call's `wire.Error`; `RUNPAT='^TestExecSrtMCPEndToEnd$' $S/bg_load.sh 5 8` → `$S/bg5.log` | 7 FAIL / 1 PASS. **Signature 2 reproduced once**: first call OK (plan 2.00 s, call 4384 ms), second call `PROBE second call error: {Code:-32603 Message:host callback timed out}` with `phase=plan took=2.035s err=context deadline exceeded`, then `where = "<nil>"`. The other 6 failures were signature 1. | Same mechanism as signature 1 |
| V15 | probe patch also sets `PlanPhaseBudget = 20s`; `$S/bg_load.sh 8 4` → `$S/bg8-uncapped.log` | 10 plans **2.42–3.72 s**, all `err=<nil>`; verify 376–660 ms; child 1.96–2.98 s. Every test whose plan ran PASSED (5). 3 failed earlier at publication. | Distribution of the overrun the 2 s cap cut off in V13 |
| V16 | `$S/bg_load.sh 12 3` (uncapped) → `$S/bg12-uncapped.log` | 6/6 FAIL **before any plan**: 2 × `exec startup probe … broker: handler subprocess timed out` (10 s cap), 4 × `PublishSet … is not loadable under its pinned interpreter` (output truncated; publication runs under the test's 30 s `boundedTestContext`, `host/daemon/m7b_deadline_test.go:12`) | Above ~8 burners the plan stops being the first component to fail. `grep -l 'not loadable' $S/runs/*.log` → 0: **CI has never shown these.** |
| V17 | `cd $S/childrun.nm6q && env -i AILANG_FS_SANDBOX=$PWD AILANG_RELAX_MODULES=1 $B run --quiet --caps "" --entry main --args-file args.json host/capsule/main.ail` (exec.ail staged at the capsule path); `find . -newer args.json` | The interpreter writes `host/capsule/.ailang/cache/compile/` (manifest `"version":"v4"`, 8 entries: `host/capsule/main` + `std/{bytes,json,list,math,option,result,string}`; 41 files, 1.0 MB) | `$B run --help` names `Compilation cache: .ailang/cache/compile/` |
| V18 | `python3 -I $S/coldwarm.py <exec.ail> <args.json> $B` (7 runs per arm; cold arm deletes `.ailang` first); same with 8 `taskpolicy -b yes` burners | fg cold **171–181**, warm **25–30** ms. bg cold **576–822**, warm **98–152** ms. bg+8 cold **2360–2523**, warm **407–564** ms. Every run asserted rc 0 and a `"plan"` answer. | Warm/cold ratio 4.4–7× at every load |
| V19 | `python3 -I $S/stdonly.py …` | std-only template (the main entry removed): 73–83 ms. Full template copied into the cache dir before each run: 25–59 ms | std is ~100 ms of the ~150 ms cold cost; the source module is ~50 ms |
| V20 | `python3 -I $S/checkwarm.py …`: `$B check host/capsule/main.ail` in a fresh root, then `run` | after `check`: **27–32 ms**; cold: 171–174 ms | `check` executes nothing and fully warms `run` |
| V21 | `run` with `AILANG_CACHE_DIR=$c`; count `manifest.json` under the root and under `$c` | root 0, `$c` 1 | v0.41.0 honours `AILANG_CACHE_DIR` for the compile cache (`strings $B` also contains it) |
| V22 | same staged root, `run` twice (cold then warm), `shasum -a 256` of each stdout | both `4b6c4d08c6da6b91…` | n = 1 input. M3 widens this to every se-tools plan fixture (AC3.3). |
| V23 | `gh api …/actions/runs/<run>/attempts/<a>/jobs`, then `gh api --allow-escape-sequences …/jobs/<id>/logs`, ANSI-stripped, into `$S/runs/` (13 job logs, each 92–161 KB) | runs 37389860178, 37475959991 (a1–a4), 37497034244, 37499638889 (a1, a2), 37504535449 (a1, a2), 37504759449 (a1, a2), 37509957933, 37509997919: **11 red, 2 green**. Red: sig-1 ×3 (`exec_e2e_srt_test.go:126 … host callback timed out`, first call 2040–2060 ms), sig-2 ×7 (`:154 where = "<nil>"`), GrantAndBudget ×1 (`:197 first call … Error:0x…`, a pointer, **in the non-race leg**). Sig-2 first calls: 2371–3362 ms with handler `duration_ms` 407–966. The `-p 1` step's first calls in the same jobs: 503–930 ms with `duration_ms` 132–324. | Two green attempts (`112401768092`, `112418771019`) are the CI control. No CI log contains `plan phase`: see V25. |
| V24 | `go test -race -c …; daemon.race.test -test.run '^TestHeadErrorsUseAPIEnvelope$' -test.v` | prints `ailang-worldd: internal error: GET /v1/head: store: read selected head: sql: database is closed`, then `--- PASS`. The subtest closes the store on purpose (`host/daemon/handlers_test.go:250–263`). | The lines after it in CI (`host-secret-detail`, `private commit read detail`) are other tests' fixtures. **Signature 2's "database is closed" is not from the failing test.** |
| V25 | `sed -n 143,148p host/daemon/workspace_test.go`; `sed -n 150,157p;195,198p host/daemon/exec_e2e_srt_test.go` | The e2e rig sets `cfg.ErrorLog = io.Discard`, so the daemon's `mcp refusal … execute plan phase` line never reaches CI. The second call ignores `wire.Error` and reads only `StructuredContent["stdout"]`. `:197` prints `Error:0x…` (the pointer, not the message). | Explains why CI showed two signatures for one cause |
| V26 | temporary `host/capsule/zz_probe_timing_test.go` (banked as `$S/zz_probe_timing_test.go.txt`): 3 × `Runner.Run`, then direct `exec.Command` of the archived binary and of `$B` in one fixed root | capsule runs 214–221 ms each (fresh root every time). Direct exec in a fixed root: 1st 168 ms, then 25–31 ms, with or without `--caps ""`. | The gap is the cache, not exec mechanics, binary path or flags |
| V27 | `grep -n 'plan capsule run' design_docs/planned/w-software-engineering-domain.md`; `grep -rn 'AC0.3' design_docs --include='*.md'` | `:181` asserts "plan capsule run ≤ 2 s". `:248` asks for "median of 5 capsule plan runs … as a V-row". The grep finds no row-134 AC0.3 result anywhere (hits are other rows' AC0.3). `git log -S'PlanPhaseBudget   = 2'` → introduced in `f3a90ee` (row 134). | The 2 s was never derived |
| V28 | `sed -n 262,300p host/coordinator/coordinator.go`; `grep -n AppendIntent host/coordinator/coordinator.go` | Before the plan, Dispatch only reads (`GetReceipt` `:288`). The first durable write is `AppendIntent` in `appendAndCommit` (`:386`), after the plan. `inFlight` is deleted on return (`:271`). | A timed-out plan leaves no receipt, so the same task id re-executes |
| V29 | `npm install @anthropic-ai/sandbox-runtime@0.0.78` in `$S/srt`; `shasum -a 256 …/dist/cli.js` | `3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2c7cd40d96e` = pin | The real-srt path runs on this macOS rig (V7 PASS, not SKIP) |
| V30 | job log header; `sed -n 316,320p scripts/verify_go.sh` | `Image: ubuntu-24.04`. The repo's own note: "on the 2-core runner … with -p 2 packages run two at a time" | **Core count not re-measured** (UNMEASURED U2) |

## 4. Root-cause verdicts

**Signature 1** (`-32603 host callback timed out` on the first call, ~2.04–2.06 s). This is the 2 s plan cap expiring (V13; CI V23 matches to the 20 ms). The overrun is CPU time spent recompiling `std/*` and the transition from cold on every call (V17, V18), plus re-hashing 100 MB (V6). The store and the callback queue play no part (V7, V6).

**Signature 2** (`where = "<nil>"` + `database is closed`). **Same root cause**, measured in V14. The first call's plan squeaked under 2 s; the second call's did not. The test reads `StructuredContent` without checking `wire.Error` (V25), so the timeout shows up as a nil field. `database is closed` is a red herring: a passing test prints it on purpose, and `go test` prints a failing package's whole stderr (V24). It is **not** a teardown/ordering race. **Row 161 is closed by this row** (M1's rig fix plus M2/M3), and the charter should mark it so.

**The GrantAndBudget variant** (V23, the non-race leg). Same shape: first call error, plan-bound. `-race` is not a precondition (V9 vs V13: the race-built daemon adds ~15 ms; the child is never race-built).

**Why CI and not M4's serial runs.** Row 93 M4 ran 70 serial task-runs with 0 timeouts. That is consistent: idle plans are 10× under the cap (V7). Only CPU starvation of the 2-core runner crosses it. In this proxy that means ≥ 5 throttled burners (V12–V14).

## 5. Options

### (a) The error the client sees

- **a1 (chosen). A typed `coordinator.PhaseTimeoutError{Phase, Budget, Elapsed}`.** `runPhase` returns it when the phase context expired **and the caller's did not**. When the caller's context is also done, the caller's deadline is reported unchanged (that is a call timeout, not a phase timeout). `Unwrap()` returns `context.DeadlineExceeded`.
  - **A2A (World-owned, client-visible now):** `dispatchError` gains a case *before* the generic deadline case. It returns `"transition plan phase exceeded its 4s budget; nothing ran or was committed; resend the same task id"`, with an exported prefix constant. This is retryable by construction (V28). A finish-phase timeout never reaches this case: it already arrives wrapped in `EffectsUnrecordedError`, which `dispatchError` matches first (`projection.go:438`).
  - **MCP (honest limit):** the released seam carries only `CallbackMessage(err)` (V3). Keeping `DeadlineExceeded` in the chain keeps the wire at `host callback timed out`. Of the four fixed strings, that is the only true one; without the unwrap it would become the less true `host callback failed`. The phase name reaches the operator line (`mcp refusal … plan phase exceeded its 4s budget after 4.01s`) now. It reaches the MCP client only when row 159 lands on ailang#1602. M1 adds a frozen-wire assertion that row 159 will flip (the same tripwire pattern as `TestSeToolsMCPPostEffectFailureIsLogged`).
- a2. A new JSON-RPC code for phase timeouts on A2A. Rejected: every other `dispatchError` outcome uses `codeInternal` or `codeInvalidParams`, and A2A clients key on the message. A one-off code is a protocol change with no consumer.
- a3. Wait for #1602 before changing anything. Rejected: A2A and the operator line can be honest today.

### (b) The bound and the contention

- **b1 (chosen, M2). Derive the plan cap from the deadline arithmetic.** The plan runs before any effect, so the cap exists only to stop a slow plan from eating the handler's budget. The handler gets `min(HandlerCap, remaining − HandlerHeadroom)` (`effectful.go:111–118`). The largest plan cap that still leaves the handler its full cap is `invokeDeadline − HandlerCap − HandlerHeadroom = 20 − 10 − (2 + 4) = 4 s`. A cap tighter than that only produces false failures; a looser one steals handler time. The value is computed from named constants and pinned by a test; it is not tuned to a machine. **Sufficiency is measured, not assumed.** At the heaviest proxy load where the rig still boots, uncapped plans peak at 3.72 s (V15). That is inside 4 s but with thin margin (7%), which is why b2 ships too.
- **b2 (chosen, M3). Stop compiling from cold on every call.** The capsule keeps a host-owned, per-`(interpreter HashRef, source sha256)` compile-cache template. It is built by `ailang check` at the capsule's own staging path, which executes no transition code (V20). Each run copies the template into a fresh per-run cache directory passed as `AILANG_CACHE_DIR` (V21). Runs never share a writable cache, and a run cannot mutate the template. Measured effect: child 5–7× faster at every load (V18). On the 8-burner proxy that is a ~2.4 s cold child → ~0.5 s, so the plan becomes ≈ verify 0.4–0.66 s + 0.5 s, i.e. roughly 3× under the derived 4 s cap. It also speeds up the finish phase, pure transitions and replay, all through the same `RunContext`. Pre-warm happens at daemon start for every registered descriptor's `(interpreter, source)`, with a lazy build on a miss. Otherwise the first call of every fresh daemon (the CI tests' first call) would still be cold.
  - **Integrity argument.** The template is produced by the pinned interpreter (verified by hash) from the published source (content-addressed). It is keyed by both, so it can only contain compile products of that exact pair. It is copied, never linked, into each run. The interpreter validates cache entries by its own `cache_key` (V17 manifest), so a stale or empty template costs speed, not correctness. AC3.3 makes "warm output = cold output" a gate over every se-tools plan fixture, which bounds the replay-determinism risk.
- b3. Hoist or cache the 100 MB re-verification (18–30% of the phase, V15). **Deferred, not chosen.** Caching by file identity changes the integrity semantics the capsule doc settled (`w-capsule-output-cap-load-flake.md` §3 "Why the alternatives lose": a TOCTOU analysis is required). After b2 the plan fits with margin without it. A follow-up row is proposed in §10.
- b4. Raise `invokeDeadline` (bounded by `writeTimeout` 30 s and pinned by tests). Rejected: it moves a client-visible contract to hide a cost b2 removes.
- b5. Raise `PlanPhaseBudget` to a guessed number (5 s, 10 s). Rejected by the row's rule: a bound and a stimulus that both scale with machine speed must be derived. b1 derives it.

### (c) The test rig

- **c1 (chosen, M1). Make the e2e tests report the cause and stop hiding errors.**
  - Every `r.call` in `exec_e2e_srt_test.go` asserts `wire.Error == nil` and prints `*wire.Error` (not the pointer).
  - The exec srt rig routes the daemon `ErrorLog` to `t.Log` through a test writer, so a failing test shows the operator line naming the phase.
  - This merges signature 2 into signature 1 and makes any recurrence diagnosable from CI alone.
- **c2 (chosen, M4). A CI evidence gate, not a retry.** The M3 PR head must pass the parallel legs on **5 of 5** consecutive CI attempts, counting every attempt. At the observed base red rate of 11/13 (V23), five greens by chance has probability (2/13)^5 ≈ 9 × 10⁻⁵. Even at a 50% base rate it is ≈ 0.03. Each attempt is recorded in the sprint log.
- c3. Move the e2e tests out of the parallel legs (they already run PASS-gated in the `-p 1` step). **Rejected as the fix, kept as a named lever.** The parallel leg is the only CI witness of the product under load, and removing it hides row 153's defect rather than fixing it. It is admissible only if M4's gate fails after M3, and then only as an attended decision.
- c4. Inflate the test's own timeouts, `t.Skip` under load, or retry until green: forbidden by the row. None is used. The tests assert no wall-clock bound of their own: the bound under test is the product's.

**Recommendation:** a1 + b1 + b2 + c1 + c2, in four landings (§6). b3 becomes a follow-up row.

## 6. Milestones

Each milestone is one PR, CI-green on its own. Gate list per milestone in §8.

### M1 — typed phase timeout + diagnosable rig (~0.4 d)

Files: `host/coordinator/errors.go`, `host/coordinator/effectful.go` (`runPhase`), `host/projection/projection.go` (`dispatchError`), `host/daemon/exec_e2e_srt_test.go`, the exec srt rig helper, and tests.

- **AC1.1** `runPhase` returns `*PhaseTimeoutError{Phase:"plan", Budget, Elapsed}` iff `pctx.Err() != nil && ctx.Err() == nil`. `errors.Is(err, context.DeadlineExceeded)` is true. The classification is a pure helper `classifyPhaseErr(phase, budget, elapsed, parentErr, phaseErr)`, unit-tested over the 4 context-state combinations with no sleeps.
- **AC1.2** The caller's own deadline is never relabelled. With the parent context already expired, the error is not a `*PhaseTimeoutError` and `dispatchError` still answers `invocation exceeded its deadline`.
- **AC1.3** A2A: a fake coordinator returning `*PhaseTimeoutError{Phase:"plan"}` → `tasks/send` answers `-32603` with the message starting `PhaseTimeoutPrefix`, naming `plan` and the budget, and saying nothing ran.
- **AC1.4** A2A retry: a plan that overruns once (a `phaseRunner` whose first plan blocks until `ctx.Done()`, under a test-only budget hook), then the **same task id** resent → `completed`, with exactly one log entry and one effect record (V28).
- **AC1.5** MCP frozen wire (row 159 tripwire): the same plan timeout over `/mcp/` answers `{"code":-32603,"message":"host callback timed out"}`. The rig's `ErrorLog` contains `plan phase exceeded its`. The test comment names row 159 / ailang#1602 as the flip.
- **AC1.6** Every `call` in `TestExecSrtMCPEndToEnd` / `TestExecSrtMCPGrantAndBudget` fails with the printed `*wire.Error` when it is non-nil. The rig's `ErrorLog` is `t.Log`.

### M2 — the derived plan cap (~0.15 d)

Files: `host/coordinator/effectful.go`, `host/daemon/workspace_exec_test.go` (or a sibling).

- **AC2.1** `PlanPhaseBudget = 4 * time.Second`, with the derivation as a comment: invokeDeadline − HandlerCap − HandlerHeadroom.
- **AC2.2** A daemon-package test asserts `coordinator.PlanPhaseBudget == invokeDeadline − coordinator.HandlerCap − coordinator.HandlerHeadroom`, the same pattern as `TestExecMaxTimeoutIsHandlerCapLessOneSecond` (`workspace_exec_test.go:127`). Changing any of the four constants without the others reds it.
- **AC2.3** A coordinator test pins the arithmetic's purpose: with a plan that uses its whole budget (fake runner blocks until `ctx.Done()` under a scaled test budget), the handler's context deadline is still ≥ `HandlerCap − ε` after the plan. This kills "cap raised past the derivation".
- **AC2.4** `FinishPhaseBudget` stays 2 s, as an explicit decision. Raising it would shrink the handler cap through `HandlerHeadroom`, and M3 removes the finish phase's cold cost too. Residual R2 records this.

### M3 — compile-cache template (~0.7 d)

Files: `host/capsule/capsule.go` (+ a small `cachetemplate.go`), `host/daemon/daemon.go` (pre-warm at start), tests.

- **AC3.1** `Runner` gains `Prewarm(interp, source) error`. It stages the source at `host/capsule/main.ail` in a scratch root and runs `<interp> check host/capsule/main.ail` with `AILANG_CACHE_DIR=<tmp>` under a bounded context (the archive check's 10 s, `host/archive/check.go:20`). On rc 0 with a non-empty cache, it atomically renames `<tmp>` to `<state>/capsule-cache/<interp-digest>/<source-sha256>/` (mode 0700, host-owned). Failure is logged and non-fatal: runs fall back to cold.
- **AC3.2** `RunContext` copies the template (if present) into a per-run directory inside the run's temp root, sets `AILANG_CACHE_DIR` to it, and leaves the template byte-unchanged afterwards. A test hashes the tree before and after a run that would rewrite it.
- **AC3.3** Determinism gate: for every se-tools package fixture's plan input (and the finish inputs of the tools that have one), warm stdout is byte-equal to cold stdout under the pinned interpreter. This runs in CI on `AILANG_BIN` (no srt needed).
- **AC3.4** Effect gate: a capsule test runs the exec plan 5× cold (no template) and 5× warm and asserts **median warm < median cold / 2**. This is a ratio, so machine speed cancels out. V18 measured 4.4–7× at all loads; the test asserts 2× to leave room.
- **AC3.5** Canary: a test asserts that the pinned interpreter writes `manifest.json` under `AILANG_CACHE_DIR` after `check`. A floor-raise that drops or moves the cache reds here instead of silently going cold. This is listed for the S8 floor-raise coupling inventory.
- **AC3.6** The daemon pre-warms, at start, every effectful descriptor in the current registry head (and on publish). With the probe environment, the first plan of a fresh daemon shows a warm child (an AC-only check, recorded as a V-row in the sprint log).

### M4 — evidence, records, row 161 (~0.25 d)

- **AC4.1** Local proxy re-run on the M3 head: `$S/bg_load.sh 8 4` (the same script, archived into the sprint's verification dir) → **0 plan timeouts**, with plan `took` values recorded. The control is the V13 run on the base tree (4/4 timeouts).
- **AC4.2** CI gate c2: 5/5 consecutive attempts green on the M3 (or M4) head. Every attempt's run id and conclusion is recorded.
- **AC4.3** Charter: row 153 LANDED, and row 161 marked closed by 153 (V14/V24 cited). Proposed follow-up rows from §10 are filed by the controller.

## 7. Test plan (what each test kills)

| Test (new unless noted) | Pkg | Kills the mutation |
|---|---|---|
| `TestClassifyPhaseErrMatrix` (AC1.1/1.2) | coordinator | MUT-RELABEL: every expiry labelled a phase timeout (a parent deadline mislabelled). MUT-NOWRAP: `Unwrap` dropped, so MCP flips to `host callback failed` |
| `TestA2APlanTimeoutIsTypedAndNamesThePhase` (AC1.3) | projection | MUT-ORDER: the new case placed after the generic deadline case; MUT-GENERIC: message lacks the phase |
| `TestA2APlanTimeoutRetrySameTaskID` (AC1.4) | coordinator / daemon | MUT-EARLY-INTENT: an intent appended before the plan, so the resend answers `NotCommitted` |
| `TestMCPPlanTimeoutFrozenWire` (AC1.5) | daemon | MUT-NOWRAP (wire side); row-159 tripwire |
| `TestExecSrtMCP*` amended (AC1.6) | daemon | MUT-SILENT-ERROR: a second-call error read as a nil field (signature 2) |
| `TestPlanPhaseBudgetIsDerived` (AC2.2) | daemon | MUT-GUESS: `PlanPhaseBudget` set to a literal ≠ the arithmetic; any one of the four constants changed alone |
| `TestSlowPlanLeavesHandlerItsCap` (AC2.3) | coordinator | MUT-STEAL: plan cap > invokeDeadline − HandlerCap − HandlerHeadroom |
| `TestCapsuleTemplateIsCopiedNotShared` (AC3.2) | capsule | MUT-SHARE: `AILANG_CACHE_DIR` pointed at the template itself |
| `TestCapsuleWarmEqualsColdForSeTools` (AC3.3) | capsule | MUT-STALE-KEY: template keyed by interpreter only, so source B is served source A's cache. The fixture includes two sources whose outputs differ. |
| `TestCapsuleWarmIsFaster` (AC3.4) | capsule | MUT-NOTEMPLATE: the template never copied in (silent cold path) |
| `TestPinnedInterpreterHonoursCacheDir` (AC3.5) | capsule | floor-raise drift (the env var dropped or renamed) |
| `TestDaemonPrewarmsRegisteredEffectful` (AC3.6) | daemon | MUT-LAZY-ONLY: no pre-warm, so the first call stays cold |

## 8. CI gate list (every milestone)

```sh
export PATH=/opt/homebrew/bin:$PATH AILANG_BIN=~/.pinned-ailang/ailang   # gate interpreter v0.41.0
./scripts/verify_ail.sh                        # no .ail change in this row; still run
go vet ./...                                   # compile fence for _test.go (go build is not one)
go build ./... && go test ./... -count=1
go test -race ./... -count=1 -timeout 8m       # CI's race leg
# with real srt (rig: WORLD_EXEC_SRT_NODE_MODULES=<pinned 0.0.78 install>):
go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -v -run '^(TestExecSrtMCPEndToEnd|TestExecSrtMCPGrantAndBudget|TestExecSrtReplayIsByteEqual)$'   # each must print --- PASS
```

Every new test name is checked with a `--- PASS: <name>` grep in the same call. A bare `-run` with no match exits 0. The tool binary (v0.52.1) is not needed: no se-tools e2e test changes behaviour here.

## 9. Conflict Surface

| File | Change | Interaction |
|---|---|---|
| `host/coordinator/effectful.go` | `runPhase` classification; `PlanPhaseBudget` 2 → 4 s | **Row 152** (exec availability decoupled from the policy summary) changes the binder build in `Dispatch`/`daemon`, not `runPhase`: no overlap. If 152 lands first, rebase only. |
| `host/coordinator/errors.go` | `PhaseTimeoutError` | none |
| `host/projection/projection.go` | `dispatchError` case + prefix const | **Row 159** later sends the same `dispatchError` message over MCP. M1's frozen-wire test is the tripwire 159 flips. |
| `host/capsule/capsule.go` (+ new file) | template copy-in, `Prewarm`, `AILANG_CACHE_DIR` | Used by the pure path, replay and the archive check's peers. **Row 154** (`-race` leg headroom, ≈ 540–556 s of 600 s) benefits, because every capsule run in the suite gets faster. No conflict. |
| `host/daemon/daemon.go` | pre-warm at start / on publish | Startup time grows by one `check` per effectful descriptor (≈ 30–180 ms each idle, V20). It runs after the listener-independent init; M3 measures it. |
| `host/daemon/exec_e2e_srt_test.go`, rig helper | error assertions, `t.Log` ErrorLog | **Row 161**: closed by this row |
| — | — | **Row 93** (World arm under load): M2 + M3 remove the plan-phase timeout from its measurement. The plan cost per call drops ~5×, and that cost counts against World's +25% wall-clock bound (R-SE-1). |

## 10. Non-goals, residuals, follow-up rows (proposed; the controller files)

- **R1 — re-verification cost (b3).** 18–30% of every capsule run hashes 100 MB. Proposed row: *w-capsule-verify-once* (clause-4). Measure the TOCTOU surface first, then decide between hoisting verification into the archive's ingest and caching by `(dev, inode, size, mtime, ctime)` plus a periodic re-hash.
- **R2 — finish phase.** It stays at 2 s (AC2.4). After M3 its child is warm too. Raising it would need re-deriving `HandlerHeadroom`. Revisit only if a measured finish overrun appears. An overrun there already surfaces as the typed `EffectsUnrecordedError`.
- **R3 — rig-only load failures beyond the plan (V16).** At 12 throttled burners the startup probe's 10 s subprocess cap and publication under the 30 s test context fail first. CI has never shown them (V16 control). Not filed. Recorded so a future flake with those signatures is recognised.
- **R4 — MCP clients still see `host callback timed out`** until row 159 / ailang#1602. That is honest (it is a timeout) but does not name the phase on the wire.

## 11. Risks

- **The template trusts interpreter-written files across runs.** Mitigated by the key `(interp digest, source sha256)`, copy-in, AC3.3's determinism gate, and the interpreter's own `cache_key` validation. Residual: an interpreter cache bug that returns stale code for an equal key. This is upstream's invariant, and AC3.3 would catch it for every shipped source.
- **4 s may still be exceeded on a worse runner.** M3 is what buys the margin (≈ 3× on the proxy). M4's 5/5 gate measures CI directly. Lever c3 exists as an attended-only fallback.
- **Undocumented cache layout.** AC3.5's canary turns a silent cold regression into a red test at the next floor-raise.

## 12. Upstream asks

None required. `AILANG_CACHE_DIR` (V21) and `check`'s cache write (V20) suffice on v0.41.0. *Optional hardening, not filed:* a documented read-only compile-cache mode (`--cache-readonly`), so the copy-in step could become a bind. Ailang#1602 (typed MCP errors) already exists and is row 159's gate.

## 13. UNMEASURED

- **U1 — Linux.** All timings are darwin/arm64. The cache mechanism is the interpreter's own, and CI's first-call numbers fit the throttled proxy (V12 vs V23), but the cold/warm ratio on linux/amd64 was not measured. M3's AC3.4 measures it in CI because it runs there.
- **U2 — CI core count.** Taken from the repo's note (V30), not re-read from a runner.
- **U3 — The CI plan-time tail.** CI logs carry no phase timings (V25). The 4 s sufficiency claim rests on the proxy (V15). M4's AC4.2 is the direct CI measurement.
- **U4 — Pre-warm cost at daemon start** for the full registry (9 se-tools sources plus pure transitions). Estimated from V20 at ≈ 30–180 ms each idle. M3 records it.

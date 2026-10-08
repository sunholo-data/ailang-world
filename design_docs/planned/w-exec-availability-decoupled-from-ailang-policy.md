# w-exec-availability-decoupled-from-ailang-policy — a failed AILANG tool-handler build unbinds `Workspace.Exec`, which never uses it (row 152)

**Status:** PLANNED. Designed in iteration 242 (designer lane claude-opus-5-5, unattended).
**Queue:** World mission row 152 `w-exec-availability-decoupled-from-ailang-policy` (clause 4). It is next after row 153 (D-WORLD-65, re-ordered by D-WORLD-68 = A).
**Estimated:** ~0.75 d of host Go and tests (M1 ~0.35 d, M2 ~0.35 d, M3 ~0.05 d). **Dependencies:** none. No `.ail` changes, no pin changes, no upstream gate.
**Measured base:** worktree `sprint/row152-exec-availability` at `fe7c2a4`. Gate interpreter `~/.pinned-ailang/ailang` = `AILANG v0.41.0`. `go1.26.6 darwin/arm64`. Probe test banked as `~/.ailang/state/mission-world-iter242/zz_probe152_test.go.txt`, output as `…/probe152.log` (`$S` below). **The probe was removed from the tree; this commit touches only this document.**

## 1. Problem

`workspaceTools.registry(episodeID)` (`host/daemon/workspace.go:307–336`) builds an episode's handler registry: eight AILANG tool effects bound to one `AilangToolHandler`, plus `Workspace.Exec`. It has four early returns, each returning an **empty** registry (V1). Exec is bound only on the success path, at `:334`. Two of those returns are AILANG-only failures:
- the module root (`sandboxRoot`);
- the AILANG handler build (`episodeHandler`): row 141's package-cache refusals, the policy write, and the `policy-tool summary` subprocess.

`execHandler` takes only `epRoot` and the exec config. It never reads the module root, the policy or the AILANG cache (V2). The probe shows the effect: when the summary subprocess cannot finish, `registry("ep1")` returns `[]`, so `Workspace.Exec` is unbound. Yet the exec handler built directly for the same episode runs `exit 0` in the ep1 worktree (V6, V7). The coordinator then refuses `workspace-exec` with R8 before its plan (V12).

The coupling has two more costs, both measured:
- **The failed build is retried on every Binder open, on the caller's clock.** It is never cached (V3). The Binder is opened for every Dispatch, pure calls included (V11). So each call in a broken episode waits the full construction budget, 3.0 s, before it starts (V6: 3.007 s then 3.018 s, 2 summary spawns for 2 calls).
- **The build holds the shared mutex `w.mu` across the subprocess.** `execHandler` takes the same mutex. Building another episode's exec handler waited **2.82 s** behind ep1's stuck summary (V8).

## 2. Measured facts (summary)

| # | Claim | V |
|---|---|---|
| F1 | `registry()` has 4 early returns, all empty. Exec is bound only after `sandboxRoot` and `episodeHandler` succeed. | V1 |
| F2 | `execHandler` depends on `episodeID`, `epRoot`, `w.exec`, `w.root` and `w.stateDir`. It does not depend on the sandbox (module root), the policy or the cache dir. | V2 |
| F3 | A summary that never answers → `registry` = `[]` after **3.0 s**, with no `Workspace.Exec`. The exec handler built directly runs `exit 0`. A non-ok summary takes the same path in 17 ms. | V6, V7 |
| F4 | A failed build is not cached. Every `registry()` call re-spawns the summary (1 → 2 spawns over 2 calls). | V3, V6 |
| F5 | The build holds `w.mu` across the subprocess, so another episode's exec construction waits behind it (2.82 s). | V3, V8 |
| F6 | **The "10 s" in the row text is a mislabel.** The effective bound is `workspaceHandlerBudget = 3 s`, both now and at M4's build `7664d3f`. `HandlerTimeoutError` prints the handler's own `execTimeout` (10 s), even when the parent context's earlier deadline is the one that fired. | V4, V5, V6, V10 |
| F7 | **The M4 `serve.log` does not show the sequence.** It has 3 lines: 0 mention `policy summary`, `registered handler` or `go1`, and 2 are the control `plan phase` lines. The only record is the prose of README finding F5. | V9 |
| F8 | The Binder (and so `registry()`) is opened for every Dispatch: inside the episode lock for effectful calls, and for pure calls too. | V11 |
| F9 | R8 is decided per transition. `workspace-exec` declares only `Workspace.Exec`, and each `ailang-*` tool declares only AILANG names. So leaving the AILANG names unbound refuses only the AILANG tools. | V12, V14 |
| F10 | A typed failure handler for the eight AILANG names would be worse for the client than R8: (1) it commits a `failed` invocation; (2) it debits the budget; (3) the record has no field for the cause; (4) the wire body is `{}`. R8 refuses before the plan with no record and no debit. | V12, V13 |
| F11 | The row-153 plan cap is derived as `invokeDeadline − HandlerCap − HandlerHeadroom = 4 s`, assuming nothing before the plan uses the invoke deadline. A 3 s registry stall before the plan breaks that assumption. | V15 |

## 3. Verification Log

Shorthand: `W=` this worktree, `S=~/.ailang/state/mission-world-iter242`. The probe ran as `PATH=/opt/homebrew/bin:$PATH AILANG_BIN=~/.pinned-ailang/ailang go test ./host/daemon/ -count=1 -v -run '^TestZZProbe152'` with `$S/zz_probe152_test.go.txt` copied to `host/daemon/zz_probe152_test.go`. Output is in `$S/probe152.log`, and both probe tests printed `--- PASS`. The probe's tool binary is a `/bin/sh` stub shaped like `fakeToolBin` (`workspace_test.go:95`): it answers `--version` with `AILANG v0.52.1`, and its `policy-tool` summary either `sleep 30`s (mode `sleep`) or prints `{"ok":false}` and exits 1 (mode `fail`). Exec is configured with the existing `stubExecSandbox` + `execProfile` (`workspace_exec_test.go:26,50`, command `where` = `pwd -P`). **Rule applied:** each null result is paired with a positive control.

| V | Command (trimmed) | Reading | Control / scope |
|---|---|---|---|
| V1 | `sed -n 304,336p host/daemon/workspace.go` | Returns `broker.Registry{}` at `:309` (`!w.enabled()`), `:313` (`episodeRoot` not ok), `:318` (`sandboxRoot` err) and `:328` (`episodeHandler` err; an `*episodeRefusal` prints its own line `:324`, anything else prints `workspace tools unavailable for episode %q: %v` `:326`). `reg[broker.EffectWorkspaceExec] = w.execHandler(episodeID, epRoot)` is at `:334`, after all four. | `grep -c 'return broker.Registry{}' host/daemon/workspace.go` → 4 |
| V2 | `sed -n 342,369p host/daemon/workspace.go` | `execHandler(episodeID, epRoot)` reads `w.exec`, `w.execH`, `w.root`, `w.stateDir`, `w.exec.{sandbox,operatorHome,maxOutputBytes}`, `profileFor(episodeID)`. `NewExecHandler` gets `Worktree: epRoot`. It takes `w.mu` (`:346`). A construction error → `execFailedHandler{err}` (`:361`). | Its arguments contain no `sandbox`/`policyPath`/`cacheDir` (they are locals of `episodeHandler`, `:401–414`) |
| V3 | `sed -n 395,437p host/daemon/workspace.go` | `episodeHandler` takes `w.mu` for its whole body (`:396–397`). It runs, in order: `MkdirAll` policy/cache dirs, `linkPackageCache`, `checkLockCoverage`, `RenderEpisodePolicy`, `writePolicy`, then `NewAilangToolHandler` under `context.WithTimeout(context.Background(), workspaceHandlerBudget)` (`:422`). It writes to the cache `w.handler[...]` only after success (`:432–435`). Every error returns before that. | The success cache is proved by the existing `TestWorkspaceRegistryRefusesEpisodesOutsideTheGrammarAndRoot` (`summaries == 1` after two calls, `workspace_test.go:218–222`) |
| V4 | `grep -rn 'workspaceHandlerBudget\|ailangToolExecTimeout' host --include='*.go'` | `workspaceHandlerBudget = 3 * time.Second` (`workspace.go:44`), used once (`:422`). `ailangToolExecTimeout = 10 * time.Second` (`handlers_ailang.go:48`), the handler's default `execTimeout` (`:328`). `NewAilangToolHandler` → `loadSummary` → `policyToolWith(ctx, …)` (`handlers_ailang.go:331,338`; `handlers_ailang_run.go:376`). | static |
| V5 | `sed -n 50,52p;153,157p;204,232p host/broker/handlers.go` | `runBounded`: `runCtx, cancel := context.WithTimeout(ctx, bounds.execTimeout)` (`:156`). On expiry it returns `&HandlerTimeoutError{Timeout: bounds.execTimeout}` (`:208,229,231`). `Error()` = `fmt.Sprintf("%v after %s", ErrHandlerTimeout, e.Timeout)` (`:51`). A parent deadline that fires first is reported with the child's 10 s. | V6 is the runtime positive |
| V6 | probe `TestZZProbe152SummaryTimeout/sleep` | call 1: `elapsed=3.007s names=[] execBound=false summaries=1`. call 2: `elapsed=3.018s names=[] execBound=false summaries=2`. Operator log, twice: `ailang-worldd: workspace tools unavailable for episode "ep1": broker: policy summary: broker: handler subprocess timed out after 10s (stderr "")`. | **Control:** `execHandler("ep1", epRoot)` called directly in the same daemon → `{"exit_code":0,…,"stdout":".../ws/ep1\n",…}`, `err=<nil>` |
| V7 | probe `…/fail` | `elapsed=17ms`/`18ms`, `names=[]`, `execBound=false`, `summaries=2`. Log: `… broker: policy summary: broker: handler subprocess failed: exit status 1 (output "{\"ok\":false}\n")`. | Same control → `exit_code 0`. Deterministic trigger of the same path with no wall-clock bound |
| V8 | probe `TestZZProbe152LockCoupling`: `registry("ep1")` in a goroutine (stub summary sleeps), 200 ms later `execHandler("ep2", ep2Root)` | `execHandler(ep2) waited 2.822s behind ep1's summary` | Cross-episode coupling through `w.mu`. The 2.82 s matches 3 s − the 0.2 s head start. |
| V9 | `cd design_docs/verification/world-row140-m4; grep -c 'policy summary' serve.log; grep -c 'registered handler' serve.log; grep -c go1 serve.log; grep -c 'plan phase' serve.log`; `grep -c 'policy summary' README.md`; `sed -n 1,11p log-tail.txt` | serve.log: **0 / 0 / 0**, control `plan phase` **2** (it has 3 lines in all). README: 1 (finding F5's prose). `log-tail.txt` shows two `go1 workspace-exec [Workspace.Exec ok]` entries (#1, #2). | The row's citation of `serve.log` is **false**; the sequence is reported only in README F5. |
| V10 | `git show 7664d3f:host/daemon/workspace.go \| grep -n workspaceHandlerBudget`; `git show 7664d3f:host/broker/handlers.go \| grep -n 'HandlerTimeoutError{Timeout'` | At M4's build: `workspaceHandlerBudget = 3 * time.Second` (`:43`), and `HandlerTimeoutError{Timeout: bounds.execTimeout}` at `:208,229,231,313,328,330` | F5's "timed out at 10 s" was the V5 mislabel. The bound was 3 s. |
| V11 | `sed -n 267,320p host/coordinator/coordinator.go` | `Dispatch` evaluates `c.cfg.Binder(call.EpisodeID, grants)` as an argument of `transitionreg.Bind` (`:316–317`) on **every** call. For effectful calls this happens after `lockEpisode` (`:283–289`). Production `Binder` = `d.binder` (`daemon.go:681`) → `broker.OpenBinder(…, d.workspace.registry(episodeID))` (`daemon.go:745`). | `grep -n 'workspace.registry' host/daemon/*.go \| grep -v _test` → only `:745` (plus the comment at `:449`) |
| V12 | `sed -n 274,279p host/coordinator/effectful.go`; `sed -n 37,44p host/coordinator/errors.go`; `grep -n 'no handler for' host/projection/projection.go`; `sed -n 21,32p $GOMODCACHE/github.com/sunholo-data/ailang@v0.47.2/serveapi/protocol/envelope.go` | R8 runs before `loadSourceAndWorld` and the plan: `if missing := bound.Unhandled(); len(missing) != 0 { return …&EffectsUnsupportedError{…} }`. Operator text: `coordinator: transition %q declares effects with no registered handler %v; effectful invocation is unavailable`. **A2A wire:** `-32603 "transition declares an effect this daemon has no handler for"` (`projection.go:462`). **MCP wire:** `CallbackMessage` sends anything that is neither capacity, `DeadlineExceeded` nor `Canceled` as `"host callback failed"`. | The existing table test `projection_test.go:828` pins the A2A string. `setools_e2e_test.go:825` pins `host callback failed` on the MCP wire (for R14). |
| V13 | `sed -n 291,308p host/broker/broker.go`; `sed -n 19,30p host/broker/record.go`; `sed -n 177,203p;243,245p host/coordinator/effectful.go`; `grep -n 'has no finish phase' packages/se-tools/se_tools/exec.ail` | A handler error writes `EffectRecord{…Allowed: true, Failed: true…}` with `BudgetAfter: decision.Remaining` (a debit), appends outcome `failed`, and returns `EffectFailedError`. `EffectRecord`'s fields are Effect, Scope, Cost, BudgetBefore/After, Allowed, Failed, Denial, RequestRef, ResultRef: **no cause**. The coordinator maps this to `StatusFailed` and commits. With one effect and no finish, the body is `{}` (`:244–245`, "denied or failed: the world block carries the status and record"). exec.ail has no finish phase (`:356`). | `Denial` is the positive control: a record field for a reason exists only for denials |
| V14 | `python3 -I -c` over `packages/se-tools/transitions.json` printing `id`, `declaredEffects` | 9 entries. `workspace-exec` → `[Workspace.Exec]`. `ailang-read` → `[Workspace.Read]`, `ailang-write`/`ailang-edit` → `[Workspace.Write]`, `ailang-check` → `[Ailang.Check]`, `ailang-run` → `[Ailang.Run, Ailang.RunEnv, Ailang.RunNet]`, `builtins-search`/`examples-search` → `[Ailang.Discover]`, `ailang-cli` → `[Ailang.CLI]` | No transition declares names from both families |
| V15 | `sed -n 36,49p;118,125p host/coordinator/effectful.go`; `grep -n 'invokeDeadline *=' host/daemon/daemon.go` | `PlanPhaseBudget = 4s` is derived in the comment as "20 s − HandlerCap − HandlerHeadroom = 20 − 10 − (2 + 4) = 4 s". `HandlerCap = 10s`, `HandlerHeadroom = Finish (2s) + PostEffect (4s)`. `handlerBudget` = `min(HandlerCap, time.Until(deadline) − HandlerHeadroom)`. `invokeDeadline = 20s` (`daemon.go:99`). | A 3 s stall before the plan leaves the handler `min(10, 20 − 3 − 4 − 6) = 7 s` when the plan uses its full budget. That is arithmetic on the constants, not a measured run (UNMEASURED U2). |
| V16 | `grep -n 'len(reg) != 0\|want empty' host/daemon/*_test.go`; `sed -n 889,915p host/daemon/workspace_layout_test.go` | Existing "empty registry" assertions: `workspace_layout_test.go:159` (module root), `:586` (package cache plants), `:874` (lock coverage), `:905` (FIFO lock, `n != 0`); `workspace_test.go:203` (episode grammar), `:277`, `:290` (one flag alone). | M1 changes the first four and must leave the last three unchanged (§6) |
| V17 | `sed -n 95,121p host/daemon/workspace_test.go`; `sed -n 449,538p host/daemon/workspace_test.go`; `sed -n 128,130p host/daemon/workspace.go` | Seams that already exist: `fakeToolBin` (a shell stub archived as the tool binary, which counts `summaries`/`dispatches`), `publishReadTool` + `readPlanRunner` (a coordinator wired to `d.binder` with a fake plan runner, so no interpreter is needed), and `stubExecSandbox`. `execStartupHook` adjusts only `ExecStartupConfig` at startup (`configureExec`), **not** handler construction, so it cannot drive the summary. | V6/V7 used these seams directly |
| V18 | `sed -n 988,991p host/daemon/daemon.go`; `grep -n 'func (d \*Daemon) stopTemplateMaintenance' host/daemon/templates.go` | `Close()` = `d.stopTemplateMaintenance(); return d.store.Close()`. Row 153's background rebuild (`templates.go:158`) is the precedent for a daemon-owned goroutine that `Close` joins. | static |
| V19 | `sed -n 984,989p docs/QUICKSTART.md` | "Each is written once per refused call, and the episode's tools stay refused (**every declared effect fails, as for a missing worktree**) until you fix the cause". This sentence becomes false after M1 (S7). | static |
| V20 | `sed -n 229,245p;283,288p .github/workflows/ci.yml`; `sed -n 290,330p scripts/verify_go.sh` | CI runs `verify_ail.sh`, then `verify_go.sh` (`go build ./...`, `go test ./... -count=1`, `go test ./... -count=1 -race -timeout 8m` under a 600 s kill), a linux srt step (`go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -v -run '^(…TestExecSrtMCPEndToEnd|TestExecSrtMCPGrantAndBudget|TestExecSrtReplayIsByteEqual)$'` with a `--- PASS` grep per name) and the subprocess-cleanup step | §8 is derived from these |

## 4. Root cause

`registry()` treats the episode as one unit: either all nine names are bound, or none are. That was correct at row 134, when the eight AILANG names were the whole registry. Rows 140 (exec) and 141 (module root, package cache) added a second family and AILANG-only checks, but the all-or-nothing return was kept. The AILANG build is also synchronous, uncached on failure and serialised behind `w.mu`. So one slow subprocess delays every call in the episode and every other episode's exec construction.

**What the row text got wrong.** The "10 s cap" is the handler's default `execTimeout`. `HandlerTimeoutError` names it even when the parent's 3 s `workspaceHandlerBudget` fired first (F6). This design does not fix that label. It lives in the shared `runBounded` timeout surface, which row 166 (broker head/tail timeout flake) also touches, so it is proposed as a follow-up row (§10 R1).

## 5. Options

### (a) Which failures leave exec bound

The deciding fact is which inputs each family's handler needs (F2):

| Failure (V1) | AILANG names | `Workspace.Exec` | Why |
|---|---|---|---|
| `!enabled()` (a flag missing) | unbound | **unbound** | Exec needs the workspace tools. `configureExec` already refuses `--exec-*` without them (`workspace.go:142–145`), so this case is only reachable with exec unconfigured. Unchanged. |
| `episodeRoot` refused (grammar, symlink, absent, not a dir) | unbound | **unbound** | `epRoot` *is* exec's worktree. Without a safe one exec has no sandbox. Unchanged (AC4.4). |
| `sandboxRoot` err (row 141 module root) | unbound | **bound** | The module root is the AILANG sandbox, and exec never reads it (F2). QUICKSTART already says "an exec profile's `root` is independent of the module root" (`QUICKSTART.md:981–982`). |
| `episodeHandler` err: `episodeRefusal` (package cache / lock coverage), mkdir, policy render/write, `NewAilangToolHandler` (summary timeout or non-ok summary) | unbound | **bound** | All are AILANG-only inputs (V3). |

- **a1 (chosen).** Bind exec iff `enabled && episodeRoot ok`. Bind the eight AILANG names iff the AILANG handler is built. On an AILANG failure the operator line is unchanged.

### (b) What the eight AILANG names bind to when their handler cannot be built

- **b1 (chosen). Leave them unregistered; R8 refuses the AILANG transitions only.** This is today's documented behaviour for the AILANG family, now scoped to that family (F9).
  - Measured client view: A2A `-32603 "transition declares an effect this daemon has no handler for"`, MCP `-32603 "host callback failed"` (V12).
  - Nothing runs, nothing is recorded, no budget is debited, and a resend is safe.
  - The operator line names the cause.
- b2. A typed failure handler, like `execFailedHandler`. **Rejected on measurement (F10).** Each call would commit a `failed` invocation and debit a budget unit for a host-side construction fault. The cause still would not reach the client: the record has no cause field and the wire body is `{}` (V13). It is strictly more state for no added information.
  - `execFailedHandler` is right for exec because exec's construction failure is per-profile configuration with no R8-level alternative (row 140's choice). It is not a template for a transient subprocess failure.
- b3. A refusal-body handler (`{"ok":false,"refused":…}`, like `ExecRefusalHandler`). Rejected: the eight AILANG effects have eight output schemas, and none declares a `refused` arm on construction. This would mean `.ail` and output-contract changes for a host fault. If client-visible cause text is wanted, that is row 159 / ailang#1602's job for every R8 refusal, not this row's.

### (c) Retry policy for a failed AILANG build

- c1. Today: a synchronous retry on every Binder open. **Rejected.** It costs up to `workspaceHandlerBudget` (3.0 s, V6) per call, pure calls included (F8). The cost is paid under the episode lock and under `w.mu`, so it also blocks other episodes (V8). It also eats into the derived plan/handler budget (F11). And the failure happens when the machine is loaded, so retrying on every call adds load at the worst moment.
- c2. Cache the failure until restart. Rejected: a transient timeout would disable the episode's AILANG tools for the daemon's lifetime.
- c3. Cache the failure for a cool-down of N seconds. Rejected: N would be a guessed number (row 153's rule: derive bounds or do not add them).
- **c4 (chosen). Per-episode single-flight build state; after a failure, retry off the caller's clock.**
  - The first build for an episode stays synchronous. That is today's healthy path (~0.1 s, per the `workspaceHandlerBudget` comment), so the first AILANG call still works.
  - The build no longer holds `w.mu` across the subprocess. Concurrent first callers of the *same* episode wait on that one build. Other episodes and exec construction do not wait at all (fixes F5).
  - After a failure, `registry()` returns at once with exec bound and the AILANG names unbound. If no retry is in flight, it starts **one** background build, bounded by `workspaceHandlerBudget` and joined by `Daemon.Close` (the V18 precedent). Success is cached, and the next Binder open binds all nine names.
  - Cost under load: at most one summary subprocess per failing episode in flight. Zero added caller latency after the first failure.
  - Trade-off, stated: after an operator fixes a row-141 layout refusal, the *first* call afterwards still gets R8 while the background retry runs. The call after it succeeds.

### (d) How the test drives the summary to time out

No wall-clock ratio and no reliance on the 3 s constant. Two pieces:
1. **A seam.** `workspaceHandlerBudget` becomes a package `var`, following the `geteuid`/`killGroup` pattern. A test lowers it.
2. **A stub that never answers.** The stub's summary blocks on a FIFO (or `sleep 30`). The outcome is then a timeout on any machine, however fast or slow: the stub never answers, so only the budget can end the build.

Assertions are on outcomes: names bound, R8 vs committed, summary spawn counts, and error kind (`errors.As(err, *broker.HandlerTimeoutError)` reaching the operator line). They are never on elapsed time. Where a test needs "while the build is in flight", the FIFO holds it open until the test releases it, so the state is reached by construction, not by a sleep.

**Recommendation:** a1 + b1 + c4 + d, in two code landings plus records (§6). The `HandlerTimeoutError` label becomes a follow-up row (R1).

**Why is this not a package? (S3)** The handler registry is host-boundary wiring (S2). The broker binds effect names to Go handlers that spawn subprocesses, and no `.ail` code can express or replace that. This row adds no surface and no kernel growth. It narrows an existing host coupling and moves no policy out of AILANG: the confinement stays in AILANG's own policy layer.

## 6. Milestones

Each milestone is one PR, CI-green on its own; gates in §8.

### M1 — two families, two bindings (~0.35 d)

Files: `host/daemon/workspace.go`, `host/daemon/workspace_test.go`, `host/daemon/workspace_layout_test.go`, `host/daemon/workspace_exec_test.go` (or a new `workspace_family_test.go`), `docs/QUICKSTART.md`.

- **AC1.1** `registry()` binds `Workspace.Exec` whenever `enabled && episodeRoot ok`. It binds the eight AILANG names only when the AILANG handler is built. On a `sandboxRoot` or `episodeHandler` failure it prints the same operator line as today and returns exactly `[Workspace.Exec]`.
- **AC1.2** The exec handler cache gets its own mutex (`execMu`). `execHandler` never takes `w.mu`.
- **AC1.3** `workspaceHandlerBudget` becomes a `var` (a test seam, with a comment). The production value stays 3 s.
- **AC1.4 — the row's pin, end to end.** The fixture is all existing seams (V17):
  - a coordinator over `d.binder`;
  - a published `ws.exec` descriptor (Access and declared effect `Workspace.Exec@worktree`), using `publishReadTool` generalised to take the effect and id;
  - a fake plan runner that answers one `Workspace.Exec` `{"command":"where"}` effect;
  - exec configured with `stubExecSandbox` + `execProfile`;
  - a tool stub whose summary never answers;
  - the budget seam set low.

  Then:
  - (i) Dispatch `ws.exec` in ep1 **commits**: `exit_code 0`, and `stdout` is the ep1 worktree path.
  - (ii) Dispatch `ws.read` in ep1 → `*coordinator.EffectsUnsupportedError` naming `[Workspace.Read]`, with zero read plans run.
  - (iii) The operator log has `workspace tools unavailable for episode "ep1": broker: policy summary:`.

  Control: the same fixture with a healthy `fakeToolBin` → both calls commit.
- **AC1.5** Row 141 refusals leave exec intact. Four existing tests are amended from "want empty" to "want exactly `[Workspace.Exec]`": `workspace_layout_test.go:159, :586, :874`, and the FIFO test's `n != 0` → `n != 1` at `:905`. Their `summaries == 0` and exact-log assertions stay unchanged.
- **AC1.6** Episode-level refusals stay empty, unedited: `workspace_test.go:203` (all 18 `refusedEpisodes`), `:277`, `:290`, and `TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect`.
- **AC1.7** Exec construction never waits on an AILANG build. The test holds ep1's summary on a FIFO, with the budget seam set to the test's own bounded context, so only the FIFO can end the build. It then calls `execHandler("ep2", …)` and `registry("ep2")`'s exec entry. Both must return **while ep1's `registry` goroutine is still blocked**: the goroutine's done-channel is asserted not closed after they return. Then the test releases the FIFO.
- **AC1.8 (S7)** Update the QUICKSTART operator-lines paragraph (V19): the episode's AILANG tools stay refused, and `workspace-exec` keeps working (it needs only the worktree). Update the file header and the `ailangToolEffects`/`workspaceEffects` comments in `workspace.go`, which say "an EMPTY registry" and "ALWAYS bound".

### M2 — single-flight build, retry off the caller's clock (~0.35 d)

Files: `host/daemon/workspace.go`, `host/daemon/daemon.go` (`Close`), tests.

- **AC2.1** Per-episode build state: `{built tool | last error | in-flight done-chan}`, keyed as today by `(episodeID, sandbox)`. `w.mu` guards only the map, never a subprocess.
- **AC2.2** Same-episode single flight: two concurrent first `registry("ep1")` calls spawn **one** summary (stub spawn count == 1).
- **AC2.3** After a failure:
  - `registry("ep1")` returns `[Workspace.Exec]` without spawning synchronously;
  - it starts at most one background retry;
  - three further calls while that retry is held on the FIFO return at once and leave the spawn count at exactly 2 (1 failed + 1 retry).
- **AC2.4** Recovery: the stub is switched to answer ok (a flag file) and the FIFO released. A test-only `w.awaitBuilds()` joins the retry. The next `registry("ep1")` binds all nine names, and a recovery operator line is printed once: `workspace tools available for episode "ep1" after retry`.
- **AC2.5** `Daemon.Close` cancels and joins in-flight retries: the retry's context derives from a daemon-owned context that `Close` cancels. The test closes the daemon with a retry blocked on the FIFO and asserts that `Close` returns and the stub's process group is gone (the FIFO writer side observes EOF/EPIPE). No goroutine outlives `Close`.
- **AC2.6** The failure operator line is printed once per *attempt*, not once per call. With 1 failure + 3 calls during the retry, there are 2 lines when the retry also fails.

### M3 — records (~0.05 d)

- **AC3.1** Charter row 152 → LANDED. The controller files the §10 follow-up rows. Note in the M4 evidence README that F5's "10 s" was the 3 s budget mislabelled (V10), with a pointer to this doc (the evidence file itself is not rewritten).

## 7. Test plan (what each test kills)

Every row is mutation-proven by the executor: apply the named production mutation, see the named assertion go red, restore.

| Test (new unless noted) | AC | Kills the mutation |
|---|---|---|
| `TestWorkspaceExecSurvivesPolicySummaryTimeout` | 1.4 | **MUT-ALL-OR-NOTHING:** restore `return broker.Registry{}` on `episodeHandler` failure → (i) red with R8 `[Workspace.Exec]`. **MUT-BIND-AILANG-ON-FAIL:** bind the eight names to `execFailedHandler`-style failures → (ii) red (commits `failed` instead of R8). |
| amended `TestWorkspaceModuleRootSymlinkOrMissingIsR8`, package-cache plant / lock-coverage / FIFO tests | 1.5 | **MUT-SANDBOX-GATES-EXEC:** exec bound only after `sandboxRoot` succeeds → registry empty, red. **MUT-EXEC-USES-SANDBOX:** `execHandler(ep, sandbox)` instead of `epRoot` → the module-root test's exec `where` stdout ≠ ep1 worktree (the amended test runs one `where` call). |
| unedited `TestWorkspaceRegistryRefusesEpisodesOutsideTheGrammarAndRoot`, `TestWorkspaceRegistryNeedsBothFlags`, `TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect` | 1.6 | **MUT-EXEC-BEFORE-EPROOT:** exec bound before the `episodeRoot` check → `registry("escape")` non-empty, red |
| `TestWorkspaceExecNeverWaitsOnAilangBuild` | 1.7 | **MUT-SHARED-LOCK:** `execHandler` takes `w.mu` again → ep2 blocks until the FIFO release; the "ep1 still in flight" assertion fires |
| `TestWorkspaceAilangBuildSingleFlight` | 2.2 | **MUT-NO-SINGLEFLIGHT:** each caller builds → spawn count 2 |
| `TestWorkspaceFailedBuildRetriesOffClock` | 2.3, 2.6 | **MUT-SYNC-RETRY:** today's synchronous retry → calls block on the FIFO (assert-returns fires) and the spawn count grows per call. **MUT-RETRY-STORM:** a retry per call → spawn count > 2. **MUT-LINE-PER-CALL:** log count > 2. |
| `TestWorkspaceFailedBuildRecovers` | 2.4 | **MUT-NEGATIVE-CACHE-FOREVER:** failure cached with no retry → still 1 name after recovery |
| `TestWorkspaceCloseJoinsRetry` | 2.5 | **MUT-LEAK:** the retry not tied to the daemon context → `Close` returns with a live stub group (the process-group probe fires) |

## 8. CI gate list (every milestone)

```sh
export PATH=/opt/homebrew/bin:$PATH AILANG_BIN=~/.pinned-ailang/ailang   # gate interpreter v0.41.0
./scripts/verify_ail.sh                        # no .ail change in this row; still run
go vet ./...                                   # compile fence for _test.go (go build is not one)
./scripts/verify_go.sh                         # go build ./...; host/evidence manifest; go test ./... -count=1; -race -timeout 8m under 600 s
go test ./host/daemon/ -count=1 -race -v -run '^(TestWorkspace|TestExecMaxTimeout)' | tee $S/m<N>-ws.log   # every new/amended name: grep "^--- PASS: <name>"
# real srt (CI's linux step runs this; locally with WORLD_EXEC_SRT_NODE_MODULES=<pinned 0.0.78 install>):
go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -v -run '^(TestExecSrtMCPEndToEnd|TestExecSrtMCPGrantAndBudget|TestExecSrtReplayIsByteEqual)$'   # each must print --- PASS
```

A bare `-run` with no match exits 0, so every new test name gets a `--- PASS: <name>` grep in the same call. The v0.52.1 tool binary is not needed: every new test uses the `/bin/sh` stub. The real-tool tests (`realToolBin`, `WORLD_TOOL_AILANG_BIN`) are unchanged by this row.

## 9. Conflict Surface

| File | Change | Interaction |
|---|---|---|
| `host/daemon/workspace.go` | registry split, `execMu`, per-episode build state, budget `var` | Rows 140/141 code. No open row edits `registry()`. |
| `host/daemon/daemon.go` | `Close` joins workspace retries | Row 153 added `stopTemplateMaintenance` here. Same pattern; append, not replace. |
| `host/daemon/workspace_layout_test.go`, `workspace_test.go` | 4 assertions amended (AC1.5), harness generalised (`publishReadTool` → effect/id) | `publishReadTool` is shared with row 153's template note (`workspace_test.go:499–507`): keep that block intact. |
| `docs/QUICKSTART.md` | operator-lines paragraph | S7; QUICKSTART is executed verbatim. Only prose changes, no command. |
| `host/broker/handlers.go` | **untouched** | The `HandlerTimeoutError` label (R1) is deliberately left to a follow-up row, so this row does not collide with row 166 (broker head/tail timeout flake). |
| — | — | **Row 93** (World arm under load): after M2 a broken episode no longer adds 3 s per call (F11). |

## 10. Non-goals, residuals, follow-up rows (proposed; the controller files)

- **R1 — `HandlerTimeoutError` names the wrong bound** (F6, V5, V10). Proposed row *w-handler-timeout-names-effective-bound* (hygiene). `runBounded` should report the effective deadline (`min(parent, execTimeout)`), or say which bound fired. Every subprocess handler shares this, which is why it is not folded in here. It misled row 140 M4's diagnosis.
- **R2 — the first build is still on the first caller's clock.** A fresh episode's first call of *any* transition waits for the first AILANG build (≤ `workspaceHandlerBudget`), because `registry()` cannot know which family the transition needs: `BinderFor(episodeID, caps)` carries no declared effects (V11). The healthy cost is ~0.1 s. The worst case is 3 s once per episode, which still erodes F11's derivation once. Proposed row *w-binder-builds-only-declared-families*: pass the descriptor's declared effects to `BinderFor` (it changes the coordinator type; there are 18 test call sites, from `grep -rn 'Binder: ' host --include='*_test.go' | wc -l`). That removes the residual entirely.
- **R3 — client-visible cause for R8.** The A2A/MCP R8 text does not say *why* the handler is absent. That is row 159 / ailang#1602's territory for MCP, and an A2A message change for every R8, so it is out of scope.

## 11. Risks

- **A background goroutine in the daemon.** It is bounded by `workspaceHandlerBudget`, cancelled and joined by `Close` (AC2.5), and has row 153's precedent (V18). A test that removes temp dirs while a retry runs would flake; AC2.5's join is what prevents that, and every test daemon is closed via `t.Cleanup` (`newWSDaemon`).
- **One extra R8 after an operator fixes a layout refusal** (c4 trade-off). It is documented in QUICKSTART (AC1.8), and the next call succeeds.
- **Exec available while AILANG tools are not** is a new mixed state. It is intended (the row's acceptance bar). The operator line still names the AILANG cause on every attempt.

## 12. Upstream asks

None.

## 13. UNMEASURED

- **U1 — The real-tool path.** All probes used the `/bin/sh` stub. That the real v0.52.1 `policy-tool summary` can exceed 3 s under load rests on README F5's prose (V9: no raw log). The code path is the same for any `NewAilangToolHandler` error (V6 vs V7), so the fix does not depend on the trigger.
- **U2 — The F11 erosion in a live run.** The 7 s handler figure is arithmetic on the constants (V15), not a measured Dispatch with a stalled registry.
- **U3 — Linux.** Probes ran on darwin/arm64. The lock and early-return findings are static (V1–V3); only the 3.0 s/2.82 s timings are platform-measured.
- **U4 — Whether the daemon logs an `EffectFailedError`'s cause anywhere.** Not checked. It only bears on rejected option b2.

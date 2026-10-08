# w-exec-availability-decoupled-from-ailang-policy — a failed AILANG tool-handler build unbinds `Workspace.Exec`, which never uses it (row 152)

**Status:** PLANNED, revision r3. **Quorum r1 BLOCKED → r2 BLOCKED 2/3 → r3 carve-out (verbatim fixes), routed to planner.** In r2, gemini-3-1-pro passed and oc-glm-5-3 and oc-kimi-k3 rejected; gpt6-1-sol was unreachable. Neither rejecting seat disputed the direction. r3 applies their fixes verbatim (§15). Earlier history: Designed in iteration 242 (designer lane claude-opus-5-5, unattended). **Quorum r1 BLOCKED → r2.** In r1, 3 of 3 present seats rejected (oc-glm-5-3, oc-kimi-k3, gemini-3-1-pro) and gpt6-1-sol was absent. All three objections landed on M2's background retry or on the probe's provenance; M1 drew no objection on its direction. r2 drops the background retry, keeps the retry synchronous, and commits the probe in-repo. §14 maps each objection to its change.
**Queue:** World mission row 152 `w-exec-availability-decoupled-from-ailang-policy` (clause 4). It is next after row 153 (D-WORLD-65, re-ordered by D-WORLD-68 = A).
**Estimated:** (r2) ~0.6 d of host Go and tests (M1 ~0.35 d, M2 ~0.2 d, M3 ~0.05 d). **Dependencies:** none. No `.ail` changes, no pin changes, no upstream gate.
**Measured base:** worktree `sprint/row152-exec-availability` at `fe7c2a4`. Gate interpreter `~/.pinned-ailang/ailang` = `AILANG v0.41.0`. `go1.26.6 darwin/arm64`. (r2) The probe is committed in-repo under `design_docs/verification/world-row152-design/` (`$P` below):
- the source, as `zz_probe152_test.go.txt` so it is never compiled;
- two run logs;
- a README with the exact re-run command.

The compiled copy under `host/daemon/` is removed after each run. No production code or test file changes in this commit.

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
| F10 | A typed failure handler for the eight AILANG names would be worse for the client than R8: (1) it commits a `failed` invocation; (2) it debits the budget; (3) the record has no field for the cause. R8 refuses before the plan with no record and no debit. **(r3)** r2's claim (4), "the wire body is `{}`", is **dropped**. Two of the eight AILANG transitions plan with `finish: true` (`check.ail:233`, `builtins_search.ail:369–370`), so their failed effect goes to a finish phase, not to the `{}` one-effect rule (V13). Reasons (1)–(3) are verified independently. | V12, V13 |
| F11 | The row-153 plan cap is derived as `invokeDeadline − HandlerCap − HandlerHeadroom = 4 s`, assuming nothing before the plan uses the invoke deadline. A 3 s registry stall before the plan breaks that assumption. | V15 |

## 3. Verification Log

Shorthand: `W=` this worktree, `P=design_docs/verification/world-row152-design` (in-repo, r2). The probe ran as `PATH=/opt/homebrew/bin:$PATH AILANG_BIN=~/.pinned-ailang/ailang go test ./host/daemon/ -count=1 -v -run '^TestZZProbe152'`, with `$P/zz_probe152_test.go.txt` copied to `host/daemon/zz_probe152_test.go` and removed afterwards (exact steps in `$P/README.md`). It was run twice on base `fe7c2a4`:
- run 1 (r1): `$P/probe-run1.log`, filtered;
- run 2 (r2): `$P/probe.log`, full `-v` output, ending `rc=0`.

Both probe tests printed `--- PASS` in both runs. The probe's tool binary is a `/bin/sh` stub shaped like `fakeToolBin` (`workspace_test.go:95`): it answers `--version` with `AILANG v0.52.1`, and its `policy-tool` summary either `sleep 30`s (mode `sleep`) or prints `{"ok":false}` and exits 1 (mode `fail`). Exec is configured with the existing `stubExecSandbox` + `execProfile` (`workspace_exec_test.go:26,50`, command `where` = `pwd -P`). **Rule applied:** each null result is paired with a positive control.

| V | Command (trimmed) | Reading | Control / scope |
|---|---|---|---|
| V1 | `sed -n 304,336p host/daemon/workspace.go` | Returns `broker.Registry{}` at `:309` (`!w.enabled()`), `:313` (`episodeRoot` not ok), `:318` (`sandboxRoot` err) and `:328` (`episodeHandler` err; an `*episodeRefusal` prints its own line `:324`, anything else prints `workspace tools unavailable for episode %q: %v` `:326`). `reg[broker.EffectWorkspaceExec] = w.execHandler(episodeID, epRoot)` is at `:334`, after all four. | `grep -c 'return broker.Registry{}' host/daemon/workspace.go` → 4 |
| V2 | `sed -n 342,369p host/daemon/workspace.go` | `execHandler(episodeID, epRoot)` reads `w.exec`, `w.execH`, `w.root`, `w.stateDir`, `w.exec.{sandbox,operatorHome,maxOutputBytes}`, `profileFor(episodeID)`. `NewExecHandler` gets `Worktree: epRoot`. It takes `w.mu` (`:346`). A construction error → `execFailedHandler{err}` (`:361`). | Its arguments contain no `sandbox`/`policyPath`/`cacheDir` (they are locals of `episodeHandler`, `:401–414`) |
| V3 | `sed -n 395,437p host/daemon/workspace.go` | `episodeHandler` takes `w.mu` for its whole body (`:396–397`). It runs, in order: `MkdirAll` policy/cache dirs, `linkPackageCache`, `checkLockCoverage`, `RenderEpisodePolicy`, `writePolicy`, then `NewAilangToolHandler` under `context.WithTimeout(context.Background(), workspaceHandlerBudget)` (`:422`). It writes to the cache `w.handler[...]` only after success (`:432–435`). Every error returns before that. | The success cache is proved by the existing `TestWorkspaceRegistryRefusesEpisodesOutsideTheGrammarAndRoot` (`summaries == 1` after two calls, `workspace_test.go:218–222`) |
| V4 | `grep -rn 'workspaceHandlerBudget\|ailangToolExecTimeout' host --include='*.go'` | `workspaceHandlerBudget = 3 * time.Second` (`workspace.go:44`), used once (`:422`). `ailangToolExecTimeout = 10 * time.Second` (`handlers_ailang.go:48`), the handler's default `execTimeout` (`:328`). `NewAilangToolHandler` → `loadSummary` → `policyToolWith(ctx, …)` (`handlers_ailang.go:331,338`; `handlers_ailang_run.go:376`). | static |
| V5 | `sed -n 50,52p;153,157p;204,232p host/broker/handlers.go` | `runBounded`: `runCtx, cancel := context.WithTimeout(ctx, bounds.execTimeout)` (`:156`). On expiry it returns `&HandlerTimeoutError{Timeout: bounds.execTimeout}` (`:208,229,231`). `Error()` = `fmt.Sprintf("%v after %s", ErrHandlerTimeout, e.Timeout)` (`:51`). A parent deadline that fires first is reported with the child's 10 s. | V6 is the runtime positive |
| V6 | probe `TestZZProbe152SummaryTimeout/sleep` → `$P/probe-run1.log`, `$P/probe.log` | run 1: call 1 `elapsed=3.007s names=[] execBound=false summaries=1`, call 2 `elapsed=3.018s … summaries=2`. Run 2: `3.009s`/`3.015s`, same names and counts. Operator log, twice: `ailang-worldd: workspace tools unavailable for episode "ep1": broker: policy summary: broker: handler subprocess timed out after 10s (stderr "")`. | **Control:** `execHandler("ep1", epRoot)` called directly in the same daemon → `{"exit_code":0,…,"stdout":".../ws/ep1\n",…}`, `err=<nil>` |
| V7 | probe `…/fail` → same files | run 1 `elapsed=17ms`/`18ms`, run 2 `26ms`/`14ms`; `names=[]`, `execBound=false`, `summaries=2`. Log: `… broker: policy summary: broker: handler subprocess failed: exit status 1 (output "{\"ok\":false}\n")`. | Same control → `exit_code 0`. Deterministic trigger of the same path with no wall-clock bound |
| V8 | probe `TestZZProbe152LockCoupling` → same files: `registry("ep1")` in a goroutine (stub summary sleeps), 200 ms later `execHandler("ep2", ep2Root)` | `execHandler(ep2) waited 2.822s` (run 1) / `2.82s` (run 2) `behind ep1's summary` | Cross-episode coupling through `w.mu`. The 2.82 s matches 3 s − the 0.2 s head start. |
| V9 | `cd design_docs/verification/world-row140-m4; grep -c 'policy summary' serve.log; grep -c 'registered handler' serve.log; grep -c go1 serve.log; grep -c 'plan phase' serve.log`; `grep -c 'policy summary' README.md`; `sed -n 1,11p log-tail.txt` | serve.log: **0 / 0 / 0**, control `plan phase` **2** (it has 3 lines in all). README: 1 (finding F5's prose). `log-tail.txt` shows two `go1 workspace-exec [Workspace.Exec ok]` entries (#1, #2). | The row's citation of `serve.log` is **false**; the sequence is reported only in README F5. |
| V10 | `git show 7664d3f:host/daemon/workspace.go \| grep -n workspaceHandlerBudget`; `git show 7664d3f:host/broker/handlers.go \| grep -n 'HandlerTimeoutError{Timeout'` | At M4's build: `workspaceHandlerBudget = 3 * time.Second` (`:43`), and `HandlerTimeoutError{Timeout: bounds.execTimeout}` at `:208,229,231,313,328,330` | F5's "timed out at 10 s" was the V5 mislabel. The bound was 3 s. |
| V11 | `sed -n 267,320p host/coordinator/coordinator.go` | `Dispatch` evaluates `c.cfg.Binder(call.EpisodeID, grants)` as an argument of `transitionreg.Bind` (`:316–317`) on **every** call. For effectful calls this happens after `lockEpisode` (`:283–289`). Production `Binder` = `d.binder` (`daemon.go:681`) → `broker.OpenBinder(…, d.workspace.registry(episodeID))` (`daemon.go:745`). | `grep -n 'workspace.registry' host/daemon/*.go \| grep -v _test` → only `:745` (plus the comment at `:449`) |
| V12 | `sed -n 274,279p host/coordinator/effectful.go`; `sed -n 37,44p host/coordinator/errors.go`; `grep -n 'no handler for' host/projection/projection.go`; `sed -n 21,32p $GOMODCACHE/github.com/sunholo-data/ailang@v0.47.2/serveapi/protocol/envelope.go` | R8 runs before `loadSourceAndWorld` and the plan: `if missing := bound.Unhandled(); len(missing) != 0 { return …&EffectsUnsupportedError{…} }`. Operator text: `coordinator: transition %q declares effects with no registered handler %v; effectful invocation is unavailable`. **A2A wire:** `-32603 "transition declares an effect this daemon has no handler for"` (`projection.go:462`). **MCP wire:** `CallbackMessage` sends anything that is neither capacity, `DeadlineExceeded` nor `Canceled` as `"host callback failed"`. | The existing table test `projection_test.go:828` pins the A2A string. `setools_e2e_test.go:825` pins `host callback failed` on the MCP wire (for R14). |
| V13 | `sed -n 291,308p host/broker/broker.go`; `sed -n 19,30p host/broker/record.go`; `sed -n 177,203p;243,245p host/coordinator/effectful.go`; (r3, corrected) `grep -n 'effectPlan(' packages/se-tools/se_tools/*.ail` and `sed` of each call's finish argument; `transitions.json` has no finish field (V23) | A handler error writes `EffectRecord{…Allowed: true, Failed: true…}` with `BudgetAfter: decision.Remaining` (a debit), appends outcome `failed`, and returns `EffectFailedError`. `EffectRecord`'s fields are Effect, Scope, Cost, BudgetBefore/After, Allowed, Failed, Denial, RequestRef, ResultRef: **no cause**. The coordinator maps this to `StatusFailed` and commits. With one effect and no finish, the body is `{}` (`:244–245`, "denied or failed: the world block carries the status and record"). **(r3) Finish is a per-plan field.** It is `true` in `check.ail:233` (`Ailang.Check`) and `builtins_search.ail:369–370` (`Ailang.Discover` `builtins_list`). It is `false` in `read.ail:168`, `write.ail:170–171`, `edit.ail:174–176`, `run.ail:387–388`, `cli.ail:309`, `examples_search.ail:170` and `exec.ail:332`. So `{}` holds for 6 of the 8 AILANG transitions only, and F10(4) is dropped. | `Denial` is the positive control: a record field for a reason exists only for denials |
| V14 | `python3 -I -c` over `packages/se-tools/transitions.json` printing `id`, `declaredEffects` | 9 entries. `workspace-exec` → `[Workspace.Exec]`. `ailang-read` → `[Workspace.Read]`, `ailang-write`/`ailang-edit` → `[Workspace.Write]`, `ailang-check` → `[Ailang.Check]`, `ailang-run` → `[Ailang.Run, Ailang.RunEnv, Ailang.RunNet]`, `builtins-search`/`examples-search` → `[Ailang.Discover]`, `ailang-cli` → `[Ailang.CLI]` | No transition declares names from both families |
| V15 | `sed -n 36,49p;118,125p host/coordinator/effectful.go`; `grep -n 'invokeDeadline *=' host/daemon/daemon.go` | `PlanPhaseBudget = 4s` is derived in the comment as "20 s − HandlerCap − HandlerHeadroom = 20 − 10 − (2 + 4) = 4 s". `HandlerCap = 10s`, `HandlerHeadroom = Finish (2s) + PostEffect (4s)`. `handlerBudget` = `min(HandlerCap, time.Until(deadline) − HandlerHeadroom)`. `invokeDeadline = 20s` (`daemon.go:99`). | A 3 s stall before the plan leaves the handler `min(10, 20 − 3 − 4 − 6) = 7 s` when the plan uses its full budget. That is arithmetic on the constants, not a measured run (UNMEASURED U2). |
| V16 | `grep -n 'len(reg) != 0\|want empty' host/daemon/*_test.go`; `sed -n 889,915p host/daemon/workspace_layout_test.go` | Existing "empty registry" assertions: `workspace_layout_test.go:159` (module root), `:586` (package cache plants), `:874` (lock coverage), `:905` (FIFO lock, `n != 0`); `workspace_test.go:203` (episode grammar), `:277`, `:290` (one flag alone). | M1 changes the first four and must leave the last three unchanged (§6) |
| V17 | `sed -n 95,121p host/daemon/workspace_test.go`; `sed -n 449,538p host/daemon/workspace_test.go`; `sed -n 128,130p host/daemon/workspace.go` | Seams that already exist: `fakeToolBin` (a shell stub archived as the tool binary, which counts `summaries`/`dispatches`), `publishReadTool` + `readPlanRunner` (a coordinator wired to `d.binder` with a fake plan runner, so no interpreter is needed), and `stubExecSandbox`. `execStartupHook` adjusts only `ExecStartupConfig` at startup (`configureExec`), **not** handler construction, so it cannot drive the summary. | V6/V7 used these seams directly |
| V18 | `sed -n 988,991p host/daemon/daemon.go`; `grep -n 'func (d \*Daemon) stopTemplateMaintenance' host/daemon/templates.go` | `Close()` = `d.stopTemplateMaintenance(); return d.store.Close()`. Row 153's background rebuild (`templates.go:158`) is the precedent for a daemon-owned goroutine that `Close` joins. | static. (r2) Cited only for the withdrawn c4: r2 adds no goroutine. |
| V19 | `sed -n 984,989p docs/QUICKSTART.md` | "Each is written once per refused call, and the episode's tools stay refused (**every declared effect fails, as for a missing worktree**) until you fix the cause". This sentence becomes false after M1 (S7). | static |
| V20 | `sed -n 229,245p;283,288p .github/workflows/ci.yml`; `sed -n 290,330p scripts/verify_go.sh` | CI runs `verify_ail.sh`, then `verify_go.sh` (`go build ./...`, `go test ./... -count=1`, `go test ./... -count=1 -race -timeout 8m` under a 600 s kill), a linux srt step (`go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -v -run '^(…TestExecSrtMCPEndToEnd|TestExecSrtMCPGrantAndBudget|TestExecSrtReplayIsByteEqual)$'` with a `--- PASS` grep per name) and the subprocess-cleanup step | §8 is derived from these |
| V21 | (r3) `grep -n 'w\.mu\b\|w\.mu\.' host/daemon/workspace*.go \| grep -v _test`; `grep -n 'handler  *map\|execH  *map\|mu  *sync' host/daemon/workspace.go`; `grep -n 'w.handler\[\|w.execH\[' host/daemon/workspace*.go`; `sed -n 342,369p;395,437p host/daemon/workspace.go` | Fields: `mu sync.Mutex` (`:99`), `handler map[string]episodeTool // episode id -> constructed handler` (`:100`), `execH map[string]episodeTool // episode id -> constructed exec handler` (`:101`). Two non-test holders. (1) `execHandler` `:346–347` guards `w.execH`: lookup `w.execH[episodeID]` with `cached.root == epRoot` (`:348`), store `w.execH[episodeID] = episodeTool{root: epRoot, h: h}` (`:367`), lazy map init (`:364–366`). (2) `episodeHandler` `:396–397` guards `w.handler`: lookup `w.handler[episodeID]` with `cached.root == sandbox` (`:398`), store `w.handler[episodeID] = episodeTool{root: sandbox, h: tool}` (`:435`). It also covers the whole build: mkdirs, `linkPackageCache`, `checkLockCoverage`, policy render/write, and `NewAilangToolHandler` (`:401–430`). Config fields (`root`, `bin`, `exec`, `stateDir`, `moduleRoot`, `packageCache`, …) are written only at startup and are read without the lock. **Exec already has its own map, `execH`, keyed by episode id.** glm's "one shared map" scenario does not hold; the coupling is the shared *mutex* only. | `workspace_layout.go:300` mentions `w.mu` in a comment only (the grep's positive control for a non-holder hit) |
| V22 | (r3) `sed -n 140,150p host/daemon/workspace.go` | `if !w.enabled() { return fmt.Errorf("the --exec-* flags need the workspace tools (--workspace-root and --tool-ailang-bin): " + "workspace-exec runs in an episode worktree") }` (`:142–145`). It runs after the all-flags-empty early return (`:139–141`). | Premise of §5(a)'s `!enabled()` row |
| V23 | (r3) V14 re-run: `python3 -I -c` printing `len(json.load(open('packages/se-tools/transitions.json')))`, then each `id` and any key containing `inish` | `len 9`. ids: ailang-read, ailang-write, ailang-edit, ailang-check, ailang-run, builtins-search, examples-search, ailang-cli, workspace-exec. No entry has a finish key. | 9 is the total; 8 AILANG + 1 exec |

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
- b2. A typed failure handler, like `execFailedHandler`. **Rejected on measurement (F10).** Each call would commit a `failed` invocation and debit a budget unit for a host-side construction fault. The record has no field for the cause (V13). (r3: r2's "the wire body is `{}`" is dropped, because `ailang-check` and `builtins-search` plan a finish phase; reasons (1)–(3) carry the rejection.) It is strictly more state for no added information.
  - `execFailedHandler` is right for exec because exec's construction failure is per-profile configuration with no R8-level alternative (row 140's choice). It is not a template for a transient subprocess failure.
- b3. A refusal-body handler (`{"ok":false,"refused":…}`, like `ExecRefusalHandler`). Rejected: the eight AILANG effects have eight output schemas, and none declares a `refused` arm on construction. This would mean `.ail` and output-contract changes for a host fault. If client-visible cause text is wanted, that is row 159 / ailang#1602's job for every R8 refusal, not this row's.

### (c) Retry policy for a failed AILANG build

Today's retry policy (c1) has two costs. r2 separates them.
- **Cross-family / cross-episode cost.** The build holds `w.mu`, so exec construction and every other episode wait behind it (F5, V8).
- **Per-call cost in a persistently broken episode.** Every Binder open re-runs the build on the caller's clock, up to 3.0 s, pure calls included (F4, F8, F11).

- c1. Today, unchanged: a synchronous retry under `w.mu`. **Rejected.** It carries both costs.
- c2. Cache the failure until restart. Rejected: a transient timeout would disable the episode's AILANG tools for the daemon's lifetime.
- c3. Cache the failure for a cool-down of N seconds. Rejected: N would be a guessed number (row 153's rule: derive bounds or do not add them).
- c4. (r1's choice, **withdrawn in r2**) Retry in the background after a failure. Quorum r1 showed it is non-deterministic in two ways:
  - **For clients.** A client that fixes a layout error and resends can still get a stale R8 if it arrives before the background build ends, and it cannot know when that is (gemini).
  - **For tests.** An asynchronously arriving operator line races the row-141 tests' exact-log assertions (glm).
- **c5 (chosen, r2). Synchronous retry on the caller's clock, as today, but never under a shared lock.**
  - Every Binder open of an episode whose AILANG handler is not built runs the build itself, and prints one operator line if it fails. These are today's per-call semantics, so a fixed layout error takes effect on the very next call. No existing log-line or summary-count assertion changes (§6, AC2.4).
  - **Per-episode single flight.** Concurrent `registry()` calls for the same episode share one in-flight build and its result. Each caller whose call ends without the AILANG names still prints its own line, as today.
  - **Locks.** `w.mu` guards only the build-state map, and the new `execMu` guards only the exec cache. Neither is held across the summary subprocess. Exec construction for any episode, and the AILANG build of *other* episodes, never wait behind a stuck build (fixes F5/V8 deterministically).
  - **Exec stays bound throughout** (M1).
  - **Cost kept, stated.** In a *persistently* broken episode, each call to that episode, `workspace-exec` included, still waits up to `workspaceHandlerBudget` (3 s) before its plan. In that case F11's erosion of the derived handler budget remains. This is the same cost as today, now confined to the broken episode. The deterministic fix is to build only the families a transition declares. r2 proposes that as one follow-up row covering both this cost and R2 (§10, R2).

### (d) How the test drives the summary to time out

No wall-clock ratio and no reliance on the 3 s constant. Two pieces:
1. **A seam.** `workspaceHandlerBudget` becomes a package `var`, following the `geteuid`/`killGroup` pattern. A test lowers it.
2. **A stub that never answers.** The stub's summary blocks on a FIFO (or `sleep 30`). The outcome is then a timeout on any machine, however fast or slow: the stub never answers, so only the budget can end the build.

Assertions are on outcomes: names bound, R8 vs committed, summary spawn counts, and error kind (`errors.As(err, *broker.HandlerTimeoutError)` reaching the operator line). They are never on elapsed time. Where a test needs "while the build is in flight", the FIFO holds it open until the test releases it, so the state is reached by construction, not by a sleep.

**Recommendation:** a1 + b1 + c5 + d, in two code landings plus records (§6). The `HandlerTimeoutError` label becomes a follow-up row (R1).

**Why is this not a package? (S3)** The handler registry is host-boundary wiring (S2). The broker binds effect names to Go handlers that spawn subprocesses, and no `.ail` code can express or replace that. This row adds no surface and no kernel growth. It narrows an existing host coupling and moves no policy out of AILANG: the confinement stays in AILANG's own policy layer.

## 6. Milestones

Each milestone is one PR, CI-green on its own; gates in §8.

### M1 — two families, two bindings (~0.35 d)

Files: `host/daemon/workspace.go`, `host/daemon/workspace_test.go`, `host/daemon/workspace_layout_test.go`, `host/daemon/workspace_exec_test.go` (or a new `workspace_family_test.go`), `docs/QUICKSTART.md`.

- **AC1.1** `registry()` binds `Workspace.Exec` whenever `enabled && episodeRoot ok`. It binds the eight AILANG names only when the AILANG handler is built. On a `sandboxRoot` or `episodeHandler` failure it prints the same operator line as today and returns exactly `[Workspace.Exec]`.
- **AC1.2** (r3, concrete structures) **`execMu` guards the pre-existing, distinct map `w.execH`** (V21). No map is split, and **no entries move**: `execH` already holds only exec handlers (`execHandler` `:348`, `:367`), and `w.handler` only AILANG tools (`:398`, `:435`). r3 changes only which mutex guards `w.execH`. `execHandler` takes `execMu` and never `w.mu`. `w.mu` keeps guarding `w.handler` (and, from M2, the AILANG build-state map of AC2.1).
- **AC1.3** `workspaceHandlerBudget` becomes a `var` (a test seam, with a comment). The production value stays 3 s.
- **AC1.4 — the row's pin, end to end.** The fixture is all existing seams (V17):
  - a coordinator over `d.binder`;
  - (r3) published `ws.exec` and `ws.read` descriptors, using `publishReadTool` generalised to take the effect and id (`ws.exec`: Access and declared effect `Workspace.Exec@worktree`);
  - a fake plan runner that answers one `Workspace.Exec` `{"command":"where"}` effect;
  - exec configured with `stubExecSandbox` + `execProfile`;
  - a tool stub whose summary never answers;
  - the budget seam set low.

  Then:
  - (i) Dispatch `ws.exec` in ep1 **commits**: `exit_code 0`, and `stdout` is the ep1 worktree path.
  - (ii) Dispatch `ws.read` in ep1 → `*coordinator.EffectsUnsupportedError` naming `[Workspace.Read]`, with zero read plans run.
  - (iii) The operator log has `workspace tools unavailable for episode "ep1": broker: policy summary:`.

  Control: the same fixture with a healthy `fakeToolBin` → both calls commit.
- **AC1.5** Row 141 refusals leave exec intact. Four existing tests are amended from "want empty" to "want exactly `[Workspace.Exec]`": `workspace_layout_test.go:159, :586, :874`, and the FIFO test's `n != 0` → `n != 1` at `:905`. Their `summaries == 0` and exact-log assertions stay unchanged. **(r3) Each amended refusal test (:159, :586, :874, :905) must additionally drive one exec call through the bound handler and assert `exit_code == 0` and `stdout ==` the episode worktree, using the existing `stubExecSandbox` + `execProfile` seams. Membership alone is insufficient.**
- **AC1.6** Episode-level refusals stay empty, unedited: `workspace_test.go:203` (all 18 `refusedEpisodes`), `:277`, `:290`, and `TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect`.
- **AC1.9** (r3) A test that fails `episodeHandler` before `writePolicy` returns exactly `[Workspace.Exec]`, **and** a dispatched exec call commits (AC1.4's coordinator fixture, `exit_code 0`, stdout = ep1 worktree).
  - The fault is planted after the daemon starts, as **a regular file at `<stateDir>/policies`** (replacing the directory if startup created it), so `os.MkdirAll` at `workspace.go:404` fails with "not a directory". This works for root as well.
  - The chmod-0555 variant kimi suggested is **not used.** Root ignores mode bits, and the repo's chmod-based tests skip under root (`host/archive/environment_test.go:163`, `archive_test.go:503`: `if os.Geteuid() == 0 { t.Skip(...) }`). A skip would make the AC vacuous under root (S6). The plant-a-file variant needs no skip.
  - Kimi's optional `cache-refusal` probe mode is skipped: AC1.5's amended package-cache tests (`:586`, `:874`) already run an exec call under a planted row-141 refusal.
- **AC1.7** Exec construction never waits on an AILANG build. The test holds ep1's summary on a FIFO, with the budget seam set to the test's own bounded context, so only the FIFO can end the build. It then calls `execHandler("ep2", …)` and `registry("ep2")`'s exec entry. Both must return **while ep1's `registry` goroutine is still blocked**: the goroutine's done-channel is asserted not closed after they return. Then the test releases the FIFO. **(r3)** In addition, **`TestWorkspaceTwoEpisodeConcurrentRegistryExecRace`**: 2 episodes × N goroutines (N = 8) call `registry(ep)` and `execHandler(ep, root)` concurrently, with one episode's summary failing and the other healthy. The assertions are made after the goroutines join, never on timing: every call returns, exec is bound in every registry, the healthy episode's eight names are all bound, and each episode's `execH` entry is one handler (pointer-equal across calls). It must pass under `-race`, and its name is in §8's grep list.
- **AC1.8 (S7)** Update the QUICKSTART operator-lines paragraph (V19): the episode's AILANG tools stay refused, and `workspace-exec` keeps working (it needs only the worktree). Update the file header and the `ailangToolEffects`/`workspaceEffects` comments in `workspace.go`, which say "an EMPTY registry" and "ALWAYS bound".

### M2 — (r2) single-flight synchronous build, no shared lock across the subprocess (~0.2 d)

Files: `host/daemon/workspace.go`, tests. **`daemon.go` is untouched:** r2 adds no goroutine and no `Close` change.

- **AC2.1** Per-episode build state: `{built tool | in-flight done-chan + result}`, keyed as today by `(episodeID, sandbox)`. `w.mu` guards only the map and is released before `NewAilangToolHandler` runs. (r3, concrete structures) The build state is a **new map `w.ailangBuild map[string]*ailangBuild`** (episode id → `{root string; done chan struct{}; tool broker.Handler; err error}`), guarded by `w.mu`. The success cache stays `w.handler`, also under `w.mu`, with unchanged key and `cached.root == sandbox` check (`:398`, `:435`). `w.execH` stays under `execMu` (AC1.2). No entries move between maps. **A failure is not cached** (today's semantics): the in-flight entry is cleared when the build returns, so the next Binder open builds again, synchronously.
- **AC2.2** Same-episode single flight. Two concurrent `registry("ep1")` calls, made while the first build is held on a FIFO, spawn **one** summary (stub spawn count == 1). Both get the build's result. If the build fails, each prints its own operator line (2 lines, one per call, as today).
- **AC2.3** Other episodes never wait. While ep1's build is held on a FIFO (budget seam set to the test's bounded context), `registry("ep2")` with a healthy summary returns all nine names. The test asserts this **before** releasing the FIFO, by checking that ep1's done-channel is still open. This is the AILANG-family counterpart of AC1.7. **(r3)** AC1.7's `TestWorkspaceTwoEpisodeConcurrentRegistryExecRace` is re-run on the M2 head under `-race` with the new build-state map in play, and is in §8's grep list.
- **AC2.4 — assertions M2 must leave unchanged** (each test is re-run with `--- PASS` grepped; none is edited by M2):
  - `workspace_test.go:207`: `summaries == 0` for all refused episodes.
  - `workspace_test.go:220`: `summaries == 1` after two `registry("ep1")` calls (the success cache).
  - `workspace_test.go:270`: `summaries + dispatches == 0` with one flag.
  - `workspace_test.go:558`: `plans == 0 && dispatches == 0 && summaries == 0` for refused episodes.
  - `workspace_layout_test.go:128`: `summaries == 1`.
  - `workspace_layout_test.go:162–165`: `summaries == 0`, `dispatches == 0`, and **exactly 1** operator line per refused call.
  - `workspace_layout_test.go:596`: `log.String() == want`, **exact**.
  - `workspace_layout_test.go:599`: `summaries == 0`.
  - `workspace_layout_test.go:869`: healthy registry with `log.Len() == 0`.
  - `workspace_layout_test.go:874–878`: `summaries == 0` and **exactly 1** line, exact or containing `want`.

  M1 changes only the registry-size halves of `:159`, `:586`, `:874` and `:905` (AC1.5). The log-line and summary-count halves listed above are unchanged by both milestones.
- **AC2.5** Per-call retry preserved. With a summary stub that fails (non-ok, no timing), three sequential `registry("ep1")` calls produce: 3 summary spawns, 3 operator lines, and `[Workspace.Exec]` each time. Then the stub is switched to ok (a flag file), and the 4th call binds all nine names with no extra line. This pins c5's "a fix takes effect on the next call" (gemini's objection) and kills any cached or async failure.

### M3 — records (~0.05 d)

- **AC3.1** Charter row 152 → LANDED. The controller files the §10 follow-up rows. Note in the M4 evidence README that F5's "10 s" was the 3 s budget mislabelled (V10), with a pointer to this doc (the evidence file itself is not rewritten).

## 7. Test plan (what each test kills)

Every row is mutation-proven by the executor: apply the named production mutation, see the named assertion go red, restore.

| Test (new unless noted) | AC | Kills the mutation |
|---|---|---|
| `TestWorkspaceExecSurvivesPolicySummaryTimeout` | 1.4 | **MUT-ALL-OR-NOTHING:** restore `return broker.Registry{}` on `episodeHandler` failure → (i) red with R8 `[Workspace.Exec]`. **MUT-BIND-AILANG-ON-FAIL:** bind the eight names to `execFailedHandler`-style failures → (ii) red (commits `failed` instead of R8). |
| amended `TestWorkspaceModuleRootSymlinkOrMissingIsR8`, package-cache plant / lock-coverage / FIFO tests | 1.5 | **MUT-SANDBOX-GATES-EXEC:** exec bound only after `sandboxRoot` succeeds → registry empty, red. **MUT-EXEC-USES-SANDBOX:** `execHandler(ep, sandbox)` instead of `epRoot` → the module-root test's exec `where` stdout ≠ ep1 worktree (the amended test runs one `where` call). |
| unedited `TestWorkspaceRegistryRefusesEpisodesOutsideTheGrammarAndRoot`, `TestWorkspaceRegistryNeedsBothFlags`, `TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect` | 1.6 | **MUT-EXEC-BEFORE-EPROOT:** exec bound before the `episodeRoot` check → `registry("escape")` non-empty, red |
| amended AC1.5 tests' exec call (r3) | 1.5 | **MUT-BIND-REFUSAL-HANDLER:** on refusal, bind `Workspace.Exec` to `execFailedHandler`/`ExecRefusalHandler` instead of the real handler → membership still passes but `exit_code == 0` fires |
| `TestWorkspaceExecSurvivesPolicyDirFault` (r3) | 1.9 | **MUT-EARLY-FAULT-EMPTIES:** keep `return broker.Registry{}` for failures *before* `NewAilangToolHandler` (mkdir/policy write) while splitting only the subprocess failure → registry empty, and the exec Dispatch answers R8 |
| `TestWorkspaceTwoEpisodeConcurrentRegistryExecRace` (r3) | 1.7, 2.3 | **MUT-EXECH-UNLOCKED:** `execHandler` reads/writes `w.execH` with no lock → `-race` reports a DATA RACE. **MUT-WRONG-LOCK:** `w.execH` guarded by `w.mu` in one path and `execMu` in the other → DATA RACE. **MUT-DOUBLE-BUILD:** cache store dropped → not pointer-equal |
| `TestWorkspaceExecNeverWaitsOnAilangBuild` | 1.7 | **MUT-SHARED-LOCK:** `execHandler` takes `w.mu` again → ep2 blocks until the FIFO release; the "ep1 still in flight" assertion fires |
| `TestWorkspaceAilangBuildSingleFlight` | 2.2 | **MUT-NO-SINGLEFLIGHT:** each caller builds → spawn count 2 |
| `TestWorkspaceAilangBuildDoesNotBlockOtherEpisodes` | 2.3 | **MUT-LOCK-ACROSS-SUBPROCESS:** `w.mu` held across `NewAilangToolHandler` (today) → `registry("ep2")` blocks until the FIFO release; the "ep1 still in flight" assertion fires |
| `TestWorkspaceFailedBuildRetriesEachCall` | 2.5 | **MUT-NEGATIVE-CACHE:** failure cached → spawn count 1 and no recovery on call 4. **MUT-ASYNC-RETRY:** r1's background retry → call 4 is still `[Workspace.Exec]` (the stale R8). **MUT-LINE-PER-ATTEMPT:** lines ≠ 3. |
| the AC2.4 list (unedited) | 2.4 | Any change to per-call log or summary semantics reds an existing exact assertion. This row is an instrument-health control, not a load-bearing claim (S6). |

## 8. CI gate list (every milestone)

```sh
export PATH=/opt/homebrew/bin:$PATH AILANG_BIN=~/.pinned-ailang/ailang   # gate interpreter v0.41.0
./scripts/verify_ail.sh                        # no .ail change in this row; still run
go vet ./...                                   # compile fence for _test.go (go build is not one)
./scripts/verify_go.sh                         # go build ./...; host/evidence manifest; go test ./... -count=1; -race -timeout 8m under 600 s
go test ./host/daemon/ -count=1 -race -v -run '^(TestWorkspace|TestExecMaxTimeout)' | tee ~/.ailang/state/<iter>/m<N>-ws.log
#   then grep "^--- PASS: <name>" for every new/amended name; (r3) the list includes
#   TestWorkspaceTwoEpisodeConcurrentRegistryExecRace and TestWorkspaceExecSurvivesPolicyDirFault
# real srt (CI's linux step runs this; locally with WORLD_EXEC_SRT_NODE_MODULES=<pinned 0.0.78 install>):
go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -v -run '^(TestExecSrtMCPEndToEnd|TestExecSrtMCPGrantAndBudget|TestExecSrtReplayIsByteEqual)$'   # each must print --- PASS
```

A bare `-run` with no match exits 0, so every new test name gets a `--- PASS: <name>` grep in the same call. The v0.52.1 tool binary is not needed: every new test uses the `/bin/sh` stub. The real-tool tests (`realToolBin`, `WORLD_TOOL_AILANG_BIN`) are unchanged by this row.

## 9. Conflict Surface

| File | Change | Interaction |
|---|---|---|
| `host/daemon/workspace.go` | registry split, `execMu`, per-episode single-flight build state (no lock across the subprocess), budget `var` | Rows 140/141 code. No open row edits `registry()`. |
| `host/daemon/daemon.go` | (r2) **untouched**: no background retry, so no `Close` change | none |
| `host/daemon/workspace_layout_test.go`, `workspace_test.go` | 4 assertions amended (AC1.5), harness generalised (`publishReadTool` → effect/id) | `publishReadTool` is shared with row 153's template note (`workspace_test.go:499–507`): keep that block intact. |
| `docs/QUICKSTART.md` | operator-lines paragraph | S7; QUICKSTART is executed verbatim. Only prose changes, no command. |
| `host/broker/handlers.go` | **untouched** | The `HandlerTimeoutError` label (R1) is deliberately left to a follow-up row, so this row does not collide with row 166 (broker head/tail timeout flake). |
| — | — | **Row 93** (World arm under load): a *persistently* broken episode still adds up to 3 s per call (F11). After M2 that cost no longer spreads to other episodes; the per-call cost is removed only by the R2 follow-up row. |

## 10. Non-goals, residuals, follow-up rows (proposed; the controller files)

- **R1 — `HandlerTimeoutError` names the wrong bound** (F6, V5, V10). Proposed row *w-handler-timeout-names-effective-bound* (hygiene). `runBounded` should report the effective deadline (`min(parent, execTimeout)`), or say which bound fired. Every subprocess handler shares this, which is why it is not folded in here. It misled row 140 M4's diagnosis.
- **R2 — (r2, widened) AILANG builds still run on the caller's clock.** Two costs remain after this row:
  - (i) A fresh episode's first call of *any* transition waits for the first AILANG build (healthy ~0.1 s, worst case 3 s).
  - (ii) **(r2)** In a *persistently* broken episode, every call, `workspace-exec` included, waits up to `workspaceHandlerBudget` (3 s) before its plan, because c5 keeps today's synchronous per-call retry (F4, F8). In both cases F11's derived handler budget is eroded by the stall (to 7 s with a full-budget plan, V15).

  Both exist because `registry()` cannot know which family the transition needs: `BinderFor(episodeID, caps)` carries no declared effects (V11). Proposed row *w-binder-builds-only-declared-families* (clause 4): pass the descriptor's declared effects to `BinderFor` and build only the families they name. Then a `workspace-exec` call never runs the AILANG build, and both costs vanish deterministically, with no async state. It changes the coordinator type, with 18 test call sites (`grep -rn 'Binder: ' host --include='*_test.go' | wc -l`), which is why it is a separate row.
- **R3 — client-visible cause for R8.** The A2A/MCP R8 text does not say *why* the handler is absent. That is row 159 / ailang#1602's territory for MCP, and an A2A message change for every R8, so it is out of scope.

## 11. Risks

- **(r2) The per-call stall in a broken episode is kept, not fixed** (c5). A client calling `workspace-exec` in an episode whose AILANG build keeps timing out gets a correct result 3 s later than in a healthy episode. Under heavy load that cost could combine with a slow plan (F11). The cost is confined to that episode (AC2.3) and is the same as today. R2's follow-up row removes it. r1's background retry removed it too, but at the price of the two non-determinisms quorum r1 found (§14).
- **Exec available while AILANG tools are not** is a new mixed state. It is intended (the row's acceptance bar). The operator line still names the AILANG cause on every attempt.

## 12. Upstream asks

None.

## 13. UNMEASURED

- **U1 — The real-tool path.** All probes used the `/bin/sh` stub. That the real v0.52.1 `policy-tool summary` can exceed 3 s under load rests on README F5's prose (V9: no raw log). The code path is the same for any `NewAilangToolHandler` error (V6 vs V7), so the fix does not depend on the trigger.
- **U2 — The F11 erosion in a live run.** The 7 s handler figure is arithmetic on the constants (V15), not a measured Dispatch with a stalled registry.
- **U3 — Linux.** Probes ran on darwin/arm64. The lock and early-return findings are static (V1–V3); only the 3.0 s/2.82 s timings are platform-measured.
- **U5 — (r2) The probe on another machine.** It is now re-runnable from `$P/README.md`. Only this rig's two runs are recorded.
- **U4 — Whether the daemon logs an `EffectFailedError`'s cause anywhere.** Not checked. It only bears on rejected option b2.

## 14. Quorum r1 resolution (r2)

Round 1: oc-glm-5-3, oc-kimi-k3 and gemini-3-1-pro rejected, and gpt6-1-sol was absent (N−1). Full JSON is in the controller's quorum artifact `w-exec-availability-decoupled-from-ailang-policy-2026-10-08T00-39-39Z.json`. The controller checked all three premises and found each true. None disputed M1's direction.

| # | Objection (seat) | Doc change |
|---|---|---|
| Q-1 | M2's AC2.6 (one line per attempt) plus an auto-started background retry would add a second line that arrives asynchronously and races the row-141 tests' exact-log assertions, which AC1.5 claims are unchanged. The doc never listed the assertions M2 invalidates, and `awaitBuilds()` existed only for AC2.4 (glm). | Background retry **dropped** (§5 c4 withdrawn, c5 chosen). The retry is synchronous per call, one line per failed call, exactly as today. New AC2.4 lists every summary-count and log-line assertion in `workspace_test.go`/`workspace_layout_test.go` by file:line, all unchanged; AC2.5 pins per-call retry. `awaitBuilds()`, the recovery line and the `Close` join are removed. |
| Q-2 | A client that fixes a layout error and retries gets a stale R8 if it arrives before the background build finishes, and cannot know when that is (gemini). | Same change. The next call after a fix runs the build itself (AC2.5, killing MUT-ASYNC-RETRY). The per-call 3 s cost this reintroduces in a persistently broken episode is stated in §5(c), §9 and §11 and moved to R2's follow-up row, the deterministic fix. |
| Q-3 | V6–V8 rested on a probe banked only on one machine; nobody else could re-run it (kimi). | The probe source, both run logs and a README with the exact re-run command are committed at `design_docs/verification/world-row152-design/`. The probe was re-run for r2 (`probe.log`, `rc=0`: 3.009/3.015 s, 26/14 ms, 2.82 s; same outcomes as run 1). V6–V8 cite both files. |

Re-estimate: M2 shrinks from ~0.35 d to ~0.2 d (no goroutine lifecycle), for a total of ~0.6 d.

## 15. Quorum r2 resolution (carve-out r3)

Round 2: gemini-3-1-pro PASS; oc-glm-5-3 and oc-kimi-k3 rejected; gpt6-1-sol unreachable. Neither rejecting seat disputed the direction, so r3 applies their fixes verbatim in substance and adds nothing else.

| # | Fix (seat) | Applied |
|---|---|---|
| G-1 | V21: enumerate every `w.mu` holder and the state it guards, the `w.handler` key, and whether/where `execHandler` caches (glm) | V21 added. Its reading (also the controller's measurement) is that exec **already has its own map** `w.execH` (`:101`, `:348`, `:367`). The "one shared map" scenario does not hold; only the mutex is shared. The fix is applied anyway. |
| G-2 | V22: `sed -n 140,150p` verifying the `--exec-*` refusal (glm) | V22 added (`:142–145`) |
| G-3 | V23: re-run V14 printing `len(transitions)` (glm) | V23 added: `len 9`, no finish key |
| G-4 | AC1.2/AC2.1 name the concrete structures and moved entries (glm) | AC1.2: `execMu` guards the pre-existing `w.execH`; no split, no moved entries. AC2.1: new `w.ailangBuild` map under `w.mu`; `w.handler` unchanged |
| G-5 | A named two-episode concurrent registry+execHandler stress test under `-race`, in §8's grep list (glm) | `TestWorkspaceTwoEpisodeConcurrentRegistryExecRace` in AC1.7/AC2.3, §7 (3 mutants) and §8 |
| K-1 | AC1.5: each amended refusal test drives one exec call and asserts `exit_code == 0` and stdout == worktree (kimi) | AC1.5 amended verbatim; §7 row (MUT-BIND-REFUSAL-HANDLER) |
| K-2 | AC1.9: fail `episodeHandler` before `writePolicy` → exactly `[Workspace.Exec]` and an exec call commits (kimi) | AC1.9 added. It plants a *file* at `<stateDir>/policies`, which works for root, instead of chmod: chmod tests in this repo skip under euid 0 (`host/archive/environment_test.go:163`, `archive_test.go:503`), which would be vacuous (S6). §7 row (MUT-EARLY-FAULT-EMPTIES). The optional probe mode is skipped (AC1.5 covers the package-cache refusal with an exec call). |
| K-3 | Correct V13's command to cover the AILANG tools' finish phases; drop F10(4) if any declares one (kimi) | V13 corrected. `ailang-check` and `builtins-search` plan `finish: true`, so **F10(4) and §5(b)'s `{}` clause are dropped**; b2's rejection rests on (1)–(3). |
| M-1 | AC1.4 states both descriptors are published (gemini, non-blocking) | AC1.4: "published `ws.exec` and `ws.read` descriptors, using `publishReadTool` generalised to take the effect and id" |

Estimate is unchanged at ~0.6 d. The added tests reuse the M1/M2 fixtures.

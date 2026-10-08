# Sprint plan — `w-exec-availability-decoupled-from-ailang-policy` (row 152, clause 4)

**Design doc:** [`w-exec-availability-decoupled-from-ailang-policy.md`](w-exec-availability-decoupled-from-ailang-policy.md) (r3 carve-out, `2bb27ad`)
**Sprint base:** `fe7c2a4` (dev) + the three design commits = branch `sprint/row152-exec-availability` @ `2bb27ad`. The design commits touch only `design_docs/` (the probe is committed as `.txt`, never compiled).
**Planner:** mission iteration 242, claude-opus-5-5 (unattended). Planner logs: `~/.ailang/state/mission-world-iter242/planner-base-*.log`.
**Estimate:** ~0.65 d (design: 0.6 d; +0.05 d for D3/D6/D7 below) over 3 commits in ONE PR, run as **one executor run** (§6). **Risk:** medium. M2 replaces a mutex-held build with a single-flight build; the only lifecycle change is in-memory.
**Scope:** host Go and tests, one QUICKSTART paragraph (S7), comments. No `.ail`, no pin, no `tools/launchd/*`, no `host/broker/*` change, no `design_docs/world-mission*.md`.

---

## 0. The short version

`registry()` returns an empty registry when the AILANG handler cannot be built, so `Workspace.Exec` (which never reads the AILANG inputs) is unbound and `workspace-exec` is refused with R8. The work:
- **M1** (C1, C2): bind exec iff `enabled && episodeRoot ok`; bind the eight AILANG names iff the AILANG build succeeds; `execMu` guards `w.execH`; `workspaceHandlerBudget` becomes a test seam.
- **M2** (C3): the AILANG build runs single-flight per episode with `w.mu` released across the subprocess; failures stay uncached (per-call retry, today's semantics).

The planner re-ran every gate on the pristine base: **all green, no pre-existing reds** (§2). The design's citations were re-read (§3). Seven findings change how the executor must work (§3 D2–D8): the most important are
1. **AC1.7's `registry("ep2")` half cannot pass at M1** (ep2's AILANG build still takes `w.mu`, held by ep1's build until M2) — it moves to M2.
2. **`exec_e2e_srt_test.go:102–104` (`execSpawns`) reads `w.execH` under `w.mu`**, and is not in the design's file list. After AC1.2 it would read under the wrong mutex — a race the CI linux `-race` leg can see (srt tests run there, `WORLD_EXEC_SRT_NODE_MODULES` is exported to `$GITHUB_ENV` before `verify_go.sh`).
3. **The amended refusal tests must reuse the one `registry()` result** for their exec call: a second `registry()` call prints a second operator line and reds the unchanged `n == 1` / exact-log halves (:165, :596, :878).

---

## 1. Environment (on every command)

```sh
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang            # AILANG v0.41.0, commit 24ee108 (gate interpreter)
export WORLD_EXEC_SRT_NODE_MODULES=$HOME/.ailang/state/mission-world-iter241/srt/node_modules   # srt 0.0.78 (pin.json), used locally by the real-srt leg
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
L=$HOME/.ailang/state/mission-world-iter242/exec    # executor logs; never /tmp
```

`go version` = `go1.26.6 darwin/arm64` (= CI's `GOTOOLCHAIN`). The v0.52.1 tool binary is **not** needed: every new test uses a `/bin/sh` stub. **Instrument rules:**
- Redirect to a file, then `echo rc=$?` on the next line. Never `| tail; echo $?`, never `PIPESTATUS` (zsh).
- A `-run` that matches nothing exits 0 and `go test` without `-v` hides SKIPs: every named test is proven by `grep -- "^--- PASS: <name> "` in the same call.
- `go build` is not a compile fence for `_test.go`: use `go vet ./...`.
- Never end a turn waiting on a background gate; poll its log in a bounded `date +%s` loop (`verify_go.sh` takes ~11 min).

---

## 2. Baseline — measured by the planner on the pristine tree (`2bb27ad`; code == `fe7c2a4`)

| Gate | Command | rc | Wall | Notes / log |
|---|---|---|---|---|
| `.ail` | `./scripts/verify_ail.sh` | 0 | 23 s | 16 identities, 40 named tests, 444 se-tools named tests, PUB011 ratchet ok (`planner-base-verify_ail.log`) |
| vet | `go vet ./...` | 0 | — | empty output (`planner-base-vet.log`) |
| full Go gate | `./scripts/verify_go.sh` | 0 | ~11 min | binary-hygiene, decision ledger, race known-positive control, build, evidence manifest (37), plain `./...` (daemon 172 s, verifygate 155 s, broker 107 s), race `./...` (daemon **191 s**, verifygate 169 s, broker 145 s, coordinator 51 s); `✓ go gate PASSED` (`planner-base-verify_go.log`) |
| focused | `go test ./host/daemon/ -count=1 -race -v -run '^(TestWorkspace\|TestExecMaxTimeout)'` | 0 | 16 s | **25 `--- PASS`**, 2 pre-existing SKIP (`TestWorkspaceSandboxIsTheEpisodeNotTheRoot`, `TestWorkspaceReadThroughProductionWiring`: need `WORLD_TOOL_AILANG_BIN`) (`planner-base-ws.log`) |
| real srt (3 named in design §8) | `go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -timeout 600s -v -run '^(TestExecSrtMCPEndToEnd\|TestExecSrtMCPGrantAndBudget\|TestExecSrtReplayIsByteEqual)$'` | 0 | 14 s | 3/3 `--- PASS` on darwin with the env above (`planner-base-srt3.log`) |
| gofmt | `gofmt -l host/ cmd/` | 0 | — | empty |
| personal email | `bash scripts/check_no_personal_email.sh` | 0 | — | |

**Pre-existing reds: none. Pre-existing flakes observed locally: none.** The CI-only `TestExecSrtMCP*` flake (D-WORLD-68 history, row 166 signature "host callback timed out") is the controller's to count; a red at a commit boundary locally is attributable to that commit.

### 2.1 Gate list per commit boundary (derived from `.github/workflows/ci.yml` go-verify job + `scripts/verify_go.sh`)

```sh
mkdir -p $L
./scripts/verify_ail.sh                                   >$L/cN-ail.log 2>&1; echo rc=$?
go vet ./...                                              >$L/cN-vet.log 2>&1; echo rc=$?
./scripts/verify_go.sh                                    >$L/cN-go.log  2>&1; echo rc=$?   # background + bounded poll; build, evidence, plain, race (8m under 600 s)
go test ./host/daemon/ -count=1 -race -v -run '^(TestWorkspace|TestExecMaxTimeout)' >$L/cN-ws.log 2>&1; echo rc=$?
for t in <every name in the commit's PASS list, §5>; do grep -q -- "^--- PASS: $t " $L/cN-ws.log || echo "NOPASS $t"; done   # must print nothing
# CI's real-srt step (all 20 names, ci.yml "Row 140 Workspace.Exec real-srt tests"), run locally with the env above:
tests="TestExecSrtGoTest TestExecSrtExitStatusAndArgv TestExecSrtRefusalsSpawnNothing TestExecSrtConfinementMatrix TestExecSrtGitHeadIsNotWritable TestExecSrtTimeoutKillsTheGroup TestExecSrtPycacheStaysOutOfTheWorktree TestExecSrtProfileEnvNeverReachesTheHostNode TestExecSrtStartupProbePasses TestExecSrtProbeRefusesAPackageOnlySandbox TestExecProbePassThroughShimFailsEveryConfinementArm TestExecProbeNoopShimFailsArm1 TestExecProbeZeroShimFailsArm5 TestExecProbeFailingToolchainFailsArm7 TestExecProbeMissingClientFailsTheRefusalArms TestQuickstartSection10ProfilesLoad TestExecStartupRefusalTable TestExecSrtMCPEndToEnd TestExecSrtMCPGrantAndBudget TestExecSrtReplayIsByteEqual"
go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -timeout 600s -v -run "^($(echo $tests | tr ' ' '|'))\$" >$L/cN-srt.log 2>&1; echo rc=$?
for t in $tests; do grep -q -- "^--- PASS: $t " $L/cN-srt.log || echo "NOSRTPASS $t"; done
# the three srt tests that read w.execH (execSpawns) ALSO under -race (D3; CI's verify_go race leg runs them on linux):
go test ./host/daemon/ -count=1 -race -v -run '^(TestExecSrtMCPEndToEnd|TestExecSrtMCPGrantAndBudget|TestExecSrtReplayIsByteEqual)$' >$L/cN-srt-race.log 2>&1; echo rc=$?
gofmt -l host/ cmd/                                       # must print nothing
git diff --numstat fe7c2a4 HEAD | awk '$1=="-"'           # must print nothing (no binaries)
bash scripts/check_no_personal_email.sh; echo rc=$?
```

Not run per commit (unaffected, recorded so the omission is deliberate): the fork-lock verbose step (`host/verifygate`, linux-only guard), the row-24 cleanup verbose step (`procbound/proctest/capsule/broker` — no file in those packages changes), the compiler reproducer, bench smoke/claims, census/gate suites, floor harness, kernel-contract and ratchet steps (no `.ail`). Report per-package wall for `host/daemon` in the race leg at every boundary: row 154 (CI race leg ≈ 540–556 s of a 600 s cap) makes any growth a finding, even when green. Expected growth is < 3 s.

---

## 3. Design citations re-verified at `2bb27ad`

**TRUE as cited:** `workspace.go` :44 (`const workspaceHandlerBudget = 3 * time.Second`), :99–101 (`mu`, `handler`, `execH`), :139–145 (V22), :309/:313/:318/:328 (four `return broker.Registry{}`), :334 (exec bound last), :346–347 and :396–397 (the two `w.mu` holders), :348/:367 (`execH`), :398/:435 (`handler`), :404 (`MkdirAll`), :422 (budget use); `workspace_layout.go:300` (comment-only `w.mu` mention); `workspace_test.go` :95 (`fakeToolBin`), :142 (`newWSDaemon`), :203/:207/:220/:270/:277/:290/:558, :451 (`readPlanRunner`), :478 (`publishReadTool`); `workspace_exec_test.go` :26 (`stubExecSandbox`), :50 (`execProfile`); `workspace_layout_test.go` :128, :159, :162–165, :586, :596, :599, :869, :874–878; `daemon.go:745` (`registry` in `binder`); `host/archive/environment_test.go:163`, `archive_test.go:503` (`if os.Geteuid() == 0 { t.Skip }`). No `t.Parallel()` exists in `host/daemon` tests (`grep -c` = 0), so package-var seams are safe.

**Drift / defects (recorded; the executor applies the planner decisions in §7, it does not edit the design):**

| # | Design says | Measured | Consequence |
|---|---|---|---|
| D1 | V16/AC1.5: the FIFO test's `n != 0` is at `workspace_layout_test.go:905` | it is at **:907** (`TestWorkspacePackageCacheLockFIFOIsRefusedPromptly`, func at :889) | cosmetic |
| D2 | AC1.7 (M1): `execHandler("ep2")` **and `registry("ep2")`'s exec entry** return while ep1's build is blocked | at M1 `episodeHandler` still takes `w.mu` for its whole body (:396–397), so `registry("ep2")` blocks behind ep1's build until M2 releases the lock | PD3: the `registry("ep2")` arm moves to C3 (`TestWorkspaceAilangBuildDoesNotBlockOtherEpisodes`, AC2.3) |
| D3 | AC1.2 file list: `workspace.go` and the workspace tests | `exec_e2e_srt_test.go:100–113` `execSpawns` does `r.d.workspace.mu.Lock()` then reads `r.d.workspace.execH["ep1"]` | must switch to `execMu` in C1 (PD4). These tests run on linux CI under `-race` (verify_go's race leg, srt env exported) and locally with the §1 env |
| D4 | §7: MUT-EXEC-USES-SANDBOX killed by "the module-root test's exec `where`" | in `TestWorkspaceModuleRootSymlinkOrMissingIsR8` `sandboxRoot` **fails**, so there is no sandbox to pass. The tests where a sandbox exists and differs from the worktree are the lock-coverage (:814, module root `tools`) and FIFO (:889, module root `tools`) tests | PD5: those two amendments are the killers; their `where` stdout must equal `<root>/ep1`, not `<root>/ep1/tools` |
| D5 | AC1.4 / §5(d): `errors.As(err, *broker.HandlerTimeoutError)` "reaching the operator line" | an operator line is text; `errors.As` needs the error value | PD6: assert `errors.As` on a direct `d.workspace.episodeHandler("ep1", epRoot)` error, plus the line `Contains "timed out"` |
| D6 | AC1.4: "`publishReadTool` generalised to take the effect and id", publishing `ws.exec` and `ws.read` | `publishReadTool` commits genesis at revision 0 and `CompareAndSetRegistryHead(…, hashref.HashRef{}, …)`; calling it twice fails | PD7: one revision carrying both descriptors |
| D7 | AC1.5: "each amended refusal test must additionally drive one exec call" with its log/summary halves unchanged | every `registry()` call on a refused episode prints one more operator line (`execJSON`, `workspace_layout_test.go:750`, calls `registry` again) | PD8: the exec call uses the handler from the **same** `reg` the test already obtained |
| D8 | AC1.7/AC2.2/AC2.3: the summary is "held on a FIFO" | a FIFO with several blocked readers (MUT-NO-SINGLEFLIGHT spawns two) releases non-deterministically, and opening it for write fails `ENXIO` once the reader was killed; AC2.2 also has no deterministic signal that the second caller has joined | PD1 (hold file) + PD2 (`waiters` counter) |
| D9 | §9: `daemon.go` untouched | `daemon.go:741–743` (`binder` doc: "empty, so R8 refuses every declared effect, unless …") becomes false after M1 | PD9: comment-only edit in C1 |
| D10 | §8 lists 3 srt names | CI's step runs **20** names with a PASS grep each (ci.yml "Row 140 Workspace.Exec real-srt tests") | §2.1 uses all 20 |
| D11 | §7: MUT-EXECH-UNLOCKED / MUT-WRONG-LOCK killed by `-race` | the race detector sees only interleavings that happen; `registry()`'s `w.mu` acquisitions add happens-before edges. Probabilistic, not deterministic | PD10: drill at `-race -count=5`; a mutant that survives 5 runs is reported as SURVIVED, never strengthened silently. MUT-DOUBLE-BUILD (pointer equality) is the deterministic tooth |

**Root/euid (AC1.9, brief item 3).** The repo's chmod-based tests skip under euid 0 (`host/archive/environment_test.go:163`, `archive_test.go:503`); `host/daemon`'s package-cache tests fake uid via the `geteuid` seam (`workspace_layout.go:115`, `withUID` at `workspace_layout_test.go:316`) and never chmod for refusal. CI's go-verify job runs on `ubuntu-latest` as the unprivileged runner user (the workflow needs `sudo` for `apt-get`/`sysctl`, ci.yml "Row 140 srt confinement matrix"). The chosen trigger — a **regular file at `<stateDir>/policies`** — makes `os.MkdirAll` (`workspace.go:404`) fail with `ENOTDIR` for every uid on darwin and linux, so no skip is needed. `<stateDir>/policies` is created only by `episodeHandler` (:401) and `configureRunCaps` (:460, only with `--run-*` flags), so it does not exist after `New` in the test's config; the test still `os.RemoveAll`s it before planting and `Lstat`s the plant as regular (control).

---

## 4. Shared test fixtures (added in C2, used by C2 and C3)

All in `host/daemon/workspace_family_test.go` (new) unless stated.

- **`gatedToolBin(t, logDir) string`** — `fakeToolBin`'s script plus per-episode summary control. For a summary request it: reads stdin (`req=$(cat)`), appends `x` to `$log/summaries` **and** `$log/summaries.$ep` where `ep=$(basename "$(pwd -P)")` (the summary runs with cwd = sandbox; these tests use no module root, so `ep` is the episode id); then `while [ -e "$log/hold.$ep" ]; do sleep 0.02; done`; then, if `$log/fail.$ep` exists, `printf '{"ok":false}\n'; exit 1`; else the healthy `fakeToolBin` answer (`fs_sandbox` = `$(pwd -P)`). Release a hold = `os.Remove(hold)`. A hold never released + a low budget = "never answers".
- **`withHandlerBudget(t, d time.Duration)`** — sets the `workspaceHandlerBudget` var, restores in `t.Cleanup` (serial, as `withUID`).
- **`publishTools(t, d, interp, src string, tools ...wsTool)`** with `type wsTool struct{ ID, Effect string }` — in `workspace_test.go`, the body of today's `publishReadTool` with one descriptor per tool in **one** revision (Access `{Effect, worktree, 0}`, DeclaredEffects `{Effect, worktree, 1}`). `publishReadTool(t, d, interp, src)` becomes `publishTools(t, d, interp, src, wsTool{"ws.read", broker.EffectWorkspaceRead})` with its signature and its row-153 `CheckSource` block (`:499–507`) intact; its two callers (`workspace_test.go:543,638`, `mcp_plan_timeout_test.go:39`) are untouched.
- **`execPlanRunner`** — like `readPlanRunner`, answers the plan phase with `{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"Workspace.Exec","scope":"worktree","cost":1,"payload":{"command":"where"}}],"finish":false,"result":null}` and counts plans.
- **`execDispatchCall(t, d, episode, task, now)`** — `readCall`'s shape with a `Workspace.Exec@worktree` grant (budget 5) and `SkillID: "ws.exec"`, `Input: {}`. (Name chosen because `execCall` exists, `workspace_exec_test.go:66`.)
- **Two coordinators** over the same `d.store` and `d.binder`: one with `execPlanRunner`, one with `readPlanRunner` (PD7).
- Operator logs that can be written from several goroutines use **`syncLog`** (`setools_e2e_test.go:792`), never `bytes.Buffer`.

---

## 5. Commits

Conventions: "MUT" = apply exactly this edit on top of the commit's own tree, run the named test(s) `-count=1 -v` (`-race` where stated), require `--- FAIL: <name>` with the stated reason, then `git checkout -- <file>` and require `git diff --quiet` (byte identity; for a mutation in a file the commit also changed, `git stash`-free: re-checkout from HEAD and compare `git diff --stat` = empty). Record each MUT (the diff line, the failing assertion's output line, the revert check) in the PR body's mutation ledger. **A MUT that does not red is a finding: report it; never weaken or strengthen the mutation to make it red.**

### C0 — re-baseline (no commit)

Run §2.1 on HEAD before any edit; compare with §2. Mismatch → STOP and report.

### C1 — M1 production: two families, two bindings — AC1.1, AC1.2, AC1.3, AC1.5, AC1.6, AC1.8

**Files:** `host/daemon/workspace.go`, `host/daemon/daemon.go` (comment only, PD9), `host/daemon/workspace_layout_test.go`, `host/daemon/exec_e2e_srt_test.go` (lock only, PD4), `docs/QUICKSTART.md`.

**Steps:**
1. `workspace.go`: `workspaceHandlerBudget` → `var` with a comment ("a test seam; production value 3 s; tests set it with `withHandlerBudget`, never concurrently with a build").
2. Add `execMu sync.Mutex` beside `execH` (comment: guards `execH` only). `execHandler` takes `execMu`, never `w.mu`. `w.mu` keeps guarding `w.handler`.
3. `registry()`: after `enabled` and `episodeRoot`, build the AILANG family (sandboxRoot + episodeHandler, printing the same operator lines as today), then **always** set `reg[broker.EffectWorkspaceExec] = w.execHandler(episodeID, epRoot)`; set the eight AILANG names only on success. PD11: AILANG family first, exec second (preserves today's line order).
4. Comments: file header (:3–15), `ailangToolEffects` (:49–54), `workspaceEffects` (:60–64), `registry` doc (:304–306) — the AILANG names are bound only when the handler is built; `Workspace.Exec` whenever the worktree resolves. `daemon.go:741–743` likewise.
5. `exec_e2e_srt_test.go:102–103`: `execMu` instead of `mu`.
6. Amend the four refusal tests (AC1.5). Each: configure exec **after** `New` with `d.workspace.exec = &workspaceExec{profile: execProfile(t), sandbox: stubExecSandbox(t, f.stateDir)}` (as `workspace_exec_test.go:75` does); keep the single `registry("ep1")` call; change "want empty" to "want exactly `[Workspace.Exec]`" (`registryNames(reg)` == `[Workspace.Exec]`); after the unchanged summary/log assertions, call `execCall(t, reg[broker.EffectWorkspaceExec])` and assert `exit_code == 0` and `strings.TrimSpace(stdout) == filepath.Join(f.root, "ep1")` (PD8). Sites: `:159` (`TestWorkspaceModuleRootSymlinkOrMissingIsR8`, every subtest), `:586` (`TestWorkspacePackageCacheRefusesUnprovisionedRegistry`, the 3 plant subtests), `:874` (`TestWorkspacePackageCacheCoversLock`, the refusal arms; the healthy arms keep `len(reg) == 0` → fail), `:907` (FIFO: the goroutine sends the registry, not its length; `n != 0` → names == `[Workspace.Exec]`).
7. `docs/QUICKSTART.md:984–985`: the operator-lines sentence says the episode's **AILANG tools** stay refused until the cause is fixed, and `workspace-exec` keeps working (it needs only the worktree). Prose only.
8. AC1.6: `workspace_test.go:203` (all 18 `refusedEpisodes`), `:277`, `:290`, `TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect` stay **unedited** (`git diff fe7c2a4 -- host/daemon/workspace_test.go` shows no hunk in those functions).

**PASS list (C1):** `TestWorkspaceModuleRootSymlinkOrMissingIsR8 TestWorkspacePackageCacheRefusesUnprovisionedRegistry TestWorkspacePackageCacheCoversLock TestWorkspacePackageCacheLockFIFOIsRefusedPromptly TestWorkspaceRegistryRefusesEpisodesOutsideTheGrammarAndRoot TestWorkspaceRegistryNeedsBothFlags TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect TestWorkspaceModuleRootIsTheSandbox TestWorkspaceExecConfiguredBindsTheEpisodeHandler TestWorkspaceExecUnbuildableHandlerFails TestWorkspaceExecUnconfiguredSpawnsNothing` + the 3 srt `-race` names.

**MUTs (C1):**

| MUT | Edit | Must red |
|---|---|---|
| MUT-SANDBOX-GATES-EXEC | in `registry`, `return broker.Registry{}` on `sandboxRoot` error (bind exec only after it) | `TestWorkspaceModuleRootSymlinkOrMissingIsR8` (registry empty) |
| MUT-ALL-OR-NOTHING-C1 | `return broker.Registry{}` on `episodeHandler` error | `TestWorkspacePackageCacheRefusesUnprovisionedRegistry`, `…CoversLock`, `…LockFIFO…` (registry empty) |
| MUT-EXEC-USES-SANDBOX | `w.execHandler(episodeID, sandbox)` on the paths where `sandbox` is known (D4) | `TestWorkspacePackageCacheCoversLock` refusal arms and `…LockFIFO…`: stdout `<root>/ep1/tools` ≠ `<root>/ep1` |
| MUT-BIND-REFUSAL-HANDLER | on AILANG failure bind `Workspace.Exec` to `execFailedHandler{err}` | all four amended tests: `execCall` returns the error / `exit_code` assertion fires |
| MUT-EXEC-BEFORE-EPROOT | bind exec (with `filepath.Join(w.root, episodeID)`) before the `episodeRoot` check, returning `reg` on refusal | `TestWorkspaceRegistryRefusesEpisodesOutsideTheGrammarAndRoot` (`registry("escape")` non-empty: unconfigured exec binds `ExecUnconfiguredHandler`) |

### C2 — M1 new tests — AC1.4, AC1.7 (exec half), AC1.9, race test

**Files:** `host/daemon/workspace_family_test.go` (new), `host/daemon/workspace_test.go` (`publishTools` generalisation only, PD7).

**Tests:**
- **`TestWorkspaceExecSurvivesPolicySummaryTimeout`** (AC1.4). `gatedToolBin`, hold file for ep1 never released, `withHandlerBudget(t, 300*time.Millisecond)`, exec configured, `publishTools` with `ws.exec` + `ws.read`, two coordinators, `ErrorLog: &syncLog{}`. (i) exec coordinator `Dispatch(execDispatchCall(…"ep1"…))` → nil error, output `exit_code` 0, `stdout` == `<root>/ep1`; (ii) read coordinator → `*coordinator.EffectsUnsupportedError` with `Effects == [Workspace.Read]` and `readPlanRunner.plans == 0`; (iii) log `Contains` `workspace tools unavailable for episode "ep1": broker: policy summary:` and `timed out`; (iv) PD6: `_, err := d.workspace.episodeHandler("ep1", <root>/ep1)` → `errors.As(err, &hte)` with `var hte *broker.HandlerTimeoutError`. **Control subtest** (`/healthy`): same fixture, no hold → both Dispatches commit (read output contains `fake:<root>/ep1`).
- **`TestWorkspaceExecSurvivesPolicyDirFault`** (AC1.9). Healthy `gatedToolBin`; after `New`: `os.RemoveAll(<stateDir>/policies)`, `writeFile(<stateDir>/policies, "plant\n")`, `Lstat` it regular (control). `registry("ep1")` names == `[Workspace.Exec]`, `summaries == 0`; exec coordinator Dispatch commits (`exit_code` 0, stdout = ep1 worktree); log `Contains "not a directory"`.
- **`TestWorkspaceExecNeverWaitsOnAilangBuild`** (AC1.7, exec half). `withHandlerBudget(t, 5*time.Second)` (PD12), hold for ep1. Goroutine: `registry("ep1")` → closes `ep1Done`. Precondition wait (bounded by `boundedTestContext`, poll 10 ms): `summaries.ep1 == 1`. Then `execHandler("ep2", <root>/ep2)` (direct) → `execCall` exit 0 in ep2; then assert `ep1Done` **not** closed (`select { case <-ep1Done: t.Fatal(…) default: }`); remove the hold; `<-ep1Done`; ep1's registry has all nine names.
- **`TestWorkspaceTwoEpisodeConcurrentRegistryExecRace`** (AC1.7 r3 / AC2.3). `fail.ep1` present, ep2 healthy, exec configured, `syncLog`. 2 episodes × 8 goroutines calling `registry(ep)` and 8 × `execHandler(ep, root)`, released together by a closed start channel; `sync.WaitGroup` join; then: every registry has `Workspace.Exec`; every ep2 registry has all nine; every ep1 registry is exactly `[Workspace.Exec]`; per episode every returned exec handler `==` `w.execH[ep].h` read under `execMu` (pointer-equal); `summaries.ep2 == 1`. No timing assertion. Run `-race`.

**PASS list (C2):** C1's list + `TestWorkspaceExecSurvivesPolicySummaryTimeout TestWorkspaceExecSurvivesPolicyDirFault TestWorkspaceExecNeverWaitsOnAilangBuild TestWorkspaceTwoEpisodeConcurrentRegistryExecRace` (and `^    --- PASS: TestWorkspaceExecSurvivesPolicySummaryTimeout/healthy `).

**Tests-first proof (once, uncommitted):** with C2's test file on top of C1's tests but **C0's `workspace.go`** (`git checkout fe7c2a4 -- host/daemon/workspace.go` plus a one-line `var` for the budget so it compiles), `…PolicySummaryTimeout`, `…PolicyDirFault` and `…NeverWaitsOnAilangBuild` must FAIL; restore and require `git diff --quiet`.

**MUTs (C2):**

| MUT | Edit | Must red |
|---|---|---|
| MUT-ALL-OR-NOTHING | `return broker.Registry{}` on `episodeHandler` error | `…PolicySummaryTimeout` (i): R8 `[Workspace.Exec]` |
| MUT-BIND-AILANG-ON-FAIL | on failure bind the eight names to `execFailedHandler{err}` | `…PolicySummaryTimeout` (ii): Dispatch commits `failed`, not R8; `plans == 1` |
| MUT-EARLY-FAULT-EMPTIES | `return broker.Registry{}` when the error is not from `NewAilangToolHandler` (e.g. `!strings.Contains(err.Error(), "policy summary")`) | `…PolicyDirFault`; `…PolicySummaryTimeout` stays green (shows the split) |
| MUT-SHARED-LOCK | `execHandler` takes `w.mu` instead of `execMu` | `…NeverWaitsOnAilangBuild`: `ep1Done` closed when ep2 returns (≈5 s) |
| MUT-DOUBLE-BUILD | delete the `w.execH[episodeID] = …` store | `…ConcurrentRegistryExecRace` pointer equality (deterministic) |
| MUT-EXECH-UNLOCKED | delete `execMu` Lock/Unlock in `execHandler` | `…ConcurrentRegistryExecRace` under `-race -count=5`: `DATA RACE` in any run (PD10) |
| MUT-WRONG-LOCK | lookup under `execMu`, store under `w.mu` | same, `-race -count=5` (PD10) |

### C3 — M2: single-flight build, no lock across the subprocess — AC2.1–AC2.5, AC1.7 registry half

**Files:** `host/daemon/workspace.go`, `host/daemon/workspace_layout.go` (comment `:300` only: a blocking open would hold the episode's in-flight build, not `w.mu`), `host/daemon/workspace_family_test.go`. **`daemon.go` untouched beyond C1's comment.**

**Steps:**
1. `type ailangBuild struct{ root string; done chan struct{}; tool broker.Handler; err error; waiters int }`; field `ailangBuild map[string]*ailangBuild` under `w.mu` (AC2.1).
2. `episodeHandler`: under `w.mu` — success cache hit (`w.handler`, unchanged key and `root == sandbox`) → return; in-flight entry with the same root → `waiters++`, unlock, `<-done`, return its `(tool, err)`; else insert a new entry, unlock. Build **without any lock** (mkdirs, `linkPackageCache`, `checkLockCoverage`, render, `writePolicy`, `NewAilangToolHandler` under `workspaceHandlerBudget`). Re-lock: on success store `w.handler[ep]`; delete the in-flight entry (if it is still this one); set `tool, err`; `close(done)`; unlock. **A failure is never cached.** The `waiters` field is read only by tests (PD2).
3. `registry()` prints one line per caller whose result is an error (unchanged code path), so N failing callers → N lines.

**Tests:**
- **`TestWorkspaceAilangBuildSingleFlight`** (AC2.2). Budget 5 s, hold ep1. Goroutine A `registry("ep1")`; wait `summaries.ep1 == 1`; goroutine B `registry("ep1")`; precondition wait until (`waiters` of `w.ailangBuild["ep1"]` read under `w.mu` ≥ 1) **or** `summaries.ep1 ≥ 2`; then `fail.ep1` created, hold removed; join both. Assert `summaries.ep1 == 1`, both registries `[Workspace.Exec]`, log has exactly 2 lines for ep1. Then a second phase on ep2 (healthy): same shape, both get nine names and the same AILANG handler pointer, `summaries.ep2 == 1`.
- **`TestWorkspaceAilangBuildDoesNotBlockOtherEpisodes`** (AC2.3 + D2's moved AC1.7 half). Budget 5 s, hold ep1; goroutine `registry("ep1")` → `ep1Done`; wait `summaries.ep1 == 1`; `registry("ep2")` returns all nine names and its exec handler runs `where` in ep2; assert `ep1Done` still open; release; join.
- **`TestWorkspaceFailedBuildRetriesEachCall`** (AC2.5). `fail.ep1`; three sequential `registry("ep1")` → each `[Workspace.Exec]`, `summaries.ep1 == 3`, 3 log lines; remove `fail.ep1`; 4th call → nine names, `summaries.ep1 == 4`, still 3 lines.
- Re-run `TestWorkspaceTwoEpisodeConcurrentRegistryExecRace` under `-race` (unchanged).
- **AC2.4 — unedited** (each must `--- PASS`, and `git diff <C2> -- host/daemon/workspace_test.go host/daemon/workspace_layout_test.go` is empty at C3): `workspace_test.go:207, :220, :270, :558`; `workspace_layout_test.go:128, :162–165, :596, :599, :869, :874–878`.

**PASS list (C3):** C2's list + `TestWorkspaceAilangBuildSingleFlight TestWorkspaceAilangBuildDoesNotBlockOtherEpisodes TestWorkspaceFailedBuildRetriesEachCall TestWorkspacePolicyRenderedPerEpisode TestWorkspacePackageCacheLinkedIntoEpisodeHome TestWorkspacePackageCacheDigestIsStamped`.

**MUTs (C3):**

| MUT | Edit | Must red |
|---|---|---|
| MUT-NO-SINGLEFLIGHT | ignore the in-flight entry (always start a new build) | `…SingleFlight`: `summaries.ep1 == 2` |
| MUT-LOCK-ACROSS-SUBPROCESS | hold `w.mu` across the build (C2's `episodeHandler`) | `…DoesNotBlockOtherEpisodes`: `ep1Done` closed when ep2 returns (≈5 s) |
| MUT-NEGATIVE-CACHE | on failure store the error in `w.ailangBuild` and do not delete it; later callers return it | `…FailedBuildRetriesEachCall`: `summaries == 1`, call 4 still `[Workspace.Exec]` |
| MUT-ASYNC-RETRY | on failure, return the failure and start `go w.episodeHandler(…)`; callers return the last recorded failure while a background build runs | `…FailedBuildRetriesEachCall` call 4 stale. If it cannot be expressed in ≤ 10 lines, record "not applied" with the reason |
| MUT-LINE-PER-ATTEMPT | waiters return `(nil, nil)`-shaped "no handler" without the error (registry prints nothing for them) | `…SingleFlight`: 1 line, not 2 |
| MUT-LOCK-ACROSS-SUBPROCESS also re-run against `TestWorkspaceExecNeverWaitsOnAilangBuild` | — | must stay GREEN (exec never takes `w.mu`): shows the two locks are independent |

---

## 6. Landing shape

**One PR, 3 commits (C1, C2, C3), each green on §2.1 at its boundary, one executor run.** M3 (charter row 152 → LANDED, follow-up rows R1/R2 from design §10, the note in `design_docs/verification/world-row140-m4/README.md` that F5's "10 s" was the 3 s budget) is **controller work**; the executor touches no `design_docs/world-mission*.md` and no `design_docs/verification/` file. The PR body carries: the per-boundary gate table (rc + daemon race-leg wall), the mutation ledger, the tests-first proof, and any SURVIVED mutant. The executor never reruns CI; it hands the PR to the controller.

If the run dies after C2, C1+C2 are a complete, landable M1 (exec decoupled; F5's lock coupling remains for the AILANG family only) and C3 can be a second run.

---

## 7. Planner decisions (the executor applies them; any deviation is reported)

| # | Ambiguity | Decision and the measurement behind it |
|---|---|---|
| PD1 | "held on a FIFO" (D8) | Per-episode **hold file** polled by the stub (`sleep 0.02` loop). Same "never answers until released" property; releases any number of spawns; survives a killed reader. The design's probe already used a stub `sleep` killed by the budget (V6: 3.007 s, group kill worked) |
| PD2 | AC2.2 needs "the second caller has joined" deterministically (D8) | `waiters int` on the in-flight entry, incremented under `w.mu`; tests poll it under `w.mu`. The precondition loop also exits on `summaries ≥ 2` so MUT-NO-SINGLEFLIGHT reds fast instead of timing out |
| PD3 | AC1.7's `registry("ep2")` arm is unreachable at M1 (D2) | Moves to C3 `TestWorkspaceAilangBuildDoesNotBlockOtherEpisodes`; C2 keeps the direct `execHandler("ep2")` arm |
| PD4 | `execSpawns` wrong lock after AC1.2 (D3) | C1 switches it to `execMu`; the 3 srt tests run under `-race` locally at every boundary (§2.1) |
| PD5 | MUT-EXEC-USES-SANDBOX's killer (D4) | the lock-coverage and FIFO amendments (module root `tools`) |
| PD6 | `errors.As` on a log line (D5) | `errors.As` on a direct `episodeHandler` error + log `Contains "timed out"` |
| PD7 | publishing two descriptors (D6); two plans | `publishTools` writes one revision with both; `publishReadTool` becomes a wrapper with identical behaviour. Two coordinators over one store/binder, each with its own fake runner (calls are sequential, so per-coordinator in-flight maps never interact). If that fails, a single runner keyed on `e.Source` (distinct sources per tool) is the fallback — report it |
| PD8 | the exec call in amended refusal tests (D7) | reuse the test's single `reg`; never call `registry`/`execJSON` again before the log assertions |
| PD9 | stale `daemon.go:741–743` comment (D9) | comment-only edit in C1 |
| PD10 | race-detector kills are probabilistic (D11) | `-race -count=5`; report SURVIVED honestly; MUT-DOUBLE-BUILD is the deterministic tooth |
| PD11 | order of the two families in `registry` | AILANG first, exec second: today's operator-line order is preserved for any test that configures a failing exec profile |
| PD12 | "budget seam set to the test's own bounded context" | **5 s** for hold tests (only a mutant ever waits it; correct code releases the hold within ms) and **300 ms** for the never-answer tests (AC1.4). Never assert elapsed time. 5 s < `boundedTestContext`'s 30 s, so a mutant reds rather than hanging |
| PD13 | AC1.9 trigger (brief item 3) | regular file at `<stateDir>/policies` (§3 root/euid paragraph); no chmod, no skip |
| PD14 | where new tests live | `workspace_family_test.go`; the `TestWorkspace*` prefix keeps them in the focused leg's `-run '^(TestWorkspace|…)'` |

---

## 8. Risks the executor must watch (report; never paper over)

1. **CI race-leg headroom (row 154).** New tests add < 3 s locally (two 300 ms budgets, ms-scale hold releases). Report `host/daemon` race wall at each boundary.
2. **Goroutine leaks into later tests.** Every test that starts a `registry` goroutine joins it before returning (release the hold first); the budget seam restore runs only after the join.
3. **`TestExecSrtMCP*` CI flake** ("host callback timed out", D-WORLD-68 history). It predates this row; the controller counts CI attempts. Locally a red there at a boundary is this commit's.
4. **Linux.** The hold-file stub runs under dash on CI (`sleep 0.02` is coreutils, fractional OK); darwin `/bin/sh` is bash 3.2. No bash-only syntax in the stub.
5. **`execMu`/`w.mu` ordering.** No code path may hold both; if one ever must, take `w.mu` first and document it. A deadlock shows as a 30 s `boundedTestContext` expiry.

---

## 9. Effort

| Commit | Estimate |
|---|---|
| C0 re-baseline | 0.05 d (~15 min, `verify_go.sh` ~11 min) |
| C1 production + amendments + QUICKSTART | 0.2 d |
| C2 new M1 tests + fixtures + drills | 0.2 d |
| C3 single-flight + tests + drills | 0.2 d |
| **Total** | **~0.65 d**; executor wall-clock ≈ 3–3.5 h (4 full gate runs ≈ 4 × 15 min, ~17 mutation drills ≈ 30 min incl. 5 s mutant waits and `-count=5` race drills) |

## 10. Not verified by the planner

- That MUT-EXECH-UNLOCKED / MUT-WRONG-LOCK red under `-race -count=5` (PD10): the executor's drill is the measurement.
- The hold-file stub under linux dash: CI is the measurement (risk 4).
- The exact operator-line text for the `policies`-file fault beyond `not a directory` (Go's `MkdirAll` wording is `mkdir …: not a directory`; the test matches the substring only).

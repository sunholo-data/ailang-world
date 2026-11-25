# Sprint plan — w-world-address-bar-verbs-v4, sprint S5d (iteration 434, row 243)

M1 ships the S3b reserve/observe step bank and pins the S5c read seam before it acquires an effect caller. M2 ships the signed-start worktree reservation, Git creation and readback, plus attended disposition in the same patch. M2's preparation entry has zero live callers. Its only new CLI entry is `ailang-worldd service runner reconcile-worktree ACTION-ID`, which requires an existing owned worktree reservation.

Authority is design r1, quorum r2 PROCEED 3/3 with all seats present, and attended D-WORLD-119=A / D-WORLD-120=A at `82f05004`. Landed S1 #606, S2 #609, S2b #611, S3 #643, S4 #651, S4b #654, S5 #657, S5b #659 and S5c #661 (`3a6e98a7`; judge PASS 87, 0 blocking). The runner admission, S8 browser, S12 herdr and finish gates remain binding.

## 0. Binding constraints

- Executor lane `codex:gpt-6.2-sol` uses the incremental patches in Appendix P with `/usr/bin/patch -p1 -E`. Controller commits the plan and the verified snapshots. The planner touches only the two plan files in the checkout.
- The cap is 1,400 changed product lines. Reservation, effect, readback, and attended disposition have one commit boundary. M1 independently supplies an inert journal with restart-tested CAS.
- The launch, worktree-add, and reserve/observe path requires an independently verified paired-process activation sprint. In M2, the attended command may read Git and dispose only the reservation it reviewed. It never writes Git, files, credentials, launches, or daemon state.
- No real services, launchctl, credential operations, or test worktrees on the real repository. Tests use fixture stores from the existing migration, the repository signer and disposable Git repositories under `t.TempDir()`. The only real-path evidence allowed is read-only fingerprinting by the controller.

## 1. Pre-flight

### 1.1 Baseline at 68c36735

`git rev-parse HEAD` returned `68c3673579b53326a2c5d4d7bc62b9d3bb7a1105`. `git status --short` was empty before the plan. The first gate sweep started on a pristine `git archive 68c36735` export before staging changes. The raw `gofmt -l cmd host` printed exactly `host/daemon/daemon.go`, and `verify_go.sh` stopped at `tracked_binary_gate.sh: FATAL: git ls-tree HEAD failed: fatal: not a git repository`. The final gate verifies that file byte-for-byte against base and labels the full check SKIPPED for repository-less export modes. It runs `verify_go` on the committed checkout.

The updated §6 gate on the pristine export completed in 120 s with zero informative failures. The final sweep ran from the archive rebuilt from the actual plan fences and completed in 198 s, also with zero informative failures. Both wrote `ALL_DONE.68c36735`. The build, vet, compile-only fence, AILANG verify, named inventories, store package and store race, mint/preset regressions, routes, scope and production caller checks passed. The route census is 54 registrations and three surface POSTs on both exports. The final caller leg found zero callers of the new preparation entry and exactly one call to the attended CLI implementation.

Both sweeps recorded the CLI package, CLI race package and QUICKSTART §9 flag test as **UNINFORMATIVE UNDER SANDBOX** because the existing fixtures bind loopback. File observations were also **UNINFORMATIVE UNDER SANDBOX**. Archive `verify_go` was **SKIPPED**. Run the final mode against the committed checkout outside the sandbox for these legs, plus the Linux CI jobs.

`LEG added-cli-race rc=0 30s` on the final archive. The five added CLI families took 28.851 s of Go package runtime under race. S5c's controller measured 389 s for the full package against 480 s; adding the measured named cost gives an indicative 418 s. The full current CLI race runtime remains **UNINFORMATIVE UNDER SANDBOX**. M1 and M2 passed `go build ./...`, touched-package vet, named acceptance and changed-file gofmt. See `verify-final-M1.log`, `verify-final-M2.log`, `gate-base-final/legs.txt` and `gate-bounded/legs.txt`.

The external gate command is `/bin/bash "$T/gates.sh" "$W" "$S/gate-controller-final" final` (§7). It includes `go test -count=1 -timeout 8m ./cmd/ailang-worldd`, the corresponding `-race` leg, the named `TestQuickstartSection9FlagsMatchTheCLI` check, `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` and read-only before/after file observations. §4.1 supplies the separate whole-package mutant command.

### 1.2 Current inputs and coupling inventory

Use `git grep -n -E 'func (SYMBOLS)\(' -- cmd host | grep -vF '_test.go:'` for unqualified function rows and `git grep -nF 'LITERAL' -- PATH` for method and SQL rows below. Raw output is in `$S/preflight.txt` and `$S/anchors.txt` (§7). Line numbers in this table are the base.

| ID | Command | Observed at 68c36735 |
|---|---|---|
| C1 | `git diff --name-only 3a6e98a7 68c36735 -- cmd host docs app`; `git diff --stat 3a6e98a7 68c36735 -- cmd host docs app` | 11 paths, 689 additions and 1 deletion. `main.go` plus `host/daemon/{daemon.go,surface_onboarding.go,surface_onboarding_test.go,surface_verb.go,surface_verb_test.go}`, `host/surfacepage/{repos.go,repos/index.html,repos/repos.js}` and `host/verbs/{lookup.go,lookup_test.go}`. |
| C2 | `git diff 3a6e98a7 68c36735 -- cmd/ailang-worldd/main.go`; `git grep -nF 'service runner install\|' -- cmd/ailang-worldd/main.go` | Hunk at `runServe`, +`PresetsDir: defaultPresetsDir()` near 651. The only S5d hunk there is the help inventory at 86. |
| C3 | `git grep -n -E 'MintOperatorStart\(\|reconstructOperatorStart\(\|authority.RecheckVerb\(\|ClaimOperatorAction\(\|ReserveOperatorStep\(\|ObserveOperatorStep\(' -- '*.go' \| grep -vF '_test.go:'` | `host/authority/mint.go:204`, `cmd/ailang-worldd/operator_start_input.go:66` and `host/store/operator_cas.go:119` are definitions. Live positive control `cmd/ailang-worldd/episode_finish_surface.go:239`. `test ! -e host/store/operator_steps.go` succeeded at base. |
| C4 | `git grep -n -E 'operatorStartFacts\(\|resolvePresetStart\(\|loadLaunchPreset\(\|presetLayout\(\|canonicalLaunchPath\(' -- cmd/ailang-worldd` | `operator_start_input.go:14,19,66,68`, `episode_preset.go:79,80,328`, `preset.go:154,329,345`, `episode_launch.go:73`. Explicit episode selects `episodes[0]` (`episode_preset.go:100`); `launchLauncher` use is conditional on interactive mode (146). `prepareMint` rejects reused layouts via `episode_preset.go:232`. |
| C5 | `git grep -n -E 'func (provisionEpisode\|Add\|AddExisting\|Remove\|Output)\(\|worktree\.(Add\|AddExisting\|Remove)\(\|var sessionGit' -- cmd/ailang-worldd/session_new.go host/worktree/worktree.go` | `session_new.go:169,405,510,512,577`; worktree `Add:63`, `AddExisting:77`, `Output:83`, `Remove:110`. M2 only calls `Add` and bounded readback; a derived pre-existing branch is a conflict. |
| C6 | `git grep -n -E 'OperatorCAS\|operatorCASRow\(\|ClaimOperatorAction\(\|operatorChangedOne\(' -- host/store/operator_cas.go`; `git grep -n -E 'OperatorActionRow\|operatorAccess\(\|operatorColumns' -- host/store/operator_actions.go` | `operator_cas.go:15,32,119,157`; `operator_actions.go:25,40,50,60`. `operator_schema.sql:1,2,34,48` identifies schema 10, action, heartbeat and step rows. The step PK is `(kind, identity)` (`:58`); reserved rows cannot carry receipt bytes. |
| C7 | `git grep -n -E 'readOperatorMint\(\|decodeOperatorMintAction\(\|operatorMintApprover\(\|NOT IN.*worktree\|type OperatorMintAttempt' -- host/store/operator_mint.go`; `git grep -nF 'func RecheckClaimedOperatorAction' -- host/authority` | Mint attempt `:22`, enrollment recheck `:69`, commit `:94`, decode/digest `:175`, `readOperatorMint:249`. The step exclusion at `:266` permits `operation,worktree,mint`. Authority recheck `host/authority/judge_operator.go:129`. |
| C8 | `git grep -n -E 'reconcile-mint\|reconcileRunnerMint\(\|readRunnerMint\(\|inspectRunnerMintToken\(\|installedRunner\(\|runnerWrite\(\|runnerDB\(\|operatorEpisodeLockPath\(\|openRunnerLock\(' -- cmd/ailang-worldd` | Mint resolver `operator_mint.go:17,38,41,73,92`, runner service help/grammar/dispatch `operator_service.go:24,222,258,275,286`, top service help/confirm/dispatch `service.go:37,171,288`. `operator_episode_lock.go:14`, `operator_runner.go:25,35,162,180`. `credentialsFor` `session.go:121` and `CanonicalPath` `host/store/credentials.go:185`. |
| C9 | `git grep -n -E 'syscall\.Flock.*LOCK_NB' -- cmd host \| grep -vF '_test.go:'` | Seven lines including one comment. Four CLI sites (`operator_runner.go:35`, `operator_service.go:251`, `service.go:158,662`) and two writer lock calls (`writer_lock_unix.go:41,86`); existing gate count is unchanged. |
| C10 | `git grep -n -E 'Test(DaemonDependencyAllowlist\|WorldBoundaryDependencyAllowlist\|ProductionGoSurface\|ProductionContextRoots\|CLICredentialWriterSites\|EveryVerbHasHelp\|QuickstartServiceSectionMatchesTheCLIHelp\|VerbTableMatchesTheCLISwitch\|CredentialsMethodsRequireADeadline\|StoreSurfaceTakesContext\|PluginCommandsMatchTheCLIHelp)\(' -- cmd host` | Daemon deps `daemon_test.go:899`; boundary deps `allowlist_world_test.go:758`. Go files and context roots `context_roots_test.go:234,129`; credential writer inventory `cli_credential_writers_test.go:18`; explicit deadline inventory `credentials_test.go:277` and reflection guard `deadline_guard_test.go:43`. Help `cmd/ailang-worldd/help_test.go:21`, service runbook `quickstart_service_test.go:15`, bar parity `verbs_parity_test.go:66`, plugin `plugin_help_test.go:52` / `cmd/world-publish/plugin_help_test.go:58`. |
| C11 | `git ls-tree --name-only 68c36735 -- cmd host docs app`; `git log --oneline --all --grep '#604' -2` | `cmd`, `docs` and `host`; `app` absent. #604 merge `65108353`. |
| C12 | `wc -l ~/.ailang/state/world-iter433/plan-scratch/st-s5d-footprint/cmd/ailang-worldd/s5c_{worktree_probe.go,files_probe.go,fd_probe_test.go}` | 59 worktree, 53 file-publication and 72 process-harness probe lines. `wc -l` plus `diff -u` on the previous bank and test reports the old 132-line measurement. |

Inventory impact from the applied diff and §6 gates:

| Inventory | S5d disposition |
|---|---|
| `TestEveryVerbHasHelp`, runner/service help, top usage, QUICKSTART and `TestQuickstartServiceSectionMatchesTheCLIHelp` | M2 extends the enumerated runner command and usage strings. There is no separate help golden for this subcommand. The worktree window test reads the QUICKSTART example, compares it with the emitted remedy, and executes that exact CLI command through `runServiceWith`. Prior `reconcile-mint` also appears in top-level help, aligning it with the already-present nested help. |
| `TestVerbTableMatchesTheCLISwitch`, both `TestPluginCommandsMatchTheCLIHelp` instances and `TestOperatorNeverDispatch` | Existing service case and plugin inventories run unchanged. The service subcommand stays in the attended service router; the bar gains no verb. |
| `TestCLICredentialWriterSites`, `TestDaemonCredentialWriterSites` | New store writes use the existing `runnerWrite` opener; no inventory change. Fixture supervisor counters are checked alongside store state. |
| `TestWorldBoundaryDependencyAllowlist`, `TestDaemonDependencyAllowlist`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestBoundaryASTWriteGuard`, `TestBrokerDependencyAllowlist` | Existing worktree, store, approval and authority imports in the CLI group; no new package or dependency exception. These run by name in §6. |
| `TestProductionGoSurface`, `TestProductionContextRoots`, `TestStoreSurfaceTakesContext`, `TestCredentialsMethodsRequireADeadline` | New files are within `cmd/host`; deadlines come from caller context or `budgetContext`. Only the credentials method test is edited (three added calls in each milestone). Schema 10 DDL stays byte-identical. |
| `TestWriterLockCrossProcessContention`, `TestWriterLockHeldIsNonMutating`, `LOCK_NB` census | The two new callers reuse `openRunnerLock` with flags `0`, never creating an episode lock on reservation or recovery. The existing helper owns the flags. No writer-lock inventory edit. The contention fixtures use a nonblocking contender or a three-second release fuse. |

### 1.3 Decisions

1. Put durable state and SQL CAS in `host/store`. Keep Git and attended CLI effects at the host boundary. These operations require credentials-store transactions and OS/Git identity, as described in design §6.
2. Persist the worktree target before Git, including action ID, physical repo, common directory, physical path, new branch and starting SHA. This is private journal evidence. The second execution refuses under the record/CAS and Git collision laws.
3. A direct success needs one uninterrupted exclusion scope, successful `worktree.Add`, and exact registration, path, common dir, ref and actual checkout HEAD. A nonzero Git return plus complete absence stops the action. A partial or uncertain effect retains its owner and reserved step. A crash can leave `running/reserved`; the attended verb accepts that row as well as `uncertain/reserved`.
4. The attended choice is stop-and-retain. When any worktree evidence remains, the operator explicitly stops the action and preserves every Git artifact. Exact-looking evidence after a crash is eligible for that stop, and requires a new reviewed path if any retained artifact conflicts. The command does not convert an interrupted effect into a successful start.

### 1.4 Premises found false or incomplete at HEAD

| ID | Command and observed evidence | Resolution |
|---|---|---|
| F1 | `git grep -nF 'launchLauncher(agent, input.hostsDir' -- cmd/ailang-worldd/episode_preset.go` → 146, under `mode == "interactive"` at 145. `git grep -nF 'headless: true' -- cmd/ailang-worldd/operator_start_input.go` → 20. | Workspace is a live layout guard in headless mode. Hosts is deliberately forwarded to the seam, but the headless path never consults it. The forwarding test injects the same read function and checks the actual seam arguments, including its real workspace-layout refusal. |
| F2 | `git archive 68c36735 \| tar -x -C "$S/base"` then `gofmt -l cmd host` → `host/daemon/daemon.go`; name/touched diff at C1. | Raw formatting baseline is red. The final script strictly checks that the pre-existing file is byte-identical to base and all changed files are gofmt clean. |
| F3 | Base `./scripts/verify_go.sh`, with pinned environment, printed `test_tracked_binary_gate: 24 passed, 0 failed` then `git ls-tree HEAD failed: fatal: not a git repository`. | Only committed-checkout full verify is informative; export modes mark the leg SKIPPED. |
| F4 | `git grep -nF "kind NOT IN ('operation','worktree','mint')" -- host/store/operator_mint.go` → 266. | Mint recovery's exclusion accepts the worktree kind without filtering its phase. The future consumer must require an observed, successful worktree receipt before mint. S5d has no path to mint. |
| F5 | `git grep -n -F 'runner install\|start\|stop\|status\|cancel' -- cmd/ailang-worldd/{main,service}.go` → `main.go:86` omitted `reconcile-mint`; `service.go:37` included it, `operator_service.go:24` describes it. | Add both attended dispositions to top usage as part of the same help line. |
| F6 | `git grep -nF 'active_child' -- cmd host` paired with `git grep -nF 'func TestOperatorServiceUpgrade' -- cmd/ailang-worldd` → zero child arm, positive upgrade and pair-fence tests. | Current upgrade tests exercise schema and rollback fences. The active-child case is an activation gate for S5f/S6. |
| F7 | `git ls-tree --name-only 68c36735 -- cmd host docs app` → no `app`; `git ... --grep '#604'` → `65108353`. | Snapshot and round trip handle absent `app`. The mobile work merged. S8 still owns explicit app-shell and Chrome sequencing. |

### 1.5 Size and CI race cost

Cap measurement after each controller commit is exactly `git diff -M --numstat 68c36735 HEAD -- cmd host docs app`. Sum adds and deletes and repeat without the pathspec for the product-to-base comparison (the separate controller plan commit contributes only the two plan files to a literal unfiltered history). The planner cannot make Git writes, so `python3 "$S/stage-diff.py" M1 M2` uses `git diff --no-index -M --numstat` against the same pristine export, also checking every unfiltered path. No rename or binary changes were measured.

| Milestone | Added | Deleted | Cumulative | Incremental patch | Unfiltered product | STATUS/DELETED/MANIFEST |
|---|---:|---:|---:|---:|---:|---|
| M1 | 274 | 6 | 280 | 280 | 280 | 5/0/5 |
| M2 | 1384 | 13 | 1397 | 1121 | 1397 | 14/0/14 |

````text
M1 increment 280
5	3	cmd/ailang-worldd/operator_start_input.go
32	3	cmd/ailang-worldd/operator_start_input_test.go
3	0	host/store/credentials_test.go
162	0	host/store/operator_cas_test.go
72	0	host/store/operator_steps.go
cumulative 274 6
unfiltered 280
M2 increment 1121
2	2	cmd/ailang-worldd/help_test.go
1	1	cmd/ailang-worldd/main.go
6	3	cmd/ailang-worldd/operator_service.go
5	3	cmd/ailang-worldd/operator_start_input.go
32	3	cmd/ailang-worldd/operator_start_input_test.go
261	0	cmd/ailang-worldd/operator_worktree.go
110	0	cmd/ailang-worldd/operator_worktree_reconcile.go
557	0	cmd/ailang-worldd/operator_worktree_test.go
1	1	cmd/ailang-worldd/service.go
11	0	docs/QUICKSTART.md
6	0	host/store/credentials_test.go
162	0	host/store/operator_cas_test.go
84	0	host/store/operator_steps.go
146	0	host/store/operator_worktree.go
cumulative 1384 13
unfiltered 1397
````

Bank + forwarding is independently green. At the first production footprint probe, M1 was 175 lines and cumulative M2 was 728 before the full tests. The compiled bill includes both matrices and the runbook. M2 exceeds the preferred ~700 increment because the reservation, add/readback, uncertainty, and attended disposition cross one durability boundary. The measured total fits the cap with the full gate, so Appendix C begins at S5e.

A literal unfiltered branch diff including the controller's two plan records adds 2623 lines to the product sum. Expected unfiltered branch totals are 2903 after M1 and 4020 after M2.

### 1.6 Portability

Darwin arm64, `git --version` 2.54.0; Apple `/usr/bin/git --version` 2.50.1. Ubuntu 24.04 is the Linux verification host. The helper parses `--porcelain -z`, evaluates symlinks through `canonicalLaunchPath` on registration paths and the reserved ancestor, and compares physical checkout and common paths. This normalizes `/var` / `/private/var` on Darwin. Tests assert file-type refusal or an unchanged snapshot, never a Darwin-specific UID or exact symlink errno. The distinct-open-file-description nonblocking flock contract is already tested in `TestWriterLockCrossProcessContention` and `TestOperatorRunnerKernelLock`. The new caller uses the same sidecar and does not add flags or depend on process-inherited flock semantics. Run Linux verification through the existing CI job; no source or test edits in `host/daemon` and `host/verifygate`.

## 2. How the executor applies a sub-milestone

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter434-row243-v4s5d
T=$HOME/.ailang/state/world-iter434/executor
S=$HOME/.ailang/state/world-iter434/plan-scratch
B=68c3673579b53326a2c5d4d7bc62b9d3bb7a1105
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5d-sprint-plan.md
mkdir -p "$T"
for name in accept.sh gates.sh snapshot.sh; do
  awk -v n="$name" '$0 == "````bash " n {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$name"
  test -s "$T/$name" && /bin/bash -n "$T/$name" || exit 2
done
for M in M1 M2; do
  awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
  test -s "$T/$M.patch" || exit 2
  (cd "$W" && /usr/bin/patch -p1 -E < "$T/$M.patch") || exit 2
  /bin/bash "$T/accept.sh" "$W" "$M" fast || exit 2
  /bin/bash "$T/snapshot.sh" "$W" "$M" || exit 2
done
````

Each diff is incremental. Apply M1 to the base and M2 to the M1 tree. If the controller has already accepted M1 and applied its commit, use its byte-identical M1 files as the starting point for M2. Exclude the two plan records from snapshots.

## 3. Milestones

### M1 — Step bank and pre-caller forwarding

- `ReserveOperatorStep` compares action ID/digest/generation/owner/state and clock range, input shape, identity uniqueness and the per-action step cap. `ObserveOperatorStep` compares exact action and reserved step independently, checks receipt size, then increments the version exactly once. Observation remains possible after creation authority expires. `OperatorStep` can read the same reservation from a reopened credentials store.
- `TestOperatorStepRestart` checks expiry, a close/reopen, duplicate reserve and observe, stale version, and the exact stored evidence. `TestOperatorStepGuards` isolates wrong keys, receipt size, action state, owner and bound. `TestCredentialsMethodsRequireADeadline` exercises every added method without a deadline.
- `TestOperatorStartForwarding` captures the hosts parameter at the real read boundary and refuses a preset physically within the supplied workspace. It also differentiates the configured physical worktree from the caller's layout root. `TestOperatorStartFacts` uses a full nonzero attempt and pins the injected clock on the returned start record.

### M2 — Reservation, Git readback and unattended proof / attended disposition

- `reserveOperatorStartWorktree` is dormant. It canonicalizes the served credentials path, validates the episode, opens the existing mode-0600 nonblocking sidecar, and holds that fd across reconstruction, the reservation, worktree add, readback and journal writeback.
- The fresh reconstruction supplies the exact S5c reviewed facts. The additional absence check uses independent Git registration/ref and filesystem evidence. The store rechecks the sealed runner proof with the current enrollment and CASs the claimed row, current runner generation/owner/protocol/readiness/catalogue and heartbeat, and absence of start, credential and prior effect rows. It atomically changes the action to `running` and reserves the identity with version 1. The credentials writer is closed before Git runs.
- The Git reader checks the pinned common directory from the repo, NUL registration records and physical path, ref, declared HEAD, .git link shape, actual common directory and independent checkout HEAD. A successful uninterrupted attempt observes the step with `created` and version 2. A Git failure with full absence stops it; otherwise uncertainty keeps the owner and advances the reserved step and action generations. A journal failure after add leaves the earlier durable reservation.
- The CLI validates the approved runner installation and a real reserved start row before acquiring the existing episode lock. It reads owner/identity again under that lock, takes the Git snapshot, presents the actual disposition and requires a typed `y`. After the prompt it repeats Git readback and the write transaction independently CASs the exact current action and step. Later/duplicate effects prevent owner release. Keeping or stopping changes only action/step rows; no Git writes or fake supervisor calls occur. No transition out of a reservation is inferred from the prompt alone.

### 3.1 Branch traces

| Boundary | Named test / observed state |
|---|---|
| Before reservation | `TestOperatorWorktreeWindow/before_reservation`. Claimed row and owner/phase unchanged, step absent, Git branch/worktree absent, files/start/session absent. Exact command refuses on that row. |
| After reservation, before invoking Git | `TestOperatorWorktreeWindow/before_add`. Running row, version-1 reserved step with exact target; Git/branch/files absent. The emitted command, no TTY and negative prompt preserve the snapshot; typed y commits the absent disposition. |
| After Git returns, before readback | `TestOperatorWorktreeWindow/before_readback`. Same reservation with one physical checkout, Git registration and exact branch/SHA. An attended stop retains them. |
| After readback, before journal observation | `TestOperatorWorktreeWindow/before_observe`. Same durable reservation and exact Git evidence. `/observe_error` also injects a database write failure and reads the preserved reservation before retrying the attended remedy. |
| Git nonzero and proven absence | `TestOperatorWorktreeWindow/git_refused`, an actual Git ref-lock conflict, yields observed stopped-absent and no branch/path/registration or mint. |
| Partial or unattributable | `TestOperatorWorktreeWindow/path_only`, `/ref_changed`, `/head_changed`, `/readback_error`, `/add_error_after_checkout` and `/redirected_checkout`. Uncertain owner and reserved identity survive. Typed-y disposition records a stop and retains the real ref/checkout/files. |
| Existing conflicts and live authority drift | `TestOperatorStartCollisions` exercises branch, detached prunable registration, directory, symlink, token/task/MCP, persisted start, credential, stale signed fact, replaced runner or catalogue, heartbeat boundary/protocol/readiness, revocation, fresh transaction expiry, and held or missing episode lock. Compare store bytes, Git refs/worktrees/files and fake calls. |
| Exact once under lock | `TestOperatorStartOnce` reads committed reservation and correct owner while Git evidence is still absent, and verifies contention at that same moment. It then observes exactly one worktree/ref, the pinned branch/SHA, one receipt and no mint/files/launch. Replay leaves the snapshot byte-identical. |
| Concurrent evidence or lock changes during recovery | `TestOperatorWorktreeDispositionCAS` injects stale action owner/digest/generation, step identity/owner/version/phase, a later or duplicate effect, and a ref appearing around the prompt. SQL triggers independently replace action fields between the step and action writes. It verifies the independent refusal and makes no further change. It also covers negative approval and contention before the prompt. |

### 3.2 Design §5.3 row assignment and §5.5 activation

| Reviewed row | Assignment |
|---|---|
| `TestOperatorStartOnce` | M2 worktree half only. Mint receipt already banked in S5b. One file set is S5e; receipt-before-spawn, confinement and one run are S6; stable browser link is S8. |
| `TestOperatorStartChanged` | S5c comparison remains authoritative; M1 seam pins, M2 invokes the reconstruction while holding the reservation lock and tests a stale effect refusal. |
| `TestOperatorStartCollisions` | M2 with real Git and store plus reused token/task/MCP absence checks. No file publication. |
| `TestOperatorWorktreeWindow` | M2 fully. The four injected boundaries, exact Git evidence and verbatim attended disposition live in the same patch. |
| `TestOperatorMintBinding`, `TestOperatorMintRollback`, `TestOperatorMintReceiptWindow` | S5b landed and smoke-tested by related existing mint-recovery gates. A live mint caller remains closed. |
| `TestOperatorFilesReconcile` | S5e with its file-specific disposition in the same patch. |
| `TestOperatorLaunchSession`, `TestOperatorSpawnRefused`, `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchExit` | S5f/S6 shared owner, readiness and drain gate, then S6 launch implementation. |
| `TestOperatorStartBrowser` | S8 bar/reconnect and fake launcher with real Chrome coverage. |
| `TestOperatorHerdrUnavailable`, `TestOperatorHerdrHandoff`, `TestOperatorHerdrWindow` | S12 with E4 and typed ticket/TTY/identity gates. |
| `TestOperatorServiceFence`, `TestOperatorServiceOwnership`, `TestOperatorServiceConfig` | §6 runs all three existing suites. M2 reaches only the attended subcommand; installed configuration and terminal fences precede disposition. |
| `TestOperatorServicePair`, `TestOperatorServiceStartup` | Existing configuration/diagnostic gates run in §6. The paired restart, schema/build compatibility, unavailable daemon, failed boot and startup timeout must all gate any future dispatch. |
| `TestOperatorServiceUpgrade` | Existing rollback and schema guard plus pair-fence in §6. `/active_child` and controller upgrade/drain handoff stay S5f/S6. |
| `TestOperatorServiceWriteScope`, `TestOperatorServiceEnvironment` | Existing fixture lifecycle and environment boundaries run by name. S5d adds no lifecycle command, process or ambient browser input. |
| `TestOperatorFinishReadiness` | Existing store test in §6. Keep §2.2 E2-equivalent finish preflight at activation and future finish claims. |
| `TestOperatorSchema10` | Existing fresh/v9/failure compatibility and DDL gate (§6), no schema change. |

## 4. Mutation table

Mutation specs, exact original/replacement spans and package holding each killing test are in `$S/mutants.json`; replay results and the exact commands are in `$S/mutation-ledger.jsonl`, linked by ID. A named §5.3 survivor blocks that row until the test or implementation is fixed. Compilation errors, static checks, and sandbox bind failures are never classified as behavioural kills.

**220/220 claimed named mutants killed** (36 M1, 184 M2). The 2 inherited equivalent controls survived their named runs and are listed separately. All 32 informative store whole-package mutants died. The other 188 positive CLI whole-package rows reach loopback fixtures and await §4.1. No claimed design-row survivor remains.

M1 non-vacuity is shown by the workspace, hosts and phase guards (N-WORKSPACEROOT-EMPTY, N-HOSTSDIR-EMPTY, O-PHASE). M2 non-vacuity is shown by a real relaxed absence, effect before checkpoint, path-only success, dropped action predicates and bypassed attended prompt (W-INITIAL-ABSENCE, W-RESERVE-ORDER, W-DESIGN-PATH-ONLY, D-STORE-OWNER/GEN/DIGEST and D-NO-PROMPT). The initial survivors O-PHASE, W-INITIAL-READERR, GIT-REG-EXACT, GIT-REG-BRANCH-LOOKUP, GIT-REF-ABSENCE, GIT-CHECKOUT-REF and the three action CAS rows were retargeted to independent fixtures and re-run killed. Rows that invert a positive admission predicate use the existing success path and verify its durable effects. The ledger points to the exact edit and command. No named kill depends on an AST, inventory, or source-string check. Two initial propagation mutations left the errors import unused. The final replay inverted the original errors.Is checks and compiled.

| ID | Production file and line | Test package | Named behavioural test | Named | Whole package |
|---|---|---|---|---|---|
| N-WORKSPACEROOT-EMPTY | `cmd/ailang-worldd/operator_start_input.go:22` | `cmd/ailang-worldd` | `TestOperatorStartForwarding` | KILLED | §4.1 |
| N-HOSTSDIR-EMPTY | `cmd/ailang-worldd/operator_start_input.go:22` | `cmd/ailang-worldd` | `TestOperatorStartForwarding` | KILLED | §4.1 |
| N-WORKSPACE-SWAP | `cmd/ailang-worldd/operator_start_input.go:22` | `cmd/ailang-worldd` | `TestOperatorStartForwarding` | KILLED | §4.1 |
| N-ATTEMPT-ZERO | `cmd/ailang-worldd/operator_start_input.go:42` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| N-CLOCK-CONST | `cmd/ailang-worldd/operator_start_input.go:24` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | KILLED | §4.1 |
| N-WORKTREE-DERIVED | `cmd/ailang-worldd/operator_start_input.go:60` | `cmd/ailang-worldd` | `TestOperatorStartForwarding` | KILLED | §4.1 |
| N-STORE-RAW-CRED | `cmd/ailang-worldd/operator_start_input.go:32` | `cmd/ailang-worldd` | `TestOperatorStartFacts` | SURVIVED (equivalent) | Equivalent control |
| N-RES-ABS | `cmd/ailang-worldd/operator_mint.go:95` | `cmd/ailang-worldd` | `TestOperatorMintRecovery` | SURVIVED (equivalent) | Equivalent control |
| R-ACCESS | `host/store/operator_steps.go:18` | `host/store` | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED |
| O-ACCESS | `host/store/operator_steps.go:38` | `host/store` | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED |
| G-ACCESS | `host/store/operator_steps.go:61` | `host/store` | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED |
| R-OP | `host/store/operator_steps.go:21` | `host/store` | `TestOperatorStepGuards/reserve_operation` | KILLED | KILLED |
| R-EMPTY-ID | `host/store/operator_steps.go:21` | `host/store` | `TestOperatorStepGuards/reserve_empty_identity` | KILLED | KILLED |
| R-EMPTY-OWNER | `host/store/operator_steps.go:21` | `host/store` | `TestOperatorStepGuards/reserve_empty_owner` | KILLED | KILLED |
| R-VERSION | `host/store/operator_steps.go:21` | `host/store` | `TestOperatorStepGuards/reserve_version` | KILLED | KILLED |
| R-CAS-ID | `host/store/operator_steps.go:27` | `host/store` | `TestOperatorStepGuards/reserve_id` | KILLED | KILLED |
| R-CAS-DIGEST | `host/store/operator_steps.go:27` | `host/store` | `TestOperatorStepGuards/reserve_digest` | KILLED | KILLED |
| R-CAS-GEN | `host/store/operator_steps.go:27` | `host/store` | `TestOperatorStepGuards/reserve_generation` | KILLED | KILLED |
| R-CAS-OWNER | `host/store/operator_steps.go:27` | `host/store` | `TestOperatorStepGuards/reserve_owner` | KILLED | KILLED |
| R-STATE | `host/store/operator_steps.go:28` | `host/store` | `TestOperatorStepGuards` | KILLED | KILLED |
| R-EXPIRY | `host/store/operator_steps.go:28` | `host/store` | `TestOperatorStepRestart` | KILLED | KILLED |
| R-STARTTIME | `host/store/operator_steps.go:28` | `host/store` | `TestOperatorStepGuards` | KILLED | KILLED |
| R-CAP | `host/store/operator_steps.go:29` | `host/store` | `TestOperatorStepGuards` | KILLED | KILLED |
| R-REPLAY | `host/store/operator_steps.go:30` | `host/store` | `TestOperatorStepRestart` | KILLED | KILLED |
| O-EMPTY | `host/store/operator_steps.go:41` | `host/store` | `TestOperatorStepGuards` | KILLED | KILLED |
| O-CAP | `host/store/operator_steps.go:41` | `host/store` | `TestOperatorStepGuards` | KILLED | KILLED |
| O-PHASE | `host/store/operator_steps.go:46` | `host/store` | `TestOperatorStepRestart` | KILLED | KILLED |
| O-CAS-KIND | `host/store/operator_steps.go:46` | `host/store` | `TestOperatorStepGuards/observe_operation` | KILLED | KILLED |
| O-CAS-IDENTITY | `host/store/operator_steps.go:46` | `host/store` | `TestOperatorStepGuards/observe_empty_identity` | KILLED | KILLED |
| O-CAS-VERSION | `host/store/operator_steps.go:46` | `host/store` | `TestOperatorStepGuards/observe_version` | KILLED | KILLED |
| O-CAS-ACTION-GEN | `host/store/operator_steps.go:47` | `host/store` | `TestOperatorStepGuards/observe_generation` | KILLED | KILLED |
| O-CAS-ACTION-DIGEST | `host/store/operator_steps.go:47` | `host/store` | `TestOperatorStepGuards/observe_digest` | KILLED | KILLED |
| O-CAS-ACTION-OWNER | `host/store/operator_steps.go:47` | `host/store` | `TestOperatorStepGuards` | KILLED | KILLED |
| O-VERSION-INC | `host/store/operator_steps.go:45` | `host/store` | `TestOperatorStepRestart` | KILLED | KILLED |
| O-RECEIPT | `host/store/operator_steps.go:45` | `host/store` | `TestOperatorStepRestart` | KILLED | KILLED |
| W-LOCK-CREATE | `cmd/ailang-worldd/operator_worktree.go:40` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/missing_lock` | KILLED | §4.1 |
| W-LOCK-RELEASE | `cmd/ailang-worldd/operator_worktree.go:44` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-RECONSTRUCT | `cmd/ailang-worldd/operator_worktree.go:45` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/stale` | KILLED | §4.1 |
| W-INITIAL-ABSENCE | `cmd/ailang-worldd/operator_worktree.go:61` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/registration` | KILLED | §4.1 |
| W-INITIAL-READERR | `cmd/ailang-worldd/operator_worktree.go:57` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/git_initial_read` | KILLED | §4.1 |
| W-FAULT-RESERVE | `cmd/ailang-worldd/operator_worktree.go:64` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_reservation` | KILLED | §4.1 |
| W-FAULT-ADD | `cmd/ailang-worldd/operator_worktree.go:87` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| W-FAULT-READBACK | `cmd/ailang-worldd/operator_worktree.go:91` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_readback` | KILLED | §4.1 |
| W-FAULT-OBSERVE | `cmd/ailang-worldd/operator_worktree.go:95` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_observe` | KILLED | §4.1 |
| W-RESERVE-ORDER | `cmd/ailang-worldd/operator_worktree.go:87` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-ADD-BRANCH | `cmd/ailang-worldd/operator_worktree.go:90` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-ADD-SHA | `cmd/ailang-worldd/operator_worktree.go:90` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-DESIGN-PATH-ONLY | `cmd/ailang-worldd/operator_worktree.go:99` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/path_only` | KILLED | §4.1 |
| W-ADD-ACK | `cmd/ailang-worldd/operator_worktree.go:99` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/add_error_after_checkout` | KILLED | §4.1 |
| W-ABSENT-GIT | `cmd/ailang-worldd/operator_worktree.go:101` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/git_refused` | KILLED | §4.1 |
| W-OBSERVE-RECEIPT | `cmd/ailang-worldd/operator_worktree.go:106` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-OBSERVE-ERR | `cmd/ailang-worldd/operator_worktree.go:110` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/observe_error` | KILLED | §4.1 |
| GIT-PORCELAIN-Z | `cmd/ailang-worldd/operator_worktree.go:188` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-CHECKOUT-HEAD | `cmd/ailang-worldd/operator_worktree.go:257` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/redirected_checkout` | KILLED | §4.1 |
| GIT-REG-EXACT | `cmd/ailang-worldd/operator_worktree.go:235` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/registration_locked` | KILLED | §4.1 |
| GIT-REG-BRANCH-LOOKUP | `cmd/ailang-worldd/operator_worktree.go:205` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/registration_foreign` | KILLED | §4.1 |
| GIT-REG-ABSENCE | `cmd/ailang-worldd/operator_worktree.go:224` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/registration` | KILLED | §4.1 |
| GIT-REF-ABSENCE | `cmd/ailang-worldd/operator_worktree.go:224` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/ref_descendant` | KILLED | §4.1 |
| GIT-CHECKOUT-REF | `cmd/ailang-worldd/operator_worktree.go:235` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/ref_read_skew` | KILLED | §4.1 |
| D-NO-PROMPT | `cmd/ailang-worldd/operator_worktree_reconcile.go:87` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-GIT-REREAD | `cmd/ailang-worldd/operator_worktree_reconcile.go:97` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/git_changed` | KILLED | §4.1 |
| D-ABSENT-DECISION | `cmd/ailang-worldd/operator_worktree_reconcile.go:84` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-EXACT-REMEDY | `cmd/ailang-worldd/operator_worktree_reconcile.go:16` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| W-STORE-EXPIRY | `host/store/operator_worktree.go:28` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/expires_at_reservation` | KILLED | §4.1 |
| W-STORE-SESSIONS | `host/store/operator_worktree.go:57` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/credential` | KILLED | §4.1 |
| W-STORE-START | `host/store/operator_worktree.go:56` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/episode` | KILLED | §4.1 |
| W-STORE-RUNNER-CATALOGUE | `host/store/operator_worktree.go:53` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/catalogue` | KILLED | §4.1 |
| W-STORE-RUNNER-GEN | `host/store/operator_worktree.go:53` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/runner` | KILLED | §4.1 |
| W-HEARTBEAT-STALE | `host/store/operator_worktree.go:54` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/heartbeat_stale` | KILLED | §4.1 |
| W-HEARTBEAT-FUTURE | `host/store/operator_worktree.go:54` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/heartbeat_future` | KILLED | §4.1 |
| W-HEARTBEAT-PROTO | `host/store/operator_worktree.go:54` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/heartbeat_protocol` | KILLED | §4.1 |
| W-HEARTBEAT-READY | `host/store/operator_worktree.go:54` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/heartbeat_readiness` | KILLED | §4.1 |
| W-STORE-NO-RESERVE | `host/store/operator_worktree.go:65` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-STORE-STATE | `host/store/operator_worktree.go:51` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-STORE-ACTION-GEN | `host/store/operator_worktree.go:64` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| D-STORE-LATER | `host/store/operator_worktree.go:138` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/later_effect` | KILLED | §4.1 |
| D-STORE-DUPLICATE | `host/store/operator_worktree.go:139` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/duplicate_step` | KILLED | §4.1 |
| D-STORE-OWNER | `host/store/operator_worktree.go:137` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/action_owner` | KILLED | §4.1 |
| D-STORE-GEN | `host/store/operator_worktree.go:137` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/action_generation` | KILLED | §4.1 |
| D-STORE-DIGEST | `host/store/operator_worktree.go:137` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/action_digest` | KILLED | §4.1 |
| D-STORE-KIND | `host/store/operator_worktree.go:137` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/action_kind` | KILLED | §4.1 |
| D-STORE-ACTION-STATE | `host/store/operator_worktree.go:137` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/action_state` | KILLED | §4.1 |
| GIT-ANCESTOR-REF | `cmd/ailang-worldd/operator_worktree.go:218` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/ref_ancestor` | KILLED | §4.1 |
| GIT-DESCENDANT-REF | `cmd/ailang-worldd/operator_worktree.go:218` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/ref_descendant` | KILLED | §4.1 |
| D-STORE-FAIL-STATE | `host/store/operator_worktree.go:135` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-STEP-OWNER | `host/store/operator_steps.go:57` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/step_owner` | KILLED | §4.1 |
| D-STEP-VERSION | `host/store/operator_steps.go:58` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/step_version` | KILLED | §4.1 |
| D-STEP-IDENTITY | `host/store/operator_steps.go:58` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/step_identity` | KILLED | §4.1 |
| D-STEP-PHASE | `host/store/operator_steps.go:58` | `cmd/ailang-worldd` | `TestOperatorWorktreeDispositionCAS/step_phase` | KILLED | §4.1 |
| W-LOCK-NB | `cmd/ailang-worldd/operator_runner.go:35` | `cmd/ailang-worldd` | `TestOperatorRunnerKernelLock` | KILLED | §4.1 |
| W-LOCK-EX | `cmd/ailang-worldd/operator_runner.go:35` | `cmd/ailang-worldd` | `TestOperatorRunnerKernelLock` | KILLED | §4.1 |
| O-CAS-STEP-OWNER | `host/store/operator_steps.go:46` | `host/store` | `TestOperatorStepGuards` | KILLED | KILLED |
| R-COUNT-OPERATION | `host/store/operator_steps.go:29` | `host/store` | `TestOperatorStepGuards` | KILLED | KILLED |
| W-STORE-LOCK-PATH | `cmd/ailang-worldd/operator_worktree.go:40` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-LOCK-NOT-RELEASED | `cmd/ailang-worldd/operator_worktree.go:44` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-SIGNED-STORE-IDENTITY | `cmd/ailang-worldd/operator_worktree.go:37` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-SIGNED-EPISODE-GRAMMAR | `cmd/ailang-worldd/operator_worktree.go:37` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-READ-RECONSTRUCTION | `cmd/ailang-worldd/operator_worktree.go:45` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-RESERVATION-RECHECK | `cmd/ailang-worldd/operator_worktree.go:77` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-CHALLENGE-LENGTH | `cmd/ailang-worldd/operator_worktree.go:74` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-EFFECT-PROOF | `cmd/ailang-worldd/operator_worktree.go:81` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-EFFECT-ATTEMPT | `cmd/ailang-worldd/operator_worktree.go:81` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-EFFECT-HOLD-ERR | `cmd/ailang-worldd/operator_worktree.go:83` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-RESERVE-ABSENT-ADD-ERROR | `cmd/ailang-worldd/operator_worktree.go:101` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/git_refused` | KILLED | §4.1 |
| W-UNCERTAIN-DEFAULT | `cmd/ailang-worldd/operator_worktree.go:100` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/path_only` | KILLED | §4.1 |
| W-TARGET-ACTION | `cmd/ailang-worldd/operator_worktree.go:156` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TARGET-REPO | `cmd/ailang-worldd/operator_worktree.go:156` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TARGET-COMMON | `cmd/ailang-worldd/operator_worktree.go:156` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TARGET-PATH | `cmd/ailang-worldd/operator_worktree.go:156` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TARGET-BRANCH | `cmd/ailang-worldd/operator_worktree.go:156` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TARGET-SHA | `cmd/ailang-worldd/operator_worktree.go:156` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TARGET-ABS | `cmd/ailang-worldd/operator_worktree.go:132` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TARGET-CLEAN | `cmd/ailang-worldd/operator_worktree.go:132` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-REPO-COMMON | `cmd/ailang-worldd/operator_worktree.go:185` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-REG-HEADER | `cmd/ailang-worldd/operator_worktree.go:198` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-REG-ABS-PATH | `cmd/ailang-worldd/operator_worktree.go:198` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-REG-COUNT | `cmd/ailang-worldd/operator_worktree.go:235` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-TARGET-DIR | `cmd/ailang-worldd/operator_worktree.go:235` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-TARGET-PHYSICAL | `cmd/ailang-worldd/operator_worktree.go:239` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-DOT-GIT-TYPE | `cmd/ailang-worldd/operator_worktree.go:249` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-LINK-CHECK | `cmd/ailang-worldd/operator_worktree.go:249` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-CHECKOUT-COMMON | `cmd/ailang-worldd/operator_worktree.go:253` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-PATH-ABSENT | `cmd/ailang-worldd/operator_worktree.go:222` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-LINK-ABSENT | `cmd/ailang-worldd/operator_worktree.go:242` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-ACCESS | `host/store/operator_worktree.go:15` | `host/store` | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED |
| W-TX-READ-CASES | `host/store/operator_worktree.go:24` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-CLAIMED | `host/store/operator_worktree.go:31` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-OWNER | `host/store/operator_worktree.go:31` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-CHALLENGE | `host/store/operator_worktree.go:31` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-IDENTITY | `host/store/operator_worktree.go:31` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-SERVED-STORE | `host/store/operator_worktree.go:38` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-START-KIND | `host/store/operator_worktree.go:38` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-HEADLESS | `host/store/operator_worktree.go:38` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-APPROVAL-BIND | `host/store/operator_worktree.go:41` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-APPROVER | `host/store/operator_worktree.go:45` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-ENROLLED | `host/store/operator_worktree.go:48` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-OWNER-HB | `host/store/operator_worktree.go:53` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-NO-PRIOR-STEP | `host/store/operator_worktree.go:55` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-NO-COMMIT | `host/store/operator_worktree.go:68` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| D-TX-ACCESS | `host/store/operator_worktree.go:109` | `host/store` | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED |
| D-TX-NO-COMMIT | `host/store/operator_worktree.go:145` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-TX-STEP-GENERATION | `host/store/operator_worktree.go:122` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/path_only` | KILLED | §4.1 |
| D-TX-RELEASE-OWNER | `host/store/operator_worktree.go:136` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-READ-ACTIVE-STATE | `host/store/operator_worktree.go:80` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-READ-OWNER | `host/store/operator_worktree.go:89` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-READ-PHASE | `host/store/operator_worktree.go:89` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-CONFIG | `cmd/ailang-worldd/operator_worktree_reconcile.go:53` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-STEP-LOAD | `cmd/ailang-worldd/operator_worktree_reconcile.go:57` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-STEP-JSON | `cmd/ailang-worldd/operator_worktree_reconcile.go:32` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-STORE-IDENTITY | `cmd/ailang-worldd/operator_worktree_reconcile.go:39` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-START-KIND | `cmd/ailang-worldd/operator_worktree_reconcile.go:39` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-HEADLESS | `cmd/ailang-worldd/operator_worktree_reconcile.go:39` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-EPISODE-ID | `cmd/ailang-worldd/operator_worktree_reconcile.go:40` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-EPISODE-BASENAME | `cmd/ailang-worldd/operator_worktree_reconcile.go:40` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-RESERVATION-ID | `cmd/ailang-worldd/operator_worktree_reconcile.go:41` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-RESERVATION-PATH | `cmd/ailang-worldd/operator_worktree_reconcile.go:41` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-RESERVATION-BRANCH | `cmd/ailang-worldd/operator_worktree_reconcile.go:41` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-REPO-PATH | `cmd/ailang-worldd/operator_worktree_reconcile.go:42` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-COMMON-PATH | `cmd/ailang-worldd/operator_worktree_reconcile.go:42` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-EPISODE-LOCK | `cmd/ailang-worldd/operator_worktree_reconcile.go:61` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-LOCKED-REREAD | `cmd/ailang-worldd/operator_worktree_reconcile.go:70` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-LOCKED-EPISODE | `cmd/ailang-worldd/operator_worktree_reconcile.go:70` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-READBACK-PRE | `cmd/ailang-worldd/operator_worktree_reconcile.go:78` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-READBACK-POST | `cmd/ailang-worldd/operator_worktree_reconcile.go:92` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-CLI-PARSE | `cmd/ailang-worldd/operator_service.go:275` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-CLI-DISPATCH | `cmd/ailang-worldd/operator_service.go:289` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| W-CANONICAL-CREDENTIALS | `cmd/ailang-worldd/operator_worktree.go:33` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-LOCK-BEFORE-RECONSTRUCTION | `cmd/ailang-worldd/operator_worktree.go:40` | `cmd/ailang-worldd` | `TestOperatorStartCollisions/held` | KILLED | §4.1 |
| W-CANONICAL-PATH-ERR | `cmd/ailang-worldd/operator_worktree.go:33` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-EPISODE-LOCK-ERR | `cmd/ailang-worldd/operator_worktree.go:40` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-GIT-PATH-ERR | `cmd/ailang-worldd/operator_worktree.go:49` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TARGET-ERR | `cmd/ailang-worldd/operator_worktree.go:53` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-LOAD-ACTION-ERR | `cmd/ailang-worldd/operator_worktree.go:70` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-PROOF-ERR | `cmd/ailang-worldd/operator_worktree.go:77` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-CREATED-FINAL | `cmd/ailang-worldd/operator_worktree.go:113` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-CREATED-RECEIPT | `cmd/ailang-worldd/operator_worktree.go:105` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-ABSENT-READ-ACK | `cmd/ailang-worldd/operator_worktree.go:101` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/git_refused` | KILLED | §4.1 |
| W-ABSENT-REPLY | `cmd/ailang-worldd/operator_worktree.go:116` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/git_refused` | KILLED | §4.1 |
| W-SIGNED-PATH-FACT | `cmd/ailang-worldd/operator_worktree.go:124` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-PHYSICAL-ANCESTOR | `cmd/ailang-worldd/operator_worktree.go:135` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-PHYSICAL-REPO | `cmd/ailang-worldd/operator_worktree.go:144` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-PHYSICAL-PATH | `cmd/ailang-worldd/operator_worktree.go:148` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-PHYSICAL-COMMON | `cmd/ailang-worldd/operator_worktree.go:152` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-COMMON-COMMAND | `cmd/ailang-worldd/operator_worktree.go:160` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-COMMON-ERR | `cmd/ailang-worldd/operator_worktree.go:185` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-LIST-ERR | `cmd/ailang-worldd/operator_worktree.go:188` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-LIST-PHYSICAL-ERR | `cmd/ailang-worldd/operator_worktree.go:201` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-REG-PHYSICAL-MATCH | `cmd/ailang-worldd/operator_worktree.go:205` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-REG-RAW-MATCH | `cmd/ailang-worldd/operator_worktree.go:205` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-REG-PINNED-ENTRY | `cmd/ailang-worldd/operator_worktree.go:207` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-REF-ERR | `cmd/ailang-worldd/operator_worktree.go:212` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-REF-EXACT-NAME | `cmd/ailang-worldd/operator_worktree.go:218` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-STAT-ERR | `cmd/ailang-worldd/operator_worktree.go:229` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-STAT-PHYSICAL-ERR | `cmd/ailang-worldd/operator_worktree.go:239` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-LSTAT-LINK-ERR | `cmd/ailang-worldd/operator_worktree.go:246` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-CHECKOUT-COMMON-ERR | `cmd/ailang-worldd/operator_worktree.go:253` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| GIT-CHECKOUT-HEAD-ERR | `cmd/ailang-worldd/operator_worktree.go:257` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| D-STORE-OPEN-ERR | `cmd/ailang-worldd/operator_worktree_reconcile.go:21` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-RESERVED-ROW-ERR | `cmd/ailang-worldd/operator_worktree_reconcile.go:28` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-PINNED-PATH-ERR | `cmd/ailang-worldd/operator_worktree_reconcile.go:35` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-SIDECAR-ERR | `cmd/ailang-worldd/operator_worktree_reconcile.go:61` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-REREAD-ERR | `cmd/ailang-worldd/operator_worktree_reconcile.go:66` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-GIT-PATH-ERR | `cmd/ailang-worldd/operator_worktree_reconcile.go:73` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-WRITE-RESULT | `cmd/ailang-worldd/operator_worktree_reconcile.go:103` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| W-TX-LOCK-NIL | `host/store/operator_worktree.go:38` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-OPEN-ERR | `host/store/operator_worktree.go:18` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-ROW-ERR | `host/store/operator_worktree.go:24` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-DECODE-ERR | `host/store/operator_worktree.go:34` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-BIND-ERR | `host/store/operator_worktree.go:41` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-ENROLMENT-ERR | `host/store/operator_worktree.go:48` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-UPDATE-ERR | `host/store/operator_worktree.go:60` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| W-TX-RESERVE-ERR | `host/store/operator_worktree.go:65` | `cmd/ailang-worldd` | `TestOperatorStartOnce` | KILLED | §4.1 |
| D-READ-UNCERTAIN | `host/store/operator_worktree.go:80` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/path_only` | KILLED | §4.1 |
| D-READ-ACTION-ERR | `host/store/operator_worktree.go:76` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-READ-DECODE-ERR | `host/store/operator_worktree.go:83` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-READ-STEP-ERR | `host/store/operator_worktree.go:93` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-TX-KIND-GUARD | `host/store/operator_worktree.go:112` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-TX-BEGIN-ERR | `host/store/operator_worktree.go:115` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-TX-UNCERTAIN-CAS-ERR | `host/store/operator_worktree.go:125` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/path_only` | KILLED | §4.1 |
| D-TX-OBSERVE-ERR | `host/store/operator_worktree.go:129` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| D-TX-ACTION-ERR | `host/store/operator_worktree.go:142` | `cmd/ailang-worldd` | `TestOperatorWorktreeWindow/before_add` | KILLED | §4.1 |
| G-STEP-ERR | `host/store/operator_steps.go:68` | `host/store` | `TestOperatorStepRestart` | KILLED | KILLED |

**Judge-433 follow-up.** The four newly meaningful prior survivors (N-WORKTREE-DERIVED, N-ATTEMPT-ZERO, N-CLOCK-CONST and N-WORKSPACEROOT-EMPTY) now fail their named seam tests. N-HOSTSDIR-EMPTY and N-WORKSPACE-SWAP also die. The two controls below were re-applied and verified to survive.

- **N-STORE-RAW-CRED** — Equivalent for canonical S5 seam input; credentialsFor already canonicalizes the path and reservation independently pins store identity. Named run rc=0.
- **N-RES-ABS** — Equivalent for mint token validation; runnerPath rejects non-absolute parents. Named run rc=0.

**Outside sandbox.** The ID list in `whole-mutant-ids.txt` contains exactly the §4.1 rows above. The previous mint resolver interleavings and process survivors stay with the paired-process harness (Appendix C).

### 4.1 Whole-package runs for the controller

Use the single scratch launcher `$S/whole-mutants.sh` outside the sandbox after the plan commit. It runs `mutate.py --whole-only` against the immutable planner stage exports and records one log per ID, with `go test -count=1 -timeout 8m -v ./PACKAGE` and the same pinned environment. CLI whole-package runs reach loopback fixtures; named test runs use no socket. Remaining rows and commands are listed with the mutation results above. Full CLI runs can take minutes per ID. The controller should let the foreground runner finish and keep its ledger.

## 5. Files

The file-by-file measured manifest appears in §1.5. M1 changes the credentials bank, its test inventory and the S5c read seam/tests. M2 adds the inert preparation and attended worktree files and updates the CLI usage and QUICKSTART. There are no package/module, schema, route, frontend, driver, launcher, daemon or verifier test changes. §6 compares the served sources and route census and includes the existing live finish authority call as a positive control for the caller fence.

## 6. Sprint gates

````bash gates.sh
#!/bin/bash
# Usage: gates.sh <tree> <logdir> <base|final|scratch> [sha]
set -u
set -o pipefail
export GIT_OPTIONAL_LOCKS=0
T=$1; L=$2; MODE=$3
B=68c3673579b53326a2c5d4d7bc62b9d3bb7a1105
REF=${GATES_GIT_REF:-$T}
SHA=${4:-$(git -C "$REF" rev-parse HEAD)}
mkdir -p "$L"
cd "$T" || exit 2
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export AILANG_STATE_DIR="$L/ailang-state" AILANG_CACHE_DIR="$L/ailang-cache"
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
mkdir -p "$AILANG_STATE_DIR"
: > "$L/legs.txt"
FAILS=
UNKNOWN=
s0=$(date +%s)
line() { echo "$*" | tee -a "$L/legs.txt"; }
leg() {
  local name=$1 cmd=$2 start rc
  start=$(date +%s)
  (set -eo pipefail; eval "$cmd") 2>&1 | tee "$L/$name.txt"
  rc=${PIPESTATUS[0]}
  if [ "$rc" -ne 0 ] && [ "${GATES_SANDBOX:-0}" = 1 ] &&
     { grep -Fq 'OBSERVATION REQUIRES CONTROLLER' "$L/$name.txt" ||
       { grep -Eq 'listen .*bind: operation not permitted' "$L/$name.txt" &&
         ! grep -E '^ *--- FAIL:' "$L/$name.txt" | grep -Fv -e '--- FAIL: TestApproveDecideWriteMessageAndSignature ' -e '--- FAIL: TestQuickstartSection9FlagsMatchTheCLI '; }; }; then
    line "LEG $name rc=$rc $(( $(date +%s)-start ))s UNINFORMATIVE UNDER SANDBOX"
    UNKNOWN="$UNKNOWN $name"
  else
    line "LEG $name rc=$rc $(( $(date +%s)-start ))s"
    [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
  fi
}
observe() {
  if [ "${GATES_SANDBOX:-0}" = 1 ]; then
    echo 'OBSERVATION REQUIRES CONTROLLER (operator file observations)'
    return 77
  fi
  local d file creds
  for d in "$HOME/Library/LaunchAgents" "$HOME/.ailang-world/bin" "$HOME/.ailang-world"; do
    if [ ! -d "$d" ]; then printf 'absent %s\n' "$d"; continue; fi
    echo "DIRECTORY $d"
    if [ "$d" = "$HOME/.ailang-world" ]; then
      find "$d" -maxdepth 1 \( -type f -o -type l \) -name '*.json' -print0
    else
      find "$d" \( -type f -o -type l \) -print0
    fi | LC_ALL=C sort -z | while IFS= read -r -d '' file; do
      /bin/ls -la -d "$file"
      if [ -L "$file" ]; then readlink "$file" | shasum -a 256; else shasum -a 256 "$file"; wc -l "$file"; fi
    done
    if [ "$d" = "$HOME/.ailang-world" ]; then
      for file in "$d/service.json" "$d/operator-runner.json"; do
        if [ ! -e "$file" ] && [ ! -L "$file" ]; then printf 'absent %s\n' "$file"; fi
      done
    fi
  done
  creds=${GATES_REAL_CREDS:-$HOME/.ailang/world/credentials.db}
  for file in "$creds" "$creds".operator-episode-*.lock; do
    if [ -f "$file" ]; then shasum -a 256 "$file"; wc -l "$file"
    else printf 'absent %s\n' "$file"; fi
  done
  echo 'Service process write scope is checked with fixture supervisor counters'
}
leg observation-before 'observe'
paths() {
  case "$MODE" in
    base) git -C "$REF" diff --name-only "$B" "$B" ;;
    scratch) cat "$GATES_CHANGED_FILE" ;;
    final) git -C "$REF" diff --name-only "$B" HEAD ;;
    *) return 2 ;;
  esac
}
leg paths 'paths'
TOUCHED=$(awk -F/ '/\.go$/ && ($1=="cmd" || $1=="host") {NF--; OFS="/"; $1=$1; print}' "$L/paths.txt" | sort -u | tr '\n' ' ')
EXPECTED='cmd/ailang-worldd host/store '
line "TOUCHED_PACKAGES: $TOUCHED"
leg touched 'if [ "$MODE" = base ]; then test -z "$TOUCHED"; else test "$TOUCHED" = "$EXPECTED"; fi'
leg status 'git -C "$REF" status --porcelain'
leg gofmt 'out=$(gofmt -l cmd host); printf "%s\n" "$out"
  if [ -n "$out" ]; then
    test "$out" = host/daemon/daemon.go
    git -C "$REF" show "$B:host/daemon/daemon.go" | cmp - "$out"
    echo BASELINE-FORMAT-DRIFT
  fi' 
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./..."
leg verify_ail './scripts/verify_ail.sh'
PACKAGES=$TOUCHED
if [ "$MODE" = base ]; then PACKAGES=$EXPECTED; fi
for p in $PACKAGES; do
  tag=${p//\//-}
  leg "pkg-$tag" "go test -count=1 -timeout 8m ./$p"
  leg "race-$tag" "/usr/bin/time -p go test -count=1 -race -timeout 8m ./$p"
done
inv() {
  local p=$1 name n t
  shift
  n=$(printf '%s\n' "$*" | tr ' ' '|')
  name="inv-${p//\//-}"
  leg "$name" "go test -v -count=1 -timeout 8m ./$p -run '^($n)\$'"
  for t in "$@"; do
    if ! grep -qE "^--- PASS: $t " "$L/$name.txt"; then
      line "NOT PASSED ($name): $t"
      FAILS="$FAILS ${name}_missing"
    fi
  done
}
inv host/approval TestApprovalImportsOnlyTheStandardLibrary TestOperatorRunnerRecheckLaw TestOperatorPureEntryPointsWriteNothing
inv host/approval/approvaltest TestApprovalTestIsTestOnly TestTestSignerAgreesWithSSHKeygen
inv host/authority TestVerifiedOperatorActionOnlyFromTheJudges TestOperatorRunnerRechecks TestOperatorJudgesWriteNothing TestOperatorStoreWritesNothing
inv host/boundary TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard
inv host/daemon TestDaemonDependencyAllowlist TestDaemonCredentialWriterSites TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestConfirmRoutesRegistered TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes
inv host/verbs TestOperatorNeverDispatch TestOperatorOtherVerbParity
inv host/broker TestBrokerDependencyAllowlist
inv cmd/ailang-worldd TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestCLICredentialWriterSites TestOperatorFlagParity TestLaunchAndStartHelp TestQuickstartSection19MatchesTheBar TestQuickstartSection20MatchesTheFinishFlag TestQuickstartServiceSectionMatchesTheCLIHelp TestPluginCommandsMatchTheCLIHelp TestHelpIsNeverBlockedByTheRegistryKeyGuard TestHomeServiceSeamOnly TestHomeServiceInstallAndPlist TestHomeServiceUpgradeRollback TestOperatorServiceFence TestOperatorServiceOwnership TestOperatorServiceConfig TestOperatorServicePair TestOperatorServiceStartup TestOperatorServiceUpgrade TestOperatorServiceUpgradePairFence TestOperatorServiceWriteScope TestOperatorServiceEnvironment
leg inv-quickstart-section9 "go test -v -count=1 -timeout 8m ./cmd/ailang-worldd -run '^TestQuickstartSection9FlagsMatchTheCLI\$'"
if [ "${GATES_SANDBOX:-0}" != 1 ] && ! grep -Fq -- '--- PASS: TestQuickstartSection9FlagsMatchTheCLI (' "$L/inv-quickstart-section9.txt"; then
  FAILS="$FAILS inv-quickstart-section9_missing"
fi
inv host/store TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL TestOperatorSchema10 TestOperatorClaimBeforeExpiry TestOperatorClaimRace TestOperatorTerminalExclusion TestOperatorFinishReadiness TestWriterLockCrossProcessContention TestWriterLockHeldIsNonMutating
route_fence() {
  local routes posts old_routes old_posts
  routes=$(grep -Ec 'mux\.(HandleFunc|Handle)\(' host/daemon/daemon.go)
  posts=$(grep -Fc 'mux.HandleFunc("POST /surface/' host/daemon/daemon.go)
  old_routes=$(git -C "$REF" show "$B:host/daemon/daemon.go" | grep -Ec 'mux\.(HandleFunc|Handle)\(')
  old_posts=$(git -C "$REF" show "$B:host/daemon/daemon.go" | grep -Fc 'mux.HandleFunc("POST /surface/')
  printf 'DAEMON_ROUTES: %s (base %s); SURFACE_POSTS: %s (base %s)\n' "$routes" "$old_routes" "$posts" "$old_posts"
  test "$routes" -gt 0
  test "$routes" = "$old_routes"
  test "$posts" -gt 0
  test "$posts" = "$old_posts"
  local f
  for f in host/daemon/daemon.go host/daemon/surface_page.go host/surfacepage/app/app.js; do
    git -C "$REF" show "$B:$f" | cmp - "$f"
    printf 'UNCHANGED: %s\n' "$f"
  done
  grep -n 'surfaceCSP' host/daemon/surface_page.go
}
leg routes 'route_fence'
inv cmd/world-publish TestPluginCommandsMatchTheCLIHelp
leg mint-quickstart "go test -v -count=1 -timeout 1m ./cmd/ailang-worldd -run '^(TestQuickstartRunnerMintRecovery|TestOperatorMintRecovery)$'"
leg start-input "go test -v -count=1 -timeout 2m ./cmd/ailang-worldd -run '^(TestOperatorStartFacts|TestOperatorStartChanged)$'"
if [ "$MODE" = base ]; then
  line 'LEG s5d-accept SKIPPED 0s (base has no S5d tests)'
  line 'LEG added-cli-race SKIPPED 0s (base has no S5d tests)'
else
  leg s5d-accept '/bin/bash "$(dirname "$0")/accept.sh" "$T" M2 fast'
  leg added-cli-race 'go test -race -v -count=1 -timeout 2m ./cmd/ailang-worldd -run "^(TestOperatorStartForwarding|TestOperatorWorktreeWindow|TestOperatorStartOnce|TestOperatorStartCollisions|TestOperatorWorktreeDispositionCAS)$"'
fi
scope_fence() {
  local allowed="$L/allowed.txt" extra
  cat > "$allowed" <<'EOF'
cmd/ailang-worldd/help_test.go
cmd/ailang-worldd/main.go
cmd/ailang-worldd/operator_service.go
cmd/ailang-worldd/operator_start_input.go
cmd/ailang-worldd/operator_start_input_test.go
cmd/ailang-worldd/operator_worktree.go
cmd/ailang-worldd/operator_worktree_reconcile.go
cmd/ailang-worldd/operator_worktree_test.go
cmd/ailang-worldd/service.go
docs/QUICKSTART.md
host/store/credentials_test.go
host/store/operator_cas_test.go
host/store/operator_steps.go
host/store/operator_worktree.go
design_docs/planned/w-world-address-bar-verbs-v4s5d-sprint-plan.md
design_docs/planned/sprint_w-world-address-bar-verbs-v4s5d.json
EOF
  printf '%s\n' cmd/ailang-worldd/operator_start_input.go | grep -Fx -f "$allowed"
  printf '%s\n' cmd/ailang-worldd/episode_preset.go | grep -Fxv -f "$allowed"
  printf '%s\n' world/transitions.ail | grep -Fxv -f "$allowed"
  extra=$(grep -Fxv -f "$allowed" "$L/paths.txt" || test "$?" -eq 1)
  printf 'UNTOUCHED_VIOLATIONS: %s\n' "${extra:-none}"
  test -z "$extra"
}
leg scope 'scope_fence'
production() {
  grep -R -nF --include='*.go' "$1" cmd host | grep -Fv '_test.go:'
}
callers() {
  production "$1(" | grep -Ev "func (\([^)]*\) )?$1\(" || :
}
one_site() {
  test "$(printf '%s\n' "$1" | wc -l | tr -d ' ')" = 1
  printf '%s\n' "$1" | grep -F "$2:"
}
caller_fence() {
  local count hits name
  production 'authority.RecheckVerb('
  production 'func MintOperatorStart('
  production ' ClaimOperatorAction('
  production 'func reconstructOperatorStart('
  for name in MintOperatorStart ClaimOperatorAction reserveOperatorStartWorktree; do
    test -z "$(callers "$name")"
    printf 'Production entry callers of %s: 0\n' "$name"
  done
  count=$(grep -R -nE --include='*.go' 'syscall\.Flock.*LOCK_NB' cmd host | grep -Fv '_test.go:' | wc -l | tr -d ' ')
  printf 'Production Flock LOCK_NB lines (with comment): %s\n' "$count"
  test "$count" = 7
  hits=$(callers provisionEpisode)
  printf '%s\n' "$hits"
  test "$(printf '%s\n' "$hits" | wc -l | tr -d ' ')" = 5
  one_site "$(production '.CommitOperatorMint(')" host/authority/mint.go
  one_site "$(production '.AbortOperatorMint(')" cmd/ailang-worldd/operator_mint.go
  one_site "$(callers operatorStartFacts)" cmd/ailang-worldd/operator_start_input.go
  production 'resolvePresetStart' | grep -F 'cmd/ailang-worldd/episode_preset.go:'
  production 'resolvePresetStart' | grep -F 'cmd/ailang-worldd/operator_start_input.go:'
  if [ "$MODE" = base ]; then
    test -z "$(callers reconstructOperatorStart)"
    test ! -e host/store/operator_steps.go
  else
    one_site "$(callers reconstructOperatorStart)" cmd/ailang-worldd/operator_worktree.go
    one_site "$(production '.ReserveOperatorWorktree(')" cmd/ailang-worldd/operator_worktree.go
    one_site "$(production '.ObserveOperatorStep(')" cmd/ailang-worldd/operator_worktree.go
    test -z "$(production '.ReserveOperatorStep(' || :)"
    one_site "$(production '.ReadOperatorWorktree(')" cmd/ailang-worldd/operator_worktree_reconcile.go
    one_site "$(production '.reconcileRunnerWorktree(')" cmd/ailang-worldd/operator_service.go
    hits=$(production '.DisposeOperatorWorktree(')
    printf '%s\n' "$hits" | grep -F 'cmd/ailang-worldd/operator_worktree.go:'
    printf '%s\n' "$hits" | grep -F 'cmd/ailang-worldd/operator_worktree_reconcile.go:'
    test "$(printf '%s\n' "$hits" | wc -l | tr -d ' ')" = 2
  fi
}
leg callers 'caller_fence'
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ] || [ "$MODE" != final ]; then
  echo 'Verify on the committed checkout (tracked-blob gate needs .git)' > "$L/verify_go.txt"
  line "LEG verify_go SKIPPED 0s (run in committed checkout)"
else
  leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN="$HOME/.pinned-ailang/ailang" ./scripts/verify_go.sh'
fi
leg observation-after 'observe'
if ! cmp -s "$L/observation-before.txt" "$L/observation-after.txt"; then
  diff -u "$L/observation-before.txt" "$L/observation-after.txt" | tee "$L/observation-diff.txt"
  FAILS="$FAILS observation_changed"
fi
line "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none} UNINFORMATIVE:${UNKNOWN:- none}"
tail -n 1 "$L/legs.txt" > "$L/ALL_DONE.$SHA"
[ -z "$FAILS" ] || exit 1
[ -z "$UNKNOWN" ] || exit 77
````

## 7. Executor directive notes

`S=$HOME/.ailang/state/world-iter434/plan-scratch` holds pristine `base` and cumulative `st-M1/st-M2`, `numstat.txt`, `preflight.txt`, `anchors.txt`, and the verification logs. The executor extracts the scripts from this plan. `accept.sh` runs named test groups without loopback and fails for a missing named PASS. It points AILANG state and cache to scratch, leaves TMPDIR/GOTMPDIR alone, and sets the pinned binary, SSH key fixture requirement, Node and tool AILANG environment.

````bash accept.sh
#!/bin/bash
set -euo pipefail
W=$1; M=$2; PART=$3
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
S="$HOME/.ailang/state/world-iter434/plan-scratch"
export AILANG_STATE_DIR="$S/accept-state" AILANG_CACHE_DIR="$S/accept-cache"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
cd "$W"
test "$PART" = fast
run() {
  local pkg=$1 names=$2 name log
  log=$(mktemp "$S/accept.XXXXXX")
  go test -v -count=1 -timeout 2m "./$pkg" -run "^($names)\$" | tee "$log"
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log"
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
case "$M" in
  M1|M2) ;;
  *) exit 2 ;;
esac
run host/store 'TestOperatorStepRestart|TestOperatorStepGuards|TestCredentialsMethodsRequireADeadline|TestOperatorSchema10'
run cmd/ailang-worldd 'TestOperatorStartForwarding|TestOperatorStartFacts|TestOperatorStartChanged'
if [ "$M" = M2 ]; then
  run cmd/ailang-worldd 'TestOperatorWorktreeWindow|TestOperatorStartOnce|TestOperatorStartCollisions|TestOperatorWorktreeDispositionCAS|TestEveryVerbHasHelp|TestQuickstartServiceSectionMatchesTheCLIHelp|TestCLICredentialWriterSites|TestOperatorMintRecovery|TestQuickstartRunnerMintRecovery'
fi
printf 'ACCEPT_OK %s %s\n' "$M" "$PART"
````

The snapshot path is `$HOME/.ailang/state/world-iter434/snap/$M`; filenames and hashes are cumulative against base. Both plan files are excluded by the pathspec. The scratch override is for rehearsal of archive trees only. Real executor snapshots use the read-only git status/ls-files branches of this script, then `shasum -a 256 -c MANIFEST.sha256`.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
B=68c3673579b53326a2c5d4d7bc62b9d3bb7a1105
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter434/snap}/$M
mkdir -p "$(dirname "$D")"
mkdir "$D"
cd "$W"
if [ -n "${SNAP_STATUS_LIST:-}" ]; then
  cp "$SNAP_STATUS_LIST" "$D/STATUS"
else
  {
    git diff --name-status --no-renames "$B" -- cmd host docs app
    git ls-files --others --exclude-standard -- cmd host docs app | sed $'s/^/A\t/'
  } | LC_ALL=C sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
fi
test -s "$D/STATUS"
awk -F '\t' '$1 != "D" {print $2}' "$D/STATUS" | while IFS= read -r p; do
  mkdir -p "$D/$(dirname "$p")"
  cp -p "$p" "$D/$p"
done
awk -F '\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
(
  cd "$D"
  while IFS="$(printf '\t')" read -r kind p; do
    if [ "$kind" != D ]; then shasum -a 256 "$p"; fi
  done < STATUS > MANIFEST.sha256
  shasum -a 256 -c MANIFEST.sha256
)
wc -l "$D"/{STATUS,DELETED,MANIFEST.sha256}
````

Use these commit-time commands. The controller may run the full gate after both commits. Do not make the commits on the executor lane.

````bash
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter434-row243-v4s5d
T=$HOME/.ailang/state/world-iter434/executor
S=$HOME/.ailang/state/world-iter434/plan-scratch
B=68c3673579b53326a2c5d4d7bc62b9d3bb7a1105
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5d-sprint-plan.md
STATE=$HOME/.ailang/state/world-iter434
mkdir -p "$T"
printf '%s\n' 'feat(store): bank operator step reservation and observation' > "$T/M1.msg"
printf '%s\n' 'feat(operator): reserve start worktree with attended disposition' > "$T/M2.msg"
bash "$STATE/commit_ms.sh" M1 "$T/M1.msg" 'AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang /bin/bash "'"$T"'/accept.sh" "'"$W"'" M1 fast'
bash "$STATE/commit_ms.sh" M2 "$T/M2.msg" 'AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang /bin/bash "'"$T"'/accept.sh" "'"$W"'" M2 fast'
/bin/bash "$T/gates.sh" "$W" "$S/gate-controller-final" final
/bin/bash "$S/whole-mutants.sh"
````

The final checkpoint was extracted with the controller's awk. Each patch applied with /usr/bin/patch -p1 -E to a fresh archive, and full-archive `diff -r` returned 0 for both milestones. Separate controller comparisons for cmd, host and docs also returned 0. The optional app directory was absent. No .orig or .rej was produced. The three fenced scripts were byte-identical to the scratch originals.

Snapshot rehearsal used the extracted script on both archive trees with read-only `GIT_DIR=$(git rev-parse --absolute-git-dir)`, `GIT_WORK_TREE=<archive>`, `GIT_OPTIONAL_LOCKS=0` and `SNAP_PARENT="$S/snap-rehearsal-bounded"`. This exercised both git diff and untracked discovery, with no injected status list. STATUS matched the measured manifest and all hashes verified. M1 was 5/0/5 and M2 was 14/0/14. Evidence is in `roundtrip-bounded.log`, `roundtrip-bounded-{M1,M2}.diff` and `snapshot-bounded-{M1,M2}.log`.

## 8. Risks

- Checkout contents alone cannot show that the runner created a worktree. Only uninterrupted success under the episode lock plus complete readback makes a created receipt; disposition after interruption is attended and preserves the external effects.
- A lost or incompatible Git read preserves ownership or yields retained-stop after review of remaining evidence. If the original repo is unreadable, restore read access to the pinned repo and re-run the printed remedy. Keep the episode lockfile in place across restart.
- The paired-process sprint must integrate the same canonical open file description across claim and effects and verify live software/configuration, quiescence, recovery and terminal drain. Until then no runner loop, browser call or terminal provision path can enter the new preparation seam.
- Test fakes use a macOS service environment with fake supervisor and transport. Git path and lock portability still require Linux CI. The known baseline daemon format drift and denied CLI socket tests are visible in §1.1.

## 9. Acceptance criteria

1. Both embedded product patches apply cumulatively with `patch -p1 -E` and byte-match their scratch trees. Filtered cumulative numstat is at most 1,400; unfiltered product scope and manifest also match.
2. M1 is independently green and effect-free from live entry points. Forwarding has named negative tests before M2 introduces the inert caller.
3. M2 passes the §5.3 worktree row tests, all four injected boundaries, store/readback invariants and the typed, verbatim disposition including owner/step conflicts. All new worktree uncertainty is covered by the attended command in M2.
4. There is no live reservation, add, mint, launch, daemon or UI dispatch. The attended entry rejects absent, stale or mismatched records and keeps all Git/file evidence. Gates distinguish absence, readback, true failure and sandbox-uninformative runs.
5. The §5.5 inventory runs by name and remaining activation gates travel with Appendix C. Controller replays whole-package mutants and runs the gate on the committed checkout and the Linux CI jobs.

## Appendix C — after this sprint (measured)

- **S5e file publication and disposition.** 53-line measured prior fd readback probe (C12); full size UNMEASURED. Include private no-replace publication, independent token/task/MCP path, type, ownership/mode/link/digest and changed-file cases, observable readback and its attended disposition in one patch. The worktree and credential remain available as durable evidence.
- **S5f/S6 paired-process consumer, owner recovery and drain.** 72-line measured prior fd probe (C12); full size UNMEASURED. S3b step bank is accounted in S5d. Keep `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorOldBuildTerminalDrain`, `TestOperatorCurrentOwnerProcessRecovery` and `TestOperatorServiceUpgrade/active_child` together. Carry `E-RES-EPISODE-CONFLICT`, `E-RES-REREAD` (the mint resolver needs an interleaving before lock), `J-OWN-UID` (alternate owner), `Z-ABSENT-ANY` and `Z-WAIT-READBUDGET` (process wait/drain). Also carry terminal `session revoke` exclusion and independent current-session/expiry/revocation checks, exact signed executable identity and confinement. No live mint before an observed-created worktree receipt and before all paired-process gates.
- **Launch, browser, finish and herdr.** Design §8 budgets S6 890 raw, S8 860 raw and S12 700 raw, full implementations UNMEASURED. Preserve E3, E5/#604 app sequencing (§1.2 C11), E4 and D-WORLD-120/E2 finish readiness. This sprint ships no spawn, run receipt, browser, herdr, process stop/drain or remote effects.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/cmd/ailang-worldd/operator_start_input.go b/cmd/ailang-worldd/operator_start_input.go
index baf81e38..5576a2e1 100644
--- a/cmd/ailang-worldd/operator_start_input.go
+++ b/cmd/ailang-worldd/operator_start_input.go
@@ -8,15 +8,17 @@ import (
 	"github.com/sunholo-data/ailang-world/host/store"
 )
 
+type presetStartReader func(context.Context, string, string, string, presetStartInput, func() int64) (*presetStartRead, error)
+
 // operatorStartFacts projects one live read of a headless start into the
 // effective mint input and the complete, name-ordered resolution a reviewed
 // action signs. It derives the task only from the signed task bytes.
 func operatorStartFacts(ctx context.Context, signed approval.OperatorAction, attempt store.OperatorMintAttempt,
-	presets, profiles, workspace, hosts string, clock func() int64) (store.OperatorMintInput, []approval.OperatorResolved, error) {
+	presets, profiles, workspace, hosts string, clock func() int64, readPreset presetStartReader) (store.OperatorMintInput, []approval.OperatorResolved, error) {
 	if signed.Kind != "start" || !signed.Options.Headless {
 		return store.OperatorMintInput{}, nil, store.ErrOperatorConflict
 	}
-	read, err := resolvePresetStart(ctx, signed.Options.Preset, presets, profiles, presetStartInput{
+	read, err := readPreset(ctx, signed.Options.Preset, presets, profiles, presetStartInput{
 		episodes: []string{signed.Episode}, workspaceRoot: workspace, hostsDir: hosts, launching: true, headless: true,
 		readTask: func() (string, []byte, error) { return "inline", []byte(signed.Task), nil },
 	}, clock)
@@ -65,7 +67,7 @@ func operatorStartFacts(ctx context.Context, signed approval.OperatorAction, att
 // production caller until the paired-process consumer lands.
 func reconstructOperatorStart(ctx context.Context, signed approval.OperatorAction, attempt store.OperatorMintAttempt,
 	presets, profiles, workspace, hosts string, clock func() int64) (store.OperatorMintInput, error) {
-	in, facts, err := operatorStartFacts(ctx, signed, attempt, presets, profiles, workspace, hosts, clock)
+	in, facts, err := operatorStartFacts(ctx, signed, attempt, presets, profiles, workspace, hosts, clock, resolvePresetStart)
 	if err != nil {
 		return store.OperatorMintInput{}, err
 	}
diff --git a/cmd/ailang-worldd/operator_start_input_test.go b/cmd/ailang-worldd/operator_start_input_test.go
index 554bbe6b..3f056ca0 100644
--- a/cmd/ailang-worldd/operator_start_input_test.go
+++ b/cmd/ailang-worldd/operator_start_input_test.go
@@ -1,12 +1,15 @@
 package main
 
 import (
+	"context"
 	"encoding/json"
 	"errors"
 	"os"
 	"path/filepath"
 	"reflect"
+	"strings"
 	"testing"
+	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/daemon"
@@ -30,7 +33,7 @@ const operatorStartTask = "# Fix fixture\nExact signed task bytes.\n"
 func newOperatorStartRig(t *testing.T) operatorStartRig {
 	t.Helper()
 	f := newPresetFixture(t)
-	r := operatorStartRig{f: f, bin: startReadFakeAgent(t, f), hosts: t.TempDir(), at: store.OperatorMintAttempt{Catalogue: 3}}
+	r := operatorStartRig{f: f, bin: startReadFakeAgent(t, f), hosts: t.TempDir(), at: store.OperatorMintAttempt{Catalogue: 3, CAS: store.OperatorCAS{ID: fixtureActionID, Generation: 7}, Runner: store.OperatorOwner{Nonce: "fixture", Generation: 11}, StepVersion: 2}}
 	r.creds = matchingHealth(t, f.r.f, f.r.profiles, f.profile).CredentialsDB
 	canonical, err := store.CanonicalPath(r.creds)
 	if err != nil {
@@ -39,7 +42,7 @@ func newOperatorStartRig(t *testing.T) operatorStartRig {
 	r.action = approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "start", Episode: "ep",
 		Store: canonical, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: operatorStartTask, Review: "fixture"}
 	r.action.Argv = approval.OperatorArgv(r.action.Kind, r.action.Episode, r.action.Options)
-	in, facts, err := operatorStartFacts(presetTestContext(t), r.action, r.at, f.dir, f.r.profiles, f.r.f.ws, r.hosts, presetClock)
+	in, facts, err := operatorStartFacts(presetTestContext(t), r.action, r.at, f.dir, f.r.profiles, f.r.f.ws, r.hosts, presetClock, resolvePresetStart)
 	if err != nil {
 		t.Fatal(err)
 	}
@@ -100,7 +103,7 @@ func TestOperatorStartFacts(t *testing.T) {
 		}
 	}
 	in := r.in
-	if in.Start.Episode != "ep" || in.Start.Base.SHA != f.base || in.Start.Preset.SHA256 != serviceHash(f.data) ||
+	if in.Start.StartedAt != time.Unix(presetClock(), 0).Format(time.RFC3339) || in.Start.Episode != "ep" || in.Start.Base.SHA != f.base || in.Start.Preset.SHA256 != serviceHash(f.data) ||
 		in.Start.Branch != "world/ep" || in.Start.BranchExisted || in.Start.Agent != "codex" || in.Start.TTLS != 43200 ||
 		in.Start.Task == nil || in.Start.Task.Path != filepath.Join(f.r.sessions, "ep.task.md") ||
 		in.Start.Task.SHA256 != serviceHash([]byte(operatorStartTask)) || in.TokenPath != filepath.Join(f.r.sessions, "ep.session") ||
@@ -227,3 +230,29 @@ func TestOperatorStartChanged(t *testing.T) {
 		})
 	}
 }
+
+// The preset daemon can configure a workspace distinct from the layout root.
+// Headless starts still forward the hosts directory through the common read seam.
+func TestOperatorStartForwarding(t *testing.T) {
+	r := newOperatorStartRig(t)
+	f := r.f
+	before := observeProvision(t, f.r, r.creds)
+	workspace := t.TempDir()
+	reader := func(ctx context.Context, name, presets, profiles string, in presetStartInput, clock func() int64) (*presetStartRead, error) {
+		if in.hostsDir != r.hosts {
+			t.Errorf("headless start forwarded hosts %q, want %q", in.hostsDir, r.hosts)
+		}
+		return resolvePresetStart(ctx, name, presets, profiles, in, clock)
+	}
+	_, facts, err := operatorStartFacts(presetTestContext(t), r.action, r.at, f.dir, f.r.profiles, workspace, r.hosts, presetClock, reader)
+	if err != nil || !reflect.DeepEqual(facts, r.action.Resolution) {
+		t.Fatalf("different layout root must retain the configured worktree: %v, %v", facts, err)
+	}
+	_, err = reconstructOperatorStart(presetTestContext(t), r.action, r.at, f.dir, f.r.profiles, f.dir, r.hosts, presetClock)
+	if err == nil || !strings.Contains(err.Error(), "inside the workspace") {
+		t.Fatalf("a preset stored inside the supplied workspace must refuse: %v", err)
+	}
+	if !reflect.DeepEqual(before, observeProvision(t, f.r, r.creds)) {
+		t.Fatal("forwarding checks changed durable state")
+	}
+}
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index 0fab29aa..82aec381 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -308,6 +308,9 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	_, errs["PendingVerbRequests"] = c.PendingVerbRequests(ctx, 1000, 0)
 	errs["ConsumeVerbRequest"] = c.ConsumeVerbRequest(ctx, "x", "y", ApproverUse{}, "", nil, "", 1000)
 	_, errs["OperatorAction"] = c.OperatorAction(ctx, "x", 1000)
+	errs["ReserveOperatorStep"] = c.ReserveOperatorStep(ctx, OperatorStepKey{}, 1000)
+	errs["ObserveOperatorStep"] = c.ObserveOperatorStep(ctx, OperatorStepKey{}, nil)
+	_, errs["OperatorStep"] = c.OperatorStep(ctx, OperatorStepKey{})
 	errs["CommitOperatorMint"] = c.CommitOperatorMint(ctx, OperatorMintInput{}, nil, nil)
 	_, errs["ReadOperatorMint"] = c.ReadOperatorMint(ctx, "x")
 	errs["AbortOperatorMint"] = c.AbortOperatorMint(ctx, OperatorMintReceipt{}, 1000)
diff --git a/host/store/operator_cas_test.go b/host/store/operator_cas_test.go
index 6d528ac5..199c887a 100644
--- a/host/store/operator_cas_test.go
+++ b/host/store/operator_cas_test.go
@@ -1,7 +1,10 @@
 package store
 
 import (
+	"bytes"
 	"errors"
+	"fmt"
+	"path/filepath"
 	"sync"
 	"testing"
 )
@@ -121,3 +124,162 @@ func TestOperatorClaimRace(t *testing.T) {
 		t.Fatal("the same V3 card was linked to two actions")
 	}
 }
+
+func TestOperatorStepRestart(t *testing.T) {
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	open := func() *Credentials {
+		c, err := OpenCredentials(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		t.Cleanup(func() { c.Close() })
+		return c
+	}
+	c := open()
+	cas, runner := operatorClaimFixture(t, c)
+	ctx := testCtx(t)
+	if err := c.ClaimOperatorAction(ctx, cas, runner, 1599); err != nil {
+		t.Fatal(err)
+	}
+	cas.Generation++
+	key := OperatorStepKey{Action: cas, Owner: runner.Nonce, Kind: "launch", Identity: "run-one", Version: 1}
+	if err := c.ReserveOperatorStep(ctx, key, 1600); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("reserve at expiry: %v", err)
+	}
+	operatorCounts(t, c, 1, 1)
+	if err := c.ReserveOperatorStep(ctx, key, 1599); err != nil {
+		t.Fatal(err)
+	}
+	if err := c.Close(); err != nil {
+		t.Fatal(err)
+	}
+	c = open()
+	if err := c.ReserveOperatorStep(ctx, key, 1599); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("reservation replay: %v", err)
+	}
+	reservation, err := c.OperatorStep(ctx, key)
+	if err != nil || reservation.Phase != "reserved" || reservation.Version != 1 || reservation.Owner != runner.Nonce || reservation.Receipt != nil {
+		t.Fatalf("lost reservation on restart: %+v %v", reservation, err)
+	}
+	receipt := []byte("observed-test-launch")
+	wrong := key
+	wrong.Version++
+	if err := c.ObserveOperatorStep(ctx, wrong, receipt); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("wrong step CAS: %v", err)
+	}
+	if got, err := c.OperatorStep(ctx, key); err != nil || got.Phase != "reserved" || got.Receipt != nil {
+		t.Fatal("failed CAS changed reservation")
+	}
+	if err := c.ObserveOperatorStep(ctx, key, receipt); err != nil {
+		t.Fatal(err)
+	}
+	if err := c.ObserveOperatorStep(ctx, key, receipt); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("observation replay: %v", err)
+	}
+	got, err := c.OperatorStep(ctx, key)
+	if err != nil || got.Phase != "observed" || got.Version != 2 || !bytes.Equal(got.Receipt, receipt) {
+		t.Fatalf("observed receipt: %+v %v", got, err)
+	}
+	key.Version++
+	if err := c.ObserveOperatorStep(ctx, key, []byte("second receipt")); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("observed phase accepted at current version: %v", err)
+	}
+	if after, err := c.OperatorStep(ctx, key); err != nil || after.Version != got.Version || !bytes.Equal(after.Receipt, got.Receipt) {
+		t.Fatalf("observed replay changed its receipt: %+v %v", after, err)
+	}
+	operatorCounts(t, c, 1, 2)
+}
+
+func TestOperatorStepGuards(t *testing.T) {
+	c := openTestCredentials(t)
+	ctx := testCtx(t)
+	cas, runner := operatorClaimFixture(t, c)
+	if err := c.ClaimOperatorAction(ctx, cas, runner, 1599); err != nil {
+		t.Fatal(err)
+	}
+	cas.Generation++
+	key := OperatorStepKey{Action: cas, Owner: runner.Nonce, Kind: "checks", Identity: "guard", Version: 1}
+	checks := []struct {
+		name             string
+		edit             func(*OperatorStepKey)
+		reserve, observe error
+	}{
+		{"id", func(k *OperatorStepKey) { k.Action.ID += "-other" }, ErrOperatorConflict, ErrOperatorConflict},
+		{"digest", func(k *OperatorStepKey) { k.Action.Digest[0] ^= 1 }, ErrOperatorConflict, ErrOperatorConflict},
+		{"generation", func(k *OperatorStepKey) { k.Action.Generation++ }, ErrOperatorConflict, ErrOperatorConflict},
+		{"owner", func(k *OperatorStepKey) { k.Owner += "-other" }, ErrOperatorConflict, ErrOperatorConflict},
+		{"empty_owner", func(k *OperatorStepKey) { k.Owner = "" }, ErrOperatorInput, ErrOperatorConflict},
+		{"operation", func(k *OperatorStepKey) { k.Kind = "operation" }, ErrOperatorInput, ErrOperatorConflict},
+		{"empty_identity", func(k *OperatorStepKey) { k.Identity = "" }, ErrOperatorInput, ErrOperatorConflict},
+		{"version", func(k *OperatorStepKey) { k.Version++ }, ErrOperatorInput, ErrOperatorConflict},
+	}
+	for _, tc := range checks {
+		t.Run("reserve_"+tc.name, func(t *testing.T) {
+			bad := key
+			tc.edit(&bad)
+			if err := c.ReserveOperatorStep(ctx, bad, 1599); !errors.Is(err, tc.reserve) {
+				t.Fatalf("reservation class: %v", err)
+			}
+			operatorCounts(t, c, 1, 1)
+		})
+	}
+	if err := c.ReserveOperatorStep(ctx, key, 999); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("clock before creation: %v", err)
+	}
+	operatorExec(t, c, "UPDATE operator_actions SET state='uncertain'")
+	if err := c.ReserveOperatorStep(ctx, key, 1599); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("uncertain action: %v", err)
+	}
+	operatorCounts(t, c, 1, 1)
+	operatorExec(t, c, "UPDATE operator_actions SET state='running'")
+	if err := c.ReserveOperatorStep(ctx, key, 1599); err != nil {
+		t.Fatal(err)
+	}
+	for _, tc := range checks {
+		t.Run("observe_"+tc.name, func(t *testing.T) {
+			bad := key
+			tc.edit(&bad)
+			if err := c.ObserveOperatorStep(ctx, bad, []byte("readback")); !errors.Is(err, tc.observe) {
+				t.Fatalf("observation class: %v", err)
+			}
+			assertReservedOperatorStep(t, c, key)
+		})
+	}
+	for _, body := range [][]byte{nil, make([]byte, MaxSurfaceApprovalArtifactBytes+1)} {
+		if err := c.ObserveOperatorStep(ctx, key, body); !errors.Is(err, ErrOperatorInput) {
+			t.Fatalf("receipt size %d: %v", len(body), err)
+		}
+		assertReservedOperatorStep(t, c, key)
+	}
+	operatorExec(t, c, "UPDATE operator_action_steps SET owner='replacement' WHERE kind='checks'")
+	if err := c.ObserveOperatorStep(ctx, key, []byte("readback")); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("stale step owner: %v", err)
+	}
+	replaced := key
+	replaced.Owner = "replacement"
+	assertReservedOperatorStep(t, c, replaced)
+	operatorExec(t, c, "UPDATE operator_action_steps SET owner='runner' WHERE kind='checks'")
+	operatorExec(t, c, "UPDATE operator_actions SET owner='replacement'")
+	if err := c.ObserveOperatorStep(ctx, key, []byte("readback")); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("stale action owner: %v", err)
+	}
+	assertReservedOperatorStep(t, c, key)
+	operatorExec(t, c, "UPDATE operator_actions SET owner='runner'")
+	for i := 1; i <= 32; i++ {
+		next := key
+		next.Identity = fmt.Sprint(i)
+		err := c.ReserveOperatorStep(ctx, next, 1599)
+		if i < 32 && err != nil || i == 32 && !errors.Is(err, ErrOperatorConflict) {
+			t.Fatalf("additional reservation %d: %v", i, err)
+		}
+	}
+	operatorCounts(t, c, 1, 33)
+}
+
+func assertReservedOperatorStep(t *testing.T, c *Credentials, key OperatorStepKey) {
+	t.Helper()
+	got, err := c.OperatorStep(testCtx(t), key)
+	if err != nil || got.Owner != key.Owner || got.Version != key.Version || got.Phase != "reserved" || got.Receipt != nil {
+		t.Fatalf("expected exact reserved step: %+v %v", got, err)
+	}
+}
diff --git a/host/store/operator_steps.go b/host/store/operator_steps.go
new file mode 100644
index 00000000..0430d011
--- /dev/null
+++ b/host/store/operator_steps.go
@@ -0,0 +1,72 @@
+package store
+
+import (
+	"context"
+	"database/sql"
+	"errors"
+)
+
+type OperatorStepKey struct {
+	Action                OperatorCAS
+	Owner, Kind, Identity string
+	Version               int64
+}
+
+// ReserveOperatorStep records identity before a future caller crosses an
+// effect boundary. The caller must first recheck its own frozen proof.
+func (c *Credentials) ReserveOperatorStep(ctx context.Context, key OperatorStepKey, now int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if key.Kind == "operation" || key.Identity == "" || key.Owner == "" || key.Version != 1 {
+		return ErrOperatorInput
+	}
+	result, err := c.db.ExecContext(ctx, `INSERT INTO operator_action_steps
+		(kind, identity, action_id, phase, version, owner, created_at)
+		SELECT ?, ?, id, 'reserved', ?, owner, ? FROM operator_actions
+		WHERE id = ? AND digest = ? AND generation = ? AND owner = ?
+		AND state IN ('claimed', 'running') AND expires_at > ? AND created_at <= ?
+		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind <> 'operation') < 32
+		ON CONFLICT(kind, identity) DO NOTHING`, key.Kind, key.Identity, key.Version, now,
+		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner, now, now, key.Action.ID)
+	return operatorChangedOne(result, err)
+}
+
+// ObserveOperatorStep records the outcome of a previously reserved effect.
+// Expiry prevents new reservations but does not prevent exact observation.
+func (c *Credentials) ObserveOperatorStep(ctx context.Context, key OperatorStepKey, receipt []byte) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if len(receipt) == 0 || len(receipt) > MaxSurfaceApprovalArtifactBytes {
+		return ErrOperatorInput
+	}
+	result, err := c.db.ExecContext(ctx, `UPDATE operator_action_steps
+		SET phase = 'observed', version = version + 1, receipt = ?
+		WHERE action_id = ? AND kind = ? AND identity = ? AND version = ? AND owner = ? AND phase = 'reserved'
+		AND EXISTS (SELECT 1 FROM operator_actions WHERE id = ? AND digest = ? AND generation = ? AND owner = ?)`,
+		receipt, key.Action.ID, key.Kind, key.Identity, key.Version, key.Owner,
+		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner)
+	return operatorChangedOne(result, err)
+}
+
+type OperatorStepReceipt struct {
+	Phase, Owner string
+	Version      int64
+	Receipt      []byte
+}
+
+// OperatorStep reads a reservation or recorded observation across a restart.
+func (c *Credentials) OperatorStep(ctx context.Context, key OperatorStepKey) (OperatorStepReceipt, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorStepReceipt{}, err
+	}
+	var step OperatorStepReceipt
+	err := c.db.QueryRowContext(ctx, `SELECT phase, coalesce(owner,''), version, receipt FROM operator_action_steps
+		WHERE action_id = ? AND kind = ? AND identity = ?`, key.Action.ID, key.Kind, key.Identity).
+		Scan(&step.Phase, &step.Owner, &step.Version, &step.Receipt)
+	if errors.Is(err, sql.ErrNoRows) {
+		err = ErrOperatorMissing
+	}
+	return step, err
+}
````

### Patch M2

````diff M2
diff --git a/cmd/ailang-worldd/help_test.go b/cmd/ailang-worldd/help_test.go
index 4a9eb60a..d0ae5963 100644
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -13,7 +13,7 @@ var everyVerb = [][]string{
 	{"tools", "list"}, {"call"}, {"why"}, {"provenance"}, {"setup"}, {"doctor"}, {"hosts", "init"},
 	{"session"}, {"session", "new"}, {"session", "list"}, {"session", "revoke"}, {"session", "mint"},
 	{"approve"}, {"service", "runner"}, {"service", "runner", "install"}, {"service", "runner", "start"},
-	{"service", "runner", "stop"}, {"service", "runner", "status"}, {"service", "runner", "cancel"}, {"service", "runner", "uninstall"}, {"approver"}, {"approver", "list"}, {"approver", "enroll"}, {"approver", "revoke"},
+	{"service", "runner", "stop"}, {"service", "runner", "status"}, {"service", "runner", "cancel"}, {"service", "runner", "reconcile-worktree"}, {"service", "runner", "uninstall"}, {"approver"}, {"approver", "list"}, {"approver", "enroll"}, {"approver", "revoke"},
 }
 
 // AC5.7: `<verb> --help` and `help <verb>` exit 0 with that verb's help on
@@ -42,7 +42,7 @@ func TestEveryVerbHasHelp(t *testing.T) {
 		"serve --db", "] health", "] head", "world get", "object get", "object find", "log get", "log range",
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
-		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
+		"service runner install|start|stop|status|cancel|reconcile-mint|reconcile-worktree|uninstall", "session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
 		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve <request-id> --webauthn-link", "approve <request-id> --webauthn-assertion <file>", "approve <request-id> --await [--await-for <duration>]", "approve --list [--decided]", "approve <request-ref> [--deny] [--daemon <host:port>]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
diff --git a/cmd/ailang-worldd/main.go b/cmd/ailang-worldd/main.go
index a533640c..99385961 100644
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -83,7 +83,7 @@ Usage:
   ailang-worldd hosts init [--out-dir <dir>] [--addr <url>] [--replace]
   ailang-worldd version [--json]
   ailang-worldd service install|status|upgrade|stop|start|uninstall [flags; upgrade: --binary FILE --pause]
-  ailang-worldd service runner install|start|stop|status|cancel|uninstall
+  ailang-worldd service runner install|start|stop|status|cancel|reconcile-mint|reconcile-worktree|uninstall
   ailang-worldd serve --db <path> [--credentials-db <path>] [--bind host:port] [--ailang-bin <path>]
                       [--workbench-still]
                       [--workspace-root <dir> --tool-ailang-bin <path>]
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index d837ef24..fd2978d0 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -21,8 +21,8 @@ import (
 
 const runnerLabel = "com.sunholo.ailang-world.operator-runner"
 const runnerSystemPath = "/usr/bin:/bin:/usr/sbin:/sbin"
-const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | uninstall
-Install, stop, cancel, orphan mint revocation and uninstall require typed y. Start uses the approved configuration.
+const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | uninstall
+Install, stop, cancel, orphan mint revocation, worktree disposition and uninstall require typed y. Start uses the approved configuration.
 Config JSON keys: home, uid, address, world, credentials, executables (absolute path to sha256 digest), preset_names, tool_dirs, backend ("headless").
 The World service must already be running. Use absolute, canonical paths and include its installed binary in executables.
 `
@@ -272,7 +272,7 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 		if len(args) != 3 || args[1] != "--config" {
 			return refuse(errors.New("expected install --config FILE"))
 		}
-	case "cancel", "reconcile-mint":
+	case "cancel", "reconcile-mint", "reconcile-worktree":
 		if len(args) != 2 || !verbs.ValidOperatorActionID(args[1]) {
 			return refuse(fmt.Errorf("expected %s ACTION-ID", verb))
 		}
@@ -286,6 +286,9 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 	if verb == "reconcile-mint" {
 		return e.reconcileRunnerMint(args[1], out, errOut)
 	}
+	if verb == "reconcile-worktree" {
+		return e.reconcileRunnerWorktree(args[1], out, errOut)
+	}
 	p := e.runnerPaths()
 	if verb == "status" {
 		if err := e.runnerStatus(out); err != nil {
diff --git a/cmd/ailang-worldd/operator_worktree.go b/cmd/ailang-worldd/operator_worktree.go
new file mode 100644
index 00000000..f579fed6
--- /dev/null
+++ b/cmd/ailang-worldd/operator_worktree.go
@@ -0,0 +1,261 @@
+package main
+
+import (
+	"context"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"os"
+	"path/filepath"
+	"strings"
+	"syscall"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+	"github.com/sunholo-data/ailang-world/host/worktree"
+)
+
+type operatorWorktreeTarget struct {
+	Action, Repo, Common, Path, Branch, SHA string
+}
+
+type operatorStartSource struct {
+	presets, profiles, workspace, hosts string
+	clock                               func() int64
+}
+
+// reserveOperatorStartWorktree is an isolated preparation seam. The
+// paired-process consumer must integrate its held episode fd before dispatch.
+func reserveOperatorStartWorktree(ctx context.Context, credentials string, signed approval.OperatorAction,
+	attempt store.OperatorMintAttempt, src operatorStartSource, checkpoint func(string) error) error {
+	canonical, err := store.CanonicalPath(credentials)
+	if err != nil {
+		return err
+	}
+	if signed.Store != canonical || !verbs.EpisodeIDPattern.MatchString(signed.Episode) {
+		return store.ErrOperatorConflict
+	}
+	lock, err := openRunnerLock(operatorEpisodeLockPath(canonical, signed.Episode), 0)
+	if err != nil {
+		return fmt.Errorf("%w: episode lock unavailable", store.ErrOperatorConflict)
+	}
+	defer lock.Close()
+	in, err := reconstructOperatorStart(ctx, signed, attempt, src.presets, src.profiles, src.workspace, src.hosts, src.clock)
+	if err != nil {
+		return err
+	}
+	git, err := sessionGit()
+	if err != nil {
+		return err
+	}
+	target, err := newOperatorWorktreeTarget(ctx, git, signed, in.Start)
+	if err != nil {
+		return err
+	}
+	initial, err := target.inspect(ctx, git)
+	if err != nil {
+		return err
+	}
+	if initial.state != treeAbsent {
+		return fmt.Errorf("%w: worktree path, branch or registration already exists", store.ErrOperatorConflict)
+	}
+	if err := checkpoint("before reservation"); err != nil {
+		return err
+	}
+	identity, _ := json.Marshal(target) // string fields only
+	var key store.OperatorStepKey
+	err = runnerWrite(canonical, func(db *store.Credentials) error {
+		row, err := db.OperatorAction(ctx, signed.ID, src.clock())
+		if err != nil {
+			return err
+		}
+		if len(row.Challenge) != 32 {
+			return approval.ErrNotVerified
+		}
+		proof, err := authority.RecheckClaimedOperatorAction(ctx, db, signed, [32]byte(row.Challenge), row.ExpiresAt, src.clock(), row)
+		if err != nil {
+			return err
+		}
+		key, err = db.ReserveOperatorWorktree(ctx, in.Attempt, string(identity), proof.OperatorConfirmation().MintProof(), src.clock)
+		return err
+	})
+	if err != nil {
+		return err
+	}
+	if err := checkpoint("before add"); err != nil {
+		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
+	}
+	addErr := worktree.Add(ctx, git, target.Repo, target.Path, target.Branch, target.SHA)
+	if err := checkpoint("before readback"); err != nil {
+		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
+	}
+	got, readErr := target.inspect(ctx, git)
+	if err := checkpoint("before observe"); err != nil {
+		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
+	}
+
+	created := addErr == nil && readErr == nil && got.state == treeExact
+	disposition := store.OperatorWorktreeUncertain
+	if addErr != nil && readErr == nil && got.state == treeAbsent {
+		disposition = store.OperatorWorktreeAbsent
+	}
+	err = runnerWrite(canonical, func(db *store.Credentials) error {
+		if created {
+			return db.ObserveOperatorStep(ctx, key, []byte("created"))
+		}
+		return db.DisposeOperatorWorktree(ctx, key, disposition, src.clock())
+	})
+	if err != nil {
+		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
+	}
+	if created {
+		return nil
+	}
+	if disposition == store.OperatorWorktreeAbsent {
+		return errors.New("worktree add failed; worktree and branch absent")
+	}
+	return fmt.Errorf("worktree outcome uncertain; %s", reconcileWorktreeCommand(key.Action.ID))
+}
+
+func operatorWorktreeFact(a approval.OperatorAction) string {
+	for _, f := range a.Resolution {
+		if f.Name == "worktree_path" {
+			return f.Value
+		}
+	}
+	return ""
+}
+
+func operatorPhysicalTarget(path string) (string, error) {
+	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
+		return "", store.ErrOperatorConflict
+	}
+	parent, err := canonicalLaunchPath(filepath.Dir(path))
+	if err != nil {
+		return "", err
+	}
+	return filepath.Join(parent, filepath.Base(path)), nil
+}
+
+func newOperatorWorktreeTarget(ctx context.Context, git string, a approval.OperatorAction, start store.EpisodeStart) (operatorWorktreeTarget, error) {
+	var target operatorWorktreeTarget
+	repo, err := filepath.EvalSymlinks(start.Repo)
+	if err != nil {
+		return target, err
+	}
+	path, err := operatorPhysicalTarget(operatorWorktreeFact(a))
+	if err != nil {
+		return target, err
+	}
+	common, err := operatorGitCommon(ctx, git, repo)
+	if err != nil {
+		return target, err
+	}
+	return operatorWorktreeTarget{a.ID, repo, common, path, start.Branch, start.Base.SHA}, nil
+}
+
+func operatorGitCommon(ctx context.Context, git, dir string) (string, error) {
+	out, err := worktree.Output(ctx, git, []string{"-C", dir, "rev-parse", "--path-format=absolute", "--git-common-dir"})
+	if err != nil {
+		return "", err
+	}
+	return filepath.EvalSymlinks(strings.TrimSuffix(out, "\n"))
+}
+
+const (
+	treeAbsent = iota
+	treePartial
+	treeExact
+)
+
+type operatorWorktreeEvidence struct {
+	state                            int
+	registration, refs, common, head string
+	mode                             os.FileMode
+	dev, inode                       uint64
+}
+
+// inspect checks the repository even when the target is absent. Path absence
+// alone cannot distinguish a prunable registration or a branch-only result.
+func (target operatorWorktreeTarget) inspect(ctx context.Context, git string) (operatorWorktreeEvidence, error) {
+	var e operatorWorktreeEvidence
+	common, err := operatorGitCommon(ctx, git, target.Repo)
+	if err != nil || common != target.Common {
+		return e, errors.New("worktree repository identity differs")
+	}
+	list, err := worktree.Output(ctx, git, []string{"-C", target.Repo, "worktree", "list", "--porcelain", "-z"})
+	if err != nil {
+		return e, err
+	}
+	ref := "refs/heads/" + target.Branch
+	want := "worktree " + target.Path + "\x00HEAD " + target.SHA + "\x00branch " + ref
+	exact := 0
+	for _, entry := range strings.Split(strings.TrimSuffix(list, "\x00\x00"), "\x00\x00") {
+		path, _, _ := strings.Cut(entry, "\x00")
+		raw, ok := strings.CutPrefix(path, "worktree ")
+		if !ok || !filepath.IsAbs(raw) {
+			return e, errors.New("invalid Git worktree registration")
+		}
+		physical, err := canonicalLaunchPath(raw)
+		if err != nil {
+			return e, err
+		}
+		if physical == target.Path || raw == target.Path || strings.Contains("\x00"+entry+"\x00", "\x00branch "+ref+"\x00") {
+			e.registration += entry + "\x00\x00"
+			if "worktree "+physical+strings.TrimPrefix(entry, path) == want {
+				exact++
+			}
+		}
+	}
+	refs, err := worktree.Output(ctx, git, []string{"-C", target.Repo, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads"})
+	if err != nil {
+		return e, err
+	}
+	for _, line := range strings.Split(refs, "\n") {
+		name, _, _ := strings.Cut(line, " ")
+		if name == ref || strings.HasPrefix(name, ref+"/") || strings.HasPrefix(ref, name+"/") {
+			e.refs += line + "\n"
+		}
+	}
+	info, err := os.Lstat(target.Path)
+	if errors.Is(err, os.ErrNotExist) {
+		if e.registration != "" || e.refs != "" {
+			e.state = treePartial
+		}
+		return e, nil
+	}
+	if err != nil {
+		return e, err
+	}
+	stat := info.Sys().(*syscall.Stat_t)
+	e.mode, e.dev, e.inode = info.Mode(), uint64(stat.Dev), stat.Ino
+	e.state = treePartial
+	if !info.IsDir() || exact != 1 || strings.Count(e.registration, "\x00\x00") != 1 || e.refs != ref+" "+target.SHA+"\n" {
+		return e, nil
+	}
+	physical, err := filepath.EvalSymlinks(target.Path)
+	if err != nil || physical != target.Path {
+		return e, errors.New("worktree physical path differs")
+	}
+	link, err := os.Lstat(filepath.Join(target.Path, ".git"))
+	if errors.Is(err, os.ErrNotExist) {
+		return e, nil
+	}
+	if err != nil {
+		return e, err
+	}
+	if !link.Mode().IsRegular() || !worktree.IsLinked(target.Path) {
+		return e, nil
+	}
+	e.common, err = operatorGitCommon(ctx, git, target.Path)
+	if err != nil || e.common != target.Common {
+		return e, nil
+	}
+	e.head, err = worktree.Output(ctx, git, []string{"-C", target.Path, "rev-parse", "HEAD", "--symbolic-full-name", "HEAD"})
+	if err == nil && e.head == target.SHA+"\n"+ref+"\n" {
+		e.state = treeExact
+	}
+	return e, nil
+}
diff --git a/cmd/ailang-worldd/operator_worktree_reconcile.go b/cmd/ailang-worldd/operator_worktree_reconcile.go
new file mode 100644
index 00000000..eadf23c3
--- /dev/null
+++ b/cmd/ailang-worldd/operator_worktree_reconcile.go
@@ -0,0 +1,110 @@
+package main
+
+import (
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"path/filepath"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+func reconcileWorktreeCommand(id string) string {
+	return "ailang-worldd service runner reconcile-worktree " + id
+}
+
+func readRunnerWorktree(credentials, id string) (store.OperatorStepKey, operatorWorktreeTarget, string, error) {
+	var target operatorWorktreeTarget
+	db, err := store.OpenCredentialsReadOnly(runnerDB(credentials))
+	if err != nil {
+		return store.OperatorStepKey{}, target, "", err
+	}
+	defer db.Close()
+	ctx, cancel := budgetContext(time.Second)
+	defer cancel()
+	key, signed, err := db.ReadOperatorWorktree(ctx, id)
+	if err != nil {
+		return key, target, "", err
+	}
+	if err := json.Unmarshal([]byte(key.Identity), &target); err != nil {
+		return key, target, "", store.ErrOperatorConflict
+	}
+	path, err := operatorPhysicalTarget(operatorWorktreeFact(signed))
+	if err != nil {
+		return key, target, "", err
+	}
+	if signed.Store != credentials || signed.Kind != "start" || !signed.Options.Headless ||
+		!verbs.EpisodeIDPattern.MatchString(signed.Episode) || filepath.Base(path) != signed.Episode ||
+		target.Action != id || target.Path != path || target.Branch != "world/"+signed.Episode ||
+		!filepath.IsAbs(target.Repo) || !filepath.IsAbs(target.Common) {
+		return key, target, "", store.ErrOperatorConflict
+	}
+	return key, target, signed.Episode, nil
+}
+
+func (e serviceEnv) reconcileRunnerWorktree(id string, out, stderr io.Writer) int {
+	remedy := reconcileWorktreeCommand(id)
+	refuse := func(err error) int {
+		return serviceMessage(stderr, "worktree recovery refused", err.Error(), remedy, exitUsage)
+	}
+	c, _, err := e.installedRunner()
+	if err != nil {
+		return refuse(err)
+	}
+	key, target, episode, err := readRunnerWorktree(c.Credentials, id)
+	if err != nil {
+		return refuse(err)
+	}
+	lock, err := openRunnerLock(operatorEpisodeLockPath(c.Credentials, episode), 0)
+	if err != nil {
+		return refuse(errors.New("episode recovery lock held or unavailable"))
+	}
+	defer lock.Close()
+	current, _, currentEpisode, err := readRunnerWorktree(c.Credentials, id)
+	if err != nil {
+		return refuse(err)
+	}
+	if current != key || currentEpisode != episode {
+		return refuse(store.ErrOperatorConflict)
+	}
+	git, err := sessionGit()
+	if err != nil {
+		return refuse(err)
+	}
+	ctx, cancel := budgetContext(5 * time.Second)
+	before, err := target.inspect(ctx, git)
+	cancel()
+	if err != nil {
+		return refuse(err)
+	}
+	disposition := store.OperatorWorktreeRetained
+	if before.state == treeAbsent {
+		disposition = store.OperatorWorktreeAbsent
+	}
+	if !e.confirm(stderr, fmt.Sprintf("Stop action %s (%s)? Keep its worktree, branch and files as found.", id, disposition)) {
+		fmt.Fprintln(stderr, "Retry at a terminal:", remedy)
+		return exitUsage
+	}
+	ctx, cancel = budgetContext(5 * time.Second)
+	after, err := target.inspect(ctx, git)
+	cancel()
+	if err != nil {
+		return refuse(err)
+	}
+	if before != after {
+		return refuse(errors.New("worktree evidence changed after review"))
+	}
+	err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
+		ctx, cancel := budgetContext(time.Second)
+		defer cancel()
+		return db.DisposeOperatorWorktree(ctx, key, disposition, e.terminal.now())
+	})
+	if err != nil {
+		return refuse(err)
+	}
+	fmt.Fprintln(out, "Worktree reservation disposed for", id, "("+disposition+")")
+	return exitOK
+}
diff --git a/cmd/ailang-worldd/operator_worktree_test.go b/cmd/ailang-worldd/operator_worktree_test.go
new file mode 100644
index 00000000..9557aef8
--- /dev/null
+++ b/cmd/ailang-worldd/operator_worktree_test.go
@@ -0,0 +1,557 @@
+package main
+
+import (
+	"context"
+	"database/sql"
+	"encoding/json"
+	"errors"
+	"io"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"syscall"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/worktree"
+)
+
+const worktreeExample = "ailang-worldd service runner reconcile-worktree oa-0123456789abcdef0123456789abcdef"
+
+type worktreeStartRig struct {
+	operatorStartRig
+	*operatorServiceFixture
+	db  *sql.DB
+	src operatorStartSource
+}
+
+func newWorktreeStartRig(t *testing.T, claude bool) *worktreeStartRig {
+	t.Helper()
+	s := operatorService(t)
+	s.mustInstall(t)
+	s.writes = nil
+	s.env.terminal.now = func() int64 { return 1012 }
+	f := newPresetFixture(t)
+	r := operatorStartRig{f: f, bin: startReadFakeAgent(t, f), creds: s.config.Credentials, hosts: t.TempDir()}
+	r.f.r.f.db = s.config.World
+	r.action = approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "start", Episode: "ep",
+		Store: r.creds, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: operatorStartTask, Review: "fixture"}
+	r.action.Argv = approval.OperatorArgv(r.action.Kind, r.action.Episode, r.action.Options)
+	if claude {
+		if err := os.Rename(filepath.Join(r.bin, "codex"), filepath.Join(r.bin, "claude")); err != nil {
+			t.Fatal(err)
+		}
+	}
+	r.editPreset(t, func(p *launchPreset) {
+		p.Daemon, _ = json.Marshal(presetDaemon{Addr: defaultEpisodeAddr, DB: s.config.World, CredentialsDB: r.creds,
+			WorkspaceRoot: r.f.r.f.ws, SessionsDir: r.f.r.sessions})
+		if claude {
+			p.Agent = "claude"
+		}
+	})
+	h := matchingHealth(t, r.f.r.f, r.f.r.profiles, r.f.profile)
+	h.DBPath = s.config.World
+	stubDaemon(t, h, nil)
+	r.at = store.OperatorMintAttempt{Catalogue: 1}
+	src := operatorStartSource{r.f.dir, r.f.r.profiles, r.f.r.f.ws, r.hosts, s.env.terminal.now}
+	in, facts, err := operatorStartFacts(presetTestContext(t), r.action, r.at, src.presets, src.profiles, src.workspace, src.hosts, src.clock, resolvePresetStart)
+	if err != nil {
+		t.Fatal(err)
+	}
+	r.in, r.action.Resolution = in, facts
+	held := claimFixtureEpisodeAction(t, r.creds, r.action.Episode, func(a *approval.OperatorAction) {
+		*a = r.action
+	})
+	if err := held.Close(); err != nil {
+		t.Fatal(err)
+	}
+	db, err := sql.Open("sqlite", runnerDB(r.creds))
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { db.Close() })
+	err = runnerWrite(r.creds, func(c *store.Credentials) error {
+		ctx := boundedTestContext(t)
+		row, err := c.OperatorAction(ctx, r.action.ID, src.clock())
+		if err != nil {
+			return err
+		}
+		h, err := c.OperatorHeartbeat(ctx)
+		r.at.CAS = store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
+		r.at.Runner = h.Owner
+		return err
+	})
+	if err != nil {
+		t.Fatal(err)
+	}
+	return &worktreeStartRig{r, s, db, src}
+}
+
+func (r *worktreeStartRig) start(t *testing.T, step func(string) error) error {
+	t.Helper()
+	ctx, cancel := context.WithTimeout(presetTestContext(t), 5*time.Second)
+	defer cancel()
+	return reserveOperatorStartWorktree(ctx, runnerDB(r.creds), r.action, r.at, r.src, step)
+}
+
+func (r *worktreeStartRig) exec(t *testing.T, q string, args ...any) {
+	t.Helper()
+	if _, err := r.db.ExecContext(boundedTestContext(t), q, args...); err != nil {
+		t.Fatal(err)
+	}
+}
+
+func (r *worktreeStartRig) path() string { return filepath.Join(r.f.r.f.ws, r.action.Episode) }
+
+func (r *worktreeStartRig) rows(t *testing.T, wantState, wantPhase, wantReceipt string, wantVersion int) {
+	t.Helper()
+	ctx := boundedTestContext(t)
+	var state, owner string
+	var generation, terminal, sessions, starts int
+	q := `SELECT state, coalesce(owner,''), generation, coalesce(terminal_at,0),
+		(SELECT count(*) FROM session_credentials), (SELECT count(*) FROM episode_starts) FROM operator_actions`
+	err := r.db.QueryRowContext(ctx, q).Scan(&state, &owner, &generation, &terminal, &sessions, &starts)
+	wantOwner, wantGeneration, wantTerminal := r.at.Runner.Nonce, r.at.CAS.Generation, 0
+	switch wantState {
+	case "running":
+		wantGeneration++
+	case "uncertain":
+		wantGeneration += 2
+	case "failed":
+		wantOwner, wantGeneration, wantTerminal = "", wantGeneration+int64(wantVersion), 1012
+	}
+	if err != nil || state != wantState || owner != wantOwner || int64(generation) != wantGeneration ||
+		terminal != wantTerminal || starts != 0 || sessions != 0 {
+		t.Fatalf("action=%s owner=%q generation=%d terminal=%d starts=%d sessions=%d: %v", state, owner, generation, terminal, starts, sessions, err)
+	}
+	var phase, receipt, identity, stepOwner, actionID string
+	var version, created int
+	err = r.db.QueryRowContext(ctx, `SELECT phase, coalesce(receipt,''), version, identity, owner, action_id, created_at
+		FROM operator_action_steps WHERE kind='worktree'`).Scan(&phase, &receipt, &version, &identity, &stepOwner, &actionID, &created)
+	if wantPhase == "" {
+		if !errors.Is(err, sql.ErrNoRows) {
+			t.Fatalf("unexpected step: %v", err)
+		}
+	} else {
+		var target operatorWorktreeTarget
+		path, pathErr := filepath.EvalSymlinks(filepath.Dir(r.path()))
+		repo, repoErr := filepath.EvalSymlinks(r.f.r.f.repo)
+		if err != nil || phase != wantPhase || receipt != wantReceipt || version != wantVersion ||
+			actionID != r.action.ID || stepOwner != r.at.Runner.Nonce || created != 1012 ||
+			json.Unmarshal([]byte(identity), &target) != nil || pathErr != nil || repoErr != nil ||
+			target != (operatorWorktreeTarget{r.action.ID, repo, filepath.Join(repo, ".git"), filepath.Join(path, r.action.Episode), "world/ep", r.f.base}) {
+			t.Fatalf("step=%q receipt=%q version=%d target=%+v owner=%q action=%q created=%d: %v", phase, receipt, version, target, stepOwner, actionID, created, err)
+		}
+	}
+	r.f.r.assertNoSessionFiles(t)
+	if _, err := os.Lstat(filepath.Join(r.bin, "ran")); !errors.Is(err, os.ErrNotExist) ||
+		len(r.writes) != 0 || r.ctl.boots != 0 || r.ctl.stops != 0 {
+		t.Fatalf("unexpected launch, terminal service effect or file write: %v", err)
+	}
+}
+
+func (r *worktreeStartRig) state(t *testing.T, wantState, wantPhase, wantReceipt string, wantVersion int, tree bool) {
+	t.Helper()
+	r.rows(t, wantState, wantPhase, wantReceipt, wantVersion)
+	if tree {
+		out, err := worktree.Output(boundedTestContext(t), r.f.r.f.git, []string{"-C", r.path(), "rev-parse", "HEAD", "--symbolic-full-name", "HEAD"})
+		if err != nil || strings.TrimSpace(out) != r.f.base+"\nrefs/heads/world/ep" ||
+			r.f.r.f.worktrees(t) != 2 || presetGit(t, r.f.r.f, "rev-parse", "refs/heads/world/ep") != r.f.base {
+			t.Fatalf("registered worktree and pinned ref: %q %v", out, err)
+		}
+	} else {
+		if _, err := os.Lstat(r.path()); !errors.Is(err, os.ErrNotExist) || r.f.r.f.worktrees(t) != 1 ||
+			strings.Contains(presetGit(t, r.f.r.f, "for-each-ref", "--format=%(refname)", "refs/heads"), "refs/heads/world/ep") {
+			t.Fatalf("unexpected worktree or branch: %v", err)
+		}
+	}
+}
+
+func (r *worktreeStartRig) run(t *testing.T, cmd string) (int, string) {
+	t.Helper()
+	argv := strings.Fields(cmd)
+	if len(argv) != 5 || argv[0] != "ailang-worldd" || argv[1] != "service" {
+		t.Fatalf("invalid worktree remedy %q", cmd)
+	}
+	return r.operatorServiceFixture.run(argv[2:]...)
+}
+
+func TestOperatorStartOnce(t *testing.T) {
+	r := newWorktreeStartRig(t, false)
+	if err := r.start(t, func(step string) error {
+		if step == "before add" {
+			r.state(t, "running", "reserved", "", 1, false)
+			contender, err := os.Open(operatorEpisodeLockPath(r.creds, r.action.Episode))
+			if err != nil {
+				t.Fatal(err)
+			}
+			defer contender.Close()
+			if err := syscall.Flock(int(contender.Fd()), syscall.LOCK_SH|syscall.LOCK_NB); !errors.Is(err, syscall.EWOULDBLOCK) {
+				t.Fatalf("episode lock not held across reservation: %v", err)
+			}
+		}
+		return nil
+	}); err != nil {
+		t.Fatal(err)
+	}
+	r.state(t, "running", "observed", "created", 2, true)
+	before := observeProvision(t, r.f.r, r.creds)
+	err := r.start(t, func(string) error { t.Fatal("replay crossed an effect boundary"); return nil })
+	if err == nil || !strings.Contains(err.Error(), "branch already exists") || !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+		t.Fatalf("replay changed state or wrong refusal: %v", err)
+	}
+}
+
+func TestOperatorStartCollisions(t *testing.T) {
+	for _, tc := range []struct{ name, why string }{
+		{"branch", "branch already exists"}, {"registration", "registration already exists"},
+		{"registration_foreign", "registration already exists"}, {"ref_ancestor", "registration already exists"}, {"ref_descendant", "registration already exists"},
+		{"git_initial_read", "fixture worktree list refusal"},
+		{"directory", "not a git worktree"}, {"symlink", "not a directory"},
+		{"token", "ep.session already exists"}, {"task", "ep.task.md already exists"}, {"mcp", "ep.mcp.json already exists"},
+		{"episode", "operator action conflicts"}, {"credential", "operator action conflicts"},
+		{"stale", "mint_start"}, {"catalogue", "operator action conflicts"},
+		{"runner", "operator action conflicts"}, {"revoked", "revoked"}, {"expires_at_reservation", "operator action expired"},
+		{"heartbeat_stale", "operator action conflicts"}, {"heartbeat_future", "operator action conflicts"},
+		{"heartbeat_protocol", "operator action conflicts"}, {"heartbeat_readiness", "operator action conflicts"},
+		{"held", "episode lock unavailable"}, {"missing_lock", "episode lock unavailable"},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := newWorktreeStartRig(t, tc.name == "mcp")
+			ctx := boundedTestContext(t)
+			lockPath := operatorEpisodeLockPath(r.creds, r.action.Episode)
+			switch tc.name {
+			case "branch":
+				presetGit(t, r.f.r.f, "branch", "world/ep", r.f.base)
+			case "registration_foreign":
+				other := filepath.Join(t.TempDir(), "other")
+				presetGit(t, r.f.r.f, "worktree", "add", "-b", "world/ep", other, r.f.base)
+				presetGit(t, r.f.r.f, "update-ref", "-d", "refs/heads/world/ep")
+			case "ref_ancestor", "ref_descendant":
+				branch := map[string]string{"ref_ancestor": "world", "ref_descendant": "world/ep/sub"}[tc.name]
+				presetGit(t, r.f.r.f, "branch", branch, r.f.base)
+			case "git_initial_read":
+				git := filepath.Join(r.bin, "git")
+				if err := os.Remove(git); err != nil {
+					t.Fatal(err)
+				}
+				script := "#!/bin/sh\nif [ \"$3 $4 $5 $6\" = 'worktree list --porcelain -z' ]; then echo 'fixture worktree list refusal' >&2; exit 1; fi\nexec \"" + r.f.r.f.git + "\" \"$@\"\n"
+				if err := os.WriteFile(git, []byte(script), 0700); err != nil {
+					t.Fatal(err)
+				}
+			case "registration":
+				presetGit(t, r.f.r.f, "worktree", "add", "--detach", r.path(), r.f.base)
+				if err := os.RemoveAll(r.path()); err != nil {
+					t.Fatal(err)
+				}
+			case "directory":
+				if err := os.Mkdir(r.path(), 0700); err != nil {
+					t.Fatal(err)
+				}
+			case "symlink":
+				if err := os.Symlink(t.TempDir(), r.path()); err != nil {
+					t.Fatal(err)
+				}
+			case "token", "task", "mcp":
+				if err := os.MkdirAll(r.f.r.sessions, 0700); err != nil {
+					t.Fatal(err)
+				}
+				suffix := map[string]string{"token": ".session", "task": ".task.md", "mcp": ".mcp.json"}[tc.name]
+				seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, "ep"+suffix), []byte(r.action.Task))
+			case "credential", "episode":
+				err := runnerWrite(r.creds, func(c *store.Credentials) error {
+					return c.MintSessionWithStart(ctx, store.SessionRow{CredentialID: serviceHash([]byte(r.action.Task)), EpisodeID: "ep",
+						CreatedAt: 1010, ExpiresAt: 1600, GrantsJSON: "[]"}, &r.in.Start)
+				})
+				if err != nil {
+					t.Fatal(err)
+				}
+				if tc.name == "episode" {
+					r.exec(t, "DELETE FROM session_credentials")
+				} else {
+					r.exec(t, "DELETE FROM episode_starts")
+				}
+			case "stale":
+				r.editPreset(t, func(p *launchPreset) { p.TTLS-- })
+			case "catalogue":
+				r.exec(t, "UPDATE operator_runner_state SET catalogue_generation=catalogue_generation+1")
+			case "runner":
+				r.at.Runner.Generation++
+			case "expires_at_reservation":
+				calls := 0
+				r.src.clock = func() int64 {
+					calls++
+					if calls >= 4 {
+						return 1600
+					}
+					return 1012
+				}
+			case "heartbeat_stale", "heartbeat_future", "heartbeat_protocol", "heartbeat_readiness":
+				q := map[string]string{"heartbeat_stale": "seen_at=1002", "heartbeat_future": "seen_at=1013",
+					"heartbeat_protocol": "protocol=2", "heartbeat_readiness": "readiness='finish'"}
+				r.exec(t, "UPDATE operator_runner_state SET "+q[tc.name])
+			case "revoked":
+				if err := runnerWrite(r.creds, func(c *store.Credentials) error {
+					return c.RevokeApprover(ctx, "fixture", 1011)
+				}); err != nil {
+					t.Fatal(err)
+				}
+			case "held":
+				r.editPreset(t, func(p *launchPreset) { p.TTLS-- })
+				held, err := openRunnerLock(lockPath, 0)
+				if err != nil {
+					t.Fatal(err)
+				}
+				defer held.Close()
+				defer time.AfterFunc(3*time.Second, func() { held.Close() }).Stop()
+			case "missing_lock":
+				if err := os.Remove(lockPath); err != nil {
+					t.Fatal(err)
+				}
+			}
+			before := observeProvision(t, r.f.r, r.creds)
+			err := r.start(t, func(string) error { return nil })
+			if err == nil || !strings.Contains(err.Error(), tc.why) || !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) ||
+				len(r.writes) != 0 || r.ctl.boots != 0 || r.ctl.stops != 0 {
+				t.Fatalf("collision %s changed state or wrong refusal: %v", tc.name, err)
+			}
+		})
+	}
+}
+
+func TestOperatorWorktreeWindow(t *testing.T) {
+	md, err := os.ReadFile("../../docs/QUICKSTART.md")
+	if err != nil || !strings.Contains(string(md), "\n"+worktreeExample+"\n") {
+		t.Fatalf("missing executable QUICKSTART example: %v", err)
+	}
+	for _, window := range []string{"before reservation", "before add", "before readback", "before observe", "git refused", "path only", "ref changed", "head changed", "readback error", "observe error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew"} {
+		t.Run(window, func(t *testing.T) {
+			r := newWorktreeStartRig(t, false)
+			initial := observeProvision(t, r.f.r, r.creds)
+			var atFailure provisionSnapshot
+			injected := errors.New("injected interruption")
+			gitPath := filepath.Join(r.bin, "git")
+			err := r.start(t, func(step string) error {
+				switch {
+				case window == "add error after checkout" && step == "before add":
+					hook := filepath.Join(r.f.r.f.repo, ".git", "hooks", "post-checkout")
+					if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0700); err != nil {
+						t.Fatal(err)
+					}
+				case window == "add error after checkout" && step == "before readback":
+					r.state(t, "running", "reserved", "", 1, true)
+					atFailure = observeProvision(t, r.f.r, r.creds)
+				case window == "git refused" && step == "before add":
+					dir := filepath.Join(r.f.r.f.repo, ".git", "refs", "heads", "world")
+					if err := os.Mkdir(dir, 0700); err != nil {
+						t.Fatal(err)
+					}
+					seedStartFixtureFile(t, filepath.Join(dir, "ep.lock"), []byte(r.action.Task))
+				case (window == "path only" || window == "ref changed" || window == "head changed" || (window == "redirected checkout" || window == "registration locked")) && step == "before readback":
+					switch window {
+					case "path only":
+						presetGit(t, r.f.r.f, "worktree", "remove", "--force", r.path())
+						presetGit(t, r.f.r.f, "branch", "-D", "world/ep")
+						if err := os.Mkdir(r.path(), 0700); err != nil {
+							t.Fatal(err)
+						}
+					case "ref changed":
+						presetGit(t, r.f.r.f, "update-ref", "refs/heads/world/ep", "feature")
+					case "registration locked":
+						presetGit(t, r.f.r.f, "worktree", "lock", r.path())
+					case "redirected checkout":
+						decoy := filepath.Join(t.TempDir(), "decoy")
+						presetGit(t, r.f.r.f, "worktree", "add", "--detach", decoy, r.f.base)
+						seedStartFixtureFile(t, filepath.Join(r.path(), ".git"), readStartFixtureFile(t, filepath.Join(decoy, ".git")))
+					case "head changed":
+						git := r.f.r.f
+						git.repo = r.path()
+						presetGit(t, git, "checkout", "--detach", "feature")
+					}
+					atFailure = observeProvision(t, r.f.r, r.creds)
+				case window == "ref read skew" && step == "before readback":
+					if err := os.Remove(gitPath); err != nil {
+						t.Fatal(err)
+					}
+					script := "#!/bin/sh\nreal=" + r.f.r.f.git + "\n" +
+						"if [ \"$3\" = for-each-ref ]; then\n" +
+						"  \"$real\" -C \"$2\" update-ref refs/heads/world/ep feature || exit 1\n" +
+						"  \"$real\" \"$@\" || exit 1\n" +
+						"  \"$real\" -C \"$2\" update-ref refs/heads/world/ep " + r.f.base + "\n" +
+						"else exec \"$real\" \"$@\"; fi\n"
+					if err := os.WriteFile(gitPath, []byte(script), 0700); err != nil {
+						t.Fatal(err)
+					}
+					atFailure = observeProvision(t, r.f.r, r.creds)
+				case window == "ref read skew" && step == "before observe":
+					if err := os.Remove(gitPath); err != nil {
+						t.Fatal(err)
+					}
+					if err := os.Symlink(r.f.r.f.git, gitPath); err != nil {
+						t.Fatal(err)
+					}
+				case window == "readback error" && step == "before readback":
+					if err := os.Rename(gitPath, gitPath+".hidden"); err != nil {
+						t.Fatal(err)
+					}
+					atFailure = observeProvision(t, r.f.r, r.creds)
+				case window == "readback error" && step == "before observe":
+					if err := os.Rename(gitPath+".hidden", gitPath); err != nil {
+						t.Fatal(err)
+					}
+				case window == "observe error" && step == "before observe":
+					r.exec(t, `CREATE TRIGGER fail_observe BEFORE UPDATE OF phase ON operator_action_steps
+						BEGIN SELECT RAISE(ABORT,'fixture observation refused'); END`)
+				case step == window:
+					return injected
+				}
+				return nil
+			})
+			switch window {
+			case "before reservation":
+				r.state(t, "claimed", "", "", 0, false)
+				if !errors.Is(err, injected) || !reflect.DeepEqual(initial, observeProvision(t, r.f.r, r.creds)) {
+					t.Fatalf("failure before reservation changed state or class: %v", err)
+				}
+				if code, out := r.run(t, worktreeExample); code != exitUsage || !strings.Contains(out, "operator action conflicts") ||
+					!reflect.DeepEqual(initial, observeProvision(t, r.f.r, r.creds)) {
+					t.Fatalf("command accepted an action without reservation: %d %s", code, out)
+				}
+				return
+			case "git refused":
+				r.state(t, "failed", "observed", "stopped: absent", 2, false)
+				if err == nil || err.Error() != "worktree add failed; worktree and branch absent" {
+					t.Fatal(err)
+				}
+				return
+			case "observe error":
+				r.state(t, "running", "reserved", "", 1, true)
+				if err == nil || !strings.Contains(err.Error(), "fixture observation refused") {
+					t.Fatal(err)
+				}
+				r.exec(t, "DROP TRIGGER fail_observe")
+			case "path only", "ref changed", "head changed", "readback error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew":
+				r.rows(t, "uncertain", "reserved", "", 2)
+				got := observeProvision(t, r.f.r, r.creds)
+				got.credentials = atFailure.credentials
+				if !reflect.DeepEqual(atFailure, got) || err == nil || !strings.Contains(err.Error(), "worktree outcome uncertain") {
+					t.Fatalf("partial result was changed or misclassified: %v", err)
+				}
+			default:
+				r.state(t, "running", "reserved", "", 1, window != "before add")
+				if !errors.Is(err, injected) {
+					t.Fatalf("lost failure class: %v", err)
+				}
+			}
+			_, remedy, ok := strings.Cut(err.Error(), "; ")
+			if !ok || remedy != worktreeExample {
+				t.Fatalf("incorrect remedy: %v", err)
+			}
+			before := observeProvision(t, r.f.r, r.creds)
+			for _, answer := range []string{"no tty", "n\n", "y\n"} {
+				r.env.terminal = testMintEnv(answer, 1012)
+				if answer == "no tty" {
+					r.env.terminal = noTerminalEnv(1012)
+				}
+				code, msg := r.run(t, remedy)
+				if answer != "y\n" {
+					_, retry, ok := strings.Cut(msg, "Retry at a terminal: ")
+					if code != exitUsage || !ok || strings.TrimSuffix(retry, "\n") != remedy ||
+						!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+						t.Fatalf("unattended recovery modified state: %d %s", code, msg)
+					}
+					remedy = strings.TrimSuffix(retry, "\n")
+					continue
+				}
+				if code != exitOK || !strings.Contains(msg, "Worktree reservation disposed for "+r.action.ID) {
+					t.Fatalf("%d %s", code, msg)
+				}
+			}
+			receipt, version := "stopped: retained", 2
+			switch window {
+			case "before add":
+				receipt = "stopped: absent"
+			case "path only", "ref changed", "head changed", "readback error", "add error after checkout", "redirected checkout", "registration locked", "ref read skew":
+				version = 3
+			}
+			r.rows(t, "failed", "observed", receipt, version)
+			got := observeProvision(t, r.f.r, r.creds)
+			got.credentials = before.credentials
+			if !reflect.DeepEqual(before, got) {
+				t.Fatal("recovery modified files, worktree registration or branch")
+			}
+			before = observeProvision(t, r.f.r, r.creds)
+			if code, msg := r.run(t, remedy); code != exitUsage || !strings.Contains(msg, "operator action conflicts") ||
+				!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+				t.Fatalf("second disposition changed the state: %d %s", code, msg)
+			}
+		})
+	}
+}
+
+func TestOperatorWorktreeDispositionCAS(t *testing.T) {
+	for _, change := range []string{"action_owner", "action_generation", "action_digest", "action_kind", "action_state", "step_identity", "step_phase", "step_owner", "step_version", "later_effect", "duplicate_step", "git_changed", "lock_held"} {
+		t.Run(change, func(t *testing.T) {
+			r := newWorktreeStartRig(t, false)
+			if err := r.start(t, func(step string) error {
+				if step == "before add" {
+					return errors.New("crash")
+				}
+				return nil
+			}); err == nil {
+				t.Fatal("expected interrupted reservation")
+			}
+			before := observeProvision(t, r.f.r, r.creds)
+			why := "operator action conflicts"
+			if change == "lock_held" {
+				fd, err := openRunnerLock(operatorEpisodeLockPath(r.creds, "ep"), 0)
+				if err != nil {
+					t.Fatal(err)
+				}
+				defer fd.Close()
+				defer time.AfterFunc(3*time.Second, func() { fd.Close() }).Stop()
+				why = "episode recovery lock held"
+			}
+			if change == "git_changed" {
+				why = "worktree evidence changed after review"
+			}
+			r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+				switch change {
+				case "action_owner", "action_generation", "action_digest", "action_kind", "action_state":
+					edit := map[string]string{"action_owner": "owner='replacement'", "action_generation": "generation=generation+1",
+						"action_digest": "digest=zeroblob(32)", "action_kind": "kind='finish'", "action_state": "state='claimed'"}
+					// Change only the action after the step update, inside the disposition transaction.
+					r.exec(t, `CREATE TRIGGER stale_action AFTER UPDATE OF phase ON operator_action_steps
+						BEGIN UPDATE operator_actions SET `+edit[change]+`; END`)
+				case "step_phase":
+					r.exec(t, "UPDATE operator_action_steps SET phase='observed',receipt=? WHERE kind='worktree'", []byte("created"))
+				case "step_identity":
+					r.exec(t, "UPDATE operator_action_steps SET identity=identity||' ' WHERE kind='worktree'")
+				case "step_owner":
+					r.exec(t, "UPDATE operator_action_steps SET owner='replacement' WHERE kind='worktree'")
+				case "step_version":
+					r.exec(t, "UPDATE operator_action_steps SET version=version+1 WHERE kind='worktree'")
+				case "later_effect":
+					r.exec(t, `INSERT INTO operator_action_steps(kind,identity,action_id,phase,version,owner,created_at)
+							SELECT 'mint', id, id, 'reserved', 1, owner, 1012 FROM operator_actions`)
+				case "duplicate_step":
+					r.exec(t, `INSERT INTO operator_action_steps SELECT kind, identity||'another', action_id, phase, version, owner, created_at, receipt
+						FROM operator_action_steps WHERE kind='worktree'`)
+				case "git_changed":
+					presetGit(t, r.f.r.f, "branch", "world/ep", r.f.base)
+				default:
+					t.Fatal("prompt opened despite lock contention")
+				}
+				before = observeProvision(t, r.f.r, r.creds)
+				return testMintEnv("y\n", 1012).openTerminal()
+			}
+			if code, msg := r.run(t, worktreeExample); code != exitUsage || !strings.Contains(msg, why) ||
+				!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+				t.Fatalf("changed evidence was accepted: %d %s", code, msg)
+			}
+		})
+	}
+}
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
index 3556272f..8aeb906b 100644
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -34,7 +34,7 @@ const serviceHelp = `usage: ailang-worldd service install|status|upgrade|stop|st
     [--examples-dir DIR] [--surface-index PATH] [--surface-follow PATH]
     [--surface-origin https://<host>.ts.net[:port]]... [--binary PATH] [--replace]
   status [--json] [--verbose]
-  runner install|start|stop|status|cancel|reconcile-mint|uninstall (see service runner --help)
+  runner install|start|stop|status|cancel|reconcile-mint|reconcile-worktree|uninstall (see service runner --help)
   upgrade [--binary FILE] [--pause]
 Install, stop and uninstall require a typed y at the controlling terminal.
 Start uses the approved plist; edited plists require confirmation.
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 03f46aad..e70616bd 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -96,6 +96,17 @@ token is retained for file reconciliation with the episode owner still
 held. A conflicting file or later effect reservation keeps the owner in
 place. Run the disposable example with `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^TestOperatorMintRecovery/revocation$' -count=1 -v`.
 
+For an owned start with an unresolved worktree reservation, run the attended
+disposition with its action ID. The disposable example uses this ID.
+
+```sh
+ailang-worldd service runner reconcile-worktree oa-0123456789abcdef0123456789abcdef
+```
+
+Review whether Git evidence is absent or retained, then type `y` to stop the
+action and release its owner. Any worktree, branch or files stay in place.
+Run the disposable example with `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^TestOperatorWorktreeWindow$' -count=1 -v`.
+
 ## Run your agents through World
 
 World launchers now tell Claude Code, Codex and pi their episode, worktree, session
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index 82aec381..7426dc3f 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -311,6 +311,9 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	errs["ReserveOperatorStep"] = c.ReserveOperatorStep(ctx, OperatorStepKey{}, 1000)
 	errs["ObserveOperatorStep"] = c.ObserveOperatorStep(ctx, OperatorStepKey{}, nil)
 	_, errs["OperatorStep"] = c.OperatorStep(ctx, OperatorStepKey{})
+	_, errs["ReserveOperatorWorktree"] = c.ReserveOperatorWorktree(ctx, OperatorMintAttempt{}, "", nil, func() int64 { return 1000 })
+	_, _, errs["ReadOperatorWorktree"] = c.ReadOperatorWorktree(ctx, "x")
+	errs["DisposeOperatorWorktree"] = c.DisposeOperatorWorktree(ctx, OperatorStepKey{}, OperatorWorktreeAbsent, 1000)
 	errs["CommitOperatorMint"] = c.CommitOperatorMint(ctx, OperatorMintInput{}, nil, nil)
 	_, errs["ReadOperatorMint"] = c.ReadOperatorMint(ctx, "x")
 	errs["AbortOperatorMint"] = c.AbortOperatorMint(ctx, OperatorMintReceipt{}, 1000)
diff --git a/host/store/operator_steps.go b/host/store/operator_steps.go
index 0430d011..2652788e 100644
--- a/host/store/operator_steps.go
+++ b/host/store/operator_steps.go
@@ -21,7 +21,15 @@ func (c *Credentials) ReserveOperatorStep(ctx context.Context, key OperatorStepK
 	if key.Kind == "operation" || key.Identity == "" || key.Owner == "" || key.Version != 1 {
 		return ErrOperatorInput
 	}
-	result, err := c.db.ExecContext(ctx, `INSERT INTO operator_action_steps
+	return reserveOperatorStep(ctx, c.db, key, now)
+}
+
+type operatorStepWriter interface {
+	ExecContext(context.Context, string, ...any) (sql.Result, error)
+}
+
+func reserveOperatorStep(ctx context.Context, db operatorStepWriter, key OperatorStepKey, now int64) error {
+	result, err := db.ExecContext(ctx, `INSERT INTO operator_action_steps
 		(kind, identity, action_id, phase, version, owner, created_at)
 		SELECT ?, ?, id, 'reserved', ?, owner, ? FROM operator_actions
 		WHERE id = ? AND digest = ? AND generation = ? AND owner = ?
@@ -41,7 +49,11 @@ func (c *Credentials) ObserveOperatorStep(ctx context.Context, key OperatorStepK
 	if len(receipt) == 0 || len(receipt) > MaxSurfaceApprovalArtifactBytes {
 		return ErrOperatorInput
 	}
-	result, err := c.db.ExecContext(ctx, `UPDATE operator_action_steps
+	return observeOperatorStep(ctx, c.db, key, receipt)
+}
+
+func observeOperatorStep(ctx context.Context, db operatorStepWriter, key OperatorStepKey, receipt []byte) error {
+	result, err := db.ExecContext(ctx, `UPDATE operator_action_steps
 		SET phase = 'observed', version = version + 1, receipt = ?
 		WHERE action_id = ? AND kind = ? AND identity = ? AND version = ? AND owner = ? AND phase = 'reserved'
 		AND EXISTS (SELECT 1 FROM operator_actions WHERE id = ? AND digest = ? AND generation = ? AND owner = ?)`,
diff --git a/host/store/operator_worktree.go b/host/store/operator_worktree.go
new file mode 100644
index 00000000..39f76344
--- /dev/null
+++ b/host/store/operator_worktree.go
@@ -0,0 +1,146 @@
+package store
+
+import (
+	"context"
+	"database/sql"
+	"errors"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// ReserveOperatorWorktree admits the first effect of a claimed start. The
+// caller holds episode exclusion and has reconstructed the reviewed inputs.
+func (c *Credentials) ReserveOperatorWorktree(ctx context.Context, at OperatorMintAttempt, identity string, proof *approval.OperatorMintProof, clock func() int64) (OperatorStepKey, error) {
+	var key OperatorStepKey
+	if err := c.operatorAccess(ctx); err != nil {
+		return key, err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return key, err
+	}
+	defer tx.Rollback()
+	now := clock()
+	row, err := operatorCASRow(ctx, tx, at.CAS)
+	if err != nil {
+		return key, err
+	}
+	if now >= row.ExpiresAt {
+		return key, ErrOperatorExpired
+	}
+	if row.State != "claimed" || row.Owner != at.Runner.Nonce || len(row.Challenge) != 32 || identity == "" {
+		return key, ErrOperatorConflict
+	}
+	action, err := decodeOperatorMintAction(row)
+	if err != nil {
+		return key, err
+	}
+	if c.lock == nil || action.Store != c.lock.dbPath || action.Kind != "start" || !action.Options.Headless {
+		return key, ErrOperatorConflict
+	}
+	b, err := proof.Bind(action, [32]byte(row.Challenge), row.ExpiresAt, row.Artifact)
+	if err != nil {
+		return key, err
+	}
+	if row.ApprovedBy != b.Format+":"+b.Name {
+		return key, ErrOperatorConflict
+	}
+	if err := operatorMintApprover(ctx, tx, proof, b, now); err != nil {
+		return key, err
+	}
+	res, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state = 'running', generation = generation + 1
+		WHERE id = ? AND digest = ? AND generation = ? AND owner = ? AND state = 'claimed'
+		AND EXISTS (SELECT 1 FROM operator_runner_state WHERE generation = ? AND owner = ? AND catalogue_generation = ?
+			AND protocol = ? AND seen_at > ? AND seen_at <= ? AND readiness IN ('start','both'))
+		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind <> 'operation')
+		AND NOT EXISTS (SELECT 1 FROM episode_starts WHERE episode = ?)
+		AND NOT EXISTS (SELECT 1 FROM session_credentials WHERE episode_id = ?)`,
+		row.ID, at.CAS.Digest[:], row.Generation, row.Owner, at.Runner.Generation, at.Runner.Nonce,
+		at.Catalogue, OperatorRunnerProtocol, now-10, now, row.ID, row.Episode, row.Episode)
+	if err := operatorChangedOne(res, err); err != nil {
+		return key, err
+	}
+	key = OperatorStepKey{Action: at.CAS, Owner: row.Owner, Kind: "worktree", Identity: identity, Version: 1}
+	key.Action.Generation++
+	if err := reserveOperatorStep(ctx, tx, key, now); err != nil {
+		return OperatorStepKey{}, err
+	}
+	return key, tx.Commit()
+}
+
+// ReadOperatorWorktree returns only an unresolved worktree reservation. The
+// CLI verifies its effect identity and re-reads it under the episode lock.
+func (c *Credentials) ReadOperatorWorktree(ctx context.Context, id string) (OperatorStepKey, approval.OperatorAction, error) {
+	var key OperatorStepKey
+	var action approval.OperatorAction
+	row, err := c.OperatorAction(ctx, id, 0)
+	if err != nil {
+		return key, action, err
+	}
+	if row.State != "running" && row.State != "uncertain" {
+		return key, action, ErrOperatorConflict
+	}
+	action, err = decodeOperatorMintAction(row)
+	if err != nil {
+		return key, action, err
+	}
+	key = OperatorStepKey{Action: OperatorCAS{id, [32]byte(row.Digest), row.Generation}, Owner: row.Owner, Kind: "worktree"}
+	err = c.db.QueryRowContext(ctx, `SELECT identity, version FROM operator_action_steps
+		WHERE action_id = ? AND kind = 'worktree' AND owner = ? AND phase = 'reserved'
+		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree'))
+		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'worktree') = 1`,
+		id, row.Owner, id, id).Scan(&key.Identity, &key.Version)
+	if errors.Is(err, sql.ErrNoRows) {
+		err = ErrOperatorMissing
+	}
+	return key, action, err
+}
+
+const (
+	OperatorWorktreeUncertain = ""
+	OperatorWorktreeAbsent    = "stopped: absent"
+	OperatorWorktreeRetained  = "stopped: retained"
+)
+
+// DisposeOperatorWorktree checks the current action and the exact reserved
+// step independently. Uncertainty retains ownership and leaves the step
+// reserved at its new version. A stop leaves permanent evidence.
+func (c *Credentials) DisposeOperatorWorktree(ctx context.Context, key OperatorStepKey, disposition string, now int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if key.Kind != "worktree" {
+		return ErrOperatorInput
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	switch disposition {
+	case OperatorWorktreeUncertain:
+		res, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET version = version + 1
+			WHERE action_id = ? AND kind = 'worktree' AND identity = ? AND owner = ? AND version = ? AND phase = 'reserved'`,
+			key.Action.ID, key.Identity, key.Owner, key.Version)
+		if err := operatorChangedOne(res, err); err != nil {
+			return err
+		}
+	case OperatorWorktreeAbsent, OperatorWorktreeRetained:
+		if err := observeOperatorStep(ctx, tx, key, []byte(disposition)); err != nil {
+			return err
+		}
+	default:
+		return ErrOperatorInput
+	}
+	res, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state = CASE WHEN ? THEN 'uncertain' ELSE 'failed' END,
+		owner = CASE WHEN ? THEN owner END, terminal_at = CASE WHEN NOT ? THEN ? END, generation = generation + 1
+		WHERE id = ? AND digest = ? AND generation = ? AND owner = ? AND kind = 'start' AND state IN ('running','uncertain')
+		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree'))
+		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'worktree') = 1`,
+		disposition == OperatorWorktreeUncertain, disposition == OperatorWorktreeUncertain, disposition == OperatorWorktreeUncertain, now,
+		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner, key.Action.ID, key.Action.ID)
+	if err := operatorChangedOne(res, err); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
````

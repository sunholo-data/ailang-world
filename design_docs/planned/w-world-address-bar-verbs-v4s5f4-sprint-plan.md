# Row 243 V4 S5f4 — activation through prepared-awaiting-S6

Iteration 442. Planner and executor lane `codex:gpt-6.2-sol`. Base `ad1a482c2a8755f616505a9b2d93f3369310ae01`. Apply M1 as one cumulative milestone.

## 0. Binding constraints and cut

**Outcome A.** The compiled rebase of iteration 441's prototype measured **+910/−7 = 917**. The activation unit in Appendix P measures **1391 changed lines** with filter `cmd host docs app`, including the attended gap and orphan remedies. The 1,400-line ordinary cap applies. The exit condition is durable `prepared-awaiting-S6`. The D-WORLD-119=A and D-WORLD-120=A rulings are in `82f05004`. The V4 r1 design and its r2 PROCEED 3/3 quorum govern this implementation. The five S5f3 revision files are accounted for in §1.4. The controller still runs the outside acceptance in §4.1.

`$W` is `/Users/voightkampff/dev/sunholo-data/.wt-world-iter442-row243-v4s5f4`, `$STATE` is `$HOME/.ailang/state/world-iter442`, and `$S` is `$STATE/plan-scratch`. Scratch test trees use `git archive` and read-only Git inspection. Every Git worktree created by tests resides inside `t.TempDir()`. All gated Go invocations pin `AILANG_BIN`, the SSH signer, Node and the tool interpreter. Use fake HOME, supervisor and service files in scratch or `t.TempDir()`, real migrations, repository test signer, and real Git in temporary repositories. Compare physical paths after resolving both names. Any scope proof uses before/after hash, stat or absence as observations. The test harness never seeds a real service path.

## 1. Pre-flight

### 1.1 Baseline at ad1a482c

Before the plan was written, `git status --short --branch` showed only `## sprint/world-iter442-row243-v4s5f4`. `git rev-parse HEAD` returned the base above. `git archive ad1a482c | tar -x -C "$S/base-tree"` produced the pristine sandbox baseline. The first gate started there before product edits. `$S/gate-base/legs.txt` records sandbox legs. The controller baseline under `$STATE/base/` is read-only. Package conclusions use complete text fail sets. The marker in `$STATE/baseline.done` must end with the exact base hash and `ALL_DONE` prefix before those results can authorize attribution. Section 4.1 gives the commands for remaining coverage.

The pristine gofmt, build, native vet, Linux vet, eleven inventories, GC probe, AILANG verifier and caller census recorded rc 0. Base acceptance printed `SKIPPED` because the activation pins are added by M1. Sandbox CLI plain and race failed to bind at the start of the whole package. The following package statuses reflect the raw records, without attributing a truncated or denied result. Controller baseline marker observed complete.

| Leg | Pristine sandbox | Controller pristine | Staged sandbox |
|---|---|---|---|
| `store-plain` | empty fail set | TestConvertedStoreMethodsM5b | empty fail set |
| `store-race` | empty fail set | empty fail set | empty fail set |
| `authority-plain` | SKIPPED (unchanged) | empty fail set | SKIPPED (compatibility pin in acceptance) |
| `authority-race` | SKIPPED (unchanged) | empty fail set | SKIPPED (compatibility pin in acceptance) |
| `verifygate-plain` | TestMissionConfigDefaultFlowIntegration | TestMissionConfigDefaultFlowIntegration | empty fail set |
| `cmd-plain` | UNINFORMATIVE UNDER SANDBOX | TRUNCATED | UNINFORMATIVE UNDER SANDBOX |
| `daemon-race` | PENDING | TRUNCATED | PENDING |

The final named fast acceptance emitted `ACCEPT_OK M1 fast` with **125 PASS lines** in `$S/accept-M1-final.raw`. The first broad acceptance timed out under co-load. The original gate transcript has two completed verifygate legs because its script was edited while that rehearsal was running. The final verifier uses the frozen script and its own raw logs. The legacy worktree window selector later failed its pre-existing Git deadline arms under load. The final acceptance gate retains the borrowed-FD and file-publication selectors, and §4.1 reruns the two longer legacy windows separately at base and stage.

The sandbox named race acceptance emitted `ACCEPT_OK M1 race` with **125 PASS lines** in `$S/accept-M1-race.raw`.

The byte-frozen gate reached `GATES_DONE fail=0 uninformative=1` and wrote `ALL_DONE`. Its only uninformative legs were whole CLI plain and race. The whole verifygate leg returned 0 with a complete 406.257 s package result. The log is `$S/gate-M1-sealed/legs.txt`.

### 1.2 Inputs and coupling inventory

`/bin/bash "$S/anchors.sh" "$W" "$S/st-M1" > "$S/anchors.raw"` records base declarations and numbered before/after text for each touched file. Base names in this table are at ad1a482c. CLI files are under `cmd/ailang-worldd/`.

| Input | Command and observed output |
|---|---|
| S5f3 landing | `git log -10 --oneline ad1a482c` returned record #688 `ad1a482c`, product merge `a50b0d72`, round-three `ed135376`, round-two `81789324`, and M3/M2/M1 `211524ae`, `b56285dc`, `7b7acdf6`. `git diff --name-only a50b0d72 ad1a482c` listed seven `design_docs/` paths. |
| Internal mode and heartbeat | The runner grep in `anchors.sh` printed `main.go:390`, `operator_runner.go:42,69,80,89,94`. The base comment at :89 is `No effect dispatch`. `$S/callers-M1.raw` shows the new internal consumer at `operator_runner.go:92` and its single loop entry at :70. |
| Non-test effect census | `git grep -nE -e 'withOperatorStartCandidate\(' -e 'reserveOperatorStartWorktree\(' -e 'publishOperatorStartFiles\(' -e 'MintOperatorStart\(' ad1a482c -- cmd host ':!*_test.go'` returned four definitions at `operator_start_candidate.go:15`, `operator_worktree.go:30`, `operator_files.go:267` and `host/authority/mint.go:204`. The positive controls are the actual test calls at `operator_start_candidate_test.go:92` and `operator_mint_test.go:82`. Stage output gives exactly four calls in `operator_preparation.go:38,129,147,152`, and one caller of the consumer in the private loop. |
| Activation tests | The four missing test-name grep at base returned 1. The same command at head with `grep -rnE` found all four definitions. At base, the positive grep returned `operator_pair_upgrade_test.go:200 TestOperatorOldBuildTerminalDrain` and `operator_upgrade_test.go:14 active_child`. |
| Durable APIs | `anchors.sh` printed `host/store/operator_actions.go:74`, `operator_cas.go:119`, `operator_steps.go:17,45`, `operator_mint.go:94,242`, `operator_files.go:67,87`, and `operator_start_recovery.go:38,148`. These locate the select, claim, reserve/observe, atomic mint, file read/write, and durable classification and disposition. |
| Lock and version checks | `anchors.sh` returned `operator_runner.go:21,42`, `operator_episode_lock.go:18`, `operator_service.go:113,178,222,246,258,398`, `operator_pair_upgrade.go:27,122,134`, `service_upgrade.go:223`, `session.go:477`. CLI and daemon GC controls in `callers.sh` report `episode_gc.go:116` and `host/daemon/auto_gc.go:232`. |
| CI race bound | `git grep -nF -- '-timeout 8m' ad1a482c -- scripts/verify_go.sh` printed :350. The controller readings report 480.134 s and 480.095 s on S5f3 PR and merge for daemon and 478.231 s for CLI on the PR. Normal merge retry was 273.859 s CLI and 358.805 s daemon. D-WORLD-121 remains OPEN under the default B ruling. |

The inventory command uses `go test -v -count=1 -timeout 2m ./<package>/ -run '^(<names>)$'`. Base results and stage receipts are under `$S/gate-base/` and `$S/gate-M1/`. The package regexes are embedded in §6.

| Inventory and base declaration | Coupling and selector |
|---|---|
| `TestDaemonDependencyAllowlist`, `host/daemon/daemon_test.go:899` | Daemon import closure outside the product diff. Run `inv-daemon`. |
| `TestWorldBoundaryDependencyAllowlist`, `host/boundary/allowlist_world_test.go:758` | CLI uses existing store, approval, authority and standard-library dependencies. Run `inv-boundary`. |
| `TestProductionGoSurface`, `host/store/context_roots_test.go:241` | New Go files are inside the scan. Run `inv-store`. |
| `TestProductionContextRoots`, `host/store/context_roots_test.go:136` | The preparation loop uses the existing `budgetContext` root with 22 s and the attended path uses two 5 s bounded passes on either side of typed confirmation. Run `go test ./host/store/ -run '^TestProductionContextRoots$'`. |
| `TestCredentialsMethodsRequireADeadline`, `host/store/credentials_test.go:277` | `NextOwnedOperatorStart` and `MarkOperatorMintUncertain` use `operatorAccess(ctx)`. Run `inv-store`. |
| `TestStoreSurfaceTakesContext`, `host/store/deadline_guard_test.go:43` | Both new credentials receiver methods take `context.Context` first. `ValidateOperatorStartSession` is pure. Run `inv-store`. |
| `TestBoundaryASTWriteGuard`, `host/boundary/allowlist_world_test.go:1323` | Writes go through the existing `runnerWrite` and store transactions. Run `inv-boundary`. |
| `TestOperatorNeverDispatch`, `host/verbs/operator_policy_test.go:69` | The app retains exact refusal classes for approver, session, service, serve, setup and publish, unrelated verbs and shell. Add an explicit `service runner reconcile-start` refusal. `callers.sh` admits only the internal start consumer and four effect calls above. Run `inv-verbs`, `TestOperatorRunnerSetupRefusals/global_address`, `TestOperatorServiceFence` and the positive process test. |
| `TestCLICredentialWriterSites`, `cmd/ailang-worldd/cli_credential_writers_test.go:18` | The exact list of writable openings is preserved through the existing `runnerWrite` seam. Run `inv-cli`. |
| `TestSessionRootBudgets`, `cmd/ailang-worldd/session_budget_test.go:9` | The exact `revokeCtx, cancelRevoke` source-text check is load-bearing. Run `inv-cli` and whole CLI plain. |
| `TestSubprocessSinkGate`, `host/verifygate/subprocess_sink_gate_test.go:223` | Existing shared witness continues to use `syncBuffer`. The added `bytes.Buffer` captures an in-process writer. Run `inv-sink` and whole verifygate. |

Guard source grep `git grep -n -i -E "context.root|cancelRevoke|deadline|subprocess sink|AST|allowlist|writer sites" -- cmd/ailang-worldd host/verifygate host/store host/boundary` is archived in `$S/guard-sources.raw`. The executed inventory selectors printed these results.

| Leg | Exact named passes at base and M1 |
|---|---|
| `inv-store` | `TestProductionContextRoots` PASS/PASS, `TestProductionGoSurface` PASS/PASS, `TestStoreSurfaceTakesContext` PASS/PASS, `TestCredentialsMethodsRequireADeadline` PASS/PASS |
| `inv-boundary` | `TestWorldBoundaryDependencyAllowlist` PASS/PASS, `TestBoundaryASTWriteGuard` PASS/PASS |
| `inv-cli` | `TestCLICredentialWriterSites` PASS/PASS, `TestSessionRootBudgets` PASS/PASS, `TestEveryVerbHasHelp` PASS/PASS |
| `inv-verbs` | `TestOperatorNeverDispatch` PASS/PASS |
| `inv-daemon` | `TestDaemonDependencyAllowlist` PASS/PASS |
| `inv-sink` | `TestSubprocessSinkGate` PASS/PASS |

Whole `go test -count=1 -timeout 20m ./host/verifygate/` is recorded separately. Base is TestMissionConfigDefaultFlowIntegration and M1 is empty fail set.

The sealed-check transcript `$S/seal-checks.raw` records the final touched-file gofmt, native build, native/Linux vet, exact context/sink controls and the executed disposable quickstart selector.

The exact commit verifier rehearsal returned `VERIFYGATE_FAILSET_OK`. See `$S/commit-verifygate-rehearsal/receipt.raw` and its base, head and added fail sets.

### 1.3 Decisions

- Land one milestone so claim, preparation, recovery, remedies and their tests arrive together. One bounded tick visits at most one owned start by cursor and at most one confirmed start through the canonical candidate selection. A blocking recovery yields that tick. Each in-tick effect checks the installed config and plist, pair, physical workspace, live heartbeat and catalogue, current proof, digest and stop state. Configuration and executable UID/mode checks use existing service admission. Readback can finish after admission closes.
- Hold the kernel runner descriptor before the loop. The effect visit keeps the canonical episode lock through the last durable read. Hold lifecycle exclusion across the stop check and the claim only. The runner's lifecycle attempt is non-blocking. Attended recovery takes lifecycle exclusion first, proves kernel quiescence within the existing drain budget, then takes the episode lock without waiting. Pass both descriptors to the real Git child and preserve the old owner's receipts at every unproved cut. Audit the new preparation caller with the same test that observes the real Git child and refuses terminal revoke and GC.
- The recovered prepared state is a read-only observation of agreeing action, owner, tree, atomic start/session and file receipts. An uncovered or drifted state remains blocked. Attended dispositions require the stopped companion and both descriptors, read before and after the prompt, and transact on the exact reviewed values. The code path continues to rely on S5f3's supported-name census and canonical record for pending-drain proof.

### 1.4 Premises found false or incomplete

| Premise | Command and observed result |
|---|---|
| Old M4 patch applies at ad1a482c | `(cd "$S/rebase" && /usr/bin/patch -p1 -E --dry-run --batch < "$HOME/.ailang/state/world-iter441/plan-scratch/M4.patch")` on a fresh export printed all paths and `1 out of 1 hunks failed` for `host/store/operator_start_recovery_test.go`. Successful hunks are the positive control. Added the rejected **new test** after the round-two tests. Their bodies remain byte-for-byte. Both package compiles printed `ok ... [no tests to run]`. |
| Reports of five touched files in rounds two and three may hide a conflict | `git diff --stat 211524ae ed135376 -- cmd host` printed `operator_pair_upgrade_test.go` 86, `operator_upgrade_test.go` 5, `service_upgrade.go` 39, `session.go` 6 and `host/store/operator_start_recovery_test.go` 36. Only the single above hunk failed. Retain the `syncBuffer`, revoke context and disposition fixes from the actual merge. |
| The in-memory heartbeat is a source for the current catalogue generation | `grep -n 'h.Owner.Generation = generation' cmd/ailang-worldd/operator_runner.go` printed :96 at base. The prototype's recovery used `h.CatalogueGeneration` without a durable read. The added untouched control passes zero in the in-memory field and requires current signed catalogue, exact CAS and physical absence. Production reads the fresh store heartbeat in that window. |
| A sealed preparation result can be inferred from heartbeat or path alone | `grep -nE 'OperatorStartPrepared|ReadOperatorFiles\(' "$S/st-M1/cmd/ailang-worldd/operator_start_recovery.go"` shows the prepared branch and exact readback gates. The live-child, expiry, token identity, metadata and revocation tests observe class and state equality. |
| Native outside-equivalent package result inside the sandbox | `grep -nF 'bind: operation not permitted' "$S/gate-base/"cmd-{plain,race}.raw` identified `TestApproveDecideWriteMessageAndSignature` in both runs. Denials invalidate both package readings. The focused process fixtures themselves use fake health, fake supervisor, injected time and bounded real re-execs. |
| The timing helper has the stated old name | `ls "$HOME/.ailang/state/world-iter441/plan-scratch/measure-race.sh"` found no file. `measure-final-race.sh` and `timing-release/` were the positive controls. The new helper is `$S/measure-race.sh`. |

### 1.5 Size and CI race cost

| Stage | Increment | Cumulative additions/deletions | Filtered | Unfiltered product | STATUS / DELETED / hashes |
|---|---:|---:|---:|---:|---:|
| Rebasing only | 917 | +910/−7 | 917 | 917 | 11/0/11 |
| M1 | 1391 | +1381/−10 | 1391 | 1391 | 14/0/14 |

Use `git diff -M --numstat ad1a482c2a8755f616505a9b2d93f3369310ae01 HEAD -- cmd host docs app` at the controller commit. `$S/numstat-M1` lists each path and the scratch base-relative calculation. `$S/unfiltered-M1.raw` records the full archive comparison. All product files are under the filter. Planning records add 2147 Markdown plus 1989 JSON lines outside the filter. The entire future branch sum is 5527 if the controller commits only these records and M1.

The focused comparison uses serial `go test -race -v -count=1` legs and `GOMAXPROCS=4` for the matched samples. The total counts each top-level body once. Child time is included in the parent. Receipts are in `$S/timing-low-load/`.

| Package | Named added body | Race seconds |
|---|---|---:|
| CLI | `TestOperatorActivationDrainRefusals` | 0.40 |
| CLI | `TestOperatorRevokeStatRefusal` | 0.08 |
| CLI | `TestOperatorFirstManagedEpisodeProcess` | 6.66 |
| CLI | `TestOperatorCurrentOwnerProcessRecovery` | 3.58 |
| CLI | `TestOperatorAttendedStartDisposition` | 2.01 |
| store | `TestOperatorPartialFilesOrphanRevoke` | 0.27 |
| store | `TestNextOwnedOperatorStart` | 0.14 |
| store | `TestOperatorMarkMintUncertain` | 0.18 |

The new child arms are included in their parent times.
- `TestOperatorActivationDrainRefusals/from_binary` 0.10 s
- `TestOperatorActivationDrainRefusals/budget_expired` 0.09 s
- `TestOperatorActivationDrainRefusals/stop_canonical` 0.10 s
- `TestOperatorActivationDrainRefusals/stopped_failure` 0.11 s
- `TestOperatorCurrentOwnerProcessRecovery/worktree` 0.13 s
- `TestOperatorCurrentOwnerProcessRecovery/identity` 0.11 s
- `TestOperatorCurrentOwnerProcessRecovery/task` 0.21 s
- `TestOperatorCurrentOwnerProcessRecovery/metadata` 0.04 s
- `TestOperatorCurrentOwnerProcessRecovery/expiry` 0.10 s
- `TestOperatorCurrentOwnerProcessRecovery/revocation` 0.09 s
- `TestOperatorCurrentOwnerProcessRecovery/partial_orphan` 0.43 s

**CLI named added total 12.73 s**. **Daemon added total 0.00 s** because it has zero files and zero tests changed. The earlier prototype alone measured 10.49 s on this base.

The focused M1 samples were `timing-final` 53.11 s, `timing-matched` 54.71 s, `timing-sealed` 37.08 s, `timing-postbaseline` 27.99 s, `timing-low-load` 12.73 s, `timing-isolated` 58.04 s, `timing-final-quiet` 56.56 s. Ambient load varied, including during the attempt after the other planner jobs finished.

The measured CLI time **exceeds the 12 s budget** in the primary reading. PENDING — controller to measure: `GOMAXPROCS=4 /bin/bash "$S/measure-race.sh" "$S/timing-outside"` on a quiet rig.

| Whole package | Pristine → M1 race seconds |
|---|---|
| `host/store` | 274.349 → 157.961 (-116.388 s) |
| `host/verbs` | 1.390 → 1.446 (+0.056 s) |
| `cmd/ailang-worldd` | **UNINFORMATIVE UNDER SANDBOX**. Command in §4.1. |
| `host/daemon` | **UNMEASURED**. Command in §4.1. |

### 1.6 Portability

Run native and Linux-target vet on all touched Go packages. Darwin supervision and the supported-name census retain their S5f3 behavior. All fake supervision and fake installation fixtures live under `t.TempDir()`. The re-exec tests use a pipe barrier and deadline-capped child context, wait acknowledgement, and an `atexit_sleep_ms=0` setting confined to the child. The production loop's first-tick test closes an injected tick channel after setup and needs no asynchronous wall-clock delay.

## 2. How the executor applies a sub-milestone

Use Bash in the sprint worktree. The controller commits the two planning records separately before `commit_ms.sh`. The executor uses a read-only Git reference and `/usr/bin/patch`.

````bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter442-row243-v4s5f4
STATE=$HOME/.ailang/state/world-iter442
S=$STATE/plan-scratch
T=$STATE/executor
B=ad1a482c2a8755f616505a9b2d93f3369310ae01
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5f4-sprint-plan.md
M=M1
mkdir -p "$T"
unset SNAP_PARENT SNAP_STATUS_LIST
git -C "$W" merge-base --is-ancestor "$B" HEAD
git -C "$W" diff --quiet "$B" -- cmd host docs app
test -z "$(git -C "$W" ls-files --others --exclude-standard -- cmd host docs app)"
cp "$S"/{bounded.py,failset.py,verify-scratch.py,callers.sh,probe.sh,probe-fixture.py,operator_gc_after_test.go,changed-M1} "$T/"
for name in accept.sh gates.sh snapshot.sh; do
  awk -v n="$name" '$0 == "````bash " n {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$name"
  test -s "$T/$name"
  /bin/bash -n "$T/$name"
done
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
(cd "$W" && /usr/bin/patch -p1 -E < "$T/$M.patch")
/bin/bash "$T/accept.sh" "$W" "$M" fast
/bin/bash "$T/snapshot.sh" "$W" "$M"
````

If acceptance returns 77 on actual bind denial, leave the exact staged tree for the outside controller acceptance and snapshot. A real failing named test stops the milestone.

## 3. Milestones

### M1 — Activate one installed start through durable preparation

The private `operator-runner` mode enters `runOperatorPreparationLoop` after installed build/config verification and kernel lock acquisition. It publishes the existing heartbeat and starts a bounded 22 s context. `consumeOperatorPreparation` inspects one owned recovery row and visits at most one confirmed start. Under the canonical episode lock it re-reads the candidate, checks the proof and configuration, takes lifecycle exclusion and claims by CAS. The worktree helper uses its existing reservation and exact Git readback with both descriptors. One atomic start/session/mint receipt precedes file publication. The file helper rechecks current session and expiry before emission. In-progress readback is allowed with no second effect.

| Pin | Witness and scope |
|---|---|
| `TestOperatorFirstManagedEpisodeProcess` | Uses real migration, signed start, fake supervisor/install and private temporary Git workspace. Confirms a pre-digest stale binding, current untouched release, non-executable and stale-heartbeat refusals with equal observations, terminal episode-lock contention, two stop boundaries, lifecycle refusal, and one real production-loop tick. The Git hook re-execs the test binary to check the inherited runner lock. Reads the expected owner and prepared window, refuses terminal revoke and GC, verifies no second effect and observes the fake agent marker absent. |
| `TestOperatorCurrentOwnerProcessRecovery` | Same fake supervisor/install and temp-root fixture. The independent process holds the episode lock while the current and stale heartbeats each refuse. After exit the readback holds exclusion and preserves evidence. A moved worktree, forged token or task bytes, session metadata drift, equality-at-expiry and missing credential each refuse. The independent `partial_orphan` arm reuses the real reservation identity to simulate the pre-observation boundary and exercises present-token refusal and attended deletion after token loss. |
| `TestOperatorOldBuildTerminalDrain` | Existing real old writer re-exec and fake supervisor/install under `t.TempDir()`. The fake process inventory is backed by the real child PID probe. Requires unchanged receipts under live contention and exact matched upgrade after exit. |
| `TestOperatorServiceUpgrade/active_child` | Existing inherited real child, fake supervisor/install under `t.TempDir()`. Checks typed approval and prompt digests, stopped admission, live child drain refusal, stopped-companion notice and compatible pair rollback. |
| `TestOperatorAttendedStartDisposition` | Real test signer, migrated stores, Git and mint abort trigger in disposable roots. Fake installed service and fake supervisor under `t.TempDir()`. A live process and explicit typed decline each refuse with evidence equal. After the typed yes only the exact reserved mint disposition changes, preserving Git, branch and files. |
| `TestOperatorPartialFilesOrphanRevoke` | Pure transaction seam in `host/store` with the repository's migrated test store and mint fixture. A stale version and a blocked or displaced delete roll back. The successful write removes the exact receipt credential while preserving the adjacent credential and start. |
| `TestNextOwnedOperatorStart`, `TestOperatorMarkMintUncertain` | Store-only cursor, exclusion of unowned or other-kind rows, transactional compare and uncertain replay. Observe full logical evidence before and after every refusal. |
| `TestOperatorActivationDrainRefusals`, `TestOperatorRevokeStatRefusal` | Fake supervisor/install in `t.TempDir()`. Distinct verified from-build, elapsed budget, aliased canonical drain path, already-stopped failed upgrade, and looped versus missing credential stat are covered. See §4.2 for the exact judge edits. |

| Durable window | Response after a restart and working pin |
|---|---|
| Claimed without reservation | Keep the old owner until process/episode exclusion, exact current heartbeat and signed catalogue, current approver and physical Git/branch/files absence. Release by CAS. Stale proof needs a fresh review and the attended stop. `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorStartRecoveryWindows/claimed`. |
| Worktree reserved | Keep reservation and direct to `reconcile-worktree`. The existing Git readback verifies registration, pinned commit and physical ref after quiescence. `TestOperatorWorktreeWindow`, `TestOperatorStartRecoveryWindows/reserved_worktree`. |
| Worktree observed without mint | Keep exact tree and direct to `reconcile-start`. Under both locks verify the exact worktree and dispose to failed only on typed approval. `TestOperatorStartRecoveryWindows/observed_worktree`. |
| Mint reserved without commit | All atomic records absent permits typed stop after worktree readback. Contradiction retains owner and receipts and may be marked uncertain. `TestOperatorAttendedStartDisposition`, `TestOperatorMarkMintUncertain`, and `TestOperatorStartRecoveryRejectsContradictoryEvidence/nonempty`. |
| Mint observed without files | Keep atomic receipt. `reconcile-mint` checks exact current credential and token evidence. Proven orphan revocation preserves the start and tree. `TestOperatorMintRecovery`, `TestOperatorStartRecoveryWindows/observed_mint`. |
| Files reserved or partial | `reconcile-files` can read exact or absent evidence. With a matching attributable directory and missing token, attended `reconcile-start` revokes only the exact orphan, with files, Git and start retained. Collision refuses before write. `TestOperatorFilesReconcile`, `TestOperatorPartialFilesOrphanRevoke`. |
| Worktree, mint and files observed | Derive `prepared-awaiting-S6`, then observe exact identity, path, branch, retained store rows and private files under episode exclusion. No replay. `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorCurrentOwnerProcessRecovery`. |

## 4. Verification and mutation ledger

### 4.1 Gates, fail sets and outside replay

A passing activation has literal PASS lines for all six Appendix C pins, including `TestOperatorServiceUpgrade/active_child`, and all seven store window arms. `$S/accept.sh` enforces those names. It prints `UNINFORMATIVE UNDER SANDBOX` and returns 77 on a bind denial. Whole CLI plain/race and optional daemon-race comparisons require valid baseline and stage files. Entire verifygate is recorded at both stages, including the documented `TestMissionConfigDefaultFlowIntegration` load case. A new case is inspected and rerun on both trees in isolation before attributing it. Do not change the row-271 real-binary cases, the listed git-fixture/setup flakes or the CI timeout in this sprint.

The fail-set parser prints `FAILSET_CONTROL host/store TestParserPositive` per invocation. Check every raw log for `panic: test timed out`. A truncated output, setup/build error, missing summary, or nonzero rc with empty set is an instrument error. Leg comparison uses `comm -13` on sorted sets after every recording leg completes. Stable base-only rows are recorded separately. The gates use text Go output throughout and no PATH shim.

Outside the sandbox run these commands after the controller baseline writes `ALL_DONE` and while the rig is idle for the timing leg. The base comparator is `$STATE/base/cmd-plain.{raw,rc,fails}`, `store-{plain,race}.{raw,rc,fails}`, and `verifygate-plain.{raw,rc,fails}`. The outside script records the missing CLI-race and verbs base legs on the same pristine archive. Set `WORLD_PLAN_DAEMON_RACE=1` only if a comparable whole daemon timing is needed. Its added named test total is zero.

````bash
STATE=$HOME/.ailang/state/world-iter442
S=$STATE/plan-scratch
D=$(mktemp -d "$S/outside.XXXXXX")
/bin/bash "$S/outside.sh" "$D/results"
GOMAXPROCS=4 /bin/bash "$S/measure-race.sh" "$D/race-timing"
/bin/bash "$S/whole-mutants.sh" "$S/whole-mutant-ids.txt" "$D/results/base-extra"
/bin/bash "$S/go-leg.sh" "$S/st-M1" "$D/startup.raw" 230 go test -v -count=1 -timeout 3m ./cmd/ailang-worldd/ -run '^TestOperatorServiceStartup$'
/bin/bash "$S/go-leg.sh" "$S/st-M1" "$D/native-census.raw" 120 go run "$S/kinfo-control.go"
grep -F 'self=true error=<nil>' "$D/native-census.raw"
for pair in "base:$S/base-tree" "head:$S/st-M1"; do
  tag=${pair%%:*}
  tree=${pair#*:}
  /bin/bash "$S/go-leg.sh" "$tree" "$D/legacy-windows-$tag.raw" 420 go test -v -count=1 -timeout 6m ./cmd/ailang-worldd/ -run '^TestOperator(WorktreeWindow|FilesReconcile)$'
done
````

Previous test setup, optional network controls and the native process census may need outside replay only if the corresponding observed log says permission denied. The test process census in the activation and drain tests is a fixture seam with a real child. The local native read-only control printed `uid scan count=0 self=false error=operation not permitted` in `$S/kinfo-control.raw`. Native enumeration is **UNINFORMATIVE UNDER SANDBOX**. The existing name-based native census and no-service-record stop refusal remain design constraints. There are zero production lines in those functions in M1.

Controller `cmd-plain.raw` hit `panic: test timed out after 20m0s`. Its 14-line fail set is truncated. `outside.sh` falls back to a fresh outside pristine recording and passes that comparator to `whole-mutants.sh`.

Controller `daemon-race.raw` also hit the 20-minute package timeout. The whole-daemon delta is unmeasured. `WORLD_PLAN_DAEMON_RACE=1` adds a fresh matched base/head leg to the outside script.

The first staged store-race reported `TestConvertedStoreMethodsM5a/cutoff/PutObject` at the cutoff. The code driving that test is outside the diff. The exact cutoff selector passed 3/3 times on each tree in `$S/cutoff-{base-tree,st-M1}.raw`. Whole retry evidence is kept separately.

### 4.2 Diff-anchored and carried edits

The generator `python3 "$S/make-mutants.py"` pins exact before/after text, stage line, package and named killer in `$S/mutants.json`. The script checks that each edit has exactly one textual anchor, creates a scratch copy and restores the original bytes after each run. Use `go test -v -count=1 -timeout 3m` for named tests and `-timeout 20m` for whole packages. A compile error or a denial cannot kill a mutant. The verbatim named `--- FAIL:` and whole-package fail-set delta decide the two ledgers. A positive test control is required before using its mutant evidence.

| Ledger | Selection | Results |
|---|---|---|
| carried | named | 5 KILLED |
| carried | whole | 5 PENDING |
| diff | named | 136 KILLED, 1 EQUIVALENT |
| diff | whole | 109 PENDING, 27 KILLED, 1 EQUIVALENT |
| preserved | named | 1 KILLED |
| preserved | whole | 1 KILLED |

Named transcripts include `$S/ledger-cli-named.jsonl`, `ledger-store-named.jsonl`, `ledger-supp-named.jsonl`, `ledger-supp-later-named.jsonl`, `ledger-readback-retry-named.jsonl`, `ledger-lock-fix-named.jsonl`, and `ledger-carried-named.jsonl`. Corrections and their exact reruns remain in those raw receipts and `ledger-remedies-named.jsonl`. Whole store and verbs receipts are in `ledger-store-whole.jsonl` and `ledger-preserved-whole.jsonl`.

`D-OWNED-LIMIT` is observationally equivalent for the result of this selector. Ordered `QueryRowContext` consumes one row. SQL keeps the explicit bound for the query plan.

Non-equivalent named survivors in the latest readings **zero**.

Review the remaining candidate rows with `python3 "$S/branch-audit.py" | grep -F UNINSTRUMENTED`. Some are continuations of an already tested condition. A stand-alone case shown there still needs its own behavioural edit and replay.

| Carried id and stage anchor | Exact source edit | Named test result |
|---|---|---|
| `J-FROM-BINARY`, `cmd/ailang-worldd/operator_pair_upgrade.go:45` | `approved[0] != serviceHash(seen) || c.Executables[e.paths.Binary] != from` → `approved[0] != serviceHash(seen)` | KILLED by `TestOperatorActivationDrainRefusals/from_binary` |
| `J-BUDGET-EXPIRED`, `cmd/ailang-worldd/operator_pair_upgrade.go:128` | `if !e.runnerNow().Before(deadline) {` → `if false {` | KILLED by `TestOperatorActivationDrainRefusals/budget_expired` |
| `J-RUNNER-STOP-CANON`, `cmd/ailang-worldd/operator_service.go:423` | `if err != nil || creds != r.Credentials {` → `if err != nil {` | KILLED by `TestOperatorActivationDrainRefusals/stop_canonical` |
| `J-REVOKE-MISSINGSTAT`, `cmd/ailang-worldd/session.go:485` | `if _, err := os.Stat(credentials); errors.Is(err, os.ErrNotExist) {` → `if _, err := os.Stat(credentials); err != nil {` | KILLED by `TestOperatorRevokeStatRefusal` |
| `O-DROP-LOADED`, `cmd/ailang-worldd/service_upgrade.go:224` | `if pair == nil || !pair.loaded {` → `if pair == nil {` | KILLED by `TestOperatorActivationDrainRefusals/stopped_failure` |

Those five checks remain in M1. The initial J-BUDGET-EXPIRED harness timed out under the mutant. The bounded direct drain test then produced the named failure with the exact judge edit. Early owner/freshness survivors were corrected by adding the untouched-release refusal arms. The attended lock-release test now prevents file finalizers from satisfying the explicit-close assertion.

The `$S/final-{diff,carried,preserved}-{named,whole}.jsonl` indexes contain the latest reading of each replayed id. The new product edits and carried ids are counted separately. `$S/diff-branch-audit.tsv` lists changed branch/order/SQL candidates with anchor coverage. Its 168 candidate lines include control structure and compound-condition continuations. 96 candidate lines overlap an exact mutation anchor. The remaining audit rows are visible for review. Do not count them as tested edits. `$S/whole-mutant-ids.txt` is the bounded outside priority list. Use §4.1 for CLI whole-package confirmation because the sandbox trips the unrelated package bind panic before reaching the operator pins.

## 5. Existing tests and helpers

| Existing test or helper | Production function before | Production function after | Behavioural mutation |
|---|---|---|---|
| `TestOperatorProcessWitness` in `operator_pair_upgrade_test.go` | Re-exec lock holder exercised by `swapUpgrade`, `waitRunnerExitUntil`, and old-terminal census | Original modes plus `git-inherited` / `owner` used by the primary activation and recovery pins. The same `syncBuffer` remains the subprocess sink. | `D-GIT-RUNNER-FD` omits the inherited runner fd. `D-READBACK-LOCK` points the recovery lock at a different episode. `J-BUDGET-EXPIRED` exercises the same drain implementation's bounded return. |
| `TestOperatorNeverDispatch` in `host/verbs/operator_policy_test.go` | App `AdmitOperator` via `admitLine`, including service-family refusal and shell rejection | Same production path and assertions plus terminal-only `reconcile-start` | `P-APP-SERVICE` reroutes service out of the never-dispatch set. |

The inherited `host/store/operator_start_recovery_test.go` adds independent functions. The two round-two bodies and all existing setup helpers remain byte-for-byte. The existing `TestOperatorServiceUpgrade`, `TestOperatorOldBuildTerminalDrain`, borrowed worktree/file refusals and file publication fixture are rerun through existing callers.

## 6. Gate script

The executor and controller extract the exact fence below. It records each raw result and prints `ALL_DONE` to the marker file last. Exit 0 means all required comparisons passed, 77 means an outside reading remains, and other values require inspecting the recorded legs. A skipped optional leg has a `LEG ... SKIPPED` record.

````bash gates.sh
#!/bin/bash
# gates.sh TREE LOGDIR base|M1 [PRISTINE-GATE-LOGDIR]
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P) || exit 2
T=$(cd "$1" && pwd -P) || exit 2
mkdir -p "$2" || exit 2
L=$(cd "$2" && pwd -P) || exit 2
M=$3
BASE_LOG=${4:-}
case "$M" in base|M1) ;; *) exit 2;; esac
test ! -e "$L/legs.txt" || { echo "Use a fresh log directory"; exit 2; }
PIN=$HOME
export WORLD_PLAN_PIN_HOME="$PIN" GIT_OPTIONAL_LOCKS=0
export AILANG_BIN="$PIN/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$PIN/.pinned-ailang-tools/v0.54.15/ailang"
export GOCACHE="$(go env GOCACHE)" GOMODCACHE="$(go env GOMODCACHE)"
export HOME="$L/home" AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
export XDG_CONFIG_HOME="$HOME/.config" XDG_CACHE_HOME="$HOME/.cache" XDG_DATA_HOME="$HOME/.local/share"
mkdir -p "$HOME" "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_TEST_OPERATOR_WITNESS WORLD_TEST_OPERATOR_WITNESS_PATH
unset WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
cd "$T" || exit 2
fail=0 unknown=0
line() { printf '%s\n' "$*" | tee -a "$L/legs.txt"; }
skip() { line "LEG $1 SKIPPED ($2)"; }
bound() { python3 "$S/bounded.py" "$@"; }
denied() { grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$1"; }
truncated() { grep -Eq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$1"; }
leg() {
  local tag=$1 start rc=0 result; shift
  start=$(date +%s)
  "$@" > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  result="rc=$rc $(( $(date +%s)-start ))s"
  if truncated "$L/$tag.raw"; then fail=1; result="$result INSTRUMENT FAILURE (truncated)"
  elif denied "$L/$tag.raw" || [ "$rc" -eq 77 ]; then unknown=1; result="$result UNINFORMATIVE UNDER SANDBOX"
  elif [ "$rc" -ne 0 ]; then fail=1; fi
  line "LEG $tag $result"
}
record() {
  local tag=$1 pkg=$2 start rc=0 parser=0 result; shift 2
  start=$(date +%s)
  bound 1230 go test -count=1 -timeout 20m "$@" "./$pkg/" > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  python3 "$S/failset.py" "$L/$tag.raw" "$rc" > "$L/$tag.fails" 2> "$L/$tag.parser" || parser=$?
  cat "$L/$tag.parser"
  result="rc=$rc parser=$parser $(( $(date +%s)-start ))s"
  if truncated "$L/$tag.raw"; then fail=1; result="$result INSTRUMENT FAILURE (truncated)"
  elif denied "$L/$tag.raw"; then unknown=1; result="$result UNINFORMATIVE UNDER SANDBOX"
  elif [ "$parser" -ne 0 ] || [ "$rc" -gt 1 ]; then fail=1; result="$result INSTRUMENT FAILURE"
  else cp "$L/$tag.fails" "$L/$tag.attributable"; fi
  line "LEG $tag $result"
}
fmt_touched() {
  local path out= bad=
  test -s "$S/changed-M1" || return 2
  while IFS= read -r path; do
    case "$path" in *.go) if [ -f "$path" ]; then out=$(gofmt -l "$path") || return 2; bad="$bad$out"; fi;; esac
  done < "$S/changed-M1"
  printf '%s\n' "$bad"
  test -z "$bad"
}
inventory() {
  local pkg=$1 names=$2 name rc=0 output
  output=$(mktemp "$L/inventory.XXXXXX") || return 2
  bound 180 go test -v -count=1 -timeout 2m "./$pkg/" -run "^($names)\$" > "$output" 2>&1 || rc=$?
  cat "$output"
  [ "$rc" -eq 0 ] || return "$rc"
  while IFS= read -r name; do grep -Fq -- "--- PASS: $name (" "$output" || return 2; done < <(printf '%s\n' "$names" | tr '|' '\n')
}
leg gofmt fmt_touched
leg build bound 300 go build ./...
leg vet bound 300 go vet ./cmd/ailang-worldd/ ./host/store/ ./host/verbs/
leg linux-vet env GOOS=linux python3 "$S/bounded.py" 300 go vet ./cmd/ailang-worldd/ ./host/store/ ./host/verbs/
leg inv-store inventory host/store 'TestProductionContextRoots|TestProductionGoSurface|TestStoreSurfaceTakesContext|TestCredentialsMethodsRequireADeadline'
leg inv-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
leg inv-cli inventory cmd/ailang-worldd 'TestCLICredentialWriterSites|TestSessionRootBudgets|TestEveryVerbHasHelp'
leg inv-verbs inventory host/verbs TestOperatorNeverDispatch
leg inv-daemon inventory host/daemon TestDaemonDependencyAllowlist
leg inv-sink inventory host/verifygate TestSubprocessSinkGate
if [ "$M" = base ]; then skip accept 'activation pins absent at base'
else leg accept /bin/bash "$S/accept.sh" "$T" M1 fast; fi
leg gc-probe env HOME="$PIN" /bin/bash "$S/probe.sh" "$T"
leg verify_ail bound 510 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
leg callers /bin/bash "$S/callers.sh" "$T" "$L/callers.census" "$M"
record store-plain host/store
record store-race host/store -race
record cmd-plain cmd/ailang-worldd
record cmd-race cmd/ailang-worldd -race
record verifygate-plain host/verifygate
record verbs-plain host/verbs
record verbs-race host/verbs -race
if [ "${WORLD_PLAN_DAEMON_RACE:-0}" = 1 ]; then record daemon-race host/daemon -race
else skip daemon-race 'no daemon product or test diff; opt in for timing evidence'; fi
# Every recording leg above finishes before comparison.
controller="${WORLD_PLAN_BASE_FAILSETS:-$S/../base}"
for tag in store-plain store-race cmd-plain cmd-race verifygate-plain verbs-plain verbs-race daemon-race; do
  if [ "$M" = base ]; then skip "compare-$tag" 'pristine recording'; continue; fi
  if [ ! -f "$L/$tag.attributable" ]; then skip "compare-$tag" 'head recording unavailable'; continue; fi
  source="${BASE_LOG:+$BASE_LOG/$tag.attributable}"
  if grep -Fq 'ALL_DONE for ad1a482c2a8755f616505a9b2d93f3369310ae01' "$controller/../baseline.done" 2>/dev/null &&
      [ -s "$controller/$tag.raw" ] && [ -f "$controller/$tag.rc" ] &&
      ! truncated "$controller/$tag.raw" && ! denied "$controller/$tag.raw" &&
      python3 "$S/failset.py" "$controller/$tag.raw" "$(cat "$controller/$tag.rc")" > "$L/$tag.controller"; then
    source="$L/$tag.controller"
  fi
  if [ ! -f "$source" ]; then
    skip "compare-$tag" 'attributable base recording unavailable'
    unknown=1
    continue
  fi
  LC_ALL=C sort -u "$source" > "$L/$tag.base"
  LC_ALL=C sort -u "$L/$tag.attributable" > "$L/$tag.head"
  LC_ALL=C comm -13 "$L/$tag.base" "$L/$tag.head" > "$L/$tag.new"
  LC_ALL=C comm -23 "$L/$tag.base" "$L/$tag.head" > "$L/$tag.base-only"
  if [ -s "$L/$tag.new" ]; then
    fail=1; line "LEG compare-$tag NEW FAIL SET (see $L/$tag.new)"
  else line "LEG compare-$tag OK"; fi
done
line "GATES_DONE fail=$fail uninformative=$unknown"
printf 'ALL_DONE\n' > "$L/ALL_DONE"
[ "$fail" -eq 0 ] || exit 1
[ "$unknown" -eq 0 ] || exit 77
````

## 7. Acceptance, snapshot and executor handoff

`accept.sh TREE M1 fast|race` succeeds only with **exit 0 and a final `ACCEPT_OK M1 fast|race <log-directory>`**. The explicit named assertions prevent a red or zero-tests-selected result from being treated as success. Commit uses fast mode with the complete environment pinned. The outside controller runs race before accepting the sprint.

````bash accept.sh
#!/bin/bash
# accept.sh TREE M1 fast|race
set -euo pipefail
T=$(cd "$1" && pwd -P)
M=$2
mode=$3
PIN=${WORLD_PLAN_PIN_HOME:-$HOME}
S="$PIN/.ailang/state/world-iter442/plan-scratch"
L=$(mktemp -d "$S/accept-XXXXXX")
export AILANG_BIN="$PIN/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$PIN/.pinned-ailang-tools/v0.54.15/ailang"
export GOCACHE="$(go env GOCACHE)" GOMODCACHE="$(go env GOMODCACHE)" GIT_OPTIONAL_LOCKS=0
export HOME="$L/home" AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
export XDG_CONFIG_HOME="$HOME/.config" XDG_CACHE_HOME="$HOME/.cache" XDG_DATA_HOME="$HOME/.local/share"
mkdir -p "$HOME" "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_TEST_OPERATOR_WITNESS WORLD_TEST_OPERATOR_WITNESS_PATH
unset WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
args=()
case "$mode" in fast) ;; race) args+=(-race);; *) exit 2;; esac
cd "$T"
test -x "$AILANG_BIN" && test -x "$WORLD_EXEC_NODE" && test -x "$WORLD_TOOL_AILANG_BIN"
run() {
  local pkg=$1 names=$2 name log="$L/${1//\//-}.raw" rc=0
  python3 "$S/bounded.py" 420 go test ${args[@]+"${args[@]}"} -v -count=1 -timeout 6m "./$pkg" -run "^($names)\$" > "$log" 2>&1 || rc=$?
  cat "$log"
  if grep -Ei 'panic: test timed out|BOUND EXCEEDED' "$log"; then echo "INSTRUMENT FAILURE: truncated $log"; return 2; fi
  if grep -Ei '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$log"; then echo "UNINFORMATIVE UNDER SANDBOX: $log"; return 77; fi
  [ "$rc" -eq 0 ] || return "$rc"
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log" || { echo "MISSING PASS $name"; return 1; }
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
activation='TestOperatorFirstManagedEpisodeProcess|TestOperatorCurrentOwnerProcessRecovery|TestOperatorAttendedStartDisposition|TestOperatorActivationDrainRefusals|TestOperatorRevokeStatRefusal'
old='TestOperatorOldBuildTerminalDrain|TestOperatorServiceUpgrade|TestOperatorPairedUpgradeRefusals|TestOperatorServiceFence|TestOperatorRunnerSetupRefusals|TestOperatorRunnerPathUID|TestOperatorTerminalRevokeExclusion|TestSessionRootBudgets|TestCLICredentialWriterSites'
borrow='TestOperatorBorrowedWorktreeGitProcess|TestOperatorFilePublicationIdentityWindows|TestOperatorFilePublicationAdmission|TestOperatorStartCandidateProcess|TestOperatorMintRecovery'
matrix='TestOperatorStartRecoveryWindows|TestOperatorStartRecoveryRejectsContradictoryEvidence|TestOperatorStartDispositionRejectsBlobOwner|TestOperatorPartialFilesOrphanRevoke|TestOperatorMarkMintUncertain|TestNextOwnedOperatorStart|TestProductionContextRoots'
[ "$M" = M1 ] || exit 2
run host/store "$matrix"
run host/authority TestOperatorPreDigestSignedStart
run host/verbs TestOperatorNeverDispatch
run cmd/ailang-worldd "$activation|$old|$borrow"
grep -Fq -- '--- PASS: TestOperatorServiceUpgrade/active_child (' "$L/cmd-ailang-worldd.raw"
grep -Fq -- '--- PASS: TestOperatorCurrentOwnerProcessRecovery/partial_orphan (' "$L/cmd-ailang-worldd.raw"
for arm in claimed reserved_worktree observed_worktree reserved_mint observed_mint reserved_files prepared; do
  grep -Fq -- "--- PASS: TestOperatorStartRecoveryWindows/$arm (" "$L/host-store.raw"
done
printf 'ACCEPT_OK %s %s %s\n' "$M" "$mode" "$L"
````

Run `/bin/bash "$T/snapshot.sh" "$W" M1` exactly after acceptance. It writes cumulative `$STATE/snap/M1/{STATUS,DELETED,MANIFEST.sha256}` plus the changed files. Planning records are excluded. The script accepts an explicit status list only for an archive rehearsal and has no dependency on the existence of optional product directories.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1) ;; *) exit 2;; esac
B=ad1a482c2a8755f616505a9b2d93f3369310ae01
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter442/snap}/$M
mkdir -p "$(dirname "$D")"
mkdir "$D"
D=$(cd "$D" && pwd -P)
cd "$W"
if [ -n "${SNAP_STATUS_LIST:-}" ]; then
  cp "$SNAP_STATUS_LIST" "$D/STATUS"
else
  {
    git diff --name-status --no-renames "$B" --
    git ls-files --others --exclude-standard | sed $'s/^/A\t/'
  } | awk -F '\t' '
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s5f4-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s5f4.json"
  ' | LC_ALL=C sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
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

Controller commit command, after committing the plan records and after comparing the exact M1 snapshot with the sealed stage. `commit_ms.sh` performs gofmt, build and touched-package vet before running this command. The verifygate wrapper prints `VERIFYGATE_FAILSET_OK` only when the complete head set adds no failure to the complete controller base set.

````bash
export W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter442-row243-v4s5f4
STATE=$HOME/.ailang/state/world-iter442
export S=$STATE/plan-scratch
/bin/bash "$STATE/commit_ms.sh" M1 "$S/M1.msg" '
  export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
  export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
  /bin/bash "$S/accept.sh" "$W" M1 fast &&
  L=$(mktemp -d "$S/commit-verifygate-M1.XXXXXX") &&
  /bin/bash "$S/commit-verifygate.sh" "$W" "$L" go test -count=1 -timeout 20m ./host/verifygate/
'
````

The final Appendix P fence was extracted from the plan and applied with `/usr/bin/patch -p1 -E` to a fresh export in `$S/rt-1ah4ZF/`. M1 `diff -r` returned 0 for `cmd`, `host`, `docs` and the full archive. The cumulative snapshot had 14 STATUS entries, 0 DELETED entries and all 14 SHA-256 entries verified byte-for-byte. See `$S/roundtrip-final.raw`. The controller's `ctl-roundtrip` path was left available for its own replay.

## 8. Scope ledger

| Scope | Disposition |
|---|---|
| First production preparation consumer | M1. One bounded confirmed-start visit, both lock descriptions, initial and current proof and store-pair reads, real Git reservation/readback, atomic mint and private file publication through prepared. See the two process pins. |
| Seven durable cut points and gap remedies | M1. See §3. All completed prep recovery is read-only. Untouched release and attended remedies preserve prior effects. |
| Upgrade and stop drain | M1 acceptance exercises the landed S5f3 paths with `active_child` and old terminal drain. New behavioural refusal pins are test-only. |
| Operator-visible surface | Adds `service runner reconcile-start ACTION-ID` help and typed prompts for stopping the action with retained evidence, revoking the exact orphan, or marking mint uncertain with owner retained. Success prints `Start action stopped; existing evidence retained: <id>` or `Mint commit uncertain; owner retained: <id>`. Refusals use `start recovery refused` with the why and recovery next step. Runner logs include `runner halted` with the service-record-unavailable why, `operator start binding stale (<field>); request fresh review`, `operator start blocked; attended recovery: <command>`, and `operator start blocked; inspect service runner status and reconcile durable evidence`. Preparation is represented by the durable readback window, with the existing readiness projection. `docs/QUICKSTART.md` adds a disposable example. |
| S6 and later | Agent launch, spawn reservation and any subsequent browser route are later reviewed sprints. This plan carries no product lines for them. |
| Scope proof | `grep -E -e '\.ail$' -e /app/ -e /tools/launchd/ -e /host/daemon/ -e /migration "$S/unfiltered-product-paths.txt"` exited 1. The positive `grep -F '/cmd/ailang-worldd/operator_runner.go'` matched its measured path. `grep -nF 'operatorSchemaVersion = 10' "$S"/{base-tree,st-M1}/host/store/operator_runner.go` printed :12 in both archives. Full comparison and manifest enumerate the product paths. |

## Appendix P — compiled patches

````diff M1
diff --git a/cmd/ailang-worldd/operator_activation_lifecycle_test.go b/cmd/ailang-worldd/operator_activation_lifecycle_test.go
new file mode 100644
index 00000000..4c103e67
--- /dev/null
+++ b/cmd/ailang-worldd/operator_activation_lifecycle_test.go
@@ -0,0 +1,118 @@
+package main
+
+import (
+	"encoding/json"
+	"errors"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+)
+
+func TestOperatorActivationDrainRefusals(t *testing.T) {
+	t.Run("from_binary", func(t *testing.T) {
+		f, _, candidate := pairedUpgradeFixture(t)
+		before := pairedUpgradeEvidence(f)
+		hash := serviceFileHash(candidate)
+		if hash == before[0] || hash == "" {
+			t.Fatal("need distinct verified build hashes")
+		}
+		pair, err := f.env.reviewRunnerUpgrade(f.env.runnerUpgradeReview(), hash, hash)
+		if pair != nil || err == nil || !strings.Contains(err.Error(), "changed during confirmation") ||
+			pairedUpgradeEvidence(f) != before {
+			t.Fatalf("approved pair crossed the from-build binding: %v", err)
+		}
+	})
+	t.Run("budget_expired", func(t *testing.T) {
+		f, s, candidate := pairedUpgradeFixture(t)
+		before := pairedUpgradeEvidence(f)
+		start := time.Unix(1000, 0)
+		elapsed := false
+		f.env.runnerNow = func() time.Time {
+			if elapsed {
+				return start.Add(f.env.wait)
+			}
+			elapsed = true
+			return start
+		}
+		candidateBefore := serviceFileHash(candidate)
+		guard, err := f.env.drainRunner(f.config)
+		if guard != nil {
+			guard.Close()
+		}
+		_, lockErr := os.Lstat(f.config.Credentials + ".operator-runner.lock")
+		if guard != nil || err == nil || !strings.Contains(err.Error(), "companion drain budget expired") || s.runnerLoaded ||
+			f.homeJob.kicks != 0 || pairedUpgradeEvidence(f) != before || serviceFileHash(candidate) != candidateBefore ||
+			!errors.Is(lockErr, os.ErrNotExist) {
+			t.Fatalf("expired drain acquired the lock or swapped: %v (%v)", err, lockErr)
+		}
+	})
+	t.Run("stop_canonical", func(t *testing.T) {
+		f, s, _ := pairedUpgradeFixture(t)
+		alias := filepath.Join(f.home, "credentials-alias")
+		if err := os.Symlink(f.config.Credentials, alias); err != nil {
+			t.Fatal(err)
+		}
+		got, gotErr := filepath.EvalSymlinks(alias)
+		want, wantErr := filepath.EvalSymlinks(f.config.Credentials)
+		if gotErr != nil || wantErr != nil || got != want || got == alias {
+			t.Fatal("need a distinct path to the same credentials")
+		}
+		r, err := readServiceRecord(f.env.paths)
+		if err != nil {
+			t.Fatal(err)
+		}
+		r.Credentials = alias
+		b, err := json.Marshal(r)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if err := f.env.atomic(f.env.paths.Record, b, 0600); err != nil {
+			t.Fatal(err)
+		}
+		before := pairedUpgradeEvidence(f)
+		code, msg := f.command("stop")
+		if code != 1 || !strings.Contains(msg, "canonical credentials lock") || !s.runnerLoaded ||
+			pairedUpgradeEvidence(f) != before {
+			t.Fatalf("stop accepted an alias in the service record: %d %s", code, msg)
+		}
+	})
+	t.Run("stopped_failure", func(t *testing.T) {
+		f, s, candidate := pairedUpgradeFixture(t)
+		s.runnerLoaded = false
+		before := pairedUpgradeEvidence(f)
+		read, calls := f.env.runnerProcesses, 0
+		f.env.runnerProcesses = func() ([]worldProcess, error) {
+			calls++
+			if calls == 2 {
+				return nil, errors.New("fixture post-drain census denied")
+			}
+			return read()
+		}
+		code, msg := f.run("upgrade", "--binary", candidate)
+		if code != 1 || !strings.Contains(msg, "terminal drain refused") || strings.Contains(msg, wantStoppedRunnerRemedy) ||
+			s.runnerLoaded || f.homeJob.kicks != 0 || pairedUpgradeEvidence(f) != before {
+			t.Fatalf("pre-stopped companion failure: %d %s", code, msg)
+		}
+	})
+}
+
+func TestOperatorRevokeStatRefusal(t *testing.T) {
+	f := operatorService(t)
+	path := filepath.Join(f.home, "stat-loop")
+	if err := os.Symlink(path, path); err != nil {
+		t.Fatal(err)
+	}
+	before := pairedUpgradeEvidence(f)
+	ctx := boundedTestContext(t)
+	id := strings.Repeat("a", 64)
+	if _, ok, err := readSessionRevokeTarget(ctx, f.config.World, path+".absent", id); ok || err != nil {
+		t.Fatalf("missing credentials control: %t %v", ok, err)
+	}
+	_, ok, err := readSessionRevokeTarget(ctx, f.config.World, path, id)
+	target, linkErr := os.Readlink(path)
+	if ok || err == nil || errors.Is(err, os.ErrNotExist) || target != path || linkErr != nil || pairedUpgradeEvidence(f) != before {
+		t.Fatalf("stat failure was treated as missing credentials: %t %v", ok, err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_activation_test.go b/cmd/ailang-worldd/operator_activation_test.go
new file mode 100644
index 00000000..75852c43
--- /dev/null
+++ b/cmd/ailang-worldd/operator_activation_test.go
@@ -0,0 +1,484 @@
+package main
+
+import (
+	"bytes"
+	"encoding/json"
+	"errors"
+	"io"
+	"os"
+	"path/filepath"
+	"reflect"
+	"runtime/debug"
+	"strings"
+	"syscall"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func operatorActivationFixture(t *testing.T) (*worktreeStartRig, runnerConfig, store.OperatorHeartbeat, *os.File) {
+	t.Helper()
+	r := newWorktreeStartRig(t, false)
+	c := r.config
+	agent := operatorFileFact(r.action, "launch_binary")
+	c.PresetNames = []string{r.action.Options.Preset}
+	c.Executables[agent] = serviceFileHash(agent)
+	p := r.env.runnerPaths()
+	config, err := json.Marshal(c)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := r.env.atomic(p.Record, config, 0600); err != nil {
+		t.Fatal(err)
+	}
+	if err := r.env.atomic(p.Plist, []byte(runnerPlist(p, c, config)), 0600); err != nil {
+		t.Fatal(err)
+	}
+	r.env.controller = r.homeJob
+	r.env.runnerProcesses = func() ([]worldProcess, error) {
+		return []worldProcess{{os.Getpid(), "ailang-worldd"}}, nil
+	}
+	record, err := readServiceRecord(r.env.paths)
+	if err != nil {
+		t.Fatal(err)
+	}
+	record.Workspace = r.src.workspace
+	body, err := json.Marshal(record)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := r.env.atomic(r.env.paths.Record, body, 0600); err != nil {
+		t.Fatal(err)
+	}
+	if err := os.Mkdir(r.f.r.sessions, 0700); err != nil {
+		t.Fatal(err)
+	}
+	r.exec(t, "UPDATE operator_actions SET state='confirmed', owner=NULL")
+	build := c.Executables[r.env.paths.Binary]
+	catalogue, err := json.Marshal(c.PresetNames)
+	if err != nil {
+		t.Fatal(err)
+	}
+	r.exec(t, "UPDATE operator_runner_state SET build=?, catalogue=?", build, catalogue)
+	h, err := runnerHeartbeat(r.creds, time.Second)
+	if err != nil {
+		t.Fatal(err)
+	}
+	held, err := runnerKernelLock(r.creds + ".operator-runner.lock")
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { held.Close() })
+	return r, c, h, held
+}
+
+func activationEvidence(t *testing.T, r *worktreeStartRig) [5]any {
+	t.Helper()
+	git := filepath.Join(r.f.r.f.repo, ".git")
+	return [5]any{serviceFileHash(r.creds), provisionFiles(t, filepath.Join(git, "refs")),
+		provisionFiles(t, filepath.Join(git, "worktrees")), provisionFiles(t, r.src.workspace), provisionFiles(t, r.f.r.sessions)}
+}
+
+func TestOperatorFirstManagedEpisodeProcess(t *testing.T) {
+	r, c, h, runner := operatorActivationFixture(t)
+	ctx := boundedTestContext(t)
+	prepare := func(acquire func(string, string) (*store.OperatorEpisodeLock, error), stops <-chan os.Signal) error {
+		cursor := ""
+		return r.env.consumeOperatorPreparation(ctx, c, h, r.src, stops, runner, &cursor, acquire)
+	}
+	snapshot := func() any {
+		return activationEvidence(t, r)
+	}
+	row := candidateRow(t, r.creds, r.src.clock())
+	legacy := r.action
+	legacy.Resolution = nil
+	for _, fact := range r.action.Resolution {
+		if fact.Name != "launch_binary_sha256" {
+			legacy.Resolution = append(legacy.Resolution, fact)
+		}
+	}
+	wire, err := approval.CanonicalOperatorAction(legacy)
+	if err != nil {
+		t.Fatal(err)
+	}
+	digest, err := approval.OperatorSignatureDigest(legacy)
+	if err != nil {
+		t.Fatal(err)
+	}
+	msg, err := approval.OperatorActionMessage(digest, [32]byte(row.Challenge), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	r.exec(t, "UPDATE operator_actions SET action=?, digest=?, artifact=?", wire, digest[:], approvaltest.NewSigner(51).Sign(approval.OperatorActionNamespace, msg))
+	before := snapshot()
+	var binding *approval.BindingError
+	err = prepare(store.OpenOperatorEpisodeLock, nil)
+	if !errors.As(err, &binding) || binding.Field != "launch_binary_sha256" || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("old signed review: %v", err)
+	}
+	r.exec(t, "UPDATE operator_actions SET action=?, digest=?, artifact=?", row.Action, row.Digest, row.Artifact)
+
+	if err := runnerWrite(r.creds, func(db *store.Credentials) error {
+		return db.ClaimOperatorAction(ctx, operatorPreparationCAS(row), h.Owner, r.src.clock())
+	}); err != nil {
+		t.Fatal(err)
+	}
+	beforeRelease := activationEvidence(t, r)
+	// The live loop only caches the lease generation, never the catalogue.
+	loopHeartbeat := h
+	loopHeartbeat.CatalogueGeneration = 0
+	cursor := ""
+	r.exec(t, "UPDATE operator_runner_state SET catalogue_generation=catalogue_generation+1")
+	beforeDrift := snapshot()
+	err = r.env.recoverOperatorStart(ctx, c, loopHeartbeat, r.src, &cursor, nil)
+	if !errors.As(err, &binding) || binding.Field != "mint_catalogue" || !reflect.DeepEqual(beforeDrift, snapshot()) {
+		t.Fatalf("untouched recovery crossed a changed catalogue: %v", err)
+	}
+	r.exec(t, "UPDATE operator_runner_state SET catalogue_generation=catalogue_generation-1")
+	cursor = ""
+	foreignHeartbeat := loopHeartbeat
+	foreignHeartbeat.Owner.Generation++
+	beforeForeign := snapshot()
+	err = r.env.recoverOperatorStart(ctx, c, foreignHeartbeat, r.src, &cursor, nil)
+	if !errors.As(err, &binding) || binding.Field != "mint_catalogue" || !reflect.DeepEqual(beforeForeign, snapshot()) {
+		t.Fatalf("foreign recovery lease released the untouched claim: %v", err)
+	}
+	r.exec(t, "UPDATE operator_runner_state SET seen_at=?", r.src.clock()-10)
+	cursor = ""
+	beforeStale := snapshot()
+	err = r.env.recoverOperatorStart(ctx, c, loopHeartbeat, r.src, &cursor, nil)
+	if !errors.As(err, &binding) || binding.Field != "mint_catalogue" || !reflect.DeepEqual(beforeStale, snapshot()) {
+		t.Fatalf("stale recovery lease released the untouched claim: %v", err)
+	}
+	r.exec(t, "UPDATE operator_runner_state SET seen_at=?", r.src.clock())
+	cursor = ""
+	if err := r.env.recoverOperatorStart(ctx, c, loopHeartbeat, r.src, &cursor, nil); err != nil {
+		t.Fatalf("untouched claim with current signed catalogue: %v", err)
+	}
+	released := candidateRow(t, r.creds, r.src.clock())
+	afterRelease := activationEvidence(t, r)
+	if released.ID != row.ID || released.Owner != "" || released.Generation != row.Generation+2 ||
+		!reflect.DeepEqual(beforeRelease[1:], afterRelease[1:]) {
+		t.Fatalf("untouched release changed physical evidence or lost the confirmed action: %+v", released)
+	}
+
+	binary := operatorFileFact(r.action, "launch_binary")
+	if err := os.Chmod(binary, 0600); err != nil {
+		t.Fatal(err)
+	}
+	before = snapshot()
+	err = prepare(store.OpenOperatorEpisodeLock, nil)
+	if err == nil || !strings.Contains(err.Error(), "(F10)") || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("non-executable unchanged bytes admitted preparation: %v", err)
+	}
+	if err := r.env.checkStartPair(c, r.src.workspace); !errors.Is(err, store.ErrOperatorConflict) {
+		t.Fatalf("installed pair accepted a non-executable file: %v", err)
+	}
+	if err := os.Chmod(binary, 0700); err != nil {
+		t.Fatal(err)
+	}
+	r.exec(t, "UPDATE operator_runner_state SET seen_at=?", r.src.clock()-10)
+	before = snapshot()
+	err = prepare(store.OpenOperatorEpisodeLock, nil)
+	if !errors.Is(err, store.ErrOperatorRunner) || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("stale heartbeat admitted preparation: %v", err)
+	}
+	r.exec(t, "UPDATE operator_runner_state SET seen_at=?", r.src.clock())
+	before = snapshot()
+	oldOwner := h.Owner
+	oldOwner.Generation++
+	if err := r.env.prepareOperatorStart(ctx, c, oldOwner, h.Build, r.src, nil, runner, store.OpenOperatorEpisodeLock); !errors.Is(err, store.ErrOperatorRunner) || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("foreign runner admitted preparation: %v", err)
+	}
+	readProcesses := r.env.runnerProcesses
+	r.env.runnerProcesses = func() ([]worldProcess, error) { return nil, nil }
+	err = prepare(store.OpenOperatorEpisodeLock, nil)
+	if err == nil || !strings.Contains(err.Error(), "did not include this process") || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("unproved terminal retirement admitted preparation: %v", err)
+	}
+	r.env.runnerProcesses = readProcesses
+	healthy := r.health.Build.ExecutableSHA256
+	r.health.Build.ExecutableSHA256 = "sha256:unapproved"
+	err = r.env.checkStartPair(c, r.src.workspace)
+	if err == nil || !strings.Contains(err.Error(), "daemon build mismatch") || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("drifted World build admitted preparation: %v", err)
+	}
+	r.health.Build.ExecutableSHA256 = healthy
+
+	before = snapshot()
+	paused := false
+	err = prepare(func(path, ep string) (*store.OperatorEpisodeLock, error) {
+		paused = true
+		_, release := operatorProcessWitness(t, "episode", path)
+		defer release()
+		return store.OpenOperatorEpisodeLock(path, ep)
+	}, nil)
+	if !paused || !errors.Is(err, store.ErrOperatorConflict) || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("terminal holder did not defeat runner before CAS: %v", err)
+	}
+	stopped := make(chan os.Signal, 1)
+	stopped <- os.Interrupt
+	err = prepare(func(path, ep string) (*store.OperatorEpisodeLock, error) {
+		t.Fatal("stopped tick entered candidate selection")
+		return nil, store.ErrOperatorConflict
+	}, stopped)
+	if !errors.Is(err, errOperatorDraining) || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("stop admitted a claim: %v", err)
+	}
+	err = prepare(func(path, ep string) (*store.OperatorEpisodeLock, error) {
+		stopped <- os.Interrupt
+		return store.OpenOperatorEpisodeLock(path, ep)
+	}, stopped)
+	if !errors.Is(err, errOperatorDraining) || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("stop after candidate selection admitted a claim: %v", err)
+	}
+	bin, err := os.Executable()
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Setenv("WORLD_TEST_OPERATOR_WITNESS", "git-inherited")
+	t.Setenv("WORLD_TEST_OPERATOR_WITNESS_PATH", runner.Name())
+	t.Setenv("GORACE", strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
+	hook := filepath.Join(r.f.r.f.repo, ".git", "hooks", "post-checkout")
+	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexec "+shellQuote(bin)+" -test.run '^TestOperatorProcessWitness$'\n"), 0700); err != nil {
+		t.Fatal(err)
+	}
+	lifecycle, err := runnerLifecycleLock(filepath.Dir(r.env.paths.Record))
+	if err != nil {
+		t.Fatal(err)
+	}
+	err = prepare(store.OpenOperatorEpisodeLock, nil)
+	lifecycle.Close()
+	if !errors.Is(err, errOperatorDraining) || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("claim crossed lifecycle exclusion: %v", err)
+	}
+	ticks := make(chan time.Time)
+	close(ticks)
+	var runLog bytes.Buffer
+	if code := r.env.runOperatorPreparationLoop(&runLog, c, h, r.src, nil, ticks, runner); code != exitOK || runLog.Len() != 0 {
+		t.Fatalf("production tick failed: %d %s", code, runLog.String())
+	}
+	closed, err := runnerLifecycleLock(filepath.Dir(r.env.paths.Record))
+	if err != nil {
+		t.Fatalf("claim retained lifecycle exclusion: %v", err)
+	}
+	closed.Close()
+	if _, err := os.Stat(runner.Name() + ".git-observed"); err != nil {
+		t.Fatalf("real Git did not observe the inherited process lock: %v", err)
+	}
+	s, _, err := readRunnerFiles(ctx, r.creds, row.ID)
+	if err != nil || s.Phase != "observed" {
+		t.Fatalf("file receipt: %s %v", s.Phase, err)
+	}
+	ro, err := store.OpenCredentialsReadOnly(runnerDB(r.creds))
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer ro.Close()
+	prepared, err := ro.ReadOperatorStartRecovery(ctx, row.ID)
+	if err != nil || prepared.Window != store.OperatorStartPrepared || prepared.Owner != h.Owner.Nonce {
+		t.Fatalf("prepared owner and evidence: %+v %v", prepared, err)
+	}
+	before = snapshot()
+	var out bytes.Buffer
+	if code := runSessionRevoke([]string{"--db", c.World, s.Session.CredentialID}, &out, &out); code != exitUsage ||
+		!strings.Contains(out.String(), "active action "+row.ID) {
+		t.Fatalf("terminal could cross prepared ownership: %d %s", code, &out)
+	}
+	report := daemon.RunEpisodeGC(ctx, daemon.EpisodeGCOptions{WorkspaceRoot: r.src.workspace, CredentialsPath: r.creds, Now: time.Unix(r.src.clock(), 0)})
+	if len(report.Removed) != 0 || len(report.Skipped) != 1 || report.Skipped[0].Reason != "non-terminal operator step" ||
+		!reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("GC altered prepared owner or evidence: %+v", report)
+	}
+	if err := prepare(store.OpenOperatorEpisodeLock, nil); !errors.Is(err, store.ErrOperatorMissing) || !reflect.DeepEqual(before, snapshot()) {
+		t.Fatalf("next tick prepared the same row twice: %v", err)
+	}
+	if _, err := os.Stat(filepath.Join(r.bin, "ran")); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("spawn crossed the S6 boundary: %v", err)
+	}
+}
+
+func TestOperatorCurrentOwnerProcessRecovery(t *testing.T) {
+	r, c, h, runner := operatorActivationFixture(t)
+	ctx := boundedTestContext(t)
+	if err := r.env.prepareOperatorStart(ctx, c, h.Owner, h.Build, r.src, nil, runner, store.OpenOperatorEpisodeLock); err != nil {
+		t.Fatal(err)
+	}
+	_, release := operatorProcessWitness(t, "owner", r.creds)
+	observe := func() any {
+		return activationEvidence(t, r)
+	}
+	recover := func() error { cursor := ""; return r.env.recoverOperatorStart(ctx, c, h, r.src, &cursor, nil) }
+	for _, seen := range []int64{r.src.clock(), r.src.clock() - 11} {
+		h.SeenAt = seen
+		r.exec(t, "UPDATE operator_runner_state SET seen_at=?", seen)
+		before := observe()
+		if err := recover(); !errors.Is(err, store.ErrOperatorConflict) || !strings.Contains(err.Error(), "live episode owner") || !reflect.DeepEqual(before, observe()) {
+			t.Fatalf("heartbeat displaced live owner: %v", err)
+		}
+	}
+	release()
+	before := observe()
+	clock := r.src.clock
+	r.src.clock = func() int64 {
+		other, err := store.OpenOperatorEpisodeLock(r.creds, r.action.Episode)
+		if other != nil {
+			other.Close()
+		}
+		if !errors.Is(err, syscall.EWOULDBLOCK) {
+			t.Fatalf("prepared readback lost episode exclusion: %v", err)
+		}
+		return clock()
+	}
+	if err := recover(); err != nil || !reflect.DeepEqual(before, observe()) {
+		t.Fatalf("quiescent prepared recovery retried an effect: %v", err)
+	}
+	r.src.clock = clock
+	s, _, err := readRunnerFiles(ctx, r.creds, r.action.ID)
+	if err != nil {
+		t.Fatal(err)
+	}
+	for _, reason := range []string{"worktree", "identity", "task", "metadata", "expiry", "revocation"} {
+		t.Run(reason, func(t *testing.T) {
+			switch reason {
+			case "worktree":
+				from, to := r.path(), r.path()+".moved"
+				if err := os.Rename(from, to); err != nil {
+					t.Fatal(err)
+				}
+				t.Cleanup(func() {
+					if err := os.Rename(to, from); err != nil {
+						t.Error(err)
+					}
+				})
+			case "identity", "task":
+				path := s.Mint.TokenPath
+				if reason == "task" {
+					path = r.in.Start.Task.Path
+				}
+				token, err := os.ReadFile(path)
+				if err != nil {
+					t.Fatal(err)
+				}
+				t.Cleanup(func() {
+					if err := os.WriteFile(path, token, 0600); err != nil {
+						t.Error(err)
+					}
+				})
+				if err := os.WriteFile(path, []byte("changed file fixture"), 0600); err != nil {
+					t.Fatal(err)
+				}
+			case "metadata":
+				r.exec(t, "UPDATE session_credentials SET project='unreviewed' WHERE credential_id=?", s.Session.CredentialID)
+				t.Cleanup(func() {
+					r.exec(t, "UPDATE session_credentials SET project=? WHERE credential_id=?", s.Session.Project, s.Session.CredentialID)
+				})
+			case "expiry":
+				r.src.clock = func() int64 { return s.Session.ExpiresAt }
+				t.Cleanup(func() { r.src.clock = func() int64 { return 1012 } })
+			case "revocation":
+				r.exec(t, "DELETE FROM session_credentials WHERE credential_id=?", s.Session.CredentialID)
+				t.Cleanup(func() {
+					if err := runnerWrite(r.creds, func(db *store.Credentials) error { return db.MintSession(ctx, s.Session) }); err != nil {
+						t.Error(err)
+					}
+				})
+			}
+			before := observe()
+			if err := recover(); !errors.Is(err, store.ErrOperatorConflict) || !reflect.DeepEqual(before, observe()) {
+				t.Fatalf("recovery crossed %s evidence: %v", reason, err)
+			}
+		})
+	}
+	t.Run("partial_orphan", func(t *testing.T) {
+		runner.Close()
+		// Restore the pre-observation receipt boundary using the real file
+		// reservation identity just produced by the runner.
+		r.exec(t, `UPDATE operator_action_steps SET phase='reserved', receipt=NULL, version=version-1
+			WHERE action_id=? AND kind='files'`, r.action.ID)
+		before := activationEvidence(t, r)
+		code, out := r.command("reconcile-start", r.action.ID)
+		if code != exitUsage || !strings.Contains(out, store.ErrOperatorConflict.Error()) ||
+			!reflect.DeepEqual(before, activationEvidence(t, r)) {
+			t.Fatalf("existing token allowed orphan disposition: %d %s", code, out)
+		}
+		if err := os.Remove(s.Mint.TokenPath); err != nil {
+			t.Fatal(err)
+		}
+		before = activationEvidence(t, r)
+		code, out = r.command("reconcile-start", r.action.ID)
+		after := activationEvidence(t, r)
+		if code != exitOK || !strings.Contains(out, "existing evidence retained") || !reflect.DeepEqual(before[1:], after[1:]) {
+			t.Fatalf("partial-file disposition: %d %s", code, out)
+		}
+		if err := runnerWrite(r.creds, func(db *store.Credentials) error {
+			_, found, err := db.ResolveSession(ctx, s.Session.CredentialID)
+			if found {
+				t.Error("orphan credential was retained")
+			}
+			return err
+		}); err != nil {
+			t.Fatal(err)
+		}
+	})
+}
+
+func TestOperatorAttendedStartDisposition(t *testing.T) {
+	r, c, h, runner := operatorActivationFixture(t)
+	ctx := boundedTestContext(t)
+	r.exec(t, `CREATE TRIGGER interrupted_mint BEFORE INSERT ON session_credentials
+		BEGIN SELECT RAISE(ABORT,'fixture mint rollback'); END`)
+	err := r.env.prepareOperatorStart(ctx, c, h.Owner, h.Build, r.src, nil, runner, store.OpenOperatorEpisodeLock)
+	if err == nil || !strings.Contains(err.Error(), "fixture mint rollback") {
+		t.Fatalf("mint rollback: %v", err)
+	}
+	ro, err := store.OpenCredentialsReadOnly(runnerDB(r.creds))
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer ro.Close()
+	id := r.action.ID
+	row, err := ro.ReadOperatorStartRecovery(ctx, id)
+	if err != nil || row.Window != store.OperatorStartMintPending {
+		t.Fatalf("durable rollback window: %s %v", row.Window, err)
+	}
+	before := activationEvidence(t, r)
+	code, output := r.command("reconcile-start", "untyped")
+	if code != exitUsage || !strings.Contains(output, "expected reconcile-start ACTION-ID") ||
+		!reflect.DeepEqual(before, activationEvidence(t, r)) {
+		t.Fatalf("invalid action identifier reached reconciliation: %d %s", code, output)
+	}
+	code, output = r.command("reconcile-start", id)
+	if code != exitUsage || !strings.Contains(output, "runner lock") || !reflect.DeepEqual(before, activationEvidence(t, r)) {
+		t.Fatalf("attended reconcile crossed live process: %d %s", code, output)
+	}
+	runner.Close()
+	// A file finalizer must not mask a missing explicit close.
+	defer debug.SetGCPercent(debug.SetGCPercent(-1))
+	terminal := r.env.terminal.openTerminal
+	r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+		return &fakeTerm{r: strings.NewReader("n\n")}, nil
+	}
+	code, output = r.command("reconcile-start", id)
+	if code != exitUsage || !reflect.DeepEqual(before, activationEvidence(t, r)) {
+		t.Fatalf("attended disposition crossed typed refusal: %d %s", code, output)
+	}
+	r.env.terminal.openTerminal = terminal
+	code, output = r.command("reconcile-start", id)
+	after := activationEvidence(t, r)
+	if code != exitOK || !strings.Contains(output, "existing evidence retained") {
+		t.Fatalf("attended disposition: %d %s", code, output)
+	}
+	before[0] = after[0]
+	if !reflect.DeepEqual(before, after) {
+		t.Fatal("disposition modified the retained worktree, branch or private files")
+	}
+	stopped, err := ro.OperatorAction(ctx, id, 0)
+	if err != nil || stopped.State != "failed" || stopped.Owner != "" {
+		t.Fatalf("rollback disposition: %s %q %v", stopped.State, stopped.Owner, err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_pair_upgrade_test.go b/cmd/ailang-worldd/operator_pair_upgrade_test.go
index c8d4960d..ab9e4971 100644
--- a/cmd/ailang-worldd/operator_pair_upgrade_test.go
+++ b/cmd/ailang-worldd/operator_pair_upgrade_test.go
@@ -32,6 +32,21 @@ func TestOperatorProcessWitness(t *testing.T) {
 	}
 	path := os.Getenv("WORLD_TEST_OPERATOR_WITNESS_PATH")
 	switch mode {
+	case "git-inherited":
+		lock := os.NewFile(4, "Git inherited runner lock")
+		defer lock.Close()
+		if got, err := lock.Stat(); err != nil {
+			t.Fatal(err)
+		} else if want, err := os.Stat(path); err != nil || !os.SameFile(got, want) {
+			t.Fatalf("Git received a different runner lock: %v", err)
+		}
+		if err := syscall.Flock(4, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
+			t.Fatal(err)
+		}
+		if err := os.WriteFile(path+".git-observed", nil, 0600); err != nil {
+			t.Fatal(err)
+		}
+		return
 	case "inherited":
 		lock := os.NewFile(3, "inherited runner lock")
 		defer lock.Close()
@@ -45,6 +60,18 @@ func TestOperatorProcessWitness(t *testing.T) {
 			t.Fatal(err)
 		}
 		return
+	case "episode":
+		lock, err := terminalEpisodeLock(path, "ep", nil)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer lock.Close()
+	case "owner":
+		lock, err := store.OpenOperatorEpisodeLock(path, "ep")
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer lock.Close()
 	default:
 		t.Fatalf("bad witness mode %q", mode)
 	}
diff --git a/cmd/ailang-worldd/operator_preparation.go b/cmd/ailang-worldd/operator_preparation.go
new file mode 100644
index 00000000..52f701fe
--- /dev/null
+++ b/cmd/ailang-worldd/operator_preparation.go
@@ -0,0 +1,210 @@
+package main
+
+import (
+	"context"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"os"
+	"path/filepath"
+	"reflect"
+	"slices"
+	"strings"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+var errOperatorDraining = errors.New("operator runner draining")
+
+func (e serviceEnv) consumeOperatorPreparation(ctx context.Context, c runnerConfig, heartbeat store.OperatorHeartbeat,
+	src operatorStartSource, stops <-chan os.Signal, process *os.File, recoveryCursor *string,
+	acquire func(string, string) (*store.OperatorEpisodeLock, error)) error {
+	select {
+	case <-stops:
+		return errOperatorDraining
+	default:
+	}
+	if err := e.recoverOperatorStart(ctx, c, heartbeat, src, recoveryCursor, stops); err != nil {
+		return err
+	}
+	return e.prepareOperatorStart(ctx, c, heartbeat.Owner, heartbeat.Build, src, stops, process, acquire)
+}
+
+func (e serviceEnv) prepareOperatorStart(ctx context.Context, c runnerConfig, owner store.OperatorOwner, build string,
+	src operatorStartSource, stops <-chan os.Signal, process *os.File,
+	acquire func(string, string) (*store.OperatorEpisodeLock, error)) error {
+	return withOperatorStartCandidate(ctx, c.Credentials, src.clock, store.OpenCredentialsReadOnly, acquire, func(held *store.OperatorEpisodeLock, row store.OperatorActionRow) error {
+		ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
+		if err != nil {
+			return err
+		}
+		defer ro.Close()
+		h, err := ro.OperatorHeartbeat(ctx)
+		if err != nil {
+			return err
+		}
+		if !preparationHeartbeat(h, owner, build, h.CatalogueGeneration, src.clock()) {
+			return store.ErrOperatorRunner
+		}
+		action, err := approval.DecodeOperatorAction(row.Action)
+		if err != nil {
+			return err
+		}
+		binary := operatorFileFact(action, "launch_binary")
+		sha, err := authority.OperatorStartLaunchDigest(action)
+		if err != nil {
+			return err
+		}
+		if !slices.Contains(c.PresetNames, action.Options.Preset) || c.Backend != "headless" ||
+			filepath.Base(binary) != "codex" || c.Executables[binary] != sha ||
+			operatorFileFact(action, "served_address") != c.Address || len(row.Challenge) != 32 {
+			return store.ErrOperatorConflict
+		}
+		if err := e.proveOldTerminalQuiet(); err != nil {
+			return err
+		}
+		newEffect := func() error {
+			select {
+			case <-stops:
+				return errOperatorDraining
+			default:
+			}
+			if serviceFileHash(binary) != sha {
+				return &approval.BindingError{Field: "launch_binary_sha256"}
+			}
+			if err := e.checkStartPair(c, src.workspace); err != nil {
+				return err
+			}
+			live, err := ro.OperatorHeartbeat(ctx)
+			if err != nil || !preparationHeartbeat(live, owner, build, h.CatalogueGeneration, src.clock()) {
+				return store.ErrOperatorRunner
+			}
+			current, err := ro.OperatorAction(ctx, row.ID, src.clock())
+			if err != nil {
+				return err
+			}
+			_, err = authority.RecheckClaimedOperatorAction(ctx, ro, action, [32]byte(row.Challenge), row.ExpiresAt, src.clock(), current)
+			return err
+		}
+		attempt := store.OperatorMintAttempt{CAS: operatorPreparationCAS(row), Runner: owner, Catalogue: h.CatalogueGeneration}
+		in, err := reconstructOperatorStart(ctx, action, attempt, src.presets, src.profiles, src.workspace, src.hosts, src.clock)
+		if err != nil {
+			return err
+		}
+		if in.Start.Agent != "codex" || operatorFileFact(action, "mcp_path") != "" {
+			return store.ErrOperatorConflict
+		}
+		if err := runnerPath(c.Home, filepath.Dir(in.TokenPath), c.UID, true, 0700); err != nil {
+			return fmt.Errorf("private session directory: %w", err)
+		}
+		// Hold lifecycle exclusion across the last stop check and claim.
+		// The exclusive episode description covers all later effects.
+		admission, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
+		if err != nil {
+			return errOperatorDraining
+		}
+		err = newEffect()
+		if err == nil {
+			err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
+				fresh, err := db.OperatorHeartbeat(ctx)
+				if err != nil || !preparationHeartbeat(fresh, owner, build, h.CatalogueGeneration, src.clock()) {
+					return store.ErrOperatorRunner
+				}
+				return db.ClaimOperatorAction(ctx, attempt.CAS, owner, src.clock())
+			})
+		}
+		admission.Close()
+		if err != nil {
+			return err
+		}
+		attempt.CAS.Generation++
+		checkpoint := func(name string) error {
+			if strings.HasPrefix(name, "before ") && name != "before readback" && name != "before observe" {
+				return newEffect()
+			}
+			return nil
+		}
+		if err := reserveOperatorStartWorktree(ctx, c.Credentials, held, action, attempt, src, checkpoint, process); err != nil {
+			return err
+		}
+		if err := newEffect(); err != nil {
+			return fmt.Errorf("%w; %s", err, reconcileStartCommand(row.ID))
+		}
+		row, err = ro.OperatorAction(ctx, row.ID, src.clock())
+		if err != nil {
+			return err
+		}
+		in.Attempt.CAS, in.Attempt.StepVersion = operatorPreparationCAS(row), 1
+		var token string
+		if err := runnerWrite(c.Credentials, func(db *store.Credentials) error {
+			key := store.OperatorStepKey{Action: in.Attempt.CAS, Owner: row.Owner, Kind: "mint", Identity: row.ID, Version: 1}
+			if err := db.ReserveOperatorStep(ctx, key, src.clock()); err != nil {
+				return err
+			}
+			var err error
+			token, err = authority.MintOperatorStart(ctx, db, ro, in, action, row, src.clock)
+			return err
+		}); err != nil {
+			return fmt.Errorf("%w; %s", err, reconcileStartCommand(row.ID))
+		}
+		if err := publishOperatorStartFiles(ctx, c, held, row.ID, token, src.clock, checkpoint); err != nil {
+			return fmt.Errorf("%w; %s", err, reconcileStartCommand(row.ID))
+		}
+		prepared, err := ro.ReadOperatorStartRecovery(ctx, row.ID)
+		if err != nil || prepared.Window != store.OperatorStartPrepared {
+			return fmt.Errorf("%w; %s", store.ErrOperatorConflict, reconcileStartCommand(row.ID))
+		}
+		return nil
+	})
+}
+
+func operatorPreparationCAS(row store.OperatorActionRow) store.OperatorCAS {
+	return store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
+}
+
+func preparationHeartbeat(h store.OperatorHeartbeat, owner store.OperatorOwner, build string, catalogue, now int64) bool {
+	return h.Owner == owner && h.CatalogueGeneration == catalogue && h.Available(build, "start", now)
+}
+
+func (e serviceEnv) checkStartPair(c runnerConfig, workspace string) error {
+	installed, _, err := e.installedRunner()
+	if err != nil || !reflect.DeepEqual(installed, c) {
+		return store.ErrOperatorConflict
+	}
+	if err := e.checkRunnerPair(c); err != nil {
+		return err
+	}
+	r, err := readServiceRecord(e.paths)
+	if err != nil {
+		return err
+	}
+	physical, err := filepath.EvalSymlinks(r.Workspace)
+	current, currentErr := filepath.EvalSymlinks(workspace)
+	if err != nil || currentErr != nil || physical != current {
+		return store.ErrOperatorConflict
+	}
+	return nil
+}
+
+func reconcileStartCommand(id string) string {
+	return "ailang-worldd service runner reconcile-start " + id
+}
+
+func recoveryWorktree(ctx context.Context, current store.OperatorStartRecovery, a approval.OperatorAction) (operatorWorktreeEvidence, error) {
+	var target operatorWorktreeTarget
+	path, err := operatorPhysicalTarget(operatorWorktreeFact(a))
+	if err != nil || json.Unmarshal([]byte(current.Worktree), &target) != nil || target.Action != a.ID || target.Path != path {
+		return operatorWorktreeEvidence{}, store.ErrOperatorConflict
+	}
+	git, err := sessionGit()
+	if err != nil {
+		return operatorWorktreeEvidence{}, err
+	}
+	seen, err := target.inspect(ctx, git)
+	if err != nil || seen.state != treeExact {
+		return operatorWorktreeEvidence{}, store.ErrOperatorConflict
+	}
+	return seen, nil
+}
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
index ac031e5f..1ea8ea43 100644
--- a/cmd/ailang-worldd/operator_runner.go
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -12,6 +12,7 @@ import (
 	"syscall"
 	"time"
 
+	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/daemon"
 	"github.com/sunholo-data/ailang-world/host/store"
 )
@@ -65,12 +66,51 @@ func runOperatorRunner(args []string, errOut io.Writer, e serviceEnv) int {
 	defer signal.Stop(stops)
 	ticker := time.NewTicker(runnerPoll)
 	defer ticker.Stop()
+	source := operatorStartSource{defaultPresetsDir(), defaultExecProfilesDir(), "", defaultHostsDir(), e.terminal.now}
+	return e.runOperatorPreparationLoop(errOut, c, h, source, stops, ticker.C, lock)
+}
+
+func (e serviceEnv) runOperatorPreparationLoop(errOut io.Writer, c runnerConfig, h store.OperatorHeartbeat,
+	src operatorStartSource, stops <-chan os.Signal, ticks <-chan time.Time, lock *os.File) int {
+	var recoveryCursor string
 	for {
+		select {
+		case <-stops:
+			return 0
+		default:
+		}
+
 		if err := e.runnerTick(c, &h); err != nil {
 			return serviceMessage(errOut, "runner halted", "pair or heartbeat unavailable", "ailang-worldd service runner status", 2)
 		}
+		record, err := readServiceRecord(e.paths)
+		if err != nil {
+			return serviceMessage(errOut, "runner halted", "service record unavailable", "ailang-worldd service runner status", 2)
+		}
+		ctx, cancel := budgetContext(22 * time.Second)
+		src.workspace = record.Workspace
+		err = e.consumeOperatorPreparation(ctx, c, h, src, stops, lock, &recoveryCursor, store.OpenOperatorEpisodeLock)
+		cancel()
+		switch {
+		case errors.Is(err, errOperatorDraining):
+			return 0
+		case err != nil && !errors.Is(err, store.ErrOperatorMissing):
+			var stale *approval.BindingError
+			var remedy operatorRecoveryNeeded
+			switch {
+			case errors.As(err, &stale):
+				fmt.Fprintf(errOut, "operator start binding stale (%s); request fresh review\n", stale.Field)
+			case errors.As(err, &remedy):
+				fmt.Fprintf(errOut, "operator start blocked; attended recovery: %s\n", remedy.command)
+			default:
+				fmt.Fprintln(errOut, "operator start blocked; inspect service runner status and reconcile durable evidence")
+			}
+		}
 		select {
-		case <-ticker.C:
+		case _, open := <-ticks:
+			if !open {
+				return 0
+			}
 		case <-stops:
 			return 0
 		}
@@ -86,7 +126,7 @@ func (e serviceEnv) runnerTick(c runnerConfig, h *store.OperatorHeartbeat) error
 	}
 	h.SeenAt = e.terminal.now()
 	h.Readiness = "start"
-	// S4 projects the configured names. No effect dispatch or GitHub probe.
+	// Publish the admitted installed preset names before the start consumer.
 	h.Catalogue, _ = json.Marshal(c.PresetNames)
 	return runnerWrite(c.Credentials, func(db *store.Credentials) error {
 		ctx, cancel := budgetContext(time.Second)
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index 69cd515b..ebc2aa04 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -21,7 +21,7 @@ import (
 
 const runnerLabel = "com.sunholo.ailang-world.operator-runner"
 const runnerSystemPath = "/usr/bin:/bin:/usr/sbin:/sbin"
-const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | reconcile-files ACTION-ID | uninstall
+const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | reconcile-files ACTION-ID | reconcile-start ACTION-ID | uninstall
 Install, stop, cancel, orphan mint revocation, worktree and file disposition and uninstall require typed y. Start uses the approved configuration.
 Config JSON keys: home, uid, address, world, credentials, executables (absolute path to sha256 digest), preset_names, tool_dirs, backend ("headless").
 The World service must already be running. Use absolute, canonical paths and include its installed binary in executables.
@@ -272,7 +272,7 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 		if len(args) != 3 || args[1] != "--config" {
 			return refuse(errors.New("expected install --config FILE"))
 		}
-	case "cancel", "reconcile-mint", "reconcile-worktree", "reconcile-files":
+	case "cancel", "reconcile-mint", "reconcile-worktree", "reconcile-files", "reconcile-start":
 		if len(args) != 2 || !verbs.ValidOperatorActionID(args[1]) {
 			return refuse(fmt.Errorf("expected %s ACTION-ID", verb))
 		}
@@ -289,6 +289,9 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 	if verb == "reconcile-worktree" {
 		return e.reconcileRunnerWorktree(args[1], out, errOut, readRunnerWorktree)
 	}
+	if verb == "reconcile-start" {
+		return e.reconcileOperatorStart(args[1], out, errOut)
+	}
 	if verb == "reconcile-files" {
 		return e.reconcileRunnerFiles(args[1], out, errOut)
 	}
diff --git a/cmd/ailang-worldd/operator_start_input.go b/cmd/ailang-worldd/operator_start_input.go
index f83eb66a..88d12a2f 100644
--- a/cmd/ailang-worldd/operator_start_input.go
+++ b/cmd/ailang-worldd/operator_start_input.go
@@ -84,8 +84,8 @@ func operatorStartFacts(ctx context.Context, signed approval.OperatorAction, att
 
 // reconstructOperatorStart requires every signed resolution fact to equal its
 // live value before any effect. A change, even one lowering authority, needs a
-// fresh review. It locks, writes, mints and launches nothing, and it has no
-// production caller until the paired-process consumer lands.
+// fresh review. It is read-only and runs inside the preparation consumer's
+// process and episode exclusion.
 func reconstructOperatorStart(ctx context.Context, signed approval.OperatorAction, attempt store.OperatorMintAttempt,
 	presets, profiles, workspace, hosts string, clock func() int64) (store.OperatorMintInput, error) {
 	in, facts, err := operatorStartFacts(ctx, signed, attempt, presets, profiles, workspace, hosts, clock, resolvePresetStart)
diff --git a/cmd/ailang-worldd/operator_start_recovery.go b/cmd/ailang-worldd/operator_start_recovery.go
new file mode 100644
index 00000000..8b7c8950
--- /dev/null
+++ b/cmd/ailang-worldd/operator_start_recovery.go
@@ -0,0 +1,299 @@
+package main
+
+import (
+	"context"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"os"
+	"path/filepath"
+	"strconv"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+type operatorRecoveryNeeded struct{ command string }
+
+func (r operatorRecoveryNeeded) Error() string { return r.command }
+
+func (e serviceEnv) recoverOperatorStart(ctx context.Context, c runnerConfig, h store.OperatorHeartbeat, src operatorStartSource, cursor *string, stops <-chan os.Signal) error {
+	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
+	if err != nil {
+		return err
+	}
+	defer ro.Close()
+	row, err := ro.NextOwnedOperatorStart(ctx, *cursor)
+	if errors.Is(err, store.ErrOperatorMissing) {
+		*cursor = ""
+		return nil
+	}
+	if err != nil {
+		return err
+	}
+	*cursor = row.ID
+	return e.recoverOwnedStart(ctx, ro, c, h, src, row, stops)
+}
+
+func (e serviceEnv) recoverOwnedStart(ctx context.Context, ro *store.Credentials, c runnerConfig, h store.OperatorHeartbeat,
+	src operatorStartSource, row store.OperatorActionRow, stops <-chan os.Signal) error {
+	if !verbs.EpisodeIDPattern.MatchString(row.Episode) {
+		return store.ErrOperatorConflict
+	}
+	lock, err := store.OpenOperatorEpisodeLock(c.Credentials, row.Episode)
+	if err != nil {
+		return fmt.Errorf("%w: live episode owner or effect", store.ErrOperatorConflict)
+	}
+	defer lock.Close()
+	current, err := ro.ReadOperatorStartRecovery(ctx, row.ID)
+	if err != nil || current.CAS != operatorPreparationCAS(row) || current.Owner != row.Owner || current.Episode != row.Episode {
+		return store.ErrOperatorConflict
+	}
+	a, err := approval.DecodeOperatorAction(row.Action)
+	if err != nil || a.Store != c.Credentials {
+		return store.ErrOperatorConflict
+	}
+	// A heartbeat can make an old review ineligible. It cannot prove an
+	// episode process has exited. The runner and episode locks cover this.
+	switch current.Window {
+	case store.OperatorStartGitPending:
+		return operatorRecoveryNeeded{reconcileWorktreeCommand(row.ID)}
+	case store.OperatorStartMintOnly:
+		return operatorRecoveryNeeded{"ailang-worldd service runner reconcile-mint " + row.ID}
+	case store.OperatorStartFilesPending:
+		return operatorRecoveryNeeded{reconcileFilesCommand(row.ID)}
+	case store.OperatorStartTreeOnly, store.OperatorStartMintPending, store.OperatorStartMintUncertain:
+		return operatorRecoveryNeeded{reconcileStartCommand(row.ID)}
+	case store.OperatorStartPrepared:
+		return verifyPreparedStart(ctx, c, current, ro, a, src.clock())
+	case store.OperatorStartUntouched:
+		if err := e.proveOldTerminalQuiet(); err != nil {
+			return err
+		}
+		catalogue, err := strconv.ParseInt(operatorFileFact(a, "mint_catalogue"), 10, 64)
+		if err != nil {
+			return &approval.BindingError{Field: "mint_catalogue"}
+		}
+		live, err := ro.OperatorHeartbeat(ctx)
+		if err != nil || !preparationHeartbeat(live, h.Owner, h.Build, catalogue, src.clock()) {
+			return &approval.BindingError{Field: "mint_catalogue"}
+		}
+		if err := e.checkStartPair(c, src.workspace); err != nil {
+			return err
+		}
+		if _, err := authority.OperatorStartLaunchDigest(a); err != nil {
+			return err
+		}
+		p, err := authority.RecheckClaimedOperatorAction(ctx, ro, a, [32]byte(row.Challenge), row.ExpiresAt, src.clock(), row)
+		if err != nil {
+			return err
+		}
+		attempt := store.OperatorMintAttempt{CAS: current.CAS, Runner: h.Owner, Catalogue: catalogue}
+		in, err := reconstructOperatorStart(ctx, a, attempt, src.presets, src.profiles, src.workspace, src.hosts, src.clock)
+		if err != nil {
+			return err
+		}
+		git, err := sessionGit()
+		if err != nil {
+			return err
+		}
+		target, err := newOperatorWorktreeTarget(ctx, git, a, in.Start)
+		if err != nil {
+			return err
+		}
+		tree, err := target.inspect(ctx, git)
+		if err != nil || tree.state != treeAbsent {
+			return store.ErrOperatorConflict
+		}
+		for _, path := range []string{in.TokenPath, in.Start.Task.Path, operatorFileFact(a, "mcp_path")} {
+			if path != "" {
+				if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
+					return store.ErrOperatorConflict
+				}
+			}
+		}
+		admit, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
+		if err != nil {
+			return errOperatorDraining
+		}
+		defer admit.Close()
+		select {
+		case <-stops:
+			return errOperatorDraining
+		default:
+		}
+		return runnerWrite(c.Credentials, func(db *store.Credentials) error {
+			live, err := db.OperatorHeartbeat(ctx)
+			if err != nil || !preparationHeartbeat(live, h.Owner, h.Build, catalogue, src.clock()) {
+				return store.ErrOperatorRunner
+			}
+			return db.DisposeOperatorStartGap(ctx, current, p.OperatorConfirmation().MintProof(), src.clock())
+		})
+	}
+	return store.ErrOperatorConflict
+}
+
+func verifyPreparedStart(ctx context.Context, c runnerConfig, current store.OperatorStartRecovery,
+	ro *store.Credentials, a approval.OperatorAction, now int64) error {
+	s, _, err := ro.ReadOperatorFiles(ctx, a.ID)
+	if err != nil || s.Phase != "observed" || s.Mint.CAS != current.CAS || s.Worktree != current.Worktree ||
+		now >= s.Session.ExpiresAt || store.ValidateOperatorStartSession(s, a) != nil {
+		return store.ErrOperatorConflict
+	}
+	if _, err := recoveryWorktree(ctx, current, a); err != nil {
+		return err
+	}
+	files, err := readStartRecoveryFiles(c, s, a)
+	if err != nil || files.Files[0].Status != "exact" || files.Files[1].Status != "exact" {
+		return store.ErrOperatorConflict
+	}
+	return nil
+}
+
+func readStartRecoveryFiles(c runnerConfig, s store.OperatorFileState, a approval.OperatorAction) (operatorFilesProof, error) {
+	want, start, err := operatorFilePaths(s, a)
+	var recorded operatorFileSet
+	if err != nil || json.Unmarshal([]byte(s.Step.Identity), &recorded) != nil || recorded.Inode == 0 {
+		return operatorFilesProof{}, store.ErrOperatorConflict
+	}
+	want.Device, want.Inode = recorded.Device, recorded.Inode
+	if recorded != want {
+		return operatorFilesProof{}, store.ErrOperatorConflict
+	}
+	addr, err := operatorFileAddress(start, operatorFileFact(a, "served_address"), c.Address)
+	if err != nil {
+		return operatorFilesProof{}, err
+	}
+	files := readOperatorFileSet(recorded, addr)
+	if files.Dir.Status != "private directory" {
+		return operatorFilesProof{}, store.ErrOperatorConflict
+	}
+	return files, nil
+}
+
+type operatorRecoveryObserved struct {
+	tree    operatorWorktreeEvidence
+	files   operatorFilesProof
+	session store.OperatorFileState
+}
+
+func (e serviceEnv) reconcileOperatorStart(id string, out, stderr io.Writer) int {
+	refuse := func(err error) int {
+		return serviceMessage(stderr, "start recovery refused", err.Error(), "stop the companion and inspect durable evidence before retrying "+reconcileStartCommand(id), exitUsage)
+	}
+	lifecycle, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
+	if err != nil {
+		return refuse(err)
+	}
+	defer lifecycle.Close()
+	c, _, err := e.installedRunner()
+	if err != nil {
+		return refuse(err)
+	}
+	if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) {
+		return refuse(errors.New("companion must be unloaded"))
+	}
+	process, err := e.waitRunnerExit(c)
+	if err != nil {
+		return refuse(err)
+	}
+	defer process.Close()
+	if err := e.proveOldTerminalQuiet(); err != nil {
+		return refuse(err)
+	}
+	ctx, cancel := budgetContext(5 * time.Second)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
+	if err != nil {
+		return refuse(err)
+	}
+	defer ro.Close()
+	row, err := ro.OperatorAction(ctx, id, 0)
+	if err != nil || !verbs.EpisodeIDPattern.MatchString(row.Episode) {
+		return refuse(store.ErrOperatorConflict)
+	}
+	episode, err := store.OpenOperatorEpisodeLock(c.Credentials, row.Episode)
+	if err != nil {
+		return refuse(err)
+	}
+	defer episode.Close()
+	reviewed, err := ro.ReadOperatorStartRecovery(ctx, id)
+	if err != nil || reviewed.CAS != operatorPreparationCAS(row) || reviewed.Owner != row.Owner {
+		return refuse(store.ErrOperatorConflict)
+	}
+	a, err := approval.DecodeOperatorAction(row.Action)
+	if err != nil || a.Store != c.Credentials {
+		return refuse(store.ErrOperatorConflict)
+	}
+	check := func(ctx context.Context) (operatorRecoveryObserved, error) {
+		var seen operatorRecoveryObserved
+		switch reviewed.Window {
+		case store.OperatorStartUntouched:
+			return seen, nil
+		case store.OperatorStartTreeOnly, store.OperatorStartMintPending, store.OperatorStartMintUncertain, store.OperatorStartFilesPending:
+			tree, err := recoveryWorktree(ctx, reviewed, a)
+			if err != nil {
+				return seen, err
+			}
+			seen.tree = tree
+			if reviewed.Window != store.OperatorStartFilesPending {
+				return seen, nil
+			}
+			seen.session, _, err = ro.ReadOperatorFiles(ctx, id)
+			if err != nil {
+				return seen, err
+			}
+			seen.files, err = readStartRecoveryFiles(c, seen.session, a)
+			if err != nil || seen.files.Files[1].Status != "absent" ||
+				(seen.files.Files[0].Status != "absent" && seen.files.Files[0].Status != "exact") {
+				return seen, store.ErrOperatorConflict
+			}
+			return seen, nil
+		default:
+			return seen, store.ErrOperatorConflict
+		}
+	}
+	before, err := check(ctx)
+	if err != nil {
+		return refuse(err)
+	}
+	cancel()
+	change := fmt.Sprintf("Stop action %s (%s) and retain its existing evidence?", id, reviewed.Window)
+	if reviewed.Window == store.OperatorStartFilesPending {
+		change = fmt.Sprintf("Stop %s and revoke orphan credential %s, preserving all files and the worktree?", id, before.session.Session.CredentialID)
+	}
+	if reviewed.Window == store.OperatorStartMintUncertain {
+		change = fmt.Sprintf("Mark %s uncertain and retain its owner and all evidence?", id)
+	}
+	if !e.confirm(stderr, change) {
+		return exitUsage
+	}
+	ctx, done := budgetContext(5 * time.Second)
+	defer done()
+	after, err := check(ctx)
+	if err != nil || before != after {
+		return refuse(store.ErrOperatorConflict)
+	}
+	if err := runnerWrite(c.Credentials, func(db *store.Credentials) error {
+		switch reviewed.Window {
+		case store.OperatorStartFilesPending:
+			return db.RecordOperatorFiles(ctx, before.session, store.OperatorFilesStoppedOrphan, nil, e.terminal.now)
+		case store.OperatorStartMintUncertain:
+			return db.MarkOperatorMintUncertain(ctx, reviewed)
+		default:
+			return db.DisposeOperatorStartGap(ctx, reviewed, nil, e.terminal.now())
+		}
+	}); err != nil {
+		return refuse(err)
+	}
+	if reviewed.Window == store.OperatorStartMintUncertain {
+		fmt.Fprintln(out, "Mint commit uncertain; owner retained:", id)
+		return exitOK
+	}
+	fmt.Fprintln(out, "Start action stopped; existing evidence retained:", id)
+	return exitOK
+}
diff --git a/cmd/ailang-worldd/operator_worktree.go b/cmd/ailang-worldd/operator_worktree.go
index 99603326..28bc8ef0 100644
--- a/cmd/ailang-worldd/operator_worktree.go
+++ b/cmd/ailang-worldd/operator_worktree.go
@@ -28,7 +28,7 @@ type operatorStartSource struct {
 
 // reserveOperatorStartWorktree borrows the caller-owned episode lock.
 func reserveOperatorStartWorktree(ctx context.Context, credentials string, held *store.OperatorEpisodeLock, signed approval.OperatorAction,
-	attempt store.OperatorMintAttempt, src operatorStartSource, checkpoint func(string) error) error {
+	attempt store.OperatorMintAttempt, src operatorStartSource, checkpoint func(string) error, process ...*os.File) error {
 	canonical, err := store.CanonicalPath(credentials)
 	if err != nil {
 		return err
@@ -88,7 +88,7 @@ func reserveOperatorStartWorktree(ctx context.Context, credentials string, held
 	if err != nil {
 		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
 	}
-	addErr := worktree.Add(ctx, git, target.Repo, target.Path, target.Branch, target.SHA, fd)
+	addErr := worktree.Add(ctx, git, target.Repo, target.Path, target.Branch, target.SHA, append([]*os.File{fd}, process...)...)
 	if _, err := held.Verify(canonical, signed.Episode); err != nil {
 		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
 	}
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 78f366a5..238c746e 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -121,6 +121,28 @@ ailang-worldd service runner reconcile-files oa-0123456789abcdef0123456789abcdef
 ```
 
 
+Create the approved preset’s session directory with mode 0700 before admitting a start.
+The companion prepares one confirmed headless start per tick. A prepared start retains its
+worktree, credential and private files while awaiting the launch transition. For a claimed start
+with a worktree-only or incomplete mint checkpoint, stop the companion and use the signed
+action's ID at the terminal. Review the proposed disposition and type `y`.
+
+```sh
+ailang-worldd service runner stop
+ailang-worldd service runner reconcile-start oa-0123456789abcdef0123456789abcdef
+```
+
+The runner log names the applicable attended command for a retained reservation. Successful
+preparation is durable as `prepared-awaiting-S6` and leaves the status readiness fields intact.
+An ambiguous mint retains the owner and every receipt. A file reservation with a verified missing
+token may be stopped by revoking its exact orphan credential. A changed binary or an older review
+without its digest requires a fresh attended review. Try the command with a fake supervisor and
+a disposable workspace using the following test.
+
+```sh
+AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang" go test ./cmd/ailang-worldd -run '^TestOperatorAttendedStartDisposition$' -count=1 -v
+```
+
 ## Run your agents through World
 
 World launchers now tell Claude Code, Codex and pi their episode, worktree, session
diff --git a/host/store/operator_files.go b/host/store/operator_files.go
index 786bc0f8..699b65cf 100644
--- a/host/store/operator_files.go
+++ b/host/store/operator_files.go
@@ -81,6 +81,7 @@ const (
 	OperatorFilesPublished       = "published"
 	OperatorFilesStoppedAbsent   = "stopped: absent"
 	OperatorFilesStoppedRetained = "stopped: retained"
+	OperatorFilesStoppedOrphan   = "stopped: orphan revoked"
 )
 
 // RecordOperatorFiles reserves an empty phase, or observes or disposes the exact reservation.
@@ -122,7 +123,7 @@ func (c *Credentials) RecordOperatorFiles(ctx context.Context, reviewed Operator
 		}
 	case "reserved":
 		switch value {
-		case OperatorFilesPublished, OperatorFilesStoppedAbsent, OperatorFilesStoppedRetained:
+		case OperatorFilesPublished, OperatorFilesStoppedAbsent, OperatorFilesStoppedRetained, OperatorFilesStoppedOrphan:
 			if err := observeOperatorStep(ctx, tx, key, []byte(value)); err != nil {
 				return err
 			}
@@ -164,7 +165,7 @@ func (c *Credentials) RecordOperatorFiles(ctx context.Context, reviewed Operator
 			return err
 		}
 	}
-	stopped := current.Phase == "reserved" && (value == OperatorFilesStoppedAbsent || value == OperatorFilesStoppedRetained)
+	stopped := current.Phase == "reserved" && (value == OperatorFilesStoppedAbsent || value == OperatorFilesStoppedRetained || value == OperatorFilesStoppedOrphan)
 	uncertain := current.Phase == "reserved" && value == OperatorFilesUncertain
 	res, err := tx.ExecContext(ctx, `UPDATE operator_actions
 		SET state=CASE WHEN ? THEN 'failed' WHEN ? THEN 'uncertain' ELSE state END,
@@ -181,6 +182,13 @@ func (c *Credentials) RecordOperatorFiles(ctx context.Context, reviewed Operator
 	if err := operatorChangedOne(res, err); err != nil {
 		return err
 	}
+	if stopped && value == OperatorFilesStoppedOrphan {
+		res, err := tx.ExecContext(ctx, `DELETE FROM session_credentials WHERE credential_id=? AND episode_id=?`,
+			current.Session.CredentialID, current.Session.EpisodeID)
+		if err := operatorChangedOne(res, err); err != nil {
+			return err
+		}
+	}
 	return tx.Commit()
 }
 
@@ -207,6 +215,12 @@ func admitOperatorFiles(ctx context.Context, tx *sql.Tx, row OperatorActionRow,
 	if err := operatorMintApprover(ctx, tx, proof, b, now); err != nil {
 		return err
 	}
+	return ValidateOperatorStartSession(current, action)
+}
+
+// ValidateOperatorStartSession compares the current persisted session with
+// the session implied by the frozen signed start and its exact mint time.
+func ValidateOperatorStartSession(current OperatorFileState, action approval.OperatorAction) error {
 	catalogue, err := strconv.ParseInt(mintFact(action, "mint_catalogue"), 10, 64)
 	if err != nil {
 		return ErrOperatorConflict
diff --git a/host/store/operator_start_recovery.go b/host/store/operator_start_recovery.go
index 13127149..e051ab47 100644
--- a/host/store/operator_start_recovery.go
+++ b/host/store/operator_start_recovery.go
@@ -195,3 +195,40 @@ func (c *Credentials) DisposeOperatorStartGap(ctx context.Context, reviewed Oper
 	}
 	return tx.Commit()
 }
+
+// NextOwnedOperatorStart scans only owned starts. A cursor lets the runner
+// inspect one action per tick without starving behind a prepared checkpoint.
+func (c *Credentials) NextOwnedOperatorStart(ctx context.Context, after string) (OperatorActionRow, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorActionRow{}, err
+	}
+	return scanOperatorAction(c.db.QueryRowContext(ctx, `SELECT `+operatorColumns+
+		` FROM operator_actions WHERE kind='start' AND owner IS NOT NULL AND id>? ORDER BY id LIMIT 1`, after))
+}
+
+// MarkOperatorMintUncertain preserves every receipt and the old owner while
+// recording the ambiguity found after an interrupted mint.
+func (c *Credentials) MarkOperatorMintUncertain(ctx context.Context, reviewed OperatorStartRecovery) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	current, row, err := readOperatorStartRecovery(ctx, tx, reviewed.CAS.ID)
+	if err != nil || current != reviewed || current.Window != OperatorStartMintUncertain {
+		return ErrOperatorConflict
+	}
+	if row.State == "uncertain" {
+		return nil
+	}
+	res, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state='uncertain', generation=generation+1
+		WHERE id=? AND digest=? AND generation=? AND owner=? AND state='running'`,
+		current.CAS.ID, current.CAS.Digest[:], current.CAS.Generation, current.Owner)
+	if err := operatorChangedOne(res, err); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
diff --git a/host/store/operator_start_recovery_test.go b/host/store/operator_start_recovery_test.go
index bb819324..aa5f8e16 100644
--- a/host/store/operator_start_recovery_test.go
+++ b/host/store/operator_start_recovery_test.go
@@ -2,6 +2,7 @@ package store
 
 import (
 	"errors"
+	"os"
 	"reflect"
 	"testing"
 )
@@ -162,3 +163,118 @@ func TestOperatorStartDispositionRejectsBlobOwner(t *testing.T) {
 		t.Fatalf("disposition with matching TEXT owner: %v", err)
 	}
 }
+
+func TestOperatorPartialFilesOrphanRevoke(t *testing.T) {
+	f := newOperatorPinFixture(t, false)
+	f.mint(t)
+	id := f.in.Attempt.CAS.ID
+	state, _, err := f.c.ReadOperatorFiles(f.ctx, id)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := f.c.RecordOperatorFiles(f.ctx, state, "file-identity", f.proof, func() int64 { return pinNow }); err != nil {
+		t.Fatal(err)
+	}
+	state, _, err = f.c.ReadOperatorFiles(f.ctx, id)
+	if err != nil {
+		t.Fatal(err)
+	}
+	other := SessionRow{CredentialID: hex64("other"), EpisodeID: f.in.Start.Episode,
+		GrantsJSON: "[]", ExpiresAt: pinNow + 1, CreatedAt: pinNow}
+	if err := f.c.MintSession(f.ctx, other); err != nil {
+		t.Fatal(err)
+	}
+	before := f.snapshot(t)
+	stale := state
+	stale.Step.Version++
+	f.assertRefused(t, before, f.c.RecordOperatorFiles(f.ctx, stale, OperatorFilesStoppedOrphan, nil, func() int64 { return pinNow }), ErrOperatorConflict)
+	operatorExec(t, f.c, `CREATE TRIGGER refuse_orphan_delete BEFORE DELETE ON session_credentials
+		WHEN OLD.credential_id != '`+other.CredentialID+`' BEGIN SELECT RAISE(IGNORE); END`)
+	f.assertRefused(t, before, f.c.RecordOperatorFiles(f.ctx, state, OperatorFilesStoppedOrphan, nil, func() int64 { return pinNow }), ErrOperatorConflict)
+	operatorExec(t, f.c, "DROP TRIGGER refuse_orphan_delete")
+	operatorExec(t, f.c, `CREATE TRIGGER move_orphan_episode AFTER UPDATE OF state ON operator_actions
+		WHEN NEW.state='failed' BEGIN UPDATE session_credentials SET episode_id='other'
+		WHERE credential_id <> '`+other.CredentialID+`'; END`)
+	f.assertRefused(t, before, f.c.RecordOperatorFiles(f.ctx, state, OperatorFilesStoppedOrphan, nil, func() int64 { return pinNow }), ErrOperatorConflict)
+	operatorExec(t, f.c, "DROP TRIGGER move_orphan_episode")
+	if err := f.c.RecordOperatorFiles(f.ctx, state, OperatorFilesStoppedOrphan, nil, func() int64 { return pinNow }); err != nil {
+		t.Fatal(err)
+	}
+	after := f.snapshot(t)
+	if after.action.State != "failed" || after.action.Owner != "" || len(after.sessions) != 1 ||
+		after.sessions[0].CredentialID != hex64("other") || !reflect.DeepEqual(before.start, after.start) {
+		t.Fatalf("orphan disposition did not preserve unrelated credential and exact start: %+v", after)
+	}
+}
+
+func TestNextOwnedOperatorStart(t *testing.T) {
+	c, ctx := openTestCredentials(t), testCtx(t)
+	add := func(n int, episode, kind string, owner bool) OperatorActionRow {
+		t.Helper()
+		req := operatorRequest(t, c, n, episode, kind)
+		row := putOperator(t, c, req, req.CreatedAt)
+		if owner {
+			signOperatorPinActionAt(t, c, row, req.Action, req.CreatedAt)
+			operatorExec(t, c, "UPDATE operator_actions SET state='claimed', owner='runner' WHERE id=?", row.ID)
+		}
+		return getOperator(t, c, row.ID, req.CreatedAt)
+	}
+	last := add(4, "last", "start", true)
+	first := add(3, "first", "start", true)
+	add(2, "unowned", "start", false)
+	add(1, "finish", "finish", true)
+	before, err := os.ReadFile(c.lock.dbPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	for _, run := range []struct{ after, want string }{{"", first.ID}, {first.ID, last.ID}, {last.ID, ""}} {
+		got, err := c.NextOwnedOperatorStart(ctx, run.after)
+		if got.ID != run.want || (run.want == "" && !errors.Is(err, ErrOperatorMissing)) || (run.want != "" && err != nil) {
+			t.Fatalf("owned start after %q: got %q (%v), want %q", run.after, got.ID, err, run.want)
+		}
+		assertGCLedgerBytes(t, c.lock.dbPath, before)
+	}
+}
+
+func TestOperatorMarkMintUncertain(t *testing.T) {
+	f := newOperatorPinFixture(t, false)
+	id := f.in.Attempt.CAS.ID
+	tree := f.reserveWorktree(t)
+	if err := f.c.ObserveOperatorStep(f.ctx, tree, []byte("created")); err != nil {
+		t.Fatal(err)
+	}
+	key := OperatorStepKey{tree.Action, tree.Owner, "mint", id, 1}
+	if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
+		t.Fatal(err)
+	}
+	pending, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
+	if err != nil || pending.Window != OperatorStartMintPending {
+		t.Fatalf("pending mint: %+v %v", pending, err)
+	}
+	f.assertRefused(t, f.snapshot(t), f.c.MarkOperatorMintUncertain(f.ctx, pending), ErrOperatorConflict)
+	if err := f.c.MintSession(f.ctx, SessionRow{CredentialID: f.in.CredentialID, EpisodeID: f.in.Start.Episode,
+		GrantsJSON: "[]", ExpiresAt: pinNow + 1, CreatedAt: pinNow}); err != nil {
+		t.Fatal(err)
+	}
+	reviewed, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
+	if err != nil || reviewed.Window != OperatorStartMintUncertain {
+		t.Fatalf("uncertain mint: %+v %v", reviewed, err)
+	}
+	before := f.snapshot(t)
+	f.assertRefused(t, before, f.c.MarkOperatorMintUncertain(f.ctx, pending), ErrOperatorConflict)
+	if err := f.c.MarkOperatorMintUncertain(f.ctx, reviewed); err != nil {
+		t.Fatal(err)
+	}
+	after := f.snapshot(t)
+	if after.action.State != "uncertain" || after.action.Owner != before.action.Owner || after.action.Generation != before.action.Generation+1 ||
+		before.steps != after.steps || !reflect.DeepEqual(before.sessions, after.sessions) || !reflect.DeepEqual(before.start, after.start) {
+		t.Fatalf("uncertain state lost evidence: %+v", after)
+	}
+	current, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := f.c.MarkOperatorMintUncertain(f.ctx, current); err != nil || !reflect.DeepEqual(after, f.snapshot(t)) {
+		t.Fatalf("uncertain disposition replay changed evidence: %v", err)
+	}
+}
diff --git a/host/verbs/operator_policy_test.go b/host/verbs/operator_policy_test.go
index f935eff8..f81eb26f 100644
--- a/host/verbs/operator_policy_test.go
+++ b/host/verbs/operator_policy_test.go
@@ -69,6 +69,7 @@ func TestOperatorGrammarPolicy(t *testing.T) {
 func TestOperatorNeverDispatch(t *testing.T) {
 	for _, line := range []string{"approver enroll --name mark-key", "approver revoke mark-key", "session mint w1", "session revoke w1",
 		"service install", "service", "serve", "setup", "ailang-worldd approver enroll", "ailang-worldd service uninstall",
+		"service runner reconcile-start oa-0123456789abcdef0123456789abcdef",
 		"world-publish approve sha256:x", "world-publish"} {
 		a := admitLine(t, line, allYes)
 		refused(t, a, line, OperatorNever, "-", "is never run from the app")
````

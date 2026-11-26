# Row 243 V4 S5f3 — operator start activation prerequisites

Iteration 441. Planner and executor lane `codex:gpt-6.2-sol`. Base `2f573488043565f480574ed42d1dacd514c6aa0b`. Apply M1, M2 and M3 in order. Appendix C measures the final sprint to the first caller.

## 0. Binding constraints and cut

**Outcome B.** The complete scratch build at `$S/st-M4` measured **2302 changed lines** from `2f573488`, against the 1,650-line complete-unit ceiling. The earlier integration already measured 1,884 before the verified orphan and incomplete-mint dispositions and the final paired-drain checks. S5f3 takes all three prerequisites in priority order for **1385 changed lines**. Appendix C assigns the production caller to exactly one subsequent sprint, S5f4, measured as **917 changed lines** on top of the cumulative M3 tree. S5f4 ends at prepared-awaiting-S6.

The V4 r1 design was reviewed with quorum r2 PROCEED 3/3. D-WORLD-119=A and D-WORLD-120=A are the attended rulings in `82f05004`. S5f2 #685 merged at `150d0d96` with the corrected borrowed-descriptor refusal fixtures from `c7c08e5b`, 1,476 changed product lines and Sonnet r2 PASS 77. #686 is a documentation record. Row 267 and D-WORLD-121 retain the CI race bound. S6 owns agent launch. The patches use schema 10 and leave routes, `.ail` sources and the shared launchd driver outside the cut.

Use real migrations, the repository signer, re-executed test binaries, and Git repositories created by the existing fixture under `t.TempDir()`. All service installations and supervisor calls in tests are fake, under disposable install roots. Use bounded process barriers and injected clocks for state and expiry comparisons. Compare physical paths on both sides. For guarded scope compare hashes, stat or absence before and after. The tests never seed a real service path, use real launchctl, or open real credentials.

`$W` is `/Users/voightkampff/dev/sunholo-data/.wt-world-iter441-row243-v4s5f3`, `$STATE` is `$HOME/.ailang/state/world-iter441`, and `$S` is `$STATE/plan-scratch`.

## 1. Pre-flight

### 1.1 Baseline at 2f573488

`git rev-parse HEAD` printed the full base above. `git status --short --branch` printed just `sprint/world-iter441-row243-v4s5f3` before writing the two plan files. The baseline gate started in `$S/base-tree`, a pristine `git archive 2f573488` export, before product stage edits. `$S/gate-base/legs.txt` and its raw logs record the results. Native build and vet, Linux vet, all five inventories, the GC probe, the AILANG verifier and caller census passed. The four whole CLI/daemon runs encountered bind denials. `$S/gate-base/ALL_DONE` was written after recording every leg and its comparison disposition. `verify_go` was an explicitly skipped optional leg. The package gates use plain-text Go output and `-timeout 20m`.

The controller marker is complete. The raw files and attributable flags are recorded under `$STATE/base/`. See §4.1 for any outside leg whose baseline is absent or truncated. The controller daemon-race raw log includes a time-out, so its partial fail list is truncated and excluded. The untruncated controller cmd-plain baseline has 12 failures. Its recorded set is the row-271 eleven plus TestHostsLaunchersExecHost.

### 1.2 Inputs and coupling inventory

The staged declaration lists are stored in `$S/anchors-M1.txt`, `$S/anchors-M2.txt`, `$S/anchors-M3.txt` and `$S/anchors-M4.txt`. The commands in `/bin/bash "$S/anchors.sh" "$W"` and the observed output in `$S/anchors-base.txt` pin these anchors. In this table, unqualified CLI file names are under `cmd/ailang-worldd/`.

| Input | Command and observed output |
|---|---|
| Base and history | `git diff --numstat 91d1cabf^1 91d1cabf -- cmd host` printed 17/5 `hosts_awareness.go`, 12/2 `mcp_stdio.go`, 4/1 `host/broker/handlers_exec.go` and 2/0 `host/broker/workspace_search_list.go` (35 additions, 8 deletions). `git diff --name-only 150d0d96 2f573488` listed nine documentation paths. |
| Runner dispatch | `git grep -nE 'runnerTick\|No effect dispatch\|PutOperatorHeartbeat' -- cmd/ailang-worldd/operator_runner.go` printed startup tick at :70, `runnerTick` at :80, heartbeat-only comment at :89 and write at :94. The internal mode entry is `main.go:386`. |
| Existing callers | `git grep -nE '(withOperatorStartCandidate\|reserveOperatorStartWorktree\|publishOperatorStartFiles\|MintOperatorStart)\(' -- cmd host ':!*_test.go'` printed the four declarations `operator_start_candidate.go:15`, `operator_worktree.go:30`, `operator_files.go:267` and `host/authority/mint.go:204`. The test positives are the actual calls under `operator_start_candidate_test.go` and `operator_mint_test.go`. |
| Process pin census | `git grep -nE 'TestOperator(FirstManagedEpisodeProcess\|OldBuildTerminalDrain\|CurrentOwnerProcessRecovery)\|active_child' -- cmd host` exited 1. The positive grep found `operator_upgrade_test.go:13 func TestOperatorServiceUpgrade`. |
| Terminal revoke | `git grep -nE 'runSessionRevoke\|openCredentialsForSession\(revokeCtx\|st.ResolveSession\(revokeCtx\|authority.Revoke\(revokeCtx' -- cmd/ailang-worldd/session.go` printed :361, :405, :422 and :431. `nl -ba .../session.go` places by-episode revoke at :413. The terminal lock grep prints mint :269 and provision `session_new.go:427`. `host/store/credentials.go:561` resolves the indexed ID, and `credentials_migrate.go:77,107` reads the world schema and legacy sessions. |
| Durable APIs | `git grep -nE 'func \(.* (ClaimOperatorAction\|NextOperatorStart\|ReserveOperatorStep\|ObserveOperatorStep\|CommitOperatorMint\|ReadOperatorMint\|ReadOperatorFiles\|ReadOperatorWorktree\|RecordOperatorFiles)' -- host/store` printed `operator_actions.go:74`, `operator_cas.go:119`, `operator_steps.go:17,45`, `operator_mint.go:94,242`, `operator_files.go:67,87` and `operator_worktree.go:73`. |
| Process vs episode lock | `git grep -nE 'runnerKernelLock\|terminalEpisodeLock\|held.Verify\|worktree.Add' -- cmd/ailang-worldd/{operator_runner,operator_episode_lock,operator_worktree,operator_files}.go` locates runner lock :54, terminal lock :18, worktree borrowed verification :39/:87/:92, inherited Git call :91 and file verification :275. The CLI and daemon GC positives are in `$S/gate-base/callers.census` and the final caller gate. |
| CI budget | `git grep -nF -- '-timeout 8m' scripts/verify_go.sh` printed the race command at :350. S5f2 merge measurements were 351.0 s CLI and 409.3 s daemon, and PR head 273.2 s and 360.3 s. Four recent slow-runner dev merges hit 480 s on attempt 1. |

| Inventory and base location | Disposition through M3 |
|---|---|
| `TestDaemonDependencyAllowlist`, `host/daemon/daemon_test.go:899` | Same dependency closure. There is no daemon patch or new daemon test. |
| `TestWorldBoundaryDependencyAllowlist`, `host/boundary/allowlist_world_test.go:758` | Same admitted imports. The Darwin census uses the existing `golang.org/x/sys/unix` CLI dependency. The binding helper uses the existing hash reference package. |
| `TestProductionGoSurface`, `host/store/context_roots_test.go:241` | The named inventory scans the added Go files. Its predicate and expected floor stay intact. |
| `TestProductionContextRoots`, `host/store/context_roots_test.go:136` | The revoke deadline root stays in `runSessionRevoke`. The new drains use the existing `budgetContext` constructor. Store and authority receive caller contexts. Run `go test ./host/store/ -run '^TestProductionContextRoots$'`. |
| `TestCredentialsMethodsRequireADeadline`, `host/store/credentials_test.go:277` | Recovery read and disposition call `operatorAccess(ctx)`, which requires a deadline. The extracted legacy read calls `requireDeadline` before opening the DB. |
| `TestStoreSurfaceTakesContext`, `host/store/deadline_guard_test.go:43` | New store receiver methods take context first. |
| `TestBoundaryASTWriteGuard`, `host/boundary/allowlist_world_test.go:1323` | Existing writer seam is reused. Disposition goes through the store transaction, and paired files through `serviceEnv.atomic`. |
| `TestOperatorNeverDispatch`, `host/verbs/operator_policy_test.go:69` | Preserves the existing browser refusal set and raw shell refusal. In S5f4 the new allowance is one internal confirmed-start consumer in `operator-runner` after admission, with no argv dispatch, spawn, browser ingress or finish claim. The production-caller census changes only at that stage. Keep `TestOperatorServiceFence` and `TestOperatorRunnerSetupRefusals/global_address`. |
| `TestCLICredentialWriterSites`, `cmd/ailang-worldd/cli_credential_writers_test.go:18` | Revoke adds a read-only handle, and writes still use `openCredentialsForSession`. S5f4 uses `runnerWrite`. The writer inventory remains exact. |

### 1.3 Decisions

- M1 uses the supplied episode as the managed key. Legacy opaque IDs are encoded into a disjoint key derived from the exact ID, with the original ID retained for credential lookup and revocation. Read the old schema read-only by indexed credential ID. Hold the lock through migration, the writer's exact row comparison and revoke. Active holders and durable owners both refuse.
- M2 makes the seven decisions from a single SQLite read snapshot. Evidence is scoped to a specific action, generation and owner. A transactional re-read precedes any disposition. Case 4 can retain contradictory evidence as a separate uncertain result. S5f4 adds the attended caller for the gaps. A real-signed pre-digest envelope verifies under its original bytes and is ineligible for the activation step.
- M3 moves the attended upgrade from a blanket installed-companion refusal to an exact-pair transition. The confirmation includes the reviewed companion config/plist hashes. Retire older terminal processes named like `ailang-worldd` and stop invoking previous copies before confirming. The Darwin census excludes this command and the exact supervised job PIDs. Under lifecycle and service exclusion verify pair and schema, stop admission and unload the companion, then require its kernel lock before replacement. A live inherited child still holds that lock. Restart World first, update and validate the companion config and plist, then resume the companion only if it was loaded. On a failed restart restore only the verified, compatible pair. A timeout stops the attempt and preserves the builds and receipts. The production bootout request and kernel-lock verification share the supervisor exit budget, with readback required before the apply step. An interrupted partial install can be stopped with the service record's canonical lock and repaired through the already attended uninstall/install flow.
- The proof covers processes using the supported World executable names. A renamed, manually retained old executable outside the installation is outside that census. Retirement of such copies is part of the operator's attended confirmation. A denied or incomplete native process census fails before replacement.

### 1.4 Premises found false or incomplete

| Premise | Current-base evidence and consequence |
|---|---|
| A paired upgrader was available | `git grep -nE 'operator runner is still installed\|cannot verify that the operator runner is unloaded' -- cmd/ailang-worldd/service_upgrade.go` printed :208 and :212. `operator_service.go:418` delegates stop to the supervisor. `service_launchctl.go` documents the 30 s exit allowance. Supervisor success and a stale heartbeat are insufficient proof of a retained child's exit. |
| Existing commands cover every durable window | The store reads at `operator_worktree.go:89`, `operator_mint.go:265,277` and `operator_files.go:67` require an appropriate reserved worktree or an observed atomic mint. `$S/inherited-named.raw` reproduces the observed-worktree-only and reserved-mint-abort rows on a copy rebased onto 2f573488. The three existing commands all refused both. |
| An executable digest implies executable permission and supported confinement | `git grep -nE 'launch_binary_sha256\|st.Mode\(\).Perm\(\)&0100\|Backend != "headless"' -- cmd/ailang-worldd/{operator_start_input,operator_service}.go` prints signed digest :72/:75, independent config backend :124 and execute check :150. S5f4 requires an attended digest, executable/UID and mode checks, the reconstructed headless Codex plan, and the canonical service pair. |
| Revoke already shares terminal mint exclusion | The terminal grep in §1.2 locates separate writable revoke entries, with the writer before both branches at `session.go:405`. M1 adds the same canonical lock and a read-only-before-writer path. |
| Explicit `episode gc --credentials-db` alone proves workspace binding | `git grep -nE 'EvalSymlinks\|SameDatabase\|explicit !=' -- cmd/ailang-worldd/episode_gc.go` printed :66–74. Binding is now required even with an explicit credential argument. N2 is a documented skip for an unbound non-default workspace. N3's persistent per-episode coordination sidecar is visible in `host/store/operator_episode_lock.go`. Both are carried design evidence. The S5f3 product makes no GC policy change. |
| Sandboxed native process census can prove retirement | `cd "$S/st-M3" && go run "$S/kinfo-control.go"` printed `uid scan count=0 self=false error=operation not permitted`. This result is **UNINFORMATIVE UNDER SANDBOX**. The process pins use a real child and fake supervisor with a fixture process inventory backed by `kill(pid, 0)`. Section 4.1 gives the native outside control. |

### 1.5 Size and CI race cost

| Stage | Increment | Cumulative +/− | Filtered | Unfiltered product | STATUS / DELETED / hashes |
|---|---:|---:|---:|---:|---:|
| M1 | 289 | +271/−18 | 289 | 289 | 4/0/4 |
| M2 | 397 | +668/−18 | 686 | 686 | 8/0/8 |
| M3 | 699 | +1346/−39 | 1385 | 1385 | 18/0/18 |
| S5f4 prototype M4 | 917 | +2256/−46 | 2302 | 2302 | scratch |

Measure with `git diff -M --numstat 2f573488043565f480574ed42d1dacd514c6aa0b HEAD -- cmd host docs app` after controller commits. The scratch equivalent uses `python3 "$S/stage-diff.py" M1 M2 M3 M4`, treating all trees as archives and running per-file `git diff --no-index -M --numstat`. No product rename or deletion is present. The unfiltered product total was independently confirmed with full archive-to-archive `git diff --no-index -M --numstat` in `$S/unfiltered-sealed-M3.raw` and `$S/unfiltered-sealed-M4.raw`. The two planning records are excluded from product snapshots.

| Stage and package | Named added test or arm | Race seconds |
|---|---|---:|
| M3 CLI | `TestOperatorTerminalRevokeExclusion` | 1.24 |
| M3 CLI | `TestOperatorTerminalRevokeLegacy` | 0.66 |
| M3 CLI | `TestOperatorRunnerPathUID` | 0.00 |
| M3 CLI | `TestOperatorOldBuildTerminalDrain` | 0.68 |
| M3 CLI | `TestOperatorServiceUpgrade/active_child` | 0.73 |
| M3 CLI | `TestOperatorPairedUpgradeRefusals` | 0.57 |
| M3 CLI | `TestOperatorProcessWitness` | 0.00 |
| S5f4 CLI | `TestOperatorFirstManagedEpisodeProcess` | 5.57 |
| S5f4 CLI | `TestOperatorCurrentOwnerProcessRecovery` | 3.18 |
| S5f4 CLI | `TestOperatorAttendedStartDisposition` | 2.21 |
| M2 store | `TestOperatorStartRecoveryWindows` | 1.61 |
| M2 store | `TestOperatorStartRecoveryRejectsContradictoryEvidence` | 1.39 |
| M2 authority | `TestOperatorPreDigestSignedStart` | 0.94 |

The focused receipts use `go test -race -v -count=1` with bounded test contexts. Top-level test bodies are counted once, and the existing upgrade test contributes only the new active-child arm. **S5f3 added CLI test time +3.88 s**, daemon +0.00 s. **S5f4 added CLI test time +10.96 s** against M3, daemon +0.00 s. Read the raw named receipts in `$S/timing-release/`. The older concurrent S5f4 trials measured 13.64 s (`timing-sealed`) and 5.55 s (`timing-final`). Charge only S5f4's new bodies against the activation increment. S5f4 also adds `TestOperatorPartialFilesOrphanRevoke` at 0.36 s in store. Whole host/store base/head package seconds were 151.124/119.183 (-31.941 s). Whole host/authority base/head package seconds were 32.409/71.132 (+38.723 s). Whole CLI and daemon race deltas require the outside comparison in §4.1.

### 1.6 Portability

Native and Linux-target vet are named gates. The real paired transition is macOS service-only. Linux tests supply the fake supervisor and process inventory; the production Linux fallback fails closed. The child harness uses a finite exec context, pipe readiness and exit acknowledgment, and disables only the child's race runtime exit delay. Git fixture repositories and all install roots are disposable. The test command resolves Git through the ordinary installed tool path with no gate-level PATH shim.

## 2. How the executor applies a sub-milestone

Run with Bash in the assigned sprint worktree. Apply M1, M2, M3 once each. The controller may already have committed the preceding snapshot.

````bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter441-row243-v4s5f3
STATE=$HOME/.ailang/state/world-iter441
S=$STATE/plan-scratch
T=$STATE/executor
B=2f573488043565f480574ed42d1dacd514c6aa0b
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5f3-sprint-plan.md
M=M1
mkdir -p "$T"
unset SNAP_PARENT SNAP_STATUS_LIST
git -C "$W" merge-base --is-ancestor "$B" HEAD
if [ "$M" = M1 ]; then
  git -C "$W" diff --quiet "$B" -- cmd host docs app
  test -z "$(git -C "$W" ls-files --others --exclude-standard -- cmd host docs app)"
fi
cp "$S"/{bounded.py,failset.py,verify-scratch.py,callers.sh,probe.sh,probe-fixture.py,operator_gc_after_test.go,changed-M1,changed-M2,changed-M3} "$T/"
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

Set `M=M2` and then `M=M3` on subsequent invocations. The executor does not apply the S5f4 scratch patch.

## 3. Milestones

### M1 — Terminal revoke exclusion and UID proof

`TestOperatorTerminalRevokeExclusion/holder` uses a real conflicting lock in both routes. `/durable_owner` closes the holder and preserves the durable signed owner. `/stale_lookup` changes the reported episode after the real read-only lookup, verifies the writer is still free in the handoff, and requires the stale-refusal class, unchanged post-fixture state and closure of the actual held descriptor. The same tests revoke a different credential while preserving the blocked credential. `TestOperatorTerminalRevokeLegacy/{id,episode,opaque}` asserts pre-split world bytes equal, split store absent under contention, and successful exact migration/revoke once released. `TestOperatorRunnerPathUID` runs actual file stat against both the process UID and an alternate expected login UID. The corrected judge-429 edit is unavailable. The carried ledger labels the reproducible runnerPath reconstruction explicitly.

### M2 — Seven durable windows and pre-digest compatibility

`TestOperatorStartRecoveryWindows` exercises the seven dispositions with real migration and signed claims. It asserts exact state for every pure read, forbids stale dispositions, and requires identical receipts, sessions and start record for the permitted state transitions. `TestOperatorStartRecoveryRejectsContradictoryEvidence` covers duplicate or foreign steps, invalid version, wrong kind, missing owner, non-created worktree receipts and a reserved mint with start/session evidence. `TestOperatorPreDigestSignedStart` checks real signed old and modern artifacts. The old artifact still verifies, then returns stale binding with the original bytes and zero new claim, start, session or mint receipt.

### M3 — Attended matched upgrade and companion stop

`TestOperatorServiceUpgrade/active_child` uses a real re-exec that inherits the exact open runner lock, a fake supervisor and fake install roots under `t.TempDir()`. It refuses without typed `y`, then proves reported unload alone cannot authorize replacement while the migrated store still holds its real signed claim. It exercises matched rollback after the candidate companion restart fails. The fake supervisor attempts real independent lock opens at daemon replacement and companion restart, proving the guard covers replacement and is released before bootstrap. `TestOperatorOldBuildTerminalDrain` uses a real re-exec with the older writer path, a fake supervisor and disposable install roots. It refuses with the old PID live, compares the binary, both configs/plists, DBs, observed workspace and private log, then upgrades the pair after acknowledged child exit. `TestOperatorPairedUpgradeRefusals` checks change at confirmation, incompatible schema, unavailable/missing process evidence, failed stop, and a stop report with a still-loaded companion. Keep the malformed-install, loaded-without-install, busy lifecycle, rollback-schema, build, and typed refusals in existing tests.

The operator can keep a healthy companion installed across `service upgrade --binary FILE` and sees the companion hashes in the typed prompt. It resumes only a previously loaded companion after daemon and pair readback. Stop and uninstall prove supervisor and kernel quiescence first. An unverifiable installed pair, unknown process census, mismatched build/store/schema, ongoing effect or unsafe path refuses. The existing `--pause` requirement for active daemon sessions still applies. The Quickstart records the observable stop, recovery and upgrade commands.

## 4. Verification and mutation ledger

### 4.1 Gates, fail sets and outside replay

| Stage | Stage tree for the patch | Receipt |
|---|---|---|
| Base | `$S/base-tree` | `$S/gate-base/legs.txt` and raw logs. |
| M1 | `$S/st-M1` | `$S/accept-M1.raw` and build, native vet and Linux vet in `$S/stage-checks.txt`. |
| M2 | `$S/st-M2` | `$S/accept-M2.raw` with all seven PASS arms and the signed legacy control. |
| M3 | `$S/st-M3` | `$S/seal-M3-accept.raw`, `$S/gate-M3/legs.txt`, `$S/seal-checks.txt`, and M3 race receipts in §1.5. |
| S5f4 prototype | `$S/st-M4` | Compiled against cumulative M3. `$S/seal-M4-accept.raw`, `$S/seal-checks.txt` and §1.5. |

Final recorded gate result in `$S/gate-M3/legs.txt`. `GATES_DONE 790s FAILS: compare-race-host-store UNINFORMATIVE: pkg-host-daemon race-host-daemon pkg-cmd-ailang-worldd race-cmd-ailang-worldd`. The outside commands below cover socket-related verification and missing comparison legs. The gate's comparison tags initially named worktree where they should have named authority. The corrected script compares authority. The explicit base/head plain and race fail-set comparison in `$S/authority-failset-replay.txt` found zero additions. An earlier co-load seal run saw the unchanged startup bootstrap and the mint used to set up the old revoke test fail. The startup refusal also reproduced in pristine 2f573488. Both focused tests passed on sealed M3, and M3 and M4 accepted in full after retry. The raw runs are under `$S/acceptance-retry/`. The initial store race added `TestConvertedStoreMethodsM5b/cutoff/AppendOutcome` with a deadline-classification difference. The same named arm passed three times under race at both pristine and M3. An untruncated whole store-race retry at M3 gave no new failures against the controller's empty base set. Raw evidence is under `$S/cutoff-retry/` and `$S/store-race-retry/`. No test or production fix was made for that load-sensitive result.

The CLI and daemon sandboxed whole-package legs are **UNINFORMATIVE UNDER SANDBOX** on bind denial. Treat `panic: test timed out` as truncated even when an earlier failure list exists. A nonzero Go result without attributable tests is an instrument failure. Compare fail sets, never just exit codes. Row 271's CLI eleven and daemon three, and the charter's load-sensitive tests including the fixture Git timeouts, are baseline candidates. For a new difference rerun the specific test on both pristine and staged exports before assigning a regression. Preserve the original raw failures and their classification.

Outside commands, run in full from the assigned worktree. `outside-whole.sh` records all four legs before it compares. Use the completed controller `base/cmd-plain` and `base/daemon-race` only when their raw captures are attributable and untruncated. Otherwise rerun the missing or truncated captures in the local pristine export. `-timeout 20m` is already in the script.

````bash
STATE=$HOME/.ailang/state/world-iter441
S=$STATE/plan-scratch
/bin/bash "$S/outside-whole.sh" "$S/base-tree" "$S/outside-base" pristine
/bin/bash "$S/outside-whole.sh" "$S/st-M3" "$S/outside-M3" head "$S/outside-base"
/bin/bash "$S/whole-mutants.sh" "$S/whole-mutant-ids.txt" "$S/outside-base"
(cd "$S/st-M3" && AILANG_BIN="$HOME/.pinned-ailang/ailang" python3 "$S/bounded.py" 90 go run "$S/native-census-proof.go" "$S")
````

The native control re-executes an `ailang-worldd.prev` alias inside scratch. It succeeds only with `NATIVE_CENSUS_OK self=true old=true retired=true`. The paired scripts use the untruncated controller cmd-plain baseline automatically. Its raw log and set are `$STATE/base/cmd-plain.raw` and `$STATE/base/cmd-plain.fails`. The daemon race baseline is truncated and the script records a fresh pristine leg. The controller comparison and mutation scripts accept ordinary text output and print a parser firing control. Bash `grep` is used for path predicates and log classification. `$S/timeout-audit.txt` records matches of `panic: test timed out` in raw logs. Any timeout arm is repeated with lower co-load before relying on its result.

### 4.2 Diff-anchored and carried edits

Diff-ledger anchors include the exact original and replacement text, stage, source line, named test selector and actual raw log. M1 edits run against M1; M2 edits against M2; M3 edits against M3. A killed row requires a named `--- FAIL:` and no bind denial, build failure or test timeout. Whole-package mutation rows additionally require the corresponding newly failing test compared with an attributable base. The whole-package controller list is deliberately limited to the high-priority CLI edges and carried IDs.

**Diff-anchored ledger.**

| ID | Stage | Source anchor | Named killer | Named result | Whole result |
|---|---|---|---|---|---|
| `D-RV-RO-FIRST` | M1 | `cmd/ailang-worldd/session.go:417` | `TestOperatorTerminalRevokeLegacy/id` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-RV-PREVIEW-ID` | M1 | `cmd/ailang-worldd/session.go:427` | `TestOperatorTerminalRevokeExclusion/holder` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-RV-OPAQUE` | M1 | `cmd/ailang-worldd/session.go:429` | `TestOperatorTerminalRevokeLegacy/opaque` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-RV-LOCK-KEY` | M1 | `cmd/ailang-worldd/session.go:434` | `TestOperatorTerminalRevokeExclusion/holder` | KILLED | OUTSIDE QUEUED |
| `D-RV-LOCK-REFUSAL` | M1 | `cmd/ailang-worldd/session.go:434` | `TestOperatorTerminalRevokeExclusion/holder` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-RV-RELEASE` | M1 | `cmd/ailang-worldd/session.go:439` | `TestOperatorTerminalRevokeExclusion/stale_lookup` | KILLED | OUTSIDE QUEUED |
| `D-RV-WRITER-ORDER` | M1 | `cmd/ailang-worldd/session.go:434` | `TestOperatorTerminalRevokeLegacy/episode` | KILLED | OUTSIDE QUEUED |
| `D-RV-ROUTE` | M1 | `cmd/ailang-worldd/session.go:447` | `TestOperatorTerminalRevokeLegacy/episode` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-RV-REREAD` | M1 | `cmd/ailang-worldd/session.go:457` | `TestOperatorTerminalRevokeExclusion/stale_lookup` | KILLED | OUTSIDE QUEUED |
| `D-RV-ROW-IDENTITY` | M1 | `cmd/ailang-worldd/session.go:465` | `TestOperatorTerminalRevokeExclusion/stale_lookup` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-RV-LEGACY-RO` | M1 | `cmd/ailang-worldd/session.go:483` | `TestOperatorTerminalRevokeLegacy/id` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-RV-MISSING-STORE` | M1 | `cmd/ailang-worldd/session.go:485` | `TestSessionRevokeAcceptsEitherOrderAndRefusesAnUnknownID` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-RV-RESOLVE-SPLIT` | M1 | `host/store/credentials.go:572` | `TestCredentialsMintAndRevokeBesideALiveWorldWriter` | KILLED | KILLED |
| `D-RV-LEGACY-SCHEMA` | M1 | `host/store/credentials_migrate.go:140` | `TestPreSplitStoreOpenModes` | KILLED | KILLED |
| `D-REC-KIND` | M2 | `host/store/operator_start_recovery.go:60` | `TestOperatorStartRecoveryRejectsContradictoryEvidence/wrong_kind` | KILLED | KILLED |
| `D-REC-OWNER` | M2 | `host/store/operator_start_recovery.go:60` | `TestOperatorStartRecoveryRejectsContradictoryEvidence/unowned` | KILLED | KILLED |
| `D-REC-LEDGER-FILTER` | M2 | `host/store/operator_start_recovery.go:64` | `TestOperatorStartRecoveryWindows/reserved_worktree` | KILLED | KILLED |
| `D-REC-STEP-OWNER` | M2 | `host/store/operator_start_recovery.go:87` | `TestOperatorStartRecoveryRejectsContradictoryEvidence/foreign` | KILLED | KILLED |
| `D-REC-STEP-DUPLICATE` | M2 | `host/store/operator_start_recovery.go:87` | `TestOperatorStartRecoveryRejectsContradictoryEvidence/duplicate` | KILLED | KILLED |
| `D-REC-STEP-VERSION` | M2 | `host/store/operator_start_recovery.go:87` | `TestOperatorStartRecoveryRejectsContradictoryEvidence/version_zero` | KILLED | KILLED |
| `D-REC-WORKTREE-RECEIPT` | M2 | `host/store/operator_start_recovery.go:113` | `TestOperatorStartRecoveryRejectsContradictoryEvidence/bad_receipt` | KILLED | KILLED |
| `D-REC-CLAIMED` | M2 | `host/store/operator_start_recovery.go:108` | `TestOperatorStartRecoveryWindows/claimed` | KILLED | KILLED |
| `D-REC-RESERVED-WORKTREE` | M2 | `host/store/operator_start_recovery.go:112` | `TestOperatorStartRecoveryWindows/reserved_worktree` | KILLED | KILLED |
| `D-REC-OBSERVED-WORKTREE` | M2 | `host/store/operator_start_recovery.go:116` | `TestOperatorStartRecoveryWindows/observed_worktree` | KILLED | KILLED |
| `D-REC-RESERVED-MINT` | M2 | `host/store/operator_start_recovery.go:123` | `TestOperatorStartRecoveryWindows/reserved_mint` | KILLED | KILLED |
| `D-REC-OBSERVED-MINT` | M2 | `host/store/operator_start_recovery.go:131` | `TestOperatorStartRecoveryWindows/observed_mint` | KILLED | KILLED |
| `D-REC-RESERVED-FILES` | M2 | `host/store/operator_start_recovery.go:133` | `TestOperatorStartRecoveryWindows/reserved_files` | KILLED | KILLED |
| `D-REC-PREPARED` | M2 | `host/store/operator_start_recovery.go:135` | `TestOperatorStartRecoveryWindows/prepared` | KILLED | KILLED |
| `D-REC-MINT-NONEMPTY` | M2 | `host/store/operator_start_recovery.go:120` | `TestOperatorStartRecoveryRejectsContradictoryEvidence/nonempty` | KILLED | KILLED |
| `D-REC-REVIEW-CAS` | M2 | `host/store/operator_start_recovery.go:161` | `TestOperatorStartRecoveryWindows/observed_worktree` | KILLED | KILLED |
| `D-REC-STOP-SCOPE` | M2 | `host/store/operator_start_recovery.go:169` | `TestOperatorStartRecoveryWindows/prepared` | KILLED | KILLED |
| `D-REC-RELEASE-PROOF` | M2 | `host/store/operator_start_recovery.go:174` | `TestOperatorStartRecoveryWindows/claimed` | KILLED | KILLED |
| `D-REC-RELEASE-DEADLINE` | M2 | `host/store/operator_start_recovery.go:178` | `TestOperatorStartRecoveryWindows/claimed` | KILLED | KILLED |
| `D-REC-RELEASE-OWNER` | M2 | `host/store/operator_start_recovery.go:190` | `TestOperatorStartRecoveryWindows/claimed` | KILLED | KILLED |
| `D-START-DIGEST-FACT` | M2 | `host/authority/operator_start_digest.go:12` | `TestOperatorPreDigestSignedStart` | KILLED | KILLED |
| `D-START-DIGEST-MISSING` | M2 | `host/authority/operator_start_digest.go:20` | `TestOperatorPreDigestSignedStart` | KILLED | KILLED |
| `D-PAIR-REVIEW` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:28` | `TestOperatorPairedUpgradeRefusals/approved_files_changed` | KILLED | OUTSIDE QUEUED |
| `D-PAIR-SCHEMA` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:71` | `TestOperatorPairedUpgradeRefusals/schema` | KILLED | OUTSIDE QUEUED |
| `D-PAIR-SUPERVISOR-STOP` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:113` | `TestOperatorPairedUpgradeRefusals/stop_error` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-SUPERVISOR-READBACK` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:116` | `TestOperatorPairedUpgradeRefusals/stop_reported-only` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-EXIT-BOUND` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:90` | `TestOperatorServiceUpgrade/active_child` | KILLED | OUTSIDE QUEUED |
| `D-PAIR-PROOF-UNKNOWN` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:160` | `TestOperatorPairedUpgradeRefusals/process_inventory` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-PROOF-CONTROL` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:171` | `TestOperatorPairedUpgradeRefusals/process_inventory` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-OLD-TERMINAL` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:167` | `TestOperatorOldBuildTerminalDrain` | KILLED | OUTSIDE QUEUED |
| `D-PAIR-CONFIG-SWAP` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:198` | `TestOperatorOldBuildTerminalDrain` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-PLIST-SWAP` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:201` | `TestOperatorOldBuildTerminalDrain` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-RESTART` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:204` | `TestOperatorOldBuildTerminalDrain` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-TYPED-Y` | M3 | `cmd/ailang-worldd/service_upgrade.go:203` | `TestOperatorServiceUpgrade/active_child` | KILLED | OUTSIDE QUEUED |
| `D-PAIR-LIFECYCLE-SCOPE` | M3 | `cmd/ailang-worldd/service_upgrade.go:206` | `TestOperatorServiceUpgradePairFence/locked` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-ROLLBACK-OLD` | M3 | `cmd/ailang-worldd/service_upgrade.go:337` | `TestOperatorServiceUpgrade/active_child` | KILLED | OUTSIDE QUEUED |
| `D-RV-LEGACY-RO-SCHEMA-CUT` | M1 | `cmd/ailang-worldd/session.go:482` | `TestOperatorTerminalRevokeLegacy/id` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-RV-WRITER-CLOSE` | M1 | `cmd/ailang-worldd/session.go:445` | `TestOperatorTerminalRevokeExclusion/stale_lookup` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-REC-MINT-OBSERVED-PHASE` | M2 | `host/store/operator_start_recovery.go:125` | `TestOperatorStartRecoveryWindows/observed_mint` | KILLED | KILLED |
| `D-REC-MINT-TRANSACTION-READ` | M2 | `host/store/operator_start_recovery.go:126` | `TestOperatorStartRecoveryWindows/observed_mint` | KILLED | KILLED |
| `D-REC-FILES-OBSERVED-RECEIPT` | M2 | `host/store/operator_start_recovery.go:134` | `TestOperatorStartRecoveryWindows/prepared` | KILLED | KILLED |
| `D-REC-TERMINAL-RESULT` | M2 | `host/store/operator_start_recovery.go:173` | `TestOperatorStartRecoveryWindows/reserved_mint` | KILLED | KILLED |
| `D-START-DIGEST-ALGO` | M2 | `host/authority/operator_start_digest.go:14` | `TestOperatorPreDigestSignedStart` | KILLED | KILLED |
| `D-PAIR-BOTH-PAIR-FILES` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:35` | `TestOperatorServiceUpgradePairFence/installed` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-UNINSTALLED-SUPERVISOR` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:36` | `TestOperatorServiceUpgradePairFence/registered` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-BOOTOUT-ORDER` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:111` | `TestOperatorOldBuildTerminalDrain` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-RETAINED-DESCRIPTOR-ERROR` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:85` | `TestOperatorServiceUpgrade/active_child` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-CANDIDATE-BINDING` | M3 | `cmd/ailang-worldd/operator_pair_upgrade.go:58` | `TestOperatorOldBuildTerminalDrain` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-DRAIN-BEFORE-SWAP` | M3 | `cmd/ailang-worldd/service_upgrade.go:237` | `TestOperatorServiceUpgrade/active_child` | KILLED | UNMEASURED UNDER SANDBOX |
| `D-PAIR-ROLLBACK-GUARD` | M3 | `cmd/ailang-worldd/service_upgrade.go:294` | `TestOperatorServiceUpgrade/active_child` | KILLED | OUTSIDE QUEUED |
| `D-PAIR-RESTART-WITH-GUARD` | M3 | `cmd/ailang-worldd/service_upgrade.go:272` | `TestOperatorOldBuildTerminalDrain` | KILLED | OUTSIDE QUEUED |

**Carried ledger.**

| ID | Stage | Source anchor | Named killer | Named result | Whole result |
|---|---|---|---|---|---|
| `Z-ABSENT-ANY` | M1 | `cmd/ailang-worldd/operator_episode_lock.go:42` | `TestOperatorTerminalCreateLockMissingParent` | KILLED | OUTSIDE QUEUED |
| `Z-WAIT-READBUDGET` | M1 | `cmd/ailang-worldd/operator_runner.go:115` | `TestOperatorWaitReadBudget` | KILLED | OUTSIDE QUEUED |
| `J-OWN-UID-reconstruction` | M1 | `cmd/ailang-worldd/operator_service.go:71` | `TestOperatorRunnerPathUID` | KILLED | OUTSIDE QUEUED |

The exact original judge-429 corrected edit is unavailable. The runnerPath-UID reconstruction is reported separately from the two exact Z edits. S5f3 leaves the S5f2 D-LOCK source lines untouched. CLI mutation whole replays are queued only for the prioritized list in `$S/whole-mutant-ids.txt`. The raw logs and fail-set parser controls are beside the named ledger files and the whole pure ledger. Replay named rows with `python3 "$S/mutate.py" named ID`.

The new branch audit `$S/diff-branch-audit.txt` links source changes to exact mutant anchors. Rows marked `audit only` include propagated I/O failures and predicates that share a transaction guard. They have no separately measured mutant in this plan. The counts describe only the emitted runnable edit set and its measured results.

## 5. Existing tests and helpers

| Existing test or helper changed | Production before | Production after | Mutation that still kills |
|---|---|---|---|
| `TestOperatorServiceUpgrade`, `cmd/ailang-worldd/operator_upgrade_test.go:13` | `serviceEnv.upgrade → rollbackUpgrade` via six ordinary upgrade/rollback arms | The six bodies still drive the same refusal paths. The added `active_child` drives the attended pair upgrade with a real retained descriptor and fake installed pair. | Existing wrong-schema arms require the previous schema gate, and the missing/substituted/link arms require old-binary identity. Named preserved guards are `P-ROLLBACK-OLDER`, `P-ROLLBACK-NO-PREVIOUS-SCHEMA`, `P-ROLLBACK-NO-CANDIDATE-SCHEMA`, `P-ROLLBACK-SUBSTITUTION`, `P-ROLLBACK-LINK` and `P-ROLLBACK-PREVIOUS-ABSENT`. The new arm kills `D-PAIR-EXIT-BOUND`, `D-PAIR-ROLLBACK-GUARD` and `D-PAIR-ROLLBACK-OLD`. The separate receipts are in `$S/named-existing.raw`. |

The product cut changes no existing test helper. `TestOperatorStartCollisions/{held,missing_lock}` and `TestOperatorFilePublicationAdmission/lock_held` still invoke their production helpers with the corrected S5f2 refusal fixtures. The measured S5f4 prototype adds the episode/current-owner and Git witness arms to the S5f3 process test helper. Its original callers `TestOperatorOldBuildTerminalDrain` and `TestOperatorServiceUpgrade/active_child` still exercise the same re-exec, exit and lock paths. The S5f4 named acceptance includes both. The Git witness exercises the actual inherited descriptor for the new caller.

## 6. Gate script

Run `/bin/bash "$S/gates.sh" "$S/st-M3" "$S/gate-M3" M3 "$S/gate-base"` for the final product stage, or pass the assigned worktree path after the controller commits. A `SKIPPED` leg carries its reason. The last marker is `ALL_DONE`. Treat the script's nonzero status together with its classified, recorded legs.

````bash gates.sh
#!/bin/bash
# gates.sh TREE LOGDIR base|M1|M2|M3|M4 [BASE-LOGDIR] [READ-ONLY-GIT-REF-TREE]
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
mkdir -p "$2" || exit 2
L=$(cd "$2" && pwd -P)
M=$3
BASE_LOG=${4:-}
REF=${5:-}
case "$M" in base|M1|M2|M3|M4) ;; *) exit 2;; esac
ORIGINAL=$HOME
export GIT_OPTIONAL_LOCKS=0
export AILANG_BIN="$ORIGINAL/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$ORIGINAL/.pinned-ailang-tools/v0.54.15/ailang"
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
export GOCACHE="$(go env GOCACHE)" GOMODCACHE="$(go env GOMODCACHE)" HOME="$L/home"
export XDG_CONFIG_HOME="$HOME/.config" XDG_CACHE_HOME="$HOME/.cache" XDG_DATA_HOME="$HOME/.local/share"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" "$HOME"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
unset WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
cd "$T" || exit 2
: > "$L/legs.txt"
FAILS=() UNKNOWN=() START=$(date +%s)
line() { printf '%s\n' "$*" | tee -a "$L/legs.txt"; }
skip() { line "LEG $1 SKIPPED 0s ($2)"; }
bounded() { python3 "$S/bounded.py" "$@"; }
denied() { grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$1"; }
leg() {
  local name=$1 rc=0 start; shift
  start=$(date +%s)
  "$@" > "$L/$name.txt" 2>&1 || rc=$?
  local result="rc=$rc $(( $(date +%s)-start ))s"
  if denied "$L/$name.txt" || [ "$rc" -eq 77 ]; then
    result="$result UNINFORMATIVE UNDER SANDBOX"
    UNKNOWN+=("$name")
  elif [ "$rc" -ne 0 ]; then
    FAILS+=("$name")
  fi
  line "LEG $name $result"
}
record() {
  local tag=$1 rc=0 parser=0; shift
  "$@" > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  python3 "$S/failset.py" "$L/$tag.raw" "$rc" > "$L/$tag.fails" || parser=$?
  if grep -En 'panic: test timed out|BOUND EXCEEDED' "$L/$tag.raw"; then
    echo 'INSTRUMENT FAILURE: truncated'; return 2
  fi
  if denied "$L/$tag.raw"; then grep -Ei '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$L/$tag.raw"; return 77; fi
  if [ "$parser" -ne 0 ] || [ "$rc" -gt 1 ] ||
    grep -Eq '\[(build|setup) failed\]|panic: test timed out|BOUND EXCEEDED' "$L/$tag.raw"; then
    echo 'INSTRUMENT FAILURE'; return 2
  fi
  cp "$L/$tag.fails" "$L/$tag.attributable"
  cat "$L/$tag.fails"
}
inventory() {
  local pkg=$1 names=$2 name output rc=0
  output=$(mktemp "$L/inventory.XXXXXX") || return
  bounded 180 go test -v -count=1 -timeout 2m "./$pkg" -run "^($names)\$" > "$output" 2>&1 || rc=$?
  cat "$output"
  [ "$rc" -eq 0 ] || return "$rc"
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$output" || return 2
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
fmt() {
  if [ "$M" != base ] && [ ! -s "$S/changed-$M" ]; then
    echo "missing touched-file inventory $S/changed-$M"; return 2
  fi
  if [ -f "$S/changed-$M" ]; then
    while IFS= read -r path; do
      case "$path" in *.go) [ ! -f "$path" ] || gofmt -l "$path" || return;; esac
    done < "$S/changed-$M"
  fi
}
check_fmt() { local out; out=$(fmt) || return; printf '%s\n' "$out"; [ -z "$out" ]; }
{ uptime; ps -A -o pid,pcpu,comm -r | head -12; go version; command -v git; } > "$L/load.txt" 2>&1
leg gofmt check_fmt
leg build bounded 240 go build ./...
leg vet bounded 240 go vet ./host/store ./host/daemon ./cmd/ailang-worldd ./host/authority
leg linux-vet env GOOS=linux python3 "$S/bounded.py" 240 go vet ./host/store ./host/daemon ./cmd/ailang-worldd ./host/authority
leg inv-store inventory host/store 'TestProductionContextRoots|TestProductionGoSurface|TestStoreSurfaceTakesContext|TestCredentialsMethodsRequireADeadline'
leg inv-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
leg inv-cli inventory cmd/ailang-worldd 'TestCLICredentialWriterSites|TestEveryVerbHasHelp'
leg inv-verbs inventory host/verbs TestOperatorNeverDispatch
leg inv-daemon inventory host/daemon TestDaemonDependencyAllowlist
if [ "$M" = base ]; then skip accept 'new pins absent at base'
else leg accept env WORLD_PLAN_PIN_HOME="$ORIGINAL" /bin/bash "$S/accept.sh" "$T" "$M" fast; fi
leg gc-probe env HOME="$ORIGINAL" /bin/bash "$S/probe.sh" "$T"
if [ "$M" = base ] || [ "$M" = M3 ] || [ "$M" = M4 ]; then
  leg verify_ail bounded 510 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
  if [ -n "$REF" ] && [ "${WORLD_PLAN_FULL_GO:-0}" = 1 ]; then
    leg verify_go record verify_go bounded 1500 python3 "$S/verify-scratch.py" "$T" "$L" verify_go "$REF"
  else skip verify_go 'set WORLD_PLAN_FULL_GO=1 with the read-only reference'; fi
else
  skip verify_ail 'AILANG tree byte-identical; base and M3 carry the verifier'
  skip verify_go 'pristine and final broad gate'
fi
leg callers /bin/bash "$S/callers.sh" "$T" "$L/callers.census" "$M"
for pkg in host/store host/authority host/daemon cmd/ailang-worldd; do
  name=${pkg//\//-}
  if { [ "$M" = M1 ] && { [ "$pkg" = host/daemon ] || [ "$pkg" = host/authority ]; }; } ||
     { [ "$M" = M2 ] && [ "$pkg" = host/daemon ]; }; then
    skip "pkg-$name" 'package files unchanged in cumulative stage'
    skip "race-$name" 'package files unchanged in cumulative stage'
    continue
  fi
  leg "pkg-$name" record "pkg-$name" bounded 1230 go test -count=1 -timeout 20m "./$pkg"
  leg "race-$name" record "race-$name" bounded 1230 go test -race -count=1 -timeout 20m "./$pkg"
done
# Every recording leg above runs before any attribution/comparison.
controller=${WORLD_PLAN_BASE_FAILSETS:-"$S/../base"}
recorded_set() {
  local tag=$1 source="$BASE_LOG/$1.attributable" outside=
  case "$tag" in
    pkg-host-store) outside=store-plain;;
    race-host-store) outside=store-race;;
    race-host-daemon) outside=daemon-race;;
    pkg-cmd-ailang-worldd) outside=cmd-plain;;
    race-cmd-ailang-worldd) outside=cmd-race;;
    verify_go) outside=verify_go;;
  esac
  if [ -n "$outside" ] && grep -Fq "ALL_DONE for 2f573488" "$controller/../baseline.done" 2>/dev/null && [ -f "$controller/$outside.fails" ] &&
      [ -s "$controller/$outside.raw" ] && [ -f "$controller/$outside.rc" ] &&
      ! grep -Eq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$controller/$outside.raw" &&
      ! denied "$controller/$outside.raw"; then
    python3 "$S/failset.py" "$controller/$outside.raw" "$(cat "$controller/$outside.rc")" > "$L/$tag.base" || return 2
  elif [ -f "$source" ]; then
    cp "$source" "$L/$tag.base"
  else
    return 77
  fi
}
new_failures() {
  local tag=$1
  LC_ALL=C sort -u "$L/$tag.base" -o "$L/$tag.base"
  LC_ALL=C sort -u "$L/$tag.attributable" -o "$L/$tag.attributable"
  echo "No-longer-failing tests (baseline only):"
  LC_ALL=C comm -23 "$L/$tag.base" "$L/$tag.attributable"
  echo "New failures:"
  LC_ALL=C comm -13 "$L/$tag.base" "$L/$tag.attributable" > "$L/$tag.new"
  cat "$L/$tag.new"
  [ ! -s "$L/$tag.new" ]
}
for tag in pkg-host-store race-host-store pkg-host-authority race-host-authority pkg-host-daemon race-host-daemon pkg-cmd-ailang-worldd race-cmd-ailang-worldd verify_go; do
  if [ "$M" = base ]; then
    skip "compare-$tag" 'recording the pristine base'
  elif [ -f "$L/$tag.attributable" ]; then
    if recorded_set "$tag"; then
      leg "compare-$tag" new_failures "$tag"
    else
      skip "compare-$tag" 'attributable base recording unavailable'
    fi
  else
    skip "compare-$tag" 'attributable stage recording unavailable'
  fi
done
line "GATES_DONE $(( $(date +%s)-START ))s FAILS: ${FAILS[*]:-none} UNINFORMATIVE: ${UNKNOWN[*]:-none}"
printf 'ALL_DONE\n' > "$L/ALL_DONE"
[ "${#FAILS[@]}" -eq 0 ] || exit 1
[ "${#UNKNOWN[@]}" -eq 0 ] || exit 77
````

## 7. Acceptance, snapshot and executor handoff

**Success signal** is exit 0 plus the printed `ACCEPT_OK <milestone> <mode> <artifact-directory>` line. The script requires real `--- PASS:` lines for all named tests and specific new subtest arms. Commit-time commands are focused and use fake transport/supervisors. Do not substitute a whole CLI or daemon run for these commands.

````bash accept.sh
#!/bin/bash
# accept.sh TREE M1|M2|M3|M4 fast|race
set -euo pipefail
T=$(cd "$1" && pwd -P)
M=$2
mode=$3
PIN=${WORLD_PLAN_PIN_HOME:-$HOME}
S="$PIN/.ailang/state/world-iter441/plan-scratch"
L=$(mktemp -d "$S/accept-XXXXXX")
export AILANG_BIN="$PIN/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$PIN/.pinned-ailang-tools/v0.54.15/ailang"
export GOCACHE="$(go env GOCACHE)" GOMODCACHE="$(go env GOMODCACHE)" GIT_OPTIONAL_LOCKS=0
export HOME="$L/home" AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
export XDG_CONFIG_HOME="$HOME/.config" XDG_CACHE_HOME="$HOME/.cache" XDG_DATA_HOME="$HOME/.local/share"
mkdir -p "$HOME" "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
unset WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
args=()
case "$mode" in fast) ;; race) args+=(-race);; *) exit 2;; esac
cd "$T"
test -x "$AILANG_BIN" && test -x "$WORLD_EXEC_NODE" && test -x "$WORLD_TOOL_AILANG_BIN"
run() {
  local pkg=$1 names=$2 name log="$L/${1//\//-}.raw" rc=0
  python3 "$S/bounded.py" 220 go test ${args[@]+"${args[@]}"} -v -count=1 -timeout 3m "./$pkg" -run "^($names)\$" > "$log" 2>&1 || rc=$?
  cat "$log"
  if grep -Ei 'panic: test timed out|BOUND EXCEEDED' "$log"; then echo "INSTRUMENT FAILURE: truncated $log"; return 2; fi
  if grep -Ei '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$log"; then echo "UNINFORMATIVE UNDER SANDBOX: $log"; return 77; fi
  [ "$rc" -eq 0 ] || return "$rc"
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log" || { echo "MISSING PASS $name"; return 1; }
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
revoke='TestOperatorTerminalRevokeExclusion|TestOperatorTerminalRevokeLegacy|TestOperatorRunnerPathUID|TestOperatorTerminalCreateLockMissingParent|TestOperatorWaitReadBudget|TestSessionRevokeAcceptsEitherOrderAndRefusesAnUnknownID|TestSessionRevokeByEpisode|TestSessionVerbsNeverCreateAStore|TestCLICredentialWriterSites'
matrix='TestOperatorStartRecoveryWindows|TestOperatorStartRecoveryRejectsContradictoryEvidence|TestProductionContextRoots|TestCredentialsMethodsRequireADeadline|TestStoreSurfaceTakesContext'
pair='TestOperatorOldBuildTerminalDrain|TestOperatorServiceUpgrade|TestOperatorPairedUpgradeRefusals|TestOperatorServiceUpgradePairFence|TestOperatorServiceStartup|TestOperatorServiceWriteScope|TestHomeServiceUpgradeSuccess|TestHomeServiceUpgradeActive'
case "$M" in
  M1) run cmd/ailang-worldd "$revoke" ;;
  M2)
    run host/store "$matrix"
    run host/authority TestOperatorPreDigestSignedStart
    for arm in claimed reserved_worktree observed_worktree reserved_mint observed_mint reserved_files prepared; do
      grep -Fq -- "--- PASS: TestOperatorStartRecoveryWindows/$arm (" "$L/host-store.raw"
    done
    ;;
  M3|M4)
    if [ "$M" = M4 ]; then
      pair="$pair|TestOperatorFirstManagedEpisodeProcess|TestOperatorCurrentOwnerProcessRecovery|TestOperatorAttendedStartDisposition"
      run host/store "$matrix|TestOperatorPartialFilesOrphanRevoke"
    fi
    run cmd/ailang-worldd "$revoke|$pair"
    grep -Fq -- '--- PASS: TestOperatorServiceUpgrade/active_child (' "$L/cmd-ailang-worldd.raw"
    ;;
  *) exit 2 ;;
esac
printf 'ACCEPT_OK %s %s %s\n' "$M" "$mode" "$L"
````

For each applied milestone `M`, run the extracted script with `/bin/bash "$T/snapshot.sh" "$W" "$M"`. A successful script prints the verified hash lines and counts. It writes cumulative `$STATE/snap/$M`. Prefix overrides below are rehearsal-only. The worktree script compares with the full base SHA and collects tracked plus untracked changes. It excludes exactly the two plan records. The filter does not require a physical `docs/`, `app/` or `cmd/` directory.

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1|M2|M3|M4) ;; *) exit 2;; esac
B=2f573488043565f480574ed42d1dacd514c6aa0b
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter441/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s5f3-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s5f3.json"
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

The controller's commit commands use the same focused selectors. The prepared message files are in `$S/commit-messages/`. Invoke the existing controller script from its assigned environment.

````bash
STATE=$HOME/.ailang/state/world-iter441
S=$STATE/plan-scratch
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter441-row243-v4s5f3
M=M1
msgfile=$S/commit-messages/$M
"$STATE/commit_ms.sh" "$M" "$msgfile" "AILANG_BIN=\$HOME/.pinned-ailang/ailang WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN=\$HOME/.pinned-ailang-tools/v0.54.15/ailang /bin/bash '$S/accept.sh' '$W' '$M' fast"
````

Set `M=M2` and `M=M3` for continuation. Expected cumulative NUMSTAT is in §1.5. The worktree and controller scratch paths remain independent. The controller may run `"$STATE/ctl_roundtrip.sh" M1 M2 M3` after extracting the plan. The planning rehearsal runs the same awk extraction and patch commands in `$S/rt-*`.

The extracted plan replayed under `$S/rt-*`. Every full archive comparison and snapshot manifest verified. `$S/roundtrip-final.log` contains the transcript. `REPLAY M1 diff -r rc=0; snapshot SHA-256 verified`. `REPLAY M2 diff -r rc=0; snapshot SHA-256 verified`. `REPLAY M3 diff -r rc=0; snapshot SHA-256 verified`. `REPLAY S5f4 prototype incremental diff -r rc=0`.

## 8. Scope ledger

| Scope item | Placement |
|---|---|
| Terminal revoke in both orders, stale lookup, legacy stores and J-OWN-UID | M1. Named tests in §3 and carried ledger in §4. |
| Seven durable windows and real-signed pre-digest compatibility | M2. Named store and authority tests in §3. |
| Paired attended transition, drain and compatible rollback | M3. Active-child and old-terminal process pins in §3. |
| First preparation caller, prepared-state restart, process recovery and runnable start-gap remedy | S5f4, 917 changed lines measured on cumulative M3. Appendix C gives the path, tests and boundary. |
| Agent start/spawn, launch receipt and launch recovery | S6 per reviewed design §8. Reviewed 890 raw / 1,202 calibrated. This planning work did not prototype S6. |
| Browser route, app and `.ail` | Outside S5f3 and S5f4. Reviewed route delivery starts at S7. |
| Judge-440 N2/N3 | Recorded in §1.4 as design evidence. They add no S5f3 product lines. |

## Appendix C — final sprint to the first caller

| Sprint | Starting stage | Measured incremental | Cumulative against 2f573488 | Endpoint |
|---|---|---:|---:|---|
| S5f3 (M1, M2, M3) | `$S/base-tree` | 1385 | 1385 | Prerequisites with four definitions and zero production callers. |
| **S5f4, activation** | **`$S/st-M3`** | **917** | **2302** | First caller, paired lifecycle exercised, prepared-awaiting-S6. |

The compiled S5f4 prototype is `$S/st-M4`. `$S/M4.patch` is relative to **this plan's last stage, M3**. Its changed-line budget is 1,400 with no exception needed. Rebase its patch onto the exact reviewed S5f3 product head before the S5f4 plan is applied. A changed S5f3 merge requires a fresh measurement. Keep all four actual process pins as activation gates. The inherited five-file core measured +410/−0 over this base. Its two observed gaps are reproduced in `$S/inherited-named.raw`, and its rebased size is **+576/−16 = 592**, including the added legacy read seam. It is separate evidence from the final prototype.

`runOperatorRunner` performs at most one confirmed-start visit per bounded tick. It holds the kernel runner lock and obtains the canonical episode lock before the candidate re-read through prepared readback. Acquire lifecycle exclusion across the stop check and claim. Verify service/configuration/build/store, the physical workspace, fresh heartbeat and catalogue generation, the signed action, current approver, the configured executable's UID, mode and digest, the supported headless Codex plan and private session directory. Recheck before each new effect. Lend the episode and runner descriptions to Git so an old child can still bar replacement. Let authorized readback complete after admission closes. On drift keep the receipt and direct to the matching attended remedy. One exact worktree reservation and Git readback precede mint. Mint uses the atomic start/session/step receipt. Publication rechecks current-session identity and expiry between file boundaries. Final state is the exact three matching receipts and owner. No spawn reservation or agent launch occurs in S5f4.

| Durable window | Classification and action |
|---|---|
| Claimed, no reservation | Under process and episode exclusion, re-read exact owner/version and prove no receipts, no start/session and no physical worktree/branch/files. Release to confirmed by exact CAS only with live reconstructed proof. A stale proof goes to attended disposition and a fresh review. |
| Reserved worktree | Preserve the reservation. Existing `service runner reconcile-worktree ACTION-ID` reads Git registration, branch, physical path and pinned commit after quiescence. |
| Observed worktree, no mint | Verify and retain the exact tree. `service runner reconcile-start ACTION-ID` offers typed terminal stop under both locks and exact re-read. No Git retry. |
| Reserved mint, no commit | An absent atomic observation, start and session proves the rollback cut after owner retirement. Attended stop preserves the worktree and reservation. A contradictory commit is uncertain, and the same command can record uncertainty while retaining the owner and receipts. |
| Observed mint, no files | Use `service runner reconcile-mint ACTION-ID` to examine the exact atomic record and current credential. With proved missing token, revoke only the orphan. Preserve the start and worktree. |
| Reserved or partial files | Use existing file reconciliation for exact/readable or absent file sets. Collision leaves uncertainty. The new attended remedy additionally handles an attributable directory with a verified missing token, atomically revokes only the receipt's orphan and preserves the file and worktree evidence. |
| All prep observed | Derive `prepared-awaiting-S6` from the current owned action, worktree, start/session and published file receipts. Re-read real path, branch, mode, identity and hashes under exclusion. Recovery records an observation and issues no second effect. |

| Activation pin | Witness for S5f4 |
|---|---|
| `TestOperatorFirstManagedEpisodeProcess` | The production consumer pauses in candidate acquisition. An independent terminal re-exec holds the episode lock and leaves the signed row confirmed with zero effect reservations. The reverse run uses real Git and its inherited lock witness, creates a single prepared result, refuses terminal revoke and GC, and observes no second preparation or agent launch. Fake daemon/install fixture under `t.TempDir()`. |
| `TestOperatorCurrentOwnerProcessRecovery` | Independent live holder defeats both current and stale heartbeat states. After exit, prepared readback is read-only and admits exact evidence. Changed file identity, session metadata, expiry and revocation each refuse with observed state equality. Fake daemon/install fixture under `t.TempDir()`. |
| `TestOperatorOldBuildTerminalDrain` | M3's old writer re-exec and fake supervisor/install fixture under `t.TempDir()`. |
| `TestOperatorServiceUpgrade/active_child` | M3's inherited real child and fake supervisor/install fixture under `t.TempDir()`. Verify typed approval, admission closure, unproved child refusal and schema-compatible pair rollback. |
| `TestOperatorAttendedStartDisposition`, `TestOperatorPartialFilesOrphanRevoke` | Real rollback receipt and attended CLI gap remedy, plus exact transactional orphan revocation while preserving the unrelated credential. Fake install and migrated stores under `t.TempDir()`. |

The fresh production census must find exactly the four effect/helper call sites in `operator_preparation.go` and its single internal-mode consumer entry at `operator_runner.go:86`. The measured effect calls are `operator_preparation.go:37,125,143,148`. The M4 positive raw receipt is `$S/callers-M4.raw`. Base and M1–M3 each contain only the four definitions. The existing browser `TestOperatorNeverDispatch` refusals for approver, session mint/revoke, service, serve, setup and world-publish remain exact, along with the unrelated-verb and shell controls. The new allowance is precisely the verified start consumer in the internal installed mode. Re-run the lock-borrow and publication-window selectors from S5f2, GC probe, focused race measurements and outside fail-set comparisons before enabling. The S5f4 prototype gate is `/bin/bash "$S/accept.sh" "$S/st-M4" M4 fast`. Use `/bin/bash "$S/callers.sh" "$S/st-M4" "$S/callers-M4.census" M4` for its positive production census.

## Appendix P — compiled patches

````diff M1
diff --git a/cmd/ailang-worldd/operator_revoke_test.go b/cmd/ailang-worldd/operator_revoke_test.go
new file mode 100644
index 0000000..7aa7663
--- /dev/null
+++ b/cmd/ailang-worldd/operator_revoke_test.go
@@ -0,0 +1,173 @@
+package main
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"errors"
+	"fmt"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorTerminalRevokeExclusion(t *testing.T) {
+	fixture := newEpisodeRun(t)
+	for _, scenario := range []string{"holder", "durable_owner", "stale_lookup"} {
+		t.Run(scenario, func(t *testing.T) {
+			r := fixture
+			r.f.db = dbPathOf(t)
+			cred, err := credentialsFor(r.f.db, "")
+			if err != nil {
+				t.Fatal(err)
+			}
+			mint := func(ep string) string {
+				t.Helper()
+				var out, stderr bytes.Buffer
+				if code := runSessionMint([]string{"--db", r.f.db, "--episode", ep, "--grant", "fs.read=repo:1"}, &out, &stderr,
+					testMintEnv("y\n", 1000)); code != exitOK {
+					t.Fatalf("mint fixture: %d %s", code, &stderr)
+				}
+				return hex64TokenRegex.FindString(stderr.String())
+			}
+			id, other := mint("ep"), mint("other")
+			var held *os.File
+			refusal := "episode lock held"
+			if scenario == "durable_owner" {
+				claimFixtureEpisode(t, cred, "ep").Close()
+				refusal = "active action " + fixtureActionID
+			} else if scenario == "holder" {
+				held, err = terminalEpisodeLock(cred, "ep", nil)
+				if err != nil {
+					t.Fatal(err)
+				}
+				defer held.Close()
+			}
+
+			snapshot := func() any {
+				return struct {
+					db                 string
+					repo, ws, sessions map[string]string
+				}{serviceFileHash(cred), provisionFiles(t, filepath.Join(r.f.repo, ".git")),
+					provisionFiles(t, r.f.ws), provisionFiles(t, r.sessions)}
+			}
+			before := snapshot()
+			check := func(code int, out, stderr *bytes.Buffer) {
+				t.Helper()
+				if code != exitUsage || !strings.Contains(stderr.String(), refusal) || out.Len() != 0 ||
+					!reflect.DeepEqual(before, snapshot()) {
+					t.Fatalf("revoke did not refuse without effects: %d %s", code, stderr)
+				}
+			}
+			switch scenario {
+			case "stale_lookup":
+				refusal = "session changed during revoke"
+				acquire := func(path, ep string, inherited *os.File) (*os.File, error) {
+					t.Helper()
+					if ep != "ep" || path != cred {
+						t.Fatalf("wrong lock target: %q %q", path, ep)
+					}
+					rw, err := store.OpenCredentials(cred)
+					if err != nil {
+						t.Fatalf("credentials writer acquired before exclusion: %v", err)
+					}
+					rw.Close()
+					credSQL(t, cred, "UPDATE session_credentials SET episode_id='moved' WHERE episode_id='ep'")
+					before = snapshot()
+					held, err = terminalEpisodeLock(path, ep, inherited)
+					return held, err
+				}
+				var out, stderr bytes.Buffer
+				code := sessionRevokeWithLock(boundedTestContext(t), []string{"--db", r.f.db, id}, &out, &stderr, acquire)
+				check(code, &out, &stderr)
+				if _, err := held.Stat(); !errors.Is(err, os.ErrClosed) {
+					t.Fatalf("revoke retained the descriptor: %v", err)
+				}
+			default:
+				for _, args := range [][]string{{id}, {"--episode", "ep"}} {
+					var out, stderr bytes.Buffer
+					code := runSessionRevoke(append([]string{"--db", r.f.db}, args...), &out, &stderr)
+					check(code, &out, &stderr)
+				}
+			}
+			var out, stderr bytes.Buffer
+			if code := runSessionRevoke([]string{"--db", r.f.db, other}, &out, &stderr); code != exitOK {
+				t.Fatalf("unrelated revoke returned %d: %s", code, &stderr)
+			}
+			if rows := r.f.rows(t); len(rows) != 1 || rows[0].CredentialID != id {
+				t.Fatalf("wrong rows after unrelated revoke: %+v", rows)
+			}
+		})
+	}
+}
+
+func TestOperatorTerminalRevokeLegacy(t *testing.T) {
+	for _, route := range []string{"id", "episode", "opaque"} {
+		t.Run(route, func(t *testing.T) {
+			world := dbPathOf(t)
+			episode, lockKey := "legacy", "legacy"
+			if route == "opaque" {
+				episode = "legacy/team"
+				lockKey = fmt.Sprintf("-%x", sha256.Sum256([]byte(episode)))
+			}
+			id := preSplitRow(t, world, episode)
+			cred, err := credentialsFor(world, "")
+			if err != nil {
+				t.Fatal(err)
+			}
+			args := []string{"--db", world, id}
+			if route != "id" {
+				args = []string{"--db", world, "--episode", episode}
+			}
+			held, err := store.OpenOperatorEpisodeLock(cred, lockKey)
+			if err != nil {
+				t.Fatal(err)
+			}
+			defer held.Close()
+			before := serviceFileHash(world)
+			var out, stderr bytes.Buffer
+			code := runSessionRevoke(args, &out, &stderr)
+			if code != exitUsage || !strings.Contains(stderr.String(), "episode lock held") || out.Len() != 0 ||
+				serviceFileHash(world) != before {
+				t.Fatalf("migrated before exclusion: %d %s", code, &stderr)
+			}
+			if _, err := os.Stat(cred); !errors.Is(err, os.ErrNotExist) {
+				t.Fatalf("split store exists after refusal: %v", err)
+			}
+			held.Close()
+			if code := runSessionRevoke(args, &out, &stderr); code != exitOK || !strings.Contains(out.String(), id) {
+				t.Fatalf("legacy revoke: %d %s", code, &stderr)
+			}
+			ro, err := store.OpenCredentialsReadOnly(cred)
+			if err != nil {
+				t.Fatal(err)
+			}
+			defer ro.Close()
+			if _, ok, err := ro.ResolveSession(boundedTestContext(t), id); err != nil || ok {
+				t.Fatalf("legacy row still exists: %v %v", ok, err)
+			}
+		})
+	}
+}
+
+func TestOperatorRunnerPathUID(t *testing.T) {
+	home, err := filepath.EvalSymlinks(t.TempDir())
+	if err != nil {
+		t.Fatal(err)
+	}
+	path := filepath.Join(home, "record")
+	if err := os.WriteFile(path, []byte("fixture"), 0600); err != nil {
+		t.Fatal(err)
+	}
+	before := serviceFileHash(path)
+	if err := runnerPath(home, path, os.Getuid(), false, 0600); err != nil {
+		t.Fatal(err)
+	}
+	if err := runnerPath(home, path, os.Getuid()+1, false, 0600); err == nil ||
+		!strings.Contains(err.Error(), "path ownership mismatch") || serviceFileHash(path) != before {
+		t.Fatalf("foreign UID accepted or evidence changed: %v", err)
+	}
+}
diff --git a/cmd/ailang-worldd/session.go b/cmd/ailang-worldd/session.go
index 672620c..71779e8 100644
--- a/cmd/ailang-worldd/session.go
+++ b/cmd/ailang-worldd/session.go
@@ -4,6 +4,7 @@ import (
 	"bufio"
 	"context"
 	"crypto/rand"
+	"crypto/sha256"
 	"errors"
 	"flag"
 	"fmt"
@@ -359,6 +360,12 @@ daemon refuses the revoked token from its next request.
 // — MUT-REVOKE-UNKNOWN-OK); a missing --db is refused, never created
 // (MUT-STORE-CREATE); --episode revokes every credential of one episode.
 func runSessionRevoke(args []string, stdout, stderr io.Writer) int {
+	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
+	defer cancel()
+	return sessionRevokeWithLock(ctx, args, stdout, stderr, terminalEpisodeLock)
+}
+
+func sessionRevokeWithLock(revokeCtx context.Context, args []string, stdout, stderr io.Writer, acquire func(string, string, *os.File) (*os.File, error)) int {
 	fs := flag.NewFlagSet("ailang-worldd session revoke", flag.ContinueOnError)
 	fs.SetOutput(stderr)
 	fs.Usage = func() {}
@@ -400,16 +407,43 @@ func runSessionRevoke(args []string, stdout, stderr io.Writer) int {
 		return exitUsage
 	}
 
-	revokeCtx, cancelRevoke := context.WithTimeout(context.Background(), 3*time.Second)
-	defer cancelRevoke()
+	target := *episode
+	var reviewed store.SessionRow
+	missing := func() int {
+		fmt.Fprintf(stderr, "ailang-worldd session revoke: refusing: no session credential %s in %s "+
+			"(already revoked, or minted in another store); nothing was revoked\n", credID, credPath)
+		return exitUsage
+	}
+	if target == "" {
+		var ok bool
+		reviewed, ok, err = readSessionRevokeTarget(revokeCtx, *dbPath, credPath, credID)
+		if err != nil {
+			fmt.Fprintf(stderr, "ailang-worldd session revoke: %v\n", err)
+			return exitFatal
+		}
+		if !ok {
+			return missing()
+		}
+		target = reviewed.EpisodeID
+	}
+	if !verbs.EpisodeIDPattern.MatchString(target) {
+		// Unmanaged legacy IDs may include path separators. Encode the exact
+		// value in a disjoint lock namespace without changing the lookup ID.
+		target = fmt.Sprintf("-%x", sha256.Sum256([]byte(target)))
+	}
+	lock, err := acquire(credPath, target, nil)
+	if err != nil {
+		fmt.Fprintf(stderr, "ailang-worldd session revoke: %v\n", err)
+		return exitUsage
+	}
+	defer lock.Close()
 	st, err := openCredentialsForSession(revokeCtx, *dbPath, credPath)
 	if err != nil {
 		fmt.Fprintf(stderr, "ailang-worldd session revoke: cannot open the credentials store %q: %v\n", credPath, err)
 		return exitFatal
 	}
-	defer func() { _ = st.Close() }()
+	defer st.Close()
 
-	ids := []string{credID}
 	if *episode != "" {
 		if err := revokeEpisodeCredentials(revokeCtx, st, *episode, stdout); err != nil {
 			fmt.Fprintf(stderr, "ailang-worldd session revoke: %v\n", err)
@@ -419,24 +453,46 @@ func runSessionRevoke(args []string, stdout, stderr io.Writer) int {
 			return exitFatal
 		}
 		return exitOK
-	} else if _, ok, err := st.ResolveSession(revokeCtx, credID); err != nil {
+	}
+	current, ok, err := st.ResolveSession(revokeCtx, credID)
+	if err != nil {
 		fmt.Fprintf(stderr, "ailang-worldd session revoke: %v\n", err)
 		return exitFatal
-	} else if !ok {
-		fmt.Fprintf(stderr, "ailang-worldd session revoke: refusing: no session credential %s in %s "+
-			"(already revoked, or minted in another store); nothing was revoked\n", credID, credPath)
+	}
+	if !ok {
+		return missing()
+	}
+	if current != reviewed {
+		fmt.Fprintf(stderr, "ailang-worldd session revoke: %v: session changed during revoke\n", store.ErrOperatorConflict)
 		return exitUsage
 	}
-	for _, id := range ids {
-		if err := authority.Revoke(revokeCtx, st, id); err != nil {
-			fmt.Fprintf(stderr, "ailang-worldd session revoke: %v\n", err)
-			return exitFatal
-		}
-		fmt.Fprintf(stdout, "revoked session credential %s\n", id)
+	if err := authority.Revoke(revokeCtx, st, credID); err != nil {
+		fmt.Fprintf(stderr, "ailang-worldd session revoke: %v\n", err)
+		return exitFatal
 	}
+	fmt.Fprintf(stdout, "revoked session credential %s\n", credID)
 	return exitOK
 }
 
+func readSessionRevokeTarget(ctx context.Context, world, credentials, id string) (store.SessionRow, bool, error) {
+	version, err := store.WorldSchemaVersion(world)
+	if err != nil {
+		return store.SessionRow{}, false, err
+	}
+	if version == store.PreSplitSchemaVersion {
+		return store.ResolveLegacySession(ctx, world, id)
+	}
+	if _, err := os.Stat(credentials); errors.Is(err, os.ErrNotExist) {
+		return store.SessionRow{}, false, nil
+	}
+	db, err := store.OpenCredentialsReadOnly(runnerDB(credentials))
+	if err != nil {
+		return store.SessionRow{}, false, err
+	}
+	defer db.Close()
+	return db.ResolveSession(ctx, id)
+}
+
 // parseGrants converts the repeatable --grant EFFECT=SCOPE:BUDGET specs into
 // []broker.Capability, preserving order.
 func parseGrants(specs []string) ([]broker.Capability, error) {
diff --git a/host/store/credentials.go b/host/store/credentials.go
index 4bee637..cebbc28 100644
--- a/host/store/credentials.go
+++ b/host/store/credentials.go
@@ -562,10 +562,14 @@ func (c *Credentials) ResolveSession(ctx context.Context, credentialID string) (
 	if err := requireDeadline(ctx); err != nil {
 		return SessionRow{}, false, err
 	}
+	return resolveSession(ctx, c.db, c.version, credentialID)
+}
+
+func resolveSession(ctx context.Context, db *sql.DB, version int, credentialID string) (SessionRow, bool, error) {
 	var row SessionRow
 	var project, digest sql.NullString
-	err := c.db.QueryRowContext(ctx,
-		sessionColumns(c.version)+` FROM session_credentials WHERE credential_id = ?;`,
+	err := db.QueryRowContext(ctx,
+		sessionColumns(version)+` FROM session_credentials WHERE credential_id = ?;`,
 		credentialID,
 	).Scan(&row.CredentialID, &row.EpisodeID, &row.GrantsJSON, &row.ExpiresAt, &row.CreatedAt, &project, &digest)
 	switch {
diff --git a/host/store/credentials_migrate.go b/host/store/credentials_migrate.go
index 647390a..8f3fa76 100644
--- a/host/store/credentials_migrate.go
+++ b/host/store/credentials_migrate.go
@@ -101,6 +101,25 @@ func WorldSchemaVersion(path string) (int, error) {
 // read-only and lock-free, so `session list` works beside a daemon of an
 // older build that has not moved them yet. It refuses any other version.
 func LegacySessions(ctx context.Context, worldPath, episodeID string, limit int) ([]SessionRow, error) {
+	db, err := openLegacySessionDB(ctx, worldPath)
+	if err != nil {
+		return nil, err
+	}
+	defer db.Close()
+	return listSessions(ctx, db, credentialsSchemaV1, episodeID, limit)
+}
+
+// ResolveLegacySession performs the same indexed lookup in the pre-split store.
+func ResolveLegacySession(ctx context.Context, worldPath, id string) (SessionRow, bool, error) {
+	db, err := openLegacySessionDB(ctx, worldPath)
+	if err != nil {
+		return SessionRow{}, false, err
+	}
+	defer db.Close()
+	return resolveSession(ctx, db, credentialsSchemaV1, id)
+}
+
+func openLegacySessionDB(ctx context.Context, worldPath string) (*sql.DB, error) {
 	if err := requireDeadline(ctx); err != nil {
 		return nil, err
 	}
@@ -112,16 +131,17 @@ func LegacySessions(ctx context.Context, worldPath, episodeID string, limit int)
 	if err != nil {
 		return nil, fmt.Errorf("store: legacy sessions %q: %w", worldPath, err)
 	}
-	defer func() { _ = db.Close() }()
 	db.SetMaxOpenConns(1)
 	var version int
 	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
+		db.Close()
 		return nil, fmt.Errorf("store: legacy sessions %q: %w", worldPath, err)
 	}
 	if version != preSplitSchemaVersion {
+		db.Close()
 		return nil, fmt.Errorf("store: legacy sessions %q: user_version %d is not the pre-split schema %d", worldPath, version, preSplitSchemaVersion)
 	}
-	return listSessions(ctx, db, credentialsSchemaV1, episodeID, limit)
+	return db, nil
 }
 
 // MigrateCredentials moves every session_credentials row of the schema-4 world
````

````diff M2
diff --git a/host/authority/operator_start_digest.go b/host/authority/operator_start_digest.go
new file mode 100644
index 0000000..ec9f41d
--- /dev/null
+++ b/host/authority/operator_start_digest.go
@@ -0,0 +1,21 @@
+package authority
+
+import (
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+)
+
+// OperatorStartLaunchDigest reads the immutable signed binding. Earlier
+// envelopes require a new review before an executable may be selected.
+func OperatorStartLaunchDigest(a approval.OperatorAction) (string, error) {
+	for _, fact := range a.Resolution {
+		if fact.Name == "launch_binary_sha256" {
+			hash, err := hashref.Parse(fact.Value)
+			if err == nil && hash.Algo() == "sha256" {
+				return fact.Value, nil
+			}
+			break
+		}
+	}
+	return "", &approval.BindingError{Field: "launch_binary_sha256"}
+}
diff --git a/host/authority/operator_start_digest_test.go b/host/authority/operator_start_digest_test.go
new file mode 100644
index 0000000..f3ebdf7
--- /dev/null
+++ b/host/authority/operator_start_digest_test.go
@@ -0,0 +1,47 @@
+package authority
+
+import (
+	"errors"
+	"reflect"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorPreDigestSignedStart(t *testing.T) {
+	for _, modern := range []bool{false, true} {
+		digest := hashref.SumSHA256([]byte("fixture executable")).String()
+		r := newOperatorStoreRigAction(t, false, func(a *approval.OperatorAction) {
+			if modern {
+				a.Resolution = append(a.Resolution, approval.OperatorResolved{Name: "launch_binary_sha256", Value: digest})
+			}
+		})
+		before := r.confirm(t)
+		ctx := boundedTestContext(t)
+		action, err := approval.DecodeOperatorAction(before.Action)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if _, err := RecheckClaimedOperatorAction(ctx, r.Credentials, action, [32]byte(before.Challenge), before.ExpiresAt, 1011, before); err != nil {
+			t.Fatalf("old resolution must still verify: %v", err)
+		}
+		got, err := OperatorStartLaunchDigest(action)
+		var stale *approval.BindingError
+		if modern && (err != nil || got != digest) {
+			t.Fatalf("modern digest = %q %v", got, err)
+		}
+		if !modern && (!errors.As(err, &stale) || stale.Field != "launch_binary_sha256" || got != "") {
+			t.Fatalf("pre-digest signature needs stale binding: %q %v", got, err)
+		}
+		after, err := r.OperatorAction(ctx, before.ID, 1011)
+		sessions, listErr := r.ListSessions(ctx, action.Episode, 1)
+		start, startErr := r.EpisodeStart(ctx, action.Episode)
+		_, stepErr := r.OperatorStep(ctx, store.OperatorStepKey{Action: r.cas, Kind: "mint", Identity: before.ID})
+		if err != nil || listErr != nil || startErr != nil || start != nil || len(sessions) != 0 ||
+			!errors.Is(stepErr, store.ErrOperatorMissing) || !reflect.DeepEqual(before, after) {
+			t.Fatalf("digest admission changed the signed row or created effects: %v %v %v %v", err, listErr, startErr, stepErr)
+		}
+	}
+}
diff --git a/host/store/operator_start_recovery.go b/host/store/operator_start_recovery.go
new file mode 100644
index 0000000..1312714
--- /dev/null
+++ b/host/store/operator_start_recovery.go
@@ -0,0 +1,197 @@
+package store
+
+import (
+	"context"
+	"database/sql"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+type OperatorStartWindow string
+
+const (
+	OperatorStartUntouched     OperatorStartWindow = "claimed without effects"
+	OperatorStartGitPending    OperatorStartWindow = "worktree reserved"
+	OperatorStartTreeOnly      OperatorStartWindow = "worktree observed"
+	OperatorStartMintPending   OperatorStartWindow = "mint reserved without commit"
+	OperatorStartMintUncertain OperatorStartWindow = "mint commit uncertain"
+	OperatorStartMintOnly      OperatorStartWindow = "mint observed"
+	OperatorStartFilesPending  OperatorStartWindow = "files reserved"
+	OperatorStartPrepared      OperatorStartWindow = "prepared-awaiting-S6"
+)
+
+type OperatorStartRecovery struct {
+	CAS            OperatorCAS
+	Owner, Episode string
+	Window         OperatorStartWindow
+	Worktree       string
+	steps          [3]startRecoveryStep
+}
+
+type startRecoveryStep struct {
+	Identity, Phase, Receipt string
+	Version                  int64
+}
+
+// ReadOperatorStartRecovery classifies a single transactional snapshot. An
+// owned window never proves the runner, worktree child or terminal is retired.
+func (c *Credentials) ReadOperatorStartRecovery(ctx context.Context, id string) (OperatorStartRecovery, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorStartRecovery{}, err
+	}
+	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
+	if err != nil {
+		return OperatorStartRecovery{}, err
+	}
+	defer tx.Rollback()
+	current, _, err := readOperatorStartRecovery(ctx, tx, id)
+	return current, err
+}
+
+func readOperatorStartRecovery(ctx context.Context, tx *sql.Tx, id string) (OperatorStartRecovery, OperatorActionRow, error) {
+	row, err := scanOperatorAction(tx.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id=?`, id))
+	if err != nil {
+		return OperatorStartRecovery{}, row, err
+	}
+	current := OperatorStartRecovery{CAS: OperatorCAS{row.ID, [32]byte(row.Digest), row.Generation}, Owner: row.Owner, Episode: row.Episode}
+	refuse := func() (OperatorStartRecovery, OperatorActionRow, error) {
+		return OperatorStartRecovery{}, row, ErrOperatorConflict
+	}
+	if row.Kind != "start" || row.Owner == "" {
+		return refuse()
+	}
+	steps, err := tx.QueryContext(ctx, `SELECT kind, identity, phase, owner, version, receipt
+		FROM operator_action_steps WHERE action_id=? AND kind<>'operation' LIMIT 4`, id)
+	if err != nil {
+		return OperatorStartRecovery{}, row, err
+	}
+	defer steps.Close()
+	for steps.Next() {
+		var step startRecoveryStep
+		var kind, owner string
+		var receipt []byte
+		if err := steps.Scan(&kind, &step.Identity, &step.Phase, &owner, &step.Version, &receipt); err != nil {
+			return OperatorStartRecovery{}, row, err
+		}
+		var i int
+		switch kind {
+		case "worktree":
+			i = 0
+		case "mint":
+			i = 1
+		case "files":
+			i = 2
+		default:
+			return refuse()
+		}
+		if owner != row.Owner || current.steps[i].Identity != "" || step.Identity == "" || step.Version < 1 {
+			return refuse()
+		}
+		step.Receipt = string(receipt)
+		current.steps[i] = step
+	}
+	if err := steps.Err(); err != nil {
+		return OperatorStartRecovery{}, row, err
+	}
+	tree, mint, files := current.steps[0], current.steps[1], current.steps[2]
+	var start, session bool
+	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM episode_starts WHERE episode=?),
+		EXISTS(SELECT 1 FROM session_credentials WHERE episode_id=?)`,
+		row.Episode, row.Episode).Scan(&start, &session)
+	if err != nil {
+		return OperatorStartRecovery{}, row, err
+	}
+
+	current.Worktree = tree.Identity
+	switch {
+	case row.State == "claimed" && tree.Phase == "" && mint.Phase == "" && files.Phase == "" && !start && !session:
+		current.Window = OperatorStartUntouched
+	case row.State != "running" && row.State != "uncertain":
+		return refuse()
+	case tree.Phase == "reserved" && mint.Phase == "" && files.Phase == "" && !start && !session:
+		current.Window = OperatorStartGitPending
+	case tree.Phase != "observed" || tree.Receipt != "created":
+		return refuse()
+	case mint.Phase == "" && files.Phase == "" && !start && !session:
+		current.Window = OperatorStartTreeOnly
+	case mint.Phase == "reserved" && mint.Identity == id && files.Phase == "":
+		// CommitOperatorMint makes the start, session and observation one
+		// transaction. No subset of those three is proof of rollback.
+		if start || session {
+			current.Window = OperatorStartMintUncertain
+		} else {
+			current.Window = OperatorStartMintPending
+		}
+	case mint.Phase == "observed" && mint.Identity == id:
+		if _, _, err := readOperatorFiles(ctx, tx, id); err != nil {
+			return refuse()
+		}
+		switch {
+		case files.Phase == "":
+			current.Window = OperatorStartMintOnly
+		case files.Phase == "reserved":
+			current.Window = OperatorStartFilesPending
+		case files.Phase == "observed" && files.Receipt == OperatorFilesPublished:
+			current.Window = OperatorStartPrepared
+		default:
+			return refuse()
+		}
+	default:
+		return refuse()
+	}
+	return current, row, nil
+}
+
+// DisposeOperatorStartGap applies either a fresh release of an untouched
+// claim or an attended stop. The caller proves process quiescence under the
+// episode lock and inspects any retained worktree before calling.
+func (c *Credentials) DisposeOperatorStartGap(ctx context.Context, reviewed OperatorStartRecovery, release *approval.OperatorMintProof, now int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	current, row, err := readOperatorStartRecovery(ctx, tx, reviewed.CAS.ID)
+	if err != nil {
+		return err
+	}
+	if current != reviewed {
+		return ErrOperatorConflict
+	}
+	a, err := decodeOperatorMintAction(row)
+	if err != nil || c.lock == nil || a.Store != c.lock.dbPath {
+		return ErrOperatorConflict
+	}
+	switch current.Window {
+	case OperatorStartUntouched, OperatorStartTreeOnly, OperatorStartMintPending:
+	default:
+		return ErrOperatorConflict
+	}
+	state, terminalAt := "failed", any(now)
+	if release != nil {
+		if current.Window != OperatorStartUntouched {
+			return ErrOperatorConflict
+		}
+		if now >= row.ExpiresAt {
+			return ErrOperatorExpired
+		}
+		binding, err := release.Bind(a, [32]byte(row.Challenge), row.ExpiresAt, row.Artifact)
+		if err != nil || row.ApprovedBy != binding.Format+":"+binding.Name {
+			return ErrOperatorConflict
+		}
+		if err := operatorMintApprover(ctx, tx, release, binding, now); err != nil {
+			return err
+		}
+		state, terminalAt = "confirmed", nil
+	}
+	result, err := tx.ExecContext(ctx, `UPDATE operator_actions SET state=?, terminal_at=?, owner=NULL, generation=generation+1
+		WHERE id=? AND digest=? AND generation=? AND owner=?`, state, terminalAt,
+		current.CAS.ID, current.CAS.Digest[:], current.CAS.Generation, current.Owner)
+	if err := operatorChangedOne(result, err); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
diff --git a/host/store/operator_start_recovery_test.go b/host/store/operator_start_recovery_test.go
new file mode 100644
index 0000000..29ba0bc
--- /dev/null
+++ b/host/store/operator_start_recovery_test.go
@@ -0,0 +1,132 @@
+package store
+
+import (
+	"errors"
+	"reflect"
+	"testing"
+)
+
+func TestOperatorStartRecoveryWindows(t *testing.T) {
+	for i, want := range []OperatorStartWindow{OperatorStartUntouched, OperatorStartGitPending,
+		OperatorStartTreeOnly, OperatorStartMintPending, OperatorStartMintOnly, OperatorStartFilesPending, OperatorStartPrepared} {
+		t.Run([]string{"claimed", "reserved_worktree", "observed_worktree", "reserved_mint", "observed_mint", "reserved_files", "prepared"}[i], func(t *testing.T) {
+			f := newOperatorPinFixture(t, false)
+			id := f.in.Attempt.CAS.ID
+			if i >= 4 {
+				f.mint(t)
+			} else if i >= 1 {
+				worktree := f.reserveWorktree(t)
+				if i >= 2 {
+					if err := f.c.ObserveOperatorStep(f.ctx, worktree, []byte("created")); err != nil {
+						t.Fatal(err)
+					}
+					if i == 3 {
+						key := OperatorStepKey{worktree.Action, worktree.Owner, "mint", id, 1}
+						if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
+							t.Fatal(err)
+						}
+					}
+				}
+			}
+			if i >= 5 {
+				read, _, err := f.c.ReadOperatorFiles(f.ctx, id)
+				if err != nil {
+					t.Fatal(err)
+				}
+				if err := f.c.RecordOperatorFiles(f.ctx, read, "files-identity", f.proof, func() int64 { return pinNow }); err != nil {
+					t.Fatal(err)
+				}
+				if i == 6 {
+					read, _, err = f.c.ReadOperatorFiles(f.ctx, id)
+					if err != nil {
+						t.Fatal(err)
+					}
+					if err := f.c.RecordOperatorFiles(f.ctx, read, OperatorFilesPublished, f.proof, func() int64 { return pinNow }); err != nil {
+						t.Fatal(err)
+					}
+				}
+			}
+
+			before := f.snapshot(t)
+			got, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
+			if err != nil || got.Window != want || got.CAS.Generation != before.action.Generation ||
+				got.Owner != before.action.Owner || !reflect.DeepEqual(before, f.snapshot(t)) {
+				t.Fatalf("read recovery = %q %v, want %q with unchanged state", got.Window, err, want)
+			}
+			stale := got
+			stale.CAS.Generation++
+			f.assertRefused(t, before, f.c.DisposeOperatorStartGap(f.ctx, stale, nil, pinNow), ErrOperatorConflict)
+			switch want {
+			case OperatorStartUntouched, OperatorStartTreeOnly, OperatorStartMintPending:
+				if want == OperatorStartUntouched {
+					f.assertRefused(t, before, f.c.DisposeOperatorStartGap(f.ctx, got, f.proof, before.action.ExpiresAt), ErrOperatorExpired)
+					if err := f.c.DisposeOperatorStartGap(f.ctx, got, f.proof, pinNow); err != nil {
+						t.Fatal(err)
+					}
+				} else if err := f.c.DisposeOperatorStartGap(f.ctx, got, nil, pinNow); err != nil {
+					t.Fatal(err)
+				}
+				after := f.snapshot(t)
+				state := "failed"
+				if want == OperatorStartUntouched {
+					state = "confirmed"
+				}
+				if after.action.State != state || after.action.Owner != "" || after.action.Generation != got.CAS.Generation+1 ||
+					before.steps != after.steps || !reflect.DeepEqual(before.start, after.start) ||
+					!reflect.DeepEqual(before.sessions, after.sessions) {
+					t.Fatalf("disposition = %s with owner %q and changed evidence", after.action.State, after.action.Owner)
+				}
+			default:
+				f.assertRefused(t, before, f.c.DisposeOperatorStartGap(f.ctx, got, nil, pinNow), ErrOperatorConflict)
+			}
+		})
+	}
+}
+
+func TestOperatorStartRecoveryRejectsContradictoryEvidence(t *testing.T) {
+	for _, kind := range []string{"unexpected", "foreign", "duplicate", "nonempty", "wrong_kind", "version_zero", "unowned", "bad_receipt"} {
+		t.Run(kind, func(t *testing.T) {
+			f := newOperatorPinFixture(t, false)
+			id := f.in.Attempt.CAS.ID
+			if kind == "nonempty" {
+				f.mint(t)
+				operatorExec(t, f.c, "UPDATE operator_action_steps SET phase='reserved', receipt=NULL WHERE kind='mint'")
+			} else if kind == "wrong_kind" {
+				operatorExec(t, f.c, "UPDATE operator_actions SET kind='finish'")
+			} else if kind == "unowned" {
+				operatorExec(t, f.c, "UPDATE operator_actions SET owner=''")
+			} else {
+				key := f.reserveWorktree(t)
+				switch kind {
+				case "version_zero":
+					operatorExec(t, f.c, "UPDATE operator_action_steps SET version=0 WHERE kind='worktree'")
+				case "bad_receipt":
+					if err := f.c.ObserveOperatorStep(f.ctx, key, []byte("other")); err != nil {
+						t.Fatal(err)
+					}
+				case "unexpected":
+					operatorExec(t, f.c, "UPDATE operator_action_steps SET kind='launch' WHERE kind='worktree'")
+				case "foreign":
+					operatorExec(t, f.c, "UPDATE operator_action_steps SET owner='foreign' WHERE kind='worktree'")
+				case "duplicate":
+					key.Identity += "-second"
+					if err := f.c.ReserveOperatorStep(f.ctx, key, pinNow); err != nil {
+						t.Fatal(err)
+					}
+				}
+			}
+			before := f.snapshot(t)
+			current, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
+			if kind == "nonempty" {
+				if err != nil || current.Window != OperatorStartMintUncertain {
+					t.Fatalf("contradictory mint commit evidence: %q %v", current.Window, err)
+				}
+				f.assertRefused(t, before, f.c.DisposeOperatorStartGap(f.ctx, current, nil, pinNow), ErrOperatorConflict)
+				return
+			}
+			if !errors.Is(err, ErrOperatorConflict) || !reflect.DeepEqual(before, f.snapshot(t)) {
+				t.Fatalf("contradictory evidence: %v", err)
+			}
+		})
+	}
+}
````

````diff M3
diff --git a/cmd/ailang-worldd/operator_pair_upgrade.go b/cmd/ailang-worldd/operator_pair_upgrade.go
new file mode 100644
index 0000000..2cce4ee
--- /dev/null
+++ b/cmd/ailang-worldd/operator_pair_upgrade.go
@@ -0,0 +1,208 @@
+package main
+
+import (
+	"encoding/json"
+	"errors"
+	"fmt"
+	"os"
+	"strconv"
+	"strings"
+	"syscall"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+type runnerUpgradePair struct {
+	config        runnerConfig
+	before, after []byte
+	loaded        bool
+}
+
+func (e serviceEnv) runnerUpgradeReview() [2]string {
+	p := e.runnerPaths()
+	return [2]string{serviceFileHash(p.Record), serviceFileHash(p.Plist)}
+}
+
+func (e serviceEnv) reviewRunnerUpgrade(approved [2]string, from, to string, schemas ...int) (*runnerUpgradePair, error) {
+	if approved != e.runnerUpgradeReview() {
+		return nil, errors.New("operator runner changed during confirmation")
+	}
+	p := e.runnerPaths()
+	_, recordErr := os.Lstat(p.Record)
+	_, plistErr := os.Lstat(p.Plist)
+	_, loadedErr := e.controller.Print(e.domain, runnerLabel)
+	if errors.Is(recordErr, os.ErrNotExist) && errors.Is(plistErr, os.ErrNotExist) {
+		if !errors.Is(loadedErr, errServiceNotLoaded) {
+			return nil, errors.New("cannot verify that the operator runner is unloaded")
+		}
+		return nil, nil
+	}
+	c, seen, err := e.installedRunner()
+	if err != nil {
+		return nil, errors.New("operator runner is still installed but its configuration and plist cannot be verified")
+	}
+	if approved[0] != serviceHash(seen) || c.Executables[e.paths.Binary] != from {
+		return nil, errors.New("operator runner changed during confirmation")
+	}
+	if loadedErr != nil && !errors.Is(loadedErr, errServiceNotLoaded) {
+		return nil, errors.New("cannot verify that the operator runner is unloaded")
+	}
+	if err := e.checkRunnerPair(c); err != nil {
+		return nil, err
+	}
+	pair := &runnerUpgradePair{config: c, before: seen, loaded: loadedErr == nil}
+	if err := pair.checkSchema(schemas...); err != nil {
+		return nil, err
+	}
+	c.Executables[e.paths.Binary] = to
+	pair.after, err = json.MarshalIndent(c, "", "  ")
+	c.Executables[e.paths.Binary] = from
+	return pair, err
+}
+
+func (p *runnerUpgradePair) checkSchema(schemas ...int) error {
+	db, err := store.OpenCredentialsReadOnly(runnerDB(p.config.Credentials))
+	if err != nil {
+		return err
+	}
+	defer db.Close()
+	for _, schema := range schemas {
+		if schema != db.Version() {
+			return errors.New("paired upgrade requires both builds to match the actual credentials schema")
+		}
+	}
+	return nil
+}
+
+func (e serviceEnv) waitRunnerExit(c runnerConfig) (*os.File, error) {
+	return e.waitRunnerExitUntil(c, e.runnerNow().Add(min(e.wait, serviceExitTimeout)))
+}
+
+func (e serviceEnv) waitRunnerExitUntil(c runnerConfig, deadline time.Time) (*os.File, error) {
+	for {
+		lock, err := runnerKernelLock(c.Credentials + ".operator-runner.lock")
+		if err == nil || !errors.Is(err, syscall.EWOULDBLOCK) {
+			return lock, err
+		}
+		remaining := deadline.Sub(e.runnerNow())
+		if remaining <= 0 {
+			return nil, errors.New("companion or retained child still holds the runner lock")
+		}
+		e.runnerSleep(min(remaining, 100*time.Millisecond))
+	}
+}
+
+func (e serviceEnv) stopRunnerJob(budget time.Duration) error {
+	_, err := e.controller.Print(e.domain, runnerLabel)
+	if errors.Is(err, errServiceNotLoaded) {
+		return nil
+	}
+	if err != nil {
+		return errors.New("cannot verify that the operator runner is unloaded")
+	}
+	if controller, ok := e.controller.(systemServiceController); ok {
+		var output string
+		output, err = controller.commandWithin(budget, "bootout", e.domain+"/"+runnerLabel)
+		if err != nil && strings.Contains(output, "Could not find service") {
+			err = errServiceNotLoaded
+		}
+	} else {
+		err = e.controller.Bootout(e.domain, runnerLabel)
+	}
+	if err != nil && !errors.Is(err, errServiceNotLoaded) {
+		return errors.New("could not stop the companion")
+	}
+	if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) {
+		return errors.New("cannot verify that the operator runner is unloaded")
+	}
+	return nil
+}
+
+func (e serviceEnv) drainRunner(c runnerConfig) (*os.File, error) {
+	budget := min(e.wait, serviceExitTimeout)
+	deadline := e.runnerNow().Add(budget)
+	if err := e.stopRunnerJob(budget); err != nil {
+		return nil, err
+	}
+	if !e.runnerNow().Before(deadline) {
+		return nil, errors.New("companion drain budget expired")
+	}
+	return e.waitRunnerExitUntil(c, deadline)
+}
+
+func (e serviceEnv) proveOldTerminalQuiet() error {
+	allowed := map[int]bool{os.Getpid(): true}
+	for _, label := range []string{e.label, runnerLabel} {
+		job, err := e.controller.Print(e.domain, label)
+		if errors.Is(err, errServiceNotLoaded) {
+			continue
+		}
+		if err != nil {
+			return errors.New("cannot inventory supervised World processes")
+		}
+		pid := 0
+		for _, line := range strings.Split(job, "\n") {
+			if raw, ok := strings.CutPrefix(strings.TrimSpace(line), "pid = "); ok {
+				pid, _ = strconv.Atoi(raw)
+				break
+			}
+		}
+		if pid <= 0 {
+			return errors.New("cannot identify a loaded World process")
+		}
+		allowed[pid] = true
+	}
+	read := e.runnerProcesses
+	if read == nil {
+		read = listWorldProcesses
+	}
+	rows, err := read()
+	if err != nil {
+		return errors.New("cannot verify old terminal process retirement")
+	}
+	self := false
+	for _, process := range rows {
+		self = self || process.pid == os.Getpid()
+		if !allowed[process.pid] && strings.HasPrefix(process.name, "ailang-worldd") {
+			return fmt.Errorf("old terminal process retirement not proved (pid %d)", process.pid)
+		}
+	}
+	if !self {
+		return errors.New("World process inventory did not include this process")
+	}
+	return nil
+}
+
+type worldProcess struct {
+	pid  int
+	name string
+}
+
+func (e serviceEnv) writeRunnerPair(pair *runnerUpgradePair, before bool) error {
+	data := pair.after
+	if before {
+		data = pair.before
+	}
+	var c runnerConfig
+	if err := json.Unmarshal(data, &c); err != nil {
+		return err
+	}
+	p := e.runnerPaths()
+	if err := e.checkRunnerConfig(c); err != nil {
+		return err
+	}
+	if err := e.checkRunnerPair(c); err != nil {
+		return err
+	}
+	if err := e.atomic(p.Record, data, 0600); err != nil {
+		return err
+	}
+	if err := e.atomic(p.Plist, []byte(runnerPlist(p, c, data)), 0600); err != nil {
+		return err
+	}
+	if pair.loaded {
+		return e.changeRunner("start", nil)
+	}
+	return nil
+}
diff --git a/cmd/ailang-worldd/operator_pair_upgrade_test.go b/cmd/ailang-worldd/operator_pair_upgrade_test.go
new file mode 100644
index 0000000..c703dcd
--- /dev/null
+++ b/cmd/ailang-worldd/operator_pair_upgrade_test.go
@@ -0,0 +1,342 @@
+package main
+
+import (
+	"bufio"
+	"bytes"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"os"
+	"os/exec"
+	"reflect"
+	"strconv"
+	"strings"
+	"syscall"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorProcessWitness(t *testing.T) {
+	mode := os.Getenv("WORLD_TEST_OPERATOR_WITNESS")
+	if mode == "" {
+		return
+	}
+	wait := func() error {
+		fmt.Println("ready")
+		_, err := io.Copy(io.Discard, os.Stdin)
+		return err
+	}
+	path := os.Getenv("WORLD_TEST_OPERATOR_WITNESS_PATH")
+	switch mode {
+	case "inherited":
+		lock := os.NewFile(3, "inherited runner lock")
+		defer lock.Close()
+		if _, err := lock.Stat(); err != nil {
+			t.Fatal(err)
+		}
+	case "old-terminal":
+		// Pre-exclusion writable terminal entry. Hold the real writer
+		// without acquiring either operator lock.
+		if err := runnerWrite(path, func(*store.Credentials) error { return wait() }); err != nil {
+			t.Fatal(err)
+		}
+		return
+	default:
+		t.Fatalf("bad witness mode %q", mode)
+	}
+	if err := wait(); err != nil {
+		t.Fatal(err)
+	}
+}
+
+func operatorProcessWitness(t *testing.T, mode, path string, files ...*os.File) (int, func()) {
+	t.Helper()
+	ctx, cancel := budgetContext(8 * time.Second)
+	bin, err := os.Executable()
+	if err != nil {
+		cancel()
+		t.Fatal(err)
+	}
+	cmd := exec.CommandContext(ctx, bin, "-test.run=^TestOperatorProcessWitness$")
+	cmd.WaitDelay = time.Second
+	cmd.ExtraFiles = files
+	cmd.Env = append(os.Environ(), "WORLD_TEST_OPERATOR_WITNESS="+mode, "WORLD_TEST_OPERATOR_WITNESS_PATH="+path,
+		"GORACE="+strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
+	stdin, err := cmd.StdinPipe()
+	if err != nil {
+		cancel()
+		t.Fatal(err)
+	}
+	stdout, err := cmd.StdoutPipe()
+	if err != nil {
+		cancel()
+		t.Fatal(err)
+	}
+	var stderr bytes.Buffer
+	cmd.Stderr = &stderr
+	if err := cmd.Start(); err != nil {
+		cancel()
+		t.Fatal(err)
+	}
+	waited := false
+	t.Cleanup(func() {
+		cancel()
+		if !waited {
+			_ = cmd.Wait()
+		}
+	})
+	reader := bufio.NewReader(stdout)
+	line, err := reader.ReadString('\n')
+	if err != nil || line != "ready\n" {
+		t.Fatalf("process barrier: %q %v (%s)", line, err, &stderr)
+	}
+	release := func() {
+		t.Helper()
+		if !waited {
+			stdin.Close()
+			io.Copy(io.Discard, reader)
+			err := cmd.Wait()
+			waited = true
+			if err != nil {
+				t.Fatalf("witness exit: %v: %s", err, &stderr)
+			}
+		}
+	}
+	return cmd.Process.Pid, release
+}
+
+// Pair fixture uses fake supervisor, HTTP transport and install roots from
+// operatorService, all under t.TempDir().
+type pairedUpgradeSupervisor struct {
+	serviceController
+	f            *operatorServiceFixture
+	runnerLoaded bool
+	bootFailure  bool
+	stopFault    string
+}
+
+func (s *pairedUpgradeSupervisor) Print(domain, label string) (string, error) {
+	if label == runnerLabel && !s.runnerLoaded {
+		return "", errServiceNotLoaded
+	}
+	return s.serviceController.Print(domain, label)
+}
+
+func (s *pairedUpgradeSupervisor) Bootout(domain, label string) error {
+	if label == runnerLabel {
+		switch s.stopFault {
+		case "error":
+			return errors.New("fixture stop failed")
+		case "reported-only":
+			return nil
+		}
+		s.runnerLoaded = false
+		return nil
+	}
+	return s.serviceController.Bootout(domain, label)
+}
+
+func (s *pairedUpgradeSupervisor) Bootstrap(domain, plist string) error {
+	if plist != s.f.env.runnerPaths().Plist {
+		return s.serviceController.Bootstrap(domain, plist)
+	}
+	c, _, err := s.f.env.installedRunner()
+	if err != nil {
+		return err
+	}
+	process, err := runnerKernelLock(c.Credentials + ".operator-runner.lock")
+	if err != nil {
+		return fmt.Errorf("companion restart could not acquire its process lock: %w", err)
+	}
+	process.Close()
+	s.runnerLoaded = true
+	if s.bootFailure {
+		s.bootFailure = false
+		return errors.New("refusing first companion boot")
+	}
+	s.f.heartbeat = store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "paired-" + serviceShort(c.Executables[s.f.env.paths.Binary])},
+		Protocol: store.OperatorRunnerProtocol, Build: c.Executables[s.f.env.paths.Binary]}
+	return s.f.env.runnerTick(c, &s.f.heartbeat)
+}
+
+func pairedUpgradeFixture(t *testing.T) (*operatorServiceFixture, *pairedUpgradeSupervisor, string) {
+	t.Helper()
+	f := operatorService(t)
+	f.mustInstall(t)
+	f.homeJob.runner = true
+	f.homeJob.onKick = func() error {
+		probe, err := runnerKernelLock(f.config.Credentials + ".operator-runner.lock")
+		if probe != nil {
+			probe.Close()
+		}
+		if !errors.Is(err, syscall.EWOULDBLOCK) {
+			t.Fatalf("daemon replacement lost its runner drain guard: %v", err)
+		}
+		f.health.Build.ExecutableSHA256 = f.homeJob.running
+		return nil
+	}
+	candidate := f.source
+	if err := os.WriteFile(candidate, []byte("candidate fixture"), 0700); err != nil {
+		t.Fatal(err)
+	}
+	s := &pairedUpgradeSupervisor{serviceController: f.homeJob, f: f, runnerLoaded: true}
+	f.env.controller = s
+	f.env.runnerProcesses = func() ([]worldProcess, error) {
+		return []worldProcess{{os.Getpid(), "ailang-worldd"}}, nil
+	}
+	return f, s, candidate
+}
+
+func pairedUpgradeEvidence(f *operatorServiceFixture) [7]string {
+	p := f.env.runnerPaths()
+	return [7]string{serviceFileHash(f.env.paths.Binary), serviceFileHash(f.env.paths.Record), serviceFileHash(f.env.paths.Plist),
+		serviceFileHash(p.Record), serviceFileHash(p.Plist), serviceFileHash(f.config.Credentials), serviceFileHash(f.config.World)}
+}
+
+func TestOperatorOldBuildTerminalDrain(t *testing.T) {
+	f, s, candidate := pairedUpgradeFixture(t)
+	pid, release := operatorProcessWitness(t, "old-terminal", f.config.Credentials)
+	read := f.env.runnerProcesses
+	f.env.runnerProcesses = func() ([]worldProcess, error) {
+		rows, _ := read()
+		if syscall.Kill(pid, 0) == nil {
+			rows = append(rows, worldProcess{pid, "ailang-worldd.prev"})
+		}
+		return rows, nil
+	}
+	before := pairedUpgradeEvidence(f)
+	record, err := readServiceRecord(f.env.paths)
+	if err != nil {
+		t.Fatal(err)
+	}
+	workspaceBefore := provisionFiles(t, record.Workspace)
+	spool := f.env.runnerPaths().Log
+	filesBefore := serviceFileHash(spool)
+	code, msg := f.run("upgrade", "--binary", candidate)
+	if code != 1 || !strings.Contains(msg, "old terminal process retirement not proved (pid "+strconv.Itoa(pid)+")") ||
+		pairedUpgradeEvidence(f) != before || !reflect.DeepEqual(workspaceBefore, provisionFiles(t, record.Workspace)) || serviceFileHash(spool) != filesBefore || !s.runnerLoaded {
+		t.Fatalf("upgrade passed a live old-build writer: %d %s", code, msg)
+	}
+	release()
+	code, msg = f.run("upgrade", "--binary", candidate)
+	if code != 0 || !s.runnerLoaded || f.health.Build.ExecutableSHA256 != serviceFileHash(candidate) {
+		t.Fatalf("retired terminal and matched pair did not upgrade: %d %s", code, msg)
+	}
+}
+
+func testOperatorUpgradeActiveChild(t *testing.T) {
+	f, s, candidate := pairedUpgradeFixture(t)
+	claimFixtureEpisode(t, f.config.Credentials, "ep").Close()
+	old := pairedUpgradeEvidence(f)
+	lock, err := runnerKernelLock(f.config.Credentials + ".operator-runner.lock")
+	if err != nil {
+		t.Fatal(err)
+	}
+	_, release := operatorProcessWitness(t, "inherited", "", lock)
+	lock.Close()
+
+	original := f.env.terminal.openTerminal
+	f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+		return &fakeTerm{r: strings.NewReader("n\n")}, nil
+	}
+	code, msg := f.run("upgrade", "--binary", candidate)
+	if code != 1 || !strings.Contains(msg, "aborted (not confirmed)") || pairedUpgradeEvidence(f) != old || !s.runnerLoaded {
+		t.Fatalf("upgrade crossed typed refusal: %d %s", code, msg)
+	}
+	f.env.terminal.openTerminal = original
+	code, msg = f.run("upgrade", "--binary", candidate)
+	if code != 1 || !strings.Contains(msg, "retained child still holds the runner lock") || pairedUpgradeEvidence(f) != old ||
+		s.runnerLoaded || serviceFileHash(candidate) == "" || serviceFileHash(f.env.paths.Binary+".prev") != "" {
+		t.Fatalf("upgrade swapped before child exit: %d %s", code, msg)
+	}
+	release()
+
+	// Restore the pre-upgrade admission state, then exercise rollback after
+	// the candidate daemon answers but its companion fails to start.
+	s.runnerLoaded, s.bootFailure = true, true
+	code, msg = f.run("upgrade", "--binary", candidate)
+	rollback := pairedUpgradeEvidence(f)
+	if code != 2 || !strings.Contains(msg, "paired upgrade rolled back") || rollback[0] != old[0] ||
+		rollback[3] != old[3] || rollback[4] != old[4] || !s.runnerLoaded ||
+		serviceFileHash(f.env.paths.Binary+".failed") != serviceFileHash(candidate) {
+		t.Fatalf("compatible pair rollback: %d %s", code, msg)
+	}
+}
+
+func TestOperatorPairedUpgradeRefusals(t *testing.T) {
+	f, supervisor, candidate := pairedUpgradeFixture(t)
+	before := pairedUpgradeEvidence(f)
+	check := func(t *testing.T, want string, evidence [7]string) {
+		t.Helper()
+		code, msg := f.run("upgrade", "--binary", candidate)
+		if code != 1 || !strings.Contains(msg, want) || pairedUpgradeEvidence(f) != evidence ||
+			!supervisor.runnerLoaded || f.homeJob.kicks != 0 || serviceFileHash(f.env.paths.Binary+".prev") != "" {
+			t.Fatalf("paired refusal %q: %d %s", want, code, msg)
+		}
+	}
+	t.Run("schema", func(t *testing.T) {
+		probe := f.env.probe
+		t.Cleanup(func() { f.env.probe = probe })
+		f.env.probe = func(path, home string) ([]byte, error) {
+			b, err := probe(path, home)
+			if err != nil || path != candidate {
+				return b, err
+			}
+			var v serviceVersion
+			if err := json.Unmarshal(b, &v); err != nil {
+				return nil, err
+			}
+			v.CredentialsSchema++
+			return json.Marshal(v)
+		}
+		check(t, "both builds to match the actual credentials schema", before)
+	})
+	t.Run("process_inventory", func(t *testing.T) {
+		read := f.env.runnerProcesses
+		t.Cleanup(func() { f.env.runnerProcesses = read })
+		f.env.runnerProcesses = func() ([]worldProcess, error) { return nil, errors.New("fixture inventory unavailable") }
+		check(t, "cannot verify old terminal process retirement", before)
+		f.env.runnerProcesses = func() ([]worldProcess, error) { return nil, nil }
+		check(t, "did not include this process", before)
+	})
+	for _, tc := range []struct{ fault, want string }{
+		{"error", "could not stop the companion"},
+		{"reported-only", "cannot verify that the operator runner is unloaded"},
+	} {
+		t.Run("stop_"+tc.fault, func(t *testing.T) {
+			supervisor.stopFault = tc.fault
+			t.Cleanup(func() { supervisor.stopFault = "" })
+			check(t, tc.want, before)
+		})
+	}
+	t.Run("approved_files_changed", func(t *testing.T) {
+		terminal := f.env.terminal.openTerminal
+		p := f.env.runnerPaths()
+		data, err := os.ReadFile(p.Record)
+		if err != nil {
+			t.Fatal(err)
+		}
+		t.Cleanup(func() {
+			f.env.terminal.openTerminal = terminal
+			if err := os.WriteFile(p.Record, data, 0600); err != nil {
+				t.Error(err)
+			}
+		})
+		var changed [7]string
+		f.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
+			if err := os.WriteFile(p.Record, append(data, '\n'), 0600); err != nil {
+				t.Fatal(err)
+			}
+			changed = pairedUpgradeEvidence(f)
+			return terminal()
+		}
+		code, msg := f.run("upgrade", "--binary", candidate)
+		if code != 1 || !strings.Contains(msg, "changed during confirmation") ||
+			pairedUpgradeEvidence(f) != changed || changed[3] == before[3] || !supervisor.runnerLoaded || f.homeJob.kicks != 0 {
+			t.Fatalf("changed approved pair: %d %s", code, msg)
+		}
+	})
+}
diff --git a/cmd/ailang-worldd/operator_process_darwin.go b/cmd/ailang-worldd/operator_process_darwin.go
new file mode 100644
index 0000000..4264764
--- /dev/null
+++ b/cmd/ailang-worldd/operator_process_darwin.go
@@ -0,0 +1,20 @@
+package main
+
+import (
+	"os"
+
+	"golang.org/x/sys/unix"
+)
+
+func listWorldProcesses() ([]worldProcess, error) {
+	rows, err := unix.SysctlKinfoProcSlice("kern.proc.uid", os.Getuid())
+	if err != nil {
+		return nil, err
+	}
+	out := make([]worldProcess, 0, len(rows))
+	for _, row := range rows {
+		p := row.Proc
+		out = append(out, worldProcess{int(p.P_pid), unix.ByteSliceToString(p.P_comm[:])})
+	}
+	return out, nil
+}
diff --git a/cmd/ailang-worldd/operator_process_other.go b/cmd/ailang-worldd/operator_process_other.go
new file mode 100644
index 0000000..7a30ebd
--- /dev/null
+++ b/cmd/ailang-worldd/operator_process_other.go
@@ -0,0 +1,9 @@
+//go:build !darwin
+
+package main
+
+import "errors"
+
+func listWorldProcesses() ([]worldProcess, error) {
+	return nil, errors.New("operator upgrade requires the macOS process inventory")
+}
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index 36d0041..69cd515 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -415,8 +415,21 @@ func (e serviceEnv) changeRunner(verb string, args []string) error {
 				return err
 			}
 		}
-		if err := e.controller.Bootout(e.domain, runnerLabel); err != nil && !errors.Is(err, errServiceNotLoaded) {
-			return errors.New("could not stop the companion")
+		if r, err := readServiceRecord(e.paths); err == nil {
+			creds, err := store.CanonicalPath(r.Credentials)
+			if err != nil || creds != r.Credentials {
+				return errors.New("cannot locate the companion's canonical credentials lock")
+			}
+			guard, err := e.drainRunner(runnerConfig{Credentials: creds})
+			if err != nil {
+				return err
+			}
+			defer guard.Close()
+		} else {
+			if err := e.stopRunnerJob(min(e.wait, serviceExitTimeout)); err != nil {
+				return err
+			}
+			return errors.New("cannot prove exit without the service record")
 		}
 		if verb == "uninstall" {
 			for _, path := range []string{p.Plist, p.Record} {
diff --git a/cmd/ailang-worldd/operator_upgrade_test.go b/cmd/ailang-worldd/operator_upgrade_test.go
index 359daaa..3a8f73c 100644
--- a/cmd/ailang-worldd/operator_upgrade_test.go
+++ b/cmd/ailang-worldd/operator_upgrade_test.go
@@ -11,6 +11,7 @@ import (
 )
 
 func TestOperatorServiceUpgrade(t *testing.T) {
+	t.Run("active_child", testOperatorUpgradeActiveChild)
 	for _, tc := range []struct {
 		name         string
 		before, next int
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
index e26eaf0..8c12054 100644
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -79,11 +79,12 @@ type serviceEnv struct {
 	client              *http.Client
 	portFree            func(string) error
 	// observe records every intended filesystem mutation, including temp paths.
-	observe     func(string)
-	probe       func(string, string) ([]byte, error)
-	runnerLogin func() (int, string, error)
-	runnerNow   func() time.Time
-	runnerSleep func(time.Duration)
+	observe         func(string)
+	probe           func(string, string) ([]byte, error)
+	runnerLogin     func() (int, string, error)
+	runnerNow       func() time.Time
+	runnerSleep     func(time.Duration)
+	runnerProcesses func() ([]worldProcess, error)
 }
 
 func realServiceEnv() serviceEnv {
diff --git a/cmd/ailang-worldd/service_launchctl.go b/cmd/ailang-worldd/service_launchctl.go
index a6e4108..38a8356 100644
--- a/cmd/ailang-worldd/service_launchctl.go
+++ b/cmd/ailang-worldd/service_launchctl.go
@@ -46,8 +46,12 @@ func launchctlBudget(verb string) time.Duration {
 // scrubs the registry credential (a direct exec here failed the AC10 census).
 // launchctl reports "Could not find service" on stderr, which a non-zero exit
 // carries in its *exec.ExitError; it is joined to stdout as CombinedOutput did.
-func (systemServiceController) command(args ...string) (string, error) {
-	ctx, cancel := budgetContext(launchctlBudget(args[0]))
+func (c systemServiceController) command(args ...string) (string, error) {
+	return c.commandWithin(launchctlBudget(args[0]), args...)
+}
+
+func (systemServiceController) commandWithin(budget time.Duration, args ...string) (string, error) {
+	ctx, cancel := budgetContext(budget)
 	defer cancel()
 	out, err := toolprobe.OutputEnv(ctx, "/bin/launchctl", "", []string{"PATH=/usr/bin:/bin:/usr/sbin:/sbin"}, args...)
 	var exit *exec.ExitError
diff --git a/cmd/ailang-worldd/service_upgrade.go b/cmd/ailang-worldd/service_upgrade.go
index d921fe9..c169e84 100644
--- a/cmd/ailang-worldd/service_upgrade.go
+++ b/cmd/ailang-worldd/service_upgrade.go
@@ -191,7 +191,12 @@ func (e serviceEnv) upgrade(candidate string, pause bool, out, errOut io.Writer)
 		}
 	}
 	recordHash := serviceFileHash(e.paths.Record)
+	runnerReview := e.runnerUpgradeReview()
 	change := fmt.Sprintf("from %s → to %s\ntool AILANG: %s (unchanged)\ninterpreter: %s (unchanged)\nplist: unchanged\nexpected downtime: 10 s drain plus a restart", serviceShort(from), serviceShort(to), v.ToolRelease, v.InterpreterRelease)
+	if runnerReview != [2]string{} {
+		change += "\ncompanion: drain and replace the matched pair; retire older terminal World commands before proceeding"
+		change += "\napproved companion config: " + runnerReview[0] + "\napproved companion plist: " + runnerReview[1]
+	}
 	if active {
 		change += "\n--pause will interrupt: " + episodes
 	}
@@ -203,14 +208,20 @@ func (e serviceEnv) upgrade(candidate string, pause bool, out, errOut io.Writer)
 		return refuse("operator runner lifecycle is busy or unavailable", "inspect service runner status and retry")
 	}
 	defer pairLock.Close()
-	for _, path := range []string{e.runnerPaths().Record, e.runnerPaths().Plist} {
-		if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
-			return refuse("operator runner is still installed", "stop and uninstall the companion; upgrade World, then review and reinstall the runner")
-		}
+	pair, err := e.reviewRunnerUpgrade(runnerReview, from, to, installed.CredentialsSchema, v.CredentialsSchema)
+	if err != nil {
+		return refuse(err.Error(), "inspect service runner status and retry")
 	}
-	if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) {
-		return refuse("cannot verify that the operator runner is unloaded", "inspect service runner status")
+	if pair != nil {
+		if err := e.proveOldTerminalQuiet(); err != nil {
+			return refuse(err.Error(), "retire old World terminal writers, then retry upgrade")
+		}
 	}
+	return e.swapUpgrade(r, source, candidate, from, to, recordHash, installed, v, pair, out, errOut)
+}
+
+func (e serviceEnv) swapUpgrade(r serviceRecord, source []byte, candidate, from, to, recordHash string,
+	installed, v serviceVersion, pair *runnerUpgradePair, out, errOut io.Writer) int {
 	lock, err := e.lock()
 	if err != nil {
 		return serviceMessage(errOut, "a service change is in progress", err.Error(), "wait and retry", 1)
@@ -221,6 +232,17 @@ func (e serviceEnv) upgrade(candidate string, pause bool, out, errOut io.Writer)
 	if serviceFileHash(e.paths.Binary) != from || serviceFileHash(e.paths.Plist) != r.PlistSHA256 || serviceFileHash(e.paths.Record) != recordHash {
 		return serviceMessage(errOut, "service changed during confirmation", "no binary was swapped", "ailang-worldd service status; retry upgrade", 1)
 	}
+	var drained *os.File
+	if pair != nil {
+		drained, err = e.drainRunner(pair.config)
+		if err != nil {
+			return serviceMessage(errOut, "companion drain refused", err.Error(), "wait for process exit, then inspect service runner status and retry", 1)
+		}
+		defer drained.Close()
+		if err := e.proveOldTerminalQuiet(); err != nil {
+			return serviceMessage(errOut, "terminal drain refused", err.Error(), "retire old terminal writers and retry", 1)
+		}
+	}
 	fresh, previous := e.paths.Binary+".new", e.paths.Binary+".prev"
 	if err = e.atomic(fresh, source, 0755); err != nil {
 		return serviceMessage(errOut, "cannot stage upgrade", err.Error(), "inspect "+fresh+"; retry upgrade", 2)
@@ -246,18 +268,38 @@ func (e serviceEnv) upgrade(candidate string, pause bool, out, errOut io.Writer)
 		updated := r
 		updated.BinarySHA256, updated.PreviousSHA256 = to, from
 		if err = e.saveUpgrade(updated, from, to, "ok"); err == nil {
-			fmt.Fprintf(out, "upgraded %s → %s\n", serviceShort(from), serviceShort(to))
-			return 0
+			if pair != nil {
+				err = drained.Close()
+				if err == nil {
+					err = e.writeRunnerPair(pair, false)
+				}
+			}
+			if err == nil {
+				fmt.Fprintf(out, "upgraded %s → %s\n", serviceShort(from), serviceShort(to))
+				return 0
+			}
 		}
 		// A healthy but unrecorded swap is not a successful upgrade.
 	}
 	// A candidate can migrate before reporting health, or while health is
 	// unavailable. The highest schema it could open is the safe comparison.
 	// Older version outputs omitted this field and cannot prove compatibility.
-	return e.rollbackUpgrade(r, from, to, candidate, installed.CredentialsSchema, v.CredentialsSchema, errOut)
+	var rollbackGuard *os.File
+	if pair != nil {
+		drained.Close()
+		rollbackGuard, err = e.drainRunner(pair.config)
+		if err != nil {
+			return serviceMessage(errOut, "paired rollback refused", err.Error(), "verify the companion and children have exited; retain both builds", 3)
+		}
+		defer rollbackGuard.Close()
+		if err := pair.checkSchema(installed.CredentialsSchema); err != nil {
+			return serviceMessage(errOut, "paired rollback refused", err.Error(), "recover the matched pair using a verified schema-compatible build", 3)
+		}
+	}
+	return e.rollbackUpgrade(r, from, to, candidate, installed.CredentialsSchema, v.CredentialsSchema, pair, rollbackGuard, errOut)
 }
 
-func (e serviceEnv) rollbackUpgrade(r serviceRecord, from, to, candidate string, previousSchema, candidateSchema int, w io.Writer) int {
+func (e serviceEnv) rollbackUpgrade(r serviceRecord, from, to, candidate string, previousSchema, candidateSchema int, pair *runnerUpgradePair, guard *os.File, w io.Writer) int {
 	tail := serviceTail(e.paths.Log, 20)
 	failed, previous := e.paths.Binary+".failed", e.paths.Binary+".prev"
 	previousInfo, previousErr := os.Lstat(previous)
@@ -288,6 +330,16 @@ func (e serviceEnv) rollbackUpgrade(r serviceRecord, from, to, candidate string,
 		if recordErr != nil {
 			return serviceMessage(w, "old build is serving but upgrade history could not be saved", recordErr.Error(), "inspect "+e.paths.Record, 3)
 		}
+		if pair != nil {
+			if err := guard.Close(); err != nil {
+				return serviceMessage(w, "old World build restored; companion restart refused", err.Error(), "inspect service runner status", 3)
+			}
+			if err := e.writeRunnerPair(pair, true); err != nil {
+				_ = e.controller.Bootout(e.domain, runnerLabel)
+				return serviceMessage(w, "old World build restored; companion restart refused", err.Error(), "inspect the approved runner configuration and service status", 3)
+			}
+			return serviceMessage(w, "paired upgrade rolled back", "the previous World service and companion configuration were restored", "read the service and runner logs; then retry upgrade", 2)
+		}
 		return serviceMessage(w, fmt.Sprintf("build %s did not answer /v1/health within 30 s; rolled back to %s, which is serving", serviceShort(to), serviceShort(from)), "new build log: "+strings.ReplaceAll(tail, "\n", "; "), "read "+e.paths.Log+"; fix and re-run service upgrade --binary "+candidate, 2)
 	}
 	stopErr := e.controller.Bootout(e.domain, e.label)
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 4ba3f89..78f366a 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -73,14 +73,21 @@ ailang-worldd service runner status
 ailang-worldd service runner stop
 ailang-worldd service runner uninstall
 ```
-Review and type `y` for install, stop and uninstall. Uninstall before upgrading or replacing the config. Spool and action records are retained.
+Review and type `y` for install, stop and uninstall. Uninstall before replacing the config. The
+companion can stay installed during an attended `ailang-worldd service upgrade --binary <file>`.
+Retire older terminal World commands before typing `y`. The upgrade verifies the installed pair
+and both builds' support for the current credentials schema. It unloads the companion and waits for
+its process lock, including any retained children, before swapping the binary. An installed but
+stopped companion stays stopped. If a restart fails, rollback requires a verified compatible pair.
+A drain refusal retains the builds and evidence. After the processes retire, inspect
+`ailang-worldd service runner status` and retry the upgrade. Spool and action records are retained.
 Terminal episode operations use private mode-0600 sidecars named
 `<creds>.operator-episode-<ep>.lock` beside the canonical credentials database.
 Keep these sidecars in place across stop and uninstall. Independent opens need
 the same file to coordinate access to the episode.
 Run the checked example with a fake supervisor and disposable stores:
 ```sh
-AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^TestOperatorServiceWriteScope$' -count=1 -v
+AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./cmd/ailang-worldd -run '^(TestOperatorServiceWriteScope|TestOperatorServiceUpgrade|TestOperatorOldBuildTerminalDrain)$' -count=1 -v
 ```
 
 For an owned start with a committed mint and a missing token, use the
````

# Row 243 V4 S5f2 — operator start activation prerequisites

Iteration 440, planner lane `codex:gpt-6.2-sol`. Base `b5f41458ae483940a8d65bad3cc3ff7df7326f50`. Apply M1, M2, M3 and M4p in order. Appendix C assigns the first production caller and its inseparable lifecycle proof to S5f3.

## 0. Binding constraints

- The reviewed V4 r1 design and quorum r2 PROCEED 3/3 govern this sprint. D-WORLD-119=A and D-WORLD-120=A are attended decisions at `82f05004`. S5f #681 at `886648ab` landed 578 changed lines and received judge Sonnet PASS 86. Attended #678/#679 precede that merge. #680 `0033ca11` and record #682 follow it.
- M1–M3 prepare the exclusion prerequisites. M4p pins dormant publication windows. S5f3 opens the preparation path only with the paired drain and four process pins. S6 owns agent execution. Routes, app, schema, `.ail` and `tools/launchd` are outside the product patches.
- Stores are created by the real migration. Start fixtures use the repository signer. Git repositories, worktrees, private inputs and fake service installations live under `t.TempDir()`. The harness uses process exits, synchronous checkpoints and finite contexts. Compare physical paths with `filepath.EvalSymlinks` on both spellings. No live service, launchctl, home-service files or real credentials are used.
- The product cap is 1,400 changed lines. The 1,650 exception applies only to the complete safe activation unit. Row 267 and open D-WORLD-121 leave measured race allowances of +12 s for `cmd/ailang-worldd` and +5 s for `host/daemon`.

## 1. Pre-flight

### 1.1 Baseline at b5f41458

`git status --short --branch` printed only the sprint branch before opening the two plan files. `git rev-parse HEAD` printed the full base above. `git show -s --oneline 82f05004` returned the attended A rulings. `git log -13 --oneline` placed `6ac6d966` and `419cca9d` before `886648ab`, and `0033ca11` and `b5f41458` after it. The command and output receipt is `$S/anchors-base.txt`. `$S` is `$HOME/.ailang/state/world-iter440/plan-scratch`.

The first gate ran against `$S/base-tree`, a pristine `git archive b5f41458` export created before M1. The baseline raw logs and per-leg statuses are in `$S/gate-base`. Native build/vet, Linux vet, five inventory legs, the GC probe, AILANG verifier and caller census passed. The broad verifier had bind denials. An initial four socket-package runs stopped during compilation on missing Go cache/toolchain artifacts. Those are instrument failures. Use the outside controller captures named in §6 for comparison. The corrected pristine archive gate at `$S/gate-base-final` ended in `ALL_DONE` with Go compile, inventory and AILANG/GC/census receipts. Local daemon and CLI whole-package legs report UNINFORMATIVE UNDER SANDBOX. The controller marker ended in `ALL_DONE`. Its raw `store-race` was clean at 39.184 s and `cmd-plain` had the eleven row-271 failures at 579.844 s. `cmd-race` and `daemon-race` printed `panic: test timed out after 10m0s` at raw lines 43 and 48. Their ten/four preliminary failures are truncated. `verify_go.raw` has four default-timeout markers and an empty attributable fail set with nonzero rc. It is an instrument failure. Section 4.1 records the 20-minute replacement commands. Use `base/cmd-plain.raw` and `base/store-race.raw` as complete outside captures, and the recorded pristine archive plain store and worktree legs for local whole mutations.

The first local whole replay of the blocking-flock omission timed out in its contender fixture. Moving contention and descriptor-count observation into a bounded child corrected the harness. The exact same omission now fails the named test in the untruncated whole store run.

### 1.2 Current inputs and coupling inventory

Every command below was run at the base. Re-run `/bin/bash "$S/anchors.sh" "$W"` for full commands, positive controls and line numbers. `cmd/` means `cmd/ailang-worldd/` in these receipts.

| Input | Command and observed output |
|---|---|
| Branch history | `git diff --name-only 886648ab HEAD -- cmd host` printed `episode.go`, `episode_launch.go`, `episode_preset.go`, `operator_parity_test.go`, `run_cmd{,_test}.go` and `host/verbs/operator.go`. |
| Mode and heartbeat | `git grep -nE 'runnerTick\|signal.Notify\|No effect dispatch\|PutOperatorHeartbeat' -- cmd/ailang-worldd/operator_runner.go` printed startup :42, runner kernel-lock acquisition :54, signal subscription :64, tick :80, heartbeat-only comment :89 and write :94. `main.go:386` is the internal `operator-runner` entry. :90 projects configured preset names. |
| Caller census | `git grep -nE '(withOperatorStartCandidate\|reserveOperatorStartWorktree\|publishOperatorStartFiles\|MintOperatorStart)\(' -- cmd host ':!*_test.go'` printed exactly four declarations, at `operator_start_candidate.go:16`, `operator_worktree.go:31`, `operator_files.go:262` and `host/authority/mint.go:204`. Test controls include `operator_start_candidate_test.go:92`, `operator_worktree_test.go:100`, `operator_files_test.go:150` and `operator_mint_test.go:82`. |
| Process pins | `git grep -nE 'TestOperator(FirstManagedEpisodeProcess\|OldBuildTerminalDrain\|CurrentOwnerProcessRecovery)\|active_child' -- cmd host` exited 1, zero matches. The positive census printed `operator_upgrade_test.go:13 TestOperatorServiceUpgrade`, `operator_service_test.go:302 TestOperatorServiceUpgradePairFence` and `operator_start_candidate_test.go:21 TestOperatorStartCandidateProcess`. |
| J01 | `git grep -nF 'catalogue IS ?' -- host/store/operator_runner.go` printed :64. The schema grep printed `host/store/operator_schema.sql:43: catalogue BLOB`. Append a nil/non-nil behavioral pair using the migrated heartbeat and an injected integer clock. |
| J05 | The candidate read/lock/read grep printed canonical path :18, read-only open :22, deferred close :26, reads :27 and :39, acquire :34 and callback :46. The observer must retain the exact returned credentials pointer and detect closure through SQLite after the visit. |
| Borrowing | `git grep -nE 'openRunnerLock\|defer lock.Close\|worktree.Add\|target.inspect' -- cmd/ailang-worldd/operator_worktree.go cmd/ailang-worldd/operator_files.go` printed independent opens/defers at worktree :40/:44 and files :269/:273, real Git add :90, worktree readback :94 and files' worktree checks :308/:395. `operator_episode_lock.go:14,31,42,48` names the key, path use, noncreating absence and ledger owner read. |
| GC wiring | The `RunEpisodeGC\|HasOpenOperatorStep\|checkWorktreeUnsafe\|removeWorktree\|worktree.*prune` grep printed CLI :107/:113, daemon :79/:85–86/:150/:158–159/:175/:196/:268/:296. `git grep -nF operatorEpisodeLock -- host/daemon/auto_gc.go cmd/ailang-worldd/episode_gc.go` exited 1. Positive controls are `operator_episode_lock.go:14,31`. |
| Lifecycle | The upgrade grep printed `service_upgrade.go:208,212,257,260`, `operator_service.go:418,449,454` and `service.go:197`. Installed/loaded companions are refused ahead of swap. The supervisor allows 30 s for exit. The heartbeat kernel lock is distinct from the episode lock. |
| Executable authority | `git grep -nE 'owner.Uid\|\.Uid\)\|checkRunnerPair\|checkRunnerConfig\|Executables\[\|launch_binary\|Backend' -- cmd/ailang-worldd/{operator_runner,operator_service,operator_start_input,episode_launch,episode_launch_text}.go` printed raw lock ownership at `operator_runner.go:31`, path ownership :71, headless cfg :124, configured executable verification :145–156, pair checks :178–210 in `operator_service.go`, and signed path/digest :74–75 in `operator_start_input.go`. |
| Store boundaries | The exported method grep printed `host/store/operator_actions.go:60,74`, `operator_cas.go:119`, `operator_steps.go:17,45`, `operator_worktree.go:13,73,108`, `operator_mint.go:242`, `operator_files.go:67,87` and `credentials.go:561`. See `anchors-base.txt` for signatures. |
| Terminal scope | `git grep -nE 'terminalEpisodeLock\(\|func runSessionRevoke\|ResolveSession\|RevokeSession' -- cmd/ailang-worldd/{session,session_new}.go` printed mint :269, revoke :361 and :422, and provision :427. `nl -ba cmd/ailang-worldd/session.go` shows the writer open at :405, by-episode path :413–421 and by-ID read and revoke :422–435. |
| CI bound | `git grep -nF -- '-timeout 8m' scripts/verify_go.sh` returned the race gate at :355. The four slow-runner dev attempts are recorded under row 267. Broad local Go package runs use 20 m. |

| Inventory and base location | Staged disposition |
|---|---|
| `TestDaemonDependencyAllowlist`, `host/daemon/daemon_test.go:899` | Same closure. GC reuses existing daemon dependencies on store, verbs and worktree. Run `inv-daemon`. |
| `TestWorldBoundaryDependencyAllowlist`, `host/boundary/allowlist_world_test.go:758` | No snapshot edit. Store adds only stdlib lock primitives. Worktree retains its bounded subprocess site. |
| `TestProductionGoSurface`, `host/store/context_roots_test.go:241` | New lock file is inside `host/store`, already scanned. |
| `TestProductionContextRoots`, `host/store/context_roots_test.go:136` | No production root added. GC passes its existing bounded context to lookup and Git. Named store inventory gate included. |
| `TestCredentialsMethodsRequireADeadline`, `host/store/credentials_test.go:277` | Existing credentials methods and first-argument contexts are preserved. Lock methods operate on a descriptor without database I/O. |
| `TestStoreSurfaceTakesContext`, `host/store/deadline_guard_test.go:43` | No new `Store` or `Credentials` methods. |
| `TestBoundaryASTWriteGuard`, `host/boundary/allowlist_world_test.go:1323` | No snapshot edit. Lock creation is a private sidecar at the store boundary. |
| `TestOperatorNeverDispatch`, `host/verbs/operator_policy_test.go:69` | Keep the `OperatorNever` browser classifications for approver, session mint/revoke, service, serve, setup and world-publish, the `OperatorOther` controls and `ErrShell` arm (:70–88). S5f3 adds a separate internal-mode census and negative process gate. It allows exactly one confirmed signed `start` consumer after runner config, pair and lifecycle admission. Keep `TestOperatorServiceFence` (:127–131) refusing both service-runner start and raw internal mode in the bar, and `TestOperatorRunnerSetupRefusals/global_address` (:70–80) refusing client address injection. Allow only the installed internal-mode queue reader. No command line from the app receives a new dispatch classification. |
| `TestCLICredentialWriterSites`, `cmd/ailang-worldd/cli_credential_writers_test.go:18` | Same writer inventory. The selector opens a read-only handle and GC only reads the ledger. |

### 1.3 Decisions

- M1 observes closure through the original handle, with a per-call opener that returns the real read-only credentials DB. The callback receives no new authority. Nullable SQL catalogue identity is tested on equal NULL and equal non-NULL values, with heartbeat time increasing.
- M2 shares the on-disk key between terminal verbs, the dormant selector and GC. The helper borrows the same verified open description without closing it. The Git child inherits it until exit. Terminal verbs retain their raw open, legacy fallbacks and ownership checks.
- M3 takes that canonical credentials-store-plus-episode lock before the authoritative ledger and fresh Git eligibility checks. The daemon supplies its configured store. CLI `--credentials-db` asserts the store paired with the physical workspace by the service record or default layout. Unknown or mismatched layouts report unavailable evidence. GC confines removal to the registered clean merged checkout. Successful Git removal unregisters that checkout. Remove the unlocked repository-wide prune and recursive-directory fallback, which could operate without the named registration or after losing the episode lock. Keep branch creation-only behavior and existing cache handling.
- Opening the canonical sidecar can create persistent coordination metadata. Tests assert no changes to the action, branch, registration, checked-out tree or private files on refusal, and the contention barrier checks that a refused pass creates no alternate episode key.
- M4p is the independent publication pin slice. The remainder of M4 is parked as a safe unit. The rebase, partial activation receipt and required S5f3 assertions are in Appendix C.

### 1.4 Premises found false or incomplete

| Premise | Measurement |
|---|---|
| The paired upgrader is available at HEAD | The §1.2 source grep shows `service_upgrade.go:208,212` refuse installed or loaded companion. `operator_service.go:418` delegates stop to Bootout. The S5f3 prototype's runner-lock wait proves process exit in a fake stop test. It leaves paired upgrade, compatible pair rollback and old terminal retirement unresolved. |
| Every receipt has an attended remedy | The phase grep gives `host/store/operator_worktree.go:89` reserved only, `operator_mint.go:265,277` observed mint with joined session and start, `operator_files_reconcile.go:32` reserved only. `git grep -nE 'RecoverOperatorClaim\|reconcile-start\|prepared-awaiting-S6\|revoke.*Locked' -- cmd host` exited 1, controlled by the named reconcile functions and claim CAS in §1.2. `$S/prototype.log` captures observed-worktree-only and reserved-mint-abort states for which all three existing commands refuse. |
| The worktree and file effect helpers can already accept the runner's fd | The §1.2 two-open readback shows separate open file descriptions. The real re-exec candidate/terminal control proves those descriptions contend. |
| Stable catalogue generation supplies the complete eligible catalogue | `operator_runner.go:89–90` emits the names slice. Stable generation solves repeated identical heartbeat invalidation. Reviewed §2.2 still requires bounded, eligible, immutable catalogue publication and freshness. |
| The signed digest alone proves executable permission and confinement | `operator_start_input.go:74–75` signs the path and digest. `operator_service.go:145–156` independently verifies the attended executable map. Match both, enforce the supported headless Sol confinement, and require fresh review for older signed envelopes missing the digest. |
| Terminal revoke shares the lock via the mint path | The `session.go:361–437` read shows separate writable revoke entries. The mint and provision controls use the lock elsewhere. Appendix C reserves both revoke routes. |
| An explicit GC credentials path already proves workspace binding | `nl -ba cmd/ailang-worldd/episode_gc.go` prints the unconditional explicit return at :52–54 before any service/default physical workspace check. An unrelated migrated DB with no matching action could clear this gate. M3 pins the conflicting explicit override and the matching positive. |

### 1.5 Size and CI race cost

| Milestone | Increment | Cumulative +/− | Filtered | Unfiltered product | STATUS / DELETED / hashes |
|---|---:|---:|---:|---:|---:|
| M1 | 73 | +66/−7 | 73 | 73 | 3/0/3 |
| M2 | 586 | +605/−42 | 647 | 647 | 12/0/12 |
| M3 | 645 | +1134/−158 | 1292 | 1292 | 17/0/17 |
| M4p | 108 | +1240/−160 | 1400 | 1400 | 18/0/18 |

Product counts are the archive equivalent of `git diff -M --numstat b5f41458ae483940a8d65bad3cc3ff7df7326f50 HEAD -- cmd host docs app` using per-file `git diff --no-index -M --numstat`. Unfiltered product counts inspect every archive file. The two plan records are excluded from all snapshots and reported separately in the JSON. No product renames or deletions.

The sealed paired fresh-archive and `st-M4p` focused receipts use `go test -race -v -count=1 -timeout 90s` with the selectors in `$S/measure-race.sh`. `$S/measure-race-sealed.log` indexes raw files in `$S/timing-v96Y8L`. Sum each top-level test once.

| Package and test | Base s | M4p s | Delta s |
|---|---:|---:|---:|
| CLI `TestOperatorStartCandidateProcess` | 0.33 | 0.16 | -0.17 |
| CLI `TestOperatorStartCandidateRefusals` | 0.61 | 0.41 | -0.20 |
| CLI `TestOperatorStartCandidateClosesCredentials` | — | 0.07 | 0.07 |
| CLI `TestOperatorBorrowedWorktreeGitProcess` | — | 2.74 | 2.74 |
| CLI `TestOperatorFilesReconcile/no_mcp` | 1.49 | 1.13 | -0.36 |
| CLI `TestEpisodeGCFailSafeAndReclaims` | 1.36 | 1.23 | -0.13 |
| CLI `TestOperatorFilePublicationIdentityWindows` | — | 1.33 | 1.33 |
| **CLI total** | **3.79** | **7.07** | **3.28** |
| Daemon `TestEpisodeGCOperatorSteps` | 0.79 | 1.91 | **1.12** |
| Store `TestOperatorCatalogueNullRefresh` and four `TestOperatorEpisodeLock*` tests | — | 0.48 | **0.48** |
| Worktree `TestRemoveRegisteredLendsDescriptor` | — | 0.39 | **0.39** |

Real Git child arms — `/kept` 0.98 s, `/before_git` 0.76 s, `/during_git` 1.00 s.

Publication arms — `/N04` 0.00 s, `/N12` 0.95 s, `/N15` 0.37 s.

New daemon arms — `/active` 0.04 s, `/bad_store_path` 0.03 s, `/invalid_episode` 0.02 s, `/lock_changed_under_review` 0.32 s, `/replaced_during_scan` 0.04 s, `/not_linked` 0.01 s, `/activation_first` 0.03 s, `/recheck_after_ledger` 0.03 s, `/registration_rechecked` 0.06 s, `/unregistered` 0.04 s, `/duplicate_registration` 0.12 s, `/gc_first` 0.38 s.

The CLI total is +3.28 s and the sum of positive CLI deltas is +4.14 s. An earlier paired under-load sample in `measure-race-final.log` measured +11.92 s CLI, +1.84 s daemon. Both paired samples satisfy the +12/+5 s limits.

The final whole-package base/head race seconds were 50.776/67.166 for store (+16.390) and 2.966/2.784 for worktree (−0.182). All four attributable fail sets were empty. The controller's separate outside store base was 39.184 s, illustrating load variation. The socket-package whole-race deltas require the outside 20-minute base and head commands in §4.1. The final local GC probe and AILANG verifier passed.

### 1.6 Portability

Native and Linux target vet are named legs. Child process pipes, Git calls and subprocess commands have explicit bounds. The new Go-only fixtures require no listening socket. Test wall time for a re-exec child excludes the Go race runtime's fixed exit delay using a child-specific `GORACE` setting. The gate must use the existing real Git from PATH, with no PATH shim.

## 2. How the executor applies a sub-milestone

Run with Bash in the sprint worktree. Apply each patch once, in order. The controller may already have committed an earlier snapshot.

````bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter440-row243-v4s5f2
STATE=$HOME/.ailang/state/world-iter440
S=$STATE/plan-scratch
T=$STATE/executor
B=b5f41458ae483940a8d65bad3cc3ff7df7326f50
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s5f2-sprint-plan.md
M=M1  # For the next unapplied milestone set M2, M3 or M4p.
mkdir -p "$T"
unset SNAP_PARENT SNAP_STATUS_LIST
git -C "$W" merge-base --is-ancestor "$B" HEAD
if [ "$M" = M1 ]; then
  git -C "$W" diff --quiet "$B" -- cmd host docs app
  test -z "$(git -C "$W" ls-files --others --exclude-standard -- cmd host docs app)"
fi
cp "$S"/{bounded.py,failset.py,verify-scratch.py,callers.sh,probe.sh,probe-fixture.py,operator_gc_after_test.go} "$T/"
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

## 3. Milestones

### M1 — Judge 439 J01 and J05

`TestOperatorCatalogueNullRefresh/null` inserts SQL NULL using the migrated heartbeat API and proves an identical subsequent catalogue leaves catalogue generation fixed while liveness advances. `/blob` is the non-NULL twin. `TestOperatorStartCandidateClosesCredentials` captures the real read-only credentials pointer, visits one confirmed row, then proves the same pointer has been closed. Preserve the read/lock/read order and unchanged action bytes. The exact original edits are replayed independently.

### M2 — Borrow the per-episode descriptor

`TestOperatorEpisodeLockIdentity`, `TestOperatorEpisodeLockAcquireRefusals`, `TestOperatorEpisodeLockFIFORefused` and `TestOperatorEpisodeLockContendsWithoutWaiting` prove the canonical key, kernel contention, mode/owner/type/link and descriptor lifecycle. `TestOperatorStartCandidateProcess` retains the real terminal process and compares the physical identities of its two path spellings. `TestOperatorBorrowedWorktreeGitProcess/kept` proves the real Git hook sees the same inode and open description using FD 3 and an independent flock. `/before_git` and `/during_git` pin substitution at the two effect boundaries, retaining the correct reservation and executable remedy. `TestOperatorFilesReconcile/no_mcp` rejects a missing borrow with complete snapshot equality and retains the real publication control. Terminal callers and attended reconciliation keep their original lock policy.

### M3 — Exclude activation and GC in either order

`TestEpisodeGCOperatorSteps/activation_first` holds the canonical lock via a symlinked store and enters the daemon pass. `/gc_first` re-executes an independent contender at the removal boundary before releasing the descriptor. The fixtures compare refs, Git registration, tracked bytes and the worktree status on refusal. `/recheck_after_ledger`, `/replaced_during_scan`, `/registration_rechecked` and `/lock_changed_under_review` check drift after enumeration. `TestEpisodeGCFailSafeAndReclaims` drives the actual CLI with a paired record, a wrong store, both workspace spellings, contention, dirty and unmerged protections and a successful scoped removal. `TestRemoveRegisteredLendsDescriptor` verifies that the actual Git remove receives the borrowed kernel description and preserves the branch. A failed Git removal yields a diagnostic and no broader cleanup.

### M4p — Optional publication race pins

`TestOperatorFilePublicationIdentityWindows/N04` swaps the readback directory before the second real open. `/N12` moves the private directory immediately before first publication under a fully minted start. `/N15` moves the verified worktree after file readback. Both publication cases assert the precise uncertain receipt, unchanged Git and file evidence, and the same mint. The signed mint fixture is reset from its real migrated bytes between cases. The per-call opener seam preserves attended behavior.

### 3.1 Reviewed §5.3 branch assignment

| Branch trace | Sprint delivery |
|---|---|
| Worktree path/branch collisions, partial Git result and readback | M2 pins descriptor provenance and time of borrow. Existing creation-only, readiness and readback tests remain. |
| Same episode becomes eligible for GC or active while inspected | M3 holds the lock across final eligibility and removal. Scan-time identity changes produce a skip. |
| Runtime pair, signed inputs, current approval or session expires/revokes between preparation effects | S5f3 admission, CAS and per-effect recheck in Appendix C. |
| Mint rollback and partial private writes | S5b–S5e3 fixtures establish the local receipt operations. Cross-effect recovery and hand-off are S5f3. |
| Spawn, child outcome, browser disconnect, TTY hand-off | S6, S8 and S12 in the reviewed split. |

### 3.2 Reviewed §5.5 and §2.2 lifecycle assignment

M1–M4p leave the production runner heartbeat-only. Appendix C requires paired build/schema identity, drain admission closure, old terminal retirement and current-owner recovery before the first consumer is reachable.

## 4. Mutation table

The listed mutations are per owning stage only. Every substitution is literal, anchored to a stage file and line, and measured against its behavioral witness. Its complete source, control, raw log, command and classification are printed in the separate diff-anchored and carried-id JSONL ledgers. A named kill needs the relevant `--- FAIL:` witness and no bind denials. For whole packages, compare the failure set against an attributable pristine baseline and require the expected added failure. Compiler errors, a missing named witness, and truncated logs are instrument failures.

| Ledger | Edits | Named killed | Named equivalent | Whole killed | Whole equivalent | Whole unmeasured / sandbox |
|---|---:|---:|---:|---:|---:|---:|
| Diff-anchored ledger | 95 | 94 | 1 | 30 | 1 | 64 |
| Carried-id ledger | 7 | 7 | 0 | 1 | 0 | 6 |

Exact substitutions and raw logs are in `$S/ledger-named-diff.jsonl`, `$S/ledger-named-carried.jsonl`, `$S/ledger-whole-diff.jsonl` and `$S/ledger-whole-carried.jsonl`. The `stage` in each row is the cumulative archive holding both that milestone's patch and its behavioral witness. The sole named equivalent, `D-LOCK-NO-NONBLOCK`, drops `O_NONBLOCK` on open of a regular single-link sidecar. It has no effect on accepted regular files; the flag remains as defense in depth for abnormal device paths. `D-LOCK-FLOCK-NONBLOCK` separately kills removing `LOCK_NB` from the kernel lock attempt with a bounded real contender.

Base J01 with its original edit passed the whole store package in 28.722 s (`$S/base-J01.raw`), then failed only `TestOperatorCatalogueNullRefresh/null` on M1. Base J05 with its original edit passed `TestOperatorStartCandidate{Process,Refusals}` in 0.596 s (`$S/base-J05.raw`), then failed `TestOperatorStartCandidateClosesCredentials` on M1. The whole CLI run under the sandbox reported bind denials. Act on the staged named kill and reproduce the whole comparison outside.

`J-OWN-UID` belongs to Appendix C. The original amended edit is unavailable (`ls -ld ~/.ailang/state/world-iter429 ~/.ailang/state/world-iter430` returned absent). It contributes zero to the staged denominators. The typed descriptor's owner check is covered by `D-LOCK-OWNER`; the existing raw runner and configuration checks remain future coverage.

`$S/whole-mutant-ids.txt` contains the 70 pending whole comparisons. A captured sandbox denial is recorded under the same ID. See §4.1 for the outside command and base.

### 4.1 Whole-package controller runs

PENDING — controller to measure outside the sandbox with the commands below. The controller baseline marker now ends in `ALL_DONE`. Run from the sprint checkout. Run the commands with Bash and an isolated scratch HOME created by the scripts. The base gate also supplies daemon plain and replaces the truncated daemon-race, CLI-race and broad verifier captures with 20-minute legs. The controller's recorded `base/{store-race,daemon-race,cmd-plain,cmd-race,verify_go}` legs take precedence when they have complete parser and raw receipts.

````bash
STATE=$HOME/.ailang/state/world-iter440
S=$STATE/plan-scratch
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter440-row243-v4s5f2
test "$(tail -1 "$STATE/baseline.done")" = ALL_DONE
WORLD_PLAN_FULL_GO=1 /bin/bash "$S/gates.sh" "$S/base-tree" "$S/outside-gate-base" base "" "$W"
WORLD_PLAN_FULL_GO=1 /bin/bash "$S/gates.sh" "$S/st-M4p" "$S/outside-gate-M4p" M4p "$S/outside-gate-base" "$W"
/bin/bash "$S/base-j05-whole.sh"
/bin/bash "$S/whole-mutants.sh" "$S/whole-mutant-ids.txt" "$S/outside-gate-base"
````

The outside script derives a separate comparable baseline directory from the raw base legs and controller captures. A valid comparison uses the strict plain `go test -count=1 -timeout 20m` baseline for the package holding the killer. The whole mutated command uses `-v` to retain the exact subtest witness. No raw log with a bind denial, timeout or build failure can report KILLED. The ID file includes every unmeasured whole substitute. Stage provenance is assigned in `mutants.json`. Controller can run subsets by passing an explicit IDs file of the same format.

Locally measured whole package results for store and worktree compare against `$S/gate-base-final`; M4p and base package guards recorded zero attributable test failures. The socket packages are uninformative in the sandbox. `$S/whole-J05.log`, `$S/base-j05-whole.log`, `$S/whole-daemon-sample.log` and `$S/gate-head-final/legs.txt` show representative bind denials. Use the commands above to complete CLI and daemon plain/race, the broad verifier and the outside mutant ledger. New and dropped failures are recorded individually.

## 5. Files

| Path | M4p cumulative +/− |
|---|---:|
| `cmd/ailang-worldd/episode_gc.go` | +14/−13 |
| `cmd/ailang-worldd/episode_gc_test.go` | +17/−4 |
| `cmd/ailang-worldd/operator_episode_lock.go` | +1/−1 |
| `cmd/ailang-worldd/operator_file_identity_test.go` | +99/−0 |
| `cmd/ailang-worldd/operator_files.go` | +9/−6 |
| `cmd/ailang-worldd/operator_files_test.go` | +22/−5 |
| `cmd/ailang-worldd/operator_start_candidate.go` | +4/−4 |
| `cmd/ailang-worldd/operator_start_candidate_test.go` | +47/−14 |
| `cmd/ailang-worldd/operator_worktree.go` | +12/−8 |
| `cmd/ailang-worldd/operator_worktree_guards_test.go` | +127/−1 |
| `cmd/ailang-worldd/operator_worktree_test.go` | +12/−1 |
| `host/daemon/auto_gc.go` | +134/−85 |
| `host/daemon/auto_gc_s5e4_test.go` | +300/−14 |
| `host/store/operator_episode_lock.go` | +73/−0 |
| `host/store/operator_episode_lock_test.go` | +270/−0 |
| `host/store/operator_runner_test.go` | +29/−0 |
| `host/worktree/worktree.go` | +12/−4 |
| `host/worktree/worktree_test.go` | +58/−0 |

Per-file base and stage hunk anchors are in Appendix P. `$S/numstat-final.txt` records every staged path and count.

## 6. Sprint gates

The raw test log is authoritative. Every parser emits a firing control and every raw log is checked for `panic: test timed out`. A timed-out package is truncated. Record all legs before comparison. `SKIPPED` legs state their dependency. Socket-denied captures are UNINFORMATIVE UNDER SANDBOX. The native build, vet, Linux vet, five inventory legs, named acceptance, package plain/race, AILANG verify, GC probe and production-caller census all have separate receipts. The final base and head are `$S/gate-base-final` and `$S/gate-head-final`. Native formatting/build/vet, Linux vet and the named acceptance were also checked per stage in `$S/check-final-M*.log` and `$S/accept-final-M*.log`. Whole package comparisons use recorded sets. A baseline-only failure disappearance is reported without being called a regression.

````bash gates.sh
#!/bin/bash
# gates.sh TREE LOGDIR base|M1|M2|M3|M4p [BASE-LOGDIR] [READ-ONLY-GIT-REF-TREE]
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
mkdir -p "$2" || exit 2
L=$(cd "$2" && pwd -P)
M=$3
BASE_LOG=${4:-}
REF=${5:-}
case "$M" in base|M1|M2|M3|M4p) ;; *) exit 2;; esac
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
  local path
  for path in host/store/operator*.go host/daemon/auto_gc*.go cmd/ailang-worldd/operator*.go cmd/ailang-worldd/episode_gc*.go host/worktree/worktree*.go; do
    [ ! -f "$path" ] || gofmt -l "$path" || return
  done
}
check_fmt() { local out; out=$(fmt) || return; printf '%s\n' "$out"; [ -z "$out" ]; }
{ uptime; ps -A -o pid,pcpu,comm -r | head -12; go version; command -v git; } > "$L/load.txt" 2>&1
leg gofmt check_fmt
leg build bounded 240 go build ./...
leg vet bounded 240 go vet ./host/store ./host/daemon ./cmd/ailang-worldd ./host/worktree
leg linux-vet env GOOS=linux python3 "$S/bounded.py" 240 go vet ./host/store ./host/daemon ./cmd/ailang-worldd ./host/worktree
leg inv-store inventory host/store 'TestProductionContextRoots|TestProductionGoSurface|TestStoreSurfaceTakesContext|TestCredentialsMethodsRequireADeadline'
leg inv-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
leg inv-cli inventory cmd/ailang-worldd 'TestCLICredentialWriterSites|TestEveryVerbHasHelp'
leg inv-verbs inventory host/verbs TestOperatorNeverDispatch
leg inv-daemon inventory host/daemon TestDaemonDependencyAllowlist
if [ "$M" = base ]; then skip accept 'new pins absent at base'
else leg accept env WORLD_PLAN_PIN_HOME="$ORIGINAL" /bin/bash "$S/accept.sh" "$T" "$M" fast; fi
leg gc-probe env HOME="$ORIGINAL" /bin/bash "$S/probe.sh" "$T"
if [ "$M" = base ] || [ "$M" = M4p ]; then
  leg verify_ail bounded 510 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
  if [ -n "$REF" ] && [ "${WORLD_PLAN_FULL_GO:-0}" = 1 ]; then
    leg verify_go record verify_go bounded 1500 python3 "$S/verify-scratch.py" "$T" "$L" verify_go "$REF"
  else skip verify_go 'set WORLD_PLAN_FULL_GO=1 with the read-only reference'; fi
else
  skip verify_ail 'AILANG tree byte-identical; base and M3 carry the verifier'
  skip verify_go 'pristine and final broad gate'
fi
leg callers /bin/bash "$S/callers.sh" "$T" "$L/callers.census" "$M"
for pkg in host/store host/worktree host/daemon cmd/ailang-worldd; do
  name=${pkg//\//-}
  if { [ "$M" = M1 ] && { [ "$pkg" = host/daemon ] || [ "$pkg" = host/worktree ]; }; } ||
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
    race-host-store) outside=store-race;;
    race-host-daemon) outside=daemon-race;;
    pkg-cmd-ailang-worldd) outside=cmd-plain;;
    race-cmd-ailang-worldd) outside=cmd-race;;
    verify_go) outside=verify_go;;
  esac
  if [ -n "$outside" ] && [ -f "$controller/$outside.fails" ] &&
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
for tag in pkg-host-store race-host-store pkg-host-worktree race-host-worktree pkg-host-daemon race-host-daemon pkg-cmd-ailang-worldd race-cmd-ailang-worldd verify_go; do
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

## 7. Executor directive notes

````bash accept.sh
#!/bin/bash
# accept.sh TREE M1|M2|M3|M4p fast|race
set -euo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
M=$2
MODE=$3
PINS=${WORLD_PLAN_PIN_HOME:-$HOME}
R=$(mktemp -d "$S/accept-XXXXXX")
export AILANG_BIN="$PINS/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$PINS/.pinned-ailang-tools/v0.54.15/ailang"
export GOCACHE="$(go env GOCACHE)" GOMODCACHE="$(go env GOMODCACHE)"
export HOME="$R/home" AILANG_STATE_DIR="$R/state" AILANG_CACHE_DIR="$R/cache"
export XDG_CONFIG_HOME="$HOME/.config" XDG_CACHE_HOME="$HOME/.cache" XDG_DATA_HOME="$HOME/.local/share"
mkdir -p "$HOME" "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
cd "$T"
opts=(-v -count=1 -timeout 90s)
case "$MODE" in
  fast) ;;
  race) opts+=(-race) ;;
  *) exit 2 ;;
esac
check() {
  local pkg=$1 selector=$2 witness=$3 result="$R/${1##*/}-${3%%/*}.raw" rc=0
  python3 "$S/bounded.py" 120 go test "${opts[@]}" "./$pkg" -run "$selector" > "$result" 2>&1 || rc=$?
  cat "$result"
  if grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$result"; then
    printf 'UNINFORMATIVE UNDER SANDBOX %s\n' "$result"
    return 77
  fi
  if grep -Fq 'panic: test timed out' "$result"; then
    printf 'INSTRUMENT FAILURE: truncated %s\n' "$result"
    return 2
  fi
  [ "$rc" -eq 0 ] && grep -Fq -- "--- PASS: $witness (" "$result"
}
check host/store '^TestOperatorCatalogueNullRefresh$' TestOperatorCatalogueNullRefresh/null
check cmd/ailang-worldd '^TestOperatorStartCandidate(ClosesCredentials|Process|Refusals)$' TestOperatorStartCandidateClosesCredentials
case "$M" in M1) exit 0;; M2|M3|M4p) ;; *) exit 2;; esac
check host/store '^TestOperatorEpisodeLock(Identity|AcquireRefusals|FIFORefused|ContendsWithoutWaiting)$' TestOperatorEpisodeLockIdentity
check cmd/ailang-worldd '^TestOperatorBorrowedWorktreeGitProcess$' TestOperatorBorrowedWorktreeGitProcess
check cmd/ailang-worldd '^TestOperatorFilesReconcile$/^no_mcp$' TestOperatorFilesReconcile/no_mcp
check host/worktree '^TestAdd(AndRemoveAWorktree|WithBranch)$' TestAddAndRemoveAWorktree
if [ "$M" = M3 ] || [ "$M" = M4p ]; then
  check cmd/ailang-worldd '^TestEpisodeGCFailSafeAndReclaims$' TestEpisodeGCFailSafeAndReclaims
  check host/daemon '^TestEpisodeGCOperatorSteps$' TestEpisodeGCOperatorSteps/gc_first
  check host/worktree '^TestRemoveRegisteredLendsDescriptor$' TestRemoveRegisteredLendsDescriptor
fi
if [ "$M" = M4p ]; then
  check cmd/ailang-worldd '^TestOperatorFilePublicationIdentityWindows$' TestOperatorFilePublicationIdentityWindows/N15
fi
````

````bash snapshot.sh
#!/bin/bash
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1|M2|M3|M4p) ;; *) exit 2;; esac
B=b5f41458ae483940a8d65bad3cc3ff7df7326f50
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter440/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s5f2-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s5f2.json"
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

The executor runs the §2 block with `B=b5f41458ae483940a8d65bad3cc3ff7df7326f50` for each milestone. `snapshot.sh` records every cumulative product file relative to B, including untracked package files, and copies them under `$HOME/.ailang/state/world-iter440/snap/<M>/`. `STATUS`, `DELETED` and `MANIFEST.sha256` use repository-relative paths. The two plan records are omitted. No target directory outside the changed path list is required.

The controller records the two plan files before copying any snapshot, since `commit_ms.sh` refuses untracked `design_docs` paths outside the snapshot. For each accepted milestone, use the following command from Bash after creating its commit message at `$T/$M.msg`. This runs the named suites with fake stores and no loopback.

````bash
M=M1  # Repeat in stage order with M2, M3 and M4p.
STATE=$HOME/.ailang/state/world-iter440
T=$STATE/executor
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter440-row243-v4s5f2
test -s "$T/$M.msg"
printf -v commit_test '%q ' env "AILANG_BIN=$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node "WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang" /bin/bash "$T/accept.sh" "$W" "$M" fast
/bin/bash "$STATE/commit_ms.sh" "$M" "$T/$M.msg" "$commit_test"
````

Rehearsal `$S/rt-OuAPrE` extracted the four patches from this plan with the controller's exact awk. A fresh archive replay returned `diff -r` rc 0 at each of M1, M2, M3 and M4p. Snapshot STATUS / DELETED / hashes were 3/0/3, 12/0/12, 17/0/17 and 18/0/18. Each SHA-256 manifest verified. The transcript is `$S/roundtrip-final.log`.

The final extracted `accept.sh` also accepted the replayed M4p archive. The separately extracted gate and snapshot scripts passed Bash syntax checks and matched the measured script bytes.

The controller round trip uses `/bin/bash "$STATE/ctl_roundtrip.sh" M1 M2 M3 M4p` once, with the default controller output path. Planner rehearsal uses only a unique `$S/rt-*` directory.

## 8. Risks and pre-existing failures

- Row 271 on this rig since `f8fb3ac1`. CLI has eleven known reds (`TestCLIAgainstSeToolsDaemon`, `TestCLINoGenesisIsNoWorldHead`, eight `TestDemoSeed*`, `TestMarketingStillRoutes`). Daemon has `TestSeToolsDescriptorLiveExactEight`, `TestSeToolsMCPVisibilityFollowsGrants` and `TestSeToolsA2AEndToEnd`.
- Load-sensitive test classes are `TestWorldContextAndStatusline`, `host/broker TestExecProtectedCheckTimeout`, `host/daemon TestSeToolsRunBenchmarkTasksEndToEnd`, `host/verifygate TestMissionConfigDefaultFlowIntegration` and `host/broker TestSeedRestartIgnoresASurvivingOldCopy`. Use the controller's recorded baseline. Keep these tests out of commit-time acceptance.
- S5f3 must address the pre-existing ability of an older terminal binary to bypass the episode lock. A stale heartbeat or absent sidecar alone cannot prove quiescence. Drain and restart limits are bounded by the attended supervisor contract.
- A different credentials DB with a readable, empty operator ledger is semantically unrelated to the workspace. CLI stops guessing for an unconfigured custom layout. Operators can use the daemon with its configured pair or supply the matching service/default workspace and store.

## 9. Acceptance criteria

1. Each stage is cumulative, within its measured budget, formatted, compiled and vetted. The plan-to-export replay and SHA-256 snapshot rehearsal match every staged product byte.
2. J01 and J05 survive their original measured base selectors and fail their new named pins with exactly the original edits.
3. The two preparation helpers operate under the caller's exact descriptor. Real Git observes that descriptor across its effect and the calling code does not issue a receipt on inode mismatch.
4. Both GC entry points use the common physical key, fresh evidence and an explicit skip reason. Neither a protected episode nor an unverifiable checkout is reclaimed.
5. `runnerTick` remains heartbeat-only at the staged HEAD. The controller verifies the outside baseline comparisons before product commit. S5f3 carries the first preparation caller plus the lifecycle and process proof below.

## Appendix C — S5f3 activation unit

The first production caller is gated on the complete M4 safe unit. The scratch rebase of the read-only iteration 439 `proto-S5f2` is at `$S/st-S5f3-measured`, using the M2 borrowed descriptor and the M3 GC fence. Its five inherited prototype source/test files were copied to scratch before rebase. `/bin/bash "$S/prototype-size.sh"` reproduced the inherited +1,057/−25 = 1,082 lines against iteration 439's base. For the rebased unit, `/bin/bash "$S/measure-safe-cut.sh"` measured **+410/−0** for the partial caller, prepared-effect wiring, stop adapter and inherited prototype tests, **+570/−14 = 584** with the prototype revoke and UID probes, and **+1,810/−174 = 1,984** cumulative against b5f41458 (filtered and unfiltered). Excluding M4p's 108 lines the same measured incomplete unit is 1,876. Even the 410-line core plus M1–M3 is 1,702. The partial rebase already exceeds 1,650 while the paired transition and the required process/reconcile proofs are still missing. The line count measures compiled code, including its incomplete tests. The full safe unit also requires the missing stages below. The rebase is excluded from Appendix P, snapshots and the product caller census. Native build, native vet and Linux-target vet succeeded. The scoped prototype preparation, drain, revoke and alternate-UID tests passed. `prototype-final.log` records 9.19 s preparation, 0.08 s synthetic drain, 0.64 s revoke, and 0.00 s alternate UID under race instrumentation, or 9.91 s of prototype test bodies. `TestS5f2PrototypePreparation/observed_worktree_without_mint` and `/reserved_mint_rollback` both reproduced checkpoints where the three installed remedies refused. The prior prototype assumed every heartbeat invalidated a catalogue; changing the receipt fixture's actual catalogue bytes restored the intended stale-generation witness at this base. The prototype's activated call chain is `runOperatorRunner → consumeOperatorPreparation → prepareOperatorStart → withOperatorStartCandidate`, with one direct `authority.MintOperatorStart` call. The packaged sprint head has neither production caller.

The reviewed sequence from §4.2 and §§3.2–3.3 is one confirmed start per bounded tick, under a single canonical episode lock from before CAS through the final checkpoint. Re-read pair, heartbeat generation, canonical workspace and fresh approval under the lock. Bind the signed executable path and digest to the attended entry and supported confinement before claim and each new effect. An older signed start missing `launch_binary_sha256` receives a stale-binding result, zero claimed state and a fresh attended review. Never synthesize a digest into an already signed envelope. Verify the current credential ID, metadata, expiry boundary, and revocation between mint and each file boundary. Prepare worktree by exclusive reservation and Git inspection, mint exactly once with its atomic start receipt, then publish and verify private files. End at durable prepared-awaiting-S6. Never launch or reserve an agent spawn in M4.

| Durable window under the kernel and episode locks | Required S5f3 disposition |
|---|---|
| Claimed with no side-effect reservation | Prove owner/process quiescent and full absence of steps/effects. Exact old-owner CAS can release only this untouched claim to confirmed if proof is fresh. Expired or revoked proof requires attended disposition and fresh review. |
| Worktree reserved, child status or Git evidence ambiguous | Retain the reservation and invoke the existing reconcile-worktree command only after process quiescence; inspect exact registration, branch, physical path and pinned commit. |
| Observed worktree, no mint reservation | Under old-owner CAS verify the exact observed worktree, absence of every later receipt, fresh proof, and schema. A permitted prepared continuation or terminal disposition must preserve the tree; no Git add retry. |
| Reserved mint without commit | Require durable proof of rollback, absent start/session/observation and exact matching owner/version. One in-place retry with renewed authority or explicit attended disposition. An ambiguous commit goes uncertain with a runnable command. |
| Observed mint, no files reservation | Read current session and the atomic receipt. With no attributable recoverable token, revoke only the credential from that receipt and stop with a repair path, preserving the start and worktree. Never generate a second credential for the action. |
| Partially published files | Existing reserved-file reconcile handles verified files and returns named uncertainty for collisions. Verify all paths, owner, mode, identity and hashes. Preserve evidence; revoke only the orphan if unrecoverable. No overwrite, blind retry or plaintext journal token. |
| All prep observed, awaiting S6 | Derive the durable status from the matching worktree, mint and file receipts and owned action. A migration requires independent evidence that this representation cannot encode the reviewed state. Re-read under lock and prove restart can adopt the fully prepared checkpoint with no second effect. Successful launch remains an S6 transition. |

| Pin carried to S5f3 | Required process/evidence proof |
|---|---|
| `TestOperatorFirstManagedEpisodeProcess` | Real runner paused before CAS; a real terminal holder defeats it and leaves confirmed with zero reservations. Reverse order, assert the runner's one prepared receipt and terminal/GC refusals. |
| `TestOperatorOldBuildTerminalDrain` | Re-exec the older lock-bypassing terminal semantics. Refuse new runner admission until the old writer has actually retired. A fake supervisor result or stale heartbeat cannot authorize mutation. Check binary, record, branch, DB and file equality after refusal. |
| `TestOperatorCurrentOwnerProcessRecovery` | With an independent current owner, observe live evidence or refuse in both fresh and stale heartbeat cases. Once owner and child quiescence are proved, re-read the signed row and current-session data. Pin expiry, revocation and identity independently. |
| `TestOperatorServiceUpgrade/active_child` | Bounded re-exec with a real active child and fake supervisor. A typed approval is required. Pin admission closure, effect finish/readback or uncertain and an executable remedy, matched binary/config/plist restart and schema-compatible rollback. Preserve pending evidence and both builds on an unproved exit. |
| `J-OWN-UID` | Judge 429's corrected before/after text is unavailable. Use a real file stat and a different expected login UID for the `runnerPath` acceptance check. The original raw lock guard requires a real wrong-owner fixture or a UID seam pinned to the process login. M2 separately tests the new descriptor UID comparison. The reconstructed `runnerPath` comparison removal failed `TestS5f3PrototypeRunnerPathUID` in a scratch-only probe (`$S/j-own-reconstruction.raw`). It remains separate from the original carried edit. |
| Terminal `session revoke` | By-ID read-only lookup, then lock then writer re-read before revoke; by-episode locks the exact supplied episode before writer/migration. Both routes refuse an active holder or durable owner and preserve unrelated credentials, branch, worktree and private files. Keep legacy-store handling and pin stale lookup. |
| Optional N04/N12/N15 | Built independently as M4p. The exact originals fail `TestOperatorFilePublicationIdentityWindows/N04`, `/N12`, and `/N15`. M4p's 108 changed lines fit the ordinary cap with M1–M3. |

On SIGTERM or an attended upgrade, stop admitting before the next CAS. Drain retained Git children or retain an attributable uncertain receipt after the bounded wait. Keep the drain budget within the existing 30-second supervisor limit. The approval deadline limits the next effect. An already-authorized child and its readback must finish or retain its exact uncertain reservation. Upgrade must stop and positively prove companion and old terminal quiescence, check signed configuration/build/store identity under lifecycle exclusion, replace and restart a matched pair, and roll back only to a pair compatible with the actual credentials schema. Keep process-heavy pins in short pipe-controlled re-execs and place the pure disposition matrix in store/authority. These checks are part of admission, including any context and writer inventory updates they require. The current browser dispatch refusal inventory remains exact.

Place the eligible-catalogue projection and bounded one-per-tick admission beside `operator_runner.go:80`. Reuse the claim CAS and transactional receipt methods in `host/store/operator_cas.go:119`, `operator_steps.go:17,45`, and `operator_mint.go:242`. Keep the durable disposition predicates with the store and authority proofs. Keep paired stop/readback/rollback in `service_upgrade.go:201–260` and `operator_service.go:398–429`. Preserve all four named process arms as activation gates and reuse the existing process-fixture seams. Test signed legacy compatibility with an artifact from the real signer over the older resolution, then compare both the untouched row and the absent effects on stale-binding refusal.

Reproduce the measured size now with `/bin/bash "$S/measure-safe-cut.sh"`. Complete the remaining paired drain and recovery before measuring the enabling stage. Its required focused selector is `go test -race -v -count=1 -timeout 90s ./cmd/ailang-worldd -run '^(TestOperatorFirstManagedEpisodeProcess|TestOperatorOldBuildTerminalDrain|TestOperatorCurrentOwnerProcessRecovery|TestOperatorServiceUpgrade)$'`. Require real `--- PASS:` witnesses for all four, including `TestOperatorServiceUpgrade/active_child`, then rerun the GC and carried ledgers. The full enabling stage remains unmeasured until it contains those named tests and the runnable attended remedies.

## Appendix P — compiled patches

### Patch M1

````diff M1
diff --git a/cmd/ailang-worldd/operator_start_candidate.go b/cmd/ailang-worldd/operator_start_candidate.go
index 944267bb..59682606 100644
--- a/cmd/ailang-worldd/operator_start_candidate.go
+++ b/cmd/ailang-worldd/operator_start_candidate.go
@@ -14,12 +14,13 @@ import (
 // visit. The runner may wire a writer here only once the process and recovery
 // gates are met. The visit must use a fresh clock and the locked row for CAS.
 func withOperatorStartCandidate(ctx context.Context, credentials string, now func() int64,
-	acquire func(string) (*os.File, error), visit func(*os.File, store.OperatorActionRow) error) error {
+	open func(string) (*store.Credentials, error), acquire func(string) (*os.File, error),
+	visit func(*os.File, store.OperatorActionRow) error) error {
 	canonical, err := store.CanonicalPath(credentials)
 	if err != nil {
 		return err
 	}
-	db, err := store.OpenCredentialsReadOnly(runnerDB(canonical))
+	db, err := open(runnerDB(canonical))
 	if err != nil {
 		return err
 	}
diff --git a/cmd/ailang-worldd/operator_start_candidate_test.go b/cmd/ailang-worldd/operator_start_candidate_test.go
index 4aac89c0..8b415bf3 100644
--- a/cmd/ailang-worldd/operator_start_candidate_test.go
+++ b/cmd/ailang-worldd/operator_start_candidate_test.go
@@ -89,7 +89,7 @@ func TestOperatorStartCandidateProcess(t *testing.T) {
 		}
 	}
 	now := func() int64 { return 1012 }
-	err = withOperatorStartCandidate(ctx, aliasPath, now, runnerKernelLock, func(*os.File, store.OperatorActionRow) error {
+	err = withOperatorStartCandidate(ctx, aliasPath, now, store.OpenCredentialsReadOnly, runnerKernelLock, func(*os.File, store.OperatorActionRow) error {
 		t.Fatal("visit ran while the real terminal process held the lock")
 		return nil
 	})
@@ -108,7 +108,7 @@ func TestOperatorStartCandidateProcess(t *testing.T) {
 	}
 	handled := errors.New("test visit stopped")
 	var acquired *os.File
-	err = withOperatorStartCandidate(ctx, aliasPath, now, func(path string) (*os.File, error) {
+	err = withOperatorStartCandidate(ctx, aliasPath, now, store.OpenCredentialsReadOnly, func(path string) (*os.File, error) {
 		var err error
 		acquired, err = runnerKernelLock(path)
 		return acquired, err
@@ -178,7 +178,7 @@ func TestOperatorStartCandidateRefusals(t *testing.T) {
 			treeBefore, filesBefore := provisionFiles(t, worktree), provisionFiles(t, files)
 			now := int64(1012)
 			var acquired bool
-			err = withOperatorStartCandidate(boundedTestContext(t), creds, func() int64 { return now },
+			err = withOperatorStartCandidate(boundedTestContext(t), creds, func() int64 { return now }, store.OpenCredentialsReadOnly,
 				func(path string) (*os.File, error) {
 					acquired = true
 					if tc.advance {
@@ -205,7 +205,7 @@ func TestOperatorStartCandidateRefusals(t *testing.T) {
 		root := t.TempDir()
 		before := provisionFiles(t, root)
 		missing := filepath.Join(root, "absent", "credentials.db")
-		err := withOperatorStartCandidate(boundedTestContext(t), missing, func() int64 { return 1012 },
+		err := withOperatorStartCandidate(boundedTestContext(t), missing, func() int64 { return 1012 }, store.OpenCredentialsReadOnly,
 			func(string) (*os.File, error) {
 				t.Fatal("canonical path failure reached the lock")
 				return nil, nil
@@ -217,7 +217,7 @@ func TestOperatorStartCandidateRefusals(t *testing.T) {
 	t.Run("missing_store", func(t *testing.T) {
 		root := t.TempDir()
 		before := provisionFiles(t, root)
-		err := withOperatorStartCandidate(boundedTestContext(t), filepath.Join(root, "credentials.db"), func() int64 { return 1012 },
+		err := withOperatorStartCandidate(boundedTestContext(t), filepath.Join(root, "credentials.db"), func() int64 { return 1012 }, store.OpenCredentialsReadOnly,
 			func(string) (*os.File, error) {
 				t.Fatal("missing store reached the lock")
 				return nil, nil
@@ -227,3 +227,32 @@ func TestOperatorStartCandidateRefusals(t *testing.T) {
 		}
 	})
 }
+
+func TestOperatorStartCandidateClosesCredentials(t *testing.T) {
+	creds := filepath.Join(t.TempDir(), "credentials.db")
+	claimFixtureEpisode(t, creds, "ep").Close()
+	credSQL(t, creds, "UPDATE operator_actions SET state='confirmed', owner=NULL")
+	before := serviceFileHash(creds)
+	ctx := boundedTestContext(t)
+	var db *store.Credentials
+	open := func(dsn string) (*store.Credentials, error) {
+		var err error
+		db, err = store.OpenCredentialsReadOnly(dsn)
+		return db, err
+	}
+	entered := false
+	err := withOperatorStartCandidate(ctx, creds, func() int64 { return 1012 }, open, runnerKernelLock,
+		func(_ *os.File, row store.OperatorActionRow) error {
+			entered = row.ID == fixtureActionID
+			return nil
+		})
+	if db != nil {
+		t.Cleanup(func() { db.Close() })
+	}
+	if err != nil || !entered || db == nil || serviceFileHash(creds) != before {
+		t.Fatalf("candidate did not visit the confirmed action read-only: %v", err)
+	}
+	if _, err := db.NextOperatorStart(ctx, 1012); err == nil || err.Error() != "sql: database is closed" {
+		t.Fatalf("candidate did not close its credentials handle: %v", err)
+	}
+}
diff --git a/host/store/operator_runner_test.go b/host/store/operator_runner_test.go
index 6a3aeb8d..6580cfba 100644
--- a/host/store/operator_runner_test.go
+++ b/host/store/operator_runner_test.go
@@ -88,6 +88,35 @@ func TestOperatorCatalogueAdmission(t *testing.T) {
 	publish(a, 2, 4)
 }
 
+func TestOperatorCatalogueNullRefresh(t *testing.T) {
+	for _, tt := range []struct {
+		name      string
+		catalogue []byte
+	}{
+		{"null", nil},
+		{"blob", []byte("eligible")},
+	} {
+		t.Run(tt.name, func(t *testing.T) {
+			c := openTestCredentials(t)
+			ctx := testCtx(t)
+			h := OperatorHeartbeat{Owner: OperatorOwner{Nonce: "fixture"}, Protocol: OperatorRunnerProtocol,
+				Build: "fixture", Readiness: "start", SeenAt: 100, Catalogue: tt.catalogue}
+			gen, err := c.PutOperatorHeartbeat(ctx, h)
+			if err != nil || gen != 1 {
+				t.Fatalf("insert heartbeat: generation %d: %v", gen, err)
+			}
+			h.Owner.Generation, h.SeenAt, h.CatalogueGeneration = gen, 101, 1
+			if gen, err := c.PutOperatorHeartbeat(ctx, h); err != nil || gen != h.Owner.Generation {
+				t.Fatalf("refresh heartbeat: generation %d: %v", gen, err)
+			}
+			got, err := c.OperatorHeartbeat(ctx)
+			if err != nil || !reflect.DeepEqual(got, h) {
+				t.Fatalf("unchanged catalogue must preserve generation across a heartbeat: got %+v, want %+v: %v", got, h, err)
+			}
+		})
+	}
+}
+
 func TestOperatorFinishReadiness(t *testing.T) {
 	c := openTestCredentials(t)
 	cas, runner := operatorClaimFixture(t, c)
````

### Patch M2

````diff M2
diff --git a/cmd/ailang-worldd/operator_episode_lock.go b/cmd/ailang-worldd/operator_episode_lock.go
index 2ffbba1c..f03ec36d 100644
--- a/cmd/ailang-worldd/operator_episode_lock.go
+++ b/cmd/ailang-worldd/operator_episode_lock.go
@@ -12,7 +12,7 @@ import (
 )
 
 func operatorEpisodeLockPath(canonical, episode string) string {
-	return canonical + ".operator-episode-" + episode + ".lock"
+	return store.OperatorEpisodeLockPath(canonical, episode)
 }
 
 func terminalEpisodeLock(creds, episode string, held *os.File) (*os.File, error) {
diff --git a/cmd/ailang-worldd/operator_files.go b/cmd/ailang-worldd/operator_files.go
index 1b4e0f38..3bea962e 100644
--- a/cmd/ailang-worldd/operator_files.go
+++ b/cmd/ailang-worldd/operator_files.go
@@ -259,18 +259,16 @@ func recheckFileApproval(ctx context.Context, credentials string, a approval.Ope
 
 // publishOperatorStartFiles is dormant until the paired-consumer gates open.
 // rawToken is the existing mint output. Neither it nor the MCP bytes enter SQL.
-func publishOperatorStartFiles(ctx context.Context, c runnerConfig, id, rawToken string, clock func() int64, checkpoint func(string) error) error {
+func publishOperatorStartFiles(ctx context.Context, c runnerConfig, held *store.OperatorEpisodeLock, id, rawToken string, clock func() int64, checkpoint func(string) error) error {
 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
 	defer cancel()
 	s, a, err := readRunnerFiles(ctx, c.Credentials, id)
 	if err != nil {
 		return err
 	}
-	lock, err := openRunnerLock(operatorEpisodeLockPath(c.Credentials, a.Episode), 0)
-	if err != nil {
+	if _, err := held.Verify(c.Credentials, a.Episode); err != nil {
 		return err
 	}
-	defer lock.Close()
 	current, a, err := readRunnerFiles(ctx, c.Credentials, id)
 	if err != nil || current != s || s.Phase != "" {
 		return store.ErrOperatorConflict
diff --git a/cmd/ailang-worldd/operator_files_test.go b/cmd/ailang-worldd/operator_files_test.go
index c183f736..8f6eaf6e 100644
--- a/cmd/ailang-worldd/operator_files_test.go
+++ b/cmd/ailang-worldd/operator_files_test.go
@@ -1,6 +1,7 @@
 package main
 
 import (
+	"context"
 	"database/sql"
 	"encoding/json"
 	"errors"
@@ -59,6 +60,13 @@ func mintOperatorFileFixture(t *testing.T, r *worktreeStartRig) string {
 	return token
 }
 
+func publishOperatorFilesFixture(t *testing.T, ctx context.Context, r *worktreeStartRig, token string, clock func() int64, checkpoint func(string) error) error {
+	t.Helper()
+	held := r.borrow(t)
+	defer held.Close()
+	return publishOperatorStartFiles(ctx, r.config, held, r.action.ID, token, clock, checkpoint)
+}
+
 func TestOperatorFilesReconcile(t *testing.T) {
 	t.Run("metadata", TestOperatorFileReadbackMetadata)
 	doc, err := os.ReadFile("../../docs/QUICKSTART.md")
@@ -94,6 +102,15 @@ func TestOperatorFilesReconcile(t *testing.T) {
 			data := [3]string{r.action.Task, token, string(readStartFixtureFile(t, mcpRef))}
 			ctx := presetTestContext(t)
 			initial := observeProvision(t, r.f.r, r.creds)
+			if tc.name == "no_mcp" {
+				err := publishOperatorStartFiles(ctx, r.config, nil, r.action.ID, token, r.src.clock, func(string) error {
+					t.Fatal("unheld publication reached an effect boundary")
+					return nil
+				})
+				if !errors.Is(err, store.ErrOperatorConflict) || !reflect.DeepEqual(initial, observeProvision(t, r.f.r, r.creds)) {
+					t.Fatalf("missing borrow changed files, Git or credentials: %v", err)
+				}
+			}
 			var baseline, want store.OperatorFileState
 			baseline, _, err := readRunnerFiles(ctx, r.creds, r.action.ID)
 			if err != nil {
@@ -147,7 +164,7 @@ func TestOperatorFilesReconcile(t *testing.T) {
 				defer syscall.Umask(mask)
 			}
 			crash := errors.New("injected file interruption")
-			err = publishOperatorStartFiles(ctx, r.config, r.action.ID, token, r.src.clock, func(point string) error {
+			err = publishOperatorFilesFixture(t, ctx, r, token, r.src.clock, func(point string) error {
 				if point == "before ep.task.md" {
 					step("reserved", "running", "", 1)
 					want, _, err = readRunnerFiles(ctx, r.creds, r.action.ID)
@@ -213,7 +230,7 @@ func TestOperatorFilesReconcile(t *testing.T) {
 				}
 				step("observed", "running", "published", 2)
 				exact := observeProvision(t, r.f.r, r.creds)
-				if err := publishOperatorStartFiles(ctx, r.config, r.action.ID, token, r.src.clock, func(string) error { return errors.New("unexpected replay") }); !errors.Is(err, store.ErrOperatorConflict) ||
+				if err := publishOperatorFilesFixture(t, ctx, r, token, r.src.clock, func(string) error { return errors.New("unexpected replay") }); !errors.Is(err, store.ErrOperatorConflict) ||
 					!reflect.DeepEqual(exact, observeProvision(t, r.f.r, r.creds)) {
 					t.Fatal("replay modified a published start")
 				}
@@ -454,7 +471,7 @@ func TestOperatorFilePublicationAdmission(t *testing.T) {
 			case "expiry_after_reservation", "credential_expiry_after_reservation":
 				now := int64(1012)
 				clock = func() int64 { return now }
-				err := publishOperatorStartFiles(ctx, r.config, r.action.ID, token, clock, func(point string) error {
+				err := publishOperatorFilesFixture(t, ctx, r, token, clock, func(point string) error {
 					if point == "before ep.task.md" {
 						now = 1600
 						if name == "credential_expiry_after_reservation" {
@@ -525,7 +542,7 @@ func TestOperatorFilePublicationAdmission(t *testing.T) {
 				want = nil
 			}
 			before := observeProvision(t, r.f.r, r.creds)
-			err := publishOperatorStartFiles(ctx, r.config, r.action.ID, token, clock, func(point string) error {
+			err := publishOperatorFilesFixture(t, ctx, r, token, clock, func(point string) error {
 				if point != "before reservation" ||
 					(name != "expired_credential" && name != "expired_action" && name != "wrong_grants" &&
 						name != "wrong_ttl" && name != "revoked_approver" && !strings.HasPrefix(name, "reserve_") && name != "future_mint") {
@@ -591,7 +608,7 @@ func TestOperatorFileAllowanceOrigin(t *testing.T) {
 		t.Fatalf("allowance provenance missing: %v", err)
 	}
 	before := observeProvision(t, r.f.r, r.creds)
-	err = publishOperatorStartFiles(ctx, r.config, r.action.ID, token, r.src.clock, func(point string) error {
+	err = publishOperatorFilesFixture(t, ctx, r, token, r.src.clock, func(point string) error {
 		t.Errorf("allowance entered publication at %s", point)
 		return nil
 	})
diff --git a/cmd/ailang-worldd/operator_start_candidate.go b/cmd/ailang-worldd/operator_start_candidate.go
index 59682606..b25c9e7f 100644
--- a/cmd/ailang-worldd/operator_start_candidate.go
+++ b/cmd/ailang-worldd/operator_start_candidate.go
@@ -3,7 +3,6 @@ package main
 import (
 	"context"
 	"fmt"
-	"os"
 	"reflect"
 
 	"github.com/sunholo-data/ailang-world/host/store"
@@ -14,8 +13,8 @@ import (
 // visit. The runner may wire a writer here only once the process and recovery
 // gates are met. The visit must use a fresh clock and the locked row for CAS.
 func withOperatorStartCandidate(ctx context.Context, credentials string, now func() int64,
-	open func(string) (*store.Credentials, error), acquire func(string) (*os.File, error),
-	visit func(*os.File, store.OperatorActionRow) error) error {
+	open func(string) (*store.Credentials, error), acquire func(string, string) (*store.OperatorEpisodeLock, error),
+	visit func(*store.OperatorEpisodeLock, store.OperatorActionRow) error) error {
 	canonical, err := store.CanonicalPath(credentials)
 	if err != nil {
 		return err
@@ -32,7 +31,7 @@ func withOperatorStartCandidate(ctx context.Context, credentials string, now fun
 	if !verbs.EpisodeIDPattern.MatchString(initial.Episode) {
 		return store.ErrOperatorConflict
 	}
-	lock, err := acquire(operatorEpisodeLockPath(canonical, initial.Episode))
+	lock, err := acquire(canonical, initial.Episode)
 	if err != nil {
 		return fmt.Errorf("%w: episode lock unavailable: %w", store.ErrOperatorConflict, err)
 	}
diff --git a/cmd/ailang-worldd/operator_start_candidate_test.go b/cmd/ailang-worldd/operator_start_candidate_test.go
index 8b415bf3..6e1fc03a 100644
--- a/cmd/ailang-worldd/operator_start_candidate_test.go
+++ b/cmd/ailang-worldd/operator_start_candidate_test.go
@@ -89,7 +89,7 @@ func TestOperatorStartCandidateProcess(t *testing.T) {
 		}
 	}
 	now := func() int64 { return 1012 }
-	err = withOperatorStartCandidate(ctx, aliasPath, now, store.OpenCredentialsReadOnly, runnerKernelLock, func(*os.File, store.OperatorActionRow) error {
+	err = withOperatorStartCandidate(ctx, aliasPath, now, store.OpenCredentialsReadOnly, store.OpenOperatorEpisodeLock, func(*store.OperatorEpisodeLock, store.OperatorActionRow) error {
 		t.Fatal("visit ran while the real terminal process held the lock")
 		return nil
 	})
@@ -107,15 +107,19 @@ func TestOperatorStartCandidateProcess(t *testing.T) {
 		t.Fatalf("terminal release: %v", err)
 	}
 	handled := errors.New("test visit stopped")
-	var acquired *os.File
-	err = withOperatorStartCandidate(ctx, aliasPath, now, store.OpenCredentialsReadOnly, func(path string) (*os.File, error) {
+	var acquired *store.OperatorEpisodeLock
+	err = withOperatorStartCandidate(ctx, aliasPath, now, store.OpenCredentialsReadOnly, func(path, episode string) (*store.OperatorEpisodeLock, error) {
 		var err error
-		acquired, err = runnerKernelLock(path)
+		acquired, err = store.OpenOperatorEpisodeLock(path, episode)
 		return acquired, err
-	}, func(fd *os.File, row store.OperatorActionRow) error {
+	}, func(lock *store.OperatorEpisodeLock, row store.OperatorActionRow) error {
+		fd, err := lock.Verify(canonical, "ep")
+		if err != nil {
+			t.Fatal(err)
+		}
 		physical, err := filepath.EvalSymlinks(fd.Name())
 		want, wantErr := filepath.EvalSymlinks(operatorEpisodeLockPath(canonical, "ep"))
-		if err != nil || wantErr != nil || physical != want || fd != acquired || row.ID != fixtureActionID {
+		if err != nil || wantErr != nil || physical != want || lock != acquired || row.ID != fixtureActionID {
 			t.Fatalf("visit did not receive exact locked candidate: %v %v", err, wantErr)
 		}
 		if err := runnerWrite(canonical, func(*store.Credentials) error { return nil }); err != nil {
@@ -179,7 +183,7 @@ func TestOperatorStartCandidateRefusals(t *testing.T) {
 			now := int64(1012)
 			var acquired bool
 			err = withOperatorStartCandidate(boundedTestContext(t), creds, func() int64 { return now }, store.OpenCredentialsReadOnly,
-				func(path string) (*os.File, error) {
+				func(path, episode string) (*store.OperatorEpisodeLock, error) {
 					acquired = true
 					if tc.advance {
 						now = 1600
@@ -188,8 +192,8 @@ func TestOperatorStartCandidateRefusals(t *testing.T) {
 						credSQL(t, creds, tc.interleave)
 						before, hash = candidateRow(t, creds, 1012), serviceFileHash(creds)
 					}
-					return runnerKernelLock(path)
-				}, func(*os.File, store.OperatorActionRow) error {
+					return store.OpenOperatorEpisodeLock(path, episode)
+				}, func(*store.OperatorEpisodeLock, store.OperatorActionRow) error {
 					t.Fatal("refused candidate reached its visit")
 					return nil
 				})
@@ -206,7 +210,7 @@ func TestOperatorStartCandidateRefusals(t *testing.T) {
 		before := provisionFiles(t, root)
 		missing := filepath.Join(root, "absent", "credentials.db")
 		err := withOperatorStartCandidate(boundedTestContext(t), missing, func() int64 { return 1012 }, store.OpenCredentialsReadOnly,
-			func(string) (*os.File, error) {
+			func(string, string) (*store.OperatorEpisodeLock, error) {
 				t.Fatal("canonical path failure reached the lock")
 				return nil, nil
 			}, nil)
@@ -218,7 +222,7 @@ func TestOperatorStartCandidateRefusals(t *testing.T) {
 		root := t.TempDir()
 		before := provisionFiles(t, root)
 		err := withOperatorStartCandidate(boundedTestContext(t), filepath.Join(root, "credentials.db"), func() int64 { return 1012 }, store.OpenCredentialsReadOnly,
-			func(string) (*os.File, error) {
+			func(string, string) (*store.OperatorEpisodeLock, error) {
 				t.Fatal("missing store reached the lock")
 				return nil, nil
 			}, nil)
@@ -241,8 +245,8 @@ func TestOperatorStartCandidateClosesCredentials(t *testing.T) {
 		return db, err
 	}
 	entered := false
-	err := withOperatorStartCandidate(ctx, creds, func() int64 { return 1012 }, open, runnerKernelLock,
-		func(_ *os.File, row store.OperatorActionRow) error {
+	err := withOperatorStartCandidate(ctx, creds, func() int64 { return 1012 }, open, store.OpenOperatorEpisodeLock,
+		func(_ *store.OperatorEpisodeLock, row store.OperatorActionRow) error {
 			entered = row.ID == fixtureActionID
 			return nil
 		})
diff --git a/cmd/ailang-worldd/operator_worktree.go b/cmd/ailang-worldd/operator_worktree.go
index f579fed6..99603326 100644
--- a/cmd/ailang-worldd/operator_worktree.go
+++ b/cmd/ailang-worldd/operator_worktree.go
@@ -26,9 +26,8 @@ type operatorStartSource struct {
 	clock                               func() int64
 }
 
-// reserveOperatorStartWorktree is an isolated preparation seam. The
-// paired-process consumer must integrate its held episode fd before dispatch.
-func reserveOperatorStartWorktree(ctx context.Context, credentials string, signed approval.OperatorAction,
+// reserveOperatorStartWorktree borrows the caller-owned episode lock.
+func reserveOperatorStartWorktree(ctx context.Context, credentials string, held *store.OperatorEpisodeLock, signed approval.OperatorAction,
 	attempt store.OperatorMintAttempt, src operatorStartSource, checkpoint func(string) error) error {
 	canonical, err := store.CanonicalPath(credentials)
 	if err != nil {
@@ -37,11 +36,9 @@ func reserveOperatorStartWorktree(ctx context.Context, credentials string, signe
 	if signed.Store != canonical || !verbs.EpisodeIDPattern.MatchString(signed.Episode) {
 		return store.ErrOperatorConflict
 	}
-	lock, err := openRunnerLock(operatorEpisodeLockPath(canonical, signed.Episode), 0)
-	if err != nil {
-		return fmt.Errorf("%w: episode lock unavailable", store.ErrOperatorConflict)
+	if _, err := held.Verify(canonical, signed.Episode); err != nil {
+		return err
 	}
-	defer lock.Close()
 	in, err := reconstructOperatorStart(ctx, signed, attempt, src.presets, src.profiles, src.workspace, src.hosts, src.clock)
 	if err != nil {
 		return err
@@ -87,7 +84,14 @@ func reserveOperatorStartWorktree(ctx context.Context, credentials string, signe
 	if err := checkpoint("before add"); err != nil {
 		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
 	}
-	addErr := worktree.Add(ctx, git, target.Repo, target.Path, target.Branch, target.SHA)
+	fd, err := held.Verify(canonical, signed.Episode)
+	if err != nil {
+		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
+	}
+	addErr := worktree.Add(ctx, git, target.Repo, target.Path, target.Branch, target.SHA, fd)
+	if _, err := held.Verify(canonical, signed.Episode); err != nil {
+		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
+	}
 	if err := checkpoint("before readback"); err != nil {
 		return fmt.Errorf("%w; %s", err, reconcileWorktreeCommand(key.Action.ID))
 	}
diff --git a/cmd/ailang-worldd/operator_worktree_guards_test.go b/cmd/ailang-worldd/operator_worktree_guards_test.go
index 4ef6a1cd..c1d96e37 100644
--- a/cmd/ailang-worldd/operator_worktree_guards_test.go
+++ b/cmd/ailang-worldd/operator_worktree_guards_test.go
@@ -1,6 +1,10 @@
 package main
 
 import (
+	"context"
+	"syscall"
+	"time"
+
 	"crypto/sha256"
 	"encoding/json"
 	"errors"
@@ -32,7 +36,7 @@ func TestOperatorWorktreeSurvivors(t *testing.T) {
 				} else {
 					a.Episode = "BAD!"
 				}
-				if err := reserveOperatorStartWorktree(presetTestContext(t), runnerDB(r.creds), a, r.at, r.src, func(string) error {
+				if err := reserveOperatorStartWorktree(presetTestContext(t), runnerDB(r.creds), nil, a, r.at, r.src, func(string) error {
 					t.Fatal("invalid start crossed a worktree boundary")
 					return nil
 				}); err != store.ErrOperatorConflict {
@@ -366,3 +370,125 @@ func TestOperatorWorktreeLockedReread(t *testing.T) {
 		})
 	}
 }
+
+func TestOperatorBorrowedWorktreeGitProcess(t *testing.T) {
+	const lockEnv = "WORLD_TEST_BORROWED_GIT_LOCK"
+	const swapEnv = "WORLD_TEST_BORROWED_GIT_SWAP"
+	if lockPath := os.Getenv(lockEnv); lockPath != "" {
+		childFD := os.NewFile(3, "inherited episode lock")
+		if childFD == nil {
+			t.Fatal("Git child has no inherited descriptor")
+		}
+		defer childFD.Close()
+		info, err := childFD.Stat()
+		path, pathErr := os.Stat(lockPath)
+		if err != nil || pathErr != nil || !os.SameFile(info, path) {
+			t.Fatalf("Git child inherited a different inode: %v %v", err, pathErr)
+		}
+		// Flock on the already held description succeeds, while an independent
+		// open of the same inode must contend even in the same process.
+		if err := syscall.Flock(3, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
+			t.Fatalf("Git child has a different open file description: %v", err)
+		}
+		if f, err := openRunnerLock(lockPath, 0); !errors.Is(err, syscall.EWOULDBLOCK) {
+			if f != nil {
+				f.Close()
+			}
+			t.Fatalf("Git child can bypass the held lock: %v", err)
+		}
+		if err := os.WriteFile(lockPath+".checked", nil, 0600); err != nil {
+			t.Fatal(err)
+		}
+		if os.Getenv(swapEnv) == "1" {
+			swapFixtureEpisodeLock(t, lockPath)
+		}
+		return
+	}
+	for _, arm := range []string{"kept", "before_git", "during_git"} {
+		t.Run(arm, func(t *testing.T) {
+			testOperatorBorrowedWorktreeGit(t, arm, lockEnv, swapEnv)
+		})
+	}
+}
+
+func swapFixtureEpisodeLock(t *testing.T, path string) {
+	t.Helper()
+	if err := os.Rename(path, path+".old"); err != nil {
+		t.Fatal(err)
+	}
+	if err := os.WriteFile(path, nil, 0600); err != nil {
+		t.Fatal(err)
+	}
+}
+
+func testOperatorBorrowedWorktreeGit(t *testing.T, arm, lockEnv, swapEnv string) {
+	t.Helper()
+	r := newWorktreeStartRig(t, false)
+	held := r.borrow(t)
+	defer held.Close()
+	ctx, cancel := context.WithTimeout(presetTestContext(t), 5*time.Second)
+	defer cancel()
+	before := observeProvision(t, r.f.r, r.creds)
+	checkpoint := func(string) error {
+		t.Fatal("unheld worktree reached an effect boundary")
+		return nil
+	}
+	err := reserveOperatorStartWorktree(ctx, r.creds, nil, r.action, r.at, r.src, checkpoint)
+	if !errors.Is(err, store.ErrOperatorConflict) || !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
+		t.Fatalf("missing borrow changed state or returned wrong class: %v", err)
+	}
+	bin, err := os.Executable()
+	if err != nil {
+		t.Fatal(err)
+	}
+	lockPath := operatorEpisodeLockPath(r.creds, r.action.Episode)
+	hook := filepath.Join(r.f.r.f.repo, ".git", "hooks", "post-checkout")
+	t.Setenv(lockEnv, lockPath)
+	if arm == "during_git" {
+		t.Setenv(swapEnv, "1")
+	}
+	t.Setenv("GORACE", strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
+	script := "#!/bin/sh\nexec " + shellQuote(bin) + " -test.run '^TestOperatorBorrowedWorktreeGitProcess$'\n"
+	if err := os.WriteFile(hook, []byte(script), 0700); err != nil {
+		t.Fatal(err)
+	}
+	var reserved provisionSnapshot
+	err = reserveOperatorStartWorktree(ctx, r.creds, held, r.action, r.at, r.src, func(step string) error {
+		if step == "before add" {
+			if arm == "before_git" {
+				swapFixtureEpisodeLock(t, lockPath)
+			}
+			reserved = observeProvision(t, r.f.r, r.creds)
+		}
+		return nil
+	})
+	if arm != "kept" {
+		after := observeProvision(t, r.f.r, r.creds)
+		if !errors.Is(err, store.ErrOperatorConflict) || !strings.Contains(err.Error(), "service runner reconcile-worktree") ||
+			reserved.credentials == "" || reserved.credentials != after.credentials ||
+			!reflect.DeepEqual(reserved.sessions, after.sessions) {
+			t.Fatalf("changed inode advanced a receipt or lacked recovery: %v", err)
+		}
+		if arm == "before_git" && !reflect.DeepEqual(reserved, after) {
+			t.Fatal("replacement before Git changed the branch or worktree")
+		}
+		if arm == "during_git" && (reserved.worktrees == after.worktrees ||
+			presetGit(t, r.f.r.f, "rev-parse", "refs/heads/world/ep") != r.f.base) {
+			t.Fatal("real Git did not complete; reserved worktree evidence required")
+		}
+		return
+	}
+	if err != nil {
+		t.Fatal(err)
+	}
+	if _, err := os.Stat(lockPath + ".checked"); err != nil {
+		t.Fatalf("real Git hook did not verify the descriptor: %v", err)
+	}
+	r.state(t, "running", "observed", "created", 2, true)
+	if f, err := openRunnerLock(lockPath, 0); !errors.Is(err, syscall.EWOULDBLOCK) {
+		if f != nil {
+			f.Close()
+		}
+		t.Fatalf("worktree helper released the caller's lock: %v", err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_worktree_test.go b/cmd/ailang-worldd/operator_worktree_test.go
index 8ab9820e..6693891b 100644
--- a/cmd/ailang-worldd/operator_worktree_test.go
+++ b/cmd/ailang-worldd/operator_worktree_test.go
@@ -97,7 +97,18 @@ func (r *worktreeStartRig) start(t *testing.T, step func(string) error) error {
 	t.Helper()
 	ctx, cancel := context.WithTimeout(presetTestContext(t), 5*time.Second)
 	defer cancel()
-	return reserveOperatorStartWorktree(ctx, runnerDB(r.creds), r.action, r.at, r.src, step)
+	held := r.borrow(t)
+	defer held.Close()
+	return reserveOperatorStartWorktree(ctx, runnerDB(r.creds), held, r.action, r.at, r.src, step)
+}
+
+func (r *worktreeStartRig) borrow(t *testing.T) *store.OperatorEpisodeLock {
+	t.Helper()
+	held, err := store.OpenOperatorEpisodeLock(r.creds, r.action.Episode)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return held
 }
 
 func (r *worktreeStartRig) exec(t *testing.T, q string, args ...any) {
diff --git a/host/store/operator_episode_lock.go b/host/store/operator_episode_lock.go
new file mode 100644
index 00000000..72f34026
--- /dev/null
+++ b/host/store/operator_episode_lock.go
@@ -0,0 +1,73 @@
+package store
+
+import (
+	"fmt"
+	"os"
+	"strings"
+	"syscall"
+)
+
+// OperatorEpisodeLockPath requires a canonical credentials path and a validated
+// episode. Terminal callers share this on-disk key with the runner and GC.
+func OperatorEpisodeLockPath(canonical, episode string) string {
+	return canonical + ".operator-episode-" + episode + ".lock"
+}
+
+// OperatorEpisodeLock attests to a non-waiting flock taken on this exact open
+// description. The owner alone closes it. The sidecar is never unlinked.
+type OperatorEpisodeLock struct{ file *os.File }
+
+// OpenOperatorEpisodeLock takes the shared kernel lock without waiting.
+func OpenOperatorEpisodeLock(credentials, episode string) (*OperatorEpisodeLock, error) {
+	canonical, err := CanonicalPath(credentials)
+	if err != nil {
+		return nil, err
+	}
+	if episode == "" || strings.ContainsAny(episode, "/\\\x00") {
+		return nil, ErrOperatorConflict
+	}
+	path := OperatorEpisodeLockPath(canonical, episode)
+	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0600)
+	if err != nil {
+		return nil, err
+	}
+	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
+		f.Close()
+		return nil, err
+	}
+	lock := &OperatorEpisodeLock{f}
+	if _, err := lock.Verify(canonical, episode); err != nil {
+		lock.Close()
+		return nil, err
+	}
+	return lock, nil
+}
+
+func (l *OperatorEpisodeLock) Close() error { return l.file.Close() }
+
+// Verify lends the exact descriptor for a synchronous effect and, if needed,
+// for exec inheritance. Its borrower must never close or unlock the file.
+func (l *OperatorEpisodeLock) Verify(canonical, episode string) (*os.File, error) {
+	return l.verify(canonical, episode, os.Getuid())
+}
+
+func (l *OperatorEpisodeLock) verify(canonical, episode string, uid int) (*os.File, error) {
+	path := OperatorEpisodeLockPath(canonical, episode)
+	if l == nil || l.file == nil || l.file.Name() != path {
+		return nil, fmt.Errorf("%w: wrong episode lock", ErrOperatorConflict)
+	}
+	opened, err := l.file.Stat()
+	if err != nil {
+		return nil, fmt.Errorf("%w: closed episode lock: %w", ErrOperatorConflict, err)
+	}
+	current, err := os.Lstat(path)
+	if err != nil || !os.SameFile(opened, current) {
+		return nil, fmt.Errorf("%w: episode lock pathname changed", ErrOperatorConflict)
+	}
+	stat := opened.Sys().(*syscall.Stat_t)
+	if !opened.Mode().IsRegular() || opened.Mode().Perm() != 0600 || int(stat.Uid) != uid ||
+		stat.Nlink != 1 {
+		return nil, fmt.Errorf("%w: unsafe episode lock", ErrOperatorConflict)
+	}
+	return l.file, nil
+}
diff --git a/host/store/operator_episode_lock_test.go b/host/store/operator_episode_lock_test.go
new file mode 100644
index 00000000..6b2308ef
--- /dev/null
+++ b/host/store/operator_episode_lock_test.go
@@ -0,0 +1,270 @@
+package store
+
+import (
+	"context"
+	"errors"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"strconv"
+	"strings"
+	"syscall"
+	"testing"
+	"time"
+)
+
+func TestOperatorEpisodeLockIdentity(t *testing.T) {
+	root := t.TempDir()
+	credentials, err := CanonicalPath(filepath.Join(root, "credentials.db"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	path := OperatorEpisodeLockPath(credentials, "ep")
+	alias := filepath.Join(t.TempDir(), "alias")
+	if err := os.Symlink(root, alias); err != nil {
+		t.Fatal(err)
+	}
+	l, err := OpenOperatorEpisodeLock(filepath.Join(alias, "credentials.db"), "ep")
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { l.Close() })
+	fd, err := l.Verify(credentials, "ep")
+	if err != nil || fd != l.file {
+		t.Fatalf("expected exact descriptor: %v", err)
+	}
+	refuse := func(t *testing.T, lease *OperatorEpisodeLock, episode string, reason string) {
+		t.Helper()
+		before, beforeErr := os.Lstat(path)
+		got, err := lease.Verify(credentials, episode)
+		after, afterErr := os.Lstat(path)
+		if got != nil || !errors.Is(err, ErrOperatorConflict) || err.Error() != ErrOperatorConflict.Error()+": "+reason ||
+			beforeErr != nil || afterErr != nil || !os.SameFile(before, after) || before.Mode() != after.Mode() {
+			t.Fatalf("lock refusal changed inode or returned wrong class: %v", err)
+		}
+	}
+	t.Run("nil", func(t *testing.T) { refuse(t, nil, "ep", "wrong episode lock") })
+	t.Run("zero", func(t *testing.T) { refuse(t, &OperatorEpisodeLock{}, "ep", "wrong episode lock") })
+	t.Run("wrong_key", func(t *testing.T) { refuse(t, l, "other", "wrong episode lock") })
+	t.Run("uid", func(t *testing.T) {
+		before, err := l.file.Stat()
+		if err != nil {
+			t.Fatal(err)
+		}
+		got, err := l.verify(credentials, "ep", os.Getuid()+1)
+		after, statErr := l.file.Stat()
+		if got != nil || !errors.Is(err, ErrOperatorConflict) || statErr != nil || !os.SameFile(before, after) {
+			t.Fatalf("wrong owner changed the inode or was admitted: %v %v", err, statErr)
+		}
+	})
+	t.Run("permissions", func(t *testing.T) {
+		if err := os.Chmod(path, 0644); err != nil {
+			t.Fatal(err)
+		}
+		refuse(t, l, "ep", "unsafe episode lock")
+		if err := os.Chmod(path, 0600); err != nil {
+			t.Fatal(err)
+		}
+	})
+	t.Run("hardlink", func(t *testing.T) {
+		alias := filepath.Join(root, "alias")
+		if err := os.Link(path, alias); err != nil {
+			t.Fatal(err)
+		}
+		refuse(t, l, "ep", "unsafe episode lock")
+		if err := os.Remove(alias); err != nil {
+			t.Fatal(err)
+		}
+	})
+	if err := os.Rename(path, path+".held"); err != nil {
+		t.Fatal(err)
+	}
+	t.Run("removed", func(t *testing.T) {
+		got, err := l.Verify(credentials, "ep")
+		_, statErr := os.Lstat(path)
+		if got != nil || !errors.Is(err, ErrOperatorConflict) || !errors.Is(statErr, os.ErrNotExist) {
+			t.Fatalf("missing pathname recreated or admitted: %v %v", err, statErr)
+		}
+	})
+	t.Run("replaced", func(t *testing.T) {
+		if err := os.WriteFile(path, nil, 0600); err != nil {
+			t.Fatal(err)
+		}
+		refuse(t, l, "ep", "episode lock pathname changed")
+	})
+	if err := l.Close(); err != nil {
+		t.Fatal(err)
+	}
+	if got, err := l.Verify(credentials, "ep"); got != nil || !errors.Is(err, ErrOperatorConflict) || !errors.Is(err, os.ErrClosed) {
+		t.Fatalf("closed descriptor admitted: %v", err)
+	}
+	if next, err := OpenOperatorEpisodeLock(credentials, "ep"); err != nil {
+		t.Fatalf("owner did not release the lock: %v", err)
+	} else {
+		next.Close()
+	}
+}
+
+func operatorLockFDCount(t *testing.T, f *os.File) int {
+	t.Helper()
+	var want syscall.Stat_t
+	if err := syscall.Fstat(int(f.Fd()), &want); err != nil {
+		t.Fatal(err)
+	}
+	dir, err := os.Open("/dev/fd")
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer dir.Close()
+	names, err := dir.Readdirnames(-1)
+	if err != nil {
+		t.Fatal(err)
+	}
+	count := 0
+	for _, name := range names {
+		n, err := strconv.Atoi(name)
+		if err != nil {
+			continue
+		}
+		var have syscall.Stat_t
+		if syscall.Fstat(n, &have) == nil && have.Dev == want.Dev && have.Ino == want.Ino {
+			count++
+		}
+	}
+	return count
+}
+
+func TestOperatorEpisodeLockAcquireRefusals(t *testing.T) {
+	root, err := filepath.EvalSymlinks(t.TempDir())
+	if err != nil {
+		t.Fatal(err)
+	}
+	credentials := filepath.Join(root, "credentials.db")
+	for _, ep := range []string{"", "sub/ep", "sub\\ep", "sub\x00ep"} {
+		got, err := OpenOperatorEpisodeLock(credentials, ep)
+		if got != nil || !errors.Is(err, ErrOperatorConflict) {
+			t.Fatalf("unsafe key %q admitted: %v", ep, err)
+		}
+	}
+	entries, err := os.ReadDir(root)
+	if err != nil || len(entries) != 0 {
+		t.Fatalf("rejected keys left effects: %v %v", entries, err)
+	}
+	if got, err := OpenOperatorEpisodeLock(filepath.Join(root, "missing", "credentials.db"), "ep"); got != nil || !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("missing store parent was created or returned wrong class: %v", err)
+	}
+	path := credentials + ".operator-episode-ep.lock" // legacy terminal key
+	target := filepath.Join(root, "target")
+	if err := os.WriteFile(target, []byte("target"), 0600); err != nil {
+		t.Fatal(err)
+	}
+	if err := os.Symlink(target, path); err != nil {
+		t.Fatal(err)
+	}
+	if got, err := OpenOperatorEpisodeLock(credentials, "ep"); got != nil || !errors.Is(err, syscall.ELOOP) {
+		t.Fatalf("expected no-follow refusal: %v", err)
+	}
+	if link, err := os.Readlink(path); err != nil || link != target {
+		t.Fatalf("symlink refusal altered lock name: %q %v", link, err)
+	}
+	if err := os.Remove(path); err != nil {
+		t.Fatal(err)
+	}
+	if err := os.WriteFile(path, []byte("unsafe"), 0644); err != nil {
+		t.Fatal(err)
+	}
+	if got, err := OpenOperatorEpisodeLock(credentials, "ep"); got != nil || !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("unsafe lock mode admitted: %v", err)
+	}
+	f, err := os.Open(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer f.Close()
+	if count := operatorLockFDCount(t, f); count != 1 {
+		t.Fatalf("refused unsafe lock leaked a descriptor: %d", count)
+	}
+	data, err := os.ReadFile(path)
+	if err != nil || string(data) != "unsafe" {
+		t.Fatalf("refusal altered lock bytes: %q %v", data, err)
+	}
+}
+
+func TestOperatorEpisodeLockFIFORefused(t *testing.T) {
+	credentials := filepath.Join(t.TempDir(), "credentials.db")
+	canonical, err := CanonicalPath(credentials)
+	if err != nil {
+		t.Fatal(err)
+	}
+	path := OperatorEpisodeLockPath(canonical, "ep")
+	if err := syscall.Mkfifo(path, 0600); err != nil {
+		t.Fatal(err)
+	}
+	before, err := os.Lstat(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	f, err := os.OpenFile(path, os.O_RDWR|syscall.O_NONBLOCK, 0600)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer f.Close()
+	// Some kernels reject flock on a FIFO before the stat checks. Exercise the
+	// regular-file check with that real FIFO descriptor on every kernel.
+	got, err := (&OperatorEpisodeLock{f}).Verify(canonical, "ep")
+	after, statErr := os.Lstat(path)
+	if got != nil || !errors.Is(err, ErrOperatorConflict) || statErr != nil || !os.SameFile(before, after) {
+		t.Fatalf("FIFO lock changed or was admitted: %v %v", err, statErr)
+	}
+	if acquired, err := OpenOperatorEpisodeLock(canonical, "ep"); acquired != nil ||
+		(!errors.Is(err, syscall.ENOTSUP) && !errors.Is(err, ErrOperatorConflict)) {
+		t.Fatalf("FIFO entered the acquisition path: %v", err)
+	}
+}
+
+func TestOperatorEpisodeLockContendsWithoutWaiting(t *testing.T) {
+	const childKey = "WORLD_TEST_STORE_EPISODE_LOCK"
+	if path := os.Getenv(childKey); path != "" {
+		canonical, err := CanonicalPath(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		witness, err := os.Open(OperatorEpisodeLockPath(canonical, "ep"))
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer witness.Close()
+		if count := operatorLockFDCount(t, witness); count != 1 {
+			t.Fatalf("descriptor census positive control: %d", count)
+		}
+		got, err := OpenOperatorEpisodeLock(path, "ep")
+		if got != nil {
+			got.Close()
+		}
+		if got != nil || !errors.Is(err, syscall.EWOULDBLOCK) {
+			t.Fatalf("unexpected independent acquisition: %v", err)
+		}
+		if count := operatorLockFDCount(t, witness); count != 1 {
+			t.Fatalf("refused contention leaked a descriptor: %d", count)
+		}
+		return
+	}
+	path := filepath.Join(t.TempDir(), "credentials.db")
+	l, err := OpenOperatorEpisodeLock(path, "ep")
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer l.Close()
+	bin, err := os.Executable()
+	if err != nil {
+		t.Fatal(err)
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
+	defer cancel()
+	child := exec.CommandContext(ctx, bin, "-test.run=^TestOperatorEpisodeLockContendsWithoutWaiting$")
+	child.Env = append(os.Environ(), childKey+"="+path, "GORACE="+strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
+	child.WaitDelay = time.Second
+	if out, err := child.CombinedOutput(); err != nil {
+		t.Fatalf("bounded contention failed: %v %s", err, out)
+	}
+}
diff --git a/host/worktree/worktree.go b/host/worktree/worktree.go
index cd79f34e..27650005 100644
--- a/host/worktree/worktree.go
+++ b/host/worktree/worktree.go
@@ -59,8 +59,9 @@ func Env(environ []string) []string {
 }
 
 // Add runs `git -C repo worktree add --detach dir`, or with branch != ""
-// `git -C repo worktree add -b branch dir`.
-func Add(ctx context.Context, git, repo, dir, branch, startPoint string) error {
+// `git -C repo worktree add -b branch dir`. An operator caller lends its
+// held episode descriptor so it also covers Git if the parent is interrupted.
+func Add(ctx context.Context, git, repo, dir, branch, startPoint string, inherited ...*os.File) error {
 	args := []string{"-C", repo, "worktree", "add"}
 	if branch != "" {
 		args = append(args, "-b", branch, dir)
@@ -70,7 +71,7 @@ func Add(ctx context.Context, git, repo, dir, branch, startPoint string) error {
 	if startPoint != "" {
 		args = append(args, startPoint)
 	}
-	return runGit(ctx, git, args, nil)
+	return runGit(ctx, git, args, nil, inherited...)
 }
 
 // AddExisting checks out an existing branch; git refuses a branch held by another worktree.
@@ -132,12 +133,13 @@ func (c *capped) Write(p []byte) (int, error) {
 // runGit is the package's one git launch. With stdout nil, git's stdout and
 // stderr both go to the error message (Add, Remove); otherwise stdout goes to
 // stdout and only stderr to the message (Output).
-func runGit(ctx context.Context, git string, args []string, stdout *capped) error {
+func runGit(ctx context.Context, git string, args []string, stdout *capped, inherited ...*os.File) error {
 	runCtx, cancel := context.WithTimeout(ctx, Timeout)
 	defer cancel()
 	cmd := exec.CommandContext(runCtx, git, args...)
 	cmd.Env = Env(os.Environ())
 	cmd.Stdin = nil
+	cmd.ExtraFiles = inherited
 	cmd.WaitDelay = waitDelay
 	if deadline, ok := runCtx.Deadline(); ok {
 		// Short callers (session metadata) must not spend another two seconds
````

### Patch M3

````diff M3
diff --git a/cmd/ailang-worldd/episode_gc.go b/cmd/ailang-worldd/episode_gc.go
index 8bd6b967..8a2f3a55 100644
--- a/cmd/ailang-worldd/episode_gc.go
+++ b/cmd/ailang-worldd/episode_gc.go
@@ -26,7 +26,7 @@ ailang-worldd serve (2m after boot, then every 6h).
 
   --workspace-root D   episode worktrees root (default ~/.ailang/world-ws)
   --exec-cache-dir D   per-episode exec caches (default ~/.ailang/world/exec-cache)
-  --credentials-db PATH   credentials for this workspace (matching home service or default workspace inferred)
+  --credentials-db PATH   assert the store paired with this workspace (home service or defaults)
   --older-than DUR     minimum age before a clean merged worktree is removed (default 72h)
   --dry-run            report what would be removed without deleting anything
   --json               emit JSON report
@@ -47,11 +47,8 @@ func defaultExecCacheDir() string {
 }
 
 // episodeGCCredentialsPath uses a store known to belong to the scanned
-// workspace. An explicit path is required for other workspace layouts.
+// workspace. An explicit override must attest to the same physical store.
 func episodeGCCredentialsPath(root, explicit string) string {
-	if explicit != "" {
-		return explicit
-	}
 	home, err := os.UserHomeDir()
 	if err != nil {
 		return ""
@@ -71,6 +68,12 @@ func episodeGCCredentialsPath(root, explicit string) string {
 	if errors.Join(absErr, rootErr, workspaceErr) != nil || root != workspace {
 		return ""
 	}
+	if explicit != "" {
+		match, err := store.SameDatabase(credentials, explicit)
+		if err != nil || !match {
+			return ""
+		}
+	}
 	return credentials
 }
 
@@ -105,14 +108,12 @@ func runEpisodeGC(args []string, stdout, stderr io.Writer) int {
 	defer cancel()
 	credentials := episodeGCCredentialsPath(*wsRoot, *creds)
 	rep := daemon.RunEpisodeGC(ctx, daemon.EpisodeGCOptions{
-		WorkspaceRoot: *wsRoot,
-		ExecCacheDir:  *cacheDir,
-		OlderThan:     *olderThan,
-		DryRun:        *dryRun,
-		Now:           time.Now(),
-		HasOpenOperatorStep: func(ctx context.Context, episode string) (bool, error) {
-			return store.EpisodeHasUnresolvedOperatorSteps(ctx, credentials, episode)
-		},
+		WorkspaceRoot:   *wsRoot,
+		ExecCacheDir:    *cacheDir,
+		OlderThan:       *olderThan,
+		DryRun:          *dryRun,
+		Now:             time.Now(),
+		CredentialsPath: credentials,
 	})
 	if *asJSON {
 		enc := json.NewEncoder(stdout)
diff --git a/cmd/ailang-worldd/episode_gc_test.go b/cmd/ailang-worldd/episode_gc_test.go
index 0c5b18b7..92d5c326 100644
--- a/cmd/ailang-worldd/episode_gc_test.go
+++ b/cmd/ailang-worldd/episode_gc_test.go
@@ -18,7 +18,7 @@ import (
 
 func runGit(t *testing.T, dir string, args ...string) {
 	t.Helper()
-	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
+	cmd := exec.CommandContext(boundedTestContext(t), "git", append([]string{"-C", dir}, args...)...)
 	cmd.Env = append(os.Environ(),
 		"GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com",
 		"GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com",
@@ -144,7 +144,8 @@ func TestEpisodeGCFailSafeAndReclaims(t *testing.T) {
 		}
 		return got
 	}
-	head, status := gitState("rev-parse", "HEAD"), gitState("status", "--porcelain=v1")
+	head, status := gitState("rev-parse", "HEAD"), gitState("--no-optional-locks", "status", "--porcelain=v1")
+	refs, registration := gitState("show-ref"), gitState("worktree", "list", "--porcelain")
 	physical := func(path string) string {
 		t.Helper()
 		got, err := filepath.EvalSymlinks(path)
@@ -184,7 +185,9 @@ func TestEpisodeGCFailSafeAndReclaims(t *testing.T) {
 			}
 		}
 		if _, err := os.Stat(cleanWT); err != nil || !found || !bytes.Equal(head, gitState("rev-parse", "HEAD")) ||
-			!bytes.Equal(status, gitState("status", "--porcelain=v1")) {
+			!bytes.Equal(status, gitState("--no-optional-locks", "status", "--porcelain=v1")) ||
+			!bytes.Equal(refs, gitState("show-ref")) || !bytes.Equal(registration, gitState("worktree", "list", "--porcelain")) ||
+			string(readStartFixtureFile(t, filepath.Join(cleanWT, "README.md"))) != "hello\n" {
 			t.Fatalf("want unchanged ep-clean and skip %q: %+v, stat %v", reason, rep, err)
 		}
 	}
@@ -193,12 +196,22 @@ func TestEpisodeGCFailSafeAndReclaims(t *testing.T) {
 	if got := episodeGCCredentialsPath(wsRoot, ""); got != "" {
 		t.Fatalf("different configured workspace selected credentials %q", got)
 	}
-	assertSkip("non-terminal operator step", "--credentials-db", creds)
+	assertSkip("operator step lookup unavailable", "--credentials-db", creds)
 	alias := filepath.Join(base, "workspace-alias")
 	if err := os.Symlink(wsRoot, alias); err != nil {
 		t.Fatal(err)
 	}
 	setRecord(wsRoot)
+	assertSkip("operator step lookup unavailable", "--credentials-db", defaultCreds)
+	held, err := store.OpenOperatorEpisodeLock(creds, "ep-clean")
+	if err != nil {
+		t.Fatal(err)
+	}
+	assertSkip("operator episode lock unavailable", "--credentials-db", creds)
+	if err := held.Close(); err != nil {
+		t.Fatal(err)
+	}
+	assertSkip("non-terminal operator step", "--credentials-db", creds)
 	if got := episodeGCCredentialsPath(alias, ""); got != creds {
 		t.Fatalf("symlink scan root selected %q, want %q", got, creds)
 	}
diff --git a/host/daemon/auto_gc.go b/host/daemon/auto_gc.go
index 6c44bd04..05f94059 100644
--- a/host/daemon/auto_gc.go
+++ b/host/daemon/auto_gc.go
@@ -2,6 +2,7 @@ package daemon
 
 import (
 	"context"
+	"errors"
 	"fmt"
 	"os"
 	"path/filepath"
@@ -10,7 +11,8 @@ import (
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/store"
-	"github.com/sunholo-data/ailang-world/host/toolprobe"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+	"github.com/sunholo-data/ailang-world/host/worktree"
 )
 
 // EpisodeGCItem describes one worktree or cache directory removed or skipped by RunEpisodeGC.
@@ -39,8 +41,11 @@ type EpisodeGCOptions struct {
 	DryRun        bool
 	Now           time.Time
 	IsActive      func(episode string) bool
-	// A missing lookup, a lookup error, or an unresolved effect bars worktree removal.
-	HasOpenOperatorStep func(ctx context.Context, episode string) (bool, error)
+	// CredentialsPath is the store paired with WorkspaceRoot. An unavailable
+	// lock or ledger bars worktree removal.
+	CredentialsPath    string
+	lookupOperatorStep func(context.Context, string, string) (bool, error)
+	removeRegistered   func(context.Context, string, string, string, *os.File) error
 }
 
 const (
@@ -77,14 +82,12 @@ func (d *Daemon) runAutoGCPass(wsRoot, cacheDir string, now time.Time) EpisodeGC
 	ctx, cancel := context.WithTimeout(context.Background(), autoGCPassBudget)
 	defer cancel()
 	rep := RunEpisodeGC(ctx, EpisodeGCOptions{
-		WorkspaceRoot: wsRoot,
-		ExecCacheDir:  cacheDir,
-		OlderThan:     autoGCOlderThan,
-		Now:           now,
-		IsActive:      d.isEpisodeActive,
-		HasOpenOperatorStep: func(ctx context.Context, episode string) (bool, error) {
-			return store.EpisodeHasUnresolvedOperatorSteps(ctx, d.credentialsPath, episode)
-		},
+		WorkspaceRoot:   wsRoot,
+		ExecCacheDir:    cacheDir,
+		OlderThan:       autoGCOlderThan,
+		Now:             now,
+		IsActive:        func(episode string) bool { return d.isEpisodeActive(episode, now.Unix()) },
+		CredentialsPath: d.credentialsPath,
 	})
 	if len(rep.Removed) > 0 && d.errLog != nil {
 		fmt.Fprintf(d.errLog, "ailang-worldd: auto-gc reclaimed %s (%d removed, %d skipped)\n",
@@ -93,8 +96,8 @@ func (d *Daemon) runAutoGCPass(wsRoot, cacheDir string, now time.Time) EpisodeGC
 	return rep
 }
 
-func (d *Daemon) isEpisodeActive(episode string) bool {
-	snap := d.activity.snapshot(time.Now().Unix())
+func (d *Daemon) isEpisodeActive(episode string, now int64) bool {
+	snap := d.activity.snapshot(now)
 	for _, s := range snap.Sessions {
 		if s.Episode == episode {
 			return true
@@ -117,7 +120,6 @@ func RunEpisodeGC(ctx context.Context, opts EpisodeGCOptions) EpisodeGCReport {
 		Skipped:   []EpisodeGCItem{},
 	}
 	removedWorktrees := map[string]bool{}
-	reposToPrune := map[string]bool{}
 
 	// 1. Scan episode worktrees in WorkspaceRoot.
 	if opts.WorkspaceRoot != "" {
@@ -133,56 +135,12 @@ func RunEpisodeGC(ctx context.Context, opts EpisodeGCOptions) EpisodeGCReport {
 					continue
 				}
 				size := gcDirSizeBytes(wtPath)
-				if opts.IsActive != nil && opts.IsActive(name) {
-					rep.Skipped = append(rep.Skipped, EpisodeGCItem{
-						Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: "active session",
-					})
-					continue
-				}
-				age := now.Sub(info.ModTime())
-				if age < opts.OlderThan {
-					rep.Skipped = append(rep.Skipped, EpisodeGCItem{
-						Kind: "worktree", Name: name, Path: wtPath, Bytes: size,
-						Reason: fmt.Sprintf("too recent (%s < %s)", age.Round(time.Minute), opts.OlderThan),
-					})
-					continue
-				}
-				if reason := checkWorktreeUnsafe(ctx, wtPath); reason != "" {
+				if reason := collectEpisodeWorktree(ctx, opts, now, name, wtPath, info); reason != "" {
 					rep.Skipped = append(rep.Skipped, EpisodeGCItem{
 						Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: reason,
 					})
 					continue
 				}
-				commonDir := gitCommonDir(ctx, wtPath)
-				reason := "operator step lookup unavailable"
-				if opts.HasOpenOperatorStep != nil {
-					unresolved, err := opts.HasOpenOperatorStep(ctx, name)
-					switch {
-					case err != nil || ctx.Err() != nil:
-					case unresolved:
-						reason = "non-terminal operator step"
-					default:
-						reason = ""
-					}
-				}
-				if reason != "" {
-					rep.Skipped = append(rep.Skipped, EpisodeGCItem{
-						Kind: "worktree", Name: name, Path: wtPath, Bytes: size, Reason: reason,
-					})
-					continue
-				}
-				if !opts.DryRun {
-					if err := removeWorktree(ctx, wtPath, commonDir); err != nil {
-						rep.Skipped = append(rep.Skipped, EpisodeGCItem{
-							Kind: "worktree", Name: name, Path: wtPath, Bytes: size,
-							Reason: "remove failed: " + err.Error(),
-						})
-						continue
-					}
-					if commonDir != "" {
-						reposToPrune[commonDir] = true
-					}
-				}
 				removedWorktrees[name] = true
 				rep.ReclaimedBytes += size
 				rep.Removed = append(rep.Removed, EpisodeGCItem{
@@ -191,11 +149,6 @@ func RunEpisodeGC(ctx context.Context, opts EpisodeGCOptions) EpisodeGCReport {
 			}
 		}
 	}
-	if !opts.DryRun {
-		for commonDir := range reposToPrune {
-			_, _ = toolprobe.Output(ctx, "git", "", "--git-dir", commonDir, "worktree", "prune")
-		}
-	}
 
 	// 2. Scan exec-cache dir for orphaned startup probe caches (wep-*) and orphaned episode caches.
 	if opts.ExecCacheDir != "" {
@@ -265,41 +218,137 @@ func RunEpisodeGC(ctx context.Context, opts EpisodeGCOptions) EpisodeGCReport {
 	return rep
 }
 
-func checkWorktreeUnsafe(ctx context.Context, wtPath string) string {
-	gitMeta, err := os.Lstat(filepath.Join(wtPath, ".git"))
-	if err != nil || gitMeta.IsDir() {
-		return "not a linked git worktree"
+func collectEpisodeWorktree(ctx context.Context, opts EpisodeGCOptions, now time.Time, name, path string, scanned os.FileInfo) string {
+	if opts.CredentialsPath == "" {
+		return "operator step lookup unavailable"
 	}
-	out, err := toolprobe.Output(ctx, "git", wtPath, "status", "--porcelain")
+	if !verbs.EpisodeIDPattern.MatchString(name) {
+		return "operator episode lock unavailable"
+	}
+	canonical, err := store.CanonicalPath(opts.CredentialsPath)
 	if err != nil {
-		return "git status failed"
+		return "operator episode lock unavailable"
 	}
-	if strings.TrimSpace(string(out)) != "" {
-		return "dirty worktree (uncommitted or untracked files)"
+	lock, err := store.OpenOperatorEpisodeLock(canonical, name)
+	if err != nil {
+		return "operator episode lock unavailable"
 	}
-	for _, ref := range []string{"@{u}", "origin/dev", "origin/main", "dev", "main"} {
-		if _, err := toolprobe.Output(ctx, "git", wtPath, "merge-base", "--is-ancestor", "HEAD", ref); err == nil {
-			return ""
+	defer lock.Close()
+	lookup := opts.lookupOperatorStep
+	if lookup == nil {
+		lookup = store.EpisodeHasUnresolvedOperatorSteps
+	}
+	unresolved, err := lookup(ctx, canonical, name)
+	if err != nil || ctx.Err() != nil {
+		return "operator step lookup unavailable"
+	}
+	if unresolved {
+		return "non-terminal operator step"
+	}
+
+	info, err := os.Lstat(path)
+	if err != nil || !info.IsDir() || !os.SameFile(info, scanned) {
+		return "worktree changed during scan"
+	}
+	if opts.IsActive != nil && opts.IsActive(name) {
+		return "active session"
+	}
+	if age := now.Sub(info.ModTime()); age < opts.OlderThan {
+		return fmt.Sprintf("too recent (%s < %s)", age.Round(time.Minute), opts.OlderThan)
+	}
+	repo, registered, err := registeredGCWorktree(ctx, path)
+	if err != nil {
+		return err.Error()
+	}
+	if reason := checkWorktreeUnsafe(ctx, registered); reason != "" {
+		return reason
+	}
+	fd, err := lock.Verify(canonical, name)
+	if err != nil {
+		return "operator episode lock unavailable"
+	}
+	if !opts.DryRun {
+		remove := opts.removeRegistered
+		if remove == nil {
+			remove = worktree.RemoveRegistered
+		}
+		if err := remove(ctx, "git", repo, registered, fd); err != nil {
+			return "remove failed: " + err.Error()
 		}
 	}
-	return "unmerged commits on HEAD"
+	return ""
 }
 
-func gitCommonDir(ctx context.Context, wtPath string) string {
-	out, err := toolprobe.Output(ctx, "git", wtPath, "rev-parse", "--path-format=absolute", "--git-common-dir")
+func registeredGCWorktree(ctx context.Context, path string) (string, string, error) {
+	none := errors.New("worktree registration unavailable or changed")
+	meta := filepath.Join(path, ".git")
+	info, err := os.Lstat(meta)
+	if err != nil || !info.Mode().IsRegular() {
+		return "", "", errors.New("not a linked git worktree")
+	}
+	physical, err := filepath.EvalSymlinks(meta)
+	if err != nil {
+		return "", "", none
+	}
+	gitDir, err := worktree.Output(ctx, "git", []string{"-C", path, "rev-parse", "--absolute-git-dir"})
+	if err != nil {
+		return "", "", none
+	}
+	back, err := os.ReadFile(filepath.Join(strings.TrimSuffix(gitDir, "\n"), "gitdir"))
 	if err != nil {
-		return ""
+		return "", "", none
 	}
-	return strings.TrimSpace(string(out))
+	backPath, err := filepath.EvalSymlinks(strings.TrimSuffix(string(back), "\n"))
+	if err != nil || backPath != physical {
+		return "", "", none
+	}
+	common, err := worktree.Output(ctx, "git", []string{"-C", path, "rev-parse", "--path-format=absolute", "--git-common-dir"})
+	if err != nil {
+		return "", "", none
+	}
+	repo, err := filepath.EvalSymlinks(strings.TrimSuffix(common, "\n"))
+	if err != nil {
+		return "", "", none
+	}
+	list, err := worktree.Output(ctx, "git", []string{"--git-dir", repo, "worktree", "list", "--porcelain", "-z"})
+	if err != nil {
+		return "", "", none
+	}
+	match := ""
+	for _, record := range strings.Split(strings.TrimSuffix(list, "\x00\x00"), "\x00\x00") {
+		first, _, _ := strings.Cut(record, "\x00")
+		dir, ok := strings.CutPrefix(first, "worktree ")
+		if !ok {
+			return "", "", none
+		}
+		candidate, err := filepath.EvalSymlinks(dir)
+		if err == nil && candidate == filepath.Dir(physical) {
+			if match != "" {
+				return "", "", none
+			}
+			match = dir
+		}
+	}
+	if match == "" {
+		return "", "", none
+	}
+	return repo, match, nil
 }
 
-func removeWorktree(ctx context.Context, wtPath, commonDir string) error {
-	if commonDir != "" {
-		if _, err := toolprobe.Output(ctx, "git", "", "--git-dir", commonDir, "worktree", "remove", wtPath); err == nil {
-			return nil
+func checkWorktreeUnsafe(ctx context.Context, wtPath string) string {
+	out, err := worktree.Output(ctx, "git", []string{"--no-optional-locks", "-C", wtPath, "status", "--porcelain"})
+	if err != nil {
+		return "git status failed"
+	}
+	if strings.TrimSpace(out) != "" {
+		return "dirty worktree (uncommitted or untracked files)"
+	}
+	for _, ref := range []string{"@{u}", "origin/dev", "origin/main", "dev", "main"} {
+		if _, err := worktree.Output(ctx, "git", []string{"-C", wtPath, "merge-base", "--is-ancestor", "HEAD", ref}); err == nil {
+			return ""
 		}
 	}
-	return os.RemoveAll(wtPath)
+	return "unmerged commits on HEAD"
 }
 
 func gcDirSizeBytes(root string) int64 {
diff --git a/host/daemon/auto_gc_s5e4_test.go b/host/daemon/auto_gc_s5e4_test.go
index 0e9f49e1..fadbc641 100644
--- a/host/daemon/auto_gc_s5e4_test.go
+++ b/host/daemon/auto_gc_s5e4_test.go
@@ -7,13 +7,26 @@ import (
 	"os"
 	"os/exec"
 	"path/filepath"
+	"strings"
 	"testing"
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/store"
+	"syscall"
 )
 
 func TestEpisodeGCOperatorSteps(t *testing.T) {
+	const contenderEnv = "WORLD_TEST_GC_CONTENDER"
+	if creds := os.Getenv(contenderEnv); creds != "" {
+		contender, err := store.OpenOperatorEpisodeLock(creds, "ep")
+		if contender != nil {
+			contender.Close()
+		}
+		if !errors.Is(err, syscall.EWOULDBLOCK) {
+			t.Fatalf("activation acquired at the GC removal boundary: %v", err)
+		}
+		return
+	}
 	root := t.TempDir()
 	ws := filepath.Join(root, "workspaces")
 	if err := os.Mkdir(ws, 0700); err != nil {
@@ -33,8 +46,12 @@ func TestEpisodeGCOperatorSteps(t *testing.T) {
 		return out
 	}
 	git(root, "init", "-b", "main")
-	git(root, "commit", "--allow-empty", "-m", "base")
-	git(root, "worktree", "add", "--detach", wt, "main")
+	if err := os.WriteFile(filepath.Join(root, "tracked"), []byte("fixture"), 0600); err != nil {
+		t.Fatal(err)
+	}
+	git(root, "add", "tracked")
+	git(root, "commit", "-m", "base")
+	git(root, "worktree", "add", "-b", "world/ep", wt, "main")
 	info, err := os.Stat(wt)
 	if err != nil {
 		t.Fatal(err)
@@ -48,11 +65,25 @@ func TestEpisodeGCOperatorSteps(t *testing.T) {
 		return got
 	}
 	path := physical(wt)
-	head, status := git(wt, "rev-parse", "HEAD"), git(wt, "status", "--porcelain=v1")
+	head, status := git(wt, "rev-parse", "HEAD"), git(wt, "--no-optional-locks", "status", "--porcelain=v1")
+	refs, registrations := git(root, "show-ref"), git(root, "worktree", "list", "--porcelain")
+	contents := func() []byte {
+		t.Helper()
+		data, err := os.ReadFile(filepath.Join(wt, "tracked"))
+		if err != nil {
+			t.Fatal(err)
+		}
+		return data
+	}
+	originalContents := contents()
 	if len(status) != 0 {
 		t.Fatalf("fixture is dirty: %q", status)
 	}
-	opts := EpisodeGCOptions{WorkspaceRoot: ws, Now: info.ModTime().Add(73 * time.Hour), OlderThan: 72 * time.Hour}
+	creds := filepath.Join(root, "credentials.db")
+	if _, err := store.UpgradeCredentials(creds); err != nil {
+		t.Fatal(err)
+	}
+	opts := EpisodeGCOptions{WorkspaceRoot: ws, CredentialsPath: creds, Now: info.ModTime().Add(73 * time.Hour), OlderThan: 72 * time.Hour}
 	checkSkip := func(t *testing.T, rep EpisodeGCReport, reason string) {
 		t.Helper()
 		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Name != "ep" ||
@@ -60,19 +91,35 @@ func TestEpisodeGCOperatorSteps(t *testing.T) {
 			t.Fatalf("want sole worktree skip %q: %+v", reason, rep)
 		}
 		if _, err := os.Stat(wt); err != nil || !bytes.Equal(head, git(wt, "rev-parse", "HEAD")) ||
-			!bytes.Equal(status, git(wt, "status", "--porcelain=v1")) {
+			!bytes.Equal(status, git(wt, "--no-optional-locks", "status", "--porcelain=v1")) ||
+			!bytes.Equal(refs, git(root, "show-ref")) || !bytes.Equal(registrations, git(root, "worktree", "list", "--porcelain")) ||
+			!bytes.Equal(originalContents, contents()) {
 			t.Fatalf("refusal changed worktree or Git evidence: %v", err)
 		}
 	}
 	t.Run("candidate", func(t *testing.T) {
 		o := opts
 		o.DryRun = true
-		o.HasOpenOperatorStep = func(context.Context, string) (bool, error) { return false, nil }
+		o.lookupOperatorStep = func(context.Context, string, string) (bool, error) { return false, nil }
+		index := filepath.Join(string(bytes.TrimSpace(git(wt, "rev-parse", "--absolute-git-dir"))), "index")
+		before, err := os.ReadFile(index)
+		if err != nil {
+			t.Fatal(err)
+		}
+		tracked := filepath.Join(wt, "tracked")
+		if err := os.Chtimes(tracked, o.Now, o.Now); err != nil {
+			t.Fatal(err)
+		}
 		rep := RunEpisodeGC(boundedTestContext(t), o)
 		if len(rep.Skipped) != 0 || len(rep.Removed) != 1 || rep.Removed[0].Reason != "clean & merged" ||
-			physical(rep.Removed[0].Path) != path {
+			physical(rep.Removed[0].Path) != path || !bytes.Equal(refs, git(root, "show-ref")) ||
+			!bytes.Equal(registrations, git(root, "worktree", "list", "--porcelain")) || !bytes.Equal(originalContents, contents()) {
 			t.Fatalf("successful negative lookup must allow the candidate: %+v", rep)
 		}
+		after, err := os.ReadFile(index)
+		if err != nil || !bytes.Equal(before, after) {
+			t.Fatalf("GC eligibility refreshed the Git index: %v", err)
+		}
 	})
 	for _, name := range []string{"protected", "error", "timeout", "canceled", "unavailable"} {
 		t.Run(name, func(t *testing.T) {
@@ -83,19 +130,21 @@ func TestEpisodeGCOperatorSteps(t *testing.T) {
 			switch name {
 			case "protected":
 				reason = "non-terminal operator step"
-				o.HasOpenOperatorStep = func(_ context.Context, episode string) (bool, error) {
+				o.lookupOperatorStep = func(_ context.Context, _, episode string) (bool, error) {
 					return episode == "ep", nil
 				}
 			case "error":
-				o.HasOpenOperatorStep = func(context.Context, string) (bool, error) {
+				o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
 					return false, errors.New("fixture store error")
 				}
 			case "timeout":
-				o.HasOpenOperatorStep = func(context.Context, string) (bool, error) {
+				o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
 					return false, context.DeadlineExceeded
 				}
+			case "unavailable":
+				o.CredentialsPath = ""
 			case "canceled":
-				o.HasOpenOperatorStep = func(context.Context, string) (bool, error) {
+				o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
 					cancel()
 					return false, nil
 				}
@@ -106,14 +155,251 @@ func TestEpisodeGCOperatorSteps(t *testing.T) {
 			return
 		}
 	}
-	t.Run("daemon_wiring", func(t *testing.T) {
-		creds := filepath.Join(root, "credentials.db")
-		if _, err := store.UpgradeCredentials(creds); err != nil {
+	t.Run("active", func(t *testing.T) {
+		d := &Daemon{credentialsPath: creds}
+		d.activity.begin("ep", opts.Now.Unix()+1, opts.Now.Unix())()
+		checkSkip(t, d.runAutoGCPass(ws, "", opts.Now), "active session")
+	})
+	t.Run("bad_store_path", func(t *testing.T) {
+		o := opts
+		o.CredentialsPath = filepath.Join(root, "missing", "credentials.db")
+		checkSkip(t, RunEpisodeGC(boundedTestContext(t), o), "operator episode lock unavailable")
+	})
+	t.Run("invalid_episode", func(t *testing.T) {
+		before := git(root, "worktree", "list", "--porcelain")
+		reason := collectEpisodeWorktree(boundedTestContext(t), opts, opts.Now, "INVALID", wt, info)
+		if reason != "operator episode lock unavailable" || !bytes.Equal(before, git(root, "worktree", "list", "--porcelain")) ||
+			!bytes.Equal(refs, git(root, "show-ref")) || !bytes.Equal(originalContents, contents()) {
+			t.Fatalf("malformed episode modified Git or returned wrong reason: %q", reason)
+		}
+	})
+	t.Run("lock_changed_under_review", func(t *testing.T) {
+		o := opts
+		canonical, err := store.CanonicalPath(creds)
+		if err != nil {
+			t.Fatal(err)
+		}
+		path := store.OperatorEpisodeLockPath(canonical, "ep")
+		o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
+			if err := os.Rename(path, path+".old"); err != nil {
+				t.Fatal(err)
+			}
+			return false, os.WriteFile(path, nil, 0600)
+		}
+		checkSkip(t, RunEpisodeGC(boundedTestContext(t), o), "operator episode lock unavailable")
+	})
+	t.Run("replaced_during_scan", func(t *testing.T) {
+		o := opts
+		alternative := filepath.Join(root, "moved")
+		o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
+			if err := os.Rename(wt, alternative); err != nil {
+				t.Fatal(err)
+			}
+			return false, os.Mkdir(wt, 0700)
+		}
+		rep := RunEpisodeGC(boundedTestContext(t), o)
+		if len(rep.Skipped) != 1 || rep.Skipped[0].Reason != "worktree changed during scan" ||
+			len(rep.Removed) != 0 || !bytes.Equal(refs, git(root, "show-ref")) {
+			t.Fatalf("reused pathname was removed or returned wrong reason: %+v", rep)
+		}
+		if err := os.Remove(wt); err != nil {
+			t.Fatal(err)
+		}
+		if err := os.Rename(alternative, wt); err != nil {
+			t.Fatal(err)
+		}
+		checkSkip(t, rep, "worktree changed during scan")
+	})
+	t.Run("not_linked", func(t *testing.T) {
+		_, _, err := registeredGCWorktree(boundedTestContext(t), root)
+		if err == nil || err.Error() != "not a linked git worktree" ||
+			!bytes.Equal(refs, git(root, "show-ref")) || !bytes.Equal(originalContents, contents()) {
+			t.Fatalf("primary repo was admitted or changed: %v", err)
+		}
+	})
+	t.Run("activation_first", func(t *testing.T) {
+		alias := filepath.Join(root, "alias.db")
+		if err := os.Symlink(creds, alias); err != nil {
+			t.Fatal(err)
+		}
+		held, err := store.OpenOperatorEpisodeLock(alias, "ep")
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer held.Close()
+		// Acquire the runner description before the daemon even scans.
+		d := &Daemon{credentialsPath: creds}
+		checkSkip(t, d.runAutoGCPass(ws, "", opts.Now), "operator episode lock unavailable")
+		canonical, err := store.CanonicalPath(creds)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if _, err := os.Lstat(store.OperatorEpisodeLockPath(canonical, "ep-wrong")); !errors.Is(err, os.ErrNotExist) {
+			t.Fatalf("contended pass created a different episode lock: %v", err)
+		}
+	})
+	t.Run("recheck_after_ledger", func(t *testing.T) {
+		o := opts
+		o.lookupOperatorStep = func(_ context.Context, _, _ string) (bool, error) {
+			if err := os.Chtimes(wt, o.Now, o.Now); err != nil {
+				t.Fatal(err)
+			}
+			return false, nil
+		}
+		checkSkip(t, RunEpisodeGC(boundedTestContext(t), o), "too recent (0s < 72h0m0s)")
+		if err := os.Chtimes(wt, info.ModTime(), info.ModTime()); err != nil {
+			t.Fatal(err)
+		}
+	})
+	t.Run("registration_rechecked", func(t *testing.T) {
+		o := opts
+		o.DryRun = true
+		decoy := filepath.Join(root, "decoy")
+		git(root, "worktree", "add", "--detach", decoy, "main")
+		defer git(root, "worktree", "remove", decoy)
+		meta := filepath.Join(wt, ".git")
+		original, err := os.ReadFile(meta)
+		if err != nil {
 			t.Fatal(err)
 		}
+		defer func() {
+			if err := os.WriteFile(meta, original, 0600); err != nil {
+				t.Error(err)
+			}
+		}()
+		foreign, err := os.ReadFile(filepath.Join(decoy, ".git"))
+		if err != nil {
+			t.Fatal(err)
+		}
+		registered := git(root, "worktree", "list", "--porcelain")
+		o.lookupOperatorStep = func(context.Context, string, string) (bool, error) {
+			return false, os.WriteFile(meta, foreign, 0600)
+		}
+		rep := RunEpisodeGC(boundedTestContext(t), o)
+		rewrite, err := os.ReadFile(meta)
+		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Reason != "worktree registration unavailable or changed" ||
+			err != nil || !bytes.Equal(rewrite, foreign) || !bytes.Equal(refs, git(root, "show-ref")) ||
+			!bytes.Equal(registered, git(root, "worktree", "list", "--porcelain")) || !bytes.Equal(originalContents, contents()) {
+			t.Fatalf("changed registration was accepted or modified: %+v, %v", rep, err)
+		}
+	})
+	t.Run("unregistered", func(t *testing.T) {
+		meta := filepath.Join(wt, ".git")
+		original, err := os.ReadFile(meta)
+		if err != nil {
+			t.Fatal(err)
+		}
+		admin := strings.TrimSpace(strings.TrimPrefix(string(original), "gitdir: "))
+		outside := filepath.Join(filepath.Dir(filepath.Dir(admin)), "retained", filepath.Base(admin))
+		if err := os.MkdirAll(filepath.Dir(outside), 0700); err != nil {
+			t.Fatal(err)
+		}
+		if err := os.Rename(admin, outside); err != nil {
+			t.Fatal(err)
+		}
+		defer func() {
+			if err := os.Rename(outside, admin); err != nil {
+				t.Error(err)
+			}
+			if err := os.WriteFile(meta, original, 0600); err != nil {
+				t.Error(err)
+			}
+		}()
+		if err := os.WriteFile(meta, []byte("gitdir: "+outside+"\n"), 0600); err != nil {
+			t.Fatal(err)
+		}
+		// The checkout still works, and Git metadata points back to it.
+		// The administrative directory is outside the live registrations.
+		if got := git(wt, "rev-parse", "HEAD"); !bytes.Equal(got, head) {
+			t.Fatalf("unregistered fixture has no original HEAD: %s", got)
+		}
+		o := opts
+		o.DryRun = true
+		rep := RunEpisodeGC(boundedTestContext(t), o)
+		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Reason != "worktree registration unavailable or changed" ||
+			!bytes.Equal(originalContents, contents()) || !bytes.Equal(refs, git(root, "show-ref")) {
+			t.Fatalf("unregistered checkout changed or misclassified: %+v", rep)
+		}
+	})
+	t.Run("duplicate_registration", func(t *testing.T) {
+		decoy := filepath.Join(root, "duplicate")
+		git(root, "worktree", "add", "--detach", decoy, "main")
+		defer git(root, "worktree", "remove", decoy)
+		admin := string(bytes.TrimSpace(git(decoy, "rev-parse", "--absolute-git-dir")))
+		back := filepath.Join(admin, "gitdir")
+		original, err := os.ReadFile(back)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer func() {
+			if err := os.WriteFile(back, original, 0600); err != nil {
+				t.Error(err)
+			}
+		}()
+		if err := os.WriteFile(back, []byte(filepath.Join(wt, ".git")+"\n"), 0600); err != nil {
+			t.Fatal(err)
+		}
+		duplicate := git(root, "worktree", "list", "--porcelain")
+		matching := 0
+		for _, line := range strings.Split(string(duplicate), "\n") {
+			if dir, ok := strings.CutPrefix(line, "worktree "); ok && physical(dir) == path {
+				matching++
+			}
+		}
+		if matching != 2 {
+			t.Fatalf("Git did not report two spellings of the registration: %s", duplicate)
+		}
+		o := opts
+		o.DryRun = true
+		rep := RunEpisodeGC(boundedTestContext(t), o)
+		if len(rep.Removed) != 0 || len(rep.Skipped) != 1 || rep.Skipped[0].Reason != "worktree registration unavailable or changed" ||
+			!bytes.Equal(originalContents, contents()) || !bytes.Equal(refs, git(root, "show-ref")) ||
+			!bytes.Equal(duplicate, git(root, "worktree", "list", "--porcelain")) {
+			t.Fatalf("duplicate registration changed or misclassified: %+v", rep)
+		}
+	})
+	t.Run("gc_first", func(t *testing.T) {
+		o := opts
+		const refusal = "fixture stops before Git removal"
+		calls := 0
+		o.removeRegistered = func(ctx context.Context, git, repo, registered string, fd *os.File) error {
+			calls++
+			// The contender runs at the last boundary before the Git child.
+			info, err := fd.Stat()
+			lockPath, pathErr := store.CanonicalPath(creds)
+			named, statErr := os.Stat(store.OperatorEpisodeLockPath(lockPath, "ep"))
+			if err != nil || pathErr != nil || statErr != nil || !os.SameFile(info, named) {
+				t.Fatalf("GC did not lend its canonical descriptor to Git: %v %v %v", err, pathErr, statErr)
+			}
+			bin, err := os.Executable()
+			if err != nil {
+				t.Fatal(err)
+			}
+			childCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
+			defer cancel()
+			child := exec.CommandContext(childCtx, bin, "-test.run=^TestEpisodeGCOperatorSteps$")
+			child.Env = append(os.Environ(), contenderEnv+"="+creds, "GORACE="+strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
+			child.WaitDelay = time.Second
+			if out, err := child.CombinedOutput(); err != nil || !bytes.Contains(out, []byte("PASS")) {
+				t.Fatalf("independent activation contender: %v %s", err, out)
+			}
+			return errors.New(refusal)
+		}
+		checkSkip(t, RunEpisodeGC(boundedTestContext(t), o), "remove failed: "+refusal)
+		if calls != 1 {
+			t.Fatalf("GC boundary fired %d times", calls)
+		}
+		contender, err := store.OpenOperatorEpisodeLock(creds, "ep")
+		if err != nil {
+			t.Fatalf("GC lock remained held after removal refusal: %v", err)
+		}
+		contender.Close()
+	})
+	t.Run("daemon_wiring", func(t *testing.T) {
 		// Invoke the synchronous daemon pass with a temporary configured store.
 		// The timers, server, launchctl and network are absent from this fixture.
 		d := &Daemon{credentialsPath: creds}
+		d.activity.begin("ep", opts.Now.Unix()-1, opts.Now.Unix()-2)()
 		rep := d.runAutoGCPass(ws, "", opts.Now)
 		if len(rep.Skipped) != 0 || len(rep.Removed) != 1 || filepath.Join(physical(filepath.Dir(rep.Removed[0].Path)), filepath.Base(rep.Removed[0].Path)) != path {
 			t.Fatalf("daemon did not pass its readable credentials store to GC: %+v", rep)
diff --git a/host/worktree/worktree.go b/host/worktree/worktree.go
index 27650005..5117bd1e 100644
--- a/host/worktree/worktree.go
+++ b/host/worktree/worktree.go
@@ -112,6 +112,12 @@ func Remove(ctx context.Context, git, repo, dir string) error {
 	return runGit(ctx, git, []string{"-C", repo, "worktree", "remove", "--force", dir}, nil)
 }
 
+// RemoveRegistered removes only the named linked worktree, without forcing
+// dirty/locked worktrees or pruning other registrations in the repository.
+func RemoveRegistered(ctx context.Context, git, commonDir, dir string, held *os.File) error {
+	return runGit(ctx, git, []string{"--git-dir", commonDir, "worktree", "remove", dir}, nil, held)
+}
+
 // capped keeps the first max bytes written to it and notes any beyond, which
 // it discards so the child never blocks on a full pipe.
 type capped struct {
diff --git a/host/worktree/worktree_test.go b/host/worktree/worktree_test.go
index fe19c9f2..8acb5002 100644
--- a/host/worktree/worktree_test.go
+++ b/host/worktree/worktree_test.go
@@ -7,6 +7,8 @@ import (
 	"path/filepath"
 	"strconv"
 	"strings"
+	"syscall"
+
 	"testing"
 	"time"
 )
@@ -182,3 +184,59 @@ func TestAddFailureIsReported(t *testing.T) {
 		t.Fatalf("err = %v, want a reported git failure", err)
 	}
 }
+
+func TestRemoveRegisteredLendsDescriptor(t *testing.T) {
+	const check = "WORLD_TEST_GC_CHILD_LOCK"
+	if path := os.Getenv(check); path != "" {
+		fd := os.NewFile(3, "GC episode lock")
+		if fd == nil {
+			t.Fatal("GC child lacks its inherited descriptor")
+		}
+		defer fd.Close()
+		got, err := fd.Stat()
+		want, statErr := os.Stat(path)
+		if err != nil || statErr != nil || !os.SameFile(got, want) {
+			t.Fatalf("GC child inherited wrong inode: %v %v", err, statErr)
+		}
+		if err := syscall.Flock(3, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
+			t.Fatalf("GC child inherited a different open description: %v", err)
+		}
+		return
+	}
+	git := gitBin(t)
+	repo := newRepo(t, git)
+	path := filepath.Join(t.TempDir(), "ep")
+	if err := Add(testCtx(t), git, repo, path, "world/ep", "HEAD"); err != nil {
+		t.Fatal(err)
+	}
+	f, err := os.Create(filepath.Join(t.TempDir(), "episode.lock"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer f.Close()
+	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
+		t.Fatal(err)
+	}
+	bin, err := os.Executable()
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Setenv(check, f.Name())
+	t.Setenv("WORLD_TEST_GC_CHILD_EXE", bin)
+	t.Setenv("WORLD_TEST_GC_REAL_GIT", git)
+	t.Setenv("GORACE", strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
+	wrapper := filepath.Join(t.TempDir(), "git-check")
+	body := "#!/bin/sh\n\"$WORLD_TEST_GC_CHILD_EXE\" -test.run '^TestRemoveRegisteredLendsDescriptor$' || exit 1\nexec \"$WORLD_TEST_GC_REAL_GIT\" \"$@\"\n"
+	if err := os.WriteFile(wrapper, []byte(body), 0700); err != nil {
+		t.Fatal(err)
+	}
+	if err := RemoveRegistered(testCtx(t), wrapper, filepath.Join(repo, ".git"), path, f); err != nil {
+		t.Fatal(err)
+	}
+	if _, err := os.Stat(path); !os.IsNotExist(err) {
+		t.Fatalf("registered worktree still present after real Git remove: %v", err)
+	}
+	if _, err := Output(testCtx(t), git, []string{"-C", repo, "show-ref", "--verify", "refs/heads/world/ep"}); err != nil {
+		t.Fatalf("registered removal changed the branch: %v", err)
+	}
+}
````

### Patch M4p

````diff M4p
diff --git a/cmd/ailang-worldd/operator_file_identity_test.go b/cmd/ailang-worldd/operator_file_identity_test.go
new file mode 100644
index 00000000..4b52e4bf
--- /dev/null
+++ b/cmd/ailang-worldd/operator_file_identity_test.go
@@ -0,0 +1,99 @@
+package main
+
+import (
+	"errors"
+	"os"
+	"path/filepath"
+	"reflect"
+	"testing"
+)
+
+func TestOperatorFilePublicationIdentityWindows(t *testing.T) {
+	t.Run("N04", func(t *testing.T) {
+		root, err := filepath.EvalSymlinks(t.TempDir())
+		if err != nil {
+			t.Fatal(err)
+		}
+		dir, moved := filepath.Join(root, "private"), filepath.Join(root, "retained")
+		if err := os.Mkdir(dir, 0700); err != nil {
+			t.Fatal(err)
+		}
+		set := operatorFileSet{Dir: dir, UID: uint64(os.Geteuid())}
+		calls := 0
+		got := readOperatorFileSetWithOpener(set, "", func(s operatorFileSet) (*os.File, operatorFileProof) {
+			calls++
+			if calls == 2 {
+				if err := os.Rename(dir, moved); err != nil {
+					t.Fatal(err)
+				}
+			}
+			return openOperatorFilesDir(s)
+		})
+		items, listErr := os.ReadDir(moved)
+		if calls != 2 || got.Dir.Status != "directory changed during readback" || listErr != nil || len(items) != 0 {
+			t.Fatalf("directory swap result %q, calls %d, changed files %v: %v", got.Dir.Status, calls, items, listErr)
+		}
+	})
+	// Restore the signed mint fixture between the two publication windows.
+	var r *worktreeStartRig
+	var token string
+	var minted []byte
+	for _, name := range []string{"N12", "N15"} {
+		if !t.Run(name, func(child *testing.T) {
+			if r == nil {
+				r = newWorktreeStartRig(t, false)
+				token = mintOperatorFileFixture(t, r)
+				if err := r.db.Close(); err != nil {
+					child.Fatal(err)
+				}
+				minted = readStartFixtureFile(t, r.creds)
+			}
+			t := child
+			ctx := presetTestContext(t)
+			path, point := r.f.r.sessions, "before ep.task.md"
+			if name == "N15" {
+				path, point = r.path(), "before observe"
+			}
+			moved := path + ".retained"
+			var before provisionSnapshot
+			err := publishOperatorFilesFixture(t, ctx, r, token, r.src.clock, func(at string) error {
+				if at != point {
+					return nil
+				}
+				if err := os.Rename(path, moved); err != nil {
+					return err
+				}
+				before = observeProvision(t, r.f.r, r.creds)
+				return nil
+			})
+			if err == nil || err.Error() != "file publication uncertain; "+filesExample {
+				t.Fatalf("checkpoint swap refusal: %v", err)
+			}
+			s, _, err := readRunnerFiles(ctx, r.creds, r.action.ID)
+			if err != nil || s.Phase != "reserved" || s.Step.Version != 2 ||
+				s.Mint.CredentialID != serviceHash([]byte(token))[len("sha256:"):] {
+				t.Fatalf("lost exact uncertainty evidence: %+v %v", s.Step, err)
+			}
+			after := observeProvision(t, r.f.r, r.creds)
+			if before.worktrees == "" || before.worktrees != after.worktrees || before.branches != after.branches ||
+				!reflect.DeepEqual(before.workspace, after.workspace) || !reflect.DeepEqual(before.sessions, after.sessions) {
+				t.Fatal("refusal changed Git or published files")
+			}
+			_, missing := os.Lstat(path)
+			items, err := os.ReadDir(moved)
+			if err != nil || !errors.Is(missing, os.ErrNotExist) || name == "N12" && len(items) != 0 {
+				t.Fatalf("moved subject files changed: %v %v", items, err)
+			}
+			if name == "N12" {
+				if err := os.Rename(moved, path); err != nil {
+					t.Fatal(err)
+				}
+				if err := os.WriteFile(r.creds, minted, 0600); err != nil {
+					t.Fatal(err)
+				}
+			}
+		}) {
+			return
+		}
+	}
+}
diff --git a/cmd/ailang-worldd/operator_files.go b/cmd/ailang-worldd/operator_files.go
index 3bea962e..96903c36 100644
--- a/cmd/ailang-worldd/operator_files.go
+++ b/cmd/ailang-worldd/operator_files.go
@@ -165,7 +165,12 @@ func readOperatorFileAt(dir *os.File, name, want string, uid uint64, open operat
 }
 
 func readOperatorFileSet(set operatorFileSet, addr string) operatorFilesProof {
-	f, dir := openOperatorFilesDir(set)
+	return readOperatorFileSetWithOpener(set, addr, openOperatorFilesDir)
+}
+
+func readOperatorFileSetWithOpener(set operatorFileSet, addr string,
+	open func(operatorFileSet) (*os.File, operatorFileProof)) operatorFilesProof {
+	f, dir := open(set)
 	p := operatorFilesProof{Dir: dir}
 	if f == nil {
 		return p
@@ -187,7 +192,7 @@ func readOperatorFileSet(set operatorFileSet, addr string) operatorFilesProof {
 			p.Files[2].Status = "unattributed file"
 		}
 	}
-	current, reopened := openOperatorFilesDir(set)
+	current, reopened := open(set)
 	if current != nil {
 		current.Close()
 	}
````

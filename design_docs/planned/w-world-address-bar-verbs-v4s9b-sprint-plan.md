# w-world-address-bar-verbs-v4 — S9b sprint plan

## 0. Scope and outcome

**Outcome B.** Row 243, iteration 449, executor `codex:gpt-6.2-sol`. Continue reviewed design r1 (quorum r2 PROCEED 3/3, all seats present) from S9 M1 and M2 at merge `c1869819`. The inherited M3 already consumes 1,400 changed lines before its listed activation gaps and the judge-448 pins. M1 and M2 deliver the proof, fence, journal and no-terminal private work under tests. Appendix P M3 rebases the remainder as an S9c activation draft. Apply only M1 and M2 from this plan.

M1 and M2 expose zero production paths to the private worker and zero new production claims. This is the release cut. Route admission, proof confirmation, readiness, lock ownership, live input reconstruction, the local journal and the worker caller must go live together in S9c. No request can enter a finish queue whose runner lacks its execution and recovery gates.

The prepared hand-off is a typed durable `local-commit-done` receipt naming the verified SHA with empty `v3_id`. S9c exposes that receipt in production when all activation gates pass. S10 owns the linked V3 prompt and proof. S11 owns push, label creation, PR, finish recording, detach and revoke. Noninteractive readiness uses the reviewed E2 label **read**, Git read and permission signal. Tests require zero remote write, label creation, PR, V3 put or finish record. No push occurs in the fixtures. The public CLI still consults F1 first for effect-bearing or invalid invocations, including surface flags and environment markers. The inherited exact no-effect help exemption stays in place. Schema 10 and AILANG source suffice.

## 1. Measurements

### 1.1 Base and provenance

- HEAD `77dca5edc8206a1791743ebb3c64512f051c52f3`, initially clean. `git diff --stat c1869819 77dca5ed -- cmd host docs app` emitted zero lines. S9 M1 `9a04073e` and M2 `c9cabdfd` match the controller's planned stages.
- The iteration-448 `M3.patch` applied to a fresh archive of this exact base with `/usr/bin/patch -p1 -E --dry-run`, no offset, fuzz or reject. Controls in `$S/inherited-dry-run.raw`. Calc uses sum of insertions and deletions, renames on.
- Design references are §2.2 lifecycle, §3.1–3.3 authority and journal, §4.3 steps 1–2, §5.4 test contract, §6–8 placement and cut, §11 E2 and §12 W15–W20, W28, W33, W37. W33 gives read-level readiness, with separate Git authentication. Actual remote write authorization remains an S11 observation. The controller's 448 record and evaluator r1 are the source of the carried pins.
- **Premises corrected.** An empty tracked diff counted a phantom path against the 257-name cap. The commit primitive accepted an unverified author header. A loose content-addressed object alone counted as resume evidence even when there was no matching World log entry. The inherited boot catalogue passed an empty workspace to its finish scan; the appendix now derives it from the installed record. Inherited `e-store-commit-once` removed SQL placeholders and did not test uniqueness. The corrected fixture or test for each is below.

### 1.2 Size and blast radius

| Stage | Incremental | Cumulative `-- cmd host docs app` | Cumulative unfiltered product |
|---|---:|---:|---:|
| M1 | 401 | +376 / -25 = 401 | 401 |
| M2 | 817 | +1192 / -26 = 1218 | 1218 |
| S9c appendix M3 | 1243 | +2383 / -78 = 2461 | 2461 |

Product counts use `git diff -M --no-index --numstat` on fresh archive trees, equivalent to `git diff -M --numstat 77dca5ed HEAD -- cmd host docs app` after the controller commits the snapshots. All product paths fall inside the filtered roots. Including the two planning artifacts would give an unfiltered final M2 total of 7567 lines. No moved file changes the count. The selected product is below the 1,400 cap.

### 1.3 Baseline and budgets

Merge CI run 46129639435 at `c1869819` recorded cmd race 501.807 s, daemon race 526.741 s, broker race 422.229 s and verifygate race 475.711 s. The CI ceiling is 900 s per package. The named added caps are 15 s cmd and 12 s daemon. Local whole-package options carry explicit 45 m (cmd/daemon) or 20 m (store/broker/verifygate/surfacepage). These are evidence timeouts, independent of CI's ceiling.

Use iteration-449 controller pristine raw logs and `.fails` from `$STATE/base` and wait for its exact `baseline.done` marker. The skipped-11 cmd race and plain/race surfacepage need `$S/base-extra`. No fixed list of pre-existing failures is substituted for an actual comparable base set. The known real binary and timing failures remain in their existing fixtures. A dropped test, package not run, truncated package, setup failure or sandbox denial supplies no pass. The parser prints its positive control and reports package-qualified fail names.

### 1.4 Production callers

| Call | Base | M1 | M2 | Appendix M3 expected |
|---|---:|---:|---:|---:|
| RequestOperatorAction | 1 | 1 | 1 | 1 |
| JudgeOperatorAction | 1 | 1 | 1 | 1 |
| store ConfirmOperatorAction | 1 | 1 | 1 | 1 |
| ClaimOperatorAction | 1 | 1 | 1 | 2 |
| runEpisodeFinishWithPrep | 0 | 0 | 0 | 1 |
| private local Git commit | 0 | 0 | 0 | 1 |
| attended worktree commit | 1 | 1 | 1 | 1 |
| worktree Push / dispatch Push / create label / view PR / create PR / PutVerbRequest | 1 each | 1 each | 1 each | 1 each |

The census prints paths for each match and requires positive recheck, terminal dispatch and installed runner loop probes. In the appendix only `consumeFinish` calls the private entry after a current verified claim. Both terminal and runner use the same episode flock. CLI words, flags and environment have no route into the private entry.

## 2. Milestones

**M1 — Defend proof scope, tool pins and the local Git boundary.** Recheck each signed kind and surface mode using canonical bytes and an enrolled fixture signer. The valid signed start refuses as an unverified finish, and the separate valid signed terminal finish refuses as surface-only. Preserve the finish-specific GH readiness pin exceptions while rejecting a changed Git or World executable. Exercise the capture name/byte/visit bounds and compare the held file with its pathname. Replacing Git configuration with different bytes during parsing refuses before the tool consumes the replacement. The Git primitive verifies input snapshot, exact tree, pinned branch, prior HEAD CAS and the independently expected author before adoption. Schema and production reachability stay at S9.

**M2 — Exercise the private check and commit journal.** Bring forward the inherited local worker and fixtures with proof, fresh effect checks, sandbox-only execution and durable ordered reservations. Refuse a lost or partial check outcome until its matching World outcome can be established. Stop on check mutation, failure or expiry. Reserve the deterministic staged tree before touching the index/ref. Reconcile only exact branch, parent, tree, message/trailer, committer and author. Preserve ambiguous reservations and do no reset. The exact typed local receipt and a zero V3 id are fixture-observed. `ObjectCommits` binds previously observed check evidence to the recorded entry as well as the content address. The store selector tests complete and select a real typed launch receipt without model execution.

**Appendix M3 — S9c activation.** Rebase the remaining inherited prep admission, heartbeat catalogue, claim, readiness, worker, taskless card and redacted projections, with the route-to-worker fixture in scratch. The appendix stays disabled in the executor. Before activation:

1. Run `TestOperatorFinishAdmission` and `TestOperatorFinishRouteWorker` outside the sandbox through the real listener and mux. Require the completed fake child, exact route review and SSHSIG/passkey proof, a claim made by the worker, one sandboxed tool invocation and one local commit. Require the exact typed receipt, no dispatch, no V3 put and no finish record. The route-to-worker fixture drives the pinned workspace-exec RPC through a fake response and observes the exact session, command and args. Run the true sandboxed-check rig alongside this fixture to validate the pinned tool subprocess.
2. Complete the integrated eligibility and drift matrix. The local seams cover HEAD, same-length untracked bytes, path replacement, config replacement, pinned base, protected and unattributed evidence. The appendix exercises degraded health, report unavailability, empty tree, changed token, exact session expiry, HEAD and untracked contents, remote/base and preset/profile drift against the worker. Require the real listener run and an independent live-agent predicate control. Refuse each before any claim or effect and record the operator remedy. Require the Git/config fence at each first use and effect boundary, including a replacement between configuration validation and subprocess use. Pin attribution of already committed paths in the base-to-HEAD delta as well as dirty inputs before enabling approval. The current private snapshot stages only dirty inputs. Every failure needs separate observable refusal and no-change controls.
3. Prove terminal-vs-runner exclusion in both directions under the same lock while the real fake child and worker turns are owned. Inject the timeout in the test and assert non-waiting acquisition plus a bounded read at ≤5 s, with the active action id. The appendix fixture controls a held episode lock before the worker's claim and attempts the same terminal lock during the worker's check. Retain the inherited full CLI F1/exclusion fixture. The listener must execute both directions and surface the active id.
4. Require absence of task, actual artifact, token and manifest markers in routine requests, Needs you, status JSON, CLI diagnostic and synchronized daemon and runner sinks. Observe the rendered signed card and explicit `Request approval` control as a positive disclosure boundary. Run exact taskless app behavior in Node. User-facing walkthrough and status changes are deferred with activation.
5. Run both selected and whole-package gates and mutant ledgers with the outside pristine fail sets. Pass the per-test race ceilings for both cmd and daemon and per-package CI ceiling. O1 remains a pre-existing one-second write budget in `runnerTick`. The new finish pin test uses config/readiness directly, and the completed worker inherits the activation rig's injected runner budget. No selected product changes the heartbeat budget.

## 3. Operator output

S9b verifies the local finish path in fixtures. The operator continues to use the existing typed terminal path until S9c activates the controls. In the appendix the operator enters `episode finish ep --preset p --confirm-on-surface` in the bar, reviews the pinned local commit, and presses `Request approval`. The review and card leave out the task editor and show the exact frozen intent. The key gesture authorizes local preparation. Readiness failure shows `blocked: github-auth` with the attended tool/login/helper remedy. The status panel shows check and commit progress, blocker, uncertainty or retained SHA. The JSON extends status with `checks` and `localCommit`. `ailang-worldd service runner finish-status ACTION-ID` shows the same summary. `TestOperatorFinishAdmission` reads the exact QUICKSTART example from the appendix and submits it to the real listener. S10 issues any PR confirmation later.

## 4. Acceptance, timing and fail sets

The fake witness SSH signer and existing passkey fixture supply authentic test approval bytes. SQLite uses the real migrations in temp credentials and World stores. No real enrollment, model, launchctl, preset, service, Chrome profile, or live developer credential path is used. Observations take hashes/counts of the fixture row/journal, HEAD and index, marker presence, and exact seam calls before and after. Absence is a legal observation. The outer test deadlines bound failed instrumentation, with clocks, budget overrides and fixture channels controlling protocol decisions. Temp ancestry under a 1777 parent is an invocation-only exception to the default TMPDIR. Comparison uses physical paths on both sides.


**Selected M1** from `$S/loc-gates-M1/accept-race.raw`.

| Cmd test | Race s |
|---|---:|
| `TestOperatorFinishToolPins` | 1.08 |
| `TestOperatorFinishProofScope` | 0.39 |
| `TestOperatorFinishCaptureBounds` | 2.62 |
| `TestOperatorFinishReadReplacement` | 0.00 |
| `TestOperatorCommitReservationFence` | 3.58 |
| `TestOperatorFinishConfigReplacement` | 0.72 |
| **Total named cmd** | **8.39** |

`RACE_BUDGET cmd/ailang-worldd n=6/6 added=8.39s cap=15s`

**Selected M2** from `$S/race-M2-final.raw`.

| Cmd test | Race s |
|---|---:|
| `TestOperatorFinishToolPins` | 7.32 |
| `TestOperatorFinishPrep` | 16.36 |
| `TestOperatorCommitReconcile` | 21.32 |
| `TestOperatorCommitFailed` | 8.22 |
| `TestOperatorFinishChecks` | 38.85 |
| `TestOperatorFinishUnsandboxed` | 5.78 |
| `TestOperatorFinishPrepExpiry` | 8.00 |
| `TestOperatorFinishStale` | 37.26 |
| `TestOperatorFinishJournalOnce` | 1.87 |
| `TestOperatorFinishManifestCapture` | 8.14 |
| `TestOperatorFinishWorldEvidence` | 0.05 |
| `TestOperatorFinishProofScope` | 0.08 |
| `TestOperatorFinishCaptureBounds` | 7.70 |
| `TestOperatorFinishReadReplacement` | 0.00 |
| `TestOperatorCommitReservationFence` | 14.94 |
| `TestOperatorFinishConfigReplacement` | 0.87 |
| **Total named cmd** | **176.76** |

`RACE_BUDGET cmd/ailang-worldd n=16/16 added=176.76s cap=15s`

**Selected M2 isolated repeat** from `$S/race-M2-sealed.raw`.

| Cmd test | Race s |
|---|---:|
| `TestOperatorFinishToolPins` | 0.09 |
| `TestOperatorFinishPrep` | 0.41 |
| `TestOperatorCommitReconcile` | 1.02 |
| `TestOperatorCommitFailed` | 0.82 |
| `TestOperatorFinishChecks` | 1.31 |
| `TestOperatorFinishUnsandboxed` | 0.37 |
| `TestOperatorFinishPrepExpiry` | 0.55 |
| `TestOperatorFinishStale` | 0.64 |
| `TestOperatorFinishJournalOnce` | 0.52 |
| `TestOperatorFinishManifestCapture` | 0.45 |
| `TestOperatorFinishWorldEvidence` | 0.05 |
| `TestOperatorFinishProofScope` | 0.10 |
| `TestOperatorFinishCaptureBounds` | 1.02 |
| `TestOperatorFinishReadReplacement` | 0.00 |
| `TestOperatorCommitReservationFence` | 1.74 |
| `TestOperatorFinishConfigReplacement` | 0.56 |
| **Total named cmd** | **9.65** |

`RACE_BUDGET cmd/ailang-worldd n=16/16 added=9.65s cap=15s`

**S9c appendix** from `$S/race-M3-appendix.raw`.

| Cmd test | Race s |
|---|---:|
| `TestOperatorFinishReadiness` | 3.57 |
| `TestOperatorFinishEligibility` | 1.03 |
| **Total named cmd** | **4.60** |

`RACE_BUDGET cmd/ailang-worldd n=2/2 added=4.60s cap=15s`

The first M2 benchmark measured 24.25 s. Two later contended attempts measured 176.76 s and 142.00 s with all sixteen test bodies passing. The isolated selected M2 repeat passed at 9.65 s, with 0.16 s for TestOperatorFinishCompletedLaunch. The pass criterion is the actual ≤15 s cmd sum and ≤12 s daemon sum on the outside acceptance leg. S9b introduces zero daemon tests and no daemon race delta from new test bodies. The S9c listener tests `TestOperatorFinishAdmission` and `TestOperatorFinishRouteWorker` require `outside.sh … M3` for their race seconds. The 12.26 s and 34.20 s in-memory worker diagnostics are plain. Neither is a listener or race acceptance measurement. The S9c executable app case TestOperatorFinishBarReview took 0.47 s under race.

| Touched package | Pristine race | Local stage race and delta |
|---|---:|---|
| host/store | 43.412 s | 171.832 s, +128.420 s |
| cmd/ailang-worldd | 501.807 s CI | Whole comparable local skip-11 base and stage pending outside. Named seconds above. |
| host/daemon (S9c) | 526.741 s CI | Pending outside. |
| host/surfacepage (S9c) | 11.794 s local | 6.747 s, -5.047 s |

## 5. Branch traces, test ownership and mutation matrix

### 5.4 Finish, F1 and remote reconciliation

| Reviewed row | Milestone, executing test and activation requirement |
|---|---|
| Public interactive typed finish | Inherited `TestOperatorFinishTTYParity`. M1 reruns the legacy terminal prompt and typed-y trace. |
| Public no tty, flags or bypass | Inherited `TestOperatorFinishF1`. M1 reruns with no downstream calls or state change. |
| Private missing, forged, swapped, expired, replayed or wrong proof | M1 `TestOperatorFinishNoProof`, `TestOperatorFinishProofScope`. M2 `TestOperatorFinishPrep` checks bound manifest and claim before execution. |
| Valid prep and local commit from completed launch | M2 `TestOperatorCommitReconcile`, `TestOperatorFinishJournalOnce` and `TestOperatorFinishCompletedLaunch` exercise the same receipts independently. S9c `TestOperatorFinishRouteWorker` drives the complete route. |
| No start, degraded/empty/running, session drift/expiry or read failure | S9c `TestOperatorFinishEligibility` has missing start, invalid path and store-read failure. The route worker adds degraded, empty, report-read and changed/expired-session controls. The live process predicate needs its own completed-state adversarial control in §2. |
| Changed preset, HEAD, untracked file, path race, remote/base or profile | M1 `TestOperatorFinishCaptureBounds`, `TestOperatorFinishReadReplacement`, `TestOperatorFinishConfigReplacement`. M2 `TestOperatorFinishManifestCapture` and `TestOperatorFinishStale`. S9c completes actual worker drift matrix. |
| Protected/unattributed or untrusted Git hooks/helpers/config | M1 `TestOperatorFinishWorktreeFence`, `TestOperatorCommitIdentity`; M2 `TestOperatorFinishManifestCapture`. M2 `TestOperatorFinishPrep` also asserts distinct signed protected and unattributed blockers. S9c tests the served and runtime eligibility boundaries. |
| Unsandboxed check including resume | M2 `TestOperatorFinishUnsandboxed`. The same class, unchanged journal and zero child calls on both arms. |
| Missing gh, expired login, SSH/HTTPS Git auth | M1 `TestOperatorFinishAuth`, `TestOperatorFinishToolPins`, outside `TestOperatorFinishAuthInstall`. S9c `TestOperatorFinishReadiness` on confirmed and resumed action. Fake authenticated gh never masks failed Git read. |
| Check fails, aborts, times out, unknown write or mutates input | M2 `TestOperatorFinishChecks`, `TestOperatorFinishWorldEvidence`. No commit or blind rerun, and matched World log evidence. |
| Definite local commit failure | M2 `TestOperatorCommitFailed`. Retained passed check, failed typed receipt, HEAD stable and zero prompt. |
| Lost commit reply, partial stage or different HEAD | M1 `TestOperatorCommitIdentity`, `TestOperatorCommitReservationFence`. M2 `TestOperatorCommitReconcile`, `TestOperatorFinishJournalOnce`. Require branch, parent, tree, message and trailer, committer and author and exact CAS. |
| Prompt size, entropy, store busy or V3 cap | S10 `TestOperatorFinishCardWriteRefusal` |
| V3 store/restart windows | S10 `TestOperatorFinishCardRestart` |
| Missing/superseded/expired/read-failed/revoked V3 approval | S10 `TestOperatorFinishCardRefusal` |
| Typed or unrelated approval reuse | S10 `TestOperatorFinishRejectReuse` |
| V3 proof for different prompt/decision | S10 `TestOperatorFinishOwnBytes` |
| Drift after V3 approval | S10 `TestOperatorFinishPreDispatchFence` |
| Equal/ancestor remote head | S11 `TestOperatorPushIdentity` |
| Diverged/failed/uncertain push | S11 `TestOperatorPushUncertain` |
| Label read/create/readback | S11 `TestOperatorLabelReconcile` |
| Matching existing PR | S11 `TestOperatorReadExistingPR` |
| Mismatching/duplicate/closed/merged PR | S11 `TestOperatorPRIdentity` |
| PR refusal pre-send | S11 `TestOperatorPRNoSend` |
| Lost PR response or empty/stale lookup | S11 `TestOperatorPRResponseLost` |
| Browser close/restart before writeback | S10 Chrome hand-off. S11 `TestOperatorFinishBrowser` verifies one PR and finish record. |
| Linked card with expired/cancelled/replaced action | S10 `TestOperatorFinishLinkedCard` |
| Prep expires after checks | M2 `TestOperatorFinishPrepExpiry`. Retained check receipt, zero commit reservation. S9c still requires explicit re-review/resume wiring. |
| Late exact or forged finish receipt | S11 `TestOperatorFinishObservedRecord` |
| Store crash/head CAS | S11 `TestOperatorFinishRecordCAS` |
| PR with failed writeback, detach or revoke | S11 `TestOperatorFinishCleanup` |

### Security controls and carried scope

| Control | Behavioural pin and affected mutant |
|---|---|
| F1 terminal entry and private capability | `TestOperatorFinishF1` (carried F1-surface), `TestOperatorFinishNoProof` (carried P-recheck), `TestOperatorFinishProofScope` (`e-P-kind` and `e-P-surface`), `TestOperatorFinishPrep` (`J-claim`, `J-manifest`). |
| Readiness before fresh or resumed claim | `TestOperatorFinishToolPins` (`e-pin-skip-git`, `e-pin-skip-all`, `e-R-hash`), S9c `TestOperatorFinishReadiness` (`R-auth`). |
| Stale inputs, bounded Git, no prior unknown check replay | `TestOperatorFinishCaptureBounds` (Judge F bounds), `TestOperatorFinishReadReplacement` (`e-F-samefile2`), `TestOperatorFinishConfigReplacement` (`D-config-swap`), `TestOperatorFinishStale` (`J-bytes`), `TestOperatorFinishUnsandboxed` (`J-unsandboxed`), `TestOperatorFinishChecks` (`J-unknown`, `J-fail-check`, `J-after-check`). |
| Reserve before commit and adopt exact evidence | `TestOperatorCommitReconcile` (`J-early-check`, `J-early-commit`, `J-observe`), `TestOperatorCommitFailed` (`J-no-ref`), `TestOperatorFinishJournalOnce` (binding-preserving `e-store-commit-once`), `TestOperatorCommitIdentity` and `TestOperatorCommitReservationFence` (Judge C pins). |
| Contention | Inherited `TestOperatorTerminalExclusion`. S9c `TestOperatorFinishRouteWorker` (`D-lock-conflict`); symmetric active private runner coverage is an activation gate. |
| Redaction and local-only boundary | M2 fixture cleanup records V3 id, zero daemon requests and zero finish-dispatch/scopes/record counts for the private worker. S9c real listener `TestOperatorFinishAdmission` (`R-prep-card`, `R-status`) and `TestOperatorFinishRouteWorker`. Dedicated projection leak mutants remain activation requirements. |

The former M2 dead surface now has executing callers in fixtures. `TestOperatorFinishManifestCapture` exercises `frozen` and `baseMatches`, `TestOperatorFinishJournalOnce` exercises `finishManifestForRow`, `NextOperatorFinish`, journal projection and schema-10 commit uniqueness. `TestOperatorFinishWorldEvidence` drives the exact stored log membership. `TestOperatorFinishCompletedLaunch` drives `OperatorCompletedLaunches`, `OperatorFinishLaunch`, `Action` and `FinishReview`. The private worker's production caller and route remain deferred together.

Three named Judge carries are behaviourally redundant on the current seam. `e-F-nofilters` removes an explicit defence from `hash-object --stdin` with no pathname. Git already uses raw bytes in that invocation. `e-F-literal` removes a prefix shared by the helpers, none of which supplies a user-controlled pathspec. Stage paths travel as NUL-delimited index-info, and the diff calls use empty `--`. `e-C-samehead` cannot reach a successful adoption because the same exact object would need itself as parent. Keep these explicit controls and re-evaluate them if S9c changes the arguments. All three are zero product delta to retain the guards. Judge-447 C-09/C-10/C-10b and C-01/C-02 stay S10, with the original paired 31- and 24-line test sketches measured in iteration 448.

Other evaluator-only survivors have the following dispositions. Scope sketches below were written in scratch and measured as source lines. They add zero selected product lines and do not certify the corresponding edits.

| Evaluator-only survivors | Disposition |
|---|---|
| e-F1-single-flag, e-F1-deferclose | F1 still refuses no-tty on the later path, and the close change was the intended must-live control. Retain F1 and terminal parity. |
| e-P-id, e-P-digest, e-P-episode | The canonical digest and stored action-byte equality enforce the same three facts. Distinct mismatched action, digest and episode inputs still refuse. |
| e-R-symref, e-R-sha, e-R-firstpath, e-R-sock-mode | Defer the 17-line `$S/carry-size/readiness.go-fragment` sketch. Separate exact read-output, shadowed tool, and AF_UNIX validation against the normal positive E2 harness. The existing E2 fixture still observes separate gh and Git calls. |
| e-F-assume, e-F-gitpart, e-C-lines, f-C-timefields | Defer the shared 19-line `$S/carry-size/git-and-commit.go-fragment` sketch. The primary capture bounds, config, path, committer and author controls are selected. |
| Invalid or duplicate edits in evaluator r1 | Use the compiling f-C variants where named. The old SQL-arity failure never counted. The corrected SQL mutant executes the exact second-intent CAS. |

### Existing tests and helpers modified

| File and existing name | Production function before and after | Behavioural pin |
|---|---|---|
| cmd/ailang-worldd/operator_finish_tree_test.go, `TestOperatorCommitIdentity` | `operatorFinishGit.commit`, `inspectCommit`, commit matching, plus branch and author variants | `f-C-branch`, `e-C-authoremail`, `D-author` |
| Appendix only cmd/ailang-worldd/operator_worktree_test.go, `newWorktreeStartRigWithTask` | Default fake start and prepared receipt still travel through the original start paths. The new overload creates a fixture with an approved HTTPS URL but a bare temp remote. | `TestOperatorFirstManagedEpisodeProcess` retains signed launch binding refusal. `D-lock-conflict` discriminates the new finish caller. |
| Appendix only cmd/ailang-worldd/operator_activation_test.go, `operatorActivationFixtureWithTask` | Existing fake service and runner context, new remote forwarding only on demand | `TestOperatorFirstManagedEpisodeProcess` and `R-auth` |
| Appendix only cmd/ailang-worldd/operator_launch_test.go, `newOperatorLaunchRigWithTask` | Same real fake-child supervisor, task/remote forwarded by new overload | `TestOperatorLaunchDispatchRechecks` plus S9c `D-lock-conflict` |

Other existing test functions are preserved. Additions to `operator_finish_auth_test.go`, `operator_finish_preflight_test.go` and appendix app fixture are new tests. No golden or pinned asset hash changes in M1/M2. S9c app JS changes are guarded by the assembly, pinned copy and executable Node tests. No new production context root or subprocess call site.

### Stem-matched selectors

- **M1** `operator_finish_auth_test.go`: `TestOperatorFinishAuth`, `TestOperatorFinishAuthInstall` (outside), `TestOperatorFinishToolPins`. `operator_finish_preflight{,_test}.go`: `TestOperatorFinishF1`, `TestOperatorFinishTTYParity`, `TestOperatorFinishNoProof`, `TestOperatorFinishProofScope`. `operator_finish_tree{,_test}.go`: `TestOperatorFinishWorktreeFence`, `TestOperatorCommitIdentity`, `TestOperatorFinishCaptureBounds`, `TestOperatorFinishConfigReplacement`, `TestOperatorFinishReadReplacement`, `TestOperatorCommitReservationFence`. Run public flag parity and terminal exclusion.
- **M2** carries the M1 selectors. Add all `operator_finish_local_test.go` tests in `CMD_NEW_M2`, `TestOperatorFinishManifestCapture`, `TestOperatorFinishWorldEvidence`, `TestOperatorFinishCompletedLaunch` and shared launch/step controls in `STORE_FINISH`.
- **Appendix** exact selections in `CMD_SERVICE`, `CMD_STEM_ROUTE`, `DAEMON_STEMS` and `APP_STEMS`. Fixture changes involve `operator_activation_test.go`, `operator_launch_test.go` and `operator_worktree_test.go`. Their full top-level stems must run in the outside leg. Docs run the exact new command in `TestOperatorFinishAdmission`. Run `TestSurfaceVerbOutcomes` for the inherited golden and the pinned app-copy checks for JS.

## 6. Gate and mutation receipts

| Gate | M1 | M2 |
|---|---|---|
| gofmt | OK 0s | OK 0s |
| build | OK 14s | OK 4s |
| vet | OK 44s | OK 2s |
| linux-vet | OK 114s | OK 6s |
| guard-store | OK 42s | OK 2s |
| guard-cli | OK 27s | OK 4s |
| guard-boundary | OK 42s | OK 2s |
| guard-daemon | OK 23s | OK 4s |
| guard-subprocess | OK 2s | OK 2s |
| guard-broker | UNINFORMATIVE UNDER SANDBOX rc=1 | UNINFORMATIVE UNDER SANDBOX rc=1 |
| caller-census | OK 3s | OK 1s |
| accept | rc=2 96s | OK 61s |
| accept-race | OK 199s | rc=1 224s |
| 1777 | OK 19s | OK 1s |
| test-compile | OK 61s | OK 41s |
| verify_ail | OK 307s | OK 85s |

The M2 local package sweep recorded one new plain store failure, TestUncertainBareCommitReconcilesByLogIndexNotHead, from the pre-existing 3 s write context. Its paired named repeat passed 2/2 at base and 2/2 at M2. See the independently recorded full repeat below. No existing store code was changed.

The first local M1 `accept` sweep hit an inherited readiness timeout and a concurrent script rewrite truncated the shell stream. Use the immutable §7 extraction receipt and the outside named leg for release. The first checkpoint M2 rehearsal found a shared fixture DSN that did not yet exist and a contended Git observation. The final fixture canonicalizes the new DSN, injects its budgets, and shares the captured Git baseline for check-only test setup.

The broker subprocess inventory attempted a loopback bind under sandbox and is `UNINFORMATIVE UNDER SANDBOX`. Both appendix listener cases show `bind: operation not permitted` in `$S/appendix-listener.raw`. The in-memory real-mux diagnostic passed the expanded matrix and initial workspace resolution in `$S/route-diagnostic/boot/raw`. The appendix native and Linux vet, build and Node syntax returned rc 0 in `$S/appendix-native/legs.txt`. It is a control of fixture wiring and supplies no listener result.

Whole cmd, daemon, broker and verifygate gate results require the outside scripts and complete controller base fail set. Unfiltered product count was also measured from the full scratch trees. The appendix has no Chrome acceptance claim.

#### Named mutation ledger

| ID | Stage | Named behavioural test | Result and raw-log directory |
|---|---|---|---|
| `e-P-kind` | M1 | `TestOperatorFinishProofScope` | KILLED, `$S/named-mutants-v1` |
| `e-P-surface` | M1 | `TestOperatorFinishProofScope` | KILLED, `$S/named-mutants-v1` |
| `e-pin-skip-git` | M1 | `TestOperatorFinishToolPins` | KILLED, `$S/named-mutants-v1` |
| `e-pin-skip-all` | M1 | `TestOperatorFinishToolPins` | KILLED, `$S/named-mutants-v1` |
| `e-R-hash` | M1 | `TestOperatorFinishToolPins` | KILLED, `$S/named-mutants-v1` |
| `e-F-cfgworktree` | M1 | `TestOperatorFinishCaptureBounds` | KILLED, `$S/named-mutants-v1` |
| `e-F-alternates` | M1 | `TestOperatorFinishCaptureBounds` | KILLED, `$S/named-mutants-v1` |
| `e-F-staged` | M1 | `TestOperatorFinishCaptureBounds` | KILLED, `$S/named-mutants-v1` |
| `e-F-digest` | M1 | `TestOperatorFinishCaptureBounds` | KILLED, `$S/named-mutants-v1` |
| `e-F-visited` | M1 | `TestOperatorFinishCaptureBounds` | KILLED, `$S/named-mutants-v1` |
| `e-F-count` | M1 | `TestOperatorFinishCaptureBounds` | KILLED, `$S/named-mutants-v1` |
| `e-F-budget` | M1 | `TestOperatorFinishCaptureBounds` | KILLED, `$S/named-mutants-v1` |
| `e-F-samefile2` | M1 | `TestOperatorFinishReadReplacement` | KILLED, `$S/named-mutants-v1` |
| `f-C-branch` | M1 | `TestOperatorCommitIdentity` | KILLED, `$S/named-mutants-v1` |
| `f-C-recheck` | M1 | `TestOperatorCommitReservationFence` | KILLED, `$S/named-mutants-v1` |
| `f-C-gottree` | M1 | `TestOperatorCommitReservationFence` | KILLED, `$S/named-mutants-v1` |
| `e-C-cas` | M1 | `TestOperatorCommitReservationFence` | KILLED, `$S/named-mutants-v1` |
| `e-C-authoremail` | M1 | `TestOperatorCommitIdentity` | KILLED, `$S/named-mutants-v1` |
| `e-store-commit-once` | M2 | `TestOperatorFinishJournalOnce` | KILLED, `$S/named-mutants-v2` |
| `D-config-swap` | M1 | `TestOperatorFinishConfigReplacement` | KILLED, `$S/named-mutants-v1` |
| `D-config-readback` | M1 | `TestOperatorFinishConfigReplacement` | KILLED, `$S/named-mutants-v1` |
| `D-name-empty` | M1 | `TestOperatorFinishCaptureBounds` | KILLED, `$S/named-mutants-v1` |
| `D-author` | M1 | `TestOperatorCommitIdentity` | KILLED, `$S/named-mutants-v1` |
| `D-author-time` | M1 | `TestOperatorCommitIdentity` | KILLED, `$S/named-mutants-v1` |
| `J-claim` | M2 | `TestOperatorFinishPrep` | KILLED, `$S/named-mutants-v2` |
| `J-manifest` | M2 | `TestOperatorFinishPrep` | KILLED, `$S/named-mutants-v2` |
| `J-unsandboxed` | M2 | `TestOperatorFinishUnsandboxed` | KILLED, `$S/named-mutants-v2` |
| `J-protected` | M2 | `TestOperatorFinishPrep` | KILLED, `$S/named-mutants-v2` |
| `J-unattributed` | M2 | `TestOperatorFinishPrep` | KILLED, `$S/named-mutants-v2` |
| `J-bytes` | M2 | `TestOperatorFinishStale` | KILLED, `$S/named-mutants-v2` |
| `J-check-version` | M2 | `TestOperatorCommitReconcile` | KILLED, `$S/named-mutants-v2` |
| `J-early-check` | M2 | `TestOperatorCommitReconcile` | KILLED, `$S/named-mutants-v2` |
| `J-unknown` | M2 | `TestOperatorFinishChecks` | KILLED, `$S/named-mutants-v2` |
| `J-fail-check` | M2 | `TestOperatorFinishChecks` | KILLED, `$S/named-mutants-v2` |
| `J-outcome-invalid` | M2 | `TestOperatorFinishChecks` | KILLED, `$S/named-mutants-v2` |
| `J-old-failure` | M2 | `TestOperatorFinishChecks` | KILLED, `$S/named-mutants-v2` |
| `J-after-check` | M2 | `TestOperatorFinishChecks` | KILLED, `$S/named-mutants-v2` |
| `J-early-commit` | M2 | `TestOperatorCommitReconcile` | KILLED, `$S/named-mutants-v2` |
| `J-no-ref` | M2 | `TestOperatorCommitFailed` | KILLED, `$S/named-mutants-v2` |
| `J-observe` | M2 | `TestOperatorCommitReconcile` | KILLED, `$S/named-mutants-v2` |
| `J-prep-expiry` | M2 | `TestOperatorFinishPrepExpiry` | KILLED, `$S/named-mutants-v2` |
| `J-input-first` | M2 | `TestOperatorFinishStale` | KILLED, `$S/named-mutants-v3` |
| `J-observed-commit` | M2 | `TestOperatorCommitReconcile` | KILLED, `$S/named-mutants-v3` |
| `J-tree-identity` | M2 | `TestOperatorCommitReconcile` | KILLED, `$S/named-mutants-v3` |
| `J-check-observe` | M2 | `TestOperatorCommitReconcile` | KILLED, `$S/named-mutants-v3` |
| `J-reserve-check` | M2 | `TestOperatorCommitReconcile` | KILLED, `$S/named-mutants-v3` |
| `J-reserve-commit` | M2 | `TestOperatorCommitReconcile` | KILLED, `$S/named-mutants-v3` |
| `D-outcome-name` | M2 | `TestOperatorFinishWorldEvidence` | KILLED, `$S/named-mutants-v2` |
| `D-outcome-result` | M2 | `TestOperatorFinishWorldEvidence` | KILLED, `$S/named-mutants-v2` |
| `D-world-membership` | M2 | `TestOperatorFinishWorldEvidence` | KILLED, `$S/named-mutants-v2` |
| `D-candidate-expiry` | M2 | `TestOperatorFinishJournalOnce` | KILLED, `$S/named-mutants-v2` |
| `D-observed-candidate` | M2 | `TestOperatorFinishJournalOnce` | KILLED, `$S/named-mutants-v2` |
| `D-launch-version` | M2 | `TestOperatorFinishCompletedLaunch` | KILLED, `$S/named-mutants-v2` |
| `R-auth` | M3 | `TestOperatorFinishReadiness` | KILLED, `$S/named-mutants-v1` |
| `R-prep-card` | M3 | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX, `$S/named-mutants-v1` |
| `R-status` | M3 | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX, `$S/named-mutants-v1` |
| `B-prep-task` | M3 | `TestOperatorFinishBarReview` | KILLED, `$S/named-mutants-v1` |
| `B-prep-kind` | M3 | `TestOperatorFinishBarReview` | KILLED, `$S/named-mutants-v1` |
| `B-prep-flag` | M3 | `TestOperatorFinishBarReview` | KILLED, `$S/named-mutants-v1` |
| `D-lock-conflict` | M3 | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX, `$S/named-mutants-v1` |
| `D-worker-drift` | M3 | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX, `$S/outside.sh` |
| `D-finish-admission` | M3 | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX, `$S/outside.sh` |
| `D-status-redaction` | M3 | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX, `$S/outside.sh` |
| `D-catalogue-workspace` | M3 | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX, `$S/outside.sh` |

Selected named ledger {'KILLED': 53}. Within that ledger the judge-448 carries give {'KILLED': 19}. Appendix named ledger {'KILLED': 4, 'UNINFORMATIVE UNDER SANDBOX': 7}. The three redundant named carries and the other evaluator-only dispositions are recorded in §5.

#### Whole-package mutation ledger

The local whole `host/store` `D-launch-version` control added only `TestOperatorFinishCompletedLaunch` to its pristine package fail set. PENDING — controller to measure the other 24 ids with `whole-mutants.sh LOGDIR M1`, `… M2` and `… M3` outside after `baseline.done`. Use each package's pristine iteration-449 fail set. No whole-package survivor count is inferred from the named runs.

| Prioritized whole-package ID | Result |
|---|---|
| `e-P-kind` | UNMEASURED outside |
| `e-P-surface` | UNMEASURED outside |
| `e-pin-skip-git` | UNMEASURED outside |
| `e-pin-skip-all` | UNMEASURED outside |
| `e-R-hash` | UNMEASURED outside |
| `e-F-digest` | UNMEASURED outside |
| `e-F-visited` | UNMEASURED outside |
| `e-F-count` | UNMEASURED outside |
| `e-F-budget` | UNMEASURED outside |
| `e-F-cfgworktree` | UNMEASURED outside |
| `e-F-alternates` | UNMEASURED outside |
| `e-F-staged` | UNMEASURED outside |
| `e-F-samefile2` | UNMEASURED outside |
| `f-C-branch` | UNMEASURED outside |
| `f-C-recheck` | UNMEASURED outside |
| `e-C-cas` | UNMEASURED outside |
| `f-C-gottree` | UNMEASURED outside |
| `e-C-authoremail` | UNMEASURED outside |
| `e-store-commit-once` | UNMEASURED outside |
| `J-early-commit` | UNMEASURED outside |
| `J-unknown` | UNMEASURED outside |
| `J-prep-expiry` | UNMEASURED outside |
| `D-launch-version` | KILLED, `$S/whole-store-449.raw` |
| `R-prep-card` | UNMEASURED outside |
| `R-status` | UNMEASURED outside |

Controller base observations available at rendering in `$STATE/base/legs.txt`.

````text
LEG store-plain rc=0 parser=0 timeouts=0 fails=0 secs=16
LEG store-race rc=0 parser=0 timeouts=0 fails=0 secs=47
LEG authority-plain rc=0 parser=0 timeouts=0 fails=0 secs=3
LEG authority-race rc=0 parser=0 timeouts=0 fails=0 secs=19
LEG agentrun-plain rc=0 parser=0 timeouts=0 fails=0 secs=16
LEG agentrun-race rc=0 parser=0 timeouts=0 fails=0 secs=20
LEG verifygate-plain rc=1 parser=0 timeouts=0 fails=1 secs=494
LEG verbs-plain rc=0 parser=0 timeouts=0 fails=0 secs=36
LEG broker-plain rc=1 parser=0 timeouts=0 fails=4 secs=629
LEG cmd-plain rc=1 parser=0 timeouts=0 fails=16 secs=1615
LEG daemon-plain rc=1 parser=0 timeouts=0 fails=8 secs=746
````

Store repeats in `$S/store-comparables`.

- Stage plain 24.704 s, empty added fail set.
- Stage race 171.832 s, empty added fail set.

### Reproduction scripts

````bash accept.sh
#!/bin/bash
# TREE M1|M2|M3 [fast|race|outside|outside-race|1777|outside-1777]
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P) || exit 2
T=$(cd "$1" && pwd -P) || exit 2
M=$2 mode=${3:-fast}
. "$S/selectors.sh"
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export GOPROXY=off GIT_OPTIONAL_LOCKS=0
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
L=$(mktemp -d "$S/accept.$M.$mode.XXXXXX") || exit 2
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" || exit 2
cd "$T" || exit 2
printf 'ACCEPT_LOGS %s\n' "$L"
kind=(-p 1)
case "$mode" in race|outside-race) kind=(-race -p 1);; fast|outside|1777|outside-1777) ;; *) exit 2;; esac
case "$M" in
  M1) added="$CMD_NEW_M1";;
  M2) added="$CMD_NEW_M1|$CMD_NEW_M2";;
  M3) added="$CMD_NEW_M3";;
  *) exit 2;;
esac
result=0
named() {
  local tag=$1 pkg=$2 selector=$3 rc=0
  shift 3
  [ "$pkg" != cmd/ailang-worldd ] || set -- "$@" -skip "$CMD_SKIP"
  python3 "$S/bounded.py" 680 go test -v -count=1 -timeout 10m "${kind[@]}" "$@" "./$pkg/" -run "^($selector)$" > "$L/$tag.raw" 2>&1 || rc=$?
  grep -E '^--- (PASS|FAIL|SKIP):|^ok|^FAIL|^panic:|^    [^ ]+\.go:' "$L/$tag.raw" || true
  if grep -Eq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$L/$tag.raw"; then
    echo "LEG $tag INSTRUMENT FAILURE (see $L/$tag.raw)"
    result=2
  elif grep -Eiq '(bind|listen tcp).*(operation not permitted|permission denied)' "$L/$tag.raw"; then
    echo "LEG $tag UNINFORMATIVE UNDER SANDBOX"
    [ "$result" -ne 0 ] || result=77
  elif ! python3 "$S/failset.py" "$L/$tag.raw" "$rc" > "$L/$tag.fails" ||
       ! grep -Eq "^--- PASS: ($selector) " "$L/$tag.raw" || [ "$rc" -ne 0 ]; then
    echo "LEG $tag FAILED rc=$rc"
    result=1
  else
    echo "LEG $tag OK"
  fi
}
timing() {
  awk -v selector="$2" -v cap="$3" -v pkg="$4" '
    BEGIN{nExpected=split(selector, tests, /[|]/)}
    $1=="---" && $2=="PASS:" && $3~("^(" selector ")$") {
      v=$4; gsub(/[()s]/,"",v); s+=v; n++; printf "RACE_TEST %s %s %s\n",pkg,$3,v
    }
    END{
      printf "RACE_BUDGET %s n=%d/%d added=%.2fs cap=%ds\n",pkg,n,nExpected,s,cap
      exit (n!=nExpected || s>cap)
    }' "$1"
}
case "$mode" in
  1777|outside-1777)
    walk='TestOperatorFinishToolPins|TestRunnerPathRefusalClasses'
    if [ "$mode" = outside-1777 ]; then
      walk="$walk|TestOperatorFinishAuthInstall|TestOperatorFinishWorktreeFence|TestOperatorCommitIdentity|TestOperatorFinishConfigReplacement"
    fi
    if [ "$M" = M3 ]; then walk="$walk|TestOperatorFinishReadiness"; fi
    named cmd-1777 cmd/ailang-worldd "$walk"
    if [ "$M" = M3 ] && [ "$mode" = outside-1777 ]; then
      named cmd-worker cmd/ailang-worldd TestOperatorFinishRouteWorker
    else
      echo 'LEG listener-1777 SKIPPED (outside.sh)'
    fi
  ;;
  race|outside-race)
    named cmd-race-named cmd/ailang-worldd "$added"
    if [ "$result" -eq 0 ]; then timing "$L/cmd-race-named.raw" "$added" 15 cmd/ailang-worldd || result=1; fi
    if [ "$M" != M1 ]; then named store-race-named host/store TestOperatorFinishCompletedLaunch; fi
    if [ "$M" = M3 ]; then
      named app-race-named host/surfacepage TestOperatorFinishBarReview
      if [ "$mode" = outside-race ]; then
        named daemon-race-named host/daemon TestOperatorFinishAdmission
        timing "$L/daemon-race-named.raw" TestOperatorFinishAdmission 12 host/daemon || result=1
        named cmd-worker cmd/ailang-worldd TestOperatorFinishRouteWorker
        # Sum the named worker wall time with the selected test total.
        cat "$L/cmd-race-named.raw" "$L/cmd-worker.raw" > "$L/cmd-race-combined.raw"
        timing "$L/cmd-race-combined.raw" "$added|TestOperatorFinishRouteWorker" 15 cmd/ailang-worldd || result=1
      else
        echo 'LEG listeners-race SKIPPED (outside.sh)'
      fi
    fi
  ;;
  fast|outside)
    if [ "$M" = M3 ]; then added="$CMD_NEW_M1|$CMD_NEW_M2|$added"; fi
    named cmd-core cmd/ailang-worldd "$added|$CMD_CORE|TestOperatorFinishNoProof|TestOperatorFinishF1|TestOperatorFinishTTYParity|TestOperatorFinishAuth|TestOperatorFlagParity|TestOperatorTerminalExclusion"
    if [ "$result" -eq 0 ]; then
      while IFS= read -r name; do
        if ! grep -Fq -- "--- PASS: $name (" "$L/cmd-core.raw"; then
          printf 'LEG %s INSTRUMENT FAILURE missing named PASS\n' "$name"
          result=2
        fi
      done < <(printf '%s\n' "$added" | tr '|' '\n')
    fi
    if [ "$M" != M1 ]; then named store-finish-stems host/store "$STORE_FINISH"; fi
    if [ "$mode" = outside ]; then
      named cmd-auth-install cmd/ailang-worldd TestOperatorFinishAuthInstall
    else
      echo 'LEG cmd-auth-install SKIPPED (attended-install fake ran in outside.sh)'
    fi
    if [ "$M" = M3 ]; then
      "$WORLD_EXEC_NODE" --check host/surfacepage/app/app.js || result=1
      named app-stems host/surfacepage "$APP_STEMS"
      if [ "$mode" = outside ]; then
        named cmd-listeners cmd/ailang-worldd "TestOperatorFinishRouteWorker|$CMD_APPENDIX_FIXTURE_STEMS|$CMD_SERVICE"
        named daemon-stems host/daemon "$DAEMON_STEMS"
      else
        echo 'LEG route and service stems SKIPPED (outside.sh)'
      fi
    fi
  ;;
esac
[ "$result" -eq 0 ] || exit "$result"
echo "ACCEPT_OK $M $mode"
````

````bash gates.sh
#!/bin/bash
# TREE LOGDIR M1|M2|M3 [local|outside|outside-base]
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P) || exit 2
T=$(cd "$1" && pwd -P) || exit 2
mkdir "$2" || exit 2
L=$(cd "$2" && pwd -P) || exit 2
M=$3 mode=${4:-local}
BASE="$HOME/.ailang/state/world-iter449/base"
B=77dca5edc8206a1791743ebb3c64512f051c52f3
. "$S/selectors.sh"
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export GOPROXY=off GIT_OPTIONAL_LOCKS=0
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
cd "$T" || exit 2
issue=0 pending=0
note() { printf '%s\n' "$*" | tee -a "$L/legs.txt"; }
skip() { note "LEG $1 SKIPPED ($2)"; }
bound() { python3 "$S/bounded.py" "$@"; }
truncated() { grep -Eq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$1"; }
denied() { grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$1"; }
run() {
  local tag=$1 rc=0 start
  shift
  start=$(date +%s)
  mkdir "$L/$tag.state" "$L/$tag.cache"
  (export AILANG_STATE_DIR="$L/$tag.state" AILANG_CACHE_DIR="$L/$tag.cache"; "$@") > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  if truncated "$L/$tag.raw"; then issue=1; note "LEG $tag INSTRUMENT FAILURE truncated rc=$rc"
  elif denied "$L/$tag.raw" || [ "$rc" -eq 77 ]; then pending=1; note "LEG $tag UNINFORMATIVE UNDER SANDBOX rc=$rc"
  elif [ "$rc" -ne 0 ]; then issue=1; note "LEG $tag rc=$rc $(( $(date +%s)-start ))s"
  else note "LEG $tag OK $(( $(date +%s)-start ))s"; fi
}
record() {
  local tag=$1 pkg=$2 rc=0 parser=0 start cap=1260 limit=20m
  shift 2
  case "$tag" in
    cmd-race) set -- "$@" -skip "$CMD_SKIP"; cap=2800; limit=45m;;
    cmd-plain|daemon-plain|daemon-race) cap=2800; limit=45m;;
  esac
  start=$(date +%s)
  mkdir "$L/$tag.state" "$L/$tag.cache"
  (export AILANG_STATE_DIR="$L/$tag.state" AILANG_CACHE_DIR="$L/$tag.cache"
   bound "$cap" go test -count=1 -timeout "$limit" "$@" "./$pkg/") > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  python3 "$S/failset.py" "$L/$tag.raw" "$rc" > "$L/$tag.fails" 2> "$L/$tag.parser" || parser=$?
  cat "$L/$tag.parser"
  if denied "$L/$tag.raw" && [ "$mode" = local ]; then pending=1; note "LEG $tag UNINFORMATIVE UNDER SANDBOX rc=$rc"
  elif truncated "$L/$tag.raw" || [ "$parser" -ne 0 ]; then issue=1; note "LEG $tag INSTRUMENT FAILURE rc=$rc parser=$parser"
  else cp "$L/$tag.fails" "$L/$tag.verified"; note "LEG $tag recorded rc=$rc $(( $(date +%s)-start ))s"; fi
  grep -E '^(ok|FAIL)[[:space:]]+github.com/' "$L/$tag.raw" >> "$L/pkg-times.txt" || true
}
fmt_touched() {
  local file list="$S/changed-$M"
  [ "$M" = M1 ] || list="$S/changed-$M-cumulative"
  while IFS= read -r file; do
    case "$file" in *.go) [ -z "$(gofmt -l "$file")" ] || { printf 'gofmt %s\n' "$file"; return 1; };; esac
  done < "$list"
}
inventory() {
  local pkg=$1 selector=$2
  bound 390 go test -v -count=1 -timeout 370s "./$pkg/" -run "^($selector)\$" -skip "$CMD_SKIP"
}
mkdir -m 1777 "$L/t1777"
mkdir -m 0700 "$L/t1777/x"
pkgs=(cmd/ailang-worldd)
if [ "$M" != M1 ]; then pkgs+=(host/store); fi
if [ "$M" = M3 ]; then pkgs+=(host/daemon host/surfacepage); fi
tags=(verifygate-plain broker-plain)
owners=(host/verifygate host/broker)
for pkg in "${pkgs[@]}"; do
  for variant in plain race; do
    case "$pkg" in cmd/ailang-worldd) short=cmd;; *) short=${pkg#host/};; esac
    tags+=("$short-$variant") owners+=("$pkg")
  done
done
if [ "$mode" = outside-base ]; then
  for i in "${!tags[@]}"; do
    t=${tags[$i]}
    case "$t" in verifygate-*|broker-*|store-*|cmd-plain|daemon-*) continue;; esac
    if [[ "$t" == *-race ]]; then record "$t" "${owners[$i]}" -race; else record "$t" "${owners[$i]}"; fi
  done
else
  args=()
  for p in "${pkgs[@]}"; do args+=("./$p"); done
  run gofmt fmt_touched
  run build bound 600 go build ./...
  run vet bound 600 go vet "${args[@]}"
  run linux-vet env GOOS=linux python3 "$S/bounded.py" 600 go vet "${args[@]}"
  run guard-store inventory host/store "$STORE_GUARDS"
  run guard-cli inventory cmd/ailang-worldd 'TestCLICredentialWriterSites|TestSessionRootBudgets|TestOperatorProductionDeadlines'
  run guard-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
  run guard-daemon inventory host/daemon 'TestDaemonCredentialWriterSites|TestDaemonDependencyAllowlist|TestDaemonCallerInventoryAliasControl|TestSurfaceVerbIsReadOnly|TestSurfaceVerbWhyIsReadOnly'
  run guard-subprocess inventory host/verifygate TestSubprocessSinkGate
  run guard-broker inventory host/broker TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential
  run caller-census bash "$S/callers.sh" "$S/base-tree" "$T" "$M"
  accept=fast; races=race
  if [ "$mode" = outside ]; then accept=outside; races=outside-race; fi
  run accept bash "$S/accept.sh" "$T" "$M" "$accept"
  run accept-race bash "$S/accept.sh" "$T" "$M" "$races"
  walk=1777; if [ "$mode" = outside ]; then walk=outside-1777; fi
  run 1777 env TMPDIR="$L/t1777/x" bash "$S/accept.sh" "$T" "$M" "$walk"
  run test-compile bound 690 go test -p 1 -run '^$' -timeout 45m ./...
  run verify_ail bound 850 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
  for i in "${!tags[@]}"; do
    t=${tags[$i]}
    if [ "$mode" = local ] && [[ "$t" == cmd-* || "$t" == daemon-* || "$t" == broker-* || "$t" == verifygate-* ]]; then
      skip "$t" 'outside sandbox required'
      continue
    fi
    if [[ "$t" == *-race ]]; then record "$t" "${owners[$i]}" -race; else record "$t" "${owners[$i]}"; fi
  done
  # All recordings precede attribution.
  base_marker="ALL_DONE for $B"
  if ! grep -Fxq "$base_marker" "$BASE/../baseline.done" 2>/dev/null; then
    note "LEG baseline-completion SKIPPED (controller baseline still recording)"
    if [ "$mode" = outside ]; then pending=1; fi
  fi
  # Compare comparable fail sets only.
  for tag in "${tags[@]}"; do
    if [ ! -f "$L/$tag.verified" ]; then skip "compare-$tag" 'stage record unavailable'; continue; fi
    source="$BASE/$tag"
    case "$tag" in cmd-race|surfacepage-*) source="$S/base-extra/$tag";; esac
    if [ ! -s "$source.raw" ] || [ ! -f "$source.rc" ] || truncated "$source.raw" ||
       ! python3 "$S/failset.py" "$source.raw" "$(cat "$source.rc")" > "$L/$tag.base"; then
      pending=1; skip "compare-$tag" 'comparable pristine fail set unavailable'; continue
    fi
    LC_ALL=C sort -u "$L/$tag.base" > "$L/$tag.base.sorted"
    LC_ALL=C sort -u "$L/$tag.verified" > "$L/$tag.head.sorted"
    LC_ALL=C comm -13 "$L/$tag.base.sorted" "$L/$tag.head.sorted" > "$L/$tag.new"
    if [ -s "$L/$tag.new" ]; then issue=1; note "LEG compare-$tag NEW FAIL SET"; else note "LEG compare-$tag OK"; fi
  done
fi
note "GATES_DONE issue=$issue pending=$pending"
printf '%s\n' ALL_DONE > "$L/ALL_DONE"
printf '%s\n' ALL_DONE
[ "$issue" -eq 0 ] || exit 1
[ "$pending" -eq 0 ] || exit 77
````

````bash snapshot.sh
#!/bin/bash
# snapshot.sh TREE M1|M2|M3; SNAP_STATUS_LIST/ SNAP_PARENT are rehearsal-only.
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1|M2|M3) ;; *) exit 2;; esac
B=77dca5edc8206a1791743ebb3c64512f051c52f3
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter449/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s9b-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s9b.json"
  ' | LC_ALL=C sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
fi
test -s "$D/STATUS"
awk -F '\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
while IFS="$(printf '\t')" read -r kind p; do
  [ "$kind" = D ] && continue
  mkdir -p "$D/$(dirname "$p")"
  cp -p "$p" "$D/$p"
done < "$D/STATUS"
(
  cd "$D"
  while IFS="$(printf '\t')" read -r kind p; do
    [ "$kind" = D ] && continue
    shasum -a 256 "$p"
  done < STATUS > MANIFEST.sha256
  shasum -a 256 -c MANIFEST.sha256
)
wc -l "$D"/{STATUS,DELETED,MANIFEST.sha256}
````

````bash selectors.sh
CMD_SKIP='^(TestCLIAgainstSeToolsDaemon|TestCLINoGenesisIsNoWorldHead|TestDemoSeedApprovalPacket|TestDemoSeedEndToEnd|TestDemoSeedManifestReferences|TestDemoSeedProofStage|TestDemoSeedReadTestCheck|TestDemoSeedRejectsToolReportRefusals|TestDemoSeedRosterKindMismatchRefuses|TestDemoSeedWhyNamedLinks|TestMarketingStillRoutes)$'
CMD_S9='TestOperatorFinishAuth|TestOperatorFinishAuthInstall|TestOperatorFinishF1|TestOperatorFinishTTYParity|TestOperatorFinishNoProof'
CMD_CORE='TestOperatorFinishWorktreeFence|TestOperatorCommitIdentity'
CMD_ACTIVATE='TestOperatorFinishReadiness|TestOperatorFinishEligibility|TestOperatorFinishPrep|TestOperatorFinishPrepExpiry|TestOperatorFinishStale|TestOperatorFinishChecks|TestOperatorFinishUnsandboxed|TestOperatorCommitFailed|TestOperatorCommitReconcile'
CMD_FINISH_STEMS='TestEpisodeFinish.*|TestCanonicalFinishArgv|TestFinishConfirmOnSurfaceRefusesBeforeTheStore|TestFinishVerbRequestFields|TestSurfaceCardLinks|TestPutFinishVerbRequestRefusals|TestAwaitVerbConfirmationEnds|TestPollVerbRequestRechecks|TestReadBackFinishDecision|TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision|TestFinishRepository|TestFinishProductionGHArgv|TestFinishPromptIsTheGolden'
CMD_FINISH_OUTSIDE='TestFinishConfirmOnSurface(OpensThePR|RefusesBeforeAnyCheck|RerunsAndTypedPrompt|KeepsTheCommit|DispatchesTheNamedDecision|ResumeReusesTheApproval)'
CMD_SERVICE='TestOperatorService.*|TestRunnerPathRefusalClasses|TestOperatorRunnerKernelLock|TestOperatorStatusDisconnected'
CMD_SCOPED='TestOperatorFlagParity|TestOperatorTerminalExclusion'
CMD_STEM_ROUTE='TestOperatorStartBrowser|TestOperatorLaunchDispatchRechecks|TestOperatorFirstManagedEpisodeProcess'
DAEMON_STEMS='TestOperator(AdmissionHTTP|ReviewRefusesForgedCatalogue|ReplayConflicts|LinksProjection|ConfirmHTTP|ConfirmProjection|ConfirmationAdmissionNeighbours|CardAssets|FinishAdmission)|TestSurfaceVerbOutcomes'
APP_STEMS='TestSurfaceApp.*|TestApprove(Page|Model).*|TestAddressHintsArePinned|TestNeedsYouListsConfirmations|TestNeedsYouConfirmationEmbeddedNewline|TestNeedsYouConfirmationsNonArray|TestVerbDecidePrefixBoundary|TestNeedsYouListsDecisions|TestOperator(BarStartReview|OtherVerbParity|LiteralTask|DuplicateTap|Reconnect|PanelLatest|PanelNetwork|LinkAllowlist|FinishBarReview)'
STORE_GUARDS='TestProductionContextRoots|TestCredentialsMethodsRequireADeadline|TestNoNewDeadlineFreeStoreReads|TestStoreDeadlineGuardDynamic|TestOperatorSchema.*|TestOperatorClaimRace|TestOperatorAdmissionBoundaries'
CMD_PINS='TestOperatorFinish(ProofScope|ToolPins|CaptureBounds|ReadReplacement|ConfigReplacement)|TestOperatorCommitReservationFence'
CMD_NEW_M1='TestOperatorFinishProofScope|TestOperatorFinishToolPins|TestOperatorFinishCaptureBounds|TestOperatorFinishReadReplacement|TestOperatorFinishConfigReplacement|TestOperatorCommitReservationFence'
CMD_NEW_M2='TestOperatorFinishPrep|TestOperatorFinishPrepExpiry|TestOperatorFinishStale|TestOperatorFinishChecks|TestOperatorFinishUnsandboxed|TestOperatorCommitFailed|TestOperatorCommitReconcile|TestOperatorFinishJournalOnce|TestOperatorFinishManifestCapture|TestOperatorFinishWorldEvidence'
CMD_NEW_M3='TestOperatorFinishReadiness|TestOperatorFinishEligibility'
STORE_FINISH='TestOperatorFinishCompletedLaunch|TestOperatorLaunchReceiptLaw|TestOperatorStepGuards'
CMD_APPENDIX_FIXTURE_STEMS='TestOperatorFirstManagedEpisodeProcess|TestOperatorPreparationRechecksHeartbeatAfterWorktree|TestOperatorStartMintRecoveryRemedies|TestOperatorCurrentOwnerProcessRecovery|TestOperatorAttendedStartDisposition|TestOperatorStartOnce|TestOperatorStartCollisions|TestOperatorWorktreeWindow|TestOperatorWorktreeDispositionCAS|TestOperatorSpawnRefused|TestOperatorLaunchDispatchRechecks|TestOperatorLaunchCrashWindow|TestOperatorLaunchExit|TestOperatorStartBrowser'
````

````bash callers.sh
#!/bin/bash
set -euo pipefail
census() (
  cd "$1"
  label=$2
  all() { find cmd host -name '*.go' ! -name '*_test.go' -print0 | xargs -0 grep -E -nH "$1" || true; }
  count() {
    local name=$1 pattern=$2 expected=$3 hits n
    hits=$(all "$pattern" | grep -vE 'func .*\{|approval\.ConfirmOperatorAction\(' || true)
    n=$(printf '%s\n' "$hits" | grep -c . || true)
    [ "$n" -eq "$expected" ] || { printf '%s\n' "$hits"; echo "Unexpected $name $n, want $expected"; exit 1; }
    if [ "$n" -gt 0 ]; then printf '%s\n' "$hits"; fi
    echo "CALLERS $label $name $n"
  }
  grep -nF 'approval.ConfirmOperatorAction(' host/authority/judge_operator.go
  grep -nF 'finishDispatchOutcome(ctx, seam,' cmd/ailang-worldd/episode_finish.go
  grep -nF 'return e.runOperatorPreparationLoop(' cmd/ailang-worldd/operator_runner.go
  count request '\.RequestOperatorAction\(' 1
  count judge '\.JudgeOperatorAction\(' 1
  count store-confirm '\.ConfirmOperatorAction\(' 1
  count claim '\.ClaimOperatorAction\(' "$3"
  count private 'runEpisodeFinishWithPrep\(' "$4"
  count attended-commit 'worktree\.CommitWithCommitter\(' 1
  count local-commit '(\.git|in\.git)\.commit\(' "$4"
  count worktree-push 'worktree\.Push\(' 1
  count dispatch-push 'seam\.Push\(' 1
  count dispatch-label 'seam\.CreateLabel\(' 1
  count dispatch-view-pr 'seam\.ViewPR\(' 1
  count dispatch-create-pr 'seam\.CreatePR\(' 1
  count v3-put '\.PutVerbRequest\(' 1
)
census "$1" base 1 0
case "$3" in M1|M2) census "$2" "$3" 1 0;; M3) census "$2" M3 2 1;; *) exit 2;; esac
echo CALLER_CENSUS_OK
````

````bash outside.sh
#!/bin/bash
# TREE M1|M2|M3. Run on the corresponding planner stage BEFORE the executor.
set -euo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
M=$2
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
L=$(mktemp -d "$S/outside.$M.XXXXXX")
echo "Outside receipts in $L"
mkdir -m 1777 "$L/t1777"
mkdir -m 0700 "$L/t1777/x"
# Each leg reports independently even if an earlier one fails.
rc=0
bash "$S/accept.sh" "$T" "$M" outside > "$L/plain.raw" 2>&1 || rc=1
bash "$S/accept.sh" "$T" "$M" outside-race > "$L/race.raw" 2>&1 || rc=1
TMPDIR="$L/t1777/x" bash "$S/accept.sh" "$T" "$M" outside-1777 > "$L/1777.raw" 2>&1 || rc=1
grep -hE '^ACCEPT_OK|^RACE_BUDGET|^LEG |^--- (PASS|FAIL)|^panic:|^INSTRUMENT FAILURE|^UNINFORMATIVE UNDER SANDBOX' "$L/"*.raw || true
echo ALL_DONE
exit "$rc"
````

````bash whole-mutants.sh
#!/bin/bash
# LOGDIR [M1|M2|M3]; run the prioritized list outside, against complete base fail sets.
set -euo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
args=()
if [ $# -gt 1 ]; then args=(--stage "$2"); fi
case "$(tail -n 1 "$HOME/.ailang/state/world-iter449/baseline.done")" in 'ALL_DONE for 77dca5edc8206a1791743ebb3c64512f051c52f3') ;; *) exit 2;; esac
exec python3 "$S/mutate.py" "$1" --ids "$S/whole-mutant-ids.txt" --whole --outside "${args[@]}"
````

````bash commit-guards.sh
#!/bin/bash
# TREE LOGDIR. Final milestone only. Preserve every named failure and compare.
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
mkdir "$2" || exit 2
L=$(cd "$2" && pwd -P)
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang" GOPROXY=off
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
mkdir "$L/state" "$L/cache"
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
cd "$T"
rc=0
python3 "$S/bounded.py" 1280 go test -count=1 -timeout 20m ./host/verifygate/ ./host/broker/ > "$L/raw" 2>&1 || rc=$?
if grep -Eq 'panic: test timed out|BOUND EXCEEDED' "$L/raw"; then
  echo 'INSTRUMENT FAILURE: truncated commit guard'; exit 2
fi
python3 "$S/failset.py" "$L/raw" "$rc" > "$L/head" || exit 2
B="$HOME/.ailang/state/world-iter449/base"
test "$(tail -1 "$B/../baseline.done")" = 'ALL_DONE for 77dca5edc8206a1791743ebb3c64512f051c52f3' || exit 2
for tag in verifygate-plain broker-plain; do
  grep -Eq 'panic: test timed out|BOUND EXCEEDED' "$B/$tag.raw" && exit 2
  python3 "$S/failset.py" "$B/$tag.raw" "$(cat "$B/$tag.rc")" > "$L/$tag.fails" || exit 2
done
LC_ALL=C sort -u "$L/"*.fails > "$L/base"
LC_ALL=C sort -u "$L/head" > "$L/head.sorted"
LC_ALL=C comm -13 "$L/base" "$L/head.sorted" > "$L/new"
cat "$L/new"
test ! -s "$L/new" || exit 1
echo COMMIT_GUARDS_OK
````

````python bounded.py
import os, signal, subprocess, sys
p = subprocess.Popen(sys.argv[2:], start_new_session=True)
try:
    sys.exit(p.wait(timeout=int(sys.argv[1])))
except subprocess.TimeoutExpired:
    os.killpg(p.pid, signal.SIGKILL)
    p.wait()
    print("BOUND EXCEEDED", file=sys.stderr)
    sys.exit(124)
````

````python failset.py
import re
import sys
from pathlib import Path

prefix = 'github.com/sunholo-data/ailang-world/'
def parse(lines):
    pending = []
    tests = set()
    packages = set()
    for line in lines:
        m = re.match(r'^--- FAIL: ([^/\s]+) \(', line)
        if m:
            pending.append(m[1])
        m = re.match(r'^(FAIL|ok)\s+(' + re.escape(prefix) + r'\S+)', line)
        if m:
            pkg = m[2].removeprefix(prefix)
            packages.add(pkg)
            if m[1] == 'FAIL':
                if not pending:
                    raise ValueError('unattributed package failure: ' + pkg)
                tests.update((pkg, t) for t in pending)
            elif pending:
                raise ValueError('test failure without matching failed package')
            pending.clear()
    if pending:
        raise ValueError('test failure with no package summary')
    return tests, packages

control, _ = parse(['--- FAIL: TestParserPositive (0.00s)', 'FAIL\t' + prefix + 'host/store\t0.001s'])
if control != {('host/store', 'TestParserPositive')}:
    sys.exit('fail-set parser control failed')
print('FAILSET_CONTROL host/store TestParserPositive', file=sys.stderr)
tests, packages = parse(Path(sys.argv[1]).read_text(errors='replace').splitlines())
if not packages or (int(sys.argv[2]) != 0 and not tests):
    sys.exit('INSTRUMENT FAILURE: no package summary or nonzero rc with empty fail set')
for pkg, name in sorted(tests):
    print(pkg, name)
````

````python verify-scratch.py
"""Run the repository verifier with its explicit scratch temp paths and real Git."""
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys

tree, log, verifier, ref = sys.argv[1:]
tree, log = Path(tree), Path(log)
dest = log / verifier
dest.mkdir()
tmp = dest / "tmp"
tmp.mkdir()
copy = dest / "tree"
# A tracked checkout's Git metadata belongs to the controller, and only the
# committed-blob read needs it. All Go tests run against the scratch copy.
shutil.copytree(tree, copy, symlinks=True, ignore=shutil.ignore_patterns(".git"))
for name in (("verify_ail.sh", "verify_world_package.sh") if verifier == "verify_ail" else
             ("verify_go.sh", "test_tracked_binary_gate.sh")):
    p = copy / "scripts" / name
    src = p.read_text()
    dst = re.sub(r'mktemp -t (\w+\.XXXXXX)', r'mktemp "$S5E4_VERIFY_TMP/\1"', src)
    if name == "verify_go.sh":
        dst = dst.replace('${TMPDIR:-/tmp}/verify-go-evidence.XXXXXX',
                          '$S5E4_VERIFY_TMP/verify-go-evidence.XXXXXX')
        if ref:
            old = 'bash scripts/tracked_binary_gate.sh --allow-png'
            if dst.count(old) != 1:
                raise SystemExit("tracked gate anchor drift")
            dst = dst.replace(old, 'S5E4_BLOB_GIT="$S5E4_GIT_READ_TREE" ' + old)
    if name == "test_tracked_binary_gate.sh":
        dst = dst.replace('mktemp -d)', 'mktemp -d "$S5E4_VERIFY_TMP/adversarial.XXXXXX")')
    if name == "verify_go.sh":
        dst = dst.replace("go test ./... -count=1", "go test ./... -count=1 -timeout 20m")
    if src == dst:
        raise SystemExit(f"{name} temp path anchor drift")
    p.write_text(dst)
if verifier == "verify_go":
    p = copy / "scripts/tracked_binary_gate.sh"
    src = p.read_text()
    a = "import hashlib, re, struct, subprocess, sys, zlib"
    if src.count(a) != 1 or src.count('["git",') != 2:
        raise SystemExit("Git observation anchor drift")
    code = src.replace(a, a + "\nimport os\nread_tree = os.environ.get('S5E4_BLOB_GIT')\ngit = ['git', '-C', read_tree] if read_tree else ['git']")
    p.write_text(code.replace('["git", "', 'git + ["'))
env = dict(os.environ, S5E4_VERIFY_TMP=str(tmp), VERIFY_GO_LEGS="plain", S5E4_GIT_READ_TREE=ref)
env.pop("GOFLAGS", None)
print(f"Verifier temp paths relocated under {tmp}; real Git committed-blob reference {ref or '(not used)'}", flush=True)
sys.exit(subprocess.call([f"./scripts/{verifier}.sh"], cwd=copy, env=env))
````

````python mutate.py
"""Run a stage-anchored compiling mutation in a scratch overlay."""
import argparse
import json
import os
import re
import subprocess
from pathlib import Path

arg = argparse.ArgumentParser()
arg.add_argument("logdir", type=Path)
arg.add_argument("--ids", type=Path, required=True)
arg.add_argument("--stage")
arg.add_argument("--whole", action="store_true")
arg.add_argument("--outside", action="store_true")
a = arg.parse_args()
root = Path(__file__).resolve().parent
state = Path.home()/".ailang/state/world-iter449"
stages = state/"plan-scratch"
base = state/"base"
out = a.logdir.resolve()
out.mkdir()
definitions = {m["id"]: m for m in json.loads((root / "mutants.json").read_text())}
skip = (root / "selectors.sh").read_text().split("CMD_SKIP='", 1)[1].split("'", 1)[0]
env = dict(os.environ)
home = Path.home()
env.update(AILANG_BIN=str(home/".pinned-ailang/ailang"), WORLD_REQUIRE_SSH_KEYGEN="1",
           WORLD_EXEC_NODE="/opt/homebrew/bin/node", WORLD_TOOL_AILANG_BIN=str(home/".pinned-ailang-tools/v0.54.15/ailang"),
           GOPROXY="off", GIT_OPTIONAL_LOCKS="0")
for name in ("GOFLAGS","WORLD_CHROME","WORLD_REQUIRE_PASSKEY_BROWSER","WORLD_PASSKEY_RECORD_DIR"):
    env.pop(name, None)
ids = [line.strip() for line in a.ids.read_text().splitlines() if line.strip() and not line.startswith("#")]
prefix = "github.com/sunholo-data/ailang-world/"
timed_out = re.compile(r"panic: test timed out|BOUND EXCEEDED")
denied = re.compile(r"(?:bind|listen tcp).*(?:operation not permitted|permission denied)", re.I)

def record(line):
    print(line, flush=True)
    with (out/"ledger.tsv").open("a") as f:
        f.write(line+"\n")

def run(m):
    mid, pkg, test = m["id"], m["pkg"], m["test"]
    if m["outside"] and not a.outside:
        return "UNINFORMATIVE UNDER SANDBOX", "outside listener"
    tree = stages / ("st-" + m["stage"])
    entry = tree / m["file"]
    source = entry.read_text()
    if source.count(m["old"]) != 1:
        return "INSTRUMENT FAILURE", "anchor did not match uniquely"
    mu = out / mid
    mu.mkdir()
    for name in ("state", "cache"):
        (mu / name).mkdir()
    child_env = dict(env, AILANG_STATE_DIR=str(mu/"state"), AILANG_CACHE_DIR=str(mu/"cache"))
    changed = mu / entry.name
    changed.write_text(source.replace(m["old"], m["new"]))
    overlay = mu/"overlay.json"
    overlay.write_text(json.dumps({"Replace": {str(entry): str(changed)}}))
    rx = r"^" + re.escape(test) + r"$"
    whole_big = pkg in ("cmd/ailang-worldd","host/daemon")
    limit, cap = ("45m", "2800") if a.whole and whole_big else ("20m", "1260") if a.whole else ("8m", "540")
    argv = ["python3", str(root/"bounded.py"), cap, "go", "test", "-p", "1", "-v", "-count=1",
            "-timeout", limit, "-overlay", str(overlay), "./"+pkg+"/"]
    if not a.whole:
        argv += ["-run", rx]
    raw = mu/"raw"
    with raw.open("w") as f:
        rc = subprocess.run(argv, cwd=tree, env=child_env, stdout=f, stderr=subprocess.STDOUT).returncode
    log = raw.read_text(errors="replace")
    if timed_out.search(log) or "[build failed]" in log or "[setup failed]" in log:
        return "INSTRUMENT FAILURE", "timed out, setup or compile failure"
    if denied.search(log):
        return "UNINFORMATIVE UNDER SANDBOX", "listener or process denied"
    parsing = subprocess.run(["python3", str(root/"failset.py"), str(raw), str(rc)], capture_output=True, text=True)
    (mu/"parsed").write_text(parsing.stderr + parsing.stdout)
    print(parsing.stderr.strip(), flush=True)
    if parsing.returncode:
        return "INSTRUMENT FAILURE", "no attributable package result"
    failed = re.search(r"^--- FAIL: " + re.escape(test) + r" \(", log, re.M)
    passed = re.search(r"^--- PASS: " + re.escape(test) + r" \(", log, re.M)
    if a.whole:
        tag = "cmd" if pkg == "cmd/ailang-worldd" else pkg.split("/")[-1]
        first = (stages/"base-extra" if tag == "surfacepage" else base) / (tag+"-plain")
        prior, prior_raw = first.with_suffix(".fails"), first.with_suffix(".raw")
        if not prior.is_file() or not prior_raw.is_file() or timed_out.search(prior_raw.read_text(errors="replace")):
            return "UNMEASURED", "no complete pristine whole-package control"
        new = set(parsing.stdout.splitlines()) - set(prior.read_text().splitlines())
        (mu/"new").write_text("\n".join(sorted(new))+"\n")
        if failed and pkg+" "+test not in new:
            return "INSTRUMENT FAILURE", "named failure also present on pristine control"
    if failed:
        return "KILLED", test
    if passed:
        return "SURVIVOR", test
    return "INSTRUMENT FAILURE", "named test did not execute"

bad = False
for mid in ids:
    m = definitions[mid]
    if a.stage and m["stage"] != a.stage:
        continue
    disposition, detail = run(m)
    record(f'{mid}\t{m["stage"]}\t{m["pkg"]}\t{disposition}\t{detail}')
    bad |= disposition != "KILLED"
record("ALL_DONE")
raise SystemExit(int(bad))
````

````python mutant-spec.py
import json
import sys
from pathlib import Path

root = Path(__file__).resolve().parent
judge = Path(sys.argv[1])
originals = {}
for path in judge.glob("*.json"):
    doc = json.loads(path.read_text())
    for item in (doc.values() if isinstance(doc, dict) else doc):
        if isinstance(item, dict) and 'id' in item:
            originals[item["id"]] = item
items = []
cmd = "cmd/ailang-worldd/"
def add(mid, stage, file, old, new, test, origin="diff", outside=False):
    if "/" not in file:
        file = cmd + file + ".go"
    pkg = "cmd/ailang-worldd" if file.startswith(cmd) else "/".join(file.split("/")[:2])
    if test == "TestOperatorFinishJournalOnce":
        pkg = "cmd/ailang-worldd"
    source = (root / f"st-{stage}" / file).read_text()
    if source.count(old) != 1:
        raise SystemExit(f"{mid}: {source.count(old)} anchors in {file}")
    line = source[:source.index(old)].count("\n") + 1
    items.append(dict(id=mid, stage=stage, file=file, old=old, new=new,
                      test=test, pkg=pkg, origin=origin, outside=outside, line=line))

def carried(mid, stage, test, old=None, new=None):
    item = originals[mid]
    add(mid, stage, item["file"], item["before"] if old is None else old,
        item["after"] if new is None else new, test, "judge-448")

carried("e-P-kind", "M1", "TestOperatorFinishProofScope")
carried("e-P-surface", "M1", "TestOperatorFinishProofScope", "if !frozen.Options.ConfirmOnSurface {", "if false && !frozen.Options.ConfirmOnSurface {")
for mid in ("e-pin-skip-git", "e-pin-skip-all", "e-R-hash"):
    carried(mid, "M1", "TestOperatorFinishToolPins")
for mid in ("e-F-cfgworktree", "e-F-alternates", "e-F-staged", "e-F-digest", "e-F-visited", "e-F-count", "e-F-budget"):
    carried(mid, "M1", "TestOperatorFinishCaptureBounds")
carried("e-F-samefile2", "M1", "TestOperatorFinishReadReplacement")
carried("f-C-branch", "M1", "TestOperatorCommitIdentity")
for mid in ("f-C-recheck", "f-C-gottree", "e-C-cas"):
    carried(mid, "M1", "TestOperatorCommitReservationFence")
carried("e-C-authoremail", "M1", "TestOperatorCommitIdentity", "GIT_AUTHOR_EMAIL=ailang-world@invalid", "GIT_AUTHOR_EMAIL=x@y")
carried("e-store-commit-once", "M2", "TestOperatorFinishJournalOnce",
        "AND (? <> 'commit' OR NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind = 'commit'))",
        "AND (? <> 'commit' OR EXISTS (SELECT 1 FROM operator_actions WHERE id = ?))")
for mid, old, new, test in [
    ("D-config-swap", "!bytes.Equal(current, data)", "false && !bytes.Equal(current, data)", "TestOperatorFinishConfigReplacement"),
    ("D-config-readback", 'current, _, err := operatorFinishRead(g.common, "config", 64<<10)', 'current, err := data, error(nil)', "TestOperatorFinishConfigReplacement"),
    ("D-name-empty", 'if tracked != "" {', "if true {", "TestOperatorFinishCaptureBounds"),
    ("D-author", 'author+">" != "author "+agent+" via AILANG World <ailang-world@invalid>"', 'false && author+">" != "author "+agent+" via AILANG World <ailang-world@invalid>"', "TestOperatorCommitIdentity"),
    ("D-author-time", "len(strings.Fields(authorTime)) != 2", "len(strings.Fields(authorTime)) < 0", "TestOperatorCommitIdentity")
]:
    add(mid, "M1", "operator_finish_tree", old, new, test)
for mid, old, new, test in [
    ("J-claim", "if !w.prep.claimed() {", "if false && !w.prep.claimed() {", "TestOperatorFinishPrep"),
    ("J-manifest", "if !w.boundManifest(w.generation) {", "if false && !w.boundManifest(w.generation) {", "TestOperatorFinishPrep"),
    ("J-unsandboxed", "if check.Unsandboxed {", "if false && check.Unsandboxed {", "TestOperatorFinishUnsandboxed"),
    ("J-protected", " || len(m.Protected) != 0", "", "TestOperatorFinishPrep"),
    ("J-unattributed", "len(m.Unattributed) != 0 || ", "", "TestOperatorFinishPrep"),
    ("J-bytes", "!reflect.DeepEqual(snap, w.manifest.Commit.Snapshot)", "false && !reflect.DeepEqual(snap, w.manifest.Commit.Snapshot)", "TestOperatorFinishStale"),
    ("J-check-version", 'Identity: r.ID + ":" + identity, Version: 1', 'Identity: r.ID + ":" + identity, Version: 2', "TestOperatorCommitReconcile"),
    ("J-early-check", 'key := w.step("checks", fmt.Sprintf("%02d", index))\n\tif err', 'key := w.step("checks", fmt.Sprintf("%02d", index))\n\t_, _ = w.run(ctx, check)\n\tif err', "TestOperatorCommitReconcile"),
    ("J-unknown", "if j.OpenCheck || len(j.Checks)", "if false && j.OpenCheck || len(j.Checks)", "TestOperatorFinishChecks"),
    ("J-fail-check", "passed := result.ExitCode == 0", "passed := result.ExitCode >= 0", "TestOperatorFinishChecks"),
    ("J-outcome-invalid", "err != nil || entry < 1", "err != nil", "TestOperatorFinishChecks"),
    ("J-old-failure", 'if old.Result != "passed" {', 'if false && old.Result != "passed" {', "TestOperatorFinishChecks"),
    ("J-after-check", "contents, err = w.inputs(ctx)\n\t\tif err != nil {\n\t\t\treturn err", "contents, err = w.inputs(ctx)\n\t\tif false && err != nil {\n\t\t\treturn err", "TestOperatorFinishChecks"),
    ("J-early-commit", 'key := w.step("commit", tree)', 'key := w.step("commit", tree)\n\t_ = w.commit(ctx, m.Commit, contents, tree)', "TestOperatorCommitReconcile"),
    ("J-no-ref", "if errors.Is(commitErr, errOperatorCommitNoRef) {", "if false && errors.Is(commitErr, errOperatorCommitNoRef) {", "TestOperatorCommitFailed"),
    ("J-observe", 'return w.observe(key, store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: sha})', '_ = sha\n\treturn nil', "TestOperatorCommitReconcile"),
    ("J-prep-expiry", 'if err := w.check(ctx, check, i+len(j.Checks)); err != nil {\n\t\t\treturn err\n\t\t}', 'if err := w.check(ctx, check, i+len(j.Checks)); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tw.fresh = func(context.Context) error { return nil }\n\t\tw.now = func() int64 { return w.prep.row.CreatedAt }', "TestOperatorFinishPrepExpiry"),
]:
    add(mid, "M2", "operator_finish_local", old, new, test)

add("J-input-first", "M2", "operator_finish_local",
    "contents, err := w.inputs(ctx)\n\tif err != nil {", "contents, err := w.inputs(ctx)\n\tif false && err != nil {",
    "TestOperatorFinishStale")
add("J-observed-commit", "M2", "operator_finish_local",
    'if j.CommitPhase == "observed" {', 'if false && j.CommitPhase == "observed" {',
    "TestOperatorCommitReconcile")
add("J-tree-identity", "M2", "operator_finish_local",
    'key := w.step("commit", tree)', 'key := w.step("commit", strings.Repeat("a", 40))',
    "TestOperatorCommitReconcile")
add("J-check-observe", "M2", "operator_finish_local",
    "if err := w.observe(key, store.OperatorFinishCheck{Name: check.Name, Result: class, Ref: obj.Hash, Entry: entry}); err != nil {",
    "_ = class\n\tif err := error(nil); err != nil {",
    "TestOperatorCommitReconcile")
for effect, identity in [("check", 'key := w.step("checks", fmt.Sprintf("%02d", index))'), ("commit", 'key := w.step("commit", tree)')]:
    add("J-reserve-"+effect, "M2", "operator_finish_local", identity+"\n\tif err := w.reserve(key); err != nil {",
        identity+"\n\tif err := error(nil); err != nil {", "TestOperatorCommitReconcile")

for mid, old, new, test in [
    ("D-outcome-name", " || outcome.Task != receipt.Name", "", "TestOperatorFinishWorldEvidence"),
    ("D-outcome-result", " || outcome.Outcome.Class != wantClass", " || outcome.Outcome.Class != wantClass && false", "TestOperatorFinishWorldEvidence"),
    ("D-world-membership", "err == nil && len(entries) == 1 && entries[0] == receipt.Entry", "err == nil && len(entries) >= 0", "TestOperatorFinishWorldEvidence"),
    ("D-candidate-expiry", "created_at<=? AND expires_at>?", "created_at<=? AND expires_at<=?", "TestOperatorFinishJournalOnce"),
    ("D-observed-candidate", "AND NOT EXISTS (SELECT 1 FROM operator_action_steps s", "AND NOT EXISTS (SELECT 1 FROM operator_action_steps s", "TestOperatorFinishJournalOnce"),
    ("D-launch-version", "s.phase='observed' AND s.version=4\n\t\tORDER BY a.created_at DESC,a.id LIMIT 16", "s.phase='observed' AND s.version>=3\n\t\tORDER BY a.created_at DESC,a.id LIMIT 16", "TestOperatorFinishCompletedLaunch")
]:
    if mid in {"D-world-membership", "D-outcome-name", "D-outcome-result"}:
        path = "operator_finish_manifest"
    else:
        path = "host/store/operator_finish.go"
    add(mid, "M2", path, old if mid != "D-observed-candidate" else "s.kind='commit' AND s.phase='observed')",
        new if mid != "D-observed-candidate" else "s.kind='commit' AND s.phase='reserved')", test)

suffix = Path(sys.argv[2])
for m in json.loads(suffix.read_text()):
    if m["id"] in {"R-auth", "R-prep-card", "R-status", "B-prep-task", "B-prep-kind", "B-prep-flag"}:
        add(m["id"], "M3", m["file"], m.get("old", m.get("before")), m.get("new", m.get("after")), m["test"], "appendix-diff", m["outside"])
add("D-lock-conflict", "M3", "operator_finish_runner", "if errors.Is(err, syscall.EWOULDBLOCK) {", "if false && errors.Is(err, syscall.EWOULDBLOCK) {", "TestOperatorFinishRouteWorker", "appendix-diff", True)

add("D-worker-drift", "M3", "operator_finish_runner",
    "if !reflect.DeepEqual(reviewed, m) ||", "if false && !reflect.DeepEqual(reviewed, m) ||",
    "TestOperatorFinishRouteWorker", "appendix-diff", True)
add("D-finish-admission", "M3", "host/daemon/surface_operator.go",
    'in.Kind == "start" && in.Headless || in.Kind == "finish"', 'in.Kind == "start" && in.Headless',
    "TestOperatorFinishAdmission", "appendix-diff", True)
add("D-status-redaction", "M3", "host/daemon/surface_operator_read.go",
    "if statusJSON {\n\t\twriteApprovalJSON(w, http.StatusOK, status)",
    "if statusJSON {\n\t\tstatus.Reason += string(row.Action)\n\t\twriteApprovalJSON(w, http.StatusOK, status)",
    "TestOperatorFinishAdmission", "appendix-diff", True)

add("D-catalogue-workspace", "M3", "operator_catalogue", "\t\tsrc.workspace = r.Workspace\n", "",
    "TestOperatorFinishRouteWorker", "appendix-diff", True)

(root / "mutants.json").write_text(json.dumps(items, indent=2) + "\n")
````

````python patches.py
import pathlib
import re
import subprocess
import sys

root = pathlib.Path(__file__).resolve().parent
pairs = [("base-tree", "st-M1")]
for stage in sys.argv[1:]:
    previous = "st-M" + str(int(stage[1:]) - 1)
    pairs += [(previous, "st-" + stage), ("base-tree", "st-" + stage)]
for old, new in pairs:
    m = new.removeprefix("st-")
    tag = m if old != "base-tree" or m == "M1" else m + "-cumulative"
    cmd = ["git", "-c", "core.quotePath=false", "diff", "--no-index", "--no-color"]
    patch, stat = [], []
    for folder in ("cmd", "host", "docs", "app"):
        a, b = f"{old}/{folder}", f"{new}/{folder}"
        if not (root / a).exists() and not (root / b).exists():
            continue
        for flags, dest in ((["-M", "--binary"], patch), (["-M", "--numstat"], stat)):
            p = subprocess.run([*cmd, *flags, "--", a, b], cwd=root, capture_output=True, text=True)
            if p.returncode not in (0, 1):
                sys.exit(p.stderr)
            for line in p.stdout.splitlines():
                if line.startswith(("diff --git ", "--- ", "+++ ")):
                    for side in ("a", "b"):
                        for stage in (old, new):
                            line = line.replace(f"{side}/{stage}/", f"{side}/")
                dest.append(line)
    (root / (tag + ".patch")).write_text("\n".join(patch) + "\n")
    (root / (tag + ".numstat")).write_text("\n".join(stat) + "\n")
    status, changed = [], []
    plus = minus = 0
    for entry in stat:
        added, deleted, path = entry.split("\t")
        plus += int(added)
        minus += int(deleted)
        path = path.replace("/dev/null => ", "").removesuffix(" => /dev/null")
        path = re.sub(r"^" + re.escape(new) + "/", "", path)
        path = path.replace("{" + old + " => " + new + "}/", "")
        path = re.sub(r"^" + re.escape(old) + "/", "", path)
        base, head = (root / old / path).is_file(), (root / new / path).is_file()
        if not base and not head:
            sys.exit(f"unresolved path {path!r}")
        status.append(f"{'M' if base and head else 'A' if head else 'D'}\t{path}")
        if head:
            changed.append(path)
    (root / ("status-" + tag)).write_text("\n".join(sorted(status, key=lambda row: row.split("\t")[1])) + "\n")
    (root / ("changed-" + tag)).write_text("\n".join(sorted(changed)) + "\n")
    print(tag, f"+{plus}", f"-{minus}", plus + minus)
````

````json mutants.json
[
  {
    "id": "e-P-kind",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_preflight.go",
    "old": "frozen.Kind != \"finish\" || ",
    "new": "",
    "test": "TestOperatorFinishProofScope",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 26
  },
  {
    "id": "e-P-surface",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_preflight.go",
    "old": "if !frozen.Options.ConfirmOnSurface {",
    "new": "if false && !frozen.Options.ConfirmOnSurface {",
    "test": "TestOperatorFinishProofScope",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 29
  },
  {
    "id": "e-pin-skip-git",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_service.go",
    "old": "name == \"gh\" ||",
    "new": "name == \"gh\" || name == \"git\" ||",
    "test": "TestOperatorFinishToolPins",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 165
  },
  {
    "id": "e-pin-skip-all",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_service.go",
    "old": "if len(c.FinishRemotes) != 0 && (name == \"gh\"",
    "new": "if len(c.FinishRemotes) != 0 || (name == \"gh\"",
    "test": "TestOperatorFinishToolPins",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 165
  },
  {
    "id": "e-R-hash",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "old": " && serviceFileHash(real) == c.Executables[real]",
    "new": "",
    "test": "TestOperatorFinishToolPins",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 42
  },
  {
    "id": "e-F-cfgworktree",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "filepath.Join(g.admin, \"config.worktree\"), ",
    "new": "",
    "test": "TestOperatorFinishCaptureBounds",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 125
  },
  {
    "id": "e-F-alternates",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "filepath.Join(g.common, \"objects\", \"info\", \"alternates\"), ",
    "new": "",
    "test": "TestOperatorFinishCaptureBounds",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 126
  },
  {
    "id": "e-F-staged",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "err != nil || staged != \"\" {",
    "new": "err != nil || false && staged != \"\" {",
    "test": "TestOperatorFinishCaptureBounds",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 185
  },
  {
    "id": "e-F-digest",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "p := store.OperatorFinishPath{Path: path, Mode: \"100644\", Digest: finishSHA(b)}",
    "new": "p := store.OperatorFinishPath{Path: path, Mode: \"100644\", Digest: \"fixed\"}",
    "test": "TestOperatorFinishCaptureBounds",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 240
  },
  {
    "id": "e-F-visited",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "visited > 16384",
    "new": "visited > 1<<30",
    "test": "TestOperatorFinishCaptureBounds",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 260
  },
  {
    "id": "e-F-count",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "len(names) > 257",
    "new": "len(names) > 1<<30",
    "test": "TestOperatorFinishCaptureBounds",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 221
  },
  {
    "id": "e-F-budget",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "\t\tbudget -= int64(len(b))\n",
    "new": "",
    "test": "TestOperatorFinishCaptureBounds",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 251
  },
  {
    "id": "e-F-samefile2",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "!os.SameFile(before, after) || ",
    "new": "",
    "test": "TestOperatorFinishReadReplacement",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 80
  },
  {
    "id": "f-C-branch",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "strings.TrimSpace(branch) != spec.Snapshot.Branch",
    "new": "strings.TrimSpace(branch) != strings.TrimSpace(branch)",
    "test": "TestOperatorCommitIdentity",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 313
  },
  {
    "id": "f-C-recheck",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "!reflect.DeepEqual(current, spec.Snapshot)",
    "new": "!reflect.DeepEqual(current, current)",
    "test": "TestOperatorCommitReservationFence",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 369
  },
  {
    "id": "f-C-gottree",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "got != tree",
    "new": "got != got",
    "test": "TestOperatorCommitReservationFence",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 358
  },
  {
    "id": "e-C-cas",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "update-ref\", spec.Snapshot.Branch, head, spec.Snapshot.Head)",
    "new": "update-ref\", spec.Snapshot.Branch, head)",
    "test": "TestOperatorCommitReservationFence",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 375
  },
  {
    "id": "e-C-authoremail",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "GIT_AUTHOR_EMAIL=ailang-world@invalid",
    "new": "GIT_AUTHOR_EMAIL=x@y",
    "test": "TestOperatorCommitIdentity",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 362
  },
  {
    "id": "e-store-commit-once",
    "stage": "M2",
    "file": "host/store/operator_steps.go",
    "old": "AND (? <> 'commit' OR NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind = 'commit'))",
    "new": "AND (? <> 'commit' OR EXISTS (SELECT 1 FROM operator_actions WHERE id = ?))",
    "test": "TestOperatorFinishJournalOnce",
    "pkg": "cmd/ailang-worldd",
    "origin": "judge-448",
    "outside": false,
    "line": 37
  },
  {
    "id": "D-config-swap",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "!bytes.Equal(current, data)",
    "new": "false && !bytes.Equal(current, data)",
    "test": "TestOperatorFinishConfigReplacement",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 160
  },
  {
    "id": "D-config-readback",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "current, _, err := operatorFinishRead(g.common, \"config\", 64<<10)",
    "new": "current, err := data, error(nil)",
    "test": "TestOperatorFinishConfigReplacement",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 159
  },
  {
    "id": "D-name-empty",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "if tracked != \"\" {",
    "new": "if true {",
    "test": "TestOperatorFinishCaptureBounds",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 207
  },
  {
    "id": "D-author",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "author+\">\" != \"author \"+agent+\" via AILANG World <ailang-world@invalid>\"",
    "new": "false && author+\">\" != \"author \"+agent+\" via AILANG World <ailang-world@invalid>\"",
    "test": "TestOperatorCommitIdentity",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 340
  },
  {
    "id": "D-author-time",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "len(strings.Fields(authorTime)) != 2",
    "new": "len(strings.Fields(authorTime)) < 0",
    "test": "TestOperatorCommitIdentity",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 340
  },
  {
    "id": "J-claim",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if !w.prep.claimed() {",
    "new": "if false && !w.prep.claimed() {",
    "test": "TestOperatorFinishPrep",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 80
  },
  {
    "id": "J-manifest",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if !w.boundManifest(w.generation) {",
    "new": "if false && !w.boundManifest(w.generation) {",
    "test": "TestOperatorFinishPrep",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 84
  },
  {
    "id": "J-unsandboxed",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if check.Unsandboxed {",
    "new": "if false && check.Unsandboxed {",
    "test": "TestOperatorFinishUnsandboxed",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 88
  },
  {
    "id": "J-protected",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": " || len(m.Protected) != 0",
    "new": "",
    "test": "TestOperatorFinishPrep",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 93
  },
  {
    "id": "J-unattributed",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "len(m.Unattributed) != 0 || ",
    "new": "",
    "test": "TestOperatorFinishPrep",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 93
  },
  {
    "id": "J-bytes",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "!reflect.DeepEqual(snap, w.manifest.Commit.Snapshot)",
    "new": "false && !reflect.DeepEqual(snap, w.manifest.Commit.Snapshot)",
    "test": "TestOperatorFinishStale",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 70
  },
  {
    "id": "J-check-version",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "Identity: r.ID + \":\" + identity, Version: 1",
    "new": "Identity: r.ID + \":\" + identity, Version: 2",
    "test": "TestOperatorCommitReconcile",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 45
  },
  {
    "id": "J-early-check",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "key := w.step(\"checks\", fmt.Sprintf(\"%02d\", index))\n\tif err",
    "new": "key := w.step(\"checks\", fmt.Sprintf(\"%02d\", index))\n\t_, _ = w.run(ctx, check)\n\tif err",
    "test": "TestOperatorCommitReconcile",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 171
  },
  {
    "id": "J-unknown",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if j.OpenCheck || len(j.Checks)",
    "new": "if false && j.OpenCheck || len(j.Checks)",
    "test": "TestOperatorFinishChecks",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 109
  },
  {
    "id": "J-fail-check",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "passed := result.ExitCode == 0",
    "new": "passed := result.ExitCode >= 0",
    "test": "TestOperatorFinishChecks",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 183
  },
  {
    "id": "J-outcome-invalid",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "err != nil || entry < 1",
    "new": "err != nil",
    "test": "TestOperatorFinishChecks",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 194
  },
  {
    "id": "J-old-failure",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if old.Result != \"passed\" {",
    "new": "if false && old.Result != \"passed\" {",
    "test": "TestOperatorFinishChecks",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 116
  },
  {
    "id": "J-after-check",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "contents, err = w.inputs(ctx)\n\t\tif err != nil {\n\t\t\treturn err",
    "new": "contents, err = w.inputs(ctx)\n\t\tif false && err != nil {\n\t\t\treturn err",
    "test": "TestOperatorFinishChecks",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 131
  },
  {
    "id": "J-early-commit",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "key := w.step(\"commit\", tree)",
    "new": "key := w.step(\"commit\", tree)\n\t_ = w.commit(ctx, m.Commit, contents, tree)",
    "test": "TestOperatorCommitReconcile",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 150
  },
  {
    "id": "J-no-ref",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if errors.Is(commitErr, errOperatorCommitNoRef) {",
    "new": "if false && errors.Is(commitErr, errOperatorCommitNoRef) {",
    "test": "TestOperatorCommitFailed",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 158
  },
  {
    "id": "J-observe",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "return w.observe(key, store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: sha})",
    "new": "_ = sha\n\treturn nil",
    "test": "TestOperatorCommitReconcile",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 224
  },
  {
    "id": "J-prep-expiry",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if err := w.check(ctx, check, i+len(j.Checks)); err != nil {\n\t\t\treturn err\n\t\t}",
    "new": "if err := w.check(ctx, check, i+len(j.Checks)); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tw.fresh = func(context.Context) error { return nil }\n\t\tw.now = func() int64 { return w.prep.row.CreatedAt }",
    "test": "TestOperatorFinishPrepExpiry",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 128
  },
  {
    "id": "J-input-first",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "contents, err := w.inputs(ctx)\n\tif err != nil {",
    "new": "contents, err := w.inputs(ctx)\n\tif false && err != nil {",
    "test": "TestOperatorFinishStale",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 120
  },
  {
    "id": "J-observed-commit",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if j.CommitPhase == \"observed\" {",
    "new": "if false && j.CommitPhase == \"observed\" {",
    "test": "TestOperatorCommitReconcile",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 103
  },
  {
    "id": "J-tree-identity",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "key := w.step(\"commit\", tree)",
    "new": "key := w.step(\"commit\", strings.Repeat(\"a\", 40))",
    "test": "TestOperatorCommitReconcile",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 150
  },
  {
    "id": "J-check-observe",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if err := w.observe(key, store.OperatorFinishCheck{Name: check.Name, Result: class, Ref: obj.Hash, Entry: entry}); err != nil {",
    "new": "_ = class\n\tif err := error(nil); err != nil {",
    "test": "TestOperatorCommitReconcile",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 201
  },
  {
    "id": "J-reserve-check",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "key := w.step(\"checks\", fmt.Sprintf(\"%02d\", index))\n\tif err := w.reserve(key); err != nil {",
    "new": "key := w.step(\"checks\", fmt.Sprintf(\"%02d\", index))\n\tif err := error(nil); err != nil {",
    "test": "TestOperatorCommitReconcile",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 171
  },
  {
    "id": "J-reserve-commit",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "key := w.step(\"commit\", tree)\n\tif err := w.reserve(key); err != nil {",
    "new": "key := w.step(\"commit\", tree)\n\tif err := error(nil); err != nil {",
    "test": "TestOperatorCommitReconcile",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 150
  },
  {
    "id": "D-outcome-name",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_manifest.go",
    "old": " || outcome.Task != receipt.Name",
    "new": "",
    "test": "TestOperatorFinishWorldEvidence",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 95
  },
  {
    "id": "D-outcome-result",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_manifest.go",
    "old": " || outcome.Outcome.Class != wantClass",
    "new": " || outcome.Outcome.Class != wantClass && false",
    "test": "TestOperatorFinishWorldEvidence",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 95
  },
  {
    "id": "D-world-membership",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_manifest.go",
    "old": "err == nil && len(entries) == 1 && entries[0] == receipt.Entry",
    "new": "err == nil && len(entries) >= 0",
    "test": "TestOperatorFinishWorldEvidence",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 99
  },
  {
    "id": "D-candidate-expiry",
    "stage": "M2",
    "file": "host/store/operator_finish.go",
    "old": "created_at<=? AND expires_at>?",
    "new": "created_at<=? AND expires_at<=?",
    "test": "TestOperatorFinishJournalOnce",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 89
  },
  {
    "id": "D-observed-candidate",
    "stage": "M2",
    "file": "host/store/operator_finish.go",
    "old": "s.kind='commit' AND s.phase='observed')",
    "new": "s.kind='commit' AND s.phase='reserved')",
    "test": "TestOperatorFinishJournalOnce",
    "pkg": "cmd/ailang-worldd",
    "origin": "diff",
    "outside": false,
    "line": 90
  },
  {
    "id": "D-launch-version",
    "stage": "M2",
    "file": "host/store/operator_finish.go",
    "old": "s.phase='observed' AND s.version=4\n\t\tORDER BY a.created_at DESC,a.id LIMIT 16",
    "new": "s.phase='observed' AND s.version>=3\n\t\tORDER BY a.created_at DESC,a.id LIMIT 16",
    "test": "TestOperatorFinishCompletedLaunch",
    "pkg": "host/store",
    "origin": "diff",
    "outside": false,
    "line": 101
  },
  {
    "id": "R-auth",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if err := prepareOperatorFinish(ctx, c, m.Remote.URL); err != nil {",
    "new": "if err := prepareOperatorFinish(ctx, c, m.Remote.URL); false && err != nil {",
    "test": "TestOperatorFinishReadiness",
    "pkg": "cmd/ailang-worldd",
    "origin": "appendix-diff",
    "outside": false,
    "line": 215
  },
  {
    "id": "R-prep-card",
    "stage": "M3",
    "file": "host/daemon/surface_operator_confirm.go",
    "old": " || a.Kind == \"finish\" && a.Options.ConfirmOnSurface",
    "new": "",
    "test": "TestOperatorFinishAdmission",
    "pkg": "host/daemon",
    "origin": "appendix-diff",
    "outside": true,
    "line": 96
  },
  {
    "id": "R-status",
    "stage": "M3",
    "file": "host/daemon/surface_operator_read.go",
    "old": "status.LocalCommit = j.Commit.SHA",
    "new": "status.LocalCommit = \"\"",
    "test": "TestOperatorFinishAdmission",
    "pkg": "host/daemon",
    "origin": "appendix-diff",
    "outside": true,
    "line": 132
  },
  {
    "id": "B-prep-task",
    "stage": "M3",
    "file": "host/surfacepage/app/app.js",
    "old": "argv[2] === 'finish' ? task === ''",
    "new": "argv[2] === 'finish' ? true",
    "test": "TestOperatorFinishBarReview",
    "pkg": "host/surfacepage",
    "origin": "appendix-diff",
    "outside": false,
    "line": 1448
  },
  {
    "id": "B-prep-kind",
    "stage": "M3",
    "file": "host/surfacepage/app/app.js",
    "old": "b.kind !== kind",
    "new": "false && b.kind !== kind",
    "test": "TestOperatorFinishBarReview",
    "pkg": "host/surfacepage",
    "origin": "appendix-diff",
    "outside": false,
    "line": 1459
  },
  {
    "id": "B-prep-flag",
    "stage": "M3",
    "file": "host/surfacepage/app/app.js",
    "old": "argv.length === 8 && argv[7] === '--ready'",
    "new": "argv.length >= 8 && argv[7] === '--ready'",
    "test": "TestOperatorFinishBarReview",
    "pkg": "host/surfacepage",
    "origin": "appendix-diff",
    "outside": false,
    "line": 1446
  },
  {
    "id": "D-lock-conflict",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if errors.Is(err, syscall.EWOULDBLOCK) {",
    "new": "if false && errors.Is(err, syscall.EWOULDBLOCK) {",
    "test": "TestOperatorFinishRouteWorker",
    "pkg": "cmd/ailang-worldd",
    "origin": "appendix-diff",
    "outside": true,
    "line": 220
  },
  {
    "id": "D-worker-drift",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if !reflect.DeepEqual(reviewed, m) ||",
    "new": "if false && !reflect.DeepEqual(reviewed, m) ||",
    "test": "TestOperatorFinishRouteWorker",
    "pkg": "cmd/ailang-worldd",
    "origin": "appendix-diff",
    "outside": true,
    "line": 253
  },
  {
    "id": "D-finish-admission",
    "stage": "M3",
    "file": "host/daemon/surface_operator.go",
    "old": "in.Kind == \"start\" && in.Headless || in.Kind == \"finish\"",
    "new": "in.Kind == \"start\" && in.Headless",
    "test": "TestOperatorFinishAdmission",
    "pkg": "host/daemon",
    "origin": "appendix-diff",
    "outside": true,
    "line": 29
  },
  {
    "id": "D-status-redaction",
    "stage": "M3",
    "file": "host/daemon/surface_operator_read.go",
    "old": "if statusJSON {\n\t\twriteApprovalJSON(w, http.StatusOK, status)",
    "new": "if statusJSON {\n\t\tstatus.Reason += string(row.Action)\n\t\twriteApprovalJSON(w, http.StatusOK, status)",
    "test": "TestOperatorFinishAdmission",
    "pkg": "host/daemon",
    "origin": "appendix-diff",
    "outside": true,
    "line": 146
  },
  {
    "id": "D-catalogue-workspace",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_catalogue.go",
    "old": "\t\tsrc.workspace = r.Workspace\n",
    "new": "",
    "test": "TestOperatorFinishRouteWorker",
    "pkg": "cmd/ailang-worldd",
    "origin": "appendix-diff",
    "outside": true,
    "line": 39
  }
]
````

````text named-mutant-ids.txt
e-P-kind
e-P-surface
e-pin-skip-git
e-pin-skip-all
e-R-hash
e-F-cfgworktree
e-F-alternates
e-F-staged
e-F-digest
e-F-visited
e-F-count
e-F-budget
e-F-samefile2
f-C-branch
f-C-recheck
f-C-gottree
e-C-cas
e-C-authoremail
e-store-commit-once
D-config-swap
D-config-readback
D-name-empty
D-author
D-author-time
J-claim
J-manifest
J-unsandboxed
J-protected
J-unattributed
J-bytes
J-check-version
J-early-check
J-unknown
J-fail-check
J-outcome-invalid
J-old-failure
J-after-check
J-early-commit
J-no-ref
J-observe
J-prep-expiry
J-input-first
J-observed-commit
J-tree-identity
J-check-observe
J-reserve-check
J-reserve-commit
D-outcome-name
D-outcome-result
D-world-membership
D-candidate-expiry
D-observed-candidate
D-launch-version
R-auth
R-prep-card
R-status
B-prep-task
B-prep-kind
B-prep-flag
D-lock-conflict
D-worker-drift
D-finish-admission
D-status-redaction
D-catalogue-workspace
````

````text whole-mutant-ids.txt
e-P-kind
e-P-surface
e-pin-skip-git
e-pin-skip-all
e-R-hash
e-F-digest
e-F-visited
e-F-count
e-F-budget
e-F-cfgworktree
e-F-alternates
e-F-staged
e-F-samefile2
f-C-branch
f-C-recheck
e-C-cas
f-C-gottree
e-C-authoremail
e-store-commit-once
J-early-commit
J-unknown
J-prep-expiry
D-launch-version
R-prep-card
R-status
````

## 7. Apply, accept and snapshot

Apply M1 and M2 in order from the sprint worktree. The controller commits the two plan artifacts separately. The executor applies product files with `/usr/bin/patch -p1 -E` and copies cumulative snapshots to `$HOME/.ailang/state/world-iter449/snap/<M>`. A successful milestone prints `ACCEPT_OK <M> <mode>` and `APPLY_ACCEPT_SNAPSHOT_OK <M> <mode>`. Require both markers and `FINAL_APPLY_RC 0`. `STATUS` lists adds, modifications and deletions without rename inference, `DELETED` contains deleted relative paths, and `MANIFEST.sha256` covers all retained snapshot files. The two planning files are excluded. The controller runs the outside legs on the planner's stage before the executor starts.


````bash
set -euo pipefail
B=77dca5edc8206a1791743ebb3c64512f051c52f3
STATE="$HOME/.ailang/state/world-iter449"
S="$STATE/plan-scratch"
WT=${WT:-/Users/voightkampff/dev/sunholo-data/.wt-world-iter449-row243-v4s9b}
PLAN=${PLAN:-$WT/design_docs/planned/w-world-address-bar-verbs-v4s9b-sprint-plan.md}
TREE=${TREE:-$WT}
TOOLS=${TOOLS:-$STATE/exec-tools}
MILESTONES=${MILESTONES:-M1 M2}
MODE=${MODE:-fast}
REHEARSAL=${REHEARSAL:-0}
export GIT_OPTIONAL_LOCKS=0
mkdir -p "$TOOLS"
grab() {
  awk -v fence="$1" '$0==fence {f=1;next} f&&$0=="````" {exit} f' "$PLAN" > "$2"
  test -s "$2"
}
for f in accept.sh gates.sh snapshot.sh selectors.sh callers.sh outside.sh whole-mutants.sh commit-guards.sh; do
  grab '````bash '"$f" "$TOOLS/$f"
done
for f in bounded.py failset.py verify-scratch.py mutate.py; do
  grab '````python '"$f" "$TOOLS/$f"
done
grab '````json mutants.json' "$TOOLS/mutants.json"
grab '````text whole-mutant-ids.txt' "$TOOLS/whole-mutant-ids.txt"
grab '````text named-mutant-ids.txt' "$TOOLS/named-mutant-ids.txt"
if [ "$REHEARSAL" = 1 ]; then
  test ! -e "$TREE"
  mkdir -p "$TREE"
  git -C "$WT" archive "$B" | tar -x -C "$TREE"
else
  git -C "$TREE" merge-base --is-ancestor "$B" HEAD
  git -C "$TREE" diff --quiet "$B" HEAD -- cmd host docs app
fi
for M in $MILESTONES; do
  case "$M" in M1|M2) ;; M3) test "$REHEARSAL" = 1;; *) exit 2;; esac
  grab '````diff '"$M" "$TOOLS/$M.patch"
  (cd "$TREE" && /usr/bin/patch -p1 -E --dry-run < "$TOOLS/$M.patch") > "$TOOLS/$M.dry.raw"
  if grep -Eiq 'offset|fuzz|reject|FAILED' "$TOOLS/$M.dry.raw"; then
    cat "$TOOLS/$M.dry.raw"; exit 1
  fi
  (cd "$TREE" && /usr/bin/patch -p1 -E -s < "$TOOLS/$M.patch")
  bash "$TOOLS/accept.sh" "$TREE" "$M" "$MODE"
  if [ "$REHEARSAL" = 1 ]; then
    case "$M" in M1) tag=M1;; *) tag="$M-cumulative";; esac
    SNAP_STATUS_LIST="$S/status-$tag" SNAP_PARENT="$TOOLS/rehearsal-snap" bash "$TOOLS/snapshot.sh" "$TREE" "$M"
    for dir in cmd host docs app; do
      if [ -d "$TREE/$dir" ] || [ -d "$S/st-$M/$dir" ]; then
        diff -r "$TREE/$dir" "$S/st-$M/$dir"
        printf 'ROUNDTRIP_OK %s %s\n' "$M" "$dir"
      fi
    done
  else
    bash "$TOOLS/snapshot.sh" "$TREE" "$M"
  fi
  printf 'APPLY_ACCEPT_SNAPSHOT_OK %s %s\n' "$M" "$MODE"
done
printf '%s\n' 'FINAL_APPLY_RC 0'
````

### Controller commands

Wait for `$STATE/baseline.done` to end with `ALL_DONE for 77dca5edc8206a1791743ebb3c64512f051c52f3` before attributing the whole-package recordings. `outside-base` measures the skipped-11 cmd race and the two surfacepage legs absent from the controller baseline. It writes only under planner scratch.

````bash
STATE="$HOME/.ailang/state/world-iter449"
S="$STATE/plan-scratch"
bash "$S/gates.sh" "$S/base-tree" "$S/base-extra" M3 outside-base
bash "$S/outside.sh" "$S/st-M1" M1
bash "$S/outside.sh" "$S/st-M2" M2
bash "$S/gates.sh" "$S/st-M1" "$S/outside-gates-M1" M1 outside
bash "$S/gates.sh" "$S/st-M2" "$S/outside-gates-M2" M2 outside
bash "$S/whole-mutants.sh" "$S/outside-mutants-M1" M1
bash "$S/whole-mutants.sh" "$S/outside-mutants-M2" M2
bash "$STATE/ctl_roundtrip.sh" M1 M2
````

For S9c activation after its gates in §2 are met. Cmd and daemon listeners are outside requirements. The Node syntax check and executable taskless review test also run on the scratch stage.

````bash
S="$HOME/.ailang/state/world-iter449/plan-scratch"
bash "$S/outside.sh" "$S/st-M3" M3
bash "$S/gates.sh" "$S/st-M3" "$S/outside-gates-M3" M3 outside
python3 "$S/mutate.py" "$S/outside-named-M3" --ids "$S/named-mutant-ids.txt" --stage M3 --outside
bash "$S/whole-mutants.sh" "$S/outside-mutants-M3" M3
````

The controller owns the Git writes below. The last commit gate records both verifygate and broker with `-timeout 20m` and compares the named fail sets. Cmd and daemon commit gates use named cases. S9b contains no JS files.

````bash
STATE="$HOME/.ailang/state/world-iter449"
S="$STATE/plan-scratch"
printf '%s\n' 'world: pin finish proof scope and local Git boundaries' > "$S/M1.commitmsg"
printf '%s\n' 'world: verify the private finish journal and exact local receipt' > "$S/M2.commitmsg"
"$STATE/commit_ms.sh" M1 "$S/M1.commitmsg" 'export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"; S="$HOME/.ailang/state/world-iter449/plan-scratch"; bash "$S/accept.sh" "$PWD" M1 outside'
"$STATE/commit_ms.sh" M2 "$S/M2.commitmsg" 'export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"; S="$HOME/.ailang/state/world-iter449/plan-scratch"; bash "$S/accept.sh" "$PWD" M2 outside && bash "$S/commit-guards.sh" "$PWD" "$S/commit-guards-M2-$(date +%s)"'
````

### Rehearsal

The exact §7 block, extracted from this plan, printed `FINAL_APPLY_RC 0` after its cumulative fast acceptance and snapshot checks. Raw receipt `/Users/voightkampff/.ailang/state/world-iter449/plan-scratch/rt-apply/final-2/receipt.raw`. See `ROUNDTRIP_OK` for byte-equal cmd, host and docs comparisons for each recorded stage. The fresh archive contains no top-level app directory. The controller's separate round-trip path is available.

## Appendix P — patches

### M1 — Defend proof scope, tool pins and the local Git boundary

````diff M1
diff --git a/cmd/ailang-worldd/operator_finish_auth_test.go b/cmd/ailang-worldd/operator_finish_auth_test.go
index cfa85c7..d5006ea 100644
--- a/cmd/ailang-worldd/operator_finish_auth_test.go
+++ b/cmd/ailang-worldd/operator_finish_auth_test.go
@@ -156,3 +156,41 @@ func TestOperatorFinishAuthInstall(t *testing.T) {
 		t.Fatalf("status %d: %s", rc, out)
 	}
 }
+
+func TestOperatorFinishToolPins(t *testing.T) {
+	f := operatorService(t)
+	f.config.FinishRemotes = []string{fixtureFinishRemote}
+	tools := f.config.ToolDirs[0]
+	trace := filepath.Join(f.home, "unexpected-readiness")
+	for _, name := range []string{"gh", "git"} {
+		writeFinishTestTool(t, &f.config, tools, name, "exit 17\n")
+	}
+	if err := f.env.checkRunnerConfig(f.config); err != nil {
+		t.Fatalf("pinned configuration: %v", err)
+	}
+	for _, path := range []string{filepath.Join(tools, "git"), f.env.paths.Binary, filepath.Join(tools, "gh")} {
+		b, err := os.ReadFile(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if err := os.WriteFile(path, []byte("#!/bin/sh\ntouch "+shellQuote(trace)+"\nexit 1\n"), 0700); err != nil {
+			t.Fatal(err)
+		}
+		if filepath.Base(path) == "gh" {
+			if err := f.env.checkRunnerConfig(f.config); err != nil {
+				t.Fatalf("finish-only gh must leave start available: %v", err)
+			}
+			if err := prepareOperatorFinish(boundedTestContext(t), f.config, fixtureFinishRemote); !errors.Is(err, operatorGitHubAuth(2)) {
+				t.Fatalf("changed readiness executable was accepted: %v", err)
+			}
+		} else if err := f.env.checkRunnerConfig(f.config); err == nil || !strings.Contains(err.Error(), "executable changed") {
+			t.Fatalf("changed shared executable %s: %v", path, err)
+		}
+		if err := os.WriteFile(path, b, 0700); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if _, err := os.Lstat(trace); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("a changed tool executed or sentinel could not be read: %v", err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_preflight.go b/cmd/ailang-worldd/operator_finish_preflight.go
index d06cc9f..96360a6 100644
--- a/cmd/ailang-worldd/operator_finish_preflight.go
+++ b/cmd/ailang-worldd/operator_finish_preflight.go
@@ -4,12 +4,15 @@ import (
 	"bytes"
 	"context"
 	"crypto/sha256"
+	"errors"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/authority"
 	"github.com/sunholo-data/ailang-world/host/store"
 )
 
+var errOperatorFinishSurface = errors.New("finish prep requires surface confirmation")
+
 // operatorFinishPrep binds a current durable claim to independently resolved
 // inputs. Only the proof returned by the runner's recheck can seal it.
 type operatorFinishPrep struct {
@@ -20,9 +23,12 @@ type operatorFinishPrep struct {
 
 func verifyOperatorFinishPrep(ctx context.Context, ro *store.Credentials, frozen approval.OperatorAction, now int64,
 	cas store.OperatorCAS, owner string) (operatorFinishPrep, error) {
-	if frozen.Kind != "finish" || !frozen.Options.ConfirmOnSurface || frozen.ID != cas.ID {
+	if frozen.Kind != "finish" || frozen.ID != cas.ID {
 		return operatorFinishPrep{}, approval.ErrNotVerified
 	}
+	if !frozen.Options.ConfirmOnSurface {
+		return operatorFinishPrep{}, errOperatorFinishSurface
+	}
 	row, err := ro.OperatorAction(ctx, cas.ID, now)
 	if err != nil {
 		return operatorFinishPrep{}, err
diff --git a/cmd/ailang-worldd/operator_finish_preflight_test.go b/cmd/ailang-worldd/operator_finish_preflight_test.go
index 7cd1b0e..ba91035 100644
--- a/cmd/ailang-worldd/operator_finish_preflight_test.go
+++ b/cmd/ailang-worldd/operator_finish_preflight_test.go
@@ -175,3 +175,67 @@ func TestOperatorFinishNoProof(t *testing.T) {
 		t.Fatal("daemon's confirmation proof skipped independent runner recheck")
 	}
 }
+
+func TestOperatorFinishProofScope(t *testing.T) {
+	for _, tc := range []struct {
+		kind, task string
+		want       error
+	}{
+		{"start", "reviewed task", approval.ErrNotVerified},
+		{"finish", "", errOperatorFinishSurface},
+	} {
+		t.Run(tc.kind, func(t *testing.T) {
+			dir, err := filepath.EvalSymlinks(t.TempDir())
+			if err != nil {
+				t.Fatal(err)
+			}
+			path := filepath.Join(dir, "credentials.db")
+			db, err := store.OpenCredentials(path)
+			if err != nil {
+				t.Fatal(err)
+			}
+			defer db.Close()
+			ctx := boundedTestContext(t)
+			const now = 1000
+			key := approvaltest.NewSigner(51)
+			if err := db.EnrolApprover(activationTestContext(t), store.ApproverRow{Name: "fixture", Format: "ssh", PublicKey: key.PublicKey().Blob, EnrolledAt: now - 1}); err != nil {
+				t.Fatal(err)
+			}
+			a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: tc.kind,
+				Episode: "ep1", Store: path, Task: tc.task, Options: approval.OperatorOptions{Preset: "fixture"},
+				Review: "scope fixture", Resolution: []approval.OperatorResolved{{Name: "manifest", Value: "fixture"}}}
+			a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+			req := store.OperatorRequest{Action: a, OperationID: strings.Repeat("2", 32),
+				Challenge: sha256.Sum256([]byte(t.Name())), CreatedAt: now, ExpiresAt: now + 600}
+			row, err := db.RequestOperatorAction(activationTestContext(t), req, now)
+			if err != nil {
+				t.Fatal(err)
+			}
+			cas := operatorPreparationCAS(row)
+			msg, err := approval.OperatorActionMessage(cas.Digest, req.Challenge, req.ExpiresAt)
+			if err != nil {
+				t.Fatal(err)
+			}
+			sig := key.Sign(approval.OperatorActionNamespace, msg)
+			proof, err := authority.JudgeOperatorAction(ctx, db, authority.OperatorJudgeInput{Action: a, Digest: cas.Digest,
+				BytesDigest: sha256.Sum256(row.Action), Challenge: req.Challenge, CreatedAt: now, ExpiresAt: req.ExpiresAt, Now: now + 1}, sig)
+			if err != nil {
+				t.Fatal(err)
+			}
+			if _, err := db.ConfirmOperatorAction(activationTestContext(t), cas, sig, proof, now+1); err != nil {
+				t.Fatal(err)
+			}
+			before, err := db.OperatorAction(ctx, row.ID, now+1)
+			if err != nil || before.State != "confirmed" {
+				t.Fatalf("signed control: %s %v", before.State, err)
+			}
+			digest := serviceFileHash(path)
+			got, err := verifyOperatorFinishPrep(ctx, db, a, now+1, cas, "")
+			seen, readErr := db.OperatorAction(ctx, row.ID, now+1)
+			if !errors.Is(err, tc.want) || got.authority != nil || readErr != nil ||
+				!reflect.DeepEqual(before, seen) || digest != serviceFileHash(path) {
+				t.Fatalf("scope refusal changed evidence: %v, row %s (%v)", err, seen.State, readErr)
+			}
+		})
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_tree.go b/cmd/ailang-worldd/operator_finish_tree.go
index fd9bf62..570f474 100644
--- a/cmd/ailang-worldd/operator_finish_tree.go
+++ b/cmd/ailang-worldd/operator_finish_tree.go
@@ -1,6 +1,7 @@
 package main
 
 import (
+	"bytes"
 	"context"
 	"errors"
 	"fmt"
@@ -65,12 +66,16 @@ func operatorFinishRead(dir, name string, maxBytes int64) ([]byte, os.FileMode,
 	}
 	file := os.NewFile(uintptr(fd), name)
 	defer file.Close()
+	return readOperatorFinishFile(file, filepath.Join(dir, name), maxBytes)
+}
+
+func readOperatorFinishFile(file *os.File, path string, maxBytes int64) ([]byte, os.FileMode, error) {
 	before, err := file.Stat()
 	if err != nil || !before.Mode().IsRegular() || before.Size() > maxBytes {
 		return nil, 0, errOperatorFinishFence
 	}
 	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
-	after, pathErr := os.Lstat(filepath.Join(dir, name))
+	after, pathErr := os.Lstat(path)
 	if err != nil || pathErr != nil || int64(len(data)) != before.Size() || int64(len(data)) > maxBytes ||
 		!os.SameFile(before, after) || before.Mode() != after.Mode() ||
 		before.Size() != after.Size() || before.ModTime() != after.ModTime() {
@@ -151,6 +156,10 @@ func (g operatorFinishGit) config(ctx context.Context, start *store.EpisodeStart
 	if remote != 1 {
 		return "", errOperatorFinishFence
 	}
+	current, _, err := operatorFinishRead(g.common, "config", 64<<10)
+	if err != nil || !bytes.Equal(current, data) {
+		return "", errOperatorFinishFence
+	}
 	return finishSHA(data), nil
 }
 
@@ -194,25 +203,18 @@ func (g operatorFinishGit) snapshot(ctx context.Context, start *store.EpisodeSta
 	if err != nil {
 		return tree, nil, errOperatorFinishFence
 	}
-	names := strings.Split(strings.TrimSuffix(tracked, "\x00"), "\x00")
+	var names []string
+	if tracked != "" {
+		names = strings.Split(strings.TrimSuffix(tracked, "\x00"), "\x00")
+	}
 	root, err := os.OpenRoot(g.dir)
 	if err != nil {
 		return tree, nil, errOperatorFinishFence
 	}
 	defer root.Close()
-	visited := 0
-	err = fs.WalkDir(root.FS(), ".", func(path string, d fs.DirEntry, err error) error {
-		visited++
-		if err != nil || ctx.Err() != nil || visited > 16384 {
-			return errOperatorFinishFence
-		}
-		if path != ".git" && !d.IsDir() && !trackedPaths[path] {
-			names = append(names, path)
-		}
-		return nil
-	})
+	names, err = operatorFinishUntracked(ctx, root.FS(), trackedPaths, names)
 	if err != nil {
-		return tree, nil, errOperatorFinishFence
+		return tree, nil, err
 	}
 	slices.Sort(names)
 	names = slices.Compact(names)
@@ -251,6 +253,24 @@ func (g operatorFinishGit) snapshot(ctx context.Context, start *store.EpisodeSta
 	return tree, contents, nil
 }
 
+func operatorFinishUntracked(ctx context.Context, root fs.FS, trackedPaths map[string]bool, names []string) ([]string, error) {
+	visited := 0
+	err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
+		visited++
+		if err != nil || ctx.Err() != nil || visited > 16384 {
+			return errOperatorFinishFence
+		}
+		if path != ".git" && !d.IsDir() && !trackedPaths[path] {
+			names = append(names, path)
+		}
+		return nil
+	})
+	if err != nil {
+		return nil, errOperatorFinishFence
+	}
+	return names, nil
+}
+
 func (g operatorFinishGit) stage(ctx context.Context, spec store.OperatorLocalCommit, contents [][]byte, env []string) (string, error) {
 	if _, err := g.git(ctx, nil, env, "read-tree", spec.Snapshot.Head); err != nil {
 		return "", err
@@ -288,7 +308,7 @@ func (g operatorFinishGit) planTree(ctx context.Context, spool string, spec stor
 	return g.stage(ctx, spec, contents, env)
 }
 
-func (g operatorFinishGit) inspectCommit(ctx context.Context, spec store.OperatorLocalCommit, tree string) (string, error) {
+func (g operatorFinishGit) inspectCommit(ctx context.Context, spec store.OperatorLocalCommit, tree, agent string) (string, error) {
 	branch, err := g.git(ctx, nil, nil, "symbolic-ref", "HEAD")
 	if err != nil || strings.TrimSpace(branch) != spec.Snapshot.Branch {
 		return "", errOperatorFinishFence
@@ -299,7 +319,7 @@ func (g operatorFinishGit) inspectCommit(ctx context.Context, spec store.Operato
 		return "", errOperatorFinishFence
 	}
 	raw, err := g.git(ctx, nil, nil, "cat-file", "commit", head)
-	if err != nil || !operatorFinishCommitMatches(raw, spec, tree) {
+	if err != nil || !operatorFinishCommitMatches(raw, spec, tree, agent) {
 		return "", errOperatorFinishFence
 	}
 	index, err := g.git(ctx, nil, nil, "diff", "--cached", "--name-only", "--no-ext-diff", "--no-textconv", head, "--")
@@ -309,13 +329,17 @@ func (g operatorFinishGit) inspectCommit(ctx context.Context, spec store.Operato
 	return head, nil
 }
 
-func operatorFinishCommitMatches(raw string, spec store.OperatorLocalCommit, tree string) bool {
+func operatorFinishCommitMatches(raw string, spec store.OperatorLocalCommit, tree, agent string) bool {
 	header, message, found := strings.Cut(raw, "\n\n")
 	lines := strings.Split(header, "\n")
 	if !found || len(lines) != 4 || lines[0] != "tree "+tree ||
 		lines[1] != "parent "+spec.Snapshot.Head || message != spec.Message {
 		return false
 	}
+	author, authorTime, ok := strings.Cut(lines[2], "> ")
+	if !ok || author+">" != "author "+agent+" via AILANG World <ailang-world@invalid>" || len(strings.Fields(authorTime)) != 2 {
+		return false
+	}
 	committer, timeFields, ok := strings.Cut(lines[3], "> ")
 	return ok && committer+">" == "committer "+spec.Committer && len(strings.Fields(timeFields)) == 2
 }
diff --git a/cmd/ailang-worldd/operator_finish_tree_test.go b/cmd/ailang-worldd/operator_finish_tree_test.go
index dbf27d1..7927390 100644
--- a/cmd/ailang-worldd/operator_finish_tree_test.go
+++ b/cmd/ailang-worldd/operator_finish_tree_test.go
@@ -1,12 +1,18 @@
 package main
 
 import (
+	"bytes"
 	"errors"
+	"fmt"
+	"io"
 	"os"
 	"path/filepath"
 	"reflect"
+	"slices"
 	"strings"
 	"testing"
+	"testing/fstest"
+	"time"
 
 	"github.com/sunholo-data/ailang-world/host/store"
 	"golang.org/x/sys/unix"
@@ -25,7 +31,7 @@ func operatorFinishGitFixture(t *testing.T, f *finishFixture) operatorFinishGit
 func TestOperatorFinishWorktreeFence(t *testing.T) {
 	f := newFinishFixture(t)
 	g := operatorFinishGitFixture(t, f)
-	ctx := boundedTestContext(t)
+	ctx := activationTestContext(t)
 	path := filepath.Join(g.common, "config")
 	config := readStartFixtureFile(t, path)
 	if _, _, err := g.snapshot(ctx, f.rep.Start); err != nil {
@@ -64,7 +70,7 @@ func TestOperatorFinishWorktreeFence(t *testing.T) {
 func TestOperatorCommitIdentity(t *testing.T) {
 	f := newFinishFixture(t)
 	g := operatorFinishGitFixture(t, f)
-	ctx := boundedTestContext(t)
+	ctx := activationTestContext(t)
 	initial, contents, err := g.snapshot(ctx, f.rep.Start)
 	if err != nil {
 		t.Fatal(err)
@@ -90,13 +96,13 @@ func TestOperatorCommitIdentity(t *testing.T) {
 	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) {
 		t.Fatalf("repo-controlled hook executed or sentinel unreadable: %v", err)
 	}
-	head, err := g.inspectCommit(ctx, spec, tree)
+	head, err := g.inspectCommit(ctx, spec, tree, f.rep.Start.Agent)
 	if err != nil || head == initial.Head {
 		t.Fatalf("exact descendant missing: %s %v", head, err)
 	}
 	raw := f.gitOutput("-C", f.wt, "cat-file", "commit", head)
-	for _, field := range []string{"parent", "tree", "message", "trailer", "committer"} {
-		bad, want := spec, tree
+	for _, field := range []string{"parent", "tree", "message", "trailer", "committer", "author-name", "author-email", "author-time"} {
+		bad, want, candidate := spec, tree, raw
 		switch field {
 		case "parent":
 			bad.Snapshot.Head = strings.Repeat("d", 40)
@@ -108,14 +114,227 @@ func TestOperatorCommitIdentity(t *testing.T) {
 			bad.Message = strings.Replace(spec.Message, "episode=ep1", "episode=other", 1)
 		case "committer":
 			bad.Committer = "Someone Else <different@invalid>"
+		case "author-name":
+			candidate = strings.Replace(raw, "author codex via AILANG World ", "author other via AILANG World ", 1)
+		case "author-email":
+			candidate = strings.Replace(raw, "ailang-world@invalid", "other@invalid", 1)
+		case "author-time":
+			headers := strings.SplitN(raw, "\n", 4)
+			identity, _, _ := strings.Cut(headers[2], "> ")
+			headers[2] = identity + "> 1"
+			candidate = strings.Join(headers, "\n")
 		}
-		if operatorFinishCommitMatches(raw, bad, want) {
+		if (strings.HasPrefix(field, "author") && candidate == raw) || operatorFinishCommitMatches(candidate, bad, want, f.rep.Start.Agent) {
 			t.Fatalf("%s mismatch observed", field)
 		}
 	}
+	f.gitOutput("-C", f.wt, "branch", "other-branch", head)
+	f.gitOutput("-C", f.wt, "symbolic-ref", "HEAD", "refs/heads/other-branch")
+	if _, err := g.inspectCommit(ctx, spec, tree, f.rep.Start.Agent); !errors.Is(err, errOperatorFinishFence) {
+		t.Fatalf("adopted a commit on the wrong branch: %v", err)
+	}
+	f.gitOutput("-C", f.wt, "symbolic-ref", "HEAD", initial.Branch)
 	f.write("extra", "partial stage")
 	f.gitOutput("-C", f.wt, "add", "--", "extra")
-	if _, err := g.inspectCommit(ctx, spec, tree); err == nil || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head {
+	if _, err := g.inspectCommit(ctx, spec, tree, f.rep.Start.Agent); err == nil || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head {
 		t.Fatal("uncertain partial stage adopted or rewrote HEAD")
 	}
 }
+
+func TestOperatorFinishCaptureBounds(t *testing.T) {
+	f := newFinishFixture(t)
+	g := operatorFinishGitFixture(t, f)
+	ctx := activationTestContext(t)
+	capture := func() store.OperatorFinishTree {
+		t.Helper()
+		snap, _, err := g.snapshot(ctx, f.rep.Start)
+		if err != nil {
+			t.Fatal(err)
+		}
+		return snap
+	}
+	refuse := func() {
+		t.Helper()
+		before := serviceFileHash(filepath.Join(g.admin, "index"))
+		if got, _, err := g.snapshot(ctx, f.rep.Start); !errors.Is(err, errOperatorFinishFence) || len(got.Paths) != 0 ||
+			serviceFileHash(filepath.Join(g.admin, "index")) != before || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
+			t.Fatalf("capture fence allowed a changed Git state: %v", err)
+		}
+	}
+	first := capture()
+	f.write("b.go", "edit\n")
+	second := capture()
+	if reflect.DeepEqual(first, second) || !slices.Contains(second.Paths, store.OperatorFinishPath{Path: "b.go", Mode: "100644", Digest: finishSHA([]byte("edit\n"))}) {
+		t.Fatal("untracked same-size edit did not change the content digest")
+	}
+	for _, path := range []string{filepath.Join(g.admin, "config.worktree"), filepath.Join(g.common, "objects", "info", "alternates")} {
+		seedStartFixtureFile(t, path, []byte("\n"))
+		refuse()
+		if err := os.Remove(path); err != nil {
+			t.Fatal(err)
+		}
+	}
+	f.gitOutput("-C", f.wt, "add", "--", "b.go")
+	refuse()
+	f.gitOutput("-C", f.wt, "read-tree", first.Head)
+
+	extra := filepath.Join(f.wt, "bounds")
+	if err := os.Mkdir(extra, 0700); err != nil {
+		t.Fatal(err)
+	}
+	for n := range 254 {
+		seedStartFixtureFile(t, filepath.Join(extra, fmt.Sprint(n)), nil)
+	}
+	if got := len(capture().Paths); got != 257 {
+		t.Fatalf("exact file cap: %d", got)
+	}
+	seedStartFixtureFile(t, filepath.Join(extra, "overflow"), nil)
+	refuse()
+	if err := os.RemoveAll(extra); err != nil {
+		t.Fatal(err)
+	}
+	for _, name := range []string{"a.go", "b.go", "gen.txt"} {
+		if err := os.Remove(filepath.Join(f.wt, name)); err != nil {
+			t.Fatal(err)
+		}
+	}
+	for _, p := range []struct {
+		name string
+		size int64
+	}{{"large", 16 << 20}, {"overflow", 1}} {
+		path := filepath.Join(f.wt, p.name)
+		seedStartFixtureFile(t, path, nil)
+		if err := os.Truncate(path, p.size); err != nil {
+			t.Fatal(err)
+		}
+		if p.name == "large" {
+			capture()
+		}
+	}
+	refuse()
+	files := fstest.MapFS{}
+	tracked := map[string]bool{}
+	for n := range 16383 {
+		name := fmt.Sprintf("fixture-%05d", n)
+		files[name] = &fstest.MapFile{}
+		tracked[name] = true
+	}
+	if names, err := operatorFinishUntracked(ctx, files, tracked, nil); err != nil || len(names) != 0 {
+		t.Fatalf("exact traversal bound including root: %v", err)
+	}
+	files["overflow"] = &fstest.MapFile{}
+	if names, err := operatorFinishUntracked(ctx, files, tracked, nil); !errors.Is(err, errOperatorFinishFence) || len(names) != 0 {
+		t.Fatalf("traversal beyond the bound: %v", err)
+	}
+}
+
+func TestOperatorFinishReadReplacement(t *testing.T) {
+	dir := t.TempDir()
+	path := filepath.Join(dir, "checked")
+	replacement := filepath.Join(dir, "replacement")
+	for _, p := range []string{path, replacement} {
+		seedStartFixtureFile(t, p, []byte("same"))
+		stamp := time.Unix(1000, 0)
+		if err := os.Chtimes(p, stamp, stamp); err != nil {
+			t.Fatal(err)
+		}
+	}
+	f, err := os.Open(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer f.Close()
+	if b, _, err := readOperatorFinishFile(f, path, 4); err != nil || string(b) != "same" {
+		t.Fatalf("held file control: %q %v", b, err)
+	}
+	if _, err := f.Seek(0, io.SeekStart); err != nil {
+		t.Fatal(err)
+	}
+	if err := os.Rename(replacement, path); err != nil {
+		t.Fatal(err)
+	}
+	if b, _, err := readOperatorFinishFile(f, path, 4); !errors.Is(err, errOperatorFinishFence) || b != nil {
+		t.Fatalf("same bytes and different inode escaped fence: %q %v", b, err)
+	}
+}
+
+func TestOperatorCommitReservationFence(t *testing.T) {
+	f := newFinishFixture(t)
+	g := operatorFinishGitFixture(t, f)
+	ctx := activationTestContext(t)
+	initial, contents, err := g.snapshot(ctx, f.rep.Start)
+	if err != nil {
+		t.Fatal(err)
+	}
+	spec := store.OperatorLocalCommit{Snapshot: initial, Message: "fixture\n", Committer: "Fixture Operator <fixture@example.invalid>"}
+	// Leave the valid object in this temp repo so an unchecked mismatch would commit it.
+	tree, err := g.stage(ctx, spec, contents, []string{"GIT_INDEX_FILE=" + filepath.Join(f.root, "plan-index")})
+	if err != nil {
+		t.Fatal(err)
+	}
+	changed := slices.Clone(contents)
+	changed[0] = []byte("different tree")
+	if err := g.commit(ctx, f.root, spec, changed, tree, f.rep.Start); !errors.Is(err, errOperatorCommitNoRef) {
+		t.Fatalf("candidate tree mismatch: %v", err)
+	}
+	wrap := func(g operatorFinishGit, script string) operatorFinishGit {
+		t.Helper()
+		path := filepath.Join(t.TempDir(), "git")
+		if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+script+"\nexec "+shellQuote(f.git)+" \"$@\"\n"), 0700); err != nil {
+			t.Fatal(err)
+		}
+		g.bin = path
+		return g
+	}
+	unexpectedPath := filepath.Join(f.wt, "changed-after-commit-tree")
+	afterObject := wrap(g, "for x; do if [ \"$x\" = commit-tree ]; then touch "+shellQuote(unexpectedPath)+"; fi; done")
+	if err := afterObject.commit(ctx, f.root, spec, contents, tree, f.rep.Start); !errors.Is(err, errOperatorFinishReview) {
+		t.Fatalf("post-object drift ignored: %v", err)
+	}
+	if err := os.Remove(unexpectedPath); err != nil {
+		t.Fatal(err)
+	}
+	if seen, _, err := g.snapshot(ctx, f.rep.Start); err != nil || !reflect.DeepEqual(initial, seen) {
+		t.Fatalf("refusal changed Git state: %v", err)
+	}
+	// Advance the pinned branch at the update-ref seam, after the final snapshot.
+	foreign := strings.TrimSpace(f.gitOutput("-C", f.wt, "commit-tree", tree, "-p", initial.Head, "-m", "foreign"))
+	refArgs := " --git-dir=" + shellQuote(g.admin) + " update-ref " + shellQuote(initial.Branch) + " " + shellQuote(foreign)
+	race := wrap(g, "for x; do if [ \"$x\" = update-ref ]; then "+shellQuote(f.git)+refArgs+"; fi; done")
+	if err := race.commit(ctx, f.root, spec, contents, tree, f.rep.Start); err == nil ||
+		strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != foreign {
+		t.Fatalf("old-value CAS overwrote a different HEAD: %v", err)
+	}
+	if _, err := g.inspectCommit(ctx, spec, tree, f.rep.Start.Agent); !errors.Is(err, errOperatorFinishFence) {
+		t.Fatalf("recovery adopted the competing commit: %v", err)
+	}
+}
+
+func TestOperatorFinishConfigReplacement(t *testing.T) {
+	f := newFinishFixture(t)
+	g := operatorFinishGitFixture(t, f)
+	ctx := activationTestContext(t)
+	path := filepath.Join(g.common, "config")
+	if _, err := g.config(ctx, f.rep.Start); err != nil {
+		t.Fatal(err)
+	}
+	swapped := filepath.Join(f.root, "replacement")
+	body := readStartFixtureFile(t, path)
+	seedStartFixtureFile(t, swapped, append(body, []byte("\n[core]\n\tfsmonitor = /no-fixture-program\n")...))
+	sentinel := filepath.Join(f.root, "outside-config-parser")
+	fake := filepath.Join(f.root, "git-config-race")
+	if err := os.WriteFile(fake, []byte("#!/bin/sh\nset -eu\n"+
+		"if [ \"$1\" = --git-dir=/dev/null ]; then mv "+shellQuote(swapped)+" "+shellQuote(path)+
+		"; else touch "+shellQuote(sentinel)+"; exit 93; fi\nexec "+shellQuote(f.git)+" \"$@\"\n"), 0700); err != nil {
+		t.Fatal(err)
+	}
+	g.bin = fake
+	before := serviceFileHash(filepath.Join(g.admin, "index"))
+	if _, _, err := g.snapshot(ctx, f.rep.Start); !errors.Is(err, errOperatorFinishFence) {
+		t.Fatalf("replacement was accepted: %v", err)
+	}
+	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) ||
+		before != serviceFileHash(filepath.Join(g.admin, "index")) || bytes.Equal(body, readStartFixtureFile(t, path)) {
+		t.Fatalf("replacement did not stop before using the live Git configuration: %v", err)
+	}
+}
````

### M2 — Exercise the private check and commit journal

````diff M2
diff --git a/cmd/ailang-worldd/operator_finish_local.go b/cmd/ailang-worldd/operator_finish_local.go
new file mode 100644
index 0000000..cad8b54
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_local.go
@@ -0,0 +1,230 @@
+package main
+
+import (
+	"context"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"reflect"
+	"slices"
+	"strings"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/checkrun"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+var (
+	errOperatorFinishTerminal  = errors.New("unsandboxed check requires terminal finish")
+	errOperatorFinishUncertain = errors.New("finish outcome uncertain; inspect durable receipts")
+	errOperatorFinishCheck     = errors.New("sandboxed check failed; local commit withheld")
+	errOperatorCommitFailed    = errors.New("local commit failed; retained the check receipts")
+)
+
+type operatorFinishWork struct {
+	prep       operatorFinishPrep
+	generation int64
+	manifest   store.OperatorFinishReview
+	git        operatorFinishGit
+	report     daemon.EpisodeReport
+	spool      string
+	read       *store.Credentials
+	fresh      func(context.Context) error
+	now        func() int64
+	budget     func(time.Duration) (context.Context, context.CancelFunc)
+	run        func(context.Context, store.StartCheck) (checkrun.Result, error)
+	record     func(context.Context, demoObject) (int64, error)
+	evidence   func(context.Context, store.OperatorFinishCheck) bool
+	commit     func(context.Context, store.OperatorLocalCommit, [][]byte, string) error
+}
+
+func (w operatorFinishWork) step(kind, identity string) store.OperatorStepKey {
+	r := w.prep.row
+	return store.OperatorStepKey{Action: operatorPreparationCAS(r), Owner: r.Owner, Kind: kind, Identity: r.ID + ":" + identity, Version: 1}
+}
+
+func (w operatorFinishWork) reserve(k store.OperatorStepKey) error {
+	ctx, cancel := w.budget(time.Second)
+	defer cancel()
+	return runnerWrite(w.manifest.Store, func(db *store.Credentials) error {
+		return db.ReserveOperatorStep(ctx, k, w.now())
+	})
+}
+
+func (w operatorFinishWork) observe(k store.OperatorStepKey, receipt any) error {
+	raw, err := json.Marshal(receipt)
+	if err != nil {
+		return err
+	}
+	ctx, cancel := w.budget(time.Second)
+	defer cancel()
+	return runnerWrite(w.manifest.Store, func(db *store.Credentials) error {
+		return db.ObserveOperatorStep(ctx, k, raw)
+	})
+}
+
+func (w operatorFinishWork) inputs(ctx context.Context) ([][]byte, error) {
+	snap, b, err := w.git.snapshot(ctx, w.report.Start)
+	if err != nil || !reflect.DeepEqual(snap, w.manifest.Commit.Snapshot) {
+		return nil, errOperatorFinishReview
+	}
+	return b, nil
+}
+
+// runEpisodeFinishWithPrep is the private entry beside terminal finish. Its
+// caller has the episode and runner locks, and has completed Git auth before
+// the claim. All effects use the exact inputs in the locally verified claim.
+func runEpisodeFinishWithPrep(ctx context.Context, w operatorFinishWork) error {
+	if !w.prep.claimed() {
+		return approval.ErrNotVerified
+	}
+	m := w.manifest
+	if !w.boundManifest(w.generation) {
+		return errOperatorFinishReview
+	}
+	for _, check := range m.Checks {
+		if check.Unsandboxed {
+			return errOperatorFinishTerminal
+		}
+	}
+	if m.Episode != w.prep.action.Episode || m.Store != w.prep.action.Store || w.report.Start == nil ||
+		len(m.Unattributed) != 0 || len(m.Protected) != 0 || len(m.Checks) > 31 || len(m.Commit.Snapshot.Paths) == 0 {
+		return errOperatorFinishFence
+	}
+	j, err := w.read.ReadOperatorFinishJournal(ctx, operatorPreparationCAS(w.prep.row))
+	if err != nil {
+		return err
+	}
+	if j.CommitPhase == "reserved" {
+		return w.reconcileCommit(j.CommitKey)
+	}
+	if j.CommitPhase == "observed" {
+		if j.Commit.State == store.OperatorLocalCommitDone {
+			return nil
+		}
+		return errOperatorCommitFailed
+	}
+	if j.OpenCheck || len(j.Checks) > len(m.Checks) {
+		return errOperatorFinishUncertain
+	}
+	for i, old := range j.Checks {
+		if old.Name != m.Checks[i].Name || !w.evidence(ctx, old) {
+			return errOperatorFinishUncertain
+		}
+		if old.Result != "passed" {
+			return errOperatorFinishCheck
+		}
+	}
+	contents, err := w.inputs(ctx)
+	if err != nil {
+		return err
+	}
+	for i, check := range m.Checks[len(j.Checks):] {
+		if err := w.fresh(ctx); err != nil {
+			return err
+		}
+		if err := w.check(ctx, check, i+len(j.Checks)); err != nil {
+			return err
+		}
+		contents, err = w.inputs(ctx)
+		if err != nil {
+			return err
+		}
+	}
+	if err := w.fresh(ctx); err != nil {
+		return err
+	}
+	tree, err := w.git.planTree(ctx, w.spool, m.Commit, contents)
+	if err != nil {
+		return err
+	}
+	contents, err = w.inputs(ctx)
+	if err != nil {
+		return err
+	}
+	if err := w.fresh(ctx); err != nil {
+		return err
+	}
+	key := w.step("commit", tree)
+	if err := w.reserve(key); err != nil {
+		return err
+	}
+	commitErr := w.commit(ctx, m.Commit, contents, tree)
+	if err := w.reconcileCommit(key); err == nil {
+		return nil
+	}
+	if errors.Is(commitErr, errOperatorCommitNoRef) {
+		budget, cancel := w.budget(3 * time.Second)
+		defer cancel()
+		if _, err := w.inputs(budget); err == nil {
+			if err := w.observe(key, store.OperatorFinishLocalCommit{State: "commit-failed"}); err == nil {
+				return errOperatorCommitFailed
+			}
+		}
+	}
+	return errOperatorFinishUncertain
+}
+
+func (w operatorFinishWork) check(ctx context.Context, check store.StartCheck, index int) error {
+	key := w.step("checks", fmt.Sprintf("%02d", index))
+	if err := w.reserve(key); err != nil {
+		return err
+	}
+	runCtx, stop := operatorChildBudget(ctx, time.Duration(check.TimeoutS*float64(time.Second)))
+	result, err := w.run(runCtx, check)
+	stop()
+	if err != nil {
+		// No trustworthy terminal readback. Retain the reservation and let
+		// an operator inspect evidence. A later runner must not retry it.
+		return fmt.Errorf("%w: sandboxed check interrupted", errOperatorFinishUncertain)
+	}
+	passed := result.ExitCode == 0
+	outcome := finishOutcome(w.report, check, result, passed)
+	raw, _ := json.Marshal(outcome)
+	canonical, _, err := parseTaskOutcome(raw)
+	if err != nil {
+		return errOperatorFinishUncertain
+	}
+	obj := taskOutcomeObject(canonical)
+	writeCtx, cancel := w.budget(3 * time.Second)
+	defer cancel()
+	entry, err := w.record(writeCtx, obj)
+	if err != nil || entry < 1 {
+		return errOperatorFinishUncertain
+	}
+	class := "passed"
+	if !passed {
+		class = "failed"
+	}
+	if err := w.observe(key, store.OperatorFinishCheck{Name: check.Name, Result: class, Ref: obj.Hash, Entry: entry}); err != nil {
+		return err
+	}
+	if !passed {
+		return errOperatorFinishCheck
+	}
+	return nil
+}
+
+func (w operatorFinishWork) reconcileCommit(key store.OperatorStepKey) error {
+	tree, ok := strings.CutPrefix(key.Identity, w.prep.row.ID+":")
+	if !ok || !operatorFinishSHA.MatchString(tree) {
+		return errOperatorFinishUncertain
+	}
+	ctx, cancel := w.budget(3 * time.Second)
+	defer cancel()
+	if _, err := w.git.config(ctx, w.report.Start); err != nil {
+		return errOperatorFinishUncertain
+	}
+	sha, err := w.git.inspectCommit(ctx, w.manifest.Commit, tree, w.report.Start.Agent)
+	if err != nil {
+		return errOperatorFinishUncertain
+	}
+	return w.observe(key, store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: sha})
+}
+
+func (w operatorFinishWork) boundManifest(generation int64) bool {
+	a, err := w.manifest.Action(generation, w.prep.action.Options.Ready)
+	return err == nil && slices.Equal(a.Resolution, w.prep.action.Resolution) && a.Review == w.prep.action.Review
+}
diff --git a/cmd/ailang-worldd/operator_finish_local_test.go b/cmd/ailang-worldd/operator_finish_local_test.go
new file mode 100644
index 0000000..e934a52
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_local_test.go
@@ -0,0 +1,408 @@
+package main
+
+import (
+	"context"
+	"crypto/sha256"
+	"errors"
+	"io"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/checkrun"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func operatorLocalFinishFixture(t *testing.T, change func(*store.OperatorFinishReview), confirmed ...bool) (*finishFixture, operatorFinishWork, *int64) {
+	t.Helper()
+	f := newFinishFixture(t)
+	g, snap := operatorFinishFixtureSnapshot(t, f)
+	return operatorLocalFinishOn(t, f, g, snap, change, confirmed...)
+}
+
+func operatorFinishFixtureSnapshot(t *testing.T, f *finishFixture) (operatorFinishGit, store.OperatorFinishTree) {
+	t.Helper()
+	g := operatorFinishGitFixture(t, f)
+	snap, _, err := g.snapshot(activationTestContext(t), f.rep.Start)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return g, snap
+}
+
+func operatorLocalFinishOn(t *testing.T, f *finishFixture, g operatorFinishGit, snap store.OperatorFinishTree, change func(*store.OperatorFinishReview), confirmed ...bool) (*finishFixture, operatorFinishWork, *int64) {
+	t.Helper()
+	ctx := activationTestContext(t)
+	path, err := store.CanonicalPath(f.health.CredentialsDB)
+	if err != nil {
+		t.Fatal(err)
+	}
+	db, err := store.OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer db.Close()
+	key := approvaltest.NewSigner(51)
+	if err := db.EnrolApprover(activationTestContext(t), store.ApproverRow{Name: "fixture", Format: "ssh", PublicKey: key.PublicKey().Blob, EnrolledAt: 900}); err != nil {
+		t.Fatal(err)
+	}
+	m := store.OperatorFinishReview{Store: path, Episode: "ep1", Preset: "fixture", Worktree: g.dir,
+		Remote: f.rep.Start.Remote, Checks: []store.StartCheck{{Name: "short", Argv: []string{"echo", "fixture"}, Exec: "echo", Cwd: ".", Sandbox: "sandboxed", TimeoutS: 15}},
+		Commit: store.OperatorLocalCommit{Snapshot: snap, Message: "Fixture finish\n\nWorld-Provenance: episode=ep1\n", Committer: "Fixture Operator <fixture@example.invalid>"}}
+	if change != nil {
+		change(&m)
+	}
+	a, err := m.Action(1, false)
+	if err != nil {
+		t.Fatal(err)
+	}
+	a.ID = fixtureActionID
+	req := store.OperatorRequest{Action: a, OperationID: strings.Repeat("1", 32), Challenge: sha256.Sum256([]byte(t.Name())), CreatedAt: 1000, ExpiresAt: 1600}
+	row, err := db.RequestOperatorAction(activationTestContext(t), req, 1000)
+	if err != nil {
+		t.Fatal(err)
+	}
+	cas := operatorPreparationCAS(row)
+	msg, err := approval.OperatorActionMessage(cas.Digest, req.Challenge, req.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	sig := key.Sign(approval.OperatorActionNamespace, msg)
+	in := authority.OperatorJudgeInput{Action: a, Digest: cas.Digest, BytesDigest: sha256.Sum256(row.Action),
+		Challenge: req.Challenge, CreatedAt: 1000, ExpiresAt: 1600, Now: 1001}
+	proof, err := authority.JudgeOperatorAction(ctx, db, in, sig)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if _, err := db.ConfirmOperatorAction(activationTestContext(t), cas, sig, proof, 1001); err != nil {
+		t.Fatal(err)
+	}
+	h := store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "finish-runner"}, Protocol: store.OperatorRunnerProtocol, Build: "fixture", Readiness: "both", SeenAt: 1001}
+	h.Owner.Generation, err = db.PutOperatorHeartbeat(activationTestContext(t), h)
+	if err != nil {
+		t.Fatal(err)
+	}
+	owner := ""
+	if len(confirmed) == 0 {
+		if err := db.ClaimOperatorAction(activationTestContext(t), cas, h.Owner, 1001); err != nil {
+			t.Fatal(err)
+		}
+		cas.Generation++
+		owner = h.Owner.Nonce
+	}
+	ro, err := store.OpenCredentialsReadOnly(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { ro.Close() })
+	now := int64(1001)
+	p, err := verifyOperatorFinishPrep(ctx, ro, a, now, cas, owner)
+	if err != nil {
+		t.Fatal(err)
+	}
+	w := operatorFinishWork{prep: p, generation: 1, manifest: m, report: f.rep, git: g, spool: f.root,
+		read: ro, now: func() int64 { return now }, budget: (serviceEnv{operatorBudget: activationTestBudget}).operatorContext}
+	w.fresh = func(ctx context.Context) error {
+		_, err := verifyOperatorFinishPrep(ctx, ro, a, now, cas, owner)
+		return err
+	}
+	w.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
+		if j := operatorFinishJournal(t, w); !j.OpenCheck || j.CommitPhase != "" {
+			t.Fatalf("check ran without its reservation: %+v", j)
+		}
+		f.checks++
+		return checkrun.Result{ExitCode: 0}, nil
+	}
+	w.record = func(context.Context, demoObject) (int64, error) {
+		f.outcomes++
+		return int64(f.outcomes), nil
+	}
+	w.evidence = func(context.Context, store.OperatorFinishCheck) bool { return f.outcomes > 0 }
+	w.commit = func(ctx context.Context, spec store.OperatorLocalCommit, b [][]byte, tree string) error {
+		return g.commit(ctx, f.root, spec, b, tree, f.rep.Start)
+	}
+	t.Cleanup(func() {
+		row, err := ro.OperatorAction(activationTestContext(t), p.row.ID, now)
+		if err != nil || row.V3ID != "" || len(f.seam.calls) != 0 || len(f.scopes) != 0 || len(f.finishRecords) != 0 || len(f.requests) != 0 {
+			t.Errorf("private preparation leaked into dispatch: V3=%q, %v, %v", row.V3ID, f.seam.calls, err)
+		}
+	})
+	return f, w, &now
+}
+
+func operatorFinishJournal(t *testing.T, w operatorFinishWork) store.OperatorFinishJournal {
+	t.Helper()
+	j, err := w.read.ReadOperatorFinishJournal(activationTestContext(t), operatorPreparationCAS(w.prep.row))
+	if err != nil {
+		t.Fatal(err)
+	}
+	return j
+}
+
+func TestOperatorFinishPrep(t *testing.T) {
+	f, w, _ := operatorLocalFinishFixture(t, nil)
+	baseline := operatorFinishJournal(t, w)
+	unheld := w
+	unheld.prep.row.Owner = ""
+	if err := runEpisodeFinishWithPrep(activationTestContext(t), unheld); !errors.Is(err, approval.ErrNotVerified) ||
+		!reflect.DeepEqual(baseline, operatorFinishJournal(t, w)) || f.checks != 0 {
+		t.Fatalf("unclaimed proof crossed the entry: %v", err)
+	}
+	original := w.manifest.Commit.Message
+	w.manifest.Commit.Message = "different local commit"
+	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorFinishReview) ||
+		!reflect.DeepEqual(baseline, operatorFinishJournal(t, w)) || f.checks != 0 {
+		t.Fatalf("other signed intent was admitted: %v", err)
+	}
+	w.manifest.Commit.Message = original
+	if !w.boundManifest(w.generation) {
+		t.Fatal("signed manifest must bind original inputs")
+	}
+	for _, kind := range []string{"protected", "unattributed"} {
+		t.Run(kind, func(t *testing.T) {
+			rig := *f
+			rig.t = t
+			rig.health.CredentialsDB = filepath.Join(t.TempDir(), "credentials.db")
+			_, blocked, _ := operatorLocalFinishOn(t, &rig, w.git, w.manifest.Commit.Snapshot, func(m *store.OperatorFinishReview) {
+				if kind == "protected" {
+					m.Protected = []string{"a.go"}
+				} else {
+					m.Unattributed = []string{"b.go"}
+				}
+			})
+			before := operatorFinishJournal(t, blocked)
+			if err := runEpisodeFinishWithPrep(activationTestContext(t), blocked); !errors.Is(err, errOperatorFinishFence) ||
+				!reflect.DeepEqual(before, operatorFinishJournal(t, blocked)) || rig.checks+rig.outcomes != 0 {
+				t.Fatalf("unsupported %s change reached execution: %v", kind, err)
+			}
+		})
+	}
+}
+
+func TestOperatorCommitReconcile(t *testing.T) {
+	f, w, _ := operatorLocalFinishFixture(t, nil)
+	sentinel := filepath.Join(f.root, "unexpected-hook")
+	for _, name := range []string{"pre-commit", "reference-transaction"} {
+		hook := filepath.Join(w.git.common, "hooks", name)
+		if err := os.WriteFile(hook, []byte("#!/bin/sh\ntouch "+shellQuote(sentinel)+"\n"), 0700); err != nil {
+			t.Fatal(err)
+		}
+	}
+	original := w.commit
+	w.commit = func(ctx context.Context, spec store.OperatorLocalCommit, b [][]byte, tree string) error {
+		j := operatorFinishJournal(t, w)
+		if j.CommitPhase != "reserved" || !strings.HasSuffix(j.CommitKey.Identity, ":"+tree) || len(j.Checks) != 1 || j.Checks[0].Result != "passed" {
+			t.Fatalf("commit without exact reserved intent: %+v", j)
+		}
+		if err := original(ctx, spec, b, tree); err != nil {
+			t.Fatalf("fixture commit did not reach its lost-reply boundary: %v", err)
+		}
+		return io.ErrUnexpectedEOF // bounded commit advanced HEAD, reply was lost
+	}
+	ctx := activationTestContext(t)
+	if err := runEpisodeFinishWithPrep(ctx, w); err != nil {
+		t.Fatal(err)
+	}
+	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("a repository hook executed or sentinel unreadable: %v", err)
+	}
+	seen := operatorFinishJournal(t, w)
+	head := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
+	if head == f.before || seen.Commit != (store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: head}) {
+		t.Fatalf("hand-off %+v, head %s", seen.Commit, head)
+	}
+	w.commit = func(context.Context, store.OperatorLocalCommit, [][]byte, string) error {
+		t.Fatal("repeated an observed commit")
+		return nil
+	}
+	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); err != nil || f.checks != 1 || f.outcomes != 1 {
+		t.Fatalf("resume: %v, checks %d, outcomes %d", err, f.checks, f.outcomes)
+	}
+}
+
+func TestOperatorCommitFailed(t *testing.T) {
+	f, w, _ := operatorLocalFinishFixture(t, nil)
+	n := 0
+	w.commit = func(context.Context, store.OperatorLocalCommit, [][]byte, string) error {
+		n++
+		if got := operatorFinishJournal(t, w); got.CommitPhase != "reserved" {
+			t.Fatal("commit attempt preceded reservation")
+		}
+		return errOperatorCommitNoRef
+	}
+	for range 2 {
+		if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorCommitFailed) {
+			t.Fatal(err)
+		}
+	}
+	j := operatorFinishJournal(t, w)
+	if n != 1 || len(j.Checks) != 1 || j.Commit.State != "commit-failed" || j.Commit.SHA != "" ||
+		strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
+		t.Fatalf("failure lost evidence or retried: %+v attempts=%d", j, n)
+	}
+}
+
+func TestOperatorFinishChecks(t *testing.T) {
+	shared := newFinishFixture(t)
+	git, snap := operatorFinishFixtureSnapshot(t, shared)
+	for _, name := range []string{"failure", "abort", "timeout", "unknown-outcome", "changed-input"} {
+		t.Run(name, func(t *testing.T) {
+			rig := *shared
+			rig.t = t
+			rig.health.CredentialsDB = filepath.Join(t.TempDir(), "credentials.db")
+			f, w, _ := operatorLocalFinishOn(t, &rig, git, snap, func(m *store.OperatorFinishReview) {
+				if name == "changed-input" {
+					next := m.Checks[0]
+					next.Name = "after-drift"
+					m.Checks = append(m.Checks, next)
+				}
+			})
+			w.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
+				f.checks++
+				switch name {
+				case "failure":
+					return checkrun.Result{ExitCode: 1}, nil
+				case "abort":
+					return checkrun.Result{}, context.Canceled
+				case "timeout":
+					return checkrun.Result{}, context.DeadlineExceeded
+				default:
+					f.write("b.go", "different bytes")
+					return checkrun.Result{ExitCode: 0}, nil
+				}
+			}
+			if name == "unknown-outcome" {
+				w.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
+					f.checks++
+					return checkrun.Result{ExitCode: 0}, nil
+				}
+				w.record = func(context.Context, demoObject) (int64, error) { return 0, nil }
+			}
+			want := map[string]error{"unknown-outcome": errOperatorFinishUncertain, "failure": errOperatorFinishCheck, "abort": errOperatorFinishUncertain,
+				"timeout": errOperatorFinishUncertain, "changed-input": errOperatorFinishReview}[name]
+			for range 2 {
+				if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, want) {
+					t.Fatal(err)
+				}
+			}
+			j := operatorFinishJournal(t, w)
+			lost := name == "abort" || name == "timeout" || name == "unknown-outcome"
+			if f.checks != 1 || j.OpenCheck != lost || j.CommitPhase != "" || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
+				t.Fatalf("check state changed after refusal: %d %+v", f.checks, j)
+			}
+		})
+	}
+}
+
+func TestOperatorFinishUnsandboxed(t *testing.T) {
+	f, w, _ := operatorLocalFinishFixture(t, func(m *store.OperatorFinishReview) { m.Checks[0].Unsandboxed = true })
+	for _, resume := range []bool{false, true} {
+		if resume {
+			if err := w.reserve(w.step("checks", "00")); err != nil {
+				t.Fatal(err)
+			}
+		}
+		before := operatorFinishJournal(t, w)
+		if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorFinishTerminal) ||
+			f.checks != 0 || !reflect.DeepEqual(before, operatorFinishJournal(t, w)) {
+			t.Fatalf("unsandboxed resume=%v: %v", resume, err)
+		}
+	}
+}
+
+func TestOperatorFinishPrepExpiry(t *testing.T) {
+	f, w, now := operatorLocalFinishFixture(t, nil)
+	w.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
+		f.checks++
+		*now = w.prep.row.ExpiresAt
+		return checkrun.Result{ExitCode: 0}, nil
+	}
+	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, store.ErrOperatorExpired) {
+		t.Fatal(err)
+	}
+	j := operatorFinishJournal(t, w)
+	if len(j.Checks) != 1 || j.CommitPhase != "" || f.checks != 1 {
+		t.Fatalf("expiry lost check receipt or advanced commit: %+v", j)
+	}
+}
+
+func TestOperatorFinishStale(t *testing.T) {
+	f, w, _ := operatorLocalFinishFixture(t, nil)
+	f.write("b.go", "changed untracked contents")
+	before := serviceFileHash(w.manifest.Store)
+	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorFinishReview) ||
+		serviceFileHash(w.manifest.Store) != before || f.checks != 0 || f.outcomes != 0 {
+		t.Fatalf("untracked drift crossed preflight: %v", err)
+	}
+	f.write("b.go", "hand\n")
+	f.gitOutput("-C", f.wt, "commit", "--allow-empty", "-qm", "foreign HEAD")
+	if err := runEpisodeFinishWithPrep(activationTestContext(t), w); !errors.Is(err, errOperatorFinishReview) ||
+		serviceFileHash(w.manifest.Store) != before || f.checks != 0 || f.outcomes != 0 {
+		t.Fatalf("HEAD drift crossed preflight: %v", err)
+	}
+}
+
+func TestOperatorFinishJournalOnce(t *testing.T) {
+	_, w, now := operatorLocalFinishFixture(t, nil, true)
+	ctx := activationTestContext(t)
+	next, err := w.read.NextOperatorFinish(ctx, *now)
+	if err != nil || next.ID != w.prep.row.ID || next.State != "confirmed" {
+		t.Fatalf("confirmed candidate: %s %s %v", next.ID, next.State, err)
+	}
+	c := runnerConfig{Credentials: w.manifest.Store, FinishRemotes: []string{w.manifest.Remote.URL}}
+	a, m, generation, err := finishManifestForRow(c, next)
+	if err != nil || !reflect.DeepEqual(m, w.manifest) || a.ID != next.ID || generation != 1 {
+		t.Fatalf("decode exact frozen intent: %d (%v)", generation, err)
+	}
+	c.FinishRemotes = nil
+	if _, _, _, err := finishManifestForRow(c, next); !errors.Is(err, approval.ErrNotVerified) {
+		t.Fatalf("unpinned remote: %v", err)
+	}
+	if _, err := w.read.NextOperatorFinish(ctx, next.ExpiresAt); !errors.Is(err, store.ErrOperatorMissing) {
+		t.Fatalf("expired confirmed candidate: %v", err)
+	}
+	h, err := w.read.OperatorHeartbeat(ctx)
+	if err != nil {
+		t.Fatal(err)
+	}
+	cas := operatorPreparationCAS(next)
+	if err := runnerWrite(w.manifest.Store, func(db *store.Credentials) error {
+		write, stop := w.budget(time.Second)
+		defer stop()
+		return db.ClaimOperatorAction(write, cas, h.Owner, *now)
+	}); err != nil {
+		t.Fatal(err)
+	}
+	w.prep.row, err = w.read.NextOperatorFinish(ctx, *now)
+	if err != nil || w.prep.row.Owner != h.Owner.Nonce || w.prep.row.Generation != cas.Generation+1 {
+		t.Fatalf("claim candidate: %s %v", w.prep.row.State, err)
+	}
+	key := w.step("commit", strings.Repeat("a", 40))
+	if err := w.reserve(key); err != nil {
+		t.Fatal(err)
+	}
+	other := w.step("commit", strings.Repeat("b", 40))
+	for _, done := range []bool{false, true} {
+		if done {
+			if err := w.observe(key, store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: strings.Repeat("c", 40)}); err != nil {
+				t.Fatal(err)
+			}
+		}
+		before := operatorFinishJournal(t, w)
+		if err := w.reserve(other); !errors.Is(err, store.ErrOperatorConflict) {
+			t.Fatalf("second distinct commit intent after observed=%t: %v", done, err)
+		}
+		if after := operatorFinishJournal(t, w); after.CommitKey.Identity != key.Identity || !reflect.DeepEqual(before, after) {
+			t.Fatalf("commit reservation changed after refusal: %+v", after)
+		}
+	}
+	if _, err := w.read.NextOperatorFinish(ctx, *now); !errors.Is(err, store.ErrOperatorMissing) {
+		t.Fatalf("completed candidate selected again: %v", err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_manifest.go b/cmd/ailang-worldd/operator_finish_manifest.go
index d5ea314..3f2d399 100644
--- a/cmd/ailang-worldd/operator_finish_manifest.go
+++ b/cmd/ailang-worldd/operator_finish_manifest.go
@@ -75,11 +75,26 @@ func finishManifestForRow(c runnerConfig, row store.OperatorActionRow) (approval
 
 func finishWorldEvidence(bound context.Context, path string, receipt store.OperatorFinishCheck) bool {
 	hash, err := hashref.Parse(receipt.Ref)
+	if err != nil || receipt.Entry < 1 {
+		return false
+	}
 	world, openErr := store.OpenReadOnly(path)
 	if openErr != nil {
 		return false
 	}
 	defer world.Close()
 	object, found, readErr := world.GetObject(bound, hash)
-	return err == nil && readErr == nil && found && hashref.SumSHA256(object.Payload) == hash
+	if readErr != nil || !found || object.SemanticID != taskOutcomeSemanticID || hashref.SumSHA256(object.Payload) != hash {
+		return false
+	}
+	_, outcome, err := parseTaskOutcome(object.Payload)
+	wantClass := "PASS"
+	if receipt.Result == "failed" {
+		wantClass = "FAIL"
+	}
+	if err != nil || outcome.Label != "episode-finish" || outcome.Task != receipt.Name || outcome.Outcome.Class != wantClass {
+		return false
+	}
+	entries, err := world.ObjectCommits(bound, hash, receipt.Entry-1, 1)
+	return err == nil && len(entries) == 1 && entries[0] == receipt.Entry
 }
diff --git a/cmd/ailang-worldd/operator_finish_manifest_test.go b/cmd/ailang-worldd/operator_finish_manifest_test.go
new file mode 100644
index 0000000..46192f5
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_manifest_test.go
@@ -0,0 +1,89 @@
+package main
+
+import (
+	"encoding/json"
+	"errors"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/checkrun"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorFinishManifestCapture(t *testing.T) {
+	f := newFinishFixture(t)
+	f.gitOutput("-C", f.wt, "branch", f.rep.Start.Base.Branch, f.before)
+	ctx := activationTestContext(t)
+	g := operatorFinishGitFixture(t, f)
+	in := operatorFinishInputs{report: f.rep, git: g, profile: &broker.ExecProfile{Protected: []string{"a.go"}}, world: f.db,
+		manifest: store.OperatorFinishReview{Base: f.rep.Start.Base, Commit: store.OperatorLocalCommit{Message: "Reviewed finish"}}}
+	m, err := in.frozen(ctx)
+	if err != nil || !reflect.DeepEqual(m.Protected, []string{"a.go"}) || !reflect.DeepEqual(m.Unattributed, []string{"b.go"}) ||
+		!strings.HasPrefix(m.Commit.Message, "Reviewed finish\n\n") || m.Commit.Snapshot.Head != f.before {
+		t.Fatalf("signed snapshot and provenance: %+v %v", m, err)
+	}
+	before := serviceFileHash(filepath.Join(g.admin, "index"))
+	in.manifest.Base.SHA = strings.Repeat("d", 40)
+	if _, err := in.frozen(ctx); !errors.Is(err, errOperatorFinishReview) || serviceFileHash(filepath.Join(g.admin, "index")) != before {
+		t.Fatalf("wrong base changed evidence: %v", err)
+	}
+}
+
+func TestOperatorFinishWorldEvidence(t *testing.T) {
+	ctx := activationTestContext(t)
+	path := filepath.Join(t.TempDir(), "world.db")
+	world, err := store.Open(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer world.Close()
+	h := func(s string) hashref.HashRef { return hashref.SumSHA256([]byte(s)) }
+	prev := store.World{Ref: h("genesis"), StateRoot: h("state"), LogHead: h("log")}
+	if err := world.PutWorld(activationTestContext(t), prev); err != nil {
+		t.Fatal(err)
+	}
+	if err := world.SelectHead(activationTestContext(t), prev.Ref); err != nil {
+		t.Fatal(err)
+	}
+	raw, err := json.Marshal(finishOutcome(ac5Report(), store.StartCheck{Name: "short", Argv: []string{"echo", "ok"}}, checkrun.Result{}, true))
+	if err != nil {
+		t.Fatal(err)
+	}
+	body, _, err := parseTaskOutcome(raw)
+	if err != nil {
+		t.Fatal(err)
+	}
+	result := taskOutcomeObject(body)
+	ref := hashref.SumSHA256(result.Payload)
+	obj := store.Object{Hash: ref, InterfaceHash: h("interface"), SemanticID: result.SemanticID, Provenance: result.Provenance, Payload: result.Payload}
+	if err := world.PutObject(activationTestContext(t), obj); err != nil {
+		t.Fatal(err)
+	}
+	r := store.OperatorFinishCheck{Name: "short", Result: "passed", Ref: ref.String(), Entry: 1}
+	if finishWorldEvidence(ctx, path, r) {
+		t.Fatal("loose object is not a check receipt")
+	}
+	next := store.World{Ref: h("next"), Revision: 1, StateRoot: h("next-state"), LogHead: h("next-log")}
+	if err := world.Commit(activationTestContext(t), store.Commit{ObservedHead: prev.Ref, Objects: []store.Object{obj}, NextWorld: next,
+		Entry: store.LogEntry{Header: store.LogHeader{EntryIndex: 1, SemanticsEpoch: 1, TransitionFn: ref,
+			Interpreter: h("interpreter"), PrevEntryHash: prev.LogHead, WrittenBy: taskOutcomeWrittenBy},
+			EntryHash: next.LogHead, TransitionRef: ref}}); err != nil {
+		t.Fatal(err)
+	}
+	if !finishWorldEvidence(ctx, path, r) {
+		t.Fatal("exact committed outcome was not readable")
+	}
+	for _, wrong := range []store.OperatorFinishCheck{
+		{Name: r.Name, Result: r.Result, Ref: r.Ref, Entry: 2},
+		{Name: "other", Result: r.Result, Ref: r.Ref, Entry: 1},
+		{Name: r.Name, Result: "failed", Ref: r.Ref, Entry: 1},
+	} {
+		if finishWorldEvidence(ctx, path, wrong) {
+			t.Fatalf("another check/entry reused a valid outcome object: %+v", wrong)
+		}
+	}
+}
diff --git a/host/store/operator_finish_test.go b/host/store/operator_finish_test.go
new file mode 100644
index 0000000..be409e6
--- /dev/null
+++ b/host/store/operator_finish_test.go
@@ -0,0 +1,73 @@
+package store
+
+import (
+	"encoding/json"
+	"errors"
+	"reflect"
+	"testing"
+)
+
+func TestOperatorFinishCompletedLaunch(t *testing.T) {
+	f, digest, prepared := preparedOperatorLaunchFixture(t, 0)
+	id := prepared.CAS.ID
+	now := func() int64 { return pinNow }
+	secret := [32]byte{1}
+	before := f.snapshot(t)
+	notCompleted := func() {
+		t.Helper()
+		if ids, err := f.c.OperatorCompletedLaunches(f.ctx); err != nil || len(ids) != 0 {
+			t.Fatalf("incomplete launch discovered: %v %v", ids, err)
+		}
+		if _, err := f.c.OperatorFinishLaunch(f.ctx, prepared.Episode); !errors.Is(err, ErrOperatorMissing) {
+			t.Fatalf("incomplete episode selected: %v", err)
+		}
+	}
+	notCompleted()
+	if !reflect.DeepEqual(before, f.snapshot(t)) {
+		t.Fatal("candidate read changed the prepared launch")
+	}
+	if err := f.c.ReserveOperatorLaunch(startRecoveryTestContext(t), prepared, secret, f.proof, now); err != nil {
+		t.Fatal(err)
+	}
+	if err := f.c.AuthorizeOperatorSpawn(startRecoveryTestContext(t), id, secret, f.proof, now); err != nil {
+		t.Fatal(err)
+	}
+	read := func() OperatorLaunchSnapshot {
+		t.Helper()
+		v, err := f.c.ReadOperatorLaunch(f.ctx, id)
+		if err != nil {
+			t.Fatal(err)
+		}
+		return v
+	}
+	receipt := OperatorLaunchReceipt{Run: id, Executable: digest, PID: 42, Birth: "kernel-start", Status: "launched"}
+	if err := f.c.RecordOperatorLaunch(startRecoveryTestContext(t), read(), receipt, nil, pinNow); err != nil {
+		t.Fatal(err)
+	}
+	notCompleted()
+	receipt.Status = "stopped"
+	if err := f.c.RecordOperatorLaunch(startRecoveryTestContext(t), read(), receipt, &EpisodeLaunch{Mode: "headless"}, pinNow); err != nil {
+		t.Fatal(err)
+	}
+	ids, err := f.c.OperatorCompletedLaunches(f.ctx)
+	got, readErr := f.c.OperatorFinishLaunch(f.ctx, prepared.Episode)
+	if err != nil || !reflect.DeepEqual(ids, []string{id}) || readErr != nil || got.Receipt != receipt || got.Key.Version != 4 {
+		t.Fatalf("completed candidate or exact receipt absent: %v %v %+v %v", ids, err, got, readErr)
+	}
+	v := OperatorFinishReview{Episode: prepared.Episode, Preset: "sol", Store: f.c.lock.dbPath, Report: "fixture"}
+	raw, err := json.Marshal(map[string]any{OperatorFinishCatalogueKey: map[string]OperatorFinishReview{v.Episode: v}})
+	if err != nil {
+		t.Fatal(err)
+	}
+	h := OperatorHeartbeat{Catalogue: raw}
+	review, err := h.FinishReview(v.Episode, v.Preset)
+	if err != nil || !reflect.DeepEqual(review, v) {
+		t.Fatalf("catalogue finish projection: %+v %v", review, err)
+	}
+	if _, err := h.FinishReview(v.Episode, "other"); !errors.Is(err, ErrOperatorRunner) {
+		t.Fatalf("foreign preset accepted: %v", err)
+	}
+	if a, err := review.Action(4, true); err != nil || a.Kind != "finish" || !a.Options.Ready || a.Task != "" {
+		t.Fatalf("taskless typed review: %+v %v", a, err)
+	}
+}
````

### M3 — S9c activation draft. Planner rehearsal only

````diff M3
diff --git a/cmd/ailang-worldd/operator_activation_test.go b/cmd/ailang-worldd/operator_activation_test.go
index 1ef9292..3f186da 100644
--- a/cmd/ailang-worldd/operator_activation_test.go
+++ b/cmd/ailang-worldd/operator_activation_test.go
@@ -43,7 +43,12 @@ func operatorActivationFixtureWithAgent(t *testing.T, script string) (*worktreeS
 
 func operatorActivationFixtureWithTask(t *testing.T, script, task string) (*worktreeStartRig, runnerConfig, store.OperatorHeartbeat, *os.File) {
 	t.Helper()
-	r := newWorktreeStartRigWithTask(t, false, activationTestContext(t), script, task)
+	return operatorActivationFixtureWithRemote(t, script, task, "")
+}
+
+func operatorActivationFixtureWithRemote(t *testing.T, script, task, remote string) (*worktreeStartRig, runnerConfig, store.OperatorHeartbeat, *os.File) {
+	t.Helper()
+	r := newWorktreeStartRigWithRemote(t, false, activationTestContext(t), script, task, remote)
 	r.env.operatorBudget = activationTestBudget
 	c := r.config
 	agent := operatorFileFact(r.action, "launch_binary")
diff --git a/cmd/ailang-worldd/operator_catalogue.go b/cmd/ailang-worldd/operator_catalogue.go
index 190d331..6e464bd 100644
--- a/cmd/ailang-worldd/operator_catalogue.go
+++ b/cmd/ailang-worldd/operator_catalogue.go
@@ -16,7 +16,7 @@ import (
 
 // Read the configured presets on runner ticks, independent of HTTP input.
 func (e serviceEnv) readOperatorCatalogue(c runnerConfig, src operatorStartSource) []byte {
-	items := make(map[string]store.OperatorStartReview)
+	items := make(map[string]any)
 	r, err := readServiceRecord(e.paths)
 	if err != nil {
 		return []byte("{}")
@@ -35,8 +35,18 @@ func (e serviceEnv) readOperatorCatalogue(c runnerConfig, src operatorStartSourc
 			}
 		}
 	}
+	if len(c.FinishRemotes) != 0 {
+		src.workspace = r.Workspace
+		finishCtx, stop := e.operatorContext(4 * time.Second)
+		items[store.OperatorFinishCatalogueKey] = e.finishCatalogue(finishCtx, c, src)
+		stop()
+	}
 	raw, _ := json.Marshal(items)
-	if ctx.Err() != nil || len(raw) > 16*approval.MaxOperatorActionBytes {
+	if len(raw) > 16*approval.MaxOperatorActionBytes {
+		delete(items, store.OperatorFinishCatalogueKey)
+		raw, _ = json.Marshal(items)
+	}
+	if len(raw) > 16*approval.MaxOperatorActionBytes {
 		return []byte("{}")
 	}
 	return raw
diff --git a/cmd/ailang-worldd/operator_finish_route_test.go b/cmd/ailang-worldd/operator_finish_route_test.go
new file mode 100644
index 0000000..47e474b
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_route_test.go
@@ -0,0 +1,316 @@
+package main
+
+import (
+	"bytes"
+	"encoding/hex"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"net/http"
+	"net/http/httptest"
+	"net/url"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/worktree"
+)
+
+func TestOperatorFinishRouteWorker(t *testing.T) {
+	f := newOperatorLaunchRigWithRemote(t, fakeLaunchMarker+"exit 0\n", "private-task-fixture-marker", fixtureFinishRemote)
+	if out := f.tick(t); !strings.Contains(out, ": stopped") {
+		t.Fatalf("the fake child never completed: %s", out)
+	}
+	ctx := activationTestContext(t)
+	start, err := f.db.EpisodeStart(ctx, f.action.Episode)
+	if err != nil || start == nil || len(start.Checks) != 1 || start.Checks[0].Unsandboxed {
+		t.Fatalf("completed start missing pinned checks: %v", err)
+	}
+	const filename = "private-path-fixture-marker"
+	change := filepath.Join(f.path(), filename)
+	seedStartFixtureFile(t, change, []byte("approved change\n"))
+	repo := f.f.r.f
+	bare := filepath.Join(t.TempDir(), "remote.git")
+	presetGit(t, repo, "clone", "-q", "--local", "--bare", "--", repo.repo, bare)
+	remoteHead := presetGit(t, repo, "--git-dir="+bare, "rev-parse", "HEAD")
+	trace := filepath.Join(f.home, "unexpected-remote-write")
+	tools := f.c.ToolDirs[0]
+	f.c.FinishRemotes = []string{fixtureFinishRemote}
+	writeFinishTestTool(t, &f.c, tools, "gh", "case \"$1 $2\" in 'auth status'|'label list') ;; 'repo view') printf 'WRITE\\n';; *) touch "+shellQuote(trace)+"; exit 96;; esac\n")
+	writeFinishTestTool(t, &f.c, tools, "git", "for arg; do case \"$arg\" in push|fetch|clone|pull|send-pack|fetch-pack|receive-pack|upload-pack|remote|submodule) touch "+shellQuote(trace)+"; exit 96;; esac; done\n"+
+		"case \" $* \" in *' ls-remote '*) exec "+shellQuote(repo.git)+" --git-dir=/dev/null ls-remote --symref -- "+shellQuote(bare)+" HEAD;; esac\n"+
+		"exec "+shellQuote(repo.git)+" \"$@\"\n")
+	paths := f.env.runnerPaths()
+	config, err := json.Marshal(f.c)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := f.env.atomic(paths.Record, config, 0600); err != nil {
+		t.Fatal(err)
+	}
+	if err := f.env.atomic(paths.Plist, []byte(runnerPlist(paths, f.c, config)), 0600); err != nil {
+		t.Fatal(err)
+	}
+
+	const origin = "http://localhost:7644"
+	log := &syncBuffer{}
+	d, err := daemon.New(ctx, daemon.Config{DBPath: f.c.World, CredentialsPath: f.creds, BindHost: daemon.DefaultBindHost,
+		SurfaceOrigins: []string{origin}, ErrorLog: log, ApprovalClock: f.src.clock})
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = d.Close() })
+	srv := httptest.NewServer(d.Handler())
+	t.Cleanup(srv.Close)
+	request := func(method, path string, input any, out any) string {
+		t.Helper()
+		var data []byte
+		if input != nil {
+			data, err = json.Marshal(input)
+			if err != nil {
+				t.Fatal(err)
+			}
+		}
+		requestCtx, cancel := f.env.operatorContext(activationTestBudget)
+		defer cancel()
+		req, err := http.NewRequestWithContext(requestCtx, method, srv.URL+path, bytes.NewReader(data))
+		if err != nil {
+			t.Fatal(err)
+		}
+		req.Header.Set("Origin", origin)
+		req.Header.Set("Sec-Fetch-Site", "same-origin")
+		req.Header.Set("Content-Type", "application/json")
+		res, err := srv.Client().Do(req)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer res.Body.Close()
+		raw, err := io.ReadAll(res.Body)
+		if err != nil || res.StatusCode != 200 || out != nil && json.Unmarshal(raw, out) != nil {
+			t.Fatalf("%s %s: %d %s (%v)", method, path, res.StatusCode, raw, err)
+		}
+		return string(raw)
+	}
+
+	rep := daemon.EpisodeReport{Episode: start.Episode, Start: start, Entries: [2]int64{1, 1},
+		Writes: daemon.ReportWrites{Written: []daemon.ReportWrite{{Path: filename, Op: "write", Count: 1}}}}
+	h := hashref.SumSHA256([]byte("selected world")).String()
+	checks, outcomes := 0, 0
+	activeID := ""
+	reportDown := false
+	seen := []string{}
+	original := finishClient
+	t.Cleanup(func() { finishClient = original })
+	finishClient = func(string) *client {
+		return memoryDaemon(func(w http.ResponseWriter, r *http.Request) {
+			seen = append(seen, r.URL.Path)
+			switch r.URL.Path {
+			case "/v1/episodes/" + start.Episode + "/report":
+				if reportDown {
+					w.WriteHeader(http.StatusServiceUnavailable)
+					return
+				}
+				_ = json.NewEncoder(w).Encode(rep)
+			case "/mcp/":
+				terminal, lockErr := openTerminalEpisodeLock(f.creds, start.Episode, nil, os.O_CREATE, 5*time.Second)
+				if terminal != nil {
+					_ = terminal.Close()
+				}
+				if !errors.Is(lockErr, store.ErrOperatorConflict) || !strings.Contains(lockErr.Error(), "active action "+activeID) {
+					t.Errorf("terminal did not see the active private finish: %v", lockErr)
+				}
+				var call struct {
+					Params struct {
+						Name      string
+						Arguments struct {
+							Command string
+							Args    []string
+						}
+					}
+				}
+				if r.Header.Get("Authorization") != "Bearer "+f.token || json.NewDecoder(r.Body).Decode(&call) != nil ||
+					call.Params.Name != "workspace-exec" || call.Params.Arguments.Command != start.Checks[0].Exec ||
+					!reflect.DeepEqual(call.Params.Arguments.Args, start.Checks[0].Args) {
+					t.Errorf("check used different argv or session")
+				}
+				checks++
+				w.Header().Set("Content-Type", "text/event-stream")
+				fmt.Fprintf(w, "data: {\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{\"content\":[{\"type\":\"text\",\"text\":%q}]}}\n\n", `{"exit_code":0}`)
+			case "/v1/head":
+				fmt.Fprint(w, h)
+			case "/v1/worlds/" + h:
+				_ = json.NewEncoder(w).Encode(wireWorld{Ref: h, Revision: 1, LogHead: h})
+			case "/v1/log/1":
+				_ = json.NewEncoder(w).Encode(wireEntry{Header: wireHeader{EntryIndex: 1, SemanticsEpoch: 1, Interpreter: h}})
+			case "/v1/commit":
+				var commit demoCommit
+				if r.Header.Get("Authorization") != "Bearer "+f.token || json.NewDecoder(r.Body).Decode(&commit) != nil ||
+					len(commit.Objects) != 1 || commit.Objects[0].SemanticID != taskOutcomeSemanticID {
+					t.Errorf("check outcome writer crossed the pinned evidence seam")
+				}
+				outcomes++
+				_ = json.NewEncoder(w).Encode(map[string]string{"selectedHead": commit.NextWorld.Ref})
+			default:
+				t.Errorf("unexpected remote, V3 or finish-record operation: %s %s", r.Method, r.URL.Path)
+				w.WriteHeader(http.StatusBadRequest)
+			}
+		})
+	}
+
+	source := f.src
+	source.workspace = "" // The real runner starts with the workspace supplied by its installed record.
+	raw := f.env.readOperatorCatalogue(f.c, source)
+	f.h.Build, f.h.Readiness, f.h.Catalogue, f.h.Owner.Generation = daemon.CurrentBuild().ExecutableSHA256, "both", raw, 0
+	entry, err := f.h.FinishReview(start.Episode, start.Preset.Name)
+	if err != nil || entry.Commit.Snapshot.Head != start.Base.SHA {
+		t.Fatalf("completed launch was not captured from the installed workspace: %v", err)
+	}
+	if err := f.env.touchRunner(f.c, &f.h); err != nil {
+		t.Fatal(err)
+	}
+	var review struct {
+		Outcome             string
+		Argv                []string
+		CatalogueGeneration int64
+	}
+	command := "episode finish " + start.Episode + " --preset " + start.Preset.Name + " --confirm-on-surface"
+	request("GET", "/surface/verb.json?q="+url.QueryEscape(command), nil, &review)
+	if review.Outcome != "review" {
+		t.Fatalf("missing review: %+v", review)
+	}
+	var status struct{ ID, Href, StatusHref, State, LocalCommit string }
+	request("POST", "/surface/operator-actions", map[string]any{
+		"argv": review.Argv, "operationId": strings.Repeat("4", 32), "catalogueGeneration": review.CatalogueGeneration, "task": "",
+	}, &status)
+	row, err := f.db.OperatorAction(ctx, status.ID, f.now)
+	if err != nil {
+		t.Fatal(err)
+	}
+	key := approvaltest.NewSigner(54)
+	seedApprover(t, f.creds, "finish-operator-fixture", key.PublicKey(), f.now-1)
+	message, err := approval.OperatorActionMessage([32]byte(row.Digest), [32]byte(row.Challenge), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	artifact := key.Sign(approval.OperatorActionNamespace, message)
+	request("POST", status.Href+"/confirm", map[string]any{"v": 1, "digest": hex.EncodeToString(row.Digest), "sshsig": string(artifact)}, &status)
+
+	activeID = row.ID
+	run := func() error {
+		turn, stop := f.env.operatorContext(activationTestBudget)
+		defer stop()
+		return f.env.consumeFinish(turn, f.c, f.h, f.src)
+	}
+	before := serviceFileHash(f.creds)
+	held := f.borrow(t)
+	blocked := run()
+	if err := held.Close(); err != nil {
+		t.Fatal(err)
+	}
+	if !errors.Is(blocked, store.ErrOperatorConflict) || !strings.Contains(blocked.Error(), row.ID) || before != serviceFileHash(f.creds) {
+		t.Fatalf("lock contention crossed admission: %v", blocked)
+	}
+	refuse := func(name string, want error) {
+		t.Helper()
+		if err := run(); !errors.Is(err, want) || before != serviceFileHash(f.creds) || checks+outcomes != 0 {
+			t.Fatalf("%s changed evidence before claim: %v", name, err)
+		}
+	}
+	originalHealth := f.health
+	f.health.DegradedEpisodes = []string{start.Episode}
+	refuse("degraded health", errOperatorFinishEligibility)
+	f.health = originalHealth
+	reportDown = true
+	refuse("report unavailable", errOperatorFinishEligibility)
+	reportDown = false
+	if err := os.Remove(change); err != nil {
+		t.Fatal(err)
+	}
+	refuse("empty worktree", errOperatorFinishReview)
+	seedStartFixtureFile(t, change, []byte("different approved input\n"))
+	refuse("untracked contents", errOperatorFinishReview)
+	seedStartFixtureFile(t, change, []byte("approved change\n"))
+
+	branch := "refs/heads/" + start.Branch
+	baseBranch := "refs/heads/" + start.Base.Branch
+	tree := presetGit(t, repo, "rev-parse", start.Base.SHA+"^{tree}")
+	foreign := presetGit(t, repo, "-c", "user.name=Fixture", "-c", "user.email=fixture@invalid",
+		"commit-tree", tree, "-p", start.Base.SHA, "-m", "new fixture HEAD")
+	for _, ref := range []string{branch, baseBranch} {
+		presetGit(t, repo, "update-ref", ref, foreign, start.Base.SHA)
+		refuse(ref, errOperatorFinishReview)
+		presetGit(t, repo, "update-ref", ref, start.Base.SHA, foreign)
+	}
+	presetGit(t, repo, "remote", "set-url", start.Remote.Name, "https://github.com/fixture/changed.git")
+	refuse("remote", errOperatorFinishReview)
+	presetGit(t, repo, "remote", "set-url", start.Remote.Name, start.Remote.URL)
+
+	_, profilePath, err := loadMintProfile(f.src.profiles, "proj")
+	if err != nil {
+		t.Fatal(err)
+	}
+	for _, path := range []string{f.f.file, profilePath} {
+		pinned := readStartFixtureFile(t, path)
+		seedStartFixtureFile(t, path, append(bytes.Clone(pinned), '\n'))
+		refuse(filepath.Base(path), errOperatorFinishReview)
+		seedStartFixtureFile(t, path, pinned)
+	}
+	token := filepath.Join(f.f.r.sessions, start.Episode+".session")
+	saved := readStartFixtureFile(t, token)
+	seedStartFixtureFile(t, token, []byte(strings.Repeat("0", len(f.token))))
+	refuse("changed session", errOperatorFinishEligibility)
+	seedStartFixtureFile(t, token, saved)
+	expired := f.src
+	expired.clock = func() int64 { return f.self.ExpiresAt }
+	if _, err := f.env.finishInputs(ctx, f.c, expired, f.db, start.Episode, nil); !errors.Is(err, errOperatorFinishEligibility) ||
+		before != serviceFileHash(f.creds) || checks+outcomes != 0 {
+		t.Fatalf("expired session changed evidence: %v", err)
+	}
+	if err := run(); err != nil {
+		t.Fatal(err)
+	}
+	row, err = f.db.OperatorAction(ctx, status.ID, f.now)
+	if err != nil {
+		t.Fatal(err)
+	}
+	j, err := f.db.ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))
+	head, headErr := worktree.Output(ctx, repo.git, []string{"-C", f.path(), "rev-parse", "HEAD"})
+	if err != nil || headErr != nil || j.Commit.SHA != strings.TrimSpace(head) || j.Commit.State != store.OperatorLocalCommitDone ||
+		row.V3ID != "" || checks != 1 || outcomes != 1 || len(j.Checks) != 1 || j.Checks[0].Result != "passed" {
+		t.Fatalf("local hand-off: %+v %v head %v, checks %d, outcomes %d, calls %v", j, err, headErr, checks, outcomes, seen)
+	}
+	if err := run(); !errors.Is(err, store.ErrOperatorMissing) || checks != 1 || outcomes != 1 {
+		t.Fatalf("durable local hand-off replayed: %v", err)
+	}
+	request("GET", status.StatusHref, nil, &status)
+	if status.LocalCommit != j.Commit.SHA {
+		t.Fatalf("redacted status lacked local commit: %+v", status)
+	}
+	var cli, cliErr syncBuffer
+	if rc := runServiceRunner([]string{"finish-status", status.ID}, &cli, &cliErr, f.env); rc != exitOK ||
+		!strings.Contains(cli.String(), j.Commit.SHA) || !strings.Contains(cli.String(), "passed") {
+		t.Fatalf("finish-status did not use the redacted receipt: %d %s %s", rc, &cli, &cliErr)
+	}
+	output := cli.String() + cliErr.String() + request("GET", status.StatusHref, nil, nil) + request("GET", "/surface/requests.json", nil, nil) +
+		request("GET", "/surface/home-summary.json", nil, nil) + log.String()
+	for _, secret := range []string{filename, "private-task-fixture-marker", "BEGIN SSH SIGNATURE", f.token, "finish_manifest"} {
+		if strings.Contains(output, secret) {
+			t.Fatalf("routine output exposed %q", secret)
+		}
+	}
+	if got := presetGit(t, repo, "--git-dir="+bare, "rev-parse", "HEAD"); got != remoteHead {
+		t.Fatalf("local runner changed the bare remote: %s", got)
+	}
+	if _, err := os.Lstat(trace); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("remote write call escaped the worker: %v", err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_runner.go b/cmd/ailang-worldd/operator_finish_runner.go
new file mode 100644
index 0000000..856433c
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_runner.go
@@ -0,0 +1,410 @@
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
+	"reflect"
+	"slices"
+	"strings"
+	"syscall"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/checkrun"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+	"github.com/sunholo-data/ailang-world/host/worktree"
+)
+
+var errOperatorFinishEligibility = errors.New("finish needs an eligible episode with a stopped launch and its original live session")
+
+type operatorFinishBusy string
+
+func (id operatorFinishBusy) Error() string {
+	return fmt.Sprintf("active action %s holds the episode turn; run ailang-worldd service runner status", string(id))
+}
+
+func (operatorFinishBusy) Unwrap() error { return store.ErrOperatorConflict }
+
+func (e serviceEnv) finishInputs(ctx context.Context, c runnerConfig, src operatorStartSource, ro *store.Credentials,
+	ep string, lock *os.File) (operatorFinishInputs, error) {
+	var out operatorFinishInputs
+	if !verbs.EpisodeIDPattern.MatchString(ep) {
+		return out, errOperatorFinishEligibility
+	}
+	start, err := ro.EpisodeStart(ctx, ep)
+	if err != nil || start == nil || len(start.Resumes) != 0 || !slices.Contains(c.PresetNames, start.Preset.Name) ||
+		!slices.Contains(c.FinishRemotes, start.Remote.URL) {
+		return out, errOperatorFinishEligibility
+	}
+	launch, err := ro.OperatorFinishLaunch(ctx, ep)
+	if err != nil || launch.Key.Version != 4 {
+		return out, errOperatorFinishEligibility
+	}
+	if birth, live := operatorProcessIdentity(launch.Receipt.PID); live && birth == launch.Receipt.Birth {
+		return out, errOperatorFinishEligibility
+	}
+	health, err := e.health(c.Address)
+	if err != nil || health.Status != "ok" || slices.Contains(health.DegradedEpisodes, ep) ||
+		health.Build.ExecutableSHA256 != c.Executables[e.paths.Binary] {
+		return out, errOperatorFinishEligibility
+	}
+	_, rep, err := fetchReport(ctx, finishClient(c.Address), ep)
+	if err != nil || rep.Episode != ep || rep.Uncommitted != 0 || rep.Finish != nil || !reflect.DeepEqual(rep.Start, start) {
+		return out, errOperatorFinishEligibility
+	}
+	presetDir, err := filepath.EvalSymlinks(src.presets)
+	if err != nil {
+		return out, errOperatorFinishReview
+	}
+	raw, _, err := operatorFinishRead(presetDir, start.Preset.Name+".json", 128<<10)
+	if err != nil || serviceHash(raw) != start.Preset.SHA256 || string(raw) != start.Preset.Bytes {
+		return out, errOperatorFinishReview
+	}
+	var preset launchPreset
+	if strictPresetJSON(raw, &preset) != nil || preset.Name != start.Preset.Name {
+		return out, errOperatorFinishReview
+	}
+	p, err := finishPreset(start.Preset.Name, src.presets)
+	physical, physicalErr := filepath.EvalSymlinks(p.WorkspaceRoot)
+	root, rootErr := filepath.EvalSymlinks(src.workspace)
+	if err != nil || physicalErr != nil || rootErr != nil || physical != root || p.Addr != c.Address {
+		return out, errOperatorFinishReview
+	}
+	if equal, _ := store.SameDatabase(p.DB, c.World); !equal {
+		return out, errOperatorFinishReview
+	}
+	creds, err := credentialsFor(p.DB, p.CredentialsDB)
+	if equal, _ := store.SameDatabase(creds, c.Credentials); err != nil || !equal {
+		return out, errOperatorFinishReview
+	}
+	profile, _, err := loadMintProfile(src.profiles, preset.Profile)
+	if err != nil {
+		return out, errOperatorFinishReview
+	}
+	for _, check := range start.Checks {
+		if check.Unsandboxed {
+			return out, errOperatorFinishTerminal
+		}
+	}
+	checks, err := validateGraderChecks(preset.Checks, preset.Name)
+	if err != nil || pinPresetChecks(checks, profile) != nil || !reflect.DeepEqual(checks, start.Checks) {
+		return out, errOperatorFinishReview
+	}
+	sessions, err := ro.ListSessions(ctx, ep, 2)
+	if err != nil || len(sessions) != 1 || sessions[0].ExpiresAt <= src.clock() ||
+		sessions[0].Project != preset.Profile {
+		return out, errOperatorFinishEligibility
+	}
+	if sessions[0].ProfileDigest != profile.Digest {
+		return out, errOperatorFinishReview
+	}
+	sessionDir, err := filepath.EvalSymlinks(p.SessionsDir)
+	if err != nil || runnerPath(c.Home, sessionDir, c.UID, true, 0700) != nil {
+		return out, errOperatorFinishEligibility
+	}
+	token, mode, err := operatorFinishRead(sessionDir, ep+".session", maxSessionFileBytes)
+	if err != nil || mode.Perm() != 0600 {
+		return out, errOperatorFinishEligibility
+	}
+	out.token = strings.TrimSpace(string(token))
+	resolved, err := authority.New(ro).ResolveContext(ctx, "Bearer "+out.token, src.clock())
+	if err != nil || resolved.Success == nil || resolved.Success.EpisodeID != ep || resolved.Success.ExpiresAt != sessions[0].ExpiresAt {
+		return out, errOperatorFinishEligibility
+	}
+
+	auth, err := planOperatorFinishAuth(c, start.Remote.URL)
+	if err != nil {
+		return out, err
+	}
+	g, err := openOperatorFinishGit(auth, filepath.Join(root, ep), start.Repo, lock)
+	if err != nil || g.dir != filepath.Join(root, ep) {
+		return out, errOperatorFinishFence
+	}
+	name, email, err := worktree.ParseCommitter(preset.Committer)
+	if err != nil || name == "" {
+		return out, errOperatorFinishFence
+	}
+	title := "World episode " + ep
+	if start.Task != nil && start.Task.Title != "" {
+		title = start.Task.Title
+	}
+	rep.ExecCache = nil
+	source, _ := json.Marshal(rep)
+	sessionWire, _ := json.Marshal(sessions)
+	launchWire, _ := json.Marshal(launch.Receipt)
+	out.git, out.report, out.profile, out.world = g, rep, profile, c.World
+	out.manifest = store.OperatorFinishReview{Store: c.Credentials, Episode: ep, Preset: start.Preset.Name, Address: c.Address,
+		Worktree: g.dir, Report: finishSHA(source), PresetDigest: start.Preset.SHA256, Profile: profile.Digest,
+		Sessions: finishSHA(sessionWire), Launch: finishSHA(launchWire), Remote: start.Remote, Base: start.Base, Checks: start.Checks,
+		Commit: store.OperatorLocalCommit{Committer: name + " <" + email + ">", Message: title}}
+	return out, nil
+}
+
+func (e serviceEnv) finishCatalogue(ctx context.Context, c runnerConfig, src operatorStartSource) map[string]store.OperatorFinishReview {
+	catalogue := make(map[string]store.OperatorFinishReview)
+	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
+	if err != nil {
+		return catalogue
+	}
+	defer ro.Close()
+	ids, err := ro.OperatorCompletedLaunches(ctx)
+	if err != nil {
+		return catalogue
+	}
+	for _, id := range ids {
+		row, err := ro.OperatorAction(ctx, id, src.clock())
+		if err != nil {
+			continue
+		}
+		in, err := e.finishInputs(ctx, c, src, ro, row.Episode, nil)
+		if err != nil {
+			continue
+		}
+		v, err := in.frozen(ctx)
+		if err == nil {
+			if _, err := v.Action(1, false); err == nil {
+				catalogue[row.Episode] = v
+			}
+		}
+	}
+	return catalogue
+}
+
+func (e serviceEnv) consumeFinish(ctx context.Context, c runnerConfig, h store.OperatorHeartbeat, src operatorStartSource) error {
+	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
+	if err != nil {
+		return err
+	}
+	defer ro.Close()
+	row, err := ro.NextOperatorFinish(ctx, src.clock())
+	if err != nil {
+		return err
+	}
+	a, m, generation, err := finishManifestForRow(c, row)
+	if err != nil {
+		return err
+	}
+	cas := operatorPreparationCAS(row)
+	journal, err := ro.ReadOperatorFinishJournal(ctx, cas)
+	if err != nil {
+		return err
+	}
+	observation := row.Owner != "" && journal.CommitPhase == "reserved"
+	verifyAt := src.clock()
+	if observation {
+		// Expiry stops new effects. Use the signed deadline for proof of what
+		// was approved, then observe only the already reserved commit.
+		verifyAt = min(verifyAt, row.ExpiresAt-1)
+	}
+	prep, err := verifyOperatorFinishPrep(ctx, ro, a, verifyAt, cas, row.Owner)
+	if err != nil {
+		return err
+	}
+	for _, check := range m.Checks {
+		if check.Unsandboxed {
+			return errOperatorFinishTerminal
+		}
+	}
+	if err := prepareOperatorFinish(ctx, c, m.Remote.URL); err != nil {
+		return err
+	}
+
+	lock, err := store.OpenOperatorEpisodeLock(c.Credentials, row.Episode)
+	if errors.Is(err, syscall.EWOULDBLOCK) {
+		return operatorFinishBusy(row.ID)
+	}
+	if err != nil {
+		return err
+	}
+	defer lock.Close()
+	fd, err := lock.Verify(c.Credentials, row.Episode)
+	if err != nil {
+		return err
+	}
+	if observation {
+		if _, err := verifyOperatorFinishPrep(ctx, ro, a, verifyAt, cas, row.Owner); err != nil {
+			return err
+		}
+		return e.recoverFinishCommit(ctx, c, src, ro, fd, prep, m, journal.CommitKey)
+	}
+	old := row.Owner != ""
+	if _, err := verifyOperatorFinishPrep(ctx, ro, a, src.clock(), cas, row.Owner); err != nil {
+		return err
+	}
+	in, err := e.finishInputs(ctx, c, src, ro, row.Episode, fd)
+	if err != nil {
+		return err
+	}
+	if !finishReviewInputsMatch(in.manifest, m) {
+		return errOperatorFinishReview
+	}
+	if !old {
+		reviewed, err := in.frozen(ctx)
+		if err != nil {
+			return errors.Join(errOperatorFinishReview, err)
+		}
+		if !reflect.DeepEqual(reviewed, m) || len(m.Protected) != 0 || len(m.Unattributed) != 0 {
+			return errOperatorFinishReview
+		}
+		admit, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
+		if err != nil {
+			return errOperatorDraining
+		}
+		if ctx.Err() == nil && e.checkStartPair(c, src.workspace) == nil && e.proveOldTerminalQuiet() == nil {
+			err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
+				write, cancel := e.operatorContext(time.Second)
+				defer cancel()
+				if _, err := verifyOperatorFinishPrep(write, db, a, src.clock(), cas, ""); err != nil {
+					return err
+				}
+				live, err := db.OperatorHeartbeat(write)
+				if err != nil || live.Owner != h.Owner || live.CatalogueGeneration != generation || !live.Available(h.Build, "finish", src.clock()) {
+					return store.ErrOperatorRunner
+				}
+				return db.ClaimOperatorAction(write, cas, h.Owner, src.clock())
+			})
+		} else {
+			err = store.ErrOperatorRunner
+		}
+		admit.Close()
+		if err != nil {
+			return err
+		}
+		cas.Generation++
+		row.Owner = h.Owner.Nonce
+	}
+	prep, err = verifyOperatorFinishPrep(ctx, ro, a, src.clock(), cas, row.Owner)
+	if err != nil {
+		return err
+	}
+
+	client := finishClient(c.Address)
+	w := operatorFinishWork{prep: prep, generation: generation, manifest: m, report: in.report, git: in.git,
+		spool: filepath.Dir(e.runnerPaths().Log), read: ro, now: src.clock, budget: e.operatorContext}
+	w.fresh = func(bound context.Context) error {
+		if err := bound.Err(); err != nil {
+			return err
+		}
+		live, err := ro.OperatorHeartbeat(bound)
+		if err != nil || live.Owner != h.Owner || row.Owner != h.Owner.Nonce || !live.Available(h.Build, "finish", src.clock()) ||
+			e.checkStartPair(c, src.workspace) != nil {
+			return store.ErrOperatorRunner
+		}
+		current, err := e.finishInputs(bound, c, src, ro, row.Episode, fd)
+		if err != nil {
+			return err
+		}
+		if !finishReviewInputsMatch(current.manifest, m) || !current.baseMatches(bound) {
+			return errOperatorFinishReview
+		}
+		_, err = verifyOperatorFinishPrep(bound, ro, a, src.clock(), cas, row.Owner)
+		return err
+	}
+	w.run = func(bound context.Context, check store.StartCheck) (checkrun.Result, error) {
+		return finishSandboxCheck(bound, client, in.token, check)
+	}
+	w.record = func(bound context.Context, outcome demoObject) (int64, error) {
+		return commitOneObject(bound, client, in.token, outcome, taskOutcomeWrittenBy, "", "outcome record")
+	}
+	w.evidence = func(bound context.Context, receipt store.OperatorFinishCheck) bool {
+		return finishWorldEvidence(bound, c.World, receipt)
+	}
+	w.commit = func(bound context.Context, spec store.OperatorLocalCommit, contents [][]byte, tree string) error {
+		return in.git.commit(bound, w.spool, spec, contents, tree, in.report.Start)
+	}
+	return runEpisodeFinishWithPrep(ctx, w)
+}
+
+func finishReviewInputsMatch(current, frozen store.OperatorFinishReview) bool {
+	return current.Remote == frozen.Remote && current.Base == frozen.Base && current.PresetDigest == frozen.PresetDigest &&
+		current.Profile == frozen.Profile && current.Sessions == frozen.Sessions && current.Launch == frozen.Launch &&
+		reflect.DeepEqual(current.Checks, frozen.Checks)
+}
+
+func (e serviceEnv) recoverFinishCommit(ctx context.Context, c runnerConfig, src operatorStartSource, ro *store.Credentials,
+	fd *os.File, prep operatorFinishPrep, m store.OperatorFinishReview, key store.OperatorStepKey) error {
+	start, err := ro.EpisodeStart(ctx, m.Episode)
+	if err != nil || start == nil || start.Remote != m.Remote || start.Base != m.Base {
+		return errOperatorFinishUncertain
+	}
+	auth, err := planOperatorFinishAuth(c, m.Remote.URL)
+	if err != nil {
+		return err
+	}
+	git, err := openOperatorFinishGit(auth, filepath.Join(src.workspace, m.Episode), start.Repo, fd)
+	if err != nil || git.dir != m.Worktree {
+		return errOperatorFinishUncertain
+	}
+	w := operatorFinishWork{prep: prep, manifest: m, report: daemon.EpisodeReport{Start: start}, git: git, budget: e.operatorContext}
+	return w.reconcileCommit(key)
+}
+
+func (e serviceEnv) printFinishStatus(id string, out io.Writer) error {
+	c, _, err := e.installedRunner()
+	if err != nil {
+		return err
+	}
+	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
+	if err != nil {
+		return errOperatorFinishEligibility
+	}
+	defer ro.Close()
+	ctx, cancel := e.operatorContext(time.Second)
+	defer cancel()
+	row, err := ro.OperatorAction(ctx, id, e.terminal.now())
+	if err != nil || row.Kind != "finish" {
+		return store.ErrOperatorMissing
+	}
+	j, err := ro.ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))
+	if err != nil {
+		return errOperatorFinishUncertain
+	}
+	fmt.Fprintln(out, "finish action", row.ID, row.State)
+	for _, check := range j.Checks {
+		fmt.Fprintf(out, "check %q: %s\n", check.Name, check.Result)
+	}
+	if j.OpenCheck || j.CommitPhase == "reserved" {
+		fmt.Fprintln(out, errOperatorFinishUncertain)
+	}
+	if j.Commit.State == store.OperatorLocalCommitDone {
+		fmt.Fprintln(out, "local commit", j.Commit.SHA, "retained; PR confirmation has not been issued")
+	} else if j.Commit.State == "commit-failed" {
+		fmt.Fprintln(out, errOperatorCommitFailed)
+	}
+	return nil
+}
+
+// An admitted check may outlast the heartbeat's admission window. The runner
+// keeps the same pinned catalogue alive while the exclusive finish turn runs.
+func (e serviceEnv) runFinishTurn(c runnerConfig, h *store.OperatorHeartbeat, src operatorStartSource,
+	stops <-chan os.Signal, ticks <-chan time.Time) error {
+	ctx, cancel := e.operatorContext(2 * time.Hour)
+	defer cancel()
+	reply := make(chan error, 1)
+	look := *h
+	go func() { reply <- e.consumeFinish(ctx, c, look, src) }()
+	for {
+		select {
+		case err := <-reply:
+			return err
+		case <-stops:
+			cancel()
+			<-reply
+			return errOperatorDraining
+		case _, live := <-ticks:
+			if !live {
+				cancel()
+				<-reply
+				return errOperatorDraining
+			}
+			_ = e.touchRunner(c, h) // A busy writer leaves the real heartbeat deadline in force.
+		}
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_runner_test.go b/cmd/ailang-worldd/operator_finish_runner_test.go
new file mode 100644
index 0000000..4c05c86
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_runner_test.go
@@ -0,0 +1,81 @@
+package main
+
+import (
+	"errors"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorFinishReadiness(t *testing.T) {
+	for _, phase := range []string{"confirmed", "resume"} {
+		t.Run(phase, func(t *testing.T) {
+			var confirm []bool
+			if phase == "confirmed" {
+				confirm = []bool{true}
+			}
+			f, work, now := operatorLocalFinishFixture(t, nil, confirm...)
+			path, err := filepath.EvalSymlinks(f.root)
+			if err != nil {
+				t.Fatal(err)
+			}
+			tools := filepath.Join(path, "approved")
+			if err := os.Mkdir(tools, 0700); err != nil {
+				t.Fatal(err)
+			}
+			trace := filepath.Join(tools, "trace")
+			c := runnerConfig{Home: path, Credentials: work.manifest.Store, UID: os.Getuid(), ToolDirs: []string{tools},
+				Executables: make(map[string]string), FinishRemotes: []string{work.manifest.Remote.URL}}
+			writeFinishTestTool(t, &c, tools, "gh", "case \"$1\" in auth|label) ;; repo) printf 'WRITE\\n';; *) exit 83;; esac\n")
+			writeFinishTestTool(t, &c, tools, "git", "printf 'remote-read\\n' >> "+shellQuote(trace)+"\nexit 1\n")
+			before := serviceFileHash(work.manifest.Store)
+			ctx := boundedTestContext(t)
+			h, err := work.read.OperatorHeartbeat(ctx)
+			if err != nil {
+				t.Fatal(err)
+			}
+			old, err := work.read.OperatorAction(ctx, work.prep.row.ID, *now)
+			if err != nil {
+				t.Fatal(err)
+			}
+			err = (serviceEnv{}).consumeFinish(ctx, c, h, operatorStartSource{clock: func() int64 { return *now }})
+			row, readErr := work.read.OperatorAction(ctx, work.prep.row.ID, *now)
+			marker, markerErr := os.ReadFile(trace)
+			beforeHead := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
+			j := operatorFinishJournal(t, work)
+			if !errors.Is(err, operatorGitHubAuth(6)) || readErr != nil || !reflect.DeepEqual(old, row) ||
+				before != serviceFileHash(work.manifest.Store) || string(marker) != "remote-read\n" || markerErr != nil ||
+				beforeHead != f.before || j.CommitPhase != "" || j.OpenCheck || len(j.Checks) != 0 || f.checks != 0 ||
+				!h.Available(h.Build, "start", *now) {
+				t.Fatalf("readiness crossed %s admission: %v, row %s, marker=%q (%v)", phase, err, row.State, marker, markerErr)
+			}
+		})
+	}
+}
+
+func TestOperatorFinishEligibility(t *testing.T) {
+	_, work, now := operatorLocalFinishFixture(t, nil)
+	ctx := boundedTestContext(t)
+	cfg := runnerConfig{Credentials: work.manifest.Store}
+	src := operatorStartSource{clock: func() int64 { return *now }}
+	for _, episode := range []string{"ep1", "../ep1"} {
+		old := serviceFileHash(work.manifest.Store)
+		_, err := (serviceEnv{}).finishInputs(ctx, cfg, src, work.read, episode, nil)
+		if !errors.Is(err, errOperatorFinishEligibility) || serviceFileHash(work.manifest.Store) != old {
+			t.Fatalf("ineligible %q changed evidence: %v", episode, err)
+		}
+	}
+	// Reads after closing the handle use the same safe eligibility refusal.
+	ro, err := store.OpenCredentialsReadOnly(work.manifest.Store)
+	if err != nil {
+		t.Fatal(err)
+	}
+	ro.Close()
+	if _, err := (serviceEnv{}).finishInputs(ctx, cfg, src, ro, "ep1", nil); !errors.Is(err, errOperatorFinishEligibility) {
+		t.Fatalf("store error was not classified as eligibility: %v", err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_launch_test.go b/cmd/ailang-worldd/operator_launch_test.go
index 91b097a..e6fcecc 100644
--- a/cmd/ailang-worldd/operator_launch_test.go
+++ b/cmd/ailang-worldd/operator_launch_test.go
@@ -43,7 +43,12 @@ func newOperatorLaunchRig(t *testing.T, script string) *operatorLaunchRig {
 
 func newOperatorLaunchRigWithTask(t *testing.T, script, task string) *operatorLaunchRig {
 	t.Helper()
-	r, c, h, runner := operatorActivationFixtureWithTask(t, script, task)
+	return newOperatorLaunchRigWithRemote(t, script, task, "")
+}
+
+func newOperatorLaunchRigWithRemote(t *testing.T, script, task, remote string) *operatorLaunchRig {
+	t.Helper()
+	r, c, h, runner := operatorActivationFixtureWithRemote(t, script, task, remote)
 	f := &operatorLaunchRig{worktreeStartRig: r, c: c, h: h, runner: runner, ctx: activationTestContext(t), now: r.src.clock()}
 	r.src.clock = func() int64 { return f.now }
 	r.env.terminal.now = r.src.clock
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
index ddd3362..6f08902 100644
--- a/cmd/ailang-worldd/operator_runner.go
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -92,6 +92,18 @@ func runOperatorRunner(args []string, errOut io.Writer, e serviceEnv) int {
 	ticker := time.NewTicker(runnerPoll)
 	defer ticker.Stop()
 	source := operatorStartSource{defaultPresetsDir(), defaultExecProfilesDir(), "", defaultHostsDir(), e.terminal.now}
+	if len(c.FinishRemotes) != 0 {
+		ctx, cancel := e.operatorContext(30 * time.Second)
+		h.Readiness = "both"
+		for _, remote := range c.FinishRemotes {
+			if err := prepareOperatorFinish(ctx, c, remote); err != nil {
+				fmt.Fprintln(errOut, err)
+				h.Readiness = "start"
+				break
+			}
+		}
+		cancel()
+	}
 	e.operatorCatalogue = func(c runnerConfig) []byte { return e.readOperatorCatalogue(c, source) }
 	return e.runOperatorPreparationLoop(errOut, c, h, source, stops, ticker.C, lock)
 }
@@ -142,6 +154,27 @@ func (e serviceEnv) runOperatorPreparationLoop(errOut io.Writer, c runnerConfig,
 				fmt.Fprintln(errOut, "operator start blocked; inspect service runner status and reconcile durable evidence")
 			}
 		}
+		if h.Readiness == "both" {
+			err := e.runFinishTurn(c, &h, src, stops, ticks)
+			var auth operatorGitHubAuth
+			var busy operatorFinishBusy
+			switch {
+			case errors.Is(err, errOperatorDraining):
+				return 0
+			case errors.As(err, &auth):
+				h.Readiness = "start"
+				if err := e.touchRunner(c, &h); err != nil {
+					return serviceMessage(errOut, "runner halted", "heartbeat unavailable", "ailang-worldd service runner status", 2)
+				}
+				fmt.Fprintln(errOut, auth)
+			case errors.As(err, &busy):
+				fmt.Fprintln(errOut, busy.Error())
+			case errors.Is(err, errOperatorFinishReview):
+				fmt.Fprintln(errOut, errOperatorFinishReview)
+			case err != nil && !errors.Is(err, store.ErrOperatorMissing):
+				fmt.Fprintln(errOut, "operator finish blocked; inspect action status and durable evidence")
+			}
+		}
 		select {
 		case _, open := <-ticks:
 			if !open {
@@ -160,12 +193,18 @@ func (e serviceEnv) runnerTick(c runnerConfig, h *store.OperatorHeartbeat) error
 	if err := e.checkRunnerPair(c); err != nil {
 		return err
 	}
-	h.Readiness = "start"
+	if h.Readiness == "" {
+		h.Readiness = "start"
+	}
 	if e.operatorCatalogue != nil {
 		h.Catalogue = e.operatorCatalogue(c)
 	} else {
 		h.Catalogue, _ = json.Marshal(c.PresetNames)
 	}
+	return e.touchRunner(c, h)
+}
+
+func (e serviceEnv) touchRunner(c runnerConfig, h *store.OperatorHeartbeat) error {
 	h.SeenAt = e.terminal.now()
 	return runnerWrite(c.Credentials, func(db *store.Credentials) error {
 		ctx, cancel := e.operatorContext(time.Second)
@@ -233,7 +272,9 @@ func (e serviceEnv) runnerStatus(w io.Writer) error {
 	}
 	fmt.Fprintln(w, "start: ready")
 	if !h.Available(build, "finish", now) {
-		fmt.Fprintln(w, "finish blocked: github-auth; repair PATH, GitHub login or Git credential helper in runner setup")
+		fmt.Fprintln(w, "finish blocked: github-auth; "+operatorFinishRemedy)
+	} else {
+		fmt.Fprintln(w, "finish: ready")
 	}
 	return nil
 }
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index 4230030..d4f697d 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -21,7 +21,7 @@ import (
 
 const runnerLabel = "com.sunholo.ailang-world.operator-runner"
 const runnerSystemPath = "/usr/bin:/bin:/usr/sbin:/sbin"
-const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | reconcile-files ACTION-ID | reconcile-start ACTION-ID | launch-status ACTION-ID | uninstall
+const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | reconcile-files ACTION-ID | reconcile-start ACTION-ID | launch-status ACTION-ID | finish-status ACTION-ID | uninstall
 Install, stop, cancel, orphan mint revocation, worktree and file disposition and uninstall require typed y. Start uses the approved configuration.
 Config JSON keys: home, uid, address, world, credentials, executables (absolute path to sha256 digest), preset_names, tool_dirs, backend ("headless").
 Optional finish readiness keys: finish_remotes (approved GitHub URLs), finish_ssh_socket (approved login agent socket), finish_https_helper ("gh" or "osxkeychain"). Pin gh, git and any separate HTTPS helper.
@@ -291,7 +291,7 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 		if len(args) != 3 || args[1] != "--config" {
 			return refuse(errors.New("expected install --config FILE"))
 		}
-	case "cancel", "reconcile-mint", "reconcile-worktree", "reconcile-files", "reconcile-start", "launch-status":
+	case "cancel", "reconcile-mint", "reconcile-worktree", "reconcile-files", "reconcile-start", "launch-status", "finish-status":
 		if len(args) != 2 || !verbs.ValidOperatorActionID(args[1]) {
 			return refuse(fmt.Errorf("expected %s ACTION-ID", verb))
 		}
@@ -314,6 +314,12 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 	if verb == "reconcile-files" {
 		return e.reconcileRunnerFiles(args[1], out, errOut)
 	}
+	if verb == "finish-status" {
+		if err := e.printFinishStatus(args[1], out); err != nil {
+			return refuse(err)
+		}
+		return exitOK
+	}
 	if verb == "launch-status" {
 		if err := e.operatorLaunchStatus(args[1], out); err != nil {
 			return refuse(err)
diff --git a/cmd/ailang-worldd/operator_worktree_test.go b/cmd/ailang-worldd/operator_worktree_test.go
index ae079e2..5cb3434 100644
--- a/cmd/ailang-worldd/operator_worktree_test.go
+++ b/cmd/ailang-worldd/operator_worktree_test.go
@@ -45,12 +45,20 @@ func newWorktreeStartRigWithAgent(t *testing.T, claude bool, ctx context.Context
 }
 
 func newWorktreeStartRigWithTask(t *testing.T, claude bool, ctx context.Context, body, task string, ttl ...int64) *worktreeStartRig {
+	t.Helper()
+	return newWorktreeStartRigWithRemote(t, claude, ctx, body, task, "", ttl...)
+}
+
+func newWorktreeStartRigWithRemote(t *testing.T, claude bool, ctx context.Context, body, task, remote string, ttl ...int64) *worktreeStartRig {
 	t.Helper()
 	s := operatorService(t)
 	s.mustInstall(t)
 	s.writes = nil
 	s.env.terminal.now = func() int64 { return 1012 }
 	f := newPresetFixture(t, ctx)
+	if remote != "" {
+		presetGit(t, f.r.f, "remote", "set-url", "origin", remote)
+	}
 	// The runner checks the agent and session directory up to its login home.
 	bin, err := os.MkdirTemp(s.config.Home, "agent-")
 	if err != nil {
@@ -79,6 +87,10 @@ func newWorktreeStartRigWithTask(t *testing.T, claude bool, ctx context.Context,
 		if claude {
 			p.Agent = "claude"
 		}
+		if remote != "" {
+			p.Committer = "Fixture Operator <fixture@example.invalid>"
+			p.Checks = p.Checks[:1]
+		}
 		if len(ttl) != 0 {
 			p.TTLS = ttl[0]
 		}
diff --git a/host/daemon/surface_operator.go b/host/daemon/surface_operator.go
index 37d5db3..ba635ad 100644
--- a/host/daemon/surface_operator.go
+++ b/host/daemon/surface_operator.go
@@ -21,7 +21,42 @@ import (
 const operatorPath = "/surface/operator-actions"
 const operatorUnavailable = "runner or reviewed preset unavailable; inspect ailang-worldd service runner status or start at a terminal"
 
-var operatorPolicy = verbs.OperatorPolicy{Start: "A", Finish: "A"} // Headless start is the only installed S7 backend.
+var operatorPolicy = verbs.OperatorPolicy{Start: "A", Finish: "A"}
+
+var errOperatorFinishAuth = errors.New("finish blocked: github-auth; repair PATH, GitHub login or Git credential helper in attended service-runner setup and re-run readiness")
+
+func supportedOperatorAction(in *verbs.OperatorIntent) bool {
+	return in != nil && (in.Kind == "start" && in.Headless || in.Kind == "finish")
+}
+
+func (d *Daemon) operatorActionForIntent(ctx context.Context, db *store.Credentials, in *verbs.OperatorIntent, task string, now int64) (store.OperatorHeartbeat, approval.OperatorAction, error) {
+	if in.Kind == "start" {
+		h, entry, err := d.operatorReview(ctx, db, in.Preset, now)
+		if err != nil {
+			return h, approval.OperatorAction{}, err
+		}
+		a, err := entry.Action(in.Episode, task, h.CatalogueGeneration)
+		return h, a, err
+	}
+	h, err := db.OperatorHeartbeat(ctx)
+	build := CurrentBuild().ExecutableSHA256
+	if err != nil || build == "" || !h.Available(build, "start", now) || h.CatalogueGeneration <= 0 {
+		return h, approval.OperatorAction{}, store.ErrOperatorRunner
+	}
+	if !h.Available(build, "finish", now) {
+		return h, approval.OperatorAction{}, errOperatorFinishAuth
+	}
+	entry, err := h.FinishReview(in.Episode, in.Preset)
+	served, pathErr := store.CanonicalPath(d.credentialsPath)
+	if err != nil || pathErr != nil || entry.Store != served {
+		return h, approval.OperatorAction{}, store.ErrOperatorRunner
+	}
+	if len(entry.Protected) != 0 || len(entry.Unattributed) != 0 {
+		return h, approval.OperatorAction{}, errors.New("protected or unattributed changes need terminal finish")
+	}
+	a, err := entry.Action(h.CatalogueGeneration, in.Ready)
+	return h, a, err
+}
 
 func (d *Daemon) operatorVerb(r *http.Request, argv []string, out *verbAnswer) bool {
 	a := verbs.AdmitOperator(argv, operatorPolicy)
@@ -29,24 +64,30 @@ func (d *Daemon) operatorVerb(r *http.Request, argv []string, out *verbAnswer) b
 	case verbs.OperatorNever:
 		return out.Class == string(verbs.Never) // Lookup already supplied the refusal.
 	case verbs.OperatorReview:
-		if a.Intent.Kind == "start" && a.Intent.Headless {
-			out.Outcome, out.Reason, out.CLI = "terminal", operatorUnavailable, a.CLI
-			ctx, cancel := d.readCtx(r)
-			defer cancel()
-			ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
-			if err != nil {
-				return true
-			}
+		if !supportedOperatorAction(a.Intent) {
+			return false
+		}
+		ctx, cancel := d.readCtx(r)
+		defer cancel()
+		ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+		if err == nil {
 			defer ro.Close()
-			h, entry, err := d.operatorReview(ctx, ro, a.Intent.Preset, d.approveNow())
-			if err == nil {
-				if action, err := entry.Action(a.Intent.Episode, "", h.CatalogueGeneration); err == nil {
-					out.Outcome, out.Reason, out.Argv, out.Text = "review", "", action.Argv, action.Review
-					out.CatalogueGeneration = h.CatalogueGeneration
-				}
+			h, action, err := d.operatorActionForIntent(ctx, ro, a.Intent, "", d.approveNow())
+			switch {
+			case err == nil:
+				out.Outcome, out.Reason, out.Argv, out.Text = "review", "", action.Argv, action.Review
+				out.CatalogueGeneration = h.CatalogueGeneration
+				return true
+			case errors.Is(err, errOperatorFinishAuth):
+				out.Outcome, out.Reason, out.CLI = "terminal", errOperatorFinishAuth.Error(), a.CLI
+				return true
 			}
-			return true
 		}
+		if a.Intent.Kind == "finish" {
+			return false // The existing route supplies its terminal fallback.
+		}
+		out.Outcome, out.Reason, out.CLI = "terminal", operatorUnavailable, a.CLI
+		return true
 	}
 	return false
 }
@@ -158,13 +199,13 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
 	}
 	var body operatorRequestBody
 	if approval.DecodeOperatorObject(raw, &body, "argv", "operationId", "catalogueGeneration", "task") != nil ||
-		len(body.Task) == 0 || len(body.Task) > approval.MaxOperatorTaskBytes || strings.ContainsRune(body.Task, 0) {
+		len(body.Task) > approval.MaxOperatorTaskBytes || strings.ContainsRune(body.Task, 0) {
 		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send argv, operationId, catalogueGeneration and a task of 1–8192 UTF-8 bytes")
 		return
 	}
 	a := verbs.AdmitOperator(body.Argv, operatorPolicy)
-	if a.Class != verbs.OperatorReview || a.Intent.Kind != "start" || !a.Intent.Headless {
-		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "review a supported headless start first")
+	if a.Class != verbs.OperatorReview || !supportedOperatorAction(a.Intent) || (a.Intent.Kind == "start") != (body.Task != "") {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "review a supported headless start with a task or a finish without a task")
 		return
 	}
 	if body.CatalogueGeneration <= 0 || !store.ValidOperatorOperationID(body.OperationID) {
@@ -189,8 +230,13 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
 			return
 		}
 		frozen, err := approval.DecodeOperatorAction(row.Action)
-		if err != nil || !slices.Equal(frozen.Argv, a.Intent.Argv) || frozen.Task != body.Task ||
-			!slices.Contains(frozen.Resolution, approval.OperatorResolved{Name: "mint_catalogue", Value: strconv.FormatInt(body.CatalogueGeneration, 10)}) {
+		argv, fact := a.Intent.Argv, "mint_catalogue"
+		if a.Intent.Kind == "finish" {
+			argv = approval.OperatorArgv("finish", a.Intent.Episode, approval.OperatorOptions{Preset: a.Intent.Preset, ConfirmOnSurface: true, Ready: a.Intent.Ready})
+			fact = "finish_catalogue"
+		}
+		if err != nil || !slices.Equal(frozen.Argv, argv) || frozen.Task != body.Task ||
+			!slices.Contains(frozen.Resolution, approval.OperatorResolved{Name: fact, Value: strconv.FormatInt(body.CatalogueGeneration, 10)}) {
 			d.operatorError(w, store.ErrOperatorConflict)
 			return
 		}
@@ -204,7 +250,7 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
 	}
 	defer wr.Close()
 	now = d.approveNow()
-	h, entry, err := d.operatorReview(ctx, wr, a.Intent.Preset, now)
+	h, action, err := d.operatorActionForIntent(ctx, wr, a.Intent, body.Task, now)
 	if err == nil && h.CatalogueGeneration != body.CatalogueGeneration {
 		err = store.ErrOperatorRunner
 	}
@@ -212,11 +258,6 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
 		d.operatorError(w, err)
 		return
 	}
-	action, err := entry.Action(a.Intent.Episode, body.Task, h.CatalogueGeneration)
-	if err != nil {
-		d.operatorError(w, err)
-		return
-	}
 	var id [16]byte
 	var challenge [32]byte
 	if _, err := rand.Read(id[:]); err != nil {
@@ -239,6 +280,8 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
 
 func (d *Daemon) operatorError(w http.ResponseWriter, err error) {
 	switch {
+	case errors.Is(err, errOperatorFinishAuth):
+		writeApprovalRefusal(w, http.StatusConflict, "RunnerUnavailable", errOperatorFinishAuth.Error())
 	case isOlderStore(err), errors.Is(err, store.ErrOperatorRunner):
 		writeApprovalRefusal(w, http.StatusServiceUnavailable, "RunnerUnavailable", operatorUnavailable+"; review again")
 	case store.IsWriterAlreadyActive(err):
diff --git a/host/daemon/surface_operator_confirm.go b/host/daemon/surface_operator_confirm.go
index 60e3f8f..95683f8 100644
--- a/host/daemon/surface_operator_confirm.go
+++ b/host/daemon/surface_operator_confirm.go
@@ -93,7 +93,8 @@ func (d *Daemon) handleOperatorConfirm(w http.ResponseWriter, r *http.Request, i
 		d.operatorError(w, pathErr)
 		return
 	}
-	if err != nil || a.Kind != "start" || !a.Options.Headless || a.ID != id || a.Episode != row.Episode || a.Kind != row.Kind || a.Store != served ||
+	supported := a.Kind == "start" && a.Options.Headless || a.Kind == "finish" && a.Options.ConfirmOnSurface
+	if err != nil || !supported || a.ID != id || a.Episode != row.Episode || a.Kind != row.Kind || a.Store != served ||
 		!bytes.Equal(row.Digest, digest[:]) || body.Digest != hex.EncodeToString(row.Digest) || len(row.Challenge) != 32 {
 		d.operatorError(w, store.ErrOperatorConflict)
 		return
diff --git a/host/daemon/surface_operator_finish_test.go b/host/daemon/surface_operator_finish_test.go
new file mode 100644
index 0000000..79df813
--- /dev/null
+++ b/host/daemon/surface_operator_finish_test.go
@@ -0,0 +1,138 @@
+package daemon
+
+import (
+	"context"
+	"crypto/sha256"
+	"encoding/hex"
+	"encoding/json"
+	"net/http"
+	"net/http/httptest"
+	"net/url"
+	"os"
+	"reflect"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorFinishAdmission(t *testing.T) {
+	f := newOperatorHTTP(t, 3*time.Minute)
+	reader, err := os.ReadFile("../../docs/QUICKSTART.md")
+	if err != nil {
+		t.Fatal(err)
+	}
+	_, usage, ok := strings.Cut(string(reader), "### Prepare a local finish in the World app")
+	_, usage, found := strings.Cut(usage, "```text\n")
+	command, _, closed := strings.Cut(usage, "\n```")
+	if !ok || !found || !closed {
+		t.Fatal("QUICKSTART finish fixture is absent")
+	}
+	query := "/surface/verb.json?q=" + url.QueryEscape(command)
+	path, err := store.CanonicalPath(f.d.credentialsPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	manifest := store.OperatorFinishReview{Store: path, Episode: "ep", Preset: "p",
+		Report: "private-manifest-fixture-marker", Commit: store.OperatorLocalCommit{Message: "signed local finish"}}
+	catalogue := map[string]json.RawMessage{}
+	if err := json.Unmarshal(f.h.Catalogue, &catalogue); err != nil {
+		t.Fatal(err)
+	}
+	catalogue[store.OperatorFinishCatalogueKey] = opJSON(map[string]store.OperatorFinishReview{"ep": manifest})
+	f.h.Catalogue = opJSON(catalogue)
+	f.publish()
+	get := func(path string) *verbAnswer {
+		t.Helper()
+		response := f.do("GET", path, nil, nil)
+		var answer verbAnswer
+		if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &answer) != nil {
+			t.Fatalf("read %d %s", response.Code, response.Body.String())
+		}
+		return &answer
+	}
+	before := f.evidence()
+	if got := get(query); got.Outcome != "terminal" || !strings.Contains(got.Reason, "blocked: github-auth") || f.evidence() != before {
+		t.Fatalf("finish readiness refusal: %+v", got)
+	}
+	f.h.Readiness = "both"
+	f.publish()
+	before = f.evidence()
+	review := get(query)
+	if review.Outcome != "review" || !strings.Contains(review.Text, manifest.Report) || f.evidence() != before {
+		t.Fatalf("review was unavailable or effected a write: %+v", review)
+	}
+
+	post := func(path string, data []byte) *httptest.ResponseRecorder { return f.do("POST", path, data, nil) }
+	req := operatorRequestBody{Argv: review.Argv, OperationID: opRequest(35).OperationID, CatalogueGeneration: review.CatalogueGeneration}
+	wrong := req
+	wrong.Task = "no task belongs to prep"
+	if got := post(operatorPath, opJSON(wrong)); got.Code != 400 || !strings.Contains(got.Body.String(), "BadRequest") || f.evidence() != before {
+		t.Fatalf("task refusal changed state: %+v", got)
+	}
+	wrong = req
+	wrong.CatalogueGeneration--
+	if got := post(operatorPath, opJSON(wrong)); got.Code != 503 || !strings.Contains(got.Body.String(), "RunnerUnavailable") || f.evidence() != before {
+		t.Fatalf("stale catalogue changed state: %+v", got)
+	}
+
+	key := approvaltest.NewPasskey(approval.AlgEdDSA, 43)
+	f.enrolPasskey("finish-key", key, opOrigin)
+	var requested operatorStatus
+	if got := post(operatorPath, opJSON(req)); got.Code != 200 || json.Unmarshal(got.Body.Bytes(), &requested) != nil ||
+		requested.Kind != "finish" || requested.State != "requested" {
+		t.Fatalf("request %+v", got)
+	}
+	row := f.row(requested.ID)
+	for _, path := range []string{"/surface/requests.json", "/surface/home-summary.json"} {
+		routine := f.do("GET", path, nil, nil)
+		if routine.Code != 200 || strings.Contains(routine.Body.String(), manifest.Report) || strings.Contains(routine.Body.String(), "artifact") ||
+			strings.Contains(routine.Body.String(), opTask) {
+			t.Fatalf("routine action list exposed protected payloads: %s HTTP %d", path, routine.Code)
+		}
+	}
+	if got := post(operatorPath, opJSON(req)); got.Code != 200 || !strings.Contains(got.Body.String(), requested.ID) {
+		t.Fatalf("retry lost same request: %+v", got)
+	}
+	artifact := approveBody(hex.EncodeToString(row.Digest),
+		key.Assert(opOrigin, "localhost", sha256.Sum256(opMessage(t, [32]byte(row.Digest), [32]byte(row.Challenge), row.ExpiresAt)), 2))
+	before = f.evidence()
+	if got := post(requested.Href+"/confirm", sshBody(hex.EncodeToString(row.Digest), []byte("y"))); got.Code != 422 || !strings.Contains(got.Body.String(), "NotApproved") || f.evidence() != before {
+		t.Fatalf("plain y changed state: %+v", got)
+	}
+	if got := post(requested.Href+"/confirm", artifact); got.Code != 200 || !strings.Contains(got.Body.String(), "confirmed") {
+		t.Fatalf("key refused: %+v", got)
+	}
+	first := f.evidence()
+	if got := post(requested.Href+"/confirm", artifact); got.Code != 200 || !strings.Contains(got.Body.String(), "already-confirmed") || f.evidence() != first {
+		t.Fatalf("replay changed durable state: %+v", got)
+	}
+	if e := f.evidence(); e.actions != 1 || e.operations != 1 || e.sessions != before.sessions || !reflect.DeepEqual(e.world, before.world) {
+		t.Fatalf("unexpected prep-route effects %+v", e)
+	}
+	row = f.row(row.ID)
+	cas := store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
+	sha := strings.Repeat("e", 40)
+	f.write(func(ctx context.Context, db *store.Credentials) error {
+		if err := db.ClaimOperatorAction(ctx, cas, f.h.Owner, apT0); err != nil {
+			return err
+		}
+		cas.Generation++
+		step := store.OperatorStepKey{Action: cas, Owner: f.h.Owner.Nonce, Kind: "commit", Identity: row.ID + ":" + sha, Version: 1}
+		if err := db.ReserveOperatorStep(ctx, step, apT0); err != nil {
+			return err
+		}
+		return db.ObserveOperatorStep(ctx, step, opJSON(store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: sha}))
+	})
+	response := f.do("GET", requested.StatusHref, nil, nil)
+	var status operatorStatus
+	if response.Code != 200 || json.Unmarshal(response.Body.Bytes(), &status) != nil || status.LocalCommit != sha ||
+		status.ConfirmHref != "" || !strings.Contains(status.Reason, "PR confirmation has not been issued") ||
+		strings.Contains(response.Body.String(), manifest.Report) || strings.Contains(f.log.String(), manifest.Report) ||
+		strings.Contains(response.Body.String(), "authenticatorData") || strings.Contains(f.log.String(), "authenticatorData") || f.row(row.ID).V3ID != "" {
+		t.Fatalf("local-only status %d %+v", response.Code, status)
+	}
+}
diff --git a/host/daemon/surface_operator_read.go b/host/daemon/surface_operator_read.go
index 47ad635..d70a51a 100644
--- a/host/daemon/surface_operator_read.go
+++ b/host/daemon/surface_operator_read.go
@@ -15,6 +15,9 @@ import (
 )
 
 type operatorStatus struct {
+	Checks      []store.OperatorFinishCheck `json:"checks,omitempty"`
+	LocalCommit string                      `json:"localCommit,omitempty"`
+
 	V            int    `json:"v"`
 	ID           string `json:"id"`
 	Episode      string `json:"episode"`
@@ -62,7 +65,7 @@ var operatorCard = template.Must(template.New("operator action").Parse(`<!doctyp
 <h1 id="operator-title" tabindex="-1">Operator action {{.Status.ID}}</h1>
 <p id="operator-state" data-state="{{.Status.State}}">{{.Status.State}}</p>
 <p id="operator-progress" role="status">Status at {{.Status.Now}}{{if .Status.Reason}} · {{.Status.Reason}}{{end}}</p>
-{{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre><h2>Task</h2><pre id="operator-task">{{.Task}}</pre>{{end}}
+{{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre>{{end}}{{if .Task}}<h2>Task</h2><pre id="operator-task">{{.Task}}</pre>{{end}}
 <section id="operator-controls" hidden aria-label="Confirm this action">
 <button id="operator-key" type="button" hidden>Confirm with passkey</button>
 <details><summary>Confirm with SSH signature</summary>
@@ -117,6 +120,29 @@ func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id s
 			return
 		}
 	}
+	if status.Href != "" && row.Kind == "finish" && row.Owner != "" {
+		j, err := ro.ReadOperatorFinishJournal(ctx, store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation})
+		if err != nil {
+			d.operatorError(w, err)
+			return
+		}
+		status.Checks = j.Checks
+		switch {
+		case j.Commit.State == store.OperatorLocalCommitDone:
+			status.LocalCommit = j.Commit.SHA
+			status.Reason = "Local commit " + j.Commit.SHA + " retained; PR confirmation has not been issued"
+		case j.OpenCheck || j.CommitPhase == "reserved":
+			status.Reason = "Finish outcome uncertain; inspect durable receipts"
+		case j.Commit.State == "commit-failed":
+			status.Reason = "Local commit failed; check receipts retained"
+		case len(j.Checks) > 0:
+			status.Reason = "Sandboxed check results recorded; see Action status (JSON)"
+		}
+	}
+	if row.Kind == "finish" && status.Reason == "runner disconnected; showing durable action status" &&
+		h.Available(CurrentBuild().ExecutableSHA256, "start", now) && !h.Available(CurrentBuild().ExecutableSHA256, "finish", now) {
+		status.Reason = errOperatorFinishAuth.Error()
+	}
 	if statusJSON {
 		writeApprovalJSON(w, http.StatusOK, status)
 		return
diff --git a/host/surfacepage/app/app.js b/host/surfacepage/app/app.js
index 0481d37..4249764 100644
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -1437,24 +1437,26 @@
   const taskValid = task => typeof task === 'string' && task.length > 0 && !task.includes('\0') &&
     new TextEncoder().encode(task).length <= 8192;
   function reviewValid(argv, catalogueGeneration) {
-    return Array.isArray(argv) && argv.length === 8 &&
-      argv[0] === 'ailang-worldd' && argv[1] === 'episode' && argv[2] === 'start' &&
-      typeof argv[3] === 'string' && EPISODE_ID.test(argv[3]) && argv[4] === '--preset' &&
-      typeof argv[5] === 'string' && PRESET_ID.test(argv[5]) && argv[5] !== 'se-tools' &&
-      argv[6] === '--launch' && argv[7] === '--headless' &&
-      Number.isSafeInteger(catalogueGeneration) && catalogueGeneration > 0;
-  }
+    if (!Array.isArray(argv) || argv[0] !== 'ailang-worldd' || argv[1] !== 'episode' ||
+        typeof argv[3] !== 'string' || !EPISODE_ID.test(argv[3]) || argv[4] !== '--preset' ||
+        typeof argv[5] !== 'string' || !PRESET_ID.test(argv[5]) || argv[5] === 'se-tools' ||
+        !Number.isSafeInteger(catalogueGeneration) || catalogueGeneration <= 0) return false;
+    if (argv[2] === 'start') return argv.length === 8 && argv[6] === '--launch' && argv[7] === '--headless';
+    return argv[2] === 'finish' && argv[6] === '--confirm-on-surface' &&
+      (argv.length === 7 || (argv.length === 8 && argv[7] === '--ready'));
+  }
+  const operatorTaskValid = (argv, task) => argv[2] === 'finish' ? task === '' : taskValid(task);
   function operatorDraft(storage) {
     try {
       const d = JSON.parse(storage.getItem(OPERATOR_DRAFT));
       const r = d?.request;
       if (typeof d?.body === 'string' && d.body && r && reviewValid(r.argv, r.catalogueGeneration) &&
-          typeof r.operationId === 'string' && OPERATION_ID.test(r.operationId) && taskValid(r.task)) return d;
+          typeof r.operationId === 'string' && OPERATION_ID.test(r.operationId) && operatorTaskValid(r.argv, r.task)) return d;
     } catch (_) { /* page stays readable */ }
     return null;
   }
-  function operatorCardHref(b, episode) {
-    if (!b || b.v !== 1 || b.kind !== 'start' || b.episode !== episode || typeof b.id !== 'string' || !OPERATOR_ID.test(b.id) ||
+  function operatorCardHref(b, episode, kind = 'start') {
+    if (!b || b.v !== 1 || b.kind !== kind || b.episode !== episode || typeof b.id !== 'string' || !OPERATOR_ID.test(b.id) ||
         !['requested', 'confirmed', 'claimed', 'running', 'succeeded', 'failed', 'expired', 'cancelled', 'uncertain'].includes(b.state) ||
         ![b.now, b.createdAt, b.expiresAt].every(Number.isSafeInteger)) return '';
     const path = OPERATOR_PATH + '/' + b.id;
@@ -1466,7 +1468,7 @@
     return { get request() { return request; }, async submit(task) {
       if (busy || !io.current()) return {};
       if (!request) {
-        if (!taskValid(task)) return { error: 'Enter a task of 1–8192 UTF-8 bytes without NUL.' };
+        if (!operatorTaskValid(panel.argv, task)) return { error: 'Enter a task of 1–8192 UTF-8 bytes without NUL.' };
         const operationId = Array.from(io.random(new Uint8Array(16)), n => n.toString(16).padStart(2, '0')).join('');
         request = { argv: panel.argv, operationId, catalogueGeneration: panel.catalogueGeneration, task };
       }
@@ -1481,7 +1483,7 @@
           redirect: 'error', credentials: 'omit', signal: io.timeout(VERB_ROUTE_BUDGET_MS), body: JSON.stringify(request)});
         const b = await res.json();
         if (!io.current()) return {};
-        const href = res.status === 200 && operatorCardHref(b, request.argv[3]);
+        const href = res.status === 200 && operatorCardHref(b, request.argv[3], request.argv[2]);
         if (!href) {
           const reason = typeof b?.error === 'string' && typeof b?.reason === 'string' ? b.error + ': ' + b.reason + ' ' : '';
           return { error: reason + 'Request outcome unavailable. Use Reconnect to read or retry this request.' };
@@ -3546,7 +3548,9 @@
     label.textContent = 'Task (1–8192 UTF-8 bytes)'; label.htmlFor = task.id;
     request.id = 'operator-request'; request.type = 'button';
     info.id = 'operator-request-info'; info.setAttribute('role', 'status');
-    area.append(label, document.createElement('br'), task, document.createElement('br'), request, info);
+    const finish = panel.argv[2] === 'finish';
+    if (!finish) area.append(label, document.createElement('br'), task, document.createElement('br'));
+    area.append(request, info);
     const live = () => current() && el.isConnected;
     const control = operatorRequestControl(panel, {storage: operatorStorage, current: live, fetch: (...args) => fetch(...args),
       random: b => crypto.getRandomValues(b), timeout: ms => AbortSignal.timeout(ms)});
@@ -3566,9 +3570,9 @@
       });
       repaint();
       request.disabled = false;
-      if (live()) (task.readOnly ? request : task).focus();
+      if (live()) (finish || task.readOnly ? request : task).focus();
     });
-    (task.readOnly ? request : task).focus();
+    (finish || task.readOnly ? request : task).focus();
   }
   const outstandingOperatorRequest = operatorDraft(operatorStorage);
   if (outstandingOperatorRequest) {
diff --git a/host/surfacepage/app_operator_test.go b/host/surfacepage/app_operator_test.go
index 183b1e3..74d07c3 100644
--- a/host/surfacepage/app_operator_test.go
+++ b/host/surfacepage/app_operator_test.go
@@ -186,3 +186,25 @@ for (const [k,v] of Object.entries({v:'1',episode:'other',kind:'finish',state:'i
   assert.equal(pick({...result,[k]:v}),'',k);
 `)
 }
+
+func TestOperatorFinishBarReview(t *testing.T) {
+	runOperatorBar(t, `
+const finishArgv = ['ailang-worldd','episode','finish','s8','--preset','p','--confirm-on-surface'];
+for(const bad of [finishArgv.slice(0,6), finishArgv.concat('--title'), finishArgv.concat('--ready','--ready')])
+  assert.equal(review({argv:bad}).review, undefined);
+const admitted = review({argv:finishArgv.concat('--ready')}).review;
+assert.deepEqual(admitted.argv, finishArgv.concat('--ready'));
+assert.equal(posts.length,0);
+const startCase = app.operatorRequestControl(admitted,io);
+assert.ok((await startCase.submit('not a task')).error);
+assert.equal(posts.length,0);
+assert.ok((await startCase.submit('')).error); // A start card answers a finish request.
+assert.deepEqual(posts[0].argv, finishArgv.concat('--ready'));
+assert.equal(posts[0].task,'');
+response = reply({...result, kind:'finish'});
+const reconnected = app.operatorRequestControl(admitted,io);
+assert.equal((await reconnected.submit('')).href,href);
+assert.deepEqual(posts[0],posts[1]);
+assert.equal(slot,null);
+`)
+}
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index dde5280..27ea2ee 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -2184,3 +2184,22 @@ choosing **Confirm with passkey**. The card shows the recorded confirmation and
 Close the page and reopen the same card URL to read the outcome and **Episode report**. Use
 **Refresh status** to reconnect if status reads fail. **Back to World** returns to the World surface.
 If the request response was lost, **Reconnect** in the review panel retries the same frozen request.
+
+### Prepare a local finish in the World app
+
+Use the Sol finish fixture with its explicit committer and fake GitHub URL.
+The fake Git tools read a bare local repository. With a stopped launch and
+one World-attributed worktree change, enter this fixture command in the bar.
+
+```text
+episode finish ep --preset p
+```
+
+Enter opens the review. Select **Request action**, review the held inputs on
+the card and confirm with the enrolled fixture passkey or SSH key. Use
+**Action status (JSON)** to read check outcomes and the retained local commit.
+The card reports that PR confirmation has yet to be issued.
+
+If the runner shows `finish blocked: github-auth`, repair its approved PATH,
+GitHub login or Git credential helper in the attended runner setup. The start
+flow remains available.
````

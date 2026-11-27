# w-world-address-bar-verbs-v4 — S6 LAUNCH sprint plan

## 0. Scope and outcome

Row 243, iteration 444, executor `codex:gpt-6.2-sol`. Reviewed design r1, quorum r2 PROCEED 3/3. D-WORLD-119 A and D-WORLD-120 A were attended at `82f05004`. S5f4 #689 merged at `c3504c56`. S6 picks up its `prepared-awaiting-S6` checkpoint. Design §8 budgets 460 + 430 = 890 raw, 1,202 calibrated.

**Outcome A by size.** The complete compiled prototype measures **1,650 changed lines** (+1,604, −46). Use the safe-unit exception because 1,650 exceeds 1,400 and fits 1,650. The gate and mutation evidence below controls execution readiness. The filter is `git diff -M --numstat b2b0ebab7bffbcf5f3105aa17e3ced452e480e35 HEAD -- cmd host docs app`. No renames occur. The planner measured each tracked edit or new path using the same Git numstat algorithm over pristine and cumulative archives, and checked the complete archive for out-of-filter product edits. The unfiltered product sum is **1,650**. The two untracked planning files are counted separately in the sprint JSON.

M1 installs the durable law and its pins. M2 ships the installed launch caller, the bounded worker and its process pins. The runner starts one approved headless Codex child, under the existing confinement and the new write-ahead reservation, without a browser route. Browser confirmation belongs to S7. F1/finish and herdr stay in their reviewed sprints. No schema migration is needed because schema 10 already admits the `launch` step, versioned observations and terminal action states. No `.ail` work is in this sprint.

## 1. Measurements

### 1.1 Base and provenance

- `git rev-parse HEAD` returned `b2b0ebab7bffbcf5f3105aa17e3ced452e480e35`. The initial `git status --short` was empty. `git log --oneline -14` recorded #691 followed by the documentation record at base. See `$S/base-history.raw`.
- `grep -nF 'S6 and later' design_docs/planned/w-world-address-bar-verbs-v4s5f4-sprint-plan.md` returned line 540. That row assigns agent launch and reservation to a later reviewed sprint. Iteration 442 controller readings and evaluator rounds were read as prior evidence. Iteration 444 measurements are recorded independently.
- `cat "$HOME/.ailang/state/world-iter429/e3.out"` and `find "$HOME/.ailang/state/world-iter429" -iname '*e3*'` report that the supplied directory is absent on this machine. The controller's checked-in readings report E3 PASS, twice-run synthetic child, two identities, the explicit HOME/PATH/UID and absent synthetic state. The E3 scope is the fake GUI supervisor. [Iteration 429 controller readings](../verification/world-iter429/readings.md)
- The controller baseline script emits `ALL_DONE for b2b0ebab7bffbcf5f3105aa17e3ced452e480e35` as its last line. §6 accepts this exact form and the bare marker. Policy is fail-set comparison for complete logs and exclusion for denied or truncated legs.

### 1.2 Guard inventory

Source commands are in `$S/anchors.sh`. `$S/guard-sources.raw` records the sweep over every CLI test file and every verifygate test file. The in-tree controls are `TestSubprocessSinkGate/detector_controls`, `TestNoParallelWriteForkPackagesOutsideVerifygate`'s known parallel write/fork case, the `TestExecGateArmsAreLoadBearing` verifier controls, and the named writer-site positive read-only census. Do not treat a guard kill as a behavioral mutation kill.

| Guard, base anchor | Selector and S6 disposition |
|---|---|
| `host/daemon/daemon_test.go:899` `TestDaemonDependencyAllowlist` | `inv-daemon`. No new daemon imports. |
| `host/boundary/allowlist_world_test.go:758` `TestWorldBoundaryDependencyAllowlist` | `inv-boundary`. The new store uses standard data and existing approval types. |
| `host/store/context_roots_test.go:136` `TestProductionContextRoots` | `context-roots`, `go test ./host/store/ -run '^TestProductionContextRoots$'`. Roots already named in `budgetContext` and `serveSignalContext` stay at their existing declarations. New launch calls use these finite seams. |
| `host/store/credentials_test.go:277` `TestCredentialsMethodsRequireADeadline` | `inv-store`. Pins all five added credentials methods against an unbounded context. |
| `host/store/deadline_guard_test.go:43` `TestStoreSurfaceTakesContext` | `inv-store`. Scans the separate World `Store` type. The launch ledger uses `Credentials`. |
| `host/boundary/allowlist_world_test.go:1323` `TestBoundaryASTWriteGuard` | `inv-boundary`. Store SQL stays under `host/store`. Run the suite's positive controls. |
| `host/verbs/operator_policy_test.go:69` `TestOperatorNeverDispatch` | `inv-verbs`. Add the new service status spelling to the existing `OperatorNever` family and the unregistered private worker to `OperatorOther`. No production grammar change. |
| `cmd/ailang-worldd/cli_credential_writers_test.go:18` `TestCLICredentialWriterSites` | `inv-cli`. The worker owns one named writer from final transactional authority through `Start` and birth capture. Old CLI writer sites stay in the census. |
| `cmd/ailang-worldd/session_budget_test.go:9` `TestSessionRootBudgets` | `inv-cli`. Existing mint/revoke root names remain pinned. |
| `host/verifygate/subprocess_sink_gate_test.go:223` `TestSubprocessSinkGate` | `inv-verifygate`. The worker uses private unlinked `*os.File` streams. Tests use `syncBuffer` for subprocess-directed logs. |
| `host/verifygate/forklocked_write_test.go:520,717` `TestVerifygateTestWritesAreForkLocked`, `TestNoParallelWriteForkPackagesOutsideVerifygate` | `inv-verifygate` and whole verifygate. The new CLI fixtures write and fork serially. |
| `host/verifygate/forklocked_write_test.go:309` `TestKernelRefusesExecOfWriterOpenFile`; `module_manifest_gate_test.go:495` `TestExecGateArmsAreLoadBearing` | Whole verifygate and `inv-verifygate`. The first is a kernel exec control and the second checks verifier arms. The source sweep found no separate list that enumerates every production CLI exec site. The known positive production sites are `episode_launch.go:556` (headless supervisor) and `host/agentrun/agentrun.go:46` (child creation). |

| Leg | Pristine archive | M1 | M2 |
|---|---|---|---|
| `context-roots` | PASS in inv-store | OK 10s | OK 7s |
| `inv-store` | OK 22s | OK 3s | OK 3s |
| `inv-boundary` | OK 15s | OK 10s | OK 11s |
| `inv-cli` | OK 12s | OK 27s | OK 28s |
| `inv-daemon` | OK 8s | OK 5s | OK 14s |
| `inv-verbs` | OK 2s | OK 2s | OK 1s |
| `inv-verifygate` | OK 81s | PASS on bounded retry; Linux kernel control skipped on macOS | PASS on bounded retry; Linux kernel control skipped on macOS |

### 1.3 Budgets and iteration-442 notes

| Note | Decision and pin |
|---|---|
| J21, tick passes `time.Second` | `TestOperatorProductionDeadlines` runs the empty production loop and observes the requested tick budget and its effective deadline. It insists on a 22 s preparation request and a 1 s heartbeat request. |
| J22, zero override returns `budgetContext(budget+time.Hour)` | Same test bounds the observed absolute deadline between its before/after clock samples. No waiting for a production deadline. Also observes the nonzero override. `TestOperatorAttendedStartDisposition` observes both requested 5 s reconciliation contexts in the real attended fixture path. |
| Git timeout | `WithCommandTimeout` clamps the override to 30 s. `TestWithCommandTimeoutClampsGitCap` uses a cancelled context and verifies both large and small values reaching the Git timeout error immediately. |
| Shared override | The 5 s launch admission and writeback contexts use the same test-only `operatorBudget`. The signed child wall cap is independent. `TestOperatorStartOnce/launch` observes three 5 s launch requests and the configured test deadline, checks the signed one-hour cap passed to the launcher, and checks inherited locks. `TestOperatorLaunchExit/capped_and_scrubbed` cancels an injected child-cap context after a real byte handshake. |
| Stuck owner | Pin the cursor's strict advancement with `TestNextOwnedOperatorStart` and the runner behavior in `TestOperatorStartOnce/launch`. In the nested two-tick replay of a reserved row, exactly one tick emits the recovery remedy, and the following tick reaches the cursor reset. Both preserve the reservation. |

### 1.4 Runtime anchor measurements

| Command or positive source control | Observation at b2b0ebab |
|---|---|
| `grep -nF OperatorStartPrepared host/store/operator_start_recovery.go` | The eighth value is line 20. Classifier sets it after matching observed worktree, mint and files, at line 135. |
| `grep -nF 'spawn crossed the S6 boundary' cmd/ailang-worldd/operator_activation_test.go` | Lines 325–326 forbid the fake `ran` marker after a real preparation tick and readback. S6 still forbids any child or launch step in that path. Exactly one marker is allowed only on the subsequent prepared-row launch tick, after the durable reservation and the worker's single dispatch grant. |
| `grep -n -e 'func TestOperatorFirstManagedEpisodeProcess' -e 'func TestOperatorCurrentOwnerProcessRecovery' -e 'func TestOperatorStartOnce' cmd/ailang-worldd/operator_*test.go` | `operator_activation_test.go:98,369`, `operator_worktree_test.go:215`. The last is an existing test for worktree retry, so its new launch half is a subtest. |
| `grep -nF 'func TestOperatorNeverDispatch' host/verbs/operator_policy_test.go` | Line 69. The service family is already refused from the app. |
| `grep -nF 'runOperatorPreparationLoop' cmd/ailang-worldd/operator_runner.go` | `runOperatorRunner` at line 65 calls the installed loop at line 92, which calls `consumeOperatorPreparation` at line 114. That consumer is at `operator_preparation.go:21`. |
| `grep -nE 'launchEpisode\(|agentrun.Run\(|keepLaunchStream\(' cmd/ailang-worldd/{episode_launch,episode_preset,run_cmd}.go` | V3 `launchEpisode` declared at `episode_launch.go:264`, CLI passes through `episode_preset.go:379`. V3 runs the headless child at line 556 and scans both streams at lines 594–595. Reuse its mode, distinct lock rule, argv, task/prompt, context, stream scrub and supervisor. The managed caller loads the exact approved executable instead of V3 PATH discovery. |
| `grep -nF launch_binary_sha256 host/authority/operator_start_digest.go cmd/ailang-worldd/{operator_start_input,operator_preparation}.go` | Digest fact constructed at `operator_start_input.go:75`, parsed at `host/authority/operator_start_digest.go:10–20`, rechecked at `operator_preparation.go:75`. Launch checks digest at reservation, worker preflight and actual dispatch, and compares it against receipt readback. |
| `grep -nF 'func WithCommandTimeout' host/worktree/worktree.go` | Line 50. Line 51 stored the raw override, and line 151 reads it on the one Git spawn path. |
| `grep -nF 'VERIFY_GO_RACE_PKG_TIMEOUT_S=900' scripts/verify_go.sh` | Line 161. The #691 CI readings are 422.9 s CLI, 452.3 s daemon, 394.0 s verifygate. These are controller CI provenance, not the local timing samples. |

Run `/bin/bash "$S/anchors.sh" "$S/base-tree"` for every reused declaration. The final diff provides the old and new hunk anchors on every touched path. The production caller census is in §6.

## 2. Milestones

Export the requested base to scratch for inspection. Execute the patches against the executor worktree, in the order below. The controller records the planner files before running the snapshot commit script, which refuses untracked design documents. The controller applies the snapshots to the commit worktree.

| Milestone | Added | Deleted | Cumulative | Product without filter |
|---|---:|---:|---:|---:|
| M1 | 489 | 23 | 512 | 512 |
| M2 | 1604 | 46 | 1650 | 1650 |

Per-path measurements are in `$S/numstat-M1` and `$S/numstat-M2`.

1. **M1 — durable launch law and bounded deadline pins.** Reserve a single immutable run ID per start, commit the prepared snapshot CAS, bind the private launch handoff and prove the current session and enrolled approval again in the dispatch transaction. The old generic effect writer refuses launch-kind bypass. Distinct receipt versions preserve process identity, lockout through readback/GC/pruning, and atomic completion with episode launch history. Install the observer for the carried deadline pins and clamp Git's cap. This stage has no production launch caller.
2. **M2 — installed launch caller and verification.** The real preparation loop resumes the prepared action under the episode lock. It rechecks the pair, health, signed configuration, real files and current credential, and reserves while holding lifecycle exclusion. Start the private World worker with the runner, episode and lifecycle descriptions and a one-use handoff over stdin. Its only effect is one supervised headless agent in the V3 isolated launch directory. Before `cmd.Start`, it takes the durable grant and holds the credentials writer across start and initial receipt. The child and worker keep the runner and episode locks through wait. Store the redacted episode history with the final receipt. Expose terminal status, add usage, and pin the non-browser restart path.

## 3. Durable windows and operator output

| Window | Restart/readback and effect law |
|---|---|
| Prepared, no launch step | A new runner may take the episode lock and revalidate the action, approval and exact live session. A stale key, session, pair or executable refuses admission. The token, task and existing mint stay intact. |
| Reserved v1, dispatch grant absent | Only the live handoff has the secret committed by reservation. A fresh runner's read is `uncertain`, preserves the reservation and never reconstructs the handoff. |
| Reserved v2, dispatch grant committed | Possible child effect even with no observed PID. Immutable lockout, read-only `uncertain`, no second launch. Positive synchronous exec refusal can advance to known no-spawn. |
| Observed positive start v3 | Receipt has action/run ID, signed executable digest, PID and kernel birth identity. Action succeeds at observation. Episode owner, GC, claim and pruning still honor the unsettled launch. A live matching birth is `launched`. PID reuse, a missing process or unavailable process identity reads `uncertain`; preserve the evidence and effect exclusion. |
| Known no-spawn v2 or v3 | Record `spawn-refused` and a failed action with no PID, birth, positive process receipt or appended launch history. |
| Completed v4 | Atomically record `stopped`, `failed` or `capped` with the immutable run identity and episode history. Read returns the same outcome after daemon restart or session expiry. |

Existing service stop and paired upgrade share lifecycle and runner kernel exclusion. A worker that outlives the polling runner retains the process fence and signed cap. The supervisor waits boundedly and sends stop only to its own worker. A failed cleanup preserves the durable unresolved window. The attended remedy is `service runner launch-status ACTION-ID`, runner log and positive process evidence before any operator disposition. Browser actions and new effects remain blocked on uncertain ownership.

`ailang-worldd service runner launch-status ACTION-ID` prints `operator launch <id>: <status>`. The same prefix logs `launched`, `stopped`, `failed`, `capped`, `spawn-refused`, `expired` and `blocked or uncertain; inspect service runner launch-status`. Recovery of an owned launch emits the existing `operator start blocked; attended recovery:` prefix with the status command. The launch ledger holds the run ID. The episode report keeps its existing transcript and exit fields. `docs/QUICKSTART.md` names the command and drain behavior. The private process receipt is local SQLite data. No new HTTP route is needed.

## 4. Acceptance, timing and fail sets

| Branch from design §5.3 | Fixture and gate |
|---|---|
| `TestOperatorStartOnce/launch` | Real Git in temp repositories, real signed start, fake supervisor, real fake script child. First tick prepares. Next crosses the reservation before supervisor admission, proves repeated readback, compares V3 launch argv and scrubbed child env, and checks exact one fake execution and run digest. |
| `TestOperatorLaunchSession` | Real migration and test signer in `t.TempDir()`. Store reservation and dispatch compare session, competing credentials, approver and exact-expiry clock. These refusals produce no process effect or process receipt. |
| `TestOperatorSpawnRefused` | Fake service supervisor and signed fake child with a missing interpreter. Actual OS spawn refuses. The run has no process identity or episode launch history. The case compares worktree and private file evidence before and after. |
| `TestOperatorLaunchCrashWindow`, `TestOperatorStartOnce/launch`, `TestOperatorLaunchReceiptLaw` | Pure birth/readback fixture, restart from a real durable reservation and read-only CAS checks. Includes fake supervisor and fake child on the full runner path. A reused PID reads uncertain. The original birth reads launched. |
| `TestOperatorLaunchExit` | Bounded fake shell children under `t.TempDir()` through the real supervisor. Nonzero is failed. Closed output descriptors produce output-failed evidence. The pipe handshake occurs after both child streams output the canary. The injected cap stops the child and the production history function scrubs stdout, stderr and argv evidence. |
| `TestOperatorStartBrowser` daemon-restart half | Fake transport and fake supervisor under `t.TempDir()`, real held fake child. Read a live process via the terminal status command as the fake daemon becomes unavailable. Advance the injected clock beyond expiry, release the child, compare the credential and verify one durable launch history entry. Browser gesture/disconnect is S7 work. |
| Service lifecycle | Existing activation/drain suite plus actual inherited lock FDs from the fake child in `TestOperatorStartOnce/launch`. No live supervisor control in the gates. |

`accept.sh` observes the real-path metadata and allows `absent`. Every change-capable fixture runs with a synthetic home, fake supervisor, a test signer, and temp stores. The only named paths written by the new process tests are in their `t.TempDir()` fixtures. The explicit 1777 run puts private test roots under an ancestor with Linux-style permissions and compares physical aliases with `filepath.EvalSymlinks`.

### Race timings

| New CLI body or subtest | `-race` body seconds |
|---|---:|
| `TestOperatorLaunchCrashWindow` | 0.00 |
| `TestOperatorLaunchExit` | 0.37 |
| `TestOperatorProductionDeadlines` | 0.36 |
| `TestOperatorSpawnRefused` | 2.48 |
| `TestOperatorStartBrowser` | 1.89 |
| `TestOperatorStartOnce/launch` | 1.90 |
Selected raw receipt `$S/timing-repeat-2.raw`. Other focused samples: timing-repeat-1 19.96 s, timing-final-quiet 47.71 s, timing-release 19.35 s, timing-after-controller 12.72 s, timing-sealed 23.35 s. The larger pre-reduction fixture draft measured 29.322 s for its complete selection.

**CLI added total 7.00 s. Daemon added total 0.00 s.** The latter package has zero production or test edits. Parent totals include their branch tests once. Measured with `$S/measure-race.sh`.

| Whole package race leg | Base → M2, delta seconds |
|---|---|
| `host/store` | 119.751 → 93.467 (-26.284) |
| `host/agentrun` | 17.713 → 17.823 (+0.110) |
| `host/worktree` | 4.526 → 60.558 (+56.032) |
| `host/verbs` | 1.582 → 1.473 (-0.109) |
| `cmd/ailang-worldd` | UNINFORMATIVE UNDER SANDBOX |
| `host/daemon` | UNINFORMATIVE UNDER SANDBOX |

### 4.1 Controller outside measurements

Use the pristine controller receipts in `$STATE/base/` only after the exact completion marker. The gate records every leg before comparison. Whole CLI and any daemon runs that attempt a sandboxed bind are **UNINFORMATIVE UNDER SANDBOX**. Do not attribute their partial fail sets. Use complete text logs, record package elapsed time and compare only added failure identities.

```bash
S=$HOME/.ailang/state/world-iter444/plan-scratch
STATE=$HOME/.ailang/state/world-iter444
/bin/bash "$S/outside.sh" "$S/outside-final"
/bin/bash "$S/whole-mutants.sh" "$S/whole-mutant-ids.txt" "$S/ledger-whole-outside.jsonl"
GOMAXPROCS=4 /bin/bash "$S/measure-race.sh" "$S/timing-outside"
```

The first script records M1/M2 CLI plain/race and M1/M2 daemon race, then compares against the controller's `cmd-plain.fails` and `daemon-race.fails`. It records a matched pristine CLI race fail set because the controller's list has no CLI race leg. The second consumes the 25 prioritized mutant IDs, using the controller base package fail sets and separately recording a matching pristine leg where needed. Check `panic: test timed out` in every raw file. Nonzero without a parsed named failure is an instrument failure. Use fresh log names on retries.

Controller baseline seal and fail-set counts:

```text
LEG store-plain rc=0 parser=0 timeouts=0 fails=0 secs=30
LEG store-race rc=0 parser=0 timeouts=0 fails=0 secs=129
LEG authority-plain rc=0 parser=0 timeouts=0 fails=0 secs=7
LEG authority-race rc=0 parser=0 timeouts=0 fails=0 secs=42
LEG worktree-plain rc=0 parser=0 timeouts=0 fails=0 secs=3
LEG verifygate-plain rc=1 parser=0 timeouts=0 fails=1 secs=425
LEG cmd-plain rc=1 parser=0 timeouts=0 fails=12 secs=1090
LEG daemon-race rc=1 parser=0 timeouts=0 fails=9 secs=1624
ALL_DONE for b2b0ebab7bffbcf5f3105aa17e3ced452e480e35
```

The outside baseline `cmd-plain.fails` additionally records `TestOperatorWorktreeDispositionCAS`, beyond the 11 CLI fixture reds. The whole verifygate base records `TestMissionConfigDefaultFlowIntegration`.

## 5. Existing tests and helpers modified

| Existing test/helper | Before | After | Production mutation still killed or positive control |
|---|---|---|---|
| `TestOperatorFirstManagedEpisodeProcess` | Installed preparation loop and readback, expected no fake marker at the old boundary | Same prep entry and no marker, plus no launch reservation/receipt. New `TestOperatorStartOnce/launch` alone allows one child after write-ahead. | P-S5-PREP-CALLER kills the retained preparation entry. D-NORMAL-CALLER and D-RESERVE-ORDER pin the subsequent launch. |
| `TestOperatorAttendedStartDisposition` | `service runner reconcile-start` for held, rejected, drifted and successfully stopped work | Same lifecycle and change checks, plus two observed 5 s contexts | J22-reconcile |
| `operatorActivationFixture` | `newWorktreeStartRigWithContext` prepares the managed fake | Same default and new opt-in fake launch script under the temp runner home | P-S5-PREP-CALLER via the old activation test, D-RESERVE-ORDER via the new launch fixture |
| `newWorktreeStartRigWithContext` | Temp Git and signed worktree fixture | Same default via `newWorktreeStartRigWithAgent`, opt-in fake agrees with paired installed build | P-DERIVED-BRANCH-REPLAY in the old worktree body, P-WORKTREE-REPLAY in the existing registration collision arm |
| `TestOperatorStartOnce` | Worktree effect run and replay refusal via `r.start` | Original worktree replay body in /worktree plus installed-loop /launch | P-DERIVED-BRANCH-REPLAY, D-NORMAL-CALLER, D-RESERVE-ORDER, D-ARGV-CONFINEMENT |
| `TestCLICredentialWriterSites` | AST census of credentials writers | Adds one bounded `executeOperatorLaunch` writer site | I-EXTRA-CLI-WRITER and the existing positive read-only count. Static controls. |
| `TestOperatorStepRestart` | Generic step reserve, observe and replay through a fixture launch-kind row | Same generic behavior with `checks`. The launch-specific reservation and observation move under typed law and are pinned by `TestOperatorLaunchReceiptLaw`. | P-RESTART-GENERIC and D-BYPASS-GENERIC-RESERVE |
| `TestCredentialsMethodsRequireADeadline` | Existing credentials methods | Same plus five typed launch CRUD calls | P-LAUNCH-DEADLINE-REFUSAL and bounded positive calls in TestOperatorLaunchReceiptLaw. |
| `newOperatorPinFixture` | Real migration, start signer and default facts | Default wrapper plus new fixture prefix for the real signed binary digest | D-RESERVE-CAS, D-RECHECK-EXPIRY |
| `TestOperatorNeverDispatch` | `AdmitOperator` service and unknown-family refusals | Also covers launch-status and the private worker | P-NEVER-DISPATCH |

## 6. Gates and mutant receipts

Set `S=$HOME/.ailang/state/world-iter444/plan-scratch` and `STATE=$HOME/.ailang/state/world-iter444`. Exact scripts are below, with helper files already in scratch. Start with the pristine archive in `$S/base-tree`. Execute the named inventories in the same stages as the patches. No PATH instrumentation or JSON-mode verifier. `LEG ... SKIPPED` is an explicit skip, and `ALL_DONE` is written after all recordings and comparisons.

**base.** `$S/gate-base/legs.txt`

```text
LEG gofmt OK 0s
LEG build OK 25s
LEG vet OK 23s
LEG linux-vet OK 77s
LEG inv-store OK 22s
LEG inv-boundary OK 15s
LEG inv-cli OK 12s
LEG inv-verbs OK 2s
LEG inv-daemon OK 8s
LEG inv-verifygate OK 81s
LEG accept SKIPPED (S6 tests absent)
LEG t1777 SKIPPED (S6 tests absent)
LEG verify_ail OK 287s
LEG callers OK 1s
LEG store-plain recorded rc=0 26s
LEG store-race recorded rc=1 117s
LEG authority-plain recorded rc=0 47s
LEG authority-race recorded rc=0 37s
LEG worktree-plain recorded rc=0 8s
LEG worktree-race recorded rc=0 6s
LEG agentrun-plain recorded rc=0 19s
LEG agentrun-race recorded rc=0 22s
LEG verifygate-plain recorded rc=1 414s
LEG cmd-plain UNINFORMATIVE UNDER SANDBOX rc=1
LEG cmd-race UNINFORMATIVE UNDER SANDBOX rc=1
LEG daemon-race UNINFORMATIVE UNDER SANDBOX rc=1
LEG compare-store-plain SKIPPED (pristine record)
LEG compare-store-race SKIPPED (pristine record)
LEG compare-authority-plain SKIPPED (pristine record)
LEG compare-authority-race SKIPPED (pristine record)
LEG compare-worktree-plain SKIPPED (pristine record)
LEG compare-worktree-race SKIPPED (pristine record)
LEG compare-agentrun-plain SKIPPED (pristine record)
LEG compare-agentrun-race SKIPPED (pristine record)
LEG compare-verifygate-plain SKIPPED (pristine record)
LEG compare-cmd-plain SKIPPED (pristine record)
LEG compare-cmd-race SKIPPED (pristine record)
LEG compare-daemon-race SKIPPED (pristine record)
GATES_DONE issue=0 uninformative=1
```

**M1-seal.** `$S/gate-M1-seal/legs.txt`

```text
LEG gofmt OK 1s
LEG build OK 49s
LEG vet OK 5s
LEG linux-vet OK 3s
LEG context-roots OK 10s
LEG inv-store OK 3s
LEG inv-boundary OK 10s
LEG inv-cli OK 27s
LEG inv-verbs OK 2s
LEG inv-daemon OK 5s
LEG inv-verifygate INSTRUMENT FAILURE (truncated) rc=1
LEG accept OK 58s
LEG t1777 OK 24s
LEG verify_ail OK 183s
LEG callers OK 1s
LEG store-plain recorded rc=0 67s
LEG store-race recorded rc=0 160s
LEG authority-plain recorded rc=0 61s
LEG authority-race recorded rc=0 113s
LEG worktree-plain recorded rc=0 13s
LEG worktree-race recorded rc=0 36s
LEG agentrun-plain recorded rc=0 28s
LEG agentrun-race recorded rc=0 29s
LEG verifygate-plain recorded rc=1 661s
LEG cmd-plain UNINFORMATIVE UNDER SANDBOX rc=1
LEG cmd-race UNINFORMATIVE UNDER SANDBOX rc=1
LEG daemon-race UNINFORMATIVE UNDER SANDBOX rc=1
LEG verbs-plain recorded rc=0 1s
LEG verbs-race recorded rc=0 5s
LEG compare-store-plain OK
LEG compare-store-race OK
LEG compare-authority-plain OK
LEG compare-authority-race OK
LEG compare-worktree-plain OK
LEG compare-worktree-race OK
LEG compare-agentrun-plain OK
LEG compare-agentrun-race OK
LEG compare-verifygate-plain OK
LEG compare-cmd-plain SKIPPED (head record unavailable)
LEG compare-cmd-race SKIPPED (head record unavailable)
LEG compare-daemon-race SKIPPED (head record unavailable)
LEG compare-verbs-plain OK
LEG compare-verbs-race OK
GATES_DONE issue=1 uninformative=1
```

**M2-seal.** `$S/gate-M2-seal/legs.txt`

```text
LEG gofmt OK 1s
LEG build OK 43s
LEG vet OK 11s
LEG linux-vet OK 2s
LEG context-roots OK 7s
LEG inv-store OK 3s
LEG inv-boundary OK 11s
LEG inv-cli OK 28s
LEG inv-verbs OK 1s
LEG inv-daemon OK 14s
LEG inv-verifygate INSTRUMENT FAILURE (truncated) rc=1
LEG accept OK 67s
LEG t1777 OK 39s
LEG verify_ail OK 151s
LEG callers OK 3s
LEG store-plain recorded rc=0 68s
LEG store-race recorded rc=0 152s
LEG authority-plain recorded rc=0 70s
LEG authority-race recorded rc=0 122s
LEG worktree-plain recorded rc=0 56s
LEG worktree-race recorded rc=0 69s
LEG agentrun-plain recorded rc=0 16s
LEG agentrun-race recorded rc=0 19s
LEG verifygate-plain recorded rc=1 569s
LEG cmd-plain UNINFORMATIVE UNDER SANDBOX rc=1
LEG cmd-race UNINFORMATIVE UNDER SANDBOX rc=1
LEG daemon-race UNINFORMATIVE UNDER SANDBOX rc=1
LEG verbs-plain recorded rc=0 0s
LEG verbs-race recorded rc=0 4s
LEG compare-store-plain OK
LEG compare-store-race OK
LEG compare-authority-plain OK
LEG compare-authority-race OK
LEG compare-worktree-plain OK
LEG compare-worktree-race OK
LEG compare-agentrun-plain OK
LEG compare-agentrun-race OK
LEG compare-verifygate-plain OK
LEG compare-cmd-plain SKIPPED (head record unavailable)
LEG compare-cmd-race SKIPPED (head record unavailable)
LEG compare-daemon-race SKIPPED (head record unavailable)
LEG compare-verbs-plain OK
LEG compare-verbs-race OK
GATES_DONE issue=1 uninformative=1
```

Release acceptance markers: ACCEPT_OK M1 fast, ACCEPT_OK M2 fast, ACCEPT_OK M2 fast. The stage release build/native vet/Linux vet receipts are `$S/M1-sealed-build.raw.rc` and `$S/M2-sealed-build.raw.rc` (0).

The isolated final commit verifygate rehearsal returned 0. Raw receipt `$S/commit-M2-verifygate-rehearsal.raw` (411.832 s, empty fail set).

Initial stage inventory legs hit the short probe timeout around the verifier's child checks. The bounded retries in `$S/inv-verifygate-M1-final.raw` and `$S/inv-verifygate-M2-final.raw` passed. Both recorded whole verifygate fail sets matched the controller's one base failure. The final gate script uses the expanded probe budget and a platform skip for the Linux-only exec kernel control.

````bash gates.sh
#!/bin/bash
# gates.sh TREE LOGDIR base|M1|M2 [PRISTINE-LOGDIR]
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P) || exit 2
CONTROL=$S/../base
T=$(cd "$1" && pwd -P) || exit 2
mkdir -p "$2" || exit 2
L=$(cd "$2" && pwd -P) || exit 2
M=$3
PRISTINE=${4:-}
case "$M" in base|M1|M2) ;; *) exit 2;; esac
test ! -e "$L/legs.txt" || exit 2
B=b2b0ebab7bffbcf5f3105aa17e3ced452e480e35
PIN=$HOME
export WORLD_PLAN_PIN_HOME="$PIN"
export AILANG_BIN="$PIN/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$PIN/.pinned-ailang-tools/v0.54.15/ailang"
export GOCACHE="$(go env GOCACHE)" GOMODCACHE="$(go env GOMODCACHE)"
export HOME="$L/home" AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
export XDG_CONFIG_HOME="$HOME/.config" XDG_CACHE_HOME="$HOME/.cache" XDG_DATA_HOME="$HOME/.local/share"
mkdir -p "$HOME" "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_TEST_OPERATOR_WITNESS WORLD_TEST_OPERATOR_WITNESS_PATH
cd "$T" || exit 2
issue=0 unknown=0
note() { printf '%s\n' "$*" | tee -a "$L/legs.txt"; }
skip() { note "LEG $1 SKIPPED ($2)"; }
bound() { python3 "$S/bounded.py" "$@"; }
denied() { grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$1"; }
truncated() { grep -Eq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$1"; }
leg() {
  local tag=$1 start rc=0; shift
  start=$(date +%s)
  "$@" > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  if truncated "$L/$tag.raw"; then issue=1; note "LEG $tag INSTRUMENT FAILURE (truncated) rc=$rc"
  elif denied "$L/$tag.raw" || [ "$rc" -eq 77 ]; then unknown=1; note "LEG $tag UNINFORMATIVE UNDER SANDBOX rc=$rc"
  elif [ "$rc" -ne 0 ]; then issue=1; note "LEG $tag rc=$rc $(( $(date +%s)-start ))s"
  else note "LEG $tag OK $(( $(date +%s)-start ))s"; fi
}
record() {
  local tag=$1 pkg=$2 start rc=0 parser=0; shift 2
  start=$(date +%s)
  bound 930 go test -count=1 -timeout 900s "$@" "./$pkg/" > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  python3 "$S/failset.py" "$L/$tag.raw" "$rc" > "$L/$tag.fails" 2> "$L/$tag.parser" || parser=$?
  cat "$L/$tag.parser"
  if denied "$L/$tag.raw"; then unknown=1; note "LEG $tag UNINFORMATIVE UNDER SANDBOX rc=$rc"
  elif truncated "$L/$tag.raw" || [ "$parser" -ne 0 ]; then issue=1; note "LEG $tag INSTRUMENT FAILURE rc=$rc parser=$parser"
  else cp "$L/$tag.fails" "$L/$tag.attributable"; note "LEG $tag recorded rc=$rc $(( $(date +%s)-start ))s"; fi
}
fmt_touched() {
  local p bad
  while IFS= read -r p; do
    case "$p" in *.go) if [ -f "$p" ]; then bad=$(gofmt -l "$p") || return 2; test -z "$bad" || { printf '%s\n' "$bad"; return 1; }; fi;; esac
  done < "$S/changed-${M/base/M1}"
}
inventory() {
  local pkg=$1 names=$2 name output rc=0
  output=$(mktemp "$L/inventory.XXXXXX") || return 2
  bound 380 go test -v -count=1 -timeout 350s "./$pkg/" -run "^($names)\$" > "$output" 2>&1 || rc=$?
  cat "$output"
  [ "$rc" -eq 0 ] || return "$rc"
  while IFS= read -r name; do grep -Fq -- "--- PASS: $name (" "$output" || return 2; done < <(printf '%s\n' "$names" | tr '|' '\n')
}
leg gofmt fmt_touched
leg build bound 300 go build ./...
leg vet bound 300 go vet ./cmd/ailang-worldd/ ./host/store/ ./host/agentrun/ ./host/worktree/ ./host/verbs/
leg linux-vet env GOOS=linux python3 "$S/bounded.py" 300 go vet ./cmd/ailang-worldd/ ./host/store/ ./host/agentrun/ ./host/worktree/ ./host/verbs/
leg context-roots bound 180 go test ./host/store/ -run '^TestProductionContextRoots$'
leg inv-store inventory host/store 'TestProductionContextRoots|TestStoreSurfaceTakesContext|TestCredentialsMethodsRequireADeadline'
leg inv-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
leg inv-cli inventory cmd/ailang-worldd 'TestCLICredentialWriterSites|TestSessionRootBudgets'
leg inv-verbs inventory host/verbs TestOperatorNeverDispatch
leg inv-daemon inventory host/daemon TestDaemonDependencyAllowlist
leg inv-verifygate inventory host/verifygate 'TestSubprocessSinkGate|TestExecGateArmsAreLoadBearing|TestNoParallelWriteForkPackagesOutsideVerifygate|TestVerifygateTestWritesAreForkLocked'
if [ "$(uname -s)" = Linux ]; then
  leg inv-exec-kernel inventory host/verifygate TestKernelRefusesExecOfWriterOpenFile
else
  skip inv-exec-kernel 'Linux ETXTBSY control'
fi
if [ "$M" = base ]; then
  skip accept 'S6 tests absent'
  skip t1777 'S6 tests absent'
else
  leg accept /bin/bash "$S/accept.sh" "$T" "$M"
  mkdir -m 1777 "$L/t1777"
  mkdir -m 0700 "$L/t1777/x"
  leg t1777 env TMPDIR="$L/t1777/x" /bin/bash "$S/accept.sh" "$T" "$M"
fi
leg verify_ail bound 510 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
leg callers /bin/bash "$S/callers.sh" "$T" "$M"
tags=(store-plain store-race authority-plain authority-race worktree-plain worktree-race agentrun-plain agentrun-race verifygate-plain cmd-plain cmd-race daemon-race verbs-plain verbs-race)
pkgs=(host/store host/store host/authority host/authority host/worktree host/worktree host/agentrun host/agentrun host/verifygate cmd/ailang-worldd cmd/ailang-worldd host/daemon host/verbs host/verbs)
for i in "${!tags[@]}"; do
  tag=${tags[$i]}
  case "$tag" in *-race) record "$tag" "${pkgs[$i]}" -race;; *) record "$tag" "${pkgs[$i]}";; esac
done
# Recording finishes for every leg before comparison.
for tag in "${tags[@]}"; do
  if [ "$M" = base ]; then skip "compare-$tag" 'pristine record'; continue; fi
  if [ ! -f "$L/$tag.attributable" ]; then skip "compare-$tag" 'head record unavailable'; continue; fi
  source=${PRISTINE:+$PRISTINE/$tag.attributable}
  if [ -f "$CONTROL/../baseline.done" ] &&
      grep -Exq "ALL_DONE( for $B)?" <(tail -n 1 "$CONTROL/../baseline.done") &&
      [ -s "$CONTROL/$tag.raw" ] && [ -f "$CONTROL/$tag.rc" ] &&
      ! truncated "$CONTROL/$tag.raw" && ! denied "$CONTROL/$tag.raw" &&
      python3 "$S/failset.py" "$CONTROL/$tag.raw" "$(cat "$CONTROL/$tag.rc")" > "$L/$tag.controller"; then
    source="$L/$tag.controller"
  fi
  if [ ! -f "$source" ]; then skip "compare-$tag" 'attributable base unavailable'; unknown=1; continue; fi
  LC_ALL=C sort -u "$source" > "$L/$tag.base"
  LC_ALL=C sort -u "$L/$tag.attributable" > "$L/$tag.head"
  LC_ALL=C comm -13 "$L/$tag.base" "$L/$tag.head" > "$L/$tag.new"
  if [ -s "$L/$tag.new" ]; then issue=1; note "LEG compare-$tag NEW FAIL SET"; else note "LEG compare-$tag OK"; fi
done
note "GATES_DONE issue=$issue uninformative=$unknown"
printf 'ALL_DONE\n' > "$L/ALL_DONE"
[ "$issue" -eq 0 ] || exit 1
[ "$unknown" -eq 0 ] || exit 77
````

````bash accept.sh
#!/bin/bash
# accept.sh TREE M1|M2 [fast|race]
set -euo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
M=$2
mode=${3:-fast}
case "$M" in M1|M2) ;; *) exit 2;; esac
case "$mode" in fast) race=;; race) race=-race;; *) exit 2;; esac
PIN=${WORLD_PLAN_PIN_HOME:-$HOME}
D=$(mktemp -d "$S/accept.XXXXXX")
CLI='TestOperatorProductionDeadlines|TestOperatorAttendedStartDisposition|TestOperatorFirstManagedEpisodeProcess|TestOperatorCurrentOwnerProcessRecovery|TestOperatorActivationDrainRefusals'
if [ "$M" = M2 ]; then
  CLI="$CLI|TestOperatorSpawnRefused|TestOperatorLaunchCrashWindow|TestOperatorLaunchExit|TestOperatorStartBrowser|TestOperatorStartOnce"
fi
observe() {
  for p in "$PIN/Library/LaunchAgents/com.sunholo.ailang-world.worldd.plist" \
    "$PIN/Library/LaunchAgents/com.sunholo.ailang-world.operator-runner.plist" \
    "$PIN/.ailang-world" "$PIN/.ailang/world-ws"; do
    if [ ! -e "$p" ] && [ ! -L "$p" ]; then printf 'absent %s\n' "$p"
    elif [ "$(uname -s)" = Darwin ]; then stat -f '%N %p %i %z %m %c' "$p"
    else stat -c '%n %a %i %s %Y %Z' "$p"; fi
  done
}
observe > "$D/before"
run() {
  local pkg=$1 names=$2 log="$D/${1##*/}.raw" rc=0
  /bin/bash "$S/go-leg.sh" "$T" "$log" 380 go test ${race:+-race} -count=1 -v -timeout 340s "./$pkg/" -run "^($names)\$" || rc=$?
  cat "$log"
  if grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$log"; then
    echo 'UNINFORMATIVE UNDER SANDBOX'; return 77
  fi
  if grep -Eq 'panic: test timed out|BOUND EXCEEDED' "$log"; then return 2; fi
  if [ "$rc" -ne 0 ]; then return "$rc"; fi
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log"
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
STORE='TestOperatorLaunchReceiptLaw|TestOperatorStartRecoveryWindows|TestProductionContextRoots|TestCredentialsMethodsRequireADeadline|TestOperatorStepRestart'
if [ "$M" = M2 ]; then STORE="$STORE|TestOperatorLaunchSession"; fi
run host/store "$STORE"
run host/worktree TestWithCommandTimeoutClampsGitCap
run cmd/ailang-worldd "$CLI"
if [ "$M" = M2 ]; then
  run host/verbs TestOperatorNeverDispatch
fi
observe > "$D/after"
diff -u "$D/before" "$D/after"
printf 'ACCEPT_OK %s %s\n' "$M" "$mode"
````

### Production caller census

```text
92:	return e.runOperatorPreparationLoop(errOut, c, h, source, stops, ticker.C, lock)
114:		err = e.consumeOperatorPreparation(ctx, c, h, src, stops, lock, &recoveryCursor, store.OpenOperatorEpisodeLock)
556:		result, runErr = agentrun.Run(runCtx, argv, dir, childEnv, outStage, errStage, wall)

LAUNCH_CALLERS base 0
```

```text
95:	return e.runOperatorPreparationLoop(errOut, c, h, source, stops, ticker.C, lock)
117:		err = e.consumeOperatorPreparation(ctx, c, h, src, stops, lock, &recoveryCursor, store.OpenOperatorEpisodeLock)
556:		result, runErr = agentrun.Run(runCtx, argv, dir, childEnv, outStage, errStage, wall)

LAUNCH_CALLERS M1 0
```

```text
95:	return e.runOperatorPreparationLoop(errOut, c, h, source, stops, ticker.C, lock)
127:		err = e.consumeOperatorPreparation(ctx, c, h, src, stops, lock, &recoveryCursor, store.OpenOperatorEpisodeLock, launch)
556:		result, runErr = agentrun.Run(runCtx, argv, dir, childEnv, outStage, errStage, wall)
cmd/ailang-worldd/operator_runner.go:118:			err := e.launchPreparedOperatorStart(ctx, c, h, src, stops, lock, errOut, ro, row, current, episode)
LAUNCH_CALLERS M2 1
212:	return e.superviseOperatorLaunch(job, stops, log, h, in.digest)
217:	cmd := exec.Command(e.paths.Binary, "operator-launch-worker", job.row.ID)
396:		return realServiceEnv().runOperatorLaunchWorker(rest[1:], stderr)
308:	if err := e.executeOperatorLaunch(proc, job, out); err != nil {
479:		run = agentrun.RunObserved
```

The existing V3 `agentrun.Run` and the existing preparation consumer are firing controls. The M2 production chain is `runOperatorRunner → runOperatorPreparationLoop → consumeOperatorPreparation → recoverOperatorStart → recoverOwnedStart → launchPreparedOperatorStart → superviseOperatorLaunch → operator-launch-worker → executeOperatorLaunch → agentrun.RunObserved`.

### Mutations

The prioritized exact edit records are in `$S/mutants.json`. The additional Go conditional probes are recorded in `$S/ledger-generated-store.jsonl`, `$S/ledger-generated-cli.jsonl` and the generated repin ledgers. `$S/branch-audit.tsv` tracks 241 changed production condition/order anchors with 153 anchors requiring a named edit or equivalence review. Those rows are in `$S/remaining-anchors.tsv`. The ledger table counts only the enumerated S6 experiment set. Old source predicates newly reached by the launch law are explicitly labeled by their file and line. The prior mutant ledgers supply no kills. For explicit draft discoveries, the launch-start log mutant survived until the fake restart asserted both lifecycle lines, and the launch-input early-return mutant survived until the test compared returned files and action identity. The competing-session fixture closes the query's early-return probe. The final named status is shown in the table.

Conditional edit candidates are in `$S/condition-proposals.json`, produced with `python3 "$S/propose-conditions.py"`. PENDING — controller to measure the remaining conditions with `python3 "$S/mutate.py" named --spec "$S/condition-proposals.json" ID` after confirming each positive selector. The remaining SQL/ordering entries need literal edits matched to their anchors in `branch-audit.tsv`. Keep the outside whole set to the prioritized 25 IDs. A source-only guard result supplies inventory coverage. A kill requires the named behavioral fail line.

| ID | Source anchor | Named pin | Named | Whole |
|---|---|---|---|---|
| J21 | `cmd/ailang-worldd/operator_runner.go:115` | `TestOperatorProductionDeadlines` | KILLED | PENDING |
| J22 | `cmd/ailang-worldd/operator_runner.go:41` | `TestOperatorProductionDeadlines` | KILLED | PENDING |
| J22-heartbeat | `cmd/ailang-worldd/operator_runner.go:157` | `TestOperatorProductionDeadlines` | KILLED | PENDING |
| J22-reconcile | `cmd/ailang-worldd/operator_start_recovery.go:275` | `TestOperatorAttendedStartDisposition` | KILLED | PENDING |
| J-git-cap | `host/worktree/worktree.go:51` | `TestWithCommandTimeoutClampsGitCap` | KILLED | KILLED |
| J-shared-override | `cmd/ailang-worldd/operator_runner.go:43` | `TestOperatorProductionDeadlines` | KILLED | PENDING |
| J-stuck-cursor | `host/store/operator_start_recovery.go:206` | `TestNextOwnedOperatorStart` | KILLED | PENDING |
| J-launch-budget | `cmd/ailang-worldd/operator_launch.go:29` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-NORMAL-CALLER | `cmd/ailang-worldd/operator_runner.go:118` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-RESERVE-ORDER | `cmd/ailang-worldd/operator_launch.go:202` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-BYPASS-GENERIC-RESERVE | `host/store/operator_steps.go:21` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-BYPASS-GENERIC-RECEIPT | `host/store/operator_steps.go:49` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-READ-PROCESS-DIGEST | `host/store/operator_launch.go:189` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-RESERVE-CAS | `host/store/operator_launch.go:64` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-SPAWN-HANDOFF | `host/store/operator_launch.go:109` | `TestOperatorLaunchReceiptLaw` | KILLED | KILLED |
| D-SPAWN-COMMIT | `host/store/operator_launch.go:118` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-RECHECK-EXPIRY | `host/store/operator_files.go:203` | `TestOperatorLaunchSession` | KILLED | KILLED |
| D-RECHECK-SESSION | `host/store/operator_files.go:222` | `TestOperatorLaunchSession` | KILLED | PENDING |
| D-RECHECK-KEY | `host/store/operator_files.go:219` | `TestOperatorLaunchSession` | KILLED | PENDING |
| D-PUBLISH-START-LOG | `cmd/ailang-worldd/operator_launch.go:469` | `TestOperatorStartBrowser` | KILLED | PENDING |
| D-PUBLISH-START-AT-EXIT | `cmd/ailang-worldd/operator_launch.go:469` | `TestOperatorStartBrowser` | KILLED | PENDING |
| D-READ-BIRTH | `cmd/ailang-worldd/operator_launch.go:330` | `TestOperatorLaunchCrashWindow` | KILLED | PENDING |
| D-PROCESS-REPLAY | `cmd/ailang-worldd/operator_start_recovery.go:49` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-ARGV-CONFINEMENT | `cmd/ailang-worldd/operator_launch.go:139` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-HOLD-EPISODE | `cmd/ailang-worldd/operator_launch.go:437` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-CAP-BINDING | `cmd/ailang-worldd/operator_launch.go:480` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-EMIT-REFUSAL | `cmd/ailang-worldd/operator_launch.go:248` | `TestOperatorSpawnRefused` | KILLED | PENDING |
| D-CAP-STATUS | `cmd/ailang-worldd/operator_launch.go:523` | `TestOperatorLaunchExit/capped_and_scrubbed` | KILLED | PENDING |
| D-NONZERO-STATUS | `cmd/ailang-worldd/operator_launch.go:526` | `TestOperatorLaunchExit/failed` | KILLED | PENDING |
| D-SCRUB-STDOUT | `cmd/ailang-worldd/operator_launch.go:509` | `TestOperatorLaunchExit/capped_and_scrubbed` | KILLED | PENDING |
| D-SCRUB-STDERR | `cmd/ailang-worldd/operator_launch.go:510` | `TestOperatorLaunchExit/capped_and_scrubbed` | KILLED | PENDING |
| D-EXIT-BIRTH-CAS | `host/store/operator_launch.go:242` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-EXIT-HISTORY-ATOMIC | `host/store/operator_launch.go:271` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-UNRESOLVED-OWNER | `host/store/operator_runner.go:105` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-UNRESOLVED-GC | `host/store/operator_episode_gc.go:24` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-UNRESOLVED-PRUNE | `host/store/operator_actions.go:189` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-WRITEBACK-EXPIRY | `host/store/operator_launch.go:225` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| P-RESTART-GENERIC | `host/store/operator_steps.go:36` | `TestOperatorStepRestart` | KILLED | KILLED |
| P-NEVER-DISPATCH | `host/verbs/operator.go:143` | `TestOperatorNeverDispatch` | KILLED | KILLED |
| D-STORE-RESERVE-COMMIT | `host/store/operator_launch.go:70` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-HANDOFF-NONCE-COMMITMENT | `host/store/operator_launch.go:41` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-STATE-SUCCESS-ON-START | `host/store/operator_launch.go:235` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-PROCESS-READ-OWNER | `host/store/operator_launch.go:179` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-RUN-ID-BINDING | `host/store/operator_launch.go:227` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-REMINT-ON-RESUME | `host/store/operator_start_recovery.go:206` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-AFTER-EXEC-NO-PROCESS | `cmd/ailang-worldd/operator_launch.go:484` | `TestOperatorSpawnRefused` | KILLED | PENDING |
| D-PRIVATE-CWD-REFUSAL | `cmd/ailang-worldd/operator_launch.go:420` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-EXEC-DIGEST-COMPARE | `cmd/ailang-worldd/operator_launch.go:447` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-SERVICE-PAIR-COMPARE | `cmd/ailang-worldd/operator_launch.go:65` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-EXACT-DEADLINE-PROOF | `cmd/ailang-worldd/operator_launch.go:62` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-EXEC-ADMISSION-CALL | `host/agentrun/agentrun.go:81` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-EXEC-OBSERVATION-CALL | `host/agentrun/agentrun.go:92` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-CLOSE-PROC-IDENTITY-CAS | `host/store/operator_launch.go:242` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| D-KNOWN-REFUSAL-NO-PROCESS | `host/store/operator_launch.go:200` | `TestOperatorSpawnRefused` | KILLED | PENDING |
| D-LIFECYCLE-LOCK | `cmd/ailang-worldd/operator_launch.go:182` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-EPISODE-LOCK-VERIFIED | `cmd/ailang-worldd/operator_launch.go:178` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| D-RAW-STREAM-PATH | `cmd/ailang-worldd/operator_launch.go:367` | `TestOperatorLaunchExit/capped_and_scrubbed` | KILLED | PENDING |
| D-OUTPUT-FAILURE-STATUS | `cmd/ailang-worldd/operator_launch.go:512` | `TestOperatorLaunchExit/failed` | KILLED | PENDING |
| D-OTHER-CREDENTIAL | `host/store/operator_launch.go:138` | `TestOperatorLaunchSession` | KILLED | PENDING |
| P-S5-PREP-CALLER | `cmd/ailang-worldd/operator_preparation.go:32` | `TestOperatorFirstManagedEpisodeProcess` | KILLED | PENDING |
| P-WORKTREE-REPLAY | `cmd/ailang-worldd/operator_worktree.go:58` | `TestOperatorStartCollisions/registration` | KILLED | PENDING |
| P-LAUNCH-DEADLINE-REFUSAL | `host/store/operator_launch.go:47` | `TestCredentialsMethodsRequireADeadline` | KILLED | PENDING |
| I-EXTRA-CLI-WRITER | `cmd/ailang-worldd/operator_launch.go:337` | `TestCLICredentialWriterSites` | KILLED | PENDING |
| P-DERIVED-BRANCH-REPLAY | `cmd/ailang-worldd/episode_preset.go:183` | `TestOperatorStartOnce/worktree` | KILLED | PENDING |
| Q-host-store-operator_launch-48 | `host/store/operator_launch.go:46` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-51 | `host/store/operator_launch.go:46` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-55 | `host/store/operator_launch.go:50` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-61 | `host/store/operator_launch.go:61` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-67 | `host/store/operator_launch.go:66` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-79 | `host/store/operator_launch.go:77` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-83 | `host/store/operator_launch.go:81` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-88 | `host/store/operator_launch.go:88` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-97 | `host/store/operator_launch.go:95` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-100 | `host/store/operator_launch.go:95` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-104 | `host/store/operator_launch.go:99` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-112 | `host/store/operator_launch.go:111` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-126 | `host/store/operator_launch.go:124` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-129 | `host/store/operator_launch.go:128` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-133 | `host/store/operator_launch.go:134` | `TestOperatorLaunchSession` | KILLED | PENDING |
| Q-host-store-operator_launch-138 | `host/store/operator_launch.go:138` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-147 | `host/store/operator_launch.go:145` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-151 | `host/store/operator_launch.go:149` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-162 | `host/store/operator_launch.go:160` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-173 | `host/store/operator_launch.go:172` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-177 | `host/store/operator_launch.go:176` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-184 | `host/store/operator_launch.go:183` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-188 | `host/store/operator_launch.go:187` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-196 | `host/store/operator_launch.go:195` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-205 | `host/store/operator_launch.go:204` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-215 | `host/store/operator_launch.go:212` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-218 | `host/store/operator_launch.go:216` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-222 | `host/store/operator_launch.go:216` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-254 | `host/store/operator_launch.go:253` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-260 | `host/store/operator_launch.go:258` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-263 | `host/store/operator_launch.go:262` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-267 | `host/store/operator_launch.go:265` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-host-store-operator_launch-271 | `host/store/operator_launch.go:270` | `TestOperatorLaunchReceiptLaw` | KILLED | PENDING |
| Q-cmd-ailang-worldd-operator_launch-56 | `cmd/ailang-worldd/operator_launch.go:56` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| Q-cmd-ailang-worldd-operator_launch-70 | `cmd/ailang-worldd/operator_launch.go:70` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| Q-cmd-ailang-worldd-operator_launch-74 | `cmd/ailang-worldd/operator_launch.go:74` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| Q-cmd-ailang-worldd-operator_launch-88 | `cmd/ailang-worldd/operator_launch.go:88` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| Q-cmd-ailang-worldd-operator_launch-98 | `cmd/ailang-worldd/operator_launch.go:98` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| Q-cmd-ailang-worldd-operator_launch-126 | `cmd/ailang-worldd/operator_launch.go:126` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| Q-cmd-ailang-worldd-operator_launch-131 | `cmd/ailang-worldd/operator_launch.go:131` | `TestOperatorStartOnce/launch` | KILLED | PENDING |
| Q-cmd-ailang-worldd-operator_launch-340 | `cmd/ailang-worldd/operator_launch.go:337` | `TestOperatorStartBrowser` | KILLED | PENDING |
| Q-cmd-ailang-worldd-operator_launch-462 | `cmd/ailang-worldd/operator_launch.go:461` | `TestOperatorStartOnce/launch` | KILLED | PENDING |

## 7. Apply, accept and snapshot

```bash
B=b2b0ebab7bffbcf5f3105aa17e3ced452e480e35
STATE=$HOME/.ailang/state/world-iter444
S=$STATE/plan-scratch
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter444-row243-v4s6
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s6-sprint-plan.md
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
export WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
for M in M1 M2; do
  awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$S/executor-$M.patch"
  test -s "$S/executor-$M.patch" || exit 1
  (cd "$W" && /usr/bin/patch -p1 -E < "$S/executor-$M.patch") || exit 1
  /bin/bash "$S/accept.sh" "$W" "$M" fast || exit 1
  /bin/bash "$S/snapshot.sh" "$W" "$M" || exit 1
done
```

The success signal of `accept.sh` is exit 0 **and** `ACCEPT_OK <milestone> <mode>` after the positive selectors and real-path observations. Cumulative snapshots are under `$STATE/snap/M1` and `$STATE/snap/M2`. Each has the complete base-relative product files plus `STATUS`, `DELETED` and `MANIFEST.sha256`. Product directories can be absent. The snapshot script discovers changed files and creates only required destination parents. The plans are excluded.

Controller commit-time test commands. Use a fresh log directory and the pinned env from §7:

```bash
"$STATE/commit_ms.sh" M1 "$S/M1.msg" "/bin/bash '$S/accept.sh' '$W' M1 fast"
"$STATE/commit_ms.sh" M2 "$S/M2.msg" "/bin/bash '$S/accept.sh' '$W' M2 fast && /bin/bash '$S/go-leg.sh' '$W' '$S/commit-M2-verifygate.raw' 930 go test -count=1 ./host/verifygate/ -timeout 900s"
```

The last command places the whole verifygate check in M2's commit gate. Its known-red base receipt is in `$STATE/base/verifygate-plain.fails`. If the whole leg fails with that base failure, retain the receipt and rerun on the controller's quiet environment before commit. Run §6 for complete fail-set comparisons. CLI and daemon commit-time whole-package calls are omitted.

The final Appendix P was extracted by the controller's awk rule and applied with `patch -p1 -E`.

```text
REHEARSAL /Users/voightkampff/.ailang/state/world-iter444/plan-scratch/rt-c8qPyP
M1 diff -r cmd rc=0
M1 diff -r host rc=0
M1 diff -r docs rc=0
M1 diff -r app SKIPPED (absent)
M1 diff -r full archive rc=0
M1 SHA manifest and byte compare OK (20 paths, 0 deleted)
M2 diff -r cmd rc=0
M2 diff -r host rc=0
M2 diff -r docs rc=0
M2 diff -r app SKIPPED (absent)
M2 diff -r full archive rc=0
M2 SHA manifest and byte compare OK (34 paths, 0 deleted)
ROUNDTRIP_OK
```

````bash snapshot.sh
#!/bin/bash
# snapshot.sh TREE M1|M2; SNAP_STATUS_LIST/ SNAP_PARENT are rehearsal-only.
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1|M2) ;; *) exit 2;; esac
B=b2b0ebab7bffbcf5f3105aa17e3ced452e480e35
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter444/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s6-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s6.json"
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

## Appendix P — patches

### M1

````diff M1
diff --git a/cmd/ailang-worldd/operator_activation_test.go b/cmd/ailang-worldd/operator_activation_test.go
index 66c459a3..25565835 100644
--- a/cmd/ailang-worldd/operator_activation_test.go
+++ b/cmd/ailang-worldd/operator_activation_test.go
@@ -96,9 +96,6 @@ func activationEvidence(t *testing.T, r *worktreeStartRig) [5]any {
 }
 
 func TestOperatorFirstManagedEpisodeProcess(t *testing.T) {
-	if operatorPreparationBudget != 22*time.Second {
-		t.Fatalf("production preparation budget = %v, want 22s", operatorPreparationBudget)
-	}
 	r, c, h, runner := operatorActivationFixture(t)
 	ctx := activationTestContext(t)
 	prepare := func(acquire func(string, string) (*store.OperatorEpisodeLock, error), stops <-chan os.Signal) error {
@@ -558,6 +555,9 @@ func TestOperatorAttendedStartDisposition(t *testing.T) {
 	}
 	before = drifted
 	r.env.terminal.openTerminal = terminal
+	var budgets []time.Duration
+	r.env.operatorDeadline = func(budget time.Duration, _ context.Context) { budgets = append(budgets, budget) }
+
 	code, output = r.command("reconcile-start", id)
 	after := activationEvidence(t, r)
 	if code != exitOK || !strings.Contains(output, "existing evidence retained") {
@@ -567,6 +567,9 @@ func TestOperatorAttendedStartDisposition(t *testing.T) {
 	if !reflect.DeepEqual(before, after) {
 		t.Fatal("disposition modified the retained worktree, branch or private files")
 	}
+	if !reflect.DeepEqual(budgets, []time.Duration{5 * time.Second, 5 * time.Second}) {
+		t.Fatalf("reconcile preflight and writeback deadlines: %v", budgets)
+	}
 	stopped, err := ro.OperatorAction(ctx, id, 0)
 	if err != nil || stopped.State != "failed" || stopped.Owner != "" {
 		t.Fatalf("rollback disposition: %s %q %v", stopped.State, stopped.Owner, err)
diff --git a/cmd/ailang-worldd/operator_deadline_test.go b/cmd/ailang-worldd/operator_deadline_test.go
new file mode 100644
index 00000000..98a2e9a1
--- /dev/null
+++ b/cmd/ailang-worldd/operator_deadline_test.go
@@ -0,0 +1,49 @@
+package main
+
+import (
+	"context"
+	"io"
+	"reflect"
+	"testing"
+	"time"
+)
+
+func TestOperatorProductionDeadlines(t *testing.T) {
+	f := operatorService(t)
+	f.mustInstall(t)
+	r, err := readServiceRecord(f.env.paths)
+	if err != nil {
+		t.Fatal(err)
+	}
+	lock, err := runnerKernelLock(f.config.Credentials + ".operator-runner.lock")
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer lock.Close()
+	for _, override := range []time.Duration{0, time.Minute} {
+		e := f.env
+		e.operatorBudget = override
+		var requested []time.Duration
+		begin := time.Now()
+		e.operatorDeadline = func(raw time.Duration, ctx context.Context) {
+			requested = append(requested, raw)
+			want := raw
+			if override != 0 {
+				want = override
+			}
+			got, ok := ctx.Deadline()
+			if !ok || got.Before(begin.Add(want)) || got.After(time.Now().Add(want)) {
+				t.Errorf("requested %v effective %v: deadline %v (%t)", raw, want, got, ok)
+			}
+		}
+		ticks := make(chan time.Time)
+		close(ticks)
+		h := f.heartbeat
+		if rc := e.runOperatorPreparationLoop(io.Discard, f.config, h, operatorStartSource{workspace: r.Workspace, clock: e.terminal.now}, nil, ticks, lock); rc != exitOK {
+			t.Fatalf("empty real loop: %d", rc)
+		}
+		if !reflect.DeepEqual(requested, []time.Duration{time.Second, 22 * time.Second}) {
+			t.Fatalf("production child contexts: %v", requested)
+		}
+	}
+}
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
index df1ff17f..4d1e55f4 100644
--- a/cmd/ailang-worldd/operator_runner.go
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -33,7 +33,10 @@ func operatorChildBudget(ctx context.Context, budget time.Duration) (context.Con
 	return context.WithTimeout(ctx, budget)
 }
 
-func (e serviceEnv) operatorContext(budget time.Duration) (context.Context, context.CancelFunc) {
+func (e serviceEnv) operatorContext(budget time.Duration) (result context.Context, done context.CancelFunc) {
+	if e.operatorDeadline != nil {
+		defer func() { e.operatorDeadline(budget, result) }()
+	}
 	if e.operatorBudget == 0 {
 		return budgetContext(budget)
 	}
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
index 1b9583ef..b3addb4c 100644
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -2,6 +2,7 @@ package main
 
 import (
 	"bytes"
+	"context"
 	"crypto/sha256"
 	"encoding/hex"
 	"encoding/json"
@@ -86,6 +87,8 @@ type serviceEnv struct {
 	runnerNow       func() time.Time
 	runnerSleep     func(time.Duration)
 	runnerProcesses func() ([]worldProcess, error)
+
+	operatorDeadline func(time.Duration, context.Context)
 }
 
 func realServiceEnv() serviceEnv {
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index 9656fdf7..db80ca2e 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -311,6 +311,11 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	_, errs["NextOperatorStart"] = c.NextOperatorStart(ctx, 1000)
 	errs["ReserveOperatorStep"] = c.ReserveOperatorStep(ctx, OperatorStepKey{}, 1000)
 	errs["ObserveOperatorStep"] = c.ObserveOperatorStep(ctx, OperatorStepKey{}, nil)
+	errs["ReserveOperatorLaunch"] = c.ReserveOperatorLaunch(ctx, OperatorStartRecovery{}, [32]byte{}, nil, nil)
+	errs["AuthorizeOperatorSpawn"] = c.AuthorizeOperatorSpawn(ctx, "x", [32]byte{}, nil, nil)
+	_, _, errs["ReadOperatorSpawnInput"] = c.ReadOperatorSpawnInput(ctx, "x")
+	_, errs["ReadOperatorLaunch"] = c.ReadOperatorLaunch(ctx, "x")
+	errs["RecordOperatorLaunch"] = c.RecordOperatorLaunch(ctx, OperatorLaunchSnapshot{}, OperatorLaunchReceipt{}, nil, 1000)
 	_, errs["OperatorStep"] = c.OperatorStep(ctx, OperatorStepKey{})
 	_, errs["ReserveOperatorWorktree"] = c.ReserveOperatorWorktree(ctx, OperatorMintAttempt{}, "", nil, func() int64 { return 1000 })
 	_, _, errs["ReadOperatorWorktree"] = c.ReadOperatorWorktree(ctx, "x")
diff --git a/host/store/episode_starts.go b/host/store/episode_starts.go
index 2650d34e..0a2692ab 100644
--- a/host/store/episode_starts.go
+++ b/host/store/episode_starts.go
@@ -131,7 +131,14 @@ func (c *Credentials) AppendEpisodeLaunch(ctx context.Context, episode string, l
 	if err != nil {
 		return err
 	}
-	defer func() { _ = tx.Rollback() }()
+	defer tx.Rollback()
+	if err := appendEpisodeLaunch(ctx, tx, episode, launch); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
+
+func appendEpisodeLaunch(ctx context.Context, tx *sql.Tx, episode string, launch EpisodeLaunch) error {
 	var data string
 	if err := tx.QueryRowContext(ctx, "SELECT start_json FROM episode_starts WHERE episode = ?", episode).Scan(&data); err != nil {
 		return err
@@ -162,7 +169,7 @@ func (c *Credentials) AppendEpisodeLaunch(ctx context.Context, episode string, l
 	if _, err := tx.ExecContext(ctx, "UPDATE episode_starts SET start_json = ? WHERE episode = ?", string(encoded), episode); err != nil {
 		return err
 	}
-	return tx.Commit()
+	return nil
 }
 
 type EpisodeResume struct {
diff --git a/host/store/operator_actions.go b/host/store/operator_actions.go
index d7777c13..176d93c8 100644
--- a/host/store/operator_actions.go
+++ b/host/store/operator_actions.go
@@ -184,9 +184,9 @@ func (c *Credentials) PruneOperatorActions(ctx context.Context, now int64) error
 	if err := expireOperatorActions(ctx, tx, now); err != nil {
 		return err
 	}
-	if _, err := tx.ExecContext(ctx, `UPDATE operator_actions SET action = NULL, challenge = NULL, artifact = NULL, approved_by = NULL
+	if _, err := tx.ExecContext(ctx, `UPDATE operator_actions AS a SET action = NULL, challenge = NULL, artifact = NULL, approved_by = NULL
 		WHERE terminal_at < ? AND NOT EXISTS
-			(SELECT 1 FROM operator_action_steps WHERE action_id = operator_actions.id AND phase = 'reserved')`, now-30*86400); err != nil {
+			(SELECT 1 FROM operator_action_steps WHERE action_id = a.id AND phase = 'reserved') AND NOT (`+operatorLaunchUnsettled+`)`, now-30*86400); err != nil {
 		return err
 	}
 	return tx.Commit()
diff --git a/host/store/operator_cas.go b/host/store/operator_cas.go
index 20db2e59..2e1410e0 100644
--- a/host/store/operator_cas.go
+++ b/host/store/operator_cas.go
@@ -146,7 +146,9 @@ func (c *Credentials) ClaimOperatorAction(ctx context.Context, cas OperatorCAS,
 	}
 	res, err := tx.ExecContext(ctx, `UPDATE operator_actions AS a SET state = 'claimed', owner = ?, generation = generation + 1
 		WHERE id = ? AND digest = ? AND generation = ? AND state = 'confirmed' AND expires_at > ? AND owner IS NULL
-		AND NOT EXISTS (SELECT 1 FROM operator_actions b WHERE b.episode_id = a.episode_id AND b.id <> a.id AND b.owner IS NOT NULL)`,
+		AND NOT EXISTS (SELECT 1 FROM operator_actions b WHERE b.episode_id = a.episode_id AND b.id <> a.id AND b.owner IS NOT NULL)
+		AND NOT EXISTS (SELECT 1 FROM operator_action_steps ls JOIN operator_actions b ON b.id=ls.action_id
+			WHERE b.episode_id=a.episode_id AND b.state='succeeded' AND ls.kind='launch' AND ls.phase='observed' AND ls.version=3)`,
 		runner.Nonce, cas.ID, cas.Digest[:], cas.Generation, now)
 	if err := operatorChangedOne(res, err); err != nil {
 		return err
diff --git a/host/store/operator_cas_test.go b/host/store/operator_cas_test.go
index 199c887a..23a8d9c1 100644
--- a/host/store/operator_cas_test.go
+++ b/host/store/operator_cas_test.go
@@ -142,7 +142,7 @@ func TestOperatorStepRestart(t *testing.T) {
 		t.Fatal(err)
 	}
 	cas.Generation++
-	key := OperatorStepKey{Action: cas, Owner: runner.Nonce, Kind: "launch", Identity: "run-one", Version: 1}
+	key := OperatorStepKey{Action: cas, Owner: runner.Nonce, Kind: "checks", Identity: "run-one", Version: 1}
 	if err := c.ReserveOperatorStep(ctx, key, 1600); !errors.Is(err, ErrOperatorConflict) {
 		t.Fatalf("reserve at expiry: %v", err)
 	}
@@ -161,7 +161,7 @@ func TestOperatorStepRestart(t *testing.T) {
 	if err != nil || reservation.Phase != "reserved" || reservation.Version != 1 || reservation.Owner != runner.Nonce || reservation.Receipt != nil {
 		t.Fatalf("lost reservation on restart: %+v %v", reservation, err)
 	}
-	receipt := []byte("observed-test-launch")
+	receipt := []byte("observed-test-checks")
 	wrong := key
 	wrong.Version++
 	if err := c.ObserveOperatorStep(ctx, wrong, receipt); !errors.Is(err, ErrOperatorConflict) {
diff --git a/host/store/operator_episode_gc.go b/host/store/operator_episode_gc.go
index 95dbb2e3..7615ca6a 100644
--- a/host/store/operator_episode_gc.go
+++ b/host/store/operator_episode_gc.go
@@ -21,6 +21,6 @@ func EpisodeHasUnresolvedOperatorSteps(ctx context.Context, path, episode string
 	err = db.QueryRowContext(ctx, `SELECT EXISTS (
 		SELECT 1 FROM operator_actions a JOIN operator_action_steps s ON s.action_id = a.id
 		WHERE a.episode_id = ? AND s.kind <> 'operation'
-		AND (s.phase = 'reserved' OR a.terminal_at IS NULL))`, episode).Scan(&unresolved)
+		AND (s.phase = 'reserved' OR a.terminal_at IS NULL OR (`+operatorLaunchUnsettled+`)))`, episode).Scan(&unresolved)
 	return unresolved, err
 }
diff --git a/host/store/operator_files.go b/host/store/operator_files.go
index 699b65cf..aa56f619 100644
--- a/host/store/operator_files.go
+++ b/host/store/operator_files.go
@@ -21,9 +21,13 @@ type OperatorFileState struct {
 }
 
 func readOperatorFiles(ctx context.Context, tx *sql.Tx, id string) (OperatorFileState, approval.OperatorAction, error) {
+	return readOperatorFilesAt(ctx, tx, id, false)
+}
+
+func readOperatorFilesAt(ctx context.Context, tx *sql.Tx, id string, launch bool) (OperatorFileState, approval.OperatorAction, error) {
 	var state OperatorFileState
 	var action approval.OperatorAction
-	mint, err := readOperatorMint(ctx, tx, id, true)
+	mint, err := readOperatorMint(ctx, tx, id, true, launch)
 	if err != nil {
 		return state, action, err
 	}
diff --git a/host/store/operator_launch.go b/host/store/operator_launch.go
new file mode 100644
index 00000000..c450a0af
--- /dev/null
+++ b/host/store/operator_launch.go
@@ -0,0 +1,268 @@
+package store
+
+import (
+	"context"
+	"crypto/sha256"
+	"database/sql"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"strings"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+type OperatorLaunchReceipt struct {
+	Run, Executable string
+	PID             int
+	Birth           string
+	Status          string
+	Exit            int
+}
+
+// An observed start with no terminal observation still excludes effects.
+const operatorLaunchUnsettled = `a.state='succeeded' AND EXISTS (
+	SELECT 1 FROM operator_action_steps ls WHERE ls.action_id=a.id
+	AND ls.kind='launch' AND ls.phase='observed' AND ls.version=3)`
+
+func operatorLaunchDigest(s string) bool {
+	value, ok := strings.CutPrefix(s, "sha256:")
+	return ok && operatorDigestString(value)
+}
+
+type OperatorLaunchSnapshot struct {
+	Key     OperatorStepKey
+	Phase   string
+	Receipt OperatorLaunchReceipt
+}
+
+func operatorLaunchKey(s OperatorStartRecovery, handoff [32]byte) OperatorStepKey {
+	return OperatorStepKey{Action: s.CAS, Owner: s.Owner, Kind: "launch", Identity: fmt.Sprintf("%s:%x", s.CAS.ID, sha256.Sum256(handoff[:])), Version: 1}
+}
+
+// ReserveOperatorLaunch is called under the episode and lifecycle locks. It
+// consumes the exact prepared snapshot and rechecks the enrolled proof and
+// current session inside the transaction that makes the reservation durable.
+func (c *Credentials) ReserveOperatorLaunch(ctx context.Context, reviewed OperatorStartRecovery, handoff [32]byte, proof *approval.OperatorMintProof, clock func() int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if clock == nil || c.lock == nil || handoff == ([32]byte{}) {
+		return ErrOperatorInput
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	current, row, err := readOperatorStartRecovery(ctx, tx, reviewed.CAS.ID)
+	now := clock()
+	if err == nil && now >= row.ExpiresAt {
+		return ErrOperatorExpired
+	}
+	if err != nil || current != reviewed || current.Window != OperatorStartPrepared {
+		return ErrOperatorConflict
+	}
+	if _, _, err := operatorSpawnInput(ctx, tx, row, proof, now, false); err != nil {
+		return err
+	}
+	if err := reserveOperatorStep(ctx, tx, operatorLaunchKey(current, handoff), now); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
+
+// ReadOperatorSpawnInput reads the frozen start and its exact credential for
+// first dispatch. It leaves all persisted evidence as it found it.
+func (c *Credentials) ReadOperatorSpawnInput(ctx context.Context, id string) (OperatorFileState, approval.OperatorAction, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorFileState{}, approval.OperatorAction{}, err
+	}
+	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
+	if err != nil {
+		return OperatorFileState{}, approval.OperatorAction{}, err
+	}
+	defer tx.Rollback()
+	s, _, err := readOperatorLaunch(ctx, tx, id)
+	if err != nil || s.Phase != "reserved" || s.Key.Version != 1 {
+		return OperatorFileState{}, approval.OperatorAction{}, ErrOperatorConflict
+	}
+	return readOperatorFilesAt(ctx, tx, id, true)
+}
+
+// AuthorizeOperatorSpawn commits one dispatch grant. The process that receives
+// the private handoff must hold the credentials writer through cmd.Start.
+func (c *Credentials) AuthorizeOperatorSpawn(ctx context.Context, id string, handoff [32]byte, proof *approval.OperatorMintProof, clock func() int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if clock == nil || c.lock == nil || handoff == ([32]byte{}) {
+		return ErrOperatorInput
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	s, row, err := readOperatorLaunch(ctx, tx, id)
+	if err != nil || s.Phase != "reserved" || s.Key != operatorLaunchKey(OperatorStartRecovery{CAS: s.Key.Action, Owner: row.Owner}, handoff) {
+		return ErrOperatorConflict
+	}
+	if _, _, err := operatorSpawnInput(ctx, tx, row, proof, clock(), true); err != nil {
+		return err
+	}
+	res, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET version=version+1
+		WHERE action_id=? AND kind='launch' AND identity=? AND phase='reserved' AND version=1 AND owner=?`,
+		id, s.Key.Identity, s.Key.Owner)
+	if err := operatorChangedOne(res, err); err != nil {
+		return err
+	}
+	return tx.Commit()
+}
+
+func operatorSpawnInput(ctx context.Context, tx *sql.Tx, row OperatorActionRow, proof *approval.OperatorMintProof, now int64, reserved bool) (OperatorFileState, approval.OperatorAction, error) {
+	files, a, err := readOperatorFilesAt(ctx, tx, row.ID, reserved)
+	if err != nil {
+		return files, a, ErrOperatorConflict
+	}
+	if err := admitOperatorFiles(ctx, tx, row, files, a, proof, now); err != nil {
+		return files, a, err
+	}
+	var other bool
+	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(
+		SELECT 1 FROM session_credentials WHERE episode_id=? AND credential_id<>?)`,
+		row.Episode, files.Session.CredentialID).Scan(&other); err != nil {
+		return files, a, err
+	}
+	if other || !operatorLaunchDigest(mintFact(a, "launch_binary_sha256")) {
+		return files, a, ErrOperatorConflict
+	}
+	return files, a, nil
+}
+
+// ReadOperatorLaunch uses one read-only snapshot and never makes a replay
+// decision from a PID, the heartbeat, or the current presence of a credential.
+func (c *Credentials) ReadOperatorLaunch(ctx context.Context, id string) (OperatorLaunchSnapshot, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorLaunchSnapshot{}, err
+	}
+	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
+	if err != nil {
+		return OperatorLaunchSnapshot{}, err
+	}
+	defer tx.Rollback()
+	s, _, err := readOperatorLaunch(ctx, tx, id)
+	return s, err
+}
+
+func readOperatorLaunch(ctx context.Context, tx *sql.Tx, id string) (OperatorLaunchSnapshot, OperatorActionRow, error) {
+	var s OperatorLaunchSnapshot
+	row, err := scanOperatorAction(tx.QueryRowContext(ctx, "SELECT "+operatorColumns+" FROM operator_actions WHERE id=?", id))
+	if err != nil {
+		return s, row, err
+	}
+	s.Key = OperatorStepKey{Action: OperatorCAS{row.ID, [32]byte(row.Digest), row.Generation}, Kind: "launch"}
+	var n int
+	var raw []byte
+	err = tx.QueryRowContext(ctx, `SELECT identity,phase,version,owner,receipt,
+		(SELECT count(*) FROM operator_action_steps WHERE action_id=? AND kind='launch')
+		FROM operator_action_steps WHERE action_id=? AND kind='launch' LIMIT 1`, id, id).
+		Scan(&s.Key.Identity, &s.Phase, &s.Key.Version, &s.Key.Owner, &raw, &n)
+	if errors.Is(err, sql.ErrNoRows) && row.Kind == "start" {
+		return s, row, ErrOperatorMissing
+	}
+	identity, prefixed := strings.CutPrefix(s.Key.Identity, id+":")
+	if err != nil || n != 1 || row.Kind != "start" || s.Key.Owner == "" || !prefixed || !operatorDigestString(identity) {
+		return s, row, ErrOperatorConflict
+	}
+	if s.Phase == "reserved" && (s.Key.Version == 1 || s.Key.Version == 2) && len(raw) == 0 && s.Key.Owner == row.Owner &&
+		(row.State == "running" || row.State == "uncertain") {
+		return s, row, nil
+	}
+	if s.Phase != "observed" || json.Unmarshal(raw, &s.Receipt) != nil || s.Receipt.Run != id ||
+		!operatorLaunchDigest(s.Receipt.Executable) {
+		return s, row, ErrOperatorConflict
+	}
+	switch s.Receipt.Status {
+	case "launched":
+		if s.Key.Version == 3 && row.State == "succeeded" && s.Receipt.PID > 0 && s.Receipt.Birth != "" {
+			return s, row, nil
+		}
+	case "spawn-refused":
+		if (s.Key.Version == 2 || s.Key.Version == 3) && row.State == "failed" && s.Receipt.PID == 0 && s.Receipt.Birth == "" {
+			return s, row, nil
+		}
+	case "stopped", "failed", "capped":
+		if s.Key.Version == 4 && row.State == "succeeded" && s.Receipt.PID > 0 && s.Receipt.Birth != "" {
+			return s, row, nil
+		}
+	}
+	return s, row, ErrOperatorConflict
+}
+
+// RecordOperatorLaunch compares the exact read-back and never replaces a
+// process identity. Recording an admitted child's outcome ignores expiry.
+func (c *Credentials) RecordOperatorLaunch(ctx context.Context, reviewed OperatorLaunchSnapshot, r OperatorLaunchReceipt, history *EpisodeLaunch, now int64) error {
+	if err := c.operatorAccess(ctx); err != nil {
+		return err
+	}
+	if c.lock == nil {
+		return ErrOperatorInput
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return err
+	}
+	defer tx.Rollback()
+	s, row, err := readOperatorLaunch(ctx, tx, reviewed.Key.Action.ID)
+	a, decodeErr := decodeOperatorMintAction(row)
+	if err != nil || decodeErr != nil || s != reviewed || r.Run != row.ID || r.Executable != mintFact(a, "launch_binary_sha256") {
+		return ErrOperatorConflict
+	}
+	var terminal string
+	switch s.Phase {
+	case "reserved":
+		switch {
+		case s.Key.Version == 2 && r.Status == "launched" && r.PID > 0 && r.Birth != "" && history == nil:
+			terminal = "succeeded"
+		case r.Status == "spawn-refused" && r.PID == 0 && r.Birth == "" && history == nil:
+			terminal = "failed"
+		default:
+			return ErrOperatorInput
+		}
+	case "observed":
+		if s.Receipt.Status != "launched" || r.PID != s.Receipt.PID || r.Birth != s.Receipt.Birth ||
+			history == nil || history.Mode != "headless" || history.Exit != r.Exit {
+			return ErrOperatorConflict
+		}
+		switch r.Status {
+		case "stopped", "failed", "capped":
+		default:
+			return ErrOperatorInput
+		}
+	}
+	raw, err := json.Marshal(r)
+	if err != nil || len(raw) > 2048 {
+		return ErrOperatorInput
+	}
+	change, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET phase='observed',version=version+1,receipt=?
+		WHERE action_id=? AND kind='launch' AND identity=? AND phase=? AND version=? AND owner=?`,
+		raw, row.ID, s.Key.Identity, s.Phase, s.Key.Version, s.Key.Owner)
+	if err := operatorChangedOne(change, err); err != nil {
+		return err
+	}
+	if terminal != "" {
+		change, err = tx.ExecContext(ctx, `UPDATE operator_actions SET state=?,owner=NULL,terminal_at=?,generation=generation+1
+			WHERE id=? AND digest=? AND generation=? AND owner=? AND state='running'`,
+			terminal, now, row.ID, row.Digest, row.Generation, s.Key.Owner)
+		if err := operatorChangedOne(change, err); err != nil {
+			return err
+		}
+	}
+	if history != nil {
+		if err := appendEpisodeLaunch(ctx, tx, row.Episode, *history); err != nil {
+			return err
+		}
+	}
+	return tx.Commit()
+}
diff --git a/host/store/operator_launch_test.go b/host/store/operator_launch_test.go
new file mode 100644
index 00000000..5ee7b307
--- /dev/null
+++ b/host/store/operator_launch_test.go
@@ -0,0 +1,92 @@
+package store
+
+import (
+	"errors"
+	"reflect"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+func TestOperatorLaunchReceiptLaw(t *testing.T) {
+	digest := "sha256:" + hex64("fixture executable")
+	f := newOperatorPinFixtureWithFacts(t, false, []approval.OperatorResolved{{Name: "launch_binary_sha256", Value: digest}}, startRecoveryTestContext(t))
+	f.mint(t)
+	id := f.in.Attempt.CAS.ID
+	for _, value := range []string{"files-identity", OperatorFilesPublished} {
+		files, _, err := f.c.ReadOperatorFiles(f.ctx, id)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if err := f.c.RecordOperatorFiles(f.ctx, files, value, f.proof, func() int64 { return pinNow }); err != nil {
+			t.Fatal(err)
+		}
+	}
+	before := f.snapshot(t)
+	prepared, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
+	if err != nil || prepared.Window != OperatorStartPrepared {
+		t.Fatalf("prepared: %+v %v", prepared, err)
+	}
+	handoff := [32]byte{1}
+	now := func() int64 { return pinNow }
+	f.assertRefused(t, before, f.c.ReserveOperatorStep(f.ctx, operatorLaunchKey(prepared, handoff), pinNow), ErrOperatorInput)
+	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, prepared, handoff, nil, now), ErrOperatorInput)
+	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, prepared, handoff, f.proof, func() int64 { return before.action.ExpiresAt }), ErrOperatorExpired)
+	wrong := prepared
+	wrong.CAS.Generation++
+	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, wrong, handoff, f.proof, now), ErrOperatorConflict)
+	if err := f.c.ReserveOperatorLaunch(f.ctx, prepared, handoff, f.proof, now); err != nil {
+		t.Fatal(err)
+	}
+	reserved, err := f.c.ReadOperatorLaunch(f.ctx, id)
+	if err != nil || reserved.Phase != "reserved" || reserved.Receipt != (OperatorLaunchReceipt{}) {
+		t.Fatalf("reserved: %+v %v", reserved, err)
+	}
+	before = f.snapshot(t)
+	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, prepared, handoff, f.proof, now), ErrOperatorConflict)
+	f.assertRefused(t, before, f.c.ObserveOperatorStep(f.ctx, reserved.Key, []byte("fake")), ErrOperatorInput)
+	f.assertRefused(t, before, f.c.AuthorizeOperatorSpawn(f.ctx, id, [32]byte{2}, f.proof, now), ErrOperatorConflict)
+	if files, action, err := f.c.ReadOperatorSpawnInput(f.ctx, id); err != nil || files.Step.Action != prepared.CAS || action.ID != id {
+		t.Fatalf("spawn input readback: %+v %s %v", files.Step, action.ID, err)
+	}
+	if err := f.c.AuthorizeOperatorSpawn(f.ctx, id, handoff, f.proof, now); err != nil {
+		t.Fatal(err)
+	}
+	before = f.snapshot(t)
+	f.assertRefused(t, before, f.c.AuthorizeOperatorSpawn(f.ctx, id, handoff, f.proof, now), ErrOperatorConflict)
+	reserved, err = f.c.ReadOperatorLaunch(f.ctx, id)
+	if err != nil || reserved.Key.Version != 2 || reserved.Phase != "reserved" {
+		t.Fatalf("dispatch grant: %+v %v", reserved, err)
+	}
+
+	r := OperatorLaunchReceipt{Run: id, Executable: digest, PID: 42, Birth: "kernel-start", Status: "launched"}
+	wrongRun := r
+	wrongRun.Run += "x"
+	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, reserved, wrongRun, nil, pinNow), ErrOperatorConflict)
+	if err := f.c.RecordOperatorLaunch(f.ctx, reserved, r, nil, before.action.ExpiresAt+1); err != nil {
+		t.Fatal(err)
+	}
+	before = f.snapshot(t)
+	live, err := f.c.ReadOperatorLaunch(f.ctx, id)
+	if err != nil || live.Receipt != r || before.action.State != "succeeded" || !reflect.DeepEqual(before, f.snapshot(t)) {
+		t.Fatalf("observed: %+v %v", live, err)
+	}
+	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, reserved, r, nil, pinNow), ErrOperatorConflict)
+	r.Status = "capped"
+	r.Exit = -1
+	history := EpisodeLaunch{Mode: "headless", Exit: r.Exit}
+	wrongBirth := r
+	wrongBirth.Birth += "x"
+	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, live, wrongBirth, &history, pinNow), ErrOperatorConflict)
+	if err := f.c.RecordOperatorLaunch(f.ctx, live, r, &history, pinNow); err != nil {
+		t.Fatal(err)
+	}
+	finished, err := f.c.ReadOperatorLaunch(f.ctx, id)
+	start, startErr := f.c.EpisodeStart(f.ctx, before.action.Episode)
+	if err != nil || startErr != nil || finished.Receipt != r || len(start.Launches) != 1 || start.Launches[0].Exit != r.Exit {
+		t.Fatalf("completed: %+v %v, history %v %v", finished, err, start, startErr)
+	}
+	if err := f.c.RecordOperatorLaunch(f.ctx, live, r, &history, pinNow); !errors.Is(err, ErrOperatorConflict) {
+		t.Fatalf("completion replay: %v", err)
+	}
+}
diff --git a/host/store/operator_mint.go b/host/store/operator_mint.go
index 0267140f..c71cd7e7 100644
--- a/host/store/operator_mint.go
+++ b/host/store/operator_mint.go
@@ -243,10 +243,10 @@ func (c *Credentials) ReadOperatorMint(ctx context.Context, id string) (Operator
 	if err := c.operatorAccess(ctx); err != nil {
 		return OperatorMintReceipt{}, err
 	}
-	return readOperatorMint(ctx, c.db, id, false)
+	return readOperatorMint(ctx, c.db, id, false, false)
 }
 
-func readOperatorMint(ctx context.Context, db operatorMintReader, id string, fileStage bool) (OperatorMintReceipt, error) {
+func readOperatorMint(ctx context.Context, db operatorMintReader, id string, fileStage, launchStage bool) (OperatorMintReceipt, error) {
 	row, err := scanOperatorAction(db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, id))
 	if err != nil {
 		return OperatorMintReceipt{}, err
@@ -263,8 +263,9 @@ func readOperatorMint(ctx context.Context, db operatorMintReader, id string, fil
 	var receipt []byte
 	err = db.QueryRowContext(ctx, `SELECT version, receipt FROM operator_action_steps WHERE kind = 'mint'
 		AND identity = ? AND action_id = ? AND phase = 'observed' AND owner = ?
-		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree','mint') AND NOT (? AND kind = 'files'))`,
-		id, id, row.Owner, id, fileStage).Scan(&r.StepVersion, &receipt)
+		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree','mint') AND NOT (? AND kind = 'files')
+		AND NOT (? AND kind = 'launch' AND phase='reserved'))`,
+		id, id, row.Owner, id, fileStage, launchStage).Scan(&r.StepVersion, &receipt)
 	if err != nil {
 		return OperatorMintReceipt{}, fmt.Errorf("mint evidence unavailable or later effect reserved: %w", err)
 	}
@@ -299,7 +300,7 @@ func (c *Credentials) AbortOperatorMint(ctx context.Context, reviewed OperatorMi
 		return err
 	}
 	defer tx.Rollback()
-	current, err := readOperatorMint(ctx, tx, reviewed.CAS.ID, false)
+	current, err := readOperatorMint(ctx, tx, reviewed.CAS.ID, false, false)
 	if err != nil {
 		return err
 	}
diff --git a/host/store/operator_runner.go b/host/store/operator_runner.go
index 98dcbe77..375383ea 100644
--- a/host/store/operator_runner.go
+++ b/host/store/operator_runner.go
@@ -102,7 +102,7 @@ func ReadOperatorEpisodeOwner(ctx context.Context, path, episode string) (string
 		return "", nil
 	}
 	var id string
-	err = db.QueryRowContext(ctx, `SELECT id FROM operator_actions WHERE episode_id=? AND owner IS NOT NULL`, episode).Scan(&id)
+	err = db.QueryRowContext(ctx, `SELECT id FROM operator_actions a WHERE episode_id=? AND (owner IS NOT NULL OR (`+operatorLaunchUnsettled+`)) LIMIT 1`, episode).Scan(&id)
 	if errors.Is(err, sql.ErrNoRows) {
 		return "", nil
 	}
diff --git a/host/store/operator_s5e4_test.go b/host/store/operator_s5e4_test.go
index cd672d0a..539772e9 100644
--- a/host/store/operator_s5e4_test.go
+++ b/host/store/operator_s5e4_test.go
@@ -23,6 +23,11 @@ type operatorPinFixture struct {
 }
 
 func newOperatorPinFixture(t *testing.T, allowance bool, contexts ...context.Context) operatorPinFixture {
+	t.Helper()
+	return newOperatorPinFixtureWithFacts(t, allowance, nil, contexts...)
+}
+
+func newOperatorPinFixtureWithFacts(t *testing.T, allowance bool, prefix []approval.OperatorResolved, contexts ...context.Context) operatorPinFixture {
 	t.Helper()
 	c, ctx := openTestCredentials(t), testCtx(t)
 	if len(contexts) > 0 {
@@ -54,6 +59,7 @@ func newOperatorPinFixture(t *testing.T, allowance bool, contexts ...context.Con
 	if err != nil {
 		t.Fatal(err)
 	}
+	a.Resolution = append(prefix, a.Resolution...)
 	proof := signOperatorPinAction(t, c, row, a)
 	in.Attempt.CAS.Digest, err = approval.OperatorSignatureDigest(a)
 	if err != nil {
diff --git a/host/store/operator_start_recovery.go b/host/store/operator_start_recovery.go
index e051ab47..11b48072 100644
--- a/host/store/operator_start_recovery.go
+++ b/host/store/operator_start_recovery.go
@@ -203,7 +203,7 @@ func (c *Credentials) NextOwnedOperatorStart(ctx context.Context, after string)
 		return OperatorActionRow{}, err
 	}
 	return scanOperatorAction(c.db.QueryRowContext(ctx, `SELECT `+operatorColumns+
-		` FROM operator_actions WHERE kind='start' AND owner IS NOT NULL AND id>? ORDER BY id LIMIT 1`, after))
+		` FROM operator_actions a WHERE kind='start' AND (owner IS NOT NULL OR (`+operatorLaunchUnsettled+`)) AND id>? ORDER BY id LIMIT 1`, after))
 }
 
 // MarkOperatorMintUncertain preserves every receipt and the old owner while
diff --git a/host/store/operator_steps.go b/host/store/operator_steps.go
index 2652788e..aaf601b2 100644
--- a/host/store/operator_steps.go
+++ b/host/store/operator_steps.go
@@ -18,7 +18,7 @@ func (c *Credentials) ReserveOperatorStep(ctx context.Context, key OperatorStepK
 	if err := c.operatorAccess(ctx); err != nil {
 		return err
 	}
-	if key.Kind == "operation" || key.Identity == "" || key.Owner == "" || key.Version != 1 {
+	if (key.Kind == "operation" || key.Kind == "launch") || key.Identity == "" || key.Owner == "" || key.Version != 1 {
 		return ErrOperatorInput
 	}
 	return reserveOperatorStep(ctx, c.db, key, now)
@@ -46,7 +46,7 @@ func (c *Credentials) ObserveOperatorStep(ctx context.Context, key OperatorStepK
 	if err := c.operatorAccess(ctx); err != nil {
 		return err
 	}
-	if len(receipt) == 0 || len(receipt) > MaxSurfaceApprovalArtifactBytes {
+	if key.Kind == "launch" || len(receipt) == 0 || len(receipt) > MaxSurfaceApprovalArtifactBytes {
 		return ErrOperatorInput
 	}
 	return observeOperatorStep(ctx, c.db, key, receipt)
diff --git a/host/worktree/command_timeout_test.go b/host/worktree/command_timeout_test.go
new file mode 100644
index 00000000..f9fd7df3
--- /dev/null
+++ b/host/worktree/command_timeout_test.go
@@ -0,0 +1,23 @@
+package worktree
+
+import (
+	"context"
+	"strings"
+	"testing"
+	"time"
+)
+
+func TestWithCommandTimeoutClampsGitCap(t *testing.T) {
+	for _, tc := range []struct{ requested, want time.Duration }{
+		{time.Hour, 30 * time.Second},
+		{2 * time.Second, 2 * time.Second},
+	} {
+		ctx, cancel := context.WithCancel(t.Context())
+		cancel()
+		ctx = WithCommandTimeout(ctx, tc.requested)
+		_, err := Output(ctx, "/usr/bin/false", []string{"noop"})
+		if err == nil || !strings.Contains(err.Error(), "timed out after "+tc.want.String()+":") {
+			t.Fatalf("git cap for %v: %v", tc.requested, err)
+		}
+	}
+}
diff --git a/host/worktree/worktree.go b/host/worktree/worktree.go
index 6a76e859..779282c3 100644
--- a/host/worktree/worktree.go
+++ b/host/worktree/worktree.go
@@ -48,7 +48,7 @@ type commandTimeoutKey struct{}
 // WithCommandTimeout overrides the per-invocation cap. The parent deadline
 // still bounds the complete operation.
 func WithCommandTimeout(ctx context.Context, timeout time.Duration) context.Context {
-	return context.WithValue(ctx, commandTimeoutKey{}, timeout)
+	return context.WithValue(ctx, commandTimeoutKey{}, min(timeout, Timeout))
 }
 
 // Env is the child environment: environ without registry, GIT_* or WORLD_SESSION variables.
````

### M2

````diff M2
diff --git a/cmd/ailang-worldd/cli_credential_writers_test.go b/cmd/ailang-worldd/cli_credential_writers_test.go
index 8ab75fae..0f52c684 100644
--- a/cmd/ailang-worldd/cli_credential_writers_test.go
+++ b/cmd/ailang-worldd/cli_credential_writers_test.go
@@ -49,7 +49,7 @@ func TestCLICredentialWriterSites(t *testing.T) {
 			})
 		}
 	}
-	want := map[string]string{"runnerWrite": "OpenCredentials", "openCredentialsForSession": "OpenCredentials", "launchEpisode": "OpenCredentials", "finishRevoke": "OpenCredentials", "putFinishVerbRequest": "OpenCredentials", "runAllowanceGrant": "OpenCredentials", "runAllowanceRevoke": "OpenCredentials"}
+	want := map[string]string{"executeOperatorLaunch": "OpenCredentials", "runnerWrite": "OpenCredentials", "openCredentialsForSession": "OpenCredentials", "launchEpisode": "OpenCredentials", "finishRevoke": "OpenCredentials", "putFinishVerbRequest": "OpenCredentials", "runAllowanceGrant": "OpenCredentials", "runAllowanceRevoke": "OpenCredentials"}
 	if !reflect.DeepEqual(sites, want) || readOnly < 20 {
 		t.Fatalf("writer sites %v (read-only opens %d), want exactly %v and at least 20 read-only opens", sites, readOnly, want)
 	}
diff --git a/cmd/ailang-worldd/main.go b/cmd/ailang-worldd/main.go
index bf7b7f06..75dd7701 100644
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -389,6 +389,12 @@ func run(args []string, stdout, stderr io.Writer) int {
 		}
 		return runOperatorRunner(rest[1:], stderr, realServiceEnv())
 	}
+	if rest[0] == "operator-launch-worker" {
+		if addrGiven {
+			return exitUsage
+		}
+		return realServiceEnv().runOperatorLaunchWorker(rest[1:], stderr)
+	}
 	switch verb := rest[0]; verb {
 	case "version":
 		return runVersion(rest[1:], stdout, stderr)
diff --git a/cmd/ailang-worldd/operator_activation_test.go b/cmd/ailang-worldd/operator_activation_test.go
index 25565835..41ec39d8 100644
--- a/cmd/ailang-worldd/operator_activation_test.go
+++ b/cmd/ailang-worldd/operator_activation_test.go
@@ -33,7 +33,12 @@ func activationTestContext(t *testing.T) context.Context {
 
 func operatorActivationFixture(t *testing.T) (*worktreeStartRig, runnerConfig, store.OperatorHeartbeat, *os.File) {
 	t.Helper()
-	r := newWorktreeStartRigWithContext(t, false, activationTestContext(t))
+	return operatorActivationFixtureWithAgent(t, "")
+}
+
+func operatorActivationFixtureWithAgent(t *testing.T, script string) (*worktreeStartRig, runnerConfig, store.OperatorHeartbeat, *os.File) {
+	t.Helper()
+	r := newWorktreeStartRigWithAgent(t, false, activationTestContext(t), script)
 	r.env.operatorBudget = activationTestBudget
 	c := r.config
 	agent := operatorFileFact(r.action, "launch_binary")
@@ -320,7 +325,10 @@ func TestOperatorFirstManagedEpisodeProcess(t *testing.T) {
 		t.Fatalf("next tick prepared the same row twice: %v", err)
 	}
 	if _, err := os.Stat(filepath.Join(r.bin, "ran")); !errors.Is(err, os.ErrNotExist) {
-		t.Fatalf("spawn crossed the S6 boundary: %v", err)
+		t.Fatalf("preparation or readback spawned before a durable S6 reservation: %v", err)
+	}
+	if _, err := ro.ReadOperatorLaunch(ctx, row.ID); !errors.Is(err, store.ErrOperatorMissing) {
+		t.Fatalf("preparation or readback changed the launch ledger: %v", err)
 	}
 }
 
diff --git a/cmd/ailang-worldd/operator_launch.go b/cmd/ailang-worldd/operator_launch.go
new file mode 100644
index 00000000..5427b022
--- /dev/null
+++ b/cmd/ailang-worldd/operator_launch.go
@@ -0,0 +1,529 @@
+package main
+
+import (
+	"bytes"
+	"context"
+	"crypto/rand"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"net/http"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"slices"
+	"strconv"
+	"strings"
+	"syscall"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/agentrun"
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+const operatorLaunchBudget = 5 * time.Second
+
+var errOperatorLaunchHandled = errors.New("operator launch tick handled")
+
+type operatorLaunchJob struct {
+	config runnerConfig
+	row    store.OperatorActionRow
+	wall   time.Duration
+	secret [32]byte
+	locks  [3]*os.File // runner, episode, lifecycle
+}
+
+type operatorLaunchInputs struct {
+	token, dir, sessions, binary, digest, prompt string
+	argv                                         []string
+	proof                                        *approval.OperatorMintProof
+	wall                                         time.Duration
+}
+
+func (e serviceEnv) launchInputs(ctx context.Context, c runnerConfig, root string, row store.OperatorActionRow,
+	files store.OperatorFileState, a approval.OperatorAction, ro *store.Credentials) (operatorLaunchInputs, error) {
+	in := operatorLaunchInputs{binary: operatorFileFact(a, "launch_binary"), sessions: filepath.Dir(files.Mint.TokenPath)}
+	digest, err := authority.OperatorStartLaunchDigest(a)
+	if err != nil {
+		return in, err
+	}
+	in.digest = digest
+	if !a.Options.Headless || a.Store != c.Credentials || c.Backend != "headless" ||
+		filepath.Base(in.binary) != "codex" || c.Executables[in.binary] != in.digest ||
+		!slices.Contains(c.PresetNames, a.Options.Preset) || operatorFileFact(a, "served_address") != c.Address {
+		return in, store.ErrOperatorConflict
+	}
+	now := e.terminal.now()
+	if now >= row.ExpiresAt || now >= files.Session.ExpiresAt {
+		return in, store.ErrOperatorExpired
+	}
+	if err := e.checkStartPair(c, root); err != nil {
+		return in, err
+	}
+	h, err := e.health(c.Address)
+	build, _ := json.Marshal(h.Build)
+	if err != nil || string(build) != operatorFileFact(a, "served_build") || h.ToolRef != operatorFileFact(a, "served_interpreter") {
+		return in, store.ErrOperatorConflict
+	}
+	project := operatorFileFact(a, "mint_project")
+	if !slices.ContainsFunc(h.ExecProfiles, func(p daemon.HealthExecProfile) bool {
+		return p.Project == project && p.Digest == operatorFileFact(a, "mint_profile")
+	}) {
+		return in, store.ErrOperatorConflict
+	}
+	commands, err := launchCommands(h, project)
+	if err != nil {
+		return in, err
+	}
+	if _, err := recoveryWorktree(ctx, store.OperatorStartRecovery{Worktree: files.Worktree}, a); err != nil {
+		return in, err
+	}
+	set, start, err := operatorFilePaths(files, a)
+	proofs, proofErr := readStartRecoveryFiles(c, files, a)
+	if err != nil || proofErr != nil || start.Agent != "codex" || proofs.Files[0].Status != "exact" || proofs.Files[1].Status != "exact" {
+		return in, store.ErrOperatorConflict
+	}
+	set.Device, set.Inode = proofs.Dir.ID.Dev, proofs.Dir.ID.Inode
+	d, _ := openOperatorFilesDir(set)
+	if d == nil {
+		return in, store.ErrOperatorConflict
+	}
+	defer d.Close()
+	tokenProof, token := readOperatorFile(d, set.Names[1], set.Hashes[1], set.UID)
+	if tokenProof.Status != "exact" {
+		return in, store.ErrOperatorConflict
+	}
+	in.token = string(token)
+	in.wall, err = time.ParseDuration(operatorFileFact(a, "launch_wall"))
+	if err != nil || in.wall <= 0 || in.wall > time.Duration(1<<63-1)-time.Minute {
+		return in, store.ErrOperatorConflict
+	}
+	in.dir = filepath.Join(filepath.Dir(in.sessions), "launch", row.Episode)
+	if physical, err := operatorPhysicalTarget(in.dir); err != nil || physical != in.dir {
+		return in, store.ErrOperatorConflict
+	}
+	if err := validateLaunchDir(in.dir, root, operatorWorktreeFact(a), in.sessions); err != nil {
+		return in, err
+	}
+	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+c.Address+"/v1/session/self", nil)
+	if err != nil {
+		return in, err
+	}
+	req.Header.Set("Authorization", "Bearer "+in.token)
+	res, err := e.client.Do(req)
+	if err != nil {
+		return in, store.ErrOperatorConflict
+	}
+	defer res.Body.Close()
+	var self daemon.SessionSelf
+	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
+	got, gotErr := filepath.EvalSymlinks(operatorWorktreeFact(a))
+	if err != nil || res.StatusCode != http.StatusOK || json.Unmarshal(body, &self) != nil ||
+		self.SecondsLeft <= 0 || self.ExpiresAt != files.Session.ExpiresAt || self.Episode != row.Episode {
+		return in, store.ErrOperatorConflict
+	}
+	peer, peerErr := filepath.EvalSymlinks(self.Worktree)
+	if gotErr != nil || peerErr != nil || got != peer {
+		return in, store.ErrOperatorConflict
+	}
+	contextText, err := launchContext(body, c.Address, in.token)
+	if err != nil {
+		return in, err
+	}
+	in.prompt = launchPrompt(a.Task, commands)
+	in.argv = launchArgv(start.Agent, start.Model, in.prompt, in.dir, "", c.Address, contextText, "", nil)
+	in.argv[0] = in.binary
+	for _, arg := range in.argv {
+		if strings.Contains(arg, in.token) {
+			return in, store.ErrOperatorConflict
+		}
+	}
+	proof, err := authority.RecheckClaimedOperatorAction(ctx, ro, a, [32]byte(row.Challenge), row.ExpiresAt, e.terminal.now(), row)
+	if err != nil {
+		return in, err
+	}
+	in.proof = proof.OperatorConfirmation().MintProof()
+	return in, nil
+}
+
+type operatorAgentRunner func(context.Context, []string, string, []string, *os.File, *os.File, time.Duration, agentrun.LaunchHooks) (agentrun.Result, error)
+
+func (e serviceEnv) launchPreparedOperatorStart(outer context.Context, c runnerConfig, h store.OperatorHeartbeat,
+	src operatorStartSource, stops <-chan os.Signal, process *os.File, log io.Writer, ro *store.Credentials,
+	row store.OperatorActionRow, current store.OperatorStartRecovery, episode *store.OperatorEpisodeLock) error {
+	ctx, cancel := e.operatorContext(operatorLaunchBudget)
+	defer cancel()
+	files, a, err := ro.ReadOperatorFiles(ctx, row.ID)
+	if err != nil {
+		return err
+	}
+	if err := e.proveOldTerminalQuiet(); err != nil {
+		return err
+	}
+	in, err := e.launchInputs(ctx, c, src.workspace, row, files, a, ro)
+	if err != nil {
+		return err
+	}
+	if err := launchLockCheck(filepath.Join(in.sessions, row.Episode+".launch.lock")); err != nil {
+		return err
+	}
+	if outer.Err() != nil {
+		return outer.Err()
+	}
+	fd, err := episode.Verify(c.Credentials, row.Episode)
+	if err != nil {
+		return err
+	}
+	lifecycle, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
+	if err != nil {
+		return errOperatorDraining
+	}
+	defer lifecycle.Close()
+	select {
+	case <-stops:
+		return errOperatorDraining
+	default:
+	}
+	job := operatorLaunchJob{config: c, row: row, wall: in.wall, locks: [3]*os.File{process, fd, lifecycle}}
+	if _, err := rand.Read(job.secret[:]); err != nil {
+		return err
+	}
+	if err := runnerWrite(c.Credentials, func(db *store.Credentials) error {
+		live, err := db.OperatorHeartbeat(ctx)
+		catalogue, catErr := strconv.ParseInt(operatorFileFact(a, "mint_catalogue"), 10, 64)
+		if err != nil || catErr != nil || !preparationHeartbeat(live, h.Owner, h.Build, catalogue, src.clock()) {
+			return store.ErrOperatorRunner
+		}
+		return db.ReserveOperatorLaunch(ctx, current, job.secret, in.proof, src.clock)
+	}); err != nil {
+		return err
+	}
+	if e.operatorSupervisor != nil {
+		return e.operatorSupervisor(job, log)
+	}
+	if ctx.Err() != nil || serviceFileHash(e.paths.Binary) != c.Executables[e.paths.Binary] {
+		return e.operatorSpawnRefusal(c, row.ID, in.digest, log)
+	}
+	return e.superviseOperatorLaunch(job, stops, log, h, in.digest)
+}
+
+func (e serviceEnv) superviseOperatorLaunch(job operatorLaunchJob, stops <-chan os.Signal, log io.Writer,
+	h store.OperatorHeartbeat, digest string) error {
+	cmd := exec.Command(e.paths.Binary, "operator-launch-worker", job.row.ID)
+	cmd.Stdin, cmd.Stderr = bytes.NewReader(job.secret[:]), log
+	cmd.ExtraFiles = job.locks[:]
+	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
+	if err := cmd.Start(); err != nil {
+		return e.operatorSpawnRefusal(job.config, job.row.ID, digest, log)
+	}
+	job.locks[2].Close() // Worker holds the same description through agent Start.
+	done := make(chan error, 1)
+	go func() { done <- cmd.Wait() }()
+	ticker := time.NewTicker(runnerPoll)
+	defer ticker.Stop()
+	cap := time.NewTimer(job.wall + time.Minute)
+	defer cap.Stop()
+	for {
+		select {
+		case err := <-done:
+			return err
+		case <-ticker.C:
+			_ = e.runnerTick(job.config, &h)
+		case <-stops:
+			_ = cmd.Process.Signal(syscall.SIGTERM)
+			stops = nil
+		case <-cap.C:
+			_ = cmd.Process.Signal(syscall.SIGTERM)
+			return errors.New("launch supervisor did not report completion")
+		}
+	}
+}
+
+func (e serviceEnv) operatorSpawnRefusal(c runnerConfig, id, digest string, log io.Writer) error {
+	r := store.OperatorLaunchReceipt{Run: id, Executable: digest, Status: "spawn-refused"}
+	if err := e.writeOperatorLaunch(c, r, nil); err != nil {
+		return err
+	}
+	fmt.Fprintf(log, "operator launch %s: spawn-refused\n", id)
+	return nil
+}
+
+func (e serviceEnv) writeOperatorLaunch(c runnerConfig, r store.OperatorLaunchReceipt, history *store.EpisodeLaunch) error {
+	ctx, cancel := e.operatorContext(operatorLaunchBudget)
+	defer cancel()
+	return runnerWrite(c.Credentials, func(db *store.Credentials) error {
+		read, err := db.ReadOperatorLaunch(ctx, r.Run)
+		if err != nil {
+			return err
+		}
+		return db.RecordOperatorLaunch(ctx, read, r, history, e.terminal.now())
+	})
+}
+
+func (e serviceEnv) runOperatorLaunchWorker(args []string, out io.Writer) int {
+	proc, stop := serveSignalContext()
+	defer stop()
+	if len(args) != 1 || !verbs.ValidOperatorActionID(args[0]) {
+		return exitUsage
+	}
+	c, b, err := e.installedRunner()
+	if err != nil || os.Getenv("WORLD_OPERATOR_CONFIG_SHA256") != serviceHash(b) ||
+		daemon.CurrentBuild().ExecutableSHA256 != c.Executables[e.paths.Binary] {
+		return exitUsage
+	}
+	ctx, cancel := e.operatorContext(operatorLaunchBudget)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(c.Credentials)
+	if err != nil {
+		return exitUsage
+	}
+	defer ro.Close()
+	row, err := ro.OperatorAction(ctx, args[0], 0)
+	if err != nil || !verbs.EpisodeIDPattern.MatchString(row.Episode) {
+		return exitUsage
+	}
+	job := operatorLaunchJob{config: c, row: row}
+	names := [3]string{c.Credentials + ".operator-runner.lock",
+		store.OperatorEpisodeLockPath(c.Credentials, row.Episode), filepath.Dir(e.paths.Record)}
+	for i, name := range names {
+		fd := os.NewFile(uintptr(3+i), name)
+		defer fd.Close()
+		st, err := fd.Stat()
+		path, pathErr := os.Lstat(name)
+		if err != nil || pathErr != nil || !os.SameFile(st, path) || syscall.Flock(int(fd.Fd()), syscall.LOCK_EX|syscall.LOCK_NB) != nil {
+			return exitUsage
+		}
+		syscall.CloseOnExec(int(fd.Fd()))
+		job.locks[i] = fd
+	}
+	if _, err := io.ReadFull(os.Stdin, job.secret[:]); err != nil {
+		return exitUsage
+	}
+	if err := e.executeOperatorLaunch(proc, job, out); err != nil {
+		fmt.Fprintf(out, "operator launch %s: %s\n", row.ID, operatorLaunchFailure(err))
+		return 1
+	}
+	return exitOK
+}
+
+func operatorLaunchFailure(err error) string {
+	var verdict *approval.NotAdmittedError
+	if errors.Is(err, store.ErrOperatorExpired) || errors.As(err, &verdict) && verdict.Verdict == approval.RefuseExpired {
+		return "expired"
+	}
+	return "blocked or uncertain; inspect service runner launch-status"
+}
+
+func operatorLaunchProcessStatus(s store.OperatorLaunchSnapshot, read func(int) (string, bool)) string {
+	if s.NoSpawn {
+		return "spawn-refused"
+	}
+	if s.Phase == "reserved" {
+		return "uncertain"
+	}
+	if s.Receipt.Status == "launched" {
+		if birth, live := read(s.Receipt.PID); !live || birth != s.Receipt.Birth {
+			return "uncertain"
+		}
+	}
+	return s.Receipt.Status
+}
+
+func (e serviceEnv) operatorLaunchStatus(id string, out io.Writer) error {
+	c, _, err := e.installedRunner()
+	if err != nil {
+		return err
+	}
+	db, err := store.OpenCredentialsReadOnly(c.Credentials)
+	if err != nil {
+		return err
+	}
+	defer db.Close()
+	ctx, cancel := e.operatorContext(operatorLaunchBudget)
+	defer cancel()
+	s, err := db.ReadOperatorLaunch(ctx, id)
+	if err != nil {
+		return err
+	}
+	fmt.Fprintf(out, "operator launch %s: %s\n", id, operatorLaunchProcessStatus(s, operatorProcessIdentity))
+	return nil
+}
+
+func launchOperatorEnvironment(c runnerConfig, token string) []string {
+	path := strings.Join(append(slices.Clone(c.ToolDirs), runnerSystemPath), ":")
+	return []string{"HOME=" + c.Home, "PATH=" + path, "GIT_TERMINAL_PROMPT=0", "GH_PROMPT_DISABLED=1", "WORLD_SESSION=" + token}
+}
+
+func operatorLaunchStage(dir string) (*os.File, error) {
+	f, err := os.CreateTemp(dir, ".operator-output-*")
+	if err != nil {
+		return nil, err
+	}
+	if err := os.Remove(f.Name()); err != nil {
+		f.Close()
+		return nil, err
+	}
+	return f, nil
+}
+
+func (e serviceEnv) executeOperatorLaunch(proc context.Context, job operatorLaunchJob, log io.Writer) error {
+	ctx, cancel := e.operatorContext(operatorLaunchBudget)
+	defer cancel()
+	c, row := job.config, job.row
+	record, err := readServiceRecord(e.paths)
+	if err != nil {
+		return err
+	}
+	ro, err := store.OpenCredentialsReadOnly(c.Credentials)
+	if err != nil {
+		return err
+	}
+	defer ro.Close()
+	files, a, err := ro.ReadOperatorSpawnInput(ctx, row.ID)
+	if err != nil {
+		return err
+	}
+	in, err := e.launchInputs(ctx, c, record.Workspace, row, files, a, ro)
+	if err != nil {
+		return err
+	}
+	lock := filepath.Join(in.sessions, row.Episode+".launch.lock")
+	if err := launchLockCheck(lock); err != nil {
+		return err
+	}
+	f, err := os.OpenFile(lock, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
+	if err != nil {
+		return err
+	}
+	identity, statErr := f.Stat()
+	defer func() {
+		if st, err := os.Lstat(lock); err == nil && statErr == nil && os.SameFile(identity, st) {
+			_ = os.Remove(lock)
+		}
+	}()
+	_, err = fmt.Fprintln(f, os.Getpid())
+	if err := errors.Join(err, f.Close()); err != nil {
+		return err
+	}
+	for _, dir := range []string{filepath.Dir(in.dir), in.dir} {
+		if err := os.Mkdir(dir, 0700); err == nil {
+			defer os.Remove(dir) // Remove an empty directory created by this run.
+		} else if !errors.Is(err, os.ErrExist) {
+			return err
+		}
+	}
+	if err := runnerPath(c.Home, in.dir, c.UID, true, 0700); err != nil {
+		return err
+	}
+	stdout, err := operatorLaunchStage(in.sessions)
+	if err != nil {
+		return err
+	}
+	defer stdout.Close()
+	stderr, err := operatorLaunchStage(in.sessions)
+	if err != nil {
+		return err
+	}
+	defer stderr.Close()
+	r := store.OperatorLaunchReceipt{Run: row.ID, Executable: in.digest, Status: "launched"}
+	var writer *store.Credentials
+	granted := false
+	started := time.Time{}
+	hooks := agentrun.LaunchHooks{Files: job.locks[:2],
+		Before: func() error {
+			var err error
+			writer, err = store.OpenCredentials(runnerDB(c.Credentials))
+			if err != nil {
+				return err
+			}
+			if err := proc.Err(); err != nil {
+				return err
+			}
+			if serviceFileHash(in.binary) != in.digest {
+				return &approval.BindingError{Field: "launch_binary_sha256"}
+			}
+			if err := validateLaunchDir(in.dir, record.Workspace, operatorWorktreeFact(a), in.sessions); err != nil {
+				return err
+			}
+			err = writer.AuthorizeOperatorSpawn(ctx, row.ID, job.secret, in.proof, e.terminal.now)
+			granted = err == nil
+			return err
+		},
+		Started: func(pid int) error {
+			defer func() { writer.Close(); writer = nil; job.locks[2].Close() }()
+			r.PID = pid
+			r.Birth, _ = operatorProcessIdentity(pid)
+			if r.Birth == "" {
+				return errors.New("process birth unavailable")
+			}
+			s, err := writer.ReadOperatorLaunch(ctx, row.ID)
+			if err != nil {
+				return err
+			}
+			started = time.Now()
+			if err := writer.RecordOperatorLaunch(ctx, s, r, nil, e.terminal.now()); err != nil {
+				return err
+			}
+			fmt.Fprintf(log, "operator launch %s: launched\n", row.ID)
+			return nil
+		},
+	}
+	run := e.operatorAgentRun
+	if run == nil {
+		run = agentrun.RunObserved
+	}
+	result, runErr := run(proc, in.argv, in.dir, launchOperatorEnvironment(c, in.token), stdout, stderr, in.wall, hooks)
+	if writer != nil {
+		writer.Close()
+	}
+	if r.PID == 0 {
+		if granted {
+			return e.operatorSpawnRefusal(c, row.ID, in.digest, log)
+		}
+		return runErr
+	}
+	if runErr != nil {
+		return runErr // the active reservation or observed receipt stays uncertain
+	}
+	r.Exit = result.ExitCode
+	r.Status = operatorLaunchExitStatus(proc, result)
+	history := operatorLaunchHistory(in, r, row.Episode, stdout, stderr, started, time.Now())
+	if history.FinishReason == "output-failed" {
+		r.Status = "failed"
+	}
+	if err := e.writeOperatorLaunch(c, r, &history); err != nil {
+		return err
+	}
+	fmt.Fprintf(log, "operator launch %s: %s\n", row.ID, r.Status)
+	return nil
+}
+
+func operatorLaunchHistory(in operatorLaunchInputs, r store.OperatorLaunchReceipt, ep string, stdout, stderr *os.File, started, ended time.Time) store.EpisodeLaunch {
+	transcript := filepath.Join(in.sessions, ep+".transcript.jsonl")
+	errPath := filepath.Join(in.sessions, ep+".agent.err")
+	state, digest, outErr := keepLaunchStream(stdout, transcript, in.token, 1)
+	_, _, errErr := keepLaunchStream(stderr, errPath, in.token, 1)
+	finish := r.Status
+	if outErr != nil || errErr != nil {
+		finish = "output-failed"
+	}
+	return store.EpisodeLaunch{Mode: "headless", StartedAt: started.UTC().Format(time.RFC3339),
+		EndedAt: ended.UTC().Format(time.RFC3339), Exit: r.Exit, FinishReason: finish,
+		Host: in.binary, ArgvRedacted: redactLaunchArgv(in.argv, in.prompt, in.token),
+		Transcript: transcript, TranscriptSHA256: digest, TranscriptState: state}
+}
+
+func operatorLaunchExitStatus(proc context.Context, result agentrun.Result) string {
+	if result.Reason == "timeout" && proc.Err() == nil {
+		return "capped"
+	}
+	if result.ExitCode != 0 && result.Reason != "signal" && proc.Err() == nil {
+		return "failed"
+	}
+	return "stopped"
+}
diff --git a/cmd/ailang-worldd/operator_launch_test.go b/cmd/ailang-worldd/operator_launch_test.go
new file mode 100644
index 00000000..55f72cb5
--- /dev/null
+++ b/cmd/ailang-worldd/operator_launch_test.go
@@ -0,0 +1,327 @@
+package main
+
+import (
+	"bytes"
+	"context"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"net/http"
+	"os"
+	"path/filepath"
+	"reflect"
+	"slices"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/agentrun"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+const fakeLaunchMarker = "#!/bin/sh\nprintf x >> \"${0%/*}/ran\"\n"
+
+type operatorLaunchRig struct {
+	*worktreeStartRig
+	c      runnerConfig
+	h      store.OperatorHeartbeat
+	runner *os.File
+	ctx    context.Context
+	db     *store.Credentials
+	self   daemon.SessionSelf
+	now    int64
+	token  string
+}
+
+func newOperatorLaunchRig(t *testing.T, script string) *operatorLaunchRig {
+	t.Helper()
+	r, c, h, runner := operatorActivationFixtureWithAgent(t, script)
+	f := &operatorLaunchRig{worktreeStartRig: r, c: c, h: h, runner: runner, ctx: activationTestContext(t), now: r.src.clock()}
+	r.src.clock = func() int64 { return f.now }
+	r.env.terminal.now = r.src.clock
+	var err error
+	f.db, err = store.OpenCredentialsReadOnly(c.Credentials)
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { f.db.Close() })
+	health := r.env.client.Transport
+	r.env.client.Transport = serviceTransport(func(req *http.Request) (*http.Response, error) {
+		if req.URL.Path != "/v1/session/self" {
+			return health.RoundTrip(req)
+		}
+		if req.URL.Host != c.Address || f.token == "" || req.Header.Get("Authorization") != "Bearer "+f.token {
+			t.Error("launch self did not use the pinned session")
+			return nil, errors.New("fixture incorrectly bound")
+		}
+		body, _ := json.Marshal(f.self)
+		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(body))}, nil
+	})
+	r.env.operatorSupervisor = func(job operatorLaunchJob, log io.Writer) error {
+		s := f.receipt(t)
+		if s.Phase != "reserved" || s.Key.Version != 1 || s.Receipt.PID != 0 {
+			t.Fatalf("supervisor crossed reservation boundary: %+v", s)
+		}
+		return r.env.executeOperatorLaunch(f.ctx, job, log)
+	}
+	if out := f.tick(t); out != "" {
+		t.Fatalf("preparation tick: %s", out)
+	}
+	f.readPrepared(t)
+	return f
+}
+
+func (f *operatorLaunchRig) readPrepared(t *testing.T) {
+	t.Helper()
+	s, err := f.db.ReadOperatorStartRecovery(f.ctx, f.action.ID)
+	if err != nil || s.Window != store.OperatorStartPrepared {
+		t.Fatalf("not prepared: %+v %v", s, err)
+	}
+	files, _, err := f.db.ReadOperatorFiles(f.ctx, f.action.ID)
+	if err != nil {
+		t.Fatal(err)
+	}
+	token, err := os.ReadFile(files.Mint.TokenPath)
+	if err != nil || len(token) != 64 {
+		t.Fatalf("real mint did not write a token: %v", err)
+	}
+	f.token = string(token)
+	f.self = daemon.SessionSelf{Episode: f.action.Episode, Worktree: f.path(), ExpiresAt: files.Session.ExpiresAt, SecondsLeft: files.Session.ExpiresAt - f.now}
+}
+
+func (f *operatorLaunchRig) tick(t *testing.T) string {
+	t.Helper()
+	ticks := make(chan time.Time)
+	close(ticks)
+	var out syncBuffer
+	if rc := f.env.runOperatorPreparationLoop(&out, f.c, f.h, f.src, nil, ticks, f.runner); rc != exitOK {
+		t.Fatalf("runner tick %d: %s", rc, out.String())
+	}
+	return out.String()
+}
+
+func (f *operatorLaunchRig) receipt(t *testing.T) store.OperatorLaunchSnapshot {
+	t.Helper()
+	s, err := f.db.ReadOperatorLaunch(f.ctx, f.action.ID)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return s
+}
+
+func (f *operatorLaunchRig) assertUnspawned(t *testing.T) {
+	t.Helper()
+	if _, err := os.Stat(filepath.Join(f.bin, "ran")); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("fake child ran: %v", err)
+	}
+	if _, err := f.db.ReadOperatorLaunch(f.ctx, f.action.ID); !errors.Is(err, store.ErrOperatorMissing) {
+		t.Fatalf("unexpected reservation or receipt: %v", err)
+	}
+}
+
+func TestOperatorSpawnRefused(t *testing.T) {
+	f := newOperatorLaunchRig(t, "#!/no-world-fixture-interpreter\nprintf unexpected\n")
+	before := activationEvidence(t, f.worktreeStartRig)
+	out := f.tick(t)
+	s := f.receipt(t)
+	after := activationEvidence(t, f.worktreeStartRig)
+	if !strings.Contains(out, ": spawn-refused") || !s.NoSpawn || s.Receipt != (store.OperatorLaunchReceipt{}) ||
+		!reflect.DeepEqual(before[1:], after[1:]) {
+		t.Fatalf("known exec refusal: %s %+v", out, s)
+	}
+	if _, err := os.Lstat(filepath.Join(f.bin, "ran")); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("exec refusal spawned a fixture: %v", err)
+	}
+	start, err := f.db.EpisodeStart(f.ctx, f.action.Episode)
+	if err != nil || len(start.Launches) != 0 {
+		t.Fatalf("exec refusal created launch history: %v %v", start, err)
+	}
+}
+
+func TestOperatorLaunchCrashWindow(t *testing.T) {
+	s := store.OperatorLaunchSnapshot{Phase: "observed", Receipt: store.OperatorLaunchReceipt{Run: fixtureActionID, PID: 42, Birth: "first", Status: "launched"}}
+	for _, birth := range []string{"reused", "first"} {
+		want := "uncertain"
+		if birth == "first" {
+			want = "launched"
+		}
+		if got := operatorLaunchProcessStatus(s, func(pid int) (string, bool) { return birth, pid == 42 }); got != want {
+			t.Fatalf("birth %q classified %q, want %q", birth, got, want)
+		}
+	}
+}
+
+func TestOperatorLaunchExit(t *testing.T) {
+	t.Run("failed", func(t *testing.T) {
+		dir := t.TempDir()
+		script := filepath.Join(dir, "child")
+		if err := os.WriteFile(script, []byte("#!/bin/sh\nexit 7\n"), 0700); err != nil {
+			t.Fatal(err)
+		}
+		stream, err := os.CreateTemp(dir, "stream")
+		if err != nil {
+			t.Fatal(err)
+		}
+		ctx := activationTestContext(t)
+		result, err := agentrun.RunObserved(ctx, []string{script}, dir, []string{}, stream, stream, time.Minute, agentrun.LaunchHooks{})
+		_ = stream.Close()
+		history := operatorLaunchHistory(operatorLaunchInputs{sessions: dir}, store.OperatorLaunchReceipt{Status: "failed"}, "ep", stream, stream, time.Time{}, time.Time{})
+		if err != nil || operatorLaunchExitStatus(ctx, result) != "failed" || result.ExitCode != 7 || history.FinishReason != "output-failed" {
+			t.Fatalf("nonzero real child: %+v %v", result, err)
+		}
+	})
+	t.Run("capped_and_scrubbed", func(t *testing.T) {
+		dir := t.TempDir()
+		script := filepath.Join(dir, "child")
+		if err := os.WriteFile(script, []byte("#!/bin/sh\nprintf 'out %s\\n' \"$WORLD_SESSION\"\nprintf 'err %s\\n' \"$WORLD_SESSION\" >&2\nprintf x >&5\nwhile :; do :; done\n"), 0700); err != nil {
+			t.Fatal(err)
+		}
+		stdout, err := operatorLaunchStage(dir)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer stdout.Close()
+		stderr, err := operatorLaunchStage(dir)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer stderr.Close()
+		r, w, err := os.Pipe()
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer r.Close()
+		defer w.Close()
+		proc := activationTestContext(t)
+		deadline, _ := proc.Deadline()
+		if err := r.SetReadDeadline(deadline); err != nil {
+			t.Fatal(err)
+		}
+		ctx, cap := context.WithCancel(proc)
+		defer cap()
+		canary := strings.Repeat("1a", 32)
+		hooks := agentrun.LaunchHooks{Files: []*os.File{nil, nil, w}, Started: func(int) error {
+			var ready [1]byte
+			_, err := io.ReadFull(r, ready[:])
+			cap()
+			return err
+		}}
+		argv := []string{script, "sensitive task"}
+		result, err := agentrun.RunObserved(ctx, argv, dir, []string{"WORLD_SESSION=" + canary}, stdout, stderr, time.Minute, hooks)
+		status := operatorLaunchExitStatus(proc, result)
+		if err != nil || status != "capped" {
+			t.Fatalf("capped launch: %+v %v", result, err)
+		}
+		if raw, err := filepath.Glob(filepath.Join(dir, ".operator-output-*")); err != nil || len(raw) != 0 {
+			t.Fatalf("unscrubbed stream retained a pathname: %v (%v)", raw, err)
+		}
+		in := operatorLaunchInputs{sessions: dir, token: canary, argv: argv, prompt: argv[1]}
+		history := operatorLaunchHistory(in, store.OperatorLaunchReceipt{Status: status}, "ep", stdout, stderr, time.Time{}, time.Time{})
+		for _, path := range []string{"ep.transcript.jsonl", "ep.agent.err"} {
+			body, err := os.ReadFile(filepath.Join(dir, path))
+			if err != nil || strings.Contains(string(body), canary) || !strings.Contains(string(body), "scrubbed") {
+				t.Fatalf("child canary reached %s or scrub control absent: %v", path, err)
+			}
+		}
+		raw, err := json.Marshal(history)
+		if err != nil || strings.Contains(string(raw), canary) || strings.Contains(string(raw), argv[1]) ||
+			history.TranscriptState != "scrubbed" || history.FinishReason != "capped" {
+			t.Fatalf("unredacted launch history: %v", err)
+		}
+	})
+}
+
+func testOperatorStartOnceLaunch(t *testing.T) {
+	t.Setenv("WORLD_AMBIENT_TEST_CANARY", "must stay out")
+	f := newOperatorLaunchRig(t, fakeLaunchMarker+"[ -e /dev/fd/3 ] && [ -e /dev/fd/4 ]\n")
+	f.assertUnspawned(t)
+	supervise := f.env.operatorSupervisor
+	f.env.operatorSupervisor = func(job operatorLaunchJob, log io.Writer) error {
+		before := f.receipt(t)
+		var nested syncBuffer
+		ticks := make(chan time.Time, 1)
+		ticks <- time.Time{}
+		close(ticks)
+		restarted := f.env
+		restarted.operatorDeadline = nil
+		if rc := restarted.runOperatorPreparationLoop(&nested, f.c, f.h, f.src, nil, ticks, f.runner); rc != exitOK ||
+			strings.Count(nested.String(), "launch-status "+f.action.ID) != 1 || f.receipt(t) != before {
+			t.Fatalf("reserved restart crossed spawn boundary: %d %s", rc, &nested)
+		}
+		return supervise(job, log)
+	}
+	var budgets []time.Duration
+	f.env.operatorDeadline = func(want time.Duration, ctx context.Context) {
+		budgets = append(budgets, want)
+		deadline, ok := ctx.Deadline()
+		if !ok || time.Until(deadline) > activationTestBudget {
+			t.Error("launch missed shared bounded context override")
+		}
+	}
+	f.env.operatorAgentRun = func(ctx context.Context, argv []string, dir string, env []string, out, stderr *os.File, wall time.Duration, hooks agentrun.LaunchHooks) (agentrun.Result, error) {
+		commands, err := launchCommands(f.health, f.in.Profile.Project)
+		body, _ := json.Marshal(f.self)
+		text, _ := launchContext(body, f.c.Address, f.token)
+		want := launchArgv("codex", f.in.Start.Model, launchPrompt(f.action.Task, commands), dir, "", f.c.Address, text, "", nil)
+		want[0] = filepath.Join(f.bin, "codex")
+		if err != nil || !slices.Equal(want, argv) || !slices.Contains(argv, "read-only") ||
+			!slices.Contains(env, "WORLD_SESSION="+f.token) || !slices.Contains(env, "HOME="+f.c.Home) ||
+			slices.ContainsFunc(env, func(v string) bool { return strings.Contains(v, "WORLD_AMBIENT_TEST_CANARY") }) ||
+			wall != time.Hour || len(hooks.Files) != 2 {
+			t.Fatal("managed launch diverged from reviewed V3 confinement or process bounds")
+		}
+		return agentrun.RunObserved(ctx, argv, dir, env, out, stderr, wall, hooks)
+	}
+	out := f.tick(t)
+	if !strings.Contains(out, ": stopped") {
+		t.Fatalf("installed tick did not complete launch: %s", out)
+	}
+	if !reflect.DeepEqual(budgets, []time.Duration{time.Second, 22 * time.Second, 5 * time.Second, 5 * time.Second, 5 * time.Second}) {
+		t.Fatalf("launch context budget calls: %v", budgets)
+	}
+	before := f.receipt(t)
+	f.tick(t)
+	ran, err := os.ReadFile(filepath.Join(f.bin, "ran"))
+	if err != nil || string(ran) != "x" || f.receipt(t) != before || before.Receipt.Run != f.action.ID ||
+		before.Receipt.Executable != operatorFileFact(f.action, "launch_binary_sha256") {
+		t.Fatalf("spawn count or identity differed: %q %v", ran, err)
+	}
+}
+
+func TestOperatorStartBrowser(t *testing.T) {
+	f := newOperatorLaunchRig(t, fakeLaunchMarker+"while [ ! -e release ]; do :; done\nexit 0\n")
+	f.env.operatorAgentRun = func(ctx context.Context, argv []string, dir string, env []string, out, stderr *os.File, wall time.Duration, hooks agentrun.LaunchHooks) (agentrun.Result, error) {
+		start := hooks.Started
+		hooks.Started = func(pid int) error {
+			defer func() { _ = os.WriteFile(filepath.Join(dir, "release"), nil, 0600) }()
+			if err := start(pid); err != nil {
+				return err
+			}
+			f.now = f.self.ExpiresAt + 1
+			f.env.client.Transport = serviceTransport(func(*http.Request) (*http.Response, error) {
+				return nil, errors.New("fixture daemon restarting")
+			})
+			var status, stderr syncBuffer
+			if rc := runServiceRunner([]string{"launch-status", f.action.ID}, &status, &stderr, f.env); rc != exitOK ||
+				status.String() != fmt.Sprintf("operator launch %s: launched\n", f.action.ID) {
+				t.Errorf("restart failed to observe live receipt: %s %s", &status, &stderr)
+			}
+			return nil
+		}
+		return agentrun.RunObserved(ctx, argv, dir, env, out, stderr, wall, hooks)
+	}
+	before, _, err := f.db.ReadOperatorFiles(f.ctx, f.action.ID)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if out := f.tick(t); !strings.Contains(out, ": launched") || !strings.Contains(out, ": stopped") {
+		t.Fatalf("launch could not finish after expiry: %s", out)
+	}
+	s, ok, err := f.db.ResolveSession(f.ctx, before.Mint.CredentialID)
+	h, histErr := f.db.EpisodeStart(f.ctx, f.action.Episode)
+	if err != nil || histErr != nil || !ok || s != before.Session || len(h.Launches) != 1 || f.receipt(t).Receipt.Run != f.action.ID {
+		t.Fatalf("restart changed start history or duplicated mint: %v %v %t", err, histErr, ok)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_preparation.go b/cmd/ailang-worldd/operator_preparation.go
index 52f701fe..d0bb7ef9 100644
--- a/cmd/ailang-worldd/operator_preparation.go
+++ b/cmd/ailang-worldd/operator_preparation.go
@@ -20,13 +20,13 @@ var errOperatorDraining = errors.New("operator runner draining")
 
 func (e serviceEnv) consumeOperatorPreparation(ctx context.Context, c runnerConfig, heartbeat store.OperatorHeartbeat,
 	src operatorStartSource, stops <-chan os.Signal, process *os.File, recoveryCursor *string,
-	acquire func(string, string) (*store.OperatorEpisodeLock, error)) error {
+	acquire func(string, string) (*store.OperatorEpisodeLock, error), launch ...operatorPreparedLaunch) error {
 	select {
 	case <-stops:
 		return errOperatorDraining
 	default:
 	}
-	if err := e.recoverOperatorStart(ctx, c, heartbeat, src, recoveryCursor, stops); err != nil {
+	if err := e.recoverOperatorStart(ctx, c, heartbeat, src, recoveryCursor, stops, launch...); err != nil {
 		return err
 	}
 	return e.prepareOperatorStart(ctx, c, heartbeat.Owner, heartbeat.Build, src, stops, process, acquire)
diff --git a/cmd/ailang-worldd/operator_process_darwin.go b/cmd/ailang-worldd/operator_process_darwin.go
index 42647641..d5121a31 100644
--- a/cmd/ailang-worldd/operator_process_darwin.go
+++ b/cmd/ailang-worldd/operator_process_darwin.go
@@ -1,6 +1,7 @@
 package main
 
 import (
+	"fmt"
 	"os"
 
 	"golang.org/x/sys/unix"
@@ -18,3 +19,14 @@ func listWorldProcesses() ([]worldProcess, error) {
 	}
 	return out, nil
 }
+
+// operatorProcessIdentity reads kernel birth time and excludes a reused PID.
+func operatorProcessIdentity(pid int) (string, bool) {
+	row, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
+	if err != nil || int(row.Proc.P_pid) != pid {
+		return "", false
+	}
+	p := row.Proc
+	const zombie = 5 // SZOMB in Darwin <sys/proc.h>
+	return fmt.Sprintf("%d.%d", p.P_starttime.Sec, p.P_starttime.Usec), p.P_stat != zombie
+}
diff --git a/cmd/ailang-worldd/operator_process_linux.go b/cmd/ailang-worldd/operator_process_linux.go
new file mode 100644
index 00000000..4fe77293
--- /dev/null
+++ b/cmd/ailang-worldd/operator_process_linux.go
@@ -0,0 +1,26 @@
+package main
+
+import (
+	"errors"
+	"os"
+	"strconv"
+	"strings"
+)
+
+func listWorldProcesses() ([]worldProcess, error) {
+	return nil, errors.New("operator upgrade requires the macOS process inventory")
+}
+
+func operatorProcessIdentity(pid int) (string, bool) {
+	stat, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
+	boot, bootErr := os.ReadFile("/proc/sys/kernel/random/boot_id")
+	end := strings.LastIndexByte(string(stat), ')')
+	if err != nil || bootErr != nil || end < 0 {
+		return "", false
+	}
+	fields := strings.Fields(string(stat[end+1:]))
+	if len(fields) <= 19 {
+		return "", false
+	}
+	return strings.TrimSpace(string(boot)) + ":" + fields[19], fields[0] != "Z"
+}
diff --git a/cmd/ailang-worldd/operator_process_other.go b/cmd/ailang-worldd/operator_process_other.go
index 7a30ebd9..ad91f103 100644
--- a/cmd/ailang-worldd/operator_process_other.go
+++ b/cmd/ailang-worldd/operator_process_other.go
@@ -1,4 +1,4 @@
-//go:build !darwin
+//go:build !darwin && !linux
 
 package main
 
@@ -7,3 +7,5 @@ import "errors"
 func listWorldProcesses() ([]worldProcess, error) {
 	return nil, errors.New("operator upgrade requires the macOS process inventory")
 }
+
+func operatorProcessIdentity(int) (string, bool) { return "", false }
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
index 4d1e55f4..efacfff4 100644
--- a/cmd/ailang-worldd/operator_runner.go
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -114,12 +114,22 @@ func (e serviceEnv) runOperatorPreparationLoop(errOut io.Writer, c runnerConfig,
 		}
 		ctx, cancel := e.operatorContext(operatorPreparationBudget)
 		src.workspace = record.Workspace
-		err = e.consumeOperatorPreparation(ctx, c, h, src, stops, lock, &recoveryCursor, store.OpenOperatorEpisodeLock)
+		launch := func(ro *store.Credentials, row store.OperatorActionRow, current store.OperatorStartRecovery, episode *store.OperatorEpisodeLock) error {
+			err := e.launchPreparedOperatorStart(ctx, c, h, src, stops, lock, errOut, ro, row, current, episode)
+			if errors.Is(err, errOperatorDraining) {
+				return err
+			}
+			if err != nil {
+				fmt.Fprintf(errOut, "operator launch %s: %s\n", row.ID, operatorLaunchFailure(err))
+			}
+			return errOperatorLaunchHandled
+		}
+		err = e.consumeOperatorPreparation(ctx, c, h, src, stops, lock, &recoveryCursor, store.OpenOperatorEpisodeLock, launch)
 		cancel()
 		switch {
 		case errors.Is(err, errOperatorDraining):
 			return 0
-		case err != nil && !errors.Is(err, store.ErrOperatorMissing):
+		case err != nil && !errors.Is(err, store.ErrOperatorMissing) && !errors.Is(err, errOperatorLaunchHandled):
 			var stale *approval.BindingError
 			var remedy operatorRecoveryNeeded
 			switch {
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index ebc2aa04..d915f23c 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -21,7 +21,7 @@ import (
 
 const runnerLabel = "com.sunholo.ailang-world.operator-runner"
 const runnerSystemPath = "/usr/bin:/bin:/usr/sbin:/sbin"
-const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | reconcile-files ACTION-ID | reconcile-start ACTION-ID | uninstall
+const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | reconcile-files ACTION-ID | reconcile-start ACTION-ID | launch-status ACTION-ID | uninstall
 Install, stop, cancel, orphan mint revocation, worktree and file disposition and uninstall require typed y. Start uses the approved configuration.
 Config JSON keys: home, uid, address, world, credentials, executables (absolute path to sha256 digest), preset_names, tool_dirs, backend ("headless").
 The World service must already be running. Use absolute, canonical paths and include its installed binary in executables.
@@ -272,7 +272,7 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 		if len(args) != 3 || args[1] != "--config" {
 			return refuse(errors.New("expected install --config FILE"))
 		}
-	case "cancel", "reconcile-mint", "reconcile-worktree", "reconcile-files", "reconcile-start":
+	case "cancel", "reconcile-mint", "reconcile-worktree", "reconcile-files", "reconcile-start", "launch-status":
 		if len(args) != 2 || !verbs.ValidOperatorActionID(args[1]) {
 			return refuse(fmt.Errorf("expected %s ACTION-ID", verb))
 		}
@@ -295,6 +295,12 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 	if verb == "reconcile-files" {
 		return e.reconcileRunnerFiles(args[1], out, errOut)
 	}
+	if verb == "launch-status" {
+		if err := e.operatorLaunchStatus(args[1], out); err != nil {
+			return refuse(err)
+		}
+		return exitOK
+	}
 	p := e.runnerPaths()
 	if verb == "status" {
 		if err := e.runnerStatus(out); err != nil {
diff --git a/cmd/ailang-worldd/operator_start_recovery.go b/cmd/ailang-worldd/operator_start_recovery.go
index f9efc8dd..00063524 100644
--- a/cmd/ailang-worldd/operator_start_recovery.go
+++ b/cmd/ailang-worldd/operator_start_recovery.go
@@ -21,7 +21,9 @@ type operatorRecoveryNeeded struct{ command string }
 
 func (r operatorRecoveryNeeded) Error() string { return r.command }
 
-func (e serviceEnv) recoverOperatorStart(ctx context.Context, c runnerConfig, h store.OperatorHeartbeat, src operatorStartSource, cursor *string, stops <-chan os.Signal) error {
+type operatorPreparedLaunch func(*store.Credentials, store.OperatorActionRow, store.OperatorStartRecovery, *store.OperatorEpisodeLock) error
+
+func (e serviceEnv) recoverOperatorStart(ctx context.Context, c runnerConfig, h store.OperatorHeartbeat, src operatorStartSource, cursor *string, stops <-chan os.Signal, launch ...operatorPreparedLaunch) error {
 	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
 	if err != nil {
 		return err
@@ -36,14 +38,20 @@ func (e serviceEnv) recoverOperatorStart(ctx context.Context, c runnerConfig, h
 		return err
 	}
 	*cursor = row.ID
-	return e.recoverOwnedStart(ctx, ro, c, h, src, row, stops)
+	return e.recoverOwnedStart(ctx, ro, c, h, src, row, stops, launch)
 }
 
 func (e serviceEnv) recoverOwnedStart(ctx context.Context, ro *store.Credentials, c runnerConfig, h store.OperatorHeartbeat,
-	src operatorStartSource, row store.OperatorActionRow, stops <-chan os.Signal) error {
+	src operatorStartSource, row store.OperatorActionRow, stops <-chan os.Signal, launch []operatorPreparedLaunch) error {
 	if !verbs.EpisodeIDPattern.MatchString(row.Episode) {
 		return store.ErrOperatorConflict
 	}
+	if _, err := ro.ReadOperatorLaunch(ctx, row.ID); !errors.Is(err, store.ErrOperatorMissing) {
+		if err != nil {
+			return err
+		}
+		return operatorRecoveryNeeded{"ailang-worldd service runner launch-status " + row.ID}
+	}
 	lock, err := store.OpenOperatorEpisodeLock(c.Credentials, row.Episode)
 	if err != nil {
 		return fmt.Errorf("%w: live episode owner or effect", store.ErrOperatorConflict)
@@ -69,6 +77,9 @@ func (e serviceEnv) recoverOwnedStart(ctx context.Context, ro *store.Credentials
 	case store.OperatorStartTreeOnly, store.OperatorStartMintPending, store.OperatorStartMintUncertain:
 		return operatorRecoveryNeeded{reconcileStartCommand(row.ID)}
 	case store.OperatorStartPrepared:
+		if len(launch) != 0 {
+			return launch[0](ro, row, current, lock)
+		}
 		return verifyPreparedStart(ctx, c, current, ro, a, src.clock())
 	case store.OperatorStartUntouched:
 		if err := e.proveOldTerminalQuiet(); err != nil {
diff --git a/cmd/ailang-worldd/operator_worktree_test.go b/cmd/ailang-worldd/operator_worktree_test.go
index dc29e81c..c0e3d565 100644
--- a/cmd/ailang-worldd/operator_worktree_test.go
+++ b/cmd/ailang-worldd/operator_worktree_test.go
@@ -15,6 +15,7 @@ import (
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/daemon"
 	"github.com/sunholo-data/ailang-world/host/store"
 	"github.com/sunholo-data/ailang-world/host/worktree"
 )
@@ -34,6 +35,11 @@ func newWorktreeStartRig(t *testing.T, claude bool, ttl ...int64) *worktreeStart
 }
 
 func newWorktreeStartRigWithContext(t *testing.T, claude bool, ctx context.Context, ttl ...int64) *worktreeStartRig {
+	t.Helper()
+	return newWorktreeStartRigWithAgent(t, claude, ctx, "", ttl...)
+}
+
+func newWorktreeStartRigWithAgent(t *testing.T, claude bool, ctx context.Context, body string, ttl ...int64) *worktreeStartRig {
 	t.Helper()
 	s := operatorService(t)
 	s.mustInstall(t)
@@ -46,6 +52,11 @@ func newWorktreeStartRigWithContext(t *testing.T, claude bool, ctx context.Conte
 		t.Fatal(err)
 	}
 	startReadFakeAgentIn(t, f, bin)
+	if body != "" {
+		if err := os.WriteFile(filepath.Join(bin, "codex"), []byte(body), 0700); err != nil {
+			t.Fatal(err)
+		}
+	}
 	f.r.sessions = filepath.Join(s.config.Home, "sessions")
 	r := operatorStartRig{f: f, bin: bin, creds: s.config.Credentials, hosts: t.TempDir()}
 	r.f.r.f.db = s.config.World
@@ -69,7 +80,14 @@ func newWorktreeStartRigWithContext(t *testing.T, claude bool, ctx context.Conte
 	})
 	h := matchingHealth(t, r.f.r.f, r.f.r.profiles, r.f.profile)
 	h.DBPath = s.config.World
+	if body != "" {
+		h.Build.ExecutableSHA256 = s.config.Executables[s.env.paths.Binary]
+		s.health = h
+	}
 	stubDaemon(t, h, nil)
+	if body != "" {
+		episodeSelfBuild = func() daemon.BuildIdentity { return h.Build }
+	}
 	r.at = store.OperatorMintAttempt{Catalogue: 1}
 	src := operatorStartSource{r.f.dir, r.f.r.profiles, r.f.r.f.ws, r.hosts, s.env.terminal.now}
 	if ctx == nil {
@@ -213,6 +231,11 @@ func (r *worktreeStartRig) run(t *testing.T, cmd string) (int, string) {
 }
 
 func TestOperatorStartOnce(t *testing.T) {
+	t.Run("worktree", testOperatorStartOnceWorktree)
+	t.Run("launch", testOperatorStartOnceLaunch)
+}
+
+func testOperatorStartOnceWorktree(t *testing.T) {
 	r := newWorktreeStartRig(t, false)
 	if err := r.start(t, func(step string) error {
 		if step == "before add" {
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
index b3addb4c..f2012daf 100644
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -88,7 +88,9 @@ type serviceEnv struct {
 	runnerSleep     func(time.Duration)
 	runnerProcesses func() ([]worldProcess, error)
 
-	operatorDeadline func(time.Duration, context.Context)
+	operatorDeadline   func(time.Duration, context.Context)
+	operatorSupervisor func(operatorLaunchJob, io.Writer) error
+	operatorAgentRun   operatorAgentRunner
 }
 
 func realServiceEnv() serviceEnv {
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 238c746e..7c99ca31 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -2083,3 +2083,17 @@ python3 scripts/surface/import_github_checks.py --repo "$PWD" --repo-id world \
 
 The card shows when the checks were captured and links to GitHub. Rerun the
 import to refresh a pending result. Receipts attach to the exact PR head.
+
+### Managed headless start launch
+
+The installed operator runner starts one confined Codex child after an approved
+start reaches its prepared checkpoint. It rechecks the session, key and signed
+executable before taking the durable run reservation. The worker holds the
+runner and episode locks. An admitted child can finish through a daemon restart.
+
+Run `ailang-worldd service runner launch-status <action-id>`. The log and status
+use the action ID and report `launched`, `stopped`, `failed`, `capped` or
+`spawn-refused`. `uncertain` preserves the effect fence for attended review of
+the process and receipts. `ailang-worldd episode report <episode>` shows the exit
+and transcript after completion. Saved streams scrub session credentials. Use the attended
+runner stop or paired upgrade to drain a live launch.
diff --git a/host/agentrun/agentrun.go b/host/agentrun/agentrun.go
index 4950804b..8c3a4235 100644
--- a/host/agentrun/agentrun.go
+++ b/host/agentrun/agentrun.go
@@ -24,7 +24,19 @@ type Result struct {
 // Stdin is /dev/null. Cancellation and the wall cap kill the entire group;
 // received terminal signals are forwarded, then bounded by the same wall cap.
 func Run(ctx context.Context, argv []string, dir string, env []string, stdout, stderr *os.File, timeout time.Duration) (Result, error) {
-	return run(ctx, argv, dir, env, stdout, stderr, nil, timeout)
+	return run(ctx, argv, dir, env, stdout, stderr, nil, timeout, LaunchHooks{})
+}
+
+// LaunchHooks issues a durable effect grant before exec and records the
+// actual child after exec. Files retain the owner's kernel exclusions.
+type LaunchHooks struct {
+	Files   []*os.File
+	Before  func() error
+	Started func(pid int) error
+}
+
+func RunObserved(ctx context.Context, argv []string, dir string, env []string, stdout, stderr *os.File, timeout time.Duration, hooks LaunchHooks) (Result, error) {
+	return run(ctx, argv, dir, env, stdout, stderr, nil, timeout, hooks)
 }
 
 // RunInteractive runs in the terminal foreground group. Terminal signals already
@@ -33,10 +45,10 @@ func RunInteractive(ctx context.Context, argv []string, dir string, env []string
 	if tty == nil {
 		return Result{ExitCode: -1, Reason: "exit"}, errors.New("agentrun: terminal required")
 	}
-	return run(ctx, argv, dir, env, tty, tty, tty, timeout)
+	return run(ctx, argv, dir, env, tty, tty, tty, timeout, LaunchHooks{})
 }
 
-func run(ctx context.Context, argv []string, dir string, env []string, stdout, stderr, tty *os.File, timeout time.Duration) (Result, error) {
+func run(ctx context.Context, argv []string, dir string, env []string, stdout, stderr, tty *os.File, timeout time.Duration, hooks LaunchHooks) (Result, error) {
 	r := Result{ExitCode: -1, Reason: "exit"}
 	if len(argv) == 0 || timeout <= 0 || stdout == nil || stderr == nil {
 		return r, fmt.Errorf("agentrun: argv, output files and positive wall cap required")
@@ -46,6 +58,7 @@ func run(ctx context.Context, argv []string, dir string, env []string, stdout, s
 	cmd := exec.CommandContext(runCtx, argv[0], argv[1:]...)
 	cmd.Dir, cmd.Env = dir, env
 	cmd.Stdout, cmd.Stderr = stdout, stderr
+	cmd.ExtraFiles = hooks.Files
 	if tty == nil {
 		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
 		cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
@@ -65,10 +78,23 @@ func run(ctx context.Context, argv []string, dir string, env []string, stdout, s
 	if err != nil {
 		return r, err
 	}
+	if hooks.Before != nil {
+		if err := hooks.Before(); err != nil {
+			release()
+			return r, err
+		}
+	}
 	if err := cmd.Start(); err != nil {
 		release()
 		return r, err
 	}
+	var receiptErr error
+	if hooks.Started != nil {
+		receiptErr = hooks.Started(cmd.Process.Pid)
+		if receiptErr != nil {
+			_ = cmd.Cancel()
+		}
+	}
 	done := make(chan struct{})
 	forwarded := make(chan bool, 1)
 	go func() {
@@ -92,6 +118,9 @@ func run(ctx context.Context, argv []string, dir string, env []string, stdout, s
 	err = procbound.Wait(cmd.Wait, timeout+time.Second, release)
 	close(done)
 	seen := <-forwarded
+	if errors.Is(err, procbound.ErrCleanupIncomplete) {
+		return r, errors.Join(receiptErr, err)
+	}
 	if cmd.ProcessState != nil {
 		r.ExitCode = cmd.ProcessState.ExitCode()
 	}
@@ -102,9 +131,6 @@ func run(ctx context.Context, argv []string, dir string, env []string, stdout, s
 	} else if seen || r.ExitCode < 0 {
 		r.Reason = "signal"
 	}
-	if errors.Is(err, procbound.ErrCleanupIncomplete) {
-		return r, err
-	}
 	if tty != nil {
 		var exit *exec.ExitError
 		if errors.As(err, &exit) || runCtx.Err() != nil {
@@ -128,11 +154,11 @@ func run(ctx context.Context, argv []string, dir string, env []string, stdout, s
 		}
 	}
 	if killErr != nil && !procbound.NothingToKill(killErr, groupGone) {
-		return r, errors.Join(err, killErr)
+		return r, errors.Join(receiptErr, err, killErr)
 	}
 	var exit *exec.ExitError
 	if errors.As(err, &exit) || runCtx.Err() != nil {
-		return r, nil
+		return r, receiptErr
 	}
-	return r, err
+	return r, errors.Join(receiptErr, err)
 }
diff --git a/host/store/operator_launch.go b/host/store/operator_launch.go
index c450a0af..206672dc 100644
--- a/host/store/operator_launch.go
+++ b/host/store/operator_launch.go
@@ -34,6 +34,7 @@ type OperatorLaunchSnapshot struct {
 	Key     OperatorStepKey
 	Phase   string
 	Receipt OperatorLaunchReceipt
+	NoSpawn bool
 }
 
 func operatorLaunchKey(s OperatorStartRecovery, handoff [32]byte) OperatorStepKey {
@@ -183,6 +184,12 @@ func readOperatorLaunch(ctx context.Context, tx *sql.Tx, id string) (OperatorLau
 		!operatorLaunchDigest(s.Receipt.Executable) {
 		return s, row, ErrOperatorConflict
 	}
+	if len(row.Action) != 0 {
+		a, err := decodeOperatorMintAction(row)
+		if err != nil || s.Receipt.Executable != mintFact(a, "launch_binary_sha256") {
+			return s, row, ErrOperatorConflict
+		}
+	}
 	switch s.Receipt.Status {
 	case "launched":
 		if s.Key.Version == 3 && row.State == "succeeded" && s.Receipt.PID > 0 && s.Receipt.Birth != "" {
@@ -190,6 +197,7 @@ func readOperatorLaunch(ctx context.Context, tx *sql.Tx, id string) (OperatorLau
 		}
 	case "spawn-refused":
 		if (s.Key.Version == 2 || s.Key.Version == 3) && row.State == "failed" && s.Receipt.PID == 0 && s.Receipt.Birth == "" {
+			s.Receipt, s.NoSpawn = OperatorLaunchReceipt{}, true
 			return s, row, nil
 		}
 	case "stopped", "failed", "capped":
diff --git a/host/store/operator_launch_test.go b/host/store/operator_launch_test.go
index 5ee7b307..7af9a654 100644
--- a/host/store/operator_launch_test.go
+++ b/host/store/operator_launch_test.go
@@ -1,6 +1,7 @@
 package store
 
 import (
+	"encoding/json"
 	"errors"
 	"reflect"
 	"testing"
@@ -8,9 +9,29 @@ import (
 	"github.com/sunholo-data/ailang-world/host/approval"
 )
 
-func TestOperatorLaunchReceiptLaw(t *testing.T) {
+func preparedOperatorLaunchFixture(t *testing.T, ttl int64) (operatorPinFixture, string, OperatorStartRecovery) {
+	t.Helper()
 	digest := "sha256:" + hex64("fixture executable")
-	f := newOperatorPinFixtureWithFacts(t, false, []approval.OperatorResolved{{Name: "launch_binary_sha256", Value: digest}}, startRecoveryTestContext(t))
+	prefix := []approval.OperatorResolved{{Name: "launch_binary_sha256", Value: digest}}
+	f := newOperatorPinFixtureWithFacts(t, false, prefix, startRecoveryTestContext(t))
+	if ttl > 0 {
+		row := getOperator(t, f.c, f.in.Attempt.CAS.ID, pinNow)
+		a, err := approval.DecodeOperatorAction(row.Action)
+		if err != nil {
+			t.Fatal(err)
+		}
+		f.in.Start.TTLS = ttl
+		facts, err := f.in.Resolution()
+		if err != nil {
+			t.Fatal(err)
+		}
+		a.Resolution = append(prefix, facts...)
+		f.proof = signOperatorPinAction(t, f.c, row, a)
+		f.in.Attempt.CAS.Digest, err = approval.OperatorSignatureDigest(a)
+		if err != nil {
+			t.Fatal(err)
+		}
+	}
 	f.mint(t)
 	id := f.in.Attempt.CAS.ID
 	for _, value := range []string{"files-identity", OperatorFilesPublished} {
@@ -22,11 +43,16 @@ func TestOperatorLaunchReceiptLaw(t *testing.T) {
 			t.Fatal(err)
 		}
 	}
-	before := f.snapshot(t)
 	prepared, err := f.c.ReadOperatorStartRecovery(f.ctx, id)
 	if err != nil || prepared.Window != OperatorStartPrepared {
 		t.Fatalf("prepared: %+v %v", prepared, err)
 	}
+	return f, digest, prepared
+}
+
+func TestOperatorLaunchReceiptLaw(t *testing.T) {
+	f, digest, prepared := preparedOperatorLaunchFixture(t, 0)
+	id, before := prepared.CAS.ID, f.snapshot(t)
 	handoff := [32]byte{1}
 	now := func() int64 { return pinNow }
 	f.assertRefused(t, before, f.c.ReserveOperatorStep(f.ctx, operatorLaunchKey(prepared, handoff), pinNow), ErrOperatorInput)
@@ -72,6 +98,27 @@ func TestOperatorLaunchReceiptLaw(t *testing.T) {
 		t.Fatalf("observed: %+v %v", live, err)
 	}
 	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, reserved, r, nil, pinNow), ErrOperatorConflict)
+	bad := r
+	bad.Executable = "sha256:" + hex64("wrong observed executable")
+	raw, err := json.Marshal(bad)
+	if err != nil {
+		t.Fatal(err)
+	}
+	operatorExec(t, f.c, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", raw)
+	corrupt := f.snapshot(t)
+	_, err = f.c.ReadOperatorLaunch(f.ctx, id)
+	f.assertRefused(t, corrupt, err, ErrOperatorConflict)
+	raw, _ = json.Marshal(r)
+	operatorExec(t, f.c, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", raw)
+	if owner, err := ReadOperatorEpisodeOwner(f.ctx, f.c.lock.dbPath, prepared.Episode); err != nil || owner != id {
+		t.Fatalf("observed child lost terminal effect exclusion: %q %v", owner, err)
+	}
+	if pending, err := EpisodeHasUnresolvedOperatorSteps(f.ctx, f.c.lock.dbPath, prepared.Episode); err != nil || !pending {
+		t.Fatalf("observed child lost GC effect exclusion: %t %v", pending, err)
+	}
+	if err := f.c.PruneOperatorActions(f.ctx, pinNow+31*86400); err != nil || !reflect.DeepEqual(before, f.snapshot(t)) {
+		t.Fatalf("pruner changed an unresolved receipt: %v", err)
+	}
 	r.Status = "capped"
 	r.Exit = -1
 	history := EpisodeLaunch{Mode: "headless", Exit: r.Exit}
@@ -90,3 +137,37 @@ func TestOperatorLaunchReceiptLaw(t *testing.T) {
 		t.Fatalf("completion replay: %v", err)
 	}
 }
+
+func TestOperatorLaunchSession(t *testing.T) {
+	f, _, ready := preparedOperatorLaunchFixture(t, 5)
+	id, secret := ready.CAS.ID, [32]byte{1}
+	now := func() int64 { return pinNow }
+
+	before := f.snapshot(t)
+	f.assertRefused(t, before, f.c.ReserveOperatorLaunch(f.ctx, ready, secret, f.proof, func() int64 { return pinNow + 5 }), ErrOperatorExpired)
+	if err := f.c.ReserveOperatorLaunch(f.ctx, ready, secret, f.proof, now); err != nil {
+		t.Fatal(err)
+	}
+	reserved, err := f.c.ReadOperatorLaunch(f.ctx, id)
+	if err != nil || reserved.Receipt != (OperatorLaunchReceipt{}) {
+		t.Fatalf("reserved: %+v %v", reserved, err)
+	}
+	before = f.snapshot(t)
+	f.assertRefused(t, before, f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, func() int64 { return pinNow + 5 }), ErrOperatorExpired)
+	operatorExec(t, f.c, "UPDATE session_credentials SET project='changed'")
+	f.assertRefused(t, f.snapshot(t), f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, now), ErrOperatorConflict)
+	operatorExec(t, f.c, "UPDATE session_credentials SET project=?", f.in.Profile.Project)
+	operatorExec(t, f.c, `INSERT INTO session_credentials
+		(credential_id,episode_id,grants_json,expires_at,created_at,project,profile_digest)
+		SELECT ?,episode_id,grants_json,expires_at,created_at,project,profile_digest FROM session_credentials`, hex64("other credential"))
+	f.assertRefused(t, f.snapshot(t), f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, now), ErrOperatorConflict)
+	operatorExec(t, f.c, "DELETE FROM session_credentials WHERE credential_id=?", hex64("other credential"))
+	operatorExec(t, f.c, "UPDATE approvers SET revoked_at=?", pinNow)
+	before = f.snapshot(t)
+	if err := f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, now); err == nil || !reflect.DeepEqual(before, f.snapshot(t)) {
+		t.Fatalf("revoked approver changed a reservation: %v", err)
+	}
+	operatorExec(t, f.c, "UPDATE approvers SET revoked_at=NULL")
+	operatorExec(t, f.c, "DELETE FROM session_credentials")
+	f.assertRefused(t, f.snapshot(t), f.c.AuthorizeOperatorSpawn(f.ctx, id, secret, f.proof, now), ErrOperatorConflict)
+}
diff --git a/host/verbs/operator_policy_test.go b/host/verbs/operator_policy_test.go
index f81eb26f..ddabed9b 100644
--- a/host/verbs/operator_policy_test.go
+++ b/host/verbs/operator_policy_test.go
@@ -70,6 +70,7 @@ func TestOperatorNeverDispatch(t *testing.T) {
 	for _, line := range []string{"approver enroll --name mark-key", "approver revoke mark-key", "session mint w1", "session revoke w1",
 		"service install", "service", "serve", "setup", "ailang-worldd approver enroll", "ailang-worldd service uninstall",
 		"service runner reconcile-start oa-0123456789abcdef0123456789abcdef",
+		"service runner launch-status oa-0123456789abcdef0123456789abcdef",
 		"world-publish approve sha256:x", "world-publish"} {
 		a := admitLine(t, line, allYes)
 		refused(t, a, line, OperatorNever, "-", "is never run from the app")
@@ -78,7 +79,7 @@ func TestOperatorNeverDispatch(t *testing.T) {
 		}
 	}
 	for _, line := range []string{"sh -c ls", "bash", "episode launch w1 --preset p --headless", "episode report w1", "call x", "log", "why head",
-		"approve mr-0123456789abcdef", "world-publishx", "commit"} {
+		"approve mr-0123456789abcdef", "world-publishx", "commit", "operator-launch-worker oa-0123456789abcdef0123456789abcdef"} {
 		argv, _ := Split(line)
 		v, _ := Lookup(argv)
 		refused(t, admitLine(t, line, allYes), line, OperatorOther, "", v.Reason)
````

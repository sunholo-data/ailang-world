# w-world-address-bar-verbs-v4 — S9c sprint plan

## 0. Scope and outcome

**Outcome B.** Row 243, iteration 450. The executor runs on `codex:gpt-6.2-sol`. M1 activates finish from the World app through its real approval routes and the operator's user-level runner. Request, enrolled-key confirmation, runner claim, pinned checks and the local Git journal all ship in the same milestone. Appendix P M2 is a compiled S9d appendix for the separable terminal **read** command `service runner finish-status ACTION-ID`. It is outside this iteration's executor and commit scope.

The release cut preserves reviewed design r1, quorum r2 PROCEED 3/3. D-WORLD-119 and D-WORLD-120 are A. D-WORLD-122 is OPEN with default B, so the walkthrough uses a fake preset. No schema or AILANG source change. The public effect-bearing CLI still meets F1 before flags, ambient inputs or the worktree. S9's exact no-effect help remains available.

The durable hand-off is a typed `local-commit-done` observation with the exact SHA and empty `v3_id`. Browser status includes `checks` and `localCommit`. S10 starts with this receipt and owns the linked V3 card, proof and pre-dispatch fence. S11 owns push, label and PR creation, finish recording, detach, revoke and the attended real-preset exit. Observe the count of affected rows, refs and calls on every fixture path. A local finish is complete for this sprint when those counts show a single commit and zero remote writes, V3 puts and finish records.

## 1. Measurements

### 1.1 Base and provenance

- `91befdab7fd4081b4dc127cc96cc4eeb1723cffa` was clean at start. `git diff --stat d9635675 91befdab -- cmd host docs app` printed no entries. The S9b M1 and M2 receipts are `5717be11` and `f11d09c7`.
- Independently dry-ran and applied iteration 449's `M3.patch` on a fresh archive with `/usr/bin/patch -p1 -E`. Both completed with no offset, fuzz or reject. The input patch measures 1,243 changed lines. See `$S/inherited-apply.log`.
- The local journal follows design §3.3's check/outcome and commit rows. Authority and phases come from §3.1–3.2, runner liveness and lock rules from §2.2, admission and effect boundaries from §4.1 and §4.3 steps 1–2. §5.4 drives the branch table, §6–8 placement and sprint split. §11 E2 and W33 require an independent noninteractive Git read after the approved GitHub reads. W15–W20, W28 and W37 anchor the effect separation and F1.
- **Rehearsal findings.** S9b's draft checked attribution only for dirty paths, routed reserved restart straight to reconciliation outside the private worker, and copied internal check evidence locators into status. M1 includes committed base-to-HEAD paths in attribution, gives reserved restart the same private entry and projects only check name and result. The draft's app export inventory also needed the new Node-driven render seam. The real-listener legs were previously unverified and remain outside gates here. The same-UID config timing window recorded by judge 449 remains within §6’s trusted-UID boundary. The inherited recovery-first selector can elect the same unresolved owned action on successive ticks. Multi-episode fairness remains unmeasured.

### 1.2 Size and cut

| Stage | Incremental | Cumulative `-- cmd host docs app` | Cumulative unfiltered product |
|---|---:|---:|---:|
| M1, land | 1399 | +1335 / −64 = 1399 | 1399 |
| M2, S9d appendix | 57 | +1386 / −66 = 1452 | 1452 |

Measured with `git diff -M --no-index --numstat` on the two archive trees. The controller's canonical landed measurement is `git diff -M --numstat 91befdab HEAD -- cmd host docs app`, summing insertions and deletions. No product files lie outside the filtered roots. Including the two planning artifacts adds 5306 lines to the unfiltered landing diff, for 6705 lines total. M1 fits the 1,400 limit with every load-bearing gate in the activation path. The CLI mirror drives the combined change past 1,400. There is no reason to use the extended cap.

### 1.3 Baseline and budgets

The pristine fail-set source is `$STATE/base` at the pinned head. The comparable 11-test-skipped cmd race leg, which the controller base recording omits, goes to `$S/base-extra`. S9b-parent sets at iteration 449 are historical evidence. Compare only complete package-qualified results with the pinned head in gates and mutation verdicts.

Merge CI 46149287254 recorded cmd race 292.3 s, daemon race 366.1 s, broker race 274.8 s, verifygate race 315.4 s and surfacepage race 5.8 s. CI limits every package to 900 s. S9c limits the complete added named cmd and daemon totals to 15 and 12 s. Each timing attempt records `uptime` and its individual test times. Up to three readings establish the minimum, so one contended over-budget reading alone cannot fail the budget. The old 1 s heartbeat write budget stays in production. New fixtures inject a clock and the existing activation budget. Whole pkg limits are 45 m cmd and daemon, 20 m verifygate and broker, and 25 m surfacepage.

### 1.4 Production callers

| Call | Base | M1 | S9d appendix |
|---|---:|---:|---:|
| RequestOperatorAction | 1 | 1 | 1 |
| JudgeOperatorAction | 1 | 1 | 1 |
| store ConfirmOperatorAction | 1 | 1 | 1 |
| ClaimOperatorAction | 1 | 2 | 2 |
| runEpisodeFinishWithPrep | 0 | 1 | 1 |
| Private local Git commit | 0 | 1 | 1 |
| Attended Git commit | 1 | 1 | 1 |
| worktree Push / dispatch Push / label create / PR view / PR create / PutVerbRequest | 1 each | 1 each | 1 each |

`callers.sh` prints actual locations and requires live controls at the terminal dispatch and the installed runner loop. The private worker call resolves to `consumeFinish`, and its local Git effect to the closure constructed inside `finishWork`. The real runner loop chooses that turn only for admitted finish readiness. The daemon keeps one request/confirm path with closed kind-specific admission.

## 2. Milestones and five activation gates

**M1 — Activate signed local finish.** Use the installed runner's allowlisted remotes and read-only catalogue. Taskless bar review freezes the request from its generation. The existing guarded route stores an inert action and the existing approval store transaction checks and records the enrolled key ceremony. The runner judges the exact current proof, checks the sandbox policy and the independent Git readiness, takes the episode flock, resolves the completed launch and original live session, and verifies Git and the frozen manifest before the claim. A current same-store runner performs the CAS once. The worker runs the private check and commit journal with fresh effect checks and exact reservation/reconciliation. It reports the redacted browser status. There is one product-bearing milestone.

1. **Real admission to worker.** `TestOperatorFinishAdmission` sends the exact QUICKSTART command over `httptest` and the daemon mux. Check readiness refusal, task refusal, stale catalogue, exact signed card, authentic passkey, duplicate request/confirmation and unchanged evidence. `TestOperatorFinishRouteWorker` completes the fake child, takes the same request route, confirms with the repo SSH signer and enters the actual runner loop. The worker must own the claim, send the pinned sandbox RPC once with the acknowledged session and argv, write one check outcome and reserve and make one real local commit. Verify exact typed receipt, one descendant, no second turn for the same receipt, all V3 rows absent, finish registry head absent and unchanged local bare remote. The integrated fixture uses a fake RPC response for workspace-exec. Run the subprocess and sandbox-profile controls `TestExecProtectedSettingsThroughHandler` and `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` alongside the route. Those controls observe the executed fake tool.
2. **Eligibility and drift.** The worker fixture refuses missing start, live launch, degraded episode, unavailable report, empty dirty tree, changed untracked bytes, HEAD/base ref, remote, preset/profile, token/expiry and unsafe repository helper. Refusals assert classes, scrubbed remedy, unchanged durable evidence and zero check/outcome calls. `TestOperatorFinishManifestCapture` puts an unattributed file in a real pre-existing commit and verifies the base-to-HEAD attribution alongside dirty paths. `TestOperatorFinishReadReplacement`, `TestOperatorFinishWorktreeFence`, `TestOperatorFinishConfigReplacement` and `TestOperatorCommitReservationFence` cover inode replacement, no-follow/fifo, staged input and config, detached or changed identities and effect-fence races. Auth refusal runs on both confirmed and owned turns and sees the fake Git authentication failure despite fake gh success.
3. **Episode exclusion.** The same route fixture takes the actual terminal lock before the worker and calls the terminal acquisition during the check while the runner owns the turn. Read the owner with the 5 s test bound. Observe conflict and the active action id in both directions, current journal hold during the check and no claim during terminal ownership. Pair with `TestOperatorTerminalExclusion` and `TestOperatorEpisodeLockContendsWithoutWaiting` for the production terminal entry, deadline and non-waiting acquisition controls. No lock is inferred from a heartbeat.
4. **Privacy and application behavior.** The signed card is the positive disclosure control. Requests and Needs-you responses, status JSON, daemon and runner diagnostic sinks exclude task, artifact, token, manifest, path and internal check-evidence markers. The Node harness renders the actual taskless finish panel, verifies focus and the explicit Request approval click, sends no request on opening, and reconnects the same operation and typed link. Exact CLI projection coverage is compiled in the optional S9d mirror.
5. **Evidence gates.** The selected, whole-package, outside and mutation instruments are in §6. Complete ledger, named and package time checks compare against the pinned pristine base. The listener and Chrome tests are outside prerequisites. Verify on the planner stage before handing it to the executor. An instrument failure or a named surviving load-bearing mutation blocks advancement.

**Appendix P M2 — S9d CLI mirror, separately scheduled.** Expose existing redacted finish receipt fields through the read-only service diagnostic and test that both of its sinks are scrubbed. The command accepts an exact action id, reads local durable receipts and reports name/result and retained SHA. It adds no claim, remote dispatch or journal write. It is optional once M1's browser status is available.

## 3. Operator output

The World bar accepts the reviewed finish command with the pinned preset. Enter opens the taskless review and stays inert. **Request approval** opens the signed card. An enrolled key approves the frozen local intent. **Refresh status** and **Action status (JSON)** expose the current action, `checks` (name and result) and `localCommit` (SHA), plus durable reasons for unavailable readiness, reserved uncertainty, definite failure and retained commit. A Git read or GitHub readiness failure gives `blocked: github-auth` and the attended PATH, login or credential-helper remedy. Eligibility failures direct the operator to the stopped launch and original session, drift asks for fresh review, and lock contention names the action id and the service runner diagnostic.

`TestOperatorFinishAdmission` reads the exact new QUICKSTART example and runs it through the listener as printed. `TestOperatorFinishRouteWorker` drives the corresponding fake launch through one local commit. `TestOperatorFinishBarReview` drives the application control in Node. S9d adds `service runner finish-status ACTION-ID` as a second read surface.

## 4. Acceptance, timing and fail sets

The fixture scopes are temp World and credentials stores with real migrations, temp repositories and bare remotes, fake pinned gh/Git readiness programs, the existing fake child and authentic test SSH/passkey signers. Runtime writes are observed as action/step/counter/ref differences. A missing file, step or registry head is a legal reading. Browser records are isolated with the pre-existing temporary profile harness. Physical paths use `filepath.EvalSymlinks` on both sides of comparisons. Tests that enter `runnerPath` are repeated with a private TMPDIR child under a 1777 scratch parent. Each new write obtains a bounded test context. The only clock-based decisions use injected clocks.

| New test | Package | Measured race seconds |
|---|---|---:|
| `TestOperatorFinishReadiness` | cmd/ailang-worldd | 2.76 |
| `TestOperatorFinishResumeEvidence` | cmd/ailang-worldd | 0.67 |
| `TestOperatorCommitFailureInputRecheck` | cmd/ailang-worldd | 1.31 |
| `TestOperatorCommitResumeReserved` | cmd/ailang-worldd | 1.85 |
| `TestOperatorFinishRouteWorker` | cmd/ailang-worldd | UNINFORMATIVE UNDER SANDBOX |
| `TestOperatorFinishAdmission` | host/daemon | UNINFORMATIVE UNDER SANDBOX |
| `TestOperatorFinishBarReview` | host/surfacepage | 0.14 |

Minimum measured local named cmd subtotal 6.59 s. The complete cmd and daemon named totals remain outside gates.

- Reading 1, 9.25 s cmd at Sat Nov 28 15:10:19 UTC 2026, 16:10  4 users, load averages: 31.13 27.70 19.50. `TestOperatorFinishResumeEvidence` 0.93 s, `TestOperatorCommitFailureInputRecheck` 2.88 s, `TestOperatorCommitResumeReserved` 1.89 s, `TestOperatorFinishReadiness` 3.55 s. Node test 0.11 s at Sat Nov 28 15:11:18 UTC 2026, 16:11  4 users, load averages: 25.26 26.65 19.66.
- Reading 2, 10.37 s cmd at Sat Nov 28 15:13:53 UTC 2026, 16:13  4 users, load averages: 22.13 25.65 20.46. `TestOperatorFinishResumeEvidence` 0.42 s, `TestOperatorCommitFailureInputRecheck` 1.51 s, `TestOperatorCommitResumeReserved` 3.90 s, `TestOperatorFinishReadiness` 4.54 s. Node test 0.15 s at Sat Nov 28 15:14:08 UTC 2026, 16:14  4 users, load averages: 24.70 26.06 20.70.
- Reading 3, 6.59 s cmd at Sat Nov 28 15:14:13 UTC 2026, 16:14  4 users, load averages: 24.49 25.99 20.71. `TestOperatorFinishResumeEvidence` 0.67 s, `TestOperatorCommitFailureInputRecheck` 1.31 s, `TestOperatorCommitResumeReserved` 1.85 s, `TestOperatorFinishReadiness` 2.76 s. Node test 0.14 s at Sat Nov 28 15:14:23 UTC 2026, 16:14  4 users, load averages: 21.55 25.32 20.53.

Supplemental in-memory race measurements use a transport overlay, with no real-listener cap verdict. `TestOperatorFinishRouteWorker` 22.91 s at Sat Nov 28 15:11:23 UTC 2026, 16:11  4 users, load averages: 25.16 26.61 19.69. `TestOperatorFinishRouteWorker` 5.40 s at Sat Nov 28 15:16:54 UTC 2026, 16:16  4 users, load averages: 8.50 18.86 18.70. `TestOperatorFinishRouteWorker` 5.22 s at Sat Nov 28 15:17:04 UTC 2026, 16:17  4 users, load averages: 7.51 18.30 18.51. `TestOperatorFinishAdmission` 1.52 s at Sat Nov 28 15:12:13 UTC 2026, 16:12  4 users, load averages: 26.82 26.76 20.13.

Whole `host/surfacepage` race: pinned base 5.903 s, local M1 19.046 s, observed delta +13.143 s. Whole cmd and daemon race deltas require the outside sweep.

Selected local readings in `$S/local-gates-M1-final/legs.txt`:

- LEG gofmt OK 0s
- LEG build OK 2s
- LEG vet OK 1s
- LEG linux-vet OK 4s
- LEG test-compile OK 19s
- LEG inventory-store OK 1s
- LEG inventory-cli OK 1s
- LEG inventory-daemon OK 1s
- LEG inventory-boundary OK 1s
- LEG subprocess-sink OK 1s
- LEG subprocess-sites UNINFORMATIVE UNDER SANDBOX rc=1
- LEG js-syntax OK 0s
- LEG verify-ail OK 56s
- LEG caller-census OK 2s
- LEG accept OK 24s
- LEG accept-race OK 30s
- LEG accept-1777 OK 6s
- LEG cmd-plain SKIPPED (outside.sh whole)
- LEG cmd-race SKIPPED (outside.sh whole)
- LEG daemon-plain SKIPPED (outside.sh whole)
- LEG daemon-race SKIPPED (outside.sh whole)
- LEG verifygate-plain SKIPPED (outside.sh whole)
- LEG broker-plain SKIPPED (outside.sh whole)
- LEG baseline-completion SKIPPED (controller recording pristine legs)
- LEG compare-cmd-plain SKIPPED (stage evidence unavailable)
- LEG compare-cmd-race SKIPPED (stage evidence unavailable)
- LEG compare-daemon-plain SKIPPED (stage evidence unavailable)
- LEG compare-daemon-race SKIPPED (stage evidence unavailable)
- LEG compare-verifygate-plain SKIPPED (stage evidence unavailable)
- LEG compare-broker-plain SKIPPED (stage evidence unavailable)
- LEG compare-surfacepage-plain OK
- LEG compare-surfacepage-race OK

Both listener files showed `bind: operation not permitted` (rc 1). The in-process mux diagnostics passed `TestOperatorFinishAdmission` and `TestOperatorFinishRouteWorker` separately. The real listeners remain prerequisites.

Reran changed-file gofmt, build, native and Linux vet plus JS syntax on the sealed M1 bytes.

M2 changed-file gofmt, build, native and Linux vet plus JS syntax succeeded. The S9d CLI field and redaction assertion also passed under the in-memory route overlay.

The socket-free cmd QUICKSTART neighbors passed on both stages. `TestQuickstartSection9FlagsMatchTheCLI` reached a loopback bind denial and stays in the outside cmd stem leg. Test discovery found all 157 named selectors.

Pinned baseline readings available at plan close:

- LEG verifygate-plain rc=1 parser=0 timeouts=0 fails=1 secs=471
- LEG broker-plain rc=0 parser=0 timeouts=0 fails=0 secs=262
- LEG surfacepage-plain rc=0 parser=0 timeouts=0 fails=0 secs=11
- LEG surfacepage-race rc=0 parser=0 timeouts=0 fails=0 secs=8
- LEG cmd-plain rc=1 parser=0 timeouts=0 fails=11 secs=387
- LEG daemon-plain rc=1 parser=0 timeouts=0 fails=3 secs=484
- LEG daemon-race rc=1 parser=0 timeouts=0 fails=5 secs=1078

- stage ok  	github.com/sunholo-data/ailang-world/host/surfacepage	10.161s
- stage ok  	github.com/sunholo-data/ailang-world/host/surfacepage	19.046s
- pristine ok  	github.com/sunholo-data/ailang-world/host/surfacepage	5.453s
- pristine ok  	github.com/sunholo-data/ailang-world/host/surfacepage	5.903s

Raw per-leg records disambiguate plain and race.

**Outside execution on M1, before the executor.** `outside.sh` takes the tree root at argument 1. Uninformative listener output is an observed sandbox bind denial here. Memory-mux overlays are supplemental fixture diagnostics. They carry no listener or race-budget verdict.

````bash outside-commands
S="$HOME/.ailang/state/world-iter450/plan-scratch"
bash "$S/outside.sh" "$S/st-M1" M1 focused
bash "$S/outside.sh" "$S/st-M1" M1 chrome
bash "$S/gates.sh" "$S/base-tree" "$S/base-extra" M1 base-extra
bash "$S/outside.sh" "$S/st-M1" M1 whole
bash "$S/named-mutants.sh" "$S/st-M1" "$S/controller-named-M1-$(date +%s)" outside M1
bash "$S/whole-mutants.sh" "$S/st-M1" "$S/whole-mutant-ids.txt" M1
````

The deferred mirror can be checked independently with `bash "$S/outside.sh" "$S/st-M2" M2 focused` and `bash "$S/named-mutants.sh" "$S/st-M2" "$S/controller-named-M2-$(date +%s)" outside M2`. Run the selected and Chrome commands separately from long base/whole sweeps. Keep the 11 known real-binary CLI tests excluded from every cmd race command. Known daemon real-binary and load-sensitive failures remain governed by the observed pinned-base fail sets. Every package leg uses an explicit timeout at or above 1.5× its measured base where available. The parser prints a firing control, rejects a nonzero exit with no test fail set, and rejects test timeouts as incomplete measurements.

## 5. Branch traces and mutation gates

### 5.1 S9c branches of design §5.4

| §5.4 row | M1 disposition and behavioral gate | Named mutation |
|---|---|---|
| Public interactive typed finish | Existing typed prompt and response in `TestOperatorFinishTTYParity` | `f1-tty-prompt` |
| Public no-tty finish / flag / environment | `TestOperatorFinishF1` proves F1 before all effect-bearing invocations | `f1-surface-bypass` |
| Missing/wrong prep proof | `TestOperatorFinishNoProof`, `TestOperatorFinishProofScope`, `TestOperatorFinishPrep`; missing, forged, expired, other episode/manifest/kind, daemon-only, wrong owner and replayed proof | `proof-no-recheck`, `proof-deadline`, `proof-bytes`, `j01`, `j02`, `j04` |
| Valid prep completed launch, sandboxed checks and commit | Real routes, worker claim, one check, one commit, typed local-only hand-off in `TestOperatorFinishRouteWorker`. The design row's V3 confirmation begins in S10 | `runner-positive-claim`, `runner-prep-caller`, `runner-commit-caller`, `runner-loop-turn` |
| Missing start, degraded or empty tree, running agent, session changed/expired, read refusal | Worker matrix in `TestOperatorFinishRouteWorker`, auth gate in `TestOperatorFinishReadiness` | `runner-elig-start`, `runner-elig-degraded`, `runner-elig-live-agent`, `runner-elig-token`, `runner-elig-report-read` |
| Stale preset, HEAD, untracked, path race, remote/base or profile | Worker matrix, `TestOperatorFinishStale`, inode and config replacement, exact base and positive attribution pin | `runner-manifest-mismatch`, `runner-drift-preset`, `git-attribution-base-delta`, `j27` |
| Protected/unattributed, unsafe hooks, helpers or special files | `TestOperatorFinishWorktreeFence`, `TestOperatorFinishConfigReplacement`, `TestOperatorFinishManifestCapture`, `TestOperatorFinishPrep` and the route fixture's unsafe helper sentinel | `git-first-use-config`, `git-attribution-base-delta`, `j06`, `j27` |
| Unsandboxed check, including resume | `TestOperatorFinishUnsandboxed` and worker pre-claim sandbox gate | `j05`, `runner-elig-check-sandbox` |
| Missing gh, login or Git helper failure | `TestOperatorFinishAuth`, `TestOperatorFinishAuthInstall`, `TestOperatorFinishReadiness` (confirmed and resume), daemon auth remedy | `ready-before-claim`, `ready-gh-not-git`, `route-finish-health` |
| Failure, abort, timeout, unknown result or changed check inputs | `TestOperatorFinishChecks` and prior-pass World lookup in `TestOperatorFinishResumeEvidence` | `j07`, `j08`, `j10` |
| Definite commit failure | `TestOperatorCommitFailed`, `TestOperatorCommitFailureInputRecheck`, retained checks, no prompt | `j13` |
| Lost commit reply, partial stage, divergent HEAD | `TestOperatorCommitReconcile`, `TestOperatorCommitResumeReserved`, `TestOperatorCommitIdentity`, `TestOperatorCommitReservationFence` | `j14`, `j16`, `j20`, `j22`, `j23`, `j24`, `commit-author-email`, `commit-match-tree`, `commit-before-reservation`, `j32`, `j33` |
| Prep expires after successful checks | `TestOperatorFinishPrepExpiry` verifies no commit and retains checks before fresh authorization | `proof-deadline` and carried journal expiry guard |

The S9c cross-cutting rows are the terminal/runner same-lock conflict (`runner-lock-busy`, `runner-lock-ref`, `terminal-lock-active-id`), origin and exact taskless route admission (`route-finish-task`, `route-operation-generation`, `route-confirm-scope`), status data minimization (`status-evidence-leak`, `status-sha`, `js-taskless-dom` and other executing Node checks), and the empty V3 boundary (real route receipt, count and GitHub command sentinel). The static census complements the behavioral controls.

### 5.2 Later branches

| Owner | Design §5.4 rows and pin names |
|---|---|
| S9d | The read-only diagnostic mirror of the same local check and commit fields, `s9d-cli-dispatch` and `s9d-cli-sha`. No §5.4 authority or execution row moved here. |
| S10 | Prompt size/entropy/store busy/cap (`TestOperatorFinishCardWriteRefusal`), lost-put/restart (`TestOperatorFinishCardRestart`), missing/superseded/expired/revoked or read-failed card (`TestOperatorFinishCardRefusal`), reject typed/foreign reuse (`TestOperatorFinishRejectReuse`), own bytes (`TestOperatorFinishOwnBytes`), F6 after PR approval (`TestOperatorFinishPreDispatchFence`), changed parent runner action at linked card (`TestOperatorFinishLinkedCard`). Judge 447 C-09/C-10/C-10b/C-01/C-02 stay with the browser hand-off. |
| S11 | Equal/ancestor remote and uncertain push (`TestOperatorPushIdentity`, `TestOperatorPushUncertain`), label (`TestOperatorLabelReconcile`), existing/mismatch/duplicate PR (`TestOperatorReadExistingPR`, `TestOperatorPRIdentity`), pre-send/lost response (`TestOperatorPRNoSend`, `TestOperatorPRResponseLost`), browser close to exact record (`TestOperatorFinishBrowser`), observed record after expiry (`TestOperatorFinishObservedRecord`), atomic used/episode heads (`TestOperatorFinishRecordCAS`), scoped cleanup (`TestOperatorFinishCleanup`). |

### 5.3 Judge-449 pins

| ID | Exact original edit's behavioral pin | Result |
|---|---|---|
| j08 | `TestOperatorFinishResumeEvidence` uses production `finishWorldEvidence` for an earlier observed pass with missing World evidence | KILLED TestOperatorFinishResumeEvidence |
| j13 | `TestOperatorCommitFailureInputRecheck` changes inputs at the reserved commit seam and requires a retained uncertain reservation | KILLED TestOperatorCommitFailureInputRecheck |
| j14 | `TestOperatorCommitResumeReserved` creates the exact descendant after reservation, reopens the store under a new worker and verifies read-only adoption | KILLED TestOperatorCommitResumeReserved |
| j16 | `TestOperatorCommitResumeReserved` first activates an unsafe Git config after reservation and observes no adoption or index/HEAD changes | KILLED TestOperatorCommitResumeReserved |
| j19 | `TestOperatorFinishWorldEvidence` commits identical valid outcome payloads under separate stores and tests the expected vs foreign SemanticID | KILLED TestOperatorFinishWorldEvidence |

Optional j03 is redundant because canonical action bytes and the immutable digest pin the id. j17 and j25 remain equivalent under the reviewed receipt/parent constraints. j17's predecessor argument and first-match query already constrain the selected matching entry in the fixture. j25 would require a Git commit to be its own verified parent. j12 is an optional additional boundary pin between tree planning and reservation. The existing expiry control remains selected. A scratch-only optional pin measured +17/−0 and killed the exact j12 edit in `TestOperatorCommitFreshBeforeReservation`. It would move the selected patch over 1,400. Both runs and the diff are in `$S/j12-size`.

### 5.4 Existing tests and helpers edited

| Existing test/helper | Before | After | Mutation control |
|---|---|---|---|
| `operatorActivationFixtureWithTask` | `runOperatorPreparationLoop` with default fake remote | Same runner loop via a variant that pins a fake finish remote | `runner-catalogue-root`, `runner-loop-turn` in the new route |
| `newOperatorLaunchRigWithTask` | `executeOperatorLaunch` under the fake-child supervisor | Same supervisor and credible completed receipt with an approved fake remote fixture | `runner-elig-live-agent` and existing `TestOperatorSpawnRefused` |
| `newWorktreeStartRigWithTask` | `operatorStartFacts` and temp worktree preparation | Same functions with an optional fake finish remote, explicit committer and one check | `runner-drift-preset`, `runner-commit-caller` |
| `TestOperatorFinishManifestCapture` | `operatorFinishInputs.frozen` for dirty snapshots and wrong base | Same method includes committed base-to-HEAD attribution | `git-attribution-base-delta` |
| `TestOperatorFinishWorldEvidence` | `finishWorldEvidence` for an exact World entry | Adds a committed object with a foreign SemanticID using a fresh fixture store | `j19` |
| `TestOperatorFinishConfigReplacement` | Config swap during `operatorFinishGit.snapshot` | Also drives the `operatorFinishGit.commit` config fence using the non-config command sentinel | `j27`, `git-first-use-config` |
| `TestSurfaceAppPure` embedded `appHarness` | Exact module export manifest | Pins the shared browser renderer export that the Node finish test executes | `js-render-export` kills `TestOperatorFinishBarReview` by execution |
| M2 only, `TestOperatorFinishRouteWorker` | M1 real claim and browser receipt | Reads the same durable receipt using the S9d CLI, verifies both new sinks | `s9d-cli-dispatch`, `s9d-cli-sha` |

### 5.5 Selectors and asset inventories

The exact per-file selectors below are also emitted in the `selectors.sh` fence. They run on each stage via local or outside paths.

| Files/asset family | Test names |
|---|---|
| New cmd tests | `TestOperatorFinishReadiness`, `TestOperatorFinishResumeEvidence`, `TestOperatorCommitFailureInputRecheck`, `TestOperatorCommitResumeReserved`, `TestOperatorFinishRouteWorker` |
| Finish preflight, journal, Git, manifest and new config-boundary controls | `TestOperatorFinishF1`, `TestOperatorFinishTTYParity`, `TestOperatorFinishNoProof`, `TestOperatorFinishProofScope`, `TestOperatorFinishAuth`, `TestOperatorFinishToolPins`, `TestOperatorFlagParity`, `TestOperatorTerminalExclusion`, `TestOperatorFinishPrep`, `TestOperatorFinishPrepExpiry`, `TestOperatorFinishStale`, `TestOperatorFinishChecks`, `TestOperatorFinishUnsandboxed`, `TestOperatorCommitFailed`, `TestOperatorCommitReconcile`, `TestOperatorFinishJournalOnce`, `TestOperatorFinishManifestCapture`, `TestOperatorFinishWorldEvidence`, `TestOperatorFinishWorktreeFence`, `TestOperatorCommitIdentity`, `TestOperatorFinishCaptureBounds`, `TestOperatorFinishReadReplacement`, `TestOperatorCommitReservationFence`, `TestOperatorFinishConfigReplacement` |
| Cmd activation and runner neighbors | `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorPreparationRechecksHeartbeatAfterWorktree`, `TestOperatorStartMintRecoveryRemedies`, `TestOperatorCurrentOwnerProcessRecovery`, `TestOperatorAttendedStartDisposition`, `TestOperatorActivationDrainRefusals`, `TestOperatorRevokeStatRefusal`, `TestOperatorServiceFence`, `TestOperatorServiceOwnership`, `TestOperatorServiceConfig`, `TestOperatorServicePair`, `TestOperatorServiceStartup`, `TestOperatorServiceUpgradePairFence`, `TestOperatorServiceCancelConfirmation`, `TestOperatorServiceWriteScope`, `TestOperatorServiceEnvironment`, `TestRunnerPathRefusalClasses`, `TestOperatorRunnerKernelLock`, `TestOperatorStatusDisconnected`, `TestOperatorFinishAuthInstall` |
| Cmd worktree and launch fixture neighbors | `TestOperatorStartOnce`, `TestOperatorStartCollisions`, `TestOperatorWorktreeWindow`, `TestOperatorWorktreeDispositionCAS`, `TestOperatorWorktreeSurvivors`, `TestOperatorWorktreeRecoveryGuards`, `TestOperatorWorktreeReadbackGuards`, `TestOperatorReadbackErrorIsNeverCreated`, `TestOperatorWorktreeLockedReread`, `TestOperatorBorrowedWorktreeGitProcess`, `TestOperatorSpawnRefused`, `TestOperatorLaunchDispatchRechecks`, `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchExit`, `TestOperatorStartBrowser`, `TestOperatorLaunchWorker`, `TestOperatorStartedRequiresBirth`, `TestOperatorSupervisorCap`, `TestOperatorFakePresetEligibility` |
| Cmd QUICKSTART, pinned examples | `TestQuickstartSection9FlagsMatchTheCLI`, `TestQuickstartSection10FlagsMatchTheCLI`, `TestQuickstartSection11FlagsMatchTheCLI`, `TestQuickstartSection18FlagsMatchTheCLI`, `TestQuickstartSection18CardMatchesThePage`, `TestQuickstartSection19MatchesTheBar`, `TestQuickstartSection20MatchesTheFinishFlag`, `TestQuickstartAttendedStepsBoxMatchesSection9`, `TestQuickstartServiceSectionMatchesTheCLIHelp`, `TestQuickstartRunnerMintRecovery`, `TestProfileInitQuickstartLineMatchesTheCLI`, `TestDemoSeedHelpPreparation` |
| Daemon operator/read/confirm/routes, QUICKSTART, Needs-you golden | `TestOperatorFinishAdmission`, `TestOperatorAdmissionHTTP`, `TestOperatorReviewRefusesForgedCatalogue`, `TestOperatorReplayConflicts`, `TestOperatorLinksProjection`, `TestOperatorConfirmHTTP`, `TestOperatorConfirmProjection`, `TestOperatorConfirmationAdmissionNeighbours`, `TestOperatorCardAssets`, `TestSurfaceVerbRouteRegistered`, `TestSurfaceVerbBoundsAndErrors`, `TestSurfaceVerbIsReadOnly`, `TestSurfaceVerbWhyIsReadOnly`, `TestSurfaceVerbOutcomes`, `TestSurfaceAppRoutes`, `TestRequestsJSONConfirmationsGolden`, `TestMCPQuickstartPayloadsVerbatim`, `TestAGUIQuickstartControl`, `TestQuickstartDecisionPacketSection`, `TestWorkbenchLiveQuickstartControl`, `TestSeToolsQuickstartPayloadsVerbatim`, `TestExecStartupRefusalTable` |
| App asset, operator JS and JS pins | `TestOperatorBarStartReview`, `TestOperatorOtherVerbParity`, `TestOperatorLiteralTask`, `TestOperatorDuplicateTap`, `TestOperatorReconnect`, `TestOperatorPanelLatest`, `TestOperatorPanelNetwork`, `TestOperatorLinkAllowlist`, `TestOperatorFinishBarReview`, `TestOperatorCardStatusDisconnected`, `TestOperatorCardLinkAllowlist`, `TestSurfaceAppMatchesWorkingCopy`, `TestSurfaceAppShellAssembly`, `TestSurfaceAppPure`, `TestSurfaceAppVerbs`, `TestSurfaceAppFollow`, `TestSurfaceAppTextIsSelectable`, `TestSurfaceScriptAssetSeams`, `TestAddressHintsArePinned`, `TestNeedsYouListsConfirmations`, `TestNeedsYouConfirmationEmbeddedNewline`, `TestNeedsYouConfirmationsNonArray`, `TestVerbDecidePrefixBoundary`, `TestNeedsYouListsDecisions` |
| Broker and QUICKSTART package controls | `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `TestExecProtectedSettingsThroughHandler`, `TestQuickstartSection10ProfilesLoad`, `TestQuickstartProtectedContract` |

Exact Go stem matches by changed path:

- `cmd/ailang-worldd/operator_activation_test.go` — `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorPreparationRechecksHeartbeatAfterWorktree`, `TestOperatorStartMintRecoveryRemedies`, `TestOperatorCurrentOwnerProcessRecovery`, `TestOperatorAttendedStartDisposition`
- `cmd/ailang-worldd/operator_catalogue.go` — Adjacent behavioral controls appear in the table.
- `cmd/ailang-worldd/operator_finish_manifest.go` — `TestOperatorFinishManifestCapture`, `TestOperatorFinishWorldEvidence`
- `cmd/ailang-worldd/operator_finish_manifest_test.go` — `TestOperatorFinishManifestCapture`, `TestOperatorFinishWorldEvidence`
- `cmd/ailang-worldd/operator_finish_resume_test.go` — `TestOperatorFinishResumeEvidence`, `TestOperatorCommitFailureInputRecheck`, `TestOperatorCommitResumeReserved`
- `cmd/ailang-worldd/operator_finish_route_test.go` — `TestOperatorFinishRouteWorker`
- `cmd/ailang-worldd/operator_finish_runner.go` — `TestOperatorFinishReadiness`
- `cmd/ailang-worldd/operator_finish_runner_test.go` — `TestOperatorFinishReadiness`
- `cmd/ailang-worldd/operator_finish_tree.go` — `TestOperatorFinishWorktreeFence`, `TestOperatorCommitIdentity`, `TestOperatorFinishCaptureBounds`, `TestOperatorFinishReadReplacement`, `TestOperatorCommitReservationFence`, `TestOperatorFinishConfigReplacement`
- `cmd/ailang-worldd/operator_finish_tree_test.go` — `TestOperatorFinishWorktreeFence`, `TestOperatorCommitIdentity`, `TestOperatorFinishCaptureBounds`, `TestOperatorFinishReadReplacement`, `TestOperatorCommitReservationFence`, `TestOperatorFinishConfigReplacement`
- `cmd/ailang-worldd/operator_launch_test.go` — `TestOperatorSpawnRefused`, `TestOperatorLaunchDispatchRechecks`, `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchExit`, `TestOperatorStartBrowser`
- `cmd/ailang-worldd/operator_runner.go` — Adjacent behavioral controls appear in the table.
- `cmd/ailang-worldd/operator_worktree_test.go` — `TestOperatorStartOnce`, `TestOperatorStartCollisions`, `TestOperatorWorktreeWindow`, `TestOperatorWorktreeDispositionCAS`
- `host/daemon/surface_operator.go` — `TestOperatorReviewRefusesForgedCatalogue`, `TestOperatorReplayConflicts`, `TestOperatorAdmissionHTTP`, `TestOperatorLinksProjection`
- `host/daemon/surface_operator_confirm.go` — `TestOperatorConfirmHTTP`, `TestOperatorConfirmProjection`, `TestOperatorConfirmationAdmissionNeighbours`
- `host/daemon/surface_operator_finish_test.go` — `TestOperatorFinishAdmission`
- `host/daemon/surface_operator_read.go` — Adjacent behavioral controls appear in the table.
- `host/surfacepage/app_operator_test.go` — `TestOperatorBarStartReview`, `TestOperatorOtherVerbParity`, `TestOperatorLiteralTask`, `TestOperatorDuplicateTap`, `TestOperatorReconnect`, `TestOperatorPanelLatest`, `TestOperatorPanelNetwork`, `TestOperatorLinkAllowlist`, `TestOperatorFinishBarReview`
- `host/surfacepage/app_test.go` — `TestSurfaceAppMatchesWorkingCopy`, `TestSurfaceAppShellAssembly`, `TestSurfaceAppPure`, `TestSurfaceAppTextIsSelectable`

Additional runners for the affected neighbors are `TestOperatorStartWalkthroughInChrome`, `TestOperatorFilePublicationAdmission` and the package inventory selectors in §6. The whole cmd race skip expression names all 11 real-binary CLI reds exactly.

No golden fixture bytes need regeneration. The attempted post-click finish refocus mutation was equivalent once the request persisted. M1 reuses the existing post-click focus rule and tests the separate taskless first focus. The Node export manifest changes because `showOperatorReview` is the application function reused by the harness. Guard inventory covers that export, daemon writer sites, route behavior, the app asset and unchanged approval/start neighbors. No production context root moved or was added, so its pin count should remain equal.

### 5.6 Mutation ledgers

The named ledger uses exact source anchors in `mutants.json`, with file and line IDs, and the carried judge edits from iteration 449. Only an exact named behavioral failure with no sandbox denial counts as a kill. The first F1 mutant checked an unparsed flag. The rerun checks the raw argv bypass. The first finish-grammar witness did not contain an invalid option at the correct argument count. The Node rerun adds that arm. The post-click refocus check was equivalent and the final app patch removes its added condition. Compile failure, missing package/control, timeout and sandbox denial have independent classifications. The short whole ledger selects the highest-impact cases and compares against `cmd-plain`, `daemon-plain` or `surfacepage-plain` at pristine 91befdab. A result only in an inventory test cannot satisfy a named behavioral row.

**Named ledger.** CONTROL_LIVE 1, KILLED 47, PENDING 2, UNINFORMATIVE UNDER SANDBOX 74.

| ID | Source anchor (stage:line) | Behavioral test | Verdict |
|---|---|---|---|
| `j01` | `cmd/ailang-worldd/operator_finish_preflight.go` (M1:26) | `TestOperatorFinishProofScope` | KILLED TestOperatorFinishProofScope |
| `j02` | `cmd/ailang-worldd/operator_finish_preflight.go` (M1:29) | `TestOperatorFinishProofScope` | KILLED TestOperatorFinishProofScope |
| `j04` | `cmd/ailang-worldd/operator_finish_local.go` (M1:80) | `TestOperatorFinishPrep` | KILLED TestOperatorFinishPrep |
| `j05` | `cmd/ailang-worldd/operator_finish_local.go` (M1:88) | `TestOperatorFinishUnsandboxed` | KILLED TestOperatorFinishUnsandboxed |
| `j06` | `cmd/ailang-worldd/operator_finish_local.go` (M1:93) | `TestOperatorFinishPrep` | KILLED TestOperatorFinishPrep |
| `j07` | `cmd/ailang-worldd/operator_finish_local.go` (M1:109) | `TestOperatorFinishChecks` | KILLED TestOperatorFinishChecks |
| `j08` | `cmd/ailang-worldd/operator_finish_local.go` (M1:113) | `TestOperatorFinishResumeEvidence` | KILLED TestOperatorFinishResumeEvidence |
| `j10` | `cmd/ailang-worldd/operator_finish_local.go` (M1:204) | `TestOperatorFinishChecks` | KILLED TestOperatorFinishChecks |
| `j13` | `cmd/ailang-worldd/operator_finish_local.go` (M1:161) | `TestOperatorCommitFailureInputRecheck` | KILLED TestOperatorCommitFailureInputRecheck |
| `j14` | `cmd/ailang-worldd/operator_finish_local.go` (M1:101) | `TestOperatorCommitResumeReserved` | KILLED TestOperatorCommitResumeReserved |
| `j16` | `cmd/ailang-worldd/operator_finish_local.go` (M1:217) | `TestOperatorCommitResumeReserved` | KILLED TestOperatorCommitResumeReserved |
| `j19` | `cmd/ailang-worldd/operator_finish_manifest.go` (M1:102) | `TestOperatorFinishWorldEvidence` | KILLED TestOperatorFinishWorldEvidence |
| `j20` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:340) | `TestOperatorCommitIdentity` | KILLED TestOperatorCommitIdentity |
| `j26` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:313) | `TestOperatorCommitIdentity` | KILLED TestOperatorCommitIdentity |
| `j22` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:344) | `TestOperatorCommitIdentity` | KILLED TestOperatorCommitIdentity |
| `j23` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:336) | `TestOperatorCommitIdentity` | KILLED TestOperatorCommitIdentity |
| `j24` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:336) | `TestOperatorCommitIdentity` | KILLED TestOperatorCommitIdentity |
| `j27` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:160) | `TestOperatorFinishConfigReplacement` | KILLED TestOperatorFinishConfigReplacement |
| `j32` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:378) | `TestOperatorCommitReservationFence` | KILLED TestOperatorCommitReservationFence |
| `j33` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:372) | `TestOperatorCommitReservationFence` | KILLED TestOperatorCommitReservationFence |
| `j35` | `cmd/ailang-worldd/operator_finish_local.go` (M1:179) | `TestOperatorFinishChecks` | CONTROL_LIVE TestOperatorFinishChecks |
| `f1-surface-bypass` | `cmd/ailang-worldd/episode_finish.go` (M1:90) | `TestOperatorFinishF1` | KILLED TestOperatorFinishF1 |
| `f1-tty-prompt` | `cmd/ailang-worldd/episode_finish_operator.go` (M1:158) | `TestOperatorFinishTTYParity` | KILLED TestOperatorFinishTTYParity |
| `proof-no-recheck` | `cmd/ailang-worldd/operator_finish_preflight.go` (M1:59) | `TestOperatorFinishNoProof` | KILLED TestOperatorFinishNoProof |
| `proof-deadline` | `cmd/ailang-worldd/operator_finish_preflight.go` (M1:44) | `TestOperatorFinishNoProof` | KILLED TestOperatorFinishNoProof |
| `commit-author-email` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:340) | `TestOperatorCommitIdentity` | KILLED TestOperatorCommitIdentity |
| `commit-match-tree` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:335) | `TestOperatorCommitIdentity` | KILLED TestOperatorCommitIdentity |
| `commit-before-reservation` | `cmd/ailang-worldd/operator_finish_local.go` (M1:151) | `TestOperatorCommitReconcile` | KILLED TestOperatorCommitReconcile |
| `proof-bytes` | `cmd/ailang-worldd/operator_finish_local.go` (M1:84) | `TestOperatorFinishPrep` | KILLED TestOperatorFinishPrep |
| `git-first-use-config` | `cmd/ailang-worldd/operator_finish_tree.go` (M1:348) | `TestOperatorFinishConfigReplacement` | KILLED TestOperatorFinishConfigReplacement |
| `git-attribution-base-delta` | `cmd/ailang-worldd/operator_finish_manifest.go` (M1:46) | `TestOperatorFinishManifestCapture` | KILLED TestOperatorFinishManifestCapture |
| `ready-before-claim` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:214) | `TestOperatorFinishReadiness` | KILLED TestOperatorFinishReadiness |
| `ready-gh-not-git` | `cmd/ailang-worldd/operator_finish_auth.go` (M1:98) | `TestOperatorFinishReadiness` | KILLED TestOperatorFinishReadiness |
| `runner-positive-claim` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:296) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-claim-cas-version` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:305) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-prep-caller` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:234) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-commit-caller` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:344) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-world-check-writer` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:338) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-claim-owner` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:306) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-start` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:41) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-resumes` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:41) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-preset-allowlist` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:41) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-remote-allowlist` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:42) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-launch-version` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:46) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-live-agent` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:49) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-health-status` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:53) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-degraded` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:53) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-build` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:54) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-report-read` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:58) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-report-episode` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:58) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-report-intents` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:58) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-recorded-finish` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:58) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-report-start` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:58) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-drift-preset` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:66) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-drift-preset-decode` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:70) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-drift-workspace` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:76) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-drift-world` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:79) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-drift-store` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:83) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-drift-profile` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:104) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-drift-pinned-check` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:96) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-sessions` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:100) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-session-expiry` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:100) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-project` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:101) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-session-file` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:108) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-token-mode` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:112) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-token` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:117) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-elig-check-sandbox` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:210) | `TestOperatorFinishReadiness` | KILLED TestOperatorFinishReadiness |
| `runner-elig-commit-identity` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:130) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-lock-busy` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:220) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-lock-ref` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:218) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-manifest-mismatch` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:278) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-manifest-protected` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:278) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-catalogue-fresh` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:293) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-heartbeat-owner` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:293) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-heartbeat-ready` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:293) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-life-quiet` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:285) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-life-pair` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:285) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-effect-clock` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:320) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-effect-inputs` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:328) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-effect-key` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:331) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-observe-after-expiry` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:203) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-old-observation` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:198) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-catalogue-approve-size` | `cmd/ailang-worldd/operator_finish_runner.go` (M1:171) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-loop-turn` | `cmd/ailang-worldd/operator_runner.go` (M1:157) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-turn-heartbeat` | `cmd/ailang-worldd/operator_runner.go` (M1:204) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-turn-sink` | `cmd/ailang-worldd/operator_runner.go` (M1:191) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-catalogue-root` | `cmd/ailang-worldd/operator_catalogue.go` (M1:39) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-catalogue-bound` | `cmd/ailang-worldd/operator_catalogue.go` (M1:45) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `runner-catalogue-preserve-start` | `cmd/ailang-worldd/operator_catalogue.go` (M1:38) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `terminal-lock-active-id` | `cmd/ailang-worldd/operator_episode_lock.go` (M1:64) | `TestOperatorTerminalExclusion` | KILLED TestOperatorTerminalExclusion |
| `route-finish-admission` | `host/daemon/surface_operator.go` (M1:31) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-start-task` | `host/daemon/surface_operator.go` (M1:209) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-finish-task` | `host/daemon/surface_operator.go` (M1:209) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-finish-health` | `host/daemon/surface_operator.go` (M1:48) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-finish-generation` | `host/daemon/surface_operator.go` (M1:256) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-finish-store` | `host/daemon/surface_operator.go` (M1:53) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-protected-and-unattributed` | `host/daemon/surface_operator.go` (M1:56) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-idempotence-kind` | `host/daemon/surface_operator.go` (M1:238) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-operation-generation` | `host/daemon/surface_operator.go` (M1:241) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-operation-generation-conflict` | `host/daemon/surface_operator.go` (M1:241) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-review-no-register` | `host/daemon/surface_operator.go` (M1:80) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-confirm-scope` | `host/daemon/surface_operator_confirm.go` (M1:96) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `route-confirm-without-surface` | `host/daemon/surface_operator_confirm.go` (M1:96) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `status-finish-projection` | `host/daemon/surface_operator_read.go` (M1:128) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `status-evidence-leak` | `host/daemon/surface_operator_read.go` (M1:135) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `status-sha` | `host/daemon/surface_operator_read.go` (M1:139) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `status-open-window` | `host/daemon/surface_operator_read.go` (M1:141) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `status-definite-failure` | `host/daemon/surface_operator_read.go` (M1:143) | `TestOperatorFinishRouteWorker` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `status-finish-auth` | `host/daemon/surface_operator_read.go` (M1:151) | `TestOperatorFinishAdmission` | UNINFORMATIVE UNDER SANDBOX real listener required |
| `js-render-export` | `host/surfacepage/app/app.js` (M1:1589) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-finish-kind` | `host/surfacepage/app/app.js` (M1:1445) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-finish-surface` | `host/surfacepage/app/app.js` (M1:1445) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-finish-ready` | `host/surfacepage/app/app.js` (M1:1446) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-finish-argv-count` | `host/surfacepage/app/app.js` (M1:1446) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-finish-empty-task` | `host/surfacepage/app/app.js` (M1:1448) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-retain-draft` | `host/surfacepage/app/app.js` (M1:1454) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-request-kind` | `host/surfacepage/app/app.js` (M1:1459) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-request-bound-link` | `host/surfacepage/app/app.js` (M1:1486) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-taskless-dom` | `host/surfacepage/app/app.js` (M1:3555) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-taskless-focus` | `host/surfacepage/app/app.js` (M1:3577) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-taskless-button` | `host/surfacepage/app/app.js` (M1:3561) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `js-taskless-post` | `host/surfacepage/app/app.js` (M1:3568) | `TestOperatorFinishBarReview` | KILLED TestOperatorFinishBarReview |
| `s9d-cli-dispatch` | `cmd/ailang-worldd/operator_service.go` (M2:318) | `TestOperatorFinishRouteWorker` | PENDING — S9d controller outside |
| `s9d-cli-sha` | `cmd/ailang-worldd/operator_finish_runner.go` (M2:384) | `TestOperatorFinishRouteWorker` | PENDING — S9d controller outside |

Compile-only route mutation preflight: 74 compiled rows. This pass does not measure behavior.

Supplemental in-memory route mutations raised their named assertions for 18 rows. Their real-listener verdicts remain outside: `runner-positive-claim`, `runner-claim-cas-version`, `runner-prep-caller`, `runner-commit-caller`, `runner-world-check-writer`, `runner-loop-turn`, `runner-elig-live-agent`, `runner-elig-token`, `runner-manifest-mismatch`, `runner-lock-busy`, `runner-lock-ref`, `runner-catalogue-root`, `route-finish-admission`, `route-finish-task`, `route-finish-generation`, `route-confirm-scope`, `status-evidence-leak`, `status-open-window`.

**Whole ledger.** PENDING — controller to measure the selected rows below with `whole-mutants.sh`. Any completed socket-free local row uses the controller’s pinned pristine fail set.

| ID | Whole test package | Verdict |
|---|---|---|
| `j08` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `j13` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `j14` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `j16` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `j19` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `j35` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `ready-before-claim` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `git-first-use-config` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `git-attribution-base-delta` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `runner-loop-turn` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `runner-positive-claim` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `runner-commit-caller` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `runner-lock-busy` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `route-finish-admission` | `host/daemon` | UNMEASURED — outside pinned-base fail-set run |
| `route-confirm-scope` | `host/daemon` | UNMEASURED — outside pinned-base fail-set run |
| `status-evidence-leak` | `cmd/ailang-worldd` | UNMEASURED — outside pinned-base fail-set run |
| `js-taskless-dom` | `host/surfacepage` | KILLED TestOperatorFinishBarReview |

## 6. Files, guards and gate scripts

The product write list comes from `status-M1` in the plan. The planner writes exactly this plan and the matching sprint JSON in the sprint worktree. Archive bases, stages, diff tools and evidence live in `$S`. Product fixtures use test-generated keys and local temporary repositories, never live presets or real launchctl. The outside admission/Chrome/launch-shaped controls still use temp stores, fake jobs and fake subprocesses. All observations of protected resources are reads.

A final integration acceptance must see the W33 fake gh label **read**, fake Git remote **read**, one authorized local check sequence and a single local commit. The isolated credential and World state, no-push sentinel, V3 row count, finish registry head, runner/daemon sinks and exact subprocess-site tests are independently observable. Scripts never install a PATH shim.

````bash accept.sh
#!/bin/bash
# TREE M1|M2 [fast|outside|outside-race|outside-1777|race|1777]
set -uo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P) || exit 2
. "$TOOLS/env.sh"
. "$TOOLS/selectors.sh"
T=$(cd "$1" && pwd -P) || exit 2
M=$2 mode=${3:-fast}
case "$M" in M1|M2) ;; *) exit 2;; esac
L=$(mktemp -d "$S/accept.$M.$mode.XXXXXX") || exit 2
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" || exit 2
cd "$T" || exit 2
echo "ACCEPT_LOGS $L"
result=0
note() { printf '%s\n' "$*" | tee -a "$L/legs.txt"; }
skip() { note "LEG $1 SKIPPED ($2)"; }
named() {
  local tag=$1 pkg=$2 selector=$3 started rc=0 parser=0
  shift 3
  started=$(date +%s)
  { date -u; uptime; } | tee "$L/$tag.load"
  python3 "$TOOLS/bounded.py" 680 go test -v -p 1 -count=1 -timeout 10m "$@" "./$pkg/" \
    -run "^($selector)$" > "$L/$tag.raw" 2>&1 || rc=$?
  echo "$rc" > "$L/$tag.rc"
  python3 "$TOOLS/failset.py" "$L/$tag.raw" "$rc" "$pkg" > "$L/$tag.fails" || parser=$?
  grep -E '^--- (PASS|SKIP|FAIL):|^FAIL|^[[:space:]]+[^ ]+\.go:' "$L/$tag.raw" || true
  if [ "$parser" -eq 77 ]; then
    note "LEG $tag UNINFORMATIVE UNDER SANDBOX ($L/$tag.raw)"
    [ "$result" -ne 0 ] || result=77
  elif grep -Fq 'panic: test timed out' "$L/$tag.raw" || [ "$parser" -ne 0 ]; then
    note "LEG $tag INSTRUMENT FAILURE ($L/$tag.raw)"
    result=2
  elif [ "$rc" -ne 0 ] || ! grep -Eq "^--- PASS: ($selector) " "$L/$tag.raw"; then
    note "LEG $tag FAILED ($L/$tag.raw)"
    result=1
  else
    note "LEG $tag OK $(( $(date +%s)-started ))s"
  fi
}
require_all() {
  local raw=$1 selector=$2 name
  while IFS= read -r name; do
    if ! grep -Fq -- "--- PASS: $name (" "$raw"; then
      note "LEG $name INSTRUMENT FAILURE missing named PASS ($raw)"
      result=2
    fi
  done < <(printf '%s\n' "$selector" | tr '|' '\n')
}
budget() {
  local pkg=$1 selector=$2 cap=$3 attempt total tag min prior=$result
  shift 3
  local totals="$L/${pkg#*/}.totals"
  : > "$totals"
  result=0
  for attempt in 1 2 3; do
    tag="${pkg#*/}-race-$attempt"
    named "$tag" "$pkg" "$selector" -race "$@"
    [ "$result" -eq 0 ] || return
    require_all "$L/$tag.raw" "$selector"
    [ "$result" -eq 0 ] || return
    awk -v pkg="$pkg" -v selector="$selector" -v read="$attempt" '
      $1=="---" && $2=="PASS:" && $3~("^(" selector ")$") {
        t=$4; gsub(/[()s]/,"",t); sum+=t; printf "RACE_TEST %s reading=%s %s %.2f\n",pkg,read,$3,t
      }
      END { printf "%.2f\n",sum }' "$L/$tag.raw" > "$L/$tag.reading"
    sed '$d' "$L/$tag.reading"
    total=$(tail -1 "$L/$tag.reading")
    echo "$total" >> "$totals"
    note "RACE_READING $pkg reading=$attempt total=${total}s load=$L/$tag.load"
    if awk -v reading="$total" -v cap="$cap" 'BEGIN { exit !(reading<=cap) }'; then break; fi
  done
  min=$(sort -n "$totals" | head -1)
  note "RACE_BUDGET $pkg min=${min}s cap=${cap}s readings=$(paste -sd, "$totals")"
  awk -v reading="$min" -v cap="$cap" 'BEGIN { exit !(reading<=cap) }' || result=1
  [ "$result" -ne 0 ] || result=$prior
}
case "$mode" in
  fast|outside)
    unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
    named cmd-core cmd/ailang-worldd "$CMD_ADDED|$CMD_FENCES|$CMD_JOURNAL|$CMD_GIT"
    [ "$result" -ne 0 ] || require_all "$L/cmd-core.raw" "$CMD_ADDED"
    "$WORLD_EXEC_NODE" --check host/surfacepage/app/app.js || result=2
    named app-stems host/surfacepage "$APP_STEMS"
    [ "$result" -ne 0 ] || require_all "$L/app-stems.raw" TestOperatorFinishBarReview
    if [ "$mode" = outside ]; then
      named cmd-listener cmd/ailang-worldd "$CMD_LISTENER"
      named cmd-stems cmd/ailang-worldd "$CMD_STEMS"
      named daemon-stems host/daemon "$DAEMON_ADDED|$DAEMON_STEMS"
      if [ "$result" -eq 0 ]; then
        require_all "$L/cmd-listener.raw" "$CMD_LISTENER"
        require_all "$L/daemon-stems.raw" "$DAEMON_ADDED"
      fi
    else
      skip cmd-listener 'real listener outside sandbox'
      skip cmd-stems 'launchd-shaped fixtures, listener, Chrome and docs controls outside'
      skip daemon-stems 'real listener outside sandbox'
    fi
    ;;
  race)
    named cmd-race-local cmd/ailang-worldd "$CMD_ADDED" -race -skip "$CMD_SKIP"
    named app-race-local host/surfacepage TestOperatorFinishBarReview -race
    skip finish-race-budget 'runner and daemon listeners required for complete readings'
    ;;
  outside-race)
    unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
    budget cmd/ailang-worldd "$CMD_ADDED|$CMD_LISTENER" 15 -skip "$CMD_SKIP"
    budget host/daemon "$DAEMON_ADDED" 12
    named app-race host/surfacepage TestOperatorFinishBarReview -race
    ;;
  1777|outside-1777)
    unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
    mkdir -m 1777 "$L/t1777" && mkdir -m 0700 "$L/t1777/inside" || exit 2
    if [ "$mode" = outside-1777 ]; then
      export TMPDIR="$L/t1777/inside"
      named cmd-1777 cmd/ailang-worldd "$CMD_LISTENER|$CMD_ADDED|TestOperatorFinishToolPins|TestOperatorFinishAuthInstall|TestRunnerPathRefusalClasses"
      named daemon-1777 host/daemon "$DAEMON_ADDED"
    else
      export TMPDIR="$L/t1777/inside"
      named cmd-1777 cmd/ailang-worldd "$CMD_ADDED|TestOperatorFinishToolPins|TestRunnerPathRefusalClasses"
      skip listener-1777 'real listener outside sandbox'
    fi
    ;;
  *) exit 2;;
esac
[ "$result" -eq 0 ] || exit "$result"
echo "ACCEPT_OK $M $mode"
````

````bash gates.sh
#!/bin/bash
# TREE LOGDIR M1|M2 local|outside|base-extra. Record all raw legs before attribution.
set -uo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P) || exit 2
. "$TOOLS/env.sh"
. "$TOOLS/selectors.sh"
T=$(cd "$1" && pwd -P) || exit 2
mkdir "$2" || exit 2
L=$(cd "$2" && pwd -P) || exit 2
M=$3 mode=${4:-local}
case "$M:$mode" in M1:local|M1:outside|M1:base-extra|M2:local|M2:outside) ;; *) exit 2;; esac
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
cd "$T" || exit 2
issue=0
pending=0
note() { printf '%s\n' "$*" | tee -a "$L/legs.txt"; }
skip() { note "LEG $1 SKIPPED ($2)"; }
bound() { python3 "$TOOLS/bounded.py" "$@"; }
instrument() { grep -Eiq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$1"; }
denied() { grep -Eiq '(bind|listen tcp|launch.{0,30}chrome).*(operation not permitted|permission denied)|sandbox: deny' "$1"; }
run() {
  local tag=$1 started rc=0
  shift
  started=$(date +%s)
  mkdir "$L/$tag.state" "$L/$tag.cache" || return 2
  { date -u; uptime; } > "$L/$tag.load"
  ( export AILANG_STATE_DIR="$L/$tag.state" AILANG_CACHE_DIR="$L/$tag.cache"; "$@" ) > "$L/$tag.raw" 2>&1 || rc=$?
  echo "$rc" > "$L/$tag.rc"
  if instrument "$L/$tag.raw"; then issue=1; note "LEG $tag INSTRUMENT FAILURE rc=$rc"
  elif [ "$mode" = local ] && { denied "$L/$tag.raw" || [ "$rc" -eq 77 ]; }; then
    pending=1; note "LEG $tag UNINFORMATIVE UNDER SANDBOX rc=$rc"
  elif [ "$rc" -ne 0 ]; then issue=1; note "LEG $tag FAILED rc=$rc"
  else note "LEG $tag OK $(( $(date +%s)-started ))s"; fi
}
record() {
  local tag=$1 pkg=$2 limit=$3 seconds=$4 started rc=0 parser=0
  shift 4
  [ "$tag" != cmd-race ] || set -- "$@" -skip "$CMD_SKIP"
  started=$(date +%s)
  mkdir "$L/$tag.state" "$L/$tag.cache" || return 2
  { date -u; uptime; } > "$L/$tag.load"
  ( export AILANG_STATE_DIR="$L/$tag.state" AILANG_CACHE_DIR="$L/$tag.cache"
    bound "$seconds" go test -count=1 -timeout "$limit" "$@" "./$pkg/" ) > "$L/$tag.raw" 2>&1 || rc=$?
  echo "$rc" > "$L/$tag.rc"
  python3 "$TOOLS/failset.py" "$L/$tag.raw" "$rc" "$pkg" > "$L/$tag.fails" 2> "$L/$tag.parser" || parser=$?
  cat "$L/$tag.parser"
  grep -E '^(ok|FAIL)[[:space:]]+github.com/' "$L/$tag.raw" >> "$L/pkg-times.txt" || true
  if [ "$parser" -eq 77 ] && [ "$mode" = local ]; then
    pending=1; note "LEG $tag UNINFORMATIVE UNDER SANDBOX rc=$rc"
  elif [ "$parser" -ne 0 ] || instrument "$L/$tag.raw"; then
    issue=1; note "LEG $tag INSTRUMENT FAILURE rc=$rc parser=$parser"
  else
    cp "$L/$tag.fails" "$L/$tag.verified"
    note "LEG $tag recorded rc=$rc $(( $(date +%s)-started ))s"
  fi
}
compare() {
  local tag=$1 pkg=$2 source="$STATE/base/$1"
  if [ "$tag" = cmd-race ]; then source="$S/base-extra/$tag"; fi
  if [ ! -f "$L/$tag.verified" ]; then skip "compare-$tag" 'stage evidence unavailable'; return; fi
  if [ ! -f "$source.rc" ] || [ ! -f "$source.raw" ] ||
     ! python3 "$TOOLS/failset.py" "$source.raw" "$(cat "$source.rc")" "$pkg" > "$L/$tag.base"; then
    pending=1; skip "compare-$tag" 'comparable pristine evidence unavailable'; return
  fi
  LC_ALL=C sort -u "$L/$tag.base" > "$L/$tag.base.sorted"
  LC_ALL=C sort -u "$L/$tag.verified" > "$L/$tag.head.sorted"
  LC_ALL=C comm -13 "$L/$tag.base.sorted" "$L/$tag.head.sorted" > "$L/$tag.new"
  if [ -s "$L/$tag.new" ]; then issue=1; note "LEG compare-$tag NEW FAIL SET"; cat "$L/$tag.new"; else note "LEG compare-$tag OK"; fi
}
gofiles() {
  local p list="$TOOLS/changed-$M"
  [ "$M" = M1 ] || list="$TOOLS/changed-M2-cumulative"
  while IFS= read -r p; do
    case "$p" in *.go) [ -z "$(gofmt -l "$p")" ] || { echo "$p"; return 1; };; esac
  done < "$list"
}
inventory() {
  local pkg=$1 selector=$2
  bound 640 go test -v -count=1 -timeout 10m "./$pkg/" -run "^($selector)$"
}
if [ "$mode" = base-extra ]; then
  if ! diff -qr "$S/base-tree" "$T" > "$L/base-compare.raw"; then
    issue=1; note "LEG pristine-cmd-race INSTRUMENT FAILURE tree mismatch"
  else
    record cmd-race cmd/ailang-worldd 45m 2900 -race
  fi
else
  packages=(./cmd/ailang-worldd ./host/daemon ./host/surfacepage)
  run gofmt gofiles
  run build bound 850 go build ./...
  run vet bound 850 go vet "${packages[@]}"
  run linux-vet env GOOS=linux python3 "$TOOLS/bounded.py" 850 go vet "${packages[@]}"
  run test-compile bound 1100 go test -run '^$' -timeout 45m ./...
  run inventory-store inventory host/store "$STORE_GUARDS"
  run inventory-cli inventory cmd/ailang-worldd "$CLI_GUARDS"
  run inventory-daemon inventory host/daemon "$DAEMON_GUARDS"
  run inventory-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
  run subprocess-sink inventory host/verifygate TestSubprocessSinkGate
  run subprocess-sites inventory host/broker "$BROKER_GUARDS"
  run js-syntax "$WORLD_EXEC_NODE" --check host/surfacepage/app/app.js
  run verify-ail bound 1000 python3 "$TOOLS/verify-scratch.py" "$T" "$L" verify_ail ""
  run caller-census bash "$TOOLS/callers.sh" "$S/base-tree" "$T" "$M"

  if [ "$mode" = local ]; then
    run accept bash "$TOOLS/accept.sh" "$T" "$M" fast
    run accept-race bash "$TOOLS/accept.sh" "$T" "$M" race
    run accept-1777 bash "$TOOLS/accept.sh" "$T" "$M" 1777
    for tag in cmd-plain cmd-race daemon-plain daemon-race verifygate-plain broker-plain; do
      skip "$tag" 'outside.sh whole'
    done
  else
    run accept bash "$TOOLS/accept.sh" "$T" "$M" outside
    run accept-race bash "$TOOLS/accept.sh" "$T" "$M" outside-race
    run accept-1777 bash "$TOOLS/accept.sh" "$T" "$M" outside-1777
    record cmd-plain cmd/ailang-worldd 45m 2900
    record cmd-race cmd/ailang-worldd 45m 2900 -race
    record daemon-plain host/daemon 45m 2900
    record daemon-race host/daemon 45m 2900 -race
    record verifygate-plain host/verifygate 20m 1700
    record broker-plain host/broker 20m 1700
  fi
  record surfacepage-plain host/surfacepage 25m 1700
  record surfacepage-race host/surfacepage 25m 1700 -race
  if ! grep -Fxq "ALL_DONE for $B" "$STATE/baseline.done" 2>/dev/null; then
    pending=1; skip baseline-completion 'controller recording pristine legs'
  fi
  for tag in cmd-plain cmd-race daemon-plain daemon-race verifygate-plain broker-plain surfacepage-plain surfacepage-race; do
    case "$tag" in cmd-*) pkg=cmd/ailang-worldd;; *) pkg=host/${tag%-*};; esac
    compare "$tag" "$pkg"
  done
fi
note "GATE_RESULT issue=$issue pending=$pending"
note ALL_DONE
printf '%s\n' ALL_DONE > "$L/ALL_DONE"
[ "$issue" -eq 0 ] || exit 1
[ "$pending" -eq 0 ] || exit 77
````

````bash snapshot.sh
#!/bin/bash
# snapshot.sh TREE M1|M2; SNAP_STATUS_LIST/ SNAP_PARENT are rehearsal-only.
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1|M2) ;; *) exit 2;; esac
B=91befdab7fd4081b4dc127cc96cc4eeb1723cffa
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter450/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s9c-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s9c.json"
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

````bash env.sh
STATE="$HOME/.ailang/state/world-iter450"
S="$STATE/plan-scratch"
B=91befdab7fd4081b4dc127cc96cc4eeb1723cffa
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export GOPROXY=off GIT_OPTIONAL_LOCKS=0
unset GOFLAGS WORLD_PASSKEY_RECORD_DIR
````

````bash selectors.sh
# Exact selection fragments for iteration 450. Selectors may be combined with |.
CMD_SKIP='^(TestCLIAgainstSeToolsDaemon|TestCLINoGenesisIsNoWorldHead|TestDemoSeedApprovalPacket|TestDemoSeedEndToEnd|TestDemoSeedManifestReferences|TestDemoSeedProofStage|TestDemoSeedReadTestCheck|TestDemoSeedRejectsToolReportRefusals|TestDemoSeedRosterKindMismatchRefuses|TestDemoSeedWhyNamedLinks|TestMarketingStillRoutes)$'
CMD_ADDED='TestOperatorFinishReadiness|TestOperatorFinishResumeEvidence|TestOperatorCommitFailureInputRecheck|TestOperatorCommitResumeReserved'
CMD_LISTENER='TestOperatorFinishRouteWorker'
DAEMON_ADDED='TestOperatorFinishAdmission'
CMD_FENCES='TestOperatorFinishF1|TestOperatorFinishTTYParity|TestOperatorFinishNoProof|TestOperatorFinishProofScope|TestOperatorFinishAuth|TestOperatorFinishToolPins|TestOperatorFlagParity|TestOperatorTerminalExclusion'
CMD_JOURNAL='TestOperatorFinishPrep|TestOperatorFinishPrepExpiry|TestOperatorFinishStale|TestOperatorFinishChecks|TestOperatorFinishUnsandboxed|TestOperatorCommitFailed|TestOperatorCommitReconcile|TestOperatorFinishJournalOnce|TestOperatorFinishManifestCapture|TestOperatorFinishWorldEvidence'
CMD_GIT='TestOperatorFinishWorktreeFence|TestOperatorCommitIdentity|TestOperatorFinishCaptureBounds|TestOperatorFinishReadReplacement|TestOperatorCommitReservationFence|TestOperatorFinishConfigReplacement'
CMD_ACTIVATION='TestOperatorFirstManagedEpisodeProcess|TestOperatorPreparationRechecksHeartbeatAfterWorktree|TestOperatorStartMintRecoveryRemedies|TestOperatorCurrentOwnerProcessRecovery|TestOperatorAttendedStartDisposition|TestOperatorActivationDrainRefusals|TestOperatorRevokeStatRefusal'
CMD_TREE='TestOperatorStartOnce|TestOperatorStartCollisions|TestOperatorWorktreeWindow|TestOperatorWorktreeDispositionCAS|TestOperatorWorktreeSurvivors|TestOperatorWorktreeRecoveryGuards|TestOperatorWorktreeReadbackGuards|TestOperatorReadbackErrorIsNeverCreated|TestOperatorWorktreeLockedReread|TestOperatorBorrowedWorktreeGitProcess'
CMD_LAUNCH='TestOperatorSpawnRefused|TestOperatorLaunchDispatchRechecks|TestOperatorLaunchCrashWindow|TestOperatorLaunchExit|TestOperatorStartBrowser|TestOperatorLaunchWorker|TestOperatorStartedRequiresBirth|TestOperatorSupervisorCap|TestOperatorFakePresetEligibility'
CMD_SERVICE='TestOperatorServiceFence|TestOperatorServiceOwnership|TestOperatorServiceConfig|TestOperatorServicePair|TestOperatorServiceStartup|TestOperatorServiceUpgradePairFence|TestOperatorServiceCancelConfirmation|TestOperatorServiceWriteScope|TestOperatorServiceEnvironment|TestRunnerPathRefusalClasses|TestOperatorRunnerKernelLock|TestOperatorStatusDisconnected|TestOperatorFinishAuthInstall'
CMD_DOC='TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartSection18FlagsMatchTheCLI|TestQuickstartSection18CardMatchesThePage|TestQuickstartSection19MatchesTheBar|TestQuickstartSection20MatchesTheFinishFlag|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartServiceSectionMatchesTheCLIHelp|TestQuickstartRunnerMintRecovery|TestProfileInitQuickstartLineMatchesTheCLI|TestDemoSeedHelpPreparation'
CMD_STEMS="$CMD_ACTIVATION|$CMD_TREE|$CMD_LAUNCH|$CMD_SERVICE|$CMD_DOC|TestOperatorFilePublicationAdmission"
DAEMON_STEMS='TestOperatorAdmissionHTTP|TestOperatorReviewRefusesForgedCatalogue|TestOperatorReplayConflicts|TestOperatorLinksProjection|TestOperatorConfirmHTTP|TestOperatorConfirmProjection|TestOperatorConfirmationAdmissionNeighbours|TestOperatorCardAssets|TestSurfaceVerbRouteRegistered|TestSurfaceVerbBoundsAndErrors|TestSurfaceVerbIsReadOnly|TestSurfaceVerbWhyIsReadOnly|TestSurfaceVerbOutcomes|TestSurfaceAppRoutes|TestRequestsJSONConfirmationsGolden|TestMCPQuickstartPayloadsVerbatim|TestAGUIQuickstartControl|TestQuickstartDecisionPacketSection|TestWorkbenchLiveQuickstartControl|TestSeToolsQuickstartPayloadsVerbatim|TestExecStartupRefusalTable'
APP_STEMS='TestOperatorBarStartReview|TestOperatorOtherVerbParity|TestOperatorLiteralTask|TestOperatorDuplicateTap|TestOperatorReconnect|TestOperatorPanelLatest|TestOperatorPanelNetwork|TestOperatorLinkAllowlist|TestOperatorFinishBarReview|TestOperatorCardStatusDisconnected|TestOperatorCardLinkAllowlist|TestSurfaceAppMatchesWorkingCopy|TestSurfaceAppShellAssembly|TestSurfaceAppPure|TestSurfaceAppVerbs|TestSurfaceAppFollow|TestSurfaceAppTextIsSelectable|TestSurfaceScriptAssetSeams|TestAddressHintsArePinned|TestNeedsYouListsConfirmations|TestNeedsYouConfirmationEmbeddedNewline|TestNeedsYouConfirmationsNonArray|TestVerbDecidePrefixBoundary|TestNeedsYouListsDecisions'
STORE_GUARDS='TestProductionContextRoots|TestProductionGoSurface|TestCredentialsMethodsRequireADeadline|TestNoNewDeadlineFreeStoreReads|TestStoreDeadlineGuardDynamic|TestOperatorFinishCompletedLaunch|TestOperatorClaimRace|TestOperatorAdmissionBoundaries|TestOperatorEpisodeLockContendsWithoutWaiting'
CLI_GUARDS='TestCLICredentialWriterSites|TestSessionRootBudgets|TestOperatorProductionDeadlines'
DAEMON_GUARDS='TestDaemonCredentialWriterSites|TestDaemonDependencyAllowlist|TestDaemonCallerInventoryAliasControl'
BROKER_GUARDS='TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential|TestExecProtectedSettingsThroughHandler|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract'
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
case "$3" in M1|M2) census "$2" "$3" 2 1;; *) exit 2;; esac
echo CALLER_CENSUS_OK
````

````bash outside.sh
#!/bin/bash
# TREE M1|M2 [focused|whole|chrome]. Run outside the sandbox on the planner stage first.
set -uo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P) || exit 2
. "$TOOLS/env.sh"
. "$TOOLS/selectors.sh"
T=$(cd "$1" && pwd -P) || exit 2
M=$2 mode=${3:-focused}
L=$(mktemp -d "$S/outside.$M.$mode.XXXXXX") || exit 2
echo "OUTSIDE_LOGS $L"
rc=0
run() {
  local tag=$1 status=0
  shift
  (date -u; uptime) > "$L/$tag.load"
  "$@" > "$L/$tag.raw" 2>&1 || status=$?
  cat "$L/$tag.raw"
  echo "LEG $tag rc=$status"
  [ "$status" -eq 0 ] || rc=$status
}
case "$mode" in
  focused)
    for stage in outside outside-race outside-1777; do
      run "$stage" bash "$TOOLS/accept.sh" "$T" "$M" "$stage"
    done
    ;;
  whole)
    run whole bash "$TOOLS/gates.sh" "$T" "$L/gates" "$M" outside
    ;;
  chrome)
    export WORLD_CHROME="${WORLD_CHROME:-/Applications/Google Chrome.app/Contents/MacOS/Google Chrome}"
    export WORLD_REQUIRE_PASSKEY_BROWSER=1 AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
    mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" || exit 2
    cd "$T" || exit 2
    run chrome python3 "$TOOLS/bounded.py" 310 go test -v -count=1 -timeout 5m ./cmd/ailang-worldd/ -run '^TestOperatorStartWalkthroughInChrome$'
    python3 "$TOOLS/failset.py" "$L/chrome.raw" "$rc" cmd/ailang-worldd > "$L/chrome.fails" || rc=$?
    grep -Fq -- '--- PASS: TestOperatorStartWalkthroughInChrome (' "$L/chrome.raw" || [ "$rc" -ne 0 ] || rc=2
    ;;
  *) exit 2;;
esac
echo ALL_DONE
exit "$rc"
````

````bash commit-guards.sh
#!/bin/bash
# TREE LOGDIR, outside sandbox, after compatible pristine logs exist.
set -uo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P) || exit 2
. "$TOOLS/env.sh"
T=$(cd "$1" && pwd -P) || exit 2
mkdir "$2" || exit 2
L=$(cd "$2" && pwd -P) || exit 2
cd "$T" || exit 2
grep -Fxq "ALL_DONE for $B" "$STATE/baseline.done" || exit 77
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" || exit 2
{ date -u; uptime; } > "$L/guard.load"
rc=0
python3 "$TOOLS/bounded.py" 2650 go test -count=1 -timeout 20m ./host/verifygate/ ./host/broker/ > "$L/guard.raw" 2>&1 || rc=$?
echo "$rc" > "$L/guard.rc"
python3 "$TOOLS/failset.py" "$L/guard.raw" "$rc" host/verifygate host/broker > "$L/guard.fails" || exit 2
for pkg in verifygate broker; do
  base="$STATE/base/$pkg-plain"
  python3 "$TOOLS/failset.py" "$base.raw" "$(cat "$base.rc")" "host/$pkg" > "$L/$pkg.base" || exit 2
done
cat "$L/verifygate.base" "$L/broker.base" | LC_ALL=C sort -u > "$L/base.fails"
LC_ALL=C sort -u "$L/guard.fails" > "$L/head.fails"
LC_ALL=C comm -13 "$L/base.fails" "$L/head.fails" > "$L/new.fails"
if [ -s "$L/new.fails" ]; then cat "$L/new.fails"; exit 1; fi
echo COMMIT_GUARDS_OK
````

````bash named-mutants.sh
#!/bin/bash
# TREE [LOGDIR] [outside|local] [M1|M2] [IDFILE].
set -euo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P)
. "$TOOLS/env.sh"
export STATE
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
T=$(cd "$1" && pwd -P)
args=()
if [ "${3:-outside}" = local ]; then args+=(--local); fi
python3 "$TOOLS/mutate.py" "$T" named "${2:-$S/named-mutants-$(date +%s)}" "${5:-$TOOLS/named-mutant-ids.txt}" --stage "${4:-M1}" "${args[@]}"
````

````bash whole-mutants.sh
#!/bin/bash
# TREE [IDFILE] [M1|M2]. Real listeners and comparable pristine TEXT logs required.
set -euo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P)
. "$TOOLS/env.sh"
export STATE
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
T=$(cd "$1" && pwd -P)
L=$(mktemp -d "$S/whole-mutants.XXXXXX")
echo "WHOLE_MUTANT_LOGS $L"
grep -Fxq "ALL_DONE for $B" "$STATE/baseline.done" || exit 77
python3 "$TOOLS/mutate.py" "$T" whole "$L/results" "${2:-$TOOLS/whole-mutant-ids.txt}" --stage "${3:-M1}"
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
"""Parse package-qualified Go TEXT fail sets. Reject incomplete evidence."""
from pathlib import Path
import re
import sys

prefix = "github.com/sunholo-data/ailang-world/"
invalid = re.compile(r"panic: test timed out|BOUND EXCEEDED|\[(?:build|setup) failed\]", re.I)
deny = re.compile(r"(?:bind|listen tcp|launch.{0,30}chrome).*(?:operation not permitted|permission denied)|sandbox: deny", re.I)

def parse(text, rc):
    if invalid.search(text):
        raise ValueError("truncated or uncompiled test package")
    if deny.search(text):
        raise PermissionError("sandbox denial")
    pending, fails, results = [], set(), set()
    for line in text.splitlines():
        failed = re.match(r"^--- FAIL: ([^/\s]+) \(", line)
        if failed:
            pending.append(failed[1])
        pkg = re.match(r"^(ok|FAIL)\s+" + re.escape(prefix) + r"(\S+)", line)
        if pkg:
            results.add(pkg[2])
            if (pkg[1] == "FAIL") != bool(pending):
                raise ValueError("package summary and test failures disagree: " + pkg[2])
            fails.update((pkg[2], test) for test in pending)
            pending.clear()
    if pending or not results or rc != 0 and not fails or rc == 0 and fails:
        raise ValueError("missing package summary or unattributed text failure")
    return fails, results

def main():
    must = {("host/store", "TestParserPositive")}
    got, _ = parse("--- FAIL: TestParserPositive (0.00s)\nFAIL\t" + prefix + "host/store\t0.001s\n", 1)
    if got != must:
        raise ValueError("firing control did not fire")
    print("FAILSET_CONTROL host/store TestParserPositive", file=sys.stderr)
    text, rc, *expected = sys.argv[1:]
    fails, results = parse(Path(text).read_text(errors="replace"), int(rc))
    if not set(expected) <= results:
        raise ValueError("missing whole package or named gate: " + ", ".join(sorted(set(expected)-results)))
    for pkg, name in sorted(fails):
        print(pkg, name)

if __name__ == "__main__":
    try:
        main()
    except PermissionError as error:
        print("UNINFORMATIVE UNDER SANDBOX:", error, file=sys.stderr)
        sys.exit(77)
    except ValueError as error:
        sys.exit("INSTRUMENT FAILURE: " + str(error))
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
"""One exact mutation per overlay, named test or whole package. No checkout writes."""
import argparse
import json
import os
from pathlib import Path
import re
import subprocess
import sys

p = argparse.ArgumentParser(__doc__)
p.add_argument("tree", type=Path)
p.add_argument("mode", choices=["named", "whole", "compile"])
p.add_argument("log", type=Path)
p.add_argument("ids", type=Path)
p.add_argument("--stage", default="M1")
p.add_argument("--local", action="store_true", help="record real-listener cases for outside")
a = p.parse_args()
tools = Path(__file__).resolve().parent
tree = a.tree.resolve()
a.log.mkdir()
all_rows = json.loads((tools/"mutants.json").read_text())
patterns = list(dict.fromkeys(a.ids.read_text().split()))
limits = {"cmd/ailang-worldd": ("45m", 2900), "host/daemon": ("45m", 2900),
          "host/surfacepage": ("25m", 1700)}
stage_to_base = {"cmd/ailang-worldd":"cmd-plain", "host/daemon":"daemon-plain", "host/surfacepage":"surfacepage-plain"}
base_dir = Path(os.environ["STATE"]) / "base"

def parsed(log, result, pkg, to):
    proc = subprocess.run([sys.executable,str(tools/"failset.py"),str(log),str(result),pkg],
                          capture_output=True,text=True)
    to.write_text(proc.stderr)
    print(proc.stderr.rstrip(),flush=True)
    return proc.returncode, set(proc.stdout.splitlines())

def trial(id, row):
    pkg, test, expect = row["package"], row["test"], row["control"]
    if row["stage"] != a.stage:
        return "SKIPPED other milestone"
    if a.local and row["access"] == "outside":
        return "UNINFORMATIVE UNDER SANDBOX real listener required"
    spot = a.log/id
    spot.mkdir()
    path = tree / row["path"]
    src = path.read_text()
    hits = list(re.finditer(re.escape(row["old"]),src))
    if not hits or row["occurrence"] >= len(hits):
        return "INSTRUMENT_FAILURE anchor drift"
    hit = hits[row["occurrence"]]
    text = src[:hit.start()] + row["new"] + src[hit.end():]
    file = spot/path.name
    file.write_text(text)
    overlay = spot/"overlay.json"
    overlay.write_text(json.dumps({"Replace":{str(path):str(file)}}))
    if a.mode == "whole":
        name=stage_to_base[pkg]
        base = base_dir/(name+".raw")
        rc_file = base_dir/(name+".rc")
        if not base.exists() or not rc_file.exists():
            return "UNMEASURED comparable whole base not recorded"
        rc, base_fails = parsed(base,int(rc_file.read_text()),pkg,spot/"base.parser")
        if rc:
            return "INSTRUMENT_FAILURE unusable baseline"
        limit, seconds = limits[pkg]
        selector = []
    else:
        base_fails = set()
        limit, seconds = "8m", 535
        selector = ["-run", "^$" if a.mode=="compile" else "^"+re.escape(test)+"$"]
    env = dict(os.environ, AILANG_STATE_DIR=str(spot/"state"), AILANG_CACHE_DIR=str(spot/"cache"))
    for k in ("AILANG_STATE_DIR","AILANG_CACHE_DIR"):
        Path(env[k]).mkdir()
    with (spot/"load").open("w") as fd:
        subprocess.run(["uptime"],stdout=fd)
    raw = spot/"go.raw"
    command = ["go","test","-count=1","-v","-timeout",limit,"-overlay",str(overlay),"./"+pkg+"/"] + selector
    (spot/"command.json").write_text(json.dumps(command))
    with raw.open("w") as fd:
        result = subprocess.run([sys.executable,str(tools/"bounded.py"),str(seconds),*command],
                                cwd=tree, env=env, stdout=fd, stderr=subprocess.STDOUT).returncode
    (spot/"rc").write_text(str(result)+"\n")
    rc, found = parsed(raw,result,pkg,spot/"parser")
    if rc == 77:
        return "UNINFORMATIVE UNDER SANDBOX denial"
    if rc:
        return "INSTRUMENT_FAILURE (compile/timeout/attribution)"
    if a.mode == "compile":
        return "COMPILED" if result == 0 else "INSTRUMENT_FAILURE compile"
    (spot/"base.fails").write_text("".join(x+"\n" for x in sorted(base_fails)))
    (spot/"head.fails").write_text("".join(x+"\n" for x in sorted(found)))
    new = found - base_fails
    if pkg + " " + test in new and re.search(r"^--- FAIL: "+re.escape(test)+r" \(",raw.read_text(errors="replace"),re.M):
        return "CONTROL_FAILED unexpected named fail" if expect=="live" else "KILLED " + test
    if new:
        return "OTHER_FAIL " + ", ".join(sorted(new))
    if expect=="live":
        return "CONTROL_LIVE " + test
    return "SURVIVED " + test

has_issue = False
has_pending = False
with (a.log/"ledger").open("w",buffering=1) as ledger:
    for id in patterns:
        row = all_rows[id]
        result = trial(id,row)
        print(id,result,flush=True)
        ledger.write(f"{id}\t{result}\n")
        has_issue |= result.startswith(("SURVIVED","INSTRUMENT_FAILURE","OTHER_FAIL","CONTROL_FAILED"))
        has_pending |= result.startswith(("UNINFORMATIVE","UNMEASURED"))
    ledger.write("ALL_DONE\n")
print("ALL_DONE")
sys.exit(1 if has_issue else 77 if has_pending else 0)
````

````json mutants.json
{
  "j01": {
    "path": "cmd/ailang-worldd/operator_finish_preflight.go",
    "old": "frozen.Kind != \"finish\" || frozen.ID != cas.ID",
    "new": "frozen.ID != cas.ID",
    "test": "TestOperatorFinishProofScope",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 26,
    "stage": "M1",
    "control": "die"
  },
  "j02": {
    "path": "cmd/ailang-worldd/operator_finish_preflight.go",
    "old": "if !frozen.Options.ConfirmOnSurface {",
    "new": "if false {",
    "test": "TestOperatorFinishProofScope",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 29,
    "stage": "M1",
    "control": "die"
  },
  "j04": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if !w.prep.claimed() {",
    "new": "if false {",
    "test": "TestOperatorFinishPrep",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 80,
    "stage": "M1",
    "control": "die"
  },
  "j05": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if check.Unsandboxed {\n\t\t\treturn errOperatorFinishTerminal",
    "new": "if check.Unsandboxed && false {\n\t\t\treturn errOperatorFinishTerminal",
    "test": "TestOperatorFinishUnsandboxed",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 88,
    "stage": "M1",
    "control": "die"
  },
  "j06": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "len(m.Unattributed) != 0 || len(m.Protected) != 0 ||",
    "new": "len(m.Unattributed) != 0 ||",
    "test": "TestOperatorFinishPrep",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 93,
    "stage": "M1",
    "control": "die"
  },
  "j07": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if j.OpenCheck || len(j.Checks)",
    "new": "if len(j.Checks)",
    "test": "TestOperatorFinishChecks",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 109,
    "stage": "M1",
    "control": "die"
  },
  "j08": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "old.Name != m.Checks[i].Name || !w.evidence(ctx, old)",
    "new": "old.Name != m.Checks[i].Name",
    "test": "TestOperatorFinishResumeEvidence",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 113,
    "stage": "M1",
    "control": "die"
  },
  "j10": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "\tif !passed {\n\t\treturn errOperatorFinishCheck\n\t}\n\treturn nil",
    "new": "\treturn nil",
    "test": "TestOperatorFinishChecks",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 204,
    "stage": "M1",
    "control": "die"
  },
  "j13": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if _, err := w.inputs(budget); err == nil {",
    "new": "if _, err := w.inputs(budget); err == nil || true {",
    "test": "TestOperatorCommitFailureInputRecheck",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 161,
    "stage": "M1",
    "control": "die"
  },
  "j14": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "return w.reconcileCommit(j.CommitKey)",
    "new": "return errOperatorFinishUncertain",
    "test": "TestOperatorCommitResumeReserved",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 101,
    "stage": "M1",
    "control": "die"
  },
  "j16": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if _, err := w.git.config(ctx, w.report.Start); err != nil {\n\t\treturn errOperatorFinishUncertain\n\t}\n\tsha",
    "new": "sha",
    "test": "TestOperatorCommitResumeReserved",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 217,
    "stage": "M1",
    "control": "die"
  },
  "j19": {
    "path": "cmd/ailang-worldd/operator_finish_manifest.go",
    "old": "object.SemanticID != taskOutcomeSemanticID || ",
    "new": "",
    "test": "TestOperatorFinishWorldEvidence",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 102,
    "stage": "M1",
    "control": "die"
  },
  "j20": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "|| len(strings.Fields(authorTime)) != 2 {",
    "new": "|| len(strings.Fields(authorTime)) < 0 {",
    "test": "TestOperatorCommitIdentity",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 340,
    "stage": "M1",
    "control": "die"
  },
  "j26": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "err != nil || strings.TrimSpace(branch) != spec.Snapshot.Branch {",
    "new": "err != nil || strings.TrimSpace(branch) == \"\\x00\" {",
    "test": "TestOperatorCommitIdentity",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 313,
    "stage": "M1",
    "control": "die"
  },
  "j22": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "return ok && committer+\">\" == \"committer \"+spec.Committer && len(strings.Fields(timeFields)) == 2",
    "new": "return ok && committer != \"\" && len(strings.Fields(timeFields)) == 2",
    "test": "TestOperatorCommitIdentity",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 344,
    "stage": "M1",
    "control": "die"
  },
  "j23": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "lines[1] != \"parent \"+spec.Snapshot.Head || message",
    "new": "message",
    "test": "TestOperatorCommitIdentity",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 336,
    "stage": "M1",
    "control": "die"
  },
  "j24": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "message != spec.Message {",
    "new": "!strings.HasPrefix(message, spec.Message[:10]) {",
    "test": "TestOperatorCommitIdentity",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 336,
    "stage": "M1",
    "control": "die"
  },
  "j27": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "err != nil || !bytes.Equal(current, data)",
    "new": "err != nil || !bytes.Equal(current, current)",
    "test": "TestOperatorFinishConfigReplacement",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 160,
    "stage": "M1",
    "control": "die"
  },
  "j32": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "\"update-ref\", spec.Snapshot.Branch, head, spec.Snapshot.Head)",
    "new": "\"update-ref\", spec.Snapshot.Branch, head)",
    "test": "TestOperatorCommitReservationFence",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 378,
    "stage": "M1",
    "control": "die"
  },
  "j33": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "if err != nil || !reflect.DeepEqual(current, spec.Snapshot) {\n\t\treturn errOperatorFinishReview",
    "new": "if false && !reflect.DeepEqual(current, spec.Snapshot) {\n\t\treturn errOperatorFinishReview",
    "test": "TestOperatorCommitReservationFence",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 372,
    "stage": "M1",
    "control": "die"
  },
  "j35": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "// No trustworthy terminal readback.",
    "new": "// No trustworthy readback.",
    "test": "TestOperatorFinishChecks",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 179,
    "stage": "M1",
    "control": "live"
  },
  "f1-surface-bypass": {
    "path": "cmd/ailang-worldd/episode_finish.go",
    "old": "\t\tif err != nil {\n\t\t\treturn fail(\"refusing: no controlling terminal\"",
    "new": "\t\tif err != nil && !strings.Contains(strings.Join(args, \" \"), \"--confirm-on-surface\") {\n\t\t\treturn fail(\"refusing: no controlling terminal\"",
    "test": "TestOperatorFinishF1",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 90,
    "stage": "M1",
    "control": "die"
  },
  "f1-tty-prompt": {
    "path": "cmd/ailang-worldd/episode_finish_operator.go",
    "old": "Open this pull request? [y/N] ",
    "new": "Open this pull request? [n/Y] ",
    "test": "TestOperatorFinishTTYParity",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 158,
    "stage": "M1",
    "control": "die"
  },
  "proof-no-recheck": {
    "path": "cmd/ailang-worldd/operator_finish_preflight.go",
    "old": "p.authority.Rechecked() &&",
    "new": "p.authority != nil &&",
    "test": "TestOperatorFinishNoProof",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 59,
    "stage": "M1",
    "control": "die"
  },
  "proof-deadline": {
    "path": "cmd/ailang-worldd/operator_finish_preflight.go",
    "old": "if now < row.CreatedAt || now >= row.ExpiresAt {",
    "new": "if now < row.CreatedAt {",
    "test": "TestOperatorFinishNoProof",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 44,
    "stage": "M1",
    "control": "die"
  },
  "commit-author-email": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "author+\">\" != \"author \"+agent+\" via AILANG World <ailang-world@invalid>\"",
    "new": "!strings.HasPrefix(author, \"author \"+agent+\" via AILANG World <\")",
    "test": "TestOperatorCommitIdentity",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 340,
    "stage": "M1",
    "control": "die"
  },
  "commit-match-tree": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "lines[0] != \"tree \"+tree ||",
    "new": "lines[0] != \"tree \"+tree && false ||",
    "test": "TestOperatorCommitIdentity",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 335,
    "stage": "M1",
    "control": "die"
  },
  "commit-before-reservation": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if err := w.reserve(key); err != nil {",
    "new": "if err := error(nil); err != nil {",
    "test": "TestOperatorCommitReconcile",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 151,
    "stage": "M1",
    "control": "die"
  },
  "proof-bytes": {
    "path": "cmd/ailang-worldd/operator_finish_local.go",
    "old": "if !w.boundManifest(w.generation) {",
    "new": "if !w.boundManifest(w.generation) && false {",
    "test": "TestOperatorFinishPrep",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 84,
    "stage": "M1",
    "control": "die"
  },
  "git-first-use-config": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "err != nil || config != spec.Snapshot.Config {",
    "new": "(err != nil || config != spec.Snapshot.Config) && false {",
    "test": "TestOperatorFinishConfigReplacement",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 348,
    "stage": "M1",
    "control": "die"
  },
  "git-attribution-base-delta": {
    "path": "cmd/ailang-worldd/operator_finish_manifest.go",
    "old": "if path != \"\" {\n\t\t\tpaths = append(paths, path)",
    "new": "if path != \"\" && false {\n\t\t\tpaths = append(paths, path)",
    "test": "TestOperatorFinishManifestCapture",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 46,
    "stage": "M1",
    "control": "die"
  },
  "ready-before-claim": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if err := prepareOperatorFinish(ctx, c, m.Remote.URL); err != nil {",
    "new": "if err := prepareOperatorFinish(ctx, c, m.Remote.URL); err != nil && false {",
    "test": "TestOperatorFinishReadiness",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 214,
    "stage": "M1",
    "control": "die"
  },
  "ready-gh-not-git": {
    "path": "cmd/ailang-worldd/operator_finish_auth.go",
    "old": "env = append(env, \"GIT_CONFIG_NOSYSTEM=1\", \"GIT_CONFIG_GLOBAL=/dev/null\")",
    "new": "if permission == \"WRITE\" { return nil }; env = append(env, \"GIT_CONFIG_NOSYSTEM=1\", \"GIT_CONFIG_GLOBAL=/dev/null\")",
    "test": "TestOperatorFinishReadiness",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 98,
    "stage": "M1",
    "control": "die"
  },
  "runner-positive-claim": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "return db.ClaimOperatorAction(write, cas, h.Owner, src.clock())",
    "new": "return nil",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 296,
    "stage": "M1",
    "control": "die"
  },
  "runner-claim-cas-version": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "cas.Generation++",
    "new": "cas.Generation += 0",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 305,
    "stage": "M1",
    "control": "die"
  },
  "runner-prep-caller": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "return runEpisodeFinishWithPrep(ctx, w)",
    "new": "_ = w\n\treturn nil",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 234,
    "stage": "M1",
    "control": "die"
  },
  "runner-commit-caller": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "return in.git.commit(bound, w.spool, spec, contents, tree, in.report.Start)",
    "new": "return nil",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 344,
    "stage": "M1",
    "control": "die"
  },
  "runner-world-check-writer": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "return commitOneObject(bound, client, in.token, outcome, taskOutcomeWrittenBy, \"\", \"outcome record\")",
    "new": "return 1, nil",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 338,
    "stage": "M1",
    "control": "die"
  },
  "runner-claim-owner": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "row.Owner = h.Owner.Nonce",
    "new": "row.Owner = \"\"",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 306,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-start": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if err != nil || start == nil || len(start.Resumes) != 0 ||",
    "new": "if err != nil || false && start == nil || len(start.Resumes) != 0 ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 41,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-resumes": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "len(start.Resumes) != 0 ||",
    "new": "len(start.Resumes) != 0 && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 41,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-preset-allowlist": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "!slices.Contains(c.PresetNames, start.Preset.Name) ||",
    "new": "!slices.Contains(c.PresetNames, start.Preset.Name) && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 41,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-remote-allowlist": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "!slices.Contains(c.FinishRemotes, start.Remote.URL) {",
    "new": "!slices.Contains(c.FinishRemotes, start.Remote.URL) && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 42,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-launch-version": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "err != nil || launch.Key.Version != 4",
    "new": "err != nil || launch.Key.Version != 4 && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 46,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-live-agent": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "live && birth == launch.Receipt.Birth",
    "new": "live && birth == launch.Receipt.Birth && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 49,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-health-status": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "health.Status != \"ok\" ||",
    "new": "health.Status != \"ok\" && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 53,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-degraded": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "slices.Contains(health.DegradedEpisodes, ep) ||",
    "new": "slices.Contains(health.DegradedEpisodes, ep) && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 53,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-build": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "health.Build.ExecutableSHA256 != c.Executables[e.paths.Binary]",
    "new": "health.Build.ExecutableSHA256 != c.Executables[e.paths.Binary] && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 54,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-report-read": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if err != nil || rep.Episode != ep ||",
    "new": "if err != nil && false || rep.Episode != ep ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 58,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-report-episode": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "rep.Episode != ep ||",
    "new": "rep.Episode != ep && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 58,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-report-intents": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "rep.Uncommitted != 0 ||",
    "new": "rep.Uncommitted != 0 && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 58,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-recorded-finish": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "rep.Finish != nil ||",
    "new": "rep.Finish != nil && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 58,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-report-start": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "!reflect.DeepEqual(rep.Start, start) {",
    "new": "!reflect.DeepEqual(rep.Start, start) && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 58,
    "stage": "M1",
    "control": "die"
  },
  "runner-drift-preset": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "serviceHash(raw) != start.Preset.SHA256 || string(raw) != start.Preset.Bytes",
    "new": "(serviceHash(raw) != start.Preset.SHA256 || string(raw) != start.Preset.Bytes) && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 66,
    "stage": "M1",
    "control": "die"
  },
  "runner-drift-preset-decode": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "strictPresetJSON(raw, &preset) != nil || preset.Name != start.Preset.Name",
    "new": "(strictPresetJSON(raw, &preset) != nil || preset.Name != start.Preset.Name) && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 70,
    "stage": "M1",
    "control": "die"
  },
  "runner-drift-workspace": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "physical != root || p.Addr != c.Address",
    "new": "(physical != root || p.Addr != c.Address) && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 76,
    "stage": "M1",
    "control": "die"
  },
  "runner-drift-world": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "store.SameDatabase(p.DB, c.World); !equal",
    "new": "store.SameDatabase(p.DB, c.World); !equal && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 79,
    "stage": "M1",
    "control": "die"
  },
  "runner-drift-store": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "err != nil || !equal {",
    "new": "(err != nil || !equal) && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 83,
    "stage": "M1",
    "control": "die"
  },
  "runner-drift-profile": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "sessions[0].ProfileDigest != profile.Digest",
    "new": "sessions[0].ProfileDigest != profile.Digest && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 104,
    "stage": "M1",
    "control": "die"
  },
  "runner-drift-pinned-check": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "!reflect.DeepEqual(checks, start.Checks)",
    "new": "!reflect.DeepEqual(checks, start.Checks) && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 96,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-sessions": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "len(sessions) != 1 ||",
    "new": "len(sessions) != 1 && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 100,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-session-expiry": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "sessions[0].ExpiresAt <= src.clock() ||",
    "new": "sessions[0].ExpiresAt <= src.clock() && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 100,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-project": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "sessions[0].Project != preset.Profile",
    "new": "sessions[0].Project != preset.Profile && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 101,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-session-file": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "runnerPath(c.Home, sessionDir, c.UID, true, 0700) != nil",
    "new": "runnerPath(c.Home, sessionDir, c.UID, true, 0700) != nil && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 108,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-token-mode": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "mode.Perm() != 0600",
    "new": "mode.Perm() != 0600 && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 112,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-token": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "err != nil || resolved.Success == nil || resolved.Success.EpisodeID != ep || resolved.Success.ExpiresAt != sessions[0].ExpiresAt",
    "new": "(err != nil || resolved.Success == nil || resolved.Success.EpisodeID != ep || resolved.Success.ExpiresAt != sessions[0].ExpiresAt) && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 117,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-check-sandbox": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if check.Unsandboxed {",
    "new": "if check.Unsandboxed && false {",
    "test": "TestOperatorFinishReadiness",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 1,
    "line": 210,
    "stage": "M1",
    "control": "die"
  },
  "runner-elig-commit-identity": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if err != nil || name == \"\" {",
    "new": "if (err != nil || name == \"\") && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 130,
    "stage": "M1",
    "control": "die"
  },
  "runner-lock-busy": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "return operatorFinishBusy(row.ID)",
    "new": "return store.ErrOperatorConflict",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 220,
    "stage": "M1",
    "control": "die"
  },
  "runner-lock-ref": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "store.OpenOperatorEpisodeLock(c.Credentials, row.Episode)",
    "new": "store.OpenOperatorEpisodeLock(c.Credentials, row.Episode+\"-other\")",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 218,
    "stage": "M1",
    "control": "die"
  },
  "runner-manifest-mismatch": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "!reflect.DeepEqual(reviewed, m) ||",
    "new": "!reflect.DeepEqual(reviewed, m) && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 278,
    "stage": "M1",
    "control": "die"
  },
  "runner-manifest-protected": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "len(m.Protected) != 0 || len(m.Unattributed) != 0 {",
    "new": "(len(m.Protected) != 0 || len(m.Unattributed) != 0) && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 278,
    "stage": "M1",
    "control": "die"
  },
  "runner-catalogue-fresh": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "live.CatalogueGeneration != generation ||",
    "new": "live.CatalogueGeneration != generation && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 293,
    "stage": "M1",
    "control": "die"
  },
  "runner-heartbeat-owner": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "err != nil || live.Owner != h.Owner || live.CatalogueGeneration",
    "new": "err != nil || live.Owner != h.Owner && false || live.CatalogueGeneration",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 293,
    "stage": "M1",
    "control": "die"
  },
  "runner-heartbeat-ready": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "!live.Available(h.Build, \"finish\", src.clock()) {",
    "new": "!live.Available(h.Build, \"finish\", src.clock()) && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 293,
    "stage": "M1",
    "control": "die"
  },
  "runner-life-quiet": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "e.proveOldTerminalQuiet() == nil {",
    "new": "(e.proveOldTerminalQuiet() == nil || true) {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 285,
    "stage": "M1",
    "control": "die"
  },
  "runner-life-pair": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "e.checkStartPair(c, src.workspace) == nil &&",
    "new": "(e.checkStartPair(c, src.workspace) == nil || true) &&",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 285,
    "stage": "M1",
    "control": "die"
  },
  "runner-effect-clock": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "|| !live.Available(h.Build, \"finish\", src.clock()) ||",
    "new": "|| !live.Available(h.Build, \"finish\", src.clock()) && false ||",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 320,
    "stage": "M1",
    "control": "die"
  },
  "runner-effect-inputs": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "!finishReviewInputsMatch(current.manifest, m) || !current.baseMatches(bound)",
    "new": "(!finishReviewInputsMatch(current.manifest, m) || !current.baseMatches(bound)) && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 328,
    "stage": "M1",
    "control": "die"
  },
  "runner-effect-key": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "_, err = verifyOperatorFinishPrep(bound, ro, a, src.clock(), cas, row.Owner)\n\t\treturn err",
    "new": "return nil",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 331,
    "stage": "M1",
    "control": "die"
  },
  "runner-observe-after-expiry": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "verifyAt = min(verifyAt, row.ExpiresAt-1)",
    "new": "verifyAt = max(verifyAt, row.ExpiresAt)",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 203,
    "stage": "M1",
    "control": "die"
  },
  "runner-old-observation": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "observation := row.Owner != \"\" && journal.CommitPhase == \"reserved\"",
    "new": "observation := row.Owner != \"\" && journal.CommitPhase == \"reserved\" && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 198,
    "stage": "M1",
    "control": "die"
  },
  "runner-catalogue-approve-size": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if _, err := v.Action(1, false); err == nil {",
    "new": "if _, err := v.Action(1, false); err == nil || true {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 171,
    "stage": "M1",
    "control": "die"
  },
  "runner-loop-turn": {
    "path": "cmd/ailang-worldd/operator_runner.go",
    "old": "if h.Readiness == \"both\" {",
    "new": "if h.Readiness == \"both\" && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 157,
    "stage": "M1",
    "control": "die"
  },
  "runner-turn-heartbeat": {
    "path": "cmd/ailang-worldd/operator_runner.go",
    "old": "if h.Readiness == \"\" {",
    "new": "if h.Readiness == \"\" || true {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 204,
    "stage": "M1",
    "control": "die"
  },
  "runner-turn-sink": {
    "path": "cmd/ailang-worldd/operator_runner.go",
    "old": "return class.Error()",
    "new": "return err.Error()",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 191,
    "stage": "M1",
    "control": "die"
  },
  "runner-catalogue-root": {
    "path": "cmd/ailang-worldd/operator_catalogue.go",
    "old": "src.workspace = r.Workspace",
    "new": "src.workspace = \"\"",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 39,
    "stage": "M1",
    "control": "die"
  },
  "runner-catalogue-bound": {
    "path": "cmd/ailang-worldd/operator_catalogue.go",
    "old": "if len(raw) > 16*approval.MaxOperatorActionBytes {",
    "new": "if len(raw) > 16*approval.MaxOperatorActionBytes && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 45,
    "stage": "M1",
    "control": "die"
  },
  "runner-catalogue-preserve-start": {
    "path": "cmd/ailang-worldd/operator_catalogue.go",
    "old": "if len(c.FinishRemotes) != 0 {",
    "new": "if len(c.FinishRemotes) != 0 && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 38,
    "stage": "M1",
    "control": "die"
  },
  "terminal-lock-active-id": {
    "path": "cmd/ailang-worldd/operator_episode_lock.go",
    "old": "active action %s; run ailang-worldd service runner status",
    "new": "runner active %s; run ailang-worldd service runner status",
    "test": "TestOperatorTerminalExclusion",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "line": 64,
    "stage": "M1",
    "control": "die"
  },
  "route-finish-admission": {
    "path": "host/daemon/surface_operator.go",
    "old": "in.Kind == \"finish\")",
    "new": "in.Kind == \"finish\" && false)",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 31,
    "stage": "M1",
    "control": "die"
  },
  "route-start-task": {
    "path": "host/daemon/surface_operator.go",
    "old": "(a.Intent.Kind == \"start\") != (body.Task != \"\")",
    "new": "(a.Intent.Kind == \"start\") && body.Task == \"\"",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 209,
    "stage": "M1",
    "control": "die"
  },
  "route-finish-task": {
    "path": "host/daemon/surface_operator.go",
    "old": "(a.Intent.Kind == \"start\") != (body.Task != \"\")",
    "new": "(a.Intent.Kind == \"start\") != (body.Task != \"\") && a.Intent.Kind != \"finish\"",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 209,
    "stage": "M1",
    "control": "die"
  },
  "route-finish-health": {
    "path": "host/daemon/surface_operator.go",
    "old": "if !h.Available(build, \"finish\", now) {",
    "new": "if !h.Available(build, \"finish\", now) && false {",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 48,
    "stage": "M1",
    "control": "die"
  },
  "route-finish-generation": {
    "path": "host/daemon/surface_operator.go",
    "old": "h.CatalogueGeneration != body.CatalogueGeneration",
    "new": "false && h.CatalogueGeneration != body.CatalogueGeneration",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 256,
    "stage": "M1",
    "control": "die"
  },
  "route-finish-store": {
    "path": "host/daemon/surface_operator.go",
    "old": "entry.Store != served {",
    "new": "entry.Store != served && false {",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 53,
    "stage": "M1",
    "control": "die"
  },
  "route-protected-and-unattributed": {
    "path": "host/daemon/surface_operator.go",
    "old": "len(entry.Protected) != 0 || len(entry.Unattributed) != 0",
    "new": "(len(entry.Protected) != 0 || len(entry.Unattributed) != 0) && false",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 56,
    "stage": "M1",
    "control": "die"
  },
  "route-idempotence-kind": {
    "path": "host/daemon/surface_operator.go",
    "old": "fact = \"finish_catalogue\"",
    "new": "fact = \"mint_catalogue\"",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 238,
    "stage": "M1",
    "control": "die"
  },
  "route-operation-generation": {
    "path": "host/daemon/surface_operator.go",
    "old": "Value: strconv.FormatInt(body.CatalogueGeneration, 10)",
    "new": "Value: strconv.FormatInt(body.CatalogueGeneration+1, 10)",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 241,
    "stage": "M1",
    "control": "die"
  },
  "route-operation-generation-conflict": {
    "path": "host/daemon/surface_operator.go",
    "old": "!slices.Contains(frozen.Resolution, approval.OperatorResolved{Name: fact, Value: strconv.FormatInt(body.CatalogueGeneration, 10)})",
    "new": "!slices.Contains(frozen.Resolution, approval.OperatorResolved{Name: fact, Value: strconv.FormatInt(body.CatalogueGeneration, 10)}) && false",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 241,
    "stage": "M1",
    "control": "die"
  },
  "route-review-no-register": {
    "path": "host/daemon/surface_operator.go",
    "old": "out.Outcome, out.Reason, out.Argv, out.Text = \"review\", \"\", action.Argv, action.Review",
    "new": "out.Outcome, out.Reason, out.Argv, out.Text = \"terminal\", \"\", action.Argv, action.Review",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 80,
    "stage": "M1",
    "control": "die"
  },
  "route-confirm-scope": {
    "path": "host/daemon/surface_operator_confirm.go",
    "old": "a.Kind == \"finish\" && a.Options.ConfirmOnSurface",
    "new": "a.Kind == \"finish\" && a.Options.ConfirmOnSurface && false",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 96,
    "stage": "M1",
    "control": "die"
  },
  "route-confirm-without-surface": {
    "path": "host/daemon/surface_operator_confirm.go",
    "old": "a.Kind == \"finish\" && a.Options.ConfirmOnSurface",
    "new": "a.Kind == \"finish\"",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 96,
    "stage": "M1",
    "control": "die"
  },
  "status-finish-projection": {
    "path": "host/daemon/surface_operator_read.go",
    "old": "if status.Href != \"\" && row.Kind == \"finish\" && row.Owner != \"\" {",
    "new": "if status.Href != \"\" && row.Kind == \"finish\" && row.Owner != \"\" && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 128,
    "stage": "M1",
    "control": "die"
  },
  "status-evidence-leak": {
    "path": "host/daemon/surface_operator_read.go",
    "old": "Name: check.Name, Result: check.Result",
    "new": "Name: check.Ref, Result: check.Result",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 135,
    "stage": "M1",
    "control": "die"
  },
  "status-sha": {
    "path": "host/daemon/surface_operator_read.go",
    "old": "status.LocalCommit = j.Commit.SHA",
    "new": "status.LocalCommit = \"\"",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 139,
    "stage": "M1",
    "control": "die"
  },
  "status-open-window": {
    "path": "host/daemon/surface_operator_read.go",
    "old": "case j.OpenCheck || j.CommitPhase == \"reserved\":",
    "new": "case (j.OpenCheck || j.CommitPhase == \"reserved\") && false:",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 141,
    "stage": "M1",
    "control": "die"
  },
  "status-definite-failure": {
    "path": "host/daemon/surface_operator_read.go",
    "old": "case j.Commit.State == \"commit-failed\":",
    "new": "case j.Commit.State == \"commit-failed\" && false:",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 143,
    "stage": "M1",
    "control": "die"
  },
  "status-finish-auth": {
    "path": "host/daemon/surface_operator_read.go",
    "old": "status.Reason = errOperatorFinishAuth.Error()",
    "new": "status.Reason = \"\"",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "line": 151,
    "stage": "M1",
    "control": "die"
  },
  "js-render-export": {
    "path": "host/surfacepage/app/app.js",
    "old": "operatorRequestControl, operatorCardHref, showOperatorReview };",
    "new": "operatorRequestControl, operatorCardHref };",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 1589,
    "stage": "M1",
    "control": "die"
  },
  "js-finish-kind": {
    "path": "host/surfacepage/app/app.js",
    "old": "return argv[2] === 'finish' &&",
    "new": "return argv[2] === 'finish' && false &&",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 1445,
    "stage": "M1",
    "control": "die"
  },
  "js-finish-surface": {
    "path": "host/surfacepage/app/app.js",
    "old": "argv[2] === 'finish' && argv[6] === '--confirm-on-surface' &&",
    "new": "argv[2] === 'finish' &&",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 1445,
    "stage": "M1",
    "control": "die"
  },
  "js-finish-ready": {
    "path": "host/surfacepage/app/app.js",
    "old": "argv.length === 8 && argv[7] === '--ready'",
    "new": "argv.length === 8",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 1446,
    "stage": "M1",
    "control": "die"
  },
  "js-finish-argv-count": {
    "path": "host/surfacepage/app/app.js",
    "old": "argv.length === 7 || (argv.length === 8 && argv[7] === '--ready')",
    "new": "argv.length >= 7",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 1446,
    "stage": "M1",
    "control": "die"
  },
  "js-finish-empty-task": {
    "path": "host/surfacepage/app/app.js",
    "old": "argv[2] === 'finish' ? task === '' : taskValid(task)",
    "new": "argv[2] === 'finish' ? true : taskValid(task)",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 1448,
    "stage": "M1",
    "control": "die"
  },
  "js-retain-draft": {
    "path": "host/surfacepage/app/app.js",
    "old": "&& operatorTaskValid(r.argv, r.task)) return d;",
    "new": "&& taskValid(r.task)) return d;",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 1454,
    "stage": "M1",
    "control": "die"
  },
  "js-request-kind": {
    "path": "host/surfacepage/app/app.js",
    "old": "b.kind !== kind || b.episode",
    "new": "b.episode",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 1459,
    "stage": "M1",
    "control": "die"
  },
  "js-request-bound-link": {
    "path": "host/surfacepage/app/app.js",
    "old": "operatorCardHref(b, request.argv[3], request.argv[2])",
    "new": "operatorCardHref(b, request.argv[3])",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 1486,
    "stage": "M1",
    "control": "die"
  },
  "js-taskless-dom": {
    "path": "host/surfacepage/app/app.js",
    "old": "if (!finish) area.append(label,",
    "new": "if (true) area.append(label,",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 3555,
    "stage": "M1",
    "control": "die"
  },
  "js-taskless-focus": {
    "path": "host/surfacepage/app/app.js",
    "old": "(finish || task.readOnly ? request : task).focus();",
    "new": "(task.readOnly ? request : task).focus();",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 3577,
    "stage": "M1",
    "control": "die"
  },
  "js-taskless-button": {
    "path": "host/surfacepage/app/app.js",
    "old": "finish ? 'Request approval' : 'Request action'",
    "new": "'Request action'",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 3561,
    "stage": "M1",
    "control": "die"
  },
  "js-taskless-post": {
    "path": "host/surfacepage/app/app.js",
    "old": "await control.submit(task.value).then",
    "new": "await control.submit('fake task').then",
    "test": "TestOperatorFinishBarReview",
    "package": "host/surfacepage",
    "access": "local",
    "occurrence": 0,
    "line": 3568,
    "stage": "M1",
    "control": "die"
  },
  "s9d-cli-dispatch": {
    "path": "cmd/ailang-worldd/operator_service.go",
    "old": "err := e.printFinishStatus(args[1], out)",
    "new": "err := error(nil)",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 318,
    "stage": "M2",
    "control": "die"
  },
  "s9d-cli-sha": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "fmt.Fprintln(out, \"local commit\", j.Commit.SHA,",
    "new": "fmt.Fprintln(out, \"local commit\", \"\",",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 384,
    "stage": "M2",
    "control": "die"
  }
}
````

````text changed-M1
cmd/ailang-worldd/operator_activation_test.go
cmd/ailang-worldd/operator_catalogue.go
cmd/ailang-worldd/operator_finish_manifest.go
cmd/ailang-worldd/operator_finish_manifest_test.go
cmd/ailang-worldd/operator_finish_resume_test.go
cmd/ailang-worldd/operator_finish_route_test.go
cmd/ailang-worldd/operator_finish_runner.go
cmd/ailang-worldd/operator_finish_runner_test.go
cmd/ailang-worldd/operator_finish_tree.go
cmd/ailang-worldd/operator_finish_tree_test.go
cmd/ailang-worldd/operator_launch_test.go
cmd/ailang-worldd/operator_runner.go
cmd/ailang-worldd/operator_worktree_test.go
docs/QUICKSTART.md
host/daemon/surface_operator.go
host/daemon/surface_operator_confirm.go
host/daemon/surface_operator_finish_test.go
host/daemon/surface_operator_read.go
host/surfacepage/app/app.js
host/surfacepage/app_operator_test.go
host/surfacepage/app_test.go
````

````text changed-M2-cumulative
cmd/ailang-worldd/operator_activation_test.go
cmd/ailang-worldd/operator_catalogue.go
cmd/ailang-worldd/operator_finish_manifest.go
cmd/ailang-worldd/operator_finish_manifest_test.go
cmd/ailang-worldd/operator_finish_resume_test.go
cmd/ailang-worldd/operator_finish_route_test.go
cmd/ailang-worldd/operator_finish_runner.go
cmd/ailang-worldd/operator_finish_runner_test.go
cmd/ailang-worldd/operator_finish_tree.go
cmd/ailang-worldd/operator_finish_tree_test.go
cmd/ailang-worldd/operator_launch_test.go
cmd/ailang-worldd/operator_runner.go
cmd/ailang-worldd/operator_service.go
cmd/ailang-worldd/operator_worktree_test.go
docs/QUICKSTART.md
host/daemon/surface_operator.go
host/daemon/surface_operator_confirm.go
host/daemon/surface_operator_finish_test.go
host/daemon/surface_operator_read.go
host/surfacepage/app/app.js
host/surfacepage/app_operator_test.go
host/surfacepage/app_test.go
````

````text status-M1
M	cmd/ailang-worldd/operator_activation_test.go
M	cmd/ailang-worldd/operator_catalogue.go
M	cmd/ailang-worldd/operator_finish_manifest.go
M	cmd/ailang-worldd/operator_finish_manifest_test.go
A	cmd/ailang-worldd/operator_finish_resume_test.go
A	cmd/ailang-worldd/operator_finish_route_test.go
A	cmd/ailang-worldd/operator_finish_runner.go
A	cmd/ailang-worldd/operator_finish_runner_test.go
M	cmd/ailang-worldd/operator_finish_tree.go
M	cmd/ailang-worldd/operator_finish_tree_test.go
M	cmd/ailang-worldd/operator_launch_test.go
M	cmd/ailang-worldd/operator_runner.go
M	cmd/ailang-worldd/operator_worktree_test.go
M	docs/QUICKSTART.md
M	host/daemon/surface_operator.go
M	host/daemon/surface_operator_confirm.go
A	host/daemon/surface_operator_finish_test.go
M	host/daemon/surface_operator_read.go
M	host/surfacepage/app/app.js
M	host/surfacepage/app_operator_test.go
M	host/surfacepage/app_test.go
````

````text status-M2-cumulative
M	cmd/ailang-worldd/operator_activation_test.go
M	cmd/ailang-worldd/operator_catalogue.go
M	cmd/ailang-worldd/operator_finish_manifest.go
M	cmd/ailang-worldd/operator_finish_manifest_test.go
A	cmd/ailang-worldd/operator_finish_resume_test.go
A	cmd/ailang-worldd/operator_finish_route_test.go
A	cmd/ailang-worldd/operator_finish_runner.go
A	cmd/ailang-worldd/operator_finish_runner_test.go
M	cmd/ailang-worldd/operator_finish_tree.go
M	cmd/ailang-worldd/operator_finish_tree_test.go
M	cmd/ailang-worldd/operator_launch_test.go
M	cmd/ailang-worldd/operator_runner.go
M	cmd/ailang-worldd/operator_service.go
M	cmd/ailang-worldd/operator_worktree_test.go
M	docs/QUICKSTART.md
M	host/daemon/surface_operator.go
M	host/daemon/surface_operator_confirm.go
A	host/daemon/surface_operator_finish_test.go
M	host/daemon/surface_operator_read.go
M	host/surfacepage/app/app.js
M	host/surfacepage/app_operator_test.go
M	host/surfacepage/app_test.go
````

````text named-mutant-ids.txt
j01
j02
j04
j05
j06
j07
j08
j10
j13
j14
j16
j19
j20
j26
j22
j23
j24
j27
j32
j33
j35
f1-surface-bypass
f1-tty-prompt
proof-no-recheck
proof-deadline
commit-author-email
commit-match-tree
commit-before-reservation
proof-bytes
git-first-use-config
git-attribution-base-delta
ready-before-claim
ready-gh-not-git
runner-positive-claim
runner-claim-cas-version
runner-prep-caller
runner-commit-caller
runner-world-check-writer
runner-claim-owner
runner-elig-start
runner-elig-resumes
runner-elig-preset-allowlist
runner-elig-remote-allowlist
runner-elig-launch-version
runner-elig-live-agent
runner-elig-health-status
runner-elig-degraded
runner-elig-build
runner-elig-report-read
runner-elig-report-episode
runner-elig-report-intents
runner-elig-recorded-finish
runner-elig-report-start
runner-drift-preset
runner-drift-preset-decode
runner-drift-workspace
runner-drift-world
runner-drift-store
runner-drift-profile
runner-drift-pinned-check
runner-elig-sessions
runner-elig-session-expiry
runner-elig-project
runner-elig-session-file
runner-elig-token-mode
runner-elig-token
runner-elig-check-sandbox
runner-elig-commit-identity
runner-lock-busy
runner-lock-ref
runner-manifest-mismatch
runner-manifest-protected
runner-catalogue-fresh
runner-heartbeat-owner
runner-heartbeat-ready
runner-life-quiet
runner-life-pair
runner-effect-clock
runner-effect-inputs
runner-effect-key
runner-observe-after-expiry
runner-old-observation
runner-catalogue-approve-size
runner-loop-turn
runner-turn-heartbeat
runner-turn-sink
runner-catalogue-root
runner-catalogue-bound
runner-catalogue-preserve-start
terminal-lock-active-id
route-finish-admission
route-start-task
route-finish-task
route-finish-health
route-finish-generation
route-finish-store
route-protected-and-unattributed
route-idempotence-kind
route-operation-generation
route-operation-generation-conflict
route-review-no-register
route-confirm-scope
route-confirm-without-surface
status-finish-projection
status-evidence-leak
status-sha
status-open-window
status-definite-failure
status-finish-auth
js-render-export
js-finish-kind
js-finish-surface
js-finish-ready
js-finish-argv-count
js-finish-empty-task
js-retain-draft
js-request-kind
js-request-bound-link
js-taskless-dom
js-taskless-focus
js-taskless-button
js-taskless-post
s9d-cli-dispatch
s9d-cli-sha
````

````text whole-mutant-ids.txt
j08
j13
j14
j16
j19
j35
ready-before-claim
git-first-use-config
git-attribution-base-delta
runner-loop-turn
runner-positive-claim
runner-commit-caller
runner-lock-busy
route-finish-admission
route-confirm-scope
status-evidence-leak
js-taskless-dom
````

## 7. Executor and snapshot contract

Run this block with Bash in the sprint worktree. The default executor applies M1 only with `/usr/bin/patch -p1 -E`; the optional M2 fence is a rehearsal for S9d. A controller or planner round trip runs the same block with `PLAN_APPLY_MODE=rehearsal`, a fresh archive under `$S/rt-apply` and both fences. The two product snapshots are cumulative. In the worktree, `snapshot.sh` writes `$STATE/snap/<M>/` from changes relative to the pinned base, excludes the two plan artifacts, and records `STATUS`, `DELETED` and `MANIFEST.sha256`. Archive rehearsal uses the measured status list after reconstructing the exact patch. An absent top-level root is handled explicitly.

The acceptance success signal is `ACCEPT_OK <M> <mode>`. The application ends each milestone with `APPLY_ACCEPT_SNAPSHOT_OK <M>` and ends the block with `FINAL_APPLY_RC 0`. Commit M1 only, after real-listener, Chrome, race and outside evidence. The controller runs the acceptance suite and separate JS syntax plus Node test leg via `accept.sh`. The commit guard executes the exact `go test -count=1 -timeout 20m ./host/verifygate/ ./host/broker/` with isolated ailang state and compares all failures against base. Commit-time cmd/daemon work uses named tests. No whole cmd or daemon run is embedded in it.

````bash commit-M1
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
S="$HOME/.ailang/state/world-iter450/plan-scratch"
bash "$S/accept.sh" "$PWD" M1 outside &&
  bash "$S/commit-guards.sh" "$PWD" "$S/commit-guards-M1-$(date +%s)"
````

````bash apply
(
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
B=91befdab7fd4081b4dc127cc96cc4eeb1723cffa
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter450-row243-v4s9c
STATE="$HOME/.ailang/state/world-iter450"
S="$STATE/plan-scratch"
PLAN="$WT/design_docs/planned/w-world-address-bar-verbs-v4s9c-sprint-plan.md"
MODE=${PLAN_APPLY_MODE:-executor}
case "$MODE" in
  rehearsal)
    R="$S/rt-apply"
    mkdir "$R"
    T="$R/export"
    TOOLS="$R/tools"
    mkdir "$T" "$TOOLS"
    git -C "$WT" archive "$B" | tar -x -C "$T"
    milestones=(M1 M2)
    ;;
  executor)
    T="$WT"
    TOOLS="$STATE/exec-tools"
    mkdir -p "$TOOLS"
    test "$(git -C "$T" rev-parse HEAD)" = "$B"
    test -z "$(git -C "$T" status --porcelain -- cmd host docs app)"
    milestones=(M1)
    ;;
  *) exit 2;;
esac
fence() {
  awk -v label="$1" '$0 == "````" label {p=1; next} p && $0 == "````" {exit} p' "$PLAN" > "$2"
  test -s "$2"
}
for name in accept.sh gates.sh snapshot.sh env.sh selectors.sh callers.sh outside.sh commit-guards.sh named-mutants.sh whole-mutants.sh; do
  fence "bash $name" "$TOOLS/$name"
done
for name in bounded.py failset.py verify-scratch.py mutate.py; do
  fence "python $name" "$TOOLS/$name"
done
fence 'json mutants.json' "$TOOLS/mutants.json"
for name in changed-M1 changed-M2-cumulative status-M1 status-M2-cumulative named-mutant-ids.txt whole-mutant-ids.txt; do
  fence "text $name" "$TOOLS/$name"
done
for M in "${milestones[@]}"; do
  fence "diff $M" "$TOOLS/$M.patch"
  (cd "$T" && /usr/bin/patch -p1 -E --dry-run < "$TOOLS/$M.patch") > "$TOOLS/$M.dry-run" 2>&1
  if grep -Eiq 'offset|fuzz|FAILED|reject' "$TOOLS/$M.dry-run"; then cat "$TOOLS/$M.dry-run"; exit 1; fi
  (cd "$T" && /usr/bin/patch -p1 -E < "$TOOLS/$M.patch")
  bash "$TOOLS/accept.sh" "$T" "$M" fast
  if [ "$MODE" = rehearsal ]; then
    status="$TOOLS/status-$M"
    [ "$M" = M1 ] || status="$TOOLS/status-M2-cumulative"
    SNAP_PARENT="$R/snap" SNAP_STATUS_LIST="$status" bash "$TOOLS/snapshot.sh" "$T" "$M"
    for d in cmd host docs app; do
      if [ -d "$T/$d" ] || [ -d "$S/st-$M/$d" ]; then
        diff -r "$T/$d" "$S/st-$M/$d" > "$R/diff-$M-$d"
        echo "ROUNDTRIP_DIFF_OK $M $d"
      else
        echo "ROUNDTRIP_ABSENT $M $d"
      fi
    done
    diff -r "$T" "$S/st-$M" > "$R/diff-$M-unfiltered"
    echo "ROUNDTRIP_DIFF_OK $M unfiltered"
  else
    bash "$TOOLS/snapshot.sh" "$T" "$M"
  fi
  echo "APPLY_ACCEPT_SNAPSHOT_OK $M"
done
echo "FINAL_APPLY_RC 0"
)
````

Verbatim rehearsal printed `APPLY_ACCEPT_SNAPSHOT_OK M1`, `APPLY_ACCEPT_SNAPSHOT_OK M2` and `FINAL_APPLY_RC 0`. Both base-relative snapshots pass SHA-256 inventory verification (M1 21/0/21, M2 22/0/22 for STATUS/DELETED/MANIFEST). `diff -r` status is 0 for cmd, host, docs and the full archive at both stages. The top-level app directory is absent on both sides. See `$S/rt-apply.receipt`.

## Appendix P. Compiled patches

### M1 — Signed no-terminal local finish

````diff M1
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
diff --git a/cmd/ailang-worldd/operator_finish_manifest.go b/cmd/ailang-worldd/operator_finish_manifest.go
index 3f2d399..4a78943 100644
--- a/cmd/ailang-worldd/operator_finish_manifest.go
+++ b/cmd/ailang-worldd/operator_finish_manifest.go
@@ -36,10 +36,25 @@ func (i operatorFinishInputs) frozen(ctx context.Context) (store.OperatorFinishR
 	if err := validateFinishChecks(v.Checks, nil, i.git.dir); err != nil {
 		return v, errOperatorFinishFence
 	}
+	committed, err := i.git.git(ctx, nil, nil, "diff", "--name-only", "-z", "--no-renames",
+		"--no-ext-diff", "--no-textconv", v.Base.SHA, tree.Head, "--")
+	if err != nil {
+		return v, errOperatorFinishFence
+	}
 	var paths []string
+	for path := range strings.SplitSeq(strings.TrimSuffix(committed, "\x00"), "\x00") {
+		if path != "" {
+			paths = append(paths, path)
+		}
+	}
 	for _, p := range tree.Paths {
 		paths = append(paths, p.Path)
 	}
+	slices.Sort(paths)
+	paths = slices.Compact(paths)
+	if len(paths) > 257 {
+		return v, errOperatorFinishFence
+	}
 	v.Protected = finishProtectedPaths(i.profile, paths)
 	v.Unattributed = finishUnattributed(paths, i.report)
 	v.Commit.Snapshot = tree
diff --git a/cmd/ailang-worldd/operator_finish_manifest_test.go b/cmd/ailang-worldd/operator_finish_manifest_test.go
index 46192f5..7441e76 100644
--- a/cmd/ailang-worldd/operator_finish_manifest_test.go
+++ b/cmd/ailang-worldd/operator_finish_manifest_test.go
@@ -26,6 +26,12 @@ func TestOperatorFinishManifestCapture(t *testing.T) {
 		!strings.HasPrefix(m.Commit.Message, "Reviewed finish\n\n") || m.Commit.Snapshot.Head != f.before {
 		t.Fatalf("signed snapshot and provenance: %+v %v", m, err)
 	}
+	f.gitOutput("-C", f.wt, "add", "b.go")
+	f.gitOutput("-C", f.wt, "-c", "user.name=Fixture", "-c", "user.email=fixture@invalid", "commit", "-m", "Unattributed fixture")
+	m, err = in.frozen(ctx)
+	if err != nil || !reflect.DeepEqual(m.Unattributed, []string{"b.go"}) || len(m.Commit.Snapshot.Paths) != 2 {
+		t.Fatalf("committed base-to-HEAD changes escaped attribution or changed staging: %+v %v", m, err)
+	}
 	before := serviceFileHash(filepath.Join(g.admin, "index"))
 	in.manifest.Base.SHA = strings.Repeat("d", 40)
 	if _, err := in.frozen(ctx); !errors.Is(err, errOperatorFinishReview) || serviceFileHash(filepath.Join(g.admin, "index")) != before {
@@ -34,6 +40,13 @@ func TestOperatorFinishManifestCapture(t *testing.T) {
 }
 
 func TestOperatorFinishWorldEvidence(t *testing.T) {
+	for _, id := range []string{taskOutcomeSemanticID, "test/wrong-semantic-id/v1"} {
+		t.Run(id, func(t *testing.T) { testOperatorFinishWorldEvidence(t, id) })
+	}
+}
+
+func testOperatorFinishWorldEvidence(t *testing.T, semanticID string) {
+	t.Helper()
 	ctx := activationTestContext(t)
 	path := filepath.Join(t.TempDir(), "world.db")
 	world, err := store.Open(path)
@@ -59,7 +72,7 @@ func TestOperatorFinishWorldEvidence(t *testing.T) {
 	}
 	result := taskOutcomeObject(body)
 	ref := hashref.SumSHA256(result.Payload)
-	obj := store.Object{Hash: ref, InterfaceHash: h("interface"), SemanticID: result.SemanticID, Provenance: result.Provenance, Payload: result.Payload}
+	obj := store.Object{Hash: ref, InterfaceHash: h("interface"), SemanticID: semanticID, Provenance: result.Provenance, Payload: result.Payload}
 	if err := world.PutObject(activationTestContext(t), obj); err != nil {
 		t.Fatal(err)
 	}
@@ -74,8 +87,8 @@ func TestOperatorFinishWorldEvidence(t *testing.T) {
 			EntryHash: next.LogHead, TransitionRef: ref}}); err != nil {
 		t.Fatal(err)
 	}
-	if !finishWorldEvidence(ctx, path, r) {
-		t.Fatal("exact committed outcome was not readable")
+	if finishWorldEvidence(ctx, path, r) != (semanticID == result.SemanticID) {
+		t.Fatalf("committed outcome semantic type was misclassified: %s", semanticID)
 	}
 	for _, wrong := range []store.OperatorFinishCheck{
 		{Name: r.Name, Result: r.Result, Ref: r.Ref, Entry: 2},
diff --git a/cmd/ailang-worldd/operator_finish_resume_test.go b/cmd/ailang-worldd/operator_finish_resume_test.go
new file mode 100644
index 0000000..94052ed
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_resume_test.go
@@ -0,0 +1,92 @@
+package main
+
+import (
+	"context"
+	"errors"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorFinishResumeEvidence(t *testing.T) {
+	f, w, _ := operatorLocalFinishFixture(t, nil)
+	ctx := activationTestContext(t)
+	if err := w.check(ctx, w.manifest.Checks[0], 0); err != nil {
+		t.Fatal(err)
+	}
+	w.evidence = func(ctx context.Context, c store.OperatorFinishCheck) bool { return finishWorldEvidence(ctx, f.db, c) }
+	before := operatorFinishJournal(t, w)
+	if err := runEpisodeFinishWithPrep(ctx, w); !errors.Is(err, errOperatorFinishUncertain) ||
+		!reflect.DeepEqual(before, operatorFinishJournal(t, w)) || len(before.Checks) != 1 || before.Checks[0].Result != "passed" ||
+		f.checks != 1 || f.outcomes != 1 || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
+		t.Fatalf("resume trusted a prior pass without matching World evidence: %v", err)
+	}
+}
+
+func TestOperatorCommitFailureInputRecheck(t *testing.T) {
+	f, w, _ := operatorLocalFinishFixture(t, nil)
+	w.commit = func(context.Context, store.OperatorLocalCommit, [][]byte, string) error {
+		seedStartFixtureFile(t, filepath.Join(f.wt, "unexpected"), []byte("changed after reservation"))
+		return errOperatorCommitNoRef
+	}
+	err := runEpisodeFinishWithPrep(activationTestContext(t), w)
+	j := operatorFinishJournal(t, w)
+	if !errors.Is(err, errOperatorFinishUncertain) || j.CommitPhase != "reserved" || j.Commit.State != "" ||
+		len(j.Checks) != 1 || f.checks != 1 || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
+		t.Fatalf("changed inputs were recorded as a definite commit failure: %v %+v", err, j)
+	}
+}
+
+func TestOperatorCommitResumeReserved(t *testing.T) {
+	f, first, now := operatorLocalFinishFixture(t, nil)
+	ctx := activationTestContext(t)
+	b, err := first.inputs(ctx)
+	if err != nil {
+		t.Fatal(err)
+	}
+	spec := first.manifest.Commit
+	tree, err := first.git.planTree(ctx, first.spool, spec, b)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := first.reserve(first.step("commit", tree)); err != nil {
+		t.Fatal(err)
+	}
+	if err := first.commit(ctx, spec, b, tree); err != nil {
+		t.Fatal(err)
+	}
+	reopened, err := store.OpenCredentialsReadOnly(first.manifest.Store)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer reopened.Close()
+	p := first.prep
+	p, err = verifyOperatorFinishPrep(ctx, reopened, p.action, *now, operatorPreparationCAS(p.row), p.row.Owner)
+	if err != nil {
+		t.Fatal(err)
+	}
+	resumed := operatorFinishWork{prep: p, manifest: first.manifest, generation: first.generation, read: reopened,
+		git: first.git, report: first.report, budget: first.budget}
+	head := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
+	index := filepath.Join(first.git.admin, "index")
+	config := filepath.Join(first.git.common, "config")
+	oldIndex, oldConfig := serviceFileHash(index), readStartFixtureFile(t, config)
+	seedStartFixtureFile(t, config, append(oldConfig, []byte("\n[credential]\n\thelper = false\n")...))
+	before := operatorFinishJournal(t, resumed)
+	if err := runEpisodeFinishWithPrep(ctx, resumed); !errors.Is(err, errOperatorFinishUncertain) ||
+		!reflect.DeepEqual(before, operatorFinishJournal(t, resumed)) || before.CommitPhase != "reserved" ||
+		serviceFileHash(index) != oldIndex || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head {
+		t.Fatalf("restart adopted a commit despite unsafe Git config: %v", err)
+	}
+	seedStartFixtureFile(t, config, oldConfig)
+	if err := runEpisodeFinishWithPrep(ctx, resumed); err != nil {
+		t.Fatal(err)
+	}
+	if j := operatorFinishJournal(t, resumed); j.Commit != (store.OperatorFinishLocalCommit{State: store.OperatorLocalCommitDone, SHA: head}) ||
+		head == f.before || serviceFileHash(index) != oldIndex || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head || f.checks != 0 {
+		t.Fatalf("restart failed to adopt the exact reserved descendant read-only: %+v", j)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_route_test.go b/cmd/ailang-worldd/operator_finish_route_test.go
new file mode 100644
index 0000000..12df2b4
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_route_test.go
@@ -0,0 +1,361 @@
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
+	activeID, statusHref := "", ""
+	reportDown := false
+	original := finishClient
+	t.Cleanup(func() { finishClient = original })
+	finishClient = func(string) *client {
+		return memoryDaemon(func(w http.ResponseWriter, r *http.Request) {
+			switch r.URL.Path {
+			case "/v1/episodes/" + start.Episode + "/report":
+				if reportDown {
+					w.WriteHeader(http.StatusServiceUnavailable)
+					return
+				}
+				_ = json.NewEncoder(w).Encode(rep)
+			case "/mcp/":
+				if body := request("GET", statusHref, nil, nil); !strings.Contains(body, "Finish outcome uncertain") {
+					t.Errorf("reserved check was missing from action status: %s", body)
+				}
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
+	var status struct {
+		ID, Href, StatusHref, State, LocalCommit string
+		Checks                                   []struct{ Name, Result string }
+	}
+	request("POST", "/surface/operator-actions", map[string]any{
+		"argv": review.Argv, "operationId": strings.Repeat("4", 32), "catalogueGeneration": review.CatalogueGeneration, "task": "",
+	}, &status)
+	row, err := f.db.OperatorAction(ctx, status.ID, f.now)
+	if err != nil {
+		t.Fatal(err)
+	}
+	card := request("GET", status.Href, nil, nil)
+	if !strings.Contains(card, filename) || !strings.Contains(card, "Prepare local commit") || strings.Contains(card, `id="operator-task"`) {
+		t.Fatal("signed review omitted its pinned manifest or rendered a task")
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
+	activeID, statusHref = row.ID, status.StatusHref
+	run := func() error {
+		turn, stop := f.env.operatorContext(activationTestBudget)
+		defer stop()
+		return f.env.consumeFinish(turn, f.c, f.h, f.src)
+	}
+	before := serviceFileHash(f.creds)
+	held, err := openTerminalEpisodeLock(f.creds, start.Episode, nil, os.O_CREATE, 5*time.Second)
+	if err != nil || held == nil {
+		t.Fatalf("terminal could not acquire the episode turn: %v", err)
+	}
+	blocked := run()
+	if err := held.Close(); err != nil {
+		t.Fatal(err)
+	}
+	if !errors.Is(blocked, store.ErrOperatorConflict) || !strings.Contains(blocked.Error(), row.ID) || before != serviceFileHash(f.creds) {
+		t.Fatalf("lock contention crossed admission: %v", blocked)
+	}
+	var runnerLog syncBuffer
+	refuse := func(name string, want error) {
+		t.Helper()
+		before := serviceFileHash(f.creds)
+		err := run()
+		_, remoteWrite := os.Lstat(trace)
+		if !errors.Is(err, want) || before != serviceFileHash(f.creds) || checks+outcomes != 0 ||
+			!errors.Is(remoteWrite, os.ErrNotExist) || finishRefusal(err) != want.Error() {
+			t.Fatalf("%s changed evidence or lost the remedy before claim: %v", name, err)
+		}
+		fmt.Fprintln(&runnerLog, finishRefusal(err))
+	}
+	launch := f.receipt(t).Receipt
+	savedLaunch, _ := json.Marshal(launch)
+	launch.PID = os.Getpid()
+	birth, live := operatorProcessIdentity(launch.PID)
+	if !live {
+		t.Fatal("live PID positive control unavailable")
+	}
+	launch.Birth = birth
+	liveLaunch, _ := json.Marshal(launch)
+	f.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", liveLaunch)
+	refuse("active agent", errOperatorFinishEligibility)
+	f.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", savedLaunch)
+	startWire, _ := json.Marshal(start)
+	f.exec(t, "DELETE FROM episode_starts WHERE episode=?", start.Episode)
+	refuse("missing start", errOperatorFinishEligibility)
+	f.exec(t, "INSERT INTO episode_starts(episode,start_json) VALUES (?,?)", start.Episode, string(startWire))
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
+	configPath := filepath.Join(repo.repo, ".git", "config")
+	approvedConfig := readStartFixtureFile(t, configPath)
+	seedStartFixtureFile(t, configPath, append(approvedConfig, []byte("\n[credential]\n\thelper = !touch "+shellQuote(trace)+"\n")...))
+	refuse("repository helper", errOperatorFinishReview)
+	seedStartFixtureFile(t, configPath, approvedConfig)
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
+	f.exec(t, "UPDATE session_credentials SET expires_at=? WHERE episode_id=?", f.now, start.Episode)
+	refuse("expired session", errOperatorFinishEligibility)
+	f.exec(t, "UPDATE session_credentials SET expires_at=? WHERE episode_id=?", f.self.ExpiresAt, start.Episode)
+
+	f.env.operatorCatalogue = func(runnerConfig) []byte { return raw }
+	fmt.Fprint(&runnerLog, f.tick(t))
+	row, err = f.db.OperatorAction(ctx, status.ID, f.now)
+	if err != nil {
+		t.Fatal(err)
+	}
+	j, err := f.db.ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))
+	head, headErr := worktree.Output(ctx, repo.git, []string{"-C", f.path(), "rev-parse", "HEAD"})
+	if err != nil || headErr != nil || j.Commit.SHA != strings.TrimSpace(head) || j.Commit.State != store.OperatorLocalCommitDone ||
+		row.V3ID != "" || checks != 1 || outcomes != 1 || len(j.Checks) != 1 || j.Checks[0].Result != "passed" {
+		t.Fatalf("local hand-off: %+v %v head %v, checks %d, outcomes %d", j, err, headErr, checks, outcomes)
+	}
+	if err := run(); !errors.Is(err, store.ErrOperatorMissing) || checks != 1 || outcomes != 1 {
+		t.Fatalf("durable local hand-off replayed: %v", err)
+	}
+	request("GET", status.StatusHref, nil, &status)
+	if row.Owner != f.h.Owner.Nonce || row.Generation == 0 || status.LocalCommit != j.Commit.SHA || len(status.Checks) != 1 ||
+		status.Checks[0].Name != start.Checks[0].Name || status.Checks[0].Result != "passed" ||
+		presetGit(t, repo, "-C", f.path(), "rev-list", "--count", start.Base.SHA+"..HEAD") != "1" {
+		t.Fatalf("runner claim, exact commit count or projected result missing: %+v", status)
+	}
+	var v3 int
+	if err := f.worktreeStartRig.db.QueryRowContext(ctx, "SELECT count(*) FROM verb_requests").Scan(&v3); err != nil || v3 != 0 {
+		t.Fatalf("V3 put crossed local finish: %d %v", v3, err)
+	}
+	world, err := store.OpenReadOnly(f.c.World)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer world.Close()
+	if _, found, err := world.GetRegistryHead(ctx, "world/episode-finish/"+start.Episode); err != nil || found {
+		t.Fatalf("finish record crossed local finish: %v %v", found, err)
+	}
+	output := request("GET", status.StatusHref, nil, nil) + request("GET", "/surface/requests.json", nil, nil) +
+		request("GET", "/surface/home-summary.json", nil, nil) + log.String() + runnerLog.String()
+	for _, secret := range []string{filename, "private-task-fixture-marker", "BEGIN SSH SIGNATURE", f.token, "finish_manifest", j.Checks[0].Ref} {
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
index 0000000..af48805
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_runner.go
@@ -0,0 +1,380 @@
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
+	"syscall"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/checkrun"
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
+	_, err = verifyOperatorFinishPrep(ctx, ro, a, verifyAt, cas, row.Owner)
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
+	w, err := e.finishWork(ctx, c, h, src, ro, fd, row, cas, a, m, generation, observation, verifyAt)
+	if err != nil {
+		return err
+	}
+	return runEpisodeFinishWithPrep(ctx, w)
+}
+
+func (e serviceEnv) finishWork(ctx context.Context, c runnerConfig, h store.OperatorHeartbeat, src operatorStartSource,
+	ro *store.Credentials, fd *os.File, row store.OperatorActionRow, cas store.OperatorCAS, a approval.OperatorAction,
+	m store.OperatorFinishReview, generation int64, observation bool, verifyAt int64) (operatorFinishWork, error) {
+	w := operatorFinishWork{generation: generation, manifest: m, spool: filepath.Dir(e.runnerPaths().Log),
+		read: ro, now: src.clock, budget: e.operatorContext}
+	if observation {
+		var err error
+		w.prep, err = verifyOperatorFinishPrep(ctx, ro, a, verifyAt, cas, row.Owner)
+		if err != nil {
+			return w, err
+		}
+		start, err := ro.EpisodeStart(ctx, m.Episode)
+		if err != nil || start == nil || start.Remote != m.Remote || start.Base != m.Base {
+			return w, errOperatorFinishUncertain
+		}
+		auth, err := planOperatorFinishAuth(c, m.Remote.URL)
+		if err != nil {
+			return w, err
+		}
+		w.git, err = openOperatorFinishGit(auth, filepath.Join(src.workspace, m.Episode), start.Repo, fd)
+		if err != nil || w.git.dir != m.Worktree {
+			return w, errOperatorFinishUncertain
+		}
+		w.report.Start = start
+		return w, nil
+	}
+	if _, err := verifyOperatorFinishPrep(ctx, ro, a, src.clock(), cas, row.Owner); err != nil {
+		return w, err
+	}
+	in, err := e.finishInputs(ctx, c, src, ro, row.Episode, fd)
+	if err != nil {
+		return w, err
+	}
+	if !finishReviewInputsMatch(in.manifest, m) {
+		return w, errOperatorFinishReview
+	}
+	if row.Owner == "" {
+		reviewed, err := in.frozen(ctx)
+		if err != nil {
+			return w, errors.Join(errOperatorFinishReview, err)
+		}
+		if !reflect.DeepEqual(reviewed, m) || len(m.Protected) != 0 || len(m.Unattributed) != 0 {
+			return w, errOperatorFinishReview
+		}
+		admit, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
+		if err != nil {
+			return w, errOperatorDraining
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
+			return w, err
+		}
+		cas.Generation++
+		row.Owner = h.Owner.Nonce
+	}
+	w.prep, err = verifyOperatorFinishPrep(ctx, ro, a, src.clock(), cas, row.Owner)
+	if err != nil {
+		return w, err
+	}
+
+	client := finishClient(c.Address)
+	w.report, w.git = in.report, in.git
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
+	return w, nil
+}
+
+func finishReviewInputsMatch(current, frozen store.OperatorFinishReview) bool {
+	return current.Remote == frozen.Remote && current.Base == frozen.Base && current.PresetDigest == frozen.PresetDigest &&
+		current.Profile == frozen.Profile && current.Sessions == frozen.Sessions && current.Launch == frozen.Launch &&
+		reflect.DeepEqual(current.Checks, frozen.Checks)
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
+				ticks = nil
+				continue
+			}
+			_ = e.touchRunner(c, h) // A busy writer leaves the real heartbeat deadline in force.
+		}
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_runner_test.go b/cmd/ailang-worldd/operator_finish_runner_test.go
new file mode 100644
index 0000000..afb565d
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_runner_test.go
@@ -0,0 +1,58 @@
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
+	for _, phase := range []string{"confirmed", "resume", "terminal"} {
+		t.Run(phase, func(t *testing.T) {
+			var confirm []bool
+			if phase == "confirmed" {
+				confirm = []bool{true}
+			}
+			var change func(*store.OperatorFinishReview)
+			var want error = operatorGitHubAuth(6)
+			wantTrace := "remote-read\n"
+			if phase == "terminal" {
+				change = func(m *store.OperatorFinishReview) { m.Checks[0].Unsandboxed = true }
+				want, wantTrace = errOperatorFinishTerminal, ""
+			}
+			f, work, now := operatorLocalFinishFixture(t, change, confirm...)
+			path := filepath.Dir(work.manifest.Store)
+			tools, err := os.MkdirTemp(path, "approved-")
+			if err != nil {
+				t.Fatal(err)
+			}
+			trace := filepath.Join(tools, "trace")
+			c := runnerConfig{Home: path, Credentials: work.manifest.Store, UID: os.Getuid(), ToolDirs: []string{tools},
+				Executables: make(map[string]string), FinishRemotes: []string{work.manifest.Remote.URL}}
+			writeFinishTestTool(t, &c, tools, "gh", "case \"$1\" in auth|label) ;; repo) printf 'WRITE\\n';; *) exit 83;; esac\n")
+			writeFinishTestTool(t, &c, tools, "git", "printf 'remote-read\\n' >> "+shellQuote(trace)+"\nexit 1\n")
+			before := serviceFileHash(work.manifest.Store)
+			ctx := activationTestContext(t)
+			h, err := work.read.OperatorHeartbeat(ctx)
+			if err != nil {
+				t.Fatal(err)
+			}
+			err = (serviceEnv{}).consumeFinish(ctx, c, h, operatorStartSource{clock: func() int64 { return *now }})
+			row, readErr := work.read.OperatorAction(ctx, work.prep.row.ID, *now)
+			marker, markerErr := os.ReadFile(trace)
+			beforeHead := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
+			j := operatorFinishJournal(t, work)
+			if !errors.Is(err, want) || readErr != nil || !reflect.DeepEqual(work.prep.row, row) ||
+				before != serviceFileHash(work.manifest.Store) || string(marker) != wantTrace || markerErr != nil && !errors.Is(markerErr, os.ErrNotExist) ||
+				beforeHead != f.before || j.CommitPhase != "" || j.OpenCheck || len(j.Checks) != 0 || f.checks != 0 ||
+				!h.Available(h.Build, "start", *now) {
+				t.Fatalf("readiness crossed %s admission: %v, row %s, marker=%q (%v)", phase, err, row.State, marker, markerErr)
+			}
+		})
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_tree.go b/cmd/ailang-worldd/operator_finish_tree.go
index 570f474..231afdc 100644
--- a/cmd/ailang-worldd/operator_finish_tree.go
+++ b/cmd/ailang-worldd/operator_finish_tree.go
@@ -345,6 +345,9 @@ func operatorFinishCommitMatches(raw string, spec store.OperatorLocalCommit, tre
 }
 
 func (g operatorFinishGit) commit(ctx context.Context, spool string, spec store.OperatorLocalCommit, contents [][]byte, tree string, start *store.EpisodeStart) error {
+	if config, err := g.config(ctx, start); err != nil || config != spec.Snapshot.Config {
+		return errOperatorFinishReview
+	}
 	name, email, err := worktree.ParseCommitter(spec.Committer)
 	if err != nil || name == "" || strings.ContainsAny(start.Agent, "\r\n<>") {
 		return errOperatorCommitNoRef
diff --git a/cmd/ailang-worldd/operator_finish_tree_test.go b/cmd/ailang-worldd/operator_finish_tree_test.go
index 7927390..dfe1c88 100644
--- a/cmd/ailang-worldd/operator_finish_tree_test.go
+++ b/cmd/ailang-worldd/operator_finish_tree_test.go
@@ -333,6 +333,10 @@ func TestOperatorFinishConfigReplacement(t *testing.T) {
 	if _, _, err := g.snapshot(ctx, f.rep.Start); !errors.Is(err, errOperatorFinishFence) {
 		t.Fatalf("replacement was accepted: %v", err)
 	}
+	spec := store.OperatorLocalCommit{Committer: "Fixture Operator <fixture@example.invalid>"}
+	if err := g.commit(ctx, f.root, spec, nil, "", f.rep.Start); !errors.Is(err, errOperatorFinishReview) {
+		t.Fatalf("unsafe configuration reached commit staging: %v", err)
+	}
 	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) ||
 		before != serviceFileHash(filepath.Join(g.admin, "index")) || bytes.Equal(body, readStartFixtureFile(t, path)) {
 		t.Fatalf("replacement did not stop before using the live Git configuration: %v", err)
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
index ddd3362..3aa3c7c 100644
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
@@ -142,6 +154,25 @@ func (e serviceEnv) runOperatorPreparationLoop(errOut io.Writer, c runnerConfig,
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
+			case err != nil && !errors.Is(err, store.ErrOperatorMissing):
+				fmt.Fprintln(errOut, finishRefusal(err))
+			}
+		}
 		select {
 		case _, open := <-ticks:
 			if !open {
@@ -153,6 +184,16 @@ func (e serviceEnv) runOperatorPreparationLoop(errOut io.Writer, c runnerConfig,
 	}
 }
 
+func finishRefusal(err error) string {
+	for _, class := range []error{errOperatorFinishReview, errOperatorFinishEligibility, errOperatorFinishTerminal,
+		errOperatorFinishFence, errOperatorFinishCheck, errOperatorFinishUncertain, errOperatorCommitFailed} {
+		if errors.Is(err, class) {
+			return class.Error()
+		}
+	}
+	return "operator finish blocked; inspect action status and durable evidence"
+}
+
 func (e serviceEnv) runnerTick(c runnerConfig, h *store.OperatorHeartbeat) error {
 	if err := e.checkRunnerConfig(c); err != nil {
 		return err
@@ -160,12 +201,18 @@ func (e serviceEnv) runnerTick(c runnerConfig, h *store.OperatorHeartbeat) error
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
@@ -233,7 +280,9 @@ func (e serviceEnv) runnerStatus(w io.Writer) error {
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
index 37d5db3..b3ae699 100644
--- a/host/daemon/surface_operator.go
+++ b/host/daemon/surface_operator.go
@@ -21,7 +21,44 @@ import (
 const operatorPath = "/surface/operator-actions"
 const operatorUnavailable = "runner or reviewed preset unavailable; inspect ailang-worldd service runner status or start at a terminal"
 
-var operatorPolicy = verbs.OperatorPolicy{Start: "A", Finish: "A"} // Headless start is the only installed S7 backend.
+var operatorPolicy = verbs.OperatorPolicy{Start: "A", Finish: "A"}
+
+var errOperatorFinishChanges = errors.New("protected or unattributed changes need terminal finish")
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
+		return h, approval.OperatorAction{}, errOperatorFinishChanges
+	}
+	a, err := entry.Action(h.CatalogueGeneration, in.Ready)
+	return h, a, err
+}
 
 func (d *Daemon) operatorVerb(r *http.Request, argv []string, out *verbAnswer) bool {
 	a := verbs.AdmitOperator(argv, operatorPolicy)
@@ -29,24 +66,30 @@ func (d *Daemon) operatorVerb(r *http.Request, argv []string, out *verbAnswer) b
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
+			case errors.Is(err, errOperatorFinishAuth), errors.Is(err, errOperatorFinishChanges):
+				out.Outcome, out.Reason, out.CLI = "terminal", err.Error(), a.CLI
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
@@ -158,13 +201,13 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
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
@@ -189,8 +232,13 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
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
@@ -204,7 +252,7 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
 	}
 	defer wr.Close()
 	now = d.approveNow()
-	h, entry, err := d.operatorReview(ctx, wr, a.Intent.Preset, now)
+	h, action, err := d.operatorActionForIntent(ctx, wr, a.Intent, body.Task, now)
 	if err == nil && h.CatalogueGeneration != body.CatalogueGeneration {
 		err = store.ErrOperatorRunner
 	}
@@ -212,11 +260,6 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
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
@@ -239,6 +282,8 @@ func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, o
 
 func (d *Daemon) operatorError(w http.ResponseWriter, err error) {
 	switch {
+	case errors.Is(err, errOperatorFinishAuth), errors.Is(err, errOperatorFinishChanges):
+		writeApprovalRefusal(w, http.StatusConflict, "RunnerUnavailable", err.Error())
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
index 0000000..852214d
--- /dev/null
+++ b/host/daemon/surface_operator_finish_test.go
@@ -0,0 +1,124 @@
+package daemon
+
+import (
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
+	card := f.do("GET", requested.Href, nil, nil)
+	if card.Code != 200 || !strings.Contains(card.Body.String(), manifest.Report) || strings.Contains(card.Body.String(), `id="operator-task"`) {
+		t.Fatal("signed card did not display exact finish review without a task")
+	}
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
+	for _, output := range []string{f.do("GET", requested.StatusHref, nil, nil).Body.String(), f.log.String()} {
+		if strings.Contains(output, manifest.Report) || strings.Contains(output, "authenticatorData") || f.row(row.ID).V3ID != "" {
+			t.Fatal("confirmed status exposed prep data or issued a V3 prompt")
+		}
+	}
+}
diff --git a/host/daemon/surface_operator_read.go b/host/daemon/surface_operator_read.go
index 47ad635..04185de 100644
--- a/host/daemon/surface_operator_read.go
+++ b/host/daemon/surface_operator_read.go
@@ -14,7 +14,15 @@ import (
 	"github.com/sunholo-data/ailang-world/host/verbs"
 )
 
+type operatorCheckStatus struct {
+	Name   string `json:"name"`
+	Result string `json:"result"`
+}
+
 type operatorStatus struct {
+	Checks      []operatorCheckStatus `json:"checks,omitempty"`
+	LocalCommit string                `json:"localCommit,omitempty"`
+
 	V            int    `json:"v"`
 	ID           string `json:"id"`
 	Episode      string `json:"episode"`
@@ -62,7 +70,7 @@ var operatorCard = template.Must(template.New("operator action").Parse(`<!doctyp
 <h1 id="operator-title" tabindex="-1">Operator action {{.Status.ID}}</h1>
 <p id="operator-state" data-state="{{.Status.State}}">{{.Status.State}}</p>
 <p id="operator-progress" role="status">Status at {{.Status.Now}}{{if .Status.Reason}} · {{.Status.Reason}}{{end}}</p>
-{{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre><h2>Task</h2><pre id="operator-task">{{.Task}}</pre>{{end}}
+{{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre>{{end}}{{if .Task}}<h2>Task</h2><pre id="operator-task">{{.Task}}</pre>{{end}}
 <section id="operator-controls" hidden aria-label="Confirm this action">
 <button id="operator-key" type="button" hidden>Confirm with passkey</button>
 <details><summary>Confirm with SSH signature</summary>
@@ -117,6 +125,31 @@ func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id s
 			return
 		}
 	}
+	if status.Href != "" && row.Kind == "finish" && row.Owner != "" {
+		j, err := ro.ReadOperatorFinishJournal(ctx, store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation})
+		if err != nil {
+			d.operatorError(w, err)
+			return
+		}
+		for _, check := range j.Checks {
+			status.Checks = append(status.Checks, operatorCheckStatus{Name: check.Name, Result: check.Result})
+		}
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
index 0481d37..f179af6 100644
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
@@ -1584,7 +1586,7 @@
 
   if (typeof module !== 'undefined') module.exports = { needsYouHTML, decidedHTML, notifications, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, finishMatchesNode, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
     adaptModel, projectHeadsModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, filePath, fileEvidenceHTML, filesTabHTML, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW },
-    VERB_ROUTE_BUDGET_MS, verbOutcome, askVerb, verbAct, verbPanel, goWith, addressGo, operatorDraft, operatorRequestControl, operatorCardHref };
+    VERB_ROUTE_BUDGET_MS, verbOutcome, askVerb, verbAct, verbPanel, goWith, addressGo, operatorDraft, operatorRequestControl, operatorCardHref, showOperatorReview };
   if (typeof document === 'undefined') return;
 
   // ================= The page =================
@@ -3532,27 +3534,31 @@
   }
   let operatorStorage;
   try { operatorStorage = window.sessionStorage; } catch (_) { /* request stays inert */ }
-  function showOperatorReview(panel, current) {
-    const draft = operatorDraft(operatorStorage);
+  function showOperatorReview(panel, current, ui = {
+    doc: document, show: showVerbPanel, storage: operatorStorage, fetch: (...args) => fetch(...args),
+    random: b => crypto.getRandomValues(b), timeout: ms => AbortSignal.timeout(ms), assign: href => location.assign(href),
+  }) {
+    const draft = operatorDraft(ui.storage);
     if (draft) panel = { ...draft.request, body: draft.body };
-    const el = showVerbPanel({ ...panel, head: 'Review ' + panel.argv[3] });
-    const body = el.querySelector('.verb-body'), area = document.createElement('div');
+    const el = ui.show({ ...panel, head: 'Review ' + panel.argv[3] });
+    const body = el.querySelector('.verb-body'), area = ui.doc.createElement('div');
     area.className = body.className;
     area.textContent = body.textContent;
     body.replaceWith(area);
-    const label = document.createElement('label'), task = document.createElement('textarea');
-    const request = document.createElement('button'), info = document.createElement('p');
+    const label = ui.doc.createElement('label'), task = ui.doc.createElement('textarea');
+    const request = ui.doc.createElement('button'), info = ui.doc.createElement('p');
     task.id = 'operator-task'; task.rows = 3; task.cols = 26;
     label.textContent = 'Task (1–8192 UTF-8 bytes)'; label.htmlFor = task.id;
     request.id = 'operator-request'; request.type = 'button';
     info.id = 'operator-request-info'; info.setAttribute('role', 'status');
-    area.append(label, document.createElement('br'), task, document.createElement('br'), request, info);
+    const finish = panel.argv[2] === 'finish';
+    if (!finish) area.append(label, ui.doc.createElement('br'), task, ui.doc.createElement('br'));
+    area.append(request, info);
     const live = () => current() && el.isConnected;
-    const control = operatorRequestControl(panel, {storage: operatorStorage, current: live, fetch: (...args) => fetch(...args),
-      random: b => crypto.getRandomValues(b), timeout: ms => AbortSignal.timeout(ms)});
+    const control = operatorRequestControl(panel, {...ui, current: live});
     const repaint = () => {
       task.readOnly = !!control.request;
-      request.textContent = control.request ? 'Reconnect' : 'Request action';
+      request.textContent = control.request ? 'Reconnect' : finish ? 'Request approval' : 'Request action';
     };
     task.value = control.request?.task || '';
     repaint();
@@ -3561,14 +3567,14 @@
       task.readOnly = true;
       await control.submit(task.value).then(result => {
         if (!live()) return;
-        if (result.href) return location.assign(result.href);
+        if (result.href) return ui.assign(result.href);
         info.textContent = result.error || '';
       });
       repaint();
       request.disabled = false;
       if (live()) (task.readOnly ? request : task).focus();
     });
-    (task.readOnly ? request : task).focus();
+    (finish || task.readOnly ? request : task).focus();
   }
   const outstandingOperatorRequest = operatorDraft(operatorStorage);
   if (outstandingOperatorRequest) {
diff --git a/host/surfacepage/app_operator_test.go b/host/surfacepage/app_operator_test.go
index 183b1e3..47440fb 100644
--- a/host/surfacepage/app_operator_test.go
+++ b/host/surfacepage/app_operator_test.go
@@ -186,3 +186,39 @@ for (const [k,v] of Object.entries({v:'1',episode:'other',kind:'finish',state:'i
   assert.equal(pick({...result,[k]:v}),'',k);
 `)
 }
+
+func TestOperatorFinishBarReview(t *testing.T) {
+	runOperatorBar(t, `
+const finishArgv = ['ailang-worldd','episode','finish','s8','--preset','p','--confirm-on-surface'];
+for(const bad of [finishArgv.slice(0,6), [...finishArgv.slice(0,6), "--dry-run"], finishArgv.concat('--title'), finishArgv.concat('--ready','--ready')])
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
+const mounted = [], clicked = [];
+let focus, link;
+const el = { isConnected:true, querySelector:()=>({replaceWith(){}}) };
+const fakeDoc = {createElement:tag=>({tag, append(...children){ mounted.push(...children); },
+  setAttribute(){}, focus(){ focus=this; }, addEventListener(name,fn){ if(name==='click') clicked.push(fn); }})};
+app.showOperatorReview(admitted, () => true, {...io, doc:fakeDoc, show:()=>el, assign:x=>{link=x;}});
+assert.equal(mounted.some(n=>n.tag==='textarea'||n.tag==='label'),false);
+assert.equal(focus.textContent,'Request approval');
+assert.equal(posts.length,2);
+await clicked[0]();
+assert.equal(posts.length,3);
+assert.equal(posts[2].task,'');
+assert.equal(focus.id,'operator-request');
+assert.equal(link,href);
+`)
+}
diff --git a/host/surfacepage/app_test.go b/host/surfacepage/app_test.go
index f9d3f0c..eb11cce 100644
--- a/host/surfacepage/app_test.go
+++ b/host/surfacepage/app_test.go
@@ -167,7 +167,7 @@ const stacks = JSON.parse(fs.readFileSync('stacks.json', 'utf8'));
 const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
 assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'applyDelta', 'barHTML', 'blendXS', 'cleanTitle', 'crumbsHTML', 'decidedHTML', 'esc', 'fetchCoherent', 'finishMatchesNode', 'filePath', 'fileEvidenceHTML', 'filesTabHTML', 'frontHTML',
   'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'notifications', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'projectHeadsModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict',
-  'VERB_ROUTE_BUDGET_MS', 'askVerb', 'verbAct', 'verbOutcome', 'verbPanel', 'goWith', 'addressGo', 'operatorCardHref', 'operatorDraft', 'operatorRequestControl'].sort());
+  'VERB_ROUTE_BUDGET_MS', 'askVerb', 'verbAct', 'verbOutcome', 'verbPanel', 'goWith', 'addressGo', 'operatorCardHref', 'operatorDraft', 'operatorRequestControl', 'showOperatorReview'].sort());
 assert.equal(app.operatorCardHref(null, 'ep'), '');
 
 assert.equal(app.serviceLineHTML(undefined), '');
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index dde5280..96a0dfe 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -2184,3 +2184,22 @@ choosing **Confirm with passkey**. The card shows the recorded confirmation and
 Close the page and reopen the same card URL to read the outcome and **Episode report**. Use
 **Refresh status** to reconnect if status reads fail. **Back to World** returns to the World surface.
 If the request response was lost, **Reconnect** in the review panel retries the same frozen request.
+
+### Prepare a local finish in the World app
+
+Use the fake finish preset with its explicit committer and fake GitHub URL.
+The fake Git tools read a bare local repository. With a stopped launch and
+one World-attributed worktree change, enter this fixture command in the bar.
+
+```text
+episode finish ep --preset p
+```
+
+Enter opens the review. Select **Request approval**, review the held inputs on
+the card and confirm with the enrolled fixture passkey or SSH key. Use
+**Action status (JSON)** to read check outcomes and the retained local commit.
+The card reports that PR confirmation has yet to be issued.
+
+If the runner shows `finish blocked: github-auth`, repair its approved PATH,
+GitHub login or Git credential helper in the attended runner setup. The start
+flow remains available.
````

### M2 — S9d CLI mirror. Planner rehearsal only

````diff M2
diff --git a/cmd/ailang-worldd/operator_finish_route_test.go b/cmd/ailang-worldd/operator_finish_route_test.go
index 12df2b4..522bb7e 100644
--- a/cmd/ailang-worldd/operator_finish_route_test.go
+++ b/cmd/ailang-worldd/operator_finish_route_test.go
@@ -345,7 +345,12 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 	if _, found, err := world.GetRegistryHead(ctx, "world/episode-finish/"+start.Episode); err != nil || found {
 		t.Fatalf("finish record crossed local finish: %v %v", found, err)
 	}
-	output := request("GET", status.StatusHref, nil, nil) + request("GET", "/surface/requests.json", nil, nil) +
+	var cli, cliErr syncBuffer
+	if rc := runServiceRunner([]string{"finish-status", status.ID}, &cli, &cliErr, f.env); rc != exitOK ||
+		!strings.Contains(cli.String(), j.Commit.SHA) || !strings.Contains(cli.String(), "passed") {
+		t.Fatalf("finish-status did not use the redacted receipt: %d %s %s", rc, &cli, &cliErr)
+	}
+	output := cli.String() + cliErr.String() + request("GET", status.StatusHref, nil, nil) + request("GET", "/surface/requests.json", nil, nil) +
 		request("GET", "/surface/home-summary.json", nil, nil) + log.String() + runnerLog.String()
 	for _, secret := range []string{filename, "private-task-fixture-marker", "BEGIN SSH SIGNATURE", f.token, "finish_manifest", j.Checks[0].Ref} {
 		if strings.Contains(output, secret) {
diff --git a/cmd/ailang-worldd/operator_finish_runner.go b/cmd/ailang-worldd/operator_finish_runner.go
index af48805..569579d 100644
--- a/cmd/ailang-worldd/operator_finish_runner.go
+++ b/cmd/ailang-worldd/operator_finish_runner.go
@@ -5,6 +5,7 @@ import (
 	"encoding/json"
 	"errors"
 	"fmt"
+	"io"
 	"os"
 	"path/filepath"
 	"reflect"
@@ -352,6 +353,41 @@ func finishReviewInputsMatch(current, frozen store.OperatorFinishReview) bool {
 		reflect.DeepEqual(current.Checks, frozen.Checks)
 }
 
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
 // An admitted check may outlast the heartbeat's admission window. The runner
 // keeps the same pinned catalogue alive while the exclusive finish turn runs.
 func (e serviceEnv) runFinishTurn(c runnerConfig, h *store.OperatorHeartbeat, src operatorStartSource,
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
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 96a0dfe..a64b152 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -2198,7 +2198,9 @@ episode finish ep --preset p
 Enter opens the review. Select **Request approval**, review the held inputs on
 the card and confirm with the enrolled fixture passkey or SSH key. Use
 **Action status (JSON)** to read check outcomes and the retained local commit.
-The card reports that PR confirmation has yet to be issued.
+The card reports that PR confirmation has yet to be issued. The same
+receipts are available from `ailang-worldd service runner finish-status ACTION-ID`.
+Replace `ACTION-ID` with the action ID shown on the card.
 
 If the runner shows `finish blocked: github-auth`, repair its approved PATH,
 GitHub login or Git credential helper in the attended runner setup. The start
````

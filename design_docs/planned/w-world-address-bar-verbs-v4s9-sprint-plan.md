# w-world-address-bar-verbs-v4 — S9 sprint plan

## 0. Scope and outcome

Row 243, iteration 448, sprint lane `codex:gpt-6.2-sol`. Design r1, quorum r2 PROCEED 3/3, absent seats []. The reviewed §8 split was 350 + 550 raw, 900 → 1,215 calibrated. The ledger resolves D-WORLD-119 and D-WORLD-120 to A. E2 W33 is read-level PASS and S8 is landed. D-WORLD-122 retains E5’s operator-chosen preset for S11.

**Outcome B.** The compiled whole spans 2799 changed lines. S9 selects M1 and M2, 1399 changed lines, at a boundary where the installed runner still has no caller of the private entry, no finish claim and no commit path. M3 in Appendix P is the S9b continuation, 1400 incremental lines. Apply M3 only in that continuation after its listed gates. M2 can be bisected independently because its manifest, fence and journal codecs have fixtures and its local Git primitive is unreachable from the service.

S9b is the remainder of reviewed S9, through §4.3 steps 1–2. The compiled M3 candidate runs the E2-equivalent read probe before claim, holds the episode lock, verifies against the current enrolled key and claims the exact frozen action. It runs the pinned sandboxed checks under individual reservations and makes one bounded local commit. The durable hand-off is an owned action with the observed typed `operator_action_steps` commit receipt `{state:"local-commit-done",sha:"…"}`, matching observed check receipts and `v3_id` empty. Retain the signed manifest, action generation/owner and staged tree OID in the receipt identity. S10 consumes those facts before reserving and issuing a V3 `VerbRequest`, obtains and reads back the linked key decision, then freezes the exact V3 review, link and Chrome approval. S11 owns push, label creation, PR creation/lookup, authoritative finish recording, detach and revoke. The only GitHub label command used by the S9 readiness path is E2's read-only `gh label list`; there are zero calls into the label-creation, push, PR or finish-record paths.

S9b reuses S2's finish envelope, S5/S7's closed admission and S7b's same-origin key-confirmation handler. It adds a taskless prep branch to the S8 bar and reuses the standalone card's existing finish-kind support. Its acceptance drives the route, sealed CLI/runner boundaries, store and executable Node functions. V3 review/link and browser approval remain S10's Chrome scope.

## 1. Measurements

### 1.1 Base and provenance

- `git rev-parse HEAD` is `9a40fe68e45cd9b5c787aecd716b047327bddb22`. Initial status was clean. `git diff --stat 85ed2234 9a40fe68 -- cmd host docs app` was empty. The selected product is measured with rename detection against independent base and cumulative exports. The controller measures that same range on committed HEAD with `git diff -M --numstat 9a40fe68 HEAD -- cmd host docs app`.
- D-WORLD-119/120 are RESOLVED at `design_docs/world-mission.md:1075–1076`, attended ledger commit `82f05004`. W33 records rc 2 with the floor PATH and PASS with Homebrew PATH, ADMIN and a noninteractive HTTPS Keychain read. No real E2 rerun occurred. S9's fake tools assert the argv and clean env, and forward the read to a bare repository made in `t.TempDir()`.
- S3's journal extension is schema 10. HEAD's overall credentials schema is **11** (`host/store/credentials.go:52`), unlike W15's historical 9 and the handoff shorthand “schema 10”. `operatorSchemaVersion` is 10 at `host/store/operator_runner.go:12`. Existing slots for checks and commit already support the §3.3 receipts, so no DDL, version bump or .ail edit is needed. Fresh and upgraded stores retain their exact schema and writer boundary. S10's V3 linkage remains a separate receipt.
- S8's card already accepts `kind: "finish"` (`host/surfacepage/operator/card.js:12`). Its app bar and S7/S7b daemon admission still require headless start. Reuse the card and generalize the two closed admission shapes in M3.

### 1.2 Guard and route inventory

| Guard or seam | Measured at base | Selected and deferred treatment |
|---|---|---|
| W16/W37 public entry | `episode_finish.go:59,63`; F1 after argument parsing at :92. Pure help exits before F1. | M1 opens F1 before parsing effectful invocations, including unknown flag markers. Standalone flag help keeps its existing no-IO contract. Typed prompt golden retains exact bytes. |
| W16 checks and local commit | Fresh unsandboxed guard :281, check fork :297–308, outcome writer :320, local `CommitWithCommitter` :378. Sandbox route is `finishSandboxCheck` :696. | M2 local Git plumbing works through the existing Git seam. M3 checks use the existing daemon workspace-exec path and outcome writer. An explicit approved Git index and expected tree are used for the local commit. |
| W17/W20 auth/dispatch | Public `Prepare` :417 after commit/body, reuse :435, typed/surface prompt :439–460, send :470, record :482, detach/revoke :487–490. `episode_finish_dispatch.go:51` Prepare does lookup + label read, :67 remote-head, :82 ancestor, :91 push, :95 label create, :100 PR view, :123 PR create. | M1 reuses checkrun and worktree's bounded launch sites for the isolated E2 Prepare/remote read. M3 invokes that probe before lock/claim/effects, on fresh and resume. |
| W15/W18 V3 waiter | `store/verb_requests.go:121,239` put/consume. `episode_finish_surface.go:53,75,168,188,214,252` argv, canonical request, put, wait, poll, decision readback. `authority/judge_verb.go:107` rechecks exact own bytes. | Live V3 path remains terminal-gated through the same helper. M3 returns before the first V3 put. Reviewed V3 byte and state extensions belong to S10. |
| W19 resume | `episode_finish_resume.go:19,79,118,133,148` mode, saved body, reusable approval, detach, revoke. Existing PR view only accepts OPEN and lacks exact identity proof. | Existing typed behavior and its tests stay active. M3 does not accept a reusable V3 approval as prep authority. |
| W28 delayed record | `host/daemon/episode_finish.go:231,234,243,276` sets the used and episode heads and evaluates scope expiry. | S11 owns the exact late-record CAS/reconciliation. S9b has no finish record or cleanup call. |
| Exclusion | `terminalEpisodeLock` already used by terminal finish :144, shared lock `operator_episode_lock.go` uses `LOCK_NB` and a 4 s bound for owner read. | Reuse `TestOperatorTerminalExclusion` and require active action id and remedy under contention. The same descriptor is lent to bounded local Git effects in M3. |
| Store, boundary and policy | `TestProductionContextRoots`, `TestCredentialsMethodsRequireADeadline`, `TestNoNewDeadlineFreeStoreReads`, `TestStoreDeadlineGuardDynamic`, `TestCLICredentialWriterSites`, `TestDaemonCredentialWriterSites`, `TestBoundaryASTWriteGuard`, `TestWorldBoundaryDependencyAllowlist`, `TestSessionRootBudgets`, `TestOperatorProductionDeadlines`. | No production context root is added. Derive individual budgets through the existing root seam. Existing inventories plus named acceptance run on all stages. |
| Subprocesses | Broker AC10 includes `host/worktree/worktree.go` and `host/checkrun/checkrun.go`. Verifygate has `TestSubprocessSinkGate`. | Both shared sites keep childenv credential scrubbing. Auth and Git data plumbing add no launch site. Existing package guards and dedicated child trace observations are required. |
| Privacy and answers | Needs-you selects id, kind, episode and times, standalone card requests details, routine list has no artifact projection. | Assert the task, full manifest and proof absent from routine lists, status JSON and logs. Standalone private review shows signed review bytes. No existing verb golden is updated. |

### 1.3 Baseline and race budgets

CI run 46106233909 on `85ed2234` reports race times of cmd **451.877 s**, daemon **484.675 s**, broker **323.388 s** and verifygate **403.341 s**. The package cap is 900 s. New named cmd tests sum to at most 15 s and new daemon tests to at most 12 s. A whole-package delta is a separate measurement against the same comparable baseline.

Read-only controller file `$STATE/base/legs.txt` is updated as legs finish. At the initial read, store, authority and agentrun plain/race, verifygate and verbs contained zero observed failures. Broker plain recorded 6 base failures, listed verbatim in `$STATE/base/broker-plain.fails`. The scripts record all head legs first and compare by test names, validate the parser control and refuse a truncated package or a nonzero rc with an empty fail set. `cmd-race` uses the exact 11-test skip selector below in every local race leg. Remaining baseline and final timings are recorded in §4.

### 1.4 Production callers

| Call | Base | M1 | M2 | S9b M3 |
|---|---:|---:|---:|---:|
| `RequestOperatorAction`, `JudgeOperatorAction`, store `ConfirmOperatorAction` | 1 each | 1 | 1 | 1 |
| `ClaimOperatorAction` | 1 | 1 | 1 | 2, the new call belongs to the finish worker |
| `runEpisodeFinishWithPrep` | 0 | 0 | 0 | 1, from `consumeFinish` |
| Attended `worktree.CommitWithCommitter` / runner local commit seam | 1 / 0 | 1 / 0 | 1 / 0 | 1 / 1 |
| `worktree.Push`, dispatch `Push`, `CreateLabel`, `ViewPR`, `CreatePR` | 1 each | 1 | 1 | 1 |
| V3 store `PutVerbRequest` | 1 | 1 | 1 | 1 |

`callers.sh` prints matched lines, checks exact counts and proves the instrument observes inherited `approval.ConfirmOperatorAction`, runner dispatch and the positive terminal remote-dispatch call. Direct push, label create and PR commands remain confined to the inherited dispatch seam. Local Git readback and a single ref CAS use the same bounded Git subprocess site.

## 2. Milestones

| Id | Scope | Cumulative numstat, cmd host docs app | Unfiltered product |
|---|---|---:|---:|
| **M1** | Public F1 parity, private proof validation and attended typed install readiness. Host subprocess seams take explicit scrubbed environments. | 621 (+599/−22) | 621 |
| **M2** | Frozen manifest types, bounded no-follow worktree digest capture, local Git staged tree and exact readback primitive, existing-schema journal and test fixtures. No executor path from the runner yet. | 1399 (+1374/−25) | 1399 |
| **M3** (S9b) | Two-kind inert request and confirm reuse, runtime auth and claim, private sandboxed-check journal, commit reservation and reconciliation, redacted status and taskless action review. | 2799 (+2725/−74), 1400 incremental | 2799 |

Cumulative sets are `$S/st-M1/`, `$S/st-M2/`, `$S/st-M3/` / status lists and patch receipts under the same scratch. The filtered count below and the unfiltered product count include insertions plus deletions. The two plan artifacts add 6045 lines outside the product filter, so the selected eventual unfiltered total including both artifacts is 7444. No stage modifies `tools/launchd`, checked-in manifests, .ail or the repository index.

M1 returns typed readiness classes 2–6. An installed configuration supplies the approved GitHub remotes, executables and noninteractive login environment. Read class 2 means the resolvable tools are absent or fail pinning, 3 is login, 4 is label **read**, 5 is repository permission read/WRITE–MAINTAIN–ADMIN and 6 independently authenticates Git `ls-remote --symref` against the pinned URL. Git config is isolated; trusted fixed arguments pin BatchMode, disabling prompts, approved SSH socket and the selected `gh` or Keychain HTTPS helper. The positive `ls-remote` fixture uses a local bare repository via a scoped FAKE git executable. The prep verifier compares phase, owner, episode, digest, generation, lifetime and exact canonical bytes, then independently reads the enrollment and artifact. The public CLI has no argv or env marker for this verifier.

M2 signs source report identity, episode/preset/profile/session and launch evidence, branch/HEAD/index and repo-config digests, exact pinned URL/base, the checks and provenance, sorted changed paths with modes/digests, attribution and protected path sets, and exact commit message/trailer and committer. Dirty traversal is bounded and neither follows symlinks nor reads special files. Git can omit an untracked FIFO from `ls-files --others` (observed in the first fence run), hence the bounded filesystem walk. Do the tentative tree build in private object/index storage, then reserve the action/check or action/tree under CAS before a real effect. Construct the commit from the exact staged tree and expected parent, then CAS the full branch ref. The receipt checker compares the parent, tree, complete message/trailer, committer and staged index. Keep divergence and partial staging without reset. Schema 10's existing step table stores the tree qualified by action id; the action digest binds the rest of the commit intent.

M3 performs all read-only eligibility/approval checks before claim and revalidates each effect boundary. A recorded check result can be adopted only against content-addressed world evidence. Unknown prior output stays fenced. Expiry prohibits any new check/commit even after a successful check; the read-only commit reconciliation path can observe the exact descendant of an already reserved operation. This slice parks checks when a successor authorization would be needed. The S10 consumer must explicitly carry the matching receipts when implementing a successor review. Never infer success from a stored trailer alone, never reconstruct a new commit on recovery, and never enter V3 remote dispatch from this stage.

## 3. Operator output

M1 adds optional attended install setup fields `finish_remotes`, `finish_ssh_socket` and `finish_https_helper`. A failed probe displays `finish blocked: github-auth (readiness step N)` and the attended PATH/login/helper remedy. Companions can still install and start. M2 has no new controls.

S9b adds the taskless prep review panel for the canonical finish tuple and the existing **Request action** button. Posted requests return the existing exact same-origin action and status links. The existing SSH/passkey controls confirm the prep envelope. Status adds `checks` and `localCommit` and projects a retained commit with PR confirmation still to be issued. `service runner finish-status ACTION-ID` shows the check results, retained SHA or uncertainty; routine service status shows finish readiness. `docs/QUICKSTART.md` includes the temporary-fixture bar input. `TestOperatorFinishAdmission` lifts that input verbatim and runs its GET and explicit request/confirm with the existing passkey fixture against the real mux, and `TestOperatorFinishBarReview` executes the corresponding app JavaScript. Full browser approval of the V3 card is S10.

## 4. Acceptance, timing and fail sets

The checked-in race ceilings are 900 s per package. The narrow named-test cap is an additional gate. `accept.sh` sums the Go-reported top-level timings, prints `RACE_BUDGET`, and fails the race mode if cmd exceeds 15 s or daemon exceeds 12 s. A shared-rig overshoot remains an overshoot even when every test assertion passes.

| New cmd test | Final M1 race s | Final M2 race s |
|---|---:|---:|
| `TestOperatorFinishAuth` | 7.57 | 6.27 |
| `TestOperatorFinishAuthInstall` | 0.44 | 0.94 |
| `TestOperatorFinishF1` | 0.22 | 0.19 |
| `TestOperatorFinishTTYParity` | 0.83 | 0.90 |
| `TestOperatorFinishNoProof` | 0.07 | 0.08 |
| `TestOperatorFinishWorktreeFence` | — | 1.13 |
| `TestOperatorCommitIdentity` | — | 3.22 |
| **Named total** | **9.13** | **12.73** |

`$S/delivery-race/M1.raw` and `M2.raw` contain `ACCEPT_OK` and the instrument's `RACE_BUDGET` markers for these exact delivered bytes. Earlier contended M2 samples were 21.43 s and 28.48 s. The final stage kept the local bare-repo Git read and used direct local fixture setup and file observations for the F1 snapshot. The unsandboxed pre-executor recheck is still required.

| Touched package | Pristine race seconds | M1 / delta | M2 / delta |
|---|---:|---:|---:|
| `host/checkrun` | 2.153 | 2.231 / +0.078 | 1.976 / −0.177 |
| `host/worktree` | 11.619 | 12.705 / +1.086 | 8.986 / −2.633 |
| `host/store` | 40.150 (controller) | — | 86.621 / +46.471 |
| `cmd/ailang-worldd` | UNMEASURED locally, CI 451.877 | UNMEASURED whole | UNMEASURED whole |
| `host/daemon` | CI 484.675 | 0 new tests | 0 new tests |

### Deferred S9b timing

| New cmd test | Final race s | Prior contended race s |
|---|---:|---:|
| `TestOperatorFinishPrep` | 0.27 | 2.30 |
| `TestOperatorCommitReconcile` | 0.94 | 14.08 |
| `TestOperatorCommitFailed` | 0.78 | 14.20 |
| `TestOperatorFinishChecks` | 1.56 | 34.50 |
| `TestOperatorFinishUnsandboxed` | 0.26 | 5.38 |
| `TestOperatorFinishPrepExpiry` | 0.49 | 18.50 |
| `TestOperatorFinishStale` | 0.42 | 3.43 |
| `TestOperatorFinishReadiness` | 1.26 | 3.49 |
| `TestOperatorFinishEligibility` | 0.26 | 2.23 |
| **Named total** | **6.24** | **98.11** |

`$S/delivery-deferred-race/sealed/accept.raw` prints `RACE_BUDGET cmd/ailang-worldd n=9 added=6.24s cap=15s` and `ACCEPT_OK M3 race`. A slower run passed the same nine assertions and totalled 98.11 s (`$S/delivery-deferred-race/cmd.raw`). A separate truncated run supplies no timing receipt. The outside timing gate remains required. `TestOperatorFinishBarReview` ran under race in 0.07 s. `TestOperatorFinishAdmission` could not bind a listener in the sandbox. Its measured route evidence and daemon time remain outside requirements.

Package times are actual command output from the controlled scratch sweep `$S/local-race-comparables` or the controller recording. The cmd whole-package race time and delta require its 11-test-skipped comparable base run. Future cmd timing must meet the S9b budget before activation. The complete local `host/surfacepage` race package passed on pristine (6.856 s) and exact M3 (7.857 s, +1.001 s). Plain also passed on both. The daemon delta requires the M3 stage/baseline gates.

PENDING — controller to measure the entire outside legs with `outside.sh` and `gates.sh` in §7. Worktree/checkrun appeared with empty sets locally. The controller completed all 12 base legs. Cmd plain has 13 failures, daemon plain has 4 and daemon race has 14. The outside daemon race baseline took 1,612.938 s on this loaded rig. CI took 484.675 s. No controller fail set was rewritten.

No wall-clock wait is used as an assertion. Clock, expiry and write budgets use fixture controls. The 1777 run passes an invocation-local TMPDIR under a 1777 parent and compares physical paths. All fakes observe their own temp paths. Outside launchd-shaped fakes run with the same fixture controller and never invoke real launchctl. There is no accepted Chrome result for S9 because no Chrome leg is claimed; optional carried browser pins are described below.

## 5. Branch traces, test ownership and mutation matrix

### 5.4 Finish, F1 and remote reconciliation

| Reviewed row | Stage and named pin | Behavioural change under test |
|---|---|---|
| Public typed finish | M1 `TestOperatorFinishTTYParity`, `TestFinishPromptIsTheGolden` | TTY-prompt strips exact legacy display; bytes and typed y are observed. |
| Public no tty, flags/env | M1 `TestOperatorFinishF1` | F1-surface/env/flag. Observe ref/HEAD/index/body/credentials hash, daemon requests, check/outcome calls and Git-seam count. |
| Private missing/forged/wrong/expired/replayed proof | M1 `TestOperatorFinishNoProof`, M3 `TestOperatorFinishPrep` | P-artifact/recheck/owner/generation/expiry, J-manifest. Real migration and SSH test signer, a valid proof is the control. Wrong manifest and episode stay read-only. |
| Valid prep and local commit | M3 `TestOperatorFinishPrep`, `TestOperatorCommitReconcile`, outside `TestOperatorFinishAdmission` | J-early-commit and R-prep-card. End at typed local receipt, no V3 request. |
| No start/degraded/empty/running/expired session/read failure | M3 `TestOperatorFinishEligibility` and the existing S6 fake launch controls | The isolated no-start, bad-id and read-error cases verify unchanged durable state. PENDING — complete runtime fixture matrix for degraded, empty, live launch and expired/changed session. |
| Stale preset, HEAD, untracked contents, path race, remote/base/profile | M2 `TestOperatorFinishWorktreeFence`, M3 `TestOperatorFinishStale` | J-bytes detects untracked or HEAD drift before any effect. Ordered path/file snapshot and config gate. PENDING — measure separate base, preset/profile and path-race mutants at the runner. |
| Protected/unattributed and Git boundary | M2 `TestOperatorFinishWorktreeFence`, `TestOperatorCommitIdentity`, M3 route/frozen manifest | G-hookconfig/filter/helper/monitor/fifo, G-hooks. Exact fence refusal, real hook sentinel absent. PENDING — add the independent origin/protected attribution arms and the worktree-config replacement race at full runtime. |
| Unsandboxed including resume | M3 `TestOperatorFinishUnsandboxed` | J-unsandboxed. Observe same journal and zero checks both with zero and existing reserved work. |
| Missing tools/login/Git auth | M1 `TestOperatorFinishAuth`, `TestOperatorFinishAuthInstall`; M3 `TestOperatorFinishReadiness` | E2-* and R-auth. Install probe precedes fixture writes. Runtime authenticated gh still fails Git before a claim on confirmed or resume. |
| Failing/aborted/timed-out/mutating check | M3 `TestOperatorFinishChecks` | J-fail-check/J-unknown. Preserve observation or reservation; no commit or blind retry. |
| Definite commit failure | M3 `TestOperatorCommitFailed` | J-no-ref. Initial refs and check evidence survive, no second attempt. |
| Lost commit reply / partial index / different HEAD | M2 `TestOperatorCommitIdentity`, M3 `TestOperatorCommitReconcile` | G-parent/tree/message/committer/index, J-early-commit. Confirm exact descendant, a bad trailer or other message fails the same exact byte comparison. |
| Prep expiry after checks | M3 `TestOperatorFinishPrepExpiry` | J-prep-expiry. Advance the injected clock on the result. Retain the acknowledged check. A later successor must explicitly carry any accepted receipts. |
| Prompt size/entropy/store busy/V3 cap | S10 `TestOperatorFinishCardWriteRefusal` | Reserve V3 only after retained local receipt. |
| V3 store/restart windows | S10 `TestOperatorFinishCardRestart` | Same V3 id/challenge/expiry and intent read-back. |
| Missing/superseded/expired/read-failed/revoked V3 approval | S10 `TestOperatorFinishCardRefusal` | No dispatch from a nonmatching proof or request. |
| Typed/other-action approval reuse | S10 `TestOperatorFinishRejectReuse` | Runner accepts only its linked V3 approval. |
| Proof for different V3 prompt/decision | S10 `TestOperatorFinishOwnBytes` | Compare durable own bytes, decision, digest and lifetime. |
| Changed HEAD/remote/body/title/target after V3 approval | S10 `TestOperatorFinishPreDispatchFence` | Final pre-dispatch fence before invoking S11. |
| Equal remote head / ancestor | S11 `TestOperatorPushIdentity` | Exact SHA, ref and normal push policy. |
| Diverged/failed-read/uncertain push | S11 `TestOperatorPushUncertain` | Fenced reconciliation from the exact remote ref. |
| Label read/create/readback | S11 `TestOperatorLabelReconcile` | Exact repo/name and lost write policy. |
| Matching existing PR | S11 `TestOperatorReadExistingPR` | Adopt exact identity. |
| Mismatching/duplicate/closed/merged PR | S11 `TestOperatorPRIdentity` | Full exact outcome. |
| PR refusal before send | S11 `TestOperatorPRNoSend` | Explicit new approval for retry. |
| Lost PR response/timeout/empty lookup | S11 `TestOperatorPRResponseLost` | Reconcile exact receipt or remain uncertain. |
| Browser close/restart before writeback | S11 `TestOperatorFinishBrowser` | S10 tests Chrome approval/link, S11 verifies one PR and one finish record. |
| Linked card vs expired/cancelled/replaced action | S10 `TestOperatorFinishLinkedCard` | Scope to linked row and live generation while keeping standalone V3. |
| Late exact/forged finish receipt | S11 `TestOperatorFinishObservedRecord` | Scope to previously authorized remote effect. |
| Store crash/head CAS | S11 `TestOperatorFinishRecordCAS` | Recover exact stored heads. |
| Finished PR / failed writeback / detach/revoke | S11 `TestOperatorFinishCleanup` | Scoped cleanup without a second remote effect. |

### Security controls across the finish rows

| Control | Acceptance observation and mutation |
|---|---|
| Internal entry census and seal | `callers.sh` observes zero selected callers and exactly one M3 caller, `consumeFinish`. The positive census matches inherited terminal dispatch and the key confirmation handler. `TestOperatorFinishF1` exercises env, extra flag and surface markers and retains current refs/index/outcomes. `TestOperatorFinishNoProof` kills P-recheck, the substitution of the daemon proof for the runner's independent proof. |
| Terminal contention | Inherited `TestOperatorTerminalExclusion` drives the public lock, active action diagnostic and zero downstream calls. The inherited deletion control is removal of the lock check in `runEpisodeFinishWith`. M3 must separately exercise the new private lock lifecycle with a fake active runner and mutate acquisition to proceed despite contention. The attempt is bounded by an injected deadline and reports the active action id. |
| Logs and list projection | M3 `TestOperatorFinishAdmission` observes that routine request and status JSON omit the unique manifest marker and proof and that the daemon sink omits them. R-status targets removal of the private status boundary. The full route gate is outside-only. Inspect Needs-you and runner logs under the completed-launch fixture before enabling M3. |
| Local-only outcome | Each private journal fixture observes its V3 id, dispatch seam calls, scope and finish-record counts. The existing public typed fixture is the positive remote-call control. The census additionally checks every push/label create/PR seam against base. |

The §5.5 readiness row is also M1/M3. The fake runtime cases make a separate Git read fail after successful gh and observe zero new claims on both confirmed and resume paths.

### Carried notes

The selected product leaves one line under the changed-line cap. Two scratch-only test sketches measure 31 added lines for the timer/POST observations and 24 for controlled reply ordering and frozen-body checks. They are scope estimates from written sketches, with no claimed Chrome passes. The Chrome harness is outside-only and the full review/link browser work belongs to S10.

| Judge carry | Disposition |
|---|---|
| C-09 | Defer the 300 ms Reconnect callback pin with the shared 31-line observation sketch. Observe the actual POST count as well as the durable action count; retries can deduplicate. |
| C-10 | Defer the 300 ms SSH-submit callback pin with those same 31 lines. Observe actual POST and answer state before the first click. |
| C-10b | Defer the actual **400 ms** passkey callback control with the same sketch. The original edit in judge-447 `chrome2.json` uses 400 ms. The judge classified it as a survivor. |
| C-01 | Defer the `() => true` review liveness edit. The shared DOM sketch is 24 lines and requires observed supersession with deterministic delivery. |
| C-02 | Defer the removed `draft?.body ??` edit with the same DOM sketch. Compare actual restored review text after a competing bar reply. |

Exact original edits are copied to `$S/carry-size/judge-original.json` from judge-447 `chrome.json` and `chrome2.json`. The carried ledger has five Judge-447 survivors and zero new S9 browser runs. No code adding an auto-gesture was found in the S8 product.

### Existing tests and fixtures modified

| File and test/helper | Before and after | Control |
|---|---|---|
| `cmd/ailang-worldd/operator_service_test.go`, `operatorService` | The same fixture runner configuration is initialized with named fields. Retains `checkRunnerConfig`, lifecycle and health fakes. New optional finish fields default to zero for all prior tests. | `TestOperatorServiceFence` ownership/refusal and `TestOperatorFinishAuthInstall` (E2-install) verify the installed route. |
| `host/surfacepage/app_operator_test.go`, M3 | Existing harness/tests are unchanged. Adds the executable B-prep-* finish case to the same helper. | Existing start review/network/link controls still execute, and finish has no task. |

No existing golden, asset hash, signer, passkey fixture or Chrome helper is edited by the selected patches.

### Stem-matched selectors

- M1 `cmd/ailang-worldd`: `CMD_FINISH_STEMS`, `CMD_FINISH_OUTSIDE`, `CMD_SERVICE`, `CMD_SCOPED`; `TestFinishPromptIsTheGolden` reads the pinned golden, and `TestOperatorFlagParity` checks the public parser inventory. Outside-only are the six episode-finish V3 listener fixtures and `TestOperatorServiceStartup`.
- M1/M2 `host/worktree`: all tests in `worktree_test.go`; also full package plain and race. `host/checkrun`: all `checkrun_test.go` tests and `TestRunEnvProfile`. New M2 files use their exact top-level test names in `CMD_CORE`. The store step/manifest files have no existing same-stem test file. The store guards include the existing claim/schema/admission coverage.
- M3 `host/daemon`: all `surface_operator*_test.go` top-level tests plus `TestSurfaceVerbOutcomes` and its unchanged golden. `host/surfacepage`: `APP_STEMS` covers `app*_test.go` through the assembly byte/pinned-copy test, verbs, follow, Needs-you fixtures and new bar test. `CMD_STEM_ROUTE` covers affected runner integration. The named acceptance script runs Node syntax plus executable tests. Exact expanded selector strings are embedded below.

## 6. Gate and mutation receipts

| Gate | M1 local | M2 local | M3 compiled appendix |
|---|---|---|---|
| gofmt, `go build ./...`, native/Linux vet | OK | OK | Build and native/Linux vet OK |
| Store deadline/root/schema, CLI/writer/budget, boundary and daemon static inventories | OK | OK | Selected scripts must be rerun on M3 for release |
| `TestSubprocessSinkGate` | OK | OK | Pending M3 sweep |
| `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` | UNINFORMATIVE UNDER SANDBOX, `httptest` bind | Same | Outside |
| Named acceptance and stems | Passed in apply rehearsal | Passed in apply rehearsal. A separate sweep hit `TestOperatorProductionDeadlines` at the existing 5 s boundary under load. Its guard run passed. | Fast acceptance and executable Node tests passed. Whole surfacepage plain and race passed on base and M3. Real listener outside |
| 1777 fixture | `ACCEPT_OK M1 1777` | `ACCEPT_OK M2 1777` | `ACCEPT_OK M3 1777`, outside recheck pending |
| `./scripts/verify_ail.sh` | OK, temp paths relocated in scratch | OK, same script and package data | Pending M3 sweep |
| Whole `host/verifygate` | Recorded, timed fixture | Recorded, same fixture | Outside attribution required |
| Whole `host/broker` | Recorded, UNINFORMATIVE UNDER SANDBOX | Recorded, UNINFORMATIVE UNDER SANDBOX | Outside against six base failures |
| Caller census | Exact base/M1 counts printed | Exact base/M2 counts printed | Exact positive private/commit call and remote counts printed |
| Race | 9.13 s named | 12.73 s named | 6.24 s final named cmd. Earlier 98.11 s. Daemon outside |
| Chrome | No run scheduled | No run scheduled | No run scheduled for this scope |

Logs are under `$S/local-gates-M1`, `$S/local-gates-M2`, `$S/activate-local`, and `$S/census-M3.raw`. The controller must run all listener tests, V3 route neighbours, the launchd-shaped fake cases, the real Git temp-fixture outside recheck, and whole package comparisons before accepting the selected snapshots. Chrome has zero time charged to these budgets.

A scratch-only real-mux recorder diagnostic, without a listener, exposed a wrong stale-catalogue class in the new test and a count assertion after its seeded commit fixture. Both were corrected and the recorder diagnostic passed. This is **UNINFORMATIVE UNDER SANDBOX** for the route gate. Its output is in `$S/route-diagnostic-r3.raw`.


The local whole verifygate comparison against the controller's clean outside base found the same `TestMissionConfigDefaultFlowIntegration` timeout at M1 and M2. An isolated sandbox control reproduces that exact timeout on the pristine export and on both stages (`$S/repeat-controls/{base-tree,st-M1,st-M2}/verifygate.raw`). The outside whole-package comparison remains an executor prerequisite. The earlier M2 race sweep also read a changed script while executing. The sealed `$S/delivery-race` and `$S/delivery-1777` receipts distinguish the final tests from that instrument failure.

The diff ledger uses concrete edits anchored in the compiled stage. Each row in the embedded mutation catalog records the source, line, exact edit and named test. Whole-package proof uses the prerecorded pristine fail sets. A noncompiling trial is excluded from the ledger. The committer edit was corrected to preserve compilation and then killed by the named test.

| Stage | Named KILLED | Named survivors or outside |
|---|---|---|
| M1 | 21. `F1-surface`, `F1-env`, `F1-flag`, `TTY-prompt`, `E2-tool-pin`, `E2-login`, `E2-label-read`, `E2-permission`, `E2-repo-read`, `E2-remote`, `E2-https`, `E2-keychain`, `E2-ssh`, `E2-gh-env`, `E2-install`, `E2-optional`, `P-artifact`, `P-recheck`, `P-expiry`, `P-owner`, `P-generation` | Zero named survivors |
| M2 | 12. `G-hooks`, `G-parent`, `G-tree`, `G-message`, `G-committer`, `G-index`, `G-hookconfig`, `G-filter`, `G-helper`, `G-monitor`, `G-fifo`, `G-symlink` | Zero named survivors |
| M3 | 12. `J-manifest`, `J-early-commit`, `J-unsandboxed`, `J-unknown`, `J-fail-check`, `J-bytes`, `J-prep-expiry`, `J-no-ref`, `R-auth`, `B-prep-task`, `B-prep-kind`, `B-prep-flag` | `R-prep-card` UNINFORMATIVE UNDER SANDBOX, `R-status` UNINFORMATIVE UNDER SANDBOX |


The receipts are `$S/mutations-selected-m1.raw`, `$S/mutations-selected-next.raw` `$S/mutations-deferred.raw` and `$S/mutation-prep-expiry.raw`. S9 selected mutations have zero observed named survivors. These runs use the final production hunks with the mutation-identified fixture cases. The last fixture-only optimizations keep the same named assertions.

The **carried** ledger is separate: C-01, C-02, C-09, C-10 and C-10b all survived Judge-447. All five are deferred with exact edits and scratch line counts in §5. No carried edit ran in this plan.

Whole-package `B-prep-task` is KILLED against the matching clean pristine surfacepage fail set. The raw parser control and named behavioural failure are in `$S/delivery-surfacepage/whole-B-prep-task`. PENDING — controller to measure the other 23 prioritized whole-package ids with `whole-mutants.sh`, plus the outside named route arm. The remaining untested independent eligibility, source/preset/profile, path replacements and worktree-config race comparisons are explicit S9b prerequisites. Additional changed-arm coverage is unmeasured where no name appears in the catalog.

### Reproduction scripts

````bash accept.sh
#!/bin/bash
# TREE M1|M2|M3 [fast|race|outside|outside-race|1777|outside-1777]
set -euo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
M=$2 mode=${3:-fast}
. "$S/selectors.sh"
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export GOPROXY=off GIT_OPTIONAL_LOCKS=0
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
L=$(mktemp -d "$S/accept.$M.$mode.XXXXXX")
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
cd "$T"
r=(-p 1)
case "$mode" in race|outside-race) r=(-race -p 1);; fast|outside|1777|outside-1777) ;; *) exit 2;; esac
case "$M" in M1) new=$CMD_S9;; M2) new="$CMD_S9|$CMD_CORE";; M3) new="$CMD_S9|$CMD_CORE|$CMD_ACTIVATE";; *) exit 2;; esac
named() {
  local tag=$1 pkg=$2 test=$3 rc=0 limit=240
  if [ "$M" = M3 ] && [ "$tag" = cmd-finish-stems ]; then limit=480; fi
  shift 3
  python3 "$S/bounded.py" "$((limit+23))" go test -v -count=1 -timeout "${limit}s" "${r[@]}" "$@" "./$pkg/" -run "^($test)\$" > "$L/$tag.raw" 2>&1 || rc=$?
  cat "$L/$tag.raw"
  if grep -Eq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$L/$tag.raw"; then
    echo "INSTRUMENT FAILURE $tag truncated"; return 2
  fi
  if grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$L/$tag.raw"; then
    echo "UNINFORMATIVE UNDER SANDBOX $tag"; return 77
  fi
  python3 "$S/failset.py" "$L/$tag.raw" "$rc" > "$L/$tag.fails"
  grep -Eq "^--- PASS: ($test) " "$L/$tag.raw"
  [ "$rc" -eq 0 ]
}
if [[ "$mode" == *1777 ]]; then
  walk='TestOperatorFinish(Auth|AuthInstall)|TestRunnerPathRefusalClasses'
  if [ "$M" = M3 ]; then walk="$walk|TestOperatorFinishReadiness"; fi
  named cmd-1777 cmd/ailang-worldd "$walk"
else
  named cmd-new cmd/ailang-worldd "$new" -skip "$CMD_SKIP"
  while IFS= read -r expected; do
    grep -Fq -- "--- PASS: $expected (" "$L/cmd-new.raw" || { echo "MISSING PASS $expected"; exit 1; }
  done < <(printf '%s\n' "$new" | tr '|' '\n')
  if [[ "$mode" != *race ]]; then
    named worktree-stems host/worktree 'TestAddAndRemoveAWorktree|TestAddWithBranch|TestAddIgnoresAmbientGitAndRegistryEnvironment|TestIsLinked|TestOutputReturnsStdoutWithoutStderr|TestOutputRefusesStdoutPastItsCap|TestAddFailureIsReported|TestRemoveRegisteredLendsDescriptor|TestRemoveRegisteredRefusesLockedOrDirty'
    named checkrun-stems host/checkrun 'TestCheckrunHeadsAndSessionScrub|TestCheckrunWallCapKillsProcessGroup|TestCheckrunDeclaresInvalidInputs|TestRunEnvProfile'
    named cmd-finish-stems cmd/ailang-worldd "$CMD_FINISH_STEMS" -skip "$CMD_FINISH_OUTSIDE"
    named cmd-service-stems cmd/ailang-worldd "$CMD_SERVICE|$CMD_SCOPED" -skip TestOperatorServiceStartup
    if [ "$M" = M3 ]; then
      "$WORLD_EXEC_NODE" --check host/surfacepage/app/app.js
      named app-stems host/surfacepage "$APP_STEMS"
    fi
  fi
fi
case "$mode" in
  outside)
    named cmd-listener cmd/ailang-worldd "$CMD_FINISH_OUTSIDE|TestOperatorServiceStartup" -skip "$CMD_SKIP"
    if [ "$M" = M3 ]; then
      named daemon-listener host/daemon "$DAEMON_STEMS"
      named cmd-route cmd/ailang-worldd "$CMD_STEM_ROUTE" -skip "$CMD_SKIP"
    fi ;;
  outside-race)
    if [ "$M" = M3 ]; then
      named daemon-new host/daemon TestOperatorFinishAdmission
    fi ;;
  fast|race|1777) echo "LEG listener-and-launchd SKIPPED (outside.sh TREE $M)" ;;
esac
if [[ "$mode" == *race ]]; then
  added=$new
  if [ "$M" = M3 ]; then added=$CMD_ACTIVATE; fi
  race_limit() {
    awk -v selector="$2" -v cap="$3" -v pkg="$4" '
      $1=="---" && $2=="PASS:" && $3~("^(" selector ")$") {
        t=$4; gsub(/[()s]/,"",t); s+=t; n++
      }
      END {
        printf "RACE_BUDGET %s n=%d added=%.2fs cap=%ds\n",pkg,n,s,cap
        exit (n==0 || s>cap)
      }' "$1"
  }
  race_limit "$L/cmd-new.raw" "$added" 15 cmd/ailang-worldd
  if [ "$mode" = outside-race ] && [ "$M" = M3 ]; then
    race_limit "$L/daemon-new.raw" TestOperatorFinishAdmission 12 host/daemon
  fi
fi
echo "LEG chrome SKIPPED (S9 route, runner and executable Node acceptance)"
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
BASE="$HOME/.ailang/state/world-iter448/base"
B=9a40fe68e45cd9b5c787aecd716b047327bddb22
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
  local tag=$1 pkg=$2 rc=0 parser=0 start cap=945 limit=900s
  shift 2
  case "$tag" in cmd-race) set -- "$@" -skip "$CMD_SKIP";; cmd-plain|daemon-plain) cap=2800; limit=45m;; esac
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
pkgs=(cmd/ailang-worldd host/checkrun host/worktree)
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
  run 1777 env TMPDIR="$L/t1777/x" bash "$S/accept.sh" "$T" "$M" 1777
  run verify_ail bound 850 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
  for i in "${!tags[@]}"; do
    t=${tags[$i]}
    if [ "$mode" = local ] && [[ "$t" == cmd-* || "$t" == daemon-* ]]; then
      skip "$t" 'outside sandbox required'
      continue
    fi
    if [[ "$t" == *-race ]]; then record "$t" "${owners[$i]}" -race; else record "$t" "${owners[$i]}"; fi
  done
  # All recordings precede attribution. Compare comparable fail sets only.
  for tag in "${tags[@]}"; do
    if [ ! -f "$L/$tag.verified" ]; then skip "compare-$tag" 'stage record unavailable'; continue; fi
    source="$BASE/$tag"
    case "$tag" in cmd-race|worktree-*|checkrun-*|surfacepage-*) source="$S/base-extra/$tag";; esac
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
B=9a40fe68e45cd9b5c787aecd716b047327bddb22
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter448/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s9-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s9.json"
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
  count local-commit '\.git\.commit\(' "$4"
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
# LOG_PARENT [M1|M2|M3]. Stage filter defaults to all available IDs.
set -euo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
L=$1
mkdir "$L"
python3 - "$S" "$L" "${2:-}" <<'PY'
import json,sys
from pathlib import Path
s,l,stage=Path(sys.argv[1]),Path(sys.argv[2]),sys.argv[3]
rows=[r for r in json.loads((s/'mutants.json').read_text()) if not stage or r['stage']==stage]
all_ids=set((s/'whole-mutant-ids.txt').read_text().split())
(l/'whole.ids').write_text(''.join(r['id']+'\n' for r in rows if r['id'] in all_ids))
(l/'outside-named.ids').write_text(''.join(r['id']+'\n' for r in rows if r['id'] not in all_ids and r['outside']))
PY
rc=0
python3 "$S/mutate.py" "$L/prioritized" "$L/whole.ids" --outside --whole | tee "$L/whole.raw" || rc=1
if [ -s "$L/outside-named.ids" ]; then
  python3 "$S/mutate.py" "$L/named" "$L/outside-named.ids" --outside | tee "$L/named.raw" || rc=1
fi
echo ALL_DONE
exit "$rc"
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
python3 "$S/bounded.py" 1250 go test -count=1 ./host/verifygate/ ./host/broker/ > "$L/raw" 2>&1 || rc=$?
if grep -Eq 'panic: test timed out|BOUND EXCEEDED' "$L/raw"; then
  echo 'INSTRUMENT FAILURE: truncated commit guard'; exit 2
fi
python3 "$S/failset.py" "$L/raw" "$rc" > "$L/head" || exit 2
B="$HOME/.ailang/state/world-iter448/base"
test "$(tail -1 "$B/../baseline.done")" = 'ALL_DONE for 9a40fe68e45cd9b5c787aecd716b047327bddb22' || exit 2
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
"""Run one test arm per named edit in a disposable copy. Outside mode adds package legs."""
import argparse
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys

s = Path(__file__).resolve().parent
cli = argparse.ArgumentParser()
cli.add_argument("out", type=Path)
cli.add_argument("ids", type=Path)
cli.add_argument("--outside", action="store_true")
cli.add_argument("--whole", action="store_true")
args = cli.parse_args()
assert args.out.is_absolute()
args.out.mkdir()
specs = {r["id"]: r for r in json.loads((s/"mutants.json").read_text())}
trees = {}
env = dict(os.environ, AILANG_BIN=os.environ["HOME"]+"/.pinned-ailang/ailang", WORLD_REQUIRE_SSH_KEYGEN="1",
           WORLD_EXEC_NODE="/opt/homebrew/bin/node", WORLD_TOOL_AILANG_BIN=os.environ["HOME"]+"/.pinned-ailang-tools/v0.54.15/ailang",
           GOPROXY="off", GIT_OPTIONAL_LOCKS="0")
for k in ("GOFLAGS","WORLD_CHROME","WORLD_REQUIRE_PASSKEY_BROWSER","WORLD_PASSKEY_RECORD_DIR"):
    env.pop(k, None)
pkg_owner = {"app.js":"host/surfacepage"}

def gate(id, kind, pkg, test, tree, extra):
    logdir = args.out/id/kind
    logdir.mkdir(parents=True)
    for d in ("state", "cache"):
        (logdir/d).mkdir()
    child = dict(env, AILANG_STATE_DIR=str(logdir/"state"), AILANG_CACHE_DIR=str(logdir/"cache"))
    raw = logdir/"raw"
    cap = (2800 if pkg in ("cmd/ailang-worldd", "host/daemon") else 944) if kind == "whole" else 235
    with raw.open("wb") as f:
        rc = subprocess.run([sys.executable, str(s/"bounded.py"), str(cap), "go","test","-count=1","-v","-timeout",
                             str(cap-20)+"s", "./"+pkg+"/", *extra], cwd=tree, env=child, stdout=f, stderr=subprocess.STDOUT).returncode
    data = raw.read_text(errors="replace")
    if re.search(r"(bind|listen tcp).*(: operation not permitted|: permission denied)",data,re.I):
        return "UNINFORMATIVE UNDER SANDBOX"
    if re.search(r"panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]",data):
        return "INSTRUMENT FAILURE"
    with (logdir/"fails").open("w") as failed, (logdir/"parser").open("w") as parser:
        ok = subprocess.run([sys.executable,str(s/"failset.py"),str(raw),str(rc)], stdout=failed,stderr=parser).returncode
    if ok:
        return "INSTRUMENT FAILURE"
    if kind == "whole":
        stem = {"cmd/ailang-worldd":"cmd", "host/daemon":"daemon", "host/surfacepage":"surfacepage"}.get(pkg, pkg.removeprefix("host/"))
        tag = stem+"-plain"
        root = s.parent/"base" if stem in ("cmd","daemon","store","broker","verifygate") else s/"base-extra"
        baseline = root/(tag+".fails" if root.name=="base" else tag+".verified")
        if not baseline.is_file():
            return "UNMEASURED BASE FAIL SET"
        if re.search(r"panic: test timed out|BOUND EXCEEDED", (root/(tag+".raw")).read_text(errors="replace")):
            return "INSTRUMENT FAILURE BASE TRUNCATED"
        base = set(baseline.read_text().splitlines())
        current = set((logdir/"fails").read_text().splitlines())
        new = current-base
        (logdir/"new").write_text(''.join(line+'\n' for line in sorted(new)))
        if pkg+" "+test not in new:
            return "SURVIVED BASE SUBTRACTION"
    if re.search(r"^--- FAIL: "+re.escape(test)+r" \(",data,re.M):
        return "KILLED"
    if re.search(r"^--- PASS: "+re.escape(test)+r" \(",data,re.M):
        return "SURVIVED"
    return "INSTRUMENT FAILURE"

failed = 0
for id in args.ids.read_text().split():
    r = specs[id]
    if r["outside"] and not args.outside:
        print(id, "named", "UNINFORMATIVE UNDER SANDBOX", flush=True)
        continue
    stage = r["stage"]
    if stage not in trees:
        src, dst = s/("st-"+stage), args.out/("copy-"+stage)
        if sys.platform == "darwin":
            subprocess.run(["/bin/cp","-cR",str(src),str(dst)],check=True)
        else:
            shutil.copytree(src,dst,symlinks=True)
        trees[stage] = dst
    tree = trees[stage]
    edit = tree/r["file"]
    old = edit.read_text()
    if old.count(r["before"]) != 1 or old[:old.index(r["before"])].count("\n")+1 != r["line"]:
        raise SystemExit(f"mutation anchor changed: {id}")
    pkg = pkg_owner.get(edit.name, str(Path(r["file"]).parent))
    edit.write_text(old.replace(r["before"],r["after"],1))
    try:
        named = gate(id, "named", pkg, r["test"], tree, ["-run","^"+r["test"]+"$"])
        print(id, "named", named, flush=True)
        failed += int(named != "KILLED")
        if args.whole and r["whole"]:
            result = gate(id, "whole", pkg, r["test"], tree, [])
            print(id, "whole", result, flush=True)
            failed += int(result != "KILLED")
    finally:
        edit.write_text(old)
print("ALL_DONE",flush=True)
sys.exit(1 if failed else 0)
````

````python mutant-spec.py
from pathlib import Path
import json

S = Path(__file__).resolve().parent
rows = []
def spec(id, m, file, test, old, new, outside=False, whole=False):
    f = S / ('st-' + m) / file
    source = f.read_text()
    if source.count(old) != 1:
        raise SystemExit(f'{id}: {f.name} anchor count {source.count(old)}')
    rows.append(dict(id=id, stage=m, file=file, test=test, before=old, after=new, outside=outside, whole=whole,
                     line=source[:source.index(old)].count('\n') + 1))

cmd = 'cmd/ailang-worldd/'
auth = cmd + 'operator_finish_auth.go'
prep = cmd + 'operator_finish_preflight.go'
ep = cmd + 'episode_finish.go'
svc = cmd + 'operator_service.go'
test = 'TestOperatorFinishAuth'
f1 = 'tty, err = env.openTerminal()'
spec('F1-surface', 'M1', ep, 'TestOperatorFinishF1', f1,
     f1 + '\n\t\tif err != nil && strings.Contains(strings.Join(args, " "), "--confirm-on-surface") { tty, err = os.Open(os.DevNull) }', whole=True)
spec('F1-env', 'M1', ep, 'TestOperatorFinishF1', f1,
     f1 + '\n\t\tif err != nil && os.Getenv("WORLD_OPERATOR_FINISH_PREP") != "" { tty, err = os.Open(os.DevNull) }', whole=True)
spec('F1-flag', 'M1', ep, 'TestOperatorFinishF1', f1,
     f1 + '\n\t\tif err != nil && strings.Contains(strings.Join(args, " "), "--operator-prep") { return exitOK }')
spec('TTY-prompt', 'M1', ep, 'TestOperatorFinishTTYParity',
     'finishPrompt(tty, scope, actualTitle, rendered, protected, pinned.Committer)', 'fmt.Fprint(tty, "")', whole=True)
spec('E2-tool-pin', 'M1', auth, test,
     'c.Executables[real] != "" && serviceFileHash(real) == c.Executables[real]',
     '(c.Executables[real] == "" || serviceFileHash(real) == c.Executables[real])', whole=True)
spec('E2-login', 'M1', auth, test, 'return operatorGitHubAuth(3)', '/* ignored failed login */')
spec('E2-label-read', 'M1', auth, test, 'return operatorGitHubAuth(4)', '/* ignored failed label read */')
spec('E2-permission', 'M1', auth, test, 'permission != "ADMIN")', 'permission != "ADMIN" && permission != "READ")')
spec('E2-repo-read', 'M1', auth, test, 'if !ok || (permission != "WRITE"', 'if false && !ok || (permission != "WRITE"')
spec('E2-remote', 'M1', auth, test,
     'fields[4] != "HEAD" {\n\t\treturn operatorGitHubAuth(6)', 'fields[4] != "HEAD" {\n\t\treturn nil', whole=True)
spec('E2-https', 'M1', auth, test, '"credential.helper="+a.Helper', '"credential.helper="')
spec('E2-keychain', 'M1', auth, test, 'helper = "!" + shellQuote(path)', 'helper = "!" + shellQuote(gh)')
spec('E2-ssh', 'M1', auth, test, ' -o BatchMode=yes -o StrictHostKeyChecking=yes', ' -o StrictHostKeyChecking=yes')
spec('E2-gh-env', 'M1', auth, test, '"GH_PROMPT_DISABLED=1", "GIT_TERMINAL_PROMPT=0"', '"GH_PROMPT_DISABLED=1", "GIT_TERMINAL_PROMPT=0", "GH_TOKEN="+os.Getenv("GH_TOKEN")')
spec('E2-install', 'M1', svc, 'TestOperatorFinishAuthInstall', 'for _, remote := range c.FinishRemotes {\n\t\t\t\tif err := prepareOperatorFinish', 'for _, remote := range c.FinishRemotes[:0] {\n\t\t\t\tif err := prepareOperatorFinish', whole=True)
spec('E2-optional', 'M1', svc, 'TestOperatorFinishAuthInstall', 'continue // The per-verb readiness gate pins finish tools.', '/* finish-only tools become fatal */')
proof = 'TestOperatorFinishNoProof'
spec('P-artifact', 'M1', prep, proof,
     'if err != nil {\n\t\treturn operatorFinishPrep{}, approval.ErrNotVerified\n\t}\n\treturn operatorFinishPrep{frozen, row, proof}',
     'if false && err != nil {\n\t\treturn operatorFinishPrep{}, approval.ErrNotVerified\n\t}\n\treturn operatorFinishPrep{frozen, row, proof}', whole=True)
spec('P-recheck', 'M1', prep, proof, ' && p.authority.Rechecked()', '', whole=True)
spec('P-expiry', 'M1', prep, proof, 'now >= row.ExpiresAt {', 'now > row.ExpiresAt {')
spec('P-owner', 'M1', prep, proof, ' || row.Owner != owner', '', whole=True)
spec('P-generation', 'M1', prep, proof, ' || row.Generation != cas.Generation', '')

tree = cmd + 'operator_finish_tree.go'
t = 'TestOperatorCommitIdentity'
spec('G-hooks', 'M2', tree, t, '"-c", "core.hooksPath=/dev/null", ', '', whole=True)
spec('G-parent', 'M2', tree, t, 'lines[1] != "parent "+spec.Snapshot.Head || ', '', whole=True)
spec('G-tree', 'M2', tree, t, 'lines[0] != "tree "+tree ||', 'false && lines[0] != "tree "+tree ||', whole=True)
spec('G-message', 'M2', tree, t, ' || message != spec.Message', ' || false && message != spec.Message', whole=True)
spec('G-committer', 'M2', tree, t, ' && committer+">" == "committer "+spec.Committer',
     ' && (committer+">" == "committer "+spec.Committer || committer != "")', whole=True)
spec('G-index', 'M2', tree, t, 'if err != nil || index != "" {', 'if err != nil || false && index != "" {', whole=True)
for token, name in [('core.hooksPath','hookconfig'),('filter.fixture.clean','filter'),('credential.helper','helper'),('core.fsmonitor','monitor')]:
    spec('G-'+name, 'M2', tree, 'TestOperatorFinishWorktreeFence', 'case name == "remote."+start.Remote.Name+".fetch":',
         'case name == "remote."+start.Remote.Name+".fetch" || strings.EqualFold(name, "' + token + '"):')
spec('G-fifo', 'M2', tree, 'TestOperatorFinishWorktreeFence',
     'path != ".git" && !d.IsDir() && !trackedPaths[path]', 'path != ".git" && !d.IsDir() && d.Type()&os.ModeNamedPipe == 0 && !trackedPaths[path]', whole=True)

spec('G-symlink', 'M2', tree, 'TestOperatorFinishWorktreeFence', 'unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_NONBLOCK | unix.O_CLOEXEC',
     'unix.O_RDONLY | unix.O_NONBLOCK | unix.O_CLOEXEC')
local = cmd + 'operator_finish_local.go'
runner = cmd + 'operator_finish_runner.go'
spec('J-manifest', 'M3', local, 'TestOperatorFinishPrep', 'if !w.boundManifest(w.generation) {', 'if false && !w.boundManifest(w.generation) {', whole=True)
spec('J-early-commit', 'M3', local, 'TestOperatorCommitReconcile',
     'key := w.step("commit", tree)', 'key := w.step("commit", tree)\n\t_ = w.commit(ctx, m.Commit, contents, tree)', whole=True)
spec('J-unsandboxed', 'M3', local, 'TestOperatorFinishUnsandboxed', 'if check.Unsandboxed {', 'if false && check.Unsandboxed {', whole=True)
spec('J-unknown', 'M3', local, 'TestOperatorFinishChecks', 'if j.OpenCheck || len(j.Checks) > len(m.Checks)', 'if false && j.OpenCheck || len(j.Checks) > len(m.Checks)')
spec('J-fail-check', 'M3', local, 'TestOperatorFinishChecks', 'passed := result.ExitCode == 0', 'passed := result.ExitCode >= 0')
spec('J-bytes', 'M3', local, 'TestOperatorFinishStale', 'if !reflect.DeepEqual(snap, w.manifest.Commit.Snapshot) {', 'if false && !reflect.DeepEqual(snap, w.manifest.Commit.Snapshot) {', whole=True)
spec('J-prep-expiry', 'M3', local, 'TestOperatorFinishPrepExpiry',
     'if err := w.check(ctx, check, i+len(j.Checks)); err != nil {\n\t\t\treturn err\n\t\t}',
     'if err := w.check(ctx, check, i+len(j.Checks)); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tw.fresh = func(context.Context) error { return nil }\n\t\tw.now = func() int64 { return w.prep.row.CreatedAt }',
     whole=True)
spec('J-no-ref', 'M3', local, 'TestOperatorCommitFailed',
     'if errors.Is(commitErr, errOperatorCommitNoRef) {', 'if false && errors.Is(commitErr, errOperatorCommitNoRef) {')
spec('R-auth', 'M3', runner, 'TestOperatorFinishReadiness', 'if err := prepareOperatorFinish(ctx, c, m.Remote.URL); err != nil {',
     'if err := prepareOperatorFinish(ctx, c, m.Remote.URL); false && err != nil {', whole=True)
spec('R-prep-card', 'M3', 'host/daemon/surface_operator_confirm.go', 'TestOperatorFinishAdmission',
     ' || a.Kind == "finish" && a.Options.ConfirmOnSurface', '', outside=True, whole=True)
spec('R-status', 'M3', 'host/daemon/surface_operator_read.go', 'TestOperatorFinishAdmission',
     'status.LocalCommit = j.Commit.SHA', 'status.LocalCommit = ""', outside=True)
spec('B-prep-task', 'M3', 'host/surfacepage/app/app.js', 'TestOperatorFinishBarReview',
     "argv[2] === 'finish' ? task === ''", "argv[2] === 'finish' ? true", whole=True)
spec('B-prep-kind', 'M3', 'host/surfacepage/app/app.js', 'TestOperatorFinishBarReview',
     "b.kind !== kind", "false && b.kind !== kind")
spec('B-prep-flag', 'M3', 'host/surfacepage/app/app.js', 'TestOperatorFinishBarReview',
     "argv.length === 8 && argv[7] === '--ready'", "argv.length >= 8 && argv[7] === '--ready'")

(S/'mutants.json').write_text(json.dumps(rows,indent=2)+'\n')
(S/'whole-mutant-ids.txt').write_text(''.join(r['id']+'\n' for r in rows if r['whole']))
(S/'mutant-local-ids.txt').write_text(''.join(r['id']+'\n' for r in rows if not r['outside']))
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
    "id": "F1-surface",
    "stage": "M1",
    "file": "cmd/ailang-worldd/episode_finish.go",
    "test": "TestOperatorFinishF1",
    "before": "tty, err = env.openTerminal()",
    "after": "tty, err = env.openTerminal()\n\t\tif err != nil && strings.Contains(strings.Join(args, \" \"), \"--confirm-on-surface\") { tty, err = os.Open(os.DevNull) }",
    "outside": false,
    "whole": true,
    "line": 89
  },
  {
    "id": "F1-env",
    "stage": "M1",
    "file": "cmd/ailang-worldd/episode_finish.go",
    "test": "TestOperatorFinishF1",
    "before": "tty, err = env.openTerminal()",
    "after": "tty, err = env.openTerminal()\n\t\tif err != nil && os.Getenv(\"WORLD_OPERATOR_FINISH_PREP\") != \"\" { tty, err = os.Open(os.DevNull) }",
    "outside": false,
    "whole": true,
    "line": 89
  },
  {
    "id": "F1-flag",
    "stage": "M1",
    "file": "cmd/ailang-worldd/episode_finish.go",
    "test": "TestOperatorFinishF1",
    "before": "tty, err = env.openTerminal()",
    "after": "tty, err = env.openTerminal()\n\t\tif err != nil && strings.Contains(strings.Join(args, \" \"), \"--operator-prep\") { return exitOK }",
    "outside": false,
    "whole": false,
    "line": 89
  },
  {
    "id": "TTY-prompt",
    "stage": "M1",
    "file": "cmd/ailang-worldd/episode_finish.go",
    "test": "TestOperatorFinishTTYParity",
    "before": "finishPrompt(tty, scope, actualTitle, rendered, protected, pinned.Committer)",
    "after": "fmt.Fprint(tty, \"\")",
    "outside": false,
    "whole": true,
    "line": 454
  },
  {
    "id": "E2-tool-pin",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": "c.Executables[real] != \"\" && serviceFileHash(real) == c.Executables[real]",
    "after": "(c.Executables[real] == \"\" || serviceFileHash(real) == c.Executables[real])",
    "outside": false,
    "whole": true,
    "line": 42
  },
  {
    "id": "E2-login",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": "return operatorGitHubAuth(3)",
    "after": "/* ignored failed login */",
    "outside": false,
    "whole": false,
    "line": 88
  },
  {
    "id": "E2-label-read",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": "return operatorGitHubAuth(4)",
    "after": "/* ignored failed label read */",
    "outside": false,
    "whole": false,
    "line": 91
  },
  {
    "id": "E2-permission",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": "permission != \"ADMIN\")",
    "after": "permission != \"ADMIN\" && permission != \"READ\")",
    "outside": false,
    "whole": false,
    "line": 94
  },
  {
    "id": "E2-repo-read",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": "if !ok || (permission != \"WRITE\"",
    "after": "if false && !ok || (permission != \"WRITE\"",
    "outside": false,
    "whole": false,
    "line": 94
  },
  {
    "id": "E2-remote",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": "fields[4] != \"HEAD\" {\n\t\treturn operatorGitHubAuth(6)",
    "after": "fields[4] != \"HEAD\" {\n\t\treturn nil",
    "outside": false,
    "whole": true,
    "line": 116
  },
  {
    "id": "E2-https",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": "\"credential.helper=\"+a.Helper",
    "after": "\"credential.helper=\"",
    "outside": false,
    "whole": false,
    "line": 101
  },
  {
    "id": "E2-keychain",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": "helper = \"!\" + shellQuote(path)",
    "after": "helper = \"!\" + shellQuote(gh)",
    "outside": false,
    "whole": false,
    "line": 69
  },
  {
    "id": "E2-ssh",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": " -o BatchMode=yes -o StrictHostKeyChecking=yes",
    "after": " -o StrictHostKeyChecking=yes",
    "outside": false,
    "whole": false,
    "line": 111
  },
  {
    "id": "E2-gh-env",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_auth.go",
    "test": "TestOperatorFinishAuth",
    "before": "\"GH_PROMPT_DISABLED=1\", \"GIT_TERMINAL_PROMPT=0\"",
    "after": "\"GH_PROMPT_DISABLED=1\", \"GIT_TERMINAL_PROMPT=0\", \"GH_TOKEN=\"+os.Getenv(\"GH_TOKEN\")",
    "outside": false,
    "whole": false,
    "line": 82
  },
  {
    "id": "E2-install",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_service.go",
    "test": "TestOperatorFinishAuthInstall",
    "before": "for _, remote := range c.FinishRemotes {\n\t\t\t\tif err := prepareOperatorFinish",
    "after": "for _, remote := range c.FinishRemotes[:0] {\n\t\t\t\tif err := prepareOperatorFinish",
    "outside": false,
    "whole": true,
    "line": 348
  },
  {
    "id": "E2-optional",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_service.go",
    "test": "TestOperatorFinishAuthInstall",
    "before": "continue // The per-verb readiness gate pins finish tools.",
    "after": "/* finish-only tools become fatal */",
    "outside": false,
    "whole": false,
    "line": 166
  },
  {
    "id": "P-artifact",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_preflight.go",
    "test": "TestOperatorFinishNoProof",
    "before": "if err != nil {\n\t\treturn operatorFinishPrep{}, approval.ErrNotVerified\n\t}\n\treturn operatorFinishPrep{frozen, row, proof}",
    "after": "if false && err != nil {\n\t\treturn operatorFinishPrep{}, approval.ErrNotVerified\n\t}\n\treturn operatorFinishPrep{frozen, row, proof}",
    "outside": false,
    "whole": true,
    "line": 46
  },
  {
    "id": "P-recheck",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_preflight.go",
    "test": "TestOperatorFinishNoProof",
    "before": " && p.authority.Rechecked()",
    "after": "",
    "outside": false,
    "whole": true,
    "line": 53
  },
  {
    "id": "P-expiry",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_preflight.go",
    "test": "TestOperatorFinishNoProof",
    "before": "now >= row.ExpiresAt {",
    "after": "now > row.ExpiresAt {",
    "outside": false,
    "whole": false,
    "line": 38
  },
  {
    "id": "P-owner",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_preflight.go",
    "test": "TestOperatorFinishNoProof",
    "before": " || row.Owner != owner",
    "after": "",
    "outside": false,
    "whole": true,
    "line": 34
  },
  {
    "id": "P-generation",
    "stage": "M1",
    "file": "cmd/ailang-worldd/operator_finish_preflight.go",
    "test": "TestOperatorFinishNoProof",
    "before": " || row.Generation != cas.Generation",
    "after": "",
    "outside": false,
    "whole": false,
    "line": 34
  },
  {
    "id": "G-hooks",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorCommitIdentity",
    "before": "\"-c\", \"core.hooksPath=/dev/null\", ",
    "after": "",
    "outside": false,
    "whole": true,
    "line": 38
  },
  {
    "id": "G-parent",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorCommitIdentity",
    "before": "lines[1] != \"parent \"+spec.Snapshot.Head || ",
    "after": "",
    "outside": false,
    "whole": true,
    "line": 316
  },
  {
    "id": "G-tree",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorCommitIdentity",
    "before": "lines[0] != \"tree \"+tree ||",
    "after": "false && lines[0] != \"tree \"+tree ||",
    "outside": false,
    "whole": true,
    "line": 315
  },
  {
    "id": "G-message",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorCommitIdentity",
    "before": " || message != spec.Message",
    "after": " || false && message != spec.Message",
    "outside": false,
    "whole": true,
    "line": 316
  },
  {
    "id": "G-committer",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorCommitIdentity",
    "before": " && committer+\">\" == \"committer \"+spec.Committer",
    "after": " && (committer+\">\" == \"committer \"+spec.Committer || committer != \"\")",
    "outside": false,
    "whole": true,
    "line": 320
  },
  {
    "id": "G-index",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorCommitIdentity",
    "before": "if err != nil || index != \"\" {",
    "after": "if err != nil || false && index != \"\" {",
    "outside": false,
    "whole": true,
    "line": 306
  },
  {
    "id": "G-hookconfig",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorFinishWorktreeFence",
    "before": "case name == \"remote.\"+start.Remote.Name+\".fetch\":",
    "after": "case name == \"remote.\"+start.Remote.Name+\".fetch\" || strings.EqualFold(name, \"core.hooksPath\"):",
    "outside": false,
    "whole": false,
    "line": 144
  },
  {
    "id": "G-filter",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorFinishWorktreeFence",
    "before": "case name == \"remote.\"+start.Remote.Name+\".fetch\":",
    "after": "case name == \"remote.\"+start.Remote.Name+\".fetch\" || strings.EqualFold(name, \"filter.fixture.clean\"):",
    "outside": false,
    "whole": false,
    "line": 144
  },
  {
    "id": "G-helper",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorFinishWorktreeFence",
    "before": "case name == \"remote.\"+start.Remote.Name+\".fetch\":",
    "after": "case name == \"remote.\"+start.Remote.Name+\".fetch\" || strings.EqualFold(name, \"credential.helper\"):",
    "outside": false,
    "whole": false,
    "line": 144
  },
  {
    "id": "G-monitor",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorFinishWorktreeFence",
    "before": "case name == \"remote.\"+start.Remote.Name+\".fetch\":",
    "after": "case name == \"remote.\"+start.Remote.Name+\".fetch\" || strings.EqualFold(name, \"core.fsmonitor\"):",
    "outside": false,
    "whole": false,
    "line": 144
  },
  {
    "id": "G-fifo",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorFinishWorktreeFence",
    "before": "path != \".git\" && !d.IsDir() && !trackedPaths[path]",
    "after": "path != \".git\" && !d.IsDir() && d.Type()&os.ModeNamedPipe == 0 && !trackedPaths[path]",
    "outside": false,
    "whole": true,
    "line": 209
  },
  {
    "id": "G-symlink",
    "stage": "M2",
    "file": "cmd/ailang-worldd/operator_finish_tree.go",
    "test": "TestOperatorFinishWorktreeFence",
    "before": "unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_NONBLOCK | unix.O_CLOEXEC",
    "after": "unix.O_RDONLY | unix.O_NONBLOCK | unix.O_CLOEXEC",
    "outside": false,
    "whole": false,
    "line": 55
  },
  {
    "id": "J-manifest",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "test": "TestOperatorFinishPrep",
    "before": "if !w.boundManifest(w.generation) {",
    "after": "if false && !w.boundManifest(w.generation) {",
    "outside": false,
    "whole": true,
    "line": 87
  },
  {
    "id": "J-early-commit",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "test": "TestOperatorCommitReconcile",
    "before": "key := w.step(\"commit\", tree)",
    "after": "key := w.step(\"commit\", tree)\n\t_ = w.commit(ctx, m.Commit, contents, tree)",
    "outside": false,
    "whole": true,
    "line": 155
  },
  {
    "id": "J-unsandboxed",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "test": "TestOperatorFinishUnsandboxed",
    "before": "if check.Unsandboxed {",
    "after": "if false && check.Unsandboxed {",
    "outside": false,
    "whole": true,
    "line": 91
  },
  {
    "id": "J-unknown",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "test": "TestOperatorFinishChecks",
    "before": "if j.OpenCheck || len(j.Checks) > len(m.Checks)",
    "after": "if false && j.OpenCheck || len(j.Checks) > len(m.Checks)",
    "outside": false,
    "whole": false,
    "line": 112
  },
  {
    "id": "J-fail-check",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "test": "TestOperatorFinishChecks",
    "before": "passed := result.ExitCode == 0",
    "after": "passed := result.ExitCode >= 0",
    "outside": false,
    "whole": false,
    "line": 188
  },
  {
    "id": "J-bytes",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "test": "TestOperatorFinishStale",
    "before": "if !reflect.DeepEqual(snap, w.manifest.Commit.Snapshot) {",
    "after": "if false && !reflect.DeepEqual(snap, w.manifest.Commit.Snapshot) {",
    "outside": false,
    "whole": true,
    "line": 73
  },
  {
    "id": "J-prep-expiry",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "test": "TestOperatorFinishPrepExpiry",
    "before": "if err := w.check(ctx, check, i+len(j.Checks)); err != nil {\n\t\t\treturn err\n\t\t}",
    "after": "if err := w.check(ctx, check, i+len(j.Checks)); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tw.fresh = func(context.Context) error { return nil }\n\t\tw.now = func() int64 { return w.prep.row.CreatedAt }",
    "outside": false,
    "whole": true,
    "line": 130
  },
  {
    "id": "J-no-ref",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_local.go",
    "test": "TestOperatorCommitFailed",
    "before": "if errors.Is(commitErr, errOperatorCommitNoRef) {",
    "after": "if false && errors.Is(commitErr, errOperatorCommitNoRef) {",
    "outside": false,
    "whole": false,
    "line": 163
  },
  {
    "id": "R-auth",
    "stage": "M3",
    "file": "cmd/ailang-worldd/operator_finish_runner.go",
    "test": "TestOperatorFinishReadiness",
    "before": "if err := prepareOperatorFinish(ctx, c, m.Remote.URL); err != nil {",
    "after": "if err := prepareOperatorFinish(ctx, c, m.Remote.URL); false && err != nil {",
    "outside": false,
    "whole": true,
    "line": 203
  },
  {
    "id": "R-prep-card",
    "stage": "M3",
    "file": "host/daemon/surface_operator_confirm.go",
    "test": "TestOperatorFinishAdmission",
    "before": " || a.Kind == \"finish\" && a.Options.ConfirmOnSurface",
    "after": "",
    "outside": true,
    "whole": true,
    "line": 96
  },
  {
    "id": "R-status",
    "stage": "M3",
    "file": "host/daemon/surface_operator_read.go",
    "test": "TestOperatorFinishAdmission",
    "before": "status.LocalCommit = j.Commit.SHA",
    "after": "status.LocalCommit = \"\"",
    "outside": true,
    "whole": false,
    "line": 132
  },
  {
    "id": "B-prep-task",
    "stage": "M3",
    "file": "host/surfacepage/app/app.js",
    "test": "TestOperatorFinishBarReview",
    "before": "argv[2] === 'finish' ? task === ''",
    "after": "argv[2] === 'finish' ? true",
    "outside": false,
    "whole": true,
    "line": 1448
  },
  {
    "id": "B-prep-kind",
    "stage": "M3",
    "file": "host/surfacepage/app/app.js",
    "test": "TestOperatorFinishBarReview",
    "before": "b.kind !== kind",
    "after": "false && b.kind !== kind",
    "outside": false,
    "whole": false,
    "line": 1459
  },
  {
    "id": "B-prep-flag",
    "stage": "M3",
    "file": "host/surfacepage/app/app.js",
    "test": "TestOperatorFinishBarReview",
    "before": "argv.length === 8 && argv[7] === '--ready'",
    "after": "argv.length >= 8 && argv[7] === '--ready'",
    "outside": false,
    "whole": false,
    "line": 1446
  }
]
````

````text whole-mutant-ids.txt
F1-surface
F1-env
TTY-prompt
E2-tool-pin
E2-remote
E2-install
P-artifact
P-recheck
P-owner
G-hooks
G-parent
G-tree
G-message
G-committer
G-index
G-fifo
J-manifest
J-early-commit
J-unsandboxed
J-bytes
J-prep-expiry
R-auth
R-prep-card
B-prep-task
````

## 7. Apply, accept and snapshot

The controller commits the two plan artifacts separately. Run the block below from the clean sprint checkout with product files at the pinned base. The executor applies M1 then M2 using only `/usr/bin/patch -p1 -E`. Set `TOOLS=$STATE/exec-tools` to keep the planner's tools and receipts separate. Executor snapshots are cumulative; `$STATE/snap/<M>/STATUS` records added, changed and deleted paths, `DELETED` contains just deleted paths, and `MANIFEST.sha256` verifies the copied files. The two plan files are excluded.

For a successful apply the script prints `ACCEPT_OK <M> <mode>` and `APPLY_ACCEPT_SNAPSHOT_OK <M> <mode>` after the cumulative snapshot. Both must print and the script must return zero. Local fast acceptance intentionally skips listener and launchd-shaped arms. The controller runs those independently on the planner stage before starting the executor.


````bash
set -euo pipefail
B=9a40fe68e45cd9b5c787aecd716b047327bddb22
STATE="$HOME/.ailang/state/world-iter448"
S="$STATE/plan-scratch"
WT=${WT:-/Users/voightkampff/dev/sunholo-data/.wt-world-iter448-row243-v4s9}
PLAN=${PLAN:-$WT/design_docs/planned/w-world-address-bar-verbs-v4s9-sprint-plan.md}
TREE=${TREE:-$WT}
TOOLS=${TOOLS:-$S}
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
for f in bounded.py failset.py verify-scratch.py mutate.py mutant-spec.py patches.py; do
  grab '````python '"$f" "$TOOLS/$f"
done
grab '````json mutants.json' "$TOOLS/mutants.json"
grab '````text whole-mutant-ids.txt' "$TOOLS/whole-mutant-ids.txt"
if [ "$REHEARSAL" = 1 ]; then
  test ! -e "$TREE"
  mkdir -p "$TREE"
  git -C "$WT" archive "$B" | tar -x -C "$TREE"
else
  git -C "$TREE" merge-base --is-ancestor "$B" HEAD
  git -C "$TREE" diff --quiet "$B" HEAD -- cmd host docs app
fi
for M in $MILESTONES; do
  case "$M" in
    M1|M2) ;;
    M3) test "$REHEARSAL" = 1 ;;
    *) exit 2 ;;
  esac
  grab '````diff '"$M" "$TOOLS/$M.patch"
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
````


### Controller commands

Use the controller's existing baseline. If `baseline.done` lacks the exact `ALL_DONE for <base>` marker, defer the whole-package comparison until it is complete. The plain fixture failures in the cmd and daemon base are compared by name. A package timed out or a parser control that fails leaves that package unreadable. Complete the missing base cmd race and plain/race worktree/checkrun/surfacepage legs using the pristine export:

````bash
S="$HOME/.ailang/state/world-iter448/plan-scratch"
STATE="$HOME/.ailang/state/world-iter448"
bash "$S/gates.sh" "$S/base-tree" "$S/base-extra" M3 outside-base
bash "$S/outside.sh" "$S/st-M1" M1
bash "$S/outside.sh" "$S/st-M2" M2
bash "$S/gates.sh" "$S/st-M1" "$S/outside-gates-M1" M1 outside
bash "$S/gates.sh" "$S/st-M2" "$S/outside-gates-M2" M2 outside
bash "$S/whole-mutants.sh" "$S/outside-mutants-M1" M1
bash "$S/whole-mutants.sh" "$S/outside-mutants-M2" M2
bash "$STATE/ctl_roundtrip.sh" M1 M2
````

M3 continuation commands, once the missing fixture proofs and size/timing limits are verified:

````bash
S="$HOME/.ailang/state/world-iter448/plan-scratch"
bash "$S/outside.sh" "$S/st-M3" M3
bash "$S/gates.sh" "$S/st-M3" "$S/outside-gates-M3" M3 outside
bash "$S/whole-mutants.sh" "$S/outside-mutants-M3" M3
````

Commit-time commands use named tests, with the final selected milestone also running both whole-package guards in isolated fixture state. Copy the stage message file below under planner scratch if it is absent. These commands are for the controller, which owns the worktree's Git writes.

````bash
STATE="$HOME/.ailang/state/world-iter448"
S="$STATE/plan-scratch"
printf '%s\n' 'world: pin finish F1, prep proof and runner auth readiness' > "$S/M1.commitmsg"
printf '%s\n' 'world: freeze local finish inputs and exact commit evidence' > "$S/M2.commitmsg"
"$STATE/commit_ms.sh" M1 "$S/M1.commitmsg" 'export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"; S="$HOME/.ailang/state/world-iter448/plan-scratch"; bash "$S/accept.sh" "$PWD" M1 outside'
"$STATE/commit_ms.sh" M2 "$S/M2.commitmsg" 'export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"; S="$HOME/.ailang/state/world-iter448/plan-scratch"; bash "$S/accept.sh" "$PWD" M2 outside && bash "$S/commit-guards.sh" "$PWD" "$S/commit-guards-M2-$(date +%s)"'
````

### Rehearsal

Extracted the §7 block from the plan with awk, compared it byte-for-byte with its scratch source, and executed that extract against a fresh `git archive 9a40fe68` export. The run returned zero. It used the default M1 and M2 product chain followed by the M3 rehearsal-only continuation. The receipt is `$S/rt-apply/sealed-1795848173/receipt.raw`.

```text
ACCEPT_OK M1 fast
ROUNDTRIP_OK M1 cmd
ROUNDTRIP_OK M1 host
ROUNDTRIP_OK M1 docs
APPLY_ACCEPT_SNAPSHOT_OK M1 fast
ACCEPT_OK M2 fast
ROUNDTRIP_OK M2 cmd
ROUNDTRIP_OK M2 host
ROUNDTRIP_OK M2 docs
APPLY_ACCEPT_SNAPSHOT_OK M2 fast
ACCEPT_OK M3 fast
ROUNDTRIP_OK M3 cmd
ROUNDTRIP_OK M3 host
ROUNDTRIP_OK M3 docs
APPLY_ACCEPT_SNAPSHOT_OK M3 fast
FINAL_APPLY_RC 0
```

Exact `diff -r` results for cmd, host and docs are byte-identical for all three milestones. The top-level app tree is absent. Cumulative snapshots verified as M1 9/0/9, M2 15/0/15 and M3 28/0/28 for STATUS, DELETED and MANIFEST. The expanded M3 app selector, including the approval model and page fixtures, passed separately on the same complete M3 stage (`$S/delivery-surfacepage/expanded-stems/raw`).

In the earlier delivery rehearsal the M3 legacy stem command hit its 263 s bound after the new cmd assertions passed. Its separate snapshot verified and the retained export matched M3. The final instrument gives the M3 legacy stems 480 s, retains the M1/M2 bounds and completed the full run above.

## Appendix P — patches

M1 and M2 are sequential and selected. M3 is an incremental compiled appendix for S9b. It must begin on the M2 product bytes.

### M1 — Install auth readiness and private proof preflight

````diff M1
diff --git a/cmd/ailang-worldd/episode_finish.go b/cmd/ailang-worldd/episode_finish.go
index a6de86e..b1e4d11 100644
--- a/cmd/ailang-worldd/episode_finish.go
+++ b/cmd/ailang-worldd/episode_finish.go
@@ -77,6 +77,21 @@ func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEn
 	if episodeFlagSetSeen != nil {
 		episodeFlagSetSeen(fs)
 	}
+	fail := func(what, why, next string) int {
+		fmt.Fprintln(stderr, presetRefusal(what, why, next))
+		return exitUsage
+	}
+	// Standalone help only describes the flags. Every attempt to finish
+	// meets F1 before consulting flags, ambient inputs, or the worktree.
+	var tty io.ReadWriteCloser
+	if len(args) != 1 || (args[0] != "--help" && args[0] != "-help" && args[0] != "-h" && args[0] != "--h") {
+		var err error
+		tty, err = env.openTerminal()
+		if err != nil {
+			return fail("refusing: no controlling terminal", "finish is an attended act", "run it in a terminal (an agent or CI can't finish an episode)")
+		}
+		defer tty.Close()
+	}
 	pos, err := verbs.ParseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, episodeFinishHelp)
@@ -85,16 +100,6 @@ func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEn
 	if err != nil {
 		return exitUsage
 	}
-	fail := func(what, why, next string) int {
-		fmt.Fprintln(stderr, presetRefusal(what, why, next))
-		return exitUsage
-	}
-	// F1 precedes even bad preset/episode validation, and all IO beyond flags.
-	tty, err := env.openTerminal()
-	if err != nil {
-		return fail("refusing: no controlling terminal", "finish is an attended act", "run it in a terminal (an agent or CI can't finish an episode)")
-	}
-	defer tty.Close()
 	if *onSurface && *dry { // row 243 V3 §2.7 step 1
 		return fail("--confirm-on-surface cannot be combined with --dry-run", "a dry run opens no pull request", "drop one of them")
 	}
diff --git a/cmd/ailang-worldd/operator_finish_auth.go b/cmd/ailang-worldd/operator_finish_auth.go
new file mode 100644
index 0000000..f9b7cfc
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_auth.go
@@ -0,0 +1,128 @@
+package main
+
+import (
+	"context"
+	"fmt"
+	"net"
+	"os"
+	"path/filepath"
+	"regexp"
+	"strings"
+	"syscall"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/checkrun"
+	"github.com/sunholo-data/ailang-world/host/worktree"
+)
+
+const operatorFinishRemedy = "repair PATH, GitHub login or Git credential helper in attended service-runner setup and re-run readiness"
+
+var operatorFinishSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)
+
+type operatorGitHubAuth int
+
+func (e operatorGitHubAuth) Error() string {
+	return fmt.Sprintf("finish blocked: github-auth (readiness step %d); %s", int(e), operatorFinishRemedy)
+}
+
+type operatorFinishAuth struct {
+	GH, Git, Repo, Remote, Home, Socket, Path, Helper string
+}
+
+// Require the first executable on the approved PATH to match its pin.
+func operatorFinishTool(c runnerConfig, dirs []string, name string) string {
+	for _, dir := range dirs {
+		path := filepath.Join(dir, name)
+		st, err := os.Stat(path)
+		if err != nil || !st.Mode().IsRegular() || st.Mode().Perm()&0111 == 0 {
+			continue
+		}
+		real, err := filepath.EvalSymlinks(path)
+		if err == nil && runnerPath(c.Home, real, c.UID, false, 0) == nil &&
+			c.Executables[real] != "" && serviceFileHash(real) == c.Executables[real] {
+			return real
+		}
+		return ""
+	}
+	return ""
+}
+
+func planOperatorFinishAuth(c runnerConfig, remote string) (operatorFinishAuth, error) {
+	repo, err := finishRepository(remote)
+	if err != nil {
+		return operatorFinishAuth{}, operatorGitHubAuth(6)
+	}
+	dirs := append(append([]string(nil), c.ToolDirs...), filepath.SplitList(runnerSystemPath)...)
+	gh, git := operatorFinishTool(c, dirs, "gh"), operatorFinishTool(c, dirs, "git")
+	if gh == "" || git == "" {
+		return operatorFinishAuth{}, operatorGitHubAuth(2)
+	}
+	socket, helper := c.FinishSSHSocket, ""
+	if strings.HasPrefix(remote, "https://") {
+		socket = ""
+		helper = "!" + shellQuote(gh) + " auth git-credential"
+		if c.FinishHTTPSHelper == "osxkeychain" {
+			path := operatorFinishTool(c, dirs, "git-credential-osxkeychain")
+			if path == "" {
+				return operatorFinishAuth{}, operatorGitHubAuth(6)
+			}
+			helper = "!" + shellQuote(path)
+		}
+	}
+	if socket != "" {
+		st, err := os.Lstat(socket)
+		if err != nil || st.Mode()&os.ModeSocket == 0 || int(st.Sys().(*syscall.Stat_t).Uid) != c.UID {
+			return operatorFinishAuth{}, operatorGitHubAuth(6)
+		}
+	}
+	return operatorFinishAuth{gh, git, repo, remote, c.Home, socket, strings.Join(dirs, ":"), helper}, nil
+}
+
+func (a operatorFinishAuth) prepare(ctx context.Context) error {
+	env := []string{"HOME=" + a.Home, "PATH=" + a.Path, "GH_PROMPT_DISABLED=1", "GIT_TERMINAL_PROMPT=0"}
+	gh := func(args ...string) (string, bool) {
+		r, err := checkrun.RunIsolated(ctx, append([]string{a.GH}, args...), a.Home, 10*time.Second, env)
+		return strings.TrimSpace(r.Stdout), err == nil && r.ExitCode == 0
+	}
+	if _, ok := gh("auth", "status", "--hostname", "github.com"); !ok {
+		return operatorGitHubAuth(3)
+	}
+	if _, ok := gh("label", "list", "-R", a.Repo, "--search", "ailang-world", "--json", "name"); !ok {
+		return operatorGitHubAuth(4)
+	}
+	permission, ok := gh("repo", "view", a.Repo, "--json", "viewerPermission", "--jq", ".viewerPermission")
+	if !ok || (permission != "WRITE" && permission != "MAINTAIN" && permission != "ADMIN") {
+		return operatorGitHubAuth(5)
+	}
+
+	env = append(env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null")
+	args := []string{"--git-dir=/dev/null", "-c", "credential.interactive=never", "-c", "credential.helper="}
+	if strings.HasPrefix(a.Remote, "https://") {
+		args = append(args, "-c", "credential.helper="+a.Helper)
+	} else {
+		if a.Socket != "" {
+			conn, err := (&net.Dialer{Timeout: time.Second}).DialContext(ctx, "unix", a.Socket)
+			if err != nil {
+				return operatorGitHubAuth(6)
+			}
+			conn.Close()
+			env = append(env, "SSH_AUTH_SOCK="+a.Socket)
+		}
+		args = append(args, "-c", "core.sshCommand=/usr/bin/ssh -F /dev/null -o BatchMode=yes -o StrictHostKeyChecking=yes")
+	}
+	out, err := worktree.OutputIsolated(ctx, a.Git, append(args, "ls-remote", "--symref", "--", a.Remote, "HEAD"), env)
+	fields := strings.Fields(out)
+	if err != nil || len(fields) != 5 || fields[0] != "ref:" || !strings.HasPrefix(fields[1], "refs/heads/") ||
+		fields[2] != "HEAD" || !operatorFinishSHA.MatchString(fields[3]) || fields[4] != "HEAD" {
+		return operatorGitHubAuth(6)
+	}
+	return nil
+}
+
+func prepareOperatorFinish(ctx context.Context, c runnerConfig, remote string) error {
+	a, err := planOperatorFinishAuth(c, remote)
+	if err != nil {
+		return err
+	}
+	return a.prepare(ctx)
+}
diff --git a/cmd/ailang-worldd/operator_finish_auth_test.go b/cmd/ailang-worldd/operator_finish_auth_test.go
new file mode 100644
index 0000000..a5d0c82
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_auth_test.go
@@ -0,0 +1,158 @@
+package main
+
+import (
+	"errors"
+	"fmt"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/worktree"
+)
+
+const fixtureFinishRemote = "https://github.com/fixture/approved.git"
+
+func writeFinishTestTool(t *testing.T, c *runnerConfig, dir, name, text string) {
+	t.Helper()
+	path := filepath.Join(dir, name)
+	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+text), 0700); err != nil {
+		t.Fatal(err)
+	}
+	c.Executables[path] = serviceFileHash(path)
+}
+
+func TestOperatorFinishAuth(t *testing.T) {
+	dir, err := filepath.EvalSymlinks(t.TempDir())
+	if err != nil {
+		t.Fatal(err)
+	}
+	git, err := exec.LookPath("git")
+	if err != nil {
+		t.Fatal(err)
+	}
+	bare := filepath.Join(dir, "bare")
+	ctx := boundedTestContext(t)
+	gitRun := func(args ...string) {
+		t.Helper()
+		if _, err := worktree.OutputIsolated(ctx, git, args, []string{"HOME=" + dir, "PATH=" + runnerSystemPath,
+			"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}); err != nil {
+			t.Fatal(err)
+		}
+	}
+	gitRun("-c", "init.defaultBranch=main", "init", "--bare", "-q", bare)
+	gitRun("--git-dir="+bare, "--work-tree="+dir, "-c", "core.bare=false", "-c", "user.name=Fixture", "-c", "user.email=test@invalid",
+		"commit", "--allow-empty", "-qm", "fixture")
+	t.Setenv("GH_TOKEN", "private-fixture-token")
+	t.Setenv("GIT_CONFIG_COUNT", "forbidden")
+	t.Setenv("AILANG_REGISTRY_API_KEY", "registry-fixture-token")
+	for _, tc := range []struct {
+		name, remote, ghFailure, permission, gitResult string
+		step                                           operatorGitHubAuth
+	}{
+		{"missing-unpinned-tool", fixtureFinishRemote, "", "WRITE", "", 2},
+		{"login", fixtureFinishRemote, "auth", "WRITE", "", 3},
+		{"label-read", fixtureFinishRemote, "label", "WRITE", "", 4},
+		{"repo-read", fixtureFinishRemote, "repo", "WRITE", "", 5},
+		{"permission", fixtureFinishRemote, "", "READ", "", 5},
+		{"https-helper", fixtureFinishRemote, "", "WRITE", "exit 1", 6},
+		{"ssh-read", "git@github.com:fixture/approved.git", "", "WRITE", "exit 1", 6},
+		{"empty", fixtureFinishRemote, "", "WRITE", "exit 0", 6},
+		{"write", fixtureFinishRemote, "", "WRITE", "", 0},
+		{"maintain", fixtureFinishRemote, "", "MAINTAIN", "", 0},
+		{"keychain", fixtureFinishRemote, "", "WRITE", "", 0},
+		{"admin", "ssh://git@github.com/fixture/approved.git", "", "ADMIN", "", 0},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			tools := filepath.Join(dir, tc.name)
+			if err := os.Mkdir(tools, 0700); err != nil {
+				t.Fatal(err)
+			}
+			trace := filepath.Join(tools, "trace")
+			c := runnerConfig{Home: dir, UID: os.Getuid(), ToolDirs: []string{tools}, Executables: map[string]string{}}
+			envGuard := `[ "${GH_TOKEN-unset}|${GIT_CONFIG_COUNT-unset}|${AILANG_REGISTRY_API_KEY-unset}" = unset\|unset\|unset ]
+[ "$GH_PROMPT_DISABLED|$GIT_TERMINAL_PROMPT" = 1\|0 ]
+`
+			writeFinishTestTool(t, &c, tools, "gh", envGuard+fmt.Sprintf(`
+case "$*" in
+  'auth status --hostname github.com') ;;
+  'label list -R fixture/approved --search ailang-world --json name') ;;
+  'repo view fixture/approved --json viewerPermission --jq .viewerPermission') ;;
+  *) printf forbidden >> %s; exit 9;;
+esac
+printf '%%s ' "$1" >> %s
+if [ "$1" = repo ]; then printf '%%s\n' %s; fi
+[ "$1" != %s ]
+`, shellQuote(trace), shellQuote(trace), shellQuote(tc.permission), shellQuote(tc.ghFailure)))
+			if tc.name == "keychain" {
+				c.FinishHTTPSHelper = "osxkeychain"
+				writeFinishTestTool(t, &c, tools, "git-credential-osxkeychain", "exit 1\n")
+			}
+			gitArgs := "--git-dir=/dev/null -c credential.interactive=never -c credential.helper= "
+			if strings.HasPrefix(tc.remote, "https://") {
+				gitArgs += "-c credential.helper=!" + shellQuote(filepath.Join(tools, "gh")) + " auth git-credential"
+				if tc.name == "keychain" {
+					gitArgs = "--git-dir=/dev/null -c credential.interactive=never -c credential.helper= -c credential.helper=!" +
+						shellQuote(filepath.Join(tools, "git-credential-osxkeychain"))
+				}
+			} else {
+				gitArgs += "-c core.sshCommand=/usr/bin/ssh -F /dev/null -o BatchMode=yes -o StrictHostKeyChecking=yes"
+			}
+			gitArgs += " ls-remote --symref -- " + tc.remote + " HEAD"
+			if tc.gitResult == "" {
+				tc.gitResult = "exec " + shellQuote(git) + " --git-dir=/dev/null ls-remote --symref -- " + shellQuote(bare) + " HEAD"
+			}
+			writeFinishTestTool(t, &c, tools, "git", envGuard+fmt.Sprintf(`
+[ "$GIT_CONFIG_NOSYSTEM|$GIT_CONFIG_GLOBAL" = 1\|/dev/null ]
+[ "${SSH_AUTH_SOCK-unset}" = unset ]
+[ "$*" = %s ] || { printf forbidden >> %s; exit 9; }
+printf 'git ' >> %s
+%s
+`, shellQuote(gitArgs), shellQuote(trace), shellQuote(trace), tc.gitResult))
+			if tc.step == 2 {
+				delete(c.Executables, filepath.Join(tools, "gh"))
+			}
+			err := prepareOperatorFinish(boundedTestContext(t), c, tc.remote)
+			if tc.step == 0 && err != nil || tc.step != 0 && (!errors.Is(err, tc.step) || !strings.Contains(err.Error(), "blocked: github-auth")) {
+				t.Fatalf("want readiness %d, got %v", tc.step, err)
+			}
+			want := map[operatorGitHubAuth]string{2: "", 3: "auth ", 4: "auth label ", 5: "auth label repo "}[tc.step]
+			if tc.step == 0 || tc.step == 6 {
+				want = "auth label repo git "
+			}
+			if got, err := os.ReadFile(trace); string(got) != want || err != nil && !(tc.step == 2 && errors.Is(err, os.ErrNotExist)) {
+				t.Fatalf("child calls %q (%v), want %q", got, err, want)
+			}
+		})
+	}
+}
+
+func TestOperatorFinishAuthInstall(t *testing.T) {
+	f := operatorService(t)
+	f.config.FinishRemotes = []string{fixtureFinishRemote}
+	tools := f.config.ToolDirs[0]
+	trace := filepath.Join(f.home, "readiness")
+	writeFinishTestTool(t, &f.config, tools, "gh", "printf ready > "+shellQuote(trace)+"\nexit 1\n")
+	writeFinishTestTool(t, &f.config, tools, "git", "exit 99\n")
+	observe := f.env.observe
+	f.env.observe = func(path string) {
+		got, err := os.ReadFile(trace)
+		if string(got) != "ready" || err != nil {
+			t.Fatalf("installation wrote %s before auth was observed", path)
+		}
+		observe(path)
+	}
+	if rc, out := f.install(t); rc != exitOK || !strings.Contains(out, "finish blocked: github-auth (readiness step 3)") {
+		t.Fatalf("install %d: %s", rc, out)
+	}
+	if err := os.Remove(filepath.Join(tools, "gh")); err != nil {
+		t.Fatal(err)
+	}
+	if rc, out := f.command("start"); rc != exitOK {
+		t.Fatalf("start after finish auth failure %d: %s", rc, out)
+	}
+	if rc, out := f.command("status"); rc != exitOK || !strings.Contains(out, "start: ready") || !strings.Contains(out, "finish blocked: github-auth") {
+		t.Fatalf("status %d: %s", rc, out)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_preflight.go b/cmd/ailang-worldd/operator_finish_preflight.go
new file mode 100644
index 0000000..d06cc9f
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_preflight.go
@@ -0,0 +1,55 @@
+package main
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// operatorFinishPrep binds a current durable claim to independently resolved
+// inputs. Only the proof returned by the runner's recheck can seal it.
+type operatorFinishPrep struct {
+	action    approval.OperatorAction
+	row       store.OperatorActionRow
+	authority *authority.VerifiedOperatorAction
+}
+
+func verifyOperatorFinishPrep(ctx context.Context, ro *store.Credentials, frozen approval.OperatorAction, now int64,
+	cas store.OperatorCAS, owner string) (operatorFinishPrep, error) {
+	if frozen.Kind != "finish" || !frozen.Options.ConfirmOnSurface || frozen.ID != cas.ID {
+		return operatorFinishPrep{}, approval.ErrNotVerified
+	}
+	row, err := ro.OperatorAction(ctx, cas.ID, now)
+	if err != nil {
+		return operatorFinishPrep{}, err
+	}
+	phase := row.State == "confirmed"
+	if owner != "" {
+		phase = row.State == "claimed" || row.State == "running"
+	}
+	if !phase || row.Owner != owner || row.Generation != cas.Generation || row.Kind != frozen.Kind ||
+		row.Episode != frozen.Episode || len(row.Challenge) != 32 || !bytes.Equal(row.Digest, cas.Digest[:]) {
+		return operatorFinishPrep{}, approval.ErrNotVerified
+	}
+	if now < row.CreatedAt || now >= row.ExpiresAt {
+		return operatorFinishPrep{}, store.ErrOperatorExpired
+	}
+	raw, err := approval.CanonicalOperatorAction(frozen)
+	if err != nil || sha256.Sum256(raw) != cas.Digest || !bytes.Equal(raw, row.Action) {
+		return operatorFinishPrep{}, approval.ErrNotVerified
+	}
+	proof, err := authority.RecheckClaimedOperatorAction(ctx, ro, frozen, [32]byte(row.Challenge), row.ExpiresAt, now, row)
+	if err != nil {
+		return operatorFinishPrep{}, approval.ErrNotVerified
+	}
+	return operatorFinishPrep{frozen, row, proof}, nil
+}
+
+func (p operatorFinishPrep) claimed() bool {
+	return p.row.Owner != "" && len(p.row.Challenge) == 32 && p.authority.Rechecked() &&
+		p.authority.Binds(p.action, [32]byte(p.row.Challenge), p.row.ExpiresAt) == nil
+}
diff --git a/cmd/ailang-worldd/operator_finish_preflight_test.go b/cmd/ailang-worldd/operator_finish_preflight_test.go
new file mode 100644
index 0000000..7cd1b0e
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_preflight_test.go
@@ -0,0 +1,177 @@
+package main
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"database/sql"
+	"errors"
+	"io"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorFinishF1(t *testing.T) {
+	f := newFinishFixture(t)
+	t.Setenv("WORLD_OPERATOR_CONFIG_SHA256", strings.Repeat("a", 64))
+	t.Setenv("WORLD_OPERATOR_FINISH_PREP", fixtureActionID)
+	head := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "--path-format=absolute", "--git-path", "HEAD"))
+	common := filepath.Join(f.root, "repo", ".git")
+	snapshot := func() []any {
+		return []any{
+			serviceFileHash(head), serviceFileHash(filepath.Join(common, "packed-refs")),
+			provisionFiles(t, filepath.Join(common, "refs")),
+			serviceFileHash(filepath.Join(filepath.Dir(head), "index")),
+			serviceFileHash(f.health.CredentialsDB),
+			serviceFileHash(filepath.Join(f.sessions, "ep1.pr-body.md")),
+		}
+	}
+	before := snapshot()
+	oldGit := sessionGit
+	t.Cleanup(func() { sessionGit = oldGit })
+	calls := 0
+	sessionGit = func() (string, error) { calls++; return f.git, nil }
+	for _, flags := range [][]string{nil, {"--confirm-on-surface"}, {"--operator-prep", fixtureActionID}, {"--yes"}, {"--confirm-on-surface", "--dry-run"}} {
+		rc, _, errOut := f.run(noTerminalEnv(1), flags...)
+		if rc != exitUsage || !strings.Contains(errOut, "refusing: no controlling terminal") ||
+			calls != 0 || len(f.requests) != 0 || len(f.seam.calls) != 0 || f.checks != 0 || f.outcomes != 0 ||
+			len(f.finishRecords) != 0 || !reflect.DeepEqual(before, snapshot()) {
+			t.Fatalf("F1 flags %q: rc %d, git %d, requests %v, dispatch %v, %s", flags, rc, calls, f.requests, f.seam.calls, errOut)
+		}
+	}
+}
+
+type operatorFinishTTY struct {
+	io.Reader
+	bytes.Buffer
+}
+
+func (t *operatorFinishTTY) Read(b []byte) (int, error) { return t.Reader.Read(b) }
+func (t *operatorFinishTTY) Close() error               { return nil }
+
+func TestOperatorFinishTTYParity(t *testing.T) {
+	f := newFinishFixture(t)
+	tty := &operatorFinishTTY{Reader: strings.NewReader("y\n")}
+	env := testMintEnv("", 1)
+	env.openTerminal = func() (io.ReadWriteCloser, error) { return tty, nil }
+	rc, _, errOut := f.run(env, "--include-unattributed", "--title", "Reviewed fixture")
+	if rc != exitOK || len(f.scopes) != 1 || len(f.finishRecords) != 1 {
+		t.Fatalf("typed finish: rc %d scopes %d records %d %s", rc, len(f.scopes), len(f.finishRecords), errOut)
+	}
+	scope, err := broker.ParseFinishApprovalScope(f.scopes[0])
+	if err != nil {
+		t.Fatal(err)
+	}
+	want := broker.FinishPromptText(scope, f.seam.title, f.seam.body, nil, "") + "Open this pull request? [y/N] "
+	if tty.String() != want {
+		t.Fatalf("terminal prompt changed\ngot %q\nwant %q", tty.String(), want)
+	}
+}
+
+func TestOperatorFinishNoProof(t *testing.T) {
+	dir, err := filepath.EvalSymlinks(t.TempDir())
+	if err != nil {
+		t.Fatal(err)
+	}
+	path := filepath.Join(dir, "creds.db")
+	db, err := store.OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer db.Close()
+	ctx := boundedTestContext(t)
+	frozen := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "finish",
+		Episode: "ep1", Store: path, Options: approval.OperatorOptions{Preset: "fixture", ConfirmOnSurface: true},
+		Review: "finish fixture", Resolution: []approval.OperatorResolved{{Name: "manifest", Value: "fixture"}}}
+	frozen.Argv = approval.OperatorArgv(frozen.Kind, frozen.Episode, frozen.Options)
+	request := store.OperatorRequest{Action: frozen, OperationID: strings.Repeat("1", 32),
+		Challenge: sha256.Sum256([]byte(t.Name())), CreatedAt: 1000, ExpiresAt: 1600}
+	row, err := db.RequestOperatorAction(ctx, request, 1000)
+	if err != nil {
+		t.Fatal(err)
+	}
+	key := approvaltest.NewSigner(51)
+	if err := db.EnrolApprover(ctx, store.ApproverRow{Name: "fixture", Format: "ssh", PublicKey: key.PublicKey().Blob, EnrolledAt: 900}); err != nil {
+		t.Fatal(err)
+	}
+	cas := operatorPreparationCAS(row)
+	refuse := func(action approval.OperatorAction, at int64, cas store.OperatorCAS, owner string, want error) {
+		t.Helper()
+		before := serviceFileHash(path)
+		got, err := verifyOperatorFinishPrep(ctx, db, action, at, cas, owner)
+		if !errors.Is(err, want) || got.claimed() || serviceFileHash(path) != before {
+			t.Fatalf("want %v without state change, got %v", want, err)
+		}
+	}
+	refuse(frozen, 1010, cas, "", approval.ErrNotVerified)
+	if (operatorFinishPrep{}).claimed() {
+		t.Fatal("zero-value capability opened the claimed gate")
+	}
+
+	msg, err := approval.OperatorActionMessage(cas.Digest, request.Challenge, request.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	sig := key.Sign(approval.OperatorActionNamespace, msg)
+	proof, err := authority.JudgeOperatorAction(ctx, db, authority.OperatorJudgeInput{Action: frozen, Digest: cas.Digest,
+		BytesDigest: sha256.Sum256(row.Action), Challenge: request.Challenge, CreatedAt: 1000, ExpiresAt: 1600, Now: 1010}, sig)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if _, err := db.ConfirmOperatorAction(ctx, cas, sig, proof, 1010); err != nil {
+		t.Fatal(err)
+	}
+	writer, err := sql.Open("sqlite", runnerDB(path))
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer writer.Close()
+	setArtifact := func(artifact []byte) {
+		t.Helper()
+		if _, err := writer.ExecContext(ctx, "UPDATE operator_actions SET artifact=? WHERE id=?", artifact, cas.ID); err != nil {
+			t.Fatal(err)
+		}
+	}
+	setArtifact(approvaltest.NewSigner(52).Sign(approval.OperatorActionNamespace, msg))
+	refuse(frozen, 1010, cas, "", approval.ErrNotVerified)
+	setArtifact(sig)
+	verified, err := verifyOperatorFinishPrep(ctx, db, frozen, 1010, cas, "")
+	if err != nil || verified.claimed() {
+		t.Fatalf("confirmed proof must precede claim: %v", err)
+	}
+	foreign := frozen
+	foreign.Episode = "ep-other"
+	foreign.Argv = approval.OperatorArgv(foreign.Kind, foreign.Episode, foreign.Options)
+	refuse(foreign, 1010, cas, "", approval.ErrNotVerified)
+	foreign = frozen
+	foreign.Resolution = []approval.OperatorResolved{{Name: "manifest", Value: "different"}}
+	refuse(foreign, 1010, cas, "", approval.ErrNotVerified)
+
+	h := store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "test-runner"}, Protocol: store.OperatorRunnerProtocol, Build: "fixture", Readiness: "both", SeenAt: 1010}
+	if h.Owner.Generation, err = db.PutOperatorHeartbeat(ctx, h); err != nil {
+		t.Fatal(err)
+	}
+	if err := db.ClaimOperatorAction(ctx, cas, h.Owner, 1011); err != nil {
+		t.Fatal(err)
+	}
+	refuse(frozen, 1011, cas, "", approval.ErrNotVerified)
+	refuse(frozen, 1011, cas, h.Owner.Nonce, approval.ErrNotVerified)
+	cas.Generation++
+	refuse(frozen, 1011, cas, "other-runner", approval.ErrNotVerified)
+	refuse(frozen, 1600, cas, h.Owner.Nonce, store.ErrOperatorExpired)
+	verified, err = verifyOperatorFinishPrep(ctx, db, frozen, 1011, cas, h.Owner.Nonce)
+	if err != nil || !verified.claimed() {
+		t.Fatalf("live held proof failed: %v", err)
+	}
+	verified.authority = proof
+	if verified.claimed() {
+		t.Fatal("daemon's confirmation proof skipped independent runner recheck")
+	}
+}
diff --git a/cmd/ailang-worldd/operator_service.go b/cmd/ailang-worldd/operator_service.go
index d915f23..4230030 100644
--- a/cmd/ailang-worldd/operator_service.go
+++ b/cmd/ailang-worldd/operator_service.go
@@ -24,19 +24,23 @@ const runnerSystemPath = "/usr/bin:/bin:/usr/sbin:/sbin"
 const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | reconcile-mint ACTION-ID | reconcile-worktree ACTION-ID | reconcile-files ACTION-ID | reconcile-start ACTION-ID | launch-status ACTION-ID | uninstall
 Install, stop, cancel, orphan mint revocation, worktree and file disposition and uninstall require typed y. Start uses the approved configuration.
 Config JSON keys: home, uid, address, world, credentials, executables (absolute path to sha256 digest), preset_names, tool_dirs, backend ("headless").
+Optional finish readiness keys: finish_remotes (approved GitHub URLs), finish_ssh_socket (approved login agent socket), finish_https_helper ("gh" or "osxkeychain"). Pin gh, git and any separate HTTPS helper.
 The World service must already be running. Use absolute, canonical paths and include its installed binary in executables.
 `
 
 type runnerConfig struct {
-	Home        string            `json:"home"`
-	UID         int               `json:"uid"`
-	Address     string            `json:"address"`
-	World       string            `json:"world"`
-	Credentials string            `json:"credentials"`
-	Executables map[string]string `json:"executables"`
-	PresetNames []string          `json:"preset_names"`
-	ToolDirs    []string          `json:"tool_dirs"`
-	Backend     string            `json:"backend"`
+	Home              string            `json:"home"`
+	UID               int               `json:"uid"`
+	Address           string            `json:"address"`
+	World             string            `json:"world"`
+	Credentials       string            `json:"credentials"`
+	Executables       map[string]string `json:"executables"`
+	PresetNames       []string          `json:"preset_names"`
+	ToolDirs          []string          `json:"tool_dirs"`
+	Backend           string            `json:"backend"`
+	FinishRemotes     []string          `json:"finish_remotes,omitempty"`
+	FinishHTTPSHelper string            `json:"finish_https_helper,omitempty"`
+	FinishSSHSocket   string            `json:"finish_ssh_socket,omitempty"`
 }
 
 func (e serviceEnv) runnerPaths() servicePaths {
@@ -142,10 +146,25 @@ func (e serviceEnv) checkRunnerConfig(c runnerConfig) error {
 			return errors.New("invalid runner preset name")
 		}
 	}
+	if c.FinishHTTPSHelper != "" && c.FinishHTTPSHelper != "gh" && c.FinishHTTPSHelper != "osxkeychain" {
+		return errors.New("invalid approved finish HTTPS helper")
+	}
+	if len(c.FinishRemotes) > 16 || c.FinishSSHSocket != "" && (!filepath.IsAbs(c.FinishSSHSocket) || filepath.Clean(c.FinishSSHSocket) != c.FinishSSHSocket) {
+		return errors.New("invalid finish readiness configuration")
+	}
+	for _, remote := range c.FinishRemotes {
+		if _, err := finishRepository(remote); err != nil || len(remote) > 512 {
+			return errors.New("invalid approved finish remote")
+		}
+	}
 	if c.Executables[e.paths.Binary] == "" {
 		return errors.New("missing installed World executable")
 	}
 	for path, sum := range c.Executables {
+		name := filepath.Base(path)
+		if len(c.FinishRemotes) != 0 && (name == "gh" || name == "git-credential-osxkeychain") {
+			continue // The per-verb readiness gate pins finish tools.
+		}
 		if err := runnerPath(home, path, uid, false, 0); err != nil {
 			return err
 		}
@@ -325,6 +344,14 @@ func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
 			if !e.confirm(errOut, change) {
 				return 1
 			}
+			ctx, done := e.operatorContext(30 * time.Second)
+			for _, remote := range c.FinishRemotes {
+				if err := prepareOperatorFinish(ctx, c, remote); err != nil {
+					fmt.Fprintln(out, err)
+					break
+				}
+			}
+			done()
 			lock, err := runnerLifecycleLock(filepath.Dir(p.Record))
 			if err != nil {
 				return refuse(errors.New("runner lifecycle is busy or unavailable"))
diff --git a/cmd/ailang-worldd/operator_service_test.go b/cmd/ailang-worldd/operator_service_test.go
index d5b7d01..5be8d49 100644
--- a/cmd/ailang-worldd/operator_service_test.go
+++ b/cmd/ailang-worldd/operator_service_test.go
@@ -50,7 +50,9 @@ func operatorService(t *testing.T) *operatorServiceFixture {
 	world, creds := r.Store, r.Credentials // fixture paths are physically named
 	f.env.runnerLogin = func() (int, string, error) { return os.Getuid(), f.home, nil }
 	o := &operatorServiceFixture{serviceFixture: f, homeJob: f.ctl, configFile: filepath.Join(f.home, "proposed.json"),
-		config:    runnerConfig{f.home, os.Getuid(), r.Address, world, creds, map[string]string{r.Binary: r.BinarySHA256}, []string{"sol"}, []string{filepath.Dir(r.Binary)}, "headless"},
+		config: runnerConfig{Home: f.home, UID: os.Getuid(), Address: r.Address, World: world, Credentials: creds,
+			Executables: map[string]string{r.Binary: r.BinarySHA256}, PresetNames: []string{"sol"},
+			ToolDirs: []string{filepath.Dir(r.Binary)}, Backend: "headless"},
 		heartbeat: store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "fake-runner"}, Protocol: store.OperatorRunnerProtocol, Build: r.BinarySHA256},
 		health:    daemon.HealthResponse{Status: "ok", DBPath: world, CredentialsDB: creds, Build: daemon.BuildIdentity{ExecutableSHA256: r.BinarySHA256}}}
 	f.ctl = &fakeServiceController{runner: true, binary: r.Binary}
diff --git a/host/checkrun/checkrun.go b/host/checkrun/checkrun.go
index 706440d..a646267 100644
--- a/host/checkrun/checkrun.go
+++ b/host/checkrun/checkrun.go
@@ -89,6 +89,15 @@ func Run(ctx context.Context, argv []string, dir string, timeout time.Duration)
 
 // RunEnv overlays the pinned profile on the scrubbed inherited environment.
 func RunEnv(ctx context.Context, argv []string, dir string, timeout time.Duration, env map[string]string, pathPrefix []string) (Result, error) {
+	return run(ctx, argv, dir, timeout, checkEnv(os.Environ(), env, pathPrefix))
+}
+
+// RunIsolated uses only the approved child environment.
+func RunIsolated(ctx context.Context, argv []string, dir string, timeout time.Duration, env []string) (Result, error) {
+	return run(ctx, argv, dir, timeout, checkEnv(env, nil, nil))
+}
+
+func run(ctx context.Context, argv []string, dir string, timeout time.Duration, env []string) (Result, error) {
 	r := Result{ExitCode: -1}
 	if len(argv) == 0 || timeout <= 0 {
 		return r, fmt.Errorf("checkrun: argv and a positive wall cap are required")
@@ -97,7 +106,7 @@ func RunEnv(ctx context.Context, argv []string, dir string, timeout time.Duratio
 	defer cancel()
 	cmd := exec.CommandContext(runCtx, argv[0], argv[1:]...)
 	cmd.Dir = dir
-	cmd.Env = checkEnv(os.Environ(), env, pathPrefix)
+	cmd.Env = env
 	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
 	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
 	cmd.WaitDelay = 100 * time.Millisecond
diff --git a/host/worktree/worktree.go b/host/worktree/worktree.go
index 779282c..4fb871f 100644
--- a/host/worktree/worktree.go
+++ b/host/worktree/worktree.go
@@ -90,6 +90,18 @@ func AddExisting(ctx context.Context, git, repo, dir, branch string) error {
 // Output runs `git args...` and returns its standard output. Stderr is kept
 // only for the error message; stdout over MaxStdout is an error.
 func Output(ctx context.Context, git string, args []string) (string, error) {
+	return gitOutput(ctx, git, args)
+}
+
+type isolatedGitEnvKey struct{}
+
+// OutputIsolated uses a trusted caller's explicit environment at the same
+// bounded Git site. Unlike ambient Git, caller-supplied policy can set GIT_*.
+func OutputIsolated(ctx context.Context, git string, args, env []string) (string, error) {
+	return gitOutput(context.WithValue(ctx, isolatedGitEnvKey{}, childenv.Scrubbed(env)), git, args)
+}
+
+func gitOutput(ctx context.Context, git string, args []string) (string, error) {
 	out := &capped{max: MaxStdout}
 	if err := runGit(ctx, git, args, out); err != nil {
 		return "", err
@@ -155,7 +167,11 @@ func runGit(ctx context.Context, git string, args []string, stdout *capped, inhe
 	runCtx, cancel := context.WithTimeout(ctx, timeout)
 	defer cancel()
 	cmd := exec.CommandContext(runCtx, git, args...)
-	cmd.Env = Env(os.Environ())
+	if approved, ok := ctx.Value(isolatedGitEnvKey{}).([]string); ok {
+		cmd.Env = approved
+	} else {
+		cmd.Env = Env(os.Environ())
+	}
 	cmd.Stdin = nil
 	cmd.ExtraFiles = inherited
 	cmd.WaitDelay = waitDelay
````

### M2 — Frozen inputs and local Git receipt primitives

````diff M2
diff --git a/cmd/ailang-worldd/operator_finish_auth.go b/cmd/ailang-worldd/operator_finish_auth.go
index f9b7cfc..e29ab0e 100644
--- a/cmd/ailang-worldd/operator_finish_auth.go
+++ b/cmd/ailang-worldd/operator_finish_auth.go
@@ -110,7 +110,7 @@ func (a operatorFinishAuth) prepare(ctx context.Context) error {
 		}
 		args = append(args, "-c", "core.sshCommand=/usr/bin/ssh -F /dev/null -o BatchMode=yes -o StrictHostKeyChecking=yes")
 	}
-	out, err := worktree.OutputIsolated(ctx, a.Git, append(args, "ls-remote", "--symref", "--", a.Remote, "HEAD"), env)
+	out, err := worktree.RunIsolated(ctx, a.Git, append(args, "ls-remote", "--symref", "--", a.Remote, "HEAD"), env, nil)
 	fields := strings.Fields(out)
 	if err != nil || len(fields) != 5 || fields[0] != "ref:" || !strings.HasPrefix(fields[1], "refs/heads/") ||
 		fields[2] != "HEAD" || !operatorFinishSHA.MatchString(fields[3]) || fields[4] != "HEAD" {
diff --git a/cmd/ailang-worldd/operator_finish_auth_test.go b/cmd/ailang-worldd/operator_finish_auth_test.go
index a5d0c82..cfa85c7 100644
--- a/cmd/ailang-worldd/operator_finish_auth_test.go
+++ b/cmd/ailang-worldd/operator_finish_auth_test.go
@@ -36,8 +36,8 @@ func TestOperatorFinishAuth(t *testing.T) {
 	ctx := boundedTestContext(t)
 	gitRun := func(args ...string) {
 		t.Helper()
-		if _, err := worktree.OutputIsolated(ctx, git, args, []string{"HOME=" + dir, "PATH=" + runnerSystemPath,
-			"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}); err != nil {
+		if _, err := worktree.RunIsolated(ctx, git, args, []string{"HOME=" + dir, "PATH=" + runnerSystemPath,
+			"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}, nil); err != nil {
 			t.Fatal(err)
 		}
 	}
diff --git a/cmd/ailang-worldd/operator_finish_manifest.go b/cmd/ailang-worldd/operator_finish_manifest.go
new file mode 100644
index 0000000..d5ea314
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_manifest.go
@@ -0,0 +1,85 @@
+package main
+
+import (
+	"context"
+	"encoding/json"
+	"fmt"
+	"slices"
+	"strconv"
+	"strings"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+type operatorFinishInputs struct {
+	manifest store.OperatorFinishReview
+	report   daemon.EpisodeReport
+	git      operatorFinishGit
+	token    string
+	profile  *broker.ExecProfile
+	world    string
+}
+
+func (i operatorFinishInputs) frozen(ctx context.Context) (store.OperatorFinishReview, error) {
+	v := i.manifest
+	tree, _, err := i.git.snapshot(ctx, i.report.Start)
+	if err != nil {
+		return v, err
+	}
+	if !i.baseMatches(ctx) || len(tree.Paths) == 0 {
+		return v, errOperatorFinishReview
+	}
+	if err := validateFinishChecks(v.Checks, nil, i.git.dir); err != nil {
+		return v, errOperatorFinishFence
+	}
+	var paths []string
+	for _, p := range tree.Paths {
+		paths = append(paths, p.Path)
+	}
+	v.Protected = finishProtectedPaths(i.profile, paths)
+	v.Unattributed = finishUnattributed(paths, i.report)
+	v.Commit.Snapshot = tree
+	v.Commit.Message += fmt.Sprintf("\n\n%s %s; %d calls; %d/%d checks passed; %d files changed.\n\n%s\n",
+		i.report.Start.Agent, i.report.Start.Model, i.report.TotalCalls, len(v.Checks), len(v.Checks), len(paths), reportTrailer(i.report, i.world))
+	if scrubMarkdown(v.Commit.Message) != nil {
+		return v, errOperatorFinishFence
+	}
+	return v, nil
+}
+
+func (i operatorFinishInputs) baseMatches(ctx context.Context) bool {
+	if _, err := i.git.config(ctx, i.report.Start); err != nil {
+		return false
+	}
+	base, err := i.git.git(ctx, nil, nil, "rev-parse", "--verify", "refs/heads/"+i.manifest.Base.Branch+"^{commit}")
+	return err == nil && strings.TrimSpace(base) == i.manifest.Base.SHA
+}
+
+func finishManifestForRow(c runnerConfig, row store.OperatorActionRow) (approval.OperatorAction, store.OperatorFinishReview, int64, error) {
+	a, err := approval.DecodeOperatorAction(row.Action)
+	if err != nil || a.Kind != "finish" || a.Episode != row.Episode || a.Store != c.Credentials {
+		return approval.OperatorAction{}, store.OperatorFinishReview{}, 0, approval.ErrNotVerified
+	}
+	var m store.OperatorFinishReview
+	generation, genErr := strconv.ParseInt(operatorFileFact(a, "finish_catalogue"), 10, 64)
+	if json.Unmarshal([]byte(operatorFileFact(a, "finish_manifest")), &m) != nil || genErr != nil || generation <= 0 ||
+		!slices.Contains(c.FinishRemotes, m.Remote.URL) {
+		return approval.OperatorAction{}, store.OperatorFinishReview{}, 0, approval.ErrNotVerified
+	}
+	return a, m, generation, nil
+}
+
+func finishWorldEvidence(bound context.Context, path string, receipt store.OperatorFinishCheck) bool {
+	hash, err := hashref.Parse(receipt.Ref)
+	world, openErr := store.OpenReadOnly(path)
+	if openErr != nil {
+		return false
+	}
+	defer world.Close()
+	object, found, readErr := world.GetObject(bound, hash)
+	return err == nil && readErr == nil && found && hashref.SumSHA256(object.Payload) == hash
+}
diff --git a/cmd/ailang-worldd/operator_finish_tree.go b/cmd/ailang-worldd/operator_finish_tree.go
new file mode 100644
index 0000000..fd9bf62
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_tree.go
@@ -0,0 +1,353 @@
+package main
+
+import (
+	"context"
+	"errors"
+	"fmt"
+	"io"
+	"io/fs"
+	"os"
+	"path/filepath"
+	"reflect"
+	"slices"
+	"strconv"
+	"strings"
+	"unicode/utf8"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/worktree"
+	"golang.org/x/sys/unix"
+)
+
+const operatorFinishDirtyLimit = 16 << 20
+
+var (
+	errOperatorFinishFence  = errors.New("worktree fence: operator input required")
+	errOperatorFinishReview = errors.New("finish inputs changed; request fresh review")
+	errOperatorCommitNoRef  = errors.New("local commit stopped before index or ref")
+)
+
+type operatorFinishGit struct {
+	bin, dir, admin, common string
+	env                     []string
+	lock                    *os.File
+}
+
+func (g operatorFinishGit) git(ctx context.Context, input io.Reader, extra []string, args ...string) (string, error) {
+	argv := []string{"--literal-pathspecs", "--git-dir=" + g.admin, "--work-tree=" + g.dir,
+		"-c", "core.hooksPath=/dev/null", "-c", "core.fsmonitor=false", "-c", "core.attributesFile=/dev/null", "-c", "commit.gpgsign=false"}
+	env := append(slices.Clone(g.env), extra...)
+	var inherited []*os.File
+	if g.lock != nil {
+		inherited = append(inherited, g.lock)
+	}
+	return worktree.RunIsolated(ctx, g.bin, append(argv, args...), env, input, inherited...)
+}
+
+// Open beneath a held directory, and compare the file to its final pathname.
+func operatorFinishRead(dir, name string, maxBytes int64) ([]byte, os.FileMode, error) {
+	fd, err := unix.Open(dir, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
+	if err != nil {
+		return nil, 0, err
+	}
+	parts := strings.Split(name, string(filepath.Separator))
+	for i, part := range parts {
+		flags := unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_NONBLOCK | unix.O_CLOEXEC
+		if i+1 < len(parts) {
+			flags |= unix.O_DIRECTORY
+		}
+		next, err := unix.Openat(fd, part, flags, 0)
+		unix.Close(fd)
+		if err != nil {
+			return nil, 0, err
+		}
+		fd = next
+	}
+	file := os.NewFile(uintptr(fd), name)
+	defer file.Close()
+	before, err := file.Stat()
+	if err != nil || !before.Mode().IsRegular() || before.Size() > maxBytes {
+		return nil, 0, errOperatorFinishFence
+	}
+	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
+	after, pathErr := os.Lstat(filepath.Join(dir, name))
+	if err != nil || pathErr != nil || int64(len(data)) != before.Size() || int64(len(data)) > maxBytes ||
+		!os.SameFile(before, after) || before.Mode() != after.Mode() ||
+		before.Size() != after.Size() || before.ModTime() != after.ModTime() {
+		return nil, 0, errOperatorFinishFence
+	}
+	return data, before.Mode(), nil
+}
+
+func operatorFinishPointer(dir, name, prefix string) (string, error) {
+	b, _, err := operatorFinishRead(dir, name, 4096)
+	path, ok := strings.CutPrefix(strings.TrimSuffix(string(b), "\n"), prefix)
+	if err != nil || !ok || path == "" || strings.ContainsAny(path, "\x00\r\n") {
+		return "", errOperatorFinishFence
+	}
+	if !filepath.IsAbs(path) {
+		path = filepath.Join(dir, path)
+	}
+	return filepath.EvalSymlinks(path)
+}
+
+func openOperatorFinishGit(auth operatorFinishAuth, dir, repo string, lock *os.File) (operatorFinishGit, error) {
+	var g operatorFinishGit
+	dir, err := filepath.EvalSymlinks(dir)
+	if err != nil {
+		return g, errOperatorFinishFence
+	}
+	repo, err = filepath.EvalSymlinks(repo)
+	if err != nil {
+		return g, errOperatorFinishFence
+	}
+	admin, err := operatorFinishPointer(dir, ".git", "gitdir: ")
+	if err != nil || filepath.Dir(admin) != filepath.Join(repo, ".git", "worktrees") {
+		return g, errOperatorFinishFence
+	}
+	common, err := operatorFinishPointer(admin, "commondir", "")
+	link, linkErr := operatorFinishPointer(admin, "gitdir", "")
+	if err != nil || linkErr != nil || common != filepath.Join(repo, ".git") || link != filepath.Join(dir, ".git") {
+		return g, errOperatorFinishFence
+	}
+	return operatorFinishGit{bin: auth.Git, dir: dir, admin: admin, common: common, lock: lock,
+		env: []string{"HOME=" + auth.Home, "PATH=" + auth.Path, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
+			"GIT_ATTR_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0", "GIT_NO_REPLACE_OBJECTS=1", "GIT_NO_LAZY_FETCH=1", "GIT_TERMINAL_PROMPT=0"}}, nil
+}
+
+// Inspect Git configuration as data before running any worktree command.
+func (g operatorFinishGit) config(ctx context.Context, start *store.EpisodeStart) (string, error) {
+	for _, path := range []string{filepath.Join(g.admin, "config.worktree"), filepath.Join(g.common, "config.worktree"),
+		filepath.Join(g.common, "objects", "info", "alternates"), filepath.Join(g.common, "info", "grafts")} {
+		if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
+			return "", errOperatorFinishFence
+		}
+	}
+	data, _, err := operatorFinishRead(g.common, "config", 64<<10)
+	if err != nil {
+		return "", errOperatorFinishFence
+	}
+	raw, err := worktree.RunIsolated(ctx, g.bin, []string{"--git-dir=/dev/null", "config", "--file", "-",
+		"--null", "--no-includes", "--list"}, g.env, strings.NewReader(string(data)))
+	if err != nil {
+		return "", errOperatorFinishFence
+	}
+	remote := 0
+	for item := range strings.SplitSeq(strings.TrimSuffix(raw, "\x00"), "\x00") {
+		name, value, _ := strings.Cut(item, "\n")
+		switch {
+		case name == "core.repositoryformatversion" && value == "0":
+		case name == "core.bare" && value == "false":
+		case slices.Contains([]string{"core.filemode", "core.logallrefupdates", "core.ignorecase", "core.precomposeunicode",
+			"user.name", "user.email", "commit.gpgsign"}, name):
+		case strings.HasPrefix(name, "branch.") && (strings.HasSuffix(name, ".remote") || strings.HasSuffix(name, ".merge")):
+		case name == "remote."+start.Remote.Name+".fetch":
+		case name == "remote."+start.Remote.Name+".url" && value == start.Remote.URL:
+			remote++
+		default:
+			return "", errOperatorFinishFence
+		}
+	}
+	if remote != 1 {
+		return "", errOperatorFinishFence
+	}
+	return finishSHA(data), nil
+}
+
+func (g operatorFinishGit) snapshot(ctx context.Context, start *store.EpisodeStart) (store.OperatorFinishTree, [][]byte, error) {
+	var tree store.OperatorFinishTree
+	config, err := g.config(ctx, start)
+	if err != nil {
+		return tree, nil, err
+	}
+	branch, err := g.git(ctx, nil, nil, "symbolic-ref", "HEAD")
+	if err != nil || strings.TrimSpace(branch) != "refs/heads/"+start.Branch {
+		return tree, nil, errOperatorFinishFence
+	}
+	head, err := g.git(ctx, nil, nil, "rev-parse", "--verify", "HEAD")
+	if err != nil || !operatorFinishSHA.MatchString(strings.TrimSpace(head)) {
+		return tree, nil, errOperatorFinishFence
+	}
+	index, _, err := operatorFinishRead(g.admin, "index", 64<<20)
+	if err != nil {
+		return tree, nil, errOperatorFinishFence
+	}
+	staged, err := g.git(ctx, nil, nil, "diff", "--cached", "--name-only", "--no-ext-diff", "--no-textconv", "HEAD", "--")
+	if err != nil || staged != "" {
+		return tree, nil, errOperatorFinishFence
+	}
+	flags, err := g.git(ctx, nil, nil, "ls-files", "-v", "-z")
+	if err != nil {
+		return tree, nil, errOperatorFinishFence
+	}
+	trackedPaths := make(map[string]bool)
+	for entry := range strings.SplitSeq(strings.TrimSuffix(flags, "\x00"), "\x00") {
+		if entry == "" {
+			continue
+		}
+		if !strings.HasPrefix(entry, "H ") {
+			return tree, nil, errOperatorFinishFence
+		}
+		trackedPaths[entry[2:]] = true
+	}
+	tracked, err := g.git(ctx, nil, nil, "diff", "--name-only", "--no-renames", "--no-ext-diff", "--no-textconv", "-z", "HEAD", "--")
+	if err != nil {
+		return tree, nil, errOperatorFinishFence
+	}
+	names := strings.Split(strings.TrimSuffix(tracked, "\x00"), "\x00")
+	root, err := os.OpenRoot(g.dir)
+	if err != nil {
+		return tree, nil, errOperatorFinishFence
+	}
+	defer root.Close()
+	visited := 0
+	err = fs.WalkDir(root.FS(), ".", func(path string, d fs.DirEntry, err error) error {
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
+		return tree, nil, errOperatorFinishFence
+	}
+	slices.Sort(names)
+	names = slices.Compact(names)
+	if len(names) > 257 {
+		return tree, nil, errOperatorFinishFence
+	}
+	tree = store.OperatorFinishTree{Head: strings.TrimSpace(head), Branch: strings.TrimSpace(branch), Index: finishSHA(index), Config: config}
+	var contents [][]byte
+	budget := int64(operatorFinishDirtyLimit)
+	for _, path := range names {
+		if path == "" {
+			continue
+		}
+		if !filepath.IsLocal(path) || filepath.Clean(path) != path || len(path) > 1024 || !utf8.ValidString(path) {
+			return store.OperatorFinishTree{}, nil, errOperatorFinishFence
+		}
+		for part := range strings.SplitSeq(path, "/") {
+			if strings.EqualFold(part, ".git") {
+				return store.OperatorFinishTree{}, nil, errOperatorFinishFence
+			}
+		}
+		b, mode, err := operatorFinishRead(g.dir, path, budget)
+		p := store.OperatorFinishPath{Path: path, Mode: "100644", Digest: finishSHA(b)}
+		switch {
+		case errors.Is(err, os.ErrNotExist):
+			p.Mode, p.Digest = "0", ""
+		case err != nil:
+			return store.OperatorFinishTree{}, nil, errOperatorFinishFence
+		case mode&0100 != 0:
+			p.Mode = "100755"
+		}
+		tree.Paths = append(tree.Paths, p)
+		contents = append(contents, b)
+		budget -= int64(len(b))
+	}
+	return tree, contents, nil
+}
+
+func (g operatorFinishGit) stage(ctx context.Context, spec store.OperatorLocalCommit, contents [][]byte, env []string) (string, error) {
+	if _, err := g.git(ctx, nil, env, "read-tree", spec.Snapshot.Head); err != nil {
+		return "", err
+	}
+	var entries strings.Builder
+	for i, p := range spec.Snapshot.Paths {
+		oid := strings.Repeat("0", 40)
+		if p.Mode != "0" {
+			out, err := g.git(ctx, strings.NewReader(string(contents[i])), env, "hash-object", "-w", "--no-filters", "--stdin")
+			oid = strings.TrimSpace(out)
+			if err != nil || !operatorFinishSHA.MatchString(oid) {
+				return "", errOperatorFinishFence
+			}
+		}
+		fmt.Fprintf(&entries, "%s %s\t%s\x00", p.Mode, oid, p.Path)
+	}
+	if _, err := g.git(ctx, strings.NewReader(entries.String()), env, "update-index", "-z", "--index-info"); err != nil {
+		return "", err
+	}
+	out, err := g.git(ctx, nil, env, "write-tree")
+	if !operatorFinishSHA.MatchString(strings.TrimSpace(out)) {
+		return "", errOperatorFinishFence
+	}
+	return strings.TrimSpace(out), err
+}
+
+func (g operatorFinishGit) planTree(ctx context.Context, spool string, spec store.OperatorLocalCommit, contents [][]byte) (string, error) {
+	dir, err := os.MkdirTemp(spool, "finish-tree.")
+	if err != nil {
+		return "", err
+	}
+	defer os.RemoveAll(dir)
+	env := []string{"GIT_INDEX_FILE=" + filepath.Join(dir, "index"), "GIT_OBJECT_DIRECTORY=" + dir,
+		"GIT_ALTERNATE_OBJECT_DIRECTORIES=" + strconv.Quote(filepath.Join(g.common, "objects"))}
+	return g.stage(ctx, spec, contents, env)
+}
+
+func (g operatorFinishGit) inspectCommit(ctx context.Context, spec store.OperatorLocalCommit, tree string) (string, error) {
+	branch, err := g.git(ctx, nil, nil, "symbolic-ref", "HEAD")
+	if err != nil || strings.TrimSpace(branch) != spec.Snapshot.Branch {
+		return "", errOperatorFinishFence
+	}
+	head, err := g.git(ctx, nil, nil, "rev-parse", "--verify", "HEAD")
+	head = strings.TrimSpace(head)
+	if err != nil || !operatorFinishSHA.MatchString(head) || head == spec.Snapshot.Head {
+		return "", errOperatorFinishFence
+	}
+	raw, err := g.git(ctx, nil, nil, "cat-file", "commit", head)
+	if err != nil || !operatorFinishCommitMatches(raw, spec, tree) {
+		return "", errOperatorFinishFence
+	}
+	index, err := g.git(ctx, nil, nil, "diff", "--cached", "--name-only", "--no-ext-diff", "--no-textconv", head, "--")
+	if err != nil || index != "" {
+		return "", errOperatorFinishFence
+	}
+	return head, nil
+}
+
+func operatorFinishCommitMatches(raw string, spec store.OperatorLocalCommit, tree string) bool {
+	header, message, found := strings.Cut(raw, "\n\n")
+	lines := strings.Split(header, "\n")
+	if !found || len(lines) != 4 || lines[0] != "tree "+tree ||
+		lines[1] != "parent "+spec.Snapshot.Head || message != spec.Message {
+		return false
+	}
+	committer, timeFields, ok := strings.Cut(lines[3], "> ")
+	return ok && committer+">" == "committer "+spec.Committer && len(strings.Fields(timeFields)) == 2
+}
+
+func (g operatorFinishGit) commit(ctx context.Context, spool string, spec store.OperatorLocalCommit, contents [][]byte, tree string, start *store.EpisodeStart) error {
+	name, email, err := worktree.ParseCommitter(spec.Committer)
+	if err != nil || name == "" || strings.ContainsAny(start.Agent, "\r\n<>") {
+		return errOperatorCommitNoRef
+	}
+	dir, err := os.MkdirTemp(spool, "finish-commit.")
+	if err != nil {
+		return errOperatorCommitNoRef
+	}
+	defer os.RemoveAll(dir)
+	got, err := g.stage(ctx, spec, contents, []string{"GIT_INDEX_FILE=" + filepath.Join(dir, "index")})
+	if err != nil || got != tree {
+		return errOperatorCommitNoRef
+	}
+	env := []string{"GIT_COMMITTER_NAME=" + name, "GIT_COMMITTER_EMAIL=" + email,
+		"GIT_AUTHOR_NAME=" + start.Agent + " via AILANG World", "GIT_AUTHOR_EMAIL=ailang-world@invalid"}
+	out, err := g.git(ctx, strings.NewReader(spec.Message), env, "commit-tree", tree, "-p", spec.Snapshot.Head, "-F", "-")
+	head := strings.TrimSpace(out)
+	if err != nil || !operatorFinishSHA.MatchString(head) {
+		return errOperatorCommitNoRef
+	}
+	current, _, err := g.snapshot(ctx, start)
+	if err != nil || !reflect.DeepEqual(current, spec.Snapshot) {
+		return errOperatorFinishReview
+	}
+	if _, err := g.git(ctx, nil, nil, "read-tree", tree); err != nil {
+		return err
+	}
+	_, err = g.git(ctx, nil, nil, "update-ref", spec.Snapshot.Branch, head, spec.Snapshot.Head)
+	return err
+}
diff --git a/cmd/ailang-worldd/operator_finish_tree_test.go b/cmd/ailang-worldd/operator_finish_tree_test.go
new file mode 100644
index 0000000..dbf27d1
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_tree_test.go
@@ -0,0 +1,121 @@
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
+	"golang.org/x/sys/unix"
+)
+
+func operatorFinishGitFixture(t *testing.T, f *finishFixture) operatorFinishGit {
+	t.Helper()
+	g, err := openOperatorFinishGit(operatorFinishAuth{Git: f.git, Home: f.root, Path: runnerSystemPath},
+		f.wt, filepath.Join(f.root, "repo"), nil)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return g
+}
+
+func TestOperatorFinishWorktreeFence(t *testing.T) {
+	f := newFinishFixture(t)
+	g := operatorFinishGitFixture(t, f)
+	ctx := boundedTestContext(t)
+	path := filepath.Join(g.common, "config")
+	config := readStartFixtureFile(t, path)
+	if _, _, err := g.snapshot(ctx, f.rep.Start); err != nil {
+		t.Fatal(err)
+	}
+	sentinel := filepath.Join(f.root, "sentinel")
+	payload := "touch " + shellQuote(sentinel)
+	for _, setting := range []string{"core.hooksPath", "filter.fixture.clean", "credential.helper", "core.fsmonitor"} {
+		f.gitOutput("-C", f.wt, "config", setting, "!"+payload)
+		if _, _, err := g.snapshot(boundedTestContext(t), f.rep.Start); !errors.Is(err, errOperatorFinishFence) {
+			t.Fatalf("%s was accepted: %v", setting, err)
+		}
+		seedStartFixtureFile(t, path, config)
+	}
+	bad := filepath.Join(f.wt, "fixture-link")
+	if err := os.Symlink(sentinel, bad); err != nil {
+		t.Fatal(err)
+	}
+	if _, _, err := g.snapshot(boundedTestContext(t), f.rep.Start); !errors.Is(err, errOperatorFinishFence) {
+		t.Fatalf("dirty symlink was accepted: %v", err)
+	}
+	if err := os.Remove(bad); err != nil {
+		t.Fatal(err)
+	}
+	if err := unix.Mkfifo(bad, 0600); err != nil {
+		t.Fatal(err)
+	}
+	if _, _, err := g.snapshot(boundedTestContext(t), f.rep.Start); !errors.Is(err, errOperatorFinishFence) {
+		t.Fatalf("dirty special file was accepted: %v", err)
+	}
+	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("untrusted program executed or sentinel unobservable: %v", err)
+	}
+}
+
+func TestOperatorCommitIdentity(t *testing.T) {
+	f := newFinishFixture(t)
+	g := operatorFinishGitFixture(t, f)
+	ctx := boundedTestContext(t)
+	initial, contents, err := g.snapshot(ctx, f.rep.Start)
+	if err != nil {
+		t.Fatal(err)
+	}
+	spec := store.OperatorLocalCommit{Snapshot: initial, Message: "fixture\n\nWorld-Provenance: episode=ep1\n", Committer: "Fixture Operator <fixture@example.invalid>"}
+	tree, err := g.planTree(ctx, f.root, spec, contents)
+	if err != nil {
+		t.Fatal(err)
+	}
+	seen, _, err := g.snapshot(ctx, f.rep.Start)
+	if err != nil || !reflect.DeepEqual(seen, initial) {
+		t.Fatalf("private staging changed production index or refs: %v", err)
+	}
+	sentinel := filepath.Join(f.root, "git-hook")
+	for _, hook := range []string{"pre-commit", "reference-transaction"} {
+		if err := os.WriteFile(filepath.Join(g.common, "hooks", hook), []byte("#!/bin/sh\ntouch "+shellQuote(sentinel)+"\n"), 0700); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if err := g.commit(ctx, f.root, spec, contents, tree, f.rep.Start); err != nil {
+		t.Fatal(err)
+	}
+	if _, err := os.Lstat(sentinel); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("repo-controlled hook executed or sentinel unreadable: %v", err)
+	}
+	head, err := g.inspectCommit(ctx, spec, tree)
+	if err != nil || head == initial.Head {
+		t.Fatalf("exact descendant missing: %s %v", head, err)
+	}
+	raw := f.gitOutput("-C", f.wt, "cat-file", "commit", head)
+	for _, field := range []string{"parent", "tree", "message", "trailer", "committer"} {
+		bad, want := spec, tree
+		switch field {
+		case "parent":
+			bad.Snapshot.Head = strings.Repeat("d", 40)
+		case "tree":
+			want = head
+		case "message":
+			bad.Message = "bad title\n\nWorld-Provenance: episode=ep1\n"
+		case "trailer":
+			bad.Message = strings.Replace(spec.Message, "episode=ep1", "episode=other", 1)
+		case "committer":
+			bad.Committer = "Someone Else <different@invalid>"
+		}
+		if operatorFinishCommitMatches(raw, bad, want) {
+			t.Fatalf("%s mismatch observed", field)
+		}
+	}
+	f.write("extra", "partial stage")
+	f.gitOutput("-C", f.wt, "add", "--", "extra")
+	if _, err := g.inspectCommit(ctx, spec, tree); err == nil || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != head {
+		t.Fatal("uncertain partial stage adopted or rewrote HEAD")
+	}
+}
diff --git a/host/store/operator_finish.go b/host/store/operator_finish.go
new file mode 100644
index 0000000..2084b9e
--- /dev/null
+++ b/host/store/operator_finish.go
@@ -0,0 +1,135 @@
+package store
+
+import (
+	"context"
+	"database/sql"
+	"encoding/json"
+	"errors"
+	"regexp"
+)
+
+type OperatorFinishCheck struct {
+	Name   string `json:"name"`
+	Result string `json:"result"`
+	Ref    string `json:"ref"`
+	Entry  int64  `json:"entry"`
+}
+
+type OperatorFinishLocalCommit struct {
+	State string `json:"state"`
+	SHA   string `json:"sha,omitempty"`
+}
+
+const OperatorLocalCommitDone = "local-commit-done"
+
+var operatorLocalOID = regexp.MustCompile(`^[0-9a-f]{40}$`)
+
+type OperatorFinishJournal struct {
+	Checks      []OperatorFinishCheck
+	Commit      OperatorFinishLocalCommit
+	CommitKey   OperatorStepKey
+	CommitPhase string
+	OpenCheck   bool
+}
+
+// ReadOperatorFinishJournal projects only typed check and commit receipts.
+// A check reservation without an observation is never permission to replay.
+func (c *Credentials) ReadOperatorFinishJournal(ctx context.Context, action OperatorCAS) (OperatorFinishJournal, error) {
+	var out OperatorFinishJournal
+	if err := c.operatorAccess(ctx); err != nil {
+		return out, err
+	}
+	rows, err := c.db.QueryContext(ctx, `SELECT s.kind,s.identity,s.phase,s.version,s.owner,s.receipt
+		FROM operator_action_steps s JOIN operator_actions a ON a.id=s.action_id
+		WHERE a.id=? AND a.generation=? AND a.digest=? AND a.kind='finish' AND s.kind IN ('checks','commit')
+		ORDER BY s.rowid`, action.ID, action.Generation, action.Digest[:])
+	if err != nil {
+		return out, err
+	}
+	defer rows.Close()
+	for rows.Next() {
+		key := OperatorStepKey{Action: action}
+		var phase string
+		var raw []byte
+		if err := rows.Scan(&key.Kind, &key.Identity, &phase, &key.Version, &key.Owner, &raw); err != nil {
+			return out, err
+		}
+		switch {
+		case key.Kind == "commit":
+			if out.CommitPhase != "" {
+				return OperatorFinishJournal{}, ErrOperatorConflict
+			}
+			if phase == "observed" {
+				if json.Unmarshal(raw, &out.Commit) != nil ||
+					!(out.Commit.State == OperatorLocalCommitDone && operatorLocalOID.MatchString(out.Commit.SHA) ||
+						out.Commit.State == "commit-failed" && out.Commit.SHA == "") {
+					return OperatorFinishJournal{}, ErrOperatorConflict
+				}
+			}
+			out.CommitKey, out.CommitPhase = key, phase
+		case phase == "reserved":
+			out.OpenCheck = true
+		default:
+			var receipt OperatorFinishCheck
+			if json.Unmarshal(raw, &receipt) != nil || receipt.Name == "" || receipt.Entry < 1 || receipt.Ref == "" ||
+				(receipt.Result != "passed" && receipt.Result != "failed") {
+				return OperatorFinishJournal{}, ErrOperatorConflict
+			}
+			out.Checks = append(out.Checks, receipt)
+		}
+	}
+	return out, rows.Err()
+}
+
+func (c *Credentials) NextOperatorFinish(ctx context.Context, now int64) (OperatorActionRow, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorActionRow{}, err
+	}
+	return scanOperatorAction(c.db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions
+		WHERE kind='finish' AND (state='confirmed' AND created_at<=? AND expires_at>? OR state IN ('claimed','running'))
+		AND NOT EXISTS (SELECT 1 FROM operator_action_steps s WHERE s.action_id=operator_actions.id AND s.kind='commit' AND s.phase='observed')
+		ORDER BY state='confirmed', created_at,id LIMIT 1`, now, now))
+}
+
+// OperatorCompletedLaunches returns bounded candidate IDs. The runner checks
+// each exact typed receipt before publishing an eligible episode.
+func (c *Credentials) OperatorCompletedLaunches(ctx context.Context) ([]string, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return nil, err
+	}
+	rows, err := c.db.QueryContext(ctx, `SELECT a.id FROM operator_actions a JOIN operator_action_steps s ON s.action_id=a.id
+		WHERE a.kind='start' AND a.state='succeeded' AND s.kind='launch' AND s.phase='observed' AND s.version=4
+		ORDER BY a.created_at DESC,a.id LIMIT 16`)
+	if err != nil {
+		return nil, err
+	}
+	defer rows.Close()
+	var ids []string
+	for rows.Next() {
+		var id string
+		if err := rows.Scan(&id); err != nil {
+			return nil, err
+		}
+		ids = append(ids, id)
+	}
+	return ids, rows.Err()
+}
+
+// OperatorFinishLaunch reads the sole completed start for an episode, with
+// the same structural/identity checks as launch-status.
+func (c *Credentials) OperatorFinishLaunch(ctx context.Context, episode string) (OperatorLaunchSnapshot, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorLaunchSnapshot{}, err
+	}
+	id := ""
+	err := c.db.QueryRowContext(ctx, `SELECT a.id FROM operator_actions a JOIN operator_action_steps s ON s.action_id=a.id
+		WHERE a.episode_id=? AND a.kind='start' AND a.state='succeeded' AND s.kind='launch' AND s.phase='observed' AND s.version=4
+		ORDER BY a.created_at DESC,a.id LIMIT 1`, episode).Scan(&id)
+	if errors.Is(err, sql.ErrNoRows) {
+		err = ErrOperatorMissing
+	}
+	if err != nil {
+		return OperatorLaunchSnapshot{}, err
+	}
+	return c.ReadOperatorLaunch(ctx, id)
+}
diff --git a/host/store/operator_finish_review.go b/host/store/operator_finish_review.go
new file mode 100644
index 0000000..24c9e30
--- /dev/null
+++ b/host/store/operator_finish_review.go
@@ -0,0 +1,72 @@
+package store
+
+import (
+	"encoding/json"
+	"fmt"
+	"strconv"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+type OperatorFinishPath struct {
+	Path, Mode, Digest string
+}
+
+type OperatorFinishTree struct {
+	Head, Branch, Index, Config string
+	Paths                       []OperatorFinishPath
+}
+
+type OperatorLocalCommit struct {
+	Snapshot           OperatorFinishTree
+	Message, Committer string
+}
+
+// OperatorFinishReview is the bounded manifest published by the runner. The
+// daemon fills only the draft decision and signed catalogue generation.
+type OperatorFinishReview struct {
+	Store, Episode, Preset, Address, Worktree string
+	Report, PresetDigest, Profile, Sessions   string
+	Launch                                    string
+	Remote                                    StartRemote
+	Base                                      StartBase
+	Checks                                    []StartCheck
+	Protected, Unattributed                   []string
+	Commit                                    OperatorLocalCommit
+}
+
+func (v OperatorFinishReview) Action(generation int64, ready bool) (approval.OperatorAction, error) {
+	manifest, err := json.Marshal(v)
+	if err != nil {
+		return approval.OperatorAction{}, err
+	}
+	intent := "draft"
+	if ready {
+		intent = "ready"
+	}
+	o := approval.OperatorOptions{Preset: v.Preset, ConfirmOnSurface: true, Ready: ready}
+	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, Kind: "finish", Episode: v.Episode, Store: v.Store,
+		Argv: approval.OperatorArgv("finish", v.Episode, o), Options: o,
+		Resolution: []approval.OperatorResolved{{Name: "finish_catalogue", Value: strconv.FormatInt(generation, 10)},
+			{Name: "finish_manifest", Value: string(manifest)}},
+		Review: fmt.Sprintf("Prepare local commit for %s (%s pull request review follows)\n%s", v.Episode, intent, manifest),
+	}
+	_, err = approval.OperatorReviewDigest(a)
+	return a, err
+}
+
+const OperatorFinishCatalogueKey = "@finish"
+
+func (h OperatorHeartbeat) FinishReview(episode, preset string) (OperatorFinishReview, error) {
+	var catalogue map[string]json.RawMessage
+	var entries map[string]OperatorFinishReview
+	if len(h.Catalogue) > 16*approval.MaxOperatorActionBytes || json.Unmarshal(h.Catalogue, &catalogue) != nil ||
+		json.Unmarshal(catalogue[OperatorFinishCatalogueKey], &entries) != nil {
+		return OperatorFinishReview{}, ErrOperatorRunner
+	}
+	v, ok := entries[episode]
+	if !ok || v.Episode != episode || v.Preset != preset {
+		return OperatorFinishReview{}, ErrOperatorRunner
+	}
+	return v, nil
+}
diff --git a/host/store/operator_steps.go b/host/store/operator_steps.go
index aaf601b..12072e9 100644
--- a/host/store/operator_steps.go
+++ b/host/store/operator_steps.go
@@ -34,9 +34,10 @@ func reserveOperatorStep(ctx context.Context, db operatorStepWriter, key Operato
 		SELECT ?, ?, id, 'reserved', ?, owner, ? FROM operator_actions
 		WHERE id = ? AND digest = ? AND generation = ? AND owner = ?
 		AND state IN ('claimed', 'running') AND expires_at > ? AND created_at <= ?
+		AND (? <> 'commit' OR NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind = 'commit'))
 		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind <> 'operation') < 32
 		ON CONFLICT(kind, identity) DO NOTHING`, key.Kind, key.Identity, key.Version, now,
-		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner, now, now, key.Action.ID)
+		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner, now, now, key.Kind, key.Action.ID, key.Action.ID)
 	return operatorChangedOne(result, err)
 }
 
diff --git a/host/worktree/worktree.go b/host/worktree/worktree.go
index 4fb871f..65841eb 100644
--- a/host/worktree/worktree.go
+++ b/host/worktree/worktree.go
@@ -95,15 +95,21 @@ func Output(ctx context.Context, git string, args []string) (string, error) {
 
 type isolatedGitEnvKey struct{}
 
-// OutputIsolated uses a trusted caller's explicit environment at the same
-// bounded Git site. Unlike ambient Git, caller-supplied policy can set GIT_*.
-func OutputIsolated(ctx context.Context, git string, args, env []string) (string, error) {
-	return gitOutput(context.WithValue(ctx, isolatedGitEnvKey{}, childenv.Scrubbed(env)), git, args)
+type isolatedGitInput struct {
+	env   []string
+	stdin io.Reader
 }
 
-func gitOutput(ctx context.Context, git string, args []string) (string, error) {
+// RunIsolated gives the bounded, scrubbed Git seam explicit stdin and an
+// optional inherited episode lock. No child reads from a terminal.
+func RunIsolated(ctx context.Context, git string, args, env []string, stdin io.Reader, inherited ...*os.File) (string, error) {
+	in := isolatedGitInput{childenv.Scrubbed(env), stdin}
+	return gitOutput(context.WithValue(ctx, isolatedGitEnvKey{}, in), git, args, inherited...)
+}
+
+func gitOutput(ctx context.Context, git string, args []string, inherited ...*os.File) (string, error) {
 	out := &capped{max: MaxStdout}
-	if err := runGit(ctx, git, args, out); err != nil {
+	if err := runGit(ctx, git, args, out, inherited...); err != nil {
 		return "", err
 	}
 	if out.over {
@@ -167,12 +173,11 @@ func runGit(ctx context.Context, git string, args []string, stdout *capped, inhe
 	runCtx, cancel := context.WithTimeout(ctx, timeout)
 	defer cancel()
 	cmd := exec.CommandContext(runCtx, git, args...)
-	if approved, ok := ctx.Value(isolatedGitEnvKey{}).([]string); ok {
-		cmd.Env = approved
+	if approved, ok := ctx.Value(isolatedGitEnvKey{}).(isolatedGitInput); ok {
+		cmd.Env, cmd.Stdin = approved.env, approved.stdin
 	} else {
 		cmd.Env = Env(os.Environ())
 	}
-	cmd.Stdin = nil
 	cmd.ExtraFiles = inherited
 	cmd.WaitDelay = waitDelay
 	if deadline, ok := runCtx.Deadline(); ok {
````

### M3 — S9b private finish, check journal and route integration (S9b draft)

````diff M3
diff --git a/cmd/ailang-worldd/operator_catalogue.go b/cmd/ailang-worldd/operator_catalogue.go
index 190d331..7ebbdb0 100644
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
@@ -35,8 +35,17 @@ func (e serviceEnv) readOperatorCatalogue(c runnerConfig, src operatorStartSourc
 			}
 		}
 	}
+	if len(c.FinishRemotes) != 0 {
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
diff --git a/cmd/ailang-worldd/operator_finish_local.go b/cmd/ailang-worldd/operator_finish_local.go
new file mode 100644
index 0000000..8b5eb80
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_local.go
@@ -0,0 +1,235 @@
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
+	if err != nil {
+		return nil, err
+	}
+	if !reflect.DeepEqual(snap, w.manifest.Commit.Snapshot) {
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
+	if _, err := w.inputs(ctx); err != nil {
+		return err
+	}
+	for i, check := range m.Checks[len(j.Checks):] {
+		if err := w.fresh(ctx); err != nil {
+			return err
+		}
+		if err := w.check(ctx, check, i+len(j.Checks)); err != nil {
+			return err
+		}
+		if _, err := w.inputs(ctx); err != nil {
+			return err
+		}
+	}
+	contents, err := w.inputs(ctx)
+	if err != nil {
+		return err
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
+	runCtx, stop := context.WithTimeout(ctx, time.Duration(check.TimeoutS*float64(time.Second)))
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
+	if err != nil {
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
+	sha, err := w.git.inspectCommit(ctx, w.manifest.Commit, tree)
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
index 0000000..fc85bc6
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_local_test.go
@@ -0,0 +1,293 @@
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
+	g := operatorFinishGitFixture(t, f)
+	ctx := boundedTestContext(t)
+	snap, _, err := g.snapshot(ctx, f.rep.Start)
+	if err != nil {
+		t.Fatal(err)
+	}
+	path, err := filepath.EvalSymlinks(f.health.CredentialsDB)
+	if err != nil {
+		t.Fatal(err)
+	}
+	db, err := store.OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer db.Close()
+	key := approvaltest.NewSigner(51)
+	if err := db.EnrolApprover(ctx, store.ApproverRow{Name: "fixture", Format: "ssh", PublicKey: key.PublicKey().Blob, EnrolledAt: 900}); err != nil {
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
+	row, err := db.RequestOperatorAction(ctx, req, 1000)
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
+	if _, err := db.ConfirmOperatorAction(ctx, cas, sig, proof, 1001); err != nil {
+		t.Fatal(err)
+	}
+	h := store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "finish-runner"}, Protocol: store.OperatorRunnerProtocol, Build: "fixture", Readiness: "both", SeenAt: 1001}
+	h.Owner.Generation, err = db.PutOperatorHeartbeat(ctx, h)
+	if err != nil {
+		t.Fatal(err)
+	}
+	owner := ""
+	if len(confirmed) == 0 {
+		if err := db.ClaimOperatorAction(ctx, cas, h.Owner, 1001); err != nil {
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
+		read: ro, now: func() int64 { return now }, budget: func(time.Duration) (context.Context, context.CancelFunc) {
+			return context.WithTimeout(t.Context(), 15*time.Second)
+		}}
+	w.fresh = func(ctx context.Context) error {
+		_, err := verifyOperatorFinishPrep(ctx, ro, a, now, cas, owner)
+		return err
+	}
+	w.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
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
+		row, err := ro.OperatorAction(boundedTestContext(t), p.row.ID, now)
+		if err != nil || row.V3ID != "" || len(f.seam.calls) != 0 || len(f.scopes) != 0 || len(f.finishRecords) != 0 {
+			t.Errorf("private preparation leaked into dispatch: V3=%q, %v, %v", row.V3ID, f.seam.calls, err)
+		}
+	})
+	return f, w, &now
+}
+
+func operatorFinishJournal(t *testing.T, w operatorFinishWork) store.OperatorFinishJournal {
+	t.Helper()
+	j, err := w.read.ReadOperatorFinishJournal(boundedTestContext(t), operatorPreparationCAS(w.prep.row))
+	if err != nil {
+		t.Fatal(err)
+	}
+	return j
+}
+
+func TestOperatorFinishPrep(t *testing.T) {
+	f, w, _ := operatorLocalFinishFixture(t, nil)
+	original := w.manifest.Commit.Message
+	w.manifest.Commit.Message = "different local commit"
+	baseline := operatorFinishJournal(t, w)
+	if err := runEpisodeFinishWithPrep(boundedTestContext(t), w); !errors.Is(err, errOperatorFinishReview) ||
+		!reflect.DeepEqual(baseline, operatorFinishJournal(t, w)) || f.checks != 0 {
+		t.Fatalf("other signed intent was admitted: %v", err)
+	}
+	w.manifest.Commit.Message = original
+	if !w.boundManifest(w.generation) {
+		t.Fatal("signed manifest must bind original inputs")
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
+			return err
+		}
+		return io.ErrUnexpectedEOF // bounded commit advanced HEAD, reply was lost
+	}
+	ctx := boundedTestContext(t)
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
+	if err := runEpisodeFinishWithPrep(boundedTestContext(t), w); err != nil || f.checks != 1 || f.outcomes != 1 {
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
+		if err := runEpisodeFinishWithPrep(boundedTestContext(t), w); !errors.Is(err, errOperatorCommitFailed) {
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
+	for _, name := range []string{"failure", "abort", "timeout", "changed-input"} {
+		t.Run(name, func(t *testing.T) {
+			f, w, _ := operatorLocalFinishFixture(t, nil)
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
+			want := map[string]error{"failure": errOperatorFinishCheck, "abort": errOperatorFinishUncertain,
+				"timeout": errOperatorFinishUncertain, "changed-input": errOperatorFinishReview}[name]
+			for range 2 {
+				if err := runEpisodeFinishWithPrep(boundedTestContext(t), w); !errors.Is(err, want) {
+					t.Fatal(err)
+				}
+			}
+			j := operatorFinishJournal(t, w)
+			if f.checks != 1 || j.CommitPhase != "" || strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD")) != f.before {
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
+		if err := runEpisodeFinishWithPrep(boundedTestContext(t), w); !errors.Is(err, errOperatorFinishTerminal) ||
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
+	if err := runEpisodeFinishWithPrep(boundedTestContext(t), w); !errors.Is(err, store.ErrOperatorExpired) {
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
+	if err := runEpisodeFinishWithPrep(boundedTestContext(t), w); !errors.Is(err, errOperatorFinishReview) ||
+		serviceFileHash(w.manifest.Store) != before || f.checks != 0 || f.outcomes != 0 {
+		t.Fatalf("untracked drift crossed preflight: %v", err)
+	}
+	f.write("b.go", "hand\n")
+	f.gitOutput("-C", f.wt, "commit", "--allow-empty", "-qm", "foreign HEAD")
+	if err := runEpisodeFinishWithPrep(boundedTestContext(t), w); !errors.Is(err, errOperatorFinishReview) ||
+		serviceFileHash(w.manifest.Store) != before || f.checks != 0 || f.outcomes != 0 {
+		t.Fatalf("HEAD drift crossed preflight: %v", err)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_finish_runner.go b/cmd/ailang-worldd/operator_finish_runner.go
new file mode 100644
index 0000000..3a5965b
--- /dev/null
+++ b/cmd/ailang-worldd/operator_finish_runner.go
@@ -0,0 +1,393 @@
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
+		sessions[0].Project != preset.Profile || sessions[0].ProfileDigest != profile.Digest {
+		return out, errOperatorFinishEligibility
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
+	if in.manifest.Remote != m.Remote || in.manifest.Base != m.Base || in.manifest.PresetDigest != m.PresetDigest ||
+		in.manifest.Profile != m.Profile || in.manifest.Sessions != m.Sessions || in.manifest.Launch != m.Launch ||
+		!reflect.DeepEqual(in.manifest.Checks, m.Checks) {
+		return errOperatorFinishReview
+	}
+	if !old {
+		reviewed, err := in.frozen(ctx)
+		if err != nil {
+			return err
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
+		if current.manifest.Profile != m.Profile || current.manifest.Sessions != m.Sessions ||
+			current.manifest.PresetDigest != m.PresetDigest || current.manifest.Launch != m.Launch ||
+			current.manifest.Remote != m.Remote || current.manifest.Base != m.Base || !current.baseMatches(bound) {
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
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
index ddd3362..ba70f4b 100644
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
@@ -142,6 +154,22 @@ func (e serviceEnv) runOperatorPreparationLoop(errOut io.Writer, c runnerConfig,
 				fmt.Fprintln(errOut, "operator start blocked; inspect service runner status and reconcile durable evidence")
 			}
 		}
+		if h.Readiness == "both" {
+			err := e.runFinishTurn(c, &h, src, stops, ticks)
+			var auth operatorGitHubAuth
+			switch {
+			case errors.Is(err, errOperatorDraining):
+				return 0
+			case errors.As(err, &auth):
+				h.Readiness = "start"
+				if err := e.touchRunner(c, &h); err != nil {
+					return serviceMessage(errOut, "runner halted", "heartbeat unavailable", "ailang-worldd service runner status", 2)
+				}
+				fmt.Fprintln(errOut, auth)
+			case err != nil && !errors.Is(err, store.ErrOperatorMissing):
+				fmt.Fprintln(errOut, "operator finish blocked; inspect action status and durable evidence")
+			}
+		}
 		select {
 		case _, open := <-ticks:
 			if !open {
@@ -160,12 +188,18 @@ func (e serviceEnv) runnerTick(c runnerConfig, h *store.OperatorHeartbeat) error
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
@@ -233,7 +267,9 @@ func (e serviceEnv) runnerStatus(w io.Writer) error {
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
index 0000000..6eee6bd
--- /dev/null
+++ b/host/daemon/surface_operator_finish_test.go
@@ -0,0 +1,134 @@
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
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorFinishAdmission(t *testing.T) {
+	f := newOperatorHTTP(t)
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
+	routine := f.do("GET", "/surface/requests.json", nil, nil)
+	if routine.Code != 200 || strings.Contains(routine.Body.String(), manifest.Report) || strings.Contains(routine.Body.String(), "artifact") {
+		t.Fatalf("routine action list exposed protected payloads: HTTP %d", routine.Code)
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
index 47ad635..a80c9a9 100644
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


# w-world-address-bar-verbs-v4 — S7 sprint plan

## 0. Scope and outcome

Row 243, iteration 445, lane `codex:gpt-6.2-sol`, base `a83d40a6f5bb7d2cfe4f3183ea1061b8f983853b`. The reviewed design is r1, quorum r2 PROCEED 3/3. D-WORLD-119 and D-WORLD-120 are A, attended at `82f05004`. S6 #695 merged at `00b4ec55`. `git log -4 --format='%h %s'` also identified the #694 fix `2b3ba7af` and the #696 record merge at the base. Design §8 budgets 540 + 330 = 870 raw and 1,175 calibrated lines.

**Outcome B.** The compiled complete candidate measures **1,892 changed lines** with rename detection. Its inert prefix measures **1,399**. M1 pins the actual launch worker and receipt preconditions. M2 ships the existing verb route's `review` outcome, inert POST, stable card, JSON status, Needs-you projection and runbook. The confirmation URL receives a 404 with the pending contents unchanged. One follow-up, **S7b**, applies the compiled M3 patch at **517 changed lines** against the M2 tree. It registers confirmation, its key and transport refusal pins, ceremony data and live process identity readback in that milestone. Browser POST activation follows that milestone's outside acceptance. The S8 browser bindings and Chrome can consume the live routes afterward.

The measurement is equivalent to `git diff -M --numstat a83d40a6 HEAD -- cmd host docs app`. The planner applied Git's rename and numstat algorithm to the two archive trees with all changed paths present. `numstat-M{1,2,3}{,-unfiltered}` holds the per-path counts and `numstat-final.raw` holds the arithmetic. Independent full-tree numstat
using `git diff --no-index -M --numstat "$S/base-tree" "$S/st-M2"` and the M3 equivalent
returned +1,344/−55 and +1,803/−89 respectively. Controller commit output is the final check. The unfiltered **product** counts match the filtered counts. The two plan documents are outside both product snapshots.

**Why the cut matters.** At base and through M2, no production caller of the credentials store's `ConfirmOperatorAction` can wake the installed runner. At M3, the HTTP judge and the atomic credentials transaction can make one pending action claimable. Before activation, verify real worker/descriptor behavior, exact proof, no state change on refusal, and one counter update on replay. An HTTP request never owns a process lifetime.

### Reviewed §5.1 rows

| Design row | Owner and acceptance |
|---|---|
| Other-verb parity | S7 M2 server via `TestOperatorAdmissionHTTP` and existing verb report/why/approve selectors. S8 reuses those responses in the bar. |
| Start review | S7 M2 GET review, then explicit inert POST, in `TestOperatorAdmissionHTTP`. S8 owns Enter binding and browser gesture. |
| Finish review | S9 owns the finish preview, submission and no-effects-until-key trace. S7 returns terminal guidance. |
| Grammar case | S1 `TestOperatorGrammarCase`, reused by M2's `AdmitOperator` call. |
| Grammar arity and missing value | S1 `TestOperatorGrammarArity`, reused by M2. |
| Grammar whole words | S1 `TestOperatorGrammarExactWords`, reused by M2. |
| Grammar quoted/equals/interleaved canonicalization | S1 `TestOperatorGrammarCanonical`, M2 review fixture uses equals and interleaving. |
| Grammar duplicate flags | S1 `TestOperatorGrammarDuplicates`, reused by M2. |
| Grammar unsupported options/false/standalone launch | S1 `TestOperatorGrammarOptions`, M2 admits headless starts only. |
| Grammar authority modes | S1 `TestOperatorGrammarPolicy`, M2 provides A/A with hand-off unavailable. S9 adds ready finish and S12 adds interactive availability. |
| Never dispatch | S1 `TestOperatorNeverDispatch` plus M2 server cases for publish, enrol/revoke, mint/revoke and service. |
| Identifier anchors | S1 grammar for episode/preset and M2 route mutation D-ID for action IDs and path boundaries. |
| Identifier flags | S1 grammar and M2 transport cases for case and escaped newline. |
| Identifier lengths | S1 grammar and M2 transport cases for short and long IDs. |
| Request limits | M2 HTTP body, task and operation grammar. M3 carries the separate 20 KiB proof-body cap with confirmation. |
| Wire types | M2 JSON object, argv array and element cases plus pure typed link/status projection. S8 owns invalid replies at the browser. |
| Strict JSON | M2 request object and M3 confirm object share the closed reader. |
| Literal task | M2 HTML escaping and redacted list/log assertions. S6 start/launch tests drive literal task→one prompt argument. S7b repeats the argument check against the browser-supplied task form. |
| Misdirected | S7 M2 Host controls and a body-read spy before admission. S7b verifies Host on confirm with valid proof. |
| Cross origin | S7 M2 request and S7b confirm, exact allowlist. |
| Fetch site | S7 M2 request and S7b confirm. |
| Media type and body bound | S7 M2 request and S7b confirm. |
| Unsupported method / OPTIONS | S7 M2 and S7b, no Access-Control grant. |
| Runner unavailable | S7 M2 terminal/direct POST refusal and S4 heartbeat arms `TestOperatorRunnerUnavailable`. |
| Catalogue admission | S7 M2 catalogue publisher, generation and missing preset refusal, plus S4 `TestOperatorCatalogueAdmission`. |
| Admission boundaries | S7 M2 seventeenth row and reused operation conflict over HTTP. S3 `TestOperatorAdmissionBoundaries` owns cross-kind episode conflict. |
| Duplicate tap | S7 M2 both operation and intent reuse in HTTP. S3 retains `TestOperatorDuplicateTap`. |
| Reconnect | S7 M2 reads operation identity across catalogue expiry. S3 `TestOperatorReconnect` verifies reopen persistence. S8 owns browser refresh/retry policy. |
| Panel network/invalid reply | S8. |
| Panel latest/abort | S8. |
| Link allowlist | S7 M2 anchored status/card/report paths. S8 owns response-link types in the client. S10 owns linked V3 card and S11 pinned HTTPS PR URL. |

### Reviewed §5.2 rows

| Design row | Owner and acceptance |
|---|---|
| Both keys | S7b M3 `TestOperatorConfirmHTTP`, using the real SSHSIG and fake authenticators. Existing law `TestOperatorConfirmBothKeys` still runs. |
| Plain y / missing or malformed proof / wrong challenge | S7b M3 `TestOperatorConfirmHTTP/Proof` and existing proof law. |
| Forged, unknown, namespace, UV, RP and origin | S7b M3 `TestOperatorConfirmHTTP/Authenticity` and existing authenticity law. |
| Revocation and counter | S7b M3 `TestOperatorConfirmHTTP/RevocationCounter` and existing counter law. |
| Other action/episode/kind/store/argv/preset/task and old artifacts | S7b M3 `TestOperatorConfirmHTTP/Binding` and existing binding law. |
| Expiry | S7b M3 `TestOperatorConfirmHTTP/Expiry`, including exact equality, before issue, after expiry and 601 s. |
| Confirm replay | S7b M3 exact artifact replay returns already-confirmed with one counter advance. Existing S3 CAS replay remains a gate. |
| Confirm crash/read/write races | S7b M3 `TestOperatorConfirmHTTP/MissingAndWriter` and the S3 `TestOperatorConfirmCrash` rollback and counter-race gates. S8 owns disconnect and retry in the browser. |
| Claim before expiry and cancel/runner loss | Existing S3 `TestOperatorClaimBeforeExpiry` and S4 `TestOperatorCancelConfirmed`. S8 owns browser disconnect during claim. |
| Claim race/owner | S3 `TestOperatorClaimRace`, S4 kernel lock tests and S5f4 live preparation. |
| Terminal exclusion | S3/S4 `TestOperatorTerminalExclusion` and installed runner/CLI activation pins. No change this sprint. |
| Runner rechecks | S2b claimed proof and S5/S6 start/effect recheck fixtures. M3 connects the existing pipeline. |
| Recover untouched | S5f4 `TestOperatorCurrentOwnerProcessRecovery` and `TestOperatorStartRecoveryWindows`. |
| Recover reserved/live/unknown child | S6 launch law, S7 M1 real worker pins, S7b live proof readback, S8 close/reconnect from the browser. |
| Status disconnected | S7 M2 `TestOperatorAdmissionHTTP/Status` and pure invalid-state fixture. S8 owns the last-known client cache on transport failure. |
| Pruning/tombstones | S3 `TestOperatorPrune` and S6 unresolved launch retention in `TestOperatorLaunchReceiptLaw`. |
| Cancel after in-flight checkpoint | S4 `TestOperatorCancelInflight` preserves effects, S6 launch reservation remains fenced. S9 and S11 own finish-phase cancellation/reconciliation. |

The base app classifies by `outcome`, as `grep -nE 'VERB_NEEDS|b.outcome ===' host/surfacepage/app/app.js`
shows at 1430–1442. A refusal cannot take the address-navigation branch. The reviewed start response
is available through the server in S7, and the panel presents it once S8 binds the new outcome.

S8 owns all `TestOperatorPanel*` browser rows, Chrome and phone viewport. S9–S11 own finish/F1. S12 owns the interactive herdr path. The reviewed §8 allocations
for S8, S9, S10, S11 and S12 are 1,161, 1,215, 1,094, 1,215 and 945 calibrated lines.
Only S7b has a compiled follow-up measurement here. Schema 10 supplies the existing tables.

## 1. Measurements

### 1.1 Base and provenance

`git rev-parse HEAD` returned `a83d40a6f5bb7d2cfe4f3183ea1061b8f983853b`; `git status --short` was empty before planning. `anchors.sh` records positive controls and base/stage line numbers in `anchors-base-tree.raw` and `anchors-st-M*.raw`. Read the reviewed design §§3.1–3.3, 4.1–4.2, 5.1–5.2, 6 and 8, the S6 sprint plan and iteration-444 readings/judges. Use the S6 judges' remeasurements as provenance for carried mutants only. The planner left iteration-444 scratch read-only.

| Command run from the base archive | Observed |
|---|---|
| `grep -nF 'GET /surface/verb.json' host/daemon/daemon.go`, then `grep -nF '/surface/operator-actions' host/daemon/*.go` | Verb at 1175, through `approvalRoute(d.handleSurfaceVerb)`. Zero operator-action path matches. |
| `grep -nF 'a.Outcome = "terminal"' host/daemon/surface_verb.go` | Line 69. S7 calls the admitted operator subset before the original navigation/classification switch. |
| `grep -nE 'func .* (RequestOperatorAction\|ConfirmOperatorAction\|OperatorAction)\(' host/store/operator_{actions,cas}.go` | Base request at 86, CAS in `operator_cas.go:42`, row read at 60. |
| `grep -nE 'func (JudgeOperatorAction\|RecheckOperatorAction\|RecheckClaimedOperatorAction)' host/authority/judge_operator.go` | 61, 77, 129. Judge comment at 59 identifies S7. |
| `grep -nE 'readRequests\|readConfirmations' host/daemon/surface_{approve,home_summary}.go` | Both paths call `readRequests` (97 and 24). V3 `readConfirmations` is separate (113). Add the redacted pending operator rows to the shared read and the summary to that list's length. |
| `grep -nE 'humanMu.Lock\|commitBudget\|OpenCredentials(ReadOnly)?\(' host/daemon/surface_confirm.go` | V3 has the human serialization and request-bound budget (129–132), read-only store then non-waiting writer (134, 192). M2 extracts the open into `openConfirmationWriter` and keeps V3 on it. M3 reuses that exact seam. |
| `grep -nE 'JudgeVerb\|ParseApproveBody' host/daemon/surface_confirm.go` and `grep -nE 'CardPasskeys\|func RequestCard' host/surfacepage/approve_model.go` | Use the same approval-body contract, Origin/fetch/media checks and passkey grouping. M3 uses V4's canonical 97-byte message and namespace, and its own judge. V3/A3 artifacts stay invalid for V4. |
| `grep -nE 'Catalogue\|Marshal' cmd/ailang-worldd/operator_runner.go` | At 165 only `json.Marshal(c.PresetNames)`. A necessary S7 change is the private runner-published snapshot of frozen launch facts. `PutOperatorHeartbeat` already stores bytes and increments catalogue generation when readiness or contents change. No browser pre-key Git lookup. |
| `grep -nE 'reconstructOperatorStart\|launchPreparedOperatorStart\|StartWorker' cmd/ailang-worldd/operator_{preparation,runner,launch}.go` | Preparation reconstructs at 92, installed loop hands off at 118, launch supervisor crosses the `StartWorker` seam at 216. Worker validates before using the admitted handoff. |
| `gofmt -l host/daemon/daemon.go` | Base prints `host/daemon/daemon.go`. M2 fixes the aligned adjacent registration comments while changing Handler. |
| `grep -nF VERIFY_GO_RACE_PKG_TIMEOUT_S=900 scripts/verify_go.sh` | Per-package bound at 161. CI timings and slow-runner bound in §4. |
| `grep -nE 'review\|OperatorLaunchWorker\|J06\|J17\|J18' "$HOME/.ailang/state/world-iter444/evaluator-r2.md"` | The real worker, cap and independent birth guards were left open there. Exact carried mutations are read from `world-iter444/judge/mut.py`. |

### 1.2 Guard and route inventory

Selectors are exact top-level names. `gates.sh` runs the grouped variants with positive PASS controls. Inventory source anchors are in `anchors-*.raw` and patch hunks in Appendix P. Scope-specific program invocations are under §7.

| Selector and result | Change or relevance |
|---|---|
| `TestDaemonDependencyAllowlist` — PASS | New daemon imports are local packages and the standard library. `$S/inv-M2/daemon.cmd` and `.raw`. |
| `TestWorldBoundaryDependencyAllowlist` — PASS | Dependency boundary remains local. `$S/inv-M2/boundary.cmd` and `.raw`. |
| `TestProductionContextRoots` — PASS | Catalogue uses operatorContext and HTTP uses readCtx or the request-bound commit budget. No new context roots. `$S/inv-M2/store.cmd` and `.raw`. |
| `TestCredentialsMethodsRequireADeadline` — PASS | New reads share operatorAccess. `$S/inv-M2/store.cmd` and `.raw`. |
| `TestStoreSurfaceTakesContext`, `TestStoreDeadlineGuardDynamic` — PASS | Lookup and list use bounded contexts. `$S/inv-M2/store.cmd` and `.raw`. |
| `TestBoundaryASTWriteGuard` — PASS | No new World writer. `$S/inv-M2/boundary.cmd` and `.raw`. |
| `TestOperatorNeverDispatch` — PASS | Original refusal law and test reused before navigation. `$S/inv-M2/verbs.cmd` and `.raw`. |
| `TestCLICredentialWriterSites` — PASS | Catalogue shares runnerWrite. `$S/inv-M2/cli.cmd` and `.raw`. |
| `TestSessionRootBudgets` — PASS | No session root changes. `$S/inv-M2/cli.cmd` and `.raw`. |
| `TestSubprocessSinkGate`, `TestExecGateArmsAreLoadBearing`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestVerifygateTestWritesAreForkLocked` — PASS | The worker fixture uses syncBuffer or io.Discard and existing exec sites. `$S/inv-M2/verifygate.cmd` and `.raw`. |
| `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` — UNINFORMATIVE UNDER SANDBOX | The existing agentrun drivers exercise both sites and the sentinel. `$S/inv-M2/broker.cmd` and `.raw`. |
| `TestDaemonCredentialWriterSites` — PASS | Writer census replaces handleSurfaceConfirm's direct open with the shared openConfirmationWriter. `$S/inv-M2/daemon.cmd` and `.raw`. |
| `TestSurfaceVerbIsReadOnly` — PASS | The method inventory adds the read-only operatorVerb adapter. `$S/inv-M2/daemon.cmd` and `.raw`. |
| `TestFrozenV1RouteTableMatchesMux`, `TestApprovalRoutesRegistered`, `TestApprovalRoutesRefuseForeignHost`, `TestConfirmRoutesRegistered`, `TestDecisionRoutesRegistered`, `TestSurfaceVerbRouteRegistered` — PASS | The guarded action family intercepts only its own paths before ServeMux cleaning. `$S/inv-M2/daemon.cmd` and `.raw`. |
| `TestSurfaceSecurityHeaders`, `TestSurfaceHomeSummary`, `TestRequestsJSONConfirmationsGolden` — PASS | The new card uses inherited CSP. Needs you adds the shared operator list and retains the V3 golden. `$S/inv-M2/daemon.cmd` and `.raw`. |
| `TestReadRoutesAndPayloadGate`, `TestGETRoutesRejectOtherMethods`, `TestProjectionRoutes_MountedAndNotProtected`, `TestAgentCardRoute_GETOnly`, `TestSurfaceAppRoutes`, `TestSurfaceScriptRoute`, `TestSurfaceModelRoutes`, `TestApprovePageRoutes`, `TestPasskeyPageRoutes`, `TestSurfaceFontRoute`, `TestWorkbenchRouteIsReadOnly`, `TestWorkbenchSecurityHeaders` — PASS | Recheck surrounding method, route and header families on the wrapped mux. `$S/inv-M2/daemon-extra.cmd` and `.raw`. |

The scoped agent-run subprocess drivers passed under the existing broker guard in M1, M2 and M3, with both driver0 and driver1 observed. The full broker selector hits a loopback denial here. The S7b dependency and writer inventories also passed with its shared process observer in `$S/m3-guard-inventory.raw`. See `$S/M*-broker-agentrun.raw` and `$S/inv-M2/broker.raw`.

### 1.3 Known-red baseline and race budgets

Use `$STATE/base` from the outside controller after the last line of `$STATE/baseline.done` matches `ALL_DONE for $B`. A package with a listener denial is UNINFORMATIVE UNDER SANDBOX even if it reports individual FAIL lines. Compare full TEXT fail sets, and retain the raw log, rc, parser and timeout observation. Baseline observations are in §4.

### 1.4 Production-caller census

`/bin/bash "$S/callers.sh" "$S/base-tree" base` and the same command with `st-M1`, `st-M2`, `st-M3`. Positive controls are the installed runner loop, launch handoff and existing `agentrun.StartWorker` call.

The base row reader census has **seven** call expressions, six in the CLI and one in host/store/operator_worktree.go:76. This corrects the inherited count. The authority judge's call to the approval law is a positive control distinct from the store CAS. Existing Host and route recorder tests in host/daemon/surface_approve_rig_test.go:108–120 also correct the premise that every daemon route test binds. The new S7 HTTP gates do bind.

````text
CALLERS base-tree RequestOperatorAction 0
CALLERS base-tree ConfirmOperatorAction 0
CALLERS base-tree JudgeOperatorAction 0
CALLERS st-M1 RequestOperatorAction 0
CALLERS st-M1 ConfirmOperatorAction 0
CALLERS st-M1 JudgeOperatorAction 0
host/daemon/surface_operator.go:230:	row, err := wr.RequestOperatorAction(ctx, store.OperatorRequest{Action: action, OperationID: body.OperationID, Challenge: challenge,
CALLERS st-M2 RequestOperatorAction 1
CALLERS st-M2 ConfirmOperatorAction 0
CALLERS st-M2 JudgeOperatorAction 0
host/daemon/surface_operator.go:235:	row, err := wr.RequestOperatorAction(ctx, store.OperatorRequest{Action: action, OperationID: body.OperationID, Challenge: challenge,
CALLERS st-M3 RequestOperatorAction 1
host/daemon/surface_operator_confirm.go:131:	row, err = wr.ConfirmOperatorAction(ctx, store.OperatorCAS{ID: id, Digest: digest, Generation: row.Generation}, artifact, proof, now)
CALLERS st-M3 ConfirmOperatorAction 1
host/daemon/surface_operator_confirm.go:107:	proof, err := authority.JudgeOperatorAction(ctx, ro, authority.OperatorJudgeInput{
CALLERS st-M3 JudgeOperatorAction 1
````

## 2. Milestones

| Milestone | Added | Deleted | Cumulative product | Increment | Unfiltered product |
|---|---:|---:|---:|---:|---:|
| M1 | +158 | −15 | 173 | 173 | 173 |
| M2 | +1,344 | −55 | 1,399 | 1,226 | 1,399 |
| M3 (S7b) | +1,803 | −89 | 1,892 | 517 | 1,892 |

**M1 — pre-reachability pins.** Its real worker entry receives existing production stdin, descriptors 3/4/5 and installed build identity from `main`. A test calls that entry with duplicated real temp lock descriptions, signed durable fixture, fake service transport and real fake child. An unrelated inode and a busy separately-opened description yield the usage refusal and identical durable evidence. Successful worker execution has one child marker and completed receipt. The worker is bounded by the fixture's injected budget. The supervisor test separately fires an injected timer channel on an owned fake worker. The receipt law rejects empty starting birth directly, and a unit at the production Started guard rejects an unreadable birth ahead of store readback.

**M2 — inert route family and projection.** Attach one family handler ahead of ServeMux canonical redirects. The four selectors through M2 are the existing verb GET, request POST, card GET and status GET. In the action family, Host precedes path/method and POST-body checks. On new requests validate the strict body, rerun the shared grammar, capture the server's current catalogue under the existing non-waiting writer, then call `RequestOperatorAction`. For lost responses, use the operation id ledger and compare the same canonical argv, exact task and original generation before returning the stored action. The new row is requested; an operation ledger entry is recorded alongside it. Keep the stable card under inherited CSP. Needs-you is the same credentials read for requests and for the home summary.

**S7b M3 — confirmation activation.** Apply the compiled follow-up delta in Appendix P only in the next sprint. It registers the method and negative tests together. Its card presents only ceremony metadata for a pending, live review. Confirm parses exactly one artifact under the same guards, verifies the immutable row and envelope through `JudgeOperatorAction`, then calls the existing store CAS with a fresh clock and the exact artifact. Counter/approver use and state are one transaction. Replays return a fresh read of the approved row for the same artifact. Process readback gains the same kernel birth check as the terminal. S7b should run the existing S6 runner fake-child gates alongside the compiled HTTP fixture before activation.

## 3. Durable windows and operator output

| Surface or window | Observable result |
|---|---|
| `GET /surface/verb.json` | `review` adds `catalogueGeneration`, canonical `argv` and the frozen preview `text`. An unsupported/unavailable action receives the current terminal/refused presentation and terminal instructions. |
| `POST /surface/operator-actions` | A requested action at `href` with `statusHref`. Same operation and intent recover the existing href; conflicting reuse is 409. The seventeenth pending request receives `QueueFull`. The evidence is one action, its operation identity, zero external effects, no runner work and zero new sessions. |
| `GET /surface/operator-actions/{id}` | Stable HTML card with escaped review and literal task. The route family refuses escaped, shortened, suffixed or redirected selectors. S7b adds the exact pending ceremony with `post`, `digest`, `messageHex`, `namespace` and origin-grouped `passkeys`. |
| `GET /surface/operator-actions/{id}/status.json` | `v,id,episode,kind,state,href,statusHref,reportHref,reason,now,createdAt,expiresAt,terminalAt,runnerSeenAt`. Missing/busy/read failure is explicit. A stale runner leaves the last durable state plus a reason and last heartbeat timestamp. A reserved/unproven process is uncertain, and an observed terminal receipt returns its outcome while keeping the same URL. S7b uses the shared live PID/birth observation for an unsettled launched receipt. |
| Body-bearing route refusals | Host 421 runs before the body read. Absent, null, duplicate or foreign Origin and invalid fetch metadata return 403. Unsupported media returns 415 and an oversized body returns 413. Wrong methods return 405 on registered actions with no Access-Control header. |
| `GET /surface/requests.json` | `operatorActions` holds the pending id, episode, kind, expiry and exact link. An unavailable projection reports `operatorActionsUnavailable`. A compatible empty projection uses the existing empty-list behavior. `GET /surface/home-summary.json` adds the number from that same read to Needs you. |
| S7b `POST /surface/operator-actions/{id}/confirm` | `result` is confirmed or already-confirmed and `approvedBy` identifies the enrolled key. One exact confirmed start becomes claimable. Wrong proof, expiry, revocation or busy credentials leaves the action and key use unchanged. |
| Logging | M2 and M3 emit no new routine confirmation or body log line. The shared Needs-you logger retains the existing read-failure behavior, and logs from the runner retain `operator launch <id>: <status>` and its existing attended recovery advice. |

### Attended launch uncertainty

The M2 QUICKSTART names the read-only `service runner launch-status ACTION-ID`, the absence of an unlock command, and the evidence to preserve. Refresh or a fresh browser operation does not settle existing uncertainty. Use the episode report and private runner evidence for attended reconciliation. Do not delete the reservation, worktree or session to test whether the prior launch happened.

**Measured choice.** The uncertainty-specific documentation occupies 5 lines in the selected product diff. The existing attended start analogue is 253 lines across CLI, CAS and acceptance test, before command registration/help and any process-quiescence evidence. Its actual ranges are recorded below. That is a scale proxy, since a safe settle decision needs positive process evidence and has no existing launch-specific implementation. The sprint chooses documentation and retains the fence.

- `cmd/ailang-worldd/operator_start_recovery.go:195–310` has 116 lines.
- `cmd/ailang-worldd/operator_activation_test.go:499–585` has 87 lines.
- `host/store/operator_start_recovery.go:148–197` has 50 lines.

## 4. Acceptance, timing and fail sets

`accept.sh` uses fixture keys, isolated stores and plans, and reads stat or SHA observations on the real-path guards. An absent path is a valid observation. The child is a fixture shell script, and only the owning test starts it. The callable fake example is documented in QUICKSTART. In S7b its proof bytes come from the operator message function and actual signer fixtures. Browser/Chrome start-through-launch remains S8.

- Nominal CI at `00b4ec55` reported CLI 291.0 s, daemon 358.6 s, verifygate 323.0 s and broker 275.5 s. Reported slow-runner daemon at `6447e776` was 452.3 s. Keep the 900 s race bound and measure added CLI and daemon selectors separately.
- For non-attributable packages, use `outside.sh` for real listeners and `measure-race.sh` for the exact added test times. `gates.sh` owns whole-package plain/race, full base comparisons and raw logs. `-race` caps are ≤15 s for CLI and ≤12 s for daemon. A recorder rehearsal is useful fixture evidence and has no authority to clear the listener gate.
- The Linux-style temp fixture run uses a 1777 outer directory and a private 0700 child. Physical comparisons resolve aliases where paths come from two sources. Only the explicit temp leg sets TMPDIR.
- Baseline known reds include row 271 real-binary CLI and daemon e2e tests, load-sensitive Git and broker fixtures and service startup gates under sandbox. Preserve those receipts and never modify their fixtures to green this sprint.

| Added named test, measured by Go -race | Result |
|---|---|
| M1 CLI TestOperatorLaunchWorker / TestOperatorStartedRequiresBirth / TestOperatorSupervisorCap | 5.63 s / 0.00 s / 0.01 s, total **5.64 s**. `m1-final-race.raw` passed. |
| M2 CLI TestOperatorLaunchWorker, including /catalogue / TestOperatorStartedRequiresBirth / TestOperatorSupervisorCap | 7.21 s (catalogue 0.41 s inside that total) / 0.00 s / 0.02 s. **7.23 s** total, package command 8.951 s. `timing-M2-shared/ailang-worldd.raw`. Earlier separate fixtures reached 37.29 s; the shipped suite reuses one prepared fixture. |
| M2 daemon TestOperatorLinksProjection | **0.00 s** measured in `timing-M2-shared/daemon.raw`. |
| M2 daemon TestOperatorAdmissionHTTP including /Status | **UNINFORMATIVE UNDER SANDBOX**. The real test's listener is denied. See the outside-race command in §7. |
| S7b daemon TestOperatorConfirmHTTP and inherited TestOperatorAdmissionHTTP | **UNINFORMATIVE UNDER SANDBOX** for real HTTP. The last real-mux recorder rehearsal in `m3-rehearse-release.raw` measured 4.49 s confirm, 5.74 s admission and 0.00 s links. These are 10.23 s of recorder evidence, with actual listener timing still required. |
| S7b CLI TestOperatorLiteralTaskArgv | 0.00 s in `m3-literal-argv.raw`. |

Explicit 1777-parent runs produced `ACCEPT_OK M1 fast` and `ACCEPT_OK M2 fast` in `m1-final-1777.raw` and `gate-M2/t1777.raw`. The M2 stage's named non-listener fast gate passed. The base gofmt check reported the pre-existing daemon.go indentation. M2's gofmt leg was OK. M2 and M3 native build, Go vet and GOOS=linux vet receipts are in `sealed-legs.txt` and `sealed-M*-*.raw`.
| Whole-package race | Base | M1 | Final touched stage | Delta from base |
|---|---:|---:|---:|---:|
| `host/store` | 53.357 s | 203.102 s | 111.216 s | +57.859 s |
| `host/approval` | 9.912 s | untouched | 2.377 s | -7.535 s |
| `host/surfacepage` | 8.632 s | untouched | 8.106 s | -0.526 s |
| `host/agentrun (S7b)` | 20.392 s | untouched | 21.549 s | +1.157 s |
| cmd/ailang-worldd | Sandbox whole command UNINFORMATIVE | same | same | PENDING — controller to measure outside |
| host/daemon | Sandbox whole command UNINFORMATIVE | untouched | same | PENDING — controller to measure outside |

Package elapsed deltas reflect independent shared-rig runs. Compare failure membership, and use the added named-test selectors to enforce the runtime caps.

Pristine archive gate receipts in `$S/gate-base/legs.txt`. The gate started before stage verification. The base `verify_ail.sh` leg passed in 217 s. Selected non-listener inventories passed. Broker and full CLI/daemon packages hit sandbox denials. The controller's independent base legs in `$STATE/base` are read-only.
- M1 `go test -count=1 ./host/verifygate/` — no new fail-set member, recorded set host/verifygate TestMissionConfigDefaultFlowIntegration. Raw `$S/M1-verifygate-plain.raw`.
- M1 `go test -count=1 ./host/broker/` — UNINFORMATIVE UNDER SANDBOX. Raw `$S/M1-broker-plain.raw`.
- M2 `go test -count=1 ./host/verifygate/` — no new fail-set member, recorded set host/verifygate TestMissionConfigDefaultFlowIntegration. Raw `$S/M2-verifygate-plain.raw`.
- M2 `go test -count=1 ./host/broker/` — UNINFORMATIVE UNDER SANDBOX. Raw `$S/M2-broker-plain.raw`.

The controller's baseline.done ends in ALL_DONE at a83d40a6. Exact sets are in `$STATE/base/*.fails`. Outside counts by leg: store-plain/race 0/0, authority-plain/race 0/0, worktree-plain 0, verbs-plain 0, verifygate-plain 1, broker-plain 4, cmd-plain 17, daemon-race 16. The latter two package samples were 1,478.816 s and 1,628.703 s. They are fail-set controls and exceed the execution race bound under that load.

A separate M2 store-race attempt recorded TestConvertedStoreMethodsM5b/cutoff/PutVerifyResult with a context deadline before its intended cutoff. The same exact selector passed 3/3 at base and M2, and the normal M2 whole store-race leg was clean. Inspect cutoff-base-tree.raw and cutoff-st-M2.raw. M2 authority-plain also recorded TestOperatorMintBinding/action_id with a deadline under load. The exact base and M2 repetitions passed 3/3 in mint-binding-*.raw. Keep both observations in the controller fail-set comparison.

The first pristine gate started against the archive before stage verification. Recorded base legs remain in gate-base. An in-flight script revision interrupted the final recording loop, so the missing ALL_DONE on that planner gate is an instrument failure. Use the sealed outside-pristine command in §7 for any base leg absent from the completed controller set.

The repeated complete M2 authority-plain leg passed in 8.794 s and its empty fail set matched the pristine controller set. The earlier timing failure remains recorded in gate-M2/authority-plain.new. The repeat evidence is m2-authority-plain-repeat.raw.

S7b's joint broker/verifygate leg recorded a broker listener denial. Compare verifygate's text receipt only where the package completed, and repeat the combined gate outside before registering confirm. See `$S/M3-guards-plain.raw`.

Existing server read/address/why behavior was exercised by `TestSurfaceVerbBoundsAndErrors`, `TestSurfaceVerbWhyIsReadOnly` and `TestWhyResult`, all PASS in `m2-verb-parity.raw`.

## 5. Existing tests and helpers modified

| Existing test/helper | Driver before and after | Named mutation for behavioral witness |
|---|---|---|
| `host/store/TestOperatorLaunchReceiptLaw` | Same migrated fixture and `RecordOperatorLaunch`, with direct wrong-birth readback and file/row evidence on rejection | J06 kills the new arm before the success control. |
| `host/daemon/TestDaemonCredentialWriterSites` | Same writer-open AST census. V3 `handleSurfaceConfirm` now reaches the existing writer through `openConfirmationWriter`. M2 inert admission and M3 CAS share it. | D-REG-REQUEST and K-NO-CAS drive the new caller behavior. Existing `TestConfirmPOSTGuards` and `TestConfirmPOST` preserve the V3 path. The AST test is an inventory. |
| `host/daemon/TestSurfaceVerbIsReadOnly` | Same standard-library/verbs import inventory and ordinary verb calls. Method list adds the read-only `operatorVerb`. | D-REVIEW and D-FORBIDDEN drive the same production entry at its new branch. |
| `cmd/ailang-worldd/TestTaskTitle` moved to `host/store` | `taskTitle` becomes `StartTaskTitle` with its body unchanged. Existing CLI and new catalogue call the same function. | C-TITLE-BEHAVIOUR removes heading normalization and the same fixture rejects it after the move. C-TITLE-REUSE pins the catalogue's use of it. |
| `operatorStartFacts` support through `resolvePresetStart` | The exact EpisodeStart constructor moves to `resolvedPreset.startRecord`. CLI assigns the same StartedAt and the snapshot omits a clock until actual execution. | C-RESOLUTION changes one output field and its comparison to the live-resolution fixture fails. |
| `operatorLaunchWorker` test from M1, extended in M2 | Adds the `/catalogue` subtest after the same positive worker execution. Reuses the fixture and compares the catalogue's exact mint and launch resolution byte inputs. | W-SAMEFILE, W-FLOCK, W-STDIN and W-FDS still drive the worker. C-PUB and C-RESOLUTION drive periodic publication. |
| `operatorLaunchProcessStatus` at S7b only | Terminal caller keeps the observed/unsettled classification through the moved store method and kernel process read. The browser projection uses the same operation. | Existing `TestOperatorLaunchCrashWindow` checks matching birth, absent PID and mismatch. |

The inert fixture has an executable fake child under its temp root. D-INERT-SPAWN's recorder
rehearsal fires the child marker assertion. The baseline runner and episode effect counts remain
zero through a positive inert POST. The full launch and supervisor fixture runs in M1. No existing
subprocess launcher is added. M1 invokes the broker-driven `agentrun.StartWorker` and `agentrun.RunObserved` sites with owned fixtures. Existing `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` continues to exercise both and its positive registry-credential sentinel. The subprocess output sink is `syncBuffer` in the fake runner and `io.Discard` in the capped supervisor fixture.

## 6. Gates and mutant receipts

**Lock FD observation.** Base `operator_launch.go:216` supplies the runner, episode and lifecycle open descriptions to `StartWorker`. Base line 434 passes only `job.locks[:2]` to the supervised agent. M1's positive worker fake child checks open descriptors 3 and 4 without writing to them, then writes one marker to its temp directory. The worker's independently opened lock refusal establishes that this is a real kernel lock, and the wrong-inode refusal establishes `SameFile`. This preserves the reviewed design's lock inheritance.

**Carried variants.** J06 copies the original empty-Birth comparison deletion exactly. J17 changes the cap addition from one minute to 1,000 hours at the injected timer call, which is the same timing mutation after introducing the factory. J18 deletes the same empty-Birth guard from the extracted Started helper. Each has a named positive selector; the guard overlap no longer masks those independent checks.

**Carried ledger.** KILLED 3. `$S/ledger-carried-named.jsonl`.

- KILLED `J06`, `J17`, `J18`.

**Diff ledger.** EQUIVALENT AT THE ROUTE 2, KILLED 15, UNINFORMATIVE UNDER SANDBOX 42. `$S/ledger-diff-named.jsonl`.

- KILLED `W-SAMEFILE`, `W-FLOCK`, `W-STDIN`, `W-FDS`, `W-PID`, `W-TIMER-STOP`, `S-LINK-ACTION`, `S-LINK-EPISODE`, `S-LINK-KIND`, `S-LINK-STATE`, `C-PUB`, `C-RESOLUTION`, `C-TITLE-REUSE`, `C-TITLE-BEHAVIOUR`, `C-GRAMMAR-NEVER`.
- Recorder rehearsal failed at the named assertion for `D-HOST-ORDER`, `D-FAMILY-CLEAN`, `D-HOST`, `D-ORIGIN-EXACT`, `D-FETCH-EXACT`, `D-MEDIA-TYPE`, `D-BODY-CAP`, `D-ID`, `D-GENERATION`, `D-OP-REPLAY-EARLY`, `S-NEEDS-SHARED`, `S-CARD-ESCAPE`, `S-PENDING-STATE`, `S-PENDING-CREATED`, `D-INERT-SPAWN`. Those remain UNINFORMATIVE at the listener boundary.
- Equivalent `JSON-RAWNUL`, because the JSON scanner already rejects a literal NUL, independently of the early raw-byte check.
- Equivalent `JSON-TAIL`, because the typed unmarshal independently rejects trailing JSON and the HTTP error class is identical.

**Follow-up ledger.** EQUIVALENT AT THE ROUTE 1, KILLED 1, UNINFORMATIVE UNDER SANDBOX 15. `$S/ledger-follow-up-named.jsonl`.

- KILLED `L-ARGV-SPLIT`.
- Recorder rehearsal failed at the named assertion for `K-REG`, `K-BODY-EXACT`, `K-ONLY-HEADLESS`, `K-APPROVED-ARTIFACT`, `K-JUDGE-CHALLENGE`, `K-NO-CAS`, `K-EXPIRE-CLASS`, `L-ARGV-SPLIT`. Those remain UNINFORMATIVE at the listener boundary.
- Equivalent `K-ONLY-START`, because DecodeOperatorAction already refuses headless finish and the separate Headless guard remains.

The original C-TITLE-REUSE attempt discovered zero tests during the test-fixture move. The sealed selector rerun killed it. That first receipt stays in `mutants-named-local`; use `mutants-named-final` for the exact final names.

Whole-package local kills are J06, C-TITLE-BEHAVIOUR and C-GRAMMAR-NEVER. The short outside replay manifest is `$S/whole-mutant-ids.txt`. The named and whole commands in §7 use the same anchored patch and Go compiler. The classifier checks expected test discovery and a firing parser control, and excludes socket denial and truncation from kills. Full whole-package results need the same base fail sets as the stage gates.

The ledger is keyed by the production file, source line and exact old/new bytes, with the killing selector recorded on every row. Artefact layout, raw output and exact invocations are in each ledger directory. Treat unmeasured branch coverage and outside mutable route observations as pending before activation.

Generated conditions named locally killed 9. `Q-M1-operator_launch-469`, `Q-M2-operator_catalogue-21`, `Q-M2-operator_catalogue-25`, `Q-M2-operator_catalogue-32`, `Q-M2-operator_catalogue-33`, `Q-M2-operator_catalogue-39`, `Q-M2-operator_catalogue-50`, `Q-M2-operator_catalogue-66`, `Q-M2-operator_catalogue-71`.
Reproduce with `python3 "$S/mutate.py" named "$S/local-conditions-ids.txt" "$S/local-conditions-check-1" "$S/condition-proposals.json"`.

`BRANCH_AUDIT 189 anchors 39 with exact ledger anchors 80 generated conditional proposals; 14 need a fixture selector`

The inventory distinguishes edited lines from the smaller executable mutation ledger. Remaining conditions and order or SQL anchors are in `$S/remaining-anchors.tsv`. Candidate single-condition edits are in `$S/condition-proposals.json`. They require selector and compiler checks before they count.

PENDING — controller to measure: `python3 "$S/branch-audit.py"` and replay each selected conditional with `python3 "$S/mutate.py" named "$S/proposed-ids.txt" "$S/conditional-outside-1" "$S/condition-proposals.json"`. Entries without a fixture need one before activation. The reported mutation receipts cover the exact enumerated experiments.

### Reproduction scripts

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
case "$M" in base|M1|M2|M3) ;; *) exit 2;; esac
test ! -e "$L/legs.txt" || exit 2
B=a83d40a6f5bb7d2cfe4f3183ea1061b8f983853b
PIN=$HOME
export WORLD_PLAN_PIN_HOME="$PIN"
export AILANG_BIN="$PIN/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$PIN/.pinned-ailang-tools/v0.54.15/ailang"
export GOCACHE="$(go env GOCACHE)" GOMODCACHE="$(go env GOMODCACHE)"
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
export XDG_CONFIG_HOME="$L/home/.config" XDG_CACHE_HOME="$L/home/.cache" XDG_DATA_HOME="$L/home/.local/share"
mkdir -p "$L/home" "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
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
pkgs_changed=(./cmd/ailang-worldd/ ./host/store/)
if [ "$M" != M1 ]; then pkgs_changed+=(./host/approval/ ./host/daemon/ ./host/surfacepage/); fi
if [ "$M" = M3 ]; then pkgs_changed+=(./host/agentrun/); fi
fmt_touched() {
  local p bad
  while IFS= read -r p; do
    case "$p" in *.go) if [ -f "$p" ]; then bad=$(gofmt -l "$p") || return 2; test -z "$bad" || { printf '%s\n' "$bad"; return 1; }; fi;; esac
  done < "$S/changed-${M/base/M2}"
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
leg vet bound 300 go vet "${pkgs_changed[@]}"
leg linux-vet env GOOS=linux python3 "$S/bounded.py" 300 go vet "${pkgs_changed[@]}"
leg context-roots bound 180 go test ./host/store/ -run '^TestProductionContextRoots$'
leg inv-store inventory host/store 'TestProductionContextRoots|TestStoreSurfaceTakesContext|TestCredentialsMethodsRequireADeadline'
leg inv-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
leg inv-cli inventory cmd/ailang-worldd 'TestCLICredentialWriterSites|TestSessionRootBudgets'
leg inv-verbs inventory host/verbs TestOperatorNeverDispatch
leg inv-daemon inventory host/daemon 'TestDaemonDependencyAllowlist|TestDaemonCredentialWriterSites|TestFrozenV1RouteTableMatchesMux|TestSurfaceSecurityHeaders|TestApprovalRoutesRegistered|TestApprovalRoutesRefuseForeignHost|TestConfirmRoutesRegistered|TestDecisionRoutesRegistered|TestSurfaceVerbRouteRegistered|TestSurfaceVerbIsReadOnly|TestRequestsJSONConfirmationsGolden|TestSurfaceHomeSummary'
leg inv-broker inventory host/broker TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential
leg inv-verifygate inventory host/verifygate 'TestSubprocessSinkGate|TestExecGateArmsAreLoadBearing|TestNoParallelWriteForkPackagesOutsideVerifygate|TestVerifygateTestWritesAreForkLocked'
if [ "$(uname -s)" = Linux ]; then
  leg inv-exec-kernel inventory host/verifygate TestKernelRefusesExecOfWriterOpenFile
else
  skip inv-exec-kernel 'Linux ETXTBSY control'
fi
if [ "$M" = base ]; then
  skip accept 'S7 tests absent'
  skip t1777 'S7 tests absent'
else
  leg accept /bin/bash "$S/accept.sh" "$T" "$M"
  mkdir -m 1777 "$L/t1777"
  mkdir -m 0700 "$L/t1777/x"
  leg t1777 env TMPDIR="$L/t1777/x" /bin/bash "$S/accept.sh" "$T" "$M"
fi
leg verify_ail bound 510 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
leg callers /bin/bash "$S/callers.sh" "$T" "$M"
tags=(store-plain store-race authority-plain authority-race worktree-plain verifygate-plain broker-plain verbs-plain cmd-plain cmd-race daemon-plain daemon-race approval-plain approval-race surfacepage-plain surfacepage-race agentrun-plain agentrun-race)
pkgs=(host/store host/store host/authority host/authority host/worktree host/verifygate host/broker host/verbs cmd/ailang-worldd cmd/ailang-worldd host/daemon host/daemon host/approval host/approval host/surfacepage host/surfacepage host/agentrun host/agentrun)
for i in "${!tags[@]}"; do
  tag=${tags[$i]}
  if [ "$M" = M1 ]; then
    case "$tag" in daemon-*|approval-*|surfacepage-*) skip "$tag" 'package untouched in M1'; continue;; esac
  fi
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
printf 'ALL_DONE\n' 
[ "$issue" -eq 0 ] || exit 1
[ "$unknown" -eq 0 ] || exit 77
````

````bash accept.sh
#!/bin/bash
# accept.sh TREE M1|M2|M3 [fast|race|outside|outside-race]
set -euo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
M=$2
mode=${3:-fast}
case "$M" in M1|M2|M3) ;; *) exit 2;; esac
case "$mode" in fast|outside) race=;; race|outside-race) race=-race;; *) exit 2;; esac
PIN=${WORLD_PLAN_PIN_HOME:-$HOME}
export AILANG_BIN="$PIN/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$PIN/.pinned-ailang-tools/v0.54.15/ailang"
unset GOFLAGS
D=$(mktemp -d "$S/accept.XXXXXX")
observe() {
  local installed
  installed=$(python3 - "$PIN/.ailang-world/service.json" <<'CODE'
import json, os, sys
try:
    with open(sys.argv[1]) as f:
        data = json.load(f)
    path = data.get("credentials", "") if isinstance(data, dict) else ""
    if isinstance(path, str) and os.path.isabs(path):
        print(path)
except (OSError, ValueError):
    pass
CODE
  )
  local default_creds="$PIN/.ailang/world/credentials.db"
  for p in "$PIN/Library/LaunchAgents/com.sunholo.ailang-world.worldd.plist" \
    "$PIN/Library/LaunchAgents/com.sunholo.ailang-world.operator-runner.plist" \
    "$PIN/.ailang-world" "$PIN/.ailang-world/service.json" "$PIN/.ailang/world-ws" \
    "$default_creds" "${installed:-$default_creds}"; do
    if [ ! -e "$p" ] && [ ! -L "$p" ]; then printf 'absent %s\n' "$p"; continue; fi
    if [ "$(uname -s)" = Darwin ]; then stat -f '%N %p %i %z %m %c' "$p"
    else stat -c '%n %a %i %s %Y %Z' "$p"; fi
    case "$p" in "$default_creds"|"$installed") ;; *)
      if [ -f "$p" ] && [ -r "$p" ]; then shasum -a 256 "$p"; fi ;;
    esac
  done
}
observe > "$D/before"
trap 'observe > "$D/after"; diff -u "$D/before" "$D/after"' EXIT
cd "$T"
run() {
  local pkg=$1 names=$2 log="$D/${1##*/}.raw" rc=0
  python3 "$S/bounded.py" 390 go test ${race:+-race} -count=1 -v -timeout 350s "./$pkg/" -run "^($names)\$" > "$log" 2>&1 || rc=$?
  cat "$log"
  if grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$log"; then
    echo 'UNINFORMATIVE UNDER SANDBOX'
    return 77
  fi
  if grep -Eq 'panic: test timed out|BOUND EXCEEDED' "$log"; then return 2; fi
  if [ "$rc" -ne 0 ]; then return "$rc"; fi
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log" || return 2
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
CLI='TestOperatorLaunchWorker|TestOperatorSupervisorCap|TestOperatorStartedRequiresBirth'
STORE='TestOperatorLaunchReceiptLaw'
if [ "$M" = M3 ]; then CLI="$CLI|TestOperatorLiteralTaskArgv"; fi
if [ "$M" != M1 ]; then
  STORE="$STORE|TestOperatorAdmissionBoundaries|TestOperatorDuplicateTap|TestOperatorRunnerUnavailable|TestTaskTitle"
fi
run cmd/ailang-worldd "$CLI"
run host/store "$STORE"
if [ "$M" != M1 ]; then
  run host/verbs TestOperatorNeverDispatch
  DAEMON='TestOperatorLinksProjection|TestSurfaceVerbIsReadOnly|TestDaemonCredentialWriterSites'
  case "$mode" in
    outside|outside-race)
      DAEMON="$DAEMON|TestOperatorAdmissionHTTP"
      [ "$M" != M3 ] || DAEMON="$DAEMON|TestOperatorConfirmHTTP"
      ;;
    *) echo "LEG listener-acceptance SKIPPED (run outside.sh for actual loopback)" ;;
  esac
  run host/daemon "$DAEMON"
fi
if [ "$M" = M3 ]; then
  run host/approval 'TestOperatorConfirmBothKeys|TestOperatorConfirmProof|TestOperatorConfirmAuthenticity|TestOperatorConfirmRevocationCounter|TestOperatorConfirmBinding|TestOperatorConfirmExpiry'
  run host/authority 'TestOperatorConfirmReplay|TestOperatorConfirmCrash|TestOperatorConfirmStoreBinding'
fi
observe > "$D/after"
diff -u "$D/before" "$D/after"
trap - EXIT
printf 'ACCEPT_OK %s %s\n' "$M" "$mode"
````

## 7. Apply, accept and snapshot

```bash
B=a83d40a6f5bb7d2cfe4f3183ea1061b8f983853b
STATE=$HOME/.ailang/state/world-iter445
S=$STATE/plan-scratch
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter445-row243-v4s7
PLAN=$W/design_docs/planned/w-world-address-bar-verbs-v4s7-sprint-plan.md
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
for M in M1 M2; do
  awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$S/executor-$M.patch"
  test -s "$S/executor-$M.patch" || exit 1
  (cd "$W" && /usr/bin/patch -p1 -E < "$S/executor-$M.patch") || exit 1
  /bin/bash "$S/accept.sh" "$W" "$M" fast || exit 1
  /bin/bash "$S/snapshot.sh" "$W" "$M" || exit 1
done
```

`accept.sh` succeeds with **exit 0 and `ACCEPT_OK <milestone> <mode>`**. In sandbox modes the listener leg prints SKIPPED. Use outside modes for the listener gates. Snapshots are cumulative from `$B`. `$STATE/snap/<M>/` contains product files at repository-relative paths and `STATUS`, `DELETED`, `MANIFEST.sha256`. Plan files are excluded. The operator's real paths are observed only.

### Controller commands

The commit script refuses untracked planning docs, so record the two authored plan files before installing the commit snapshots. Its Go build, vet and positive named checks run outside the sandbox.

```bash
/bin/bash "$S/outside.sh" "$S/st-M2" "$S/outside-M2-1" M2
/bin/bash "$S/measure-race.sh" "$S/st-M2" M2 "$S/outside-race-M2-1" outside
/bin/bash "$S/gates.sh" "$S/base-tree" "$S/outside-pristine-1" base
/bin/bash "$S/gates.sh" "$S/st-M1" "$S/outside-gate-M1-1" M1 "$S/outside-pristine-1"
/bin/bash "$S/gates.sh" "$S/st-M2" "$S/outside-gate-M2-1" M2 "$S/outside-pristine-1"
/bin/bash "$S/whole-mutants.sh" "$S/outside-named-1" "$S/mutant-diff-ids.txt" named
/bin/bash "$S/whole-mutants.sh" "$S/outside-named-carried-1" "$S/mutant-carried-ids.txt" named
/bin/bash "$S/whole-mutants.sh" "$S/outside-whole-1"
"$STATE/commit_ms.sh" M1 "$S/M1.msg" "export AILANG_BIN='$AILANG_BIN' WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE='$WORLD_EXEC_NODE' WORLD_TOOL_AILANG_BIN='$WORLD_TOOL_AILANG_BIN'; /bin/bash '$S/accept.sh' '$W' M1 fast"
"$STATE/commit_ms.sh" M2 "$S/M2.msg" "export AILANG_BIN='$AILANG_BIN' WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE='$WORLD_EXEC_NODE' WORLD_TOOL_AILANG_BIN='$WORLD_TOOL_AILANG_BIN'; /bin/bash '$S/accept.sh' '$W' M2 outside && /bin/bash '$S/commit-guards.sh' '$W' '$S/commit-M2-guards-1'"
"$STATE/ctl_roundtrip.sh" M1 M2
```

The M2 commit prerequisite includes `go test -count=1 ./host/verifygate/ ./host/broker/`, as executed by `commit-guards.sh`. That wrapper checks the completed base receipts, the parser control and absence of new failures. CLI and daemon commit checks are named only. Some packages in planner gates require a listener or process environment and need outside fail-set runs. Use `gates.sh "$S/st-M3" "$S/outside-gate-M3-1" M3 "$S/outside-pristine-1"`, `outside.sh "$S/st-M3" "$S/outside-M3-1" M3` and `whole-mutants.sh "$S/outside-named-M3-1" "$S/mutant-follow-up-ids.txt" named` for S7b before registering confirm.

````bash snapshot.sh
#!/bin/bash
# snapshot.sh TREE M1|M2; SNAP_STATUS_LIST/ SNAP_PARENT are rehearsal-only.
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1|M2|M3|M4) ;; *) exit 2;; esac
B=a83d40a6f5bb7d2cfe4f3183ea1061b8f983853b
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter445/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s7-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s7.json"
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

### Rehearsal

The sealed export passed a native build and vet after applying the exact three fences.

````text
REHEARSAL /Users/voightkampff/.ailang/state/world-iter445/plan-scratch/rt-4
M1 diff -r cmd rc=0
M1 diff -r host rc=0
M1 diff -r docs rc=0
M1 diff -r app SKIPPED (absent on both sides)
M1 diff -r full archive rc=0
M1 SNAPSHOT_OK 5 paths
M2 diff -r cmd rc=0
M2 diff -r host rc=0
M2 diff -r docs rc=0
M2 diff -r app SKIPPED (absent on both sides)
M2 diff -r full archive rc=0
M2 SNAPSHOT_OK 29 paths
M3 diff -r cmd rc=0
M3 diff -r host rc=0
M3 diff -r docs rc=0
M3 diff -r app SKIPPED (absent on both sides)
M3 diff -r full archive rc=0
M3 SNAPSHOT_OK 38 paths
ROUNDTRIP_OK
````

## Appendix P — patches

M1 and M2 are the S7 product milestones. M3 is the measured S7b follow-up and is outside the S7 apply loop.

### M1

````diff M1
diff --git a/cmd/ailang-worldd/main.go b/cmd/ailang-worldd/main.go
index 75dd7701..aca0cb7a 100644
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -393,7 +393,7 @@ func run(args []string, stdout, stderr io.Writer) int {
 		if addrGiven {
 			return exitUsage
 		}
-		return realServiceEnv().runOperatorLaunchWorker(rest[1:], stderr)
+		return realServiceEnv().runOperatorLaunchWorker(rest[1:], stderr, operatorWorkerInput{daemon.CurrentBuild().ExecutableSHA256, os.Stdin, os.NewFile})
 	}
 	switch verb := rest[0]; verb {
 	case "version":
diff --git a/cmd/ailang-worldd/operator_launch.go b/cmd/ailang-worldd/operator_launch.go
index aceffdf0..e3c17073 100644
--- a/cmd/ailang-worldd/operator_launch.go
+++ b/cmd/ailang-worldd/operator_launch.go
@@ -223,8 +223,15 @@ func (e serviceEnv) superviseOperatorLaunch(job operatorLaunchJob, stops <-chan
 	go func() { done <- cmd.Wait() }()
 	ticker := time.NewTicker(runnerPoll)
 	defer ticker.Stop()
-	cap := time.NewTimer(job.wall + time.Minute)
-	defer cap.Stop()
+	timer := e.operatorSupervisorCap
+	if timer == nil {
+		timer = func(d time.Duration) (<-chan time.Time, func()) {
+			t := time.NewTimer(d)
+			return t.C, func() { t.Stop() }
+		}
+	}
+	cap, stopCap := timer(job.wall + time.Minute)
+	defer stopCap()
 	for {
 		select {
 		case err := <-done:
@@ -234,7 +241,7 @@ func (e serviceEnv) superviseOperatorLaunch(job operatorLaunchJob, stops <-chan
 		case <-stops:
 			_ = cmd.Process.Signal(syscall.SIGTERM)
 			stops = nil
-		case <-cap.C:
+		case <-cap:
 			_ = cmd.Process.Signal(syscall.SIGTERM)
 			return errors.New("launch supervisor did not report completion")
 		}
@@ -262,7 +269,13 @@ func (e serviceEnv) writeOperatorLaunch(c runnerConfig, r store.OperatorLaunchRe
 	})
 }
 
-func (e serviceEnv) runOperatorLaunchWorker(args []string, out io.Writer) int {
+type operatorWorkerInput struct {
+	build   string
+	handoff io.Reader
+	file    func(uintptr, string) *os.File
+}
+
+func (e serviceEnv) runOperatorLaunchWorker(args []string, out io.Writer, input operatorWorkerInput) int {
 	proc, stop := serveSignalContext()
 	defer stop()
 	if len(args) != 1 || !verbs.ValidOperatorActionID(args[0]) {
@@ -270,7 +283,7 @@ func (e serviceEnv) runOperatorLaunchWorker(args []string, out io.Writer) int {
 	}
 	c, b, err := e.installedRunner()
 	if err != nil || os.Getenv("WORLD_OPERATOR_CONFIG_SHA256") != serviceHash(b) ||
-		daemon.CurrentBuild().ExecutableSHA256 != c.Executables[e.paths.Binary] {
+		input.build != c.Executables[e.paths.Binary] {
 		return exitUsage
 	}
 	ctx, cancel := e.operatorContext(operatorLaunchBudget)
@@ -288,7 +301,7 @@ func (e serviceEnv) runOperatorLaunchWorker(args []string, out io.Writer) int {
 	names := [3]string{c.Credentials + ".operator-runner.lock",
 		store.OperatorEpisodeLockPath(c.Credentials, row.Episode), filepath.Dir(e.paths.Record)}
 	for i, name := range names {
-		fd := os.NewFile(uintptr(3+i), name)
+		fd := input.file(uintptr(3+i), name)
 		defer fd.Close()
 		st, err := fd.Stat()
 		path, pathErr := os.Lstat(name)
@@ -298,7 +311,7 @@ func (e serviceEnv) runOperatorLaunchWorker(args []string, out io.Writer) int {
 		syscall.CloseOnExec(int(fd.Fd()))
 		job.locks[i] = fd
 	}
-	if _, err := io.ReadFull(os.Stdin, job.secret[:]); err != nil {
+	if _, err := io.ReadFull(input.handoff, job.secret[:]); err != nil {
 		return exitUsage
 	}
 	if err := e.executeOperatorLaunch(proc, job, out); err != nil {
@@ -453,10 +466,8 @@ func (e serviceEnv) executeOperatorLaunch(proc context.Context, job operatorLaun
 		},
 		Started: func(pid int) error {
 			defer func() { writer.Close(); writer = nil; job.locks[2].Close() }()
-			r.PID = pid
-			r.Birth, _ = operatorProcessIdentity(pid)
-			if r.Birth == "" {
-				return errors.New("process birth unavailable")
+			if err := requireOperatorProcessBirth(&r, pid, operatorProcessIdentity); err != nil {
+				return err
 			}
 			s, err := writer.ReadOperatorLaunch(ctx, row.ID)
 			if err != nil {
@@ -500,6 +511,17 @@ func (e serviceEnv) executeOperatorLaunch(proc context.Context, job operatorLaun
 	return nil
 }
 
+// requireOperatorProcessBirth stops the Started hook before receipt readback
+// when the kernel cannot identify the child it just started.
+func requireOperatorProcessBirth(r *store.OperatorLaunchReceipt, pid int, birth func(int) (string, bool)) error {
+	r.PID = pid
+	r.Birth, _ = birth(pid)
+	if r.Birth == "" {
+		return errors.New("process birth unavailable")
+	}
+	return nil
+}
+
 func operatorLaunchHistory(in operatorLaunchInputs, r store.OperatorLaunchReceipt, ep string, stdout, stderr *os.File, started, ended time.Time) store.EpisodeLaunch {
 	transcript := filepath.Join(in.sessions, ep+".transcript.jsonl")
 	errPath := filepath.Join(in.sessions, ep+".agent.err")
diff --git a/cmd/ailang-worldd/operator_worker_test.go b/cmd/ailang-worldd/operator_worker_test.go
new file mode 100644
index 00000000..7ea92079
--- /dev/null
+++ b/cmd/ailang-worldd/operator_worker_test.go
@@ -0,0 +1,117 @@
+package main
+
+import (
+	"bytes"
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
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func TestOperatorLaunchWorker(t *testing.T) {
+	f := newOperatorLaunchRig(t, "#!/bin/sh\n: >&3 || exit 3\n# The locks are open read/write. Inspect their presence without writing.\n: >&4 || exit 4\n"+strings.TrimPrefix(fakeLaunchMarker, "#!/bin/sh\n"))
+	e := &f.env
+	config, err := os.ReadFile(e.runnerPaths().Record)
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Setenv("WORLD_OPERATOR_CONFIG_SHA256", serviceHash(config))
+	e.operatorSupervisor = func(job operatorLaunchJob, log io.Writer) error {
+		input := operatorWorkerInput{build: f.h.Build}
+		run := func(other *os.File, handoff []byte, out io.Writer) int {
+			input.handoff = bytes.NewReader(handoff)
+			input.file = func(n uintptr, name string) *os.File {
+				fd := job.locks[n-3]
+				if n == 3 && other != nil {
+					fd = other
+				}
+				dup, err := syscall.Dup(int(fd.Fd()))
+				if err != nil {
+					t.Fatal(err)
+				}
+				return os.NewFile(uintptr(dup), name)
+			}
+			return e.runOperatorLaunchWorker([]string{job.row.ID}, out, input)
+		}
+		before := activationEvidence(t, f.worktreeStartRig)
+		foreign, err := os.CreateTemp(t.TempDir(), "unrelated")
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer foreign.Close()
+		occupied, err := os.OpenFile(f.creds+".operator-runner.lock", os.O_RDWR, 0)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer occupied.Close()
+		for _, wrong := range []*os.File{foreign, occupied} {
+			if rc := run(wrong, job.secret[:], io.Discard); rc != exitUsage ||
+				!reflect.DeepEqual(before, activationEvidence(t, f.worktreeStartRig)) {
+				t.Fatalf("bad inherited descriptor crossed reservation, rc=%d", rc)
+			}
+		}
+		rc := run(nil, job.secret[:], log)
+
+		if rc != 0 {
+			t.Fatalf("worker exited %d", rc)
+		}
+		return nil
+	}
+	f.tick(t)
+	marker, err := os.ReadFile(filepath.Join(f.bin, "ran"))
+	s := f.receipt(t)
+	if err != nil || string(marker) != "x" || s.Receipt.Status != "stopped" || s.Receipt.Birth == "" {
+		t.Fatalf("worker child or fds 3/4: marker=%q receipt=%+v err=%v", marker, s.Receipt, err)
+	}
+}
+
+func TestOperatorStartedRequiresBirth(t *testing.T) {
+	var r store.OperatorLaunchReceipt
+	if err := requireOperatorProcessBirth(&r, 42, func(int) (string, bool) { return "", false }); err == nil || r.Birth != "" {
+		t.Fatalf("Started can advance without a kernel birth: %+v %v", r, err)
+	}
+	if err := requireOperatorProcessBirth(&r, 42, func(pid int) (string, bool) {
+		if pid != 42 {
+			t.Errorf("wrong PID: %d", pid)
+		}
+		return "kernel-start", true
+	}); err != nil || r.PID != 42 || r.Birth != "kernel-start" {
+		t.Fatalf("Started did not retain the birth: %+v %v", r, err)
+	}
+}
+
+func TestOperatorSupervisorCap(t *testing.T) {
+	dir := t.TempDir()
+	f, err := os.CreateTemp(dir, "lock")
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer f.Close()
+	path := filepath.Join(dir, "worker")
+	if err := os.WriteFile(path, []byte("#!/bin/sh\nwhile :; do :; done\n"), 0700); err != nil {
+		t.Fatal(err)
+	}
+	e := serviceEnv{}
+	e.paths.Binary = path
+	fired, stopped := make(chan time.Time, 1), false
+	fired <- time.Time{}
+	wall := 2 * time.Minute
+	e.operatorSupervisorCap = func(got time.Duration) (<-chan time.Time, func()) {
+		if got != wall+time.Minute {
+			t.Errorf("supervisor cap duration %s", got)
+		}
+		return fired, func() { stopped = true }
+	}
+	job := operatorLaunchJob{wall: wall, locks: [3]*os.File{f, f, f}}
+	err = e.superviseOperatorLaunch(job, nil, io.Discard, store.OperatorHeartbeat{}, "")
+	if err == nil || !strings.Contains(err.Error(), "did not report completion") || !stopped || errors.Is(err, os.ErrClosed) {
+		t.Fatalf("supervisor did not honor its injected cap: %v", err)
+	}
+}
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
index f2012daf..50a414c0 100644
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -88,9 +88,10 @@ type serviceEnv struct {
 	runnerSleep     func(time.Duration)
 	runnerProcesses func() ([]worldProcess, error)
 
-	operatorDeadline   func(time.Duration, context.Context)
-	operatorSupervisor func(operatorLaunchJob, io.Writer) error
-	operatorAgentRun   operatorAgentRunner
+	operatorDeadline      func(time.Duration, context.Context)
+	operatorSupervisor    func(operatorLaunchJob, io.Writer) error
+	operatorAgentRun      operatorAgentRunner
+	operatorSupervisorCap func(time.Duration) (<-chan time.Time, func())
 }
 
 func realServiceEnv() serviceEnv {
diff --git a/host/store/operator_launch_test.go b/host/store/operator_launch_test.go
index 7af9a654..8ae99135 100644
--- a/host/store/operator_launch_test.go
+++ b/host/store/operator_launch_test.go
@@ -86,6 +86,9 @@ func TestOperatorLaunchReceiptLaw(t *testing.T) {
 	}
 
 	r := OperatorLaunchReceipt{Run: id, Executable: digest, PID: 42, Birth: "kernel-start", Status: "launched"}
+	wrongBirthAtStart := r
+	wrongBirthAtStart.Birth = ""
+	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, reserved, wrongBirthAtStart, nil, pinNow), ErrOperatorInput)
 	wrongRun := r
 	wrongRun.Run += "x"
 	f.assertRefused(t, before, f.c.RecordOperatorLaunch(f.ctx, reserved, wrongRun, nil, pinNow), ErrOperatorConflict)
````

### M2

````diff M2
diff --git a/cmd/ailang-worldd/episode_preset.go b/cmd/ailang-worldd/episode_preset.go
index 860bbf88..d3de397a 100644
--- a/cmd/ailang-worldd/episode_preset.go
+++ b/cmd/ailang-worldd/episode_preset.go
@@ -219,7 +219,7 @@ func resolvePresetStart(ctx context.Context, name, dir, profilesDir string, inpu
 		if err != nil {
 			return nil, err
 		}
-		taskRecord = &store.StartTask{Path: path, SHA256: serviceHash(taskData), Title: taskTitle(taskData)}
+		taskRecord = &store.StartTask{Path: path, SHA256: serviceHash(taskData), Title: store.StartTaskTitle(taskData)}
 	}
 	req, err := prepareMint(mintOptions{verb: "episode start", episode: ep, db: p.Daemon.DB, credsDB: p.Daemon.CredentialsDB,
 		wsRoot: p.Daemon.WorkspaceRoot, repo: p.Preset.Repo, branch: branch, grants: p.Grants, ttl: ttl,
@@ -246,10 +246,8 @@ func resolvePresetStart(ctx context.Context, name, dir, profilesDir string, inpu
 		}
 	}
 	req.startPoint, req.existingBranch = p.BaseSHA, existing
-	req.start = &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: ep, StartedAt: time.Unix(now(), 0).Format(time.RFC3339),
-		Preset: store.StartPreset{Name: name, Path: p.Path, SHA256: p.SHA256, Bytes: p.File}, Repo: p.Preset.Repo,
-		Remote: store.StartRemote{Name: p.Preset.Remote, URL: p.RemoteURL}, Base: store.StartBase{Branch: p.Preset.Base, SHA: p.BaseSHA},
-		Branch: branch, BranchExisted: existing, Agent: agent, Model: model, Overrides: overrides, Task: taskRecord, Checks: p.Checks, TTLS: ttl}
+	req.start = p.startRecord(ep, branch, agent, model, ttl, existing, overrides, taskRecord)
+	req.start.StartedAt = time.Unix(now(), 0).Format(time.RFC3339)
 	var prompt strings.Builder
 	fmt.Fprintf(&prompt, "episode %s from preset %s (%s, %s)\n  repo      %s, remote %s = %s\n  base      %s @ %s (%s)\n",
 		ep, name, p.SHA256, p.Path, p.Preset.Repo, p.Preset.Remote, p.RemoteURL, p.Preset.Base, p.BaseSHA, p.Lag)
@@ -286,6 +284,14 @@ func resolvePresetStart(ctx context.Context, name, dir, profilesDir string, inpu
 		task: taskData, launchMode: mode, launchWall: wall, launchBinary: binary}, nil
 }
 
+// startRecord is the exact mint input shared by the live read and the runner catalogue.
+func (p *resolvedPreset) startRecord(ep, branch, agent, model string, ttl int64, existing bool, overrides []string, task *store.StartTask) *store.EpisodeStart {
+	return &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: ep,
+		Preset: store.StartPreset{Name: p.Name, Path: p.Path, SHA256: p.SHA256, Bytes: p.File}, Repo: p.Preset.Repo,
+		Remote: store.StartRemote{Name: p.Preset.Remote, URL: p.RemoteURL}, Base: store.StartBase{Branch: p.Preset.Base, SHA: p.BaseSHA},
+		Branch: branch, BranchExisted: existing, Agent: agent, Model: model, Overrides: overrides, Task: task, Checks: p.Checks, TTLS: ttl}
+}
+
 func readPresetStartTask(task string) (string, []byte, error) {
 	path, err := filepath.Abs(task)
 	if err != nil {
diff --git a/cmd/ailang-worldd/operator_catalogue.go b/cmd/ailang-worldd/operator_catalogue.go
new file mode 100644
index 00000000..190d3316
--- /dev/null
+++ b/cmd/ailang-worldd/operator_catalogue.go
@@ -0,0 +1,83 @@
+package main
+
+import (
+	"context"
+	"encoding/json"
+	"net/url"
+	"path/filepath"
+	"strings"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// Read the configured presets on runner ticks, independent of HTTP input.
+func (e serviceEnv) readOperatorCatalogue(c runnerConfig, src operatorStartSource) []byte {
+	items := make(map[string]store.OperatorStartReview)
+	r, err := readServiceRecord(e.paths)
+	if err != nil {
+		return []byte("{}")
+	}
+	h, err := e.health(c.Address)
+	if err != nil || h.Status != "ok" || h.Build.ExecutableSHA256 != c.Executables[e.paths.Binary] {
+		return []byte("{}")
+	}
+	ctx, cancel := e.operatorContext(3 * time.Second)
+	defer cancel()
+	for _, name := range c.PresetNames {
+		p, err := loadLaunchPreset(ctx, name, src.presets, src.profiles, r.Workspace)
+		if err == nil {
+			if entry, ok := operatorCatalogueEntry(ctx, c, r.Workspace, p, h); ok {
+				items[name] = entry
+			}
+		}
+	}
+	raw, _ := json.Marshal(items)
+	if ctx.Err() != nil || len(raw) > 16*approval.MaxOperatorActionBytes {
+		return []byte("{}")
+	}
+	return raw
+}
+
+func operatorCatalogueEntry(ctx context.Context, c runnerConfig, root string, p *resolvedPreset, h daemon.HealthResponse) (store.OperatorStartReview, bool) {
+	world, worldErr := store.SameDatabase(c.World, p.Daemon.DB)
+	creds, credErr := credentialsFor(p.Daemon.DB, p.Daemon.CredentialsDB)
+	workspace, workspaceErr := filepath.EvalSymlinks(p.Daemon.WorkspaceRoot)
+	physical, physicalErr := filepath.EvalSymlinks(root)
+	if ctx.Err() != nil || !world || worldErr != nil || credErr != nil || creds != c.Credentials || p.Daemon.Addr != c.Address ||
+		workspaceErr != nil || physicalErr != nil || physical != workspace || p.Preset.Agent != "codex" || p.ProfileDigest == "" {
+		return store.OperatorStartReview{}, false
+	}
+	// Git remote configuration can contain credentials. Never publish these
+	// forms to a browser, even when accepted at the terminal.
+	if strings.Contains(p.RemoteURL, "://") {
+		u, err := url.Parse(p.RemoteURL)
+		if err != nil || u.RawQuery != "" || u.Fragment != "" || u.User != nil && (u.Scheme != "ssh" || u.User.String() != "git") {
+			return store.OperatorStartReview{}, false
+		}
+	}
+	_, wall, err := launchMode(p.Preset, true, false)
+	binary, binErr := launchBinary(p.Preset.Agent)
+	sessions, pathErr := operatorPhysicalTarget(filepath.Join(p.Daemon.SessionsDir, "snapshot.session"))
+	sum := c.Executables[binary]
+	if err != nil || binErr != nil || pathErr != nil || sum == "" || serviceFileHash(binary) != sum {
+		return store.OperatorStartReview{}, false
+	}
+	profile := approval.ProfileRef{Project: p.Preset.Profile, Digest: p.ProfileDigest}
+	req := &mintRequest{credPath: creds, prof: &broker.ExecProfile{Project: profile.Project, Digest: profile.Digest}}
+	if checkPresetDaemon(h, req, p, c.Address) != nil {
+		return store.OperatorStartReview{}, false
+	}
+	out := store.OperatorStartReview{Start: *p.startRecord("", "", p.Preset.Agent, p.Preset.Model, p.Preset.TTLS, false, []string{}, nil),
+		Profile: profile, Store: c.Credentials, Workspace: p.Daemon.WorkspaceRoot, Sessions: filepath.Dir(sessions),
+		Address: c.Address, Interpreter: h.ToolRef, Binary: binary, BinaryDigest: sum, Wall: wall.String(), Committer: displayCommitter(p.Preset.Committer)}
+	build, _ := json.Marshal(h.Build)
+	out.Build = string(build)
+	for _, g := range p.Grants {
+		out.Grants = append(out.Grants, approval.Grant{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget})
+	}
+	return out, true
+}
diff --git a/cmd/ailang-worldd/operator_catalogue_test.go b/cmd/ailang-worldd/operator_catalogue_test.go
new file mode 100644
index 00000000..dbd3010b
--- /dev/null
+++ b/cmd/ailang-worldd/operator_catalogue_test.go
@@ -0,0 +1,35 @@
+package main
+
+import (
+	"reflect"
+	"testing"
+)
+
+func testOperatorCataloguePublication(t *testing.T, f *operatorLaunchRig) {
+	t.Helper()
+	r, c, h := f.worktreeStartRig, f.c, f.h
+	r.env.operatorCatalogue = func(c runnerConfig) []byte { return r.env.readOperatorCatalogue(c, r.src) }
+	if err := r.env.runnerTick(c, &h); err != nil {
+		t.Fatal(err)
+	}
+	seen, err := runnerHeartbeat(c.Credentials, activationTestBudget)
+	if err != nil || seen.CatalogueGeneration != 2 || !seen.Available(c.Executables[r.env.paths.Binary], "start", r.src.clock()) {
+		t.Fatalf("published review generation: %d %v", seen.CatalogueGeneration, err)
+	}
+	preview, err := seen.StartReview(r.action.Options.Preset)
+	if err != nil {
+		t.Fatalf("missing signed start inputs: %v; catalogue %q", err, seen.Catalogue)
+	}
+	a, err := preview.Action(r.action.Episode, r.action.Task, seen.CatalogueGeneration)
+	if err != nil {
+		t.Fatal(err)
+	}
+	for i := range r.action.Resolution {
+		if r.action.Resolution[i].Name == "mint_catalogue" {
+			r.action.Resolution[i].Value = "2"
+		}
+	}
+	if !reflect.DeepEqual(a.Resolution, r.action.Resolution) || !reflect.DeepEqual(a.Argv, r.action.Argv) {
+		t.Fatalf("published review does not round-trip to exact launch authority:\n%+v\n%+v", a.Resolution, r.action.Resolution)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_runner.go b/cmd/ailang-worldd/operator_runner.go
index efacfff4..ddd33622 100644
--- a/cmd/ailang-worldd/operator_runner.go
+++ b/cmd/ailang-worldd/operator_runner.go
@@ -92,6 +92,7 @@ func runOperatorRunner(args []string, errOut io.Writer, e serviceEnv) int {
 	ticker := time.NewTicker(runnerPoll)
 	defer ticker.Stop()
 	source := operatorStartSource{defaultPresetsDir(), defaultExecProfilesDir(), "", defaultHostsDir(), e.terminal.now}
+	e.operatorCatalogue = func(c runnerConfig) []byte { return e.readOperatorCatalogue(c, source) }
 	return e.runOperatorPreparationLoop(errOut, c, h, source, stops, ticker.C, lock)
 }
 
@@ -159,10 +160,13 @@ func (e serviceEnv) runnerTick(c runnerConfig, h *store.OperatorHeartbeat) error
 	if err := e.checkRunnerPair(c); err != nil {
 		return err
 	}
-	h.SeenAt = e.terminal.now()
 	h.Readiness = "start"
-	// Publish the admitted installed preset names before the start consumer.
-	h.Catalogue, _ = json.Marshal(c.PresetNames)
+	if e.operatorCatalogue != nil {
+		h.Catalogue = e.operatorCatalogue(c)
+	} else {
+		h.Catalogue, _ = json.Marshal(c.PresetNames)
+	}
+	h.SeenAt = e.terminal.now()
 	return runnerWrite(c.Credentials, func(db *store.Credentials) error {
 		ctx, cancel := e.operatorContext(time.Second)
 		defer cancel()
diff --git a/cmd/ailang-worldd/operator_worker_test.go b/cmd/ailang-worldd/operator_worker_test.go
index 7ea92079..3d574c0b 100644
--- a/cmd/ailang-worldd/operator_worker_test.go
+++ b/cmd/ailang-worldd/operator_worker_test.go
@@ -70,6 +70,7 @@ func TestOperatorLaunchWorker(t *testing.T) {
 	if err != nil || string(marker) != "x" || s.Receipt.Status != "stopped" || s.Receipt.Birth == "" {
 		t.Fatalf("worker child or fds 3/4: marker=%q receipt=%+v err=%v", marker, s.Receipt, err)
 	}
+	t.Run("catalogue", func(t *testing.T) { testOperatorCataloguePublication(t, f) })
 }
 
 func TestOperatorStartedRequiresBirth(t *testing.T) {
diff --git a/cmd/ailang-worldd/service.go b/cmd/ailang-worldd/service.go
index 50a414c0..5acd7b03 100644
--- a/cmd/ailang-worldd/service.go
+++ b/cmd/ailang-worldd/service.go
@@ -92,6 +92,7 @@ type serviceEnv struct {
 	operatorSupervisor    func(operatorLaunchJob, io.Writer) error
 	operatorAgentRun      operatorAgentRunner
 	operatorSupervisorCap func(time.Duration) (<-chan time.Time, func())
+	operatorCatalogue     func(runnerConfig) []byte
 }
 
 func realServiceEnv() serviceEnv {
diff --git a/cmd/ailang-worldd/task_title.go b/cmd/ailang-worldd/task_title.go
deleted file mode 100644
index fb5f3e87..00000000
--- a/cmd/ailang-worldd/task_title.go
+++ /dev/null
@@ -1,36 +0,0 @@
-package main
-
-import "strings"
-
-// taskTitle derives a short title for new starts; stored titles stay unchanged.
-func taskTitle(data []byte) string {
-	title := ""
-	for _, line := range strings.Split(string(data), "\n") {
-		title = strings.TrimSpace(line)
-		if title != "" {
-			break
-		}
-	}
-	if strings.HasPrefix(title, "#") {
-		title = strings.TrimSpace(strings.TrimLeft(title, "#"))
-	} else {
-		if len(title) >= 5 && strings.EqualFold(title[:5], "Task:") {
-			title = strings.TrimSpace(title[5:])
-		}
-		for _, sep := range []string{". ", "; ", " ("} {
-			if i := strings.Index(title, sep); i >= 0 {
-				title = title[:i]
-			}
-		}
-	}
-	title = strings.Join(strings.Fields(title), " ")
-	runes := []rune(title)
-	if len(runes) > 72 {
-		cut := 72
-		if i := strings.LastIndex(string(runes[:72]), " "); i >= 0 {
-			cut = len([]rune(string(runes[:72])[:i]))
-		}
-		title = string(runes[:cut]) + "…"
-	}
-	return title
-}
diff --git a/cmd/ailang-worldd/task_title_test.go b/cmd/ailang-worldd/task_title_test.go
deleted file mode 100644
index 76de6a76..00000000
--- a/cmd/ailang-worldd/task_title_test.go
+++ /dev/null
@@ -1,27 +0,0 @@
-package main
-
-import (
-	"strings"
-	"testing"
-)
-
-func TestTaskTitle(t *testing.T) {
-	for _, tc := range []struct{ name, input, want string }{
-		{"455", "Task: make `world-publish`'s \"store is locked\" refusal say what to do when the World service holds the store; explain the recovery.", "make `world-publish`'s \"store is locked\" refusal say what to do when…"},
-		{"heading", "# Heading", "Heading"},
-		{"blanks", "\n \n##  A   heading\ttext\nbody", "A heading text"},
-		{"sentence", strings.Repeat("word ", 40) + ". Next.", strings.TrimSpace(strings.Repeat("word ", 14)) + "…"},
-		{"word", strings.Repeat("界", 200), strings.Repeat("界", 72) + "…"},
-		{"empty", "", ""},
-		{"whitespace", " \n\t", ""},
-		{"label", "Task:", ""},
-		{"sentence-end", "tAsK: Fix this. Then that", "Fix this"},
-		{"paren", "Fix this (details)", "Fix this"},
-	} {
-		t.Run(tc.name, func(t *testing.T) {
-			if got := taskTitle([]byte(tc.input)); got != tc.want {
-				t.Fatalf("%q, want %q", got, tc.want)
-			}
-		})
-	}
-}
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 7c99ca31..21e7b126 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -2097,3 +2097,27 @@ use the action ID and report `launched`, `stopped`, `failed`, `capped` or
 the process and receipts. `ailang-worldd episode report <episode>` shows the exit
 and transcript after completion. Saved streams scrub session credentials. Use the attended
 runner stop or paired upgrade to drain a live launch.
+
+### Review one headless start on the local surface
+
+Check the paired headless runner with `ailang-worldd service runner status`.
+1. URL-encode `episode start demo --preset sol --launch --headless` as `q`
+   in `GET /surface/verb.json`. The read-only `review` returns `argv`, preview
+   `text` and `catalogueGeneration`. An unavailable runner gives terminal guidance.
+2. Same-origin `POST /surface/operator-actions` with `Content-Type: application/json`.
+   Send that `argv` and `catalogueGeneration`, `task` (1–8192 bytes of literal UTF-8)
+   and `operationId` (`openssl rand -hex 16`). Keep the operation id for retry.
+3. Open the response's `href` at `GET /surface/operator-actions/{id}`. Poll its
+   `statusHref` at `GET /surface/operator-actions/{id}/status.json` for the durable
+   state. `GET /surface/requests.json` lists live requests in `operatorActions`.
+   The home summary counts them from that read. Run this fake example with loopback:
+
+```sh
+AILANG_BIN="$HOME/.pinned-ailang/ailang" go test -count=1 ./host/daemon/ -run '^TestOperatorAdmissionHTTP$'
+```
+
+For a missing definitive launch receipt, run `ailang-worldd service runner
+launch-status ACTION-ID`. There is no unlock command for `uncertain`. Preserve
+the action id, runner log, launch receipt, worktree and session records for
+attended reconciliation. Check the episode report for existing work. Only
+positive evidence can settle the reservation. Refresh and retry preserve it.
diff --git a/host/approval/operator_action.go b/host/approval/operator_action.go
index 148bc6d5..8a9c3996 100644
--- a/host/approval/operator_action.go
+++ b/host/approval/operator_action.go
@@ -234,6 +234,30 @@ func DecodeOperatorAction(b []byte) (OperatorAction, error) {
 	return a, nil
 }
 
+// DecodeOperatorObject reads a closed, single JSON object without duplicate
+// keys or case-folded aliases. The caller supplies its body and field bounds.
+func DecodeOperatorObject(b []byte, dst any, keys ...string) error {
+	if !utf8.Valid(b) {
+		return ErrOperatorUTF8
+	}
+	if bytes.IndexByte(b, 0) >= 0 {
+		return ErrOperatorNUL
+	}
+	if err := strictObject(b); err != nil {
+		return err
+	}
+	var fields map[string]json.RawMessage
+	if err := json.Unmarshal(b, &fields); err != nil {
+		return err
+	}
+	for name := range fields {
+		if !slices.Contains(keys, name) {
+			return ErrOperatorUnknownKey
+		}
+	}
+	return json.Unmarshal(b, dst)
+}
+
 // strictObject walks b's tokens (host/daemon strictRosterJSON's form): one JSON object (D4), no key
 // twice in any object (D5) and nothing after it (D6). encoding/json alone keeps the last of a repeated
 // key and stops reading after the first value.
diff --git a/host/daemon/daemon.go b/host/daemon/daemon.go
index 9f97d8ce..2ebec163 100644
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -287,7 +287,7 @@ type Config struct {
 	// probe is left out — named on the error log and in /v1/health — while
 	// every other profile serves; the --exec-profile flags keep their
 	// all-or-nothing rule.
-	ExecProfilesDir     string
+	ExecProfilesDir string
 	// PresetsDir (row 269) is a read-only directory of <name>.json launch
 	// presets surfaced on GET /surface/presets.json; World never writes here.
 	PresetsDir          string
@@ -1172,7 +1172,7 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("POST /surface/confirm/{id}", d.approvalRoute(d.handleSurfaceConfirm))
 	mux.HandleFunc("GET /surface/confirm.js", d.approvalRoute(d.handleSurfaceConfirmScript))
 	mux.HandleFunc("GET /surface/confirmations/{id}", d.approvalRoute(d.handleSurfaceConfirmation))
-	mux.HandleFunc("GET /surface/verb.json", d.approvalRoute(d.handleSurfaceVerb)) // row 243 V1: the bar's verbs, read-only
+	mux.HandleFunc("GET /surface/verb.json", d.approvalRoute(d.handleSurfaceVerb))         // row 243 V1: the bar's verbs, read-only
 	mux.HandleFunc("GET /surface/profiles.json", d.approvalRoute(d.handleSurfaceProfiles)) // row 269 Phase 1: read-only exec profiles
 	mux.HandleFunc("GET /surface/presets.json", d.approvalRoute(d.handleSurfacePresets))   // row 269 Phase 1: read-only launch presets
 	mux.HandleFunc("GET /surface/repos", d.approvalRoute(d.handleSurfaceReposPage))        // row 269 Phase 1: Repos & Presets UI
@@ -1185,7 +1185,7 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("POST /a2a/", d.projection.A2A)
 	mux.HandleFunc("POST /mcp/", d.projection.MCP)
 	mux.HandleFunc("POST /agui/", d.handleAGUI)
-	return trackActivityRequests(NewSessionMiddleware(d.resolver, d.credentialBudget, d.writeInternalError).Wrap(d.isProtected, d.surfaceFontPaths(mux)))
+	return trackActivityRequests(NewSessionMiddleware(d.resolver, d.credentialBudget, d.writeInternalError).Wrap(d.isProtected, d.surfaceFontPaths(d.operatorRoutes(mux))))
 }
 
 // isProtected reports whether a request must carry a valid session credential
diff --git a/host/daemon/surface_approve.go b/host/daemon/surface_approve.go
index 54e344bc..0b245e24 100644
--- a/host/daemon/surface_approve.go
+++ b/host/daemon/surface_approve.go
@@ -147,6 +147,16 @@ func (d *Daemon) readRequests(ctx context.Context, now int64, out *surfacepage.R
 	for _, row := range decided {
 		out.Decided = append(out.Decided, surfacepage.Decided(storedRequest(row)))
 	}
+	actions, err := ro.PendingOperatorActions(ctx, now)
+	if err != nil && !isOlderStore(err) {
+		out.OperatorActionsUnavailable = "the operator actions could not be read"
+	}
+	for _, row := range actions {
+		p := operatorStatusFor(row, now)
+		out.OperatorActions = append(out.OperatorActions, surfacepage.PendingOperatorAction{
+			ID: p.ID, Episode: p.Episode, Kind: p.Kind, ExpiresAt: p.ExpiresAt, Href: p.Href,
+		})
+	}
 	return nil
 }
 
diff --git a/host/daemon/surface_approve_post_guard_test.go b/host/daemon/surface_approve_post_guard_test.go
index 87dde271..317af377 100644
--- a/host/daemon/surface_approve_post_guard_test.go
+++ b/host/daemon/surface_approve_post_guard_test.go
@@ -165,7 +165,7 @@ func TestDaemonCredentialWriterSites(t *testing.T) {
 			})
 		}
 	}
-	if want := map[string]string{"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials", "keyArm": "OpenCredentials", "handleSurfaceConfirm": "OpenCredentials"}; !reflect.DeepEqual(sites, want) || readOnly < 5 {
+	if want := map[string]string{"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials", "keyArm": "OpenCredentials", "openConfirmationWriter": "OpenCredentials"}; !reflect.DeepEqual(sites, want) || readOnly < 5 {
 		t.Fatalf("writer sites %v (read-only opens %d), want exactly %v and at least 5 read-only opens", sites, readOnly, want)
 	}
 }
diff --git a/host/daemon/surface_confirm.go b/host/daemon/surface_confirm.go
index 24b39fac..de0d9d6b 100644
--- a/host/daemon/surface_confirm.go
+++ b/host/daemon/surface_confirm.go
@@ -77,6 +77,12 @@ func (d *Daemon) confirmStoreError(w http.ResponseWriter, r *http.Request, err e
 	}
 }
 
+// openConfirmationWriter is the daemon's shared non-waiting writer seam
+// for key confirmation and inert operator-action admission.
+func (d *Daemon) openConfirmationWriter() (*store.Credentials, error) {
+	return store.OpenCredentials(d.credentialsPath)
+}
+
 // consumeRefusal is C15's answer: the row raced (409), the approver or the counter changed since
 // C11 (422), else confirmStoreError.
 func (d *Daemon) consumeRefusal(w http.ResponseWriter, r *http.Request, row store.VerbRequestRow, err error) {
@@ -189,7 +195,7 @@ func (d *Daemon) handleSurfaceConfirm(w http.ResponseWriter, r *http.Request, o
 		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the "+what+" does not confirm request "+id+": "+err.Error())
 		return
 	}
-	wr, err := store.OpenCredentials(d.credentialsPath) // C13: non-waiting, before any write
+	wr, err := d.openConfirmationWriter() // C13: non-waiting, before any write
 	if err != nil {
 		d.confirmStoreError(w, r, err)
 		return
diff --git a/host/daemon/surface_home_summary.go b/host/daemon/surface_home_summary.go
index c0b7e39d..3e785066 100644
--- a/host/daemon/surface_home_summary.go
+++ b/host/daemon/surface_home_summary.go
@@ -28,7 +28,10 @@ func (d *Daemon) handleHomeSummary(w http.ResponseWriter, r *http.Request) {
 		if err := d.readDecisionList(ctx, now, &reqs); err != nil {
 			d.writeInternalErrorLog(r, err)
 		}
-		needs += len(reqs.Decisions)
+		needs += len(reqs.Decisions) + len(reqs.OperatorActions)
+		if reqs.OperatorActionsUnavailable != "" {
+			needs = -1
+		}
 	}
 	var edge *surfacepage.NowEdge
 	if d.follow != nil {
diff --git a/host/daemon/surface_operator.go b/host/daemon/surface_operator.go
new file mode 100644
index 00000000..d0b6280f
--- /dev/null
+++ b/host/daemon/surface_operator.go
@@ -0,0 +1,257 @@
+package daemon
+
+import (
+	"context"
+	"crypto/rand"
+	"encoding/hex"
+	"errors"
+	"io"
+	"net/http"
+	"path"
+	"slices"
+	"strconv"
+	"strings"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+const operatorPath = "/surface/operator-actions"
+const operatorUnavailable = "runner or reviewed preset unavailable; inspect ailang-worldd service runner status or start at a terminal"
+
+var operatorPolicy = verbs.OperatorPolicy{Start: "A", Finish: "A"} // Headless start is the only installed S7 backend.
+
+func (d *Daemon) operatorVerb(r *http.Request, argv []string, out *verbAnswer) bool {
+	a := verbs.AdmitOperator(argv, operatorPolicy)
+	switch a.Class {
+	case verbs.OperatorOther:
+		return false
+	case verbs.OperatorNever, verbs.OperatorUsage:
+		out.Outcome, out.Reason, out.CLI = "refused", a.Reason, a.CLI
+	case verbs.OperatorReview:
+		out.Outcome, out.Reason, out.CLI = "terminal", "finish preparation requires an attended terminal", a.CLI
+		if a.Intent.Kind == "start" && a.Intent.Headless {
+			out.Reason = operatorUnavailable
+			ctx, cancel := d.readCtx(r)
+			defer cancel()
+			ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+			if err != nil {
+				return true
+			}
+			defer ro.Close()
+			h, entry, err := d.operatorReview(ctx, ro, a.Intent.Preset, d.approveNow())
+			if err == nil {
+				if action, err := entry.Action(a.Intent.Episode, "", h.CatalogueGeneration); err == nil {
+					out.Outcome, out.Reason, out.Argv, out.Text = "review", "", action.Argv, action.Review
+					out.CatalogueGeneration = h.CatalogueGeneration
+				}
+			}
+		}
+	default:
+		out.Outcome, out.Reason, out.CLI = "terminal", a.Reason, a.CLI
+	}
+	return true
+}
+
+func (d *Daemon) operatorReview(ctx context.Context, db *store.Credentials, preset string, now int64) (store.OperatorHeartbeat, store.OperatorStartReview, error) {
+	h, err := db.OperatorHeartbeat(ctx)
+	if err != nil {
+		return h, store.OperatorStartReview{}, err
+	}
+	build := CurrentBuild().ExecutableSHA256
+	if build == "" || !h.Available(build, "start", now) || h.CatalogueGeneration <= 0 {
+		return h, store.OperatorStartReview{}, store.ErrOperatorRunner
+	}
+	entry, err := h.StartReview(preset)
+	served, pathErr := store.CanonicalPath(d.credentialsPath)
+	if err == nil && (pathErr != nil || entry.Store != served) {
+		err = store.ErrOperatorRunner
+	}
+	return h, entry, err
+}
+
+// Handle the whole route family before ServeMux path cleaning. Malformed
+// action paths receive a refusal and cannot become redirect destinations.
+func (d *Daemon) operatorRoutes(next http.Handler) http.Handler {
+	family := func(p string) bool { return p == operatorPath || strings.HasPrefix(p, operatorPath+"/") }
+	operator := d.approvalRoute(d.handleOperatorRoutes)
+	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		if family(r.URL.Path) || family(path.Clean(r.URL.Path)) {
+			operator(w, r)
+			return
+		}
+		next.ServeHTTP(w, r)
+	})
+}
+
+func (d *Daemon) handleOperatorRoutes(w http.ResponseWriter, r *http.Request, origins []string) {
+	if r.URL.Path != r.URL.EscapedPath() || path.Clean(r.URL.Path) != r.URL.Path || r.URL.RawQuery != "" {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "use the exact operator action path")
+		return
+	}
+	if r.URL.Path == operatorPath {
+		if r.Method == http.MethodPost {
+			d.handleOperatorRequest(w, r, origins)
+			return
+		}
+	} else {
+		id, suffix, hasSuffix := strings.Cut(strings.TrimPrefix(r.URL.Path, operatorPath+"/"), "/")
+		if !verbs.ValidOperatorActionID(id) {
+			writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "an operator action id is oa- and 32 lowercase hex digits")
+			return
+		}
+		switch {
+		case !hasSuffix || suffix == "status.json":
+			if r.Method == http.MethodGet {
+				d.handleOperatorRead(w, r, id, hasSuffix)
+				return
+			}
+		default:
+			writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no operator action route")
+			return
+		}
+	}
+	writeApprovalRefusal(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "use GET for a card/status or POST to request an action")
+}
+
+// readOperatorPOST is reached after the Host and method checks, before any body read.
+func readOperatorPOST(w http.ResponseWriter, r *http.Request, origins []string, cap int64) ([]byte, bool) {
+	origin, fetch, media := r.Header.Values("Origin"), r.Header.Values("Sec-Fetch-Site"), r.Header.Values("Content-Type")
+	switch {
+	case len(r.Header.Values("Authorization")) != 0 || len(r.Header.Values("Cookie")) != 0:
+		writeApprovalRefusal(w, http.StatusForbidden, "Forbidden", "use the operator card without a session or cookie")
+	case len(origin) != 1 || !allowedApprovalOrigin(r.Header.Get("Origin"), origins):
+		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", "use an exact approval origin")
+	case len(fetch) > 1 || !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")):
+		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", "Sec-Fetch-Site must be same-origin when sent")
+	case len(media) != 1 || !jsonMediaType(r.Header.Get("Content-Type")):
+		writeApprovalRefusal(w, http.StatusUnsupportedMediaType, "UnsupportedMediaType", "use application/json")
+	default:
+		raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, cap))
+		var tooLarge *http.MaxBytesError
+		switch {
+		case errors.As(err, &tooLarge):
+			writeApprovalRefusal(w, http.StatusRequestEntityTooLarge, "TooLarge", "the operator action body exceeds its byte limit")
+		case err != nil:
+			writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "the operator action body could not be read")
+		default:
+			return raw, true
+		}
+	}
+	return nil, false
+}
+
+type operatorRequestBody struct {
+	Argv                []string `json:"argv"`
+	OperationID         string   `json:"operationId"`
+	CatalogueGeneration int64    `json:"catalogueGeneration"`
+	Task                string   `json:"task"`
+}
+
+func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, origins []string) {
+	raw, ok := readOperatorPOST(w, r, origins, 64<<10)
+	if !ok {
+		return
+	}
+	var body operatorRequestBody
+	if approval.DecodeOperatorObject(raw, &body, "argv", "operationId", "catalogueGeneration", "task") != nil ||
+		len(body.Task) == 0 || len(body.Task) > approval.MaxOperatorTaskBytes || strings.ContainsRune(body.Task, 0) {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send argv, operationId, catalogueGeneration and a task of 1–8192 UTF-8 bytes")
+		return
+	}
+	a := verbs.AdmitOperator(body.Argv, operatorPolicy)
+	if a.Class != verbs.OperatorReview || a.Intent.Kind != "start" || !a.Intent.Headless {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "review a supported headless start first")
+		return
+	}
+	if body.CatalogueGeneration <= 0 || !store.ValidOperatorOperationID(body.OperationID) {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "use a positive catalogue generation and 32 lowercase hex operationId")
+		return
+	}
+
+	d.humanMu.Lock()
+	defer d.humanMu.Unlock()
+	ctx, cancel := d.readCtx(r)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	defer ro.Close()
+	now := d.approveNow()
+	if row, err := ro.OperatorActionForOperation(ctx, body.OperationID, now); !errors.Is(err, store.ErrOperatorMissing) {
+		if err != nil {
+			d.operatorError(w, err)
+			return
+		}
+		frozen, err := approval.DecodeOperatorAction(row.Action)
+		if err != nil || !slices.Equal(frozen.Argv, a.Intent.Argv) || frozen.Task != body.Task ||
+			!slices.Contains(frozen.Resolution, approval.OperatorResolved{Name: "mint_catalogue", Value: strconv.FormatInt(body.CatalogueGeneration, 10)}) {
+			d.operatorError(w, store.ErrOperatorConflict)
+			return
+		}
+		writeApprovalJSON(w, http.StatusOK, operatorStatusFor(row, now))
+		return
+	}
+	wr, err := d.openConfirmationWriter()
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	defer wr.Close()
+	now = d.approveNow()
+	h, entry, err := d.operatorReview(ctx, wr, a.Intent.Preset, now)
+	if err == nil && h.CatalogueGeneration != body.CatalogueGeneration {
+		err = store.ErrOperatorRunner
+	}
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	action, err := entry.Action(a.Intent.Episode, body.Task, h.CatalogueGeneration)
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	var id [16]byte
+	var challenge [32]byte
+	if _, err := rand.Read(id[:]); err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	if _, err := rand.Read(challenge[:]); err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	action.ID = "oa-" + hex.EncodeToString(id[:])
+	row, err := wr.RequestOperatorAction(ctx, store.OperatorRequest{Action: action, OperationID: body.OperationID, Challenge: challenge,
+		CreatedAt: now, ExpiresAt: now + int64(approval.OperatorActionLifetime/time.Second)}, now)
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	writeApprovalJSON(w, http.StatusOK, operatorStatusFor(row, now))
+}
+
+func (d *Daemon) operatorError(w http.ResponseWriter, err error) {
+	switch {
+	case isOlderStore(err), errors.Is(err, store.ErrOperatorRunner):
+		writeApprovalRefusal(w, http.StatusServiceUnavailable, "RunnerUnavailable", operatorUnavailable+"; review again")
+	case store.IsWriterAlreadyActive(err):
+		w.Header().Set("Retry-After", "1")
+		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Busy", "credentials writer busy; retry")
+	case errors.Is(err, store.ErrOperatorMissing):
+		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no operator action")
+	case errors.Is(err, store.ErrOperatorCapacity):
+		writeApprovalRefusal(w, http.StatusConflict, "QueueFull", "sixteen actions are already pending")
+	case errors.Is(err, store.ErrOperatorConflict):
+		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "operator action changed or conflicts with existing work")
+	case errors.Is(err, store.ErrOperatorInput), errors.Is(err, approval.ErrOperatorTooLarge):
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "invalid reviewed action")
+	default:
+		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Disconnected", "operator actions could not be read or written; retry")
+	}
+}
diff --git a/host/daemon/surface_operator_read.go b/host/daemon/surface_operator_read.go
new file mode 100644
index 00000000..0bdf31ca
--- /dev/null
+++ b/host/daemon/surface_operator_read.go
@@ -0,0 +1,126 @@
+package daemon
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"errors"
+	"html/template"
+	"net/http"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+type operatorStatus struct {
+	V            int    `json:"v"`
+	ID           string `json:"id"`
+	Episode      string `json:"episode"`
+	Kind         string `json:"kind"`
+	State        string `json:"state"`
+	Href         string `json:"href,omitempty"`
+	StatusHref   string `json:"statusHref,omitempty"`
+	ReportHref   string `json:"reportHref,omitempty"`
+	ConfirmHref  string `json:"confirmHref,omitempty"`
+	Reason       string `json:"reason,omitempty"`
+	Now          int64  `json:"now"`
+	CreatedAt    int64  `json:"createdAt"`
+	ExpiresAt    int64  `json:"expiresAt"`
+	TerminalAt   int64  `json:"terminalAt,omitempty"`
+	RunnerSeenAt int64  `json:"runnerSeenAt,omitempty"`
+}
+
+func operatorStatusFor(row store.OperatorActionRow, now int64) operatorStatus {
+	out := operatorStatus{V: 1, ID: row.ID, Episode: row.Episode, Kind: row.Kind, State: row.State, Now: now,
+		CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, TerminalAt: row.TerminalAt}
+	switch row.State {
+	case "requested", "confirmed", "claimed", "running", "awaiting-pr-confirmation", "succeeded", "failed", "expired", "cancelled", "uncertain":
+	default:
+		out.State, out.Reason = "unavailable", "unknown stored action state"
+		return out
+	}
+	if !verbs.ValidOperatorActionID(row.ID) || !verbs.EpisodeIDPattern.MatchString(row.Episode) ||
+		(row.Kind != "start" && row.Kind != "finish") {
+		out.State, out.Reason = "unavailable", "invalid stored action identity"
+		return out
+	}
+	out.Href = operatorPath + "/" + row.ID
+	out.StatusHref = out.Href + "/status.json"
+	if row.State == "succeeded" || row.State == "failed" || row.State == "uncertain" {
+		out.ReportHref = "/v1/episodes/" + row.Episode + "/report"
+	}
+
+	return out
+}
+
+var operatorCard = template.Must(template.New("operator action").Parse(`<!doctype html>
+<html lang="en"><meta charset="utf-8"><title>Operator action</title>
+<body><h1>{{.Status.ID}}</h1><p>{{.Status.State}}</p>{{if .Status.Reason}}<p>{{.Status.Reason}}</p>{{end}}
+{{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre><h2>Task</h2><pre>{{.Task}}</pre>{{end}}
+{{if .Status.StatusHref}}<p><a href="{{.Status.StatusHref}}">Action status (JSON)</a></p>{{end}}
+{{if .Status.ReportHref}}<p><a href="{{.Status.ReportHref}}">Episode report</a></p>{{end}}
+</body></html>`))
+
+func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id string, statusJSON bool) {
+	ctx, cancel := d.readCtx(r)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	defer ro.Close()
+	now := d.approveNow()
+	row, err := ro.OperatorAction(ctx, id, now)
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	status := operatorStatusFor(row, now)
+	h, heartbeatErr := ro.OperatorHeartbeat(ctx)
+	status.RunnerSeenAt = h.SeenAt
+	if (heartbeatErr != nil || !h.Available(CurrentBuild().ExecutableSHA256, row.Kind, now)) && status.Reason == "" {
+		status.Reason = "runner disconnected; showing durable action status"
+	}
+	if status.Href != "" && row.Kind == "start" && (row.State == "running" || row.State == "succeeded" || row.State == "uncertain" || row.State == "failed") {
+		launch, err := ro.ReadOperatorLaunch(ctx, id)
+		if err == nil {
+			switch {
+			case launch.NoSpawn:
+				status.State = "spawn-refused"
+			case launch.Phase == "reserved" || launch.Receipt.Status == "launched":
+				status.State = "uncertain"
+			default:
+				status.State = launch.Receipt.Status
+			}
+			if status.ReportHref == "" {
+				status.ReportHref = "/v1/episodes/" + row.Episode + "/report"
+			}
+			if status.State == "uncertain" {
+				status.Reason = "launch fenced; inspect ailang-worldd service runner launch-status " + id
+			}
+		} else if !errors.Is(err, store.ErrOperatorMissing) {
+			d.operatorError(w, err)
+			return
+		}
+	}
+	if statusJSON {
+		writeApprovalJSON(w, http.StatusOK, status)
+		return
+	}
+	card := struct {
+		Status       operatorStatus
+		Review, Task string
+	}{Status: status}
+	if len(row.Action) != 0 && status.Href != "" {
+		a, err := approval.DecodeOperatorAction(row.Action)
+		digest := sha256.Sum256(row.Action)
+		if err != nil || a.ID != id || a.Episode != row.Episode || a.Kind != row.Kind || !bytes.Equal(digest[:], row.Digest) {
+			d.operatorError(w, store.ErrOperatorConflict)
+			return
+		}
+		card.Review, card.Task = a.Review, a.Task
+	}
+	setSurfaceHeaders(w)
+	_ = operatorCard.Execute(w, card)
+}
diff --git a/host/daemon/surface_operator_test.go b/host/daemon/surface_operator_test.go
new file mode 100644
index 00000000..7350c645
--- /dev/null
+++ b/host/daemon/surface_operator_test.go
@@ -0,0 +1,437 @@
+package daemon
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"database/sql"
+	"encoding/json"
+	"fmt"
+	"io"
+	"net"
+	"net/http"
+	"net/http/httptest"
+	"net/url"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"sync"
+	"sync/atomic"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+const opTask = "# Task\n<script>alert(1)</script> ; $(literal) &"
+const opOrigin = "http://localhost:7644"
+
+type operatorLog struct {
+	mu sync.Mutex
+	b  bytes.Buffer
+}
+
+func (l *operatorLog) Write(p []byte) (int, error) {
+	l.mu.Lock()
+	defer l.mu.Unlock()
+	return l.b.Write(p)
+}
+
+func (l *operatorLog) String() string {
+	l.mu.Lock()
+	defer l.mu.Unlock()
+	return l.b.String()
+}
+
+type operatorHTTP struct {
+	approvalRig
+	handler http.Handler
+	server  *httptest.Server
+	now     atomic.Int64
+	h       store.OperatorHeartbeat
+	fake    string
+	log     operatorLog
+}
+
+func newOperatorHTTP(t *testing.T) *operatorHTTP {
+	t.Helper()
+	f := &operatorHTTP{approvalRig: newApprovalRig(t)}
+	f.now.Store(apT0)
+	f.d.approveNow = f.now.Load
+	f.d.errLog = &f.log
+	f.d.cfg.SurfaceOrigins = []string{opOrigin}
+	root := t.TempDir()
+	f.fake = filepath.Join(root, "ran")
+	if err := os.WriteFile(filepath.Join(root, "fake-codex"), []byte("#!/bin/sh\nprintf x > \"${0%/*}/ran\"\n"), 0700); err != nil {
+		t.Fatal(err)
+	}
+	physical, err := store.CanonicalPath(f.d.credentialsPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	build := CurrentBuild().ExecutableSHA256
+	entry := store.OperatorStartReview{Store: physical, Workspace: root, Sessions: root, Binary: filepath.Join(root, "fake-codex"), Wall: "1m",
+		Start: store.EpisodeStart{Schema: "world/episode-start/v1", Agent: "codex", Model: "fixture", TTLS: 60,
+			Preset: store.StartPreset{Name: "sol", SHA256: "sha256:fixture"}, Overrides: []string{}}}
+	f.h = store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "fixture"}, Protocol: store.OperatorRunnerProtocol, Build: build,
+		Readiness: "start", SeenAt: apT0, Catalogue: opJSON(map[string]store.OperatorStartReview{"sol": entry})}
+	f.publish()
+	f.handler = f.d.Handler()
+	f.server = httptest.NewServer(f.handler)
+	t.Cleanup(f.server.Close)
+	return f
+}
+
+func opJSON(v any) []byte {
+	b, err := json.Marshal(v)
+	if err != nil {
+		panic(err)
+	}
+	return b
+}
+
+func (f *operatorHTTP) publish() {
+	f.write(func(ctx context.Context, c *store.Credentials) error {
+		gen, err := c.PutOperatorHeartbeat(ctx, f.h)
+		f.h.Owner.Generation = gen
+		return err
+	})
+}
+
+func (f *operatorHTTP) do(method, path string, body []byte, edit func(*http.Request)) *httptest.ResponseRecorder {
+	f.t.Helper()
+	req, err := http.NewRequestWithContext(boundedTestContext(f.t), method, f.server.URL+path, bytes.NewReader(body))
+	if err != nil {
+		f.t.Fatal(err)
+	}
+	if method == "POST" {
+		req.Header.Set("Origin", opOrigin)
+		req.Header.Set("Content-Type", "application/json")
+		req.Header.Set("Sec-Fetch-Site", "same-origin")
+	}
+	if edit != nil {
+		edit(req)
+	}
+	client := *f.server.Client()
+	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
+	res, err := client.Do(req)
+	if err != nil {
+		f.t.Fatal(err)
+	}
+	defer res.Body.Close()
+	out := httptest.NewRecorder()
+	out.Code, out.HeaderMap = res.StatusCode, res.Header
+	if _, err := io.Copy(out.Body, res.Body); err != nil {
+		f.t.Fatal(err)
+	}
+	return out
+}
+
+type opEvidence struct {
+	db, world                              [32]byte
+	actions, operations, effects, sessions int
+}
+
+func (f *operatorHTTP) evidence() (out opEvidence) {
+	f.t.Helper()
+	fileHash := func(p string) [32]byte {
+		b, err := os.ReadFile(p)
+		if err != nil {
+			f.t.Fatal(err)
+		}
+		return sha256.Sum256(b)
+	}
+	out.db, out.world = fileHash(f.d.credentialsPath), fileHash(f.d.cfg.DBPath)
+	u := url.URL{Scheme: "file", Path: f.d.credentialsPath, RawQuery: "mode=ro"}
+	db, err := sql.Open("sqlite", u.String())
+	if err != nil {
+		f.t.Fatal(err)
+	}
+	defer db.Close()
+	err = db.QueryRowContext(boundedTestContext(f.t), `SELECT (SELECT count(*) FROM operator_actions),
+		(SELECT count(*) FROM operator_action_steps WHERE kind='operation'),
+		(SELECT count(*) FROM operator_action_steps WHERE kind<>'operation'), (SELECT count(*) FROM session_credentials)`).
+		Scan(&out.actions, &out.operations, &out.effects, &out.sessions)
+	if err != nil || out.effects != 0 || out.sessions != 0 {
+		f.t.Fatalf("a runner or effect crossed the browser boundary: %+v, %v", out, err)
+	}
+	if _, err := os.Stat(f.fake); !os.IsNotExist(err) {
+		f.t.Fatalf("unexpected fake child marker: %v", err)
+	}
+	return out
+}
+
+func opRequest(n int) operatorRequestBody {
+	o := approval.OperatorOptions{Preset: "sol", Headless: true}
+	return operatorRequestBody{Argv: approval.OperatorArgv("start", fmt.Sprintf("ep-%d", n), o),
+		Task: opTask, OperationID: fmt.Sprintf("%032x", n), CatalogueGeneration: 1}
+}
+
+func (f *operatorHTTP) request(body operatorRequestBody) operatorStatus {
+	f.t.Helper()
+	r := f.do("POST", operatorPath, opJSON(body), nil)
+	var s operatorStatus
+	if err := json.Unmarshal(r.Body.Bytes(), &s); r.Code != 200 || err != nil || s.ID == "" || s.Href == "" {
+		f.t.Fatalf("request: %d %s (%v)", r.Code, r.Body, err)
+	}
+	return s
+}
+
+func (f *operatorHTTP) refuse(before opEvidence, r *httptest.ResponseRecorder, status int, class string) {
+	f.t.Helper()
+	requireRefusal(f.t, class, r, status, class)
+	if r.Header().Get("Location") != "" || before != f.evidence() {
+		f.t.Fatal("refusal changed durable evidence or built a redirect")
+	}
+}
+
+type operatorBodySpy struct{ reads int }
+
+func (b *operatorBodySpy) Read([]byte) (int, error) { b.reads++; return 0, io.EOF }
+func (b *operatorBodySpy) Close() error             { return nil }
+
+func (f *operatorHTTP) sub(name string, run func(*testing.T)) {
+	f.t.Run(name, func(t *testing.T) {
+		old := f.t
+		f.t = t
+		defer func() { f.t = old }()
+		run(t)
+	})
+}
+
+func TestOperatorAdmissionHTTP(t *testing.T) {
+	f := newOperatorHTTP(t)
+	valid := opJSON(opRequest(1))
+	before := f.evidence()
+	refuse := func(path, method string, raw []byte, code int, class string, edit func(*http.Request)) {
+		f.refuse(before, f.do(method, path, raw, edit), code, class)
+	}
+	hdr := func(k, v string) func(*http.Request) {
+		return func(r *http.Request) {
+			if v == "" {
+				r.Header.Del(k)
+			} else {
+				r.Header.Set(k, v)
+			}
+		}
+	}
+	f.sub("Misdirected", func(t *testing.T) {
+		refuse(operatorPath, "POST", valid, 421, "MisdirectedRequest", func(r *http.Request) { r.Host = "localhost:7644.evil" })
+		spy := &operatorBodySpy{}
+		r := httptest.NewRequest("POST", operatorPath, spy)
+		r = r.WithContext(context.WithValue(r.Context(), http.LocalAddrContextKey, &net.TCPAddr{Port: 7644}))
+		r.Host = "localhost:7644.evil"
+		w := httptest.NewRecorder()
+		f.handler.ServeHTTP(w, r)
+		f.refuse(before, w, 421, "MisdirectedRequest")
+		if spy.reads != 0 {
+			t.Fatal("foreign Host reached body read")
+		}
+	})
+	f.sub("CrossOrigin", func(t *testing.T) {
+		for _, origin := range []string{"", "null", opOrigin + ".evil", "https://localhost:7644", "http://localhost:7645"} {
+			refuse(operatorPath, "POST", valid, 403, "CrossOrigin", hdr("Origin", origin))
+		}
+		refuse(operatorPath, "POST", valid, 403, "CrossOrigin", func(r *http.Request) { r.Header.Add("Origin", opOrigin) })
+	})
+	f.sub("FetchSite", func(t *testing.T) {
+		for _, site := range []string{"same-site", "cross-site", "same-origin, cross-site", "Same-Origin"} {
+			refuse(operatorPath, "POST", valid, 403, "CrossOrigin", hdr("Sec-Fetch-Site", site))
+		}
+		refuse(operatorPath, "POST", valid, 403, "CrossOrigin", func(r *http.Request) { r.Header.Add("Sec-Fetch-Site", "cross-site") })
+	})
+	f.sub("MediaAndBody", func(t *testing.T) {
+		refuse(operatorPath, "POST", valid, 415, "UnsupportedMediaType", hdr("Content-Type", "text/plain"))
+		refuse(operatorPath, "POST", valid, 415, "UnsupportedMediaType", func(r *http.Request) { r.Header.Add("Content-Type", "text/plain") })
+		refuse(operatorPath, "POST", bytes.Repeat([]byte(" "), (64<<10)+1), 413, "TooLarge", nil)
+	})
+	f.sub("UnsupportedMethodAndPath", func(t *testing.T) {
+		id := "/oa-" + strings.Repeat("a", 32)
+		for _, path := range []string{operatorPath, operatorPath + id, operatorPath + id + "/status.json"} {
+			refuse(path, "OPTIONS", nil, 405, "MethodNotAllowed", nil)
+		}
+		for _, path := range []string{operatorPath + id + "/", operatorPath + id + "?next=https://elsewhere.invalid", operatorPath + "/" + id,
+			operatorPath + "/x/.." + id, "/surface/x/../operator-actions" + id, operatorPath + "/%6fa-" + strings.Repeat("a", 32), operatorPath + strings.ToUpper(id),
+			operatorPath + id + "x", operatorPath + id[:len(id)-1], operatorPath + id + "%0a"} {
+			refuse(path, "GET", nil, 400, "BadRequest", nil)
+		}
+		refuse(operatorPath+id+"/status.json/x", "GET", nil, 404, "NotFound", nil)
+		refuse(operatorPath+id+"/confirm", "POST", valid, 404, "NotFound", nil)
+	})
+	f.sub("StrictJSONAndLimits", func(t *testing.T) {
+		invalid := [][]byte{nil, []byte("[]"), []byte(`{"argv":"start"}`), []byte(`{"argv":[42]}`), append(bytes.Clone(valid), []byte(" {}")...),
+			append([]byte(`{"task":"duplicate",`), valid[1:]...), append([]byte(`{"Task":"alias",`), valid[1:]...), append([]byte(`{"other":true,`), valid[1:]...),
+			bytes.Replace(valid, []byte("Task"), []byte{0xff}, 1), bytes.Replace(valid, []byte("Task"), []byte{0}, 1),
+			bytes.Replace(valid, []byte("Task"), []byte(`\u0000`), 1)}
+		for _, path := range invalid {
+			refuse(operatorPath, "POST", path, 400, "BadRequest", nil)
+		}
+		for _, text := range []string{"", strings.Repeat("x", (8<<10)+1)} {
+			b := opRequest(1)
+			b.Task = text
+			refuse(operatorPath, "POST", opJSON(b), 400, "BadRequest", nil)
+		}
+		for _, id := range []string{"", strings.Repeat("A", 32), "x" + opRequest(1).OperationID} {
+			b := opRequest(1)
+			b.OperationID = id
+			refuse(operatorPath, "POST", opJSON(b), 400, "BadRequest", nil)
+		}
+	})
+	f.sub("ExplicitHold", func(t *testing.T) {
+		refuse(operatorPath, "POST", valid, 403, "Forbidden", hdr("Authorization", "Bearer session"))
+		refuse(operatorPath, "POST", valid, 403, "Forbidden", hdr("Cookie", "session=x"))
+	})
+	checkBar := func(line, want string) {
+		t.Helper()
+		result := f.do("GET", "/surface/verb.json?q="+url.QueryEscape(line), nil, nil)
+		var a verbAnswer
+		if err := json.Unmarshal(result.Body.Bytes(), &a); err != nil || result.Code != 200 || a.Outcome != want ||
+			(want == "review" && (a.CatalogueGeneration != 1 || !reflect.DeepEqual(a.Argv, opRequest(1).Argv))) || before != f.evidence() {
+			t.Fatalf("bar %q: %d %s (%v)", line, result.Code, result.Body, err)
+		}
+	}
+	checkBar("episode start ep-1 --preset=sol --headless --launch", "review")
+	for _, line := range []string{"world-publish", "approver enroll", "approver revoke", "session mint", "session revoke", "service runner start"} {
+		checkBar(line, "refused")
+	}
+	checkBar("approve mr-0123456789abcdef", "navigate")
+	checkBar("episode report ep-1", "result")
+	for _, v := range []struct {
+		generation int64
+		preset     string
+	}{{2, "sol"}, {1, "absent"}} {
+		b := opRequest(1)
+		b.CatalogueGeneration, b.Argv[5] = v.generation, v.preset
+		refuse(operatorPath, "POST", opJSON(b), 503, "RunnerUnavailable", nil)
+	}
+	s := f.request(opRequest(1))
+	if after := f.evidence(); after.world != before.world || after.actions != 1 || after.operations != 1 || s.State != "requested" {
+		t.Fatalf("inert request evidence: %+v status=%+v", after, s)
+	}
+	f.sub("Status", func(t *testing.T) { testOperatorStatus(f, s) })
+	if duplicate := f.request(opRequest(1)); duplicate.ID != s.ID || f.evidence().actions != 1 {
+		t.Fatal("same operation did not recover the existing row")
+	}
+	b := opRequest(1)
+	b.OperationID = opRequest(90).OperationID
+	if identical := f.request(b); identical.ID != s.ID || f.evidence().actions != 1 {
+		t.Fatal("identical intent from a different operation created a row")
+	}
+	before = f.evidence()
+	b.Task += " changed"
+	refuse(operatorPath, "POST", opJSON(b), 409, "Conflict", nil)
+	f.now.Store(apT0 + 10)
+	checkBar("episode start ep-1 --preset sol --headless --launch", "terminal")
+	refuse(operatorPath, "POST", opJSON(opRequest(2)), 503, "RunnerUnavailable", nil)
+	if recovered := f.request(opRequest(1)); recovered.ID != s.ID || before != f.evidence() {
+		t.Fatal("lost request response did not survive catalogue expiry")
+	}
+	f.now.Store(apT0)
+	f.write(func(ctx context.Context, db *store.Credentials) error {
+		entry, err := f.h.StartReview("sol")
+		if err != nil {
+			return err
+		}
+		for i := 2; i <= 16; i++ {
+			b := opRequest(i)
+			a, err := entry.Action(fmt.Sprintf("ep-%d", i), b.Task, 1)
+			if err != nil {
+				return err
+			}
+			a.ID = "oa-" + b.OperationID
+			_, err = db.RequestOperatorAction(ctx, store.OperatorRequest{Action: a, OperationID: b.OperationID,
+				Challenge: sha256.Sum256([]byte(a.ID)), CreatedAt: apT0, ExpiresAt: apT0 + 600}, apT0)
+			if err != nil {
+				return err
+			}
+		}
+		return nil
+	})
+	before = f.evidence()
+	if before.actions != 16 {
+		t.Fatalf("pending cap positive control: %d", before.actions)
+	}
+	refuse(operatorPath, "POST", opJSON(opRequest(17)), 409, "QueueFull", nil)
+	f.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.CancelUnstartedOperatorAction(ctx, s.ID, apT0)
+	})
+	if r := f.do("GET", "/surface/requests.json", nil, nil); r.Code != 200 || strings.Contains(r.Body.String(), s.ID) {
+		t.Fatalf("cancelled row in Needs you: %d %s", r.Code, r.Body)
+	}
+}
+
+func testOperatorStatus(f *operatorHTTP, s operatorStatus) {
+	t := f.t
+	defer f.now.Store(apT0)
+	before := f.evidence()
+	get := func(p string) *httptest.ResponseRecorder { return f.do("GET", p, nil, nil) }
+	card := get(s.Href)
+	if card.Code != 200 || card.Header().Get("Content-Security-Policy") != surfaceCSP ||
+		card.Header().Get("Cache-Control") != "no-store" || !strings.Contains(card.Body.String(), "&lt;script&gt;alert(1)&lt;/script&gt;") ||
+		strings.Contains(card.Body.String(), "<script>") || !strings.Contains(card.Body.String(), s.StatusHref) {
+		t.Fatalf("stable protected card: %d %s", card.Code, card.Body)
+	}
+	for _, target := range []string{s.Href, s.StatusHref} {
+		f.refuse(before, f.do("GET", target, nil, func(r *http.Request) { r.Host = "other.invalid:7644" }), 421, "MisdirectedRequest")
+	}
+	out := get(s.StatusHref)
+	var current operatorStatus
+	if err := json.Unmarshal(out.Body.Bytes(), &current); out.Code != 200 || err != nil || current.ID != s.ID ||
+		current.State != "requested" || current.ExpiresAt != apT0+600 || current.RunnerSeenAt != apT0 {
+		t.Fatalf("status: %d %s (%v)", out.Code, out.Body, err)
+	}
+	checkNeeds := func(want int) {
+		t.Helper()
+		r := get("/surface/requests.json")
+		sum := get("/surface/home-summary.json")
+		var needs surfacepage.RequestsJSON
+		var home surfacepage.HomeSummary
+		if json.Unmarshal(r.Body.Bytes(), &needs) != nil || json.Unmarshal(sum.Body.Bytes(), &home) != nil || len(needs.OperatorActions) != want ||
+			len(home.Items) == 0 || home.Items[0].Value != fmt.Sprint(want) ||
+			want == 1 && needs.OperatorActions[0].Href != s.Href {
+			t.Fatalf("Needs you diverged from summary: %s / %s", r.Body, sum.Body)
+		}
+		for _, text := range []string{r.Body.String(), sum.Body.String(), out.Body.String(), f.log.String()} {
+			if strings.Contains(text, "alert(1)") || strings.Contains(text, "challenge") {
+				t.Fatalf("private review or ceremony leaked: %q", text)
+			}
+		}
+	}
+	checkNeeds(1)
+	f.now.Store(apT0 - 1)
+	checkNeeds(0)
+	f.now.Store(apT0 + 10)
+	if r := get(s.StatusHref); r.Code != 200 || !strings.Contains(r.Body.String(), "runner disconnected") || !strings.Contains(r.Body.String(), `"state":"requested"`) {
+		t.Fatalf("last durable status during disconnect: %d %s", r.Code, r.Body)
+	}
+	f.now.Store(apT0 + 600)
+	checkNeeds(0)
+	if r := get(s.StatusHref); r.Code != 200 || !strings.Contains(r.Body.String(), `"state":"expired"`) || f.evidence() != before {
+		t.Fatalf("expiry projection altered durable evidence: %d %s", r.Code, r.Body)
+	}
+}
+
+func TestOperatorLinksProjection(t *testing.T) {
+	r := store.OperatorActionRow{ID: "oa-" + strings.Repeat("a", 32), Episode: "episode-1", Kind: "start", State: "succeeded"}
+	if s := operatorStatusFor(r, 1); s.ReportHref != "/v1/episodes/episode-1/report" || s.StatusHref != operatorPath+"/"+r.ID+"/status.json" {
+		t.Fatalf("exact link control: %+v", s)
+	}
+	for _, id := range []string{r.ID + "\n", r.ID + "/x", "x" + r.ID, strings.ToUpper(r.ID), "https://github.com/other/repo/pull/1?x=y", "//elsewhere.invalid/"} {
+		bad := r
+		bad.ID = id
+		if s := operatorStatusFor(bad, 1); s.Href != "" || s.StatusHref != "" || s.ReportHref != "" || s.State != "unavailable" {
+			t.Fatalf("untrusted action link: %+v", s)
+		}
+	}
+	for _, bad := range []store.OperatorActionRow{
+		{ID: r.ID, Episode: "../x", Kind: r.Kind, State: r.State},
+		{ID: r.ID, Episode: r.Episode, Kind: r.Kind, State: "other"},
+		{ID: r.ID, Episode: r.Episode, Kind: "other", State: r.State},
+	} {
+		if s := operatorStatusFor(bad, 1); s.Href != "" || s.StatusHref != "" || s.ReportHref != "" || s.State != "unavailable" {
+			t.Fatalf("unsafe or unclassified state projection: %+v", s)
+		}
+	}
+}
diff --git a/host/daemon/surface_verb.go b/host/daemon/surface_verb.go
index 52633802..208e3a1c 100644
--- a/host/daemon/surface_verb.go
+++ b/host/daemon/surface_verb.go
@@ -16,25 +16,26 @@ const VerbSchema = "world/surface-verb/v1"
 // verbAnswer is one answer of the verb route. Every answer carries v, schema, q, argv, verb, class
 // and outcome; the rest belong to their outcome.
 type verbAnswer struct {
-	V       int      `json:"v"`
-	Schema  string   `json:"schema"`
-	Q       string   `json:"q"`
-	Argv    []string `json:"argv"`
-	Verb    string   `json:"verb"`
-	Class   string   `json:"class"`
-	Outcome string   `json:"outcome"`
-	Reason  string   `json:"reason,omitempty"`
-	CLI     string   `json:"cli,omitempty"`
-	Href    string   `json:"href,omitempty"`
-	Text    string   `json:"text,omitempty"`
-	JSON    string   `json:"json,omitempty"`
-	Stderr  string   `json:"stderr,omitempty"`
-	Exit    int      `json:"exit,omitempty"`
+	V                   int      `json:"v"`
+	Schema              string   `json:"schema"`
+	Q                   string   `json:"q"`
+	Argv                []string `json:"argv"`
+	Verb                string   `json:"verb"`
+	Class               string   `json:"class"`
+	Outcome             string   `json:"outcome"`
+	Reason              string   `json:"reason,omitempty"`
+	CLI                 string   `json:"cli,omitempty"`
+	Href                string   `json:"href,omitempty"`
+	Text                string   `json:"text,omitempty"`
+	JSON                string   `json:"json,omitempty"`
+	Stderr              string   `json:"stderr,omitempty"`
+	Exit                int      `json:"exit,omitempty"`
+	CatalogueGeneration int64    `json:"catalogueGeneration,omitempty"`
 }
 
 // handleSurfaceVerb is GET /surface/verb.json?q=<line>, behind approvalRoute: the bar's line through
 // host/verbs' one grammar. It reads (episode report, the /v1 route's own builder; why, host/why's walk) and writes nothing:
-// no store write, no credentials store, no process. An operator's typo is a 200 refusal.
+// no store write or process. An operator's typo is a 200 refusal.
 func (d *Daemon) handleSurfaceVerb(w http.ResponseWriter, r *http.Request, _ []string) {
 	query := r.URL.Query()
 	if !query.Has("q") {
@@ -50,6 +51,9 @@ func (d *Daemon) handleSurfaceVerb(w http.ResponseWriter, r *http.Request, _ []s
 	}
 	v, args := verbs.Lookup(argv)
 	a.Argv, a.Verb, a.Class, a.CLI, a.Reason = argv, v.Words, string(v.Class), v.CLI, v.Reason
+	if d.operatorVerb(r, argv, &a) {
+		return
+	}
 	switch {
 	case v.Class == verbs.Address:
 		a.Outcome = "address"
diff --git a/host/daemon/surface_verb_test.go b/host/daemon/surface_verb_test.go
index c93d6529..4526734a 100644
--- a/host/daemon/surface_verb_test.go
+++ b/host/daemon/surface_verb_test.go
@@ -128,8 +128,8 @@ func TestSurfaceVerbBoundsAndErrors(t *testing.T) {
 	assertErrorClass(t, requestRecorder(t, s.d, "GET", "/v1/episodes/w1/report", nil), 500, "Internal")
 }
 
-// VS1: host/verbs imports no I/O package; the route's file reaches no store, credentials store or
-// process; approve's id pattern is host/store's.
+// VS1: host/verbs stays pure. The bar delegates operator review to one
+// read-only adapter, and approve's id pattern is host/store's.
 func TestSurfaceVerbIsReadOnly(t *testing.T) {
 	importsOf := func(path string) (map[string]bool, *ast.File) {
 		t.Helper()
@@ -176,7 +176,7 @@ func TestSurfaceVerbIsReadOnly(t *testing.T) {
 		}
 		return true
 	})
-	if want := map[string]bool{"readCtx": true, "episodeReport": true, "writeInternalErrorLog": true, "verbEpisodeReport": true, "readDeadline": true, "verbWhy": true, "verbProfileList": true, "verbPresetList": true}; !reflect.DeepEqual(calls, want) {
+	if want := map[string]bool{"readCtx": true, "episodeReport": true, "writeInternalErrorLog": true, "verbEpisodeReport": true, "readDeadline": true, "verbWhy": true, "verbProfileList": true, "verbPresetList": true, "operatorVerb": true}; !reflect.DeepEqual(calls, want) {
 		t.Errorf("surface_verb.go reaches d.%v, want exactly %v", calls, want)
 	}
 	for _, id := range []string{"mr-0123456789abcdef", "mr-0123456789ABCDEF", "mr-0123456789abcde", "mr-0123456789abcdef0", "xr-0123456789abcdef", "mr-012345678 abcdef"} {
diff --git a/host/store/operator_actions.go b/host/store/operator_actions.go
index 176d93c8..58cc2502 100644
--- a/host/store/operator_actions.go
+++ b/host/store/operator_actions.go
@@ -37,6 +37,9 @@ type OperatorRequest struct {
 
 var operatorOperationID = regexp.MustCompile(`^[0-9a-f]{32}$`)
 
+// ValidOperatorOperationID is the exact browser retry key grammar.
+func ValidOperatorOperationID(id string) bool { return operatorOperationID.MatchString(id) }
+
 func (c *Credentials) operatorAccess(ctx context.Context) error {
 	if err := requireDeadline(ctx); err != nil {
 		return err
@@ -68,6 +71,44 @@ func (c *Credentials) OperatorAction(ctx context.Context, id string, now int64)
 	return row, err
 }
 
+// OperatorActionForOperation recovers an inert request after a lost response.
+func (c *Credentials) OperatorActionForOperation(ctx context.Context, id string, now int64) (OperatorActionRow, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return OperatorActionRow{}, err
+	}
+	var action string
+	err := c.db.QueryRowContext(ctx, `SELECT action_id FROM operator_action_steps WHERE kind='operation' AND identity=?`, id).Scan(&action)
+	if errors.Is(err, sql.ErrNoRows) {
+		return OperatorActionRow{}, ErrOperatorMissing
+	}
+	if err != nil {
+		return OperatorActionRow{}, err
+	}
+	return c.OperatorAction(ctx, action, now)
+}
+
+// PendingOperatorActions supplies the redacted Needs-you projection.
+func (c *Credentials) PendingOperatorActions(ctx context.Context, now int64) ([]OperatorActionRow, error) {
+	if err := c.operatorAccess(ctx); err != nil {
+		return nil, err
+	}
+	rows, err := c.db.QueryContext(ctx, `SELECT id, episode_id, kind, created_at, expires_at FROM operator_actions
+		WHERE state = 'requested' AND created_at <= ? AND expires_at > ? ORDER BY created_at, id LIMIT 16`, now, now)
+	if err != nil {
+		return nil, err
+	}
+	defer rows.Close()
+	var out []OperatorActionRow
+	for rows.Next() {
+		r := OperatorActionRow{State: "requested"}
+		if err := rows.Scan(&r.ID, &r.Episode, &r.Kind, &r.CreatedAt, &r.ExpiresAt); err != nil {
+			return nil, err
+		}
+		out = append(out, r)
+	}
+	return out, rows.Err()
+}
+
 // NextOperatorStart reads the oldest eligible start without assigning ownership.
 // The caller must take the canonical episode lock and compare the row again
 // before claiming it.
@@ -88,7 +129,7 @@ func (c *Credentials) RequestOperatorAction(ctx context.Context, request Operato
 		return OperatorActionRow{}, err
 	}
 	a := request.Action
-	if !operatorOperationID.MatchString(request.OperationID) ||
+	if !ValidOperatorOperationID(request.OperationID) ||
 		request.CreatedAt <= 0 || now < request.CreatedAt ||
 		request.ExpiresAt-request.CreatedAt > int64(approval.OperatorActionLifetime/time.Second) ||
 		c.lock == nil || a.Store != c.lock.dbPath {
diff --git a/host/store/operator_review.go b/host/store/operator_review.go
new file mode 100644
index 00000000..345b2caa
--- /dev/null
+++ b/host/store/operator_review.go
@@ -0,0 +1,71 @@
+package store
+
+import (
+	"crypto/sha256"
+	"encoding/json"
+	"fmt"
+	"path/filepath"
+	"strings"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// OperatorStartReview is a runner-published snapshot. The daemon fills only
+// the episode and literal task, using the same mint resolution as the runner.
+type OperatorStartReview struct {
+	Start                                 EpisodeStart
+	Grants                                []approval.Grant
+	Profile                               approval.ProfileRef
+	Store, Workspace, Sessions            string
+	Address, Build, Interpreter           string
+	Binary, BinaryDigest, Wall, Committer string
+}
+
+func (h OperatorHeartbeat) StartReview(preset string) (OperatorStartReview, error) {
+	var catalogue map[string]OperatorStartReview
+	if len(h.Catalogue) > 16*approval.MaxOperatorActionBytes || json.Unmarshal(h.Catalogue, &catalogue) != nil {
+		return OperatorStartReview{}, ErrOperatorRunner
+	}
+	v, ok := catalogue[preset]
+	if !ok || v.Start.Preset.Name != preset || v.Start.Agent != "codex" {
+		return OperatorStartReview{}, ErrOperatorRunner
+	}
+	return v, nil
+}
+
+// Action derives the only browser-start envelope from the published inputs.
+// The initially empty ID is assigned only by the request writer.
+func (v OperatorStartReview) Action(episode, task string, generation int64) (approval.OperatorAction, error) {
+	start := v.Start
+	start.Episode, start.Branch = episode, "world/"+episode
+	start.Task = &StartTask{
+		Path:   filepath.Join(v.Sessions, episode+".task.md"),
+		SHA256: fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(task))), Title: StartTaskTitle([]byte(task)),
+	}
+	in := OperatorMintInput{Attempt: OperatorMintAttempt{Catalogue: generation}, Start: start, Grants: v.Grants,
+		Profile: v.Profile, TokenPath: filepath.Join(v.Sessions, episode+".session")}
+	mint, err := in.Resolution()
+	if err != nil {
+		return approval.OperatorAction{}, err
+	}
+	worktree := filepath.Join(v.Workspace, episode)
+	facts := []approval.OperatorResolved{{Name: "launch_binary", Value: v.Binary},
+		{Name: "launch_binary_sha256", Value: v.BinaryDigest}, {Name: "launch_wall", Value: v.Wall}, {Name: "mcp_path"}}
+	facts = append(facts, mint...)
+	facts = append(facts, approval.OperatorResolved{Name: "served_address", Value: v.Address},
+		approval.OperatorResolved{Name: "served_build", Value: v.Build}, approval.OperatorResolved{Name: "served_interpreter", Value: v.Interpreter},
+		approval.OperatorResolved{Name: "worktree_path", Value: worktree})
+	var review strings.Builder
+	fmt.Fprintf(&review, "Episode %s from %s (%s)\nRepository %s\nRemote %s = %s\nBase %s @ %s\nNew worktree %s on %s\nHeadless %s %s, cap %s\nCommitter %s\nProfile %s %s\nSession %ds\n",
+		episode, start.Preset.Name, start.Preset.SHA256, start.Repo, start.Remote.Name, start.Remote.URL, start.Base.Branch, start.Base.SHA, worktree, start.Branch,
+		start.Agent, start.Model, v.Wall, v.Committer, v.Profile.Project, v.Profile.Digest, start.TTLS)
+	for _, check := range start.Checks {
+		fmt.Fprintf(&review, "Check %s (%s)\n", check.Name, check.Sandbox)
+	}
+	for _, grant := range v.Grants {
+		fmt.Fprintf(&review, "Grant %s %s budget %d\n", grant.Effect, grant.Scope, grant.Budget)
+	}
+	o := approval.OperatorOptions{Preset: start.Preset.Name, Headless: true}
+	return approval.OperatorAction{Schema: approval.OperatorActionSchema, Kind: "start", Episode: episode, Store: v.Store,
+		Argv: approval.OperatorArgv("start", episode, o), Options: o, Task: task, Review: review.String(), Resolution: facts}, nil
+}
diff --git a/host/store/start_task_title.go b/host/store/start_task_title.go
new file mode 100644
index 00000000..351dcb40
--- /dev/null
+++ b/host/store/start_task_title.go
@@ -0,0 +1,36 @@
+package store
+
+import "strings"
+
+// StartTaskTitle derives a short title for new starts; stored titles stay unchanged.
+func StartTaskTitle(data []byte) string {
+	title := ""
+	for _, line := range strings.Split(string(data), "\n") {
+		title = strings.TrimSpace(line)
+		if title != "" {
+			break
+		}
+	}
+	if strings.HasPrefix(title, "#") {
+		title = strings.TrimSpace(strings.TrimLeft(title, "#"))
+	} else {
+		if len(title) >= 5 && strings.EqualFold(title[:5], "Task:") {
+			title = strings.TrimSpace(title[5:])
+		}
+		for _, sep := range []string{". ", "; ", " ("} {
+			if i := strings.Index(title, sep); i >= 0 {
+				title = title[:i]
+			}
+		}
+	}
+	title = strings.Join(strings.Fields(title), " ")
+	runes := []rune(title)
+	if len(runes) > 72 {
+		cut := 72
+		if i := strings.LastIndex(string(runes[:72]), " "); i >= 0 {
+			cut = len([]rune(string(runes[:72])[:i]))
+		}
+		title = string(runes[:cut]) + "…"
+	}
+	return title
+}
diff --git a/host/store/start_task_title_test.go b/host/store/start_task_title_test.go
new file mode 100644
index 00000000..3f494e87
--- /dev/null
+++ b/host/store/start_task_title_test.go
@@ -0,0 +1,27 @@
+package store
+
+import (
+	"strings"
+	"testing"
+)
+
+func TestTaskTitle(t *testing.T) {
+	for _, tc := range []struct{ name, input, want string }{
+		{"455", "Task: make `world-publish`'s \"store is locked\" refusal say what to do when the World service holds the store; explain the recovery.", "make `world-publish`'s \"store is locked\" refusal say what to do when…"},
+		{"heading", "# Heading", "Heading"},
+		{"blanks", "\n \n##  A   heading\ttext\nbody", "A heading text"},
+		{"sentence", strings.Repeat("word ", 40) + ". Next.", strings.TrimSpace(strings.Repeat("word ", 14)) + "…"},
+		{"word", strings.Repeat("界", 200), strings.Repeat("界", 72) + "…"},
+		{"empty", "", ""},
+		{"whitespace", " \n\t", ""},
+		{"label", "Task:", ""},
+		{"sentence-end", "tAsK: Fix this. Then that", "Fix this"},
+		{"paren", "Fix this (details)", "Fix this"},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			if got := StartTaskTitle([]byte(tc.input)); got != tc.want {
+				t.Fatalf("%q, want %q", got, tc.want)
+			}
+		})
+	}
+}
diff --git a/host/surfacepage/approve_model.go b/host/surfacepage/approve_model.go
index 7ddc6e23..05899110 100644
--- a/host/surfacepage/approve_model.go
+++ b/host/surfacepage/approve_model.go
@@ -72,8 +72,10 @@ type RequestsJSON struct {
 	DecisionsTruncated bool              `json:"decisionsTruncated"`
 	// Confirmations are the pending verb requests a confirmation card confirms (row 243 V3 §2.5);
 	// below credentials schema 9, ConfirmationsUnavailable says why there are none.
-	Confirmations            []PendingConfirmation `json:"confirmations"`
-	ConfirmationsUnavailable string                `json:"confirmationsUnavailable"`
+	Confirmations              []PendingConfirmation   `json:"confirmations"`
+	ConfirmationsUnavailable   string                  `json:"confirmationsUnavailable"`
+	OperatorActions            []PendingOperatorAction `json:"operatorActions,omitempty"`
+	OperatorActionsUnavailable string                  `json:"operatorActionsUnavailable,omitempty"`
 }
 
 // PendingConfirmation is one pending verb request: a summary decoded from its stored bytes, or,
diff --git a/host/surfacepage/operator.go b/host/surfacepage/operator.go
new file mode 100644
index 00000000..336d294b
--- /dev/null
+++ b/host/surfacepage/operator.go
@@ -0,0 +1,11 @@
+package surfacepage
+
+// PendingOperatorAction is one Needs-you row. Only an anchored, same-origin
+// card route can appear in Href.
+type PendingOperatorAction struct {
+	ID        string `json:"id"`
+	Episode   string `json:"episode"`
+	Kind      string `json:"kind"`
+	ExpiresAt int64  `json:"expiresAt"`
+	Href      string `json:"href,omitempty"`
+}
````

### M3 — S7b follow-up

````diff M3
diff --git a/cmd/ailang-worldd/operator_catalogue_test.go b/cmd/ailang-worldd/operator_catalogue_test.go
index dbd3010b..d59a44c3 100644
--- a/cmd/ailang-worldd/operator_catalogue_test.go
+++ b/cmd/ailang-worldd/operator_catalogue_test.go
@@ -2,6 +2,7 @@ package main
 
 import (
 	"reflect"
+	"strings"
 	"testing"
 )
 
@@ -33,3 +34,12 @@ func testOperatorCataloguePublication(t *testing.T, f *operatorLaunchRig) {
 		t.Fatalf("published review does not round-trip to exact launch authority:\n%+v\n%+v", a.Resolution, r.action.Resolution)
 	}
 }
+
+func TestOperatorLiteralTaskArgv(t *testing.T) {
+	task := "# Task\n<script>alert(1)</script> ; $(literal) &"
+	prompt := launchPrompt(task, nil)
+	argv := launchArgv("codex", "fixture", prompt, t.TempDir(), "", "127.0.0.1:7644", "", "", nil)
+	if len(argv) < 2 || argv[0] != "codex" || argv[len(argv)-1] != prompt || !strings.HasPrefix(prompt, task) {
+		t.Fatalf("browser task did not remain one literal prompt argument: %q", argv)
+	}
+}
diff --git a/cmd/ailang-worldd/operator_launch.go b/cmd/ailang-worldd/operator_launch.go
index e3c17073..fe3ca974 100644
--- a/cmd/ailang-worldd/operator_launch.go
+++ b/cmd/ailang-worldd/operator_launch.go
@@ -330,18 +330,7 @@ func operatorLaunchFailure(err error) string {
 }
 
 func operatorLaunchProcessStatus(s store.OperatorLaunchSnapshot, read func(int) (string, bool)) string {
-	if s.NoSpawn {
-		return "spawn-refused"
-	}
-	if s.Phase == "reserved" {
-		return "uncertain"
-	}
-	if s.Receipt.Status == "launched" {
-		if birth, live := read(s.Receipt.PID); !live || birth != s.Receipt.Birth {
-			return "uncertain"
-		}
-	}
-	return s.Receipt.Status
+	return s.ProcessStatus(read)
 }
 
 func (e serviceEnv) operatorLaunchStatus(id string, out io.Writer) error {
@@ -546,3 +535,5 @@ func operatorLaunchExitStatus(proc context.Context, result agentrun.Result) stri
 	}
 	return "stopped"
 }
+
+func operatorProcessIdentity(pid int) (string, bool) { return agentrun.ProcessIdentity(pid) }
diff --git a/cmd/ailang-worldd/operator_process_darwin.go b/cmd/ailang-worldd/operator_process_darwin.go
index d5121a31..42647641 100644
--- a/cmd/ailang-worldd/operator_process_darwin.go
+++ b/cmd/ailang-worldd/operator_process_darwin.go
@@ -1,7 +1,6 @@
 package main
 
 import (
-	"fmt"
 	"os"
 
 	"golang.org/x/sys/unix"
@@ -19,14 +18,3 @@ func listWorldProcesses() ([]worldProcess, error) {
 	}
 	return out, nil
 }
-
-// operatorProcessIdentity reads kernel birth time and excludes a reused PID.
-func operatorProcessIdentity(pid int) (string, bool) {
-	row, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
-	if err != nil || int(row.Proc.P_pid) != pid {
-		return "", false
-	}
-	p := row.Proc
-	const zombie = 5 // SZOMB in Darwin <sys/proc.h>
-	return fmt.Sprintf("%d.%d", p.P_starttime.Sec, p.P_starttime.Usec), p.P_stat != zombie
-}
diff --git a/cmd/ailang-worldd/operator_process_linux.go b/cmd/ailang-worldd/operator_process_linux.go
deleted file mode 100644
index 4fe77293..00000000
--- a/cmd/ailang-worldd/operator_process_linux.go
+++ /dev/null
@@ -1,26 +0,0 @@
-package main
-
-import (
-	"errors"
-	"os"
-	"strconv"
-	"strings"
-)
-
-func listWorldProcesses() ([]worldProcess, error) {
-	return nil, errors.New("operator upgrade requires the macOS process inventory")
-}
-
-func operatorProcessIdentity(pid int) (string, bool) {
-	stat, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
-	boot, bootErr := os.ReadFile("/proc/sys/kernel/random/boot_id")
-	end := strings.LastIndexByte(string(stat), ')')
-	if err != nil || bootErr != nil || end < 0 {
-		return "", false
-	}
-	fields := strings.Fields(string(stat[end+1:]))
-	if len(fields) <= 19 {
-		return "", false
-	}
-	return strings.TrimSpace(string(boot)) + ":" + fields[19], fields[0] != "Z"
-}
diff --git a/cmd/ailang-worldd/operator_process_other.go b/cmd/ailang-worldd/operator_process_other.go
index ad91f103..7a30ebd9 100644
--- a/cmd/ailang-worldd/operator_process_other.go
+++ b/cmd/ailang-worldd/operator_process_other.go
@@ -1,4 +1,4 @@
-//go:build !darwin && !linux
+//go:build !darwin
 
 package main
 
@@ -7,5 +7,3 @@ import "errors"
 func listWorldProcesses() ([]worldProcess, error) {
 	return nil, errors.New("operator upgrade requires the macOS process inventory")
 }
-
-func operatorProcessIdentity(int) (string, bool) { return "", false }
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 21e7b126..86cf78f6 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -2121,3 +2121,23 @@ launch-status ACTION-ID`. There is no unlock command for `uncertain`. Preserve
 the action id, runner log, launch receipt, worktree and session records for
 attended reconciliation. Check the episode report for existing work. Only
 positive evidence can settle the reservation. Refresh and retry preserve it.
+
+#### Confirm with an enrolled key
+
+Verify the pending card's review and literal task. For SSHSIG, save
+`messageHex` as bytes with `xxd -r -p` and sign with
+`ssh-keygen -Y sign -n world-operator-action@sunholo -f <enrolled-key>`.
+For a passkey, use SHA-256 of those bytes as the WebAuthn challenge,
+`userVerification: "required"`, and the current origin's `rp` and `allow`.
+
+On the approval origin, send JSON to the card's
+`POST /surface/operator-actions/{id}/confirm` with `v: 1`, `digest` and
+one `sshsig` string or `webauthn` assertion in the approval-card format.
+The response names the enrolled key and action link. Exact replay returns
+`already-confirmed`. A ready runner may then claim the one approved action.
+
+Run the fake-key example with the repository fixture keys and temporary stores:
+
+```sh
+AILANG_BIN="$HOME/.pinned-ailang/ailang" go test -count=1 ./host/daemon/ -run '^TestOperatorConfirmHTTP$'
+```
diff --git a/host/agentrun/process_darwin.go b/host/agentrun/process_darwin.go
new file mode 100644
index 00000000..042a176e
--- /dev/null
+++ b/host/agentrun/process_darwin.go
@@ -0,0 +1,18 @@
+package agentrun
+
+import (
+	"fmt"
+
+	"golang.org/x/sys/unix"
+)
+
+// ProcessIdentity reads kernel birth time and excludes a reused PID.
+func ProcessIdentity(pid int) (string, bool) {
+	row, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
+	if err != nil || int(row.Proc.P_pid) != pid {
+		return "", false
+	}
+	p := row.Proc
+	const zombie = 5 // SZOMB in Darwin <sys/proc.h>
+	return fmt.Sprintf("%d.%d", p.P_starttime.Sec, p.P_starttime.Usec), p.P_stat != zombie
+}
diff --git a/host/agentrun/process_linux.go b/host/agentrun/process_linux.go
new file mode 100644
index 00000000..1ab62ad6
--- /dev/null
+++ b/host/agentrun/process_linux.go
@@ -0,0 +1,21 @@
+package agentrun
+
+import (
+	"os"
+	"strconv"
+	"strings"
+)
+
+func ProcessIdentity(pid int) (string, bool) {
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
diff --git a/host/agentrun/process_other.go b/host/agentrun/process_other.go
new file mode 100644
index 00000000..255b2d92
--- /dev/null
+++ b/host/agentrun/process_other.go
@@ -0,0 +1,5 @@
+//go:build !darwin && !linux
+
+package agentrun
+
+func ProcessIdentity(int) (string, bool) { return "", false }
diff --git a/host/daemon/surface_operator.go b/host/daemon/surface_operator.go
index d0b6280f..682f7217 100644
--- a/host/daemon/surface_operator.go
+++ b/host/daemon/surface_operator.go
@@ -108,12 +108,17 @@ func (d *Daemon) handleOperatorRoutes(w http.ResponseWriter, r *http.Request, or
 				d.handleOperatorRead(w, r, id, hasSuffix)
 				return
 			}
+		case suffix == "confirm":
+			if r.Method == http.MethodPost {
+				d.handleOperatorConfirm(w, r, id, origins)
+				return
+			}
 		default:
 			writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no operator action route")
 			return
 		}
 	}
-	writeApprovalRefusal(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "use GET for a card/status or POST to request an action")
+	writeApprovalRefusal(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "use GET for a card/status or POST to request or confirm an action")
 }
 
 // readOperatorPOST is reached after the Host and method checks, before any body read.
diff --git a/host/daemon/surface_operator_confirm.go b/host/daemon/surface_operator_confirm.go
new file mode 100644
index 00000000..cfbb8f45
--- /dev/null
+++ b/host/daemon/surface_operator_confirm.go
@@ -0,0 +1,146 @@
+package daemon
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"encoding/hex"
+	"encoding/json"
+	"errors"
+	"net/http"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+func operatorCeremony(ctx context.Context, db *store.Credentials, row store.OperatorActionRow, href string) (string, error) {
+	active, err := db.ActivePasskeys(ctx)
+	if err != nil {
+		return "", err
+	}
+	var keys []surfacepage.CardPasskey
+	for _, e := range active {
+		keys = append(keys, surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})
+	}
+	m, err := approval.OperatorActionMessage([32]byte(row.Digest), [32]byte(row.Challenge), row.ExpiresAt)
+	if err != nil {
+		return "", err
+	}
+	request := surfacepage.RequestCard(surfacepage.StoredRequest{}, keys)
+	wire, err := json.MarshalIndent(struct {
+		Post      string                     `json:"post"`
+		Digest    string                     `json:"digest"`
+		Message   string                     `json:"messageHex"`
+		Namespace string                     `json:"namespace"`
+		Passkeys  []surfacepage.CardPasskeys `json:"passkeys"`
+	}{href, hex.EncodeToString(row.Digest), hex.EncodeToString(m), approval.OperatorActionNamespace, request.Passkeys}, "", "  ")
+	return string(wire), err
+}
+
+type operatorConfirmationAnswer struct {
+	operatorStatus
+	Result     string `json:"result"`
+	ApprovedBy string `json:"approvedBy"`
+}
+
+func (d *Daemon) handleOperatorConfirm(w http.ResponseWriter, r *http.Request, id string, origins []string) {
+	raw, ok := readOperatorPOST(w, r, origins, surfacepage.MaxApproveBodyBytes)
+	if !ok {
+		return
+	}
+	var fields map[string]json.RawMessage
+	if approval.DecodeOperatorObject(raw, &fields, "v", "digest", "webauthn", "sshsig") != nil || len(fields) != 3 {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send a v1 digest and one key artifact")
+		return
+	}
+	body, err := surfacepage.ParseApproveBody(raw)
+	if err != nil {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send a v1 digest and one key artifact")
+		return
+	}
+	artifact := []byte(body.WebAuthn)
+	if body.SSHSIG != nil {
+		artifact = []byte(*body.SSHSIG)
+	}
+
+	d.humanMu.Lock()
+	defer d.humanMu.Unlock()
+	ctx, cancel := context.WithTimeout(r.Context(), d.commitBudget)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	defer ro.Close()
+	now := d.approveNow()
+	row, err := ro.OperatorAction(ctx, id, now)
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	a, err := approval.DecodeOperatorAction(row.Action)
+	digest := sha256.Sum256(row.Action)
+	if err != nil || a.Kind != "start" || !a.Options.Headless || a.ID != id || a.Episode != row.Episode || a.Kind != row.Kind ||
+		!bytes.Equal(row.Digest, digest[:]) || body.Digest != hex.EncodeToString(row.Digest) || len(row.Challenge) != 32 {
+		d.operatorError(w, store.ErrOperatorConflict)
+		return
+	}
+	if row.ApprovedBy != "" {
+		if !bytes.Equal(row.Artifact, artifact) {
+			d.operatorError(w, store.ErrOperatorConflict)
+			return
+		}
+		writeApprovalJSON(w, http.StatusOK, operatorConfirmationAnswer{operatorStatusFor(row, now), "already-confirmed", row.ApprovedBy})
+		return
+	}
+	if row.State == "expired" || now < row.CreatedAt {
+		operatorConfirmExpired(w)
+		return
+	}
+	if row.State != "requested" {
+		d.operatorError(w, store.ErrOperatorConflict)
+		return
+	}
+	proof, err := authority.JudgeOperatorAction(ctx, ro, authority.OperatorJudgeInput{
+		Action: a, Digest: [32]byte(row.Digest), BytesDigest: digest, Challenge: [32]byte(row.Challenge),
+		CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Now: now,
+	}, artifact)
+	var judge *authority.JudgeError
+	switch {
+	case errors.Is(err, approval.ErrOperatorLifetime):
+		operatorConfirmExpired(w)
+		return
+	case errors.As(err, &judge) && judge.Step == authority.JudgeStepStore:
+		d.operatorError(w, err)
+		return
+	case err != nil:
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the enrolled key did not confirm this exact action")
+		return
+	}
+
+	wr, err := d.openConfirmationWriter()
+	if err != nil {
+		d.operatorError(w, err)
+		return
+	}
+	defer wr.Close()
+	now = d.approveNow()
+	row, err = wr.ConfirmOperatorAction(ctx, store.OperatorCAS{ID: id, Digest: digest, Generation: row.Generation}, artifact, proof, now)
+	switch {
+	case errors.Is(err, store.ErrOperatorExpired):
+		operatorConfirmExpired(w)
+	case isApproverRace(err), isCounterRefusal(err):
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the enrolled key changed during confirmation")
+	case err != nil:
+		d.operatorError(w, err)
+	default:
+		writeApprovalJSON(w, http.StatusOK, operatorConfirmationAnswer{operatorStatusFor(row, now), "confirmed", row.ApprovedBy})
+	}
+}
+
+func operatorConfirmExpired(w http.ResponseWriter) {
+	writeApprovalRefusal(w, http.StatusConflict, "Expired", "the action is outside its signed lifetime; review again")
+}
diff --git a/host/daemon/surface_operator_confirm_test.go b/host/daemon/surface_operator_confirm_test.go
new file mode 100644
index 00000000..36042e93
--- /dev/null
+++ b/host/daemon/surface_operator_confirm_test.go
@@ -0,0 +1,230 @@
+package daemon
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"database/sql"
+	"encoding/hex"
+	"encoding/json"
+	"net/http"
+	"net/http/httptest"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+func (f *operatorHTTP) row(id string) store.OperatorActionRow {
+	f.t.Helper()
+	ro, err := store.OpenCredentialsReadOnly(f.d.credentialsPath)
+	if err != nil {
+		f.t.Fatal(err)
+	}
+	defer ro.Close()
+	row, err := ro.OperatorAction(boundedTestContext(f.t), id, 0)
+	if err != nil {
+		f.t.Fatal(err)
+	}
+	return row
+}
+
+func opMessage(t *testing.T, digest, nonce [32]byte, expires int64) []byte {
+	t.Helper()
+	m, err := approval.OperatorActionMessage(digest, nonce, expires)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return m
+}
+
+func TestOperatorConfirmHTTP(t *testing.T) {
+	f := newOperatorHTTP(t)
+	key := approvaltest.NewPasskey(approval.AlgEdDSA, 43)
+	signer := approvaltest.NewSigner(45)
+	f.enrolSSH("fixture-ssh", signer)
+	f.enrolPasskey("fixture-passkey", key, opOrigin)
+	s := f.request(opRequest(1))
+	row := f.row(s.ID)
+	digest, nonce := [32]byte(row.Digest), [32]byte(row.Challenge)
+	d := hex.EncodeToString(row.Digest)
+	m := opMessage(t, digest, nonce, row.ExpiresAt)
+	challenge := sha256.Sum256(m)
+	sig := func(key *approvaltest.Signer, namespace string, msg []byte) []byte {
+		return sshBody(d, key.Sign(namespace, msg))
+	}
+	signed := sig(signer, approval.OperatorActionNamespace, m)
+	passkey := approveBody(d, key.Assert(opOrigin, "localhost", challenge, 2))
+	before := f.evidence()
+	post := func(body []byte) *httptest.ResponseRecorder { return f.do("POST", s.Href+"/confirm", body, nil) }
+	refuse := func(b []byte, status int, class string) { f.refuse(before, post(b), status, class) }
+	f.sub("Transport", func(t *testing.T) {
+		for _, tc := range []struct {
+			key, value, class string
+			code              int
+		}{
+			{"Origin", "", "CrossOrigin", 403}, {"Origin", "null", "CrossOrigin", 403}, {"Origin", opOrigin + ".invalid", "CrossOrigin", 403},
+			{"Sec-Fetch-Site", "same-site", "CrossOrigin", 403}, {"Content-Type", "text/plain", "UnsupportedMediaType", 415},
+			{"Authorization", "Bearer y", "Forbidden", 403}, {"Cookie", "y=y", "Forbidden", 403},
+		} {
+			r := f.do("POST", s.Href+"/confirm", signed, func(r *http.Request) {
+				if tc.value == "" {
+					r.Header.Del(tc.key)
+				} else {
+					r.Header.Set(tc.key, tc.value)
+				}
+			})
+			f.refuse(before, r, tc.code, tc.class)
+		}
+		f.refuse(before, f.do("POST", s.Href+"/confirm", signed, func(r *http.Request) { r.Host = "other.invalid:7644" }), 421, "MisdirectedRequest")
+		for _, method := range []string{"GET", "OPTIONS", "PUT"} {
+			f.refuse(before, f.do(method, s.Href+"/confirm", signed, nil), 405, "MethodNotAllowed")
+		}
+		refuse(bytes.Repeat([]byte(" "), (20<<10)+1), 413, "TooLarge")
+	})
+	f.sub("Proof", func(t *testing.T) {
+		for _, body := range [][]byte{[]byte("y"), opJSON(map[string]any{"v": 1, "digest": d}),
+			opJSON(map[string]any{"v": 1, "digest": d, "y": true}), append([]byte(`{"v":1,`), passkey[1:]...), append([]byte(`{"sshsig":null,`), passkey[1:]...)} {
+			refuse(body, 400, "BadRequest")
+		}
+		refuse(sshBody(d, []byte("malformed")), 422, "NotApproved")
+		other := nonce
+		other[0] ^= 1
+		refuse(sig(signer, approval.OperatorActionNamespace, opMessage(t, digest, other, row.ExpiresAt)), 422, "NotApproved")
+	})
+	f.sub("Authenticity", func(t *testing.T) {
+		refuse(sig(approvaltest.NewSigner(47), approval.OperatorActionNamespace, m), 422, "NotApproved")
+		refuse(sig(signer, approval.VerbNamespace, m), 422, "NotApproved")
+		old, err := approval.VerbMessage(digest, nonce, row.ExpiresAt)
+		if err != nil {
+			t.Fatal(err)
+		}
+		refuse(sig(signer, approval.VerbNamespace, old), 422, "NotApproved")
+		refuse(sig(signer, approval.Namespace, m), 422, "NotApproved")
+		for _, proof := range [][]byte{
+			key.Assert(opOrigin, "localhost", challenge, 2, approvaltest.Edit{BadSignature: true}),
+			key.Assert(opOrigin, "localhost", challenge, 2, approvaltest.Edit{Clear: 0x04}),
+			key.Assert("http://localhost:7645", "localhost", challenge, 2),
+			key.Assert(opOrigin, "elsewhere.invalid", challenge, 2),
+			approvaltest.NewPasskey(approval.AlgEdDSA, 44).Assert(opOrigin, "localhost", challenge, 2),
+		} {
+			refuse(approveBody(d, proof), 422, "NotApproved")
+		}
+	})
+	f.sub("Binding", func(t *testing.T) {
+		var a map[string]json.RawMessage
+		if err := json.Unmarshal(row.Action, &a); err != nil {
+			t.Fatal(err)
+		}
+		for field, value := range map[string]any{"id": "oa-" + strings.Repeat("a", 32), "episode": "other", "kind": "finish",
+			"store": "/other/credentials.db", "argv": []string{"different"}, "options": approval.OperatorOptions{Preset: "other", Headless: true}, "task": "changed task",
+		} {
+			match := string(opJSON(field)) + ":" + string(a[field])
+			replace := string(opJSON(field)) + ":" + string(opJSON(value))
+			alternate := bytes.Replace(row.Action, []byte(match), []byte(replace), 1)
+			if bytes.Equal(row.Action, alternate) {
+				t.Fatalf("binding control did not change %s", field)
+			}
+			refuse(sig(signer, approval.OperatorActionNamespace, opMessage(t, sha256.Sum256(alternate), nonce, row.ExpiresAt)), 422, "NotApproved")
+		}
+	})
+	f.sub("Expiry", func(t *testing.T) {
+		for _, now := range []int64{row.CreatedAt - 1, row.ExpiresAt, row.ExpiresAt + 1} {
+			f.now.Store(now)
+			refuse(signed, 409, "Expired")
+		}
+		f.now.Store(apT0)
+		db, err := sql.Open("sqlite", f.d.credentialsPath)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer db.Close()
+		change := func(exp int64) {
+			if _, err := db.ExecContext(boundedTestContext(t), "UPDATE operator_actions SET expires_at=? WHERE id=?", exp, s.ID); err != nil {
+				t.Fatal(err)
+			}
+			before = f.evidence()
+		}
+		change(row.ExpiresAt + 1)
+		refuse(sig(signer, approval.OperatorActionNamespace, opMessage(t, digest, nonce, row.ExpiresAt+1)), 409, "Expired")
+		change(row.ExpiresAt)
+	})
+	f.sub("RevocationCounter", func(t *testing.T) {
+		refuse(approveBody(d, key.Assert(opOrigin, "localhost", challenge, 1)), 422, "NotApproved")
+		revoked := approvaltest.NewSigner(46)
+		f.enrolSSH("revoked-ssh", revoked)
+		f.write(func(ctx context.Context, c *store.Credentials) error {
+			return c.RevokeApprover(ctx, "revoked-ssh", apT0)
+		})
+		before = f.evidence()
+		refuse(sig(revoked, approval.OperatorActionNamespace, m), 422, "NotApproved")
+	})
+	f.sub("UnsupportedKinds", func(t *testing.T) {
+		for i, kind := range []string{"start", "finish"} {
+			a, err := approval.DecodeOperatorAction(row.Action)
+			if err != nil {
+				t.Fatal(err)
+			}
+			a.ID, a.Episode, a.Kind = "oa-"+opRequest(i+30).OperationID, "other-"+kind, kind
+			a.Options.Headless = false
+			if kind == "finish" {
+				a.Task = ""
+			}
+			a.Argv = approval.OperatorArgv(kind, a.Episode, a.Options)
+			var stored store.OperatorActionRow
+			f.write(func(ctx context.Context, c *store.Credentials) error {
+				var err error
+				stored, err = c.RequestOperatorAction(ctx, store.OperatorRequest{
+					Action: a, OperationID: opRequest(i + 30).OperationID, Challenge: nonce, CreatedAt: apT0, ExpiresAt: apT0 + 600,
+				}, apT0)
+				return err
+			})
+			before = f.evidence()
+			body := sshBody(hex.EncodeToString(stored.Digest), signer.Sign(approval.OperatorActionNamespace,
+				opMessage(t, [32]byte(stored.Digest), nonce, stored.ExpiresAt)))
+			f.refuse(before, f.do("POST", operatorPath+"/"+a.ID+"/confirm", body, nil), 409, "Conflict")
+		}
+	})
+	f.sub("MissingAndWriter", func(t *testing.T) {
+		f.refuse(before, f.do("POST", operatorPath+"/oa-"+strings.Repeat("0", 32)+"/confirm", signed, nil), 404, "NotFound")
+		held, err := store.OpenCredentials(f.d.credentialsPath)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer held.Close()
+		refuse(signed, 503, "Busy")
+	})
+	card := f.do("GET", s.Href, nil, nil)
+	if card.Code != 200 || !strings.Contains(card.Body.String(), hex.EncodeToString(m)) || !strings.Contains(card.Body.String(), approval.OperatorActionNamespace) {
+		t.Fatalf("card did not offer the exact ceremony: %d %s", card.Code, card.Body)
+	}
+	f.now.Store(row.ExpiresAt - 1)
+	expect := func(r *httptest.ResponseRecorder, id, result, key string) {
+		t.Helper()
+		var out operatorConfirmationAnswer
+		if err := json.Unmarshal(r.Body.Bytes(), &out); err != nil || r.Code != 200 || out.Result != result || out.ApprovedBy != key ||
+			out.ID != id || out.State != "confirmed" {
+			t.Fatalf("confirmation %s: %d %s (%v)", result, r.Code, r.Body, err)
+		}
+	}
+	expect(post(passkey), s.ID, "confirmed", "webauthn:fixture-passkey")
+	before = f.evidence()
+	expect(post(passkey), s.ID, "already-confirmed", "webauthn:fixture-passkey")
+	if f.evidence() != before || !strings.Contains(f.credState(), "[fixture-passkey 2") {
+		t.Fatal("replay changed durable evidence or advanced the counter twice")
+	}
+	refuse(signed, 409, "Conflict")
+	f.now.Store(apT0)
+	second := f.request(opRequest(2))
+	r2 := f.row(second.ID)
+	body := sshBody(hex.EncodeToString(r2.Digest), signer.Sign(approval.OperatorActionNamespace,
+		opMessage(t, [32]byte(r2.Digest), [32]byte(r2.Challenge), r2.ExpiresAt)))
+	expect(f.do("POST", second.Href+"/confirm", body, nil), second.ID, "confirmed", "ssh:fixture-ssh")
+	for _, text := range []string{f.do("GET", second.StatusHref, nil, nil).Body.String(), f.do("GET", "/surface/requests.json", nil, nil).Body.String(), f.log.String()} {
+		if strings.Contains(text, "BEGIN SSH SIGNATURE") || strings.Contains(text, "alert(1)") || strings.Contains(text, hex.EncodeToString(m)) {
+			t.Fatalf("private task or artifact reached a routine read or log: %s", text)
+		}
+	}
+}
diff --git a/host/daemon/surface_operator_read.go b/host/daemon/surface_operator_read.go
index 0bdf31ca..deb28875 100644
--- a/host/daemon/surface_operator_read.go
+++ b/host/daemon/surface_operator_read.go
@@ -7,6 +7,7 @@ import (
 	"html/template"
 	"net/http"
 
+	"github.com/sunholo-data/ailang-world/host/agentrun"
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/store"
 	"github.com/sunholo-data/ailang-world/host/verbs"
@@ -57,6 +58,9 @@ var operatorCard = template.Must(template.New("operator action").Parse(`<!doctyp
 <html lang="en"><meta charset="utf-8"><title>Operator action</title>
 <body><h1>{{.Status.ID}}</h1><p>{{.Status.State}}</p>{{if .Status.Reason}}<p>{{.Status.Reason}}</p>{{end}}
 {{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre><h2>Task</h2><pre>{{.Task}}</pre>{{end}}
+{{if .Ceremony}}<h2>Confirm this action</h2>
+<p>Verify the exact review and task. Submit a passkey assertion or SSH signature over this action using the ceremony below.</p>
+<pre id="operator-ceremony">{{.Ceremony}}</pre>{{end}}
 {{if .Status.StatusHref}}<p><a href="{{.Status.StatusHref}}">Action status (JSON)</a></p>{{end}}
 {{if .Status.ReportHref}}<p><a href="{{.Status.ReportHref}}">Episode report</a></p>{{end}}
 </body></html>`))
@@ -85,14 +89,7 @@ func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id s
 	if status.Href != "" && row.Kind == "start" && (row.State == "running" || row.State == "succeeded" || row.State == "uncertain" || row.State == "failed") {
 		launch, err := ro.ReadOperatorLaunch(ctx, id)
 		if err == nil {
-			switch {
-			case launch.NoSpawn:
-				status.State = "spawn-refused"
-			case launch.Phase == "reserved" || launch.Receipt.Status == "launched":
-				status.State = "uncertain"
-			default:
-				status.State = launch.Receipt.Status
-			}
+			status.State = launch.ProcessStatus(agentrun.ProcessIdentity)
 			if status.ReportHref == "" {
 				status.ReportHref = "/v1/episodes/" + row.Episode + "/report"
 			}
@@ -109,8 +106,8 @@ func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id s
 		return
 	}
 	card := struct {
-		Status       operatorStatus
-		Review, Task string
+		Status                 operatorStatus
+		Review, Task, Ceremony string
 	}{Status: status}
 	if len(row.Action) != 0 && status.Href != "" {
 		a, err := approval.DecodeOperatorAction(row.Action)
@@ -120,6 +117,13 @@ func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id s
 			return
 		}
 		card.Review, card.Task = a.Review, a.Task
+		if row.State == "requested" && len(row.Challenge) == 32 && row.CreatedAt <= now {
+			card.Ceremony, err = operatorCeremony(ctx, ro, row, status.Href+"/confirm")
+			if err != nil {
+				d.operatorError(w, err)
+				return
+			}
+		}
 	}
 	setSurfaceHeaders(w)
 	_ = operatorCard.Execute(w, card)
diff --git a/host/daemon/surface_operator_test.go b/host/daemon/surface_operator_test.go
index 7350c645..9ffc0ce5 100644
--- a/host/daemon/surface_operator_test.go
+++ b/host/daemon/surface_operator_test.go
@@ -258,7 +258,6 @@ func TestOperatorAdmissionHTTP(t *testing.T) {
 			refuse(path, "GET", nil, 400, "BadRequest", nil)
 		}
 		refuse(operatorPath+id+"/status.json/x", "GET", nil, 404, "NotFound", nil)
-		refuse(operatorPath+id+"/confirm", "POST", valid, 404, "NotFound", nil)
 	})
 	f.sub("StrictJSONAndLimits", func(t *testing.T) {
 		invalid := [][]byte{nil, []byte("[]"), []byte(`{"argv":"start"}`), []byte(`{"argv":[42]}`), append(bytes.Clone(valid), []byte(" {}")...),
diff --git a/host/store/operator_launch.go b/host/store/operator_launch.go
index 206672dc..84f62df2 100644
--- a/host/store/operator_launch.go
+++ b/host/store/operator_launch.go
@@ -208,6 +208,22 @@ func readOperatorLaunch(ctx context.Context, tx *sql.Tx, id string) (OperatorLau
 	return s, row, ErrOperatorConflict
 }
 
+// ProcessStatus classifies the persisted receipt with a current identity read.
+func (s OperatorLaunchSnapshot) ProcessStatus(read func(int) (string, bool)) string {
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
 // RecordOperatorLaunch compares the exact read-back and never replaces a
 // process identity. Recording an admitted child's outcome ignores expiry.
 func (c *Credentials) RecordOperatorLaunch(ctx context.Context, reviewed OperatorLaunchSnapshot, r OperatorLaunchReceipt, history *EpisodeLaunch, now int64) error {
````

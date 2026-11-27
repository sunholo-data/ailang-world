# w-world-address-bar-verbs-v4 — S7b sprint plan

## 0. Scope and outcome

Row 243, iteration 446, lane `codex:gpt-6.2-sol`. Design r1 has quorum r2 PROCEED 3/3, with all seats present. D-WORLD-119 and D-WORLD-120 are A, attended at `82f05004`. S7 landed as #700, merge `6836c9a9`. M1 activates the first browser POST that can make one bounded headless start claimable. Route registration, the ceremony, the reviewed refusals and live process birth readback ship together. The installed runner retains the claim, mint, worktree, process and recovery lifetimes.

M1 measures **+1191/−72 = 1263** cumulatively, **1263** unfiltered product lines. The cap is 1,400. The two plan files are excluded from executor snapshots and product numstat. The stage is a full archive of the base. The planner used `git diff --no-index -M --numstat base-tree st-M1` and checked every changed path under `cmd host docs app`. The controller checks `git diff -M --numstat 88fea1e9 HEAD -- cmd host docs app` after commit.

Scope remains one headless start with existing schema 10, existing enrollment, S2/S3 proof/CAS, and S6 launch. S8 owns app bar/standalone browser UI, remembered status on disconnect and Chrome/mobile acceptance. S9–S11 own finish and F1, and S12 owns interactive hand-off. The reviewed S8–S12 budgets are respectively 1,161, 1,215, 1,094, 1,215 and 945 calibrated lines. There is no additional S7b product milestone parked behind the route. Design §6 places browser/OS identity in Go and credentials transactions in the store.

## 1. Measurements

### 1.1 Base and provenance

`B=88fea1e9a8422b12a059eb1ecfc9c10c5b7b4e7d`, with a clean product worktree before planning. `git diff --stat 6836c9a9 88fea1e9 -- cmd host docs app` returned no entries. `git show --stat eb762a2e` confirmed +193/−33 across six files, 226 lines. S7 M1 was `80af7f36`, M2 `e27c93ef`, followed by `eb762a2e`. Read the V4 design §§3.1–3.3, 4.1–4.2, 5.2, 6 and 8, S7 plan and JSON including Appendix P M3, iteration-445 log/readings and both judges. The 445 state is read-only.

| Base origin | Observed at 88fea1e9 |
|---|---|
| `host/daemon/surface_operator.go:73–118` | The action family is wrapped in `approvalRoute` before ServeMux path cleaning. Only exact request and GET card/status dispatch. Confirm suffix falls through to 404. |
| `host/daemon/surface_approve.go:76–86` and `surface_operator.go:120–145` | Existing Host check runs first. Exact origin count/membership, fetch-site grammar when present, JSON media type, authorization/cookie, bounded reader and shared refusal writer are reusable. |
| `host/daemon/surface_confirm.go:97–207` | V3 supplies the bounded request context, human writer lock and fail-fast writer open. V4 reuses the writer seam and uses the V4 domain/signature law. |
| `host/store/operator_cas.go:42`, `host/authority/judge_operator.go:61,66` | Store and judge exist. Only the intervening approval-law call is production-wired. |
| `cmd/ailang-worldd/operator_preparation.go:115`, `operator_runner.go:96`, `operator_launch.go:216` | Installed preparation owns claim, the runner loop owns the handoff and the supervisor starts the worker. |
| `scripts/verify_go.sh:161` | Native whole-package race bound is 900 s. Merge #700 run 46046770928 measured cmd 284.9 s, daemon 363.2 s, verifygate 324.9 s, broker 279.0 s. PR run 46044386173 measured 503.351 s and 513.203 s for cmd and daemon respectively. |

### 1.2 Guard and route inventory

`selectors.sh` and `gates.sh` contain the exact invocations. Inventory tests require an actual `--- PASS:` for each selected test. The existing verb and requests goldens have byte-identical fixtures in the base and stage. The actual confirmation route is also inventoried by a valid positive request, a near-match 404, wrong-method 405, malformed 400, and inherited headers on the card. The verifygate scan found no independent operator URI table. It found the subprocess, boundary/fork and CI manifest/timeout test families. See `inventory-scan.txt` and the whole verifygate leg.

| Package / selector | Reading |
|---|---|
| `host/daemon` `TestDaemonDependencyAllowlist`, `TestDaemonCredentialWriterSites`, `TestFrozenV1RouteTableMatchesMux`, `TestApprovalRoutesRegistered`, `TestApprovalRoutesRefuseForeignHost`, `TestConfirmRoutesRegistered`, `TestDecisionRoutesRegistered`, `TestSurfaceVerbRouteRegistered`, `TestSurfaceVerbIsReadOnly` | PASS (`$S/local-M1/inv-daemon.raw`) |
| `host/daemon` `TestSurfaceVerbBoundsAndErrors`, `TestSurfaceVerbWhyIsReadOnly`, `TestSurfaceVerbWhyOutcomes`, `TestWhyResult` | PASS (four tests in `$S/verb-parity-final.raw`) |
| `host/daemon` `TestSurfaceSecurityHeaders`, `TestSurfaceHomeSummary`, `TestReadRoutesAndPayloadGate`, `TestGETRoutesRejectOtherMethods`, `TestProjectionRoutes_MountedAndNotProtected`, `TestAgentCardRoute_GETOnly`, `TestSurfaceAppRoutes`, `TestSurfaceScriptRoute`, `TestSurfaceModelRoutes`, `TestApprovePageRoutes`, `TestPasskeyPageRoutes`, `TestSurfaceFontRoute`, `TestWorkbenchRouteIsReadOnly`, `TestWorkbenchSecurityHeaders` | PASS (`$S/local-M1/inv-daemon.raw`) |
| `host/daemon` `TestSurfaceVerbOutcomes`, `TestRequestsJSONConfirmationsGolden`; `host/surfacepage` `TestSurfaceAppVerbs`, `TestNeedsYouConfirmationEmbeddedNewline`, `TestNeedsYouConfirmationsNonArray`, `TestVerbDecidePrefixBoundary` | PASS (six named golden tests in `$S/golden-and-compile/golden.raw`) |
| `host/store` `TestProductionContextRoots`, `TestStoreSurfaceTakesContext`, `TestCredentialsMethodsRequireADeadline`, `TestStoreDeadlineGuardDynamic` | PASS (`$S/local-M1/inv-store.raw`). The new timeout is a child of `r.Context()`. No new production context root. |
| `host/boundary` `TestWorldBoundaryDependencyAllowlist`, `TestBoundaryASTWriteGuard` | PASS (`$S/local-M1/inv-boundary.raw`) |
| `cmd/ailang-worldd` `TestCLICredentialWriterSites`, `TestSessionRootBudgets` | PASS (`$S/local-M1/inv-cli.raw`) |
| `host/verbs` `TestOperatorNeverDispatch`, verb golden reader `TestOperatorOtherVerbParity` | PASS (both tests in `$S/inventory-extra/verbs.raw`) |
| `host/daemon` `TestSurfacePage`, `TestSurfaceIndex`, `TestSurfaceFollowRoutes`, `TestSurfaceFollowRoutesLeakNothing`, `TestSurfaceFileReadsOnlyPublishedLandedOutputs`, `TestDecidePageAssets`, `TestConfirmPageAssetsServed`, `TestConfirmationCardServed`, `TestApprovePOSTCrossOrigin`, `TestApprovePOSTNoSessionNoCookie`, `TestConfirmPOSTGuards`, `TestWorkbenchLiveScriptRoute`, `TestWorkbenchErrorPageInert`, `TestWorkbenchStillAllRegions`, `TestWorkbenchTimelineBound`, `TestAGUIRejectsGET`, `TestAGUIInputErrors` | PASS (17 tests in `$S/inventory-extra/headers.raw`) |
| `host/daemon` `TestDaemonReadDeadline`, `TestWorkbenchReadDeadline`, `TestReadCtxCancelledAfterHandler`, `TestInternalErrorsAreSanitized` | PASS (four tests in `$S/inventory-extra/read.raw`) |
| `host/verifygate` `TestSubprocessSinkGate`, `TestExecGateArmsAreLoadBearing`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestVerifygateTestWritesAreForkLocked`, Linux-only `TestKernelRefusesExecOfWriterOpenFile` | PASS (`$S/local-M1/inv-verifygate.raw`). Linux runtime control is skipped on Darwin and remains a CI arm. |
| `host/broker` `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` | UNINFORMATIVE UNDER SANDBOX. Controller runs the outside gate.. The fake-launch test reuses the existing child and syncBuffer. No new exec site. |

### 1.3 Known-red baseline and race budgets

Controller pristine evidence is `$STATE/base`, accepted after `tail -n 1 "$STATE/baseline.done"` equals `ALL_DONE for $B`. It records store, authority and agentrun plain/race, verifygate/verbs/broker/cmd plain, and daemon plain/race. The added whole cmd race base is `$S/outside-pristine/cmd-race.attributable`; see §7. Every comparison records TEXT first, scans for truncation, exercises the parser control and compares test identities in `*.fails`. A nonzero rc with an empty fail set is an instrument failure. Comparison and time receipts are in §4. Known baseline failures are left with their existing owners.

Local real-binary baseline controls include `TestCLIAgainstSeToolsDaemon`, `TestCLINoGenesisIsNoWorldHead`, the eight `TestDemoSeed*` and `TestMarketingStillRoutes`, and daemon `TestSeToolsDescriptorLiveExactEight`, `TestSeToolsMCPVisibilityFollowsGrants`, `TestSeToolsA2AEndToEnd`. Known load shapes include world statusline, exec hosts, episode revoke, worktree/file publication, seed commit/preflight/old-copy, broker protected timeouts/creation, SE Tools fake benchmark, A2A deadline, registry phase, verifygate mission integration and 30 s git fixtures. `TestOperatorServiceStartup/gate/*` requires outside. Use the measured fail sets and a same-selector base re-run for any distinct new failure; these names do not excuse a regression.

### 1.4 Premises found false and caller census

`$S/rebase.raw` records applying the extracted 517-line S7 M3 to a fresh B archive with `/usr/bin/patch -p1 -E`. The only reject was the second hunk of `cmd/ailang-worldd/operator_catalogue_test.go` at the old end of `testOperatorCataloguePublication`. The verbose diagnostic observed hunk #1 at line 6 with fuzz 2 and offset +4, and hunk #2 rejected at line 38. The landed BuildDigest, Agent and RemoteCredentials subtests occupy the former EOF. Keep the file byte-identical to B, dropping the prototype's added import, and place the literal argv assertion in the real `TestOperatorStartBrowser` fixture, which signs the same browser task before real preparation and fake-child launch. The `surface_operator.go` route hunk applied at line 104 with offset −4, and the `surface_operator_test.go` test hunk applied at line 320 with offset +62. The former retains the S7 method-refusal text for existing routes, and the latter replaces only the cold confirm check with the invalid `confirm/x` assertion. Other hunks applied at their prototype contexts. The landed supported/never classification, independent writer contexts, forged catalogue guards and replay conflict cases retain their pins. A fresh comparison of the iteration-445 stage trees reproduced 517 lines with rename detection (`$S/prototype-numstat.raw`).

The prototype's foreign Host test used a real listener with a valid proof but had no body-read witness. It had a missing approval guard test for the served-store identity of a forged stored row, no concurrent request witness, no native Linux runtime test, and no exact J06/J22 pins. The rebased candidate adds those tests. The schema-10 request API refuses a foreign-store fixture before the handler sees it. The forged-row test inserts via the normal request API and then mutates only the fixture DB bytes/digest. A deleted heartbeat requires generation zero at re-registration. The prototype's pure launch argument fixture did not drive the launch worker, and S6's `TestOperatorStartBrowser` uses a preconfirmed synthetic row. The plan runs independent real HTTP and real fake worker fixtures, with the same task bytes and a claimability readback.

The premise that every daemon route test binds loopback was false. The existing recorder, the new projection and admission-neighbour tests and the planned positive-control probe use the real mux without a socket. The actual `TestOperatorConfirmHTTP` and `TestOperatorAdmissionHTTP` require a listener. The sandbox denied `httptest.NewServer` bind in `listener.raw`; those listener legs are UNINFORMATIVE UNDER SANDBOX. The optional recorder overlay provides a separate measurement only.

| Production calls excluding definitions and tests | B | M1 |
|---|---|---|
| `approval.ConfirmOperatorAction` | One in `host/authority/judge_operator.go:66` | Same |
| `(*store.Credentials).ConfirmOperatorAction` | 0 | One in `handleOperatorConfirm`, `host/daemon/surface_operator_confirm.go:141` |
| `authority.JudgeOperatorAction` | 0 | One in `handleOperatorConfirm`, `host/daemon/surface_operator_confirm.go:117` |
| `RequestOperatorAction` (positive control) | One in `surface_operator.go:226` | One in `surface_operator.go:231` |
| `ClaimOperatorAction` | One in `cmd/ailang-worldd/operator_preparation.go:115` | Same |
| Base 404, request positive control | 404 `NotFound`; request control 200 (`$S/base-probe/test.raw`) | POST dispatch via exact `case suffix == "confirm"`; socket acceptance required |

`$S/callers.raw` is the source census, with independent controls for `StartWorker` and the installed runner call. Guard and context analysis excludes test files.

## 2. Milestones

| Milestone | Product change and acceptance | Expected cumulative numstat |
|---|---|---|
| **M1** | Register `POST /surface/operator-actions/{id}/confirm`, proof and CAS, all §5 refusal pins, card ceremony and process identity status. Run the S6 real fake worker controls, schema-10 proof/CAS and projection gates before activation. Update QUICKSTART. | **+1191/−72 = 1263**; increment equals cumulative. |

The complete slice fits under the cap. One commit keeps the route, status readback, key ceremony, refusal matrix and operator note together. Appendix P is a unified diff from the current full stage tree, not the judged 445 snapshot.

## 3. Durable windows and operator output

Review → client sends one inert request → stable card → exact key proof → confirm CAS → response → separate installed runner can claim. The handler evaluates Host and the path/method before reading the proof. It checks Origin, Fetch Metadata and the media type and uses the 20 KiB proof cap. Typed canonical bytes and recomputed digest, exact store identity and bound artifact are checked against the pending row. It re-reads the clock before CAS. Existing store code persists approval, approver use and passkey counter in one transaction. Two HTTP requests reach the shared human lock and one updates the counter. An exact artifact retry reads the approved row and returns `already-confirmed`; another artifact returns `Conflict`. The response has no worker/session handle. Cancellation expiry and runner freshness still gate a subsequent claim.

| Surface | Observable |
|---|---|
| GET card, while pending and issued | HTML-escaped review and literal task with a `Confirm this action` section. `operator-ceremony` contains `post`, `digest`, `messageHex`, `namespace`, and `passkeys` grouped by `origin`, `rp` and `allow`. It inherits no-store/CSP and exact same-origin links. Expired, pre-issue and used cards omit ceremony. |
| POST confirm | `result` = `confirmed` or `already-confirmed`, `approvedBy` = named enrolled key, and normal status fields `v,id,episode,kind,state,href,statusHref,reportHref,confirmHref,reason,now,createdAt,expiresAt,terminalAt,runnerSeenAt` subject to the existing omitempty shape. Only this POST returns the two extra fields. Wrong methods at the new suffix return 405 without a grant. |
| GET status and card state | A matching live kernel birth shows `launched`; a missing/mismatched birth or reservation shows `uncertain`. `spawn-refused` is the positive non-spawn receipt, and terminal receipt status remains `stopped`, `failed` or `capped`. `reason` carries the existing attended fence or runner-disconnect note. Links remain validated against typed identities. |
| Requests, Needs you, logs | A confirmed row drops out of pending reads and counts. The routine JSON and daemon log sink have no task, signature or assertion. This change adds no log line. Existing runner logs still report action ID and `launched` or terminal result. |

The kernel identity reader moves into `host/agentrun`. On Darwin it reads sysctl; on Linux it uses procfs boot ID plus start ticks and refuses missing, truncated, invalid and zombie readings. The unsupported platform reports unavailable. Polling is observational. It does not authorize repeat execution or clear a reservation. Reuse the positive live fake-child observation in `TestOperatorStartBrowser` and the outside named status/claim selectors.

## 4. Acceptance, timing and fail sets

Run `accept.sh st-M1 M1 fast` under sandbox and `outside.sh st-M1 <new-logdir> M1` outside. Record one `-race -v -count=1` duration for each new daemon test and a before/after for the modified launch case. Pure and fixture times include actual test execution and are distinct from the package/rig clocks. Three named daemon tests must sum to ≤12 s. Cmd has no new top-level test; the modified launch case must add ≤10 s measured base-to-head. The full CI per-package bound stays 900 s.

| Package / race leg | B | M1 | Delta / disposition |
|---|---:|---:|---|
| host/store whole | 302.99 s | 255.41 s | -47.583 s. Matched GOMAXPROCS=2 rerun. |
| host/agentrun whole | 23.05 s | 28.30 s | +5.252 s. Matched GOMAXPROCS=2 Darwin rerun. |
| cmd/ailang-worldd `TestOperatorStartBrowser` | 13.28 s | 7.77 s | -5.51 s measured, 0 new top-level tests. Whole-package delta pending outside. |
| host/daemon listener and whole race | Controller baseline | UNINFORMATIVE UNDER SANDBOX | Actual total and delta from `outside.sh` and `gates.sh ... outside`. |

Current recorder diagnostic `$S/final-daemon-race.raw` measured `TestOperatorConfirmHTTP` 4.87 s, `TestOperatorConfirmProjection` 1.96 s, `TestOperatorConfirmationAdmissionNeighbours` 0.39 s for 7.22 s total. The earlier overloaded 31.72 s HTTP-recorder reading led to a single open read-only evidence connection and one logical SQLite snapshot per refusal. Per-write budgets are fixture-injected and independent. The receiver consumes real socket requests only in the outside leg.

Local `ACCEPT_OK M1 fast` was recorded.
The focused portability leg recorded `ACCEPT_OK M1 1777`.
The matched 1777 fixture selector passed at M1 after a base Git timeout. The controller cmd-plain base records TestOperatorStartCollisions and TestOperatorWorktreeWindow. The first local acceptance sweep hit the inherited `TestOperatorWorktreeWindow/registration_locked` and the first 1777 sweep hit `TestOperatorWorktreeLockedReread/generation` with a 30 s Git fixture timeout. The isolated registration arm passed at base and M1. Whole store and agentrun plain/race added zero failures and local whole verifygate passed. The completed pristine verifygate leg had one failure. See `$S/partial-failsets.tsv`. The broker listener makes its sandbox readings uninformative.
The sealed 1777 attempt also observed a real service record replacement at 20:46:25; the credentials-store and LaunchAgent observations matched. The source is unattributed. The selected service fixture resolves a `t.TempDir()` home, installs a fake controller and uses in-process HTTP transport. The observation is a failed isolation receipt. Use a stable outside observation window for activation.
The initial local gate stream was interrupted by an edit to its script after its recording legs completed. It has no final marker. The independent raw guard and fail-set receipts above survive that instrument failure. Execute the sealed scripts unchanged for the controller gate.
Whole verifygate on the sealed candidate recorded the same TestMissionConfigDefaultFlowIntegration fail set as its pristine leg, with zero new failures. The earlier whole leg was clean. See `$S/sealed-verifygate/test.fails` and `$STATE/base/verifygate-plain.fails`.
Sealed gofmt, build, Darwin and Linux vet, and test compilation passed in `$S/sealed-validation/legs.txt`. Linux agentrun tests cross-compiled for amd64. Native CI still executes real procfs.
The contention diagnostic `$S/final-race/diagnostic-contention.raw` ran the recorder three times at GOMAXPROCS=1 (4.43 s, 5.64 s and 4.91 s summed test seconds). Transport timing remains an outside measurement.

Controller baseline is still collecting. The completed raw legs so far include store/authority/agentrun plain and race and verifygate/verbs/broker plain. Wait for its exact completion marker before treating the combined fail-set comparison as final.


Known base observations at planning cutoff and stage comparison are in `$S/local-M1`, with baseline `$STATE/base`. Controller runs the outside stage sweep and compares normalized fail sets including the additional cmd-race base. Treat a panic timeout as a truncated package. Active GPU and fixture load can widen individual runs without establishing a regression. Outside commit and acceptance must pass their named assertions. No load-shaped red is cleared by editing a fixture unrelated to S7b.

## 5. Refusal, effect, test and mutation matrix

For every actual refusal, assert status and refusal class, absence of Location/CORS on confirm, unchanged row SHA and row count, unchanged passkey counter SHA and count, unchanged credentials/world file SHA, zero effect receipts, zero starts/sessions and an absent fake-child marker. The positive two-key, replay, claim and fake-worker controls establish that each fixture can progress when authorized. The recorder and HTTP fixtures share the real migrated SQLite store schema, daemon mux, repo SSH signer and generated authenticator assertions. S3 and S6 test-only faults create isolated evidence under `t.TempDir()`.

| Reviewed §5.2 row | Disposition and assertion | Named authority/route mutant |
|---|---|---|
| Both keys | M1 `TestOperatorConfirmHTTP/{ConcurrentReplay,SSHSIG}`, plus `TestOperatorConfirmBothKeys` and `TestOperatorJudgeBothKeys`. Correct approver, durable row and artifact. | K-JUDGE-AS-SUCCESS, K-NO-CAS, K-ARTIFACT-FORM |
| Plain y, missing/malformed, wrong challenge | M1 `TestOperatorConfirmHTTP/Proof` and `TestOperatorConfirmProof`. 400 for body grammar, 422 for wrong signature/challenge. | K-BODY-EXACT, K-BODY-PARSER, K-JUDGE-CHALLENGE |
| Forged/unknown, namespace, UV/RP/origin | M1 `TestOperatorConfirmHTTP/Authenticity` and `TestOperatorConfirmAuthenticity`. Enrolled-key forgery, unknown SSH/passkey, wrong namespace, invalid UV, RP or origin each rejected. | A-SSH-NS, A-SSH-ENROLL, A-UV, A-RP, A-ORIGIN, A-SIGNATURE |
| Revocation / counter | M1 `TestOperatorConfirmHTTP/RevocationCounter` and `TestOperatorConfirmRevocationCounter`. Revoked signer, zero/decreased and unchanged nonzero counter. | A-REVOKED, A-COUNTER |
| Other action/episode/kind/store/argv/preset/task and old artifact | M1 `TestOperatorConfirmHTTP/{Binding,Authenticity,StoredIdentity,UnsupportedKinds}` and `TestOperatorConfirmBinding`. Same challenge with each alternate digest, rejected; V3/A3 cannot confer authority. | A-MESSAGE-DIGEST, K-ROW-ID, K-ROW-EP, K-ROW-KIND, K-ROW-DIGEST, K-SERVED-STORE |
| Expired/equal/pre-issue/601s | M1 `TestOperatorConfirmHTTP/Expiry` and `TestOperatorConfirmExpiry`. A fresh time at the write catches expiry during judgment. | K-EXPIRE-PRE, K-EXPIRE-CLASS, K-FRESH-CLOCK, A-LIFETIME |
| Approved replay / lost response | M1 `TestOperatorConfirmHTTP/ConcurrentReplay` and `TestOperatorConfirmReplay`. One counter advance, exact replay, different artifact refusal. | K-APPROVED-ARTIFACT, K-APPROVED-RESULT, K-CONFIRM-LOCK, C-REPLAY-WRITE |
| Missing/busy/counter CAS/death | M1 `TestOperatorConfirmHTTP/{MissingAndWriter,Rollback}`, `TestOperatorConfirmCrash`, `TestOperatorConfirmReadOnly`, `TestOperatorConfirmStoreBinding`. SQL abort after use and counter CAS fault remain rolled back; held writer retry; missing ID 404. | K-BUSY-STATUS, K-CAS-COUNTER-CLASS, C-SPLIT-TX, C-STATE-WRITE |
| Expiry/cancel/readiness before claim | Existing `TestOperatorClaimBeforeExpiry`, `TestOperatorCancelConfirmed`, plus M1 `TestOperatorConfirmHTTP/ConcurrentReplay` removes heartbeat before confirm and restores one only to prove claimability. | R-CLAIM-READY, R-CLAIM-EXPIRY |
| Two runners or stale generation | Existing `TestOperatorClaimRace` and `TestOperatorRunnerKernelLock` run at activation. | R-CLAIM-GENERATION |
| Terminal exclusion | Existing store and cmd `TestOperatorTerminalExclusion` plus `TestOperatorLaunchWorker`. | R-CLAIM-OWNER, W-FLOCK |
| Forged row after claim | Existing `TestOperatorRunnerRechecks`, `TestOperatorStartChanged` from S5 and `TestOperatorLaunchDispatchRechecks`. | R-APPROVER |
| Restart before reserve | Existing `TestOperatorCurrentOwnerProcessRecovery`, `TestOperatorStartRecoveryWindows`, and `TestOperatorLaunchCrashWindow`. | L-RESERVED |
| Restart after reserve/live or unknown child | M1 `TestOperatorConfirmProjection` with real kernel birth and wrong birth, plus `TestOperatorStartBrowser`, `TestOperatorLaunchReceiptLaw` and S6 fake-child gates. | J06, L-BIRTH, L-PID-LIVE |
| Status disconnected/invalid | Existing `TestOperatorStatusDisconnected` and `TestOperatorLinksProjection`; M1 `TestOperatorConfirmProjection`, `TestOperatorConfirmationAdmissionNeighbours`. Client's remembered status is in S8. | J18, J22, L-BIRTH |
| Prune/cap/tombstone | Existing `TestOperatorPrune`, `TestOperatorAdmissionBoundaries`, `TestOperatorLaunchReceiptLaw` and S7 pending projection. | R-PRUNE-UNCERTAIN |
| Cancel after checkpoint | Existing `TestOperatorCancelInflight` and `TestOperatorLaunchCrashWindow`. | R-CANCEL-EFFECT, L-RESERVED |

| Additional security row | Test and named mutant |
|---|---|
| Exact Host before read with a valid key proof | `TestOperatorConfirmHTTP/Transport`, read spy behind the real mux plus actual HTTP 421. D-HOST-ORDER, D-HOST-EXACT |
| Missing/null/substring/duplicated Origin and fetch-site/media guards | `TestOperatorConfirmHTTP/Transport`. D-ORIGIN-COUNT, D-ORIGIN-EXACT, D-FETCH-COUNT, D-FETCH-EXACT, D-MEDIA-COUNT, D-MEDIA-TYPE, D-AUTH, D-COOKIE |
| Separate proof cap/wrong method/OPTIONS | `TestOperatorConfirmHTTP/Transport` and the existing request's 64 KiB sibling. K-SIZELIMIT, K-METHOD, K-REG, D-NO-CORS |
| Whole typed id/path/no redirect | `TestOperatorConfirmHTTP/Transport` and modified `TestOperatorAdmissionHTTP/UnsupportedMethodAndPath`. D-PATH-ESCAPED, D-PATH-CLEAN, D-FAMILY-CLEAN, D-QUERY, D-ID |
| Literal browser task and privacy | `TestOperatorConfirmHTTP/{Ceremony,SSHSIG}` verifies escaped HTML and no payload in routines/logs. `TestOperatorStartBrowser` observes the actual fake child argument vector and checks exactly one prompt argument. L-ARGV-SPLIT, K-CEREMONY-CHALLENGE, K-LOG-PROOF |
| Never dispatch bar tuples | `TestOperatorNeverDispatch` and `TestOperatorAdmissionHTTP` bar rows for world-publish, approver enroll/revoke, session mint/revoke and service. D-NEVER |
| Handler releases writer, installed runner owns execution | `TestOperatorConfirmHTTP/ConcurrentReplay` verifies writer serialization before CAS, released after response, no effects and later claim. `TestOperatorLaunchWorker` and `TestOperatorLaunchExit` run the later fake child. K-CONFIRM-LOCK, K-NO-CAS |

### Carried judge-445 notes

Exact edits come from `world-iter445/judge/mut/specs1.py` and the compilable J22 revision in `specs2.py`, as described in the evaluator reports. The J06 projection assignment is retained in the daemon to run the exact original edit even with the shared status method. The first J22 deletion appeared as NOCOMPILE in round 1; the judge's exact compiling `false &&` edit is the measured one.

| Judge id | Disposition |
|---|---|
| J06 | Pin `TestOperatorConfirmProjection`, `NoSpawn` → `spawn-refused` in status JSON. |
| J22 | Pin `TestOperatorConfirmProjection`, digest drift → 409 without leaked ceremony or mutation. |
| J02 | Pin `TestOperatorConfirmationAdmissionNeighbours`, zero submitted catalogue generation refuses 400 before review. |
| J10 | Pin `TestOperatorConfirmationAdmissionNeighbours`, zero heartbeat generation fails the live review seam. The HTTP result remains terminal. |
| J18 | Pin `TestOperatorConfirmationAdmissionNeighbours`, unreadable operator rows make Needs you `unknown`. |
| J19 | Pin original request lock in `TestOperatorConfirmationAdmissionNeighbours`. The confirm lock is separately mutated by K-CONFIRM-LOCK and the two overlapping confirm requests are measured in `ConcurrentReplay`. |

### Existing tests and fixtures modified

| Path / symbol | Production behaviour before and after | Mutation still detected |
|---|---|---|
| `host/daemon/surface_approve_rig_test.go` `approvalRig.write` | Each fixture write retains its own context. S7b callers inject a bounded minute; existing callers retain the five-second default. | D-MEDIA-COUNT and J18 exercise the rig with durable readback |
| `host/daemon/surface_operator_test.go` `newOperatorHTTP`, `newOperatorSurface` | Same migrated store, pinned clock, heartbeat and real mux. S7b opts into per-write/request budgets before the first commit. | J19, D-HOST-ORDER and K-REG exercise both sides of the shared fixture |
| `host/daemon/surface_operator_test.go` `TestOperatorAdmissionHTTP/UnsupportedMethodAndPath` | Whole exact path family, previous cold confirm 404 probe becomes 404 at invalid `confirm/x`; actual POST and 405 are verified in the new test. | K-REG on the positive confirm, D-QUERY/D-PATH-CLEAN on exact-path refusals |
| `cmd/ailang-worldd/operator_launch_test.go` `TestOperatorStartBrowser`, `newOperatorLaunchRig` | Fresh real preparation, then observed fake child and live-status-on-restart. The fixture can specify the browser's literal task while preserving the original defaults of other callers. | L-ARGV-SPLIT and L-BIRTH; W-FLOCK on the shared worker |
| `cmd/ailang-worldd/operator_activation_test.go` `operatorActivationFixtureWithAgent` | Same signed action, installed paired runner and runner lock; the new task parameter flows through the shared fixture. | L-ARGV-SPLIT, R-APPROVER via preparation/worker controls |
| `cmd/ailang-worldd/operator_worktree_test.go` `newWorktreeStartRigWithAgent` | Same worktree and frozen resolution setup, now accepts an injected task in the new caller. | L-ARGV-SPLIT, W-FLOCK via the actual fake worker |

The catalogue tests, including all round-2 subtests, are byte-identical and are driven through `TestOperatorLaunchWorker/catalogue`. Store launch and authority gate tests are reused in full; no pre-existing tests are deleted.

### Stem-matched selectors

| Touched path family | Names included in §7's named gates |
|---|---|
| `host/daemon/surface_operator*.go` | `TestOperatorAdmissionHTTP`, `TestOperatorLinksProjection`, `TestOperatorReviewRefusesForgedCatalogue`, `TestOperatorReplayConflicts`, `TestOperatorConfirmHTTP`, `TestOperatorConfirmProjection`, `TestOperatorConfirmationAdmissionNeighbours`. `surface_operator_read` has no separate test file. |
| `cmd/ailang-worldd/operator_launch*.go` | `TestOperatorSpawnRefused`, `TestOperatorLaunchDispatchRechecks`, `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchExit`, `TestOperatorStartBrowser`. |
| `cmd/ailang-worldd/operator_activation*_test.go` | `TestOperatorFirstManagedEpisodeProcess`, `TestOperatorPreparationRechecksHeartbeatAfterWorktree`, `TestOperatorStartMintRecoveryRemedies`, `TestOperatorCurrentOwnerProcessRecovery`, `TestOperatorAttendedStartDisposition`, `TestOperatorActivationDrainRefusals`, `TestOperatorRevokeStatRefusal`. |
| `cmd/ailang-worldd/operator_worktree*_test.go` | `TestOperatorStartOnce`, `TestOperatorStartCollisions`, `TestOperatorWorktreeWindow`, `TestOperatorWorktreeDispositionCAS`, `TestOperatorWorktreeSurvivors`, `TestOperatorWorktreeRecoveryGuards`, `TestOperatorWorktreeReadbackGuards`, `TestOperatorReadbackErrorIsNeverCreated`, `TestOperatorWorktreeLockedReread`, `TestOperatorBorrowedWorktreeGitProcess`. |
| `cmd/ailang-worldd/operator_process_{darwin,linux,other}.go` | No matching test files. `TestOperatorProcessWitness`, `TestOperatorStartedRequiresBirth` and `TestOperatorStartBrowser` cover callers. |
| `host/agentrun/process_{darwin,linux,other}*.go` | Native `TestDarwinProcessIdentity`, `TestLinuxProcessIdentity`, `TestUnsupportedProcessIdentity`. Linux CI must execute the procfs live and synthetic branches. |
| `host/store/operator_launch*.go` | `TestOperatorLaunchReceiptLaw`, `TestOperatorLaunchSession`, `TestOperatorLaunchProcessStatus`. |
| `host/daemon/surface_approve_rig_test.go` | No test functions in this helper stem. `TestOperatorConfirmHTTP`, the recorder tests and the existing approval inventory drive it. |
| `docs/QUICKSTART.md` and nearby goldens | Existing verb and requests golden readers in §1.2; example runs the named `TestOperatorConfirmHTTP` on the same fake fixture. |

`TestOperatorStartChanged` is a S5 caller of the unchanged start seam rather than a stem match. Drive it alongside `TestOperatorLaunchDispatchRechecks`.

## 6. Gates and mutant receipts

`$S/mutants.json` is the executable catalogue. It records exact stage path and line, replacement, behavioural test, package and provenance. `mutant-diff-ids.txt` covers the changed Go branches and ordering/projection edges. `mutant-activation-ids.txt` covers the inherited guards and CAS/runner behaviours that become reachable. `mutant-carried-ids.txt` holds the six exact judge edits. `whole-mutant-ids.txt` prioritizes <=25 for whole-package runs. An original report's result is provenance only. The named and whole executions use source overlays with a firing parser control and retain exact raw commands.

| Ledger | Measured disposition |
|---|---|
| Own named, `$S/ledger-own.jsonl` | 120 edits (80 diff-anchored, 40 activation dependencies). EQUIVALENT 1, KILLED 32, RECORDER_KILLED 76, UNMEASURED ON THIS PLATFORM 11. |
| Carried named, `$S/ledger-carried.jsonl` | KILLED 6. J02, J06, J10, J18, J19 and J22 are exact 445 edits. |
| Prioritized whole, `$S/whole-local-final/ledger.jsonl` | 6 local rows. KILLED 6. Remaining whole rows are outside or native Linux. |

No measured, non-equivalent named survivor remains. The first recorder instruments exposed J18, D-PATH-CLEAN and K-READ-BEFORE-LOCK evidence gaps. The current fixtures drive the operator-only read failure, a noncanonical confirm suffix, and the exact concurrent subtest assertion. D-HOST-ORDER now reads only at the confirmation probe so the inert request positive control can still run. The raw-log decoder preserves bytes from binary fixture text.

Run listener mutants outside; the in-sandbox recorder receipt is marked RECORDER_KILLED even when it detects the correct failure. `K-ONLY-START` is an equivalence control: canonical V4 finish cannot have `Headless: true`, and the adjacent headless guard still refuses valid finish rows. Any other surviving named edit blocks acceptance. Native Linux procfs requires Linux CI. The synthetic parser probe in `$S/linux-parser-probe.raw` separately runs the fault inputs against the same parser on Darwin. It does not substitute for the real `/proc` self check. Compilation across GOOS is recorded separately.

The whole gate suite records baseline-like outputs before comparison. Each whole verifygate and broker leg has its own state/cache directories. The verify_ail wrapper runs the repository verifier with only its mktemp destinations relocated under scratch. No PATH shims and no GOFLAGS JSON. `gates.sh` records completion last with `ALL_DONE`; status 77 retains an unresolved sandbox/baseline indication.

### Reproduction scripts


````bash gates.sh
#!/bin/bash
# gates.sh TREE LOGDIR M1 [local|outside|outside-base]
# The controller records an outside pristine cmd race leg with this same timeout.
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P) || exit 2
T=$(cd "$1" && pwd -P) || exit 2
mkdir "$2" || exit 2
L=$(cd "$2" && pwd -P) || exit 2
M=$3
mode=${4:-local}
[ "$M" = M1 ] || exit 2
case "$mode" in local|outside|outside-base) ;; *) exit 2;; esac
B=88fea1e9a8422b12a059eb1ecfc9c10c5b7b4e7d
CONTROL="$S/../base"
. "$S/selectors.sh"
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_TEST_OPERATOR_WITNESS WORLD_TEST_OPERATOR_WITNESS_PATH
cd "$T" || exit 2
issue=0
unknown=0
note() { printf '%s\n' "$*" | tee -a "$L/legs.txt"; }
skip() { note "LEG $1 SKIPPED ($2)"; }
bound() { python3 "$S/bounded.py" "$@"; }
denied() { grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$1"; }
truncated() { grep -Eq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$1"; }
leg() {
  local tag=$1 start rc=0
  shift
  start=$(date +%s)
  mkdir -p "$L/$tag.state" "$L/$tag.cache"
  (export AILANG_STATE_DIR="$L/$tag.state" AILANG_CACHE_DIR="$L/$tag.cache"; "$@") > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  if truncated "$L/$tag.raw"; then issue=1; note "LEG $tag INSTRUMENT FAILURE (truncated) rc=$rc"
  elif denied "$L/$tag.raw" || [ "$rc" -eq 77 ]; then unknown=1; note "LEG $tag UNINFORMATIVE UNDER SANDBOX rc=$rc"
  elif [ "$rc" -ne 0 ]; then issue=1; note "LEG $tag rc=$rc $(( $(date +%s)-start ))s"
  else note "LEG $tag OK $(( $(date +%s)-start ))s"; fi
}
record() {
  local tag=$1 pkg=$2 start rc=0 parser=0 cap=1260 limit=1200s
  shift 2
  case "$tag" in
    *-race) cap=950; limit=900s ;;
    cmd-plain|daemon-plain) cap=2800; limit=45m ;;
  esac
  start=$(date +%s)
  mkdir -p "$L/$tag.state" "$L/$tag.cache"
  (export AILANG_STATE_DIR="$L/$tag.state" AILANG_CACHE_DIR="$L/$tag.cache"
   bound "$cap" go test -count=1 -timeout "$limit" "$@" "./$pkg/") > "$L/$tag.raw" 2>&1 || rc=$?
  printf '%s\n' "$rc" > "$L/$tag.rc"
  python3 "$S/failset.py" "$L/$tag.raw" "$rc" > "$L/$tag.fails" 2> "$L/$tag.parser" || parser=$?
  cat "$L/$tag.parser"
  if denied "$L/$tag.raw"; then unknown=1; note "LEG $tag UNINFORMATIVE UNDER SANDBOX rc=$rc"
  elif truncated "$L/$tag.raw" || [ "$parser" -ne 0 ]; then issue=1; note "LEG $tag INSTRUMENT FAILURE rc=$rc parser=$parser"
  else cp "$L/$tag.fails" "$L/$tag.attributable"; note "LEG $tag recorded rc=$rc $(( $(date +%s)-start ))s"; fi
  grep -E '^(ok|FAIL)[[:space:]]+github.com/' "$L/$tag.raw" >> "$L/pkg-times.txt" || true
}
fmt_touched() {
  local p bad
  while IFS= read -r p; do
    case "$p" in *.go) bad=$(gofmt -l "$p") || return 2; test -z "$bad" || { printf '%s\n' "$bad"; return 1; };; esac
  done < "$S/changed-M1"
}
inventory() {
  local pkg=$1 names=$2 output rc=0
  output=$(mktemp "$L/inventory.XXXXXX") || return 2
  bound 380 go test -v -count=1 -timeout 350s "./$pkg/" -run "^($names)\$" > "$output" 2>&1 || rc=$?
  cat "$output"
  [ "$rc" -eq 0 ] || return "$rc"
  while IFS= read -r name; do grep -Fq -- "--- PASS: $name (" "$output" || return 2; done < <(printf '%s\n' "$names" | tr '|' '\n')
}
pkgs=(./cmd/ailang-worldd/ ./host/agentrun/ ./host/daemon/ ./host/store/)
tags=(store-plain store-race agentrun-plain agentrun-race verifygate-plain broker-plain cmd-plain cmd-race daemon-plain daemon-race)
owners=(host/store host/store host/agentrun host/agentrun host/verifygate host/broker cmd/ailang-worldd cmd/ailang-worldd host/daemon host/daemon)
if [ "$mode" = outside-base ]; then
  record cmd-race cmd/ailang-worldd -race
  note ALL_DONE
  test -f "$L/cmd-race.attributable"
  exit $?
fi
leg gofmt fmt_touched
leg build bound 440 go build ./...
leg vet bound 350 go vet "${pkgs[@]}"
leg linux-vet env GOOS=linux python3 "$S/bounded.py" 440 go vet "${pkgs[@]}"
leg inv-store inventory host/store "$STORE_INVENTORY"
leg inv-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
leg inv-cli inventory cmd/ailang-worldd 'TestCLICredentialWriterSites|TestSessionRootBudgets'
leg inv-daemon inventory host/daemon "$DAEMON_INVENTORY"
leg inv-daemon-headers inventory host/daemon "$DAEMON_HEADER_INVENTORY"
leg inv-daemon-read inventory host/daemon "$DAEMON_READ_INVENTORY"
leg inv-verbs inventory host/verbs "$VERBS_INVENTORY"
leg inv-verifygate inventory host/verifygate "$VERIFY_INVENTORY"
leg inv-broker inventory host/broker "$BROKER_INVENTORY"
if [ "$(uname -s)" = Linux ]; then
  leg inv-exec-kernel inventory host/verifygate TestKernelRefusesExecOfWriterOpenFile
else
  skip inv-exec-kernel 'Linux runtime control'
fi
leg verify_ail bound 650 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
leg callers /bin/bash "$S/callers.sh" "$S/base-tree" "$T"
if [ "$mode" = outside ]; then
  leg accept /bin/bash "$S/accept.sh" "$T" "$M" outside
  leg accept-race /bin/bash "$S/accept.sh" "$T" "$M" outside-race
else
  leg accept /bin/bash "$S/accept.sh" "$T" "$M" fast
fi
mkdir -m 1777 "$L/t1777"
mkdir -m 0700 "$L/t1777/x"
leg t1777 env TMPDIR="$L/t1777/x" /bin/bash "$S/accept.sh" "$T" "$M" 1777
for i in "${!tags[@]}"; do
  tag=${tags[$i]}
  if [ "$mode" = local ]; then
    case "$tag" in cmd-*|daemon-*) skip "$tag" 'outside sandbox required'; continue;; esac
  fi
  case "$tag" in *-race) record "$tag" "${owners[$i]}" -race;; *) record "$tag" "${owners[$i]}";; esac
done
# All legs have been recorded. Compare test names, never return-code differences.
for tag in "${tags[@]}"; do
  if [ ! -f "$L/$tag.attributable" ]; then skip "compare-$tag" 'head record unavailable'; continue; fi
  source=
  if [ -f "$S/../baseline.done" ] && [ "$(tail -n 1 "$S/../baseline.done")" = "ALL_DONE for $B" ] &&
      [ -s "$CONTROL/$tag.raw" ] && [ -f "$CONTROL/$tag.rc" ] &&
      ! truncated "$CONTROL/$tag.raw" && ! denied "$CONTROL/$tag.raw" &&
      python3 "$S/failset.py" "$CONTROL/$tag.raw" "$(cat "$CONTROL/$tag.rc")" > "$L/$tag.controller"; then
    source="$L/$tag.controller"
  fi
  if [ "$tag" = cmd-race ] && [ -f "$S/outside-pristine/cmd-race.attributable" ]; then
    source="$S/outside-pristine/cmd-race.attributable"
  fi
  if [ ! -f "$source" ]; then skip "compare-$tag" 'final pristine baseline unavailable'; unknown=1; continue; fi
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
# accept.sh TREE M1 [fast|race|outside|outside-race|1777]
set -euo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
M=$2
mode=${3:-fast}
[ "$M" = M1 ] || exit 2
race=
case "$mode" in fast|outside|1777) ;; race|outside-race) race=-race;; *) exit 2;; esac
. "$S/selectors.sh"
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_TEST_OPERATOR_WITNESS WORLD_TEST_OPERATOR_WITNESS_PATH
D=$(mktemp -d "$S/accept.XXXXXX")
export AILANG_STATE_DIR="$D/state" AILANG_CACHE_DIR="$D/cache"
mkdir -p "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
observe() {
  local p
  for p in "$HOME/Library/LaunchAgents/com.sunholo.ailang-world.worldd.plist" \
    "$HOME/Library/LaunchAgents/com.sunholo.ailang-world.operator-runner.plist" \
    "$HOME/.ailang-world" "$HOME/.ailang-world/service.json" "$HOME/.ailang/world-ws" \
    "$HOME/.ailang/world/credentials.db"; do
    if [ ! -e "$p" ] && [ ! -L "$p" ]; then printf 'absent %s\n' "$p"; continue; fi
    if [ "$(uname -s)" = Darwin ]; then stat -f '%N %p %i %z %m %c' "$p"
    else stat -c '%n %a %i %s %Y %Z' "$p"; fi
    case "$p" in */credentials.db) ;; *)
      if [ -f "$p" ] && [ -r "$p" ]; then shasum -a 256 "$p"; fi ;;
    esac
  done
}
observe > "$D/before"
trap 'observe > "$D/after"; diff -u "$D/before" "$D/after"' EXIT
cd "$T"
run() {
  local tag=$1 pkg=$2 names=$3 log rc=0
  log="$D/$tag.raw"
  python3 "$S/bounded.py" 940 go test ${race:+-race} -count=1 -v -timeout 900s "./$pkg/" -run "^($names)\$" > "$log" 2>&1 || rc=$?
  cat "$log"
  if grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$log"; then
    echo "LEG $tag UNINFORMATIVE UNDER SANDBOX"
    return 77
  fi
  if grep -Eq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$log"; then return 2; fi
  if [ "$rc" -ne 0 ]; then return "$rc"; fi
  while IFS= read -r name; do
    grep -Fq -- "--- PASS: $name (" "$log" || return 2
  done < <(printf '%s\n' "$names" | tr '|' '\n')
}
if [ "$mode" = 1777 ]; then
  run cli-1777 cmd/ailang-worldd 'TestOperatorStartBrowser|TestOperatorLaunchWorker|TestOperatorSupervisorCap|TestOperatorStartedRequiresBirth'
else
  run agentrun host/agentrun "$PROCESS_TEST"
  run store host/store "$STORE_STEMS|$STORE_LAWS"
  run authority host/authority "$AUTHORITY_LAWS"
  run approval host/approval "$APPROVAL_LAWS"
  run verbs host/verbs "$VERBS_INVENTORY"
  run cli cmd/ailang-worldd "$CLI_STEMS|$CLI_S6"
  run golden host/surfacepage "$PAGE_GOLDEN"
  case "$mode" in
    outside|outside-race) run daemon host/daemon "$DAEMON_STEMS|$DAEMON_LISTENER|$DAEMON_GOLDEN" ;;
    *) run daemon host/daemon "$DAEMON_STEMS|$DAEMON_GOLDEN"
       echo "LEG listener-acceptance SKIPPED (UNINFORMATIVE UNDER SANDBOX; use outside.sh)" ;;
  esac
fi
observe > "$D/after"
diff -u "$D/before" "$D/after"
trap - EXIT
printf 'ACCEPT_OK %s %s\n' "$M" "$mode"
````
## 7. Apply, accept and snapshot

The controller checkpoints these two plan files. The executor works in the designated worktree and applies M1 once. Use this invocation verbatim, from the checkout root. It uses the already prepared scratch support files and re-extracts the scripts that the executor runs.

````bash
set -euo pipefail
B=88fea1e9a8422b12a059eb1ecfc9c10c5b7b4e7d
STATE="$HOME/.ailang/state/world-iter446"
S="$STATE/plan-scratch"
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter446-row243-v4s7b
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v4s7b-sprint-plan.md"
export GIT_OPTIONAL_LOCKS=0
cd "$W"
git cat-file -e "$B^{commit}"
test -z "$(git diff --name-only "$B" -- cmd host docs app)"
test -z "$(git ls-files --others --exclude-standard -- cmd host docs app)"
X=$(mktemp -d "$S/executor.XXXXXX")
grab() { awk -v f="$1" '$0 == f {a=1; next} a && $0 == "````" {exit} a' "$PLAN"; }
for tool in accept.sh gates.sh snapshot.sh; do
  grab "````bash $tool" > "$X/$tool"
  test -s "$X/$tool"
  cmp "$X/$tool" "$S/$tool"
done
grab '````diff M1' > "$X/M1.patch"
test -s "$X/M1.patch"
/usr/bin/patch -p1 -E < "$X/M1.patch"
bash "$S/accept.sh" "$W" M1 fast
bash "$S/snapshot.sh" "$W" M1
(cd "$STATE/snap/M1" && shasum -a 256 -c MANIFEST.sha256)
````

`accept.sh` succeeds with exit 0 and prints **`ACCEPT_OK M1 <mode>`** after all positive test selectors and before/after observations. The fast mode runs every stem-matched selector and leaves socket acceptance to the controller. The 1777 mode runs the four actual fake-worker and birth controls, including the modified browser-task launch. Both outside modes separately run the full stem selectors and listener tests. The snapshot is cumulative at `$HOME/.ailang/state/world-iter446/snap/M1/` and includes `STATUS`, `DELETED`, and `MANIFEST.sha256`, followed by every added/changed repo-relative file. It excludes the two plan documents. Expect 22 STATUS rows, 1 deletion and 21 manifest entries. There is no dependency on the presence of any particular product directory.

### Controller commands

Run these outside before invoking the commit script. Use fresh receipt directories for retries. The whole-package cmd/race baseline has no controller leg in the supplied series.

````bash
B=88fea1e9a8422b12a059eb1ecfc9c10c5b7b4e7d
STATE="$HOME/.ailang/state/world-iter446"
S="$STATE/plan-scratch"
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter446-row243-v4s7b
test "$(tail -n 1 "$STATE/baseline.done")" = "ALL_DONE for $B"
bash "$STATE/ctl_roundtrip.sh" M1
bash "$S/outside.sh" "$S/st-M1" "$S/outside-M1" M1
bash "$S/gates.sh" "$S/base-tree" "$S/outside-pristine" M1 outside-base
bash "$S/gates.sh" "$S/st-M1" "$S/outside-gates-M1" M1 outside
grep -v '^P-LINUX-' "$S/all-own-ids.txt" > "$S/darwin-own-ids.txt"
grep -v '^P-LINUX-' "$S/whole-mutant-ids.txt" > "$S/darwin-whole-ids.txt"
bash "$S/whole-mutants.sh" "$S/st-M1" "$S/outside-named" "$S/darwin-own-ids.txt" named
bash "$S/whole-mutants.sh" "$S/st-M1" "$S/outside-judge-named" "$S/mutant-carried-ids.txt" named
bash "$S/whole-mutants.sh" "$S/st-M1" "$S/outside-whole" "$S/darwin-whole-ids.txt" whole
````

For whole mutations in a package without a controller baseline, first record its pristine leg. `host/approval` uses `$S/mutants-pristine/approval-plain.{raw,rc,fails}` from the planner. `host/daemon`, `cmd/ailang-worldd`, `host/store`, `host/authority` and `host/agentrun` use `$STATE/base/{daemon,cmd,store,authority,agentrun}-plain.fails`. Complete the Linux mutation rows in native Linux CI with the `P-LINUX-*` selector entries, and run `go test -race -count=1 ./host/agentrun/ -run '^TestLinuxProcessIdentity$'` on Linux. The normal Linux CI whole package already executes that test. On the Linux CI copy of these instruments, `grep '^P-LINUX-' "$S/all-own-ids.txt" > "$S/linux-named-ids.txt"`, then `bash "$S/whole-mutants.sh" "$S/st-M1" "$S/linux-named" "$S/linux-named-ids.txt" named`. For the prioritized Linux whole arm use the corresponding filter over `whole-mutant-ids.txt` and mode `whole`. The Linux platform result remains separate from macOS controller evidence. Preserve pinned `AILANG_BIN` in the environment. No sandbox service-start result can substitute for the outside fail set.

The commit command below uses named tests for cmd and daemon, including the two real-HTTP selectors. It invokes the whole verifygate/broker command with isolated state/cache and compares its recorded fail set to the outside baseline.

````bash
STATE="$HOME/.ailang/state/world-iter446"
S="$STATE/plan-scratch"
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter446-row243-v4s7b
cat > "$S/M1.msg" <<'MSG'
feat(operator): activate key confirmation for one confined headless start

Ship guarded confirmation, exact ceremony and refusal/replay evidence with
live process identity readback. The installed runner owns launch lifetime.
MSG
bash "$STATE/commit_ms.sh" M1 "$S/M1.msg" \
  'export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"; S="$HOME/.ailang/state/world-iter446/plan-scratch"; bash "$S/accept.sh" "$PWD" M1 outside && bash "$S/commit-guards.sh" "$PWD" "$S/commit-guards-M1"'
````

After commit, expect `NUMSTAT 1263` with renames enabled. Verify the full normalized diff against `$S/numstat-M1`. A skipped, truncated, sandbox-denied or newly failing outside leg leaves release acceptance pending.

### Rehearsal

`roundtrip.sh` extracted Appendix P from the plan file and applied it with `/usr/bin/patch -p1 -E` to a fresh export at `/Users/voightkampff/.ailang/state/world-iter446/plan-scratch/rt-4`. M1 `diff -r` returned 0 for cmd, host, docs and the full archive. The app directory is skipped if absent on both sides. The cumulative snapshot verified 21 SHA-256 entries and all 22 STATUS rows. Receipt `$S/roundtrip-final.raw` ends in `ROUNDTRIP_OK`.

````bash snapshot.sh
#!/bin/bash
# snapshot.sh TREE M1|M2; SNAP_STATUS_LIST/ SNAP_PARENT are rehearsal-only.
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1) ;; *) exit 2;; esac
B=88fea1e9a8422b12a059eb1ecfc9c10c5b7b4e7d
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter446/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s7b-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s7b.json"
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

### M1 — confirmation activation

````diff M1
diff --git a/cmd/ailang-worldd/operator_activation_test.go b/cmd/ailang-worldd/operator_activation_test.go
index 41ec39d..1ef9292 100644
--- a/cmd/ailang-worldd/operator_activation_test.go
+++ b/cmd/ailang-worldd/operator_activation_test.go
@@ -38,7 +38,12 @@ func operatorActivationFixture(t *testing.T) (*worktreeStartRig, runnerConfig, s
 
 func operatorActivationFixtureWithAgent(t *testing.T, script string) (*worktreeStartRig, runnerConfig, store.OperatorHeartbeat, *os.File) {
 	t.Helper()
-	r := newWorktreeStartRigWithAgent(t, false, activationTestContext(t), script)
+	return operatorActivationFixtureWithTask(t, script, operatorStartTask)
+}
+
+func operatorActivationFixtureWithTask(t *testing.T, script, task string) (*worktreeStartRig, runnerConfig, store.OperatorHeartbeat, *os.File) {
+	t.Helper()
+	r := newWorktreeStartRigWithTask(t, false, activationTestContext(t), script, task)
 	r.env.operatorBudget = activationTestBudget
 	c := r.config
 	agent := operatorFileFact(r.action, "launch_binary")
diff --git a/cmd/ailang-worldd/operator_launch.go b/cmd/ailang-worldd/operator_launch.go
index e3c1707..fe3ca97 100644
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
diff --git a/cmd/ailang-worldd/operator_launch_test.go b/cmd/ailang-worldd/operator_launch_test.go
index c55d72c..91b097a 100644
--- a/cmd/ailang-worldd/operator_launch_test.go
+++ b/cmd/ailang-worldd/operator_launch_test.go
@@ -38,7 +38,12 @@ type operatorLaunchRig struct {
 
 func newOperatorLaunchRig(t *testing.T, script string) *operatorLaunchRig {
 	t.Helper()
-	r, c, h, runner := operatorActivationFixtureWithAgent(t, script)
+	return newOperatorLaunchRigWithTask(t, script, operatorStartTask)
+}
+
+func newOperatorLaunchRigWithTask(t *testing.T, script, task string) *operatorLaunchRig {
+	t.Helper()
+	r, c, h, runner := operatorActivationFixtureWithTask(t, script, task)
 	f := &operatorLaunchRig{worktreeStartRig: r, c: c, h: h, runner: runner, ctx: activationTestContext(t), now: r.src.clock()}
 	r.src.clock = func() int64 { return f.now }
 	r.env.terminal.now = r.src.clock
@@ -357,8 +362,20 @@ func testOperatorStartOnceLaunch(t *testing.T) {
 }
 
 func TestOperatorStartBrowser(t *testing.T) {
-	f := newOperatorLaunchRig(t, fakeLaunchMarker+"while [ ! -e release ]; do :; done\nexit 0\n")
+	// Same browser-supplied task as TestOperatorConfirmHTTP. The launch rig
+	// signs and freezes it before preparation reads the durable row.
+	task := "# Task\n<script>alert(1)</script> ; $(literal) &"
+	f := newOperatorLaunchRigWithTask(t, fakeLaunchMarker+"while [ ! -e release ]; do :; done\nexit 0\n", task)
 	f.env.operatorAgentRun = func(ctx context.Context, argv []string, dir string, env []string, out, stderr *os.File, wall time.Duration, hooks agentrun.LaunchHooks) (agentrun.Result, error) {
+		promptArgs := 0
+		for _, arg := range argv {
+			if strings.Contains(arg, task) {
+				promptArgs++
+			}
+		}
+		if promptArgs != 1 || !strings.HasPrefix(argv[len(argv)-1], task) {
+			t.Errorf("browser task did not survive launch as one literal prompt argument: %q", argv)
+		}
 		start := hooks.Started
 		hooks.Started = func(pid int) error {
 			defer func() { _ = os.WriteFile(filepath.Join(dir, "release"), nil, 0600) }()
diff --git a/cmd/ailang-worldd/operator_process_darwin.go b/cmd/ailang-worldd/operator_process_darwin.go
index d5121a3..4264764 100644
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
index 4fe7729..0000000
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
index ad91f10..7a30ebd 100644
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
diff --git a/cmd/ailang-worldd/operator_worktree_test.go b/cmd/ailang-worldd/operator_worktree_test.go
index c0e3d56..ae079e2 100644
--- a/cmd/ailang-worldd/operator_worktree_test.go
+++ b/cmd/ailang-worldd/operator_worktree_test.go
@@ -40,6 +40,11 @@ func newWorktreeStartRigWithContext(t *testing.T, claude bool, ctx context.Conte
 }
 
 func newWorktreeStartRigWithAgent(t *testing.T, claude bool, ctx context.Context, body string, ttl ...int64) *worktreeStartRig {
+	t.Helper()
+	return newWorktreeStartRigWithTask(t, claude, ctx, body, operatorStartTask, ttl...)
+}
+
+func newWorktreeStartRigWithTask(t *testing.T, claude bool, ctx context.Context, body, task string, ttl ...int64) *worktreeStartRig {
 	t.Helper()
 	s := operatorService(t)
 	s.mustInstall(t)
@@ -61,7 +66,7 @@ func newWorktreeStartRigWithAgent(t *testing.T, claude bool, ctx context.Context
 	r := operatorStartRig{f: f, bin: bin, creds: s.config.Credentials, hosts: t.TempDir()}
 	r.f.r.f.db = s.config.World
 	r.action = approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: fixtureActionID, Kind: "start", Episode: "ep",
-		Store: r.creds, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: operatorStartTask, Review: "fixture"}
+		Store: r.creds, Options: approval.OperatorOptions{Preset: "p", Headless: true}, Task: task, Review: "fixture"}
 	r.action.Argv = approval.OperatorArgv(r.action.Kind, r.action.Episode, r.action.Options)
 	if claude {
 		if err := os.Rename(filepath.Join(r.bin, "codex"), filepath.Join(r.bin, "claude")); err != nil {
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 21e7b12..ef7c1a8 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -2121,3 +2121,41 @@ launch-status ACTION-ID`. There is no unlock command for `uncertain`. Preserve
 the action id, runner log, launch receipt, worktree and session records for
 attended reconciliation. Check the episode report for existing work. Only
 positive evidence can settle the reservation. Refresh and retry preserve it.
+
+#### Confirm with an enrolled key
+
+Review the exact action and task. The pending card's `operator-ceremony` JSON
+gives the same-origin `post` URL, `digest`, `messageHex`, SSHSIG `namespace`
+and passkey `origin`, `rp` and `allow` credential IDs. For SSHSIG, save that
+JSON in a private directory as `ceremony.json`, then use your enrolled key.
+
+```sh
+jq -r .messageHex ceremony.json | xxd -r -p > action.message
+ssh-keygen -Y sign -n world-operator-action@sunholo -f "$ENROLLED_KEY" action.message
+jq -n --arg digest "$(jq -r .digest ceremony.json)" --rawfile signature action.message.sig \
+  '{v: 1, digest: $digest, sshsig: $signature}' > confirmation.json
+```
+
+POST that JSON with `Content-Type: application/json` from the card's exact
+origin to its `post` URL. A passkey uses SHA-256 of the message bytes as the
+challenge, the card's RP and allow-list for that origin, and
+`userVerification: "required"`. Send the serialized assertion as `webauthn`
+in place of `sshsig`. The POST returns `result`, `approvedBy` and the action
+links. The same artifact can be retried after a lost response and returns
+`already-confirmed`.
+
+A ready runner claims the action after confirmation. Poll `statusHref`.
+`launched` means the live PID and its birth identity match the receipt.
+`spawn-refused` means the runner recorded no spawn. A missing or mismatched
+identity reports `uncertain` and retains the fence described above.
+
+Run a fake confirmation with both key types, real same-origin transport,
+temporary stores and the runner absent. The fixture constructs both payloads
+for the exact bytes shown on the card and verifies one counter advance on replay.
+
+```sh
+AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 \
+WORLD_EXEC_NODE=/opt/homebrew/bin/node \
+WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang" \
+go test -count=1 ./host/daemon/ -run '^TestOperatorConfirmHTTP$'
+```
diff --git a/host/agentrun/process_darwin.go b/host/agentrun/process_darwin.go
new file mode 100644
index 0000000..de1df0c
--- /dev/null
+++ b/host/agentrun/process_darwin.go
@@ -0,0 +1,21 @@
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
+	if pid <= 0 {
+		return "", false
+	}
+	row, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
+	if err != nil || int(row.Proc.P_pid) != pid {
+		return "", false
+	}
+	p := row.Proc
+	const zombie = 5 // SZOMB in Darwin <sys/proc.h>
+	return fmt.Sprintf("%d.%d", p.P_starttime.Sec, p.P_starttime.Usec), p.P_stat != zombie
+}
diff --git a/host/agentrun/process_darwin_test.go b/host/agentrun/process_darwin_test.go
new file mode 100644
index 0000000..8c78989
--- /dev/null
+++ b/host/agentrun/process_darwin_test.go
@@ -0,0 +1,22 @@
+package agentrun
+
+import (
+	"os"
+	"testing"
+)
+
+func TestDarwinProcessIdentity(t *testing.T) {
+	pid := os.Getpid()
+	birth, live := ProcessIdentity(pid)
+	if !live || birth == "" {
+		t.Fatalf("live process missing birth: %q %v", birth, live)
+	}
+	if again, live := ProcessIdentity(pid); !live || again != birth {
+		t.Fatalf("unstable birth: %q %v", again, live)
+	}
+	for _, pid := range []int{0, -1, 1 << 30} {
+		if birth, live := ProcessIdentity(pid); live || birth != "" {
+			t.Fatalf("invalid or absent PID %d: %q %v", pid, birth, live)
+		}
+	}
+}
diff --git a/host/agentrun/process_linux.go b/host/agentrun/process_linux.go
new file mode 100644
index 0000000..2594d80
--- /dev/null
+++ b/host/agentrun/process_linux.go
@@ -0,0 +1,33 @@
+package agentrun
+
+import (
+	"io/fs"
+	"os"
+	"strconv"
+	"strings"
+)
+
+func ProcessIdentity(pid int) (string, bool) {
+	return linuxProcessIdentity(os.DirFS("/proc"), pid)
+}
+
+func linuxProcessIdentity(proc fs.FS, pid int) (string, bool) {
+	if pid <= 0 {
+		return "", false
+	}
+	stat, err := fs.ReadFile(proc, strconv.Itoa(pid)+"/stat")
+	boot, bootErr := fs.ReadFile(proc, "sys/kernel/random/boot_id")
+	end := strings.LastIndexByte(string(stat), ')')
+	if err != nil || bootErr != nil || end < 0 {
+		return "", false
+	}
+	fields := strings.Fields(string(stat[end+1:]))
+	if len(fields) <= 19 || len(strings.TrimSpace(string(boot))) == 0 {
+		return "", false
+	}
+	ticks, err := strconv.ParseUint(fields[19], 10, 64)
+	if err != nil {
+		return "", false
+	}
+	return strings.TrimSpace(string(boot)) + ":" + strconv.FormatUint(ticks, 10), fields[0] != "Z"
+}
diff --git a/host/agentrun/process_linux_test.go b/host/agentrun/process_linux_test.go
new file mode 100644
index 0000000..4ef19f8
--- /dev/null
+++ b/host/agentrun/process_linux_test.go
@@ -0,0 +1,62 @@
+package agentrun
+
+import (
+	"errors"
+	"os"
+	"strings"
+	"testing"
+	"testing/fstest"
+)
+
+func TestLinuxProcessIdentity(t *testing.T) {
+	stat := func(state, ticks string) string {
+		return "123 (a task ) with right parentheses) " + state + strings.Repeat(" 0", 18) + " " + ticks
+	}
+	for _, tc := range []struct {
+		name, stat, boot, want, broken string
+		pid                            int
+		live                           bool
+	}{
+		{"live", stat("S", "17"), "boot\n", "boot:17", "", 123, true},
+		{"zombie", stat("Z", "17"), "boot", "boot:17", "", 123, false},
+		{"partial stat read", stat("S", "17"), "boot", "", "123/stat", 123, false},
+		{"partial boot read", stat("S", "17"), "boot", "", "sys/kernel/random/boot_id", 123, false},
+		{"missing stat", "", "boot", "", "", 123, false},
+		{"missing boot", stat("S", "17"), "", "", "", 123, false},
+		{"blank boot", stat("S", "17"), " ", "", "", 123, false},
+		{"short stat", "123 (short) S", "boot", "", "", 123, false},
+		{"no delimiter", "S" + strings.Repeat(" 0", 18) + " 17", "boot", "", "", 123, false},
+		{"invalid ticks", stat("S", "bad"), "boot", "", "", 123, false},
+		{"negative pid", stat("S", "17"), "boot", "", "", -1, false},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			proc := fstest.MapFS{}
+			if tc.stat != "" {
+				proc["123/stat"] = &fstest.MapFile{Data: []byte(tc.stat)}
+				proc["-1/stat"] = proc["123/stat"]
+			}
+			if tc.boot != "" {
+				proc["sys/kernel/random/boot_id"] = &fstest.MapFile{Data: []byte(tc.boot)}
+			}
+			if got, live := linuxProcessIdentity(procFixture{proc, tc.broken}, tc.pid); got != tc.want || live != tc.live {
+				t.Fatalf("identity = %q, %t, want %q, %t", got, live, tc.want, tc.live)
+			}
+		})
+	}
+	if birth, live := ProcessIdentity(os.Getpid()); !live || birth == "" {
+		t.Fatalf("read actual procfs: %q %v", birth, live)
+	}
+}
+
+type procFixture struct {
+	fstest.MapFS
+	broken string
+}
+
+func (f procFixture) ReadFile(name string) ([]byte, error) {
+	b, err := f.MapFS.ReadFile(name)
+	if name == f.broken {
+		return b, errors.New("fixture: partial filesystem read")
+	}
+	return b, err
+}
diff --git a/host/agentrun/process_other.go b/host/agentrun/process_other.go
new file mode 100644
index 0000000..255b2d9
--- /dev/null
+++ b/host/agentrun/process_other.go
@@ -0,0 +1,5 @@
+//go:build !darwin && !linux
+
+package agentrun
+
+func ProcessIdentity(int) (string, bool) { return "", false }
diff --git a/host/agentrun/process_other_test.go b/host/agentrun/process_other_test.go
new file mode 100644
index 0000000..c56b85b
--- /dev/null
+++ b/host/agentrun/process_other_test.go
@@ -0,0 +1,14 @@
+//go:build !darwin && !linux
+
+package agentrun
+
+import (
+	"os"
+	"testing"
+)
+
+func TestUnsupportedProcessIdentity(t *testing.T) {
+	if birth, live := ProcessIdentity(os.Getpid()); live || birth != "" {
+		t.Fatalf("unsupported process identity: %q %v", birth, live)
+	}
+}
diff --git a/host/daemon/surface_approve_rig_test.go b/host/daemon/surface_approve_rig_test.go
index 71ef019..3cb9084 100644
--- a/host/daemon/surface_approve_rig_test.go
+++ b/host/daemon/surface_approve_rig_test.go
@@ -31,8 +31,9 @@ import (
 const apT0, apLoopback = int64(1800000000), "127.0.0.1:7644"
 
 type approvalRig struct {
-	t *testing.T
-	d *Daemon
+	t           *testing.T
+	d           *Daemon
+	writeBudget time.Duration
 }
 
 func newApprovalRig(t *testing.T) approvalRig {
@@ -54,7 +55,11 @@ func (a approvalRig) write(f func(ctx context.Context, c *store.Credentials) err
 		a.t.Fatal(err)
 	}
 	defer func() { _ = c.Close() }()
-	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	budget := a.writeBudget
+	if budget == 0 {
+		budget = 5 * time.Second
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), budget)
 	defer cancel()
 	if err := f(ctx, c); err != nil {
 		a.t.Fatal(err)
diff --git a/host/daemon/surface_operator.go b/host/daemon/surface_operator.go
index 0d008f3..37d5db3 100644
--- a/host/daemon/surface_operator.go
+++ b/host/daemon/surface_operator.go
@@ -104,6 +104,11 @@ func (d *Daemon) handleOperatorRoutes(w http.ResponseWriter, r *http.Request, or
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
diff --git a/host/daemon/surface_operator_confirm.go b/host/daemon/surface_operator_confirm.go
new file mode 100644
index 0000000..60e3f8f
--- /dev/null
+++ b/host/daemon/surface_operator_confirm.go
@@ -0,0 +1,156 @@
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
+	"strings"
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
+	passkeys := surfacepage.DecisionPasskeys(keys)
+	wire, err := json.MarshalIndent(struct {
+		Post      string                     `json:"post"`
+		Digest    string                     `json:"digest"`
+		Message   string                     `json:"messageHex"`
+		Namespace string                     `json:"namespace"`
+		Passkeys  []surfacepage.CardPasskeys `json:"passkeys"`
+	}{href, hex.EncodeToString(row.Digest), hex.EncodeToString(m), approval.OperatorActionNamespace, passkeys}, "", "  ")
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
+	if err != nil || (body.WebAuthn != nil && body.WebAuthn[0] != '{') {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send a v1 digest and one key artifact")
+		return
+	}
+	artifact := []byte(body.WebAuthn)
+	if body.SSHSIG != nil {
+		if strings.HasPrefix(strings.TrimSpace(*body.SSHSIG), "{") {
+			writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send an armored SSH signature in sshsig")
+			return
+		}
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
+	served, pathErr := store.CanonicalPath(d.credentialsPath)
+	if pathErr != nil {
+		d.operatorError(w, pathErr)
+		return
+	}
+	if err != nil || a.Kind != "start" || !a.Options.Headless || a.ID != id || a.Episode != row.Episode || a.Kind != row.Kind || a.Store != served ||
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
index 0000000..3fc407d
--- /dev/null
+++ b/host/daemon/surface_operator_confirm_test.go
@@ -0,0 +1,718 @@
+package daemon
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"database/sql"
+	"encoding/base64"
+	"encoding/hex"
+	"encoding/json"
+	"errors"
+	"html"
+	"io"
+	"net/http"
+	"net/http/httptest"
+	"net/url"
+	"os"
+	"sync"
+	"sync/atomic"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/agentrun"
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
+	f := newOperatorHTTP(t, time.Minute)
+	testOperatorConfirm(t, f)
+}
+
+type confirmEvidence struct {
+	row, counters, credentials, world  [32]byte
+	actions, operations, countersCount int
+	fixtureCount                       int64
+}
+
+// Observe raw durable state, so logical expiry and replay cannot mask a write.
+func (f *operatorHTTP) confirmationSnapshot(db *sql.DB, id string) confirmEvidence {
+	f.t.Helper()
+	var out confirmEvidence
+	var row, counters string
+	var effects, starts, sessions int
+	err := db.QueryRowContext(boundedTestContext(f.t), `SELECT
+		(SELECT json_array(id, episode_id, kind, state, hex(review), hex(digest), hex(action), hex(challenge),
+			created_at, expires_at, generation, owner, approved_by, hex(artifact), terminal_at, v3_id)
+			FROM operator_actions WHERE id=?),
+		(SELECT json_group_array(json_array(name, sign_count)) FROM
+			(SELECT name, sign_count FROM webauthn_credentials ORDER BY name)),
+		(SELECT count(*) FROM webauthn_credentials),
+		coalesce((SELECT sign_count FROM webauthn_credentials WHERE name='fixture-passkey'), 0),
+		(SELECT count(*) FROM operator_actions),
+		(SELECT count(*) FROM operator_action_steps WHERE kind='operation'),
+		(SELECT count(*) FROM operator_action_steps WHERE kind<>'operation'),
+		(SELECT count(*) FROM episode_starts), (SELECT count(*) FROM session_credentials)`, id).
+		Scan(&row, &counters, &out.countersCount, &out.fixtureCount, &out.actions, &out.operations, &effects, &starts, &sessions)
+	if err != nil || effects != 0 || starts != 0 || sessions != 0 {
+		f.t.Fatalf("confirmation crossed the effect boundary: effects=%d starts=%d sessions=%d (%v)", effects, starts, sessions, err)
+	}
+	if _, err := os.Lstat(f.fake); !os.IsNotExist(err) {
+		f.t.Fatalf("unexpected fake child marker: %v", err)
+	}
+	hash := func(path string) [32]byte {
+		f.t.Helper()
+		b, err := os.ReadFile(path)
+		if err != nil {
+			f.t.Fatal(err)
+		}
+		return sha256.Sum256(b)
+	}
+	out.row, out.counters = sha256.Sum256([]byte(row)), sha256.Sum256([]byte(counters))
+	out.credentials, out.world = hash(f.d.credentialsPath), hash(f.d.cfg.DBPath)
+	return out
+}
+
+func (f *operatorHTTP) confirmSQL(query string, args ...any) {
+	f.t.Helper()
+	// These are fixture mutations. Confirmation still uses the normal writer DSN.
+	dsn := url.URL{Scheme: "file", Path: f.d.credentialsPath, RawQuery: "_pragma=synchronous%28OFF%29"}
+	db, err := sql.Open("sqlite", dsn.String())
+	if err != nil {
+		f.t.Fatal(err)
+	}
+	defer db.Close()
+	if _, err = db.ExecContext(boundedTestContext(f.t), query, args...); err != nil {
+		f.t.Fatal(err)
+	}
+}
+
+func testOperatorConfirm(t *testing.T, f *operatorHTTP) {
+	var lockCheck atomic.Bool
+	var expireBetweenReads atomic.Int64
+	f.d.approveNow = func() int64 {
+		if lockCheck.Load() && f.d.humanMu.TryLock() {
+			f.d.humanMu.Unlock()
+			f.t.Error("confirm reached its store read outside the human writer lock")
+		}
+		if expireBetweenReads.Load() > 0 && expireBetweenReads.Add(-1) == 0 {
+			return f.now.Load() + 1
+		}
+		return f.now.Load()
+	}
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
+	dbURL := url.URL{Scheme: "file", Path: f.d.credentialsPath, RawQuery: "mode=ro"}
+	evidenceDB, err := sql.Open("sqlite", dbURL.String())
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer evidenceDB.Close()
+	before := f.confirmationSnapshot(evidenceDB, s.ID)
+	post := func(body []byte) *httptest.ResponseRecorder {
+		return f.do("POST", s.Href+"/confirm", body, func(r *http.Request) { r.Host = "localhost:7644" })
+	}
+	checkRefusal := func(r *httptest.ResponseRecorder, status int, class string) {
+		f.t.Helper()
+		requireRefusal(f.t, class, r, status, class)
+		for k := range r.Header() {
+			if k == "Location" || strings.HasPrefix(k, "Access-Control-") {
+				f.t.Fatalf("confirmation refusal emitted %s", k)
+			}
+		}
+		if f.confirmationSnapshot(evidenceDB, s.ID) != before {
+			f.t.Fatalf("confirmation refusal %s changed the row, counter, or effect evidence", class)
+		}
+	}
+	refuse := func(b []byte, status int, class string) { checkRefusal(post(b), status, class) }
+	f.sub("Transport", func(t *testing.T) {
+		for _, tc := range []struct {
+			key, value, class string
+			code              int
+		}{
+			{"Origin", "", "CrossOrigin", 403}, {"Origin", "null", "CrossOrigin", 403}, {"Origin", opOrigin + ".invalid", "CrossOrigin", 403},
+			{"Sec-Fetch-Site", "same-site", "CrossOrigin", 403}, {"Sec-Fetch-Site", "Same-Origin", "CrossOrigin", 403},
+			{"Sec-Fetch-Site", "same-origin, cross-site", "CrossOrigin", 403}, {"Content-Type", "text/plain", "UnsupportedMediaType", 415},
+			{"Authorization", "Bearer y", "Forbidden", 403}, {"Cookie", "y=y", "Forbidden", 403},
+		} {
+			r := f.do("POST", s.Href+"/confirm", signed, func(r *http.Request) {
+				if tc.value == "" {
+					r.Header.Del(tc.key)
+				} else {
+					r.Header.Set(tc.key, tc.value)
+				}
+			})
+			checkRefusal(r, tc.code, tc.class)
+		}
+		checkRefusal(f.do("POST", s.Href+"/confirm", signed, func(r *http.Request) { r.Host = "other.invalid:7644" }), 421, "MisdirectedRequest")
+		for _, header := range []string{"Origin", "Sec-Fetch-Site", "Content-Type"} {
+			code, class := 403, "CrossOrigin"
+			if header == "Content-Type" {
+				code, class = 415, "UnsupportedMediaType"
+			}
+			checkRefusal(f.do("POST", s.Href+"/confirm", signed, func(r *http.Request) {
+				r.Header.Add(header, r.Header.Get(header))
+			}), code, class)
+		}
+		req := fixtureConfirmRequest(s.Href+"/confirm", bytes.NewReader(signed))
+		spy := &confirmReadSpy{ReadCloser: req.Body}
+		req.Host, req.Body = "localhost:7644.evil", spy
+		w := httptest.NewRecorder()
+		f.handler.ServeHTTP(w, req)
+		checkRefusal(w, 421, "MisdirectedRequest")
+		if spy.reads != 0 {
+			t.Fatalf("foreign Host read a valid proof body %d times", spy.reads)
+		}
+		for _, method := range []string{"GET", "OPTIONS", "PUT"} {
+			checkRefusal(f.do(method, s.Href+"/confirm", signed, func(r *http.Request) {
+				if method == "OPTIONS" {
+					r.Header.Set("Origin", opOrigin)
+					r.Header.Set("Access-Control-Request-Method", "POST")
+					r.Header.Set("Access-Control-Request-Headers", "content-type")
+				}
+			}), 405, "MethodNotAllowed")
+		}
+		padded := append(bytes.Clone(signed), bytes.Repeat([]byte(" "), (20<<10)+1-len(signed))...)
+		refuse(padded, 413, "TooLarge")
+		for _, path := range []string{strings.TrimSuffix(s.Href, s.ID) + strings.ToUpper(s.ID),
+			s.Href + "a", strings.TrimSuffix(s.Href, s.ID) + s.ID[:len(s.ID)-1], s.Href + "%0a",
+			operatorPath + "//" + s.ID, operatorPath + "/x/../" + s.ID, "/surface/x/../operator-actions/" + s.ID,
+			strings.Replace(s.Href, "/oa-", "/%6fa-", 1),
+		} {
+			checkRefusal(f.do("POST", path+"/confirm", signed, nil), 400, "BadRequest")
+		}
+		checkRefusal(f.do("POST", s.Href+"/confirm?next=elsewhere", signed, nil), 400, "BadRequest")
+		checkRefusal(f.do("POST", s.Href+"/confirm/", signed, nil), 400, "BadRequest")
+		checkRefusal(f.do("POST", s.Href+"/confirm/x", signed, nil), 404, "NotFound")
+	})
+	f.sub("Proof", func(t *testing.T) {
+		for _, body := range [][]byte{[]byte("y"), opJSON(map[string]any{"v": 1, "digest": d}),
+			opJSON(map[string]any{"v": 1, "digest": d, "y": true}), append([]byte(`{"v":1,`), passkey[1:]...), append([]byte(`{"sshsig":null,`), passkey[1:]...)} {
+			refuse(body, 400, "BadRequest")
+		}
+		refuse(bytes.Replace(signed, []byte(`"v":1`), []byte(`"v":2`), 1), 400, "BadRequest")
+		refuse(approveBody(d, []byte("null")), 400, "BadRequest")
+		refuse(sshBody(d, key.Assert(opOrigin, "localhost", challenge, 2)), 400, "BadRequest")
+		refuse(sshBody(d, []byte("malformed")), 422, "NotApproved")
+		refuse(sshBody(strings.Repeat("0", 64), signer.Sign(approval.OperatorActionNamespace, m)), 409, "Conflict")
+		other := nonce
+		other[0] ^= 1
+		refuse(sig(signer, approval.OperatorActionNamespace, opMessage(t, digest, other, row.ExpiresAt)), 422, "NotApproved")
+	})
+	f.sub("Authenticity", func(t *testing.T) {
+		refuse(sig(approvaltest.NewSigner(47), approval.OperatorActionNamespace, m), 422, "NotApproved")
+		forged := bytes.Clone(m)
+		forged[0] ^= 1
+		refuse(sig(signer, approval.OperatorActionNamespace, forged), 422, "NotApproved")
+		refuse(sig(signer, approval.VerbNamespace, m), 422, "NotApproved")
+		old, err := approval.VerbMessage(digest, nonce, row.ExpiresAt)
+		if err != nil {
+			t.Fatal(err)
+		}
+		refuse(sig(signer, approval.VerbNamespace, old), 422, "NotApproved")
+		oldMint, err := approval.SignedMessage(digest, nonce, row.ExpiresAt)
+		if err != nil {
+			t.Fatal(err)
+		}
+		refuse(sig(signer, approval.Namespace, oldMint), 422, "NotApproved")
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
+		f.now.Store(row.ExpiresAt - 1)
+		expireBetweenReads.Store(2)
+		refuse(signed, 409, "Expired")
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
+			before = f.confirmationSnapshot(evidenceDB, s.ID)
+		}
+		change(row.ExpiresAt + 1)
+		refuse(sig(signer, approval.OperatorActionNamespace, opMessage(t, digest, nonce, row.ExpiresAt+1)), 409, "Expired")
+		change(row.ExpiresAt)
+	})
+	f.sub("RevocationCounter", func(t *testing.T) {
+		for _, count := range []uint32{0, 1} {
+			refuse(approveBody(d, key.Assert(opOrigin, "localhost", challenge, count)), 422, "NotApproved")
+		}
+		revoked := approvaltest.NewSigner(46)
+		f.enrolSSH("revoked-ssh", revoked)
+		f.write(func(ctx context.Context, c *store.Credentials) error {
+			return c.RevokeApprover(ctx, "revoked-ssh", apT0)
+		})
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		refuse(sig(revoked, approval.OperatorActionNamespace, m), 422, "NotApproved")
+	})
+	f.sub("UnsupportedKinds", func(t *testing.T) {
+		for i, kind := range []string{"start", "finish", "foreign"} {
+			a, err := approval.DecodeOperatorAction(row.Action)
+			if err != nil {
+				t.Fatal(err)
+			}
+			a.ID, a.Episode, a.Kind = "oa-"+opRequest(i+30).OperationID, "other-"+kind, kind
+			a.Options.Headless = false
+			if kind == "foreign" {
+				a.Kind, a.Options.Headless = "start", true
+			}
+			if kind == "finish" {
+				a.Task = ""
+			}
+			a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+			var stored store.OperatorActionRow
+			f.write(func(ctx context.Context, c *store.Credentials) error {
+				var err error
+				stored, err = c.RequestOperatorAction(ctx, store.OperatorRequest{
+					Action: a, OperationID: opRequest(i + 30).OperationID, Challenge: nonce, CreatedAt: apT0, ExpiresAt: apT0 + 600,
+				}, apT0)
+				return err
+			})
+			if kind == "foreign" {
+				a.Store += ".foreign"
+				wire, err := approval.CanonicalOperatorAction(a)
+				if err != nil {
+					t.Fatal(err)
+				}
+				digest := sha256.Sum256(wire)
+				f.confirmSQL("UPDATE operator_actions SET action=?, digest=? WHERE id=?", wire, digest[:], a.ID)
+				stored.Digest = digest[:]
+			}
+			before = f.confirmationSnapshot(evidenceDB, s.ID)
+			body := sshBody(hex.EncodeToString(stored.Digest), signer.Sign(approval.OperatorActionNamespace,
+				opMessage(t, [32]byte(stored.Digest), nonce, stored.ExpiresAt)))
+			checkRefusal(f.do("POST", operatorPath+"/"+a.ID+"/confirm", body, nil), 409, "Conflict")
+		}
+	})
+	f.sub("StoredIdentity", func(t *testing.T) {
+		for _, edit := range []struct{ change, restore string }{
+			{"episode_id='drift'", "episode_id=" + "'" + row.Episode + "'"},
+			{"kind='finish'", "kind='start'"},
+			{"digest=zeroblob(32)", "digest=X'" + d + "'"},
+		} {
+			f.confirmSQL("UPDATE operator_actions SET "+edit.change+" WHERE id=?", s.ID)
+			before = f.confirmationSnapshot(evidenceDB, s.ID)
+			proof := signed
+			if edit.change == "digest=zeroblob(32)" {
+				proof = sshBody(strings.Repeat("0", 64), signer.Sign(approval.OperatorActionNamespace, m))
+			}
+			refuse(proof, 409, "Conflict")
+			f.confirmSQL("UPDATE operator_actions SET "+edit.restore+" WHERE id=?", s.ID)
+		}
+		a, err := approval.DecodeOperatorAction(row.Action)
+		if err != nil {
+			t.Fatal(err)
+		}
+		a.ID = "oa-" + strings.Repeat("b", 32)
+		other, err := approval.CanonicalOperatorAction(a)
+		if err != nil {
+			t.Fatal(err)
+		}
+		different := sha256.Sum256(other)
+		f.confirmSQL("UPDATE operator_actions SET action=?, digest=? WHERE id=?", other, different[:], s.ID)
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		refuse(sshBody(hex.EncodeToString(different[:]), signer.Sign(approval.OperatorActionNamespace,
+			opMessage(t, different, nonce, row.ExpiresAt))), 409, "Conflict")
+		f.confirmSQL("UPDATE operator_actions SET action=?, digest=? WHERE id=?", row.Action, row.Digest, s.ID)
+		f.confirmSQL("PRAGMA ignore_check_constraints=ON; UPDATE operator_actions SET challenge=zeroblob(31) WHERE id=?", s.ID)
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		refuse(signed, 409, "Conflict")
+		card := f.do("GET", s.Href, nil, nil)
+		if card.Code != 200 || strings.Contains(card.Body.String(), "operator-ceremony") {
+			t.Fatalf("card offered proof for an invalid challenge: %d %s", card.Code, card.Body)
+		}
+		f.confirmSQL("UPDATE operator_actions SET challenge=? WHERE id=?", row.Challenge, s.ID)
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+	})
+	f.sub("Cancelled", func(t *testing.T) {
+		other := f.request(opRequest(29))
+		f.write(func(ctx context.Context, c *store.Credentials) error {
+			return c.CancelUnstartedOperatorAction(ctx, other.ID, apT0+1)
+		})
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		cancelled := f.row(other.ID)
+		checkRefusal(f.do("POST", other.Href+"/confirm", sshBody(hex.EncodeToString(cancelled.Digest), []byte("malformed")), nil), 409, "Conflict")
+	})
+	f.sub("MissingAndWriter", func(t *testing.T) {
+		path := f.d.credentialsPath
+		if err := os.Rename(path, path+".offline"); err != nil {
+			t.Fatal(err)
+		}
+		r := post(signed)
+		if err := os.Rename(path+".offline", path); err != nil {
+			t.Fatal(err)
+		}
+		checkRefusal(r, 503, "Disconnected")
+		checkRefusal(f.do("POST", operatorPath+"/oa-"+strings.Repeat("0", 32)+"/confirm", signed, nil), 404, "NotFound")
+		held, err := store.OpenCredentials(f.d.credentialsPath)
+		if err != nil {
+			t.Fatal(err)
+		}
+		refuse(signed, 503, "Busy")
+		if err := held.Close(); err != nil {
+			t.Fatal(err)
+		}
+		f.confirmSQL("ALTER TABLE operator_actions RENAME COLUMN action TO unreadable_action")
+		readFault := f.evidence()
+		f.refuse(readFault, post(signed), 503, "Disconnected")
+		f.confirmSQL("ALTER TABLE operator_actions RENAME COLUMN unreadable_action TO action")
+		after := f.confirmationSnapshot(evidenceDB, s.ID)
+		if before.row != after.row || before.counters != after.counters || before.countersCount != after.countersCount {
+			t.Fatal("read failure changed the pending action or counters")
+		}
+		before = after
+		f.confirmSQL("UPDATE approvers SET enrolled_at='unreadable' WHERE name='fixture-passkey'")
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		refuse(passkey, 503, "Disconnected")
+		f.confirmSQL("UPDATE approvers SET enrolled_at=? WHERE name='fixture-passkey'", apT0-10)
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+	})
+	f.sub("Rollback", func(t *testing.T) {
+		f.confirmSQL(`CREATE TRIGGER interrupt_confirm BEFORE UPDATE OF approved_by ON operator_actions
+			BEGIN SELECT RAISE(ABORT, 'fixture: interrupted after passkey use'); END`)
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		refuse(passkey, 503, "Disconnected")
+		f.confirmSQL("DROP TRIGGER interrupt_confirm")
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		f.confirmSQL(`CREATE TRIGGER concurrent_counter BEFORE UPDATE OF last_used_at ON approvers
+			BEGIN UPDATE webauthn_credentials SET sign_count=2 WHERE name='fixture-passkey'; END`)
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		refuse(passkey, 422, "NotApproved")
+		f.confirmSQL("DROP TRIGGER concurrent_counter")
+		f.confirmSQL(`CREATE TRIGGER concurrent_revocation BEFORE UPDATE OF last_used_at ON approvers
+			WHEN NEW.name='fixture-passkey' BEGIN
+				UPDATE approvers SET revoked_at=NEW.last_used_at WHERE name=NEW.name;
+				SELECT RAISE(IGNORE);
+			END`)
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		refuse(passkey, 422, "NotApproved")
+		f.confirmSQL("DROP TRIGGER concurrent_revocation")
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+	})
+	f.sub("Ceremony", func(t *testing.T) {
+		requireOperatorCeremony(t, f.do("GET", s.Href, nil, nil), s, d, m, key)
+		f.now.Store(row.CreatedAt - 1)
+		card := f.do("GET", s.Href, nil, nil)
+		if card.Code != 200 || strings.Contains(card.Body.String(), "operator-ceremony") {
+			t.Fatalf("card offered a ceremony before issue: %d %s", card.Code, card.Body)
+		}
+		f.now.Store(apT0)
+		f.confirmSQL("UPDATE approvers SET enrolled_at='unreadable' WHERE name='fixture-passkey'")
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		checkRefusal(f.do("GET", s.Href, nil, nil), 503, "Disconnected")
+		f.confirmSQL("UPDATE approvers SET enrolled_at=? WHERE name='fixture-passkey'", apT0-10)
+		f.confirmSQL("UPDATE operator_actions SET created_at=-2, expires_at=-1 WHERE id=?", s.ID)
+		f.now.Store(-2)
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		checkRefusal(f.do("GET", s.Href, nil, nil), 503, "Disconnected")
+		f.confirmSQL("UPDATE operator_actions SET created_at=?, expires_at=? WHERE id=?", row.CreatedAt, row.ExpiresAt, s.ID)
+		f.now.Store(apT0)
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+	})
+	expect := func(r *httptest.ResponseRecorder, id, result, key string) {
+		t.Helper()
+		var out operatorConfirmationAnswer
+		if err := json.Unmarshal(r.Body.Bytes(), &out); err != nil || r.Code != 200 || out.Result != result || out.ApprovedBy != key ||
+			out.ID != id || out.State != "confirmed" {
+			t.Fatalf("confirmation %s: %d %s (%v)", result, r.Code, r.Body, err)
+		}
+	}
+	f.sub("ConcurrentReplay", func(t *testing.T) {
+		f.confirmSQL("DELETE FROM operator_runner_state")
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		f.now.Store(row.ExpiresAt - 1)
+		lockCheck.Store(true)
+		defer lockCheck.Store(false)
+		// Two handlers start together. Both must return; only one may advance the counter.
+		var wg sync.WaitGroup
+		var replies [2]*httptest.ResponseRecorder
+		ready := make(chan struct{})
+		for i := range replies {
+			wg.Add(1)
+			go func() {
+				defer wg.Done()
+				<-ready
+				replies[i] = post(passkey)
+			}()
+		}
+		close(ready)
+		wg.Wait()
+		if !f.d.humanMu.TryLock() {
+			t.Fatal("confirm retained the writer lock after responding")
+		}
+		f.d.humanMu.Unlock()
+		if replies[0] == nil || replies[1] == nil {
+			t.Fatal("confirmation request did not return")
+		}
+		if strings.Contains(replies[0].Body.String(), "already-confirmed") {
+			replies[0], replies[1] = replies[1], replies[0]
+		}
+		expect(replies[0], s.ID, "confirmed", "webauthn:fixture-passkey")
+		expect(replies[1], s.ID, "already-confirmed", "webauthn:fixture-passkey")
+		before = f.confirmationSnapshot(evidenceDB, s.ID)
+		expect(post(passkey), s.ID, "already-confirmed", "webauthn:fixture-passkey")
+		if f.confirmationSnapshot(evidenceDB, s.ID) != before || before.fixtureCount != 2 {
+			t.Fatal("replay changed durable evidence or advanced the counter twice")
+		}
+		refuse(signed, 409, "Conflict")
+		after := f.row(s.ID)
+		if after.State != "confirmed" || after.Generation != row.Generation || after.Owner != "" ||
+			after.ApprovedBy != "webauthn:fixture-passkey" || f.evidence().operations != before.operations {
+			t.Fatalf("handler took a runner step or failed to make the row claimable: %+v", after)
+		}
+		f.h.SeenAt = f.now.Load()
+		f.h.Owner.Generation = 0 // Simulate a fresh runner registering after the POST.
+		f.publish()
+		f.write(func(ctx context.Context, c *store.Credentials) error {
+			cas := store.OperatorCAS{ID: s.ID, Digest: digest, Generation: after.Generation}
+			return c.ClaimOperatorAction(ctx, cas, f.h.Owner, f.h.SeenAt)
+		})
+		if got := f.row(s.ID); got.State != "claimed" || got.Owner != f.h.Owner.Nonce {
+			t.Fatalf("a subsequent ready runner could not claim once: %+v", got)
+		}
+		f.evidence()
+		f.now.Add(1)
+	})
+	f.sub("SSHSIG", func(t *testing.T) {
+		second := f.request(opRequest(2))
+		r2 := f.row(second.ID)
+		body := sshBody(hex.EncodeToString(r2.Digest), signer.Sign(approval.OperatorActionNamespace,
+			opMessage(t, [32]byte(r2.Digest), [32]byte(r2.Challenge), r2.ExpiresAt)))
+		expect(f.do("POST", second.Href+"/confirm", body, nil), second.ID, "confirmed", "ssh:fixture-ssh")
+		if !f.d.humanMu.TryLock() {
+			t.Fatal("successful response retained the writer lock")
+		}
+		f.d.humanMu.Unlock()
+		if card := f.do("GET", second.Href, nil, nil); card.Code != 200 || strings.Contains(card.Body.String(), "operator-ceremony") {
+			t.Fatalf("confirmed card exposed an active ceremony: %d %s", card.Code, card.Body)
+		}
+		for _, text := range []string{f.do("GET", second.StatusHref, nil, nil).Body.String(), f.do("GET", "/surface/requests.json", nil, nil).Body.String(), f.log.String()} {
+			if strings.Contains(text, "BEGIN SSH SIGNATURE") || strings.Contains(text, "alert(1)") || strings.Contains(text, hex.EncodeToString(m)) || strings.Contains(text, "authenticatorData") {
+				t.Fatalf("private task or artifact reached a routine read or log: %s", text)
+			}
+		}
+	})
+}
+
+type confirmReadSpy struct {
+	io.ReadCloser
+	reads int
+}
+
+func (b *confirmReadSpy) Read(p []byte) (int, error) {
+	b.reads++
+	return b.ReadCloser.Read(p)
+}
+
+func fixtureConfirmRequest(target string, body io.Reader) *http.Request {
+	req := httptest.NewRequest(http.MethodPost, opOrigin+target, body)
+	req.Header.Set("Origin", opOrigin)
+	req.Header.Set("Content-Type", "application/json")
+	req.Header.Set("Sec-Fetch-Site", "same-origin")
+	return req
+}
+
+func requireOperatorCeremony(t *testing.T, card *httptest.ResponseRecorder, s operatorStatus, digest string, message []byte, key *approvaltest.Passkey) {
+	t.Helper()
+	body := card.Body.String()
+	_, after, found := strings.Cut(body, `<pre id="operator-ceremony">`)
+	wire, _, closed := strings.Cut(after, "</pre>")
+	var got struct {
+		Post, Digest, MessageHex, Namespace string
+		Passkeys                            []surfacepage.CardPasskeys
+	}
+	if card.Code != 200 || card.Header().Get("Cache-Control") != "no-store" ||
+		card.Header().Get("Content-Security-Policy") == "" || !found || !closed || strings.Contains(body, "<script>alert(1)</script>") ||
+		!strings.Contains(body, "&lt;script&gt;alert(1)&lt;/script&gt;") ||
+		json.Unmarshal([]byte(html.UnescapeString(wire)), &got) != nil ||
+		got.Post != s.Href+"/confirm" || got.Digest != digest || got.MessageHex != hex.EncodeToString(message) ||
+		got.Namespace != approval.OperatorActionNamespace || len(got.Passkeys) != 1 ||
+		got.Passkeys[0].Origin != opOrigin || got.Passkeys[0].RP != "localhost" ||
+		len(got.Passkeys[0].Allow) != 1 || got.Passkeys[0].Allow[0] != base64.RawURLEncoding.EncodeToString(key.CredentialID()) {
+		t.Fatalf("card changed the task or key ceremony: %d %s", card.Code, body)
+	}
+}
+
+func TestOperatorConfirmProjection(t *testing.T) {
+	f := newOperatorSurface(t, nil, time.Minute)
+	raw := f.postOperator(opRequest(1))
+	var s operatorStatus
+	if err := json.Unmarshal(raw.Body.Bytes(), &s); err != nil || raw.Code != 200 || s.ID == "" {
+		t.Fatalf("projection request fixture: %d %s", raw.Code, raw.Body)
+	}
+	read := func(path string) *httptest.ResponseRecorder { return f.serve("GET", path, apLoopback, nil) }
+	original := f.row(s.ID).Digest
+	f.confirmSQL("UPDATE operator_actions SET digest=zeroblob(32) WHERE id=?", s.ID)
+	f.refuse(f.evidence(), read(s.Href), 409, "Conflict")
+	f.confirmSQL("UPDATE operator_actions SET digest=? WHERE id=?", original, s.ID)
+	if r := read(s.Href); r.Code != 200 || !strings.Contains(r.Body.String(), "operator-ceremony") {
+		t.Fatalf("unmodified review did not render: %d %s", r.Code, r.Body)
+	}
+
+	// Publish the shape of an S6 launch receipt in the migrated fixture. Only
+	// the identity read is live. Polling must have no write or child effect.
+	executable := "sha256:" + strings.Repeat("a", 64)
+	generation := f.row(s.ID).Generation
+	f.confirmSQL("UPDATE operator_actions SET action=NULL, state='failed', terminal_at=? WHERE id=?", apT0, s.ID)
+	receipt := store.OperatorLaunchReceipt{Run: s.ID, Executable: executable, Status: "spawn-refused"}
+	f.confirmSQL(`INSERT INTO operator_action_steps(action_id,kind,identity,phase,version,owner,receipt,created_at)
+		VALUES(?, 'launch', ?, 'observed', 2, 'fixture', ?, ?)`, s.ID, s.ID+":"+strings.Repeat("a", 64), opJSON(receipt), apT0)
+	status := func(want string) {
+		t.Helper()
+		before := f.row(s.ID)
+		r := read(s.StatusHref)
+		var got operatorStatus
+		if err := json.Unmarshal(r.Body.Bytes(), &got); err != nil || r.Code != 200 || got.State != want ||
+			got.ReportHref != "/v1/episodes/"+s.Episode+"/report" || got.ID != s.ID ||
+			!bytes.Equal(opJSON(f.row(s.ID)), opJSON(before)) || before.Generation != generation {
+			t.Fatalf("status readback want %s: %d %s", want, r.Code, r.Body)
+		}
+	}
+	status("spawn-refused")
+	receipt.PID = os.Getpid()
+	birth, live := agentrun.ProcessIdentity(receipt.PID)
+	if !live || birth == "" {
+		t.Skip("live birth unavailable on this platform")
+	}
+	f.confirmSQL("UPDATE operator_actions SET state='succeeded' WHERE id=?", s.ID)
+	receipt.Birth, receipt.Status = birth, "launched"
+	f.confirmSQL("UPDATE operator_action_steps SET receipt=?, version=3 WHERE kind='launch'", opJSON(receipt))
+	status("launched")
+	receipt.Birth += "-foreign"
+	f.confirmSQL("UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", opJSON(receipt))
+	status("uncertain")
+}
+
+func TestOperatorConfirmationAdmissionNeighbours(t *testing.T) {
+	f := newOperatorSurface(t, nil, time.Minute)
+	checked := 0
+	f.d.approveNow = func() int64 {
+		checked++
+		if f.d.humanMu.TryLock() {
+			f.d.humanMu.Unlock()
+			t.Error("request reached its store read outside the human writer lock")
+		}
+		return apT0
+	}
+	bad := opRequest(1)
+	bad.CatalogueGeneration = 0
+	before := f.evidence()
+	f.refuse(before, f.postOperator(bad), 400, "BadRequest")
+	if r := f.postOperator(opRequest(1)); r.Code != 200 || checked < 2 {
+		t.Fatalf("request positive control did not exercise lock: %d (%d checks)", r.Code, checked)
+	}
+	f.d.approveNow = f.now.Load
+	ctx := boundedTestContext(t)
+	ro, err := store.OpenCredentialsReadOnly(f.d.credentialsPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer ro.Close()
+	if _, _, err = f.d.operatorReview(ctx, ro, "sol", apT0); err != nil {
+		t.Fatalf("catalogue control: %v", err)
+	}
+	f.confirmSQL("UPDATE operator_runner_state SET catalogue_generation=0")
+	if h, _, err := f.d.operatorReview(ctx, ro, "sol", apT0); !errors.Is(err, store.ErrOperatorRunner) || h.CatalogueGeneration != 0 {
+		t.Fatalf("invalid generation crossed review: %v, %v", h, err)
+	}
+	bar := "/surface/verb.json?q=" + url.QueryEscape("episode start ep-2 --preset sol --launch --headless")
+	if r := f.serve("GET", bar, apLoopback, nil); r.Code != 200 || !strings.Contains(r.Body.String(), `"outcome":"terminal"`) {
+		t.Fatalf("invalid generation reached the bar review: %d %s", r.Code, r.Body)
+	}
+	summary := func() surfacepage.HomeSummary {
+		t.Helper()
+		r := f.serve("GET", "/surface/home-summary.json", apLoopback, nil)
+		var out surfacepage.HomeSummary
+		if r.Code != 200 || json.Unmarshal(r.Body.Bytes(), &out) != nil || len(out.Items) != 1 || out.Items[0].Label != "Needs you" {
+			t.Fatalf("home summary: %d %s", r.Code, r.Body)
+		}
+		return out
+	}
+	if got := summary().Items[0].Value; got != "1" {
+		t.Fatalf("healthy pending read did not show one row: %s", got)
+	}
+	f.confirmSQL("UPDATE operator_actions SET expires_at='unreadable' WHERE episode_id='ep-1'")
+	before = f.evidence()
+	r := f.serve("GET", "/surface/requests.json", apLoopback, nil)
+	var requests surfacepage.RequestsJSON
+	if r.Code != 200 || json.Unmarshal(r.Body.Bytes(), &requests) != nil || requests.Unavailable != "" || requests.OperatorActionsUnavailable == "" {
+		t.Fatalf("fixture must fail only the operator pending read: %d %s", r.Code, r.Body)
+	}
+	if got := summary().Items[0].Value; got != "unknown" || before != f.evidence() {
+		t.Fatalf("unavailable operator read looked like an empty queue or changed evidence: %s", got)
+	}
+}
diff --git a/host/daemon/surface_operator_read.go b/host/daemon/surface_operator_read.go
index 0bdf31c..3ebd00a 100644
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
@@ -85,13 +89,10 @@ func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id s
 	if status.Href != "" && row.Kind == "start" && (row.State == "running" || row.State == "succeeded" || row.State == "uncertain" || row.State == "failed") {
 		launch, err := ro.ReadOperatorLaunch(ctx, id)
 		if err == nil {
-			switch {
-			case launch.NoSpawn:
+			if launch.NoSpawn {
 				status.State = "spawn-refused"
-			case launch.Phase == "reserved" || launch.Receipt.Status == "launched":
-				status.State = "uncertain"
-			default:
-				status.State = launch.Receipt.Status
+			} else {
+				status.State = launch.ProcessStatus(agentrun.ProcessIdentity)
 			}
 			if status.ReportHref == "" {
 				status.ReportHref = "/v1/episodes/" + row.Episode + "/report"
@@ -109,8 +110,8 @@ func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id s
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
@@ -120,6 +121,13 @@ func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id s
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
index cc40d62..f536bd6 100644
--- a/host/daemon/surface_operator_test.go
+++ b/host/daemon/surface_operator_test.go
@@ -19,6 +19,7 @@ import (
 	"sync"
 	"sync/atomic"
 	"testing"
+	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/store"
@@ -55,17 +56,21 @@ type operatorHTTP struct {
 	log     operatorLog
 }
 
-func newOperatorHTTP(t *testing.T) *operatorHTTP {
+func newOperatorHTTP(t *testing.T, budget ...time.Duration) *operatorHTTP {
 	t.Helper()
-	f := newOperatorSurface(t, nil)
+	f := newOperatorSurface(t, nil, budget...)
 	f.server = httptest.NewServer(f.handler)
 	t.Cleanup(f.server.Close)
 	return f
 }
 
-func newOperatorSurface(t *testing.T, edit func(*store.OperatorStartReview)) *operatorHTTP {
+func newOperatorSurface(t *testing.T, edit func(*store.OperatorStartReview), budget ...time.Duration) *operatorHTTP {
 	t.Helper()
 	f := &operatorHTTP{approvalRig: newApprovalRig(t)}
+	if len(budget) != 0 {
+		f.approvalRig.writeBudget = budget[0]
+		f.d.commitBudget, f.d.readDeadline = budget[0], budget[0]
+	}
 	f.now.Store(apT0)
 	f.d.approveNow = f.now.Load
 	f.d.errLog = &f.log
@@ -320,7 +325,7 @@ func TestOperatorAdmissionHTTP(t *testing.T) {
 			refuse(path, "GET", nil, 400, "BadRequest", nil)
 		}
 		refuse(operatorPath+id+"/status.json/x", "GET", nil, 404, "NotFound", nil)
-		refuse(operatorPath+id+"/confirm", "POST", valid, 404, "NotFound", nil)
+		refuse(operatorPath+id+"/confirm/x", "POST", valid, 404, "NotFound", nil)
 	})
 	f.sub("StrictJSONAndLimits", func(t *testing.T) {
 		invalid := [][]byte{nil, []byte("[]"), []byte(`{"argv":"start"}`), []byte(`{"argv":[42]}`), append(bytes.Clone(valid), []byte(" {}")...),
diff --git a/host/store/operator_launch.go b/host/store/operator_launch.go
index 206672d..84f62df 100644
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
diff --git a/host/store/operator_launch_process_test.go b/host/store/operator_launch_process_test.go
new file mode 100644
index 0000000..2d733b4
--- /dev/null
+++ b/host/store/operator_launch_process_test.go
@@ -0,0 +1,33 @@
+package store
+
+import "testing"
+
+func TestOperatorLaunchProcessStatus(t *testing.T) {
+	for _, tc := range []struct {
+		name, phase, status, birth, want string
+		live, noSpawn, reads             bool
+	}{
+		{"no-spawn", "observed", "", "", "spawn-refused", false, true, false},
+		{"reserved", "reserved", "", "", "uncertain", false, false, false},
+		{"live", "observed", "launched", "same", "launched", true, false, true},
+		{"reused", "observed", "launched", "different", "uncertain", true, false, true},
+		{"gone", "observed", "launched", "same", "uncertain", false, false, true},
+		{"stopped", "observed", "stopped", "", "stopped", false, false, false},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			snapshot := OperatorLaunchSnapshot{Phase: tc.phase, NoSpawn: tc.noSpawn,
+				Receipt: OperatorLaunchReceipt{PID: 42, Birth: "same", Status: tc.status}}
+			read := false
+			got := snapshot.ProcessStatus(func(pid int) (string, bool) {
+				if pid != 42 {
+					t.Fatalf("looked up pid %d", pid)
+				}
+				read = true
+				return tc.birth, tc.live
+			})
+			if got != tc.want || read != tc.reads {
+				t.Fatalf("status %q, read=%v, want %q, read=%v", got, read, tc.want, tc.reads)
+			}
+		})
+	}
+}
````

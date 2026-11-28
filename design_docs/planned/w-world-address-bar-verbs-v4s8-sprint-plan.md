# w-world-address-bar-verbs-v4 — S8 sprint plan

## 0. Scope and outcome

Row 243, World iteration 447, on `codex:gpt-6.2-sol`. The reviewed design is V4 r1, quorum r2 PROCEED 3/3 with every seat present. Apply §4.1 and §4.2 step 5 within the §3 persistence and confirmation law. Design §8 allocates 860 raw and 1,161 calibrated lines. **Outcome A** keeps both milestones in S8 because the complete compiled product fits below 1,400 changed lines, including JS, CSS, HTML and tests. The size and snapshots in §2 are base-relative.

Chrome uses a **FAKE Sol preset**, a **FAKE launcher**, the repository's fake child, a CDP virtual authenticator and temporary migrated stores. E5's testable half runs here. Real-preset eligibility for Mark's approved `V4_PRESET_FILE` is owed before S11 and the attended exit. The planner and executor use fixture paths only. S9–S11 own finish and the V3 PR hand-off. S12 owns interactive browser activation under the existing attended decision.

## 1. Measurements

### 1.1 Base and provenance

- `git rev-parse HEAD` returned `0e32a8d7de0bb524678a3c443400cf05d6ec47eb`. `git diff --stat 26b6947f 0e32a8d7 -- cmd host docs app` returned zero lines. S7's review, inert POST, card/status and Needs-you summary and S7b's proof, refusal, replay and launch projection are present.
- HEAD's app script has 3,899 lines and zero `operator-actions` references. At `app.js:1430`, HEAD's `VERB_NEEDS` excludes `review`. Feeding a schema-valid S7 response to HEAD's `verbOutcome` produces `the verb route answered something this page cannot read`. The address input and current page remain in place and the browser sends no POST.
- `git show -s --oneline 65108353` resolves the #604 merge. `git diff --stat 65108353^1 65108353 -- host/surfacepage/app/{app.js,app.css,index.html}` shows +55/−5 across the three files. The local merge commit is 2026-11-18 09:48:29 UTC. Controller PR metadata records mergedAt 09:48:30 UTC. The post-merge-to-HEAD asset diff is +7,453/−420 across 23 shell assets. For the three shared files the additions/deletions are app.js 2,556/268, app.css 825/151 and index.html 39/1. M2 sits next to #604's document-level Enter/mobile queue guard on the current global keyboard handler, and uses the existing verb sheet above the phone content. Its outcome dispatcher and `go/showVerbPanel` hunks are outside #604's `phoneWidth`, queue repaint, wheel and resize hunks. V4 adds separate card assets. The app CSS and index shell have zero changed lines.
- E5/W25 ran `TestOperatorFakePresetEligibility` with `runPreset check --presets-dir <fixture> --exec-profiles-dir <fixture> p`, using the fake daemon health match. It printed `base dev @ 61e2e46` and `check short: sandboxed: true` in the named plain receipt. A fixture Git SHA can differ between runs. The test compares the resolved SHA and requires the sandboxed check.
- E5/W32 used the same interactive preset with fake model name `gpt-6.2-sol`. Default resolution was interactive. The bar policy refused implicit interactive activation, and the explicit headless form resolved with cap `1m0s`. The review states `Headless codex gpt-6.2-sol, cap 1m0s`. The Go test decoded the signed action field `"headless":true`, checked the canonical argv and a real test SSH signature over its digest, then re-resolved the same signed bindings through the runner. The Node bar test executes the corresponding JS outcome. **PENDING — controller to measure the Chrome panel and signed-database-row observation** with `outside.sh M2`.

### 1.2 Guard and route inventory

The exact anchored test selectors are in the `selectors.sh` fence and `gates.sh` calls below. The stem-matched selectors also run in `accept.sh`. `PAGE_STEMS` includes the golden readers for `verb-route.golden.json` and `requests-confirmations.golden.json`. `DAEMON_GOLDEN` drives the same server wire from the original route suite. No golden was regenerated.

| Gate selector | Behaviour pinned |
|---|---|
| `DAEMON_CORE` and `DAEMON_INVENTORY` | `daemon_test.go` full stem, frozen V1 table, daemon dependency/writer/constructor and alias inventories, Host, methods, route registrations, surface shell/script/model, card/report/read route boundaries and verb JSON bounds. |
| `DAEMON_HEADER_INVENTORY`, `DAEMON_READ_INVENTORY` | HTML security headers, A2/A3/V3 card and asset routes, origins/cookies, passkey, font, follow, workbench, read deadlines and sanitized errors. |
| `DAEMON_APPROVAL` | All existing `surface_operator*_test.go` top-level names, `TestOperatorCardAssets`, verb-outcome and Needs-you golden. Both carried test groups use the actual loopback listener. |
| `PAGE_STEMS`, `PAGE_INVENTORY`, `PAGE_CARD`, `PAGE_BAR` | Every `app*_test.go` test, embedded/app shell parity, script seams, preview bundle, font, confirms/assets/goldens and executable new card/bar cases. |
| `STORE_INVENTORY`, `VERBS_INVENTORY` | `TestProductionContextRoots`, `TestStoreSurfaceTakesContext`, `TestCredentialsMethodsRequireADeadline`, `TestStoreDeadlineGuardDynamic`, `TestOperatorNeverDispatch`, `TestOperatorOtherVerbParity`. |
| `TestWorldBoundaryDependencyAllowlist\|TestBoundaryASTWriteGuard` and `TestCLICredentialWriterSites\|TestSessionRootBudgets` | Boundary write policy, caller/dependency and session root budgets. |
| `VERIFY_INVENTORY`, `BROKER_INVENTORY` | `TestSubprocessSinkGate`, `TestExecGateArmsAreLoadBearing`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestVerifygateTestWritesAreForkLocked`, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`. Native Linux additionally runs `TestKernelRefusesExecOfWriterOpenFile`. Both whole packages run on each stage. |

The final executable page inventories recorded 22 named passes on M1 and 30 on M2 with zero fail or skip in `$S/inventories-final.bbQExN`. The static daemon, store, CLI, boundary, verbs and verifygate groups passed on both stages. The real route, page header and broker listener-dependent groups are UNINFORMATIVE UNDER SANDBOX. Their outside selectors are part of `accept.sh` and `gates.sh`. Section 6 links each raw result.

The existing page CSP has `script-src 'self'`, `connect-src 'self'`, default deny, no script inline permission and an existing style-inline directive inherited by other surfaces. The new card uses external same-origin JS and CSS. The assets have the existing inert asset CSP, `nosniff` and explicit content types. The rendered card keeps the original surface headers, Host gate and template escaping. M1 adds GET asset registrations. M2 uses existing routes for every request and effect.

### 1.3 Baseline and race budgets

`$STATE/base` belongs to the controller. Accept its fail sets only after `tail -n 1 "$STATE/baseline.done"` prints `ALL_DONE for $B`. Plain cmd and daemon use the controller's existing pinned-tool fail sets. The extra pristine surfacepage plain/race and cmd race receipts are under `$S/outside-pristine`. The 11 row-271 CLI tests in `CMD_SKIP` are excluded from every whole local cmd race comparison, with the same exclusion in the base. The three known daemon real-binary reds and the load-sensitive git, broker, registry and service tests keep their existing baselines. An additional failing test needs a separate base/head same-selector receipt before attribution. Timeout truncation, a missing package result and a nonzero rc with no failing test are instrument failures.

`.github/workflows/ci.yml` invokes `verify_go.sh` plain and race legs without `WORLD_CHROME` or `WORLD_REQUIRE_PASSKEY_BROWSER`. That script's per-package race bound is 900 s. The required-browser flag changes a missing dependency from counted skip to failure. CI skips the new CDP test. The controller's commit-time leg outside the sandbox pins Chrome to `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`, requires the passkey browser, and counts actual named PASS/FAIL with zero SKIP. A CI green with a skipped browser case carries zero Chrome credit. The last slow CI merge recorded cmd 511.1 s, daemon 539.8 s, verifygate 491.9 s and broker 422.1 s. The PR lane recorded 289.7/367.3/321.6/280.0 s respectively.

### 1.4 Premises and production caller census

A real socket probe for `127.0.0.1:0` was denied with `Operation not permitted`. All real-listener and Chrome outcomes in this plan are **UNINFORMATIVE UNDER SANDBOX**. The existing S7b carried projection and neighbouring admission tests originally used a recorder. M1 moves those cases to the existing real `httptest` fixture before the controller re-runs all six edits. J10 checks terminal response over the listener before asserting the helper's typed refusal. The app asset is a composed bundle. `TestSurfaceAppPure` pins its exact CommonJS exports and therefore needed a corresponding test update.

| Production caller | Pristine | M1 | M2 | Location |
|---|---:|---:|---:|---|
| `RequestOperatorAction` | 1 | 1 | 1 | `host/daemon/surface_operator.go:231` |
| `JudgeOperatorAction` | 1 | 1 | 1 | `host/daemon/surface_operator_confirm.go:117` |
| Store `ConfirmOperatorAction` | 1 | 1 | 1 | `host/daemon/surface_operator_confirm.go:141` |
| `ClaimOperatorAction` | 1 | 1 | 1 | `cmd/ailang-worldd/operator_preparation.go:115` |

`census-M1.raw` and `census-M2.raw` print `CALLER_CENSUS_OK`, the real pure-approval caller, the existing worker subprocess site and installed runner call as positive controls. No schema, `.ail`, production context root or production subprocess site is added. The M2 test clock is a construction parameter for the existing approval clock. Its nil default uses the system clock.

## 2. Milestones

| Commit | Increment | Cumulative insertions/deletions | Counted paths and unfiltered product | STATUS / DELETED / MANIFEST |
|---|---:|---:|---:|---:|
| M1 | 510 | +484/−26 | 510 / 510 | 11 / 0 / 11 |
| M2 | 759 | +1,231/−38 | 1,269 / 1,269 | 16 / 0 / 16 |

| Milestone | Complete behaviour |
|---|---|
| M1 | Standalone stable card, local asset registration, explicit key gestures and reconnectable read-only status. E15 store transaction drift and E19 synthetic Darwin process-row pins. The prior test controls for J02/J06/J10/J18/J19/J22 are now real-listener checks. A requested action can enter the card through the existing durable link. |
| M2 | Address bar Enter opens the bounded review; Request action sends one private request; Reconnect reuses frozen bytes across reload; the exact allowed card link activates. Add fake preset parity, Chrome bad payload and fake confined start, phone viewport, focus, readback and runnable QUICKSTART. |

M1 and M2 each run cumulative selectors before snapshot. M1 uses the S7 read/write API and is operable independently. Product measurements include every changed file in the scratch exports. After commit the controller checks `git diff -M --numstat $B HEAD -- cmd host docs app` and its unfiltered counterpart. No later split is required by the measured S8 size.

## 3. Operator output

- The review panel shows **Review ep** with the canonical headless model, cap, repository, remote/base, workspace/branch, committer, profile, TTL, checks and grants. **Task** is a literal 1–8192 UTF-8 byte field. Enter focuses Task. **Request action** has its own click, and empty or oversized input leaves the panel and evidence unchanged. A request timeout or malformed reply leaves the task frozen and offers **Reconnect**. GET failures retain the entered command, display **Press Enter to reconnect**, and return focus to the unchanged bar. A superseded Enter or request response cannot navigate.
- The card shows the exact stored **Review** and **Task**, the state and last successful status time. Valid states include requested, confirmed, claimed, running, awaiting PR confirmation, terminal/uncertain and the launch projection launched, stopped, capped or spawn-refused. It renders the requested/expiry time and when present the runner and completion time and the server's reason. A failed poll shows **Disconnected** with the last successful timestamp. Reload and **Refresh status** re-read the same typed status address.
- On a pending card, **Confirm with passkey** is available for the enrolled current origin. The separate **Confirm with SSH signature** disclosure presents the exact S7b message hex and namespace with **Armored signature** and **Submit signature**. The UI gives the confirmation or refusal class. It focuses the card heading on load and the answer on key submission. **Action status (JSON)**, **Episode report** when available and **Back to World** are fixed-family typed paths. The validated status wire carries `v`, `id`, `episode`, `kind`, `state`, `href`, `statusHref`, optional `reportHref` and `reason`, `now`, `createdAt`, `expiresAt` and optional `terminalAt` and `runnerSeenAt`. Confirmation additionally reads `result` and `approvedBy`. `confirmHref` is reserved for the later finish flow. Finish and PR links have no producer in S8.
- QUICKSTART's two literal fixture blocks are read and entered by `TestOperatorStartWalkthroughInChrome`. It verifies an inert GET, a lost-response reconnect across reload, one pending operation, a refused `y` payload, one passkey counter advance, tab closure and reopen, one start, one registered worktree, one credential, one observed pinned confined fake launch and the back link.

## 4. Acceptance, timing and fail sets

`$S` is `$HOME/.ailang/state/world-iter447/plan-scratch`. `accept.sh` works in the supplied tree. The green sandbox marker `ACCEPT_OK <M> fast` covers named local controls and explicitly reports skipped outside legs. The controller requires `ACCEPT_OK <M> outside` for both milestones and the two Chrome count lines for M2. The 1777 leg gives its own `TMPDIR` only to the invocation under the 1777 parent. Go writes in new fixture tests have bounded contexts. Chrome steps poll explicit DOM or database conditions under deadlines; no unbounded waits are used. Worktree identity is compared via physical path. Every test-side write is under its temporary fixture or scratch.

| Package | Pristine race s | M1 race s / delta | M2 race s / delta |
|---|---:|---:|---:|
| `host/store` | 49.523 | 126.947 / +77.424 | 53.021 / +3.498 |
| `host/agentrun` | 19.790 | 29.626 / +9.836 | 20.800 / +1.010 |
| `host/surfacepage` | 29.355 | 14.192 / -15.163 | 15.866 / -13.489 |
| `host/daemon` | PENDING | PENDING | PENDING |
| `cmd/ailang-worldd` | PENDING | no cmd hunk in M1 | PENDING |

These are package samples under the rig's observed load. Compare their fail sets with the finished controller baseline. The baseline for surfacepage is the local pristine export. A sampler delta reflects the full package and can vary independently of named test time.

| New or changed named test | Planner race evidence |
|---|---|
| cmd `TestOperatorFakePresetEligibility` | 2.18 s, 3.24 s, 0.55 s, 0.50 s, 0.51 s named. Contended sweep 22.06 s. |
| cmd `TestOperatorStartWalkthroughInChrome` | CI skips because Chrome is unconfigured. Required outside plain browser run, no CI race addition. |
| daemon `TestOperatorCardAssets` | UNINFORMATIVE UNDER SANDBOX. Measure named real listener. |
| daemon `TestOperatorConfirmProjection`, `TestOperatorConfirmationAdmissionNeighbours` | Changed transport. UNINFORMATIVE UNDER SANDBOX. Include their deltas in the outside daemon budget. |

The isolated named cmd maximum is 3.24 s on the measured fixture.
Daemon added total is PENDING — measure at most 12 s outside, including the two existing tests moved to listeners. The cmd budget is 10 s. The contended sample exceeded it. Repeat the named timing with the whole-package and mutation work separated, and preserve both receipts.
The planner's explicit 1777 child directory run printed `ACCEPT_OK M2 1777` in `$S/measure-final-1777/receipt.raw`. The outside 1777 arm adds the real fixtures, including Chrome and the actual fake runner.

## 5. Branch traces, test ownership and mutation matrix

### 5.1 Bar, schema and request

| Reviewed row | Implementation and named pin | Candidate edit or existing control |
|---|---|---|
| Existing address/report/why/approve | M2 `host/surfacepage TestOperatorOtherVerbParity`, existing `TestSurfaceAppVerbs` and server `TestSurfaceVerbOutcomes` | Reclassify a legacy outcome as review. Route golden unedited. |
| Valid headless start | M2 `TestOperatorBarStartReview`, real `TestOperatorStartWalkthroughInChrome` reads zero action/start/session rows after Enter, then one requested row on click | B-review-no-writes, X-enter-post, X-submit-event. |
| Valid finish | S9 owns effect-free finish review and prep. Existing supported headless route is limited to start. | `TestOperatorGrammarOptions`, `TestOperatorConfirmHTTP/UnsupportedKinds`. |
| Case mismatch | S1 `TestOperatorGrammarCase`; M2 additionally validates canonical strings | B-review-episode, B-review-preset, B-id-case. |
| Arity/missing episode/preset/flag | S1 `TestOperatorGrammarArity`; M2 refuses mismatching wire | B-review-count and B-review-preset-flag. |
| Prefix/substring command/flag | S1 `TestOperatorGrammarExactWords`; M2 validates entire canonical argv | B-review-root, B-review-launch and B-review-headless. |
| Quoted/equals/interleaved spellings | S1 `TestOperatorGrammarCanonical`; M2 consumes the returned canonical argv | B-review-check-call; positive argv preserved in `TestOperatorLiteralTask`. |
| Duplicate flags | S1 `TestOperatorGrammarDuplicates` | Existing parser uniqueness. |
| Unsupported options, false, standalone launch | S1 `TestOperatorGrammarOptions`; M2 `TestOperatorPanelNetwork` fails sealed wrong argv | B-review-binary/kind/reserved-preset/headless. |
| Denied authority mode | S1 `TestOperatorGrammarPolicy`; M2 `TestOperatorFakePresetEligibility` plus implicit-mode browser refusal | X-clock and Browser zero-state control. |
| Forbidden verbs and shell | Existing `TestOperatorNeverDispatch` and new JS `TestOperatorOtherVerbParity` assert refusal panel before navigation | Existing per-family policy mutants. |
| ID leading/trailing anchors | M1/M2 `TestOperatorCardLinkAllowlist`, `TestOperatorLinkAllowlist`; S1 `TestOperatorIdentifierAnchors` | C-path-start/end, B-id-start/end. |
| ID multiline/case flags | Same executable tests, S1 `TestOperatorIdentifierFlags` | C-path-icase/multiline, B-id-case/lines. |
| ID bounds | Same executable tests, S1 `TestOperatorIdentifierLength` | C-id-length, C-episode-length, B-id-length. |
| Request/body/task bounds | S7 `TestOperatorAdmissionHTTP`; M2 zero, NUL and UTF-8 cap in `TestOperatorBarStartReview` | B-task-count/nul/bytes/gate. |
| Wire types | M2 `TestOperatorPanelNetwork`, `TestOperatorLinkAllowlist`; M1 `TestOperatorCardLinkAllowlist` | Individual B-review and C-status shape/type edits. |
| Strict JSON | Existing `TestOperatorAdmissionHTTP`, `TestOperatorConfirmHTTP/Proof` | Existing parser/refusal class and no writes. |
| Literal task | M2 `TestOperatorLiteralTask`/Chrome, existing `TestOperatorStartBrowser` and `TestOperatorConfirmHTTP/SSHSIG` | B-operation-freeze, X-task-markup, X-task-unescaped, X-review-text, X-status-task-leak and X-status-proof-leak; one prompt argument, escaped task and absent task/proof in reads and logs. |
| Bad Host | S7 `TestOperatorAdmissionHTTP`; M1 `TestOperatorCardAssets` also refuses foreign Host for each asset | R-host-js/css and route-specific refusal evidence. |
| Origin, fetch-site, media, body cap | S7 `TestOperatorAdmissionHTTP`, S7b `TestOperatorConfirmHTTP/Transport`; new browser uses same-origin guarded POST with credentials omitted | B-transport-credentials/type/budget; C-credentials/fetch-bound. |
| Wrong method, no CORS | S7 `TestOperatorAdmissionHTTP` and S7b Transport; M1 `TestOperatorCardAssets` refuses POST | R-method-js/css. |
| Runner absent or incompatible | S7 `TestOperatorReviewRefusesForgedCatalogue`, `TestOperatorAdmissionHTTP`, store `TestOperatorRunnerUnavailable` | Inherited readiness and guarded terminal outcome. |
| Catalogue expired, generation, identity, missing, store busy | S7 `TestOperatorReviewRefusesForgedCatalogue`/`TestOperatorReplayConflicts`; M1 real `TestOperatorConfirmationAdmissionNeighbours`; M2 `TestOperatorPanelNetwork` | J02, J10, B-review-generation-number/positive. |
| Pending cap, active kind, conflicting operation | S3 `TestOperatorAdmissionBoundaries`, S7 `TestOperatorAdmissionHTTP` | Existing count/unchanged row/request equality proofs. |
| Same operation / duplicate tap | S3 durable dedup plus M2 `TestOperatorDuplicateTap` and Chrome durable counts | B-operation-busy/mark-busy/key. |
| Lost response, refresh/restart | S3 durable tombstone and M2 `TestOperatorReconnect`; Chrome drops the real POST response before reload | B-operation-reuse/freeze/clear, X-restore-draft/freeze-ui/freeze-inflight/reconnect-label. |
| GET/POST timeout and invalid reply | M2 `TestOperatorPanelNetwork`, M1 `TestOperatorCardStatusDisconnected` | B-transport-budget, B-request-response-status, B-storage-refusal and status body guards. |
| Superseded Enter/fetch | M2 `TestOperatorPanelLatest` and existing `TestSurfaceAppVerbs` | B-bar-stale, B-operation-stale-before/after. |
| Link allowlist | M1 `TestOperatorCardLinkAllowlist` and M2 `TestOperatorLinkAllowlist`, real back link in Chrome | Each C-path and C-report edit, each B-id and B-href edit, B-status-href, X-back. Canonical relative paths encode the only acceptable scheme/origin/port policy. Flight into a V3 or remote PR link is owned by S9–S11. |

### 5.2 Card, key, status and claim

| Reviewed row | Disposition and named pin |
|---|---|
| Good passkey and SSH | M1 `TestOperatorCardStatusDisconnected` executes exact key payloads, digest, challenge, UV and RP. Its Node recording browser receives a message and assertion produced by `approval.OperatorActionMessage` and `approvaltest`. M2 Chrome drives the real virtual passkey and checks counter and challenge. S7b `TestOperatorConfirmHTTP/SSHSIG` and `TestOperatorConfirmBothKeys` cover real fixture proof formats. C-passkey-hash/uv/rp, C-assertion and C-sshsig. |
| Plain y, missing/malformed proof or bad challenge | M1 DOM class/refusal, M2 Chrome y in the SSH field shows `NotApproved` and identical fixture DB bytes/counters. Existing S7b `TestOperatorConfirmHTTP/Proof` checks raw y, missing, malformed challenge and every no-change snapshot. C-failure-class, C-digest. |
| Forged/unknown/namespace/UV/RP/origin | S7b `TestOperatorConfirmHTTP/Authenticity` plus M1 origin/RP selection tests. C-passkey-origin/rp. |
| Revocation and counter | S7b `TestOperatorConfirmHTTP/RevocationCounter`; E15 M1 `TestOperatorBackupEligibilityDrift` checks signed binding readback, rolled-back key timestamp/counter and untouched pending row. |
| Foreign action/episode/kind/store/argv/preset/task | S7b `TestOperatorConfirmHTTP/Binding` and `TestOperatorConfirmStoreBinding`. M1 rejects a ceremony POST path differing from its own card. C-ceremony-identity. |
| Expiry equality, elapsed, issue/TTL | Existing `TestOperatorConfirmHTTP/Expiry`; M1 card polls last known state, expired removes controls. C-no-unconfirmed-controls. |
| Replayed proof/confirm loss | Existing `TestOperatorConfirmHTTP/ConcurrentReplay`, `TestOperatorConfirmReplay`. M1 confirmation response validates result and re-reads status. C-confirm-result/approver/shape. |
| Missing, busy, CAS/writer loss/restart | S7b `TestOperatorConfirmHTTP/MissingAndWriter`, `/Rollback`, `TestOperatorConfirmCrash`; new card leaves last successful state at disconnect. C-disconnect-value. |
| Expired/cancelled/offline before claim | S3/S7b `TestOperatorClaimBeforeExpiry`, `TestOperatorConfirmHTTP/Cancelled`. Runner remains the sole claimant. |
| Two runners/stale generation | S3 `TestOperatorClaimRace`; current owner/admission unchanged. |
| Terminal exclusion | S3/S6 `TestOperatorTerminalExclusion`; runner lock unchanged. |
| Forged durable row after claim | S5 `TestOperatorRunnerRechecks` and `TestOperatorStartChanged`; existing proof rechecked before runner effects. |
| Restart before reserve | `TestOperatorStartOnce/worktree`, `TestOperatorStartRecoveryWindows`, `TestOperatorStartMintRecoveryRemedies` from S5/S6. |
| Restart after reserve | `TestOperatorLaunchCrashWindow`, `TestOperatorLaunchProcessStatus`; M1 Darwin fake KinfoProc reaches E19 without creating a zombie process. |
| Status offline, read-error, unknown | M1 `TestOperatorCardStatusDisconnected`, `TestOperatorCardLinkAllowlist`; S7b `TestOperatorStatusDisconnected`. C-status-state/time, C-poll-refresh, C-disconnect-value. |
| Pruning/cap | `TestOperatorPrune`, `TestOperatorAdmissionBoundaries` from S3. |
| Cancel after reservation | `TestOperatorCancelInflight` from S3. |

### Carried notes

- **E15** adds 75 test lines in `host/store`. It uses a migrated test credentials store, a registration from `approvaltest`, a verified proof and a backup eligibility change in that fixture. It observes the complete row, key timestamp/count and action/step counts before and after refusal, then confirms the same proof against the matching row as control.
- **E19** adds 23 changed lines across the Darwin reader and test. It uses a Darwin ABI fake `unix.KinfoProc` injected into the same production classifier. The test feeds a runnable and zombie row for the same PID/birth. Production still calls the same sysctl. The fake is deterministic and Darwin-tagged. Native Linux continues exercising its existing readback test.
- **J06** `status.State = "spawn-refused"` → `"launched"`, **J22** wraps the digest comparison in `false && …` (the valid r2 edit). Both run `TestOperatorConfirmProjection` on a real listener. The r1 J22 removal leaves a dead variable and was not the compilable judge arm.
- **J02** drops the positive POST generation check, **J10** drops the catalogue zero guard, **J18** changes `needs = -1` to `needs += 0`, **J19** removes the request's human lock and defer. All run the exact edits from judge-445 on `TestOperatorConfirmationAdmissionNeighbours` and the whole daemon package. See `whole-mutants.sh`.
- E16/E29/E10/E24/E27 keep their judged near-equivalent or shadowed classification. Their lines receive zero changes. The original source edits and per-stage line anchors are in the mutation catalog.

### Existing tests and fixtures modified

| File and symbol | Before and after | Behavioural mutant retained |
|---|---|---|
| `host/agentrun/process_darwin_test.go` | Existing running-process identity test unchanged. Add a fake-kernel row test using the same classifier with fixed birth. | E19 exact zombie predicate removal. |
| `host/daemon/surface_operator_confirm_test.go`, `TestOperatorConfirmProjection` | Same S7b digest, spawn-refused and birth projection controls over the real mux/listener. | J06 and J22. |
| Same file, `TestOperatorConfirmationAdmissionNeighbours` | Same request/refusal/summary cases over the real mux/listener, atomic lock-witness switch before probes, typed direct catalogue control retained. | J02, J10, J18 and J19. |
| `host/surfacepage/app_test.go`, `TestSurfaceAppPure` | Same executable legacy harness and exact export manifest expanded by the three new tested seams. Also reject a null operator response. | B-status-null executes the null-response branch. |

No existing fixture byte, golden, passkey driver, fixture recorder or test helper is edited. The existing real-mux fixtures and `startPasskeyBrowser`, `virtualPasskey` and CDP tab, keyboard, viewport and watch commands are reused. No Chrome fork is added.

### Stem-matched selectors

`daemon.go` selects every `daemon*_test.go` case through `DAEMON_CORE`. `surface_operator_read.go` and the changed confirm test select all `surface_operator*_test.go` names through `DAEMON_APPROVAL`. `app.js` and `app_test.go` select every `app*_test.go` case through `PAGE_STEMS|PAGE_BAR`. The new card tests use `PAGE_CARD`. `process_darwin.go` selects TestDarwinProcessIdentity and TestDarwinZombieIdentity on Darwin. The store backup test selects TestOperatorBackupEligibilityDrift with the store inventories. For the new CLI file run TestOperatorFakePresetEligibility and TestOperatorStartWalkthroughInChrome; the existing fake launch/receipt control and CLI boundaries are named separately. `DAEMON_GOLDEN` and `PAGE_GOLDEN` carry existing wire and asset/golden references. Run all these selectors on every cumulative stage that contains their tests.

## 6. Gate and mutation receipts

| Gate | M1 local | M2 local |
|---|---|---|
| `accept` | OK 35s | OK 10s |
| `accept-race` | PENDING — identical gate selector in §6 | OK 77s |
| `agentrun-plain` | recorded rc=0 24s | recorded rc=0 16s |
| `agentrun-race` | recorded rc=0 30s | recorded rc=0 21s |
| `broker-plain` | UNINFORMATIVE UNDER SANDBOX rc=1 | UNINFORMATIVE UNDER SANDBOX rc=1 |
| `build` | OK 48s | OK 27s |
| `callers` | OK 1s | OK 0s |
| `cmd-plain` | PENDING — identical gate selector in §6 | SKIPPED (outside sandbox required) |
| `cmd-race` | PENDING — identical gate selector in §6 | SKIPPED (outside sandbox required) |
| `compare-agentrun-plain` | SKIPPED (final pristine baseline unavailable) | SKIPPED (final pristine baseline unavailable) |
| `compare-agentrun-race` | SKIPPED (final pristine baseline unavailable) | SKIPPED (final pristine baseline unavailable) |
| `compare-broker-plain` | SKIPPED (head record unavailable) | SKIPPED (head record unavailable) |
| `compare-cmd-plain` | PENDING — identical gate selector in §6 | SKIPPED (head record unavailable) |
| `compare-cmd-race` | PENDING — identical gate selector in §6 | SKIPPED (head record unavailable) |
| `compare-daemon-plain` | SKIPPED (head record unavailable) | SKIPPED (head record unavailable) |
| `compare-daemon-race` | SKIPPED (head record unavailable) | SKIPPED (head record unavailable) |
| `compare-store-plain` | SKIPPED (final pristine baseline unavailable) | SKIPPED (final pristine baseline unavailable) |
| `compare-store-race` | SKIPPED (final pristine baseline unavailable) | SKIPPED (final pristine baseline unavailable) |
| `compare-surfacepage-plain` | SKIPPED (final pristine baseline unavailable) | SKIPPED (final pristine baseline unavailable) |
| `compare-surfacepage-race` | SKIPPED (final pristine baseline unavailable) | SKIPPED (final pristine baseline unavailable) |
| `compare-verifygate-plain` | SKIPPED (final pristine baseline unavailable) | SKIPPED (final pristine baseline unavailable) |
| `daemon-plain` | SKIPPED (outside sandbox required) | SKIPPED (outside sandbox required) |
| `daemon-race` | SKIPPED (outside sandbox required) | SKIPPED (outside sandbox required) |
| `gofmt` | OK 1s | OK 11s |
| `inv-boundary` | OK 3s | OK 19s |
| `inv-broker` | UNINFORMATIVE UNDER SANDBOX rc=1 | UNINFORMATIVE UNDER SANDBOX rc=1 |
| `inv-cli` | OK 7s | OK 17s |
| `inv-daemon` | SKIPPED (route/header tests UNINFORMATIVE UNDER SANDBOX) | SKIPPED (route/header tests UNINFORMATIVE UNDER SANDBOX) |
| `inv-daemon-static` | OK 40s | OK 6s |
| `inv-exec-kernel` | SKIPPED (Linux runtime control) | SKIPPED (Linux runtime control) |
| `inv-page` | OK 3s | OK 8s |
| `inv-store` | OK 2s | OK 12s |
| `inv-verbs` | OK 1s | OK 3s |
| `inv-verifygate` | OK 73s | OK 136s |
| `linux-vet` | OK 7s | OK 41s |
| `store-plain` | recorded rc=0 19s | recorded rc=0 133s |
| `store-race` | recorded rc=0 135s | recorded rc=0 55s |
| `surfacepage-plain` | recorded rc=0 20s | recorded rc=0 4s |
| `surfacepage-race` | recorded rc=0 15s | recorded rc=0 7s |
| `t1777` | SKIPPED (standalone M1 adds no runner-path test) | OK 37s |
| `verify_ail` | OK 231s | OK 219s |
| `verifygate-plain` | recorded rc=1 586s | recorded rc=1 477s |
| `vet` | OK 3s | OK 61s |

Raw logs are `$S/gates-M1.raw`, `$S/gates-M2.raw` and the per-leg files in `$S/gates-local-M1` and `$S/gates-local-M2`. The real asset/header/listener selectors are named in `outside.sh`. `inv-daemon-static` locally executes the route read-only AST and production caller inventories. Guard runs are instrumented by the fixture's state/cache roots. Use the exact `outside` gate command for socket-dependent receipts and final fail-set comparison.

Controller baseline status at receipt update is PENDING — require `$STATE/baseline.done` ending in `ALL_DONE` before fail-set attribution.

Whole verifygate on both M1 and M2 captured only `TestMissionConfigDefaultFlowIntegration/valid_ledger` with its bounded fixture timing out. The same named selector reproduced the failure at pristine and M1, and passed at M2. These receipts are in `$S/verifygate-known-control`. The controller's original pristine verifygate leg passed. The final outside fail-set check remains required.

The frozen stages repeated touched-file gofmt, build, native and Linux vet, and both JS syntax checks. Every seal marker passed in `$S/seal-M1/raw` and `$S/seal-M2/raw`. The published patches were regenerated from these trees and byte-compared to the rehearsed patches.

The outside mutator tests the repo's served JS by executing it in Node or in Chrome. A killed label requires the named behavioural test's `--- FAIL` and a usable completed log. There is no substitute source-grep kill. Record named controls and the prioritized whole-package list separately. A survivor is reported by id. The carry ledger is separate from the diff ledger.

| Ledger | KILLED | SURVIVED | Real listener/Chrome pending | Other pending |
|---|---:|---:|---:|---:|
| diff (171) | 145 | 0 | 26 | 0 |
| carried (8) | 2 | 0 | 6 | 0 |

Named receipts are `$S/mutations-named-r1/ledger.jsonl` and `$S/mutations-named-r2/ledger.jsonl`, with follow-ups in `$S/mutations-named-r3/ledger.jsonl`, `$S/mutations-card-final/ledger.jsonl` and `$S/mutations-card-whole-final/ledger.jsonl`. The final classification uses the latest named result. Both original id lists and the exact edits with production line anchors are in `$S/mutant-{diff,carried}-ids.txt` and `$S/mutants.json`.

Final local survivors: none. Initial `B-draft-id-type` survived the numeric-only control. The final test also supplies an array that coerces to valid hex, and its exact artifact is in the rerun. Initial `B-operation-busy` and `B-operation-mark-busy` exposed a Node microtask control that could exit with an unresolved second promise. The test now observes the send count before awaiting it, and the harness fails on any unsettled assertion.

Prioritized whole-package receipt `$S/mutations-whole-local/ledger.jsonl`: `E15` KILLED, `E19` KILLED, `C-confirm-shape` KILLED, `C-report-object` KILLED, `C-passkey-hash` KILLED, `B-operation-reuse` KILLED, `B-operation-stale-after` KILLED, `B-href-exact` KILLED, `B-review-headless` KILLED, `B-operation-busy` KILLED, `B-operation-mark-busy` KILLED. The limited outside list is `$S/whole-mutant-ids.txt`, including the six exact real-listener carried ids. Each receives a named and whole-package run. The outside script also runs every remaining browser/route mutant by its named selector.

Outside carried disposition: J06, J22, J02, J10, J18, J19 are UNINFORMATIVE UNDER SANDBOX and await the recorded real-listener run. E15 and E19 have deterministic fixture tests.

### Reproduction scripts

````bash gates.sh
#!/bin/bash
# gates.sh TREE LOGDIR M1|M2 [local|outside|outside-base]
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P) || exit 2
T=$(cd "$1" && pwd -P) || exit 2
mkdir "$2" || exit 2
L=$(cd "$2" && pwd -P) || exit 2
M=$3
mode=${4:-local}
case "$M/$mode" in M1/local|M2/local|M1/outside|M2/outside|M2/outside-base) ;; *) exit 2;; esac
B=0e32a8d7de0bb524678a3c443400cf05d6ec47eb
BASE="$S/../base"
. "$S/selectors.sh"
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
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
  local tag=$1 pkg=$2 start rc=0 parser=0 cap=1300 limit=1200s
  shift 2
  case "$tag" in
    cmd-race) cap=945; limit=900s; set -- "$@" -skip "$CMD_SKIP" ;;
    *-race) cap=945; limit=900s ;;
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
  local p bad list="$S/changed-$M"
  [ "$M" = M1 ] || list="$S/changed-M2-cumulative"
  while IFS= read -r p; do
    case "$p" in *.go) bad=$(gofmt -l "$p") || return 2; [ -z "$bad" ] || { printf '%s\n' "$bad"; return 1; };; esac
  done < "$list"
}
inventory() {
  local pkg=$1 names=$2 output rc=0
  output=$(mktemp "$L/inventory.XXXXXX") || return 2
  bound 390 go test -v -count=1 -timeout 370s "./$pkg/" -run "^($names)\$" > "$output" 2>&1 || rc=$?
  cat "$output"
  [ "$rc" -eq 0 ] || return "$rc"
  while IFS= read -r name; do grep -Fq -- "--- PASS: $name (" "$output" || return 2; done < <(printf '%s\n' "$names" | tr '|' '\n')
}
if [ "$mode" = outside-base ]; then
  record surfacepage-plain host/surfacepage
  record surfacepage-race host/surfacepage -race
  record cmd-race cmd/ailang-worldd -race
  [ "$issue$unknown" = 00 ] || { note ALL_DONE; exit 2; }
  note ALL_DONE
  exit 0
fi
pkgs=(./host/agentrun ./host/daemon ./host/store ./host/surfacepage)
tags=(store-plain store-race agentrun-plain agentrun-race surfacepage-plain surfacepage-race verifygate-plain broker-plain daemon-plain daemon-race)
owners=(host/store host/store host/agentrun host/agentrun host/surfacepage host/surfacepage host/verifygate host/broker host/daemon host/daemon)
if [ "$M" = M2 ]; then
  pkgs+=(./cmd/ailang-worldd)
  tags+=(cmd-plain cmd-race)
  owners+=(cmd/ailang-worldd cmd/ailang-worldd)
fi
leg gofmt fmt_touched
leg build bound 440 go build ./...
leg vet bound 390 go vet "${pkgs[@]}"
leg linux-vet env GOOS=linux python3 "$S/bounded.py" 450 go vet "${pkgs[@]}"
leg inv-store inventory host/store "$STORE_INVENTORY"
leg inv-boundary inventory host/boundary 'TestWorldBoundaryDependencyAllowlist|TestBoundaryASTWriteGuard'
leg inv-cli inventory cmd/ailang-worldd 'TestCLICredentialWriterSites|TestSessionRootBudgets'
leg inv-verbs inventory host/verbs "$VERBS_INVENTORY"
leg inv-page inventory host/surfacepage "$PAGE_STEMS|$PAGE_CARD|$PAGE_INVENTORY"
leg inv-verifygate inventory host/verifygate "$VERIFY_INVENTORY"
leg inv-broker inventory host/broker "$BROKER_INVENTORY"
if [ "$(uname -s)" = Linux ]; then
  leg inv-exec-kernel inventory host/verifygate TestKernelRefusesExecOfWriterOpenFile
else
  skip inv-exec-kernel 'Linux runtime control'
fi
leg verify_ail bound 650 python3 "$S/verify-scratch.py" "$T" "$L" verify_ail ""
leg callers bash "$S/callers.sh" "$S/base-tree" "$T" "$M"
mkdir -m 1777 "$L/t1777"
mkdir -m 0700 "$L/t1777/x"
if [ "$mode" = outside ]; then
  leg inv-daemon inventory host/daemon "$DAEMON_INVENTORY|$DAEMON_HEADER_INVENTORY|$DAEMON_READ_INVENTORY"
  leg accept bash "$S/accept.sh" "$T" "$M" outside
  leg accept-race bash "$S/accept.sh" "$T" "$M" outside-race
  leg t1777 env TMPDIR="$L/t1777/x" bash "$S/accept.sh" "$T" "$M" outside-1777
else
  skip inv-daemon 'route/header tests UNINFORMATIVE UNDER SANDBOX'
  leg inv-daemon-static inventory host/daemon 'TestDaemonDependencyAllowlist|TestDaemonCredentialWriterSites|TestDaemonCallerInventoryAliasControl|TestDaemonConstructorCallerInventory|TestSurfaceVerbIsReadOnly|TestSurfaceVerbWhyIsReadOnly'
  leg accept bash "$S/accept.sh" "$T" "$M" fast
  if [ "$M" = M2 ]; then
    leg accept-race bash "$S/accept.sh" "$T" "$M" race
    leg t1777 env TMPDIR="$L/t1777/x" bash "$S/accept.sh" "$T" "$M" 1777
  else skip t1777 'standalone M1 adds no runner-path test'; fi
fi
for i in "${!tags[@]}"; do
  tag=${tags[$i]}
  if [ "$mode" = local ]; then
    case "$tag" in cmd-*|daemon-*) skip "$tag" 'outside sandbox required'; continue;; esac
  fi
  case "$tag" in *-race) record "$tag" "${owners[$i]}" -race;; *) record "$tag" "${owners[$i]}";; esac
done
for tag in "${tags[@]}"; do
  if [ ! -f "$L/$tag.attributable" ]; then skip "compare-$tag" 'head record unavailable'; continue; fi
  source=
  case "$tag" in
    cmd-race|surfacepage-*) source="$S/outside-pristine/$tag.attributable" ;;
    *)
      if [ "$(tail -n 1 "$S/../baseline.done" 2>/dev/null)" = "ALL_DONE for $B" ] &&
         [ -s "$BASE/$tag.raw" ] && ! denied "$BASE/$tag.raw" && ! truncated "$BASE/$tag.raw" &&
         python3 "$S/failset.py" "$BASE/$tag.raw" "$(cat "$BASE/$tag.rc")" > "$L/$tag.controller"; then
        source="$L/$tag.controller"
      fi ;;
  esac
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
# accept.sh TREE M1|M2 [fast|race|outside|outside-race|1777|outside-1777]
set -euo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
T=$(cd "$1" && pwd -P)
M=$2
mode=${3:-fast}
case "$M" in M1|M2) ;; *) exit 2;; esac
case "$mode" in fast|race|outside|outside-race|1777|outside-1777) ;; *) exit 2;; esac
. "$S/selectors.sh"
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
LOG=$(mktemp -d "$S/accept.$M.$mode.XXXXXX")
export AILANG_STATE_DIR="$LOG/state" AILANG_CACHE_DIR="$LOG/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
cd "$T"
race=(-p 1)
case "$mode" in race|outside-race) race=(-race -p 1);; esac
go_named() {
  local package=$1 names=$2 tag=$3 rc=0
  local f="$LOG/$tag.raw"
  local options=("${race[@]}")
  if [ "$package" = cmd/ailang-worldd ] && [ "${race[0]}" = -race ]; then options+=(-skip "$CMD_SKIP"); fi
  python3 "$S/bounded.py" 185 go test -v -count=1 -timeout 170s "${options[@]}" "./$package/" -run "^($names)\$" > "$f" 2>&1 || rc=$?
  cat "$f"
  if grep -Eq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]' "$f"; then
    echo "INSTRUMENT FAILURE $tag truncated"; return 2
  fi
  if grep -Eiq '(bind|listen tcp).*(: operation not permitted|: permission denied)' "$f"; then
    echo "UNINFORMATIVE UNDER SANDBOX $tag"; return 77
  fi
  python3 "$S/failset.py" "$f" "$rc" > "$LOG/$tag.fails"
  while IFS= read -r n; do grep -Fq -- "--- PASS: $n (" "$f" || { echo "MISSING PASS $tag $n"; return 1; }; done < <(printf '%s\n' "$names" | tr '|' '\n')
  [ "$rc" -eq 0 ]
}
chrome() {
  local names=$1 tag=$2 f="$LOG/$2.receipt" rc=0
  export WORLD_CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' WORLD_REQUIRE_PASSKEY_BROWSER=1
  go_named cmd/ailang-worldd "$names" "$tag" > "$f" 2>&1 || rc=$?
  cat "$f"
  local ran skipped
  ran=$(grep -Ec "^--- (PASS|FAIL): ($names) " "$f" || true)
  skipped=$(grep -Ec "^--- SKIP: ($names) " "$f" || true)
  echo "CHROME $tag RAN=$ran SKIP=$skipped"
  unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
  [ "$ran" -gt 0 ] && [ "$skipped" -eq 0 ] && [ "$rc" -eq 0 ]
}
if [ "$mode" != 1777 ] && [ "$mode" != outside-1777 ]; then
  "$WORLD_EXEC_NODE" --check host/surfacepage/operator/card.js
  "$WORLD_EXEC_NODE" --check host/surfacepage/app/app.js
  names="$PAGE_CARD|$PAGE_STEMS"
  [ "$M" = M1 ] || names="$names|$PAGE_BAR"
  go_named host/surfacepage "$names" page
  go_named host/store "TestOperatorBackupEligibilityDrift|$STORE_INVENTORY" store
  process="$PROCESS_TEST"
  [ "$(uname -s)" != Darwin ] || process="$process|TestDarwinZombieIdentity"
  go_named host/agentrun "$process" agentrun
  go_named host/verbs "$VERBS_INVENTORY" verbs
fi
if [ "$M" = M2 ]; then
  go_named cmd/ailang-worldd TestOperatorFakePresetEligibility preset
fi
case "$mode" in
 outside|outside-race|outside-1777)
   go_named host/daemon "$DAEMON_APPROVAL|$DAEMON_CORE|$DAEMON_INVENTORY|$DAEMON_HEADER_INVENTORY|$DAEMON_READ_INVENTORY" daemon
   if [ "$M" = M2 ]; then
     go_named cmd/ailang-worldd 'TestOperatorStartBrowser|TestCLICredentialWriterSites|TestSessionRootBudgets' cmd
     if [ "$mode" != outside-race ]; then
       chrome "$CHROME_NEW" operator-chrome
       if [ "$mode" = outside ]; then chrome "$CHROME_CONTROLS" chrome-controls; fi
     fi
   fi;;
 *)
   echo "LEG daemon SKIPPED (UNINFORMATIVE UNDER SANDBOX; outside.sh $M)"
   echo "LEG chrome SKIPPED (outside.sh $M; skip never counts as Chrome evidence)";;
esac
printf 'ACCEPT_OK %s %s\n' "$M" "$mode"
````

````bash snapshot.sh
#!/bin/bash
# snapshot.sh TREE M1|M2; SNAP_STATUS_LIST/ SNAP_PARENT are rehearsal-only.
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
W=$1
M=$2
case "$M" in M1|M2) ;; *) exit 2;; esac
B=0e32a8d7de0bb524678a3c443400cf05d6ec47eb
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter447/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s8-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s8.json"
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
# Test-family and boundary selectors. The stem selectors include overlapping siblings once.
CLI_STEMS='TestOperatorSpawnRefused|TestOperatorLaunchDispatchRechecks|TestOperatorStartChanged|TestOperatorLaunchCrashWindow|TestOperatorLaunchExit|TestOperatorStartBrowser|TestOperatorFirstManagedEpisodeProcess|TestOperatorPreparationRechecksHeartbeatAfterWorktree|TestOperatorStartMintRecoveryRemedies|TestOperatorCurrentOwnerProcessRecovery|TestOperatorAttendedStartDisposition|TestOperatorActivationDrainRefusals|TestOperatorRevokeStatRefusal|TestOperatorStartOnce|TestOperatorStartCollisions|TestOperatorWorktreeWindow|TestOperatorWorktreeDispositionCAS|TestOperatorWorktreeSurvivors|TestOperatorWorktreeRecoveryGuards|TestOperatorWorktreeReadbackGuards|TestOperatorReadbackErrorIsNeverCreated|TestOperatorWorktreeLockedReread|TestOperatorBorrowedWorktreeGitProcess'
CLI_S6='TestOperatorLaunchWorker|TestOperatorSupervisorCap|TestOperatorStartedRequiresBirth|TestOperatorTerminalExclusion|TestOperatorRunnerKernelLock|TestOperatorStatusDisconnected|TestOperatorProcessWitness'
DAEMON_STEMS='TestOperatorReviewRefusesForgedCatalogue|TestOperatorReplayConflicts|TestOperatorLinksProjection|TestOperatorConfirmProjection|TestOperatorConfirmationAdmissionNeighbours'
DAEMON_LISTENER='TestOperatorAdmissionHTTP|TestOperatorConfirmHTTP'
DAEMON_GOLDEN='TestSurfaceVerbOutcomes|TestRequestsJSONConfirmationsGolden'
PAGE_GOLDEN='TestSurfaceAppVerbs|TestNeedsYouConfirmationEmbeddedNewline|TestNeedsYouConfirmationsNonArray|TestVerbDecidePrefixBoundary'
DAEMON_INVENTORY='TestDaemonDependencyAllowlist|TestDaemonCredentialWriterSites|TestFrozenV1RouteTableMatchesMux|TestSurfaceSecurityHeaders|TestApprovalRoutesRegistered|TestApprovalRoutesRefuseForeignHost|TestConfirmRoutesRegistered|TestDecisionRoutesRegistered|TestSurfaceVerbRouteRegistered|TestSurfaceVerbIsReadOnly|TestSurfaceVerbBoundsAndErrors|TestSurfaceVerbWhyIsReadOnly|TestSurfaceVerbWhyOutcomes|TestWhyResult|TestSurfaceHomeSummary|TestReadRoutesAndPayloadGate|TestGETRoutesRejectOtherMethods|TestProjectionRoutes_MountedAndNotProtected|TestAgentCardRoute_GETOnly|TestSurfaceAppRoutes|TestSurfaceScriptRoute|TestSurfaceModelRoutes|TestApprovePageRoutes|TestPasskeyPageRoutes|TestSurfaceFontRoute|TestWorkbenchRouteIsReadOnly|TestWorkbenchSecurityHeaders'
STORE_STEMS='TestOperatorLaunchReceiptLaw|TestOperatorLaunchSession|TestOperatorLaunchProcessStatus'
STORE_LAWS='TestOperatorAdmissionBoundaries|TestOperatorClaimBeforeExpiry|TestOperatorClaimRace|TestOperatorTerminalExclusion|TestOperatorStartRecoveryWindows|TestOperatorPrune|TestOperatorCancelInflight|TestOperatorCancelConfirmed|TestOperatorRunnerUnavailable'
STORE_INVENTORY='TestProductionContextRoots|TestStoreSurfaceTakesContext|TestCredentialsMethodsRequireADeadline|TestStoreDeadlineGuardDynamic'
AUTHORITY_LAWS='TestOperatorJudgeBothKeys|TestOperatorConfirmReplay|TestOperatorConfirmCrash|TestOperatorConfirmReadOnly|TestOperatorConfirmStoreBinding|TestOperatorRunnerRechecks'
APPROVAL_LAWS='TestOperatorConfirmBothKeys|TestOperatorConfirmProof|TestOperatorConfirmAuthenticity|TestOperatorConfirmRevocationCounter|TestOperatorConfirmBinding|TestOperatorConfirmExpiry'
VERIFY_INVENTORY='TestSubprocessSinkGate|TestExecGateArmsAreLoadBearing|TestNoParallelWriteForkPackagesOutsideVerifygate|TestVerifygateTestWritesAreForkLocked'
BROKER_INVENTORY='TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential'
case "$(uname -s)" in
  Darwin) PROCESS_TEST=TestDarwinProcessIdentity ;;
  Linux) PROCESS_TEST=TestLinuxProcessIdentity ;;
  *) PROCESS_TEST=TestUnsupportedProcessIdentity ;;
esac

DAEMON_HEADER_INVENTORY='TestSurfacePage|TestSurfaceIndex|TestSurfaceFollowRoutes|TestSurfaceFollowRoutesLeakNothing|TestSurfaceFileReadsOnlyPublishedLandedOutputs|TestDecidePageAssets|TestConfirmPageAssetsServed|TestConfirmationCardServed|TestApprovePOSTCrossOrigin|TestApprovePOSTNoSessionNoCookie|TestConfirmPOSTGuards|TestWorkbenchLiveScriptRoute|TestWorkbenchErrorPageInert|TestWorkbenchStillAllRegions|TestWorkbenchTimelineBound|TestAGUIRejectsGET|TestAGUIInputErrors'
DAEMON_READ_INVENTORY='TestDaemonReadDeadline|TestWorkbenchReadDeadline|TestReadCtxCancelledAfterHandler|TestInternalErrorsAreSanitized'
VERBS_INVENTORY='TestOperatorNeverDispatch|TestOperatorOtherVerbParity'
DAEMON_CORE='TestIsLoopbackHostMirrorsSketchPredicate|TestNewRefusesNonLoopbackBind|TestBoundedWaitsAndBodyLimit|TestBoundedGracefulShutdownDrainsInFlightRequest|TestDaemonShutdownIsBoundedByTheD7Constant|TestNewReleasesWriterAuthorityOnLateFailure|TestDrainHardClosesWhenDeadlineExpires|TestHealthAndHeadRoundTrip|TestRunAnnouncesResolvedListenAddress|TestNewBootstrapsEpochRegistryIdempotently|TestReleaseFromVersion|TestDaemonDependencyAllowlist|TestAilangProtocolAdmissionIsNarrow|TestProjectionRoutes_MountedAndNotProtected|TestCardDenial_ByteIdenticalToCommitMiddleware|TestCardDenial_ExactClassAndMessagePerKind|TestAgentCardRoute_GETOnly|TestAgentCard_ViaDaemon_SessionScoped|TestA2AAuthorizedSkillNoCoordinator|TestA2ACardInvocationDescription|TestDaemonInvokeConfig|TestAgentCard_ViaDaemon_AbsentHeadZeroSkills'
DAEMON_APPROVAL="$DAEMON_STEMS|$DAEMON_LISTENER|$DAEMON_GOLDEN|TestOperatorCardAssets"
PAGE_STEMS="TestSurfaceAppMatchesWorkingCopy|TestSurfaceAppShellAssembly|TestSurfaceAppPure|TestSurfaceAppTextIsSelectable|TestSurfaceAppFollow|TestAddressHintsArePinned|TestNeedsYouListsDecisions|TestNeedsYouListsConfirmations|$PAGE_GOLDEN"
PAGE_CARD='TestOperatorCardStatusDisconnected|TestOperatorCardLinkAllowlist'
PAGE_BAR='TestOperatorBarStartReview|TestOperatorOtherVerbParity|TestOperatorLiteralTask|TestOperatorDuplicateTap|TestOperatorReconnect|TestOperatorPanelNetwork|TestOperatorPanelLatest|TestOperatorLinkAllowlist'
CHROME_NEW=TestOperatorStartWalkthroughInChrome
CHROME_CONTROLS='TestConfirmPageConfirmsWithAPasskeyInChrome|TestVerbBarEpisodeReportInChrome'
CMD_SKIP='^(TestCLIAgainstSeToolsDaemon|TestCLINoGenesisIsNoWorldHead|TestDemoSeedApprovalPacket|TestDemoSeedEndToEnd|TestDemoSeedManifestReferences|TestDemoSeedProofStage|TestDemoSeedReadTestCheck|TestDemoSeedRejectsToolReportRefusals|TestDemoSeedRosterKindMismatchRefuses|TestDemoSeedWhyNamedLinks|TestMarketingStillRoutes)$'
DAEMON_INVENTORY="$DAEMON_INVENTORY|TestDaemonCallerInventoryAliasControl|TestDaemonConstructorCallerInventory|TestEpisodeReportRESTRoute|TestEpisodeStartReadProjectionsAndSessionRoutesCannotWrite|TestLiveScriptContract|TestMarketingCaptureInventoryMutation|TestWorkbenchLiveScriptTag"
PAGE_INVENTORY='TestConfirmPageAssets|TestApprovePageRendersEveryCanonicalField|TestSurfaceScriptAssetSeams|TestPreviewScriptMatchesEmbeddedBundle|TestSurfaceFontEmbed|TestConfirmCardGolden|TestSurfaceModelGolden|TestSurfaceStackModelGolden|TestSurfaceAppMatchesWorkingCopy|TestSurfaceAppShellAssembly'
````

````bash callers.sh
#!/bin/bash
# callers.sh BASE_TREE STAGE_TREE
set -euo pipefail
census() (
  cd "$1"
  local label=$2 req=1 claim=1 conf=$3
  src() { find cmd host -name '*.go' ! -name '*_test.go' -print0 | xargs -0 grep -nHE "$1"; }
  grep -nF 'approval.ConfirmOperatorAction(' host/authority/judge_operator.go
  grep -nF 'agentrun.StartWorker(' cmd/ailang-worldd/operator_launch.go
  grep -nF 'return e.runOperatorPreparationLoop(' cmd/ailang-worldd/operator_runner.go
  for pair in "RequestOperatorAction:$req" "ConfirmOperatorAction:$conf" "JudgeOperatorAction:$conf" "ClaimOperatorAction:$claim"; do
    f=${pair%:*}; n=${pair#*:}
    calls=$(src "\.$f\\(" | grep -vF 'approval.ConfirmOperatorAction(' || true)
    printf '%s\n' "$calls"
    count=$(printf '%s\n' "$calls" | grep -c . || true)
    test "$count" -eq "$n"
    printf 'CALLERS %s %s %s\n' "$label" "$f" "$count"
  done
)
census "$1" base 1
census "$2" "${3:-stage}" 1
echo CALLER_CENSUS_OK
````

````bash outside.sh
#!/bin/bash
# outside.sh M1|M2 [TREE] -- controller only, sequential and isolated.
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
M=$1
T=${2:-$S/st-$M}
L=$(mktemp -d "$S/outside.$M.XXXXXX")
issue=0
for mode in outside outside-race; do
  bash "$S/accept.sh" "$T" "$M" "$mode" > "$L/$mode.raw" 2>&1 || issue=1
  cat "$L/$mode.raw"
done
mkdir -m 1777 "$L/t1777"
mkdir -m 0700 "$L/t1777/x"
TMPDIR="$L/t1777/x" bash "$S/accept.sh" "$T" "$M" outside-1777 > "$L/t1777.raw" 2>&1 || issue=1
cat "$L/t1777.raw"
echo "OUTSIDE $M $L issue=$issue"
echo ALL_DONE
exit "$issue"
````

````bash commit-guards.sh
#!/bin/bash
# commit-guards.sh TREE LOGDIR, outside. Attribute only comparable clean base fail sets.
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
B=0e32a8d7de0bb524678a3c443400cf05d6ec47eb
mkdir "$2" || exit 2
L=$(cd "$2" && pwd -P)
cd "$1" || exit 2
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
rc=0
python3 "$S/bounded.py" 1280 go test -count=1 ./host/verifygate/ ./host/broker/ -p 1 -timeout 20m > "$L/raw" 2>&1 || rc=$?
cat "$L/raw"
if grep -Eiq 'panic: test timed out|BOUND EXCEEDED|\[(build|setup) failed\]|(bind|listen tcp).*(: operation not permitted|: permission denied)' "$L/raw"; then
  echo "INSTRUMENT FAILURE guard result truncated or sandboxed"
  exit 2
fi
python3 "$S/failset.py" "$L/raw" "$rc" > "$L/all.fails" || exit 2
[ "$(tail -n 1 "$S/../baseline.done")" = "ALL_DONE for $B" ] || exit 2
for pkg in verifygate broker; do
  awk -v p="github.com/sunholo-data/ailang-world/host/$pkg" '
    /^(ok|FAIL)[[:space:]]/ {if ($2==p) {printf "%s%s\n", lines,$0; exit} lines="";next}
    {lines=lines $0 "\n"}
  ' "$L/raw" > "$L/$pkg.raw"
  test -s "$L/$pkg.raw" || exit 2
  grc=0; grep -q "^FAIL[[:space:]]" "$L/$pkg.raw" && grc=1
  base="$S/../base/$pkg-plain"
  grep -Eq 'panic: test timed out|BOUND EXCEEDED' "$base.raw" && exit 2
  python3 "$S/failset.py" "$base.raw" "$(cat "$base.rc")" > "$L/$pkg.base" || exit 2
  python3 "$S/failset.py" "$L/$pkg.raw" "$grc" > "$L/$pkg.head" || exit 2
  LC_ALL=C comm -13 "$L/$pkg.base" "$L/$pkg.head" > "$L/$pkg.new"
  [ ! -s "$L/$pkg.new" ] || { cat "$L/$pkg.new"; exit 1; }
done
echo COMMIT_GUARDS_OK
````

````bash whole-mutants.sh
#!/bin/bash
# Controller outside runner. First run the prioritized whole-package ids, then the remaining named outside ids.
set -uo pipefail
S=$(cd "$(dirname "$0")" && pwd -P)
L=$1
mkdir "$L" || exit 2
issue=0
python3 "$S/mutate.py" "$L/whole" "$S/whole-mutant-ids.txt" --outside --whole || issue=1
python3 - "$S" <<'PY' > "$L/remaining-ids.txt"
from pathlib import Path
import json,sys
r=Path(sys.argv[1]); already=set((r/"whole-mutant-ids.txt").read_text().splitlines())
for key,item in json.loads((r/"mutants.json").read_text()).items():
    if item["outside"] and key not in already:
        print(key)
PY
if [ -s "$L/remaining-ids.txt" ]; then
  python3 "$S/mutate.py" "$L/named" "$L/remaining-ids.txt" --outside || issue=1
fi
echo "WHOLE_MUTANTS issue=$issue"
echo ALL_DONE
exit "$issue"
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
#!/usr/bin/env python3
"""Apply one anchored edit in an isolated archive copy and test real behaviour."""
import argparse
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import time

p = argparse.ArgumentParser()
p.add_argument("ledger_dir")
p.add_argument("ids")
p.add_argument("--outside", action="store_true")
p.add_argument("--whole", action="store_true")
a = p.parse_args()
root = Path(__file__).resolve().parent
out = Path(a.ledger_dir).resolve()
out.mkdir()
mutants = json.loads((root / "mutants.json").read_text())
ids = [x.strip() for x in Path(a.ids).read_text().splitlines() if x.strip()]
if len(set(ids)) != len(ids) or (a.whole and len(ids) > 25):
    sys.exit("invalid mutant list")
env = dict(os.environ, AILANG_BIN=str(Path.home()/".pinned-ailang/ailang"),
           WORLD_REQUIRE_SSH_KEYGEN="1", WORLD_EXEC_NODE="/opt/homebrew/bin/node",
           WORLD_TOOL_AILANG_BIN=str(Path.home()/".pinned-ailang-tools/v0.54.15/ailang"))
for var in ("GOFLAGS","WORLD_CHROME","WORLD_REQUIRE_PASSKEY_BROWSER","WORLD_PASSKEY_RECORD_DIR","WORLD_TEST_OPERATOR_WITNESS"):
    env.pop(var,None)

def bound(tree, exp, cap, args, raw):
    start = time.monotonic()
    with raw.open("w") as f:
        rc = subprocess.call([sys.executable,str(root/"bounded.py"),str(cap),*args],
                             cwd=tree,env=exp,stdout=f,stderr=subprocess.STDOUT)
    return rc,round(time.monotonic()-start,3),raw.read_text(errors="replace")

def classify(m, rc, text, ran):
    test = re.escape(m["test"])
    failure = re.findall(r"(?m)^--- FAIL: "+test+r" \(",text)
    success = re.findall(r"(?m)^--- PASS: "+test+r" \(",text)
    skipped = re.findall(r"(?m)^--- SKIP: "+test+r" \(",text)
    if ran or "Chrome" in m["test"]:
        print(f'CHROME {m["id"]} RAN={len(failure)+len(success)} SKIP={len(skipped)}',flush=True)
    if re.search(r"(bind|listen tcp).*(: operation not permitted|: permission denied)",text,re.I):
        return "UNINFORMATIVE UNDER SANDBOX"
    if re.search(r"panic: test timed out|BOUND EXCEEDED|\[(?:build|setup) failed\]",text):
        return "INSTRUMENT FAILURE"
    if failure and rc:
        return "KILLED"
    if skipped or not success:
        return "INSTRUMENT FAILURE"
    if success:
        return "SURVIVED"
    return "INSTRUMENT FAILURE"

def execute(m, tree, leg):
    exp = dict(env)
    for name, suffix in (("AILANG_STATE_DIR","state"),("AILANG_CACHE_DIR","cache")):
        folder = out / (m["id"]+"-"+leg+"-"+suffix)
        folder.mkdir()
        exp[name]=str(folder)
    pattern = ["-run", "^"+m["test"]+"$"] if leg == "named" else []
    if "Chrome" in m["test"] and leg == "named":
        exp["WORLD_CHROME"]="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
        exp["WORLD_REQUIRE_PASSKEY_BROWSER"]="1"
    raw = out/(m["id"]+"-"+leg+".raw")
    cap, limit = (185,"170s") if leg == "named" else (1320,"20m")
    rc, seconds, text = bound(tree,exp,cap,["go","test","-v","-p","1","-count=1","-timeout",limit,*pattern,"./"+m["pkg"]+"/"],raw)
    label=classify(m,rc,text,False)
    with (out/(m["id"]+"-"+leg+".fails")).open("w") as fail:
        prc=subprocess.call([sys.executable,str(root/"failset.py"),str(raw),str(rc)], stdout=fail)
    if prc and label!="UNINFORMATIVE UNDER SANDBOX":
        label="INSTRUMENT FAILURE"
    return dict(label=label,rc=rc,seconds=seconds,raw=str(raw))

seen=set()
bad=False
with (out/"ledger.jsonl").open("w") as ledger:
    for id in ids:
        m = mutants[id]
        inp = dict(id=id,stage=m["stage"],file=m["file"],line=m["line"],test=m["test"],kind=m["kind"])
        if m["outside"] and not a.outside:
            inp["named"] = {"label": "UNINFORMATIVE UNDER SANDBOX"}
        elif id=="E19" and sys.platform!="darwin":
            inp["named"]={"label":"SKIPPED (Darwin fake KinfoProc ABI)"}
        else:
            tree = out/("tree-"+m["stage"])
            if m["stage"] not in seen:
                shutil.copytree(root/("st-"+m["stage"]),tree,symlinks=True)
                seen.add(m["stage"])
            file = tree/m["file"]
            src = file.read_text()
            i = src.index(m["pre"]) if m["pre"] else 0
            pos = src.index(m["old"],i)
            if src[:pos].count("\n")+1 != m["line"]:
                raise ValueError("stale diff anchor "+id)
            try:
                file.write_text(src[:pos]+m["new"]+src[pos+len(m["old"]):])
                inp["named"] = execute(m,tree,"named")
                if a.whole:
                    inp["whole"] = execute(m,tree,"whole")
            finally:
                file.write_text(src)
        print(json.dumps(inp),file=ledger,flush=True)
        kind=inp["named"]["label"]
        print(id,kind,inp.get("whole",{}).get("label",""),flush=True)
        if kind in ("SURVIVED","INSTRUMENT FAILURE") or inp.get("whole",{}).get("label") in ("SURVIVED","INSTRUMENT FAILURE"):
            bad=True
print("ALL_DONE",flush=True)
sys.exit(int(bad))
````

````python mutant-spec.py
from pathlib import Path
import json

root = Path(__file__).resolve().parent
items = {}
def spec(id, stage, file, old, new, pkg, test, outside=False, pre="", kind="diff"):
    assert id not in items
    source = (root / ("st-" + stage) / file).read_text()
    start = source.index(pre) if pre else 0
    if old not in source[start:]:
        raise ValueError((id, file, old))
    line = source[:source.index(old, start)].count("\n")+1
    items[id] = dict(id=id,stage=stage,file=file,old=old,new=new,pkg=pkg,test=test,
                     outside=outside, pre=pre,kind=kind,line=line)
def card(id, old, new, test="TestOperatorCardLinkAllowlist", pre=""):
    spec("C-"+id,"M1","host/surfacepage/operator/card.js",old,new,"host/surfacepage",test,pre=pre)
def dom(id, old, new, pre=""):
    card(id,old,new,"TestOperatorCardStatusDisconnected",pre)
def bar(id, old,new,test="TestOperatorPanelNetwork",pre=""):
    spec("B-"+id,"M2","host/surfacepage/app/app.js",old,new,"host/surfacepage",test,pre=pre or "const OPERATOR_DRAFT")
def href(id, old, new):
    bar(id,old,new,"TestOperatorLinkAllowlist","const OPERATOR_DRAFT")
def chrome(id, file, old,new,pre=""):
    spec("X-"+id,"M2",file,old,new,"cmd/ailang-worldd","TestOperatorStartWalkthroughInChrome",True,pre)
def carry(id, file, old,new,pkg,test,outside=False):
    spec(id,"M1",file,old,new,pkg,test,outside,kind="carried")

D="host/daemon/"
A="host/surfacepage/app/app.js"
for key, old, new, test in [
 ("assets-js",'"GET /surface/operator-card.js"','"GET /surface/operator-card-disabled.js"',"TestOperatorCardAssets"),
 ("assets-css",'"GET /surface/operator-card.css"','"GET /surface/operator-card-disabled.css"',"TestOperatorCardAssets"),
 ("host-js","d.approvalRoute(operatorScript)","func(w http.ResponseWriter, r *http.Request) { operatorScript(w, r, nil) }","TestOperatorCardAssets"),
 ("host-css","d.approvalRoute(operatorStyle)","func(w http.ResponseWriter, r *http.Request) { operatorStyle(w, r, nil) }","TestOperatorCardAssets"),
 ("method-js",'"GET /surface/operator-card.js"','"/surface/operator-card.js"',"TestOperatorCardAssets"),
 ("method-css",'"GET /surface/operator-card.css"','"/surface/operator-card.css"',"TestOperatorCardAssets"),
]:
    spec("R-"+key,"M1",D+"daemon.go",old,new,"host/daemon",test,True)
for key, old, new in [
 ("media-js",'"text/javascript; charset=utf-8", surfacepage.OperatorScript','"text/plain; charset=utf-8", surfacepage.OperatorScript'),
 ("media-css",'"text/css; charset=utf-8", surfacepage.OperatorStyle','"text/plain; charset=utf-8", surfacepage.OperatorStyle'),
 ("template-js",'src="/surface/operator-card.js"','src="/surface/operator-card-other.js"'),
 ("template-css",'href="/surface/operator-card.css"','href="/surface/operator-card-other.css"'),
]:
    spec("R-"+key,"M1",D+"surface_operator_read.go",old,new,"host/daemon","TestOperatorCardAssets",True)
for key, old, new in [
 ("path-start","/^\\/surface\\/operator-actions","/\\/surface\\/operator-actions"),
 ("path-end","{32})(?![\\s\\S])/","{32})/"),
 ("path-icase","{32})(?![\\s\\S])/","{32})(?![\\s\\S])/i"),
 ("path-multiline","{32})(?![\\s\\S])/","{32})(?![\\s\\S])/m"),
 ("id-length","{32})(?![\\s\\S])/","{31,33})(?![\\s\\S])/"),
 ("path-object","typeof path === 'string'","path !== undefined"),
 ("episode-start","/^[a-z0-9][a-z0-9-]","/[a-z0-9][a-z0-9-]"),
 ("episode-end","{0,63}(?![\\s\\S])/","{0,63}/"),
 ("episode-icase","{0,63}(?![\\s\\S])/","{0,63}(?![\\s\\S])/i"),
 ("episode-multiline","{0,63}(?![\\s\\S])/","{0,63}(?![\\s\\S])/m"),
 ("episode-length","{0,63}","{0,64}"),
 ("no-status","!s || s.v","s.v"),
 ("status-version","s.v !== 1","false"),
 ("status-id","s.id !== pathID(path)","false"),
 ("status-card-link","s.href !== path","false"),
 ("status-self-link","s.statusHref !== path + '/status.json'","false"),
 ("status-episode","!episodeID.test(s.episode)","false"),
 ("status-kind","!['start', 'finish'].includes(s.kind)","false"),
 ("status-state","!states.has(s.state)","false"),
 ("status-time","![s.now, s.createdAt, s.expiresAt].every(Number.isSafeInteger)","false"),
 ("status-optional-time","![s.terminalAt, s.runnerSeenAt].every(n => n === undefined || Number.isSafeInteger(n))","false"),
 ("status-reason","(s.reason !== undefined && typeof s.reason !== 'string')","false"),
 ("ceremony-exists","return c && c.post","return c.post"),
 ("ceremony-identity","c.post === path + '/confirm'","true"),
 ("ceremony-digest","hex64.test(c.digest)","true"),
 ("ceremony-message","hexBytes.test(c.messageHex)","true"),
 ("passkey-group","if (!group ||","if ("),
 ("passkey-list-shape","!Array.isArray(group.allow)","false"),
 ("passkey-list-empty","group.allow.length === 0","false"),
]:
    card(key,old,new)
for key, old,new in [
 ("passkey-origin","k.origin === origin && k.rp === hostname","k.rp === hostname"),
 ("passkey-rp","k.origin === origin && k.rp === hostname","k.origin === origin"),
 ("passkey-hash","digest('SHA-256'","digest('SHA-384'"),
 ("passkey-uv","userVerification: 'required'","userVerification: 'preferred'"),
 ("passkey-rpid","rpId: group.rp","rpId: 'wrong'"),
 ("method","method: 'POST'","method: 'GET'"),
 ("digest","digest: ceremony.digest","digest: '0'.repeat(63)"),
 ("sshsig","sshsig: $('operator-signature').value","sshsig: ''"),
 ("assertion","key.toJSON()","{}"),
 ("redirect","redirect: 'error'","redirect: 'follow'"),
 ("credentials","credentials: 'omit'","credentials: 'same-origin'"),
 ("cache","cache: 'no-store'","cache: 'default'"),
 ("fetch-bound","signal: AbortSignal.timeout(15000)","signal: undefined"),
 ("status-url","read(path + '/status.json')","read(path)"),
 ("confirm-url","read(path + '/confirm'","read(path + '/other'"),
 ("failure-class","body.error + ' '","''"),
 ("disconnect-value","'Disconnected. ' + last","'Disconnected. '"),
 ("status-time-shown","'Status at ' + s.now","''"),
 ("status-final","$('operator-state').textContent = state;","$('operator-state').textContent = 'requested';"),
 ("confirm-shape","&& statusFor(body, path) &&","&& true &&"),
 ("confirm-result","['confirmed', 'already-confirmed'].includes(body.result)","true"),
 ("confirm-approver","typeof body.approvedBy === 'string'","true"),
 ("confirm-approver-empty","body.approvedBy !== ''","true"),
 ("no-unconfirmed-controls","!ceremony || state !== 'requested'","!ceremony"),
 ("report-hide","$('operator-report').hidden = !s.reportHref;","$('operator-report').hidden = false;"),
 ("report-clear","else $('operator-report').removeAttribute('href');",""),
 ("title-focus","$('operator-title').focus();",""),
 ("answer-focus","answer.focus();",""),
 ("init-refresh","\n  refresh();\n})();","\n})();"),
 ("poll-refresh","polling = setTimeout(refresh, 5000);","polling = setTimeout(() => {}, 5000);"),
]:
    dom(key,old,new)
card("ceremony-digest-type","typeof c.digest === 'string'","true")
card("ceremony-message-type","typeof c.messageHex === 'string'","true")
card("status-episode-type", "typeof s.episode !== 'string'", "false")
for key, extra in [
 ("report-prefix","String(s.reportHref).startsWith(report)"),("report-suffix","String(s.reportHref).endsWith(report)"),
 ("report-case","String(s.reportHref).toLowerCase() === report"),("report-object","String(s.reportHref) === report"),
 ("report-query","s.reportHref === report + '?q'"),("report-fragment","s.reportHref === report + '#f'"),
 ("report-scheme","s.reportHref === 'http://user@localhost' + report"),("report-host","s.reportHref === '//localhost' + report"),
 ("report-port","s.reportHref === 'https://localhost:2' + report")]:
    card(key,"s.reportHref === report ? report",f"(s.reportHref === report || {extra}) ? report")
for key, old, new in [
 ("review-shape","Array.isArray(argv) && argv.length === 8","argv.length === 8"),
 ("review-count","argv.length === 8","true"),
 ("review-binary","argv[0] === 'ailang-worldd'","true"),
 ("review-root","argv[1] === 'episode'","true"),
 ("review-kind","argv[2] === 'start'","true"),
 ("review-episode","EPISODE_ID.test(argv[3])","true"),
 ("review-preset-flag","argv[4] === '--preset'","true"),
 ("review-preset","PRESET_ID.test(argv[5])","true"),
 ("review-reserved-preset","argv[5] !== 'se-tools'","true"),
 ("review-launch","argv[6] === '--launch'","true"),
 ("review-headless","argv[7] === '--headless'","true"),
 ("review-generation-number","Number.isSafeInteger(catalogueGeneration)","true"),
 ("review-generation-positive","catalogueGeneration > 0","true"),
 ("review-check-call","return reviewValid(b.argv, b.catalogueGeneration)","return true"),
 ("freeze-write","io.storage.setItem(OPERATOR_DRAFT, JSON.stringify({ body: draft?.body ?? panel.body, request }));",""),
 ("storage-refusal","return { error: 'Request could not be saved for reconnect. Restore browser storage and retry.' };",""),
 ("transport-post","method: 'POST', headers: {'Content-Type': 'application/json'}","method: 'GET', headers: {'Content-Type': 'application/json'}"),
 ("transport-type","'Content-Type': 'application/json'","'Content-Type': 'text/plain'"),
 ("transport-target","io.fetch(OPERATOR_PATH,","io.fetch('/surface/verb.json',"),
 ("transport-redirect","redirect: 'error'","redirect: 'follow'"),
 ("transport-credentials","credentials: 'omit'","credentials: 'same-origin'"),
 ("transport-budget","io.timeout(VERB_ROUTE_BUDGET_MS)","io.timeout(0)"),
 ("request-response-status","res.status === 200 && operatorCardHref","operatorCardHref"),
 ("request-error-class","b.error + ': ' + b.reason","'unreadable'"),
]:
    bar(key,old,new)
for key,old,new in [
 ("draft-text-type","typeof d?.body === 'string'","true"),
 ("draft-text-empty","&& d.body && r","&& r"),
 ("draft-id-type","typeof r.operationId === 'string'","true"),
 ("draft-id-shape","OPERATION_ID.test(r.operationId)","true"),
 ("draft-task","taskValid(r.task)","true"),
]:
    bar(key,old,new)
bar("get-retry","if (act.retry && fx.retry)","if (false)")
bar("get-retry-shape","this page cannot read', retry: true","this page cannot read', retry: false")
bar("get-retry-timeout","did not answer in time', retry: true","did not answer in time', retry: false")
bar("get-retry-http","HTTP ' + status, retry: true","HTTP ' + status, retry: false")
bar("get-retry-offline","String(e)), retry: true","String(e)), retry: false")
bar("other-verbs","if (b.outcome === 'review')","if (b.outcome === 'review' || b.outcome === 'result')","TestOperatorOtherVerbParity")
bar("status-null","if (!b || b.v !== 1 || b.kind","if (b.v !== 1 || b.kind","TestSurfaceAppPure")
for key, old,new,test in [
 ("review-no-writes","if (act.review) return fx.review(act.review);","if (act.review) { fx.clear(); return fx.review(act.review); }","TestOperatorBarStartReview"),
 ("task-count","task.length > 0","true","TestOperatorBarStartReview"),
 ("task-nul","!task.includes('\\0')","true","TestOperatorBarStartReview"),
 ("task-bytes","encode(task).length <= 8192","encode(task).length <= 8193 + 1","TestOperatorBarStartReview"),
 ("task-gate","if (!taskValid(task))","if (false)","TestOperatorBarStartReview"),
 ("operation-key","io.random(new Uint8Array(16))","io.random(new Uint8Array(15))","TestOperatorLiteralTask"),
 ("operation-reuse","let request = draft?.request;","let request;","TestOperatorReconnect"),
 ("operation-freeze","body: JSON.stringify(request)","body: JSON.stringify({...request,task})","TestOperatorReconnect"),
 ("operation-clear","io.storage.removeItem(OPERATOR_DRAFT);","","TestOperatorReconnect"),
 ("operation-busy","if (busy || !io.current())","if (!io.current())","TestOperatorDuplicateTap"),
 ("operation-mark-busy","busy = true;","busy = false;","TestOperatorDuplicateTap"),
 ("operation-release","busy = false;\n      }","busy = true;\n      }","TestOperatorPanelNetwork"),
 ("operation-stale-before","if (busy || !io.current())","if (busy)","TestOperatorPanelLatest"),
 ("operation-stale-after","if (!io.current()) return {};","if (false) return {};","TestOperatorPanelLatest"),
 ("bar-stale","if (!current()) return;\n    return verbAct","if (false) return;\n    return verbAct","TestOperatorPanelLatest"),
]:
    bar(key,old,new,test)
for key, old,new in [
 ("id-type","typeof b.id !== 'string'","false"),
 ("id-start","/^oa-[0-9a-f]","/oa-[0-9a-f]"),
 ("id-end","{32}(?![\\s\\S])/","{32}/"),
 ("id-case","{32}(?![\\s\\S])/","{32}(?![\\s\\S])/i"),
 ("id-lines","{32}(?![\\s\\S])/","{32}(?![\\s\\S])/m"),
 ("id-length","{32}(?![\\s\\S])/","{31,33}(?![\\s\\S])/"),
 ("id-grammar","!OPERATOR_ID.test(b.id)","false"),
 ("status-v","b.v !== 1 || b.kind","false || b.kind"),
 ("status-kind","b.kind !== 'start'","false"),
 ("status-episode","b.episode !== episode","false"),
 ("status-state",".includes(b.state) ||",".includes(b.state) && false ||"),
 ("status-time","![b.now, b.createdAt, b.expiresAt].every(Number.isSafeInteger)","false"),
 ("href-exact","b.href === path","true"),
 ("href-substring","b.href === path","String(b.href).includes(path)"),
 ("href-case","b.href === path","String(b.href).toLowerCase() === path"),
 ("href-object","b.href === path","String(b.href) === path"),
 ("href-userinfo","b.href === path","(b.href === path || b.href === 'http://user@localhost' + path)"),
 ("href-query","b.href === path","(b.href === path || b.href === path + '?q')"),
 ("href-fragment","b.href === path","(b.href === path || b.href === path + '#f')"),
 ("href-scheme","b.href === path","(b.href === path || b.href === 'https://localhost' + path)"),
 ("href-host","b.href === path","(b.href === path || b.href === '//localhost' + path)"),
 ("href-port","b.href === path","(b.href === path || b.href === 'http://localhost:8' + path)"),
 ("status-href","b.statusHref === path + '/status.json'","true"),
]:
    href(key,old,new)
for key, old, new in [
 ("back", 'href="/surface">Back to World', 'href="/surface?broken">Back to World'),
 ("task-markup", '{{.Task}}</pre>', '{{.Review}}</pre>'),
]:
    chrome(key,D+"surface_operator_read.go",old,new)
for key, old, new in [
 ("submit-event","request.addEventListener('click',","request.addEventListener('change',"),
 ("review-text","area.textContent = body.textContent;","area.textContent = '';"),
 ("focus-task","(task.readOnly ? request : task).focus();","(task.readOnly ? task : request).focus();"),
 ("restore-draft","if (outstandingOperatorRequest) {","if (false) {"),
 ("navigate-card","if (result.href) return location.assign(result.href);","if (result.href) return location.assign('/surface');"),
 ("keyboard","'dialog, #compose-btn, .state-branch, #verb:has(#operator-task)'","'dialog, #compose-btn, .state-branch'"),
 ("freeze-ui","task.readOnly = !!control.request;","task.readOnly = false;"),
 ("freeze-inflight","task.readOnly = true;","task.readOnly = false;"),
 ("reconnect-label","control.request ? 'Reconnect' : 'Request action'","'Request action'"),
]:
    chrome(key,A,old,new,"function showOperatorReview" if key not in ("keyboard",) else "")
chrome("enter-post",A,"task.value = control.request?.task || '';",
       "task.value = control.request?.task || ''; void control.submit('fixture task');","function showOperatorReview")
chrome("status-proof-leak",D+"surface_operator_read.go","if statusJSON {",
       "if statusJSON { status.Reason = string(row.Artifact)")
chrome("status-task-leak",D+"surface_operator_read.go","if statusJSON {",
       "if statusJSON { status.Reason = string(row.Action)")
card_go=(root/"st-M2"/(D+"surface_operator_read.go")).read_text()
start=card_go.index("var operatorCard")
stop=card_go.index("</body></html>",start)+len("</body></html>")
old=card_go[start:stop]
raw_lookup='.Funcs(template.FuncMap{"rawTask": func(s string) template.HTML { return template.HTML(s) }}).Parse('
chrome("task-unescaped",D+"surface_operator_read.go",old,
       old.replace(".Parse(",raw_lookup,1).replace("{{.Task}}","{{rawTask .Task}}"))
chrome("clock",D+"daemon.go","approveNow := cfg.ApprovalClock","approveNow := (func() int64)(nil)")
# The original iteration-445 mutation source is judge/mut/specs1.py and specs2.py (J22).
carry("J02",D+"surface_operator.go","body.CatalogueGeneration <= 0 || !store.Valid","!store.Valid","host/daemon","TestOperatorConfirmationAdmissionNeighbours",True)
carry("J06",D+"surface_operator_read.go",'status.State = "spawn-refused"','status.State = "launched"',"host/daemon","TestOperatorConfirmProjection",True)
carry("J10",D+"surface_operator.go"," || h.CatalogueGeneration <= 0 {"," {","host/daemon","TestOperatorConfirmationAdmissionNeighbours",True)
carry("J18",D+"surface_home_summary.go","needs = -1","needs += 0","host/daemon","TestOperatorConfirmationAdmissionNeighbours",True)
carry("J19",D+"surface_operator.go","\td.humanMu.Lock()\n\tdefer d.humanMu.Unlock()\n","","host/daemon","TestOperatorConfirmationAdmissionNeighbours",True)
carry("J22",D+"surface_operator_read.go","|| !bytes.Equal(digest[:], row.Digest)","|| (false && !bytes.Equal(digest[:], row.Digest))","host/daemon","TestOperatorConfirmProjection",True)
carry("E15","host/store/operator_cas.go","backup != binding.BackupEligible","false","host/store","TestOperatorBackupEligibilityDrift")
carry("E19","host/agentrun/process_darwin.go","p.P_stat != zombie","true","host/agentrun","TestDarwinZombieIdentity")
(root/"mutants.json").write_text(json.dumps(items,indent=2)+"\n")
(root/"mutant-diff-ids.txt").write_text("\n".join(k for k,v in items.items() if v["kind"]=="diff")+"\n")
(root/"mutant-carried-ids.txt").write_text("\n".join(k for k,v in items.items() if v["kind"]=="carried")+"\n")
````

````python patches.py
import pathlib
import re
import subprocess
import sys

root = pathlib.Path(__file__).resolve().parent
for old, new in (("base-tree", "st-M1"), ("st-M1", "st-M2"), ("base-tree", "st-M2")):
    m = new.removeprefix("st-")
    tag = m if old == "st-M1" or m == "M1" else m + "-cumulative"
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

````text whole-mutant-ids.txt
E15
E19
J06
J22
J02
J10
J18
J19
R-assets-js
R-host-js
C-confirm-shape
C-report-object
C-passkey-hash
B-operation-reuse
B-operation-stale-after
B-href-exact
B-review-headless
B-operation-busy
````

## 7. Apply, accept and snapshot

Run from the supplied clean sprint worktree. Apply M1 then M2. During executor use, the cumulative snapshot is `$STATE/snap/<M>`. The snapshot excludes exactly the two plan paths, carries `STATUS` and `DELETED`, and verifies `MANIFEST.sha256` over every copied file. The controller runs the per-milestone commit gate below. Rehearsal uses the same block with `REHEARSAL=1` and a fresh archive.

````bash
set -euo pipefail
B=0e32a8d7de0bb524678a3c443400cf05d6ec47eb
STATE="$HOME/.ailang/state/world-iter447"
S="$STATE/plan-scratch"
WT=${WT:-/Users/voightkampff/dev/sunholo-data/.wt-world-iter447-row243-v4s8}
PLAN=${PLAN:-$WT/design_docs/planned/w-world-address-bar-verbs-v4s8-sprint-plan.md}
TREE=${TREE:-$WT}
TOOLS=${TOOLS:-$S}
MILESTONES=${MILESTONES:-M1 M2}
MODE=${MODE:-fast}
REHEARSAL=${REHEARSAL:-0}
export GIT_OPTIONAL_LOCKS=0
mkdir -p "$TOOLS"
grab() {
  awk -v fence="$1" '$0 == fence {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$2"
  test -s "$2"
}
for tool in accept.sh gates.sh snapshot.sh selectors.sh callers.sh outside.sh commit-guards.sh whole-mutants.sh; do
  grab '````bash '"$tool" "$TOOLS/$tool"
done
for tool in bounded.py failset.py verify-scratch.py mutate.py mutant-spec.py patches.py; do
  grab '````python '"$tool" "$TOOLS/$tool"
done
grab '````text whole-mutant-ids.txt' "$TOOLS/whole-mutant-ids.txt"
if [ "$REHEARSAL" = 1 ]; then
  test ! -e "$TREE"
  mkdir -p "$TREE"
  git -C "$WT" archive "$B" | tar -x -C "$TREE"
else
  git -C "$TREE" merge-base --is-ancestor "$B" HEAD
fi
for M in $MILESTONES; do
  case "$M" in M1|M2) ;; *) exit 2;; esac
  grab '````diff '"$M" "$TOOLS/$M.patch"
  (cd "$TREE" && /usr/bin/patch -p1 -E -s < "$TOOLS/$M.patch")
  bash "$TOOLS/accept.sh" "$TREE" "$M" "$MODE"
  if [ "$REHEARSAL" = 1 ]; then
    case "$M" in M1) tag=M1;; M2) tag=M2-cumulative;; esac
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


The worker's completion marker is `APPLY_ACCEPT_SNAPSHOT_OK <M> <mode>`, with `ACCEPT_OK <M> <mode>` printed before the snapshot. A missing marker or nonzero rc is an unfinished milestone.

### Controller commands

After the executor snapshots both stages, use unique, fresh logs and run the outside gates serially. This leg is the real listener and Chrome gate. The snapshot-backed commit script checks gofmt, build and vet. JS syntax and executable JS tests are in `accept.sh`. Run the stage and baseline commands before attribution, sequential with mutation runs.

````bash
set -euo pipefail
B=0e32a8d7de0bb524678a3c443400cf05d6ec47eb
STATE="$HOME/.ailang/state/world-iter447"
S="$STATE/plan-scratch"
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
bash "$S/gates.sh" "$S/base-tree" "$S/outside-pristine" M2 outside-base
L=$(mktemp -d "$S/pristine-daemon-race.XXXXXX")
mkdir "$L/state" "$L/cache"
(
  cd "$S/base-tree"
  export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
  unset GOFLAGS WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
  go test -race -v -p 1 -count=2 -timeout 170s ./host/daemon/ \
    -run '^(TestOperatorConfirmProjection|TestOperatorConfirmationAdmissionNeighbours)$'
) | tee "$L/raw"
bash "$S/outside.sh" M1 "$S/st-M1"
bash "$S/outside.sh" M2 "$S/st-M2"
bash "$S/gates.sh" "$S/st-M1" "$S/outside-gates-M1" M1 outside
bash "$S/gates.sh" "$S/st-M2" "$S/outside-gates-M2" M2 outside
bash "$S/whole-mutants.sh" "$S/outside-mutations"
bash "$STATE/ctl_roundtrip.sh" M1 M2
````

Use this exact named-test command at each snapshot commit. M2 includes the measured fail-set guard which executes the actual `go test -count=1 ./host/verifygate/ ./host/broker/` command in isolated state. Whole daemon/cmd legs stay with the separately recorded gates above.

````bash
STATE="$HOME/.ailang/state/world-iter447"
S="$STATE/plan-scratch"
"$STATE/commit_ms.sh" M1 "$S/M1.commitmsg" 'export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"; S="$HOME/.ailang/state/world-iter447/plan-scratch"; bash "$S/accept.sh" "$PWD" M1 outside'
"$STATE/commit_ms.sh" M2 "$S/M2.commitmsg" 'export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1 WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang" WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1; S="$HOME/.ailang/state/world-iter447/plan-scratch"; bash "$S/accept.sh" "$PWD" M2 outside && bash "$S/commit-guards.sh" "$PWD" "$S/commit-guards-M2-$(date +%s)"'
````

### Rehearsal

The §7 block was extracted with awk, compared with the exact source, and executed verbatim against a fresh archive under `$S/rt-apply/final/export`. The receipt is `$S/rt-apply/final/receipt.raw`.

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
```

All six `diff -r` comparisons for cmd, host and docs were byte-identical. The root has no top-level app folder. Snapshot manifests verified against both cumulative stages.

## Appendix P — patches

### M1 — Standalone card and carried pins

````diff M1
diff --git a/host/agentrun/process_darwin.go b/host/agentrun/process_darwin.go
index de1df0c..670711c 100644
--- a/host/agentrun/process_darwin.go
+++ b/host/agentrun/process_darwin.go
@@ -8,10 +8,14 @@ import (
 
 // ProcessIdentity reads kernel birth time and excludes a reused PID.
 func ProcessIdentity(pid int) (string, bool) {
+	return darwinProcessIdentity(pid, unix.SysctlKinfoProc)
+}
+
+func darwinProcessIdentity(pid int, read func(string, ...int) (*unix.KinfoProc, error)) (string, bool) {
 	if pid <= 0 {
 		return "", false
 	}
-	row, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
+	row, err := read("kern.proc.pid", pid)
 	if err != nil || int(row.Proc.P_pid) != pid {
 		return "", false
 	}
diff --git a/host/agentrun/process_darwin_test.go b/host/agentrun/process_darwin_test.go
index 8c78989..1cad53b 100644
--- a/host/agentrun/process_darwin_test.go
+++ b/host/agentrun/process_darwin_test.go
@@ -3,6 +3,8 @@ package agentrun
 import (
 	"os"
 	"testing"
+
+	"golang.org/x/sys/unix"
 )
 
 func TestDarwinProcessIdentity(t *testing.T) {
@@ -20,3 +22,18 @@ func TestDarwinProcessIdentity(t *testing.T) {
 		}
 	}
 }
+
+func TestDarwinZombieIdentity(t *testing.T) {
+	var row unix.KinfoProc
+	row.Proc.P_pid = 123
+	row.Proc.P_starttime.Sec = 456
+	read := func(string, ...int) (*unix.KinfoProc, error) { return &row, nil }
+	row.Proc.P_stat = 2 // SRUN in Darwin <sys/proc.h>
+	if birth, live := darwinProcessIdentity(123, read); birth != "456.0" || !live {
+		t.Fatalf("live control: %q %t", birth, live)
+	}
+	row.Proc.P_stat = 5 // SZOMB in Darwin <sys/proc.h>
+	if birth, live := darwinProcessIdentity(123, read); birth != "456.0" || live {
+		t.Fatalf("a zombie was live: %q %t", birth, live)
+	}
+}
diff --git a/host/daemon/daemon.go b/host/daemon/daemon.go
index 2ebec16..1ceff23 100644
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -1171,6 +1171,8 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("GET /surface/confirm/{id}", d.approvalRoute(d.handleSurfaceConfirmPage))
 	mux.HandleFunc("POST /surface/confirm/{id}", d.approvalRoute(d.handleSurfaceConfirm))
 	mux.HandleFunc("GET /surface/confirm.js", d.approvalRoute(d.handleSurfaceConfirmScript))
+	mux.HandleFunc("GET /surface/operator-card.js", d.approvalRoute(operatorScript))
+	mux.HandleFunc("GET /surface/operator-card.css", d.approvalRoute(operatorStyle))
 	mux.HandleFunc("GET /surface/confirmations/{id}", d.approvalRoute(d.handleSurfaceConfirmation))
 	mux.HandleFunc("GET /surface/verb.json", d.approvalRoute(d.handleSurfaceVerb))         // row 243 V1: the bar's verbs, read-only
 	mux.HandleFunc("GET /surface/profiles.json", d.approvalRoute(d.handleSurfaceProfiles)) // row 269 Phase 1: read-only exec profiles
diff --git a/host/daemon/surface_operator_browser_test.go b/host/daemon/surface_operator_browser_test.go
new file mode 100644
index 0000000..d7590a1
--- /dev/null
+++ b/host/daemon/surface_operator_browser_test.go
@@ -0,0 +1,46 @@
+package daemon
+
+import (
+	"bytes"
+	"net/http"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+func TestOperatorCardAssets(t *testing.T) {
+	f := newOperatorHTTP(t, time.Minute)
+	before := f.evidence()
+	s := f.request(opRequest(1))
+	seeded := f.evidence()
+	if seeded.actions != before.actions+1 || seeded.operations != before.operations+1 {
+		t.Fatalf("pending row control: %+v to %+v", before, seeded)
+	}
+	card := f.do("GET", s.Href, nil, nil)
+	for _, asset := range []struct {
+		path, media string
+		body        []byte
+	}{
+		{"/surface/operator-card.js", "text/javascript; charset=utf-8", surfacepage.OperatorScript},
+		{"/surface/operator-card.css", "text/css; charset=utf-8", surfacepage.OperatorStyle},
+	} {
+		if !bytes.Contains(card.Body.Bytes(), []byte(asset.path)) {
+			t.Fatalf("card omitted %s", asset.path)
+		}
+		r := f.do("GET", asset.path, nil, nil)
+		if r.Code != 200 || !bytes.Equal(r.Body.Bytes(), asset.body) || r.Header().Get("Content-Type") != asset.media ||
+			r.Header().Get("X-Content-Type-Options") != "nosniff" || r.Header().Get("Content-Security-Policy") != "default-src 'none'" {
+			t.Fatalf("asset %s: %d %v", asset.path, r.Code, r.Header())
+		}
+		if r := f.do("POST", asset.path, nil, nil); r.Code != http.StatusMethodNotAllowed {
+			t.Fatalf("POST %s: %d", asset.path, r.Code)
+		}
+		if r := f.do("GET", asset.path, nil, func(r *http.Request) { r.Host = "other.invalid" }); r.Code != http.StatusMisdirectedRequest {
+			t.Fatalf("bad Host %s: %d", asset.path, r.Code)
+		}
+	}
+	if after := f.evidence(); after != seeded {
+		t.Fatalf("asset reads changed action or operation counts: %+v to %+v", seeded, after)
+	}
+}
diff --git a/host/daemon/surface_operator_confirm_test.go b/host/daemon/surface_operator_confirm_test.go
index 3fc407d..51ad30f 100644
--- a/host/daemon/surface_operator_confirm_test.go
+++ b/host/daemon/surface_operator_confirm_test.go
@@ -608,13 +608,13 @@ func requireOperatorCeremony(t *testing.T, card *httptest.ResponseRecorder, s op
 }
 
 func TestOperatorConfirmProjection(t *testing.T) {
-	f := newOperatorSurface(t, nil, time.Minute)
-	raw := f.postOperator(opRequest(1))
+	f := newOperatorHTTP(t, time.Minute)
+	raw := f.do("POST", operatorPath, opJSON(opRequest(1)), nil)
 	var s operatorStatus
 	if err := json.Unmarshal(raw.Body.Bytes(), &s); err != nil || raw.Code != 200 || s.ID == "" {
 		t.Fatalf("projection request fixture: %d %s", raw.Code, raw.Body)
 	}
-	read := func(path string) *httptest.ResponseRecorder { return f.serve("GET", path, apLoopback, nil) }
+	read := func(path string) *httptest.ResponseRecorder { return f.do("GET", path, nil, nil) }
 	original := f.row(s.ID).Digest
 	f.confirmSQL("UPDATE operator_actions SET digest=zeroblob(32) WHERE id=?", s.ID)
 	f.refuse(f.evidence(), read(s.Href), 409, "Conflict")
@@ -659,23 +659,29 @@ func TestOperatorConfirmProjection(t *testing.T) {
 
 func TestOperatorConfirmationAdmissionNeighbours(t *testing.T) {
 	f := newOperatorSurface(t, nil, time.Minute)
-	checked := 0
+	var checked atomic.Int64
+	var checking atomic.Bool
+	checking.Store(true)
 	f.d.approveNow = func() int64 {
-		checked++
-		if f.d.humanMu.TryLock() {
-			f.d.humanMu.Unlock()
-			t.Error("request reached its store read outside the human writer lock")
+		if checking.Load() {
+			checked.Add(1)
+			if f.d.humanMu.TryLock() {
+				f.d.humanMu.Unlock()
+				t.Error("request reached its store read outside the human writer lock")
+			}
 		}
 		return apT0
 	}
+	f.server = httptest.NewServer(f.handler)
+	t.Cleanup(f.server.Close)
 	bad := opRequest(1)
 	bad.CatalogueGeneration = 0
 	before := f.evidence()
-	f.refuse(before, f.postOperator(bad), 400, "BadRequest")
-	if r := f.postOperator(opRequest(1)); r.Code != 200 || checked < 2 {
-		t.Fatalf("request positive control did not exercise lock: %d (%d checks)", r.Code, checked)
+	f.refuse(before, f.do("POST", operatorPath, opJSON(bad), nil), 400, "BadRequest")
+	if r := f.do("POST", operatorPath, opJSON(opRequest(1)), nil); r.Code != 200 || checked.Load() < 2 {
+		t.Fatalf("request positive control did not exercise lock: %d (%d checks)", r.Code, checked.Load())
 	}
-	f.d.approveNow = f.now.Load
+	checking.Store(false)
 	ctx := boundedTestContext(t)
 	ro, err := store.OpenCredentialsReadOnly(f.d.credentialsPath)
 	if err != nil {
@@ -686,16 +692,16 @@ func TestOperatorConfirmationAdmissionNeighbours(t *testing.T) {
 		t.Fatalf("catalogue control: %v", err)
 	}
 	f.confirmSQL("UPDATE operator_runner_state SET catalogue_generation=0")
-	if h, _, err := f.d.operatorReview(ctx, ro, "sol", apT0); !errors.Is(err, store.ErrOperatorRunner) || h.CatalogueGeneration != 0 {
-		t.Fatalf("invalid generation crossed review: %v, %v", h, err)
-	}
 	bar := "/surface/verb.json?q=" + url.QueryEscape("episode start ep-2 --preset sol --launch --headless")
-	if r := f.serve("GET", bar, apLoopback, nil); r.Code != 200 || !strings.Contains(r.Body.String(), `"outcome":"terminal"`) {
+	if r := f.do("GET", bar, nil, nil); r.Code != 200 || !strings.Contains(r.Body.String(), `"outcome":"terminal"`) {
 		t.Fatalf("invalid generation reached the bar review: %d %s", r.Code, r.Body)
 	}
+	if h, _, err := f.d.operatorReview(ctx, ro, "sol", apT0); !errors.Is(err, store.ErrOperatorRunner) || h.CatalogueGeneration != 0 {
+		t.Fatalf("invalid generation crossed review: %v, %v", h, err)
+	}
 	summary := func() surfacepage.HomeSummary {
 		t.Helper()
-		r := f.serve("GET", "/surface/home-summary.json", apLoopback, nil)
+		r := f.do("GET", "/surface/home-summary.json", nil, nil)
 		var out surfacepage.HomeSummary
 		if r.Code != 200 || json.Unmarshal(r.Body.Bytes(), &out) != nil || len(out.Items) != 1 || out.Items[0].Label != "Needs you" {
 			t.Fatalf("home summary: %d %s", r.Code, r.Body)
@@ -707,7 +713,7 @@ func TestOperatorConfirmationAdmissionNeighbours(t *testing.T) {
 	}
 	f.confirmSQL("UPDATE operator_actions SET expires_at='unreadable' WHERE episode_id='ep-1'")
 	before = f.evidence()
-	r := f.serve("GET", "/surface/requests.json", apLoopback, nil)
+	r := f.do("GET", "/surface/requests.json", nil, nil)
 	var requests surfacepage.RequestsJSON
 	if r.Code != 200 || json.Unmarshal(r.Body.Bytes(), &requests) != nil || requests.Unavailable != "" || requests.OperatorActionsUnavailable == "" {
 		t.Fatalf("fixture must fail only the operator pending read: %d %s", r.Code, r.Body)
diff --git a/host/daemon/surface_operator_read.go b/host/daemon/surface_operator_read.go
index 3ebd00a..5e298f8 100644
--- a/host/daemon/surface_operator_read.go
+++ b/host/daemon/surface_operator_read.go
@@ -10,6 +10,7 @@ import (
 	"github.com/sunholo-data/ailang-world/host/agentrun"
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
 	"github.com/sunholo-data/ailang-world/host/verbs"
 )
 
@@ -55,14 +56,25 @@ func operatorStatusFor(row store.OperatorActionRow, now int64) operatorStatus {
 }
 
 var operatorCard = template.Must(template.New("operator action").Parse(`<!doctype html>
-<html lang="en"><meta charset="utf-8"><title>Operator action</title>
-<body><h1>{{.Status.ID}}</h1><p>{{.Status.State}}</p>{{if .Status.Reason}}<p>{{.Status.Reason}}</p>{{end}}
-{{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre><h2>Task</h2><pre>{{.Task}}</pre>{{end}}
-{{if .Ceremony}}<h2>Confirm this action</h2>
-<p>Verify the exact review and task. Submit a passkey assertion or SSH signature over this action using the ceremony below.</p>
-<pre id="operator-ceremony">{{.Ceremony}}</pre>{{end}}
+<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
+<title>Operator action</title><link rel="stylesheet" href="/surface/operator-card.css">
+<script defer src="/surface/operator-card.js"></script></head><body>
+<h1 id="operator-title" tabindex="-1">Operator action {{.Status.ID}}</h1>
+<p id="operator-state" data-state="{{.Status.State}}">{{.Status.State}}</p>
+<p id="operator-progress" role="status">Status at {{.Status.Now}}{{if .Status.Reason}} · {{.Status.Reason}}{{end}}</p>
+{{if .Review}}<h2>Review</h2><pre>{{.Review}}</pre><h2>Task</h2><pre id="operator-task">{{.Task}}</pre>{{end}}
+<section id="operator-controls" hidden aria-label="Confirm this action">
+<button id="operator-key" type="button" hidden>Confirm with passkey</button>
+<details><summary>Confirm with SSH signature</summary>
+<p>Review the action and sign the exact message with an enrolled SSH key in the ceremony's namespace.</p>
+{{if .Ceremony}}<pre id="operator-ceremony">{{.Ceremony}}</pre>{{end}}
+<label for="operator-signature">Armored signature</label><textarea id="operator-signature" rows="4"></textarea>
+<button id="operator-sign" type="button">Submit signature</button></details></section>
+<p id="operator-answer" role="status" tabindex="-1"></p>
+<button id="operator-refresh" type="button">Refresh status</button>
 {{if .Status.StatusHref}}<p><a href="{{.Status.StatusHref}}">Action status (JSON)</a></p>{{end}}
-{{if .Status.ReportHref}}<p><a href="{{.Status.ReportHref}}">Episode report</a></p>{{end}}
+<p><a id="operator-report" {{if .Status.ReportHref}}href="{{.Status.ReportHref}}"{{else}}hidden{{end}}>Episode report</a></p>
+<p><a id="operator-back" href="/surface">Back to World</a></p>
 </body></html>`))
 
 func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id string, statusJSON bool) {
@@ -132,3 +144,11 @@ func (d *Daemon) handleOperatorRead(w http.ResponseWriter, r *http.Request, id s
 	setSurfaceHeaders(w)
 	_ = operatorCard.Execute(w, card)
 }
+
+func operatorScript(w http.ResponseWriter, _ *http.Request, _ []string) {
+	writeSurfaceAsset(w, "text/javascript; charset=utf-8", surfacepage.OperatorScript)
+}
+
+func operatorStyle(w http.ResponseWriter, _ *http.Request, _ []string) {
+	writeSurfaceAsset(w, "text/css; charset=utf-8", surfacepage.OperatorStyle)
+}
diff --git a/host/store/operator_backup_test.go b/host/store/operator_backup_test.go
new file mode 100644
index 0000000..fa27868
--- /dev/null
+++ b/host/store/operator_backup_test.go
@@ -0,0 +1,75 @@
+package store
+
+import (
+	"errors"
+	"reflect"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+type storeOperatorProof struct {
+	signed *approval.SignedOperatorAction
+}
+
+func (p storeOperatorProof) OperatorConfirmation() *approval.SignedOperatorAction { return p.signed }
+
+func TestOperatorBackupEligibilityDrift(t *testing.T) {
+	c := openTestCredentials(t)
+	ctx := testCtx(t)
+	key := approvaltest.NewPasskey(approval.AlgEdDSA, 82)
+	a, p := testPasskey("fixture", 82)
+	challenge := [32]byte{1}
+	enrolled, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: a.Name, Origin: p.Origin, Challenge: challenge,
+		Response: key.Register(p.Origin, p.RPID, challenge), Now: a.EnrolledAt})
+	if err != nil {
+		t.Fatal(err)
+	}
+	a.PublicKey, p.CredentialID, p.Alg, p.AAGUID = enrolled.Passkey.Key.Canonical(), key.CredentialID(), approval.AlgEdDSA, enrolled.Passkey.AAGUID[:]
+	if err := c.EnrolPasskey(testCtx(t), a, p); err != nil {
+		t.Fatal(err)
+	}
+	req := operatorRequest(t, c, 1, "drift", "start")
+	row := putOperator(t, c, req, 1000)
+	cas := OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
+	hash, err := approval.OperatorActionMessageHash(cas.Digest, req.Challenge, req.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	artifact := key.Assert(p.Origin, p.RPID, hash, 2)
+	signed, err := approval.ConfirmOperatorAction(approval.OperatorConfirmInput{Action: req.Action, Challenge: req.Challenge, Artifact: artifact,
+		Format: a.Format, Approver: enrolled, Admit: approval.AdmitInput{Now: 1001, CreatedAt: req.CreatedAt, ExpiresAt: req.ExpiresAt,
+			StoredDigest: cas.Digest, BytesDigest: cas.Digest}})
+	if err != nil {
+		t.Fatal(err)
+	}
+	setBackup := func(eligible bool) {
+		t.Helper()
+		if _, err := c.db.ExecContext(testCtx(t), `UPDATE webauthn_credentials SET backup_eligible = ? WHERE name = ?`, eligible, a.Name); err != nil {
+			t.Fatal(err)
+		}
+	}
+	setBackup(true)
+	read := func() [2]int64 {
+		t.Helper()
+		var seen [2]int64
+		if err := c.db.QueryRowContext(ctx, `SELECT COALESCE(last_used_at, 0), sign_count FROM approvers JOIN webauthn_credentials USING(name) WHERE name = ?`, a.Name).
+			Scan(&seen[0], &seen[1]); err != nil {
+			t.Fatal(err)
+		}
+		return seen
+	}
+	before := read()
+	if _, err := c.ConfirmOperatorAction(testCtx(t), cas, artifact, storeOperatorProof{signed}, 1001); !errors.Is(err, ErrApproverKeyChanged) {
+		t.Fatalf("backup drift: %v", err)
+	}
+	if !reflect.DeepEqual(row, getOperator(t, c, row.ID, 1001)) || read() != before {
+		t.Fatal("backup drift wrote the action or key counter")
+	}
+	operatorCounts(t, c, 1, 1)
+	setBackup(false)
+	if confirmed, err := c.ConfirmOperatorAction(testCtx(t), cas, artifact, storeOperatorProof{signed}, 1001); err != nil || confirmed.State != "confirmed" || read() != [2]int64{1001, 2} {
+		t.Fatalf("matching control: state=%s key=%v err=%v", confirmed.State, read(), err)
+	}
+}
diff --git a/host/surfacepage/operator/card.css b/host/surfacepage/operator/card.css
new file mode 100644
index 0000000..2884359
--- /dev/null
+++ b/host/surfacepage/operator/card.css
@@ -0,0 +1,5 @@
+body { font: 1rem/1.5 system-ui, sans-serif; margin: 1rem auto; padding: 0 1rem; max-width: 48rem; overflow-wrap: anywhere; }
+pre { white-space: pre-wrap; padding: .6rem; background: #f1f3f5; color: #17202a; }
+textarea { box-sizing: border-box; width: 100%; }
+button, summary, a { min-height: 2.75rem; }
+button { margin: .25rem 0; }
diff --git a/host/surfacepage/operator/card.js b/host/surfacepage/operator/card.js
new file mode 100644
index 0000000..be8014a
--- /dev/null
+++ b/host/surfacepage/operator/card.js
@@ -0,0 +1,104 @@
+(() => {
+  'use strict';
+  const cardPath = /^\/surface\/operator-actions\/(oa-[0-9a-f]{32})(?![\s\S])/;
+  const episodeID = /^[a-z0-9][a-z0-9-]{0,63}(?![\s\S])/;
+  const hex64 = /^[0-9a-f]{64}(?![\s\S])/;
+  const hexBytes = /^(?:[0-9a-f]{2})+(?![\s\S])/;
+  const pathID = path => typeof path === 'string' ? cardPath.exec(path)?.[1] : undefined;
+  const states = new Set(['requested', 'confirmed', 'claimed', 'running', 'awaiting-pr-confirmation',
+    'succeeded', 'failed', 'expired', 'cancelled', 'uncertain', 'launched', 'stopped', 'capped', 'spawn-refused']);
+  function statusFor(s, path) {
+    if (!s || s.v !== 1 || s.id !== pathID(path) || s.href !== path || s.statusHref !== path + '/status.json' ||
+        typeof s.episode !== 'string' || !episodeID.test(s.episode) || !['start', 'finish'].includes(s.kind) ||
+        !states.has(s.state) || ![s.now, s.createdAt, s.expiresAt].every(Number.isSafeInteger) ||
+        ![s.terminalAt, s.runnerSeenAt].every(n => n === undefined || Number.isSafeInteger(n)) ||
+        (s.reason !== undefined && typeof s.reason !== 'string')) return null;
+    const report = '/v1/episodes/' + s.episode + '/report';
+    return { ...s, reportHref: s.reportHref === report ? report : '' };
+  }
+  function ceremonyFor(c, path) {
+    return c && c.post === path + '/confirm' && typeof c.digest === 'string' && hex64.test(c.digest) &&
+      typeof c.messageHex === 'string' && hexBytes.test(c.messageHex) ? c : null;
+  }
+  async function passkeyProof(c, origin, hostname, crypto, credentials, signal) {
+    const group = Array.isArray(c.passkeys) && c.passkeys.find(k => k.origin === origin && k.rp === hostname);
+    if (!group || !Array.isArray(group.allow) || group.allow.length === 0) throw new Error('No passkey for this origin');
+    const challenge = await crypto.subtle.digest('SHA-256', Uint8Array.from(c.messageHex.match(/../g), x => parseInt(x, 16)));
+    const allowCredentials = group.allow.map(id => ({ type: 'public-key', id: Uint8Array.from(atob(id.replace(/-/g, '+').replace(/_/g, '/')), ch => ch.charCodeAt(0)) }));
+    const key = await credentials.get({ publicKey: { challenge, rpId: group.rp, allowCredentials, userVerification: 'required', timeout: 15000 }, signal });
+    return { webauthn: key.toJSON() };
+  }
+  if (typeof module !== 'undefined') module.exports = { pathID, statusFor, ceremonyFor, passkeyProof };
+  if (typeof document === 'undefined') return;
+
+  const $ = id => document.getElementById(id);
+  const path = location.pathname;
+  if (!pathID(path)) return;
+  let ceremony = null;
+  try { ceremony = ceremonyFor(JSON.parse($('operator-ceremony')?.textContent), path); } catch (_) { /* status stays readable */ }
+  let pending = false, polling, last = $('operator-progress').textContent, state = $('operator-state').dataset.state;
+  const controls = $('operator-controls'), progress = $('operator-progress'), answer = $('operator-answer');
+  function buttons() {
+    controls.hidden = !ceremony || state !== 'requested';
+    $('operator-key').hidden = !Array.isArray(ceremony?.passkeys) ||
+      !ceremony.passkeys.some(k => k.origin === location.origin && k.rp === location.hostname);
+    for (const id of ['operator-key', 'operator-sign', 'operator-refresh']) $(id).disabled = pending;
+  }
+  async function read(url, options) {
+    const res = await fetch(url, { ...options, redirect: 'error', credentials: 'omit', cache: 'no-store',
+      signal: AbortSignal.timeout(15000) });
+    return { code: res.status, body: await res.json() };
+  }
+  async function refresh() {
+    if (pending) return;
+    pending = true;
+    clearTimeout(polling);
+    buttons();
+    try {
+      const { code, body } = await read(path + '/status.json');
+      const s = code === 200 && statusFor(body, path);
+      if (!s) throw new Error('Invalid status');
+      state = s.state;
+      $('operator-state').textContent = state;
+      last = progress.textContent = [s.reason, 'Status at ' + s.now, 'Requested ' + s.createdAt, 'Expires ' + s.expiresAt,
+        s.runnerSeenAt && 'Runner seen ' + s.runnerSeenAt, s.terminalAt && 'Completed ' + s.terminalAt].filter(Boolean).join(' · ');
+      $('operator-report').hidden = !s.reportHref;
+      if (s.reportHref) $('operator-report').href = s.reportHref;
+      else $('operator-report').removeAttribute('href');
+    } catch (_) {
+      progress.textContent = 'Disconnected. ' + last + '. Use Refresh status to reconnect.';
+    } finally {
+      pending = false;
+      buttons();
+      polling = setTimeout(refresh, 5000);
+    }
+  }
+  async function confirm(proof) {
+    if (pending || !ceremony || state !== 'requested') return;
+    pending = true;
+    buttons();
+    answer.textContent = 'Waiting for confirmation…';
+    answer.focus();
+    try {
+      const artifact = await proof();
+      const { code, body } = await read(path + '/confirm', { method: 'POST', headers: { 'Content-Type': 'application/json' },
+        body: JSON.stringify({ v: 1, digest: ceremony.digest, ...artifact }) });
+      answer.textContent = code === 200 && statusFor(body, path) && ['confirmed', 'already-confirmed'].includes(body.result) && typeof body.approvedBy === 'string' && body.approvedBy !== ''
+        ? 'Confirmed by ' + body.approvedBy
+        : 'Not confirmed (HTTP ' + code + '): ' + (typeof body?.error === 'string' ? body.error + ' ' : '') +
+          (typeof body?.reason === 'string' ? body.reason : 'unreadable confirmation response');
+    } catch (_) {
+      answer.textContent = 'Confirmation unavailable. Read the action status before trying again.';
+    } finally {
+      pending = false;
+      await refresh();
+    }
+  }
+  $('operator-key').addEventListener('click', () => confirm(() => passkeyProof(ceremony, location.origin, location.hostname,
+    crypto, navigator.credentials, AbortSignal.timeout(15000))));
+  $('operator-sign').addEventListener('click', () => confirm(() => ({ sshsig: $('operator-signature').value })));
+  $('operator-refresh').addEventListener('click', refresh);
+  buttons();
+  $('operator-title').focus();
+  refresh();
+})();
diff --git a/host/surfacepage/operator_card.go b/host/surfacepage/operator_card.go
new file mode 100644
index 0000000..6d3b2a5
--- /dev/null
+++ b/host/surfacepage/operator_card.go
@@ -0,0 +1,11 @@
+package surfacepage
+
+import _ "embed"
+
+// OperatorScript and OperatorStyle are the standalone operator card assets.
+//
+//go:embed operator/card.js
+var OperatorScript []byte
+
+//go:embed operator/card.css
+var OperatorStyle []byte
diff --git a/host/surfacepage/operator_test.go b/host/surfacepage/operator_test.go
new file mode 100644
index 0000000..1d14c05
--- /dev/null
+++ b/host/surfacepage/operator_test.go
@@ -0,0 +1,168 @@
+package surfacepage
+
+import (
+	"context"
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/hex"
+	"encoding/json"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+func runOperatorScript(t *testing.T, asset []byte, harness string) {
+	t.Helper()
+	node := os.Getenv("WORLD_EXEC_NODE")
+	if node == "" {
+		var err error
+		node, err = exec.LookPath("node")
+		if err != nil {
+			t.Skip("node is required to execute the served script")
+		}
+	}
+	path := filepath.Join(t.TempDir(), "asset.js")
+	if err := os.WriteFile(path, asset, 0600); err != nil {
+		t.Fatal(err)
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
+	defer cancel()
+	cmd := exec.CommandContext(ctx, node, "-e", harness, path)
+	if out, err := cmd.CombinedOutput(); err != nil {
+		t.Fatalf("%v\n%s", err, out)
+	}
+}
+
+func TestOperatorCardStatusDisconnected(t *testing.T) {
+	const origin, rp = "http://localhost:7644", "localhost"
+	digest := sha256.Sum256([]byte("operator-card fixture"))
+	challenge := [32]byte{9}
+	message, err := approval.OperatorActionMessage(digest, challenge, 1600)
+	if err != nil {
+		t.Fatal(err)
+	}
+	hash := sha256.Sum256(message)
+	key := approvaltest.NewPasskey(approval.AlgES256, 93)
+	fixture, err := json.Marshal(map[string]any{
+		"digest": hex.EncodeToString(digest[:]), "message": hex.EncodeToString(message), "hash": hex.EncodeToString(hash[:]),
+		"credential": base64.RawURLEncoding.EncodeToString(key.CredentialID()),
+		"assertion":  json.RawMessage(key.Assert(origin, rp, hash, 2)),
+	})
+	if err != nil {
+		t.Fatal(err)
+	}
+	runOperatorScript(t, OperatorScript, "const fixture = "+string(fixture)+";\n"+`
+const fs = require('node:fs'), vm = require('node:vm'), assert = require('node:assert/strict');
+const path = '/surface/operator-actions/oa-' + 'a'.repeat(32), origin = 'http://localhost:7644';
+const nodes = {}, events = {}, sent = [];
+const node = id => nodes[id] ??= { textContent: '', dataset: {state: 'requested'},
+  focus() { events.focus = id; }, addEventListener(k, f) { events[id] = f; }, removeAttribute() { delete this.href; } };
+const status = {v:1, id:path.split('/').pop(), episode:'e-1', kind:'start', href:path, statusHref:path+'/status.json',
+  state:'requested', now:1001, createdAt:1000, expiresAt:1600};
+const c = {post:path+'/confirm', digest:fixture.digest, messageHex:fixture.message, passkeys:[{origin:'http://localhost:7645',rp:'localhost',allow:['Ag']},
+  {origin,rp:'not-localhost',allow:['Aw']},{origin,rp:'localhost',allow:[fixture.credential]}]};
+node('operator-ceremony').textContent = JSON.stringify(c);
+const get = { ...status, now:1002 };
+let response = get, code = 200, poll, asked;
+const context = {document:{getElementById:node}, location:{pathname:path,origin,hostname:'localhost'}, atob, AbortSignal,
+  crypto:require('node:crypto').webcrypto, navigator:{credentials:{get: async options => {
+    asked = options;
+    return {toJSON: () => fixture.assertion};
+  }}}, setTimeout: f => (poll = f, 1), clearTimeout: () => {},
+  fetch: async (url, options) => { sent.push([url, options]); if (!response) throw Error('offline');
+    return {status:code,json:async()=>response}; }};
+vm.runInNewContext(fs.readFileSync(process.argv[1], 'utf8'), context);
+(async () => {
+  assert.equal(events.focus,'operator-title');
+  assert.equal(sent.filter(([,o])=>o.method==='POST').length,0);
+  await new Promise(setImmediate);
+  assert.match(node('operator-progress').textContent, /Status at 1002/);
+  assert.equal(node('operator-controls').hidden,false);
+  assert.equal(node('operator-key').hidden,false);
+  assert.equal(sent[0][0],path+'/status.json');
+  assert.equal(sent[0][1].redirect,'error');
+  assert.equal(sent[0][1].credentials,'omit');
+  assert.equal(sent[0][1].cache,'no-store');
+  assert(sent[0][1].signal instanceof AbortSignal);
+  response = null; await poll();
+  assert.match(node('operator-progress').textContent,/Disconnected.*1002/);
+  assert.equal(node('operator-state').textContent,'requested');
+  response = { ...get, state:'invented' }; await poll();
+  assert.match(node('operator-progress').textContent,/Disconnected/);
+  response = { error:'BadRequest', reason:'send a v1 digest and one key artifact' }; code = 400;
+  node('operator-signature').value = 'y';
+  await events['operator-sign']();
+  const [u, p] = sent.find(([,o])=>o.method==='POST');
+  assert.equal(u,path+'/confirm');
+  assert.deepEqual(JSON.parse(p.body),{v:1,digest:c.digest,sshsig:'y'});
+  assert.equal(events.focus,'operator-answer');
+  assert.match(node('operator-answer').textContent,/Not confirmed \(HTTP 400\): BadRequest/);
+  for (const [refused, statusCode] of [['NotApproved',422],['Expired',409]]) {
+    response = {error:refused,reason:'fixture refusal'}; code = statusCode;
+    await events['operator-sign']();
+    assert(node('operator-answer').textContent.includes(refused));
+  }
+  response = {...get,state:'expired'}; code = 200; await poll();
+  assert.equal(node('operator-controls').hidden,true);
+  response = get; await poll();
+  code = 200;
+  for (const edit of [{result:'fake'},{approvedBy:[]},{approvedBy:''},{state:'invented'}]) {
+    response = {...get, result:'confirmed', approvedBy:'webauthn:fixture', ...edit};
+    await events['operator-sign']();
+    assert.match(node('operator-answer').textContent,/Not confirmed/);
+  }
+  response = {...get, state:'confirmed', result:'confirmed', approvedBy:'webauthn:fixture'};
+  await events['operator-key']();
+  assert.equal(asked.publicKey.userVerification,'required');
+  assert.equal(asked.publicKey.rpId,'localhost');
+  assert.equal(Buffer.from(asked.publicKey.allowCredentials[0].id).toString('base64url'),fixture.credential);
+  assert.equal(Buffer.from(asked.publicKey.challenge).toString('hex'),fixture.hash);
+  assert.deepEqual(JSON.parse(sent.filter(([,o])=>o.method==='POST').at(-1)[1].body),{v:1,digest:c.digest,webauthn:fixture.assertion});
+  assert.match(node('operator-answer').textContent,/Confirmed by webauthn:fixture/);
+  assert.equal(node('operator-controls').hidden,true);
+  response = {...get,state:'launched',reportHref:'/v1/episodes/e-1/report'}; await poll();
+  assert.equal(node('operator-report').href,response.reportHref);
+  assert.equal(node('operator-state').textContent,'launched');
+  response = {...response,reportHref:{toString:()=>'/v1/episodes/e-1/report'}}; await poll();
+  assert.equal(node('operator-report').href,undefined);
+  assert.equal(node('operator-report').hidden,true);
+})().catch(e=>{ console.error(e); process.exitCode=1; });
+`)
+}
+
+func TestOperatorCardLinkAllowlist(t *testing.T) {
+	runOperatorScript(t, OperatorScript, `
+const assert = require('node:assert/strict'), {pathID,statusFor,ceremonyFor,passkeyProof} = require(process.argv[1]);
+const id = 'oa-'+'a'.repeat(32), p = '/surface/operator-actions/'+id;
+const s = {v:1,id,episode:'e-1',kind:'start',state:'launched',now:1,createdAt:1,expiresAt:601,href:p,statusHref:p+'/status.json'};
+const h = '/v1/episodes/e-1/report';
+assert.equal(pathID(p),id);
+assert.equal(statusFor({...s,reportHref:h}, p).reportHref,h);
+for (const bad of [p+'x','x'+p,p+'\n', 'x\n'+p, p+'\nx',p.toUpperCase(),p+'A',p.slice(0,-1),p+'a',p+'?x',p+'#x',p+'/status.json',
+  '//localhost'+p,'https://user@localhost:7644'+p,{toString:()=>p}]) assert.equal(pathID(bad),undefined,JSON.stringify(bad));
+for (const bad of [h+'x','x'+h,h+'\n',h.toUpperCase(),h+'?q',h+'#f','//localhost'+h,'https://localhost:2'+h,
+  'http://user@localhost'+h,'/surface/confirm/vr-0123456789abcdef','https://github.com/sunholo-data/ailang-world/pull/604','/v1/episodes/other/report',{toString:()=>h}]) assert.equal(statusFor({...s,reportHref:bad},p).reportHref,'');
+for (const episode of ['A', '_a', '-a', 'a_b', '', 'a'.repeat(65), 'a\n','!a','!\na']) assert.equal(statusFor({...s,episode},p),null);
+for (const [k,v] of Object.entries({v:'1',episode:['e-1'],id:id+'x',href:p+'#',statusHref:{href:p},state:'ok',now:'1',expiresAt:null,createdAt:[],runnerSeenAt:'1',terminalAt:'2',reason:1,kind:'unknown'}))
+  assert.equal(statusFor({...s,[k]:v},p),null,k);
+const c = {post:p+'/confirm',digest:'b'.repeat(64),messageHex:'abcd'};
+assert.equal(ceremonyFor(c,p),c);
+assert.equal(ceremonyFor(null,p),null);
+assert.equal(statusFor(null,p),null);
+(async()=>{
+ const origin='http://localhost:7644', rp='localhost';
+ const get=()=>assert.fail('credential requested without enrolled fixture');
+ for (const passkeys of [null,{},[],[{origin,rp:'other',allow:['AQ']}],
+   [{origin:'http://localhost:1',rp,allow:['AQ']}],[{origin,rp,allow:[]}],[{origin,rp,allow:{}}]]) {
+  await assert.rejects(passkeyProof({...c,passkeys},origin,rp,require('node:crypto').webcrypto,{get}),/No passkey/);
+ }
+})().catch(e=>{console.error(e);process.exitCode=1;});
+for (const [k,v] of [['post',c.post+'#'],['post',{toString:()=>c.post}],['digest','A'.repeat(64)],['digest',1],['digest',[c.digest]],['messageHex','a'],['messageHex',12],['messageHex',[c.messageHex]]])
+  assert.equal(ceremonyFor({...c,[k]:v},p),null,k);
+`)
+}
````

### M2 — App bar, fake browser start and usage

````diff M2
diff --git a/cmd/ailang-worldd/operator_browser_chrome_test.go b/cmd/ailang-worldd/operator_browser_chrome_test.go
new file mode 100644
index 0000000..b92d7a6
--- /dev/null
+++ b/cmd/ailang-worldd/operator_browser_chrome_test.go
@@ -0,0 +1,388 @@
+package main
+
+import (
+	"bytes"
+	"context"
+	"encoding/base64"
+	"encoding/json"
+	"fmt"
+	"io"
+	"net/http"
+	"net/http/httptest"
+	"os"
+	"path/filepath"
+	"reflect"
+	"slices"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/agentrun"
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+const browserStart = "episode start ep --preset p --launch --headless"
+const browserTask = "# Task\n<script>alert(1)</script> ; $(literal) &"
+const browserFakeChild = fakeLaunchMarker + "while [ ! -e release ]; do :; done\nexit 0\n"
+
+func browserPreset(t *testing.T, r *worktreeStartRig, c runnerConfig) (store.OperatorStartReview, *resolvedPreset) {
+	t.Helper()
+	r.editPreset(t, func(p *launchPreset) {
+		p.Model = "gpt-6.2-sol"
+		p.Launch = &presetLaunch{Mode: "interactive", DeadlineS: 60}
+		p.Checks = p.Checks[:1]
+	})
+	ctx, cancel := r.env.operatorContext(activationTestBudget)
+	defer cancel()
+	p, err := loadLaunchPreset(ctx, "p", r.src.presets, r.src.profiles, r.src.workspace)
+	if err != nil {
+		t.Fatal(err)
+	}
+	entry, ok := operatorCatalogueEntry(ctx, c, r.src.workspace, p, r.health)
+	if !ok {
+		t.Fatal("the fake installed preset did not enter the catalogue")
+	}
+	return entry, p
+}
+
+func publishBrowserCatalogue(t *testing.T, r *worktreeStartRig, c runnerConfig, h store.OperatorHeartbeat, entry store.OperatorStartReview) store.OperatorHeartbeat {
+	t.Helper()
+	catalog, err := json.Marshal(map[string]store.OperatorStartReview{"p": entry})
+	if err != nil {
+		t.Fatal(err)
+	}
+	h.Build = daemon.CurrentBuild().ExecutableSHA256
+	h.Owner.Generation = 0 // Publish a new fake owner generation for the real HTTP daemon's build.
+	r.env.operatorCatalogue = func(runnerConfig) []byte { return catalog }
+	if err := r.env.runnerTick(c, &h); err != nil {
+		t.Fatal(err)
+	}
+	fresh, err := runnerHeartbeat(c.Credentials, activationTestBudget)
+	if err != nil || !fresh.Available(h.Build, "start", r.src.clock()) || fresh.CatalogueGeneration <= 0 {
+		t.Fatalf("fake catalogue did not reach the read boundary: %v", err)
+	}
+	return fresh
+}
+
+func TestOperatorFakePresetEligibility(t *testing.T) {
+	t.Setenv("HOME", t.TempDir())
+	r, c, h, _ := operatorActivationFixtureWithTask(t, browserFakeChild, browserTask)
+	entry, p := browserPreset(t, r, c)
+	h = publishBrowserCatalogue(t, r, c, h, entry)
+	var out, stderr syncBuffer
+	if rc := runPreset([]string{"check", "--presets-dir", r.src.presets, "--exec-profiles-dir", r.src.profiles, "p"}, &out, &stderr); rc != exitOK {
+		t.Fatalf("fake preset check: %d %s", rc, &stderr)
+	}
+	t.Logf("preset check: %s", out.String())
+	if !strings.Contains(out.String(), "check short: sandboxed: ") || strings.Contains(out.String(), "UNSANDBOXED") ||
+		!strings.Contains(out.String(), "base dev @ "+p.BaseSHA[:7]) {
+		t.Fatal("resolved preset parity failed")
+	}
+	mode, _, err := launchMode(p.Preset, false, false)
+	if err != nil || mode != "interactive" {
+		t.Fatalf("preset default: %s %v", mode, err)
+	}
+	words, err := verbs.Split(browserStart)
+	if err != nil {
+		t.Fatal(err)
+	}
+	policy := verbs.OperatorPolicy{Start: "A", Finish: "A"}
+	explicit, implicit := verbs.AdmitOperator(words, policy), verbs.AdmitOperator(words[:len(words)-1], policy)
+	if explicit.Intent == nil || !explicit.Intent.Headless || implicit.Class != verbs.OperatorTerminal || implicit.Intent != nil {
+		t.Fatal("only explicit headless may override the interactive fixture")
+	}
+	mode, wall, err := launchMode(p.Preset, explicit.Intent.Headless, false)
+	if err != nil || mode != "headless" || wall != time.Minute {
+		t.Fatalf("explicit override: %s %s %v", mode, wall, err)
+	}
+	a, err := entry.Action("ep", browserTask, h.CatalogueGeneration)
+	if err != nil {
+		t.Fatal(err)
+	}
+	a.ID = fixtureActionID
+	wire, err := approval.CanonicalOperatorAction(a)
+	var decoded approval.OperatorAction
+	if err == nil {
+		decoded, err = approval.DecodeOperatorAction(wire)
+	}
+	if err != nil || !decoded.Options.Headless || !slices.Equal(decoded.Argv, explicit.Intent.Argv) ||
+		!strings.Contains(decoded.Review, "Headless codex gpt-6.2-sol, cap 1m0s") || !strings.Contains(string(wire), `"headless":true`) {
+		t.Fatalf("signed review mode: %v %v", err, decoded.Options)
+	}
+	digest, err := approval.OperatorSignatureDigest(decoded)
+	if err != nil {
+		t.Fatal(err)
+	}
+	challenge := [32]byte{1}
+	msg, err := approval.OperatorActionMessage(digest, challenge, 1600)
+	if err != nil {
+		t.Fatal(err)
+	}
+	key := approvaltest.NewSigner(89)
+	proof, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: decoded, Challenge: challenge, Expiry: 1600, Now: 1010,
+		Artifact: key.Sign(approval.OperatorActionNamespace, msg), Format: approval.FormatSSH,
+		Approver: approval.Approver{Name: "fixture", Format: approval.FormatSSH, Key: key.PublicKey(), EnrolledAt: 900}})
+	if err != nil || proof.Binds(decoded, challenge, 1600) != nil {
+		t.Fatalf("the mode was not bound by fixture proof: %v", err)
+	}
+
+	ctx, cancel := r.env.operatorContext(activationTestBudget)
+	defer cancel()
+	if _, err := reconstructOperatorStart(ctx, decoded, store.OperatorMintAttempt{Catalogue: h.CatalogueGeneration},
+		r.src.presets, r.src.profiles, r.src.workspace, r.src.hosts, r.src.clock); err != nil {
+		t.Fatalf("same signed headless start did not re-resolve: %v", err)
+	}
+	t.Log("mode=interactive; explicit=headless; cap=1m0s; headless is in canonical signed bytes")
+}
+
+func TestOperatorStartWalkthroughInChrome(t *testing.T) {
+	t.Setenv("HOME", t.TempDir())
+	b := startPasskeyBrowser(t)
+	r, c, h, runner := operatorActivationFixtureWithTask(t, browserFakeChild, browserTask)
+	entry, _ := browserPreset(t, r, c)
+	h = publishBrowserCatalogue(t, r, c, h, entry)
+	writeCtx, cancel := r.env.operatorContext(activationTestBudget)
+	_, err := r.db.ExecContext(writeCtx, "DELETE FROM operator_action_steps; DELETE FROM operator_actions")
+	cancel()
+	if err != nil {
+		t.Fatal(err)
+	}
+	log := &syncBuffer{}
+	d, err := daemon.New(activationTestContext(t), daemon.Config{DBPath: c.World, CredentialsPath: c.Credentials,
+		BindHost: daemon.DefaultBindHost, ErrorLog: log, ApprovalClock: r.src.clock})
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { d.Close() })
+	srv := httptest.NewServer(d.Handler())
+	t.Cleanup(func() { srv.CloseClientConnections(); srv.Close() })
+	origin := strings.Replace(srv.URL, "127.0.0.1", "localhost", 1)
+	tab := map[string]any{"authenticator": "uv", "inject": "getRecorder", "viewport": map[string]any{"width": 390, "height": 844}}
+	b.do(t, nil, "tab", tab)
+	virtualPasskey(t, b, c.Credentials, origin, "fixture-sol", 93, r.src.clock()-10)
+	wait := func(js string) { t.Helper(); b.do(t, nil, "waitFor", map[string]any{"expression": js, "ms": 10000}) }
+	open := func(path string) {
+		t.Helper()
+		b.do(t, nil, "openCard", map[string]any{"url": origin + path, "raw": true})
+	}
+	tap := func(id string) { t.Helper(); b.do(t, nil, "tap", map[string]any{"id": id}) }
+	fill := func(id, text string) {
+		t.Helper()
+		wait(fmt.Sprintf(`(() => { const e=document.getElementById(%q); e.value=%q; e.dispatchEvent(new Event('input')); return true; })()`, id, text))
+	}
+	key := func(k string) { t.Helper(); b.do(t, nil, "key", map[string]any{"key": k}) }
+	enter := func(line string) {
+		t.Helper()
+		wait(fmt.Sprintf(`(() => { const e=document.getElementById('addr'); e.focus(); e.value=%q; return document.activeElement===e; })()`, line))
+		key("Enter")
+	}
+	counts := func(want [4]int) {
+		t.Helper()
+		var got [4]int
+		err := r.db.QueryRowContext(activationTestContext(t), `SELECT (SELECT count(*) FROM operator_actions),
+  (SELECT count(*) FROM episode_starts),(SELECT count(*) FROM session_credentials),
+  (SELECT sign_count FROM webauthn_credentials WHERE name='fixture-sol')`).Scan(&got[0], &got[1], &got[2], &got[3])
+		if err != nil || got != want {
+			t.Fatalf("actions/starts/sessions/counter: %v want %v (%v)", got, want, err)
+		}
+	}
+	open("/surface")
+	wait(`Boolean(window.World)`)
+	noEffects := activationEvidence(t, r)
+	enter(strings.TrimSuffix(browserStart, " --headless"))
+	wait(`Boolean(document.querySelector('#verb .verb-body')?.textContent.includes('terminal') || document.getElementById('hmsg').textContent)`)
+	if !reflect.DeepEqual(activationEvidence(t, r), noEffects) {
+		t.Fatal("the preset's implicit interactive mode inserted an action")
+	}
+	doc, err := os.ReadFile(filepath.Join("..", "..", "docs", "QUICKSTART.md"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	_, example, present := strings.Cut(string(doc), "<!-- operator-browser-fixture -->\n")
+	line, rest, lineOK := strings.Cut(strings.TrimPrefix(example, "```text\n"), "\n```")
+	_, rest, taskStart := strings.Cut(rest, "```text\n")
+	task, _, taskEnd := strings.Cut(rest, "\n```")
+	if !present || !lineOK || !taskStart || !taskEnd || line != browserStart || task != browserTask {
+		t.Fatal("QUICKSTART fixture did not match the actual preset and literal task")
+	}
+	enter(line)
+	wait(`document.activeElement?.id==='operator-task' && document.querySelector('#verb .verb-body')?.textContent.includes('Headless codex gpt-6.2-sol')`)
+	counts([4]int{0, 0, 0, 1})
+	key("n")
+	wait(`document.activeElement?.id==='operator-task' && document.getElementById('operator-task').value==='n'`)
+	fill("operator-task", "")
+	if !reflect.DeepEqual(activationEvidence(t, r), noEffects) {
+		t.Fatal("Enter changed the start evidence")
+	}
+	wait(`innerWidth===390 && innerHeight===844 && document.documentElement.scrollWidth===390 && document.getElementById('operator-task').getBoundingClientRect().right<=390`)
+	key("Tab")
+	wait(`document.activeElement?.id==='operator-request'`)
+	key("Enter")
+	wait(`document.getElementById('operator-request-info')?.textContent.includes('1–8192 UTF-8 bytes')`)
+	counts([4]int{0, 0, 0, 1})
+	fill("operator-task", task)
+	wait(`(() => {
+  const original=window.fetch.bind(window);
+  window.fetch=async (path,options) => {
+    const response=await original(path,options);
+    if(path==='/surface/operator-actions' && options?.method==='POST') {
+      window.__fixtureTaskFrozen = document.getElementById('operator-task').readOnly;
+      throw Error('fixture dropped POST response');
+    }
+    return response;
+  };
+  return true;
+ })()`)
+	tap("operator-request")
+	wait(`document.getElementById('operator-request-info')?.textContent.includes('Request outcome unavailable') && window.__fixtureTaskFrozen===true`)
+	counts([4]int{1, 0, 0, 1})
+	open("/surface")
+	wait(`document.activeElement?.id==='operator-request' && document.getElementById('operator-request').textContent==='Reconnect' && document.getElementById('operator-task').readOnly`)
+	tap("operator-request")
+	wait(`document.activeElement?.id==='operator-title' && document.getElementById('operator-key')?.disabled===false && !document.getElementById('operator-controls').hidden`)
+	counts([4]int{1, 0, 0, 1})
+	ctx := activationTestContext(t)
+	ro, err := store.OpenCredentialsReadOnly(c.Credentials)
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { ro.Close() })
+	var id string
+	var operations int
+	if err := r.db.QueryRowContext(ctx, "SELECT id,(SELECT count(*) FROM operator_action_steps WHERE kind='operation') FROM operator_actions").Scan(&id, &operations); err != nil || operations != 1 {
+		t.Fatalf("reconnect operation count=%d: %v", operations, err)
+	}
+	row, err := ro.OperatorAction(ctx, id, r.src.clock())
+	if err != nil {
+		t.Fatal(err)
+	}
+	r.action, err = approval.DecodeOperatorAction(row.Action)
+	if err != nil || r.action.Task != browserTask || !r.action.Options.Headless {
+		t.Fatalf("literal signed task or mode: %v", err)
+	}
+	wait(fmt.Sprintf(`document.getElementById('operator-task').textContent===%q && !document.getElementById('operator-task').querySelector('script') && document.documentElement.scrollWidth===390`, browserTask))
+	key("Tab")
+	wait(`document.activeElement?.id==='operator-key'`)
+	wait(`document.querySelector('#operator-controls details').open = true`)
+	fill("operator-signature", "y")
+	before, err := os.ReadFile(c.Credentials)
+	if err != nil {
+		t.Fatal(err)
+	}
+	tap("operator-sign")
+	wait(`document.getElementById('operator-answer')?.textContent.includes('NotApproved') && !document.getElementById('operator-key').disabled`)
+	counts([4]int{1, 0, 0, 1})
+	after, err := os.ReadFile(c.Credentials)
+	if err != nil || !bytes.Equal(before, after) {
+		t.Fatal("plain y changed the pending row or key counter")
+	}
+	tap("operator-key")
+	wait(`document.getElementById('operator-answer')?.textContent==='Confirmed by webauthn:fixture-sol'`)
+	counts([4]int{1, 0, 0, 2})
+	var web struct {
+		Gets       []cardGet
+		Violations []string
+	}
+	b.do(t, &web, "waitFor", map[string]any{"expression": `({gets:window.__gets||[],violations:window.__cspv||[]})`})
+	challenge, err := approval.OperatorActionMessageHash([32]byte(row.Digest), [32]byte(row.Challenge), row.ExpiresAt)
+	if err != nil || len(web.Gets) != 1 || web.Gets[0].UserVerification != "required" ||
+		web.Gets[0].Challenge != base64.RawURLEncoding.EncodeToString(challenge[:]) || len(web.Violations) != 0 {
+		t.Fatalf("browser ceremony: %+v", web)
+	}
+	b.do(t, nil, "tab", map[string]any{"viewport": tab["viewport"]}) // close the approved page before running a single fake launch
+
+	f := &operatorLaunchRig{worktreeStartRig: r, c: c, h: h, runner: runner, ctx: ctx, db: ro, now: r.src.clock()}
+	if out := f.tick(t); out != "" {
+		t.Fatalf("prepare: %s", out)
+	}
+	f.readPrepared(t)
+	health := r.env.client.Transport
+	r.env.client.Transport = serviceTransport(func(req *http.Request) (*http.Response, error) {
+		if req.URL.Path != "/v1/session/self" {
+			return health.RoundTrip(req)
+		}
+		if req.URL.Host != c.Address || req.Header.Get("Authorization") != "Bearer "+f.token {
+			t.Error("session read escaped its fixture binding")
+		}
+		data, _ := json.Marshal(f.self)
+		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(data))}, nil
+	})
+	// Only the fixture launcher is entered. Keep the fake child alive until the stored birth is read.
+	r.env.operatorSupervisor = func(job operatorLaunchJob, out io.Writer) error { return r.env.executeOperatorLaunch(ctx, job, out) }
+	r.env.operatorAgentRun = func(ctx context.Context, argv []string, dir string, env []string, out, errs *os.File, wall time.Duration, hooks agentrun.LaunchHooks) (agentrun.Result, error) {
+		prompts := 0
+		for _, arg := range argv {
+			if strings.Contains(arg, browserTask) {
+				prompts++
+			}
+		}
+		if !slices.Contains(argv, "read-only") || !strings.HasPrefix(argv[len(argv)-1], browserTask) || prompts != 1 ||
+			!slices.Contains(env, "WORLD_SESSION="+f.token) || wall != time.Minute {
+			t.Error("fake launch did not preserve confinement or the literal prompt")
+		}
+		start := hooks.Started
+		hooks.Started = func(pid int) error {
+			defer func() {
+				if err := os.WriteFile(filepath.Join(dir, "release"), nil, 0600); err != nil {
+					t.Error(err)
+				}
+			}()
+			if err := start(pid); err != nil {
+				return err
+			}
+			if got := f.receipt(t).ProcessStatus(agentrun.ProcessIdentity); got != "launched" {
+				t.Errorf("live birth readback: %s", got)
+			}
+			return nil
+		}
+		return agentrun.RunObserved(ctx, argv, dir, env, out, errs, wall, hooks)
+	}
+	if out := f.tick(t); !strings.Contains(out, ": launched") || !strings.Contains(out, ": stopped") {
+		t.Fatalf("fake launch: %s", out)
+	}
+	p := "/surface/operator-actions/" + id
+	open(p)
+	wait(`document.getElementById('operator-state')?.textContent==='stopped' && document.getElementById('operator-report')?.getAttribute('href')==='/v1/episodes/ep/report' &&
+   document.activeElement?.id==='operator-title' && !document.getElementById('operator-refresh').disabled && document.documentElement.scrollWidth===390`)
+	tap("operator-refresh")
+	wait(`document.getElementById('operator-state')?.textContent==='stopped' && !document.getElementById('operator-refresh').disabled`)
+	counts([4]int{1, 1, 1, 2})
+	worktrees := presetGit(t, r.f.r.f, "worktree", "list", "--porcelain")
+	physical, err := filepath.EvalSymlinks(r.path())
+	ran, readErr := os.ReadFile(filepath.Join(r.bin, "ran"))
+	matches := 0
+	for _, line := range strings.Split(worktrees, "\n") {
+		path, ok := strings.CutPrefix(line, "worktree ")
+		if !ok {
+			continue
+		}
+		linked, e := filepath.EvalSymlinks(path)
+		if e != nil {
+			t.Fatal(e)
+		}
+		if linked == physical {
+			matches++
+		}
+	}
+	if err != nil || readErr != nil || string(ran) != "x" || strings.Count(worktrees, "worktree ") != 2 ||
+		matches != 1 || f.receipt(t).Receipt.Run != id {
+		t.Fatalf("durable worktree or launch: %v %v %q (%d matching registrations)", err, readErr, ran, matches)
+	}
+	for _, endpoint := range []string{"/surface/requests.json", p + "/status.json"} {
+		var text string
+		b.do(t, &text, "waitFor", map[string]any{"expression": fmt.Sprintf("fetch(%q).then(r=>r.text())", endpoint)})
+		if strings.Contains(text, "alert(1)") || strings.Contains(text, "authenticatorData") || strings.Contains(text, "messageHex") {
+			t.Fatal("routine read revealed task or proof")
+		}
+	}
+	if strings.Contains(log.String(), "alert(1)") || strings.Contains(log.String(), "authenticatorData") {
+		t.Fatal("daemon log revealed task or proof")
+	}
+	tap("operator-back")
+	wait(`location.pathname==='/surface' && !location.search && !location.hash && Boolean(document.getElementById('addr'))`)
+	b.do(t, nil, "tab", map[string]any{})
+}
diff --git a/host/daemon/daemon.go b/host/daemon/daemon.go
index 1ceff23..5f529d1 100644
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -335,6 +335,9 @@ type Config struct {
 	// on that stream were measured deadlocking Run against an io.Pipe.
 	ErrorLog io.Writer
 
+	// ApprovalClock supplies Unix seconds to the approval routes. Nil uses the system clock.
+	ApprovalClock func() int64
+
 	// templateRebuildGate is a test seam (row 153 AC3.7): when non-nil, the
 	// background template rebuild awaits a receive on it before each rebuild.
 	// Production leaves it nil.
@@ -831,6 +834,10 @@ func newDaemon(ctx, phaseParent context.Context, registryBudget time.Duration, c
 		_ = s.Close()
 		return nil, &StartupError{Stage: StageConfig, Detail: err.Error()}
 	}
+	approveNow := cfg.ApprovalClock
+	if approveNow == nil {
+		approveNow = func() int64 { return time.Now().Unix() }
+	}
 	d := &Daemon{
 		follow:    follow,
 		startedAt: time.Now().Unix(),
@@ -838,7 +845,7 @@ func newDaemon(ctx, phaseParent context.Context, registryBudget time.Duration, c
 		aguiBudget: aguiRunBudget, aguiBodyBound: aguiBodyBound, aguiTick: aguiTick,
 		aguiStop: make(chan struct{}), aguiSlots: make(chan struct{}, aguiCap),
 		readDeadline: readDeadline, errLog: resolveErrorLog(cfg.ErrorLog), finishNow: func() int64 { return time.Now().Unix() },
-		approveNow:   func() int64 { return time.Now().Unix() },
+		approveNow:   approveNow,
 		scanPageSize: integrityScanPageSize, scanRowBudget: integrityScanRowBudget,
 		scanTimeBudget: integrityScanTimeBudget, workspace: workspace,
 		resolver: authority.New(authority.ReadOnlyAt(credPath)), credentialsPath: credPath,
diff --git a/host/surfacepage/app/app.js b/host/surfacepage/app/app.js
index fdba53d..f7aa3bb 100644
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -1427,15 +1427,84 @@
 
   // ---------- The bar's verbs (row 243 V1): GET /surface/verb.json decides; only a schema-valid address navigates ----------
   const VERB_ROUTE_BUDGET_MS = 15000;   // above the daemon's 10 s read deadline; a named constant, not a setting
-  const VERB_NEEDS = { address: [], result: ['cli'], navigate: ['href'], terminal: ['cli', 'reason'], refused: ['reason'] };
-  const UNREADABLE = { say: 'the verb route answered something this page cannot read' };
+  const VERB_NEEDS = { address: [], result: ['cli'], navigate: ['href'], terminal: ['cli', 'reason'], refused: ['reason'], review: ['text'] };
+  const OPERATOR_DRAFT = 'world.operator-request.v1';
+  const OPERATOR_PATH = '/surface/operator-actions';
+  const OPERATOR_ID = /^oa-[0-9a-f]{32}(?![\s\S])/;
+  const EPISODE_ID = /^[a-z0-9][a-z0-9-]{0,63}(?![\s\S])/;
+  const PRESET_ID = /^[a-z0-9][a-z0-9-]{0,50}(?![\s\S])/;
+  const OPERATION_ID = /^[0-9a-f]{32}(?![\s\S])/;
+  const taskValid = task => typeof task === 'string' && task.length > 0 && !task.includes('\0') &&
+    new TextEncoder().encode(task).length <= 8192;
+  function reviewValid(argv, catalogueGeneration) {
+    return Array.isArray(argv) && argv.length === 8 &&
+      argv[0] === 'ailang-worldd' && argv[1] === 'episode' && argv[2] === 'start' &&
+      typeof argv[3] === 'string' && EPISODE_ID.test(argv[3]) && argv[4] === '--preset' &&
+      typeof argv[5] === 'string' && PRESET_ID.test(argv[5]) && argv[5] !== 'se-tools' &&
+      argv[6] === '--launch' && argv[7] === '--headless' &&
+      Number.isSafeInteger(catalogueGeneration) && catalogueGeneration > 0;
+  }
+  function operatorDraft(storage) {
+    try {
+      const d = JSON.parse(storage.getItem(OPERATOR_DRAFT));
+      const r = d?.request;
+      if (typeof d?.body === 'string' && d.body && r && reviewValid(r.argv, r.catalogueGeneration) &&
+          typeof r.operationId === 'string' && OPERATION_ID.test(r.operationId) && taskValid(r.task)) return d;
+    } catch (_) { /* page stays readable */ }
+    return null;
+  }
+  function operatorCardHref(b, episode) {
+    if (!b || b.v !== 1 || b.kind !== 'start' || b.episode !== episode || typeof b.id !== 'string' || !OPERATOR_ID.test(b.id) ||
+        !['requested', 'confirmed', 'claimed', 'running', 'succeeded', 'failed', 'expired', 'cancelled', 'uncertain'].includes(b.state) ||
+        ![b.now, b.createdAt, b.expiresAt].every(Number.isSafeInteger)) return '';
+    const path = OPERATOR_PATH + '/' + b.id;
+    return b.href === path && b.statusHref === path + '/status.json' ? path : '';
+  }
+  function operatorRequestControl(panel, io) {
+    let draft = operatorDraft(io.storage), busy = false;
+    let request = draft?.request;
+    return { get request() { return request; }, async submit(task) {
+      if (busy || !io.current()) return {};
+      if (!request) {
+        if (!taskValid(task)) return { error: 'Enter a task of 1–8192 UTF-8 bytes without NUL.' };
+        const operationId = Array.from(io.random(new Uint8Array(16)), n => n.toString(16).padStart(2, '0')).join('');
+        request = { argv: panel.argv, operationId, catalogueGeneration: panel.catalogueGeneration, task };
+      }
+      try {
+        io.storage.setItem(OPERATOR_DRAFT, JSON.stringify({ body: draft?.body ?? panel.body, request }));
+      } catch (_) {
+        return { error: 'Request could not be saved for reconnect. Restore browser storage and retry.' };
+      }
+      busy = true;
+      try {
+        const res = await io.fetch(OPERATOR_PATH, {method: 'POST', headers: {'Content-Type': 'application/json'},
+          redirect: 'error', credentials: 'omit', signal: io.timeout(VERB_ROUTE_BUDGET_MS), body: JSON.stringify(request)});
+        const b = await res.json();
+        if (!io.current()) return {};
+        const href = res.status === 200 && operatorCardHref(b, request.argv[3]);
+        if (!href) {
+          const reason = typeof b?.error === 'string' && typeof b?.reason === 'string' ? b.error + ': ' + b.reason + ' ' : '';
+          return { error: reason + 'Request outcome unavailable. Use Reconnect to read or retry this request.' };
+        }
+        io.storage.removeItem(OPERATOR_DRAFT);
+        return { href };
+      } catch (_) {
+        return { error: 'Request outcome unavailable. Use Reconnect to read or retry this request.' };
+      } finally {
+        busy = false;
+      }
+    } };
+  }
+  const UNREADABLE = { say: 'the verb route answered something this page cannot read', retry: true };
   // verbOutcome maps the route's answer (HTTP status, body text) to one action; anything else is a refusal.
   function verbOutcome(status, text) {
-    if (status !== 200) return { say: 'the verb route did not answer: HTTP ' + status };
+    if (status !== 200) return { say: 'the verb route did not answer: HTTP ' + status, retry: true };
     let b = null;
     try { b = JSON.parse(text); } catch (_) { return UNREADABLE; }
     const need = b && typeof b === 'object' && b.schema === 'world/surface-verb/v1' && b.v === 1 && Object.prototype.hasOwnProperty.call(VERB_NEEDS, b.outcome) ? VERB_NEEDS[b.outcome] : null;
     if (!need || !need.every(k => typeof b[k] === 'string' && b[k] !== '')) return UNREADABLE;
+    if (b.outcome === 'review') return reviewValid(b.argv, b.catalogueGeneration)
+      ? { review: { argv: b.argv, catalogueGeneration: b.catalogueGeneration, body: b.text } } : UNREADABLE;
     if (b.outcome === 'address') return { address: true };
     if (b.outcome === 'refused') return typeof b.cli === 'string' && b.cli !== '' ? { panel: { head: '$ ' + b.cli, body: b.reason, copy: b.cli } } : { say: b.reason };
     if (b.outcome === 'navigate') return b.href.startsWith('/surface/approve/mr-') || b.href.startsWith('/surface/decide/sha256:') || b.href.startsWith('/surface/confirm/vr-') ? { assign: b.href } : { say: 'refused: unexpected target' };
@@ -1462,8 +1531,8 @@
       return job.superseded ? { superseded: true } : verbOutcome(res.status, text);   // an answer that beat the abort is still dropped
     } catch (e) {
       if (job.superseded) return { superseded: true };
-      if (job.late) return { say: 'the verb route did not answer in time' };
-      return { say: 'the verb route is unreachable: ' + (e && e.message ? e.message : String(e)) };
+      if (job.late) return { say: 'the verb route did not answer in time', retry: true };
+      return { say: 'the verb route is unreachable: ' + (e && e.message ? e.message : String(e)), retry: true };
     } finally {
       io.clearTimeout(timer);
       if (verbJob === job) verbJob = null;
@@ -1473,9 +1542,11 @@
   // a navigation or a panel clears the bar; a refusal says why and leaves the view and the input.
   function verbAct(act, raw, fx) {
     if (act.superseded) return;
+    if (act.review) return fx.review(act.review);
     if (act.address) { fx.clear(); return fx.address(raw); }
     if (act.assign) { fx.clear(); return fx.assign(act.assign); }
     if (act.panel) { fx.clear(); return fx.panel(act.panel); }
+    if (act.retry && fx.retry) return fx.retry(act.say);
     fx.say(act.say);
   }
   // goWith is the bar's Enter over fx (design A4): every Enter owns the bar. An empty line supersedes a pending ask and goes
@@ -1486,7 +1557,7 @@
     if (String(raw == null ? '' : raw).trim() === '') { dropVerb(); return fx.out(); }
     const act = await fx.ask(raw);
     if (!current()) return;
-    return verbAct(act, raw, { ...fx.act, address: r => fx.address(r, current) });
+    return verbAct(act, raw, { ...fx.act, address: r => fx.address(r, current), review: p => fx.act.review(p, current) });
   }
   // addressGo is the address bar's own path (M4a's grammar), reached only through the route's address answer; after the
   // index loads it acts only for the current Enter.
@@ -1513,7 +1584,7 @@
 
   if (typeof module !== 'undefined') module.exports = { needsYouHTML, decidedHTML, notifications, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, finishMatchesNode, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
     adaptModel, projectHeadsModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, filePath, fileEvidenceHTML, filesTabHTML, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW },
-    VERB_ROUTE_BUDGET_MS, verbOutcome, askVerb, verbAct, verbPanel, goWith, addressGo };
+    VERB_ROUTE_BUDGET_MS, verbOutcome, askVerb, verbAct, verbPanel, goWith, addressGo, operatorDraft, operatorRequestControl, operatorCardHref };
   if (typeof document === 'undefined') return;
 
   // ================= The page =================
@@ -3426,6 +3497,7 @@
   }
   // go is the bar's Enter: empty is home (no fetch); anything else asks the verb route first.
   function go(raw) {
+    if ($('operator-task')) $('verb').remove();
     return goWith(raw, { out, address: goAddress,
       ask: q => askVerb(q, { fetch: (u, o) => fetch(u, o), AbortController, setTimeout: (f, ms) => setTimeout(f, ms), clearTimeout: t => clearTimeout(t) }),
       act: { assign: h => {
@@ -3434,7 +3506,12 @@
           : camera.active ? (location.hash === '#heads/activity' ? '#heads/activity' : '#heads')
           : location.hash === '#activity' ? '#activity' : '';
         location.assign(h + suffix);
-      }, panel: showVerbPanel, say, clear: () => { if ($('addr').value === raw) $('addr').value = ''; } } });
+      }, review: showOperatorReview, panel: showVerbPanel, say,
+      retry: message => {
+        say(message + '. Press Enter to reconnect.');
+        if ($('addr').value === raw) $('addr').focus({preventScroll: true});
+      },
+      clear: () => { if ($('addr').value === raw) $('addr').value = ''; } } });
   }
   function showVerbPanel(p) {
     const old = $('verb'), el = verbPanel(document, p);
@@ -3443,11 +3520,60 @@
     closeBtn.className = 'verb-close';
     closeBtn.textContent = 'Close';
     closeBtn.setAttribute('aria-label', 'Close command output');
-    closeBtn.addEventListener('click', () => el.remove());
+    closeBtn.addEventListener('click', () => {
+      el.remove();
+      if (p.argv) $('addr').focus({preventScroll: true});
+    });
     el.append(closeBtn);
     if (old) old.replaceWith(el); else $('stage').appendChild(el);
     const b = el.querySelector('button[data-copy]');
     if (b) b.addEventListener('click', () => { try { navigator.clipboard.writeText(b.dataset.copy).then(() => say('Copied command'), () => say('Copy it from the panel')); } catch (_) { say('Copy it from the panel'); } });
+    return el;
+  }
+  let operatorStorage;
+  try { operatorStorage = window.sessionStorage; } catch (_) { /* request stays inert */ }
+  function showOperatorReview(panel, current) {
+    const draft = operatorDraft(operatorStorage);
+    if (draft) panel = { ...draft.request, body: draft.body };
+    const el = showVerbPanel({ ...panel, head: 'Review ' + panel.argv[3] });
+    const body = el.querySelector('.verb-body'), area = document.createElement('div');
+    area.className = body.className;
+    area.textContent = body.textContent;
+    body.replaceWith(area);
+    const label = document.createElement('label'), task = document.createElement('textarea');
+    const request = document.createElement('button'), info = document.createElement('p');
+    task.id = 'operator-task'; task.rows = 3; task.cols = 26;
+    label.textContent = 'Task (1–8192 UTF-8 bytes)'; label.htmlFor = task.id;
+    request.id = 'operator-request'; request.type = 'button';
+    info.id = 'operator-request-info'; info.setAttribute('role', 'status');
+    area.append(label, document.createElement('br'), task, document.createElement('br'), request, info);
+    const live = () => current() && el.isConnected;
+    const control = operatorRequestControl(panel, {storage: operatorStorage, current: live, fetch: (...args) => fetch(...args),
+      random: b => crypto.getRandomValues(b), timeout: ms => AbortSignal.timeout(ms)});
+    const repaint = () => {
+      task.readOnly = !!control.request;
+      request.textContent = control.request ? 'Reconnect' : 'Request action';
+    };
+    task.value = control.request?.task || '';
+    repaint();
+    request.addEventListener('click', async () => {
+      request.disabled = true;
+      task.readOnly = true;
+      await control.submit(task.value).then(result => {
+        if (!live()) return;
+        if (result.href) return location.assign(result.href);
+        info.textContent = result.error || '';
+      });
+      repaint();
+      request.disabled = false;
+      if (live()) (task.readOnly ? request : task).focus();
+    });
+    (task.readOnly ? request : task).focus();
+  }
+  const outstandingOperatorRequest = operatorDraft(operatorStorage);
+  if (outstandingOperatorRequest) {
+    const turn = barTurn;
+    showOperatorReview(outstandingOperatorRequest.request, () => turn === barTurn);
   }
   function goAddress(raw, current) { return addressGo(raw, current, { parse: SEAMS.parseAddress, loadIndex, out, openState, say }); }
 
@@ -3804,7 +3930,7 @@
   });
   const onRow = row => Object.keys(DISP).filter(id => DISP[id].row === row).sort((a, b) => DISP[a].x - DISP[b].x);
   document.addEventListener('keydown', e => {
-    if (e.target && (e.target.tagName === 'SELECT' || e.target.tagName === 'INPUT' || e.target.closest?.('dialog, #compose-btn, .state-branch'))) return;
+    if (e.target && (e.target.tagName === 'SELECT' || e.target.tagName === 'INPUT' || e.target.closest?.('dialog, #compose-btn, .state-branch, #verb:has(#operator-task)'))) return;
     const k = e.key; let done = true;
     if (e.target.closest?.('.camera-angles') && k === 'Enter') return;
     if (e.target.closest?.('.stack-disclosure summary, .file-path, .file-evidence button') && k === 'Enter') return;
diff --git a/host/surfacepage/app_operator_test.go b/host/surfacepage/app_operator_test.go
new file mode 100644
index 0000000..183b1e3
--- /dev/null
+++ b/host/surfacepage/app_operator_test.go
@@ -0,0 +1,188 @@
+package surfacepage
+
+import "testing"
+
+func runOperatorBar(t *testing.T, script string) {
+	t.Helper()
+	runOperatorScript(t, AppScript, `
+const assert = require('node:assert/strict'), app = require(process.argv[1]);
+const argv = ['ailang-worldd','episode','start','s8','--preset','p','--launch','--headless'];
+const b = {schema:'world/surface-verb/v1',v:1,outcome:'review',argv,catalogueGeneration:7,text:'Headless codex gpt-6.2-sol'};
+const review = (change={}) => app.verbOutcome(200,JSON.stringify({...b,...change}));
+const id = 'oa-'+'a'.repeat(32), href = '/surface/operator-actions/'+id;
+const result = {v:1, id, href, statusHref:href+'/status.json', kind:'start', episode:'s8', state:'requested', now:1, createdAt:1, expiresAt:601};
+const reply = body => ({status:200,json:async()=>body});
+let slot = null, posts = [], randoms = 0, current = true, response = reply(result);
+const storage = { getItem:()=>slot, setItem:(key,value)=>{slot=value;}, removeItem:()=>{slot=null;} };
+const io = {storage, random: x=>(randoms++, x.fill(randoms)), current:()=>current, timeout:ms=>(assert.equal(ms,15000), 'timeout'),
+  fetch: async (path,opts)=>{ assert.equal(path,'/surface/operator-actions'); assert.equal(opts.signal,'timeout');
+    assert.equal(opts.method,'POST'); assert.equal(opts.headers['Content-Type'],'application/json');
+    assert.equal(opts.redirect,'error'); assert.equal(opts.credentials,'omit'); posts.push(JSON.parse(opts.body));
+    if(!response) throw Error('offline'); return response; }};
+const control = () => app.operatorRequestControl(review().review,io);
+const deferred = () => { let resolve; const promise = new Promise(r=>resolve=r); return {promise,resolve}; };
+let completed = false;
+process.once('beforeExit', () => {
+  if (!completed) { console.error('operator bar assertion did not settle'); process.exitCode=1; }
+});
+(async()=>{`+script+`})().then(() => { completed=true; }, e => { completed=true; console.error(e); process.exitCode=1; });
+`)
+}
+
+func TestOperatorBarStartReview(t *testing.T) {
+	runOperatorBar(t, `
+const actions = [];
+const fx = {ask:async()=>review(),act:{review:(p,live)=>actions.push([p,live()]),
+  panel:()=>assert.fail('ordinary panel'),assign:()=>assert.fail('navigated'),clear:()=>assert.fail('cleared input'),say:assert.fail}};
+await app.goWith(argv.join(' '),fx);
+assert.equal(posts.length,0);
+assert.equal(randoms,0);
+assert.deepEqual(actions,[[{argv,catalogueGeneration:7,body:b.text},true]]);
+assert.equal((await control().submit('')).error.includes('1–8192'),true);
+assert.equal((await control().submit('é'.repeat(4097))).error.includes('1–8192'),true);
+assert.equal((await control().submit('a\0b')).error.includes('without NUL'),true);
+assert.equal(posts.length,0);
+assert.equal((await control().submit('é'.repeat(4096))).href,href);
+`)
+}
+
+func TestOperatorOtherVerbParity(t *testing.T) {
+	runOperatorBar(t, `
+const v = x=>app.verbOutcome(200,JSON.stringify({...b,...x}));
+assert.deepEqual(v({outcome:'address'}),{address:true});
+for (const cli of ['episode report e','why sha256:x']) {
+  assert.deepEqual(v({outcome:'result',cli,text:'body'}),{panel:{head:'$ '+cli,body:'body',copy:''}});
+}
+const approve = '/surface/approve/mr-0123456789abcdef';
+assert.deepEqual(v({outcome:'navigate',href:approve}),{assign:approve});
+for(const cli of ['world-publish','approver enroll','approver revoke','session mint','session revoke','service']) {
+  assert.deepEqual(v({outcome:'refused',cli,reason:'blocked'}),{panel:{head:'$ '+cli,body:'blocked',copy:cli}});
+}
+`)
+}
+
+func TestOperatorLiteralTask(t *testing.T) {
+	runOperatorBar(t, `
+const task = '# Task\n<script>alert(1)</script> ; $(literal) &';
+assert.equal((await control().submit(task)).href,href);
+assert.deepEqual(posts,[{argv,operationId:'01'.repeat(16),catalogueGeneration:7,task}]);
+`)
+}
+
+func TestOperatorDuplicateTap(t *testing.T) {
+	runOperatorBar(t, `
+const pending = deferred(), r = control();
+response = pending.promise;
+const first = r.submit('literal');
+const second = r.submit('different');
+assert.equal(posts.length,1);
+assert.deepEqual(await second,{});
+pending.resolve(reply(result));
+assert.deepEqual(await first,{href});
+assert.equal(randoms,1);
+assert.equal(slot,null);
+`)
+}
+
+func TestOperatorReconnect(t *testing.T) {
+	runOperatorBar(t, `
+response = null;
+assert.match((await control().submit('frozen')).error,/Reconnect/);
+const draft = app.operatorDraft(storage);
+assert.equal(draft.request.task,'frozen');
+assert.equal(draft.body,b.text);
+response = reply(result);
+assert.deepEqual(await control().submit('edited'),{href});
+assert.deepEqual(posts,[draft.request,draft.request]);
+assert.equal(randoms,1);
+assert.equal(slot,null);
+`)
+}
+
+func TestOperatorPanelLatest(t *testing.T) {
+	runOperatorBar(t, `
+const old = deferred(), now = deferred(), acted = [];
+const fx = {ask:q=>q==='older'?old.promise:now.promise,act:{review:p=>acted.push(p.body),
+  assign:h=>acted.push(h),clear:()=>{},say:assert.fail}};
+const a = app.goWith('older',fx), n = app.goWith('newer',fx);
+now.resolve(review()); await n;
+old.resolve({assign:href}); await a;
+assert.deepEqual(acted,[b.text]);
+current = false;
+assert.deepEqual(await control().submit('stale'),{});
+assert.equal(posts.length,0);
+current = true;
+const pending = deferred(); response = pending.promise;
+const post = control().submit('frozen'); current = false;
+pending.resolve(reply(result));
+assert.deepEqual(await post,{});
+assert.notEqual(slot,null);
+`)
+}
+
+func TestOperatorPanelNetwork(t *testing.T) {
+	runOperatorBar(t, `
+const calls = [];
+let trigger;
+const fx = {ask: raw=>app.askVerb(raw, {AbortController,setTimeout:f=>(trigger=f,1),clearTimeout:()=>{},
+  fetch:()=>new Promise(()=>{})}), act:{retry:x=>calls.push(x),say:()=>assert.fail('no reconnect offered'),clear:()=>assert.fail('cleared'),assign:()=>assert.fail('navigated')}};
+const ask = app.goWith('start',fx); trigger(); await ask;
+assert.match(calls[0],/in time/);
+await app.goWith('start', {...fx,ask:async()=>review({argv:[]})});
+assert.match(calls[1],/cannot read/);
+await app.goWith('start', {...fx,ask:async()=>app.verbOutcome(503,'')});
+assert.match(calls[2],/HTTP 503/);
+await app.goWith('start', {...fx,ask:()=>app.askVerb('start',{AbortController,
+ setTimeout:()=>1,clearTimeout:()=>{},fetch:async()=>{throw Error('offline')}})});
+assert.match(calls[3],/unreachable/);
+for (const change of [{argv:{}},{argv:{...argv,length:8}},{argv:[...argv,'tail']},{argv:['ailang-worldd', 'episode','start',{},'--preset','p','--launch','--headless']},
+  {argv:argv.slice(1)},{argv:argv.with(7,'--launch')},{argv:argv.with(0,'sh')},{argv:argv.with(1,'service')},
+  {argv:argv.with(2,'finish')},{argv:argv.with(3,'bad\n')},{argv:argv.with(4,'--other')},{argv:argv.with(5,3)},
+  {argv:argv.with(5,'se-tools')},{argv:argv.with(5,'A')},{argv:argv.with(5,'a'.repeat(52))},{argv:argv.with(6,'--other')},
+  {catalogueGeneration:0},{catalogueGeneration:-1},{catalogueGeneration:0.5},{catalogueGeneration:'1'},{text:1},{text:''}])
+  assert.match(review(change).say,/cannot read/);
+for (const d of [null, {}, {body:1,request:{argv,catalogueGeneration:1,operationId:'01'.repeat(16),task:'x'}},
+  {body:'',request:{argv,catalogueGeneration:1,operationId:'01'.repeat(16),task:'x'}}, {body:'review',request:null},
+  {body:'review',request:{argv,catalogueGeneration:1,operationId:1,task:'x'}},
+  {body:'review',request:{argv,catalogueGeneration:1,operationId:['01'.repeat(16)],task:'x'}},
+  {body:'review',request:{argv,catalogueGeneration:1,operationId:'x'.repeat(32),task:'x'}},
+  {body:'review',request:{argv,catalogueGeneration:1,operationId:'01'.repeat(16),task:''}}]) {
+  slot=JSON.stringify(d); assert.equal(app.operatorDraft(storage),null);
+}
+slot=null;
+const save=storage.setItem;
+storage.setItem=()=>{ throw Error('disabled storage'); };
+const r = control();
+assert.match((await r.submit('original')).error,/could not be saved/);
+assert.equal(posts.length,0);
+storage.setItem=save;
+response = null;
+assert.match((await r.submit('original')).error,/Reconnect/);
+response = reply({...result,state:['requested']});
+assert.match((await r.submit('will not replace')).error,/Reconnect/);
+response = {status:409,json:async()=>result};
+assert.match((await r.submit('again')).error,/Reconnect/);
+response = {status:409,json:async()=>({error:'Conflict',reason:'changed'})};
+assert.match((await r.submit('again')).error,/Conflict: changed.*Reconnect/);
+assert.equal(randoms,1);
+assert(posts.every(x=>JSON.stringify(x)===JSON.stringify(posts[0])));
+`)
+}
+
+func TestOperatorLinkAllowlist(t *testing.T) {
+	runOperatorBar(t, `
+const pick = b => app.operatorCardHref(b,'s8');
+assert.equal(pick(result),href);
+for (const bad of [href+'x',href+'a',href.slice(0,-1),'x'+href,href+'\n',href.toUpperCase(),href+'#f',href+'?q',href+'/status.json',
+  'https://localhost'+href,'http://localhost:8'+href,'http://user@localhost'+href,'//localhost'+href,{toString:()=>href}]) {
+  assert.equal(pick({...result,href:bad}),'',JSON.stringify(bad));
+  assert.equal(pick({...result,statusHref:bad===result.statusHref ? result.href : bad}),'',JSON.stringify(bad));
+}
+for (const bad of [id+'x','x'+id,id+'\n','x\n'+id,id+'\nx',id.toUpperCase(),id+'a',id.slice(0,-1),{toString:()=>id}]) {
+  const h = '/surface/operator-actions/'+bad;
+  assert.equal(pick({...result,id:bad,href:h,statusHref:h+'/status.json'}),'');
+}
+for (const [k,v] of Object.entries({v:'1',episode:'other',kind:'finish',state:'invented',now:'1',createdAt:[],expiresAt:null}))
+  assert.equal(pick({...result,[k]:v}),'',k);
+`)
+}
diff --git a/host/surfacepage/app_test.go b/host/surfacepage/app_test.go
index 322ce59..f9d3f0c 100644
--- a/host/surfacepage/app_test.go
+++ b/host/surfacepage/app_test.go
@@ -167,7 +167,8 @@ const stacks = JSON.parse(fs.readFileSync('stacks.json', 'utf8'));
 const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
 assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'applyDelta', 'barHTML', 'blendXS', 'cleanTitle', 'crumbsHTML', 'decidedHTML', 'esc', 'fetchCoherent', 'finishMatchesNode', 'filePath', 'fileEvidenceHTML', 'filesTabHTML', 'frontHTML',
   'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'notifications', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'projectHeadsModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict',
-  'VERB_ROUTE_BUDGET_MS', 'askVerb', 'verbAct', 'verbOutcome', 'verbPanel', 'goWith', 'addressGo'].sort());
+  'VERB_ROUTE_BUDGET_MS', 'askVerb', 'verbAct', 'verbOutcome', 'verbPanel', 'goWith', 'addressGo', 'operatorCardHref', 'operatorDraft', 'operatorRequestControl'].sort());
+assert.equal(app.operatorCardHref(null, 'ep'), '');
 
 assert.equal(app.serviceLineHTML(undefined), '');
 // Row 201 A3: Needs you is the pending requests, each a link to its card; a refused row shows why, with no link.
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index ef7c1a8..dde5280 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -2159,3 +2159,28 @@ WORLD_EXEC_NODE=/opt/homebrew/bin/node \
 WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang" \
 go test -count=1 ./host/daemon/ -run '^TestOperatorConfirmHTTP$'
 ```
+
+### Start one headless episode from the browser
+
+This walkthrough runs on the temporary fixture in `TestOperatorStartWalkthroughInChrome` with a fake
+Sol preset called `p`, a fake launcher and a virtual passkey. Open the fixture's World surface at
+`/surface`. Type the following in the address bar and press Enter.
+
+<!-- operator-browser-fixture -->
+```text
+episode start ep --preset p --launch --headless
+```
+
+Check the **Review ep** panel for the preset, repository, base, worktree, headless model, time cap,
+profile, checks and grants. The task field receives focus. Paste the fixture task into **Task**.
+
+```text
+# Task
+<script>alert(1)</script> ; $(literal) &
+```
+
+Choose **Request action** to open its durable card. Check the signed review and literal task before
+choosing **Confirm with passkey**. The card shows the recorded confirmation and the last read status.
+Close the page and reopen the same card URL to read the outcome and **Episode report**. Use
+**Refresh status** to reconnect if status reads fail. **Back to World** returns to the World surface.
+If the request response was lost, **Reconnect** in the review panel retries the same frozen request.
````


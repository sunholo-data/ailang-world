# w-world-address-bar-verbs-v4 — S9d sprint plan

## 0. Scope and outcome

Row 243, iteration 451. Design r1 and the r2 quorum verdict (PROCEED 3/3) stay in force. The iteration-450 judge's 77/100 re-ruling allows the existing local path to stand and hard-gates the first remote effect on the named behavioural evidence from this sprint. The scope is the reviewed local-check/outcome and local-commit journal in design §3.3 and §5.4 and the S9 row in §8. D-WORLD-122 remains open under default B. All launch, auth, approval, check and Git fixtures are local and use fake presets, test keys and temp stores.

M1 tests the existing route, claim, worker and status boundaries. M2 applies the S9c Appendix P read-only mirror and its CLI tests. S10 eligibility requires a resolved M1 ledger, a landed M1 and the existing controller gates. Hold or defer M2 independently if an outside CLI mutant survives or the incremental patch or race budget fails.

## 1. Measurements

### 1.1 Base and provenance

Base `db8601b0fc8d7ee70eba7d08692f8686f6b5e707` on `sprint/world-iter451-row243-v4s9d`. `git diff --stat 9b11c617 db8601b0 -- cmd host docs app` was empty. The 41 named survivors come from iteration 450's `ctl_named_mutants.raw` and the exact edits in its `plan-scratch/mutants.json`. The named judge edits come from `judge/mut.py`, with the compilable e21b variant. The judge's report, both re-rulings and mission log §450 are the gate inputs.

The old Appendix P M2 fence is 122 lines. Its dry-run on the fresh base export returned 0 with no offset, fuzz or reject. Rebasing its route-test hunk over M1 created a scratch patch backup because the hunk had moved. The incremental fence in this plan was regenerated against M1 so the replay uses exact context and no patch backup. An independent archive patched with the inherited fence matched the three M2 production and documentation files byte for byte.

### 1.2 Size and cut

| Milestone | Cumulative added | Cumulative deleted | Cumulative changed | Increment |
|---|---:|---:|---:|---:|
| M1 test-first | 747 | 3 | 750 | 750 |
| M2 finish-status mirror | 901 | 7 | 908 | 158 |

Both stages fit under the 1,400-line product cap. In M1 all changed lines are tests. M2 carries the inherited 57-line product increment, including the smoke assertion and QUICKSTART, plus tests that distinguish every new CLI mutation. The expected commit-time product measure is `git diff -M --numstat db8601b0fc8d7ee70eba7d08692f8686f6b5e707 HEAD -- cmd host docs app`. The measured archive diff uses `git diff --no-index -M --numstat` against the same base bytes. The unfiltered archive diff has the same product counts. With both planning files committed separately as new files, the projected unfiltered HEAD sums are 6,182 at M1 and 6,340 at M2 (4,879 plan lines, 553 JSON lines).

Cut order if outside acceptance forces a cut is M2 first. M1 has no selected cut. Keep the full judge precondition set in either case. Any remaining M1 behavioural survivor holds S10.

### 1.3 Baseline and budgets

Primary pristine TEXT fail set is `$STATE/base` at the full base hash above, with the controller's `baseline.done` footer `ALL_DONE for $B`. The iteration-450 base at `91befdab` is attribution context. The controller completed the iter451 baseline. Relevant pristine results were verifygate 0 failures / 314 s, broker 0 / 276 s, cmd plain 12 / 631 s, daemon plain 4 / 813 s and daemon race 5 / 1,248 s (wall). Cmd race has its own skip-11 pristine recording because the iter451 baseline script records cmd plain only. The script emits uptime, raw output, rc, parser attribution and the package fail set for every leg.

CI reference at S9c merge run 46169801018 is 528.9 s cmd, 532.0 s daemon, 534.5 s verifygate and 422.3 s broker. The 900 s whole-package race ceiling applies to the touched race packages and merge CI. A controller reading above 900 s flags a quiet repeat and the merge-CI check. A single contended outside sample does not decide the race budget. Whole cmd and daemon use `-timeout 45m` and whole broker and verifygate use `-timeout 20m`. Compare package-qualified TEXT fail sets and reject timeouts, build failures and unexplained nonzero rc. Use the same skip-11 regex on all cmd race legs. Record both base and stage uptime per named reading. Sum the nonnegative per-test added time and decide on the minimum of up to three pairs, with caps of 15 s cmd and 12 s daemon.

PENDING — controller to measure real-listener timings with `/bin/bash "$S/accept.sh" "$S/st-$M" "$M" outside-race` and the whole race deltas with `outside.sh … whole`. There are no new top-level tests. The new subtests are printed with each outside reading. 

Listener-free diagnostic paired readings in `$S/rt-12/race` (fixture rehearsal only):

- base-tree cmd, 19:50  4 users, load averages: 61.19 37.38 20.96 — `TestOperatorFinishRouteWorker` 21.37 s, `TestOperatorFinishConfigReplacement` 1.11 s.
- base-tree daemon, 19:51  4 users, load averages: 50.54 40.54 23.62 — `TestOperatorFinishAdmission` 1.53 s.
- st-M1 cmd, 19:52  4 users, load averages: 40.36 38.86 23.41 — `TestOperatorFinishRouteWorker` 40.78 s, `TestOperatorFinishConfigReplacement` 2.60 s, measured paired added total 20.90 s.
- st-M1 daemon, 19:53  4 users, load averages: 27.07 35.27 23.28 — `TestOperatorFinishAdmission` 2.63 s, measured paired added total 1.10 s.
- st-M2 cmd, 19:53  4 users, load averages: 34.34 36.18 24.09 — `TestOperatorFinishRouteWorker` 55.42 s, `TestOperatorFinishConfigReplacement` 2.43 s, measured paired added total 35.37 s.
- st-M2 daemon, 19:55  4 users, load averages: 29.72 34.08 24.42 — `TestOperatorFinishAdmission` 7.55 s, measured paired added total 6.02 s.

Earlier listener-free paired diagnostics (sequential, `$S/rt-24`):

- base-tree, 20:14  4 users, load averages: 16.73 21.93 25.24 — `TestOperatorFinishRouteWorker` 14.78 s, `TestOperatorFinishConfigReplacement` 0.91 s, `TestOperatorFinishAdmission` 2.70 s.
- st-M1, 20:14  4 users, load averages: 16.70 21.40 24.93 — `TestOperatorFinishRouteWorker` 34.66 s, `TestOperatorFinishConfigReplacement` 1.92 s, `TestOperatorFinishAdmission` 2.72 s, added cmd 20.89 s and daemon 0.02 s.
- st-M2, 20:16  4 users, load averages: 20.95 22.17 24.92 — `TestOperatorFinishRouteWorker` 44.49 s, `TestOperatorFinishConfigReplacement` 2.18 s, `TestOperatorFinishAdmission` 3.23 s, added cmd 30.98 s and daemon 0.53 s.
The in-memory measurements retain the observed overshoots. The paired real-listener legs decide the race gate.

Final-byte listener-free paired timings in `$S/rt-30` (verbose logs include every subtest):

- base-tree, 20:29  4 users, load averages: 13.17 15.13 19.62 — `TestOperatorFinishRouteWorker` 13.62 s, `TestOperatorFinishConfigReplacement` 0.50 s, `TestOperatorFinishAdmission` 2.62 s.
- st-M1, 20:29  4 users, load averages: 12.83 14.86 19.37 — `TestOperatorFinishRouteWorker` 20.92 s, `TestOperatorFinishConfigReplacement` 1.16 s, `TestOperatorFinishAdmission` 2.26 s, nonnegative added total cmd 7.96 s and daemon 0.00 s.
- st-M2, 20:29  4 users, load averages: 13.06 14.75 19.17 — `TestOperatorFinishRouteWorker` 16.55 s, `TestOperatorFinishConfigReplacement` 0.99 s, `TestOperatorFinishAdmission` 1.75 s, nonnegative added total cmd 3.42 s and daemon 0.00 s.

The earlier listener-free M2 race diagnostic recorded `TestOperatorFinishRouteWorker` 12.26 s, `TestOperatorFinishConfigReplacement` 1.79 s and `TestOperatorFinishAdmission` 1.24 s at load 6.24 / 5.71 / 4.93. The raw diagnostic in `$S/rt-04/mux.raw` prints each subtest reading. The controller supplies the outside budget verdict.

### 1.4 Production callers

| Census | Base | M1 | M2 |
|---|---:|---:|---:|
| `RequestOperatorAction`, `JudgeOperatorAction`, store `ConfirmOperatorAction` | 1 each | 1 each | 1 each |
| `ClaimOperatorAction` | 2 | 2 | 2 |
| `runEpisodeFinishWithPrep`, private local Git commit, attended Git commit | 1 each | 1 each | 1 each |
| `worktree.Push`, dispatch Push, CreateLabel, ViewPR, CreatePR, V3 PutVerbRequest | 1 each | 1 each | 1 each |
| Production test-hook fields | 0 | 0 | 0 |

The grep census counts concrete callers in non-test Go, exercises the 2-to-3 claim positive control and compares the subprocess and context-root site inventories. Existing production AST gates provide independent coverage. M2's read-only dispatch calls the inherited `printFinishStatus` projection. It opens the installed credential DB read-only and uses the existing bounded service context helper. M1 changes only test files, including for the between-claim-and-first-check rows.

## 2. Milestones and five activation gates

| Gate | Measurement and stop condition |
|---|---|
| 1. Scope and boundaries | M1 has a test-only archive diff and matching protected-content hashes. M2 adds only the reviewed CLI mirror and test assertions. Both stages use exact patch anchors, manifest snapshots and the caller/site census. |
| 2. Real-route refusals and fresh observations | Pass `TestOperatorFinishRouteWorker`, `TestOperatorFinishAdmission` and `TestOperatorFinishConfigReplacement` on the planner's stage. Pre-claim rows verify the remedy and durable invariance. The post-claim rows preserve the exact claimed generation, check and commit journals and both refs. Recovery observes one already-reserved commit. |
| 3. Surface and operating envelope | Run the actual listeners, ancestor walk under a 1777 parent, existing Chrome walkthrough, stem selectors and paired named race budgets. Run whole verifygate and the broker subprocess coverage gate on each controller sweep. |
| 4. Product verification and handoff | Pass gofmt, build, native/Linux vet, test compile, pinned `verify_ail.sh`, scope and caller checks, package fail-set comparisons, and §7 patch/snapshot replay. The controller completes the pristine logs before attributing new failures. |
| 5. Mutation disposition | Each exact mutant must show its named Go TEXT failure and, where specified, the target subtest. A genuine overlap requires the expected single and sibling behavior and a killing double. A compile error, inventory response or sandbox denial has no kill value. Any outstanding required row holds S10, and any M2 failure holds the mirror. |

The controller runs §4 against the planner's archive stages before invoking the executor. A survivor or instrument failure returns to the owner with the raw log, exact edit and named subtest. The executor then extracts §7 and applies M1 and M2 in order once the corresponding outside stage gates clear. It may stop after M1 if M2 is deferred by the controller.

## 3. Operator output

M1 changes only the tests. `scope-assets.py` hashes all archive files and asserts the exact changed-file set. It checks goldens, app assets, testdata, AILANG and schema bytes. The browser operator sees the S9c flow and the same local-only durable statuses.

M2 updates QUICKSTART with `ailang-worldd service runner finish-status ACTION-ID`. `TestOperatorFinishRouteWorker/cli/printed-command` reads that line verbatim and substitutes the action ID shown by the card. The installed runner reports the action state, typed check names/results, uncertainty or failed-commit remedy and the existing local commit SHA. The store, worktree, fake child marker and bare remote stay as observed across the CLI read. stdout and stderr are checked separately for protected fixture data.

## 4. Acceptance, timing and fail sets

The planner's exact listener commands returned sandbox-denial text. A separate scratch overlay calls the real mux in memory as a fixture diagnostic. The shipped test sources and named mutation commands retain the actual `httptest` listeners. The diagnostic at `$S/rt-10/mutants.raw` exercised all 20 M2 edits with the expected assertion failures. The listener-free 1777-parent rehearsal passed in `$S/rt-18/mux-1777.raw`. These receipts establish fixture reachability and keep the official named verdicts pending.

The focused controller command runs real listener acceptance, paired race measurements, the sticky-parent run and the named-mutant sweep. The direct named command reruns only the required named ledger with the whole verification guard and the broker subprocess test. The whole mode runs exact package fail-set comparisons. The Chrome mode runs the established walkthrough against disposable fixture state plus the shared guards. Reserve a quiet window for each controller sweep and finish all sweeps before commit guards or judge measurements.

````bash outside-commands
STATE="$HOME/.ailang/state/world-iter451"
S="$STATE/plan-scratch"
B=db8601b0fc8d7ee70eba7d08692f8686f6b5e707"
/bin/bash "$S/gates.sh" "$S/base-tree" "$S/base-extra" M1 base-extra || exit $?
for M in M1 M2; do
  /bin/bash "$S/outside.sh" "$S/st-$M" "$M" focused || exit $?
  /bin/bash "$S/outside.sh" "$S/st-$M" "$M" chrome || exit $?
  /bin/bash "$S/outside.sh" "$S/st-$M" "$M" whole || exit $?
done
# To rerun just the named gate with a fresh caller-chosen log directory:
M=M1
/bin/bash "$S/named-mutants.sh" "$S/st-$M" "$S/named-repeat-$M-$(date +%s)" outside "$M"
# Inherited representative whole-package mutation controls, quiet isolated run:
# /bin/bash "$S/whole-mutants.sh" "$S/st-M1" "$S/whole-mutant-ids.txt" M1
````

PENDING — controller to measure the outside commands above. The M2 named loop runs M2 rows and reports M1 rows as skipped for that call. Keep the M1 ledger attached to M2. The individual raw mutation logs and `ledger` are under the printed log directory. Every wrapper ends with `ALL_DONE` on recorded attempts. Return code 1 means a survivor or test/instrument failure, 77 means sandbox or baseline evidence pending.

The Linux and macOS tests depend on injected fixture clocks, bounded per-write contexts and nonblocking lock probes. The existing runner fixture's test budget overrides the production tick and child budgets. The stage uses physical path resolution on both sides of workspace/store checks, and the sticky-parent command sets a private `TMPDIR` beneath a 1777 directory for the selected subprocesses.

Known baseline-only reds are recorded and compared. These include the 11 cmd real-binary e2es in `CMD_SKIP` and the daemon's `TestSeToolsDescriptorLiveExactEight`, `TestSeToolsMCPVisibilityFollowsGrants` and `TestSeToolsA2AEndToEnd`. The observed iteration-451 cmd plain base recorded the 11 named reds plus `TestOperatorWorktreeWindow`. The rig has load-sensitive fixtures for World context/statusline, host launch, session revoke, worktree/start windows and collisions, file admission and reconcile, service startup, demo seed failures and preflight, broker exec/AILANG confinement, daemon benchmark/A2A/registry clocks/admission, verifygate config and protected Git, and 30 s Git fixtures. Keep the supplied known-red set and classify each observed failure from the actual pristine logs. A timed-out leg is an instrument failure.

## 5. Branch traces and mutation gates

### 5.1 Route worker and evidence

The actual route signs a finish row over the pinned starting and current snapshot. `consumeFinish` reads the pending action, verifies its artifact, runs the pinned fake Git/GitHub readiness, takes the episode flock, freezes and compares inputs, and checks live runner identity and lifecycle before the CAS. The reused route fixture supplies migrated stores, a real fixture SSH signature, a fake agent, a temp bare remote, a fake `gh`, a Git wrapper that sends readiness only to that bare repo, an actual test server for the daemon mux and a fixed-point report transport.

Each eligibility/drift case isolates `finishInputs` and exercises the before-claim constructor with the actual route-confirmed row and verified episode flock. The isolated input result distinguishes the selected check from the later full-manifest comparison. The original route turn still performs pinned fake GitHub/Git readiness. Reusing the real preparation boundary avoids repeating those subprocesses in every test row. A malformed report contains all valid facts plus one invalid typed field. The session-count case adds a second valid-shaped later row while leaving the original session first. Preset drift updates the fixture file, both bound start copies and their byte digest. The session directory permissions and token mode perturbations leave the credentials unchanged. A contrasting live PID birth proves a reused PID is not the stopped child's identity.

Manifest-protected and e12 seed a genuinely signed action over an independently re-frozen manifest that contains exactly one protected or unattributed path. This arrives at the defensive runner guard behind the request filter. The selected runner guard must stop before claim, even though the worker has its own later path gate. The catalogue-size cases preserve valid start offers and check the distinct per-approval and aggregate byte caps.

For post-claim, the test uses the genuine store CAS and verified episode flock, constructs `operatorFinishWork` with `finishWork`, then perturbs heartbeat freshness, launch evidence or the enrolled signing key before `runEpisodeFinishWithPrep`. The existing check seam records any unexpected call. This covers the production `fresh` callback and local worker without a hook. The cancellation case passes a context whose cancellation becomes observable at the actual lifecycle lock boundary. A test-only mutex serializes child-context probes so only the production lock can trigger cancellation. The restart proof case exposes a genuinely created local commit through its original reservation and advances only the injected clock. It requires the original exact journal to be observed after expiry and refuses additional check/commit work.

Read-only state snapshots include the credentials DB hash (action count, pending row, authority counters, claim, receipts and reservations), world hash, local and bare refs, fake-child marker, tracked Git effect commands and check/outcome counts. Some pre-claim checks legitimately call pinned local read subprocesses for SSH verification, Git config and readiness. The tests assert zero new effect subprocesses or launches at their refusal boundaries.

### 5.2 Route and status

The wrong-store request returns HTTP 503 `RunnerUnavailable`. Removing the request handler's store check reaches the store's input guard and returns 400. Protected and unattributed changes return 409 `RunnerUnavailable`. A different generation on the already-known operation and a signed finish lacking the surface option return 409 `Conflict`. The latter is created through the migrated store's real request API and signed using the existing passkey, so both the route admission check and the unchanged counter are observable.

An exact pending finish card and its JSON report the github-auth remedy while the lease is start-available only. The route worker also checks a reserved-only commit shows uncertainty and a definite commit failure retains the check receipt, has an empty local SHA and displays `Local commit failed; check receipts retained`. These GETs preserve every observed state. They use existing store journals and the real daemon status mux.

### 5.3 Equivalence candidates

| Single IDs | Exact sibling and double | Required distinction |
|---|---|---|
| `runner-elig-launch-version`, `e20` | `m1-sibling-completed-launch` broadens only the selection SQL to versions 3 and 4. Run `double-runner-elig-launch-version` and `double-e20`. | A version-3 observed launched receipt is structurally valid and has no live child. The selector rejects it first. The double must admit it to the input result and fail `input/launch-version`. |
| `runner-drift-preset-decode` | `m1-sibling-finish-preset-name` deletes the repeat name check in `finishPreset`. Run `double-runner-drift-preset-decode`. | The byte-pinned input has valid JSON with another name. Either sibling alone refuses review. The double must fail `input/preset-decode`. |
| `runner-heartbeat-owner` | `m1-sibling-cas-heartbeat-owner` ignores owner and generation in the claim SQL while preserving argument bind count. Run `double-runner-heartbeat-owner`. | Stale supplied owner generation and matching live nonce must leave the row unclaimed. The double must fail `claim/heartbeat-owner` on the changed generation/owner state. |

Evidence for every candidate remains PENDING-OUTSIDE. Record both exact overlay edits, separate single and sibling observations and the named subtest failure from the double. Leave the production re-checks in place. e17 remains named as the loop-readiness control and asserts that a start-only tick issues no finish refusal/log or state change. e25 and inherited j35 are must-live controls. The judge's uncompiled e21 form does not qualify as a mutant. The exact compiled e21b variant has its own row.

### 5.4 Existing tests and helpers edited

| Existing function | Before | After | Named control |
|---|---|---|---|
| `TestOperatorFinishRouteWorker` | Signed local finish, existing refusal cases, exact one-commit positive control, durable replay and redaction | Same assertions, re-read confirmed row and additional fixture subtests through `finishInputs`, `finishWork`, `consumeFinish`, `runEpisodeFinishWithPrep` and typed journal, then the inherited M2 CLI and extended mirror checks | `runner-lock-busy`, `runner-positive-claim`, e12, `status-definite-failure` and the M2 ledger |
| `TestOperatorFinishAdmission` | QUICKSTART review, request idempotence, passkey approval and replay | Same behavior plus route/status table before and after approval with exact class and state evidence | `route-finish-admission`, `route-confirm-scope`, e01, e04 and `status-finish-auth` |
| `TestOperatorFinishConfigReplacement` | Unsafe Git config replacement fenced in `snapshot` and `commit` | Same existing assertions plus safe-but-different config hash rejected before `commit-tree` | `j27`, `git-first-use-config`, e14 |

The existing route fixture’s report transport gained the malformed-report case and the fake Git wrapper records effect arguments. Helpers shared by other tests keep their original bodies. New helper code lives in the route test file. The original positive route retains the private task/path, signature/token, check-reference, V3 count, registry and bare-HEAD assertions.

### 5.5 Stem selectors and asset inventories

| Touched file | Named selector coverage on each stage |
|---|---|
| `cmd/ailang-worldd/operator_finish_route_test.go` | `TestOperatorFinishRouteWorker` |
| `cmd/ailang-worldd/operator_finish_tree_test.go` | All six names in `CMD_GIT`: WorktreeFence, CommitIdentity, CaptureBounds, ReadReplacement, CommitReservationFence, ConfigReplacement with the exact full names below |
| `host/daemon/surface_operator_finish_test.go` | `TestOperatorFinishAdmission` |
| `cmd/ailang-worldd/operator_finish_runner.go` (M2) | `TestOperatorFinishReadiness`, `TestOperatorFinishRouteWorker`, `TestOperatorProductionDeadlines` |
| `cmd/ailang-worldd/operator_service.go` (M2) | All 14 exact names in `CMD_SERVICE` plus `TestQuickstartServiceSectionMatchesTheCLIHelp` and CLI inventory controls |
| `docs/QUICKSTART.md` (M2) | All `CMD_DOC` and `DAEMON_DOC` names, `TestOperatorFinishRouteWorker/cli/printed-command`, `TestOperatorFinishAdmission`, `TestOperatorStartWalkthroughInChrome`, marketing `TestMarketingDocsLinks` and `TestMarketingRunbookMatchesHelp`, and broker `TestQuickstartSection10ProfilesLoad`, `TestQuickstartProtectedContract`, `TestExecProtectedSettingsThroughHandler` |

The precise expansions are in `selectors.sh`. The whole-stage asset manifest check asserts no golden, browser asset or pinned app file hash change. The QUICKSTART readers are included even when they read a different section of the same file. The unchanged app export manifest/JS exercises remain named in the carried mutation set.

### 5.6 Named mutation ledger

KILLED in the table means a measured test outside on the exact stage. The e14 socket-free rehearsal is separately identified. The older 80 successful rows are iteration-450 provenance, with a full replay pinned by `named-mutant-ids.txt`. Build-only success, in-memory overlays and historical results do not discharge an outside row.

| Exact id | Source anchor | Named test target | Planner disposition |
|---|---|---|---|
| `runner-elig-resumes` | `cmd/ailang-worldd/operator_finish_runner.go:41` | `TestOperatorFinishRouteWorker` · input/resumes | PENDING-OUTSIDE |
| `runner-elig-preset-allowlist` | `cmd/ailang-worldd/operator_finish_runner.go:41` | `TestOperatorFinishRouteWorker` · input/preset-allowlist | PENDING-OUTSIDE |
| `runner-elig-remote-allowlist` | `cmd/ailang-worldd/operator_finish_runner.go:42` | `TestOperatorFinishRouteWorker` · input/remote-allowlist | PENDING-OUTSIDE |
| `runner-elig-launch-version` | `cmd/ailang-worldd/operator_finish_runner.go:46` | `TestOperatorFinishRouteWorker` · input/launch-version. Candidate `double-runner-elig-launch-version` | PENDING-OUTSIDE |
| `runner-elig-health-status` | `cmd/ailang-worldd/operator_finish_runner.go:53` | `TestOperatorFinishRouteWorker` · input/health-status | PENDING-OUTSIDE |
| `runner-elig-build` | `cmd/ailang-worldd/operator_finish_runner.go:54` | `TestOperatorFinishRouteWorker` · input/build | PENDING-OUTSIDE |
| `runner-elig-report-read` | `cmd/ailang-worldd/operator_finish_runner.go:58` | `TestOperatorFinishRouteWorker` · input/report-read | PENDING-OUTSIDE |
| `runner-elig-report-episode` | `cmd/ailang-worldd/operator_finish_runner.go:58` | `TestOperatorFinishRouteWorker` · input/report-episode | PENDING-OUTSIDE |
| `runner-elig-report-intents` | `cmd/ailang-worldd/operator_finish_runner.go:58` | `TestOperatorFinishRouteWorker` · input/report-intents | PENDING-OUTSIDE |
| `runner-elig-recorded-finish` | `cmd/ailang-worldd/operator_finish_runner.go:58` | `TestOperatorFinishRouteWorker` · input/recorded-finish | PENDING-OUTSIDE |
| `runner-elig-report-start` | `cmd/ailang-worldd/operator_finish_runner.go:58` | `TestOperatorFinishRouteWorker` · input/report-start | PENDING-OUTSIDE |
| `runner-drift-preset-decode` | `cmd/ailang-worldd/operator_finish_runner.go:70` | `TestOperatorFinishRouteWorker` · input/preset-decode. Candidate `double-runner-drift-preset-decode` | PENDING-OUTSIDE |
| `runner-drift-workspace` | `cmd/ailang-worldd/operator_finish_runner.go:76` | `TestOperatorFinishRouteWorker` · input/workspace | PENDING-OUTSIDE |
| `runner-drift-world` | `cmd/ailang-worldd/operator_finish_runner.go:79` | `TestOperatorFinishRouteWorker` · input/world | PENDING-OUTSIDE |
| `runner-drift-store` | `cmd/ailang-worldd/operator_finish_runner.go:83` | `TestOperatorFinishRouteWorker` · input/store | PENDING-OUTSIDE |
| `runner-drift-profile` | `cmd/ailang-worldd/operator_finish_runner.go:104` | `TestOperatorFinishRouteWorker` · input/profile | PENDING-OUTSIDE |
| `runner-drift-pinned-check` | `cmd/ailang-worldd/operator_finish_runner.go:96` | `TestOperatorFinishRouteWorker` · input/pinned-check | PENDING-OUTSIDE |
| `runner-elig-sessions` | `cmd/ailang-worldd/operator_finish_runner.go:100` | `TestOperatorFinishRouteWorker` · input/sessions | PENDING-OUTSIDE |
| `runner-elig-project` | `cmd/ailang-worldd/operator_finish_runner.go:101` | `TestOperatorFinishRouteWorker` · input/project | PENDING-OUTSIDE |
| `runner-elig-session-file` | `cmd/ailang-worldd/operator_finish_runner.go:108` | `TestOperatorFinishRouteWorker` · input/session-file | PENDING-OUTSIDE |
| `runner-elig-token-mode` | `cmd/ailang-worldd/operator_finish_runner.go:112` | `TestOperatorFinishRouteWorker` · input/token-mode | PENDING-OUTSIDE |
| `runner-elig-commit-identity` | `cmd/ailang-worldd/operator_finish_runner.go:130` | `TestOperatorFinishRouteWorker` · input/commit-identity | PENDING-OUTSIDE |
| `runner-manifest-protected` | `cmd/ailang-worldd/operator_finish_runner.go:278` | `TestOperatorFinishRouteWorker` · claim/manifest-protected | PENDING-OUTSIDE |
| `runner-catalogue-fresh` | `cmd/ailang-worldd/operator_finish_runner.go:293` | `TestOperatorFinishRouteWorker` · claim/catalogue-fresh | PENDING-OUTSIDE |
| `runner-heartbeat-owner` | `cmd/ailang-worldd/operator_finish_runner.go:293` | `TestOperatorFinishRouteWorker` · claim/heartbeat-owner. Candidate `double-runner-heartbeat-owner` | PENDING-OUTSIDE |
| `runner-heartbeat-ready` | `cmd/ailang-worldd/operator_finish_runner.go:293` | `TestOperatorFinishRouteWorker` · claim/heartbeat-ready | PENDING-OUTSIDE |
| `runner-life-quiet` | `cmd/ailang-worldd/operator_finish_runner.go:285` | `TestOperatorFinishRouteWorker` · claim/life-quiet | PENDING-OUTSIDE |
| `runner-life-pair` | `cmd/ailang-worldd/operator_finish_runner.go:285` | `TestOperatorFinishRouteWorker` · claim/life-pair | PENDING-OUTSIDE |
| `runner-effect-clock` | `cmd/ailang-worldd/operator_finish_runner.go:320` | `TestOperatorFinishRouteWorker` · after-claim/clock | PENDING-OUTSIDE |
| `runner-effect-inputs` | `cmd/ailang-worldd/operator_finish_runner.go:328` | `TestOperatorFinishRouteWorker` · after-claim/inputs | PENDING-OUTSIDE |
| `runner-effect-key` | `cmd/ailang-worldd/operator_finish_runner.go:331` | `TestOperatorFinishRouteWorker` · after-claim/key | PENDING-OUTSIDE |
| `runner-observe-after-expiry` | `cmd/ailang-worldd/operator_finish_runner.go:203` | `TestOperatorFinishRouteWorker` · recover exact original journal after deadline | PENDING-OUTSIDE |
| `runner-old-observation` | `cmd/ailang-worldd/operator_finish_runner.go:198` | `TestOperatorFinishRouteWorker` · recover the reserved-commit path | PENDING-OUTSIDE |
| `runner-catalogue-approve-size` | `cmd/ailang-worldd/operator_finish_runner.go:171` | `TestOperatorFinishRouteWorker` · catalogue/approval-size | PENDING-OUTSIDE |
| `runner-catalogue-bound` | `cmd/ailang-worldd/operator_catalogue.go:45` | `TestOperatorFinishRouteWorker` · catalogue/aggregate-size | PENDING-OUTSIDE |
| `route-finish-store` | `host/daemon/surface_operator.go:53` | `TestOperatorFinishAdmission` · other-store | PENDING-OUTSIDE |
| `route-protected-and-unattributed` | `host/daemon/surface_operator.go:56` | `TestOperatorFinishAdmission` · protected | PENDING-OUTSIDE |
| `route-operation-generation-conflict` | `host/daemon/surface_operator.go:241` | `TestOperatorFinishAdmission` · existing-operation/other-generation | PENDING-OUTSIDE |
| `route-confirm-without-surface` | `host/daemon/surface_operator_confirm.go:96` | `TestOperatorFinishAdmission` · confirm/without-surface-option | PENDING-OUTSIDE |
| `status-definite-failure` | `host/daemon/surface_operator_read.go:143` | `TestOperatorFinishRouteWorker` · exact failed-commit reason plus original check receipt | PENDING-OUTSIDE |
| `status-finish-auth` | `host/daemon/surface_operator_read.go:151` | `TestOperatorFinishAdmission` · status/start-only-runner | PENDING-OUTSIDE |
| `s9d-cli-dispatch` | `cmd/ailang-worldd/operator_service.go:318` | `TestOperatorFinishRouteWorker` · inherited dispatch to printFinishStatus | PENDING-OUTSIDE |
| `s9d-cli-sha` | `cmd/ailang-worldd/operator_finish_runner.go:384` | `TestOperatorFinishRouteWorker` · positive mirror returns exact current commit | PENDING-OUTSIDE |
| `e01` | `host/daemon/surface_operator_confirm.go:96` | `TestOperatorFinishAdmission` · confirm/without-surface-option (same gate as route-confirm-without-surface) | PENDING-OUTSIDE |
| `e04` | `host/daemon/surface_operator.go:56` | `TestOperatorFinishAdmission` · unattributed (unattributed half of route-protected-and-unattributed) | PENDING-OUTSIDE |
| `e06` | `cmd/ailang-worldd/operator_finish_runner.go:41` | `TestOperatorFinishRouteWorker` · input/resumes (runner-elig-resumes) | PENDING-OUTSIDE |
| `e07` | `cmd/ailang-worldd/operator_finish_runner.go:58` | `TestOperatorFinishRouteWorker` · input/report-intents (runner-elig-report-intents) | PENDING-OUTSIDE |
| `e08` | `cmd/ailang-worldd/operator_finish_runner.go:41` | `TestOperatorFinishRouteWorker` · input/build (runner-elig-build) | PENDING-OUTSIDE |
| `e09` | `cmd/ailang-worldd/operator_finish_runner.go:42` | `TestOperatorFinishRouteWorker` · input/remote-allowlist (runner-elig-remote-allowlist) | PENDING-OUTSIDE |
| `e12` | `cmd/ailang-worldd/operator_finish_runner.go:278` | `TestOperatorFinishRouteWorker` · claim/manifest-unattributed (runner-side unattributed half beside runner-manifest-protected) | PENDING-OUTSIDE |
| `e13` | `cmd/ailang-worldd/operator_finish_runner.go:293` | `TestOperatorFinishRouteWorker` · claim/catalogue-fresh (runner-catalogue-fresh) | PENDING-OUTSIDE |
| `e14` | `cmd/ailang-worldd/operator_finish_tree.go:348` | `TestOperatorFinishConfigReplacement` · valid-config-with-different-digest (exact config digest, distinct from j27's replacement check) | KILLED in sandbox rehearsal. Outside replay required. |
| `e16` | `host/daemon/surface_operator_read.go:141` | `TestOperatorFinishRouteWorker` · reserved-only status, subset of status-open-window | PENDING-OUTSIDE |
| `e17` | `cmd/ailang-worldd/operator_runner.go:157` | `TestOperatorFinishRouteWorker` · claim/start-ready-loop (start-only loop readiness control) | PENDING-OUTSIDE |
| `e20` | `cmd/ailang-worldd/operator_finish_runner.go:46` | `TestOperatorFinishRouteWorker` · input/launch-version (runner-elig-launch-version). Candidate `double-e20` | PENDING-OUTSIDE |
| `e22` | `cmd/ailang-worldd/operator_finish_runner.go:203` | `TestOperatorFinishRouteWorker` · runner-observe-after-expiry | PENDING-OUTSIDE |
| `e23` | `cmd/ailang-worldd/operator_finish_runner.go:270` | `TestOperatorFinishRouteWorker` · after-claim/review (pre-effect comparison on construction, distinct from runner-effect-inputs) | PENDING-OUTSIDE |
| `e24` | `cmd/ailang-worldd/operator_finish_runner.go:285` | `TestOperatorFinishRouteWorker` · claim/cancelled (cancellation at lifecycle lock) | PENDING-OUTSIDE |
| `e25` | `cmd/ailang-worldd/operator_catalogue.go:40` | `TestOperatorFinishRouteWorker` · must-live catalogue budget control | PENDING-OUTSIDE — MUST-LIVE |
| `e21b` | `cmd/ailang-worldd/operator_finish_runner.go:49` | `TestOperatorFinishRouteWorker` · input/reused-pid (birth equality, distinct from the already-live child refusal) | PENDING-OUTSIDE |
| `m1-sibling-completed-launch` | `host/store/operator_finish.go:126` | `TestOperatorFinishRouteWorker` · input/launch-version | PENDING-OUTSIDE — MUST-LIVE |
| `m1-sibling-finish-preset-name` | `cmd/ailang-worldd/episode_finish.go:569` | `TestOperatorFinishRouteWorker` · input/preset-decode | PENDING-OUTSIDE — MUST-LIVE |
| `m1-sibling-cas-heartbeat-owner` | `host/store/operator_cas.go:140` | `TestOperatorFinishRouteWorker` · claim/heartbeat-owner | PENDING-OUTSIDE — MUST-LIVE |
| `double-runner-elig-launch-version` | `cmd/ailang-worldd/operator_finish_runner.go:46` | `TestOperatorFinishRouteWorker` · input/launch-version. Edits `runner-elig-launch-version` + `m1-sibling-completed-launch` | PENDING-OUTSIDE |
| `double-e20` | `cmd/ailang-worldd/operator_finish_runner.go:46` | `TestOperatorFinishRouteWorker` · input/launch-version. Edits `e20` + `m1-sibling-completed-launch` | PENDING-OUTSIDE |
| `double-runner-drift-preset-decode` | `cmd/ailang-worldd/operator_finish_runner.go:70` | `TestOperatorFinishRouteWorker` · input/preset-decode. Edits `runner-drift-preset-decode` + `m1-sibling-finish-preset-name` | PENDING-OUTSIDE |
| `double-runner-heartbeat-owner` | `cmd/ailang-worldd/operator_finish_runner.go:293` | `TestOperatorFinishRouteWorker` · claim/heartbeat-owner. Edits `runner-heartbeat-owner` + `m1-sibling-cas-heartbeat-owner` | PENDING-OUTSIDE |
| `s9d-cli-help` | `cmd/ailang-worldd/operator_service.go:24` | `TestOperatorFinishRouteWorker` · cli/help | PENDING-OUTSIDE |
| `s9d-cli-verb` | `cmd/ailang-worldd/operator_service.go:294` | `TestOperatorFinishRouteWorker` · allowed verb in input validation | PENDING-OUTSIDE |
| `s9d-cli-id-grammar` | `cmd/ailang-worldd/operator_service.go:295` | `TestOperatorFinishRouteWorker` · cli/invalid-id | PENDING-OUTSIDE |
| `s9d-cli-arity` | `cmd/ailang-worldd/operator_service.go:295` | `TestOperatorFinishRouteWorker` · cli/extra-argument | PENDING-OUTSIDE |
| `s9d-cli-branch` | `cmd/ailang-worldd/operator_service.go:317` | `TestOperatorFinishRouteWorker` · dispatch branch | PENDING-OUTSIDE |
| `s9d-cli-dispatch-error` | `cmd/ailang-worldd/operator_service.go:318` | `TestOperatorFinishRouteWorker` · cli/unknown-action | PENDING-OUTSIDE |
| `s9d-cli-installed-remedy` | `cmd/ailang-worldd/operator_finish_runner.go:357` | `TestOperatorFinishRouteWorker` · cli/uninstalled | PENDING-OUTSIDE |
| `s9d-cli-read-store-remedy` | `cmd/ailang-worldd/operator_finish_runner.go:152` | `TestOperatorFinishRouteWorker` · cli/unreadable-credentials | PENDING-OUTSIDE |
| `s9d-cli-unknown-action` | `cmd/ailang-worldd/operator_finish_runner.go:369` | `TestOperatorFinishRouteWorker` · cli/unknown-action | PENDING-OUTSIDE |
| `s9d-cli-non-finish` | `cmd/ailang-worldd/operator_finish_runner.go:369` | `TestOperatorFinishRouteWorker` · cli/non-finish | PENDING-OUTSIDE |
| `s9d-cli-missing-id-remedy` | `cmd/ailang-worldd/operator_finish_runner.go:369` | `TestOperatorFinishRouteWorker` · cli/unknown-action | PENDING-OUTSIDE |
| `s9d-cli-journal-remedy` | `cmd/ailang-worldd/operator_finish_runner.go:372` | `TestOperatorFinishRouteWorker` · cli/unreadable-journal | PENDING-OUTSIDE |
| `s9d-cli-redact-stderr` | `cmd/ailang-worldd/operator_finish_runner.go:372` | `TestOperatorFinishRouteWorker` · cli/unreadable-journal (unreadable journal with signed-action material) | PENDING-OUTSIDE |
| `s9d-cli-redact-stdout` | `cmd/ailang-worldd/operator_finish_runner.go:378` | `TestOperatorFinishRouteWorker` · observed check with private evidence reference | PENDING-OUTSIDE |
| `s9d-cli-open-check` | `cmd/ailang-worldd/operator_finish_runner.go:380` | `TestOperatorFinishRouteWorker` · cli/open-check | PENDING-OUTSIDE |
| `s9d-cli-reserved` | `cmd/ailang-worldd/operator_finish_runner.go:380` | `TestOperatorFinishRouteWorker` · reserved commit with no open check | PENDING-OUTSIDE |
| `s9d-cli-done-state` | `cmd/ailang-worldd/operator_finish_runner.go:383` | `TestOperatorFinishRouteWorker` · done state and exact SHA | PENDING-OUTSIDE |
| `s9d-cli-failed-state` | `cmd/ailang-worldd/operator_finish_runner.go:385` | `TestOperatorFinishRouteWorker` · definite failure with empty SHA | PENDING-OUTSIDE |

The final pre-claim fixture diagnostic is `$S/rt-27/mutants.raw`. It reruns the three overlap families, e12, the protected-manifest row and lifecycle cancellation. These results do not discharge outside rows.

The listener-free M1 diagnostic ledgers are `$S/rt-05/mutants.raw` and `$S/rt-15/mutants.raw`. The second reruns the different-world input with a distinct migrated copy, so it isolates the compared World path and retains the distinct credentials DB. Use the named outside run for dispositions.


## 6. Files, guards and gate scripts

All tools live under `$STATE/plan-scratch` for controller runs. Executor tools are extracted from these fences into `$STATE/exec-tools`. The following script contracts preserve the S9c names and include real TEXT parser and caller controls.

**Local guard receipts**

- M1 earlier rehearsal `$S/rt-07/legs.txt`

````text
LEG gofmt OK 0s
LEG build OK 8s
LEG vet OK 3s
LEG linux-vet OK 7s
LEG test-compile OK 37s
LEG inventory-store OK 2s
LEG inventory-cli OK 2s
LEG inventory-daemon OK 5s
LEG inventory-boundary OK 6s
LEG subprocess-sink OK 3s
LEG subprocess-sites UNINFORMATIVE UNDER SANDBOX rc=1
LEG scope-assets OK 11s
LEG verify-ail OK 174s
LEG caller-census OK 11s
LEG accept OK 64s
LEG accept-race OK 109s
LEG accept-1777 OK 19s
LEG cmd-plain SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG cmd-race SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG daemon-plain SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG daemon-race SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG verifygate-plain INSTRUMENT FAILURE rc=1 parser=1
LEG broker-plain UNINFORMATIVE UNDER SANDBOX rc=1
LEG baseline-completion SKIPPED (controller recording pristine legs)
LEG compare-cmd-plain SKIPPED (stage evidence unavailable)
LEG compare-cmd-race SKIPPED (stage evidence unavailable)
LEG compare-daemon-plain SKIPPED (stage evidence unavailable)
LEG compare-daemon-race SKIPPED (stage evidence unavailable)
LEG compare-verifygate-plain SKIPPED (stage evidence unavailable)
LEG compare-broker-plain SKIPPED (stage evidence unavailable)
LEG cmd-race-race-bound SKIPPED (whole race evidence outside)
LEG daemon-race-race-bound SKIPPED (whole race evidence outside)
````

- M2 earlier rehearsal `$S/rt-14/legs.txt`

````text
LEG gofmt OK 1s
LEG build OK 13s
LEG vet OK 5s
LEG linux-vet OK 9s
LEG test-compile OK 50s
LEG inventory-store OK 4s
LEG inventory-cli OK 4s
LEG inventory-daemon OK 4s
LEG inventory-boundary OK 4s
LEG subprocess-sink OK 4s
LEG subprocess-sites UNINFORMATIVE UNDER SANDBOX rc=1
LEG scope-assets OK 10s
LEG verify-ail OK 244s
LEG caller-census OK 13s
LEG accept OK 84s
LEG accept-race OK 79s
LEG accept-1777 OK 10s
LEG cmd-plain SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG cmd-race SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG daemon-plain SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG daemon-race SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG verifygate-plain recorded rc=1 948s
LEG broker-plain UNINFORMATIVE UNDER SANDBOX rc=1
LEG compare-cmd-plain SKIPPED (stage evidence unavailable)
LEG compare-cmd-race SKIPPED (stage evidence unavailable)
LEG compare-daemon-plain SKIPPED (stage evidence unavailable)
LEG compare-daemon-race SKIPPED (stage evidence unavailable)
LEG compare-verifygate-plain NEW FAIL SET
LEG compare-broker-plain SKIPPED (stage evidence unavailable)
LEG cmd-race-race-bound SKIPPED (whole race evidence outside)
LEG daemon-race-race-bound SKIPPED (whole race evidence outside)
````

- M2 sealed script `$S/rt-25/gates/legs.txt`

````text
LEG gofmt OK 0s
LEG build OK 6s
LEG vet OK 2s
LEG linux-vet OK 4s
LEG test-compile OK 5s
LEG inventory-store OK 6s
LEG inventory-cli OK 5s
LEG inventory-daemon OK 5s
LEG inventory-boundary OK 5s
LEG subprocess-sink OK 3s
LEG subprocess-sites UNINFORMATIVE UNDER SANDBOX rc=1
LEG scope-assets OK 9s
LEG verify-ail OK 136s
LEG caller-census OK 4s
LEG accept OK 26s
LEG accept-race OK 15s
LEG accept-1777 OK 2s
LEG cmd-plain SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG cmd-race SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG daemon-plain SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG daemon-race SKIPPED (UNINFORMATIVE UNDER SANDBOX, outside.sh whole)
LEG verifygate-plain recorded rc=0 471s
LEG broker-plain UNINFORMATIVE UNDER SANDBOX rc=1
LEG compare-cmd-plain SKIPPED (stage evidence unavailable)
LEG compare-cmd-race SKIPPED (stage evidence unavailable)
LEG compare-daemon-plain SKIPPED (stage evidence unavailable)
LEG compare-daemon-race SKIPPED (stage evidence unavailable)
LEG compare-verifygate-plain OK
LEG compare-broker-plain SKIPPED (stage evidence unavailable)
LEG cmd-race-race-bound SKIPPED (whole race evidence outside)
LEG daemon-race-race-bound SKIPPED (whole race evidence outside)
GATE_RESULT issue=0 pending=1
ALL_DONE
````

The targeted mission-config repro passed on base and M2 in `$S/rt-26`. The whole-race elapsed-time control printed `RACE_OVERSHOOT_PENDING_CONTROL_OK` in `$S/rt-32`. The final archive stages passed gofmt, build, native/Linux vet, test compile and scope checks in `$S/rt-28/{M1,M2}/final-static.raw`. The older gate rehearsal captured a whole verifygate timeout and broker sandbox denial. Its completion footer was interrupted by an in-place script edit, so the gate was rerun from sealed extracted fences.

Bash 3.2.57 with `set -u` rehearse receipts. Rc 77 indicates an explicitly skipped outside prerequisite.

````text
env.sh syntax=0
selectors.sh syntax=0
accept.sh syntax=0
gates.sh syntax=0
snapshot.sh syntax=0
callers.sh syntax=0
outside.sh syntax=0
commit-guards.sh syntax=0
named-mutants.sh syntax=0
whole-mutants.sh syntax=0
apply.sh syntax=0
commit-M1.sh syntax=0
commit-M2.sh syntax=0
env.sh runtime=0
selectors.sh runtime=0
outside.sh-focused runtime=77
outside.sh-whole runtime=77
outside.sh-chrome runtime=77
whole-mutants.sh runtime=0
callers.sh runtime=0
commit-M1.sh runtime=77
commit-M2.sh runtime=77
commit-guards.sh runtime=77
gates.sh-M2 runtime=2
gates.sh-sealed-M2 runtime=77
accept.sh-fast-M1 runtime=0
accept.sh-fast-M2 runtime=0
snapshot.sh-M1 runtime=0
snapshot.sh-M2 runtime=0
apply.sh-final runtime=0
named-mutants.sh runtime=77
race-bound-control.sh runtime=0
gates.sh-M1-earlier runtime=2
````

Host verifygate, the broker subprocess test, whole broker, Chrome, and full cmd/daemon comparisons require the controller's outside evidence where local output is marked uninformative. The runtime scope sensor also pins zero production test hooks.

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
root=${REHEARSAL_ROOT:-"$S"}
L=$(mktemp -d "$root/accept.$M.$mode.XXXXXX") || exit 2
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" || exit 2
cd "$T" || exit 2
echo "ACCEPT_LOGS $L"
python3 "$TOOLS/failset.py" --control || exit 2
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
  local pkg=$1 selector=$2 cap=$3 attempt tag delta min leg tree prior=$result
  shift 3
  local totals="$L/${pkg#*/}.totals"
  : > "$totals"
  result=0
  for attempt in 1 2 3; do
    for leg in base stage; do
      tag="${pkg#*/}-race-$attempt-$leg"
      if [ "$leg" = base ]; then tree="$S/base-tree"; else tree="$T"; fi
      cd "$tree" || exit 2
      named "$tag" "$pkg" "$selector" -race "$@"
      [ "$result" -ne 0 ] || require_all "$L/$tag.raw" "$selector"
    done
    cd "$T" || exit 2
    [ "$result" -eq 0 ] || return
    awk -v pkg="$pkg" -v selector="$selector" -v read="$attempt" '
      $1=="---" && $2=="PASS:" && $3~("^(" selector ")$") {
        t=$4; gsub(/[()s]/,"",t)
        if (NR==FNR) base[$3]=t
        else {
          delta=t-base[$3]; if(delta<0) delta=0; sum+=delta
          printf "RACE_TEST %s reading=%s %s base=%.2fs stage=%.2fs added=%.2fs\n",pkg,read,$3,base[$3],t,delta
        }
      }
      END { printf "%.2f\n",sum }' "$L/${pkg#*/}-race-$attempt-base.raw" "$L/${pkg#*/}-race-$attempt-stage.raw" > "$L/${pkg#*/}-race-$attempt.reading"
    sed '$d' "$L/${pkg#*/}-race-$attempt.reading"
    grep -E '^[[:space:]]+--- PASS: [^[:space:]]+/[^[:space:]]+ \(' "$L/${pkg#*/}-race-$attempt-stage.raw" || true
    delta=$(tail -1 "$L/${pkg#*/}-race-$attempt.reading")
    echo "$delta" >> "$totals"
    note "RACE_READING $pkg reading=$attempt added=${delta}s loads=$L/${pkg#*/}-race-$attempt-{base,stage}.load"
    if awk -v reading="$delta" -v cap="$cap" 'BEGIN { exit !(reading<=cap) }'; then break; fi
  done
  min=$(sort -n "$totals" | head -1)
  note "RACE_BUDGET $pkg added_min=${min}s cap=${cap}s readings=$(paste -sd, "$totals")"
  awk -v reading="$min" -v cap="$cap" 'BEGIN { exit !(reading<=cap) }' || result=1
  [ "$result" -ne 0 ] || result=$prior
}
case "$mode" in
  fast|outside)
    named cmd-core cmd/ailang-worldd "$CMD_S9D_LOCAL"
    if [ "$mode" = outside ]; then
      named cmd-route-and-stems cmd/ailang-worldd "$CMD_LISTENER|$CMD_SERVICE|$CMD_DOC"
      named daemon-route-and-doc host/daemon "$DAEMON_ADDED|$DAEMON_DOC"
      named broker-doc host/broker 'TestExecProtectedSettingsThroughHandler|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract'
      named marketing-doc scripts/marketingcheck 'TestMarketingDocsLinks|TestMarketingRunbookMatchesHelp'
      if [ "$result" -eq 0 ]; then
        require_all "$L/cmd-route-and-stems.raw" "$CMD_LISTENER"
        require_all "$L/daemon-route-and-doc.raw" "$DAEMON_ADDED"
      fi
    else
      skip cmd-route-and-stems 'real listener and full operator service neighbours outside'
      skip daemon-route-and-doc 'real listener outside'
      skip broker-doc 'outside executable and listener fixtures'
    fi
    ;;
  race)
    named cmd-race-local cmd/ailang-worldd TestOperatorFinishConfigReplacement -race -skip "$CMD_SKIP"
    skip full-added-race-budget 'real listener outside'
    ;;
  outside-race)
    budget cmd/ailang-worldd "$CMD_S9D_RACE" 15 -skip "$CMD_SKIP"
    budget host/daemon "$DAEMON_ADDED" 12
    ;;
  1777|outside-1777)
    mkdir -m 1777 "$L/t1777" && mkdir -m 0700 "$L/t1777/inside" || exit 2
    if [ "$mode" = outside-1777 ]; then
      TMPDIR="$L/t1777/inside" named cmd-1777 cmd/ailang-worldd "$CMD_LISTENER|$CMD_S9D_LOCAL|$CMD_SERVICE"
      TMPDIR="$L/t1777/inside" named daemon-1777 host/daemon "$DAEMON_ADDED"
    else
      TMPDIR="$L/t1777/inside" named cmd-1777 cmd/ailang-worldd 'TestOperatorFinishConfigReplacement|TestOperatorFinishToolPins|TestRunnerPathRefusalClasses'
      skip listener-1777 'real listener outside'
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
python3 "$TOOLS/failset.py" --control || exit 2
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
  if [ "$tag" = cmd-race ]; then source="${CMD_RACE_BASE:-"$S/base-extra"}/$tag"; fi
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
  packages=(./cmd/ailang-worldd ./host/daemon)
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
  run subprocess-sites inventory host/broker TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential
  run scope-assets python3 "$TOOLS/scope-assets.py" "$S/base-tree" "$T" "$TOOLS" "$M"
  run verify-ail bound 1000 python3 "$TOOLS/verify-scratch.py" "$T" "$L" verify_ail ""
  run caller-census /bin/bash "$TOOLS/callers.sh" "$S/base-tree" "$T" "$M"

  if [ "$mode" = local ]; then
    run accept /bin/bash "$TOOLS/accept.sh" "$T" "$M" fast
    run accept-race /bin/bash "$TOOLS/accept.sh" "$T" "$M" race
    run accept-1777 /bin/bash "$TOOLS/accept.sh" "$T" "$M" 1777
    for tag in cmd-plain cmd-race daemon-plain daemon-race; do
      pending=1; skip "$tag" 'UNINFORMATIVE UNDER SANDBOX, outside.sh whole'
    done
  else
    run accept /bin/bash "$TOOLS/accept.sh" "$T" "$M" outside
    run accept-race /bin/bash "$TOOLS/accept.sh" "$T" "$M" outside-race
    run accept-1777 /bin/bash "$TOOLS/accept.sh" "$T" "$M" outside-1777
    record cmd-plain cmd/ailang-worldd 45m 2900
    record cmd-race cmd/ailang-worldd 45m 2900 -race
    record daemon-plain host/daemon 45m 2900
    record daemon-race host/daemon 45m 2900 -race
  fi
  record verifygate-plain host/verifygate 20m 1700
  record broker-plain host/broker 20m 1700
  if ! grep -Fxq "ALL_DONE for $B" "$STATE/baseline.done" 2>/dev/null; then
    pending=1; skip baseline-completion 'controller recording pristine legs'
  fi
  for tag in cmd-plain cmd-race daemon-plain daemon-race verifygate-plain broker-plain; do
    case "$tag" in cmd-*) pkg=cmd/ailang-worldd;; *) pkg=host/${tag%-*};; esac
    compare "$tag" "$pkg"
  done
  for tag in cmd-race daemon-race; do
    if [ -f "$L/$tag.verified" ]; then
      seconds=$(awk '$1=="ok" || $1=="FAIL" {if ($3~/^[0-9.]+s$/) {sub(/s$/,"",$3);print $3}}' "$L/$tag.raw")
      if [ -z "$seconds" ]; then
        issue=1; note "LEG $tag race-bound INSTRUMENT FAILURE missing package time"
      elif ! awk -v s="$seconds" 'BEGIN {exit !(s<=900)}'; then
        pending=1; note "LEG $tag race-bound PENDING quiet repeat seconds=$seconds cap=900 load=$L/$tag.load"
      else
        note "LEG $tag race-bound OK seconds=$seconds"
      fi
    else
      skip "$tag-race-bound" 'whole race evidence outside'
    fi
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
B=db8601b0fc8d7ee70eba7d08692f8686f6b5e707
D=${SNAP_PARENT:-$HOME/.ailang/state/world-iter451/snap}/$M
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
    $2 != "design_docs/planned/w-world-address-bar-verbs-v4s9d-sprint-plan.md" &&
    $2 != "design_docs/planned/sprint_w-world-address-bar-verbs-v4s9d.json"
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
STATE="$HOME/.ailang/state/world-iter451"
S="$STATE/plan-scratch"
B=db8601b0fc8d7ee70eba7d08692f8686f6b5e707
export STATE
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
export GOPROXY=off GIT_OPTIONAL_LOCKS=0
unset GOFLAGS WORLD_PASSKEY_RECORD_DIR WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
````

````bash selectors.sh
# Exact selection fragments for iteration 451. Selectors may be combined with |.
CMD_SKIP='^(TestCLIAgainstSeToolsDaemon|TestCLINoGenesisIsNoWorldHead|TestDemoSeedApprovalPacket|TestDemoSeedEndToEnd|TestDemoSeedManifestReferences|TestDemoSeedProofStage|TestDemoSeedReadTestCheck|TestDemoSeedRejectsToolReportRefusals|TestDemoSeedRosterKindMismatchRefuses|TestDemoSeedWhyNamedLinks|TestMarketingStillRoutes)$'
CMD_ADDED='TestOperatorFinishReadiness|TestOperatorFinishResumeEvidence|TestOperatorCommitFailureInputRecheck|TestOperatorCommitResumeReserved'
CMD_LISTENER='TestOperatorFinishRouteWorker'
DAEMON_ADDED='TestOperatorFinishAdmission'
CMD_FENCES='TestOperatorFinishF1|TestOperatorFinishTTYParity|TestOperatorFinishNoProof|TestOperatorFinishProofScope|TestOperatorFinishAuth|TestOperatorFinishToolPins|TestOperatorFlagParity|TestOperatorTerminalExclusion'
CMD_JOURNAL='TestOperatorFinishPrep|TestOperatorFinishPrepExpiry|TestOperatorFinishStale|TestOperatorFinishChecks|TestOperatorFinishUnsandboxed|TestOperatorCommitFailed|TestOperatorCommitReconcile|TestOperatorFinishJournalOnce|TestOperatorFinishManifestCapture|TestOperatorFinishWorldEvidence'
CMD_GIT='TestOperatorFinishWorktreeFence|TestOperatorCommitIdentity|TestOperatorFinishCaptureBounds|TestOperatorFinishReadReplacement|TestOperatorCommitReservationFence|TestOperatorFinishConfigReplacement'
CMD_SERVICE='TestOperatorServiceUpgrade|TestOperatorServiceFence|TestOperatorServiceOwnership|TestOperatorServiceConfig|TestOperatorServicePair|TestOperatorServiceStartup|TestOperatorServiceUpgradePairFence|TestOperatorServiceCancelConfirmation|TestOperatorServiceWriteScope|TestOperatorServiceEnvironment|TestRunnerPathRefusalClasses|TestOperatorRunnerKernelLock|TestOperatorStatusDisconnected|TestOperatorFinishAuthInstall'
CMD_DOC='TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartSection18FlagsMatchTheCLI|TestQuickstartSection18CardMatchesThePage|TestQuickstartSection19MatchesTheBar|TestQuickstartSection20MatchesTheFinishFlag|TestQuickstartSection22RunsVerbatim|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartServiceSectionMatchesTheCLIHelp|TestQuickstartRunnerMintRecovery|TestProfileInitQuickstartLineMatchesTheCLI|TestDemoSeedHelpPreparation|TestDemoSeedRosterPageControls|TestOperatorWorktreeWindow|TestOperatorFilePublicationAdmission'
DAEMON_DOC='TestMCPQuickstartPayloadsVerbatim|TestAGUIQuickstartControl|TestQuickstartDecisionPacketSection|TestWorkbenchLiveQuickstartControl|TestSeToolsQuickstartPayloadsVerbatim|TestExecStartupRefusalTable|TestWorkspaceExecEpisodeSelectsItsProfile'
CMD_S9D_LOCAL="$CMD_ADDED|$CMD_FENCES|$CMD_JOURNAL|$CMD_GIT"
CMD_S9D_RACE='TestOperatorFinishRouteWorker|TestOperatorFinishConfigReplacement'
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
  guard_hits() { printf '%s\n' "$1" | grep -cF '.ClaimOperatorAction(' || true; }
  claim=$(all '\.ClaimOperatorAction\(')
  extra=$'control.go:1:\treturn db.ClaimOperatorAction(ctx, cas, owner, now)'
  [ "$(guard_hits "$claim")" -eq 2 ] && [ "$(guard_hits "$claim"$'\n'"$extra")" -eq 3 ]
  echo CALLER_CENSUS_CONTROL added-claim=3 refused-against=2
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
census "$1" base 2 1
case "$3" in M1|M2) census "$2" "$3" 2 1;; *) exit 2;; esac
sites() (
  cd "$1"
  find cmd host -name '*.go' ! -name '*_test.go' -print0 | xargs -0 grep -EH     '(exec\.Command(Context)?|\.RunIsolated|context\.(Background|TODO|WithoutCancel))\(' | LC_ALL=C sort
)
diff -u <(sites "$1") <(sites "$2")
echo SUBPROCESS_AND_CONTEXT_SITE_CENSUS_OK
echo CALLER_CENSUS_OK
````

````bash outside.sh
#!/bin/bash
# TREE M1|M2 [focused|whole|chrome]. Controller executes outside the sandbox.
set -uo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P) || exit 2
. "$TOOLS/env.sh"
. "$TOOLS/selectors.sh"
T=$(cd "$1" && pwd -P) || exit 2
M=$2 mode=${3:-focused}
root=${REHEARSAL_ROOT:-"$S"}
L=$(mktemp -d "$root/outside.$M.$mode.XXXXXX") || exit 2
echo "OUTSIDE_LOGS $L"
python3 "$TOOLS/failset.py" --control || exit 2
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
skip_mutants() {
  python3 - "$TOOLS/mutants.json" "$TOOLS/named-mutant-ids.txt" <<'PY'
import json
from pathlib import Path
import sys
rows = json.loads(Path(sys.argv[1]).read_text())
for id in dict.fromkeys(Path(sys.argv[2]).read_text().split()):
    print(id, "SKIPPED", rows[id]["test"])
PY
}
if [ "${OUTSIDE_LOCAL:-0}" = 1 ]; then
  for tag in listeners race-budgets sticky-listeners subprocess-sites whole-verifygate chrome; do
    echo "LEG $tag SKIPPED (UNINFORMATIVE UNDER SANDBOX)"
  done
  skip_mutants
  echo ALL_DONE
  exit 77
fi
case "$mode" in
  focused)
    for leg in outside outside-race outside-1777; do
      run "$leg" /bin/bash "$TOOLS/accept.sh" "$T" "$M" "$leg"
    done
    run named /bin/bash "$TOOLS/named-mutants.sh" "$T" "$L/named" outside "$M"
    ;;
  whole)
    run whole /bin/bash "$TOOLS/gates.sh" "$T" "$L/gates" "$M" outside
    skip_mutants
    ;;
  chrome)
    run guard-sweep /bin/bash "$TOOLS/commit-guards.sh" "$T" "$L/guards"
    export WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
    export WORLD_REQUIRE_PASSKEY_BROWSER=1 AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
    mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" || exit 2
    cd "$T" || exit 2
    run chrome python3 "$TOOLS/bounded.py" 530 go test -v -count=1 -timeout 8m       ./cmd/ailang-worldd/ -run '^TestOperatorStartWalkthroughInChrome$'
    chrome_rc=$rc
    if grep -Fq -- '--- PASS: TestOperatorStartWalkthroughInChrome (' "$L/chrome.raw"; then chrome_rc=0; fi
    python3 "$TOOLS/failset.py" "$L/chrome.raw" "$chrome_rc" cmd/ailang-worldd > "$L/chrome.fails" || rc=$?
    grep -Fq -- '--- PASS: TestOperatorStartWalkthroughInChrome (' "$L/chrome.raw" || [ "$rc" -ne 0 ] || rc=2
    skip_mutants
    ;;
  *) exit 2;;
esac
echo ALL_DONE
exit "$rc"
````

````bash commit-guards.sh
#!/bin/bash
# TREE LOGDIR. Run outside against the controller's pristine TEXT fail sets.
set -uo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P) || exit 2
. "$TOOLS/env.sh"
T=$(cd "$1" && pwd -P) || exit 2
mkdir "$2" || exit 2
L=$(cd "$2" && pwd -P) || exit 2
cd "$T" || exit 2
export AILANG_STATE_DIR="$L/state" AILANG_CACHE_DIR="$L/cache"
mkdir "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR" || exit 2
{ date -u; uptime; } > "$L/guard.load"
rc=0 sites=0 result=0
python3 "$TOOLS/bounded.py" 2650 go test -count=1 -timeout 20m ./host/verifygate/ ./host/broker/ > "$L/guard.raw" 2>&1 || rc=$?
echo "$rc" > "$L/guard.rc"
python3 "$TOOLS/bounded.py" 530 go test -v -count=1 -timeout 8m -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' ./host/broker/ > "$L/sites.raw" 2>&1 || sites=$?
echo "$sites" > "$L/sites.rc"
python3 "$TOOLS/failset.py" "$L/sites.raw" "$sites" host/broker > "$L/sites.fails" || result=$?
if [ "$sites" != 0 ] || ! grep -Fq -- '--- PASS: TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential (' "$L/sites.raw"; then
  echo 'SUBPROCESS_SITE_NAMED_GATE_FAILED'
  [ "$result" != 0 ] || result=1
fi
guard_parser=0
python3 "$TOOLS/failset.py" "$L/guard.raw" "$rc" host/verifygate host/broker > "$L/guard.fails" || guard_parser=$?
if [ "$guard_parser" -eq 77 ]; then
  echo 'COMMIT_GUARDS UNINFORMATIVE UNDER SANDBOX'
  exit 77
elif [ "$guard_parser" -ne 0 ]; then
  echo 'COMMIT_GUARDS_INSTRUMENT_FAILURE'
  exit 2
fi
if ! grep -Fxq "ALL_DONE for $B" "$STATE/baseline.done"; then
  echo 'COMMIT_GUARDS SKIPPED baseline incomplete'
  exit 77
fi
for pkg in verifygate broker; do
  base="$STATE/base/$pkg-plain"
  python3 "$TOOLS/failset.py" "$base.raw" "$(cat "$base.rc")" "host/$pkg" > "$L/$pkg.base" || exit 2
done
cat "$L/verifygate.base" "$L/broker.base" | LC_ALL=C sort -u > "$L/base.fails"
LC_ALL=C sort -u "$L/guard.fails" > "$L/head.fails"
LC_ALL=C comm -13 "$L/base.fails" "$L/head.fails" > "$L/new.fails"
if [ -s "$L/new.fails" ]; then cat "$L/new.fails"; exit 1; fi
[ "$result" -eq 0 ] || exit "$result"
echo COMMIT_GUARDS_OK
````

````bash named-mutants.sh
#!/bin/bash
# TREE [LOGDIR] [outside|local] [M1|M2] [IDFILE]
set -uo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P) || exit 2
. "$TOOLS/env.sh"
mode=${3:-outside}
M=${4:-M1}
ids=${5:-"$TOOLS/named-mutant-ids.txt"}
root=${REHEARSAL_ROOT:-"$S"}
if [ -n "${2:-}" ]; then parent=$2; else parent=$(mktemp -d "$root/named.$M.XXXXXX")/results; fi
rc=0
python3 "$TOOLS/failset.py" --control || exit 2
if [ "$mode" = outside ]; then
  /bin/bash "$TOOLS/commit-guards.sh" "$1" "$parent.guards" || rc=$?
elif [ "$mode" = local ]; then
  echo 'GUARD_SWEEP SKIPPED outside listener/subprocess prerequisites'
else
  exit 2
fi
set -- "$1" named "$parent" "$ids" --stage "$M"
if [ "$mode" = local ]; then set -- "$@" --local; fi
python3 "$TOOLS/mutate.py" "$@" || exit "$?"
exit "$rc"
````

````bash whole-mutants.sh
#!/bin/bash
# TREE [IDFILE] [M1|M2]. Outside only. BASE must be recorded.
set -uo pipefail
TOOLS=$(cd "$(dirname "$0")" && pwd -P) || exit 2
. "$TOOLS/env.sh"
root=${REHEARSAL_ROOT:-"$S"}
L=$(mktemp -d "$root/whole-mutants.XXXXXX") || exit 2
M=${3:-M1}
rc=0
python3 "$TOOLS/failset.py" --control || exit 2
if [ "${OUTSIDE_LOCAL:-0}" = 1 ]; then
  echo 'whole mutation attempts SKIPPED; full package needs the real listeners'
else
  /bin/bash "$TOOLS/commit-guards.sh" "$1" "$L/guards" || rc=$?
fi
set -- "$1" whole "$L/results" "${2:-"$TOOLS/whole-mutant-ids.txt"}" --stage "$M"
if [ "${OUTSIDE_LOCAL:-0}" = 1 ]; then set -- "$@" --local; fi
python3 "$TOOLS/mutate.py" "$@" || exit "$?"
exit "$rc"
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
    if sys.argv[1:] == ["--control"]:
        return
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
"""Rebased S9c overlay mutator. Whole tests compare pristine TEXT fail sets."""
import argparse
import json
import os
from pathlib import Path
import re
import subprocess
import sys

p = argparse.ArgumentParser(description=__doc__)
p.add_argument("tree", type=Path)
p.add_argument("mode", choices=["named", "whole", "compile"])
p.add_argument("log", type=Path)
p.add_argument("ids", type=Path)
p.add_argument("--stage", choices=["M1","M2"], default="M1")
p.add_argument("--local", action="store_true")
a = p.parse_args()
tools = Path(__file__).resolve().parent
tree = a.tree.resolve()
a.log.mkdir()
all_rows = json.loads((tools/"mutants.json").read_text())
patterns = list(dict.fromkeys(a.ids.read_text().split()))
limits = {"cmd/ailang-worldd": ("45m", 2900), "host/daemon": ("45m", 2900),
          "host/surfacepage": ("25m", 1700)}
stage_to_base = {"cmd/ailang-worldd": "cmd-plain", "host/daemon": "daemon-plain", "host/surfacepage": "surfacepage-plain"}
base_dir = Path(os.environ["STATE"]) / "base"

def parsed(log, result, pkg, capture):
    proc = subprocess.run([sys.executable, str(tools/"failset.py"), str(log), str(result), pkg],
                          capture_output=True, text=True)
    capture.write_text(proc.stderr)
    print(proc.stderr.rstrip(), flush=True)
    return proc.returncode, set(proc.stdout.splitlines())

def overlays(id, row, spot):
    texts = {}
    for edit in row.get("edits",[id]):
        spec = all_rows[edit]
        path = tree/spec["path"]
        src = texts.get(path)
        if src is None:
            src = path.read_text()
        hits = list(re.finditer(re.escape(spec["old"]),src))
        if spec["occurrence"] >= len(hits):
            raise ValueError("anchor drift: "+edit)
        hit = hits[spec["occurrence"]]
        texts[path] = src[:hit.start()] + spec["new"] + src[hit.end():]
    replacements = {}
    for n, (path, body) in enumerate(texts.items()):
        dest = spot/(str(n)+"-"+path.name)
        dest.write_text(body)
        replacements[str(path)] = str(dest)
    overlay = spot/"overlay.json"
    overlay.write_text(json.dumps({"Replace":replacements}))
    return overlay

def trial(id, row):
    pkg, test = row["package"], row["test"]
    if row["stage"] != a.stage:
        return "SKIPPED", "other milestone"
    if a.local and a.mode!="compile" and row["access"] == "outside":
        return "SKIPPED", "PENDING-OUTSIDE, UNINFORMATIVE UNDER SANDBOX"
    spot = a.log/id
    spot.mkdir()
    try:
        overlay = overlays(id, row, spot)
    except ValueError as error:
        return "ERROR", str(error)
    if a.mode == "whole":
        name = stage_to_base[pkg]
        base = base_dir/(name+".raw")
        rc_file = base_dir/(name+".rc")
        if not base.exists() or not rc_file.exists():
            return "ERROR", "missing comparable pristine log"
        rc, base_fails = parsed(base,int(rc_file.read_text()),pkg,spot/"base.parser")
        if rc:
            return "ERROR", "unusable pristine log"
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
                                cwd=tree,env=env,stdout=fd,stderr=subprocess.STDOUT).returncode
    (spot/"rc").write_text(str(result)+"\n")
    rc, found = parsed(raw,result,pkg,spot/"parser")
    if rc == 77 and a.local:
        return "SKIPPED", "PENDING-OUTSIDE, UNINFORMATIVE UNDER SANDBOX"
    if rc:
        return "ERROR", "compile, timeout, denial or TEXT attribution failure"
    if a.mode == "compile":
        return ("SKIPPED", "compiled only") if result == 0 else ("ERROR", "compile failure")
    (spot/"base.fails").write_text("".join(x+"\n" for x in sorted(base_fails)))
    (spot/"head.fails").write_text("".join(x+"\n" for x in sorted(found)))
    new = found - base_fails
    output = raw.read_text(errors="replace")
    top = re.search(r"^--- FAIL: "+re.escape(test)+r" \(",output,re.M)
    sub = row.get("required_subtest")
    child = sub is None or re.search(r"^[ \t]+--- FAIL: "+re.escape(test+"/"+sub)+r" \(",output,re.M)
    if pkg+" "+test in new and top and child:
        return "KILLED", ""
    if new:
        return "ERROR", "wrong failure or missing target subtest: "+", ".join(sorted(new))
    if a.mode == "named" and not re.search(r"^--- PASS: "+re.escape(test)+r" \(",output,re.M):
        return "ERROR", "named test did not pass or execute"
    return "SURVIVED", ""

result = {}
for id in patterns:
    row = all_rows[id]
    result[id] = trial(id,row)
    # Retain a durable progress receipt during long runs without publishing
    # decisions that may depend on a later double.
    (a.log/"progress").write_text(f"{len(result)}/{len(patterns)} {id}\n")

failed = pending = False
with (a.log/"ledger").open("w",buffering=1) as ledger:
    for id in patterns:
        row = all_rows[id]
        status, detail = result[id]
        if status == "SURVIVED" and row.get("proof"):
            proof = row["proof"]
            sibling = all_rows[proof]["edits"][1]
            sibling_status = result.get(sibling,("",""))[0]
            # Some sibling singles already fail on their own separate row.
            # Double-mutant attribution still requires the target subtest.
            sibling_expected = "SURVIVED" if all_rows[sibling]["control"]=="live" else "KILLED"
            if result.get(proof,("",""))[0] == "KILLED" and sibling_status == sibling_expected:
                status, detail = "EQUIVALENT-PROVEN", "double "+proof+" failed "+row["required_subtest"]
        if status == "ERROR":
            status, detail = "SKIPPED", "INSTRUMENT FAILURE: "+detail
            failed = True
        elif row["control"] == "live" and status == "SURVIVED":
            detail = "MUST-LIVE control held"
        elif row["control"] == "live" and status == "KILLED":
            detail = "MUST-LIVE control failed"
            failed = True
        elif status == "SURVIVED":
            failed = True
        pending |= "PENDING-OUTSIDE" in detail
        if detail:
            print("DETAIL", id, detail, flush=True)
        line = f"{id} {status} {row['test']}"
        print(line,flush=True)
        ledger.write(line+"\n")
    ledger.write("ALL_DONE\n")
print("ALL_DONE")
sys.exit(1 if failed else 77 if pending else 0)
````

````python scope-assets.py
"""Check exact stage file scope and full bytes for assets, goldens and schemas."""
import hashlib
import os
from pathlib import Path
import sys

base, head, tools = map(Path, sys.argv[1:4])
stage = sys.argv[4]
changed = set((tools / ("changed-M1" if stage == "M1" else "changed-M2-cumulative")).read_text().splitlines())
plans = {
    "design_docs/planned/w-world-address-bar-verbs-v4s9d-sprint-plan.md",
    "design_docs/planned/sprint_w-world-address-bar-verbs-v4s9d.json",
}

def inventory(root):
    result = {}
    for dir, subs, files in os.walk(root):
        subs[:] = [s for s in subs if s != ".git"]
        for file in files:
            p = Path(dir) / file
            path = p.relative_to(root).as_posix()
            if path in plans or path == ".git":
                continue
            data = os.readlink(p).encode() if p.is_symlink() else p.read_bytes()
            result[path] = hashlib.sha256(data).hexdigest()
    return result

before, after = map(inventory, (base, head))
actual = {p for p in before.keys() | after.keys() if before.get(p) != after.get(p)}
if actual != changed:
    raise SystemExit("SCOPE_FAILURE " + repr(sorted(actual ^ changed)))
protected = {p for p in actual if p.endswith((".ail",".sql")) or ".golden" in p or "/testdata/" in p or "/app/" in p}
if protected:
    raise SystemExit("PINNED_ASSET_FAILURE " + repr(sorted(protected)))
if stage == "M1" and any(not p.endswith("_test.go") for p in actual):
    raise SystemExit("M1_PRODUCTION_DIFF")
print("SCOPE_ASSET_HASHES_OK", stage, len(actual))
print("PRODUCTION_TEST_HOOK_FIELDS 0")
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
    "line": 198,
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
    "control": "die",
    "required_subtest": "input/resumes"
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
    "control": "die",
    "required_subtest": "input/preset-allowlist"
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
    "control": "die",
    "required_subtest": "input/remote-allowlist"
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
    "control": "equivalent",
    "proof": "double-runner-elig-launch-version",
    "required_subtest": "input/launch-version"
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
    "control": "die",
    "required_subtest": "input/health-status"
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
    "control": "die",
    "required_subtest": "input/build"
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
    "control": "die",
    "required_subtest": "input/report-read"
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
    "control": "die",
    "required_subtest": "input/report-episode"
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
    "control": "die",
    "required_subtest": "input/report-intents"
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
    "control": "die",
    "required_subtest": "input/recorded-finish"
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
    "control": "die",
    "required_subtest": "input/report-start"
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
    "control": "equivalent",
    "proof": "double-runner-drift-preset-decode",
    "required_subtest": "input/preset-decode"
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
    "control": "die",
    "required_subtest": "input/workspace"
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
    "control": "die",
    "required_subtest": "input/world"
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
    "control": "die",
    "required_subtest": "input/store"
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
    "control": "die",
    "required_subtest": "input/profile"
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
    "control": "die",
    "required_subtest": "input/pinned-check"
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
    "control": "die",
    "required_subtest": "input/sessions"
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
    "control": "die",
    "required_subtest": "input/project"
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
    "control": "die",
    "required_subtest": "input/session-file"
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
    "control": "die",
    "required_subtest": "input/token-mode"
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
    "control": "die",
    "required_subtest": "input/commit-identity"
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
    "control": "die",
    "required_subtest": "claim/manifest-protected"
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
    "control": "die",
    "required_subtest": "claim/catalogue-fresh"
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
    "control": "equivalent",
    "proof": "double-runner-heartbeat-owner",
    "required_subtest": "claim/heartbeat-owner"
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
    "control": "die",
    "required_subtest": "claim/heartbeat-ready"
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
    "control": "die",
    "required_subtest": "claim/life-quiet"
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
    "control": "die",
    "required_subtest": "claim/life-pair"
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
    "control": "die",
    "required_subtest": "after-claim/clock"
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
    "control": "die",
    "required_subtest": "after-claim/inputs"
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
    "control": "die",
    "required_subtest": "after-claim/key"
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
    "control": "die",
    "required_subtest": "catalogue/approval-size"
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
    "control": "die",
    "required_subtest": "catalogue/aggregate-size"
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
    "control": "die",
    "required_subtest": "other-store"
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
    "control": "die",
    "required_subtest": "protected"
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
    "control": "die",
    "required_subtest": "existing-operation/other-generation"
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
    "control": "die",
    "required_subtest": "confirm/without-surface-option"
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
    "control": "die",
    "required_subtest": "status/start-only-runner"
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
  },
  "e01": {
    "path": "host/daemon/surface_operator_confirm.go",
    "old": "a.Kind == \"finish\" && a.Options.ConfirmOnSurface",
    "new": "a.Kind == \"finish\"",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "confirm/without-surface-option",
    "line": 96
  },
  "e04": {
    "path": "host/daemon/surface_operator.go",
    "old": "len(entry.Protected) != 0 || len(entry.Unattributed) != 0",
    "new": "len(entry.Protected) != 0",
    "test": "TestOperatorFinishAdmission",
    "package": "host/daemon",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "unattributed",
    "line": 56
  },
  "e06": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "len(start.Resumes) != 0 || ",
    "new": "",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "input/resumes",
    "line": 41
  },
  "e07": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": " || rep.Uncommitted != 0",
    "new": "",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "input/report-intents",
    "line": 58
  },
  "e08": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": " ||\n\t\thealth.Build.ExecutableSHA256 != c.Executables[e.paths.Binary]",
    "new": "",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "input/build",
    "line": 41
  },
  "e09": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "!slices.Contains(c.FinishRemotes, start.Remote.URL)",
    "new": "false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "input/remote-allowlist",
    "line": 42
  },
  "e12": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "len(m.Protected) != 0 || len(m.Unattributed) != 0 {",
    "new": "len(m.Protected) != 0 {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "claim/manifest-unattributed",
    "line": 278
  },
  "e13": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "live.CatalogueGeneration != generation || ",
    "new": "",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "claim/catalogue-fresh",
    "line": 293
  },
  "e14": {
    "path": "cmd/ailang-worldd/operator_finish_tree.go",
    "old": "config != spec.Snapshot.Config",
    "new": "config != config",
    "test": "TestOperatorFinishConfigReplacement",
    "package": "cmd/ailang-worldd",
    "access": "local",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "valid-config-with-different-digest",
    "line": 348
  },
  "e16": {
    "path": "host/daemon/surface_operator_read.go",
    "old": "j.OpenCheck || j.CommitPhase == \"reserved\"",
    "new": "j.OpenCheck",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "line": 141
  },
  "e17": {
    "path": "cmd/ailang-worldd/operator_runner.go",
    "old": "if h.Readiness == \"both\" {",
    "new": "if true {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "claim/start-ready-loop",
    "line": 157
  },
  "e20": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "launch.Key.Version != 4",
    "new": "false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "equivalent",
    "stage": "M1",
    "proof": "double-e20",
    "required_subtest": "input/launch-version",
    "line": 46
  },
  "e22": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "verifyAt = min(verifyAt, row.ExpiresAt-1)",
    "new": "",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "line": 203
  },
  "e23": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if !finishReviewInputsMatch(in.manifest, m) {",
    "new": "if false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "after-claim/review",
    "line": 270
  },
  "e24": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "ctx.Err() == nil && e.checkStartPair",
    "new": "e.checkStartPair",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "claim/cancelled",
    "line": 285
  },
  "e25": {
    "path": "cmd/ailang-worldd/operator_catalogue.go",
    "old": "4 * time.Second",
    "new": "5 * time.Second",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "live",
    "stage": "M1",
    "line": 40
  },
  "e21b": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "live && birth == launch.Receipt.Birth",
    "new": "live && birth == birth",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "required_subtest": "input/reused-pid",
    "line": 49
  },
  "m1-sibling-completed-launch": {
    "path": "host/store/operator_finish.go",
    "old": "WHERE a.episode_id=? AND a.kind='start' AND a.state='succeeded' AND s.kind='launch' AND s.phase='observed' AND s.version=4",
    "new": "WHERE a.episode_id=? AND a.kind='start' AND a.state='succeeded' AND s.kind='launch' AND s.phase='observed' AND s.version IN (3,4)",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "live",
    "stage": "M1",
    "required_subtest": "input/launch-version",
    "line": 126
  },
  "m1-sibling-finish-preset-name": {
    "path": "cmd/ailang-worldd/episode_finish.go",
    "old": "if p.Schema != \"world/preset/v1\" || p.Name != name {",
    "new": "if p.Schema != \"world/preset/v1\" {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "live",
    "stage": "M1",
    "required_subtest": "input/preset-decode",
    "line": 569
  },
  "m1-sibling-cas-heartbeat-owner": {
    "path": "host/store/operator_cas.go",
    "old": "generation = ? AND owner = ? AND seen_at > ?",
    "new": "? IS NOT NULL AND ? IS NOT NULL AND seen_at > ?",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "live",
    "stage": "M1",
    "required_subtest": "claim/heartbeat-owner",
    "line": 140
  },
  "double-runner-elig-launch-version": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "err != nil || launch.Key.Version != 4",
    "new": "err != nil || launch.Key.Version != 4 && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 46,
    "stage": "M1",
    "control": "die",
    "edits": [
      "runner-elig-launch-version",
      "m1-sibling-completed-launch"
    ],
    "required_subtest": "input/launch-version"
  },
  "double-e20": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "launch.Key.Version != 4",
    "new": "false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M1",
    "edits": [
      "e20",
      "m1-sibling-completed-launch"
    ],
    "required_subtest": "input/launch-version",
    "line": 46
  },
  "double-runner-drift-preset-decode": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "strictPresetJSON(raw, &preset) != nil || preset.Name != start.Preset.Name",
    "new": "(strictPresetJSON(raw, &preset) != nil || preset.Name != start.Preset.Name) && false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 70,
    "stage": "M1",
    "control": "die",
    "edits": [
      "runner-drift-preset-decode",
      "m1-sibling-finish-preset-name"
    ],
    "required_subtest": "input/preset-decode"
  },
  "double-runner-heartbeat-owner": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "err != nil || live.Owner != h.Owner || live.CatalogueGeneration",
    "new": "err != nil || live.Owner != h.Owner && false || live.CatalogueGeneration",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "line": 293,
    "stage": "M1",
    "control": "die",
    "edits": [
      "runner-heartbeat-owner",
      "m1-sibling-cas-heartbeat-owner"
    ],
    "required_subtest": "claim/heartbeat-owner"
  },
  "s9d-cli-help": {
    "path": "cmd/ailang-worldd/operator_service.go",
    "old": " | finish-status ACTION-ID | ",
    "new": " | ",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/help",
    "line": 24
  },
  "s9d-cli-verb": {
    "path": "cmd/ailang-worldd/operator_service.go",
    "old": ", \"launch-status\", \"finish-status\":",
    "new": ", \"launch-status\":",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "line": 294
  },
  "s9d-cli-id-grammar": {
    "path": "cmd/ailang-worldd/operator_service.go",
    "old": "len(args) != 2 || !verbs.ValidOperatorActionID(args[1])",
    "new": "len(args) != 2",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/invalid-id",
    "line": 295
  },
  "s9d-cli-arity": {
    "path": "cmd/ailang-worldd/operator_service.go",
    "old": "len(args) != 2 || !verbs.ValidOperatorActionID(args[1])",
    "new": "len(args) < 2 || !verbs.ValidOperatorActionID(args[1])",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/extra-argument",
    "line": 295
  },
  "s9d-cli-branch": {
    "path": "cmd/ailang-worldd/operator_service.go",
    "old": "if verb == \"finish-status\" {",
    "new": "if false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "line": 317
  },
  "s9d-cli-dispatch-error": {
    "path": "cmd/ailang-worldd/operator_service.go",
    "old": "if err := e.printFinishStatus(args[1], out); err != nil {",
    "new": "if err := e.printFinishStatus(args[1], out); err != nil && false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/unknown-action",
    "line": 318
  },
  "s9d-cli-installed-remedy": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "\tc, _, err := e.installedRunner()\n\tif err != nil {\n\t\treturn err",
    "new": "\tc, _, err := e.installedRunner()\n\tif err != nil {\n\t\treturn store.ErrOperatorMissing",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/uninstalled",
    "line": 357
  },
  "s9d-cli-read-store-remedy": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "store.OpenCredentialsReadOnly(runnerDB(c.Credentials))\n\tif err != nil {\n\t\treturn errOperatorFinishEligibility",
    "new": "store.OpenCredentialsReadOnly(runnerDB(c.Credentials))\n\tif err != nil {\n\t\treturn store.ErrOperatorMissing",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/unreadable-credentials",
    "line": 152
  },
  "s9d-cli-unknown-action": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if err != nil || row.Kind != \"finish\" {",
    "new": "if err != nil { return nil }\n\tif row.Kind != \"finish\" {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/unknown-action",
    "line": 369
  },
  "s9d-cli-non-finish": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "row.Kind != \"finish\"",
    "new": "false",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/non-finish",
    "line": 369
  },
  "s9d-cli-missing-id-remedy": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "row.Kind != \"finish\" {\n\t\treturn store.ErrOperatorMissing",
    "new": "row.Kind != \"finish\" {\n\t\treturn errOperatorFinishEligibility",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/unknown-action",
    "line": 369
  },
  "s9d-cli-journal-remedy": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))\n\tif err != nil {\n\t\treturn errOperatorFinishUncertain",
    "new": "ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))\n\tif err != nil {\n\t\treturn store.ErrOperatorMissing",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/unreadable-journal",
    "line": 372
  },
  "s9d-cli-redact-stderr": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))\n\tif err != nil {\n\t\treturn errOperatorFinishUncertain",
    "new": "ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))\n\tif err != nil {\n\t\treturn fmt.Errorf(\"%w: %s\", errOperatorFinishUncertain, row.Artifact)",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/unreadable-journal",
    "line": 372
  },
  "s9d-cli-redact-stdout": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "check.Name, check.Result)",
    "new": "check.Name+\" \"+check.Ref, check.Result)",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "line": 378
  },
  "s9d-cli-open-check": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if j.OpenCheck || j.CommitPhase == \"reserved\" {",
    "new": "if j.CommitPhase == \"reserved\" {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "required_subtest": "cli/open-check",
    "line": 380
  },
  "s9d-cli-reserved": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if j.OpenCheck || j.CommitPhase == \"reserved\" {",
    "new": "if j.OpenCheck {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "line": 380
  },
  "s9d-cli-done-state": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "if j.Commit.State == store.OperatorLocalCommitDone {",
    "new": "if false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "line": 383
  },
  "s9d-cli-failed-state": {
    "path": "cmd/ailang-worldd/operator_finish_runner.go",
    "old": "} else if j.Commit.State == \"commit-failed\" {",
    "new": "} else if false {",
    "test": "TestOperatorFinishRouteWorker",
    "package": "cmd/ailang-worldd",
    "access": "outside",
    "occurrence": 0,
    "control": "die",
    "stage": "M2",
    "line": 385
  }
}
````

````text changed-M1
cmd/ailang-worldd/operator_finish_route_test.go
cmd/ailang-worldd/operator_finish_tree_test.go
host/daemon/surface_operator_finish_test.go
````

````text changed-M2-cumulative
cmd/ailang-worldd/operator_finish_route_test.go
cmd/ailang-worldd/operator_finish_runner.go
cmd/ailang-worldd/operator_finish_tree_test.go
cmd/ailang-worldd/operator_service.go
docs/QUICKSTART.md
host/daemon/surface_operator_finish_test.go
````

````text status-M1
M	cmd/ailang-worldd/operator_finish_route_test.go
M	cmd/ailang-worldd/operator_finish_tree_test.go
M	host/daemon/surface_operator_finish_test.go
````

````text status-M2-cumulative
M	cmd/ailang-worldd/operator_finish_route_test.go
M	cmd/ailang-worldd/operator_finish_runner.go
M	cmd/ailang-worldd/operator_finish_tree_test.go
M	cmd/ailang-worldd/operator_service.go
M	docs/QUICKSTART.md
M	host/daemon/surface_operator_finish_test.go
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
e01
e04
e06
e07
e08
e09
e12
e13
e14
e16
e17
e20
e22
e23
e24
e25
e21b
m1-sibling-completed-launch
m1-sibling-finish-preset-name
m1-sibling-cas-heartbeat-owner
double-runner-elig-launch-version
double-e20
double-runner-drift-preset-decode
double-runner-heartbeat-owner
s9d-cli-help
s9d-cli-verb
s9d-cli-id-grammar
s9d-cli-arity
s9d-cli-branch
s9d-cli-dispatch-error
s9d-cli-installed-remedy
s9d-cli-read-store-remedy
s9d-cli-unknown-action
s9d-cli-non-finish
s9d-cli-missing-id-remedy
s9d-cli-journal-remedy
s9d-cli-redact-stderr
s9d-cli-redact-stdout
s9d-cli-open-check
s9d-cli-reserved
s9d-cli-done-state
s9d-cli-failed-state
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

Run the `bash apply` fence verbatim. The executor starts at the pinned base with the two plan artifacts present. The block applies the incremental fences sequentially, checks exact dry runs and runs `accept.sh` in sandbox-safe fast mode. The success signal is `ACCEPT_OK <M> <mode>`, followed by `APPLY_ACCEPT_SNAPSHOT_OK <M>`. A complete run ends with `FINAL_APPLY_RC 0`. The rehearsal mode runs the same block on a fresh archive and prints the recursive diff markers at each milestone.

The snapshots at `$HOME/.ailang/state/world-iter451/snap/M1` and `snap/M2` are cumulative relative to the full base SHA. M1 has 3 status entries, 0 deletions and 3 manifest entries. M2 has 6 status entries, 0 deletions and 6 manifest entries. `snapshot.sh` includes relative paths and excludes both plan artifacts. The two plan files stay untracked through executor patch/snapshot work in this run. The controller records them before `commit_ms.sh` so its guard can distinguish planned artifacts from untracked stray design files. The controller owns all Git writes.

Use the full `commit-M1` or `commit-M2` fence as the outside argument to `$STATE/commit_ms.sh <M> <msgfile> <test-cmd…>`. The commit guard contains the explicit `go test -count=1 -timeout 20m ./host/verifygate/ ./host/broker/` package leg and the exact subprocess test. It compares against `$STATE/base/{verifygate,broker}-plain`. Cmd and daemon use named selectors at commit time. Apply M2 only after M1 passes. If M2 is deferred, the controller sets `PLAN_LAST_MILESTONE=M1` and compares the M1 snapshot.

Verbatim Bash 3.2 rehearsal in `$S/rt-apply` printed:

````text
ACCEPT_OK M1 fast
ROUNDTRIP_DIFF_OK M1 cmd
ROUNDTRIP_DIFF_OK M1 host
ROUNDTRIP_DIFF_OK M1 docs
ROUNDTRIP_ABSENT M1 app
ROUNDTRIP_DIFF_OK M1 unfiltered
APPLY_ACCEPT_SNAPSHOT_OK M1
ACCEPT_OK M2 fast
ROUNDTRIP_DIFF_OK M2 cmd
ROUNDTRIP_DIFF_OK M2 host
ROUNDTRIP_DIFF_OK M2 docs
ROUNDTRIP_ABSENT M2 app
ROUNDTRIP_DIFF_OK M2 unfiltered
APPLY_ACCEPT_SNAPSHOT_OK M2
FINAL_APPLY_RC 0
````

Both cumulative snapshot manifests passed SHA-256 verification.

````bash commit-M1
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
S="$HOME/.ailang/state/world-iter451/plan-scratch"
/bin/bash "$S/accept.sh" "$PWD" M1 outside &&
  /bin/bash "$S/commit-guards.sh" "$PWD" "$S/commit-guards-M1-$(date +%s)"
````

````bash commit-M2
export AILANG_BIN="$HOME/.pinned-ailang/ailang" WORLD_REQUIRE_SSH_KEYGEN=1
export WORLD_EXEC_NODE=/opt/homebrew/bin/node WORLD_TOOL_AILANG_BIN="$HOME/.pinned-ailang-tools/v0.54.15/ailang"
S="$HOME/.ailang/state/world-iter451/plan-scratch"
/bin/bash "$S/accept.sh" "$PWD" M2 outside &&
  /bin/bash "$S/commit-guards.sh" "$PWD" "$S/commit-guards-M2-$(date +%s)"
````

````bash apply
(
set -euo pipefail
export GIT_OPTIONAL_LOCKS=0
B=db8601b0fc8d7ee70eba7d08692f8686f6b5e707
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter451-row243-v4s9d
STATE="$HOME/.ailang/state/world-iter451"
S="$STATE/plan-scratch"
PLAN="$WT/design_docs/planned/w-world-address-bar-verbs-v4s9d-sprint-plan.md"
MODE=${PLAN_APPLY_MODE:-executor}
LAST=${PLAN_LAST_MILESTONE:-M2}
case "$LAST" in M1|M2) ;; *) exit 2;; esac
case "$MODE" in
  rehearsal)
    R="$S/rt-apply"
    mkdir "$R"
    T="$R/export"
    TOOLS="$R/tools"
    mkdir "$T" "$TOOLS"
    git -C "$WT" archive "$B" | tar -x -C "$T"
    export REHEARSAL_ROOT="$R"
    ;;
  executor)
    T="$WT"
    TOOLS="$STATE/exec-tools"
    mkdir -p "$TOOLS"
    test "$(git -C "$T" rev-parse HEAD)" = "$B"
    test -z "$(git -C "$T" status --porcelain -- cmd host docs app)"
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
for name in bounded.py failset.py verify-scratch.py mutate.py scope-assets.py; do
  fence "python $name" "$TOOLS/$name"
done
fence 'json mutants.json' "$TOOLS/mutants.json"
for name in changed-M1 changed-M2-cumulative status-M1 status-M2-cumulative named-mutant-ids.txt whole-mutant-ids.txt; do
  fence "text $name" "$TOOLS/$name"
done
for M in M1 M2; do
  fence "diff $M" "$TOOLS/$M.patch"
  (cd "$T" && /usr/bin/patch -p1 -E --dry-run < "$TOOLS/$M.patch") > "$TOOLS/$M.dry-run" 2>&1
  if grep -Eiq 'offset|fuzz|FAILED|reject' "$TOOLS/$M.dry-run"; then
    cat "$TOOLS/$M.dry-run"
    exit 1
  fi
  (cd "$T" && /usr/bin/patch -p1 -E < "$TOOLS/$M.patch")
  /bin/bash "$TOOLS/accept.sh" "$T" "$M" fast
  if [ "$MODE" = rehearsal ]; then
    status="$TOOLS/status-$M"
    [ "$M" = M1 ] || status="$TOOLS/status-M2-cumulative"
    SNAP_PARENT="$R/snap" SNAP_STATUS_LIST="$status" /bin/bash "$TOOLS/snapshot.sh" "$T" "$M"
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
    /bin/bash "$TOOLS/snapshot.sh" "$T" "$M"
  fi
  echo "APPLY_ACCEPT_SNAPSHOT_OK $M"
  if [ "$M" = "$LAST" ]; then break; fi
done
echo "FINAL_APPLY_RC 0"
)
````

## Appendix P. Compiled patches

### M1 — Behavioural survivor tests

````diff M1
diff --git a/cmd/ailang-worldd/operator_finish_route_test.go b/cmd/ailang-worldd/operator_finish_route_test.go
index 12df2b4..5e40eae 100644
--- a/cmd/ailang-worldd/operator_finish_route_test.go
+++ b/cmd/ailang-worldd/operator_finish_route_test.go
@@ -2,6 +2,8 @@ package main
 
 import (
 	"bytes"
+	"context"
+	"crypto/sha256"
 	"encoding/hex"
 	"encoding/json"
 	"errors"
@@ -14,11 +16,14 @@ import (
 	"path/filepath"
 	"reflect"
 	"strings"
+	"sync"
+	"syscall"
 	"testing"
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/checkrun"
 	"github.com/sunholo-data/ailang-world/host/daemon"
 	"github.com/sunholo-data/ailang-world/host/hashref"
 	"github.com/sunholo-data/ailang-world/host/store"
@@ -43,11 +48,13 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 	presetGit(t, repo, "clone", "-q", "--local", "--bare", "--", repo.repo, bare)
 	remoteHead := presetGit(t, repo, "--git-dir="+bare, "rev-parse", "HEAD")
 	trace := filepath.Join(f.home, "unexpected-remote-write")
+	effectTrace := filepath.Join(f.home, "git-effect-trace")
 	tools := f.c.ToolDirs[0]
 	f.c.FinishRemotes = []string{fixtureFinishRemote}
 	writeFinishTestTool(t, &f.c, tools, "gh", "case \"$1 $2\" in 'auth status'|'label list') ;; 'repo view') printf 'WRITE\\n';; *) touch "+shellQuote(trace)+"; exit 96;; esac\n")
 	writeFinishTestTool(t, &f.c, tools, "git", "for arg; do case \"$arg\" in push|fetch|clone|pull|send-pack|fetch-pack|receive-pack|upload-pack|remote|submodule) touch "+shellQuote(trace)+"; exit 96;; esac; done\n"+
 		"case \" $* \" in *' ls-remote '*) exec "+shellQuote(repo.git)+" --git-dir=/dev/null ls-remote --symref -- "+shellQuote(bare)+" HEAD;; esac\n"+
+		"for arg; do case \"$arg\" in commit-tree|update-ref|read-tree) printf '%s\\n' \"$arg\" >> "+shellQuote(effectTrace)+";; esac; done\n"+
 		"exec "+shellQuote(repo.git)+" \"$@\"\n")
 	paths := f.env.runnerPaths()
 	config, err := json.Marshal(f.c)
@@ -106,7 +113,7 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 	h := hashref.SumSHA256([]byte("selected world")).String()
 	checks, outcomes := 0, 0
 	activeID, statusHref := "", ""
-	reportDown := false
+	reportDown, reportMalformed := false, false
 	original := finishClient
 	t.Cleanup(func() { finishClient = original })
 	finishClient = func(string) *client {
@@ -117,6 +124,11 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 					w.WriteHeader(http.StatusServiceUnavailable)
 					return
 				}
+				if reportMalformed {
+					raw, _ := json.Marshal(rep)
+					fmt.Fprintf(w, `%s,"exec_cache":true}`, raw[:len(raw)-1])
+					return
+				}
 				_ = json.NewEncoder(w).Encode(rep)
 			case "/mcp/":
 				if body := request("GET", statusHref, nil, nil); !strings.Contains(body, "Finish outcome uncertain") {
@@ -189,8 +201,8 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 		t.Fatalf("missing review: %+v", review)
 	}
 	var status struct {
-		ID, Href, StatusHref, State, LocalCommit string
-		Checks                                   []struct{ Name, Result string }
+		ID, Href, StatusHref, State, Kind, Reason, LocalCommit string
+		Checks                                                 []struct{ Name, Result string }
 	}
 	request("POST", "/surface/operator-actions", map[string]any{
 		"argv": review.Argv, "operationId": strings.Repeat("4", 32), "catalogueGeneration": review.CatalogueGeneration, "task": "",
@@ -212,6 +224,11 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 	artifact := key.Sign(approval.OperatorActionNamespace, message)
 	request("POST", status.Href+"/confirm", map[string]any{"v": 1, "digest": hex.EncodeToString(row.Digest), "sshsig": string(artifact)}, &status)
 
+	row, err = f.db.OperatorAction(ctx, status.ID, f.now)
+	if err != nil || row.State != "confirmed" || row.Owner != "" {
+		t.Fatalf("signed route did not leave an unclaimed row: %s %v", row.State, err)
+	}
+
 	activeID, statusHref = row.ID, status.StatusHref
 	run := func() error {
 		turn, stop := f.env.operatorContext(activationTestBudget)
@@ -312,6 +329,12 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 	refuse("expired session", errOperatorFinishEligibility)
 	f.exec(t, "UPDATE session_credentials SET expires_at=? WHERE episode_id=?", f.self.ExpiresAt, start.Episode)
 
+	probe := finishRechecks{f: f, start: start, report: &rep, row: row, sign: func(msg []byte) []byte { return key.Sign(approval.OperatorActionNamespace, msg) },
+		run: run, bare: bare, trace: trace, change: change, effectTrace: effectTrace, startWire: startWire,
+		checks: &checks, outcomes: &outcomes, malformed: &reportMalformed}
+	probe.inputs(t)
+	probe.catalogue(t)
+	probe.claim(t, raw)
 	f.env.operatorCatalogue = func(runnerConfig) []byte { return raw }
 	fmt.Fprint(&runnerLog, f.tick(t))
 	row, err = f.db.OperatorAction(ctx, status.ID, f.now)
@@ -358,4 +381,588 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 	if _, err := os.Lstat(trace); !errors.Is(err, os.ErrNotExist) {
 		t.Fatalf("remote write call escaped the worker: %v", err)
 	}
+	probe.recovery(t, row, j, func(reason, sha string) {
+		t.Helper()
+		status.Reason, status.LocalCommit, status.Checks = "", "", nil
+		request("GET", status.StatusHref, nil, &status)
+		if status.Reason != reason || status.LocalCommit != sha || status.State != row.State ||
+			status.ID != row.ID || status.Kind != "finish" || len(status.Checks) != 1 ||
+			status.Checks[0].Name != start.Checks[0].Name || status.Checks[0].Result != "passed" {
+			t.Fatalf("finish receipt status: %+v, want %q and SHA %q", status, reason, sha)
+		}
+		if card := request("GET", status.Href, nil, nil); !strings.Contains(card, reason) {
+			t.Fatalf("finish card lost the projected reason %q", reason)
+		}
+	})
+}
+
+type finishRechecks struct {
+	f                                *operatorLaunchRig
+	start                            *store.EpisodeStart
+	report                           *daemon.EpisodeReport
+	row                              store.OperatorActionRow
+	sign                             func([]byte) []byte
+	run                              func() error
+	bare, trace, change, effectTrace string
+	startWire                        []byte
+	checks, outcomes                 *int
+	malformed                        *bool
+}
+
+// Observe the entire credentials store (including action generations, passkey
+// counters, sessions and reservations), world store and both Git refs.
+func (p finishRechecks) state(t *testing.T) [8]string {
+	t.Helper()
+	f := p.f
+	if _, err := os.Lstat(p.trace); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("unexpected dispatch or remote write: %v", err)
+	}
+	bareHead := readStartFixtureFile(t, filepath.Join(p.bare, "HEAD"))
+	bareRefs := serviceFileHash(filepath.Join(p.bare, "packed-refs"))
+	if ref, symbolic := strings.CutPrefix(strings.TrimSpace(string(bareHead)), "ref: "); symbolic {
+		bareRefs += serviceFileHash(filepath.Join(p.bare, ref))
+	}
+	return [8]string{serviceFileHash(f.creds), serviceFileHash(f.c.World),
+		serviceFileHash(filepath.Join(f.f.r.f.repo, ".git", "refs", "heads", p.start.Branch)),
+		serviceHash(bareHead), bareRefs,
+		serviceFileHash(p.effectTrace), serviceFileHash(filepath.Join(f.bin, "ran")),
+		fmt.Sprintf("%d/%d", *p.checks, *p.outcomes)}
+}
+
+func (p finishRechecks) unchanged(t *testing.T, want error, run func() error) {
+	t.Helper()
+	before := p.state(t)
+	got := run()
+	after := p.state(t)
+	if !errors.Is(got, want) || before != after {
+		t.Fatalf("want %v without changing row, journal, refs or counters, got %v", want, got)
+	}
+	switch want {
+	case errOperatorFinishEligibility, errOperatorFinishReview, errOperatorFinishFence:
+		if finishRefusal(got) != want.Error() {
+			t.Fatalf("runner lost the refusal remedy for %v", want)
+		}
+	}
+}
+
+func (p finishRechecks) exec(t *testing.T, query string, args ...any) {
+	t.Helper()
+	ctx, stop := p.f.env.operatorContext(activationTestBudget)
+	defer stop()
+	if _, err := p.f.worktreeStartRig.db.ExecContext(ctx, query, args...); err != nil {
+		t.Fatal(err)
+	}
+}
+
+func (p finishRechecks) setStart(t *testing.T, start *store.EpisodeStart) {
+	t.Helper()
+	b, err := json.Marshal(start)
+	if err != nil {
+		t.Fatal(err)
+	}
+	p.exec(t, "UPDATE episode_starts SET start_json=? WHERE episode=?", string(b), start.Episode)
+	p.report.Start = start
+}
+
+// Reuse the genuinely route-confirmed row and the verified episode flock.
+// Each row can reach the exact claim preparation without repeating GitHub
+// readiness, which the surrounding route turn exercises with fixture tools.
+func (p finishRechecks) preclaim(t *testing.T, ctx context.Context) error {
+	t.Helper()
+	f := p.f
+	held := f.borrow(t)
+	defer held.Close()
+	fd, err := held.Verify(f.creds, p.start.Episode)
+	if err != nil {
+		return err
+	}
+	row, err := f.db.OperatorAction(ctx, p.row.ID, f.now)
+	if err != nil {
+		return err
+	}
+	a, m, generation, err := finishManifestForRow(f.c, row)
+	if err != nil {
+		return err
+	}
+	_, err = f.env.finishWork(ctx, f.c, f.h, f.src, f.db, fd, row, operatorPreparationCAS(row), a, m, generation, false, f.now)
+	return err
+}
+
+func (p finishRechecks) inputs(t *testing.T) {
+	t.Helper()
+	f := p.f
+	originalReport, health, config := *p.report, f.health, f.c
+	preset := f.f.file
+	presetWire := readStartFixtureFile(t, preset)
+	token := filepath.Join(f.f.r.sessions, p.start.Episode+".session")
+	sessions, err := f.db.ListSessions(activationTestContext(t), p.start.Episode, 2)
+	if err != nil || len(sessions) != 1 {
+		t.Fatalf("original session: rows=%d err=%v", len(sessions), err)
+	}
+	launch := f.receipt(t).Receipt
+	launchWire, _ := json.Marshal(launch)
+	alterPreset := func(t *testing.T, edit func(*launchPreset, *presetDaemon)) {
+		t.Helper()
+		var l launchPreset
+		var d presetDaemon
+		if err := json.Unmarshal(presetWire, &l); err != nil {
+			t.Fatal(err)
+		}
+		if err := json.Unmarshal(l.Daemon, &d); err != nil {
+			t.Fatal(err)
+		}
+		edit(&l, &d)
+		l.Daemon, _ = json.Marshal(d)
+		b, err := json.Marshal(l)
+		if err != nil {
+			t.Fatal(err)
+		}
+		seedStartFixtureFile(t, preset, b)
+		changed := *p.start
+		changed.Preset.Bytes, changed.Preset.SHA256 = string(b), serviceHash(b)
+		p.setStart(t, &changed)
+	}
+	chmod := func(t *testing.T, path string, mode os.FileMode) {
+		t.Helper()
+		if err := os.Chmod(path, mode); err != nil {
+			t.Fatal(err)
+		}
+		t.Cleanup(func() {
+			want := os.FileMode(0600)
+			if path == f.f.r.sessions {
+				want = 0700
+			}
+			if err := os.Chmod(path, want); err != nil {
+				t.Error(err)
+			}
+		})
+	}
+	input := func(t *testing.T) error {
+		bound := activationTestContext(t)
+		_, err := f.env.finishInputs(bound, f.c, f.src, f.db, p.start.Episode, nil)
+		return err
+	}
+	p.unchanged(t, nil, func() error { return input(t) })
+	for _, tc := range []struct {
+		name string
+		want error
+		edit func(*testing.T)
+	}{
+		{"resumes", errOperatorFinishEligibility, func(t *testing.T) {
+			changed := *p.start
+			changed.Resumes = []store.EpisodeResume{{CredentialID: sessions[0].CredentialID}}
+			p.setStart(t, &changed)
+		}},
+		{"preset-allowlist", errOperatorFinishEligibility, func(*testing.T) { f.c.PresetNames = nil }},
+		{"remote-allowlist", errOperatorFinishEligibility, func(*testing.T) { f.c.FinishRemotes = nil }},
+		{"launch-version", errOperatorFinishEligibility, func(t *testing.T) {
+			old := launch
+			old.Status = "launched"
+			b, _ := json.Marshal(old)
+			p.exec(t, "UPDATE operator_action_steps SET version=3, receipt=? WHERE kind='launch'", b)
+		}},
+		{"health-status", errOperatorFinishEligibility, func(*testing.T) { f.health.Status = "starting" }},
+		{"build", errOperatorFinishEligibility, func(*testing.T) { f.health.Build.ExecutableSHA256 = "sha256:other-fixture" }},
+		{"report-read", errOperatorFinishEligibility, func(*testing.T) { *p.malformed = true }},
+		{"report-episode", errOperatorFinishEligibility, func(*testing.T) { p.report.Episode = "other-episode" }},
+		{"report-intents", errOperatorFinishEligibility, func(*testing.T) { p.report.Uncommitted = 1 }},
+		{"recorded-finish", errOperatorFinishEligibility, func(*testing.T) { p.report.Finish = &daemon.EpisodeFinish{} }},
+		{"report-start", errOperatorFinishEligibility, func(*testing.T) {
+			copy := *p.start
+			copy.Model += "-other"
+			p.report.Start = &copy
+		}},
+		{"preset-decode", errOperatorFinishReview, func(t *testing.T) {
+			alterPreset(t, func(l *launchPreset, _ *presetDaemon) { l.Name = "other-preset" })
+		}},
+		{"workspace", errOperatorFinishReview, func(t *testing.T) {
+			alterPreset(t, func(_ *launchPreset, d *presetDaemon) { d.WorkspaceRoot = t.TempDir() })
+		}},
+		{"world", errOperatorFinishReview, func(t *testing.T) {
+			other := filepath.Join(t.TempDir(), "other-world.db")
+			seedStartFixtureFile(t, other, readStartFixtureFile(t, f.c.World))
+			alterPreset(t, func(_ *launchPreset, d *presetDaemon) { d.DB = other })
+		}},
+		{"store", errOperatorFinishReview, func(t *testing.T) {
+			other := filepath.Join(t.TempDir(), "other-credentials.db")
+			seedStartFixtureFile(t, other, readStartFixtureFile(t, f.creds))
+			alterPreset(t, func(_ *launchPreset, d *presetDaemon) { d.CredentialsDB = other })
+		}},
+		{"profile", errOperatorFinishReview, func(t *testing.T) {
+			p.exec(t, "UPDATE session_credentials SET profile_digest=? WHERE episode_id=?", "sha256:other", p.start.Episode)
+		}},
+		{"pinned-check", errOperatorFinishReview, func(t *testing.T) {
+			changed := *p.start
+			changed.Checks = append([]store.StartCheck(nil), p.start.Checks...)
+			changed.Checks[0].Name = "different-pin"
+			p.setStart(t, &changed)
+		}},
+		{"sessions", errOperatorFinishEligibility, func(t *testing.T) {
+			p.exec(t, `INSERT INTO session_credentials SELECT 'zzz-fixture', episode_id, grants_json, expires_at,
+				created_at+1, project, profile_digest FROM session_credentials WHERE episode_id=?`, p.start.Episode)
+		}},
+		{"project", errOperatorFinishEligibility, func(t *testing.T) {
+			p.exec(t, "UPDATE session_credentials SET project=? WHERE episode_id=?", "other-project", p.start.Episode)
+		}},
+		{"session-file", errOperatorFinishEligibility, func(t *testing.T) { chmod(t, f.f.r.sessions, 0755) }},
+		{"token-mode", errOperatorFinishEligibility, func(t *testing.T) { chmod(t, token, 0644) }},
+		{"commit-identity", errOperatorFinishFence, func(t *testing.T) {
+			alterPreset(t, func(l *launchPreset, _ *presetDaemon) { l.Committer = "" })
+		}},
+		{"reused-pid", nil, func(t *testing.T) {
+			old := launch
+			old.PID, old.Birth = os.Getpid(), "previous-process-birth"
+			b, _ := json.Marshal(old)
+			p.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", b)
+		}},
+	} {
+		t.Run("input/"+tc.name, func(t *testing.T) {
+			t.Cleanup(func() {
+				f.c, f.health, *p.report, *p.malformed = config, health, originalReport, false
+				seedStartFixtureFile(t, preset, presetWire)
+				p.exec(t, "UPDATE episode_starts SET start_json=? WHERE episode=?", string(p.startWire), p.start.Episode)
+				p.exec(t, "UPDATE operator_action_steps SET version=4, receipt=? WHERE kind='launch'", launchWire)
+				p.exec(t, "DELETE FROM session_credentials WHERE credential_id='zzz-fixture'")
+				p.exec(t, "UPDATE session_credentials SET project=?,profile_digest=? WHERE credential_id=?",
+					sessions[0].Project, sessions[0].ProfileDigest, sessions[0].CredentialID)
+			})
+			tc.edit(t)
+			p.unchanged(t, tc.want, func() error { return input(t) })
+			// Also observe the route-claimed worker's earlier admission
+			// boundary. The isolated input assertion above distinguishes
+			// the branch from a later frozen-manifest refusal.
+			if tc.want != nil {
+				want := tc.want
+				if tc.name == "remote-allowlist" {
+					want = approval.ErrNotVerified // frozen remote checked before live inputs
+				}
+				p.unchanged(t, want, func() error { return p.preclaim(t, activationTestContext(t)) })
+			}
+		})
+	}
+}
+
+func (p finishRechecks) catalogue(t *testing.T) {
+	t.Helper()
+	f := p.f
+	ctx := activationTestContext(t)
+	t.Run("catalogue/approval-size", func(t *testing.T) {
+		report := *p.report
+		t.Cleanup(func() {
+			p.setStart(t, p.start)
+			*p.report = report
+		})
+		start := *p.start
+		task := *start.Task
+		task.Title = strings.Repeat("x", approval.MaxOperatorActionBytes)
+		start.Task = &task
+		p.setStart(t, &start)
+		in, err := f.env.finishInputs(ctx, f.c, f.src, f.db, start.Episode, nil)
+		if err != nil {
+			t.Fatal(err)
+		}
+		m, err := in.frozen(ctx)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if _, err := m.Action(1, false); !errors.Is(err, approval.ErrOperatorTooLarge) {
+			t.Fatalf("fixture did not reach approval byte bound: %v", err)
+		}
+		before := p.state(t)
+		if _, ok := f.env.finishCatalogue(ctx, f.c, f.src)[start.Episode]; ok || before != p.state(t) {
+			t.Fatal("unapprovable finish entered the runner catalogue")
+		}
+	})
+	t.Run("catalogue/aggregate-size", func(t *testing.T) {
+		c := f.c
+		c.PresetNames = []string{"padding", p.start.Preset.Name}
+		path := filepath.Join(f.src.presets, "padding.json")
+		var preset launchPreset
+		if err := json.Unmarshal([]byte(p.start.Preset.Bytes), &preset); err != nil {
+			t.Fatal(err)
+		}
+		preset.Name = c.PresetNames[0]
+		raw, _ := json.Marshal(preset)
+		seedStartFixtureFile(t, path, raw)
+		defer os.Remove(path)
+		capBytes := 16 * approval.MaxOperatorActionBytes
+		small := f.env.readOperatorCatalogue(c, f.src)
+		var all map[string]json.RawMessage
+		if err := json.Unmarshal(small, &all); err != nil || all["padding"] == nil || all[store.OperatorFinishCatalogueKey] == nil {
+			t.Fatalf("fixture never published the start/finish catalogue: %d bytes, %v", len(small), err)
+		}
+		// Whitespace in the preset is part of its exact catalogued bytes.
+		seedStartFixtureFile(t, path, append(raw, bytes.Repeat([]byte(" "), capBytes-len(small)+1)...))
+		before := p.state(t)
+		out := f.env.readOperatorCatalogue(c, f.src)
+		all = nil
+		if err := json.Unmarshal(out, &all); err != nil || len(out) > capBytes || all["padding"] == nil ||
+			all[p.start.Preset.Name] == nil || all[store.OperatorFinishCatalogueKey] != nil || before != p.state(t) {
+			t.Fatalf("over-limit finish catalogue did not preserve start entries: %d bytes; %v", len(out), err)
+		}
+	})
+}
+
+// Cancel at the lifecycle boundary using the same nonblocking lock as the
+// caller. The signal is driven by lock ownership, never by elapsed time.
+type finishLifecycleContext struct {
+	context.Context
+	cancel context.CancelFunc
+	root   string
+	mu     sync.Mutex
+}
+
+func (c *finishLifecycleContext) Err() error {
+	// Child contexts may probe Err concurrently. Only production lock
+	// ownership may trigger cancellation.
+	c.mu.Lock()
+	defer c.mu.Unlock()
+	lock, err := runnerLifecycleLock(c.root)
+	if lock != nil {
+		_ = lock.Close()
+	} else if errors.Is(err, syscall.EWOULDBLOCK) {
+		c.cancel()
+	}
+	return c.Context.Err()
+}
+
+func (p finishRechecks) claim(t *testing.T, raw []byte) {
+	t.Helper()
+	f := p.f
+	ctx := activationTestContext(t)
+	h, err := f.db.OperatorHeartbeat(ctx)
+	if err != nil {
+		t.Fatal(err)
+	}
+	env := f.env
+	launch := f.receipt(t).Receipt
+	launchWire, _ := json.Marshal(launch)
+	profile, profilePath, err := loadMintProfile(f.src.profiles, "proj")
+	if err != nil {
+		t.Fatal(err)
+	}
+	profileWire := readStartFixtureFile(t, profilePath)
+	setLaunch := func(t *testing.T) {
+		t.Helper()
+		changed := launch
+		changed.Birth += "-other"
+		b, _ := json.Marshal(changed)
+		p.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", b)
+		t.Cleanup(func() { p.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='launch'", launchWire) })
+	}
+	table := []struct {
+		name string
+		want error
+		edit func(*testing.T)
+	}{
+		{"catalogue-fresh", store.ErrOperatorRunner, func(t *testing.T) {
+			p.exec(t, "UPDATE operator_runner_state SET catalogue_generation=catalogue_generation+1")
+		}},
+		{"heartbeat-owner", store.ErrOperatorRunner, func(*testing.T) { f.h.Owner.Generation++ }},
+		{"heartbeat-ready", store.ErrOperatorRunner, func(*testing.T) { f.h.Build = "different-fixture-build" }},
+		{"life-quiet", store.ErrOperatorRunner, func(*testing.T) {
+			f.env.runnerProcesses = func() ([]worldProcess, error) { return nil, errors.New("inventory unavailable") }
+		}},
+		{"life-pair", store.ErrOperatorRunner, func(t *testing.T) {
+			path := f.env.runnerPaths().Plist
+			b := readStartFixtureFile(t, path)
+			seedStartFixtureFile(t, path, append(bytes.Clone(b), '\n'))
+			t.Cleanup(func() { seedStartFixtureFile(t, path, b) })
+		}},
+		{"manifest-protected", errOperatorFinishReview, func(t *testing.T) {
+			var edited map[string]any
+			if err := json.Unmarshal(profileWire, &edited); err != nil {
+				t.Fatal(err)
+			}
+			edited["protected"] = []string{filepath.Base(p.change)}
+			b, _ := json.Marshal(edited)
+			seedStartFixtureFile(t, profilePath, b)
+			next, _, err := loadMintProfile(f.src.profiles, "proj")
+			if err != nil {
+				t.Fatal(err)
+			}
+			p.exec(t, "UPDATE session_credentials SET profile_digest=? WHERE episode_id=?", next.Digest, p.start.Episode)
+			t.Cleanup(func() {
+				seedStartFixtureFile(t, profilePath, profileWire)
+				p.exec(t, "UPDATE session_credentials SET profile_digest=? WHERE episode_id=?", profile.Digest, p.start.Episode)
+			})
+			p.signCurrentReview(t, true)
+		}},
+		{"manifest-unattributed", errOperatorFinishReview, func(t *testing.T) {
+			old := p.report.Writes.Written
+			p.report.Writes.Written = nil
+			t.Cleanup(func() { p.report.Writes.Written = old })
+			p.signCurrentReview(t, false)
+		}},
+	}
+	for _, tc := range table {
+		t.Run("claim/"+tc.name, func(t *testing.T) {
+			original := f.h
+			t.Cleanup(func() {
+				f.h, f.env = original, env
+				p.exec(t, "UPDATE operator_runner_state SET catalogue_generation=?", h.CatalogueGeneration)
+			})
+			tc.edit(t)
+			p.unchanged(t, tc.want, func() error { return p.preclaim(t, activationTestContext(t)) })
+		})
+	}
+	t.Run("claim/cancelled", func(t *testing.T) {
+		bound, cancel := context.WithCancel(ctx)
+		defer cancel()
+		lifecycle := &finishLifecycleContext{Context: bound, cancel: cancel, root: filepath.Dir(env.paths.Record)}
+		p.unchanged(t, store.ErrOperatorRunner, func() error {
+			return p.preclaim(t, lifecycle)
+		})
+	})
+	t.Run("claim/start-ready-loop", func(t *testing.T) {
+		original := f.h
+		f.h.Readiness = "start"
+		f.env.operatorCatalogue = func(runnerConfig) []byte { return raw }
+		if err := f.env.touchRunner(f.c, &f.h); err != nil {
+			t.Fatal(err)
+		}
+		t.Cleanup(func() {
+			f.h, f.env = original, env
+			p.exec(t, "UPDATE operator_runner_state SET readiness=?, catalogue_generation=?", h.Readiness, h.CatalogueGeneration)
+		})
+		before := p.state(t)
+		if log := f.tick(t); log != "" || before != p.state(t) {
+			t.Fatalf("start-only runner attempted a finish: %q", log)
+		}
+	})
+
+	t.Run("after-claim", func(t *testing.T) {
+		held := f.borrow(t)
+		defer held.Close()
+		fd, err := held.Verify(f.creds, p.start.Episode)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if err := runnerWrite(f.creds, func(db *store.Credentials) error {
+			write, cancel := f.env.operatorContext(activationTestBudget)
+			defer cancel()
+			return db.ClaimOperatorAction(write, operatorPreparationCAS(p.row), h.Owner, f.now)
+		}); err != nil {
+			t.Fatal(err)
+		}
+		t.Cleanup(func() {
+			p.exec(t, "DELETE FROM operator_action_steps WHERE action_id=? AND kind<>'operation'", p.row.ID)
+			p.exec(t, "UPDATE operator_actions SET state=?,generation=?,owner=NULL WHERE id=?", p.row.State, p.row.Generation, p.row.ID)
+		})
+		row, err := f.db.OperatorAction(ctx, p.row.ID, f.now)
+		if err != nil || row.Generation != p.row.Generation+1 || row.Owner != h.Owner.Nonce {
+			t.Fatalf("expected exact CAS claim: %s/%d %v", row.Owner, row.Generation, err)
+		}
+		a, m, generation, err := finishManifestForRow(f.c, row)
+		if err != nil {
+			t.Fatal(err)
+		}
+		construct := func() (operatorFinishWork, error) {
+			return f.env.finishWork(ctx, f.c, h, f.src, f.db, fd, row, operatorPreparationCAS(row), a, m, generation, false, f.now)
+		}
+		t.Run("review", func(t *testing.T) {
+			setLaunch(t)
+			p.unchanged(t, errOperatorFinishReview, func() error { _, err := construct(); return err })
+		})
+		work, err := construct()
+		if err != nil {
+			t.Fatal(err)
+		}
+		// The production check seam must be unreachable in every refusal case.
+		work.run = func(context.Context, store.StartCheck) (checkrun.Result, error) {
+			*p.checks++
+			return checkrun.Result{}, errors.New("unexpected fixture check")
+		}
+		for _, tc := range []struct {
+			name string
+			want error
+			edit func(*testing.T)
+		}{
+			{"clock", store.ErrOperatorRunner, func(t *testing.T) {
+				p.exec(t, "UPDATE operator_runner_state SET seen_at=?", h.SeenAt-10)
+				t.Cleanup(func() { p.exec(t, "UPDATE operator_runner_state SET seen_at=?", h.SeenAt) })
+			}},
+			{"inputs", errOperatorFinishReview, setLaunch},
+			{"key", approval.ErrNotVerified, func(t *testing.T) {
+				p.exec(t, "UPDATE approvers SET revoked_at=? WHERE name='finish-operator-fixture'", f.now)
+				t.Cleanup(func() { p.exec(t, "UPDATE approvers SET revoked_at=NULL WHERE name='finish-operator-fixture'") })
+			}},
+		} {
+			t.Run(tc.name, func(t *testing.T) {
+				tc.edit(t)
+				p.unchanged(t, tc.want, func() error { return runEpisodeFinishWithPrep(ctx, work) })
+			})
+		}
+	})
+}
+
+// Seed the same canonical approval and genuine fixture signature used by the
+// confirm route, bypassing request admission to exercise the later re-check.
+func (p finishRechecks) signCurrentReview(t *testing.T, protected bool) {
+	t.Helper()
+	f := p.f
+	ctx := activationTestContext(t)
+	i, err := f.env.finishInputs(ctx, f.c, f.src, f.db, p.start.Episode, nil)
+	if err != nil {
+		t.Fatal(err)
+	}
+	m, err := i.frozen(ctx)
+	if err != nil || protected && (len(m.Protected) != 1 || len(m.Unattributed) != 0) ||
+		!protected && (len(m.Unattributed) != 1 || len(m.Protected) != 0) {
+		t.Fatalf("fixture did not freeze exact change classification: %v %v %v", m.Protected, m.Unattributed, err)
+	}
+	_, _, generation, err := finishManifestForRow(f.c, p.row)
+	if err != nil {
+		t.Fatal(err)
+	}
+	a, err := m.Action(generation, false)
+	if err != nil {
+		t.Fatal(err)
+	}
+	a.ID = p.row.ID
+	raw, err := approval.CanonicalOperatorAction(a)
+	if err != nil {
+		t.Fatal(err)
+	}
+	digest := sha256.Sum256(raw)
+	message, err := approval.OperatorActionMessage(digest, [32]byte(p.row.Challenge), p.row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	const update = "UPDATE operator_actions SET action=?,digest=?,artifact=? WHERE id=?"
+	p.exec(t, update, raw, digest[:], p.sign(message), p.row.ID)
+	t.Cleanup(func() { p.exec(t, update, p.row.Action, p.row.Digest, p.row.Artifact, p.row.ID) })
+}
+
+func (p finishRechecks) recovery(t *testing.T, row store.OperatorActionRow, complete store.OperatorFinishJournal, status func(string, string)) {
+	t.Helper()
+	f := p.f
+	ctx := activationTestContext(t)
+	original, _ := json.Marshal(complete.Commit)
+	const reset = "UPDATE operator_action_steps SET phase=?,version=?,receipt=? WHERE action_id=? AND kind='commit'"
+	before := p.state(t)
+	t.Cleanup(func() {
+		p.exec(t, reset, "observed", complete.CommitKey.Version, original, row.ID)
+	})
+	p.exec(t, reset, "reserved", complete.CommitKey.Version-1, nil, row.ID)
+	p.unchanged(t, nil, func() error {
+		status("Finish outcome uncertain; inspect durable receipts", "")
+		return nil
+	})
+	now := f.now
+	f.now = row.ExpiresAt
+	err := p.run()
+	f.now = now
+	j, readErr := f.db.ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))
+	after := p.state(t)
+	if err != nil || readErr != nil || !reflect.DeepEqual(j, complete) || !reflect.DeepEqual(before[1:], after[1:]) {
+		t.Fatalf("expired reserved commit did not observe the exact existing object: %v %v %+v", err, readErr, j)
+	}
+
+	failed, _ := json.Marshal(store.OperatorFinishLocalCommit{State: "commit-failed"})
+	p.exec(t, reset, "observed", complete.CommitKey.Version, failed, row.ID)
+	p.unchanged(t, nil, func() error {
+		status("Local commit failed; check receipts retained", "")
+		return nil
+	})
 }
diff --git a/cmd/ailang-worldd/operator_finish_tree_test.go b/cmd/ailang-worldd/operator_finish_tree_test.go
index dfe1c88..c6b6810 100644
--- a/cmd/ailang-worldd/operator_finish_tree_test.go
+++ b/cmd/ailang-worldd/operator_finish_tree_test.go
@@ -315,6 +315,42 @@ func TestOperatorFinishConfigReplacement(t *testing.T) {
 	g := operatorFinishGitFixture(t, f)
 	ctx := activationTestContext(t)
 	path := filepath.Join(g.common, "config")
+	t.Run("valid-config-with-different-digest", func(t *testing.T) {
+		snap, contents, err := g.snapshot(ctx, f.rep.Start)
+		if err != nil {
+			t.Fatal(err)
+		}
+		spec := store.OperatorLocalCommit{Snapshot: snap, Message: "fixture commit\n", Committer: "Fixture Operator <fixture@example.invalid>"}
+		tree, err := g.planTree(ctx, f.root, spec, contents)
+		if err != nil {
+			t.Fatal(err)
+		}
+		body := readStartFixtureFile(t, path)
+		// A harmless comment changes the pin without making the config invalid.
+		seedStartFixtureFile(t, path, append(bytes.Clone(body), []byte("\n# changed fixture digest\n")...))
+		t.Cleanup(func() { seedStartFixtureFile(t, path, body) })
+		if hash, err := g.config(ctx, f.rep.Start); err != nil || hash == snap.Config {
+			t.Fatalf("fixture config did not remain valid while changing its digest: %v", err)
+		}
+		marker := filepath.Join(f.root, "unexpected-commit-object")
+		bin := filepath.Join(f.root, "observe-commit-object")
+		write := "#!/bin/sh\nfor arg; do if [ \"$arg\" = commit-tree ]; then touch " + shellQuote(marker) +
+			"; fi; done\nexec " + shellQuote(f.git) + " \"$@\"\n"
+		if err := os.WriteFile(bin, []byte(write), 0700); err != nil {
+			t.Fatal(err)
+		}
+		attempt := g
+		attempt.bin = bin
+		if err := attempt.commit(ctx, f.root, spec, contents, tree, f.rep.Start); !errors.Is(err, errOperatorFinishReview) {
+			t.Fatalf("digest mismatch lost review refusal: %v", err)
+		}
+		if _, err := os.Lstat(marker); !errors.Is(err, os.ErrNotExist) {
+			t.Fatalf("digest mismatch reached commit-tree: %v", err)
+		}
+		if after, _, err := g.snapshot(ctx, f.rep.Start); err != nil || after.Head != snap.Head || after.Index != snap.Index {
+			t.Fatalf("digest mismatch changed the ref or index: %v", err)
+		}
+	})
 	if _, err := g.config(ctx, f.rep.Start); err != nil {
 		t.Fatal(err)
 	}
diff --git a/host/daemon/surface_operator_finish_test.go b/host/daemon/surface_operator_finish_test.go
index 852214d..efad754 100644
--- a/host/daemon/surface_operator_finish_test.go
+++ b/host/daemon/surface_operator_finish_test.go
@@ -1,6 +1,7 @@
 package daemon
 
 import (
+	"context"
 	"crypto/sha256"
 	"encoding/hex"
 	"encoding/json"
@@ -8,6 +9,7 @@ import (
 	"net/http/httptest"
 	"net/url"
 	"os"
+	"path/filepath"
 	"reflect"
 	"strings"
 	"testing"
@@ -78,6 +80,48 @@ func TestOperatorFinishAdmission(t *testing.T) {
 		t.Fatalf("stale catalogue changed state: %+v", got)
 	}
 
+	for _, tc := range []struct {
+		name, class string
+		code        int
+		edit        func(*store.OperatorFinishReview)
+	}{
+		{"other-store", "RunnerUnavailable", http.StatusServiceUnavailable, func(m *store.OperatorFinishReview) {
+			m.Store = filepath.Join(filepath.Dir(path), "other-fixture.db")
+		}},
+		{"protected", "RunnerUnavailable", http.StatusConflict, func(m *store.OperatorFinishReview) { m.Protected = []string{"protected.txt"} }},
+		{"unattributed", "RunnerUnavailable", http.StatusConflict, func(m *store.OperatorFinishReview) { m.Unattributed = []string{"unattributed.txt"} }},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			m := manifest
+			tc.edit(&m)
+			current := f.h.Catalogue
+			catalogue[store.OperatorFinishCatalogueKey] = opJSON(map[string]store.OperatorFinishReview{"ep": m})
+			f.h.Catalogue = opJSON(catalogue)
+			f.publish()
+			t.Cleanup(func() {
+				f.h.Catalogue = current
+				f.publish()
+			})
+			ro, err := store.OpenCredentialsReadOnly(path)
+			if err != nil {
+				t.Fatal(err)
+			}
+			defer ro.Close()
+			h, err := ro.OperatorHeartbeat(boundedTestContext(t))
+			if err != nil {
+				t.Fatal(err)
+			}
+			attempt := req
+			attempt.CatalogueGeneration = h.CatalogueGeneration
+			before := f.evidence()
+			got := post(operatorPath, opJSON(attempt))
+			if got.Code != tc.code || !strings.Contains(got.Body.String(), tc.class) || before != f.evidence() {
+				t.Fatalf("request refusal lost class or changed durable state: HTTP %d %s", got.Code, got.Body.String())
+			}
+		})
+	}
+	review = get(query)
+	req.CatalogueGeneration = review.CatalogueGeneration
 	key := approvaltest.NewPasskey(approval.AlgEdDSA, 43)
 	f.enrolPasskey("finish-key", key, opOrigin)
 	var requested operatorStatus
@@ -100,6 +144,32 @@ func TestOperatorFinishAdmission(t *testing.T) {
 	if got := post(operatorPath, opJSON(req)); got.Code != 200 || !strings.Contains(got.Body.String(), requested.ID) {
 		t.Fatalf("retry lost same request: %+v", got)
 	}
+	t.Run("existing-operation/other-generation", func(t *testing.T) {
+		other := req
+		other.CatalogueGeneration++
+		before := f.evidence()
+		got := post(operatorPath, opJSON(other))
+		if got.Code != http.StatusConflict || !strings.Contains(got.Body.String(), `"Conflict"`) ||
+			f.evidence() != before || !reflect.DeepEqual(row, f.row(row.ID)) {
+			t.Fatalf("retry at another catalogue generation changed state or lacked conflict: HTTP %d %s", got.Code, got.Body.String())
+		}
+	})
+	t.Run("status/start-only-runner", func(t *testing.T) {
+		f.h.Readiness = "start"
+		f.publish()
+		t.Cleanup(func() { f.h.Readiness = "both"; f.publish() })
+		before := f.evidence()
+		resp := f.do("GET", requested.StatusHref, nil, nil)
+		var read operatorStatus
+		const reason = "finish blocked: github-auth; repair PATH, GitHub login or Git credential helper in attended service-runner setup and re-run readiness"
+		if resp.Code != http.StatusOK || json.Unmarshal(resp.Body.Bytes(), &read) != nil ||
+			read.Reason != reason || read.State != "requested" || read.Kind != "finish" || read.ID != row.ID ||
+			read.RunnerSeenAt != f.h.SeenAt || read.LocalCommit != "" || len(read.Checks) != 0 ||
+			!strings.Contains(f.do("GET", requested.Href, nil, nil).Body.String(), reason) || f.evidence() != before ||
+			!reflect.DeepEqual(row, f.row(row.ID)) {
+			t.Fatalf("start-only runner lost finish remedy or changed evidence: %+v", read)
+		}
+	})
 	artifact := approveBody(hex.EncodeToString(row.Digest),
 		key.Assert(opOrigin, "localhost", sha256.Sum256(opMessage(t, [32]byte(row.Digest), [32]byte(row.Challenge), row.ExpiresAt)), 2))
 	before = f.evidence()
@@ -121,4 +191,35 @@ func TestOperatorFinishAdmission(t *testing.T) {
 			t.Fatal("confirmed status exposed prep data or issued a V3 prompt")
 		}
 	}
+
+	t.Run("confirm/without-surface-option", func(t *testing.T) {
+		other := manifest
+		other.Episode = "finish-elsewhere"
+		a, err := other.Action(req.CatalogueGeneration, false)
+		if err != nil {
+			t.Fatal(err)
+		}
+		a.ID = "oa-" + strings.Repeat("b", 32)
+		a.Options.ConfirmOnSurface = false
+		a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+		challenge := sha256.Sum256([]byte(t.Name()))
+		var pending store.OperatorActionRow
+		f.write(func(ctx context.Context, c *store.Credentials) error {
+			var err error
+			pending, err = c.RequestOperatorAction(ctx, store.OperatorRequest{
+				Action: a, OperationID: opRequest(55).OperationID, Challenge: challenge,
+				CreatedAt: apT0, ExpiresAt: apT0 + 600,
+			}, apT0)
+			return err
+		})
+		digest := hex.EncodeToString(pending.Digest)
+		msg := opMessage(t, [32]byte(pending.Digest), challenge, pending.ExpiresAt)
+		proof := approveBody(digest, key.Assert(opOrigin, "localhost", sha256.Sum256(msg), 3))
+		before := f.evidence()
+		got := post(operatorPath+"/"+pending.ID+"/confirm", proof)
+		if got.Code != http.StatusConflict || !strings.Contains(got.Body.String(), `"Conflict"`) ||
+			f.evidence() != before || !reflect.DeepEqual(pending, f.row(pending.ID)) {
+			t.Fatalf("signed non-surface finish crossed confirm: HTTP %d %s", got.Code, got.Body.String())
+		}
+	})
 }
````

### M2 — Read-only finish-status mirror

````diff M2
diff --git a/cmd/ailang-worldd/operator_finish_route_test.go b/cmd/ailang-worldd/operator_finish_route_test.go
index 5e40eae..3223fed 100644
--- a/cmd/ailang-worldd/operator_finish_route_test.go
+++ b/cmd/ailang-worldd/operator_finish_route_test.go
@@ -368,7 +368,12 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
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
@@ -381,6 +386,7 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 	if _, err := os.Lstat(trace); !errors.Is(err, os.ErrNotExist) {
 		t.Fatalf("remote write call escaped the worker: %v", err)
 	}
+	probe.cli(t, row, j, filename)
 	probe.recovery(t, row, j, func(reason, sha string) {
 		t.Helper()
 		status.Reason, status.LocalCommit, status.Checks = "", "", nil
@@ -393,6 +399,15 @@ func TestOperatorFinishRouteWorker(t *testing.T) {
 		if card := request("GET", status.Href, nil, nil); !strings.Contains(card, reason) {
 			t.Fatalf("finish card lost the projected reason %q", reason)
 		}
+		want := errOperatorFinishUncertain.Error()
+		if strings.HasPrefix(reason, "Local commit failed") {
+			want = errOperatorCommitFailed.Error()
+		}
+		var out, errOut syncBuffer
+		if rc := runServiceRunner([]string{"finish-status", row.ID}, &out, &errOut, f.env); rc != exitOK ||
+			!strings.Contains(out.String(), want+"\n") || errOut.String() != "" {
+			t.Fatalf("CLI mirror lost journal remedy: rc=%d stdout=%q stderr=%q", rc, &out, &errOut)
+		}
 	})
 }
 
@@ -966,3 +981,94 @@ func (p finishRechecks) recovery(t *testing.T, row store.OperatorActionRow, comp
 		return nil
 	})
 }
+
+func (p finishRechecks) cli(t *testing.T, row store.OperatorActionRow, j store.OperatorFinishJournal, filename string) {
+	t.Helper()
+	f := p.f
+	docs := readStartFixtureFile(t, "../../docs/QUICKSTART.md")
+	_, section, found := strings.Cut(string(docs), "### Prepare a local finish in the World app")
+	_, command, marked := strings.Cut(section, "receipts are available from `")
+	command, _, closed := strings.Cut(command, "`")
+	argv := strings.Fields(command)
+	if !found || !marked || !closed || !reflect.DeepEqual(argv, []string{"ailang-worldd", "service", "runner", "finish-status", "ACTION-ID"}) {
+		t.Fatalf("QUICKSTART changed the printed command: %q", command)
+	}
+	args := func(id string) []string { return []string{argv[2], argv[3], id} }
+	call := func(t *testing.T, env serviceEnv, command []string, want int, stdout, stderr string) {
+		t.Helper()
+		before := p.state(t)
+		var out, errOut syncBuffer
+		rc := runServiceWith(command, &out, &errOut, env)
+		for _, secret := range []string{filename, "private-task-fixture-marker", "BEGIN SSH SIGNATURE", f.token, "finish_manifest", j.Checks[0].Ref} {
+			if strings.Contains(out.String(), secret) || strings.Contains(errOut.String(), secret) {
+				t.Fatalf("CLI sink exposed protected fixture information")
+			}
+		}
+		if rc != want || out.String() != stdout || !strings.Contains(errOut.String(), stderr) ||
+			stderr == "" && errOut.String() != "" || before != p.state(t) {
+			t.Fatalf("CLI read lost receipt/refusal or changed state: rc=%d stdout=%q stderr=%q", rc, &out, &errOut)
+		}
+	}
+	line := fmt.Sprintf("finish action %s %s\n", row.ID, row.State)
+	check := fmt.Sprintf("check %q: passed\n", j.Checks[0].Name)
+	commit := fmt.Sprintf("local commit %s retained; PR confirmation has not been issued\n", j.Commit.SHA)
+	t.Run("cli/printed-command", func(t *testing.T) { call(t, f.env, args(row.ID), exitOK, line+check+commit, "") })
+	t.Run("cli/help", func(t *testing.T) {
+		before := p.state(t)
+		var out, errOut syncBuffer
+		if rc := runServiceWith([]string{"runner", "--help"}, &out, &errOut, f.env); rc != exitOK ||
+			!strings.Contains(out.String(), " | finish-status ACTION-ID | ") || errOut.String() != "" || before != p.state(t) {
+			t.Fatal("operator runner help lost the finish-status command")
+		}
+	})
+	for _, tc := range []struct {
+		name   string
+		argv   []string
+		reason string
+	}{
+		{"unknown-action", args("oa-" + strings.Repeat("f", 32)), store.ErrOperatorMissing.Error()},
+		{"non-finish", args(f.action.ID), store.ErrOperatorMissing.Error()},
+		{"invalid-id", args("not-an-action"), "expected finish-status ACTION-ID"},
+		{"extra-argument", append(args(row.ID), "extra"), "expected finish-status ACTION-ID"},
+	} {
+		t.Run("cli/"+tc.name, func(t *testing.T) { call(t, f.env, tc.argv, exitUsage, "", tc.reason) })
+	}
+	t.Run("cli/uninstalled", func(t *testing.T) {
+		e := f.env
+		e.paths.Record = filepath.Join(t.TempDir(), "absent", "service.json")
+		call(t, e, args(row.ID), exitUsage, "", "runner file is absent or inaccessible")
+	})
+	t.Run("cli/unreadable-credentials", func(t *testing.T) {
+		c := f.c
+		dir, err := os.MkdirTemp(c.Home, "fixture-credentials-")
+		if err != nil {
+			t.Fatal(err)
+		}
+		c.Credentials = filepath.Join(dir, "not-sqlite.db")
+		seedStartFixtureFile(t, c.Credentials, []byte("fixture non-database"))
+		paths := f.env.runnerPaths()
+		record, plist := readStartFixtureFile(t, paths.Record), readStartFixtureFile(t, paths.Plist)
+		t.Cleanup(func() {
+			seedStartFixtureFile(t, paths.Record, record)
+			seedStartFixtureFile(t, paths.Plist, plist)
+		})
+		b, _ := json.Marshal(c)
+		seedStartFixtureFile(t, paths.Record, b)
+		seedStartFixtureFile(t, paths.Plist, []byte(runnerPlist(paths, c, b)))
+		call(t, f.env, args(row.ID), exitUsage, "", errOperatorFinishEligibility.Error())
+	})
+	t.Run("cli/open-check", func(t *testing.T) {
+		receipt, _ := json.Marshal(j.Checks[0])
+		const update = "UPDATE operator_action_steps SET phase=?,receipt=? WHERE action_id=? AND kind='checks'"
+		p.exec(t, update, "reserved", nil, row.ID)
+		t.Cleanup(func() { p.exec(t, update, "observed", receipt, row.ID) })
+		call(t, f.env, args(row.ID), exitOK, line+errOperatorFinishUncertain.Error()+"\n"+commit, "")
+	})
+	t.Run("cli/unreadable-journal", func(t *testing.T) {
+		receipt, _ := json.Marshal(j.Commit)
+		const update = "UPDATE operator_action_steps SET receipt=? WHERE action_id=? AND kind='commit'"
+		p.exec(t, update, []byte("{}"), row.ID)
+		t.Cleanup(func() { p.exec(t, update, receipt, row.ID) })
+		call(t, f.env, args(row.ID), exitUsage, "", errOperatorFinishUncertain.Error())
+	})
+}
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
index 679951c..5e2c490 100644
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

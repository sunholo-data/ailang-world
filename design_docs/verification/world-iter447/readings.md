# World iteration 447 — readings (row 243 V4 S8)

All readings are first-party controller measurements unless marked otherwise. Paths under `~/.ailang/state/world-iter447/` are
rig-local scratch and are not published.

## Gate 0–2
- Kill switch absent; `gh` = `sunholo-voight-kampff`; billing tripwire CLEAN.
- `mission_directives.sh --issue 635`: 0 allowlisted directives (33 comments; LEGACY watermark). Mission-scoped watermark
  `mission-world-last-seen` = `2026-10-27T09:03:42Z` (older of the two is the issue-scoped `2026-08-17T06:15:12Z`; both read).
- Step 6a `mission_gate0_self_notices.sh` rc 0: 0 crash notices (control 107:4 ok; `#574` watermark absent → DEGRADED single-watermark
  read). 3 transcript-scan notices in the window, fires `world-1795730668-41608` (2 hits), `world-1795769251-75944` (1 hit) and
  `world-1795786215-30397` (3 hits) — the same three iteration 446 reported; no new one. Key names are in the private controlplane notice only.
- `mission-world` inbox: 0 unread messages addressed to this mission (the unread list is fleet-wide controlplane traffic).
- Local `dev` = `origin/dev` = `0e32a8d7` (no fast-forward needed). No open own-account PRs.
- Skill drift: the pin worktree matched origin on all 14 skill files. The symlink-resolved copy (V1 main checkout, 27 commits behind)
  differed from origin on `resources/gate-5-retro.md` only (one line: origin's idle-fingerprint input list adds "open nightly-rig issues
  and their failure reports"). Gate 5 read origin's version.
- Ledger valid, 108 rows, 0 OPEN. Clause map: 1.0 clauses 1–7 MET, none UNMET.
- Dev CI on `0e32a8d7`: `checks=3`, all `success`.
- Pick re-measured at `0e32a8d7`: `app.js` 3,899 lines with 0 `operator-actions` references; no S8 plan, branch, PR or worktree
  (`git branch -a`, `gh pr list --search S8`, `git worktree list`). Quorum reader rc 0: V4 r2 `…T10-24-51Z` PROCEED 3/3,
  `absent_reviewers` `[]`.
- Design §8 says S8 waits for "119 A/B, S7, E5 and #604 sequencing". 119 = A (attended); S7/S7b landed; **#604 MERGED
  2026-11-18T09:48:30Z as `65108353`** (`gh pr view 604`), so the sequencing gate is resolved against post-merge dev; **E5** (§11:
  "E5 requires Mark's chosen preset") has never run — `grep -rln V4_PRESET_FILE design_docs/` finds only the design itself, and the rig's
  presets directory holds seven presets, none named for Sol. The loop cannot choose Mark's preset, so S8 was routed against a FAKE Sol
  preset with E5's loop-runnable half (W25/W32) done in planning, and D-WORLD-122 asks Mark to name the preset before S11.
- Weekly external-issue sweep not due (iteration 424 ran it after the 2026-11-23 rotation). `#635` has 33 comments, created
  2026-11-23T07:34Z (after the Monday 07:00 local boundary): no rotation.

## Routing
- `resolve-role-spawn.sh`: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`; executor `recipe codex:gpt-6.2-sol
  declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`. `derive-planner-lane.sh` same as the planner line.
- Codex probe rc 0 (`ok`). Codex usage 14.0% used / 25.9% allowed (`ailang mission quota --bucket codex`).

## Baseline (pristine `0e32a8d7`, outside the sandbox, isolated `AILANG_STATE_DIR`, while the planner ran)
store-plain 0 · store-race 0 · authority-plain 0 · authority-race 0 · agentrun-plain 0 · agentrun-race 0 · verifygate-plain 0 ·
verbs-plain 0 · broker-plain 4 (`TestExecHandlerWorktreeDelta`, `TestExecProtectedCheckTimeout`, `TestExecProtectedDetectsIgnoredCreation`,
`TestExecProtectedPhaseOrder`) · cmd-plain 14 (the 11 row-271 real-binary tests + `TestOperatorFilePublicationAdmission`,
`TestOperatorFilesReconcile`, `TestOperatorWorktreeReadbackGuards`) · daemon-plain 14 · daemon-race 10. `ALL_DONE for 0e32a8d7…`.

## Planner (`codex:gpt-6.2-sol`, ~77 min, rc 0, 727,578 tok)
- Outcome A: M1 510, M1+M2 1,269 changed lines (`cmd host docs app`). Plan + JSON committed as `cc16e985`.
- The planner executed its own §7 block verbatim against a fresh archive (`APPLY_ACCEPT_SNAPSHOT_OK M1 fast`, `… M2 fast`, six
  `ROUNDTRIP_OK`) — the iteration-446 carry, honoured.
- Controller round trip (`ctl_roundtrip.sh M1 M2` on a fresh `git archive 0e32a8d7`): `diff -r` rc 0 for cmd, host and docs at both stages.
- Planner-reported (not first-party): caller census 1 → 1 → 1 for `RequestOperatorAction`, `JudgeOperatorAction`, store
  `ConfirmOperatorAction`, `ClaimOperatorAction`; 171 diff mutants (145 named kills, 26 outside-only); E15/E19 killed; every route and
  Chrome result UNINFORMATIVE UNDER SANDBOX.

## Executor round 1 (`codex:gpt-6.2-sol`, ~6 min, rc 0, 48,945 tok)
- Ran the plan's §7 block verbatim with `TOOLS=$STATE/exec-tools` (a parameter the block already declares, so planner scratch stayed
  untouched — 17/17 planner tool hashes unchanged afterwards). `ACCEPT_OK M1 fast`, `APPLY_ACCEPT_SNAPSHOT_OK M1 fast`, `ACCEPT_OK M2 fast`,
  `APPLY_ACCEPT_SNAPSHOT_OK M2 fast`. Snapshots M1 11/0/11 and M2 16/0/16, `shasum -c` clean, all files `cmp`-equal to the planner's stages.

## Commit gate
- **M1** `commit_ms.sh M1 … accept.sh M1 outside`: `ACCEPT_OK M1 outside`, committed **`86f40c81`**, NUMSTAT 510.
- **M2 attempt 1** (`accept.sh M2 outside` + `commit-guards.sh`): REFUSED. `TestOperatorStartWalkthroughInChrome`
  `operator_browser_chrome_test.go:193: EOF`, `driver: signal: killed` at 120.05 s (`CHROME operator-chrome RAN=1 SKIP=0`). Reproduced
  alone on a quiet rig (load 2.1): same line, 120.23 s. Existing Chrome controls on the same rig: `TestConfirmPageConfirmsWithAPasskeyInChrome`
  and `TestVerbBarEpisodeReportInChrome` PASS at base (1.39 s, 1.14 s) and the latter PASS on the M2 tree (1.23 s) — harness and page load fine.

## Diagnosis of the walkthrough (scratch copy `diag/tree` of the planner's st-M2, instrumentation there only)
1. Timestamps: driver up T+9.9 s, fixtures done T+17.2 s, then `open("/surface")` never returned. A Go `http.Get` of `/surface` from the
   same fixture answered 200 in 1.5 ms. A CDP trace in the driver showed `SEND Page.navigate {"url":"http://localhost:59937/surface"}` with
   no `RECV` until the kill. **Cause: the test ran `t.Setenv("HOME", t.TempDir())` BEFORE `startPasskeyBrowser`, so Chrome launched under an
   empty HOME.** Moving the Setenv after the browser start: the walkthrough proceeds (F1).
2. After the dropped POST and reload: dump `activeElement=root, operator-request='Reconnect', readOnly=true,
   sessionStorage=[world.operator-request.v1]` — the restore worked, but `boot()`'s last line (`app.js:4022`) moved focus to `#root` after
   the restore focused Reconnect (F2). Adding `!$('verb').contains(document.activeElement)` → `activeElement=operator-request`.
3. On the card: computed `<pre>` white-space `pre` (card.css says `pre-wrap`), `scrollWidth` 1,025 at `innerWidth` 390, and the page's
   CSP-violation recorder held `style-src-elem http://localhost:…/surface/operator-card.css`: the inherited `surfaceCSP` is
   `style-src 'unsafe-inline'` only (`host/daemon/surface_page.go:19`), the same reason the app shell inlines `app.css`
   (`host/surfacepage/app.go:102–109`). Inlining the card CSS: no violation, `scrollWidth` 375 (F3).
4. Two card-page checks asserted `scrollWidth===390`; with a vertical scrollbar it is 375 with no overflowing element (F4, F5).
5. `json: cannot unmarshal object into Go value of type string`: the driver evaluates `JSON.stringify(<expr>)`, so a Promise-valued
   `fetch(…).then(r=>r.text())` is `{}` (F6). Two-step read via `window.__read ?? null`.
6. With F1–F6 applied to the scratch copy: **`--- PASS: TestOperatorStartWalkthroughInChrome (15.23s)`**.

## Executor rounds 2 and 3 (`codex:gpt-6.2-sol`)
- Round 2 (~5 min, rc 0, 90,445 tok) applied F1–F6 as specified; `ACCEPT_OK M2 fast` (137 PASS / 0 FAIL / 0 SKIP incl. subtests);
  `TestOperatorCardAssets` changed from asserting the external stylesheet route to asserting exactly one inline `<style>` equal to
  `OperatorStyle` and no `rel="stylesheet"`; the `GET /surface/operator-card.css` route and `operatorStyle` handler removed; `surfaceCSP`
  unchanged. Controller diff of snap-r2 vs the plan's M2: exactly the six fixes (60 changed lines across 5 files).
- **Controller defect in its own F2 spec, caught on review:** `#verb` is created on demand (`app.js:1577`) and removed (`:3500`, `:3945`),
  so on an ordinary load `$('verb')` is null and `.contains` would throw at the end of every `boot()`. Round 3 (rc 0, 44,953 tok) changed
  that single token to `$('verb')?.contains(…)`; snap-r3 differs from snap-r2 in `app.js` only.

## Commit gate, M2 final
- `commit_ms_r3.sh M2 … accept.sh M2 outside && commit-guards.sh`: **`ACCEPT_OK M2 outside`**, `CHROME operator-chrome RAN=1 SKIP=0`,
  `CHROME chrome-controls RAN=2 SKIP=0`; then `commit-guards.sh` REFUSED: broker `.new` empty (no new broker fails), verifygate `.new`
  = `TestFixtureDiscrimination`, `TestMissionConfigDefaultFlowIntegration` (`context deadline exceeded` in the fixture's script run).
- Attribution: `go list -deps -test ./host/verifygate/` names no other World package (positive control: it lists verifygate itself), and
  S8 changes no script or `.ail`. Alone, head (`.wt-world-iter447-row243-v4s8`): both PASS (57.1 s); alone, base (`.wt-world-iter447-base`):
  3 FAIL (80.0 s). Negative control fired → load, not S8.
- Committed on the `ACCEPT_OK M2 outside` + Chrome receipts for that exact snapshot (manifest re-verified by `commit_ms`), the 445/446
  precedent: **M2 `8e341d9b`**, NUMSTAT 1,268 (`cmd host docs app`).

## Controller outside sweep (`ctl_outside.sh`, on a fresh `git archive 8e341d9b`)
- Trimmed `gates.sh … M2 outside` (the planner's legs minus the whole-package `record`/`compare` legs, skipped for time — CI's
  whole-package Linux plain and race jobs and `commit-guards.sh` cover them): gofmt, build, vet, linux vet, inv-store, inv-boundary,
  inv-cli, inv-verbs, inv-page, inv-verifygate, inv-broker, verify_ail (203 s), callers, inv-daemon all OK; `accept` (outside) rc 1;
  `accept-race` OK (213 s); `t1777` OK (16 s).
- The `accept` rc 1 is Chrome failing to START: `the driver did not start: … timed out waiting for DevToolsActivePort` (the shared
  driver's 30 s wait; load ≈4.8). The `t1777` leg ran the same walkthrough and it PASSED. The failed start orphaned a Chrome (PID 9268,
  profile `…/TestOperatorStartWalkthroughInChrome3155169754/…`, started 04:20:58) — the driver never learned its PID, so nothing reaped it.
  The controller identified it by its profile path and stopped it. A second, unrelated orphan from a 16-day-old
  `TestPasskeyPageRecordsTheFixtures` run (PID 5200 and helpers) is still alive; not ours, reported, not killed.
- Pristine daemon `-race` pair (`TestOperatorConfirmProjection|TestOperatorConfirmationAdmissionNeighbours`, count 2): rc 0, 12 s.
- `whole-mutants.sh`: E15 KILLED (named and whole), E19 KILLED (named and whole); then stopped by the controller (rc 143) because each
  remaining id runs whole `host/daemon` (~19 min). The six carried ids were left to the judge.

## Gate 3b
- PR #705 on `8e341d9b`. Attempt 1: `go host build + test gate` RED on `host/daemon TestA2AInvokeDeadlineClosesSocket` only
  (`invoke_e2e_test.go:412`: the follow-up normal call is refused with "invocation outcome is not confirmed; resend the same task id");
  race and verify gates green. Attempt 2 (re-run of the failed job): the same single test RED. Locally alone at head `-count=3`: ok 6.8 s.
- Attribution: S8 touches no A2A/invocation file; its only `host/daemon` production change is the `GET /surface/operator-card.js`
  route and `Config.ApprovalClock` (defaults to the same `time.Now`), used only by `surface_approve.go:93,250`. S8's daemon test files sort
  after `invoke_e2e_test.go`, so they run after it. **Negative control on CI:** `gh workflow run CI --ref dev` at base `0e32a8d7`
  (run 46104673930) dispatched in the same window as head attempt 3: base host job green; head attempt 3 host job green. With dev's own
  push run at `0e32a8d7` (green), base is 2/2 and head 1/3 on this test. No other CI run in the last 40 failed. Recorded as a flake
  (instances: iteration 445's local whole-package run, and two CI attempts here) — no S8 mechanism found.
- PR CI 3/3 at attempt 3 → `mergeable` MERGEABLE/CLEAN, `origin/dev` unmoved at `0e32a8d7` → merged `--match-head-commit 8e341d9b` →
  **`85ed2234`**.

## Judge (Agent `sonnet`, background, ~54 min, 124,233 tok) — `evaluator-r1.md`
- **PASS 85** (27/22/18/9/9), 0 blocking. Reachability: no page path creates or confirms an action without an explicit click; census 1
  each (S8 adds none); every new `location.assign` comes from the anchored allowlist after `io.current()`.
- Own mutants: 29 specs + 3 Chrome-only + 3 Chrome DOM; 24 killed, 9 survived, 1 did not compile; must-live control survived. Every JS
  kill by a test that executes the JS. Planner outside-only ids re-run on the walkthrough: 6/6 killed.
- Carried: E15, E19 killed; J02/J06/J10/J18/J19 killed on the real listener; the original J22 edit does not compile, variant J22b killed.
- Chrome: walkthrough PASS twice (2.87 s, 2.64 s), 0 SKIP, also under `-race`. Added race ≈2.4 s `host/daemon`, ≈6.3 s `cmd/ailang-worldd`.
- Non-blocking: a delayed auto-click on reload or on the card (C-09/C-10/C-10b) survives the walkthrough (no such product code: a test
  gap); several Chrome-only DOM wiring mutants survive (C-01/S8-25 liveness argument, C-02/S8-26 draft body); the card's
  `state !== 'requested'` guard (S8-14) is defence in depth; the harness's Chrome start is the weak point under load; an outstanding draft
  overrides a freshly reviewed panel. One daemon `-race` selector fail (`TestWorkspaceProtectedPolicyRequests`), passes alone 4.6 s, no
  S8 code in its path.
- Numstat 1,268 (`cmd host docs app`; plan 1,269); unfiltered 7,130 (plan docs). Judge worktree clean at the end. GPU job (`pilot5_run.sh`,
  `train2.py`) and its restore loops alive after the judge.
- Merge CI on `85ed2234`: 3/3 at attempt 1 (run 46106233909) — all three checks success. **LANDED.**

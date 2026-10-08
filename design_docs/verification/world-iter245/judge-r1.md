# Evaluator r1 — iteration 245, row 149 w-workbench-live-and-polish (head 964d3b3, base b702d73)

SCORE 91/100 — PASS (threshold 70). Zero BLOCKING findings.

| Category | Score |
|---|---|
| Correctness | 28/30 |
| Test quality / non-vacuity | 21/25 |
| Design fidelity | 19/20 |
| Safety & conflict surface | 14/15 |
| Docs / usability | 9/10 |

## What I checked and how
- Walked the whole `git diff b702d73..964d3b3` (49 files). Frozen surface (§8): NO hunk in `isProtected`, `acceptedWorkbenchKeys`, `supportedWorkbenchQuery`, `workbenchErrorTemplate`, `host/agui/*`, `host/daemon/agui.go`, `/v1` routes, `tools/launchd/*`, `cmd/`. The only deletion in daemon/workbench.go is the CSP const line.
- AC0 (hunk listing): pre-existing test files changed = `workbench_test.go` (the two CSP literals only), `render_test.go` (+3 appended census lines only), `deadline_guard_test.go` (+1 fixture, PD3). Everything else is new files. PASS.
- Every AC (AC1.1-AC6.2) has an implementation and a named test (list in my notes: TestTokensBothThemes ... TestWorkbenchLiveChromeDrill). Tests green here: workbench/store/broker `-count=1` ok (broker 104 s); new daemon tests `-race -count=5` green; CI step command reproduced locally with `-race`: `--- PASS: TestLiveScriptAssetNonEmpty / TestLiveScriptPure / TestLiveScriptContract`.
- CI cannot silently skip the node harness: step lives in `go-verify` after the step that exports `WORLD_EXEC_NODE` to GITHUB_ENV (ci.yml:216), has no `if:`, requires `WORLD_EXEC_NODE` set + executable, and greps `^--- PASS: <name> ` for all three tests (a SKIP or zero-test run fails). A bad `WORLD_EXEC_NODE` is `t.Fatal`, not skip; the skip only fires locally when no node exists.
- No test writes outside t.TempDir(): node harness files, Chrome profile/HOME/TMPDIR, golden-update output all under TempDir; QUICKSTART/daemon.go are read-only reads.
- Timing: no `time.Sleep` anywhere in the new tests. Guards: 5 s, 10 s, 20 s, 30 s (all >= 2 s). Stimuli: aguiBudget 300 ms (assert status+frames only, never elapsed time), 5 s (empty-store arm, commit after first poll via aguiPollReads), 1 s Chrome. No oracle assumes one bounded run drains the log (nonempty-race arm asserts the exact delta list `[]`/`[6]` for a 300 ms run that starts at the max index, which is the intended shape). JS-off 1 s window is a labelled control, Chrome-only (not CI).
- Security: CSP literal is exactly the D2 string; `/workbench/live.js` sets nosniff, `text/javascript`, `default-src 'none'`, no-store, GET-only. Escape probe (temp test, deleted): hostile `"><script>..<img onerror>` in decision rows (effect/scope/requester/status/decidedBy/unavailable), graph relation/target labels, live head/hash, and `javascript:` hrefs in all href slots => page still has exactly ONE `<script` (the live.js tag), no `<img`, no `href="javascript`, all text entity-escaped, javascript: hrefs sanitized to `/workbench`. Script gadgets: all same-origin GETs are html (nosniff), `writeJSON` (object/array JSON, application/json, no side effects), `/v1/head` daemon-minted hash text/plain, agent card, live.js; matches design §8/§9, accepted residual (design item 9.1), not caused by this row beyond what the design priced in.
- End-to-end on a real daemon: I built `ailang-worldd`, ran the QUICKSTART section 11 commands VERBATIM (see deviation D). Then ran the real live.js in a node `vm` with node's real `fetch` against the live daemon (stub DOM): POST /agui/ lastIndex 0 -> backlog entry delivered -> GET /workbench re-fetch -> all five regions swapped, footer never replaced, status read "Live: checked ..., no new entries after entry 2; next check within 18 s", re-POST with lastIndex 2 after RUN_FINISHED, and a commit made 6 s in produced a second refresh. Works against real chunked SSE framing.

## BLOCKING findings
None.

## NON-BLOCKING findings
N1. `Live.Head` wiring is unpinned. host/daemon/workbench.go (`page.Live.Head = latest[0].EntryHash.String()`). Delete that line (mutant O16) or use the oldest hash (O13): `go test ./host/daemon -run '^TestWorkbenchLive'` stays GREEN; the page prints "(head )". The daemon test's `Contains(live region, c.Entry.EntryHash)` is satisfied by the row list, not the sentence. The pure template test pins the sentence with a literal Head. Fix: assert "(head <hash>)" in the daemon gap test.
N2. Decisions "as of entry N" is pinned only at N=0 (host/daemon/workbench_decisions_test.go): mutant `page.Decisions.Cursor = 0` survives (O3).
N3. live.js weak pins (all survive `TestLiveScriptPure`): (O1) the D5 positive-assertion status text ("Live: checked ..., no new entries after entry N; next check within 18 s") is never asserted (tests only match /paused:/ and /refresh failed/); (O2) the `response.headers.get("Retry-After")` wiring (header ignored -> wait still 1 s) because the lifecycle rig only ever sends "1" which equals the fallback; (O5) backoff reset on a terminal run.
N4. Selected-entry graph node `Href` is unpinned (O12: blank href survives TestWorkbenchGraphWiring); `aria-live="polite"` on the live section is unpinned (O14: removed, whole host/workbench package green).
N5. Deviation C-D1 side effect: when the home graph read fails the page renders BOTH "UNAVAILABLE: graph checked edges could not be read" and the stale "graph: select an entry" line (`{{if .Graph.Nodes}}{{else}}...` also fires), with `viewBox="0 0 640 0"` (Height 0). Cosmetic; the JS swap would install it too.
N6. CI step has `timeout-minutes: 2`; it includes the `-race` compile of host/daemon. Likely warm from earlier steps on the same runner (I measured 25 s wall locally cold-ish), but there is no slack proof on a loaded 2-vCPU runner.
N7. `gofmt -l` flags host/store/deadline_guard_test.go; it is ALREADY unformatted at b702d73 (base copy also listed); the new `"LogEntriesLatest"` entry is merely not realigned. Pre-existing, not a regression.
N8. Existing TestWorkbenchTimelinePaging is 27.6 s (base) vs 29.6 s (head) under `-race` (+7 %, extra per-render reads: 1 latest + <=40+40 approvals + <=15 graph). No CI risk seen; `-count=4` of the daemon workbench set exceeds a 110 s timeout purely because of that test at base speed.
N9. Chrome drill (AC6.2) could NOT be reproduced here: with WORLD_CHROME set, Chrome in my sandbox never issues the page GET in EITHER arm (js-off instrument control fails identically: "waiting for GET cursor 0"), so it is environmental; not counted as a code defect and not a CI gate (V27). The node-vm real-server run above covers the same lifecycle.

## Mutant table
All run as: apply to production code in my worktree, run ONLY the named test (package filter), observe, restore from cp backup, verify shasum (every row printed `restored-ok`; `git status` clean at end). Each red was an assertion failure (checked the FAIL line), not a build failure; JS mutants were syntactically valid.

Design-listed (§11) — 35 distinct mutants across all six milestones, all RED:
| Mutant | Test | Result |
|---|---|---|
| M-n dark theme drops --warnbg (M1) | TestTokensBothThemes | red |
| M-o empty-log branch renders nothing (M1) | TestRenderLiveRegion/-1 | red |
| M-k live rows as `<h3>entry N</h3>` (M1) | TestWorkbenchTimelineBound (daemon) and TestRenderLiveRegion | red (first try ran in wrong pkg = 0 tests; re-run in host/daemon red) |
| M-l svg xmlns (M1/M4) | TestRenderEmitsOnlyLocalLinks | red |
| M-m extra section after walk (M1) | TestProvenanceWalkStaysLast | red |
| footer loses role=status (M1, extra) | TestRenderQuietFooter | red |
| M-i LogEntriesLatest ASC (M2) | TestLogEntriesLatestNewestFirst | red |
| M-j cursor = latest+1 (M2) | TestWorkbenchLiveCursorIsResumable | red (404 arms) |
| M-ad empty cursor 0 (M2) | TestWorkbenchLiveCursorIsResumable/empty-live-entry-zero | red |
| M-h cursor = len-1 (variant of dense-scan, M2) | TestWorkbenchLiveCursorCrossesGap | red |
| M-a CSP +'unsafe-inline' (M3) | TestWorkbenchSecurityHeaders | red |
| M-b CSP drops connect-src (M3) | TestWorkbenchSecurityHeaders | red |
| M-c inline `<script>` in page (M3) | TestWorkbenchLiveScriptTag | red |
| M-d script in error template (M3) | TestWorkbenchErrorPageInert (400/404/500) | red |
| M-e route drops nosniff (M3) | TestWorkbenchLiveScriptRoute | red |
| M-f route accepts any method (M3) | TestWorkbenchLiveScriptRoute (405 arm) | red |
| M-g live.js added to isProtected (M3) | TestWorkbenchLiveScriptRoute | red |
| M-v splitter keeps unterminated tail (M3) x2 forms | TestLiveScriptPure | red |
| M-w accepts non-replace delta (M3) | TestLiveScriptPure | red |
| M-x schema drift v1->v2 (M3) | TestLiveScriptContract | red |
| M-y rename aria-label timeline (M3) | TestLiveScriptContract (empty/home/entry/object) | red |
| M-ae footer re-added to REGIONS (M3) | TestLiveScriptContract | red |
| M-af skip status re-assert (M3) | TestLiveScriptPure | red |
| M-ag throw / silent-skip on missing region (M3) 2 forms | TestLiveScriptPure | red, red |
| M-ah graph region omitted on ?object= (M3/M4) | TestLiveScriptContract/object | red |
| M-ai Retry-After undefined -> NaN (M3) | TestLiveScriptPure | red |
| M-q graph links unstored target (M4) | TestGraphUnavailableNodeHasNoLink | red |
| M-r row pitch 20 < box 28 (M4) | TestGraphNodesDisjointAndInBounds | red |
| M-s reversed object order (M4) | TestGraphDeterministic, TestGraphGolden | red |
| M-t golden truncated to 0 bytes (M4) | TestGraphGolden anti-vacuity arm | red |
| M-u grade label in graph (M4) | TestGraphCarriesNoGrade | red |
| graph height off (M4, extra) | TestGraphGolden + Disjoint | red |
| M-p unrendered field in DecisionRow (M5) | TestWorkbenchViewFieldsAllRender | red |
| M-p DecisionRow omitted from census (M5) | TestWorkbenchNewViewTypesInCensus | red |
| M-z newest decision not kept (M5) | TestFoldApprovals | red |
| M-aa walk bound +5 heads / summaries 20->200 (M5) | TestRecentApprovalsBounded | red, red |
| M-ab malformed -> 500 (M5) | TestRecentApprovalsMalformed (daemon: semantic/json/missing-request) | red |
| M-ac `<form><button>` in pane (M5) | TestWorkbenchDecisionsReadOnly (workbench AND daemon copy) | red, red |

Own mutants (derived from what the diff ships) — 16 run, 9 killed, 7 SURVIVED:
| Own mutant | Test run | Result |
|---|---|---|
| O4 drop `fresh` class on new rows | TestLiveScriptPure | red |
| O6 live limit 10 -> 11 | TestWorkbenchLiveCursorCrossesGap (limits spy) | red |
| O7 home graph newest 5 -> 6 | TestWorkbenchGraphWiring | red |
| O10 visible() ignores visibilityState | TestLiveScriptPure | red |
| O11 drop "graph: select an entry" | TestWorkbenchGraphWiring | red |
| O15 Decisions.Truncated unset | TestWorkbenchDecisionsPane | red |
| cursor = len(latest)-1 (gap) | TestWorkbenchLiveCursorCrossesGap | red |
| CSP literal drift / region list / approvals bound | (covered by M-a/b, M-ae, M-aa rows) | red |
| O1 status text "Live: checked..." -> "Live" | TestLiveScriptPure | SURVIVED |
| O2 Retry-After header ignored | TestLiveScriptPure | SURVIVED |
| O3 Decisions.Cursor = 0 | decisions tests | SURVIVED |
| O5 no backoff reset on terminal | TestLiveScriptPure | SURVIVED |
| O12 selected-entry graph Href blank | TestWorkbenchGraphWiring | SURVIVED |
| O13/O16 Live.Head wrong / unset | TestWorkbenchLive* | SURVIVED |
| O14 aria-live removed | whole host/workbench pkg | SURVIVED |
Note: O14 against the daemon package was inconclusive (full package exceeded my 120 s harness cap); the workbench-package run is conclusive for the template pin.

## Deviation verdicts
- (B) `<script src="/workbench/live.js" defer>` added in render.go although the M3 file list omitted it: ACCEPT. Design D2/AC3.2 require that exact tag and §8 lists render.go among edited files; M-c/M-d mutants and TestWorkbenchLiveScriptTag prove it is the single `<script`.
- (C-D1) home-page graph read errors -> sanitized UNAVAILABLE graph at 200, selected-entry errors stay 500: ACCEPT WITH CAVEAT. Design D6 is silent on graph errors; the behaviour is explicit, sanitized (tested: no "private" text), logged, and pinned by TestWorkbenchGraphStoreError, and every other store error (latest, approvals) still goes to the 500 page. Caveat N5 (double message). Slight asymmetry with D7's "I/O -> 5xx" rule; defensible because the page is useful without the graph.
- (C-D2) extra daemon read-only companion test for M-ac: ACCEPT. Verified it is load-bearing: the M-ac mutant is red in both packages.
- (D) QUICKSTART rehearsal stopped at /dev/tty in the executor sandbox: I COMPLETED the non-browser part. The session mint succeeds under a pty (`script -q /dev/null`, 'y' typed after a 3 s delay; typing immediately is swallowed by the pty echo). Verified verbatim: build, serve (`listening`), kill, mint, serve, the python generator (extracted from the doc by awk and run unmodified), genesis commit (selectedHead returned), `curl --fail /workbench` 200 with the D2 CSP, `data-live-cursor="0"`, second commit, `log get 1` = entryIndex 1. Doc caveat: `unset AILANG_REGISTRY_API_KEY` must be in the same shell as the mint (a fresh shell inherits it and mint refuses; the doc's single-terminal flow is fine). The visible-browser observation itself was not performed (see N9); my node-vm real-server run stands in for it.

## Gaps vs "reproduce before you assert"
- Chrome AC6.2 drill: could not reproduce (environmental, N9).
- Safari/Firefox behaviour: not measured by anyone (design says unmeasured).

---
# ROUND 2 (narrow) — head a5672da (r1 head 964d3b3 + one commit)

SCORE 92/100 — PASS. Zero blocking findings. (r1 91; +1 because the one thing r1 could not measure, AC6.2 on a real Chrome, is now measured green.)

## Change list verified
`git diff --stat 964d3b3..a5672da` = exactly one file, `host/daemon/workbench_live_chrome_test.go` (+7/-3): "home" removed from the created dirs, HOME removed from the stripped env keys and from the appended overrides, a comment added. XDG_CACHE_HOME/TMPDIR/TMP/TEMP are still relocated into t.TempDir(); `--user-data-dir=<TempDir>/profile` still present. Nothing else (no production file, no oracle line).

## (1) Drill, measured (worktree detached at a5672da)
`WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" go test -count=1 -v -run '^TestWorkbenchLiveChromeDrill$' ./host/daemon` => PASS (2.79 s):
- js-on: `trace GET /workbench`, `trace POST /agui/ lastIndex=0`, `committed entry 1 after POST cursor 0`, `trace GET /workbench` (re-fetch), `trace POST /agui/ lastIndex=1`. This is the full designed oracle (POST k, commit k+1, re-fetch GET, later POST k+1) observed from a real browser.
- js-off: page GET observed, zero POSTs for 1 s (instrument control).
So r1's N9 ("UNMEASURED, Chrome never GETs") is RESOLVED: the drill is measured green on this rig after the commit, and r1's red was the HOME relocation.

## (2) HOME hypothesis, independent test with controls
Python loopback server logging GETs (script in a separate dir, `python3 -I`), Chrome `--headless=new --no-first-run --disable-background-networking --user-data-dir=<tmp>` in each run; only HOME differed. Two repetitions each:
- HOME=real: requests logged = 2, 2 (page + favicon).
- HOME=<empty temp dir>: requests logged = 0, 0 (waited 8 s each).
Hypothesis CONFIRMED (control = real HOME, positive; the zero is not a dead-server artifact because the same server logged 2 in the sibling runs). I did not determine WHY (likely macOS keychain/Crashpad/updater lookup under HOME); only the effect is measured.

## (3) Does not relocating HOME matter?
- Writes outside TempDir: measured with a marker file around one Chrome run with real HOME. Chrome touched `~/Library/Application Support/Google/RLZ/RlzStore.plist` and `~/Library/Application Support/Google/Chrome/Crashpad/settings.dat` (the real-profile Crashpad/RLZ stores; `--user-data-dir` does not cover them). Other changed files in my `find` (gcloud logs, herdr log, spaces plist) were unrelated processes. It is Chrome's own housekeeping, not test data: no page content, DB, cookies or profile go there (those are in `--user-data-dir`, a TempDir). The test is env-gated (`WORLD_CHROME`), local-only, never in CI; the sprint rule targets test fixtures/DBs/profiles, and those remain inside t.TempDir(). Verdict: acceptable, non-blocking; it is a narrowing of the "nothing outside TempDir" stance that deserves the code comment it already has, plus a sentence in the QUICKSTART drill paragraph would be honest (suggestion only).
- Oracles: unweakened. The diff touches only env/dir setup, no assertion, wait, guard or trace line; the js-on oracle is intact and the drill went green with the full trace above. CI unaffected (skipped without WORLD_CHROME; step list unchanged).

## Carried-forward r1 non-blocking findings (unchanged, by name)
N1 `Live.Head` wiring unpinned; N2 Decisions "as of entry N" pinned only at N=0; N3 live.js weak pins (O1 status text, O2 Retry-After header wiring, O5 backoff reset); N4 selected-entry graph Href and `aria-live` unpinned; N5 home-graph-error shows both UNAVAILABLE and the stale "select an entry" line (+ viewBox height 0); N6 CI node step 2-min timeout includes -race compile; N7 pre-existing gofmt on deadline_guard_test.go; N8 TimelinePaging +7% under -race. N9 is superseded by this round (resolved). Deviation verdicts B, C-D1, C-D2, D unchanged.

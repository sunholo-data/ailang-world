# Sprint plan — `w-workbench-live-and-polish`

**Iteration 245 · queue row 149 · clause-5 (release row, D-WORLD-64)**

Design: [w-workbench-live-and-polish.md](w-workbench-live-and-polish.md), r2 narrow-refinement carve-out after r1/r2 BLOCKED; proceed without round 3. Base code `b702d73`, detached worktree `6fcfe57` plus design-doc commits only. Binding: CLAUDE.md, coding-standards.md S2/S3/S6/S7 and row-148 [AG-UI contract](../implemented/w-worldd-agui-event-stream.md) §1/D1–D6. ~2.1 executor days; C1 0.25d, C2 0.3d, C3 0.65d, C4 0.35d, C5 0.35d, C6 0.2d; all <=1 day.

Plan only: exactly two new planning files. Execution: one PR, six milestones M1-M6 and six commits C1-C6; executor never runs git. Controller commits from cumulative .snap/M1/ through .snap/M6/ snapshots. Host rendering/store/broker views and additive script route; no .ail, dependencies, npm, protocol changes, actions, row150/151 or mission edits.

The EXECUTOR does not run git, including add/commit/stash/checkout. The controller commits C1..C6 from cumulative `.snap/M<k>/` snapshots (M1–M6). Each boundary includes red-first evidence, named green test manifest, command rc/timing, mutation ledger and exact restore evidence. Suggested executor packing stays design §5: M1+M2; M3; M4; M5+M6, with a snapshot for each milestone. No implementation in this planner turn.

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export PATH=/opt/homebrew/bin:$PATH
```

Any test writing outside its own t.TempDir() is forbidden. DBs, fixture outputs, fake HOME/state/cache roots, harnesses, Chrome profiles and subprocess logs must be inside t.TempDir; never seed $HOME or ~/.ailang. Read-only pinned binary and fixtures may be read. Golden update writes only under the generating test’s TempDir and emits labelled base64/SHA-256 to stdout before cleanup; executor decodes captured bytes into tracked paths after the test returns (PD7). No test writes tracked goldens. Planner writes only its two deliverables; no git write command.

Baseline: `go vet ./...` planner rc=0, empty output. Controller runs full gates outside sandbox on pristine b702d73; Go baseline log was incomplete at inspection. No complete baseline gate result or incomplete-log quotation is claimed. Paths: `~/.ailang/state/world-iter245/base_ail.log`, `base_go.log`; controller attaches completed evidence. Known base flake: row-166 host/broker head/tail timeout under full parallel verification. Any `bind: operation not permitted` is **UNINFORMATIVE UNDER SANDBOX**, never pass or fail.

AC0 is required at every milestone: inside pre-existing daemon/workbench test functions/helpers only the two V3 CSP literals may change; only appended entries may change `workbenchViewTypes`. All new tests use new functions/files. Controller audits hunks against b702d73 and merge head. Store fixture exception is the separately numbered PD3, outside AC0’s packages.

Every newly authored feature oracle must fail on pre-milestone code before implementation; run tests with no-match checks and capture the precise red. Existing regression controls stay green (PD5). Missing API compilation is build-red only; mutation drills must fail their named assertions on compiling code. No test skip, harness timeout or sandbox denial is proof.

Numbered planner decisions:

**PD1.** Design M-p names an omitted census type but its named field test only examines listed types (V13, confirmed in render_test.go). Preserve TestWorkbenchViewFieldsAllRender for the design’s parenthetical mutant (listed type with an unrendered field). Add TestWorkbenchNewViewTypesInCensus in host/workbench/live_test.go: recursively enumerate the exported struct types reachable through Page.Live, Page.Graph and Page.Decisions (including slice/array/pointer row types), require each in workbenchViewTypes, and fail on an empty census. Drill both omission and unrendered-field forms separately for every new view/row type at C1, C4 and C5. This corrects the omitted-type mapping rather than changing the existing test body.

**PD2.** Design M-ab names TestRecentApprovalsMalformed, but a broker-only test cannot kill a daemon 500 response mutant. Define that exact test name in BOTH host/broker/approvals_view_test.go (ErrApprovalChainMalformed for wrong semantic id/undecodable payload) and host/daemon/workbench_decisions_test.go (same malformed chain yields 200 plus UNAVAILABLE and the offending ref). Run both exact package commands; the daemon assertion is the M-ab kill.

**PD3.** Adding exported Store.LogEntriesLatest makes TestStoreDeadlineGuardDynamic fail its reflected-method/fixture exact census. Design sections 4/8 omit this dependency. Permit exactly one entry in host/store/deadline_guard_test.go guardFixtures: "LogEntriesLatest": simple(1). No allow-list or guard weakening. AC0 applies to daemon/workbench tests, so this store fixture addition does not violate it. Capture the missing-fixture red after adding the method and before adding the entry.

**PD4.** AC6.2’s 3 s JS-off observation exceeds the brief’s approximately 2 s total deliberate-wait ceiling. Keep this explicitly non-load-bearing control but reduce its window to 1 s AFTER the recorded page GET; use a timer/channel select, not Sleep. JS-on uses channel handshakes and the designed 1 s run budget, with no extra deliberate waits. Every stuck guard stays >=2 s (20 s Chrome waits, 30 s process context). This is a timing-only correction, not a weaker behavioral oracle; the JS-on request trace remains mandatory.

**PD5.** Several ACs ask existing tests to stay green, not become red, and AC6.1 has no named test. Preserve pre-existing regression bodies under AC0. Add new TestRenderLiveRegion/contract/census assertions for red-first feature proof; record existing green tests as controls. Add binary-free TestWorkbenchLiveQuickstartControl in host/daemon/workbench_live_chrome_test.go before docs edits; require the new heading, runnable URL, four pane descriptions, JS-off explanation, 16-slot/hidden-tab posture, commit payload construction and route doc comment. It must fail on C5 docs. This documentation control is instrument-health only; the actual verbatim rehearsal and Chrome drill still discharge S7. New TestProvenanceWalkStaysLast includes existence/order of all new sections before the walk, and TestWorkbenchErrorPageInert includes a success-script positive control, so their unchanged base invariants do not pass red-first vacuously. The new C6 Chrome drill checks the shared absent-docs precondition before launching; its red-first failure is documentation, while its green acceptance still requires the real trace.

**PD6.** Design section 4 locates the CI step at M3 while the brief asks to specify any M6 workflow edit. Install the sole new verbose PASS-loop step in C3 alongside the node harness; M6 makes NO workflow edit. Render the graph and decisions region shells unconditionally in C1 (no graph/approval data logic yet) so the five-selector per-variant contract can pass at C3. Populate them in C4/C5 without changing selectors. No GraphView/DecisionsView fields are introduced until their respective census and rendering land.

**PD7.** D6’s golden-update flag convention would write repository testdata, conflicting with the user’s absolute TempDir rule. Generate under t.TempDir and emit labelled base64 plus SHA-256 before cleanup; executor decodes the captured stdout into exactly the two named golden files after test exit. Copying from TempDir after return would fail because cleanup removes it. Normal tests never update fixtures, and the flag does not waive missing/empty golden controls on ordinary runs.

**PD8.** Design §5 repeats full build/plain/race gates at every milestone, while this brief requires the narrowest package gates and flags the ~600 s CI Go-leg budget. Use the listed narrow package tests plus go vet ./... at C1-C6; retain the full verify_ail/verify_go and explicit four-package race gate at final handoff. No script, runner timeout or production bound change; report focused and full gate durations. This is a gate-scheduling refinement, all design-wide checks remain mandatory.

Frozen surface (no delivered hunk; temporary mutations restored before snapshot):

- `GET /v1/health`
- `GET /v1/head`
- `GET /v1/worlds/{ref}`
- `GET /v1/objects/{ref}`
- `GET /v1/objects/by-semantic-id/{name...}`
- `GET /v1/log/{index}`
- `GET /v1/log`
- `GET /v1/registry/{name...}`
- `POST /v1/commit`
- `GET /v1/receipts/{id}`
- `host/daemon/handlers.go and all handlers for the ten /v1 patterns (no hunk)`
- `host/daemon/daemon.go isProtected body`
- `host/daemon/workbench.go acceptedWorkbenchKeys and supportedWorkbenchQuery`
- `host/agui/* (including schema and stream goldens)`
- `host/daemon/agui.go`
- `every D7 transport constant and TestBoundedWaitsAndBodyLimit`
- `host/daemon/workbench.go workbenchErrorTemplate`
- `tools/launchd/*`
- `design_docs/world-mission*.md`

Only the paths named below are executor scope. No .ail, package.json/npm or new dependency. An AILANG package cannot set HTTP CSP, own the route/embed/browser asset or render host HTML; existing log/broker semantics remain authoritative (D9). Workbench production stays pure, no net/http/store/broker imports; effectful reads and HTTP stay at the boundary.

CI and local test policy:

Node is mandatory in the CI Go job using existing `WORLD_EXEC_NODE`. `TestLiveScriptPure` can `t.Skip` locally only when neither configured node nor PATH node exists; a configured invalid executable fails. Binary-free `TestLiveScriptAssetNonEmpty` and `TestLiveScriptContract` always run and refuse empty extraction/assets. The verbose step requires all three exact top-level PASS lines, so absent runtime, skipped harness and zero tests fail CI (row-160 lesson). The AILANG-only job does not execute these tests; the required Go job supplies JS runtime evidence. Local C3 and final acceptance also require node PASS.

`TestWorkbenchLiveChromeDrill` skips only when WORLD_CHROME is unset in ordinary runs/CI. It is local-only: controller/executor must set a valid local executable, run the command below outside socket sandbox and retain its PASS transcript for C6. Invalid configured Chrome fails. No CI Chrome gate.

**Exact workflow edit: C3 only**, `.github/workflows/ci.yml`, job `go-verify`, new step **Row 149 workbench live script tests, verbose (node required)**, immediately after **Row 140 Workspace.Exec real-srt tests, verbose (linux bubblewrap)** and before **go build + test gate**. Existing WORLD_EXEC_NODE export is in scope. One step, timeout-minutes: 2, no other workflow edits. M6 edits no workflow. Run block:

```sh
set -euo pipefail
: "${WORLD_EXEC_NODE:?WORLD_EXEC_NODE is required}"
[ -x "$WORLD_EXEC_NODE" ]
log="$RUNNER_TEMP/workbench-live-verbose.log"
go test -race -count=1 -timeout 60s -v ./host/workbench ./host/daemon -run '^(TestLiveScriptPure|TestLiveScriptAssetNonEmpty|TestLiveScriptContract)$' >"$log" 2>&1
cat "$log"
for t in TestLiveScriptPure TestLiveScriptAssetNonEmpty TestLiveScriptContract; do
  grep -q -- "^--- PASS: $t " "$log" || { echo "$t did not PASS (SKIP or zero tests)"; exit 1; }
done
```

Timing: CI Go leg runs -race on small loaded linux runners with ~600 s process budget. Use existing per-instance aguiBudget/aguiTick seams, mostly 300 ms/5–10 ms; no production bound edits. Empty-store test cancels after delivery instead of spending its 5 s allowance. No new test sleeps or deliberately waits more than approximately 2 s total; Chrome JS-off control reduced to 1 s. All stuck guards >=2 s. Use handshakes/cancellation/joins and report focused test and daemon race wall time; do not increase CI budgets.

Row-148 resume rule: Row 148 bounded runs do not drain long logs. Every all-entry oracle must re-POST at RUN_FINISHED.result.lastIndex until the independently expected maximum; discarded tails never advance the cursor. No single-run completeness assumption.

## C1 / M1 — Tokens, live/footer regions and view census

Production estimate/budget: ~120 added LOC, including JS/CSS/templates; ≤~250. One controller commit.

Exact files:

- `host/workbench/render.go`
- `host/workbench/live_test.go`
- `host/workbench/render_test.go`

Tests FIRST: all new feature assertions below must be red on the pre-milestone snapshot; existing controls green as PD5 specifies.

- **AC1.1 — TestTokensBothThemes:** `TestTokensBothThemes`: parse the `<style>` block. `:root` and the dark `:root` declare the same 11 token names (V21). Null control: the test fails if it finds 0 tokens in either. Assert exactly --bg --card --ink --mut --line --acc --accbg --ok --okbg --warn --warnbg in both themes, including light/dark bg values from the mockups; no external stylesheet.
- **AC1.2 — TestRenderLiveRegion:** `TestRenderLiveRegion`: `Live{Cursor:7, Recent:[7,6]}` renders `data-live-cursor="7"`, rows newest first, no `<h3>entry `, and inside the live section "Log at entry 7 (head …)". `Cursor:-1` renders `data-live-cursor="-1"` and "The log is empty" inside the live section.
- **AC1.3 — TestRenderQuietFooter:** `TestRenderQuietFooter`: both states render a footer that is outside `<main>` and holds `<span data-live-status>Live updates: off`. The log-position sentence is **not** in the footer (it is in the swapped live section, D5).
- **AC1.4 — TestProvenanceWalkStaysLast:** `TestProvenanceWalkStaysLast`: the rendered page contains exactly one `"</section>\n</main>"`, and it closes the provenance walk. Red-first positive control: require all three NEW sections live/world graph/decisions exist exactly once before the provenance walk, so this new test fails on base markup even though the old walk is already last.
- **AC1.5 — TestWorkbenchNewViewTypesInCensus, TestWorkbenchViewFieldsAllRender:** `workbenchViewTypes` gains `LiveView`. `TestWorkbenchViewFieldsAllRender` passes. Append LiveView and every new struct row type it holds, not only the top view. PD1 checks actual census membership recursively.
- **AC1.6 — TestRenderLiveRegion plus all existing host/workbench tests (green controls):** All pre-existing `host/workbench` tests pass, with AC0. AC0 allows only census appends and (at C3) the two CSP literals in pre-existing daemon/workbench tests. Preserve grade bytes, closed links, timeline headings/counts and provenance markup. Render all five swap regions on every 200 variant, using empty graph/decisions shells until C4/C5 (PD6).

Implementation only after red proof: Add eleven light/dark tokens, LiveView and rows, newest-first live markup, positive log-position sentence, no-JS status footer after main, and graph/decisions shells before provenance. Leave grade markup and existing markup contracts intact.

Acceptance commands (environment above; named PASS required):

```sh
go test -count=1 -v ./host/workbench
go vet ./...
```

Mutation drills (apply one form at a time and require the named assertion red):

| Mutant | Isolated edit | Exact test / assertion |
|---|---|---|
| M-m | new section placed after the provenance walk | TestProvenanceWalkStaysLast — `TestProvenanceWalkStaysLast` |
| M-n | dark theme missing a token | TestTokensBothThemes — `TestTokensBothThemes` |
| M-o | live section's empty-log branch renders nothing | TestRenderLiveRegion — `TestRenderLiveRegion` (`Cursor:-1` arm) |
| M-p | any of the three new view types (`LiveView` AC1.5, `GraphView` AC4.6, `DecisionsView` AC5.4, or their row types) not appended to the census | TestWorkbenchNewViewTypesInCensus / TestWorkbenchViewFieldsAllRender — `TestWorkbenchViewFieldsAllRender` (V13: an unlisted type's unrendered field is not caught, so each AC's append is what arms it; the mutant is: append the type, leave one field unrendered) PD1: omitted-type form reds the new membership test; listed-but-unrendered form reds the unchanged field test. Repeat every new type at C1/C4/C5. |

Exact mutation commands:

```sh
# M-m
go test -count=1 -race -v ./host/workbench -run '^(TestProvenanceWalkStaysLast)$'
# M-n
go test -count=1 -race -v ./host/workbench -run '^(TestTokensBothThemes)$'
# M-o
go test -count=1 -race -v ./host/workbench -run '^(TestRenderLiveRegion)$'
# M-p
go test -count=1 -race -v ./host/workbench -run '^(TestWorkbenchNewViewTypesInCensus|TestWorkbenchViewFieldsAllRender)$'
```

**Controller commit message:** `feat(workbench): add theme tokens and live status regions`

## C2 / M2 — Newest-first log read and live cursor wiring

Production estimate/budget: ~90 added LOC, including JS/CSS/templates; ≤~250. One controller commit.

Exact files:

- `host/store/log_latest.go`
- `host/store/log_latest_test.go`
- `host/store/deadline_guard_test.go`
- `host/daemon/daemon.go`
- `host/daemon/workbench.go`
- `host/daemon/workbench_live_test.go`

Tests FIRST: all new feature assertions below must be red on the pre-milestone snapshot; existing controls green as PD5 specifies.

- **AC2.1 — TestLogEntriesLatestNewestFirst:** `TestLogEntriesLatestNewestFirst` (entries 0, 1, 5): limit 10 gives [5,1,0], limit 1 gives [5], an empty store gives []. Each row deep-equals `GetLogEntry`. Reuse scanLogAfterRow without editing log_after.go or store.go; SQL is ORDER BY entry_index DESC LIMIT ?. One indexed query at limit 10 feeds both cursor/head and recent rows. No dense scan.
- **AC2.2 — TestLogEntriesLatestBounds, TestStoreDeadlineGuardDynamic:** Limit 0 or 501 gives `InvalidLimitError`. No deadline gives `ErrNoDeadline`. A quarantined store is refused. Also test canceled/expired contexts. Add only the PD3 fixture entry; run the full store census and bounds tests.
- **AC2.3 — TestWorkbenchLiveCursorCrossesGap:** `TestWorkbenchLiveCursorCrossesGap`: a daemon with commits at 0, 1, 5 (REST, as row 148 AC3.14) renders `data-live-cursor="5"`. An empty store renders `-1`. Construct gaps through the existing authenticated REST commit fixture; require 0/1/5 commit statuses 200 before asserting cursor. Cover newest 10 with >10 entries and unchanged oldest-first timeline pagination.
- **AC2.4 — TestWorkbenchLiveCursorIsResumable:** `TestWorkbenchLiveCursorIsResumable`: POST `/agui/` with the rendered cursor as `state.lastIndex` gets 200, not 404. It uses `aguiBudget` 300 ms (stimulus) and a 10 s client context (stuck guard). **Empty-store arm (r1):** on an empty store, the cursor read from the rendered `data-live-cursor` (`-1`) is POSTed over an `httptest.NewServer`; the response is 200; after the run's first poll (the `aguiPollReads` pattern of row 148's `TestAGUISeesDirectStoreCommit`) the test commits entry 0 and asserts a `CUSTOM` with `"entryIndex":0` arrives on that same stream (V29a measured exactly this). Stuck guard 5 s per wait. Set per-instance aguiTick=5–10 ms. Nonempty run budget=300 ms, client guard=10 s. Empty-store arm retains 5 s budget as a safety allowance, first-poll channel handshake, 5 s wait guards, then cancel/join immediately after complete CUSTOM+STATE_DELTA for entry 0; never wait out the budget. Both arms use a real httptest.NewServer. For any longer log oracle, repeatedly resume at RUN_FINISHED.result.lastIndex; never require a bounded run to drain it. Capture status before frames; include a render-to-POST commit race arm.
- **AC2.5 — TestWorkbenchLiveStoreError:** `TestWorkbenchLiveStoreError`: a seam failing `LogEntriesLatest` gives the existing constant 500 HTML, with no store text. Assert no secret sentinel in HTML and existing constant error page; no script on error pages.

Implementation only after red proof: Add LogEntriesLatest with quarantine/deadline/limit checks and DESC indexed SQL; reuse scanLogAfterRow. Add readStore method; call once at limit 10 to populate Live cursor/head/recent on all success variants. Empty cursor -1. Add PD3 fixture only.

Acceptance commands (environment above; named PASS required):

```sh
go test -count=1 -v ./host/store -run '^(TestLogEntriesLatestNewestFirst|TestLogEntriesLatestBounds|TestStoreDeadlineGuardDynamic)$'
go test -count=1 -v ./host/daemon -run '^TestWorkbenchLive(CursorCrossesGap|CursorIsResumable|StoreError)$'
go test -count=1 ./host/store ./host/daemon ./host/workbench
go vet ./...
```

Mutation drills (apply one form at a time and require the named assertion red):

| Mutant | Isolated edit | Exact test / assertion |
|---|---|---|
| M-h | cursor via dense `GetLogEntry` scan | TestWorkbenchLiveCursorCrossesGap — `TestWorkbenchLiveCursorCrossesGap` |
| M-i | `LogEntriesLatest` ascending | TestLogEntriesLatestNewestFirst — `TestLogEntriesLatestNewestFirst` |
| M-j | cursor rendered as latest+1 | TestWorkbenchLiveCursorIsResumable — `TestWorkbenchLiveCursorIsResumable` (404) |
| M-k | live rows rendered as `<h3>entry N</h3>` | TestWorkbenchTimelineBound — `TestWorkbenchTimelineBound` (measured, V12 m4) |
| M-ad | empty-log cursor rendered as `0` instead of `-1` | TestWorkbenchLiveCursorIsResumable — `TestWorkbenchLiveCursorIsResumable` (empty-store arm: status 404 ≠ 200; premise measured, V29 control) |

Exact mutation commands:

```sh
# M-h
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchLiveCursorCrossesGap)$'
# M-i
go test -count=1 -race -v ./host/store -run '^(TestLogEntriesLatestNewestFirst)$'
# M-j
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchLiveCursorIsResumable)$'
# M-k
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchTimelineBound)$'
# M-ad
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchLiveCursorIsResumable)$'
```

**Controller commit message:** `feat(store): wire newest-first log entries into workbench`

## C3 / M3 — Embedded live script, route and exact CSP

Production estimate/budget: ~230 added LOC, including JS/CSS/templates; ≤~250. One controller commit.

Exact files:

- `host/workbench/live.js`
- `host/workbench/embed.go`
- `host/workbench/live_script_test.go`
- `host/daemon/workbench.go`
- `host/daemon/daemon.go`
- `host/daemon/workbench_test.go`
- `host/daemon/workbench_live_test.go`
- `.github/workflows/ci.yml`

Tests FIRST: all new feature assertions below must be red on the pre-milestone snapshot; existing controls green as PD5 specifies.

- **AC3.1 — TestWorkbenchLiveScriptRoute plus TestWorkbenchSecurityHeaders, TestWorkbenchReadDeadline and all assertWorkbenchSecurityHeaders callers (green controls after the two literal edits):** The two V3 literals become the D2 literal. `TestWorkbenchSecurityHeaders` and `assertWorkbenchSecurityHeaders` callers pass. Pin exact D2 CSP in both pre-existing literal locations; never derive test expectation from production symbol. No other pre-existing test body edits.
- **AC3.2 — TestWorkbenchLiveScriptTag:** `TestWorkbenchLiveScriptTag`: the success page contains exactly one `<script` and it is `<script src="/workbench/live.js" defer></script>`. No `on[a-z]+=` attribute appears.
- **AC3.3 — TestWorkbenchErrorPageInert:** `TestWorkbenchErrorPageInert`: 400/404/500 workbench pages contain no `<script`. Red-first positive control in this new test: first require the success page to carry the new external script, then verify each 400/404/500 response remains inert. Base code fails the success control; after implementation M-d fails only the error-page assertion.
- **AC3.4 — TestWorkbenchLiveScriptRoute:** `TestWorkbenchLiveScriptRoute`: GET returns 200, the D2 headers (literal values), and a body byte-equal to `workbench.LiveScript` with `len > 0`. POST gets 405. `isProtected` is equal (both `false`) for `GET /workbench/live.js` and `GET /workbench` (V30). Pin Content-Type text/javascript; charset=utf-8, nosniff, no-store, and script-route CSP default-src 'none'. Assert false/false GET posture directly with protected POST /v1/commit positive control. Keep isProtected byte-unchanged.
- **AC3.5 — TestLiveScriptContract:** `TestLiveScriptContract` (binary-free): from the embedded script, extract the `STATE_SCHEMA` literal and assert it equals `agui.StateSchema`. Extract the one-line `REGIONS` JSON array; assert it has exactly 5 selectors and none names `footer` or `live status` (D1 step 4). **Table-driven per variant (r2):** for each row of {home on an empty store, home on a seeded store, `?from=0&entry=0`, `?object=<stored ref>`} (the V32 variants), the daemon page is 200 and every selector's `aria-label` occurs exactly once. Extract the POST path and assert it equals `/agui/`. Null control: the test fails if any extraction finds nothing or the variant table is empty. Test lives in host/daemon/workbench_live_test.go and always runs without node/Chrome. Extract and pin all FIVE exact D1 selectors (not only array length), schema and /agui/ POST path; no successful empty extraction. Include footer exclusion and variant counts. Graph/decisions object variant shells must remain through C4/C5.
- **AC3.6 — TestLiveScriptPure (node) plus TestLiveScriptAssetNonEmpty (binary-free companion):** `TestLiveScriptPure` (node): runs `node` (from `WORLD_EXEC_NODE` or `PATH`) on the script plus a harness (the script exports `splitFrames`/`applyDelta`/`nextBackoff`/`applySwap`/ `renderStatus` only when `module` exists, and starts only when `document` exists). It feeds `host/agui/testdata/stream_fixture.golden` split at **every** byte offset into two chunks. The event list equals the whole-buffer parse, and the cursor never exceeds the last fully delivered entry. `applyDelta` with an `add` op throws. Backoff is 1, 2, 4 … capped at 30. A missing or unparsable `Retry-After` gives the first step, 1 (D8). **Footer arm (r1):** the harness builds a minimal stub document (plain objects with `querySelector`/`replaceWith`/`textContent`; no DOM library), rewrites the status span to the script's live status, then calls `applySwap` with a fresh document whose footer holds the server-rendered `Live updates: off. Reload to refresh.` span, and asserts the span still carries `renderStatus(state)` and that the footer stub's `replaceWith` was never called. A second sub-arm makes the re-assert itself load-bearing: the span holds the server "off" text, script state is `paused`, and after `applySwap` the span must read `renderStatus(paused)`. **Missing-region sub-arm (r2):** a stub fresh document that omits `section[aria-label="decisions"]` (and, separately, an old document that omits it): `applySwap` does not throw, swaps the other four regions, does not touch the missing one, and the span reads "paused: page layout changed, reload required". Locally the test SKIPs without node. In CI a new verbose step runs it and fails unless `--- PASS: TestLiveScriptPure` appears (the ci.yml PASS-loop pattern). Test lives in host/workbench/live_script_test.go; TempDir harness and read-only row-148 golden. Feed Buffer bytes through streaming TextDecoder at every cut (do not treat byte cuts as character cuts); assert complete frames only, last complete delta cursor, no duplication/loss against the whole-buffer event list, and discarded severed tails. Hand-author a terminal RUN_FINISHED fixture because the row-148 golden excludes terminal frames; verify next request uses its lastIndex. Use fake timers/fetch/AbortController to assert one active visible-tab POST, hidden abort/resume, 300 ms debounce with one refresh in flight, non-200 refresh keeps nodes, 503/network/RUN_ERROR paused status, immediate successful terminal re-POST, and no re-POST after missing-region layout pause. Keep footer live and paused sub-arms and both one-sided missing-region sub-arms. Module exports only the named pure seams when module exists; no startup without document. Node command timeout=30 s. Missing node skips locally only; an explicit invalid WORLD_EXEC_NODE fails. Companion TestLiveScriptAssetNonEmpty always asserts embedded nonzero asset and all export seams; TestLiveScriptContract always binds the real daemon variants. CI must require their PASS AND the node test PASS (see step below).

Implementation only after red proof: Embed live.js as LiveScript in embed.go. Add GET /workbench/live.js handleWorkbenchScript and exact success script tag. Widen CSP by only script-src self and connect-src self; edit exactly the two test literals. Follow D1/D5/D8 with fetch+ReadableStream, complete-frame state, server-HTML swaps, footer reassertion, visible-tab abort, exponential backoff and layout pause. Add sole CI step below.

Acceptance commands (environment above; named PASS required):

```sh
go test -count=1 -v ./host/daemon -run '^(TestWorkbenchSecurityHeaders|TestWorkbenchReadDeadline|TestWorkbenchLiveScriptTag|TestWorkbenchErrorPageInert|TestWorkbenchLiveScriptRoute|TestLiveScriptContract|TestFrozenV1RouteTableMatchesMux)$'
go test -count=1 -v ./host/workbench -run '^(TestLiveScriptPure|TestLiveScriptAssetNonEmpty)$'
set -euo pipefail
: "${WORLD_EXEC_NODE:?WORLD_EXEC_NODE is required}"
[ -x "$WORLD_EXEC_NODE" ]
log="$(mktemp)"
trap 'rm -f "$log"' EXIT
go test -race -count=1 -timeout 60s -v ./host/workbench ./host/daemon -run '^(TestLiveScriptPure|TestLiveScriptAssetNonEmpty|TestLiveScriptContract)$' >"$log" 2>&1
cat "$log"
for t in TestLiveScriptPure TestLiveScriptAssetNonEmpty TestLiveScriptContract; do
  grep -q -- "^--- PASS: $t " "$log" || { echo "$t did not PASS (SKIP or zero tests)"; exit 1; }
done
go test -count=1 ./host/workbench ./host/daemon
go vet ./...
```

Mutation drills (apply one form at a time and require the named assertion red):

| Mutant | Isolated edit | Exact test / assertion |
|---|---|---|
| M-a | CSP gains `'unsafe-inline'` in `script-src` | TestWorkbenchSecurityHeaders — `TestWorkbenchSecurityHeaders` (literal) |
| M-b | CSP drops `connect-src 'self'` | TestWorkbenchSecurityHeaders / TestWorkbenchRefusalBranches / TestWorkbenchReadDeadline — `TestWorkbenchSecurityHeaders`, `assertWorkbenchSecurityHeaders` callers |
| M-c | an inline `<script>` added to the page | TestWorkbenchLiveScriptTag — `TestWorkbenchLiveScriptTag` (count) |
| M-d | script tag added to the error template | TestWorkbenchErrorPageInert — `TestWorkbenchErrorPageInert` |
| M-e | script route drops `nosniff` or serves `text/plain` | TestWorkbenchLiveScriptRoute — `TestWorkbenchLiveScriptRoute` |
| M-f | script route accepts POST | TestWorkbenchLiveScriptRoute — `TestWorkbenchLiveScriptRoute` (405 arm) |
| M-g | `/workbench/live.js` added to `isProtected` | TestWorkbenchLiveScriptRoute — `TestWorkbenchLiveScriptRoute` (posture arm) |
| M-v | script splitter parses the unterminated tail | TestLiveScriptPure — `TestLiveScriptPure` (byte-offset arm) |
| M-w | script accepts a non-`replace` delta | TestLiveScriptPure — `TestLiveScriptPure` (`add` arm) |
| M-x | script's schema string drifts from `agui.StateSchema` | TestLiveScriptContract — `TestLiveScriptContract` |
| M-y | template renames `aria-label="timeline"` | TestLiveScriptContract — `TestLiveScriptContract` (selector arm) |
| M-ae | `footer[aria-label="live status"]` re-added to `REGIONS` | TestLiveScriptContract — `TestLiveScriptContract` (footer-exclusion arm) |
| M-af | `applySwap` skips the status re-assert | TestLiveScriptPure — `TestLiveScriptPure` (footer arm, second sub-arm: span still reads the server "off" text) |
| M-ag | `applySwap` throws on, or silently skips, a one-sided missing region | TestLiveScriptPure — `TestLiveScriptPure` (missing-region sub-arm: throw, or span not "paused: page layout changed…") |
| M-ai | `Retry-After` missing → wait 0 or NaN | TestLiveScriptPure — `TestLiveScriptPure` (backoff arm: `nextBackoff(undefined)` and `nextBackoff("x")` = 1) |

Exact mutation commands:

```sh
# M-a
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchSecurityHeaders)$'
# M-b
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchSecurityHeaders|TestWorkbenchRefusalBranches|TestWorkbenchReadDeadline)$'
# M-c
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchLiveScriptTag)$'
# M-d
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchErrorPageInert)$'
# M-e
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchLiveScriptRoute)$'
# M-f
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchLiveScriptRoute)$'
# M-g
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchLiveScriptRoute)$'
# M-v
go test -count=1 -race -v ./host/workbench -run '^(TestLiveScriptPure)$'
# M-w
go test -count=1 -race -v ./host/workbench -run '^(TestLiveScriptPure)$'
# M-x
go test -count=1 -race -v ./host/daemon -run '^(TestLiveScriptContract)$'
# M-y
go test -count=1 -race -v ./host/daemon -run '^(TestLiveScriptContract)$'
# M-ae
go test -count=1 -race -v ./host/daemon -run '^(TestLiveScriptContract)$'
# M-af
go test -count=1 -race -v ./host/workbench -run '^(TestLiveScriptPure)$'
# M-ag
go test -count=1 -race -v ./host/workbench -run '^(TestLiveScriptPure)$'
# M-ai
go test -count=1 -race -v ./host/workbench -run '^(TestLiveScriptPure)$'
```

**Controller commit message:** `feat(workbench): follow AG-UI with an embedded live script`

## C4 / M4 — Deterministic server-rendered SVG graph

Golden generation command (PD7, separate from normal acceptance):

```sh
go test -count=1 -v ./host/workbench -run '^TestGraphGolden$' -update-workbench-golden
```

Test writes only in t.TempDir, logs WORLD_GRAPH_GOLDEN <basename> <sha256> <base64> before cleanup. Executor captures the two stdout markers, decodes with strict base64 validation, verifies SHA-256 and writes only host/workbench/testdata/graph_home.golden and graph_selected.golden. Then rerun without update flag; byte equality, node count and nonzero-svg anti-vacuity must PASS. Never compare the generated file to itself as ordinary acceptance.

Production estimate/budget: ~200 added LOC, including JS/CSS/templates; ≤~250. One controller commit.

Exact files:

- `host/workbench/graph.go`
- `host/workbench/graph_test.go`
- `host/workbench/testdata/graph_home.golden`
- `host/workbench/testdata/graph_selected.golden`
- `host/workbench/render.go`
- `host/workbench/render_test.go`
- `host/daemon/workbench.go`
- `host/daemon/workbench_graph_test.go`

Tests FIRST: all new feature assertions below must be red on the pre-milestone snapshot; existing controls green as PD5 specifies.

- **AC4.1 — TestGraphGolden:** `TestGraphGolden`: the two fixtures render byte-equal to their goldens, which are generated by the `-update-workbench-golden` flag. **Anti-vacuity:** it fails if a golden is 0 bytes, lacks `<svg`, or has a `<rect` count ≠ the fixture's node count. PD7: generate only through -update-workbench-golden, writing the two produced files under the generating test’s t.TempDir. Before cleanup, emit labelled base64 bytes plus SHA-256 to stdout; executor captures and decodes them into the two allowed tracked paths after the test exits. Normal tests only read goldens; missing/zero goldens are red before generation. No test writes tracked fixtures or disables cleanup.
- **AC4.2 — TestGraphDeterministic:** `TestGraphDeterministic`: two layouts and renders are byte-equal, and so is a layout of a deep copy. Repeat layout/render enough times (e.g. 32) with shared/distinct targets and deep copies, and compare first-appearance node order explicitly to an independently authored expectation. Repetition alone is not a map-order oracle.
- **AC4.3 — TestGraphNodesDisjointAndInBounds:** `TestGraphNodesDisjointAndInBounds`: for fixtures of 1–5 entries with shared and distinct targets, every box lies inside the viewBox and no two boxes intersect. Pin box 280x28, pitch 36, margin 16, x=16/344, width 640, height 32+36*max(rows); assert integer deterministic geometry, no overlap and bounds.
- **AC4.4 — TestGraphUnavailableNodeHasNoLink:** `TestGraphUnavailableNodeHasNoLink`: an unstored target renders `UNAVAILABLE`, and the graph's `<a` count equals the count of available nodes. Anchor count is AVAILABLE nodes across BOTH entry and object columns (available entries plus checked available targets); assert unavailable target label exists and its specific node has no anchor.
- **AC4.5 — TestGraphCarriesNoGrade:** `TestGraphCarriesNoGrade`: the graph contains none of `PROVEN|TESTED|ATTESTED|CLAIMED`.
- **AC4.6 — TestWorkbenchGraphWiring, TestWorkbenchNewViewTypesInCensus, TestWorkbenchViewFieldsAllRender and TestRenderEmitsOnlyLocalLinks (green controls):** `TestWorkbenchGraphWiring` (daemon): the home page graph shows the newest ≤ 5 entries. The selected-entry page graph shows exactly that entry's 3 edges. Pre-existing `TestRenderEmitsOnlyLocalLinks` and the body-wide count tests pass (V11). **Census (r2):** `GraphView` (and every struct type it holds, e.g. its node and edge row types) is appended to `workbenchViewTypes`, and `TestWorkbenchViewFieldsAllRender` continues to pass (V13). Home takes newest <=5 live entries, <=15 checkedEdge GetObject reads; selected entry reuses its three checked edges; object variant says graph: select an entry. Only transitionFn/interpreter/transitionRef; exclude stateRoot/committedBy/evidence. Append GraphView and all graph row types to census (PD1); run existing body-wide link/count tests untouched.

Implementation only after red proof: Add pure GraphEntry/LayoutGraph/GraphView with integer geometry and stable first-appearance ordering. Render SVG without xmlns, grades or absolute URLs; links use workbenchHref. Populate home/selected views using existing checked edges; object region remains a select-entry placeholder. Generate goldens using TempDir output and labelled base64/SHA-256 stdout, then materialize the captured bytes (PD7).

Acceptance commands (environment above; named PASS required):

```sh
go test -count=1 -v ./host/workbench -run '^(TestGraph.*|TestWorkbenchNewViewTypesInCensus|TestWorkbenchViewFieldsAllRender|TestRenderEmitsOnlyLocalLinks)$'
go test -count=1 -v ./host/daemon -run '^(TestWorkbenchGraphWiring|TestWorkbenchTimeline.*|TestLiveScriptContract)$'
go test -count=1 ./host/workbench ./host/daemon
go vet ./...
```

Mutation drills (apply one form at a time and require the named assertion red):

| Mutant | Isolated edit | Exact test / assertion |
|---|---|---|
| M-l | `<svg xmlns="http://www.w3.org/2000/svg">` | TestRenderEmitsOnlyLocalLinks — `TestRenderEmitsOnlyLocalLinks` (measured, V12 m1) |
| M-q | graph links an unstored target | TestGraphUnavailableNodeHasNoLink — `TestGraphUnavailableNodeHasNoLink` |
| M-r | graph row pitch < box height | TestGraphNodesDisjointAndInBounds — `TestGraphNodesDisjointAndInBounds` |
| M-s | object order from map iteration | TestGraphDeterministic / TestGraphGolden — `TestGraphDeterministic`, `TestGraphGolden` |
| M-t | golden file truncated to 0 bytes | TestGraphGolden — `TestGraphGolden` (anti-vacuity arm) |
| M-u | graph renders a grade label | TestGraphCarriesNoGrade — `TestGraphCarriesNoGrade` |
| M-ah | a new region rendered only on some variants (e.g. graph omitted on `?object=`) | TestLiveScriptContract — `TestLiveScriptContract` (per-variant table row) |

Exact mutation commands:

```sh
# M-l
go test -count=1 -race -v ./host/workbench -run '^(TestRenderEmitsOnlyLocalLinks)$'
# M-q
go test -count=1 -race -v ./host/workbench -run '^(TestGraphUnavailableNodeHasNoLink)$'
# M-r
go test -count=1 -race -v ./host/workbench -run '^(TestGraphNodesDisjointAndInBounds)$'
# M-s
go test -count=1 -race -v ./host/workbench -run '^(TestGraphDeterministic|TestGraphGolden)$'
# M-t
go test -count=1 -race -v ./host/workbench -run '^(TestGraphGolden)$'
# M-u
go test -count=1 -race -v ./host/workbench -run '^(TestGraphCarriesNoGrade)$'
# M-ah
go test -count=1 -race -v ./host/daemon -run '^(TestLiveScriptContract)$'
```

**Controller commit message:** `feat(workbench): render deterministic checked-edge SVG graphs`

## C5 / M5 — Bounded approvals chain and read-only decisions pane

Production estimate/budget: ~180 added LOC, including JS/CSS/templates; ≤~250. One controller commit.

Exact files:

- `host/broker/approvals_view.go`
- `host/broker/approvals_view_test.go`
- `host/workbench/render.go`
- `host/workbench/decisions_test.go`
- `host/workbench/render_test.go`
- `host/daemon/workbench.go`
- `host/daemon/workbench_decisions_test.go`

Tests FIRST: all new feature assertions below must be red on the pre-milestone snapshot; existing controls green as PD5 specifies.

- **AC5.1 — TestFoldApprovals:** `TestFoldApprovals`: request A, request B, decide A deny, request C, decide C approve gives the newest-first order [C approve, A deny, B pending]. A request decided twice shows the newest decision. Order by newest relevant head encounter, including newest duplicate decision; not request timestamps or map iteration.
- **AC5.2 — TestRecentApprovalsBounded:** `TestRecentApprovalsBounded`: a 200-head chain does ≤ 40 head-object GetObject reads and ≤ 2×kept-summaries additional body reads (≤ 80 total on the counting fake) and sets Truncated. (Checked against D7: the walk reads ≤ 40 head objects; each of the ≤ 20 kept summaries needs its request body for effect/scope/requester/cost and, when decided, its decision body for `DecidedBy`, so ≤ 40 + 2×20 = 80. The registry head is one `GetRegistryHead`, not a `GetObject`.) Assert <=20 summaries and one registry-head read separately. Fold bounded heads before fetching request/decision bodies; never fetch bodies for discarded summaries. Track head and body reads separately; test duplicates, truncate at 40, and no-head empty state. Test broker I/O propagation and daemon sanitized 5xx in new TestRecentApprovalsStoreError / TestWorkbenchDecisionsStoreError, without altering the producer.
- **AC5.3 — TestRecentApprovalsMalformed (both broker and daemon; PD2):** `TestRecentApprovalsMalformed`: a wrong-semantic-id head gives `ErrApprovalChainMalformed`. The workbench renders the UNAVAILABLE line with status 200. PD2 supplies the same name in both packages; include undecodable payload and missing referenced object, never confuse malformed chain with transport I/O. Read-only walk bounded even on cycles.
- **AC5.4 — TestWorkbenchDecisionsPane, TestWorkbenchNewViewTypesInCensus, TestWorkbenchViewFieldsAllRender:** `TestWorkbenchDecisionsPane`: a real `HumanHandler` request plus `DecideApproval` on a daemon store renders a checked request link and "deny by <decidedBy>". With no head it renders "No approval requests recorded". **Census (r2):** `DecisionsView` (and `DecisionRow`) is appended to `workbenchViewTypes`, and `TestWorkbenchViewFieldsAllRender` continues to pass (V13). Use real broker.NewHumanHandler(d.store) and DecideApproval; assert effect/scope/requester/cost, deny by decidedBy, checked ?object= link, as of entry N and Truncated text; empty positive state. Append DecisionsView and DecisionRow to census (PD1). No broker types/imports in workbench production.
- **AC5.5 — TestWorkbenchDecisionsReadOnly:** `TestWorkbenchDecisionsReadOnly`: the pane contains no `<form`, no `<button` and no `method=`. Restrict the test extraction to the decisions section, require it exists before checking absence of form/button/method=; whole-page forms elsewhere are not this criterion.

Implementation only after red proof: Add ApprovalReader, bounded RecentApprovals and pure FoldApprovals in approvals_view.go, reusing producer wire types without editing approve.go. Bound 40 heads/20 summaries/80 object reads. Map into DecisionsView/DecisionRow; malformed chain becomes 200 UNAVAILABLE, I/O follows existing sanitized error path; checked read-only links and as-of entry text.

Acceptance commands (environment above; named PASS required):

```sh
go test -count=1 -v ./host/broker -run '^(TestFoldApprovals|TestRecentApprovals.*)$'
go test -count=1 -v ./host/daemon -run '^(TestRecentApprovalsMalformed|TestWorkbenchDecisions.*|TestLiveScriptContract)$'
go test -count=1 -v ./host/workbench -run '^(TestWorkbenchDecisionsReadOnly|TestWorkbenchNewViewTypesInCensus|TestWorkbenchViewFieldsAllRender)$'
go test -count=1 ./host/broker ./host/workbench ./host/daemon
go vet ./...
```

Mutation drills (apply one form at a time and require the named assertion red):

| Mutant | Isolated edit | Exact test / assertion |
|---|---|---|
| M-z | fold: oldest decision wins | TestFoldApprovals — `TestFoldApprovals` |
| M-aa | approvals walk unbounded (> 40 head reads) or body reads not bounded by kept summaries (> 80 total) | TestRecentApprovalsBounded — `TestRecentApprovalsBounded` |
| M-ab | malformed chain → 500 page | TestRecentApprovalsMalformed — `TestRecentApprovalsMalformed` PD2: exact daemon test requires status 200; the broker test separately requires ErrApprovalChainMalformed. |
| M-ac | decisions pane gains a form/button | TestWorkbenchDecisionsReadOnly — `TestWorkbenchDecisionsReadOnly` |

Exact mutation commands:

```sh
# M-z
go test -count=1 -race -v ./host/broker -run '^(TestFoldApprovals)$'
# M-aa
go test -count=1 -race -v ./host/broker -run '^(TestRecentApprovalsBounded)$'
# M-ab
go test -count=1 -race -v ./host/daemon -run '^(TestRecentApprovalsMalformed)$'
go test -count=1 -race -v ./host/broker -run '^TestRecentApprovalsMalformed$'
# M-ac
go test -count=1 -race -v ./host/daemon -run '^(TestWorkbenchDecisionsReadOnly)$'
```

**Controller commit message:** `feat(workbench): show bounded read-only approval decisions`

## C6 / M6 — S7 usage rehearsal and local Chrome drill

Production estimate/budget: ~10 added LOC, including JS/CSS/templates; ≤~250. One controller commit.

Exact files:

- `docs/QUICKSTART.md`
- `host/daemon/daemon.go`
- `host/daemon/workbench_live_chrome_test.go`

Tests FIRST: all new feature assertions below must be red on the pre-milestone snapshot; existing controls green as PD5 specifies.

- **AC6.1 — TestWorkbenchLiveQuickstartControl plus verbatim QUICKSTART rehearsal:** QUICKSTART gains "Open the live workbench": open `http://127.0.0.1:7644/workbench`, what the live list, footer, graph and decisions pane mean, the JS-off behaviour, the 16-slot posture (D8), and a commit appearing without reload. It is executed verbatim, with the transcript in the PR. It uses the next free section number (V23 collision noted, not fixed here). The `Handler` route doc comment names `GET /workbench/live.js`. PD5 binary-free docs control fails before docs edit. Use next free section number 11 (inspect headings again before writing), do not fix existing section-10 collisions. Include complete existing genesis/commit payload construction and commands by a checked reference or verbatim recipe, and a second commit whose observed browser appearance proves no reload. Rehearse against an isolated disposable daemon/store rooted under a temporary directory, never the user’s live store. Controller retains -v and actual manual transcript in PR/snapshot evidence, no extra tracked artifact.
- **AC6.2 — TestWorkbenchLiveChromeDrill (WORLD_CHROME required for local acceptance):** `TestWorkbenchLiveChromeDrill` (runs only with `WORLD_CHROME=<path>`; the executor runs it and pastes the `-v` transcript; it is **not** a CI gate, V27). Setup: `httptest.NewServer` over `d.Handler()`, wrapped by a recorder that tees `POST /agui/` bodies and counts `GET /workbench`; `d.aguiBudget = 1s`; Chrome `--headless=new --user-data-dir=t.TempDir()`. Oracle: a POST with `lastIndex == k`, then (after the test's `d.store.Commit` of entry k+1) a `GET /workbench` re-fetch, **and** a later POST with `lastIndex == k+1`. The JS-off arm (`--blink-settings=scriptEnabled=false`) sees the page GET and zero POSTs in a 3 s window. This arm is an instrument control, not load-bearing. PD4 reduces only JS-off control window to 1 s after page GET. No DOM-dump or virtual-time oracle; server-side trace is mandatory. Use real httptest server, 1 s aguiBudget, per-instance 10 ms tick, POST-observed commit handshake, each guard=20 s, Chrome subprocess guard=30 s. Tee bodies without consuming handler input. For long-log arms follow terminal cursor across re-POSTs. Chrome profile, HOME/cache/TMPDIR and logs all under t.TempDir, process canceled and reaped in cleanup. Missing WORLD_CHROME skips in ordinary CI; explicit invalid path fails; local C6 acceptance requires PASS. M6 does not change ci.yml (PD6). Red-first documentation precondition (PD5): this new drill first requires the new Open the live workbench heading/route documentation via the shared docs control before launching Chrome. C5 code already has the live feature, so the C6 pre-milestone red is the absent usage documentation; after docs land the full browser trace must PASS. Do not attribute a skip to red-first proof.

Implementation only after red proof: Add QUICKSTART section 11 Open the live workbench and Handler route doc comment. Execute printed happy path including payload construction and a new visible commit without reload; run env-gated local Chrome request-trace drill. No CI Chrome dependency or workflow edit.

Acceptance commands (environment above; named PASS required):

```sh
go test -count=1 -v ./host/daemon -run '^TestWorkbenchLiveQuickstartControl$'
curl --fail --silent --show-error http://127.0.0.1:7644/workbench >/dev/null
set -euo pipefail
: "${WORLD_CHROME:?Set WORLD_CHROME to the local Chrome executable}"
log="$(mktemp)"
trap 'rm -f "$log"' EXIT
go test -count=1 -timeout 90s -v ./host/daemon -run '^TestWorkbenchLiveChromeDrill$' >"$log" 2>&1
cat "$log"
grep -q -- '^--- PASS: TestWorkbenchLiveChromeDrill ' "$log"
go test -count=1 ./host/daemon
set -euo pipefail
: "${WORLD_EXEC_NODE:?WORLD_EXEC_NODE is required}"
[ -x "$WORLD_EXEC_NODE" ]
log="$(mktemp)"
trap 'rm -f "$log"' EXIT
go test -race -count=1 -timeout 60s -v ./host/workbench ./host/daemon -run '^(TestLiveScriptPure|TestLiveScriptAssetNonEmpty|TestLiveScriptContract)$' >"$log" 2>&1
cat "$log"
for t in TestLiveScriptPure TestLiveScriptAssetNonEmpty TestLiveScriptContract; do
  grep -q -- "^--- PASS: $t " "$log" || { echo "$t did not PASS (SKIP or zero tests)"; exit 1; }
done
go vet ./...
```

Mutation drills (apply one form at a time and require the named assertion red):

| Mutant | Isolated edit | Exact test / assertion |
|---|---|---|
| None assigned by design | Local S7 rehearsal and documentation control | AC6.1/AC6.2 still mandatory |

Exact mutation commands:

```sh
# No design mutants assigned to M6.
```

**Controller commit message:** `docs(workbench): rehearse live usage with a local Chrome drill`

## Gates and handoff

Per-milestone gate list is exactly the acceptance commands above: narrowest affected packages plus `go vet ./...`. The design’s broad build/full-test/race gate is required at final handoff; the brief’s narrow package gates avoid repeating full loaded-runner suites six times. At every boundary require all new AC tests and unchanged package controls. Final commands:

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export PATH=/opt/homebrew/bin:$PATH
./scripts/verify_ail.sh
./scripts/verify_go.sh
go test -race -count=1 ./host/workbench ./host/daemon ./host/store ./host/broker
```

Apply each mutant alone on its Cn snapshot, require the named behavioral assertion red on compiling code, restore exact bytes without git and rerun green. Run every alternative form (M-e nosniff/MIME; M-aa heads/body bounds; M-ag throw/silent skip; M-p omission/unrendered across every type). M-b additionally checks every existing assertWorkbenchSecurityHeaders caller. M-s requires independent node-order expectation, not accidental map randomness. Record assertion/rc/restore; SKIP, zero tests, build failure, harness timeout, unrelated failure and sandbox socket refusal are not kills. SURVIVED blocks acceptance. Temporary protected-surface mutations M-d/M-g/M-l are restored, never delivered.

Mutation M-p is listed once in C1 but MUST repeat at C4 and C5 for each new type and row type; no coverage deferral. M-ah runs at C4 when graph variant rendering exists, and repeat its decisions-region conditional form at C5. M-b runs all unchanged literal-helper callers as well as the listed commands. Root JSON binds the exact named test commands.

S7 execution uses a controller-provided disposable daemon at 127.0.0.1:7644 wholly in a rehearsal temporary directory; execute QUICKSTART section 11 verbatim, including the existing commit payload construction/second commit, and retain actual request trace/manual transcript with PR evidence. Curl health alone does not prove S7. Chrome drill makes the no-reload/terminal-cursor behavior observable; docs grep/control is only instrument health.

Risks and controller work:

- Node harness cannot silently skip on any CI Go leg. WORLD_EXEC_NODE is already exported by row-140 setup; invalid configured path fails, absent node skips only locally. Binary-free companion asset/contract tests always run; explicit PASS loop catches SKIP and zero-match. AILANG-only CI leg does not execute Go/JS, and cannot substitute for the required Go job.
- Chrome is local-only and env-gated; local acceptance requires configured executable and named PASS on a socket-capable host. Generic CI has no Chrome prerequisite. Avoid dump-dom/virtual-time oracle; tee requests and reap browser. No socket failure is pass/fail evidence under sandbox.
- Known base host/broker head/tail timeout (row 166) can appear under full-parallel verify_go.sh. Preserve raw full-gate red; compare complete controller baseline and focused unchanged-package runs. Do not turn a retry into a claim that the first gate passed or loosen unrelated broker behavior.
- Self CSP admits same-origin non-nosniff text/plain bodies; GET /v1/head currently contains daemon-minted hashes. Keep exact row-149 literal and escaping checks; future echo routes and hash-source hardening are controller findings, not this implementation.
- Approvals are registry/object writes, not log entries: the decisions pane refreshes on log-triggered page swaps only. No invented CUSTOM event or approval action.
- Each planned milestone is <=250 added production LOC, counting JS/template CSS/HTML and excluding tests/docs/goldens. M3 budget 230 leaves little slack: keep the five pure seams and share state helpers. Executor must measure each boundary; if an actual implementation exceeds ~250, split its work into named reviewable substeps M3a (parser/state) and M3b (lifecycle/route), retaining one cumulative C3 commit and all AC/mutation evidence; disclose the split and actual LOC, never hide it in another milestone.

Controller owns commits/snapshots, completed baseline/flake adjudication, socket-capable final gates, CI, S7/Chrome transcript and design §9 follow-ups. Findings are recorded for controller action without modifying mission docs or frozen handlers. Compare frozen surface against base and merge head.

**Coverage census:** 6 milestones/commits; 30 numbered ACs (AC1.1–1.6, AC2.1–2.5, AC3.1–3.6, AC4.1–4.6, AC5.1–5.5, AC6.1–6.2), plus cross-cutting AC0; 35 mutants M-a…M-ai; 8 numbered planner decisions. Nothing deferred.

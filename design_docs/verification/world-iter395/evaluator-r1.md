# Evaluator report r1 - iteration 395, row 243 V1 (`w-world-address-bar-verbs`), head 4962fbea, base db2b7423

## Verdict: PASS - 91/100, no BLOCKING findings

| Category | Score |
|---|---|
| Correctness vs design and plan section 9 | 28/30 |
| Tests: non-vacuous, kill mutants, cover operands | 25/30 |
| Security / fail-closed | 19/20 |
| Coding standards (pure core, effects at the edge) | 9/10 |
| Scope and size | 10/10 |
| **Total** | **91/100** |

## Size and scope
- `git diff --numstat db2b7423 HEAD -- cmd host docs`: 1759 insertions + 150 deletions = **1909** (cap 2000). Matches plan 1.5.
- `git diff --name-status`: 34 paths, all plan section 5 plus the two plan files. Nothing under tools/launchd, world/, packages/, scripts/, cmd/world-publish, any `.ail`, or skills.
- `host/verbs` imports only errors, flag, fmt, io, regexp, strings, time (flag.Duration type only), unicode/utf8; no I/O, no clock, no randomness.
- No new CORS header (`Access-Control` appears only in a test assertion), no `t.Parallel`, no new production `exec.Command` (the only `os/exec` is in new test files that launch node), no new POST arm. `surfaceCSP` unchanged. Only one mux line added in `daemon.go`, so `/v1` is still the 13 routes (TestFrozenV1RouteTableMatchesMux passes).
- Row 201: `approve.js`, `passkey.js`, static `surface.js` and the decide/approve handlers have zero diff. The `approve` CLI (`cmd/ailang-worldd/approve.go`) changes by one import line and one call-site repoint (`parseInterleaved` -> `verbs.ParseInterleaved`). Same for approver.go, episode*.go, preset.go, session*.go.
- Pre-existing production files outside design section 8: `episode_markdown.go` (+4 -4, plan P1: `daemon.SortedCounts` / `WriteReportRefusals` names). Behaviour identical; the three finish-markdown tests pass unedited.
- ParseInterleaved call sites: base `parseInterleaved(` had **10** call sites (approve, approver, episode, episode_finish, episode_launch, episode_report, preset, session, session_new x2); HEAD has **10** `verbs.ParseInterleaved` and 0 old. All still green (cmd whole package passes).

## Pre-existing `_test.go` edit audit
Exactly one file existing at base is edited: `host/surfacepage/app_test.go` (+2 -1). It only appends the five new export names (`VERB_ROUTE_BUDGET_MS`, `askVerb`, `verbAct`, `verbOutcome`, `verbPanel`) to the `TestSurfaceAppPure` export-list assertion (`.sort()`ed). Nothing weakened; the list still pins the full export set. No other pre-existing test changed.

## Test runs (mine, after ALL_DONE existed; AILANG_BIN pinned; logs in judge/)
| Command | Result |
|---|---|
| `go test -count=1 ./host/verbs/ ./host/surfacepage/` (WORLD_EXEC_NODE) | ok / ok |
| `go test -race -count=1 ./host/verbs/ ./host/surfacepage/` | ok / ok |
| `go test -count=1 ./host/daemon/` | rc=0, 255s |
| `go test -count=1 ./cmd/ailang-worldd/` | rc=0, 124s |
| `go vet ./host/... ./cmd/...` | rc=0 |
| `go test -run '^$' ./...` | rc=0 |
| live Chrome leg (TestVerbBarEpisodeReportInChrome, TestNeedsYou x2, WORLD_CHROME + WORLD_REQUIRE_PASSKEY_BROWSER=1) | rc=0, all 3 PASS |
Not run (per directive): verify_go.sh, `go test ./...`. Controller's gates summary read: `GATES_DONE ... FAILS: none`.

## Mutation testing
Ledger: `$HOME/.ailang/state/world-iter395/judge/mutants.tsv`, **78 lines** (every run and re-run appended; per-run logs in `judge/logs/`). Worktree `git status --short` clean at the end. Each mutant: edit asserted to change the file, `go vet` of the package passed before reading the result (no COMPILEERR/INVALID occurred).

`comm -23 my-ids plan-ids` (excluding my PLAN-* re-runs and -cmd/-WHOLE/-surfacepage suffix re-runs): 60 ids NOT in the plan's 120-row table, e.g. X-3B-CLASS X-3B-REST X-EP-LEN X-EP-LEN2 X-REF-LEN X-APPROVE-HREF X-CTRL-1F X-SHELL-PIPE X-SHELL-BACKTICK X-QUOTE-JOIN X-MAXWORDS-CONST X-REASON-NEVER X-R-VERBFIELD X-R-CLI X-R-EXIT X-R-NAVCLASS X-R-QECHO X-R-STDERR-NL X-R-EPBOUND-GE X-E-START X-E-STARTERR X-E-FINISHERR X-E-EXECCACHE X-P-START X-P-UNCOMMITTED X-P-ATTEMPTS and 30 `X-JS-*` ids (full list in judge/my-ids.txt vs judge/plan-ids.txt).

Controls: must-die `CTL-MUSTDIE-LONG` (Split bound +1) KILLED by TestSplit; must-live `CTL-MUSTLIVE-COMMENT`, `CTL-JS-MUSTLIVE`, `CTL-CH-MUSTLIVE` (comment-only edits) SURVIVED as required.

### Plan table sample re-run (17 rows, all as the plan says; killer named)
S-SPACE (TestLookup,TestQuoteRoundTrip,TestSplit), S-DQ-ESC-BS (TestSplit), L-NEVER (TestLookup), L-4B-EPISODE (TestLookup), L-NAV-WIDE (TestLookup), F-FOLLOW (TestEpisodeReportFlags), R-Q (TestSurfaceVerbBoundsAndErrors), R-HEADERS (TestSurfaceVerbBoundsAndErrors,TestSurfaceVerbOutcomes), E-TIMEOUT (TestEpisodeReportRESTRoute), R-JSON-INDENT (TestVerbRouteReportIsTheCLIs), J-404-FALLBACK, J-SUPERSEDE-CHECK, J-BODY-RACE, J-BUDGET-USE, J-ANY-FALLBACK (all TestSurfaceAppVerbs, node), JW-GO-ROUTE (TestVerbBarEpisodeReportInChrome), JW-GO-EMPTY SURVIVED (named survivor).

### My own mutants
KILLED, each by a behavioural test (no inventory/golden-only/hash kill; the killers are TestSplit, TestLookup, TestQuoteRoundTrip, TestApproveTarget, TestEpisodeReportFlags, TestSurfaceVerbOutcomes, TestSurfaceVerbBoundsAndErrors, TestEpisodeReportShowsBindSeedFailure, TestEpisodeStartReadProjectionsAndSessionRoutesCannotWrite, TestEpisodeReportPrintsExecCache, TestSurfaceAppVerbs, TestVerbBarEpisodeReportInChrome):
X-CTRL-1F, X-SHELL-PIPE, X-SHELL-BACKTICK, X-QUOTE-JOIN, X-3B-CLASS, X-3B-REST, X-EP-LEN, X-EP-LEN2, X-REF-LEN, X-APPROVE-HREF, X-ADMIT-JSON-ALLOWED, X-BAR-REFUSAL-TEXT, X-R-VERBFIELD, X-R-CLI, X-R-EXIT, X-R-STDERR-NL, X-R-QECHO, X-R-NAVCLASS, X-E-EXECCACHE, X-E-START (killed only in the WHOLE host/daemon run, by TestEpisodeStartReadProjectionsAndSessionRoutesCannotWrite; the narrow run and cmd run both survived), X-P-ATTEMPTS, X-JS-ABORT-ONLY, X-JS-NOABORT, X-JS-NAV-OR, X-JS-NAV-DECIDE-ONLY, X-JS-TERM-AT, X-JS-HEAD-PROMPT, X-JS-STDERR-EXIT, X-JS-STDERR-SUFFIX, X-JS-STATUS, X-JS-OBJ, X-JS-NEEDS-TERM, X-JS-LATE-MSG, X-JS-REFUSED-SAY, X-JS-PANELCLEAR, X-JS-ASSIGNCLEAR, X-JS-COPYBTN, X-JS-ENC, X-JS-GO-ADDRESS-RAW (node), X-JS-CLEAR-GUARD, X-JS-CLEAR-NEVER, X-JS-PANEL-REPLACE (Chrome).

SURVIVED (each confirmed on the whole package of the file's tests; non-blocking, all are test-coverage gaps in correct code):
- X-MAXWORDS-CONST (`MaxWords` 32 -> 33): the TooMany tests build their input from the constant, so the design's 32 is pinned only by the error text. (MaxLine 1024 is pinned by the golden's 1025-byte case.) Whole host/verbs run.
- X-REASON-NEVER (drop `; at a terminal: <cli>` from the five never-verb reasons): survives host/verbs (whole), narrow host/daemon, narrow cmd and WHOLE cmd (whole host/daemon not run for this mutant). The page shows a `refused` answer as only `say(reason)`, so this is operator-visible text with no assertion (the golden has no never word; plan V1-D5 asserts it "already ends at a terminal").
- X-R-EPBOUND-GE (`>=` on the 256-byte episode bound): survives WHOLE host/daemon. Off-by-one at exactly 256 bytes; unreachable for realistic names.
- X-P-START, X-P-UNCOMMITTED (the `start:` block and the `uncommitted intents` line of the moved printer): survive WHOLE host/daemon (start) and WHOLE cmd (both). The base cmd tests never asserted `start:` either, so this is a pre-existing gap carried by the move, not a regression.
- X-E-STARTERR, X-E-FINISHERR (extraction ignoring a start/finish read error): survive narrow daemon and narrow cmd only; not run whole (time). Same family as the named survivor E-ASb (error arms of a store read that cannot be made to fail).
- X-JS-FINALLY-NULL (`if (verbJob === job) verbJob = null` -> unconditional): survives WHOLE host/surfacepage. Real gap: with it, an older superseded ask finishing would null the newer job, so a third submit would not abort the second.
- X-JS-NEEDS-ADDR (`result` no longer requires `cli`): survives WHOLE surfacepage; a `result` answer without `cli` would render head `$ undefined`.
- X-JS-ESC-VERB (Escape no longer closes the panel) and X-JS-ENTER-REFOCUS: survive WHOLE surfacepage and the Chrome set; Escape-closes-panel and refocus have no test.
- Equivalent / cosmetic: X-JS-STOPCATCH, X-JS-ROLE, X-JS-HEADERS.
- JW-GO-EMPTY (plan named survivor): SURVIVED, reproduced.

## Assessments
**No silent fallback (my own reading, backed by the node harness):** `verbOutcome` returns `{address:true}` only when status is 200, JSON parses, `schema==='world/surface-verb/v1'`, `v===1`, outcome is an own key of `VERB_NEEDS`, and the outcome's members are non-empty strings. Every other path returns `{say}` (non-2xx, network error, non-JSON, schema/v/outcome wrong, `constructor` prototype key, null/array, missing/empty members, two result members, stderr without exit 1), the timeout (`job.late`) says `did not answer in time`, and an abort without `superseded`/`late` says unreachable. `verbAct` calls `fx.say` only on that path; only the address, assign and panel arms call `fx.clear` (and only when the input still equals the submitted line). I found no path where failure navigates or changes the view. Mutants J-404-FALLBACK, J-ANY-FALLBACK, X-JS-STATUS, X-JS-OBJ, X-JS-NEEDS-TERM were all killed by the node harness. Minor: the page does not check that the answer's `q` equals what it sent; benign (the route is same-origin and host-guarded).

**V1-D1 (superseded ask): sound.** A newer `askVerb` marks the old job `superseded` and aborts it; the abort rejects both awaits through the shared `stop` promise, and the after-await `job.superseded ? {superseded:true}` check (P12) guards the case where the old body resolved in the same turn as the new submit. I proved both pieces matter: X-JS-NOABORT, X-JS-ABORT-ONLY and J-SUPERSEDE-CHECK are all killed by TestSurfaceAppVerbs case A/B. Gap (non-blocking): an empty submit goes to `out()` without calling `askVerb`, so it does not supersede a pending verb ask, and an older `address`/panel answer can still act afterwards; and X-JS-FINALLY-NULL shows the third-submit case is untested.
**V1-D2 (injected clock, budget not a tunable): sound.** `VERB_ROUTE_BUDGET_MS=15000` is a const, exported read-only; the harness checks timer ms, that `clearTimeout` runs on every path, and the stalled-header and stalled-body timeouts (J-BUDGET-USE, J-BODY-RACE, J-LATE killed).

**Named survivors.** `E-ASb` (start/finish read failure after an expired deadline would answer 503 instead of 500 under a broader `reportTimedOut`): honestly described. The extraction is correct: `reportBuildError` wraps only the builder error and `reportTimedOut` keys on it, matching the base handler line by line; unreachable because the builder fails first under an expired deadline. The operator sees no difference except a 500 vs 503 class on a pathological store failure. Acceptable. `JW-GO-EMPTY`: honestly described; reproduced. With the route healthy the operator sees the same view (empty `q` answers `address`, home); with the route down an empty Enter would say "unreachable" instead of going home. Acceptable, non-blocking.

**Route:** read-only (`surface_verb.go` has no store write, no credentials store, no spawn; it reads through `d.episodeReport`, the same builder as `/v1`), sits behind `approvalRoute` (Host guard, 421 test), GET only, `writeApprovalJSON` (no-store, nosniff), missing `q` = 400, 1025-byte line = 200 refusal via `Split`'s first statement, internal errors are logged and show only `internal store failure` or `read deadline (10s) exceeded` (R-ERR-LEAK killed). `/v1` mapping preserved: 503 for a builder timeout (E-TIMEOUT killed), 500 otherwise, 404 when not found. Low note: `q` is echoed unbounded in the answer for an over-long line (up to the server's header limit); no amplification concern at 1 MB.

**"Byte-identical to the CLI":** VR3 (`verb_report_e2e_test.go`) runs the real CLI `run([... "episode","report", ...])` and the real route against one `daemon.New` seeded with real journal intents, for two episodes with different shapes (one with a plan refusal, one with a handler refusal), in table and `--json` form plus the not-found case. R-JSON-INDENT killed by it. The printer is one shared function (`daemon.WriteEpisodeReportText`), so equality is structural as well as tested.

**Grammar:** Split/Quote refusals and round trip are well covered (X-CTRL-1F, X-SHELL-*, X-QUOTE-JOIN killed). `Lookup`: the 27-word shadow set is checked against the CLI switch by an AST test (TestVerbTableMatchesTheCLISwitch, passes); bare/unknown `episode` and three-word `world` are refused (L-4B-* killed); `ailang-worldd` strip works. One expectation worth knowing: the bar profile refuses `--markdown`, `--follow`, `--interval`, `--addr`, `-h`, so those fall to a refusal text that is covered (F-* killed).

**Non-vacuity per milestone:** M1c: F-PI and F-PI-CLI (plan) plus my X-3B/EP rows show TestParseInterleaved/TestEpisodeReportFlags are M1c's own, and the 10 call sites are protected by the existing cmd suite (plan F-PI-CLI killed by a pre-existing cmd test, TestApproverRevokeListsWhatItApproved). M3a: X-R-VERBFIELD/CLI/EXIT/QECHO and R-Q/R-HEADERS are killed by tests added in M3a/M3b (TestSurfaceVerbBoundsAndErrors, TestSurfaceVerbOutcomes), and the extraction (E-TIMEOUT, E-WRAP) is held by the pre-existing TestEpisodeReportRESTRoute. M4a: every JS mutant above is killed by TestSurfaceAppVerbs (node), plus the Chrome file for the wiring (JW-GO-ROUTE, X-JS-CLEAR-*, X-JS-PANEL-REPLACE). Examples requested: shadowed word returns address (L-NEVER/L-OTHER killed), step 4b falls through (L-4B-* killed), page navigates on non-2xx (J-404-FALLBACK killed), drop the abort (X-JS-NOABORT killed), drop the q bound (S-LONG/R-Q killed), route prints a different report (R-JSON-INDENT killed).

**Intermediate commits:** not individually re-run by me (the controller ran the sweep); the plan's per-milestone numstat sums (203/510/697/905/1248/1491/1762/1909) match the head total I measured.

## BLOCKING findings
None.

## Non-blocking findings (ranked)
1. Never-verb refusal text loses its CLI line unchecked (X-REASON-NEVER survives every suite): add a golden row for one never word.
2. X-JS-FINALLY-NULL and the empty-submit-does-not-supersede gap in V1-D1: add a three-submit harness case.
3. `MaxWords` constant is not pinned independently of the constant (X-MAXWORDS-CONST).
4. Result answers without `cli` are accepted by the page (X-JS-NEEDS-ADDR).
5. Escape-closes-panel and the Enter refocus have no test.
6. Moved printer's `start:` block and `uncommitted intents:` line have no test (pre-existing gap; X-P-START, X-P-UNCOMMITTED).
7. X-R-EPBOUND-GE: exact-256-byte boundary untested.
8. Real refusal text is a header `say` that fades in 2.6 s (known residual RS-SAY).

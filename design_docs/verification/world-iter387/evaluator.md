# Evaluator r1 — row 201 A3a sprint 4 (base bed9d135, head 6723807a)

VERDICT: PASS  — score 90/100. No BLOCKING findings.

| Category | Pts | Score |
|---|---|---|
| Correctness vs design / plan §9 | 30 | 27 |
| Tests (non-vacuous, kill mutants) | 30 | 27 |
| Security / fail-closed | 20 | 19 |
| Coding standards | 10 | 9 |
| Scope / size | 10 | 8 |

## Size and scope
numstat `cmd host`: +354 / -35 = 389 (cap 2000). Files: 3 cmd tests/driver, QUICKSTART, app.js, app_test.go, plus the two plan files (matches §5). Only production file is host/surfacepage/app/app.js; approve.js, surface.js, passkey.js, daemon, approval, CSP, routes untouched. No innerHTML/eval/document.cookie/t.Parallel added in code. Worktree status clean at end.
(-2 scope: the plan/json add ~1500 lines of docs, outside the stated unit but expected.)

## Review
- requestPermission is called only in askNotify, wired solely to the `#notify-me` click. ASK-LOAD mutant (call on load) killed.
- Notification only when permission==='granted', first read raises none, seen-set stops repeats, refused and recorded excluded. All pinned by pure test plus XA4.
- Body/title/tag are plain text passed to the Notification constructor. Decided list and needsYouHTML escape every field (esc).
- Click opens '/surface/approve/'+id, id is not escaped there, but request ids are constrained by store schema to `mr-[0-9a-f]*` (length 19) and location.assign is same-origin; not exploitable. Non-blocking: href could be built defensively (regex check or encodeURIComponent).
- N/Enter: handler is after the INPUT early-return so it does not fire in inputs; Enter on now-card links returns early so as not to double-activate. Show-all bounded by the daemon's 20 `decided` entries (slice of array).
- Deviations (i)-(v) sound and honestly bounded: first read silent, UTC hh:mm (documented in QUICKSTART), N built in app.js, no new timer (P3 documented limitation).

## Mutants (ledger: ~/.ailang/state/world-iter387/judge/mutants.tsv, 36 lines; outputs out_<id>.txt). Node cmd: `WORLD_EXEC_NODE=$(command -v node) go test -count=1 -vet=off ./host/surfacepage/ -run TestSurfaceAppPure`; Chrome cmd: `WORLD_CHROME=".../Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -v -count=1 -vet=off ./cmd/ailang-worldd/ -run <test>`.
Plan §4 rows sampled (all KILLED, killer in brackets): ESC-AGENT, ESC-ID-HREF, ESC-RECORDED, ESC-UNAVAIL, DEC-ESC-EP(isode), DEC-ESC-AB, DEC-FEW, DEC-BTN, NT-SEEN, NT-REFUSED, NT-RECORDED, NT-GRANTED, NT-FIRST, NT-HREF, NT-TAG, NT-SEEN-UPD, NT-GUARD, NH-DENIED [all TestSurfaceAppPure]; LR-PERM, LR-SEEN, LR-CLICK-OPEN, LR-CLICK-FOCUS, ASK-LOAD, ASK-WIRE [TestNeedsYouNotifiesOnlyWhenAskedOnce]; NKEY-LOWER, NKEY-UPPER, NKEY-OUT, ENTER-NOW, DEC-TOGGLE, DEC-NOWCARD [TestNeedsYouListsPendingRequests]; CARD-INNERHTML (approve.js textContent->innerHTML) [TestApprovePageShowsAHostileValueAsText].
Every kill is by a named behavioural test (pure node assertions or live-Chrome), none rests on an inventory/static test; the pure-test mutants fail only in TestSurfaceAppPure, not a static banned-token test.
Own mutants not in the plan table: operand-level variants NT-REFUSED/NT-RECORDED separate operands, NT-GUARD-ARR (`!Array.isArray` operand), LR-CLICK-FOCUS, NKEY-OUT, NH-DENIED, ESC-UNAVAIL: all killed (>=12 beyond the plan's exact edits).
Controls: must-die ESC-AGENT KILLED. Must-live: CTRL-LIVE-comment, CTRL-LIVE-DECIDED-CONST, CTRL-LIVE-CHROME (comment/no-op edits) SURVIVED as expected.
SURVIVOR (non-blocking): DEC-GUARD-ARR (`Array.isArray(r.decided)` -> truthiness) survives; the daemon always sends an array, so impact is nil but the guard's Array.isArray operand is unpinned (-2 tests, -1 other).
Not re-run by me: plan-table rows in M9c body/title/plural/time sub-edits beyond NT-HREF/NT-TAG, and NP-NOAPI (named equivalent).

## Test runs
- `go test ./host/surfacepage/` with WORLD_EXEC_NODE: ok rc=0; `-race` surfacepage: ok rc=0
- host/daemon named subset: ok rc=0; `go vet ./host/... ./cmd/...` rc=0; `go test -run '^$' ./...` rc=0
- Live Chrome leg (rc=0, 95s): PASS TestApproveCardNetworkDeadlines, TestApproveCardPasskeyDeadline, TestApprovePageInChrome, TestApprovePageShowsAHostileValueAsText, TestApprovePageRefusesASwappedRequest, TestApprovePageSSHSIGInChrome, TestNeedsYouListsPendingRequests, TestNeedsYouNotifiesOnlyWhenAskedOnce, TestPasskeyPageAbortsAPendingCeremony, TestPasskeyPageInChrome.
- Full ./cmd/ailang-worldd non-Chrome suite and verify_go.sh not run (controller's sweep).

## Non-blocking findings
1. app.js loadRequests: notification href not validated/encoded (safe by id schema; defence in depth).
2. DEC-GUARD-ARR survivor above.
3. Notification-body `expires` time is UTC, same as Decided; stated in QUICKSTART.

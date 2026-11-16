# Evaluator report, iteration 392, row 201 A4 sprint 4 ("the surface"), head cbb3783a

## Verdict: PASS, 88/100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design and plan section 9 | 28 / 30 |
| Tests: non-vacuous, kill the mutants | 23 / 30 |
| Security / fail-closed | 19 / 20 |
| Coding standards | 9 / 10 |
| Scope and size | 9 / 10 |

(Scope: 1 point off for 1,947 vs the 2,000 cap leaving little room; standards: 1 point off for the stderr-only failure of `approve --list`'s second half. Neither is a defect.)

## Size and scope (measured)
- `git diff --numstat 799660ac HEAD -- cmd host docs`: 1,935 insertions + 12 deletions = **1,947** (`-- cmd host` alone: 1,939). Cap 2,000. Matches the plan.
- `git diff --name-status`: 34 paths, all in plan section 5 plus the two plan files. Nothing under `tools/`, `scripts/`, `world/`, `packages/`, `cmd/world-publish`. No edit to `passkey.js`, `approve.js`, `host/surfacepage/{approve,passkey}/` (checked by name filter: empty). `host/approval` gains only `decision_herr_test.go` (+37, test only). No `t.Parallel` in any touched test. `daemon.go` +2 lines: only the two new `GET /surface/decide...` routes, both through `d.approvalRoute` (same Host/origin guard as the other approval routes); no `/v1` route, `surfaceCSP` and CORS untouched (`setSurfaceHeaders` reused; the page tests assert the CSP and no ACAO).

## Test runs (AILANG_BIN pinned, WORLD_EXEC_NODE=/opt/homebrew/bin/node), after ALL_DONE existed
| Command | rc |
|---|---|
| `go test -count=1 ./host/approval/... ./host/surfacepage/` | 0 (3 ok) |
| `go test -race -count=1 ./host/surfacepage/` | 0 |
| `go vet ./host/... ./cmd/...` | 0 |
| `go test -run '^$' ./...` (compile fence) | 0 |
| `go test -count=1 ./host/daemon/` (266 s) | 0 |
| `go test -count=1 ./cmd/ailang-worldd/` (128 s) | 0 |
| live Chrome leg (6 named tests, `WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1`) | 0, 6 PASS, 0 SKIP (TestApproveDecideWebAuthnLinkInChrome, TestApprovePageInChrome, TestDecidePageApproveWithPasskey, TestDecidePageFailsClosed, TestDecidePageSSHSIGPaste, TestPasskeyPageInChrome) |

Logs: `$HOME/.ailang/state/world-iter392/judge/run-*.log`. The controller's own sweep (`gates/legs.txt`) was all rc=0, 69 named tests PASS. My worktree is clean (`git status --short` empty).

## What the card trusts (own analysis, decide.js read in full)
- The page never accepts a hash from the server: P and Q are SHA-256 of the decoded `packetBytes`/`requestBytes`; it requires `Q == ref` (location), `sha256(packetBytes) == s.packet`, and `packet.requestRef == ref`. The 91-byte message is built locally from P, Q, the verb of the clicked button and the server's `window.max` (shown in the card text as "expires ..."). The verb is the button's, so a hostile or tampered GET cannot make the page sign a verb or request other than the displayed one. A tampered GET (hash mismatch, requestRef mismatch, unlisted/missing/wrongly typed key, non-safe integer, bad base64url, malformed window) goes to `refuse` before any control is unhidden and before any `credentials.get`.
- The expiry is the one server-controlled input: a hostile daemon could choose an expiry in the window, but the number is signed into the message, shown, and the POST re-runs the law (window bound), so the worst case is a shorter/longer validity inside the law's own bound, not another subject.
- The POST arms are unchanged (no diff in `DecideBySignature`/`DecideByAssertion`); `decidedBy` still derived by the daemon; the daemon still opens credentials read-only (`OpenCredentialsReadOnly` in the new `decisionPasskeys`); the credential-writer inventory test passes.
- Escaping: card text is `textContent` only; the single `href` is built from `encodeURIComponent`. *Needs you* values are all through `esc()` (quote, `<`, `>`, `&`); every one of scope, requester, refusal, reserved owner and href is mutation-killed (below). The link target is `/surface/decide/sha256:<esc(ref)>` and the route accepts only 64 lowercase hex, so a hostile ref is a dead link, not an injection. Notification `href` is `'/surface/decide/' + x.request` (unescaped, used as a navigation target not HTML; the ref comes from the daemon's own store). Fine.
- One fail-closed UX wrinkle (non-blocking): `shown()` calls `toISOString()`, which throws RangeError for a safe integer beyond +-8.64e12 s (`createdAt`, `deadlineAt`, `now`). `load()` then rejects, the page stays in `loading` with every control hidden/disabled and no message. Safe (no signing), unhelpful.

## The listing
`PendingDecisions` = one `RecentApprovals(40)` (at most 20 summaries), one packet scan of 101 (the 101st is only a sentinel), `ValidatePacket` once for each request carried by exactly one packet, `SubjectRefusal` for display, then the pure `FoldDecisions`: pending `Human.Approve` requests with at least one scanned packet; `packet`/`deadline` only for exactly one; decided and packet-less requests absent; `decisionsTruncated` for either bound; `decisions` never null. A store error is logged and yields `[]` (mint lists intact), so an unreadable world store fails safe (nothing offered). `approve --list`: `approveListDecisions` opens `--db` with `store.OpenReadOnly` and calls the same `daemon.PendingDecisions`; there is no `http`/dial anywhere in `approve_list_decisions.go` (read in full; imports: context, fmt, io, daemon, store, surfacepage). An unreadable `--db` is one stderr line, exit 0 unchanged. Note (non-blocking): exit 0 with only a stderr line means a script cannot tell "no decisions" from "could not read"; this is a documented decision (S4-D3) and the mint list is still printed. `GET /surface/decisions/{ref}` now needs a readable credentials store even for decided or reserved subjects (500 otherwise), the design's fail-closed choice.

## Plan section 1.4 and the named residuals
- Premises P1-P9: all verified true against the diff (`--list` never dialed at base; `passkey.js` hash pin untouched; `TestDecisionSubjectJSON`, `TestDecisionRoutesRegistered`, `TestSurfaceAppRoutes` edited as stated; JG-W15-OPERAND is a narrow-run survivor, killed by `TestSignCountLaw` in the whole `host/approval` run, which passes).
- **RS-VERB** (CLI `--webauthn-link` opens the frozen A2 page reading "Approve request decide-deny-..." for a deny): honestly described. The operator sees the A2 page's approve-wording plus the request text `decide-deny-<16 hex>` and the signed hash; the CLI prints DENY in capitals; the signature binds the verb byte, so a deny cannot be replayed as an approve (DM5/W6 rows, XD4 mismatched-verb paste 422). A passkey touch for a deny that reads "Approve" is a real UX hazard, but it is confined to the CLI-link path, the phone's path is now the card (which shows its own verb, asserted by XD1), and fixing it needs edits to frozen A2/A3 files. Acceptable non-blocking; Mark's call as the plan says.
- **RS-RIG** (Chrome drives the card against a test rig's POST, not the daemon's): honestly described (S4-D11). The rig uses the real `ParseDecideBody`, `DecisionRefusal` and judges and A3a guards; the daemon's POST arm is unchanged and covered by `TestDecidePOST*` and sprint 3's live-Chrome CLI test. The gap (the card's JSON against the daemon's real POST end to end) is real but small because the card posts `{v,verb,packet,request,expiry,webauthn|sshsig}` which the unchanged daemon grammar parses in both. Acceptable; first real run is the attended step 4.
- Named survivors the plan carries (`L-VALIDATE-ERR`, `CL-LIST-ERR`, `P-PASSKEYS-OLDER`, `J-TIMEOUT`, `J-CREDS`...): reasons are plausible; my `JD-GO-LIST-LISTERR-SILENT` re-confirms `CL-LIST-ERR`.

## Pre-existing `_test.go` edits (files existing at the base)
Only three, each minimal and no assertion weakened beyond the change:
1. `host/daemon/surface_app_test.go` +4 -2: adds `/surface/decide`, `/surface/decide/`, `/surface/decide.js/x` to the 404 list and the card page and script to the 421 list (more assertions, none removed).
2. `host/daemon/surface_decide_routes_test.go` +2 -2: the route list gains the two card routes; `GET /surface/decide/{ref}` leaves the 405 method list (it is now the card). The 405 loop that also asserts no CORS headers no longer covers GET for that path; `TestDecidePageAssets` asserts the card's headers and no ACAO instead. Required by the change.
3. `host/daemon/surface_decide_subject_test.go` +1 -1: adds `"passkeys"` to the exact key list.
No other pre-existing test file is touched.

## Mutation testing
Ledger: `$HOME/.ailang/state/world-iter392/judge/mutants.tsv`, 73 lines (every run and re-run). Harness `mut.py`: asserts the anchor occurs once, `git diff --stat` non-empty, `go vet <pkg>` clean, runs the command, restores with `git checkout`. Env copied verbatim (WORLD_CHROME etc. on the decide.js rows).

### (a) 17 sampled plan rows (spread M7a-M7f), all behaviourally KILLED, none by an inventory or hash-pin test
JG-BODY-PACKET (TestApproveDecideArtifactPostsTheSubjectsDigests), JG-SIG-CAP (TestApproveDecideArtifactReadBounds), JG-HERR-TOP (TestDecideByAssertionRefusesAMissingMessageAtW6), L-TRUNC-SCAN (TestDecideModelTruncation), L-SCAN-BOUND (TestSurfaceRequestsDecisionsTruncated), CL-TRUNC (TestApproveListDecisionLines), CL-OPEN-ERR (TestApproveListShowsDecisions), N-NOTE-SEEN and N-REF-ESC (TestNeedsYouListsDecisions), P-PAGE-BYTES (TestDecidePageAssets: HTTP status/bytes), P-PASSKEYS-GROUP (TestDecisionSubjectPasskeys), J-EXPIRY, J-FRESH, J-VERB-TEXT, J-BTN-DENY (TestDecidePageApproveWithPasskey, live Chrome), J-HASH-Q, J-REQREF (TestDecidePageFailsClosed, live Chrome). I confirmed the output of two kills are real assertion failures (J-EXPIRY: card posts a 422 "response was made for ..." state; P-PAGE-BYTES: "200, or not the bytes of index.html").

### (b) My own mutants (ids `JD-*`, none in the plan's table)
`comm -23 my_ids plan_ids` = 45 ids, `comm -12` = 0 (own ids file `judge/my_ids.txt`; 44 valid, 1 INVALID anchor `JD-APP-NOTE-SEENSET`, discarded).

Controls: **must-die** `JD-CTRL-DIE` (`return nil, ...` in FoldDecisions) KILLED by TestDecideModelListsPendingDecisionsOnly and TestDecideModelListsATwoPacketRequestRefused. **Must-live** `JD-CTRL-LIVE-JS` (comment added to decide.js) and `JD-GO-DEC-LIVE-100` (`decisionListScan-1` -> `100`, equivalent) SURVIVED as they must.

KILLED (killing test, all behavioural):
- JD-JS-NOTOBJ-ARRAY: TestDecidePageFailsClosed. JD-JS-DECIDED-KEEP and JD-JS-DENYTEXT: TestDecidePageApproveWithPasskey + TestDecidePageSSHSIGPaste. JD-JS-ORIGIN, JD-JS-ORIGINNOTE, JD-JS-POSTSIG-VERB: TestDecidePageSSHSIGPaste. JD-JS-UVREQ, JD-JS-CHALLENGE-RAW: TestDecidePageApproveWithPasskey.
- JD-APP-ESC-SCOPE, -ESC-REQUESTER, -ESC-REFUSAL, -ESC-RESERVED, -ESC-HREF, -COUNT, -MORE, -NOTE-RESERVED, -NOTE-REFUSAL, -LINKPREFIX, -NOTE-HREF, -RESERVED-FIRST: TestNeedsYouListsDecisions.
- JD-GO-PRINT-TRUNC, JD-GO-PRINT-DENYHINT: TestApproveListDecisionLines. JD-GO-LIST-OPENERR-SILENT: TestApproveListShowsDecisions. JD-GO-DEC-ANDTRUNC: TestDecideModelTruncation. JD-GO-DEC-NOW, JD-GO-DEC-CARRYING: TestSurfaceRequestsListsDecisions. JD-GO-SCRIPT-HEADERS, JD-GO-PAGE-HEADERS (drop `setSurfaceHeaders`): TestDecidePageAssets (HTTP header assertions, behavioural).
- None of the kills rests only on an inventory/hash-pin/static test; no re-run with an exclusion was needed.

SURVIVED (each after a WHOLE-package run with the plan's env; two raw "KILLED" ledger lines are reclassified, see notes):
| Id | Edit | Whole-package result | Assessment |
|---|---|---|---|
| JD-JS-FRESH-READY | after the pre-sign re-read, `return true` instead of "state is ready" | SURVIVED (cmd whole + surfacepage whole; only the baseline-red TestMarketingChromeSixShots failed, red unmutated for want of WORLD_TOOL_AILANG_BIN) | Gap: a re-read that finds the request decided/refused/expired still runs the ceremony. Law re-run at POST refuses it, so non-blocking. |
| JD-JS-WHATHASH | card's "deny" message-hash text shows the approve hash | SURVIVED (same) | Gap: the displayed deny hash is never asserted (the approve hash is, per A2/A3 tests). Display only. |
| JD-JS-NOTOBJ-NULL | drop `obj === null` operand | SURVIVED (same) | A `null` JSON body makes `Object.keys` throw; page stays "loading" with no controls (fail-closed, untested message). |
| JD-JS-NOW, JD-JS-WINOBJ, JD-JS-STATUS, JD-JS-B64TYPE, JD-JS-CREDOMIT | drop one operand of the `s.now`/`s.window`, `res.status !== 200` and `typeof s !== 'string'` checks; `credentials: 'include'` | SURVIVED under the narrow 3-Chrome-test run only; whole-package NOT run (time budget) | Provisional: analysis says each fails safe (hash check absorbs B64TYPE; window/now errors throw before any control; STATUS only changes the message). Whole-package confirmation not done. |
| JD-APP-DS-ARRAY | `Array.isArray(r.decisions)` -> truthiness | SURVIVED (host/surfacepage whole) | A non-array `decisions` value is not tested. |
| JD-GO-LIST-EARLYRET | restore `return exitOK` after "no pending mint requests" | SURVIVED (cmd whole; the only red was the known flake TestPasskeyPageAbortsAPendingCeremony/get-late, and the narrow run with that test was green) | **Real test gap**: no test lists decisions with zero mint requests, which is the exact case the production edit exists for. The code is correct as written. |
| JD-GO-LIST-LISTERR-SILENT | drop the "cannot list the decisions" stderr line | SURVIVED (cmd whole) | = plan's named `CL-LIST-ERR`. |
| JD-GO-DEC-SCANSLICE | validate the sentinel (101st) packet too | SURVIVED (daemon whole) | Gap: the 101st packet's request can be listed although outside the bound; the flag is still set and the POST re-discovers. Display only. |
| JD-GO-DEC-RECENT | `RecentApprovals(..., 40)` -> 10 | SURVIVED (daemon whole) | Gap: the 40-head read bound is not pinned (the 20-summary truncation is). |

## Per-milestone non-vacuity (against each milestone's own production diff)
- M7b listing: drop Decided/Effect/zero-packet guards, drop either truncation operand (and `||`->`&&`), scan bound 100, deadline/reserved/refusal fields, `Now`, `Carrying` all killed; null-vs-`[]` and unreadable-store covered by `TestSurfaceRequestsDecisionsUnreadable`.
- M7c `--list`: truncation-only line, deny hint, open-error line killed; the zero-mint-request path (EARLYRET) and the list-error line are not.
- M7d *Needs you*: skipping `esc()` on any of scope, requester, refusal, reserved owner or href is killed; count, truncation note, reserved/refused notification suppression, link prefix and notification href killed.
- M7e `Passkeys`: group/rp/revoked covered by `TestDecisionSubjectPasskeys` (P-PASSKEYS-GROUP killed); store error is a 500.
- M7f/g/h card: hash checks (Q, P), requestRef, expiry bytes, verb byte/text, challenge hash, user verification, origin matching, SSH verb routing, "Already decided" text all killed in live Chrome; routes/headers killed by `TestDecidePageAssets`. Surviving: FRESH-READY, WHATHASH, NOTOBJ-NULL (above).

## Findings
BLOCKING: none.

Non-blocking:
1. `cmd/ailang-worldd/approve_list_decisions_test.go`: no test of `approve --list` with zero mint requests and a pending decision (`JD-GO-LIST-EARLYRET` survives the whole package). Reproduce: add `return exitOK` after `cmd/ailang-worldd/approve.go:274`'s `fmt.Fprintln(stdout, "no pending mint requests")`; `go test ./cmd/ailang-worldd/ -skip '^TestMarketingChrome'` stays green apart from the known get-late flake.
2. `host/surfacepage/decide/decide.js:231-234` (`fresh`): the re-read that returns non-ready is not exercised (`JD-JS-FRESH-READY` survives).
3. `decide.js:193`: the displayed deny message hash is unasserted (`JD-JS-WHATHASH`); `decide.js:76-77` `null` object and `:86` `toISOString` RangeError leave the page stuck on "loading" without a message (fail-closed).
4. `host/daemon/surface_decide_list.go`: the sentinel-exclusion slice and the 40-head bound are unpinned (`JD-GO-DEC-SCANSLICE`, `JD-GO-DEC-RECENT`).
5. `approve --list` returns exit 0 with a stderr line when the world store is unreadable (documented, S4-D3).
6. Five JS operand mutants (NOW, WINOBJ, STATUS, B64TYPE, CREDOMIT) were run only against the narrow Chrome suite, not the whole package, for time.
7. Ledger hygiene: `-WHOLE` lines for FRESH-READY/WHATHASH/NOTOBJ-NULL/EARLYRET are raw "KILLED" by a baseline-red or flaky test; the `-RECLASS` lines record the real verdict.

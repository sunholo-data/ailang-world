# Evaluator report r1 — row 201 A3a sprint 3 (iteration 386), head 2936964e, base 82fb9afa

**Verdict: PASS — 90/100.** No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design and plan §9 | 28/30 |
| Tests: non-vacuous, kill the mutants | 25/30 |
| Security / fail-closed | 19/20 |
| Coding standards | 10/10 |
| Scope and size | 8/10 |

Scope-and-size note: the size is 1,843 by the stated unit, under the 2,000 cap. Two points are deducted because the plan's §5 table cannot be checked file by file here (see non-blocking 5) and because the plan files, the sprint json and QUICKSTART sit outside the numstat unit.

## Size
`git diff --numstat 82fb9afa HEAD -- cmd host`: **1810 insertions + 33 deletions = 1843**. The cap is 2,000.

The changed-file list is the 10 commits' set. Nothing is touched under `tools/`, `world/`, `packages/`, `scripts/`, `cmd/world-publish`, `host/approval`, `passkey.js` or the `/surface/passkey` route.

`surfaceCSP` is unchanged: `git diff -S surfaceCSP` is empty for the production files.

New production code has no exec site, no `/v1` route, no CORS header and no `t.Parallel`. The only `context.Background()` hits in the diff are in test files.

## Correctness and security review
**SSHSIG arm (`surface_approve.go`).**
- R11 dispatches on `body.SSHSIG != nil` and calls `authority.JudgeSSHSignature(ctx, ro, row, artifact, now)` on the read-only handle.
- The writer is opened only after the judge returns nil error. Every refusal returns before R13, so no row is left.
- The row is recorded with `format`, the armored bytes, `signed.ApprovedBy()` and the request's digest, at `d.approveNow()`.
- The store-step error still routes to the sanitized 500/503.
- Mutants G05 (skip the law), G07 (record before the law), G06 (judge at creation time) and G24 (never call the SSH judge) are all killed behaviourally.

**Page routes.**
- `GET /surface/approve/{id}` and `GET /surface/approve.js` are registered through `approvalRoute`.
- Dropping the wrapper (G18, G19) is killed by `TestApprovalRoutesRefuseForeignHost` and `TestSurfaceAppRoutes/404`.
- `setSurfaceHeaders` is kept, and the CSP is unchanged.

**`approve.js`.**
- It decodes the request bytes and hashes them itself (`sha256(bytes)` against `card.digest`).
- It checks `card.id === id`.
- It parses UTF-8 with `fatal: true`.
- It walks every leaf against the `FIELDS` table. An unlisted, mistyped or missing field refuses the card; null and non-safe ints are refused too.
- It builds the 90-byte message itself and signs `sha256(message)`.
- All DOM writes use `textContent` and `createElement`; the static test bans `innerHTML`, `eval` and the like.
- Retry fires only on `503` with `error === 'Busy'`, at most 3 times, 1 s apart.
- A POST with no answer within 10 s becomes the `unknown` state with "outcome unknown". The page offers a status check and never re-POSTs or switches format.
- `status()` reads `requests.json`.
- A passkey timeout ignores a late result and sends nothing.

**`app.js` `needsYouHTML`** escapes `episode`, `refused`, `unavailable` and, by inspection, `agent`, `id` and `recorded` too. Numbers (`grants`, `expires`) are interpolated raw but come from the daemon's typed JSON. The tests cover only some of the escapes (see non-blocking 1).

**Deviations.**
- XA2 arm (b) answering 409, because R9 precedes R11, is sound and matches T4.
- The page showing the POST's answer, with the status check reading `requests.json`, is sound because the card JSON carries no recorded approval.
- The cut to sprint 4 (XA4, the Decided list, the `N` key) is the sanctioned first step of design §9's fallback.
- The never-cut list (PG2, PG3, PG3b, XA1, XA2 arm (c), XA5, J36) is all built.
- Appendix C states the measured reason: 1,843 lines now, about 2,143 if the moved part were built.

## Test runs (my own, in `.wt-world-iter386-judge`)
`AILANG_BIN=$HOME/.pinned-ailang/ailang`, Chrome env unset except where stated.

| Run | rc |
|---|---|
| `go test -count=1 -timeout 20m ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/` (surfacepage 1.1 s, daemon 276 s, cmd 125 s; log `judge/base.txt`) | **0** |
| Live Chrome: `-run '^(TestApprovePage\|TestApproveCard\|TestNeedsYou\|TestPasskeyPage)'` with `WORLD_REQUIRE_PASSKEY_BROWSER=1` (log `judge/chrome_base.txt`) | **0** |
| `go test -race -count=1 ./host/surfacepage/` | ok (the zsh `$PIPESTATUS` echo came out empty; the `ok` line is the evidence) |
| `go vet ./host/... ./cmd/...` | 0 |
| `go test -run '^$' ./...` | 0 |

- The Chrome run passed every test: `TestApproveCardNetworkDeadlines` (79 s), `TestApproveCardPasskeyDeadline`, `TestApprovePageInChrome`, `TestApprovePageRefusesASwappedRequest`, `TestApprovePageSSHSIGInChrome`, `TestNeedsYouListsPendingRequests`, `TestPasskeyPageAbortsAPendingCeremony` and `TestPasskeyPageInChrome`.
- `TestPasskeyPageRecordsTheFixtures` is a SKIP by design (plan P4).
- No load-sensitive flake appeared.
- `git status --short` was clean after every mutant and at the end.

## Mutants (my own; ledger `$HOME/.ailang/state/world-iter386/judge/mutants.tsv`, 71 lines, one per run including re-runs; per-run outputs are `judge/out-<id>.txt`)
All kills below are by behavioural tests. Kills by a call-site inventory alone: none observed. The only grep-style kill is C03b, labelled below.

**Go (host/daemon, cmd)**

| Id | Mutation | Result |
|---|---|---|
| G01 | J36: drop the `ActivePasskeys` error check | KILLED by `TestRequestCardPasskeyReadError` |
| G02 | J36: drop the `!isOlderStore` operand | KILLED by `TestRequestCardJSON` |
| G03 | SSH dispatch `if false` | KILLED by `TestApprovePOSTRecordsAVerifiedSSHSIG` |
| G04 | judge dispatch inverted | KILLED by `TestApprovePOSTRecordsAVerifiedSSHSIG`, `TestApprovePOSTRefusalArms/R11-sshsig-*`, `TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey`, and others |
| G05 | `signed, _ =` (skip the SSH law) | KILLED by `TestApprovePOSTRefusalArms` (`R11-sshsig-garbage`, `-store`, `-unenrolled`) |
| G06 | judge at `row.CreatedAt` | KILLED by `TestApprovePOSTRecordsAVerifiedSSHSIG` and `…ViaTailnetOriginWithNoPasskey` |
| G07 | record before the law (SSH) | KILLED by `TestApprovePOSTRefusalArms` (`R11-sshsig-garbage`, `-store`, `-unenrolled`) |
| G08 | refusal text unused var | COMPILEERR (not counted); redone as G08b |
| G08b | refusal wording "passkey's assertion" for SSH | KILLED by `TestApprovePOSTRefusalArms/R11-sshsig-garbage` |
| G09 | record format forced to webauthn | KILLED by `TestApprovePOSTRecordsAVerifiedSSHSIG` and `…ViaTailnetOriginWithNoPasskey` |
| G10 | `nil` artifact | COMPILEERR; redone as G10b |
| G10b | passkey artifact `{}` | KILLED by `TestApprovePOSTRecordsAVerifiedPasskey` and others |
| G11 | SSH artifact taken from `WebAuthn` | KILLED by `TestApprovePOSTRecordsAVerifiedSSHSIG` and others |
| G12 | store step `"never"` | KILLED by `TestApprovePOSTRefusalArms/R11-sshsig-store`, `/R11-store-v5` |
| G13 | page id check `if false` | COMPILEERR (unused import); redone as G13b |
| G13b | page id check `false &&` | KILLED by `TestApprovePageRoutes` |
| G14 | page without `setSurfaceHeaders` | KILLED by `TestApprovePageRoutes` |
| G15 | script without `setSurfaceHeaders` | KILLED by `TestApprovePageRoutes` |
| G16 | script without Content-Type | KILLED by `TestApprovePageRoutes` |
| G17 | page serves the script bytes | KILLED by `TestApprovePageRoutes` |
| G18 | page route bypasses `approvalRoute` | KILLED by `TestSurfaceAppRoutes/404` |
| G18b | same, run on its own | KILLED by `TestApprovalRoutesRefuseForeignHost` |
| G19 | script route bypasses `approvalRoute` | KILLED by `TestSurfaceAppRoutes/404` |
| G19b | same, run on its own | KILLED by `TestApprovalRoutesRefuseForeignHost` and `TestApprovalRoutesRegistered` |
| G20 | MUST-LIVE: comment edit | SURVIVED, as required |
| G21 | page route removed | KILLED by `TestApprovePageRoutes` and `TestSurfaceAppRoutes` |
| G22 | page bad-id status 404 | KILLED by `TestApprovePageRoutes` |
| G23 | POST route bypasses the Host check | KILLED by `TestApprovalRoutesRegistered`, `TestApprovePOSTCrossOrigin`, and others |
| G24 | SSH judge never called | KILLED by `TestApprovePOSTRecordsAVerifiedSSHSIG` and others |
| G25 | SSH arm sets the webauthn format label | KILLED by `TestApprovePOSTRecordsAVerifiedSSHSIG` and others |
| G26 | card route ignores the error | COMPILEERR, dropped (G01 covers it) |
| G27 | e2e record format | KILLED by `TestSurfaceSSHSIGThenAwaitMints` (cmd) |

The ledger holds 71 lines: 5 are COMPILEERR, 1 is the must-live control and 1 is a vacuous run (C03c).

G02 (killed by `TestRequestCardJSON`) and G13/G13b are about the card/page handler's own checks, so they count as extra production-branch mutants beyond the plan's table.

**JS: `app.js` (`TestSurfaceAppPure` in `host/surfacepage`)**

| Id | Mutation | Result |
|---|---|---|
| A01 | `esc(q.episode)` removed | KILLED by `TestSurfaceAppPure` |
| A02 | `esc(q.agent)` removed | **SURVIVED** |
| A03 | `esc(q.id)` removed from the href | **SURVIVED** |
| A04 | `esc(q.refused)` removed | KILLED |
| A05 | `esc(q.recorded)` removed | **SURVIVED** |
| A06 | `Array.isArray` guard weakened | SURVIVED |
| A07 | clamp floor `Math.max(0,…)` weakened | SURVIVED |
| A08 | heading colour swapped | KILLED |
| A09 | empty text ignores `unavailable` | KILLED |
| A10 | plural boundary | KILLED |
| A11 | link id attribute | KILLED |
| A12 | `loadRequests` catch removed | SURVIVED (also A12's test run used `^TestSurfaceApp`) |
| A13 | MUST-LIVE: comment edit | SURVIVED, as required |

**JS: `approve.js`, static tests (`host/surfacepage`)**

| Id | Mutation | Result |
|---|---|---|
| V01 | drop the `task` FIELDS entry | KILLED by `TestApprovePageFieldTagControls` and `TestApprovePageRendersEveryCanonicalField` |
| V02 | change the `ttl_s` type | KILLED by `TestApprovePageRendersEveryCanonicalField` |
| V03 | change a label text | SURVIVED (label text is unasserted) |

**JS: `approve.js`, live Chrome**

| Id | Mutation | Result |
|---|---|---|
| C01 | skip the digest check | KILLED by `TestApprovePageRefusesASwappedRequest` |
| C02 | unshown-field guard off | KILLED by `TestApprovePageRefusesASwappedRequest` |
| C03 | `dd.innerHTML` | SURVIVED every Chrome test |
| C03b | same, run against the daemon tests | KILLED by `TestApprovePageRoutes` (static banned-token grep, which the test labels "not load-bearing") |
| C03c | same, static test skipped | invalid: zero tests ran; ignore |
| C04 | `what.innerHTML` | SURVIVED the Chrome tests |
| C05 | render only the first leaf per field | KILLED by `TestApprovePageInChrome` |
| C06 | expiry display `card.expires` | **SURVIVED** |
| C07 | unknown outcome shown as recorded | KILLED by `TestApproveCardNetworkDeadlines/post` |
| C08 | unknown outcome retries the POST | KILLED by `TestApproveCardNetworkDeadlines/post` |
| C09 | retry also on 500 (Busy operand kept) | SURVIVED, equivalent: the `Busy` operand still gates it |
| C10 | drop the `Busy` operand | KILLED by `TestApproveCardNetworkDeadlines/post` |
| C11 | `RETRIES = 4` | KILLED by `TestApproveCardNetworkDeadlines/post` |
| C12 | `RETRY_MS = 0` | KILLED by `TestApproveCardNetworkDeadlines/post` |
| C13 | `NET_MS` 100 s | KILLED by `TestApproveCardNetworkDeadlines/get` and `/post` |
| C14 | unknown outcome not styled as an error | SURVIVED (cosmetic) |
| C15 | passkey timer ×10 | KILLED by `TestApproveCardPasskeyDeadline` |
| C16 | POST digest taken from `card.challenge` | KILLED by `TestApprovePageInChrome` |
| C17 | message expiry +1 | KILLED by `TestApprovePageSSHSIGInChrome` |
| C18 | MUST-LIVE: comment edit | SURVIVED, as required |
| C19 | SSH POST digest from the challenge | KILLED by `TestApprovePageSSHSIGInChrome` |
| C20 | drop the `expires < 0` operand | KILLED by `TestApprovePageRefusesASwappedRequest` |
| C21 | retry on any 5xx | KILLED by `TestApproveCardNetworkDeadlines/post` |

Plan table rows re-confirmed by my own runs: MUT-J36 (G01), the SSH dispatch, judge dispatch, skip-law, now, record-before-law, store-step, REC-FORMAT and route rows, JS-DIGEST (C01), the retry rows (C10 to C12, C21), JS-UNKNOWN (C07), JS-POST-DIGEST (C16), JS-PASSKEY-TIMER (C15), JS-NET-DEADLINE (C13), JS-WALK-UNLISTED (C02), the APP-PURE rows and the `esc` row (A01). That is more than 15 of the plan's rows.

## Per-milestone non-vacuity
- **M6a:** removing the J36 error check (G01) is killed by `TestRequestCardPasskeyReadError`.
- **M5d:** recording before the law (G07) fails `TestApprovePOSTRefusalArms`, and the e2e arm fails on a wrong format (G27).
- **M6b:** dropping the card route from `approvalRoute` (G18) fails the Host tests. Skipping the digest check (C01) fails the Chrome XA2 test. Retry on non-Busy or on any 5xx (C10, C21) fails XA5.
- **M7a:** removing the `loadRequests` calls or the nowCard hook is killed (plan APP-LOAD, APP-NOWCARD; my A09 to A11).
- **M6c:** dropping a FIELDS entry or changing its type fails the `TestApprovePage*` field tests (V01, V02).

## Findings
**BLOCKING:** none.

**Non-blocking**
1. `TestSurfaceAppPure` leaves `esc(q.agent)`, `esc(q.id)` in the href and `esc(q.recorded)` unasserted (A02, A03, A05 survive). The code is correct by inspection, but a regression there would be an unguarded XSS vector. Fixture: `agent: '<x>'` or `recorded: '<x>'`. File: `host/surfacepage/app_test.go:180`.
2. No Chrome test renders a hostile field value. Reintroducing `innerHTML` in `approve.js` is caught only by the static banned-token grep (`surface_approve_page_test.go:69`, which labels itself non-load-bearing). Removing `textContent` for `dd` (C03) or `what` (C04) survives XA1 and XA2.
3. `approve.js:197` shows the expiry via `toLocaleString` and nothing asserts it (C06 survives). The same goes for the FIELDS label text (V03) and the unknown-outcome error styling (C14).
4. A07 (the expiry clamp) and A12 (the `loadRequests` catch) survive. These are low-impact.
5. I did not diff the sprint worktree against plan §5 line by line. The changed-file list matches §5's shape: tests, the 4 `host/daemon` files, `surfacepage` approve files, `app.js`, QUICKSTART, and the 2 plan files.
6. A passkey card with no other credentials: if the page is opened at an origin with no enrolled group, only the SSH path is offered. That is the designed behaviour and is tested by `TestApprovePageSSHSIGInChrome` (JS-NOGROUP-NOTE).

## Artefacts
- Ledger: `$HOME/.ailang/state/world-iter386/judge/mutants.tsv` (71 lines).
- Logs and per-run outputs: `$HOME/.ailang/state/world-iter386/judge/`.
- Worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter386-judge`: `git status --short` clean, no commits or pushes.

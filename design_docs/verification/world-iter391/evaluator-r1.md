# Evaluator report r1 — iteration 391, row 201 A4 sprint 3 (head e3907bbf, base 7004f591)

## Verdict: PASS — 91 / 100, no BLOCKING findings

| Category | Score |
|---|---|
| Correctness vs design and plan §9 | 28 / 30 |
| Tests: non-vacuous, kill mutants, cover operands | 26 / 30 |
| Security / fail-closed | 19 / 20 |
| Coding standards | 10 / 10 |
| Scope and size discipline | 8 / 10 |
| **Total** | **91 / 100** |

Scope/size lost 2: the cut of the whole surface (S3-D1) goes beyond the design's priced fallback (see findings), though it is honestly bounded and recorded.

## BLOCKING findings
None.

## Non-blocking findings
1. **Mint-law W15 operand survives (pre-existing A2 gap, code moved by the refactor).** `passkey_assert.go` W15 `(ad.SignCount != 0 || p.SignCount != 0) && ad.SignCount <= p.SignCount` -> drop `|| p.SignCount != 0`: SURVIVED (X-W15-OPERAND; `go test -run '^(TestDecideByAssertion|TestApproveByAssertion)' ./host/approval`). A mint assertion with counter 0 over a stored counter >0 would be accepted. The line is byte-identical to base (moved, not changed), so not a regression; sprint 4 should add the arm (and the writer's `usePasskey`, untouched, has its own tests).
2. **`hErr` placement is not pinned (X-HERR-TOP SURVIVED).** Moving `if hErr != nil {return ErrChallenge}` to the top of `assertionSteps` changes no test. The plan's A-W6-HERR (deletion) is killed. The refusal order differs only when the message hash cannot be computed AND an earlier W-step also fails; the admission / daemon window checks run first, so this is not reachable via the entry points. The "unchanged refusal order" claim holds for every reachable input (A2 matrix tests green, unedited) but the placement itself is unpinned.
3. **CLI `--signature`/`--webauthn-assertion` body `packet` is unasserted.** X-CLI-BODY-PACKET (`Packet: ...(o.ref...)` in `decideByArtifact`'s body) SURVIVED, even with `WORLD_CHROME` set: the fake daemon does not compare the body packet. Fail-closed in production (the real daemon refuses a packet that is not the request's), so a test gap, not a hole.
4. **Two read caps unpinned** (X-CLI-ASR-CAP, X-CLI-SIG-CAP survived: cap replaced by 1<<30). Low risk; `readRegularFileCapped` itself is tested.
5. **Plan-named survivors confirmed:** K0-FORMAT (genuinely equivalent: `row.Approver == format:name` is already checked by `checkKeyDecisionRow`), X-KA-BINDS (defence in depth, `DecideByAssertion` echoes its inputs, so Binds is redundant for both key arms), X-KA-RACE422 (`isApproverRace` branch not reachable through the POST, unit-pinned by `TestIsApproverRaceOperands`, JM01 killed by it).
6. **Chrome gate flake, not a regression:** the controller's gate sweep reported `TestPasskeyPageAbortsAPendingCeremony/get-late` red (aborted at 240000.1 ms vs 240000, under verify_go load). It is an A2 test, untouched here; run alone with the plan's env it passes (`ok ... 3.676s`), and `TestApproveDecideWebAuthnLinkInChrome` + `TestPasskeyPageInChrome` PASS in my run.
7. **S3-D1 cut:** goes further than design §9's priced fallback (listing and *Needs you* also moved). It is recorded in the plan (§1.3, Appendix C, priced) and every arm shipped is reachable end to end by an operator without the cut parts (terminal / `--sign-with` / `--write-message`->`--signature` / `--webauthn-link`->`--webauthn-assertion`; QUICKSTART §18 walks each). Operators cannot discover pending decisions from the surface until sprint 4 (they take the ref from the packet page), which is the real cost of the cut.
8. Minor CLI nit: `--write-message` combined with `--signature` is refused by the F1 mode count (at most one key flag), so no silent precedence exists; fine.

## Deny on the passkey page (plan §1.4 P1) — NOT blocking; named, documented defect for sprint 4
- CLI (`--webauthn-link`, deny): prints `DENY request sha256:<ref> with a passkey (expires <t>, message hash <16 hex>): ... the page reads "Approve request decide-deny-<16 hex of Q>" for either verb: check it shows message hash <16 hex>`; the link's `request` text is `decide-deny-<16hex>`, its challenge is SHA-256 of the 91-byte DENY message.
- Page (`passkey.js`, unchanged): always reads "Approve request decide-deny-…" plus the message hash (the Chrome test asserts the page shows `decide-approve-<16hex>` and `message hash <8 bytes hex>`).
- Can an operator be led to sign a verb or request they did not intend? The signature binds the verb, packet, request and expiry through the hash; the daemon recomputes the hash from the verb the CLI/POST carries (mutants X-KA-VERB, X-KA-REQUEST, X-KA-EXPIRY, PLAN-D-ASSERT-VERB/PQ all killed), so an assertion for deny posted as approve (or vice versa) is 422 with nothing decided or moved. The residual is social: a hostile link could carry any challenge and any `request` text; this is inherent to A2's page and the operator's guard is the printed hash. Documented in `approve --help` text and QUICKSTART §18 ("Approve request decide-approve-<16 hex> (or decide-deny-…) whichever verb you chose: check that it shows the message hash the CLI printed"). The wording ("Approve") on a deny is confusing, not exploitable; sprint 4's card replaces it.

## Mutation testing (ledger: `$HOME/.ailang/state/world-iter391/judge/mutants.tsv`, 67 lines; per-run logs `judge/mut_<id>.log`; summary `judge/ledger_summary.txt`)
Env copied verbatim where the plan's acceptance uses it: cmd-package runs after the first batch carried `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1`. No kill below rests on an inventory/static test: every kill is by a behavioural test named (my regexes exclude the inventory tests).

Controls: **must-live** X-CONTROL-LIVE (comment-only edit) SURVIVED as required. **must-die** X-CONTROL-DIE (W14 signature check disabled) KILLED by TestApproveByAssertionOrder, TestDecideByAssertionSteps.

My own mutants (ids absent from the plan's ledger; `comm -23 mine plan` = all 41 non-PLAN ids, `comm -12` empty), result and killer:
- X-W7-ORIGIN KILLED TestApproveByAssertionOrder, TestDecideByAssertionSteps; X-W8-CROSS same two; X-W8-TOP TestApproveByAssertionOrder; X-W10-RP, X-W11-UP, X-W12-UV, X-W13-BE, X-W5-TYPE each both; X-W15-LE TestApproveByAssertionOrder; X-W15-OPERAND SURVIVED; X-HERR-TOP SURVIVED
- X-J-VERB, X-J-REQUEST, X-J-NOW, X-J-EXPIRY KILLED TestJudgeDecisionAssertion
- X-K-USEFMT-A/-B KILLED TestRecordKeyDecisionPasskeyUseMatchesTheFormat; X-K3-EXACT (exact-retry skipping useApprover) KILLED TestRecordKeyDecisionRevokedSince
- X-KA-VERB, X-KA-EXPIRY, X-KA-REQUEST, X-KA-PUBKEY, X-KA-ARTIFACT KILLED by TestDecidePOSTPasskey (+ RetryAfterAFailedWorldWrite, AcrossRestart, SSHSIG/ BusyWriter variants); X-KA-BINDS, X-KA-RACE422 SURVIVED (named equivalent / unreachable)
- X-CLI-MSG-ORDER, X-CLI-MSG-VERB, X-CLI-WINDOW, X-CLI-LINK-RP/-ORIGIN/-HASH16/-OP, X-CLI-BODY-VERB, X-CLI-BODY-REQUEST KILLED by TestApproveDecideWebAuthnLinkAndAssertion / TestApproveDecideWriteMessageAndSignature (+ SignedBySSHKeygen, the Chrome test for BODY-REQUEST); X-CLI-BODY-PACKET(-CHROME), X-CLI-ASR-CAP, X-CLI-SIG-CAP SURVIVED
- X-JM46-INCOMPLETE-RECORD (test mutated: Plist empty) KILLED by TestApproveDecideDefaultDaemonIsTheService itself at `approve_decide_jgaps_test.go:31` "the temp service record is not read: http://127.0.0.1:7644" — the test fatals before running the CLI (HOME is a temp dir), so it never dials the real service.
- Compile errors recorded as COMPILEERR, not mutants, and rewritten: X-KA-PUBKEY (first form), PLAN-C-LINK-REQUEST (first form), X-CLI-MSG-VERB (first form), X-JM46 (first form); INVALID(anchor) rows for the first X-CLI-MSG-ORDER/-VERB.

Plan-row sample (re-run by me, 22 rows across M6a-M6g, killer named): PLAN-A-W6-HERR (TestApproveByAssertionOrder, TestDecideByAssertionSteps), PLAN-D-ASSERT-VERB (TestDecideByAssertionSteps), PLAN-D-ASSERT-PQ, PLAN-D-PKUSE-COUNT (TestDecideByAssertionDefersTheCounter), PLAN-K4-DELETE (TestRecordKeyDecisionCounter + IdempotentPasskeyRetry), PLAN-K4-EXACT (TestRecordKeyDecisionIdempotentPasskeyRetry), PLAN-JM22-FILTER (TestLandKeyDecisionLandsOnlyItsRequest), PLAN-K0-APPROVER (TestRecordKeyDecisionExactMatchNeedsTheSameApprover), PLAN-K0-FORMAT SURVIVED (equivalent, as the plan says), PLAN-J-A-REVOKED (TestJudgeDecisionAssertion), PLAN-KA-COUNTER-422 / PLAN-KA-PASSKEYUSE / PLAN-SD-ARM / PLAN-DM-FORMATS (TestDecidePOSTPasskey and the retry tests), PLAN-C-WM-EXCL (TestApproveDecideWriteMessageAndSignature), PLAN-C-LINK-HASH (TestApproveDecideWebAuthnLinkAndAssertion + the Chrome test), PLAN-C-LINK-REQUEST, PLAN-C-ASR-PARSE (TestApproveDecideWebAuthnLinkAndAssertion), PLAN-JM46 (TestApproveDecideDefaultDaemonIsTheService), PLAN-JM01 (TestIsApproverRaceOperands). All behaved as the plan expects.

Per-milestone non-vacuity: counter check dropped in RecordKeyDecision (K4-DELETE) killed; exact-retry skipping useApprover (X-K3-EXACT) killed; exact-retry NOT skipping the counter (K4-EXACT) killed; verb dropped from the decision message the CLI builds (X-CLI-MSG-VERB) killed; verb dropped from the link's request text (C-LINK-REQUEST) killed; swapped verb/packet/request in the judge or daemon killed; `decidedBy` cannot come from the body (grammar is strict: TestDecidePOSTPasskey posts `decidedBy` -> 400; the value is `signed.DecidedBy()` = `webauthn:<enrolled name>`); `hErr` swap: see non-blocking finding 2 (survives, unreachable).

## Review answers (design/security)
- Passkey arm cannot decide without the law/judge/Binds: `keyArm` runs `JudgeDecisionAssertion` (ParseAssertionResponse, store lookup by credential ID, stored key, `DecideByAssertion` P2-P5/P3w/W1-W14 over the decision hash), then `Binds`, then a non-waiting writer; a mint assertion is 422 with 0 decisions, no evidence, counter unmoved (asserted in TestDecidePOSTPasskey). A decision assertion cannot mint (hash differs; A2 mint tests unedited and green). Approve-posted-as-deny is refused by the verb in the hash.
- W15 deferral: counter judged inside `RecordKeyDecision`'s transaction via the untouched `usePasskey`, off the exact match only (`!exact && use.Passkey != nil`); a retry of the identical artifact of an unlanded row lands once without moving the counter; another assertion at the same counter is refused; a landed row refuses (K0); approver re-check (K3) and the K0 approver operand still run on the exact path (X-K3-EXACT and PLAN-K0-APPROVER killed). A failed world write leaves the evidence row and the moved counter; the retry lands exactly once (TestDecidePOSTPasskeyRetryAfterAFailedWorldWrite, ...AcrossRestart, TestDecidePOSTCounterZeroAssertionLandsOnce all green). A refused request moves nothing. A counter refusal is 422, a busy writer is 503 with 0 decisions (TestDecidePOSTBusyWriter).
- No `/v1` route, CORS header, `surfaceCSP`, exec site, clock read in a law, or `t.Parallel` added (grep of the production diff: none); `host/approval` stays standard-library-only (TestApprovalImportsOnlyTheStandardLibrary green).
- `--write-message` writes exactly the 91-byte message O_EXCL to a new file (C-WM-EXCL, C-WM-MSG killed; real `ssh-keygen` sibling green); `--signature/--webauthn-assertion` require `--expiry` and post only the verb/packet/request the CLI checked.

## Pre-existing `_test.go` edit audit (`git diff 7004f591 HEAD`, files existing at the base)
Only two, both A4 sprint-2 tests that pinned the cut:
- `cmd/ailang-worldd/approve_decide_key_test.go`: a header comment line; `TestApproveDecideSignWith`'s final loop (5 lines: `--signature` and `--write-message` refused with "the other key flags do not decide it") deleted — that behaviour is exactly what this sprint reverses; the new flags are asserted in `approve_decide_artifact_test.go`.
- `cmd/ailang-worldd/approve_decide_rig_test.go`: comment; the `--webauthn-link` row of `TestApproveDecideFlagRules` ("there is nothing to await") replaced by `--webauthn-assertion x` -> "go with --expiry <t>"; the `--await`, `--signature+--sign-with` (F1) and `--deny/--daemon` rows are intact.
No A2/A3 test edited; no assertion weakened beyond the cut's removal.

## Numstat (`git diff --numstat 7004f591 HEAD`)
`-- cmd host docs`: 1423 insertions + 70 deletions = **1,493**; `-- cmd host`: 1394 + 70 = **1,464** (matches the plan; cap 2,000).
Scope: `git diff --name-status` is exactly plan §5 plus the two plan files; nothing under tools/launchd, world/, packages/, scripts/, cmd/world-publish, skills, passkey.js, approve.js, host/surfacepage/{app,approve,passkey}/.

## Test runs (AILANG_BIN=$HOME/.pinned-ailang/ailang; logs in `$HOME/.ailang/state/world-iter391/judge/t_*.log`)
- `go test -count=1 -timeout 30m ./host/approval/... ./host/authority/ ./host/store/ ./host/sshsign/ ./host/surfacepage/` rc 0 (all ok)
- `./host/daemon/` rc 0 (274.8 s); `./cmd/ailang-worldd/` rc 0 (153.0 s)
- `./host/broker/ -run '^(TestDecide|TestDecision|TestReserved|TestPublishApproval|TestApprovalContextIdentity|TestEverySubprocess)'` ok
- `go test -race -count=1 ./host/approval/ ./host/authority/ ./host/store/` ok
- `go vet ./host/... ./cmd/...` rc 0; `go test -run '^$' ./...` no compile failures
- Live Chrome (`WORLD_CHROME=... WORLD_REQUIRE_PASSKEY_BROWSER=1`): TestApproveDecideWebAuthnLinkInChrome and TestPasskeyPageInChrome PASS; the controller's red TestPasskeyPageAbortsAPendingCeremony passes alone.
- Worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter391-judge`: `git status --short` clean after the last mutant. No go command ran before ALL_DONE.e3907bbf appeared.

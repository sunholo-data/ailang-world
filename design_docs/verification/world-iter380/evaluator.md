# Evaluator report r1 — iteration 380, row 201 A2 sprint 2 (base 3b77c294, head 1fe7c86d)

## Verdict: PASS — 93/100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design and plan §9 | 28/30 |
| Tests: non-vacuous, kill mutants, flag combinations | 29/30 |
| Security / fail-closed | 19/20 |
| Coding standards | 9/10 |
| Scope and size | 8/10 (see N1: size-unit footnote) / counted 8 |
Total: 93

## Review findings (correctness / security)
- Assertion law (host/approval/passkey_assert.go) is in design order: P1 Admit, then checkApprover (P2-P5, P3w), then W1-W15. The credential ID (W2) is checked before any signed content is parsed. The counter (W15) comes after the signature (W14), so it is judged only on a genuine signature. The counter law is `(ad!=0 || p!=0) && ad<=p`, i.e. both-zero passes and anything else must strictly advance.
- checkApprover is shared with ApproveBySignature (format parameter). P3w only applies to the webauthn law, so SSH behaviour is unchanged (mutant P3w-scope killed by the SSH-order tests).
- PasskeyUse() returns a clone and is ok only for a verified webauthn approval. The returned SignedApproval holds a clone of the credential ID (aliasing mutants killed).
- Schema 6: DDL matches passkeysSQL (a drift test exists); upgrade v1/v3/v4/v5 to v6 is idempotent (mutants upg6 and tables6 killed). EnrolPasskey: S1-S10 validation, holder-naming uniqueness pre-checks that include revoked rows, one transaction, deadline required. Lookups exclude revoked rows from ActivePasskeys and cap at 64.
- Consume: C13 (format xor Passkey) is placed before the tx and checked against the schema version (C15). The counter UPDATE runs inside the mint transaction after useApprover and is guarded by name AND credential_id AND (count < new OR both zero). RowsAffected is checked, and the C16/C17/C18 diagnostics are produced by a re-read inside the same tx. Any refusal returns and the deferred rollback undoes the request update. No path mints a passkey approval without a verified signature and a strictly advancing counter (or the both-zero case).
- Authority A6: ApprovedBy is derived from rm.Signed and the passkey binding from PasskeyUse(); the caller cannot set it.
- No fail-open path, TOCTOU or aliasing found.

## Non-blocking findings
- N1: size measured with the plan's unit (`git diff --numstat 3b77c294 HEAD -- cmd host`) = 1779 insertions + 64 deletions = **1843** (<= 2000; matches the plan's 1,843). Not 10/10 because the margin is only 157 lines.
- N2: `ApproveByAssertion` trusts `in.Approver.Passkey.Key` to be a valid parsed key (a zero Key would panic or fail in verify); it relies on NewPasskeyApprover. A nil Passkey is refused by P3w, but a zero-valued Key is not guarded. Fail-closed in practice (verify errors), minor.
- N3: the first run of my own mutants P1, C17 and A6-count did not compile (unused variable); I re-ran them in compiling form (below). Not a code defect.

## Mutants (my harness: ~/.ailang/state/world-iter380/judge/h/mut.py; each asserted `git diff --stat` non-empty before the test result, file restored after each; `go test -count=1 <pkg>`)
All 66 compiled mutants listed in ~/.ailang/state/world-iter380/judge/mutants.txt; summary (all KILLED unless noted). Each row is one file:branch; the kill test is the first failing test.
| Mutant | Location | Result (killing test) |
|---|---|---|
| P2 name check off | signed.go checkApprover | KILLED (TestApproveByAssertionOrder, TestApproveBySignatureOrder) |
| P3w drop format conjunct | signed.go | KILLED (TestApproveByAssertionOrder) |
| P3w drop nil-Passkey conjunct | signed.go | KILLED (TestApproveByAssertionOrder) |
| P3w scope (applies to ssh too) | signed.go | KILLED (TestApproveBySignatureOrder, TestSignedApprovalBinds) |
| P4 revoked off | signed.go | KILLED |
| P5 `<` to `<=` | signed.go | KILLED |
| B8 PasskeyUse format guard off | signed.go | KILLED (TestSignedApprovalPasskeyUse) |
| PasskeyUse clone removed | signed.go | KILLED |
| P1 admit off | passkey_assert.go | KILLED (TestApproveByAssertionOrder) |
| W2 credential ID | passkey_assert.go | KILLED |
| W3 userHandle: drop HasUserHandle; `&&` to `\|\|` | passkey_assert.go | KILLED x2 (AcceptsChrome et al.) |
| W5 type | passkey_assert.go | KILLED |
| W6 drop err disjunct | passkey_assert.go | KILLED |
| W7 origin | passkey_assert.go | KILLED |
| W8 drop HasTopOrigin / drop CrossOrigin | passkey_assert.go | KILLED x2 |
| W10 rpIdHash | passkey_assert.go | KILLED |
| W11 UP off | passkey_assert.go | KILLED |
| W12 UV off; UV refusal only when !BE | passkey_assert.go | KILLED x2 |
| W13 BE off; one-way; BS instead of BE | passkey_assert.go | KILLED x3 |
| W14 signature check off | passkey_assert.go | KILLED |
| W15 off; `<=` to `<`; `\|\|` to `&&`; both-zero exemption dropped | passkey_assert.go | KILLED x4 (TestSignCountLaw) |
| SA-count (p.SignCount instead of ad.SignCount); SA-alias | passkey_assert.go | KILLED x2 |
| upg6 (skip v6 upgrade); tables6 (table list) | credentials.go | KILLED x2 (TestCredentialsUpgradeToV6) |
| S1, S2, S5-used, S5-revoked, S6-ceiling, S6-floor, S8, S9-rp, S10-max | passkeys.go checkPasskeyRows | KILLED x9 (TestEnrolPasskeyRefusals) |
| pre-check credential ID; pre-check public key | passkeys.go | KILLED x2 |
| lookups: include revoked; cap 65; order DESC; miss-returns-ok | passkeys.go | KILLED x4 (TestPasskeyLookups / consume test) |
| C13 off; C13 one-way | mint_requests.go | KILLED x2 (TestPasskeyConsumeCompleteness) |
| C14 drop negative; C14 max+1 | mint_requests.go | KILLED x2 |
| C15 version off | mint_requests.go | KILLED |
| C16-C18 call off | mint_requests.go | KILLED |
| UP drop credential_id guard; `<` to `<=`; both-zero branch dropped; RowsAffected check weakened | passkeys.go usePasskey | KILLED x4 (TestPasskeyConsumeAdvancesTheCounter) |
| C16 case off; C17 case off (re-run) | passkeys.go | KILLED x2 (TestPasskeyConsumeChecksTheRowInTheTransaction) |
| A6 off; A6 counter fixed to 1 (re-run); A6 credential altered | authority/mint.go | KILLED x3 (TestMintForRequestDerivesPasskeyUse) |
| **Must-die control** W15 core `<= p.SignCount` to `< 0` | passkey_assert.go | KILLED (TestSignCountLaw, TestApproveByAssertionOrder) |
| **Must-live control** `format: FormatWebAuthn` to `in.Approver.Format` (equivalent: P3w already forces webauthn) | passkey_assert.go | SURVIVED (as expected) |
Total: ~64 mutants, 0 unexpected survivors; the plan's §4 table (139 rows) was sampled by 48 of its rows spread across all 8 sub-milestones (M3b1 ~30, M3b2/M3b3 via the same law tests, M4a 2, M4b 11, M4c 4, M5a 14, M5b 3). Worktree clean afterwards (`git status --short` empty, HEAD 1fe7c86d).

## Non-vacuity (pre-milestone checkout, `go test -run <names> -v`; output in judge/nonvac.txt)
- M3b1 at 75cb07e9: TestApproveByAssertionAcceptsChrome|TestSignedApprovalPasskeyUse: "no tests to run". 
- M4b at 0f0319d9: TestEnrolPasskeyRefusals: "no tests to run".
- M5a at 17039568: TestPasskeyConsumeAdvancesTheCounter|TestPasskeyConsumeCompleteness: "no tests to run".
All named tests exist only at/after their milestone and then kill the mutants above. HEAD restored to 1fe7c86d.

## Test runs (rc)
- `go test -count=1 ./host/approval/... ./host/store/ ./host/authority/ ./cmd/ailang-worldd/`: all ok (approval 0.7s, approvaltest 0.6s, store 11.7s, authority 0.7s, cmd/ailang-worldd 111.8s); rc 0. (judge/test1.txt)
- `go test -race -count=1 ./host/approval/... ./host/store/ ./host/authority/`: all ok, rc 0. (judge/race.txt)
- `go vet ./...`: rc 0. `go test -run '^$' ./...`: rc 0, no build failures.
- AILANG_BIN=$HOME/.pinned-ailang/ailang; no verify_go.sh, no full-repo test run.

## Scope
- `git diff --name-status` matches plan §5 (host/approval, approvaltest, host/store, host/authority, two cmd/ailang-worldd test fixtures, plan + sprint json). Nothing under tools/launchd, world/, packages/, scripts/, cmd/world-publish or skills.
- host/approval and host/approval/approvaltest import only the standard library and the host/approval package (no third-party/store imports).
- Size by plan's unit (`-- cmd host`): 1779 + 64 = 1843 (<= 2000).

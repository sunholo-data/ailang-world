# Evaluator report, iteration 390, row 201 A4 sprint 2 ("the keys"), base 8338eebb, head e4510375

## Verdict: PASS, 88/100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design and plan §9 | 27/30 |
| Tests (non-vacuous, kill mutants, cover operands) | 24/30 |
| Security / fail-closed | 19/20 |
| Coding standards | 10/10 |
| Scope and size | 8/10 |

Size is 1,983 by the unit `git diff --numstat 8338eebb HEAD -- cmd host`: insertions 1929 + deletions 54 = 1983 (cap 2,000). Production-only (non-test) is 710+28 = 738. The 1983 claim is exact. Only 17 lines of headroom, so the scope score is 8: the P2 and inventory edits sit inside that number.

## Test runs (all rc=0; outputs under $HOME/.ailang/state/world-iter390/judge/t_*.txt)
- `go test -count=1 ./host/approval/... ./host/authority/ ./host/store/ ./host/sshsign/ ./host/surfacepage/`: ok, rc=0
- `./cmd/ailang-worldd/`: ok 185 s, rc=0
- `./host/daemon/`: ok 294 s, rc=0
- `./host/broker/ -run '^(TestDecide|TestDecision|TestReserved|TestPublishApproval|TestApprovalContextIdentity|TestEverySubprocessSite)'`: ok, rc=0
- `-race` on approval, authority, store, sshsign: ok, rc=0
- `go vet ./host/... ./cmd/...`: clean. `go test -run '^$' ./...`: no compile failures.
- Worktree `git status --short` clean after every mutant and at the end.

## P2 audit: every edit to an existing `_test.go` (`git diff -U0 8338eebb HEAD` per file)
All edits match §1.4 P2. No assertion was changed or weakened, with one pin-form caveat (item 6).
1. The downgrade fixtures get exactly one added token, `DROP TABLE key_decisions`, in the statement list before the older `user_version`. That is 8 sites in cmd (await_refusals x2, passkey_refusals, passkey_test, sig_test, approver_label, approver_passkey, approver_rig), 4 in host/daemon (post_test x2, read_test x2) and 2 in host/authority/judge_test.go. No expected string or `want` changed.
2. The table-set pins (`approvers_test`, `mint_requests_test`, `credentials_test`) add `"key_decisions"` to the sorted list.
3. `credentials_test.go` changes the version pin 7 to 8 (an exact-number pin, correct) and adds closed-store entries for the three new methods.
4. The additions with no removals are the `help_test.go` usage string, a `decision_law_test.go` J07 row (plus gofmt realignment of comments), the `context_roots_test.go` deadline inventory (3 new `Background` sites, each with a reason) and the `surface_approve_post_guard_test.go` writer-site map (adds `keyArm: OpenCredentials`, i.e. 1 to 2 sites).
5. `surface_decide_discovery_test.go` only adds the SSHSIG-arm refusal (`requireKeyRefused`) beside the existing typed asserts.
6. `surface_approvals_test.go` replaces the literal 7 with `credentialsSchemaVersion` in two places. This is a mild weakening of the pin form (it can no longer fail on an unintended bump), but `credentials_test.go` still pins 8 literally. Non-blocking.

## Review notes (design checks, with evidence)
- The key arm sits in S1-D9's slot: law (`DecisionRefusal`) under `humanMu`, then `keyArm` (read-only handle, `JudgeDecisionSSH`, `Binds`, non-waiting `OpenCredentials`, `RecordKeyDecision`), then `DecideApproval`, then `LandKeyDecision`. The evidence row is committed before the decision and linked after it. A failed `DecideApproval` leaves the unlanded row, and the retry replaces or keeps it (RD13a: `TestDecidePOSTRetryAfterAFailedWorldWrite`, killed D-DECIDE-ERR and D-LAND in the plan).
- `decidedBy` is derived from the verified key (`signed.DecidedBy()`, "ssh:<name>"). The body grammar uses `DisallowUnknownFields`, so `decidedBy` and `webauthn` are 400 BadRequest, never reaching an arm (JM50 and D-DECIDEDBY killed by `TestDecidePOSTSSHSIG`). Both `typed` and `sshsig` together is refused (M-TWOFORMATS). The cut is honestly bounded: no webauthn body can reach any arm.
- The verb is inside the 91-byte message, so an approve signature cannot be posted as deny (`TestDecisionBindsTheVerb`, A-MSG-VERB). A mint-namespace signature is refused (A-NAMESPACE re-run KILLED by `TestMintSignatureIsNotADecision`, behavioural). Mint callers still sign under `world-approval@sunholo`: `Sign` calls `SignIn(approval.Namespace)`, and JM29b (Sign switched to the decision namespace) is killed by `TestSignWithSSHKeygenEndToEnd` with real ssh-keygen.
- K0 for SSHSIG compares artifact, packet ref, verb, format, approver and expiry, keyed by request ref, and runs `useApprover` on both the exact-match and replace paths (S-K3-DEL killed by `TestRecordKeyDecisionRevokedSince`). A landed row refuses (ErrKeyDecisionLanded). The `ON CONFLICT ... WHERE decision_ref IS NULL` backs it.
- S2-D2 (`--daemon` plus `/v1/health`): `SameDatabase` is checked for both `--db` and the credentials store before anything is written, and `--daemon` is not restricted to loopback. A daemon that lies in `/v1/health` could receive the typed proof header. Impact is bounded: the proof is one-use, per request and verb, and consumed only by the real local daemon. Non-blocking, since the operator chooses `--daemon`.
- The CLI refuses mismatched bytes before any prompt or proof: the request hash, the packet hash, the strict decode and `packet.requestRef` are each killed by `TestApproveDecideRefusesAMismatchedSubject` (C-PKT-HASH re-run, JM10). `--sign-with` signs `DecisionMessage(verb, P, Q, expiry)` under the decision namespace (C-SIGN-NAMESPACE re-run KILLED by the real-ssh-keygen `TestApproveDecideSignWithSSHKeygen`; JM34 digest swap killed).
- Scope: no change under tools/launchd, world/, packages/, scripts/, world-publish, passkey.js or approve.js. No new `/v1` route, no CORS or `surfaceCSP` change, no `t.Parallel`, and `host/approval` is stdlib-only with no clock read (`Now` is an input). The only new exec site is `sshsign.SignIn` (the old `Sign` site, re-parametrised), covered by `TestEverySubprocessSite...`. Only the files in plan §5 plus the two plan files changed (`git diff --name-status`). Intermediate commits were not individually built. The final tree is green.

## Non-blocking findings
1. **`LandKeyDecision`'s `request_ref` filter is untested.** JM22 (and JM22b with the daemon tests) dropped `request_ref = ?` from the UPDATE and SURVIVED. Tests only ever hold one pending row. The production code is correct. Fix: add a two-pending-rows land test. `host/store/key_decisions.go`, LandKeyDecision.
2. **K0's approver and format operands are untested.** S-K0-APPROVER and S-K0-FORMAT survive; the plan lists them as equivalent. They are equivalent through the daemon (the artifact names its own key), but not at the store API: the same artifact with a different `ApproverUse` would be kept as "exact".
3. **`isApproverRace` operands are untested.** JM01-JM04 (Revoked, KeyChanged, Unknown, ClockBehind each dropped) all SURVIVED, and D-RACE survives in the plan. The plan calls the branch unreachable through the POST. A concurrent `approver revoke` between the read-only judge and the writer does reach it, and the writer fails closed via `useApprover`. Only the status mapping (422 vs store error) is unpinned.
4. `JudgeDecisionSSH`'s `LastUsedAt` pass-through (JM27), the `isDecisionRef` prefix and length operands (JM06, JM06b: the SQL GLOB catches them) and the CLI's default service address (JM46) are unpinned.
5. D-BINDS (`Binds` in `keyArm`) is survivable by construction (equivalent). It is defence in depth, as the design wants.

## Mutant table (ledger: $HOME/.ailang/state/world-iter390/judge/mutants.tsv, 68 lines; plan ids list judge/plan_ids.txt, 125 ids)
`comm -12` (mine vs plan) = 17 ids, all deliberate re-runs of plan rows. `comm -23` = 43 ids not in the plan table (the JM* rows); every id was run at least once. Two harness bugs were fixed mid-run and logged: one invalid run (LAW-D5-KEY, wrong package, no tests ran) is marked INVALID in the ledger and re-run. All other mutants were checked non-empty (`git diff --stat`) and `go vet` clean before the result was read. Every KILLED row names a behavioural test (none is killed only by an inventory test).

Plan rows re-run (17):
| id | result (killing test) |
|---|---|
| A-NAMESPACE | KILLED TestMintSignatureIsNotADecision |
| A-BINDS-EXPIRY | KILLED TestSignedDecisionBinds |
| J-REVOKED | KILLED TestJudgeDecisionSSH |
| S-K0-APPROVER | SURVIVED (plan says equivalent) |
| S-K3-DEL | KILLED TestRecordKeyDecisionRevokedSince |
| S-K5-WHERE | SURVIVED (plan says equivalent) |
| S-UPGRADE8 | KILLED TestCredentialsSchema8KeyDecisions |
| J07-HASSUFFIX | KILLED TestReservedApprovalScope |
| LAW-D5-KEY | KILLED TestDecideRefusesReservedScopesByKey (daemon) |
| D-LAND | KILLED TestDecidePOSTSSHSIG |
| D-DECIDEDBY | KILLED TestDecidePOSTSSHSIG |
| D-BINDS | SURVIVED (plan says equivalent) |
| C-PKT-HASH | KILLED TestApproveDecideRefusesAMismatchedSubject |
| C-SIGN-NAMESPACE | KILLED TestApproveDecideSignWithSSHKeygen |
| C-SSHSIGN-NS | KILLED TestApproveDecideSignWithSSHKeygen |
| J31-ANCHOR | KILLED TestDecideRefGrammarIsAnchored |
| H-LOCK | KILLED TestDecideConcurrentPOSTsDecideOnce |

My own mutants (ids not in the plan):
| id | edit | result |
|---|---|---|
| JM05 | drop the K2 approver operand | KILLED TestRecordKeyDecisionGrammar |
| JM06 | isDecisionRef `len==71` to `>=71` | SURVIVED (equivalent: hex64 fixes the length) |
| JM06b | drop the `sha256:` prefix operand | SURVIVED (the SQL GLOB also catches it) |
| JM13 | drop `found &&` (**must-live control**) | SURVIVED, as expected |
| JM14 | DecidedBy format to a literal "ssh" (**must-live control**) | SURVIVED, as expected |
| JM17 | drop the artifact-length low bound | KILLED TestRecordKeyDecisionGrammar |
| JM18 | accept only the approve verb | KILLED TestRecordKeyDecisionReplaceUnlanded |
| JM19 | LandKeyDecision decision-ref grammar | KILLED TestRecordKeyDecisionGrammar |
| JM07 | upgrade DDL drops the `packet_ref` CHECK | KILLED TestCredentialsSchema8KeyDecisions |
| JM08 | fresh DDL drops the `verb` CHECK | KILLED TestCredentialsSchema8KeyDecisions |
| JM16 | upgrade DDL drops the `decision_ref` CHECK | KILLED TestCredentialsSchema8KeyDecisions |
| JM20 | upgrade DDL drops the `artifact` CHECK | KILLED TestCredentialsSchema8KeyDecisions |
| JM21 | fresh DDL drops the `request_ref` GLOB | KILLED TestCredentialsSchema8KeyDecisions |
| JM22 / JM22b | LandKeyDecision drops the `request_ref` filter (store, then daemon tests) | SURVIVED (finding 1) |
| JM23 | `Use()` ok to false | KILLED TestDecideBySignatureApprover |
| JM24 | signature verified over Request in the packet slot | KILLED TestDecideBySignatureApprover, TestSignedDecisionBinds |
| JM25 | Binds request compared to packet | KILLED TestDecideBySignatureApprover, TestSignedDecisionBinds |
| JM11 | unknown-key step to Law | KILLED TestJudgeDecisionSSH |
| JM26 | judge drops EnrolledAt | KILLED TestJudgeDecisionSSH |
| JM27 | judge drops LastUsedAt | SURVIVED (finding 4) |
| JM09 / JM09b | grammar drops the request hex operand (surfacepage, then daemon tests) | surfacepage SURVIVED; daemon KILLED TestDecidePOSTGuards |
| JM28 / JM28b | grammar drops the packet hex operand | surfacepage SURVIVED; daemon KILLED TestDecidePOSTGuards |
| JM29 | Sign with the decision namespace (sshsign pkg only) | SURVIVED in sshsign; JM29b killed in cmd by TestSignWithSSHKeygenEndToEnd (so not an inventory-only kill) |
| JM10 | CLI ignores `s.Refusal` | KILLED TestApproveDecideRefusesAMismatchedSubject |
| JM34 | CLI swaps the P and Q digests in the message | KILLED TestApproveDecideSignWith, TestApproveDecideSignWithSSHKeygen |
| JM43 / JM44 / JM45 | render drops or changes escalations, policy, entry | KILLED TestApproveDecideTyped |
| JM46 | default daemon address off | SURVIVED (finding 4) |
| JM35 | drop the `*list` operand in the deny rule | KILLED TestApproveDecideFlagRules |
| JM47 | `--deny` sends the approve verb | KILLED TestApproveDecideSignWith, TestApproveDecideTyped |
| JM01-JM04 | `isApproverRace` four operands dropped | all SURVIVED (finding 3) |
| JM12 | judge input Request replaced by Packet | KILLED TestDecidePOSTSSHSIG (+3 more) |
| JM48 | judge Now to 0 | KILLED TestDecidePOSTSSHSIG (+3 more) |
| JM49 | store-error step mapping | KILLED TestDecidePOSTSSHSIGRefusals |
| JM50 | **must-die control**: key arm returns `ssh:evil` | KILLED TestDecidePOSTSSHSIG |

## Per-milestone non-vacuity (production-diff mutants)
- DecideBySignature namespace: killed (A-NAMESPACE, behavioural).
- The K0 match: verb, packet, artifact and expiry operands killed (plan rows, sampled with S-K0-APPROVER); the approver and format operands are equivalent through the daemon (finding 2).
- `useApprover` skipped on the exact-match path: killed (S-K3-DEL).
- DDL CHECKs: dropped on fresh and on upgrade, all killed (JM07, JM08, JM16, JM20, JM21).
- `LandKeyDecision` skipped: killed (D-LAND). `decidedBy` from the body: impossible (strict grammar); the substitute JM50 is killed.
- CLI packet-hash check skipped: killed (C-PKT-HASH). Mint namespace in `signDecision`: killed (C-SIGN-NAMESPACE, with real ssh-keygen).

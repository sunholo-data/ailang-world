# Evaluator audit: A4 r2 narrow-refinement carve-out (iteration 388)

Doc: .wt-world-iter388-a4/design_docs/planned/w-world-approvals-by-key-a4.md (HEAD 41bf9eef; r1 2c09758b). Read-only audit.

## Score: 84/100. Carve-out: VALID (with caveats). Plannable: yes.
Correctness vs code 26/30, carve-out validity 21/25, security/R4 19/20, test non-vacuity 11/15, plannability 7/10.

## 1. Carve-out validity: VALID
Round-2 artifact (12-53-54Z): 2/2 present (gemini-4-argon, gpt6-1-sol), glm and kimi absent (unreachable). Two blocking objections, both on one surface: the judge's W15 runs on the read-only handle before K0, so a counter-N retry dies at W15 and K0 is unreachable.
- gemini-4-argon: concrete proposed_fix (look up the unlanded row before W15, skip W15 on byte-equal artifact, K0 skips useApprover/usePasskey, update RD13(b)/KD7). Disputes direction: no, it is an ordering/completeness defect.
- gpt6-1-sol: concrete proposed_fix (two-path judge, defer counter to the writer transaction, exact match bypasses only the counter, recheck revocation transactionally on both paths, mint judge unchanged, restart HTTP test with a different counter-7 assertion still failing). Disputes direction: no.
- Gate (a) and (b) both hold for every remaining objection. The "after the one re-quorum" precondition holds (round 2 was the re-quorum).
- The reviewers' fixes differ on one point (gemini skips useApprover on the exact path; gpt requires the recheck). Choosing the stricter is a legitimate reconciliation, not controller invention: gemini's text is a fix proposal whose objection concerns W15 ordering, not revocation; gpt's recheck does not contradict any gemini objection, and it is a strict superset in safety (the exact path still refuses a revoked approver). No objection is overridden. r2 quotes gpt's paragraph verbatim at doc lines 487-493.
- Caveats (non-blocking): (i) the "Concretely:" paragraph (doc 495-503) and RD14's revoked-approver arm and mutant column are controller-authored elaboration of the verbatim text (the exact-match field list, "returns the verified assertion's counter", SSHSIG note). They are consistent with the quoted fixes but are interpretation. (ii) Quorum was 2 of 4 reviewers; r2 itself was not re-quorumed (the evaluator audit stands in, as the doc states).

## 2. Does r2 close the round-2 defect? Yes, on the traced paths.
Handler order (doc 462-470 plus 495-503): D1-D10 -> read-only judge W1-W14 + P2-P5 (no W15) -> non-waiting writer -> RecordKeyDecision tx { K0 reload approver+passkey, compare unlanded row } -> exact: skip usePasskey / non-exact: W15+usePasskey, K5 -> DecideApproval -> LandKeyDecision. Code check: W15 is the last step after W14 (passkey_assert.go:88), so deferring it is cleanly separable.
(a) counter-7 retry after failed world write: judge passes (W15 gone), D6 not triggered (undecided), K0 exact match (byte-equal artifact + subject fields), no counter move, lands once. CLOSED.
(b) different assertion, counter 7, stored 7: no exact match, usePasskey SQL `sign_count < 7` false -> ErrSignCountRegression, 0 decisions (KD7 2nd arm, RD14). CLOSED.
(c) counter-0 replay: landed -> D6 409 before judge; unlanded -> exact match lands the single decision; different counter-0 assertion replaces the unlanded row (usePasskey accepts 0/0). At most one decision. OK.
(d) approver revoked between attempts: refused at the read-only P2-P5 and again in-tx. OK (but see N1 on testability).
(e) SSHSIG: no counter; full verification, exact path re-checks the approver. OK.

Consistency sweep (stale statements; none say the read-only judge runs W15 explicitly, but several are now imprecise):
- N2 doc 383 (D11): "counter not advancing -> 422 with the law's reason" is listed under the format arm/judge; under r2 it can only arise in the writer transaction. The HTTP status/text for that late refusal is unspecified (RD14 says "must still fail" with no code).
- N3 doc 305 and DM7 (doc 747): "W1-W15 steps through one unexported function" and "W6-W15 each refused" - the shared function must now support deferring W15 (a flag or split); not stated, and "W15 and usePasskey exactly as the mint consume enforces them" (doc 501) misdescribes the consume, which has only usePasskey (the SQL is equivalent to W15: both refuse ad<=stored unless both zero, so the duplication is redundant but harmless).
- N4 doc 441-442 (Judge bullet) says "the same steps ... as judge.go" with no mention that W15 is omitted; doc 462-466 "Order inside the handler" and the RecordKeyDecision step list (K1,K2,K3+K4,K5) do not name K0 or the exact-match flag. The r2 text is only in the bullet at 483-509. A planner must assemble it.
- N5 sizing/gates not updated: header doc 9-10 and 39 (sprint 2 "+200 for K0, KD5-KD8, RD13") omit RD14 (restart test + revoked arm, roughly +100 lines); M5 list (doc 828) says RD10-RD13 (the -run '^TestDecidePOST' pattern does match RD14's name, so the gate still runs it); fallback cut (doc 831) does not move RD14 with the webauthn arm; M4 lists KD1-KD7 while the exact-match logic sits in M4 but RD14 is in M5. No sprint is wrong by more than ~5%.
- §3 (doc 641-643, 669) and T15 reference K0 correctly (T15 lists KD7, KD8, RD13, not RD14: trivial).

## 3. Round 1 closure: yes
(1) D4 complete discovery: D4a/D4b, Discovery{Complete,...}, RD9b-RD9d, V33/V34 present (doc 346-370, 397-399). (2) key_decisions DDL CHECKs (packet_ref, decision_ref `sha256:*`, doc 447-453) and the unlanded-row retry (KD5-KD8). Round 2 raised nothing on either.

## 4. V-row re-measurement (worktree, own commands)
- V3 MATCH: DecideApproval defs/callers approve.go:181,:188,:191,:696, publish_op.go:153; control test callers = 5.
- V4 MATCH: approvalRequestWire fields effect/scope/cost/requester/now at approve.go:66-72; DisallowUnknownFields :322 (doc's `:320-327` range); record.go:132-140 Hash SumSHA256(payload).
- V5 MATCH: main.go:310 MintAttendedApproval; episode_finish.go:173 "operator:episode-finish","operator:terminal"; defs approve.go:686, publish_op.go:125.
- V6 MATCH: only daemon.go:954 EffectHumanApprove; PollApproval no match (control: the Approve line found by the same grep).
- V11 MATCH: approve.go:544 Decision!="approve", :578 Effect, :583 Cost!=PublishCost, PublishCost=1 at registry_publish.go:31; no DecidedBy in lines 517-650 (control: same range has 5 "Decision" hits; DecidedBy appears elsewhere only at :77, :214).
- V32 MATCH: passkeys.go UPDATE with `(sign_count < ? OR (sign_count = 0 AND ? = 0))`, C16-C18 cases, passkey_assert.go:88 W15.
- V33 MATCH: decisionpacket_validate.go:34, :190, :194, :195; reader workbench_packet.go:78.
7 rows re-run, 0 mismatches (one minor range note: `:320-327` vs the actual `:322` token).

## 5. R4 by construction: holds
validatePublishApproval (approve.go:517) reads Decision, Effect, Cost, Scope only; no DecidedBy (grep above with control). Layers: L1 (reserved-scope refusal in the decision law, D5) and L2 (new check refusing decidedBy with ssh:/webauthn: prefix). Each alone stops a key decision from authorizing a publish: L2 works because the handler derives decidedBy from SignedDecision.DecidedBy() = format:name (formats only ssh/webauthn, V19), so every key decision carries the prefix; L1 refuses the request before any decision exists. Tests: R4a (reserved classifier; near-miss controls), R4b (L2 alone: key-decided publish-shaped request refused; must-die: L2 deleted), R4c (must-live: typed stamp / free-text name still publishes, registry_publish_test arms unedited), R4d (L1 compiled out, L2 still refuses; typed 403), R4e (key POST 403; L1 out -> decision lands, consumer refuses). Both layers have a must-die and the must-live control exists. Residual: L2 is prefix-based on free text; a typed `world-publish --decided-by ssh:x` is also refused (disclosed, V19 measured zero users).

## 6. Plannability: good
Three sprints (~1,950/1,800/1,500, V31 calibration ~1.05x), per-milestone gates with concrete `-run` patterns, fallback cuts named with never-cut items. P1 (91-byte message in Chrome crypto.subtle) is a first step of sprint 3 on the controller's machine; dischargeable but needs Chrome (RK8 discloses). K0/exact-match is under-specified in structure (N3-N5) but a planner can resolve it from the r2 bullet.

## Findings
Blocking: none.
Non-blocking (doc line -> evidence):
- N1 (test non-vacuity, most significant) doc 761 RD14: the revoked-approver arm cannot kill its stated mutant "exact-match path skipping the transactional approver re-check" end to end, because the read-only judge's P2-P5 already refuses a revoked approver before the writer opens. Needs a store-level test of RecordKeyDecision (K0 exact path with approver revoked after the judge) or a seam between judge and writer. Same for the passkey-reload (credential changed) arm.
- N2 doc 383: late counter refusal status/text unspecified.
- N3 doc 305, 501, 747: shared W1-W15 function needs a defer-W15 form; "as the mint consume enforces them" misdescribes (consume has no W15).
- N4 doc 441, 462-470, RecordKeyDecision steps: K0 and the no-W15 judge not integrated into the Judge bullet or the handler-order paragraph; r2 lives in one superseded-marked bullet.
- N5 doc 9-10, 39, 828, 831, 669: RD14 missing from header, size estimate, M5 list, fallback cut and T15.
- N6 doc 495-503, 761: controller-authored elaboration beyond verbatim text (see section 1 caveat); acceptable but should be labelled as the controller's reading.

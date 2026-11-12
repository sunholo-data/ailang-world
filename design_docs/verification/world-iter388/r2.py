import sys,re
p=sys.argv[1]; s=open(p).read(); n0=len(s)
def rep(old,new):
    global s
    c=s.count(old)
    if c!=1: sys.exit(f"anchor count {c} for: {old[:80]!r}")
    s=s.replace(old,new)
# 1. status line
rep("- **Status**: Planned. Designed in iteration 388, **r1**. r0 was **BLOCKED** at quorum round 1\n  (2/2 present; §15).",
    "- **Status**: Planned. Designed in iteration 388, **r2** (r2 by the controller's narrow-refinement\n  carve-out: quorum round 2 BLOCKED 2/2 on one surface, the K0 retry order; both reviewers' fixes applied\n  in their own text, marked **[r2]**; §15). r0 was **BLOCKED** at quorum round 1\n  (2/2 present; §15).")
# 2. K0 bullet
start=s.index("  - **K0 (new, inside the handler, before K3/K4).**")
end=s.index("  - **Otherwise** (no row, or a different artifact)")
old=s[start:end]
new='''  - **K0 (new, inside the handler, before K3/K4).** **[r2] Superseded by the two-path judge below:
    r1 put the judge's W15 counter check before K0, so for a hardware-counter passkey the retry
    died at W15 (stored `N`, assertion `N`) before K0 could recognise the artifact — K0 was
    unreachable on exactly the path it was added for (quorum round 2, both reviewers).**
  - **[r2] The two-path decision judge (gpt6-1-sol's fix, verbatim):** "Replace the retry paragraph
    with an explicit two-path decision judge: verify the signature, challenge binding, origin, RP
    hash, UV, and approver identity/status first, but defer counter acceptance. On the non-waiting
    credentials writer, transactionally reload the approver and passkey and compare the existing
    unlanded row against the exact artifact and all signed subject fields. Only an exact match to
    previously committed evidence may bypass counter advancement; otherwise enforce the existing W15
    and usePasskey rules. Recheck revocation transactionally even on the exact-match path. Leave the
    mint judge unchanged."
    Concretely: `JudgeDecisionAssertion` runs W1–W14 and P2–P5 on the read-only handle and does **not**
    run W15; it returns the verified assertion's counter. Inside `RecordKeyDecision`'s transaction:
    K0 reloads the approver (`useApprover`'s active/unrevoked check runs on both paths) and the
    passkey row, then — **exact match** (an unlanded row, `decision_ref IS NULL`, whose `artifact` is
    byte-equal to the POSTed one and whose `verb`, `packet_ref`, `request_ref`, `expiry`, `format`
    and `approver` equal the verified ones) → skip `usePasskey`, no counter move, go to
    `broker.DecideApproval` + `LandKeyDecision`; **no exact match** → W15 and `usePasskey` exactly as
    the mint consume enforces them, then K5. SSHSIG has no counter: its exact-match path still
    re-checks the approver. gemini-4-argon's equivalent fix ("check for an existing unlanded
    `key_decisions` row … BEFORE running W15 in the judge … skips W15 … and K0 then skips
    `useApprover`/`usePasskey`") is the same ordering; the one difference — gemini also skips
    `useApprover` — is resolved by gpt6-1-sol's "Recheck revocation transactionally even on the
    exact-match path", which is the stricter of the two and is kept. The mint judge (`judge.go`) is
    unchanged.
    The law D1–D10 still runs first on every retry, and D6 still refuses if the decision already
    landed.
'''
s=s[:start]+new+s[end:]
# 3. KD7 row: add HTTP requirement note
rep("| KD7 **[r1]** | `TestRecordKeyDecisionIdempotentRetry` |",
    "| KD7 **[r1] [r2]** | `TestRecordKeyDecisionIdempotentRetry` |")
rep("| RD13 **[r1]** | `TestDecidePOSTRetryAfterAFailedWorldWrite` (`host/daemon`) |",
    "| RD13 **[r1] [r2]** | `TestDecidePOSTRetryAfterAFailedWorldWrite` (`host/daemon`) |")
# insert RD14 after RD13 row line
i=s.index("| RD13 **[r1] [r2]**"); j=s.index("\n",i)
rd14=("\n| RD14 **[r2]** | `TestDecidePOSTRetryAfterAFailedWorldWriteAcrossRestart` (`host/daemon`) | gpt6-1-sol's test, verbatim: "
      "\"commits stored counter 7, injects a world-write failure, restarts the daemon, and retries the identical assertion successfully without moving the counter; a different assertion carrying counter 7 must still fail\"; "
      "and gemini-4-argon's: \"the full HTTP/judge path succeeds on a counter-7 retry when `webauthn_credentials.sign_count` is already 7 in SQLite\". Plus: the approver revoked between the failed attempt and the identical retry → refused, 0 decisions | "
      "the read-only judge running W15 before the exact-match check (the round-2 defect: the retry dies at W15); the exact-match path skipping the transactional approver re-check (a revoked approver's retry lands); exact match on the artifact alone without the signed subject fields |")
s=s[:j]+rd14+s[j:]
# 4. section 15
rep("- **Round 2:** *(controller fills)*",
'''- **Round 2 (2026-11-12T12:53:54Z): BLOCKED 2/2 present (gemini-4-argon, gpt6-1-sol); oc-glm-5-3 and oc-kimi-k3 absent (ollama 429 weekly limit, unreachable).** Artifact `/Users/voightkampff/.ailang/state/mission-quorum-world/w-world-approvals-by-key-a4-2026-11-12T12-53-54Z.json`. Both objections land on ONE surface, §2.7's retry order (r1's K0 sat after the judge's W15, so a counter-N retry could never reach it); round 1's discovery fix drew no objection. Neither disputed the direction; both carried a concrete `proposed_fix`.
  - **r2, by the controller's narrow-refinement carve-out (gate 2), not a designer run:** §2.7's K0 bullet replaced by gpt6-1-sol's two-path judge in its own words (W15 deferred to the writer transaction; exact match bypasses only the counter; revocation re-checked on both paths; mint judge unchanged); gemini-4-argon's equivalent ordering cited, its `useApprover` skip not taken because gpt6-1-sol's stricter revocation re-check is kept; KD7/RD13 marked; RD14 added with both reviewers' test text verbatim plus a revoked-approver arm. No other section changed. An independent evaluator audits this carve-out before any plan.''')
open(p,"w").write(s); print("ok", n0, len(s))

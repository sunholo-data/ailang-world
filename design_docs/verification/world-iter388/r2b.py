import sys
p=sys.argv[1]; s=open(p).read()
def rep(old,new):
    global s
    c=s.count(old)
    if c!=1: sys.exit(f"anchor count {c}: {old[:70]!r}")
    s=s.replace(old,new)
rep("    Concretely: `JudgeDecisionAssertion` runs",
    "    Concretely (the controller's reading of that fix, audited VALID — §15): `JudgeDecisionAssertion` runs")
rep("→ W15 and `usePasskey` exactly as\n    the mint consume enforces them, then K5.",
    "→ W15 (the A2 law's counter rule, `passkey_assert.go:88`) and `usePasskey`, then K5. The A2\n    verifier gains a defer-W15 form for this judge only; the mint path keeps calling the full W1–W15.\n    A counter refused here, in the writer, answers **422** with the same `ErrSignCountRegression` text\n    the read-only judge uses, and writes nothing (the transaction rolls back).")
rep("- **M5 the key POST arms and the key CLI** (≈ 230 + 520 **[r1]**): RD10–RD13, KD8,",
    "- **M5 the key POST arms and the key CLI** (≈ 230 + 520 **[r1]**, + ≈ 60 **[r2]**): RD10–RD14, KD8,")
rep("(KD7, KD8, RD13) |", "(KD7, KD8, KD9, RD13, RD14) **[r2]**: W15 deferred to the writer so the exact-match retry is reachable |")
# N1: replace RD14's revoked arm with a store-level KD9
rep(" Plus: the approver revoked between the failed attempt and the identical retry → refused, 0 decisions |",
    " |")
rep("the exact-match path skipping the transactional approver re-check (a revoked approver's retry lands); exact match",
    "exact match")
i=s.index("| KD8 **[r1]**"); j=s.index("\n",i)
kd9=("\n| KD9 **[r2]** | `TestRecordKeyDecisionExactMatchRechecksTheApprover` (`host/store`) | an unlanded row from counter 7; the approver revoked **after** the read-only judge and **before** `RecordKeyDecision` (store level: the judge is not in this test, so nothing else can refuse) → the exact-match call refuses with the revoked-approver error, counter still 7, row still unlanded, 0 decisions | the exact-match path skipping `useApprover` (gemini-4-argon's variant: a revoked approver's retry lands) |")
s=s[:j]+kd9+s[j:]
rep("RD14 added with both reviewers' test text verbatim plus a revoked-approver arm.",
    "RD14 added with both reviewers' test text verbatim. After the independent audit (evaluator sonnet, CARVE-OUT VALID 84/100, zero blocking): its N1 (an end-to-end revoked-approver arm cannot kill its mutant, because the read-only judge refuses a revoked approver first) moved that arm to the store-level KD9; N2/N3 the writer's late counter refusal is a 422 with the judge's text, and the A2 verifier gains a defer-W15 form; N5 RD14/KD9 added to M5 and T15; N6 the controller's elaboration is labelled. N4 (fold K0 into the Judge and handler-order bullets) is left to the planner: the two-path bullet sits directly under them.")
open(p,"w").write(s); print("ok")

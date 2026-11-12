MISSION-ROLE: evaluator

You are the independent EVALUATOR (generator ≠ judge) for Ailang World mission iteration 388,
unattended. You did not write the document you are judging; the designer was claude-opus-5-5 and the
controller wrote the r2 change. Your job is an AUDIT of a narrow-refinement carve-out, scored /100.

READ-ONLY. Work in `/Users/voightkampff/dev/sunholo-data/.wt-world-iter388-a4` (branch
`design/world-iter388-row201-a4`, base `b0fbcd3e`). Do NOT edit any tracked file and do NOT run any
git write. Scratch output goes to `/Users/voightkampff/.ailang/state/world-iter388/judge/`.

## The document and its history
- `design_docs/planned/w-world-approvals-by-key-a4.md` — the A4 child design of queue row 201.
  Commits: r0 `1fdcade6`, r1 `2c09758b` (designer), r2 = HEAD (controller). See them with
  `git -C <wt> log --oneline -4` and `git -C <wt> diff 2c09758b HEAD`.
- Parent (attended, rulings D-WORLD-89 R1–R5 = yes): `design_docs/planned/w-world-approvals-by-key.md`.
- Quorum round 1 artifact: `/Users/voightkampff/.ailang/state/mission-quorum-world/w-world-approvals-by-key-a4-2026-11-12T12-44-28Z.json`
- Quorum round 2 artifact: `/Users/voightkampff/.ailang/state/mission-quorum-world/w-world-approvals-by-key-a4-2026-11-12T12-53-54Z.json`
  Read reviewer verdicts at `.reviewers[].result` and absences at `.synthesis.absent_reviewers`.

## The carve-out rule you are auditing (gate 2, quoted)
"AFTER the one re-quorum, IF **every** remaining blocking objection (a) carries a concrete
reviewer-authored `proposed_fix` AND (b) does NOT dispute the design DIRECTION (only completeness /
determinism / attribution / a scope-cut of a non-core helper), the controller MAY make a **bounded
2nd revision that applies the reviewers' VERBATIM fixes** (their own text — never a
controller-invented resolution, never overriding an objection) and route straight to sprint-planner
… If ANY remaining objection disputes the direction, or lacks a concrete fix, or would need
controller judgment to resolve → park."

## What to produce
1. **Carve-out validity** (VALID / INVALID with reasons): for each round-2 objection, does it carry a
   concrete fix, does it dispute direction, and does r2 apply that fix in the reviewer's own text?
   The two reviewers' fixes differ in one point (gemini also skips `useApprover` on the exact-match
   path; gpt6-1-sol requires a transactional revocation re-check). The controller kept the stricter
   one. Judge whether that choice is "controller-invented" or a legitimate reconciliation of two
   reviewer texts, and say why.
2. **Does r2 actually close the round-2 defect?** Trace the handler order in §2.7 for (a) a
   counter-7 passkey retry after a failed world write, (b) a different assertion with counter 7,
   (c) counter-0 replay, (d) an approver revoked between attempts, (e) SSHSIG. Check that every
   other section that describes the judge or the handler order (§2.7 Judge bullet, "Order inside the
   handler", §3, §4 T15, §5, §6 rows JD*/KD7/RD13/RD14, §9 gates and `-run` patterns) is consistent
   with r2. Any section that still says the read-only judge runs W15 is a finding.
3. **Did r1 close round 1?** Briefly confirm the two round-1 objections are addressed (round 2 did
   not repeat them).
4. **Independent re-measurement:** re-run at least 5 §14 V-rows (include V3–V6, V11, V32, V33) with
   their own commands and say whether the observed output matches the doc. Pair any empty result with
   a known-positive control in the same command.
5. **R4 by construction:** confirm from code that `validatePublishApproval` (host/broker/approve.go
   ~517) does not read `DecidedBy`, and that the design's two layers (L1 reserved-scope refusal, L2
   publish-consumer refusal of `ssh:`/`webauthn:` decidedBy) would each alone stop a key decision
   from authorizing a publish. Name the test rows that prove each and their must-die/must-live
   controls.
6. **Plannability:** are the sprint sizes (≈1,950 / 1,800 / 1,500) and fallback cuts concrete enough
   for a planner? Is any obligation left that a planner cannot discharge (e.g. browser premises P1)?
7. **Score /100** (correctness of the design vs code 30, carve-out validity 25, security/R4 20,
   test plan non-vacuity 15, plannability 10) with **blocking** findings (would make the plan wrong)
   separated from **non-blocking** ones. For each finding give file:line in the doc and the evidence.

## Process rules (you run headless; nothing will re-invoke you)
- Do NOT background anything; run commands in the foreground. A `Monitor` is an event stream, not a
  wait; never end your turn announcing that you are waiting. If something cannot finish, mark it
  UNMEASURED and finish the report.
- Use absolute paths or `cd <abs> && …` in the same command.
- Budget ~25 minutes.
- Write your full report to `/Users/voightkampff/.ailang/state/world-iter388/judge/audit.md` AND return
  it as your final reply (score, VALID/INVALID, blocking list, non-blocking list, V-rows re-run with
  results).

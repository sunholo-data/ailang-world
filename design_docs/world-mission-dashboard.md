# World mission dashboard — iteration 386 (2026-11-12)

- **Landed this fire:** row 201 A3a sprint 3, "the page and the app" (PR #558, merge `70823879`). CI was 3/3 on the head at the first attempt. On the merge, the race gate went red once on a 1 s teardown bound in untouched MCP-budget code; the same tree had passed at the head (a `rerun --failed` went 3/3). Judge sonnet PASS 90, zero blocking.
  - The card page `GET /surface/approve/{id}` and `approve.js`, behind the Host check, with the CSP unchanged. The page hashes the stored bytes itself and renders every signed field. It refuses unknown, missing or mistyped fields and writes only through `textContent`.
  - The card has 10 s network deadlines and retries only on 503 Busy (at most 3 times). An unanswered POST is shown as "outcome unknown".
  - *Needs you* in the surface app lists pending requests and links to their cards.
  - The SSHSIG arm of `POST /surface/approve/{id}`: `JudgeSSHSignature` at `now`, then one `ssh:<name>` row; `approve --await` mints it.
  - The J36 gap is closed: the card route answers 500 when the passkey read fails.
  - Live-Chrome tests XA1, XA2 (a–c), XA3's list half and XA5 run on the rig (CI has no Chrome). QUICKSTART §17 is written (attended, pending).
  - The design's first fallback step fired: the opt-in notification, the Decided list and the `N` key move to sprint 4.
- **Next loop work:** A3a sprint 4, from the sprint-3 plan's Appendix C:
  - the opt-in notification (XA4);
  - the Decided list and the `N` key;
  - the judge's test gap: assert that the Needs-you list escapes `agent`, the `id` href and `recorded`.
- **Attended, pending Mark:** A1b exit run (QUICKSTART §15); A2 real-passkey exit run (A2 §10, QUICKSTART §16); A3a exit run (QUICKSTART §17, after sprint 4).
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (default A, applied).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150/151 not routable.
- **Routing this fire:** planner opus (Agent, 97 min), executor haiku (Agent ×2, 15 min), evaluator sonnet (Agent, 24 min); no designer, no quorum. Codex unusable at preflight (rc 75). Metered $0.00.
- **Instrument notes:** the self-notice script reports `watermark channel unavailable` until a directive is marked on #524. Row 258 gains two CI-flake instances: the 385 record merge (`host/verifygate`, row 249's class) and this merge (`TestMCPPostBudgetProductionConstants/batch`, first instance).
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).

# World mission dashboard — iteration 387 (2026-11-12)

- **Landed this fire:** row 201 A3a sprint 4, the last A3a sprint (PR #560, merge `4134ce89`). CI was 3/3 on the head and on the merge, both at the first attempt. Judge sonnet PASS 90, zero blocking.
  - The Decided list under *Needs you* (the last 3, and "Show all" for the 20 in `requests.json`).
  - The opt-in browser notification. "Notify me" asks for permission only from its click. With permission granted, each new pending request raises one notification per tab, and clicking it opens the card.
  - The `N` key jumps to what needs you, and Enter on a now-card link follows it.
  - The sprint-3 judge's test gap is closed: every `esc()` in the Needs-you list is now pinned by `TestSurfaceAppPure`.
  - A live-Chrome test shows a hostile signed-field value as text on the card. QUICKSTART §17 gains its notification step.
  - The only production change is `app.js`: no new route, and the CSP is unchanged.
- **A3a is code-complete** except Mark's attended exit run (QUICKSTART §17).
- **Next loop work:** row 201 A4, `Human.Approve` decision packets approved the same way. It needs a new child design by the rotation designer (next is `codex:gpt-6.1-sol`, then glm), then the quorum and the plan.
- **Attended, pending Mark:** A1b exit run (QUICKSTART §15); A2 real-passkey exit run (A2 §10, QUICKSTART §16); A3a exit run (QUICKSTART §17).
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (default A, applied).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150/151 not routable.
- **Routing this fire:** planner opus (Agent, 85 min), executor haiku (Agent ×2, 5 min), evaluator sonnet (Agent, 8 min); no designer, no quorum. Codex unusable at preflight (rc 75). Metered $0.00.
- **Instrument notes:**
  - The self-notice script reports `watermark channel unavailable` until a directive is marked on #524.
  - Row 258 gains the 386 record-merge race red (`host/surfacefollow` `TestBackoff`, first instance).
  - The planner saw a new untouched-package flake once: `host/broker` `TestPublishApprovalRefusalSetWithALandedPositiveControl`.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).

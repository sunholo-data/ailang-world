# World mission dashboard — iteration 388 (2026-11-12)

- **Done this fire:** the row 201 A4 child design, r2: [`w-world-approvals-by-key-a4.md`](planned/w-world-approvals-by-key-a4.md). No code.
  - An operator approves or denies an agent's `Human.Approve` request (the row-172 decision packet) with a typed `y`, an enrolled SSH key or a passkey, from the CLI (`approve <request-ref> [--deny]`) or the surface.
  - The gap was measured first: no production code calls `broker.DecideApproval`, so these packets could never be decided.
  - The daemon decides and derives who decided. A key decision can never authorize a publish: there are two separate fences, because the publish check never reads who decided.
  - Three sprints, about 5,250 lines.
- **Quorum:** blocked in rounds 1 and 2, with 2 of 3 seats present (gemini-4-argon and gpt6-1-sol; glm and kimi absent on ollama's weekly limit). r2 was made under the controller's carve-out. The independent audit (sonnet) found the carve-out VALID, 84/100, zero blocking. Metered $0.36.
- **New decision:** D-WORLD-115 asks whether to keep the A4 decision message without the parent's stored random challenge. Default A, applied from iteration 389.
- **Next loop work:** the A4 sprint-1 planner ("the law and the typed root"), with no re-design.
- **Attended, pending Mark:** the A1b exit run (QUICKSTART §15), the A2 real-passkey exit run (A2 §10, QUICKSTART §16) and the A3a exit run (QUICKSTART §17).
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (default A, applied), D-WORLD-115 (default A, from 389).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:** designer opus (Agent; r0 17 min, r1 4.5 min), evaluator sonnet (Agent, 2 min). No planner or executor. Codex was over ration (probe rc 75).
- **Instrument notes:**
  - The self-notice script reports `watermark channel unavailable` until a directive is marked on #524.
  - Row 258 gains the 387 record-PR race red (`marketingcheck` shell-wait).
  - macOS has no `timeout`, so one quorum launch exited rc 127 having run nothing.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).

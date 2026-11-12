# World mission dashboard — iteration 389 (2026-11-12)

- **Done this fire:** row 201 A4 sprint 1 LANDED ([#563](https://github.com/sunholo-data/ailang-world/pull/563), merge `338ff821`, CI 3/3 head and merge first try). Judge (sonnet) PASS 88/100, zero blocking.
  - A `Human.Approve` decision packet now has a pure refusal ladder (D1–D10), a complete request-wide packet discovery, and the 91-byte `world-decision/v1` message.
  - A request is decided at most once. A key-made decision can never authorize a publish (fences L1 and L2).
  - The daemon serves `GET /surface/decisions/{ref}` and accepts a typed decision at `POST /surface/decide/{ref}` behind the A3a guards. It derives `decidedBy` itself and consumes the verb-bound operator proof first.
  - The typed CLI (`approve <request-ref> [--deny]`) moved to sprint 2 by the design's priced fallback, so an operator cannot decide from a terminal yet. The sprint is 1,985 of 2,000 lines.
- **Judge gaps carried to sprint 2 (non-blocking):** J07 (the L1 effect arm), J31/J32 (the D1 regex anchor and case at the daemon; reproduced by the controller), J33 (the GET's D4a mapping).
- **Next loop work:** the A4 sprint-2 planner ("the keys"), typed CLI first. Expect the webauthn fallback cut.
- **Attended, pending Mark:** the A1b exit run (QUICKSTART §15), the A2 real-passkey exit run (A2 §10, QUICKSTART §16) and the A3a exit run (QUICKSTART §17).
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (default A, applied), D-WORLD-115 (default A, applied from 389).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:** planner opus (Agent, 105 min), executor haiku (Agent ×2; run A resumed once after ending its turn on a wait), evaluator sonnet (Agent, 50 min). No designer. Codex was over ration (probe rc 75). Metered $0.00.
- **Instrument notes:**
  - The rig ran at load 40–210. Local `verify_go` hit five 10-minute package timeouts; every red was on a known list or green alone, and clean-runner CI was the gate.
  - `TestWorkspaceProtectedPolicyWiring` (host/daemon) and `TestExecProtectedCheckTimeout/reserve-after-before` (host/broker) are new load reds, each green alone.
  - The self-notice script still reports `watermark channel unavailable` until a directive is marked on #524.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).

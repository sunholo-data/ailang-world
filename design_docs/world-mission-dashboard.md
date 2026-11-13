# World mission dashboard — iteration 390 (2026-11-12)

- **Done this fire:** row 201 A4 sprint 2 LANDED ([#566](https://github.com/sunholo-data/ailang-world/pull/566), merge `9ea7b45f`, CI 3/3 head first try; merge 3/3 on attempt 2 after the known filter-control flake). Judge (sonnet) PASS 88/100, zero blocking.
  - An operator can approve or deny a `Human.Approve` decision packet at a terminal: `ailang-worldd approve <sha256:request-ref> [--deny]` checks the subject's hashes, takes `y`, and posts the verb-bound proof.
  - An enrolled SSH key can decide too: `--sign-with` signs the 91-byte decision message under `world-decision@sunholo`. The daemon judges the key, records it in credentials schema 8 `key_decisions`, decides, then lands the row.
  - A key decision still cannot authorize a publish (L1 + L2). Concurrent POSTs decide once, and a busy writer gets 503 with 0 decisions.
  - The passkey arm and `--write-message`/`--signature` on a request ref moved to sprint 3 by the design's priced fallback. The sprint is 1,983 of 2,000 lines.
- **Judge gaps carried to sprint 3 (non-blocking):** JM22 (`LandKeyDecision`'s request filter; reproduced by the controller), K0's approver/format operands at the store API, JM01–JM04 (`isApproverRace` → 422), JM27, JM46.
- **Next loop work:** the A4 sprint-3 planner. The webauthn arm and CD6–CD7 come first, then the surface (decision listing, card, live Chrome, QUICKSTART §18). It is ≈ 2,060 lines, so expect its own fallback cut.
- **Attended, pending Mark:** the A1b exit run (QUICKSTART §15), the A2 real-passkey exit run (QUICKSTART §16) and the A3a exit run (QUICKSTART §17).
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (default A, applied), D-WORLD-115 (default A, applied).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:** planner opus (Agent, 77 min), executor haiku (Agent ×3, all clean), evaluator sonnet (Agent, 18 min). No designer. Codex was unusable at preflight (probe rc 75). Metered $0.00.
- **Instrument notes:**
  - The rig ran at load 16–30. Every leg of the §6 gates was green, and `verify_go` plain gave rc 0 on its first run.
  - The self-notice script still reports `watermark channel unavailable` until a directive is marked on #524.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).

# World mission dashboard — iteration 306 (2026-10-27)

- **This fire:** row 150. R0 re-confirmed the shot bottoms (V59), and M7c0 (tall shot geometry) is BUILT in PR #373.
- **R0:** the rig's Chrome moved from 156.0.8078.47 to .48. Two fresh seeds gave exactly the V51 bottoms (packet 3564.56, graph 2042.14, provenance 2881.17; delta 0.00). A third seed, run by the judge, matched.
- **M7c0:** `Shots` heights are now packet 3700, graph 2100, provenance 3000, pinned by `TestMarketingShotGeometry`. The packet image now contains the validated packet and all four roster rows.
- **Judge:** sonnet PASS 96/100, zero blocking; 11/11 of its own `png.go` mutants killed.
- **Roles (Agent tool):** executor codex gpt-6.1-sol through an Agent supervisor running the pinned recipe; evaluator sonnet. Designer and planner were not needed (the doc and the prescriptive plan exist).
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81) is 156→150→151.
- **Attended rulings acknowledged:** D-WORLD-91/92/93 (rows 188/151), provenance `attended`. Row 151 now builds on D-WORLD-92/93.
- **Row 156:** 22/24; D-WORLD-86 holds M5c/M5d.
- **Row 150:** 16/22 built. Pending: M7c2a, M7c2b, M7c3, M7c4, M8a, M8b. Controller steps: R1, P0, P3. Attended: A1.
- **Next:** M7c2a, then M7c2b.
- **New carried note:** `-race` on `scripts/marketingcheck` reds about 1 run in 100 in a supervisor test (EPERM). The hypothesis is that `kill_group` catches only `ProcessLookupError`. It is M7a code and not fixed yet.
- **Cosmetic residuals (Mark, if wanted):** the phone title is clipped with horizontal overflow, the provenance payload clips horizontally, and the tall image proportions are unusual.
- **Parked on Mark:** D-WORLD-86. A: drop the alias pre-checks and rely on pin+M5e (recommended); B: defence in depth. Default is HOLD.
- **Cost:** metered $0.00; codex and Anthropic quota only.
- **Record/evidence:** world-mission-log.md entry 306; verification/world-iter306/.

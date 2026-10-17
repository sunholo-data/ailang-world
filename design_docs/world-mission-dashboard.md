# World mission dashboard — iteration 278 (2026-10-17)

- **This fire:** row 156 M0b (PR #299). The protected-write matrix ran on Linux CI: V5's 16 arms under no deny, a directory-only deny and a literal-file deny (48 rows). Linux agrees with the design's §4.4 outcome for outcome; only the error codes differ from macOS. Banked as design V14/V15 with both JSON outputs.
- **Finding:** on Linux only directory-form patterns are enforced; base-name matches (`*_test.go`) stay writable, and listing existing files as literals still lets new test files be created (6/6). That confirms the design's detect-only path for Linux.
- **CI:** the first run measured correctly but left fixture trees that three repo-wide Go scans parsed (red); fixed by deleting each tree after its snapshot (gated); re-run green 3/3 with identical rows.
- **Judge:** independent sonnet r1 **PASS 85**, zero blocking; survivors carried to M1a.
- **Latest GitHub release:** `v1.0.0-rc.1` (pre-release, 2026-10-16). Proof/compiler pin AILANG 0.41.0; tool pin 0.52.1.
- **Clauses:** **7 of 7 MET.** Before `v1.0.0` (D-WORLD-81): finish row 156 (2/19), then row 150 (9/17), then row 151.
- **Parked on Mark:** **D-WORLD-82** (row 156 design decisions D-156-1..4; defaults A, A, A, B; the plan follows the defaults).
- **Next:** row 156 M1a.
- **Routing/cadence:** executor codex gpt-6.1-sol (recipe), judge Agent sonnet; no designer/planner (ready-made milestone). 4h throttle.
- **Metered:** $0. No shared skill, harness or routing-policy edit.
- **Evidence:** log entry 278; `design_docs/verification/world-row156-design-2026-10-17/m0b_*.json`; role artifacts in `~/.ailang/state/world-iter278/`.

# World mission dashboard — iteration 307 (2026-10-27)

- **This fire:** row 150. M7c2a, the capture primitives, is BUILT: a Chrome-version floor (≥155), HTTP readiness (`/v1/health` then `/v1/head`), a bounded six-shot `pass` that stops at the first bad shot, and `marketingcheck routes`, whose rows carry the M7c0 heights.
- **Judge:** sonnet PASS 88/100, zero blocking; the plan's 3 named mutants were killed. 13 of the judge's 26 diff-anchored mutants survived (mostly `pass` refusal branches). The one M7c2b depends on (stop at the first non-ok shot) was closed by a controller test.
- **Confirmed defect:** `kill_group` lets macOS `killpg` EPERM escape as `SupervisorError` (a standalone repro hit it 5 times in 600 calls; base `-race` was red 1 time in 4). A real capture on a Mac could fail on a good shot. **M7a-fix goes next**, before M7c2b. Linux CI is unaffected.
- **Roles (Agent tool):** executor codex gpt-6.1-sol, through an Agent supervisor running the pinned recipe; evaluator sonnet. Designer and planner were not needed (the doc and the prescriptive plan exist).
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81) is 156 → 150 → 151.
- **Attended rulings acknowledged:** D-WORLD-95/96/97 (the all-agent-work roadmap), provenance `attended`. They are attended-track rows and leave the release order unchanged.
- **Row 156:** 22/24; D-WORLD-86 holds M5c/M5d.
- **Row 150:** 17/22 built. Pending: M7a-fix (proposed), M7c2b, M7c3, M7c4, M8a, M8b. Controller steps: R1, P0, P3. Attended: A1.
- **Next:** M7a-fix, then M7c2b.
- **Carried for M7c2b:** `pass` creates `--out-dir` before it validates routes; untested refusal branches; `chrome_version` timeout path.
- **Cosmetic residuals (for Mark, if wanted):** the phone title is clipped with horizontal overflow, the provenance payload clips horizontally, and the tall image proportions are unusual.
- **Parked on Mark:** D-WORLD-86. A: drop the alias pre-checks and rely on pin+M5e (recommended); B: defence in depth. Default is HOLD.
- **Cost:** metered $0.00; codex and Anthropic quota only.
- **Record/evidence:** world-mission-log.md entry 307; verification/world-iter307/.

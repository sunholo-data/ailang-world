
# World mission dashboard — iteration 311 (2026-10-28)

- **This fire:** row 150. M7c4 is BUILT (PR #394). The capture shell takes a frozen inventory of the public store before pass A and after the last pass and refuses if anything changed; diagnostics record `inventory: equal`. `TestMarketingChromeSixShots` is the real-Chrome drill for R1.
- **Judge:** sonnet PASS 92/100, zero blocking; the plan mutant was killed at the named assertion.
- **Controller:** every gate re-run outside the executor's sandbox (green, including `-race`). The executor's fake daemon could not tell where an inventory ran relative to the passes, so 3 of 9 drill mutants survived; a test-only fix (`1ca31e5e`) keys the fake head change to a pass, and all 9 are now killed.
- **Real Chrome (156.0.8078.48), smoke only:** the full path runs in 15–20 s. **Made-up 100x100 rectangles PASSED `check` end to end**, so R1's C8 evidence must come from genuinely measured rectangles; a misplaced-rectangle negative control belongs in R1.
- **Roles (Agent tool):** executor codex gpt-6.1-sol through an Agent supervisor running the pinned recipe; evaluator sonnet. Designer and planner not needed.
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81) is 156 → 150 → 151.
- **Row 156:** 22/24; D-WORLD-86 holds M5c/M5d.
- **Row 150:** 21/23 built. Pending: M8a, M8b. Controller steps: R1, P0, P3. Attended: A1.
- **Next:** R1 (measured rectangles, two real passes, threshold V-row), then M8a → M8b.
- **Carried (judge, non-blocking):** the diagnostics line written before the compare is caught only incidentally; the drill's rectangles precondition is weaker than `LoadRectangles`; the manifest assertion only checks it parses; `v0.41.0` hardcoded. Also from 310/309: fetch/unmarshal and CLI usage paths untested; the guard cannot see writes that bypass the log or an A→B→A revert; A1 lifetime and `ready()` stop branch untested.
- **Flake:** `TestMarketingSupervisorFakeDaemonCleanup` red once in a full package run at load ~6 (3/3 green alone); the M7a load-sensitive class.
- **Cosmetic residuals (for Mark, if wanted):** the phone title is clipped with horizontal overflow, the provenance payload clips horizontally, and the tall image proportions are unusual.
- **Parked on Mark:** D-WORLD-86. A: drop the alias pre-checks and rely on pin+M5e (recommended); B: defence in depth. Default is HOLD.
- **Cost:** metered $0.00; codex and Anthropic quota only.
- **Record/evidence:** world-mission-log.md entry 311; verification/world-iter311/.

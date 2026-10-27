# World mission dashboard — iteration 308 (2026-10-27)

- **This fire:** row 150. M7a-fix is BUILT. `kill_group` now swallows macOS `killpg` EPERM as well as ESRCH, so a good shot or a clean daemon stop no longer ends as `SupervisorError [Errno 1]`. Any other `OSError` still propagates.
- **Premise corrected:** the planner's check of `man 2 kill` showed that for a group, EPERM means *any* member could not be signalled. It does not mean all members are dead. The fix is scoped as a workaround for the observed post-reap case, and the existing death assertions stay the guard.
- **Judge:** sonnet PASS 88/100, zero blocking; the plan's 3 named mutants were killed. Judge's controlled comparison: fixed 0/12 vs base 2/12 EPERM-signature reds.
- **Controller:** `-race ./scripts/marketingcheck` 25/25 green on a quiet rig (base 1/3 red). Opt-in stress row: 200 trials, raw EPERM 3–54 per run, 0 wrapper exceptions. A test-only fix raised its wall bound from 25 s to 60 s (measured).
- **New flake class:** M7a's fake-Chrome fixture can red under heavy load (`fake must record child and grandchild: []` at load ~37). It is not EPERM and is carried.
- **Roles (Agent tool):** planner and executor codex gpt-6.1-sol, through Agent supervisors running the pinned recipe; evaluator sonnet. Designer not needed.
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81) is 156 → 150 → 151.
- **Attended ruling acknowledged:** D-WORLD-98 (row 219, home service), provenance `attended`. It is an attended-track row and leaves the release order unchanged.
- **Row 156:** 22/24; D-WORLD-86 holds M5c/M5d.
- **Row 150:** 18/23 built. Pending: M7c2b, M7c3, M7c4, M8a, M8b. Controller steps: R1, P0, P3. Attended: A1.
- **Next:** M7c2b, then M7c3.
- **Carried for M7c2b:** `pass` creates `--out-dir` before it validates routes; untested refusal branches; `chrome_version` timeout path.
- **Cosmetic residuals (for Mark, if wanted):** the phone title is clipped with horizontal overflow, the provenance payload clips horizontally, and the tall image proportions are unusual.
- **Parked on Mark:** D-WORLD-86. A: drop the alias pre-checks and rely on pin+M5e (recommended); B: defence in depth. Default is HOLD.
- **Cost:** metered $0.00; codex and Anthropic quota only.
- **Record/evidence:** world-mission-log.md entry 308; verification/world-iter308/.

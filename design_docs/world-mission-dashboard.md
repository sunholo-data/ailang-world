# World mission dashboard — iteration 318 (2026-10-29)

- **This fire:** dev was red at Gate 1, so the red outranked rule (d). Row 232 LANDED (PR #410, merge `793cc67c`; merge check-runs 3/3 success). Dev is green again.
- **The red:** dev `df5a0671` (#402's merge) failed the go race gate on `reset-holder` (`Why = ""`). It was row 232's flake, not #402, and the rerun was green.
- **The fix:** a reverse-leg dial that fails at connect with ECONNRESET is now SILENT (sticky refusal). ECONNREFUSED, EINVAL and timeouts still record nothing.
- **Linux measured first-party** (throwaway CI probe): with nothing listening, every dial got ECONNREFUSED (7.3M of 7.3M). Against the reset holder the dial fails at connect with ECONNRESET. T7 failed at least 7 of 300 at base and passed 300 of 300 at the fixed head, both under load.
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81): 150 → 151 → v1.0.0.
- **Row 150:** 23/23 milestones built; waits on the **attended A1** (controller P0 ≤24 h before, P3 before pass A).
- **Row 151:** needs row 188 N2d and N4a, after the attended rows 219 → 217 → 218.
- **Rows 231 → 230 → 233:** hygiene, barred by rule (d). New row 233 is R5: a foreign holder that exits while the poll is dialling refuses startup about 3–13% of the time after a port collision.
- **Roles:**
  - designer opus (Agent tool), quorum r1 BLOCKED, then r2 PROCEED at N-1 with gpt6-1-sol absent (unreachable);
  - planner opus, executor haiku;
  - judge sonnet: **PASS 92**, zero blocking.
- **Parked on Mark:** **D-WORLD-102**, now asking about 231 → 230 (+233). Recommend A, default B. **Attended actions:** schedule row 150's A1; the 219 → 217 → 218 track unblocks 151.
- **Next (loop):** a block report, unless D-WORLD-102 = A or dev goes red.
- **Cost:** $0.61 metered (quorum); otherwise Anthropic quota only.
- **Record and evidence:** world-mission-log.md entry 318; verification/world-iter318/.

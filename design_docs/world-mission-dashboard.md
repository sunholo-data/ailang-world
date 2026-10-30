# World mission dashboard — iteration 327 (2026-10-30)

- **This fire:** Gate-1 dev red → row 236 built. `TestEpisodeFinishScopeExpiry` raced the test's clock read against the finish handler's later one across a second edge (dev `4d43d8e3`, race gate attempt 1; attempt 2 green). Fix: an unexported `finishNow` clock seam on `Daemon` + exact-edge tests on a pinned clock. PR #436; judge sonnet PASS 93.
- **Dev:** green at `4d43d8e3` on the rerun (3/3). Row 236 lands when #436's merge commit is 3/3 green (landing note).
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81): 150 → 151 → v1.0.0.
- **Row 150:** 23/23 milestones built. It waits on the **attended A1**: the controller runs P0 ≤24 h before and P3 before pass A.
- **Row 151:** needs row 188 N2d and N4a, after the attended rows 219 → 217 → 218 (all still `[NEXT]`).
- **Rows 231 → 230 → 233:** hygiene, barred by rule (d) until D-WORLD-102 = A. Rows 235 and 237 (new: the approve guard's `Record != nil` arm is unpinned) are hygiene outside D-WORLD-102's scope.
- **Roles this fire:** designer opus (Agent; quorum r1 and r2 BLOCKED on evidence, carve-out applied) · planner opus · executor haiku · judge sonnet. All foreground Agent-tool spawns.
- **Parked on Mark:** **D-WORLD-102**: may the loop take 231 → 230 → 233 while the release rows wait? Recommend A, default B. **Attended actions:** schedule row 150's A1; the 219 → 217 → 218 track unblocks 151. Attended PR #433 (census probe cap) would also remove one of the two load flakes the local full gate hit this fire.
- **Next (loop):** row 236's landing note, then a block report unless D-WORLD-102 = A or dev goes red.
- **Cost:** $0.6052 metered (two quorum rounds); Anthropic quota (opus controller/designer/planner, haiku executor, sonnet judge).
- **Record and evidence:** world-mission-log.md entry 327; verification/world-iter327/; index regenerated with `rotate-log --index-only`.

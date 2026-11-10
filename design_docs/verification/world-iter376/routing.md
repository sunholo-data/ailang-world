# Iteration 376 — routing (resolver outputs, verbatim)

- designer: not spawned. There was no re-design: sprint 2 builds from the quorum-reviewed `w-world-approvals-by-key-a1b.md` (r2) and the sprint-1 plan's Appendix A, as iteration 375's Next named.
- quorum: not run. The doc has two quorum artifacts (`w-world-approvals-by-key-a1b-2026-11-10T05-10-46Z.json` and `…T05-17-00Z.json`) plus the ratified carve-out r2.
- planner: `derive-planner-lane.sh` → `opus fail-closed:env-pin`; env `agent-tool opus`. Spawned as Agent `opus`: 1 run, 31.5 min, 459,045 tok.
- executor: env `recipe claude:claude-haiku-5-5`. Spawned as Agent `haiku`, following the 375 precedent. 4 runs, no git writes:
  - run 1, M3a+M3b: 98,185 tok, 2.0 min.
  - run 2, M3c+M3d: 99,595 tok, 2.8 min.
  - run 3, M3e+M3f: 125,768 tok, 3.7 min.
  - run 4, M3g: 110,032 tok, 3.0 min.
  - Run 5 (M4a/M4b) was not run, because the plan's size gate stopped it at 1,896 > 1,683.
- evaluator: env `agent-tool sonnet`. Spawned as Agent `sonnet` on its own worktree at `3ea0bea9`: 1 round, 15.4 min, 92,745 tok.
- Generator ≠ judge: the executor was haiku, the planner opus and the judge sonnet.
- The driver noted codex unusable at preflight (probe rc 75, quota unreadable).
- Metered spend: $0.00 (no quorum round), against a $5 budget.

# Routing — iteration 381

| Role | Resolver | Spawned as | Duration | Tokens |
|---|---|---|---|---|
| Controller | session | `claude-opus-5-5` | — | not reported |
| Designer | — | not spawned (A2 r2 covers sprint 3) | — | — |
| Quorum | — | not re-run (artifact `w-world-approvals-by-key-a2-2026-11-10T21-14-57Z.json`, doc unchanged) | — | — |
| Planner | `derive-planner-lane.sh` → `opus fail-closed:env-pin`; resolver `agent-tool opus fail-closed:env-pin` | Agent `opus` | 181.2 min | 536,858 |
| Executor | `recipe claude:claude-haiku-5-5 declared:provider-pin` | Agent `haiku`, 4 runs | 8.6 + 5.2 + 2.6 + 13.6 min | 49,537 + 45,854 + 40,458 + 50,472 |
| Evaluator | `agent-tool sonnet declared:alias-pin` | Agent `sonnet`, own worktree | 30.2 min | 105,937 |

Generator ≠ judge: executor haiku, judge sonnet, planner opus. Codex unusable at preflight (probe rc 75, quota unreadable); no role needed it. Metered $0.00.

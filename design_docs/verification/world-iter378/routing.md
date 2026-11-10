# Routing evidence — iteration 378

base=1fce6999ebf875299254a1d77003e9e3351a3347@2026-11-10T16:50:36Z (Gate 1, `mission-base.sh record gate1`)

| Role | Resolver output | Ran as | Tokens | Wall | Notes |
|---|---|---|---|---|---|
| Controller | session | `claude-opus-5-5` | not reported | — | |
| Designer | not spawned | — | — | — | design r2 covers the sprint; no re-design |
| Quorum | not run | — | — | — | `mc_quorum_gate2 w-world-approvals-by-key-a1b` rc 0 (artifact 2026-11-10T05-10-46Z) |
| Planner | `agent-tool opus fail-closed:env-pin`; `derive-planner-lane.sh` → `opus fail-closed:env-pin` | Agent `opus` | 397,156 + 452,035 (S4 revision, SendMessage resume) | 38.0 + 13.5 min | |
| Executor | `recipe claude:claude-haiku-5-5 declared:provider-pin` | Agent `haiku` (372–377 precedent; hook accepted) | 38,466 / 41,522 / 42,852 / 47,709 / 45,582 | 1.7 / 1.2 / 1.3 / 1.1 / 10.8 min | 5 runs: T1+D1, D2+D3, S1+S2, S3, S4 |
| Evaluator | `agent-tool sonnet declared:alias-pin` | Agent `sonnet` | 79,233 | 13.9 min | own worktree at `5e85f8bb`; generator ≠ judge |

Driver preflight note: codex `gpt-6.1-sol` unusable (probe rc 75, quota unreadable) for planner/executor seeds; no role needed it.
Metered total: $0.00 (no quorum, no codex, no metered lane).

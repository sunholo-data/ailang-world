# Iteration 377 — routing evidence

base=545fdaaecf6f23dd58072675964754ad716049f9@2026-11-10T12:49:45Z (Gate 1; Gate 3 re-read identical)

| Role | Env value | Path used | Model | Tokens | Duration | Outcome |
|---|---|---|---|---|---|---|
| Controller | session | — | claude-opus-5-5 | not reported | — | — |
| Designer | claude:claude-opus-5-5 (rotation) | not spawned | — | — | — | design r2 reviewed (2 quorum artifacts); no re-design |
| Quorum | — | not run | — | — | — | doc already reviewed |
| Planner | opus (`derive-planner-lane.sh` → `opus fail-closed:env-pin`) | Agent tool | opus | 526,305 | 50.1 min | plan + JSON, 102 mutants run in scratch, patch protocol |
| Executor | claude:claude-haiku-5-5 (recipe) | Agent tool `haiku` (iteration 372–376 precedent) | haiku | 34,027 (run 1, M4a) + 46,598 (resume: M4b, G1, M4c1) + 144,994 (M4c2–M4e) + 46,724 (M5a–M5c) | 0.6 + 1.7 + 3.7 + 1.7 min | 11 sub-milestones applied, all acceptance rc 0 |
| Evaluator | sonnet | Agent tool, own worktree at 9f014e0d | sonnet | 89,520 | 9.4 min | PASS 92/100, zero blocking |

Driver preflight: codex `gpt-6.1-sol` unusable for planner and executor (probe rc 75, quota UNREADABLE). Metered total: $0.00.

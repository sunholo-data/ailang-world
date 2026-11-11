# Routing evidence — iteration 380

base=3b77c294000ca40c2fe77e312ae0e8ce0f71a030@2026-11-11T00:50:05Z (Gate 1, `mission-base.sh record gate1`); sprint worktree base re-measured at Gate 3: same SHA.

| Role | Resolver output | Ran as | Tokens | Wall | Notes |
|---|---|---|---|---|---|
| Controller | session | `claude-opus-5-5` | not reported | — | |
| Designer | — | not spawned | — | — | the A2 child design r2 covers sprint 2; no re-design (iteration 379's route) |
| Quorum | `mc_quorum_gate2 w-world-approvals-by-key-a2` rc 0 (`…a2-2026-11-10T21-14-57Z.json`) | not re-run | — | — | doc unchanged since its r2 carve-out |
| Planner | `agent-tool opus fail-closed:env-pin`; `derive-planner-lane.sh` → `opus fail-closed:env-pin` | Agent `opus` | 478,824 | 93.8 min | compiled patches in Appendix P; 139 mutants (135 killed, 4 named equivalents) |
| Executor | `recipe claude:claude-haiku-5-5 declared:provider-pin` | Agent `haiku` (372–379 precedent; hook accepted) | 46,214 / 49,619 / 40,775 / 44,161 | 1.5 / 5.0 / 1.3 / 1.4 min | 4 runs of 2 sub-milestones; no hand edits, no git writes |
| Evaluator | `agent-tool sonnet declared:alias-pin` | Agent `sonnet` | 73,582 | 11.1 min | own worktree `.wt-world-iter380-judge` at `1fe7c86d`; generator ≠ judge (haiku/opus vs sonnet) |

Driver preflight note: codex `gpt-6.1-sol` unusable (probe rc 75, quota unreadable) for the planner/executor seeds; no role needed it.
Metered total: $0.00 (no quorum, no codex, no pi). Ceiling $5.

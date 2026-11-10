# Routing evidence — iteration 379

base=98374909268b9eab35382c5bf1779a3784edbf5b@2026-11-10T20:50:02Z (Gate 1, `mission-base.sh record gate1`)

| Role | Resolver output | Ran as | Tokens | Wall | Notes |
|---|---|---|---|---|---|
| Controller | session | `claude-opus-5-5` | not reported | — | |
| Designer | `recipe claude:claude-opus-5-5 declared:provider-pin` | Agent `opus` (318–378 precedent; hook accepted) | 320,153 (r0) + 365,888 (r1, SendMessage resume) | 22.6 + 4.8 min | rotation file reads `claude:claude-opus-5-5` last-used; the resolver's answer was followed (as at 373/375). One design doc: r0 + one protocol-mandated revision (diet held). r2 = controller carve-out |
| Quorum r1 | `ailang design-quorum --author claude:claude-opus-5-5` | seats oc-kimi-k3, gpt6-1-sol, gemini-4-argon; reserve oc-glm-5-3; claude-sonnet-5@claude-p benched (author vendor) | 108,047 in / 7,393 out | — | BLOCKED 3/3, $0.4386, absent_reviewers [] |
| Quorum r2 | same | same seats | 124,095 in / 14,972 out | — | BLOCKED 2/3 (gemini PASS), $0.5997, absent_reviewers [] |
| Planner | `agent-tool opus fail-closed:env-pin`; `derive-planner-lane.sh` → `opus fail-closed:env-pin` | Agent `opus` | 445,890 | 58.3 min | |
| Executor | `recipe claude:claude-haiku-5-5 declared:provider-pin` | Agent `haiku` (372–378 precedent; hook accepted) | 39,424 / 37,756 | 1.7 / 0.7 min | 2 runs: M0+M1, M2+M3a |
| Evaluator | `agent-tool sonnet declared:alias-pin` | Agent `sonnet` | 102,582 + 106,877 (round 2, SendMessage resume) | 6.6 + 0.4 min | own worktree `.wt-world-iter379-judge`; generator ≠ judge |

Driver preflight note: codex `gpt-6.1-sol` unusable (probe rc 75, quota unreadable) for planner/executor seeds; no role needed it.
Metered total: $1.04 (quorum r1 $0.4386 + r2 $0.5997; no codex, no pi). Ceiling $5.

# Iteration 375 — routing (resolver outputs, verbatim)

- designer: `recipe claude:claude-opus-5-5 declared:provider-pin`. The rotation file read `claude:claude-opus-5-5`; the resolver was followed (codex over ration, probe rc 75). Spawned as Agent `opus` (the 318–374 precedent).
  - r0: 18.1 min, 259,187 tok.
  - r1: by `SendMessage` resume, 3.4 min, 298,123 tok as reported.
  - r2: the controller's narrow-refinement carve-out (no designer run).
- quorum: author `claude:claude-opus-5-5`, so Anthropic was benched. Seats gpt6-1-sol, gemini-4-argon, oc-glm-5-3 (reserve oc-kimi-k3). No absent reviewers in either round.
  - r1: BLOCKED 3/3, $0.3757 (`w-world-approvals-by-key-a1b-2026-11-10T05-10-46Z.json`).
  - r2: BLOCKED 2/3, gemini PASS, $0.3611 (`…-2026-11-10T05-17-00Z.json`).
- planner: `agent-tool opus fail-closed:env-pin`; `derive-planner-lane.sh` → `opus fail-closed:env-pin`. Spawned as Agent `opus`: 1 run, 34.6 min, 350,528 tok.
- executor: `recipe claude:claude-haiku-5-5 declared:provider-pin`. Spawned as Agent `haiku` (precedent), 3 runs:
  - M1a+M1b: 104,801 tok, 3.2 min.
  - M1c+M2a: 109,031 tok, 10.6 min.
  - M2b+M2c: 87,422 tok, 3.4 min.
- evaluator: `agent-tool sonnet declared:alias-pin`. Spawned as Agent `sonnet` on its own worktree at `048537c3`: 1 round, 15.7 min, 71,477 tok. Generator ≠ judge: the executor was haiku, the designer and planner opus, the judge sonnet.
- Metered spend: $0.74 (two quorum rounds), against a $5 budget.

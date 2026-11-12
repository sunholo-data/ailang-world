# Iteration 387 routing evidence
base=4134ce898548ec7b9c16c89f952544cdb5ea5d7a@2026-11-12T10:58:15Z
- resolve-role-spawn: planner `agent-tool opus fail-closed:env-pin`; executor `recipe claude:claude-haiku-5-5 declared:provider-pin` (spawned as Agent `haiku`, as 374–386); evaluator `agent-tool sonnet declared:alias-pin`. Designer not resolved or spawned (A3 r2 covers A3a sprint 4).
- derive-planner-lane: `opus fail-closed:env-pin`.
- Controller `claude-opus-5-5` (tok: not reported) · designer not spawned · planner `opus` Agent (395,119 tok, 84.6 min) · executor `haiku` Agent ×2 (45,509 tok, 1.6 min; 45,004 tok, 3.3 min) · evaluator `sonnet` Agent (71,557 tok, 8.4 min).
- Driver routing note: codex `gpt-6.1-sol` probe rc 75 (quota unreadable) → planner opus, executor haiku. No role needed codex.
- Metered: $0.00 (no quorum; the A3 doc's quorum artifacts are from iteration 383). Judge independence: cross-model (opus planner, haiku executor, sonnet judge).

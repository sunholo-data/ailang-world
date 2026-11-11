# Iteration 384 routing evidence
base=0a25c0d6c6c9af6cdf81da75270f41af41a44bad@2026-11-11T20:06:31Z
- resolve-role-spawn: planner `agent-tool opus fail-closed:env-pin`; executor `recipe claude:claude-haiku-5-5 declared:provider-pin` (spawned as Agent `haiku`, as 374–382); evaluator `agent-tool sonnet declared:alias-pin`. Designer not spawned (A3 r2 covers A3a sprint 1).
- derive-planner-lane: `opus fail-closed:env-pin`.
- Controller `claude-opus-5-5` (tok: not reported) · designer not spawned · planner `opus` Agent (452,011 tok, 69 min) · executor `haiku` Agent ×2 (66,873 tok, 7.9 min; 39,719 tok, 5.8 min) · evaluator `sonnet` Agent (77,541 tok, 31 min).
- Driver routing note: codex `gpt-6.1-sol` probe rc 75 (quota unreadable) → planner opus, executor haiku. No role needed codex.
- Metered: $0.00 (no quorum; the A3 doc's quorum artifacts are from iteration 383). Judge independence: cross-model (opus planner, haiku executor, sonnet judge).

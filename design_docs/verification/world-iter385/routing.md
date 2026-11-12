# Iteration 385 routing evidence
base=a8c666554f284651de33d2f6a2c85a16f5b72c02@2026-11-12T00:06:27Z
- resolve-role-spawn: designer `recipe claude:claude-opus-5-5 declared:provider-pin` (not spawned: A3 r2 covers A3a sprint 2); planner `agent-tool opus fail-closed:env-pin`; executor `recipe claude:claude-haiku-5-5 declared:provider-pin` (spawned as Agent `haiku`, as 374–384); evaluator `agent-tool sonnet declared:alias-pin`.
- derive-planner-lane: `opus fail-closed:env-pin`.
- Controller `claude-opus-5-5` (tok: not reported) · designer not spawned · planner `opus` Agent (466,125 tok, 77 min) · executor `haiku` Agent ×2 (52,970 tok, 5.9 min; 45,017 tok, 7.3 min) · evaluator `sonnet` Agent (84,797 tok, 8.6 min).
- Driver routing note: codex `gpt-6.1-sol` probe rc 75 (quota unreadable) → planner opus, executor haiku. No role needed codex.
- Metered: $0.00 (no quorum; the A3 doc's quorum artifacts are from iteration 383). Judge independence: cross-model (opus planner, haiku executor, sonnet judge).

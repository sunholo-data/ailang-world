# Iteration 386 routing evidence
base=82fb9afa9086eed2f781bb0d4fe54d4bae3b70d7@2026-11-12T04:07:19Z
- resolve-role-spawn: planner `agent-tool opus fail-closed:env-pin`; executor `recipe claude:claude-haiku-5-5 declared:provider-pin` (spawned as Agent `haiku`, as 374–385); evaluator `agent-tool sonnet declared:alias-pin`. Designer not resolved or spawned (A3 r2 covers A3a sprint 3).
- derive-planner-lane: `opus fail-closed:env-pin`.
- Controller `claude-opus-5-5` (tok: not reported) · designer not spawned · planner `opus` Agent (496,816 tok, 97.3 min) · executor `haiku` Agent ×2 (50,422 tok, 10.4 min; 69,352 tok, 4.8 min) · evaluator `sonnet` Agent (114,178 tok, 24.1 min).
- Driver routing note: codex `gpt-6.1-sol` probe rc 75 (quota unreadable) → planner opus, executor haiku. No role needed codex.
- Metered: $0.00 (no quorum; the A3 doc's quorum artifacts are from iteration 383). Judge independence: cross-model (opus planner, haiku executor, sonnet judge).

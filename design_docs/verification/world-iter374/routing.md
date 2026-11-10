# Iteration 374 — routing (resolver outputs, verbatim)

- designer: `recipe claude:claude-opus-5-5 declared:provider-pin` — NOT SPAWNED: the r3 child design already covers M4–M6 ("built from this same doc without re-design", design §0), and the quorum artifacts exist (`mc_quorum_gate2` rc 0, r1 2026-11-09T21-08-17Z, r2 2026-11-09T21-20-35Z).
- planner: `agent-tool opus fail-closed:env-pin`; `derive-planner-lane.sh` → `opus fail-closed:env-pin`. Spawned Agent `opus`, 1 run, 28.9 min, 426,444 tok.
- executor: `recipe claude:claude-haiku-5-5 declared:provider-pin`. Spawned Agent `haiku` (the 318–373 precedent), 5 runs: M3b+M4a 131,590 tok / 5.2 min; M4b+M4c 109,383 / 7.1; M5a+M5b 155,339 / 13.3; M6a+M6b 178,683 / 10.2; M6c+M6d 156,258 / 31.2.
- evaluator: `agent-tool sonnet declared:alias-pin`. Spawned Agent `sonnet` on an isolated detached worktree at `86aa2f9c`, 1 round, 21.7 min, 98,492 tok. Generator ≠ judge: executor haiku, planner opus, judge sonnet.
- Metered spend: $0.00 (no quorum round, no pi/codex lane).

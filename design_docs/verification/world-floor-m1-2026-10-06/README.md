# Row 93 AC1.2 re-run on tool binary v0.52.1 (2026-10-06)

This is a reproduction of row 135 M0's AC1.2 re-run (`../world-row135-m0-2026-10-03/`). It was made when the row-93 M1+M3 landing was re-checked against `dev` after row 140. The floor grader's `$TOOL` stays at v0.52.1: `~/.pinned-ailang-tools/v0.52.1/ailang`, sha256 `0dd70a1d00360be0b7b4c667054c71fdc4ea87aa3deb79710cf56f18ffe1a8f5`, `AILANG v0.52.1` commit `c68ded4`. That is the same binary the World arm executes.

Command (attended executor, darwin arm64):

```
python3 scripts/floor/agreement.py --repo <ailang> --commit 76a5aef65c81ac2fb4657fafaeff4a89e9faaf75 \
  --tool-bin ~/.pinned-ailang-tools/v0.52.1/ailang \
  --alt-bin ~/.pinned-ailang-tools/v0.51.0/ailang --alt-bin ~/.pinned-ailang/ailang \
  --alt-bin ~/.pinned-ailang/ailang.v0.30.0 \
  --expect-tool-sha256 0dd70a1d00360be0b7b4c667054c71fdc4ea87aa3deb79710cf56f18ffe1a8f5 \
  --out grader-agreement.json
```

The run exited rc 0 and printed `AC1.2 PASS`.

| Run | Tool | Rows | Agree | Disagreements (typed) | Unexplained |
|---|---|---|---|---|---|
| `world-floor-m1-2026-10-03` | v0.51.0 | 1357 | 1351 (0.9956) | 6: 3 version, 2 stdlib-skew, 1 grade-timeout | 0 |
| `world-row135-m0-2026-10-03` | v0.52.1 | 1357 | 1350 (0.9948) | 7: 4 version, 2 stdlib-skew, 1 grade-timeout | 0 |
| **this run** | v0.52.1 | **1419** | **1410 (0.9937)** | 9: 4 version, 2 stdlib-skew, 3 grade-timeout | **0** |

Every grade executed exactly the pinned digest `0dd70a1d…`.

**The same 1357 rows reproduce exactly.** This run has the same 1357 row paths as the 2026-10-03 v0.52.1 run, and its regraded `stdout_ok` is identical on every one of them (1350 agree in both).

**62 rows were banked since.** The ailang checkout's `eval_results/` gained 62 eligible core rows after 2026-10-03, all under `rotation/os-rolling/v0.52.0`. 60 of them agree. The other 2 are `banked_grade_timeout`: banked fails that the port now grades as passes. They are:
- `effect_tracking_io_fs` by `motoko-local-qwen3-8-27b-microrag`;
- `csv_to_json_converter` by `opencode-qwen3-8-27b`.

**Corpus pin unchanged.** ailang `origin/dev` `b7028a7c…` (2026-10-06, 47 commits past `76a5aef6…`) changes nothing under `benchmarks/`, the agent task template, or the grader (`runner.go`, `agent_validation.go`, `agent_runner_multi.go`, `agent_prompt.go`, `spec.go`). The only change under `internal/eval_harness/` is one new test file, `agent_runner_multi_cost_test.go`.

`grader-agreement.rows.jsonl` holds the per-row records.

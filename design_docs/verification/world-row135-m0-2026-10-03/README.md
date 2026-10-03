# Row 135 M0: floor grader agreement on tool binary v0.52.1 (2026-10-03)

The tool-binary pin moved from v0.51.0 to v0.52.1 (row 135 M0). The World arm and the floor grader use the same `$TOOL` (row-93 design P12), so the row-93 grader-agreement check AC1.2 was re-run on the new binary before `scripts/floor/agreement.py`'s pinned digest was changed.

Command (attended executor, darwin arm64):

```
python3 scripts/floor/agreement.py --repo <ailang> --commit 76a5aef65c81ac2fb4657fafaeff4a89e9faaf75 \
  --tool-bin ~/.pinned-ailang-tools/v0.52.1/ailang \
  --alt-bin ~/.pinned-ailang-tools/v0.51.0/ailang --alt-bin ~/.pinned-ailang/ailang \
  --alt-bin ~/.pinned-ailang/ailang.v0.30.0 \
  --expect-tool-sha256 0dd70a1d00360be0b7b4c667054c71fdc4ea87aa3deb79710cf56f18ffe1a8f5 \
  --out grader-agreement.json
```

The run exited with rc 0 and printed `AC1.2 PASS`.

| Measure | Result |
|---|---|
| Rows | 1357 regraded (1272 banked passes, 85 banked fails) |
| Agreement | 1350 agree (**0.9948**) |
| Disagreements | 7 |
| Executed tool digest | exactly `0dd70a1d0036…` |

The 7 disagreements break down by cause:
- `version`: 4
- `banked_stdlib_skew_v37`: 2
- `banked_grade_timeout`: 1
- `unexplained`: **0**

**The baseline.** The banked v0.51.0 run is in `../world-floor-m1-2026-10-03/`: 1351/1357 agree (0.9956), with 6 disagreements of which 3 are `version`.

**The new disagreement.** It is an `api_call_json` solution by `opencode-qwen3-6-35b-a3b-mxfp8`, banked on v0.33.1:
- v0.52.1 refuses it at type-check: `TC_MATCH_001: constructor pattern 'Network' does not name any constructor in scope`.
- v0.51.0 reproduces the banked pass, so the cause is `version`.

**The four row-135 tasks.** Agreement per task, on v0.52.1 and on v0.51.0:

| Task | v0.52.1 | v0.51.0 |
|---|---|---|
| `pipeline` | 62/62 | 62/62 |
| `cli_args` | 83/83 | 83/83 |
| `prompt_injection` | 18/19 | 18/19 |
| `api_call_json` | 81/82 | 82/82 (the one row lost is the new disagreement above) |

`grader-agreement.rows.jsonl` holds the per-row records.

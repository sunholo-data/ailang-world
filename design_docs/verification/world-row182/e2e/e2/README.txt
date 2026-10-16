Ground truth (E2) for one question about a small benchmark run (row-182 end-to-end check, 2026-10-16), produced by the controller; the answering walk neither produced nor saw it.
- grader-rows.json: the graded rows of the run (codex; shell and World arms; tasks fizzbuzz canary and config_file_parser): class, error category, grader stdout/stderr heads, solution_sha256.
- grader-clean-dir-rule.txt: the grader source: how a solution is placed and run.
- rerun-world.json: the World arm's final benchmark/solution.ail (from the store's Workspace.Write requests; sha256 equals the graded one) re-executed with the grader's argv in a clean directory and in a directory that also holds the app_config.json the agent wrote.
- shell/ + rerun-shell.txt: the shell arm's final solution.ail (sha256 equals the graded one) and the app_config.json it left in its worktree, re-executed the same two ways.

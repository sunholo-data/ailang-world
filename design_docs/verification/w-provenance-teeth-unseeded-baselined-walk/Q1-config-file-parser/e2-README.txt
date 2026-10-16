Ground-truth bundle (E2) for one question about the 2026-10-14 row-93 FINAL. Produced by the controller; neither answering arm produced or saw it.
- phase2-config-rows.txt: the FINAL's phase-2 log lines for task config_file_parser (both agents, both arms, runs 1-3).
- grader-rows.json: the graded rows for that task: class, error_category, grader stdout/stderr heads, solution_sha256.
- grader-clean-dir-rule.txt: the grader source (scripts/floor/grade.py at the pre-diagnosis commit): how a solution is placed and run.
- rerun.json + rerun/: for each codex World run, the final benchmark/solution.ail and the app_config.json written in that run (taken from the World store's Workspace.Write effect requests), re-executed with the grader's argv on the grading release (AILANG v0.41.0) twice: in a clean directory ("clean") and in a directory that also holds that run's app_config.json ("with_app_config"). solution_sha256 there equals grader-rows.json's.

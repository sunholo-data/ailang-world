Ground truth (E2) for one question about trial episode tps-b (platform-source, 2026-10-16), produced by the controller by re-execution after the episode; neither answering arm produced or saw it.
- grader-212-episode.stdout: the operator's ADK 2.12 grader output (`.venv-adk212/bin/python -m pytest -q -p no:cacheprovider -m adk_contract tests/unit`, cwd backend) on the episode's final worktree; exit rc=2.
- base-212-suite.stdout: the same command, same 2.12 venv, on the repository at the episode's starting commit e5b0461e (before the episode added anything); rc at the end.
- collect-errors-episode.txt / collect-errors-base.txt: the `ERROR tests/...` collection-error lines of each run, sorted. They are byte-identical (15 each) and none names the four new test files.
- new-files-212-individually.txt: each of the four new test files run ALONE on the same 2.12 venv on the final worktree (tail of output + rc). Each collects and fails 1 test (rc=1), i.e. at runtime, not collection.
Note: the 2.12 venv's editable install points at the episode worktree's backend package; the episode changed no production code, only tests and one doc.

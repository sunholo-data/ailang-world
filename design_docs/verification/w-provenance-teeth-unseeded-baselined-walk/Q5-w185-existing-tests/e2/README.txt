Ground truth (E2) for one question about trial episode w185 (ailang-world, 2026-10-16), produced by the controller from git and re-execution after the episode; neither answering arm produced or saw it. Base = a08ca14 (the episode's starting commit); final = the episode's committed worktree.
- all-changed-files.txt: `git diff --name-status a08ca14 HEAD` (M = modified, A = added).
- modified-test-files.diff: the full diff of every pre-existing _test.go the change touched (session_new_test.go, worktree_test.go). It contains 0 removed lines; every hunk inserts new code after the closing brace of an existing function.
- test-func-counts.txt: number of top-level `func Test…` in each of those files, base vs final.
- base-tests-referencing-changed-behaviour.txt: lines in the base commit's _test.go files that mention the refusal text ("is not a git worktree", "pick another episode") or isGitWorktree; isGitWorktree still exists in the final code (it now delegates to host/worktree.IsLinked).
- go-test-final.txt: `go test -count=1 ./cmd/ailang-worldd ./host/worktree ./host/daemon` on the final worktree, unsandboxed (rc at the end).

Ground truth (E2) for one question about trial episode tps-a (platform-source, 2026-10-16), produced by the controller; neither answering arm produced or saw it.
- git-status.txt: `git status --porcelain --untracked-files=all` of the episode worktree after the episode (and after the graders).
- untracked-files.txt: every untracked path with mtime|size|path.
- store-invocations.json: every log entry of the episode's World store (archival copy, sha256 fddd2e59…): entry index, writtenBy, episode, skill (World tool) and its input.
- transcript-tool-calls.json: every tool call in the agent's own stream-json transcript, in order.
- timeline.txt: the agent's start/finish and the graders' run times.

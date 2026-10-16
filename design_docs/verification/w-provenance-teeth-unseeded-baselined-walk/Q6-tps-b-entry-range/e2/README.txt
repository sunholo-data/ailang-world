Ground truth (E2) for one question about trial episode tps-b (2026-10-16), produced by the controller after the episodes; neither answering arm produced or saw it.
Context: two trial episodes, tps-b (codex, platform-source) and w185 (Claude Code, ailang-world), were launched by the operator at the same moment against ONE World daemon and store.
- episode-records.json: the runner's own record of each episode: start/finish time, log_range (the runner reads the store head before spawning the agent → from = head+1, and after it exits → to = head), and the World-Provenance trailer (from `provenance --since <from> --episode <ep>`).
- trial-py-excerpt.txt: the runner code that computes log_range, the trailer, and the PR body's "World entries" line.
- store-entries-82-100.json: from the archival copy of the store (sha256 9d72ad32…), each log entry 82–100 with its episode, the World tool (skill) and the tool input. Entries 82–88 were the operator's warm-up calls made before launch.

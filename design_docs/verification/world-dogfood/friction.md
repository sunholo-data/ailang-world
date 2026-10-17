# World-on-World dogfood — friction log (row 187)

Each entry: what happened, its cost, and the row that owns the fix. Newest last.

| # | date | episode | friction | cost | owner |
|---|---|---|---|---|---|
| 1 | 2026-10-16 | every | Each task needs a human-typed `session new` at a real TTY. That is correct for authority, but it is a context switch to a terminal per task. | ~1 min of Mark's attention per task | row 189 (`episode start`), row 188 (desktop approval) |
| 2 | 2026-10-16 | tps-c | Minting needs the daemon STOPPED (`session new` takes the store writer lock), so a new task cannot be minted while another episode runs. | serializes tasks; daemon restart per mint | row 189 (mint through or beside a serving daemon) |
| 3 | 2026-10-16 | w185 | The exec sandbox denies `/dev/ptmx`, so `TestSessionNewReadsTheAnswerFromThePtyWhenStdinIsDevNull` fails sandboxed. The operator added a line to the task text. | one operator note per task touching session_new | row 187 (a profile-level known-failure list) |
| 4 | 2026-10-16 | wdocs | The exec sandbox denies binding a loopback port, so `httptest` tests (e.g. `TestQuickstartSection9FlagsMatchTheCLI`) cannot run in an episode. They passed unsandboxed afterwards. | the agent cannot self-check those tests | row 187 (opt-in loopback for test profiles) |
| 5 | 2026-10-16 | wdocs | To read large files the agent wrote a throwaway `zz_peek.ail` line-reader and ran it, because `ailang-read` returns a whole file. | an extra program and compile cache per large read | candidate: ranged reads (`offset`/`limit`) in `ailang-read` |
| 6 | 2026-10-16 | tps-b + w185 | Two episodes on one daemon interleave entries, so the PR body's unfiltered "World entries N–M" range includes the other episode's entries (row 114 Q6). | misleading PR text | `scripts/trial/trial.py` pr-body (use the episode-filtered range) |
| 7 | 2026-10-16 | — | `ailang-worldd` refuses every verb (even read-only `why`) while `AILANG_REGISTRY_API_KEY` is in the shell. This is by design, but the operator's shell carries the key. | every command needs `env -u` | operator environment (keep the key out of the login shell) |
| 8 | 2026-10-17 | tps-c | World has no delete tool, and `workspace-exec` takes no per-call env vars, so a one-off generator step (regenerating a golden behind an env flag) had to be done by temporarily rewriting the test file. | an extra write/run/restore cycle; a delete is impossible | candidate: `Workspace.Delete` / `ailang-rm`, and profile-declared env flags |

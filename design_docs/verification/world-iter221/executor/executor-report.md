# Iteration 221 executor acceptance report

PASS: exact unchanged base and candidate acceptance profiles completed sequentially; all four rc0.

Actual executor: codex:gpt-6.1-sol (native role declaration). Role tokens: not reported.

No source/test/fixture edits, commits, pushes, external messages or merges. Fresh detached worktrees retain exact HEAD and empty porcelain. Main untracked folder and attended worktree untouched.

| Command | Start UTC | End UTC | rc |
|---|---|---|---|
| base ./scripts/verify_ail.sh | 2026-10-02T11:28:58.172707+00:00 | 2026-10-02T11:29:06.038476+00:00 | 0 |
| base ./scripts/verify_go.sh | 2026-10-02T11:29:06.038933+00:00 | 2026-10-02T11:32:48.477272+00:00 | 0 |
| verify ./scripts/verify_ail.sh | 2026-10-02T11:32:48.478144+00:00 | 2026-10-02T11:32:54.215519+00:00 | 0 |
| verify ./scripts/verify_go.sh | 2026-10-02T11:32:54.215853+00:00 | 2026-10-02T11:37:04.331852+00:00 | 0 |

All named AIL stages present: 16 identities, 40 tests, nine package steps. Both Go profiles contain hygiene, ledger, floor, armed positive race control, build, exact evidence37, and matching full plain/race passing package sets. No deadline reached.

Binary/version/sha256, Go/platform, PID/PGID/deadline and exact heads are preserved in preflight.json and four command JSON files. Full stage assertion details are in executor-report.json.

Raw directory: /tmp/world-iter221-executor-heg34rbs

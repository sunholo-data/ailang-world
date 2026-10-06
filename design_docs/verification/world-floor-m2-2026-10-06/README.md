# Row 93 M2 canaries (2026-10-06)

These canaries were run by the attended executor on darwin arm64 with `claude` 2.1.291 and `codex-cli` 0.159.2, using the tuning models only: `claude-haiku-4-5` and `gpt-6-luna`. Five model calls were spent. Design: `design_docs/planned/w-resident-agent-non-inferiority-floor-run.md` §6 M2, V-M2-1..5.

| Dir | What | Result |
|---|---|---|
| `canary-shell-r1/` | `canary.py shell`, both agents, on the **drafted** argv (Claude `--permission-mode bypassPermissions`) | codex **PASS**. Claude **FAIL** on AC2.5: the Write tool wrote outside the worktree, while Bash was refused by the sandbox |
| `canary-shell-r2/` | `canary.py shell --agents claude` on the **amended** argv (`acceptEdits`, V-M2-3) | Claude **PASS**: AC2.2 shell half, AC2.5, AC2.6 and the inside-write positive control |
| `probes/` | the `acceptEdits` probe (Claude), the forced outside-`echo` probe (codex), and the offline `web_search` key checks | V-M2-1, V-M2-3, V-M2-4 |
| `canary-mcp-list/` | `canary.py mcp-list`: `codex mcp list` under each arm's overrides (no model call) | World arm `["world"]`, shell arm `[]` |
| `canary-world/` | `canary.py world`, run after the attended mint (§4.5) | **pending**: needs sessions `fl-tune-cc` and `fl-tune-cx` |

Each `summary.json` records the redacted argv (the prompt appears only as a digest), the init or item counts, the file checks, and its own AC2.7 token grep. The prompts appear in plain text in `summary.json` and in the gzipped transcripts. These evidence dirs hold no session token: the shell arms have none, and `canary-world/` is grepped against both session files before it is written.

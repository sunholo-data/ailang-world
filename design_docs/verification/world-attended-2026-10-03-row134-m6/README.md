# Row 134 M6: attended smoke (2026-10-03)

Mark ran the two terminal-fenced steps (`world-publish transitions`, `session mint`) with stdin from `/dev/tty`. The controller ran the rest, against a scratch store, with the row-134 branch binaries, the PIN at v0.41.0 and the TOOL at v0.51.0.

| Step | Result |
|---|---|
| publish | `world/transition-registry/v1` head present (8 transitions) |
| mint | episode `ep1`, 6 grants × 50, TTL 14400 s |
| `tools/list` (curl) | exactly the 8 tools (`tools-list.json`) |
| `ailang-read` (curl) | `ok:true`, content, `world.effects[0]` ok with a record ref (`call-read.json`) |
| **pi** 0.85.1, `--no-builtin-tools`, World tools only | read → edit → check (0 errors) → run; stdout `hello from pi`; 9 s (`pi.txt`) |
| **Claude Code**, `--tools "" --strict-mcp-config`, World tools only | read → edit → check → run; stdout `hello from claude`; 9 s (`claude.txt`) |
| World log | 10 entries: genesis + 9 tool calls (1 curl + 4 pi + 4 Claude); 9 effect intents / 9 outcomes (`log-*.json`) |

Snags found on the first run, now handled by `tools/attended/se_smoke.sh`:
- the API-key guard on daemon start;
- the TTY fence in embedded terminals;
- `/v1/commit` being session-gated, so the genesis commit moves after the mint;
- SSE responses on `/mcp/`;
- stopping the daemon by PID;
- pi's adapter cache warm-up;
- pi blocking on a non-TTY stdin.

No session token appears in this directory (grep-checked).

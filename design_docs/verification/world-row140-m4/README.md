# Row 140 M4 — attended run on three real projects (2026-10-06)

Attended by Mark Edmondson. Mark ran the two TTY-fenced steps himself: the publish (9 descriptors, revision 1, head `sha256:973c777e…`) and minting three sessions (`go1`, `ts1`, `py1`, `--preset se-tools`). Claude ran everything else from `docs/QUICKSTART.md` §9–§10.

**Rig:**
- macOS 26.6.2 arm64.
- `ailang-worldd` built from `dev` `7664d3f`, pin AILANG v0.41.0, tool binary v0.52.1.
- srt 0.0.78 (`cli.js` sha256 `3c3092bd…d96e`), node v26.0.0.
- Store `/tmp/m4-world`, workspace root `/tmp/m4-ws`.
- Each episode is a detached worktree of the real checkout:
  - `go1`: `~/dev/sunholo-data/ailang`;
  - `ts1`: `~/dev/markedmondson1234/TwilightGame`, with `npm ci`;
  - `py1`: `~/dev/sunholo-data/sunholo-platform`, with `uv sync --frozen`.
- The profiles as served are in this directory, with `$HOME` redacted. `serve.log` is the daemon log; `log-tail.txt` is the store's log (11 entries).

## Results

**Clients:**
- The go1 calls ran through **Claude Code** as the MCP client (`claude -p --tools "" --strict-mcp-config`, only `mcp__world__workspace-exec` allowed).
- The ts1 retry and the py1 calls were direct `/mcp/` JSON-RPC posts with each episode's bearer session.

| Episode | Call | Result |
|---|---|---|
| go1 | `test ./internal/lexer/` (Claude Code), 1st | `timed_out: true`, 9011 ms, 0 bytes: cold first touch of a fresh worktree (finding F4) |
| go1 | the same, 2nd | **exit 0, 873 ms**, `ok github.com/sunholo-data/ailang/internal/lexer 0.251s` |
| ts1 | `typecheck` (Claude Code), 1st | MCP `host callback timed out`; daemon: `execute plan phase: context deadline exceeded` (finding F3) |
| ts1 | `typecheck`, retry | **exit 0, 6797 ms**, `tsc --noEmit` on the whole project |
| py1 | `test-file -k bq tests/test_cli_bq.py`, 1st | exit 126, `.venv/bin/python: Operation not permitted` (finding F1: the profile's `read_roots` missed uv's version-alias dir) |
| py1 | the same, after the profile fix | **exit 0, 523 ms**; argv `… pytest -k bq tests/test_cli_bq.py`: `-k` is emitted as two items (`form: "sep"`) |
| py1 | `install requests` (unprofiled) | refused: `command "install" is not a command of exec profile "sunholo-cli"`; nothing ran |
| py1 | `test-file -p evil …` (unlisted flag) | refused: `argument 0 "-p" refused: the command does not list this flag`; nothing ran |
| py1 | `test-file tests/test_m4_decoy.py` | **exit 0, test passed**. The test fails if it can read `$HOME/.m4-decoy`. Control: the same test run unsandboxed **FAILED** (`AssertionError: DECOY-READ…`). Afterwards the decoy on the host was unchanged |

The provenance trailer reads `World-Provenance: store=/tmp/m4-world/world.db episode=py1 entries=4-10`.

## Findings

- **F1 (doc bug, fixed in this commit):** QUICKSTART §10's Python profile adds only the *realpath* of uv's interpreter dir to `read_roots`. The venv's `python` links to uv's version-alias dir (`cpython-3.12-macos-aarch64-none` → `cpython-3.12.13-…`), so the alias path stayed fenced and exec failed (126). Design V24 measured that both the link and its realpath are needed. The script now adds both.
- **F2 (doc, fixed):**
  - §10 guessed the TwilightGame and sunholo-platform paths wrong. They are now variables the reader must set.
  - §9 told the operator to type `publish 8 transitions` (it is 9).
  - The two TTY-fenced commands need `< /dev/tty` in an IDE terminal. The fence's STOP line names the reason but not this fix.
- **F3 (product, filed as a row):** the plan phase runs `exec.ail` in the pinned interpreter. Under load on this rig it overran its deadline and the MCP client saw a bare `host callback timed out`. This is the same signature as CI's one-off `TestExecSrtMCPEndToEnd` failure. A retry passed.
- **F4 (operational, documented):** the first exec in a freshly created worktree can overrun the 9 s cap from cold disk alone. `go test` measured 14 s → 2.3 s → 0.3 s unsandboxed with 0.4 s CPU, while Spotlight (`mds`) and `syspolicyd` were indexing the new trees. This is R-140-9 in practice: the result is a typed `timed_out`, and a second call passes. Operators should warm each worktree once (run the command unsandboxed) before the session.
- **F5 (product, filed as a row):** for `go1`, building the episode's workspace handlers first ran the AILANG policy summary, which timed out at 10 s (`broker: policy summary: … handler subprocess timed out after 10s`). That left `Workspace.Exec` with no handler for the call (`declares effects with no registered handler`). Exec does not use the AILANG policy, so it should not be unavailable when that step fails.

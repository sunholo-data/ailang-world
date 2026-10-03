# First dogfood: Daneel #201 through AILANG World (2026-10-03, attended)

**Task:** Daneel issue #201. Replace hand-rolled string `length`/`find` helpers with the `std/string` natives across Daneel's `tools/*.ail` modules.

**Setup:**
- **World:** store `~/.ailang/world-dogfood/world/` (kept, because it is real clause-5 corpus). Tool binary v0.51.0, interpreter v0.41.0.
- **Workspace:** the workspace root is a git worktree of `sunholo-data/daneel`, and the episode is **`tools`**, so the sandbox is Daneel's module root (the workaround behind row 141).
- **Attended steps:** Mark ran the transitions publish and the session mint (300/300/300/100/100/100 calls, 4 h).
- **Agent:** Claude Code with `--tools "" --strict-mcp-config`, using World's MCP tools only (prompt in `prompt.md`).

**Result:**
- **9 modules changed, net −106 lines.** Every module still passes `ai_check`, with verify counts unchanged.
- **Behaviour confirmed unchanged:** the natives were probed directly, and a 14-case equivalence script over the edited exports passed, with a planted failure caught.
- **One module skipped:** `daneel_sources` imports a registry package the confined tools can't resolve (row 141).
- **Daneel's own guard** (`tools/ci.sh` with the CI-pinned AILANG v0.50.0) reports `guard: ok`: 52 + 23 modules verified, 24 test suites passed (`daneel-guard.txt`).
- **PR:** sunholo-data/daneel#325 — **merged by Mark (2026-10-03T10:35:12Z 687fac6d)**, labelled `ailang-world` (D-WORLD-60).

**What World recorded:** 117 log entries, one per call (`log-summary.json`):

| Calls | Kind |
|---|---|
| 40 | edit |
| 35 | run |
| 23 | cli |
| 11 | check |
| 4 | write |
| 2 | search |
| 1 | read |
| 1 | genesis |

Each call carries its plan and effect record.

**World gaps found (filed):**
- **Row 141:** a module-root subdirectory sandbox, and registry packages under confinement.
- **To row 138/140:** there is no delete tool (the agent's scratch files had to be removed by hand), and no ranged read (it wrote a grep helper because a file was too large to read whole).

**Environment note:** the `ailang` on PATH is a dirty dev build (`v0.52.0-17-…-dirty`). It rejects Daneel's unchanged `ci_gate.ail`, so the guard was run with the pinned v0.50.0 release only.

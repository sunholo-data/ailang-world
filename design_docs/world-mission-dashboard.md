# World mission dashboard — iteration 432 (2026-11-25)

- **Done this fire:** row 243 **V4 sprint 5b LANDED** (PR #659, merge `5d3b0581`, merge CI 3/3; judge **PASS 84**, 0 blocking).
  - **Last lock gap closed:** the terminal `session mint` now takes the per-episode lock and refuses an operator-owned episode by exact action id, consuming nothing.
  - **Atomic operator mint:** a sealed proof from the runner recheck; the episode start, session credential and mint receipt are written in one transaction; a lost handoff never mints twice. It has **no production caller** yet.
  - **Attended resolver:** `ailang-worldd service runner reconcile-mint ACTION-ID` keeps a matching token, or revokes only an orphaned credential after you type `y`. Documented in QUICKSTART.
  - 1,400 changed lines (exactly the cap). Nothing dispatches yet.
- **Split:** worktree reservation, task/token/MCP file publication and the paired-process consumer moved to **S5c** (measured floor 378 lines).
- **Before your next `service upgrade` of the home service:** still back up `credentials.db`.
- **DECISIONS FOR MARK:** none open (ledger 107 rows, 0 OPEN).
- **Next loop work:** **V4 S5c**. Carried: seven untested defence-in-depth guards the judge found, owed before the mint gets a caller (one reproduced: the transaction's `start`-kind check).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** controller `claude:claude-opus-5-5`; planner and executor `codex:gpt-6.2-sol` (479k / 60k tok); judge Agent `sonnet`, cross-vendor. Metered $0.
- **Dev CI:** 3/3 green at `f8d8e3ad` at start. PR #659 3/3 on attempt 1 (row 267 did not hit; count stays 5).
- **Harness:** row 268 records fleet ticket `skill:mission-base-relative-path-records-driver-repo` (blocking none).
- **Open PRs:** #617 attended UI work, #650 row 228 search (neither touched by the loop).
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).

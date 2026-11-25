# World mission dashboard — iteration 433 (2026-11-25)

- **Done this fire:** row 243 **V4 sprint 5c LANDED** (PR #661, merge `3a6e98a7`, merge CI 3/3 on attempt 2; judge **PASS 87**, 0 blocking).
  - **Guard tests:** five of the seven untested mint guards the last judge found now have killing tests, and the two legacy preset parity gaps are closed. Tests only.
  - **Start check before any effect:** a signed headless start is now re-read against its live inputs (the store, the episode start record with its preset and task, profile, project, grants, catalogue, token path, served build and interpreter, worktree path, launch binary and wall clock, MCP path), and any changed fact is refused by name before anything is locked, written or launched. It has **no production caller** yet.
  - 401 changed lines. Nothing dispatches yet.
- **Split:** worktree reservation (with its attended disposition), file publication and the paired-process consumer moved to **S5d / S5e / S5f** (316 measured probe lines; full sizes unmeasured).
- **Before your next `service upgrade` of the home service:** still back up `credentials.db`.
- **DECISIONS FOR MARK:** none open (ledger 107 rows, 0 OPEN).
- **Next loop work:** **V4 S5d**. Carried: a test that the workspace/hosts directories reach the start check (the judge's top finding), and two guard tests that need two processes, both owed before the mint gets a caller.
- **New queue row (attended, during this fire):** row 269 `w-ui-profile-preset-onboarding` (design r0, #662) — Phase 1 read-only view can land independently; Phase 2 writes wait on row 243 V4. Not yet ranked against row 243; the next fire's Gate 2 reads it.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** codex timed out at the driver probe, so planner Agent `opus` (272k tok), executor `claude-haiku-5-5` (subscription), judge Agent `sonnet` (99k tok) — all Anthropic, independence model-level only (flagged). Metered $0.
- **Dev CI:** 3/3 green at `606f6e2c` at start. PR #661 3/3 on attempt 1. The merge's race job timed out in `host/daemon` and `host/verifygate` (row 267, instance 6; neither touched); attempt 2 green.
- **Open PRs:** #617 attended UI work (not touched by the loop). #650 (row 228 search, attended) merged by its owner during this fire.
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).

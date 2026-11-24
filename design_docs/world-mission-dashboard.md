# World mission dashboard — iteration 429 (2026-11-24)

- **Done this fire:** row 243 **V4 sprint 4 LANDED** (PR #651, merge `ee8ef1c0`, merge CI 3/3; judge r1 **FAIL 72** → r2/r3 **PASS 81**, 0 blocking).
  - The **opt-in operator-runner companion**: `ailang-worldd service runner install|start|stop|status|cancel|uninstall`, its own LaunchAgent label, private files, approved PATH/HOME, no stored tokens, a liveness-only heartbeat and an inert (unserved) catalogue. Nothing executes yet: no step dispatch, mint, launch, finish or route.
  - A **schema-compatible home rollback**: automatic `service upgrade` rollback now refuses a previous binary it cannot show can open the migrated store, keeping both builds for you instead of stopping the service. This addresses S3's v10 rollback risk once the upgrading binary carries the guard.
  - 1,539 changed lines, over the design's 1,400 cap by test-only authority pins (judge r1's blocking finding) and a Linux portability fix — deliberate, design §8.
- **Before your next `service upgrade` of the home service:** still back up `credentials.db` — the guard runs only in a binary that has it, and the loop did not measure which binary performs the rollback during an upgrade from today's install.
- **Install of the companion is attended and optional.** Nothing in the loop installs it; it does nothing useful until S5–S8.
- **DECISIONS FOR MARK:** none open (ledger 107 rows, 0 OPEN).
- **Next loop work:** **V4 S4b** — terminal exclusion (start/launch/finish refuse while the runner owns an episode) and the startup timeout, compiled in S4's plan; then S5–S12 in design §8 order.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** controller `claude:claude-opus-5-5`; planner and executor `codex:gpt-6.2-sol` (511k / 53k+82k+67k tok over three rounds); judge Agent `sonnet` (three rounds), cross-vendor. Metered $0. Codex window at start `used=42.0% allowed=59.6%`.
- **Open PRs:** #617 attended UI work (not the loop's; left alone). Attended #646–#649 merged during the sprint, no shared files.
- **Dev CI:** 3/3 green at merge `ee8ef1c0`. Red elsewhere on Linux CI recently (1–2 sightings each, not yet a flake ledger): `host/daemon` race timeout (also on dev `212a87ba`), `TestFinishConfirmOnSurfaceDispatchesTheNamedDecision`, marketingcheck Chrome timeout.
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).

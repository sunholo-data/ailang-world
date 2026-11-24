# World mission dashboard — iteration 428 (2026-11-24)

- **Done this fire:** row 243 **V4 sprint 3 LANDED** (PR #643, merge `e83d1717`, judge **PASS 86**, 0 blocking).
  - Credentials **schema 10**: `operator_actions`, `operator_runner_state` and `operator_action_steps`, including the finish shape so S9–S11 need no second one-way bump.
  - The durable operator-action store law: admission with dedup, a 16-action pending cap, one active action per episode, strict expiry, prune with tombstones; confirm CAS consuming the key-verified proof; claim CAS.
  - Unregistered: no route, no production caller, no `.ail`. 1,396 changed lines (cap 1,400).
- **Why it moved:** your attended rulings **D-WORLD-119 = A** and **D-WORLD-120 = A** (`82f05004`) ended the 18-fire rule-(d) run. Acknowledged; nothing re-asked.
- **Before your next `service upgrade` of the home World service:** schema 10 is one-way. If that upgrade's health check fails after the migration, the automatic rollback restores the v9 binary, which refuses the v10 store — the service stops (data intact, candidate kept as `.failed`). Back up `credentials.db` first. S4's upgrade seam owns a v10-safe rollback.
- **DECISIONS FOR MARK:** none open (ledger 107 rows, 0 OPEN).
- **Next loop work:** **V4 sprint 4** — runner install, isolated process/status/upgrade seam, heartbeat (E1/E3), carrying S3b (132-line step-receipt slice), the rollback fix and judge-428's test gaps. Then S5–S12 in design §8 order.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** controller `claude:claude-opus-5-5`; planner and executor `codex:gpt-6.2-sol` (479k / 80k tok); judge Agent `sonnet` (137k tok), cross-vendor. Metered $0. Codex window `used=41.0% allowed=56.3%`.
- **Open PRs:** #617 attended UI work (not the loop's; left alone). Attended agy PRs #641/#642 merged mid-sprint, no shared files.
- **Dev CI:** 3/3 green at merge `e83d1717`.
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).

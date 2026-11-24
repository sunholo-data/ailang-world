# World mission dashboard — iteration 430 (2026-11-24)

- **Done this fire:** row 243 **V4 sprint 4b LANDED** (PR #654, merge `fa3ad190`, merge CI 3/3; judge **PASS 85**, 0 blocking).
  - **Terminal exclusion:** `episode start` (plain and preset), `resume`, `launch` and the F1 terminal `finish` now take a per-episode lock and refuse at once, naming the exact action id, while an operator action owns that episode — even after the lock is released or the lease expires. Unrelated episodes still run.
  - **Bounded runner startup:** `service runner start` waits (bounded) for a fresh, compatible heartbeat generation and stops a job it started if none arrives.
  - **Pins** for the S4 guards the S4 judge found untested. 978 changed lines.
  - Nothing dispatches yet: no runner consumer, no route.
- **New files you may notice:** `credentials.db.operator-episode-<id>.lock` beside the credentials store, one per started episode (empty, 0600). Not yet in QUICKSTART — carried.
- **Before your next `service upgrade` of the home service:** still back up `credentials.db`.
- **DECISIONS FOR MARK:** none open (ledger 107 rows, 0 OPEN).
- **Next loop work:** **V4 S5** — pinned-start read seam and minted-start proof (design §5.3), carrying: `session new` and the mint-request `approve` path still provision without the lock (latent until a runner consumer exists); an attended resolver for an uncertain action; S4c-replay before S7; S3b.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** controller `claude:claude-opus-5-5`; planner and executor `codex:gpt-6.2-sol` (381k / 56k tok); judge Agent `sonnet`, cross-vendor. Metered $0.
- **Dev CI:** 3/3 green at merge `fa3ad190`. Fire started on a red dev race gate: `host/daemon`'s 8-minute package timeout on a slow runner, third instance → filed as **row 267** (hygiene, under rule (d)); rerun green. Row 266 (`TestFinishConfirmOnSurfaceDispatchesTheNamedDecision`) is the other open race flake.
- **Open PRs:** #617 attended UI work, #650 row 228 search (neither touched by the loop). #653 (broker) merged during the sprint, no shared files.
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).

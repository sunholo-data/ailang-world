# World mission dashboard — iteration 431 (2026-11-25)

- **Done this fire:** row 243 **V4 sprint 5 LANDED** (PR #657, merge `ab520385`, merge CI 3/3; judge **PASS 86**, 0 blocking).
  - **Lock gap closed:** `session new` and the signed mint-request `approve` now take the per-episode lock before provisioning and refuse an operator-owned episode by exact action id, consuming nothing.
  - **Read seam:** a read-only preset-start resolver shared by the terminal `episode start --preset` and the future operator start. Terminal behaviour unchanged.
  - **Confirm replay is exact:** a replay with different artifact bytes is refused (S4c).
  - QUICKSTART now documents the `credentials.db.operator-episode-<id>.lock` sidecars. 743 changed lines. Nothing dispatches yet.
- **Split:** the atomic minted start (design §5.3) moved to **S5b** — measured 1,432 lines against the 1,400 cap.
- **Before your next `service upgrade` of the home service:** still back up `credentials.db`.
- **DECISIONS FOR MARK:** none open (ledger 107 rows, 0 OPEN).
- **Next loop work:** **V4 S5b** — atomic minted-start proof and recovery, an attended resolver before any new uncertain outcome, the paired-process lock consumer. Carried: `session mint` still mints without the episode lock (latent until a runner consumer exists).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** controller `claude:claude-opus-5-5`; planner and executor `codex:gpt-6.2-sol` (398k / 65k tok); judge Agent `sonnet`, cross-vendor. Metered $0.
- **Dev CI:** 3/3 green at `cd06d950` at start. PR #657's race gate hit row 267 (`host/daemon` 8-minute timeout on a slow runner) — **instance 5**, rerun green; the sprint touches no `host/daemon` code.
- **Open PRs:** #617 attended UI work, #650 row 228 search (neither touched by the loop).
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).

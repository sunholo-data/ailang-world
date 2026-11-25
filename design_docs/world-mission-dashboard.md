# World mission dashboard — iteration 434 (2026-11-25)

- **Done this fire:** row 243 **V4 sprint 5d LANDED** (PR #666, merge `bb214516`, merge CI 3/3 on attempt 2; judge **PASS 84**, 0 blocking).
  - **Step bank:** the runner can now reserve and observe a numbered step of a signed action (S3b), with restart and guard tests.
  - **Worktree reservation:** a signed start can reserve its worktree, create it, and check what Git really did, all while holding the episode lock. It has **no production caller** yet.
  - **A way out of every uncertain result:** `ailang-worldd service runner reconcile-worktree ACTION-ID` shows what Git holds, asks for a typed `y`, and records keep or stop. It never deletes a worktree, branch or file. The text it prints is run verbatim by a test, and QUICKSTART documents it.
  - 1,397 changed lines (cap 1,400). Nothing dispatches yet.
- **Found this fire (row 271):** 14 real-binary end-to-end tests still expect 10 se-tools after row 228 (#650, attended) added `workspace-search`/`workspace-list`. CI skips these tests, so it stays green, but local `verify_go` is red at every head since `f8fb3ac1`. One open question for the fix: should a Read-only grant list the two new tools?
- **Before your next `service upgrade` of the home service:** still back up `credentials.db`.
- **DECISIONS FOR MARK:** none open (ledger 107 rows, 0 OPEN).
- **Next loop work:** **V4 S5e**. Its first milestone adds the judge's missing guard tests on the reservation path (store checks, Git readback shape, reconcile tamper checks) before anything calls it. Then file publication with its attended disposition. Row 271 is small and keeps local gates readable.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** planner and executor `codex:gpt-6.2-sol` (435k + 73k tok), judge Agent `sonnet` (120k tok): cross-vendor. Metered $0.
- **Dev CI:** race job red at the start on two attended row-269 commits (row 267 instances 7 and 8, `host/daemon` 8-minute timeout on slow runners); attempt 2 green. PR #666 3/3 on attempt 1. The merge's host-build job failed once on a known load-sensitive `host/broker` test (package untouched); attempt 2 green.
- **Open PRs:** #617 attended UI work (not touched by the loop).
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).

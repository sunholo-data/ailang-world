# World mission dashboard — iteration 438 (2026-11-26)

- **Done this fire:** row 243 **V4 sprint 5e4 LANDED** (PR #676, merge `e554db6b`). PR CI was 3/3 on the first attempt. The judge gave **PASS 88** with 0 blocking findings.
  - **The pre-caller pins are in.** All sixteen store survivor edits the last two judges carried are now killed by `host/store` tests. They include N33/N34, so `reconcile-mint` is now pinned to refuse once files are reserved.
  - **Episode GC no longer removes an operator worktree it still needs.** The daemon's 6-hourly pass and `ailang-worldd episode gc` both skip a worktree whose episode has an unresolved operator step. If the step lookup is unavailable, they skip too.
  - 929 lines, under the 1,400 cap. Nothing dispatches yet: no production caller.
- **Judge's open points (non-blocking):**
  - The GC lookup's tests cover only worktree steps. Dropping `files` or `mint` from its query would go unnoticed; the controller reproduced this. The query is correct today, and the pin is carried into S5f.
  - A home with no operator-capable credentials store (absent, older than schema 10, or locked) now reclaims no worktrees. This is fail-safe and the reason is printed.
- **Instrument note:** under rig load the controller's base runs of `cmd/ailang-worldd -race` and `verify_go` hit Go's 10-minute timeout. They were discarded and CI arbitrated those legs.
- **Fleet fixes received:** two harness tickets are resolved upstream (rows 268 and 270). Row 270's `rotate-log --workdir` reaches this rig only once the fleet `ailang` is rebuilt.
- **Before your next `service upgrade` of the home service:** back up `credentials.db`.
- **DECISIONS FOR MARK:** none open (ledger 107 rows, 0 OPEN).
- **Next loop work:** **V4 S5f**, the paired-process consumer and drain with the first production caller. It carries:
  - the GC lookup's `files`/`mint` test rows
  - N04/N12/N15 (82 lines, parked)
  - `E-RES-EPISODE-CONFLICT`, `E-RES-REREAD`, J-OWN-UID, Z-ABSENT-ANY, Z-WAIT-READBUDGET and the four process pins

  Row 271 (14 local-only e2e tests expecting 10 se-tools) is still small.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** planner and executor `codex:gpt-6.2-sol` (465k + 60k tok). The judge was Agent `sonnet` (112k tok), a different vendor. Metered $0.
- **CI headroom:** on PR #676 the race job took 262.9 s for `cmd/ailang-worldd` and 355.8 s for `host/daemon`, against the 480 s bound (row 267).
- **Open PRs:** #617 is attended UI work that the loop does not touch.
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).

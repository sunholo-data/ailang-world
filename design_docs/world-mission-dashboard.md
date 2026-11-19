# World mission dashboard — iteration 413 (2026-11-19)

- **Done this fire:** nothing built. This is the fourth rule-(d) block report in a row (410–413). No row can move without an answer from Mark.
  - Row 243 V4 sprints 3–12 all need **D-WORLD-119** (9–11 also need **D-WORLD-120**).
  - Release rows 150 and 151 are not routable. 150 waits on your attended A1. 151 waits on row 188's N2d and N4a, which are not built.
- **DECISIONS FOR MARK (open, unchanged since iteration 406):**
  - **D-WORLD-119:** may the app start and launch a confined episode through the runner after a key approval? Recommendation A. Every remaining V4 sprint waits on it.
  - **D-WORLD-120:** may finish run without a terminal (key for local prep, then key for the exact PR)? Recommendation A. Sprints 9–11 also wait on it.
- **Why every idle fire still costs a full run:** the codex probe has failed on every fire since 410, which counts as "degraded", and a degraded fire may not record the idle fingerprint. The fleet ticket `skill:idle-record-barred-on-degraded-no-pick-fire` is open (count now 3; queue row 263). The codex failure looks like the fleet's D-FLEET-47.
- **Next loop work:** with D-WORLD-119 = A or B, **V4 sprint 3** (schema 10 for operator actions, receipts and CAS). Unanswered, another block report.
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved.
- **Open PRs #615** (`ui/finish-card-review`) **and #617** (`ui/live-branch-heads`, new today) are attended UI work, not the loop's; left alone.
- **Routing this fire:** controller opus; designer, planner and executor not needed; judge sonnet via Agent on the record. The codex lane probed rc 1. Metered $0.
- **Dev CI:** 3/3 green at `31413192`.
- **Bookkeeping thread:** #574 (rotates Monday 2026-11-23).

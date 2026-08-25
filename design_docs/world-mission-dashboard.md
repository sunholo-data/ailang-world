# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in the charter STATUS + `world-mission-log.md`.*

**Updated**: 2026-08-25 (iteration 123) · `dev` @ `3dda87e` · CI **GREEN** (`checks=2`, both success)

## Where we are

- **Queue item 14** `w-workbench-read-only` is **COMPLETE — 11 of 11 milestones**, all **32**
  mutation rows discharged by measurement. Doc + sprint plan moved to `design_docs/implemented/`.
- **`WB.K` (drill 4/4) killed M24–M28 with no survivor** — the second clean drill milestone of four.
- **The §7h retrospective judge pass owed by iteration 122 is discharged**: evaluator `sonnet`,
  **93/100 PASS, zero blocking**, twelve named targets across §7h and §7i, none refuted.
- **Queue row 5** `w-mcp-projection` is now the **queue head**: unblocked at source
  (`ailang#764` CLOSED) and ordered here by Mark's **"Finish 14"** (`D-WORLD-25` arm B), which item
  14 has now satisfied. Sole remaining blocker on **M4**, the reference-agent value gate.
- **Latest upstream release**: AILANG **v0.33.2**; pinned `.ail` compiler v0.30.0 (separate axis).

## In flight / next

- **NEXT: row 5** `w-mcp-projection`. Its **first milestone is the toolchain precondition**, not a
  decision: v0.33.2 declares `go 1.26.6` while CI pins `GOTOOLCHAIN: go1.25.6`; the repo's own
  canary clears the move (`go1.25.6` rc=0, `go1.26.5` rc=1 miscompile, `go1.26.6` rc=0).
  Then rows 38, 37, 36, 35, 34, 32, 33, item 22, row 31.
- **Three of four landing legs read GREEN on a mutation that landed inside a comment** (§7i(a)).
  Only a query against the file's **parsed form** refused it. Shape-specific: a substitution's
  two-sided count predicate already refuses; the exposure is the one-sided predicate an *insertion*
  forces, and M24 was the catalogue's only insertion-shaped row.
- **Row 34 now carries seven hunks.** Sixth and seventh added this iteration by the evaluator:
  `supportedWorkbenchQuery`'s pair-composition guards (`workbench.go:72`/`:75`, `&&` → `||`) both
  **SURVIVE** with an empty red set, and both are live on the function's own truth table.

## Loop posture

- **Cadence**: launchd `dev.ailang.mission-world`, staggered vs the V1 loop.
- **Routing**: controller `claude:claude-opus-5`; evaluator **`sonnet`** (ran this iteration —
  generator≠judge satisfied, and iteration 122's capacity gap is closed retrospectively); executor
  chain **codex → opus** (D-WORLD-20 suspended the DeepSeek lane). Drill milestones are
  controller-work by construction (§7f(b): loopback-binding classification arm).
- **⚠ The RUNNING shared mission-control skill is 27 lines behind `origin/dev`** (3,757 vs 3,784,
  hunk `origin:1108–1134` — V1 iteration 274's mutation-intended-effect rule). Read from origin and
  applied anyway. **World cannot repair it**: the V1 checkout is off-limits from this repo. For V1.
- **Quota / cost**: metered **$0.00** of $5 (opus ×1 controller, sonnet ×1 evaluator). Fable +
  designer rotation unspent a **19th** consecutive iteration.
- **Bookkeeping issue**: **#89** (week of 2026-08-24; predecessor #68). Open asks: **0**.

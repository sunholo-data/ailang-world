# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in the charter STATUS + `world-mission-log.md`.*

**Updated**: 2026-08-25 (iteration 121) · `dev` @ `2e7154b` · CI **GREEN** (`checks=2`, both success)

## Where we are

- **Mark answered the ordering fork**: `#89` @ `2026-08-24T23:14:21Z`, verbatim **"Finish 14"** —
  **D-WORLD-25 arm B**. Item 14 completes before row 5. **Open asks are back to ZERO.**
- **Queue item 14** `w-workbench-read-only`: `[IN-SPRINT]`, **9 of 11** landed (`WB.I` this
  iteration). `WB.J` and `WB.K` remain — both **controller-work** (their classification arm binds
  loopback; the sandboxed executor lane denies it).
- **Queue row 5** `w-mcp-projection`: **UNBLOCKED** since iter-120, deliberately **not started** —
  it is next after item 14 closes. Sole remaining blocker on **M4**, the value gate.
- **Latest upstream release**: AILANG **v0.33.2**. **Pinned `.ail` compiler**: v0.30.0 at
  `/tmp/ailang-v0300/ailang` (a separate axis from the Go library dependency).

## In flight / next

- **NEXT: `WB.J`** (10 of 11 — mutation drill 3/4, discharges M10–M13, M22, M23, M29–M32), then
  `WB.K`, then rows 38, 37, 36, 35, 34, 32, 33, item 22, row 31, then **row 5**.
- **`WB.I` found a mutant that cannot be killed.** M9's guard is masked by a sibling guard on the
  same request path, so no single-site mutation of it is detectable. Recorded SURVIVED per the
  protocol; residue stays with **queue row 37**, not absorbed.
- **Row 5 carries a toolchain precondition, not a redesign**: v0.33.2 declares `go 1.26.6`; CI pins
  `GOTOOLCHAIN: go1.25.6`. The repo's own canary already clears the move.

## Loop posture

- **Cadence**: launchd `dev.ailang.mission-world`, staggered vs the V1 loop.
- **Routing**: controller `claude:claude-opus-5`; executor chain **codex → opus** (D-WORLD-20
  suspended the DeepSeek lane); evaluator Sonnet (**97/100 PASS, zero blocking** this iteration).
- **Quota / cost**: metered **$0.00** of $5 (opus ×1 controller, sonnet ×1 evaluator; no executor,
  planner or designer lane spent). Fable + designer rotation unspent a **17th** consecutive iteration.
- **Bookkeeping issue**: **#89** (week of 2026-08-24; predecessor #68). Open asks: **0**.

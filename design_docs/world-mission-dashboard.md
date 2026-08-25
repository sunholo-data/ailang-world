# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in the charter STATUS + `world-mission-log.md`.*

**Updated**: 2026-08-25 (iteration 122) · `dev` @ `b0d973c` · CI **GREEN** (`checks=2`, both success)

## Where we are

- **Queue item 14** `w-workbench-read-only`: `[IN-SPRINT]`, **10 of 11** landed (`WB.J` this
  iteration). Only `WB.K` remains — **controller-work** (its classification arm binds loopback,
  which the sandboxed executor lane denies).
- **`WB.J` is the first drill milestone of four to end with NO surviving mutant.** All ten rows
  discharged: M10–M13, M22, M23, M29–M32.
- **Queue row 5** `w-mcp-projection`: **UNBLOCKED** — `ailang#764` now reads **CLOSED** (re-run as a
  command, not transcribed). Deliberately not started: Mark's **"Finish 14"** (`D-WORLD-25` arm B)
  puts it immediately after item 14 closes. Sole remaining blocker on **M4**, the value gate.
- **Latest upstream release**: AILANG **v0.33.2**; pinned `.ail` compiler v0.30.0 (separate axis).

## In flight / next

- **NEXT: `WB.K`** (11 of 11 — drill 4/4, discharges M24–M28, plus full gates and final
  acceptance), then rows 38, 37, 36, 35, 34, 32, 33, item 22, row 31, then **row 5**.
- **`WB.K` owes a retrospective judge pass over §7h** — see the gap below.
- **`M12` is killed by a hardcoded literal; the assertion its row names cannot detect it.** The
  count comparison reads `workbench.WorkbenchPageLimit` — the constant M12 mutates — so both sides
  move together. Pin is live but carries a decorative member; residue → **row 34**.
- **The §7d(c) tautology tell (iter-113) was never swept.** Swept now: **one** hit tree-wide.

## Loop posture

- **Cadence**: launchd `dev.ailang.mission-world`, staggered vs the V1 loop.
- **Routing**: controller `claude:claude-opus-5`; executor chain **codex → opus** (D-WORLD-20
  suspended the DeepSeek lane).
- **⚠ GAP: the evaluator lane did NOT run this iteration** — agent spawning was disabled for the
  session, so **generator≠judge was not satisfied**. A capacity gap (Standing rule 8), not a
  judgment one: resumes on any iteration with the lane available, no human answer needed. Not
  harmless — the iter-119 and iter-121 judges each found a real survivor I had missed (**2 of 2**).
- **Quota / cost**: metered **$0.00** of $5 (opus ×1 controller only). Fable + designer rotation
  unspent an **18th** consecutive iteration.
- **Bookkeeping issue**: **#89** (week of 2026-08-24; predecessor #68). Open asks: **0**.

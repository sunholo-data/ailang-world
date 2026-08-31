# Mission Dashboard — Ailang World

*Snapshot, overwritten each iteration. History: `world-mission.md` STATUS + `world-mission-log.md`.*
**As of** 2026-08-31, iteration **141** · `dev` = [`d195717`](https://github.com/sunholo-data/ailang-world/commit/d195717) · CI **GREEN** (2/2 on the merge commit)

## Just landed
- **Row 51** — the floor-raise coupling-inventory gate was blind to an **asymmetric addition**: a 7th
  coupled-file row in one of the two homes passed with no count moving. Now a symmetric
  set-difference + per-home duplicate guard + a `>= 6` instrument-failure floor.
  PR [#108](https://github.com/sunholo-data/ailang-world/pull/108) · evaluator `sonnet` **97/100 PASS**,
  zero blocking · **6 real mutants were invisible to the old test** (7 of 9 arms survived it).

## Next up (ready, gated on nothing)
1. **Row 52** — wiring test attributes a step to the previous one under key reorder.
2. **Row 53** — a quorum reviewer can be silenced by its own review content.
3. **Row 59** *(new)* — a `grep -c` cannot prove an assertion is *live*: wrap row 51's own new block
   in `if false {}` and its acceptance criterion still reads green. Row 49's defect one layer up.

## Parked on Mark — 2 open decisions, both one word
- **`D-WORLD-28`** — how should `verify_go.sh` guarantee its nested race-control module can execute?
  **A (recommended)**: fail closed unless `ACTIVE_GO` is at-or-above the root module floor. **B**: keep
  ambient auto-selection, pin `racecontrol/go.mod` independently.
- **`D-WORLD-29`** — after the whitespace-tolerant rewrite of `shellAssignmentValues`, should a single
  *indented* assignment be **ACCEPTED** (**A**, recommended) or **REJECTED** (**B**)? Answering also
  amends ratified queue row 50, which is why the loop may not settle it.
*Default if unanswered: rows 48 and 50 stay parked and the queue advances past them.*

## Loop cadence + routing
Nightly launchd, one queue item per iteration. Controller `opus` · designer **rotation**
(`claude:claude-fable-5` ↔ `pi:ollama/deepseek-v4-flash:0731-cloud` — the pi lane is now **2-for-2**
with typed `ok` verdicts and non-empty diffs) · planner `opus` · executor `codex:gpt-5.6-sol` ·
evaluator `sonnet`. generator≠judge held four ways this iteration.

## Quota posture
`metered = $0.3427` of the `$5` iteration ceiling, all of it design quorum. Every other lane is a
subscription/flat-rate bucket. Pinned `ailang` v0.30.0 at `/tmp/ailang-v0300`.

## Standing rig caveat
`scripts/verify_go.sh` is **FLAKY on this rig**, not deterministically red — 4 runs in ~90 min gave
2 red / 2 green over 3 different failing sets while dev CI stayed green. An `rc=0` criterion for it
is unsafe in both directions; use a set comparison. Row **58**.

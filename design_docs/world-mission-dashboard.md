# Mission Dashboard — AILANG World

_Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `world-mission-log.md`._

**As of:** 2026-09-03 · iter 150 · `dev` = [`a7b58dd`](https://github.com/sunholo-data/ailang-world/commit/a7b58dd) · CI green (3/3 on the merge commit)

## Last iteration
**Row 67 LANDED** — PR [#114](https://github.com/sunholo-data/ailang-world/pull/114) → squash `a7b58dd` ·
Gate 3b GREEN on the merge commit (`present=3 == expected=3`, all `success`) · **HARNESS** ·
metered **$0.00** (controller-authored direct fix; no sub-agent spawned).

`dev` was RED — the first red here this loop did not write. Two fine branches met: one added a third CI
job with no Go in it (`launchd-drivers`, bash 3.2), the other carried a gate whose one `wantJobs` constant
answered both *which jobs exist* and *how many Go pins there must be*. Now two lists, so a new job must be
classified in both edits, a job classified OUT is asserted to carry **zero** Go pins, and pins are
attributed **per job** — closing a fail-open the count-only gate always had.

## Goal distance
**58 of 70 rows closed** (57 of 66 at iter-149; row 67 opened and closed this iteration, rows 68/69/70
open from its own measurements). Row 50 parked on `D-WORLD-31`.

## Next picks
1. **Row 57** `w-approvals-spine-prints-a-green-no-pending` — queue head, ungated.
2. **Row 69** `w-heartbeat-script-absent` — the skill's mandated per-gate stamp is `No such file or
   directory` here, so rule 7's attribution contract rides on controller discipline. Fleet-owned port.
3. **Row 58** the known `verify_go.sh` flake. Then 59–66, 68, 70, then 39.

## Loop + routing
Controller `claude:claude-opus-5` · designer ROTATION, last used **`claude:claude-fable-5`** (next:
`pi:ollama/deepseek-v4-flash:0731-cloud`) · planner `opus` · executor `codex:gpt-5.6-sol` · evaluator
`sonnet` (generator≠judge). Iter-150 spawned **no** role — a ~40-line test-file change decided by a
mutation drill the controller runs anyway; no independent judge saw it, stated as the narrowing it is.

## Parked on Mark
**`D-WORLD-31`** — one word. Ship `D-WORLD-29`'s rule A as ratified, or hold row 50 for the fixture
migration. Nothing else is blocked; the queue advances either way.

## Standing reds — owned elsewhere, none is a World failure
- **`verify_go.sh` RED on the rig:** (a) the fleet arm — World's driver copy is behind fleet HEAD, i.e.
  *the fleet must commit*; (b) row 58's known FLAKE. **CI is unaffected** (the fleet arm loud-skips there).
- **The running `mission-control` skill is byte-identical to `origin/dev`** (`cmp` against the resolved
  symlink target) — iter-149's 187-line staleness is **CLOSED**.
- **`tools/launchd/mission-heartbeat.sh` is absent here** (row 69) — stamps made via V1's absolute path.

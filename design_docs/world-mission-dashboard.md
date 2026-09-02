# Mission Dashboard — AILANG World

_Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS) and
`world-mission-log.md`._

**As of:** 2026-09-02 · iteration 148 · `dev` = [`14036ee`](https://github.com/sunholo-data/ailang-world/commit/14036ee) · CI green (2/2)

## In flight / just landed
- **LANDED iter-148 — row 54 `w-driver-copy-stale-and-the-drift-gate-compares-it-to-itself`**
  (PRODUCT/HARNESS boundary → **HARNESS**). PR [#111](https://github.com/sunholo-data/ailang-world/pull/111)
  → squash `14036ee`, Gate 3b green on the merge commit. The `D-WORLD-DRIVER-1` gate compared the
  driver copy to itself, so **11 commits / 705 differing lines** of staleness read green. Adds a
  fleet-comparison arm (+ `--driver-fleet-check`). Evaluator **86/100**, one blocking finding fixed
  in-sprint.
- **Expected and intended:** `./scripts/verify_go.sh` on the RIG is now **RED** until the fleet
  lands a current driver. That red means "the fleet must commit", never "absorb it into World".
  CI is unaffected (the arm loud-skips; measured in the ubuntu job log, not assumed).

## Next picks
1. **55** `w-dispatch-lever-parser-false-reds-on-valid-yaml` · ~0.3d
2. **56** `w-canary-fence-blind-to-a-skipped-canary`
3. **57** `w-approvals-spine-prints-a-green-no-pending-under-the-row-it-just-listed`
then 58–64, then **39** `w-session-authority`. Row **50** stays parked on `D-WORLD-31`.

## Parked on Mark
- **`D-WORLD-31`** (row 50) — ONE WORD: ship the ratified rule A as-is, or hold row 50 for the
  fixture migration? Unchanged since iter-146. **No new ask this iteration.**
- `design_docs/AI-EMPLOYEE.md` remains an attended draft that says it "steers nothing until placed
  on #1" — still not treated as a directive.

## Loop cadence + routing
- Fires ~2-hourly via `launchd`; one backlog item per iteration.
- designer **rotation** (`claude:claude-fable-5` → `pi:ollama/deepseek-v4-flash:0731-cloud`);
  this iteration used **deepseek** (flat-rate, $0) — pointer now advanced to deepseek.
- planner `opus` (`derive-planner-lane.sh` → `opus fail-closed:env-pin`) · executor
  `codex:gpt-5.6-sol` · evaluator `sonnet` (generator≠judge).

## Quota / spend
- **metered this iteration: $0.26775** of the $5 ceiling — all of it quorum reviewers
  (2 full-strength rounds + one restored reviewer). Everything else rode quota buckets or flat-rate.
- Decision ledger: **18 rows, `--check` valid, ONE OPEN** (`D-WORLD-31`).

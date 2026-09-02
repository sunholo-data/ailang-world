# Mission Dashboard — AILANG World

_Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `world-mission-log.md`._

**As of:** 2026-09-02 · iter 149 · `dev` = [`165b9fd`](https://github.com/sunholo-data/ailang-world/commit/165b9fd) · CI green (2/2)

## Last iteration
**Row 55 LANDED** — PR [#112](https://github.com/sunholo-data/ailang-world/pull/112) → squash `165b9fd` ·
Gate 3b GREEN on the merge commit (`present=2 == expected=2`, both `success`) · evaluator `sonnet`
**95/100**, zero blocking · **HARNESS** · metered **$0.2211** of $5.

The dispatch-lever gate red-lit three shapes of *valid, lever-declaring* YAML — including the standard
remedy for a famous Actions footgun (quoting `on` so YAML 1.1 does not read it as `true`). Quoted and
flow-style keys are now GREEN; a tab-indented trigger became a loud TYPED failure instead of one falsely
claiming the block was absent. Line scan, no YAML parser, no new dependency.

## Goal distance
**57 of 66 rows closed** (56 of 64 at iter-148; row 55 closes, rows 65/66 open from this iteration's
planner and judge findings). Row 50 parked on `D-WORLD-31`.

## Next picks
1. **Row 56** `w-canary-fence-blind-to-a-skipped-canary` — queue head, ungated.
2. **Row 65** `w-go-build-is-not-a-compile-fence-for-a-test-file` — new; invalidates a mutant-BUILDS
   fence used across this repo's plans *and* prescribed by the shared skill.
3. **Row 57** `w-approvals-spine-prints-a-green-no-pending`. Then 58–64, 66, then 39.

## Loop + routing
Controller `claude:claude-opus-5` · designer ROTATION, last used **`claude:claude-fable-5`** (next:
`pi:ollama/deepseek-v4-flash:0731-cloud`) · planner `opus` (`opus fail-closed:env-pin`) · executor
`codex:gpt-5.6-sol` · evaluator `sonnet`, own worktree (generator≠judge).

## Parked on Mark
**`D-WORLD-31`** — one word. Ship `D-WORLD-29`'s rule A as ratified, or hold row 50 for the fixture
migration. Nothing else is blocked; the queue advances either way.

## Standing reds — both owned elsewhere, neither is a World failure
- **`verify_go.sh` RED on the rig:** (a) the fleet arm — the World driver copy is **759** diff-lines
  behind the fleet (705 yesterday; it grows every fire), i.e. *the fleet must commit*; (b) row 58's
  known FLAKE `TestHandlerTimeoutKillsTheWholeProcessGroup`, attributed by mechanism this iteration.
  **CI is unaffected** — the fleet arm loud-skips there.
- **The running `mission-control` skill is 187 lines behind origin** (was 147) — V1's stale checkout,
  which World cannot fix. World reads the delta and follows origin's rules.

# Mission Dashboard — Ailang World

*Snapshot, overwritten each iteration. History: `world-mission.md` STATUS + `world-mission-log.md`.*
**As of** 2026-09-01, iteration **143** · `dev` = [`a362624`](https://github.com/sunholo-data/ailang-world/commit/a362624) · CI **GREEN** (2/2)

## This iteration (REFUTATION — nothing shipped here; the fix is upstream)
- **Row 53 ROUTED, World-side complete** → [`ailang#941`](https://github.com/sunholo-data/ailang/issues/941).
  The row blamed the reviewer's JSON escaping for quoting `"GOOS"`/`"GOARCH"`. Measured on the
  artifact: those literals are **escaped correctly** — the break is past the 200-char window the
  error publishes, and the offset that would locate it is dropped on the way out.
- **The class is bigger than the row.** Sweep of **193 quorum artifacts / 389 reviewer rows** across
  both missions: **6 of 6 `invalid` absences carry a recoverable `"verdict":"reject"`** already
  inside the published window (negative control on the 15 `budget`/`unreachable` absences: **0**).
  **Three of those six synthesised `proceed` with zero external reviewers present** — the same three
  vacuous passes the shared skill already documents. A discarded reject became a green.

## Next picks (ready, ungated)
1. **Row 54** — the drift gate that guards frozen core compares the driver copy to *itself*; the
   copy is 8 commits / 430 lines behind the fleet. World's own `verify_go.sh` is in scope.
2. **Row 55** — the row-47 lever parser false-reds on three forms of valid YAML.
3. **Rows 56, 57** — canary fence blind to a skipped canary; approvals spine prints "no pending".
4. **Rows 58, 59** — `verify_go.sh` is flaky on this rig; a `grep -c` cannot prove an assertion live.
Then row **39**.

## Loop cadence + routing
Controller `claude:claude-opus-5`. Designer rotation `claude:claude-fable-5` ⇄
`pi:ollama/deepseek-v4-flash:0731-cloud` — **rotation NOT advanced this iteration** (nothing was
authored), so deepseek is still next. Planner `opus`, executor `codex:gpt-5.6-sol`, evaluator
`sonnet`.

## Parked on Mark — THREE open asks, all unchanged
- **`D-WORLD-28`** — how should `verify_go.sh` guarantee its nested race-control module can execute?
- **`D-WORLD-29`** — should a single *indented* shell assignment be ACCEPTED or REJECTED?
- **`D-WORLD-30`** — row-52 fix: **LINE SCAN** (recommended) or **YAML PARSE**?
Rows **48**, **50**, **52** are the items these block. No new ask this iteration.

## Quota / spend posture
`metered=$0.00` of the $5 ceiling — no quorum was owed and no metered lane ran. Billing tripwire
**CLEAN**; nested `claude` calls go through `claude-sub`, which strips the API keys by construction.

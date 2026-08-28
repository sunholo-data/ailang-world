# Mission Dashboard — AILANG World

> 30-second control context. OVERWRITTEN every iteration; history lives in the charter STATUS
> block and `design_docs/world-mission-log.md`. Namespaced path — never write the bare
> `design_docs/mission-dashboard.md`, which is fleet-shared.

**Last iteration:** 136 · 2026-08-28 · **`P47` LANDED**

## Latest landing
- PR [#103](https://github.com/sunholo-data/ailang-world/pull/103) → squash
  [`2a01c43`](https://github.com/sunholo-data/ailang-world/commit/2a01c43). Queue row **47 CLOSED**.
- `ci.yml` now declares `workflow_dispatch:`; `host/verifygate/dispatch_lever_gate_test.go` pins it
  repo-wide with an anti-vacuity floor. A dropped webhook to `dev` is now recoverable — for the
  **tip of a named ref**, never an arbitrary SHA, and it buys a verdict on a commit, not a
  mergeable PR.
- Gate 3b GREEN on the merge commit: `present=2 == expected=2` (enumerated from `ci.yml`'s own
  `jobs:` block; 1 workflow file, so the enumeration is complete), `notdone=0`, `notgreen=0`,
  `runs=1`, parent control `checks=2` rev-parsed.

## In flight / next
- **NEXT:** rows **48**, **49**, **50**, **51**, **52**, **53**, **54**, **55**, then **39**.
- New row **55** this iteration: the dispatch-lever parser false-reds on valid YAML shapes
  (quoted `"on":`, flow-style, tab indentation). Fails loud, never a silent pass.
- Row **54** still open and still live: `launchd` runs THIS repo's driver copy, which is stale
  against the fleet's, and `verify_go.sh`'s drift gate compares that copy to itself.
  FLEET-owned (`D-WORLD-DRIVER-1`) — reported, not touched.

## Loop cadence + routing
- Controller `claude:claude-opus-5`. Designer **rotation advanced to
  `pi:ollama/deepseek-v4-flash:0731-cloud`** — its FIRST authoring run under the 2026-08-28
  amendment, verdict `ok` twice (authoring + revision), **$0** (flat-rate).
- Planner `opus` (lane derived verbatim: `opus fail-closed:env-pin`). Executor
  `codex:gpt-5.6-sol`, `metered=$0`. Evaluator `sonnet` **82/100 PASS, zero blocking**
  (generator≠judge holds: executor is OpenAI's).
- Quorum: **two rounds, both BLOCKED at full strength** (`absent_reviewers` EMPTY both times,
  3/3 external present); closed under the ratified narrow-refinement carve-out.

## Quota / spend
- `metered=$0.1900` of `$5` (quorum reviewers only; both provider lanes were flat-rate or quota).
- Billing tripwire **CLEAN**, re-checked before every nested provider call.

## Parked on Mark
- **ZERO OPEN ASKS.** Decision ledger: 14 rows, **0 OPEN**.

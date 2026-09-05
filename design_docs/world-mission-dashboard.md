# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` and `world-mission-log.md`.*

**As of** 2026-09-06, iteration **158** · bookkeeping issue
[#107](https://github.com/sunholo-data/ailang-world/issues/107)

## Where the mission is

- **Latest landed**: row **62** (PR [#119](https://github.com/sunholo-data/ailang-world/pull/119)
  -> squash `dcf534f`, CI green 3/3 on the merge commit) — the CI flag scan now reads the guarded
  step's own mapping key and its VALUE, so an explicit `continue-on-error: false` and a comment
  or `echo` merely naming the flag no longer red, while `true`, an empty value and a `${{ }}`
  expression still refuse.
- **Goal distance**: charter clause-2 gate-hardening queue. Rows 58–62 landed on five
  consecutive iterations; head moves to row **63**.
- **In flight**: none. No open PRs, no stale worktrees.
- **New rows filed**: none this iteration.

## Up next (banked, ready)

1. **63** — `w-locator-derivation-refusals-are-unpinned-and-undeclared` (clause-2, same gate)
2. **64–66** — clause-2 gate-hardening rows
3. **68–78** — gate/infra rows, incl. **76** (verify_go.sh rc=1 at base on the fleet-owned drift
   arm) and **69** (heartbeat/directives/resolver/pi-runner scripts absent in this repo)
4. **81**, **82**, **83**, then **39**

**Parked for design review**: rows **79**, **80** (by their own text — not human-blocked).

## Loop cadence + routing

- Controller `claude:claude-opus-5`. Designer rotation pointer unchanged at
  `pi:ollama/deepseek-v4-flash:0731-cloud` — no designer spawned. (The rotation itself was
  amended attended on 2026-09-05 to add `codex:gpt-6-astra` as a third, fable-class entry.)
- Recent picks have been controller-authored direct fixes on ~0.1–0.2d rows: **no independent
  judge runs on those**, stated explicitly each time, with landed-and-restored mutation drills as
  the compensating discipline.
- Verify gate = `verify_ail.sh` + `go build/vet/test` with `AILANG_BIN` set. `verify_go.sh` is
  **rc=1 at base** on the FLEET-OWNED driver-drift arm (row 76), re-measured this iteration and
  now naming fleet HEAD `f516881a` — a fleet commit clears it, never a World edit.
- The shared main checkout sits behind `origin/dev` between iterations by construction (every
  landing goes by worktree + PR). Mission state is read from origin; no reconcile is attempted
  without a human standing authorisation.

## Parked on Mark

**Nothing.** Decision ledger: **18 rows, `--check` valid, ZERO OPEN**. Five consecutive
iterations have asked for nothing.

## Quota posture

`metered=$0.00` of the $5 per-iteration ceiling; every lane used was a subscription/quota bucket.
Billing tripwire CLEAN (`ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN` both unset).

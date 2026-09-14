# Mission Dashboard — Ailang World

_Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` and `world-mission-log.md`._

**As of**: iteration **174** — 2026-09-14 · repo `sunholo-data/ailang-world` · branch `dev` @ `50103fd` (+ M2 `70b19e0` + record)

## Where the mission is

- **Latest landing**: row **73** — `scripts/gate0_self_notices.sh`, Gate 0's second, no-authority
  read for the driver's own crash notices (PR [#137](https://github.com/sunholo-data/ailang-world/pull/137)
  → squash `50103fd`, CI 2/2 green on the merge; evaluator round 2 **PASS 98/100 zero blocking**
  after a round-1 blocking finding was fixed). Its first live run read **rc=1 — A FIRE DIED** on
  iteration 173's own death notice (`#129`, 2026-09-08T07:33:01Z), which this iteration credited.
- **Iteration 173 was an orphan**: it designed row 73, then its slot was killed at gate-3 and the
  kill switch held the mission off for 30 fires (2026-09-09 → 2026-09-14). Credited in the log.
- **Queue census**: `census: 43 closed / 90 rows (tagged-open 4, untagged 43) [LANDED 42, ROUTED 1, RULED OUT 0; PARKED 3, IN-SPRINT 1, NEXT 0]`
  — moved by exactly one this iteration. Re-run:
  `bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79`
- **Milestones shipped**: M1 (semantic world library) · M2 (`ailang-worldd` daemon) — both complete.
- **Queue head**: rows **69** (fleet-gated), **74–78**, **81–90**, then **39**. Row **65** deferred per `D-WORLD-33`.

## Loop cadence and routing

- launchd `dev.ailang.mission-world`; kill switch `~/.ailang/state/mission-world.disabled` (armed).
- Bookkeeping issue rotates weekly (see `~/.ailang/state/mission-world-gh-issue`). Verify profile
  `ailang-code`; pinned binary `~/.pinned-ailang/ailang` **v0.30.0**.
- Gate 0 now runs TWO reads: `mission_directives.sh` (allowlisted humans) and
  `scripts/gate0_self_notices.sh` (the loop's own crash notices, no authority).
- This iteration's roles: designer none (orphan's fable doc) · planner `opus` · executor
  `pi:ollama/deepseek-v4-flash:0731-cloud` · evaluator `sonnet` ×2. Generator ≠ judge on model and vendor.

## Waiting on Mark

**Nothing.** Decision ledger: **22 rows, ZERO OPEN**. Fleet proposal filed as sunholo-data/ailang#1160 (a proposal, not an ask).

## Quota posture

Metered **$0.00** of the $5 ceiling this fire (no quorum ran). Quota: opus planner 229k tok, sonnet
judge 314k tok over two rounds, ollama executor flat-rate. Iteration 173's quorum had cost $0.215.

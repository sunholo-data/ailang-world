# Mission Dashboard — Ailang World

_Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` and `world-mission-log.md`._

**As of**: iteration **172** — 2026-09-08 · repo `sunholo-data/ailang-world` · branch `dev` @ `82d7ef7`

## Where the mission is

- **Latest landing**: row **72** — `scripts/queue_census.sh`, the queue-closure census instrument
  (PR [#136](https://github.com/sunholo-data/ailang-world/pull/136) → squash `82d7ef7`, CI 2/2 green
  on the merge, evaluator **PASS 87/100 zero blocking**).
- **Queue census (the loop's own progress metric, now measured rather than carried)**:
  `census: 41 closed / 89 rows (tagged-open 4, untagged 44) [LANDED 40, ROUTED 1, RULED OUT 0; PARKED 3, IN-SPRINT 1, NEXT 0]`
  read at the pre-record charter. 44 rows are UNTAGGED because their tag is not in the leading
  position; five of them (16, 18, 66, 68, 70) are in fact closed and are reported by number rather
  than guessed. Re-run it yourself:
  `bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79`
- **Milestones shipped**: M1 (semantic world library) · M2 (`ailang-worldd` daemon) — both complete.
- **Queue head**: rows **69** (fleet-gated), **73–78**, **81–89**, then **39**.
  Row **65** deferred per `D-WORLD-33`.

## Loop cadence and routing

- launchd `dev.ailang.mission-world`; kill switch `~/.ailang/state/mission-world.disabled` (armed).
- Bookkeeping issue **#129** (rotates weekly). Verify profile `ailang-code`; pinned binary
  `~/.pinned-ailang/ailang` **v0.30.0**.
- This iteration's roles: designer `claude:claude-fable-5-1` (rotation) · planner `opus` ·
  executor `pi:ollama/deepseek-v4-flash:0731-cloud` · evaluator `pi:ollama/minimax-m3:cloud`.
  Generator ≠ judge held on model **and** vendor.

## Waiting on Mark

**Nothing.** Decision ledger: **22 rows, ZERO OPEN**. This iteration asks for no decision.

## Quota posture

Metered **$0.24** of the $5 iteration ceiling — all of it design-quorum reviewers across two rounds
(`gpt5-6-sol` $0.169, `gemini-3-1-pro` $0.072). Both ollama lanes are flat-rate; fable and opus are
subscription buckets.

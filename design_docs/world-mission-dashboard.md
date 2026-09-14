# Mission Dashboard — Ailang World

_Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), the status archive, the log._

**As of**: iteration **175** — 2026-09-14 · repo `sunholo-data/ailang-world` · branch `dev` @ `1a0d2a7` (+ record)

## Where the mission is

- **Latest landing**: row **74** — `scripts/check_no_personal_email.sh`, a real personal-email gate
  over the loop-written surface (six mission docs + `scripts/`), every exclusion clause anchored,
  exit-2 floor, 10-arm suite, `mission_answer.sh` default → GitHub noreply, two CI steps
  (PR [#139](https://github.com/sunholo-data/ailang-world/pull/139) → squash `1a0d2a7`, CI 2/2
  green on the merge; evaluator r2 **PASS 95/100 zero blocking** after r1 FAIL 77 with one
  blocking finding — four exclusion clauses were substring matches; the fleet's copy still is).
- **The charter, log, archives, dashboard and index are now scanned by CI on every push** — write
  provenance verdicts, never addresses; spell address-shaped non-persons with ` AT `.
- **Queue census**: `census: 44 closed / 91 rows (tagged-open 4, untagged 43) [LANDED 43, ROUTED 1, RULED OUT 0; PARKED 3, IN-SPRINT 1, NEXT 0]`
  — moved by one closed and one row. Re-run:
  `bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79`
- **Queue head**: rows **69** (fleet-gated), **75–78**, **81–91**, then **39**. Row **65** deferred per `D-WORLD-33`.

## Loop cadence and routing

- launchd `dev.ailang.mission-world`; kill switch `~/.ailang/state/mission-world.disabled` (armed).
- Bookkeeping issue `#138` (rotated this morning; `~/.ailang/state/mission-world-gh-issue`). Verify
  profile `ailang-code`; pinned binary `~/.pinned-ailang/ailang` **v0.30.0**.
- Gate 0 runs TWO reads: `mission_directives.sh` (allowlisted humans) and
  `scripts/gate0_self_notices.sh` (the loop's own crash notices, no authority) — rc=0 this fire.
- This iteration's roles: designer `pi:ollama/deepseek-v4-flash:0731-cloud` (rotation entry
  `codex:gpt-6-astra` ration-blocked) · planner `opus` · executor `pi:ollama/deepseek-v4-flash:0731-cloud`
  · evaluator `sonnet` ×2. Generator ≠ judge on model and vendor.

## Waiting on Mark

**Nothing.** Decision ledger: **22 rows, ZERO OPEN**. Fleet proposal (Gate 5): anchor every clause of
the fleet's `check_no_personal_email.sh` exclusion list (a proposal, not an ask).

## Quota posture

Metered **$0.128** of the $5 ceiling this fire (four quorum reviewer bills). Quota: opus planner 168k
tok, sonnet judge 298k over two rounds, ollama designer + executor flat-rate.

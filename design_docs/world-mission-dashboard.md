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
- **Queue head — REGROOMED 2026-09-21 (Mark, attended); see the GROOMED block atop the charter Queue**:
  **8** (`SM.D`, clause-7, **ATTENDED-ONLY — surface to Mark, never run in-loop**) → **39**
  (session authority, clause-3, hard 1.0 blocker) → **92** (NEW: clause-5 value demonstration —
  what World is FOR) → **34/35/38** (workbench) → **40** → **27** → **93** (NEW: clause-4 floor
  run) → **75** → daemon robustness **22–26, 32**.
  **Maintenance, sweep-only, NOT iteration work**: 69, 76, 77, 78, 83, 84, 85, 86, 88, 89, 90, 91
  and 28, 29, 30, 33, 36, 37, 81, 87. Rows 76/77/78/84/89 belong to the SHARED harness in
  `sunholo-data/ailang` — route upstream, do not work here. Row **65** deferred per `D-WORLD-33`.
  Why: 27 of 34 open rows were clause-2 catch-all, 13 of them pure mission-harness, while the
  clause-5 value proof and the clause-4 floor run were **not rows at all**.

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

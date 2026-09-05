# Mission Dashboard — Ailang World

**Snapshot: 2026-09-05, after iteration 155.** Overwritten every iteration; history lives in
`design_docs/world-mission.md` (STATUS block), `world-mission-status-archive.md` and
`world-mission-log.md`.

## State
- `dev` == `origin/dev` == `d353ef1`. CI **GREEN 3/3** on that merge commit
  (`ailang-code verify gate`, `go host build + test gate`, `launchd drivers (bash 3.2)`).
- Verify gate = **two legs**: `AILANG_BIN=~/.pinned-ailang/ailang ./scripts/verify_ail.sh`
  and `go build ./... && go test ./... -count=1`. Pin is **AILANG v0.30.0 / `e37b370`**.
  `scripts/verify_go.sh` is **rc=1 at base** on a FLEET-OWNED driver-drift arm (row 76) — do not
  use it as a gate; use the two legs.
- Bookkeeping issue **#107** (prev #89). Decision ledger: **18 rows, ZERO OPEN**.

## Just landed (iteration 155)
- **Row 59 LANDED** and **row 50 LANDED as its consequence** — PR #116, rebase-merged so the three
  milestone commits survive. Evaluator `sonnet` **PASS 96/100**, zero blocking.
- The toolchain pins moved out of `run.sh`'s code into a data-only fixture with a bounded fail-loud
  parser; `coding-standards.md` S6 now carries the rule: *a criterion that greps for an assertion
  measures that somebody TYPED it; only a mutation measures that anybody RUNS it.*

## Next picks
1. **Row 60** — queue head.
2. **Row 61**, then **62–66**, **68–78**, then **39**.
3. **Rows 79 / 80** (Astra evidence-applicability, requirement-change vertical) — `[PARKED —
   DESIGN REVIEW]`, queued attended 2026-09-05. They need their own quorum and Mark's approval;
   by their own text they do not reorder existing release work.

## Parked on Mark
**Nothing.** The ledger has zero open decisions.

## Loop / routing
Controller `claude:claude-opus-5` · planner `opus` (`fail-closed:env-pin`) · executor
`codex:gpt-5.6-sol` · evaluator `sonnet` · designer = **rotation**, amended attended 2026-09-05 to
`codex:gpt-6-astra` → `pi:ollama/deepseek-v4-flash:0731-cloud` (astra takes the fable slot; fable is
now astra's fallback).

## Quota posture
Iteration 155 spent **`metered=$0.00`** of the $5 ceiling — every lane was a quota/subscription
bucket. No quorum round was owed.

## Known local gaps (not defects in this repo's code)
- `mission_directives.sh`, `mission-heartbeat.sh`, `resolve-role-spawn.sh` and `mission_pi_run.sh`
  are **absent here** and are reached by ABSOLUTE PATH into the V1 checkout (row 69).
- `tools/launchd/mission-control.sh.tmp.astra` is an untracked **fleet** artifact in this checkout
  (frozen core) — left alone deliberately.

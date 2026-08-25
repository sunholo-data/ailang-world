# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in the charter STATUS + `world-mission-log.md`.*

**Updated**: 2026-08-25 (iteration 125) · `dev` @ `2e44e3e` · CI **GREEN** (`checks=2`, both success)

## Latest release / state
- No release gate in flight. Item 14 `w-workbench-read-only` **COMPLETE** (11/11, iter-123).
- Queue head **row 5 `w-mcp-projection`** (clause 6) is in design — **not** sprintable as it stands.

## In-flight / next
- **NOW OWED (nothing blocked on a human): split #2 on `w-mcp-projection`** — carve `P6.B-A2A`
  out behind new queue row **39 `w-session-authority`**, re-quorum the `P6.T`/`P6.D`/`P6.V`
  remainder ONCE, then sprint **`P6.T`**.
- **`P6.T`** (toolchain floor `go1.25.6` → `go1.26.6`, ~0.1d, independently mergeable) has drawn
  **zero objections across all four quorum rounds** and is the first thing to land once the
  remainder clears.
- **Child doc `w-mcp-dispatch-projection.md`** (new, 173 lines) BLOCKED on upstream
  [`ailang#885`](https://github.com/sunholo-data/ailang/issues/885) — re-measured 2026-08-25: OPEN,
  0 comments, latest upstream release still `v0.33.2`. NOT quorum-cleared; quorum at pick time.
- **New row 39 `w-session-authority`** (~0.5–0.8d, needs a design doc, gated on nothing): the repo
  has **no inbound credential → session resolution at all** (`Bearer` 0, session-lookup functions
  0 across `host/`, same-scope control 181). Gates `P6.B-A2A` and nothing earlier.

## Loop cadence + routing
- Controller `claude:claude-opus-5` · planner `opus` · executor `codex:gpt-5.6-sol` → `opus`
  (deepseek SUSPENDED, `D-WORLD-20`) · evaluator `sonnet` (generator≠judge).
- **Designer rotation is stuck on one usable authoring lane.** Next entry `codex:gpt-5.6-sol`
  probed **rc=0 (healthy)** this iteration and was skipped on **judge-independence** — it is one
  of the two quorum reviewers. gemini is read-only under `CapRemoteSandbox` and cannot author.
  So Fable authors every doc and the pointer cannot advance. Loop cannot fix this; human can.

## Parked on Mark
- **Nothing.** Decision ledger: **13 rows, 0 OPEN** (`scripts/mission_decisions.sh --check`).
- `D-WORLD-26` **RESOLVED** this iteration — ARM A, `Authorization: Bearer` (Mark, attended, `#89`,
  2026-08-25T19:06:41Z, verbatim comment `A`).

## Quota / cost posture
- Metered this iteration: **$0.213298** (one quorum round, both reviewers present) against the
  $5 ceiling. Quota lanes: opus (controller), fable (designer, 1 bounded run).
- `w-mcp-projection` has cost **four quorum rounds** to find four surfaces of differing readiness
  — a scoping signal, surfaced for the human.

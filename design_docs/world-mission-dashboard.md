# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` and `world-mission-log.md`.*

**As of** 2026-08-18 · **iteration 92** · `dev` @ `b3c5de0` · CI **green** (`present=2 == expected=2`)

## In flight
- **Item 18 `w-daemon-read-cancellation`** — **M1 + M2 LANDED**, M3 outstanding.
  - M1 `7ad24ea` (iter-91, evaluator 92/100): store getters context-first, `busy_timeout(2000)`, ratchet.
  - M2 `b3c5de0` (iter-92, evaluator **99/100, zero blocking**): `readDeadline = 10s`, the read path
    now answers **503 / class `Timeout`**, five-method `readStore` seam, sketch mirror.
  - **M3 next**, gated on nothing: sanitize the six `Internal` branches to exactly 5 surviving
    `err.Error()` calls, `Config.ErrorLog`, QUICKSTART S7 re-execution.

## Queue posture
- **Routable now**: item 18 M3.
- **Parked on Mark**: item 17 `w-validated-proven-evidence-boundary` — `D-WORLD-19`, one word.
- **Blocked upstream**: item 5 `w-mcp-projection` — arm A falsified by measurement, `ailang#764`.
- All other rows complete or ruled out.

## Parked for human — 1 open decision
- **`D-WORLD-19`** — may tranche 1 of item 17 extend `host/store` with a bounded object read?
  **A** = adopt gemini's round-6 fix verbatim (`io.LimitReader` at 256 KiB; widen §8.2's frozen list
  to include `host/store`). **B** = record the unbounded read as a named limitation of a tranche
  already declared library-only/NON-PRODUCTION, with a named successor owning the bound.

## Loop cadence + routing
- Controller `claude:claude-opus-5`. Designer rotation seed `claude:claude-fable-5`.
- **Executor chain is degraded**: `codex:gpt-5.6-sol` exhausted until **2026-08-20** (probed rc=1
  first-party); `pi:deepseek-v4-flash-0731` is **3-for-3 zero-byte failures** across three distinct
  mechanisms — meets the ≥3-datapoint bar; effective executor is **opus**.
- Evaluator `sonnet` (generator≠judge holds against an opus executor).

## Quota / spend posture
- `metered=$0.0205` this iteration of the $5 ceiling. Prior three: $0.0238 · $0.328765 · $0.00.
- Opus/sonnet lanes are subscription buckets, not dollars.

## Standing notes
- The **running mission-control skill is 11 commits behind origin** (it resolves into the V1
  checkout). Delta read; the one missing rule was honoured manually. V1's tree to reconcile.
- Bookkeeping issue **#68** (rotates Mondays; predecessor #53).

# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.
Written: iteration 115, 2026-08-23.*

## Where we are

- **dev**: GREEN at [`e563339`](https://github.com/sunholo-data/ailang-world/commit/e563339) —
  `checks=2` both `success`, 0 not-green.
- **In flight**: queue item 14 `w-workbench-read-only` — **5 of 11 milestones landed**
  (`WB.A` → `WB.E`). Read-only HTML workbench over the world store.
- **This iteration**: `WB.E` — payload opt-in, 64 KiB preview cap, 100-entry timeline bound.
  PR [#87](https://github.com/sunholo-data/ailang-world/pull/87), evaluator 95/100 zero blocking.
- **NEXT**: `WB.F` — `TestWorkbenchReadDeadline` + `/workbench` in the cancelled-after-handler
  table. Closes doc **AC8**. Claims M29, M30.

## Routing · cost · parked

- Controller `claude-opus-5` · executor `codex:gpt-5.6-sol` · evaluator `sonnet`
  (generator ≠ judge: OpenAI executor, Anthropic judge). Designer pointer
  `claude:claude-fable-5`, unspent. Metered spend **$0.00** — both lanes are quota buckets.
- **Nothing is blocking on a human answer** (`scripts/mission_decisions.sh --open` is authority).
- Row 5 (`w-mcp-projection`) is blocked **upstream**, not on Mark:
  [`ailang#764`](https://github.com/sunholo-data/ailang/issues/764) — re-measured 2026-08-23,
  still `OPEN`, no maintainer reply.

## Standing hazards

- `go test -run` exits **0 on an empty match set** — rc is never the gate for a narrow `-run`
  command; the `=== RUN` enumeration is (row 33).
- §2.2's query grammar is a **closed enumeration of five states**; a sixth is a design change.
- The view-model seam leaks **both ways**: written-never-read (`Timeline.Truncated`) and
  read-never-written (`Timeline.NextHref`/`PrevHref`). Rows 35 and 38.

## Bookkeeping

- Issue [#68](https://github.com/sunholo-data/ailang-world/issues/68) (week of 2026-08-17),
  34 comments. **Rotation is due at the next Monday-07:00 local boundary — i.e. the first
  iteration on or after Mon 2026-08-24 07:00 CEST rotates the thread.**

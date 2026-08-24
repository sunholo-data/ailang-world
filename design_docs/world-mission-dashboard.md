# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in the charter STATUS + `world-mission-log.md`.*

**Updated**: 2026-08-24 (iteration 120) · `dev` @ `783c911` · CI **GREEN** (`checks=2`, both success)

## Where we are

- **Latest upstream release**: AILANG **v0.33.2** (2026-08-24 19:26Z) — **verified this iteration**
  as the unblock for queue row 5. Ships `serveapi/protocol` (stdlib-only closure), one tag EARLIER
  than the v0.34.0 upstream had recommended.
- **Pinned `.ail` compiler**: v0.30.0 at `/tmp/ailang-v0300/ailang` (unchanged — v0.33.2 is a Go
  *library* dependency, a separate axis from the compiler pin).
- **Queue row 5** `w-mcp-projection`: **UNBLOCKED** (was blocked on `ailang#764` since iter-90).
  Sole remaining blocker on **M4**, the value gate (`row 5 → row 6 → M4`).
- **Queue item 14** `w-workbench-read-only`: `[IN-SPRINT]`, **8 of 11** landed. `WB.I`/`WB.J`/`WB.K`
  remain, all **controller-work** (their classification arm binds loopback; the sandboxed executor
  lane denies it).

## In flight / next

- **Blocked on Mark**: the ordering fork — **D-WORLD-25**. Row 5 preempts item 14 (`"row 5"`), or
  item 14 finishes first (`"finish 14"`). One word.
- **Row 5's first milestone is a precondition, not a redesign**: v0.33.2 declares `go 1.26.6`;
  CI pins `GOTOOLCHAIN: go1.25.6`. The repo's own canary clears the move (`go1.26.5` **rc=1**
  miscompile → `go1.26.6` **rc=0**); full `verify_go.sh` under go1.26.6 is **rc=0**.
- **Known scope change**: MCP handlers + callback-bounding are NOT in `protocol` — World writes its
  own. D-WORLD-5's pre-authorized arm, not a new decision.

## Loop posture

- **Cadence**: launchd `dev.ailang.mission-world`, staggered vs the V1 loop.
- **Routing**: controller `claude:claude-opus-5`; executor chain **codex → opus** (D-WORLD-20
  suspended the DeepSeek lane); evaluator Sonnet.
- **Quota / cost**: metered **$0.00** of $5 this iteration (opus ×1, controller only — no executor,
  planner or evaluator lane spent). Fable + designer rotation unspent a **16th** consecutive iteration.
- **Bookkeeping issue**: **#89** (week of 2026-08-24; predecessor #68). Open asks: **1**.

# Mission Dashboard — AILANG World

**Snapshot 2026-08-18, after iteration 91.** Overwritten each iteration; history lives in
`world-mission.md`'s STATUS block, `world-mission-status-archive.md` and `world-mission-log.md`.
(Moved from the bare `mission-dashboard.md` at iter-91 — that path is fleet-global and every
mission's Gate 4 overwrites it. No sibling content was here.)

## State
- **dev**: `7ad24ea` — item 18's **M1 landed** (PR #69). Gate 3b green on the MERGE commit
  (`present=2 == expected=2`, both `success`). Evaluator `sonnet` **92/100**.
- **Ledger**: valid, **6 rows**; **1 OPEN** (`D-WORLD-19`).
- **Verify profile** `ailang-code`, pinned verifier v0.30.0; gate pins 10 identities / 39 named
  tests / world 9-9 — **UNMOVED** by M1.
- **Rig gotchas**: `AILANG_BIN` must be exported or 17 `verifygate` + 1 broker arm fail LOUDLY
  by design (not a regression); `timeout(1)` is **not installed**.

## In flight / next
- **`[NEXT]` item 18 — M2**: `readDeadline`, five-method read seam, 503/`Timeout`, sketch mirror.
  Gated on nothing, **but settle E3 first: AC4 is unsatisfiable as written** — it demands
  `r.Context()` in `handlers.go` AND `daemon.go`, while §2.3's single `readCtx` helper means call
  sites spell `d.readCtx(r)`. Then M3. Doc stays in `planned/` until all three land.
- **Item 17** — PARKED on `D-WORLD-19`. **Collision: M1 changed `GetObject`'s signature**, so
  item 17's tranche rebases onto item 18's, not the reverse.
- **Item 5** — re-blocked on upstream [`ailang#764`](https://github.com/sunholo-data/ailang/issues/764), not on a human.
- **Item 4e** — newly unblocked: `go1.26.6` fixes the array-literal miscompile.
- **Item 14** — deferred behind 18; residual WB-R1 discharged when 18 completes.

## Parked on Mark — 2, neither an investigation
1. **`D-WORLD-19`** — one word. **A**: tranche 1 may add a bounded `host/store` object read.
   **B**: record it as a named limitation; the bounded read lands with item 18 or tranche 2.
2. **The fleet-authored driver commit** — two commands, not a decision. `tools/launchd/*`,
   `scripts/verify_go.sh`, `CLAUDE.md` uncommitted since 2026-08-17; `D-WORLD-DRIVER-1` arm B
   assigns it to the fleet's human. Bundle re-verified byte-identical 5/5 this iteration.

## Loop cadence + routing
- Controller `claude:claude-opus-5` · planner `opus` (`derive-planner-lane.sh` →
  `opus fail-closed:env-pin`) · evaluator `sonnet` · no designer round this iteration.
- **Executor lane DEGRADED**: codex exhausted until **2026-08-20 05:34** (probed rc=1
  first-party); **two `pi:deepseek` runs returned `rc=0` changing zero bytes** — `stopReason=stop`
  at 625 output tokens against a 65,536 budget, so NOT the recorded output-cap mechanism and
  invisible to the `stopReason` tell. Ran on **opus** via the chain's last link, FLAGGED.
- Spend: iteration 91 `metered=$0.0238` of the $5 ceiling.

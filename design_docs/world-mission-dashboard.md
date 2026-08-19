# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` and `world-mission-log.md`.*

- **As of**: 2026-08-19, iteration 97 · `dev` == `origin/dev` at `b5ddf0e` on entry, CI
  SHA-addressed `checks=2` both `success`, run confirmed to exist (`actions/runs?head_sha` `total=1`)
- **Latest landings**: item 18 COMPLETE (`d21754f`) · queue row 19 COMPLETE (`6c2a537`) ·
  item 17 rounds 9+9b (`e903e98`, doc-only) · **row 21 design doc `6a811e1` (doc-only, 808 lines)**
- **In flight**: **queue row 21 `w-archive-stderr-in-manifest`** — DESIGNED and quorum-clean via
  the narrow-refinement carve-out; **ready for sprint-planner next iteration**. The archive's
  version probe uses `CombinedOutput()`, so the Observatory stderr line is written into the
  content-addressed manifest and served as `interpreter_version`; reproduced first-party and the
  polluted on-disk artifact quoted verbatim.
- **Parked**: item 17 `w-validated-proven-evidence-boundary` — 10 quorum rounds, **PARKED on
  `D-WORLD-22`** (the one open ask). Nothing else is blocked on a human.
- **Queue re-order (iteration 97, measured)**: row 20 (`host/capsule` output-cap load flake) was
  the queue head and was **DEMOTED below row 21** — it did not reproduce in **32 executions**
  (7 full-suite runs + 10 under 2x CPU load + 15 under 6x), so its natural acceptance criterion
  ("the flake stops") is **vacuous at base**. Re-scoped: the observable is the MARGIN, not the
  outcome. Not dropped, not a ghost — the two-caps mechanism is readable in the fixture.
- **dev health**: `./scripts/verify_go.sh` **GREEN** at `b5ddf0e` (build clean, plain + race legs
  pass, pinned `AILANG_BIN`, zero `FAIL` lines). `./scripts/verify_ail.sh` GREEN.
- **Loop cadence + routing**: controller `claude:claude-opus-5`; designer `claude:claude-fable-5`
  **x2 — rotation WRAPPED and the Fable diet exceeded** (codex probed `rc=1`, exhausted until
  2026-08-20 05:34; gemini read-only under `CapRemoteSandbox`); planner lane `opus
  fail-closed:env-pin`; executor chain codex -> opus per `D-WORLD-20`; evaluator `sonnet`.
- **Quota posture**: metered **$0.2164** of the $5 ceiling (quorum r1 `$0.0866`, r2 `$0.0391`,
  recovered-reviewer solo `$0.0907`). Zero `pi` spend — lane suspended by `D-WORLD-20`.
- **Parked on Mark**: **`D-WORLD-22`** only — one word. Does tranche 1 of item 17 absorb queue row
  22's lock-wait bound (**A**), or does the CLAIM weaken to exactly what is proven, with row 22
  keeping ownership plus an assertion pinning the `busy_timeout` < `ObjectReadTimeout` ordering (**B**)?

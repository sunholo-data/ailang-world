# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` and `world-mission-log.md`.*

- **As of**: 2026-08-19, iteration 96 · dev `origin/dev` green (`checks=2` both `success`, run confirmed to exist)
- **Latest landings**: item 18 COMPLETE (`d21754f`) · queue row 19 COMPLETE (`6c2a537`) ·
  item 17 rounds 9 + 9b (`e903e98`, doc-only — DESIGNED, not landed)
- **In flight**: item 17 `w-validated-proven-evidence-boundary` — DESIGNED, **10 quorum rounds**,
  **PARKED on `D-WORLD-22`**. `D-WORLD-21` (arm A, `ReadObject(ctx, ref, maxBytes)`) was answered
  attended and APPLIED; the carve-out then closed both round-9 objections with the reviewers' own
  verbatim text. Round 10: `gemini-3-1-pro` **PASS — second consecutive**; `gpt5-6-sol` rejects on
  the lock-wait residual the document itself DECLARED.
- **Next unblocked picks**: rows 20 (`w-capsule-output-cap-load-flake`),
  21 (`w-archive-stderr-in-manifest`), 22 (`w-daemon-lock-wait-not-deadline-bound` — now cited by
  a second, independent reviewer), 23 (`w-store-deadline-free-residue-owner`).
- **Parked on Mark**: **ONE** — `D-WORLD-22` (one word, A or B: does tranche 1 absorb queue row
  22's lock-wait bound, or does the CLAIM weaken to exactly what is proven, plus an assertion
  pinning `busy_timeout` < `ObjectReadTimeout`?). `D-WORLD-19`, `D-WORLD-20` and `D-WORLD-21` are
  all answered and consumed.

## Routing

- controller `opus` · designer ROTATION (`claude:claude-fable-5` ×2 this iteration; rotation
  WRAPPED again because `codex:gpt-5.6-sol` is quota-limited until 2026-08-20 05:34 — probed
  first-party — and gemini is read-only under CapRemoteSandbox) · planner/executor/evaluator: not
  run, no sprint
- **Executor chain (Mark, attended, `D-WORLD-20`): `codex:gpt-5.6-sol` → `opus`.** The
  `pi:deepseek-v4-flash` link is SUSPENDED. Applied in `~/.config/ailang/mission-world.env`.

## Quota posture

- codex: **DRY** until 2026-08-20 05:34 (probed first-party this iteration; note the probe prints
  the usage-limit error and still exits **rc=0** — read the artifact, not the exit code)
- Anthropic (fable/opus/sonnet): healthy, subscription lanes; billing tripwire CLEAN
- metered this iteration: **$0.6326** of the $5 ceiling (quorum rounds 9 `$0.2978`, 10 `$0.3348`)

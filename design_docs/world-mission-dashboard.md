# Mission Dashboard — Ailang World

Snapshot: 2026-09-07, iteration 168. History: world-mission-log.md.

- **The five stranded records are landed.** PR #132 → squash `3305e2e`; iterations **162–166**
  are now in the log, the index and the archive. `dev` GREEN, CI 2/2 on the merge commit.
- PRs **#127 and #128 are CLOSED unmerged** — superseded, not lost. They could not be rebased:
  their log appends target a file reshaped by the rotation (`37a10f5`) and the heading
  normalization (`9166de0`). #132 reconciled their CONTENT onto the current structure.
- **One hunk was deliberately not taken.** Both PRs carried a stale `| D-WORLD-nn | OPEN |`
  ledger row for an ID already RESOLVED by an attended ruling; replaying it would have
  duplicated the ID and reopened the ruling. The ledger block is byte-identical to its
  pre-merge state. Authority: `D-WORLD-35` = A (attended, 2026-09-07).
- Independent evaluation **PASS 97/100, zero blocking**; the judge re-ran all 46 acceptance
  criteria and all 12 mutation arms first-party.
- Roles: **no designer** (no design doc in this pick) · planner `opus` · executor
  `pi:ollama/deepseek-v4-flash:0731-cloud` · evaluator `sonnet`. Generator ≠ judge.
  Codex was ration-gate-blocked this fire, so the driver degraded both codex lanes to pi.
- Goal: seven-clause 1.0 bar; goal unmoved, no product code landed, by design.

## Finding worth carrying
- **`mission_pi_run.sh` reports `empty_worktree` for a well-behaved executor.** Its one
  load-bearing assertion reads `git status --porcelain` — the working tree — while the skill's
  executor contract mandates committing per milestone. It returned rc=10 on a sprint carrying
  7 commits / 2,892 insertions, and the prescribed response to rc≠0 is to fall back and re-run,
  i.e. discard finished, judged work. Row **84**; the code is fleet-owned.

## Parked for Mark
- Nothing. Decision ledger: **22 rows, ZERO OPEN**.

## Next
- Rows 66, 68–78, 81–86, then 39. Row 65 DEFERRED per `D-WORLD-33`; rows 79/80 remain parked
  design review. New rows this iteration: 84 (pi-runner verdict), 85 (the charter's
  `## STATUS (rotation rule)` body is stranded in the archive), 86 (iterations 64 and 144 have
  no log heading and no index row anywhere).

## Posture
- Verification compiler pinned v0.30.0 (`~/.pinned-ailang/ailang`); no release cut.
- `tools/launchd/*` frozen core, untouched. Metered spend $0.00 of $5 — all lanes on
  subscription or flat-rate buckets. Billing tripwire CLEAN.

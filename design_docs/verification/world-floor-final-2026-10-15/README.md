# Row 93 FINAL (2026-10-15, attended): verdict HOLDS

Pre-registered in `preregistration.json` (`4e8299c`, digest `f46aae32…`). Authority: D-WORLD-70/72/74, the
2026-10-15 known-deltas smoke (`world-floor-smoke-2026-10-15`, C1–C3 HOLD) and Mark's attended instruction to
run it. Protocol (§4.11 as amended by D-WORLD-74): Phase 1 shell ×3, eligibility sealed and committed
(`9b6f6a8`, tag `row93-final3-eligibility-sealed`, pushed before any World run), then Phase 2 paired ×3, where
each task runs shell and World back to back in counterbalanced order. Δ and O come from the paired rows only.
Verdict: `report.py verdict` → `verdict.json`, `summary.md`.

| agent | eligible | Δ (bar ≥ −1/50) | O (bar ≤ 1/4) | agent verdict |
|---|---|---|---|---|
| claude (2.1.300, pinned) | yes (69/69) | 0 (69/69 vs 69/69) | −17.9% | passes |
| codex (gpt-6.1-sol, codex-cli 0.162.0) | yes (69/69) | 0 (69/69 vs 69/69) | −5.7% | passes |

**Overall: HOLDS.** Every one of the 414 graded cells passed (Phase 1: 2 agents × 23 tasks × 3; Phase 2:
2 agents × 2 arms × 23 × 3). Harness faults 0, capacity re-runs 0, aborts 0 in all six phase records; all 18
canaries passed. Drift (paired shell vs Phase 1 shell) +8.4% claude, −3.5% codex, both within 10%.

## What changed vs the 2026-10-14 FINAL (which FAILED on codex Δ)

Exactly the fixes listed in the preregistration's `known_issues_from_the_2026_10_14_final`, nothing else:

- **World wording `row93-direct`** (was `row135-rw`). Codex World input tokens per task, median: 127,950 vs
  shell 128,821 = **0.99×** (2026-10-14: 2.03×). Codex O went from +24.5% to −5.7%.
- **Corpus `77c9e97`** (sunholo-data/ailang#1837): `config_file_parser` says the program itself must write
  `app_config.json`. It passed 6/6 for codex in Phase 2 (2026-10-14: World 0/3, shell 2/3).
- **Client pin:** Claude Code 2.1.300 copied to the rig's `clibin/claude`, first on PATH, with
  `DISABLE_AUTOUPDATER=1`, for both phases. Each of the six runs logged `claude --version` = 2.1.300 at its start and passed the
  AC5.1 digest check, which covers `cli_versions`.
- Paired counterbalanced order (D-WORLD-74), as on 2026-10-14.

## Notes

- **Spend.** Claude $34.89 (Phase 1 $14.88, Phase 2 $20.01). Codex on the existing subscription session.
- **Times.** Phase 1 15:09–15:56, Phase 2 15:56–17:23 CEST.
- **Rig.** `~/.ailang/state/floor-final-2026-10-15` (fresh `final-prepare` build, se-tools registry revision 1,
  six sessions minted by Mark). Transcripts stay on the rig and are hashed per row (`transcript_sha256`).
- Logs: `phase1.log`, `phase2.log`.

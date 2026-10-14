# Row 93 FINAL (2026-10-14, attended): verdict FAILS

Pre-registered in `preregistration.json` (`fd321b6`, digest `d079339e…`; authority D-WORLD-70/72/74). Protocol
(§4.11 as amended by D-WORLD-74): Phase 1 shell ×3, eligibility sealed and committed (`4f15e17`), then Phase 2
paired ×3, where each task runs shell and World back to back (World first iff task index + run is odd). Δ and O
come from the paired rows only. Verdict: `report.py verdict` → `verdict.json`, `summary.md`.

| agent | eligible | Δ (bar ≥ −1/50) | O (bar ≤ 1/4) | agent verdict |
|---|---|---|---|---|
| claude (2.1.299, pinned) | yes | 0 (69/69 vs 69/69) | −21.6% | passes |
| codex (gpt-6.1-sol) | yes | **−2/69 = −2.90%** (68/69 vs 66/69) | +24.5% | **fails Δ** |

**Overall: FAILS** (every agent must pass). Harness faults 0, capacity re-runs 0, unclassified 0 in all six
phase records. Drift (paired shell vs Phase 1 shell) is +0.7% for claude and −0.6% for codex, so it is not
confounded. Interleaving removed the arm-order artifact: the pilots' "+6–8 s World overhead" did not recur as
a gate failure, and codex O passed.

## The Δ failure is one task

Codex `config_file_parser`: shell 2/3, World 0/3. Every other task is 3/3 in both arms. One fewer World fail
would give Δ = −1/69 = −1.45%, which passes.

The task says "Create a JSON config file named app_config.json … implement a function to load and validate
it". The grader runs only `benchmark/solution.ail` in a clean directory, so a passing program must write the
file itself. All four failing rows (World r1/r2/r3, shell r3) wrote `app_config.json` as a workspace file
(World via `ailang-write`, shell via a file change) and had the program only read it. At grade time it is
absent ("Configuration file not found" / "open app_config.json: no such file"). Passing rows called
`writeFile` inside the program. Claude wrote it in-program in all 9 of its runs of this task.

This is a real agent behaviour under the preregistered rule, not a harness fault: the verdict stands as FAILS.
That World's first-class write tool makes the workspace-file reading more likely is plausible
(World 3/3 vs shell 1/3 here, 0/3 in Phase 1) but N is small. Any remedy changes the protocol or the
corpus and needs an attended ruling before another FINAL.

## Notes on this run

- **Claude Code pin.** Claude Code auto-updated to 2.1.300 at 20:25, between the seal and Phase 2. Phase 1's
  last Claude row ended 20:16 on 2.1.299. Phase 2 ran a copy of 2.1.299 at the rig's `clibin/claude`, first on
  PATH, with `DISABLE_AUTOUPDATER=1`, so AC5.1 (`cli_versions` in the digest) matched. Every phase record names
  commit `4f15e17`. The operator's own install was not modified.
- **Seal tag.** The seal commit was first made as `fcd40a7`, tagged `row93-final2-eligibility-sealed` and
  pushed before any World run. It was then rebased over two loop landings (#260, #261) to `4f15e17`; the
  evidence is byte-identical.
- **Spend.** Claude $36.86 (Phase 1 $15.79, Phase 2 $21.08). Codex on the existing subscription session.
- Logs: `phase1.log`, `phase2.log`. Transcripts are on the rig (`~/.ailang/state/floor-final-2026-10-14/transcripts/`)
  and hashed per row (`transcript_sha256`).

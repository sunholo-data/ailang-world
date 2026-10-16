# Row 114 — unseeded, baselined provenance walks (clause 5)

Protocol: `design_docs/planned/w-provenance-teeth-unseeded-baselined-walk.md` r4, frozen by D-WORLD-76 = A (attended 2026-10-16). Plan: `…-sprint-plan.md`. Run attended 2026-10-16 (Mark attending); the controller is the attended session, which had seen the Q1 diagnosis and was neither B, W nor V.

## Harvest log

| id | question | cohort | status |
|---|---|---|---|
| Q1 | Why did codex's World arm fail `config_file_parser` 0/3 in the 2026-10-14 FINAL, while its shell arm passed 2/3? | K1 | **attempted 2026-10-16 — B UNVERIFIED, W UNVERIFIED** (`Q1-config-file-parser/`) |
| Q2 | Why was codex's World overhead +26.6% in the 2026-10-09 FINAL? | K3 | candidate (pre-diagnosis commit not yet measured) |
| C1/C2/C3 | iter-206/207/210 CI and suite incidents | none | rejected NOT-RUN(NO-CORPUS); C1 also contaminated |

## Q1 accounting

| | B (baseline) | W (World walk) |
|---|---|---|
| channel | pi, `ailang_read` (+list) / `ailang_check`; B root = `git archive fd321b6` + raw rig logs | Claude Code 2.1.300 `sonnet`, Bash limited to 7 World client verbs on a served copy of K1 |
| answer time | 99.0 s, 44 tool calls | 140.1 s, 55 tool calls |
| envelope (issue → verified) | 128.9 s (no threshold) | **147.9 s** (bar ≤ 300 s) |
| verifier | UNVERIFIED | UNVERIFIED |
| spend | $0.68 (OpenRouter) + V $0.06 | $0.52 + V $0.03 |

Integrity: K1 and original sha256 `bf1d4db5…` unchanged; served `head` `sha256:1710ddd7…` and 495 entries unchanged; W 0 non-allowlisted calls, 0 denials; question sha256 `61d44f24…` in both prompts; initial-context scans 0 hits; B root scan 0 hits over 4866 files with a positive control (2 hits). No ratio is stated: neither arm verified (AC10).

Deviations from the plan, all recorded in `pins` before reveal:
- **B binary v0.54.11** (attended ruling). policy-tool in v0.41.0 and v0.52.1 refuses directory listing ("is a directory"); v0.54.11 lists, with the same escape and `run` refusals.
- **B provider OpenRouter** (`z-ai/glm-5.3`, the plan's pinned fallback). zai was out of balance (429 code 1113, 3/3 attempts); OpenRouter worked after the operator raised its key limit.
- **Claude Code 2.1.300** from a pinned copy, `DISABLE_AUTOUPDATER=1` (the ambient install had moved to 2.1.301).
- pi needs stdin closed (`< /dev/null`); otherwise it waits silently.

## Findings

**F-1 (World capture gap — why the walk failed).** W found, in two of three codex episodes, that codex wrote `app_config.json` with `ailang-write`, then a solution that reads it, and that `ailang-run` printed the expected line. It could not say why grading failed: the store records World tool calls only, so it holds no grader verdict and nothing about how grading runs (only `benchmark/solution.ail`, in an empty directory). It also did not find the r2 episode's attempt. The facts of the mechanism are in the store; the context that makes them a cause is not. B, holding the repo (with `scripts/floor/grade.py`) and raw logs, did connect them. Candidate follow-up (the one scoped row §11 allows): record the grading outcome and environment next to an episode's effects, or let a walk declare an external evidence source.

**F-2 (controller defect in E2).** The plan's E2 (b) covers the codex World solutions only; the rubric, written by this controller, asks for the causal claim for both arms. V_B therefore confirmed B's World-arm explanation (re-execution: "not found" in a clean dir, the expected line with the run's `app_config.json`) but marked the shell half unsupported. E2 was not changed after reveal, so both verifiers judged against the same bundle. A later attempt needs E2 to cover every arm the question names (here also the codex shell solutions).

**F-3 (baseline strength).** With read and list only (no grep), the baseline answered in 99 s with a World-arm mechanism its verifier confirmed. The design's expectation that a no-grep baseline is slow (§12) did not hold for this question: the raw transcripts in B's root carry the same writes the store does.

Row 114 stays open: 0 of ≥3 qualifying successes.

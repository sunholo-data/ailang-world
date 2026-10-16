# Row 114 — unseeded, baselined provenance walks (clause 5)

Protocol: `design_docs/planned/w-provenance-teeth-unseeded-baselined-walk.md` r4, frozen by D-WORLD-76 = A (attended 2026-10-16). Plan: `…-sprint-plan.md`. Run attended 2026-10-16 (Mark attending); the controller is the attended session, which had seen the Q1 diagnosis and was neither B, W nor V.

## Harvest log

| id | question | cohort | status |
|---|---|---|---|
| Q1 | Why did codex's World arm fail `config_file_parser` 0/3 in the 2026-10-14 FINAL, while its shell arm passed 2/3? | K1 | **attempted 2026-10-16 — B UNVERIFIED, W UNVERIFIED** (`Q1-config-file-parser/`) |
| Q2 | Why was codex's World overhead +26.6% in the 2026-10-09 FINAL? | K3 | candidate (pre-diagnosis commit not yet measured) |
| Q3 | In trial episode tps-a (platform-source, row 183), why does the worktree hold 47 untracked files — incl. zz_ls.ail, zz_slice.ail and .ailang/ — when the task asked for 4 test files + 1 doc edit? | K4 (trial store, sha256 fddd2e59…) | **attempted 2026-10-16 — B VERIFIED (88.7 s), W INVALID(seeding) (UNVERIFIED on its answer; 66.8 s)** (`Q3-tps-a-untracked/`) |
| Q4 | In trial episode tps-b, why did the ADK 2.12 grader exit rc=2 with collection errors while the sub-doc says every 2.12 failure of the new tests is runtime, not collection? | K5 (trial store, sha256 9d72ad32…) | **attempted 2026-10-16 — B VERIFIED (67.0 s), W PASS (54.6 s)** (`Q4-tps-b-212-collection/`) |
| Q5 | In trial episode w185, why did the change modify pre-existing test code rather than only add tests — which changed, and was any assertion weakened? | K5 | **attempted 2026-10-16 — B VERIFIED (239.9 s), W UNVERIFIED (198.9 s)** (`Q5-w185-existing-tests/`) |
| Q6 | Why does tps-b's generated PR body give its World entries as 89–154 while its trailer gives 96–154 — what are 89–95? | K5 | **attempted 2026-10-16 — B VERIFIED (87.5 s), W PASS (25.3 s)** (`Q6-tps-b-entry-range/`) |
| Q7 | In trial episode w185, why does the change derive an episode's degraded status from its directory each read instead of recording it on the session as row 185 says? | K5 | **attempted 2026-10-16 — B VERIFIED (92.6 s), W UNVERIFIED (94.7 s)** (`Q7-w185-degraded-derivation/`) |
| Q8 | In trial episode tps-b, why does the sub-doc say ResilientLlm's retry loop catches Exception and the BaseException handler the plan cites is in first-token-timeout cleanup? | K5 | **attempted 2026-10-16 — B VERIFIED (116.8 s), W PASS (39.0 s)** (`Q8-tps-b-baseexception-plan/`) |
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

## Q3 findings (2026-10-16)

**F-4 (World capture gap: tool side effects).** `ailang-run`/`ailang-check` wrote 41 files of AILANG compile cache into the episode worktree (`.ailang/cache/compile/`), but the store records only the tool INPUTS and outputs, not those filesystem writes; explicit `Workspace.Write` effects are recorded, the toolchain's own writes are not. The walk therefore could not say where `.ailang/` came from, while the baseline read the files and their mtimes directly. Candidate row: record a worktree change summary (paths + hashes created/modified) per tool invocation, or per episode at its end.

**F-5 (W channel confinement hole).** Claude Code 2.1.300 persisted an oversized World output to `~/.claude/projects/<cwd>/…/tool-results/<id>.txt` and, in `--permission-mode dontAsk` with a verb allowlist, still auto-allowed a read-only `grep … | head` on that file chained with an allowed verb. The M1 controls did not cover this case. Fix before the next attempt: add explicit denies (`--disallowedTools` for `Bash(grep:*)`, `Bash(head:*)`, `Bash(cat:*)`, `Read`) and a control that tries exactly this; or keep outputs small enough not to be persisted.

**F-6 (baseline).** With read/list only and the raw transcript in its root, the baseline answered the full cause in 66.6 s, verified in 88.7 s — again fast where the evidence is in the logs.

## K5 (Q4–Q6, 2026-10-16)

Cohort K5 = the trial store after row 183 episodes tps-b (codex, platform-source) and w185 (Claude Code, ailang-world — row 185 built inside World), which ran concurrently on one daemon. Questions harvested while preparing their draft PRs (platform-source#71, ailang-world#291); questions, E2 bundles, rubrics, preambles and pins sealed before any arm ran (`pins`, `exposure-log` in each Q dir). B's root: both final worktrees + both base commits (no `.git`, no venvs) + raw logs, 11950 files; diagnosis-phrase scan with a positive control (the PR bodies: 3 + 3 hits) found 52 generic hits, none mentioning either episode, all identical in `base/` (reviewed). W: the A1 deny list on every run; the pre-reveal control denied 8/8 forbidden commands including F-5's case.

| | Q4 B | Q4 W | Q5 B | Q5 W | Q6 B | Q6 W |
|---|---|---|---|---|---|---|
| answer time | 53.9 s | 40.1 s | 230.6 s | 186.7 s | 70.9 s | 14.7 s |
| envelope (issue → verified) | 67.0 s | **54.6 s** | 239.9 s | 198.9 s | 87.5 s | **25.3 s** |
| verifier | VERIFIED | VERIFIED | VERIFIED | UNVERIFIED | VERIFIED | VERIFIED |
| tool calls | 20 | 17 | 35 | 114 | 12 | 7 |
| spend | $0.34 | ~$0.21 (sub) | $1.59 | ~$0.89 (sub) | $0.33 | ~$0.06 (sub) |

Integrity, every question: K5 archival sha256 unchanged; served head `sha256:abe8303c…` unchanged; question bytes identical in both prompts; pins verify; 0 non-allowlisted W calls succeeded.

**F-7 (`why` hides which file an edit touched).** Q5's walk ran `why 198` on the `ailang-edit` of `host/worktree/worktree_test.go`, but `why` prints a truncated input preview and canonical JSON orders `new_text` before `path`, so the path never appears; the walk then said only one pre-existing test file was edited. The record is complete; the display drops its key field. Candidate fix (small): `why` prints `path` (and other short identifying fields) ahead of the truncated preview for workspace and AILANG tools.

**F-8 (baseline).** The read/list-only baseline verified on all three K5 questions, as on Q3. The walk was faster on every K5 question (by 12 s, 41 s and 62 s), and on Q6 — a question about World's own entry numbers — 3.5× faster; but the baseline is not failing, so the comparison shows speed, not answers only World can give.

## K5 second harvest (Q7, Q8, 2026-10-16)

Harvested after Q4–Q6 were revealed, from the same two PRs; sealed before any arm (`pins` per Q dir; exposure appended to the K5 `exposure-log`). Same B root (re-scanned for the new diagnosis phrases: control 3 hits, root 0), same W allowlist and deny list. One change, pinned: W's client is the F-7 build (`9eaf83d`, PR #292), so `why` shows an elided input's short fields such as `path`.

| | Q7 B | Q7 W | Q8 B | Q8 W |
|---|---|---|---|---|
| answer time | 83.9 s | 86.0 s | 107.9 s | 29.8 s |
| envelope | 92.6 s | 94.7 s | 116.8 s | **39.0 s** |
| verifier | VERIFIED | UNVERIFIED | VERIFIED | VERIFIED |
| tool calls | 24 | 47 | 16 | 13 |
| spend | $1.06 | ~$0.36 (sub) | $0.37 | ~$0.13 (sub) |

**F-9 (World does not record the agent's own reasoning).** Q7's answer — the agent avoided a store schema bump — is in the store only as a code comment inside an `ailang-edit`'s `new_text`; the walk fetched 11 input payloads, not that one, and answered "the records do not state the reason". The baseline read it from the agent's transcript and final summary, which World never records. Candidate row: commit the agent's final summary (the trial runner already captures it) as an object linked to the episode's outcome, so a walk can find the stated rationale.

## Tally

Row 114: **3 of ≥3 qualifying successes — Q4 (54.6 s), Q6 (25.3 s), Q8 (39.0 s)**, each with a same-question live baseline that ALSO verified (67.0 s, 87.5 s, 116.8 s). Attempts that did not qualify: Q1, Q3 (INVALID seeding), Q5, Q7 — all preserved. Clause 5's closure is Mark's ruling (D-WORLD-80).

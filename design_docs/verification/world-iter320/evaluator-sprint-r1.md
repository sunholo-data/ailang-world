# Iteration 320 — row 234 sprint, independent evaluator round 1

- Judge: Agent tool, model `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), fresh foreground context, own worktree `.wt-world-iter320-eval` at `ebc220e0` (restored byte-identical, `git status` empty at the end).
- Reviewed: `785b3bd1..ebc220e0` (design `5407c160`, plan `391ef1aa`, M1 `63fd4f81`, M2 `ebc220e0`).
- Harness usage: 45,248 tokens, 14 tool uses, 657 s.

## Verdict: PASS 91/100, zero blocking

| Criterion | Score |
|---|---|
| Correctness of the fix | 30/30 |
| Guard purpose + mutation evidence | 21/25 |
| Test non-vacuity | 19/20 |
| Scope / evidence discipline | 14/15 |
| Code quality | 7/10 |

## Verified by the judge

- Under a `GOTMPDIR` root containing `x134x` and a control root, the six `requireF15` callers, T1, T2 and `TestExecToolExit134ChecksNode` pass; T2's own precondition proves the root was honoured.
- Reverting only the guard to `strings.Contains(out, "134")` reds 6 callers under the 134 root, incl. both CI signatures; T2 reds under the control root too (it forces its own root).
- A launched reply is always a marshalled `execResult` carrying `exit_code` (`handlers_exec.go`), so the strict two-key decode rejects it; the post-launch 134 refusal (`handlers_exec.go:326`) moves the spawn counter.
- Plan mutants M1, M3, M7 each red their named rows. Own mutants: X2 (`after < before`), X3, X4 killed; **X1 (missing `ok` accepted) and X5 (missing `refused` accepted) survived**.
- Only `host/broker/exec_tools_test.go` changed outside `design_docs/`; gofmt and vet clean.

## Findings and controller follow-through

- **NB1 (X1/X5 survive).** Controller reproduced both survivals (T1 rc 0 with each mutant), added T1 rows `ok-missing` `{"refused":"…"}` and `refused-missing` `{"ok":false}` (commit `227b494b`), then re-applied each mutant: X1 → only `ok-missing` red, X5 → only `refused-missing` red; file restored byte-identical (`cmp`).
- **NB2 (race leg rc 1 in the judge's shell).** `TestEpisodeLiveReplayThreeArmsAndEvidence` needs `AILANG_BIN` set to the pinned binary; with it set the plain run passed. Environment, not the change; the executor's race run (with the repo's env) was rc 0, 0 DATA RACE, and CI runs it with the pin.
- **NB3 (sibling-sweep wording).** The design's §7 sentence is already scoped to literal-token `strings.*` absence checks and leaves the variable-token class open as R3. No change.

## Round 2

Not run: the one code change adds two table rows proven to kill the two surviving mutants; nothing else scored changed.

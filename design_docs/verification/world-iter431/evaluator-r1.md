# Evaluator report r1 — row 243 V4 S5 (6967f0c2..2ddfb651, base cd06d950)

## Verdict: PASS — 86/100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design/plan §9 | 26/30 |
| Tests | 24/30 |
| Security / fail-closed | 18/20 |
| Coding standards | 9/10 |
| Scope and size | 9/10 |

## Measurements
- Size unit `git diff -M --numstat cd06d950 HEAD -- cmd host docs app` = **743** (cap 1,400). Product files match plan §5; `.ail`, app.js, daemon.go, surface files untouched. Controller `callers.txt`: provisionEpisode callers = 5 (start, preset, resume pass held lock; approve and session new pass nil); `resolvePresetStart` has exactly one caller (terminal preset start, episode_preset.go:328); no production caller of Claim/Reserve/ObserveOperator*; no new route.
- Gate sweep (controller): all legs rc=0, `FAILS: none`, observation-before == observation-after (diff clean).
- My runs: `go test -count=1 ./host/store ./host/authority` rc=0; targeted cmd run (Operator|Episode|Preset|Launch|Finish|Session|Approve|Quickstart|Help|CredentialWriters) ok 69 s; `-race -run OperatorProvision|PresetStartRead|OperatorTerminalExclusion` ok 11 s; `go vet ./cmd/... ./host/...` rc=0; `go test -run '^$' ./...` clean.
- CI cost: new tests `OperatorProvision|PresetStartRead` = 6.3 s wall (PresetStartReadRefusals 4.6 s). No real sleeps found in the new tests. Whole package in the controller sweep was 235 s without -race.

## M1 — exclusion
Enumerated callers (grep): `provisionEpisode` at approve.go:398, episode.go:377, episode_preset.go:365, episode_resume.go:215, session_new.go:660. Order inside it: terminal fence (skipped only for signed approve) -> `terminalEpisodeLock` (non-waiting, durable owner read, exact action id in message, `exitUsage`) -> prompt -> writer/worktree/mint. Nothing is consumed before the refusal: in approve, the signature/assertion checks, daemon health and re-derivation run before, but the request consume, approver use and passkey counter all happen inside the mint transaction after the lock. The test asserts request row and approver rows are DeepEqual after refusal. `held != nil` returns before any second open, so there is no flock self-collision (J-HELD-INVERT, J-*-PASS-NIL killed). `defer lock.Close()` is on every return path (J-NOCLOSE killed by approve_releases/session_new reacquire). A durable owner still refuses after descriptor close with an expired lease (test closes the synthetic holder first). A declined attended terminal leaving the private sidecar is acceptable (stable inode, documented in QUICKSTART).
Gap (non-blocking, carried): `session mint` (session.go:299, `authority.Mint`) still mints an episode credential with no episode lock. It adds no worktree, S5 dispatches nothing, and Appendix C names it ("direct terminal session mint/revoke are outside provisionEpisode"), but it is not listed as a hard S5b/S6 activation gate with a test.

## M2 — read seam
Diff read in full: prompt, refusal strings, ordering (preset load -> launch preflight -> episode grammar -> overrides -> Git -> collisions -> task read once -> prepareMint -> creds version -> health -> THEN lock -> provision) preserved; lock still acquired after all reads. Task read exactly once via callback; the parity test rewrites the file after resolution (inside openTerminal) and asserts the original bytes reach copy/prompt/start row. Resolver opens no writer/lock (test asserts no sidecar, DeepEqual observeProvision before/after). Only bounded probes (git, launch binary lookup). No new context root, no sleeps. Minor: `fs.Lookup("hosts-dir")`/`launch`/`headless` nil checks dropped (plan notes the flags are always defined).
Weakness: parity tests compare resolver output against terminal flow of the same new code and fixed expectations, not a captured pre-refactor transcript. Pre-existing preset tests (unchanged) are the real regression control, and they pass. Two mutants survive the whole package (below).

## M3 — replay
Byte comparison precedes the `ApprovedBy != ""` early return (operator_cas.go:58), after the digest/id conflict check. Exact replay still returns the stored row (test DeepEquals it) and the changed-bytes case asserts `ErrOperatorConflict` and an unchanged row (J-REPLAY-NOCHECK/NILERR/WRONGERR killed by TestOperatorConfirmReplay in host/authority; the isolated host/store package does not kill them — plan discloses this).

## S5/S5b split
Honest: 743 + 689 calibrated = 1,432 > 1,400. Appendix C carries the atomic mint/worktree/file recovery rows (§3.2 maps every §5.3 test to S5b/S6/S8/S12), the attended resolver, and the paired-process lock consumer/drain gate. Nothing in §5.3's S5 share is silently dropped; `session mint` is only mentioned in passing.

## Linux portability
Read the new tests: ELOOP asserted via `syscall.ELOOP.Error()` (same text both OSes), flock via independent open descriptions, mode 0750 chmod for nonprivate sessions (fails as non-root only; fine on CI), "no such file" substring matches Linux text. No macOS-only assertion found.

## Pre-existing edit audit
approve.go/episode.go/episode_resume.go/session_new.go: one-token signature changes. episode_preset.go: refactor (above). operator_store_test.go: pure additions (replay mismatch + unchanged row). operator_cas.go: +3 lines. No assertion weakened.

## Mutation ledger: `$HOME/.ailang/state/world-iter431/judge/mutants.tsv` (46 lines; scripts in judge/bin/mut.py, per-run logs runs-*.log)
Own mutants not in plan ids (`comm -23 judge/my2.txt judge/plan-ids.txt`, 25 ids incl. reruns; 20 distinct edits): J-LOCK-DISABLE (must-die) KILLED by TestOperatorProvisionEpisodeExclusion/{approve,session_new}; J-HELD-INVERT KILLED same + TestPresetStartReadInteractive; J-REFUSE-EXITOK, J-REFUSE-NOERR (drops action id), J-LOCK-EPISODE-WRONG KILLED by Exclusion/{approve,session_new}; J-NOCLOSE KILLED by Exclusion/{approve_releases,session_new}; J-PRESET-PASS-NIL KILLED by TestStartLaunchHeadlessMintsAndLaunchesOnce, TestPresetStartReadParity/Interactive and TestEpisodeResume*; J-RESUME-PASS-NIL KILLED by TestEpisodeResumeChainUsesCurrentCredential et al; J-START-PASS-NIL KILLED by TestEpisodeStartProvisionsAndPrintsTheLaunchWithTheTokenInTwoFilesOnly; J-REPLAY-NOCHECK/NILERR/WRONGERR KILLED by TestOperatorConfirmReplay; J-TITLE-EMPTY KILLED by TestPresetStartReadParity; J-TTL-IGNORED KILLED by Parity + Refusals/bad_ttl; J-MCP-MODE KILLED by TestPresetStartReadInteractive; J-WALL-ZERO2, J-BINARY-EMPTY2 KILLED by TestPresetStartReadHeadless; J-TASKREAD-NIL-HEADLESS KILLED by TestStartLaunchPreFenceTaskBinaryAndDirectory/task. All kills are behavioural tests (none static/AST). Two initial COMPILEERR rows (J-WALL-ZERO, J-BINARY-EMPTY unused var) were re-edited (…2). **SURVIVORS (whole package rc=0): J-PRINT-AGENT (printLaunch agent blanked), J-DOUBLE-OVERRIDE-APPEND (persisted `Overrides` duplicated), J-QUICKSTART-b (docs sentence removed; no test asserts it).** Must-live control J-NEUTRAL-CTX (8 s -> 9 s) survived as expected.

Plan §4 UNMEASURED rows re-measured (named test, 17 rows; all KILLED, all behavioural): L-SNEW, L-APPROVE (Exclusion subtests), L-DUP (TestStartLaunchHeadlessMintsAndLaunchesOnce), L-OPEN-ERR, L-REFUSAL-ID, L-RELEASE (approve_releases), L-TRANSFER-RESUME (TestEpisodeResumeChainUsesCurrentCredential), R-READ-MODE (PresetStartReadHeadless), R-EP-GRAMMAR (Refusals/bad_episode), R-OVERRIDE-BRANCH, R-CLEAR-MODEL, R-TASK-READ (Parity), R-COLLISION-MCP (Refusals/mcp), R-DAEMON-BUILD (Refusals/build), R-FLAG-HOSTS (Interactive), R-LAUNCH-ONCE (Headless...Once), P-ARTIFACT-CHECK (TestOperatorConfirmReplay). Worktree restored; `git status --short` clean.

## Per-milestone non-vacuity
M1: with only the lock fallback disabled (J-LOCK-DISABLE) the Exclusion test fails in both session_new and approve subtests. M3: TestOperatorConfirmReplay fails with only the byte check removed. M2: terminal-path breaks (title, ttl, mcp mode, wall, binary, headless task) are each killed.

## Non-blocking findings (ranked)
1. `session mint` remains lock-free; not a stated activation gate (Appendix C mentions it only as "outside provisionEpisode"). Add to the S5b/S6 gate list.
2. Parity tests do not assert persisted `Overrides` or the printed `agent:` line independently (survivors above); parity is partly new-code-vs-new-code.
3. QUICKSTART sidecar paragraph has no behavioural test (only plan gate grep).
4. Replay kill needs host/authority; host/store alone passes the mutants (disclosed).
5. Load-time: not measured under -race whole package locally (controller sweep non-race 235 s; headroom OK given only +6 s added).

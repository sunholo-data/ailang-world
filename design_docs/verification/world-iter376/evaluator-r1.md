# Iteration 376 — judge round 1 (sonnet, Agent tool, own worktree at `3ea0bea9`)

**PASS 88/100, zero blocking.** Correctness 28/30, tests 25/30, security 19/20, standards 9/10, scope 7/10.

- **Mutants:** 94 run, each against its named killer test. 86 killed; 5 equivalent or no-op; 3 real gaps.
  - Equivalent or no-op: N-ROW-FORMAT, N6-second, N10-ERR, N10-NL (planner P3: the two checks are redundant), and N12-exit (the judge's own no-op).
  - Real gaps: **F-CREDS** (`approver.go:188`), **AS-CREDFOR** (`:203`) and **F-DB** (`:174`). No approver test exercises `--credentials-db` or the default `--db` (no test sets `HOME`).
- **Mutants that survived their named test but were killed by another:**
  - RL-TIMER-D, killed by CF1/CA9.
  - D-list, killed by `TestApproverList`.
  - X-OWNHELP, killed by the help tests.
- **Judge's harness note:** its first N-LOCK and K-LOCK attempts proved nothing (one did not compile, one ran with an invalid context). Re-run with a deadline context and a real lock held, both were killed.
- **Security properties confirmed:**
  - revoke marks the row and never deletes it;
  - no lock is held during the prompt;
  - the prompt is bounded at exactly 120 s, and a nil timer gives `exitFatal`;
  - `y`+EOF, a read error, a partial line, a timeout and a non-`y` answer all refuse and write nothing;
  - FIFO, non-regular and oversize key files are refused without blocking;
  - a private key file is refused without echoing its bytes;
  - names follow the strict pattern (Q2), and names and keys are never reused (Q3);
  - `readTermLine` is unchanged (R8);
  - exit codes use the named constants.
- **Q4:** the session listing on revoke (K8/K9) is deferred to Appendix B with M4/M5.
- **D-WORLD-109:** no `exec`, no `UserHomeDir`, no `~/.ssh`; QUICKSTART §15 is tagged attended-pending and was not run.
- **Scope:** the size gate was applied as written (the judge re-measured 1,896). The diff stays inside plan §5. M3 came in 16% over its estimate.
- **Gates in the judge's worktree:**
  - `go build ./... && go vet ./...`: rc 0.
  - `go test` of cmd, approval, store and authority: ok (cmd took 291 s).
  - `-race` on the approver tests: ok.
  - `-v`: 0 SKIP.
- **Clean status:** `git status` was empty after the last restore.
- **Carry to sprint 3:** tests for `--credentials-db` and the default `--db` (F-CREDS, AS-CREDFOR, F-DB).

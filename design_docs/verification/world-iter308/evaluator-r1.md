# Evaluation report: iteration 308, row 150, M7a-fix (round 1)

Tree: /Users/voightkampff/dev/sunholo-data/.eval-world-iter308 at 7475ec3c (base 0363b3ae). Machine load ranged 6 to 48 during the run.

## VERDICT: PASS, score 88/100
EVALUATION_SCORE: 88/100

| Rubric item | Score |
|---|---|
| Acceptance criteria / plan conformance | 26/30 |
| Tests and gates | 24/25 |
| Mutation / non-vacuity | 18/20 |
| Correctness of premise and risk analysis | 10/10 |
| Code quality, hygiene | 10/15 |

Deductions:
- Controller bound change 25 s to 60 s (Python) and 35 s to 70 s (Go) deviates from the plan text. It is measured and justified (see below), but the plan section was not amended.
- The stress row's pass condition is an observation when raw EPERM is 0 (plan-sanctioned).
- Redundant `Trials != 200` / `WallSeconds >= 60` checks.
- One `-race` package failure observed under load (see non-blocking finding 1).

## Blocking findings
None.

## Non-blocking findings
1. **One `-race` package run red under heavy load.** `go test -race ./scripts/marketingcheck -count=1 -timeout 240s` was red on run 1 with "shell wait daemon must be dead after bounded cleanup (pid 49360)" at supervisor_test.go:716. The package took 217 s of the 240 s timeout while load was 40 to 48. Run 2 was green in 60 s. I only kept the tail of run 1, so the full failure text was lost.
   - It is a process-death assertion in the unchanged `FakeDaemonCleanup/shell-wait` test, not a PermissionError refusal. It is a timing flake under extreme load, not evidence the fix is wrong.
   - Controlled comparison, `-race` of `FakeDaemonCleanup/(shell-wait|interrupted-shot)` x12 each:
     - fixed tree: 0/12 red.
     - base tree: 2/12 red, both "terminated trap must exit 143 after supervised cleanup: exit status 1" (the EPERM refusal signature).
   - The fix measurably removes the EPERM class. The controller's 25-run campaign should expect, and record, load-induced timing reds as a separate class. Repro: `go test -race ./scripts/marketingcheck -count=1 -timeout 240s` while load is above 30.
2. **Stress row is an observation when raw EPERM is 0.** The Go test does not fail on `raw_eperm==0`; it logs UNINFORMATIVE to stderr via `t.Log`. This is plan-sanctioned (mutant D5 below). The deterministic UNIT row is the primary guard.
3. **The `len(lines) != len(names)` guard is nearly redundant.** The duplicate and missing-case checks catch most stray lines. Only a stray line with a unique, unexpected case name would pass without the guard. Defensive, not load-bearing.
4. **Redundant checks in the stress row.** `Trials != 200` and `WallSeconds >= 60` overlap with the Python exit status and the in-Python deadline.
5. **Bound change not recorded in the plan.** 7475ec3c's 60 s and 70 s differ from the plan's 25 s and 35 s. Measured 32.4 s and 37.3 s at load 6 to 8 (this run), plus the controller's ~31 s. 60 s gives about 1.6x headroom, which is justified at normal rig load. It would not survive load of about 40 or more. The plan's "<30 s target" is missed on this rig; report it as such.
6. **Unrelated `sleep 60` processes.** `pgrep -f 'sleep 60'` returned PIDs during my run (65644, 31995, ...). Their parent is PID 432, `/bin/sh .../integrations/herdr/ailang-world/bin/world-watch`, so they are not from the Go tests. After my stress runs, pgrep's non-empty result was attributable to world-watch only. A PID is real; silence would not have proven anything.

## Conformance to the plan (derived from `git diff 0363b3ae 7475ec3c`)
- **Files:** changed files are the two executor files plus the plan MD and JSON (planner commit). No stray files, no `__pycache__` (`ls scripts | grep -c pycache` = 0).
- **Production diff:** exactly the planned handler, `except (ProcessLookupError, PermissionError):` plus the one-line comment `# macOS may report EPERM for dead descendants after the leader is reaped.` The docstring and the `os.killpg` call are untouched. The comment density matches the file.
- **Existing regression tests:** the test-file diff is purely additive (237 insertions, 0 deletions), so `FakeChromeTimeoutKillsGroup` and `FakeDaemonCleanup` are byte-unchanged.
- **UNIT test:** matches the plan.
  - Resolves python3 with `LookPath`, runs `-B -c`, imports through importlib as `marketing_capture`, with a 30 s context.
  - Monkeypatches `os.killpg` for the 3 cases and records `(pid, sig)` before raising.
  - Prints 3 JSON rows, with the 3-row guard and the exact messages ("must call killpg once with the expected pid and signal", "must swallow PermissionError EPERM", "must swallow ProcessLookupError ESRCH", "must propagate other OSError EINVAL unchanged", "inline kill_group probe must emit all three cases").
- **STRESS test:**
  - Opt-in skip message matches exactly.
  - 200 trials, 0.02 s settle, nested-shell session, 3 alternating raw/wrapper pairs, 600/600 counts, try/finally cleanup, file sinks, UNINFORMATIVE message, and the exact two assertion messages.
  - Deviation: bounds 60 s / 70 s instead of 25 s / 35 s (finding 5).

## Gate table (outside any sandbox)
| Command | rc | Key line |
|---|---|---|
| `gofmt -l scripts/marketingcheck` | 0 | empty output |
| `go vet ./scripts/marketingcheck` | 0 | clean |
| `go test -run '^$' ./...` (compile fence) | 0 | every package `ok ... [no tests to run]`; 0 non-ok lines |
| `python3 -B` ast.parse of `marketing_capture.py` | 0 | ok (py_compile avoided to prevent bytecode) |
| `go test -run '^TestMarketingKillGroupTolerance$' ... -v` | 0 | PASS, 3 subtests PASS (0.16 s) |
| `go test -run '^TestMarketingSupervisorFakeChromeTimeoutKillsGroup$' ... -v` | 0 | PASS (8.20 s) |
| `go test -run '^TestMarketingSupervisorFakeDaemonCleanup$' ... -v` | 0 | PASS, all 10 subtests (14.3 s) |
| `WORLD_KILLPG_STRESS=1 ... Stress` run 1 | 0 | trials 200, wrapper_calls 600, raw_eperm 15, raw_esrch 582, raw_success 3, wrapper_exceptions 0, wall 32.43 s |
| same, run 2 | 0 | trials 200, wrapper_calls 600, raw_eperm 54, raw_esrch 546, raw_success 0, wrapper_exceptions 0, wall 37.26 s |
| `go test -race ./scripts/marketingcheck -count=1 -timeout 240s` run 1 | FAIL | `shell wait daemon must be dead after bounded cleanup (pid 49360)` (finding 1), 217 s, load 40+ |
| same, run 2 | 0 | ok, 60.3 s |
| `-race` focused x12, fixed tree (shell-wait + interrupted-shot) | 0/12 red | |
| `-race` focused x12, base tree | 2/12 red | "terminated trap must exit 143 after supervised cleanup" |

## Mutant ledger
Each mutant was applied to one file, tests untouched otherwise, then restored from a scratch copy with `cmp` (identical).

| Mutant | Edit | Command | Result |
|---|---|---|---|
| M7a-fix-eperm (plan) | handler to `except ProcessLookupError:` | UNIT | KILLED: `supervisor_test.go:104: kill_group must swallow PermissionError EPERM` (only eperm-swallowed red, assertion reached) |
| M7a-fix-broad (plan) | handler to `except OSError:` | UNIT | KILLED: `:112: kill_group must propagate other OSError EINVAL unchanged` (only einval-propagates red) |
| M7a-fix-noop (plan) | body of `kill_group` to `pass` | UNIT | KILLED: `:99: kill_group must call killpg once with the expected pid and signal` on all 3 subtests |
| M7a-fix-noop (plan), existing test | same | `FakeChromeTimeoutKillsGroup` | KILLED: `:498: grandchild must be dead after bounded cleanup (pid 82080)` and `:499: child must be dead ...`. The `:483` verdict error (CleanupFailed) also fired, but the named grandchild assertion was reached as well. |
| D-eperm-stress | eperm mutant, run STRESS twice | STRESS | KILLED 2/2: 100 and 124 wrapper PermissionErrors, `:252 real kill_group must raise nothing after leader reap`. The first exception came in trial 0 and trial 4, so it reds almost immediately and reliably. |
| D1 | drop `calls.append` in the unit probe | UNIT | KILLED: `:99` on all 3 subtests |
| D3 | remove the 3-line-count guard, add a stray duplicate-case line | UNIT | KILLED by the duplicate check (`duplicate "extra"`) |
| D3b | stray line with the guard present | UNIT | KILLED by the guard (`:78`) |
| D3c (reasoned, not run) | guard removed, stray line with a unique unexpected case name | UNIT | would SURVIVE (finding 3); non-load-bearing |
| P1 | handler to `except PermissionError:` only | UNIT | KILLED: `:108 must swallow ProcessLookupError ESRCH` |
| P2 | re-wrap other OSError as a new `OSError` | UNIT | KILLED: `:112` (identity `Same` false); the `same` check is load-bearing |
| D5 | stress raw control replaced with `pass` (raw EPERM 0, raw_success 600) | STRESS | SURVIVED: PASS with only a UNINFORMATIVE log line. Load-bearing for informativeness but by plan design (the Go test never fails on a zero raw-EPERM count). The unit row covers correctness. |
| D6 | stress `range(200)` to `range(100)` | STRESS | KILLED: `:249 stress must complete 200 trials and 600 wrapped post-reap kills` |

## Stress-row assessment
- It is a real guard when EPERM is provoked: 100 to 124 of 600 raw calls hit EPERM under the eperm mutant, and the fixed tree gave 0 wrapper exceptions twice.
- On the fixed tree raw EPERM was 15 and 54 of 600, so EPERM is provoked reliably on this host, but only the combination "UNINFORMATIVE" log plus a count reveals vacuity. It is not a hard gate.
- The pass/fail on the fixed tree was reliable at load 6 to 8 (2/2 green, 32 to 37 s). The 60 s bound is justified at that load. It carries no margin at load of about 40 or more, where a spurious red is plausible.

## Premise check (man 2 kill / killpg, read myself)
`man 2 killpg` lists EPERM as "one or more of the target processes has an effective user ID different from that of the sending process" and ESRCH as "no process can be found in the process group". It does not say EPERM means no signallable member remains, as the plan's F8 states. The plan correctly does not claim that.

Is swallowing EPERM safe here?
- `spawn` uses `start_new_session=True` with an inherited uid and no credential override, so every group member is normally the same uid. EPERM on a same-uid group is then most plausibly the macOS post-reap quirk the iteration-307 repro shows: raw EPERM appears only after the leader is reaped (here 15 to 124 of 600 raw calls), and never raised an exception in the wrapper once swallowed.
- Concrete residual scenarios where swallowing could hide a real failure:
  - (a) A setuid or credential-changing descendant (for example `sudo` or a helper) stays alive and un-signalable, so EPERM from the group kill is swallowed and the process leaks. Not present in current callers: the Chrome and daemon helpers are the same uid.
  - (b) A live same-uid member that for some reason cannot be signalled (an uninterruptible kernel state) would also be silently ignored.
- Existing tests would still catch cleanup failures of this kind in the fake fixtures. `FakeChromeTimeoutKillsGroup` asserts grandchild and child PIDs dead via `kill -0`, and `FakeDaemonCleanup` asserts the daemon and its shells dead. These assertions are independent of the exception handling, as the noop mutant shows. Real Chrome and daemon processes are not covered by an unmutated death check, so (a) and (b) would not be caught in production runs. I found no concrete current-caller scenario, and I judge the risk low and adequately disclosed by the plan's F8 caveat.

## Hygiene
- `git status --porcelain` is empty; no `__pycache__`.
- The Go tests create no orphans: the stress row's try/finally reaps every session, and the only `sleep 60` PIDs seen belong to the unrelated `world-watch` script.

## Final `git status --porcelain` for the eval worktree
(empty; both source files restored and `cmp`-verified)

EVALUATION_ROUND: 1
EVALUATION_VERDICT: PASS score=88

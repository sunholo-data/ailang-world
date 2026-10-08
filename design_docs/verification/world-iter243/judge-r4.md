# Judge r4 — fix-3 follow-up 942ff3e (branch sprint/row148-slowbody-oracle)

## Verdict: PASS — 92/100, zero BLOCKING findings
Diff `e7f3422..942ff3e` is one file, `host/daemon/agui_test.go` (+66/-30). Production code unchanged.

| Rubric | Max | Score |
|---|---|---|
| Correctness (prod unchanged) | 30 | 28 |
| Test quality / AC coverage | 25 | 23 |
| Mutation resistance | 15 | 13 |
| Design / standards | 15 | 14 |
| Docs | 10 | 9 |
| Conflict surface | 5 | 5 |
| Total | 100 | 92 |

## 1. r3 items, each addressed; no oracle weaker
1. SeesDirectStoreCommit tight bound: waitEntry and the "commit late" check are now `2*tick + 1 s`; run budget 5 s. Addressed.
2. ShutdownEndsRuns 150 ms: now 2 s, and budget raised 1 s to 5 s so natural expiry cannot satisfy the oracle. This is stronger than before (M-g below).
3. `aguiWait` 150 ms: now 2 s. Addressed.
4. Goroutine-baseline and GlobalCap slot loops: now 2 s; GlobalCap held-run budget 10 s so held runs outlive the deadlines; the 17th-request guard is 2 s. Addressed.
5. GoldenRun / ResumeExact single-run byte comparisons:
   - GoldenRun budget 2 s; still ONE run compared byte-for-byte to golden + `RUN_FINISHED` (twice). Byte-exact, unchanged.
   - ResumeExact now drains and compares replays with `aguiReplayBytes`: first run's RUN_STARTED + snapshot verbatim, every entry frame verbatim, the final `RUN_FINISHED` verbatim, and each intermediate run's envelope checked separately (same RUN_STARTED, and a snapshot rebuilt from the previous run's cursor and logHead). Header, state (with forged `logHead`) and matching-cursor paths must be byte-identical to each other. Still byte-exact where it matters; only scheduler-dependent run cuts are factored out.
- Full-page test: budget 2 s, tick now 5 s, still requires the first run to deliver at least 101 entries, so it can only pass via the immediate re-read.
- One deliberate change in `aguiDrain`: it no longer stops at the first empty run; it resumes until the cursor reaches the last wanted index (or `want` is empty), still ending in exact `DeepEqual`. Sound (an empty run under load is legitimate), at the cost of a broken server that never delivers burning up to 200 runs (slow failure, not a false pass), bounded by the test timeout.

## 2. Mutation re-runs (`-race`, cp backup, shasum restore verified, tree clean)
| Mutant | Result |
|---|---|
| O19 (no immediate re-read on full page) | KILLED (`TestAGUIGoldenRun/full-page`, 10.6 s) |
| M-g (`RegisterOnShutdown` removed) | KILLED (`TestAGUIShutdownEndsRuns/false`, 5 s run, i.e. does not pass by natural expiry) |
| M-o (poll only once, no tail) | KILLED (`TestAGUISeesDirectStoreCommit`: "entry0 not polled within two ticks plus scheduler allowance") |

## 3. Load run
12 `yes` burners: `GOMAXPROCS=2 go test -race -count=3 -cpu 2 ./host/daemon -run AGUI` — rc 0, wall 90 s (package 89.0 s). `pgrep -x yes` = 0 afterwards.

## Findings
BLOCKING: none.
NON-BLOCKING:
1. AGUI -race wall time grew (about 30 s per pass under 2-CPU load, roughly 45 s per pass earlier at 16 CPUs unloaded for count=2): the 5 s/10 s budgets and 2 s guards are paid only by tests that hold runs open (O19 drill itself costs 10 s when it fails). Fine for the 600 s CI go leg but keep an eye on it.
2. The `aguiDrain` slow-failure mode in section 1 (diagnosis delayed, not wrong).
3. Carried: O2 watcher-join, O17 sleep clamp, O18 25 s guard survivors from r1.

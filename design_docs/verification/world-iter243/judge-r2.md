# Judge r2 — row 148 (iteration 243), HEAD 20f2f29 (test-only fix round)

## Verdict: PASS — 92/100, zero BLOCKING findings
(r1 was 90. r1 survivor O19 is now killed; CI-red root cause addressed.)

| Rubric | Max | Score |
|---|---|---|
| Correctness (production code unchanged since r1) | 30 | 28 |
| AC coverage / test quality | 25 | 23 |
| Mutation resistance | 15 | 13 |
| Design / standards | 15 | 14 |
| Docs / usability | 10 | 9 |
| Conflict-surface hygiene | 5 | 5 |
| Total | 100 | 92 |

Diff d2d1815..20f2f29 touches only `host/daemon/agui_test.go` and `design_docs/verification/world-iter243/judge-r1.md`; no production file changed.

## 1. Drain oracle review (`aguiDrain`) — sound
- Per run it requires: HTTP 200, body ends `\n\n`, no `RUN_ERROR`, exactly one `RUN_FINISHED` and it is the LAST event, `result.lastIndex` present and EQUAL to the last complete delta/snapshot cursor (`aguiCompleteCursor`), so the terminal cursor is cross-checked, not trusted.
- Concatenated indexes across runs are compared with `reflect.DeepEqual` to `want`, so dups, gaps and misordering fail. It stops when a run delivers zero entries or the cursor reaches the last wanted index; either way the exact-equality check runs, so a premature empty run cannot hide missing entries. Bound: 200 runs, then a fatal naming delivered vs wanted.
- Resume uses header (`Last-Event-ID`) or a `state` rewrite depending on the starting mode, so both D4 paths are exercised.
- No AC weakened: AC3.2 (`ResumeExact`) keeps header-vs-state byte equality, the conflict/malformed/404 arms and matching-cursor check, and adds exact drains for both modes; AC3.14, 3.2b, 3.3 (joined runs), 3.9 (now also drains per auth header), 3.12 (severance) all still assert the same sequences. `TestAGUISlowReaderCutIsResumable` shrank 5000 to 1000 entries and grew each entry about 5x (about 22 KB, roughly 22 MB backlog), so it still exceeds socket buffering and still requires a cut with cursor < 999; the resume half now uses the real handler via the drain instead of the cancelling recorder. Not weakened.
- New full-page check: budget 500 ms, tick 1 s, 205 entries; drain must equal 0..204 AND the first run must deliver at least 101 entries. With a 1 s tick, 101+ in one run is only possible via the immediate re-read on a full page, so it pins D5's rule.

## 2. Mutation re-runs at 20f2f29 (`-race`, cp backup, shasum restore verified, tree clean)
| Mutant | Result |
|---|---|
| O19 (no immediate re-read on full page) | KILLED (`TestAGUIGoldenRun/full-page`) |
| M-a-skip (cursor skips the last row of a full page; the executor's M-a-page equivalent, my first M-a-page edit did not match and is discarded) | KILLED (`full-page`) |
| M-c (tail via `GetLogEntry(cursor+1)`) | KILLED (`TestAGUIRESTGapCrossed`; also ResumeExact/GoldenRun) |
| M-v (`id:` before CUSTOM frame) | KILLED (`TestAGUISeveranceResumable`) |
| O21 (drop last of page, cursor not advanced) | SURVIVED, but equivalent: the dropped row is re-read next iteration, so no behaviour change. Not a finding. |

## 3. Remaining wall-clock / single-run assumptions (NON-BLOCKING)
1. `TestAGUIGoldenRun` (300 ms budget) still compares ONE run byte-for-byte to golden + `RUN_FINISHED` (twice), and `TestAGUIResumeExact` still compares a header run and a state run byte-for-byte for equality (300 ms budget). These are 5-entry runs, so they only fail if the process is descheduled for more than 300 ms between handler entry and the first read, or between the two runs. Residual flake risk only on a badly starved runner; fix if CI shows it: raise those two budgets to 2 s (the run still ends on budget, cost about 2 s each per k), or compare drained concatenations.
2. `TestAGUIRESTGapCrossed` uses a 30 ms budget but now drains, so it is safe. `TestAGUISlowBodyNeverTruncates`/`DisconnectReturns`/`ShutdownEndsRuns` are timing tests by design and assert upper bounds with slack.
3. Under `-race`, `TestAGUISlowReaderCutIsResumable` writes a roughly 22 MB backlog; it was 17.9 s before. Not re-timed individually; whole AGUI set below.

## 4. Load runs
Machine has 16 CPUs, so 8 `yes` is mild load; I also ran a harsher case.
| Command | rc | Wall |
|---|---|---|
| 8x `yes`: `go test -race -count=2 ./host/daemon -run AGUI` | 0 | 45 s (package 44.3 s) |
| 16x `yes` + `GOMAXPROCS=2 -cpu 2`: same, `-count=2` | 0 | 49 s (package 48.6 s) |
All `yes` processes killed afterwards (`pgrep -x yes` = 0 both times).

## Findings
BLOCKING: none.
NON-BLOCKING: item 3.1 above (two 300 ms single-run byte comparisons); carried r1 items 1 (O2 watcher-join, O17 sleep clamp, O18 25 s guard survivors; O19 now closed), 2, 3, 5 unchanged. I did not re-run the full unfiltered `go test -race ./host/daemon` or the other packages this round; the production code is byte-identical to r1, where they passed.

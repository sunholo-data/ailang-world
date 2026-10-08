# Judge r3 — test-only follow-up e7f3422 (PR #233)

> Banked by the controller from the judge's printed final message. The judge's own untracked file
> did not survive its checkout to the r4 head. The content below is the judge's verdict and findings.

**Verdict: PASS, 91/100, zero BLOCKING.** Diff `38db793..e7f3422` changes only `host/daemon/agui_test.go` (+77/−19); production code is unchanged.

Rubric: correctness 28/30 · test quality 22/25 · mutation resistance 13/15 · design 14/15 · docs 9/10 · conflict surface 5/5.

## Review
- **Accepting a complete 408 is consistent with D1/§6.** The invariant is a complete APIError, a stream ending in RUN_FINISHED/RUN_ERROR, or zero bytes. The 100 ms body bound racing the 150 ms scratch-server deadlines is a real race between two legitimate outcomes. The new arm rejects any `data:` byte, and a non-empty reply must parse as a 408 with `application/json`, Content-Length equal to the body length, and valid JSON with class `SlowBody` and nothing after it, so a truncated 408 fails. A connection reset is accepted as severance; a client-deadline timeout is not.
- **Anchoring is still proven.** The second arm raises the body bound to 5 s and delays the body past `aguiBudget`. An entry-anchored run must then emit RUN_FINISHED with zero `LogEntriesAfter` polls, which a poll-counting wrapper checks. This is structural, not a wall-time bound.
- **The `server-deadline-zero` arm forces the zero-byte path.** The body bound is 5 s, ReadTimeout 150 ms and WriteTimeout 50 ms, so the write deadline has expired before the stalled read fails. Any byte fails the arm.

## Mutants (race; cp backup, shasum-verified restore)
| Mutant | Result |
|---|---|
| M-q (clock anchored after the body read) | RED |
| M-r (body bound disabled) | RED (400 instead of 408) |
| O22 (truncated 408 body) | RED |
| O23 (25 s guard on the 408 write removed) | SURVIVED (needs a 25 s wait; as in r1) |

## Stress
With 12 `yes` burners on 16 CPUs, `go test -race -count=10 ./host/daemon -run TestAGUISlowBodyNeverTruncates` gave rc 0 in 10 s, and gave rc 0 in 11 s with `GOMAXPROCS=2 -cpu 2`. No burners remained afterwards.

## NON-BLOCKING — remaining tight wall-clock assertions
1. `TestAGUISeesDirectStoreCommit`: within `2*tick + 50 ms`. Suggest `2*tick + 500 ms`.
2. `TestAGUIShutdownEndsRuns`: terminal within 150 ms of `Shutdown()`. Suggest 1 s.
3. `aguiWait` (150 ms). Suggest 1 s.
4. 150 ms `runtime.Gosched` loops in `TestAGUIClientDisconnectReturns` and `TestAGUIGlobalCap`.
5. `TestAGUIGoldenRun` and `TestAGUIResumeExact` single runs compared at a 300 ms budget (r2 finding).

All of these were addressed by fix-3 (`942ff3e`); see judge-r4.md.

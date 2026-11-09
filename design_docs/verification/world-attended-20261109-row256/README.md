# Row 256: attended attribution under D-WORLD-111 = B (2026-11-09)

Row 256 was parked at iteration 362 (independent FAIL 85). Its full Go suite showed three failures beyond the two permitted baseline diagnostics. D-WORLD-111 = B (attended, 2026-11-09) changed the bar: every extra failure must either reproduce at the pinned base under paired same-environment arms (at least 100 each), or get a causal in-scope repair. Unmeasured or candidate-only failures still block.

Mark asked for the row to be finished outside the loop ("we can do it outside of the loop if it hurries things up a bit"). This directory is that session's evidence.

- **Base:** `dev` at `ad13d835`.
- **Candidate:** the row 256 implementation, rebased cleanly onto it (`b9d06456`, since rebased again onto `dc907675`).
- **Host:** the Studio, at a 1-minute load of about 4–9.
- **Interpreter:** `AILANG_BIN` was the pinned v0.41.0 interpreter.

## 1. Full suite

`go test -count=1 ./...` returned rc=0 on both the candidate and the base. `paired-sequential-summary.txt` has the result. Neither of iteration 362's two "permitted baseline" failures occurs with the pinned interpreter.

## 2. Paired arms, 100 runs per arm per test

Arms alternated: base, then candidate. Each arm used a compiled `host/broker` test binary, both in the same environment (`measure256.sh`).

| test (iteration 362's extra failure) | base | candidate | attribution |
|---|---|---|---|
| `TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce` | 0/100 | 0/100 | Did not recur. That makes 0 in 800 runs, counting iteration 362's 600. |
| `TestWorktreeSnapshotRefusesAConfiguredFilter` | 2/100 | 0/100 | Reproduces at base. This is the existing CI flake (queue row 200). |
| `TestExecProbeNet6ForeignDialIsRetriedNotRefused` | 0/100 | 5/100 | Candidate-only in this measurement. Attributed in §3. |

## 3. The net6 failure is an external dialer, and it reproduces at base

All net6 failures are the `raw-twice-proxy` subcase, with an identical refusal:

```
inconclusive after 3 attempts: every attempt's window held an accept that was not the child's (foreign 1, 2, 1)
```

That subcase plans foreign dials on attempts 1 and 2 only. So the `1` in attempt 3 is an accept that the rig did not make.

1. **The dialer.** A witness listener on a random loopback port (`listen_probe.py`, log in `scanner-witness.log`) received `GET / … User-Agent: SunholoHome/1.0` every 31–32 s. It came from the Studio's start page, which sweeps local listening ports. This is the "rig's start page" dialer that row 250 itself names (design V3).
2. **Why only the candidate.** In §2 a pair took about 10.5 s, so the 32 s sweep phase-locked onto every third pair's second arm. The failing iterations were 81, 84, 87, 93 and 96.
3. **Randomized arms.** `net6jitter.sh` ran 100 runs per arm in random order, with random gaps and the witness running. Base failed 0/100 and candidate 1/100. The single failure (18:17:49–51) overlapped a scan at 18:17:50. None of the 175 runs that missed a scan failed, on either side.
4. **Aligned arms, unmodified code.** `net6aligned.py` started each arm a fixed offset before the predicted scan, sweeping 0.2–1.8 s (`scan-aligned-sweep.tsv`):
   - At offset 1.0 s, **base and candidate both fail, with the identical signature**.
   - At every other offset, both pass.

   So the failure reproduces at the pinned base in paired same-environment arms, from the same external cause.
5. **Code reading.** The candidate's change to `exec_probe.go` adds a deadline timestamp in `fail()` and an `OutOfTime()` reader. It does not change `net6Leg`, `window`, classification or the retry bound.

**Verdict under D-WORLD-111 = B: criterion MET.** Every extra failure either reproduces at the base or did not recur in at least 100 paired runs per arm. No failure is candidate-only once the external dialer's phase is controlled.

## 4. Follow-up, not a row 256 blocker

Any process that sweeps loopback ports can fail `raw-twice-proxy`, because that subcase spends all three attempts. That's row 250's test design. It's for row 258, the hygiene triage, alongside the row 200 flake. A natural fix: the rig counts the accepts it made itself, and the test tolerates an outside one on the last attempt only if that attempt also carries a matching witness hit.

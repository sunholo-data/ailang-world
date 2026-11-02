# Evaluator r1, iteration 334, row 235 (worktree .judge-wt-iter334 @ 065c58cb)

Score: 88/100 -> PASS (no BLOCKING finding)
Rubric: correctness 32/35, test quality 26/30, honesty 12/15, code quality 9/10, scope 9/10.

## Q5 scope and pinned diff
- sha256 of host/daemon/decisionpacket_daemon_test.go = 877013b8...450383 = JSON fixed_sha256_prototype. daemon.go sha e6b70d7e... unchanged.
- `git diff --name-only a0439a2b 065c58cb | grep -v ^verification/world-iter334/` -> only host/daemon/decisionpacket_daemon_test.go. numstat 177/17 equals the JSON plan_diff_stat.

## Q1 CI failure mode
- Base test (a0439a2b file) + `time.Sleep(300ms)` before the second Dispatch, plain -count=20: 20/20 FAIL, each `:270 unwrapped probe max = 1, want 2` (and `:263 serialized probe max = 1` still logged = H2 vacuity).
- Fixed file, plain -count=20 and -race -count=20: rc=0. The `late-second` subtest is that gap (300 ms anchored on first entry). With lateSecond raised to 3 s (T8), plain and race -count=3 still pass.

## Q2 mutants (-count=20 each; "fails" counts parent + 2 subtests; M1-class: both subtests, 40 leaf fails + 20 parent)
| Mutant | Substitution | plain | race | killer |
|---|---|---|---|---|
| M1 | daemon.go serializedHandler: Lock/Unlock removed | KILLED 20/20 both subtests | KILLED 20/20 | :372 `serialized probe ran 2 of 2 executions without the daemon mutex held` |
| M1b | unlock before inner Execute | KILLED | KILLED | :372 same |
| M2 | test: gate.h = probe (bare probe in serialized arm) | KILLED | KILLED | :372 same |
| M5 | wrapper with a different `new(sync.Mutex)` | KILLED | KILLED | :372 same |
| M4 | bare probe + `held` removed (TryLock oracle off, arrival gate on) | KILLED | KILLED | :375 `serialized probe max = 2, want 1` (independent second signal works) |
| M7 | bare probe, held and arrived both off | KILLED | KILLED | :375 max = 2 |
| M8 | mutex released by a 5 ms goroutine (partial hold) | KILLED | KILLED | :375 max = 2 (TryLock sees it held on entry; max catches it) |
| M3 | awaitSecond timer 200 ms instead of probeWaitBound | KILLED 20/20 (late-second only; `together` passes) | KILLED 20/20 | :369 `serialized arm measured nothing ...` |
| BASE300 | base test + 300 ms sleep | FAIL 20/20 (the bug) | n/a | :270 |
| FIXED | no mutation | PASS 20/20 | PASS 20/20 | n/a |
Seam-test mutants (run -count=3): T1 ignore timeout channel -> KILLED (:422 "Execute did not return after the injected timeout within 5s", plain + race); T2 change why text -> KILLED (:426); T3 drop `p.missed = why` -> KILLED (:424).
Surviving mutants (gaps, see findings): T4 arm-level `serialized arm measured nothing` Fatalf -> Logf: rc=0. T6 same for CONTROL: rc=0. T7 delete `time.Sleep(stagger)` (late-second degenerates into together): rc=0.

## Q3 labels, clause by clause
- TryLock "deterministic": supported. Any wrapper not holding sh.mu around the inner call fails TryLock on the FIRST entry, whether or not a second arrives (M1, M1b, M2, M5 all die at :372 without relying on the arrival). Only external humanMu holder could false-fail; grep shows no other taker in this test file. Honest.
- Arrival gate: gate counts at Execute entry outside the mutex; the first entry waits on gate.second, so "second arrived while first still inside Execute" holds by construction (or `missed` is set). Prose that it does not observe Lock is accurate.
- Anchored stagger "empirical": honest. The fixed test passes at any gap (T8); the gap only matters to M3 (killed 40/40 here) and to AC-1.
- TestConcurrencyProbeNamesAMissedArrival: genuinely tests the probe's timeout seam deterministically (T1/T2/T3 killed). It does NOT test the arm-level checks; doc wording "fails by name" is true of the probe, not demonstrated for the arm Fatalf (T4/T6 survive).
- Sentences the code/diff does not support: doc header "About +126/-17 LOC" (actual +177/-17, JSON is right); "unheld==0 ... final diff :367" (actual :371-372); both stale from r2 (NON-BLOCKING).

## Q4 bounds
Probe first-entry wait: select on wait / timeout (10 s) / ctx (30 s). Serialized hold: 50 ms or ctx. Stagger entry wait: 10 s / ctx. dispatchPair: 15 s / ctx -> t.Fatal. Seam test: 5 s timers on both waits, buffered timeout channel. No unbounded wait found; a missing second dispatch yields "...did not arrive within 10s", never "max = 1" (M3 output shows it). After t.Fatal in dispatchPair goroutines are bounded by ctx, cancelled by t.Cleanup.

## Findings
1. NON-BLOCKING: arm-level `measured nothing` assertions (serialized and CONTROL) are unexercised. T4 (Fatalf -> Logf) and T6 pass rc=0. If the second dispatch ever fails to arrive and that check regresses, serialized arm silently passes at max = 1 (the original H2). Only the probe's `missed` string is tested.
2. NON-BLOCKING: nothing asserts the stagger really happened (T7 survives). The late-second subtest is only meaningful via M3-style kills; acknowledged empirical in the doc.
3. NON-BLOCKING: design doc stale numbers (+126 vs +177 LOC; line :367 vs :371/:372; says the final diff differs from prototype by one comment line while r3 added ~50 lines).
4. NON-BLOCKING: `missed`/`unheld`/`max` are read unlocked after wg completion; safe (happens-before via wg), -race clean; the t.Fatal timeout path reads nothing. Not a defect.
5. NON-BLOCKING: M3-style failure text says "within 10s" from the constant even when a mutated timer fired earlier; only reachable under mutation.

## Not verified
- Full `go test ./...`, verify_ail.sh, package-wide -race (controller's results at verification/world-iter334/controller-fullgate.txt taken as given).
- The actual Linux CI environment / the original red (only the mechanism by delay insertion on darwin/arm64).
- Behaviour under heavy CPU load (did not run a loaded soak; relied on the controller's 200/200).
- Mutants were -count=20 only for the table's rows; seam mutants only -count=3.
- Whether row 235's wording (full text in world-mission.md) is met beyond the AC mapping I checked.

## Worktree state
All mutated files restored from saved originals; sha256 of both files equals originals (877013b8... / e6b70d7e...). `git status --porcelain` empty.

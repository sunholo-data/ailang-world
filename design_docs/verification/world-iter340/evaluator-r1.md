# Evaluator r1 - iteration 340, row 242 (sprint/world-iter340-net6-hoist @ c1706638)

Verdict: PASS, 93/100 (pass line 70). Rig load 6-7 (uptime 15:23-15:27).

## Score
| Category | Pts | Score | Basis |
|---|---|---|---|
| Tests pass | 20 | 20 | Net6 tests 5/5 each at -count=5; exec-probe/srt/startup tests ok; full ./host/broker -v: rc 0, 1096 PASS / 39 SKIP / 0 FAIL |
| Lint | 10 | 10 | gofmt -l empty, go vet ./host/broker clean |
| Acceptance criteria | 30 | 27 | Mutant moving the bind after a slow step reds T1 and T2. Two gaps/residuals below |
| Code quality | 15 | 13 | Small, clear diff; AST test is intentionally strict (flags correct-but-refactored hoists) |
| Documentation | 15 | 14 | Design/plan/evidence thorough; design admits plant-error path UNMEASURED |
| Design fidelity | 10 | 9 | Matches section 3.1-3.3/4 |
| **Total** | 100 | **93** | |

## Production review
- Binding all N listeners after plantExecProbe succeeds and before the fan-out, with one deferred close loop after wg.Wait(): every listener is closed on every path after the bind (normal, early returns inside probeExecProfile, bind failure carried in l.err; close is nil-ln safe). No close happens before an owner's Run C: the only close is the parent defer, which runs after wg.Wait().
- In-process error semantics: probeExecProfile/probe changes are pure parameter threading; the error text and ordering of arm results are untouched. plantExecProbe failure still returns (nil, err) before any bind. Deferred closes run before env.cleanup (LIFO), harmless.
- Scope: only host/broker/exec_probe.go, host/broker/exec_probe_net6_test.go, .github/workflows/ci.yml (two test names added) plus design/evidence docs. 0 .ail files changed.

## Blocking findings
None.

## Non-blocking notes
1. M8 survives (leak on the plantExecProbe error path if someone binds before it and defers after). The design marks this path "UNMEASURED dynamically" (error table), so it is a known residual, not an introduced defect, and production is correct as written. Repro: apply M8 below, run `go test -count=1 -run '^TestExecProbeNet6' ./host/broker/` -> rc 0.
2. M4 survives, but it is an equivalent mutant: closing each listener after its own probe returns is safe (a closed port is not a live net6 port and the owner's Run C is already done). Not a gap.
3. T2 (AST scan) is stricter than the behaviour: a semantically correct hoist into a helper (M7) reds it. This fits "impossible by construction" but means a legitimate refactor must edit the test.
4. TestExecProbePassThroughShimFailsEveryConfinementArm flaked (accepts 4) 2 times in 8 runs on this loaded rig; the other 6 plus the full-package run passed. Consistent with design R2 (foreign dialer). Single-profile test, so the change does not alter its exposure materially; not counted as a regression, but the rate looked above the 5% quoted (small n, load 6-7).
5. T1 uses fd-scan via Getsockname and darwin-specific semantics (no SO_ACCEPTCONN); it is portable-ish on linux but untested there.

## Mutant table (production diff; test run alone, -count=1 unless noted)
| Mutant | Test expected | rc | Survived? |
|---|---|---|---|
| M1 bind moved into the goroutine (before probeExecProfile) | T2 | 1 | no (T2: call site in func literal) |
| M1b bind back inside probe() at the original site (after pre-probe Git), loop makes placeholders | T1 | 1 | no ("net6 port ... not a live listener") |
| M3 defer close body removed (leak) | T1 | 1 | no ("listeners left ... [52269 52270]") |
| M4 close each listener right after its own probe, defer removed (undesigned) | none expected (equivalent) | 0 (x5) | yes, equivalent, safe |
| M5 thread net6[0] to every profile (undesigned) | T1 | 1 | no (only one distinct net6 port) |
| M6 close loop breaks after first listener (undesigned) | T1 | 1 | no |
| M7 bind loop moved into helper bindAll called before fan-out (undesigned, semantically equal) | T2 | 1 | no (false positive by design) |
| M8 bind before plantExecProbe, defer after (leak on plant-error return) (undesigned) | none | 0 | YES, accepted residual (design "UNMEASURED") |
Pristine file restored and verified with cmp after every mutant (and after the base-code check).

## Vacuity checks
- Base `host/broker/exec_probe.go` (e22e0794) with the new test file kept: T1 FAILS ("pick 0: net6 port 55283 was not a live listener (slow profile in flight=true...)"); T2 FAILS ("probe (not before a go statement)"). Both are non-vacuous.
- T2 cannot pass on a reintroduced hazard: M1, M1b, M7 all red by T2 or T1; the only way past T2 is a single call site in ProbeExecProfilesEach before the first go, which is the fix.
- T1 also asserts that the first pick landed while the slow profile's Git was in flight, so it cannot pass silently without the slow-step condition.

## Numbers re-derived
- broker test census, `go test -count=1 -v ./host/broker/` at c1706638: 1096 PASS, 39 SKIP, 0 FAIL (203 s). Matches the banked 1096/39/0 (controller-gates-summary.txt). Top-level PASS lines: 338 (matches the 338 in the A2 census file; its 30 top-level SKIPs not separately recounted).
- T1+T2: 10 pass lines for -count=5 (mine) against 20 in the controller's -count=10 (A1); consistent.
- Mutation 27/27 (A3b) was not re-run in full; I ran 8 own mutants instead (above), all outcomes as expected apart from the two noted survivors.
- Controller A5 "ok=34 FAIL=4" is explained in its own file as a leftover artifact (store re-run rc 0); not re-run per instruction.

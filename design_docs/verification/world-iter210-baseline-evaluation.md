# Evaluation — ADMIN observational baseline (iter 210), independent Sonnet judge

Base commit verified: `git rev-parse HEAD` = 9b2cfac20f770811b113d62e987530f7e6299660 (detached evaluator worktree, clean).

## Verdict: PASS (narrow) — 84/100
Bankable ONLY as: "one observational, filesystem-sandboxed diagnosis-archaeology baseline of a local
verify_go.sh red, attempt 2 of 2 (attempt 1 = NOT-RUN/TRANSPORT)". NOT a completed World walk, NOT a
product-completion claim, NO bar credit.

## Artifacts inspected (all under .ailang/state/iter210-baseline-inputs/)
preflight.txt; baseline-report.md; baseline-report-2.md; baseline-transcript-20260930T102633Z.txt;
baseline.sb; sandbox-controls.json; sandbox-controls-2.json; base-verify-go.log; base-verify-go-2.log;
base-verify-ail.log; pkgproj-focused.log; baseline-corpus/{context.json,question.txt,operation.log,go.mod,go.sum,.github,cmd,host,scripts}.
Source read: host/pkgproj/iface_test.go:165-240, host/pkgproj/iface.go:55-195, scripts/verify_go.sh:15-20,300-320.

## Independent verification of the factual diagnosis
1. Corpus fidelity: `diff -rq` of corpus host/, scripts/, cmd/, .github/, go.mod, go.sum against the 9b2cfac
   worktree = no differences. corpus operation.log is byte-identical to base-verify-go.log (`diff` empty).
   context.json SHA matches the reviewed base.
2. Failure: base-verify-go.log:74-77 (and identical position in -2): `--- FAIL: TestQueryInterfaceReturnsWhileADescendantHoldsStdout (3.00s)`,
   `iface_test.go:232: instrument: descendant pid not recorded ... no such file or directory`, FAIL host/pkgproj; only failing package; :88 final FAIL.
   Cited line numbers in report-2 (74-77, 59-60, 88, 14-54, 55-58) all check out against operation.log.
3. Source mechanism confirmed: test shell double `sleep 30 &; echo $! > pidFile; wait` (iface_test.go:219-224);
   tinyBounds 3s/1s (:182); QueryTimeoutError check (:226-229) precedes the pid read (:230-232), so the
   error type check passed; no readiness handshake exists. Production: deadline set before Run, WaitDelay set (iface.go:169-174), QueryTimeoutError on deadline (:189-190); production bounds 30s/2s (:63-73). Report-2 is accurate here, and correctly labels "PID absent because shell did not reach its write" as inference.
4. Extra corroboration the report did NOT use (my own reading): the failing test took 3.00s, whereas every
   passing run takes ~4.00s (3s deadline + 1s WaitDelay held open by the live descendant; pkgproj-focused.log, and my rerun). Returning at 3.00s with no pid file implies no descendant held the inherited stdout, i.e. the shell never forked/wrote before the deadline (a ~3s spawn stall), which supports the instrument/scheduling reading and argues against a WaitDelay product bug. Still an inference; cause of the stall (load, macOS exec/scan latency) unproved.
5. Second identical full gate: base-verify-go-2.log vs base-verify-go.log differ only in addresses, goroutine ids, timings, and the random TempDir suffix
   (3786632707 vs 2016941257) — same test, same line 232, same ENOENT, same 3.00s, same stopping point (plain leg; race leg and final
   success never reached; verify_go.sh set -euo pipefail :20, plain leg :310-312). VERIFIED: 2/2 full-gate reds at the same instrument.
6. Focused five: pkgproj-focused.log shows 5/5 PASS with live descendant (alive: true), ~4.0s each. I independently reran
   `go test ./host/pkgproj -run TestQueryInterfaceReturnsWhileADescendantHoldsStdout -count=5 -v` in this worktree: 5/5 PASS, 4.00-4.02s, descendant alive, ok 20.36s. VERIFIED.
   Caveat: pkgproj-focused.log carries no command header/SHA; its provenance is by controller assertion only (my rerun substitutes for it).
7. verify_ail: base-verify-ail.log ends "verify gate PASSED: 16 required identities verified, 40 named tests pass"; not contradicted.
   Result overall: base gate is a GREEN ail leg + a plain-Go-leg RED that reproduced 2/2 full, 0/10 isolated (5 controller + 5 mine).

## Timer / fence controls
- Sandbox: baseline.sb is `(allow default)(deny network*)(deny file-read* /Users /private/tmp /tmp)(allow file-read* corpus)(deny file-write*)(allow file-write* /dev/null)`.
  Attempt-1 profile (deny default) could not pass its own positive control (sandbox-controls.json rc=-6) and is correctly declared
  NOT-RUN; no question bytes were obtained (baseline-report.md). Attempt-2 controls (sandbox-controls-2.json; transcript lines 1-8):
  positive read of question.txt rc0; negative reads of ailang-world/.git/config, ~/.ailang, /private/tmp all rc1 "Operation not permitted". Controls precede the timer.
- Timer: transcript prints START 10:26:33.427670Z / monotonic 1803667432099791 at the question reveal (transcript line 16). Report-2 elapsed 82.938396125s:
  I recomputed 1803750370495916-1803667432099791 = 82,938,396,125 ns — arithmetic exact; wall stamps agree to <1ms.
- Every diagnostic command in the transcript goes through `sandbox-exec -f baseline.sb`, with 5s/10s subprocess timeouts. The two path errors (rg rc71, env.go rc1) are honestly disclosed and not used.

## Findings / weaknesses (none blocking for narrow banking)
W1 (moderate). The STOP stamp (10:27:56.366Z / monotonic stop) is NOT in the transcript; only START is machine-recorded there. The stop is asserted in
   the report. The elapsed number is therefore internally consistent but not independently transcript-verifiable; treat as approximate (~83s), one sample.
W2 (moderate). The fence is command-level: it constrains only commands routed through sandbox-exec. It cannot constrain the diagnosing agent's memory, and
   the diagnosing agent/controller had seen this red beforehand (attempt-1 report references "the local gate red"). So "blinded" = filesystem-blinded, not cognitively blind. Any timing comparison must carry that.
W3 (minor). pkgproj-focused.log lacks command/SHA header (see 6).
W4 (minor). Attempt-1 to attempt-2 profile changed between attempts; attempt 1 sandbox claim was never assessed as a control — fine, since it is declared NOT-RUN.
W5 (minor). Correct stance in report-2 ("no product regression demonstrated") is appropriately weak; do not upgrade it to "no product regression exists" — race leg never ran locally on this SHA, and the 3s stall itself is unexplained.
W6 (note). The spawn-latency comment at iface_test.go:178-181 shows this exact test has flaked under full-suite load before (row-114 candidate per memory); the baseline is consistent with, but does not prove, a recurrence of that flake class.

## Claim limits
- ONE observational baseline, one diagnostician run, n=1 timing. Not a walk, not a design/quorum result, not a regression proof.
- Proved: named failing test and line; reproducible 2/2 in full gate; 0/10 in isolation; instrument-setup failure point (PID file absent after accepted QueryTimeoutError).
- Unproved: why the shell did not write the pid (load vs macOS exec latency vs other); local -race leg outcome; remote per-test timing; production behaviour under the intended stimulus on this run.
- No product-completion claim. No bar credit. No World walk credit. The pending revised design and quorum are outside this review.
- Recommendation for the ledger entry: bank as "observational baseline, attempt 2, sandbox-fenced, elapsed ≈83s (stop stamp untranscribed)"; row-114 flake candidate stays a candidate.

## Score /100
Diagnosis correctness 27/30; controls & timer 20/25 (−W1 stop stamp, −W2 cognitive blind gap); corpus/transcript validity 20/20; claim discipline 12/15; reproducibility of gate/focused claims 5/10 → bumped by my own rerun; total 84/100. **PASS for narrow baseline banking.**

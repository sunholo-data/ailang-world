# Evaluator r1 - iteration 359, row 250 (adc19bd7, impl d29c4155)

Verdict: PASS, 91/100 (pass line 70). Rig load 3-5 (uptime 19:30).

## Score
| Category | Pts | Score | Basis |
|---|---|---|---|
| Tests pass | 20 | 20 | Required set -count=2: 0 FAIL (all Net6, PassThroughShim, ReverseLeg tests twice). Controller's full ./host/broker ok (210 s) |
| Lint | 10 | 10 | gofmt -l host/broker empty; go vet ./host/broker clean |
| Acceptance criteria | 30 | 27 | Foreign-dial retry, child/silent still refuse, exhaustion refuses, reverse verdicts unchanged. One test gap (survivor mutant X3) |
| Code quality | 15 | 13 | One shared readProbeLine/judgeLeg/attempt bound; clear comments. net6Leg takes 8 params and a *int cursor |
| Documentation | 15 | 14 | Design r3, plan, evidence thorough and traceable |
| Design fidelity | 10 | 7 | Matches design, but retry-token re-check is unpinned (see N1) |
| **Total** | 100 | **91** | |

## Production review
Q1 Fail-closed. Reading exec_probe.go net6Leg (~l.669-694), window (~l.775-800), judgeLeg (~l.190-200), readProbeLine (~l.860-890):
- A leg passes only if its token is a refusal AND judgeLeg returns revPass, i.e. its window held NO child, silent or foreign accept; or a later attempt does. Child wins over everything (judgeLeg case 1), and a non-refusal token fails before retry (revNoToken). So a real reach (child line, whole) fails at once on attempt 1.
- Pending accepts: window waits up to execProbeWindowWait (1.5 s) and counts pending as silent -> refuse. Split child line: readProbeLine reads byte-wise until the line, so split reads are fine; a late line (>500 ms) is silent -> refuse. Reader errors / EOF with zero bytes / strict prefix -> silent -> refuse. Only a non-empty non-prefix byte string is foreign.
- Forged nonce: the child nonce is the child's own argument, so only the sandboxed child (or the R2 relay case) can produce it; a wrong nonce is foreign (retried), which is the stated R2 residual, bounded by execProbeAttempts=3 (attempt==3 returns inconclusive; no off-by-one: test X10 reds).
- Accept between windows: windows are contiguous ([from,next)), `*from = next`, so a late accept falls into the next leg's window and still counts (it can fail the next leg, never be lost). After the final (proxy) window the 200 ms sleep precedes it, same exposure as HEAD. Accepts before mark() (control + pre-run foreign dials) are excluded, but the child cannot predate its own run. Control: waits for a classified net6Control, so control cannot be satisfied by a foreign dialer (as HEAD's count==1 could).
- No path found where a genuine reach passes that HEAD would fail other than R2. One test-only gap, see N1.

Q2 Reverse leg. judgeReverse -> judgeLeg(tokOK, seen, retryTok) has identical case order (child, no-token, silent, foreign||EADDRINUSE). classifyReverseDial now uses readProbeLine(hello, want+"\n"): buffer capacity and stop at '\n' unchanged; new early stop at the first non-prefix byte only shortens a buffer that is already non-empty and not equal to want+"\n", so the switch (equal / len>0 / clean EOF / silent) yields the same verdict for every input; the hello is still written, deadline unchanged. Constant rename execProbeReverseAttempts -> execProbeAttempts (=3, same value). Reverse-leg tests all pass x2.

Q3 see mutant table. Q4: net6Dial in exec_probe_net6_dial_test.go does net.DialTimeout to the listener port from the test process while the shim client blocks on go-N, so the dial goes through the real accept loop and classifier; no mock. MUT-RAW-COUNT is red on TestExecProbeNet6ForeignDialIsRetriedNotRefused (A5-mutants.log l.207, expected_rc=1 actual_rc=1; the PassThroughShim test is expected green there). Q5: all four names appear exactly once, on the verbose real-srt list (ci.yml l.247 `tests=`), which the step loop asserts `--- PASS: $t ` for. Q6: diff outside design_docs: .github/workflows/ci.yml (1 line), host/broker/exec_probe.go, exec_probe_net6_dial_test.go (new), exec_probe_net6_test.go (regexp for the new nonce suffix), exec_probe_test.go (expected strings, reverseLeg arg). No .ail files, nothing under tools/launchd/. Impl files at adc19bd7 are identical to d29c4155 (A1-hashes, git diff empty for host/.github).

## Blocking findings
None.

## Non-blocking notes
1. N1, survivor mutant X3: deleting `tok = r.token(0)` in net6Leg's retry loop (host/broker/exec_probe.go ~l.689) keeps all Net6 + PassThroughShim tests green (rc 0). Effect: after a foreign window, a retry run whose token is not a refusal (e.g. WORLD-PROBE-CONNECTED, ERROR) would be judged on the stale attempt-1 refusal and pass if its window is empty. Production code is correct; the test suite does not pin the per-attempt token re-check, and HEAD would have failed such a run by token. Suggest a ForeignDialIsRetried variant whose second run emits a non-refusal token.
2. Window/hello share one constant (execProbeHelloWait 500 ms) as the read deadline; a loaded rig delaying a child's line past 500 ms makes arm 6 refuse (fails closed, delays a pass only). Documented in the constant's comment.
3. Controller A7 verify_go rc=1 is attributed to its TMPDIR override; confirmed (below). Evidence directory is very large (80 files, ~13k lines), mostly designer/executor transcripts.

## Mutant table (own mutants, scratch copies of HEAD, run `go test -count=1 -run '^(TestExecProbeNet6.*|TestExecProbePassThroughShimFailsEveryConfinementArm)$'`)
| Mutant | rc | Killed by |
|---|---|---|
| X3 drop `tok = r.token(0)` in retry loop | 0 | SURVIVED (N1) |
| X8 ignore silent accepts in net6Leg's reverseSeen (`silent: false`) | 1 | TestExecProbeNet6ChildOrSilentAcceptStillFails |
| X9 retry runs op "connect" regardless of leg op | 1 | TestExecProbeNet6ForeignDialIsRetriedNotRefused |
| X10 attempt bound off by one (`attempt > execProbeAttempts`) | 1 | TestExecProbeNet6ForeignEveryAttemptRefuses |
(First attempt at applying X3-X10 with BSD sed did not apply; verified the diff with diff(1) after redoing with perl, and each scratch differed from HEAD only on the mutated line. Scratch dir deleted.)

## Evidence honesty spot-checks
- "A4 ... 19 tests rc=0 pass=1": the required set I ran myself passes (twice); the A4 logs list the same names.
- "A5 runner match=177 mismatch=0": A5-mutants.log l.355 reads "177 named-test rows; 0 mismatches; 21 scenarios"; MUT-RAW-COUNT rows read as quoted.
- A7 verify_go rc=1 attributed to TMPDIR: A7b shows MOD010 with the override and ok without, and 0 files under cmd/ changed. I re-ran `TestDemoProofNoVerifiedSetRefuses` with default TMPDIR: ok. A7c (verify_go.sh, plain leg, default TMPDIR) rc=0. Attribution holds. Caveat: the full -race leg was not re-run by the controller (CI-only).
- A1 hashes: git diff d29c4155..adc19bd7 over host/ and .github/ is empty.

## What I ran
- `git diff --name-only 1835d8be adc19bd7 | grep -v ^design_docs/` -> 5 files listed above; grep for .ail/launchd -> none.
- `go test -count=2 -v -run '^(TestExecProbeNet6.*|TestExecProbePassThroughShimFailsEveryConfinementArm|TestExecProbeReverseLeg.*)$' ./host/broker/` -> ok 57.7 s, every test PASS twice, 0 FAIL.
- `gofmt -l host/broker` -> empty; `go vet ./host/broker` -> clean.
- Mutants X3, X8, X9, X10 on scratch copies -> see table.
- `go test -count=1 -run ^TestDemoProofNoVerifiedSetRefuses$ ./cmd/ailang-worldd/` (AILANG_BIN pinned) -> ok.
- Not run: the full ./host/broker suite and the linux real-srt CI step (relied on controller's A7-broker ok and CI).

# Evaluator r1 — iteration 318, row 232 (w-exec-probe-reset-holder-false-pass), head a828c450

Verdict: PASS, 92/100. (correctness 28/30, test quality & non-vacuity 24/25, design conformance 19/20, evidence integrity 13/15, scope & docs 8/10)

## What I verified (commands in the worktree .wt-world-iter318-eval, AILANG_BIN=$HOME/.pinned-ailang/ailang)
- Change set: `git diff --stat df5a0671 a828c450` = production `host/broker/exec_probe.go` (+25/-8: `dialErrIsSilent` = `errors.Is(err, syscall.ECONNRESET)`, `pollLoopback` dial seam, `else if dialErrIsSilent(err) { silent.Store(true) }`), tests (+93), ci.yml (1 line), QUICKSTART (1 clause), parent-doc amendment notes, design/plan/evidence.
- `go vet ./host/broker/` clean; `go test ./host/broker/ -run 'TestExecProbeReverse|TestExecProbeForeign' -count=1` ok.
- T9 determinism: `go test ./host/broker/ -run 'TestExecProbeReversePollRecordsConnectReset|TestExecProbeReverseDialErrIsSilent' -count=200` ok (0.8 s). T9 is timing-free: first dial happens before the first stop check, injected dial, extra dials return ECONNREFUSED.
- Fail-closed (P20): the only production behaviour change is that one more event sets `silent`. `judgeReverse` order is unchanged (child -> REACHED first; no-token -> NoToken; silent -> refuse; foreign -> retry), so a real reach still wins and the reverse leg can only refuse more, never pass more. Verified in source (`exec_probe.go:129-141`).
- No false refusals: ECONNREFUSED, EINVAL, ECONNABORTED, EADDRNOTAVAIL, timeouts, os.ErrDeadlineExceeded, nil are all asserted non-silent in T8 (real refused and real expired-context dials are controls checked first); wrapped errors (`fmt.Errorf %w` of OpError->SyscallError->Errno) and bare Errno are silent (`errors.Is` unwraps).
- CI: ci.yml step 8 (`Row 140 ... linux`) adds `TestExecProbeReverseDialErrIsSilent TestExecProbeReversePollRecordsConnectReset` to `$tests`; the same list feeds both the `-run "^(...)$"` regex and the `for t in $tests; grep -q "^--- PASS: $t "` loop, so a SKIP or a missing test fails the step. Names match the test functions exactly.
- Docs vs code: QUICKSTART §10 arm-6 "(or that resets the host's dial)" matches `dialErrIsSilent`; parent header and §4.2 amendment notes match the diff.
- Linux evidence: `linux-errno-probe-run40842980481.txt` shows linux/amd64 errno `connection reset by peer(0x68) is-ECONNRESET=true` against the rst holder and `connection refused(0x6f)` only with no listener (idle and under `yes` load), and 7 `--- FAIL: ...SilentHolderRefusesSticky/reset-holder` at base under load (grep -c). `ac-linux-load-run40849609705.txt` (found in the executor worktree .wt-world-iter318-row232, NOT in the a828c450 tree): `PASS_LINES=300 FAIL_LINES=0 t7 rc=0` with 8 `yes` on 4 cores, HEAD 8b6defdb "AC-LINUX-LOAD at fixed head a828c450". `git diff --stat 8b6defdb a828c450 -- host/` is empty, so the tested production+test code is identical to the judged head. The `FAIL_LINES` text I grepped on that log is only the echo of the command itself, not a failure.
- Controller local verify_go reds (`controller-verify-go-a828c450.txt`): `TestDoctorQuickIsFast` (4.088 s timing) in cmd/ailang-worldd and `TestProtectedGitFilterPreflight/configured-process-include`, with a recorded load average of 5-8 and Firefox/git at the top of CPU; both re-pass `-count=1` x3 (ok ok ok). Neither touches host/broker exec_probe; attribution to load is credible (timing assertion plus a git-subprocess preflight). Not independently re-run (shared rig).

## Mutation table I ran (each against `go test ./host/broker/ -run 'TestExecProbeReverse|TestExecProbeForeign' -count=1`, restored from /tmp backup)
| Mutant | Result (failing rows) |
|---|---|
| M22 dial errors ignored (= base behaviour) | RED: T9 connect-reset |
| M23 ECONNRESET never silent (`&& false`) | RED: T8 op-reset, bare-reset, wrapped-reset; T9 connect-reset |
| M24 `err == syscall.ECONNRESET` | RED: T8 op-reset, wrapped-reset; T9 connect-reset; bare-reset green (control) |
| M25 ECONNREFUSED is silent | RED: T8 bare-reset, op-reset, wrapped-reset, op-refused, real-refused; T9; ForeignListener control/chatty and Chatty-retry rows |
| M26 any non-nil err silent | RED: T8 addrnotavail, deadline-exceeded, op-aborted, op-einval, op-refused, real-refused, timeout; ForeignListener control/chatty; Chatty rows |
| M27 EINVAL silent | RED: T8 op-einval; T9 einval-then-refused |
| M28 timeout silent | RED: T8 timeout, deadline-exceeded |
| M28b ErrDeadlineExceeded silent | RED: T8 deadline-exceeded only (timeout row stays green) |
| M29 connect reset sets foreign | RED: T9 connect-reset |
| OWN1 ECONNABORTED also silent | RED: T8 op-aborted |
| OWN2 reset sets `child` (would turn a reset into REACHED) | RED: T9 connect-reset |
| OWN3 the new `else if` branch disabled (`false &&`) | RED: T9 connect-reset |
| OWN4 seam bypassed (poll ignores injected dial) | RED: T9 all three rows |
(My first M23 and M26 attempts failed to compile via an unused-import error of my own making; I reran them with compiling forms, results above.) Restore proven: `shasum -a 256 host/broker/exec_probe.go` = 390c74e841f26a1283743c4dc95efe6d31baf3ccc44b77e64d6af8283e8a23cd before and after; test file never mutated.

## Findings
1. NON-BLOCKING: evidence for AC-LINUX-LOAD and the controller verify log is not in the judged commit a828c450 (they live untracked in .wt-world-iter318-row232). `git diff --stat df5a0671 a828c450` lists neither. Commit them before/with merge or the docs' evidence claims are unreproducible from the tree.
2. NON-BLOCKING: the design's §6 killer summary table still lists M23..M28 for T8 without M28b (executor report deviation 6, `grep -n M28b design_docs/planned/w-exec-probe-reset-holder-false-pass.md`). Cosmetic.
3. NON-BLOCKING: the executor skipped the plan's focused base run (deviation 4) and its mutation-log files are in /tmp; I re-ran the drill myself, so no residual gap.
4. NON-BLOCKING residual risk: Linux idle rst-holder probe shows connect-time resets are rare per dial (1742 of ~644k dials, 30 s); a listen window (execProbeListenMS) with few dials could in principle see none. AC-LINUX-LOAD's 300/300 PASS under 2x-nproc `yes` load on one kernel (7.0.0-1017-azure) is the only direct evidence; the design already names this (single kernel) as a limit. A miss would degrade to the old false-pass, not to a false refusal.
5. NON-BLOCKING: connect-time timeouts against a saturated reset holder (143-179 per probe) are not silent by design; consistent with the stated fail-open-to-old-behaviour boundary, and T8 asserts it.

No BLOCKING findings. Worktree left byte-identical: `git status --porcelain` empty (plus this untracked report under design_docs/verification/world-iter318/, which is excluded from the "empty" claim only because the task asked me to write it).

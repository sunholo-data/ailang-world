# Evaluator report: row 156 M4c (885af2d1), iteration 286

VERDICT PASS  SCORE 84/100

Breakdown: correctness/spec fidelity 27/30; test quality and non-vacuity 19/30 (5 spec'd negative controls plus 2 controller-found mutants killed; 7 of my 12 independently designed mutants survived); gates/verification 15/15; scope/standards/design 13/15; code clarity/cost 10/10.

## Gates (env -u WORLD_TOOL_AILANG_BIN AILANG_BIN=$HOME/.pinned-ailang/ailang)
- go vet ./... rc 0
- TestExecProbeArm8 (no srt env): SKIP, rc 0 (the skip is expected without srt)
- TestExecProbeArm8FailsUnsandboxed (6 subtests), TestExecProbeArm8Absent, TestExecProbeProtectedSample (11 subtests), TestExecProbeArm8Wired (2 subtests): PASS
- Real srt darwin (WORLD_EXEC_SRT_NODE_MODULES, WORLD_EXEC_NODE): TestExecProbeArm8 PASS (0.91 s, "arm8 directory refusal; darwin basename oracle passed"); TestExecSrtStartupProbePasses PASS
- go test -race -count=1 ./host/broker/... rc 0 (180 s) with TMPDIR outside the worktree. The first attempt, with TMPDIR inside the worktree as instructed, failed one unrelated test (TestReplayReturnsTheRecordedPublishResultWithZeroNetworkCalls: "registry credential file is inside the working tree"). That is an artifact of my TMPDIR placement, not of the commit.
- Linux/bubblewrap behaviour: UNMEASURED-LOCAL. That covers the linux branch with real bwrap, the base-name-write-succeeds claim, and the EBUSY/EXDEV acceptance.

## BLOCKING findings
None.

## NON-BLOCKING findings
1. Guards in the new code that no test kills (mutants below). Several are explicitly spec'd behaviours, so a follow-up should add rows to TestExecProbeArm8FailsUnsandboxed:
   - Linux success-token requirement (spec: "Linux basename write requires success token") is unkilled. A linux row whose client prints REFUSED should fail with "missing success token".
   - The `!ok` host-file-deleted check on the refusal path is unkilled. A client that deletes the planted file and prints REFUSED EPERM would pass, mutated or not.
   - The darwin `absent(target)` check on creation is unkilled. A client that creates the file but prints REFUSED would pass.
   - `target == basename` on the linux branch is unkilled. No linux row uses a directory-form pattern, so the directory leg on linux (must be refused, not written) is untested. On real Linux this is the only thing separating the dir leg from the detect leg.
   - The EBUSY/EXDEV acceptance is unkilled. Nothing asserts execRefusalCodes("linux").
   - The `checked.err`/`checked.intact == nil` fail-closed guards, the base-name parent-avoids-protected-dir loop, and the final `matchProtected` guard in the sample builder are unkilled. They are defensive, and the controller's survivors list should name them as accepted residuals.
2. (Controller point 1) execRefusalCodes("linux") now accepts EBUSY/EXDEV for ALL arms (2-4 and the others that use it), not only arm 8. V14 supports the codes as real bubblewrap refusals, and arms 2-4 each also check host state (outside file absent and similar), so the widening cannot turn a leak into a pass. It is still a behaviour change to the legacy arms, beyond arm-8 scope, and no test pins it. It is safe but should be called out; an arm-8-only list would be tighter.
3. (Controller point 2) execProbeProtectedSample brute-forces runes 32..MaxRune per class glob. It is correct: it uses the grammar's own folded path.Match plus a final matchProtected check, and '/' is skipped. The cost is small for realistic patterns (under 1 ms). It is bounded but not trivial for classes that match only very high runes: a single `x[U+10FFF0-U+10FFFF]` takes about 345 ms, an unsatisfiable class about 225 ms before the explicit error, and 20 such classes in one 256-byte pattern take about 5.5 s of startup. Startup-only and the input is operator-authored, so it is non-blocking, but a cap or class-aware construction would be better. "[^a]" yields a space-named specimen; it works.
4. The real-srt TestExecSrtStartupProbePasses was not extended with a protected set, so the production path with a passing arm 8 is never exercised on real srt. TestExecProbeArm8 calls probeProtected directly, and TestExecProbeArm8Wired only proves arm 8 runs and fails under a pass-through shim.
5. A missing `git` binary now fails startup as arm8-protected-write when `protected` is declared (the error names the arm). Acceptable, fail-closed.
6. Only the first directory-form and first base-name pattern are probed (matches the spec). The "control" write before "BASE" in the plant loop is behaviourally redundant (mutant H survives, equivalent).

## Mutants
Controller's 7 mutants were not re-run individually, apart from J (a new total-count variant). All other rows are my own designed mutants. Each run was `go test -count=1 ./host/broker/ -run 'TestExecProbe|TestExecStartup|TestExecRefus'`, and each file was restored by copying the saved bytes and confirmed with cmp.

| Mutant | Change | Result |
|---|---|---|
| K-NO-EBUSY-EXDEV | drop "EBUSY","EXDEV" from the linux codes | SURVIVED (ok) |
| A-LINUX-NO-WROTE-TOKEN | `if r.token(0) != "WORLD-PROBE-WROTE"` -> `if false` | SURVIVED |
| D-NO-OK-CHECK | drop `!ok` from the refusal-path check | SURVIVED |
| E-NO-ABSENT | drop `!absent(target)` from darwin creation | SURVIVED |
| L-PARENT-NO-AVOID | `if !matchProtected(dirs, parent)` -> `if true` | SURVIVED |
| B-NO-INTACT-NIL | drop `checked.intact == nil` | SURVIVED |
| C-LINUX-ERR-IGNORED | drop `checked.err != nil` | SURVIVED |
| F-SEEN-NOOP | `if seen[directory]` -> `if false` | SURVIVED (near-equivalent) |
| G-DOTS | drop the final matchProtected from the sample guard | SURVIVED |
| H-ONLY-CONTROL | plant only "BASE" | SURVIVED (equivalent) |
| I-LINUX-ALL-TARGETS | `goos == "linux" && target == basename` -> `goos == "linux"` | SURVIVED |
| J-TOTAL | `total: len(order)` -> `total: 0` | KILLED: --- FAIL TestExecProbeArm8Wired/protected |

Net: 1 of 12 killed, 11 survived. Of those, A, D, E, I and K are spec'd behaviours and therefore real test gaps (finding 1); the rest are defensive or equivalent.

# Evaluator report round 2: row 156 M4c (9edf1162)

VERDICT PASS  SCORE 91/100

Breakdown: correctness/spec 28/30; test quality and non-vacuity 25/30; gates 15/15; scope/standards 13/15; clarity/cost 10/10.

Delta from 885af2d1 is test-only (exec_probe_test.go, +34); exec_probe.go is byte-identical (cmp against saved bytes, clean after every mutant).

## Gates
- go vet ./host/broker/ rc 0.
- TestExecProbe*: linux-detect, linux-no-token, linux-checker-error, linux-directory-write, darwin-create, darwin-create-forged and darwin-refuses-all PASS; linux-ebusy-passes and linux-exdev-passes SKIP on darwin, by name (t.Skip with a reason).
- Not re-run this round: the real srt and full -race gates. Round 1 measured them on the same exec_probe.go and the delta is test rows only. The controller reports real srt 2/2 PASS and broker -race rc 0 at 9edf1162. Linux/bubblewrap: UNMEASURED-LOCAL.

## Mutants (one at a time, `-run TestExecProbe`, restored with cmp)
| Mutant | Change | Result |
|---|---|---|
| A | linux success-token check -> `if false` | KILLED: TestExecProbeArm8FailsUnsandboxed/linux-no-token |
| C-single | drop `checked.err != nil` | SURVIVES, EQUIVALENT (see below) |
| C-double | keep only `*checked.intact` | PANIC (nil deref) in linux-checker-error; not counted as a kill |
| E | drop `!absent(target)` on darwin creation | KILLED: darwin-create-forged |
| I | `goos=="linux" && target==basename` -> `goos=="linux"` | KILLED: linux-directory-write |
| D | drop `!ok` | SURVIVES, EQUIVALENT |
| N1 | darwin creation leg disabled (`if false`) | KILLED: darwin-create and darwin-create-forged |
| N2 | linux leg drops `*checked.intact` (detection result ignored) | KILLED: linux-detect |
| N3 | `got != "BASE"` -> `got == ""` | KILLED: forged-token |
| N4 | checker given base "HEAD" instead of the captured SHA | KILLED: linux-detect (genuine assertion failure; the first attempt was a compile error and was redone with `_ = base`) |

## Verifying the executor's equivalence claims
- C-single equivalent: CONFIRMED. In checkProtected (exec_protected_check.go) every `result.err = ...` return happens before `result.intact` is assigned. So err != nil implies intact == nil, and the `checked.intact == nil` clause alone still fails the arm. The linux-checker-error row (asserts WorktreeFilterRefused in Why) therefore guards the observable behaviour (a checker error fails the arm, with the filter-refusal text), not the redundant clause. The double drop panics, which is correctly not claimed as a kill.
- D equivalent: CONFIRMED. fileText returns ("", false) for a missing file, and "" != "BASE", so `got != "BASE"` already fails the arm when the sandbox deletes the file.
- B (drop `checked.intact == nil`): equivalent for the same reason as C-single, because the err guard covers it.
- H (plant only "BASE"): equivalent. The control write is only a pre-write proof that the path is writable, not a verdict input.
- F (dedupe `seen` disabled): near-equivalent, since the repeated pattern only plants another file with the same behaviour.

## Judgement on the host-GOOS seam and the skip
- The skip of linux-ebusy-passes and linux-exdev-passes is honest: it is named, carries a reason, and the executor and controller state it. It is a CI-only kill of mutant K (the EBUSY/EXDEV list). I did not measure the rows passing on a linux host, so that is UNMEASURED-LOCAL.
- NON-BLOCKING finding: execFileRefusal reads hostGOOS(), not probeProtected's goos argument. This is a seam inconsistency. The linux oracle is testable on darwin for the checker and token legs, but not for the errno list. Threading goos into the refusal check, or a direct unit assertion on execRefusalCodes("linux"), would make K killable on darwin. That is a cheap follow-up, not a landing blocker.
- Residual NON-BLOCKING items from round 1 remain: arm-wide EBUSY/EXDEV widening (safe, since host-state checks still guard), the rune brute-force cost for high-rune classes (about 345 ms per class), and TestExecSrtStartupProbePasses not carrying a protected set. The darwin real-srt arm-8 test covers the passing path.

No blocking findings. The five spec'd survivors from round 1 (A, E, I plus the equivalent D and C) are resolved. K is covered only on linux CI.

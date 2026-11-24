# Evaluator r3 - iter 429 V4 S4, M5 4fc443cd

## Verdict: PASS. Score 81/100 (unchanged from r2). No BLOCKING findings.
Correctness 27, Tests 23, Security 16, Standards 8, Scope 7 (size now 1,539 by the stated unit; overage is a Linux-CI portability fix plus its test, justified; no change to the r2 trade).

## Diff (verified)
`git diff --stat 976f4028 4fc443cd`: operator_service.go 3+/3-, operator_service_test.go +29. Nothing else.

## (1) Does the reorder change behaviour? No acceptance change.
In `runnerPath` the ownership check is still first; the writable-by-others check (`Perm()&0022`) and the type check (`wantDir && !IsDir || !wantDir && !IsRegular`) are independent predicates that each return an error and nothing between them mutates state. A path is accepted only if both pass, before and after, so the accepted set is identical: no path refused before can now pass. Only the message changes when a path fails both: a wrong-type entry that is also writable (symlink on Linux with mode 0777, a writable non-dir component, a writable non-regular target) now reports "wrong type", where it used to report "writable by others". Existing test `writable tool dir` uses a directory (passes the type check) so still reports "writable by others". The private-mode check is unchanged. The controller's reading is consistent with the code: Linux lstat of a symlink is 0777, which hit the old first check.

## (2) New test non-vacuous on macOS?
Yes for what macOS can pin. R3-TYPE-DROP (`if wantDir && false {`, first edit with `if false` did not compile and is logged COMPILEERR) on the WHOLE `./cmd/ailang-worldd/`: KILLED by TestOperatorServiceOwnership/symlink_tool and TestRunnerPathRefusalClasses. The 0666 regular-file arm pins the "writable by others" class, and the valid target-directory arm guards against over-refusal. The ORDER mutant is not killable on macOS (symlink bits 0755 there); I did not run it, because its result is predetermined by the platform. UNMEASURED here; the Linux CI run is the only discriminator.

## (3) Whole package plain at 4fc443cd
`go vet ./cmd/ailang-worldd/` rc=0. `go test -count=1 ./cmd/ailang-worldd/` rc=0, ok 224.7 s (log judge/r3-plain.txt).

## (4) Linux CI verdict
`$HOME/.ailang/state/world-iter429/ci_pr651_head4fc443cd.verdict` did not exist when I finished: UNMEASURED. The order question (the actual cause of the Linux failure) remains unconfirmed until that run reports. Not blocking on my side, since the reorder cannot widen acceptance and the macOS result is green.

## Ledger
`$HOME/.ailang/state/world-iter429/judge/mutants.tsv` (R3-PLAIN-WHOLE, R3-TYPE-DROP x2). Eval worktree clean (`git status --short` empty) at detached 4fc443cd.

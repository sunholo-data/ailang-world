# Executor report: w-exec-refusal-check-134-substring (row 234, iteration 320)

**Worktree:** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter320-f15`
**Branch:** `sprint/world-iter320-f15-tempdir`
**Base:** HEAD `391ef1aa` (the plan commit) on top of `5407c160` (the design commit). The Go tree
equals the design's measurement base `785b3bd1`: `git diff --stat 785b3bd1 HEAD -- host` is empty.
**Plan:** `design_docs/planned/w-exec-refusal-check-134-substring-sprint-plan.md`
**Design:** `design_docs/planned/w-exec-refusal-check-134-substring.md`

## Base sha256 check (plan §3)

| File | Expected | Observed |
|---|---|---|
| `host/broker/exec_tools_test.go` | `20179c17…e63` | match |
| `host/broker/exec_tools_round2_test.go` | `fe2deb2a…d7` | match (unchanged) |
| `host/broker/handlers_exec.go` | `d8545725…d0` | match (unchanged) |

Pinned toolchain: `go1.26.6 darwin/arm64`; `AILANG v0.41.0` (`$HOME/.pinned-ailang/ailang`).
`GOTMPDIR` was unset in the ambient shell and set inline only by the AC commands.

## Milestone M1: T2 red-first, the fix, T1, T3, AC-A, AC-A-ROOT, mutation drill

| Step | Command | rc | Result |
|---|---|---|---|
| AC-B (M0) | `go test ./host/broker/ -run '^TestExecToolF15RefusalUnderTempRoot134$' -count=10 -v` | 1 | 10 FAIL, 0 PASS, 0 "temp root not honoured". First line: `not a pre-launch refusal: {"ok":false,"refused":"exec node /private/var/folders/…/w234-x134x-…/TestExecToolF15RefusalUnderTempRoot134…/001/node …` (the precondition held) |
| G1 | `gofmt -l host/broker/` | 0 | empty output |
| G2 | `go vet ./host/broker/` | 0 | clean |
| G3 | `go test -run '^$' ./host/broker/` | 0 | `ok … [no tests to run]` (compile fence) |
| Diff | `git diff --stat` | - | `exec_tools_test.go | 95 (+91 −4)` |
| Fixed sha256 | `shasum -a 256` | - | `475fb831692647b87b00690c1142ae2aa6efc9aa8e41df236c2ac2682e039be6`, equal to the planner's prototype (byte-identical edits) |
| AC-A, `tmp-x134x` | `GOTMPDIR=/tmp/w320x/tmp-x134x go test ./host/broker/ -run "$RUN9" -count=20 -v` | 0 | top-PASS 180, sub-PASS 320, FAIL 0, SKIP 0, names 9 |
| AC-A, `tmp-ctl` | the same with `GOTMPDIR=/tmp/w320x/tmp-ctl` | 0 | top-PASS 180, sub-PASS 320, FAIL 0, SKIP 0, names 9 |
| AC-A-ROOT | prebuilt `broker.test`, `GOTMPDIR` = mode-555 root, `-test.run` six callers, `-test.timeout=2m` | 1 | elapsed 1 s; top=6, all=11, leaf=10, mkdir=10, timeout=0, run=11 |

Mutation drill (AC-C). The snapshot was taken after G1–G3 and AC-A; `shasum -a 256 -c` printed `OK` after every restore against the fixed sha `475fb831…039be6`.

| Mutant | Edit (one line, exact plan form) | Named killer rows | Result | Restore |
|---|---|---|---|---|
| M1 MUT-PATH-SUBSTRING | `if after != before \|\| strings.Contains(string(out), "134")` | T1 `refusal-temp-root-134`, T1 `trailing-whitespace`; T2 | killed: both T1 rows FAIL, T2 FAIL (rc 1) | OK |
| M2 MUT-NO-SPAWN-GUARD | `if false && after != before` | T1 `post-launch-f15`; T3 `broken-node`; control T3 `command-aborted` PASS | killed: `post-launch-f15` FAIL, `broken-node` FAIL, `command-aborted` PASS (rc 1) | OK |
| M3 MUT-LOOSE-DECODE | `DisallowUnknownFields()` line deleted | T1 `refusal-plus-exit-code`; control T1 `launched-exit-134` PASS | killed: `refusal-plus-exit-code` FAIL, `launched-exit-134` PASS (rc 1) | OK |
| M4 MUT-OK-IGNORED | `if false && (r.OK == nil \|\| *r.OK) {` | T1 `ok-true` | killed (rc 1) | OK |
| M5 MUT-REFUSED-OPTIONAL | `if r.Refused == nil {` | T1 `refused-empty` | killed (rc 1) | OK |
| M6 MUT-ACCEPT-ALL | `if true { return "" }` inserted above the decoder | T1's 7 reject rows; T3 `broken-node` and `command-aborted` | killed: all 7 reject rows FAIL, both T3 FAIL (rc 1) | OK |
| M7 MUT-NO-EOF-CHECK | `if false && err != io.EOF` | T1 `concatenated-json`, T1 `trailing-garbage` | killed: both FAIL (rc 1) | OK |

Every killer rc was 1, with the named rows red and the named controls green. There were no build
failures and no protocol errors. The file was restored from `/tmp/w320x/exec_tools_test.go.snap`
after each mutant.

## Milestone M2: AC-D and evidence

| Leg | Command | rc | Result |
|---|---|---|---|
| plain | `GOTMPDIR=/tmp/w320x/tmp-x134x go test ./host/broker/ -count=1` | 0 | `ok  github.com/sunholo-data/ailang-world/host/broker  234.810s` |
| race | `GOTMPDIR=/tmp/w320x/tmp-x134x go test -race ./host/broker/ -count=1` | 0 | `ok … 308.852s`; `DATA RACE` count 0 |
| full gate | `go build ./... && go test ./...` (ambient GOTMPDIR unset) | 0 | 35 `ok` packages, 0 `FAIL` lines; `host/broker` ok 328.977s |
| verify_ail | `./scripts/verify_ail.sh` | 0 | `verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`; `PUB011 ratchet passed` |

No failure occurred in any leg, so no known-flake re-run was needed. No `exec_tools*` test failed.

Post-M2 checks: `shasum -a 256 -c /tmp/w320x/fixed.sha` = `OK`;
`exec_tools_round2_test.go` and `handlers_exec.go` unchanged (sha256 as in §Base);
`git diff --name-only` = `host/broker/exec_tools_test.go` only; no `design_docs/world-mission*` file
was touched.

## Deviations

1. **Commits made by the executor.** Plan §1 constraint 1 and the sprint JSON say no git writes by
   the executor (the controller commits). The launching brief asked for two commits (M1 and M2),
   so the executor committed, with the named files only and the Co-Authored-By trailer. The
   controller should check this on the branch.
2. **Base HEAD.** The plan says base HEAD `5407c160`; the worktree was at `391ef1aa` (the plan
   commit on top of it). The Go tree is the same as `785b3bd1` (see above), so the base is
   unchanged in substance.
3. **AC-D background driver stopped at its time limit.** The first driver ran plain, race and the
   full gate in one background job with a 600 s limit. Plain (235 s) and race (314 s) completed
   with rc 0. The job was stopped during the full gate, which had produced no output. The full gate
   and `verify_ail.sh` were rerun alone, from scratch, with no code change in between
   (full 438 s rc 0; verify_ail 47 s rc 0). The plain and race logs come from the first run.
4. **Full-gate package count.** The planner's `git archive` prototype reported 32 `ok` packages; the
   real worktree reports 35 `ok` packages, 0 `FAIL`. The plan's expectation (no failure outside
   `exec_tools*`) holds.
5. **Plan §M1.2 text vs file.** The fixed file is byte-identical to the planner's prototype
   (sha256 `475fb831…039be6`), so no edit deviates from the plan's code blocks.

Nothing in the design was found false. Design §1 (the `134` substring guard) and the fix in §2 were
applied as written.

## Evidence files (this directory)

`executor-acb-m0.txt`, `executor-aca-tmp-x134x.txt`, `executor-aca-tmp-ctl.txt`,
`executor-acaroot.txt`, `executor-g3.txt`, `executor-mut-M1a.txt` … `executor-mut-M7.txt`,
`executor-acd-plain.txt`, `executor-acd-race.txt`, `executor-acd-full.txt`,
`executor-acd-full-summary.txt`, `executor-acd-ail.txt`, `executor-acd-status.txt`.

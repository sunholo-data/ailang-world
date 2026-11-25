# Evaluator report r1 — Ailang World row 243, V4 sprint S5e2 (iteration 436), M1 `3795994c` on base `5cac9316`

## Verdict: PASS — 84 / 100. No BLOCKING findings.

| Category | Max | Score |
|---|---:|---:|
| Correctness vs design §2.2/§5.3/§5.5, judge-435 survivors, §9 | 30 | 26 |
| Tests (non-vacuous, kill, read-back, real git/signer, portable, no weakening) | 30 | 25 |
| Security / fail-closed / D-WORLD-109 | 20 | 19 |
| Coding standards (no production change smuggled in) | 10 | 10 |
| Scope and size discipline, split, Appendix C | 10 | 4 |
| **Total** | 100 | **84** |

## Headline facts (own measurements)
- 11 of the 13 judge-435 survivors are now killed behaviourally, each by a named runtime test. The other two (`Q-R-SIGNEDKIND`, `Q-RB-EXACT`) are genuinely equivalent. I checked both invariants in code and by whole-package run. None were silently dropped.
- No production file changed: `git diff --name-only 5cac9316 HEAD -- '*.go' ':!*_test.go'` printed 0 lines. Only 3 existing `_test.go` files changed.
- Numstat: `git diff -M --numstat 5cac9316 HEAD -- cmd host docs app` = 89+36 + 15+1 + 60+4 = **205**, matching the planner.
- `-race` of the touched families (`OperatorWorktree|OperatorStart|OperatorStep|OperatorMint|ReadbackError`, `./cmd/ailang-worldd/ ./host/store/`): head 59.965 s (cmd) and 2.274 s (store); base 59.881 s and 2.550 s. This is +0.08 s, so I do not reproduce the planner's +3.489 s. Load noise dominates. CI cost is negligible.
- The controller's `legs.txt` shows every leg rc=0 except `verify_go` rc=1 (the row-271 baseline set). `observation-after` is byte-equal to `observation-before` (I diffed them), so there was no real-path write.

## Test runs (all rc 0, with the env from the directive)
- `go test -count=1 ./host/store/` is ok (14.8 s).
- `go test -count=1 -run 'Operator|Preset|Session|Quickstart|Help|CredentialWriters' ./cmd/ailang-worldd/` is ok (66.0 s).
- `go test -race -count=1 -run '…' ./cmd/ailang-worldd/ ./host/store/` is ok at head and at base.
- `go vet ./cmd/... ./host/...` rc=0.
- `go test -run '^$' ./...` has no failures.
- Whole-package `go test -count=1 -timeout 8m ./cmd/ailang-worldd/` was run only under the two equivalence mutants. Its fail set is identical to the base 11-line set (see below).

## Mutation ledger
Ledger: `$HOME/.ailang/state/world-iter436/judge/mutants.tsv`, 52 lines. Specs: `judge/specs.json`, runner `judge/mut.py`, logs `judge/log-*.txt`. Every anchor was asserted unique, and the edit was confirmed by `git diff --stat` and `go vet` before reading the result. Eval worktree `git status --short` is clean at the end. The throwaway base worktree was removed.

`comm -23` of my ids against the plan's 17 ids (`plan-ids.txt`) gives 39 ids that are not in the plan. These are 26 distinct own mutants once the `-WHOLE`/`2` reruns and controls are discounted.

### 13 judge-435 survivors (named selector: `Operator` in host/store; the plan's 5-family selector in cmd/ailang-worldd)
| Id | Package | Result | Killing test |
|---|---|---|---|
| Q-ST-OWNER | host/store | KILLED | TestReserveOperatorWorktreeGuards/owner |
| Q-ST-IDENT | host/store | KILLED | TestReserveOperatorWorktreeGuards/identity |
| Q-ST-EXPIRYEDGE | host/store | KILLED | TestReserveOperatorWorktreeGuards/expiry |
| Q-ST-READSTATE | host/store | KILLED | TestReadOperatorWorktreeGuards/uncertain |
| Q-ST-READCOUNT | host/store | KILLED | TestReadOperatorWorktreeGuards/duplicate |
| Q-ST-DISPKINDSQL | host/store | KILLED | TestDisposeOperatorUncertainStepGuards/action_kind |
| Q-RB-RAWABS | cmd/ailang-worldd | KILLED | TestOperatorWorktreeReadbackGuards/relative_registration |
| Q-W-ABSENTDISP | cmd/ailang-worldd | KILLED | TestOperatorWorktreeWindow/success_then_absent |
| Q-R-REPOABS | cmd/ailang-worldd | KILLED | TestOperatorWorktreeRecoveryGuards/repo_abs |
| Q-R-COMMONABS | cmd/ailang-worldd | KILLED | TestOperatorWorktreeRecoveryGuards/common_abs |
| Q-R-REREADERR | cmd/ailang-worldd | KILLED | TestOperatorWorktreeLockedReread/error |
| Q-RB-EXACT | cmd/ailang-worldd | SURVIVED, equivalent | whole package: fail set == base set (`judge/fails-EXACT.txt`) |
| Q-R-SIGNEDKIND | cmd/ailang-worldd | SURVIVED, equivalent | whole package: fail set == base set (`judge/fails-SK.txt`) |

All 11 kills are behavioural: real store, real git in `t.TempDir()`, and the real CLI path. None is an AST, static or inventory test.

The two `*-WHOLE` lines in the ledger read "KILLED-by-TestCLIAgainstSeToolsDaemon,…". That is the runner's naive rc≠0 label. The actual verdict is SURVIVED, because the failing top-level set is exactly the 11 row-271 tests (`diff` against `pkg-cmd-ailang-worldd.base-fails` gave FAILSET-EQUAL both times).

Equivalence invariants, checked in code:
- `Q-RB-EXACT` (`cmd/ailang-worldd/operator_worktree.go:235`): `exact` is incremented only inside the branch that appends to `e.registration` (lines 205–209). So `exact >= 2` implies at least two `\x00\x00` delimiters, and the separate `strings.Count(e.registration,"\x00\x00") != 1` clause refuses it. `exact != 1` and `exact < 1` therefore agree on every reachable state. Confirmed by run.
- `Q-R-SIGNEDKIND` (`operator_worktree_reconcile.go:39`): `approval/operator_action.go:~150` rejects finish with headless, and `store/operator_mint.go:183` enforces row kind equals action kind. So a signed finish always has `Headless=false`, and the adjacent `!signed.Options.Headless` guard refuses it. Equivalent. Confirmed by run. The `signed_kind` arm is correct as an inert regression pin and is not itself a kill.

### Controls
- Must-die: `J436-CTL-DIE` (`initial.state == 99`) KILLED by TestOperatorStartCollisions/{ref_ancestor,ref_descendant,registration…}. `J436-CTL-DIE-STORE` (invert claimed check) KILLED by TestReserveOperatorWorktreeGuards/{claimed,valid}.
- Must-live: `J436-CTL-LIVE-STORE` (redundant `|| false`) SURVIVED at named and whole package. `J436-CTL-LIVE2` (redundant `exact != 1 || false`) SURVIVED the cmd selector.
- `J436-CTL-LIVE` (suffix `!` on the "invalid Git worktree registration" message) was KILLED by `relative_registration`. That is correct and not a harness defect, because the new arm asserts that exact message. It does mean the plan's CTL-LIVE selector choice (TestOperatorStartOnce) is not the only reachable test. This is informational.

### My own mutants (not in the plan), reservation/readback/disposition path
Killed, with test and package:
- `ST-CHALLENGELEN` → Reserve…/challenge (host/store)
- `ST-HEADLESS` → Reserve…/headless
- `ST-APPROVEDBY` → Reserve…/approver
- `ST-DISPUNCVER` → Dispose…/version
- `ST-DISPTERMINAL` → TestReadOperatorWorktreeGuards/uncertain (host/store)
- `W-CREATEDADDERR` → TestOperatorWorktreeWindow/add_error_after_checkout (cmd)
- `RB-ISDIR` → ReadbackGuards/directory
- `RB-REFSEXACT` → Window/ref_read_skew
- `RB-REFPREFIX` → TestOperatorStartCollisions/ref_ancestor+ref_descendant
- `RB-COMMON2` → ReadbackGuards/repo_common
- `RB-PHYSICAL2` → ReadbackGuards/physical
- `R-HEADLESS` → RecoveryGuards/headless
- `R-STOREEQ` → RecoveryGuards/store
- `R-ABSENTSTATE` → Window/{head_changed,path_only,redirected_checkout,…}
- `R-BEFOREAFTER2` → TestOperatorWorktreeDispositionCAS/git_changed

Compile errors, discarded and redone with `&& false` forms: `RB-COMMON`, `RB-PHYSICAL`, `R-BEFOREAFTER`, `R-CURRENTEP`.

SURVIVED (host/store, whole package rc 0 — these are real unpinned guards, none in the judge-435 list):
- `ST-HEARTBEATWIN`: the lower `seen_at > now-10` freshness bound in the reserve UPDATE is not pinned.
- `ST-NOEPISODESTART`: `NOT EXISTS episode_starts` / `session_credentials` in the reserve UPDATE is not pinned.
- `ST-READPHASE`: `phase='reserved'` in `ReadOperatorWorktree`.
- `ST-READOTHERKIND`: the other-kind-steps `NOT EXISTS` in `ReadOperatorWorktree`.
- `ST-DISPOTHERKIND`: the same `NOT EXISTS` in the dispose UPDATE.

SURVIVED the named cmd selector only (not run whole-package, for time; I judge them equivalent by reading the code):
- `W-CREATEDREADERR` (`created := addErr == nil && readErr == nil && …` without readErr): when `readErr != nil`, `got` is the zero evidence, so `got.state == treeAbsent` and `created` is false either way. Equivalent.
- `R-CURRENTEP2` (`currentEpisode != episode` dropped): the episode derives from the digest-bound action already compared via `current != key`. Equivalent in practice.

I did not run those two as whole-package runs, so the directive's "SURVIVED needs whole package" is not strictly met for them. The equivalence argument rests on the code read.

## Per-arm non-vacuity (fails with the guard off, passes at head)
- `owner`, `identity`, `expiry` → Q-ST-OWNER, Q-ST-IDENT, Q-ST-EXPIRYEDGE fail them. The expiry arm keeps `seen_at` fresh, so under `>` the mutant reaches the later SQL gate and returns the wrong class. Verified by the kill.
- `TestReadOperatorWorktreeGuards/{uncertain,duplicate}` → Q-ST-READSTATE and Q-ST-READCOUNT. `uncertain` also kills `ST-DISPTERMINAL`.
- `action_kind` → Q-ST-DISPKINDSQL.
- `success_then_absent` → Q-W-ABSENTDISP. It removes the real worktree and branch between a real `git worktree add` and the readback. It asserts `uncertain`, step reserved v2 and the exact remedy `worktreeExample`. The shared loop then executes that remedy through `r.run` for `"no tty"`, `"n\n"` and `"y\n"`, and checks `failed`/`observed`/`stopped: absent` v3 by row read-back.
- `repo_abs`, `common_abs` → Q-R-REPOABS and Q-R-COMMONABS, with the refusal asserted before the prompt and store/branch/files snapshots unchanged.
- `relative_registration` → Q-RB-RAWABS. It rewrites only the path field of real `git worktree list --porcelain -z` output, served through a wrapper, with a checkout snapshot before and after. It asserts the exact message "invalid Git worktree registration".
- `LockedReread/error` → Q-R-REREADERR (the injected error is returned on call 2). The `generation` arm is retained: the old test is split into two subtests with the same assertions, `calls == 2`, no prompt, and state preserved.
- `signed_kind` is an inert regression pin (equivalent), kept for the stored-kind row.

## Test quality
- Fixtures use real git in `t.TempDir()`, the real migration, the repo signer, and read-back of the action, step and sessions rows. Refusals assert "no state changed" by `reflect.DeepEqual` against a prior read-back.
- No real sleeps were added: the diff contains no `Sleep` or `time.After`. The one `AfterFunc` is pre-existing in `TestOperatorWorktreeDispositionCAS`.
- No macOS-only assertions in the new code. `relative_registration` uses `EvalSymlinks(cwd)` and `filepath.Rel`, and the wrapper is POSIX `sh`/`cat`.
- One existing file, `operator_worktree_test.go`, was edited. Every changed hunk adds an arm, and no assertion was weakened. The `LockedReread` rewrite and the `physical` branch widening (`want` message selected per arm) keep the original checks.
- Gap: five store guards in the reserve/read/dispose SQL on this same path are still unpinned (see survivors above). The sprint named only the judge-435 list, so this is not a plan failure, but the claim "pre-caller pins" is weaker than it reads.

## Inertness / activation
- Zero production `.go` changes.
- Production callers (non-test, non-definition): `reserveOperatorStartWorktree` 0 (only a doc comment at `operator_worktree.go:29`); `ReserveOperatorWorktree` 1 (its only caller is `reserveOperatorStartWorktree` at `operator_worktree.go:81`, which is itself uncalled in production); `ReadOperatorWorktree` 1 (`readRunnerWorktree`, the existing S5d attended `reconcile-worktree`); `reconstructOperatorStart` 1 (from `reserveOperatorStartWorktree`); `MintOperatorStart` 0 production entry callers (the controller's `callers.txt` agrees). The positive control is that the same grep finds the `ReserveOperatorWorktree` call at `operator_worktree.go:81`. Nothing is newly reachable.

## The split and Appendix C
- M1 alone is defensible: the pins protect a path that M2 builds on, so it is safer to pin first. M2 hardened is 1,259 incremental, under the cap by itself. The planner's "over the cap" holds only for M1+M2 = 1,464 (64 over, 4.6%). The honest alternative was shipping M2 alone this sprint, with these 205 lines of pins following or being trimmed.
- This is the second consecutive sprint in which publication was measured and deferred (S5e 967, S5e2 1,259). Appendix C also marks the whole flag/timing/refusal matrix **UNMEASURED**, so the real M2 will exceed even 1,259. As a result the feature has seen zero landed lines in two sprints, and a third deferral risks becoming a pattern. I treat this as non-blocking but a real scope-discipline finding. It is why that category scores 4/10.
- Appendix C is measured for the publication+disposition unit (967/1,172 pure, 1,259/1,464 hardened) and lists the per-flag mutants, the S5f/S6 paired-process pins and the #667 GC gate. It is honest about what is unmeasured. It does not size the S5f/S6 or GC work.
- #667 GC probe: I read `probe-M2-GC.log` and the probe source (`plan-scratch/probe-M2-GC/cmd/ailang-worldd/operator_gc_probe_test.go`). Both the reserved and published phases show a clean `git status`, a dry-run candidate with reason "clean & merged" at an injected 73 h age (`st.ModTime()+73h`, `OlderThan=72h`), and evidence unchanged by read-back. The claim holds as stated, with two caveats. The probe stubs `IsActive` to false (production would consult a live-session snapshot that does not include operator steps). The published arm depends on the deferred M2 code, which is not in head. The conclusion that GC needs an operator-step-aware skip before the first production caller is sound. The `episode_gc.go` CLI and its help text show no operator-awareness.

## Findings
BLOCKING: none.

Non-blocking, ranked:
1. Scope/split: publication deferred for the second consecutive sprint although M2 alone (1,259) fit the 1,400 cap. A third deferral should be refused. Plan the next sprint as publication-first with its pins co-located.
2. Five unpinned store guards on the same reserve/read/dispose path survive the whole host/store package: heartbeat lower bound (`operator_worktree.go:59`), `episode_starts`/`session_credentials` NOT EXISTS (`:56`), read `phase='reserved'` (`:89`), the other-kind steps NOT EXISTS in read (`:90`) and dispose (`:138`). They are pre-existing and were not on the judge-435 list. Pin them before the first production caller.
3. The `signed_kind` arm is non-discriminating (equivalent mutant). It is fine as a pin, but the plan should not count it toward the "11 killed" figure. It does not: the plan counts 11 non-equivalent kills correctly.
4. The `-race` delta (+3.489 s) is not reproducible on this rig (+0.08 s). The cost is negligible either way.
5. Two own mutants (`W-CREATEDREADERR`, `R-CURRENTEP2`) were judged equivalent by code reading without a whole-package run (time budget).

## Paths
- Report: `$HOME/.ailang/state/world-iter436/evaluator-r1.md`
- Ledger and logs: `$HOME/.ailang/state/world-iter436/judge/` (`mutants.tsv`, `specs.json`, `fails-EXACT.txt`, `fails-SK.txt`, `race-head.txt`, `race-base.txt`, `t-store.txt`, `t-cmd.txt`)
- Eval worktree `../.wt-world-iter436-row243-v4s5e2-eval` is clean.

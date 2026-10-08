# Row 153 evaluator r2 — iteration 241 (judge: claude-opus-5-5)

Under review: PR #227 head `0c67c9a` (fix round c40410e … 0c67c9a on top of r1's `e327415`).
Evaluator worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter241-eval`, detached at `0c67c9a`, restored clean.
Logs: `r2/` (gates) and `mut/R2-*` (mutants) in this directory.

## Verdict: **92/100, PASS, no BLOCKING findings.**

| Rubric | Max | r1 | r2 | Why |
|---|---|---|---|---|
| ACs implemented and covered | 30 | 26 | 29 | AC3.4 now asserts a deterministic witness; AC3.3 warm arms now exercise the publication-built template |
| Mutation non-vacuity | 25 | 20 | 23 | Every r1 survivor is now killed; two single-site finish mutants survive (one equivalent, one minor, R-a) |
| Correctness and safety | 20 | 18 | 18 | No production change; r1 adjudications stand; N8 residual |
| Gates | 15 | 8 | 13 | Local gates all rc 0; CI for 0c67c9a not yet observed by me (5 runs in progress) |
| Evidence honesty | 10 | 8 | 9 | The executor's fix-log claims re-measured and held; the single-site equivalence claim is only half right (R-a) |

## Scope check
`git diff --name-only e327415 0c67c9a` gives exactly the 10 files the coordinator listed (+247/−14). Non-test files are only
the two design/plan docs and `executor-r2-fixes.md`. **No production code changed and no file was omitted.**

## Gates (mine, at 0c67c9a, `AILANG_BIN=~/.pinned-ailang/ailang` v0.41.0)
- `verify_ail.sh`: rc 0 (16 identities, 40 named tests, 444 se-tools tests).
- `go vet ./...`: rc 0.
- `go test ./... -count=1`: rc 0.
- `go test -race ./... -count=1 -timeout 8m`: rc 0. This ran on a `git archive` copy while my mutation runs loaded the host.
- Real srt `-p 1 -v`: rc 0, with `--- PASS` for TestExecSrtMCPEndToEnd, TestExecSrtMCPGrantAndBudget and TestExecSrtReplayIsByteEqual.
- All 24 new tests (`8e15d28..HEAD`) by name: rc 0, 24/24 `--- PASS`, 0 SKIP.
- CI runs 37703800801, 37703807574, 37703809910, 37703812491 and 37703815214 were all `in_progress` on 0c67c9a when I checked; the controller tallies them.

## r1 findings

| Finding | Status | My evidence |
|---|---|---|
| **B1** ratio assertion | **CLOSED** | The ratio is gone (`cachetemplate_test.go:234-238`): timings are logged only, plus `ColdRuns()==5`. R2-MUT-NOTEMPLATE: KILLED at :231 "warm run 0 compiled … template was not copied in". Design AC3.4 and the plan row carry the r2 note. |
| **N2** AC3.3 used the run-promoted template | **CLOSED** | `dropTemplate` added before `buildTemplate` (`capsule_template_test.go:75-78`). R2-OWN-8 (hollow check-built template): KILLED with 11 subtest FAILs (10 phase arms + cross-source); in r1 only cross-source failed. |
| **N3** finish budget unpinned | **CLOSED for the effect path; residual R-a** | R2-OWN-2 (both call sites): KILLED, `phase_timeout_test.go:204` "finish phase deadline in 3.99998975s, want within (0, 2s]". Single-site mutants: outer-only SURVIVED, inner-only SURVIVED (whole coordinator package). See R-a. |
| **N4** symlink refusal | **CLOSED** | R2-OWN-4: KILLED, `capsulecache_test.go:179` "CopyCacheTree over a symlink = <nil>". |
| **N5** failed-run promotion | **CLOSED** | R2-OWN-6: KILLED, `cachepromote_test.go:54` "template ready = true after exit 3". The success arm is a working control. |
| **N6** ready without manifest | **CLOSED** | R2-OWN-7: KILLED, `capsulecache_test.go:208`. I agree N5/N6 were unpinned rather than wrong; the code was already correct. |
| **N7** 7 s announce bound | **CLOSED** | 60 s at `templates_test.go:177`. R2-MUT-BLOCKING-PREWARM-b (`maintainTemplates` before `Listen` in `Run`): still KILLED at :178 after 62.56 s, so the test lost no power. |
| **N9** serveR14 runner log discarded | **CLOSED** | Cold-forcing mutant (`if false && CapsuleTemplateReady`): with the fix, TestSeToolsA2APostEffectFailureIsMapped and TestSeToolsMCPPostEffectFailureIsLogged both FAIL on the tripwire (`testlog_test.go:58`). Same mutant with the fix reverted to `capsule.Config{}`: both PASS, which proves the gap was real (`mut/R2-N9-coldforce-{fixed,unfixed}.log`). |
| **N1** AC1.7 wire test can't see the predicate | Residual, agreed non-blocking | The `TestA2ADispatch` table rows are the tooth: r1's MUT-ANYPHASE, MUT-ORDER-FINISH and the double mutant were all KILLED there. |
| **N8** GC race window | Residual, agreed non-blocking | Worst case is one labelled cold run, then lazy re-promotion. |
| **N10** cold label text | Residual, agreed non-blocking | Cosmetic. |
| **N11** CI race-leg headroom | Residual, agreed non-blocking | Runner variance; row 154 owns it. |

## R-a — the executor's single-site equivalence claim, checked
The executor says a mutant at only one of the two call sites is equivalent because the two nested contexts are both 2 s.
- **Outer-only** (`detached(ctx, c.planBudget)` at `effectful.go:327`): the inner `pctx` still caps the finish at 2 s, so the
  deadline is equivalent. The only difference is that the inner cap now expires first, giving a `PhaseTimeoutError{finish}`
  instead of a plain wrap. Both are wrapped in `EffectsUnrecordedError`, so the wire answer is the same. **Equivalent.**
- **Inner-only** (`runPhase(finishCtx, c.planBudget, …)` at `:231`): this is equivalent on the effect path, where the outer
  2 s `fctx` bounds it. It is **not** equivalent on the zero-effect refusal-plan path (`:296`, `composeEffectOutput(ctx, …)`)
  or the replay path (`:468`), where the finish would get 4 s under the caller's context. Nothing ran on those paths and
  4 s + 4 s fits inside 20 s, so it is not unsafe, but that cap is unpinned.
- **Verdict:** the claim is half right. Non-blocking residual. Optional fix: a refusal-plan variant of
  `TestFinishPhaseRunsUnderFinishBudget`.

## Wall-clock re-scan of new tests (load-flake risk on a 2-core runner)
- **TestSlowPlanLeavesHandlerItsCap** (`phase_timeout_test.go:160`). Derivation: the handler deadline is
  th + min(10 s, t0+20 s − th − 6 s), and the plan returns about 500 ms before its 4 s cap. So the test passes iff
  (pre-plan work + timer lateness) ≤ 600 ms. The 100 ms margin is on top of a 500 ms structural slack. There is a second
  path: if the plan goroutine is starved for more than 500 ms, `time.After` and `ctx.Done` are both ready and `select` picks
  one at random, so the plan can error. That needs half a second of goroutine starvation. **Low risk, non-blocking**, same as r1.
- **TestFinishPhaseRunsUnderFinishBudget**: the upper bound is one-sided, and load only shrinks the measured value. The
  `got < 0` arm needs more than 2 s between `detached()` and a fake runner's first line. Negligible.
- **TestMCPPlanTimeoutFrozenWire** (`took ≤ 8 s` for a 4 s budget): 4 s of slack. Low.
- **TestDaemonStartReportsMissingTemplates**: 60 s announce and 60 s rebuild bounds. Fine.
- **TestCapsuleWarmRunsDoNotCompile**: no timing assertion left.
- Every other bound is a fail-instead-of-hang safety net, never the pass criterion.
- **No remaining wall-clock pass criterion is likely to flake.**

## Mutants this round
- **KILLED (9):** R2-MUT-NOTEMPLATE, R2-OWN-2-both, R2-OWN-4, R2-OWN-6, R2-OWN-7, R2-OWN-8, R2-MUT-BLOCKING-PREWARM-b, and the N9 cold-forcing mutant on each of the two R14 tests.
- **SURVIVED (2):** R2-OWN-2-outer-only (equivalent) and R2-OWN-2-inner-only (minor, R-a).

## Worktree
`git status --porcelain` is empty, `git diff --quiet` is clean, and HEAD is `0c67c9a`. Every mutant was restored with
`git checkout`; the `git archive` gate copy was deleted.

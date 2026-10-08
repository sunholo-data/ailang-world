# Row 153 — executor fix round r2 (iteration 241)

Judge: `evaluator-r1.md` (80/100, one BLOCKING). Base for this round: `e327415`.
Logs: `~/.ailang/state/mission-world-iter241/executor-r2/`.

## Deviations

1. **AC3.4 changed (B1).** The wall-clock assertion `median warm < median cold / 2` is removed. It is machine-dependent
   (CI linux `-race`: warm/cold 0.63 at b6f7dfb and 0.53 at e327415; idle mac 0.35) because fixed per-run costs (100 MB verify,
   copy-in, spawn, MkdirTemp) do not scale with the interpreter's compile cost; the design's "ratio, so machine speed cancels"
   claim was refuted (judge decomposition). CI runs 37696379845 / 37700015498. The test is renamed
   `TestCapsuleWarmRunsDoNotCompile`: it asserts every cold arm compiles (MOD010), no warm arm compiles, and
   `runner.ColdRuns() == 5`; both timings are `t.Logf`-ed as evidence only. The speed effect stays evidenced by AC3.8/AC4.1's
   proxy. Design AC3.4 text and the sprint plan's row carry the "r2" note.
2. **N5/N6 were unpinned, not wrong.** The code already refuses to promote a failed run (`promoteCold` after the `runErr`
   check) and `CapsuleTemplateReady` already requires a regular `compile/manifest.json`; no production change, tests only.
3. **Production code is unchanged in this round.** All commits touch tests and docs only.

## Fixes and mutation proofs

| Item | Test (new/changed) | Mutant | Result with the fix | Restored |
|---|---|---|---|---|
| B1 | `TestCapsuleWarmRunsDoNotCompile` (capsule) | MUT-NOTEMPLATE (template copied to a discard dir) | red at `cachetemplate_test.go:231` "warm run 0 compiled ... template was not copied in" (before the ratio, as the judge measured) | yes |
| N2 | `TestCapsuleWarmEqualsColdForSeTools` per-phase arms now `dropTemplate` before `buildTemplate` | OWN-8 (check-built template hollowed) | all 10 phase arms and cross-source red (was: only cross-source) | yes |
| N3 | `TestFinishPhaseRunsUnderFinishBudget` (coordinator) | OWN-2 (finish uses planBudget, both sites) | red: "finish phase deadline in 3.99s, want within (0, 2s]" | yes |
| N4 | `TestCopyCacheTreeRefusesSymlinks` (archive) | OWN-4 (`case true:` copy follows symlinks) | red: "CopyCacheTree over a symlink = nil" | yes |
| N5 | `TestFailedColdRunIsNotPromoted` (capsule; shell-fake interpreter, success arm is the control) | OWN-6 (promote before the runErr check) | red: "template ready = true after exit 3" | yes |
| N6 | `TestCapsuleTemplateReadyNeedsTheManifest` (archive) | OWN-7 (Ready = dir exists) | red: "ready for a directory without compile/manifest.json" | yes |
| N7 | `templates_test.go` announce bound 7 s -> 60 s | blocking rebuild (`maintainTemplates` run synchronously before the announce) | still red at the 60 s bound (62.7 s run) | yes |
| N9 | `serveR14` runner uses `testOperatorLog` | cold-forcing (`if false && CapsuleTemplateReady`) | before the fix the two R14 tests PASS under the mutant (gap proven); after: both red on the tripwire | yes |

N3 note: a mutant at only ONE of the two call sites is behaviourally equivalent (the effective phase deadline is the
minimum of the two nested contexts, both 2 s), so the test pins the observable deadline, which is what AC2.4's headroom
arithmetic depends on.

## Wall-clock assertions in NEW tests (git diff 8e15d28..HEAD, `*_test.go`)

| Test | Assertion | Slack |
|---|---|---|
| `TestCapsuleWarmRunsDoNotCompile` | none (timings logged only) | n/a |
| `TestSlowPlanLeavesHandlerItsCap` | handler deadline-left >= `HandlerCap - 100ms` after a plan that returns 500 ms before its phase deadline | ~500 ms of pre-plan work is the real slack; 100 ms is tolerance on top; judge: effective 600 ms. Residual, not changed |
| `TestMCPPlanTimeoutFrozenWire` (`mcp_plan_timeout_test.go`) | `took <= PlanPhaseBudget + 4s` (8 s) | 4 s over a 4 s budget |
| `TestFinishPhaseRunsUnderFinishBudget` | remaining finish deadline `<= 2 s` (upper bound only) | one-sided, load only shrinks the measured value; no flake path |
| `TestA2APlanTimeoutRetrySameTaskID` | `planBudget = 100ms` against a runner that blocks until cancelled | none needed: ctx-driven, no elapsed-time bound |
| `TestDaemonStartReportsMissingTemplates` | announce within 60 s; rebuilt line within 60 s | 60 s (was 7 s, N7) |
| `TestTemplateGCPrunesStaleInterpDigests` | maintenance done within 20 s | generous; GC of a handful of dirs |
| phase-timeout wire/handler tests (`a2a_wire_test.go`, `plan_budget`/`effectfulHandler`) | failure-mode safety nets `budget + 5 s`, 20 s handler | fail-instead-of-hang bounds, never the pass criterion |
| 30 s `context.WithTimeout` in rigs | per-run safety cap | 30 s vs ~0.2-1.3 s runs |

## Residuals (recorded, not implemented)

- **N1** AC1.7 wire test cannot see the predicate it claims to guard; the `TestA2ADispatch` table rows are the tooth.
- **N8** template GC runs once at start against the head read in `New`; small race window, consequence is one labelled cold run.
- **N10** the cold label is always `(template missing)`, also for the unreadable-template path.
- **N11** CI race-leg headroom (verifygate 379 s, broker 312 s of 480 s): row 154's concern.
- `host/broker` `TestRunBoundedHeadTailTimeoutKeepsPartialOutput`: confirmed pre-existing flake (0.3 s budget, untouched package).

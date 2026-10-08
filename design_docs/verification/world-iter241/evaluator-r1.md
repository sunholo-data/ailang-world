# Row 153 evaluator r1 — iteration 241 (judge: claude-opus-5-5; executor: claude-sonnet)

Under review: PR #227 head `e327415` (C1 b5a4a25 … C6 e327415) vs plan base `8e15d28`.
Evaluator worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter241-eval` (detached `e327415`; restored clean, see end).
Logs: this directory (`verify_ail.log`, `vet.log`, `test.log`, `race.log`, `srt.log`, `newtests.log`, `decomp-*.log`,
`warmfaster-race-bg8.log`, `broker-*.log`, `ci-e327415-113060961481.log`, `mut/*.{diff,log}`).

## Verdict

**Score 80/100 — PASS on the rubric, NOT LANDABLE as-is: one BLOCKING finding (B1) reds CI's race leg on 2 of 2 attempts.**
B1 is a one-test fix; nothing else blocks.

| Rubric | Max | Score | Why |
|---|---|---|---|
| AC implementation + coverage | 30 | 26 | All ACs implemented; AC3.4 as specified is wrong (B1); AC3.3 warm arms do not exercise the publication-built template (N2) |
| Test non-vacuity (mutation) | 25 | 20 | Every safety-critical named mutant re-killed by me; 4 own mutants survive (N3–N6); AC1.7 wire test alone does not discriminate (N1) |
| Correctness / safety | 20 | 18 | Retry semantics, template integrity, EXDEV, racer, non-blocking startup all hold; minor GC race / discarded log (N8, N9) |
| Gates | 15 | 8 | Local: all green. CI: race leg red 2/2 (B1) — "nothing lands red" |
| Evidence honesty | 10 | 8 | Claims re-measured and held (broker baseline, MOD010, proxy); executor flagged the ratio margin (F4) but shipped it |

## Gates (mine, at `e327415`, `PATH=/opt/homebrew/bin:$PATH AILANG_BIN=~/.pinned-ailang/ailang` = v0.41.0 `24ee108`)

| Gate | rc | Log |
|---|---|---|
| `./scripts/verify_ail.sh` | 0 ("16 required identities verified, 40 named tests pass, 444 se-tools named tests pass") | verify_ail.log |
| `go vet ./...` | 0 | vet.log |
| `go test ./... -count=1` | 0 (all packages ok, incl. host/broker) | test.log |
| `go test -race ./... -count=1 -timeout 8m` (git-archive copy of e327415) | 0 (capsule 56.3 s, daemon 200.2 s, broker 154.8 s) | race.log |
| real srt `-p 1 -v` (`WORLD_EXEC_SRT_NODE_MODULES=…/srt/node_modules`) | 0; `--- PASS:` TestExecSrtMCPEndToEnd, TestExecSrtMCPGrantAndBudget, TestExecSrtReplayIsByteEqual (0 SKIP) | srt.log |
| 20 new tests by name, `-v` | 0; 20/20 `--- PASS: <name>`, 0 SKIP | newtests.log |
| **CI run 37700015498 (e327415), job 113060961481** | **failure** — only `TestCapsuleWarmIsFaster` in the `-race` leg | ci-e327415-113060961481.log |
| CI run 37696379845 (b6f7dfb), job 113049026570 | failure — same test, same leg | ../pr-b6f-113049026570.log |

Both CI runs: non-race leg all ok; in both parallel legs `host/daemon` ok, so TestExecSrtMCP* (row 161's signatures) passed under the
parallel legs 4/4 legs observed, and passed in the `-p 1` srt step (`--- PASS: TestExecSrtMCPEndToEnd (7.81s)` etc., log lines 440–700).

## BLOCKING

### B1 — AC3.4's wall-clock ratio is not machine-invariant; `TestCapsuleWarmIsFaster` reds CI's race leg 2/2
`host/capsule/cachetemplate_test.go:232-236` (`if mw >= mc/2 { t.Fatalf(...) }`).

Evidence (instrument = CI job logs, banked):
- b6f7dfb race leg: cold [571 727 625 541 682] ms, median 625; warm [333 345 432 413 394], median 394 → warm/cold 0.63.
- e327415 race leg: cold [1227 1110 1007 1089 1056] ms, median 1089; warm [580 580 623 567 602], median 580 → 0.53.
- Both runs: the MOD010 assertions at :222 (every cold run compiled) and :228 (no warm run compiled) did **not** fire — the
  failure is at :235. The template works on linux/amd64; only the ratio fails. Non-race leg passed both times.

Decomposition (my probe, temporary `zz_eval_timing_test.go`, deleted; `decomp-*.log`): a warm run is mostly fixed cost that a
ratio cannot cancel.

| Arm | verify (100 MB sha256) | copy-in | warm total | cold total | warm/cold |
|---|---|---|---|---|---|
| plain, idle (M4 Max) | 38–45 ms | 4–5 ms | 68–74 ms | 197–201 ms | 0.35 |
| `-race`, idle | 53–58 ms | 5–7 ms | 88–90 ms | 214–222 ms | 0.41 |
| `-race`, `taskpolicy -b` + 8 bg burners | 564–579 ms | 88–133 ms | 1.29–1.40 s | 3.27–4.03 s | 0.35–0.42 |
| CI ubuntu amd64 `-race` (2 runs) | — | — | 394 / 580 ms | 625 / 1089 ms | 0.63 / 0.53 |

warm/cold = (F + w)/(F + c) with F = verify + copy-in + spawn + MkdirTemp/RemoveAll. It cancels machine speed only if F, w and c
scale together; they don't (F is hash/I/O in the race-built test binary, c−w is the interpreter's compile CPU). The design's
claim "This is a ratio, so machine speed cancels out" (design AC3.4) is refuted. **The controller's hypothesis is correct**, and
it is exactly the load-fragile-timing class the row exists to remove.

Fix (verified sufficient): delete the `mw >= mc/2` assertion (keep the `t.Logf` of both arms as descriptive evidence). The
deterministic MOD010 witness already kills MUT-NOTEMPLATE before the ratio is reached — re-run by me: `MUT-NOTEMPLATE: KILLED …
cachetemplate_test.go:229: warm run 0 compiled (stderr "WARNING MOD010 …")` (`mut/MUT-NOTEMPLATE.log`). Optionally strengthen with
`runner.ColdRuns()` == 5 after the loop (counts exactly the cold arm). Do NOT replace it with a child-segment-only timing ratio: it
needs a new production seam to expose the child's wall, and it is still a wall-clock comparison under contention. Record the AC3.4
amendment as a deviation citing these four rows; the speed effect stays evidenced by AC3.8/AC4.1's proxy (child 0.50–0.82 s warm
vs 1.96–2.98 s cold) — evidence, not a gate.

## NON-BLOCKING

- **N1 — AC1.7's wire test cannot see the predicate it claims to guard.** `host/projection/a2a_wire_test.go:600-623`. In production
  a finish overrun never produces a `*PhaseTimeoutError`: `fctx = detached(ctx, c.finishBudget)` (`effectful.go:327`) and
  `runPhase`'s `pctx = WithTimeout(fctx, c.finishBudget)` have the same budget, the parent's deadline is earlier, so
  `ctx.Err() != nil` and `classifyPhaseErr` returns the plain wrap. My double mutant (predicate dropped AND case moved above
  `EffectsUnrecordedError`): `MUT-DOUBLE-ANYPHASE-ORDERFINISH: SURVIVED` against `TestA2AFinishOverrunStaysEffectsUnrecorded`;
  `MUT-DOUBLE-table: KILLED` by `TestA2ADispatch` rows `FinishPhaseTimeout_bare` / `EffectsUnrecorded_wrapping_plan_timeout`
  (`projection_test.go:844-849`). Safety holds (the table rows are the tooth); the wire test is a regression guard for the
  unrecorded wrap (MUT-UNWRAP-FINISH), not for the predicate. Fix: correct the test's comment/ledger attribution; optionally a
  wire arm whose finish runner returns `&coordinator.PhaseTimeoutError{Phase:"finish"}` directly.
- **N2 — AC3.3's per-phase warm arms exercise the RUN-promoted template, not the publication-built one.**
  `host/coordinator/capsule_template_test.go:69-83`: the cold run lazily promotes its cache, so `buildTemplate`'s `PromoteTemplate`
  hits the racer rule and discards the check-built tmp. Measured with own mutant OWN-8 (check-built template hollowed: modules
  removed before promotion): all 10 phase arms `--- PASS`, only `cross-source` fails (`mut/OWN-8-…log`). Check-built equivalence
  is therefore gated for 1 of 10 inputs (exec plan) plus the e2e rigs. Fix: `rig.dropTemplate(e)` between the cold run and
  `rig.buildTemplate(e)`.
- **N3 — finish budget at the call sites is unpinned (own mutant OWN-2 survives).** `effectful.go:231,327`: replacing
  `c.finishBudget` by `c.planBudget` (finish gets 4 s, breaking the HandlerHeadroom arithmetic AC2.4 relies on) passes the whole
  coordinator and projection packages. Pre-row-153 the literal constant was used; the field indirection opened the gap. Fix: a
  mirror of `TestSlowPlanLeavesHandlerItsCap` asserting the finish context's deadline ≤ `FinishPhaseBudget`.
- **N4 — `CopyCacheTree`'s symlink refusal is untested (OWN-4 survives archive + capsule packages).** `capsulecache.go:153-170`;
  the doc comment claims "a template can never smuggle a path out". Fix: a test with a symlink in the source tree → error.
- **N5 — OWN-6 survives**: promoting a FAILED cold run (`promoteCold` moved before the `runErr` check, `capsule.go:287-292`) passes
  the capsule package. Low impact (PromoteCopy still needs a manifest; the interpreter validates entries — I observed
  `CACHE_INVALID … reason=ARTIFACT_INVALID; recom[piling]` in OWN-8's log).
- **N6 — OWN-7 survives**: `CapsuleTemplateReady` checking dir existence instead of `compile/manifest.json` passes the capsule
  package; a corrupt template would then run compiling but unlabelled. Low impact (promotion is a whole-dir rename).
- **N7 — other new wall-clock bounds.** `templates_test.go:174` (`Run announced nothing within 7s`) is gratuitous: the gate holds
  the rebuild indefinitely, so MUT-BLOCKING-PREWARM-b blocks forever and any bound kills it — raise to ~60 s with no loss of
  power. (It passed on CI 2/2, but daemon New under the -race leg is the 10 s-startup-probe regime of V16.) Lower risk:
  `phase_timeout_test.go:160` (effective slack = 600 ms of pre-plan work), `mcp_plan_timeout_test.go:53` (4 s slack).
- **N8 — template GC runs once at start against the head read in `New`**; a set published (new interpreter) between `New` and
  the maintenance goroutine could have its fresh digest pruned (`templates.go:98-104,170-178`). Window is small; consequence is a
  labelled cold run + lazy re-promotion. Note only.
- **N9 — `serveWith` (`setools_e2e_test.go:735`) builds its runner with `capsule.Config{}` (Log nil → io.Discard)**, so the AC3.6
  tripwire does not cover runs through that coordinator. Pass `errLog` into `capsule.Config{Log: …}`.
- **N10 — cold label is always `(template missing)`**, also for the unreadable-template path (`cachetemplate.go:42-50`); the design
  said `missing|rebuilding`. Cosmetic.
- **N11 — CI race-leg headroom**: e327415's race leg had verifygate 379 s, broker 312 s (8 m = 480 s per package). Runner
  variance vs b6f7dfb (216 s / 179 s); row 154's concern, not this row's code.

## Mutation re-runs (mine; each a one-line edit, named test `-count=1 -v`, file restored by `git checkout`)

| Mutant | Test | Result |
|---|---|---|
| MUT-ANYPHASE (drop `&& phase.Phase == "plan"`) | TestA2ADispatch + finish wire test | KILLED (table row FinishPhaseTimeout_bare; wire test PASSES) |
| MUT-ORDER-FINISH (case above unrecorded) | same | KILLED (table row; wire test PASSES — N1) |
| MUT-ORDER (case after generic deadline) | TestA2APlanTimeoutIsTypedAndNamesThePhase | KILLED (`"invocation exceeded its deadline"`) |
| MUT-RELABEL | TestClassifyPhaseErrMatrix | KILLED (:41) |
| MUT-NOWRAP (Unwrap → Canceled) | same | KILLED (:33) |
| MUT-SOFT-PUBLISH | TestPublishWithoutATemplateIsRefused | KILLED (:255) |
| MUT-SHARE | TestCapsuleTemplateIsCopiedNotShared | KILLED (:172) |
| MUT-STALE-KEY | TestCapsuleWarmEqualsColdForSeTools | KILLED (cross-source :95) |
| MUT-SILENT-COLD | TestCapsuleColdRunIsLabelled | KILLED (:190) |
| MUT-CROSS-DEVICE-RENAME | TestCapsulePromotionSurvivesCrossDevice | KILLED (:115) |
| MUT-NOTEMPLATE | TestCapsuleWarmIsFaster | KILLED at the MOD010 tooth (:229), before the ratio |
| OWN-1 PromoteCopy renames the run dir straight in | TestCapsulePromotionSurvivesCrossDevice | KILLED |
| OWN-2 finish uses planBudget | whole coordinator + projection pkgs | **SURVIVED** (N3) |
| OWN-3 classifyPhaseErr parent/phase args swapped | TestA2APlanTimeoutRetrySameTaskID | KILLED (:85) |
| OWN-4 CopyCacheTree follows symlinks | whole archive + capsule pkgs | **SURVIVED** (N4) |
| OWN-6 promote failed cold runs | whole capsule pkg | **SURVIVED** (N5) |
| OWN-7 Ready = dir exists | whole capsule pkg | **SURVIVED** (N6) |
| OWN-8 check-built template hollow | TestCapsuleWarmEqualsColdForSeTools | KILLED only by cross-source; 10 phase arms PASS (N2) |
| DOUBLE (ANYPHASE+ORDER-FINISH) | finish wire test / TestA2ADispatch | SURVIVED / KILLED (N1) |

## Adjudications

- **Candidate finding (controller): CONFIRMED and BLOCKING** — see B1. Second CI instance at e327415 (0.53) added to the first
  (0.63). Correct fix: drop the ratio, keep MOD010 (+ optional ColdRuns count); MUT-NOTEMPLATE stays killed (measured).
- **Broker baseline red (executor F1): CONFIRMED pre-existing.** `host/broker` is byte-identical between `990a6da` and `e327415`
  (`git diff --quiet 990a6da e327415 -- host/broker` → identical; the requested `checkout 990a6da -- host/broker` was therefore a
  no-op, done and restored). The test passed in my full plain and race runs; it reproduced 1/30 under `taskpolicy -b` + 8 bg
  burners (`broker-bgload.log`: `handlers_capture_test.go:122: partial stdout = ""`), 0/40 under 24 default-QoS burners. Its
  executor C0 log shows the same signature at base. A 300 ms exec budget on an untouched package — not row 153.
- **Linux AILANG_CACHE_DIR / MOD010**: confirmed working on ubuntu/amd64 from both CI logs (MOD010 asserts at :222/:228 silent,
  failure at :235; `TestPinnedInterpreterHonoursCacheDir` and `TestCapsuleWarmEqualsColdForSeTools` in ok packages, both legs).
- **V28 retry semantics**: plan runs before any durable write (`effectful.go:285` precedes `appendAndCommit`); AC1.4 test asserts
  0 dispatch/0 spend/untouched after the timeout and exactly 1/1 + one entry after the same-id resend; MUT-EARLY-INTENT per run-1
  ledger. The "resend the same task id" text is reachable only for `Phase=="plan"`.
- **Template integrity**: key = (interp digest, sha256(source)) (`capsulecache.go:73-75`); copy-in, never shared (MUT-SHARE killed);
  EXDEV-safe promotion (copy into `capsule-cache/.tmp-*`, same-dir rename; MUT-CROSS-DEVICE-RENAME/OWN-1 killed); racer = success
  (`isRacer`); GC keeps daemon pin + every head interpreter and skips when a source is unreadable (fail-safe).
- **Startup never blocks serving**: `reportTemplates` is stat-only in `New`; rebuild starts in `Run` after `Listen` + announce
  (`daemon.go:1092-1094`); MUT-BLOCKING-PREWARM(-b) killed per run-2 ledger.
- **No silent cold run**: labelled + counted (`cachetemplate.go:49-50`), tripwire in exec-srt and `newWSDaemon` rigs (gap N9).
- **Executor deviations run-1 1–11, run-2 1–6**: all accepted; each measured or behaviour-preserving. Run-1 #4 (MUT-STALE-KEY reds
  only cross-source) is the visible edge of N2. Run-2 #1 (status for every descriptor) is the safer superset.
- **AC4.2 outlook**: daemon (row 161's tests) passed in all 4 parallel legs observed post-M3 and in the srt step; the only red is
  B1. With B1 fixed, 5/5 is plausible but unmeasured (2 attempts, both with the B1 red).

## Worktree restoration
`git status --porcelain` → empty; `git diff --quiet` → clean; HEAD `e327415e63e9755a9859406949579a5fe14d60ff`. The temporary
probe test was deleted; the race run used a `git archive` copy under this directory (removed after use).

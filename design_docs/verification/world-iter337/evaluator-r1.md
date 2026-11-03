# Evaluator r1 — row 239 w-marketing-still-probe-commit (iter 337), head 9e88e94f

VERDICT: PASS — 91/100

| Criterion | Pts | Score |
|---|---|---|
| Correctness vs acceptance | 35 | 33 |
| Test strength / mutation | 20 | 17 |
| Scope & conflict surface | 15 | 15 |
| Evidence integrity | 15 | 13 |
| Standards & maintainability | 15 | 13 |

## Commands re-run (pins set unless stated; all via run_bounded.sh.txt)
- `go test -count=1 -v -run '^TestMarketingStillRoutes$' ./cmd/ailang-worldd`: rc 0, PASS 8.2 s, 4 STILL lines, all `repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1`.
- CI-shaped control (only AILANG_BIN): rc 0, `--- SKIP: TestMarketingStillRoutes`.
- Six-shot drill (WORLD_CHROME + WORLD_MARKETING_RECTS=world-iter312/r1-rects.json + both pins): rc 0, `--- PASS: TestMarketingChromeSixShots (22.3s)`; pass-rule grep (`STILL .* repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1 `) = 4; fail grep (`repeat_equal=false|script_tags=[1-9]`) = 0.
- `./scripts/verify_ail.sh`: rc 0 ("verify gate PASSED: 16 identities, 40 named, 503 se-tools").
- `go test ./...` full: first run rc 1 with three unrelated failures (host/daemon TestSeToolsRunBenchmarkTasksEndToEnd, host/verifygate TestMissionConfigDefaultFlowIntegration "fixture timed out", scripts/marketingcheck supervisor cleanup "exceeded bound"). I was running mutation tests concurrently, so these are load flakes. Isolated re-run: verifygate, marketingcheck and cmd/ailang-worldd ok. host/daemon failed again in a 4-package run; that test passes in isolation at HEAD and at base `2f293570`, and the whole `./host/daemon` package alone is ok (rc 0, 343 s). No host/ file is in the diff. The executor's own full gate (executor-fullgate.txt) was rc 0. Conclusion: not caused by this change, but a clean single `go test ./...` was not reproduced under my load; there is a pre-existing load-sensitive flake (NON-BLOCKING).
- gofmt -l empty; go vet ./cmd/ailang-worldd ok. `git diff 2f293570 --name-only -- . ':!design_docs'` = exactly the two test files.

## Item 1: pin-gated vs Chrome-gated
The helper needs no Chrome. It compares the default-daemon HTML with the still-daemon HTML, with `checkMarketingStill` run on a real demo-seeded store. I measured it catching every production mutant that varies the still page (see below). It has four independent assertions: A0 (default tag count = 1), A1 (repeat equality across 2 s), A2 (zero `<script`), A3 (still equals default minus the tag). The iter305 probe only logged (per design V20), so this is stronger than the row asked. It also runs in the drill on that drill's own store before capture. The CI-shaped control skips cleanly, so CI stays green without Chrome. The drill gate catches a removed call via the STILL count (=4); MI1 was derived, not run, by the authors. It was not run by me either (the drill gate is count-based). Acceptable.

## Mutation (copy of worktree at /tmp/j337/m; each restored, verified with cmp against the worktree)
Production mutants (against the committed test), all caught unless noted:
- P1 render.go: tag always emitted: rc 1, A2 + A3 on all 4 routes.
- P3 still varies only on the payload=1 route (counter): rc 1, A1 + A3 on exactly that route.
- P4 deterministic extra bytes in still only: rc 1, A3 only (A1 correctly blind to it).
- P5 counter/2 on still: rc 1, A1 + A3 on all routes.
- P6 Unix-second suffix on still: rc 1, A1 + A3.
- P7 Minute() suffix on still: rc 1, A3 (a still-only addition is always caught by A3). A minute-granularity drift present in both modes would flake rather than fail deterministically (inherent).
- P8 live.js tag drift (`defer`→`async`): rc 1, A0 control "default HTML has 0 live.js tags" + A3.
- P9 still emits an extra `<script>1</script>`: rc 1, A2 + A3.
- P10 still page changes only after 1.5 s: rc 1, A1 only on all 4 routes (A3 passes; this isolates the repeat check). P10 plus the test's `time.Sleep(2s)` removed: rc 0, SURVIVED (expected; shows the sleep is load-bearing, and nothing guards the sleep itself).
Test-helper mutants:
- T6 still daemon built with WorkbenchStill=false: rc 1, A2 + A3.
- T1 routes[:3]: rc 1, the len==4 fatal.
- T8 ProofRequiredIdentities nil: rc 1, daemon.New fails (the config is exercised, not decorative).
- T3 `repeat = Equal(a,a)`: rc 0 SURVIVED. T9 `second=first`: SURVIVED. T10 A0 check disabled: SURVIVED. T11 final loop over `routes[:1]`: SURVIVED. T4 sleep removed: SURVIVED with the stateless P-mutants.
  These self-weakening edits are survivors only because there is no meta-test. T11 is partly covered by the drill's STILL count=4, but not by the stand-alone test; the len==4 guard does not cover the assertion loop.
- T5 (compare against the wrong body) did not compile in my patch (unused vars). The not-run caveat applies; P4/T6 exercise the same A3 line.

## Item 4: orchestration doc
`git diff 2f293570 9e88e94f -- w-marketing-capture-orchestration.md` is exactly 1 line replaced (+1/−1, numstat 1 1). The replaced line (§4.4 item 1) is byte-identical to the fenced paste block in the designs: `grep -cF "<line 221>"` of the design at `f8bc4a32` = 1, at `ca61e793` = 1, and at HEAD = 1. Both quorum rounds (r1 and r2) therefore reviewed the exact text. The decision not to run the orchestration doc's own quorum is sound on content: the text is unchanged since the design quorum, and the design explicitly leaves the choice to the controller (row 238 V27). The text is a replacement of a log/command pass rule by one that matches the committed test, and its figures check out (the drill count = 4 reproduced above). Caveat (NON-BLOCKING): the reviewed artifact was the design doc, not the orchestration doc itself, so reviewers did not see the surrounding §4.4 context; the single-line scope makes this immaterial.

## Item 5: process claims
- r1: gpt6-1-sol reject, gemini-4-argon pass, oc-glm-5-3 reject — matches design header/§6 ("gpt6-1-sol and oc-glm-5-3 REJECT, gemini PASS"). The objections in §6 map 1:1 to the JSON strongest_objection fields.
- r2: gpt6-1-sol reject, gemini pass, oc-glm-5-3 pass; synthesis "blocked" on the sol objection — matches §7. The §7 claim "does not dispute direction (accepts a separate implementation if justified)" is a fair reading of the objection ("not why those assertions cannot be reused or extracted"). The controller's in-session "pass" and carve-out is a process shortcut: the synthesis was BLOCKED and r3 went to planning without a third quorum. It is disclosed honestly as a "controller narrow-refinement carve-out", and §2.6 supplies the reviewer's own fix. NON-BLOCKING, but it is a judgment call the next reviewer should see.
- "no direction dispute" for r1: the r1 objections are evidential, none disputes the direction. Accurate.

## Findings
BLOCKING: none.
NON-BLOCKING:
1. The 2 s sleep and the final assertion loop are unguarded: removing the sleep, restricting the loop, or neutering A1/A0 passes in the stand-alone test (T3/T4/T9/T10/T11; `/tmp/j337/mut.sh` recipes above). Only the drill's STILL count=4 backstops a loop-shrink. A tiny assertion that the loop covers len(routes)==4 (e.g. count assertions executed) would close it.
2. Duplicate-by-design with `workbench_still_test.go` (§2.6), with three invariants to keep synchronized, documented in the file header. Acceptable and honestly justified.
3. Full `go test ./...` showed load-sensitive flakes in unrelated packages (host/daemon SeTools e2e, verifygate fixture timeout, marketingcheck supervisor bound) when the machine was loaded by concurrent runs; unchanged by this diff, and all pass in isolation or in executor's run. Worth a separate row, not this one.
4. Design r2 was REJECTed by one reviewer and proceeded without a third quorum round (disclosed).
5. MI1 (drill call deleted) remains derived rather than run (disclosed by authors).

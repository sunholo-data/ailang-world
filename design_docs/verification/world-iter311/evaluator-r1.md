# M7c4 evaluator report, round 1 (worktree .eval-world-iter311 @ 1ca31e5e)

VERDICT: PASS, score 92/100

Rubric: Tests pass 20/20; Lint 10/10; Acceptance criteria 29/30; Code quality 13/15; Documentation 12/15 (no doc/plan status edits in the diff; the controller owns those; plan M7c4 still reads "pending"); Design fidelity 8/10 (A1 sound; the Go rect mirror is weaker than LoadRectangles). Conditional categories not triggered.

Diff scope: exactly the 3 stated files (capture_marketing.sh, marketingcheck/capture_test.go, cmd/ailang-worldd/marketing_chrome_test.go). demo_test.go untouched. Worktree restored: `git status --porcelain` empty, `git diff --quiet` CLEAN after all mutations.

## Gates (AILANG_BIN / WORLD_TOOL_AILANG_BIN pinned)
- gofmt -l cmd scripts host: rc0, no output
- go vet ./cmd/ailang-worldd ./scripts/marketingcheck: rc0
- bash -n scripts/capture_marketing.sh: rc0
- go test -run '^$' ./... (compile fence): rc0
- TestProductionGoSurface|TestProductionContextRoots ./host/store: rc0
- TestSubprocessSinkGate ./host/verifygate: rc0
- TestMarketingCaptureFakeChromeSixShots -v: rc0, 13 subtests PASS incl. inventory, inventory-pass-b, inventory-single-pass (15 s)
- go test ./scripts/marketingcheck (whole package): rc0, 33.5 s. TestMarketingSupervisorFakeDaemonCleanup did not flake.
- Drill, WORLD_CHROME unset: rc0, SKIP.

## Completion semantics (measured)
- (b) unset -> SKIP (rc0); WORLD_CHROME=/nonexistent -> FAIL rc1 at :58 "must be an executable file".
- (c) Chrome set, WORLD_MARKETING_RECTS unset -> FAIL :75 "never guess rectangles"; nonexistent path -> FAIL :91; `{}` -> FAIL :94 (invalid schema/shot set); packet missing a roster name -> FAIL :111. Missing WORLD_TOOL_AILANG_BIN -> FAIL :71; AILANG_BIN passed as the tool pin -> FAIL :118 (version mismatch). All precede seeding/building. Order follows plan (Chrome stat, chrome-version, pins, rects), with the pin --version probes after rects; they still precede seeding.
- (d) inventory subtest passes; M7c4-inventory-wiring killed (see table).
- (a),(e): outside this measurement.

Bonus measurement: with real Chrome and a synthetic (not measured) rectangles file (0,0,100,100 per content rect), the full drill ran end to end and PASSED in 25 s: rig seed, builds, capture_marketing.sh --verify-twice, capture-manifest, diagnostics "inventory: equal". So the plumbing is good for R1. It is not evidence, because the rectangles were made up.

## Red-first
Base capture_marketing.sh (90ef2671) with the new test: FAIL on six-shots (missing inventory-before.json, capture_test.go:480), inventory, inventory-pass-b, inventory-single-pass ("capture must compare inventories", :366). Red as required. The drill test is a skip-by-default drill with no base counterpart; its precondition arms were exercised directly above.

## Mutants (shell production diff), all via TestMarketingCaptureFakeChromeSixShots
| Mutant | Result | Assertion |
|---|---|---|
| M7c4-inventory-wiring: delete `inventory compare` step | KILLED | inventory, inventory-pass-b, inventory-single-pass: "capture must compare inventories" (:366) |
| delete `inventory before` | KILLED (8 subtests) | :366, :389, :458 |
| delete `inventory after` | KILLED (8 subtests) | same |
| inventory before moved to after pass A | KILLED | inventory, inventory-single-pass :366 |
| inventory after+compare right after pass A even with --verify-twice (the pre-A1 design) | KILLED | inventory-pass-b :366 (this is the reason that subtest exists) |
| compare given before.json twice | KILLED | inventory x3 :366 |
| `inventory after` writes to inventory-before.json | KILLED (8) | :366/:389/:458 |
| diagnostics line drops `inventory: equal` | KILLED | six-shots :484 |
| `[[ verify_twice ]] \|\| exit 0` moved above inventory after | KILLED | inventory-single-pass :366 |
| compare only inside verify-twice | KILLED | inventory-single-pass :366 |
| inventory after pointed at wrong URL | KILLED (8) | :366/:389/:458 |
| diagnostics printf (with "inventory: equal") moved before compare | PARTIAL: killed only via single-pass (stray diagnostics file); SURVIVES in verify-twice mode. The mismatch subtests return before checking diagnostics absence. Non-blocking finding F1. |
| `inventory compare ... \|\| true` | SURVIVED, equivalent: `fail` calls `exit 1`, which is unaffected by `\|\|` |
| compare --before/--after swapped | SURVIVED, equivalent: comparison is symmetric |
| inventory before moved after `routes` | SURVIVED, equivalent: no daemon-state change between them |

Review check M7c4-rects-required: the guard is test code at marketing_chrome_test.go:73-76 (`t.Fatal`, not Skip). Confirmed by running (c) above. Changing it to Skip would turn the (c) arms green-skip, a manual check as the plan says.

## A1 (single-pass compares after pass A)
Sound. The inventories bracket every pass actually run, in both modes. Version-B is checked before the after-inventory, so a version change refuses first. The inventory-pass-b subtest proves the post-B placement matters (a head change keyed to pass-b is only caught by an inventory taken after pass B), and the inventory-single-pass arm covers the no-flag mode. The fake head change is keyed to existence of `pass-a-*.args`/`pass-b-*.args` files, so it genuinely discriminates ordering (the before-moved and after-moved mutants die on it). Plan wording ("after pass B") is respected for the only mode that has a pass B.

## Test code review
- Fake daemon: serves /v1/log paging and /v1/objects/{ref}; only /v1/head varies. It does not discriminate entry/object differences at the shell level. That is covered by M7c3's own tests. Fine.
- Rig helper: seed daemon, health, publish tools, mint, second daemon for `demo seed`, then `srv.Close()`/`d.Close()` before return, so the capture shell's serve opens an unlocked db. Verified by the real run succeeding. demo_test.go not edited.
- Sinks: all subprocesses use os.Create files under t.TempDir; contexts bounded (WithTimeout on t.Context(), boundedTestContext). TestSubprocessSinkGate green.

## Findings
1. non-blocking, marketing_chrome_test.go:80-112: the rects mirror is weaker than LoadRectangles (no unknown-shot, bounds, min size checks). A bad-but-well-formed file (phone content w=99999) passes the preconditions, runs the whole seed/build/capture (~18 s here, minutes with real timings), then fails at the final `check`. Violates "unloadable by LoadRectangles semantics -> fail early" only in degree.
2. non-blocking, marketing_chrome_test.go:140-146: capture-manifest check is only "parses as a JSON object" (`{}` would pass); VerifyCaptureManifest is not importable from package main. Verification is delegated to the shell's `check` rc 0. Acceptable, but worth noting.
3. non-blocking, capture_test.go:~360-366 (F1): no assertion that a mismatch run leaves no `diagnostics.txt`/"inventory: equal"; the diagnostics-before-compare mutant survives in verify-twice mode.
4. non-blocking, marketing_chrome_test.go:118: pin release strings (v0.41.0 for the kernel) are hardcoded in the test; brittle on a kernel pin bump. The tool uses daemon.ToolBinaryRelease.
5. non-blocking: go test -timeout 660s (plan's rehearsal command) equals the inner capture bound (660 s) while rig + builds add on top; a slow run dies on the go test timeout rather than the inner bound. Plan-specified; R1 should use a larger -timeout if the capture nears its limit.

No blocking findings.

## Premise / risk for R1
- The guard sees head, entries, and object hashes before/after both passes. It cannot see in-process daemon-state changes not exposed via /v1/head, /v1/log, /v1/objects; and nothing about whether rectangles are correct (a bogus 0,0,100,100 set passed on real Chrome).
- R1 must supply genuinely measured, in-bounds rectangles for all 6 shots; thresholds in png.go are still provisional (M8a applies the V-row).
- Real Chrome here was 25 s end to end with fake rects; no environmental blocker seen (Chrome version ok, both pins ok, seed ok).
- The unattended drill is skip-by-default, so PR CI never exercises it; only the fake-Chrome shell test guards the wiring in CI.

EVALUATION_SCORE: 92/100
EVALUATION_ROUND: 1
EVALUATION_VERDICT: PASS score=92

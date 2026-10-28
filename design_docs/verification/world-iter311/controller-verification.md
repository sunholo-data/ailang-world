# Iteration 311 — row 150 M7c4 — controller verification

All runs below are the controller's, outside the codex sandbox, with
`AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0) and
`WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang`. Load was 2.5–6.

## Baseline (pristine worktree at base `90ef2671`, before the directive was sent)

- gofmt clean; `go vet ./cmd/ailang-worldd ./scripts/marketingcheck` rc 0; `bash -n scripts/capture_marketing.sh` rc 0.
- Compile fence `go test -run '^$' ./...` rc 0 (34 packages).
- `TestMarketingCaptureFakeChromeSixShots` ok; `TestProductionGoSurface|TestProductionContextRoots` ok; `TestSubprocessSinkGate` ok; `TestMarketingChromeSixShots` "no tests to run" (file absent).
- `go test ./scripts/marketingcheck` **rc 1 at base**: `TestMarketingSupervisorFakeDaemonCleanup` — "shell trap cleanup exceeded bound: 29.7s" (load ~6, inside the package run). Re-run alone with the identical `-run` 3 times at load ~3.6: **3/3 ok**. Pre-existing and load-sensitive (the M7a fixture class iteration 308 carried); named in the directive as not the executor's to fix.

## Executor diff (`31a4929b`), read in full

- `scripts/capture_marketing.sh` (+10/−5): `inventory before` after readiness, before `routes`/pass A; pass B and the version-B comparison moved inside `if [[ $verify_twice == yes ]]`; `inventory after` + `inventory compare` after the last pass in both modes; `inventory: equal` added to diagnostics.txt.
- `scripts/marketingcheck/capture_test.go` (+38/−6): fake daemon serves `/v1/log` (paged) and `/v1/objects/{ref}`; `inventory` subtest.
- `cmd/ailang-worldd/marketing_chrome_test.go` (new, 246 lines): `TestMarketingChromeSixShots` + `newMarketingChromeRig`. Preconditions all precede seeding. `demo_test.go` untouched.
- No write outside the worktree in the diff (rule 5 of codex-lane-false-greens): every path is `t.TempDir()`.

## Gates on the executor tree

All green outside the sandbox: gofmt, vet, `bash -n`, compile fence, SKIP with `WORLD_CHROME` unset, FAIL at `marketing_chrome_test.go:58` with `WORLD_CHROME=/nonexistent`, fake-Chrome 11/11 subtests, surface/context roots, sink gate.

## Red-first

Base `capture_marketing.sh` with the new test: `six-shots` red (`inventory-before.json: no such file`), `inventory` red at "capture must compare inventories: rc=0".

## Mutation drill (`controller-mutants.py.txt`, anchors asserted unique, `bash -n` checked, sha restored)

Before the controller fix: 6 of 9 killed. Survivors: before-inventory moved after pass A; after-inventory moved before pass B; single-pass mode skipping the compare. Cause: the executor's fake daemon changed `/v1/head` after the first `/v1/log?from=0` read, so it could not tell when the inventories ran relative to the passes.

**Controller test-only fix `1ca31e5e`:** the fake daemon's head changes once a fake-Chrome shot of the named pass has run (`CAPTURE_HEAD_CHANGE=pass-a|pass-b`, globbing the `<pass>-<shot>.args` files fake Chrome already writes into the same directory). New subtests `inventory-pass-b` and `inventory-single-pass`. After it, **9 of 9 killed**; the plan's M7c4-inventory-wiring dies at "capture must compare inventories".

## Completion semantics

- (b) SKIP with `WORLD_CHROME` unset; FAIL with `/nonexistent`.
- (c) With real Chrome set: rects unset → FAIL `:75`; rects `{schema, shots:{}}` → FAIL `:94`; rects path missing → FAIL `:91`; `AILANG_BIN` empty → FAIL `:71`.
- (d) `inventory` passes; plan mutant killed.
- (e) judge sonnet r1 PASS 92 (`evaluator-r1.md`).

## Real-Chrome smoke (not C8 evidence; R1 is the measured run)

Chrome 156.0.8078.48.
- Placeholder 40x40 rectangles (`real-chrome-smoke-placeholder40.txt`): real seed, real `serve`, 6/6 shots in each pass, inventory before/after compared equal, then `check` refused `rect-too-small workbench-light.png: c1`. rc 1 in 20 s, as expected.
- Made-up 100x100 rectangles at (0,0) for every content rect (`real-chrome-smoke-madeup100.txt`): **PASS end to end**, rc 0 in 15 s. Reproduces the judge's finding. `check` does not reject guessed rectangles, so R1's C8 claim rests on the rectangles being genuinely measured.

## Final tree `1ca31e5e` (`final-gates.txt`)

gofmt clean, vet rc 0, `bash -n` rc 0, fence rc 0, SKIP, fake-Chrome ok, `-race` fake-Chrome ok, surface/context ok, sink gate ok, `go test ./scripts/marketingcheck` ok (the flake did not fire), `go test ./cmd/ailang-worldd` ok (144 s).

## Judge findings, re-measured by the controller

- Diagnostics line moved before the compare: killed only by `single-pass` (incidentally, because the line now runs in single-pass mode). Reproduced; non-blocking, carried.
- Made-up rectangles pass on real Chrome: reproduced (above).

# Evaluator report r1 - iteration 316, row 150 M8b (commit 3f185a9a)

Verdict: PASS, 93/100. No blocking findings.

Scores: doc correctness 31/35 (missing-file hint, rects path unguarded); scope and protected paths 20/20; mutation drill 14/15; gates 14/15 (full/race are controller-owned, unassessed); report honesty 14/15 (minor count slip). Total 93.

## Findings
- NON-BLOCKING docs/QUICKSTART.md:1656-1663: WORLD_MARKETING_RECTS moved from R1 file to $DEMO_DIR/rects.json. TRUE per orchestration P3 (final bullet) and §4.4 "(a) accepted, never the R1 file". Prose says "the controller supplies" the file but never names how it is produced (the V60 block-box probe, §4.4 P0 step 3 / P3). An operator following §13 verbatim hits a missing file until the controller supplies it. That is a behaviour change versus base (base pointed at an existing R1 file, which §4.4 says is wrong for capture), so it is an honest reflection of §4.4, not a regression in correctness. Suggest one clause naming the V60 probe.
- NON-BLOCKING: no test guards the rects path. Mutant reverting line 1663 to design_docs/verification/world-iter312/r1-rects.json: TestMarketingRunbookMatchesHelp and TestMarketingDocumentedServeSelectsStill still PASS (survives).
- NON-BLOCKING docs/QUICKSTART.md:1698-1706: accurate but terse; P0 runs on a disposable test store, P3 on release store - sentence does not distinguish. Not false.
- NON-BLOCKING report: says "All 52 invocations", ledger had 51 initial entries (self-disclosed in entry 52). Entry 25/36/40/41 socket tests were "UNINFORMATIVE UNDER SANDBOX"; in my environment all four PASS (below). Report honestly did not claim PASS and lists unresolved P0, P3, A1, C8, promotion, C10. No overclaim found.

## Clause checks (TRUE as written)
- "P0 (a) now requires structural match against R1 (§4.4 (a1)-(a4))": TRUE.
- "phone C8'(iii') requires the disclosed overflow defect present, without magnitude": TRUE (§4.4 phone, §7 C8').
- Checker sha pins: v64-p0check-main.go.txt 0a9a6795...927c and v66-phonecheck-main.go.txt 31576922...0b4c recomputed with shasum, match §4.4. Files exist in design_docs/verification/world-iter315/.
- "r3 column / r3 field, ACCEPT/REFUSE, exit 0 either way": TRUE (§4.4 source pins paragraph, V72).
- "P0 within 24 h; P3 before pass A": TRUE.
- README 127-128, 187-190: row 156 CLOSED iteration 312, D-WORLD-86 = A drops M5c/M5d, relies on v0.54.15 pin + M5e (charter line 6197, D-WORLD-86 row 1005). D-WORLD-81 ordering (150, 151, then v1.0.0) TRUE. Remaining README/website/captions/schema: no other stale 156/open-criteria sentence found; website page "not yet in the site until promoted" still true.

## Scope / C11
- git diff --stat 8bfaa1ab..3f185a9a: README.md, docs/QUICKSTART.md only (both allowlisted). Diff of world, go.mod, go.sum, tools/launchd, design_docs/world-mission.md: empty.
- /tmp/mc316 observe on my worktree, compare vs obs-before.json rc 0; cmp byte-identical.

## Mutations (worktree .eval-world-iter316, restored byte-identically via cmp)
1. M8b-C11-redrill: `if false && b.Lines != a.Lines` at scripts/marketingcheck/main.go:427. Reds exactly one subtest: TestMarketingProtectedObservationComparator/digest-equal-line-change (observation_test.go:160 "line-count mismatch must refuse: want Code=observation-lines, got <nil>"); all other 10 subtests PASS. Restored, re-run PASS. KILLED at the named assertion.
2. QUICKSTART capture line `--rectangles` -> `--rectangls`: TestMarketingRunbookMatchesHelp/capture FAIL (docs_test.go:404,409). KILLED.
3. website marketing-capture.md image packet.png -> packetx.png: TestMarketingDocsLinks FAIL (docs_test.go:282; M8a-filename sub-case also fails at :322 since mutation target gone). KILLED.
4. Revert WORLD_MARKETING_RECTS to R1 path: SURVIVED (see non-blocking finding).

## Gate table (my worktree, AILANG_BIN and WORLD_TOOL_AILANG_BIN set)
| command | rc | named PASS |
|---|---|---|
| TestMarketingProtectedObservationComparator | 0 | PASS (3.87s) |
| TestDemoSeedEndToEnd (cmd/ailang-worldd) | 0 | PASS (21.20s) |
| TestDemoSeedWhyNamedLinks | 0 | PASS (33.87s) |
| TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive | 0 | PASS |
| TestMarketingCaptureFakeChromeSixShots | 0 | PASS (57.84s) |
| TestMarketingShotGeometry | 0 | PASS |
| TestMarketingPNGStructure | 0 | PASS |
| TestMarketingDocsLinks | 0 | PASS |
| TestMarketingRunbookMatchesHelp | 0 | PASS |
| cmd/ailang-worldd ^TestQuickstart | 0 | PASS: AttendedStepsBoxMatchesSection9, Section9/10/11FlagsMatchTheCLI |
Not run by me (controller-owned): go test ./..., -race suites, verify_ail.sh (executor: rc 0 with nine pre-existing property skips, not a zero-skip claim).

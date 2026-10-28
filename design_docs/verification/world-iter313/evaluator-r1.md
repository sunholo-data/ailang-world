# Evaluator r1 - row 150 M8a (iteration 313), commit f6c8d8d4

## Verdict: PASS, 82 / 100 (no BLOCKING findings)

| Criterion | Score |
|---|---|
| Completeness vs plan M8a and orchestration section 5 | 27 / 30 |
| Correctness and truthfulness | 22 / 25 |
| Test quality and non-vacuity | 15 / 25 |
| Scope and hygiene | 10 / 10 |
| Honesty of executor report | 8 / 10 (executor's report not seen by me; scored on diff vs directive, no overclaim found) |

## Gates (all run by me in EVAL_WT; gofmt -l empty; go vet rc 0)
| Command | rc | PASS lines |
|---|---|---|
| TestMarketingDocsLinks (scripts/marketingcheck) | 0 | --- PASS: TestMarketingDocsLinks (0.16s) |
| TestMarketingRunbookMatchesHelp | 0 | --- PASS (1.89s) |
| TestQuickstartDecisionPacketSection (host/daemon) | 0 | --- PASS (0.03s) |
| cmd/ailang-worldd set (ServeWorkspaceFlags, ServeProofFlags, ServeWorkbenchStill, DemoSeedHelpPreparation, MarketingDocumentedServeSelectsStill, Quickstart*) | 0 | all 9 PASS |
| go test ./scripts/marketingcheck (whole package) | 0 | ok 53.0s |
| TestProductionGoSurface (host/store) | 0 | ok |
| go run ./scripts/marketingcheck docs --root . | 0 | - |
No timing-sensitive re-runs were needed.

## Scope
Exactly the ten allowed files changed (git diff --stat). No .png added. go.mod, go.sum, world/*.ail, tools/launchd, world-mission.md untouched. png.go change is a comment only (V62 says keep all eight values); png_test.go change is the single word docs -> bogus. Both match orchestration V62 and section 5.

## Mutation table (all restored; shasum byte-identity confirmed after each)
Harness: /tmp/ev313/mut.sh. D = TestMarketingDocsLinks, R = TestMarketingRunbookMatchesHelp.

| Mutant | Edit | Cmd | rc | Result |
|---|---|---|---|---|
| M8a-filename (plan) | page src packet.png -> packet-v2.png | D | 1 | KILLED real-repo: "docs-image packet-v2.png: documented image target must match promoted inventory" |
| M8a-flag (plan) | QUICKSTART 13 --workbench-still -> --workbench-static | R | 1 | KILLED M8a-flag and serve-preview: "documented serve argv must parse and select still mode: serve preview lacks --workbench-still" (killed by the string-presence pre-check, not by the parser leg; see F4) |
| M8a-flag, in-process | same edit | TestMarketingDocumentedServeSelectsStill | 1 | KILLED |
| capture-flag | drop --marketingcheck-bin from capture line | R | 1 | KILLED capture: "lacks --marketingcheck-bin" |
| placeholder | skip VerifyCaptureManifest | D | 1 | KILLED promoted-dir: "want refusal manifest-mismatch, got <nil>" |
| captions disclosure loop off | main.go | D | 1 | KILLED refuse-captions-disclosure |
| captions alt-marker off | main.go | D | 1 | KILLED refuse-alt-route |
| schema 64-hex check off | main.go | D | 1 | KILLED refuse-schema-hash |
| README embed check off | main.go | D | 1 | KILLED refuse-readme-embed |
| link stat off | main.go | D | 1 | KILLED refuse-readme-broken-link |
| README section links off | main.go | D | 1 | KILLED refuse-readme-broken-link |
| CLI capture-manifest-without-require off | main.go | D | 1 | KILLED cli |
| row dims check off | main.go | D | 1 | KILLED promoted-dir |
| attended: prefix `yes \|` on capture line | QUICKSTART | R | 1 | KILLED attended-only |
| attended: drop "typed by a person" | QUICKSTART | R | 1 | KILLED attended-only |
| schema const -> v2 (doc edit) | schema.json | D+R | 1 | KILLED real-repo |
| page alt vs captions mismatch | captions Alt shortened | D+R | 1 | KILLED |
| page alt-marker check off | main.go | D | 0 | SURVIVED |
| page alt==captions check off | main.go | D | 0 | SURVIVED |
| shots min/maxItems check off | main.go | D | 0 | SURVIVED |
| schema const check off | main.go | D | 0 | SURVIVED |
| schema top-level required (schema/frozen-head/shots) off, and reduced to {schema} | main.go | D | 0 | SURVIVED |
| schema item-required check off | main.go | D | 0 | SURVIVED |
| page markdown-image check (docsPageEmbed) off | main.go | D | 0 | SURVIVED (equivalent: the /img/v1 count check also catches it) |
| page import-useBaseUrl check off | main.go | D | 0 | SURVIVED |
| page image count check reduced | main.go | D | 0 | SURVIVED |
| duplicate-image (seen) check off | main.go | D | 0 | SURVIVED |
| page relative-link check off | main.go | D | 0 | SURVIVED |
| QUICKSTART heading check off | main.go | D | 0 | SURVIVED |
| README QUICKSTART-link check off | main.go | D | 0 | SURVIVED |
| README light/packet link-required check off | main.go | D | 0 | SURVIVED |
| README no-Screenshots-section refusal off | main.go | D | 0 | SURVIVED |
| captions one-Alt/non-empty-body check off | main.go | D | 0 | SURVIVED |
| captions duplicate-section check off | main.go | D | 0 | SURVIVED |
| captions extra-PNG-section check off | main.go | D | 0 | SURVIVED |
| CheckPromotedAssets: CheckDirTotal off | main.go | D | 0 | SURVIVED |
| CheckPromotedAssets: CheckStructure off | main.go | D | 0 | SURVIVED |
| CheckPromotedAssets: regular-file (symlink) check off | main.go | D | 0 | SURVIVED |
| CheckPromotedAssets: extra-shot refusal off | main.go | D | 0 | SURVIVED |
| docs-usage refusal in CheckPromotedAssets off | main.go | D | 0 | SURVIVED (CLI layer also gates) |
| page: replace static-capture note with "live views" | page md | D+R | 0 | SURVIVED |
| page: drop separate-pins sentence | page md | D+R | 0 | SURVIVED |
| page: "No approval decision appears" -> "approval granted" | page md | D+R | 0 | SURVIVED |
| README: "static captures" -> "live screenshots" | README | D+R | 0 | SURVIVED |
| QUICKSTART 13: delete P0-within-24h sentence | QUICKSTART | D+R | 0 | SURVIVED |
| QUICKSTART 13: preview bind 7645 -> 7999 (capture line bind differs) | QUICKSTART | D+R | 0 | SURVIVED |
Note: several of my first-pass "rc 1" results were compile errors (unused variable); I re-ran those with `false &&` guards and report the re-run results above.

## Non-vacuity
The named tests are not vacuous: both plan mutants and my four reach the named assertions, the real-repo check reads real files, and the capture-flag equality is against live `--help` output. The 64-zeros placeholder is refused. The in-suite M8a-flag proves the guard.

## Findings
1. NON-BLOCKING (test gap, plan C10 wording): plan C10 requires docs to "disclose static capture, demo-authored claim, evidence limits and separate pins", but only captions.md is gated for those strings. The website page and README text can be edited to "live views" / "live screenshots", or to drop the pins and the no-approval sentence, with every test green (mutants above). Fix: run the docsDisclosures list (or a page-specific subset) over the page text and require "static capture" in the README Screenshots section.
2. NON-BLOCKING (survivors in the schema check): const, minItems/maxItems, top-level required and item-required are asserted nowhere with a negative arm; only the hex scan and a const edit on the real file are covered. Fix: add schema refusal arms (const changed, maxItems 7, `required` missing frozen-head, item missing route) mirroring refuse-schema-hash.
3. NON-BLOCKING (survivors in asset checks): --assets mode has no negative arm for a corrupt PNG (CheckStructure), a symlink, an extra PNG, or an oversize directory, although the directive lists these. Fix: add promoted-dir arms for each (symlinked shot, sixth-plus extra `x.png`, bit-flipped CRC, directory over 8 MiB).
4. NON-BLOCKING (M8a-flag fidelity): the plan asks for the mutant to reach "the actual parser". It is killed by the `--workbench-still` string pre-check, and the shell-flag-set equality check would also catch it. The parser leg (`--bind not-a-hostport`, exit 1, "is not host:port") is exercised only on the happy path; no mutant is killed by it alone. The in-process cmd test does exercise the real flag parser. Fix: add a subtest where the preview carries `--workbench-still` plus a mis-valued flag (for example `--bind` consuming the next flag) so only the parse leg can red.
5. NON-BLOCKING (page/README survivors): page link resolution, useBaseUrl import, duplicate-image, README-section-missing, README light/packet-link-required and QUICKSTART-link checks are untested by any negative arm. Fix: one overlay subtest each.
6. NON-BLOCKING (public wording): README says the screenshots "are static captures ... published with the v1.0.0 release capture" in the present tense while both linked files 404 on GitHub until promotion. The website page has the right "not yet in the site" note, but when the site deploys it will show six broken images. Fix: README "will be published with ... (not yet promoted)"; ideally keep the page unlinked from the sidebar until promotion.
7. NON-BLOCKING (accuracy nit): phone alt says "first timeline entry" (singular) while V50 measured entries 0-1; the light alt says "first entries". Align the phone alt to "first timeline entries" or drop the count. The gate cannot catch this since alts are free text.
8. NON-BLOCKING (runbook completeness): 13 names the P0 probe but not P3 on the release store, and does not mention that orchestration section 5 records phone C8'(iii) as open. Truthful but incomplete as a runbook; point to orchestration section 4.4/C12 as a whole. Also the preview bind and capture-line bind are not tied (mutant qs-wrong-bind survives).
9. NON-BLOCKING (minor): alt-marker check is case-sensitive ("Session", "Token" pass); `[a](b "title")` links are not parsed; `len(want) != 13` hard-codes the help flag count.

Truthfulness check of the public sentences: captions, page, README and 13 were compared with the parent section 2/3 shot table, orchestration 4.1/4.4/8, V42/V50/V62/V63, the capture shell's real serve argv and `--help`; "approvals chain: pending" is a real rendered string (workbench_packet_test.go). No false claim found; tall framing sizes match Shots in png.go; provenance uses the narrowed claim verbatim.

## UNMEASURED
- Website build and rendering (optional; not run).
- Real-asset release mode on actual Chrome PNGs: they do not exist until the attended capture.
- The executor's own report/snapshots under ~/.ailang/state were not read; honesty score is inferred.

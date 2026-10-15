# Executor repair r1 — World iteration268 / row150 M2a

Terminal implementation handoff, ready for independent round2; no self-issued evaluator PASS.

Repair tree `/Users/voightkampff/dev/sunholo-data/.repair-world-iter268`, branch `repair/world-iter268-m2a-r1`, base candidate `67ee2efe5b2fdca3dcab4eaeda4849a8caea1c8c`. Only three files changed: cmd/ailang-worldd/demo_seed.go, demo_test.go, docs/QUICKSTART.md. Previous source/evaluator worktrees untouched. No git writes, authority mint/publish, attended confirmations, subagents or duplicate broad suite. Actual declared role model native GPT-6.1 Sol; token usage unknown. Native transport deviation and same-model-fresh-context evaluation policy remain as controller-recorded.

Read the independent evaluator-r1.md/.json at `/Users/voightkampff/dev/sunholo-data/.eval-world-iter268/design_docs/verification/world-iter268/` (FAIL65/100,3blocking). Dispositions:

- **R1-B1:** Removed forbidden direct host/registry import. The CLI reads only the existing public registry wire shape into a local anonymous struct and json.Unmarshal, preserving payload-address verification, decode refusal and exact archived-release nomination. No copied registry publisher or new authority. Both existing dependency allowlist and recording-writer tests pass unchanged. Changed epoch decode guard has a compiling named wrong-refusal mutation kill; pristine full preflight still passes.
- **R1-B2:** New section12 exports its own `$DEMO_BIN=$DEMO_DIR/ailang-worldd`, builds that artifact and uses it for mint, serve and demo seed. It explicitly transfers the same variable values to the serving terminal. The complete QUICKSTART prefix before section12, including section11, is byte-identical to the candidate. No old parser/test/expectation changed. New preparation assertion checks all five artifact export/build/command strings; a changed command variable fails the named wiring assertion. The old section11 CLI binding test passes unchanged.
- **R1-B3:** Removed os/exec and the seed's direct version subprocess entirely. Existing daemon startup already archives bytes and obtains bounded, credential-scrubbed verbatim --version output. After matching local pin hash and read-only authority, the seed reads the existing archive sidecar with archive.New(db).ReadManifest(pinRef), verifies Hash and Size against matched local bytes, and compares Version exactly to health plus the v0.41.0 prefix. This is read-only and creates no subprocess/site, store writer, directory or authority. Existing AC10 nine-site driven/scrub audit passes unchanged. No actual key or sentinel was needed in this repair role because the offending launcher no longer exists.

## Verification

All final commands ran after mutations were terminal and source was restored. All five allowed source/document files had SHA256 asserted unchanged before/after each command. Same expected v0.41.0 and separate v0.52.1 environment pins; only GOOGLE_CLOUD_PROJECT, OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_RESOURCE_ATTRIBUTES and OTLP_GOOGLE_CLOUD_PROJECT removed from child gate environments. Literal argv, wall times and named test evidence are in repair-gates.json and matching logs.

| Required named tests | rc | Wall seconds |
|---|---:|---:|
| TestQuickstartSection11FlagsMatchTheCLI/TestWorldBoundaryDependencyAllowlist/TestWorldBoundaryRecordingWriter | 0 | 1.937 |
| TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential | 0 | 5.131 |
| TestDemoSeedAuthenticatedCommits | 0 | 1.709 |
| TestDemoSeedPreflightRefusals | 0 | 32.841 |
| TestDemoSeedHelpPreparation | 0 | 0.646 |
| TestSessionMint_RequiresControllingTerminal | 0 | 0.659 |
| TestRESTCommitRefusesDecisionPacketObjects | 0 | 0.687 |
| TestDemoSeedCommitFailures | 0 | 3.308 |

Each required top-level name appears PASS with no SKIP. The combined first command is the exact controller3-regression command (90s Go bound), followed by exact AC10 (90s), then the original6focused commands with their prior120s/45s bounds and150s outer bound. Preflight now has37 refusal cases; unchanged real daemon acceptance remains5 authenticated objects/entries/worlds and incomplete no-manifest boundary. No complete seed/release claim.

The added missing-archive test was run BEFORE replacing the direct launcher: it compiled and failed on five actual log entries where zero were required (archive-test-first.log). Five archive controls now cover missing sidecar, forged hash, wrong size, changed version and matching forged health/sidecar with unsupported release. Existing forged-health-version and readonly-db controls also pass.

repair-mutations.json: **7/7 compiling assertion kills** with exact restoration: local epoch decode; documented artifact command wiring; archive missing/read refusal; archive Hash; archive Size; archive/health version equality; release prefix. Metadata omission mutants reach actual log-inventory change where specified. Missing-sidecar guard omission is killed by its exact expected refusal assertion, with the subsequent identity guard still fail-closed; no inventory kill is claimed for that arm. First release-prefix mutation accidentally removed the last strings use and produced build-red; its raw log and repair-mutations-r1.json are retained and it was NOT counted. Corrected compiling mutant preserved the import and produced actual log entries5want0 (mutant-archive-release-r2.log).

All **12 protected paths** retain exact SHA256+newline equality, including recursive launchd testdata; protected scoped git diff, diffcheck and formatting checks pass. repair-observations.json records this and the unchanged complete QUICKSTART prefix. Relative-path copies of all five allowlisted source files are in source/.

Known nonblocking r1 rubric deductions (function size and absent changelog outside this five-file scope) were not broadened into repairs. Existing low-level operational-error coverage residuals remain scoped as in the original report. Full repaired-head Go/AIL, CI and independent round2 judgment remain controller-owned; a baseline or focused PASS is not substituted for those gates.

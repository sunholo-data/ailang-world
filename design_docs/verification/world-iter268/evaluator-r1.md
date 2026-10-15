# Independent evaluation r1 — iteration 268 / row 150 / M2a

**FAIL — 65/100; 3 blocking findings.** Reviewed commit `67ee2efe5b2fdca3dcab4eaeda4849a8caea1c8c`, base `39066a46c1ac13c994bacb162b8d4b754d5ec083`. No merge approval. This report judges only M2a's five implementation files; M1a/M1b are prior work and M2b onward intentionally pending.

Judge: fresh native Agent GPT-6.1 Sol, a different agent from the generator. **judge-independence: same-model-fresh-context**. Preferred Sonnet attempt unavailable per supplied routing evidence; no Astra, inbox reads or agent delegation. Actual token usage: **unknown**. Measurements are darwin/arm64, Go1.26.6, released capsule/proof pin v0.41.0 and separate workspace-tool pin v0.52.1. The evaluator skill/rubric, authoritative mission-control and verification rules (especially3j/3n/3o), charter, CLAUDE, coding standards, DESIGN§1/14, r4 design and both plan artifacts governed this review. r1–r3 remain BLOCKED; r4 inherits the iteration266 narrow-refinement disposition, not an invented quorum PASS.

## Blocking findings

1. **R1-B1 — direct forbidden registry dependency.** `cmd/ailang-worldd/demo_seed.go:21` imports `host/registry`. The existing production attribution scanner explicitly forbids direct registry imports in this CLI. Independently ran `go test -run '^TestWorldBoundary(DependencyAllowlist|RecordingWriter)$' ./host/boundary -count=1 -timeout 120s -v`: **rc1,1.384s**. `TestWorldBoundaryDependencyAllowlist` names this import; `TestWorldBoundaryRecordingWriter/recording_cmd_ailang-worldd` also fails because this import masks its intended path attribution. Attribution control: remove only this import and replace its simple `Decode` call with equivalent local public-wire `json.Unmarshal`: **both tests PASS,rc0,1.157s**. Restored original bytes. Repair the CLI dependency, preserving the allowlist and tests.
2. **R1-B2 — new documentation breaks existing command binding.** The new literal serve invocation at `docs/QUICKSTART.md:1289` becomes a second command consumed by the existing §11 parser, which scans the remaining document without stopping at the next heading. Independently ran `go test -run '^TestQuickstartSection11FlagsMatchTheCLI$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`: **rc1,0.597s**, named assertion `QUICKSTART §11 holds 2 serve lines, want 1`. With only QUICKSTART restored to base in my isolated tree, the same test **PASS,rc0,0.641s**. Restore the source afterward. A narrowly scoped next-heading cut in the existing test is appropriate if its file is explicitly added to repair scope; do not drop the old assertion or the runnable new recipe.
3. **R1-B3 — new process inherits registry publish authority.** `cmd/ailang-worldd/demo_seed.go:77` launches `exec.CommandContext(...).Output()` with no explicit scrubbed environment. Independently ran `go test -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' ./host/broker -count=1 -timeout 120s -v`: **rc1,2.402s**, exact assertion `10 files carry subprocess sites ... but 9 have drivers`. This is also an observed credential-inheritance defect, not only a driver-count mismatch: the temporary `TestEvaluatorR1VersionProbeInheritsCredential` drove the actual `demoPreflight` branch with a test health server and test pin, replaced `AILANG_REGISTRY_API_KEY` with **only a sentinel**, and the child wrote that sentinel. **rc0,1.261s** confirms this observation; it is not a passing product acceptance test. No real credential was read or dumped. Exact test source and log are banked. Reuse the existing archived interpreter identity path, or explicitly scrub the child with existing `childenv` and add a driven launch-site test; a count/allowlist exception alone is insufficient.

The existing boundary rule concerns direct attribution; the daemon's existing transitive registry dependency does not authorize this new direct import. The credential issue is separate from the read-only grant inspection clarification.

## Named read-only preflight target

The clarification is sound in isolation and implemented as documented. Health identity is matched before `store.OpenReadOnly`; existing authority resolution binds the opaque token to episode `demo`; existing `EffectSpend` debits the current budgets; existing broker predicates check four live grants; bounded `ScanUnreadableWorlds(...,"",1).Scanned` detects any stored World row. Deferred Close executes on return from `demoGrants`, before domain commits. There is no new writable-store handle, Mint path, endpoint or authority in this read-only path. Every public write remains an authenticated daemon REST commit. R1-B3 introduces ambient subprocess authority elsewhere and still blocks the milestone.

I independently killed both the composite nonempty guard and the world-row scan on the **real unselected World / empty selected head / empty log** fixture. Both produce the named inventory assertion showing **5 actual log entries, want0**. This is not the selected-head fixture whose later CAS refusal masks an omitted preflight guard. Bootstrap registry/session rows are accepted in the happy path. No generic all-object emptiness claim is made.

## Focused behavior and restoration

All six original focused commands passed independently, with named tests discovered and no skips:

| Exact test (`go test -run '^NAME$' PACKAGE -count=1 -timeout 120s -v`) | Package | rc | Seconds |
|---|---|---:|---:|
| TestDemoSeedAuthenticatedCommits | ./cmd/ailang-worldd | 0 | 3.023 |
| TestDemoSeedPreflightRefusals | ./cmd/ailang-worldd | 0 | 34.961 |
| TestDemoSeedHelpPreparation | ./cmd/ailang-worldd | 0 | 0.640 |
| TestDemoSeedCommitFailures | ./cmd/ailang-worldd | 0 | 3.881 |
| TestSessionMint_RequiresControllingTerminal | ./cmd/ailang-worldd | 0 | 0.752 |
| TestRESTCommitRefusesDecisionPacketObjects | ./host/daemon | 0 | 1.420 |

The JSON ledgers contain literal argv with `-timeout`, `120s` as separate arguments. Each child preserves global messaging/pins and removes only GOOGLE_CLOUD_PROJECT, OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_RESOURCE_ATTRIBUTES and OTLP_GOOGLE_CLOUD_PROJECT. Each command has a150s outer subprocess bound.

The real daemon acceptance verifies five authenticated domain commits, exact payload/interface hashes, real stored worlds/objects, carrying entries, distinct candidates, unchanged local production source bytes, frozen headers, previous-entry links and observed-head CAS. Help tests execute emitted help forms and assert both pins/flags; preparation documents explicit worktree/descriptor setup and attended publish then typed-y mint. I did not execute either attended confirmation, production mint or publication. A successful S1 stage truthfully exits1 with S2–S6 incomplete and no manifest; no complete demo/capture/release is claimed.

After mutations, independently restored tests passed again: authenticated1.768s, full preflight31.715s, help0.626s, commit failures4.835s. After the later production-revert/attribution controls, happy/help/commit-failure group passed again,rc0,7.658s. Every one of the five implementation files was finally compared byte-for-byte against the reviewed commit; hashes are in `evaluator-r1/restoration.json`. Temporary probe test removed. No source-worktree edits, git checkout/reset/commit/push or protected-file restoration occurred.

## Independent mutations and diff inventory

**20 distinct compiling production mutants killed by named assertions.** No build-red, skip, zero-test, panic or timeout was counted. Full commands/logs are in `evaluator-r1/mutation-ledger.json` and `attack-ledger.json`. An initial commit-response mapping expected `no retries`, but its actual assertion was `commit refusal boundary`; that first expected-assertion mismatch is preserved and was not counted as a predicted kill. I re-ran the identical compiling mutant with the correct named assertion, which failed as predicted. This is20 distinct mutants, not21.

| Diff refusal/support hunk | Independent mutant / named assertion |
|---|---|
| CLI required fields; positional refusal | flags, positional → actual public inventory changed |
| Existing-manifest refusal | manifest → actual public inventory changed |
| Health/database identity | identity → wrong refusal |
| Interpreter hash/version | pin-hash, pin-version → actual public inventory changed |
| Session denial; episode binding | session-denial → wrong refusal; episode → inventory changed |
| Live grant predicates; remaining budget | grant-check, spent-budget → actual public inventory changed |
| Any World row; composite public emptiness | world-scan, unselected-world → actual log entries5want0 |
| Epoch payload integrity; nomination | epoch-hash → wrong refusal; epoch-nomination → inventory changed |
| Before-each-commit CAS | CAS → extra commit request changes public inventory |
| Session-authenticated commit | bearer → real daemon SessionAbsent at successful-stage assertion |
| Commit non200; success-response binding | commit-status, corrected commit-response → commit refusal boundary |
| Frozen header identity; carrying reference | header → frozen header identity mismatch; carrying → no real carrying entry |
| Main dispatch/help, new seed/runner integration | whole own-M2a-production revert below → all3new AC tests red |
| User documentation | base-doc attribution arm → existing§11 regression disappears; R1-B2 |
| Direct dependency and new subprocess | local-wire attribution arm → boundary reds disappear; sentinel observation and R1-B1/B3 |

Rule3o control: removed only the two new M2a production Go files and replaced main.go with its base bytes, leaving the new tests and **all earlier M1a/M1b code** in place. `go test -run '^TestDemoSeed(AuthenticatedCommits|PreflightRefusals|HelpPreparation)$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` compiles and fails all three named tests (**rc1,55.944s**). Each AC test therefore detects this milestone's own production contribution. The inherited mint-TTY and decision-packet tests are regression controls, not claimed as new M2a deliverable killers.

Additional diff branches are explicitly inventoried rather than silently certified. Banked executor logs/29-mutant ledger cover readonly-open, required GET absence/error, unreadable pin, epoch decode, absent source, unknown subcommand and transport uncertainty; I inspected their recorded assertions and exercised their restored fixtures in the32-case/3-outcome runs, but do not relabel all29 as my own mutations. Individual low-level SQLite resolver/spend/scan failures, disappearing-CWD filepath.Abs errors and every inherited worldReader call-site transport error were not independently forced. The fixed primitive commit structs cannot represent a Marshal error through supported input; resolver nil/nil is outside its success-or-denial contract. The90s context budget was not mutation-tested. These are scoped coverage residuals, not an all-error-path or all-platform claim.

## Full gates and CI

No duplicate broad suite was launched. The pristine baseline evidence is useful only as a baseline, never as candidate certification. I read exact-head terminal controller artifacts:

- candidate AIL: **rc0,27.31s**,300s bound, source before/after `67ee2ef`.
- candidate `bash scripts/verify_go.sh`: **rc1,366.33s**,900s bound, source before/after `67ee2ef`. Build passed; plain tests failed on the three mechanisms above (boundary has two failing test names); later race gates were not reached.

PR267's real pull_request CI is pinned to the reviewed SHA. Final observation: **AIL SUCCESS; Go IN_PROGRESS**, two expected jobs. No remote PASS is invented. Controller owns subsequent repaired-head full gates, remote2/2 and merge. The present exact-head hard failures already require FAIL; a pending remote job cannot override them. `evaluator-r1/ci-observed.json` records the observed status.

Scoped quality checks: gofmt prints no filenames,rc0; `go vet ./cmd/ailang-worldd`,rc0; five-file code/doc diffcheck,rc0. Raw evidence whitespace is preserved; whole-evidence diffcheck may show expected raw-log trailing spaces and was not sanitized or confused with the code diffcheck.

## Full rubric

| Category | Score | Basis |
|---|---:|---|
| Tests pass | 0/20 | HARD FAIL: exact-head full Go gate and independent named regressions |
| Lint clean | 10/10 | Affected-package vet and formatting clean; no make lint target in repo profile |
| M2a acceptance | 30/30 | Three narrow behavior ACs independently exercised; no future-stage credit |
| Code quality | 10/15 | Capped-5 for new functions >50lines; no new file >800 or TODO/HACK/FIXME |
| Documentation | 10/15 | No changelog0/5; CLI generator/runbook5/5; honest incremental design state5/5 |
| Design fidelity | 5/10 | Read-only target sound; forbidden dependency and ambient authority regression |
| Conditional regression/performance | N/A | No compiler infrastructure or performance sprint |
| **Total** | **65/100** | **FAIL,3blocking** |

Pending later milestones are intentional and do not incur incomplete-sprint deductions. A repaired commit needs a fresh round; this report certifies neither that repair nor any release product.

EVALUATION_RESULT: fail
EVALUATION_SCORE: 65/100
EVALUATION_ROUND: 1
EVALUATION_REPORT_PATH: design_docs/verification/world-iter268/evaluator-r1.json
FEEDBACK_SUMMARY: Remove forbidden registry dependency, repair QUICKSTART section-boundary regression, and eliminate or scrub/test the new interpreter subprocess.

# Independent evaluator round2 — World268 / row150 / M2a

**FAIL — 70/100; 1 blocking finding.** Subject `7ba691b732d327f6b38a16223152b876e7d197df`, base `39066a46c1ac13c994bacb162b8d4b754d5ec083`, previous candidate `67ee2efe5b2fdca3dcab4eaeda4849a8caea1c8c`. The three r1 blockers are closed; the CI fixture dependency below blocks this repaired head.

Independent native Agent judge `gpt-6.1-sol`, fresh context, neither author nor prior judge. **judge-independence: same-model-fresh-context FLAG**, permitted by D-WORLD-48/current Gate3 after native Sonnet Unknown model and Anthropic/OpenRouter ration limits documented in routing.json. Tokens actually used: unknown. Local measurements are darwin/arm64, Go1.26.6, released v0.41.0 capsule/proof pin and separate v0.52.1 tool pin. No Linux certification is inferred.

## Blocking finding

**R2-B1 — `cmd/ailang-worldd/demo_test.go:394`: ordinary CI cannot satisfy the new pin-mismatch fixture.** The test requires `WORLD_TOOL_AILANG_BIN` and calls Fatal if absent. Independently inspected `.github/workflows/ci.yml`: the Go job installs/exports `AILANG_BIN` only and has no WORLD_TOOL_AILANG_BIN provisioning. Independently ran on pristine exact7ba, removing only that tool variable beyond the standard four telemetry removals:

```sh
go test ./cmd/ailang-worldd -run '^TestDemoSeedPreflightRefusals$/^pin-mismatch$' -count=1 -timeout=90s -v
```

**rc1,2.740s**, exact assertion `demo_test.go:394: required tool pin absent`. Raw evidence: [ci-env-unset-toolpin.log](evaluator-r2/ci-env-unset-toolpin.log). This is a newly introduced test/CI portability failure, not evidence that the CLI's hash rejection is wrong. A full local PASS with both pins set cannot override it. Repair the fixture with byte-distinct temporary pin data, which this hash refusal reaches before any execution; preserve refusal/inventory assertions and verify both environment arms. No skip or workflow modification is needed.

## Prior findings and authority target

- **R1-B1 CLOSED:** forbidden direct host/registry import removed. Existing dependency allowlist and recording-writer tests pass unchanged. The anonymous JSON wire shape preserves the former Decode behavior; fresh decode mutant compiles and dies on wrong-refusal assertion.
- **R1-B2 CLOSED:** all 75,634 prior QUICKSTART bytes are preserved as an exact prefix. New section12 uses its private DEMO_BIN export/build/mint/serve/seed consistently and transfers variables to the serving terminal. Both old §11 binding and `TestQuickstartDecisionPacketSection` pass unchanged; the latter independently rc0,1.175s. New wiring mutation fails its named assertion.
- **R1-B3 CLOSED:** seed's direct subprocess is removed, not exempted from the audit. Existing nine-site driven credential-scrub audit passes. `archive.New(matchedDB).ReadManifest(pinRef)` only reads; startup archive already probes with timeout and childenv.Scrubbed. Fresh missing/hash/size/version-equality/release-prefix mutants compile and fail intended assertions. No credential probe was necessary because the offending launcher no longer exists.

Read-only clarification is sound within its stated scope: health DB identity first, bounded authority resolution binds the opaque credential to episode demo, EffectSpend subtracts prior charges, broker predicates check required live grants, and ScanUnreadableWorlds(empty-afterRef,1).Scanned rejects a real stored unselected World even with empty head/log. Deferred read-only Close happens before REST commits. No writable store, Mint or new endpoint is added. Every public write retains bearer and observed-head CAS. Bootstrap session/registry rows remain allowed; no all-object emptiness claim. Manifest checks bind Hash/Size to local hash-matched bytes and Version to health plus supported release; they reuse the existing startup sidecar trust, not new signature/authenticity authority.

## Own focused gates and mutations

All eight focused commands passed with both pins set; top-level names were discovered, no skips:

| Command/test group | rc | Seconds |
|---|---:|---:|
| regressions | 0 | 5.788 |
| subprocess-audit | 0 | 33.292 |
| TestDemoSeedAuthenticatedCommits | 0 | 4.303 |
| TestDemoSeedPreflightRefusals | 0 | 52.208 |
| TestDemoSeedHelpPreparation | 0 | 0.717 |
| TestDemoSeedCommitFailures | 0 | 3.493 |
| TestSessionMint_RequiresControllingTerminal | 0 | 1.096 |
| TestRESTCommitRefusesDecisionPacketObjects | 0 | 1.461 |

Exact argv are in [focused-gates.json](evaluator-r2/focused-gates.json). Gates use explicit 90s/120s Go bounds and150s outer bounds; children remove only GOOGLE_CLOUD_PROJECT, OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_RESOURCE_ATTRIBUTES, OTLP_GOOGLE_CLOUD_PROJECT. The single-variable CI experiment is separately named above.

Fresh exact-head drills: **34 distinct compiling assertion kills**, including one diagnostic-only wrong-refusal mutation of readonly-open (not a bypass proof). Each attack records landed SHA256, build rc0, named assertion and byte-exact restore. The initial commit-response expected-assertion label was wrong: actual assertion is `commit refusal boundary`; the same mutant was rebuilt and redriven, counted once. A manifest replacement-count0 instrument refusal happened before any write and is excluded. No build-red, panic, timeout or missing test counts as a kill; the old executor29/r1judge20 are not claimed as my current proof.

**Own-M2a diff non-vacuity:** removed the two new production files and restored only main.go's production delta to the M2a base, leaving current tests/docs and M1a/M1b intact. Build rc0. All three named M2a acceptance tests fail; extra commit-failure test fails too. Earlier `TestSessionMint_RequiresControllingTerminal` remains PASS, correctly classed as regression protection; earlier REST packet prohibition is likewise not new M2a behavior credit. [Own-revert log](evaluator-r2/M2a-own-revert.log).

Diff enumeration is exhaustive by path:181 base→head entries (five implementation, two plan clarifications,174 verification/routing),107 previous-candidate→head entries. Full patches/path inventory are banked; no omitted implementation file is assumed from controller orientation. Refusal/support inventory explicitly lists unprovoked filepath.Abs, SQLite resolver/spend/scan and per-call read-error propagation, defensive resolver nil result, and nonrepresentable fixed-primitive JSON marshal errors. These are limitations, not blanket unreachable claims or current independent kills.

Three measured **nonblocking coverage survivors**: replacing90s overall budget with24h leaves all four TestDemoSeed top-level tests PASS; replacing StateRoot derivation with hash(o.Hash) and replacing entryHash with hash(append(encoded,0)) each leaves authenticated test PASS after build0. These identify unpinned supporting behavior; they do not establish violation of a normative canonical hash formula. No suite-wide/sole-killer claim is made for them. Queue them under3n rather than widening this repair. Details/raw logs in JSON and refusal-support inventory.

After every attack, source was restored from external byte copies. Final restored authenticated/preflight/help/commit-failure tests pass, formatting/vet/scoped code-doc diffcheck pass. All five file hashes equal reviewed HEAD.44 protected tracked paths retain head/base SHA256 and newline equality; none was written/restored. Raw-log whitespace remains unaltered.

## Full signals, rubric and limits

Exact-head PR267 CI observation is banked: AILANG job SUCCESS and Go job IN_PROGRESS at observation (no pending result called PASS). Controller exact7ba AIL gate rc0 is copied; fullGo status at final snapshot is in JSON. I did not duplicate full suites. The independently reproduced ordinary-CI test failure already makes Tests Pass0/HARD FAIL, regardless of a both-pin local result. Historical base green and old67ee fullGo red remain scoped to those heads.

| Category | Score | Basis |
|---|---:|---|
| Tests pass | 0/20 | HARD FAIL: new test fails ordinary-CI environment |
| Lint clean | 10/10 | affected-package vet, formatting and scoped diffcheck green |
| M2a acceptance | 30/30 | three narrow behavior criteria independently exercised with both pins; CI defect separately blocks |
| Code quality | 10/15 | capped-5 new functions >50lines; no >800line file or new TODO/HACK/FIXME |
| Documentation | 10/15 | no changelog0/5; runnable generator/runbook5/5; honest incremental status5/5 |
| Design fidelity | 10/10 | repaired architecture and named read-only clarification match narrow M2a |
| Conditional compiler/performance | N/A | neither category triggered |
| **Total** | **70/100 FAIL** | **1blocking** |

Only M2a is judged. M1a/M1b were already landed; M2b+ remain intentionally pending, without false incomplete-sprint deductions. Quorum r1–r3 stay BLOCKED; r4 inherits iteration266 narrow-refinement disposition, not a new quorum PASS. No whole seed/demo, proof stage, images, capture, release, publication, attended mint or promotion credit. No git writes, inbox calls, production publish/mint, confirmations or subagents were used. All own commands are terminal and source restored.

EVALUATION_RESULT: fail
EVALUATION_SCORE: 70/100
EVALUATION_ROUND: 2
EVALUATION_REPORT_PATH: design_docs/verification/world-iter268/evaluator-r2.json
FEEDBACK_SUMMARY: Replace tool-pin environment dependence in pin-mismatch fixture; preserve all assertions and verify ordinary-CI unset arm.

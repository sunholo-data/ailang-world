# M4b executor report — iteration 275

**Status: functional implementation and all requested gates passed; strict milestone acceptance remains INCOMPLETE for the two findings below.** No success manifest was written.

## Files changed

- `cmd/ailang-worldd/demo_seed.go`
- `cmd/ailang-worldd/demo_test.go`

Authoritative copies retain their relative paths beneath this directory. `why_test.go` needed no edit. No git write operations were performed. `why.go`, `demo.go`, `demo_proof.go`, `main.go`, host, packages, world, and all other tracked files are unchanged.

## Behavior

S6 calls the actual CLI producer, worldReader.locate/buildChain. Its pure validator checks coordinator first, then OK, nonzero links, effect presence, and world/entry/record/input/plan/effect/output presence and truth. Errors identify the entry and failed condition. Ordinary REST replacement reaches coordinator refusal before packet count.

S6 requires exactly one packet, content-checks it over REST, decodes its canonical packet bytes, and binds its proposal and S5 index to the existing daemon packet validator via its workbench GET verdict. That existing validator supplies carrying-commit membership and approval-chain/request binding. Manifest resolution separately fetches packet, request, proposal, input, invocation record, effects, effect requests/results, and output by hash over REST. This adds no second why walker.

Assembly returns only an in-memory public structure: schema/version, proof/tool releases, frozen head, entries ordered by index and their actual object refs, goal/candidate/source/claim refs, four evidence refs, packet/request/proposal, invocation/transition/interpreter/effect refs, and four raw why JSON SHA256 digests. JSON bytes match runWhy --json formatting including its newline. References are independently reread and checked, including world and entry hash recomputation and object payload hashes. Public-byte scanning rejects the session token, credential/key markers, absolute paths and timestamps; schema/pin and role allowlists also prevent arbitrary free-text public fields. Raw input, reports, credentials, paths and capture times are not serialized into the assembly.

The final seed return is always an error: `S1-S6 checked; roster checks M5c pending; no success manifest written`. No success output path was added. Existing S4/S5 integration assertions were updated to this boundary. The isolated prior S5-only helper/control retains its truthful S6-incomplete error.

NB2 now returns a named S1 payload decode error for each of the goal and two candidate payloads. NB3 requires a present matching test-entry header.entryIndex and a present integral A2A metadata.entry_index. Prior positive A2A unit fixtures now explicitly carry index zero; missing-invocation controls carry a valid index so the intended older guard remains observable.

## Gates

Both required pins were exported for all pinned gates. All runs occurred in the sandbox; no bind denial or other sandbox denial occurred. Only the two requested full TestDemo invocations were run. No repository-wide suite or verify scripts were run.

| Command | rc | Named evidence |
|---|---:|---|
| `go vet ./cmd/ailang-worldd` | 0 | No diagnostics; named PASS not applicable |
| `gofmt -l cmd/ailang-worldd` | 0 | Empty output; named PASS not applicable |
| `go test -run '^TestDemoSeedWhyNamedLinks$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | 0 | `--- PASS: TestDemoSeedWhyNamedLinks (11.20s)` |
| `go test -run '^TestDemoSeedManifestReferences$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | 0 | `--- PASS: TestDemoSeedManifestReferences (10.63s)` |
| `go test -run '^TestWhyTargetFormsResolveToOneEntry$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` | 0 | `--- PASS: TestWhyTargetFormsResolveToOneEntry (0.04s)` |
| `go test -run '^TestCLIAgainstSeToolsDaemon$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` | 0 | `--- PASS: TestCLIAgainstSeToolsDaemon (9.15s)` |
| `go test -run ^TestDemo ./cmd/ailang-worldd -count=1 -timeout 500s -v` | 0 | 26 named top-level PASS; zero FAIL/SKIP; 154.38 s |
| `WORLD_TOOL_AILANG_BIN= go test -run ^TestDemo ./cmd/ailang-worldd -count=1 -timeout 500s -v` | 0 | 23 named top-level PASS; zero FAIL; only tool-gated SKIP below; 51.12 s |

The full pinned gate preceded the final test-only assertion that gives the real plain-REST mutant its exact coordinator diagnostic. The restored real-stack focused run exercised that final assertion and passed; the subsequent full unset gate used the final files. No production behavior changed after the full pinned gate.

Pinned full-suite named evidence:

```text
--- PASS: TestDemoProofHTTPProducerRoundTrip (1.42s)
--- PASS: TestDemoProofSourceHashRefusal (1.11s)
--- PASS: TestDemoProofAdapterGuards (0.00s)
--- PASS: TestDemoProofWrongKeyCannotProve (1.32s)
--- PASS: TestDemoProofNoVerifiedSetRefuses (1.39s)
--- PASS: TestDemoProofUncertainCommitStops (4.62s)
--- PASS: TestDemoProofBounds (1.74s)
--- PASS: TestDemoProofKeyIORefusals (0.00s)
--- PASS: TestDemoProofKeyDisclosure (0.00s)
--- PASS: TestDemoSeedAuthenticatedCommits (1.43s)
--- PASS: TestDemoSeedPreflightRefusals (50.62s)
--- PASS: TestDemoSeedHelpPreparation (0.00s)
--- PASS: TestDemoSeedCommitFailures (3.74s)
--- PASS: TestDemoSeedReadTestCheck (10.56s)
--- PASS: TestDemoSeedRejectsToolReportRefusals (5.91s)
--- PASS: TestDemoSeedToolGuardControls (0.00s)
--- PASS: TestDemoSeedProofStage (42.31s)
--- PASS: TestDemoSeedRequiredIdentities (0.00s)
--- PASS: TestDemoSeedApprovalControls (0.00s)
--- PASS: TestDemoSeedApprovalBudgetControl (0.00s)
--- PASS: TestDemoSeedApprovalPacket (16.16s)
--- PASS: TestDemoSeedWhyNamedLinks (5.61s)
--- PASS: TestDemoSeedManifestReferences (5.58s)
--- PASS: TestDemoSeedM4bBoundaryControls (0.00s)
--- PASS: TestDemoSeedPacketControls (0.01s)
--- PASS: TestDemoSeedPacketObjectControls (0.04s)
```

Unset-tool full-suite named evidence:

```text
--- PASS: TestDemoProofHTTPProducerRoundTrip (1.21s)
--- PASS: TestDemoProofSourceHashRefusal (1.03s)
--- PASS: TestDemoProofAdapterGuards (0.00s)
--- PASS: TestDemoProofWrongKeyCannotProve (1.25s)
--- PASS: TestDemoProofNoVerifiedSetRefuses (1.37s)
--- PASS: TestDemoProofUncertainCommitStops (3.40s)
--- PASS: TestDemoProofBounds (1.03s)
--- PASS: TestDemoProofKeyIORefusals (0.00s)
--- PASS: TestDemoProofKeyDisclosure (0.00s)
--- PASS: TestDemoSeedAuthenticatedCommits (1.16s)
--- PASS: TestDemoSeedPreflightRefusals (36.97s)
--- PASS: TestDemoSeedHelpPreparation (0.00s)
--- PASS: TestDemoSeedCommitFailures (2.86s)
--- PASS: TestDemoSeedRejectsToolReportRefusals (0.00s)
--- PASS: TestDemoSeedToolGuardControls (0.00s)
--- PASS: TestDemoSeedRequiredIdentities (0.00s)
--- PASS: TestDemoSeedApprovalControls (0.00s)
--- PASS: TestDemoSeedApprovalBudgetControl (0.00s)
--- PASS: TestDemoSeedWhyNamedLinks (0.04s)
--- PASS: TestDemoSeedManifestReferences (0.03s)
--- PASS: TestDemoSeedM4bBoundaryControls (0.00s)
--- PASS: TestDemoSeedPacketControls (0.02s)
--- PASS: TestDemoSeedPacketObjectControls (0.04s)
```

Unset-tool skips (each is intentional tool gating; none counts as integrated acceptance):

```text
--- SKIP: TestDemoSeedReadTestCheck (0.00s)
    --- SKIP: TestDemoSeedRejectsToolReportRefusals/real-refusal-arms (0.00s)
--- SKIP: TestDemoSeedProofStage (0.02s)
--- SKIP: TestDemoSeedApprovalPacket (0.01s)
    --- SKIP: TestDemoSeedWhyNamedLinks/real-stack (0.02s)
    --- SKIP: TestDemoSeedManifestReferences/real-stack (0.02s)
```

Additional focused controls passed: TestDemoSeedM4bBoundaryControls, TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls, TestDemoSeedToolGuardControls. Restored mutation logs include the required named top-level PASS for each targeted run, never zero-test or skipped green.

## Mutation ledger

All mutations were actual source edits followed by the named test. Every completed mutation was restored from an external scratch copy, checked with cmp, and rerun green. No git restoration was used. Locations below are the observed file:line at mutation time; pre-mutation source copies preserve those revisions. Four initial compile failures are retained as non-evidence, followed by corrected reruns. `_ = local` no-ops in those corrected mutants prevent Go unused-local errors while the actual guard is replaced with `if false {`.

| Mutant | Exact command | Red rc and observed assertion | Restored |
|---|---|---|---|
| coordinator | `go test -run ^TestDemoSeedWhyNamedLinks$/controls/coordinator$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1624: coordinator must be true: entry 1: ok must be true | cmp equal; rc 0; --- PASS: TestDemoSeedWhyNamedLinks (0.03s) |
| why-ok | `go test -run ^TestDemoSeedWhyNamedLinks$/controls/ok$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1624: ok must be true: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedWhyNamedLinks (0.02s) |
| zero-links | `go test -run ^TestDemoSeedWhyNamedLinks$/controls/zero-links$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1624: zero links refused: entry 1: expected effect link must exist and be true | cmp equal; rc 0; --- PASS: TestDemoSeedWhyNamedLinks (0.02s) |
| link-false | `go test -run ^TestDemoSeedWhyNamedLinks$/controls/effect-false$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1624: expected effect link must exist and be true: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedWhyNamedLinks (0.02s) |
| link-absent | `go test -run ^TestDemoSeedWhyNamedLinks$/controls/effect-absent$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | UNINFORMATIVE compile attempt; rc 1; cmd/ailang-worldd/demo_seed.go:629:3: declared and not used: found | cmp equal; rc 0; --- PASS: TestDemoSeedWhyNamedLinks (0.02s) |
| link-absent | `go test -run ^TestDemoSeedWhyNamedLinks$/controls/effect-absent$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1624: expected effect link must exist and be true: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedWhyNamedLinks (0.02s) |
| metadata-index | `go test -run ^TestDemoSeedM4bBoundaryControls$/metadata-index$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1876: entry_index must be present and integral: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedM4bBoundaryControls (0.00s) |
| s1-decode | `go test -run ^TestDemoSeedM4bBoundaryControls$/s1-decode/0$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1862: S1 payload 0 decode refused: S1-S4 recorded; S5 test invocation unavailable: required resource absent: /v1/log/0; no success manifest written | cmp equal; rc 0; --- PASS: TestDemoSeedM4bBoundaryControls (0.00s) |
| test-index | `go test -run ^TestDemoSeedM4bBoundaryControls$/test-index$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1871: test entry entryIndex must be present and match: S5 refused: demo-approval: tasks/send refused (HTTP 200); no success manifest written | cmp equal; rc 0; --- PASS: TestDemoSeedM4bBoundaryControls (0.00s) |
| public-bytes | `go test -run ^TestDemoSeedManifestReferences$/controls/session-token$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1835: public bytes contain: public manifest schema/version/pins must match | cmp equal; rc 0; --- PASS: TestDemoSeedManifestReferences (0.02s) |
| manifest-schema | `go test -run ^TestDemoSeedManifestReferences$/controls/schema-schema$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1817: public manifest schema/version/pins must match: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedManifestReferences (0.03s) |
| manifest-role | `go test -run ^TestDemoSeedManifestReferences$/controls/role$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1824: public manifest role must be a declared reference name: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedManifestReferences (0.02s) |
| manifest-ref | `go test -run ^TestDemoSeedManifestReferences$/controls/absent-ref$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | UNINFORMATIVE compile attempt; rc 1; cmd/ailang-worldd/demo_seed.go:784:3: declared and not used: valid | cmp equal; rc 0; --- PASS: TestDemoSeedManifestReferences (0.03s) |
| manifest-ref | `go test -run ^TestDemoSeedManifestReferences$/controls/absent-ref$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1768: every declared manifest reference resolves to matching content: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedManifestReferences (0.02s) |
| why-hash | `go test -run ^TestDemoSeedManifestReferences$/controls/why-hash$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | UNINFORMATIVE compile attempt; rc 1; cmd/ailang-worldd/demo_seed.go:804:7: declared and not used: raw | cmp equal; rc 0; --- PASS: TestDemoSeedManifestReferences (0.04s) |
| why-hash | `go test -run ^TestDemoSeedManifestReferences$/controls/why-hash$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1829: why raw JSON hash mismatch: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedManifestReferences (0.03s) |
| packet-count | `go test -run ^TestDemoSeedPacketControls$/packet-count$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1889: entry 1: exactly one decision packet required: packet: does not resolve | cmp equal; rc 0; --- PASS: TestDemoSeedPacketControls (0.04s) |
| packet-anchor | `go test -run ^TestDemoSeedPacketControls$/carrying-entry$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1899: entry 7: packet carrying entry/proposal mismatch: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedPacketControls (0.00s) |
| s6-fail-closed | `go test -run ^TestDemoSeedM4bBoundaryControls$/s6-fail-closed$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1879: S1-S6 checked; roster checks M5c pending; no success manifest written: <nil> | cmp equal; rc 0; --- PASS: TestDemoSeedM4bBoundaryControls (0.00s) |
| packet-object | `go test -run ^TestDemoSeedPacketObjectControls$/packet-object$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | UNINFORMATIVE compile attempt; rc 1; cmd/ailang-worldd/demo_seed.go:847:7: declared and not used: ok | cmp equal; rc 0; --- PASS: TestDemoSeedPacketObjectControls (0.03s) |
| packet-object | `go test -run ^TestDemoSeedPacketObjectControls$/packet-object$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1961: packet: does not resolve: coordinator: not a canonical world/decision-packet/v1 encoding: EOF | cmp equal; rc 0; --- PASS: TestDemoSeedPacketObjectControls (0.03s) |
| packet-decode | `go test -run ^TestDemoSeedPacketObjectControls$/packet-decode$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1961: not a canonical world/decision-packet/v1 encoding: entry 1: packet carrying entry/proposal mismatch | cmp equal; rc 0; --- PASS: TestDemoSeedPacketObjectControls (0.04s) |
| M4b-plain-rest | `go test -run ^TestDemoSeedWhyNamedLinks$/controls/approval-stimulus$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1609: entry 2: coordinator must be true | cmp equal; rc 0; --- PASS: TestDemoSeedWhyNamedLinks (0.02s) |
| M4b-link | `go test -run ^TestDemoSeedWhyNamedLinks$/controls/accepted-report$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1598: entry 1: expected effect link must exist and be true | cmp equal; rc 0; --- PASS: TestDemoSeedWhyNamedLinks (0.01s) |
| M4b-ref | `go test -run ^TestDemoSeedManifestReferences$/controls$ ./cmd/ailang-worldd -count=1 -timeout 45s -v` | rc 1; demo_test.go:1780: every declared manifest reference resolves to matching content: every declared manifest reference resolves to matching content: input sha256:d0ee1796b2a9b04b1c11e34b0bf6d14af861bcff46f301886b127104217a00f4 | cmp equal; rc 0; --- PASS: TestDemoSeedManifestReferences (0.02s) |
| M4b-plain-rest-real-stack | `go test -run ^TestDemoSeedWhyNamedLinks$/real-stack$ ./cmd/ailang-worldd -count=1 -timeout 180s -v` | rc 1; demo_test.go:1674: coordinator must be true for every real seed entry: ailang-worldd demo seed: S6 refused: entry 9: coordinator must be true; no success manifest written | cmp equal; rc 0; --- PASS: TestDemoSeedWhyNamedLinks (5.19s) |

Plan details: M4b-link removes only effect from a decoded real CLI report, preserving coordinator/OK/other links. M4b-ref replaces a public input ref with a syntactically valid absent SHA256 ref and reaches the positive reference-resolution assertion. The pin-free M4b-plain-rest stimulus is an ordinary REST object. The additional **real-stack** M4b-plain-rest mutant replaces the actual seed S5 A2A dispatch with a valid ordinary REST goal-object commit (no forbidden packet semantic ID), rereads and compares that committed entry, and lets actual S6 validate it. It fails at `coordinator must be true for every real seed entry`, with S6 naming entry 9, before any packet-count failure. Permanent negative subtests retain ordinary REST refusal and effect omission.

Controls also compare immutable REST content before/after the mismatched-read proxy and compare the preceding coordinator entry before/after the ordinary REST commit, preventing daemon behavior from hiding those checks.

## Protected-file observation

Before/after SHA256 and newline-count inventories in observation-before.json and observation-after.json are identical (cmp rc 0). The scoped read-only `git diff --exit-code -- world go.mod go.sum tools/launchd design_docs/world-mission.md` returned 0. `git diff --check` returned 0. Only the two allowlisted files appear in the read-only changed-path listing.

## Findings and unfinished acceptance

1. **Pin resolution boundary:** manifest interpreter hashes are cross-checked against `/v1/health`; source, transition, invocation, evidence and effect hashes are fetched from the object store. The frozen REST surface has no route returning archived interpreter/tool binary bytes by hash. The proof pin bytes/archive identity are already checked by preflight, and the configured tool release is checked by the daemon/test setup. This does **not** discharge a strict demand to fetch and rehash every archived binary pin over REST. The assembly includes tool/proof releases and interpreter hash, while the tool binary hash remains within its referenced recorded output. No host/API file was changed to work around the allowlist. A strict binary-pin-bytes requirement needs a follow-up design/API decision.

2. **Exhaustive branch mutation demand:** the ledger discharges 18 distinct semantic guards/refusals plus all three plan mutants (including an additional actual-seed plain-REST run). Ordinary newly added transport/codec error-propagation returns have not each been independently bypassed and mutation-checked. Therefore the requested “EVERY guard/refusal branch” standard is not claimed complete. Existing real-stack and packet object/decode controls cover those paths only partially. Further focused error-propagation controls/mutations remain unfinished.

3. Manifest publication and M5c roster-grade checks intentionally remain pending; all tested seeded paths still refuse success and write no manifest. This is the required boundary, not an acceptance claim for later milestones.

No unrelated gate failures were observed. No sandbox result was mislabeled PASS/FAIL. No unfinished work is labeled complete.

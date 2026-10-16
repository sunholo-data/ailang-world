# M4b r2 executor report — iteration 275

Updated: 2026-10-16T16:19:45.278516+00:00

**Status: COMPLETE for F1–F5 and the requested gates/mutants.**

The first snapshot and report were confirmed present by 16:14:58 UTC, within 6 minutes 26 seconds of the first recorded task time (16:08:32 UTC). Final files are recopied after validation.

## Scope

Only `cmd/ailang-worldd/demo_seed.go` and `cmd/ailang-worldd/demo_test.go` were modified. Copies retain those relative paths beneath this report. No git write operation, full repository suite, verify script, host/API change, or plan edit was performed. All logs and mutation source copies are in `logs/`.

## Findings addressed

- F1: validator enforces exactly one frozen-head, goal, candidate-a, candidate-b, source, claim, each of test-invocation/recorded-effect/authored-claim/proof-envelope, packet, request and proposal. Why digests carry read/test/check/approval roles, each required exactly once with distinct entries. Per-run refs carry their entry index and must match the exact invocation, transition, interpreter, why links and effect request/result multiset resolved from that entry. Errors name missing/extra roles. Pin-free controls assemble from a unit daemon, independently assert the exact role multiset, then remove and duplicate every required static and per-run reference and each why role in table-driven subtests. The real-stack arm independently asserts the exact assembled role multiset too.
- F2: pin-free S6 proxy lists two distinct canonical packets and requires the exactly-one-packet error.
- F3: pin-free S6 controls serve workbench HTML with no carrying-entry anchor and with entry #99 instead of #4, requiring the packet carrying-entry/proposal refusal.
- F4: pin-free boundary control first proves S6 succeeds, then calls demoSeedApproval through the same unit daemon/proxy and requires the final S1–S6/M5c-pending error. This reaches the actual seed return with the tool pin unset; no fallback to testing the error helper alone is used.
- F5: host/coordinator is in the project import block.

The unit daemon stores canonical content and coordinator-shaped invocation fixtures. Domain/transition content is committed over REST. Canonical packet fixture bytes are inserted only while the daemon is closed (the public writer correctly forbids coordinator-only packet minting), then the daemon restarts and every manifest object/world/entry is resolved over real REST. The proxy supplies synthetic interpreter health, A2A completion and controlled workbench/listing responses. This tests seed orchestration without claiming the unit fixture runs the real tool or coordinator; the pinned real-stack arm supplies that integration coverage.

## Gates

Both pins are set for pinned commands: `AILANG_BIN=$HOME/.pinned-ailang/ailang`, `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.52.1/ailang`. Unset-mode commands remove WORLD_TOOL_AILANG_BIN from the environment (the requested empty assignment has the same test gating behavior).

| Command | Tool pin | rc | Named evidence |
|---|---|---:|---|
| `go vet ./cmd/ailang-worldd` | pinned | 0 | no diagnostics; named PASS not applicable |
| `gofmt -l cmd/ailang-worldd` | pinned | 0 | empty output; named PASS not applicable |
| `go test -run '^TestDemoSeedWhyNamedLinks$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | pinned | 0 | --- PASS: TestDemoSeedWhyNamedLinks (7.27s) |
| `go test -run '^TestDemoSeedManifestReferences$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | pinned | 0 | --- PASS: TestDemoSeedManifestReferences (10.74s) |
| `go test -run '^TestDemo' ./cmd/ailang-worldd -count=1 -timeout 500s -v` | pinned | 0 | 26 named top-level PASS lines, reproduced below |
| `go test -run '^TestDemo' ./cmd/ailang-worldd -count=1 -timeout 500s -v` | unset | 0 | 23 named top-level PASS lines, reproduced below |

demo-pinned named PASS evidence:

```text
--- PASS: TestDemoProofHTTPProducerRoundTrip (2.62s)
--- PASS: TestDemoProofSourceHashRefusal (2.86s)
--- PASS: TestDemoProofAdapterGuards (0.00s)
--- PASS: TestDemoProofWrongKeyCannotProve (2.46s)
--- PASS: TestDemoProofNoVerifiedSetRefuses (1.66s)
--- PASS: TestDemoProofUncertainCommitStops (4.42s)
--- PASS: TestDemoProofBounds (1.33s)
--- PASS: TestDemoProofKeyIORefusals (0.00s)
--- PASS: TestDemoProofKeyDisclosure (0.00s)
--- PASS: TestDemoSeedAuthenticatedCommits (1.76s)
--- PASS: TestDemoSeedPreflightRefusals (56.05s)
--- PASS: TestDemoSeedHelpPreparation (0.00s)
--- PASS: TestDemoSeedCommitFailures (4.48s)
--- PASS: TestDemoSeedReadTestCheck (4.90s)
--- PASS: TestDemoSeedRejectsToolReportRefusals (4.37s)
--- PASS: TestDemoSeedToolGuardControls (0.00s)
--- PASS: TestDemoSeedProofStage (33.89s)
--- PASS: TestDemoSeedRequiredIdentities (0.00s)
--- PASS: TestDemoSeedApprovalControls (0.00s)
--- PASS: TestDemoSeedApprovalBudgetControl (0.00s)
--- PASS: TestDemoSeedApprovalPacket (14.06s)
--- PASS: TestDemoSeedWhyNamedLinks (5.65s)
--- PASS: TestDemoSeedManifestReferences (7.81s)
--- PASS: TestDemoSeedM4bBoundaryControls (0.07s)
--- PASS: TestDemoSeedPacketControls (0.07s)
--- PASS: TestDemoSeedPacketObjectControls (0.02s)
```

demo-pinned SKIPs (not counted as acceptance):

```text
none
```

demo-unset named PASS evidence:

```text
--- PASS: TestDemoProofHTTPProducerRoundTrip (1.68s)
--- PASS: TestDemoProofSourceHashRefusal (2.17s)
--- PASS: TestDemoProofAdapterGuards (0.00s)
--- PASS: TestDemoProofWrongKeyCannotProve (3.21s)
--- PASS: TestDemoProofNoVerifiedSetRefuses (2.37s)
--- PASS: TestDemoProofUncertainCommitStops (4.43s)
--- PASS: TestDemoProofBounds (1.36s)
--- PASS: TestDemoProofKeyIORefusals (0.00s)
--- PASS: TestDemoProofKeyDisclosure (0.00s)
--- PASS: TestDemoSeedAuthenticatedCommits (1.73s)
--- PASS: TestDemoSeedPreflightRefusals (56.09s)
--- PASS: TestDemoSeedHelpPreparation (0.00s)
--- PASS: TestDemoSeedCommitFailures (4.45s)
--- PASS: TestDemoSeedRejectsToolReportRefusals (0.00s)
--- PASS: TestDemoSeedToolGuardControls (0.00s)
--- PASS: TestDemoSeedRequiredIdentities (0.00s)
--- PASS: TestDemoSeedApprovalControls (0.00s)
--- PASS: TestDemoSeedApprovalBudgetControl (0.00s)
--- PASS: TestDemoSeedWhyNamedLinks (0.04s)
--- PASS: TestDemoSeedManifestReferences (2.19s)
--- PASS: TestDemoSeedM4bBoundaryControls (0.08s)
--- PASS: TestDemoSeedPacketControls (0.07s)
--- PASS: TestDemoSeedPacketObjectControls (0.03s)
```

demo-unset SKIPs (not counted as acceptance):

```text
--- SKIP: TestDemoSeedReadTestCheck (0.00s)
    --- SKIP: TestDemoSeedRejectsToolReportRefusals/real-refusal-arms (0.00s)
--- SKIP: TestDemoSeedProofStage (0.02s)
--- SKIP: TestDemoSeedApprovalPacket (0.02s)
    --- SKIP: TestDemoSeedWhyNamedLinks/real-stack (0.02s)
    --- SKIP: TestDemoSeedManifestReferences/real-stack (0.02s)
```

## Mutation ledger

Every mutant is an actual source edit, restored from `/tmp/world275-r2-clean/` in a finally block; both allowlisted files are compared with `cmp` after each restore. The same named test is rerun on restored source. No skip, zero-test run, panic, compile failure, or earlier guard is counted as a kill. The complete source-edit before/after strings and all failing assertion lines are in `logs/mutations.json`; raw red/restored logs and mutated source copies are retained.

| Mutant | Command | Tool pin | rc / failing assertion | Restored |
|---|---|---|---|---|
| O6-unset | `go test -run '^TestDemoSeedManifestReferences$/controls$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | unset | rc 1; demo_test.go:1973: manifest completeness: role proposal requires 1 references, got 0 | cmp seed rc 0; cmp test rc 0; rerun rc 0; --- PASS: TestDemoSeedManifestReferences (2.12s) |
| O6-pinned-real | `go test -run '^TestDemoSeedManifestReferences$/real-stack$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | pinned | rc 1; demo_test.go:2087: manifest completeness: real seed assembly must contain every required role: ailang-worldd demo seed: S6 refused: manifest completeness: role packet requires exactly 1 reference, got 0; no success manifest written | cmp seed rc 0; cmp test rc 0; rerun rc 0; --- PASS: TestDemoSeedManifestReferences (5.52s) |
| O10-unset | `go test -run '^TestDemoSeedManifestReferences$/controls$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | unset | rc 1; demo_test.go:1973: manifest completeness: role recorded-effect requires 1 references, got 0 | cmp seed rc 0; cmp test rc 0; rerun rc 0; --- PASS: TestDemoSeedManifestReferences (2.41s) |
| O10-pinned-real | `go test -run '^TestDemoSeedManifestReferences$/real-stack$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | pinned | rc 1; demo_test.go:2087: manifest completeness: real seed assembly must contain every required role: ailang-worldd demo seed: S6 refused: manifest completeness: role test-invocation requires exactly 1 reference, got 0; no success manifest written | cmp seed rc 0; cmp test rc 0; rerun rc 0; --- PASS: TestDemoSeedManifestReferences (5.41s) |
| completeness-unset | `go test -run '^TestDemoSeedManifestReferences$/controls$/completeness$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | unset | rc 1; demo_test.go:1986: manifest completeness: role frozen-head: <nil>; demo_test.go:1991: manifest completeness: role frozen-head: <nil>; 122 total named assertion failures in log | cmp seed rc 0; cmp test rc 0; rerun rc 0; --- PASS: TestDemoSeedManifestReferences (2.55s) |
| O5-two-packets-unset | `go test -run '^TestDemoSeedPacketControls$/two-packets$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | unset | rc 1; demo_test.go:2141: entry 4: exactly one decision packet required: <nil> | cmp seed rc 0; cmp test rc 0; rerun rc 0; --- PASS: TestDemoSeedPacketControls (0.04s) |
| O8-anchor-unset | `go test -run '^TestDemoSeedPacketControls$/(missing-anchor\|wrong-anchor)$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | unset | rc 1; demo_test.go:2151: entry 4: packet carrying entry/proposal mismatch: <nil>; demo_test.go:2151: entry 4: packet carrying entry/proposal mismatch: <nil> | cmp seed rc 0; cmp test rc 0; rerun rc 0; --- PASS: TestDemoSeedPacketControls (0.06s) |
| O7-seed-return-unset | `go test -run '^TestDemoSeedM4bBoundaryControls$/s6-fail-closed$' ./cmd/ailang-worldd -count=1 -timeout 180s -v` | unset | rc 1; demo_test.go:2133: S1-S6 checked; roster checks M5c pending; no success manifest written: <nil> | cmp seed rc 0; cmp test rc 0; rerun rc 0; --- PASS: TestDemoSeedM4bBoundaryControls (0.10s) |

## Unfinished / limits

None of the requested F1–F5 fixes, gates or mutants remains unfinished. M5c roster grading and manifest publication remain outside this task; seed still refuses success.

Development fixture setup initially hit the existing coordinator-only packet restriction, then the sole-writer lock, then noncanonical hashes in the older generic commit fixture. These were corrected in the test fixture (closed-daemon insertion and canonical commit hashing). Those development runs are not gate passes or mutation evidence. No sandbox denial occurred.

Final verification: both snapshot copies and both mutation scratch copies match the worktree byte-for-byte (`cmp` rc 0 for all four comparisons). `git diff --check` rc 0; changed tracked paths are exactly the two allowlisted files. Hashes and command results are in `logs/final-verification.json`. All started commands have completed.

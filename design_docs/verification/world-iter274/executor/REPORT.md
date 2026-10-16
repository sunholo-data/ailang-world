# Iteration 274 — row 150 M4a executor report

Completed in the assigned worktree. Started 2026-10-16 11:29:26 UTC; report written 2026-10-16 11:48:23 UTC.

## Files changed

- `cmd/ailang-worldd/demo_seed.go` — snapshot at the same relative path; SHA-256 `d234cdb450ad88872debc011ef8260e572886e0d548851c8cbf28e890db2134d`.
- `cmd/ailang-worldd/demo_test.go` — snapshot at the same relative path; SHA-256 `ef2b583d89dd02367d6efb186ac95f8b2eaf234ce80aa26a44ea40307f847ba1`.

Only these two worktree files changed. No git write operations, host/package/kernel edits, full repository suite, verify_go.sh or verify_ail.sh were performed.

## S5 behavior

Successful S4 now returns its actual producer envelope reference to the seed. The seed retains real S2/S3 results and the S3 test entry index, fetches that entry over REST to obtain its invocation-record content hash, and sends exactly one `/a2a/` `tasks/send`: task `demo-approval`, skill `demo-approve`.

The input contains the goal and two proposal description strings from the S1 objects actually committed, and four ordered kind/ref pairs: S3 test invocation record, S2 read effect record, S1 demo-authored claim object, S4 authenticated envelope. Roster length and kinds are checked before dispatch. No packet is fabricated, no decision verb is called, and no invocation is retried.

Approval runs through the existing real pinned capsule → coordinator → bound broker invoker → daemon HumanHandler → packet mint/commit path. After S5 the command still exits with an error: `S1-S5 recorded; S6 incomplete; no success manifest written; inspect history before any further action`. Public manifest serialization remains outside M4a.

## Acceptance and controls

`TestDemoSeedApprovalPacket` publishes the checked-in demo descriptor with the production transition publisher, alongside se-tools in the existing test harness. The positive arm uses unchanged checked-in approve source and the real runner, coordinator and daemon handler. It verifies one S5 entry/effect/packet; reads every linked object over REST and recomputes its content hash; verifies e1, plan, input, output, effect record, framed effect-request bytes, effect result and pending approval request; checks the published descriptor access cost is zero; and compares all four roster refs and goal/candidate descriptions with independently read committed stages.

The daemon packet page invokes its real packet validator; the test requires its first carrying entry to equal S5 and its approvals-chain status to be pending. Named negative arms are `single-approval/no-decision-verb-packet-stays-pending`, `single-approval/seed-invokes-s5-once`, `no-human-budget`, and `extra-effect-refused`. Missing Human.Approve budget leaves zero packets, zero A2A approval calls and unchanged empty REST history. The valid two-effect temp module preserves the descriptor schemas and declared capability, has distinct e1/e2, typechecks under the proof pin, and is refused with no packet, spend or S5 commit. Its operator diagnostic is `coordinator: effect plan violates L1: 2 effects; a v1 plan carries at most 1`.

Budget observation uses the existing REST content-addressed accounting read path: read the immutable effect records immediately before S5 and after it, subtract actual allowed Human.Approve costs from the fixture’s initial minted grant, then cross-check those independently observed balances against the S5 broker record’s BudgetBefore/BudgetAfter (also fetched by hash). Both checks require an exact delta of one. There is no standalone REST grant-balance endpoint; none was invented or added.

Pin-free coverage is provided by `TestDemoSeedApprovalControls` (stage-ref assembly, wrong roster length/kind, refusal before dispatch, missing S3 entry, S5 refusal, S6 incomplete text) and `TestDemoSeedApprovalBudgetControl` (rejects deltas zero, two and minus one). Tests write fixture files only under t.TempDir(); the shared real-tool rig also sets TMPDIR to t.TempDir so production capsule/source-check scratch roots stay there. All affected gates were rerun after adding this containment. The new real integration test reuses cliSeBins for unset-tool skipping and keeps the required proof-pin failure behavior.

## Required gates

Configured-pin commands used AILANG_BIN=$HOME/.pinned-ailang/ailang (v0.41.0) and WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.52.1/ailang. No required gate encountered a sandbox denial. Empty outputs below are non-test gates, not vacuous test acceptance.

| Gate | rc | Wall seconds | Named result |
|---|---:|---:|---|
| `go vet ./cmd/ailang-worldd` | 0 | 1.76 | empty output |
| `gofmt -l cmd/ailang-worldd` | 0 | 0.04 | empty output |
| `go test -run '^TestDemoSeedApprovalPacket$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` | 0 | 39.84 | --- PASS: TestDemoSeedApprovalPacket (33.35s) |
| `go test -run '^TestDemo' ./cmd/ailang-worldd -count=1 -timeout 400s -v` | 0 | 209.87 | 21 top-level PASS, zero SKIP; exact lines below |
| `WORLD_TOOL_AILANG_BIN= go test -run '^TestDemo' ./cmd/ailang-worldd -count=1 -timeout 400s -v` | 0 | 84.84 | 18 top-level PASS; expected tool-dependent skips listed below |
| `go test -run '^TestDaemonInvocationMintsDecisionPacket$' ./host/daemon -count=1 -timeout 45s -v` | 0 | package 0.769 | --- PASS: TestDaemonInvocationMintsDecisionPacket (0.02s) |
| `go test -run '^TestRESTCommitRefusesDecisionPacketObjects$' ./host/daemon -count=1 -timeout 45s -v` | 0 | package 0.553 | --- PASS: TestRESTCommitRefusesDecisionPacketObjects (0.01s) |
| `go test -run '^TestDecisionPacketCodec$' ./host/coordinator -count=1 -timeout 45s -v` | 0 | package 0.297 | --- PASS: TestDecisionPacketCodec (0.00s) |

Exact configured-pin full Demo gate PASS lines:

```
--- PASS: TestDemoProofHTTPProducerRoundTrip (2.10s)
--- PASS: TestDemoProofSourceHashRefusal (2.04s)
--- PASS: TestDemoProofAdapterGuards (0.01s)
--- PASS: TestDemoProofWrongKeyCannotProve (1.83s)
--- PASS: TestDemoProofNoVerifiedSetRefuses (2.76s)
--- PASS: TestDemoProofUncertainCommitStops (5.50s)
--- PASS: TestDemoProofBounds (2.06s)
--- PASS: TestDemoProofKeyIORefusals (0.00s)
--- PASS: TestDemoProofKeyDisclosure (0.00s)
--- PASS: TestDemoSeedAuthenticatedCommits (2.72s)
--- PASS: TestDemoSeedPreflightRefusals (55.76s)
--- PASS: TestDemoSeedHelpPreparation (0.00s)
--- PASS: TestDemoSeedCommitFailures (4.32s)
--- PASS: TestDemoSeedReadTestCheck (8.44s)
--- PASS: TestDemoSeedRejectsToolReportRefusals (7.64s)
--- PASS: TestDemoSeedToolGuardControls (0.00s)
--- PASS: TestDemoSeedProofStage (58.24s)
--- PASS: TestDemoSeedRequiredIdentities (0.00s)
--- PASS: TestDemoSeedApprovalControls (0.00s)
--- PASS: TestDemoSeedApprovalBudgetControl (0.00s)
--- PASS: TestDemoSeedApprovalPacket (25.21s)
```

Unset-tool gate is an expected environment-control result (rc 0), not S5 integration acceptance. Its exact skips:

```
--- SKIP: TestDemoSeedReadTestCheck (0.00s)
--- SKIP: TestDemoSeedRejectsToolReportRefusals/real-refusal-arms (0.00s)
--- SKIP: TestDemoSeedProofStage (0.03s)
--- SKIP: TestDemoSeedApprovalPacket (0.05s)
```

Its pin-free approval tests both printed top-level PASS. Configuring WORLD_TOOL_AILANG_BIN to the v0.41.0 proof binary was separately rejected (expected rc 1) at `demo_test.go:1292`, `want release AILANG v0.52.1`; see [wrong-tool-release.log](evidence/wrong-tool-release.log).

## Full mutation ledger

Every accepted mutation below was actually applied, executed, restored from the original file bytes, and rerun green with a named top-level PASS. Each mutant exited rc 1 at the stated assertion, rather than skipping or selecting zero tests. Logs include full diagnostics; observed line numbers are those at execution time. The S4 guard was rerun after strengthening its committed REST before/after observation.

### roster-length-guard

Replace the len(roster) != 4 guard with if false.

Command: `go test -run '^TestDemoSeedApprovalControls$/^wrong-roster-length$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (1.16 s). Observed failing assertion:

```
demo_test.go:1181: wrong roster length must refuse: <nil>
```

Restored-and-green: rc 0 (0.69 s), `--- PASS: TestDemoSeedApprovalControls (0.00s)`.

[Mutant log](evidence/roster-length-guard-mutant.log), [restored log](evidence/roster-length-guard-restored.log).

### roster-kind-guard

Replace the roster kind guard with if false.

Command: `go test -run '^TestDemoSeedApprovalControls$/^wrong-roster-kind$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (1.13 s). Observed failing assertion:

```
demo_test.go:1189: wrong roster kind must refuse: <nil>
```

Restored-and-green: rc 0 (0.67 s), `--- PASS: TestDemoSeedApprovalControls (0.00s)`.

[Mutant log](evidence/roster-kind-guard-mutant.log), [restored log](evidence/roster-kind-guard-restored.log).

### input-refusal-before-dispatch

Replace the input-construction error guard with if false.

Command: `go test -run '^TestDemoSeedApprovalControls$/^invalid-roster-before-dispatch$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (1.12 s). Observed failing assertion:

```
demo_test.go:1196: invalid roster must refuse before dispatch: calls=1 err=S1-S4 recorded; S5 refused: demo-approval: tasks/send refused (HTTP 200); S6 incomplete; no success manifest written; inspect history before any further action
```

Restored-and-green: rc 0 (0.68 s), `--- PASS: TestDemoSeedApprovalControls (0.00s)`.

[Mutant log](evidence/input-refusal-before-dispatch-mutant.log), [restored log](evidence/input-refusal-before-dispatch-restored.log).

### missing-test-entry-guard

Replace the missing S3-entry read error guard with if false.

Command: `go test -run '^TestDemoSeedApprovalControls$/^missing-test-entry-before-dispatch$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (1.14 s). Observed failing assertion:

```
demo_test.go:1210: missing S3 entry must refuse before dispatch: calls=1 err=S1-S4 recorded; S5 refused: demo-approval: tasks/send refused (HTTP 404); S6 incomplete; no success manifest written; inspect history before any further action
```

Restored-and-green: rc 0 (0.67 s), `--- PASS: TestDemoSeedApprovalControls (0.00s)`.

[Mutant log](evidence/missing-test-entry-guard-mutant.log), [restored log](evidence/missing-test-entry-guard-restored.log).

### s5-refusal-guard

Replace the S5 A2A error guard with if false.

Command: `go test -run '^TestDemoSeedApprovalControls$/^approval-refusal-stops$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (1.92 s). Observed failing assertion:

```
demo_test.go:1217: S5 refusal must stop without retry: calls=1 err=S1-S5 recorded; S6 incomplete; no success manifest written; inspect history before any further action
```

Restored-and-green: rc 0 (0.66 s), `--- PASS: TestDemoSeedApprovalControls (0.00s)`.

[Mutant log](evidence/s5-refusal-guard-mutant.log), [restored log](evidence/s5-refusal-guard-restored.log).

### s6-incomplete-refusal

Replace the S6-incomplete refusal error with nil.

Command: `go test -run '^TestDemoSeedApprovalControls$/^s6-incomplete-refusal$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (1.14 s). Observed failing assertion:

```
demo_test.go:1223: S6 incomplete must refuse success: <nil>
```

Restored-and-green: rc 0 (0.68 s), `--- PASS: TestDemoSeedApprovalControls (0.00s)`.

[Mutant log](evidence/s6-incomplete-refusal-mutant.log), [restored log](evidence/s6-incomplete-refusal-restored.log).

### s4-refusal-guard

Replace the post-S4 producer error guard with if false.

Command: `go test -run '^TestDemoSeedProofStage$/^safe-refusal$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (5.96 s). Observed failing assertion:

```
demo_test.go:1070: seed uncertain/refused S4 must stop with one write and no retry/reseed/manifest: writes=1 envelopes=0 err="ailang-worldd demo seed: S1-S4 recorded; S5 refused: demo-approval: invocation or effect record reference absent; S6 incomplete; no success manifest written; inspect history before any further action\n"
```

Restored-and-green: rc 0 (5.41 s), `--- PASS: TestDemoSeedProofStage (4.74s)`.

[Mutant log](evidence/s4-refusal-guard-mutant.log), [restored log](evidence/s4-refusal-guard-restored.log).

### hardcoded-s4-ref

Replace the S4 envelope ref with the hash of a hard-coded fabricated-envelope string. The wrong value is read back from committed S5 input.

Command: `go test -run '^TestDemoSeedApprovalPacket$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (14.95 s). Observed failing assertion:

```
demo_test.go:1496: four roster refs must equal preceding real stages: got [{Kind:test-invocation Ref:sha256:39af1933232af04b2240d96e4796f7e6ff9ca49b6ec3e6111170bf1d2d5414af} {Kind:recorded-effect Ref:sha256:fcf52e12e8dfab2d7b6721cc0aaa05273eb4d15d3626629623b2d8d3b90b9364} {Kind:authored-claim Ref:sha256:6d438b4f7537389edbb0918bf93edf866820e443dc566864e0bfdee6c8a0ac00} {Kind:proof-envelope Ref:sha256:f170778d3dc1ce9b3a3bdfa9ec1b38ece9936ab7fc727233097c85b890136e38}] want [{Kind:test-invocation Ref:sha256:39af1933232af04b2240d96e4796f7e6ff9ca49b6ec3e6111170bf1d2d5414af} {Kind:recorded-effect Ref:sha256:fcf52e12e8dfab2d7b6721cc0aaa05273eb4d15d3626629623b2d8d3b90b9364} {Kind:authored-claim Ref:sha256:6d438b4f7537389edbb0918bf93edf866820e443dc566864e0bfdee6c8a0ac00} {Kind:proof-envelope Ref:sha256:c1cf451dfd677970c7aad1fe780a369b4441f15cb566a3c43c77f2248e53e8b4}]
```

Restored-and-green: rc 0 (14.49 s), `--- PASS: TestDemoSeedApprovalPacket (13.81s)`.

[Mutant log](evidence/hardcoded-s4-ref-mutant.log), [restored log](evidence/hardcoded-s4-ref-restored.log).

### budget-greater-equal

Widen the exact budget assertion predicate from before-after == 1 to before-after >= 1. The independent delta-two control reds.

Command: `go test -run '^TestDemoSeedApprovalBudgetControl$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (1.14 s). Observed failing assertion:

```
demo_test.go:1275: budget delta must equal exactly 1; accepted delta 2
```

Restored-and-green: rc 0 (0.68 s), `--- PASS: TestDemoSeedApprovalBudgetControl (0.00s)`.

[Mutant log](evidence/budget-greater-equal-mutant.log), [restored log](evidence/budget-greater-equal-restored.log).

### seed-double-s5

Add a second identical demo-approval A2A dispatch in the seed. The daemon rejects the duplicate but the explicit request count still reds.

Command: `go test -run '^TestDemoSeedApprovalPacket$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (14.99 s). Observed failing assertion:

```
demo_test.go:1532: seed must invoke demo-approval exactly once: calls=2
```

Restored-and-green: rc 0 (14.93 s), `--- PASS: TestDemoSeedApprovalPacket (13.78s)`.

[Mutant log](evidence/seed-double-s5-mutant.log), [restored log](evidence/seed-double-s5-restored.log).

### M4a-extra-effect

Use the valid copied two-effect temp module in the positive arm too (full TestDemoSeedApprovalPacket command). Exactly-one assertion reds with entries=0/packets=0 and the real coordinator L1 refusal in its diagnostic.

Command: `go test -run '^TestDemoSeedApprovalPacket$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (15.11 s). Observed failing assertion:

```
demo_test.go:1401: exactly one S5 carrying entry/effect/packet required: entries=0 packets=0; seed=ailang-worldd demo seed: S1-S4 recorded; S5 refused: demo-approval: tasks/send refused (HTTP 200); S6 incomplete; no success manifest written; inspect history before any further action
```

Restored-and-green: rc 0 (14.41 s), `--- PASS: TestDemoSeedApprovalPacket (13.75s)`.

[Mutant log](evidence/M4a-extra-effect-mutant.log), [restored log](evidence/M4a-extra-effect-restored.log).

### human-budget-preflight-guard

Disable the existing budget admission guard with if false (retain _ = snap so the mutant compiles). The daemon’s later refusal cannot mask the extra S5 call and ten committed history entries.

Command: `go test -run '^TestDemoSeedApprovalPacket$/^no-human-budget$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (6.29 s). Observed failing assertion:

```
demo_test.go:1389: missing Human.Approve budget must refuse before any packet or dispatch and preserve REST history: packets=0 calls=1 entries=10 head=sha256:3549988759c6d0686fcd5bc72b86034a3e1e9d46476f332d9deab822ab003833 before= read=<nil> err=ailang-worldd demo seed: S1-S5 recorded; S6 incomplete; no success manifest written; inspect history before any further action
```

Restored-and-green: rc 0 (5.15 s), `--- PASS: TestDemoSeedApprovalPacket (4.00s)`.

[Mutant log](evidence/human-budget-preflight-guard-mutant.log), [restored log](evidence/human-budget-preflight-guard-restored.log).

### s4-refusal-guard-with-rest-snapshot

Repeat the post-S4 if-false mutation after adding the REST snapshot captured immediately before the envelope write. The refusal assertion also requires committed history to match that snapshot.

Command: `go test -run '^TestDemoSeedProofStage$/^safe-refusal$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

Mutant rc 1 (6.03 s). Observed failing assertion:

```
demo_test.go:1076: seed uncertain/refused S4 must stop with one write and no retry/reseed/manifest: writes=1 envelopes=0 err="ailang-worldd demo seed: S1-S4 recorded; S5 refused: demo-approval: invocation or effect record reference absent; S6 incomplete; no success manifest written; inspect history before any further action\n"
```

Restored-and-green: rc 0 (5.43 s), `--- PASS: TestDemoSeedProofStage (4.76s)`.

[Mutant log](evidence/s4-refusal-guard-with-rest-snapshot-mutant.log), [restored log](evidence/s4-refusal-guard-with-rest-snapshot-restored.log).

## Findings and unfinished work

- No M4a implementation, mandatory gate or accepted mutation control remains unfinished. S6 and manifest serialization are deliberately not implemented, as required by the contract.
- The first duplicate-invocation experiment red at an earlier packet-carrier assertion because the test proxy overwrote its pre-S5 snapshot on the second call. The proxy now retains its first snapshot; the rerun reds at the explicit calls==1 assertion. The initial result is retained in [seed-double-s5-initial-masked.log](evidence/seed-double-s5-initial-masked.log) and is not counted as the call-count discharge.
- The first experiment disabling the inherited budget guard did not compile because snap became unused. That invalid mutant is not counted as evidence. The compile-preserving if-false rerun reached the named budget refusal/history assertion; the initial build result is retained in [human-budget-preflight-initial-invalid.log](evidence/human-budget-preflight-initial-invalid.log).
- The existing demo.go help still describes an S1-only incremental build. It predates this change and is outside the allowlist; it was not edited. Runtime refusal text now accurately identifies S1-S5/S6.
- No unrelated gate failure was observed. The configured wrong-tool failure and accepted mutation failures were intentional negative controls.
- A final read-only process-status probe (`ps -eo pid,etime,command`) returned `operation not permitted: ps`: UNINFORMATIVE UNDER SANDBOX. It was not a gate, is neither a pass nor a fail, and does not affect the successfully completed test gates.

Raw gate logs and machine-readable ledgers are in [evidence/](evidence/). No session credential or proof-key file was copied into the snapshot.

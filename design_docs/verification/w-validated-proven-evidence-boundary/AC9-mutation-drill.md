# AC9 full mutation re-drill

Environment: repo root `/Users/voightkampff/dev/sunholo-data/.wt-iter108-pef`; Go commands used `PATH=/opt/homebrew/bin:$PATH`, `GOTOOLCHAIN=go1.25.6`, and Go tests used `AILANG_BIN=/tmp/ailang-v0300/ailang`.

## M1 — arbitrary AILANG receipt authority

- Exact edit: `world/types.ail`, in `gradeOf`'s body only, `ProofReceipt(_) => CLAIMED` → `ProofReceipt(_) => PROVEN` (the identical ensures-clause arm was untouched).
- LANDED: before `323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f`; mutant `52686198384cbbaa51bf8f013a08e83881ce9f28bb1522498b3f984f7e86adac` (different).
- BUILDS: `GOTOOLCHAIN=go1.25.6 go build ./...` returned rc=0 before the test verdict was read.
- Exact mutant gate command: `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` (repo-root cwd), rc=1.
- Observed RED: `✗ world/types.ail: verify.counterexample == 1`.
- Enumerated RED SET: direct repo-root enumeration with `AILANG_BIN=/tmp/ailang-v0300/ailang /tmp/ailang-v0300/ailang test --format json world/` returned `failed_tests=1`: `gradeCode_test_7`, assertion `test 0: expected 1, got 4`. In addition, the gate's contract-verification arm reported the `gradeOf` counterexample above. Thus the named inline killer is one of two independent gate observations, and the inline named-test red set itself has size 1. A `-skip` check is not applicable to the overall AILANG gate because the contract arm also kills.
- Named control GREEN after restore, same gate command: rc=0, `10/10 required world/ identities verified`, `all 40 required named tests pass (failed_tests=0)`, projection `4/4`, and world package `9/9` passed.
- Restore verified byte-identical: backup `323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f`; restored file `323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f`.
- §6 divergence: wording differs (`expected 1, got 4`, not `got 4, want 1`), and `verify_ail.sh` stops at the contract counterexample before its named-test leg; the separately enumerated inline test nevertheless kills exactly as predicted.

## M17 — named-manifest removal (round A snapshot; not re-run)

- Source: `.snap/S2/`, as directed. Round A removed one literal required evidence-test name from the isolated copy of `scripts/verify_go.sh` while leaving the observed synthetic test present.
- LANDED/BUILDS/restore: already drilled and independently verified by the controller in round A; not mutated or re-run in this round, per directive. The S2 snapshot retains the landed manifest gate and its isolated test.
- Exact killer command represented by the isolated test: `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6 go test ./host/verifygate -run '^TestEvidenceNamedManifestRejectsUnpinnedTest$' -count=1`.
- Observed RED from the round-A drill: `evidence test set differs from REQUIRED_EVIDENCE_TESTS`, with the still-observed removed literal reported as extra. Enumerated RED SET: `TestEvidenceNamedManifestRejectsUnpinnedTest` only (isolated gate self-mutation test); sole killer: yes. Its pristine synthetic control requires `all 37 required top-level evidence tests passed exactly once`.
- Control GREEN/restored: round-A controller verified both; this row is a restatement, not a new empirical run.

## M2 — invalid proof ref

- Exact edit: `host/evidence/validator.go`, `if _, err := hashref.Parse(reportRef.String()); err != nil` → `if _, err := hashref.Parse(reportRef.String()); false && err != nil`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `b936c2fe5ca4bf969a2e0ad38255d8dba0a7013d4a2a76526d55a564802b7247` (different).
- BUILDS: `GOTOOLCHAIN=go1.25.6 go build ./...` rc=0 before reading test results.
- Exact mutant command: `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6 go test -json ./host/evidence ./host/store ./host/verifygate -count=1`, rc=1.
- Observed RED: `authority_test.go:114: invalid-ref guard: got missing; want invalid_ref`.
- Enumerated RED SET: `{host/evidence.TestInvalidProofRefIsRefused}`; size 1; named killer is sole killer. Mutant check `go test ./host/evidence ./host/store ./host/verifygate -count=1 -skip '^TestInvalidProofRefIsRefused$'` returned rc=0.
- Named pristine control: `go test ./host/evidence ./host/store ./host/verifygate -count=1 -run '^TestInvalidProofRefIsRefused$'` returned rc=0.
- Restore verified byte-identical: backup and restored hashes both `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: only the source spelling drifted (`hashref.Parse` inline rather than the table's `refErr` local); behavior and failure text matched.

## M3 — missing proof report

- Exact edit: `host/evidence/validator.go`, `if payload == nil` → `if false && payload == nil` (current implementation spelling of §6's `if !ok`).
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `8596ff9b796089cc33af0943dc459ac1fad48028fa5be6949250de80bf4c4e70` (different).
- BUILDS: `GOTOOLCHAIN=go1.25.6 go build ./...` rc=0 before reading results.
- Exact mutant command: `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6 go test -json ./host/evidence ./host/store ./host/verifygate -count=1`, rc=1.
- Observed RED: `authority_test.go:121: missing-object guard: got hash_mismatch; want missing`.
- Enumerated RED SET: `{host/evidence.TestMissingProofReportIsRefused}`; size 1; sole killer. The same module-wide command with `-skip '^TestMissingProofReportIsRefused$'` returned rc=0.
- Named pristine control: same module-wide package command with `-run '^TestMissingProofReportIsRefused$'` returned rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: table predicts `got malformed; want missing`; observed successor is `hash_mismatch`. The source anchor also drifted from `!ok` to `payload == nil`.

## M5 — payload hash integrity

- Exact edit: `host/evidence/validator.go`, `if got := hashref.SumSHA256(payload); got != reportRef` → `if got := hashref.SumSHA256(payload); false && got != reportRef`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `75afe7d7c6bb81609bcb65d7c668454dd769213fbcdea2fdb588b9549000841d`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant enumeration command: pinned/AILANG-set `go test -json ./host/evidence ./host/store ./host/verifygate -count=1`, rc=1.
- Observed RED: `recomputed payload-hash guard: got unauthenticated_report; want hash_mismatch`.
- Enumerated RED SET: `{host/evidence.TestPayloadHashMismatchIsRefused}`; size 1; sole killer. Mutant `go test ./host/evidence -count=1 -skip '^TestPayloadHashMismatchIsRefused$'` rc=0.
- Named pristine control: `go test ./host/evidence -count=1 -run '^TestPayloadHashMismatchIsRefused$'` rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: successor reason is `unauthenticated_report` as predicted; only assertion punctuation/prefix differs.

## M7 — wrong semantic ID

- Exact edit: `host/evidence/validator.go`, `if meta.SemanticID != ProofSemanticID` → `if false && meta.SemanticID != ProofSemanticID`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `ac678f68cbf5866f754c91b09a7504b5f1f8c9a2b63153ba2338bf212d134184`.
- BUILDS: pinned `go build ./...` completed rc=0 before enumeration (empty build log and successful command stage).
- Exact mutant enumeration command: pinned/AILANG-set `go test -json ./host/evidence ./host/store ./host/verifygate -count=1`; JSON completed all three packages and enumerated the red below. The enclosing tool session ended immediately after the 27.026s verifygate package completion, before its trailing shell rc print; the JSON contains the package completions and failing-test event.
- Observed RED: `semantic-ID guard: got unauthenticated_report; want wrong_semantic_id`.
- Enumerated RED SET: `{host/evidence.TestWrongSemanticIDIsRefused}`; size 1; sole killer. Mutant `go test ./host/evidence -count=1 -skip '^TestWrongSemanticIDIsRefused$'` rc=0.
- Named pristine control: `go test ./host/evidence -count=1 -run '^TestWrongSemanticIDIsRefused$'` rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: behavior matched; only assertion prefix/punctuation differs.

## Clean stop / resume point

- Completed 6/27 live rows: M1, M2, M3, M5, M7, and the directed S2 restatement of M17.
- Remaining, in prescribed order: M18, M19; M8, M9, M10, M11, M12, M13, M14, M16, M20, M21, M23, M24, M26, M27, M29, M30; then real-store M4, M22, M25.
- No nine-row boundary was reached, so no cumulative `.snap/S3/` checkpoint was due.
- Tool-runtime note: a combined three-package test takes about 27 seconds because `host/verifygate` itself takes ~27 seconds. Some enclosing shell invocations were ended immediately after that package completion before trailing shell output; mutations were restored and hash-verified immediately in the next command. M7 records this explicitly.

## Round C continuation

## M18 — projection drift

- Exact edit: `world/types.ail`, inserted the source-only comment `// AC9 M18/M19 projection freshness mutation.` immediately after `module world/types`; the package projection was deliberately left unchanged.
- LANDED: before `323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f`; mutant `a3525190baee63f7a44e7a7c1a6367d1d0ddb857a04baa91307eb81f9dae5afe` (different).
- BUILDS: `GOTOOLCHAIN=go1.25.6 go build ./...` rc=0 before the gate verdict was read.
- Exact mutant command: `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` from the repo root, rc=1.
- Observed RED: `✗ projection mismatch: world/types.ail=a3525190baee63f7a44e7a7c1a6367d1d0ddb857a04baa91307eb81f9dae5afe packages/world-core/world/types.ail=323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f`.
- Enumerated RED SET: the full gate ran through AILANG legs 1 and 2 green, then stopped at world-package step 3/9 on the single projection hash comparison above; size 1; the named step-3 arm is the sole killer.
- Named pristine control: immediately before this AILANG batch, the same repo-root gate command returned rc=0 with `10/10` identities, `40` named tests, `4/4 projection hashes`, and world-package `9/9` green.
- Restore verified byte-identical: deferred only across the immediately-following M19 half of this same AILANG batch; the canonical-file restore hash is recorded in M19, after which the mandatory post-batch pristine gate is run.
- §6 divergence: exact wording is `projection mismatch: ...`, not the table's `projection hash mismatch: world/types.ail`; the live diagnostic also includes both full hashes and both paths.

## M19 — stale ready packet

- Exact edit: after M18's canonical comment insertion, `packages/world-core/world/types.ail` was rebuilt to the same bytes by inserting the identical comment after `module world/types`; `scripts/world_package_ready_packet.golden.json` remained at sha256 `fc9d23b2c9b66149c5bd8af9cc25e6c5c29cd14227a5bb472801aeb6c2419da6`.
- LANDED: projected file before `323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f`; mutant `a3525190baee63f7a44e7a7c1a6367d1d0ddb857a04baa91307eb81f9dae5afe` (different), equal to the mutated canonical source.
- BUILDS: `GOTOOLCHAIN=go1.25.6 go build ./...` rc=0 before the gate verdict was read.
- Exact mutant command: `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` from the repo root, rc=1.
- Observed RED: `✗ ready packet differs byte-for-byte from golden`; the diff changed `contentHash` `7616498e…` → `c67a9735…`, `tarballBytes` `7883` → `7916`, and `tarballSHA256` `2f18c5e8…` → `64a0bd06…`, while `interfaceHash` remained unchanged.
- Enumerated RED SET: the full gate ran AILANG legs 1 and 2 green and world-package steps 1–8 green, then stopped on the single byte comparison at step 9/9; size 1; the named step-9 arm is the sole killer.
- Named pristine control: immediately after restore, the same repo-root gate returned rc=0 with `10/10` identities, `40` named tests, `4/4 projection hashes`, and world-package `9/9` green.
- Restore verified byte-identical: canonical backup/restored both `323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f`; projection backup/restored both `323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f`. This also completes M18's deferred restore verification.
- §6 divergence: none; the failure text matched exactly and the diff had the documented three changed packet fields with stable `interfaceHash`.

## M8 — wrong interface

- Exact edit: `host/evidence/validator.go`, `if meta.InterfaceHash != proofInterfaceHash` → `if false && meta.InterfaceHash != proofInterfaceHash`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `db86b583116c41f53864842385182c61bf147a81f885a6823c2fae511eb7a56c` (different).
- BUILDS: `GOTOOLCHAIN=go1.25.6 go build ./...` rc=0 before reading test results.
- Exact mutant command: `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6 go test -json ./host/evidence ./host/store ./host/verifygate -count=1`, rc=1.
- Observed RED: `authority_test.go:144: interface-hash guard: got unauthenticated_report; want wrong_interface`.
- Enumerated RED SET: `{host/evidence.TestWrongInterfaceIsRefused}`; size 1; named killer is sole killer. This was the round-C narrowing control, enumerated across all three packages: `host/store` passed and `host/verifygate` explicitly passed (`ok .../host/verifygate 27.729s`).
- Named pristine control: after restore, `go test ./host/evidence -run '^TestWrongInterfaceIsRefused$' -count=1` returned rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- Round-C narrowing: control validated. Subsequent Go-mutant red sets are enumerated over host/evidence + host/store; host/verifygate excluded — see this round-C control.
- §6 divergence: behavior matched; only the assertion prefix/punctuation differs.

## M9 — malformed report

- Exact edit: live anchor `host/evidence/envelope_codec.go`, `if err != nil` immediately after `decoded, err := DecodeProofReportV1(report)` → `if false && err != nil`, allowing the zero decoded report to continue.
- LANDED: before `6710e18d0f1b7a28b9ec955d58d754f6d4c665d41c991a83a2b20422b1f36dee`; mutant `a87fda80932fcc2b0bd27999819cd10519bc5f90ae8208cd7a01f2a13156afc6`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `authority_test.go:154: strict report-decode guard: got unauthenticated_report; want malformed`.
- Enumerated RED SET: `{host/evidence.TestMalformedProofReportIsRefused}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: `go test ./host/evidence -run '^TestMalformedProofReportIsRefused$' -count=1` rc=0.
- Restore verified byte-identical: backup/restored `6710e18d0f1b7a28b9ec955d58d754f6d4c665d41c991a83a2b20422b1f36dee`.
- §6 divergence: material source-anchor drift: the load-bearing strict report-decode guard is now in `host/evidence/envelope_codec.go`, not the table's `report_codec.go`; observed failure text otherwise matched.

## M10 — MAC authentication

- Exact edit: `host/evidence/validator.go`, `if len(envelope.MAC) != sha256.Size || !hmac.Equal(want, envelope.MAC)` → `if false && (len(envelope.MAC) != sha256.Size || !hmac.Equal(want, envelope.MAC))`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `63d433837cbfce35604b393f036028639bff2caa3047d5b513b3bb050ee1bd84`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `absent-MAC authentication guard: got ; want unauthenticated_report` and `wrong-MAC authentication guard: got ; want unauthenticated_report`.
- Enumerated RED SET: `{host/evidence.TestOtherwisePerfectReportWithoutMACIsUnauthenticated, host/evidence.TestOtherwisePerfectReportWithWrongMACIsUnauthenticated}`; size 2; both named killers, each explained by its absent-tag/wrong-tag stimulus reaching a sealed result. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine controls: one regex-selected run of both named tests returned rc=0 after restore.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: the tests do not emit the table's `report resolved PROVEN; want unauthenticated_report`; their shared `requireReason` observes the sealed result as an empty unsupported reason and prints `got ; want unauthenticated_report`.

## M11 — subject mismatch

- Exact edit: `host/evidence/validator.go`, `if report.Subject != expectedSubject` → `if false && report.Subject != expectedSubject`.
- LANDED: before `d9de07e…`; mutant `c7330e3b303469c8f4790393b504c5b44196a78f7d3b4b57c1aec869c5269dfb`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `authority_test.go:176: subject-binding guard: got tool_mismatch; want subject_mismatch`.
- Enumerated RED SET: `{host/evidence.TestMismatchedProofSubjectIsRefused}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named-only test rc=0 after restore.
- Restore verified byte-identical: backup/restored full hash `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: none beyond assertion prefix/punctuation.

## M14 — proof incomplete

- Exact edit: `host/evidence/validator.go`, `if incomplete` → `if false && incomplete`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `6ecf5aacc06a4e9ff22f52c54105f9f9f4d76088dd3c8ade1604da44506d166d`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `authority_test.go:204: required-identity guard: got ; want proof_incomplete`.
- Enumerated RED SET: `{host/evidence.TestIncompleteProofReportIsRefused}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named-only test rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: the test emits `got ; want proof_incomplete`, not the table's `incomplete proof resolved PROVEN; want proof_incomplete`; it observes the sealed result through the unsupported-reason helper.

## M16 — seal bypass

- Exact edit: live file `host/evidence/validator.go` (the table's `grade.go` does not exist), added exported `func ResolveHashRef(hashref.HashRef) ResolvedGrade { return ResolvedGradeProven }`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `d500947326e6337c1ae69b6c6c1306c2023f5eb9722e3bfdd57588eacf39d695`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `public authority surface changed: added: [func ResolveHashRef] removed: []`.
- Enumerated RED SET: `{host/evidence.TestPublicAuthoritySurfaceIsFrozen}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named-only test rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: material file drift (`grade.go` is absent; authority types now live in `validator.go`) and observed first failure is the exact exported-surface inventory delta, not `public authority surface exposes non-sealed PROVEN ingress`, because inventory comparison precedes the return-type scan.

## M20 — cross-validator binding

- Exact edit: `host/evidence/validator.go`, `if sealed.mintedBy != v.id` → `if false && sealed.mintedBy != v.id`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `c17bc0ad59a89f7148310f8ab64d2a80da6026562bc946303c78e5ad531792d4`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `binding check: got <nil>; want ErrForeignSeal`; sibling pin: `distinct NewValidator calls with the same key shared a mint identity`.
- Enumerated RED SET: `{host/evidence.TestAttackerChosenValidatorCannotMintForHostAuthority, host/evidence.TestValidatorMintIdentitiesAreDistinct}`; size 2; named killer is one of 2. The first observes foreign resolution directly; the second independently requires distinct validator instances not to cross-resolve. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named killer rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: wider red set than the table names; named assertion reports the missing `ErrForeignSeal` (`got <nil>`) rather than `foreign seal resolved ResolvedGradeProven`.

## M21 — zero-value mint validity

- Exact edit: `host/evidence/validator.go`, current `if v == nil || v.id == nil || sealed.mintedBy == nil` → `if false && (v == nil || v.id == nil || sealed.mintedBy == nil)`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `18beadb090d1a5b0935442d86fa4e2df8ac218e9814fb00dd38e7e12fcd16cfa`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `mint-validity check: got <nil>; want ErrUnmintedAuthority`.
- Enumerated RED SET: `{host/evidence.TestZeroValueForgeryCannotResolve}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named-only test rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: current guard additionally protects a nil `*Validator`; named failure reports missing error (`got <nil>`) rather than the table's resolved-grade wording.

## M23 — deadline arming

- Exact edit: `host/evidence/validator.go`, `context.WithTimeout(ctx, v.compiler.ObjectReadTimeout)` → `context.WithCancel(ctx)`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `b33f23b1d8fdc066055f20721276a9b2e7b27ae1e95563f3f46ffb2e31bd6af2`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `blocked read exceeded the test-side watchdog of 54.982958ms (derived from the measured decoy hold)` after logging derived timeout `2.749147ms`.
- Enumerated RED SET: `{host/evidence.TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: same real-store named test rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: this run hit the documented watchdog hang mode rather than returning later and printing the table's `blocked read sealed...` branch.

## M24 — oversize translation

- Exact edit: `host/evidence/validator.go`, `if errors.As(err, &tooLarge)` → `if false && errors.As(err, &tooLarge)`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `e81cc8a322f999cef0a296e3317adc280ec3f0e2531389a6c2c5a36c27434066`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `got operational store-read error; want oversize: evidence: read proof report: store: object payload is 349339 bytes; maximum is 262144`.
- Enumerated RED SET: `{host/evidence.TestOversizeProofReportIsRefused}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named-only test rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: behavior matched; live assertion includes the wrapped operational error detail.

## M27 — depth-bomb refusal

- Exact edit: `host/evidence/proposal_codec.go`, the `if err != nil` immediately after `d.Decode(&v)` → `if false && err != nil`.
- LANDED: before `ed757fd9c7c25f039c7d125019dbbd37af022089abbbfac04750b52a21e71abb`; mutant `44a06aaccc6332c8f0449b0a116287d6a23fddb0647f32428da831566267232b`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `depth-bomb refusal = evidence: decode refused: malformed: trailing JSON; want the stdlib scanner's own max-depth refusal`.
- Enumerated RED SET: `{host/evidence.TestNestingDepthBombWithinByteCapIsRefused}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named-only test rc=0.
- Restore verified byte-identical: backup/restored `ed757fd9c7c25f039c7d125019dbbd37af022089abbbfac04750b52a21e71abb`.
- §6 divergence: material observable drift. Swallowing the first decode error does not return success: the following `d.Token()` produces a typed `trailing JSON` refusal, and the test distinguishes it from the required stdlib max-depth refusal. The table predicts `depth bomb decoded as ClaimedEvidence`.

## M29 — required validator identities

- Exact edit: `host/evidence/validator.go`, `if len(requiredIdentities) == 0` → `if false && len(requiredIdentities) == 0`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `82cbda4dc8f28114312f4d8c5b324a7e15de38cce9d776cb68651a9d164e1d21`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `NewValidator accepted empty required identities; want ErrInvalidValidatorConfig`.
- Enumerated RED SET: `{host/evidence.TestConstructorRefusesEmptyRequiredIdentities}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named-only test rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: none.

## M26 — timeout ordering

- Exact edit: `host/evidence/validator.go`, `if window > 0 && cfg.ObjectReadTimeout <= window` → `if false && (window > 0 && cfg.ObjectReadTimeout <= window)`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `a9e5cc53c48a76f39f1898c44f11bc471fb28497d8d68b756639125b9ca014fe`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `ordering refusal did not name runtime values: <nil>`; real-store arm `positive unordered timeout: got <nil>; want ErrUnorderedTimeouts and not ErrInvalidValidatorConfig`; wrapper arm `forwarding wrapper lost real-store wait bound: got <nil>; want ErrUnorderedTimeouts`.
- Enumerated RED SET: `{TestConstructorNamesActuallyUsedUnorderedTimeouts, TestConstructorPinsBusyTimeoutBelowObjectReadTimeout, TestReaderWaitBoundsCannotBeLostThroughWrapper/forwarding-real-store}`; top-level size 3. The first two are named killers; the third is the integration consumer of the same ordering refusal. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine controls: regex-selected isolated and real-store named tests rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: wider red set (the AC22 forwarding-real-store subtest also pins this branch); live AC18 wording differs from the table's quoted sentence.

## M30 — mandatory wait-bound metadata

- Exact edit: `host/evidence/validator.go`, replaced the `window < 0` branch's `return nil, fmt.Errorf(...)` with `window = 0`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `6eb0cf23248b7c81bdf101937aa26686c1c20cb11c060296bb6a73de3bd006d7`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `NewValidator accepted unknown BusyTimeout: <nil>` and `NewValidator accepted wrapper with unknown wait bound: <nil>; want ErrInvalidValidatorConfig`.
- Enumerated RED SET: `{TestConstructorRefusesUnknownBusyTimeout, TestReaderWaitBoundsCannotBeLostThroughWrapper/unknown}`; top-level size 2; both named by §6 (isolated killer plus integration arm). Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: isolated named test rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: none; integration arm supplies the second expected red.

## M4 — oversize pre-materialization guard

- Exact edit: live file `host/store/read_object.go`, `if meta.PayloadLength > maxBytes` → `if false && meta.PayloadLength > maxBytes`.
- LANDED: before `a25f41ec604872622cd694e1f3e3941fc6a64a225bce9bf72ba1548fd9c3128c`; mutant `08c1b4326b1006f12149e2d32174e8250ea5a55a2b486f58d8a4e23082c82f56`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: named integration arm `oversize envelope sealed; want oversize`; direct store arm `ReadObject oversize error = <nil>; want *ObjectTooLargeError`.
- Enumerated RED SET: `{host/evidence.TestOversizeProofReportIsRefused, host/store.TestReadObjectProbeOmitsPayloadAndGuardsMaterialization}`; size 2; named killer is one of 2. The evidence arm proves the end-to-end authority consequence; the store arm directly pins typed refusal/pre-materialization. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: real-store evidence named test rc=0.
- Restore verified byte-identical: backup/restored `a25f41ec604872622cd694e1f3e3941fc6a64a225bce9bf72ba1548fd9c3128c`.
- §6 divergence: source moved from table's `host/store/store.go` to `host/store/read_object.go`; red set also contains the direct store test not named in §6.

## M22 — read detached from context

- Exact edit: live `host/store/read_object.go`, replaced supplied `ctx` with `context.Background()` at all four context-taking calls: `s.db.Conn`, `conn.BeginTx`, probe `tx.QueryRowContext`, and payload `tx.QueryRowContext`.
- LANDED: before `a25f41ec604872622cd694e1f3e3941fc6a64a225bce9bf72ba1548fd9c3128c`; mutant `4f14f0296fd6454a8def1bc1f5cf21a00aa94d0181ad38e5567b1723ee457dd0`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `blocked read exceeded the test-side watchdog of 61.79975ms (derived from the measured decoy hold)` after derived timeout `3.089987ms`.
- Enumerated RED SET: `{host/evidence.TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: same real-store named test rc=0.
- Restore verified byte-identical: backup/restored `a25f41ec604872622cd694e1f3e3941fc6a64a225bce9bf72ba1548fd9c3128c`.
- §6 divergence: source moved from `host/store/store.go` to `read_object.go`; this run hit the explicitly documented watchdog hang mode rather than the table's later sealed-result assertion.

## M25 — one-snapshot probe/payload

- Exact edit: live `host/store/read_object.go`, removed `conn.BeginTx`, deferred rollback, and commit; changed both `tx.QueryRowContext` calls to `conn.QueryRowContext`, yielding two bare autocommit statements on the same reserved connection.
- LANDED: before `a25f41ec604872622cd694e1f3e3941fc6a64a225bce9bf72ba1548fd9c3128c`; mutant `d4090084388e143eabb2d5e1b7d3e61d34fa59327a6a01b29749f6149dae5cf9`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `read_object_test.go:141: payload diverged from probed row: probe=100 bytes, payload=350 bytes; want one snapshot`.
- Enumerated RED SET: `{host/store.TestConcurrentMutationCannotDesyncProbeAndPayload/concurrent_mutation}`; top-level size 1; sole killer. Its no-write control subtest remained green in the same enumeration. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: top-level named test rc=0 after restore.
- Restore verified byte-identical: backup/restored `a25f41ec604872622cd694e1f3e3941fc6a64a225bce9bf72ba1548fd9c3128c`.
- §6 divergence: source moved from `host/store/store.go` to `read_object.go`; actual injected new payload is 350 bytes, not the table's illustrative 300 bytes.

## M13 — failed proof

- Exact edit: `host/evidence/validator.go`, `if proofFailed` → `if false && proofFailed`.
- LANDED: before `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `19a37f4e8911fb7dcb96c016db60c1d6cdf5d4c40d91c809f6aa8f6e3c460f5f`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `authority_test.go:195: proof-success guard: got proof_incomplete; want proof_failed`.
- Enumerated RED SET: `{host/evidence.TestFailedProofReportIsRefused}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named-only test rc=0.
- Restore verified byte-identical: backup/restored `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: none beyond assertion prefix/punctuation.

## M12 — tool mismatch

- Exact edit: `host/evidence/validator.go`, `if toolMismatch` → `if false && toolMismatch`.
- LANDED: before full hash `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`; mutant `6ca671868f6b1db5c1e58d77c865bcad02a274a948eb09e9f7ec7d165b70dd1c`.
- BUILDS: pinned `go build ./...` rc=0 before verdict.
- Exact mutant command: pinned/AILANG-set `go test -json ./host/evidence ./host/store -count=1`, rc=1.
- Observed RED: `compiler-identity guard: got proof_incomplete; want tool_mismatch`.
- Enumerated RED SET: `{host/evidence.TestMismatchedProofToolIsRefused}`; size 1; sole killer. Red set enumerated over host/evidence + host/store; host/verifygate excluded — see the round-C control.
- Named pristine control: named-only test rc=0.
- Restore verified byte-identical: backup/restored full hash `d9de07eb6087605c5d54d256b0022c04865db8e4a1f307410dce2ec29b12b4d7`.
- §6 divergence: none beyond assertion prefix/punctuation.

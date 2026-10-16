AILANG World row 150 M3b — iteration 273, fix round r2.

Completed against HEAD 72bd59a after reading the independent r1 evaluation. All three requested guards are strengthened, the unused key-file write is removed, all six mutation drills fail at new named assertions and pass after restoration, and all four gates pass. Nothing remains unfinished.

Changes:

- `cmd/ailang-worldd/demo_seed.go`: names the exact `["sameRef","servesEntry"]` configuration and uses it in the S4 call. A single forwarding function variable (`demoSeedProve = demoProve`) lets the integration test observe the actual adapters and arguments. This seam is necessary because HTTP cannot reveal the client’s timeout fields; production still calls the same producer with the same inputs and behavior.
- `cmd/ailang-worldd/demo_test.go`: adds `TestDemoSeedRequiredIdentities` and the seed integration subcontrols `required-identities` and `http-bounds`. Both configured and actual call requirements must match exactly. Both real reader/writer adapters must have 10 s context and HTTP client bounds; the transport also checks actual request deadlines and counts real object reads and the envelope write. The seed envelope is decoded to assert both verified identities, and validated using a validator requiring both identities. Checks stdout, stderr, manifest path and any unexpected manifest contents for raw/encoded key material; still requires the manifest to be absent.
- `cmd/ailang-worldd/demo_proof_test.go`: adds `TestDemoProofKeyDisclosure` to check returned error text for short/long keys. Shared disclosure checks cover raw bytes, lowercase hex, uppercase hex, standard base64, and raw URL base64 for the full key and malformed file bytes. Uncertain-commit error checks use the same helper. Removes only the unused `freshProofKey` TempDir/write; key writes actually read by tests are retained.

Only the three files above are changed and copied under their relative paths in this directory. `demo_proof.go` was mutated during drills and restored byte-for-byte; it has no final diff. `changed-files.json` records snapshot SHA-256 hashes, checked against the final workspace.

Gates ran with `AILANG_BIN=$HOME/.pinned-ailang/ailang` and `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.52.1/ailang`. Exact commands, return codes, logs and named top-level PASS lines follow. Both test gates ran nonzero tests with zero SKIP lines.

| Gate | rc | Result | Log |
| --- | ---: | --- | --- |
| `go vet ./cmd/ailang-worldd` | 0 | PASS; no output | [gate-vet.log](gate-vet.log) |
| `gofmt -l cmd/ailang-worldd` | 0 | PASS; no output | [gate-gofmt.log](gate-gofmt.log) |
| `go test -run '^TestDemo' ./cmd/ailang-worldd -count=1 -timeout 400s -v` | 0 | PASS; 18 named top-level tests; zero SKIP | [gate-demo.log](gate-demo.log) |
| `go test -run '^TestProducerEndToEndMintsEnvelopeValidatorAccepts$' ./host/evidence -count=1 -timeout 120s -v` | 0 | PASS; 1 named top-level tests; zero SKIP | [gate-evidence.log](gate-evidence.log) |

Named top-level PASS lines from the final gates:

```text
--- PASS: TestDemoProofHTTPProducerRoundTrip (1.11s)
--- PASS: TestDemoProofSourceHashRefusal (0.97s)
--- PASS: TestDemoProofAdapterGuards (0.00s)
--- PASS: TestDemoProofWrongKeyCannotProve (1.15s)
--- PASS: TestDemoProofNoVerifiedSetRefuses (1.24s)
--- PASS: TestDemoProofUncertainCommitStops (3.14s)
--- PASS: TestDemoProofBounds (0.96s)
--- PASS: TestDemoProofKeyIORefusals (0.00s)
--- PASS: TestDemoProofKeyDisclosure (0.00s)
--- PASS: TestDemoSeedAuthenticatedCommits (1.06s)
--- PASS: TestDemoSeedPreflightRefusals (32.46s)
--- PASS: TestDemoSeedHelpPreparation (0.00s)
--- PASS: TestDemoSeedCommitFailures (2.69s)
--- PASS: TestDemoSeedReadTestCheck (4.21s)
--- PASS: TestDemoSeedRejectsToolReportRefusals (4.15s)
--- PASS: TestDemoSeedToolGuardControls (0.00s)
--- PASS: TestDemoSeedProofStage (30.72s)
--- PASS: TestDemoSeedRequiredIdentities (0.00s)
--- PASS: TestProducerEndToEndMintsEnvelopeValidatorAccepts (0.57s)
```

Mutation evidence is in [mutation-ledger.json](mutation-ledger.json), including exact mutations, commands, new assertion failure lines, restoration hashes and per-drill logs. Every mutant returned rc 1; every restored run returned rc 0 with named top-level PASS lines and no SKIP.

| Mutant | New assertion | Red rc | Restored rc |
| --- | --- | ---: | ---: |
| [required-identities-call-drop-servesEntry](required-identities-call-drop-servesEntry-red.log) | `seed-required-identities:` | 1 | 0 |
| [required-identities-config-drop-servesEntry](required-identities-config-drop-servesEntry-red.log) | `seed-required-identities:` | 1 | 0 |
| [seed-client-without-timeouts](seed-client-without-timeouts-red.log) | `seed-s4-http-bounds:` | 1 | 0 |
| [seed-client-without-http-client-timeout](seed-client-without-http-client-timeout-red.log) | `seed-s4-http-bounds:` | 1 | 0 |
| [seed-client-without-request-context-timeout](seed-client-without-request-context-timeout-red.log) | `seed-s4-http-bounds:` | 1 | 0 |
| [key-hex-leaked-in-error](key-hex-leaked-in-error-red.log) | `seed-key-disclosure:` | 1 | 0 |

The hex-error mutant failed at `demo_proof_test.go:949` on `returned-error` and `demo_test.go:1011` on CLI `stderr`, both with `seed-key-disclosure: ... contains hex-lower key`. Dropping servesEntry at the actual seed call failed at `demo_test.go:966`, independently of the named configuration control. Timeout mutants removed both bounds and each bound separately; all failed at `demo_test.go:977`.

Hygiene was verified by inspecting the removed unused write and the retained, read key-file writes, plus the final TestDemo gate. No hygiene mutant was specified.

`git diff --check` passed. No git write operations were run. No full suite, verify_go.sh or verify_ail.sh was run. There were no sandbox denials. The pre-existing untracked evaluator directory was left untouched.

Completed at 2026-10-16T07:01:27.478379+00:00.

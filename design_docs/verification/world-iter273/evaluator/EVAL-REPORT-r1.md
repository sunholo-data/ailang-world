# Independent evaluation, row 150 M3b (HEAD 72bd59a, base 3e22bd7)

Score: 88/100 (pass >= 70). Verdict: PASS. BLOCKING findings: none.

## Gates (run myself, pinned AILANG_BIN / WORLD_TOOL_AILANG_BIN)
- go vet ./cmd/ailang-worldd: rc 0; gofmt -l: empty.
- go test -run '^TestDemo' ./cmd/ailang-worldd -v: ok 82.8s, 0 SKIP. Top-level PASS for WrongKeyCannotProve, NoVerifiedSetRefuses, UncertainCommitStops, Bounds, KeyIORefusals, SeedProofStage (7 arms), AdapterGuards (incl. checker-readfile).
- host/evidence TestProducerEndToEndMintsEnvelopeValidatorAccepts: PASS.
- Allowlist: exactly the 4 cmd files + design_docs/verification/world-iter273/**. No producer/validator/codec edits.

## Score
- Contract/ACs C3 (3 named tests, positive+negative arms, real HTTP, REST read-back): 28/30
- Non-vacuity / mutation kill: 22/25 (two survivors, see non-blocking)
- Correctness of S4 wiring and key semantics: 19/20
- M3a fold-ins (10s bound, 90s bound+cancel, checker-readfile control): 8/10
- Scope/hygiene/evidence: 11/15 (executor ledger not independently re-verified; freshProofKey does a pointless TempDir write)

## Correctness review
- demoReadProofKey mirrors daemon.loadProofKey (open, stat, regular, LimitReader 33, >32 / <32 refusal). Errors name only the path-less prefix and byte COUNT, never bytes.
- S4 runs after S2/S3, uses hashref.SumSHA256(source) (the S1-committed object) and ["sameRef","servesEntry"]. It always returns an error ("S1-S4 recorded; S5-S6 incomplete; no success manifest written"); no manifest is written (asserted by os.Stat in every arm).
- Key refusals happen before any S4 write (writes==0, history revision == S3 revision asserted).
- Uncertain commit: one write, no retry/reseed; store.IsUncertain survives wrapping; lost-ack arm reads the envelope back over REST (daemon cannot mask: payload and revision+1 compared; both unit-level and CLI-level with a Hijack proxy).

## Mutation table (mine; restored after each)
| Mutant | Test | Result | Asserting line |
|---|---|---|---|
| plan M3b-key (validator uses different key) | WrongKeyCannotProve | killed | demo_proof_test.go:706 |
| plan M3b-null-verified (positive source = plain) | NoVerifiedSetRefuses | killed | demo_proof_test.go:732 |
| plan M3b-retry (writer retries once) | UncertainCommitStops | killed | demo_proof_test.go:797 (got 2) |
| open-err guard `if false` | SeedProofStage/missing | killed (message differs: stat: invalid argument) | demo_test.go:1000 |
| stat-err guard | KeyIORefusals/stat | killed | demo_proof_test.go:907 |
| IsRegular guard (`false &&`) | SeedProofStage/directory | killed | demo_test.go:1000 |
| read-err guard | KeyIORefusals/read | killed | demo_proof_test.go:907 |
| >32 guard | SeedProofStage/long | killed | demo_test.go:1000 |
| <32 guard | SeedProofStage/short | killed | demo_test.go:1000 |
| LimitReader 33 -> 32 | SeedProofStage/long | killed | demo_test.go:1000 |
| c.timeout removed | Bounds | killed | demo_proof_test.go:837 |
| c.http.Timeout removed | Bounds | killed | demo_proof_test.go:837 |
| whole bound 2x | Bounds | killed | demo_proof_test.go:842 |
| demoProofContext without deadline | Bounds | killed | demo_proof_test.go:842 |
| demoProve ignores caller ctx | Bounds | killed | demo_proof_test.go:872 |
| checker-readfile message | AdapterGuards/checker-readfile | killed | demo_proof_test.go:583 |
| seed key-refusal guard | SeedProofStage/missing | killed | demo_test.go:1000 |
| seed prove-error guard | SeedProofStage/safe-refusal | killed | demo_test.go:995 |
| required identities -> nil | SeedProofStage/success | killed | demo_test.go:984 |
| wrong source ref | SeedProofStage/success | killed | demo_test.go:984 |
| skip key load | SeedProofStage/missing | killed | demo_test.go:1000 |
| final return nil | SeedProofStage/success | killed | demo_test.go:984 |
| key bytes (%s) leaked in short-key error | SeedProofStage/short | killed | demo_test.go:962 |
| required identities -> ["sameRef"] only | SeedProofStage/success | SURVIVED | none |
| seed uses newClient+10s instead of demoProofClient | Bounds | SURVIVED (near-equivalent) | none |
| key bytes leaked as hex (%x) | SeedProofStage/short | SURVIVED | none |

## Non-blocking findings
1. Dropping "servesEntry" from the seed's required list survives: nothing asserts the CLI path refuses a source lacking servesEntry (the minted envelope carries it anyway). Add a seed-level control or assert the required slice.
2. TestDemoProofBounds asserts demoProofClient(), not that demoSeedDomain actually uses it; swapping to an equivalent hand-built client survives. Weak but behaviourally near-equivalent.
3. Key-disclosure check only searches raw bytes (and the 32-byte key); a hex/base64 leak would survive.
4. freshProofKey now writes an unused proof.key into a TempDir (harmless noise).
5. Executor REPORT/ledger claims were not re-verified beyond my own independent drills.

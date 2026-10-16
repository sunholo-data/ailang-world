# M3b executor report — iteration 273, row 150

Completed M3b in the designated worktree. Started 2026-10-16 06:36:49 UTC; report written 2026-10-16T06:46:42.013795+00:00. No git write operations were performed. Only the four allowlisted checkout files changed.

## Files changed and snapshotted

- `cmd/ailang-worldd/demo_proof.go` — SHA256 `9c0b1f8f474cda2e2db2534c427aad94a95db0c98fc05140a43dd80f11ed4f6b`.
- `cmd/ailang-worldd/demo_proof_test.go` — SHA256 `5386980fd1bf5b707f2849dd52b40a0d073c71ad6fdfd8fabfcf85b917c99942`.
- `cmd/ailang-worldd/demo_seed.go` — SHA256 `846fce2d302a7d1ce75f0a7e52e78816ec0ec8390f11168108bc7b814c74bb11`.
- `cmd/ailang-worldd/demo_test.go` — SHA256 `608cd26fd517d3835aea3a84ad00bf94d67f1c95891af46eefbd2fd20d5d1c57`.

## S4 behavior

After S1 domain commits and successful real S2 read / S3 test and ai_check invocations, seed loads the operator key from o.key. It opens the key once, stats the open file, requires a regular file, reads at most 33 bytes, and requires exactly 32 bytes. Errors contain operation/size information and no key contents. Test keys are independently generated with crypto/rand and written under t.TempDir with mode 0600.

Seed calls demoProve with the production HTTP reader/writer, health-archived v0.41.0 checker, the content ref of the unchanged logepoch source committed by S1, and required identities sameRef and servesEntry. The production client carries both 10s request-context and HTTP-client timeouts, below 12s object bounds. The proof context has a 90s whole bound and propagates cancellation.

A successful S4 yields an envelope ref, but seed still returns “S1-S4 recorded; S5-S6 incomplete; no success manifest written”. No success manifest is created. Every S4 failure stops. store.UncertainError survives wrapping and tells the operator to inspect history, never retry or reseed. No new writable-store shortcut, encoder, checker runner, key fixture, producer/validator change, or key output was introduced.

## Verification coverage

- Real successful HTTP mint: an independent wrong key is unsupported as unauthenticated_report and cannot Resolve PROVEN (UNAVAILABLE); the same key and exact subject Resolve PROVEN.
- Real standalone no-contract scratch module: checked and proved with the actual v0.41.0 pin; ErrVacuousCheck, zero ref, zero writer calls and unchanged REST history. The production logepoch positive control mints sameRef and servesEntry.
- Safe HTTP refusal, uncertain-before-send, and lost acknowledgement after a successful real commit: exactly one write attempt, preserved uncertainty where appropriate, zero returned ref on refusal, no retry/reseed/success. REST object and head reads distinguish absent writes from an actually committed envelope.
- Full seed command from an empty public world with published se-tools: real S1-S4, same-key PROVEN envelope, refusal to finish, no manifest file, no key bytes on stdout/stderr. Missing, directory, 31-byte and 33-byte keys refuse before every S4 write; committed REST inventory/history is checked. Full-command safe-error and lost-ack controls also assert one attempt and no manifest.
- Named checker-readfile negative control, production 10s client bound assertion, whole 90s context deadline assertion, and blocking-reader cancellation with context.Canceled and zero writes.

## Environment and required gates

All commands ran in the sandbox with AILANG_BIN=$HOME/.pinned-ailang/ailang and WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.52.1/ailang. Bounded executable preflight reported v0.41.0/24ee108 and v0.52.1/c68ded4, both rc=0. The pinned language prompt and builtins were consulted before authoring the deliberate no-contract test module. No proof/tool pin substitution was made.

### vet: PASS

Command: `go vet ./cmd/ailang-worldd`

rc=0; wall=0.4s; log: `vet-gate.log`.

No named PASS line applies to this non-test command. Output was empty.

### gofmt: PASS

Command: `gofmt -l cmd/ailang-worldd`

rc=0; wall=0.01s; log: `gofmt-gate.log`.

No named PASS line applies to this non-test command. Output was empty as required by gofmt.

### wrong-key: PASS

Command: `go test -run '^TestDemoProofWrongKeyCannotProve$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

rc=0; wall=1.81s; log: `wrong-key-gate.log`.

```text
--- PASS: TestDemoProofWrongKeyCannotProve (1.17s)
```

### no-verified: PASS

Command: `go test -run '^TestDemoProofNoVerifiedSetRefuses$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

rc=0; wall=1.88s; log: `no-verified-gate.log`.

```text
--- PASS: TestDemoProofNoVerifiedSetRefuses (1.24s)
```

### uncertain: PASS

Command: `go test -run '^TestDemoProofUncertainCommitStops$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`

rc=0; wall=3.83s; log: `uncertain-gate.log`.

```text
--- PASS: TestDemoProofUncertainCommitStops (3.19s)
```

### demo: PASS

Command: `go test -run '^TestDemo' ./cmd/ailang-worldd -count=1 -timeout 400s -v`

rc=0; wall=84.54s; log: `demo-gate.log`.

```text
--- PASS: TestDemoProofHTTPProducerRoundTrip (1.11s)
--- PASS: TestDemoProofSourceHashRefusal (0.97s)
--- PASS: TestDemoProofAdapterGuards (0.00s)
--- PASS: TestDemoProofWrongKeyCannotProve (1.17s)
--- PASS: TestDemoProofNoVerifiedSetRefuses (1.26s)
--- PASS: TestDemoProofUncertainCommitStops (3.16s)
--- PASS: TestDemoProofBounds (0.97s)
--- PASS: TestDemoProofKeyIORefusals (0.00s)
--- PASS: TestDemoSeedAuthenticatedCommits (1.08s)
--- PASS: TestDemoSeedPreflightRefusals (33.25s)
--- PASS: TestDemoSeedHelpPreparation (0.00s)
--- PASS: TestDemoSeedCommitFailures (2.72s)
--- PASS: TestDemoSeedReadTestCheck (4.18s)
--- PASS: TestDemoSeedRejectsToolReportRefusals (3.90s)
--- PASS: TestDemoSeedToolGuardControls (0.00s)
--- PASS: TestDemoSeedProofStage (30.13s)
```

### producer: PASS

Command: `go test -run '^TestProducerEndToEndMintsEnvelopeValidatorAccepts$' ./host/evidence -count=1 -timeout 120s -v`

rc=0; wall=1.13s; log: `producer-gate.log`.

```text
--- PASS: TestProducerEndToEndMintsEnvelopeValidatorAccepts (0.56s)
```

All selected tests executed. No SKIP, zero-test selection, sandbox-denied gate or bound exhaustion was accepted as green.

## Full mutation ledger

Each mutation was applied alone, the named test was run red, the original file was restored from an in-memory scratch copy (never git), and the same command was rerun green. Every red below is an intended assertion failure, not a build failure, panic, skip or timeout. Raw red/restored logs and the complete edit strings are preserved in mutation-ledger.json and the correspondingly named *.log files.

The regular-file guard mutant also adds `_ = info` solely to keep the intentionally disabled guard compiling. The three plan mutants use exactly the specified key substitution, consistent no-contract source/hash substitution, and one additional POST with the same valid object after transport uncertainty.

### key-open

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
if err != nil {
		return [32]byte{}, fmt.Errorf("proof key: open: -> if false {
		return [32]byte{}, fmt.Errorf("proof key: open:
```

Command (red and restored): `go test -run ^TestDemoSeedProofStage$/^missing$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: key missing refusal must precede every S4 write.

```text
demo_test.go:1000: key missing refusal must precede every S4 write: writes=0 envelopes=0 err="ailang-worldd demo seed: S1-S3 recorded; S4 key refused: proof key: stat: invalid argument; no success manifest written\n" want="proof key: open:"
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoSeedProofStage (4.19s)
```

### key-stat

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
if err != nil {
		return key, fmt.Errorf("proof key: stat: -> if false {
		return key, fmt.Errorf("proof key: stat:
```

Command (red and restored): `go test -run ^TestDemoProofKeyIORefusals$/^stat$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: key stat error must refuse without bytes.

```text
demo_proof_test.go:907: key stat error must refuse without bytes: err=<nil> zero=false
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofKeyIORefusals (0.00s)
```

### key-regular

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
if !info.Mode().IsRegular() { -> _ = info
 if false {
```

Command (red and restored): `go test -run ^TestDemoSeedProofStage$/^directory$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: key directory refusal must precede every S4 write.

```text
demo_test.go:1000: key directory refusal must precede every S4 write: writes=0 envelopes=0 err="ailang-worldd demo seed: S1-S3 recorded; S4 key refused: proof key: read: read /var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/TestDemoSeedProofStagedirectory2578374642/003: is a directory; no success manifest written\n" want="not a regular file"
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoSeedProofStage (4.28s)
```

### key-read

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
if err != nil {
		return key, fmt.Errorf("proof key: read: -> if false {
		return key, fmt.Errorf("proof key: read:
```

Command (red and restored): `go test -run ^TestDemoProofKeyIORefusals$/^read$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: key read error must refuse without bytes.

```text
demo_proof_test.go:907: key read error must refuse without bytes: err=<nil> zero=false
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofKeyIORefusals (0.00s)
```

### key-long

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
if len(raw) > len(key) { -> if false {
```

Command (red and restored): `go test -run ^TestDemoSeedProofStage$/^long$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: key long refusal must precede every S4 write.

```text
demo_test.go:1000: key long refusal must precede every S4 write: writes=1 envelopes=1 err="ailang-worldd demo seed: S1-S4 recorded; S5-S6 incomplete; no success manifest written; inspect history before any further action\n" want="holds more than 32 bytes"
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoSeedProofStage (4.42s)
```

### key-short

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
if len(raw) < len(key) { -> if false {
```

Command (red and restored): `go test -run ^TestDemoSeedProofStage$/^short$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: key short refusal must precede every S4 write.

```text
demo_test.go:1000: key short refusal must precede every S4 write: writes=1 envelopes=1 err="ailang-worldd demo seed: S1-S4 recorded; S5-S6 incomplete; no success manifest written; inspect history before any further action\n" want="holds 31 bytes"
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoSeedProofStage (4.28s)
```

### seed-key-refusal

File: `cmd/ailang-worldd/demo_seed.go`.

Edit:

```go
if err != nil {
		return hashref.HashRef{}, fmt.Errorf("S1-S3 recorded; S4 key refused: -> if false {
		return hashref.HashRef{}, fmt.Errorf("S1-S3 recorded; S4 key refused:
```

Command (red and restored): `go test -run ^TestDemoSeedProofStage$/^missing$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: key missing refusal must precede every S4 write.

```text
demo_test.go:1000: key missing refusal must precede every S4 write: writes=1 envelopes=1 err="ailang-worldd demo seed: S1-S4 recorded; S5-S6 incomplete; no success manifest written; inspect history before any further action\n" want="proof key: open:"
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoSeedProofStage (4.29s)
```

### seed-proof-refusal

File: `cmd/ailang-worldd/demo_seed.go`.

Edit:

```go
if err != nil {
		return hashref.HashRef{}, fmt.Errorf("S1-S3 recorded; S4 refused: -> if false {
		return hashref.HashRef{}, fmt.Errorf("S1-S3 recorded; S4 refused:
```

Command (red and restored): `go test -run ^TestDemoProofUncertainCommitStops$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: S4 refusal must stop without success or reseed.

```text
demo_proof_test.go:800: S4 refusal must stop without success or reseed: ref= err=S1-S4 recorded; S5-S6 incomplete; no success manifest written; inspect history before any further action
demo_proof_test.go:800: S4 refusal must stop without success or reseed: ref= err=S1-S4 recorded; S5-S6 incomplete; no success manifest written; inspect history before any further action
demo_proof_test.go:800: S4 refusal must stop without success or reseed: ref= err=S1-S4 recorded; S5-S6 incomplete; no success manifest written; inspect history before any further action
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofUncertainCommitStops (3.27s)
```

### checker-readfile

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
if err != nil {
		return evidence.CheckerConfig{}, fmt.Errorf("proof checker: read archived interpreter: -> if false {
			return evidence.CheckerConfig{}, fmt.Errorf("proof checker: read archived interpreter:
```

Command (red and restored): `go test -run ^TestDemoProofAdapterGuards$/^checker-readfile$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: checker-readfile must name the archived interpreter read failure.

```text
demo_proof_test.go:583: checker-readfile must name the archived interpreter read failure: proof checker: archived pin bytes hash to sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855; health pins sha256:c7cd479e66a549cb24e986a7e40a3d8dc5a5f0eb9ac55edad754a23dc405e6d5
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofAdapterGuards (0.00s)
```

### M3b-key

File: `cmd/ailang-worldd/demo_proof_test.go`.

Edit:

```go
same := proofValidator(t, c, h, key) // M3b-key: -> same := proofValidator(t, c, h, different) // M3b-key:
```

Command (red and restored): `go test -run ^TestDemoProofWrongKeyCannotProve$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: same-key positive arm must Resolve PROVEN.

```text
demo_proof_test.go:706: same-key positive arm must Resolve PROVEN for the exact subject: valid=false proven=false grade=0
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofWrongKeyCannotProve (1.20s)
```

### M3b-null-verified

File: `cmd/ailang-worldd/demo_proof_test.go`.

Edit:

```go
source := demoProofSource(t) // M3b-null-verified: -> source := plain // M3b-null-verified:
```

Command (red and restored): `go test -run ^TestDemoProofNoVerifiedSetRefuses$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: positive proof-control must mint required identities.

```text
demo_proof_test.go:732: positive proof-control must mint required identities sameRef and servesEntry: ref= calls=0 err=evidence: ai-check run verified nothing: 0 results records from non-empty output
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofNoVerifiedSetRefuses (1.28s)
```

### M3b-retry

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
// The request may already have reached the daemon: the outcome is unknown. -> _, _, _ = w.c.doAuth(ctx, http.MethodPost, "/v1/commit", bytes.NewReader(data), "Bearer "+w.token)
		// The request may already have reached the daemon: the outcome is unknown.
```

Command (red and restored): `go test -run ^TestDemoProofUncertainCommitStops$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: write attempt count must equal one.

```text
demo_proof_test.go:797: write attempt count must equal one: got 2
demo_proof_test.go:797: write attempt count must equal one: got 2
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofUncertainCommitStops (3.19s)
```

### http-context-bound

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
c.timeout = demoProofHTTPTimeout -> c.timeout = 0
```

Command (red and restored): `go test -run ^TestDemoProofBounds$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: production client/adapters must carry the 10s HTTP bound.

```text
demo_proof_test.go:837: production client/adapters must carry the 10s HTTP bound below the 12s object bound
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofBounds (0.99s)
```

### http-client-bound

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
c.http.Timeout = demoProofHTTPTimeout -> c.http.Timeout = 0
```

Command (red and restored): `go test -run ^TestDemoProofBounds$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: production client/adapters must carry the 10s HTTP bound.

```text
demo_proof_test.go:837: production client/adapters must carry the 10s HTTP bound below the 12s object bound
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofBounds (0.99s)
```

### whole-bound

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
return context.WithTimeout(ctx, demoProofWholeBound) -> return context.WithCancel(ctx)
```

Command (red and restored): `go test -run ^TestDemoProofBounds$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: whole proof context must carry the 90s deadline.

```text
demo_proof_test.go:842: whole proof context must carry the 90s deadline
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofBounds (1.00s)
```

### cancel-propagation

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
ctx, cancel := demoProofContext(ctx) -> ctx, cancel := demoProofContext(context.Background())
```

Command (red and restored): `go test -run ^TestDemoProofBounds$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: cancelled producer did not stop.

```text
demo_proof_test.go:872: cancelled producer did not stop
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoProofBounds (1.00s)
```

### key-read-cap

File: `cmd/ailang-worldd/demo_proof.go`.

Edit:

```go
io.LimitReader(f, 33) -> io.LimitReader(f, 32)
```

Command (red and restored): `go test -run ^TestDemoSeedProofStage$/^long$ ./cmd/ailang-worldd -count=1 -timeout 120s -v`.

Observed red rc=1; intended assertion: key long refusal must precede every S4 write.

```text
demo_test.go:1000: key long refusal must precede every S4 write: writes=1 envelopes=1 err="ailang-worldd demo seed: S1-S4 recorded; S5-S6 incomplete; no success manifest written; inspect history before any further action\n" want="holds more than 32 bytes"
```

Restored-and-green rc=0:

```text
--- PASS: TestDemoSeedProofStage (4.22s)
```

## Protected paths and scope

Before/after SHA256 and newline inventories match for all 12 observed protected files: observation-before.json == observation-after.json. The exact protected-path git diff command returned rc=0 with empty output. git diff --check returned rc=0 with empty output.

The full repository suite, verify_go.sh and verify_ail.sh were not run, as directed. No files outside the allowlist were modified in the checkout. Snapshots contain only source, reports and test logs; no key/session files or credential bytes were copied.

## Unfinished work and findings

No M3b implementation or required gate remains unfinished. S5/S6 approval/why and later roster/manifest milestones remain intentionally incomplete; seed refuses completion accordingly.

No blocking implementation finding or required-gate sandbox denial was observed. UNINFORMATIVE UNDER SANDBOX: an optional read-only ps process-list diagnostic returned “operation not permitted”; it was not used as a pass/fail gate and no workaround was attempted.

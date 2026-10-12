# Sprint plan — `w-bounded-z3-report-producer`

**Iteration 252 · queue row 26 · clause-2 · v1.0.0 release prerequisite (D-WORLD-69 = A)**

Design: [w-bounded-z3-report-producer.md](w-bounded-z3-report-producer.md). Quorum r1 BLOCKED, r2 BLOCKED, r3 = narrow-refinement carve-out applying the reviewers' verbatim fixes; `gemini-3-1-pro` passed both rounds. Worktree `.wt-world-iter251-row26-design`, branch `design/w-bounded-z3-report-producer`, HEAD `4dfdba2`. Code base is `ba79d1e` plus design-doc commits only. Binding: CLAUDE.md and design_docs/coding-standards.md (effects at the host boundary, non-vacuous gates, pinned binary). The design's §7 ACs (AC5, AC20, AC23–AC30) and §8 mutants (M15, M28, M31–M39) are reproduced here in full. **~2.3 executor days.**

**Plan only.** The planner writes exactly two files: this one and `sprint_w-bounded-z3-report-producer.json`. **The EXECUTOR runs no git command** (no add, commit, stash, checkout or branch). By default the executor implements M1→M8 in order in this worktree. Every milestone's gate list must be green at its own boundary before the next milestone starts. At each boundary the executor records a boundary report: changed files, the rc and wall time of each gate, the red-first evidence, the mutation ledger and restore evidence. **The controller commits** from the executor's worktree state. If the controller tells the executor to STOP at milestone boundaries, the controller commits C1…C8 one per boundary. Otherwise it commits from the final tree. Either way every milestone boundary is internally consistent (all gates green), so per-milestone commits are possible.

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang   # v0.41.0, sha256 1a67b014…5b9f (re-measured, §Baseline)
export PATH=/opt/homebrew/bin:$PATH
```

## Baseline (planner-measured at 4dfdba2, first-party)

| Gate | Result |
|---|---|
| `go vet ./host/evidence/...` | rc=0, 0.23 s |
| `go test ./host/evidence -count=1` | `ok` 1.6 s; 37 unique terminal top-level passes |
| `go test -race ./host/evidence -count=1` | `ok` 6.1 s |
| `go test -json ./host/evidence -count=1 > f; ./scripts/verify_go.sh --evidence-manifest-check f 37` | `✓ all 37 required top-level evidence tests passed exactly once`, rc=0. With 38: `observed_unique=37 exact_required=38`, rc=1 (the isolated mode is a working per-milestone gate) |
| `go test -count=1 ./host/verifygate -run '^(TestEvidenceNamedManifestRejectsUnpinnedTest\|TestVerifyGoProductCheckSequence\|TestSubprocessSinkGate)$' -v` | 3× PASS, 1.4 s |
| `go test -count=1 ./host/verifygate` (full) | `ok` 157.5 s |
| `AILANG_BIN=… ./scripts/verify_ail.sh` | rc=0, 33 s; `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 444 se-tools named tests pass` (log `~/.ailang/state/world-iter252/planner_base_ail.log`) |
| `AILANG_BIN=… ./scripts/verify_go.sh` | rc=0, **10 m 18 s wall** locally; `✓ all 37 required…`, `✓ go gate PASSED` (log `~/.ailang/state/world-iter252/planner_base_go.log`) |

`verify_go.sh` takes more than 600 s here. Run it with a timeout of at least 15 min (or let the controller run it). A 600 s tool cap is not a gate result.

## Design defects found by the planner (each re-measured first-party; each resolved by a numbered PD below)

1. **§9 / Files are wrong that "the verifygate self-mutation gate is unchanged by an additive edit".** The literal `37` and the 37-name list are hard-coded in two more files that run in `go test ./...`:
   - `host/verifygate/evidence_manifest_gate_test.go:16-54` holds `syntheticEvidenceTests`, and `:178` holds `"all 37 required top-level evidence tests passed exactly once"`.
   - `host/verifygate/mission_config_gate_test.go:466-467` holds the repaired-script needles `…(37 exact top-level tests)` and `check_evidence_manifest "$evidence_json" 37`.

   Command: `/usr/bin/grep -rn "37 exact\|evidence_json\" 37\|all 37 required" --include='*.go' --include='*.sh' .` → those 4 Go hits plus verify_go.sh:295/:302. If the manifest edit misses them, it reds `go test ./host/verifygate`. → **PD1**.
2. **§9 is wrong that "`TestPublicAuthoritySurfaceIsFrozen` must still pass unedited".** That test is an EXACT manifest of every exported declaration in package evidence (`frozenPublicSurface`, authority_test.go:286-340, compared with `reflect.DeepEqual`). `NewProducer`, `Producer`, `CheckerConfig`, `ProducerConfig`, `ObjectWriter`, the 8 new sentinels, `MaxSourceBytes` and `Producer.GenerateProof` each red it. → **PD2**.
3. **AC29's `durableAfterCommitHook` stimulus cannot be reached.** The hook is an unexported package var in `host/store` (durable.go:46; only `host/store/*_test.go` reference it). Go gives `host/evidence` tests no path to it, and the design's own Non-Goals forbid any store change. → **PD3**.
4. **AC29's oversize stimulus ("256 identities at the codec's 1 KiB string cap") is the wrong size.** `EncodeAuthenticatedEnvelope` first runs `DecodeProofReportV1(e.Report)`, which refuses any report over `MaxBytes`. So that report is refused by the codec before the `ErrEnvelopeTooLarge` branch runs, and **M37 would be unkillable**. Measured with a scratch module using `replace` to this worktree:
   - 256×1024 → `report 263235 … decode refused: oversize`;
   - 200×1000 → `report 200923 envelope 267963 … over true`;
   - 190×1024 → `report 195453 envelope 260669 … over false`.

   → **PD4**.
5. **AC25/M39 is wrong that a real `ai-check` run takes "multi-second".** Measured on the pinned binary: verifiedpass `ai-check` is 0.02–0.05 s wall (`/usr/bin/time -p` → `real 0.02`). Also, the version probe shares `ExecTimeout`, so a tiny timeout against the real binary would expire in the probe, not in the run. → **PD5**.
6. **AC26's stderr arm cannot overflow on its own with the real binary.**
   - Under the design's own staging (`staged/module.ail`, `AILANG_RELAX_MODULES=1`), every fixture writes a ~200-byte `MOD010` warning to stderr. Stdout is always larger (250–947 bytes), so no cap can overflow stderr while stdout stays under it.
   - The committed design fixtures' `.stderr` files are all 0 bytes (gate layout, no relax).

   → **PD5**.
7. **D4 step 4 (`verify.available`) can never be pinned by a real fixture.** The only real `available=false` shapes (loaderror, and a measured type-error module) also have `check.passed=false`, which is checked first. Measured: `{"p":false,"a":false}` for both. Z3 is resolved even with `PATH=/nonexistent` (measured `available:true`). Step 9 has the same problem for the status filter, because every real non-verified status also bumps its own counter. This is the design's own "a guard nothing can pin is not a guard" class. → **PD6**.
8. **AC20/AC23's decoy classification measures the wrong interval.** It times all of `GenerateProof` against 2× the bound. The producer hashes the 100,594,162-byte binary and runs two subprocesses after the read, so even a correct mint takes ≫ 2× a 2.6 ms bound. A correct program would red as the mutant. Also, the producer derives `readCtx`/`writeCtx` *before* its wrapper runs, so a wrapper that sleeps to settle the decoy uses up the deadline and passes vacuously. → **PD7**.
9. **The design's M1 says "the manifest update lands with M5, not before".** But the manifest is an EXACT set (verify_go.sh:97-109), so every intermediate commit would red `verify_go.sh`. → **PD1** (consistent at every boundary, as the brief requires).

The inherited premises listed below were re-verified, so they are not taken on trust.

**Re-verified:** binary identity (`AILANG v0.41.0`, sha256 `1a67b014…5b9f`, 100,594,162 bytes, `--version` stderr 0 bytes); the regen reproduces all 7 committed fixtures exactly apart from `duration`; `TMPDIR` redirect leaves 0 entries; `PutObject`/`ReadObject` wrap ctx errors with `%w`; `SetMaxOpenConns(1)`; default `BusyTimeout()` 2 s; the `procbound` API; capsule's shapes; that `host/` is not in `verify_ail.sh` ROOTS; that host/evidence imports nothing that would cycle with `host/archive`; and that `TestSubprocessSinkGate` refuses `bytes.Buffer`/`strings.Builder` bound to `cmd.Stdout/Stderr` in a function that calls `Start`.

## Planner decisions (numbered; each binding on the executor)

**PD1 — manifest consistent at EVERY boundary.**

- Each milestone adds its own new top-level test names, in alphabetical position, to ALL of:
  - (a) `REQUIRED_EVIDENCE_TESTS` in `scripts/verify_go.sh`;
  - (b) `syntheticEvidenceTests` in `host/verifygate/evidence_manifest_gate_test.go`.
- Each milestone sets the count N, in the same edit, in ALL of:
  - verify_go.sh banner (:295), `(N exact top-level tests)`;
  - verify_go.sh call (:302), `check_evidence_manifest "$evidence_json" N`;
  - evidence_manifest_gate_test.go:178, `"all N required top-level evidence tests passed exactly once"`;
  - the repaired-arm needles of `requireVerifyGoProductSequence`.
- **M1 restructures `requireVerifyGoProductSequence` once.** It mirrors the existing `verTokAnchor` pattern:
  - declare `evidenceBanner := \`echo "── focused host/evidence named-manifest gate (37 exact top-level tests)"\`` and `evidenceCall := \`check_evidence_manifest "$evidence_json" 37\``;
  - override both inside `if repaired { … }` with N;
  - replace the two literal needles in `needles` with `evidenceBanner, evidenceCall`.
- **The embedded `baseVerifyGoProductFixture` (lines 423-424) stays 37 forever.** It is a historical 565d0b2 snapshot.
- Counts: M1=47, M2=50, M3=51, M4=54, M5=56, M6=58, M7=61, M8=63 (26 new tests).
- The `TestAttackerChosenValidatorCannotMintForHostAuthority` entry must remain `syntheticEvidenceTests[0]` (the "missing" subtest drops index 0 and greps that exact name). Every new name sorts after it.

**PD2 — `frozenPublicSurface` extended per milestone, each addition reviewed against AC1.**

- Add, in the milestone that introduces the symbol:
  - M1: `"func NewProducer"`, `"type CheckerConfig"`, `"type ObjectWriter"`, `"type Producer"`, `"type ProducerConfig"`, `"value ErrInvalidProducerConfig"`;
  - M2: `"value ErrExecTimeout"`, `"value ErrOutputLimit"`;
  - M4: `"value ErrUntrustedTool"`;
  - M5: `"value ErrUnverifiedRun"`, `"value ErrVacuousCheck"`;
  - M6: `"value ErrEnvelopeTooLarge"`, `"value ErrRequiredIdentityMissing"`;
  - M7: `"method Producer.GenerateProof"`, `"value MaxSourceBytes"`.
- That is 15 additions. No other exported identifier may be introduced; every helper is unexported.
- AC1 review recorded in the boundary report: no new function returns `ResolvedGrade`, `ResolutionResult` or `ValidatedEvidence`. `TestPublicAuthoritySurfaceIsFrozen`'s reflective checks stay unedited. The only edit is to the list.

**PD3 — uncertain-write stimulus.** `TestProducerSurfacesUncertainWriteOutcome` uses `uncertainWriter{s *store.Store; calls *atomic.Int32}`:

- `PutObject` calls the REAL `s.PutObject(ctx, o)` (a real commit), then returns `&store.UncertainError{Op: "put object", Cause: context.DeadlineExceeded}`;
- `BusyTimeout()` forwards `s.BusyTimeout()`.

This is a wrapper over the real store (AC30's wrapper precedent), not a fake store. It is the only way to reach the classification without editing `host/store`.

**PD4 — oversize stimulus.** `TestProducerRefusesEnvelopeBeyondValidatorReadCap`:

- **Over arm:** 200 identities, each `fmt.Sprintf("id%03d", i) + strings.Repeat("x", 995)` (1000 bytes, sorted, unique). Report 200,923 B (< MaxBytes) → envelope 267,963 B (> MaxBytes) → `ErrEnvelopeTooLarge`.
- **Near-bound control arm:** 190 identities × 1024 bytes (`"id%03d"` + 1019 `x`). Envelope 260,669 B mints, and the test asserts `len(envelope) <= evidence.MaxBytes`.

Both use the real codec. Measured numbers are above.

**PD5 — scripted pinned executables for the stderr, exec-timeout and admission arms (capsule precedent, `host/capsule/cleanup_test.go:211-227`).**

- The test writes `#!/bin/sh` scripts into `t.TempDir()` and pins them in `CheckerConfig` by their own bytes (`Compiler: hashref.SumSHA256(scriptBytes)`, `CompilerVersion` = the script's `--version` line).
- The real binary still drives the stdout arm (M3), the version/byte arms (M4), fixtures (M5) and everything in M7/M8.
- This is the measured fix for defects 5 and 6. The design's "a REAL `ai-check` run … under 16 bytes" holds only for stdout.

**PD6 — derived single-field arms for guards no real fixture can pin.**

- `TestProducerInterpretationMatchesEveryFixtureShape` adds three arms. Each is built in-test by JSON-decoding the committed `verifiedpass.json` into `map[string]any`, flipping exactly ONE field, and re-marshalling:
  - `check.passed=false`;
  - `verify.available=false`;
  - `verify.results[0].status="unknown"`.
- Each would otherwise MINT `twicePlusOne`, so each pins its branch.
- The real `loaderror` arm also asserts the refusal detail `check.passed`. Neutering the passed branch then reds on detail even though the available branch refuses with the same sentinel.
- These are labelled derivations of a real fixture, never hand-typed tool output.
- The two literal inputs `""` (empty stdout) and `"{"` (unparseable) are trivial non-tool inputs and are allowed as literals.

**PD7 — decoy arms time the STORE CALL and settle inside the derived budget.** AC20/AC23 use a forwarding wrapper over the real store. The wrapper's `BusyTimeout()` forwards the real value, which is 0 for `openEvidenceStoreWithoutBusyWait`. When armed, the wrapper does the following on entry to `ReadObject` (AC23) or `PutObject` (AC20):

1. call the existing `startDecoyRead(t, s, decoyRef)`;
2. `time.Sleep(settle)` with `settle = hold/40`;
3. assert the decoy is still in flight (`select { case <-decoyDone: instrument failure; default: }`);
4. assert `ctx.Err() == nil` (deadline not already spent; otherwise count it as an instrument retry);
5. time ONLY the real store call (`callElapsed`) and return its result.

The derived bound is `hold/20` (= settle + hold/40). So the decoy hold is ≥ 20× the bound, as the design requires.

Classification:

| Outcome | Verdict |
|---|---|
| refusal with `errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, store.ErrNoDeadline)`, zero ref, no write | PASS |
| mint with `callElapsed < 2×bound` | the store call won the connection; wait for the decoy, then retry (≤5; exhaustion is a loud instrument failure) |
| mint with `callElapsed ≥ 2×bound` | the M31/M32 mutant signature; immediate red |

The test-side watchdog on the whole `GenerateProof` is `20*hold + 60*time.Second`. It catches only a true hang, because a correct mint includes hashing the binary and two subprocesses. The opposite timeout is set to `3*time.Second + 2*hold`: `ObjectWriteTimeout` in AC23 and `ObjectReadTimeout` in AC20. This stops the decoy from also starving the un-decoyed leg on a slow CI runner (hold measured at 2.63 s there; realstore_test.go:22-25).

**PD8 — milestone split and order.** The design's M1–M5 is re-cut into M1–M8, each ≤ ~260 LOC including tests:

| Plan | Design | Content |
|---|---|---|
| M1 | M1 | constructor, config, pins |
| M2 | M3 part a | bounded runner core (`runBounded`, `readCapped`, `killGroup`) |
| M3 | M3 part b | staging + `ai-check` invocation (`runCheck`) + fixtures |
| M4 | M2 | executable verification (`verifyTool`) |
| M5 | M4 part a | strict interpretation (`interpretCheck`) + revalidation |
| M6 | M4 part b | identities (`requireIdentities`) + mint (`mint`) |
| M7 | M5 part a | store legs + `GenerateProof` composition |
| M8 | M5 part b | real-store decoys + full re-drill + final gates |

Design M3 runs BEFORE design M2 because the version probe must reuse the bounded runner. `TestSubprocessSinkGate` forbids the easy `bytes.Buffer` probe, and the design's Deferred list says "both streams stay capped". Stage tests at M2–M6 call the unexported stage functions from an internal test file. `GenerateProof` lands at M7 and composes them, so no stub is ever committed.

**PD9 — AC5 constructor test spellings.** The design abbreviates them as `...X`. The resolved names are:
`TestProducerConstructorRefusesNonPositiveObjectReadTimeout`, `TestProducerConstructorRefusesNonPositiveObjectWriteTimeout`, `TestProducerConstructorRefusesNonPositiveExecTimeout`, `TestProducerConstructorRefusesNonPositiveMaxOutputBytes`, `TestProducerConstructorRefusesNilReader`, `TestProducerConstructorRefusesNilWriter`, `TestProducerConstructorRefusesEmptyExecPath`, `TestProducerConstructorRefusesUnsetCompilerIdentity`, `TestProducerConstructorRefusesEmptyCompilerVersion`.

**PD10 — the empty `requiredIdentities` refusal lives in ONE place.** It is in `requireIdentities` (sentinel `ErrRequiredIdentityMissing`, detail `no required identities named`). `GenerateProof` reaches it after interpretation. There is no duplicate early guard: two guards with one sentinel make one of them unkillable. A bounded binary run for a caller error is acceptable cost.

**PD11 — AC30's "ordered arm … produces end-to-end"** is discharged at M7 by `TestProducerEndToEndMintsEnvelopeValidatorAccepts`. That test constructs on the DEFAULT-window real store (2 s) with 3 s read/write timeouts, which is the ordered configuration. AC30's own test (M1) asserts that the ordered construction succeeds. No test body is edited after its milestone.

**PD12 — no `t.Parallel()` anywhere in host/evidence tests.** The package now both writes files and forks (`TestNoParallelWriteForkPackagesOutsideVerifygate`, and the shared global `killGroup`). No test may bind a `bytes.Buffer`/`strings.Builder` to `cmd.Stdout`/`cmd.Stderr` (`TestSubprocessSinkGate`). Tests that exec use `cmd.Output()` or the producer's own runner.

## Frozen / must-not-change

- `host/store/**`: no file changes at all.
- `host/evidence/validator.go`, `report_codec.go`, `envelope_codec.go`, `proposal_codec.go`, `json.go`, `types.go`: no hunk.
- `host/evidence/realstore_test.go`, `codec_test.go`, `m7b_deadline_test.go`: no hunk. New tests reuse their helpers.
- `host/evidence/authority_test.go`: only `frozenPublicSurface` list additions (PD2).
- `host/procbound`, `host/proctest`, `host/childenv`, `host/archive`, `host/capsule`, `host/broker`: no hunk.
- Every `.ail` under `world/`, `packages/`, `design_docs/`: no hunk. Fixture `.ail` files go only under `host/evidence/testdata/aicheck/`.
- `scripts/verify_ail.sh`, `.github/workflows/ci.yml`, `tools/launchd/*`, `design_docs/world-mission*.md`: no hunk.
- `host/verifygate/*`: only the PD1 edits in `evidence_manifest_gate_test.go` and `mission_config_gate_test.go`. `baseVerifyGoProductFixture` is untouched.
- `design_docs/verification/world-iter251/**`: read-only source of the fixture modules.

## Test write rule

Every DB, scratch root, script, pid file and fresh `ai-check` output goes under the test's own `t.TempDir()`. Tests never write tracked files, never write `$HOME` or `~/.ailang`, and never set `AILANG_CACHE_DIR` to a shared path. The pinned binary and the committed fixtures are read-only inputs.

Tracked fixtures are written ONLY by `host/evidence/testdata/aicheck/regen.sh`, run by the executor (never by a test). Every pid read from a fixture file goes through `proctest.ReadPid`, and every leftover process through `proctest.ReapOnCleanup`. **Never signal a pid that has not been checked** (the iter-190/191 `kill(-1)` lesson).

Binary-dependent tests read `AILANG_BIN` and call `t.Fatal` when it is unset or not `AILANG v0.41.0`. They NEVER `t.Skip` (mirror `host/archive/check_test.go:79-81` and verifygate `requirePinned`). A skipped manifest test reds `verify_go.sh` by design.

## Mutation protocol

1. Apply one mutant alone to the boundary tree.
2. It must compile: `go vet ./host/evidence/...` rc=0.
3. Run its exact command. The named test must FAIL on its named assertion. Build failure, panic outside a recovered-arm helper, SKIP, zero tests matched, an unrelated failure or a harness timeout are NOT kills.
4. Restore the exact original bytes without git (keep a copy in `$TMPDIR`) and confirm with `cmp` that the file is byte-identical to the copy.
5. Re-run the command and confirm green.

Record each mutant as `mutant | file:line | command | rc | failing assertion text | restore cmp rc | green rc`. A SURVIVED mutant blocks the milestone.

At M8, re-drill every design-numbered mutant (M15, M28, M31–M39) once more against the final tree.

## Shared per-milestone gate template (N = that milestone's count; NEW = its new test names joined by `|`)

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang; export PATH=/opt/homebrew/bin:$PATH
go vet ./host/evidence/... ./host/verifygate/...                       # compile fence incl. _test.go (go build is NOT one)
log="$(mktemp)"; go test -count=1 -v ./host/evidence -run "^(NEW)$" >"$log" 2>&1; rc=$?; cat "$log"; [ $rc -eq 0 ]
for t in $(printf '%s' "NEW" | tr '|' ' '); do grep -q -- "^--- PASS: $t " "$log" || { echo "$t did not PASS"; exit 1; }; done
go test ./host/evidence -count=1
go test -race ./host/evidence -count=1
ev="$(mktemp)"; go test -json ./host/evidence -count=1 >"$ev"; ./scripts/verify_go.sh --evidence-manifest-check "$ev" N
go test -count=1 ./host/verifygate -run '^(TestEvidenceNamedManifestRejectsUnpinnedTest|TestVerifyGoProductCheckSequence|TestSubprocessSinkGate|TestNoParallelWriteForkPackagesOutsideVerifygate)$'
```

The executor runs this under bash, and the bash `for` splits on spaces. In zsh `for t in $X` does NOT split, so use bash.

## Code shape (binding names; one production file `host/evidence/proof_producer.go`)

Exported:

- `ObjectWriter` interface: `PutObject(context.Context, store.Object) error`, `BusyTimeout() time.Duration`.
- `CheckerConfig{ExecPath string; Compiler hashref.HashRef; CompilerVersion string}`.
- `ProducerConfig{ObjectReadTimeout, ObjectWriteTimeout, ExecTimeout time.Duration; MaxOutputBytes int64}`.
- `Producer` (unexported fields `key *[32]byte`, `reader ObjectReader`, `writer ObjectWriter`, `checker CheckerConfig`, `cfg ProducerConfig`).
- `NewProducer(key [32]byte, reader ObjectReader, writer ObjectWriter, checker CheckerConfig, cfg ProducerConfig) (*Producer, error)`.
- `(*Producer).GenerateProof(ctx context.Context, sourceRef hashref.HashRef, requiredIdentities []string) (hashref.HashRef, error)`.
- `const MaxSourceBytes = 1 << 20`.
- Sentinels (`errors.New("evidence: …")`): `ErrInvalidProducerConfig`, `ErrUntrustedTool`, `ErrVacuousCheck`, `ErrRequiredIdentityMissing`, `ErrUnverifiedRun`, `ErrEnvelopeTooLarge`, `ErrExecTimeout`, `ErrOutputLimit`.
- Reused: `ObjectReader`, `ErrUnorderedTimeouts`, `procbound.ErrCleanupBacklog`, `procbound.ErrCleanupIncomplete`.

Unexported:

- `const pipeCloseGrace = time.Second`; `const versionProbeMaxBytes = 4 << 10`; `const aiCheckZ3Timeout = "5s"`; `const stagedModule = "staged/module.ail"`.
- `var killGroup = func(pgid int) error { return syscall.Kill(-pgid, syscall.SIGKILL) }`.
- `func readCapped(pipe io.Reader, limit int64) ([]byte, error)`.
- `func (p *Producer) runBounded(ctx context.Context, dir string, env []string, limit int64, args ...string) ([]byte, error)`.
- `func (p *Producer) runCheck(ctx context.Context, source []byte) ([]byte, error)`.
- `func (p *Producer) verifyTool(ctx context.Context) error`.
- `func interpretCheck(stdout []byte) ([]string, error)`.
- `func requireIdentities(verified, required []string) error`.
- `func (p *Producer) mint(subject hashref.HashRef, verified []string) ([]byte, hashref.HashRef, error)`.
- `func (p *Producer) readSource(ctx context.Context, ref hashref.HashRef) ([]byte, error)`.

Test files:

- `host/evidence/producer_test.go` (`package evidence_test`; M1 constructor tests).
- `host/evidence/producer_internal_test.go` (`package evidence`; M2–M6 stage tests; it may swap `killGroup`).
- `host/evidence/producer_realstore_test.go` (`package evidence_test`; AC30 at M1, then M7/M8; reuses `realstore_test.go`/`authority_test.go` helpers `openEvidenceStore`, `openEvidenceStoreWithoutBusyWait`, `openEvidenceStoreWithBusyWindow`, `measureGetObject`, `startDecoyRead`, `minDecoyHold`, `decoyHoldRatio`, `boundedTestContext`, `fakeReader`, `testKey`).

---

## C1 / M1 — Constructor, config, ordering pins (design M1)

**Production budget ~90 LOC; tests ~170.** Files:

- `host/evidence/proof_producer.go` (new)
- `host/evidence/producer_test.go` (new)
- `host/evidence/producer_realstore_test.go` (new)
- `host/evidence/authority_test.go` (PD2 list only)
- `scripts/verify_go.sh` (PD1, N=47)
- `host/verifygate/evidence_manifest_gate_test.go` (PD1)
- `host/verifygate/mission_config_gate_test.go` (PD1 restructure + 47)

**Implementation:** `NewProducer` makes these refusals, in this order. Each one is a separate `if`, and each detail names only its own field:

1. `reader == nil` → `"%w: reader is nil"`;
2. `writer == nil` → `"%w: writer is nil"`;
3. `checker.ExecPath == ""` → `"%w: ExecPath is empty"`;
4. `checker.Compiler.IsZero()` → `"%w: Compiler identity is unset"`;
5. `checker.CompilerVersion == ""` → `"%w: CompilerVersion is empty"`;
6. `cfg.ObjectReadTimeout <= 0`, `cfg.ObjectWriteTimeout <= 0`, `cfg.ExecTimeout <= 0` → `"%w: <Field> must be positive, got %s"`;
7. `cfg.MaxOutputBytes <= 0` → `"%w: MaxOutputBytes must be positive, got %d"`.

All of these use `ErrInvalidProducerConfig`.

Then, for `(seam, w)` in `[("reader", reader.BusyTimeout()), ("writer", writer.BusyTimeout())]`:

- `w < 0` → `"%w: %s BusyTimeout is unknown (%s)"` with `ErrInvalidProducerConfig`;
- `w > 0 && cfg.ObjectReadTimeout <= w` → `"%w: ObjectReadTimeout %s must exceed %s BusyTimeout %s"` with `ErrUnorderedTimeouts`;
- `w > 0 && cfg.ObjectWriteTimeout <= w` → the same text with `ObjectWriteTimeout`.

These are separate `if`s, not a compound condition. Finally copy the key exactly like validator.go:91-92 and copy `checker`/`cfg` by value.

**Tests first.** On the base tree they are build-red (missing API). The red-first proof for each is its mutation drill below on compiling code.

- **AC5 constructor arm**, in `producer_test.go`, with helper `validProducerArgs()`:
  - Valid arguments: `fakeReader{}` and a new `fakeWriter{busy}` (`PutObject` returns nil); `CheckerConfig{ExecPath: "/nonexistent/ailang", Compiler: testCompiler, CompilerVersion: "AILANG v0.41.0"}`; `ProducerConfig{3s, 3s, time.Minute, 1<<20}`.
  - Helper `requireConfigRefusal(t, field string, build func() error)` recovers a panic into `t.Fatalf("NewProducer panicked: %v; want ErrInvalidProducerConfig naming %s")`. It requires `errors.Is(err, evidence.ErrInvalidProducerConfig)` and `strings.Contains(err.Error(), field)`.
  - Tests:
    - `TestProducerConstructorRefusesNonPositiveObjectReadTimeout` (0 and −1s)
    - `TestProducerConstructorRefusesNonPositiveObjectWriteTimeout` (0, −1s)
    - `TestProducerConstructorRefusesNonPositiveExecTimeout` (0, −1s)
    - `TestProducerConstructorRefusesNonPositiveMaxOutputBytes` (0, −1)
    - `TestProducerConstructorRefusesNilReader` (detail `reader`)
    - `TestProducerConstructorRefusesNilWriter` (detail `writer`)
    - `TestProducerConstructorRefusesEmptyExecPath` (detail `ExecPath`)
    - `TestProducerConstructorRefusesUnsetCompilerIdentity` (zero `hashref.HashRef{}`, detail `Compiler`)
    - `TestProducerConstructorRefusesEmptyCompilerVersion` (detail `CompilerVersion`)
  - Each test also asserts the valid-args control constructs (`err == nil`, non-nil `*Producer`).
- **AC30**, `TestProducerConstructorPinsBusyTimeoutBelowBothTimeouts` in `producer_realstore_test.go`. All windows come from REAL stores. Arms:
  - **ordered:** default store `s` (2 s) as both seams, read=write=3 s → constructs.
  - **read-unordered:** reader `s` (2 s), writer `openEvidenceStoreWithoutBusyWait` (0), read 1 s, write 3 s → `ErrUnorderedTimeouts`, not `ErrInvalidProducerConfig`, detail contains `ObjectReadTimeout`.
  - **write-unordered:** reader no-busy, writer `s`, read 3 s, write 1 s → `ErrUnorderedTimeouts`, detail `ObjectWriteTimeout`.
  - **writer-window-only:** reader no-busy (0), writer `s` (2 s), read=write=1 s → `ErrUnorderedTimeouts`. This proves the writer's window is read.
  - **unknown:** writer wrapper `busyWriter{s, -time.Nanosecond}` → `ErrInvalidProducerConfig`.
  - **wrapper-preserves-non-default-window:** custom store at a 1500 ms window (`openEvidenceStoreWithBusyWindow`). First assert `custom.BusyTimeout() == 1500ms` (instrument check). Reader = `openEvidenceStoreWithoutBusyWait` (0), writer = `busyWriter{custom, custom.BusyTimeout()}`. Read=write=1800 ms → constructs (ORDERED against the forwarded 1500 ms, UNORDERED against the 2000 ms default, so a wrapper that stopped forwarding and returned the default would refuse). Read=write=1400 ms → `ErrUnorderedTimeouts`. This mirrors `TestReaderWaitBoundsCannotBeLostThroughWrapper`. `busyWriter{s *store.Store; busy time.Duration}` forwards `PutObject` to `s` and returns `busy` from `BusyTimeout()`.

**Acceptance:** the shared template with N=47 and NEW = the 10 names above.

**Mutation drills (M1):**

| Mutant | Isolated edit in `proof_producer.go` | Must red |
|---|---|---|
| M1-a…i | each of the nine refusals → `if false && <cond>` | its own `TestProducerConstructorRefuses…` test |
| M1-j | read-side pin `cfg.ObjectReadTimeout <= w` → `false && …` | AC30 read-unordered arm |
| M1-k | write-side pin `cfg.ObjectWriteTimeout <= w` → `false && …` | AC30 write-unordered arm |
| M1-l | writer window read → `0` (i.e. `writer.BusyTimeout()` replaced by `time.Duration(0)`) | AC30 writer-window-only arm |
| M1-m | `w < 0` → `false && w < 0` | AC30 unknown arm |

Commands: `go test -count=1 -v ./host/evidence -run '^TestProducerConstructor'`.

**Controller commit message:** `feat(evidence): bounded proof producer constructor with per-field refusals and both ordering pins (row 26 M1)`

## C2 / M2 — Bounded process-group runner core (design M3 part a; AC25, AC26-stderr)

**Production ~110 LOC; tests ~150.** Files:

- `proof_producer.go`
- `host/evidence/producer_internal_test.go` (new, `package evidence`)
- authority_test.go (PD2)
- the PD1 trio (N=50)

**Implementation** (capsule.go:230-391 shape, adapted). In `runBounded`:

1. `execCtx, cancel := context.WithTimeout(ctx, p.cfg.ExecTimeout)` ← **the M39 line**.
2. `cmd := exec.CommandContext(execCtx, p.checker.ExecPath, args...)`; `cmd.Dir = dir`; `cmd.Env = env`; `cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}`; `cmd.Cancel = func() error { if cmd.Process == nil { return nil }; return killGroup(cmd.Process.Pid) }`.
3. Take both `StdoutPipe` and `StderrPipe`.
4. `release, err := procbound.Admit()` BEFORE `Start`.
5. On a `Start` error, call `release()`. If `execCtx` has expired → `ErrExecTimeout`; otherwise wrap the error.
6. `closer := time.AfterFunc(p.cfg.ExecTimeout+pipeCloseGrace, close both pipes)`.
7. Drain both pipes concurrently through `readCapped(pipe, limit)`. An overflow on either stream triggers `killOnce.Do(func(){ killErr = killGroup(cmd.Process.Pid) })`.
8. `runErr := procbound.Wait(cmd.Wait, time.Until(started.Add(p.cfg.ExecTimeout+2*pipeCloseGrace)), release)`.

Classification, in order:

| Condition | Result |
|---|---|
| overflow (either stream) | `fmt.Errorf("%w: %s exceeded %d bytes", ErrOutputLimit, stream, limit)`. Join `killFailure` only if `killErr != nil && !procbound.NothingToKill(killErr, func() bool { return procbound.GroupEmpty(pid) })`. Join `runErr` if `errors.Is(runErr, procbound.ErrCleanupIncomplete)`. The stdout bytes are NOT returned (nil). |
| `errors.Is(execCtx.Err(), context.DeadlineExceeded)` | `ErrExecTimeout` (same joins) |
| read error | wrapped |
| `runErr` is `*exec.ExitError` | IGNORED — exit codes are advisory (D4) |
| any other `runErr` | wrapped |
| otherwise | return the stdout bytes |

`readCapped` is capsule's verbatim, including the overflow decision `if int64(len(data)) > limit` ← **the M33 line**.

**Tests first** (`producer_internal_test.go`). Helpers:

- `scriptedProducer(t, body string, cfg ProducerConfig) (*Producer, dir string)` writes `dir/ailang` containing `#!/bin/sh\nif [ "$1" = "--version" ]; then echo "AILANG scripted"; exit 0; fi\n` + body, mode 0755. It pins `Compiler: hashref.SumSHA256(bytes)`, `CompilerVersion: "AILANG scripted"`, and constructs through `NewProducer` with an in-file `nopSeams` type (`ReadObject` → zero/nil/nil, `PutObject` → nil, `BusyTimeout` → 0).
- `swapKillGroup(t) *killRecorder` records pgids, forwards to the original, and restores on `t.Cleanup`.

Tests:

- **AC26 `TestProducerRefusesStderrOverflowIndependently`**
  - **Overflow arm:** body `printf '{}'; head -c 64 /dev/zero | tr '\0' e >&2; exec sleep 30`; cap 16; ExecTimeout 10 s. Requires `errors.Is(err, ErrOutputLimit)` and `!errors.Is(err, ErrExecTimeout)`, the detail contains `stderr`, the returned bytes are nil, elapsed < 5 s (the overflow kill ended the child), and the kill recorder count is exactly 1.
  - **Independence control:** the same body without `exec sleep 30`, cap 1024 → returns exactly `{}`, err nil.
- **AC25 `TestProducerExecTimeoutKillsTheWholeProcessGroup`**
  - Body `sleep 30 & echo $! > "$PIDDIR/grandchild.pid"; echo $$ > "$PIDDIR/child.pid"; wait`, with the env passed to `runBounded` carrying `PIDDIR=<t.TempDir()>` and `PATH`. ExecTimeout 300 ms, cap 1024.
  - Run `runBounded` in a goroutine; `select` on it with `time.After(20*time.Second)` → on the watchdog, `t.Fatalf("runBounded did not return within 20s; want ErrExecTimeout after ExecTimeout=300ms")`.
  - Read both pids with `proctest.ReadPid` and register `proctest.ReapOnCleanup` for both.
  - Assert: `errors.Is(err, ErrExecTimeout)`; not `procbound.ErrCleanupIncomplete`; the kill recorder saw exactly one pgid, equal to the child pid; `proctest.DeadWithin(t, grandchild, 2*time.Second)`; the first kill time is ≥ 300 ms after the call.
- **AC25 `TestProducerAdmissionRefusesBacklog`**
  - Fill `procbound` to `MaxOutstanding` (capsule's loop, releases deferred).
  - Body `touch "$PIDDIR/ran"`.
  - Assert `errors.Is(err, procbound.ErrCleanupBacklog)` and that `ran` does not exist.
  - After releasing, `procbound.Outstanding()` equals its pre-test value.

N=50; NEW = `TestProducerAdmissionRefusesBacklog|TestProducerExecTimeoutKillsTheWholeProcessGroup|TestProducerRefusesStderrOverflowIndependently`.

**Mutation drills:**

| Mutant | Edit | Must red |
|---|---|---|
| **M33** | `readCapped`: `if int64(len(data)) > limit` → `if false && int64(len(data)) > limit` | `TestProducerRefusesStderrOverflowIndependently`. Observed form: the returned error is `ErrExecTimeout` or the bytes `{}`, not `ErrOutputLimit`. |
| **M39** | `context.WithTimeout(ctx, p.cfg.ExecTimeout)` → `context.WithTimeout(ctx, time.Hour)` | `TestProducerExecTimeoutKillsTheWholeProcessGroup` (watchdog red; the cleanup reaps the pids) |
| M2-a | overflow kill body `killErr = killGroup(...)` → `_ = 0` (no kill) | stderr test elapsed < 5 s assertion |
| M2-b | `procbound.Admit()` → `func() (func(), error) { return func() {}, nil }()` | `TestProducerAdmissionRefusesBacklog` (`ran` exists / no error) |
| M2-c | `Setpgid: true` → `Setpgid: false` | exec-timeout test (grandchild survives, or `ErrCleanupIncomplete`) |

**Controller commit message:** `feat(evidence): bounded process-group runner with independent per-stream caps (row 26 M2)`

## C3 / M3 — Staging + `ai-check` invocation + committed fixtures (design M3 part b; AC26-stdout)

**Production ~45 LOC; tests ~80; fixtures generated.** Files:

- `proof_producer.go`, `producer_internal_test.go`
- NEW `host/evidence/testdata/aicheck/{verifiedpass,refuted,encerror,mixed,skippedshape,loaderror,nocontracts}.ail`
- the same names `.json`
- `regen.sh`
- the PD1 trio (N=51)

**Fixture generation (executor, never a test, never hand-typed):**

1. Extract each module mechanically from the committed generator:
   ```sh
   G=design_docs/verification/world-iter251/gen_aicheck_fixtures.sh; D=host/evidence/testdata/aicheck; mkdir -p "$D"
   for n in verifiedpass refuted encerror mixed skippedshape loaderror nocontracts; do
     awk -v n="$n" 'index($0,"fixtures/" n ".ail\" <<")>0{f=1;next} f&&/^AIL$/{f=0} f' "$G" > "$D/$n.ail"
   done
   wc -c "$D"/*.ail   # planner-measured: 229 210 281 349 531 151 249 (verifiedpass…nocontracts); 0 bytes = extraction failed
   ```
2. Write `host/evidence/testdata/aicheck/regen.sh` (bash, `set -euo pipefail`). It `cd`s to its own dir and uses `BIN="${AILANG_BIN:?AILANG_BIN must name the pinned v0.41.0 binary}"`. It asserts `"$BIN" --version 2>/dev/null | head -1` equals `AILANG v0.41.0`, else exit 1. It makes `scratch=$(mktemp -d)` with a trap to remove it and `mkdir "$scratch/fixtures"`. For each of the 7 names it copies `<n>.ail` to `$scratch/fixtures/`, then runs `( cd "$scratch" && "$BIN" ai-check -timeout 5s "fixtures/<n>.ail" ) > "<n>.json" 2>/dev/null || true`. The rc is data; the comment says so. It writes NO `.rc` or `.stderr` files.
3. Run `AILANG_BIN=$HOME/.pinned-ailang/ailang bash host/evidence/testdata/aicheck/regen.sh`.
4. Instrument check (planner-run: all 7 MATCH):
   ```sh
   for n in verifiedpass refuted encerror mixed skippedshape loaderror nocontracts; do
     a=$(jq -S -c 'del(.verify.results[]?.duration)' host/evidence/testdata/aicheck/$n.json)
     b=$(jq -S -c 'del(.verify.results[]?.duration)' design_docs/verification/world-iter251/fixtures/$n.json)
     [ "$a" = "$b" ] && echo "$n MATCH" || { echo "$n DIFF"; exit 1; }
   done
   ```

**Implementation, `runCheck`:**

1. `root, err := os.MkdirTemp("", "evidence-producer-*")`; `defer os.RemoveAll(root)`.
2. `mkdir root/staged`, `root/tmp`, `root/cache`.
3. Write the source to `root/staged/module.ail` (0o644).
4. Build `env := append(childenv.Scrubbed(os.Environ()), "AILANG_RELAX_MODULES=1", "TMPDIR="+filepath.Join(root,"tmp"), archive.CacheDirEnv+"="+filepath.Join(root,"cache"))`. The appended assignments come last, so an ambient value cannot win.
5. Return `p.runBounded(ctx, root, env, p.cfg.MaxOutputBytes, "ai-check", "-timeout", aiCheckZ3Timeout, stagedModule)`.

Importing `host/archive` is safe: no cycle, measured via `go list -deps`.

**Tests first.** Helper `pinnedChecker(t) CheckerConfig`:

- `bin := os.Getenv("AILANG_BIN")`; empty → `t.Fatal("AILANG_BIN unset: the producer's binary arms need the pinned v0.41.0 release; never skip")`.
- `out, err := exec.Command(bin, "--version").Output()`; the first line must be `AILANG v0.41.0`, else `t.Fatalf`.
- `data, _ := os.ReadFile(bin)` → `Compiler: hashref.SumSHA256(data)`.

**AC26 `TestProducerRefusesStdoutOverflowIndependently`**, three arms:

- **(a) real binary.** `t.Setenv("TMPDIR", t.TempDir())`. Producer with `pinnedChecker`, cap 16, ExecTimeout 2 min. `runCheck(ctx, verifiedpass.ail bytes)` → `errors.Is(err, ErrOutputLimit)`, nil bytes.
- **(b) control.** Same, cap `1<<20` → the JSON has `"twicePlusOne"` with `"status": "verified"`, and after return the `TMPDIR` dir has 0 entries (scratch root removed, child temp state confined).
- **(c) scripted stdout-only overflow.** Body `head -c 64 /dev/zero | tr '\0' o`, no stderr, cap 16 → `ErrOutputLimit`, and the detail contains `stdout`.

Planner measurement: the producer-staged verifiedpass run writes stdout 365 B and stderr 200 B (an `MOD010` relax warning). So arm (a) overflows both streams, and the per-stream independence proof for stdout is arm (c).

N=51; NEW = `TestProducerRefusesStdoutOverflowIndependently`.

**Mutation drills:**

| Mutant | Edit | Must red |
|---|---|---|
| **M33** | as in M2 | `TestProducerRefusesStdoutOverflowIndependently` arms (a) and (c). Design text predicted "373-byte"; the measured staged stdout is 365 B, and the failure observed is 17 bytes returned with err nil. |
| M3-a | `defer os.RemoveAll(root)` deleted | control arm (b): `TMPDIR` not empty |

**Not drillable, declared:**

- The `AILANG_RELAX_MODULES=1` and `TMPDIR` assignments. Measured: the pinned binary auto-relaxes `MOD010` in temp dirs and finds its solver with an inherited `TMPDIR`. So removing either does not change the verdict under tests that run inside temp dirs.
- Their rationale stays the design's V19/V26. Record "not drillable, measured" in the ledger, not "killed".

**Controller commit message:** `feat(evidence): stage and run pinned ai-check; commit pinned-binary fixtures + regen.sh (row 26 M3)`

## C4 / M4 — Pinned executable byte + version verification (design M2; AC24)

**Production ~45 LOC; tests ~110.** Files: `proof_producer.go`, `producer_internal_test.go`, authority_test.go (PD2), the PD1 trio (N=54).

**Implementation, `verifyTool(ctx)`:**

1. `data, err := os.ReadFile(p.checker.ExecPath)`; on error → `"%w: read executable %s: %v"` with `ErrUntrustedTool`.
2. `got, err := hashref.Sum(p.checker.Compiler.Algo(), data)`; on error → `ErrUntrustedTool` wrap.
3. **The M35 line:** `if got.Digest() != p.checker.Compiler.Digest()` → `"%w: executable %s bytes hash to %s; want %s"`.
4. Only then probe: `out, err := p.runBounded(ctx, "", childenv.Scrubbed(os.Environ()), versionProbeMaxBytes, "--version")`. This captures stdout only (no `2>&1` anywhere), and both streams are capped. On error → `"%w: version probe: %v"` with `ErrUntrustedTool`, wrapping err with `%w` as well so `ErrExecTimeout`/`ErrOutputLimit` stay reachable.
5. `line := strings.TrimRight(strings.SplitN(string(out), "\n", 2)[0], "\r")`.
6. **The M36 line:** `if line != p.checker.CompilerVersion` → `"%w: version line %q; want %q"`.

`ReadFile` was chosen per the design's Deferred list (capsule shape, a bounded 100,594,162-byte transient).

**Tests first:**

- **`TestProducerRefutesTamperedExecutableBytes`**
  - (a) Real binary path with `Compiler: hashref.SumSHA256([]byte("not the pinned binary"))` → `ErrUntrustedTool`, detail contains `hash to` and the expected digest.
  - (b) Scripted body whose `--version` branch is preceded by `touch "$(dirname "$0")/executed"`. Note: the `touch` must sit BEFORE the version branch so ANY execution marks it. It is pinned to a WRONG ref → `ErrUntrustedTool`, and `executed` is absent. This proves bytes are checked before any execution.
- **`TestProducerRefutesWrongCompilerVersion`** — real binary, correct `Compiler`, `CompilerVersion: "AILANG v9.9.9"` → `ErrUntrustedTool`, detail contains `"AILANG v0.41.0"` and `"AILANG v9.9.9"`. Control arm: correct version → nil.
- **`TestProducerRefusesMissingExecutable`** — `ExecPath: filepath.Join(t.TempDir(), "absent")` → `ErrUntrustedTool`, detail contains `read executable`.

N=54; NEW = `TestProducerRefusesMissingExecutable|TestProducerRefutesTamperedExecutableBytes|TestProducerRefutesWrongCompilerVersion`.

**Mutation drills:**

| Mutant | Edit | Must red |
|---|---|---|
| **M35** | `if got.Digest() != p.checker.Compiler.Digest()` → `if false && got.Digest() != …` | `TestProducerRefutesTamperedExecutableBytes` arm (a): got nil, want `ErrUntrustedTool` |
| **M36** | `if line != p.checker.CompilerVersion` → `if false && line != …` | `TestProducerRefutesWrongCompilerVersion` |
| M4-a | move the byte check below the probe (order swap) | tampered arm (b): `executed` exists |
| M4-b | ReadFile error branch → `if false && err != nil` | `TestProducerRefusesMissingExecutable` detail `read executable` |

**Controller commit message:** `feat(evidence): per-run pinned executable byte check before a bounded --version probe (row 26 M4)`

## C5 / M5 — Strict `ai-check` interpretation + fresh-binary revalidation (design M4 part a; AC27)

**Production ~70 LOC; tests ~160.** Files: `proof_producer.go`, `producer_internal_test.go`, authority_test.go (PD2), the PD1 trio (N=56).

**Implementation, `interpretCheck`.** The wire struct is unexported: `check.passed bool`; `verify.available bool`, `counterexample int`, `skipped int`, `errors int`, and `results []{function string, status string}`. **There is deliberately NO field for the `verify.verified` counter** (D4 step 8). Parse with `encoding/json.Unmarshal`; unknown fields are tolerated. Fixed order, one `if` per step, these exact detail fragments:

1. `len(bytes.TrimSpace(stdout)) == 0` → `ErrVacuousCheck` `"empty ai-check stdout"`;
2. unmarshal error → `ErrUnverifiedRun` `"unparseable ai-check JSON: %v"`;
3. `!d.Check.Passed` → `ErrUnverifiedRun` `"check.passed is not true"`;
4. `!d.Verify.Available` → `ErrUnverifiedRun` `"verify.available is not true"`;
5. `d.Verify.Errors > 0` → `ErrUnverifiedRun` `"verify.errors=%d"` ← **M38**;
6. `d.Verify.Counterexample > 0` → `ErrUnverifiedRun` `"verify.counterexample=%d"`;
7. `d.Verify.Skipped > 0` → `ErrUnverifiedRun` `"verify.skipped=%d"`;
8. `verified` = `r.Function` for each record with `r.Status == "verified"`;
9. ONE guard, `if len(verified) == 0` → `ErrVacuousCheck`, with detail `"0 results records from non-empty output"` when `len(results)==0`, else `"0 verified among %d results records"` ← **M34**;
10. `sort.Strings(verified)`; return it. No dedupe: the codec refuses duplicates loudly at mint.

**Tests first:**

**`TestProducerInterpretationMatchesEveryFixtureShape`** is a table over the committed fixtures (`os.ReadFile("testdata/aicheck/<n>.json")`), each asserting sentinel + detail fragment (+ verified set):

| Fixture | Expected |
|---|---|
| verifiedpass | `["twicePlusOne"]`, nil |
| refuted | `ErrUnverifiedRun`, `verify.counterexample=1` |
| encerror | `ErrUnverifiedRun`, `verify.errors=1` |
| mixed | `ErrUnverifiedRun`, `verify.errors=1` |
| skippedshape | `ErrUnverifiedRun`, `verify.skipped=1` |
| loaderror | `ErrUnverifiedRun`, `check.passed is not true` |
| nocontracts | `ErrVacuousCheck`, `0 results records from non-empty output` |
| literal `""` | `ErrVacuousCheck`, `empty ai-check stdout` |
| literal `"{"` | `ErrUnverifiedRun`, `unparseable` |

PD6 derived arms from verifiedpass:

| Flipped field | Expected |
|---|---|
| `check.passed=false` | `ErrUnverifiedRun`, `check.passed` |
| `verify.available=false` | `ErrUnverifiedRun`, `verify.available` |
| `results[0].status="unknown"` | `ErrVacuousCheck`, `0 verified among 1 results records` |

**`TestFixturesRevalidateAgainstFreshBinaryOutput`** (`pinnedChecker`, `t.Fatal` on absence):

- For each of the 7 committed `.ail`: run `p.runCheck(ctx, src)` with cap 1<<20 and ExecTimeout 2 min.
- Fresh verdict = `interpretCheck(fresh)`; committed verdict = `interpretCheck(committed .json)`.
- Require the same sentinel identity (or both nil) and the same verified set.
- A pin move that changes JSON shape reds here. The fixtures are never trusted as frozen truth.

N=56; NEW = `TestFixturesRevalidateAgainstFreshBinaryOutput|TestProducerInterpretationMatchesEveryFixtureShape`.

**Mutation drills** (all killed by `TestProducerInterpretationMatchesEveryFixtureShape` unless stated):

| Mutant | Edit | Must red |
|---|---|---|
| **M34** | `if len(verified) == 0` → `if false && len(verified) == 0` | nocontracts arm (verified=[] with nil err, not `ErrVacuousCheck`) |
| **M38** | `if d.Verify.Errors > 0` → `if false && d.Verify.Errors > 0` | mixed arm (mints `["good"]`), encerror arm (`ErrVacuousCheck`) |
| M5-a | counterexample branch neutered | refuted arm |
| M5-b | skipped branch neutered | skippedshape arm |
| M5-c | check.passed branch neutered | loaderror detail + derived passed arm |
| M5-d | verify.available branch neutered | derived available arm |
| M5-e | `r.Status == "verified"` → `true` | derived status arm |
| M5-f | empty-stdout branch neutered | `""` arm (becomes `ErrUnverifiedRun`) |
| M5-g | fixture tamper: edit committed `verifiedpass.json` `"status": "verified"` → `"error"` (restore after) | `TestFixturesRevalidateAgainstFreshBinaryOutput` |

**Controller commit message:** `feat(evidence): strict ai-check JSON interpretation with one anti-vacuity guard (row 26 M5)`

## C6 / M6 — Required identities + MAC-tagged canonical mint (design M4 part b; AC28, AC29 size arm)

**Production ~50 LOC; tests ~120.** Files: `proof_producer.go`, `producer_internal_test.go`, authority_test.go (PD2), the PD1 trio (N=58).

**Implementation:**

`requireIdentities(verified, required)`:

- `len(required) == 0` → `ErrRequiredIdentityMissing` `"no required identities named"` (PD10);
- for each `id`: **the M15 line** `if !containsIdentity(verified, id)` → `"%w: required identity %q not verified"`.

`mint(subject, verified)`:

1. `report := ProofReportV1{Schema: ProofReportSchemaV1, Subject: subject, Compiler: p.checker.Compiler, CompilerVersion: p.checker.CompilerVersion, Verified: verified, Errors: 0, Counterexamples: 0, CheckPassed: true, ProofSucceeded: true}`.
2. `rb, err := EncodeProofReportV1(report)`.
3. `mac := hmac.New(sha256.New, p.key[:]); mac.Write(rb)`.
4. `env, err := EncodeAuthenticatedEnvelope(AuthenticatedEnvelope{Report: rb, MAC: mac.Sum(nil), MACValid: true})`.
5. **The M37 line:** `if len(env) > MaxBytes` → `"%w: envelope is %d bytes; validator read cap is %d"` with `ErrEnvelopeTooLarge`.
6. Return `env, hashref.SumSHA256(env), nil`.

**Tests first:**

- **AC28 `TestProducerRequiresCallerNamedIdentities`** (verified set from `interpretCheck(verifiedpass.json)`):
  - **control:** `["twicePlusOne"]` → nil. Then `p.mint(testSubjectRef, verified)` succeeds; `DecodeAuthenticatedEnvelope(env)` succeeds; `hmac.Equal(HMAC(key, env.Report), env.MAC)`; the decoded report's `Verified == ["twicePlusOne"]`, `Subject`/`Compiler`/`CompilerVersion` as configured; `ref == hashref.SumSHA256(env)`.
  - **absent:** `["absentIdentity","twicePlusOne"]` → `ErrRequiredIdentityMissing`, detail contains `absentIdentity`.
  - **empty:** `nil` and `[]string{}` → `ErrRequiredIdentityMissing`, detail `no required identities named`.
- **AC29 `TestProducerRefusesEnvelopeBeyondValidatorReadCap`** (PD4):
  - over arm 200×1000 → `ErrEnvelopeTooLarge`, nil bytes, zero ref;
  - control 190×1024 mints, with `len(env) <= MaxBytes` and `len(env) > MaxBytes-4096`. The near-bound control proves the stimulus straddles the cap.

N=58; NEW = `TestProducerRefusesEnvelopeBeyondValidatorReadCap|TestProducerRequiresCallerNamedIdentities`.

**Mutation drills:**

| Mutant | Edit | Must red |
|---|---|---|
| **M15** | `if !containsIdentity(verified, id)` → `if false && !containsIdentity(verified, id)` | `TestProducerRequiresCallerNamedIdentities` absent arm |
| **M37** | `if len(env) > MaxBytes` → `if false && len(env) > MaxBytes` | `TestProducerRefusesEnvelopeBeyondValidatorReadCap` over arm |
| M6-a | empty-required branch neutered | empty arm |
| M6-b | `mac.Write(rb)` → `mac.Write(rb[:len(rb)-1])` | control arm HMAC assertion |

**Controller commit message:** `feat(evidence): required identities and MAC-tagged canonical envelope mint bounded by the validator read cap (row 26 M6)`

## C7 / M7 — Store legs + `GenerateProof` (design M5 part a; AC5 e2e, AC29 write/readability)

**Production ~70 LOC; tests ~200.** Files: `proof_producer.go`, `producer_realstore_test.go`, authority_test.go (PD2), the PD1 trio (N=61).

**Implementation:**

`readSource`:

1. `if _, err := hashref.Parse(ref.String()); err != nil` → wrapped error `"evidence: invalid source ref"`.
2. `readCtx, cancel := context.WithTimeout(ctx, p.cfg.ObjectReadTimeout)` ← **M31 line**; `defer cancel()`.
3. `_, src, err := p.reader.ReadObject(readCtx, ref, MaxSourceBytes)`; on error → `fmt.Errorf("evidence: read source %s: %w", ref, err)`.
4. `src == nil` → `"evidence: source object %s not found"`.
5. `got, err := hashref.Sum(ref.Algo(), src)`; `got.Digest() != ref.Digest()` → `"evidence: source %s bytes hash to %s"`.
6. Return `src`.

`GenerateProof`:

1. `src` ← `readSource`.
2. `p.verifyTool(ctx)`.
3. `stdout` ← `p.runCheck(ctx, src)`.
4. `verified` ← `interpretCheck(stdout)`.
5. `requireIdentities(verified, requiredIdentities)`.
6. `env, ref` ← `p.mint(sourceRef, verified)`.
7. `writeCtx, cancel := context.WithTimeout(ctx, p.cfg.ObjectWriteTimeout)` ← **M32 line**.
8. **The M28 line:** `if err := p.writer.PutObject(writeCtx, store.Object{Hash: ref, InterfaceHash: InterfaceHashV1, SemanticID: ProofSemanticID, Provenance: "bounded-ai-check-producer", Payload: env}); err != nil` → `return hashref.HashRef{}, fmt.Errorf("evidence: write proof envelope: %w", err)` (an `*store.UncertainError` stays `errors.As`-reachable; never retried).
9. Return `ref, nil`.

Every failure path returns the zero `HashRef`.

**Tests first** (`producer_realstore_test.go`). Helpers:

- `pinnedChecker(t)`, duplicated in this external-package file (~15 lines).
- `countingWriter{s *store.Store; calls atomic.Int32; fail error}` — `PutObject` increments, returns `fail` if non-nil WITHOUT calling the store, else forwards; `BusyTimeout` forwards.
- `putSource(t, s, src) hashref.HashRef` — `store.Object{Hash: SumSHA256(src), InterfaceHash: hashref.SumSHA256([]byte("test/ail-source")), SemanticID: "test/ail-source", Provenance: "producer-test", Payload: src}`.

Tests:

- **AC5 `TestProducerEndToEndMintsEnvelopeValidatorAccepts`** — default-window store `s` (2 s), cfg `{3s, 3s, 2*time.Minute, 1<<20}`, source = `testdata/aicheck/verifiedpass.ail`. Subtests run in order:
  - **mint:** `GenerateProof(context.Background(), srcRef, ["twicePlusOne"])` → ref, err nil, `calls==1`. Then `v := NewValidator(testKey, s, CompilerConfig{Compiler: pinned.Compiler, CompilerVersion: "AILANG v0.41.0", ObjectReadTimeout: 3s}, ["twicePlusOne"])`; `v.ValidateProof(boundedTestContext(t), ref, srcRef).Validated()` → ok; `v.Resolve(sealed).Proven()` → `(ResolvedGradeProven, true)`.
  - **idempotent:** a second call returns the identical ref, err nil (A1).
  - **tampered-compiler:** `Compiler = SumSHA256("not the pinned binary")` → `ErrUntrustedTool`, zero ref, calls 0.
  - **absent-identity:** `["absentIdentity"]` → `ErrRequiredIdentityMissing`, zero, calls 0.
  - **source-hash-mismatch:** reader wrapper `tamperReader{s}` returns the real meta with the payload's last byte flipped → error contains `hash to`, zero, calls 0.
  - **missing-source:** ref of never-stored bytes → error contains `not found`, zero, calls 0.
  - **oversize-source:** store a `MaxSourceBytes+1` object → `errors.As(err, new(*store.ObjectTooLargeError))`, zero, calls 0.
  - **write-error:** `countingWriter{fail: errors.New("injected not-committed")}` → `errors.Is` the injected error, zero ref.
- **AC29 `TestProducerMintsValidatorReadableEnvelope`** — after a mint:
  - `s.ReadObject(boundedTestContext(t), ref, evidence.MaxBytes)` → `meta.SemanticID == evidence.ProofSemanticID`, `meta.InterfaceHash == evidence.InterfaceHashV1`, `hashref.SumSHA256(payload) == ref`;
  - `DecodeAuthenticatedEnvelope` + `DecodeProofReportV1(env.Report)` → `Subject == srcRef`, `Compiler == pinned`, `CompilerVersion == "AILANG v0.41.0"`, `Verified == ["twicePlusOne"]`, `Errors == 0`, `Counterexamples == 0`, `CheckPassed`, `ProofSucceeded`;
  - HMAC(testKey, env.Report) equals env.MAC; `calls == 1`.
- **AC29 `TestProducerSurfacesUncertainWriteOutcome`** (PD3) — `uncertainWriter` → `errors.As(err, new(*store.UncertainError))`, zero ref, `calls == 1` (no blind retry). The object DOES exist in the store afterwards (a real commit happened). Assert that through `s.GetObject`, which is the reconcile path the caller owns.

N=61; NEW = `TestProducerEndToEndMintsEnvelopeValidatorAccepts|TestProducerMintsValidatorReadableEnvelope|TestProducerSurfacesUncertainWriteOutcome`.

**Mutation drills:**

| Mutant | Edit | Must red |
|---|---|---|
| **M28** | `…PutObject(writeCtx, …); err != nil` → `…; false && err != nil` | e2e write-error arm (non-zero ref returned) and `TestProducerSurfacesUncertainWriteOutcome` |
| M7-a | delete the `p.verifyTool(ctx)` call in `GenerateProof` | e2e tampered-compiler arm (mints) |
| M7-b | source digest compare → `false && …` | e2e source-hash-mismatch arm |
| M7-c | `MaxSourceBytes` in the `ReadObject` call → `MaxSourceBytes<<10` | e2e oversize-source arm |
| M7-d | delete the `requireIdentities` call | e2e absent-identity arm |
| M7-e | `SemanticID: ProofSemanticID` → `"x"` | e2e mint arm (validator refuses wrong_semantic_id) and `TestProducerMintsValidatorReadableEnvelope` |

**Controller commit message:** `feat(evidence): GenerateProof composes derived-deadline source read, pinned run, strict mint and bounded envelope write (row 26 M7)`

## C8 / M8 — Real-store pool decoys + full re-drill + final gates (design M5 part b; AC20, AC23)

**Production 0 LOC; tests ~190.** Files: `producer_realstore_test.go`, the PD1 trio (N=63).

**Tests first** (PD7; both copy the `GOMAXPROCS(2)` precondition block from realstore_test.go:145-160 verbatim; store = `openEvidenceStoreWithoutBusyWait`):

- **AC23 `TestRealStoreBlockedSourceReadReturnsWithinObjectReadTimeout`**
  1. Source = verifiedpass. Decoy = `putEvidenceObject(t, s, []byte(strings.Repeat("d", 256<<20)))`. `hold := measureGetObject(...)`; `hold < minDecoyHold` → instrument failure.
  2. `bound := hold/decoyHoldRatio`; `settle := hold/40`.
  3. Producer: `ObjectReadTimeout: bound`, `ObjectWriteTimeout: 3*time.Second + 2*hold`, `ExecTimeout: 2*time.Minute`, `MaxOutputBytes: 1<<20`, reader = `decoyReader{s, armed}`, writer = a fresh `countingWriter{s}` per attempt. A refusal additionally requires `calls == 0`.
  4. **No-decoy control first:** unarmed `GenerateProof(context.Background(), …)` mints.
  5. Then up to 5 armed attempts with the PD7 classification. The refusal must satisfy `errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, store.ErrNoDeadline)` with a zero ref.
- **AC20 `TestRealStoreBlockedEnvelopeWriteReturnsWithinObjectWriteTimeout`**
  1. Same setup, but `ObjectWriteTimeout: bound`, `ObjectReadTimeout: 3*time.Second + 2*hold`, writer = `decoyWriter{s, armed}`.
  2. **No-decoy control:** unarmed mint.
  3. Armed attempts use a SECOND source: `verifiedpass.ail` bytes + `"\n"`, which gives a different `sourceRef` and the same verdict. Same inputs mint byte-identical envelopes, so the control already stored the first source's envelope. `decoyWriter` records the `o.Hash` it was handed on the armed call (the would-be ref). A refusal requires `errors.Is(err, context.DeadlineExceeded)`, `!errors.Is(err, store.ErrNoDeadline)`, a zero returned ref, and `s.GetObject(boundedTestContext(t), wouldBe)` → `ok == false` (no object written). A won-connection attempt (mint with store call < 2×bound) leaves that object stored, so each retry uses a further distinct source: append one more `"\n"` per attempt.
  4. The criterion text "holds NO SQLite write lock; claims nothing about lock contention" goes into the test's doc comment verbatim.

N=63; NEW = `TestRealStoreBlockedEnvelopeWriteReturnsWithinObjectWriteTimeout|TestRealStoreBlockedSourceReadReturnsWithinObjectReadTimeout`.

**Mutation drills:**

| Mutant | Edit | Must red |
|---|---|---|
| **M31** | `context.WithTimeout(ctx, p.cfg.ObjectReadTimeout)` → `context.WithTimeout(ctx, time.Hour)` | AC23: `blocked source read minted after exceeding ObjectReadTimeout: store call=…` |
| **M32** | `context.WithTimeout(ctx, p.cfg.ObjectWriteTimeout)` → `context.WithTimeout(ctx, time.Hour)` | AC20: `blocked envelope write stored after exceeding ObjectWriteTimeout: store call=…` |

**Full re-drill at M8:** M15, M28, M31, M32, M33, M34, M35, M36, M37, M38, M39 once more on the final tree, each with its own milestone's command.

**Final gates** (in order; record rc and wall time; run verify_go with ≥ 15 min allowance):

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang; export PATH=/opt/homebrew/bin:$PATH
go vet ./...
go test -count=1 ./host/verifygate                         # full, ~158 s locally (PD1 edits exercised end to end)
go test -race ./host/evidence -count=1
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh   # must stay byte-for-byte unaffected; no .ail moved in its ROOTS
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh    # banner "(63 exact top-level tests)", "✓ all 63 required…", "✓ go gate PASSED"
git diff --stat ba79d1e -- host/store scripts/verify_ail.sh .github tools/launchd world packages   # controller: must be empty
```

**Controller commit message:** `test(evidence): real-store pool decoys bound both producer store legs; manifest 63 (row 26 M8)`

---

## Ordered gate list (whole sprint)

1. M1: `go vet ./host/evidence/... ./host/verifygate/...` → named `-v` run with the PASS loop (10 tests) → `go test ./host/evidence -count=1` → `go test -race ./host/evidence -count=1` → `--evidence-manifest-check … 47` → focused verifygate (4 tests).
2. M2: same template, N=50.
3. M3: regen.sh + jq MATCH check, then the template, N=51.
4. M4: template, N=54.
5. M5: template, N=56.
6. M6: template, N=58.
7. M7: template, N=61.
8. M8: template, N=63, then full re-drill, then the final gates: `go vet ./...` → `go test -count=1 ./host/verifygate` → `go test -race ./host/evidence -count=1` → `verify_ail.sh` → `verify_go.sh` → the frozen-surface diffstat.

## Risks

- **Runtime.** CI's go-verify leg is near its budget: verify_go.sh took 10 m 18 s locally at base, and the race leg has an 8 m inner timeout. The new tests add:
  - two more 256 MiB decoys;
  - about 12 real-binary runs, each hashing 100,594,162 bytes.

  The executor reports `go test -race ./host/evidence` wall time at M7 and M8. If it exceeds ~90 s locally, that is a controller finding, not something the executor fixes by weakening a test.
- **Decoy flake.** PD7 copies the landed AC16 discipline and fixes the measured interval. Retry exhaustion is a loud instrument failure, never a pass.
- **ETXTBSY.** Scripted executables are written then exec'd. host/evidence has no `t.Parallel` (PD12), which keeps concurrent forks off the write window (capsule precedent).
- **Fixture drift.** A pin move reds `TestFixturesRevalidateAgainstFreshBinaryOutput`. Regenerate with regen.sh in the same commit as the pin bump.
- **Over-budget milestone.** If any milestone exceeds ~250 production LOC, the executor stops and reports actual LOC. It never hides overflow in a neighbour.

## Coverage census

- 8 milestones / commits C1–C8.
- Design ACs covered: AC5 (both arms), AC20, AC23, AC24, AC25, AC26, AC27, AC28, AC29, AC30. That is 10 criteria, 26 named top-level tests, manifest 37 → 63.
- Design mutants: M15, M28, M31–M39 (11), each mapped to a milestone and named test, plus 33 planner drills (M1-a…m, M2-a…c, M3-a, M4-a…b, M5-a…g, M6-a…b, M7-a…e).
- 12 planner decisions; 9 design defects recorded.
- Nothing deferred.

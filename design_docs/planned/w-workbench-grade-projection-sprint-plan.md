# Sprint plan — `w-workbench-grade-projection`

**Iteration 254 · queue row 105 · clause-5 · v1.0.0 release prerequisite (D-WORLD-69 = A)**

Design: [w-workbench-grade-projection.md](w-workbench-grade-projection.md). Quorum r1 BLOCKED (two measured objections, applied), r2 BLOCKED → narrow-refinement carve-out with the reviewers' verbatim fixes. Worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter253-row105-design`, branch `design/w-workbench-grade-projection`, planned at HEAD `b8581e0`. The code base is `edb1b63` plus design-doc commits only. Binding: CLAUDE.md and design_docs/coding-standards.md (S2 effects at the host boundary, S6 non-vacuous gates and load-bearing criteria need mutations, S7 usage surfaces ship usage docs). **~1.6 executor days, 5 milestones (the design's 3, re-cut so each is ≤ ~250 changed LOC — PD9).**

**Plan only.** The planner writes exactly two files: this one and `sprint_w-workbench-grade-projection.json`. **The EXECUTOR runs no git command** (no add, commit, stash, checkout, branch, or `git diff` writes; read-only `git diff --stat` / `git status` are allowed). The executor implements C1→C5 in order in this worktree. Every gate in a milestone's list must be green at its own boundary before the next milestone starts. At each boundary the executor writes a boundary report: changed files with LOC, the rc and wall time of each gate, the red-first evidence, the mutation ledger and restore evidence. **The controller commits** (one commit per boundary if it told the executor to stop at boundaries, otherwise from the final tree). Every boundary is internally consistent, so per-milestone commits are bisectable (PD3).

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang   # AILANG v0.41.0 — the ONLY binary; package host/daemon needs it (PD5)
export PATH=/opt/homebrew/bin:$PATH
```

Run gate scripts under **bash**. The Bash tool's default shell is zsh, and zsh `for t in $X` does not split words.

---

## Baseline (planner-measured at `b8581e0`, first-party, 2026-10-12)

Logs are under `~/.ailang/state/iter254-plan/`.

| Gate | Result |
|---|---|
| `go vet ./...` | rc=0 |
| `go test ./host/daemon ./cmd/ailang-worldd -count=1` **without** `AILANG_BIN` | **rc=1**. The only failure is `TestDaemonStartReportsMissingTemplates`, which says "AILANG_BIN is unset: this test runs the pinned interpreter's plan phase and never skips". This is an environment requirement, not a base defect (PD5). `cmd/ailang-worldd` ok 3.0 s. |
| same, **with** `AILANG_BIN=$HOME/.pinned-ailang/ailang` | **rc=0** (160 s) — see the addendum |
| `go test -race ./host/daemon -count=1` without `AILANG_BIN` | rc=1, the same single test |
| `AILANG_BIN=… ./scripts/verify_ail.sh` | **rc=0**: `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 444 se-tools named tests pass` |
| `AILANG_BIN=… ./scripts/verify_go.sh` | **rc=1**. The non-race leg passed every package (host/daemon ok 210 s, cmd/ailang-worldd ok 13 s). The `-race -timeout 8m` leg failed ONLY in `host/verifygate`: `panic: test timed out after 8m0s` inside `TestSeToolsGateArmsAreLoadBearing` (module_manifest_gate_test.go:472 → runGateAt). The machine load average was 13 during the run. host/daemon `-race` was ok in 358 s. See the addendum for the isolated rerun. |

**Load note (affects PD8).** Under verify_go's `-race -timeout 8m`, host/daemon already takes 358 s of its 480 s budget. All new tests in this sprint together must add **< 15 s** to `go test -race ./host/daemon`. The executor records the wall time at C4 and C5.

---

## Design defects found by the planner (each measured first-party; each resolved by a PD)

1. **The seam cannot fabricate validator results (AC4 mapping table, Testing Strategy).** The design's table test drives `objectGradeView` "through a fake `proofValidator` for all 12 `UnsupportedReason` tokens plus one unknown token". But `evidence.ValidationResult{validated, hasValue, unsupported, err}` and `evidence.ResolutionResult{grade, ok, err}` have ONLY unexported fields (`sed -n '96,115p;200,209p' host/evidence/validator.go`). Code in package `daemon` therefore cannot build a result that carries a reason token, a validated seal or an error. It can build only the zero value. An unknown token can never come out of a `ValidationResult`. → **PD1**.
2. **The reason sentences contain apostrophes ("daemon's"), and `html/template` renders `'` as `&#39;`.** A test that searches the body for the raw sentence would red on a correct program. Measured with a scratch `go run` of the same `{{.}}` text context: the rendered line equals `"<p>GRADE UNAVAILABLE — " + html.EscapeString(sentence) + "</p>"` for all three probed sentences (`true` ×3). → **PD2**.
3. **The design leaves identities-without-key silently unconfigured.** D2 says "the daemon never degrades to a half-configured validator" and lists the half-set arms (key without `--ailang-bin`, key without identities). But D6's code only runs inside `if cfg.ProofKeyFile != ""`. So `--proof-required-identity X` with no key file starts serving the constant, and the operator's flag is silently ignored. → **PD4** (a seventh refusal arm).
4. **S7 + AC7's "exactly six files".** Two new user-facing serve flags need usage docs beyond `--help` (coding-standards S7). The design touches no `docs/QUICKSTART.md`. There is also a trap here: `section11ServeLines` (cmd/ailang-worldd/quickstart_layout_test.go:28) cuts from `### 11. Projects with a module root…` to **end of file** with no terminating cut, and `TestQuickstartSection11FlagsMatchTheCLI` requires exactly ONE `/tmp/ailang-worldd serve ` line. The one that exists today is line 969 (`awk` measurement in PD7). So any `/tmp/ailang-worldd serve` line added to "## 11. Open the live workbench" (line 1061, after it) reds that test. → **PD7**.
5. **The design's M1 exit gate omits `AILANG_BIN` and uses `go build ./...` as a compile fence.** host/daemon cannot pass without `AILANG_BIN` (Baseline). `go build` does not compile `_test.go`. → **PD5** (every gate exports AILANG_BIN; the fence is `go vet ./...`).
6. **`envelopeFor` is unexported test code in package `evidence_test`** (`host/evidence/authority_test.go:48`), so package `daemon` cannot import it. This is not a design error (the design says "mirrors"), but the executor must re-implement the recipe with the EXPORTED codecs: `evidence.EncodeProofReportV1`, `crypto/hmac`+`sha256`, `evidence.EncodeAuthenticatedEnvelope(evidence.AuthenticatedEnvelope{Report, MAC, MACValid})`. All three are exported, and so are the fields. → **PD6**.

Premises verified as stated at HEAD (no defect): `daemon.go:688 d.interpreterVersion = m.Version` (and :687 `d.arch, d.interpreterHash = a, ref`); validator API at `validator.go:110/111/114/120/208/209/211`; `NewValidator` refusal order (`:63-93`): reader → timeout → compiler → version → identities → window; `MaxSemanticIDPage = 500` (objects_by_semantic_id.go:13), with the limit check at :34 admitting 101; `checkReadOrdering(s.BusyTimeout(), readDeadline)` at daemon.go:643, `readDeadline = 10 * time.Second` at :144; `fakeInterpreter(t, dir, version)` at daemon_test.go:42 (`printf '%s' '<version>'`). The discriminating fixture `"AILANG v0.41.0 \r\nbuild abc\n"` survives the shell verbatim: a scratch script printed bytes `… 0 sp \r \n b u i l d sp a b c \n`. `TestWorkbenchObjectGrade` is at workbench_test.go:974. The construction site is `workbench.go:477`. `HashRef{}` → `ValidateProof` returns `invalid_ref` because `HashRef{}.String()==""` and `Parse("")` refuses (hashref.go:79-83, :127-129). `PutObject` has no payload size cap (store.go:543-575), so a >256 KiB candidate is storable (it is the oversize arm). `EncodeProofReportV1` accepts a trailing-space `CompilerVersion` (only `checkString` ≤1024 bytes, json.go:78) and requires `Verified` strictly sorted and unique (report_codec.go:86-96).

---

## Planner decisions

**PD1 — the seam is a SCRIPTED GRADER that composes REAL validators. Unknown tokens are tested on a pure function.** Here `proofValidator` is the design's seam:

```go
type proofValidator interface {
	ValidateProof(ctx context.Context, reportRef, expectedSubject hashref.HashRef) evidence.ValidationResult
	Resolve(sealed evidence.ValidatedEvidence) evidence.ResolutionResult
}
```

`*evidence.Validator` satisfies it. The test double (test file only) is:

```go
// scriptedGrader returns ONLY results a real *evidence.Validator produced, or the zero value.
type scriptedGrader struct {
	v        *evidence.Validator // validates (built over a scriptedReader)
	ref      *hashref.HashRef    // when non-nil, replaces the candidate ref passed to v
	resolver proofValidator      // when non-nil, Resolve is delegated here (a DIFFERENT real validator = foreign seal)
	zero     bool                // when true, ValidateProof returns evidence.ValidationResult{} (no verdict)
}
type scriptedReader struct {
	meta    store.ObjectMeta
	payload []byte
	err     error
}
func (r scriptedReader) ReadObject(context.Context, hashref.HashRef, int64) (store.ObjectMeta, []byte, error) { return r.meta, r.payload, r.err }
func (r scriptedReader) BusyTimeout() time.Duration { return 0 }
```

Every one of the 12 tokens is produced by a real validator (the table is in C5). The unknown-token arm runs on the production pure function `gradeReasonSentence(evidence.UnsupportedReason("future_reason"))`, because `UnsupportedReason` is an exported string type. The zero-value result gets its own production branch (PD10). Nothing fakes a seal.

**PD2 — every sentence assertion goes through the HTML escaper.** Test helper: `func wantGradeUnavailable(reason string) string { return "<p>GRADE UNAVAILABLE — " + html.EscapeString(reason) + "</p>" }` (import `html`). The PROVEN line is the literal `<p><span>PROVEN</span></p>`.

**PD3 — bisectability.** Every boundary C1…C5 leaves `go vet ./...`, `go test ./host/daemon ./cmd/ailang-worldd -count=1` and `go test -race ./host/daemon -count=1` green (with `AILANG_BIN`). Nothing is deferred across a boundary. No milestone adds an unused import, no milestone leaves a failing test, and no gate counts move.

**PD4 — the identities-without-key refusal.** `configureGrader` first refuses `cfg.ProofKeyFile == "" && len(cfg.ProofRequiredIdentities) > 0` at `StageConfig` with detail `--proof-required-identity needs --proof-key-file`. This is the seventh arm of `TestProofKeyWiringRefusesEveryHalfConfiguration` (planner mutant M13). The CLI also refuses an empty identity value (`is empty`), mirroring `run-allow-caps`.

**PD5 — gate environment.** Every gate command runs with `AILANG_BIN=$HOME/.pinned-ailang/ailang` exported. The compile fence is `go vet ./...`, never `go build ./...`. `verify_go.sh` needs ≥ 15 min of allowance (its race leg alone exceeds 8 min under load). The executor runs it once at C5 with a 20-minute timeout. If the only red is the `host/verifygate` `-race` 8-minute timeout seen at baseline, the executor records it as BASE (cite the Baseline row and the addendum) and does not "fix" it.

**PD6 — envelope minting in package daemon uses only exported codecs** (defect 6):

```go
func mintGradeEnvelope(t *testing.T, key [32]byte, report evidence.ProofReportV1) []byte {
	t.Helper()
	raw, err := evidence.EncodeProofReportV1(report)
	if err != nil { t.Fatalf("EncodeProofReportV1: %v", err) }
	m := hmac.New(sha256.New, key[:])
	_, _ = m.Write(raw)
	env, err := evidence.EncodeAuthenticatedEnvelope(evidence.AuthenticatedEnvelope{Report: raw, MAC: m.Sum(nil), MACValid: true})
	if err != nil { t.Fatalf("EncodeAuthenticatedEnvelope: %v", err) }
	return env
}
```

The envelope is stored with the producer's identities (proof_producer.go:495-497): `store.Object{Hash: hashref.SumSHA256(env), InterfaceHash: evidence.InterfaceHashV1, SemanticID: evidence.ProofSemanticID, Provenance: "bounded-ai-check-producer", Payload: env}`.

**PD7 — QUICKSTART (S7) lands at C5, and it contains no `/tmp/ailang-worldd serve` line** (defect 4). Measurement: `awk '/^### 11. Projects with a module root/{f=1} f && /^\/tmp\/ailang-worldd serve /{print NR": "$0}' docs/QUICKSTART.md` → exactly `969: /tmp/ailang-worldd serve --db /tmp/se-world/world.db --ailang-bin $PIN \`. Any new line with that prefix anywhere after line 914 makes it two and reds `TestQuickstartSection11FlagsMatchTheCLI`. The new block uses `"$GRADE/ailang-worldd" serve`. The exact text is in C5. The executor executes the block verbatim once, with the pinned binary, and records its output. AC7's file list therefore grows from six to **seven** (`docs/QUICKSTART.md`). `requireWorkbenchLiveQuickstart` (workbench_live_chrome_test.go:21-41) requires that the existing §11 strings stay. The edit only APPENDS.

**PD8 — no load-sensitive timing.** No new test uses `time.Sleep`, a sub-second timeout, an elapsed-deadline stimulus or `t.Parallel()`. Every operational stimulus is a returned error value: `errSentinelInternal`, `context.DeadlineExceeded` (returned, not waited for), `&store.ReferenceIndexUnavailableError{}`, `&store.ObjectTooLargeError{}`. The only bounds are `boundedTestContext` (30 s, m7b_deadline_test.go:10) and the handler's `readDeadline` (10 s) — both are ceilings on a correct run, never the stimulus. `t.Parallel()` is forbidden because host/daemon is a write+fork package (`TestNoParallelWriteForkPackagesOutsideVerifygate`, forklocked_write_test.go:717). Measured: 0 `t.Parallel()` today in host/daemon and cmd/ailang-worldd tests (control: 19 in host/verifygate). Wall-time budget: see the Baseline load note.

**PD9 — milestone re-cut (design M1–M3 → C1–C5, each ≤ ~250 changed LOC including tests).**

| Plan | Design | Content | est. prod / test LOC |
|---|---|---|---|
| C1 | M1 (daemon half) + M3's cross-check | Config fields, `grader` field, `configureGrader`, `loadProofKey`, `proofCompilerVersionLine`, the `proofValidator` seam, shared test rig | 70 / 175 |
| C2 | M1 (CLI half) | two serve flags, three usage spots | 30 / 55 |
| C3 | M2 | `objectGradeView`, the reason constants and map, the `workbench.go` construction site | 90 / 90 |
| C4 | M3 (real-store half) | the 8 real-store reason tests incl. the truncation boundary | 0 / 230 |
| C5 | M3 (seam half) + gates | mapping table, 3 sanitized tests, the verdict sweep, QUICKSTART, the full mutation re-drill, the final gates | 0 / 220 + ~25 doc |

The cross-check moves from design M3 to C1 because everything it needs (`proofCompilerVersionLine`, a configured daemon) exists at C1. That lets M11 be drilled at the earliest boundary.

**PD10 — a validator result with no verdict is an operational error.** If `res.Err() == nil`, `res.Validated()` is false AND `res.Unsupported()` is false (the zero value), the loop returns `errGradeNoVerdict` (sanitized 500). It never silently counts the result as a subject mismatch. A real validator cannot produce this state. The scripted grader can (`zero: true`), so the branch is pinnable (planner mutant M14), not dead.

**PD11 — `TestProofKeyWiringConstructsTheValidator` (planner-added positive control).** The design delegates the wiring's positive observable to AC3, which lands at C3. Without this test the C1 mutants M15 (drop `d.grader = v`) and M16 survive at their own boundary. It asserts: configured → `d.grader.(*evidence.Validator)` ok; unconfigured → `d.grader == nil`.

**PD12 — the wiring lives in `workbench_grade.go` as `(*Daemon).configureGrader(cfg Config) error`; `New` calls it.** This matches the design's Files list ("`loadProofKey`, `proofCompilerVersionLine` … in workbench_grade.go"). It keeps daemon.go's diff at ~20 LOC. The design's mutant M3 therefore edits `workbench_grade.go`, not `daemon.go`.

**PD13 — the counting wrapper counts only proof-ID scans.** `countingProofScanStore` counts `ObjectsBySemanticID` calls with `id == evidence.ProofSemanticID`. Today no workbench path calls `ObjectsBySemanticID` at all (only `handlers.go:551`, the `/v1` route), but scoping the count keeps AC1's zero honest if a future pane adds a different-ID scan.

**PD14 — repo-wide enumerating guards: measured, and none moves.** Each item below gives the command and its output (run at `b8581e0`):

| Guard | Does this sprint touch it? | Measurement |
|---|---|---|
| **AC10 subprocess-site census** (`host/broker/registry_publish_test.go:1077` `enumerateSubprocessSites`, `:1135` `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`) | **No.** It parses NON-test `.go` under host/+cmd/ for `exec.Command`/`exec.CommandContext` calls. The new production code makes none (`os.Open`, `io.LimitReader`, HMAC and JSON only). The ONE new `exec.CommandContext` is in `workbench_grade_test.go` (the cross-check runs the fake's `--version`); test files are skipped (`strings.HasSuffix(path, "_test.go")` → return nil). **The file is NOT frozen.** If the executor finds a reason to add a production launch, it MUST add a driver to the `:1210` map in the same milestone. | `grep -rn "exec\.Command" host cmd --include="*.go" \| grep -v _test.go \| wc -l` → `12` (unchanged after the sprint: re-run at C5, must still be 12). Control: `:1210 "host/evidence/proof_producer.go": {driveEvidenceProducer}` (the census sees evidence's site). |
| **Evidence manifest** (`REQUIRED_EVIDENCE_TESTS`, verify_go.sh:27; `:322` banner, `:329` call `64`; verifygate `evidence_manifest_gate_test.go:205`, `mission_config_gate_test.go:454-455`) | **No.** No test is added to `host/evidence`. All new tests live in `host/daemon` and `cmd/ailang-worldd`. | `grep -rn "64 exact\|evidence_json\" 64\|all 64 required" --include='*.go' --include='*.sh' host scripts \| wc -l` → `5`; must still be 5 with an empty diff in those files at C5. |
| **`TestPublicAuthoritySurfaceIsFrozen`** (authority_test.go:427, `frozenPublicSurface` at :289) | **No.** No `host/evidence` file changes, so the exact exported list is unchanged. | C5: `git diff --quiet -- host/evidence` rc=0. |
| **Context-root census** `TestProductionContextRoots` (host/store/context_roots_test.go:71; exact map, host/daemon pins are only `daemon.go\|drain` and `workspace.go\|episodeHandler`) | **No — as long as no new production code calls `context.Background/TODO/WithoutCancel`.** `objectGradeView` uses the request ctx; `configureGrader`/`loadProofKey` take no ctx. **FORBIDDEN in production files this sprint: `context.Background()`, `context.TODO()`, `context.WithoutCancel`.** | measured map lines 41, 46 |
| **Deadline-free store reads** `TestNoNewDeadlineFreeStoreReads` (context_read_test.go:376; regex `.(GetObject\|…\|ReadObject\|…)(\s*context.(Background\|TODO)()`) | **No** — same rule as above | regex at :369 |
| **Production Go surface** `TestProductionGoSurface` (context_roots_test.go:176) | **No** — the new production file lives under host/ | — |
| **Fork-locked parallel tripwire** `TestNoParallelWriteForkPackagesOutsideVerifygate` | **No** — no `t.Parallel()` (PD8) | 0 today |
| **Subprocess sink gate** `TestSubprocessSinkGate` (whole repo incl. tests) | **No** — the cross-check uses `cmd.Output()` / `CommandContext(...).Output()`; it never binds a `bytes.Buffer`/`strings.Builder` to `cmd.Stdout/Stderr` | — |
| **Boundary dependency allowlist** `TestWorldBoundaryDependencyAllowlist` (host/boundary; group `cmd/ailang-worldd` forbids `host/registry` prefix, `net/http/httptest`, `net/http/httputil`, cloud SDKs in its own files' imports) | **No.** main.go adds no import. host/daemon's new import `host/evidence` adds exactly ONE package to the cmd closure, and it is not a forbidden prefix. | `comm -23 <(go list -deps ./host/evidence\|sort) <(go list -deps ./cmd/ailang-worldd\|sort)` → `github.com/sunholo-data/ailang-world/host/evidence` only (control: closure size 261) |
| **`host/boundary` file-count pin** `wantFileCount = 1` | **No** — no file in host/boundary | — |
| **daemon.go AST parsers** `TestMCPMountMethodSource`, `TestMCPDaemonBudgetWiring` (mcp_test.go:19,115 parse daemon.go's `Handler` mux and `projection.New` literal); `requireWorkbenchLiveQuickstart` reads daemon.go's `// Handler` doc | **No** — the daemon.go edits are confined to `Config`, the `Daemon` struct and one call in `New`. Do NOT touch `Handler()`, its doc comment, or the `projection.New(projection.Config{…})` literal. | — |
| **CLI help/usage tests** `TestEveryVerbHasHelp` (help_test.go), `TestServeWorkspaceFlags`, `TestServeExecFlagsReachTheDaemonConfig`, `TestQuickstartSection10/11FlagsMatchTheCLI` (all `Contains`/regex over the `usage` constant; `serve flags:` block) | **Additive only.** The new flag docs go inside the `serve flags:` block, indented like their neighbours, so `(?m)^\s+--proof-key-file\b` matches. No existing line changes. | — |
| **Reflection over `Config` / flag enumeration** | **None exists.** | `grep -rn "VisitAll\|NumField\|reflect.TypeOf(Config\|reflect.TypeOf(daemon.Config" --include=*_test.go host/daemon cmd/ailang-worldd host/verifygate` → no output (control: `grep -c frozenPublicSurface host/evidence/authority_test.go` → 4) |
| **`.ail`/module manifest/verify_ail ROOTS** | **No** `.ail` change | verify_ail baseline rc=0 |

---

## Frozen / must-not-change

- `host/evidence/**`, `host/store/**`, `host/workbench/**` (render.go, live.js, all tests), `host/boundary/**`, `host/verifygate/**`, `host/broker/**`: no hunk.
- `scripts/verify_go.sh`, `scripts/verify_ail.sh`, `.github/workflows/ci.yml`, `tools/launchd/*`, every `.ail`, `design_docs/world-mission*.md`: no hunk.
- `host/daemon/workbench_test.go` (incl. `TestWorkbenchObjectGrade`), `handlers_test.go`, `daemon_test.go`, `read_deadline_test.go`, `m7b_deadline_test.go`: no hunk. New tests REUSE their helpers (`newHandlerDaemon`, `fakeInterpreter`, `boundedTestContext`, `requestRecorder`, `workbenchTestObject`, `errSentinelInternal`, `internalDetailSentinel`, `testErrorLog`).
- `host/daemon/workbench.go`: the ONLY hunk is the construction site (C3). Lines 40-44 (the constant block) stay byte-identical.
- `host/daemon/daemon.go`: only the two Config fields, the `grader` struct field and the one `configureGrader` call. No edit to `Handler()`, `projection.New`, `releaseFromVersion` or `checkReadOrdering`.
- The changed file set at C5 is exactly: `host/daemon/{workbench_grade.go (new), workbench_grade_test.go (new), daemon.go, workbench.go}`, `cmd/ailang-worldd/{main.go, main_test.go}`, `docs/QUICKSTART.md`.

## Test write rule

Every DB, key file and fake interpreter goes under the test's own `t.TempDir()`. Tests never write tracked files, `$HOME` or `~/.ailang`. The pinned binary is read-only and is only RUN, with `--version`, in the cross-check's optional observation arm. No `t.Parallel()`, no `time.Sleep`.

## Mutation protocol

1. Apply ONE mutant alone to the boundary tree. Copy the file first: `cp F "$TMPDIR/F.orig"`.
2. It must compile: `go vet ./host/daemon ./cmd/ailang-worldd` rc=0. A build failure is NOT a kill.
3. Run the mutant's exact command. The named test must FAIL on its named assertion. A build failure, a panic, a SKIP, zero tests matched, an unrelated failure or a harness timeout are NOT kills.
4. Restore without git: `cp "$TMPDIR/F.orig" F && cmp F "$TMPDIR/F.orig"` (rc=0).
5. Re-run the command → green.

Record each mutant as `mutant | file:line | command | rc | failing assertion text | restore cmp rc | green rc`. A SURVIVED mutant blocks the milestone. C5 re-drills EVERY mutant M1–M18 once more on the final tree.

## Shared per-milestone gate template (NEW = that milestone's new test names joined by `|`)

```bash
#!/usr/bin/env bash
set -u
export AILANG_BIN=$HOME/.pinned-ailang/ailang; export PATH=/opt/homebrew/bin:$PATH
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter253-row105-design
go vet ./... || exit 1                                                   # compile fence incl. _test.go
log="$(mktemp)"; go test -count=1 -v ./host/daemon ./cmd/ailang-worldd -run "^(NEW)$" >"$log" 2>&1; rc=$?; cat "$log"; [ $rc -eq 0 ] || exit 1
for t in $(printf '%s' "NEW" | tr '|' ' '); do grep -q -- "^--- PASS: $t " "$log" || { echo "$t did not PASS"; exit 1; }; done
go test ./host/daemon ./cmd/ailang-worldd -count=1 || exit 1
go test -race ./host/daemon -count=1 || exit 1
go test -count=1 ./host/daemon -run '^TestWorkbenchObjectGrade$' -v | grep -q '^--- PASS: TestWorkbenchObjectGrade ' || exit 1
git diff --quiet -- host/evidence host/store host/workbench host/boundary host/verifygate host/broker scripts world .github tools || { echo "frozen tree touched"; exit 1; }
git diff --quiet -- host/daemon/workbench_test.go || { echo "workbench_test.go touched"; exit 1; }
```

The `-run` anchor matches top-level names in both packages. A name absent from one package simply matches nothing there; the PASS-line loop proves each name ran.

## Code shape (binding names)

**`host/daemon/workbench_grade.go`** (new; package daemon). Imports: `context`, `errors`, `fmt`, `io`, `os`, `strings`, `host/evidence`, `host/hashref`, `host/workbench`. (C1 imports only what C1 uses. C3 adds `errors`/`workbench` when it adds the projection. Go refuses unused imports, and `go vet` catches it.)

```go
const gradeCandidateScanLimit = 100   // C3
const proofKeySize = 32               // C1

const ( // C1
	gradeConfigIdentitiesWithoutKey = "--proof-required-identity needs --proof-key-file"
	gradeConfigNoInterpreter        = "the workbench grade validator needs the pinned interpreter; --ailang-bin is unset"
	gradeConfigKeyRefused           = "cannot load the proof key file"
	gradeConfigValidatorRefused     = "cannot construct the workbench grade validator"
)

const ( // C3 — frozen reason sentences (design D4 table, byte-exact)
	gradeReasonNoStoredEvidence    = "no stored evidence: the store holds no object with semantic ID world/proof-report/v1"
	gradeReasonBoundElsewhere      = "stored proof envelopes exist; none is bound to this subject"
	gradeReasonToolMismatch        = "subject-bound evidence authenticates under the daemon's key but names a compiler other than the daemon's pinned interpreter"
	gradeReasonProofFailed         = "subject-bound evidence authenticates under the daemon's key but records a failed proof outcome"
	gradeReasonProofIncomplete     = "subject-bound evidence authenticates under the daemon's key but does not verify every identity the daemon requires"
	gradeReasonUnauthenticated     = "a stored proof envelope does not authenticate under the daemon's key"
	gradeReasonStructuralCandidate = "a stored candidate object does not satisfy the proof envelope's identity, content or schema checks"
	gradeReasonNoneValidated       = "stored proof envelopes exist and none validates for this object"
	gradeScanTruncation            = " (scan bounded to the first 100 stored candidates)"
)

var errGradeNoVerdict       = errors.New("daemon: the grade validator returned no verdict")             // C3
var errGradeResolverRefused = errors.New("daemon: the grade resolver refused the validator's own seal") // C3

type proofValidator interface { /* PD1 */ }                                     // C1
func proofCompilerVersionLine(verbatim string) string                           // C1
func loadProofKey(path string) ([proofKeySize]byte, error)                      // C1
func (d *Daemon) configureGrader(cfg Config) error                              // C1
var gradeReasonSentences = map[evidence.UnsupportedReason]string{ /* 11 entries */ } // C3
func gradeReasonSentence(reason evidence.UnsupportedReason) string               // C3
func (d *Daemon) objectGradeView(ctx context.Context, ref hashref.HashRef) (workbench.GradeView, error) // C3
```

No new exported identifier anywhere except the two `daemon.Config` fields `ProofKeyFile`, `ProofRequiredIdentities`.

**`host/daemon/workbench_grade_test.go`** (new; `package daemon`). The helpers below are added at C1 and reused by later milestones. Their bodies are prescribed in C1.

`gradeTestVersion` (const), `gradeTestIdentity` (const), `producerVersionLine`, `wantGradeUnavailable`, `gradeRig` (struct), `writeProofKeyFile`, `newGradeDaemon`, `gradeDaemonConfig`, `mintGradeEnvelope`, `gradeReport`, `storeGradeEnvelope`, `storeGradeSubject`, `gradePage`.

---

## C1 — Validator wiring, key loading, version derivation (design M1 daemon half + the M3 cross-check)

**Budget: prod ~70 LOC, tests ~175 LOC.** Files:
- `host/daemon/workbench_grade.go` (new)
- `host/daemon/daemon.go` (Config fields, `grader` field, one call in `New`)
- `host/daemon/workbench_grade_test.go` (new)

**Implementation.**

1. In `daemon.go` `Config`, insert immediately BEFORE the `// ErrorLog receives …` comment:

```go
	// ProofKeyFile (`--proof-key-file`, row 105 D6) is an operator file holding
	// exactly 32 bytes: the MAC key the workbench grade validator authenticates
	// stored proof envelopes under. Empty leaves the daemon unconfigured: every
	// object page renders objectGradeUnavailableReason and reads no evidence.
	// Setting it needs AilangBin and at least one ProofRequiredIdentities entry;
	// startup refuses every half-configuration. The key is never logged.
	ProofKeyFile string
	// ProofRequiredIdentities (`--proof-required-identity`, repeatable) are the
	// identities every validated proof must verify. Refused without ProofKeyFile.
	ProofRequiredIdentities []string
```

2. In `daemon.go` `Daemon`, insert immediately AFTER the `templates templateBook` line's block (after the `arch`/`interpreterHash`/`templates` group):

```go
	// grader is the workbench grade validator (row 105): nil unless
	// --proof-key-file is set, when New wires the *evidence.Validator it built.
	// It is an interface so a test can compose validator outcomes a real store
	// cannot produce; production always holds the real validator.
	grader proofValidator
```

3. In `daemon.go` `New`, insert immediately AFTER the closing `}` of the `if cfg.AilangBin != "" { … }` block (currently daemon.go:698) and BEFORE `if invokeDeadline <= 0 || …`:

```go
	if err := d.configureGrader(cfg); err != nil {
		return nil, err
	}
```

4. `workbench_grade.go` (C1 content):

```go
// proofValidator is the daemon's validator seam (row 105 D2): *evidence.Validator
// satisfies it, and New wires nothing else.
type proofValidator interface {
	ValidateProof(ctx context.Context, reportRef, expectedSubject hashref.HashRef) evidence.ValidationResult
	Resolve(sealed evidence.ValidatedEvidence) evidence.ResolutionResult
}

// proofCompilerVersionLine mirrors the producer's derivation exactly
// (host/evidence/proof_producer.go:328): the RAW first line, trailing "\r"
// stripped, and deliberately NOT releaseFromVersion's TrimSpace.
func proofCompilerVersionLine(verbatim string) string {
	return strings.TrimRight(strings.SplitN(verbatim, "\n", 2)[0], "\r")
}

// loadProofKey reads exactly proofKeySize bytes from a regular file. It stats
// the OPEN file (no path TOCTOU) and reads at most proofKeySize+1 bytes.
// Errors name the path and sizes, never the key bytes.
func loadProofKey(path string) ([proofKeySize]byte, error) {
	var key [proofKeySize]byte
	f, err := os.Open(path)
	if err != nil {
		return key, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return key, err
	}
	if !info.Mode().IsRegular() {
		return key, fmt.Errorf("%s is not a regular file", path)
	}
	raw, err := io.ReadAll(io.LimitReader(f, proofKeySize+1))
	if err != nil {
		return key, err
	}
	if len(raw) > proofKeySize {
		return key, fmt.Errorf("%s holds more than %d bytes; a proof key is exactly %d", path, proofKeySize, proofKeySize)
	}
	if len(raw) < proofKeySize {
		return key, fmt.Errorf("%s holds %d bytes; a proof key is exactly %d", path, len(raw), proofKeySize)
	}
	copy(key[:], raw)
	return key, nil
}

// configureGrader wires the workbench grade validator (row 105 D6), or leaves
// d.grader nil when no key file is configured. Every half-configuration refuses.
func (d *Daemon) configureGrader(cfg Config) error {
	if cfg.ProofKeyFile == "" {
		if len(cfg.ProofRequiredIdentities) > 0 {
			return d.abort(StageConfig, gradeConfigIdentitiesWithoutKey, nil)
		}
		return nil
	}
	if d.interpreterHash.IsZero() {
		return d.abort(StageConfig, gradeConfigNoInterpreter, nil)
	}
	key, err := loadProofKey(cfg.ProofKeyFile)
	if err != nil {
		return d.abort(StageConfig, gradeConfigKeyRefused, err)
	}
	v, err := evidence.NewValidator(key, d.store, evidence.CompilerConfig{
		Compiler:          d.interpreterHash,
		CompilerVersion:   proofCompilerVersionLine(d.interpreterVersion),
		ObjectReadTimeout: readDeadline,
	}, cfg.ProofRequiredIdentities)
	if err != nil {
		return d.abort(StageConfig, gradeConfigValidatorRefused, err)
	}
	d.grader = v
	return nil
}
```

Each `if` is separate, and the `>` and `<` arms have separate messages, so each is separately killable (the design's M2 "`!= 32` → `< 32`" is our M2a; its "`<=` twin" is M2b).

**Shared test rig (C1, `workbench_grade_test.go`).** Imports at C1: `context`, `crypto/hmac`, `crypto/sha256`, `errors`, `html`, `net/http`, `os`, `os/exec`, `path/filepath`, `strings`, `testing`, `host/evidence`, `host/hashref`, `host/store`. Add each import only when the first use lands, because go vet reds on unused imports.

```go
// The discriminating version fixture (design V43): the producer's rule gives
// "AILANG v0.41.0 " and releaseFromVersion's TrimSpace gives "AILANG v0.41.0".
const gradeTestVersion = "AILANG v0.41.0 \r\nbuild abc\n"
const gradeTestIdentity = "world/types.ail"

// producerVersionLine is the TEST's own copy of proof_producer.go:328 — never
// the production proofCompilerVersionLine — so the two halves are independent.
func producerVersionLine(raw string) string {
	return strings.TrimRight(strings.SplitN(raw, "\n", 2)[0], "\r")
}

func wantGradeUnavailable(reason string) string {
	return "<p>GRADE UNAVAILABLE — " + html.EscapeString(reason) + "</p>"
}

type gradeRig struct {
	key      [32]byte
	compiler hashref.HashRef // SumSHA256(fake interpreter bytes)
	version  string          // PRODUCER-side: producerVersionLine(stdout of the fake's own --version run)
	rawVersion string        // that raw stdout
	bin, keyPath string
}

func writeProofKeyFile(t *testing.T, dir string, b []byte) string // os.WriteFile(filepath.Join(dir, "proof.key"), b, 0o600); returns path

// gradeDaemonConfig builds a configured Config over a fresh fake interpreter
// (gradeTestVersion) and a fresh 32-byte key file holding strings.Repeat("K", 32)
// (so a leak of the key bytes into an error is greppable), and fills the rig.
// The rig's version comes from RUNNING the fake: exec.CommandContext(boundedTestContext(t), bin, "--version").Output().
func gradeDaemonConfig(t *testing.T, identities []string) (Config, gradeRig)
	// Config{DBPath: filepath.Join(dir, "world.db"), BindHost: DefaultBindHost, AilangBin: bin,
	//        ProofKeyFile: keyPath, ProofRequiredIdentities: identities, ErrorLog: testErrorLog(t)}

// newGradeDaemon = gradeDaemonConfig + New + t.Cleanup(d.Close); t.Fatalf on New error.
func newGradeDaemon(t *testing.T, identities []string) (*Daemon, gradeRig)

func gradeReport(rig gradeRig, subject hashref.HashRef) evidence.ProofReportV1 {
	return evidence.ProofReportV1{Schema: evidence.ProofReportSchemaV1, Subject: subject, Compiler: rig.compiler,
		CompilerVersion: rig.version, Verified: []string{gradeTestIdentity}, CheckPassed: true, ProofSucceeded: true}
}
func mintGradeEnvelope(t *testing.T, key [32]byte, report evidence.ProofReportV1) []byte // PD6 body
func storeGradeEnvelope(t *testing.T, d *Daemon, env []byte) hashref.HashRef             // PD6 object, d.store.PutObject(boundedTestContext(t), …)
func storeGradeSubject(t *testing.T, d *Daemon, label string) hashref.HashRef            // workbenchTestObject(label, "test/grade-subject", hashref.SumSHA256([]byte("interface-"+label))) + PutObject
func gradePage(t *testing.T, d *Daemon, ref hashref.HashRef) (int, string)               // requestRecorder(t, d, http.MethodGet, "/workbench?object="+ref.String(), nil) → Code, Body.String()
```

**Tests first (red-first evidence = the C1 mutation drills on compiling code; on the pre-milestone tree they are build-red, which is NOT a kill).**

- **`TestProofKeyWiringRefusesEveryHalfConfiguration`** (AC2 + PD4). This is one table-driven test with `t.Run` per arm. Each arm builds its Config with `gradeDaemonConfig(t, []string{gradeTestIdentity})` and then edits ONE thing. It calls `New(boundedTestContext(t), cfg)`, requires `err != nil`, `errors.As(err, &se)` with `se *StartupError`, `se.Stage == StageConfig`, and `se.Detail ==` the arm's exact detail. It also asserts `!strings.Contains(err.Error(), strings.Repeat("K", 32))` (the rig's key bytes; the 31/33-byte arms write `K` runs of those lengths). Arms:
  1. `missing-file`: `ProofKeyFile = filepath.Join(dir, "absent.key")` → detail `gradeConfigKeyRefused`, `errors.Is(err, fs.ErrNotExist)`.
  2. `directory`: `ProofKeyFile = t.TempDir()` → detail `gradeConfigKeyRefused`, `strings.Contains(err.Error(), "is not a regular file")`.
  3. `31-bytes`: file of 31 bytes → detail `gradeConfigKeyRefused`, contains `"holds 31 bytes; a proof key is exactly 32"`.
  4. `33-bytes`: file of 33 bytes → detail `gradeConfigKeyRefused`, contains `"holds more than 32 bytes"`.
  5. `no-ailang-bin`: `AilangBin = ""` → detail `gradeConfigNoInterpreter` exactly, `se.Err == nil`.
  6. `no-identities`: `ProofRequiredIdentities = nil` → detail `gradeConfigValidatorRefused`, `errors.Is(err, evidence.ErrInvalidValidatorConfig)`, contains `"required identities are empty"`.
  7. `identities-without-key`: `ProofKeyFile = ""` → detail `gradeConfigIdentitiesWithoutKey` exactly.
  - A **control** arm at the top: the unedited config constructs (`err == nil`; close it). This proves that each refusal comes from its one edit.
- **`TestProofKeyWiringConstructsTheValidator`** (PD11): configured → `v, ok := d.grader.(*evidence.Validator)`, `ok && v != nil`. Unconfigured (`newHandlerDaemon(t)`) → `d.grader == nil`.
- **`TestWorkbenchGradeCompilerVersionDerivationsAgree`** (AC3 cross-check, design "M11"). Steps:
  - `d, rig := newGradeDaemon(t, []string{gradeTestIdentity})`.
  - Instrument checks: `rig.rawVersion == gradeTestVersion` (the shell delivered the fixture verbatim) and `d.interpreterVersion == rig.rawVersion` (the V38–V40 chain: archive verbatim → daemon unchanged).
  - Producer side: `p := producerVersionLine(rig.rawVersion)`; require `p == "AILANG v0.41.0 "`.
  - Daemon side: `dv := proofCompilerVersionLine(d.interpreterVersion)`; require `dv == p` with message `daemon-side %q != producer-side %q`.
  - Discrimination control: `releaseFromVersion(d.interpreterVersion) == "AILANG v0.41.0"` and `!= p`. This proves the fixture separates the two spellings.
  - Optional observation subtest `pinned-binary`: if `os.Getenv("AILANG_BIN") == ""` → `t.Skip("observation only")`; else run `$AILANG_BIN --version` via `exec.CommandContext(boundedTestContext(t), bin, "--version").Output()` and require `proofCompilerVersionLine(out) == producerVersionLine(out) == "AILANG v0.41.0"`.

**Acceptance (C1):** the shared template with NEW = `TestProofKeyWiringRefusesEveryHalfConfiguration|TestProofKeyWiringConstructsTheValidator|TestWorkbenchGradeCompilerVersionDerivationsAgree`.

**Mutation drills (C1).** Command for all of them: `go test -count=1 -v ./host/daemon -run '^(TestProofKeyWiringRefusesEveryHalfConfiguration|TestProofKeyWiringConstructsTheValidator|TestWorkbenchGradeCompilerVersionDerivationsAgree)$'`

| Mutant | Edit (workbench_grade.go) | Must red |
|---|---|---|
| M2a | `if len(raw) > proofKeySize {` → `if false && len(raw) > proofKeySize {` | `…/33-bytes` (a 33-byte read is then `len < 32` false → key accepted → `New` succeeds: "New accepted …") |
| M2b | `if len(raw) < proofKeySize {` → `if false && len(raw) < proofKeySize {` | `…/31-bytes` |
| M3 | `if d.interpreterHash.IsZero() {` → `if false && d.interpreterHash.IsZero() {` | `…/no-ailang-bin` (detail becomes `gradeConfigValidatorRefused`) |
| M11 | `proofCompilerVersionLine` body → `return strings.TrimSpace(strings.SplitN(verbatim, "\n", 2)[0])` | `TestWorkbenchGradeCompilerVersionDerivationsAgree` (`daemon-side "AILANG v0.41.0" != producer-side "AILANG v0.41.0 "`) |
| M13 | `if len(cfg.ProofRequiredIdentities) > 0 {` → `if false && len(…) > 0 {` | `…/identities-without-key` |
| M15 | delete the line `d.grader = v` (Go accepts the unused `v`? NO, so the mutant is `d.grader = nil; _ = v`) | `TestProofKeyWiringConstructsTheValidator` |
| M16 | `if !info.Mode().IsRegular() {` → `if false && !info.Mode().IsRegular() {` | `…/directory` (the read then fails with "is a directory", not "is not a regular file") |

**Controller commit message:** `feat(daemon): workbench grade validator wiring — proof key file, required identities, producer-mirroring version line (row 105 C1)`

---

## C2 — Serve flags and usage (design M1 CLI half)

**Budget: prod ~30 LOC, tests ~55 LOC.** Files: `cmd/ailang-worldd/main.go`, `cmd/ailang-worldd/main_test.go`.

**Implementation (main.go).**

1. Package doc comment: after the line `//	                    [--exec-max-output-bytes N]` add
   `//	                    [--proof-key-file <file> --proof-required-identity NAME ...]`
2. `usage` constant, serve synopsis: after `                      [--exec-max-output-bytes N]` add
   `                      [--proof-key-file <file> --proof-required-identity NAME ...]`
3. `usage` constant, `serve flags:` block: after the `--exec-max-output-bytes N` entry (its last line is `passes N bytes (default 67108864, 64 MiB)`) and BEFORE the blank line preceding `Exit codes:`, add exactly:

```
  --proof-key-file <file>
                       a file holding exactly 32 bytes: the MAC key the
                       workbench grade validator authenticates stored proof
                       envelopes under. Needs --ailang-bin and at least one
                       --proof-required-identity; read once at startup and
                       never logged. Unset, every object page reads "no
                       canonical host projection" and reads no evidence
  --proof-required-identity NAME
                       an identity every validated proof must verify
                       (repeatable; needs --proof-key-file)
```

4. `runServe`: after the `exec-max-output-bytes` `fs.Func` block, add:

```go
	proofKeyFile := fs.String("proof-key-file", "", "a file holding exactly 32 bytes: the workbench grade validator's MAC key")
	var proofIdentities []string
	fs.Func("proof-required-identity", "an identity every validated proof must verify (repeatable; needs --proof-key-file)", func(v string) error {
		if v == "" {
			return fmt.Errorf("is empty")
		}
		proofIdentities = append(proofIdentities, v)
		return nil
	})
```

   In the `daemon.Config{…}` literal, append `ProofKeyFile: *proofKeyFile, ProofRequiredIdentities: proofIdentities`.

**Tests first — `TestServeProofFlags`** (main_test.go; reuses `captureServeConfig` from quickstart_exec_test.go, same package). The arms run in this order:
1. **Help:** `run([]string{"serve", "--help"}, …)` exit 0. Cut at `"serve flags:"`. `(?m)^\s+--proof-key-file\b` and `(?m)^\s+--proof-required-identity\b` both match.
2. **Real startup refusal (before any capture):** `run([]string{"serve", "--db", filepath.Join(dir, "world.db"), "--bind", "127.0.0.1:0", "--proof-key-file", key, "--proof-required-identity", "world/types.ail"}, …)` where `key` is a 32-byte file in `t.TempDir()` and NO `--ailang-bin` → `exitFatal`, and stderr contains `the workbench grade validator needs the pinned interpreter; --ailang-bin is unset`. New refuses before Listen, so no socket is bound.
3. **Wiring:** `got := captureServeConfig(t)`. `run(serve --db /x/world.db --proof-key-file /k/proof.key --proof-required-identity a --proof-required-identity b)` exit 0. `got.ProofKeyFile == "/k/proof.key"` and `reflect.DeepEqual(got.ProofRequiredIdentities, []string{"a", "b"})`.
4. **Unset:** `*got = daemon.Config{}`. `run(serve --db /x/world.db)` → `got.ProofKeyFile == ""` and `got.ProofRequiredIdentities == nil`.
5. **Empty identity:** `run(serve --db /x/world.db --proof-required-identity "")` → `exitUsage`, and stderr contains `is empty`.

**Acceptance (C2):** template with NEW = `TestServeProofFlags`. Also `go test -count=1 ./cmd/ailang-worldd -run '^(TestEveryVerbHasHelp|TestServeWorkspaceFlags|TestServeExecFlagsReachTheDaemonConfig|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI)$' -v` all PASS.

**Mutation drills (C2).** Command: `go test -count=1 -v ./cmd/ailang-worldd -run '^TestServeProofFlags$'`

| Mutant | Edit (main.go) | Must red |
|---|---|---|
| M17a | drop `ProofKeyFile: *proofKeyFile` from the Config literal | arm 3 (`got.ProofKeyFile == ""`) |
| M17b | `if v == "" {` → `if false && v == "" {` | arm 5 (exit 0, not usage) |
| M17c | `proofIdentities = append(proofIdentities, v)` → `proofIdentities = []string{v}` | arm 3 (DeepEqual `[b]`) |

**Controller commit message:** `feat(ailang-worldd): serve --proof-key-file and --proof-required-identity (row 105 C2)`

---

## C3 — The projection and the construction site (design M2)

**Budget: prod ~90 LOC, tests ~90 LOC.** Files: `host/daemon/workbench_grade.go`, `host/daemon/workbench.go`, `host/daemon/workbench_grade_test.go`.

**Implementation.** Add to workbench_grade.go the C3 constants, the two error vars, and:

```go
var gradeReasonSentences = map[evidence.UnsupportedReason]string{
	evidence.UnsupportedUnauthenticatedReport: gradeReasonUnauthenticated,
	evidence.UnsupportedToolMismatch:          gradeReasonToolMismatch,
	evidence.UnsupportedProofFailed:           gradeReasonProofFailed,
	evidence.UnsupportedProofIncomplete:       gradeReasonProofIncomplete,
	evidence.UnsupportedInvalidRef:            gradeReasonStructuralCandidate,
	evidence.UnsupportedMissing:               gradeReasonStructuralCandidate,
	evidence.UnsupportedOversize:              gradeReasonStructuralCandidate,
	evidence.UnsupportedHashMismatch:          gradeReasonStructuralCandidate,
	evidence.UnsupportedWrongSemanticID:       gradeReasonStructuralCandidate,
	evidence.UnsupportedWrongInterface:        gradeReasonStructuralCandidate,
	evidence.UnsupportedMalformed:             gradeReasonStructuralCandidate,
}

// gradeReasonSentence maps a non-mismatch validator reason to its frozen
// sentence; any token outside the map renders the branch invariant.
func gradeReasonSentence(reason evidence.UnsupportedReason) string {
	if sentence, ok := gradeReasonSentences[reason]; ok {
		return sentence
	}
	return gradeReasonNoneValidated
}

// objectGradeView projects the object page's grade (row 105 D4). Unconfigured,
// it is the constant with zero reads. Configured, it scans at most
// gradeCandidateScanLimit proof envelopes (fetching one extra row to prove
// truncation) and renders PROVEN only through ValidateProof → Resolve.
func (d *Daemon) objectGradeView(ctx context.Context, ref hashref.HashRef) (workbench.GradeView, error) {
	if d.grader == nil {
		return workbench.NewGradeUnavailable(objectGradeUnavailableReason), nil
	}
	candidates, err := d.reads.ObjectsBySemanticID(ctx, evidence.ProofSemanticID, "", gradeCandidateScanLimit+1)
	if err != nil {
		return workbench.GradeView{}, err
	}
	truncated := len(candidates) == gradeCandidateScanLimit+1
	if truncated {
		candidates = candidates[:gradeCandidateScanLimit]
	}
	if len(candidates) == 0 {
		return workbench.NewGradeUnavailable(gradeReasonNoStoredEvidence), nil
	}
	firstFailure := ""
	for _, candidate := range candidates {
		res := d.grader.ValidateProof(ctx, candidate.Hash, ref)
		if err := res.Err(); err != nil {
			return workbench.GradeView{}, fmt.Errorf("daemon: validate proof candidate %s: %w", candidate.Hash, err)
		}
		if sealed, ok := res.Validated(); ok {
			resolved := d.grader.Resolve(sealed)
			if _, proven := resolved.Proven(); !proven {
				return workbench.GradeView{}, fmt.Errorf("%w: candidate %s: %v", errGradeResolverRefused, candidate.Hash, resolved.Err())
			}
			return workbench.GradeView{Available: true, Label: workbench.GradePROVEN}, nil
		}
		reason, unsupported := res.Unsupported()
		if !unsupported {
			return workbench.GradeView{}, fmt.Errorf("%w: candidate %s", errGradeNoVerdict, candidate.Hash)
		}
		if reason != evidence.UnsupportedSubjectMismatch && firstFailure == "" {
			firstFailure = gradeReasonSentence(reason)
		}
	}
	sentence := gradeReasonBoundElsewhere
	if firstFailure != "" {
		sentence = firstFailure
	}
	if truncated {
		sentence += gradeScanTruncation
	}
	return workbench.NewGradeUnavailable(sentence), nil
}
```

**`workbench.go` construction site.** Inside the `?object=` branch, immediately BEFORE `page.Object = &workbench.ObjectView{`, insert:

```go
		grade, err := d.objectGradeView(ctx, ref)
		if err != nil {
			d.writeWorkbenchStoreError(w, r, ctx, err)
			return
		}
```

Replace `Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason),` with `Grade: grade,`. Lines 40-44 stay untouched. Note that `objectGradeUnavailableReason` keeps a use: it is used in `objectGradeView`.

**Tests first.**

- **`TestWorkbenchGradeUnconfiguredServesTheConstantWithoutScanning`** (AC1). Test-file type: `countingProofScanStore{readStore; calls *atomic.Int32}` overrides `ObjectsBySemanticID(ctx, id, after, limit)`, increments when `id == evidence.ProofSemanticID`, and forwards (PD13). Steps:
  - `d := newHandlerDaemon(t)`.
  - PutObject a proof-labelled object `workbenchTestObject("grade-unconfigured-proof", "world/proof-report/v1", evidence.InterfaceHashV1)` and a subject.
  - Set `d.reads = countingProofScanStore{readStore: d.reads, calls: &n}`.
  - For both refs: `gradePage` → 200, `strings.Count(body, wantGradeUnavailable(objectGradeUnavailableReason)) == 1`, no `<span>PROVEN</span>`.
  - Finally `n.Load() == 0` (failure text `ObjectsBySemanticID(proof) calls = %d, want 0`).
  - **Control (same test):** a configured `newGradeDaemon` with the same wrapper on one object page gives `n.Load() == 1`. This proves the counter can fire.
- **`TestWorkbenchGradeProvenFromSubjectBoundEnvelope`** (AC3). Steps:
  - `d, rig := newGradeDaemon(t, []string{gradeTestIdentity})`; `subject := storeGradeSubject(t, d, "grade-proven-subject")`.
  - `env := mintGradeEnvelope(t, rig.key, gradeReport(rig, subject))`; `envRef := storeGradeEnvelope(t, d, env)`.
  - Instrument: `rig.version == "AILANG v0.41.0 "` (producer-side, discriminating).
  - **Premise:** `res := d.grader.ValidateProof(boundedTestContext(t), envRef, subject)`; `sealed, ok := res.Validated()` must be ok; `d.grader.Resolve(sealed).Proven()` must be ok.
  - Subject page: 200; `strings.Count(body, "<p><span>PROVEN</span></p>") == 1`; no `GRADE UNAVAILABLE`; no `verdict-`; none of `<span>TESTED</span>`, `<span>ATTESTED</span>`, `<span>CLAIMED</span>`.
  - **Honesty control:** the ENVELOPE's own page (`gradePage(t, d, envRef)`) shows `wantGradeUnavailable(gradeReasonBoundElsewhere)` exactly once. An object never wears another object's proof.

**Acceptance (C3):** template with NEW = `TestWorkbenchGradeUnconfiguredServesTheConstantWithoutScanning|TestWorkbenchGradeProvenFromSubjectBoundEnvelope`. Also `sed -n '40,44p' host/daemon/workbench.go` is byte-identical to base (`git diff -U0 host/daemon/workbench.go` shows only the construction-site hunk).

**Mutation drills (C3).** Command: `go test -count=1 -v ./host/daemon -run '^(TestWorkbenchGradeUnconfiguredServesTheConstantWithoutScanning|TestWorkbenchGradeProvenFromSubjectBoundEnvelope)$'`

| Mutant | Edit (workbench_grade.go) | Must red |
|---|---|---|
| M1 | move the `candidates, err := d.reads.ObjectsBySemanticID(...)` line and its `if err != nil {…}` ABOVE the `if d.grader == nil` block | Unconfigured (`calls = 1, want 0`) |
| M4 | `d.grader.ValidateProof(ctx, candidate.Hash, ref)` → `d.grader.ValidateProof(ctx, candidate.Hash, candidate.Hash)` | Proven (subject page shows the bound-elsewhere reason) |
| M6 | `Label: workbench.GradePROVEN` → `Label: workbench.GradeTESTED` | Proven (`<p><span>PROVEN</span></p>` count 0) |
| M10 | the PROVEN view → `workbench.GradeView{Available: true, Label: workbench.GradePROVEN, HasVerdict: true, Verdict: workbench.VerdictPass}` | Proven (exact line count 0 / `verdict-` present) |
| M11 (re-drill) | as C1 | Proven AND the cross-check (tool_mismatch → no PROVEN) |

**Controller commit message:** `feat(daemon): workbench object page projects PROVEN from subject-bound proof envelopes (row 105 C3)`

---

## C4 — Real-store reason states and the truncation boundary (design M3, real-store half)

**Budget: tests ~230 LOC, no production change.** File: `host/daemon/workbench_grade_test.go`.

Each test constructs its PREMISE and checks it directly, THEN asserts the exact sentence with `strings.Count(body, wantGradeUnavailable(X)) == 1` and status 200. Everything runs on `d, rig := newGradeDaemon(t, []string{gradeTestIdentity})` and `subject := storeGradeSubject(t, d, "<test-specific label>")`. `ctx := boundedTestContext(t)`.

- **`TestWorkbenchGradeReasonNoStoredEvidence`**: premise `d.store.ObjectsBySemanticID(ctx, evidence.ProofSemanticID, "", 1)` returns 0 rows → `gradeReasonNoStoredEvidence`.
- **`TestWorkbenchGradeReasonBoundToOtherSubjects`**: envelope for `other := hashref.SumSHA256([]byte("grade-other-subject"))` under `rig.key`. Premise: `d.grader.ValidateProof(ctx, envRef, subject).Unsupported()` == `subject_mismatch` → `gradeReasonBoundElsewhere`. No qualifier: `!strings.Contains(body, html.EscapeString(gradeScanTruncation))`.
- **`TestWorkbenchGradeReasonUnauthenticatedEnvelope`**: envelope for `subject` minted under `foreign := [32]byte{9, 9, 9}`. Premise reason `unauthenticated_report` → `gradeReasonUnauthenticated`. Also assert the body does NOT contain `html.EscapeString(gradeReasonBoundElsewhere)` (M7's observable).
- **`TestWorkbenchGradeReasonToolMismatch`**, three subtests on separate daemons:
  1. `compiler`: report `Compiler = hashref.SumSHA256([]byte("another compiler"))` → `gradeReasonToolMismatch` (premise `tool_mismatch`).
  2. `trimspace-version`: report `CompilerVersion = releaseFromVersion(rig.rawVersion)` (= `"AILANG v0.41.0"`, the TrimSpace spelling) → `gradeReasonToolMismatch`. This pins the drift the design's Risks row names.
  3. `mismatch-sorted-first` (PD for M12): one tool-mismatch envelope `T` for `subject`, plus foreign envelopes for subjects `fmt.Sprintf("grade-foreign-%d", i)`, `i = 0..63`. Stop minting at the first foreign envelope whose `ref.String() < T-ref.String()`, and store that one. If no such `i` exists, the arm fails as an instrument failure (`t.Fatal("no foreign envelope sorts before the tool-mismatch candidate")`). Premise: `d.store.ObjectsBySemanticID(ctx, evidence.ProofSemanticID, "", 10)` returns the foreign ref FIRST. Assert the sentence is `gradeReasonToolMismatch`.
- **`TestWorkbenchGradeReasonProofFailed`**: report `CheckPassed = false` → premise `proof_failed` → `gradeReasonProofFailed`.
- **`TestWorkbenchGradeReasonProofIncomplete`**: report `Verified = []string{"world/other.ail"}` → premise `proof_incomplete` → `gradeReasonProofIncomplete`.
- **`TestWorkbenchGradeReasonStructuralCandidate`**, three subtests, each on a fresh daemon:
  1. `wrong-interface`: a valid envelope for `subject` stored with `InterfaceHash = hashref.SumSHA256([]byte("not the envelope interface"))` (still `SemanticID = evidence.ProofSemanticID`).
  2. `malformed`: payload `[]byte("not an envelope")` with the correct identities.
  3. `oversize`: payload `bytes.Repeat([]byte("x"), evidence.MaxBytes+1)` with the correct identities.

  Each premise is the direct `ValidateProof` reason (`wrong_interface` / `malformed` / `oversize`). Each asserts `gradeReasonStructuralCandidate`.
- **`TestWorkbenchGradeReasonScanTruncationIsDisclosed`** (round 2 boundary, one daemon):
  - Store exactly 100 foreign envelopes (subjects `grade-trunc-%03d`, `i = 0..99`).
  - Premise: `len(d.store.ObjectsBySemanticID(ctx, ProofSemanticID, "", 101)) == 100`.
  - Page → `wantGradeUnavailable(gradeReasonBoundElsewhere)` exactly, and the body does NOT contain `html.EscapeString(gradeScanTruncation)`.
  - Store the 101st (`grade-trunc-100`). Premise `== 101`.
  - Page → `wantGradeUnavailable(gradeReasonBoundElsewhere + gradeScanTruncation)` exactly once.
  - Record this test's wall time under `-race` in the boundary report (PD8 budget).

**Acceptance (C4):** template with NEW = the 8 names. Also, under `-race`, `go test -race -count=1 -v ./host/daemon -run '^(TestWorkbenchGrade|TestProofKeyWiring)'` with total wall < 15 s recorded.

**Mutation drills (C4).** Command: `go test -count=1 -v ./host/daemon -run '^TestWorkbenchGradeReason'`

| Mutant | Edit (workbench_grade.go) | Must red |
|---|---|---|
| M7 | map entry `evidence.UnsupportedUnauthenticatedReport: gradeReasonUnauthenticated` → `: gradeReasonBoundElsewhere` | `TestWorkbenchGradeReasonUnauthenticatedEnvelope` |
| M8a | delete the `if truncated { sentence += gradeScanTruncation }` block | `…ScanTruncationIsDisclosed` (101 arm lacks the qualifier) |
| M8b | `truncated := len(candidates) == gradeCandidateScanLimit+1` → `truncated := len(candidates) >= gradeCandidateScanLimit` | `…ScanTruncationIsDisclosed` (the 100 arm carries the qualifier) |
| M12 | `if reason != evidence.UnsupportedSubjectMismatch && firstFailure == "" {` → `if firstFailure == "" {` | `…ToolMismatch/mismatch-sorted-first` (the sentence becomes the fallback `gradeReasonNoneValidated`) |
| M7b | map entry `UnsupportedToolMismatch: gradeReasonToolMismatch` → `: gradeReasonProofFailed` | `…ToolMismatch` |

**Controller commit message:** `test(daemon): every real-store grade reason and the exactly-100 truncation boundary (row 105 C4)`

---

## C5 — Seam table, sanitized paths, verdict sweep, QUICKSTART, final gates (design M3, seam half)

**Budget: tests ~220 LOC, docs ~25 lines.** Files: `host/daemon/workbench_grade_test.go`, `docs/QUICKSTART.md`.

**Seam helpers (test file):** `scriptedGrader`, `scriptedReader` (PD1). Also `scriptedRig(t)`, which returns:
- a configured `d`, `rig`, `subject`;
- ONE stored placeholder candidate `cand` (`workbenchTestObject("grade-scripted-candidate", evidence.ProofSemanticID, evidence.InterfaceHashV1)` via PutObject), so the scan is non-empty;
- `newV(reader scriptedReader, key [32]byte) *evidence.Validator` = `evidence.NewValidator(key, reader, evidence.CompilerConfig{Compiler: rig.compiler, CompilerVersion: rig.version, ObjectReadTimeout: 10 * time.Second}, []string{gradeTestIdentity})`.

`good := mintGradeEnvelope(t, rig.key, gradeReport(rig, subject))`, `goodRef := hashref.SumSHA256(good)`, `goodMeta := store.ObjectMeta{InterfaceHash: evidence.InterfaceHashV1, SemanticID: evidence.ProofSemanticID}`.

- **`TestWorkbenchGradeReasonMapsEveryUnsupportedReason`** (AC4 table). For each row: set `d.grader = <grader>`, GET the subject page, assert status and the exact sentence. Also assert that no OTHER sentence constant from the D4 list appears (loop over all 9 sentences + `objectGradeUnavailableReason`; only the expected one may be present).

  Notation: `newV(r)` means `newV(r, rig.key)`; `reader{m, p, e}` means `scriptedReader{meta: m, payload: p, err: e}`; `ptr(x)` is a local `func ptr(h hashref.HashRef) *hashref.HashRef { return &h }`. Envelopes in arms 9–12 are `mintGradeEnvelope(t, rig.key, <gradeReport(rig, subject) with ONE field changed>)`, served by `reader{goodMeta, thatEnvelope, nil}` with `ref: ptr(SumSHA256(thatEnvelope))`.

  | # | arm | grader | expect |
  |---|---|---|---|
  | 1 | invalid_ref | `{v: newV(reader{goodMeta, good, nil}), ref: &hashref.HashRef{}}` | 200 structural |
  | 2 | missing | `{v: newV(reader{goodMeta, nil, nil}), ref: &goodRef}` | 200 structural |
  | 3 | oversize | `{v: newV(reader{err: &store.ObjectTooLargeError{Size: evidence.MaxBytes + 1, MaxBytes: evidence.MaxBytes}})}` | 200 structural |
  | 4 | hash_mismatch | `{v: newV(reader{goodMeta, good, nil})}` (no ref override: `cand.Hash != goodRef`) | 200 structural |
  | 5 | wrong_semantic_id | `{v: newV(reader{store.ObjectMeta{InterfaceHash: evidence.InterfaceHashV1, SemanticID: "other/id"}, good, nil}), ref: &goodRef}` | 200 structural |
  | 6 | wrong_interface | `{v: newV(reader{store.ObjectMeta{InterfaceHash: hashref.SumSHA256([]byte("x")), SemanticID: evidence.ProofSemanticID}, good, nil}), ref: &goodRef}` | 200 structural |
  | 7 | malformed | `{v: newV(reader{goodMeta, []byte("{}"), nil}), ref: ptr(hashref.SumSHA256([]byte("{}")))}` | 200 structural |
  | 8 | unauthenticated_report | `{v: newV(reader{goodMeta, good, nil}, [32]byte{7}), ref: &goodRef}` | 200 unauthenticated |
  | 9 | subject_mismatch | payload = envelope for another subject, `ref: ptr(its hash)` | 200 bound-elsewhere |
  | 10 | tool_mismatch | payload = envelope with another compiler | 200 tool |
  | 11 | proof_failed | envelope `Errors: 1` | 200 proof-failed |
  | 12 | proof_incomplete | envelope `Verified: ["world/other.ail"]` | 200 incomplete |
  | 13 | validated + resolved | `{v: newV(reader{goodMeta, good, nil}), ref: &goodRef}` | 200 `<p><span>PROVEN</span></p>` |
  | 14 | validated + resolve refused | as 13 plus `resolver: newV(reader{}, rig.key)` (a DIFFERENT real validator → `ErrForeignSeal`) | 500 `>Internal<` |
  | 15 | no verdict (PD10) | `{zero: true}` | 500 `>Internal<` |

  Before each arm, the test calls the arm's grader's `ValidateProof(ctx, cand.Hash, subject)` directly and asserts the premise token (`Unsupported()` reason, or `Validated()` ok for 13/14, or neither for 15). This is the "a reason is never asserted without its premise" rule.

  A pure-function sub-check runs inside the same test: `gradeReasonSentence(r)` for the 11 mapped tokens equals the table's sentence, and `gradeReasonSentence(evidence.UnsupportedReason("future_reason")) == gradeReasonNoneValidated` (the unknown-token arm).

- **`TestWorkbenchGradeScanFailureIsSanitized`** (AC5). Test-file wrapper: `proofScanFailingStore{readStore; err error}` overrides `ObjectsBySemanticID` to return `nil, s.err` when `id == evidence.ProofSemanticID`, and forwards otherwise. Configured daemon, `d.errLog = &bytes.Buffer{}`. Arms:
  - `errSentinelInternal` → 500, body contains `>Internal<`, body lacks `internalDetailSentinel`, errLog contains `internalDetailSentinel`.
  - `context.DeadlineExceeded` (a returned value — no waiting) → 503 `>Timeout<`.
  - `&store.ReferenceIndexUnavailableError{}` → 503 `>ReferenceIndexUnavailable<`.
  - Every arm: the body lacks `GRADE UNAVAILABLE`.
  - **Control:** with the wrapper's `err` nil, the page is 200.
- **`TestWorkbenchGradeValidationReadErrorIsSanitized`** (AC5): the grader is `{v: newV(scriptedReader{err: errSentinelInternal})}`. Premise: `ValidateProof(...).Err() != nil`. Page → 500 `>Internal<`, body lacks `internalDetailSentinel`, errLog contains it.
- **`TestWorkbenchGradeValidatorRefusalIsSanitized`** (AC5): arm 14's grader. Premise: `Validated()` ok, and `foreign.Resolve(sealed).Err()` is `evidence.ErrForeignSeal`. Page → 500, body lacks `<span>PROVEN</span>`, errLog contains `refused the validator's own seal`.
- **`TestWorkbenchGradeNeverRendersVerdictOrUnmintedLabels`** (AC6): re-drive table arms 1–13 (factor the table into a test-file func `gradeSeamArms(t, rig, subject, cand) []gradeSeamArm` shared with the mapping test). For every 200 arm: the body contains no `verdict-`, no `<span>TESTED</span>`, no `<span>ATTESTED</span>`, no `<span>CLAIMED</span>`. For every unavailable arm: no `<span>PROVEN</span>`. For arm 13: exactly `<p><span>PROVEN</span></p>` once.

**QUICKSTART (PD7).** In `docs/QUICKSTART.md` `## 11. Open the live workbench`, insert after the `- **Decisions:** …` bullet paragraph (before the blank line preceding `**JavaScript disabled or blocked:**`):

```markdown
- **Object grade:** an object page shows `PROVEN` only when the daemon holds a
  proof-envelope validator (`serve --proof-key-file <32-byte file>
  --proof-required-identity NAME`, both needing `--ailang-bin`) and a stored
  `world/proof-report/v1` envelope, authenticated under that key, names the object
  as its subject. Otherwise the page says `GRADE UNAVAILABLE —` and one reason: no
  validator configured (`no canonical host projection…`), no stored envelopes
  (`no stored evidence…`), envelopes bound to other objects, an envelope that does
  not authenticate, a different compiler, a failed proof, a missing identity, or a
  malformed candidate. At most 100 candidates are checked, and a longer scan says
  so. Nothing in this release writes envelopes, so a configured daemon shows
  `no stored evidence` until a producer does.
```

Append at the END of section 11 (end of file) a subsection:

````markdown
**Check the grade validator wiring** (disposable store, port 7645, from the repository root):

```bash
export PIN=$HOME/.pinned-ailang/ailang
GRADE=$(mktemp -d)
go build -o "$GRADE/ailang-worldd" ./cmd/ailang-worldd
head -c 32 /dev/urandom > "$GRADE/proof.key"
"$GRADE/ailang-worldd" serve --db "$GRADE/world.db" --ailang-bin "$PIN" --bind 127.0.0.1:7645 \
  --proof-key-file "$GRADE/proof.key" --proof-required-identity world/types.ail &
GRADE_PID=$!
for i in $(seq 1 60); do curl --silent --fail http://127.0.0.1:7645/v1/health >/dev/null && break; sleep 1; done
HEAD=$(curl --silent --fail http://127.0.0.1:7645/v1/registry/world/epoch-registry/v1 | python3 -c 'import json,sys; print(json.load(sys.stdin)["head"])')
curl --silent --fail "http://127.0.0.1:7645/workbench?object=$HEAD" | grep -o 'GRADE UNAVAILABLE — [^<]*'
kill $GRADE_PID; wait $GRADE_PID 2>/dev/null; rm -rf "$GRADE"
```

Expected: `GRADE UNAVAILABLE — no stored evidence: the store holds no object with semantic ID world/proof-report/v1`.
````

The executor runs that bash block VERBATIM once (bash, from the worktree root, with port 7645 free) and pastes its output into the boundary report. If the output differs, STOP and report; do not edit the expected line to match. Then it runs `go test -count=1 ./cmd/ailang-worldd -run '^(TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI)$'` and `go test -count=1 ./host/daemon -run '^TestWorkbenchLiveQuickstartControl$'`. Both must be green.

**Acceptance (C5):**
1. The template with NEW = `TestWorkbenchGradeReasonMapsEveryUnsupportedReason|TestWorkbenchGradeScanFailureIsSanitized|TestWorkbenchGradeValidationReadErrorIsSanitized|TestWorkbenchGradeValidatorRefusalIsSanitized|TestWorkbenchGradeNeverRendersVerdictOrUnmintedLabels`.
2. The QUICKSTART block's verbatim output, plus the three QUICKSTART tests above.
3. PD14 re-measures: exec-site count `12`; manifest-literal count `5`; `git diff --quiet -- host/evidence` rc=0.
4. `git diff --stat` names exactly the seven files of §Frozen.
5. `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` rc=0 (base: 16 identities / 40 named / 444 se-tools).
6. `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` with a ≥ 20-minute timeout, rc=0. A red ONLY in `host/verifygate` under `-race -timeout 8m` with `panic: test timed out` is the measured BASE condition (PD5); the controller re-runs it.
7. **The full mutation re-drill M1–M18** on the final tree, one at a time. The ledger lists all of them.

**Mutation drills (C5, new).** Command: `go test -count=1 -v ./host/daemon -run '^TestWorkbenchGrade(ReasonMapsEveryUnsupportedReason|ScanFailureIsSanitized|ValidationReadErrorIsSanitized|ValidatorRefusalIsSanitized|NeverRendersVerdictOrUnmintedLabels)$'`

| Mutant | Edit (workbench_grade.go) | Must red |
|---|---|---|
| M5 | replace the Resolve block with `return workbench.GradeView{Available: true, Label: workbench.GradePROVEN}, nil` directly after `if sealed, ok := res.Validated(); ok {` (keep `_ = sealed`) | `…ValidatorRefusalIsSanitized` (200 PROVEN) + table arm 14 |
| M9 | the scan `if err != nil { return workbench.GradeView{}, err }` → `if err != nil { return workbench.NewGradeUnavailable(gradeReasonNoStoredEvidence), nil }` | `…ScanFailureIsSanitized` (200 with a reason) |
| M14 | `if !unsupported {` → `if false && !unsupported {` | table arm 15 (200 bound-elsewhere instead of 500) |
| M7c | map entry `UnsupportedMalformed: gradeReasonStructuralCandidate` → delete the entry | table arm 7 (fallback sentence) |
| M18 | `if err := res.Err(); err != nil {` → `if err := res.Err(); false && err != nil {` | `…ValidationReadErrorIsSanitized` |

**Controller commit message:** `test(daemon): seam-driven grade reason table, sanitized operational paths, PROVEN-only sweep; docs: QUICKSTART object grade (row 105 C5)`

---

## Mutation index (design → plan)

| Design | Plan mutant | Milestone |
|---|---|---|
| M1 unconfigured-scan | M1 | C3 |
| M2 key-size | M2a, M2b | C1 |
| M3 unset-interpreter | M3 (in `configureGrader`, PD12) | C1 |
| M4 subject binding | M4 | C3 |
| M5 Resolve bypass | M5 | C5 |
| M6 label substitution | M6 | C3 |
| M7 reason-map swap | M7 (+ M7b, M7c) | C4/C5 |
| M8 qualifier dropped / `>=` | M8a, M8b | C4 |
| M9 operational→reason | M9 | C5 |
| M10 verdict channel | M10 | C3 |
| (AC3 cross-check, "M11") | M11 | C1, re-drilled C3 |
| planner | M12 mismatch-skip, M13 identities-without-key, M14 no-verdict, M15 grader assignment, M16 non-regular, M17a-c CLI, M18 read-error | C1/C2/C4/C5 |

## Risks

| Risk | Mitigation |
|---|---|
| host/daemon `-race` already uses 358 s of the 480 s verify_go budget | PD8: no sleeps; the C4/C5 wall time is recorded; the 101-envelope test is the only bulk test (~101 small PutObjects) |
| `verifygate -race` 8-minute timeout under load (BASE) | PD5: classify via the addendum; the controller re-runs verify_go |
| Haiku edits a frozen helper file "to reuse" a helper | §Frozen + the template's `git diff --quiet` lines |
| An apostrophe sentence asserted raw | PD2 helper only; a raw `strings.Contains(body, gradeReasonToolMismatch)` is a plan violation |
| A QUICKSTART serve line with a `/tmp/ailang-worldd serve` prefix | PD7; the C5 acceptance runs both section tests |

## Baseline addendum (isolated reruns, with `AILANG_BIN`, planner-measured)

| Gate (AILANG_BIN exported) | Result |
|---|---|
| `go test ./host/daemon ./cmd/ailang-worldd -count=1` | **rc=0**, 160 s wall |
| `go test -race ./host/daemon -count=1` | **rc=0**, 197 s wall |
| `go test -race -timeout 8m ./host/verifygate -count=1` (isolated) | **rc=0**, `ok … 201.110s` |

The baseline `verify_go.sh` red was therefore a load-induced timeout: verifygate ran under `-race` concurrently with every other package at load average ~13. It is not a base defect in the tree. If it recurs at C5 with the same signature (`panic: test timed out after 8m0s` in host/verifygate only), the executor records it as BASE and the controller re-runs `verify_go.sh` on a quieter machine (PD5).

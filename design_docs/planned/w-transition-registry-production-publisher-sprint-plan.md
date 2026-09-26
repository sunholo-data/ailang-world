# Sprint plan — w-transition-registry-production-publisher (queue row 107)

**Design**: [w-transition-registry-production-publisher.md](w-transition-registry-production-publisher.md) (final; quorum r1 3/3 BLOCKED → r2 2/1 → carve-out + astra re-run).
**Base**: `dev b2fb92e` + the designer's prototype (14 Go files: 4 modified, 10 new) as uncommitted edits in `.planner-wt-iter195`.
**Scope**: Go host code only. **No `.ail` file changes**, so the `scripts/verify_ail.sh` constants and the S8 ledger stay as they are.
**Planner**: iter-195 sprint-planner. Every result below was observed in this session unless it is marked UNMEASURED.
**Bottom line**: 7 milestones (6 code + 1 docs). **Simulated in dependency order (M1 → M1+M2 → … → M6) on a clean copy of `b2fb92e`: every boundary passed `go vet` rc=0, compile rc=0 and targeted tests rc=0.** A few things need the controller's attention before execution:
(1) the design's milestone ORDER does not compile, so this plan reorders it (§2);
(2) two milestones go over the ~150 LOC rule. Each is a measured merge (§2);
(3) **a real prototype defect showed up.** `archive.CheckSource` returns a platform-dependent verdict: on macOS every module-declared source is refused unless it says `module entry`. There is a one-line measured fix (§4). This is the only production-code change this plan asks for.

---

## 1. Gates observed (step 2, run once in the planner worktree, pinned `AILANG_BIN=$HOME/.pinned-ailang/ailang` = AILANG v0.41.0 / Commit 24ee108)

| Gate | rc | Lines (verbatim) |
|---|---|---|
| `go vet ./...` | **0** | — |
| `go test ./... -run '^$' -count=1` (compiles every `_test.go`) | **0** | 23 × `ok … [no tests to run]` |
| `go test ./host/transitionreg ./host/archive ./host/store ./host/daemon ./cmd/world-publish ./host/broker -count=1 -timeout 600s` | **0** | `ok host/transitionreg 9.845s` · `ok host/archive 9.343s` · `ok host/store 7.220s` · `ok host/daemon 9.325s` · `ok cmd/world-publish 8.842s` · `ok host/broker 70.697s` |

Transcript: `.plan-scratch/gates.txt`. This agrees with the controller's full run on the same prototype (vet rc=0; 23 ok / 0 FAIL).

**Planner full run at the FINAL simulated state** (prototype + the §4 fix + every new test in this plan): `go vet ./...` rc=0. `go test ./... -count=1` gave 22 ok and 1 FAIL. The failure was `host/pkgproj TestQueryInterfaceReturnsWhileADescendantHoldsStdout` ("descendant pid not recorded"). That package is untouched by this row, and the run overlapped with this planner's concurrent mutation runs. Re-run in isolation it passed twice (`ok host/pkgproj 4.244s`; `ok host/pkgproj 12.546s` together with `ok host/daemon 3.368s`). **Classified as a load flake, not this row.** The executor's final full gate must still come back 23/23 on a quiet machine (`.plan-scratch/final-full.txt`).

---

## 2. Milestone partition (reordered for compile order; each boundary measured)

### Why the design's order cannot be applied as written

The design lists M1 = `publish.go` first. But `PublishSet` calls `verifySources`/`verifyLoadable`/`PublisherArchiveRequiredError` (verify.go, design M4), `verifyEpochs` (epoch.go, design M3), `buildNextRevision`/`refuseConflictingSameIDs` (collide.go, design M2), and reads `r.arch` (the `transitionreg.go` field). `verify.go` in turn calls `arch.CheckSource` (check.go, design M5).

So `publish.go` compiles only after all of these have landed. The real dependency chain is:

**check.go → (transitionreg.go field + verify.go) → epoch.go → (collide.go + publish.go) → daemon tests → CLI.**

Unused unexported functions compile in Go, and `go vet` does not flag them. That lets verify.go and epoch.go land ahead of their caller.

### Why two milestones are measured merges

- **M4 = collide.go + publish.go (187 LOC).** collide.go has no test that can run before `PublishSet` exists: every prototype test of the same-ID rule and the genesis construction goes through `PublishSet`. Landing collide.go alone would be an untested milestone (S6 vacuity). Landing it inside M2 or M3 gives 170 or 199 LOC with the same vacuity.
- **M6 = the whole CLI (268 LOC).**
  - Every prototype CLI test drives the full verb through `run`.
  - The FROZEN `flagNames` exact-set (AC24(b)), the `run` dispatch case, and the context-root pin must all land atomically with `runTransitions`.
  - A function-level split *would* compile: M6a = types + `readManifest` + `buildChanges` + `epochInList` + `pinnedTransitionFn` + requirement helpers = 152 LOC; M6b = `runTransitions` + `pinnedInterpreter` + `publishSetErrorLine` + main.go = 116 LOC. But M6a would then have **zero** killing tests. That split is **UNMEASURED** and not recommended.

LOC = `grep -vE '^\s*(//|$)'` code lines (the design's V32 rule), re-measured.

| New M | = design M | Deliverable | Prod LOC | Boundary result (simulated, pinned binary) |
|---|---|---|---|---|
| **M1** | M5 | `host/archive/check.go` (+ §4 relax fix) + AC10 driver hunk + NEW `host/archive/check_test.go` | 79 | vet 0 · compile 0 · `ok host/archive`, `ok host/verifygate`, `ok host/broker`, `ok host/store` (**after** the §4 fix; before it: `--- FAIL: TestCheckSourceUnderPinnedInterpreter`, `.plan-scratch/ms-M1.txt`) |
| **M2** | M4 | `transitionreg.go` `arch` field + `host/transitionreg/verify.go` | 85 (83 + 2) | vet 0 · compile 0 · `ok host/transitionreg 1.188s` · `ok host/store` |
| **M3** | M3 | `host/transitionreg/epoch.go` + NEW `epoch_test.go` (direct derivation test) | 116 | vet 0 · compile 0 · `ok host/transitionreg 0.991s` · `ok host/store` |
| **M4** | M1+M2 (**merge**) | `host/transitionreg/collide.go` + `publish.go` + full prototype `publish_test.go`/`verify_test.go` + AC-EPOCH-TWIN unit test | 187 (83 + 104) | vet 0 · compile 0 · `ok host/transitionreg 4.832s` · `ok host/store` |
| **M5** | M8 | `host/daemon/registry_publisher_test.go` (prototype + AC-EPOCH-TWIN daemon test) | 0 (test-only) | vet 0 · compile 0 · `ok host/daemon 31.002s` · `ok host/store` |
| **M6** | M6+M7 (**merge**) | `cmd/world-publish/transitions.go` + `main.go` hunks + `transitions_test.go` + context-root pin | 268 (256 + 12) | vet 0 · compile 0 · `ok cmd/world-publish 14.567s` · `ok host/store 9.908s` |
| **M7** | M9 | `docs/QUICKSTART.md` §6 (text in §7.3) | ~0 code / ~40 doc | attended; see §7.3 |

Runner used for the simulation: `.plan-scratch/runms.sh`. Per-boundary transcripts: `.plan-scratch/ms-M{1..6}.txt`. The runner copies the milestone's files onto a clean `git archive b2fb92e` tree and then runs:
- `go vet ./...`
- `go test ./... -run '^$'`
- the milestone's package tests.

### M1 — bounded interpreter check (design M5)

- **Files**:
  - `host/archive/check.go` (new, prototype file **plus the §4 one-line env fix**): `checkTimeout`, `checkBuffer` (+`Write`/`String`), `CheckResult`, `(*Archive).CheckSource`, `outputHead`.
  - `host/broker/registry_publish_test.go`: the **two hunks** at `:1184` (drivers map entry `"host/archive/check.go": driveArchiveCheckSource` + comment) and `:1267` (`func driveArchiveCheckSource`). They move WITH check.go, because the AC10 census reds with "has no driver" for any subprocess file without one (V31a).
  - NEW `host/archive/check_test.go` (banked: `.plan-scratch/milestones/M1/host/archive/check_test.go`).
- **Why the new test file**: the design's M5 gate `go test ./host/archive -run 'TestCheck'` matches **zero** tests in the prototype (no `TestCheck*` exists in `host/archive`), so it was vacuous. The new file has:
  - `TestCheckSourceReportsTheArchivedInterpretersVerdict`: fakes. Pass → Passed + staged bytes seen as `entry.ail`; refuse → !Passed with captured output, nil error; unarchived ref → error.
  - `TestCheckSourceUnderPinnedInterpreter`: `AILANG_BIN`-gated, `t.Fatal` when unset, never skip. Checks: standalone module → Passed "No errors found"; import-bearing → !Passed "LDR001"; garbage → !Passed.
- **Inventories at this boundary**:
  - AC10 drivers map +1 entry (file census and driver map move together).
  - The verifygate `TestSubprocessSinkGate` is green: `checkBuffer` is the mandated sync sink.
  - `contextRootPins` unchanged: CheckSource derives from the caller's ctx and has no `Background` root.
  - `TestProductionGoSurface` is green: the new file is under `host/`.
- **Tests that must pass**: the two new `TestCheckSource*`, `TestSubprocessSinkGate`, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `TestProductionContextRoots`, `TestProductionGoSurface`.
- **Mutations killed at this boundary (measured, §6)**: M1-MUT-a, M1-MUT-b, M1-MUT-c.
- **Gate**: `go vet ./host/archive/ && AILANG_BIN=… go test ./host/archive/ ./host/verifygate/ -count=1 -run 'TestCheckSource|TestSubprocessSinkGate' && AILANG_BIN=… go test ./host/broker/ -count=1 -run TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential && go test ./host/store/ -count=1`

### M2 — publisher constructor + source verification (design M4)

- **Files**:
  - `host/transitionreg/transitionreg.go`: the prototype hunks (import `host/archive`; `arch *archive.Archive` field + doc-comment lines).
  - `host/transitionreg/verify.go` (new): `CanonicalSchema`, `NewPublisher`, `PublisherArchiveRequiredError`, `TransitionSourceAbsentError`, `TransitionSourceInvalidError`, `EnsureSourceLoadable`, `verifySources`, `(*StoreReader).verifyLoadable`.
- **Test files at this boundary** (banked under `.plan-scratch/milestones/M2/`):
  - `verify_test.go` = `testFakeRelease` + `fakeInterpreterScript` + `TestEnsureSourceLoadableStagesHermetically`.
  - `publish_test.go` = `TestCanonicalSchemaExportMatchesValidate` only (imports `testing`).
- **Must pass**: `TestEnsureSourceLoadableStagesHermetically`, `TestCanonicalSchemaExportMatchesValidate`, and the whole existing `host/transitionreg` suite.
- **Mutations killed here**: MUT-8 core arm (`EnsureSourceLoadable` ignores `!result.Passed` → `--- FAIL: TestEnsureSourceLoadableStagesHermetically`, measured). MUT-2 and the `PublishSet` arm of MUT-8 are killed at M4. `verifySources`/`verifyLoadable` have no caller until M4.
- **Gate**: `go vet ./host/transitionreg/ && go test ./host/transitionreg/ -count=1 -run 'TestEnsure|TestCanonicalSchema' && go test ./host/transitionreg/ ./host/store/ -count=1`

### M3 — epoch derivation/validation (design M3)

- **Files**: `host/transitionreg/epoch.go` (new): `EpochRegistryAbsentError`, `InterpreterEpochMismatchError`, `releaseFromManifest`, `(*StoreReader).EpochsForInterpreter`, `epochRegistry`, `verifyEpochs`, `epochIn`.
- **Test files at this boundary** (banked `.plan-scratch/milestones/M3/`):
  - `verify_test.go` = M2 content + `publisherStore` + `TestEpochsForInterpreterMatchesBootstrapRelease` + `TestReleaseFromManifestMirrorsDaemonReduction`. Imports gain `hashref`, `registry`, `store`.
  - `publish_test.go` unchanged from M2.
  - NEW `epoch_test.go` with `TestEpochsForInterpreterRefusesAbsentAndUnnominated`, planner-written (text in §5.3). **Why**: without it, M3's boundary kills no mutation from the table. The prototype's refusal tests all need `PublishSet`, and MUT-10 **survives** `TestEpochsForInterpreterMatchesBootstrapRelease`, because a one-epoch registry returns `[1]` either way.
- **Must pass**: the three tests above + `TestEnsureSourceLoadableStagesHermetically` + the existing suite.
- **Mutations killed here (measured)**:
  - MUT-10 (`if candidate == release || true`) → `--- FAIL: TestEpochsForInterpreterRefusesAbsentAndUnnominated`.
  - MUT-11 absent-head arm (synthetic epoch-1 registry) → same test fails.
  - MUT-9 is killed at M4 (its call site is `PublishSet`).
- **Gate**: `go vet ./host/transitionreg/ && go test ./host/transitionreg/ -count=1 -run 'TestEpochs|TestRelease' && go test ./host/transitionreg/ ./host/store/ -count=1`

### M4 — PublishSet + genesis/same-ID rule (design M1+M2, measured merge)

- **Files**:
  - `host/transitionreg/collide.go` (new): `EmptyGenesisError`, `SameIDConflictError`, `refuseConflictingSameIDs`, `descriptorBytes`, `buildNextRevision`.
  - `host/transitionreg/publish.go` (new): `SetResult`, `CurrentRevision`, `PublishSet`, `sameEntries`.
- **Test files at this boundary** (banked `.plan-scratch/milestones/M4/`):
  - `publish_test.go` and `verify_test.go` = the prototype files **verbatim**.
  - `epoch_test.go` = the M3 test + **AC-EPOCH-TWIN unit test** `TestPublishSetEpochTwinInterpretersKeepDistinctPins`. Import adds `hashref`. Text in §5.1.
- **Must pass**:
  - `TestGenesisCannotUseBuildNextWithZeroExpectedHead`
  - `TestPublishSetGenesisThenGrowth`
  - `TestPublishSetIdempotentRepublishIsNoOp`
  - `TestPublishSetRefusesAbsentTransitionSource`
  - `TestPublishSetRefusesEmptyAndRemovalOnlyGenesis`
  - `TestPublishSetCASRetryMergesWinner`
  - `TestPublishSetSecondConflictSurfacesTypedError`
  - `TestPublishSetSameIDConflictOnCASRetryRefuses`
  - `TestPublishSetCASRetrySameIDIdenticalBytesIsNoOp`
  - `TestPublishSetRefusesWithoutArchive`
  - `TestPublishSetRefusesEpochNotNominatingTheInterpreter`
  - `TestPublishSetRefusesDefaultEpochOne`
  - `TestPublishSetRefusesAbsentEpochRegistry`
  - `TestPublishSetRefusesUnloadableTransitionSource`
  - `TestPublishSetEpochTwinInterpretersKeepDistinctPins`
- **Mutations killed here**: MUT-1, MUT-2 (unit), MUT-3 (unit), MUT-4, MUT-7, MUT-8 (unit), MUT-9 (unit ×2), MUT-10 (unit), MUT-11 (unit ×3), MUT-12, **TWIN-MUT** (measured, §5).
- **Gate**: `go vet ./host/transitionreg/ && go test ./host/transitionreg/ ./host/store/ -count=1`

### M5 — daemon acceptance (design M8, test-only)

- **Files**: `host/daemon/registry_publisher_test.go` = prototype + appended `TestEpochTwinInterpretersBothListedOnCard` (banked `.plan-scratch/milestones/M5/…`; text in §5.2). It moves **before** the CLI because it depends only on the M4 core. That is the P9 seam: the daemon test publishes through `NewPublisher(…).PublishSet`, not through the verb.
- **Must pass**: `TestPublishedTransitionsAppearOnAgentCard`, `TestPublisherRefusalKeepsCardEmpty`, `TestEpochRefusalKeepsCardUnchanged`, `TestEpochTwinInterpretersBothListedOnCard`. The recorder harness is used: no sockets, no wall-clock asserts.
- **Mutations killed here**: MUT-2 (daemon), MUT-9 (daemon), TWIN-MUT (daemon, measured). Also the row's "card ignores the capability filter" regression, via the negative arm.
- **Gate**: `go vet ./host/daemon/ && go test ./host/daemon/ -count=1 -run 'TestPublish|TestEpochRefusal|TestEpochTwin|TestCard' && go test ./host/daemon/ ./host/store/ -count=1`

### M6 — the operator verb (design M6+M7, measured merge)

- **Files**:
  - `cmd/world-publish/transitions.go` (new, all of it).
  - `cmd/world-publish/main.go`: all prototype hunks, which are atomic together:
    - 3 `options` fields
    - **`flagNames` +3** (`ailang-bin`, `interpreter-ref`, `manifest`: 15 → 18 entries, still sorted)
    - 3 `fs.StringVar`
    - 1 usage line
    - the `case "transitions"` dispatch
  - `cmd/world-publish/transitions_test.go` (new, verbatim).
  - `host/store/context_roots_test.go`: **+1 pin** `"cmd/world-publish/transitions.go|runTransitions|Background": 1` (15 → 16 entries).
- **Why these move together**:
  - Without the pin, `TestProductionContextRoots` reds (V15).
  - With the pin but no `runTransitions`, the `DeepEqual` also reds.
  - Without the flagNames move, AC24(b) reds.
- **Must pass**:
  - `TestTransitionsVerbHappyPathAndIdempotence`
  - `TestTransitionsVerbEpochArms` (3 subtests)
  - `TestTransitionsVerbRefusesGarbageSourceBeforePutObject`
  - `TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter` (AILANG_BIN-gated)
  - `TestPublishErrorLineNamesConflictingID`
  - `TestTransitionsVerbFences` (6 subtests)
  - the whole existing package, including AC24(b) and `TestWorldCoreManifestMatchesTheCommittedGolden`
  - `TestProductionContextRoots`, `TestProductionGoSurface`
- **Mutations killed here**: MUT-3 (CLI), MUT-5 ×2, MUT-6, MUT-8 (CLI), MUT-10 (CLI), MUT-11 (CLI). The design's V29 kill list is measured on this exact code; the planner did not re-run it.
- **Gate**: `go vet ./cmd/world-publish/ && AILANG_BIN=… go test ./cmd/world-publish/ ./host/store/ -count=1`, then the FULL gate: `go vet ./... && AILANG_BIN=… go test ./... -count=1` → 23 ok / 0 FAIL, and `AILANG_BIN=… ./scripts/verify_ail.sh`, which must be unchanged-green because no `.ail` is touched.

### M7 — QUICKSTART §6 (design M9 / Residual 5): exact text in §7.3.

---

## 3. Deviations: prototype vs design text (one line each)

1. **Milestone order** runs check → verify → epoch → publish+collide → daemon → CLI, not the design's publish-first order. publish.go does not compile before its callees exist (§2, measured).
2. **M4 merges design M1+M2 (187 LOC) and M6 merges design M6+M7 (268 LOC)**, over the ~150 LOC rule. The alternatives either leave a milestone with no killing test (S6) or need stub code the prototype does not have (§2).
3. **The design's M5 gate `-run 'TestCheck'` ran zero tests in `host/archive`.** The planner adds `check_test.go` (2 tests, measured) so the gate is non-vacuous.
4. **`CheckSource` stages `entry.ail`, not the capsule's `host/capsule/main.ail`** (capsule.go:32). The design's "capsule shape" wording is therefore about *hermetic import resolution* only, not the module path. That matches Residual 8's disclaimer: invocation compatibility is not established. The capsule also runs with env = `AILANG_FS_SANDBOX` only, while the check runs with the scrubbed parent env.
5. **The check's verdict is platform-dependent without the §4 fix.** This is measured, the design did not observe it, and V26(c)'s "module-declared → exit 0" was only true from an interactive shell whose `$PWD` kept the `/var/folders` alias.
6. **`transitions.go` has 256 code LOC (V32 agrees); the design's M6/M7 "~130/~126" split is not a compile-clean cut.** The shell-first half references the materialisation functions.
7. **Design M6's acceptance says "missing manifest usage"; the prototype's subtest is `absent manifest is usage`.** Same behaviour.
8. **Design Residual 5 says "QUICKSTART row (M6) … same commit as M3/M4".** The milestone numbers are stale after the r2 split. This plan puts it in M7, after the verb exists.
9. **The prototype's new files are not `gofmt`-clean** (trailing-newline and alignment drift: `gofmt -l` lists 11 of the 14 files, plus the pre-existing `host/store/schema_version_test.go`). CI does not gate gofmt (`.github/workflows/ci.yml` has no fmt step). The executor may run `gofmt -w` on the NEW files. This is whitespace only, but it realigns every line of the `contextRootPins` map, so leave `context_roots_test.go` as the +1-line hunk.
10. **Design Residual 4 (`--dry-run` for the verb)** is left out of this sprint. `dry-run` is already a frozen flag, but `runTransitions` ignores it. Adding it is an executor-judgement follow-up, not in the prototype, and not planned.

---

## 4. Finding: `CheckSource`'s verdict depends on the host path. Fix it in M1 (one line).

**Observed.** With the prototype `check.go` unchanged, the planner's `TestCheckSourceUnderPinnedInterpreter` failed at the M1 boundary:
```
--- FAIL: TestCheckSourceUnderPinnedInterpreter
    standalone module: check = {Output:… Error: Error MOD010: module 'check/entry' doesn't match file path 'entry'.
      Fix: use --relax-modules flag or set AILANG_RELAX_MODULES=1 … Passed:false}
```

**Root cause (measured by hand with the pinned v0.41.0):**
- `os.MkdirTemp("")` gives `/var/folders/…/T/world-source-check-*`.
- The child process's cwd resolves to `/private/var/folders/…`. `exec` does not pass a matching `$PWD`, so the interpreter's MOD010 "auto-relaxed for temporary directory" rule does not fire.
- Fresh dir under `/private$TMPDIR`: `check` of `module check/entry` gives **rc=1 (MOD010)**, and so does `module host/capsule/main`. Only `module entry` gives rc=0.
- From an interactive shell `cd`'d via the `/var/folders` alias, the same file gives rc=0. That is how V26(c) was observed.
- **Caution**: a second `check` in the same dir can pass from the `.ailang` cache. Always probe in a fresh dir.
- With `AILANG_RELAX_MODULES=1`, `check/entry`, `host/capsule/main` and `entry` all give rc=0. The import-bearing source still gives **rc=1 LDR001**, so the hermetic import check keeps its teeth.
- On Linux CI (`/tmp`, no alias) the prototype probably passes. That is the defect: **the same bytes and the same interpreter get a different verdict on each OS**, and Mark's Mac would refuse every realistic QUICKSTART source.

**Fix (the only production-code change in this plan; banked `.plan-scratch/check-relax.diff`):**
```go
-	cmd.Env = childenv.Scrubbed(os.Environ())
+	// AILANG_RELAX_MODULES=1: the staged file is entry.ail, so a declared
+	// module path never matches it; without the relax flag the verdict depends
+	// on whether the interpreter recognises the scratch root as a temp dir
+	// (macOS resolves it under /private/var and does not) — measured, iter-195.
+	cmd.Env = append(childenv.Scrubbed(os.Environ()), "AILANG_RELAX_MODULES=1")
```

**Measured after the fix:**
- `ok host/archive`, `ok host/broker` (AC10 still scrubs the credential), `ok host/verifygate`.
- The full M1–M6 simulation stayed green.
- Removing the line again reds `TestCheckSourceUnderPinnedInterpreter` (M1-MUT-c).

**Why relax rather than staging at `host/capsule/main.ail`:** staging at the capsule path would make publication *require* `module host/capsule/main`. That would decide row 106's source/entry convention (Residuals 1 and 8), which this row explicitly does not own. Relaxing the module-path match keeps the check about parse, types and hermetic imports, which is exactly its documented scope.

**Controller decision needed**: accept this one-line deviation, or send it back to the designer. The planner recommends accepting it. It is measured, it is the interpreter's own suggested remedy, and it adds no authority surface.

---

## 5. AC-EPOCH-TWIN: prototyped and measured

**Mutation to kill** (the design's suggestion, implemented concretely):
- In `verifyEpochs`, keep `firstByRelease map[string]hashref.HashRef`.
- When a descriptor's derived release was already seen with a different interpreter, overwrite `d.Interpreter` with the first one.
- This "publisher substitutes the nominated release's first archived interpreter for the descriptor's own pin".

**Result:**
- **SURVIVES the entire current prototype suite.** `go test ./host/transitionreg ./host/daemon ./cmd/world-publish -skip Twin` gives `ok` ×3. No existing test ever publishes two interpreters in one set.
- **KILLED by both new tests:**
  - `--- FAIL: TestPublishSetEpochTwinInterpretersKeepDistinctPins`: "published pins = (sha256:1a62…, sha256:1a62…), want each descriptor's OWN interpreter (sha256:1a62…, sha256:bfdc…)".
  - `--- FAIL: TestEpochTwinInterpretersBothListedOnCard`: "stored entries … tools.twin Interpreter:sha256:43f8… want … pinned to sha256:aec4…".
- **The new tests are load-bearing.** The first daemon draft published the twins in two separate `PublishSet` calls, and that draft did *not* kill the mutation. The final text publishes both in ONE set, which was re-measured as killed.
- Both tests pass on the unmutated code (`--- PASS` ×2, observed with `-v`).
- **Arm (b)**, a non-nominated release, is refused with `*InterpreterEpochMismatchError` (Release named, Nominating empty). Head ref/revision and entry count are unchanged, and at daemon level the card still lists exactly 2 skills.

Assigned to **M4** (unit) and **M5** (daemon).

### 5.1 Unit test: append to `host/transitionreg/epoch_test.go` at M4 (add import `github.com/sunholo-data/ailang-world/host/hashref`)

```go
// TestPublishSetEpochTwinInterpretersKeepDistinctPins is AC-EPOCH-TWIN
// (quorum r2, gpt6-astra; adopted in the form that pins D1's true
// behaviour). The epoch check is ADVISORY release nomination: two archived
// interpreters with DIFFERENT hashes and an IDENTICAL first --version line
// are both epoch-eligible under the same nomination — and each published
// descriptor keeps its OWN Interpreter HashRef (the authoritative pin, D1),
// never a release-level stand-in, so card, invocation and replay can never
// conflate them. (b) A release the registry does not nominate is refused
// with the typed mismatch and the head is unchanged.
func TestPublishSetEpochTwinInterpretersKeepDistinctPins(t *testing.T) {
	const release = "TWIN-FAKE v1"
	s, arch, first := publisherStore(t, release, release)
	second, err := arch.Archive(fakeInterpreterScript(t, release+"\nCommit: twin-second", 0))
	if err != nil {
		t.Fatalf("archive twin interpreter: %v", err)
	}
	if first == second {
		t.Fatal("premise: the twin interpreters must have different hashes")
	}
	ctx := context.Background()
	pub := NewPublisher(s, arch)
	for _, ref := range []hashref.HashRef{first, second} {
		epochs, got, err := pub.EpochsForInterpreter(ctx, ref)
		if err != nil || got != release || len(epochs) != 1 || epochs[0] != 1 {
			t.Fatalf("EpochsForInterpreter(%s) = (%v, %q, %v), want ([1], %q)", ref, epochs, got, err, release)
		}
	}

	alpha := storedSourceDescriptor(t, s, first, "tools.alpha")
	beta := storedSourceDescriptor(t, s, second, "tools.beta")
	res, err := pub.PublishSet(ctx, []Change{{ID: alpha.ID, Descriptor: &alpha}, {ID: beta.ID, Descriptor: &beta}})
	if err != nil || res.Revision != 1 || res.Unchanged {
		t.Fatalf("twin publish = (%+v, %v), want revision 1 with both entries", res, err)
	}
	snap, err := NewReader(s).ReadSnapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := snap.List()
	if len(got) != 2 || got[0].ID != "tools.alpha" || got[1].ID != "tools.beta" {
		t.Fatalf("twin entries = %+v, want tools.alpha and tools.beta", got)
	}
	if got[0].Interpreter != first || got[1].Interpreter != second {
		t.Fatalf("published pins = (%s, %s), want each descriptor's OWN interpreter (%s, %s): the epoch nominates a release, the HashRef stays authoritative",
			got[0].Interpreter, got[1].Interpreter, first, second)
	}

	other, err := arch.Archive(fakeInterpreterScript(t, "TWIN-OTHER v2", 0))
	if err != nil {
		t.Fatal(err)
	}
	gamma := storedSourceDescriptor(t, s, other, "tools.gamma")
	_, err = pub.PublishSet(ctx, []Change{{ID: gamma.ID, Descriptor: &gamma}})
	var mismatch *InterpreterEpochMismatchError
	if !errors.As(err, &mismatch) || mismatch.ID != "tools.gamma" || mismatch.Release != "TWIN-OTHER v2" || len(mismatch.Nominating) != 0 {
		t.Fatalf("non-nominated release = %v, want *InterpreterEpochMismatchError for tools.gamma with no nominating epoch", err)
	}
	after, err := NewReader(s).ReadSnapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if after.Head != snap.Head || after.Revision != 1 || len(after.List()) != 2 {
		t.Fatalf("after refusal: head %s rev %d entries %d, want unchanged (%s, 1, 2)", after.Head, after.Revision, len(after.List()), snap.Head)
	}
}
```

### 5.2 Daemon test: append to `host/daemon/registry_publisher_test.go` at M5 (no new imports)

The full file is banked at `.plan-scratch/milestones/M5/host/daemon/registry_publisher_test.go`. The appended function:

```go
// TestEpochTwinInterpretersBothListedOnCard is AC-EPOCH-TWIN at daemon
// level: a second archived interpreter whose bytes differ from the daemon's
// own but whose first --version line is the same release is epoch-eligible
// under the daemon's bootstrapped nomination; both published skills appear on
// the card and each stored descriptor keeps its own Interpreter pin. A
// non-nominated release is refused and head + card stay unchanged.
func TestEpochTwinInterpretersBothListedOnCard(t *testing.T) {
	d, dbPath := newPublisherDaemon(t)
	arch := archive.New(dbPath)
	own := daemonInterpreterRef(t, d)
	writeFake := func(name, versionOut string) hashref.HashRef {
		t.Helper()
		path := filepath.Join(t.TempDir(), name)
		script := "#!/bin/sh\ncase \"$1\" in\n  --version) printf '" + versionOut + "';;\nesac\nexit 0\n"
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
		ref, err := arch.Archive(path)
		if err != nil {
			t.Fatalf("archive %s: %v", name, err)
		}
		return ref
	}
	twin := writeFake("daemon-twin-ailang", daemonFakeRelease+"\\nCommit: twin\\n")
	if twin == own {
		t.Fatal("premise: the twin interpreter must hash differently from the daemon's own")
	}
	// pinned is one honest descriptor pinning interp; its source object is stored.
	pinned := func(id string, interp hashref.HashRef) transitionreg.Change {
		payload := []byte("transition source for " + id)
		if err := d.store.PutObject(store.Object{
			Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("daemon-test/transition-source")),
			SemanticID: "daemon-test/transition-source", Provenance: "registry_publisher_test", Payload: payload,
		}); err != nil {
			t.Fatal(err)
		}
		desc := transitionreg.Descriptor{
			ID: id, TransitionFn: hashref.SumSHA256(payload), Interpreter: interp, SemanticsEpoch: 1,
			InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`),
			Access:          transitionreg.EffectRequirement{Effect: "world.apply", Scope: "world", Cost: 1},
			DeclaredEffects: []transitionreg.EffectRequirement{{Effect: "world.apply", Scope: "world", Cost: 1}},
			Title:           "title-" + id,
		}
		return transitionreg.Change{ID: id, Descriptor: &desc}
	}
	publish := func(changes ...transitionreg.Change) error {
		_, err := transitionreg.NewPublisher(d.store, arch).PublishSet(context.Background(), changes)
		return err
	}
	// Both twins in ONE publish: the set is where a release-level stand-in
	// for the descriptor's own pin would conflate them.
	if err := publish(pinned("tools.own", own), pinned("tools.twin", twin)); err != nil {
		t.Fatalf("publish on the daemon's own interpreter and its same-release twin: %v", err)
	}
	snap, err := transitionreg.NewReader(d.store).ReadSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	entries := snap.List()
	if len(entries) != 2 || entries[0].Interpreter != own || entries[1].Interpreter != twin {
		t.Fatalf("stored entries = %+v, want tools.own pinned to %s and tools.twin pinned to %s", entries, own, twin)
	}
	auth := mintSessionGrants(t, d, "ep-twin", "world.apply")
	if card := cardSkillIDs(t, d, "Bearer "+auth); len(card.Skills) != 2 {
		t.Fatalf("card = %+v, want both twin-pinned skills", card.Skills)
	}

	other := writeFake("daemon-other-ailang", "DAEMON-OTHER v2\\n")
	var mismatch *transitionreg.InterpreterEpochMismatchError
	if err := publish(pinned("tools.other", other)); !errors.As(err, &mismatch) {
		t.Fatalf("non-nominated release = %v, want *InterpreterEpochMismatchError", err)
	}
	after, err := transitionreg.NewReader(d.store).ReadSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if after.Head != snap.Head || len(after.List()) != 2 {
		t.Fatalf("head after refusal = %s (%d entries), want unchanged %s (2)", after.Head, len(after.List()), snap.Head)
	}
	if card := cardSkillIDs(t, d, "Bearer "+auth); len(card.Skills) != 2 {
		t.Fatalf("card after refusal = %+v, want unchanged (2 skills)", card.Skills)
	}
}
```

### 5.3 M3's direct derivation test: new file `host/transitionreg/epoch_test.go` at M3 (imports `context`, `errors`, `testing`)

```go
// TestEpochsForInterpreterRefusesAbsentAndUnnominated is M3's direct pin of
// the derivation (before PublishSet exists): an interpreter whose release no
// epoch nominates derives ZERO epochs even though epoch 1 exists (no
// "every epoch" / default-to-1 derivation — MUT-10), and a store with no
// epoch-registry head is the typed EpochRegistryAbsentError, never a
// synthetic epoch-1 registry (MUT-11's absent-head arm).
func TestEpochsForInterpreterRefusesAbsentAndUnnominated(t *testing.T) {
	s, arch, ref := publisherStore(t, testFakeRelease, "OTHER-RELEASE v1")
	epochs, release, err := NewPublisher(s, arch).EpochsForInterpreter(context.Background(), ref)
	if err != nil {
		t.Fatalf("EpochsForInterpreter: %v", err)
	}
	if release != testFakeRelease || len(epochs) != 0 {
		t.Fatalf("derived (%v, %q), want no epoch for release %q (epoch 1 nominates another release)", epochs, release, testFakeRelease)
	}

	absentStore, absentArch, absentRef := publisherStore(t, testFakeRelease, "")
	_, _, err = NewPublisher(absentStore, absentArch).EpochsForInterpreter(context.Background(), absentRef)
	if !errors.As(err, new(*EpochRegistryAbsentError)) {
		t.Fatalf("absent epoch registry = %v, want *EpochRegistryAbsentError (never a default)", err)
	}

	if _, _, err := NewReader(s).EpochsForInterpreter(context.Background(), ref); !errors.As(err, new(*PublisherArchiveRequiredError)) {
		t.Fatalf("bare-reader derivation = %v, want *PublisherArchiveRequiredError", err)
	}
}
```

`host/archive/check_test.go` (M1) is banked verbatim at `.plan-scratch/milestones/M1/host/archive/check_test.go`. Its behaviour is described in §2 M1.

---

## 6. Mutation table per milestone

"Measured here" means this planner executed it (mutate → targeted test → revert). "V29" means it was measured by the designer on this exact code and not re-run by the planner.

| M | ID | Mutation | Killer test(s) | Status |
|---|---|---|---|---|
| M1 | M1-MUT-a | `CheckSource` reports a refusal as `Passed: true` | `TestCheckSourceReportsTheArchivedInterpretersVerdict`, `TestCheckSourceUnderPinnedInterpreter` (+ later `TestEnsureSourceLoadableStagesHermetically`, `TestPublishSetRefusesUnloadableTransitionSource`) | **KILLED, measured here** |
| M1 | M1-MUT-b | check child env not scrubbed (`os.Environ()` leaks the registry credential) | `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` | **KILLED, measured here** |
| M1 | M1-MUT-c | the §4 relax line removed (= the prototype as-is) | `TestCheckSourceUnderPinnedInterpreter` (macOS; probably survives on Linux `/tmp`, UNMEASURED there) | **KILLED on this rig, measured here** |
| M2 | MUT-8 (core) | `EnsureSourceLoadable` ignores `!Passed` | `TestEnsureSourceLoadableStagesHermetically` | **KILLED, measured here** |
| M3 | MUT-10 | derivation ignores the release (every epoch returned) | `TestEpochsForInterpreterRefusesAbsentAndUnnominated` | **KILLED, measured here** (survives the prototype's M3-reachable tests) |
| M3 | MUT-11 (absent arm) | absent epoch-registry head → synthetic epoch-1 registry | `TestEpochsForInterpreterRefusesAbsentAndUnnominated` | **KILLED, measured here** |
| M4 | MUT-1 | skip the captured-head CAS | `TestPublishSetGenesisThenGrowth` | V29 |
| M4 | MUT-2 (unit) | both presence arms removed | `TestPublishSetRefusesAbsentTransitionSource` | V29 |
| M4 | MUT-3 (unit) | unchanged early-return removed | `TestPublishSetIdempotentRepublishIsNoOp` | V29 |
| M4 | MUT-4 | genesis `Revision: 2` | `TestPublishSetGenesisThenGrowth` | V29 |
| M4 | MUT-7 | genesis via `BuildNext` over rev 0 | `TestPublishSetGenesisThenGrowth`, `TestPublishSetIdempotentRepublishIsNoOp` | V29 |
| M4 | MUT-8 (PublishSet) | `verifyLoadable` call removed (+ the CLI arm at M6) | `TestPublishSetRefusesUnloadableTransitionSource` | V29 |
| M4 | MUT-9 | `verifyEpochs` call removed | `TestPublishSetRefusesEpochNotNominatingTheInterpreter`, `TestPublishSetRefusesDefaultEpochOne` (+ M5 daemon arm) | V29; controller re-measured the `return nil` variant (4 reds) |
| M4 | MUT-10 (unit) | as above | `TestPublishSetRefusesDefaultEpochOne` | V29 |
| M4 | MUT-11 | default-to-1 | `TestPublishSetRefusesDefaultEpochOne`, `…EpochNotNominating…`, `…AbsentEpochRegistry` | V29 |
| M4 | MUT-12 | `refuseConflictingSameIDs` call removed | `TestPublishSetSameIDConflictOnCASRetryRefuses` | V29 |
| M4 | **TWIN-MUT** | first same-release interpreter substituted for the descriptor's own pin | `TestPublishSetEpochTwinInterpretersKeepDistinctPins` | **KILLED, measured here; SURVIVES the prototype suite** |
| M5 | MUT-2 (daemon) | as above | `TestPublisherRefusalKeepsCardEmpty` | V29 |
| M5 | MUT-9 (daemon) | as above | `TestEpochRefusalKeepsCardUnchanged` | V29 |
| M5 | **TWIN-MUT** | as above | `TestEpochTwinInterpretersBothListedOnCard` | **KILLED, measured here** |
| M5 | row example | card ignores the capability filter (landed row-40 code) | `TestPublishedTransitionsAppearOnAgentCard` negative arm | regression killer (design) |
| M6 | MUT-3 (CLI) | as above | `TestTransitionsVerbHappyPathAndIdempotence` | V29 |
| M6 | MUT-5 | skip `requireAttendedOperator` | `TestTransitionsVerbFences/ci_environment_stops`, `/wrong_phrase_stops` | V29 |
| M6 | MUT-6 | unverified `--interpreter-ref` accepted | `TestTransitionsVerbFences/unverifiable_interpreter_ref_refused` | V29 |
| M6 | MUT-8 (CLI) | `EnsureSourceLoadable` call removed from `pinnedTransitionFn` | `TestTransitionsVerbRefusesGarbageSourceBeforePutObject` | V29 |
| M6 | MUT-10/11 (CLI) | as above | `TestTransitionsVerbEpochArms` | V29 |

**Executor obligation**: re-run MUT-1…MUT-12 and TWIN-MUT once against the landed files, and record the `--- FAIL:` lines in the sprint doc's Verification Log. Tally target: **12/12 + TWIN + M1-MUT-a/b/c**. Two notes on writing mutations: a mutation that leaves an unused variable or import is killed by the compiler, not a test, so do not count it (S6). Examples: `if true` over `candidate` needs `|| true`, and dropping `childenv.Scrubbed` needs a retained reference.

---

## 7. Execution notes for the executor

### 7.1 Inputs (the controller banks these OUTSIDE the worktree, with the prototype tarball)

- The prototype's 14 files, from the controller's tarball.
- `.plan-scratch/check-relax.diff`: the §4 one-line fix to `host/archive/check.go`.
- `.plan-scratch/milestones/M{1..5}/…`: the exact per-milestone test-file states. These are `_test.go` files only.
  - M2 and M3 hold trimmed `verify_test.go`/`publish_test.go` with correct imports.
  - M4's are the prototype verbatim, plus `epoch_test.go`.
  - M5's daemon test is prototype + twin.
- `.plan-scratch/runms.sh`: the boundary runner used for the simulation. Paths must be adapted.

### 7.2 Procedure

- Work in a fresh worktree off `dev`. Apply M1 … M7 in order. At each boundary run that milestone's gate from §2 with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (mandatory: the CLI, archive and broker arms `t.Fatal` without it).
- Use `go vet` / `go test -run '^$'` as the compile fence, **never `go build`** (it skips `_test.go`), and never `verify_go.sh`.
- **Snapshot each milestone OUTSIDE the worktree**, e.g. `~/.ailang/state/world-iter195/exec-snap/M<n>/`. Never snapshot to an in-tree `.snap/`: `TestProductionGoSurface` walks every non-test `.go` file in the tree (it skips only `.git`/`vendor`) and reds on copies. This planner hit the same constraint and deleted its non-test scratch copies before finishing. Note that `/tmp` is wiped on reboot.
- Commit per milestone, or once at the end: executor's choice. Either way the final tree must pass `go vet ./...`, `AILANG_BIN=… go test ./... -count=1` (23 ok / 0 FAIL; re-run a lone `host/pkgproj` descendant-pid failure in isolation before calling it red), and `AILANG_BIN=… ./scripts/verify_ail.sh` (unchanged).
- **Do not** touch `tools/launchd/*`, any `.ail`, or `scripts/verify_ail.sh` constants.

### 7.3 M7: exact QUICKSTART text

Insert this as a new section between "## 5. Stop" and the closing `---` block of `docs/QUICKSTART.md`.

The section is **attended by construction**: `world-publish transitions` and `ailang-worldd session mint` both require a controlling terminal and a typed confirmation. **The executor must NOT drive those fences with a pty wrapper (`script`, `expect`, …)**. That is precisely the unattended-publish path MUT-5 exists to forbid. The executor runs only the non-TTY steps (building, writing the source and manifest, `serve`, `curl` of the unauthenticated-refusal case if any). The section header stays marked *pending attended run* until Mark executes it verbatim. At that point he removes the marker (S7: unexecuted text is a claim, not a quickstart).

```markdown
## 6. Publish a transition skill to the A2A card *(attended — pending first verbatim run)*

The transition registry is written by an **attended, local** operator verb, never by the daemon
and never over the network. It needs single-writer authority, so the daemon must be **stopped**
(step 5). Capture the daemon's pinned interpreter first — every published descriptor pins it:

```bash
/tmp/ailang-worldd health   # before step 5: note "interpreter_ref"
```

A publishable source must be a **self-contained** module: publication runs the pinned
interpreter's `check` on it in an empty scratch root, so a module importing `world/*` is refused
with `LDR001` (that is why no landed `world/*.ail` is publishable yet).

```bash
cat > /tmp/echo.ail <<'EOF'
module quickstart/echo

export func echo(x: string) -> string { x }
EOF
cat > /tmp/transitions.json <<'EOF'
[{"id": "tools.echo", "title": "Echo", "description": "quickstart transition",
  "transitionFnFile": "/tmp/echo.ail",
  "inputSchema": {"type": "object"}, "outputSchema": {"type": "object"},
  "access": {"effect": "world.apply", "scope": "world", "cost": 1},
  "declaredEffects": [{"effect": "world.apply", "scope": "world", "cost": 1}]}]
EOF
go build -o /tmp/world-publish ./cmd/world-publish
/tmp/world-publish transitions --store /tmp/world-demo.db --manifest /tmp/transitions.json \
  --interpreter-ref <interpreter_ref from health>
```

Type the confirmation phrase when asked. Output: `semantics epoch 1 derived from
world/epoch-registry/v1 for interpreter release "…"` then `published transition registry revision
1 (head sha256:…)`. Running it again prints `transition registry UNCHANGED at revision 1` — an
identical republish writes nothing. `semanticsEpoch` is omitted on purpose: it is derived from the
epoch registry the daemon bootstrapped, never defaulted.

Mint a session that holds the skill's capability, restart the daemon, and read the card:

```bash
/tmp/ailang-worldd session mint --db /tmp/world-demo.db --episode quickstart \
  --grant world.apply=world:10 --out /tmp/qs-session
/tmp/ailang-worldd serve --db /tmp/world-demo.db --ailang-bin /tmp/ailang-v0300/ailang &
curl -s -H "Authorization: Bearer $(cat /tmp/qs-session)" http://127.0.0.1:7644/.well-known/agent.json
```

The card lists `tools.echo` / `Echo`. A session minted without `world.apply` gets the same 200
with **zero** skills — the card is capability-filtered per session. What publication does and does
not establish: the source object exists and passes standalone `check` under the pinned
interpreter, and its epoch is one the registry nominates for that interpreter's release; it does
**not** establish that the transition can be invoked (`/a2a/` refuses every invocation until the
invocation coordinator lands).
```

Notes for whoever runs it:
- (a) `--ailang-bin` must be the same binary the daemon pinned, so the epoch release matches. The earlier sections of the quickstart still say `/tmp/ailang-v0300/ailang`, which is stale since the pin moved to `$HOME/.pinned-ailang/ailang`. That is a pre-existing S7 drift outside this row. Mark may substitute the current pin in *both* places when he runs it.
- (b) `--out` writes the credential at mode 0600. Confirm the `session mint` token file format is exactly the Bearer token before relying on `$(cat …)`. This is UNMEASURED by this planner.
- (c) The `module quickstart/echo` source was measured to `check` rc=0 under the pinned v0.41.0 **only with the §4 fix** when staged by `CheckSource` on macOS. Without it, it is refused with MOD010.

### 7.4 Sprint-doc bookkeeping the executor owns

- Residual 7 stays true.
- Update Residual 5's milestone reference to M7.
- Record in the Verification Log:
  - the §4 finding, with the fresh-dir probe transcript (strict rc=1 / relax rc=0 / LDR001 still rc=1);
  - the TWIN measurement;
  - the full-gate numbers.
- The Decision-for-Mark (A/B phrase) is unchanged. The row lands under A.

---

## 8. Planner-session artifacts (`.plan-scratch/`, `_test.go` / text only; every non-test `.go` copy was deleted before finishing)

`gates.txt` (step-2 gates) · `ms-M1.txt` … `ms-M6.txt` (per-boundary simulation; `ms-M1.txt` is the pre-fix FAIL) · `final-full.txt` (final-state full run incl. the pkgproj load flake) · `check-relax.diff` · `milestones/…` (per-milestone test files) · `runms.sh`.

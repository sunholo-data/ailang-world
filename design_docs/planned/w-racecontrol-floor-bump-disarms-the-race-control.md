# w-racecontrol-floor-bump-disarms-the-race-control

**Status**: Planned
**Target**: current iteration (World iter-137 / P48)
**Priority**: P1
**Estimated**: ~0.1d (~1h)
**Dependencies**: None
**Planner-Lane**: codex-ok

## Problem Statement

`design_docs/verification/w-race-gate-blindspot/racecontrol/` is the second of the repo's
three nested `go.mod` modules, and the only one that hosts a **known-positive control** for
`verify_go.sh`'s race-detector gate. `scripts/verify_go.sh:229` runs
`cd .../racecontrol && go run -race .`, and `:232` requires `WARNING: DATA RACE` in the output
or it FATALs `"the race detector is not armed; every 0-races result in this gate is void"`.
The module deliberately contains a data race (`main.go`: 4 goroutines incrementing `shared`)
so the control is supposed to fire on every run.

That control is **disarmed by a blanket module-floor bump**, the same defect class queue row
42 just fixed on the *repro* module. Row 42's own Systemic-Issue Audit row **V17** examined
`racecontrol/` and declared it a *non-instance* — *"driven only by `verify_go.sh:229` `go run
-race .` under the **default** toolchain, so a floor bump there cannot disarm anything"* — and
round 2's quorum reviewer objected to that trailing clause as an unmeasured premise. Measured
first-party as **V24** in `w-canary-control-does-not-survive-a-floor-raise.md`, and reproduced
verbatim at HEAD for this item (V1, V2):

- Inside the module the ambient toolchain is **`go1.26.4`** (the Homebrew base install), *not*
  the root-selected `go1.26.6` — because `racecontrol/go.mod` declares `go 1.22`, which is
  *below* the ambient base, so `GOTOOLCHAIN=auto` has nothing to push up to and Go falls back
  to the base binary.
- Base `go run -race .` → rc=1 with **2** `WARNING: DATA RACE` lines (the known-positive fires).
- Bump the directive `go 1.22` → `go 1.26.6`: with `GOTOOLCHAIN=auto` the control *still* fires
  — **only because `go1.26.6` is already in the local toolchain cache** and `auto` silently
  switched toolchains to satisfy the new floor.
- Add the one variable that removes that rescue, `GOTOOLCHAIN=local`: rc=1, output is exactly
  `go: go.mod requires go >= 1.26.6 (running go 1.26.4; GOTOOLCHAIN=local)`, **0** `WARNING:
  DATA RACE` lines — so `verify_go.sh:232` FATALs and the race-detector gate is dead. The
  control fails *before* execution; there is no race warning to count.

**The disarming is real, and its only concealment is a warm toolchain cache** — the exact
condition a fresh runner does not have. Row 42 filed this as **queue row 48** precisely because
a non-instance that an audit *publishes as measured* is "an unbound claim wearing a
measurement's clothes."

**Quorum round 1 BLOCKED at full strength** (all three external reviewers present) on this
draft's original direction — anchoring the invariant to the **oldest `KNOWN_GOOD`** toolchain.
The block and its resolutions are recorded in [Quorum History](#quorum-history). Two of the
three objections force the redesign in this revision:

1. **DIRECTIONAL (gpt5-6-sol):** `run.sh`'s `KNOWN_GOOD` list does **not** constrain the
   ambient toolchain used at `verify_go.sh:229` (a bare `go run -race .` under nested
   auto-selection). Bounding the floor `<= go1.24.9` proves nothing about whether the control
   is executable under the toolchain that actually runs it. The fix must **bind execution**
   with machinery, not assert against an unrelated list.
2. **PREMISE (oc-glm-5-2):** whether a `LOAD-BEARING` comment containing `go 1.22` breaks
   `moduleGoFloor`. Measured and recorded in the Verification Log (V8); a post-comment
   acceptance arm is added to AC7.

The third (non-blocking) reviewer, **gemini**, corrected a syntax error: `go env` takes
space-separated keys (`go env GOTOOLCHAIN GOVERSION`), not a comma-separated shorthand. Fixed
throughout this revision.

**Impact:**
- A blanket floor raise (`go mod tidy -go=…`, an IDE "update module", or a human matching the
  root `go 1.26.6`) silently destroys the *only* first-party positive control for the race
  detector, the instrument that stands behind every `-race` result in the go gate.
- The failure is **loud** (`verify_go.sh` FATALs) — which is why this is a row, not a blocker —
  but it is a **false** loud failure: it blames the race detector when the actual cause is a
  floor mismatch, and on a fresh runner it would fail every go gate run for the wrong reason.
- Row 42 explicitly did **not** claim systemic completeness; this item closes the gap it
  measured and deferred.

## Goals

**Primary Goal:** Bind the `racecontrol` module floor to the **root module's pinned toolchain**
— the toolchain the race control is now *made to run under* — so a floor raise that would
disarm the race-detector known-positive control reds the static gate before any
`go run -race .` runs, and so the runtime lane executes the control under a **known,
deny-list-vetted toolchain** rather than whatever nested auto-selection happens to pick.

**Success Metrics:**
- `verify_go.sh`'s race leg invokes the control with `GOTOOLCHAIN="$ACTIVE_GO"` (the root
  module's `go env GOVERSION`), binding execution to the same toolchain the deny-list just
  vetted — no nested auto-selection can differ.
- The `racecontrol` floor is bound by a static test in the gating lane (`host/verifygate`),
  asserting `racecontrol floor <= root module floor`, so CI reds on a disarming floor bump even
  though the runtime lane only garbles it.
- The static test also pins the execution-binding needle, so reverting the `verify_go.sh` edit
  (back to a bare `go run -race .`) is a named RED — the exact objection-1 hole cannot
  silently return.
- Every refusal/anti-vacuity branch of the new test has a named RED mutant; the disarming bump
  (floor above the root toolchain) is a named mutant.
- All three nested `go.mod` floors are bound after this item (root by existing tests, `repro`
  by row 42's Test A, `racecontrol` by this item) — the systemic audit closes out.
- The `racecontrol` floor directive is byte-unchanged (`go 1.22`); the only production edit to
  `racecontrol/go.mod` is a `LOAD-BEARING` comment.

## Quorum History

**Round 1 — BLOCKED at full strength** (all 3 external reviewers present). Decision: the round
is **not** a waiver, and the two blocking objections are accepted. The third is a non-blocking
syntax correction. This revision resolves all three; the redesign is grounded in first-party
measurements recorded in the [Verification Log](#verification-log).

| Reviewer | Objection | Verdict | Resolution (measured first-party) |
|----------|-----------|---------|-----------------------------------|
| gpt5-6-sol | **DIRECTIONAL** — anchoring the floor to the oldest `KNOWN_GOOD` (`go1.24.9`) is rejected. `run.sh`'s list does not constrain the ambient toolchain used at `verify_go.sh:229`; `floor <= go1.24.9` does not prove the control is executable under every possible host base. Re-design around machinery that actually binds execution. | **Blocking** | Re-anchored to the **root module floor** (`go1.26.6`) and added an **execution-binding production edit**: `verify_go.sh:229` now runs `GOTOOLCHAIN="$ACTIVE_GO" go run -race .`, where `ACTIVE_GO` is the root module's `go env GOVERSION` (V3). Under `GOTOOLCHAIN=auto`, `ACTIVE_GO >= root floor` always, so `racecontrol floor <= root floor` implies `racecontrol floor <= ACTIVE_GO` — the floor is satisfiable by the very toolchain that runs the control, under every possible host base (V2, V3). All `KNOWN_GOOD`-anchor claims are dropped; no claim that `KNOWN_GOOD` constrains this path remains. |
| oc-glm-5-2 | **PREMISE** — does a `LOAD-BEARING` comment containing `go 1.22` break `moduleGoFloor`'s exactly-one-`go `-line floor? | **Blocking** | Measured (V8): `moduleGoFloor` uses `strings.HasPrefix(line, "go ")`; a comment line begins with `//` at column 0, so it is ignored. A probe file containing the comment plus the real directive yields exactly one awk-equivalent column-0 `go ` line (`floors=[go1.22] count=1`), and the original file was restored byte-identically (sha256 `ab782f11…`). Recorded in the Verification Log (V8) and a **post-comment acceptance arm** is added to AC7. |
| gemini | **SYNTAX (non-blocking)** — write `go env GOTOOLCHAIN GOVERSION`, not a comma-separated shorthand. | Non-blocking, accepted | Corrected. All `go env` invocations in this revision use space-separated keys; confirmed valid (V3: `auto go1.26.6`). |

The round-1 reviewers' substantive direction (bind execution, don't assert against an unrelated
list) is now the design; the round-2 surface (row 44 CI inertness) is untouched by this
revision.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|----------|-----------------|-----------|----------|-------------|
| Anchor the invariant to the **root module floor** (`go1.26.6`, read from the root `go.mod` `go` line via `moduleGoFloor`), not the oldest `KNOWN_GOOD` and not a literal | `ACTIVE_GO` (the toolchain that runs the control after the execution-binding edit) is always `>= root floor` under `GOTOOLCHAIN=auto`; `KNOWN_GOOD`'s oldest `go1.24.9` never touches `verify_go.sh:229`'s ambient toolchain (round-1 R1). The root floor is the repo's own statement of its pinned toolchain — already the anchor of `TestGoToolchainPinsAgreeAndMatchJobList` and `PINNED` — so the bound derives from the checked value, not an authored list | agent (round-1-block-settled) | design | low |
| Bind as `<=` (floor at-or-below root toolchain), not `<` | Equality still arms the control (measured: floor `go1.26.6` under `GOTOOLCHAIN=go1.26.6` fires, V6); strict order over-constrains and adds no safety | agent (evidence-settled) | design | low |
| **Execution-binding production edit to `verify_go.sh`**: race leg runs `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` | This is the machinery that actually binds execution (round-1 R1). Without it, nested auto-selection inside `racecontrol/` can still pick a toolchain different from `ACTIVE_GO` (the `go1.26.4` base, V3), and the floor-to-root bound would be vacuous w.r.t. execution; with it, the control runs under the deny-list-vetted root toolchain and the static floor bound is meaningful | agent (round-1-block-settled) | design | low |
| **Not** option B (pin `GOTOOLCHAIN=local`) | `GOTOOLCHAIN=local` removes the cache rescue but does **not** prevent the disarm (floor > bound toolchain still fails loudly); it is a concealment fix, not a prevention | agent | design | low |
| New test lives in `host/verifygate/toolchain_pin_gate_test.go`, reusing `moduleGoFloor`/`shellAssignmentValues`/`repoRoot` | Same home and helpers as row 42's Test A; the file already owns run.sh/floor binding; `go/version` already imported | agent | design | low |
| **One-line production edit to `verify_go.sh`** is in scope (unlike the round-0 draft's "no production edit") | The round-1 DIRECTIONAL block makes the execution-binding edit the crux of the fix; the deny-list (lines 217–224) and the FATAL/arming logic (lines 226–236) are untouched | agent (round-1-block-settled) | design | low |

## Design Freeze

- **One test file modified, no test file created**: the new test appends to
  `host/verifygate/toolchain_pin_gate_test.go`, reusing `moduleGoFloor` and `repoRoot`;
  `go/version` is already imported (V9).
- **One-line production edit to `verify_go.sh:229`**: the race leg changes
  `go run -race .` → `GOTOOLCHAIN="$ACTIVE_GO" go run -race .`. `ACTIVE_GO` is already
  captured at `:217` (`ACTIVE_GO=$(go env GOVERSION)`), so no new variable, no new line, no
  change to the deny-list (`:218–224`) or the arming/FATAL logic (`:226–236`).
- **Comment-only production edit**: `racecontrol/go.mod` gains a `LOAD-BEARING` comment block
  above the directive; **the `go 1.22` line itself is not moved**.
- **NOT touched**: `ci.yml`, `run.sh`, `repro/`, the root `go.mod`, any `host/store` assertion.
- **Anchor is derived, not authored**: the invariant reads the root `go.mod` floor via
  `moduleGoFloor` and the `racecontrol` floor via `moduleGoFloor`, and the test only compares
  them (row 43's control-derived-from-checked-value refutation, binding here as in row 42).
- The invariant does not call `requirePinned` or read `AILANG_BIN`: a static text scan must
  run in any lane (row 42's P9 split).
- The `racecontrol/go.mod` comment must not collide with any needle the new test counts. The
  test's only comment-adjacent needle — the execution-binding substring
  `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` — is scanned in `verify_go.sh`, a different file,
  and the fence comment spells the channel as `GOTOOLCHAIN=$ACTIVE_GO` (no quotes, no `go run
  -race .` suffix), so no token the test scans collides.

## Deferred Decisions

- **Option B (pin `GOTOOLCHAIN=local`)** — agent may revisit later. Measured reasoning: it
  removes the *cache rescue* that conceals the disarm but leaves the *disarm* itself (floor >
  bound toolchain still reds `go run -race .`); it is a loudness improvement, not the prevention
  this row specifies. If a future item wants it, it is a one-line edit to `verify_go.sh`'s race
  leg and belongs in row 44's lane, not here.
- **An explicit repo-wide "oldest supported toolchain" policy** — would give future items a
  single maintained constant instead of deriving the bound from the root floor. Today the root
  floor is the repository's own statement of the pinned toolchain; if the repo ever adds an
  explicit oldest-toolchain policy, re-anchoring is a three-line change.

## Solution Design

### Overview

The disarming is a *static* condition — a floor directive above the toolchain that actually
runs the control — so the fix has **two** parts, one that makes the runtime lane bind execution
and one that statically enforces the bound in CI:

1. **Runtime binding (production, 1 line)**: `verify_go.sh:229` runs the control under
   `GOTOOLCHAIN="$ACTIVE_GO"`, where `ACTIVE_GO = go env GOVERSION` in the root module. This
   forces the race control to execute under the same toolchain the deny-list (`:218–224`) just
   vetted, removing the nested auto-selection that let an unrelated ambient base (`go1.26.4`)
   change which toolchain ran the control. Under `GOTOOLCHAIN=auto` the root module resolves
   `ACTIVE_GO >= root floor` always, so this binds execution to a toolchain at-or-above the
   root floor under every possible host base.
2. **Static gate (test)**: `TestRaceControlFloorStaysBelowRootToolchain` in `host/verifygate`
   asserts `racecontrol floor <= root module floor`, and pins the execution-binding needle so
   the round-1 hole cannot silently return. Because part 1 guarantees `ACTIVE_GO >= root floor`,
   this static bound is a genuine proof that the control is executable under the toolchain that
   actually runs it — answering the round-1 DIRECTIONAL objection with machinery.

The test is byte-for-byte the shape of row 42's `TestReproModuleFloorStaysBelowKnownBad
Toolchains` (same `moduleGoFloor`/`version.Compare` machinery), with a different anchor because
the two nested modules have different jobs and because round 1 rejected the `KNOWN_GOOD` anchor.

### Architecture

The repository maintains exactly **three** `go.mod` modules (census, V1): the **root** (`go
1.26.6`, the protection), **`repro`** (`go 1.22`, the miscompile reproducer, bound by row 42's
Test A to oldest `KNOWN_BAD`), and **`racecontrol`** (`go 1.22`, the race-detector
known-positive, currently *unbound* — this item). The race control is driven only by
`verify_go.sh:229` (V4, the only reference anywhere). `verify_go.sh` captures `ACTIVE_GO =
$(go env GOVERSION)` in the root module (`:217`), refuses to proceed under the deny-listed
toolchains `{go1.26.0 … go1.26.5}` (`:218–224`), then — after this item — runs the control as
`GOTOOLCHAIN="$ACTIVE_GO" go run -race .` from the `racecontrol/` directory (`:229`).

The execution chain is now a single, statically knowable toolchain: **`ACTIVE_GO`**, which is
`>= root floor` (`go1.26.6`) in all cases because the root `go.mod` declares `go 1.26.6` and
`GOTOOLCHAIN=auto` never resolves below the declaring module's floor. Therefore **`racecontrol
floor <= root floor` implies `racecontrol floor <= ACTIVE_GO`** — the floor is satisfiable by
the very toolchain that runs the control, for every possible host base. That chain is the
evidence-settled answer to round-1 R1: the invariant is anchored to the toolchain the control
is *made to run under*, not to an unrelated `KNOWN_GOOD` list.

**Components:**
1. **`verify_go.sh:229` one-line edit** — `race_control_output="$(cd
   design_docs/verification/w-race-gate-blindspot/racecontrol && GOTOOLCHAIN="$ACTIVE_GO" go
   run -race . 2>&1)"`. The control's execution is bound to the root toolchain.
2. **`TestRaceControlFloorStaysBelowRootToolchain`** (`host/verifygate`) — the static gate.
   Reads `racecontrol/go.mod`'s floor via `moduleGoFloor`; reads the root `go.mod`'s floor via
   `moduleGoFloor`; validates both with `version.IsValid`; asserts
   `version.Compare(racecontrolFloor, rootFloor) <= 0`, else `t.Fatalf` naming the exact
   consequence. Then reads `scripts/verify_go.sh` and asserts the execution-binding needle
   `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` occurs exactly once (the round-1 hole cannot
   return by reverting part 1).
3. **`racecontrol/go.mod` `LOAD-BEARING` comment** — the IDE-bump human's first sight, the
   same fence row 42 added to `repro/go.mod`, reworded for the root-toolchain anchor.

### Implementation Plan

**Phase 1: Static gate** (~0.5h)
- [ ] Append `TestRaceControlFloorStaysBelowRootToolchain` to
      `host/verifygate/toolchain_pin_gate_test.go` (reuses `moduleGoFloor` and `repoRoot`;
      `go/version` already imported).
- [ ] Run the AC1 run-existence form, the named mutations (M1–M8), and the pristine control.

**Phase 2: Runtime + fence** (~0.35h)
- [ ] Edit `scripts/verify_go.sh:229`: `go run -race .` → `GOTOOLCHAIN="$ACTIVE_GO" go run
      -race .`; confirm the deny-list and FATAL logic are untouched.
- [ ] Add the `LOAD-BEARING` comment block to `racecontrol/go.mod`; confirm the `go 1.22`
      directive is byte-unchanged (AC7).

**Phase 3: Hygiene + docs** (~0.1h)
- [ ] `go vet ./host/verifygate/`, `gofmt -l host/verifygate/`; confirm no test-name collision;
      note the row-48 linkage in `world-mission.md` if the iteration record needs it (doc-only).

## Files to Modify/Create

**Modified files:**
- `scripts/verify_go.sh` (+0 LOC net, 1 line altered at `:229`) — the race leg's invocation
  gains `GOTOOLCHAIN="$ACTIVE_GO"`; deny-list (`:218–224`) and arming/FATAL logic
  (`:226–236`) unchanged.
- `host/verifygate/toolchain_pin_gate_test.go` (+~50 LOC) — the new
  `TestRaceControlFloorStaysBelowRootToolchain`; no new imports, no name collisions (V9).
- `design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod` (+~8 LOC comment block)
  — `LOAD-BEARING` fence; the `go 1.22` directive byte-unchanged (sha256 `ab782f11…` base, V1).

**New files:** none.

No other files. `ci.yml`, `run.sh`, `repro/`, the root `go.mod`, any `host/store` assertion —
untouched.

## Examples

### Example 1: the disarming a floor bump causes, and what the static gate makes it red instead

**The attack (measured, V2 — the disarm under the execution-bound runtime):**

```
$ ACTIVE_GO=$(go env GOVERSION)          # root module: go1.26.6 (V3)
$ (cd design_docs/verification/w-race-gate-blindspot/racecontrol \
     && GOTOOLCHAIN="$ACTIVE_GO" go run -race .)
   # rc=1, 2× WARNING: DATA RACE            <- bound to ACTIVE_GO, the known-positive fires
$ sed -i '' 's/^go 1\.22$/go 1.27.0/' go.mod
$ GOTOOLCHAIN=go1.26.6 go run -race .
   # rc=1, output: go: go.mod requires go >= 1.27.0 (running go 1.26.6; GOTOOLCHAIN=go1.26.6)
   # 0× WARNING: DATA RACE                  <- verify_go.sh:232 FATALs, control is dead
```

**After this item**, the same edit reds the *static* gate in CI (the mutation M1), before any
`go run -race .`:

```
$ GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1
   # racecontrol module floor "go1.27.0" is above the root module floor "go1.26.6": ...
   #   --- FAIL
```

### Example 2: the new test, in the shape it is designed to guard

```
racecontrol/go.mod (comment-only edit):
// Deliberately a SEPARATE, NESTED module: `go build ./...` / `go test ./...` in the
// repository root module must NOT pick this race-detector control up. It is a
// diagnostic artifact, not host code, and it deliberately contains a data race.
//
// The `go 1.22` line below is LOAD-BEARING: it must stay at or below the repository
// root module floor (`go 1.26.6`, the `go` line of the root go.mod), because
// scripts/verify_go.sh runs this control under GOTOOLCHAIN=$ACTIVE_GO derived from
// the root module's GOVERSION; a floor above that toolchain makes `go run -race .`
// refuse before it can fire (verify_go.sh then FATALs "the race detector is not
// armed"). Enforced by TestRaceControlFloorStaysBelowRootToolchain (host/verifygate).
// Do not let `go mod tidy -go=…` or an IDE floor-bump touch it.
module ailang-world/verification/race-detector-control

go 1.22
```

## Success Criteria / Acceptance Criteria

Each AC carries its vacuity self-test and its **observed result on the unmodified tree at
`f1b5d1a`** (Verification Log rows cited).

- **AC1 — the test exists, RUNS, and passes on the post-sprint tree, in run-existence form.**
  `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run
  '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1 -v` → rc=0 with exactly 1 top-level
  `=== RUN` line and 1 `--- PASS`; a paired nonsense pattern (`-run
  TestNoSuchRaceControlFloorTestZZZ`) prints `[no tests to run]`, proving the instrument says
  so rather than passing vacuously. **Base @f1b5d1a: the verbatim named command → `no tests to
  run`, rc=0 (V9)** — the naive form greens at base measuring nothing; the `=== RUN` clauses
  are the repair and red at base (0 of 1).
- **AC2 — the pristine tree is green by construction and the eight RED mutations red it.**
  Base verdict on the prototype logic = `PASS: racecontrol floor go1.22 <= root floor
  go1.26.6; execution needle count=1` (V5); each of M1, M3–M8 reds on its named branch (V5);
  M2 is the equality GREEN control (V6); every arm restored sha256-byte-identical, porcelain 0
  (V1).
- **AC3 — M1 is the disarming edit and reds the static gate.** `racecontrol/go.mod` `go 1.22`
  → `go 1.27.0` (above the root toolchain `go1.26.6`) → `TestRaceControlFloorStaysBelowRoot
  Toolchain` FAILs on the bound clause (prototype verdict V5: `RACE-BOUND-FAIL`). The same edit
  *without* the gate is the disarm under the execution-bound runtime (V2: `0` DATA RACE,
  `verify_go.sh:232` would FATAL).
- **AC4 — M2 is the equality GREEN control with a runtime proof.** `racecontrol/go.mod` `go
  1.22` → `go 1.26.6` (exactly the root module floor): the gate stays GREEN (V5) **and** the
  control still fires at equality — `GOTOOLCHAIN=go1.26.6 go run -race .` from `racecontrol/`
  → rc=1, 2× `WARNING: DATA RACE` (V6). Had the runtime proof failed, `<=` would be the wrong
  bound and the strict form would return with that measurement as its justification; it did not
  fail.
- **AC5 — environment independence.** `env -u AILANG_BIN GOTOOLCHAIN=go1.26.6 go test
  ./host/verifygate/ -run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1` → rc=0
  with AC1's run-existence clauses; no network, no solver, no pinned binary (static text scan,
  V9's lane; the sibling static-lane run is env-independent at base, V9).
- **AC6 — hygiene.** `go vet ./host/verifygate/` rc=0 and `gofmt -l host/verifygate/` prints
  nothing. **Base: both green (V9)** — green-at-base by design; they measure the sprint's own
  edit.
- **AC7 — the fence comment landed, the directive did not move, and the post-comment floor
  still parses.** `grep -c "LOAD-BEARING"`
  `design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod` → 1 AND `grep -c '^go
  1.22$' …/racecontrol/go.mod` → 1 AND `shasum -a 256 …/racecontrol/go.mod` → the post-comment
  hash is the base hash plus the comment block (directive unchanged; base hash `ab782f11…`,
  V1). **Post-comment acceptance arm (round-1 R2):** `moduleGoFloor` on the commented file
  still returns exactly one valid floor (`floors=[go1.22] count=1`, V8) — the `//` comment
  line is ignored by `strings.HasPrefix(line, "go ")` and cannot break the exactly-one-floor
  refusal. **Base: `LOAD-BEARING` count = 0 (V8)** — the fence is AC7-red at base by design.
- **AC8 — systemic closure: all three `go.mod` floors are bound.** Census = exactly 3 modules
  (V1); root bound by `TestGoToolchainPinsAgreeAndMatchJobList`/`TestMiscompileInstrument
  ProbesPinnedToolchain`, `repro` bound by row 42's Test A, `racecontrol` bound by this item
  (V9: all existing pin tests green at base).

Explicitly rejected as an AC: "the full verify gate is green" alone — `go build ./... && go
test ./...` must of course pass on the final tree, but a package-wide `ok` can print while the
named test never ran (V9's `no tests to run` is rc=0); AC1's run-existence form is the binding
version.

## Conflict Surface

This change adds exactly **one** executable line to production code (`verify_go.sh:229`), which
*strengthens* the runtime lane rather than altering the deny-list or arming logic; and it adds
no executable line to the parsers/typecheckers/codegen. The conflict surface is entirely in the
shared files the static text gate reads and the shared helpers it reuses.

- **`verify_go.sh`** — the race leg is now *edited* (not merely observed) by this item, at
  exactly `:229`. `ACTIVE_GO` (`:217`) and the deny-list (`:218–224`) are read, not changed; the
  arming/FATAL block (`:226–236`) is unchanged. The execution-binding needle
  `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` is the only `verify_go.sh` line the new test scans;
  row 44's surface (CI `continue-on-error`) is untouched.
- **`racecontrol/go.mod`** — read by this test (`moduleGoFloor`) and executed by
  `verify_go.sh:229`; written by nobody. A *moved* file (M8) reds this test's read floor and
  breaks `verify_go.sh` at runtime; the test names the path, so a move must move the fence too.
- **Root `go.mod`** — read by this test (`moduleGoFloor`, as the anchor) and by
  `TestGoToolchainPinsAgreeAndMatchJobList` (floor == ci.yml pin) and
  `TestMiscompileInstrumentProbesPinnedToolchain` (PINNED). Shared fate on a *deleted* or
  *malformed* `go` line (M5, M6); the sibling reds on its own clauses. A root-floor *raise* above
  `racecontrol`'s floor keeps this test green (the bound is `<=`), which is correct: the control
  is still runnable under the higher root toolchain.
- **`run.sh`'s `KNOWN_GOOD=` line** — **no longer read by this test** (round-1 R1 dropped the
  `KNOWN_GOOD` anchor). It remains read by `TestMiscompileInstrumentProbesPinnedToolchain`
  (reads `KNOWN_GOOD`, `KNOWN_BAD`, `PINNED`). This test's execution bound no longer shares fate
  with `KNOWN_GOOD` edits.
- **Shared helpers** (`moduleGoFloor`, `repoRoot`) — a helper edit reds every consumer at once;
  already this file's accepted pattern (row 42).
- **`repro/` and the root module's other machinery** — not read by this test beyond the root
  `go.mod` `go` line; existing bindings are unaffected (V9 confirms all existing pin tests green
  at base).

### Programs that MUST still work

1. `host/verifygate/toolchain_pin_gate_test.go` → `TestReproModuleFloorStaysBelowKnownBad
   Toolchains` (row 42's Test A) — still green at base (V9).
2. … → `TestMiscompileInstrumentProbesPinnedToolchain` — still green at base (V9); unaffected
   by the `KNOWN_GOOD` anchor being dropped (that test owns `KNOWN_GOOD` on its own).
3. `scripts/verify_go.sh` race leg (lines 226–235) — after the `:229` edit, `go run -race .`
   still fires 2× `WARNING: DATA RACE` from `racecontrol/` under `GOTOOLCHAIN="$ACTIVE_GO"`
   (V2), and the FATAL/arming block is unchanged (V4).
4. `design_docs/verification/w-race-gate-blindspot/racecontrol/main.go` — `go run -race .`
   still fires 2× `WARNING: DATA RACE` under `go1.26.4`, `go1.26.6`, and the equality floor
   (V2, V6).
5. `design_docs/verification/w-race-gate-blindspot/repro/go.mod` — unchanged; row 42's binding
   intact (V9).

### What deliberately changes

At runtime, exactly one invocation: the race control is now run under an explicit
`GOTOOLCHAIN="$ACTIVE_GO"` instead of nested auto-selection. The deliberate effect is that a
`racecontrol` floor bump above the root toolchain is now a **red static gate** instead of a
silent (cache-concealed) or wrong-reason loud (fresh-runner) disarm, and that the runtime lane
executes the control under a known, deny-list-vetted toolchain. `racecontrol/go.mod`'s directive
is byte-unchanged.

## Testing Strategy

**Unit / gate tests:**
- `TestRaceControlFloorStaysBelowRootToolchain` — the static gate, with run-existence form
  (AC1) and environment independence (AC5).

**Mutation testing (non-vacuity):**
- M1–M8 as enumerated below; one named RED mutant for every refusal and anti-vacuity floor;
  M2 as the equality GREEN control with a runtime arming proof (AC4).

**Regression-surface:**
- All five "Programs that MUST still work" entries re-run green at base (V9) and re-confirmed
  by the sprint.

**Manual / runtime lane:**
- `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` from `racecontrol/` → control fires (V2); at the
  equality floor `GOTOOLCHAIN=go1.26.6 go run -race .` → control fires (V6).

## Non-Goals

- **Not** pinning `GOTOOLCHAIN=local` in `verify_go.sh` — option B; it removes concealment but
  not the disarm, and belongs in row 44's lane. [Deferred, see Deferred Decisions]
- **Not** altering `verify_go.sh`'s deny-list (`:218–224`) or its arming/FATAL logic
  (`:226–236`). Only the single invocation at `:229` gains the `GOTOOLCHAIN` binding. [Out of
  scope]
- **Not** binding the floor to `run.sh`'s `KNOWN_GOOD` list in any form — round-1 R1 rejected
  that anchor; the bound is the root module floor, derived from the toolchain that actually
  runs the control. [Quorum-round-1-settled]
- **Not** rewriting the race control's runtime behaviour; `main.go` unchanged. [Out of scope]
- **Not** a runtime pre-flight that compares the floor to the live toolchain. [A static gate in
  CI is the gating lane; the runtime lane's CI inertness is row 44's declared problem]
- **Not** re-binding the `repro` module (already bound by row 42's Test A). [Already shipped]

## Timeline / Milestones

**Milestones:**
- **M1 — gate + runtime binding land**: `verify_go.sh:229` gains `GOTOOLCHAIN="$ACTIVE_GO"`;
  `TestRaceControlFloorStaysBelowRootToolchain` added, AC1/AC5 green, M1–M8 red on named
  branches (AC2/AC3/AC4), pristine control green.
- **M2 — fence lands**: `racecontrol/go.mod` `LOAD-BEARING` comment added, directive
  byte-unchanged (AC7), post-comment floor parses (V8).
- **M3 — hygiene + audit closure**: vet/gofmt clean (AC6), no collisions (V9), systemic
  closure confirmed (AC8).

**Effort:** Phase 1 ~0.5h, Phase 2 ~0.35h, Phase 3 ~0.1h → ~1h total (matches the ~0.1d row
estimate).

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| The anchor (root module floor = `go1.26.6`) is a proxy for "the toolchain that runs the control"; if `ACTIVE_GO` ever resolved *below* the root floor the static bound would overstate the guarantee | Low | `ACTIVE_GO = go env GOVERSION` in the root module under `GOTOOLCHAIN=auto` never resolves below the declaring module's floor (V3); the execution-binding edit makes `ACTIVE_GO` the control's actual runtime toolchain, so the bound is a proof, not a proxy |
| A base toolchain *above* the root floor (future `go1.30`) raises `ACTIVE_GO` above `go1.26.6` | None | The static bound is `racecontrol floor <= root floor`; since `ACTIVE_GO >= root floor`, the control still runs under any higher `ACTIVE_GO` — the bound is monotone-safe in the host-base direction |
| The static test reads `verify_go.sh` text, so an execution binding spelled differently (computed, indirection) escapes it | Low | The exact needle `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` is pinned at count=1 (M7); any re-spelling reds the count — the honest bound, matching row 42's declared residual |
| M5/M6's root-`go.mod` mutants could be mistaken for live root-module changes | Low | They are declared instrument-failure floors on the anchor; the root module's own tests red on the same edits, so they measure this test's validity/comparison arms, not a real drift |
| The `LOAD-BEARING` comment could be misread as colliding with a scanned needle | Low | The only comment-adjacent needle lives in `verify_go.sh` (a different file) and the fence spells the channel `GOTOOLCHAIN=$ACTIVE_GO` (no quotes, no `go run -race .`), so no collision (V8) |

### Declared residual

The static gate cannot know the *live* `ACTIVE_GO` on a given runner; it binds the floor to the
root module floor (`go1.26.6`), which is `<= ACTIVE_GO` in all cases because the root module
resolves `ACTIVE_GO` at-or-above its own floor under `GOTOOLCHAIN=auto`. The only true residual
is the pre-existing boundary where the root module itself cannot resolve `ACTIVE_GO` at all — a
host whose base Go is below `go1.26.6` with no network to fetch it, which fails the entire root
module before this control runs, independent of this item. The runtime `verify_go.sh` FATAL
stays the loud backstop for any residual case.

## Axiom Compliance

**Scoring**

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +1 | The gate is a deterministic static text scan; the same tree gives the same verdict in any lane |
| A2: Replayability | 0 | No replay/session surface |
| A3: Effect Legibility | 0 | No effect/IO change to AILANG programs; one invocation-line strengthening in the shell gate |
| A4: Explicit Authority | 0 | No capability change |
| A5: Bounded Verification | +1 | Replaces a runtime-only, cache/auto-selection-dependent signal with a bounded static gate in CI that binds execution to a known toolchain |
| A6: Safe Concurrency | 0 | No concurrency change (the race control's own race is untouched and still fires) |
| A7: Machines First | +1 | A static gate is cheap, deterministic, and runnable by CI without a warm toolchain cache; the runtime lane now pins the toolchain it executes |
| A8: Minimal Syntax | 0 | No syntax change |
| A9: Cost Visibility | 0 | No cost surface |
| A10: Composability | +1 | Reuses row 42's helpers and shape; composes with the existing three-module binding |
| A11: Structured Failure | +1 | The disarm is now a named, attributed RED (a floor mismatch against the bound toolchain) instead of a wrong-reason FATAL |
| A12: System Boundary | 0 | No boundary change |

**Net Score: +5** ✅ Proceed to implementation

**Hard Violation Check**
- [x] A1 (Determinism): no implicit nondeterminism — the gate reads files and compares versions
- [x] A3 (Effects): no hidden side effects — comment-only production edit plus a one-line
      invocation binding in the gate's own shell script
- [x] A4 (Authority): no ambient access granted
- [x] A7 (Machines First): the static gate serves CI and fresh runners; the runtime lane
      executes under an explicit toolchain

## Non-Vacuity — named RED mutation for every refusal and anti-vacuity floor

Production side mutated (`racecontrol/go.mod` floor, the root `go.mod` floor, `verify_go.sh`'s
execution binding, file placement) — never the test helpers. Every assertion maps to a mutant.
The instrument-failure floors row 42's shape carries — *exactly one `go ` line* (M4) and *every
floor token valid* (M3, M6) — are enumerated here as branches with their own mutants per the
iter-108 floor-pinning rule: the anti-vacuity floors are a refusal about the *measurement*, not
about the input, and they get one arm each. M7 is the round-1 R1 regression arm: reverting the
execution binding re-opens the exact hole the quorum blocked, so it is a named RED.

| # | Exact edit | Expected RED (single test name) | Shape / branch |
|---|---|---|---|
| M1 | `racecontrol/go.mod` `go 1.22` → `go 1.27.0` (above the root floor `go1.26.6`) | bound clause: floor `go1.27.0` above root floor `go1.26.6` | **threat-shaped: the disarming edit under the execution-bound runtime (V2: 0 DATA RACE)** — the runtime lane garbles it, the static gate must name it |
| M2 | `racecontrol/go.mod` `go 1.22` → `go 1.26.6` (exactly equal to the root floor) | **GREEN** (equality control) | boundary-pair, not an arm; runtime arming proof at equality (V6) |
| M3 | `racecontrol/go.mod` `go 1.22` → `go banana` | floor-validity floor: `racecontrol floor "gobanana" invalid` | malformed-floor-shaped: a floor `version.Compare` would misorder is itself a disarm |
| M4 | delete the `go 1.22` line from `racecontrol/go.mod` | exactly-one-go-line floor: `moduleGoFloor` fatal `found 0 go lines, want 1` | anti-vacuity floor: a control with no floor has no bound |
| M5 | root `go.mod` `go 1.26.6` → `go 1.20` (below the `racecontrol` floor `go1.22`) | bound clause: `racecontrol floor go1.22` above root floor `go1.20` | anchor-direction anti-vacuity: proves the root-anchor comparison is enforced, not a tautology |
| M6 | root `go.mod` `go 1.26.6` → `go banana` | root-floor-validity floor: `root floor "gobanana" invalid` | anchor-validity anti-vacuity: the anchor is validated, not silently trusted |
| M7 | `verify_go.sh:229` `GOTOOLCHAIN="$ACTIVE_GO" go run -race .` → `go run -race .` (revert the binding) | execution-binding needle count=0, want 1 | **round-1-R1 regression: re-opens the nested-auto-selection hole the quorum blocked** — the machinery that binds execution must not silently revert |
| M8 | `git mv racecontrol/go.mod racecontrol/go.mod.moved` (or move the dir) | read floor: `moduleGoFloor` `os.ReadFile` fatal — the test names its target by path | placement-shaped: a moved control must move the gate in the same edit |

**Green control for all arms:** the unmutated post-sprint tree passes AC1/AC4/AC5, and every arm
ends restored sha256-byte-identical with `git status --porcelain` empty (V1/V2 recipe). The
sibling `TestMiscompileInstrumentProbesPinnedToolchain` and row 42's `TestReproModuleFloor
StaysBelowKnownBadToolchains` stay green at base (V9) and are re-confirmed by the sprint on the
post-sprint tree (AC3's re-confirm clause).

## Milestones (repeat for clarity)

- **M1 — gate + runtime binding land** (AC1, AC2, AC3, AC4, AC5): `verify_go.sh:229` edit + new
  test in `host/verifygate`.
- **M2 — fence lands** (AC7): `LOAD-BEARING` comment, directive unchanged, post-comment floor
  parses.
- **M3 — hygiene + closure** (AC6, AC8): vet/gofmt clean, three-module systemic closure.

## Verification Log

All rows run first-party by the designer at `f1b5d1a` (clean tree, porcelain 0 before and after
every arm), shell `zsh`, `PATH=/opt/homebrew/bin:$PATH`, darwin/arm64, 2026-08-28. KP =
known-positive control carried in the same call. The mutation prototypes ran a standalone Go
replica of the proposed test's logic against the real `racecontrol/go.mod`, root `go.mod`, and
`verify_go.sh` (base) and against temp mutated copies (M1–M8), deleted after the run — the
sprint re-runs every arm against the landed test. Rows V1–V4, V6, V8, V9 are the round-1 quorum
re-measurements; every `go env` invocation uses the space-separated `go env GOTOOLCHAIN
GOVERSION` form (gemini's syntax correction).

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Worktree is `f1b5d1a`, clean; module census is exactly 3; the `racecontrol` floor is `go 1.22` with the V24 hash; worktree `go.mod` untouched by the round | `git rev-parse HEAD`; `git status --porcelain \| wc -l`; `find . -name go.mod -not -path './.git/*'`; `shasum -a 256 design_docs/verification/w-race-gate-blindspot/racecontrol/go.mod` | `f1b5d1a79d89252…`; `0` (planned doc + log are the only non-clean entries, untouched); exactly `./go.mod`, `…/repro/go.mod`, `…/racecontrol/go.mod`; sha256 `ab782f11db0f7f259f73dd55a58eaf5a30b871bb79bd98bacbe964d50efc025b` (identical to row 42's V24 restored hash) |
| V2 | **THE FINDING re-measured under the execution-bound runtime**: the disarming edit is a floor above the bound toolchain, and the bound toolchain is `ACTIVE_GO` | `ACTIVE_GO=$(go env GOVERSION)`; `race_control_output="$(cd racecontrol && GOTOOLCHAIN="$ACTIVE_GO" go run -race . 2>&1)"` (base floor `go 1.22`); then `sed go 1.22→go 1.27.0` and re-run the same command; `cp` restore + sha256 + porcelain | base: rc=1, **2** `WARNING: DATA RACE` (the known-positive fires under the binding); floor `go 1.27.0`: rc=1, output `go: go.mod requires go >= 1.27.0 (running go 1.26.6; GOTOOLCHAIN=go1.26.6)`, **0** DATA RACE → `verify_go.sh:232` would FATAL. Restore byte-identical `ab782f11…`, porcelain 0; post-restore control re-fires (rc=1, 2 DATA RACE). (Also confirmed: floor `go 1.26.6` == `ACTIVE_GO` still fires **2** DATA RACE — equality is armed, so the old V24 edit stops being a disarm once execution is bound) |
| V3 | **The execution-binding chain**: `ACTIVE_GO` (root `go env GOVERSION`) = root floor, and the nested module's ambient toolchain differs; `go env` takes space-separated keys (gemini's fix) | `go env GOTOOLCHAIN GOVERSION` at root; `go version` at root; `go env GOTOOLCHAIN GOVERSION` inside `racecontrol/` | root: `auto go1.26.6` (two-arg space form rc=0), root `go version go1.26.6 darwin/arm64`; inside `racecontrol/`: `go1.26.4 auto` — the nested auto-selection would pick `go1.26.4`, which the execution-binding edit removes; `ACTIVE_GO = go1.26.6 = root go.mod floor`, so `racecontrol floor <= root floor` ⇒ `floor <= ACTIVE_GO` |
| V4 | `racecontrol` is driven only by `verify_go.sh:229`; the race leg is lines 226–235 and the deny-list/capture are `:217–224` | `git grep -n "racecontrol" -- '*.go' '*.sh' '*.yml'` (excluding design docs); `grep -n "ACTIVE_GO=\|go run -race\|case \\"\\$ACTIVE_GO\\"" scripts/verify_go.sh` | exactly one hit: `scripts/verify_go.sh:229` `cd .../racecontrol && go run -race . 2>&1`; `:217` `ACTIVE_GO=$(go env GOVERSION)`, `:218` `case "$ACTIVE_GO"`, `:226` race-leg header, `:232` `grep -q 'WARNING: DATA RACE'` else FATAL |
| V5 | **Prototype: base GREEN and M1–M8 verdicts** (the test's logic, run against real + mutated files) | standalone replica of `moduleGoFloor`+`version.Compare` over `racecontrol/go.mod`, root `go.mod`, `verify_go.sh` (base) and temp mutated copies (M1–M8) | base `PASS: racecontrol floor go1.22 <= root floor go1.26.6; execution needle count=1`; M1 `RACE-BOUND-FAIL` (`go1.27.0` > `go1.26.6`); M2 `PASS` (`go1.26.6` == `go1.26.6`); M3 `racecontrol floor "gobanana" invalid`; M4 `found 0 go lines, want 1`; M5 `RACE-BOUND-FAIL` (`go1.22` > `go1.20`); M6 `root floor "gobanana" invalid`; M7 `execution needle count=0, want 1`; M8 `open …/go.mod: no such file or directory` |
| V6 | **M2 runtime proof: the control is armed at the equality floor** | `GOTOOLCHAIN=go1.26.6 go run -race .` from `racecontrol/` with floor `go 1.26.6`; `GOTOOLCHAIN=go1.26.6 go version` | rc=1, **2** `WARNING: DATA RACE` — the equality floor still builds and fires the control, so `<=` is the right bound |
| V7 | Version comparison semantics for every token the gate compares | standalone `go/version`: `Compare`/`IsValid` on the bound and mutant tokens | `Compare(go1.22,go1.26.6)=-1`; `Compare(go1.27.0,go1.26.6)=1`; `Compare(go1.26.6,go1.26.6)=0`; `Compare(go1.22,go1.20)=1`; `IsValid(gobanana)=false` — `go/version` is stdlib, no new dependency |
| V8 | **Round-1 R2 premise: a `LOAD-BEARING` comment containing `go 1.22` does not break `moduleGoFloor`**; and the base has no fence yet | `moduleGoFloor`-replica over the real `racecontrol/go.mod` and over a temp file containing the proposed comment block plus the real directive; `grep -c "LOAD-BEARING" racecontrol/go.mod` | real file: `floors=[go1.22] count=1`; comment+directive probe: `floors=[go1.22] count=1` (exactly one awk-equivalent column-0 `go ` line; the `//` comment is ignored by `strings.HasPrefix(line, "go ")`); base `LOAD-BEARING` count = **0** (AC7 base-RED); probe file discarded, original restored byte-identically |
| V9 | Base gate lane: all existing pin tests green; the new name is unclaimed and its scoped run reports `no tests to run` (the AC1 vacuity trap); no-AILANG_BIN static lane works; hygiene/collision baselines | `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run 'TestReproModuleFloorStaysBelowKnownBadToolchains\|TestCanaryDeclaresPositiveArmOnly\|TestMiscompileInstrumentProbesPinnedToolchain\|TestMiscompileInstrumentStepIsGatedInCI\|TestGoToolchainPinsAgreeAndMatchJobList' -count=1`; `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestRaceControlFloorStaysBelowRootToolchain$' -count=1`; `env -u AILANG_BIN GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestMiscompileInstrumentProbesPinnedToolchain$' -count=1`; `grep -rn "TestRaceControlFloorStaysBelowRootToolchain" host/` | existing: `ok github.com/sunholo-data/ailang-world/host/verifygate 0.276s`; new name: `ok ... [no tests to run]`, rc=0 (green-at-base measuring nothing); no-AILANG_BIN scoped rc=0 (static lane is env-independent); no existing test with that name |

## Related Documents

- [`../implemented/w-canary-control-does-not-survive-a-floor-raise.md`](../implemented/w-canary-control-does-not-survive-a-floor-raise.md) —
  row 42: the item that measured **V24** (this item's source evidence), built the exact test
  shape (`TestReproModuleFloorStaysBelowKnownBadToolchains`, `moduleGoFloor`,
  `shellAssignmentValues`, `repoRoot`) and helpers this item reuses, and **explicitly deferred
  `racecontrol/` to queue row 48** because its correct binding ("buildable by whatever
  toolchain verify_go.sh happens to run") was not statically knowable there. This item decides
  that binding by *making* the toolchain statically knowable (`GOTOOLCHAIN="$ACTIVE_GO"`) and
  closes the deferral.
- [`../implemented/w-race-gate-blindspot.md`](../implemented/w-race-gate-blindspot.md) — where
  the canary and the nested-module pattern (including `racecontrol/` as the race-detector
  known-positive) were born.
- [`../implemented/w-setup-go-pin-unguarded.md`](../implemented/w-setup-go-pin-unguarded.md) —
  row 41: built the file this item extends and the root-floor binding this item's anchor (the
  root `go.mod` floor) sits alongside.
- [`../implemented/w-miscompile-instrument-inert-in-ci.md`](../implemented/w-miscompile-instrument-inert-in-ci.md) —
  row 44: owns `run.sh`'s CI inertness and `ci.yml:172`; the declared reason a *static* gate is
  this item's gating lane rather than the runtime lane.
- `design_docs/world-mission.md` — queue row 48 (this item); rows 42 (the sibling binding),
  43 (floor-raise coupling inventory), 44, 45, 46.
- `design_docs/verification/w-race-gate-blindspot/` — `racecontrol/` (this item's target),
  `repro/` (row 42's target), `run.sh` (whose `KNOWN_GOOD` anchor round 1 rejected; the
  `PINNED`/root-floor chain remains row 41/42's concern).

## Future Work

- **Option B (`GOTOOLCHAIN=local`)** — a loudness improvement that removes the cache rescue; if
  adopted, it belongs in `verify_go.sh`'s race leg (row 44's lane) as a one-line change
  complementing this static gate and execution binding.
- **An explicit repo-wide "oldest supported toolchain" policy** — would give future items a
  single maintained anchor to bind all nested modules to, instead of deriving the bound from the
  root module floor.

---
**DESIGN_DOC_PATH**: `design_docs/planned/w-racecontrol-floor-bump-disarms-the-race-control.md`

# w-shell-assignment-parser-drops-an-indented-assignment

> **SUPERSEDED.** The attended ruling `D-WORLD-31` (2026-09-03) folded this work into
> `w-load-bearing-criteria-need-a-mutation-not-a-grep.md`; do not fund row 50 separately. Row 50's
> defect—a deny-list silently narrowed by an indented assignment—is provably gone because
> `KNOWN_BAD`, `KNOWN_GOOD`, and `PINNED` now live in a data-only fixture whose every non-blank,
> non-comment line must match the column-0-anchored record grammar. An indented assignment cannot
> exist there without `TestToolchainPinFixtureIsDataOnly` redding (AC2, discharged by mutations
> M1/M2/M8).

**Status**: **SUPERSEDED under attended ruling `D-WORLD-31` (2026-09-03)** by
`w-load-bearing-criteria-need-a-mutation-not-a-grep.md`.

> **⚠ READ THIS BEFORE IMPLEMENTING ANY OF THE DESIGN BELOW.** Two quorum rounds both BLOCKED. The
> surviving objection (`gpt5-6-sol`, round 2) disputed the design DIRECTION, and the controller
> MEASURED its premise TRUE: **leading whitespace is not syntax in bash.** A script setting `COL0=`,
> a space-indented `INDENTED=` and a tab-indented `TABBED=` at top level prints all three values,
> rc=0. So the earlier design's central rationale — that an indented assignment is a
> "conditional/nested shadow" and must be counted as a DEVIATION rather than extracted — is
> **FALSE**, and the two-sided invariant built on it rejects valid shell.
>
> The item was parked for the human. The human has now ruled (D-WORLD-29, attended 2026-09-01,
> VERBATIM — normative, not to be argued with): **"ANSWERED — A — ACCEPT the indented assignment.
> One whitespace-tolerant scan: trim leading spaces/tabs, ignore comments, extract NAME=\"...\"
> regardless of indentation, require len(values) == 1. B's stated rationale ('an indented
> assignment is a conditional/nested shadow') is measured FALSE, and defence-in-depth bought with a
> false premise is not worth a loud red on a benign re-indent."** This document is revised to that
> ruling: option A, with every trace of the rejected option B removed — not annotated, not kept as
> an alternative.
>
> Two further findings survive EITHER rule and must not be lost: (1) `oc-glm-5-2` (restored from
> `absent`, verdict PASS) observes that declared residual 1 — an `export KNOWN_BAD="…"` inside the
> same `if` — leaves the IDENTICAL silent-narrowing hole open; (2) `gemini-3-1-pro` is right that
> the pre-existing failure message hardcodes `KNOWN_BAD` for all four call sites, so a `KNOWN_GOOD`
> or `PINNED` deviation would emit a false variable name; template it. Both are carried as named
> residuals below (residual 1 and a note under Call sites), not absorbed into this change.
>
> Everything else in this doc stands and was verified: the reuse audit, the declared-residual
> list, the anti-vacuity floors, and every command — each of which the controller re-executed
> verbatim after the round-1 repair. The controller's first-party measurements at HEAD this
> iteration (ARM 0 / ARM A / ARM B at `cb73cab`) are recorded in the Verification Log and re-derived
> under ruling A in the mutation table.
**Target**: current iteration (World iter-140 / queue row 50)
**Priority**: P2 (latent gap, not a live defect)
**Estimated**: ~0.1d (~1h)
**Dependencies**: None (charter row 50 is gated on nothing)
**Planner-Lane**: codex-ok
**Filed by**: the `sonnet` evaluator of `P42`, 2026-08-27 (iter-131), non-blocking finding 2

## Problem Statement

`shellAssignmentValues` (`host/verifygate/toolchain_pin_gate_test.go:215`) is the single scanner
that reads `KNOWN_GOOD` / `KNOWN_BAD` / `PINNED` out of the miscompile instrument
`design_docs/verification/w-race-gate-blindspot/run.sh` for both gate tests. It matches with
`strings.HasPrefix(line, name+"=\"")` — **column-0 anchored** — so an assignment preceded by any
whitespace is not merely mis-parsed, it is *not counted at all* (V-A, V-H2).

**This is a latent gap, not a live bug.** Today `run.sh` is flat and sets each name exactly once,
at column 0 (V-B), and the pristine consumers pass (V-C). Nothing is currently wrong. The value of
this row is closing a **silent failure mode**, not fixing broken behavior.

What makes the gap worth a row is the **asymmetry**, measured first-party by the controller in
this worktree:

- Every *other* malformed shape fails **LOUDLY** through the instrument-failure floors. A
  single-quoted `KNOWN_BAD='…'` contributes 0 values and trips the `count=0, want 1` floor. A
  commented-out `# KNOWN_BAD=` is correctly ignored.
- The one **SILENT** case is a **second, indented** assignment beside the valid column-0 one —
  exactly the shape a refactor produces when someone wraps a narrowing override in a conditional:

  ```bash
  KNOWN_BAD="go1.26.0 go1.26.3 go1.26.4 go1.26.5"   # the counted, column-0 assignment
  if [ -n "${NARROW_DENYLIST:-}" ]; then
    KNOWN_BAD="go1.26.0"                             # INVISIBLE to shellAssignmentValues
  fi
  ```

  Measured (V-D, ARM A): with this mutation landed (sha-asserted, valid-shell-asserted, effect
  asserted via gates G1 = 1 beside G2 = 1 — gate commands printed under the mutation table), **both
  consumers stay GREEN** — rc=0, RUN=2 PASS=2. The count stays 1, the pin checks all read the
  top-level value, and `TestReproModuleFloorStaysBelowKnownBadToolchains` binds the repro-module
  floor against the WIDE deny-list `go1.26.0 go1.26.3 go1.26.4 go1.26.5` while `run.sh` would at
  runtime narrow it to `go1.26.0`. The gate would then certify a floor bound against a deny-list
  the instrument no longer probes.

**THE ITEM (charter row 50, verbatim intent):** count assignments with leading whitespace
tolerated and require the total to be 1 — so a refactor that adds a conditional branch reds
loudly instead of narrowing the deny-list in silence — and add one arm per shape to the drill.

## Goals

**Primary Goal:** make a second, whitespace-indented `NAME="…"` assignment in `run.sh` a **named
RED** in both consumers of `shellAssignmentValues`, by extracting values regardless of indentation
and requiring `len(values) == 1` — without weakening any shape that fails loudly today
(single-quoted, empty) and without changing which line supplies the counted value.

**Success Metrics:**
- The V-D / ARM A mutant (second indented `KNOWN_BAD=` beside a valid column-0 one) reds BOTH
  consumers with the pre-existing `assignment count=2, want 1` message, where today it passes 2/2.
- The V-E / ARM B arm (indented-only) is **GREEN** (rc=0, RUN=2, PASS=2) — a deliberate, ratified
  behaviour change: today's red fires because the scanner cannot SEE a real assignment, not because
  anything is wrong.
- All three names (KNOWN_GOOD, KNOWN_BAD, PINNED) and both consumers inherit the fix through the
  one shared helper; no call site re-implements the scan (V-H1, V-I), and no new assertion is added
  at any call site.
- Commented assignments — column-0 `# KNOWN_BAD="…"` AND indented `   # KNOWN_BAD="…"` — remain
  ignored, proven by unit-test arms, not assumed.
- What the parser still cannot see is **declared** — in a code comment, pinned by unit-test arms,
  and named in an acceptance criterion — never assumed.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|----------|-----------------|-----------|----------|-------------|
| The invariant is **one whitespace-tolerant scan**: trim leading spaces/tabs, ignore comments, extract `NAME="…"` regardless of indentation, require `len(values) == 1` | The naive reading of the row text is the CORRECT one. Under it, the addition mutant (a second indented assignment beside the column-0 one) yields count=2 and REDS via the pre-existing `assignment count=%d, want 1` check — the row's whole purpose is served. An indented-only assignment yields count=1 and supplies its value, which is what bash actually does. Consequence: today's loud red on the indented-only shape disappears — and it is measured to fire because the scanner cannot SEE a real assignment, not because anything is wrong. The rejected alternative (a two-sided invariant counting indented occurrences as deviations) was built on a premise the controller measured FALSE | **human** (D-WORLD-29, attended 2026-09-01) | design | low |
| Fix lives in the shared helper as **one scan** (`func shellAssignmentValues(lines []string, name string) []string`), signature unchanged, no second fact returned | One scanner over the prefix grammar means the tolerance and the anchor cannot disagree (V-H1 shows the helper is already the only value scan; keep it that way). No call site changes, no new assertion | agent | design | low |
| `export`/`local`/unquoted/here-doc forms are **declared residuals**, not covered | Covering them opens token-normalization scope (arbitrary declaration prefixes) on a ~0.1d row; each is loud in the replacement direction (column-0 count drops to 0) and pinned as a residual by a unit arm so a future cover is a deliberate change, not drift | agent | design | low |

## Design Freeze

The design DIRECTION is **human-ratified** under `D-WORLD-29` (attended 2026-09-01): option A is
accepted verbatim. No further human ratification is required; the remaining decisions are
agent-resolvable inside the sprint. (No quorum triggers fire: no shared-machinery override — the
helper's only consumers are the two tests in the same file — no cost/KPI/schema surface, all
premises verified in-repo.)

## Deferred Decisions

- Whether `run.sh`'s three-name census itself (the hand-maintained
  `[]string{"KNOWN_BAD=", "KNOWN_GOOD=", "PINNED="}` control list at `:245`) should be derived or
  extended is out of scope; this row changes how each name is scanned, not which names exist.

## Solution Design

### Overview

Modify `shellAssignmentValues` (the single definition, `toolchain_pin_gate_test.go:215-227`,
V-A/V-I) to scan each line with leading whitespace trimmed and extract the value regardless of
indentation. The signature is **unchanged** — `func shellAssignmentValues(lines []string, name
string) []string` — and there is no second return value:

```go
// shellAssignmentValues returns the values of NAME="…" assignments, in order,
// regardless of leading whitespace (spaces/tabs). A commented assignment begins
// with `#` after trimming and can never begin NAME + "=\"", so it is ignored.
// DECLARED RESIDUAL: this scanner cannot see `export NAME="…"`, `local NAME="…"`,
// an unquoted/single-quoted NAME=…, a mid-line assignment after another token
// (`if …; then NAME="x"; fi` on one line), an assignment inside a MULTI-LINE branch
// that never executes (`if false; then\n    NAME="x"\nfi` — COUNTED and its value
// SUPPLIED, though bash leaves NAME unset), or NAME="…" text inside a here-doc body
// (which it would COUNT — a false RED, the loud direction; run.sh has no here-docs, V-K).
// Each residual is pinned by an arm of TestShellAssignmentValuesShapes.
func shellAssignmentValues(lines []string, name string) []string {
	prefix := name + "=\""
	var values []string
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if !strings.HasPrefix(trimmed, prefix) {
			continue
		}
		rest := strings.TrimPrefix(trimmed, prefix)
		if end := strings.IndexByte(rest, '"'); end >= 0 {
			values = append(values, rest[:end])
		}
	}
	return values
}
```

**Comment shapes, reasoned from the trimmed line's first byte and then TESTED, not assumed:**
after `TrimLeft(line, " \t")`, a commented assignment — `# KNOWN_BAD="…"`, `#   KNOWN_BAD="…"`,
or `    # KNOWN_BAD="…"` — begins with the byte `#`, which can never begin the prefix
`NAME + "=\""` (every scanned name starts with an ASCII letter), so `HasPrefix` is false and the
line contributes nothing. This is the same first-byte argument in both the column-0 and indented
positions, and it gets explicit unit arms (d)/(e) rather than resting on the argument.

### Call sites and inheritance

The helper has exactly one definition and four call sites (V-A), all in the same file; no other
file, script, or workflow reads the `run.sh` assignments (V-H1), and no other code constructs the
`name+"=\""` scan (V-I). Every consumer therefore inherits the fix from the one helper edit. There
is **no new assertion at any call site** — the four existing `assignment count=%d, want 1` checks
and their messages are **byte-unchanged** and are what does all of the work:

| Site | Consumer | Name |
|---|---|---|
| `:260` | `TestMiscompileInstrumentProbesPinnedToolchain` | KNOWN_GOOD |
| `:261` | same | KNOWN_BAD |
| `:262` | same | PINNED |
| `:349` | `TestReproModuleFloorStaysBelowKnownBadToolchains` | KNOWN_BAD |

The four existing `assignment count=%d, want 1` checks and their messages are **byte-unchanged**
(AC3 pins the V-E messages verbatim). Note for the sprint: those pre-existing messages hardcode
`KNOWN_BAD` for all four sites (`gemini-3-1-pro`'s finding) — a `KNOWN_GOOD` or `PINNED`
deviation would emit a false variable name. That is pre-existing and out of scope for this row;
template it as a separate queue note, not here.

### New unit test: one arm per shape

`TestShellAssignmentValuesShapes` (same file, table-driven over inline `[]string` lines — **no
on-disk fixture needed**), asserting `len(values)` exactly per arm:

| # | Arm (lines) | want len(values) | Pins |
|---|---|---|---|
| 1 | `NAME="v"` | 1 (`["v"]`) | baseline extraction |
| 2 | `  NAME="v"` (spaces) | 1 (`["v"]`) | whitespace tolerated; value EXTRACTED |
| 3 | `\tNAME="v"` (tab) | 1 (`["v"]`) | tabs are whitespace too |
| 4 | `# NAME="v"` and `#   NAME="v"` | 0 | column-0 comment ignored (shape d) |
| 5 | `    # NAME="v"` | 0 | indented comment ignored (shape e) |
| 6 | `NAME='v'` | 0 | single-quoted still contributes 0 → `count=0, want 1` floor keeps firing |
| 7 | `NAME="v"` + `  NAME="w"` | 2 (`["v","w"]`) | the V-D / ARM A silent arm, now visible → count=2 reds |
| 8 | `export NAME="v"` | 0 | DECLARED residual, pinned |
| 9 | `NAME=$OTHER` (unquoted) | 0 | DECLARED residual, pinned |
| 10 | `NAME="v` (no closing quote) | 0 | pre-existing extraction behavior pinned unchanged |
| 11 | `if false; then` + `    NAME="v"` + `fi` (multi-line, never taken) | 1 (`["v"]`) | DECLARED residual 6, pinned: counted and supplied though bash leaves NAME unset |

### What the new parser still CANNOT see (declared residuals)

Each entry below is (i) written into the helper's `DECLARED RESIDUAL` comment, (ii) pinned by a
unit arm above, and (iii) carried by AC8 — a named gap, not an assumed one:

1. **`export KNOWN_BAD="…"` / `local KNOWN_BAD="…"`** — not covered. The trimmed line starts
   with `e`/`l`, not the name. *Replacement* direction is loud (column-0 count drops to 0 → the
   existing `count=0, want 1` floor). *Addition* direction (an indented or top-level `export`
   beside the plain assignment) remains silent — accepted: covering it means normalizing
   arbitrary declaration prefixes, out of proportion for this row. This is `oc-glm-5-2`'s
   non-blocking objection, and it leaves the IDENTICAL silent-narrowing hole open under EITHER
   candidate rule — carried verbatim as a declared, named residual. Pinned by arm 8.
2. **Unquoted / variable-valued `KNOWN_BAD=$OTHER`** — contributes 0, unchanged. Loud on
   replacement via the same floor; silent on addition. Pinned by arm 9.
3. **Single-quoted `KNOWN_BAD='…'`** — contributes 0, **deliberately unchanged**: the row
   depends on this staying a loud `count=0, want 1` replacement floor. Pinned by arm 6.
4. **Mid-line assignment after another token** (`if …; then KNOWN_BAD="x"; fi` on one line) —
   invisible; the trimmed line starts with `if`. Declared only (no arm; the shape space of
   "assignment somewhere mid-line" is unbounded and every bounded probe of it would be
   decorative).
6. **Assignment inside a MULTI-LINE branch that never executes** — e.g.
   `if false; then` / `    KNOWN_BAD="…"` / `fi` as three lines, with no column-0 assignment
   anywhere. Under the REJECTED rule B this redded loudly (column-0 count 0 → `count=0, want 1`).
   Under ratified rule A it is **counted as 1 and its value is SUPPLIED**, so the two consumers
   pass and the toolchain floor is bound against a deny-list that **does not exist at runtime**.
   Measured first-party by the controller at `cb73cab` (V-ARMC1/V-ARMC2 below), including the
   runtime control: `bash -c 'if false; then KNOWN_BAD="x"; fi; echo "${KNOWN_BAD:-<UNSET>}"'`
   prints `<UNSET>`. **This is a real narrowing of loudness that ruling A accepts**, and it is
   strictly wider than residual 4: residual 4 covers only the ONE-LINE form, while the
   multi-line form is the shape a refactorer actually writes. It is declared rather than fixed
   because closing it needs branch-reachability analysis, which is out of proportion for a ~0.1d
   row — and because the row's own defect (a SECOND assignment silently narrowing the list) is
   fully closed either way, measured at `count=2, want 1`. Pinned by arm 11.
5. **Here-doc body text** matching `NAME="…"` at (possibly indented) line start — after this
   change it would be *counted* (count=2), i.e. a **false RED**, which is the loud direction of
   the asymmetry and therefore acceptable. `run.sh` contains no here-docs today (V-K).

## Files to Modify/Create

- `host/verifygate/toolchain_pin_gate_test.go` — the ONLY file touched (~+15/-5 LOC): the helper
  body becomes one whitespace-tolerant scan with the residual comment; signature unchanged; **no
  call-site edits, no new assertions**; new `TestShellAssignmentValuesShapes` table test.

No new files, no fixtures, no new dependencies, **no `.ail` file in scope** (this row proposes no
AILANG code), no script or workflow edits, nothing under `tools/launchd/**`.

**S3 "why is this not a package?":** this is gate instrumentation — a test-only edit to the host
verify surface (`host/verifygate/`), touching no kernel type, no transition law, and exporting no
API. There is no package-shaped surface here; the change cannot land anywhere else.

## Examples

**Before (today, measured V-D / ARM A):** land the conditional-narrowing mutant, then run cmd V-C
(Verification Log, printed in full there) → rc=0, RUN=2 PASS=2.
The floor test certifies `repro/go.mod`'s floor against all four deny-listed toolchains while the
script would probe one.

**After (this item):** the same mutant → rc=1, and BOTH consumers name the cause with the
pre-existing message:

```
--- FAIL: TestMiscompileInstrumentProbesPinnedToolchain
    …run.sh: KNOWN_BAD assignment count=2, want 1
--- FAIL: TestReproModuleFloorStaysBelowKnownBadToolchains
    instrument failure: …run.sh: KNOWN_BAD assignment count=2, want 1
```

**Deliberately changed (V-E / ARM B shape):** an indented-ONLY assignment is now **GREEN**
(rc=0, RUN=2, PASS=2) — the value is extracted and supplied, which is what bash actually does.
This is the ratified consequence of ruling A: today's red fires because the scanner cannot SEE a
real assignment, not because anything is wrong. Nothing else got quieter — single-quoted and empty
shapes still red with the byte-unchanged `count=0, want 1` messages.

## Acceptance Criteria

Every AC names a command and an observable that can FAIL. Run-existence is asserted on `=== RUN`
/ `--- PASS` counts, never exit codes alone (an rc=0 `[no tests to run]` is a FALSE PASS). Shell
recipes follow the rig rules: capture with `cmd > /tmp/out 2>&1; rc=$?` (zsh `${PIPESTATUS[0]}`
is empty), never `|| echo 0` inside `$(...)`, `export PATH=/opt/homebrew/bin:$PATH` first.
**Doc-wide convention:** every command in this document is copy-paste runnable under zsh exactly
as printed, which is why runnable commands live in fenced code blocks and never in markdown
table cells — a bare `|` breaks a table cell, and the `\|` escape that fixes the rendering
silently corrupts the command (in ERE and Go regexp `\|` matches a *literal* pipe); that
table-cell escaping is exactly why an earlier revision's commands were not runnable as printed.

- **AC1 (run-existence, unit arms).**
  `go test ./host/verifygate/ -run TestShellAssignmentValuesShapes -v > /tmp/ac1 2>&1; rc=$?` →
  rc=0 AND `grep -c -- '--- PASS' /tmp/ac1` ≥ 11 (one per table arm, subtests counted) AND
  `grep -c '=== RUN' /tmp/ac1` ≥ 1. A `[no tests to run]` line anywhere in `/tmp/ac1` FAILS this AC.
- **AC2 (the silent arm is closed — the V-D / ARM A mutant reds).** With the canonical addition
  mutant landed in `run.sh` (protocol below; effect gates G1 = 1 AND G2 = 1 — gate commands in the
  fenced block under the mutation table):
  `go test ./host/verifygate/ -run 'TestMiscompileInstrumentProbesPinnedToolchain|TestReproModuleFloorStaysBelowKnownBadToolchains' -v > /tmp/ac2 2>&1; rc=$?`
  → rc=1, `grep -c -- '--- FAIL' /tmp/ac2` = 2, and
  `grep -c 'KNOWN_BAD assignment count=2, want 1' /tmp/ac2` = 2 (both consumers, by name). The
  row's entire purpose is served by the pre-existing message, with NO new assertion at any call
  site.
- **AC3 (REVERSED by the ruling — the indented-ONLY arm is GREEN).** With the indented-ONLY
  mutant landed (effect gates: G1 = 1 AND G2 = 0): the same test command captured to `/tmp/ac3`
  (`> /tmp/ac3 2>&1; rc=$?`) → **rc=0, RUN=2 PASS=2 FAIL=0**. This is intended and ratified: the
  value is extracted and supplied, which is what bash actually does; today's red fired because the
  scanner could not SEE a real assignment, not because anything was wrong.
- **AC4 (commented shapes stay green).** With mutant (d) landed (a column-0 `# KNOWN_BAD="…"`
  line added; effect gate G5 moves 0→1) the two consumers pass
  rc=0 RUN=2 PASS=2; likewise mutant (e) (indented comment; G6 moves 0→1). Both are
  LANDED-asserted by sha256 like every other arm — a green arm proves ignoring,
  only if the bytes demonstrably changed.
- **AC5 (single-quoted floor keeps firing).** With mutant (c) landed (column-0 double-quoted
  KNOWN_BAD replaced by `KNOWN_BAD='…'`; effect gates G3 0→1 and G4 1→0): rc=1 with
  `KNOWN_BAD assignment count=0, want 1` at both sites.
- **AC6 (hygiene).** `gofmt -l host/ cmd/ > /tmp/ac6; rc=$?` → `/tmp/ac6` is **0 bytes**
  (`wc -c` = 0); `go vet ./... ` rc=0; `go build ./...` rc=0.
- **AC7 (gates — `.ail` floor by exit code, Go gate by SET comparison, plus the narrowest gate
  that can actually fail).** Three parts, under the pinned `AILANG_BIN=/tmp/ailang-v0300/ailang`:
  1. `./scripts/verify_ail.sh` rc=0 — green at pristine base (controller-measured), so the
     exit-code criterion stands.
  2. `./scripts/verify_go.sh` — an **rc=0 criterion is FORBIDDEN here**, because the gate is RED
     at pristine base on this rig (an already-red gate measures the repo, not the change; full
     measurement in Non-Goals). The criterion is a SET comparison: the set of `--- FAIL` test
     names after the change is **identical to the base set below**, AND `host/verifygate`
     reports `ok` in the same run. Print both sets:

     ```zsh
     ./scripts/verify_go.sh > /tmp/ac7-go 2>&1; rc=$?
     grep -E '^--- FAIL' /tmp/ac7-go | sort        # must match the base set, name for name
     grep -E '^(ok|FAIL).*host/verifygate' /tmp/ac7-go   # must be an `ok` line
     ```

     **The base set is measured AT DRILL START, in the same worktree and at the same commit —
     never quoted from a doc** (`oc-glm-5-2`, quorum round 3, verbatim alternative fix; adopted
     because it is also the only form robust to queue row 58's finding that this gate is FLAKY
     on this rig, not deterministically red — four runs within ~90 minutes on trees differing by
     one test file gave FAIL counts of 4, 1, 1 and 0). The previously hardcoded set was measured
     at `9c0ad0b`, and this doc targets `cb73cab`; those are different commits, so a planner
     executing the old criterion could fail the AC for reasons unrelated to this change, or pass
     it vacuously if the base set shrank. So:

     ```zsh
     # BEFORE any edit, in the drill worktree, at the same commit the drill will run on:
     ./scripts/verify_go.sh > /tmp/ac7-base 2>&1; echo "base rc=$?"
     grep -E '^--- FAIL' /tmp/ac7-base | sort > /tmp/ac7-base-set
     wc -l < /tmp/ac7-base-set        # record this number in the sprint evidence row
     ```

     The criterion is then `diff /tmp/ac7-base-set <(grep -E '^--- FAIL' /tmp/ac7-go | sort)`
     being EMPTY. Because the gate is flaky, take the base reading TWICE and, if the two differ,
     record both and treat only names present in BOTH as the stable base — a name that moves
     between two pristine runs is the gate's flake, not this diff's doing, and attributing it
     here would be exactly the co-occurrence error rule 3d names. Any name appearing or
     disappearing relative to the stable base fails this AC.
  3. The gate that CAN fail for this diff — the narrowest gate that can fail is preferred over
     the widest that looks thorough: `go test ./host/verifygate/ > /tmp/ac7-narrow 2>&1; rc=$?`
     → rc=0.

  And **0 `.ail` files touched** (this row proposes no AILANG code):

  ```zsh
  git diff --name-only "$(git merge-base dev HEAD)"..HEAD > /tmp/ac7; rc=$?
  grep -c '\.ail$' /tmp/ac7; rc=$?   # 0 matches, rc=1 (grep's legitimate-zero exit; do NOT wrap in || echo 0)
  wc -l < /tmp/ac7                   # >= 1 changed path — the same file is its own positive control
  ```
- **AC8 (residuals are declared, not assumed).**
  `grep -c 'DECLARED RESIDUAL' host/verifygate/toolchain_pin_gate_test.go` increases by ≥1 over
  the baseline of 3 (V-J), AND `grep -n 'export' host/verifygate/toolchain_pin_gate_test.go`
  hits the helper's residual comment, AND unit arms 6/8/9/10/11 exist (subtest names in `/tmp/ac1`).
- **AC9 (per-name, per-consumer inheritance — no new assertion).**
  `grep -c 'shellAssignmentValues' host/verifygate/toolchain_pin_gate_test.go` = 5 (1 def + 4
  calls, V-A); the helper signature is unchanged
  (`func shellAssignmentValues(lines []string, name string) []string`); and `git diff` for this
  sprint shows no hunk touching any of the four `assignment count=%d, want 1` sites
  (`:263-271`, `:350-355`). All four consumers inherit the fix from the single helper edit.
- **AC10 (pristine control).** On the untouched fixture (sha256 equal to V-B/V-F's
  `b80109aa5788…` baseline): the two consumers pass rc=0 RUN=2 PASS=2 before AND after the drill.

## Conflict Surface

**The reuse question, audited — not waived:** does this repo already carry shell-parsing
machinery (a parser, tokenizer, library, or dependency) that this row should route to instead of
extending `shellAssignmentValues`? Measured repo-wide (V-L, V-M, V-N; commands and outputs in the
Verification Log):

- **Dependencies (V-L):** `go.mod` declares exactly ONE direct dependency,
  `modernc.org/sqlite v1.54.0`; every other module line is `// indirect`. No shell-parsing
  library (mvdan.cc/sh or any other) exists anywhere in the module graph.
- **Repo-wide source search (V-M):** the pattern
  `mvdan|sh/syntax|shellwords|shlex|shell.?lex|shell.?pars|tokeni[sz]e` over every `*.go` and
  `*.mod` file returns **4 hits**, ALL of them the single symbol `tokenize` in
  `host/runbook/runbook_stageb_test.go` (doc comment `:473`, definition `:478`, calls `:609` and
  `:636`); **zero** hits for every library name. Same-scope known-positive control: `HasPrefix`
  returns **59** hits, so the instrument can see a positive.
- **What the single hit is (V-N):** `host/runbook/runbook_stageb_test.go:473-478` — "tokenize
  splits ONE runbook command line into argv, honouring double quotes". It is (i) a **different
  grammar** — argv splitting of a command line, not detection of a `NAME="…"` assignment at line
  start; (ii) **unexported**; and (iii) in a **different package's `_test.go` file**, so it is
  not importable by `host/verifygate` even if the grammar matched.

**Conclusion: nothing to reuse.** One adjacent tokenizer exists, with a different grammar,
unexported, in another package's test file. This row stays a deliberately narrow prefix grammar
local to one gate test, and its complete consumer set is enumerated in-file (V-A, V-H1). Under
ruling A the change is strictly SMALLER than the rejected alternative (one helper body, no
call-site edits), so the reuse surface is unchanged and the conclusion holds a fortiori.

Two adjacent scanners in the SAME file are deliberately NOT touched and must stay byte-identical:
`moduleGoFloor` (`:86`, whose column-0 `"go "` anchor is load-bearing —
`w-racecontrol-floor-bump-disarms-the-race-control` V8 depends
on comments being invisible to it) and `keyedValues` (`:70`, YAML `key:` grammar, a different
language). An AC-grade check: `git diff` for this sprint shows no hunk overlapping either
function.

## Testing Strategy — Non-Vacuity: named mutation, one arm per SHAPE

**Why an ADDITION arm is mandatory:** a removal mutant proves a check FIRES; only an addition
mutant proves the check LOOKS. This row exists because the gate's checks all fired correctly on
removals and replacements (V-E) while never looking at additions (V-D). Arm (a) is the addition.

Per the doc convention, the table cells reference named gate commands G1–G9; the commands
themselves are in the fenced block below the table, runnable as printed.

The controller measured ARM 0 / ARM A / ARM B at HEAD (`cb73cab`) under the REJECTED rule B. Under
ruling A the post-fix expectations below are **PREDICTED** — derived from the controller's
measurements and the ratified implementation, but not yet run under A. Rows marked **PREDICTED**
must be RUN by the planner, not inherited as measurements.

| Arm | Mutant (target) | Landed/effect gate | Expected post-fix (under A) | Kills which mutation / proves what |
|---|---|---|---|---|
| (a) | **ADDITION**: second `KNOWN_BAD="go1.26.0"` indented inside an `if [ -n "${NARROW_DENYLIST:-}" ]` block, beside the valid column-0 line (V-D / ARM A verbatim) | G8 sha256 moves off `b80109aa5788…`; G7 rc=0; G1 = 1 AND G2 = 1 | **RED**, both consumers, `KNOWN_BAD assignment count=2, want 1` ×2 (AC2) — **PREDICTED** | The row's silent arm — proves the tolerant scan now LOOKS at lines the old scan never saw, and that the pre-existing `count != 1` check is what closes it |
| (b) | indented-ONLY assignment (V-E / ARM B verbatim) | G1 = 1, G2 = 0; G7 rc=0 | **GREEN**, rc=0 RUN=2 PASS=2 (AC3) — **PREDICTED** | Proves the ratified behaviour change: the value is extracted and supplied, which is what bash actually does |
| (c) | single-quoted `KNOWN_BAD='…'` replacing the column-0 line | G3 0→1, G4 1→0; G7 rc=0 | **RED**, `count=0, want 1` ×2 (AC5) — **PREDICTED** | Proves quote grammar unchanged; the loud single-quote floor survives |
| (d) | commented assignment at column 0 (`# KNOWN_BAD="go1.0"` added) | G5 0→1; G8 sha moved; G7 rc=0 | **GREEN**, rc=0 RUN=2 PASS=2 (AC4) — **PREDICTED** | Proves `#` at trimmed-first-byte is ignored in the column-0 position |
| (e) | commented assignment, indented (`   # KNOWN_BAD="go1.0"` added) | G6 0→1; G8 sha moved; G7 rc=0 | **GREEN**, rc=0 RUN=2 PASS=2 (AC4) — **PREDICTED** | Proves `#` is ignored in the indented position too — the tolerance did not start counting comments |
| (f) | **neuter the tolerance itself**: in the helper, `strings.TrimLeft(line, " \t")` → `line` (Go-side) | test-file sha moved (G8 on the test file); G9 **rc=0** BEFORE reading any test result (a plain build does NOT compile test files and cannot serve here) | **RED**: `TestShellAssignmentValuesShapes` arms 2/3/7 fail (indented values want 1/1/2, got 0) — **PREDICTED** | Proves the new machinery is non-vacuous: remove the tolerance and a named test says so |
| (g1) | neuter the pre-existing `count != 1` Errorf in the probe test (`:263-271`) | test-file sha moved; G9 rc=0 | With shell arm (a) simultaneously landed and ONLY `TestMiscompileInstrumentProbesPinnedToolchain` run: the previously-expected `count=2, want 1` failure is GONE from that test's output (the floor test still reds — read WHICH test failed) — **PREDICTED** | Proves the pre-existing count check is load-bearing for the silent arm, not shadowed by its neighbor |
| (g2) | neuter the pre-existing `count != 1` Fatalf in the floor test (`:350-355`) | same G9 gate | With arm (a) landed and ONLY `TestReproModuleFloorStaysBelowKnownBadToolchains` run: rc=0 — the floor test alone goes green — **PREDICTED** | Proves the floor is a floor: without it the deny-list read certifies a shadowed list |
| (g3) | neuter a pre-existing floor: `len(badAssignments) != 1` Fatalf in the floor test (`:350`) | same G9 gate | With shell arm (c) landed and ONLY the floor test run: the `count=0, want 1` instrument failure is GONE (the later `must contain at least one toolchain` floor now catches it instead — record WHICH message fired) — **PREDICTED** | One arm per FLOOR, not per branch: proves each refusal is individually alive, and documents which floor is next in line when one dies |
| (g4) | neuter the probe test's known-positive control loop (`:245-249`, drop `"KNOWN_BAD="` from the control list) | same G9 gate | With a `run.sh` stripped of every `KNOWN_BAD` line landed and ONLY the probe test run: the `does not contain known-positive control` Fatalf is GONE; the failure shifts to `count=0, want 1` — record the shift — **PREDICTED** | The control floor is the instrument's own health check — systematically the last thing anyone pins |

**Gate commands** referenced by the table above and by ACs 2–5 (repo root, zsh, runnable as
printed):

```zsh
RS=design_docs/verification/w-race-gate-blindspot/run.sh
grep -cE '^[[:space:]]+KNOWN_BAD=' "$RS"      # G1 — indented-assignment count (effect gate for the drill)
grep -cE '^KNOWN_BAD=' "$RS"                  # G2 — column-0 assignment count
grep -cE "^KNOWN_BAD='" "$RS"                 # G3 — single-quoted column-0 count
grep -cE '^KNOWN_BAD="' "$RS"                 # G4 — double-quoted column-0 count
grep -cE '^# *KNOWN_BAD=' "$RS"               # G5 — column-0 commented count
grep -cE '^[[:space:]]+# *KNOWN_BAD=' "$RS"   # G6 — indented commented count
bash -n "$RS"                                 # G7 — mutant is valid shell (rc=0)
shasum -a 256 "$RS"                           # G8 — landed/restore assertion (run on the test file instead for Go-side arms)
go vet ./host/verifygate/                     # G9 — Go-side mutant compiles (rc=0)
```

A `grep -c` gate that legitimately counts 0 exits rc=1 — read the COUNT, never the exit code,
and never wrap in `|| echo 0`.

**Drill protocol (binding on the sprint):** every mutant asserted LANDED by sha256 before/after
(G8); the mutant's INTENDED EFFECT asserted against the system's own view (a `grep -c` that must
move — per-arm gates above), never against the file's bytes alone; Go-side mutants gate on
G9 rc=0 BEFORE any test result is read; shell-side mutants gate on
G7 rc=0; restore from a `cp` BACKUP, never `git checkout --`; restore verified
byte-identical by sha256 (`b80109aa5788…` for run.sh; the post-fix test-file sha recorded at
drill start for Go arms); the pristine control (AC10) re-run before AND after EVERY arm; the
verdict of every arm is the NAME of the failing test and the message that fired, never an exit
code alone.

## Non-Goals

- No change to `run.sh` (the fixture is mutated only inside the drill, restored byte-identical).
- No change to `moduleGoFloor`, `keyedValues`, `requireToolchainNamePin`, or any other scanner.
- No general shell parser: this remains a deliberately narrow prefix grammar with its residuals
  now *named*. Covering `export`/`local`/mid-line/here-doc shapes is out of scope (declared
  residuals 1/2/4/5).
- No new assertion at any call site, and no change to the four pre-existing
  `assignment count=%d, want 1` messages (they are byte-unchanged and do all of the work).
- No `.ail` code, no `scripts/` edits, no CI workflow edits, nothing under `tools/launchd/**`.
- If the sprint uncovers a genuinely separate defect, it is FILED as a "for the queue, not this
  sprint" note in the sprint log — not absorbed here.
- **For the queue, not this sprint — `verify_go.sh` is RED at pristine base on this rig**
  (controller-measured, untouched worktree at dev = `9c0ad0b`, TWO arms: once under concurrent
  load and once alone; rc=1 in both, with the SAME four failures both times):
  `TestEpisodeLiveReplayThreeArmsAndEvidence` (host/broker),
  `TestHandlerTimeoutKillsTheWholeProcessGroup` (host/broker),
  `TestF1PinnedInterpreterHashMismatchRefusedBeforeExec` (host/capsule),
  `TestFixtureEpisodeReplaysBitForBit` (host/replay). Three of the four share ONE mechanism:
  `cannot obtain --version from archived interpreter … timed out after 10s`. Mechanism isolated:
  a freshly copied binary at a new path costs **1336 ms** on first exec vs **96 ms** for the same
  binary at its pinned path (**372 ms** on a cached second exec) — macOS first-exec
  code-signing/provenance assessment (`com.apple.provenance` xattr present; `spctl -a -t exec`
  rejects) — on top of a per-invocation observatory retention cleanup over a 513 MB DB. This is
  a RIG-LOCAL environmental cost, not a code defect: dev CI is GREEN on this same commit (2/2
  checks) on the Linux runners, and `./scripts/verify_ail.sh` is rc=0 at base. AC7 is written
  around this fact (set comparison, never exit code); fixing the rig cost is a queue item, not
  this sprint's scope — the next iteration inherits a finding, not a mystery.

## Timeline / Milestones

Single milestone, ~0.1d: (1) helper body + unit table (~30 min); (2) mutation drill, all arms
with restores (~30 min); (3) gates + ACs recorded (~15 min). Under ruling A the change is
strictly SMALLER than the rejected alternative — one helper body, no call-site edits, no new
assertions — so the milestone is correspondingly lighter.

## Risks & Mitigations

- **Risk:** the indented-ONLY shape going GREEN is read as a regression by a future reviewer.
  **Mitigation:** it is the ratified consequence of ruling A (D-WORLD-29, attended 2026-09-01),
  recorded verbatim in the status block and AC3; today's red fired because the scanner could not
  SEE a real assignment, not because anything was wrong.
- **Risk:** a benign re-indent that ADDS a second assignment (column-0 + indented) now reds with
  `count=2, want 1`. **Mitigation:** that is the row's entire purpose — a refactor that adds a
  conditional branch reds loudly instead of narrowing the deny-list in silence. The message is
  the pre-existing one, byte-unchanged.
- **Risk:** the whitespace-tolerant scan starts counting here-doc body text or comments.
  **Mitigation:** comments are excluded by the first-byte argument (a `#` can never begin
  `NAME + "=\""`), pinned by unit arms (d)/(e); here-doc text would be a false RED (the loud
  direction, acceptable) and `run.sh` has no here-docs at baseline (V-K).
- **Risk:** drill arm (g1)/(g2) misread because the OTHER consumer's red masks the neutered one.
  **Mitigation:** the protocol runs each (g) arm against ONLY the test containing the neutered
  floor and records which message fired.

## Verification Log

Rows V-A..V-G are the controller's first-party measurements (iteration 140, this worktree),
reused. Rows V-H..V-K were measured by the designer in this worktree, 2026-08-31. Rows V-L..V-N
are the conflict-surface reuse audit (designer, 2026-08-31; results independently re-derived by
the controller in this worktree). Rows V-ARM0 / V-ARMA / V-ARMB are the controller's first-party
measurements at HEAD this iteration (`cb73cab`), taken under the REJECTED rule B; their
conclusions are re-derived under ruling A in the mutation table. Negative results carry a
known-positive control in the same call, same path. Per the doc-wide convention, **no runnable
command sits in a table cell**: each row's command is in the matching fenced block under
"Commands", byte-for-byte as executed.

**Transcription-defect note (quorum round 1, objection upheld in part):** the previous revision
printed the V-B, V-C, and V-K commands inside table cells with markdown-escaped pipes (`\|`). In
ERE — and in Go's `-run` regexp — `\|` matches a *literal* pipe, so those commands AS PRINTED
could not have produced the recorded observations (measured: the escaped V-K form returns rc=1
with 0 hits; the true-alternation form returns rc=0 with the 3 recorded hits at lines 24/25/26).
The observations themselves are TRUE: the controller independently re-derived them, and the
repaired V-B, V-C, and V-K commands were re-run in this worktree on 2026-08-31 with the outputs
recorded below. The defect was transcription — a bare `|` breaks a markdown table cell, and the
escape that fixes the rendering corrupts the command — not fabrication. The same audit caught
two count-vs-enumeration slips (V-A said "4 hits" while listing 5 grep lines; V-H1 said 9 while
listing 10) — in every case the enumeration and the conclusion stand, the counts are corrected
below, and the structural fix is that no runnable command lives in a table cell anymore.

| # | Claim | Observed (command in block V-x below) |
|---|---|---|
| V-A | `shellAssignmentValues` is defined once, column-0 anchored, with exactly 4 call sites | 5 grep lines = 1 def + 4 calls: def `host/verifygate/toolchain_pin_gate_test.go:215`; calls `:260` (KNOWN_GOOD), `:261` (KNOWN_BAD), `:262` (PINNED) in `TestMiscompileInstrumentProbesPinnedToolchain`, `:349` (KNOWN_BAD) in `TestReproModuleFloorStaysBelowKnownBadToolchains`. Body matches `strings.HasPrefix(line, name+"=\"")` |
| V-B | `run.sh` sets each name exactly once, at column 0 — the gap is LATENT | lines 24 / 25 / 26 respectively, rc=0 (re-run 2026-08-31 with the repaired command) |
| V-C | Pristine baseline: both consumers pass | rc=0, RUN=2 PASS=2 FAIL=0 (re-run 2026-08-31 with the repaired command) |
| V-D | **The silent arm, measured under B**: second indented `KNOWN_BAD="go1.26.0"` inside an `if [ -n "${NARROW_DENYLIST:-}" ]` block beside the valid column-0 line → both consumers GREEN | **rc=0, RUN=2 PASS=2 FAIL=0** — the gate is blind under B; the floor test binds the repro floor against the WIDE list `go1.26.0 go1.26.3 go1.26.4 go1.26.5` while run.sh would narrow to `go1.26.0`. Mutant landed sha256 `b80109aa5788…` → `2a703e885a6f…`; G7 rc=0; G1 = 1, G2 = 1. **Under ruling A this mutant reds** via `KNOWN_BAD assignment count=2, want 1` at both consumers (PREDICTED — see mutation arm (a)) |
| V-E | **The loud control, measured under B**: the same assignment indented with NO column-0 one | **rc=1, RUN=2 PASS=0 FAIL=2**: `toolchain_pin_gate_test.go:267: …: KNOWN_BAD assignment count=0, want 1` and `:351: instrument failure: …: KNOWN_BAD assignment count=0, want 1` — the asymmetry is measured, not argued. G1 = 1, G2 = 0. **Under ruling A this arm is GREEN** (rc=0, RUN=2, PASS=2) — deliberate, ratified (PREDICTED — see mutation arm (b)) |
| V-F | Restore verified byte-identical; pristine control re-passes | sha256 back to `b80109aa5788…`; rc=0 RUN=2 PASS=2; porcelain 0 lines |
| V-G | Gate baseline on the untouched worktree (build/vet/fmt) | rc=0; rc=0; 0 bytes. (`verify_go.sh` / `verify_ail.sh` baselines are controller-recorded — see AC7 and Non-Goals for the verify_go.sh base-red measurement) |
| V-H1 | No file outside the gate test reads the `run.sh` assignments — the helper fix covers every reader | **10** hits (controller re-derived 2026-08-31; an earlier revision said 9 while listing 10 line numbers — the count was a transcription slip, the enumeration and conclusion unchanged), ALL in `host/verifygate/toolchain_pin_gate_test.go` (`:245,:261,:267,:283,:301,:349,:351,:355,:360,:367`); 0 hits in `cmd/`, `scripts/`, `.github/` — `verify_go.sh` and `ci.yml` never re-scan the names |
| V-H2 | Helper body read: value ends at first `"`, append only when a closing quote exists | `prefix := name + "=\""`; `HasPrefix(line, prefix)`; `IndexByte(rest, '"')` guard — arm 10 of the unit table pins this unchanged (read, not grepped: `toolchain_pin_gate_test.go:215-227`) |
| V-I | The `name+"=\""` scan is not re-implemented anywhere else in the file; adjacent scanners use different grammars and stay untouched | 6 hits: `:31` (`"go"` version canon), `:86` (`"go "` module floor — column-0 anchor is load-bearing per the racecontrol doc's V8), `:189` (`"toolchain "`), `:219` (the helper — the known-positive control), `:511`/`:524` (YAML `- name:`) — none scans `NAME="` |
| V-J | Consumer checks read first-hand: 3 Errorf `count…want 1` sites `:263-271`, non-empty checks `:279-284`, floor-test Fatalf floors `:350-355`; `DECLARED RESIDUAL` baseline count for AC8 | Checks as described (the `strings.Fields` split at `:274-278`/`:353` consumes only column-0 values — unaffected); DECLARED RESIDUAL baseline = 3 (read `:229-369` plus block V-J's grep) |
| V-K | Baseline has ZERO indented name-matches and ZERO here-docs in `run.sh` — the new checks are green at baseline; residual 5 is currently unoccupied | Exactly lines 24/25/26 hit, all column-0, rc=0 — zero indented; only `:22` (`dirname`) hits — zero `<<` here-doc markers (re-run 2026-08-31 with the repaired commands) |
| V-L | The module graph carries NO shell-parsing dependency | Exactly ONE direct dependency: `modernc.org/sqlite v1.54.0`; every other module line is `// indirect` |
| V-M | Repo-wide, the only shell-parsing-shaped symbol is one test-file `tokenize`; the instrument can see a positive | 4 hits, ALL the single symbol `tokenize` in `host/runbook/runbook_stageb_test.go` (doc comment `:473`, def `:478`, calls `:609`, `:636`); zero hits for every library name. Same-scope control `HasPrefix`: 59 hits |
| V-N | The single V-M hit is a DIFFERENT grammar, unexported, in another package's `_test.go` — nothing to reuse | "tokenize splits ONE runbook command line into argv, honouring double quotes … a tiny, total subset of shell"; `func tokenize(t *testing.T, line string, session map[string]string) []string` — argv splitting, not line-start assignment detection; unexported; not importable by `host/verifygate` |
| V-ARM0 | **Pristine control at `cb73cab`** (before AND after the batch), package `host/verifygate`, arms `-run 'TestMiscompileInstrumentProbesPinnedToolchain\|TestReproModuleFloorStaysBelowKnownBadToolchains'` | vet rc=0, test rc=0, RUN=2 PASS=2 FAIL=0 |
| V-ARMA | **The silent arm at `cb73cab`** — a SECOND, INDENTED `KNOWN_BAD="go1.26.0"` added beside the column-0 one, inside `if true; then … fi` | mutant LANDED (sha256 differs), `bash -n` rc=0, intended effect asserted against the system's own view (indented `KNOWN_BAD="` lines 0 → 1), vet rc=0, test rc=0, RUN=2 PASS=2 FAIL=0 — **the gate is blind under B** while the deny-list is genuinely narrowed at runtime (bash takes the last assignment). Under ruling A this mutant reds via `count=2, want 1` (PREDICTED) |
| V-ARMB | **The loud arm at `cb73cab`** — the assignment indented ONLY, no column-0 copy | mutant LANDED, `bash -n` rc=0, intended effect col0 1 → 0 and indented 0 → 1, vet rc=0, test rc=1, RUN=2 PASS=0 FAIL=2, with `KNOWN_BAD assignment count=0, want 1` at BOTH consumer sites (`toolchain_pin_gate_test.go:267` and `:351`). Under ruling A this arm is GREEN (rc=0, RUN=2, PASS=2) — deliberate, ratified (PREDICTED) |
| V-ARMC1 | **The conditional-branch shape under the CURRENT (base) helper**, at `cb73cab`: the ONLY `KNOWN_BAD` moved inside a never-taken `if false; then` / `fi` multi-line branch, no column-0 copy | mutant LANDED (sha256 differs), `bash -n` rc=0, intended effect col0 1→0 and indented 0→1, `go vet` rc=0 read BEFORE the test result; **test rc=1, RUN=2 PASS=0 FAIL=2**, `KNOWN_BAD assignment count=0, want 1` at `:267` and `:351`. Runtime control in the same batch: `bash -c 'if false; then KNOWN_BAD="x"; fi; echo "${KNOWN_BAD:-<UNSET>}"'` prints **`<UNSET>`**, so today's red is CORRECT for this shape |
| V-ARMC2 | **The same shape under the RATIFIED rule-A helper** (the Overview body applied verbatim to a scratch worktree), at `cb73cab` | helper LANDED (sha256 differs), intended effect asserted against the system's own view (`TrimLeft` present = 1), `go vet` rc=0 read BEFORE the test result; **test rc=0, RUN=2 PASS=2 FAIL=0** — the gate accepts a deny-list bash never sets. Paired control in the same batch, proving the helper is not simply permissive: the row's own SILENT arm under the SAME helper is **rc=1, RUN=2 PASS=0 FAIL=2**, `KNOWN_BAD assignment count=2, want 1` at `:269` and `:353`. Both files restored byte-identical (sha256 re-asserted), worktree porcelain 0. **This is what residual 6 declares** |
| V-ARM-GO | **`./scripts/verify_go.sh` base behaviour at `cb73cab`**, TWICE on the IDENTICAL pristine worktree (`git status --porcelain` = 0 between runs), answering `oc-glm-5-2`'s round-3 premise objection by measurement rather than forwarding it | **run 1: rc=1**, `--- FAIL` set = `{TestHandlerTimeoutKillsTheWholeProcessGroup}` (count 1). **run 2: rc=0**, `--- FAIL` set = **empty** (count 0). Known-positive control in the same call: 19 `ok`/`---` lines, so the empty set is a measurement and not a broken instrument. **Two conclusions.** (a) The doc's previously hardcoded 4-name base set (measured at `9c0ad0b`) matches **NEITHER** run at `cb73cab` — the objection is CONFIRMED. (b) The gate is **FLAKY, not deterministically red**: one tree, two consecutive runs, two different verdicts — corroborating queue row 58 first-party. The stable base set (names present in BOTH runs) is therefore **EMPTY**, which is why AC7 part 2 now takes its base reading at drill start and treats only names common to two pristine runs as the base |

### Commands (runnable as printed)

Repo-root-relative, zsh, `export PATH=/opt/homebrew/bin:$PATH` first. Observed outputs are
restated as `#` comments so each block stays copy-paste safe.

```zsh
# V-A
grep -rn "shellAssignmentValues" host/ --include='*.go'
# observed (re-run 2026-08-31): 5 lines — :215 (def), :260, :261, :262, :349
```

```zsh
# V-B — true alternation; the earlier table cell's `\|` form matches a literal pipe (0 hits)
grep -nE 'KNOWN_BAD=|KNOWN_GOOD=|PINNED=' design_docs/verification/w-race-gate-blindspot/run.sh
# observed (re-run 2026-08-31), rc=0:
#   24:KNOWN_BAD="go1.26.0 go1.26.3 go1.26.4 go1.26.5"
#   25:KNOWN_GOOD="go1.26.6 go1.25.6 go1.24.9"
#   26:PINNED="go1.26.6"   # the toolchain go.mod pins; TestMiscompileInstrumentProbesPinnedToolchain
```

```zsh
# V-C — the two-consumer test command (also the test command of V-D/V-E/V-F and ACs 2–5);
# true alternation in the -run selector (Go regexp treats `\|` as a literal pipe too)
go test ./host/verifygate/ -run 'TestMiscompileInstrumentProbesPinnedToolchain|TestReproModuleFloorStaysBelowKnownBadToolchains' -v > /tmp/vc 2>&1; rc=$?
grep -c '=== RUN' /tmp/vc; grep -c -- '--- PASS' /tmp/vc
# observed at pristine baseline (re-run 2026-08-31): rc=0, RUN=2, PASS=2, FAIL=0
```

```zsh
# V-D / V-E / V-F — mutation arms: land the mutant, assert gates G1/G2/G7/G8 (gate block above),
# then run the V-C block. Observations are the controller's iteration-140 first-party
# measurements, restated in the table rows; no new observation is claimed here.
# V-F restore check, additionally:
git status --porcelain
# observed: 0 lines
```

```zsh
# V-ARM0 / V-ARMA / V-ARMB — controller first-party at HEAD this iteration (cb73cab),
# package host/verifygate, arms -run 'TestMiscompileInstrumentProbesPinnedToolchain|TestReproModuleFloorStaysBelowKnownBadToolchains'
go vet ./host/verifygate/
go test ./host/verifygate/ -run 'TestMiscompileInstrumentProbesPinnedToolchain|TestReproModuleFloorStaysBelowKnownBadToolchains' -v > /tmp/varm 2>&1; rc=$?
grep -c '=== RUN' /tmp/varm; grep -c -- '--- PASS' /tmp/varm; grep -c -- '--- FAIL' /tmp/varm
# observed: ARM0 — vet rc=0, test rc=0, RUN=2 PASS=2 FAIL=0 (pristine, before AND after the batch)
#           ARMA — mutant LANDED, bash -n rc=0, indented KNOWN_BAD=" lines 0->1, vet rc=0,
#                  test rc=0, RUN=2 PASS=2 FAIL=0 (gate blind under B)
#           ARMB — mutant LANDED, bash -n rc=0, col0 1->0 and indented 0->1, vet rc=0,
#                  test rc=1, RUN=2 PASS=0 FAIL=2, `KNOWN_BAD assignment count=0, want 1`
#                  at :267 and :351
# Both mutants restored byte-identical (sha256 re-asserted against the pre-mutation backup).
```

```zsh
# V-G
go build ./...
go vet ./...
gofmt -l host/ cmd/
# observed: rc=0; rc=0; empty output (0 bytes)
```

```zsh
# V-H1 (positive control: the known test-file sites must hit)
grep -rn 'KNOWN_BAD' host/ cmd/ scripts/ .github/
# observed: 10 hits, all in host/verifygate/toolchain_pin_gate_test.go; 0 elsewhere
```

```zsh
# V-I (positive control: :219, the helper itself, must hit)
grep -n 'HasPrefix' host/verifygate/toolchain_pin_gate_test.go
# observed: 6 hits — :31, :86, :189, :219, :511, :524
```

```zsh
# V-J (the read is Read/editor, not a command; the AC8 baseline is)
grep -c 'DECLARED RESIDUAL' host/verifygate/toolchain_pin_gate_test.go
# observed: 3
```

```zsh
# V-K — true alternation; the earlier `\|` form was measured rc=1, 0 hits (see note above)
grep -nE '^[[:space:]]+(KNOWN_GOOD|KNOWN_BAD|PINNED)="|^(KNOWN_GOOD|KNOWN_BAD|PINNED)="' design_docs/verification/w-race-gate-blindspot/run.sh
grep -nE '<<|dirname' design_docs/verification/w-race-gate-blindspot/run.sh
# observed (re-run 2026-08-31): first — rc=0, exactly lines 24/25/26, all column-0, zero indented
# (the 3 column-0 lines are the positive control); second — rc=0, only
# 22:cd "$(dirname "$0")/repro" || exit 1  (the positive control); zero `<<` here-doc markers
```

```zsh
# V-L — everything filtered out is `// indirect`, so what remains is the direct-dependency set
grep -vE '// indirect' go.mod
# observed (2026-08-31): module line, `go 1.26.6`, `require modernc.org/sqlite v1.54.0`, and an
# emptied `require ( … )` block — exactly ONE direct dependency, no shell-parsing library
```

```zsh
# V-M — repo-wide shell-parsing machinery search + same-scope known-positive control
grep -rnE 'mvdan|sh/syntax|shellwords|shlex|shell.?lex|shell.?pars|tokeni[sz]e' --include='*.go' --include='*.mod' .
grep -rn 'HasPrefix' --include='*.go' --include='*.mod' . | wc -l
# observed (2026-08-31): search — rc=0, 4 hits, all `tokenize` in
# host/runbook/runbook_stageb_test.go (:473 doc comment, :478 def, :609 and :636 calls);
# control — 59
```

```zsh
# V-N — what the single V-M hit is
sed -n '473,478p' host/runbook/runbook_stageb_test.go
# observed (2026-08-31): the doc comment "tokenize splits ONE runbook command line into argv,
# honouring double quotes … deliberately a tiny, total subset of shell" and the unexported
# signature func tokenize(t *testing.T, line string, session map[string]string) []string
```

## Related Documents

- `design_docs/planned/w-racecontrol-floor-bump-disarms-the-race-control.md` — same gate file,
  same fixture; its V8 is why `moduleGoFloor`'s column-0 anchor must NOT inherit this change.
- `design_docs/verification/w-race-gate-blindspot/run.sh` — the fixture whose contract this row
  hardens (P42 lineage; the row's finder).
- `design_docs/coding-standards.md` S6 (honest gates: this row is an S6 instance — a check that
  passes vacuously on the one shape it never looks at).

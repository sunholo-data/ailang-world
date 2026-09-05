# w-load-bearing-criteria-need-a-mutation-not-a-grep

**Status**: **PLANNED** — 2026-09-04, HEAD `79d80d9`. Queue row 59 (defect A) + queue row 50
(defect B), folded into ONE design per the attended ruling `D-WORLD-31` (2026-09-03): row 50's
"option B fixture migration" is funded here, not as a doc of its own, because both rows are the
same class demonstrated with the same construction, and buying it twice is waste.

> **⚠ READ THIS BEFORE IMPLEMENTING.** This document exists to kill a specific vacuous-proof
> pattern, and it must not ship a fresh instance of the pattern it kills. Row 51's gate discharged
> a load-bearing claim with `grep -c 'disagree on their site set' … ≥ 3` — a count of the
> assertion's own message text. A controller reproduced first-party that wrapping the entire
> assertion block in `if false { … }` leaves that grep reading 3 while the test returns `--- PASS`
> under a mutant that should be RED. **The shipped code is correct; the proof obligation is
> vacuous.** The generalisation this document binds: *a criterion that greps for an assertion
> measures that somebody typed it; only a mutation measures that anybody runs it.* Every acceptance
> criterion and every mutation row below is written to that rule. Do not "fix" this document by
> adding a `grep -c <assertion text> == N` discharge — that is the exact defect it exists to kill.

**Target**: current iteration (World iter-154 / queue rows 59 + 50)
**Priority**: P1 (row 59 is a live vacuous gate; row 50 is a latent silent-narrowing hole)
**Estimated**: ~1.5d across three ≤1d milestones
**Dependencies**: none (both rows are gated on nothing; row 50's doc is superseded by this one)
**Planner-Lane**: codex-ok

---

## Problem Statement

Two defects, one cure, deliberately bought together.

**DEFECT A (queue row 59) — a static grep cannot prove an assertion is live.** Row 51 shipped a
gate whose acceptance criterion AC2 discharged the claim "this new assertion block is
load-bearing" with `grep -c 'disagree on their site set' host/verifygate/floor_raise_inventory_test.go`
≥ 3 — i.e. a count of the assertion's own message text. Read first-party at
`design_docs/planned/w-inventory-test-blind-to-asymmetric-addition.md:378`, AC2 in fact reads:
"**AC2 (set-equality assertion is present and load-bearing).** `grep -c 'disagree on their site
set' host/verifygate/floor_raise_inventory_test.go` ≥ 3 (…). Load-bearing proof is mutation arm
N1: with the set-equality neutered … and ARM A1 landed, the test goes GREEN (see Mutation
Drill)." So AC2 carried a REAL mutation arm as its second sentence. The row-59 finding still stands
— the grep half reads 3 under an `if false { … }` wrapper, so the FIRST sentence is a
self-sufficient-LOOKING vacuous discharge, and that is what a reader discharges. The failure mode
is a criterion whose first clause reads as a complete discharge, not a criterion with no mutation
anywhere — which makes the rule SHARPER rather than weaker. The shipped code is correct; the proof
obligation is vacuous. The generalisation to bind: *a criterion that greps for an assertion
measures that somebody typed it; only a mutation measures that anybody runs it.*

**DEFECT B (queue row 50) — the same class in shell.** `host/verifygate/toolchain_pin_gate_test.go`
reads the toolchain deny-list out of a shell script with a column-0-anchored text scan
(`shellAssignmentValues`, V-1). Ratified rule A (D-WORLD-29) relaxes that anchor to accept an
INDENTED assignment. Measured consequence: an assignment placed inside `if false; then … fi` is
then COUNTED by the gate while bash leaves the name UNSET at runtime — so rule A closes one silent
hole by opening a fresh instance of defect A's class. Option B's cure is to stop parsing code at
all: move the three names into a DECLARATIVE, DATA-ONLY fixture holding one `NAME="…"` record per
name, sourced at runtime by the script and scanned by the gates. "Read the code and guess what
runs" becomes "read one declarative record".

Both are the same class: a static text scan standing in for a reachability proof. The single cure
is a rule that a load-bearing claim must be discharged by a mutation, plus a mechanical gate that
rejects the vacuous SHAPE where it can, plus a fixture that makes the shell side data-only by
construction.

---

## Goals

1. **Bind the general rule** (Q1): a "this assertion is load-bearing" criterion must be discharged
   by a MUTATION that makes the test red, never by a static count of the assertion's own text.
   State it in one mechanically-applicable sentence and say exactly which artifacts it binds.
2. **Make the rule fire** (Q2): add a mechanical gate in `host/verifygate` that rejects the vacuous
   SHAPE — an AST reachability pass for Go test assertions, a fixture-shape gate for the shell
   data, and a bounded runtime test that proves run.sh actually EXECUTES the fixture — so the rule
   is enforced, not merely written in prose.
3. **Migrate the toolchain deny-list to a data-only fixture** (Q3): move `KNOWN_BAD` / `KNOWN_GOOD`
   / `PINNED` out of `run.sh`'s code into a declarative fixture whose every non-blank non-comment
   line matches a single anchored assignment regexp, making the column-0 anchor COMPLETE by
   construction and killing rule A's `if false; then` residual outright.
4. **Close row 50 as a consequence** (Q4): state the exact condition under which the original
   defect (a deny-list silently narrowed by an indented assignment) is provably gone, and make it
   an acceptance criterion.
5. **Keep the instrument's logic intact**: rewrite only run.sh's three data lines, never its
   control flow.

---

## Non-Goals

- **Not a general shell parser.** This design does not attempt to parse arbitrary bash. It only
  asserts that ONE fixture file is data-only by construction. Anything that is not a `NAME="…"`
  record in that file is a red, not a parse attempt.
- **No new go.mod dependency.** `go/ast`, `go/parser`, `go/token` are already imported in the
  gate files (V-4). Adding a second direct dependency to `go.mod` (currently exactly one,
  `modernc.org/sqlite v1.54.0`) is a decision, not a detail, and this design does not make it.
- **No rewrite of run.sh's logic — with one bounded exception.** Only its three data lines move to
  the fixture, and the `source` of that fixture is replaced by a bounded fail-loud parser (D4)
  that reads the fixture line by line, validates each record against an inert whitelist, and
  refuses loudly on anything else. This exception was added in response to the round-2 objection
  (gpt5-6-sol): direct sourcing of the fixture is NOT safe, because the fixture's own grammar
  accepts executable shell expansion inside a quoted value (e.g. `KNOWN_BAD="$(touch /tmp/pwned)"`),
  so sourcing it executes the file as code. The probe loop, the platform polarity, the fail-loud
  floors, the `saw_pinned_ok` guard, and the `uname` kernel reads all stay exactly where they are.
- **No change to the row-51 gate's shipped code.** The row-51 gate is correct; its AC2 *proof
  obligation* is vacuous. This design does not touch `floor_raise_inventory_test.go`'s assertions;
  it changes the RULE under which such a claim may be discharged going forward.
- **No attempt to statically detect a runtime-false `if cond { }`.** That residual is declared
  (Q2), not hidden, and is out of scope for a static gate.
- **No change to `scripts/verify_go.sh`.** It is rc=1 at pristine base on a FLEET-OWNED
  driver-drift arm (V-6, queue row 76) and is not an acceptance command here.

---

## High-Impact Decisions

### D1 — The rule (Q1)

**The rule, one sentence:** *A criterion that claims an assertion is load-bearing must be
discharged by a MUTATION that makes the test red, never by a static count of the assertion's own
text.*

**What it binds:** BOTH artifact kinds. (a) Design-doc acceptance criteria — any AC that claims a
gate or assertion is load-bearing must name a mutant and the single assertion that fires when the
mutant lands. (b) Sprint-plan test-plan rows — any row that claims a gate is load-bearing must
carry a mutation arm, not a `grep -c` of the assertion's own message. A `grep` may appear only as
an INSTRUMENT-HEALTH control, explicitly labelled as such (R4).

### D2 — Where it binds so it can actually fire (Q2)

A rule written only in prose in `coding-standards.md` cannot fire. This design chooses **(b): prose
PLUS a mechanical gate**, with two arms and a bounded runtime test:

- **Arm 1 — AST reachability for Go test assertions.** A pass (the `canaryAssertionShapeProblems`
  pattern already in `toolchain_pin_gate_test.go`) parses a named test function and asserts a named
  `t.Errorf` / `t.Fatalf` is REACHABLE at function scope — present as a real call expression in a
  real `if` body, not merely as a string. It already rejects `t.Skip`/`t.Skipf`/`t.SkipNow`, early
  `return`, `goto`, and build constraints (V-12, the A5–A8 arms — verified first-party at HEAD,
  V-16). This is the enforcement of the general rule for Go assertions, and it is DEMONSTRATED, not
  merely cited: the controller landed the M9 mutant and observed the red (V-14).

  **Honest coverage boundary (X1).** Arm 1 adds NO new mechanical enforcement for Go assertions.
  It cites an EXISTING gate that guards exactly ONE canary file
  (`host/store/toolchain_canary_test.go`). So after this design ships, the general rule fires
  mechanically on (a) the shell fixture (Arm 2 + the runtime test, AC4) and (b) that one canary,
  and is PROSE-ONLY for every other Go assertion in the repo. Widening it to a repo-wide AST gate
  is NOT in scope for this design — that would be a much larger row (a repo-wide AST reachability
  pass over every `t.Errorf`/`t.Fatalf` in the tree) and is recorded as a follow-up queue row, not
  funded here. This document does not oversell the coverage: the mechanical net is two files wide.
- **Arm 2 — fixture-shape gate for the shell data.** A new test asserts the toolchain fixture is
  data-only BY CONSTRUCTION: every non-blank non-comment line matches a single anchored assignment
  regexp. This makes the column-0 anchor COMPLETE rather than conventional, and kills rule A's
  `if false; then` residual outright (Q3).

**What each candidate CANNOT see — declared, not hidden:**

- An `if false { … }` wrapper is trivially detectable by Arm 1 (the assertion if-stmt count drops
  to 0) and by Arm 2 (the wrapped line fails the anchored regexp).
- A `source` line placed AFTER an early exit in run.sh is invisible to `bash -n` and to a
  `grep -c` of the source line's presence — both read the text, not the execution order. The
  runtime test (AC4) closes this half of the objection: it executes the prologue through the
  fixture-load point AS WRITTEN and reads the sentinel after it, so a source line that is
  unreachable (wrapped in `if false`, or placed after an `exit`) leaves the sentinel unset and the
  test reds with `run.sh did not execute toolchain_pins.conf` (M10).
- An `if cond { … }` with a RUNTIME-false `cond` is NOT statically detectable by either arm. A
  condition keyed on a commit marker nobody writes, or on an env var that is never set, disables
  the assertion with every byte intact. This residual is real and is declared here: the mutation
  rule catches it only when a controller actually lands the mutant and observes the red. No static
  gate can close it; claiming otherwise would be the vacuity this document exists to kill.

**Candidate evaluated and NOT shipped as a hard gate — the sprint-plan JSON scan.** A scan of
sprint-plan JSON test-plan rows for a `grep -c` used as the sole discharge of a load-bearing claim
was considered. It CAN see a `grep -c` in a test-plan row. It CANNOT see whether that grep is the
*sole* discharge, whether the claim is genuinely load-bearing, or whether the grep is a legitimate
instrument-health control. It would therefore either false-red on legitimate controls or be
trivially evadable by rewording the claim. It is adopted as a PROSE rule in `coding-standards.md`
S6 (D3), not as a mechanical gate.

### D3 — Prose home (Q5)

`design_docs/coding-standards.md` S6 ("Honest gates") is the natural home for the prose rule. It
is declared BINDING ON ALL CODE by CLAUDE.md (V-5). The rule is added as a new S6 sub-clause: a
load-bearing claim must be discharged by a mutation, never by a count of the assertion's own text.

### D4 — The fixture (Q3)

**Path:** `design_docs/verification/w-race-gate-blindspot/toolchain_pins.conf`

**Grammar — DATA ONLY, nothing else:**

```
# Toolchain pin fixture — DATA ONLY. One NAME="value" record per line.
# A record is: NAME="value" optionally followed by whitespace and a # comment.
# Nothing else is legal: no control flow, no substitution, no sourcing, no indentation.
KNOWN_BAD="go1.26.0 go1.26.3 go1.26.4 go1.26.5"
KNOWN_GOOD="go1.26.6 go1.25.6 go1.24.9"
PINNED="go1.26.6"   # trailing comment allowed
```

The anchored record regexp is `^([A-Z_]+)="([^"]*)"(\s+#.*)?$`. Every non-blank non-comment line
must match it. **This grammar is a STATIC SHAPE CHECK, not a safety property**: it is the GATE's
shape check, and it alone does not make the file safe to source — a value like
`KNOWN_BAD="$(touch /tmp/pwned)"` matches the regexp yet is executable. Safety comes from the
bounded fail-loud parser in run.sh (below), which additionally validates each value against an
inert whitelist: toolchain-token characters plus spaces — the character class
`[A-Za-z0-9._+:/ -]` (letters, digits, `.`, `_`, `+`, `:`, `/`, space, hyphen) — and refuses
anything else. The trailing comment on the PINNED line (V-3) is handled by the `(\s+#.*)?` arm —
the grammar accepts it, so the migration does NOT have to drop it.

**The load-bearing property:** because every non-blank non-comment line must match a column-0
anchored regexp, an indented assignment is structurally impossible — it would fail the regexp and
the gate would red. Rule A's `if false; then` residual is killed outright, not narrowed: a line
`  KNOWN_BAD="…"` (indented) or `if false; then` / `fi` is not a valid record and the gate reds.

**How run.sh consumes it:** run.sh does NOT source the fixture. Direct sourcing is unsafe — the
fixture's own grammar accepts executable shell expansion inside a quoted value, so sourcing it
executes the file as code (V-17). Instead run.sh loads the fixture with a bounded fail-loud parser
BEFORE its `cd`, using a path derived from `$0`:

```bash
# Bounded fail-loud parser for the data-only fixture. Accepts ONLY the three
# names KNOWN_BAD, KNOWN_GOOD, PINNED; values must be inert toolchain tokens
# (character class [A-Za-z0-9._+:/ -]) plus spaces. Refuses loudly on anything
# else. This is a parser, not a source: it never evaluates the file as code.
conf="$(dirname "$0")/toolchain_pins.conf"
seen_bad=0; seen_good=0; seen_pinned=0
while IFS= read -r line || [ -n "$line" ]; do
  case "$line" in
    ''|\#*) continue ;;                       # skip blanks and # comments
  esac
  # THE REGEX MUST LIVE IN A VARIABLE. Written inline, the embedded `"` characters
  # are a SYNTAX ERROR on bash 3.2 (macOS system bash, the version the "launchd
  # drivers (bash 3.2)" CI job pins) — measured, V-19.
  RE='^([A-Z_]+)="([^"]*)"([[:space:]]+#.*)?$'
  BADCHARS='[^A-Za-z0-9._+:/ -]'
  if [[ "$line" =~ $RE ]]; then
    name="${BASH_REMATCH[1]}"; value="${BASH_REMATCH[2]}"
    case "$name" in
      KNOWN_BAD)  seen_bad=$((seen_bad+1)) ;;
      KNOWN_GOOD) seen_good=$((seen_good+1)) ;;
      PINNED)     seen_pinned=$((seen_pinned+1)) ;;
      *) echo "toolchain_pins.conf: unknown name '$name' (only KNOWN_BAD, KNOWN_GOOD, PINNED allowed)" >&2; exit 1 ;;
    esac
    # BADCHARS MUST stay in a variable: the inline form `=~ [^A-Za-z0-9._+:/ -]` is a
    # bash 3.2 SYNTAX ERROR (rc=2 `syntax error near \`-]\'`) on THIS rig, including on a
    # good value -- so every M11/M12/M13 arm would read "rejected" while the parser
    # rejected everything. Same class as V-19, one call site over. See V-20.
    if [[ "$value" =~ $BADCHARS ]]; then
      echo "toolchain_pins.conf: value for '$name' contains a disallowed character (only toolchain tokens and spaces allowed)" >&2; exit 1
    fi
    printf -v "$name" '%s' "$value"
  else
    echo "toolchain_pins.conf: malformed line (must be NAME=\"value\" at column 0): $line" >&2; exit 1
  fi
done < "$conf"
# Refuse unless all three names were seen exactly once (this also rejects duplicates).
if [ "$seen_bad" -ne 1 ] || [ "$seen_good" -ne 1 ] || [ "$seen_pinned" -ne 1 ]; then
  echo "toolchain_pins.conf: expected exactly one each of KNOWN_BAD, KNOWN_GOOD, PINNED (got $seen_bad/$seen_good/$seen_pinned)" >&2; exit 1
fi
cd "$(dirname "$0")/repro" || exit 1
```

The parser reads line by line, skips blanks and `#` comments, matches the anchored record regexp,
and REFUSES (exit non-zero, loud message) on an unknown name, a duplicate name, a malformed line,
or a value outside the whitelist; it assigns with `printf -v`; after the loop it refuses unless all
three names were seen. Because the parser never evaluates the file as code, a hostile value like
`KNOWN_BAD="$(touch /tmp/pwned)"` is rejected before it can execute. The runtime reads the SAME
data the gates scan — one source of truth.

**What happens to `shellAssignmentValues`:** it is KEPT (it is a correct column-0 scanner for a
data-only file) but REPOINTED: its four call sites (V-2) change their first argument from run.sh's
lines to the fixture's lines. It is not deleted, because it remains the correct way to read a
data-only fixture; it is not used on run.sh's own non-fixture text anymore. All four call sites in
V-2 are updated to read `toolchain_pins.conf`.

### D5 — Row 50 closes as a consequence (Q4)

Row 50's original defect is a deny-list silently narrowed by an indented assignment. It is provably
gone when the three names live in a data-only fixture whose every non-blank non-comment line
matches the anchored record regexp — because an indented assignment cannot exist in such a file
without the fixture-shape gate redding. That condition is AC2 below. Row 50 closes with this
work, not by a separate sprint (D-WORLD-31).

---

## Solution Design

### The general rule (Q1, D1)

One sentence, mechanically applicable: *a criterion that claims an assertion is load-bearing must
be discharged by a MUTATION that makes the test red, never by a static count of the assertion's
own text.* It binds design-doc acceptance criteria and sprint-plan test-plan rows alike. A grep may
appear only as an instrument-health control, explicitly labelled.

### The fixture (Q3, D4)

`design_docs/verification/w-race-gate-blindspot/toolchain_pins.conf` holds the three names as
data-only records. run.sh parses it with a bounded fail-loud parser (never sources it); the gates
scan it. The column-0 anchor becomes complete by construction.

### The fixture-shape gate (Q2, D2 Arm 2)

A new test in `host/verifygate/toolchain_pin_gate_test.go` (or a sibling file) reads the fixture
and asserts:

1. Every non-blank non-comment line matches `^([A-Z_]+)="([^"]*)"(\s+#.*)?$`.
2. The three required names `KNOWN_BAD`, `KNOWN_GOOD`, `PINNED` each appear exactly once.
3. The fixture is sourced by run.sh (a `source`/`.` of the fixture path is present in run.sh).

### The runtime fixture-execution test (Q2, D2 Arm 2, AC4)

`TestRunShExecutesToolchainPinFixture` is the load-bearing discharge for "run.sh actually loads
the fixture at runtime" — the property `bash -n` and a textual `grep` of the source line cannot
prove. In a `t.TempDir()` it:

1. Copies the run.sh prologue through the fixture-load point (the `source` line and everything
   before it) into a scratch `run.sh`.
2. Copies a fixture holding SENTINEL pin values (e.g. `KNOWN_BAD="sentinel-bad"`,
   `KNOWN_GOOD="sentinel-good"`, `PINNED="sentinel-pinned"`) into the same directory.
3. Appends a command that prints the three loaded variables.
4. Executes the scratch script under a Go `context.WithTimeout` of ≤30s, and asserts the sentinel
   values are observed.

It never runs the repo's own run.sh, never fetches toolchains, and never runs the repro module.
It `t.Fatalf`s LOUDLY if `bash` is unavailable — a silent skip would be the vacuous pass this
document exists to kill. Because the sentinel is read AFTER the prologue executes as written, a
`source` line that is unreachable (wrapped in `if false`, or placed after an early `exit`) leaves
the sentinel unset and the test reds with `run.sh did not execute toolchain_pins.conf` (M10).

### The AST reachability arm (Q2, D2 Arm 1)

The existing `canaryAssertionShapeProblems` pass (V-12) is the enforcement pattern for Go test
assertions. It is retained and, where a new load-bearing assertion is added, the design doc and
sprint plan must name the mutant that makes it red. No new AST code is required for this design
itself; the pattern already exists and is the precedent the rule cites. Its Go-side enforcement
is DEMONSTRATED, not merely cited: the controller landed the M9 mutant and observed the red
(V-14).

### Prose rule (Q5, D3)

Add a sub-clause to `coding-standards.md` S6: a load-bearing claim must be discharged by a
mutation, never by a count of the assertion's own text.

---

## Files to Modify/Create

| # | File | Change |
|---|---|---|
| 1 | `design_docs/verification/w-race-gate-blindspot/toolchain_pins.conf` | **CREATE** — the data-only fixture (D4). |
| 2 | `design_docs/verification/w-race-gate-blindspot/run.sh` | **MODIFY** — delete the three data lines (V-3, lines 24–26); add the bounded fail-loud parser (D4) that loads the fixture before the `cd`. No other logic changes. |
| 3 | `host/verifygate/toolchain_pin_gate_test.go` | **MODIFY** — repoint the four `shellAssignmentValues` call sites (V-2) **and the fifth reader, the known-positive control loop at line 315 (V-22)**, from run.sh to the fixture; add the fixture-shape gate test (AC1) and the runtime fixture-execution test `TestRunShExecutesToolchainPinFixture` (AC4). |
| 4 | `design_docs/coding-standards.md` | **MODIFY** — add the S6 sub-clause (D3). |
| 5 | `design_docs/planned/w-shell-assignment-parser-drops-an-indented-assignment.md` | **SUPERSEDED** — row 50's doc is folded into this one (D-WORLD-31); mark it superseded, do not fund it separately. |

---

## Conflict Surface

This change touches a fixture and a gate that queue rows 43, 44, 48, 50, 51 and 52 all key on.
Every test and every frozen message string that could be disturbed is named here.

**Tests that read run.sh's data lines and MUST be repointed (V-2):**

- `TestMiscompileInstrumentProbesPinnedToolchain` — reads `KNOWN_GOOD`, `KNOWN_BAD`, `PINNED` via
  `shellAssignmentValues` at lines 330–332. After the migration these read the fixture.
- `TestReproModuleFloorStaysBelowKnownBadToolchains` — reads `KNOWN_BAD` via `shellAssignmentValues`
  at line 419. After the migration it reads the fixture.
- **`TestMiscompileInstrumentProbesPinnedToolchain`'s known-positive control loop at line 315 — the
  FIFTH reader, which V-2 did not enumerate because it does not call `shellAssignmentValues`** (V-22).
  It scans run.sh's raw text for the literal strings `KNOWN_BAD=`, `KNOWN_GOOD=` and `PINNED=` and
  `t.Fatalf`s `instrument failure: … does not contain known-positive control %q` when any is missing.
  Deleting run.sh's data lines therefore REDS this test at the MS1 boundary, breaking bisectability,
  unless the loop is repointed at the fixture in the same commit. Repoint it — do NOT delete it: it
  is the instrument-health control that proves the scan can see a positive.

**Tests that read run.sh's STRUCTURE and must NOT be disturbed (their needles stay in run.sh):**

- `TestMiscompileInstrumentStepIsGatedInCI` — pins `uname -s`, `uname -m`, the `steps:` anchor, the
  step name `Measure compiler reproducer (platform-conditional, gated)`, and the absence of
  `continue-on-error` in the miscompile step. None of these move to the fixture; run.sh keeps them.
- `TestRaceControlFloorStaysBelowRootToolchain` — pins `verify_go.sh`'s P1/P2 needles and the
  `racecontrol`/root floor ordering. Untouched.
- `TestGoToolchainPinsAgreeAndMatchJobList` — pins ci.yml pins and the module floor. Untouched.

**The bounded parser changes run.sh's control flow — re-checked.** The `source` line is replaced by
a read-loop parser (D4), so run.sh's structure is no longer byte-identical to today. Re-checking
every test that pins run.sh's structure: `TestMiscompileInstrumentStepIsGatedInCI` pins `uname -s`,
`uname -m`, the `steps:` anchor, the step name, and the absence of `continue-on-error` — none of
these are in the replaced source line, so it is NOT disturbed. `TestRaceControlFloorStaysBelowRootToolchain`
pins `verify_go.sh`'s P1/P2 needles and the `racecontrol`/root floor ordering — untouched.
`TestGoToolchainPinsAgreeAndMatchJobList` pins ci.yml pins and the module floor — untouched. The
`saw_pinned_ok` guard and the `INSTRUMENT FAILURE: the PINNED toolchain` message stay in run.sh and
are not disturbed. The parser is added BEFORE the `cd`, in the same position the `source` occupied,
so no downstream needle moves.

**Frozen message strings that must NOT be renamed:**

- `"INSTRUMENT FAILURE: the PINNED toolchain"` — the fail-loud guard message in run.sh, pinned by
  `TestMiscompileInstrumentProbesPinnedToolchain`. It stays in run.sh.
- `"saw_pinned_ok"` — the OK-flag site, pinned at ≥3 occurrences in run.sh. It stays in run.sh.
- `"disagree on their site set"` — the row-51 gate's message in `floor_raise_inventory_test.go`
  (lines 165/170/174, verified first-party at HEAD, V-15). This design does NOT rename it; it
  changes the RULE under which a claim about it may be discharged. Renaming it would break the
  row-51 gate's own (correct) code.
- `"KNOWN_BAD"`, `"KNOWN_GOOD"`, `"PINNED"` — the three names. They move to the fixture but keep
  their exact spellings; the gate tests and run.sh's `saw_pinned_ok` guard reference them by name.

**`floor_raise_inventory_test.go` (row 51, V-1):** its assertions are correct and are NOT touched.
The defect is that row 51's AC2 discharged a load-bearing claim with `grep -c 'disagree on their
site set' ≥ 3`. This design does not edit that file; it establishes the rule so future claims are
discharged by mutation.

**`coding-standards.md` S8 (row 43, V-5):** the floor-raise coupling inventory table and its
`floor_raise_inventory_test.go` binding are untouched. Adding an S6 sub-clause does not move any
S8 row; the S8 table's six rows and the `## S8` heading must remain byte-identical (the inventory
gate pins them).

**Identifiers that must NOT be renamed:** `KNOWN_BAD`, `KNOWN_GOOD`, `PINNED`, `saw_pinned_ok`,
`shellAssignmentValues`, `disagree on their site set`, `INSTRUMENT FAILURE: the PINNED toolchain`,
`Measure compiler reproducer (platform-conditional, gated)`.

---

## Acceptance Criteria

Each AC names a command a reviewer can run and its expected rc/output, baselined against a
PRISTINE tree. None is red at base (V-6, V-7). None uses `grep -c <assertion text> == N` as the
discharge of a load-bearing claim.

- **AC1 — the fixture exists and is data-only by construction.**
  `go test ./host/verifygate/ -run 'TestToolchainPinFixtureIsDataOnly' -v 2>&1 | tee /dev/stderr
  | grep -q -- '--- PASS: TestToolchainPinFixtureIsDataOnly'` → rc=0, AND the same run must NOT
  print `no tests to run`. **The bare `-run` form is VACUOUS at base and was wrong here (V-21):**
  `go test -run '<a name that does not exist>'` exits **rc=0** printing `[no tests to run]`, so the
  original discharge passed before the test was written -- this document's own thesis, committed by
  the document. Baseline under the hardened form: **rc=1** at base, rc=0 once the test lands.
- **AC2 — row 50's defect is provably gone (Q4).** The three names live in `toolchain_pins.conf`
  and the fixture-shape gate `TestToolchainPinFixtureIsDataOnly` passes (AC1). The mutation table
  M1/M2/M8 discharges the load-bearing claim that an indented or control-flow-wrapped assignment
  cannot survive in the fixture. No grep-based discharge is used.
- **AC3 — the gate tests read the fixture, not run.sh's code.**
  `go test ./host/verifygate/ -run 'TestMiscompileInstrumentProbesPinnedToolchain|TestReproModuleFloorStaysBelowKnownBadToolchains'`
  → rc=0. (Baseline: both pass at base; they must still pass after repointing to the fixture.)
- **AC4 — run.sh actually EXECUTES the fixture at runtime (bounded).**
  `go test ./host/verifygate/ -run 'TestRunShExecutesToolchainPinFixture' -v 2>&1 | tee /dev/stderr
  | grep -q -- '--- PASS: TestRunShExecutesToolchainPinFixture'` → rc=0, AND the same run must NOT
  print `no tests to run` (V-21, exactly as AC1 -- the bare `-run` form is rc=0 at base).
  Baseline under the hardened form: **rc=1** at base, rc=0 once the test lands. The test copies the run.sh prologue through the fixture-load point and a sentinel
  fixture into `t.TempDir()`, appends a command that prints the loaded variables, and executes it
  under a Go `context.WithTimeout` of ≤30s. It asserts the sentinel values are observed, and
  `t.Fatalf`s LOUDLY if `bash` is unavailable (a silent skip is the vacuous pass this document
  exists to kill). It never runs the repo's own run.sh, never fetches toolchains, and never runs
  the repro module. The Go test MUST intercept `exec.Command` errors — a bash syntax error from a
  truncated `fi`, or an unbound variable from `set -u` if the source is skipped — and
  unconditionally emit the required `run.sh did not execute toolchain_pins.conf` substring in
  `t.Fatalf`, so M10's exact assertion matches regardless of the specific bash failure mode.
  `bash -n` is retained only as a syntax/instrument-health control, NOT as the load-bearing
  discharge: it is genuinely instrument health (it checks the scratch script parses), and the
  load-bearing discharge is the runtime execution itself, so it is not a vacuous grep standing in
  for a mutation.
- **AC5 — the whole Go suite is green.**
  `go build ./... && go test ./...` → rc=0. (Baseline: green at base, V-6/V-7.)
- **AC6 — the compile fence for _test.go files holds.**
  `go vet ./host/...` → rc=0. (V-7: `go build ./...` is NOT a compile fence for a `_test.go` file;
  `go vet` is the correct fence. Baseline: green at base, V-7.)
- **AC7 — the prose rule is present.**
  `grep -c 'load-bearing' design_docs/coding-standards.md` → ≥ 1 in the S6 sub-clause. (This is an
  instrument-health control on the prose rule's presence, not a load-bearing discharge: it only
  checks that the prose sentence exists at all, and no mutation would catch "the prose sentence
  exists" — the load-bearing claim, that the rule fires, is discharged by the mutation table, not
  by this grep. Baseline: the S6 sub-clause does not exist at base, so this AC is green only after
  the prose rule lands.)

---

## Testing Strategy (Non-Vacuity)

Named mutation table — one row per mutation, each naming the mutant, the venue, and THE SINGLE
ASSERTION THAT MUST FIRE (its exact message substring). Every row is a mutation a person can land
and revert. No row uses `grep -c <assertion text> == N` as the discharge of a load-bearing claim.

| # | Mutant | Venue | The single assertion that must fire (exact substring) |
|---|---|---|---|
| M1 | Indent the `KNOWN_BAD` line in the fixture (`  KNOWN_BAD="…"`) | fixture-shape gate | `does not match the anchored assignment grammar` |
| M2 | Wrap the `KNOWN_BAD` line in `if false; then … fi` in the fixture | fixture-shape gate | `does not match the anchored assignment grammar` |
| M3 | Delete the `KNOWN_BAD` line from the fixture | `TestMiscompileInstrumentProbesPinnedToolchain` | `KNOWN_BAD assignment count=0, want 1` |
| M4 | Add a second `KNOWN_BAD` line to the fixture | `TestMiscompileInstrumentProbesPinnedToolchain` | `KNOWN_BAD assignment count=2, want 1` |
| M5 | Change `KNOWN_BAD` to include the pinned floor | `TestMiscompileInstrumentProbesPinnedToolchain` | `incorrectly labels the pinned toolchain` |
| M6 | Change `PINNED` to a selection mode (`go1.26.6+auto`) | `requireToolchainNamePin` | `is a toolchain-selection mode, not a pin` |
| M7 | Change `PINNED` to a different floor | `TestMiscompileInstrumentProbesPinnedToolchain` | `want go.mod floor` (V-23: the row previously quoted `PINNED=…, want go.mod floor`, which can never appear literally — the format string is `PINNED=%q, want go.mod floor %q`, so the `…` is not a real substring. Match the invariant tail.) |
| M8 | Add a non-assignment line to the fixture (`echo hello`) | fixture-shape gate | `does not match the anchored assignment grammar` |
| M9 | Wrap the canary assertion `if rows[0].field != "stateRoot" { t.Fatalf(…) }` in `if false { … }` | `TestCanaryDeclaresPositiveArmOnly` (via `canaryAssertionShapeProblems`) | `assertion if-stmt count=0, want exactly 1` |
| M10 | Wrap the fixture `source` line in `if false; then … fi` in run.sh | `TestRunShExecutesToolchainPinFixture` (runtime) | `run.sh did not execute toolchain_pins.conf` |
| M11 | Command substitution in a value: `KNOWN_BAD="$(touch sentinel)"` in the fixture | bounded fail-loud parser in run.sh (AC4) | `value for 'KNOWN_BAD' contains a disallowed character` — and the sentinel file is ABSENT afterwards (the test observes "did not execute" by asserting the sentinel file does not exist) |
| M12 | Backtick expansion in a value: `` KNOWN_BAD="`touch sentinel`" `` in the fixture | bounded fail-loud parser in run.sh (AC4) | `value for 'KNOWN_BAD' contains a disallowed character` — and the sentinel file is ABSENT afterwards |
| M13 | An arbitrary fourth name: `PATH="/nonsense"` in the fixture | bounded fail-loud parser in run.sh (AC4) | `unknown name 'PATH' (only KNOWN_BAD, KNOWN_GOOD, PINNED allowed)` — and `PATH` is NOT clobbered (the test asserts the shell's `PATH` is unchanged) |

M9 is the demonstration of the general rule (D1): a mutant that wraps an assertion in `if false`
must make the gate red. It is the same class as the row-51 defect, applied to the canary that
already has the AST reachability arm. **M9 is DISCHARGED, not proposed**: the controller landed it
in the design worktree and measured the red first-party (V-14) — it is the document's own
self-satisfaction proof. M1/M2/M8 are the fixture-shape gate's own red arms; M3–M7 are the existing
gate tests' red arms, now reading the fixture; M10 is the runtime test's red arm (AC4); M11–M13 are
the bounded parser's red arms (AC4) — each is REJECTED with non-zero status AND must not execute or
alter state, and the test observes "did not execute" by asserting the sentinel file is absent
afterwards.

**Declared residual (not hidden):** an `if cond { }` with a RUNTIME-false `cond` is not statically
detectable by any arm. The mutation rule catches it only when a controller lands the mutant and
observes the red. No static gate closes it.

---

## Timeline / Milestones

At most 3 milestones, each ≤1 day, each independently landable with its own commit.

- **M1 (≤1d) — the fixture + repoint.** Create `toolchain_pins.conf`; delete the three data lines
  from run.sh and add the `source`; repoint the four `shellAssignmentValues` call sites to the
  fixture; add the runtime fixture-execution test `TestRunShExecutesToolchainPinFixture` (AC4).
  Commit. AC3, AC4, AC5, AC6 green.
- **M2 (≤1d) — the fixture-shape gate.** Add `TestToolchainPinFixtureIsDataOnly` (AC1, AC2) with
  the M1/M2/M8 red arms. Commit. AC1, AC2 green.
- **M3 (≤1d) — the prose rule + row-50 closure.** Add the S6 sub-clause to `coding-standards.md`;
  mark row 50's doc superseded; record row 50's closure condition (AC2). Commit. AC7 green.

---

## Risks & Mitigations

- **Risk: the fixture-shape gate false-reds on a legitimate re-indent.** Mitigation: the fixture
  is data-only by construction; there is no legitimate reason to indent a record. A re-indent is a
  red by design, and the message names the offending line.
- **Risk: run.sh's `source` breaks the `set -uo pipefail` strictness.** Mitigation: the fixture is
  data-only and sets exactly three names; sourcing it under `set -u` is safe because it references
  no unset variables. AC4's runtime test executes the prologue under the script's own `set -uo
  pipefail` and observes the sentinel; `bash -n` remains only as a syntax/instrument-health
  control.
- **Risk: a runtime-false `if cond { }` residual is mistaken for a closed hole.** Mitigation: the
  residual is DECLARED in D2 and the Testing Strategy, not hidden. The mutation rule is the only
  thing that catches it, and that is stated.
- **Risk: someone "fixes" this document by adding a `grep -c <assertion text> == N` discharge.**
  Mitigation: the header and R4 forbid it; shipping one would be self-refuting.
- **Risk: `go build ./...` is used as the compile fence for a `_test.go` mutant.** Mitigation: AC6
  uses `go vet ./host/...` (V-7, queue row 65).
- **Risk: `verify_go.sh` is used as an acceptance command.** Mitigation: it is rc=1 at pristine
  base on a FLEET-OWNED driver-drift arm (V-6, queue row 76); AC5 uses the two-leg substitute
  `go build ./... && go test ./...`.

---

## Verification Log

Every claim about the current codebase carries its command and observed output. Empty or negative
results are paired with a known-positive control in the same command. All measured at HEAD
`79d80d9` on 2026-09-04.

| # | Claim | Command | Observed output |
|---|---|---|---|
| V-1 | `shellAssignmentValues` is column-0 anchored, value-terminated at the first `"` | `sed -n '285,297p' host/verifygate/toolchain_pin_gate_test.go` | `func shellAssignmentValues(lines []string, name string) []string { prefix := name + "=\""; … if strings.HasPrefix(line, prefix) { … if end := strings.IndexByte(rest, '"'); end >= 0 { … } } }` — `HasPrefix` on the raw line, first-`"` termination. |
| V-2 | Exactly 4 call sites of `shellAssignmentValues`, all in that one file | `grep -rn "shellAssignmentValues" --exclude-dir=.git . \| grep -v design_docs` | Definition at `:285`; call sites at `:330` (KNOWN_GOOD), `:331` (KNOWN_BAD), `:332` (PINNED) in `TestMiscompileInstrumentProbesPinnedToolchain`, and `:419` (KNOWN_BAD) in `TestReproModuleFloorStaysBelowKnownBadToolchains`. |
| V-3 | run.sh data lines at column 0; trailing comment on PINNED; `set -uo pipefail`; `cd "$(dirname "$0")/repro"`; executable; 6809 bytes; 18 control-flow openers; 64 indented non-comment lines | `sed -n '24,26p' design_docs/verification/w-race-gate-blindspot/run.sh`; `grep -nE '^\s*(if\|for\|while\|case\|function\|probe\()' …/run.sh`; `grep -cE '^\s+[^#]' …/run.sh`; `wc -c …/run.sh`; `ls -l …/run.sh` | Lines 24–26 at column 0: `KNOWN_BAD="go1.26.0 go1.26.3 go1.26.4 go1.26.5"`, `KNOWN_GOOD="go1.26.6 go1.25.6 go1.24.9"`, `PINNED="go1.26.6"   # trailing comment present`; `set -uo pipefail` and `cd "$(dirname "$0")/repro"` present; 18 openers; 64 indented non-comment lines; 6809 bytes; `-rwxr-xr-x` (executable). *(Re-measured at HEAD; the mission's iter-146 figures of 17/65 have since drifted to 18/64.)* |
| V-4 | `go/ast` + `go/parser` already imported in the gate files; AST gates exist repo-wide; go.mod has exactly one direct dependency | `sed -n '1,12p' host/verifygate/toolchain_pin_gate_test.go`; `sed -n '1,10p' host/verifygate/subprocess_sink_gate_test.go`; `grep -rln "go/ast" host/broker host/evidence`; `cat go.mod` | toolchain_pin_gate_test.go imports `go/ast`, `go/build/constraint`, `go/parser`, `go/token`, `go/version`; subprocess_sink_gate_test.go imports `go/ast`, `go/parser`, `go/token`, `go/types`; AST gates in `host/broker/invoke_boundary_test.go`, `host/broker/handlers_parallel_guard_test.go`, `host/broker/registry_publish_test.go`, `host/evidence/authority_test.go`; go.mod `require modernc.org/sqlite v1.54.0` (one direct dep; the rest are indirect). |
| V-5 | `coding-standards.md` is 122 lines, sections S1..S8; S6 is "Honest gates"; S8 is the floor-raise inventory | `wc -l design_docs/coding-standards.md`; `grep -nE '^## S[0-9]' design_docs/coding-standards.md` | 122 lines; `## S6 — Honest gates` and `## S8 — The floor-raise coupling inventory (added 2026-08-27, row 43)` present. |
| V-6 | 291 tracked files; `verify_go.sh` rc=1 at pristine base on the driver-drift arm; `go build ./...` rc=0 | `git ls-files \| wc -l`; `./scripts/verify_go.sh; echo rc=$?`; `go build ./...; echo rc=$?` | 291; `verify_go.sh` rc=1 (`✗ AILANG_BIN is unset — host/replay tests would t.Skip() silently…` — the driver-drift arm short-circuits its Go legs, queue row 76); `go build ./...` rc=0. |
| V-7 | `go build ./...` is not a compile fence for a `_test.go` file | `go build ./...; echo rc=$?` | rc=0 (a type error in a `_test.go` file would still exit 0; the correct fence is `go vet ./host/...` or `go test -run '^$' ./host/...` — queue row 65). |
| V-8 | Row 51's AC2 discharged a load-bearing claim with `grep -c 'disagree on their site set' ≥ 3` | `grep -n 'disagree on their site set' design_docs/planned/w-inventory-test-blind-to-asymmetric-addition.md` | `:378` — `**AC2 (set-equality assertion is present and load-bearing).** grep -c 'disagree on their site set' host/verifygate/floor_raise_inventory_test.go ≥ 3 …` — the vacuous discharge this design kills. |
| V-9 | The canary assertion shape is `if rows[0].field != "stateRoot" { t.Fatalf(…) }` | `grep -n 'stateRoot\|if rows\|t.Fatalf\|POSITIVE ARM ONLY' host/store/toolchain_canary_test.go` | `:10 // POSITIVE ARM ONLY…`; `:52 if rows[0].field != "stateRoot" {`; `:53 t.Fatalf("Field=%q want %q", rows[0].field, "stateRoot")` — the assertion the AST reachability arm (M9) guards. |
| V-10 | Row 50's doc exists and ratified option A under D-WORLD-29 | `head -60 design_docs/planned/w-shell-assignment-parser-drops-an-indented-assignment.md` | Status `REVISED under human ruling D-WORLD-29 (attended 2026-09-01) — direction ratified as option A`; the doc is superseded by this one (D-WORLD-31). |
| V-11 | `floor_raise_inventory_test.go` pins the S8 table and the `## S8` heading | `cat host/verifygate/floor_raise_inventory_test.go` | `TestFloorRaiseInventoryNamesEveryCoupledFile` asserts the `## S8` heading, the six S8 rows, and the six inventory-block rows; the S8 table and heading must remain byte-identical. |
| V-12 | `canaryAssertionShapeProblems` EXISTS and is an AST pass | controller first-party, iteration 154: `grep -rn "canaryAssertionShapeProblems" --exclude-dir=.git . \| grep -v design_docs` | `host/verifygate/toolchain_pin_gate_test.go:441` (doc comment), `:445` (func decl), `:632` (the single call site). Doc comment, verbatim from `:441-444`: "canaryAssertionShapeProblems parses the canary test source and reports any deviation from the required assertion shape: TestToolchainCanary must exist exactly once and contain exactly one top-level `if` whose condition is rows[0].field != \"stateRoot\". That if-body must contain exactly one direct t.Fatalf expression statement, rather than merely a descendant call." Mechanism read at `:445-500`: `parser.ParseFile` → find FuncDecls named `TestToolchainCanary` → iterate `funcs[0].Body.List` for TOP-LEVEL `*ast.IfStmt` whose Cond is a BinaryExpr NEQ over `rows[0].field` and the literal `"stateRoot"` → require exactly 1 → then require exactly one DIRECT `t.Fatalf` ExprStmt in that if-body. It also counts `t.Skip/t.Skipf/t.SkipNow` and requires 0. |
| V-13 | `TestCanaryDeclaresPositiveArmOnly` exists and targets `host/store/toolchain_canary_test.go` | controller first-party, iteration 154: `sed -n '625,644p' host/verifygate/toolchain_pin_gate_test.go` | Verbatim: `func TestCanaryDeclaresPositiveArmOnly(t *testing.T) { canaryPath := filepath.Join(repoRoot, "host", "store", "toolchain_canary_test.go"); raw, err := os.ReadFile(canaryPath); if err != nil { t.Fatal(err) }; src := string(raw); if problems := canaryAssertionShapeProblems(src); len(problems) > 0 { ... } ... also asserts count("GOTOOLCHAIN")==0 and the "POSITIVE ARM ONLY" marker ... }`. So the pass IS invoked against `host/store/toolchain_canary_test.go` by a named test. |
| V-14 | M9 is DISCHARGED: the controller landed the mutant and measured the red | controller first-party, iteration 154, in the design worktree | Pristine control BEFORE: `go test ./host/verifygate/ -run 'TestCanaryDeclaresPositiveArmOnly' -count=1` → rc=0, "ok ... 0.355s". Mutant: the block `if rows[0].field != "stateRoot" { t.Fatalf("Field=%q want %q", ...) }` in `host/store/toolchain_canary_test.go` wrapped verbatim in `if false { ... }`, landed by sha256 `a23cfa79419ae691` → `6ddd8ca5209a3d37` (first 16 hex of `shasum -a 256`). Compile fence read BEFORE any test result: `go vet ./host/store/ ./host/verifygate/` → rc=0 (go vet, NOT go build — queue row 65: `go build ./...` exits 0 on a type error in a `_test.go`). Result: `go test ...` → rc=1, `--- FAIL` count = 1, and the message that fired contains the exact substring `assertion if-stmt count=0, want exactly 1`. Restored: sha256 back to `a23cfa79419ae691` (byte-identical), test rc=0 again, worktree porcelain shows only the design doc. |
| V-15 | The row-51 gate's `"disagree on their site set"` message lives at lines 165/170/174 of `floor_raise_inventory_test.go` (NOT 306/311/315 as an earlier draft claimed) | designer first-party, iteration 154: `grep -n "disagree on their site set" host/verifygate/floor_raise_inventory_test.go` | `:165`, `:170`, `:174` — the three `t.Errorf` sites. This corrects the citation-hygiene defect (objection 2): the earlier `(V-1, lines 306/311/315)` pointer named a row that did not carry the claim and line numbers that were wrong. |
| V-16 | The A5–A8 arms of `canaryAssertionShapeProblems` reject `t.Skip`/`t.Skipf`/`t.SkipNow`, early `return`, `goto`, and build constraints | designer first-party, iteration 154: `grep -n "ReturnStmt\|BranchStmt\|IsGoBuild\|IsPlusBuild\|SkipNow" host/verifygate/toolchain_pin_gate_test.go` | `:502` (Skip/Skipf/SkipNow selector set), `:522` (`*ast.ReturnStmt`), `:545` (`*ast.BranchStmt` with `token.GOTO`), `:581` (`constraint.IsGoBuild`/`IsPlusBuild`). This verifies the D2 clause the controller flagged as unverified; it is now carried with a command, not as an unverified assertion. |
| V-17 | The fixture's own grammar accepts executable shell expansion inside a quoted value, so sourcing it executes the file as code — the "data-only by construction" claim is TRUE of the STATIC SCAN and FALSE of the RUNTIME CONSUMPTION | controller first-party, iteration 154: probe file `f.conf` with three lines — `KNOWN_BAD="$(touch /tmp/fixprobe_iter154/PWNED)"`, `KNOWN_GOOD="go1.26.6"`, `PATH="/nonsense"`; command 1 `grep -cE '^[A-Z_]+="[^"]*"([[:space:]]+#.*)?$' f.conf`; command 2 `( set -uo pipefail; . ./f.conf; echo "KNOWN_BAD=[$KNOWN_BAD] PATH_now=[$PATH]" )`; negative control `echo hello` → same grep → 0 | command 1 observed `3` — ALL THREE hostile lines match the doc's own grammar; command 2 observed `KNOWN_BAD=[] PATH_now=[/nonsense]` and the file `/tmp/fixprobe_iter154/PWNED` WAS CREATED — `source` executed the command substitution and an arbitrary fourth name (PATH) was clobbered; negative control observed `0`, so the grammar is not matching everything. Conclusion: the gate reads the file as data; bash reads it as code. |
| V-18 | Row 51's AC2 carried a REAL mutation arm as its second sentence — the document's earlier claim that it discharged the claim WITH the grep overstates the defect | `sed -n '378p' design_docs/planned/w-inventory-test-blind-to-asymmetric-addition.md` | `**AC2 (set-equality assertion is present and load-bearing).** grep -c 'disagree on their site set' host/verifygate/floor_raise_inventory_test.go ≥ 3 (…). Load-bearing proof is mutation arm N1: with the set-equality neutered … and ARM A1 landed, the test goes GREEN (see Mutation Drill).` — AC2 carried a real mutation arm as its second sentence. The row-59 finding still stands: the grep half reads 3 under an `if false { … }` wrapper, so the FIRST sentence is a self-sufficient-LOOKING vacuous discharge, and that is what a reader discharges. This makes the rule SHARPER rather than weaker: the failure mode is a criterion whose first clause reads as a complete discharge, not a criterion with no mutation anywhere. |
| V-19 | The parser snippet as first drafted is a **bash 3.2 SYNTAX ERROR**, and it fails in the direction that fakes a pass: all three reject-arms M11/M12/M13 would read "rejected" while the GOOD fixture also fails | controller first-party, iteration 154, on `GNU bash, version 3.2.57(1)-release (arm64-apple-darwin25)` — the version the `launchd drivers (bash 3.2)` CI job pins. Ran the drafted parser against four fixtures: a good one, M11 (`KNOWN_BAD="$(touch …/PWNED)"`), M13 (`PATH="/nonsense"`), M1 (indented `KNOWN_BAD`) | **Inline form** (`if [[ "$line" =~ ^([A-Z_]+)="([^"]*)"…$ ]]`): all four rc=2 with `syntax error in conditional expression: unexpected token ')'` — including the GOOD fixture, so the reject-arms would have been vacuously green. **Regex-in-a-variable form** (`RE='…'; if [[ "$line" =~ $RE ]]`): good rc=0 parsing all three names with the trailing comment handled; M11 rc=1 `value for 'KNOWN_BAD' contains a disallowed character` with the `PWNED` sentinel **NOT created** (the value was never executed); M13 rc=1 `unknown name 'PATH'`; M1 rc=1 `malformed line`. The snippet above now carries the variable form and a comment saying why. **Controller factual correction applied directly** — it overrides no objection and resolves nothing that was contested. |
| V-20 | **V-19 fixed the RECORD regex and left the WHITELIST regex inline, so the snippet was STILL a bash 3.2 syntax error — instance 2 of V-19's own class, one call site over** | controller first-party, iteration 155, `/bin/bash --version` = `GNU bash, version 3.2.57(1)-release (arm64-apple-darwin25)`. Two scratch scripts differing only in the regex's placement, each with the GOOD value `value="go1.26.6"`: inline `if [[ "$value" =~ [^A-Za-z0-9._+:/ -] ]]` vs `BADCHARS='[^A-Za-z0-9._+:/ -]'; if [[ "$value" =~ $BADCHARS ]]`. Both `bash -n` and executed. Surfaced by the sprint-planner and reproduced by the controller before acting (rule 3f). | **Inline**: `bash -n` **rc=2**, `syntax error in conditional expression` / ``syntax error near `-]' `` — and identically rc=2 when EXECUTED, on the GOOD value. So M11/M12/M13 would each have read "rejected" while the parser rejected everything, including the fixture it is meant to accept. **Variable**: `bash -n` **rc=0**, execution rc=0 printing `ACCEPT`. Positive control: a knowingly-broken script (`if [ x`) reports `syntax error: unexpected end of file`, so `bash -n` discriminates. Doc corrected to the variable form with `BADCHARS` declared beside `RE`. **Guard the helper, miss the call site** — aimed at the fix for that very shape. |
| V-21 | **AC1 and AC4 were discharged by a command that is rc=0 BEFORE the test exists — the document committing its own thesis for the third time** | controller first-party, iteration 155, on a pristine tree at `32369dc`. Paired run: `go test ./host/verifygate/ -run 'TestToolchainPinFixtureIsDataOnly'` (the not-yet-written test named by AC1) and `go test ./host/verifygate/ -run 'TestMiscompileInstrumentProbesPinnedToolchain'` (a test that exists, as the known-positive control). Surfaced by the sprint-planner, reproduced by the controller before acting. | Nonexistent test: **rc=0**, `ok … [no tests to run]`. Existing test: **rc=0**, `ok`. The two are indistinguishable by exit code, so AC1's and AC4's stated baseline (*"green only after the fixture lands"*) was **false** — both were green at base, vacuously. This is precisely the `grep -c`-as-discharge defect this document exists to kill, wearing `go test -run`'s clothes. Both ACs hardened to require a `--- PASS: <TestName>` line and to refuse on `no tests to run`; under the hardened form the baseline is genuinely rc=1. |
| V-22 | **A FIFTH reader of run.sh's data lines that V-2 did not enumerate, because it does not call `shellAssignmentValues`** | controller first-party, iteration 155: `grep -n "KNOWN_BAD\|KNOWN_GOOD\|PINNED" host/verifygate/toolchain_pin_gate_test.go \| grep -v shellAssignmentValues`, plus `sed -n '308,320p'` for context, plus a negative control on an absent name (`grep -c "KNOWN_UGLY"`). Surfaced by the sprint-planner, reproduced by the controller. | `toolchain_pin_gate_test.go:315`: `for _, control := range []string{"KNOWN_BAD=", "KNOWN_GOOD=", "PINNED="} { if !strings.Contains(src, control) { t.Fatalf("instrument failure: %s does not contain known-positive control %q", …) } }` — a raw-text scan of run.sh, invisible to a `shellAssignmentValues` grep. Negative control `KNOWN_UGLY` → **0**, so the grep was not matching everything. Consequence: deleting run.sh's data lines REDS `TestMiscompileInstrumentProbesPinnedToolchain` at the MS1 boundary and breaks bisectability. Conflict Surface and the Files table corrected; the loop is REPOINTED at the fixture, never deleted — it is the instrument-health control. **V-2's enumeration was anchored to a function name rather than to the fact it claimed**, which is the same shape as its own subject. |
| V-23 | M7's quoted assertion substring cannot appear literally | controller first-party, iteration 155: read the format string at `toolchain_pin_gate_test.go:365`. | The `t.Errorf` is `PINNED=%q, want go.mod floor %q` — so `PINNED=…, want go.mod floor` is not a substring of any real output (the `…` is prose, not a wildcard). Mutation row M7 corrected to match the invariant tail `want go.mod floor`. The other 12 rows' substrings were confirmed present. |

---

## Quorum Verification Log

Round 1 of this document's review quorum: **3/3 present** (absent_reviewers `[]`), all three
REJECT, metered **$0.08615**. The three objections, one line each, and what this revision did
about each:

1. **gpt5-6-sol (REJECT)** — the design recreates the vacuous-proof pattern it kills: runtime
   consumption of `toolchain_pins.conf` was validated only by `bash -n` and a textual source-line
   check. **This revision:** replaced AC4/M10 with the bounded runtime test
   `TestRunShExecutesToolchainPinFixture` (Go `context.WithTimeout` ≤30s, runs `bash` on a COPY in
   `t.TempDir()`, never the repo's own run.sh, fails LOUDLY if `bash` is unavailable), and
   addressed the "after an early exit" half explicitly (the runtime test reads the sentinel after
   the prologue executes as written, so an unreachable source line reds).
2. **gemini-3-1-pro (REJECT)** — D2 cited `canaryAssertionShapeProblems` as an existing AST pass
   via a V-1 pointer that carried no such claim. **This revision:** swept every `(V-n)` pointer,
   added V-12 (existence + mechanism), V-13 (the named test), V-16 (the A5–A8 arms, verified with
   a command), and V-15 (corrected the wrong `floor_raise_inventory_test.go` line numbers).
3. **oc-glm-5-2 (REJECT)** — M9 was asserted, not discharged. **This revision:** added V-14 (the
   controller landed M9 and measured the red first-party) and restated M9 as DISCHARGED, the
   document's own self-satisfaction proof; D2 now says the Go-side arm is DEMONSTRATED, not merely
   cited.

This section is how the next reader knows the objections were answered rather than ignored.

---

## Quorum Verification Log — Round 2 (CONTROLLER CARVE-OUT)

Round 2 of this document's review quorum ran at FULL STRENGTH: **3/3 present** (absent_reviewers
`[]`), metered **$0.10906** (cumulative **$0.19522** across both rounds). `gemini-3-1-pro` **PASS**
(its first pass — the round-1 revision worked), `gpt5-6-sol` **REJECT**, `oc-glm-5-2` **REJECT**.
This second revision was applied under the mission-control **CONTROLLER CARVE-OUT**: both remaining
objections carry a concrete reviewer-authored fix and neither reverses the design's direction, so
the reviewers' verbatim fixes were applied with no third quorum round. The applied set is the four
fixes below; nothing was invented and nothing was overridden.

1. **gpt5-6-sol (REJECT)** — the central construction has a hole: the fixture's own grammar accepts
   executable shell expansion inside a quoted value (e.g. `KNOWN_BAD="$(touch /tmp/pwned)"`), so
   sourcing the fixture executes it as code; the controller confirmed this first-party (V-17). **This
   pass:** replaced direct sourcing with a bounded fail-loud parser in run.sh (D4), tightened the
   value whitelist to inert toolchain tokens, revised the "no control-flow change" non-goal to
   permit the parser, added mutation rows M11–M13, and re-checked the Conflict Surface for the
   parser's control-flow change.
2. **oc-glm-5-2 (REJECT)** — AC2 ships the defect it exists to kill: a `grep -cE` of the fixture's
   own shape as the discharge for "row 50's defect is provably gone". **This pass:** deleted the
   grep row from AC2 and restated it per the reviewer's verbatim prose, and swept the whole AC list
   for the same shape (AC4's `bash -n` and AC7's `grep -c 'load-bearing'` are genuinely instrument
   health, kept with an explicit why-not-a-discharge clause).
3. **gemini-3-1-pro (PASS)** — volunteered a concrete refinement to AC4: a truncated `fi` from a
   multi-line mutant would cause a bash syntax error rather than a clean execution. **This pass:**
   took it — AC4 now states the Go test must intercept `exec.Command` errors and unconditionally
   emit the `run.sh did not execute toolchain_pins.conf` substring in `t.Fatalf`.
4. **Controller correction (accuracy, measured first-party)** — the Problem Statement overstated
   the row-51 defect: AC2 carried a real mutation arm as its second sentence. **This pass:**
   corrected the Problem Statement, added V-18 with the quoted AC2 text, and noted the rule is
   SHARPER rather than weaker.

This section is how the next reader knows the objections were answered rather than ignored, and
that the carve-out applied the reviewers' own words.

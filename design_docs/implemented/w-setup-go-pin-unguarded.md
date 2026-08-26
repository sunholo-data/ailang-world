# w-setup-go-pin-unguarded — the two toolchain-pin kinds can drift apart green, and the go1.26 instrument has never probed the live pin

**Status**: Planned
**Date**: 2026-08-26
**Queue item**: 41, `w-setup-go-pin-unguarded` (clause-2, drill-surfaced, evaluator-corroborated)
**Estimated**: ≤0.2 day (revision round 1, objection 2 option A: the one-line production edit
became a ~10-line pinned-probe guard block in the *same single file*; one new test file in an
existing package, ~170 LOC; **no `ci.yml` edit, no `scripts/` edit, no other package
touched**)
**Designer**: `pi:ollama/kimi-k3:cloud` (rotation, iteration 128)
**Toolchain boundary**: every command below was run first-party in this worktree at `fd490ca`
(V1, clean tree), shell `zsh`. **Revision round 1** added rows V27–V35 at the same commit
(tree then clean apart from this untracked document); all are POSIX-safe one-liners, and V32
is controller-run with provenance marked. The module floor is active: ambient `go version` reports
`go1.26.6 darwin/arm64` because `go.mod:3` auto-selects it (V2); where a leg named the
pin explicitly, `GOTOOLCHAIN=go1.26.6` was exported. No AILANG (`.ail`) source is written or
changed by this design; the pinned `ailang` binary is not exercised.

> **Thesis:** `P6.T` (squash `8b196c3`) raised the Go floor `go1.25.6 → go1.26.6` at all four
> `ci.yml` pin sites plus `go.mod`, and its own mutation drill recorded the two `setup-go`
> sites as SURVIVED: nothing, locally or in CI, turns a `go-version:` regression red, because
> no Go test in the repo scans either pin kind (C3) and `GOTOOLCHAIN`'s auto-switch covers a
> downgraded runner silently. The repair is one static consistency test — line-anchored key
> extraction with quote-agnostic normalization, per-kind counts tied to the enumerated job
> list, cross-file agreement with the `go.mod` floor — mirroring
> `TestZ3PinDeclaredOnceAndInstalledInBothJobs` in structure and in declared honesty, plus a
> small edit to the mission's own miscompile instrument: the pin joins the known-good list, a
> `PINNED=` constant names it, and a fourth fail-loud block in the script's own voice refuses
> to print the success banner unless the pinned toolchain itself reported OK — closing the
> SKIPPED hole (an unfetchable pin counted toward neither control while the banner printed,
> P5/V33) and making "the instrument probes the live pin" a runtime-enforced property of any
> run **whose exit code is consulted** — which is the attended/local lane, and is **not** CI:
> `ci.yml:172` invokes `run.sh` under `continue-on-error: true`, so no exit code it produces
> has ever reached a build verdict (quorum round 2, `gpt5-6-sol`; scoped rather than
> force-passed, and see V37 + queue row 44 for the measured reason this is worse than the
> reviewer alleged). Test A and Test B are ordinary Go tests executed by `verify_go.sh`'s
> `go test ./...` legs (V16), so **they** gate CI normally; it is only `run.sh`'s own exit
> code that CI discards. The pin's OK is measured, not
> assumed (V32: go1.26.6 fetchable → `OK`; the go1.26.5 negative control still fires `BUG`).
> The instrument's loud-fail machinery already exists (C6); what was missing is a scan that
> keeps the list honest with the floor and a guard that keeps the banner honest with the run.

## The finding in one paragraph

`MUT-TOOLCHAIN-REGRESS` arms 3 and 3′ mutated `actions/setup-go@v5`'s `go-version: '1.26.6'`
back to `'1.25.6'` at `.github/workflows/ci.yml:28` and `:109`; both mutants LANDED
(occurrence count `4 → 3`, a query against the file's own view) and both SURVIVED
(`w-mcp-projection.md:726`, V25), independently confirmed by the sprint evaluator
(`git grep -n "GOTOOLCHAIN\|go-version" host/ cmd/` → 0). This design session reproduced both
arms first-party at `fd490ca`: each mutant landed (`1.26.6` occurrences 4→3, mutated line
read back as `'1.25.6'`), the strongest existing `ci.yml` scanner stayed `ok`
(`TestZ3PinDeclaredOnceAndInstalledInBothJobs` PASS — it guards Z3 pins, not toolchain pins),
the evaluator's zero-scan instrument still read 0 hits against a firing same-scope control,
and both restores were sha256-byte-identical (V10, V11). The mechanism is what makes this
silent rather than merely unguarded: with `GOTOOLCHAIN: go1.26.6` intact, a `setup-go` that
installed 1.25.6 is overridden by the toolchain line on every `go` invocation — CI stays
green and simply pays a second toolchain download per job; a `GOTOOLCHAIN`-only regression
*does* red (`go.mod requires go >= 1.26.6`, drill arm 2, V25), so the exposed drift shape is
exactly **setup-go-only**, in either direction of drift. Meanwhile
`design_docs/verification/w-race-gate-blindspot/run.sh:25` — an executable file (V23), the
mission's first-party instrument for the go1.26 array-literal miscompilation, invoked
non-gating by CI at `ci.yml:172` (V17) — carries `KNOWN_GOOD="go1.25.6 go1.24.9"`: its
known-good arm has never probed `go1.26.6`, the live pin (C4, V6). Its own controls (C6, V8)
guarantee it cannot silently certify an *empty* arm — but their three checks watch
`ran`/`saw_bad`/`saw_good` (`:77–:91`, V33): a SKIPPED pin counts toward neither control
while the success banner prints, and nothing makes the *list* keep pace with the floor. Two
holes — list drift, and a banner that can print without the pin ever probed — are the second
half of this item; the production edit closes both, in the script's own fail-loud voice.

## Premises

Each premise is one or more Verification Log rows; a claim without a row does not appear here.

- **P1 — `ci.yml` has exactly two jobs, each carrying one `GOTOOLCHAIN` env and one
  `setup-go` `go-version`, and all four read 1.26.6 today** (C1 reproduced, V3; job
  declaration enumeration measured independently, V4). The pins agree *now* — the defect is
  that nothing keeps them agreeing.
- **P2 — no Go source or test scans either pin kind.** Zero hits for
  `go-version|GOTOOLCHAIN` across `--include='*.go'`, same-call known-positive control at 12
  (C3 reproduced, V5). The drill survivors are reproduced first-party (V10, V11); the
  exposure is current at `fd490ca`, not historical.
- **P3 — a setup-go-only regression is invisible to every dynamic channel in the job.**
  `GOTOOLCHAIN: go1.26.6` makes every `go` command auto-switch to go1.26.6 regardless of what
  `setup-go` installed; `verify_go.sh`'s deny-list (observer
  `ACTIVE_GO=$(go env GOVERSION)` `:217` feeding the `case` deny-set `:218–:224`, V35; first
  read V22) observes only the ACTIVE toolchain, which auto-switch makes go1.26.6; the job's `go version` step
  (`ci.yml:111–114`) prints but asserts nothing (V17); and `verify_go.sh` runs `go test`
  without `-v`, so nothing observational reaches CI anyway — an assertion is the only channel
  (carried from `w-boundary-gate-tree-mutation` V16c, restated at V16 here). Mechanism argued
  AND the drill measured the arm green (V25).
- **P4 — a `GOTOOLCHAIN`-only (or both-kinds) regression reds dynamically via the module
  floor** (`go.mod requires go >= 1.26.6`, drill arms 1–2, V25), so this item designs for the
  silent shape: setup-go-only drift, in either direction (downgrade covered by a second
  download; an unnoticed *upgrade* float is symmetric).
- **P5 — `run.sh` is executable, tracked mode 100755, and already fails LOUDLY on an empty
  instrument:** exit 1 unless ≥1 KNOWN_BAD reported BUG and ≥1 KNOWN_GOOD reported OK; an
  unfetchable toolchain is SKIPPED and counts toward neither (C6 reproduced by reading, V8;
  the three fail-loud blocks check `ran`, `saw_bad`, `saw_good` at `:77–:91`, V33; exec bit
  V23; invoked non-gating with `continue-on-error: true` at `ci.yml:172`, V17). **None of the
  three checks asks whether the *pinned* toolchain ran** — a SKIPPED pin leaves every block
  green while the banner prints. That is objection 2's structural half, measured true of the
  base script, and the Decision closes it in the script's own voice.
- **P6 — `run.sh:25`'s KNOWN_GOOD predates the `P6.T` pin** and contains no `1.26.6` token
  (C4 reproduced, V6, with a firing same-file control). The instrument certifying
  "known-good" has never *carried* the toolchain the repo runs; the pin's OK was first put on
  measured footing by the controller's three-arm run (V32 — pin `OK`, negative control `BUG`),
  not by any standing gate.
- **P7 — the precedent exists, is green at base, and carries its own residual declaration and
  line-exact repair:** `TestZ3PinDeclaredOnceAndInstalledInBothJobs` at
  `host/verifygate/ail_binary_gate_test.go:668` (V24), PASS at `fd490ca` (V12). Its doc
  comment records (a) a static text scan sees text, never whether a step RUNS, and (b) a
  substring `strings.Count` SURVIVED a suffix-shaped mutation while counting whole trimmed
  lines did not. Both lessons are load-bearing in the Decision below.
- **P8 — `go.mod` declares the same floor, with no `toolchain` directive** (C2 reproduced,
  V26), and the reproducer lives in a **separate nested module** declaring `go 1.22` that the
  root module must not pick up (C5 reproduced, V7). The nested module and the in-module
  canary belong to queue row 42 and are not read further by this design.

### Design Freeze

- **One new file**: `host/verifygate/toolchain_pin_gate_test.go`, `package verifygate`,
  reusing the package-level `repoRoot` (`ail_binary_gate_test.go:27`) — no redeclaration;
  all helper names chosen disjoint from the measured package inventory (V21).
- **One production FILE, still one file**: `run.sh` — `KNOWN_GOOD` gains `go1.26.6`, a
  `PINNED=` constant appears beside the lists, a `saw_pinned_ok` flag joins
  `saw_bad`/`saw_good`/`ran` at `:27–:29`, the `OK*)` arm of `probe()` sets it for the pinned
  toolchain, and a fourth fail-loud block in the file's own voice guards the banner (~10
  lines; base shape and namespace measured at V31/V33). **No `ci.yml` edit** —
  the pins already agree at base (P1); if this design also edited `ci.yml`, AC2's
  red-before-green ordering could not attribute the new tests' teeth to the tests. Flipping
  `ci.yml:172`'s `continue-on-error: true` to gating is **a follow-up, not this item**, on
  exactly this rule — stated in Deferred Scope.
- Neither test calls `requirePinned`/`runGate` or reads `AILANG_BIN` (V14): a static text
  scan must run in any lane, solver or none — exactly like the Z3 precedent, which carries no
  `requirePinned` call (V24).
- No `scripts/verify_go.sh`, `host/store/toolchain_canary_test.go`,
  `host/verifygate/module_manifest_gate_test.go`, `packages/world-core/`, or docs edits:
  rows 42/43 territory (Scope Fence).
- All expectation constants (job names, counts) are **hand-maintained literals**, and all
  *values* are extracted from the artifacts and compared **cross-file**; nothing derives its
  expectation from the artifact it checks (row 43's evaluator lesson, applied).
- No new module dependency (no YAML library); line-oriented parsing plus zero-needles, with
  the exotic-form residual declared rather than patched.

## Decision — one consistency test over `ci.yml`+`go.mod`, one instrument test over `run.sh`, one one-line edit

### Test A — `TestGoToolchainPinsAgreeAndMatchJobList(t)`

1. Read `repoRoot/.github/workflows/ci.yml`. **Known-positive controls first** (mirroring the
   Z3 test's instrument-failure block): the file must contain `ailang-verify:`,
   `go-verify:`, `uses: actions/setup-go@v5`, and `./scripts/verify_go.sh` — counts measured
   1/1/2/1 (V19). A missing control is `t.Fatalf("instrument failure: …")`: a zero-count
   assertion below must never be satisfiable by reading the wrong file.
2. **Enumerate the job list**: lines after the trimmed `jobs:` line matching
   `^  [a-z0-9-]+:$` (exactly two-space indent). Measured shape at base: exactly two, at
   `ci.yml:17` and `:98` (V4); the same shaped lines earlier in the file (`push:`,
   `pull_request:`) precede `jobs:` and are excluded by the seen-jobs flag. Assert set
   equality with the hand-maintained constant `{"ailang-verify", "go-verify"}` — this is the
   "COUNT … is what the job list implies" clause made executable, and the constant's comment
   carries the reviewer rule: *a third job moves this set in the same edit or the test reds.*
3. **Extract per pin kind** with `pinValues(lines, key)`: trim each line, split on the FIRST
   `:`, require exact key equality (so a `go-version-file:` line can never match `go-version`
   — M7), then normalize the value: trim space, strip one matching pair of single OR double
   quotes (M6), prepend `go` if absent. Assert counts: `GOTOOLCHAIN` keyed lines == 2;
   `go-version` keyed lines == 2; `strings.Count(src, "uses: actions/setup-go@")` == 2 — the
   keyed-count tie makes an *unpinned* setup-go step a red (M4), not a silent float. Assert
   all four normalized values equal. **Explicitly not the mechanism:** any
   `strings.Count(src, "1.26.6")`-style lumped substring probe is kind-blind (it cannot
   separate `go1.26.6` from `'1.26.6'`) and is the exact prefix-needle class the Z3 suffix
   repair retired (P7) — right as a drill landing predicate, wrong as a guard.
4. **Zero-needle**: `strings.Count(src, "go-version-file")` == 0 — the alternative
   `setup-go` pin form would bypass keyed extraction silently; adjacent same-call positive
   needle `strings.Count(src, "go-version:")` ≥ 2 keeps the zero honest (M7).
5. **`go.mod` agreement**: read `repoRoot/go.mod`; exactly one `^go ` line; its value,
   normalized identically, must equal the common `ci.yml` value; and zero `^toolchain `
   lines (a toolchain directive is a hidden floor override — M9; none exists today, V26).
   *Scope note, stated rather than smuggled:* the queue text binds "every `GOTOOLCHAIN` value
   and every `setup-go` `go-version`" to each other; including `go.mod` is a 6-line extension
   implementing the item's own exposure sentence ("the two kinds of pin … drift apart") at
   the floor both kinds exist to serve, and it is what makes the NEXT floor raise red unless
   `run.sh` follows (M8). It absorbs nothing from rows 42/43.
6. **Workflow enumeration**: `filepath.Glob(repoRoot/.github/workflows/*)`, mapped through
   `filepath.Base` on every element, must equal exactly `["ci.yml"]` — measured one file at
   base (V20). A second workflow carrying its own `setup-go` pin is a cross-file addition
   this `ci.yml`-only reader cannot see (M14).
   **`filepath.Base` is load-bearing, not tidiness** (quorum round 2, `gemini-3-1-pro`,
   applied verbatim under the narrow-refinement carve-out): `filepath.Glob` returns each
   match with the pattern's directory prefix intact, so a slice equality against `["ci.yml"]`
   would have failed **unconditionally, on every run, at base** — a test that can never pass,
   which no amount of mutation drilling would have diagnosed as a spec defect rather than as
   a repo defect. Confirmed first-party by the controller (V36) rather than accepted on the
   reviewer's authority.

### Test B — `TestMiscompileInstrumentProbesPinnedToolchain(t)`

1. Read `repoRoot/design_docs/verification/w-race-gate-blindspot/run.sh`. Controls:
   `KNOWN_BAD=`, `KNOWN_GOOD=`, `PINNED=`, and the exact shebang line `#!/usr/bin/env bash`
   all present (assignment sites measured at `:24`/`:25`, V23; list contents verbatim at V31;
   `bash -n` rc=0, **V26** — V18 measures actionlint's absence, not this; mis-citation caught
   in review, corrected in round 1). **Shebang needle decision, made and kept exact**
   (objection 1(b)): the control stays `^#!/usr/bin/env bash$`, whose bytes are measured —
   `od -c` shows `# ! / u s r / b i n / e n v   b a s h \n`, anchored count **1** with a
   firing same-file control (V27). The reviewer's false-RED scenario (the script carrying
   `#!/bin/bash`) is refuted as a description of today by that measurement; keeping the exact
   form means a future rewrite of line 1 reds — and rewriting an executable instrument's
   interpreter line *is* drift this control exists to catch, made executable as M18. The
   rejected relaxation ("starts with `#!`, contains `bash`") buys protection against a
   benign-rewrite false RED by paying with interpreter-drift false GREENs; this design chooses
   a red against a real edit over a green against real drift.
2. Parse the `KNOWN_GOOD="…"` and `KNOWN_BAD="…"` assignment lines (line-anchored; single
   assignment each, V23 — contents verbatim at V31) into token lists, and the `PINNED="…"`
   line into a scalar under the same rule.
3. Read `go.mod` **independently of Test A** (no shared locals — each test stands alone) and
   require the normalized floor token (`go` + floor value) ∈ KNOWN_GOOD, **`PINNED=` == that
   token, and `PINNED=` ∈ KNOWN_GOOD** (the probe loop can only set the flag for a version it
   iterates). If a future floor
   raise moves `go.mod` (and `ci.yml`, or Test A reds) but forgets `run.sh`, this reds;
   a raise that updates the list but forgets `PINNED=` reds on the equality alone (M17):
   *the instrument must carry and enforce the toolchain the repo actually pins* — currently
   `go1.26.6`, whose floor membership is absent today (M10, the base state) and whose OK is a
   measured fact (V32).
4. Require the floor token ∉ KNOWN_BAD (M11: the live pin labelled bad).
5. Both lists non-empty (M13) — the runtime loud-fail machinery (P5) covers an empty-arm
   *run*; this covers an empty-arm *edit before any run happens*.
6. `os.Stat` mode & 0111 ≠ 0: CI invokes `./design_docs/verification/w-race-gate-blindspot/
   run.sh` directly (`ci.yml:172`, V17); a dropped exec bit is a silent instrument loss
   (M12).
7. **The pinned-probe guard exists as text**: `strings.Count(src, "saw_pinned_ok") >= 3`
   (declaration beside `saw_bad`; set inside the `OK*)` arm; test in the guard block — the
   three sites named in the production edit, shape measured at V33) **and**
   `strings.Contains(src, "INSTRUMENT FAILURE: the PINNED toolchain")`. A static scan proves
   presence and binding only; the guard's firing behaviour is AC6's runtime trip, thereafter
   exercised on every CI invocation (non-gating, `ci.yml:172`, V17). Neutering the guard
   reds HERE (M16).

### The production edit — one file, ~10 lines: the pin in the list, a `PINNED` constant, a fourth fail-loud block

All edits in `run.sh`; landing sites and variable namespace measured at V31/V33 (no existing
`PINNED`/`saw_pinned` identifier anywhere in the file).

```
# :25 (base contents verbatim, V31):
KNOWN_GOOD="go1.25.6 go1.24.9"
→
KNOWN_GOOD="go1.26.6 go1.25.6 go1.24.9"
PINNED="go1.26.6"   # the toolchain go.mod pins; TestMiscompileInstrumentProbesPinnedToolchain
                    # binds it to the `go` line — a floor raise that forgets it reds (M17).

# beside the flags at :27–:29:
saw_pinned_ok=0

# the OK*) arm of probe()'s case (:50–:53) — one test added, list-position-independent:
OK*) [ "$expect" = GOOD ] && saw_good=1; [ "$tc" = "$PINNED" ] && saw_pinned_ok=1 ;;

# a fourth fail-loud block after :87–:91, before the RESULT banner, in the file's own voice:
if [ "$saw_pinned_ok" -eq 0 ]; then
  echo "INSTRUMENT FAILURE: the PINNED toolchain ($PINNED) never reported OK — it was"
  echo "SKIPPED (unfetchable) or is absent from the probe lists. A success banner that"
  echo "never probed the pin is a false clean; refusing to print it."
  exit 1
fi

# and one banner line so the non-gating CI log (the only channel while ci.yml:172 stands) shows it:
echo "  pinned toolchain ($PINNED) reported OK  : yes"
```

Prepend the pinned version and keep the historic goods — they still probe the pre-bug lineage,
and row 42 owns any restructuring of the known-bad arm. The probe loop otherwise needs no
edit: it iterates the lists as written (V8), and the flag keys on `$tc`, not on list position,
so a row-42 KNOWN_BAD restructure cannot silently defeat the guard. If the pinned toolchain
ever reports BUG (a genuinely broken future pin), the flag never sets and the guard trips
loudly; the label-confusion mutant (M11) is red statically via Test B. **Offline-runner cost,
stated because turning a skip into a hard failure has one:** before this edit an offline
runner printed RESULT with the pin SKIPPED; after, it exits 1 naming the pin. The cost is
bounded twice: the pin is the floor `go.mod:3` demands of every build of this repo (V2, V26),
so a runner that cannot obtain go1.26.6 cannot build the repo at all — the new failure bites
exactly where nothing else could have run either; and CI invokes `run.sh` non-gating
(`continue-on-error: true`, `ci.yml:172`, V17), so the loud exit marks the LOG, never the
build. The trade is deliberate: the banner was the defect.
**And the CI half of that sentence is now measured, not assumed, and it is worse than this
doc previously implied** (controller, quorum round 2, V37): `run.sh` is **already** exiting
non-zero on every CI run — `INSTRUMENT FAILURE (or GOOD NEWS): no known-affected toolchain
reproduced the defect` on **10 of the last 10** `dev` CI runs, with zero `RESULT:` banners
and the control (`== go1.26 local-array-literal miscompilation reproduction ==`) firing every
time — because the miscompilation is **darwin/arm64-specific** while CI is `ubuntu-latest`,
so all four KNOWN_BAD toolchains report `OK` there and `saw_bad` stays 0. `continue-on-error:
true` has been swallowing that for at least ten merges. So this item's new guard is a real
guard in the lane that reads exit codes, and joins an **already-inert** instrument in the lane
that does not. That is a **pre-existing defect at HEAD which this design neither introduces
nor can fix inside 0.15 day** — and whose naive remedy (flip `continue-on-error`) would red
`dev` immediately, for the platform reason above rather than for any pin. It is therefore
filed as **queue row 44 `w-miscompile-instrument-inert-in-ci`** on its own first-party
evidence, per the decomposition rule's "a pre-existing defect surfaced by a reviewer is a
QUEUE ROW, not a revision". Whether `ci.yml:172` should become gating is **row 44's**
question, not this item's, and it cannot be answered by flipping one flag.

### What the gate CANNOT catch — declared residual, in this doc AND in the code comments

The Z3 precedent is trusted *because* it narrows its own claim; both new tests carry the
same narrowing, verbatim in their doc comments. Required comment on Test A:

```
// DECLARED RESIDUAL (mirrors the Z3 precedent at ail_binary_gate_test.go: this is a STATIC
// text scan over YAML. It sees the pin TEXT, never whether the setup-go step RUNS: a
// step-level `if:` whose expression is always false at runtime — e.g. one keyed on a commit
// marker nobody writes — disables the install with every counted byte intact, and no
// actionlint runs anywhere in this repo (measured: 0 hits over *.yml/*.sh, V18), so not even
// a literal `if: false` is flagged. It cannot see what version the runner actually INSTALLED
// (setup-go cache, image drift): the job's `go version` step prints without asserting, and
// verify_go.sh's deny-list observes only the ACTIVE toolchain, which the surviving
// GOTOOLCHAIN pin makes go1.26.6 anyway. And it parses lines, not YAML: a pin folded into a
// flow-style `with: {…}` drops the keyed count (a RED here, not a silent pass), but an exotic
// form that preserved the counts AND smuggled a value past quote-stripping would pass — the
// hand-maintained constants above, not this parser, are the bar.
```

Required comment on Test B:

```
// DECLARED RESIDUAL: this is a STATIC text scan. It proves the LIST and PINNED carry the
// floor token and that the pinned-guard machinery EXISTS as text (the saw_pinned_ok sites
// and the failure message). The guard's firing is proven at sprint time by AC6's guard-trip
// run and exercised on every CI invocation of run.sh (non-gating, continue-on-error: true,
// ci.yml:172) — the round-1 SKIPPED hole (a banner printed with the pin unprobed) is closed
// in run.sh itself, not merely narrowed here. What remains open by scope: nothing WATCHES
// the non-gating log — a loud failure nobody reads is loud only on inspection; flipping
// ci.yml:172 to gating is the named follow-up in Deferred Scope, paired with OD-1.
```

Two further residual statements: (a) this guard's CI channel is the **go-verify** job alone —
`verify_go.sh` runs `go test ./...` in both its plain and `-race` legs (V16), while job 1
executes `./scripts/verify_ail.sh` at `:96` — inside its V4-measured span (job lines `:17` vs
`:98`; invocation sites `:96` vs `:165`, V28) — and never executes Go tests; the guard
therefore verifies job 1's pins *from* job 2's lane, which is sound for a text scan and
stated so.
(b) `verify_go.sh` runs without `-v` (V16), so nothing the tests *log* reaches CI — the
assertions are the channel; the tests carry no `t.Logf`-only evidence.

## Alternatives rejected

1. **A shell grep inside `scripts/verify_go.sh`.** The established home for `ci.yml` static
   structure is `host/verifygate` (Z3 precedent); Go assertions give per-kind counting,
   quoting control and named subfailures that a grep cannot; and `verify_go.sh` is
   load-bearing gate machinery this 0.15-day item should not grow. Rejected.
2. **`host/runbook`.** It already owns one `ci.yml` text scan — the world-publish reachability
   fence (`TestNoCIStepOrScriptReachesThePublishEntrypoint`, `runbook_stageb_test.go:329`,
   targets line `:339`, paired known-positive controls `:363`/`:372`, final assertion
   `:376–:381`; measured at V29, Conflict Surface) — a different concern lineage; pin
   consistency belongs beside the Z3 pin test. Rejected.
3. **A real YAML parser** (`gopkg.in/yaml.v3`): a new module dependency to read four pins in
   a 196-line file (`wc -l` → 196, V30) whose pin lines are key-shaped. Line-anchored parsing plus the
   `go-version-file` zero-needle plus the keyed-count tie carry the rigour; the exotic-form
   residue is declared above rather than paid for with a dependency. Rejected.
4. **Assert the installed version dynamically in the `go version` step** (`test "$(go env
   GOVERSION)" = go1.26.6`): (i) it is a `ci.yml` edit, outside the Design Freeze's one-line
   production bound; (ii) exactly when it matters (setup-go regressed, `GOTOOLCHAIN`
   covering) `GOVERSION` reports go1.26.6 *anyway*, so a faithful version of this assertion
   must re-read `ci.yml`'s text inside the workflow — the static scan re-implemented in shell.
   Carried as Open Decision OD-1, follow-up only.
5. **Deriving expectations from the artifacts under check** (job count from the file used
   without a constant; `run.sh` checked against `ci.yml`'s pin instead of `go.mod`): row 43's
   evaluator refuted self-derived controls as vacuous by construction. Here the job-set and
   counts are hand-maintained constants asserted *equal to* a derived enumeration (neither
   derives the other), and both cross-artifact value checks compare independently-authored
   texts (`ci.yml`↔`go.mod`, `run.sh`↔`go.mod`). Rejected as a mechanism; retained as a rule.

## Ordering

Gated on nothing. No ordering coupling to rows 42/43: this item neither moves the canary's
known-bad control (row 42) nor publishes the floor-raise coupling map (row 43); it *uses* one
fact row 43 will document (the `go.mod` floor is coupled to `ci.yml` and to `run.sh`), and
says so as an explicit dependency of nothing — the coupling is enforced here, mapped there.
**The next floor-raiser's obligation, one sentence** (the "every future module floor" half of
objection 2): since `run.sh` is bound to whatever `go.mod` declares, a floor raise lands only
together with the candidate's *probe* — run `run.sh` against it first (candidate reports `OK`,
a KNOWN_BAD entry still reports `BUG`; the V32 three-arm shape) and move `PINNED=` and
`KNOWN_GOOD` in the same change; Test B reds the raise otherwise (M8, M17), and the guard
turns an unprobeable new pin loud, not green, on the first CI run.

## Files to Create/Modify

- **CREATE** `host/verifygate/toolchain_pin_gate_test.go` (~170 LOC after round 1: Tests A
  and B, the `pinValues` extractor, the `PINNED` scalar parse and guard needles, the
  declared-residual doc comments above). No other package member is
  modified; helper names disjoint from the measured inventory (V21).
- **MODIFY** `design_docs/verification/w-race-gate-blindspot/run.sh`: the KNOWN_GOOD line, the
  `PINNED=` constant, the `saw_pinned_ok` flag plus its case-arm set, the fourth fail-loud
  block, one banner line (~10 lines, this one file; landing sites measured V31/V33).

No other files. `ci.yml`, `go.mod`, every `scripts/*.sh`, `host/store/`,
`host/verifygate/module_manifest_gate_test.go` — all untouched.

## Conflict Surface

- **`TestZ3PinDeclaredOnceAndInstalledInBothJobs`** (`ail_binary_gate_test.go:668`) — same
  package, same target file, overlapping known-positive controls (`ailang-verify:`,
  `go-verify:`). Deliberate: the two tests guard disjoint pins (solver vs toolchain) with the
  same instrument; their control needles are *duplicated in each test*, not shared, so one
  helper's edit cannot blind both. Shared fate is correct, not accidental: a `ci.yml` edit
  that drops a job reds both (Z3's install count 2 and this design's job-set constant 2 are
  the same structural fact), and an edit that removes the `verify_go.sh` invocation reds this
  design's control and the runbook fence (next bullet) together — as it should.
- **`host/runbook/runbook_stageb_test.go`
  `TestNoCIStepOrScriptReachesThePublishEntrypoint` (`:329–:382`, V29)** — the repo's *other*
  existing `ci.yml` reader, fencing `world-publish` reachability; one of its two measured
  same-call known-positive controls (`controlCI := strings.Count(ciYML, "verify_go.sh")` ≥ 1,
  `:363`; the second, `verify_world_package.sh` in `verify_ail.sh`, `:372`) is one of this
  design's control needles.
  No needle overlap otherwise; no shared helpers; noted so a reviewer knows the read exists.
- **`scripts/verify_go.sh` toolchain deny-list (comment block `:214–:216`, observer
  `ACTIVE_GO=$(go env GOVERSION)` `:217`, `case` deny-set `:218–:224` — numbered read V35;
  earlier read V22)** — dynamic, ACTIVE-
  toolchain, negative-set (go1.26.0–1.26.5), and self-documented ("Future go1.26.6 … not
  covered here; the canary … is the version-agnostic detector", `:215–:216`). Disjoint channel from a
  static declared-pin equality check: the deny-list never reads `ci.yml`, this test never
  executes `go`. No contradiction and no edit; a floor raise does not invalidate it (those
  versions remain miscompilers).
- **`host/store/toolchain_canary_test.go` (V22)** — runtime miscompilation detector under the
  active compiler; row 42 owns its known-bad-control redesign and its stale
  "pinned-good Go 1.25.6" comment (`:7`). This design touches nothing in `host/store/`.
- **Row 42's redesign surface inside `run.sh`** — the guard keys on `$PINNED` and the `OK*)`
  arm only (V33 shape), so row 42's anticipated KNOWN_BAD-arm restructure (OD-2) cannot
  silently defeat the pinned requirement; confusion and emptiness still red (M11/M13).
- **`host/verifygate/` package members** — `evidence_manifest_gate_test.go` and the
  row-43-fenced `module_manifest_gate_test.go` share the package; the measured function
  inventory (V21) contains no `pinValues`/Test A/Test B names, so no collision; `repoRoot` is
  package-level (`:27`) and reused, not redeclared.
- **Placement decision** — `host/verifygate` over `host/runbook` or `scripts/`: precedent
  (Z3), zero-config reach (`go test ./...` in both `verify_go.sh` legs, V16), and an
  environment-independent lane (no `requirePinned`, V14). The alternative homes are rejected
  above with reasons.
- **C3's measured zero (V5)** also certifies there is no existing assertion anywhere in the
  Go tree that these pins stay *unscanned* — nothing contradicts the new checks.

## Systemic-Issue Audit

Is an unguarded, drift-capable executable pin a pattern? Enumeration of **every** tracked
`1.25.6`/`1.26.6` string (C7 reproduced, V9): exactly four live pin sites, all in `ci.yml`
(P1), plus the `go.mod:3` floor (P8). Every other hit is documentary — `bench/BASELINE.md`
prose, `implemented/` design docs, the queue/history text of `world-mission.md`, tracked
`sprint_*.json` artifacts, `docs/SELF_MOD_PUBLISH.md:120`, and **comments/advisory strings in
live files** (`verify_go.sh:215,222`; `verify_world_package.sh:226`; canary `:7`) — not
executable pins. Executable toolchain-version homes beyond those five: the `run.sh` known
lists (this item binds them to the floor and now requires the pinned toolchain itself to have
reported OK before the banner prints) and the deny-list (negative-set, explicitly not a
floor pin; nothing to drift). No sixth executable home exists; no evaluator or controller
sweep found one. After this item lands, every executable toolchain pin is either cross-checked
by Test A/B or explicitly disclaimed in its own comment. The fix is correctly local.

## Deferred Scope (cut to hold ≤0.15 day)

- The dynamic installed-version assertion in `ci.yml`'s `go version` step (OD-1) — a `ci.yml`
  edit; needs its own design sentence about the auto-switch mechanism (P3).
- `actionlint` adoption (V18: absent repo-wide) — a tooling decision above this item's pay
  grade.
- **Flipping `continue-on-error: true` at `ci.yml:172` to gating** (named by the round-1
  reviewer) — **a follow-up, not this item**: it is a `ci.yml` edit the Design Freeze
  excludes, and gating a network-dependent instrument turns every transient toolchain-fetch
  failure into a CI red, which wants a flake-rate measurement of its own; queue it paired
  with OD-1, whose dynamic assertion is the stronger gating form. Until then the channel is
  the non-gating log, declared as such in Test B's residual comment.
- Anything in rows 42/43's file sets, per the Scope Fence — named as dependencies-of-nothing
  in Ordering, not absorbed.
- Globbing additional workflow files beyond the constant-one assertion (Decision A6) — the
  constant reds force the revisit; building the general case now is speculative.

## Acceptance Criteria

Each AC carries its vacuity self-test and its **observed result on the unmodified tree at
`fd490ca`**, run this session (Verification Log rows cited).

- **AC1 — both tests exist, RUN, and pass on the post-sprint tree, in run-existence form.**
  `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run
  'TestGoToolchainPinsAgreeAndMatchJobList|TestMiscompileInstrumentProbesPinnedToolchain'
  -count=1 -v` → rc=0 with exactly 2 top-level `=== RUN` lines (one per name) and 2
  `--- PASS`; and a paired nonsense pattern (`TestNoSuchToolchainPinTestZZZ`) in the same
  invocation style prints `[no tests to run]`, proving the instrument says so rather than
  passing vacuously. **Base @fd490ca: ran verbatim → `testing: warning: no tests to run`,
  `ok … [no tests to run]`, rc=0** (V13) — the naive form ("command greens") is GREEN AT BASE
  measuring nothing; the `=== RUN`/nonsense-control clauses are the repair, and they red at
  base (0 of 2). *Not vacuous:* its red-at-base is measured, not hypothetical.
- **AC2 — test-first ordering: Test B reds before the run.sh edit and greens after, while
  Test A greens throughout.** Sprint evidence = two recorded runs: (i) test file applied,
  `run.sh` untouched → Test B FAILS naming the first missing piece (floor token ∉ KNOWN_GOOD,
  `PINNED=` absent, guard absent — all three are the measured base state, V6/V34), Test A
  PASSES (the pins agree at base, P1 — a Test A red at this stage means the test misreads a
  consistent file); (ii) the `run.sh` edit applied → both PASS. **Base @fd490ca:** the pair
  cannot be executed at base (no test file); its premise is measured instead — `grep -c
  "1\.26\.6" run.sh` → **0** with a firing same-file control `"1\.26\.5"` → **1** (C4, V6).
  *Not vacuous:* if `run.sh` already probed the pin, AC1 alone could not distinguish teeth
  from luck; only the red-then-green pair discharges this. A sprint reporting only the final
  green has NOT discharged AC2.
- **AC3 — the eighteen named mutations red the named tests on the post-sprint tree, each
  restored byte-identically (per-arm sha256, house recipe), pristine control green after
  every arm.** (M15–M18 were added in revision round 1; each names Test B.) M1/M2 carry special status: they are the `P6.T` drill's recorded SURVIVORS
  (`w-mcp-projection.md:726`, V25) and they SURVIVE at `fd490ca` today — reproduced
  first-party this session: mutated `:28`/`:109` landed (`1.26.6` occurrences 4→3, line read
  back `'1.25.6'`), the strongest existing `ci.yml` scanner stayed `ok`, the evaluator's
  zero-scan read 0 against a firing control, restores byte-identical (V10, V11).
  *Not vacuous:* M1/M2's survival at base is two independent measurements old and one
  session-fresh; the same mutants must print Test A's red after this item.
- **AC4 — the tests are environment-independent.** `env -u AILANG_BIN go test
  ./host/verifygate/ -run '<A>|<B>' -count=1` → rc=0 with AC1's run-existence clauses; no
  network, no solver, no pinned binary. **Base @fd490ca:** the same command → `no tests to
  run`, rc=0 (V13, red under the repaired form); and the package *without* `AILANG_BIN`
  FATALs its shim-arm tests via `requirePinned` (V14) — the measured reason either new test
  calling `requirePinned` reds AC4 while AC1 stays green on a rig that has the variable:
  exactly the split this AC exists to catch.
- **AC5 — hygiene of the new file:** `go vet ./host/verifygate/` rc=0 and `gofmt -l
  host/verifygate/` prints nothing. **Base @fd490ca: both green (V15) — green-at-base by
  DESIGN:** they measure the sprint's own new file, not repo state; they fail only if the new
  file is malformed, and are listed so the sprint's final-tree gate list is complete rather
  than assumed.
- **AC6 — the pinned-probe guard exists, reds at base by absence, and FIRES at runtime on the
  post-sprint tree** (revision round 1, objection 2 option A). Sprint evidence = three
  recorded results: (i) an attended `run.sh` on the post-sprint tree → rc=0, output shows
  `go1.26.6  expect=GOOD got: OK (rc=0)` and the banner's pinned line — the OK claim's
  precedent, with the negative control firing, is measured at V32; (ii) **guard-trip**,
  rehearsing exactly the reviewer's SKIPPED scenario (and the offline runner): temporarily
  set `PINNED="go1.99.9"` and place `go1.99.9` first in KNOWN_GOOD → the run prints its
  `SKIPPED` line, then `INSTRUMENT FAILURE: the PINNED toolchain (go1.99.9) never reported
  OK`, rc=1, **no RESULT banner**; restore sha256-byte-identical, house recipe. **Base
  @fd490ca:** the mechanism is absent — `grep -c 'PINNED\|saw_pinned' run.sh` → **0** with
  same-file control `grep -c 'KNOWN_' run.sh` → **4** (V34) — so every clause of this AC reds
  on the unmodified tree. *Not vacuous:* absence-at-base is measured against a firing
  control; the trip run proves the guard fires rather than merely compiles; and M16/M17 red
  Test B if the guard or its floor binding is neutered later.

**Round-1 fix → guard mapping** (the mission's own guardrail — a fix applied in review changes
the design and must change what guards it):

- Pinned-guard block + `saw_pinned_ok` flag → absent: AC6 reds at base (measured, V34);
  neutered later: **M16** reds Test B step 7; firing proven by AC6's trip run and exercised by
  every CI invocation of `run.sh`.
- `PINNED=` constant + its floor binding → absent: Test B steps 1–3 red at AC2's red stage;
  stale after a floor raise: **M17** reds Test B step 3.
- `KNOWN_GOOD` prepend → absent: **M10** is the measured base state (V6); AC2's
  red-before-green discharges it.
- Exact-shebang needle (kept exact per objection 1's stated choice) → dropped or relaxed, and
  interpreter drift would green: repaired by **M18**, the reviewer's own `#!/bin/bash`
  scenario as a named RED mutation.
- Citation corrections (V18→V26 at Test B step 1; V3→V28; "V21-adjacent"→V29 plus the
  Conflict-Surface runbook bullet; V19/V20→V30; V23→V31 at OD-2; both deny-list ranges→V35)
  → these change no mechanism and need no mutation; each replacement row carries command plus
  observed output, with same-call controls where absence-shaped (V27, V34), and AC1–AC6
  re-execute every assertion the citations support, so a mis-bound citation surfaces as an
  unexplained red rather than a quiet wrong reference.

Explicitly rejected as an AC: "the drills M1/M2 were previously run" — history proves
survival-at-base, not teeth-after-change; AC3 re-runs them. And: a package-wide `ok` alone —
`go test` prints `ok` for a package whose named tests never ran (V13).

## Non-Vacuity — named RED mutation for every added assertion

Production side mutated (`ci.yml`, `go.mod`, `run.sh` incl. its shebang line, workflows dir,
file mode) — never the
test helper. Assertion coverage: A2←M5, A3←M3, A4←M1/M2/M4/M6, A5←M7, A6←M8, A7←M9,
workflow-enumeration←M14, B1←M15/M18, B2←M10/M13, B3←M11, B4←M12, B5←M13, B6←M8/M17,
B7←M16.

| # | Exact edit | Expected RED (single test name) | Shape |
|---|---|---|---|
| M1 | `ci.yml:28` `'1.26.6'` → `'1.25.6'` (`P6.T` arm 3 — SURVIVES at base, V10) | `TestGoToolchainPinsAgreeAndMatchJobList`: setup-go kind values disagree / cross-kind mismatch vs the `go1.26.6` GOTOOLCHAIN values | threat-shaped: the measured silent survivor |
| M2 | `ci.yml:109` `'1.26.6'` → `'1.25.6'` (arm 3′ — SURVIVES at base, V11) | `TestGoToolchainPinsAgreeAndMatchJobList`, same assertion | threat-shaped: the measured silent survivor |
| M3 | `ci.yml:21` `GOTOOLCHAIN: go1.26.6` → `go1.25.6`, setup-go pins untouched | `TestGoToolchainPinsAgreeAndMatchJobList`: GOTOOLCHAIN kind internal disagreement | drift-shaped: the runtime floor already trips this (drill arm 2, V25) — the static scan must red the same site so both channels name it |
| M4 | delete the `go-version: '1.26.6'` line at `ci.yml:28` (unpinned setup-go step floats to the runner default) | `TestGoToolchainPinsAgreeAndMatchJobList`: keyed `go-version` count 2→1 while `uses: actions/setup-go@` stays 2 — the count tie fires | removal: proves the check FIRES |
| M5 | append a third job (copy of `go-verify`, renamed, identical `go1.26.6`/`'1.26.6'` pins) | `TestGoToolchainPinsAgreeAndMatchJobList`: enumerated job set `[ailang-verify go-verify <new>]` ≠ the constant set — **and the observed pin counts move 2→3 in the red message** | **ADDITION** (enumeration-completeness): counts move, not merely the verdict |
| M6 | `ci.yml:28` → `go-version: "1.25.6"` (same regression, double quotes) | `TestGoToolchainPinsAgreeAndMatchJobList`: value disagreement — proves the extractor strips BOTH quote styles; a needle on `'1.25.6'` misses this entirely | the Z3 suffix-repair lesson applied to quoting |
| M7 | add `go-version-file: 'go.mod'` inside job 1's setup-go `with:` block | `TestGoToolchainPinsAgreeAndMatchJobList`: zero-needle `go-version-file` 0→1; exact-key split also keeps it out of the `go-version` count | **ADDITION** of the alternative pin form |
| M8 | `go.mod:3` `go 1.26.6` → `go 1.27.0`, `ci.yml` untouched | REDS TWO tests by design — `TestGoToolchainPinsAgreeAndMatchJobList` (cross-file floor mismatch) AND `TestMiscompileInstrumentProbesPinnedToolchain` (KNOWN_GOOD lacks `go1.27.0` *and* `PINNED=` ≠ the new floor). Double coverage is intended: two artifacts guard the same floor for different consumers | drift-shaped: the next `P6.T`-without-this-item |
| M9 | add `toolchain go1.25.6` directive to `go.mod` | `TestGoToolchainPinsAgreeAndMatchJobList`: `^toolchain ` zero fires | addition of a hidden floor override |
| M10 | `run.sh:25` KNOWN_GOOD lacking `go1.26.6` — **this is the `fd490ca` BASE STATE itself** | `TestMiscompileInstrumentProbesPinnedToolchain` B2 red at base (AC2): `KNOWN_GOOD … does not probe the pinned toolchain go1.26.6` | threat-shaped; red-at-base measured via C4/V6 |
| M11 | insert `go1.26.6` into KNOWN_BAD | `TestMiscompileInstrumentProbesPinnedToolchain` B3 disjointness | confusion-shaped: the live pin labelled bad |
| M12 | `chmod -x design_docs/verification/w-race-gate-blindspot/run.sh` | `TestMiscompileInstrumentProbesPinnedToolchain` B4 exec-bit (CI invokes `./run.sh` directly, `ci.yml:172`, V17) | mode-shaped |
| M13 | `KNOWN_GOOD=""` | `TestMiscompileInstrumentProbesPinnedToolchain` B2+B5; independently `run.sh`'s own controls exit 1 at runtime (C6, V8) — two channels | empty-instrument-shaped |
| M14 | add `.github/workflows/other.yml` carrying its own setup-go pin | `TestGoToolchainPinsAgreeAndMatchJobList`: workflow-glob constant (exactly one workflow file, V20) fires | **cross-file ADDITION** a `ci.yml`-only reader cannot see |
| M15 | delete the `PINNED=` line from post-edit `run.sh` | `TestMiscompileInstrumentProbesPinnedToolchain` step 1 control / step 2 scalar parse — the constant the guard hangs on is gone (B1) | removal: the round-1 binding |
| M16 | neuter the guard — delete the `if [ "$saw_pinned_ok" -eq 0 ] …` block OR the `saw_pinned_ok=1` set in the `OK*)` arm | `TestMiscompileInstrumentProbesPinnedToolchain` step 7: `saw_pinned_ok` site count drops below 3 / the failure-message needle is gone — **the round-1 reviewer's SKIPPED hole, re-opened, reds** (B7) | removal of runtime enforcement |
| M17 | stale PINNED: post-edit `PINNED="go1.26.6"` → `PINNED="go1.25.6"` (a token that IS in KNOWN_GOOD — membership and non-empty clauses all pass) | `TestMiscompileInstrumentProbesPinnedToolchain` step 3: `PINNED=` ≠ `go.mod` floor — **the next floor-raise's forgetful shape, rehearsed** (B6) | drift-shaped |
| M18 | rewrite `run.sh:1` `#!/usr/bin/env bash` → `#!/bin/bash` (objection 1's scenario, made executable) | `TestMiscompileInstrumentProbesPinnedToolchain` step 1: exact-shebang control fails (B1) — the chosen risk demonstrated: a red against a real edit, never a green against interpreter drift | benign-looking rewrite |

Green control for all arms: the unmutated post-sprint tree passes AC1/AC4, and every arm ends
restored sha256-byte-identical with `git status --porcelain` empty — the house recipe this
session's V10/V11 replay already ran twice at base.

## Open Decisions

- **OD-1 — dynamic installed-version assertion in the go-verify job's `go version` step**
  (Alternative 4). *Controller default if no one answers: not this item; file as a follow-up
  note.* It is a `ci.yml` edit (Design Freeze) and its faithful form re-implements the static
  scan inside the workflow.
- **OD-2 — should Test B pin KNOWN_BAD's exact four-token set** (`go1.26.0 go1.26.3 go1.26.4
  go1.26.5` — the assignment line verbatim at V31; V23 measured only the site)? *Default: NO.* Row 42 owns known-bad placement and may move that arm into
  the nested repro module; an exact-set assertion here would fight that item. Non-empty
  (B5) plus disjoint-from-floor (B3) is the honest bound at row 41's scope.

## Verification Log

All rows run by the designer at `fd490ca` in this worktree, shell `zsh`. Controller premises
C1–C7 (iteration 128) were reproduced first-party and are marked; where a controller
parenthetical disagreed with the reproduction, the reproduction is recorded and the
disagreement stated. KP = known-positive control carried in the same call. Rows V27–V35 were
added in **revision round 1** at the same commit (tree then clean apart from this untracked
document; all POSIX-safe one-liners); V32 was run by the **controller** this session and is
carried in with provenance marked — the designer did not re-run it.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Worktree is `fd490ca`, clean, and all mutation replays below ended restored | `git rev-parse HEAD && git status --porcelain \| wc -l` (re-run after each arm) | `fd490cafe5292c02…`, `0`; `0` after every arm |
| V2 | Toolchain boundary: the floor auto-selects go1.26.6 | `go version`; `go env GOVERSION`; `GOTOOLCHAIN=go1.26.6 go env GOVERSION` | `go1.26.6 darwin/arm64` all three |
| V3 (C1) | `ci.yml`: two jobs; pin lines at 21/26/28 and 102/107/109, all 1.26.6 | `grep -n "^jobs:\|GOTOOLCHAIN\|go-version\|uses: actions/setup-go" .github/workflows/ci.yml` | `16:jobs:`, `21 GOTOOLCHAIN: go1.26.6`, `26 setup-go@v5`, `28 go-version: '1.26.6'`, `102 GOTOOLCHAIN: go1.26.6`, `107 setup-go@v5`, `109 go-version: '1.26.6'` — matches C1's substantive sites (job-name lines read separately at V4) |
| V4 | The job list is exactly {ailang-verify, go-verify} | `awk '/^jobs:/{p=1;next} p && /^  [a-z0-9-]+:$/{print NR": "$0}' ci.yml`; `grep -c "runs-on:" ci.yml` | `17: ailang-verify:`, `98: go-verify:` — exactly two; `runs-on:` count 2 |
| V5 (C3) | NO Go file scans either pin; the control distribution is 7/4/1 | `grep -rn "go-version\|GOTOOLCHAIN" --include='*.go' . \| wc -l`; KP `grep -rc "ci.yml\|workflows" --include='*.go' .` (files with hits) | scan **0**; KP total **12** across `ail_binary_gate_test.go`:7, `runbook_stageb_test.go`:4, `registry.go`:1 (a comment, `nomination workflows.`) — **row 41's "five in ail_binary_gate_test.go" parenthetical is off; the file carries seven`ci.yml`\|`workflows` lines; the total 12 stands** |
| V6 (C4) | `run.sh` never names the live pin | `grep -c "1\.26\.6" run.sh`; KP same file `grep -c "1\.26\.5" run.sh` | **0** (the KNOWN_GOOD gap); KP **1** — instrument sees the file |
| V7 (C5) | The reproducer is a separate nested module | `cat design_docs/verification/w-race-gate-blindspot/repro/go.mod` | `module ailang-world/verification/go1_26-arraylit-miscompile`, `go 1.22`, header comment: root module must NOT pick it up |
| V8 (C6) | `run.sh` fails loudly on empty arms; SKIPPED counts toward neither control | read `run.sh` in full | exits 1 unless ≥1 KNOWN_BAD=BUG and ≥1 KNOWN_GOOD=OK; unfetchable toolchain → `SKIPPED`, returns 0 without setting either `saw_*` flag; `KNOWN_BAD=` at `:24`, `KNOWN_GOOD=` at `:25` |
| V9 (C7) | The only other tracked 1.25.6/1.26.6 strings are non-pins | `git grep -n "1\.25\.6\|1\.26\.6" -- .` | four live `ci.yml` sites (V3) + `go.mod:3`; everything else: `bench/BASELINE.md` prose, `implemented/` docs, `world-mission.md` queue/history text, tracked `sprint_*.json`, `docs/SELF_MOD_PUBLISH.md:120`, and **comments/advisory strings in live files** (`verify_go.sh:215,222`; `verify_world_package.sh:226`; canary `:7`). **C7's "BASELINE.md or implemented/ doc" parenthetical under-classifies; the fuller taxonomy is this row and the Systemic-Issue Audit's** |
| V10 | **Drill arm 3 replay at base: `:28` mutant SURVIVES everything that exists today** | sha256 baseline `e078eb6e…`; python sets `:28` → `'1.25.6'` (asserted pre-image); landing by occurrence count + line read-back; `go test -run TestZ3Pin… -count=1`; evaluator instrument `git grep -c "GOTOOLCHAIN\|go-version" -- host/ cmd/`; `cp` restore, sha256 compare | landed: `1.26.6` occurrences **4→3**, :28 read back `go-version: '1.25.6'`; scanner `ok … 0.434s` (**PASS = survived**); evaluator grep rc=1, **0 hits** (**survived**; same-scope control `verify_go.sh` fired 6/6 files in this session's batch); restore **byte-identical** `e078eb6e…`, porcelain 0 (V1) |
| V11 | **Drill arm 3′ replay at base: `:109` mutant SURVIVES equally** | same recipe against `:109`, KP control inline this time | landed 4→3, :109 read back `'1.25.6'`; scanner `ok … 0.269s`; evaluator grep **0 hits**, KP control **6** files with hits — the zero is a measurement, not a broken instrument; restore byte-identical, porcelain 0 |
| V12 | The precedent is green at base | `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run 'TestZ3PinDeclaredOnceAndInstalledInBothJobs' -count=1 -v` | `=== RUN` / `--- PASS (0.00s)` / `ok … 0.300s` |
| V13 | **The vacuity trap AC1 is repaired against, measured at base** | `go test ./host/verifygate/ -run 'TestGoToolchainPinsAgreeAndMatchJobList\|TestMiscompileInstrumentProbesPinnedToolchain' -count=1 -v`; nonsense control `-run TestNoSuchToolchainPinTestZZZ` | both: `testing: warning: no tests to run`, `ok … [no tests to run]`, **rc=0** — a bare "command greens" AC would be green at base while the tests do not exist |
| V14 | Package environment contract: unset AILANG_BIN FATALs the shim-arm tests | `env -u AILANG_BIN go test ./host/verifygate/ -count=1` | `TestAcceptContractRejectsASolverlessGate` et al. FAIL with `AILANG_BIN is unset — … Never skip: a silent skip here is the false-green class…` — **why the new tests must not `requirePinned`** |
| V15 | Hygiene baselines | `GOTOOLCHAIN=go1.26.6 go vet ./host/verifygate/`; `gofmt -l host/verifygate/ \| wc -l` | rc=0; `0` |
| V16 | Where the new tests execute in CI; the assertion-not-log rule holds | `grep -n "go test" scripts/verify_go.sh` | `:245` `go test -json ./host/evidence`, `:258` `go test ./... -count=1`, `:262` the `-race -timeout 8m` leg — **no `-v` anywhere**, so only failures (assertions) reach CI logs |
| V17 | `run.sh`'s CI invocation is direct-exec and non-gating; the `go version` step asserts nothing | `sed -n '111-114,166-176p' ci.yml` | `:111-114` `go version` / `go env GOVERSION`, print only; `:172` `run: ./design_docs/verification/w-race-gate-blindspot/run.sh` under `continue-on-error: true`, `timeout-minutes: 15` |
| V18 | No actionlint anywhere (the Z3 residual's linter backstop does not exist here) | `grep -rn "actionlint" --include='*.yml' --include='*.sh' . \| wc -l`; KP same-scope `grep -c "run:" ci.yml` | **0**; KP **13** |
| V19 | Known-positive control needle counts in `ci.yml` | `grep -c` per needle | `ailang-verify:` **1**, `go-verify:` **1**, `actions/setup-go@v5` **2**, `verify_go.sh` **1**, `verify_ail.sh` **3** |
| V20 | `ci.yml` is the only workflow file | `ls .github/workflows/` | `ci.yml` only (basis of the workflow-glob constant, Decision A6/M14) |
| V21 | Package function inventory — no name collisions for the new helpers | `grep -h "^func " host/verifygate/*_test.go` | 49 declarations incl. `findRepoRoot` (`:31`), package var `repoRoot` (`:27`); no `pinValues`, no Test A/B names |
| V22 | The two adjacent mechanisms this design must not collide with | read `scripts/verify_go.sh:200-235`; read `host/store/toolchain_canary_test.go` | deny-list `:216-:224` = ACTIVE-toolchain negative set (go1.26.0–1.26.5) with the "Future go1.26.6 … not covered here" comment `:215`; canary is 47 lines, runtime-only, stale "pinned-good Go 1.25.6" comment `:7` (row 42's) |
| V23 | `run.sh` mode and pin-list site | `test -x …/run.sh`; `git ls-files -s …/run.sh`; `grep -n "KNOWN_BAD=\|KNOWN_GOOD=" …/run.sh` | `EXECUTABLE yes`; mode `100755`; `:24` / `:25` |
| V24 | The precedent's exact seat and structure | `grep -n "func TestZ3PinDeclaredOnceAndInstalledInBothJobs" host/verifygate/ail_binary_gate_test.go`; read `:600-:780` | `:668`; doc comment records the DECLARED RESIDUAL (text scan ≠ step runs) and the line-exact substring repair; the test body carries no `requirePinned` call |
| V25 | The drill record this item repairs | `grep -n "MUT-TOOLCHAIN-REGRESS" design_docs/planned/w-mcp-projection.md` | `:726` — arm 1 (`go.mod`) reds, arm 2 (`ci.yml:21/:102`) reds with the floor message verbatim, **arm 3 (`ci.yml:28`/`:109`, setup-go sites): SURVIVED, "no command locally or in CI turns it red… an honest gap… guarding it is a separate queue row"** |
| V26 | `go.mod` floor line; no toolchain directive; `run.sh` parses | `cat go.mod \| head -8`; `grep -n "^go \|^toolchain" go.mod`; `bash -n …/run.sh` | `module github.com/sunholo-data/ailang-world`, `3:go 1.26.6`, no `toolchain` line (C2); `bash -n rc=0` |
| V27 | The shebang's exact bytes — objection 1's control string, measured (round 1) | `head -1 …/run.sh \| od -c`; `grep -c '^#!/usr/bin/env bash$' …/run.sh`; KP same file `grep -c '^KNOWN_GOOD=' …/run.sh` | bytes `# ! / u s r / b i n / e n v   b a s h \n` — exactly `#!/usr/bin/env bash`; anchored count **1**; KP **1** (the control fires, so the 1 is a measurement). The `#!/bin/bash` false-RED scenario is refuted as a description of today; it survives only as named mutation M18 |
| V28 | Which job executes which script — citation repair for residual (a) (round 1) | `grep -n 'verify_ail\.sh\|verify_go\.sh' .github/workflows/ci.yml` | `:10` a comment, `:96` `run: ./scripts/verify_ail.sh` (inside job 1's span — job lines `:17`/`:98`, V4), `:150` a comment, `:165` `run: ./scripts/verify_go.sh` (inside job 2). V3's grep never named the scripts; this measures the placement |
| V29 | The runbook `world-publish` fence, measured — replaces "V21-adjacent" (round 1) | read `host/runbook/runbook_stageb_test.go:329–:382` | `TestNoCIStepOrScriptReachesThePublishEntrypoint` at `:329`; zero-enumeration `t.Fatal` `:336–:338`; targets `ci.yml` + `scripts/*.sh` at `:339`; TWO same-call known-positive controls — `controlCI := strings.Count(ciYML, "verify_go.sh")` `:363`, `controlScript := strings.Count(verifyAIL, "verify_world_package.sh")` `:372`; final assertion `:376–:381` |
| V30 | `ci.yml`'s line count — the "196-line file" claim Alternatives 3 cited V19/V20 for (round 1) | `wc -l .github/workflows/ci.yml` | `196 .github/workflows/ci.yml` — the number V19/V20 never carried |
| V31 | `run.sh`'s list lines verbatim — the tokens OD-2 quotes; V23 recorded only site numbers (round 1) | `sed -n '24,25p' …/run.sh` | `:24 KNOWN_BAD="go1.26.0 go1.26.3 go1.26.4 go1.26.5"`; `:25 KNOWN_GOOD="go1.25.6 go1.24.9"` |
| V32 | **The pin's OK is a measured fact** — controller-run this session, carried in verbatim; designer did not re-run (round 1) | three arms, identical `repro/` source, one variable `GOTOOLCHAIN`, each built then executed | `go1.26.6` (live pin): build ok → **`OK`**, rc=0; `go1.26.5` (known-bad negative control): build ok → **`BUG: Field="" want "stateRoot"`**, rc=0 — **the discriminator still fires**; `go1.25.6` (historic known-good): build ok → **`OK`**, rc=0. Labelling go1.26.6 KNOWN_GOOD asserts this measurement, not an assumption |
| V33 | `run.sh`'s full numbered shape — the edit's landing sites and namespace (round 1) | read `…/run.sh` in full, numbered | shebang `:1`; lists `:24`/`:25`; flags `saw_bad=0 saw_good=0 ran=0` `:27–:29`; `probe()` `:31–:55` with the SKIPPED path `:34–:38` returning 0 touching no flag, empty-output instrument failure `:46–:49`, `OK*)`/`BUG*)` arms `:50–:53`; the THREE fail-loud blocks `:77–:80` (`ran`), `:81–:86` (`saw_bad`), `:87–:91` (`saw_good`); RESULT banner `:92–:94`; **no existing `PINNED`/`saw_pinned` identifier** |
| V34 | The round-1 mechanism is ABSENT at base — AC6's red-at-base, with control (round 1) | `grep -c 'PINNED\|saw_pinned' …/run.sh`; KP same file `grep -c 'KNOWN_' …/run.sh` | **0**; KP **4** — the absence is measured against a firing control, not assumed |
| V35 | Deny-list exact seats — repairs the two fuzzy ranges (P3's `:216–:224`; Conflict Surface's `:214–:224`) with a numbered read (round 1) | `awk 'NR>=213 && NR<=226 {print NR": "$0}' scripts/verify_go.sh` | comment block `:214–:216` ("measured set" `:214`; "Future go1.26.6 … not covered here; the canary … version-agnostic detector" `:215–:216`); observer `ACTIVE_GO=$(go env GOVERSION)` `:217`; `case` deny-set go1.26.0–go1.26.5 `:218–:224`; `:225` blank; `:226` the next section banner |

| V36 | **`filepath.Glob` returns directory-prefixed paths** — quorum round 2, `gemini-3-1-pro`'s objection, measured rather than accepted on authority | standalone `go run` over a fixture dir containing exactly one `ci.yml`: `filepath.Glob(filepath.Join(dir,".github","workflows","*"))`, printed `%#v`, then the same slice mapped through `filepath.Base` | `results=[]string{"/tmp/w128glob/.github/workflows/ci.yml"}` — the prefix IS present, so the pre-fix Test A step 6 would have failed unconditionally at base; `basenames=[]string{"ci.yml"}` — the applied fix produces the compared value. Objection CONFIRMED |
| V37 | **`run.sh` is ALREADY inert in CI, and the cause is the host platform, not the pin** — quorum round 2, `gpt5-6-sol`'s objection, measured; the finding is larger than the objection and is queue row 44, not a revision | per `dev` CI run: `gh api repos/…/actions/jobs/<id>/logs`, then `grep -c "INSTRUMENT FAILURE (or GOOD NEWS)"`, `grep -c "RESULT: reproduction confirmed"`, and same-log control `grep -c "miscompilation reproduction"`; swept over the last 10 runs | **10 of 10** runs (`fd490ca`, `1cc8cf4`, `699f592`, `592a221`, `8b196c3`, `797211a`, `c706840`, `fcf18fa`, `2e44e3e`, `612828b`): `instrument_failure=1`, `result_banner=0`, control `=1` — the script ran every time and failed every time. CI log for `fd490ca` shows all four KNOWN_BAD reporting `OK` on `ubuntu-latest`. **Platform two-arm control, one variable, identical `repro/` source:** `go1.26.5` on darwin/arm64 → `BUG: Field="" want "stateRoot"` (V32); on linux/amd64 → `OK`. `verify_go.sh:214` already records the deny-set as "the measured set … **on darwin/arm64**". So `saw_bad=0` on every CI run, the script correctly exits 1, and `continue-on-error: true` discards it |

**Revision-round-1 citation audit** (objection 1(c)): every `Vnn` citation outside the log's
own row labels was re-read against the row it cites — **79 sites checked, 8 sites corrected**
(7 distinct defects; the deny-list fix spans two sites): (1) Test B step 1's `bash -n` cited
V18, which measures actionlint's absence — corrected to V26, the row that ran `bash -n`;
(2) Test B step 1's shebang control cited nothing — the gap closed at V27 (exact bytes,
anchored count, same-call control), and the needle kept exact by decision recorded at step 1;
(3) residual (a)'s "ailang-verify runs verify_ail.sh" cited V3, whose grep never named the
scripts — corrected to V28; (4) Alternatives 2's "V21-adjacent" — not a row — corrected to
V29; (5) the Conflict Surface runbook bullet's control-needle claim was uncited — now V29;
(6) Alternatives 3's "196-line file" cited V19/V20, neither of which measured a length —
corrected to V30; (7) OD-2's four-token quote cited V23, which recorded only line numbers —
tokens now verbatim at V31; (8) P3 and the Conflict Surface carried two different deny-list
ranges, neither matching a numbered read exactly — both now cite V35. Every other site
(thesis, boundary note, finding, P1–P8, the Freeze, both tests' steps, both residual
comments, the alternatives, ordering, the rest of the conflict surface, systemic audit,
deferred scope, AC1–AC5, the mutation table, the open decisions) cited the row that actually
measures the fact.


**Quorum verification log — round 2 disposition (controller, narrow-refinement carve-out).**
Round 2 BLOCKED with both reviewers PRESENT and EXTERNAL (`absent_reviewers` empty, metered
$0.1407). Both objections were MEASURED first-party before disposition, not forwarded:

- `gemini-3-1-pro` (`filepath.Glob` prefix) — **CONFIRMED by measurement (V36)** and applied
  **verbatim** under the carve-out: Test A step 6 now maps every match through
  `filepath.Base`. Concrete reviewer-authored fix; disputes no design direction; the pre-fix
  form could never have passed, so this is a spec defect the mutation drill would have
  mis-attributed to the repo. Its guard is unchanged (M14 still reds on a second workflow
  file) — the carve-out changed the extractor, not the assertion, so no new AC is owed
  (mission guardrail iter-98 checked and discharged: the fix is *inside* an assertion M14
  already reds, not a new branch).
- `gpt5-6-sol` (the guard sits behind `continue-on-error: true`) — **CONFIRMED, and the
  measurement made it strictly larger than filed (V37).** The reviewer argued a *hypothetical*
  green over a failing instrument; the measured state is that the instrument has been failing
  on **10 of the last 10** `dev` runs already, for a **platform** reason (darwin-only
  miscompilation vs an `ubuntu-latest` runner) that has nothing to do with any toolchain pin.
  That is a defect this doc **fails to fix**, not one it **introduces** — so per the
  decomposition rule's clause (c) it is filed as **queue row 44**
  `w-miscompile-instrument-inert-in-ci` on its own evidence, and the doc's claim is **scoped**
  to the lane whose exit code is actually consulted (thesis + §"The one production line").
  It is explicitly **not** force-passed: nothing the reviewer said is overridden, and the
  remedy it implies (flip `ci.yml:172` to gating) is refused here on measurement — it would
  red `dev` on the next push, for the platform reason, not for a pin. This is a controller
  ROUTING call and is deliberately **not** `needs-human-review`: filing it as one would
  manufacture a decision the human does not have (standing rule 8).

**What this disposition does NOT claim.** The carve-out is not a pass. Row 41's static tests
(A and B) gate CI normally because `verify_go.sh` runs `go test ./...` (V16); `run.sh`'s own
runtime guard gates only the attended/local lane until row 44 is answered. Both halves are
stated in the thesis so no downstream reader inherits the wider claim.

## Related Documents

- [`w-mcp-projection.md`](w-mcp-projection.md) — `P6.T` (squash `8b196c3`, five single-line
  edits, two files) and the AC15 drill row at `:726` whose arm-3 SURVIVORS are this item's M1/M2;
  its round-5 N7 sweep surfaced the `run.sh:25` known-good gap folded in here as Test B.
- `design_docs/world-mission.md` queue rows 41 (this item), 42 (`w-canary-control-…`) and 43
  (`w-floor-raise-coupling-inventory`) — Scope Fence neighbours, named and not absorbed.
- [`../implemented/w-boundary-gate-tree-mutation.md`](../implemented/w-boundary-gate-tree-mutation.md)
  — the format this doc follows and the mission's running rule this item inherits: *a green
  result must be unable to mean "the check never ran."*
- `design_docs/verification/w-race-gate-blindspot/` — the instrument `run.sh` (modified by one
  line here), its nested `go 1.22` repro module (row 42's, untouched), and why
  go1.26.0–1.26.5 is deny-listed.
- The precedent: `TestZ3PinDeclaredOnceAndInstalledInBothJobs`,
  `host/verifygate/ail_binary_gate_test.go:668` — structure, controls, residual honesty, and the
  line-exact repair this design mirrors.

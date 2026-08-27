# w-canary-control-does-not-survive-a-floor-raise — the known-bad arm survives only because one unguarded `go 1.22` line escapes the module floor

**Status**: Planned
**Date**: 2026-08-27
**Queue item**: 42, `w-canary-control-does-not-survive-a-floor-raise` (clause-2, drill-surfaced,
evaluator-confirmed)
**Estimated**: ≤0.2 day (two small tests appended to an existing test file, ~55 LOC + one
stdlib import; two comment-only edits; one doc-table-row text replacement proposed verbatim
here and applied by the sprint; **no `run.sh` edit, no `ci.yml` edit, no `scripts/` edit, no
assertion in `host/store/` changed**)
**Designer**: `claude-fable-5` (design-doc-creator, iteration 130)
**Revision**: round-1 quorum BLOCKED (2 of 3 rejecting, full strength); both objections ACCEPTED
and applied 2026-08-27 — A (`gpt5-6-sol`): the floor bound is `<=`, not `<`, and the boundary is
now a measured pair M2a/M2b (V18, V19); B (`oc-glm-5-2`): the sibling's M3 verdict is measured,
not reasoned (V20), and every other reasoned-not-run claim about existing code was measured or
relabeled (V23's audit). One revision pass; one re-quorum followed.
**Round 2**: 2 of 3 reviewers PASSED and both remaining comments localised onto a single
surface — the one `racecontrol/` claim V23 relabeled rather than measured. Resolved by the
controller under the narrow-refinement carve-out, applying `gpt5-6-sol`'s own
`proposed_fix` (second branch: *"explicitly defer it without claiming the fix is
systemically complete"*) after running the experiment its first branch specified: the
prediction is **REFUTED** (V24) and is filed as queue row 48. No controller-invented
design was substituted for the reviewer's text, and no third designer run was spent.
**Toolchain boundary**: every command below was run first-party in this worktree at `2f727c7`
(clean tree; porcelain 0 re-checked after every mutation arm), shell `zsh`,
`PATH=/opt/homebrew/bin:$PATH`, darwin/arm64. The root floor is active: `go version` at the
repo root reports `go1.26.6` because `go.mod:3` auto-selects it; **inside the nested repro
module the ambient default is the Homebrew base install `go1.26.4`** — itself deny-listed —
which is why `run.sh`'s `-gcflags` control prints `BUG` at base (V9). All four controller
measurements from the iteration-130 brief were **re-run, and all four reproduce** (V3–V6). No
AILANG (`.ail`) source is written or changed by this design; the pinned `ailang` binary is not
exercised.

> **Thesis:** since `P6.T` raised `go.mod` to `go 1.26.6`, the in-module canary's known-bad
> arm is structurally unsatisfiable — `GOTOOLCHAIN=go1.26.5 go test ./host/store/` reds with
> the module-floor message before `TestToolchainCanary` compiles (V4), and will after **every**
> future floor raise, permanently. The mission's compensating control is the nested repro
> module, which escapes the floor for exactly one reason: `repro/go.mod` declares `go 1.22`
> (V8) — **and nothing anywhere binds that line** (V7). A blanket floor raise, a
> `go mod tidy -go=1.26.6`, or an IDE "update module" bumps it silently, and the rehearsal
> proves what happens next: with the repro floor at 1.26.6, **all four `KNOWN_BAD` probes
> print `SKIPPED (toolchain unavailable: go: go.mod requires go >= 1.26.6 …)`, `saw_bad`
> stays 0, and `run.sh` exits 1 via `INSTRUMENT FAILURE (or GOOD NEWS)`** — loud, but for the
> wrong reason, and only in a lane CI discards (`ci.yml:172` `continue-on-error: true`, where
> the script *already* exits 1 on every run for the darwin-only platform reason — row 44's
> double mask) (V10). The repair is one static invariant test in the established home
> (`host/verifygate/toolchain_pin_gate_test.go`, reusing its own `moduleGoFloor` and
> `shellAssignmentValues` helpers): **the repro module's `go` directive must stay at or below
> the oldest toolchain in `run.sh`'s `KNOWN_BAD` list** — the exact functional requirement: at
> equality the instrument is measured still armed (V19) — compared with
> `go/version.Compare`, never string ordering (`go1.9` sorts *above* `go1.26.0` as a string —
> measured, V11). Plus the honest re-labeling the row asks for: a fence test and a
> `POSITIVE ARM ONLY` comment on the in-module canary so nobody re-adds an arm that cannot
> run, and a corrected `AC15 MUT-CANARY-BLIND` row that names the repro instead of quoting an
> observable only the repro can produce.

## The finding in one paragraph

The known-positive control at `host/store/toolchain_canary_test.go` passes under the pinned
toolchain (`GOTOOLCHAIN=go1.26.6 go test ./host/store/ -run '^TestToolchainCanary$'` → rc=0,
V3), but its named known-bad arm is vacuous: `GOTOOLCHAIN=go1.26.5` on the same command reds
with `go: go.mod requires go >= 1.26.6 (running go 1.26.5; GOTOOLCHAIN=go1.26.5)` as the
**entire output** — no test binary is compiled, `toolchain_canary_test.go:40-42` never
executes, and the failure text never mentions the canary (V4). The red is the module floor
`P6.T` raised, not the miscompilation detector, and it recurs on every future floor raise.
The working known-bad arm lives in the nested repro module
(`design_docs/verification/w-race-gate-blindspot/repro`, its own `go 1.22` module): from that
directory `GOTOOLCHAIN=go1.26.5 go run .` prints `BUG: Field="" want "stateRoot"` while
`go1.26.6` and `go1.25.6` print `OK` — **all three with rc=0**; the reproducer PRINTS its
verdict and never exits non-zero (V5). `run.sh` already drives that arm with three anti-vacuity
floors (`saw_bad`/`saw_good`/`saw_pinned_ok`), and
`TestMiscompileInstrumentProbesPinnedToolchain` (`toolchain_pin_gate_test.go:190`) already
binds `run.sh`'s lists to the **root** floor — but nothing binds the **repro** module's own
`go` directive to anything (V6, V7). The whole compensation therefore hangs on one unprotected
line, and the rehearsal (mutate `go 1.22` → `go 1.26.6`, run everything, restore
byte-identical) shows the disarmed end-state exactly as the queue row predicted: every
`KNOWN_BAD` probe SKIPPED, `saw_bad=0`, rc=1 with the *upstream-fixed-it* message — a
wrong-reason loud failure in an unread lane (V10). The generalisation the row records: *a
control embedded in the artifact it controls inherits that artifact's constraints, so raising
a floor can silently disarm the instrument that proves the floor was needed* — and the escape
hatch (a separate module) is only as durable as the guard on its floor, which today is none.

## Premises

Each premise is one or more Verification Log rows; a claim without a row does not appear here.

- **P1 — the positive arm is real**: the in-module canary passes under the pinned toolchain,
  rc=0 (V3, reproducing the controller's measurement 1 — 0.232s here vs the controller's
  0.344s, same shape).
- **P2 — the in-module known-bad arm is structurally unsatisfiable**: `GOTOOLCHAIN=go1.26.5`
  reds on the module floor with the floor message as the entire output; the canary's Fatalf
  at `:40-42` never runs (V4, reproducing measurement 2 verbatim). This holds for every
  deny-listed toolchain and after every future floor raise.
- **P3 — the nested repro escapes the floor and discriminates**: `go1.26.5` → `BUG`,
  `go1.26.6`/`go1.25.6` → `OK`, all rc=0 — the verdict is printed, never exit-coded (V5,
  reproducing measurement 3 verbatim).
- **P4 — the mechanism exists; the map and guard do not**: `run.sh:24-26` carries
  `KNOWN_BAD="go1.26.0 go1.26.3 go1.26.4 go1.26.5"`, `KNOWN_GOOD="go1.26.6 go1.25.6
  go1.24.9"`, `PINNED="go1.26.6"` plus four fail-loud floors;
  `TestMiscompileInstrumentProbesPinnedToolchain` binds those lists to the ROOT floor via
  `moduleGoFloor` (`:35`) and `shellAssignmentValues` (`:166`) (V6). A repo-wide sweep for any
  binding on `repro/go.mod` finds exactly one mention of `repro` in executable-adjacent files —
  a comment in the canary's doc header (V7).
- **P5 — the unguarded line**: `repro/go.mod:6` is `go 1.22`, sha256 `287cc106…` (V8). Today
  `go1.22 < go1.26.0` and every deny-listed toolchain builds the reproducer (V9: base
  `run.sh` rc=0, all four BAD arms `BUG`, all three GOOD arms `OK`, all three floors fired).
- **P6 — the threat is rehearsed, not hypothesised**: with `repro/go.mod` mutated to
  `go 1.26.6`, a single `GOTOOLCHAIN=go1.26.5 go build` reds on the floor, and the full
  `run.sh` prints `SKIPPED (toolchain unavailable: go: go.mod requires go >= 1.26.6 …)` for
  **all four** KNOWN_BAD arms *and* for `go1.25.6`/`go1.24.9`, leaves `saw_bad=0`, and exits 1
  through the `INSTRUMENT FAILURE (or GOOD NEWS)` block — whose text blames toolchain
  availability or an upstream fix, neither of which is what happened. Restore was
  sha256-byte-identical, porcelain 0 (V10).
- **P7 — string ordering is the named trap**: `version.Compare("go1.9","go1.26.0")` = -1
  while the string comparison `"go1.9" < "go1.26.0"` is **false**; `go/version` is stdlib
  since 1.22, and `version.IsValid("go1.26.0x")` is false, so malformed tokens are detectable
  rather than silently misordered (V11).
- **P8 — the mislabeled record and the stale comment are current**: `w-mcp-projection.md:727`
  says "run the committed canary under deny-listed `go1.26.5` | repro prints `BUG…`" — the
  committed canary named, the repro's observable quoted (V13); the canary's own doc comment
  still says "the pinned-good Go 1.25.6 does not" (`:7`), two floor-raises stale (V12).
- **P9 — the new tests' home is measured**: both existing pin tests pass scoped at base, the
  new test names print `no tests to run` (the AC1 vacuity trap, red under the repaired form),
  a scoped run without `AILANG_BIN` passes (the lane is environment-independent), and
  vet/gofmt are clean on both touched packages (V14, V15). No function-name collisions in
  `host/verifygate` (V16).

### Design Freeze

- **One test file modified, no test file created**: the two new tests append to
  `host/verifygate/toolchain_pin_gate_test.go` — the file that already owns `run.sh`/floor
  binding — reusing `moduleGoFloor`, `shellAssignmentValues`, and `repoRoot`; one stdlib
  import added (`go/version`). No new module dependency.
- **Comment-only production edits**: `host/store/toolchain_canary_test.go` (doc comment
  replaced; **no assertion touched**) and `repro/go.mod` (one comment block above the
  directive; **the `go 1.22` line itself is not moved**).
- **One doc row replaced**: `design_docs/planned/w-mcp-projection.md:727` — replacement text
  proposed verbatim below; the sprint applies it.
- **NOT touched**: `run.sh` (row 44's surface — its runtime behaviour is not this item),
  `ci.yml`, `scripts/*`, `racecontrol/`, the root `go.mod`, any `host/store` assertion.
- Neither new test calls `requirePinned` or reads `AILANG_BIN`: static text scans must run in
  any lane (P9's measured split, V14).
- **No control derives its expectation from the value it checks** (row 43's evaluator
  refutation, binding here): the repro floor and the `KNOWN_BAD` list are independently
  authored artifacts; the test only *compares* them. Nothing generates one from the other.
- The comment fence and the comment must not collide: the new canary doc comment deliberately
  avoids the literal token `GOTOOLCHAIN`, because Test B's zero-needle counts it.

## Decision — one floor-bound invariant test, one canary fence test, comment-only production edits

The known-bad arm **stays in the nested repro module** — the queue row's "most likely" answer,
confirmed: it is the only placement that survives floor raises by construction, and `P41`
already wired `run.sh` to drive it with anti-vacuity floors. What this item adds is the
missing binding (the property that keeps the placement working) and the honest labels.

### Test A — `TestReproModuleFloorStaysBelowKnownBadToolchains(t)`

In `host/verifygate/toolchain_pin_gate_test.go`, after the existing tests:

1. `reproFloor := moduleGoFloor(t, repoRoot/design_docs/verification/w-race-gate-blindspot/
   repro/go.mod)` — the existing helper already fatals unless exactly one `go ` line exists
   and normalizes to `go1.22` form. `version.IsValid(reproFloor)` must hold, else
   `t.Fatalf("instrument failure: …")`.
2. Read `run.sh`; `shellAssignmentValues(lines, "KNOWN_BAD")` must yield exactly one
   assignment (fatal otherwise — each test stands alone; the sibling test's identical clause
   is duplicated, not shared state), whose fields are non-empty: an empty deny list leaves
   nothing to stay below and is an instrument failure here, whatever the runtime floors say.
3. Every `KNOWN_BAD` token must satisfy `version.IsValid`, else
   `t.Fatalf("instrument failure: KNOWN_BAD token %q is not a valid Go version; version.Compare
   would misorder it", tc)` — a malformed token silently compares below every valid version,
   which would corrupt the minimum.
4. Compute the oldest token with `version.Compare` (P7: never string ordering — `go1.9` vs
   `go1.22` is the measured trap, V11) and assert
   **`version.Compare(reproFloor, oldest) <= 0`** — at or below. The failure message names the
   consequence of the floor rising ABOVE the oldest deny-listed toolchain: every deny-listed
   probe SKIPs, `saw_bad` stays 0, and `run.sh` reds for the wrong reason (the V10 rehearsal,
   cited in the message).

**Why `<=` and not `<` — round-1 history**: the first draft enforced strict ordering; the
quorum rejected it (`gpt5-6-sol`, blocking, accepted): *"Enforcing spare headroom as a hard
gate violates the minimal-frozen-core axiom and can block valid future deny-list changes
without preventing an actual disarmament."* The bound the gate now enforces is exactly the
functional requirement, `reproFloor <= oldest KNOWN_BAD`: at equality the oldest deny-listed
toolchain still satisfies the module directive and still builds and runs the reproducer —
measured, not reasoned (V19: with the floor at `go 1.26.0`, `GOTOOLCHAIN=go1.26.0 go run .`
prints `BUG: Field="" want "stateRoot"`, rc=0, no SKIP), so a strict gate would red a tree
that is in fact armed — a false positive. The slack that exists today (`go1.22` vs
`go1.26.0`) is an observation about the current tree, never an enforced property: prepending
an older affected toolchain to `KNOWN_BAD` needs no repro edit for as long as the floor stays
at or below the new oldest token.

### Test B — `TestCanaryDeclaresPositiveArmOnly(t)`

Same file. A fence so the next reader does not re-add an arm that cannot run:

1. Read `repoRoot/host/store/toolchain_canary_test.go`. **Known-positive control first**:
   `strings.Count(src, "stateRoot") >= 2` (base value 3, V12) — if the canary's assertion has
   moved out of this file, `t.Fatalf("instrument failure: …")` before trusting any zero below.
2. **Zero-needle**: `strings.Count(src, "GOTOOLCHAIN") == 0` (base 0, V12) — an in-module arm
   forced onto another toolchain is the tell of exactly the dead control this row retires; the
   floor makes it structurally red (P2), so its only effect is to waste the next sprint and
   re-confuse the record. Error message points at the nested repro module.
3. **Marker needle**: `strings.Contains(src, "POSITIVE ARM ONLY")` — the required comment
   below exists (base 0, V12; this is the clause that makes AC2's red-before-green real).

### The production edit — comments and one doc row; no executable line changes

**(a) `host/store/toolchain_canary_test.go`** — replace the doc comment at `:5-7` (whose
"pinned-good Go 1.25.6" is two floor-raises stale, V12) with:

```go
// TestToolchainCanary preserves the compiler-regression shape from
// design_docs/verification/w-race-gate-blindspot/repro. Go 1.26.0 through
// 1.26.5 miscompile it on darwin/arm64; the go.mod floor (currently 1.26.6)
// does not.
//
// POSITIVE ARM ONLY: this test asserts that the ACTIVE toolchain compiles the
// shape correctly — nothing more. It cannot carry a known-bad arm: the module
// floor rejects every deny-listed toolchain before this file compiles
// (`go: go.mod requires go >= 1.26.6`), so an arm forced onto a deny-listed
// toolchain reds for the floor, not for the miscompilation — and will after
// every future floor raise, permanently. The known-bad arm lives in the nested
// `go 1.22` module at design_docs/verification/w-race-gate-blindspot/repro,
// driven by run.sh; TestReproModuleFloorStaysBelowKnownBadToolchains
// (host/verifygate) keeps that module buildable by every deny-listed
// toolchain. Do not re-add a known-bad arm here.
```

(The comment deliberately never spells the environment-variable token: Test B's zero-needle
counts it, and the fence and the comment must not collide — Design Freeze.)

**(b) `design_docs/verification/w-race-gate-blindspot/repro/go.mod`** — one comment block
above the directive, at the mutation site an IDE-bump human will actually see:

```
// The `go 1.22` line below is LOAD-BEARING: it must stay at or below the
// oldest toolchain in ../run.sh's KNOWN_BAD list, or every deny-listed probe
// SKIPs and the instrument is disarmed. Enforced by
// TestReproModuleFloorStaysBelowKnownBadToolchains (host/verifygate). Do not
// let `go mod tidy -go=…` or an IDE floor-bump touch it.
```

**(c) `design_docs/planned/w-mcp-projection.md:727`** — replace the `MUT-CANARY-BLIND` row
with the following, verbatim (the sprint applies it; base state measured at V13):

```
| AC15 `MUT-CANARY-BLIND` | run the NESTED REPRO — `design_docs/verification/w-race-gate-blindspot/repro`, its own `go 1.22` module — under deny-listed `go1.26.5`: `GOTOOLCHAIN=go1.26.5 go run .` from that directory | the repro PRINTS its verdict: `BUG: Field="" want "stateRoot"` on stdout with **rc=0** — read the output, never the exit code (re-run 2026-08-27). The committed in-module canary asserts the POSITIVE arm only: since `P6.T` raised the root floor to `go 1.26.6`, `GOTOOLCHAIN=go1.26.5 go test ./host/store/` reds with `go: go.mod requires go >= 1.26.6` before `TestToolchainCanary` compiles — a floor red, not a detector red. (*Corrected 2026-08-27, row 42: the previous text named the committed canary while quoting an observable only the repro can produce.*) |
```

### What the gate CANNOT catch — declared residual

- **The directive is not the build.** Test A proves the *declared floor* stays at or below
  the oldest deny-listed toolchain; it cannot prove a deny-listed toolchain still *builds* the reproducer on a given rig
  (network fetch, toolchain cache), nor that the miscompilation still reproduces (darwin-only;
  CI's linux runner reports `OK` on every bad arm — row 44). The runtime proof lane is
  `run.sh`'s `saw_bad` floor, which is attended/local only while `ci.yml:172` stands.
- **Source-level floor requirements are invisible.** If `repro/main.go` ever grows a language
  construct needing > go1.22, every older toolchain fails the build with a *non-floor* error
  → SKIPPED → `saw_bad=0`, while Test A stays green (the directive did not move). Declared,
  not patched: the compensating control is again `run.sh`'s own floor, in its lane.
- **Test B fences one token.** A re-added bad arm that avoids the literal `GOTOOLCHAIN` string
  (e.g. exotic env plumbing) slips the zero-needle; the module floor then reds it at runtime —
  loudly but confusingly, which is precisely the state this row documents. The fence buys
  cheap early detection of the *likely* re-add, not completeness.
- **Prose is unguarded.** The `repro/go.mod` comment and the corrected `AC15` row are text no
  test binds (a needle on a planned/ doc's table row would be vacuous ceremony); the canary
  comment is the one prose artifact given a needle (Test B), because it is the one a future
  sprint-executor reads before writing code.
- **The wrong-reason failure mode survives in degraded form**: if the invariant is ever broken
  on a tree that skips `go test` (hand-run `run.sh` only), the observer still sees the
  misleading `INSTRUMENT FAILURE (or GOOD NEWS)` text. Rewording `run.sh`'s block to name a
  floor mismatch is a `run.sh` edit — row 44's surface — deliberately not taken here.

## Alternatives rejected

1. **Move the known-bad arm back in-module** (subprocess under a pinned-bad toolchain inside
   `host/store`): structurally unsatisfiable after `P6.T` — the floor rejects the toolchain
   before any test code runs (P2, V4). That is the finding, not a fix. Rejected.
2. **Lower or carve out the root floor** so bad toolchains can build the root module: defeats
   the very protection `P6.T` landed; the floor is the defect's fix, not its bug. Rejected.
3. **Derive one side from the other** — generate `repro/go.mod`'s directive from the
   `KNOWN_BAD` minimum, or assert list membership computed *from* the floor: a control derived
   from the value it checks is vacuous by construction (row 43's evaluator refutation, cited
   by the brief as binding). Both artifacts stay independently authored; the test compares.
   Rejected as mechanism, kept as rule.
4. **String/lexicographic version comparison**: `"go1.9" < "go1.26.0"` is false as strings
   while true as versions — measured (V11). `go/version` (stdlib since 1.22) with an
   `IsValid` pre-check is the whole cost. Rejected.
5. **A runtime pre-flight in `run.sh`** (compare the repro floor to the lists before probing):
   duplicates the static gate in exactly the lane whose exit code CI discards and where the
   script already exits 1 every run for the platform reason (row 44, V10's double mask); the
   static test gates CI through `verify_go.sh`'s `go test ./...` legs. Rejected.
6. **A `toolchain go1.22` directive in `repro/go.mod`**: a hidden override the root-module
   test pattern already treats as a defect shape, and useless here — the `GOTOOLCHAIN`
   environment variable run.sh sets takes precedence, **measured** (V22: with
   `toolchain go1.24.9` appended to `repro/go.mod`, `GOTOOLCHAIN=go1.26.5 go run .` still
   prints `BUG: Field="" want "stateRoot"` — the env var won; a directive win would have
   printed `OK`). Rejected.

## Ordering

Gated on nothing. Neighbours named and not absorbed: **row 43** will publish the floor-raise
coupling inventory and should *cite* this binding as one of the couplings (enforced here,
mapped there); **row 44** owns everything about `run.sh`'s CI inertness and `ci.yml:172` —
this item changes no runtime behaviour of the instrument; **rows 45/46** untouched. The next
floor-raiser's obligation, one sentence: raising the root floor never requires touching
`repro/go.mod` — if a raise reds Test A, the correct fix is updating `KNOWN_BAD`/`KNOWN_GOOD`
per row 41's obligation (probe first, then move the lists), never bumping the repro floor to
chase it.

## Files to Create/Modify

- **MODIFY** `host/verifygate/toolchain_pin_gate_test.go` (+~55 LOC: Test A, Test B, import
  `go/version`; helpers `moduleGoFloor`/`shellAssignmentValues`/`repoRoot` reused, none
  redeclared; no name collisions at base, V16).
- **MODIFY** `host/store/toolchain_canary_test.go` — doc comment `:5-7` replaced with the
  `POSITIVE ARM ONLY` block above; **zero assertion changes**.
- **MODIFY** `design_docs/verification/w-race-gate-blindspot/repro/go.mod` — comment block
  only; the `go 1.22` directive byte-unchanged.
- **MODIFY** `design_docs/planned/w-mcp-projection.md` — the `:727` row replaced verbatim.

No other files. `run.sh`, `ci.yml`, root `go.mod`, `scripts/*`, `racecontrol/` — untouched.

## Conflict Surface

- **`TestMiscompileInstrumentProbesPinnedToolchain`** (`toolchain_pin_gate_test.go:190`) —
  reads the same `KNOWN_BAD` assignment line, binds the lists to the **root** floor; Test A
  binds the **repro** floor *at or below* the lists' oldest token. Disjoint assertions over
  one shared line, both fates **measured at design time, not reasoned**: shared fate on a
  deleted `KNOWN_BAD=` line — the sibling reds with three clauses (`:200` control needle,
  `:220` count=0, `:236` non-empty), rc=1 (V21) — and split verdicts on a malformed token —
  under M3 (`go1.26.5` → `go1.26.5x`) the sibling is GREEN: scoped `--- PASS` rc=0, and
  package-wide with `AILANG_BIN` set the run is rc=0 with 0 failures, base and M3 alike (V20,
  reproducing the controller's quorum-triage measurement). The sprint re-confirms both on the
  post-sprint tree (AC3).
- **Shared helpers** (`moduleGoFloor`, `shellAssignmentValues`, `normalizeToolchainPin`) — an
  edit that breaks a helper reds every consumer at once; that is already this file's pattern
  (both existing tests share them), accepted rather than duplicated.
- **`host/store/toolchain_canary_test.go`** — comment-only edit; `TestToolchainCanary`'s
  assertions untouched, so its green under the pinned toolchain (V3) is unaffected; AC5's
  vet/gofmt covers the edit. Test B now reads this file cross-package (read-only, by path).
- **`run.sh`** — read by two tests, written by none; row 44's remediation surface is intact.
- **`w-mcp-projection.md`** — row 43 declared this doc's file table historical; the single-row
  `:727` edit collides with nothing row 43 plans.
- **`racecontrol/`** — the repo's third module, untouched; see Systemic-Issue Audit for why it
  does not need this invariant.

## Systemic-Issue Audit

Is "a floor raise silently disarms a nested control" a pattern? Census: exactly three
`go.mod` files exist (root, `repro/`, `racecontrol/` — V7). The **repro** module is the only
artifact in the repository that is deliberately built by toolchains OLDER than the root floor
— that is its whole design — so it is the only place this invariant is meaningful.
**`racecontrol/`** (also `go 1.22`, V17) is driven by `verify_go.sh:229` under the *default*
toolchain only. Round 2 rejected the "a floor-raise bump there is harmless" prediction as an
unverified premise doing load-bearing scoping work, and the controller **measured it** rather
than re-reasoning it (V24) — **the prediction is REFUTED in general and true only by
accident of this rig**. Inside `racecontrol/` the ambient toolchain is the Homebrew base
`go1.26.4`, *not* the root-selected `go1.26.6`; bumping its directive to `go 1.26.6` makes
`GOTOOLCHAIN=auto` silently switch toolchains, which succeeds here **only because go1.26.6 is
already in the local toolchain cache**. With one variable changed — `GOTOOLCHAIN=local`, i.e.
no download available — the same bump yields `go: go.mod requires go >= 1.26.6 (running
go 1.26.4; GOTOOLCHAIN=local)` and **zero** `WARNING: DATA RACE` lines, so `verify_go.sh:232`
would FATAL with *"the race detector is not armed"*. That is the identical class this row is
about — a nested diagnostic control disarmed by a floor bump — on a **second** module.
**It is therefore DEFERRED, explicitly, and this item does NOT claim the fix is systemically
complete**: `racecontrol`'s correct binding is not the one Test A enforces (it has no
`KNOWN_BAD` list to stay below, and its real requirement — "buildable by whatever toolchain
`verify_go.sh` happens to run" — is not statically knowable), so binding it here would be
controller-invented design rather than the reviewer's fix. It is filed on its own first-party
evidence as **queue row 48** (`w-racecontrol-floor-bump-disarms-the-race-control`), per the
skill's rule that a pre-existing defect surfaced by a reviewer is a queue row, not a revision.
Mitigating and recorded: the failure there is LOUD (`verify_go.sh` FATALs) rather than silent,
which is why it is a row rather than a blocker. The **root** module *is* the floor. The
in-module canary was the one instance of a control embedded in the artifact it controls; the
sweep for other toolchain-forcing test arms in Go sources returns zero (V7's grep saw only the
canary's path comment). After this item, the one meaningful instance is bound, the one
tempting re-add site is fenced and labeled, and the two non-instances are named here with
reasons. The fix is correctly local.

## Deferred Scope

- **Row 44 `w-miscompile-instrument-inert-in-ci`** — `ci.yml:172` `continue-on-error: true`
  discards `run.sh`'s exit code, and the script currently exits 1 in CI on 10-of-10 runs
  because the miscompile is darwin-only while CI is linux/amd64. This item's V10 rehearsal
  shows the disarmed state would be indistinguishable in that lane (same rc=1, same block) —
  which is exactly why the guard here is a *static test in the gating lane*, and why nothing
  in this item touches `run.sh` or `ci.yml`. Named, not absorbed.
- **Row 48 `w-racecontrol-floor-bump-disarms-the-race-control`** (NEW, filed by this
  iteration from round 2's objection) — the second nested diagnostic module has the same
  class of exposure, measured in V24: bumping `racecontrol/go.mod` to the root floor
  disarms `verify_go.sh`'s race-detector known-positive whenever no toolchain download is
  available. Deferred rather than absorbed because its correct binding is a different
  property from Test A's, and inventing it here would be controller-authored design. This
  item explicitly does NOT claim systemic completeness.
- **Row 45** (`GOTOOLCHAIN` normalizer accepts malformed values) and **row 46** (worldd CLI
  stderr race) — out of scope, named per the brief.
- **OD-1's extension** (bind the oldest `KNOWN_GOOD` token too) — deferred with a default, see
  Open Decisions.
- Rewording `run.sh`'s `INSTRUMENT FAILURE (or GOOD NEWS)` block to distinguish a floor
  mismatch from toolchain unavailability — a `run.sh` edit; belongs with row 44's lane work if
  taken at all.

## Acceptance Criteria

Each AC carries its vacuity self-test and its **observed result on the unmodified tree at
`2f727c7`**, run this session (Verification Log rows cited).

- **AC1 — both tests exist, RUN, and pass on the post-sprint tree, in run-existence form.**
  `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run
  'TestReproModuleFloorStaysBelowKnownBadToolchains|TestCanaryDeclaresPositiveArmOnly'
  -count=1 -v` → rc=0 with exactly 2 top-level `=== RUN` lines and 2 `--- PASS`; a paired
  nonsense pattern (`-run TestNoSuchReproFloorTestZZZ`) prints `[no tests to run]`, proving
  the instrument says so rather than passing vacuously. **Base @2f727c7: the verbatim command
  → `testing: warning: no tests to run`, `ok … [no tests to run]`, rc=0 (V14)** — the naive
  form is green at base measuring nothing; the `=== RUN` clauses are the repair and red at
  base (0 of 2).
- **AC2 — test-first ordering: Test B reds before the comment edit and greens after, while
  Test A greens throughout.** Sprint evidence = two recorded runs: (i) tests applied, canary
  comment untouched → **Test B FAILS on the `POSITIVE ARM ONLY` marker needle** (base marker
  count 0 against firing same-file control `stateRoot`=3, V12) while **Test A PASSES** (the
  invariant holds at base: `go1.22` is below `go1.26.0`, well within the `<=` bound —
  measured on the prototype, V18's green control — a Test A red at this stage means the test
  misreads a conformant tree); (ii) comment edit applied → both PASS. *Not vacuous:* Test A
  cannot red-before-green on a conformant base; its teeth are AC3's mutation arms, and Test
  B's marker supplies the genuine red-first leg. **Base: premise measured (V12).**
- **AC3 — the seven named RED mutations (M1, M2a, M3–M7) red the named tests on the
  post-sprint tree, and the M2b equality control stays GREEN**, each arm restored
  byte-identically (sha256, house recipe), porcelain 0 after every arm, pristine control
  green between arms. M1 carries special status: its *runtime* consequence was rehearsed at
  base this session — all four KNOWN_BAD probes SKIPPED, `saw_bad=0`, `run.sh` rc=1 on the
  wrong-reason block, restore `287cc106…` byte-identical (V10) — so the sprint's M1 run must
  show the *static* test now reds the same edit the runtime lane only garbles. M2a and M2b
  were both RUN at design time against a verbatim prototype of Test A (V18 RED, V19 GREEN +
  armed-at-equality runtime proof); the sprint re-runs both against the landed test. The
  sibling `TestMiscompileInstrumentProbesPinnedToolchain`'s M3 verdict is **MEASURED green at
  design time (V20)**; the sprint RE-CONFIRMS it on the post-sprint tree. **Base: the bound
  clause itself is measured via the prototype (V18/V19); the threat half is measured (V10).**
- **AC4 — environment independence.** `env -u AILANG_BIN go test ./host/verifygate/ -run
  '<A>|<B>' -count=1` → rc=0 with AC1's run-existence clauses; no network, no solver, no
  pinned binary. **Base: the scoped lane is measured working without `AILANG_BIN` on the
  existing pin test — rc=0 (V14)**; a new test calling `requirePinned` would split AC4 red
  from AC1 green, exactly the defect this AC exists to catch.
- **AC5 — hygiene**: `go vet ./host/verifygate/ ./host/store/` rc=0 and `gofmt -l
  host/verifygate/ host/store/` prints nothing. **Base: both green (V15) — green-at-base by
  design**; they measure the sprint's own edits and are listed so the final-tree gate list is
  complete rather than assumed.
- **AC6 — the `AC15` row correction landed verbatim.**
  `grep -c "run the committed canary under deny-listed" design_docs/planned/w-mcp-projection.md`
  → **0** and `grep -c "the repro PRINTS its verdict" design_docs/planned/w-mcp-projection.md`
  → **1**, with the replacement row byte-equal to §(c) above. **Base @2f727c7: 1 and 0
  respectively (V13)** — both clauses red at base.
- **AC7 — the repro comment landed and the directive did not move.**
  `grep -c "LOAD-BEARING" design_docs/verification/w-race-gate-blindspot/repro/go.mod` → 1
  AND `grep -c '^go 1.22$' …/repro/go.mod` → 1. **Base: 0 and 1 (V8)** — first clause red,
  second green, so the pair proves the comment arrived *without* the directive moving.

Explicitly rejected as an AC: "the full verify gate is green" alone — `verify_ail.sh` +
`go build ./... && go test ./...` must of course pass on the final tree (and a bare
`verify_go.sh` without `AILANG_BIN` is rc=1 at base by design), but a package-wide `ok` can
print while the named tests never ran (V14's `no tests to run` is rc=0); AC1's run-existence
form is the binding version.

## Non-Vacuity — named RED mutation for every added assertion

Production side mutated (`repro/go.mod`, `run.sh`'s list line, the canary file, file
placement) — never the test helpers. Assertion coverage: A1-floor-read←M1/M2a (via
`moduleGoFloor`'s fatal on a moved file, M7 shape), A2-single-assignment/non-empty←M6,
A3-validity←M3, A4-bound-compare (`<=`)←M1/M2a with M2b as A4's equality GREEN control,
B1-positive-control←M7, B2-zero-needle←M4, B3-marker←M5.

| # | Exact edit | Expected RED (single test name) | Shape |
|---|---|---|---|
| M1 | `repro/go.mod:6` `go 1.22` → `go 1.26.6` (the blanket floor raise / `go mod tidy -go=` / IDE bump) | `TestReproModuleFloorStaysBelowKnownBadToolchains`: floor `go1.26.6` above oldest KNOWN_BAD `go1.26.0` | **threat-shaped: the rehearsed disarming (V10)** — runtime lane garbles it, static lane must name it |
| M2a | `repro/go.mod:6` `go 1.22` → `go 1.26.1` (the smallest version STRICTLY ABOVE the oldest deny-listed token) | same test: bound clause fires (`Compare > 0`) — **RUN at design time on the prototype: rc=1, `--- FAIL`, message names `go1.26.1` above `go1.26.0` (V18)** | boundary-shaped: proves `<=` is enforced at all, rather than being a tautology |
| M3 | `run.sh:24` token `go1.26.5` → `go1.26.5x` | same test: `t.Fatalf("instrument failure: KNOWN_BAD token "go1.26.5x" is not a valid Go version…")` — `IsValid` measured false for this shape (V11); sibling pin test **MEASURED green at design time (V20)**, re-confirmed by the sprint (AC3) | malformed-token-shaped: the silent mis-compare pre-empted |
| M4 | insert a line containing `GOTOOLCHAIN` into `host/store/toolchain_canary_test.go` (the tell of a re-added in-module bad arm) | `TestCanaryDeclaresPositiveArmOnly`: zero-needle 0→1 | **ADDITION**: the fence fires on the re-add's first token |
| M5 | delete the `POSITIVE ARM ONLY` comment block from the canary file | `TestCanaryDeclaresPositiveArmOnly`: marker needle absent | removal: the label this item exists to add, guarded |
| M6 | delete the `KNOWN_BAD=` line from `run.sh` | `TestReproModuleFloorStaysBelowKnownBadToolchains`: assignment count 0≠1 fatal (shared fate: the sibling pin test reds too — **measured, V21**: rc=1, three clauses fire; both consume the line) | empty-instrument-shaped |
| M7 | `git mv host/store/toolchain_canary_test.go host/store/toolchain_canary_moved_test.go` | `TestCanaryDeclaresPositiveArmOnly`: `os.ReadFile` fatal — the fence names its target by path, and a moved canary must move the fence in the same edit | placement-shaped: keeps the positive control honest |

**M2b — the equality GREEN control (not a mutation arm; M2a's boundary pair).**
`repro/go.mod:6` `go 1.22` → `go 1.26.0`, exactly equal to the oldest `KNOWN_BAD` token.
Two halves, BOTH RUN at design time (V19) and re-run by the sprint: (i) the gate —
`TestReproModuleFloorStaysBelowKnownBadToolchains` stays **GREEN** (rc=0, `--- PASS` on the
prototype); (ii) the armed-at-equality runtime proof, which makes the green evidence rather
than assertion — from `design_docs/verification/w-race-gate-blindspot/repro`,
`GOTOOLCHAIN=go1.26.0 go run .` → rc=0, output `BUG: Field="" want "stateRoot"` — the
equality floor still builds and prints a verdict, no SKIP. Had (ii) failed, `<=` would be
the wrong bound and the strict form would return with that measurement as its justification;
it did not fail. Restore `287cc106…` byte-identical, porcelain clean, per arm.

Green control for all arms: the unmutated post-sprint tree passes AC1/AC4, and every arm ends
restored sha256-byte-identical with `git status --porcelain` empty — the recipe V10 already
ran once at base (`287cc106…` before and after).

The two unguarded artifacts are declared, not mutated: the `repro/go.mod` comment (AC7's grep
is the check; no test binds prose there) and the `AC15` row text (AC6's greps are the check).
Neither is an assertion; the residual section carries both.

## Open Decisions

- **OD-1 — should Test A also bind the oldest `KNOWN_GOOD` token** (today `go1.24.9`) at or
  above the repro floor (the same `<=` bound, mirrored)? *Controller default if no one
  answers: NO.* The pinned arm is already bound
  to the root floor by the sibling test; a lost historic-good arm degrades to SKIPPED with
  `saw_good` still satisfiable by the pin, and a total good-arm loss trips `saw_good=0` at
  runtime. The KNOWN_BAD side is the one with no other guard anywhere — that asymmetry is the
  item. (If taken later, it is a three-line extension of Test A, same mechanism.)
- **OD-2 — should Test A pin `KNOWN_BAD`'s exact token set?** *Default: NO*, for row 41's
  OD-2 reason unchanged: the instrument's list must stay free to evolve (upstream fixes,
  new affected patches); non-empty + all-valid + at-or-above-the-repro-floor is the honest
  bound at this row's scope.

## Verification Log

All rows run first-party by the designer at `2f727c7` (clean tree), shell `zsh`,
`PATH=/opt/homebrew/bin:$PATH`, darwin/arm64, 2026-08-27. V1–V17 during authoring; V18–V23
during the round-1 revision pass, same HEAD, same rig. The controller's four brief
measurements are marked (C1–C4) and **all four reproduced**; V20 additionally reproduces the
controller's quorum-triage M3 measurement (C5). Where an incidental reading differed
(timing), it is recorded. KP = known-positive control carried in the same call. V18/V19 ran
the Decision's Test A verbatim as a temporary prototype file inside `host/verifygate`
(reusing the package's own `moduleGoFloor`/`shellAssignmentValues`/`repoRoot`), deleted
before any package-wide row and before the porcelain checks — the sprint re-runs both arms
against the landed test.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Worktree is `2f727c7`, clean; every mutation arm below ended restored | `git rev-parse HEAD && git status --porcelain \| wc -l` (re-run after V10) | `2f727c7ed8ca…`, `0`; `0` after the arm |
| V2 | Toolchain boundary: root floor auto-selects the pin | `go version` at repo root | `go1.26.6 darwin/arm64` (inside `repro/`, the ambient default is `go1.26.4` — see V9's banner) |
| V3 (C1) | Positive control: in-module canary passes under the pin | `GOTOOLCHAIN=go1.26.6 go test ./host/store/ -run '^TestToolchainCanary$' -count=1` | rc=0, `ok … host/store 0.232s` (controller: 0.344s — same shape, timing differs); `-v` form: 1 `=== RUN`, 1 `--- PASS` |
| V4 (C2) | THE FINDING: known-bad arm reds on the floor, canary never executes | `GOTOOLCHAIN=go1.26.5 go test ./host/store/ -run '^TestToolchainCanary$' -count=1` | rc=1; **entire output**: `go: go.mod requires go >= 1.26.6 (running go 1.26.5; GOTOOLCHAIN=go1.26.5)` — no test binary, no mention of the canary. Reproduced verbatim |
| V5 (C3) | The nested repro escapes the floor and prints (never exit-codes) its verdict | from `repro/`: `GOTOOLCHAIN=<tc> go run .` for go1.26.5 / go1.26.6 / go1.25.6 | `go1.26.5 rc=0 out=BUG: Field="" want "stateRoot"`; `go1.26.6 rc=0 out=OK`; `go1.25.6 rc=0 out=OK` — rc=0 all three. Reproduced verbatim |
| V6 (C4) | The mechanism exists; the binding to the ROOT floor exists; nothing binds the repro floor | read `run.sh` in full, numbered; read `toolchain_pin_gate_test.go` in full | `KNOWN_BAD="go1.26.0 go1.26.3 go1.26.4 go1.26.5"` `:24`, `KNOWN_GOOD="go1.26.6 go1.25.6 go1.24.9"` `:25`, `PINNED="go1.26.6"` `:26`; four fail-loud floors (`ran` `:80`, `saw_bad` `:84`, `saw_good` `:90`, `saw_pinned_ok` `:95`); `TestMiscompileInstrumentProbesPinnedToolchain` at `:190` binds the lists to `moduleGoFloor(t, go.mod)` (`:239`); helpers `moduleGoFloor` `:35`, `shellAssignmentValues` `:166`. **No clause anywhere reads `repro/go.mod`** |
| V7 | Nothing in any executable-adjacent file binds or reads `repro/go.mod`; module census | `git grep -n "w-race-gate-blindspot/repro\|repro/go.mod" -- '*.go' '*.sh' '*.yml'`; `find . -name go.mod -not -path './.git/*'` | exactly one hit: `host/store/toolchain_canary_test.go:6` — a path in a comment. Three modules: root, `repro/`, `racecontrol/` |
| V8 | The unguarded line, verbatim, with its hash | `cat -n repro/go.mod`; `shasum -a 256 repro/go.mod` | `:4 module ailang-world/verification/go1_26-arraylit-miscompile`, `:6 go 1.22` (header comment: root must not pick it up); sha256 `287cc1066b1d…`. No `LOAD-BEARING` comment (AC7 base: 0) |
| V9 | Base positive control: the instrument works end-to-end today | `./design_docs/verification/w-race-gate-blindspot/run.sh` at base | rc=0; banner `host: darwin/arm64 default toolchain: go1.26.4`; all four BAD arms `BUG: Field="" want "stateRoot"`, all three GOOD arms `OK`; `RESULT: reproduction confirmed` + all three floor lines incl. `pinned toolchain (go1.26.6) reported OK` |
| V10 | **THE REHEARSAL: a floor raise on `repro/go.mod` disarms every KNOWN_BAD probe and reds `run.sh` for the wrong reason** | `cp` backup; `sed` `go 1.22`→`go 1.26.6`; single probe `GOTOOLCHAIN=go1.26.5 go build` in `repro/`; full `run.sh`; `cp` restore; sha256 + porcelain | probe rc=1 `go: go.mod requires go >= 1.26.6 (running go 1.26.5; …)`; full run rc=1: **all four KNOWN_BAD `SKIPPED (toolchain unavailable: go: go.mod requires go >= 1.26.6 …)`**, `go1.25.6`/`go1.24.9` also SKIPPED, only `go1.26.6` ran (`OK`), `-gcflags` control flips to `OK`/`OK`, then `INSTRUMENT FAILURE (or GOOD NEWS): no known-affected toolchain reproduced the defect…` — blames availability/upstream, not the floor raise that caused it. Restore byte-identical `287cc106…`, porcelain 0 |
| V11 | The version-comparison trap and the malformed-token detector, measured | standalone `go run`: `version.Compare("go1.9","go1.26.0")`, string `<`, `IsValid("go1.26.0x")`, `IsValid("go1.22")`, `Compare("go1.22","go1.26.0")` | `-1`; string `<` **false** (lexicographic misorders); IsValid(`go1.26.0x`) **false**; IsValid(`go1.22`) true; `-1`. `go/version` is stdlib — no new dependency |
| V12 | Canary-file needle base state and the stale comment | `grep -c` per needle on `host/store/toolchain_canary_test.go`; read file | `GOTOOLCHAIN` **0** (B2 base-green), `stateRoot` **3** (KP fires), `POSITIVE ARM ONLY` **0** (B3 base-RED — AC2's lever); `:7` still says "the pinned-good Go 1.25.6 does not" (two floor-raises stale); Fatalf at `:40-42` |
| V13 | The mislabeled AC15 row, located and counted | `grep -n "MUT-CANARY-BLIND" design_docs/planned/w-mcp-projection.md`; `grep -c` both needles | row at `:727`: "run the committed canary under deny-listed `go1.26.5` \| repro prints `BUG…`"; old-needle count **1**, new-phrase (`the repro PRINTS its verdict`) count **0** — AC6 reds at base both ways |
| V14 | Gate-lane base state: existing tests green scoped; the AC1 vacuity trap; the no-AILANG_BIN lane | scoped `go test -v` on the two existing pin tests; `-run '^TestReproModuleFloorStaysBelowKnownBadToolchains$'`; `env -u AILANG_BIN go test -run '^TestMiscompileInstrumentProbesPinnedToolchain$'` | 2× `--- PASS`, rc=0; new name → `testing: warning: no tests to run`, `ok … [no tests to run]`, **rc=0** (bare "command greens" is vacuous at base); scoped no-AILANG_BIN run rc=0 (static lane is env-independent) |
| V15 | Hygiene baselines on both touched packages | `GOTOOLCHAIN=go1.26.6 go vet ./host/verifygate/ ./host/store/`; `gofmt -l host/verifygate/ host/store/ \| wc -l` | rc=0; `0` |
| V16 | No test-name collision for the two new tests | `grep -h "^func " host/verifygate/*_test.go \| grep -i "repro\|floor\|canary"` | only `moduleGoFloor` (reused helper); neither new name exists |
| V17 | `racecontrol/` is the non-instance (Systemic audit) | `cat racecontrol/go.mod`; `git grep -n "racecontrol" -- '*.go' '*.sh' '*.yml'` | `go 1.22`, deliberate-race header; driven only by `scripts/verify_go.sh:229` `go run -race .` under the **default** toolchain. ~~so a floor bump there cannot disarm anything~~ **— that trailing clause is REFUTED by measurement in round 2; see V24.** |
| V18 | **M2a RUN: the `<=` bound is enforced, not a tautology** | prototype Test A in place; green control on pristine tree first; then `sed` `go 1.22`→`go 1.26.1` in `repro/go.mod`; `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestReproModuleFloorStaysBelowKnownBadToolchains$' -count=1 -v` | green control: rc=0, `--- PASS`. M2a arm: **rc=1, `--- FAIL`**, message: `repro module floor "go1.26.1" is above the oldest KNOWN_BAD toolchain "go1.26.0": every deny-listed probe SKIPs, saw_bad stays 0, and run.sh reds for the wrong reason (the V10 rehearsal)`. Restore `287cc106…` byte-identical, porcelain clean |
| V19 | **M2b RUN: GREEN at equality, and the instrument is genuinely armed there** | `sed` → `go 1.26.0` (equal to oldest KNOWN_BAD); same scoped test; then from `repro/`: `GOTOOLCHAIN=go1.26.0 go run .` | test: **rc=0, `--- PASS`**. Runtime: **rc=0, output `BUG: Field="" want "stateRoot"`** — builds and prints a verdict at equality, no SKIP; `<=` is the right bound (a SKIP here would have restored `<` with this row as justification). Restore `287cc106…`, porcelain clean |
| V20 (C5) | **Objection B's measurement: the sibling is GREEN under M3** — measured, not reasoned | pristine scoped control; `sed` `run.sh:24` `go1.26.5`→`go1.26.5x` (landed-proof `grep -c 'go1.26.5x'` → 1); scoped `go test -run '^TestMiscompileInstrumentProbesPinnedToolchain$' -count=1 -v`; package-wide `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -count=1` base and under M3 | pristine scoped: rc=0, `--- PASS`. **M3 scoped: rc=0, `--- PASS`.** Package-wide base: **rc=0, 0 `--- FAIL` lines** (`ok … 47.699s`); package-wide under M3: **rc=0, 0 `--- FAIL` lines** (`ok … 46.434s`). Reproduces the controller's quorum-triage measurement, incl. its attribution note: the same package-wide command WITHOUT `AILANG_BIN` was rc=1 with 17 failures on the RESTORED tree too — the documented `AILANG_BIN`-unset base condition, not M3. Restore `c9e2916c…`, porcelain clean |
| V21 | M6's shared-fate half measured: the sibling reds on a deleted `KNOWN_BAD=` line | delete the `KNOWN_BAD=` line from `run.sh`; scoped sibling run as in V20 | **rc=1, `--- FAIL`**, three clauses fire: `:200` `does not contain known-positive control "KNOWN_BAD="`, `:220` `KNOWN_BAD assignment count=0, want 1`, `:236` `KNOWN_BAD must contain at least one toolchain`. Restore `c9e2916c…`, porcelain clean |
| V22 | Alternatives №6's precedence claim measured: `GOTOOLCHAIN` beats a `toolchain` directive | append `toolchain go1.24.9` to `repro/go.mod`; from `repro/`: `GOTOOLCHAIN=go1.26.5 go run .` | rc=0, output `BUG: Field="" want "stateRoot"` — go1.26.5 ran (a directive win would print `OK` via go1.24.9). Restore `287cc106…`, porcelain clean. Bonus grounding: `verify_go.sh:256` carries the `go test ./... -count=1` leg cited in Alternatives №5 |
| V23 | The two revision audits, auditable | (a) strictness sweep: `grep -n -i 'strict' design_docs/planned/w-canary-control-does-not-survive-a-floor-raise.md`; (b) reasoned-not-run audit rule: every claim of a specific runtime verdict of EXISTING executable code under a specific condition must be run or relabeled | (a) pre-revision **11** line hits, all restating the strict bound; post-revision **4**, none stating the enforced bound: two in the round-1-history paragraph, one in M2b's contingency sentence, one describing M2a's mutation value ("smallest version STRICTLY ABOVE"). (b) **4 hits**: sibling-green-under-M3 → measured (V20); sibling-red-under-M6 → measured (V21); `GOTOOLCHAIN`-beats-directive → measured (V22); `racecontrol` bump-harmless → relabeled as a prediction with its reason, in Systemic-Issue Audit (mutating a third module + full `verify_go.sh` is outside this row's mutation budget; the claim is scoping rationale, not an enforced property) |
| V24 (round-2 objection, controller-measured) | **`racecontrol/`'s floor-bump-harmless prediction is REFUTED: the bump disarms the race-detector known-positive control whenever no toolchain download is available** | `cp` backup + sha256; base arm = the exact `verify_go.sh:229` invocation (`cd .../racecontrol && go run -race .`) plus `go version` / `go env GOTOOLCHAIN,GOVERSION` **from inside that module**; then `sed` `go 1.22`→`go 1.26.6` and re-run; then the one-variable arm `GOTOOLCHAIN=local`; `cp` restore + sha256 + porcelain | **base**: inside the module `go version` = **`go1.26.4`** (Homebrew base, *not* the root-selected go1.26.6), `GOTOOLCHAIN=auto`, `GOVERSION=go1.26.4`; `go run -race .` rc=1 with **2** `WARNING: DATA RACE` (the gate's known-positive fires). **bumped to `go 1.26.6`**: `go version` inside the module becomes **`go1.26.6`** — auto silently switched toolchains — and the control still fires (rc=1, **2** DATA RACE) **only because go1.26.6 is already in the local toolchain cache**. **bumped + `GOTOOLCHAIN=local`** (one variable): rc=1, output is exactly `go: go.mod requires go >= 1.26.6 (running go 1.26.4; GOTOOLCHAIN=local)`, **0** DATA RACE lines → `verify_go.sh:232` would FATAL *"the race detector is not armed"*. Restored sha256 `ab782f11db0f7f259f73dd55a58eaf5a30b871bb79bd98bacbe964d50efc025b` byte-identical; post-restore control re-fires (rc=1, 2 DATA RACE); porcelain shows only this untracked doc. Census re-confirmed complete: `find . -name go.mod -not -path './.git/*'` → exactly 3 |

## Related Documents

- [`../implemented/w-setup-go-pin-unguarded.md`](../implemented/w-setup-go-pin-unguarded.md) —
  row 41 (`P41`, PR #97 → `8e3c8cd`): built the file this item extends, the helpers it
  reuses, and the root-floor binding this item completes on the repro side; its OD-2
  (do not pin KNOWN_BAD's exact set) is inherited here unchanged.
- [`w-mcp-projection.md`](w-mcp-projection.md) — `P6.T` raised the floor that created this
  finding; its `AC15 MUT-CANARY-BLIND` row at `:727` is corrected by this item (§Production
  edit (c)).
- `design_docs/world-mission.md` queue rows 42 (this item), 43 (floor-raise coupling
  inventory — will map the binding enforced here), 44 (`run.sh` CI inertness — the lane
  caveat this doc declares rather than fixes), 45, 46.
- `design_docs/verification/w-race-gate-blindspot/` — the instrument (`run.sh`, untouched),
  the repro module (one comment added, directive untouched), and `racecontrol/` (the audited
  non-instance).
- [`../implemented/w-race-gate-blindspot.md`](../implemented/w-race-gate-blindspot.md) — where
  the canary and the nested-module pattern were born, including the original
  `MUT-CANARY-BLIND` refutation history this item's correction extends.

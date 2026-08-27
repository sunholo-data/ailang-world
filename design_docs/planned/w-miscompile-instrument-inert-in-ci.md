# w-miscompile-instrument-inert-in-ci — on linux the instrument can only ever refuse, its second floor fires before any assertion runs, and the workflow is built to discard refusals

**Status**: Planned
**Date**: 2026-08-27
**Revision**: 2 (2026-08-27) — quorum round 2 BLOCKED 3/3 on ONE shared objection (M2's
pin redded against M1's own comment bytes); closed by the ratified narrow-refinement
carve-out applying `gpt5-6-sol`'s proposed_fix VERBATIM, with the demanded verification
row measured first-party as **V19** (two arms + a scoping control + three RED mutations).
See §Quorum round 2.
**Revision**: 1 (2026-08-27) — quorum round 1 BLOCKED at full strength (3/3 present,
absent_reviewers EMPTY); disposition: R1's fail-closed platform `case`, complete-coverage
floor and linux/amd64-scoped prose applied verbatim (proof rows P17/P18 added); R2's
`[ "$saw_bad" -ne 0 ]` applied verbatim; R3 closed by premise row P16 **plus** the
mechanism itself — kernel-read `uname` probe and a wiring-test pin (M7) — not just the
premise. Details and pricing in §Quorum round 1.
**Queue item**: 44, `w-miscompile-instrument-inert-in-ci` (clause-2, quorum-surfaced at
row 41 round 2, controller-measured, queued iter-128)
**Estimated**: ~0.35 day — three files:
`design_docs/verification/w-race-gate-blindspot/run.sh` (+~45/−8 executable lines: one
kernel-read platform-polarity block with a fail-closed refuse arm, one complete-coverage
floor; darwin behaviour byte-stable on the measured profile — the banner prints the
identical string (P16) and P8's 7/7-ran profile passes the new floor),
`.github/workflows/ci.yml` (−4/+7: delete the one `continue-on-error`, rename the step,
replace its 2-line comment with the platform contract),
`host/verifygate/toolchain_pin_gate_test.go` (+~45: one new test function carrying two
pins — gating and override-channel; ~3 comment lines reworded inside an existing
doc-comment, **zero assertion bodies touched**).
Zero other files. No `.ail` touched; single-gate (Go) change.
**Designer**: `kimi-k3:cloud` (design-doc-creator, iteration 133)
**Toolchain boundary**: every command below ran first-party in this worktree
(`.wt-world-iter133`, bash, darwin/arm64, macOS 26.5.2) at `80c2bd2`, tree clean,
2026-08-27. GitHub reads used `gh` against `sunholo-data/ailang-world` (visibility PUBLIC).
At repo root `go` resolves to `go1.26.6` via the `go.mod` floor (`GOTOOLCHAIN=auto`);
the rig's base toolchain is affected — inside the nested repro module `run.sh`'s banner
printed `default toolchain: go1.26.4`, matching `OUTPUTS.md` §5 — so the attended darwin
run below is exactly the developer-rig case. The attended `run.sh` fetched per-probe
toolchains over the network; nothing it wrote entered the tree.

> **Thesis:** `run.sh` asks one question — *does a KNOWN-BAD toolchain miscompile the
> reproducer HERE?* — with a single global expectation, and the measured answer on
> linux/amd64 is permanently "no", so its second refusal floor is the only floor CI has
> ever reached and the two assertion floors behind it have never evaluated (P7). The
> existing `ci.yml` comment says the refusal "records evidence"; the sibling doc already
> falsified that theory of recording: *a step log nothing requires anyone to read is not
> a record* (P12), and a swallowed refusal on every observed dev run is the
> measurement, not a metaphor (P4, P6). The honest repair is not choosing louder silence:
> make the expectation platform-conditional over the two MEASURED platforms —
> darwin/arm64 must SEE the defect, linux/amd64 must NOT see it, and any third platform
> must REFUSE rather than inherit a polarity (P18) — with the platform itself read from
> the kernel, not from an overridable variable (P16) — add a floor that refuses unless
> EVERY configured known-bad probe ran, and then delete `continue-on-error: true` so the
> exit code carries its own meaning again, with a static test salting the ground the
> flag stood on. The row's own
> generalisation is the acceptance bar: `continue-on-error: true` converts an
> instrument's loudest possible output into silence, so the more carefully an instrument
> fails, the more completely the flag hides it — the fix must make careful failure
> *binding*, not merely audible.

## The finding in one paragraph

The mission's first-party miscompile instrument has exited 1 on every observed dev CI
run while GitHub reported the step `success` — verified by me on the current dev HEAD
run `33069999131` (`80c2bd2`): the step conclusion is `success`, the same log carries
`INSTRUMENT FAILURE (or GOOD NEWS)` exactly once, `RESULT: reproduction confirmed` zero
times, and its own same-log control banner once (P3, P4). The cause is not malfunction:
the defect is measured darwin/arm64-only (`scripts/verify_go.sh`'s deny-list comment,
P9; the iteration-46 AC6 linux/amd64 measurement, P12; my own attended darwin run
reproduces all four BUGs while CI's linux run shows seven OKs, P5/P8 — a two-arm
platform control on identical source with one variable), so on linux `saw_bad` can never
honestly become 1, and the floor order means the linux leg asserts *nothing at all*: it
would exit 1 identically if `repro/` stopped compiling, if the pinned toolchain
regressed, or if the probe lists were hollowed out (P7). `continue-on-error: true` at
the step then discards the refusal — which, per V-H's declared residual (P11), is the
named follow-up this doc discharges. The design question the row poses is what the
instrument should *assert* on a platform where the defect does not reproduce; the answer
below is: assert the measured platform truth — "no KNOWN-BAD toolchain may reproduce
HERE", with every configured probe completing — bind it to floors that are reached on
linux/amd64, refuse on unmeasured platforms, and gate it.

## Premises

`VERIFIED-BY-ME` = run by the designer in this worktree/at this rig today, output
observed first-hand. `INHERITED` = controller-measured, cited, and where load-bearing
re-derived here (noted on the row). Controller rows V-A…V-I as filed.

| # | Claim | Command (verbatim) | Observed | Status |
|---|---|---|---|---|
| P1 | `continue-on-error` occurs exactly once in the repo's only workflow, at the miscompile step (V-A) | `grep -n 'continue-on-error' .github/workflows/ci.yml`; `grep -c '^      - name:' .github/workflows/ci.yml`; `ls .github/workflows/` | `170:        continue-on-error: true`; `15`; `ci.yml` | VERIFIED-BY-ME (V-A re-derived: one hit, 15 named steps as enumerated control, single workflow so the enumeration is complete) |
| P2 | The step block's shape and its pre-existing comment (V-G) | `sed -n '167,172p' .github/workflows/ci.yml` | comment `# run.sh exits 1 when no affected toolchain reproduces, which is the good` / `# "linux/amd64 is clean" outcome; this records evidence and is not a gate.`, then `- name: Measure compiler reproducer on linux/amd64 (non-gating)`, `continue-on-error: true`, `timeout-minutes: 15`, `run: ./design_docs/verification/w-race-gate-blindspot/run.sh` | VERIFIED-BY-ME (V-G comment byte-matches the controller's quote) |
| P3 | Current dev HEAD run: step reports `success` to the API | `gh api repos/:owner/:repo/actions/runs/33069999131/jobs --jq '.jobs[] \| .steps[] \| select(.name \| test("reproducer";"i")) \| {name,conclusion}'` plus `/runs/33069999131` | step conclusion `success`; run `head_sha=80c2bd2452eb…`, `head_branch=dev`, `event=push`, conclusion `success`; step #8, 12:06:38→12:07:25Z (~47 s, the flake-relevant duration) | VERIFIED-BY-ME (V-B first half re-derived) |
| P4 | Same run's log: refusal present, banner absent, control fires | `gh run view 33069999131 --log` then `grep -c` on three needles | `INSTRUMENT FAILURE (or GOOD NEWS)` = 1; `RESULT: reproduction confirmed` = 0; `go1.26 local-array-literal miscompilation reproduction` = 1; also `##[error]Process completed with exit code 1.` = 1 | VERIFIED-BY-ME (V-B second half: GitHub's green and the script's refusal coexist on one run) |
| P5 | The linux probe table (V-C) | same log, lines after `expected BUG (affected):` | `go1.26.0/1.26.3/1.26.4/1.26.5 expect=BAD got: OK (rc=0)`; `go1.26.6 + go1.25.6 + go1.24.9 expect=GOOD got: OK (rc=0)` — ran=7, saw_bad=0, saw_good=1, saw_pinned_ok=1 | VERIFIED-BY-ME (byte-match to V-C) |
| P6 | Streak + fetch-flake measurement | per run in `33069999131 33069076984 33042694782 33042059824 33027555185 33007276859`: `gh run view $run --log` + 5 greps | every run: failure-banner=1, RESULT=0, control=1, probes=7, SKIPPED=0 → 6/6 runs, 42/42 toolchain probes fetched | VERIFIED-BY-ME for 6/6; the 10/10 streak is VERIFIED BY CONTROLLER (iters 128–130), INHERITED here and consistent with this window. **Meta-finding**: my first sweep reported 0/0/0 because `gh` ran outside a git repo and `2>/dev/null` hid rc=1 — the impossible 0 on the known-positive control caught it. Recorded because it is this row's class: an empty grep is a claim, and the control is what makes it a measurement. |
| P7 | Floor order makes the linux leg vacuous today (V-D — "the fact the row does not state") | `awk 'NR>=79 && NR<=104'` on `design_docs/verification/w-race-gate-blindspot/run.sh`; P4's log sequence | floors at :80 `ran==0`, :84 `saw_bad==0`, :90 `saw_good==0`, :95 `saw_pinned_ok==0`, each `exit 1`; the linux log ends with the saw_bad-floor text — floors 3–4 never evaluated. So today linux asserts nothing: repro bit-rot, a pinned regression, or a vanished pin all exit 1 identically, already discarded | VERIFIED-BY-ME (read + log order corroborates) |
| P8 | Darwin arm of the platform control | `./design_docs/verification/w-race-gate-blindspot/run.sh` (attended, this rig) | rc=0; all four KNOWN-BAD print `BUG: Field="" want "stateRoot"`; go1.26.6/1.25.6/1.24.9 `OK`; `-N` → OK, `-l` → BUG; banner `RESULT: reproduction confirmed, and both controls fired.` + 3 detail lines | VERIFIED-BY-ME — identical source, one variable (GOOS/GOARCH) against P5 |
| P9 | The deny-set is platform-scoped in its own comment; floors of the two modules | `grep -n 'deny-list is the measured set' scripts/verify_go.sh`; `head -5 go.mod`; `cat design_docs/verification/w-race-gate-blindspot/repro/go.mod` | `214:# This deny-list is the measured set: go1.26.0-go1.26.5 on darwin/arm64.`; root floor `go 1.26.6`; repro floor `go 1.22` with its load-bearing comment | VERIFIED-BY-ME (row's "verify_go.sh:214" confirmed verbatim) |
| P10 | Coupled static tests: binds and base state (V-E) | read `host/verifygate/toolchain_pin_gate_test.go:14-311`; `go test ./host/verifygate/ -run '…4 names…' -count=1 -v` | `TestMiscompileInstrumentProbesPinnedToolchain` requires: `KNOWN_BAD=`/`KNOWN_GOOD=`/`PINNED=` each column-0-assigned exactly once; exactly one `#!/usr/bin/env bash`; KNOWN_GOOD ∋ go.mod floor; PINNED == floor; PINNED ∈ KNOWN_GOOD; floor ∉ KNOWN_BAD; exec bit; `saw_pinned_ok` count ≥ 3; literal `INSTRUMENT FAILURE: the PINNED toolchain`. Plus `TestReproModuleFloorStaysBelowKnownBadToolchains` binds repro floor ≤ oldest KNOWN_BAD. All four scoped tests `--- PASS`, `ok … 0.125s` | VERIFIED-BY-ME (every bind re-read from source, not transcribed from the row) |
| P11 | The instrument's own declared residual — the follow-up this doc IS (V-H) | `sed -n '183,191p' host/verifygate/toolchain_pin_gate_test.go` | verbatim: *"What remains open by scope: nothing WATCHES the non-gating log — a loud failure nobody reads is loud only on inspection; flipping ci.yml:172 to gating is the named follow-up in Deferred Scope, paired with OD-1."*; and `w-setup-go-pin-unguarded.md` §Deferred Scope verbatim: gating it is *"a follow-up, not this item … gating a network-dependent instrument turns every transient toolchain-fetch failure into a CI red, which wants a flake-rate measurement of its own; queue it paired with OD-1"* | VERIFIED-BY-ME (read); the demanded flake measurement is P6: 42/42 fetches, 0 SKIPs |
| P12 | The linux/amd64 question is already answered, and the sibling doc already falsified "the log is a record" (V-I) | `sed -n '204,218p' design_docs/implemented/w-race-gate-blindspot.md`; `sed -n '405,420p' bench/BASELINE.md` | §Deferred Scope: *"ANSWERED 2026-08-04 (iteration 46, AC6): linux/amd64 is NOT affected"* (run `30872300227`, all four affected toolchains OK), with the honesty note that this cuts *against* that doc's motivation; `:217` verbatim: *"a step log nothing requires anyone to read is not a record."* — BASELINE.md records the historical linux table under the historical step name | VERIFIED-BY-ME |
| P13 | The `ci.yml:172` citation census has grown past the controller's nine (V-F, corrected) | `grep -c 'ci\.yml:172'` per file repo-wide; `grep -rc 'ci\.yml:170'` | 25 lines across 8 files: `w-setup-go-pin-unguarded.md` 14, `w-canary-control-does-not-survive-a-floor-raise.md` 4, `toolchain_pin_gate_test.go` 2, the two sprint plans 1+1, `world-mission.md` (row 44 itself!) 1, log 1, status-archive 1; `ci.yml:170` cited by 0 files | VERIFIED-BY-ME — controller's "nine" measured before rows 43/45-51 docs landed; the class grew in a week, which is row 31's thesis instantiated in this premise |
| P14 | Pricing premise for a darwin arm; new test name is free | `gh repo view --json visibility --jq .visibility`; `go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1` | `PUBLIC`; `ok … [no tests to run]` rc=0 (the vacuity base for AC1) | VERIFIED-BY-ME |
| P15 | A verifygate test red IS a CI-gate red | `grep -n 'go test' scripts/verify_go.sh` | `:258` `go test ./... -count=1`, with the `-race` leg following at :262-266, run from the repo root inside the gated `go-verify` step | VERIFIED-BY-ME — needed for the mutation venues below |
| P16 | The polarity switch's input has NO override present in-tree today (quorum R3's premise), the override mechanism is REAL, and the kernel-read replacement is immune to it | `grep -n 'GOOS\|GOARCH' .github/workflows/ci.yml; echo rc=$?`; `grep -c 'GOTOOLCHAIN' .github/workflows/ci.yml`; `grep -n 'GOOS\|GOARCH' design_docs/verification/w-race-gate-blindspot/run.sh`; `GOOS=windows GOARCH=amd64 go env GOOS GOARCH`; `GOOS=windows GOARCH=amd64 uname -s`; `GOOS=windows GOARCH=amd64 uname -m` | ci.yml: **zero hits** (grep rc=1), with the same call's `GOTOOLCHAIN` control printing `2` — the grep is live; run.sh: exactly one hit, `:61` `echo "host: $(go env GOOS)/$(go env GOARCH)   default toolchain: …"` — a display-only READ, quoted in full, no assignment anywhere; mechanism measured real: under `GOOS=windows GOARCH=amd64`, `go env` prints `windows`/`amd64` while `uname` still prints `Darwin`/`arm64`; and the uname-derived pair equals the `go env` pair on this rig (`MATCHES go env on this rig`, both `darwin/arm64`) | VERIFIED-BY-ME — R3's premise is TRUE TODAY and was UNPINNED; the revision pins it: the probe below reads the kernel, the :61 read is replaced by `$host_pair` (same printed string on both verified platforms), and the resulting zero counts in both files are asserted by the wiring test (mutation M7 proves the pin's teeth) |
| P17 | Constraint 1 under R1's fail-closed `*)` arm: every CI host maps into the listed `linux/amd64` arm, so the refuse arm cannot red dev | `grep -n 'runs-on' .github/workflows/ci.yml`; `gh run view 33069999131 --log` then greps for `Image:`, `Linux-x64`, `Linux/x86_64`, `linux/amd64` | both jobs `runs-on: ubuntu-latest` (`:19`, `:100`) — no ARM or macOS label exists anywhere in the workflow; the run's log carries `Image: ubuntu-24.04` and THREE independent in-run platform reports agreeing the machine is Linux on x64: setup-go's cache key `setup-go-Linux-x64-ubuntu24-go-1.26.6-…`, `✓ compiler pinned by exact bytes: AILANG v0.30.0 on Linux/x86_64`, and `go1.26.6 linux/amd64` ×7 → `uname -s`/`uname -m` there resolve `linux`/`x86_64` → host_pair `linux/amd64` → the listed arm; replaying P5's measured profile (ran=7, ran_bad=4 of bad_expected=4, saw_bad=0, saw_good=1, saw_pinned_ok=1) through the revised floors passes all of them → exit 0 | VERIFIED-BY-ME — a non-ubuntu-x64 runner arrives only via a deliberate, visible `runs-on` edit, and then the refuse arm's red is loud and NAMED, not silent |
| P18 | The verified platform-contract set is EXACTLY two pairs — R1's demanded proof row for the fail-closed case | `sed -n '214p' scripts/verify_go.sh`; `grep -n 'NOT affected' design_docs/implemented/w-race-gate-blindspot.md` | `# This deny-list is the measured set: go1.26.0-go1.26.5 on darwin/arm64.`; and `:207` *"…(iteration 46, AC6): linux/amd64 is NOT affected.\*\* All four affected toolchains return `OK` on …"* plus `:435` the AC6 row `MET`. No premise in this doc or its siblings measures any other host pair — darwin/amd64, linux/arm64, windows/* are all unverified | VERIFIED-BY-ME — hence the `case` lists exactly two arms and refuses everything else, instead of asserting a cleanliness no measurement supports |

## Options considered

**E0 — flip `continue-on-error: true` off, nothing else.** *Cost:* one line. *Catches:*
nothing — it only moves the refusal from the log to the check-status column. *What reds
it:* the next dev push — measured, not argued: on linux `saw_bad` is permanently 0 (P5,
P8, P12), floor 2 fires, dev goes red with no recourse. Constraint 1 kills it on
arrival; enumerated only because the row's headline could be misread as this.

**A — platform-conditional expectation in `run.sh`, then gate.** The expectation is a
function of the MEASURED platform facts, not a constant: on `darwin/arm64` a KNOWN-BAD
probe MUST report BUG (unchanged semantics); on `linux/amd64` NONE may (a BUG there
means the defect escaped its measured set on the one platform this repo's CI builds
on); every OTHER platform refuses — the verified contract set is exactly two pairs
(P18) — and the platform itself is read from the kernel because `go env` honours the
env-var form of the platform tokens (P16). One new floor — *every configured KNOWN-BAD
probe actually ran* (`ran_bad` vs `bad_expected`, counted at probe entry so it tracks
the list) — sits between `ran` and the polarity floors, so a partial SKIP on the
negative arm is a refusal, not a quieter certification (quorum R1's second half; priced
in §Quorum round 1). The existing `saw_good`/`saw_pinned_ok` floors — today unreachable
on linux (P7) — become the linux leg's positive half, and the flip is safe: today's
linux profile (ran=7, ran_bad=4 of 4, saw_bad=0, saw_good=1, saw_pinned_ok=1, P5)
passes every floor and exits 0 (P17); P8's darwin profile passes unchanged. *Cost:*
~0.25 d; the run.sh tail rewrites. *Catches:* repro compile-rot, disarmament via ANY
known-bad fetch failure, defect-spread to linux/amd64 (an assertion no artefact
currently makes), an unverified platform drifting into CI, pinned regression and
vanished-pin (newly reachable on linux). *Misses:* any darwin-side change (no darwin
runner; Declared residuals). *What reds it:* M1–M7; green-on-arrival proven by AC4(a)
and AC6 against P5/P8.

**B — move/add the arm to a darwin CI runner.** Full two-sided detection in CI.
*Cost, priced honestly:* the repo is PUBLIC (P14), so `macos` runners carry no cash
multiplier today (the 10× minute multiplier binds only if it ever goes private). The
real prices: macOS hosted capacity queues measurably longer than ubuntu, making the
detection the job's long pole; `TestGoToolchainPinsAgreeAndMatchJobList` asserts the
workflow's job set is exactly `[ailang-verify, go-verify]` and pins every
GOTOOLCHAIN/go-version/setup-go count to that cardinality — a third job is a coupled
edit to a passing gate; and it answers a question different from this row's — it
certifies the platform CI does **not** build on, leaving the linux leg asserting
nothing. *What reds it:* an upstream fix (good red) and macos image/queue churn (bad
red). Ruled out as **this** item: strictly larger than the fix the row needs, and the
darwin detection surface is already served attended (P8) plus deny-list/canary
(P9, P10). Compatible with A; named in Deferred Scope as a human-priced decision.

**C — linux leg asserts only the known-good half.** Delete/quiet the `saw_bad` floors on
linux; keep `ran`, `saw_good`, `saw_pinned_ok`. *Cost:* ~5 lines fewer than A. *Misses,
fatally:* the defect-spread event. Either KNOWN-BAD probes still run and their BUG is
required to *not* fail (loud text, green exit — the exact silence class this row exists
to kill), or they stop running and the step stops measuring the thing it is named for.
Both are A minus five lines minus the only new assertion linux needs. C ⊂ A; the delta
is the point.

**D — leave it non-gating, add a WATCHER that reads the log.** *Cost:* a new
authenticated API consumer with its own flake domain, grepping log text that is not an
API (P6's meta-finding: even an `rc`-level failure is invisible without a control; a
text grep has two more silent-failure shapes). *Catches:* what it greps for, *after*
the merge. *Misses:* bindingness — dev stays green on the offending push by
construction. And V-H (P11) already named this option's place: the watcher is what you
build when you *cannot* flip the flag; a watcher with the platform truth in it is A
re-implemented in a worse layer with weaker teeth.

## Decision

**Option A:** teach `run.sh` the measured platform asymmetry (fail-closed on
unverified platforms), add the complete-coverage floor on the known-bad arm, then make
the step gate — with a static wiring test so neither the flag nor the override channel
can quietly return.

Why the others died: E0 violates constraint 1 (measured, P5+P7: next push reds). C is A
minus the only new assertion linux needs; wherever the KNOWN-BAD probes still run, C
re-introduces the swallow-one-event shape (P4's class) one level down. B certifies the
wrong platform for CI's purposes, couples a third job into a cardinality-pinned gate,
and is priced for a decision the mission has not said it wants — it is a *detection*
upgrade for a platform already covered attended, while this row is a *binding* defect.
D is the row's own generalisation re-litigated: an instrument's careful exit code is
already the loudest signal it owns; a watcher is a second instrument whose silence modes
are strictly worse than the first's.

A alone satisfies constraint 2 without a framework: every added assertion has a named
mutation that reds the shipped gate (M1–M7), and the linux leg stops being vacuous —
after this change it asserts four things it asserts none of today: EVERY configured
known-bad probe ran to completion, none of them reproduced HERE, the reproducer still
passes under known-good toolchains, and the pin was probed OK.

**Coupled-code disposition (V-E, each named):**

- `TestMiscompileInstrumentProbesPinnedToolchain` (`toolchain_pin_gate_test.go`): **zero
  assertion changes.** The run.sh edit keeps every bound token: `KNOWN_BAD=`/
  `KNOWN_GOOD=`/`PINNED=` single column-0 assignments byte-untouched; single shebang;
  `saw_pinned_ok` sites stay 3 (declaration :31, probe-set :54, guard :95); the literal
  `INSTRUMENT FAILURE: the PINNED toolchain` stays byte-identical in its guard. The
  only touch is inside the test's own doc-comment: the two stale `ci.yml:172` lines and
  the now-discharged residual sentence (P11) are reworded — comment bytes, no code.
- `TestReproModuleFloorStaysBelowKnownBadToolchains`: untouched; the `KNOWN_BAD`
  assignment line it parses is untouched; repro/go.mod untouched.
- `host/store/toolchain_canary_test.go` and `TestCanaryDeclaresPositiveArmOnly`:
  untouched (positive arm under the pinned toolchain; unchanged behaviour on linux).
- `scripts/verify_go.sh`: untouched. Its deny-list (:214) stays the darwin-scoped,
  attended-rig detector; this design changes nothing it asserts.
- **Scope-out (rows 45, 48, 49, 50, 51):** this design edits none of their surfaces —
  not `normalizeToolchainPin` (45), not `racecontrol/` (48), not the canary fence or
  `toolchain_canary_test.go` itself (49), not `shellAssignmentValues` (50 — the new
  TOP-LEVEL assignments stay at column 0 in deference to it; the `case`-internal
  `host_os=`/`host_arch=`/`expect_defect=` arms are necessarily indented, like every
  other conditional body in the file, and invisible to a column-0 reader), not the
  inventory test (51).
  File-level overlap exists only for row 45 (same test file); see M2 conflict note.

## Milestones

**M1 — `run.sh` learns the platform truth (still non-gating; zero CI risk).**
One file: `design_docs/verification/w-race-gate-blindspot/run.sh`.

- Declarations block (after `ran=0` at :32); top-level assignments kept column-0
  (row-50 deference). Two quorum-driven properties are baked in here: the platform is
  read from the KERNEL, never `go env`, because `go env GOOS` honours a `$GOOS`
  environment variable (measured at P16 — R3); and the polarity `case` is FAIL-CLOSED,
  because the verified contract set is exactly two pairs (P18 — R1). The mechanism fix
  is R3's own second half (an un-overridable probe, normalised to Go's naming); the
  polarities and refuse-arm text are R1's proposed fix verbatim:

```bash
ran_bad=0       # KNOWN-BAD probes that actually built and ran (row 44)
bad_expected=0  # KNOWN-BAD probes CONFIGURED, counted at probe entry — tracks the list

# Platform probe (row 44; quorum round-1 R3): read the kernel, never `go env` —
# `go env GOOS` honours a $GOOS environment variable (measured: design doc P16), so an
# overridable variable must not be the polarity's sole input. `uname` has no such
# override channel. Normalized to Go's naming so the deny-set comparisons stay in one
# vocabulary; unknown values fall into the refuse arm below.
case "$(uname -s)" in
	Darwin) host_os=darwin ;;
	Linux)  host_os=linux ;;
	*)      host_os=unknown ;;
esac
case "$(uname -m)" in
	arm64|aarch64) host_arch=arm64 ;;
	x86_64|amd64)  host_arch=amd64 ;;
	*)             host_arch=unknown ;;
esac
host_pair="$host_os/$host_arch"

# Platform expectation (row 44; quorum round-1 R1): the verified contract set is
# EXACTLY two pairs (design doc P18). darwin/arm64: the deny-list's measured set — a
# KNOWN-BAD probe MUST report BUG or this script can no longer see the defect.
# linux/amd64: measured NOT affected (iteration-46 AC6) — no KNOWN-BAD probe may report
# BUG on the one platform CI builds on. Any other host has no verified contract and
# refuses; extend the measurement set before trusting it, do not inherit a polarity.
case "$host_pair" in
	darwin/arm64) expect_defect=1 ;;
	linux/amd64) expect_defect=0 ;;
	*) echo "INSTRUMENT FAILURE: no verified platform contract for $host_pair"; exit 1 ;;
esac
```

- In `probe()`, two lines: at entry,
  `[ "$expect" = "BAD" ] && bad_expected=$((bad_expected + 1))` — a SKIP still counts as
  expected, which is the point; and beside `ran=$((ran + 1))`,
  `[ "$expect" = "BAD" ] && ran_bad=$((ran_bad + 1))`. (`set -e` is absent — verified
  `:20` is `set -uo pipefail` — so the `&&` one-liner form is safe.)
- The banner read at :61 swaps to the already-computed pair — printed string identical
  on both verified platforms (P16) — and the file's last read of the overridable
  channel leaves with it:
  `echo "host: $host_pair   default toolchain: $(go version | awk '{print $3}')"`.
- Floor block replaces :84-89 (the unconditional `saw_bad` floor) and the banner at
  :101-104 is split per platform; floors :80-83 and :90-100 stay byte-identical. The
  disarmament floor is now COMPLETE COVERAGE (R1's second half: a logged partial SKIP
  must not yield certification), the hollowed-list case (`bad_expected==0`) keeps the
  old floor's teeth on the Deferred-Scope narrowing shape, and the clean-platform alarm
  tests the invariant `-ne 0` rather than exact match (R2, verbatim):

```bash
if [ "$bad_expected" -eq 0 ] || [ "$ran_bad" -ne "$bad_expected" ]; then
	echo "INSTRUMENT FAILURE: $ran_bad of $bad_expected KNOWN-BAD probes completed —"
	echo "a partial (or hollowed-out) negative arm cannot certify $host_pair. Every"
	echo "configured KNOWN-BAD entry must run; a SKIP on this arm is a refusal, not a pass."
	exit 1
fi
if [ "$expect_defect" -eq 1 ] && [ "$saw_bad" -eq 0 ]; then
	echo "INSTRUMENT FAILURE (or GOOD NEWS): no known-affected toolchain reproduced the"
	echo "defect. Either the toolchains were unavailable, or upstream fixed it — in which"
	echo "case re-derive the pin decision in design_docs/ rather than trusting this pass."
	exit 1
fi
if [ "$expect_defect" -eq 0 ] && [ "$saw_bad" -ne 0 ]; then
	echo "INSTRUMENT FAILURE (PLATFORM ALARM): a KNOWN-BAD toolchain reported BUG on"
	echo "$host_pair — measured NOT affected (iteration-46 AC6). The defect escaped its"
	echo "measured set: treat as a toolchain incident and re-derive the pin."
	exit 1
fi
# … saw_good floor and saw_pinned_ok floor unchanged …
if [ "$expect_defect" -eq 1 ]; then
	echo "RESULT: reproduction confirmed, and both controls fired."
	# (three existing detail lines unchanged)
else
	echo "RESULT: $host_pair clean — no KNOWN-BAD toolchain reproduced here, matching"
	echo "the iteration-46 AC6 measurement; all $bad_expected known-bad and $ran total"
	echo "probes ran, and the known-good and pinned ($PINNED) toolchains both reported OK."
fi
```

  The darwin RESULT banner stays byte-identical so `OUTPUTS.md` §5's transcript remains
  truthful; the linux banner is deliberately a *different* sentence, so the mission's
  historical greps (`RESULT: reproduction confirmed` = 0) stay meaningful for
  pre-landing runs. Darwin-behaviour honesty for this revision: on the measured profile
  (P8, 7/7 probes ran) exit code, banner text, and floor texts are byte-stable; the ONE
  attended-rig change is that a KNOWN-BAD SKIP on darwin now refuses at the coverage
  floor instead of the `(or GOOD NEWS)` floor — one loud refusal for another, and zero
  SKIPs exist in any observed run (P6).
- Rehearsals (AC3/AC4/AC5) + hygiene (AC7) prove it before anything gates on it.

**M2 — make the refusal binding, and salt the flag's ground.**
Two files.

- `.github/workflows/ci.yml`: delete the `continue-on-error: true` line; rename the step
  to `Measure compiler reproducer (platform-conditional, gated)`; replace the two
  comment lines (P2) with the contract, citing the landed doc by its implemented path:

```yaml
      # Platform-conditional and GATED (row 44): darwin/arm64 must reproduce the defect;
      # linux/amd64 (here) must NOT (measured, iteration-46 AC6); every configured
      # known-bad probe must run to completion and both known-good controls must fire.
      # Any other host refuses — the verified contract set is exactly two pairs.
      # exit 1 is a refusal again, not evidence nobody reads.
      # Design: design_docs/implemented/w-miscompile-instrument-inert-in-ci.md
      - name: Measure compiler reproducer (platform-conditional, gated)
        timeout-minutes: 15
        run: ./design_docs/verification/w-race-gate-blindspot/run.sh
```

  (`timeout-minutes: 15` kept; the historical step-name quotes at
  `w-race-gate-blindspot.md:209` and `bench/BASELINE.md:412` stay — they are true of the
  historical runs they cite.)
- `host/verifygate/toolchain_pin_gate_test.go`: append one test, and reword the ~3
  comment lines named in the disposition above. The test carries TWO pins: the gating
  pin (its original remit) and the override-channel pin (quorum R3's closure — the
  premise row P16 records today; this pin stops tomorrow from undoing the mechanism).
  Verbatim sketch:

```go
// TestMiscompileInstrumentStepIsGatedInCI pins the row-44 wiring on two channels that
// must not silently return. (1) `continue-on-error: true` converts an instrument's
// loudest possible output into silence, so its occurrence count in the workflow is
// asserted ZERO — re-introducing it anywhere here needs its own row. (2) The
// instrument's platform polarity reads the KERNEL (`uname`); `go env` honours the
// env-var form of the platform tokens (measured in the design doc, P16), so NEITHER
// the workflow nor run.sh may mention either token — asserted zero in both files —
// and `uname -s` is asserted present in run.sh so the kernel probe cannot quietly
// revert to the overridable channel.
// DECLARED RESIDUAL: a step-level `if:` that never evaluates true disables the step
// with this text intact (no actionlint runs in this repo — P41 V18); and these are
// byte-substring pins — a computed assignment (`eval`, string concatenation) evades
// them; the mechanism's runtime immunity is why that gap is acceptable (design doc,
// residuals 2 and 3).
func TestMiscompileInstrumentStepIsGatedInCI(t *testing.T) {
	workflowPath := filepath.Join(repoRoot, ".github", "workflows", "ci.yml")
	raw, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	if count := strings.Count(src, miscompileReproducerPath); count != 1 {
		t.Fatalf("instrument failure: ci.yml count(%q)=%d, want exactly 1 — the step this test pins must exist", miscompileReproducerPath, count)
	}
	// Scope the flag check to the miscompile STEP's own block (quorum round-2 R1:
	// "inspect only the miscompile step for continue-on-error … rather than banning
	// legitimate GOOS/GOARCH use globally"). A flag on an unrelated step is that
	// step's business; V19's scoping control proves this boundary holds.
	lines := strings.Split(src, "\n")
	start := -1
	for i, l := range lines {
		if strings.Contains(l, miscompileReproducerPath) {
			for j := i; j >= 0; j-- {
				if strings.HasPrefix(strings.TrimSpace(lines[j]), "- name:") {
					start = j
					break
				}
			}
			break
		}
	}
	if start < 0 {
		t.Fatalf("instrument failure: could not locate the miscompile step block in ci.yml")
	}
	end := len(lines)
	for j := start + 1; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "- name:") {
			end = j
			break
		}
	}
	for i := start; i < end; i++ {
		if strings.Contains(lines[i], "continue-on-error") {
			t.Errorf("ci.yml:%d re-introduces %q in the miscompile step — row 44: a swallowed refusal is how this instrument died the first time", i+1, "continue-on-error")
		}
	}
	runRaw, err := os.ReadFile(filepath.Join(repoRoot, miscompileReproducerPath))
	if err != nil {
		t.Fatal(err)
	}
	runSrc := string(runRaw)
	// Require the kernel read (quorum R3), and reject only EXECUTABLE uses of the
	// overridable channel — a comment may NAME the channel it exists to warn about.
	// The round-1 form banned the bare tokens repo-wide and therefore redded against
	// its own documentation on arrival (V19 arm A); this is round-2 R1's fix verbatim.
	for _, need := range []string{"uname -s", "uname -m"} {
		if !strings.Contains(runSrc, need) {
			t.Errorf("run.sh no longer reads the kernel via %q — quorum R3's fix reverted; the polarity must not come from an overridable variable", need)
		}
	}
	for i, line := range strings.Split(runSrc, "\n") {
		code := line
		if idx := strings.Index(code, "#"); idx >= 0 {
			code = code[:idx]
		}
		for _, bad := range []string{"go env GOOS", "go env GOARCH"} {
			if strings.Contains(code, bad) {
				t.Errorf("run.sh:%d executable use of %q — the polarity must not read an overridable channel", i+1, bad)
			}
		}
		if strings.Contains(code, "host_pair") && strings.Contains(code, "go env") {
			t.Errorf("run.sh:%d derives host_pair from `go env` — quorum R3", i+1)
		}
	}
}
```

  with `const miscompileReproducerPath = "design_docs/verification/w-race-gate-blindspot/run.sh"`
  beside it. (`repoRoot` is a package-level var of package `verifygate`, defined at
  `ail_binary_gate_test.go:27` and already used by the sibling test's ci.yml read;
  `os` and `strings` are imported at the file head — P10's read.) Post-landing the pin holds against the
  shipped bytes — proven, not argued, at **V19**: the exact test above run against the
  exact M1 `run.sh` and M2 `ci.yml` gives `=== RUN`=1 / `--- PASS`=1 / rc=0, while the
  round-1 form reds on `run.sh:37` against the identical bytes.
  Conflict note: row 45's item edits `normalizeToolchainPin` near the top of this same
  file; this appends at the end — disjoint hunks, flagged for whichever lands second.
- Whole-gate green + the merge observation (AC6) close it.

## Acceptance criteria

Each carries its base observation on the unmodified tree at `80c2bd2`, measured this
session. No criterion asserts a bare line count that can rot; counts here are 0/1
predicates over live files, each with a same-call known-positive control.

- **AC1 — the wiring test exists and RUNS (run-existence form).**
  `go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1 -v`
  → rc=0 with exactly one `=== RUN   TestMiscompileInstrumentStepIsGatedInCI` and one
  `--- PASS`; paired nonsense pattern `-run '^TestNoSuchWiringGateZZZ$'` prints
  `[no tests to run]`. **Base:** the verbatim command → `ok … [no tests to run]`, rc=0
  (P14) — the naive form is green at base measuring nothing; the `=== RUN` clause is the
  binding form.
- **AC2 — the flag is gone FROM THIS STEP and the step remains (property, not line
  number; scoped per quorum round-2 R1, which forbids a repo-wide ban).** In one
  invocation: `awk` the miscompile step's own block (from its `- name:` to the next
  `- name:`) and require `grep -c 'continue-on-error'` over that block → `0`, **and**
  `grep -c 'design_docs/verification/w-race-gate-blindspot/run.sh'
  .github/workflows/ci.yml` → `1` as the same-call control proving the greps are live
  and the step still exists. A `continue-on-error` on any OTHER step is out of scope and
  must NOT fail this criterion — that boundary is measured, not assumed (V19's scoping
  control: the flag planted on `worldd benchmark smoke gate` leaves the pin GREEN).
  **Base:** the block-scoped count is `1` and the path control is `1` (P1/P2).
  Teeth: M1.
- **AC3 — the coupled static tests are untouched and green.** The two function BODIES are
  byte-stable against base:
  `diff <(git show 80c2bd2:host/verifygate/toolchain_pin_gate_test.go | awk '/^func TestMiscompileInstrumentProbesPinnedToolchain/,/^}/') <(awk '/^func TestMiscompileInstrumentProbesPinnedToolchain/,/^}/' host/verifygate/toolchain_pin_gate_test.go)`
  → empty, and the same for `TestReproModuleFloorStaysBelowKnownBadToolchains`; plus
  `go test ./host/verifygate/ -run 'TestMiscompileInstrumentProbesPinnedToolchain|TestReproModuleFloorStaysBelowKnownBadToolchains|TestGoToolchainPinsAgreeAndMatchJobList|TestCanaryDeclaresPositiveArmOnly' -count=1` → rc=0.
  **Base:** diff empty by identity; the four tests `--- PASS` at 0.125 s (P10).
- **AC4 — darwin behavioural legs, attended (this rig's class).**
  (a) `./design_docs/verification/w-race-gate-blindspot/run.sh` → rc=0 and the darwin
  banner present. (b) Guard-trip: mutate `if [ "$host_pair" = "darwin/arm64" ]` to
  `"darwin/amd64"` (every platform in clean-mode), run → rc=1 whose output names the
  PLATFORM ALARM text; restore, assert byte-identity by sha256, re-run → rc=0. **Base for
  (a):** rc=0 with 4×`BUG: Field="" want ` + `"stateRoot"` (P8). This is row 41's AC6
  guard-trip pattern applied to the new polarity floor.
- **AC5 — the instrument floors still fire, platform-independently (rehearsal).**
  `mv design_docs/verification/w-race-gate-blindspot/repro/main.go{,.MUT}` → run.sh rc=1
  via the `no toolchain ran at all` floor; restore byte-identical. This is the V-D
  scenario ("stops compiling") turning from invisible-exit-1 into a named refusal.
  **Base:** today the same mutation exits 1 with floor 1's message — and CI discards
  it identically to the platform-clean refusal, which is the defect (P4); after M2 the
  gated step turns that refusal into a red, and the linux-clean profile exits 0 with
  the new banner instead (AC6).
- **AC6 — constraint 1 discharged on the real gate (observed, not argued).**
  The landing PR's check run AND the first dev run after merge: step conclusion
  `success` while its log carries the linux RESULT sentence exactly once,
  `INSTRUMENT FAILURE` zero times, and the control banner once —
  `gh run view <run> --log | grep -c 'RESULT: linux/amd64 clean'` = 1 with the two
  side-condition greps in the same call. **Base:** on run `33069999131` the equivalents
  measure step `success` WITH refusal=1, banner=0, control=1 (P3/P4) — the precise
  contradiction this item removes; a green landing reading proves the gate took.
- **AC7 — hygiene.** `bash -n design_docs/verification/w-race-gate-blindspot/run.sh` →
  rc=0; `stat -f '%Sp' design_docs/verification/w-race-gate-blindspot/run.sh` →
  `-rwxr-xr-x`; `gofmt -l host/verifygate/` → empty; `go vet ./host/verifygate/` → rc=0;
  `go test ./host/verifygate/ -count=1` → rc=0. **Base:** exec bit verified by `ls`
  (`-rwxr-xr-x`, P10 reading); scoped tests pass (P10); full-package run happens at
  sprint time (it runs the verify_ail shim arms; ~60 s).

## Named RED mutations

Venue matters: `verify_go.sh` runs `go test ./... -count=1` inside the gated `go-verify`
job (P15), so a verifygate RED and a step-exit RED are both CI-gate reds. Each arm is
proved LANDED before its result is read (a rc=0 on a mutation that never applied is a
false green), executed against the shipped files, then restored byte-identical by sha256
(house recipe), porcelain re-checked after each arm.

| # | Mutation | File | What it neuters | Predicted result | Landed-proof before reading the result |
|---|---|---|---|---|---|
| M1 | Insert `        continue-on-error: true` under the renamed step | `.github/workflows/ci.yml` | the row-44 flip itself | `TestMiscompileInstrumentStepIsGatedInCI` RED in the verify_go leg with the offending line number; AC2's grep → 1 | `grep -n 'continue-on-error' .github/workflows/ci.yml` prints the inserted line (count 0→1) |
| M2 | Re-label the `case` arm `darwin/arm64) expect_defect=1 ;;` → `linux/amd64) expect_defect=1 ;;` (defect expected on the wrong platform; the real defect platform left unlisted) | `run.sh` | the platform-polarity discrimination | step RED on linux CI: `saw_bad=0` with `expect_defect=1` fires the `(or GOOD NEWS)` floor — AND attended darwin RED via the fail-closed `*)` arm (`no verified platform contract for darwin/arm64`, rehearsed as AC4(c)): a polarity lie now fails on BOTH platforms instead of riding the swallow on one | `grep -n 'linux/amd64) expect_defect=1' run.sh` shows the flipped arm in place |
| M3 | `mv repro/main.go repro/main.go.MUT` | `repro/` | every probe's ability to run | step RED via floor 1 (`ran==0`) — the event that today exits 1 invisibly | `cd repro && GOTOOLCHAIN=go1.26.6 go build .` → rc≠0 ("no Go files") shown and quoted |
| M4 | repro/go.mod `go 1.22` → `go 1.27` | `repro/go.mod` | EVERY probe's ability to build (an explicit `GOTOOLCHAIN=<v>` refuses a floor above it) | TWO-layer red: `TestReproModuleFloorStaysBelowKnownBadToolchains` reds statically (1.27 > oldest KNOWN_BAD 1.26.0), and step-side every probe SKIPs → floor 1 (`ran==0`) reds. Predicted-floor corrected in revision 1 BY MEASUREMENT (the old row said the `saw_pinned_ok` floor): /tmp copy of `repro/` with floor `go 1.27` under the probe's own invocation `GOTOOLCHAIN=go1.26.6 go build .` → `go: go.mod requires go >= 1.27 (running go 1.26.6; GOTOOLCHAIN=go1.26.6)`, rc=1 — for any listed toolchain | `grep '^go ' design_docs/verification/w-race-gate-blindspot/repro/go.mod` prints `go 1.27` |
| M5 | `chmod -x run.sh` | `run.sh` | direct CI invocation | step exit 126 RED (no flag left to swallow it), and the existing exec-bit assertion reds in the verify_go leg | `stat -f '%Sp' design_docs/verification/w-race-gate-blindspot/run.sh` → `-rw-r--r--` |
| M6 | Append ` go1.99.99` to KNOWN_BAD (a configured known-bad entry that cannot fetch) | `run.sh` | the assumption that configured coverage == completed coverage | step RED via the coverage floor (`4 of 5 KNOWN-BAD probes completed … refusal`); the static layer stays GREEN by design — assignment-count, pin-membership, and oldest-floor binds (P10) are all unaffected — because these teeth are behavioural. Base MEASURED in a /tmp copy of `repro/`: `GOTOOLCHAIN=go1.99.99 go build .` → `go: download go1.99.99 for darwin/arm64: toolchain not available`, rc=1, i.e. the probe prints SKIPPED exactly as the floor requires | `grep -n '^KNOWN_BAD=' run.sh` shows the appended `go1.99.99` |
| M7 | Revert the polarity's input: `host_pair="$host_os/$host_arch"` → `host_pair="$(go env GOOS)/$(go env GOARCH)"` | `run.sh` | R3's mechanism fix — puts the overridable channel back on the runtime path | `TestMiscompileInstrumentStepIsGatedInCI` RED in the verify_go leg. **MEASURED at V19 (M7a), not predicted:** `rc=1, === RUN=1, --- PASS=0, --- FAIL=1`, `run.sh:51 executable use of "go env GOOS"` and the same for `GOARCH` | `grep -c 'host_pair="\$(go env GOOS)' run.sh` → 1 (count 0→1) |
| M8 | Neuter the kernel read: `case "$(uname -s)" in` → `case "$(echo Darwin)" in` | `run.sh` | the un-overridable probe itself, while leaving `host_pair` shaped correctly | `TestMiscompileInstrumentStepIsGatedInCI` RED. **MEASURED at V19 (M7b):** `rc=1, === RUN=1, --- PASS=0, --- FAIL=1`, `run.sh no longer reads the kernel via "uname -s"`. Note the behavioural layer alone would MISS this on a darwin rig — the mutant still resolves `darwin` there — which is why the static pin is the load-bearing arm | `grep -c 'uname -s' run.sh` → 0 (count 1→0) |

Green control for all arms: the unmutated post-sprint tree passes AC1–AC7. M3–M6 are
rehearseable on this darwin rig with semantics identical to the linux runner (M6's base
needed network and is measured above); M2's attended-side rehearsal is AC4(c) and the
linux-side prediction is profile-replayed at P17. **M1, M7 and M8 are not predictions at
all — all three were EXECUTED at V19 against the exact proposed bytes**, each asserted
LANDED by count before its result was read and restored byte-identical by sha256.

**What the round-2 revision deliberately STOPPED claiming teeth over.** Revision 1's M7
mutated `env: GOOS: linux` into ci.yml and predicted a RED. Under the scoped pin that
mutation is **GREEN, correctly** — after R3's mechanism fix the polarity reads the
kernel, so a `GOOS` in the environment cannot move the step's behaviour (P16), and
quorum round-2 R1 was right that freezing the token repo-wide buys nothing while
banning legitimate use. The reversion risk R3 actually named is a change to the
POLARITY'S INPUT, and that is what M7/M8 now mutate. This is rule 3i applied to our own
table: a "kills which mutation" cell that no longer predicts the right result is worse
than no cell, and the fix is to re-derive the cell, not to keep the design that made it
true.

## Deferred Scope

- **A darwin CI arm (Option B)** — priced in Options, not filed as a row by this doc.
  A is compatible with it; whether CI should *detect* on darwin (vs attend locally) is a
  cost decision above this item's pay grade. If ever taken, it couples into
  `TestGoToolchainPinsAgreeAndMatchJobList`'s job-set cardinality (P15-adjacent; named
  so the next designer is not surprised).
- **OD-1 (dynamic installed-version assertion in the `go version` step)** — stays open,
  now *unpaired*: the pairing existed because the flip's blocker was an unmeasured flake
  rate (`w-setup-go-pin-unguarded.md` §Deferred Scope, quoted at P11); the measurement is
  now on record (P6: 42/42 probes, 0 SKIPs), the flip discharges here, and OD-1's
  assertion is orthogonal to it.
- **The `ci.yml:172` citation census (P13)** — **row 31's**, explicitly. 25 lines across
  8 files, and *this item's own ci.yml edit renumbers the file again* — refreshing
  `:172`→`:170` by hand would mint the next stale numbers (row 43's P7: line numbers
  rot; literals do not). The durable fix is row 31's resolver, not a piecemeal sweep;
  the only exceptions are the two code-comment lines inside a function whose doc-comment
  this milestone must touch anyway (they would otherwise become *false residual
  statements*, not merely stale numbers). Historical step-name quotes at
  `w-race-gate-blindspot.md:209` and `bench/BASELINE.md:412` stand as true history.
- **KNOWN_BAD membership length is unbound** — deleting three of four list members passes
  every static test here (assignment count and oldest-floor binding survive). A latent
  narrowing, adjacent to but distinct from row 50 (assignment *shape* vs list
  *membership*); named, scoped out, and partly covered here behaviourally by the
  coverage floor: a hollowed (`bad_expected==0`) or any incompletely-run list now reds.
- **Partial KNOWN-BAD SKIPs on the negative arm** — DISCHARGED in revision 1 (quorum
  R1, second half: "a logged partial SKIP must not yield certification"). The coverage
  floor makes `ran_bad != bad_expected` — ANY partial or hollowed known-bad coverage —
  a refusal, so certification can never again ride on probes that did not run. The
  price this buys (one transient fetch SKIP now reds dev) is measured and accepted in
  §Quorum round 1 and Declared residual 6, not merely asserted. KNOWN_GOOD-side
  partial SKIPs remain declared (residual 4) and were unobjected.

## Declared residuals

1. **CI remains darwin-blind by construction.** No darwin runner exists, so a *new*
   darwin-only miscompile — this shape in a post-1.26.6 toolchain, or any other shape —
   is invisible to this gate. The detectors for that surface remain the attended `run.sh`
   (P8), the darwin-scoped deny-list (P9), and the pinned-toolchain canary (P10) — and
   rows 45/48/49 sharpen exactly those seams. This item changes that surface not at all;
   it says so rather than implying coverage.
2. **The wiring test is a text assertion over YAML.** It cannot see a step-level `if:`
   that never evaluates true, a YAML anchor/flow-style that hides the step's body from
   line greps, or `timeout-minutes: 0` games — the same class P41's
   `TestGoToolchainPinsAgreeAndMatchJobList` declares in its own doc-comment (no
   actionlint in this repo, its V18 — re-read here in the same file, P10). run.sh's
   behavioural floors are the second layer: a step that never runs
   leaves no linux RESULT sentence, which AC6's read would catch on inspection — and
   inspection is once again the weak link, honestly named. Revision 1 added the test a
   second pin (zero occurrences of the two platform-token strings in ci.yml AND
   run.sh, plus positive `uname -s` presence); it inherits this residual exactly — a
   computed or obfuscated assignment (`eval`, `"GO""OS"`) would evade the substring
   pin, and no behavioural floor backstops that channel (after R3's fix the env is
   harmless at runtime, so nothing red would fire). That residual is WHY the mechanism
   itself had to change, not just the premise.
3. **The platform probe trusts the kernel and the runner image's PATH — nothing
   else.** Revision 1 replaced `go env GOOS`/`go env GOARCH` with normalised
   `uname -s`/`uname -m` after quorum R3 measured the override channel live (P16:
   `GOOS=windows GOARCH=amd64 go env GOOS GOARCH` → `windows`/`amd64`; same env,
   `uname` → `Darwin`/`arm64`). What remains trusted, honestly: an EMULATED runner
   reports the emulated arch — consistent rather than wrong, since CI's own builds then
   also target the emulated pair; and `uname` resolves via the runner image's PATH, so
   a subverted image could still lie — supply-chain territory, outside this item's
   threat model. The wiring test pins the anti-reversion half (M7): the platform
   tokens cannot re-enter either file, and the probe cannot quietly stop being
   kernel-read. No misreporting or emulated runner is in use today (P17: image
   `ubuntu-24.04` with three agreeing in-run platform reports).
4. **`saw_good` is satisfiable by the pinned toolchain alone** (pinned ∈ KNOWN_GOOD), so
   a run where 1.25.6 and 1.24.9 both SKIP can still green on linux. Chosen: the pin is
   the coverage that matters to CI; the redundant-good floor would buy marginal
   assurance at real flake surface. The SKIP lines remain in the log table.
5. **Nothing here detects an upstream fix on darwin in CI.** The `(or GOOD NEWS)`
   branch's instruction ("re-derive the pin decision in design_docs/") still fires only
   attended. If the mission wants that surfaced automatically, that is Option B, priced.
6. **The instrument stays network-dependent for six of seven probes.** A durable
   toolchain-fetch outage now reds dev instead of going silent — chosen (S6: a disarmed
   detector must not pass), on the strength of P6's measurement, and the step's
   15-minute ceiling plus measured ~47 s runtime (P3) leave headroom. Revision 1 (R1's
   complete-coverage floor) sharpens this by one notch: a SINGLE known-bad fetch SKIP
   now reds, not just a durable outage. The price is measured, not waved: 0 SKIPs in 42
   observed probes across 6 runs (controller streak 10/10 consistent); the honest
   rule-of-three bound on zero observed events is ≈7% per probe at 95%, and the step is
   idempotent and re-runnable. Accepted as the price of a certificate that means
   complete coverage — the exact bargain the sibling doc demanded a flake measurement
   for before gating (P11).

## Quorum round 1

Full-strength review: 3/3 present, absent_reviewers EMPTY. Verdicts: R1 REJECT
(blocking), R2 PASS-with-catch (non-blocking, applied anyway), R3 REJECT (blocking).
This section is what round 2 should re-read the revision against.

### R1 — gpt5-6-sol — REJECT — **applied verbatim**

*Objection (one line):* two verified observations (darwin/arm64 affected, linux/amd64
clean) were generalised by a silent `expect_defect=0` default into "every non-darwin
platform is clean", and `ran_bad>=1` let a mostly-skipped negative arm certify.

- **Fail-closed polarity — the reviewer's `case` verbatim**, including the refuse-arm
  message text (`*) echo "INSTRUMENT FAILURE: no verified platform contract for
  $host_pair"; exit 1`). Every prose instance of "every other platform" / "beyond
  darwin/arm64" is rewritten to linux/amd64-scoped wording (thesis, Option A, the M1
  comments, and the PLATFORM ALARM text itself; the revision is grep-clean on both
  phrases). Proof rows added as demanded: **P18** (contract set is exactly two pairs;
  `verify_go.sh:214` + `w-race-gate-blindspot.md:207/:435`, both re-run first-party)
  and **P17**, which closes the constraint-1 consequence of the refuse arm: since an
  unlisted platform now EXITS 1, P17 proves from run 33069999131's own log (`Image:
  ubuntu-24.04`; setup-go cache key `setup-go-Linux-x64-ubuntu24-…`; first-party
  `AILANG v0.30.0 on Linux/x86_64`; `go1.26.6 linux/amd64` ×7) plus the
  `runs-on: ubuntu-latest` census (:19, :100) that the CI host maps into the listed
  `linux/amd64` arm, and replays P5's measured profile through the revised floors to
  exit 0 — the refuse arm cannot red the next push.
- **Complete coverage — taken, priced with measurements, not preference.** `ran_bad`
  must now equal `bad_expected` (counted at probe entry so it tracks the list), with
  `bad_expected==0` refused alongside so the hollowed-list teeth survive. The price —
  one transient fetch SKIP reds dev — is priced against evidence: P6 measures 42/42
  probes fetched and 0 SKIPs across 6 runs (controller streak 10/10 consistent); the
  honest rule-of-three bound on zero observed events is ≈7%/probe at 95%; the step
  re-runs in ~47 s of a 15-min ceiling (P3). This is precisely the flake-rate
  measurement the sibling doc demanded before gating (P11). The old Deferred Scope
  bullet that *accepted* partial SKIPs is discharged rather than defended, because the
  alternative — green after a 3/4 run — re-creates this row's own defect class one
  level down: a certificate saying more than the probes measured. Teeth rehearsed by
  M6 with a measured base (`go1.99.99` → SKIPPED → refusal, rc=1 quoted in the row).

### R2 — gemini-3-1-pro — PASS — **catch applied verbatim**

*Catch (one line):* `saw_bad` must not be able to act as a count the exact-match floor
would miss. The alarm condition is now the reviewer's `[ "$saw_bad" -ne 0 ]` verbatim.
Measured context: in today's generator `saw_bad` is a LATCH (`run.sh:55`,
`BUG*) [ "$expect" = BAD ] && saw_bad=1 ;;` — it can only ever be 0 or 1), so the
`-ne 0` form binds a future count-form edit rather than fixing a present bug; the same
invariant logic is why the coverage floor uses `-ne` against `bad_expected`, never `-lt`.

### R3 — oc-glm-5-2 — REJECT — **applied: premise row AND mechanism closed**

*Objection (one line):* the polarity's sole input, `go env GOOS`/`GOARCH`, honours the
env-var form of the platform tokens, and no premise verified that neither the workflow
nor the script sets them.

- **Premise row P16 added with the reviewer's exact commands, plus same-call
  controls.** Observed: ci.yml zero occurrences (grep rc=1; the `GOTOOLCHAIN`=2
  known-positive control proves the grep live in the same call); run.sh exactly ONE
  occurrence — the :61 banner READ, quoted in full in the row, verified not an
  assignment. The premise was TRUE TODAY and UNPINNED — that combination was the whole
  objection — so the revision pins it instead of re-asserting it (below).
- **Mechanism closed via the second half of the reviewer's own proposed_fix**, not
  merely the premise row. Option 1 (compute before any assignment) was a no-op — no
  assignment exists in either file; option 2 (a probe that cannot be overridden) is the
  closure: the polarity block now reads `uname -s`/`uname -m`, normalised to Go naming.
  Both halves of the mechanism claim are measured at P16: `GOOS=windows GOARCH=amd64
  go env GOOS GOARCH` prints `windows`/`amd64` (the channel is real), while the same
  env leaves `uname` printing `Darwin`/`arm64` (the new probe is immune), and the
  normalised pair equals the `go env` pair on this rig, so the banner now prints
  `$host_pair` — the identical string — removing the file's last read of the channel.
- **The pin (so tomorrow cannot undo it):** M2's wiring test now also asserts zero
  occurrences of the two platform-token strings in ci.yml AND run.sh, plus positive
  `uname -s` presence in run.sh — a revert to the old probe, or any new env assignment
  of the tokens, reds the verify_go leg (mutation M7, landed-proof included). **What
  it still cannot see** (residuals 2/3): a computed/obfuscated assignment
  (`eval "GO""OS"=…`) evades the substring pin; `uname` still trusts the runner
  image's PATH; an emulated runner reports the emulated arch — consistent with what CI
  builds for. Named, not defended further.

*Cross-check after revising:* every acceptance criterion and mutation row was re-read
against the revised design — M2's rehearsal moved onto the `case` arms (AC4(b)/(c)),
M4's predicted step-side floor was corrected by a measured /tmp sandbox base (it is
floor 1, `ran==0`, not the pinned floor), M6/M7 were added for the two new teeth, and
constraint 1 is carried by P17 plus AC6's green-on-arrival form.

## Quorum round 2 — and the carve-out that closed it

**Verdict: BLOCKED, 3/3 REJECT on ONE shared objection.** `absent_reviewers` was
NON-EMPTY (`oc-glm-5-2`, reason `invalid`) and was **restored before anything was acted
on**, per the shared skill's absent-reviewer rule: re-run alone with a raised cap
(`ailang design-review --reviewer oc-glm-5-2 --max-cost-usd 0.40`, $0.0232) → **reject**,
the same objection. So the round carries no hole; it is a genuine 3/3.

**Why it was recorded absent is itself worth the line, because the rule does not cover
it.** The absent-reviewer rule is written for a reviewer that is *unreachable* or
*over-budget*. This one answered, at cost, and its answer was discarded by the JSON
decoder — `invalid character 'G' after object key:value pair` — because the objection
QUOTED THE VERY STRING LITERALS UNDER DISCUSSION (`\"GOOS\"` and `\"GOARCH\"` inside a
`strings.Contains` snippet). The verdict was nonetheless recoverable verbatim from the
error payload, which is how the hole was seen before the re-run confirmed it. **A
reviewer can be silenced by the CONTENT of its review**, and a doc about string-literal
pins is exactly the doc that provokes it. Filed as a charter row, not fixed here.

**The shared objection, reproduced first-party against the exact proposed bytes before
any action.** M2's wiring test banned the substrings `GOOS`/`GOARCH` repo-wide; M1's own
`run.sh` comment block names both, because the comment exists to explain why the channel
is forbidden. The gate would have redded **on arrival**, not merely under mutation —
a constraint-1 violation authored by the doc against itself.

**Disposition: the ratified narrow-refinement carve-out** (ratified for this mission at
iteration 44, so the first-use gate does not apply). Every remaining blocking objection
(a) carried a concrete reviewer-authored `proposed_fix` and (b) disputed no design
DIRECTION — all three reviewers accept Option A; the defect is an internal contradiction
between two milestones. `gpt5-6-sol`'s fix is applied **verbatim** and is the superset:
require `uname -s`/`uname -m`; reject only EXECUTABLE uses of `go env GOOS`/`go env
GOARCH` and `host_pair` assignments derived from them; scope the `continue-on-error`
check to the miscompile step rather than banning the tokens globally. `gemini-3-1-pro`'s
alternative (strip the tokens from the comment) is **subsumed, not overridden** — it is
named inside `gpt5-6-sol`'s own *"Alternatively"* clause, and the scoped form removes the
conflict without deleting the warning the comment exists to carry. `oc-glm-5-2`'s
objection is the same defect and is satisfied by the same fix.

### V19 — the verification row `gpt5-6-sol` demanded, measured rather than promised

Its fix ends: *"Add a verification row running the exact proposed test against the exact
proposed run.sh bytes and recording one `=== RUN`, one PASS, and rc=0."* Done, in a
scratch module over a fixture tree holding the exact M1 `run.sh` and M2 `ci.yml` bytes
(`bash -n` rc=0 on the proposed script). Read by counting `=== RUN`/`--- PASS`/`--- FAIL`,
never by exit code — the negative control below is why.

| arm | what it runs | result |
|---|---|---|
| **A — the round-1 form** (repo-wide token ban) | identical fixture bytes | **rc=1, `=== RUN`=1, `--- PASS`=0, `--- FAIL`=1**, message `run.sh:37 re-introduces "GOOS"` — the objection reproduced |
| **B — `gpt5-6-sol`'s scoped fix, verbatim** | identical fixture bytes | **rc=0, `=== RUN`=1, `--- PASS`=1, `--- FAIL`=0** — exactly the row demanded |
| **negative control** | `-run '^TestNoSuchArmZZZ$'` | **rc=0**, `testing: warning: no tests to run` — proof that rc alone is not the reading |
| **M7a** revert `host_pair` to `"$(go env GOOS)/$(go env GOARCH)"` | arm B | **RED**, `run.sh:51 executable use of "go env GOOS"` (+GOARCH) |
| **M7b** replace `case "$(uname -s)" in` with `case "$(echo Darwin)" in` | arm B | **RED**, `run.sh no longer reads the kernel via "uname -s"` |
| **M1** insert `continue-on-error: true` INSIDE the miscompile step | arm B | **RED**, `ci.yml:172 re-introduces "continue-on-error" in the miscompile step` |
| **scoping control** insert `continue-on-error: true` in the UNRELATED `worldd benchmark smoke gate` step | arm B | **GREEN (rc=0, PASS=1)** — the boundary `gpt5-6-sol` asked for holds; an unrelated step is not frozen out |

Every mutation was asserted **LANDED by grep/count before its result was read** (M7a
`host_pair="$(go env GOOS)` count 0→1; M7b `uname -s` count 1→0; both flag arms
`continue-on-error` count 0→1), restored from a `/tmp` backup, and re-verified
**byte-identical by sha256** after each arm; a green control re-ran on the restored tree
and passed. The scoping control is the load-bearing one: without it, arm B's greens are
equally consistent with a test that stopped looking at `ci.yml` at all.

**Residual this carve-out adds, stated rather than left favourable:** V19 exercises the
assertion LOGIC against the exact bytes in a scratch module; it does not prove the test
compiles inside package `verifygate` against the real `repoRoot`. That is the executor's
AC1, and it is the one thing V19 deliberately cannot certify.

## Related Documents

- `design_docs/implemented/w-race-gate-blindspot.md` — the instrument's origin; its
  §Deferred Scope answers the linux/amd64 question this design is built on.
- `design_docs/implemented/w-setup-go-pin-unguarded.md` — the pin gate and Test B, whose
  declared residual this item discharges, and whose Deferred Scope named the flake
  measurement this item now supplies.
- `design_docs/implemented/w-canary-control-does-not-survive-a-floor-raise.md` — row 42's
  floor-binding for the same repro module; scoped neighbour.
- `design_docs/implemented/w-floor-raise-coupling-inventory.md` — structural sibling;
  its P7 ("line numbers rot; literals do not") is applied to the citation deferral here.
- Queue rows 31, 45, 48, 49, 50, 51 (world-mission.md) — deferred-to and scoped-out
  neighbours, each named in situ.

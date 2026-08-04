# w-race-gate-blindspot — the gate that is never run, and the compiler underneath it

**Status**: **IMPLEMENTED 2026-08-04 (iteration 46) — ITEM COMPLETE.** RG.A landed via PR #36 →
squash `f19acac`, dev CI green both jobs (SHA-addressed on the merge commit, step logs read to prove
the `-race` leg actually ran); evaluator `sonnet` **PASS 96/100, zero blocking**. `4e/OD-1`
DISCHARGED (floor at 1.25.6); `4e/OD-2` DISCHARGED (reproducer filed upstream as
[`golang/go#80706`](https://github.com/golang/go/issues/80706)); `4e/OD-3` remains DECLINED as
primary; `4e/OD-4` remains OPEN and is a cost question only. See §Implementation record.
~~INVESTIGATION COMPLETE, mechanism IDENTIFIED. Remediation **PARKED for human
ratification** (two decisions, §Open Decisions).~~ Clause-1. Item 4e in the charter queue.
**Authorship**: controller-authored (iteration 40, 2026-07-30) rather than produced by
`design-doc-creator` — see §Provenance for why, and note that the charter's iter-26 guardrail
obligation for such a doc (a Premise Verification Log and a Conflict Surface) is discharged below.
**Reproduction artifact**: `design_docs/verification/w-race-gate-blindspot/` (committed, re-runnable,
with its own known-positive controls).

---

## The finding in one paragraph

Item 4e was filed as *"a test that fails deterministically has sat in landed code for four
milestones, because the gate that would see it is never run"* — a `host/store` test that passes
without `-race` and fails 5/5 with it, plus a second test that hangs, with **zero `DATA RACE`
warnings** and the mechanism explicitly labelled UNKNOWN. The standing hypothesis was
`modernc.org/sqlite`'s heavy `unsafe` usage interacting with the race detector's altered memory
layout. **That hypothesis is refuted.** The cause is a **Go compiler code-generation regression in
go1.26.x** (present in 1.26.0 through **1.26.5**, the latest stable, on darwin/arm64) which
miscompiles a specific source shape that `host/store/scan.go` uses at exactly two sites — the two
functions the two symptoms belong to. It is **not** `-race`-specific: a default `go build` of a
52-line dependency-free program produces a **corrupt string header** (nil data pointer, garbage
length) for the same shape. `-gcflags=all=-N` fixes it, `-gcflags=all=-l` does not, and go1.25.6 /
go1.24.9 are clean. **One mechanism explains both of 4e's symptoms**: disabling optimizations for the
single package `host/store` clears the failure *and* the hang.

The consequential part is not the test, but it must be stated more carefully than the first draft of
this doc stated it. At HEAD, `go.mod` declares `go 1.26.4` and CI is configured with
`go-version-file: go.mod`. Read from the latest dev CI run's own step log, `actions/setup-go`
resolved this to **`go version go1.26.4 linux/amd64`** ("Setup go version spec 1.26.4" → "Successfully
set up Go version 1.26.4"), so CI **did** build this repository with an affected *release* — but on
**linux/amd64**, an architecture for which the miscompilation is **UNDETERMINED** (P11). The
demonstrated defect is bounded to **go1.26.0–1.26.5 on darwin/arm64**. What historical CI jobs and
individual developer builds resolved is **not claimed** beyond that one measured run. Which makes
AC6 — running the fixture in CI — the decisive experiment for CI's real exposure, rather than
something to reason about.

## Motivation — why this is worth a ratification gate rather than a quiet fix

This mission's bar is honest, non-vacuous gates. Item 4e is the **sixth** instance of its signature
shape and the first where the gate cannot fail *because it is never invoked*. The investigation turned
that into something sharper: the reason nobody could add the missing gate was that **the compiler was
wrong**, so `-race` looked untrustworthy and the standing advice was *"just add `-race` to CI" is not
yet known-good*. That advice was correct and for the wrong reason. Once the toolchain is sound, a
full-repo `-race` leg is **green in 179 s with zero data races** (measured, §PVL P9). The blocker was
never the race detector.

Two further reasons this is not a small cleanup:

1. **Default builds are exposed, not just `-race` builds.** The repository's non-`-race` tests pass
   today, which means the two production sites happen to produce correct values on the paths the
   tests cover. That is luck, not safety: the identical shape in a slightly different struct layout
   yields a corrupt string header on the same compiler.
2. **There is nothing to upgrade forward to.** go1.26.5 is the newest stable release and it is
   affected. The only remediations available are to move *back* to 1.25.x or to change the source.

## Premises (hard constraints — each verified in the Premise Verification Log)

- **P1** The two 4e symptoms reproduce at HEAD `8ed04c0`.
- **P2** `Ref` being empty is not evidence of anything (it is never assigned on that path).
- **P3** It is the optimizer, scoped to one package; not inlining.
- **P4** One switch clears both symptoms.
- **P5** It reproduces with no dependencies and no `-race`.
- **P6** It is a go1.26 regression, unfixed as of go1.26.5.
- **P7** Exactly two production sites carry the shape.
- **P8a** CI's *configuration* requests the `go.mod` version.
- **P8b** The *resolved* toolchain in the latest dev CI run was `go1.26.4 linux/amd64`.
- **P8c** That historical CI jobs and local builds *universally* used 1.26.x is **NOT CLAIMED**.
- **P9** A `-race` leg is green and costs ~179 s once the toolchain is sound.
- **P10** `modernc.org/sqlite` v1.54.0 requires only `go 1.25.0`, so a downgrade is dependency-feasible.
- **P11** Architecture scope beyond darwin/arm64 is **UNDETERMINED**.
- **P12** Neither `go.mod` directive can pin an *older* toolchain; only `GOTOOLCHAIN` can.

### Design Freeze (the sprint must not renegotiate these)

- The **mechanism claim is bounded**: "go1.26.0–1.26.5 miscompile this shape on darwin/arm64."
  No claim is made about other architectures, other shapes, or the precise SSA pass. The
  `late_fuse`/`generic_cse`/`prove` reading is a **sub-agent lead, unverified by the controller**,
  and must be labelled as such anywhere it is restated.
- **No claim may state or imply that the repository's default builds are currently producing wrong
  values.** They are not, on covered paths. The claim is that the compiler is capable of it and that
  nothing in the gate set would notice.
- The reproduction fixture stays a **nested, separate Go module** so the root module's `./...` never
  builds it. Verified: `go list ./...` returns exactly 10 packages with the fixture present.
- Any canary test added later must **fail loudly** on an affected toolchain. A canary that skips, or
  that passes by asserting "either correct or a known-bad version", is precisely the vacuous gate
  this item exists to eliminate.

## Decision 1 — the diagnosis is a toolchain defect, not a store defect

`host/store/scan.go` is correct Go. Nothing in it should change **for correctness reasons**. The
tests that fail are correct tests. The remediation belongs at the toolchain boundary.

Rejected alternative: *rewrite the two `scan.go` sites to dodge the shape* (hoist `fields` to package
scope, or copy `fields[i]` into a named local — both measured to avoid the defect). Rejected as the
**primary** remedy because it treats a compiler bug as a style rule, leaves every *other* shape the
same compiler gets wrong unguarded, and encodes a workaround nobody can later recognise as one. It
remains available as an optional belt (§Open Decisions OD-3).

## Decision 2 — pin the toolchain, and keep ENFORCEMENT separate from ASSERTION

Both `go.mod` directives are **floors, not ceilings** (measured — P12). With go1.26.4 installed,
neither `go 1.25.6` nor `toolchain go1.25.6` stops the local 1.26.4 from being used: under
`toolchain go1.25.6`, `go version` still reports 1.26.4, the built binary still reproduces the
defect, and `go version -m` stamps it **`go1.26.4`**. Only `GOTOOLCHAIN` selects an exact toolchain.
So a pin has to be enforced outside `go.mod` or it is decorative:

1. `go.mod` — lower the `go` directive to `1.25.6`. This does **not** pin anything; it is required
   only so a pinned 1.25.6 toolchain will consent to build the module at all.
2. `scripts/verify_go.sh` — **sets nothing at all.** It reads the toolchain that is actually active
   and fails loudly if that version is affected:

   ```sh
   ACTIVE_GO=$(go env GOVERSION)      # observe; never assign
   case "$ACTIVE_GO" in
     go1.26.0|go1.26.1|go1.26.2|go1.26.3|go1.26.4|go1.26.5)
       echo "verify_go.sh: FATAL: active toolchain $ACTIVE_GO miscompiles host/store/scan.go's" >&2
       echo "  array-literal shape (see design_docs/verification/w-race-gate-blindspot/). Pin a" >&2
       echo "  known-good toolchain, e.g. GOTOOLCHAIN=go1.25.6." >&2
       exit 1 ;;
   esac
   ```

   **Two earlier drafts of this step were both wrong, in the same direction, and the reviewers caught
   both.** Draft 1 exported `GOTOOLCHAIN=go1.25.6` unconditionally, which overrides a hostile
   `GOTOOLCHAIN=go1.26.4` handed in by a caller, so the assertion could never see a bad toolchain and
   **AC2 could never fail** — a gate that cannot fail, reintroduced inside the remedy for a gate that
   cannot fail (`gpt5-6-sol`, round 1). Draft 2 softened it to assign-if-unset, which still means a
   developer on a 1.26.4 rig with `GOTOOLCHAIN` unset gets a **silent fallback** to 1.25.6, a green
   check, and a false belief that their own environment is sound — violating the no-silent-fallback
   axiom and contradicting this section's own next sentence (`gemini-3-1-pro`, round 2). The form above
   sets nothing, so what it reports is what the caller actually has.

   Consequence, accepted deliberately: on an unpinned 1.26.x rig `verify_go.sh` now **fails**. That is
   the correct, loud, actionable behaviour — the pin belongs in the environment, and a verifier's job
   is to say so.
3. `.github/workflows/ci.yml` — an explicit `go-version: 1.25.6` (or `GOTOOLCHAIN` in the job env)
   rather than `go-version-file: go.mod`, since the file's floor semantics would let a newer runner
   image pick 1.26.x again.

The division of labour is deliberate: **CI and the developer's environment enforce; the script
verifies.** A verifier that also silently sets the thing it verifies is not a verifier.

## Decision 3 — the canary is the durable form of this finding

Prose decays; a committed test does not. The follow-up sprint adds a canary that carries the
miscompiled shape and asserts the correct value, so that **re-introducing an affected toolchain reds
the suite** instead of silently restoring the exposure. The canary and the pin must land in the
**same change**: the canary reds by design on the current toolchain, so it cannot land first.

## Milestones (deliberately small; total ~0.25–0.5d, gated on ratification)

### RG.A — pin + canary + `-race` leg (~0.25–0.5d, BLOCKED on OD-1)

- `go.mod` `go 1.25.6`; a **read-only** affected-version assertion in `scripts/verify_go.sh` (it sets
  no toolchain — Decision 2); explicit Go version in `ci.yml`; `go version` printed and archived
  before every build leg.
- `host/store/toolchain_canary_test.go` — the miscompiled shape, asserting the correct value.
- A `-race` leg in `scripts/verify_go.sh` and the CI `go-verify` job, bounded, with `AILANG_BIN` set.
- Run the reproduction fixture once in CI to settle the linux/amd64 question (AC6).
- Owns **AC1, AC2, AC2b, AC3, AC4, AC5, AC6, AC7** — enumerated, not a range (iter-29 rule), and this
  is the only milestone, so every AC has exactly one owner and RG.A holds the files that can fail each.

## Files to Create/Modify

| Path | Change |
|---|---|
| `go.mod` | `go 1.26.4` → `go 1.25.6` |
| `scripts/verify_go.sh` | `GOTOOLCHAIN` pin + loud assertion; add the bounded `-race` leg |
| `.github/workflows/ci.yml` | explicit Go version; `-race` leg in `go-verify` |
| `host/store/toolchain_canary_test.go` | NEW — the canary |
| `design_docs/verification/w-race-gate-blindspot/**` | **already landed this iteration** (fixture) |

## Conflict Surface (MANDATORY — every landed behaviour this design could collide with)

| Landed behaviour | Collision risk | Resolution |
|---|---|---|
| **CI `go-verify` job** (`w-world-library-m1` M6) selects the toolchain via `go-version-file: go.mod` | Changing the `go` directive silently changes CI's compiler — intended here, but it also means any future bump to the directive re-exposes the repo | The pin is asserted in `verify_go.sh` with a loud failure, so the exposure cannot return silently |
| **`scripts/verify_go.sh`** anti-false-green guard for `AILANG_BIN` | Adding a second guard must not weaken the first | New assertion is additive; the existing `AILANG_BIN` check is untouched |
| **CF-MJC-1** (MJ.C's two `maxRecoveryPages` bound tests, ~50 s) | A `-race` leg makes them the dominant cost — `host/broker` is 176.6 s of the 179 s total | Stated, not hidden. Injectable-bound decision is CF-MJC-1's, escalated as OD-2 |
| **`host/store/scan.go`** (landed `w-store-durability` SD.A) | Tempting to "fix" the source | Decision 1: source is correct; no correctness-motivated change. OD-3 keeps the belt optional |
| **`bench/BASELINE.md`** comparability rules (item 4f) | A toolchain change invalidates every baseline number in the file | **Flagged**: the pin must be recorded in `BASELINE.md` as a condition change, and the amortisation rows (already pinned to M3.C idle-rig figures, CF-MJC-2) re-derived on the new toolchain. Item 4f owns the mechanism |
| **`modernc.org/sqlite` v1.54.0** | Might require go ≥ 1.26 | Refuted (P10): its own `go.mod` says `go 1.25.0` |
| **Replay determinism / archived interpreter binaries** (`host/replay`, Decision 7 of M1) | Bit-for-bit replay resolves an archived `ailang` binary; a *host* compiler change does not touch archived artifacts, but a reader could confuse the two | No collision: the pinned AILANG interpreter (`v0.30.0`, `e37b370`) is independent of the Go toolchain that builds the host. Stated so the distinction is on the record |

## Systemic-Issue Audit

This is the **sixth** instance of *a gate that cannot fail* and the second in which the *instrument*,
not the code, was the defect (after iter-37's zsh-glob-mangled `grep` and iter-38's structurally-zero
skip count). It adds a genuinely new member to the family: **the compiler is an instrument too.**
Every gate in this repository — `ai-check`, `go test`, `verify_go.sh`, the mutation protocol itself —
reports through a compiler, and a mutation-testing discipline is only as sound as the code generator
underneath it. A named RED mutation that "reds as predicted" on a miscompiling toolchain proves less
than it appears to.

## Deferred Scope

- ~~Determining whether **linux/amd64** (what CI runs) is affected.~~ **ANSWERED 2026-08-04
  (iteration 46, AC6): linux/amd64 is NOT affected.** All four affected toolchains return `OK` on
  `ubuntu-latest` (run `30872300227`, merge commit `f19acac`, step *"Measure compiler reproducer on
  linux/amd64 (non-gating)"*), plus both known-good toolchains — so `run.sh` reports `INSTRUMENT
  FAILURE (or GOOD NEWS): no known-affected toolchain reproduced the defect`, which on amd64 is the
  GOOD NEWS branch. **Consequence, recorded because it cuts AGAINST this doc's own motivation
  section**: CI was never building this repository through the miscompilation, so the historical
  blast radius was bounded to local darwin/arm64 developer builds the whole time. The pin remains
  correct — it is what makes the two environments agree and what stops a developer shipping from a
  rig whose gates cannot be trusted — but the "default builds are exposed" argument was broader than
  the evidence now supports, and that is stated rather than left as a favourable silence. Recorded
  permanently in `bench/BASELINE.md` too: a step log nothing requires anyone to read is not a record.
- Naming the responsible SSA pass with certainty. A lead exists; verifying it needs compiler-source
  work with no bearing on the remediation.
- Auditing the `ailang` interpreter repository (`sunholo-data/ailang`) for the same exposure. It is a
  Go project on the same rig, so it plausibly shares it — but it is a **different mission's
  repository** and the charter forbids working in it. Routed as a cross-mission notification.

## Acceptance Criteria

- **AC1** `go.mod` declares a Go version with no known miscompilation of the fixture's shape, and
  `design_docs/verification/w-race-gate-blindspot/run.sh` still reports both controls firing.
- **AC2** `scripts/verify_go.sh` exits **non-zero with a named error** when the active toolchain is an
  affected version, demonstrated in a **default 1.26.4 environment with no explicit `GOTOOLCHAIN`
  override** (this rig's default — the case a developer actually hits), and demonstrated green under
  `GOTOOLCHAIN=go1.25.6`. **Both legs must appear in the milestone's evidence**; the first is the one
  that proves the assertion can fail at all.
- **AC2b** Every build leg in `scripts/verify_go.sh` and in the CI `go-verify` job **prints and
  archives `go version`** before building, so the resolved toolchain is a recorded fact in every run
  rather than something a later reader has to infer from configuration (the P8b/P8c distinction).
- **AC3** `host/store/toolchain_canary_test.go` **fails** under `GOTOOLCHAIN=go1.26.4` and **passes**
  under the pinned toolchain — both legs shown, not asserted.
- **AC4** `go test ./... -count=1 -race` is green with `AILANG_BIN` set, and its wall-clock cost is
  recorded in the log entry (expected ~179 s; `host/broker` dominant).
- **AC5** The CI `go-verify` job runs the `-race` leg, and the **step log is read** to confirm it
  actually ran (not merely that the job was green — the iter-38 rule).
- **AC6** The fixture is executed in CI at least once, so the **linux/amd64** question in Deferred
  Scope is answered with a measurement rather than left open.
- **AC7** `bench/BASELINE.md` records the toolchain change as a condition change, per the Conflict
  Surface row.

## Non-Vacuity — the named RED mutation for every gate (S6)

| Mutation | File (production/test) | Predicted result |
|---|---|---|
| **`MUT-CANARY-BLIND`** — change the canary's expected value to `""` (what an affected toolchain produces) | test | Canary passes on BOTH toolchains → proves the canary's assertion, not its presence, is what discriminates |
| **`MUT-PIN-REMOVED`** — delete the `GOTOOLCHAIN` export from `verify_go.sh` | production (gate script) | `verify_go.sh` goes RED on this rig via the canary; if it stays green the pin was decorative |
| **`MUT-VERSION-ASSERT-DEAD`** — short-circuit the affected-version comparison in `verify_go.sh` to always-false | production (gate script) | Script green under `GOTOOLCHAIN=go1.26.4` → proves the assertion carries the teeth, and guards against the **iter-36 `MUT-DDL-DRIFT` failure mode** where a gate reds by a different mechanism than the documented one |
| **`MUT-RACE-LEG-DROPPED`** — remove `-race` from the CI leg | production (CI) | Neither suite reds → confirms only the leg's presence, not any assertion, provides `-race` coverage; this is a **drift check, not a kernel proof**, and is labelled so per the iter-30 rule |

Note explicitly, per the iter-30 rule: `MUT-CANARY-BLIND` and `MUT-RACE-LEG-DROPPED` edit **test/CI**
files, so they are drift checks. The only production-side teeth here are in the gate scripts, because
**the artifact under test is the build configuration**, not kernel behaviour. That is stated rather
than dressed up as a kernel proof.

## Open Decisions (ESCALATED — these do NOT have controller defaults; RG.A is blocked until answered)

- **OD-1 (ratification-class) — lower the repository's Go toolchain from 1.26.4 to 1.25.6?**
  Recommended: **yes.** The newest stable Go miscompiles a shape present in landed durability code,
  there is no fixed release to move forward to, and the only dependency constraint
  (`modernc.org/sqlite`) needs just 1.25.0. Cost: the `go` directive is a compatibility promise, and
  the repo forgoes 1.26 language/stdlib features (none are currently used). This is a build-toolchain
  change affecting all future work, so it is Mark's call and not the controller's.
- **OD-2 — file the reproducer upstream at `golang/go`?** Recommended: **yes**, it is a clean 52-line
  regression against the latest stable. Parked because it is a **public post to a third-party
  project**, which is outside anything this loop has been authorised to do; the charter's language-gap
  rule routes to `sunholo-data/ailang`, and this is not an AILANG gap. Repro is ready to paste.
- **OD-3 — additionally change the two `scan.go` sites to a shape the affected compiler gets right?**
  Recommended: **no** for correctness (Decision 1), **defer** as a belt. Worth revisiting only if
  OD-1 is declined, in which case it becomes the *only* available mitigation.
- **OD-4 (CF-MJC-1, inherited) — make `maxRecoveryPages` injectable?** With a `-race` leg,
  `host/broker` costs 176.6 s of the 179 s total. Injecting the bound costs a small production change
  and ~0.001 s. Not decided here; surfaced with the number attached.

## Provenance — why this doc is controller-authored

The skill routes a new-doc item to `design-doc-creator` on the rotation designer. Every load-bearing
claim here is a **measurement the controller ran first-party this iteration**, and the charter's own
iter-105/iter-27 guardrails say a finding must not be laundered through another author into an
"established fact". Handing these measurements to a designer to restate would have added exactly that
hop without adding design content, since the remediation space is narrow and fully costed. The
obligation the charter attaches to a doc not produced by `design-doc-creator` (iter-26 guardrail — a
Premise Verification Log and a Conflict Surface) is discharged above. **This deviation is recorded in
the iteration-40 log's routing-evidence row rather than left implicit**, and the doc still goes
through the pick-time quorum, which is the independent-eyes gate that actually matters here.

## Premise Verification Log (live evidence, 2026-07-30, iteration 40)

Full captured output: `design_docs/verification/w-race-gate-blindspot/OUTPUTS.md`.

| # | Premise | Evidence | Provenance |
|---|---|---|---|
| **P1** | Both 4e symptoms reproduce at HEAD | `TestScanUnreadableLogKeysetResumes` rc=1 under `-race` (`Field:` empty), rc=0 without; `TestScanUnreadableWorldsFindsPoison` `panic: test timed out after 1m0s` under `-race`, `ok 0.184s` without | **VERIFIED BY ME** |
| **P2** | `Rows[0].Ref` empty is not evidence | `scan.go:76-79` — `ScanUnreadableLog` never assigns `Ref`; empty in every run | **VERIFIED BY ME (code read)** — refutes an iter-38 inference |
| **P3** | Optimizer, one package; not inlining | `-race -gcflags='…/host/store=-N'` → ok; `-race -gcflags='all=-l'` → FAIL | **VERIFIED BY ME** |
| **P4** | One switch clears both symptoms | Same `…/host/store=-N` flag makes the failing test AND the hanging test pass | **VERIFIED BY ME** |
| **P5** | Reproduces with no deps and no `-race` | 52-line program, `go vet` clean, cache cleared: plain `go build` → `len(Field)=4334851712`, SIGSEGV | **VERIFIED BY ME** |
| **P6** | go1.26 regression, unfixed at 1.26.5 | 1.24.9 OK · 1.25.6 OK · 1.26.0/.3/.4/.5 BUG. go1.26.5 is the latest stable per `go.dev/dl?mode=json` | **VERIFIED BY ME** |
| **P7** | Exactly two production sites | `grep -rn ':= \[\.\.\.\]' --include='*.go' \| grep -v _test.go` → `scan.go:74`, `scan.go:112`, with a known-positive control in the same call | **VERIFIED BY ME** |
| **P8a** | CI configuration requests the `go.mod` version | `ci.yml:58`: `go-version-file: go.mod`; `go.mod`: `go 1.26.4` | **VERIFIED BY ME (code read)** |
| **P8b** | The resolved CI toolchain was `go1.26.4 linux/amd64` | Step log of run `30483249118`: "Setup go version spec 1.26.4", "Successfully set up Go version 1.26.4", `go version go1.26.4 linux/amd64` | **VERIFIED BY ME (CI step log read, not inferred from config)** |
| **P8c** | Historical CI jobs / local builds universally used 1.26.x | — | **NOT CLAIMED.** One run measured; the rest is unestablished, and the `go` directive being a floor (P12) means the config alone cannot establish it |
| **P9** | `-race` leg green, ~179 s | `go test ./... -count=1 -race` rc=0, 10/10 packages ok, 179 s, 0 `DATA RACE` (control: 10 `ok` lines) under `GOTOOLCHAIN=go1.25.6` + `AILANG_BIN` | **VERIFIED BY ME** |
| **P10** | Downgrade is dependency-feasible | `modernc.org/sqlite@v1.54.0/go.mod`: `go 1.25.0` | **VERIFIED BY ME** |
| **P11** | Non-arm64 scope undetermined | `GOARCH=amd64` binary cannot execute here: `bad CPU type in executable`, no Rosetta | **VERIFIED BY ME that it is UNDETERMINED** |
| **P12** | `go.mod` cannot pin an older toolchain; only `GOTOOLCHAIN` can | With `toolchain go1.25.6` in `go.mod` and go1.26.4 local: `go version` → `go1.26.4`, program → `BUG`, and `go version -m ./t1` → **`./t1: go1.26.4`**. Controls in the same run: explicit `GOTOOLCHAIN=go1.25.6` → `OK`; no toolchain line → `BUG` | **VERIFIED BY ME** — added because quorum r1 flagged the claim as unlogged, and it **refutes** that reviewer's counter-claim (see §Quorum verification log) |
| **L1** | `late_fuse`/`generic_cse`/`prove` interaction; string ptr/len stored one slot late | SSA dumps read by a sub-agent | **UNVERIFIED, inherited from a sub-agent — re-check before relying. Not load-bearing for any decision here** |

## Quorum verification log

**Round 1** — `2026-07-30T00:04:26Z`, reviewers `gpt5-6-sol` + `gemini-3-1-pro`, **both present**
(`absent_reviewers: []`), controller verdict `pass`, cap raised to $0.25/reviewer per the iter-36
process fix. Outcome **BLOCKED** (rc=3), total metered cost **$0.056**. Artifact:
`.ailang/state/mission-quorum/w-race-gate-blindspot-2026-07-30T00-04-26Z.json`.

Two blocking objections, **both actioned in this one revision**; each was re-measured first-party
before being applied or rejected, per the charter's "a reviewer's finding is a claim" discipline.

1. **`gpt5-6-sol` — ACCEPTED IN FULL, and it found a real defect.** *"The proposed toolchain guard is
   internally self-defeating: `scripts/verify_go.sh` unconditionally exports `GOTOOLCHAIN=go1.25.6`
   before checking the resolved version, while AC2 requires the script to fail when invoked with
   `GOTOOLCHAIN=go1.26.4`. The export overrides that hostile input, so the script resolves 1.25.6 and
   AC2 cannot pass as written."* Correct. The remedy for a gate-that-cannot-fail contained a
   gate that could not fail. ~~Decision 2 now uses assign-if-unset (`: "${GOTOOLCHAIN:=go1.25.6}"`) so
   a caller's explicit value wins, enforcement is separated from assertion, and AC2 requires **both**
   legs to be demonstrated.~~ **SUPERSEDED at round 2** — assign-if-unset was still a silent fallback;
   the landed form sets nothing at all. See objection 4 below.

2. **`gemini-3-1-pro` — HALF ACCEPTED, HALF REFUTED BY MEASUREMENT.** The objection was that Decision
   2 rests on an unverified claim ("neither does a `toolchain` directive"), that the claim is missing
   from the Premise Verification Log, and that it is *"factually false — the Go 1.21+ toolchain
   directive explicitly forces the download and use of the exact specified compiler version,
   overriding local toolchains"*, so the design *"unnecessarily rejects the native, robust go.mod
   fix"*.
   - **The procedural half is valid and is fixed**: the claim was indeed load-bearing and unlogged.
     It is now **P12**, with evidence.
   - **The factual half is REFUTED.** Measured directly: `go.mod` containing `toolchain go1.25.6` with
     go1.26.4 installed → `go version` reports **go1.26.4**, the built program reproduces the defect,
     and `go version -m` stamps the binary **`go1.26.4`**. Run in the same session, explicit
     `GOTOOLCHAIN=go1.25.6` → `OK` (positive control) and no toolchain line → `BUG` (negative
     control). The directive raises a floor; it does not cap. Had this "native, robust go.mod fix"
     been adopted on the reviewer's authority, **the pin would have been decorative and the repo would
     have stayed on the miscompiling toolchain while the doc claimed it was pinned** — the same
     false-green class, arrived at by deferring to a reviewer instead of measuring. Recorded here
     rather than silently dropped.

**Round 2** — `2026-07-30T00:07:14Z`, same two reviewers, **both present**, controller verdict `pass`.
Outcome **BLOCKED** again, cost **$0.067**; cumulative metered **$0.123**. Artifact:
`.ailang/state/mission-quorum/w-race-gate-blindspot-2026-07-30T00-07-14Z.json`. Both objections were
again correct, and both landed on the revision I had just made:

3. **`gpt5-6-sol` — ACCEPTED. My own new P12 refuted my own headline claim.** *"The load-bearing
   exposure claim is unverified and internally contradicts P12: the doc says every CI run, every
   `scripts/verify_go.sh` run, and every local build used an affected compiler, but P8 proves only
   that `go.mod` declared 1.26.4 and CI referenced that file. A `go` directive is a floor, as the doc
   itself establishes, so it cannot prove the resolved toolchain."* Exactly right: adding P12 to
   satisfy round 1 made the round-1 exposure paragraph self-contradictory, and I did not notice.
   Fixed by their prescription — the paragraph is rewritten, and P8 is split into **P8a**
   (configuration, code read) / **P8b** (resolved toolchain) / **P8c** (**NOT CLAIMED**). Their fix
   said P8b was *"UNDETERMINED until a job log records `go version`"*; rather than record it as
   undetermined I **went and read the job log**, so P8b is now a measurement:
   `go version go1.26.4 linux/amd64` in run `30483249118`. Their "print and archive `go version`
   before every build leg" is adopted as **AC2b**.

4. **`gemini-3-1-pro` — ACCEPTED IN FULL, and it caught me contradicting myself one line later.**
   *"The proposed `verify_go.sh` modification (`: "${GOTOOLCHAIN:=go1.25.6}"`) violates the 'no silent
   fallbacks' axiom and contradicts its own stated principle ('A verifier that also silently sets the
   thing it verifies is not a verifier'). If a developer's environment has 1.26.4 installed and an
   unset GOTOOLCHAIN, the script will silently fallback to 1.25.6, pass its own check, and exit
   green, leaving the developer falsely believing…"* Correct. I wrote the principle and violated it in
   the same decision. Their fix is applied verbatim: the script **sets nothing**, reads
   `go env GOVERSION`, and fails loudly on an affected version; AC2 now requires the failure to be
   demonstrated in a **default 1.26.4 environment with no override**.

**Disposition — narrow-refinement carve-out APPLIED (bounded 2nd revision).** After the one revision
and one re-quorum the charter allows, every remaining blocking objection (a) carried a concrete
**reviewer-authored `proposed_fix`** and (b) disputed no part of the design **direction** — both were
honesty-of-claim and enforcement-mechanism defects, i.e. the completeness/determinism class the
carve-out names. Both fixes are applied as the reviewers wrote them; nothing was overridden and no
controller-invented resolution was substituted. Note the carve-out's first-use ratification clause is
long discharged in this mission (`w-m1-ailang-hardening`, `w-effect-journal`). Recorded in the
iteration-40 log's routing-evidence row.

**Standing rule 2 is not engaged**: this doc's sprint (RG.A) does **not** proceed on controller
authority regardless, because **OD-1 is a ratification-class decision parked for Mark**. What lands
this iteration is the investigation, the fixture and this ratification packet — no build change.

**The pattern worth keeping from these two rounds**: four blocking objections, and **three of them
were against text I wrote to satisfy an earlier objection**. Each revision introduced a defect of the
same family it was fixing — a silent fallback inside a no-silent-fallback remedy, a
self-contradicting exposure claim created by adding the premise that contradicted it. A revision is
not a smaller change than an original; it deserves the same adversarial read.

## Related Documents

- `design_docs/verification/w-race-gate-blindspot/` — the committed reproduction fixture
- `design_docs/implemented/w-store-durability.md` — SD.A, where `scan.go` landed
- `design_docs/implemented/w-effect-journal.md` — CF-MJC-1 (`host/broker` `-race` cost)
- `bench/BASELINE.md` — comparability conditions; item 4f
- `design_docs/coding-standards.md` — S6 (honest, non-vacuous gates)

---

## Implementation record (RG.A — iteration 46, 2026-08-04)

Landed as **PR #36 → squash `f19acac`** on `dev`, CI green on **both** jobs, SHA-addressed on the
merge commit and the step logs read rather than the badge trusted (run `30872300227`): the `-race`
banner appears exactly once with the plain-leg banner as its known-positive control, the
race-detector control emits 2 `WARNING: DATA RACE` lines, 21 `ok` lines, and the single `FAIL`
substring in the log is inside the fixture's own expected `INSTRUMENT FAILURE (or GOOD NEWS)`
message, not a test. Evaluator `sonnet` **PASS 96/100, zero blocking**, having independently
reproduced AC3 both legs, AC2 leg 1, AC4, `MUT-CANARY-BLIND` and `MUT-RACE-LEG-DROPPED` with its own
sha256 proofs.

Five commits, each independently green at its boundary, reconstructed by the controller from the
executor's cumulative `.snap/M<k>/` snapshots (the codex sandbox forbids git writes) and proven
byte-identical to the executor's final tree by `shasum -c`, 8/8 OK:
`77ce069` (M1) · `d66e12a` (M2) · `08fdb83` (M3) · `252ab51` + `01e6a01` (controller corrections).

### Acceptance criteria

| AC | Verdict | Evidence |
|---|---|---|
| **AC1** | MET | `GOTOOLCHAIN=go1.25.6 go build ./...` rc=1 before the floor change (`go.mod requires go >= 1.26.4`), rc=0 after. Fixture re-run first-party at `7550ee9`: 1.26.0/.3/.4/.5 `BUG`, 1.25.6/1.24.9 `OK`, `-N` `OK`, `-l` `BUG`, **both controls fired** |
| **AC2** | MET | Leg 1 (the one that proves the assertion can fail at all): default 1.26.4 env, `GOTOOLCHAIN` unset → rc=1, named FATAL, and **zero** `── go build` lines while the known-positive `── AILANG_BIN=` line fires in the same output. Leg 2: under the pin, rc=0 end-to-end |
| **AC2b** | MET | `go version` before the control, the build and each test leg; CI additionally prints `go env GOVERSION`. Both read from the merge-commit step log as `go version go1.25.6 linux/amd64` |
| **AC3** | MET | Canary **FAILS** under `GOTOOLCHAIN=go1.26.4` (`Field="" want "stateRoot"`), **PASSES** under `go1.25.6` — both legs run by the controller and again by the evaluator |
| **AC4** | MET | `go test ./... -count=1 -race` rc=0, 10/10 `ok`, **zero data races**, run four times across controller and evaluator plus once in CI |
| **AC5** | MET | Merge-commit step log read: the `-race` leg banner present once, with the plain leg present as its control |
| **AC6** | MET | **linux/amd64 is NOT affected** — see §Deferred Scope |
| **AC7** | MET | `bench/BASELINE.md` records the toolchain move as a condition change; no benchmark number re-measured. Item 4f owns re-derivation |

### Named RED mutations

All five legs applied with an asserted single-site replacement and a sha256 before/after/restored
triple, byte-identical restore verified each time — a mutation not proven **applied** is not a
mutation.

| Mutation | Result | Note |
|---|---|---|
| `MUT-CANARY-BLIND` | `ok` on 1.26.4, **FAIL** on 1.25.6 | **The doc's prediction ("passes on BOTH toolchains") is REFUTED** — independently by planner, executor, controller and evaluator. The SPLIT still establishes what the mutation was for: the canary's *assertion*, not its presence, is what discriminates |
| `MUT-PIN-REMOVED` | both legs red | **Re-pointed**: the doc targets "the `GOTOOLCHAIN` export in `verify_go.sh`", which after Decision 2's round-2 revision does not exist — the script sets nothing. Re-pointed at `go.mod`'s floor, predicted result preserved |
| `MUT-VERSION-ASSERT-DEAD` | reaches `── go build` | Assertion's teeth removed ⇒ the gate proceeds on an affected toolchain. Production discriminator |
| `MUT-RACE-LEG-DROPPED` A | green | Presence-only drift check, labelled as such |
| `MUT-RACE-LEG-DROPPED` B | red, named | Dropping `-race` from the known-positive control reds the gate with *"the race detector is not armed; every 0-races result in this gate is void"*. **This is what upgrades the criterion from a drift check to a production discriminator** |

### Honest tally

**4 of 8** criteria are production discriminators (AC1, AC2, AC3, AC4 — AC4 only because the
known-positive control makes it one). AC2b, AC5, AC6, AC7 are records and drift checks. The
evaluator argued 4/8 is defensible and could be read as 3/8 under the strictest counting; it is
recorded here at 4/8 with that dissent attached rather than silently rounded up.

### What was added beyond the design

The doc specified a `-race` leg. It did not specify how anyone would know the detector was armed —
so a green `0 races` would have been unfalsifiable, which is this mission's signature defect wearing
the remedy's clothes. `racecontrol/` is a nested module (invisible to `./...`; `go list ./...` still
reports exactly 10 packages) holding a deliberate data race that the gate runs FIRST and aborts on
if no `WARNING: DATA RACE` appears.

### The number that would not hold still

The doc says the `-race` leg costs ~179 s. Six measurements at one commit gave `host/broker` at
69 / 76.9 / 96.6 / 120.7 / 175.3 s on darwin/arm64 and 131.8 s in CI — a **2.54x spread** at nominal
load. The first correction replaced the single figure with "a range, 69-121 s"; the evaluator's
first independent run measured 175.3 s and falsified it. **A range you stopped measuring at is just
a wider single number.** `bench/BASELINE.md` therefore states the sample and its spread, says no
upper bound is established, and derives only what follows: the leg roughly doubles the gate, and
CI's `timeout-minutes: 25` is sized against an unknown tail with expiry routing to `4e/OD-4` rather
than being silently raised.

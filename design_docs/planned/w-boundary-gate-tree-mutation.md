# w-boundary-gate-tree-mutation — the gate must not be able to poison the tree it guards

**Status**: In sprint — **`BG.A` LANDED 2026-08-06 (iter-58)**, PR #47 → squash `278f102`. `BG.B`
and `BG.C` remain; the doc stays in `planned/` until all three land.
**Date**: 2026-08-06
**Queue item**: 10, `w-boundary-gate-tree-mutation` (promoted ahead of `SM.B2a` at iter-56, on the
strength of the SIGKILL-residue measurement below)
**Estimated**: ≤1 day (one test file, no production code; see Deferred Scope for what was cut to
hold that line)
**Designer**: controller-designated (iter-56)
**Toolchain boundary**: every command in this design was run in a worktree of `origin/dev`
@ `deeb804` with `GOTOOLCHAIN=go1.25.6` — the operative toolchain, since `scripts/verify_go.sh`
FATALs on go1.26.0–go1.26.5. No AILANG syntax is written or changed by this design (the AIL arm
changes only its Go-side reads), so the pinned `ailang` binary is not exercised here.

> **Thesis:** `host/boundary`'s teeth-proof writes forbidden imports into live production source.
> A SIGKILL mid-arm makes that write permanent, invisible to `go build`, and indistinguishable
> from a real boundary violation. The repair routes the mutant through `go list -overlay` plus an
> overlay-aware read, so the gate never writes a byte under the repository root — which closes the
> crash-residue mode and the concurrency window with one mechanism.

**Quorum revision (round 1, iter-56).** Both reviewers blocked on one thing: the earlier draft's
AC1/M3 detection mechanism, a per-arm timing-dependent polling detector whose green state was
"zero observed changes." They were right, and the irony deserves naming: a probabilistic detector
that can pass without ever observing the threat is the *same defect class this item exists to
fix* — "zero observations" and "the detector never ran" are the same result, and a search that
found nothing is a claim, not a fact (this mission's oldest recorded lesson). That mechanism is
**deleted, not supplemented**, and replaced by two complementary layers, one per reviewer's
proposed fix: **structural write confinement + an AST guard** (deterministic by construction,
filesystem-independent — Decision 7–8) and a **deterministic `ModTime`+sha backstop**
(controller-measured, C1 — Decision 9). Nothing else changed direction: the `go list -overlay`
decision, the alternatives-rejected analysis, the ordering argument, and the Verification Log
stand as reviewed.

**Quorum round 2 (iter-56) — BLOCKED again, on two NEW and narrower objections; resolved by the
controller narrow-refinement carve-out.** Both reviewers were present (`gpt5-6-sol` $0.0611,
`gemini-3-1-pro` $0.0253); **neither disputed the design DIRECTION**, both supplied a concrete
`proposed_fix`, and both objections were about determinism/completeness — so the carve-out applies
and the fixes below are the reviewers' **verbatim** remedies, not controller-invented resolutions.
Round 1's fix had, ironically, left a *second* instance of the very defect it removed:

1. **`gpt5-6-sol` — AC4/M5 was still probabilistic.** Killing at "≥10 randomized offsets" never
   established that any kill landed *after* an arm created and consumed its overlay, so a clean
   tree could mean the threat was never exercised. AC4 and M5 are rewritten to its proposed
   bounded synchronization protocol: an arm-specific ready marker, artifact verification, then
   `SIGKILL`, with a missing marker / timeout / absent artifact / early exit **failing** the
   criterion, and one run per arm instead of ten randomized kills.
2. **`gemini-3-1-pro` — P9's audit had not searched for its own threat model.** V15 grepped
   `os.WriteFile` only, while Decision 8 names `os.OpenFile`, `os.Create` and `os.Rename` as
   mutators. The controller ran gemini's widened search verbatim (V15a) and then, going beyond it,
   measured the property directly instead of reading 25 call sites (V15b): the suite produces
   **1** dirty observation with the boundary gate enabled and **0** without it. The premise
   survives — but it had genuinely been unverified, and it is the reviewer who caught that.

The through-line across both rounds is one sentence: **a green result must be unable to mean
"the check never ran."** Round 1 removed a detector that could; round 2 removed an acceptance
test that could; and the audit backing P9 is now a measurement rather than a search.

---

## The finding in one paragraph

`TestWorldBoundaryDependencyAllowlist` proves it has teeth by mutating three other packages'
production files **on disk** (`mutateAndRestore`,
`host/boundary/allowlist_world_test.go:190`), running the guard, and restoring. The restore is a
`defer` plus an explicit write — correct for every path that returns or `t.Fatal`s (V13), and
**never run on SIGKILL**. Reproduced first-party in this design's worktree (V10–V11): a `kill -9`
delivered the instant `host/store/store.go` changed left the mutant on disk permanently
(`rc=-9`, residue sha `d08c3f46…` — byte-identical to the gate's own mutant), after which
`go build ./...` → **rc=0** (the mutant import is a real package; the build gate is blind) while
the boundary gate itself → **rc=1**, accusing the innocent `host/store/store.go` of a
`forbidden registry/HTTP/cloud dependency "net/http/httputil"` — a message byte-identical in shape
to a real violation. This loop is kill-heavy by design: sub-agent lanes end in `kill`/`kill -9`,
Gate 3b polls have 30-minute deadlines, a launchd watchdog reclaims slots at 6h, and the repo's
**own** gate script SIGKILLs the process group when the `-race` leg exceeds 600s
(`scripts/verify_go.sh:113`, V17). Residue would land in the shared main checkout, which Critical
Principle 0 forbids the next iteration from stashing or resetting away — so the failure needs
attended intervention. Worse, the next queued milestone (`SM.B2a`) adds ~780 LOC of network code:
a residue-red "forbidden dependency" appearing during that sprint would accuse exactly the code
whose job is to touch the network, and nothing in the message distinguishes the two.

## Premises

- **P1 — `defer` does not run on SIGKILL, so the restore is best-effort, not guaranteed.**
  Measured, not assumed: watcher-timed `kill -9` → permanent residue (V10).
- **P2 — the residue compiles, so the build gate cannot see it.** `go build ./...` → rc=0 on the
  poisoned tree (V11). Any acceptance criterion of the form "`go build` passes" is therefore
  vacuous for this defect.
- **P3 — the poisoned gate accuses an innocent file with a real-violation message.** rc=1,
  `host/store/store.go: forbidden registry/HTTP/cloud dependency "net/http/httputil"` (V11).
- **P4 — the mutation window is real and concurrently observable.** Inherited from iter-55
  (5/90 samples per file against a 0/200 idle control; one real firing through
  `cmd/ailang-worldd/cli_test.go`'s subprocess build) and independently re-confirmed here: the
  V10 kill-timing watcher caught the on-disk mutant from another process on its only attempt.
- **P5 — `go list -deps` honors `-overlay` on go1.25.6, and the overlay closure is byte-identical
  to the on-disk-poison closure.** Directive §4 marked this interaction unverified; it is now
  verified (V7, V8, V12): the overlay run enumerates the same 229 packages as the physically
  poisoned tree — `diff` empty.
- **P6 — `parser.ParseFile(fset, realPath, mutantSrc, ImportsOnly)` parses the mutant bytes while
  attributing positions to the real path** (V14), so the attribution scan can consume the overlay
  and keep its RED message byte-identical.
- **P7 — the per-group `extraForbidden` asymmetry (iter-55, PR #45) is load-bearing** and pinned
  by `TestBareNetHTTPExemptionIsPerGroup` (`allowlist_world_test.go:322`). This design does not
  touch `forbiddenImportPrefixes`, `extraForbidden`, or that test.
- **P8 — the overlay costs nothing measurable**: `go list -deps` 0.038s with and without
  `-overlay` (V9).
- **P9 — the boundary gate is the only test in the repo that writes tracked production files in
  the live tree.** **This premise was BLOCKED as unverified in quorum round 2 and is now measured
  two independent ways.** The original V15 searched `os.WriteFile` alone, while Decision 8's own
  threat model names `os.OpenFile`, `os.Create` and `os.Rename` — `gemini-3-1-pro` caught that the
  audit had not looked for the functions the design calls dangerous, and was right. (a) The widened
  search runs 25 hits across 13 test files with the boundary gate as its known-positive; the sole
  `Create/OpenFile/Rename` outside `host/boundary` is `cli_test.go:272`, which is `t.TempDir()`-
  rooted (V15a). (b) More decisively, the property itself was measured rather than inferred from
  call sites: sampling `git status --porcelain` at 20 Hz across the whole suite yields **1** dirty
  observation with the boundary gate enabled (` M cmd/ailang-worldd/main.go` — the sampler fires)
  and **0** with only `TestWorldBoundaryDependencyAllowlist` skipped (V15b). The `host/capsule`
  `os.WriteFile(fixture.path, …)` is TempDir-rooted via `archive.New(filepath.Join(t.TempDir(),
  "world.db"))` (C2). So the AST guard's `host/boundary`-only scope is honest — and now on
  evidence that no grep could have supplied.
- **P10 — nanosecond `ModTime` advances on every rewrite of the live target, including a
  back-to-back write-and-restore.** Controller-measured this iteration (C1: unchanged in 0/200
  back-to-back trials, 0/200 with a 1 ms gap; single-write control fired 200/200; no-write control
  held equal), on APFS/darwin only. ext4 (CI job 2, `ubuntu-latest`) is a stated assumption, not a
  measurement — which is why the `ModTime` check is a backstop behind a filesystem-independent
  primary, never the guarantee.

### Design Freeze

- After this change, **no code in `host/boundary` writes any byte under the repository root** —
  enforced structurally, not by observation: the confined writer synchronously rejects any
  `repoRoot` destination, and the AST guard forbids bypasses (Decision 7–8). Mutant files and the
  overlay JSON live in `t.TempDir()`.
- The pristine-tree assertions (enumeration non-empty, real tree clean, AIL markers absent) keep
  running against the **real** tree, unchanged. Only the mutation arms move to the overlay.
- The per-group `extraForbidden` asymmetry and `TestBareNetHTTPExemptionIsPerGroup` survive
  verbatim.
- RED messages keep their exact shape: same relative path, same
  `forbidden registry/HTTP/cloud dependency %q` text (guaranteed by P6).
- One file changes: `host/boundary/allowlist_world_test.go`. No production code, no new
  dependencies, no `scripts/*` or `ci.yml` edits.
- "`go build` passes" is explicitly **not** an acceptance criterion (P2).

## Decision — substitute via `go list -overlay` + an overlay-aware read; never write the live tree

Replace `mutateAndRestore` with `mutateViaOverlay`:

1. Read the baseline, assert the `"import (\n"` anchor count is exactly 1, build the mutant
   bytes, assert `sha256(mutant) != sha256(baseline)` — all identical to today (`:271–:277`).
2. Write the mutant to `t.TempDir()` and write an overlay JSON
   `{"Replace": {"<abs repo path>": "<abs mutant path>"}}` beside it. Absolute paths on both
   sides (verified form, V7).
3. `goListDeps` gains an optional overlay argument: when set, the args become
   `["list", "-deps", "-overlay=<json>", patterns...]`. `cmd.Dir` stays `root` — the real tree's
   root from `repoRoot` is now *correct by construction*, because the substitution is declared to
   the toolchain rather than relocated on disk (this dissolves the `runtime.Caller` constraint
   that a copy-based design must fight).
4. The attribution scan in `checkGoGroup` reads file content through an overlay-aware helper: for
   a replaced path it calls `parser.ParseFile(fset, path, mutantBytes, parser.ImportsOnly)`;
   for every other path, exactly today's call. The RED still names the real relative path (P6).
5. `checkAILGroup` takes the same read indirection: a replaced `world/types.ail` resolves to the
   mutant bytes, everything else to `os.ReadFile`. (The AIL arm never used `go list`; only its
   read moves.)
6. **Overlay-consumed control (new, non-vacuous):** after the arm's RED, the arm asserts that
   `goListDeps` *with the overlay* returns a closure containing `group.mutantImport`. Without
   this, a silently broken `-overlay` flag would leave the gate looking toothy — the parse scan
   reds through our own indirection regardless — while the toolchain half of the instrument is
   dead. Measured green-side: the overlay closure contains `net/http/httputil` and grows
   160 → 229 for `host/store` (the import drags in the whole `net/http` subtree, not one
   package), and is diff-identical to the physically poisoned closure (V7, V12).
7. **Confined writer (new — the primary confinement, deterministic by construction, quorum round
   1):** every write the mutation harness performs goes through one writer helper,
   `confinedWrite(dst, data)`. It resolves `dst` to an absolute, symlink-evaluated path and
   **synchronously rejects with an error any destination inside `repoRoot`**, permitting only
   destinations beneath the arm's `t.TempDir()`. A rejection returns before a byte is written —
   there is no window because there is no write, on any filesystem, under any scheduler. A
   deterministic recording-writer test injects a recording implementation and asserts that every
   destination attempted across all arms resolves beneath the temp root and none beneath
   `repoRoot` — and that the recorded set is **non-empty** (each Go arm writes at least the mutant
   file and the overlay JSON), so an arm that recorded nothing reds as vacuous instead of passing.
   Confinement is thereby proven by an exercised rejection path (M3) and a non-empty recorded
   destination set — not by absence of observations.
8. **AST write-guard (new, quorum round 1):** a test in `host/boundary` parses every `.go` file in
   the package (`go/parser` + `go/ast`) and reds if any call to `os.WriteFile`, `os.OpenFile`,
   `os.Create`, or `os.Rename` appears outside the confined writer's single permitted function.
   Deliberately an AST walk over identifiers, **not a grep needle**: this repo's prior textual
   self-guard (`TestSchemaVersionLedgerIsIndependent`, `host/store`) was found vacuous at iter-54
   because its positive needle matched its own check line, and the repair was anchoring; an AST
   walk cannot match its own source text at all, so that failure mode is structurally absent
   rather than patched. Scope is `host/boundary` alone — honest because it is the repo's only
   live-tree mutator, a premise blocked as unverified in quorum round 2 and now measured rather
   than grepped (P9, V15a/V15b, C2).
9. **Deterministic runtime backstop (new, quorum round 1): `ModTime` + sha equality.** Each arm
   captures `os.Stat(livePath)`'s nanosecond `ModTime` and the target's sha256 before `check()`
   and asserts both are **exactly equal** after. No concurrency, no timing window: a
   write-and-restore regression inherently rewrites the file and advances `mtime_ns`
   (controller-measured, C1/P10 — the hardest case, back-to-back write+restore with no gap, was
   detected 200/200). **Stated filesystem assumption:** nanosecond mtime granularity, measured on
   APFS/darwin only; standard on ext4 (CI job 2) but not measured there — which is exactly why
   this is the *backstop* and the writer + AST guard are the guarantee. It exists to catch a
   regression that somehow reaches the live tree by a route the AST deny-list does not name (M7).

### What the gate proves, honestly

Today the teeth-proof is "the checker reds when the **live tree** contains a forbidden import,"
demonstrated by actually poisoning the live tree. After this change it is "the checker reds when
fed the real tree **plus one declared file substitution** applied by the toolchain's own overlay
mechanism." The equivalence of those two inputs is not asserted — it is measured: the overlay
closure and the on-disk-poison closure are byte-identical, 229 packages, `diff` empty (V12), and
the parse scan consumes the identical mutant bytes under the identical path (V14). One residual
gap is stated rather than gated: if a **future** check inside `checkGoGroup`/`checkAILGroup`
reads the disk directly (bypassing both `go list` and the read helper), the mutation arms will not
exercise it with the mutant. That is carried as a code comment on the read helper plus a
review-time rule, not a fake AC — the new AST guard (Decision 8) covers *writes* only; extending
its walker to disk-direct reads inside the checkers is out of 1-day scope (Deferred Scope).

### Alternatives rejected

1. **Mutate a scratch copy of the tree** (the queue row's preferred repair). Also fixes both
   failure modes, but: (i) `repoRoot` bakes the real root in at compile time (`:60–:66`), so every
   helper must be re-threaded with a second root; (ii) the copy would live in `t.TempDir()` —
   `/tmp` on the CI runner — a location-sensitivity class this mission has already paid for;
   (iii) its equivalence claim is "these 141 files are byte-identical," which must itself be
   asserted per run, and is strictly weaker than one declared substitution the toolchain applies;
   (iv) it copies ~3.4 MB per run for no gain. Overlay achieves the same both-modes fix with zero
   copies. This was the fallback had the overlay verification failed; it did not (V7).
2. **Serialise the suite (`-p 1`)**: closes only the concurrency window, leaves the SIGKILL
   residue untouched (the kill lands regardless of what else runs), and slows every
   `go test ./...` invocation. Rejected.
3. **Advisory lock taken by other packages' tests**: participants-only (the launchd watchdog and
   lane deadlines are not participants), leaves the residue untouched, and spreads mechanism
   across many files. Rejected.

## Ordering vs `SM.B2a`

This item **should** land before `SM.B2a`, and was promoted for that reason; it is not a hard
**must** — no `SM.B2a` acceptance criterion depends on it. The two concrete couplings: (a)
`SM.B2a` lengthens the overlap window (~780 LOC of network code onto the 76s-under-`-race`
`host/broker` suite), raising the probability of the latent `cli_test.go` false red; (b) far more
important, the residue's failure message accuses a production file of a **network-boundary
violation** — during the one sprint whose job is to add network code, a residue red is
indistinguishable from `SM.B2a` actually violating the boundary, and would misdirect its executor
and evaluator. Landing this first removes that entire confusion class for the cost of one day.

## Files to Create/Modify

- `host/boundary/allowlist_world_test.go` (~+150/−55 LOC) — replace `mutateAndRestore` with
  `mutateViaOverlay` (all writes routed through the confined writer); add the confined writer and
  its deterministic recording-writer test; add the AST write-guard test; add the overlay parameter
  to `goListDeps`; add the overlay-aware read helper used by `checkGoGroup`/`checkAILGroup`; add
  the per-arm `ModTime`+sha backstop; add the overlay-consumed closure assertion. All other tests
  in the file are untouched. (Budget note: the LOC estimate grew ~+55 versus the pre-quorum draft,
  but every added line is straight-line synchronous code replacing concurrency plumbing — still
  one test file, still ≤1 day; nothing was cut from Deferred Scope to absorb it.)

No other files change. `scripts/verify_go.sh`, `scripts/verify_ail.sh`, and `.github/workflows/ci.yml`
run unmodified.

## Conflict Surface

- `scripts/verify_go.sh:113` — the `-race` leg SIGKILLs its **process group** at 600s (V17). Today
  that kill can land mid-mutation and poison the tree from inside the repo's own gate; after this
  change there is nothing on disk to poison. No edit needed — the script benefits without change.
- `.github/workflows/ci.yml` — `go-verify` invokes `verify_go.sh`; `ailang-verify` runs
  `verify_ail.sh` over `world/*.ail`, which this design stops writing entirely. Neither job
  changes.
- `cmd/ailang-worldd/cli_test.go:125` (V21) — builds a subprocess binary from live source; the one
  recorded real firing of the window red-lit `TestCLIRealSubprocessEpisode` through this site.
  After this change the window does not exist. No edit.
- `host/daemon/daemon_test.go:700–768` (V19) — an independent, read-only `go list -deps`
  module-root allowlist over `./host/daemon/... ./cmd/ailang-worldd/...`. It is a window *reader*
  today, but its dot-rule classifies `net/http/httputil` as stdlib → allowed, so it was never a
  false-red victim of this particular mutant. Unaffected either way once the window is gone.
- `TestBareNetHTTPExemptionIsPerGroup` + per-group `extraForbidden` (iter-55) — untouched and must
  stay green; collapsing the asymmetry already fails that test by design.
- `-race` legs — this design now adds **no concurrency of any kind**: the quorum-blocked polling
  detector is deleted; the confined writer is synchronous, the backstop is one `os.Stat`+sha pair
  per arm, and the AST guard is a single-threaded parse of one package's files. The boundary
  package remains ~0.4s (V13), so the 600s `-race` budget is unaffected.
- Worktrees/CI checkouts — overlay JSON paths are derived from `repoRoot` at runtime as absolute
  paths, so the mechanism is checkout-location-independent (verified from a detached worktree,
  V7).

## Systemic-Issue Audit

Is live-tree mutation a pattern? **No — this is the only instance, and that is now a measurement
rather than a search.** The original audit here grepped `os.WriteFile` alone; `gemini-3-1-pro`
blocked it in quorum round 2 for not searching `os.OpenFile`/`os.Create`/`os.Rename` — the very
functions Decision 8 names as threats — and the objection was accepted rather than argued. Two
replacements: the widened search returns 25 hits across 13 test files with `host/boundary` as its
known-positive, and its only `Create/OpenFile/Rename` outside `host/boundary` is
`cli_test.go:272`, `t.TempDir()`-rooted (V15a); and the property itself was measured end-to-end by
sampling `git status --porcelain` at 20 Hz across the whole suite — **1** dirty observation with
the boundary gate enabled, **0** with only its mutation test skipped (V15b). The one suspicious
production-path write elsewhere (`host/capsule/capsule_test.go:97`) targets a `t.TempDir()`-rooted
archive copy, verified by reading `archiveExecutable` (C2). So the fix is correctly local.

One adjacent latent gap was found while reading `checkGoGroup` and is **deliberately not fixed
here** (see `10/OD-1`): the `go list -deps` closure is asserted non-empty and logged, but **never
scanned against `forbiddenImportPrefixes`** — the forbidden check runs only on the group's own
files' *direct* imports (`:138–:162`). To say it plainly, because a reader of a "dependency
allowlist" gate would reasonably assume otherwise: **the enumeration's only current role is
anti-vacuity.** `checkGoGroup` (`:130`) calls `goListDeps` and uses the result solely for the
`len(deps) == 0` guard (`:135`); the returned closure is never compared against
`forbiddenImportPrefixes`/`extraForbidden` (controller re-confirmed first-party this iteration).
A transitive route (a protected package importing an in-repo helper that imports
`net/http/httputil`) would pass today. ~~Feasibility baseline for a future scan is measured: 0
forbidden-prefix hits in all three closures~~ — **REFUTED at iter-57 (V16a/V16b), and the
refutation makes the transitive route *actual* rather than hypothetical.** `cmd/ailang-worldd`'s
closure contains `host/registry`, which IS a forbidden prefix (`:53`), reached via
`host/daemon:51`. So a closure scan would **red today** — but on **legitimate** code: `host/registry`
is the *interpreter epoch* registry, not the *package* registry the entry targets, exactly the name
collision iter-53 predicted and iter-57 measured. **`10/OD-1` therefore cannot be implemented until
that collision is resolved**, which strengthens the deferral below rather than weakening it. Bare
`net/http` still appears only in `cmd/ailang-worldd`'s (1), matching its documented loopback
exception (V16).

## Deferred Scope (cut to hold ≤1 day)

- Closure scanning against forbidden prefixes (`10/OD-1`) — changes what the gate *enforces*, not
  how it mutates; deserves its own measured design.
- A permanent CI kill-harness test (`10/OD-2`) — the kill demonstration runs as sprint evidence
  instead.
- An executable guard against future disk-direct *reads* inside the checkers — the AST
  write-guard (Decision 8) supplies the walker; extending its deny-list to `os.ReadFile`-class
  calls outside the read helper is a small follow-up, carried here as a comment + review rule to
  hold ≤1 day.
- Any change to `scripts/verify_go.sh` or CI.

## Acceptance Criteria

Each AC carries the mission's vacuity self-test: *would it pass identically if the protected thing
did not work?*

> **⚠ SPRINT-PLANNED AT ITER-57 — THE PLAN SUPERSEDES SEVERAL ACs BELOW, AND THE PLAN WINS ON
> THOSE POINTS.** Plan: `.ailang/state/sprints/w-boundary-gate-tree-mutation.{plan.json,handoff.md}`
> (planner `opus`, lane derived **fail-closed `opus missing-script`** — `tools/launchd/derive-planner-lane.sh`
> is absent from this checkout; `MISSION_PLANNER_MODEL=opus` independently agrees). Milestones
> **BG.A** (AC2, AC3, AC4, AC5 · M1, M2, M4, M5) → **BG.B** (AC1a · M3, M6) → **BG.C**
> (AC1b, AC6′ · M7); the partition is complete — 7 criteria, 7 mutations, none dropped or
> double-assigned. Applying this document's own vacuity self-test, the planner judged **four** of
> the six ACs below vacuity-capable as written, and rewrote them:
>
> - **AC6 → AC6′.** The `≤2× 0.435s` bound is **superseded**. The constant was transcribed from a
>   different worktree at a different cache warmth, so it does not survive rule 3e(i) ("a control is
>   only a control if it runs from a tree in the state the baseline was in"). Controller measurement
>   on **unchanged** code at `e9c8c85`: fresh-worktree first run **0.664s / 0.621s** (n=2, two
>   independent worktrees), warm steady state **~0.480s** (n=9). So zero code change already sits at
>   **1.43–1.53×** of the AC's own constant cold and 1.09–1.17× warm — leaving ~1.31× of headroom,
>   not 2×. Worse than the false-red risk: the noise band consumes ~76% of the budget, so a **green**
>   AC6 could not have failed informatively. The planner found a third defect the controller missed —
>   the units are ambiguous by **1.32×** (go-*reported* 0.479s vs wall-clock 0.631s median, n=5;
>   0.435s is a go-reported figure but the wording is about the command completing) — and a fourth:
>   what AC6 nominally protects is a **600s** `-race` budget against a **0.5s** package, 1200× of
>   headroom, so it could not fail for the change either. **AC6′** is a paired same-session ratio
>   (one worktree, equal warmth, ≥8 interleaved A/B pairs swapping only this file and asserting its
>   sha256 changed; `median(wall_B)/median(wall_A) ≤ 1.50`) plus a `median(wall_B) ≤ 3.0s` ceiling,
>   **asserted on wall-clock**, said explicitly so the unit cannot be re-ambiguated. Its noise floor
>   was measured rather than assumed: 8 interleaved pairs on **identical** code returned **1.0079**
>   against a true ratio of 1.000, pooled spread 1.058 — so 1.50 is ~8.5× the spread.
> - **AC1 → AC1a + AC1b.** The AST guard is vacuity-capable in a mode Decision 8 never considered:
>   `host/boundary` holds **exactly one `.go` file** and it is a `_test.go`, so an empty `ParseDir`,
>   a filter dropping `_test.go`, or a selector bug each yield "zero violations, green". Decision 8
>   defends the *self-match* mode at length and is silent on the *empty-enumeration* mode. The plan
>   requires an exact non-empty file count, a known-positive (the walker must **find** the permitted
>   `os.WriteFile` and report its line), and deny-list completeness. The `ModTime` backstop is
>   silently disarmable on CI (see V16c and the ext4 note in C1) and is made **fail-loud**: a 20/20
>   granularity probe whose failure is a **test failure** naming both `st_dev`s, plus sha256, size,
>   mode and **inode** asserted unconditionally — inode closing a rename route that both of this
>   document's stated observables miss.
> - **AC2 strengthened.** Measured by the planner: `go list -overlay` **silently ignores** a
>   `Replace` key matching no file — rc=0, base closure, **no stderr**. Asserting "the overlay
>   closure contains `mutantImport`" is therefore not enough on its own; the plan asserts on the
>   closure `checkGoGroup` actually **consumed**, plus a negative half. Free — `checkGoGroup` already
>   computes and discards it.
> - **AC4 instrument fix.** `git status --porcelain` reports **untracked** files, so an in-tree ready
>   marker reds on the harness's own artifact. The marker moves outside `repoRoot` and the residue
>   assertion is path-scoped to the four targets. (The round-2 fail-closed rewrite of AC4 itself
>   stands — the planner called it this document's best work.)
>
> AC3 and AC5 are sound as written; the controller separately verified AC3's `:217–:222` citation
> first-party (it is exactly the two RED-fidelity assertions — the doc never checked it). Estimate:
> the doc's **≤1 day of effort** holds (velocity: 12 landed feat/fix commits, median **363**
> insertions; the closest analogue is `1761a9c`, a single-test-file change to *this same file*,
> +56/−7, one iteration), but **elapsed is 2–3 iterations**, because measured cadence is ≤1
> milestone/iteration and **4 of the 7 mutations cannot run in the executor sandbox** (M5 needs
> subprocess SIGKILL + git inspection; M6/M7 re-arm live-tree writes; AC6′ needs a file swap out of
> git history) — the controller pass is the critical path. LOC re-estimated **+150 → +250**, all of
> it spent making this document's own ACs non-vacuous.

> **⚠ `BG.A` LANDED AT ITER-58 (PR #47 → squash `278f102`, dev CI: `go host build + test gate`
> SUCCESS SHA-addressed on the merge commit; `ailang-code verify gate` blocked by a **declared
> GitHub Actions outage** — see the charter STATUS stamp). Executor `opus` (the driver env's
> `MISSION_EXECUTOR_MODEL`, NOT the plan's assumed `codex:gpt-5.6-sol` — so the plan's `S-7`
> no-git-writes/snapshot rule did not apply and BG.A is one ordinary commit). Evaluator `sonnet`
> **PASS 89/100, round 1, zero blocking findings**.** **AC2, AC3, AC4 and AC5 are DISCHARGED**;
> `AC1a`/`AC1b`/`AC6′` and mutations `M3`/`M6`/`M7` remain with `BG.B`/`BG.C`. Measured, not
> inherited: AC2's four numbers per Go arm are `host/store` **160/0 → 229/1**, `host/replay`
> 162/0 → 231/1, `cmd/ailang-worldd` 233/0 → 234/1 — the planner's `PV-3`/`PV-8` prediction exactly,
> asserted on the closure `checkGoGroup` RETURNED rather than on a second `go list`. `M1` and
> `M2(b)` were re-run first-party by the controller under the house recipe (anchor count **1**,
> differing sha256 asserted before believing, control arm first, byte-identical restore verified):
> `M2(b)` — the overlay `Replace` KEY naming no real file, i.e. the **silent** failure — reds with
> `overlay closure=160, baseline closure=160 -- the toolchain half of the gate is dead`. **`M5`, the
> AC4 kill harness, was run by the controller on all four arms with its NEGATIVE CONTROL (rule 3d)**:
> marker awaited under a fixed timeout → mutant + overlay JSON verified present → overlay verified to
> map real target → temp mutant → process verified **alive** → `SIGKILL` (`rc=-9`); result **0**
> changed target sha256s and **0** `git status --porcelain` lines on all four arms, against the SAME
> kill on the **base** harness which returned `killed_while_mutating=host/store/store.go`,
> `RESIDUE=YES`, ` M host/store/store.go`. Outcomes differ, so the green measures the mechanism and
> not the environment. AC4's fail-closed property was proven in both directions: armed-but-never-
> killed → `panic: test timed out`, rc=1 (a timeout FAILS); an in-repo marker path → rejected with
> `resolves inside repoRoot`, and the marker file was never created.
>
> **ONE DEFECT IN THIS DOCUMENT AND ITS PLAN, FOUND DURING EXECUTION AND NOT BEFORE.** `go/parser`'s
> `readSource` tests `src != nil` on the **interface**, so a typed nil `[]byte` is a NON-nil
> interface and is handed back as an **empty source** — every unreplaced file then parses as
> `expected 'package', found 'EOF'`, and *a checker that cannot read the tree finds no forbidden
> imports*. Both this doc and the plan write the read helper as
> `parser.ParseFile(fset, path, <bytes-or-nil>, …)`, which is the shape that produces it. It was
> observed live, is isolated in `parseSrc`, and is recorded here because a future implementer
> following the written wording will hit it. It is this item's own spine arriving inside the fix:
> a green that means "the check never ran".
>
> **CARRY TO `BG.B` — the plan's write-site count is now WRONG BY ONE, and the missing one would red
> the guard `BG.B` installs.** The plan says "route BG.A's **two** per-arm writes (mutant file,
> overlay JSON) through `confinedWrite`"; it was written before the AC4 barrier existed as code, and
> the barrier adds a **third** direct `os.WriteFile(absMarker, …)`. Measured at `278f102`: **3**
> `os.WriteFile` sites (`:383` marker, `:428` mutant, `:439` overlay JSON), **0**
> `OpenFile`/`Create`/`Rename`, with a firing known-positive control (`os.ReadFile` = **4**).
> Decision 8's guard reds on any of the four names outside the single permitted site, so leaving the
> marker write direct makes `BG.B` red on `BG.A`'s own landed code. Routing it through
> `confinedWrite` is CORRECT rather than an exemption — the marker is *required* to resolve outside
> `repoRoot`, which is exactly what `confinedWrite` permits, so the confined writer also becomes the
> enforcement point for the AC4 marker-path rule and replaces the bespoke `insideRepo` check at
> `:367–:373`. Raised by the evaluator (`NB-2`), reproduced by the controller before being recorded;
> the plan artifact is corrected in place and carries a `controller_corrections` entry.
>
> **A DELIBERATE DEVIATION FROM THE PLAN'S LITERAL SIGNATURE, AND IT IS THE RIGHT ONE.** The plan
> specifies `checkGoGroup(root, group, overlay string)` — a single JSON path serving both halves.
> The executor used a two-field `overlay{jsonPath, replace}` so the toolchain half and the read half
> are **separately disarmable**, and the reason is `AC2`'s own falsifiability: with one string,
> dropping `-overlay` would also disarm the import scan, so `M2(a)` would red at **AC3** instead of
> at AC2 and the toolchain half would go untested — *a mutation shaped to the check tests the check,
> not the threat* (iter-54's spine). The observed `M2(a)`/`M2(b)` messages confirm it. The evaluator
> scored this **−5 on design fidelity** as an undocumented departure; the controller records the
> opposite verdict on the merits — the deviation is what makes AC2 non-vacuous, it is documented in
> a code comment on the type, and the plan is what is now corrected.

- **AC1 — all mutation-harness writes are structurally confined; a live-tree write is
  synchronously rejected.** Every write the mutation harness performs is routed through the
  confined writer; a write to any destination inside `repoRoot` is rejected synchronously with an
  error before a byte lands — deterministic, because there is no window when there is no write.
  The recording-writer test observes only temp-root destinations across every arm, and the AST
  guard rejects direct `os.WriteFile`/`os.OpenFile`/`os.Create`/`os.Rename` bypasses in
  `host/boundary` outside the writer. As a deterministic runtime backstop, each arm asserts the
  live target's nanosecond `ModTime` and sha256 are exactly equal before/after `check()` (C1/P10;
  APFS-measured, ext4 granularity is a stated assumption). *Not vacuous:* routing today's
  live-path write through the writer is rejected immediately (M3); a writer-bypassing direct
  write reds both the AST guard and the backstop (M6); a write-and-restore via a mechanism off
  the AST deny-list reds the backstop alone (M7). The green state is proven by the exercised
  rejection path and the recorded destination set — not by "zero observations," which was the
  defect the quorum blocked.
- **AC2 — the overlay demonstrably reaches `go list`.** Each Go arm asserts the overlay closure
  contains `group.mutantImport`. *Not vacuous:* removing `-overlay` (or pointing it at `{}`) makes
  this assertion red (M2) while the parse-scan RED would otherwise mask the breakage.
- **AC3 — RED fidelity is preserved.** Each arm still asserts the guard reds *and* that the error
  names the exact relative path (existing `:217–:222` assertions retained verbatim over the new
  mechanism). *Not vacuous:* a mutant with a benign import makes the arm red at
  "mutation … passed boundary guard" (M1), and the AIL arm likewise (M4).
- **AC4 — the threat is replayed against a PROVEN-ARMED target, via a bounded synchronization
  protocol rather than randomized timing.** *(Rewritten by the round-2 carve-out, applying
  `gpt5-6-sol`'s proposed fix verbatim — see "Quorum round 2" below.)* Each mutation arm carries a
  test-only, environment-controlled barrier immediately **after** the mutant and overlay JSON are
  written and validated, but **before** the checker consumes them. Per arm, the external harness
  must: select the arm deterministically; wait with a fixed timeout for that arm's ready marker;
  **verify the mutant and overlay artifacts exist and that the overlay maps the real target to the
  temp mutant**; only then send `SIGKILL`; then assert the tracked target hashes and
  `git status --porcelain` are unchanged. **A missing marker, a timeout, an absent artifact, or an
  early process exit must FAIL the criterion** — never pass it. Record **one run per arm**, not ten
  randomized kills. *Not vacuous:* under the old formulation every kill could land before any arm
  armed, and "0 dirty lines" would then mean *the threat was never exercised* — this item's own
  defect class. The barrier makes "the kill landed after the overlay was armed" an asserted
  precondition rather than a probabilistic hope. The historical 1/1 residue observation on the OLD
  implementation (V10–V11) is retained as background and is explicitly **not** treated as a per-run
  known-positive control for the new harness.
- **AC5 — the iter-55 asymmetry survives.** `TestBareNetHTTPExemptionIsPerGroup` and all
  `TestWorldBoundaryNullCases` pass unchanged; `extraForbidden` fields are not edited. *Not
  vacuous:* that test fails by construction if the asymmetry is collapsed (landed teeth, PR #45).
- **AC6 — cost stays flat.** `go test ./host/boundary/ -count=1` completes in ≤2× the measured
  0.435s green baseline (V13), keeping the `-race` budget untouched. The quorum-blocked polling
  detector — the main runtime threat to this bound — is deleted; what remains is one synchronous
  path check per write, one `os.Stat`+sha pair per arm, and one single-package AST parse in the
  guard test, all sub-millisecond-class. *Not vacuous:* a design that, e.g., re-enumerated the
  dependency closure per write or AST-parsed the whole repo in the guard would blow this bound.

Explicitly rejected as an AC: "`go build ./...` passes" — measured rc=0 on the poisoned tree
(V11); it protects nothing here.

## Non-Vacuity — named RED mutation for every added assertion

| # | Exact edit | Expected RED | Shape |
|---|---|---|---|
| M1 | In the mutant generator, replace `group.mutantImport` with `"fmt"` for one arm | that arm: `mutation in host/store/store.go passed boundary guard` | check-shaped: proves the guard, not the harness, produces the RED |
| M2 | Drop `-overlay=<json>` from the arm's `goListDeps` call (leave everything else) | AC2's closure assertion: `overlay closure does not contain "net/http/httputil"` | check-shaped: proves the overlay-consumed control is live |
| M3 | Route today's live-path mutation through the confined writer: `confinedWrite(livePath, mutant)` before `check()` | immediate synchronous rejection: `confined writer: destination inside repoRoot: host/store/store.go` — the arm reds before any byte is written | **threat-shaped + deterministic**: today's defect aimed at the writer; no observation window exists because no write occurs |
| M4 | Point the AIL arm's override at the *baseline* bytes of `world/types.ail` | `mutation in world/types.ail passed boundary guard` | check-shaped, AIL arm parity |
| M5 | No code edit — the AC4 barrier protocol, run **once per arm**: wait for that arm's ready marker under a fixed timeout, assert the mutant + overlay JSON exist and that the overlay maps the real target to the temp mutant, `SIGKILL`, then compare tracked hashes and `git status --porcelain` | on **current** code: residue (`M host/store/store.go`, measured 1/1, V10); on new code: `git status --porcelain` = 0 lines for every arm. **A missing marker / timeout / absent artifact / early exit FAILS the criterion rather than passing it** | **threat-shaped + armed-target-proven**: the crash replayed against a target proven to be in the protected window, so a green cannot mean "never exercised" |
| M6 | Reintroduce direct `os.WriteFile(livePath, mutant, perm)` + deferred restore around `check()` (today's exact code shape, bypassing the writer) | the AST guard: `direct os.WriteFile outside confinedWrite at allowlist_world_test.go:<line>`; at runtime the backstop also reds — `mtime_ns` advanced (C1: single write fired 200/200) | **threat-shaped**: the current defect verbatim, caught statically and dynamically |
| M7 | Write-and-restore the live target via a mechanism off the AST deny-list (an `exec.Command("cp", …)` pair around `check()`) | the backstop alone: `live target mtime_ns changed across arm` (C1: back-to-back write+restore left `ModTime` unchanged in 0/200 trials) | check-shaped: proves the backstop is live independently of the AST guard |

Green control for all of the above: unmutated `GOTOOLCHAIN=go1.25.6 go test ./host/boundary/
-count=1` → `ok … 0.435s`, all three baseline shas unchanged, `git status --porcelain` → 0 lines
(V13).

## Open Decisions

- **`10/OD-1` — should the dependency closure itself be scanned against
  `forbiddenImportPrefixes`/`extraForbidden`?** Today only direct imports are checked; a
  transitive forbidden dependency is invisible (measured feasibility baseline: a scan would be
  green on the current tree, V16). **Controller default if no human answers: out of scope here;
  file as a follow-up queue note.** It changes enforcement semantics (interacts with the
  `extraForbidden` asymmetry and the `cmd/ailang-worldd` loopback exception) and deserves its own
  design with its own mutations.
- **`10/OD-2` — does the AC4/M5 kill-harness become a permanent CI test or sprint evidence
  only?** **Controller default: sprint evidence only**, recorded in the sprint's verification log.
  A permanent subprocess-kill test adds timing-dependent machinery to CI to re-prove a property
  that is structural after this change (no write exists to interrupt) and is already pinned
  deterministically by AC1's writer rejection, AST guard, and `ModTime`+sha backstop.

## Verification Log

All rows measured by the designer at `origin/dev` @ `deeb804` in the detached worktree
`.wt-iter56-design` (plus a sibling scratch dir `.wt-iter56-overlay-scratch`, since removed),
shell `zsh`, `GOTOOLCHAIN=go1.25.6` exported for every `go` invocation. KP = known-positive
control carried in the same call.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Worktree is `deeb804`, clean | `git rev-parse HEAD && git status --porcelain \| wc -l` | `deeb804e458f…`, `0` |
| V2 | Gate file is 351 lines; helpers at `repoRoot:60`, `goListDeps:72`, `enumerateAIL:96`, `checkGoGroup:130`, `checkAILGroup:165`, `mutateAndRestore:190` | `wc -l` + `grep -n 'func …' host/boundary/allowlist_world_test.go` | `351`; line numbers as listed (refines directive §2/§3e, which cited `:70`/`:93`) |
| V3 | `go.mod` has zero `replace` directives (KP: `module` line present) | `grep -c '=>' go.mod; grep -c '^module' go.mod` | `0` (grep rc=1); KP `1` (`module github.com/sunholo-data/ailang-world`, `go 1.25.6`) |
| V4 | Repo is 141 tracked files / 3,385,213 tracked bytes | `git ls-files \| wc -l; git ls-files -z \| xargs -0 stat -f%z \| awk '{s+=$1} END {print s}'` | `141`, `3385213` |
| V5 | Baseline closures 160/162/233; `net/http/httputil` in none (KP: `modernc.org/sqlite` in each) | `go list -deps <p> \| wc -l` + `grep -c` per pattern | `./host/store/… total=160 httputil=0 sqlite_control=1`; replay `162/0/1`; worldd `233/0/1` |
| V6 | Gate-identical mutants generated: anchor `"import (\n"` count 1 in `store.go` and `main.go`; shas differ; baselines match directive | python: count anchor, sha256 both | `store.go: anchor_count=1 baseline_sha=40315426a1760b9a mutant_sha=d08c3f462260d3f5`; `main.go: anchor_count=1 baseline_sha=118440f9c026de6e mutant_sha=ea3d53ff1dc5c466` |
| V7 | **`go list -deps` honors `-overlay`** (absolute-path JSON, run from repo root): mutant visible, live tree untouched (KP pair: without overlay 0, with overlay 1) | `go list -deps ./host/store/... \| grep -c httputil` then same with `-overlay=…/overlay.json`; `shasum`; `git status --porcelain \| wc -l` | without: `0` (rc=1); with: rc=0, httputil `1`, closure `229` lines; live shas `40315426…`/`118440f9…` unchanged; git `0` |
| V8 | Overlay also works for the `cmd/ailang-worldd` group | `go list -deps -overlay=… ./cmd/ailang-worldd/... \| grep -c httputil` | `1` |
| V9 | Overlay adds no measurable latency | `time go list -deps ./host/store/...` with/without `-overlay` | `0.038 total` both |
| V10 | **SIGKILL mid-mutation leaves permanent residue** (re-run of directive §3b, first-party) | `go test -c ./host/boundary/`; a python kill-timing watcher hashes `store.go` in a tight loop, `kill -9` on first diff | `killed_while_mutating=host/store/store.go`, `test_binary_rc=-9`, `store.go_sha_after=d08c3f462260d3f5` = the gate's own mutant sha, `residue=YES` (1/1 attempt) |
| V11 | Poisoned tree: dirty to git, **invisible to `go build`**, and the gate accuses the innocent file | `git status --porcelain`; `go build ./...; echo rc`; `go test ./host/boundary/ -run TestWorldBoundaryDependencyAllowlist -count=1` | ` M host/store/store.go`; `build_rc=0`; `gate_rc=1` with `host/store/store.go: forbidden registry/HTTP/cloud dependency "net/http/httputil"` |
| V12 | **Overlay closure ≡ on-disk-poison closure** | `go list -deps ./host/store/...` on poisoned tree → file; `diff` vs overlay run's output | `IDENTICAL closures (poisoned-on-disk vs overlay)` — diff empty, both 229 packages |
| V13 | Restore + green control: byte-identity, clean tree, whole boundary package green | `cp` saved baseline; `shasum`; `git status --porcelain \| wc -l`; `go test ./host/boundary/ -count=1` | `40315426a1760b9a`; `0`; `ok github.com/sunholo-data/ailang-world/host/boundary 0.435s` |
| V14 | `parser.ParseFile(fset, realPath, mutantSrc, ImportsOnly)` parses the mutant under the real path | 20-line probe program parsing the V6 mutant bytes under the live `store.go` path | first import printed: `"net/http/httputil"` |
| V15 | ~~The boundary gate is the only live-tree mutator~~ — **SUPERSEDED by V15a/V15b. The original row searched `os.WriteFile` ONLY, while this design's own threat model (Decision 8) names `os.OpenFile`, `os.Create` and `os.Rename` as mutators — so it did not search for the functions the design itself calls dangerous.** Objection raised by `gemini-3-1-pro` in quorum round 2 and **correct**; retained here rather than silently rewritten, because "the audit did not look for the thing it was auditing" is the same class this document exists to fix | (original) `grep -rn os.WriteFile --include='*_test.go' …` | narrow — conclusion happened to survive, but was **unverified** as written |
| V15a | **Gemini's widened search, run verbatim by the controller** (KP: the known mutator must appear) | `grep -rnE 'os\.(WriteFile\|OpenFile\|Create\|Rename\|Remove)\|ioutil\.WriteFile' --include='*_test.go' host/ cmd/` | **25 hits across 13 test files**; KP `host/boundary/` = **3** (instrument sees the known positive). The only `Create/OpenFile/Rename` outside `host/boundary` is `cmd/ailang-worldd/cli_test.go:272`, whose `path := filepath.Join(t.TempDir(), "large.json")` is TempDir-rooted (`:271`) |
| V15b | **P9 measured directly rather than inferred from 25 call sites** — the whole suite mutates no tracked file except via the boundary gate. *(The 20 Hz poll below is a one-off external MEASUREMENT instrument run by the controller; it is not the per-arm polling detector deleted in round 1, and nothing like it is proposed for the committed test.)* | two arms in a clean worktree, `git status --porcelain` sampled at 20 Hz throughout: **KP control** `go test ./... -count=1`; **arm** the same with `-skip 'TestWorldBoundaryDependencyAllowlist'` | KP control: **1** distinct dirty observation — ` M cmd/ailang-worldd/main.go` — so the sampler fires; arm: **0** distinct dirty observations. Both arms end with a clean tree. Both also exit `rc=1` on the **same** pre-existing cause (`TestEpisodeLiveReplayThreeArmsAndEvidence`: `AILANG_BIN must name the pinned released interpreter` — `verify_go.sh`'s deliberate fail-loud guard, which a bare `go test` bypasses); identical in both arms, so the comparison is unaffected, and CI is green at `deeb804`. **P9 is therefore a measurement, not a grep** |
| V16 | ~~No forbidden prefix in any baseline closure~~ — **the first half is REFUTED at iter-57; see V16a.** The bare-`net/http` half stands | `go list -deps <p> \| grep -cE '<forbidden prefixes>'` and `grep -c '^net/http$'` per group | store `0/0`, replay `0/0`, worldd `0/1` — the **second** number of each pair is right, the **first is wrong**. The error hid because one cell bundled two different questions and only one of them was ever checked against a firing control |
| V16a | **`cmd/ailang-worldd`'s baseline closure DOES contain a forbidden prefix** — `github.com/sunholo-data/ailang-world/host/registry`, which is `forbiddenImportPrefixes[3]` (`:53`) — reached transitively via `cmd/ailang-worldd → host/daemon → host/registry` (`daemon.go:51`, a direct import). Raised by the **iter-57 planner**, reproduced first-party by the controller before being recorded | `go list -deps <p> \| grep -c '^github.com/sunholo-data/ailang-world/host/registry$'` per group, with KP `modernc.org/sqlite` in the same call | worldd **1** (of 233), store **0** (of 160), replay **0** (of 162); KP fired **1 / 1 / 1**. **So this document's "0 forbidden-prefix hits in all three closures" and "a scan would be green on the current tree" are FALSE** |
| V16b | **The red that `10/OD-1` would produce is a FALSE POSITIVE — which makes the deferral STRONGER, not weaker.** `host/registry` is the **interpreter epoch** registry (`world/epoch-registry/v1`, `w-world-library-m1` Decision 5), unrelated to the *package* registry the forbidden entry targets. iter-53 predicted exactly this (*"a name collision that will produce a false positive the moment anything legitimately needs epoch metadata"*); at iter-57 it is **measured**, and `host/daemon` is the thing that legitimately needs it | read `host/registry/registry.go:1–12`; `grep -n 'host/registry' host/daemon/*.go` (KP: `daemon.go` carries 3 in-repo `host/` imports) | package doc confirms the epoch registry; `daemon.go:51` is a direct import. **`10/OD-1` cannot be implemented until the collision is resolved** — a closure scan today reds legitimate code, so shipping it as-is would install a gate whose red means nothing |
| V16c | **Nothing a PASSING Go test emits reaches CI.** The gate's `ENUMERATION`/`MUTATION`/`RESTORE` `t.Logf` lines have therefore **never appeared in a CI log**, and "loud but non-gating" is a contradiction in CI. Raised by the **iter-57 planner**, reproduced first-party | paired arms on `TestWorldBoundaryDependencyAllowlist`, same worktree: **A** = CI's exact form (no `-v`); **B** = identical **+ `-v`**, as the known-positive | A: rc=0, output is the single line `ok … 0.580s`, matching lines **0**. B (KP): rc=0, matching lines **12** (the ENUMERATION rows dump the full 160/162/233 closures). `verify_go.sh:100` is `go test ./... -count=1`, and its `-race` leg builds `["go","test","./...","-count=1","-race","-timeout","8m"]` — **neither carries `-v`**. Consequence for this design: any observable it wants CI to see must be an **assertion**, never a log line |
| V17 | The repo's own gate SIGKILLs a process group on `-race` timeout | `grep -n 'killpg\|SIGKILL' scripts/verify_go.sh` | `113: os.killpg(os.getpgid(p.pid), signal.SIGKILL)` |
| V18 | `replay.go` has the same unique anchor; the AIL enumeration is 4 files incl. `world/types.ail` | python anchor count; `ls world/*.ail \| wc -l; ls world/types.ail` | `1`; `4`; `world/types.ail` |
| V19 | `host/daemon`'s go-list allowlist is read-only and dot-rule-classifies `net/http/httputil` as stdlib (never a victim of this mutant) | read `daemon_test.go:700–768` (`allowedDepModules`, `isStdlibImportPath`, `disallowedDeps`) | module-root allowlist, stdlib exempt via first-segment-has-no-dot; no writes |
| V20 | No `go.work` (KP: `go.mod` present in same call) | `ls go.work; ls go.mod` | `No such file or directory`; `go.mod` |
| V21 | The subprocess source build is at `cli_test.go:125` | `grep -n 'build := exec.CommandContext' cmd/ailang-worldd/cli_test.go` | `125:` (refines directive's `:128`, which is inside the same function) |
| V22 | **All three mutants applied simultaneously compile and vet clean** — the charter/queue/dashboard description "deliberately NON-COMPILING import" is false (KP pair: the boundary gate itself reds on the same poisoned tree; restore is byte-identical ×3, tree clean) | python applies the gate-identical mutant (`:276` shape, anchor count asserted =1, shas asserted differing) to all three `mutantFile`s at once; `go build ./...`; `go vet ./cmd/... ./host/store/... ./host/replay/...`; gate run; restore + sha check + `git status --porcelain` | `build_rc=0`, `vet_rc=0`; KP `gate_rc=1` with `host/store/store.go: forbidden registry/HTTP/cloud dependency "net/http/httputil"`; restored shas `40315426…`/`8d8d9fef…`/`118440f9…` byte-identical; porcelain = only this untracked design doc |

### Controller-provided rows (provenance: controller measurement, this iteration — NOT designer rows)

| # | Claim | Method | Observed |
|---|---|---|---|
| C1 | Nanosecond `ModTime` deterministically detects a write-and-restore of the live target (`os.stat().st_mtime_ns`, APFS, darwin/arm64) | 200-trial arms plus KP controls in the same run | back-to-back write+restore: `ModTime` unchanged in **0/200**; with a 1 ms gap: **0/200**; no-write control: unchanged = True (comparator works); single-write control: changed **200/200** (detector fires). **Not measured on ext4** (CI job 2, `ubuntu-latest`) — carried as AC1's stated filesystem assumption, and the reason the `ModTime` check is a backstop, not the primary |
| C2 | V15 re-verified: `host/boundary/allowlist_world_test.go` is the repo's ONLY live-tree mutator; `host/capsule`'s `fixture.path` resolves under `archive.New(filepath.Join(t.TempDir(), "world.db"))` | controller first-party re-read | confirmed — the AST guard's `host/boundary`-only scope is honest |

**Correction to mission records (V22):** the charter queue row and the dashboard both describe
the mutants as "a deliberately NON-COMPILING `net/http/httputil` import". That is false and now
measured false first-party: with all three mutants applied simultaneously, `go build ./...` →
rc=0 and `go vet` over all three groups → rc=0 — the gate's own inline comment already says
"compiling HTTP import" (`allowlist_world_test.go:276`). This matters to the harm model: the
residue is **invisible to the build**, which is worse than a build break (P2). It also means
iter-55's transient `could not import net/http/httputil` CI failure cannot be explained by "the
file is broken" — it is a build-graph/timing artifact of the file *changing* under a concurrent
build. No sentence in this design describes the mutant as non-compiling; the landing sprint
should carry this wording correction to the charter/queue/dashboard as bookkeeping.

**Refinements to the iter-56 directive found by these measurements** (none change its
conclusions): `goListDeps` is at `:72` not `:70` and `enumerateAIL` at `:96` not `:93` (V2); the
subprocess build is at `cli_test.go:125` not `:128` (V21); a poisoned/overlaid `host/store`
closure is **229** packages, not 161 — the mutant import drags in the whole `net/http` subtree
(V7, V12). One directive-unverified claim is now verified in the affirmative: `-overlay`
interacts correctly with `go list -deps` (§4 option 4 → V7/V8/V12). Additionally noted while
reading the checker (not a directive claim): the deps closure is never scanned against the
forbidden prefixes — direct imports only — recorded as `10/OD-1`.

## Related Documents

- [w-self-mod-vertical.md](w-self-mod-vertical.md) — `SM.B2a`, the milestone this item was
  promoted ahead of; see "Ordering vs `SM.B2a`".
- PR #45 → `1761a9c` (iter-55) — the per-group `extraForbidden` asymmetry this design must
  preserve (P7, AC5).
- [design_docs/coding-standards.md](../coding-standards.md) — S6 non-vacuous gates; the AC
  vacuity self-tests above apply it.
- `design_docs/verification/w-race-gate-blindspot/` — why go1.26.0–1.26.5 is denied and
  `GOTOOLCHAIN=go1.25.6` is the operative toolchain for every measurement here.

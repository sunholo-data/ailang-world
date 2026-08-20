# w-archive-stderr-in-manifest — stderr in the version manifest: the first PERSISTED instance of the stderr-merge class

**Status**: DESIGNED — quorum-clean, READY FOR sprint-planner (iteration 97, 2026-08-19).
Rounds 1 and 2 both returned BLOCKED; round 2 additionally recovered `gpt5-6-sol` from a
`budget` absence by a solo re-run at a raised cap (it returned reject). All four objections
across the two rounds were applied — rounds 1's inline (§Quorum verification log) and round 2's
VERBATIM under the shared skill's narrow-refinement carve-out (§Round-2 quorum …), which routes
the doc to sprint-planner without a third round. **No re-review is outstanding.**
**Charter clause**: clause-2. **Queue row 21** of `design_docs/world-mission.md` is the
specification of record — locate it by row number
(`grep -nE '^21\. ' design_docs/world-mission.md`), never by line number: the charter prepends
a STATUS stamp per iteration, so any line citation into it rots by construction (the original
`:3495` here already reads `3530` one iteration later).
**Authorship**: design-doc-creator role, iteration 97, against HEAD `b5ddf0e` (branch `dev`,
working tree clean). Controller measurements from the iteration-97 brief were **re-derived
first-party in this session** before use; where a re-derivation refined a controller number, the
refinement is stated in place (see §Audit, note on "exactly 2"). Round-1 exception: the P9b
rig-wide sweep and the page-bound test timings are attributed VERIFIED BY CONTROLLER in place —
the P17 grep among them was additionally re-run first-party.
**Verified against**: pinned `ailang` v0.30.0 at `/tmp/ailang-v0300/ailang`; pinned Go toolchain
`go1.25.6` (`scripts/verify_go.sh` denies the rig's ambient go1.26.4).
**Estimated**: 0.5–1 d (§Estimate; within the mission's 3–4 d guardrail with room to spare).

---

## The defect in one paragraph

`host/archive/archive.go:391` (`func probeVersion`) captures the pinned interpreter's
`--version` with `cmd.CombinedOutput()`. Every `ailang` invocation on this rig currently writes
one line to **stderr** — `2026/08/19 19:15:20 Observatory: 309MB (warn threshold: 200MB)` —
because `~/.ailang/state` (309 MB and growing between two runs today) exceeds a 200 MB warn
threshold. So the "version" string `probeVersion` returns has a wall-clock-stamped log line as
its **first line**. That string is persisted three ways: (1) into the sidecar
`manifest.json` next to the archived interpreter (`archive.go:289` → `Manifest.Version` →
`writeManifest`); (2) from the manifest into `GET /v1/health` as `interpreter_version`
(`daemon.go:450` → `:569`); and (3) — **not named by the queue row, found during this design** —
`daemon.go:452` reduces the same manifest string via `releaseFromVersion`, which returns the
**first non-empty line**, so the daemon's epoch-registry candidate becomes the Observatory log
line instead of `AILANG v0.30.0`, and `registry.Bootstrap(s, release)` writes that candidate
into **epoch 1 of the epoch registry as a content-addressed store object** (`registry.go:117`).
Because the polluted line carries a timestamp and a size that both drift (controller measured
308 MB at 19:12:38; this session measured 309 MB at 19:15:20), two archives of the **same
binary bytes** produce **different** version strings, different epoch-registry candidate
strings, and therefore different registry objects — in a system whose whole design premise is
that identity is content-addressed and replay is deterministic.

## Why this instance is different in kind

Iteration 89 closed three sites of the stderr-merge class in `scripts/verify_ail.sh` /
`host/verifygate`. All three were **transient**: a bad parse in a process that then exited. This
one **outlives the process** twice over:

- The sidecar `manifest.json` is read back on every subsequent daemon start. The idempotent
  early-return in `Archive()` (`archive.go:253-267`) returns **without re-probing and without
  rewriting the sidecar** when the bytes are already archived — so a manifest polluted once is
  served polluted **forever**, including by builds that carry the fix (P10). The fix alone does
  not clean this rig; §Decision 2 exists because of this.
- The epoch-registry epoch-1 revision is an immutable content-addressed store object whose
  head `Bootstrap` **refuses to move**: an existing head that disagrees with the recomputed
  bootstrap revision is a fatal startup error, "never silently overwritten"
  (`registry.go:126-135`). Today **zero** stores exist (P9), so nothing is bricked yet — but the
  first durable store bootstrapped from the polluted manifest would bake the log line into its
  epoch 1, and *then* the fix (which changes the candidate string) would make the fixed daemon
  **refuse to start** against it. There is a window in which the repair is free; this design
  takes it.

## Premises (hard constraints; each carries a Premise Verification Log row)

- P1 — the pinned binary's `--version` writes 168 bytes to stdout (first line
  `AILANG v0.30.0`) and 63 bytes to stderr (the Observatory line), exit 0. Measured with
  separate files, not `2>&1`.
- P2/P3 — `probeVersion` merges the streams and the merged bytes reach `Manifest.Version` on
  disk.
- P4/P5 — the manifest string is served by `/v1/health` verbatim and reduced (first non-empty
  line) into the epoch-registry candidate, which `Bootstrap` persists as a store object.
- P6/P7 — the complete non-test subprocess surface of `host/` + `cmd/` is **five** `os/exec`
  sites, each read individually (§Audit). No alternate spawn spellings exist (control in the
  same call, same scope).
- P8 — `ailang publish --dry-run` prints its parseable lines to **stdout** and the Observatory
  line to **stderr** (measured live against `packages/world-core`), so switching
  `host/pkgproj` to stdout-only parsing is safe.
- P9 — the on-disk population today: **exactly one artifact tree exists on this rig** —
  `/private/tmp/world-demo.db.artifacts` (real path; `/tmp` is a symlink to it — see the
  instrument-failure note after the PVL), rig-wide sweep in P9b. Its single interpreter
  manifest **is** polluted (quoted in the PVL, byte-identical in kind to the M3 executor's
  observation), and the QUICKSTART demo store `/tmp/world-demo.db` **does not exist**.
- P10 — the idempotent path never re-probes: `archive.go:253-267` returns `ref, nil` before
  Step 5.
- P11 — base gates at HEAD `b5ddf0e`: `go build ./...` rc=0 and
  `go test ./host/archive/... ./host/pkgproj/... ./host/daemon/... -count=1` all ok under the
  pinned toolchain; `./scripts/verify_ail.sh` green at base (P11b). The FULL `verify_go.sh`
  base status is **UNDETERMINED** pending a clean re-run — the run behind the draft's RED
  claim was contaminated by rig load and named the wrong test (P11b addendum).
- P12 — the fixtures that must keep passing: `host/archive/archive_test.go`'s `fakeInterpreter`
  writes **stdout only**, and `host/daemon/daemon_test.go:583` asserts
  `health.InterpreterVersion` equals the fake's version **verbatim** — both remain green under
  a stdout-only probe.
- P13 — `host/pkgproj`'s exec seam (`CrossCheck`) has **zero** test coverage today: the
  package's three test functions (`TestGoldenHashes`, `TestGoldenTarballBytes`,
  `TestCompareFailsLoudlyOnEveryDisagreement`) never spawn a process. The T3 test below is the
  first.
- P14 — in the pinned toolchain (go1.25.6), `(*Cmd).Output()` populates
  `*ExitError.Stderr` when `c.Stderr` is nil (read in `$GOROOT/src/os/exec/exec.go:1003`), so
  the error path keeps its diagnostics without the merge. *(After round 1 this is load-bearing
  for the `pkgproj` site only; the archive probe uses explicit buffers — Decision 1.)*
- P16 — `CrossCheck`'s complete caller population is one production caller — the inline Go
  helper in `scripts/verify_world_package.sh`, executed under `run_bounded 120` (a real bound:
  `p.wait(timeout=…)` + SIGKILL on the process group, exit 124) — plus one test-only caller.
  Positive enumeration and the read of `run_bounded` are in the PVL.
- P17 — `host/archive` today contains **zero bounded execution of any kind** (one grep hit,
  `:56`, prose in a comment; known-positive control fires in three sibling `host/` files), and
  nothing bounds the daemon's startup archival sequence (`daemon.go New()`, `:431-462`, read —
  the only `context.WithTimeout` in `daemon.go` is the CLI client's, `:678`). The repo's
  bound-constant idiom is `readDeadline` (`daemon.go:119-128`, field-from-constant `:284-290`)
  and `busyTimeoutMillis` (`writer_lock.go:173-179`).
- P18 — `store.Open` **creates** a missing database (`store.go:218`), and QUICKSTART's first
  command serves against a fresh `/tmp/world-demo.db` — so AC4's live form is runnable despite
  the store's absence (P9b).

### Design Freeze (the sprint must not renegotiate these)

- The archived executable's **bytes and content address are never touched**. Only sidecar
  manifests and the code that writes/reads them change.
- `host/broker`'s deliberate stderr merge (`handlers.go:110`, `cmd.Stderr = cmd.Stdout`) is
  **out of scope and must not be "fixed"** — §Audit site 3 explains why narrowing it would
  blind `classifyPublisherResult`.
- `tools/launchd/*` is frozen core; nothing in this design goes near it. No `.ail` files
  change; language gaps (the Observatory-on-stderr behaviour itself) route upstream, not here.

## The audit — positive enumeration of every non-test subprocess site in `host/` + `cmd/`

Per the iteration-96 guardrail, this audit is established by **positive enumeration, not by a
negative grep**. The queue row's own control ("0 `.Output()` calls in that package") is the
negative form this mission has ruled out and is **not relied on here**. Method: enumerate every
`exec.Command`/`exec.CommandContext` construction in non-test `host/` + `cmd/` (PVL P6 —
five hits, all in `host/`, zero in `cmd/`), cross-check that the five files holding them are
exactly the five non-test files importing `"os/exec"`, check alternate spawn spellings with a
known-positive control **in the same grep call over the same scope**, and read each site's
output handling individually (PVL P7).

| # | Site | Construction | Stream handling | Downstream of the bytes | Verdict |
|---|------|--------------|-----------------|-------------------------|---------|
| 1 | `host/archive/archive.go:384` (`probeVersion`) | `exec.Command(execPath, "--version")` | `CombinedOutput()` (`:391`) | **Persisted**: `Manifest.Version` on disk → `/v1/health` verbatim → `releaseFromVersion` first line → epoch-registry epoch-1 candidate in the store. Error path interpolates merged `out` into the error (`:393`). | **THE DEFECT — and a pre-existing UNBOUNDED wait.** Fix: bounded `exec.CommandContext` + separate buffers (Decision 1) + sidecar self-heal (Decision 2). |
| 2 | `host/pkgproj/pkgproj.go:213` (`CrossCheck`) | `exec.Command(ailangBin, "publish", "--dry-run")` | `CombinedOutput()` (`:219`) | **Parsed**: merged bytes go straight into `parseDryRun(out)`. The line-anchored regex cannot match the Observatory line itself, but a stderr line in dry-run shape would be parsed as data (T3 demonstrates this), duplicate-line detection sees both streams, and the `missing dry-run … in output: %s` diagnostic quotes the merged noise. Not persisted. | **SAME CLASS, transient-parse severity** (the iteration-89 kind). Fixed here because it is the same one-line shape and P8 verifies the fix is safe. |
| 3 | `host/broker/handlers.go:93` (`runBounded`) | `exec.CommandContext` | `StdoutPipe()` + `cmd.Stderr = cmd.Stdout` (`:110-113`) — a **deliberate** merge into one bounded, overflow-detected read | Three callers: git commit (`handlers_git.go:52`, opaque payload), model runs (`handlers_model.go:82`, opaque payload), and registry publish (`registry_publish.go:502`), where `classifyPublisherResult` scans the merged text for the publisher's own status markers — messages a CLI may emit on **either** stream. | **INTENTIONAL, KEEP.** Nothing here is persisted as structured identity; narrowing to stdout-only could make the classifier miss a stderr-emitted `publish blocked:` marker and misclassify a failure. Deliberate non-change (Conflict Surface). |
| 4 | `host/capsule/capsule.go:154` | `exec.CommandContext` | Separate `StdoutPipe()` / `StderrPipe()`, drained independently | stdout and stderr kept distinct in `Result` | **CORRECT.** No change. |
| 5 | `host/replay/replay.go:327` | `exec.CommandContext` | Separate `bytes.Buffer`s (`cmd.Stdout = &stdout; cmd.Stderr = &stderr`) | stdout is the replay output; stderr appears **only** in the error detail | **CORRECT — and it is the model Decision 1 adopts.** No change. |

**Refinement of the controller's measurement 4**: the controller's "exactly 2" counted
`CombinedOutput()` calls specifically; this enumeration is of **all five** non-test
`exec.Command*` constructions in `host/` + `cmd/`, of which sites 1 and 2 use
`CombinedOutput()` — the counts agree, this table is just the closed form the guardrail asks
for. Closure: the five sites live in exactly the five non-test files importing `"os/exec"`
(P6b), and a same-call grep for `os.StartProcess`/`syscall.Exec` alongside `exec.Command`
returned only the five known `exec.Command` lines — the alternations found nothing while the
control alternation found all five, in one invocation over one scope (P6c).

## Decision 1 — BOUNDED, stream-separated probes at both `CombinedOutput` sites

Two defects are corrected at the probe, not one. The stream merge is the queue row's defect.
The **unbounded wait is pre-existing**: today's `cmd.CombinedOutput()` at `archive.go:391`
waits forever on a hung interpreter and always has — `host/archive` contains zero bounded
execution of any kind (P17) while three sibling `host/` packages use `exec.CommandContext`
(the known-positive control). This design does not introduce that defect — but **Decision 2
would have widened its blast radius**: it adds a probe call to the idempotent path, which
today returns without executing anything and which is reached on every daemon startup —
placing the pre-existing unbounded wait on the one path whose hang blocks `New()` entirely,
since nothing bounds the startup archival sequence (P17). The round-1 quorum caught this
(§Quorum verification log, objection 1); the fix below closes the pre-existing unboundedness
*and* makes Decision 2 safe to ship.

`probeVersion` becomes a method on `Archive` carrying a documented bound, in the
`host/replay` model the audit already names as correct (§Audit site 5: `exec.CommandContext`
+ separate `bytes.Buffer`s):

```go
// probeTimeout bounds the ELAPSED TIME of one "<interpreter> --version"
// execution. The probe runs at daemon startup (Decision 2's idempotent-path
// heal), where nothing else bounds the wait (New() carries no deadline, P17),
// so without this constant a hung interpreter blocks startup forever — the
// bounded-waits axiom applied to the archive. 10 s is ~2 orders of magnitude
// above the measured probe (rc=0 in well under 1 s for 168 bytes of stdout,
// P1) yet far below the verify gate's enclosing test budgets (8 m race /
// 600 s outer, P17), so the labelled KindExecFailure always fires before any
// enclosing kill does.
const probeTimeout = 10 * time.Second
```

```go
func (a *Archive) probeVersion(execPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.probeTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, execPath, "--version")
	cmd.Env = childenv.Scrubbed(os.Environ()) // Decision 4, w-self-mod-vertical: preserved
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("%s --version: timed out after %v: %w (stdout: %q, stderr: %q)",
				execPath, a.probeTimeout, context.DeadlineExceeded, stdout.String(), stderr.String())
		}
		return "", fmt.Errorf("%s --version: %w (stdout: %q, stderr: %q)",
			execPath, err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}
```

- `a.probeTimeout` is a **field set from the constant by `New`** — the `daemon.go`
  `readDeadline` pattern verbatim (`daemon.go:284-290`, P17): the wiring stays assertable and
  the value shrinks in tests, which is the only way T2's deadline arm can exercise the SHIPPED
  timeout branch without a ten-second test.
- Success path: `Manifest.Version` is the verbatim **stdout** of `--version` — 168 bytes on
  the pinned binary, first line `AILANG v0.30.0`, deterministic for fixed binary bytes.
- Error path: both streams are captured separately and labelled in the error — nothing is lost
  relative to today's merged interpolation at `:393`. On deadline expiry the error wraps
  `context.DeadlineExceeded` explicitly, and because `ReplayError` unwraps (`archive.go:104`),
  the unchanged `KindExecFailure` wrapping at the caller (`archive.go:290-298`) leaves
  `errors.Is(err, context.DeadlineExceeded)` true end-to-end — the explicit timeout cause the
  objection requires, assertable through the public error chain.
- Ordering (the queue-row-22 lesson, stated so this doc does not reproduce it): `probeTimeout`
  has no sibling constant it silently depends on. The daemon's bounds block
  (`daemon.go:115-128`) bounds serving-path waits only; `New()`'s startup sequence carries no
  deadline (P17) — which is *why* the probe must bring its own bound rather than assume an
  enclosing one. The one enclosing bound on any measured path is the verify gate's test budget
  (`-timeout 8m` race under a 600 s outer wait, `verify_go.sh:150-153`), and 10 s ≪ 480 s
  holds with no cross-constant assertion machinery beyond T2's deadline arm, which reds if the
  probe regresses to unbounded.

**`CrossCheck` (`pkgproj.go:213/219`) — the objection's second arm, taken as
justify-and-verify.** Its complete caller population is positively enumerated (P16): **one**
production caller — the inline Go helper in `scripts/verify_world_package.sh`, executed at
`:183` as `run_bounded 120 …` — and **one** test-only caller
(`host/broker/registry_publish_test.go:1269`). `run_bounded` was read, not assumed
(`verify_world_package.sh:38-53`): Python `p.wait(timeout=t)`, on expiry
`os.killpg(…, SIGKILL)` with `start_new_session=True` — the kill reaches the **process
group**, so the `ailang` grandchild cannot outlive the 120 s bound — exit 124. That is a
verified caller-supplied bound. The bound-both alternative is **rejected deliberately**:
adding a second, in-process constant *under* the shell's 120 s would mint exactly the
queue-row-22 defect — two unconditional constants whose required ordering (in-process
< 120 s) crosses a Go/shell boundary where no test can pin it. Instead, `CrossCheck`'s doc
comment gains the explicit obligation: *callers must supply the bound; the sole production
caller does, via `run_bounded 120`.* Its exec still converts for the stream split:

- `host/pkgproj/pkgproj.go:219`: `cmd.Output()`, `parseDryRun` consumes stdout only (safe per
  P8), and the error return at `:221` interpolates `out` + `ee.Stderr` (P14) instead of the
  merged bytes.
- Doc comments that say "combined output" (`archive.go:380-382`) and "verbatim `ailang
  --version` output" (`daemon.go:380-381` on `HealthResponse.InterpreterVersion`) are updated
  to say stdout (and, for the probe, bounded) — the daemon comment is a **comment-only**
  change; `daemon_test.go:583` keeps passing because its fake writes stdout only (P12).

## Decision 2 — already-written manifests: repair by SELF-HEAL on the idempotent path, and why now

**The population** (P9; rig-wide sweep P9b, verified by controller): exactly one artifact
tree exists on this rig —
`/private/tmp/world-demo.db.artifacts/interpreters/sha256/e9746…3fb5/` — and its manifest's
`version` field begins `2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\n`.
The demo store `/tmp/world-demo.db` does not exist (the artifacts tree is orphaned), so **no
epoch registry currently holds a polluted candidate**.

**Why "the fix prevents new pollution" is not enough**: the idempotent path (P10) returns
before Step 5 — it neither re-probes nor rewrites. A fixed daemon started against the existing
artifacts tree reads the polluted sidecar, serves it on `/v1/health`, and — worse — bootstraps
any **fresh** store with the polluted, timestamp-bearing candidate. The pollution would
graduate from a disposable sidecar into an immutable content-addressed store object, and from
that moment the clean candidate becomes a `Bootstrap` divergence error (`registry.go:126-135`):
repairing later means refusing to start against every store created in between. Repair is
cheap now and expensive after the first durable store exists; **now** is the decision.

**Mechanism**: on the idempotent path (existing sidecar, hashes agree), run the (fixed)
`probeVersion` against the archived executable; if the fresh stdout differs from the stored
`Manifest.Version`, rewrite the sidecar with the fresh value (all other fields recomputed the
same way Step 5 computes them). Properties:

- **Convergent**: `--version` stdout is deterministic for fixed binary bytes (P1: 168
  identical bytes across runs), so at most one rewrite ever happens per artifact; after that
  the compare is equal and the path is a true no-op again.
- **Content addressing untouched**: the sidecar is not part of the content address (the hash
  covers the executable bytes; `manifest.json` sits **beside** `ailang` in the digest-named
  directory). The hash-mismatch refusal (`archive.go:255-263`) is evaluated before the heal and
  is unchanged.
- **Probe failure on this path fails loud** (`KindExecFailure`, same as the fresh path). Today
  the idempotent path executes nothing, so this is a new startup failure mode — deliberately:
  an archived interpreter that cannot answer `--version` cannot serve capsule or replay either,
  and a daemon that pins it should say so at startup, not at first use. *(Controller default —
  flagged for the reviewer; the alternative is heal-best-effort/keep-old-on-failure.)*
  **The probe this path runs is the BOUNDED one (Decision 1)**: the new startup execution
  costs at most `probeTimeout`, never an unbounded hang — without that bound this decision
  would have placed a pre-existing unbounded wait on every daemon start, which is round 1's
  objection 1. A hang past the deadline surfaces as `KindExecFailure` wrapping
  `context.DeadlineExceeded`, not as a startup that never returns.
- **The one real polluted manifest needs no manual step**: the next daemon start (or
  `Archive()` call) against that tree heals it. AC4 exercises exactly this, against the real
  artifact.

**Rejected alternative — delete `/tmp/world-demo.db.artifacts` by hand**: repairs this rig
once, leaves the class open on every other rig and every future recurrence (any future
stderr-emitting condition, not just Observatory), and leaves nothing testable behind.

## Decision 3 — the red tests, and why they cannot pass vacuously

**The hazard named by the queue row**: a test that shells out to the **real** `ailang` sees the
Observatory line only while `~/.ailang/state` exceeds 200 MB — on a clean rig it would pass
with `CombinedOutput()` still in place. Environment-dependent red is vacuous green elsewhere.
All three tests therefore use **fake emitters that write a known line to stderr
unconditionally** — the stderr emission is in the fixture, not in the environment.

- **T1 (archive, the load-bearing AC)**: a `fakeInterpreter` variant (extending the existing
  helper at `archive_test.go:22`) whose script runs
  `echo "FAKE-STDERR-MARKER should never reach the manifest" 1>&2` before printing its version
  to stdout. Archive it; `ReadManifest`; assert **both** (a)
  `m.Version == "<exact stdout string>"` and (b) `!strings.Contains(m.Version, "FAKE-STDERR-MARKER")`.
  Under the named red mutation (restore `CombinedOutput()` at the probe), the marker lands in
  `m.Version` and **both** assertions fail. It cannot pass vacuously: the fake always writes
  stderr regardless of any rig state, and assertion (a) pins the full value, so even a
  reordered or partial merge fails it.
- **T2 (self-heal)**: archive a fake once, overwrite its sidecar's `version` with
  `"2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\n<stdout>"` (the real
  pollution shape), call `Archive()` again on the same bytes, assert the sidecar now carries
  the clean stdout. Second arm: with a clean sidecar, re-`Archive()` and assert the version is
  unchanged (convergence, no rewrite churn). Third arm: replace the archived executable's
  sidecar-adjacent binary with a non-executable file and assert the idempotent-path probe
  failure is a `KindExecFailure` (the fail-loud decision is asserted, not assumed).
  **Fourth arm (deadline — round 1)**: shrink the `probeTimeout` field (e.g. 200 ms — the
  `readDeadline` shrink-in-tests pattern, P17), point the archive at a fake interpreter whose
  script sleeps far past the bound (e.g. 10 s, 50×) before answering, call `Archive()`, and
  assert (a) it returns within a deterministic upper bound on the measured wall clock (e.g.
  < 5 s — above the 200 ms bound with scheduler headroom, far under the 10 s sleep, so an
  unbounded probe cannot satisfy it), and (b) the error is a `KindExecFailure` with
  `errors.Is(err, context.DeadlineExceeded)` true — the explicit timeout cause asserted
  through the `ReplayError` chain (`Unwrap`, `archive.go:104`). Under the unbounded mutation
  the probe outlives the deadline, succeeds after the fake's finite sleep, and BOTH assertions
  fail — a deterministic ~10 s red, an assertion failure rather than a hang.
- **T3 (pkgproj, first coverage of the exec seam — P13)**: the test builds a real minimal
  package dir, computes the expected local hashes with the package's own exported functions
  (`ContentHash`/`CreateTarball`/`TarballHash`/`InterfaceHash`), then **generates** a fake
  `ailang` script that prints a well-formed dry-run block embedding those hashes on **stdout**
  — and prints one **regex-matching** line, `  Tarball: 999 bytes (sha256:<17 hex>...)`, on
  **stderr**. Under `.Output()`, `parseDryRun` sees one Tarball line and `CrossCheck`
  succeeds. Under the red mutation (restore `CombinedOutput()` at `pkgproj.go:219`),
  `parseDryRun` sees **two** Tarball lines and returns `duplicate dry-run Tarball line` — the
  test asserts success, so it reds. This is the sharper hazard demonstrated, not just the
  cosmetic one: merged stderr can **inject parseable data**, not merely noise.

**Mutation protocol (per the mission's non-vacuity rule)**: each red run must (1) apply the
named mutation, (2) show `go build ./...` rc=0 on the mutated tree — a mutant that does not
build proves nothing, (3) prove the mutation landed (`git diff --name-only` listing exactly the
one mutated file, plus the diff hunk showing `CombinedOutput` restored), (4) show the named
test failing, (5) revert and show the tree clean. Results go in the sprint's Verification Log.

## Milestone A — the whole item (~0.5–1 d, single milestone)

1. Decision 1 at both sites + comment updates (archive.go, daemon.go): bounded
   `CommandContext` probe with separate buffers, `probeTimeout` constant + field wiring in
   `New`, pkgproj `.Output()`. ~25 lines of code.
2. Decision 2 self-heal on the idempotent path. ~20 lines.
3. T1, T2 (four arms, incl. the deadline arm), T3 with their red-mutation runs recorded.
   ~180–230 lines of test.
4. AC4a fixture check (real manifest bytes) + AC4b live check against the real orphaned tree.

## Files to Create/Modify

- `host/archive/archive.go` — `probeVersion` → bounded `exec.CommandContext` + separate
  buffers + labelled error streams; `probeTimeout` constant + `Archive` field wired in `New`;
  idempotent path gains the self-heal; mechanics/doc comments updated (~45 lines changed).
- `host/archive/archive_test.go` — `fakeInterpreter` stderr variant; T1; T2 (four arms, incl.
  the deadline arm with a blocking fake).
- `host/pkgproj/pkgproj.go` — `CrossCheck` exec → `.Output()` + labelled error streams; doc
  comment records the caller-supplied-bound obligation (P16) (~10 lines changed).
- `host/pkgproj/pkgproj_test.go` — T3 (generated fake CLI + duplicate-injection red).
- `host/daemon/daemon.go` — comment-only: `HealthResponse.InterpreterVersion` doc says stdout.

No new files, no new packages, no `.ail` changes, no `tools/launchd` changes.

## Conflict Surface

- **Decision 6 archival semantics (`w-worldd-m2`, implemented)** — "identical bytes are an
  idempotent no-op success" gains one qualification: the no-op now verifies (and if stale,
  rewrites) the **sidecar**, never the artifact bytes. The hash-mismatch refusal and the
  atomic-rename write path are untouched. The archived executable remains 0o555 and
  byte-identical — T2's clean-sidecar arm asserts the no-op stays a no-op.
- **`GET /v1/health` consumers** — `interpreter_version` changes value on polluted rigs (from
  231 merged bytes to 168 stdout bytes). Enumerated consumers: `docs/QUICKSTART.md` does not
  print or assert `interpreter_version` (PVL P15); `host/daemon/daemon_test.go:583` asserts the
  fake's verbatim stdout and keeps passing (P12). No other reader of the field exists in-repo
  (the health struct is enumerated in the same PVL row's grep, whose positives — the test
  assertions — prove the pattern and scope).
- **Epoch registry (`w-log-epoch-decision`, RATIFIED)** — replay identity is pinned by
  `interpreter: HashRef`, **not** by the version string, so nothing in the ratified log format
  moves. What changes is the epoch-1 **candidate string** for stores bootstrapped after the
  fix: `AILANG v0.30.0` instead of a timestamped log line. `Bootstrap` divergence: fatal only
  for a store whose epoch 1 already holds the polluted candidate — **population today: zero
  stores exist at all** (P9). This is the window argument in Decision 2; the sprint should land
  before any durable store is created.
- **`host/broker`'s intentional merge** — deliberately NOT changed (§Audit site 3). Marker
  classification in `registry_publish.go` reads the merged stream because publisher status
  lines may arrive on either stream; a well-meaning "same fix here" would be a regression. This
  is recorded so a future auditor doesn't re-flag site 3 as a missed instance.
- **Iteration-89 sites** (`scripts/verify_ail.sh`, `host/verifygate`) — already closed,
  untouched; this item is the persisted sibling of that class, not a reopening.
- **Upstream (`sunholo-data/ailang`)** — the Observatory-on-stderr behaviour of the released
  binary is the *stimulus*, not the defect; the host must be correct against any stderr chatter.
  Whether the interpreter should log Observatory warnings on `--version` at all is an upstream
  question; if the sprint wants it raised, it routes as an issue + mission-control note per the
  frozen-core rule — **no local workaround exists in this design** (we do not suppress or parse
  around the line; we stop merging it into data).

## Systemic-Issue Audit

This is instance four of the stderr-merge class (three closed in iteration 89, transient; this
one persisted) — the class question "is this a one-off?" is answered by the §Audit table: the
complete non-test subprocess surface is five sites, two carried the defect shape, both are
fixed by this design, two already handle streams correctly, and one merges deliberately with
the reason recorded. The class is closed **by positive enumeration**, not by a zero-hit grep,
and the table is re-runnable in one command (PVL P6) for any future re-audit.

## Acceptance Criteria (each can fail; base status stated per command)

All commands run with `GOTOOLCHAIN=go1.25.6` and `AILANG_BIN=/tmp/ailang-v0300/ailang`
exported, from the repo root.

- **AC1 (the load-bearing criterion)** — T1 exists and passes:
  `go test ./host/archive/ -run 'Stderr' -count=1 -v` shows T1 running (not "no tests to run")
  and passing. **At base**: the test does not exist — `grep -c 'FAKE-STDERR-MARKER'
  host/archive/archive_test.go` returns 0 (so the post-change run measures the change, not the
  repo). Non-vacuity is carried by AC2, since `-run` with no match is green by construction.
- **AC2 (named RED mutation, archive)** — with `CombinedOutput()` restored at the probe:
  `go build ./...` rc=0 on the mutated tree, mutation proven landed (`git diff --name-only` =
  exactly `host/archive/archive.go`, hunk shown), then T1 **and** T2's heal arm FAIL; revert
  shown clean. **At base**: not runnable (T1/T2 don't exist); at sprint end this is the
  decisive red.
- **AC3 (named RED mutation, pkgproj)** — same protocol at `pkgproj.go:219`: mutant builds,
  mutation proven, T3 fails on `duplicate dry-run Tarball line`; revert clean. **At base**: not
  runnable (T3 doesn't exist; the seam has zero coverage, P13).
- **AC4 (the real artifact heals) — re-scoped in round 1** because the store is ABSENT
  (P9b: `ls /tmp/world-demo.db` → No such file), so "start the daemon against the existing
  orphaned tree" as originally worded presumed a db that does not exist. Two forms, ranking
  stated:
  - **AC4a (primary — deterministic unit form)**: T2 arm 1 seeded with the REAL polluted
    manifest bytes, pinned as a test fixture verbatim from this rig (the exact `version`
    string is quoted in PVL P9b); assert the heal rewrites it to the clean stdout. Primary
    because it is hermetic and reproduces on any rig — the acceptance tooth does not depend
    on this rig's `/private/tmp` surviving until sprint time.
  - **AC4b (live form — kept, and now more informative)**: runnable despite the absent db
    because `store.Open` **creates** a missing store (P18). `ailang-worldd serve --db
    /tmp/world-demo.db --ailang-bin /tmp/ailang-v0300/ailang` creates the store fresh;
    `Archive()` hits the EXISTING orphaned tree's idempotent path (same binary bytes → same
    digest `e9746…`), heals the sidecar; and `registry.Bootstrap` then writes the CLEAN
    candidate into the fresh store's epoch 1 — witnessing Decision 2's window argument
    end-to-end, which the original wording never exercised. Assert: `curl -s
    localhost:<port>/v1/health` shows `interpreter_version` beginning `AILANG v0.30.0`, and
    the on-disk manifest's `version` no longer contains `Observatory`. **At base**: FAILS —
    that manifest's `version` field begins
    `2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\n` (quoted in full, PVL
    P9b), and the idempotent path would serve it verbatim (P10).
- **AC5 (gates)** — `./scripts/verify_ail.sh` green, and the Go gates green **on the three
  touched packages** (`go build ./...` rc=0 plus
  `go test ./host/archive/... ./host/pkgproj/... ./host/daemon/... -count=1` ok). **At base**:
  `go build ./...` rc=0 and all three packages ok (PVL P11); `verify_ail.sh` GREEN (P11b).
  **The full `verify_go.sh` base status is UNDETERMINED** (P11b addendum: the draft's RED
  claim rested on a run contaminated by rig load, and it named the wrong test; a clean re-run
  is in flight) — so full-gate green is NOT an acceptance criterion this item can own either
  way: if the clean re-run reds, the sprint must neither launder that pre-existing red nor
  absorb fixing it (it stays routed to the controller as a base finding); if it greens, the
  sprint re-measures the full gate at sprint time rather than assuming this doc's baseline. A
  gate green at base measures the repo, so AC5's information content is "still green after
  the change" — the change-specific teeth are AC1–AC4 and AC7.
- **AC6 (audit holds post-change)** — re-run the enumeration
  `grep -rn --include='*.go' -E 'exec\.Command(Context)?\(' host/ cmd/ | grep -v '_test.go'`
  and confirm the same five sites, with sites 1 and 2 now reading via `.Output()` (shown by
  quoting the two call sites). **At base**: sites 1 and 2 read via `CombinedOutput()` (PVL P2,
  P8a). This is stated as a positive re-enumeration, not as "0 CombinedOutput remaining".
- **AC7 (bounded probe — round 1)** — T2's deadline arm exists and passes: a fake interpreter
  blocking past the shrunk `probeTimeout` yields `KindExecFailure` with
  `errors.Is(err, context.DeadlineExceeded)` true, and `Archive()` returns inside the arm's
  asserted wall bound. **At base**: not runnable, and the defect it guards is measured present —
  `host/archive` contains zero bounded execution (PVL P17: one grep hit, prose in a comment,
  with the known-positive control firing in three sibling packages). Named RED mutation: drop
  the bound (`exec.CommandContext(ctx, …)` → `exec.Command(…)`); the unbounded probe then
  outlives the shrunk deadline, succeeds after the fake's finite sleep, and both the
  wall-clock and the `DeadlineExceeded` assertions fail — a deterministic ~10 s assertion-red,
  not a hang.

## Non-Vacuity — the named RED mutations

| Gate | Named mutation | Expected red | Proof obligations |
|------|----------------|--------------|-------------------|
| T1 (+T2 heal arm) | Restore `out, err := cmd.CombinedOutput()` in `probeVersion` | T1: both assertions fail (marker present in `Version`; equality broken). T2 heal arm: heal writes the merged string, equality with clean stdout fails | mutant `go build ./...` rc=0; `git diff --name-only` = the one file; diff hunk quoted; revert to clean tree shown |
| T3 | Restore `CombinedOutput()` in `CrossCheck` | `parseDryRun` returns `duplicate dry-run Tarball line`; `CrossCheck` errors; test asserting success fails | same protocol |
| T2 fail-loud arm | (Not a mutation — a fixture arm) unexecutable archived binary | `KindExecFailure` at the idempotent path | n/a (asserted directly) |
| T2 deadline arm (AC7) | Drop the bound in `probeVersion`: `exec.CommandContext(ctx, …)` → `exec.Command(…)` | Probe outlives the shrunk deadline, succeeds after the fake's finite sleep; wall-clock and `DeadlineExceeded` assertions both fail (deterministic ~10 s assertion-red, not a hang) | same protocol (mutant builds; diff hunk quoted; revert shown clean) |

The queue row's acceptance criterion — "a test that REDS when stderr is merged back in" — is
T1 under the first mutation; everything else is the audit's remediation carrying the same
standard.

## Estimate

**0.5–1 d, inside the row's ~0.5 d at the low end and inside the mission's 3–4 d guardrail at
either end.** The code delta is small (bounded probe + one `.Output()` conversion ~25 lines,
self-heal ~20
lines, comments); the bulk is T3's generated-fixture plumbing (computing local hashes to embed
in the fake CLI script) and the three mutation-protocol runs with their evidence capture. If
T3's fixture fights back (e.g., tarball hashing is path-order-sensitive in a fresh temp dir),
the honest ceiling is 1 d — that is above the row's 0.5 d estimate and is said here in those
words rather than rounded down; it does not approach the guardrail. Round 1 adds the bounded
probe (constant + field + buffers, ~15 lines) and the deadline arm (one fixture, one arm) —
inside the same envelope; the estimate does not move.

## Premise Verification Log (all rows first-party this session — iteration 97, HEAD `b5ddf0e`, 2026-08-19 — unless marked)

| Row | Claim | Command | Observed |
|-----|-------|---------|----------|
| P1 | `--version` stdout/stderr split; pollution is rig-state, not world-state | `/tmp/ailang-v0300/ailang --version >/tmp/iter97_stdout.txt 2>/tmp/iter97_stderr.txt; echo rc=$?; wc -c < each; head -1; cat stderr` | rc=0; stdout **168** bytes, first line `AILANG v0.30.0`; stderr **63** bytes = `2026/08/19 19:15:20 Observatory: 309MB (warn threshold: 200MB)`. Controller measured 308 MB at 19:12:38 the same day — the number **moved between two runs three minutes apart**, which is the determinism argument in one row. Separate files per the brief's warning; no `2>&1`. |
| P2 | The probe merges streams | Read `host/archive/archive.go:383-396` | `out, err := cmd.CombinedOutput()` at `:391`; error path `:393` interpolates merged `out`; env scrubbed via `childenv.Scrubbed` at `:390` (Decision 4 of w-self-mod-vertical — preserved by this design). |
| P3 | Merged bytes are persisted | Read `archive.go:286-309`, `:366-378` | `:289` `version, err := probeVersion(finalPath)`; `:299-305` `Manifest{… Version: version …}`; `writeManifest` writes indented JSON sidecar. |
| P4 | Manifest → daemon → wire | Read `host/daemon/daemon.go:439-452`, `:378-383`, `:563-575` | `:445` `ReadManifest`; `:450` `d.interpreterVersion = m.Version`; `:451` `release = releaseFromVersion(m.Version)`; `:380-382` field doc + `json:"interpreter_version"`; `:569` served in `handleHealth`. |
| P5 | First line becomes the persisted epoch-registry candidate | Read `daemon.go:529-540`, `host/registry/registry.go:101-160` | `releaseFromVersion` returns the **first non-empty line**; `Bootstrap(s, release)` builds `Registry{Epochs: [{Epoch: 1, Candidates: []string{releaseString}}]}`, `PutObject` + `SetRegistryHead`; existing head naming different bytes → error `existing head %q diverges from bootstrap revision %q` (`registry.go:130-134`) — the divergence hazard in Decision 2 is this code path, read, not inferred. |
| P6 | Complete positive enumeration: five non-test `exec.Command*` sites in `host/` + `cmd/`, none elsewhere | `grep -rn --include='*.go' -E 'exec\.Command(Context)?\(' host/ cmd/ \| grep -v '_test.go'` | Exactly 5 lines: `host/broker/handlers.go:93`, `host/archive/archive.go:384`, `host/capsule/capsule.go:154`, `host/pkgproj/pkgproj.go:213`, `host/replay/replay.go:327`. (Zero in `cmd/` — a scoped statement of this positive listing, not a bare grep-zero.) |
| P6b | Closure by import surface | `grep -rn --include='*.go' '"os/exec"' host/ cmd/ tools/ \| grep -v '_test.go'` | Exactly the same 5 files. Any `os/exec` spawn requires this import; each importing file's call sites are the ones read in P7. |
| P6c | No alternate spawn spellings — with a known-positive control **in the same call, same scope** | `grep -rn --include='*.go' -E 'os\.StartProcess\|syscall\.Exec\b\|exec\.Command' host/ cmd/ \| grep -v '_test.go'` | Returns exactly the 5 `exec.Command` lines from P6 and nothing else: the control alternation fires on all five known sites in the same invocation over the same paths, so the zero for the other two alternations is a measured zero, not a vacuous one. |
| P7 | Per-site stream handling (audit verdicts) | Read `handlers.go:88-135` (+ callers `handlers_git.go:52`, `handlers_model.go:82`, `registry_publish.go:490-530`, `classifyPublisherResult` `:659-688`), `capsule.go:154-190`, `replay.go:324-347` | Broker: `cmd.Stderr = cmd.Stdout` deliberate; `classifyPublisherResult` scans merged text for markers (`publish blocked:` etc.). Capsule: separate pipes. Replay: separate buffers, stderr only in error detail. As tabled in §Audit. |
| P8 | pkgproj parses the merged bytes; dry-run data is on stdout | Read `pkgproj.go:203-235`; live: `cd packages/world-core && /tmp/ailang-v0300/ailang publish --dry-run >out 2>err` | `:219` `CombinedOutput` → `:223` `parseDryRun(out)`. Live run rc=0: **stdout** carries `Tarball: 7856 bytes (sha256:5823edcf…)`, `Content hash:`, `Interface hash:` lines; **stderr** carries only the Observatory line. So `.Output()` keeps every parsed line and drops only the noise. |
| P9 | On-disk population: 1 polluted manifest, 0 stores | `find /tmp/world-demo.db.artifacts -type f`; `cat …/manifest.json`; `ls -la /tmp/world-demo.db*`; `find . -name '*.artifacts'` (repo) | Two files in the tree (`ailang` + `manifest.json`, digest `e9746fef…b8f3fb5`); manifest `version` = `"2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\nAILANG v0.30.0\nCommit: e37b370\n…"` — polluted, matching the M3 executor's live observation. `/tmp/world-demo.db` **absent** (the glob resolves to the artifacts dir only); repo-local find: none. **Correction recorded**: this session's first probe of the store was `strings /tmp/world-demo.db \| grep -c Observatory` → `0`, which was **vacuous — the file does not exist**; the `ls` above is what established the truth. The iteration-96 spine, caught in the act again. |
| P9b | Rig-wide artifact-tree population: exactly one | `find <root> -maxdepth 6 -name '*.artifacts' -type d` over `/private/tmp`, `$HOME/.ailang`, `$HOME/dev`, `$HOME/Library/Application Support` — **real paths, never symlinked roots** (instrument-failure note below); plus `ls /tmp/world-demo.db` | **VERIFIED BY CONTROLLER (iteration 97)** — the round-1 quorum rejected this row's earlier UNVERIFIED form; the completed sweep replaces it. Exactly **one** artifact tree exists on this rig: `/private/tmp/world-demo.db.artifacts`, holding exactly one interpreter slot and exactly two files (`…/interpreters/sha256/e9746fef8570bc42b8cc52c0e88b7088468a5d2bd38bb8c42e27e5859b8f3fb5/{ailang,manifest.json}`). Its manifest IS polluted, verbatim: `"version": "2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\nAILANG v0.30.0\nCommit: e37b370\nFull:   e37b370d1d7a9c4e7136b319e38bec4d5f2bd9a0\nBuilt:  2026-07-19T09:27:00Z\n\nThe AI-First Programming Language\nCopyright (c) 2025-2026\n"` — pinned as the AC4a fixture. The three `$HOME`-scope zeros are measurements, not absences of evidence: the SAME instrument at the SAME depth returns the `/private/tmp` tree. The store `/tmp/world-demo.db` is ABSENT (`ls` → No such file) — AC4 re-scoped accordingly. |
| P10 | Idempotent path never re-probes | Read `archive.go:250-267` | On sidecar-exists + hash-equal: comment "Identical bytes already archived: idempotent no-op success" and `return ref, nil` — before Step 5's probe. |
| P11 | Base gates green (pinned toolchain) | `GOTOOLCHAIN=go1.25.6 go build ./...; go test ./host/archive/... ./host/pkgproj/... ./host/daemon/... -count=1` (with `AILANG_BIN` set) | `build rc=0`; `ok host/archive 2.984s`, `ok host/pkgproj 0.272s`, `ok host/daemon 3.897s`. |
| P11b | Durable gates at base | `./scripts/verify_go.sh`; `./scripts/verify_ail.sh` (both with pin + `AILANG_BIN`) | See the Verification Log addendum (§below): `verify_ail.sh` GREEN; full `verify_go.sh` **UNDETERMINED** after the round-1 correction (wrong test named in the draft; contaminated run; clean re-run in flight). |
| P12 | Existing fixtures stay green under stdout-only | Read `archive_test.go:19-70`, `daemon_test.go:580-586` | `fakeInterpreter` script writes version to stdout only (no stderr writes in the helper); `daemon_test.go:583-585` asserts `health.InterpreterVersion` equals the fake's verbatim version. Both compatible with `.Output()`. |
| P13 | pkgproj exec seam untested | `grep -n 'CrossCheck\|func Test' host/pkgproj/pkgproj_test.go` | Three test funcs (`TestGoldenHashes:14`, `TestGoldenTarballBytes:32`, `TestCompareFailsLoudlyOnEveryDisagreement:65`); zero references to `CrossCheck`; no fake-binary fixture exists in the package (T3 creates the first). |
| P14 | `.Output()` captures stderr into `ExitError.Stderr` | `grep -n 'Stderr was nil' $(GOTOOLCHAIN=go1.25.6 go env GOROOT)/src/os/exec/exec.go` | Hit at `exec.go:1003` in the **pinned** toolchain's source (`go version go1.25.6 darwin/arm64` confirmed in the same call) — the error-path design is read from the exact stdlib that will compile it. |
| P15 | QUICKSTART does not surface `interpreter_version`; daemon tests do | `grep -n 'interpreter_version\|InterpreterVersion' docs/QUICKSTART.md host/daemon/*_test.go` | Hits only in `daemon_test.go:583-585`; zero in QUICKSTART — one call over both paths, with the daemon-test positives proving the pattern fires; the QUICKSTART zero is scoped to a file explicitly named in the same invocation. |
| P16 (round 1) | `CrossCheck` caller population and its bound | `grep -rn "CrossCheck" --include='*.go' .` (positive enumeration) + read `scripts/verify_world_package.sh:38-53` (`run_bounded`) and `:183` (the call) | Callers: the inline Go helper embedded in `verify_world_package.sh` (heredoc `:176`, executed at `:183` as `run_bounded 120 "$tmp_proj" env -u AILANG_REGISTRY_API_KEY go run "$tmp_helper" "$AILANG_BIN"`), and test-only `host/broker/registry_publish_test.go:1269`. `run_bounded` read: Python `p.wait(timeout=t)`; on expiry `os.killpg(os.getpgid(p.pid), SIGKILL)` with `start_new_session=True` — the **process group** dies, so the `ailang` grandchild cannot outlive the 120 s bound — exit 124. Caller-supplied bound: VERIFIED, not assumed. |
| P17 (round 1) | Zero bounded execution in `host/archive`; no deadline encloses daemon startup; repo bound-constant idiom | `grep -n "context\|Timeout\|CommandContext" host/archive/archive.go` (control: the same tokens fire in `host/capsule/capsule.go`, `host/replay/replay.go`, `host/broker/handlers.go`); `grep -n "WithTimeout\|WithDeadline\|WithCancel" host/daemon/daemon.go host/archive/archive.go`; read `daemon.go:115-128`, `:284-290`, `:431-462`; read `writer_lock.go:170-179`; `grep -n "timeout" scripts/verify_go.sh` | `archive.go`: exactly **one** hit, `:56`, prose in a comment ("context fields") — zero bounded execution, a measurement (controller-measured, re-run first-party this session) because the known-positive control fires in the three sibling files. `daemon.go`: the only `WithTimeout` is the CLI client's (`:678`; constant documented `:109`); `New()`'s startup sequence (`:431-462`, read) carries no deadline. Idiom: `readDeadline` constant `:119-128` with field-from-constant `:284-290` ("assertable wiring, shrinkable in tests"); `busyTimeoutMillis` `writer_lock.go:173-179`. Gate budgets: race arm `-timeout 8m` under `p.wait(timeout=600)` (`verify_go.sh:150-153`). |
| P18 (round 1) | The daemon creates a missing store | `grep -n "creates" host/store/store.go`; `grep -n "world-demo.db" docs/QUICKSTART.md` | `store.go:218`: "Open opens (or creates) the SQLite database at path"; `QUICKSTART.md:11`'s first command serves against a fresh `/tmp/world-demo.db`. AC4b is runnable with the db absent. |

### Instrument failure, recorded (round 1) — `find` silently declines a symlinked root

The controller's first P9b sweep ran `find /tmp -name 'manifest.json' -path '*artifacts*'`
and got **0** — and its known-positive control, `find /tmp -maxdepth 2 -name '*.artifacts'`,
**also got 0** — while `ls /tmp/world-demo.db.artifacts` succeeded in the same breath. Cause:
on this rig `/tmp` is a symlink (`lrwxr-xr-x /tmp -> private/tmp`), and `find` in its default
`-P` mode does not follow a symlink handed to it as the traversal **root** — the entire search
tree was empty, so check and control both returned a confident zero. The control did its job:
a zeroed control is the documented "instrument broken" signal, and the fix was to hand `find`
the REAL path, `/private/tmp`. (A symlink as an *intermediate* path component is resolved by
kernel path lookup — which is why P9's direct-path `find /tmp/world-demo.db.artifacts -type f`
was valid all along; the symlink is only fatal at the root.) Generalised: **a path that is a
symlink is a scope, and `find` silently declines it** — rule 3a(i-d), "scope the control to
the same path as the check", arriving through the filesystem instead of through a bad
directory name. Same class as the iteration-96 spine and this doc's own P9 correction; every
`find`-based premise row in this doc therefore names a real path.

## Related Documents

- `design_docs/world-mission.md` queue row 21 (`:3495`) — specification of record.
- `design_docs/planned/w-log-epoch-decision.md` (RATIFIED D1) — epoch-registry format;
  identity-by-HashRef is why the version string is metadata, and why the candidate string still
  matters (nomination surface).
- `design_docs/implemented/w-worldd-m2.md` — Decision 6 (startup archival) whose idempotent
  path Decision 2 amends; Decision 3 (frozen route table) whose `/v1/health` field this cleans.
- `design_docs/planned/w-self-mod-vertical.md` Decision 4 — the `childenv.Scrubbed`
  environment at both fixed sites is preserved verbatim.
- Iteration-89 log entry (`design_docs/world-mission-log.md`) — the three transient sites of
  this class; this doc is the persisted fourth.
- `design_docs/implemented/w-race-gate-blindspot.md` — why every command here pins
  `GOTOOLCHAIN=go1.25.6`.

---

## Verification Log addendum — base gate runs (P11b)

- `./scripts/verify_go.sh` at `b5ddf0e`: **UNDETERMINED — round-1 correction of this row's own
  first draft** (controller re-measurement, rule 3b(v)). The draft reported the gate RED with
  `TestRecoverCommitStopsAtPageBound` hanging; both halves of that claim needed correction:
  1. **Wrong test named.** The gate's own panic trace names the adjacent sibling,
     `TestRecoverEffectStopsAtPageBound` (`recover_test.go:488`, via `recoverEffectPending`,
     `recover.go:232`). The test the draft named, `TestRecoverCommitStopsAtPageBound`
     (`recover_test.go:473`), **passes in isolation** — controller-measured twice:
     `go test ./host/broker/ -run '^TestRecoverCommitStopsAtPageBound$' -count=1 -timeout
     600s` → rc=0, ok, at 224.699 s and 217.668 s. It is slow (`maxRecoveryPages = 1 << 20`
     iterations), not hung.
  2. **Contaminated instrument.** The gate run behind the RED claim executed while 64 orphaned
     CPU spinners from an unrelated controller load experiment held the rig at load average
     ~110 for over an hour, spanning the run. A clean re-run is in flight; until it reports,
     the full-gate base status is **UNDETERMINED, not RED**.
  What survives either outcome: the three packages this design touches are green at base
  (P11, direct runs), `verify_ail.sh` is green at base, and AC5 is scoped so the sprint
  neither launders nor absorbs any pre-existing full-gate red. Also true regardless: the two
  page-bound tests (~220 s each per the controller's isolated runs, sequential within
  `host/broker`) are a real hazard for any gate arm carrying a per-package or race timeout —
  the race arm runs `-timeout 8m` under a 600 s outer wait (`verify_go.sh:150-153`), and
  ~440 s of that budget is these two tests alone, before race overhead.
  **Recorded trap, reproduced and caught** (kept from the draft — still true): this session's
  first reading of the gate was `./scripts/verify_go.sh 2>&1 | tail -15; echo rc=$?` — which
  printed `rc=0`, the exit code of `tail`, over a run whose last line was `FAIL`. The
  iteration-96 log already names this exact `| tail` artifact; the FAIL text, not the piped
  rc, is the result.
- `./scripts/verify_ail.sh` at `b5ddf0e`: **GREEN** (rc=0) — `world package gate PASSED: 9/9
  steps performed non-zero work`, `verify gate PASSED: 10 required identities verified, 39
  named tests pass`. Note the gate's own first output line was
  `2026/08/19 19:26:47 Observatory: 309MB (warn threshold: 200MB)` — the stderr chatter this
  design stops merging into data, live in the base run itself.

---

## Quorum verification log

**Round 1 — BLOCKED** (both reviewers present, `absent_reviewers: []`, metered $0.0866). Two
blocking objections, each carrying a reviewer-authored `proposed_fix`, neither disputing the
design direction. Both applied in this revision; the revision is bounded to them plus one
controller premise correction.

### Objection 1 — `gpt5-6-sol` (REJECT): unbounded subprocess wait on the startup path

> "Decision 2 introduces an unbounded subprocess wait on the formerly no-op archive/startup
> path: `probeVersion` uses `exec.Command(...).Output()` without a context or timeout. A hung
> archived interpreter can therefore block daemon startup indefinitely, directly violating the
> bounded-waits axiom."

`proposed_fix`:

> "Revise Decisions 1–2 so version probing is bounded: make `probeVersion` accept a context or
> create a documented fixed timeout, invoke `exec.CommandContext`, preserve separate
> stdout/stderr capture, and return `KindExecFailure` with an explicit timeout cause when the
> deadline expires. Add a T2 arm using a fake interpreter that blocks beyond the deadline and
> assert Archive/daemon startup returns within a deterministic upper bound. Apply the same
> bounded execution rule to `CrossCheck` or explicitly justify and verify its existing
> caller-supplied bound."

**Resolution — applied in full.** Decision 1 rewritten: documented `probeTimeout` constant
(10 s, value reasoned in the comment), field-from-constant wiring in the `readDeadline` idiom
(P17), `exec.CommandContext` with separate `bytes.Buffer`s (the `host/replay` model), explicit
`context.DeadlineExceeded` cause surviving the `KindExecFailure` wrap (`ReplayError.Unwrap`,
`archive.go:104`). T2 gains the deadline arm (blocking fake, wall-clock upper bound asserted);
AC7 and a non-vacuity mutation row pin it. `CrossCheck` is resolved via the reviewer's second
arm — justify and verify: its complete caller population is one production caller under
`run_bounded 120` (SIGKILL on the process group, exit 124 — read, not assumed; P16) plus one
test-only caller; the bound-both alternative was rejected because an in-process constant under
the shell's 120 s would recreate queue row 22's unlinked-constants defect across a Go/shell
boundary no test can pin. **One precision the doc does not inherit**: the unboundedness is
pre-existing — today's `CombinedOutput()` at `archive.go:391` is exactly as unbounded as the
draft's `.Output()`, and `host/archive` has zero bounded execution at base (P17, measured with
a firing control). What was genuinely new in the draft is placement: Decision 2 put a probe on
the idempotent startup path, which today executes nothing — widening the blast radius of the
pre-existing defect onto daemon startup. The reviewer is right about the consequence and the
remedy; Decision 1 now states the provenance in exactly those terms, and the fix closes the
pre-existing unbounded wait as well as the widened placement.

### Objection 2 — `gemini-3-1-pro` (REJECT): an explicitly unverified premise (P9b)

> "The document contains an explicitly unverified premise regarding the on-disk population of
> artifact trees. Design docs cannot defer premise verification to the sprint; all claims must
> be verified by the author prior to approval."

`proposed_fix`:

> "Complete the background `find` command for P9b, read the result, and replace the unverified
> row with the actual observed output and count of artifact trees."

**Resolution — applied in full.** P9b replaced with the completed sweep (attributed VERIFIED
BY CONTROLLER, iteration 97); the UNVERIFIED marking is gone. Result: exactly one artifact
tree rig-wide (`/private/tmp/world-demo.db.artifacts`, one interpreter slot, two files), its
manifest polluted — quoted verbatim in the row and pinned as the AC4a fixture — and the store
`/tmp/world-demo.db` ABSENT. That absence forced the AC4 re-scope the sweep exposed: the unit
form (T2 arm 1 seeded with the real manifest bytes) is promoted to primary as AC4a, and the
live form is kept as AC4b — still runnable because `store.Open` creates a missing db (P18),
and now also witnessing a clean epoch-1 bootstrap into the fresh store. The sweep's first
instrument FAILED — `find` over the `/tmp` symlink returned a confident zero from both the
check and its known-positive control — and that failure is recorded in full after the PVL
(§Instrument failure): same class as the iteration-96 spine; every `find` premise row now
names a real path.

### Controller premise correction applied in the same round (rule 3b(v))

The draft addendum's claim that `verify_go.sh` is RED at base named the wrong test
(`TestRecoverCommitStopsAtPageBound`, which passes in ~220 s; the trace names
`TestRecoverEffectStopsAtPageBound`) and rested on a run contaminated by 64 orphaned CPU
spinners (load ~110). P11, AC5, and the addendum now state what is actually established:
touched packages green, `verify_ail.sh` green, full gate UNDETERMINED pending a clean re-run;
AC5's launder/absorb scoping is unchanged and survives either outcome. The ~220 s page-bound
tests are retained as a named hazard for per-package/race timeouts.

---

## Round-2 quorum and the narrow-refinement carve-out (controller-authored, iteration 97)

**Round 2 verdict: BLOCKED, then recovered to BLOCKED-with-both-reviewers-rejecting.** The
synthesis printed `blocked` with `absent_reviewers: [{"model":"gpt5-6-sol","reason":"budget"}]`.
Per the shared skill's rule, an absent reviewer was re-run alone with a raised cap
(`ailang design-review --reviewer gpt5-6-sol --max-cost-usd 0.30`, `$0.09069`, 16,122 in /
336 out) — and it returned **reject**. This is the documented self-selecting trap firing exactly
as described: the reviewer dropped on `budget` because *the doc had grown* (444 → 676 lines) in
the revision **its own round-1 objection drove**, so the eye that closed was the one most
load-bearing. Recorded here because a `proceed`-shaped synthesis with a named hole is not a pass,
and this one was not read as one.

Both surviving objections are **carve-out eligible**: each carries a concrete reviewer-authored
`proposed_fix`, and neither disputes the design DIRECTION (`.Output()`-class stream separation +
bounded probe + a red test). They are applied here VERBATIM, in the reviewers' own terms, and the
doc routes to sprint-planner without a third quorum round.

### Objection R2-a — `gemini-3-1-pro` (present, reject, `$0.03906`)

> "The design introduces a factual contradiction regarding Decision 2's self-heal mechanism. It
> claims that after a heal, 'the compare is equal and the path is a true no-op again.' However,
> the mechanism specifies that it unconditionally 'run[s] the (fixed) `probeVersion` against the
> archived executable' on the idempotent path to obtain the fresh stdout for that comparison.
> Unconditionally spawning an OS process on daemon startup permanently penalizes the hot path and
> breaks the existing early-return no-op behavior, which violates the 'true no-op' claim."

> `proposed_fix`: "Revise Decision 2 to perform a conditional self-heal without unconditional
> process spawning. On the idempotent path, read the existing `Manifest.Version`; ONLY run
> `probeVersion` and execute the rewrite if the string is demonstrably polluted (e.g.,
> `!strings.HasPrefix(m.Version, "AILANG v")` or `strings.Contains(m.Version, "Observatory")`).
> This preserves a true no-op (zero process executions) for healthy or already-healed artifacts."

**UPHELD and APPLIED. Decision 2 is amended to a CONDITIONAL self-heal.** On the idempotent path
the manifest is read first; `probeVersion` runs **only** when the stored string fails the
well-formedness predicate. The predicate is `!strings.HasPrefix(m.Version, "AILANG v")` —
positive-shape, not a denylist for `"Observatory"`, because the polluting line is one instance of
a class (any stderr chatter) and matching the *known* pollutant would leave the next one
undetected. A healthy or already-healed artifact therefore performs **zero** process executions,
which restores the literal early-return no-op the doc claimed and did not have.

Note what this also does to R1-a: with a conditional heal, the bounded probe never executes on the
hot path for healthy artifacts at all, so the two reviewers' fixes **converge** — `probeTimeout`
bounds the wait that remains, and the condition removes the wait entirely from the common case.
The doc must state this convergence rather than presenting the bound as the whole remedy.

### Objection R2-b — `gpt5-6-sol` (recovered from `budget`-absence, reject, `$0.09069`)

> "The rollout-safety premise that 'zero stores exist at all' is unverified. P9/P9b establish only
> that `/tmp/world-demo.db` is absent and inventory artifact trees in selected roots; they do not
> enumerate SQLite store files or inspect registry heads. This is load-bearing: if any existing
> store was bootstrapped from a polluted manifest, self-healing the sidecar before
> `registry.Bootstrap` changes the candidate and can make daemon startup fail with the documented
> divergence error."

> `catch`: "Do not infer a rig-wide store population from one missing QUICKSTART database or from
> an artifact-tree sweep. The related claim that exactly one artifact tree exists 'on this rig'
> also exceeds the listed, depth-limited search scope; narrow it to the searched roots unless a
> genuinely exhaustive inventory is performed."

> `proposed_fix` (option 1, the arm taken): "Replace P9/P9b's global claims with scoped claims,
> then … add a verification-log row that positively enumerates every configured/discoverable
> worldd database path and inspects each database's registry head and epoch-1 candidate, recording
> raw paths and results…"

**UPHELD and APPLIED, and the objection lands on a CONTROLLER-supplied number.** The phrase
"exactly one artifact tree exists **on this rig**" was written by the controller and handed to the
designer under a `VERIFIED BY CONTROLLER` heading. It was measured by
`find <root> -maxdepth 6 -name '*.artifacts' -type d` over four roots — so the scope is *four
roots at depth ≤ 6*, and "on this rig" is a strictly wider sentence than the command supports.
This is rule 3b(ix) — a count is only true inside the scope it was taken in — committed by the
same controller directive that instructed the designer to "scope every count into the sentence it
appears in". Third consecutive iteration in which a reviewer or the designer refutes a
controller-supplied number.

**P9b is REPLACED by the following, correctly scoped.** Searched roots: `/private/tmp`,
`$HOME/.ailang`, `$HOME/dev/sunholo-data`, `$HOME/Library/Application Support`, each at
`-maxdepth 6`. Within that scope, and stated as a claim about that scope:

| Measurement | Command | Result |
|---|---|---|
| Artifact trees | `find <root> -maxdepth 6 -name '*.artifacts' -type d` | **1**: `/private/tmp/world-demo.db.artifacts` |
| Its companion store | `test -f /private/tmp/world-demo.db` | **ABSENT** — the tree is ORPHANED |
| SQLite files | `find <root> -maxdepth 6 -name '*.db' -type f` | **53** |
| Of those, worldd stores | `test -d "${db}.artifacts"` per file, positively, one by one | **0** |
| Control | `test -d /private/tmp/world-demo.db.artifacts` | fires — the `-d` test can see a tree that exists |

The discriminator is the repo's own definition at `host/archive/archive.go:46-48`: a worldd store
is `<store>.db` with an adjacent `<store>.db.artifacts` tree. The 53 databases are coordinator,
brain, observatory and eval stores belonging to other tools; **none** has an artifact tree, so
none is a worldd store. The enumeration is POSITIVE — every `.db` was tested individually rather
than a pattern being asserted to match nothing.

**Consequence, which is stronger than the claim it replaces:** within the searched roots there is
**no worldd store at all**, hence none bootstrapped from a polluted manifest, hence the
`registry.Bootstrap` divergence hazard the reviewer names **cannot be triggered by any store in
scope**. The rollout-safety conclusion survives; what changes is that it is now a measurement with
its boundary stated, and it no longer claims anything about paths nobody searched.

**Residual, DECLARED rather than absorbed** (this mission's round-8 lesson: state residuals, and
accept that a disclosure is a surface an objection can attach to — a document is not penalised for
the defects it hides, and hiding this one would be strictly worse). The sweep is scoped, not
exhaustive: a worldd store outside those four roots, or deeper than 6 levels, would not have been
seen. The reviewer's option (2) — making rollout safe *without* any zero-store assumption, by
checking the target store's registry state before healing and refusing to silently overwrite a
polluted epoch — is **NOT applied here** and is recorded as owed. It is a genuine design
obligation for the tranche that wires this into a store with existing epochs, and folding it in
now would widen a ~0.5–1 d item into registry-migration work. **The sprint MUST NOT claim
rollout safety for stores outside the searched roots**; AC4's wording is bounded to the one
orphaned tree measured above.

### Base-gate status — RESOLVED, and the draft's base-red was controller-induced

The doc's P11b recorded the full `verify_go.sh` as RED at HEAD, and then as UNDETERMINED. It is
neither. Measured at `b5ddf0e` after removing the contamination: **`./scripts/verify_go.sh`
PASSES end to end** — `✓ go gate PASSED: build clean, plain and race tests pass with pinned
AILANG_BIN`, **zero** `FAIL` lines in the whole log, `host/broker` `ok` at **84.597 s** on the
plain leg and **183.918 s** on the race leg.

The earlier red was an artifact of **64 CPU spinners the controller started for an unrelated
load experiment and never killed** — `kill $SPINNERS` did not reach them, and the rig sat at load
average ~110 for over an hour, spanning both the controller's first gate run and the designer's.
Two independent roles therefore reported the same "base red", and it was a fact about the
observer. Two things travel with this: the failing test named in the draft
(`TestRecoverCommitStopsAtPageBound`) was the wrong sibling — the panic trace names
`TestRecoverEffectStopsAtPageBound` (`recover_test.go:488`) — and the named test PASSES in
isolation at 224.699 s / 217.668 s. The **real** durable finding is that the two page-bound tests
each run ~1,048,576 iterations and cost ~220 s apiece under load, which is a genuine hazard for
any per-package timeout; that is kept. AC5's launder/absorb scoping is unchanged and correct
either way.

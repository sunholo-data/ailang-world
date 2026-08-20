# w-capsule-output-cap-load-flake

**Queue:** World mission row 20  
**Status:** planned  
**Scope:** `host/capsule` Go host code and tests only; no `.ail` change  
**Estimate:** 1.0 day (revision-round-1 liveness coverage makes the prior 0.75 day estimate too small; the filed 0.25 day remains insufficient)

## 1. Problem statement

`TestF6OutputCapKillsChildBeyondOnePipeBuffer` currently asks one real interpreter process to race two limits in one `Config`: emit more than 1 KiB before a 5 s execution deadline. One of two historical full-suite runs returned `*TimeoutError` after 19.33 s rather than `*OutputLimitError`. That observation is real but weak (`n=2`), and the failure does not reproduce in the larger measured sample. Therefore this design does **not** use disappearance of the flake as evidence.

The filed remedy—raise the 5 s clock “well past plausible load”—is refuted as filed. The child execution measured by the controller is about 15 ms idle, already roughly 330 times below the clock. Meanwhile `Runner.Run` resolves and hashes a 91,826,738-byte interpreter before it creates the timeout context. The test's outer stopwatch includes that unbounded preparation as well. Raising the clock makes an already-large child margin larger, leaves the dominant pre-context term unbounded, and remains probabilistic.

The required property is narrower and stronger:

> Once output collection is given an over-limit stream, overflow causes the child-kill request and `OutputLimitError` wins over a simultaneous deadline, independent of child throughput, OS pipe capacity, CPU load, and wall-clock timing.

This is effects-at-the-boundary host behavior (coding standard S2), not package policy and not a pure World invariant. It therefore belongs in `host/capsule`, not in `world/` or an AILANG package. No Z3 contract or `.ail` artifact is appropriate. The slim-kernel/package-first rule is preserved because this change adds no kernel or user-facing surface.

The separately owned process-group-wide overflow-kill issue is explicitly out of scope. This item neither changes `cmd.Process.Kill()` to group kill nor edits `host/broker`.

## 2. Premise verification log

Every timing below is descriptive, not an acceptance threshold. Commands run from the repository root with the load-bearing environment shown.

| Claim | Establishing command | Observed result | Scope |
|---|---|---|---|
| The required Go toolchain is active. | `export PATH=/opt/homebrew/bin:$PATH; export GOTOOLCHAIN=go1.25.6; go version` | `go version go1.25.6 darwin/arm64` | This worktree and this shell invocation; not CI or another host. |
| Capsule tests have a silent-skip path unless the binary is set. | `sed -n '1,32p' host/capsule/capsule_test.go` | `pinnedBinary` calls `t.Skip` when `AILANG_BIN` is empty and `t.Skipf` when unusable. | Source at HEAD `47e12cc`; establishes code path, not runtime skip count. |
| The pinned interpreter is 91,826,738 bytes. | `export AILANG_BIN=/tmp/ailang-v0300/ailang; stat -f '%z %N' "$AILANG_BIN"` | `91826738 /tmp/ailang-v0300/ailang` | That file on this darwin/arm64 rig at measurement time. Re-derived first-party; agrees with the controller. |
| Resolve and full-file verification precede timeout creation. | `sed -n '126,158p' host/capsule/capsule.go; sed -n '222,246p' host/capsule/capsule.go; sed -n '388,414p' host/archive/archive.go` | `Resolve` then `verifyExecutable`/`os.ReadFile`/hash occur before `context.WithTimeout`. | Static control flow at HEAD only. Re-derived first-party; agrees with the controller. |
| Whole-file hashing is material and outside `ExecTimeout`. | `export AILANG_BIN=/tmp/ailang-v0300/ailang; for i in 1 2 3 4 5; do /usr/bin/time -p sh -c 'shasum -a 256 "$1" >/dev/null' sh "$AILANG_BIN"; done` | real time `0.20, 0.18, 0.18, 0.18, 0.18` s. | Standalone SHA-256 utility over this binary on this rig; it is not a direct timing of Go `verifyExecutable`. Re-derived independently; the absolute values differ from the controller's Go probe, but support only the shared structural conclusion that a large unbounded read/hash exists. |
| The named failure does not reproduce in a small fresh base run. | `export PATH=/opt/homebrew/bin:$PATH GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang; go test ./host/capsule -run '^TestF6OutputCapKillsChildBeyondOnePipeBuffer$' -count=5 -v` | 5 `RUN`, 5 `PASS`, 0 skips; package 4.472 s. | Five executions at HEAD on this unloaded rig. This neither estimates a flake rate nor refutes the historical failure. It independently confirms that non-reproduction cannot be an AC. |
| The controller's larger non-reproduction sample is correctly characterized. | Controller command record supplied with this item: 7 full-suite + 10 at 2x oversubscription + 15 at 6x, all with the pinned toolchain where required and the binary set. | 0 failures in 32 executions at `b5ddf0e`. | Controller's iteration-97 rigs and loads only. Relied on as supplied; the local 5-run arm above re-derives the direction, not the 32-run count. |
| The controller's phase split identifies verification as the dominant, unbounded term. | Controller's package-local probe at HEAD `47e12cc`, five runs per idle/loaded arm, `GOTOOLCHAIN=go1.25.6`, `AILANG_BIN` set, zero skips asserted. | Resolve 3–12 us; verify 37.4–46.4 ms idle and 44.2–314.8 ms loaded; full Run 57.6–59.4 ms idle and 123.8–268.5 ms loaded; all OVERFLOW. | One fixture, darwin/arm64 16-core rig, stated loads. Relied on as controller-supplied. Static ordering and independent file/hash measurements were re-derived above; the exact phase timings were not reproduced by the standalone hash command and are not generalized beyond that scope. |
| The test stopwatch includes preparation that `ExecTimeout` does not govern. | `sed -n '276,306p' host/capsule/capsule_test.go` together with the Run ordering command above. | `start := time.Now()` precedes `New(...).Run(...)`; context creation occurs inside `Run` after resolve, verify, and source staging. | Static control flow at HEAD. Re-derived first-party; agrees with the controller. |
| Base and filed-base capsule sources are identical. | `git rev-parse e4ba56d:host/capsule/capsule.go HEAD:host/capsule/capsule.go; git rev-parse e4ba56d:host/capsule/capsule_test.go HEAD:host/capsule/capsule_test.go` | Matching pairs: `39f453…` and `93fa23…`. | Those two files only, between filed base and HEAD. |
| The mutation targets `killOnce.Do`, `readCapped`, and `errOutputLimit` exist in the base source. The same call includes negative and positive controls. | `grep -cE 'zzNoSuchSymbolIter99' host/capsule/capsule.go; grep -cE 'func \(r \*Runner\) Run' host/capsule/capsule.go; grep -cE 'killOnce\.Do' host/capsule/capsule.go; grep -cE 'func readCapped' host/capsule/capsule.go; grep -cE 'errOutputLimit' host/capsule/capsule.go` | Counts `0, 1, 1, 1, 6`; target locations include `killOnce.Do` line 193, `readCapped` line 237, and the `errOutputLimit` declaration line 78. | Controller-verified at HEAD `47e12cc`; `host/capsule/capsule.go` only. The negative control establishes absence detection and the same-path positive control establishes pattern visibility. |
| M6's named existing integration test exists in the base source, and `readCapped`'s `io.LimitReader` cite is exact. | `grep -n 'func TestF6OutputCapReturnsStructuredOverflow' host/capsule/capsule_test.go; grep -c '^func Test' host/capsule/capsule_test.go; grep -c 'func TestZZQuxNotARealTestName' host/capsule/capsule_test.go; grep -n 'io.LimitReader' host/capsule/capsule.go` | `capsule_test.go:233 func TestF6OutputCapReturnsStructuredOverflow`; same-path positive control **7** `^func Test`; same-path negative control on a fresh literal **0**; `capsule.go:238 data, err := io.ReadAll(io.LimitReader(pipe, limit+1))`. | `host/capsule/{capsule,capsule_test}.go` only, at `5a38ac0`. **CORRECTION CARRIED INTO §3: the `io.LimitReader` call is at `:238`, not `:237`.** `:237` is the `func readCapped(` declaration line — the adjacent premise row cites `:237` correctly for the FUNCTION and §3 transcribed that number onto the CALL. Off-by-one, found only by running the command the objection asked for; the substantive claim (`readCapped` already uses `io.LimitReader`) is TRUE. The negative control fires, so the absence detection is live and the positives are not vacuous. |
| No capsule pipe is explicitly closed on overflow, and no non-test host command sets `WaitDelay`. | `grep -n 'Close()' host/capsule/capsule.go; rg -n 'WaitDelay|SysProcAttr' host --glob '!**/*_test.go'; grep -nE 'stdoutPipe|stderrPipe' host/capsule/capsule.go; grep -cE 'killOnce' host/capsule/capsule.go` | The first command has no pipe-close hit. `WaitDelay` has one repo-wide comment at `host/archive/archive.go:74` saying the general case is out of scope; `SysProcAttr` has two non-test hits, including capsule line 160. The added third command returns **exactly four** hits — `stdoutPipe` at `:168` and `stderrPipe` at `:172` (creation, from `cmd.StdoutPipe()`/`cmd.StderrPipe()`), `:197` and `:198` (drain use) — and nowhere else; the same-path control `killOnce` returns **2**, so the four is a measurement rather than a broken pattern. | Controller-verified at HEAD `47e12cc` and re-measured at `5a38ac0` for the round-3 revision; the `SysProcAttr` hits and the `killOnce` count are positive visibility controls. The stdout/stderr sub-claim is now ESTABLISHED BY THE ROW'S OWN COMMAND rather than asserted in prose, which is what the round-2 premise objection asked for. Together these establish that current production liveness comes from `ExecTimeout` cancellation and group SIGKILL, not drain-local cleanup. |

## 3. Chosen design: deterministic output-collection core

Extract the post-`cmd.Start` output lifecycle from `Runner.Run` into an unexported helper. The helper receives:

- stdout and stderr `io.Reader`s;
- the byte limit;
- the already-created execution context;
- a narrow child boundary with `Kill() error` and `Wait() error` (or equivalently injected kill/wait functions).

It owns the existing two drains, the once-only kill on `errOutputLimit`, waiting, and the error-precedence decision. `Runner.Run` remains responsible for archive resolution, verification, staging, context creation, command construction, pipes, and `Start`; after `Start` it adapts the real process to the helper. Production semantics and exported APIs do not change.

The seam has an explicit liveness precondition: **the caller owns the finite cleanup bound and must arrange that cancellation makes both supplied readers and `Wait` return**. In production, `Runner.Run` supplies that property through its existing `context.WithTimeout`, `exec.CommandContext`, group-wide `cmd.Cancel` SIGKILL, and `os/exec` pipe/process cleanup. The helper itself neither makes an arbitrary `io.Reader` cancellable nor makes `Wait() error` bounded. It may wait for both drains and then call `Wait` exactly as the current code does. Tests must therefore exercise the precondition rather than silently supplying already-terminated dependencies.

Replace the timing-bearing assertion in `TestF6OutputCapKillsChildBeyondOnePipeBuffer` with deterministic package-local tests of this helper:

1. **Overflow + expired context:** stdout is a `bytes.Reader` containing exactly `limit+1` bytes; stderr is empty. A fake child's `Kill` increments a counter and marks the fake killed. Its `Wait` reads that independently written killed state and returns an explicit “waited without kill” sentinel if false. Pass an already-expired context. Assert captured stdout is exactly `limit`, `Kill` was called exactly once, `Wait` observed the kill, and the result is `*OutputLimitError`, never `*TimeoutError`. No goroutine writes bytes, no pipe fills, and no duration is asserted.
2. **Within-limit control:** provide exactly `limit` bytes and a live context. Assert byte identity, zero kills, one wait, and the fake wait result. This makes an unconditional-kill implementation fail.
3. **Two-stream once control:** make both prefilled readers exceed the limit. Assert the child boundary records exactly one kill. This guards the `sync.Once` behavior without scheduler timing.
4. **Caller-supplied liveness control:** use readers whose `Read` calls signal entry and then block on a test-owned release channel, plus a fake `Wait` that signals entry and blocks on a separate test-owned release channel. Run the helper in a goroutine; after both reader-entry signals, the test closes the reader-release channel; after the wait-entry signal, it closes the wait-release channel; then it asserts the helper returns with the expected non-overflow result. A short test-only watchdog fails with a phase-specific diagnostic instead of hanging the suite, but elapsed duration is not a product assertion. This arm can fail when the release contract or drain-before-wait ordering is absent; unlike the finite `bytes.Reader` arms, it does not smuggle termination in as an intrinsic property of the fake.

Keep `TestF6OutputCapReturnsStructuredOverflow` as the real interpreter integration coverage with its small output. It proves production wiring reaches an `OutputLimitError`; the new tests prove the beyond-buffer kill/precedence algorithm. Delete the old outer stopwatch and the 64 KiB throughput fixture from the named test. The test no longer claims elapsed wall time measures `ExecTimeout`.

The fake's observables are deliberately written at the child boundary, not alongside the assertion: the helper calls `Kill`; the fake `Kill` method writes `killCount` and killed state; the fake `Wait` method independently reads killed state; assertions read the count and wait result. The output observable reads bytes produced before helper invocation. Thus the test cannot manufacture “kill succeeded” in the same branch that asserts it.

### Why the alternatives lose

**Bound or hoist verification.** Moving `context.WithTimeout` above verification would redefine `ExecTimeout` from execution time to total preparation-plus-execution time and make a slow disk/hash legitimately beat output. Hoisting verification into `New` requires an API error path, caching policy, and a TOCTOU analysis: verification must still bind the bytes eventually executed. Either is a larger security/semantics item and neither makes a real child's time-to-first-overflow deterministic. The unbounded phase deserves a separate budget/API item if total request latency must be bounded.

**Measure only the governed region.** Moving `start` to immediately before `cmd.Start`, or returning phase timestamps, removes the current contamination. It still asserts scheduler- and interpreter-dependent elapsed time against a wall clock, and instrumentation added solely for the test expands production surface. It diagnoses a race rather than making the property deterministic.

**Raise the deadline or enlarge output.** Both merely change odds. A larger deadline attacks the term with the greatest measured margin; more output still depends on interpreter production, pipe behavior, and scheduling. Neither can support a non-vacuous base-red gate.

The selected seam tests the causal state machine with pre-existing bytes and an independently observed kill. Load can delay the test process itself, but cannot change which input is available or which branch is correct; there is no competing real-time deadline.

### Conflict Surface

`readCapped` already uses `io.LimitReader` at `host/capsule/capsule.go:238` (the function is declared at `:237`); the design does not replace a missing standard-library cap. `io.LimitReader` supplies only a bounded byte view. It does not report that byte `limit+1` existed as a typed capsule overflow, kill the child exactly once across two concurrent drains, coordinate both drains with `Wait`, or enforce the domain precedence “overflow outranks deadline.” `context.WithCancel`/`WithTimeout` supplies a cancellation signal, and `exec.CommandContext` can act on it, but neither makes an arbitrary `io.Reader` cancellable nor combines stdout, stderr, kill, wait, and capsule error translation into this state machine. The unexported helper is required to compose those standard-library primitives and existing process effects; it does not attempt to supersede them or promise a new cleanup bound.

## 4. Acceptance criteria

All commands run at repository root with:

```sh
export PATH=/opt/homebrew/bin:$PATH
export GOTOOLCHAIN=go1.25.6
export AILANG_BIN=/tmp/ailang-v0300/ailang
```

The test commands must assert discovery in the same call; bare `go test -run` is forbidden because no-match exits zero. Each AC is red on the unchanged tree, rather than merely accompanied by a green repository gate.

### AC1 — overflow causally kills and outranks an already-expired context

Command:

```sh
out="$(mktemp)"; go test ./host/capsule -run '^TestOutputCollectionOverflowKillsAndOutranksDeadline$' -count=1 -v >"$out" 2>&1
rc=$?; grep -q '^=== RUN   TestOutputCollectionOverflowKillsAndOutranksDeadline$' "$out" && grep -q '^--- PASS: TestOutputCollectionOverflowKillsAndOutranksDeadline ' "$out" && [ "$rc" -eq 0 ]
```

Expected after change: success, exactly one named test discovered and passed; its assertions require `limit+1 -> limit`, one kill observed by `Wait`, `OutputLimitError`, and not `TimeoutError`. **At base:** fails because the named test does not exist, so the `RUN` grep is false even though Go's no-match exit is zero.

### AC2 — the non-overflow control forbids unconditional kill

Command:

```sh
out="$(mktemp)"; go test ./host/capsule -run '^TestOutputCollectionAtLimitDoesNotKill$' -count=1 -v >"$out" 2>&1
rc=$?; grep -q '^=== RUN   TestOutputCollectionAtLimitDoesNotKill$' "$out" && grep -q '^--- PASS: TestOutputCollectionAtLimitDoesNotKill ' "$out" && [ "$rc" -eq 0 ]
```

Expected after change: success; exact-limit bytes round-trip, `killCount == 0`, and `Wait` runs once. **At base:** fails on the missing `RUN` identity.

### AC3 — simultaneous stream overflow issues one kill

Command:

```sh
out="$(mktemp)"; go test ./host/capsule -run '^TestOutputCollectionTwoOverflowsKillOnce$' -count=1 -v >"$out" 2>&1
rc=$?; grep -q '^=== RUN   TestOutputCollectionTwoOverflowsKillOnce$' "$out" && grep -q '^--- PASS: TestOutputCollectionTwoOverflowsKillOnce ' "$out" && [ "$rc" -eq 0 ]
```

Expected after change: success with `killCount == 1`. **At base:** fails on the missing `RUN` identity.

### AC4 — the old throughput/stopwatch oracle is removed

Command:

```sh
python3 - <<'PY'
from pathlib import Path
s = Path('host/capsule/capsule_test.go').read_text()
start = s.index('func TestF6OutputCapKillsChildBeyondOnePipeBuffer')
end = s.find('\nfunc ', start + 5)
body = s[start:] if end < 0 else s[start:end]
for forbidden in ('time.Now()', 'time.Since(', 'elapsed >= clock', 'dbl("0123456789abcdef"'):
    assert forbidden not in body, forbidden
PY
```

Expected after change: success. **At base:** fails on all four forbidden timing/throughput tokens. This AC is intentionally scoped to the named test body, so unrelated legitimate timing tests do not satisfy or red it.

### AC5 — caller-supplied releases bound blocking drains and wait

Command:

```sh
out="$(mktemp)"; go test ./host/capsule -run '^TestOutputCollectionCallerReleaseUnblocksReadersAndWait$' -count=1 -v >"$out" 2>&1
rc=$?; grep -q '^=== RUN   TestOutputCollectionCallerReleaseUnblocksReadersAndWait$' "$out" && grep -q '^--- PASS: TestOutputCollectionCallerReleaseUnblocksReadersAndWait ' "$out" && [ "$rc" -eq 0 ]
```

Expected after change: success; both reads are observed blocked before the caller releases them, `Wait` is observed blocked before its independent caller release, and the helper then returns. The test has a bounded failure watchdog for each handshake and final return, so a missing release path or incorrect sequencing fails rather than hanging. It asserts no elapsed-time performance threshold. **At base:** fails because the named test does not exist and therefore the `RUN` grep is false.

### Completion gates (not acceptance evidence by themselves)

After AC1–AC5 and the mutation drill pass, run both repository gates:

```sh
GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh
GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh
```

These are mandatory regression gates, but are not counted as ACs: both are expected green at base and therefore cannot attribute success to this change. `verify_go.sh` supplies build, plain tests, race tests, toolchain canary, and its own non-vacuity controls. `verify_ail.sh` must still run despite no `.ail` edits because it is the mission gate.

## 5. Mutation table

Run every mutant with its named AC's exact discovery assertion, then restore the file byte-identically and rerun the pristine control. Mutants must compile; use boolean neutering where possible so an unused import cannot masquerade as a semantic kill.

| Mutation | Row it must kill | Observable moved | Which write the observable reads |
|---|---|---|---|
| **M1 (neuter overflow kill):** change `if errors.Is(*dstErr, errOutputLimit) {` to `if false && errors.Is(*dstErr, errOutputLimit) {`. | AC1 | `killCount` becomes 0; fake `Wait` returns “waited without kill”; AC1 fails before accepting the error kind. | The observable reads the counter/state written only inside fake child's `Kill`; `Wait` independently reads that killed state. |
| **M2 (force kill):** change the same guard to `if true || errors.Is(*dstErr, errOutputLimit) {`. | AC2 | Exact-limit input changes `killCount` from 0 to 1. | Counter write occurs only in fake `Kill`, invoked by the helper; AC2 reads it after helper return. This is M1's required dual arm. |
| **M3 (deadline wins):** move/check `context.DeadlineExceeded` before the overflow-return branch. | AC1 | Returned type changes from `*OutputLimitError` to `*TimeoutError` with the context pre-expired before the call. | Error kind is written by the helper's return branch; AC1 reads it through `errors.As`. Context state is written by `cancel()` before helper invocation, not beside the branch. |
| **M4 (remove once-only guard):** replace `killOnce.Do(func() { ... })` with the direct kill call while retaining all imports (for example `_ = killOnce; ...`). | AC3 | Two prefilled over-limit streams move `killCount` from 1 to 2. | Each helper invocation of fake `Kill` writes the counter; AC3 reads the final counter after both drains complete. |
| **M5 (truncate without overflow):** in `readCapped`, return `data[:limit], nil` for `len(data) > limit`. | AC1 | Kill count becomes 0 and returned error ceases to be `OutputLimitError`. | `readCapped` writes `dstErr`; the helper's overflow guard reads it. Fake `Kill` is the sole counter writer; helper return is the error-kind writer. |
| **M6 (production helper wiring severed):** after `cmd.Start`, bypass the new helper and return a fixed successful `Result` (retain the helper and imports so the mutant compiles). | Existing `TestF6OutputCapReturnsStructuredOverflow` | Integration result changes from `OutputLimitError` to nil. | The real interpreter writes the pipe in the pristine arm; `readCapped` writes overflow and the helper writes the returned `OutputLimitError`. The mutant bypasses that write, so the existing integration assertion reads nil. This proves `Runner.Run` reaches the tested helper; it does not make a new process-group-kill claim. |
| **M7 (wait before drains finish):** move `runErr := child.Wait()` before the drain wait, leaving the later drain wait in place and retaining all declarations so the mutant compiles. | AC5 | The fake `Wait` is observed before both blocked readers have been released/completed; AC5 fails its ordered handshake (and its watchdog prevents a suite hang). | Reader-entry signals are written only by the blocking readers; wait-entry is written only by fake `Wait`. AC5 reads those independent events before releasing each phase. This kills loss of the current drain-before-wait sequencing and proves the liveness contract is exercised rather than assumed. |

M6 is a design warning as much as a drill: an injected fake proves the algorithm only if production wiring remains covered. The executor must not record M6 as killed by compilation failure, global timeout, missing test discovery, or a skipped capsule test. Its command asserts `=== RUN`, `--- FAIL`, zero `SKIP`, and uses the pinned environment above.

## 6. Implementation sequence and cost

1. Extract the unexported output-collection helper without changing `Runner.Run` behavior, and document its caller-owned liveness precondition beside the helper.
2. Replace the named flaky timing fixture with the three deterministic state-machine tests, add the controlled blocking-reader/blocking-wait liveness arm, and retain the existing real-interpreter structured-overflow integration test.
3. Run AC1–AC5 at base/head as specified and execute M1–M7 with byte-identical restoration controls.
4. Run `gofmt`, the focused package tests with explicit discovery/skip assertions, then both full gates with `GOTOOLCHAIN=go1.25.6` and `AILANG_BIN=/tmp/ailang-v0300/ailang`.

Estimated effort is **1.0 day**: 0.25 day extraction and finite state-machine tests, 0.25 day controlled blocking fixtures plus liveness/debugging allowance, 0.25 day mutation drill and restoration controls, and 0.25 day full plain/race gates and remediation allowance. Half 1 makes the item larger than 0.75 day because the blocking protocol must be synchronized, watchdog-bounded, race-clean, and mutation-proven. The filed 0.25 day estimate fits only a deadline edit.

## 7. Non-goals and residuals

- No process-group-wide overflow kill change; that named neighbouring queue item owns it.
- No `host/broker` change.
- No total-`Run` latency budget and no redefinition of `ExecTimeout`.
- No cache or hoist of interpreter verification and no weakening of hash-before-exec.
- No claim that the historical 1-of-2 failure is reproduced, assigned a stable rate, or fixed by load sampling.
- No `.ail`, Z3, kernel, package, CLI, or documentation-surface change beyond this design and the later Go implementation.

**RESIDUAL CORRECTED POST-SPRINT (iteration 100, evaluator `sonnet` 93/100, reproduced first-party
by the controller).** This document filed the surviving integration test's stale NAME as the residual
of retiring the oracle. That undersells it, and the correction matters because it is the same shape
this mission keeps recording: **the retired fixture was doing TWO jobs and only one of them was bad.**
Job 1 — a wall-clock kill oracle started before `New(...).Run(...)`, i.e. outside the region
`ExecTimeout` governs — is the flake, and retiring it is this item's whole point. Job 2 — emitting
64 KiB so the child genuinely BLOCKED in `write()` on a full pipe — was real coverage, and it left
with the oracle because nobody separated them. Measured: this rig's pipe capacity is **65536 bytes**
(a Go probe writing to an undrained `os.Pipe` blocks after 65536); the surviving
`TestF6OutputCapKillsChildBeyondOnePipeBuffer` emits `repeat(32)` = 32x16 + newline = **513 bytes**
against a **64-byte** cap; and `readCapped` reads only `limit+1` bytes before stopping, so filling the
pipe requires total output to exceed its remaining capacity, which 513 bytes cannot. **So no arm in
the final suite — fake or real-interpreter — drives a live child blocked in `write()` being unblocked
by `cmd.Process.Kill()`**, which is exactly the scenario `capsule.go`'s own "F6 must not decay into
F5" comment describes. Non-blocking for the landed sprint, and that is measured rather than asserted:
M-A is behaviour-preserving, confirmed by M6's kill set matching the predicted 7 tests across 2
packages exactly, so the kill path is untouched. Owner: **queue row 25
`w-capsule-blocked-child-kill-coverage`**, filed by iteration 100 in the same charter edit as this
correction. NOT queue row 24, which owns the cleanup boundary and the non-group-wide overflow kill —
a different mechanism. The generalisable point, worth more than the instance: **when you retire an
instrument, enumerate every property it was supplying, not just the one that made you retire it.**

The following lifecycle debt is **owed and not absorbed here**: a production boundary such as `CloseOutput() error`, `Wait(context.Context) error`, an explicit bounded cleanup context, and `errors.Join` surfacing of kill/close/wait failures. Its named owner is **queue row 24 `w-host-subprocess-cleanup-boundary`**, filed by iteration 99 in the same charter edit as this document's routing. (Queue row 21 *declared* this residual but LANDED on 2026-08-20, so it cannot own follow-on work — naming a completed row as an owner is exactly the defect queue row 23 exists to record.) The related non-group-wide overflow-kill defect remains owned by its separately filed queue row. Adding either behavior here would change `host/capsule` subprocess lifecycle and reorder separately owned work, so this tranche only states and tests the existing caller-supplied liveness contract. This is an explicit residual, not a claim that arbitrary readers or `Wait()` are intrinsically bounded.

## 8. Quorum verification log

Revision round 1 had both reviewers present; there is no absent-reviewer hole. `gemini-3-1-pro` returned a blocking **PREMISE** objection, and `gpt5-6-sol` returned a blocking **DESIGN** objection; pick-time quorum was therefore blocked (metered cost `$0.048867`).

- Applied the premise reviewer's proposed fix verbatim in substance and in both requested parts: the Premise Verification Log now proves `killOnce.Do`, `readCapped`, and `errOutputLimit` with the supplied same-path negative and positive controls, and the dedicated **Conflict Surface** explains why the existing `io.LimitReader` plus context primitives do not supply the composed state machine.
- Applied design objection half 1 in full: the seam contract names the caller as liveness-bound owner; AC5 contains blocking readers and blocking `Wait`, caller-controlled releases, discovery assertions, and a bounded failure escape; M7 kills loss of the sequencing that makes the arm meaningful.
- Declared design objection half 2 as an owned residual rather than absorbing or dropping it: **queue row 24 `w-host-subprocess-cleanup-boundary`** owns `WaitDelay`/bounded cleanup and the requested `CloseOutput`, context-aware wait, cleanup context, and joined cleanup errors; the separately filed overflow-kill row retains its own scope. **(Round 2's `gpt5-6-sol` `catch` was correct that this bullet said row 21 while §7 said row 24 — row 21 LANDED on 2026-08-20 and cannot own follow-on work. Corrected here; §7 was the right half. Row 24 re-asserted OPEN by command at round 3, with LANDED row 21 as the firing control.)**

### Round 3 — resolved by ratified ruling, not by a third quorum

Round 2 blocked on two objections, both with `absent_reviewers` empty. Neither survives, and neither
was re-litigated:

- **`gpt5-6-sol` (DESIGN/SCOPE) — answered by `D-WORLD-23` arm A, ratified attended by Mark
  2026-08-20T08:01:31Z (`#68`, one-word comment `A`).** Its `proposed_fix` folds a bounded cleanup
  path — `CloseOutput`, a context-aware wait, an explicit cleanup context, `WaitDelay`, joined
  cleanup errors — into this tranche. That is the exact shape the ruling now governs: **the tranche
  keeps scope, weakens its claim to precisely what it proves, and records the residual with a named
  OPEN owner.** §7 already does all three; the ruling makes it a disposition rather than an ask. The
  objection is SUSTAINED as substantively correct — the helper's drains and `Wait()` are not
  intrinsically bounded, and this document does not claim they are. What it claims, and what its
  arms test, is narrower: the existing caller-supplied liveness contract, and the state machine
  composed over it. Production liveness continues to come from `ExecTimeout` cancellation and the
  group SIGKILL, measured in §2 and unchanged by this tranche.
- **`gemini-3-1-pro` (PREMISE) — MEASURED first-party rather than forwarded (rule 3f), all three
  claims TRUE, and its `proposed_fix` applied verbatim in both requested parts.** A new §2 row
  establishes `TestF6OutputCapReturnsStructuredOverflow` (`capsule_test.go:233`) and the
  `io.LimitReader` call, each with same-path positive and negative controls in the same call; the
  final §2 row's establishing command now carries
  `grep -nE 'stdoutPipe|stderrPipe' host/capsule/capsule.go`, exactly as asked, so its Scope column's
  sub-claim is established by command instead of by prose. **The measurement carried a correction the
  objection did not predict:** `io.LimitReader` is at `:238`, while §3 cited `:237` — that is the
  `func readCapped(` declaration line, transcribed from the adjacent premise row's correct citation
  of the FUNCTION. §3 is corrected. This is the "a value you transcribed is not a measurement" class,
  and it was invisible until the command the reviewer asked for was actually run.

**No third quorum was purchased.** Under the ratified ruling the scope objection has a standing
disposition, and the premise objection is narrow-refinement carve-out eligible on its own terms
(concrete reviewer-authored `proposed_fix`, no dispute of design DIRECTION). **Iteration 98's
guardrail — a carve-out fix must acquire its own acceptance criterion and mutation before routing —
is checked and does not bind here, which is stated rather than assumed:** every round-3 edit is
DOCUMENTATION ONLY (three premise-log facts, one line number, one owner name). No acceptance
criterion, mutation, seam, or line of Go changes, so there is no new behavior for an AC to cover and
nothing for a mutant to neuter. The guardrail binds a fix that changes what the code DOES; this one
changes only what the document can PROVE it already knew.

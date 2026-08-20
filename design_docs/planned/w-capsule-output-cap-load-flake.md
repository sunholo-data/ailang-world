# w-capsule-output-cap-load-flake

**Queue:** World mission row 20  
**Status:** planned  
**Scope:** `host/capsule` Go host code and tests only; no `.ail` change  
**Estimate:** 0.75 day (the filed 0.25 day is too small for a non-vacuous redesign, mutation drill, and both full gates)

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

## 3. Chosen design: deterministic output-collection core

Extract the post-`cmd.Start` output lifecycle from `Runner.Run` into an unexported helper. The helper receives:

- stdout and stderr `io.Reader`s;
- the byte limit;
- the already-created execution context;
- a narrow child boundary with `Kill() error` and `Wait() error` (or equivalently injected kill/wait functions).

It owns the existing two drains, the once-only kill on `errOutputLimit`, waiting, and the error-precedence decision. `Runner.Run` remains responsible for archive resolution, verification, staging, context creation, command construction, pipes, and `Start`; after `Start` it adapts the real process to the helper. Production semantics and exported APIs do not change.

Replace the timing-bearing assertion in `TestF6OutputCapKillsChildBeyondOnePipeBuffer` with deterministic package-local tests of this helper:

1. **Overflow + expired context:** stdout is a `bytes.Reader` containing exactly `limit+1` bytes; stderr is empty. A fake child's `Kill` increments a counter and marks the fake killed. Its `Wait` reads that independently written killed state and returns an explicit “waited without kill” sentinel if false. Pass an already-expired context. Assert captured stdout is exactly `limit`, `Kill` was called exactly once, `Wait` observed the kill, and the result is `*OutputLimitError`, never `*TimeoutError`. No goroutine writes bytes, no pipe fills, and no duration is asserted.
2. **Within-limit control:** provide exactly `limit` bytes and a live context. Assert byte identity, zero kills, one wait, and the fake wait result. This makes an unconditional-kill implementation fail.
3. **Two-stream once control:** make both prefilled readers exceed the limit. Assert the child boundary records exactly one kill. This guards the `sync.Once` behavior without scheduler timing.

Keep `TestF6OutputCapReturnsStructuredOverflow` as the real interpreter integration coverage with its small output. It proves production wiring reaches an `OutputLimitError`; the new tests prove the beyond-buffer kill/precedence algorithm. Delete the old outer stopwatch and the 64 KiB throughput fixture from the named test. The test no longer claims elapsed wall time measures `ExecTimeout`.

The fake's observables are deliberately written at the child boundary, not alongside the assertion: the helper calls `Kill`; the fake `Kill` method writes `killCount` and killed state; the fake `Wait` method independently reads killed state; assertions read the count and wait result. The output observable reads bytes produced before helper invocation. Thus the test cannot manufacture “kill succeeded” in the same branch that asserts it.

### Why the alternatives lose

**Bound or hoist verification.** Moving `context.WithTimeout` above verification would redefine `ExecTimeout` from execution time to total preparation-plus-execution time and make a slow disk/hash legitimately beat output. Hoisting verification into `New` requires an API error path, caching policy, and a TOCTOU analysis: verification must still bind the bytes eventually executed. Either is a larger security/semantics item and neither makes a real child's time-to-first-overflow deterministic. The unbounded phase deserves a separate budget/API item if total request latency must be bounded.

**Measure only the governed region.** Moving `start` to immediately before `cmd.Start`, or returning phase timestamps, removes the current contamination. It still asserts scheduler- and interpreter-dependent elapsed time against a wall clock, and instrumentation added solely for the test expands production surface. It diagnoses a race rather than making the property deterministic.

**Raise the deadline or enlarge output.** Both merely change odds. A larger deadline attacks the term with the greatest measured margin; more output still depends on interpreter production, pipe behavior, and scheduling. Neither can support a non-vacuous base-red gate.

The selected seam tests the causal state machine with pre-existing bytes and an independently observed kill. Load can delay the test process itself, but cannot change which input is available or which branch is correct; there is no competing real-time deadline.

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

### Completion gates (not acceptance evidence by themselves)

After AC1–AC4 and the mutation drill pass, run both repository gates:

```sh
GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh
GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh
```

These are mandatory regression gates, but are not counted as ACs: both are expected green at base and therefore cannot attribute success to this change. `verify_go.sh` supplies build, plain tests, race tests, toolchain canary, and its own non-vacuity controls. `verify_ail.sh` must still run despite no `.ail` edits because it is the mission gate.

## 5. Mutation table

Run every mutant with AC1–AC3's exact discovery assertion, then restore the file byte-identically and rerun the pristine control. Mutants must compile; use boolean neutering so an unused import cannot masquerade as a semantic kill.

| Mutation | Row it must kill | Observable moved | Which write the observable reads |
|---|---|---|---|
| **M1 (neuter overflow kill):** change `if errors.Is(*dstErr, errOutputLimit) {` to `if false && errors.Is(*dstErr, errOutputLimit) {`. | AC1 | `killCount` becomes 0; fake `Wait` returns “waited without kill”; AC1 fails before accepting the error kind. | The observable reads the counter/state written only inside fake child's `Kill`; `Wait` independently reads that killed state. |
| **M2 (force kill):** change the same guard to `if true || errors.Is(*dstErr, errOutputLimit) {`. | AC2 | Exact-limit input changes `killCount` from 0 to 1. | Counter write occurs only in fake `Kill`, invoked by the helper; AC2 reads it after helper return. This is M1's required dual arm. |
| **M3 (deadline wins):** move/check `context.DeadlineExceeded` before the overflow-return branch. | AC1 | Returned type changes from `*OutputLimitError` to `*TimeoutError` with the context pre-expired before the call. | Error kind is written by the helper's return branch; AC1 reads it through `errors.As`. Context state is written by `cancel()` before helper invocation, not beside the branch. |
| **M4 (remove once-only guard):** replace `killOnce.Do(func() { ... })` with the direct kill call while retaining all imports (for example `_ = killOnce; ...`). | AC3 | Two prefilled over-limit streams move `killCount` from 1 to 2. | Each helper invocation of fake `Kill` writes the counter; AC3 reads the final counter after both drains complete. |
| **M5 (truncate without overflow):** in `readCapped`, return `data[:limit], nil` for `len(data) > limit`. | AC1 | Kill count becomes 0 and returned error ceases to be `OutputLimitError`. | `readCapped` writes `dstErr`; the helper's overflow guard reads it. Fake `Kill` is the sole counter writer; helper return is the error-kind writer. |
| **M6 (production helper wiring severed):** after `cmd.Start`, bypass the new helper and return a fixed successful `Result` (retain the helper and imports so the mutant compiles). | Existing `TestF6OutputCapReturnsStructuredOverflow` | Integration result changes from `OutputLimitError` to nil. | The real interpreter writes the pipe in the pristine arm; `readCapped` writes overflow and the helper writes the returned `OutputLimitError`. The mutant bypasses that write, so the existing integration assertion reads nil. This proves `Runner.Run` reaches the tested helper; it does not make a new process-group-kill claim. |

M6 is a design warning as much as a drill: an injected fake proves the algorithm only if production wiring remains covered. The executor must not record M6 as killed by compilation failure, global timeout, missing test discovery, or a skipped capsule test. Its command asserts `=== RUN`, `--- FAIL`, zero `SKIP`, and uses the pinned environment above.

## 6. Implementation sequence and cost

1. Extract the unexported output-collection helper without changing `Runner.Run` behavior.
2. Replace the named flaky timing fixture with the three deterministic tests and retain the existing real-interpreter structured-overflow integration test.
3. Run AC1–AC4 at base/head as specified and execute M1–M6 with byte-identical restoration controls.
4. Run `gofmt`, the focused package tests with explicit discovery/skip assertions, then both full gates with `GOTOOLCHAIN=go1.25.6` and `AILANG_BIN=/tmp/ailang-v0300/ailang`.

Estimated effort is **0.75 day**: 0.25 day extraction and tests, 0.25 day mutation drill and instrument controls, 0.25 day full plain/race gates and remediation allowance. The filed 0.25 day estimate fits only a deadline edit; it does not cover the deterministic seam and non-vacuous evidence required by the measurements and coding standard S6.

## 7. Non-goals and residuals

- No process-group-wide overflow kill change; that named neighbouring queue item owns it.
- No `host/broker` change.
- No total-`Run` latency budget and no redefinition of `ExecTimeout`.
- No cache or hoist of interpreter verification and no weakening of hash-before-exec.
- No claim that the historical 1-of-2 failure is reproduced, assigned a stable rate, or fixed by load sampling.
- No `.ail`, Z3, kernel, package, CLI, or documentation-surface change beyond this design and the later Go implementation.

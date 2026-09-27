# Sprint plan — w-capsule-blocked-child-kill-coverage (iteration 200)

**Authority:** World mission queue row 25, clause-2, and [row 20 §7, “RESIDUAL CORRECTED POST-SPRINT”](w-capsule-output-cap-load-flake.md). **Scope:** one test-only milestone, about 0.5 day. **Landing decision:** one commit after M1 and its mutations pass. No new design document or production edit is planned. The planner made no git writes; only this plan and its JSON manifest remain changed.

Use this prefix for every executor command below:

```sh
export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang
export PATH="$(dirname "$AILANG_BIN"):$PATH"
```

## 1. Pristine measurements and boundary

On this worktree, a throwaway Go probe opened an undrained `os.Pipe`, set the write descriptor nonblocking, wrote 4096-byte chunks, and counted **positive** write returns until `EAGAIN`. It measured **65,536 bytes** (a first draft of the probe incorrectly added `-1` from the `EAGAIN` return and printed 65,535; corrected and rerun). The probe was deleted. Re-measure in the new test, because the bound must scale with the pipe on the machine actually running it. `readCapped` uses `io.LimitReader(pipe, limit+1)`; for a 64-byte limit it drains at most 65 bytes. The existing real-interpreter `repeat(32)` output is 513 bytes, less than the measured pipe capacity, and therefore cannot witness the blocked writer.

| Pristine command after the prefix | Reading |
|---|---|
| `go vet ./host/capsule` | rc=0 |
| `go test ./host/capsule -count=1 -v` | rc=0; 23 top-level `PASS`, no skips; package 25.782 s |
| `go test ./host/capsule -count=5 -run '^TestBlockedChildOverflowKill$'` | rc=0 **[no tests to run]**; prerequisite absent at base, not a green stability result |
| `go test ./host/capsule -count=1 -race` | rc=0; package 32.061 s |
| `./scripts/verify_ail.sh` | rc=0; 16/16 identities, 40 named tests, 9/9 world-package steps, PUB011 ratchet passed |

No baseline gate is red. `host/capsule` does not bind sockets, so these readings are informative here. If another gate later fails only at `bind: operation not permitted`, label that result **UNINFORMATIVE UNDER SANDBOX** and rerun outside; do not count it as pass or fail.

## 2. M1 — child-owned blocked-write witness

Add `TestBlockedChildOverflowKill` in `host/capsule/cleanup_test.go` (or a nearby test file). Use the existing `archiveExecutable` / scripted-interpreter fixture pattern: a small archived `/bin/sh` interpreter answers `--version`, then `exec`s a helper mode of the current Go test executable. The helper is a **real child process** in the capsule's process group, not a `collectOutput` fake. Its helper mode records its own PID to a file, makes **one blocking `write(1, payload)` syscall**, records `short-write` if that syscall returns fewer than the requested bytes, and records `wrote-all` only after the entire syscall succeeds. An unexpected write error must record a separate failure marker. Exit directly from helper mode so Go test harness output cannot become the stimulus. Embed safely quoted marker path, payload length, and helper executable path in the generated shell wrapper; set a helper-mode environment variable there, because `Runner.Run` replaces `cmd.Env` with its two sandbox variables. The ordinary test mode returns without side effects.

At test start, measure undrained `os.Pipe` capacity using nonblocking writes until `EAGAIN` and close both ends. Reject an invalid/overflowed measurement. Set `limit = 64` and compute payload length as at least `4 * measuredCapacity`, with an explicit lower bound greater than `measuredCapacity + limit + 1`; on this rig this is **262,144 bytes** versus at most **65,601 bytes** that the pipe plus `readCapped` can accept. Print the measured capacity and selected size in test diagnostics. The child's single write cannot complete while the reader stops after 65 bytes: there is nowhere for the remaining bytes to go. The child's PID marker is created before the write; receipt of more than 64 output bytes proves it entered that syscall. On normal return there can be neither a `wrote-all` nor a `short-write` marker. Those child-owned facts and byte bounds, rather than duration, witness a live child blocked inside `write()` when overflow kill occurs. A short successful write is an explicit test failure; it must never count as the blocked case.

Run with `Config{MaxOutputBytes: limit, ExecTimeout: 120 * time.Second}`. Assert `errors.As(err, *OutputLimitError)` and **no** `*TimeoutError`; both captured streams are at most `limit` bytes; the PID file exists; `wrote-all`, `short-write`, and the unexpected-error marker are absent; `proctest.DeadWithin(t, pid, time.Second)` (or equivalent) confirms the child is dead. Do **not** assert on `killErr` or exact error identity: row 115 owns the approximately 1% `EPERM` ambiguity in `TestOrdinaryOverflowCarriesNoESRCH`, and this test must not inherit that flake. Row 117 is also out of scope.

Start `Run` in a test goroutine and use an independent, generous **30-second liveness watchdog**, well before `ExecTimeout=120s`. If the watchdog fires, signal the recorded child process group with an independent `syscall.Kill(-pid, SIGKILL)` cleanup path, reap/finish the goroutine, then **fail the test unconditionally**. The watchdog is only a bounded failure escape for M-kill-off; elapsed time is never a pass criterion or kill oracle. Use a deferred cleanup path as well so a failed assertion cannot leave a writer alive. Avoid `t.Fatal` in the worker goroutine. Choose a watchdog large enough for ordinary CI startup; it must remain below the F5 deadline. The executor should confirm its cleanup path under the no-kill mutation.

The existing `TestF6OutputCapKillsChildBeyondOnePipeBuffer` stays as a real-interpreter wiring control. A second real-interpreter blocked-write arm is **feasible in principle** if the language program can emit a capacity-scaled payload and expose its own PID and completion marker under the capsule sandbox, but is **not worth this half-day sprint**: it adds language/FS behavior and less direct child-state evidence. The scripted real-process arm above exercises production `Runner.Run`, pipes, cap, and kill. If the helper executable or sandbox prevents this design, measure that failure first, then use a tiny compiled Go helper as the archived fake interpreter; do not weaken the blocked-write witness or quietly add production code.

## 3. Mutation ledger

After the positive M1 test passes, apply **one mutation at a time** in an isolated worktree-local copy or by saving the touched file's original bytes. Record `shasum -a 256` for each touched file before mutation; run the named killer with `go test ./host/capsule -count=1 -run '^TestBlockedChildOverflowKill$'`; restore by byte copy and require the same SHA-256 before the next mutation. No mutation may remain in the landing diff. A mutation only counts as killed when the **new named test** fails for the stated child-state reason, not because the harness could not start.

| Mutation | One required failing test and diagnostic |
|---|---|
| **M-kill-off:** make the overflow `killOnce` body / `killGroup` a no-op, leaving the 120 s F5 deadline as the only production kill. | `TestBlockedChildOverflowKill` must hit its 30 s failure watchdog while the child has a PID but no completion marker; independent cleanup kills and reaps it. A return caused by F5 cannot count. |
| **M-drain-all:** change `readCapped` to drain to EOF instead of `limit+1` before deciding overflow. | `TestBlockedChildOverflowKill` must fail: the child finishes the one write and records `wrote-all` (and the result may also lose its typed overflow). Ensure the marker is checked even if another assertion already fails. |
| **M-kill-direct-only:** replace group kill with direct PID kill. | **Not a required killer for this test.** The tested writer is the direct child, so this mutation can pass; `TestOverflowKillReachesForkedGrandchild` from row 24 owns group reach. Do not claim this mutation is killed by the new test. |
| **513-byte control:** replace the new capacity-scaled payload with the existing 513-byte `repeat(32)` scale, then run M-kill-off. | `TestBlockedChildOverflowKill` must fail its blocked-write invariant (for example, `wrote-all` appears); 513 bytes must not satisfy the new witness. This control establishes that the old arm cannot stand in for M1. |

## 4. Acceptance and scope

After restoring all mutants, run these exact gates with the prefix above: `go vet ./host/capsule`; `go test ./host/capsule -count=1`; `go test ./host/capsule -count=5 -run '^TestBlockedChildOverflowKill$'`; `go test ./host/capsule -count=1 -race`; `./scripts/verify_ail.sh`. Require five actual executions of the new test in the stability gate, not merely rc=0 with `[no tests to run]`. Controller repeats outside the sandbox. Inspect `git diff --check` and the final changed-file list. Only the new Go test file should change in the executor's implementation; this sprint does not change `capsule.go`, the group-kill mechanism, row 115's `EPERM` handling, or row 117. If a production change proves necessary, report the failing measurement and re-plan rather than silently expanding M1.

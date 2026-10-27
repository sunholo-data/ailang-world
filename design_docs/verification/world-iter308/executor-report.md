# M7a-fix executor report — iteration 308, row 150

Only two source files changed:

- `scripts/marketing_capture.py`: only the prescribed exception tuple and one-line macOS comment in `kill_group`.
- `scripts/marketingcheck/supervisor_test.go`: added `TestMarketingKillGroupTolerance` and `TestMarketingKillGroupStress`.

No git writes, other source/doc/plan changes, real Chrome, real daemon, or external network use. Read-only HEAD check: `97736493`. The whole pre-existing Go test body, starting at `func pythonSupervisor`, was compared byte-for-byte with a scratch copy of base and is unchanged, including both required regression tests. Production bytes were checked against the exact prescribed replacement.

**UNINFORMATIVE UNDER SANDBOX:** every runtime result below involving killpg, signalling, or sockets has this label, including the stubbed killpg UNIT observations under the directive's blanket rule. Quoted PASS/FAIL strings are literal test output, not outside-sandbox acceptance. UNIT uses only injected exceptions and does not send real signals. Static vet/format/compile checks are independent of process-signalling acceptance.

## Red first

Wrote the tolerance test before editing production, then ran:

```sh
go test -run '^TestMarketingKillGroupTolerance$' ./scripts/marketingcheck -count=1 -timeout 90s -v
```

Observed rc **1**; only `TestMarketingKillGroupTolerance/eperm-swallowed` emitted FAIL. Exact assertion line:

```text
supervisor_test.go:104: kill_group must swallow PermissionError EPERM: {Case:eperm-swallowed Calls:[[424242 9]] Outcome:propagated Type:PermissionError Errno:0x5f06dac52748 Same:true}
```

`esrch-swallowed` and `einval-propagates` both emitted PASS. After the production edit, the same command returned rc **0**, with the exact top-level `--- PASS: TestMarketingKillGroupTolerance` and all three subtest PASS lines. UNINFORMATIVE UNDER SANDBOX as specified above.

## Gate ledger

| Command | rc | Key output / observation |
|---|---:|---|
| `go test -race ./scripts/marketingcheck -count=1 -timeout 240s` (unchanged base, before any edit) | 1 | `TestMarketingCapturePassFakeChrome/ready`; `panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: bind: operation not permitted`; package 9.291s. UNINFORMATIVE UNDER SANDBOX. |
| `go test -run '^TestMarketingKillGroupTolerance$' ./scripts/marketingcheck -count=1 -timeout 90s -v` (red first) | 1 | EPERM assertion quoted above; ESRCH/EINVAL controls emitted PASS. UNINFORMATIVE UNDER SANDBOX. |
| `go test -run '^TestMarketingKillGroupTolerance$' ./scripts/marketingcheck -count=1 -timeout 90s -v` (fixed, focused gate) | 0 | `--- PASS: TestMarketingKillGroupTolerance (0.05s)` and all three subtests; package 0.219s. UNINFORMATIVE UNDER SANDBOX. |
| `go test -run '^TestMarketingKillGroupStress$' ./scripts/marketingcheck -count=1 -timeout 120s -v` (environment variable absent) | 0 | `--- SKIP: TestMarketingKillGroupStress (0.00s)`; exact skip text below. Repeated after adding the final deadline-overrun guard, also rc 0 and the same SKIP. No acceptance claimed. |
| `go test -run '^TestMarketingSupervisorFakeChromeTimeoutKillsGroup$' ./scripts/marketingcheck -count=1 -timeout 90s -v` | 0 | `--- PASS: TestMarketingSupervisorFakeChromeTimeoutKillsGroup (8.09s)`; package 8.397s. UNINFORMATIVE UNDER SANDBOX. |
| `go test -run '^TestMarketingSupervisorFakeDaemonCleanup$' ./scripts/marketingcheck -count=1 -timeout 120s -v` | 0 | `--- PASS: TestMarketingSupervisorFakeDaemonCleanup (7.75s)`; all 13 subtests emitted PASS, including interrupted-shot and shell-hup; package 8.195s. UNINFORMATIVE UNDER SANDBOX. |
| `go vet ./scripts/marketingcheck` | 0 | Empty output; repeated after final stress deadline guard, rc 0 and empty output. |
| `gofmt -l scripts/marketingcheck` | 0 | Empty output; repeated after final stress deadline guard, rc 0 and empty output. |
| `python3 -m py_compile scripts/marketing_capture.py` | 0 | Empty output. Removed the newly generated `scripts/__pycache__/marketing_capture.cpython-314.pyc` and now-empty `scripts/__pycache__/`. |
| `go test -race ./scripts/marketingcheck -count=1 -timeout 240s` (fixed) | 1 | Same base socket refusal in `TestMarketingCapturePassFakeChrome/ready`: `panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: bind: operation not permitted`; package 9.379s. UNINFORMATIVE UNDER SANDBOX. |
| `WORLD_KILLPG_STRESS=1 go test -run '^TestMarketingKillGroupStress$' ./scripts/marketingcheck -count=1 -timeout 120s -v` | not run | Controller-owned outside-sandbox macOS run. No raw EPERM frequency or stress acceptance claimed. |

Exact default stress skip message:

```text
opt-in: set WORLD_KILLPG_STRESS=1; real process signalling must run outside the sandbox to probe macOS EPERM
```

Additional static verification: both inline Python programs parsed successfully using `ast.parse` without executing them; the existing Go test body and the exact production replacement passed byte comparisons. Initial source-status inventory contained no Python bytecode; cleanup removed only newly created bytecode.

Historical controller inputs, not measurements made here: pristine formatting/vet/compile rc 0; non-race whole package rc 0 (24.6s); race red 1 of 3 at base in `TestMarketingSupervisorFakeDaemonCleanup`, with `{"verdict": "SupervisorError", "detail": "[Errno 1] Operation not permitted"}`. The local base and fixed race attempts instead stopped earlier at sandbox socket binding. This is a reported gate limitation; no workaround or unrelated edit was made. The separate non-race whole-package command is not in the milestone's focused gate block and was not rerun.

## Mutant ledger

Commands abbreviated in this table, expanded exactly here:

```sh
# UNIT
go test -run '^TestMarketingKillGroupTolerance$' ./scripts/marketingcheck -count=1 -timeout 90s -v
# TIMEOUT
go test -run '^TestMarketingSupervisorFakeChromeTimeoutKillsGroup$' ./scripts/marketingcheck -count=1 -timeout 90s -v
```

All edits were to production `kill_group` only, with tests untouched during each mutant. No compile failure, panic, harness timeout, missing interpreter, or SKIP is counted as an assertion observation. All runtime mutant observations are UNINFORMATIVE UNDER SANDBOX under the directive; the real-process TIMEOUT leg specifically requires controller repetition.

| Mutant | Command / observed rc | Named red assertion line | Restore |
|---|---|---|---|
| `M7a-fix-eperm`: narrow handler back to `except ProcessLookupError:` | UNIT / 1 | `supervisor_test.go:104: kill_group must swallow PermissionError EPERM: {Case:eperm-swallowed Calls:[[424242 9]] Outcome:propagated Type:PermissionError Errno:0x71bf8ecda748 Same:true}` | Copied `snap/M7a-fix/scratch/eperm-fixed-marketing_capture.py.txt` back; `cmp` rc 0, empty output; restored UNIT rc 0 with all three subtests emitting PASS. |
| `M7a-fix-broad`: `except OSError:` | UNIT / 1 | `supervisor_test.go:112: kill_group must propagate other OSError EINVAL unchanged: {Case:einval-propagates Calls:[[424242 9]] Outcome:swallowed Type: Errno:<nil> Same:false}` | Copied `snap/M7a-fix/scratch/broad-fixed-marketing_capture.py.txt` back; `cmp` rc 0, empty output; restored UNIT rc 0 with all three subtests emitting PASS. |
| `M7a-fix-noop`: body becomes `pass` | UNIT / 1 | `supervisor_test.go:99: kill_group must call killpg once with the expected pid and signal: {Case:eperm-swallowed Calls:[] Outcome:swallowed Type: Errno:<nil> Same:false}`; same call assertion also fired independently in ESRCH and EINVAL subtests | Copied `snap/M7a-fix/scratch/noop-fixed-marketing_capture.py.txt` back; `cmp` rc 0, empty output; restored UNIT rc 0 with all three subtests emitting PASS. |
| `M7a-fix-noop`: same production mutant, additional process leg | TIMEOUT / 1 | `supervisor_test.go:498: grandchild must be dead after bounded cleanup (pid 79776)` | Same byte-identical restore; TIMEOUT rerun recorded below. UNINFORMATIVE UNDER SANDBOX. |

The no-op TIMEOUT leg reached the required grandchild assertion, beyond the earlier nonfatal `CleanupFailed` verdict error (`Chrome could not be reaped after KILL`). It also printed `child must be dead after bounded cleanup (pid 79366)` and completed in 23.26s; package 26.339s. The earlier no-op process observation independently reached `supervisor_test.go:496: grandchild must be dead after bounded cleanup (pid 76264)` (before the two-line stress deadline guard was added), also rc 1. These are actual assertion observations, not an inference from the earlier verdict mismatch. Both are UNINFORMATIVE UNDER SANDBOX.

In the subsequent round, the restored TIMEOUT rerun returned rc **1**, with the distinct fixture assertion `supervisor_test.go:495: fake must record child and grandchild: [81876]` (test 8.17s, package 12.589s). This is a reported finding, not a mutant kill or a successful restoration test. Production `cmp` remained rc 0, and restored UNIT rc 0 with all three subtests emitting PASS. An isolated unchanged TIMEOUT rerun was started to check this finding; its result is recorded below. UNINFORMATIVE UNDER SANDBOX.

The isolated unchanged TIMEOUT rerun returned rc **0**: `--- PASS: TestMarketingSupervisorFakeChromeTimeoutKillsGroup (8.80s)`; package 9.178s. The prior rc 1 remains a finding. At the isolated rerun, `uptime` reported load averages **9.08 8.26 9.00**; the attempted read-only process listing was denied (`operation not permitted: ps`). Load is context, not a proven explanation of the missing PID. All process observations remain UNINFORMATIVE UNDER SANDBOX.

An initial complete mutant round used one immutable fixed scratch copy: EPERM UNIT rc 1 / restored rc 0; broad UNIT rc 1 / restored rc 0; no-op UNIT rc 1 and TIMEOUT rc 1 / restored UNIT rc 0 and TIMEOUT rc 0 (`--- PASS: TestMarketingSupervisorFakeChromeTimeoutKillsGroup (8.10s)`, package 8.583s). Every restore `cmp` returned 0. The table above records the subsequent round with a separate freshly saved scratch file for each mutant.

## Deviations, findings, and remaining controller work

- Both base and fixed package race attempts are blocked by sandbox socket permissions. Reported as a finding; no related source or tests were modified.
- The later restored timeout fixture recorded only one expected descendant PID, as quoted above. The test was preserved unchanged; the earlier successful restoration result is retained separately, without relabeling this rc 1.
- During the additional per-mutant-scratch round, one EPERM repeat was restored while its tool command was still pending; it returned rc 0 and is discarded as contaminated, not counted as a mutant result. The EPERM mutation was then repeated sequentially to completion, reached the named assertion at rc 1, restored with `cmp` rc 0, and reran at rc 0. The original complete round was unaffected.
- The last stress-test refinement makes an elapsed wall time greater than 25 seconds explicitly fail, even if the last cleanup crossed the deadline after completing all trials. Static Python parsing, vet, formatting, and the default skip were rechecked afterward.
- Opt-in stress and the 25-consecutive-run macOS race campaign remain controller work. No sandbox result establishes process death or proves that EPERM means no live member remains. Swallowing a genuinely denied live member remains the plan's stated limitation.
- Final snapshots preserve relative paths; the Go snapshot has `.go.txt` suffix to avoid package discovery. Only the two source files above are source changes; report, scratch copies, and snapshots are under `snap/M7a-fix/`.

Final verification: both snapshot/source `cmp` commands returned rc 0; the entire existing Go test body and exact prescribed Python diff were verified again; no `__pycache__` remained under `scripts/`; no `.go` files exist under the snapshot tree; `git diff --check` returned rc 0. Read-only final `git diff --stat` lists exactly the two authorized source files (Python 3 changed lines, Go 237 added lines). Snapshot source paths are `snap/M7a-fix/scripts/marketing_capture.py` and `snap/M7a-fix/scripts/marketingcheck/supervisor_test.go.txt`.

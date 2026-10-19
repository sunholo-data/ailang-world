Implemented and snapshotted M3a and M3b. Acceptance remains incomplete: the required focused gate reproduces an existing failure at base, and real srt is unavailable.

| Milestone | File | + / − LOC |
|---|---|---:|
| M3a | `host/broker/exec_settings.go` | +2 / −0 |
| M3a | `host/broker/exec_settings_test.go` | +142 / −0 |
| M3a | `host/broker/handlers_exec.go` | +1 / −1 |
| M3b | `host/broker/handlers_exec_m2_test.go` | +46 / −0 |
| M3b | `host/broker/handlers_exec_srt_test.go` | +141 / −0 |

M3a totals **146 changed LOC**; M3b totals **187**. Full contents are saved in `.snap/M3a/` and cumulative `.snap/M3b/`. Both existing goldens remain byte-identical. `.scratch284/` was deleted. No Git writes occurred.

All Go gates used `WORLD_TOOL_AILANG_BIN` unset and worktree-local `TMPDIR`/`GOTMPDIR`.

| Gate command | rc / result |
|---|---|
| `go vet ./...` | **0**, both milestones |
| `go test -count=1 ./host/broker/ -run '^(TestExecSettingsProtectedRendered\|TestExecSettingsGoldenDarwin\|TestExecSettingsGoldenLinux\|TestExecSettingsNoProtectedIsByteIdentical\|TestRenderExecSettings.*\|TestWriteAndVerifyExecSettings\|TestExecProtectedSettingsThroughHandler\|TestExecHandler.*)$' -v -timeout 300s` | **1**, both runs; all new always-on tests PASS |
| `go test -count=1 ./host/broker/ -run '^TestExecSrtProtectedWrites$' -v -timeout 300s` | **0, SKIP**: srt unset |
| `gofmt -l host/broker/exec_settings.go host/broker/exec_settings_test.go host/broker/handlers_exec.go host/broker/handlers_exec_m2_test.go host/broker/handlers_exec_srt_test.go` | **0**, no output |
| `go test -count=1 ./host/broker/ -run '^(TestExecSettings.*\|TestRenderExecSettings.*\|TestWriteAndVerifyExecSettings)$' -v -timeout 300s` | **0** |
| `go test -count=1 ./host/broker/ -run '^TestExecHandlerTimeoutKillsTheGroup$' -v -timeout 300s` against untouched base files | **1** |

The focused gate fails at `handlers_exec_m2_test.go:686`: **“specific child readiness before real kill”**. The isolated base run fails at that same assertion. It was left unchanged.

Each valid mutation ran independently, followed by restoration:

| Mutant | Killer test/subtest | Observed red assertion | Restored green |
|---|---|---|---|
| MUT-NO-SRT-PROTECTED | `TestExecSettingsProtectedRendered` | `:305: protected denyWrite suffix = …` | Yes, rc 0 |
| MUT-PROTECTED-BEFORE-GITHOOKS | `TestExecSettingsGoldenLinux` | `:376: protected linux golden bytes/order differ` | Yes, rc 0 |
| MUT-LINUX-GLOB | `TestExecSettingsGoldenLinux` | `:376: protected linux golden bytes/order differ` | Yes, rc 0 |
| MUT-SRT-EMPTY-BYTES | `TestExecSettingsNoProtectedIsByteIdentical` | `:395: absent/empty protected changed legacy golden bytes` | Yes, rc 0 |
| MUT-HANDLER-PROTECTED-DROPPED | `TestExecProtectedSettingsThroughHandler/true` | `:895: handler protected denyWrite suffix missing or unexpected` | Yes, rc 0 |
| MUT-NO-SRT-PROTECTED | `TestExecSrtProtectedWrites/modify` | **UNMEASURED**: SKIP; no failure line | No; restored run also SKIP |

All five measured mutants returned rc 1 at assertions. A separate configured-corrupt-pin control returned **rc 1**, reporting `configured srt CLI pin is invalid`; it did not skip.

No scope or LOC deviations. Real Darwin/Linux confinement and the real `modify` mutation kill remain **UNMEASURED**. The full broker race suite was not run, as instructed: **UNINFORMATIVE UNDER SANDBOX**, reserved for the controller.
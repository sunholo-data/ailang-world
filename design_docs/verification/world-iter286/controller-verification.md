# Controller verification — iteration 286 (row 156 M4c)

All commands outside the codex sandbox, `WORLD_TOOL_AILANG_BIN` unset, `AILANG_BIN` = pinned v0.41.0.

| Tree | go vet ./... | named arm-8 tests | real srt (darwin) TestExecProbeArm8 + TestExecSrtStartupProbePasses | go test -race ./host/broker/... |
|---|---|---|---|---|
| base c9ab49b | rc 0 | absent (new) | — | rc 0 (194 s) |
| executor run 1 | rc 0 | rc 0 | rc 0, 2/2 PASS | rc 0 (195 s) |
| after fix 1 (= 885af2d) | rc 0 | rc 0, 21 PASS | rc 0, 2/2 PASS | rc 0 (186 s) |
| after fix 2 (= 9edf116) | rc 0 | rc 0, 25 PASS + 2 named SKIP (linux-only EBUSY/EXDEV rows) | rc 0, 2/2 PASS | rc 0 (184 s) |

Controller mutants (exec_probe.go, restored by byte copy + cmp):

| Mutant | Change | Before fix 1 | After fix 1 |
|---|---|---|---|
| MUT-ARM8-UNWIRED | `pp.probeProtected(ctx, hostGOOS())` removed from probeExecProfile | SURVIVED (all TestExecProbe*/TestExecSrt* ok, real srt) | KILLED: TestExecProbeArm8Wired/protected |
| MUT-ARM8-NO-DARWIN-CREATE | `if false && goos == "darwin" && basename != ""` | SURVIVED | KILLED: TestExecProbeArm8FailsUnsandboxed/darwin-create |

Commit reconstruction: executor snapshot `.snap/M4c/` cmp-identical to the tree; `shasum -c` of the three files OK after commit 885af2d.
Linux/bubblewrap arm 8: UNMEASURED-LOCAL; PR CI (ubuntu) runs the pin-free rows including the linux-only EBUSY/EXDEV positive rows.

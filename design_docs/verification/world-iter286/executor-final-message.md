# Executor final messages — iteration 286 (codex gpt-6.1-sol)

## Run 1 (M4c, 102,571 tok)
Implemented M4c in the three allowed files. Added `execProbeOrder`, `execProbeProtectedSample`, `probeProtected`, and `protectedProbe`: conditional arm 8, scratch Git baseline, refusal-token and host-byte checks, Darwin creation checks, and Linux detection through `checkProtected`.

Diff: **283 insertions, 7 deletions**, exceeding the approximate 245-LOC budget. No M5 work, checker edits, or Git write operations. Scratch was removed.

Every gate used:

```sh
env -u WORLD_TOOL_AILANG_BIN AILANG_BIN="$HOME/.pinned-ailang/ailang" \
  TMPDIR="$PWD/.scratch-iter286" GOTMPDIR="$PWD/.scratch-iter286"
```

| Gate command | rc | PASS/SKIP evidence |
|---|---:|---|
| `go test -count=1 ./host/broker/... -run '^TestExecProbeArm8$' -timeout 120s -v` | 0 | `--- SKIP: TestExecProbeArm8` — real srt unset |
| `go test -count=1 ./host/broker/... -run '^TestExecProbeArm8FailsUnsandboxed$' -timeout 120s -v` | 0 | `--- PASS: TestExecProbeArm8FailsUnsandboxed`; all four subtests PASS |
| `go test -count=1 ./host/broker/... -run '^TestExecProbeArm8Absent$' -timeout 120s -v` | 0 | `--- PASS: TestExecProbeArm8Absent` |
| `go test -count=1 ./host/broker/... -run '^TestExecProbeProtectedSample$' -timeout 120s -v` | 0 | `--- PASS: TestExecProbeProtectedSample`; all eleven subtests PASS |
| `go test -count=1 ./host/broker/... -run '^TestExecProbe\|^TestExecSrtStartupProbePasses$' -timeout 300s` | 1 | **UNINFORMATIVE under sandbox:** existing network test encountered `bind: operation not permitted` |
| `go vet ./...` | 0 | Clean; no output |

The passing negative-control subtests were `pass-through`, `missing-client`, `forged-token`, and `linux-detect`. Sample subtests passing were `literal`, `star`, `question`, `class`, `class-star`, `class-question`, `unicode`, `negated`, `mixed-case`, `directory`, and `nested`.

Each mutation ran independently from saved bytes, followed by restoration and the same named test passing:

| Mutant | Changed site | Killer test/subtest | Red assertion output | Restored |
|---|---|---|---|---|
| MUT-ARM8-RC-ONLY | `exec_probe.go:712` | `TestExecProbeArm8FailsUnsandboxed/missing-client` | `missing-client: arm8-protected-write must fail: token ""` | PASS, rc 0 |
| MUT-ARM8-NO-HOST-CHECK | `exec_probe.go:712` | `TestExecProbeArm8FailsUnsandboxed/forged-token` | `forged-token: arm8-protected-write must fail: host unchanged=false` | PASS, rc 0 |
| MUT-ARM8-LINUX-NO-DETECT | `exec_probe.go:706` | `TestExecProbeArm8FailsUnsandboxed/linux-detect` | `linux-detect: arm8-protected-write must fail: linux detection must report intact=false` | PASS, rc 0 |
| MUT-ARM8-ALWAYS | `exec_probe.go:596` | `TestExecProbeArm8Absent` | `legacy seven-arm order: got [... arm8-protected-write] ...; no arm8 for absent set` | PASS, rc 0 |
| MUT-ARM8-SAMPLE | `exec_probe.go:619` | `TestExecProbeProtectedSample/class` | `probe specimen must match declared pattern "spec[0-9]?.go"` | PASS, rc 0 |

All changed sites above are in `host/broker/exec_probe.go`. No compile error, skip, timeout, or socket denial was counted as a kill.

Residual validation: real srt remains unmeasured; the broader socket-dependent result and full race suite require the controller’s outside-sandbox rerun. No additional defect outside the allowlist was identified.

Full-content snapshots, verified byte-identical:

- [.snap/M4c/host/broker/exec_probe.go](.snap/M4c/host/broker/exec_probe.go)
- [.snap/M4c/host/broker/exec_probe_test.go](.snap/M4c/host/broker/exec_probe_test.go)
- [.snap/M4c/host/broker/exec_probe_srt_test.go](.snap/M4c/host/broker/exec_probe_srt_test.go)
## Fix round 1 (controller survivors, 62,422 tok)
Added startup-path coverage for arm 8 and the seven/eight-arm totals, plus `darwin-create` and the passing `darwin-refuses-all` control. Updated `.snap/M4c/host/broker/exec_probe_test.go`. Production remains byte-identical to turn start; `.scratch-iter286/` was deleted. No Git writes.

Both gates used the required environment:
- Focused tests: **rc 0**, no skips.
- `go vet ./...`: **rc 0**, silent.

```text
--- PASS: TestExecProbeArm8Wired
--- PASS: TestExecProbeArm8FailsUnsandboxed
--- PASS: TestExecProbeArm8Absent
--- PASS: TestExecProbeProtectedSample
```

The startup test uses `probeNode`’s named skip discipline when Node is absent. It ran here without real srt.

All mutants returned **rc 1** on their killer; every restored killer returned **rc 0 / PASS**. Below, `Fails/` means `TestExecProbeArm8FailsUnsandboxed/`.

| Mutant | Killer | Red assertion line | Restored green |
|---|---|---|---|
| MUT-ARM8-RC-ONLY | `Fails/missing-client` | `missing-client: arm8-protected-write must fail: token ""` | PASS, rc 0 |
| MUT-ARM8-NO-HOST-CHECK | `Fails/forged-token` | `forged-token: arm8-protected-write must fail: host unchanged=false` | PASS, rc 0 |
| MUT-ARM8-LINUX-NO-DETECT | `Fails/linux-detect` | `linux-detect: arm8-protected-write must fail: linux detection must report intact=false` | PASS, rc 0 |
| MUT-ARM8-ALWAYS | `TestExecProbeArm8Absent` | `empty protection invoked arm8` | PASS, rc 0 |
| MUT-ARM8-SAMPLE | `TestExecProbeProtectedSample/class` | `probe specimen must match declared pattern "spec[0-9]?.go": sample="world-probe.spec" error=<nil>` | PASS, rc 0 |
| MUT-ARM8-UNWIRED | `TestExecProbeArm8Wired/protected` | `startup failed arms [arm2-write-outside arm3-read-fence arm4-srt-default-write arm6-network]: arm8-protected-write present must be true` | PASS, rc 0 |
| MUT-ARM8-NO-DARWIN-CREATE | `Fails/darwin-create` | `darwin-create: arm8-protected-write must fail: protected creation` | PASS, rc 0 |
## Fix round 2 (judge r1 survivors, 61,309 tok)
Added six rows in [exec_probe_test.go](host/broker/exec_probe_test.go:324), preserving existing assertions. Saved the [snapshot](.snap/fix2/host/broker/exec_probe_test.go). Production remained byte-identical; every restoration passed `cmp` (rc 0). Scratch removed; no Git writes performed.

Both gates passed with the required environment:

- Focused `go test`: **rc 0**
- `go vet ./...`: **rc 0**

```text
PASS: TestExecProbeArm8Wired
PASS: TestExecProbeArm8FailsUnsandboxed
PASS: TestExecProbeArm8Absent
PASS: TestExecProbeProtectedSample
SKIP: TestExecProbeArm8FailsUnsandboxed/linux-ebusy-passes
SKIP: TestExecProbeArm8FailsUnsandboxed/linux-exdev-passes
```

| Mutant | rc | Result |
|---|---:|---|
| A — LINUX-NO-WROTE-TOKEN | 1 | KILLED |
| C — single drop | 0 | SURVIVED — equivalent |
| C — double drop | 1 | PANIC — not a kill |
| E — NO-ABSENT | 1 | KILLED |
| I — LINUX-ALL-TARGETS | 1 | KILLED |
| K — NO-EBUSY-EXDEV | 0 | UNMEASURED-LOCAL (CI-only; skipped) |
| RC-ONLY | 1 | KILLED |
| NO-HOST-CHECK | 1 | KILLED |
| LINUX-NO-DETECT | 1 | KILLED |
| ALWAYS | 1 | KILLED |
| SAMPLE | 1 | KILLED |
| UNWIRED | 1 | KILLED |
| NO-DARWIN-CREATE | 1 | KILLED |

C cannot meet the requested kill: every checker error returns `intact=nil`, so the single-drop variant returns the same diagnostic through its remaining nil guard. Removing both guards causes a nil dereference.

K exposes the specified platform seam: `execFileRefusal` uses the host GOOS. Both rows skip on this Darwin host with the exact reason: “execFileRefusal uses the host GOOS; EBUSY/EXDEV rows run on linux hosts (CI)”.
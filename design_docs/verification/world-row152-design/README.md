# Row 152 design probe (iteration 242)

Evidence for `design_docs/planned/w-exec-availability-decoupled-from-ailang-policy.md` V6–V8.

- `zz_probe152_test.go.txt` — the probe, a `package daemon` test. It is stored as `.go.txt` so `go build`/`go test` never compile it.
- `probe-run1.log` — run 1 (r1 of the design, base `fe7c2a4`), filtered to `PROBE`/result lines.
- `probe.log` — run 2 (r2, same base), full `go test -v` output plus `rc=`.

The probe uses only existing test seams: `fakeToolBin`'s shape, `stubExecSandbox`, `execProfile`, `execCall`, `newWSFixture`, `mustWSDaemon`. It needs no real srt and no v0.52.1 tool binary.

Re-run it from the repo root on any base where those seams exist:

```sh
export PATH=/opt/homebrew/bin:$PATH AILANG_BIN=~/.pinned-ailang/ailang   # gate interpreter v0.41.0
cp design_docs/verification/world-row152-design/zz_probe152_test.go.txt host/daemon/zz_probe152_test.go
go test ./host/daemon/ -count=1 -v -run '^TestZZProbe152' 2>&1 | grep -E 'PROBE|^(---|ok|FAIL)'
rm host/daemon/zz_probe152_test.go    # never commit the compiled copy
```

What to expect on the base:

- `mode=sleep`: each `registry("ep1")` takes about `workspaceHandlerBudget` (3 s), with `names=[] execBound=false`, and `summaries` grows 1 → 2.
- The operator line says `timed out after 10s`. That is the mislabel the design notes in §4 and R1.
- `mode=fail`: the same empty registry, in tens of ms.
- `execHandler-direct`: `exit_code 0` in the ep1 worktree (positive control).
- `LockCoupling`: `execHandler(ep2)` waits about 2.8 s behind ep1's stuck summary.

After row 152 lands, `names` must include `Workspace.Exec` in both modes, and the LockCoupling wait must be ~0.

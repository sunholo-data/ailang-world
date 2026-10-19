# Iteration 285 — controller verification (row 156 M4a + M4b)

All commands outside the codex sandbox, `WORLD_TOOL_AILANG_BIN` unset, `AILANG_BIN=~/.pinned-ailang/ailang` (v0.41.0).

| Step | Head | Result |
|---|---|---|
| Baseline `go vet ./...` | `fd7b325` | rc 0 |
| Baseline `go test -race -count=1 ./host/broker/...` | `fd7b325` | rc 0 (144 s) |
| M4a alone: vet + `TestProtectedGitFilterPreflight` (5 rows) | M4a snapshot | rc 0 |
| M4a mutants MUT-PROTECTED-FILTER-PREFLIGHT / -ERROR-OPEN (controller-run) | M4a | both FAIL on `unsafe progression: status spawns=1` (marker created under the first); restored green |
| M4b gates as delivered | M4b snapshot | `RootAndBound/cwd` and `/argv` FAIL: fake git logs `$PWD`=/private/var/…, test expected t.TempDir() /var/… (macOS symlink; executor's worktree TMPDIR hid it) → test-only fix (EvalSymlinks) → rc 0 |
| M4b controller mutants: per-call fresh context, HEAD unchecked, `dir: ""`, intact on parse error | `167ec9a` | 4/4 FAIL on their assertions; restored green |
| Snapshot reconstruction | M4b | `shasum -c` OK on all three files before the controller fix |
| `go vet ./...` | `167ec9a` | rc 0 |
| `go test -race -count=1 ./host/broker/...` | `167ec9a` | rc 0 |
| `go test -count=1 ./...` | `167ec9a` | rc 1: 5 `RootAndBound` rows timed out inside the 2 500 ms phase under load (judge NB2, reproduced) |
| Judge r1 survivor (`_ = protectedGitPreflight(...)`) | `167ec9a` | SURVIVED `^TestProtected` (reproduced) |
| Fix `d4ff957`: `filter-refused` row; non-bound rows lift the bound | `d4ff957` | survivor now FAILS `filter-refused`; per-call-context mutant still FAILS `shared-phase-bound` |
| 6 concurrent `-race` copies of `^TestProtected` | `d4ff957` | 6/6 rc 0 |
| `go vet ./...` | `d4ff957` | rc 0 |
| `go test -race -count=1 ./host/broker/... ./host/daemon/...` | `d4ff957` | rc 0 |
| `go test -count=1 ./...` | `d4ff957` | rc 1 in three packages this PR does not touch (`TestCLIRealSubprocessEpisode`, `TestWorkspaceExecEpisodeSelectsItsProfile`, `TestProducerExecTimeoutKillsTheWholeProcessGroup`), run while the judge's 8-way stress loaded the machine; each passes in isolation at `d4ff957` (1.3 s, 1.2 s, 5.3 s). PR CI is the verdict. |

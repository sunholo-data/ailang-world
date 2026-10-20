# Iteration 289 — controller verification (outside the sandbox, darwin/arm64)

Environment: `AILANG_BIN=~/.pinned-ailang/ailang` (v0.41.0), `WORLD_TOOL_AILANG_BIN` unset unless stated, default TMPDIR.

## Baseline, pristine `0ffbb9c9` (detached worktree)
- `go vet ./host/daemon/` rc 0
- `go test -count=1 ./host/daemon/ -run '^TestWorkspaceWriteProtectedPolicyTool$' -v` rc 0 SKIP (pin unset); rc 0 PASS 1.45 s (v0.54.15)
- `go test -race -count=1 ./host/daemon/` rc=0 (229.055s)

## M5e at `1b61eba9`
- G1 real pin v0.54.15, `-count=3`: rc 0, 42 subtest PASS (14 × 3), 0 SKIP, 0 FAIL
- G2 pin unset: rc 0, `--- SKIP: TestWorkspaceWriteProtectedAliasesPolicyTool` (discovered)
- `go vet ./...` rc 0; `gofmt -l host/daemon/workspace_protected_test.go` empty
- G5 existing lexical `TestWorkspaceWriteProtectedPolicyTool` on v0.54.15: rc 0 PASS
- `go test -race -count=1 ./host/daemon/` rc=0 (261.979s)

## Drills (one at a time; landed + builds asserted; restored from cp backups, sha256 verified)
| Drill | Hunk | Result |
|---|---|---|
| D3 tool regression | test-only temporary `ToolBinaryRelease = "AILANG v0.52.1"` + `WORLD_TOOL_AILANG_BIN=v0.52.1` (response tool digest `0dd70a1d…` = v0.52.1) | FAIL: the 5 symlink arms, `hardlink-write`, `hardlink-edit`, `no-protected-hardlink`; PASS: the 5 `no-protected-symlink-*` and `no-protected-singlelink` — exactly as predicted |
| D1 nil protected set | `workspace.go:642` `RenderEpisodePolicy(sandbox, nil)` | FAIL: the 5 protected symlink arms; PASS: hard-link arms (refusal is unconditional) and all controls |

Executor's in-sandbox drills (D1, D2 ×3 precondition neuters, D3) are in `executor-final-message.md`; the judge's own mutants are in `evaluator-r1.md`.

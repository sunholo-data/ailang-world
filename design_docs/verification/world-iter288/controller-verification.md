# Iteration 288 — controller verification (row 156 M5b)

All gates ran outside the sandbox with `WORLD_TOOL_AILANG_BIN` unset (except the real-tool arm) and `AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0).

## Baseline at pristine 0ee1945 (before the executor)
| Gate | Result |
|---|---|
| `go vet ./...` | rc 0 |
| `go test -race -count=1 ./host/broker/...` | rc 0 (211 s) |
| `go test -race -count=1 ./host/daemon/...` | rc 0 (279 s) |

`./cmd/...` was not in the baseline; see the TMPDIR note below.

## On the executor's tree (committed as 3f8af26)
| Gate | Result |
|---|---|
| `gofmt -l host/` | empty |
| `go vet ./...` | rc 0 |
| six named daemon tests (`TestEpisodePolicyProtectedRebased`, `TestWorkspaceProtectedPolicyWiring`, `TestWorkspaceProtectedPolicyRequests`, `TestWorkspaceWriteProtectedPolicyTool`, `TestWorkspacePolicyRenderedPerEpisode`, `TestWorkspaceModuleRootIsTheSandbox`) | rc 0, 5 PASS, 1 SKIP (real-tool arm, pin unset) |
| real-tool arm, `WORLD_TOOL_AILANG_BIN=~/.pinned-ailang-tools/v0.54.15/ailang` | rc 0, PASS: write and edit to `tests/x.ail` and `a_test.go` refused with `fs_deny_write`, host bytes unchanged, unprotected write succeeds |
| `go test -race -count=1 ./host/broker/...` | rc 0 (258 s) |
| `go test -race -count=1 ./host/daemon/...` | rc 0 (278 s) |
| `go test -count=1 ./cmd/...` with a scratch TMPDIR | rc 1: `TestDemoProofNoVerifiedSetRefuses` MOD010 (module path vs a non-temp-shaped TMPDIR) |
| `go test -count=1 ./cmd/...` with the default TMPDIR | rc 0 (worldd 60 s, world-publish 35 s) — the only variable was TMPDIR's location |

## Controller mutants (diff-anchored, one at a time, cp-restored, byte-identical restore confirmed)
| Mutant | Result |
|---|---|
| C1 single-flight join ignores the protected set | SURVIVED every executor test → controller added `TestWorkspaceProtectedPolicyJoinRequiresSameSet` (in 3f8af26): killed 3/3 (`:234` waiters=1), passes 3/3 under `-race` |
| C2 module root taken relative to the workspace root, not the episode | killed: `TestEpisodePolicyProtectedRebased` (5 subtests) and `TestWorkspaceProtectedPolicyRequests` |
| C3 `RenderEpisodePolicy(sandbox, nil)` under the real tool | killed: `TestWorkspaceWriteProtectedPolicyTool:254` (write to `tests/x.ail` returned ok:true) |
| C4 handler cache stores no protected set | killed: `TestWorkspaceProtectedPolicyWiring/rebind:191` (summaries 2 then 3) |
| C5 `AilangToolConfig.Protected` dropped | survives the real-tool arm (write/edit use the base policy, expected); killed by `TestWorkspaceProtectedPolicyRequests:227` (run variant lost the entries), confirming the executor's MUT-CONFIG-DROP claim |

## Judge r1 survivors, reproduced (rc 0 against the family + protected selector)
- M9 (join ignores the module root): accepted residual — an episode's sandbox is derived per call from startup-static config, so a same-episode root mismatch cannot happen at runtime (judge r2 agreed).
- M10 (in-flight build does not record its set): fixed test-only in `c8ab0c3` (`TestWorkspaceProtectedPolicyJoinSameSetShares`, killed 3/3 at `:270` "summaries = 2, want 1").

## Final tree c8ab0c3
| Gate | Result |
|---|---|
| `gofmt -l host/` | empty |
| `go vet ./...` | rc 0 (a first attempt read rc 1 because a zsh `rm -f` glob with no match broke the `&&` chain before vet ran — instrument failure, re-run) |
| `go test -race -count=1 ./host/daemon/...` | rc 0 (270 s) |
| both Join tests, `-race -count=3` | PASS |

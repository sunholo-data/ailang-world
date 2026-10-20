# Iteration 287 — controller verification (row 156 M5a)

All gates ran outside the sandbox with `WORLD_TOOL_AILANG_BIN` unset and `AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0).

## Baseline at pristine 0adeb70 (before the executor)
| Gate | Result |
|---|---|
| `go vet ./...` | rc 0 (2 s) |
| `go test -race -count=1 ./host/broker/...` | rc 0 (213 s) |
| `go test -race -count=1 ./host/daemon/...` | rc 0 (272 s) |

## On the executor's tree (committed as 6350cb3)
| Gate | Result |
|---|---|
| `gofmt -l host/` | empty |
| `go vet ./...` | rc 0 |
| six named broker tests (`TestEpisodePolicyProtected*`, `TestRenderEpisodePolicyHasExactlyTheDesignKeySet`, `TestRenderRunVariantKeySetsPerClass`) | rc 0, 6 PASS, 0 FAIL, 0 SKIP |
| three named daemon tests | rc 0, 3 PASS |
| `go test -race -count=1 ./host/broker/...` | rc 0 (185 s) |
| `go test -race -count=1 ./host/daemon/...` | rc 0 (230 s) |
| `go test -count=1 ./cmd/...` | rc 0 (57 s) |

## Controller mutants (diff-anchored, one at a time, cp-restored, byte-identical restore confirmed)
| Mutant | Killed by |
|---|---|
| C1 `RenderRunVariant` base short-circuit passes nil | `TestEpisodePolicyProtectedAllVariants/base` |
| C2 `tomlString` writes `\` raw | `TestEpisodePolicyProtectedGolden/escapes` |
| C3 `tomlString` writes control runes raw | `TestEpisodePolicyProtectedGolden/escapes` |
| C4 protected appended before `.git/**` and the built-in denies | Golden, AllVariants (all 4), ThroughHandler/true |
| C5 handler keeps the caller's slice (no copy) | `TestEpisodePolicyProtectedThroughHandler/true` |

## Judge r1 N1, reproduced, then fixed (3852dbc)
`case unicode.IsControl(r)` → `case r < 0x20` in `tomlString` survived every M5a test (rc 0), because the test decodes with `encoding/json`, which accepts a raw DEL. Test-only fix: `/escapes` requires `\u007F` and no raw C0/DEL rune. After the fix, `r < 0x20` fails at :751 ("DEL is not escaped"), `r < 0x7f` also fails; 6/6 named tests PASS.

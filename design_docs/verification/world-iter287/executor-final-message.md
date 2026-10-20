Implemented M5a: protected patterns reach every policy variant, cached handler variants, and `VerifyRunCaps`. TOML escaping is safe; nil/empty protection preserves legacy bytes.

Changed **8 allowed files, +314/−26 lines (340 total)**:

| File | + / − |
|---|---:|
| [handlers_ailang.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter287-row156-m5a/host/broker/handlers_ailang.go) | 10 / 4 |
| [handlers_ailang_run.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter287-row156-m5a/host/broker/handlers_ailang_run.go) | 35 / 10 |
| [handlers_ailang_run_test.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter287-row156-m5a/host/broker/handlers_ailang_run_test.go) | 261 / 4 |
| [handlers_ailang_test.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter287-row156-m5a/host/broker/handlers_ailang_test.go) | 3 / 3 |
| [workspace.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter287-row156-m5a/host/daemon/workspace.go) | 1 / 1 |
| [workspace_layout.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter287-row156-m5a/host/daemon/workspace_layout.go) | 1 / 1 |
| [workspace_layout_test.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter287-row156-m5a/host/daemon/workspace_layout_test.go) | 2 / 2 |
| [workspace_test.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter287-row156-m5a/host/daemon/workspace_test.go) | 1 / 1 |

Tests passed, with `WORLD_TOOL_AILANG_BIN` unset:

| Test | PASS | Subtests |
|---|---:|---:|
| `TestEpisodePolicyProtectedGolden` | 1 | 2 |
| `TestEpisodePolicyProtectedAllVariants` | 1 | 4 |
| `TestEpisodePolicyNoProtectedIsByteIdentical` | 1 | 4 |
| `TestEpisodePolicyProtectedThroughHandler` | 1 | 2 |
| `TestRenderEpisodePolicyHasExactlyTheDesignKeySet` | 1 | 0 |
| `TestRenderRunVariantKeySetsPerClass` | 1 | 8 |
| `TestWorkspacePolicyRenderedPerEpisode` | 1 | 0 |
| `TestWorkspaceModuleRootIsTheSandbox` | 1 | 0 |
| `TestWorkspacePackageCacheLinkedIntoEpisodeHome` | 1 | 0 |
| `TestBrokerDependencyAllowlist` | 1 | 0 |
| `TestDaemonDependencyAllowlist` | 1 | 4 |

`go vet ./...` passed. Focused race checks passed **11 tests and 24 subtests**.

Each mutant ran independently. Assertion lines below refer to `handlers_ailang_run_test.go`.

| Mutant | Exact killer test/subtest | Observed failing assertion | Restored green |
|---|---|---|---|
| MUT-NO-POLICY-APPEND | `TestEpisodePolicyProtectedGolden` | 717: differs from independent golden | Yes |
| MUT-VARIANT-MISSED (Env) | `TestEpisodePolicyProtectedAllVariants/Env` | 765: entries are not exact trailing denies | Yes |
| MUT-VARIANT-MISSED (Net) | `TestEpisodePolicyProtectedAllVariants/Net` | 765: entries are not exact trailing denies | Yes |
| MUT-POLICY-QUOTE | `TestEpisodePolicyProtectedGolden/quote` | 726: missing quote escape; 728: invalid array | Yes |
| MUT-POLICY-EMPTY-BYTES | `TestEpisodePolicyNoProtectedIsByteIdentical` | 797, 811: changed frozen legacy bytes | Yes |
| MUT-ALIAS | `TestEpisodePolicyNoProtectedIsByteIdentical/base` | 801: mutated backing array | Yes |
| MUT-HANDLER-DROP | `TestEpisodePolicyProtectedThroughHandler` | 886: Env policy lacks protected entries | Yes |
| MUT-VERIFY-DROP | `TestEpisodePolicyProtectedThroughHandler` | 899: base policy lacks protected entries | Yes |

Deviations and remaining verification:

- Full unfiltered race gate **not run**: the binding plan requires a contained runner because legacy tests write real HOME and `/tmp/claude`. Focused race results do not discharge that gate.
- Early Go runs used the default external build cache. Later runs used worktree-local `GOCACHE`; scratch was deleted.
- Used the permitted independent basic-string-array decoder because importing BurntSushi required a missing `go.sum` entry.

No dependency files, `.ail` files, or design docs changed. No Git write commands ran.
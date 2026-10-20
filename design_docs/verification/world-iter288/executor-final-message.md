Implemented M5b only. Episode policies now use canonical profile protection, rebased to the module root, with fail-closed binding reads and protected-set-aware caching.

Files changed: **329 LOC** including deletions.

- [exec_protected.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter288-row156-m5b/host/broker/exec_protected.go): +6
- [workspace.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter288-row156-m5b/host/daemon/workspace.go): +42/−16
- [workspace_protected_test.go](/Users/voightkampff/dev/sunholo-data/.wt-world-iter288-row156-m5b/host/daemon/workspace_protected_test.go): +265

| Named test | Result |
|---|---|
| `TestEpisodePolicyProtectedRebased` | PASS; 7 subtests |
| `TestWorkspaceProtectedPolicyWiring` | PASS; 6 subtests |
| `TestWorkspaceProtectedPolicyRequests` | PASS, including Declassify run variant |
| `TestWorkspaceWriteProtectedPolicyTool` | SKIP; pin unset |
| Both required existing regression tests | PASS |

Gate results, counting top-level tests:

- `go vet ./...`: PASS.
- Required daemon selector: **5 PASS, 1 SKIP**.
- Required broker selector: **19 PASS, 1 SKIP**.
- Additional focused broker race run: **19 PASS, 1 SKIP**.
- Additional daemon acceptance/concurrency race run: **12 PASS, 1 SKIP**.

Mutation assertion locations below refer to the new test file.

| Mutant | Exact killer test/subtest | Observed assertion | Restored green |
|---|---|---|---|
| MUT-NO-REBASE | `TestEpisodePolicyProtectedRebased/under-sub` | L119: policy bytes differ | Yes |
| MUT-POLICY-WRONG-EPISODE | `TestWorkspaceProtectedPolicyWiring/episode-selection` | L182: policy bytes differ | Yes |
| MUT-NO-POLICY-APPEND, config retained | `TestWorkspaceProtectedPolicyRequests` | L218: policy bytes differ | Yes |
| MUT-NO-POLICY-APPEND, neither set | `TestWorkspaceProtectedPolicyRequests` | L218: policy bytes differ | Yes |
| MUT-CONFIG-DROP | `TestWorkspaceProtectedPolicyRequests` | L227: run variant lost protected entries | Yes |
| MUT-STALE-CACHE | `TestWorkspaceProtectedPolicyWiring/rebind` | L187: policy bytes differ | Yes |
| MUT-FAIL-OPEN | `TestWorkspaceProtectedPolicyWiring/binding-error` | L158: binding error left handler bound | Yes |
| MUT-RAW-PATTERNS, reparsing canonical input | None — **UNKILLED** | No failing assertion | Yes |

Reparsing canonical input is idempotent for these fixtures; no test distinguishes it. Production feeds only `profile.Protected` directly into rebasing.

The full combined race gate remains **NOT RUN**: legacy broker tests write real HOME and `/tmp/claude` before optional-engine checks. The plan’s contained runner is required. Real-tool enforcement also remains unverified here because the optional pin was unset.

Execution deviation: I contained `TMPDIR`/`GOTMPDIR` but omitted `GOCACHE`; Go used its external default build cache. No source-scope deviations occurred. Scratch files were removed, all mutants restored, and `git diff --check` passed.
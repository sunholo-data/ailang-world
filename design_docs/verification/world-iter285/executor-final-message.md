Implemented and saved `.snap/M4a/` and cumulative `.snap/M4b/`. **Acceptance is incomplete because of two inherited preflight helper gaps**, detailed below. No repository Git writes occurred; scratch files were removed.

| Milestone | Files changed | Changed LOC |
|---|---|---:|
| M4a | New [checker](/Users/voightkampff/dev/sunholo-data/.wt-world-iter285-row156-m4/host/broker/exec_protected_check.go), new [tests](/Users/voightkampff/dev/sunholo-data/.wt-world-iter285-row156-m4/host/broker/exec_protected_check_test.go), [§4.5 annotation](/Users/voightkampff/dev/sunholo-data/.wt-world-iter285-row156-m4/design_docs/planned/w-exec-protected-test-paths.md) | 130 |
| M4b | Checker and tests only | 202 |

All three named tests passed, with every listed subtest passing:

| Test | PASS subtests |
|---|---|
| `TestProtectedGitFilterPreflight` | `configured-clean`, `configured-process`, `configured-process-include`, `no-filter`, `query-error` |
| `TestProtectedChecker` | `clean`, `modified`, `deleted`, `untracked`, `ignored`, `index-flag` (both h/S), `head-moved`, `malformed` |
| `TestProtectedCheckerRootAndBound` | `cwd`, `shared-phase-bound`, `argv`, `nonzero`, `output-failure`, `stderr-failure`, `missing-head` |

Every mutant returned rc 1 on its named assertion; every restored run returned rc 0. The two M4a mutants were also rerun against M4b.

| Mutant | Exact edit | Killer | Observed red assertion excerpt | Restored green |
|---|---|---|---|---|
| `MUT-PROTECTED-FILTER-PREFLIGHT` | `len(filters) > 0` → `< 0` | Preflight `/configured-clean` | `unsafe progression: status spawns=1 marker error=<nil>` | Yes, both milestones |
| `MUT-PROTECTED-FILTER-ERROR-OPEN` | Query-error return → `return nil` | Preflight `/query-error` | `unsafe progression: status spawns=1` | Yes, both milestones |
| `MUT-NO-INDEX-FLAG-CHECK` | Skip loop iteration `i == 2` | Checker `/index-flag` | `index-flag verdict: … changes:[] err:<nil>` | Yes |
| `MUT-HEAD-UNCHECKED` | `if head != base` → `if false` | Checker `/head-moved` | `head-moved verdict: … changes:[] err:<nil>` | Yes |
| `MUT-IGNORED-NO` | Skip loop iteration `i == 1` | Checker `/ignored` | `ignored verdict: … changes:[] err:<nil>` | Yes |
| `MUT-CHECK-PARSE-OPEN` | Parse-error condition → `if false` | Checker `/malformed` | `malformed accepted: … err:<nil>` | Yes |
| `MUT-CHECK-CWD` | Remove data command’s `dir: root` | RootAndBound `/cwd` | `cwd/env missing "cwd=…"` | Yes |
| `MUT-CHECK-DOUBLE-BUDGET` | Apply phase context only to preflight; data calls retain individual bounds | RootAndBound `/shared-phase-bound` | `shared-phase-bound accepted: … changes:[] err:<nil>` | Yes |

All measured gates used `WORLD_TOOL_AILANG_BIN` unset, the readable pinned `AILANG_BIN`, and worktree-local temporary directories. Results are advisory.

| Gate | rc |
|---|---:|
| `go vet ./host/broker/` | 0 |
| `go test -count=1 ./host/broker/... -run '^TestProtectedGitFilterPreflight$' -timeout 120s -v` | 0 |
| `go test -count=1 ./host/broker/... -run '^TestProtectedChecker$' -timeout 120s -v` | 0 |
| `go test -count=1 ./host/broker/... -run '^TestProtectedCheckerRootAndBound$' -timeout 120s -v` | 0 |
| `go test -race -count=1 ./host/broker/ -run '^TestProtected' -timeout 300s` | 0 |
| `gofmt -l host/broker/` | 0; empty output |
| `git diff --check` | 0 |
| Full `go test -race ./host/broker/...` | UNMEASURED: PF6 requires a contained runner |

The required read-only reuse of `configuredGitFilters` leaves two deviations:

- Its preflight runner retains a **96 MiB stdout cap**, rather than 8 MiB. All three data calls enforce 8 MiB.
- It discards failed-query stdout before treating rc 1 as “no filters.” A separate rc-1-with-output probe **failed**, showing `status spawns=1`. That probe was restored; the required rc-2 query-error case passes.

Fixing these gaps requires changing the forbidden helper file or relaxing mandatory helper reuse. No unrelated files or tests were changed to conceal them.
Implemented M2a and M2b. Snapshots are complete at `.snap/M2a/` and cumulative `.snap/M2b/`; final files match them byte-for-byte.

| Milestone | File | Changed LOC |
|---|---|---:|
| M2a | `host/broker/exec_profile.go` | +30/−1 |
| M2a | `host/broker/exec_profile_test.go` | +89 |
| M2b | `host/daemon/exec_startup_test.go` | +48/−3 |
| M2b | `docs/QUICKSTART.md` | +4 |

Totals: **M2a 120 LOC; M2b 55 LOC**.

`null` is explicitly refused, including null array elements. The loader reuses `ParseProtectedPatterns` and stores both the canonical list and `ProtectedDigest`. The v1 schema and raw-input profile digest remain unchanged. Startup coverage now has an exact 51-row census.

| Gate | Result |
|---|---|
| Focused broker tests | PASS |
| Startup refusal table | PASS; all four new rows executed |
| Broker race suite | **UNINFORMATIVE UNDER SANDBOX**, plus fixture findings |
| Daemon race suite | PASS |
| `go vet ./...` | PASS |
| `gofmt -l host/` | FINDING: unchanged `host/store/deadline_guard_test.go` |

The broker suite encountered denied `$HOME` writes and `ps` access. Existing credential fixtures also reject temporary paths inside this worktree. These were recorded without workarounds. The existing real-srt-only startup row skipped.

All four compiling mutations were killed and restored green. Exact commands and unabridged failure lines are in the [evidence report](/Users/voightkampff/dev/sunholo-data/.wt-world-iter283-row156-m2/.snap/MISSION_REPORT.md).

| Mutant | Exact edit | Test | Red assertion excerpt | Restored green |
|---|---|---|---|---|
| MUT-PROFILE-NO-CHECK | `ParseProtectedPatterns(patterns)` → `patterns, nil` | `TestExecProfileProtectedUnknownFormRefused` | `protected[0] names refused src/*.ts: err = <nil>` | Yes |
| MUT-PROFILE-TYPE | Non-array guard’s error return → `return p, nil` | `TestExecProfileProtected/wrong-type` | `wrong field type refused: err = <nil>, want an *ExecProfileError naming protected` | Yes |
| MUT-PROFILE-COMPAT | Version literal v1 → v2 | `TestExecProfileProtected/compatibility` | `v1 compatibility: version = "world/exec-profile/v2"` | Yes |
| MUT-STARTUP-PROTECTED-IGNORED | `ParseProtectedPatterns(patterns)` → `patterns, nil` | `TestExecStartupRefusalTable/protected-slash-glob` | `does not name "protected[0]" (QUICKSTART §10 row protected-slash-glob)` | Yes |

Process caveats: `AILANG_BIN` was absent, so I downloaded and hash-verified the official pinned interpreter inside `.snap/`. Initial focused runs used Go’s default cache before I redirected it into the worktree; required legacy probes also exercised hard-coded `/tmp/claude` paths. No git write commands were issued against this worktree.
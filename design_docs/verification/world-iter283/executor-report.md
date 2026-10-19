# Iteration 283 row 156: M2a and M2b executor evidence

Only the four authorized source/doc files changed. No git write commands were issued against this worktree; existing tests manage their disposable fixture repositories.

## Changes and snapshots

| Milestone | File | Added | Deleted | Changed LOC |
|---|---|---:|---:|---:|
| M2a | host/broker/exec_profile.go | 30 | 1 | 31 |
| M2a | host/broker/exec_profile_test.go | 89 | 0 | 89 |
| M2b | host/daemon/exec_startup_test.go | 48 | 3 | 51 |
| M2b | docs/QUICKSTART.md | 4 | 0 | 4 |

M2a: 120 changed LOC. M2b: 55 changed LOC. `.snap/M2a/` contains the two M2a files; `.snap/M2b/` contains all four. Each source was compared byte-for-byte against its snapshot after mutation restoration.

## Design decisions

- Strict JSON array of strings; null, other non-arrays, and non-string elements (including null) are refused.
- Reuse ParseProtectedPatterns for grammar, NFC, fold-key deduplication and sorting. Its indexed pattern refusal is wrapped as an ExecProfileError.
- Store both Protected and ProtectedDigest; use existing protectedDigest, with the empty set represented by no digest.
- Keep world/exec-profile/v1 and SHA-256 of exact profile input bytes.
- Startup rows use the existing fake tool and srt shim; assert StageConfig, field/pattern names, no daemon, no sandbox archive/probe marker/scratch, no policy summary or tool dispatch, and no decoys left behind. The node/tool version checks remain the existing startup preflight.
- Each t.Run owns its context and cleanup through start -> newWSDaemon -> boundedTestContext (context.WithTimeout, 30 seconds; t.Cleanup(cancel)); no shared deadline.
- Exact census is now 51 entries with bidirectional Go/docs identity checks. QUICKSTART edits contain only the four refusal rows.

## Tests first

M2a tests first failed to compile specifically because Protected and ProtectedDigest did not exist (logs/M2a-tests-first.log); then both named tests ran and passed (logs/M2a-green.log).

M2b tests first failed on `exact refusal census: QUICKSTART=47 Go=51, want 51 each` and the four missing documented rows (logs/M2b-tests-first.log). Adding the four matching refusal rows made the table pass (logs/M2b-green.log).

## Required mutation proof

One compiling production mutation at a time; every mutation was restored from its milestone snapshot without git and the same named command passed again. No mutation run was a panic, timeout, compile failure, skip or zero-test result.

| Mutant | Exact edit | Test command | Restored green |
|---|---|---|---|
| MUT-PROFILE-NO-CHECK | `p.Protected, err = ParseProtectedPatterns(patterns) => p.Protected, err = patterns, nil` | `go test -count=1 ./host/broker/... -run ^TestExecProfileProtectedUnknownFormRefused$ -timeout 120s -v` | yes |
| MUT-PROFILE-TYPE | `return nil, profileErr("protected", "must be an array of strings") => return p, nil` | `go test -count=1 ./host/broker/... -run ^TestExecProfileProtected/wrong-type$ -timeout 120s -v` | yes |
| MUT-PROFILE-COMPAT | `const ExecProfileVersion = "world/exec-profile/v1" => const ExecProfileVersion = "world/exec-profile/v2"` | `go test -count=1 ./host/broker/... -run ^TestExecProfileProtected/compatibility$ -timeout 120s -v` | yes |
| MUT-STARTUP-PROTECTED-IGNORED | `p.Protected, err = ParseProtectedPatterns(patterns) => p.Protected, err = patterns, nil` | `go test -count=1 ./host/daemon/... -run ^TestExecStartupRefusalTable/protected-slash-glob$ -timeout 120s -v` | yes |

Exact assertion failure lines (full, unabridged):

### MUT-PROFILE-NO-CHECK

```text
exec_profile_test.go:295: protected[0] names refused src/*.ts: err = <nil>
```

### MUT-PROFILE-TYPE

```text
exec_profile_test.go:267: wrong field type refused: err = <nil>, want an *ExecProfileError naming protected
exec_profile_test.go:267: wrong field type refused: err = <nil>, want an *ExecProfileError naming protected
exec_profile_test.go:267: wrong field type refused: err = <nil>, want an *ExecProfileError naming protected
exec_profile_test.go:267: wrong field type refused: err = <nil>, want an *ExecProfileError naming protected
```

### MUT-PROFILE-COMPAT

```text
exec_profile_test.go:277: v1 compatibility: version = "world/exec-profile/v2"
```

### MUT-STARTUP-PROTECTED-IGNORED

```text
exec_startup_test.go:482: refusal "daemon startup failed at config: the exec configuration is refused: exec startup probe failed for profile \"proj\" (4 of 7 arms): arm2-write-outside: a write to a sibling episode: token \"WORLD-PROBE-WROTE\" (want WORLD-PROBE-REFUSED EPERM|EACCES|EROFS), on the host: present=true (exit 0; stderr \"\"); arm3-read-fence: a read of decoy /Users/voightkampff/dev/sunholo-data/.wt-world-iter283-row156-m2/.snap/tmp/TestExecStartupRefusalTableprotected-slash-glob4287823883/001/home/.ailang-worldd-exec-probe-decoy: token \"WORLD-PROBE-READ\", want WORLD-PROBE-REFUSED EPERM|EACCES|EROFS (exit 0; stderr \"\"); arm4-srt-default-write: a write to /tmp/claude/world-exec-probe-wep-fee35d81e92f9c2c: token \"WORLD-PROBE-WROTE\", on the host: present=true (exit 0; stderr \"\"); arm6-network: a raw connect to World's live listener: token \"WORLD-PROBE-CONNECTED\", accepts 2 (want a refusal and still 1; exit 0; stderr \"\"); a connect through srt's proxy environment: token \"WORLD-PROBE-CONNECTED\", accepts 3 (want a refusal and still 1; exit 0; stderr \"\"); reverse leg: the host reached the listener the child bound on 127.0.0.1:50255" does not name "protected[0]" (QUICKSTART §10 row protected-slash-glob)
```

## Execution environment and limitations

WORLD_TOOL_AILANG_BIN was unset for every test. AILANG_BIN was absent in the supplied shell, and the PATH binary was a dirty v0.54.15 build. Downloaded the official v0.41.0 interpreter under .snap/bin and verified both hashes against pins/ailang-releases.json:

- tarball: b08f3cdea137222ea8df71c53d1a1353ef7c942fd2a307de5d7a17b8fb598e0b
- binary: 1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f

TMPDIR points to .snap/tmp for test/build temporary files and t.TempDir fixtures. Mutation and final gate runs use GOCACHE=.snap/go-cache. Required legacy probe tests still contain hard-coded /tmp/claude control writes; those existing fixtures were not changed. No new fixture depends on /tmp. Existing probe-package-only skips without WORLD_EXEC_SRT_NODE_MODULES; all four new protected rows executed without skips and required no real srt.

`gofmt -l host/` lists host/store/deadline_guard_test.go, which is byte-identical to HEAD; recorded the HEAD formatting diff in logs/gofmt-baseline.diff and left it unchanged as an out-of-scope FINDING.

## Full-suite findings so far

The broker race run exited 1 and is **UNINFORMATIVE UNDER SANDBOX**, not a regression verdict. Full raw output: logs/gate-broker-race.log.

- TestExecSrtConfinementMatrix: `handlers_exec_srt_test.go:271: open /Users/voightkampff/.world-m2-exec-decoy-…: operation not permitted`.
- TestOrdinaryOverflowCarriesNoESRCH and TestExecHandlerTimeoutKillsTheGroup depend on ps-based process observations. Direct diagnostic ps is denied (`operation not permitted`); no process-state conclusion can be drawn here.
- Existing registry credential fixtures refuse paths under the working tree. Keeping TMPDIR inside this worktree to honor the directive exposes that location conflict (e.g. `publish_op_test.go:158: broker: registry credential file "…/registry.key": file is inside the working tree …`). No fixture, credential rule or test was changed to work around it.

## Final gate results

| Command | Result | Elapsed seconds |
|---|---|---:|
| `go test -count=1 ./host/broker/... -run ^TestExecProfileProtected -timeout 120s -v` | PASS | 1.35 |
| `go test -count=1 ./host/daemon/... -run ^TestExecStartupRefusalTable$ -timeout 120s -v` | PASS | 48.67 |
| `go test -race -count=1 ./host/broker/...` | UNINFORMATIVE UNDER SANDBOX; legacy fixture FINDINGS | 239.85 |
| `go test -race -count=1 ./host/daemon/...` | PASS | 295.66 |
| `go vet ./...` | PASS | 15.37 |
| `gofmt -l host/` | FINDING: unchanged host/store/deadline_guard_test.go listed | 1.29 |

All gates completed and all result files were read before finishing. No gate is pending. The daemon race suite exited 0; no race diagnostic was printed. The focused daemon table passed with its existing real-srt-only row skipped, and all four new protected rows passed.

Plan/process deviations and limitations: AILANG_BIN was not supplied, so the verified official pin was downloaded locally; required legacy probes exercise their fixed /tmp/claude control paths, and initial focused runs used Go's default build cache before mutation/final gates switched to the worktree cache. Full broker green and whole-host formatting green could not be established here for the recorded reasons; these need the controller's outside-sandbox rerun. No out-of-scope fixes, git staging, commit, branch, or PR operations were performed.

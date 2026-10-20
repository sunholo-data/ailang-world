# Evaluator report: row 156 M5b, commit 3f8af268 (base 0ee1945). Judge: Sonnet 5.5

## Verdict: PASS, 88/100. No blocking findings.

| Line | Max | Score | Note |
|---|---|---|---|
| Correctness | 30 | 27 | The wiring is correct. Protected-set-aware cache and single-flight with a bounded retry. Two judgment calls are listed below. |
| Test quality / non-vacuity | 25 | 20 | Every production hunk has a killer except two survivors (S1, S2). The real-tool arm covers only write/edit, not the CLI/FS ops the plan names. |
| Plan / design adherence | 20 | 17 | Tests are in package daemon, not host/broker/exec_protected_tools_test.go (justified). Pin is v0.54.15, not v0.52.1. The plan's broker gate selectors for Tool/Requests would run zero tests in host/broker. |
| Code quality | 15 | 14 | Small and clear. Wrapper is minimal. |
| Evidence integrity | 10 | 10 | The executor's claims reproduce. |

## Gates (WORLD_TOOL_AILANG_BIN unset unless noted, AILANG_BIN=~/.pinned-ailang/ailang)
- gofmt -l host/: empty.
- go vet ./...: rc 0.
- Focused daemon selector, 7 named tests: 6 PASS and 1 SKIP (TestWorkspaceWriteProtectedPolicyTool, pin unset). rc 0.
- Real-tool arm with WORLD_TOOL_AILANG_BIN=~/.pinned-ailang-tools/v0.54.15/ailang: PASS (1.46s), not SKIP.
- `go test -race -count=1 ./host/daemon/...`: ok 275.6s, rc 0.
- `go test -race -count=1 ./host/broker/...`: ok 240.8s, rc 0.
- Both -race runs were done in a byte-identical copy of the worktree, with scratch TMPDIR/GOTMPDIR, so I could mutate in parallel.
- `go test ./cmd/...` with the default TMPDIR: not run. The diff does not touch cmd/.
- Diff scope: exactly the three files in the brief. Nothing else changed.

## Mutants (mine, one at a time; a kill is a named test RED on an assertion)
| Mutant | Result / killer |
|---|---|
| M1 wrapper returns canonical (no rebase) | killed: TestEpisodePolicyProtectedRebased/under-sub, /ancestor, /ancestor-above, /disjoint, /no-cross-episode-leak; TestWorkspaceProtectedPolicyRequests |
| M2 wrapper rebases against "." | killed: the same set as M1 |
| M3 ignore binding read error | killed: Wiring/binding-error |
| M4 profileFor with empty binding | killed: Wiring/rebind; JoinRequiresSameSet |
| M5 use the default profile (e.profile) | killed: Rebased/no-cross-episode-leak; Wiring/episode-selection; Wiring/rebind; JoinRequiresSameSet |
| M6 rel forced to "." | killed: Rebased/under-sub, /ancestor, /ancestor-above, /disjoint, /no-cross-episode-leak; Requests |
| M7 cache ignores the set | killed: Wiring/rebind; JoinRequiresSameSet |
| M8 join ignores the set | killed: JoinRequiresSameSet |
| M9 join ignores the root (`b.root != sandbox` removed) | SURVIVED (S1) |
| M10 build.protected not recorded | SURVIVED (S2) |
| M11 episodeTool.protected not stored | killed: Wiring/rebind |
| M12 RenderEpisodePolicy(sandbox, nil) | killed: Rebased (6 subtests), Wiring/episode-selection, Wiring/rebind, JoinRequiresSameSet, Requests |
| M13 drop AilangToolConfig.Protected | killed: Requests (run-variant assertion) |
| M15 pass nil to buildAilangHandler | killed: Rebased, Wiring/episode-selection and /rebind, JoinRequiresSameSet, Requests |
| M16 skip the bindingFor read | killed: Wiring/rebind, Wiring/binding-error, JoinRequiresSameSet |
| M17 ignore the filepath.Rel error | SURVIVED, equivalent (Rel cannot fail for a sandbox under the root; error path unreachable) |
| M14 retry returns b.tool instead of recursing | build error, so not a valid mutant. Not re-run. |

The M18 mutant (removing slices.Clone) was a no-change mutant and was not run. Clone only guards against aliasing, so it is equivalent here.

Executor mutants MUT-NO-REBASE, MUT-POLICY-WRONG-EPISODE, MUT-NO-POLICY-APPEND, MUT-CONFIG-DROP, MUT-STALE-CACHE and MUT-FAIL-OPEN are all covered by M1/M6, M5, M12, M13, M7/M11 and M3/M16 above. I did not run MUT-RAW-PATTERNS: reparsing canonical patterns is idempotent on these fixtures, so I accept it as unkillable.

## Non-blocking findings
- S1: The root-mismatch branch of the join (an in-flight build at another module root) has no test. M9 survived the Workspace/Episode/Handler/Family suites. The branch is older than M5b but was edited by it, so a test is cheap.
- S2: No test shows that two concurrent callers with the SAME non-empty protected set still share one build. M10 only costs extra builds, not correctness, but dedupe under a non-empty set is untested.
- The real-tool arm (TestWorkspaceWriteProtectedPolicyTool) covers write and edit of two lexical patterns. The plan also names CLI and FS operations and says symlink aliases are a documented upstream residual. The test comment and report should say write/edit only.
- The plan's broker gate commands `./host/broker/... -run '^TestWorkspace...'` would run zero tests, because the tests live in host/daemon. A zero-test run is not a pass. The reproduced gates here use the daemon package. Update the plan's gate text.
- Bound episode with a profileFor refusal (digest mismatch or profile not loaded) gets no protection on its AILANG tools, because the policy is today's bytes. Exec is refused there, so no exec bypass exists. But a stale-digest episode keeps unprotected AILANG write/edit. Fail-open versus fail-closed was not stated by the design, so this is a policy call to record, not a defect.
- The test-only seam `d.workspace.exec = &workspaceExec{...}` is assigned directly after construction. It is fine under -race because it is set before use.

## Judgments on the controller-sanctioned choices
1. Tests in package daemon: agree. The fixtures (fakeToolBin, realToolBin, the registry) are in package daemon. A broker-package copy would duplicate the fixtures or require exports. The allowlist file host/broker/exec_protected_tools_test.go is simply not created. The plan text needs a one-line amendment.
2. Exec not configured gives nil protected, so today's bytes. Agree; it is exactly §4.3 item 1. TestWorkspaceProtectedPolicyWiring/no-exec, /no-protected and /no-profile pin the byte identity.
3. Binding read error gives no AILANG handlers for that episode (fail closed). Agree, and acceptable. It mirrors execHandler, which already treats the same error as unavailable or refused. A store error that fails every episode request was previously absorbed by the AILANG tools alone; now it is consistent. This is a real behaviour change only when exec is configured AND the binding store errors or has conflicting credentials. Both are already broken states for exec. Logged on the operator line, and tested for both error kinds.
4. profileFor refusal gives nil protected. Defensible (see the non-blocking finding about stale digests). It matches the "no profile renders today's bytes" rule.
5. Empty after rebase gives nil. Agree. Disjoint patterns render the same bytes as no protection, and the cache key compares slices.Equal, where nil and empty are both treated as empty.
6. Pin v0.54.15 instead of v0.52.1. Agree: the real-tool arm passes at v0.54.15, which is the pinned role tool in host/pins/ailang-releases.json.

## Concurrency
- The binding and profile read and the rebase happen before w.mu, on read-only data. -race is green over the full suites.
- Retry after `<-b.done` is bounded: it recurses once per mismatched in-flight build. Each recursion waits for a build to finish, and progress is guaranteed because the builder always closes done and removes its map entry. It is not an unbounded loop unless the binding flips continuously; the stack depth stays small in that case.
- No deadlock: w.mu is released before every wait.
- No data race. The ailangBuild.protected field is written before the build is published under w.mu.

## Cleanup
Mutated files were restored from backups. `git status --porcelain` in the eval worktree is empty (checked below), and the scratch dir is deleted.

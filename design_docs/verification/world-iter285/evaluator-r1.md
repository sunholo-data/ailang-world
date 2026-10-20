# Evaluator report, iter-285, row 156 M4a + M4b (commits c0258ca6, 167ec9af)

Verdict: FAIL on one BLOCKING finding; numeric score 78/100 (would pass >= 70 on score alone). The code is correct and every plan-named mutant is killed. The block is one test gap on the central no-filter guard (finding 1). It is a small fix: add a checker-level refusal subtest. Then PASS is expected.

Score: correctness vs spec 36/40, test strength/mutation 15/25, gates green 9/15 (flake, finding 2), scope 10/10, docs/residual honesty 8/10.

## Gates (pinned binary, WORLD_TOOL_AILANG_BIN unset)
- `go vet ./host/broker/`: ok.
- `TestProtectedGitFilterPreflight`: 5/5 PASS at 167ec9af and at c0258ca6 alone (M4a checked in my own detached worktree).
- `TestProtectedChecker`: 8/8 PASS. `TestProtectedCheckerRootAndBound`: 7/7 PASS.
- `go test -race -count=1 -run '^TestProtected'`: ok (16 s).
- Allowlist respected. M4a touches the two Go files and the design doc §4.5. M4b touches only the two Go files. No production caller and no result field were added.

## Spec conformance
Conforming:
- Preflight runs first. It reuses `configuredGitFilters` and `WorktreeFilterRefused`.
- Exactly three data legs run, with the plan's argv. The ignored leg is `strings.Fields(protectedIgnoredCommand)`.
- Every leg gets `-c core.fsmonitor=false` and the candidate pathspecs.
- `dir: root`, `env: gitEnv()`, and 8 MiB caps on both stdout and stderr.
- One `context.WithTimeout(2500 ms)` covers the preflight and all three legs, because the preflight's own 30 s timeout is nested under the shared ctx.
- `parseProtectedCandidates` is reused; there is no duplicate parser.
- Intact semantics hold: HEAD==base with zero changes is intact; `head != base` appends `HEAD`/`head-moved` even with empty paths; any error or empty head gives `intact` nil.

## Mutants
| Mutant | Edit | Killer | Result |
|---|---|---|---|
| MUT-PROTECTED-FILTER-PREFLIGHT | `len(filters)>0` -> `false` | `TestProtectedGitFilterPreflight/configured-{clean,process,process-include}` | KILLED: "unsafe progression: status spawns=1 marker error=<nil>" |
| MUT-PROTECTED-FILTER-ERROR-OPEN | query error -> `return nil` | `.../query-error` | KILLED (spawns=1) |
| MUT-NO-INDEX-FLAG-CHECK | `streams[2]=nil` before parse | `TestProtectedChecker/index-flag` | KILLED: "index-flag verdict ... changes:[]" |
| MUT-HEAD-UNCHECKED | `if head != base` -> `if false` | `TestProtectedChecker/head-moved` | KILLED |
| MUT-IGNORED-NO | `streams[1]=nil` before parse | `TestProtectedChecker/ignored` | KILLED |
| MUT-CHECK-PARSE-OPEN | `err != nil \|\| head==""` -> `false` | `TestProtectedChecker/malformed` (also `missing-head`) | KILLED |
| MUT-CHECK-CWD | drop `dir: root` | `TestProtectedChecker/*` (wrong repo, HEAD mismatch) and `RootAndBound/cwd`,`/argv` | KILLED |
| MUT-CHECK-CWD (env) | `env: nil` | `RootAndBound/cwd`,`/argv` | KILLED (env missing PATH=...) |
| MUT-CHECK-DOUBLE-BUDGET | `runBounded(context.Background(),...)` for the legs | `RootAndBound/shared-phase-bound` | KILLED ("accepted") |
| own: drop pathspecs, status leg | | `RootAndBound/argv`, `/cwd` | KILLED |
| own: drop pathspecs, ignored leg | | `RootAndBound/argv`, `/cwd` | KILLED |
| own: drop pathspecs, flags leg | | `RootAndBound/argv`, `/cwd` | KILLED |
| own: ignored snapshot -> nil | | `TestProtectedChecker/ignored` | KILLED |
| own: swap streams[1]/[2] into parser | | `TestProtectedChecker/*` | KILLED |
| own: remove stdout 8 MiB cap | | `RootAndBound/output-failure` | KILLED |
| own: remove stderr 8 MiB cap | | `RootAndBound/stderr-failure` | KILLED |
| own: fsmonitor=true on leg 3 | | `RootAndBound/argv`, `/cwd` | KILLED |
| own: head=="" treated as intact | | `RootAndBound/missing-head` | KILLED |
| own: `intact := true` | | `TestProtectedChecker/modified` etc. | KILLED |
| own: phase 2500 ms -> 25 s | | `RootAndBound/shared-phase-bound` | KILLED |
| own: runBounded error -> continue | | `nonzero`, `output-failure`, `shared-phase-bound` | KILLED |
| own: preflight moved after the legs | | `RootAndBound/argv`, `/cwd` | KILLED |
| own: head-moved only when paths present | | `TestProtectedChecker/head-moved` | KILLED |
| own: `intact` set on the error path | | `malformed`, `missing-head` | KILLED |
| **own: checkProtected ignores the preflight result** (`_ = protectedGitPreflight(...)`) | exec_protected_check.go:20-22 | none | **SURVIVED**: all 3 tests PASS |

Note: my first NO-INDEX-FLAG and IGNORED-NO attempts inserted the nil after the parse call, so they were no-ops. I reran them correctly (rows above).

## BLOCKING
1. **The preflight verdict inside `checkProtected` is not guarded** (exec_protected_check.go:20-22). The preflight tests call `protectedGitPreflight` directly. The checker tests never configure a filter. All fake-git rows answer the config query with rc 1, so they exercise only the no-filter path. A checker that runs the preflight but ignores its result, or swallows the error, passes every test. That runs the data legs and executes the host clean filter, which is the PF8 "no host filter execution" requirement.
   - Fix: add a `TestProtectedChecker/filter-refused` subtest. It uses a real repo with `filter.probe.clean = touch marker; cat` and a `.gitattributes` match. It asserts `checkProtected` returns an error containing `WorktreeFilterRefused`, `intact` is nil, the marker is absent, and a wrapping fake git logs no ` status `.
   - Reproduce: replace lines 20-22 with `_ = protectedGitPreflight(ctx, git, root)`, then run `go test -count=1 ./host/broker/ -run '^TestProtected'`. It stays green.

## NON-BLOCKING
2. **Flake risk (exec_protected_check_test.go:160+ and :200+).** Several subtests run the production 2500 ms const over 4 real or fake spawns. `RootAndBound/cwd` and `/argv` also write the full env and fork `/usr/bin/env` each call. Under load these hit the phase deadline.
   - Evidence 1: my first, accidentally no-op mutant run, which was the original code, failed 5 RootAndBound subtests with phase timeouts while other builds ran.
   - Evidence 2: 6 concurrent `-race` test binaries (`-test.count=2`) gave 1 of 6 failures, `RootAndBound/cwd (2.50s)`.
   - Mitigation: make the timeout a package var (const production value) that non-bound tests raise to about 30 s, or give the bound row the 2500 ms and everything else a larger phase. The shared-phase-bound row itself is robust: a slower runner only times out sooner, and the sleep 0.9 x 4 vs 2.5 s margin is sound.
3. **Residual (a), preflight cap 96 MiB instead of 8 MiB.** Acceptable for these milestones. The query is `git config --get-regexp` bounded by the shared 2.5 s phase, and the plan made worktree_delta.go read-only. Record it in the V24 resize or the follow-up row.
4. **Residual (b), rc 1 with stdout reads as "no filters".** Acceptable. `git config --get-regexp` never prints on rc 1. `runWorktreeGit` returns nil on error, so the guard in `configuredGitFilters` is dead code. Row-184 code, so upstream or later work.
5. **Unbounded diagnostic size (exec_protected_check.go:49).** `%q` of up to 8 MiB of stderr goes into the error string. Cap the quoted stderr to about 4 KiB, as `runWorktreeGit` does.
6. **M4a's "no status spawn" is test scaffolding.** The test itself runs status if the preflight returns nil. That is fine given there is no caller, but finding 1 shows the same wiring must be asserted at the checker.
7. **Docs.** The §4.5 annotation in the design doc is accurate and mentions the coupling.

# Product evaluation — queue row 134 (software-engineering domain), draft PR #180

Candidate: branch `attended/row134`, head `bb6b2b20618e54057a9869c32beaa9e7bb1e5493`, 9 commits over `origin/dev`
(62 files, +10518/-122). Evaluated 2026-10-03 in an own detached worktree (removed afterwards).
Design: `design_docs/planned/w-software-engineering-domain.md` (RATIFIED, D-WORLD-53). Rubric: sprint-evaluator
`scoring_rubric.md` applied literally, with the repo's `verify_*` gates standing in for `make test` / `make lint`.

## Verdict: PASS, 81/100, no hard fail, no blocking finding

Implementation fidelity to the ratified design is high. All mandated gates are green and I reproduced them first-hand.
14/14 of my own production mutants were killed. I found one real confinement gap (`ailang-cli` `fmt --write`, section 5)
that is not blocking against the ratified text but that I recommend Mark rule on before M6 mints grants.

## 1. Rubric table

| Category | Max | Score | Basis |
|---|---|---|---|
| Tests pass (HARD FAIL) | 20 | 20 | `verify_ail.sh` rc=0, `verify_go.sh` rc=0 (section 3) |
| Lint clean | 10 | 10 | `gofmt -l host cmd` empty; `go vet ./...` rc=0; `verify_ail.sh` ai-check clean (PUB011 ratchet held) |
| Acceptance criteria (HARD FAIL <50%) | 30 | 27 | about 90% of M1-M5, M7 ACs have a named, passing test; two partials (section 4); M6 is attended and not claimed |
| Code quality | 15 | 5 | mechanical deductions: one test file over 800 lines (`host/broker/handlers_ailang_test.go`, 940) = -5; functions over 50 lines (`Coordinator.Replay` about 130, `dispatchEffectful` about 80, `composeEffectOutput`, `composeRunResult`, `AilangToolHandler.Execute`, `executeRun`) = -5 (cap); zero TODO/FIXME/HACK in new code |
| Documentation | 15 | 10 | QUICKSTART section 9 added and mechanically bound (`TestQuickstartSection9FlagsMatchTheCLI`, `TestSeToolsQuickstartPayloadsVerbatim`); design doc folded V62-V66 and R-SE-12. Deducted: the design header still reads "PLANNED ... Not implemented" (no status/implementation-state note), no per-AC evidence table. No CHANGELOG or examples convention in this repo |
| Design fidelity | 10 | 9 | see section 4. Deviations are all measured and recorded (V62-V66) |
| Regression-surface / perf | n/a | n/a | not triggered (no AILANG compiler paths; no perf goal). The Conflict Surface (design section 10) is filled in |
| **Total** | 100 | **81** | pass >= 70 |

Hard-fail check: Tests pass = yes. ACs met >= 50% = yes. Regression-surface not triggered. Perf not triggered.
No hard fail.

## 2. CI state (`gh pr checks 180`, read-only)

- `ailang-code verify gate`: pass (52 s)
- `go host build + test gate`: pass (13 m 21 s)

## 3. Gates run by me (exact mandatory profile)

Env: `PATH=$HOME/.pinned-ailang:$PATH AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_PKG_AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.51.0/ailang`

| Gate | rc | Time | Notes |
|---|---|---|---|
| `./scripts/verify_ail.sh` | 0 | 16 s | "16 required identities verified, 40 named tests pass, 312 se-tools named tests pass"; all 8 `packages/se-tools/se_tools/*.ail` ai-checked |
| `./scripts/verify_go.sh` | 0 | 486 s | exactly 2 `WARNING: DATA RACE` lines (its positive control) |
| `gofmt -l host cmd` | clean | n/a | no output |
| `go vet ./...` | 0 | n/a | no output |
| Explicit real-binary e2e: `go test -count=1 -run 'TestSeTools|TestWorkspace|TestAilang' ./host/daemon/ ./host/broker/ -v` | ok / ok | 34.4 s / 9.5 s | 26 tests, all `--- PASS`, zero SKIP (RAN, not skipped): TestSeToolsE2EBinaryVersions, MCPEndToEnd, A2AEndToEnd, A2APostEffectFailureIsMapped, MCPVisibilityFollowsGrants, QuickstartPayloadsVerbatim, ExamplesSearchWithoutCorpusRefuses; the full TestWorkspace* set; TestAilangToolConfinementMatrix, RunOutcomeShapesOnTheToolBinary, RunArgvPutsDashDash..., etc. |

Independent baseline for the mutation run: `go test ./host/broker ./host/coordinator ./host/authority ./host/store` all ok before any mutant.

## 4. Per-milestone and AC assessment

Reading the code against the design:

- **M1 plan law**: `host/coordinator/effectplan.go` implements L1-L4, L6, L7 and L-ARGS (L-ARGS lives in the `.ail` plans via `argKeyAllowed`, held by `TestSeToolsPlansAdmittedByPlanLaw` and `TestSeToolsPathPredicateRefusesBeforeAnyEffect`). Predicates mirror `design_docs/sketches/effectplan.ail`; `TestEffectPlanSketchDrift` evaluates every sketch row through the Go mirrors. AC1.1/AC1.2 met. Mutant F (finish keys) killed by the drift test and `TestParseFinishLaws`.
- **M2a/M2b effectful dispatch**: `effectful.go` has per-phase caps (2 s plan, 2 s finish, handler `min(10 s, remaining-6 s)`), R8 `Unhandled()` before the plan runs, one effect via `Bound.Request` (TR.C holds), post-effect work on `WithoutCancel` plus 4 s bound, and every post-effect failure typed `EffectsUnrecordedError`; `dispatchError` maps it for MCP and A2A. v2 record carries `plan` and ordered `effects`. AC2.1-AC2.6 each have a named test (TestEffectDispatchProbe, UndeclaredTriple, UnregisteredDeclared, ConflictUnrecorded, HandlerCapCommitsFailed, CallerCancelledAfterEffect, Replay with nil-registry replay session). AC2.5b (rc 3, `limit.reason=="timeout"`) is in `TestAilangRunOutcomeShapesOnTheToolBinary` on the real binary.
- **M3 budget**: `store.EffectSpend` sums over the `effect:<ep>:` intent range, grouped by the intent's own Effect/Scope; seeding is `Budget - spent` floored at 0 under a per-episode lock. AC3.1, AC3.3, AC3.4 are dispatch-level tests. **Partial (AC3.2):** survive-reopen is proven at the store level (`TestEffectSpendSurvivesReopenAndReadOnly`) but there is no coordinator-level "three dispatches, reopen, denied" test.
- **M4 handlers and wiring**: `handlers_ailang.go` matches section 4.3: fixed op allowlist per effect, case-insensitive `op` alias refusal (a hardening beyond the design), scope check, `--` before the run path, four run-outcome shapes with no default, minimal child env (no registry credential reaches the tool), `AILANG_CACHE_DIR` outside the worktree, tool binary archived, hash re-verified before every execution (`verifiedTool`), exact 7-key rendered policy, policy/cache/store/archive/examples all refused inside the root (D4 rule). AC4.1-AC4.5 covered (`TestAilangToolConfinementMatrix` real-binary matrix incl. worktree topology and cache row; `TestWorkspace*`).
- **M5 package**: 8 contracted modules, `transitions.json` publishes 8 descriptors (`TestSeToolsManifestPublishes`), A2A and MCP e2e including visibility by grants and a mapped post-effect failure. AC5.1-AC5.4 met. Note AC5.4's "injected R14 over /a2a/" is exercised as an injected post-effect failure (`TestSeToolsA2APostEffectFailureIsMapped`), same typed error and mapping.
- **M7**: QUICKSTART section 9 added; `TestQuickstartSection9FlagsMatchTheCLI` binds flags to `serve --help` and the mint parser.
- **M6** (attended smoke): correctly not claimed; QUICKSTART says "pending first verbatim run".
- **Break fixes** (mint grant expiry, builtins text inventory, `--examples-dir`): each is measured (V64-V66), pinned by a test, and recorded in the design. The mint fix is correct in both layers (CLI sets expiry, `authority.Mint` clamps; mutant E killed).

Design fidelity deviations (minor, all documented): the examples corpus is an operator flag and not built in (V65 corrects V17); `builtins-search` loses signature/description (V64); `ReplayGrant` reproduces recorded denials with a synthetic grant (a workable but slightly clever shim); replay has **no production caller** (library plus tests only, as AC2.6 requires, but there is no operator verb to run a replay).

## 5. Security spot-check

Attempted paths: read/write outside the episode worktree, `.git`, policy self-modification, budget overrun, broker bypass. I read the handler/daemon/coordinator code and drove the real v0.51.0 `policy-tool` under the exact rendered policy in a `git worktree` fixture (sibling secret file, symlink out, `.git` pointer file).

Held (no bypass found):
- `..`, absolute, symlink-out reads: `refused ... escapes sandbox`; `x/../../esc` write refused; sibling-episode secret unchanged.
- `.git` and `.GIT` writes refused; `.gitmodules` refused by the deny glob; `tree ..`, `tree /`, `check /etc/hosts`, `policy_check /etc/hosts`, `iface /etc/hosts`, `test ../x.ail`, `fmt ../outside.ail` refused; unknown flags refused per op (`docs_search --dir` not admitted).
- Policy and cache live outside the root and are refused at startup if inside (`handlers_ailang.go` `CheckPolicyOutsideRoot` plus `daemon/workspace.go` `resolveWorkspaceRoot`); mutant J killed.
- Episode comes from the session, never the tool arguments; grammar plus `EvalSymlinks(root/episode) == root/episode` (`workspace.go` `episodeRoot`); mutants D, N killed.
- Every tool effect goes `Bound.Request` -> `BoundInvoker` -> broker; denial and budget are recorded; seeded budgets survive restart; the lock serialises same-episode calls (mutant I killed).
- Child environment is `HOME`, `PATH=/usr/bin:/bin`, `LANG`, `LC_ALL`, `AILANG_CACHE_DIR` (+ `AILANG_EXAMPLES`): the session credential cannot reach a tool.

Findings:

1. **`ailang-cli` `fmt` with `flags:{"write":""}` bypasses `fs_deny_write` and the grant split (P2, non-blocking, needs a ruling).**
   Measured on v0.51.0 under the exact rendered policy:
   `{"op":"fmt","path":".claude/x.ail","flags":{"write":""}}` and the same for `.ailang/y.ail` return `ok:true`, and the files were rewritten, although `.claude/**` and `.ailang/**` are in `fs_deny_write`. A path outside the sandbox is refused. `fmt` is in the policy summary's `cli` list, which the handler uses as the `Ailang.CLI` allowlist verbatim (`host/broker/handlers_ailang.go:loadSummary`, `opsFor`). Consequences: (a) the deny-glob invariant that R-SE-4 and design section 4.3 rely on does not hold for existing `.ail` files under denied globs (impact is low: it only reformats an already-present `.ail`); (b) a session holding only `Ailang.CLI` can mutate worktree files, which `Workspace.Write` is meant to gate, so the "6 grants cover 8 tools" budget and capability story has an undocumented write path. It is an upstream policy-tool gap rather than a World confinement defect, but the design says World adds only the allowlist, and this allowlist admits a mutating op with a mutating flag. Suggested fix, non-prod-affecting: refuse `fmt --write` (or drop `fmt`) in the handler and in `cli.ail`, plus an upstream issue and a new residual; also add a regression row to the AC4.1 matrix. No test currently covers it.
2. Reading the `.git` pointer file through `ailang-read` returns the absolute path of the operator's main repository gitdir (information disclosure only; the pointer is write-protected, V54). Informational.
3. `handlerBudget` can go negative when under 6 s remain (`effectful.go`); the handler context then expires immediately. A pre-dispatch expiry returns the raw error (nothing durable), and a post-dispatch one is recorded `failed`, so this is safe but untested. Informational.

## 6. Non-vacuity: my own mutants (14, five milestones; each run against the relevant packages, restored, `git diff --quiet` clean after every one)

| # | Mutant (production) | Result | Killed by |
|---|---|---|---|
| A | handler scope check removed (`handlers_ailang.go`) | KILLED | TestAilangToolRefusesAliasedOpKeyScopeAndRunPayloadKeys |
| B | EffectSpend seeding skipped (`effectful.go`) | KILLED | TestEffectBudgetPersistsAcrossDispatches, TestEffectConcurrentSameEpisodeBudget, TestEffectReplay |
| C | v2 record drops `effects` refs (`plan.go`) | KILLED | TestEffectDispatchProbe, TestEffectHandlerCapCommitsFailed |
| D | episode `EvalSymlinks == want` check removed (`workspace.go`) | KILLED | TestWorkspaceRegistryRefuses..., TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect |
| E | mint grant-expiry clamp removed (`authority/mint.go`) | KILLED | TestMint_GrantNeverOutlivesItsSession |
| F | finish key law (`world`/`effects`) disabled (`effectplan.go`) | KILLED | TestEffectPlanSketchDrift, TestParseFinishLaws |
| G | run branch omits `--` | KILLED | TestAilangRunArgvPutsDashDash..., TestAilangRunOutcomeShapesOnTheToolBinary |
| H | replay binder built Live (`binder.go`) | KILLED | TestEffectReplay, TestOpenReplayBinderNeverDispatches, TestBoundInvokerUnhandled |
| I | per-episode lock a no-op (`effectful.go`) | KILLED | TestEffectConcurrentSameEpisodeBudget |
| J | `CheckPolicyOutsideRoot` never refuses | KILLED | TestCheckPolicyOutsideRootRefusesInsideAndSymlinkedIn, TestNewAilangToolHandlerRefusals, TestWorkspaceRootContainingDaemonStateRefusesStartup |
| K | `op`-alias refusal removed | KILLED | TestAilangToolRefusesAliasedOpKey... |
| L | `.ailang/**` dropped from rendered deny list | KILLED | TestRenderEpisodePolicyHasExactlyTheDesignKeySet, TestAilangToolConfinementMatrix |
| M | op allowlist accepts any op | KILLED | TestAilangToolOpAllowlistIsFixedPerEffect, TestAilangToolReadEffectRefusesWriteOp |
| N | episode grammar check removed | KILLED | TestWorkspaceRegistryRefuses..., TestWorkspaceEpisodeEscapeIsR8BeforeAnyEffect |

Result: 14 killed, 0 survivors. Two harness errors on my side were corrected before counting: my first versions of I (wrong splice point, survived) and J (build failure) were invalid mutants; both were rewritten and re-run, and both were killed.

## 7. Blocking findings

None.

## 8. Non-blocking notes (ordered by value)

1. `ailang-cli` `fmt --write` gap (section 5, finding 1): needs an owner ruling; file upstream and add a matrix row.
2. AC3.2 is proven at store level only; add a coordinator-level reopen test (design wording: "the same across a store reopen").
3. Code-quality hygiene: `Replay` (about 130 lines) and `dispatchEffectful` should be split; `handlers_ailang_test.go` (940 lines) should be divided.
4. Design doc header still says "PLANNED ... Not implemented"; update the status to reflect M1-M5/M7 implemented and M6 pending when landing.
5. Replay has no production caller or operator verb (only tests); clause 1/5 value claims that depend on operators replaying effectful invocations are not yet exercisable outside tests.
6. M6 (attended smoke with pi and Claude Code, real publish, `/v1/log` resolution of record refs) is outstanding by design; clause 4's World arm and clause 5's corpus depend on it, not on this PR.

## 9. Independence statement

Judge: Claude Sonnet (fresh context, wrote none of this code, read no executor reasoning beyond commit messages and the design). Executors: Claude Opus/Sonnet agents. Controller: Claude Opus, attended with the owner present. This is the same vendor across author, controller and judge; flagged per D-WORLD-48 as a preference, not a blocker. I did not push, modify the PR, run `ailang messages` or fleet CLIs, or touch any real store or `world-publish`. My worktrees were removed; the only file I added to the main checkout is this report (uncommitted).

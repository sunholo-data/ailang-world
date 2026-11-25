# Evaluator report r1 - row 243 V4 S5c (iteration 433)

Head `272d8591` (plan `987cfaf6`, M1 `5e884ac5`, M2 `272d8591`), base `606f6e2c`. Judged in `.wt-world-iter433-row243-v4s5c-eval`. Time budget was tight, so some survivors were run with named tests only; this is stated per row.

## Verdict: PASS - 87 / 100, no BLOCKING findings

| Category | Score |
|---|---|
| Correctness vs design / plan §9 | 26 / 30 |
| Tests (non-vacuity, kills, portability) | 24 / 30 |
| Security / fail-closed | 19 / 20 |
| Coding standards | 9 / 10 |
| Scope and size discipline | 9 / 10 |
| **Total** | **87** |

## BLOCKING findings
None.

## Non-blocking findings (ranked)
1. **M2 does not pin that the `workspace` and `hosts` parameters of `operatorStartFacts` are forwarded to the seam.** `workspaceRoot: workspace` -> `""` (`operator_start_input.go:20`) SURVIVED the WHOLE `cmd/ailang-worldd` package (ledger `N-WORKSPACEROOT-EMPTY`, rc=0). `hostsDir: hosts` -> `""` survived the named run (`N-HOSTSDIR-EMPTY`; not re-run whole, because the fixture is codex-only and the seam reads `hostsDir` only for non-codex launchers, `episode_preset.go:146`). Swapping the two (`N-WORKSPACE-SWAP`) also survived named. Inert today, but when S5d adds a caller, a wrong root silently drops the preset-inside-workspace layout refusal. Add a test (preset inside the workspace root refuses, and a non-codex agent needs `hosts`) before the first production caller.
2. Mutants that survived the named run and that I judge equivalent or fixture-bound (named-only, not whole-package): `N-STORE-RAW-CRED` (skip `CanonicalPath` of `req.credPath`; `credPath` is already symlink-free per `session_new.go:267`, so equivalent), `N-WORKTREE-DERIVED` (`req.dir` vs workspace/episode, equal in the fixture), `N-ATTEMPT-ZERO` (the fixture attempt has only `Catalogue`), `N-CLOCK-CONST` (the clock does not reach a signed fact), `N-RES-ABS` (`IsAbs` drop; redundant with the later `runnerPath` and base checks). None is a signed-field comparison.
3. Reconstruction does not re-derive from live input: `ID`, `Argv`, `Options`, `Review`. These are deterministic functions of the signed fields (`OperatorArgv`, `Validate()`), not live inputs. The task is deliberately taken from the signed bytes, so a changed user task file is never read (plan §3.1). Acceptable, but the future catalogue (S7) must sign exactly the 12 names; `TestOperatorStartFacts` pins them.
4. CI cost: the controller's `-race` leg for `cmd/ailang-worldd` was **389 s** (limit 480 s; S5b reading was about 347 s). About 6 s of the delta is the 19 arms (0.3 s each); the rest is likely load. Headroom is about 90 s. No real sleeps in the new tests (grep for `time.Sleep` returned nothing).
5. Planner's NV-M2 mutant (`return in, err` early) fails `go vet` (unreachable code). I used `N-LOOP-DISABLED2` (`if false && signed.Resolution[i] != fact`) as the non-vacuity proof instead; it is KILLED.

## M1 audit
- Survivors re-measured against the package that holds the test (ledger `$J/judge/mutants.tsv`):
  - `E-TX-STORE-ID` KILLED by `TestOperatorMintScope/store_copy` (`host/authority`)
  - `E-TX-KIND` KILLED by `TestOperatorMintScope/finish_kind` (`host/authority`)
  - `E-RES-STORE-ID` KILLED by `TestOperatorMintRecovery/store_identity`
  - `E-RES-TOKENPATH-BASE` KILLED by `TestOperatorMintRecovery/token_name`
  - `E-RES-REGULAR` KILLED by `TestOperatorMintTokenReadback/fifo`; `E-RES-REGULAR-DIR` KILLED by `/directory`
  - `J-PRINT-AGENT` KILLED by `TestPresetStartReadHeadless`; `J-DOUBLE-OVERRIDE-APPEND` KILLED by `TestPresetStartReadParity`
  - All five killers are behavioural (stored-row rewrite, `VACUUM INTO` copy, real FIFO/dir), not inventory/string tests. The fifo/directory arms also assert the exact message `mint token is not a regular file`, which is what separates the type refusal from the later mode refusal.
- **F1 is true.** `E-TX-KIND` dies only through refusal order: with no runner row, the mutant falls through to `ErrOperatorRunner` instead of `ErrOperatorConflict`. With a current runner the `a.Options.Headless` guard refuses a finish anyway (`Validate` forbids a headless finish), so the mutant is equivalent there. The test comment says so.
- The 1 deleted line is the `[]string{...}` list of `TestOperatorMintTokenReadback`, replaced by the same list plus `"fifo"`. No assertion weakened.
- Portability: `syscall.Mkfifo`, `O_NONBLOCK` open of a writer-less FIFO and `Mode().IsRegular()` behave the same on Linux and Darwin; no uid/mode value is macOS-specific. `VACUUM INTO` is plain SQLite.
- `E-RES-EPISODE-CONFLICT` / `E-RES-REREAD` honestly unkillable in one process (an interleaving between the resolver's first read and its lock). They are carried in Appendix C S5f/S6 with "must kill before `MintOperatorStart` gets a production caller". Not run by me (plan ledger only).

## M2 signed-field enumeration vs reconstruction
Signed `OperatorAction` fields: `Schema`, `ID`, `Kind`, `Episode`, `Store`, `Argv`, `Options`, `Resolution` (12 facts), `Task`, `Review`.
- `Kind`/`Options.Headless`: guard, `ErrOperatorConflict` (mutants `N-KIND-DROP-GUARD`, `N-HEADLESS-INPUT-DROP-GUARD` KILLED by `/finish`, `/interactive`).
- `Store`: live `CanonicalPath(credPath)` compared to signed; `BindingError{store}` (`N-STORE-SKIP2` KILLED by `/store`).
- `Episode`, `Options.Preset`, `Task`: inputs to the live seam. Their effect is signed through `mint_start` (full `EpisodeStart` hash: repo, remote, base branch/SHA, preset sha/bytes, branch, agent, model, overrides, task path+digest, checks, TTL).
- `Resolution`: all 12 names compared in order with exact struct equality; length difference is `BindingError{resolution}`; field is `fact.Name`. Each fact has its own arm or independent fixture value (`launch_binary`, `launch_wall`, `mcp_path`, `mint_catalogue`, `mint_grants`, `mint_profile`, `mint_project`, `mint_start`, `mint_token_path`, `served_build`, `served_interpreter`, `worktree_path`). Facts are derived from the live seam read, not from the signed record.
- Gap: the `interpreter`, `mcp_path` and `worktree_path` arms edit the SIGNED value (the live value is unchanged, because the fixture has no way to change it), so those arms prove the comparison exists, not the value source. The value source is pinned by `TestOperatorStartFacts` (`N-MCP-FROM-BINARY`, `N-BUILD-TOOLREF`, `N-TOKENPATH-EMPTY` KILLED there).
- Exactness: struct `!=`, no prefix/case/normalisation. `N-NAME-ONLY` and `N-VALUE-LEN-ONLY` (weakened comparisons) both KILLED by many `TestOperatorStartChanged` arms.
- Refusal read-back: every arm checks `observeProvision` equality, plus no lock path, no `ran` marker and no `ws/ep` worktree. Class AND state are asserted (field via `errors.As`, class via `errors.Is`).

## Inertness
Production callers (excluding tests): `reconstructOperatorStart` 0; `operatorStartFacts` 1 (inside reconstruct, positive control); `MintOperatorStart` 0; `ClaimOperatorAction`/`ReserveOperatorStep`/`ObserveOperatorStep` 0; `CommitOperatorMint` 1 (`host/authority/mint.go:219`, inside the uncalled `MintOperatorStart`); `AbortOperatorMint` 1 (resolver); `resolvePresetStart` 2 (existing terminal path and reconstruct). Controller `callers.txt` agrees. New file has no write, lock, exec or network call. `.ail`, `app`, `host/daemon`, `daemon.go`, `surface_page.go` unchanged (`git diff --stat` empty); `routes.txt` rc=0.

## Split assessment
401 of 1,400 by `git diff -M --numstat 606f6e2c HEAD -- cmd host docs app`; I measured: M2 test 229, M2 production 81, M1 `operator_mint_test.go` 42/1, `host/authority/operator_mint_test.go` 40, `episode_preset_read_test.go` 8 (sum 401; 91 at M1 by the same files). The cut is conservative, but defensible: worktree reservation, file publication and the consumer each need injected-failure matrices and an attended remedy in the same patch, and the plan has probe measurements (59/53/72/132) showing they are unmeasured as a whole. Appendix C carries S5d (worktree + attended disposition in the same patch), S5e (publication), S5f/S6 (the four named process pins, the two interleaving survivors, terminal `session revoke`, credential revalidation), and `J-OWN-UID`, `Z-ABSENT-ANY`, `Z-WAIT-READBUDGET`. Nothing silently dropped. Mild under-scope is possible (headroom 999 lines), but the plan's reasoning is evidence-based.

## Mutant table (ledger `/Users/voightkampff/.ailang/state/world-iter433/judge/mutants.tsv`, 45 lines; per-run logs in `judge/runs/`)
Own-only ids (`comm -23` vs plan ids): 31 (not counting the 11 shared ids re-measured).
Controls: **must-die** `N-LOOP-DISABLED2` KILLED (`TestOperatorStartChanged` base/build/catalogue/...); **must-live** `N-SWAP-EQ-OPERANDS-LIVE` SURVIVED (named run only).

| Id | Result | Killing test (package) |
|---|---|---|
| E-TX-STORE-ID, E-TX-KIND | KILLED | `TestOperatorMintScope/store_copy`, `/finish_kind` (host/authority) |
| E-RES-STORE-ID, E-RES-TOKENPATH-BASE, E-RES-REGULAR, E-RES-REGULAR-DIR | KILLED | `TestOperatorMintRecovery/store_identity`, `/token_name`, `TestOperatorMintTokenReadback/fifo`, `/directory` (cmd/ailang-worldd) |
| J-PRINT-AGENT, J-DOUBLE-OVERRIDE-APPEND | KILLED | `TestPresetStartReadHeadless`, `TestPresetStartReadParity` (cmd/ailang-worldd) |
| X-COMPARE, X-LEN | KILLED | `TestOperatorStartChanged/preset_ttl`, `/extra_fact` |
| NV-M2 | COMPILEERR (vet) | replaced by `N-LOOP-DISABLED2` |
| N-LOOP-DISABLED2, N-NAME-ONLY, N-VALUE-LEN-ONLY, N-FIELD-CONST | KILLED | `TestOperatorStartChanged` (many arms) |
| N-STORE-SKIP2 | KILLED | `TestOperatorStartChanged/store` |
| N-SKIP-LAST, N-SKIP-FIRST | KILLED | `/worktree_path`, `/launch_binary` |
| N-KIND-DROP-GUARD, N-HEADLESS-INPUT-DROP-GUARD | KILLED | `/finish`, `/interactive` |
| N-TASK-LIVE-DRIFT, N-MCP-FROM-BINARY, N-TOKENPATH-EMPTY | KILLED | `TestOperatorStartFacts` |
| N-LAUNCHING-FALSE, N-PROFILE-FROM-SIGNED-FREE, N-BUILD-TOOLREF | KILLED | `/launch_binary`+Facts, `/profile`+Facts, `/build`+Facts |
| N-RES-NLINK | KILLED | `TestOperatorMintTokenReadback/link` |
| N-TX-HEADLESS-DROP, N-TX-STORE-ALWAYS-NE | KILLED | `TestOperatorMintBinding/start_options`, `/artifact`... (host/authority) |
| N-WORKSPACEROOT-EMPTY | **SURVIVED (whole package, rc=0)** | finding 1 |
| N-HOSTSDIR-EMPTY, N-WORKSPACE-SWAP, N-STORE-RAW-CRED, N-WORKTREE-DERIVED, N-ATTEMPT-ZERO, N-CLOCK-CONST, N-RES-ABS, N-SWAP-EQ-OPERANDS-LIVE | SURVIVED (named run only) | findings 1, 2 |
| N-STORE-SKIP, N-LOOP-DISABLED, N-RES-EPISODE-PATTERN | COMPILEERR (first try) | rewritten as `...2` where relevant |
| N-CLOCK-CONST, N-TX-KIND-INVERT | INVALID first try | N-CLOCK-CONST re-run as above; `N-TX-KIND-INVERT` not retried (anchor occurs twice) |

Every KILLED row was produced by `--- FAIL:` of the named behavioural test in the package holding it; none relied on an inventory/AST/string test. Plan rows not re-measured by me (X-*, remaining 30) were replayed by the controller, not by me.

## Non-vacuity
- M1: each M1 test fails against its survivor mutant (above) and passes at head (`t1.txt`, `t2.txt`).
- M2: `TestOperatorStartChanged` fails with the comparison loop disabled (`N-LOOP-DISABLED2`), with only the names compared (`N-NAME-ONLY`) and with the first/last fact skipped.

## Test runs (rc)
- Controller sweep `gates-final-272d8591/legs.txt`: all legs rc=0, GATES_DONE, no FAILS; `observation-before.txt` == `observation-after.txt` (real LaunchAgents, `~/.ailang-world`, credentials store unchanged); `verify_go` rc=0; race `cmd/ailang-worldd` 389 s.
- Mine: `go test -count=1 ./host/store/ ./host/authority/` rc=0; `go test -count=1 -run 'Operator|Preset|Session|Quickstart|Help|CredentialWriters' ./cmd/ailang-worldd/` rc=0 (34.1 s); `go test -race -count=1 -run 'OperatorStart|OperatorMint|PresetStartRead' ./cmd/ailang-worldd/ ./host/authority/` rc=0 (28.4 s, 10.6 s); `go vet ./cmd/... ./host/...` rc=0; `go test -run '^$' ./...` clean (no non-ok lines).
- Worktree `git status --short` clean after all mutants; no scratch files added; no real-path writes, no `launchctl`, `service`, `session` or `approve` commands run.

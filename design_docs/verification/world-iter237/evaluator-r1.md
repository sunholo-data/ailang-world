VERDICT: PASS 92/100

# Row 136 `w-mcp-effects-unrecorded-is-untyped`: evaluator round 1 (iteration 237)

- Judge: claude-opus-5-5, independent. The executor was claude-sonnet-5-5.
- Subject: `sprint/row136-mcp-typed-errors` @ `18a4e9c`. Range `0ce6e43..18a4e9c`: plan `f85812e`, M1 `e480ddc`, M2 `7a0ede1`, M3 `18a4e9c`. Base is origin/dev `1bf36bd`.
- Worktree: my own detached checkout at `/Users/voightkampff/dev/sunholo-data/.wt-world-iter237-eval`. It was clean after every mutant and has been removed.
- Evidence (logs, mutant script, CI log): `~/.ailang/state/mission-world-iter237-evidence/evaluator/`.

## Rubric

| Category | Pts | Score | Basis |
|---|---|---|---|
| Tests pass | 20 | 18 | Locally, 1 FAIL out of 2612 RUN, in `host/broker`, a package this row does not touch. It passed when re-run alone. On PR CI, the first `-race` run failed `TestExecSrtMCPEndToEnd` (row 140, untouched). The `--failed` re-run of the same SHA went green. Both are pre-existing flake classes. −2 because a green CI result needed a re-run. |
| Lint / vet | 10 | 10 | `go vet ./...` rc=0 |
| Acceptance criteria | 30 | 28 | Every code AC (AC1.1–AC2.6, plus the plan's AC1.1a/AC2.3a) is implemented and kills its mutant. AC3.1: the issue is filed (ailang#1602, OPEN, `area:serveapi`, and its ask matches design §M3) and the docs link it. The plan gives the charter follow-up row (Gate-1 predicate) and the `ailang messages` note to the controller, and neither is in this range yet. −2 until the controller records them. |
| Code quality | 15 | 13 | The change is minimal and pure-core-clean: `writtenBy(surface)`, a strict `validateCall`, and one id function. −2 because 3 of the 4 relabelled MCP log sites (`tools/list` ×2 and the `tools/call -` admission refusal) have no test, so own-mutants survive there (NB1). |
| Documentation | 15 | 13 | The operator-line sweep (V15) and the provenance sweep (`coordinator:mcp`, `mcp:` ids) are consistent, M3's warning links ailang#1602, and no `<ISSUE-URL>` is left in shipped docs. −2 because the design doc's Status line still says "Nothing implemented yet" (NB3). The repo has no CHANGELOG, so nothing is deducted for one. |
| Design fidelity | 10 | 10 | Option C only. D1–D6 are measured and justified. No MCP wire change, no option D, no pin bump. |
| **Total** | 100 | **92** | PASS (≥70, no hard fail) |

## BLOCKING findings

None.

## NON-BLOCKING findings

1. **Three of the four relabelled MCP refusal sites are untested.** `host/projection/mcp.go:76`, `:81` (`tools/list -`) and `:107` (`tools/call -`, the admission refusal) were changed in M1, but no test asserts any of their lines.
   - Repro: rewrite `:107` to `logRefusal("a2a","mcp invoke","-",err)`, or `:76/:81` to `logRefusal("a2a","mcp tools","-",err)`, then run `go test ./host/projection/`. Both survive: 142 RUN, 0 FAIL.
   - The design's ACs only name the `:132` site, so this is not an AC miss. It is still a load-bearing gap for the format M1 documents. A one-line table test over a failing `allowedDescriptors` closes it.
2. **CI does not protect several mutant kills.** The `TestSeTools*` and `TestCLIAgainstSeToolsDaemon` e2e tests SKIP without `WORLD_TOOL_AILANG_BIN`, and CI never sets it (plan P3/D1). That leaves three gaps in CI:
   - AC1.1's "ref resolves to an allowed Workspace.Read EffectRecordV1" assertion. The CI companion AC1.1a checks only the cause text, not the `sha256:` refs.
   - AC2.1b (the `mcp:` id end to end under R14).
   - The `cliwalk_e2e_test.go:232` regex (MUT-CLIWALK / D2).

   Each mutant is still killed in CI by another test, except MUT-CLIWALK, which is test-only. The controller's recommended follow-up row (a tool binary in CI) should be filed.
3. **The design doc's Status line is stale.** `design_docs/planned/w-mcp-effects-unrecorded-is-untyped.md:3` still reads "Nothing implemented yet". It should be updated to Implemented, and the doc and plan moved together on landing.
4. **Controller-owned M3 residue is not yet in the range.** Three things remain:
   - The charter follow-up row with the verbatim Gate-1 predicate, with `<N>`=1602.
   - The `ailang messages send mission-control` note.
   - The iteration-log record of both (AC3.1).

   This must be done at record time.
5. **The fake floor fixture still writes the old tag.** `scripts/floor/fakeworld.py:50,75` still writes `coordinator:a2a` and `a2a:` ids. It is a fake with no consumer that checks the tag, so this is harmless. It now models the pre-136 shape for what are notionally MCP calls. Cosmetic.
6. **CI flake noted.** `TestExecSrtMCPEndToEnd` failed at `exec_e2e_srt_test.go:154` (`where = "<nil>"`) on PR run 37461985123 and passed on re-run. The test discards `wire.Error` on its second call, which hides the cause. This belongs to row 140's test, not this row.

## Mutant table

The required mutants and the killers I ran. I applied each one with `perl -0pi`, checked that `git diff` was non-empty, ran the killer, restored the file and checked for a clean tree. Every row left 0 dirty files.

| Mutant | Killer test (failing line) | Killed? | CI-protected? |
|---|---|---|---|
| MUT-MCP-LOG-NOCAUSE (`mcp.go:132` logs `errors.New("dispatch failed")`) | `TestMCPRefusalLineNamesSurfaceAndCause` (`mcp_conformance_test.go:403`) | YES | **yes** |
| MUT-MCP-LOG-NOCAUSE | `TestSeToolsMCPPostEffectFailureIsLogged` (`setools_e2e_test.go:826`) | YES | no (e2e, needs the tool binary) |
| MUT-MCP-SURFACE (b), Call literal only → `SurfaceA2A` | `TestMCPBatchAccounting` ×2, `TestMCPBatchConformance` ×2, `TestMCPAbsentPrecheckSuccessfulInvocation` (`mcp_conformance_test.go:183`, `:356`); 7 FAIL / 40 RUN | YES | **yes** |
| MUT-MCP-SURFACE (b) | `TestSeToolsMCPPostEffectFailureIsLogged` (`:826`, the D3 id-in-cause check) and `TestSeToolsSurfaceTaggedInLog` | YES | no |
| MUT-SURFACE-DEFAULT (`validateCall` defaults `""` → a2a) | `TestDispatchRefuses/R1_empty_surface` (`coordinator_test.go:425`) | YES | **yes** |
| MUT-WRITTENBY-CONST (`planCommit` `WrittenBy: "coordinator:a2a"`) | `TestPlanWrittenByNamesSurface` (`coordinator_test.go:208`) | YES | **yes** |
| MUT-WRITTENBY-CONST | `TestSeToolsSurfaceTaggedInLog` (`setools_e2e_test.go:897`) | YES | no |
| MUT-REST-ACCEPTS-MCP (`restInvocationID` also accepts `mcp:`) | `TestCommitInvocationReceiptIdentity/foreign_namespaces_are_refused` (`commit_budget_test.go:358`, 200 ≠ 400) | YES | **yes** |
| MUT-A2A-ID-FORMAT (`episode:surface:task`) | `TestNewRefusesMissingSeamsAndCaps` (`coordinator_test.go:846`) and `TestReconcileRefusesDamagedRecord`; 8 FAIL | YES | **yes** |
| MUT-WHY-PREFIX (drop `mcp:` from `parseWhyTarget`) | `TestWhyResolvesMCPInvocationID` (`why_test.go:259`) | YES | **yes** |
| OWN: A2A Call literal tagged `SurfaceMCP` (`projection.go:363`) | `TestA2ADispatch/success_and_R13_reconciled` and `TestA2A_AbsentPrecheckPublicationDispatch` (`projection_test.go:997`) | YES (in host/projection; host/daemon CI tests alone do NOT kill it: 411 RUN, 0 FAIL) | **yes** |
| OWN: `validateCall` refuses only `""` (accepts any non-empty surface) | `TestDispatchRefuses/R1_unknown_surface` | YES | **yes** |
| OWN: `mcp.go:107` admission refusal labelled `a2a`/`mcp invoke` | none (`./host/projection` 142 RUN, 0 FAIL) | **SURVIVED** | — (NB1) |
| OWN: `mcp.go:76,81` `tools/list` labelled `a2a`/`mcp tools` | none (142 RUN, 0 FAIL) | **SURVIVED** | — (NB1) |

The e2e-only kills that CI does not run are: the NOCAUSE ref-resolution assertion, the AC2.1b R14 end-to-end prefix, AC2.3's live-daemon header check (CI's AC2.3a covers the mutant), and MUT-CLIWALK (from the executor's run, not re-run by me).

## Gates (final tree `18a4e9c`, v0.41.0 pin at `~/.pinned-ailang/ailang`)

- `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh`: **rc=0**. "16 required identities verified, 40 named tests pass, 444 se-tools named tests pass."
- `go vet ./...`: **rc=0**.
- `AILANG_BIN=… go test -race -count=1 -v ./host/... ./cmd/...`: rc=1, with **2612 `=== RUN`**, 2555 PASS, 56 SKIP and **1 `--- FAIL`**.
  - The failure is `host/broker` `TestRunBoundedHeadTailTimeoutKeepsPartialOutput` (`handlers_capture_test.go:122`, partial stdout `""`).
  - `host/broker` is not in the diff and does not import coordinator or projection. Re-run alone with `-race -count=1 -run TestRunBoundedHeadTail`, it passes (rc=0). The executor saw the same flake at M1.
  - Every touched package (coordinator, projection, daemon, cmd/ailang-worldd) is `ok`.
- PR #219 CI on `18a4e9c`: the `ailang-code verify gate` passed. `go host build + test gate` failed on the first run (`TestExecSrtMCPEndToEnd`, row-140 test, untouched; the race-control DATA RACE warnings are the intentional `racecontrol` canary). The `gh run rerun --failed` of run 37461985123 completed **success**.

## Checks 3–6

- **Compatibility.**
  - An old-shape store still replays, reconciles and walks. `InvocationID(SurfaceA2A, ep, task)` is byte-identical to the old `"a2a:"+ep+":"+task`, and `TestNewRefusesMissingSeamsAndCaps` pins this.
  - Reconcile (`committed()`) and Replay compare stored strings. No code derives a prefix.
  - The only non-test readers of `WrittenBy` are display code: logtail, why, workbench and the handlers' JSON. None of them matches `coordinator:a2a`.
  - The why fixture's A2A path (`invoke` → `invokeAs(SurfaceA2A)`) commits the old bytes, and the `why a2a:ep1:x` and `log tail` tests pass.
  - The A2A wire is unchanged: `a2a_wire_test.go` is untouched and green, and its `metadata.invocation_id` is still `a2a:<ep>:<task>`.
  - A2A operator bytes are unchanged: every A2A site passes `"a2a"` plus `req.Method`, which gives the identical format.
- **No missed production caller of strict Surface validation.**
  - `\bCall\{` gives 6 hits: 2 production (`projection.go:363` A2A, `mcp.go:130` MCP), 3 test literals, and 1 `[]Call{` built from the helper. All 5 literals set `Surface`.
  - `.Dispatch(` has exactly two production callers, the same two.
  - `InvocationID(` has 3 production calls (`coordinator.go:267`, `mcp.go:132`, `projection.go:374`), each with the correct surface.
- **Scope.**
  - The diff touches no `tools/launchd/*`, no skills, and no `go.mod`/`go.sum`.
  - There is no MCP wire change: the frozen `-32603 "host callback failed"` is asserted by AC1.1.
  - There is no idempotency key.
  - `<ISSUE-URL>` appears only in the plan's instructions, not in shipped docs, and `mcp-and-a2a.md:207` links ailang#1602.
- **Coding standards.**
  - Surface tagging is pure in `plan.go` and `coordinator.go`.
  - Logging stays at the projection host boundary.
  - No local workaround of `mcphttp`: no middleware, and no use of `structuredContent` as an error channel. The gap is routed upstream (#1602).
  - No `.ail` changes.

EVALUATION_RESULT: pass
EVALUATION_SCORE: 92/100
EVALUATION_ROUND: 1
EVALUATION_REPORT_PATH: ~/.ailang/state/mission-world-iter237-evidence/evaluator-r1.md

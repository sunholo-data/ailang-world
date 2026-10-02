# Independent PRODUCT evaluation - row 108 (MCP dispatch projection), attended 2026-10-02

**Verdict: PASS - 82/100.** No hard fail. CI for the candidate is still PENDING (see below), so remote CI is unmeasured and gets no credit beyond the local gate.

Candidate: `a12fa14` (branch `attended/row108-rebased`, draft PR #176) on base `daba577`. Scored diff: `git diff daba577 a12fa14` (25 files, +2524/-63). Own worktree `/Users/voightkampff/dev/sunholo-data/.eval-attended-108` (HEAD asserted `a12fa14745855d27ed6efc94934469e67698b341`, outside the repo dir, removed afterwards, main checkout clean).

Independence: judge = Claude Sonnet via Agent tool, fresh context; code author = codex gpt-6.1-sol; controller = Claude Opus (attended). The judge wrote no production code.

## Gate results (own runs, pinned AILANG v0.41.0, AILANG_BIN and WORLD_PKG_AILANG_BIN set)

| Command | rc | Time |
|---|---|---|
| `./scripts/verify_ail.sh` | 0 ("16 required identities verified, 40 named tests pass") | 5 s |
| `./scripts/verify_go.sh` | **0** | 246 s |
| `gofmt -l` on all changed .go files | empty | - |
| `go vet ./...` | 0, empty | - |

verify_go.sh stages all ran: tracked-binary hygiene, mission-input gate, race known-positive control (2 expected `DATA RACE` lines = positive control), `go build`, 37-test evidence manifest, plain `go test ./... -count=1`, and `go test ./... -count=1 -race -timeout 8m`; every package `ok` (slowest: broker 118 s, daemon 127 s). The three failures the iter-220 judge hit (pkgproj iface_test.go:232, cmd cli_test.go:211, broker cleanup_test.go:200) did not recur here; the base-flake tranche (D-WORLD-51) plausibly explains this, but a single green run is a measurement of this run, not proof of flake elimination. I did not re-run on base `daba577` because nothing failed.

## Hard-fail check

- Mandatory Go profile green: yes. AIL gate green: yes.
- Frozen core touched (`tools/launchd/*`): no (diff list contains none). Skills copied: no.
- Production Go: no vacuous gate observed (mutants below kill). No effect-in-core violation seen: MCP adapter is a host-boundary file (`host/projection/mcp.go`), pure mapping in `mcpname.go`.
- Tracked-binary hygiene gate: passed inside verify_go.
- Result: **no hard fail.**

## Mutants (own, production-code mutants, each restored; `git diff --quiet` confirmed, final worktree clean)

| AC | Mutant | Result |
|---|---|---|
| MAP | `EncodeMCPName` replacer replaced by identity (`name := id`) | KILLED: TestMCPNameRoundTrip, TestMCPNameRefusal, TestMCPSurfaceRefusalRecovers |
| ABSENCE | `allowedDescriptors` catch-all: `if !hasHead {` (drops typed `errors.As`) | KILLED: TestMCPAbsentPrecheckSnapshotFailure |
| INVOKE/ADMIT | `Invoke` readmission skipped (`if false && !allowed`) | KILLED: TestMCPListedThenTrulyAbsent (note: TestMCPInvokeListedThenRevoked did NOT fail alone; killing rests on one test) |
| ROUTE | mount `"POST /mcp/"` -> `"/mcp/"` | KILLED: TestMCPMountMethodSource |
| SCHEMA | missing-`type` default `"object"` not inserted | KILLED: TestMCPSchemaNormalization |

5/5 killed, 0 survivors, all failures verified as test failures (not build failures). Not attempted: BOUNDS/BATCH timing mutants (the 33 s production-constant legs are expensive; the iter-220 judge recorded the author's aggregate mutants at 32.61 s / 34.62 s REDs).

## Per-AC assessment

| AC | Assessment |
|---|---|
| DEP | go.mod pins ailang v0.47.2 (narrow protocol seam: `serveapi/protocol`, `hostcall`, `mcphttp`); `go build`/vet clean; dependency gate part of verify. Supported. |
| MAP | Reversible escaping `_`->`_u`, `.`->`_d`, `/`->`_s`; strict decoder rejects bad escapes; >64-byte whole-surface refusal; killed by mutant. Supported. |
| SCHEMA | `normalizeMCPSchema` preserves raw bytes when `type:object`, inserts type via RawMessage (no float64 round trip), refuses non-object. Killed by mutant. Supported. |
| CARRIER | Shared resolver, 401 mapped to AuthorizationError with constant messages; TestMCPDenialMatrix / MountedBearerDenial present. Supported. |
| ADMIT | `Tools` reads fresh `allowedDescriptors` per call under `maxWait`. TestMCPListExactSetPerSession. Supported. |
| ABSENCE | `RegistryHeadAbsentError` typed error in transitionreg + `errors.As` in projection; empty only if precheck absent AND typed absence; post-read ctx check added. Killed by mutant. Supported. |
| INVOKE | Object-only args, fresh admission per Invoke, membership check, coordinator nil check, `mintTask` after admission, real `coord.Dispatch`. Supported. |
| ITEM | Per-Invoke fresh task id and snapshot (batch accounting test present). Supported by code reading and passing test; not mutated by me. |
| BATCH | Upstream handler owns arrays/versions; World adds no batch logic; conformance tests with honest calibration decorators. Supported. |
| BOUNDS | Config validation (positive scalars, InvokeWait < WriteWait), aggregate deadline in `Handler.MCP`, inner budgets, production-constant test (`TestMCPPostBudgetProductionConstants`, 3/10/20/30 s, slots 8) passed in the green run. Supported; not mutated by me. |
| WIRE | Thin context wrapper over `mcphttp.NewHandler`; source-ownership test + SSE golden. Supported. |
| ROUTE | `POST /mcp/` registered in `Daemon.Handler`; AST test + GET 405 test; killed by mutant. Supported. |
| CROSS | TestMCPCrossSurfaceExactSet / TestMCPToolsListMatchesAgentCard on one daemon. Supported. |
| FIXTURE | TestMCPAmbientExportsAbsent, TestMCPQuickstartPayloadsVerbatim; QUICKSTART and HOST_CHANGELOG updated. Supported. |

## Literal rubric

| Category | Points | Reason |
|---|---:|---|
| Tests | 20/20 | Exact mandatory profile (verify_ail rc=0, verify_go rc=0 incl. race) green; new tests kill 5/5 mutants. |
| Lint | 10/10 | gofmt empty, go vet rc0. |
| Acceptance | 28/30 | All 14 ACs have functional, mutation-backed or passing-test support. -2: remote CI still pending; BOUNDS/BATCH/ITEM timing legs rest on author evidence plus one green local run, not my own mutants. |
| Code quality | 0/15 | Literal -5 each for three touched files over 800 lines: `host/daemon/daemon.go` 966, `host/daemon/daemon_test.go` 1343, `host/projection/projection_test.go` 1553. All were already oversized on base; the row adds modest lines to them (daemon.go +7). New MCP code is small and well factored (mcp.go 143, mcpname.go 87). |
| Documentation | 15/15 | HOST_CHANGELOG entry, QUICKSTART MCP section, design status; comments accurate. |
| Design fidelity | 9/10 | Thin released-seam adapter exactly as designed; World owns only admission, mapping, deadlines. -1: minor non-idiomatic bundling in `daemon.go` Config literal (one-line multi-field) and `Handler.MCP` comment says "single aggregate request deadline" while ownership split relies on InvokeWait. |
| **Total** | **82/100** | Pass threshold 70; no hard fail. |

## CI state (`gh pr checks 176`, read-only, at evaluation end)

- `ailang-code verify gate`: pass (46 s)
- `go host build + test gate`: **pending** (run 36992316211). Unmeasured; PASS here is conditional on this job going green (it runs the same verify_go profile that I measured rc=0 locally).

## Blocking findings

None.

## Non-blocking notes

1. Code-quality 0/15 is inherited: daemon.go, daemon_test.go and projection_test.go were already >800 lines on base. Splitting them is a separate refactor row, not row 108 scope.
2. Single-test kill on the Invoke readmission mutant (only TestMCPListedThenTrulyAbsent fails; TestMCPInvokeListedThenRevoked stays green). Suggest strengthening the latter, row-level residual.
3. A single green full-profile run does not prove the base flakes (pkgproj descendant pid, CLI daemon announce, broker overflow EPERM) are gone; the iter-220 judge saw them red on the pre-tranche base. Watch the candidate-SHA CI job and consider repeated runs before merge.
4. Residuals already recorded by earlier judges stand: slow-body 30 s socket truncation (no delivery claim), callback-capacity / nil-coordinator supporting-hunk survivors, guide initialize recipe absent.
5. I scored row 108's diff only; the D-WORLD-51 base-flake commits were not evaluated beyond noting they sit under a green profile.

## Disposition

PASS 82/100, no hard fail, no blocking findings. Merge-readiness requires only the pending `go host build + test gate` check on PR #176 to finish green.

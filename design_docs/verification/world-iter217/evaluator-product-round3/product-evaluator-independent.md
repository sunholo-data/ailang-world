# Independent PRODUCT evaluation — iteration 217, round 3 (FINAL)

**FAIL — 59/100 (exact 58.545), hard test failure.** Tested immutable full candidate `741e3c371de3e558243e4ee77a09eeaddbac50a3`. This is an independent cross-vendor product judgment, not controller acceptance. Round 1 FAIL 49/100 and Round 2 FAIL 54/100 remain preserved as independent evidence, never acceptance.

**Judge model**: pi:ollama/deepseek-v4-pro:0813-cloud (cross-vendor fallback flag; generator/executor/controller were GPT6.1Sol). Native token usage is unknown (`null`), not zero. No subagents, no production fixes, no remote messages, no pushes, no shared-harness writes. Own detached worktree `HEAD` verified exact before reads; only this report directory is untracked.

## Verdict summary

The Round 3 repair is **test-only** and correctly eliminates the MCP cooperative-slot recovery race (independently reproduced green normal×10 and race×10). It does **not** touch production code, budgets, pins, or contracts. However, the controller's exact **unchanged full Go profile is still rc1** (`scripts/verify_go.sh`, `go test ./... -count=1`), so the hard test gate fails and the sprint remains unaccepted. This is the final allowed round; the verdict is FAIL and the row is parked, not a request for a fourth judge.

## Blocking evidence

Controller marker `controller-final-profiles-terminal.json` arrived with `terminal:true`, `candidate_sha=741e3c371de3e558243e4ee77a09eeaddbac50a3`, and `runtime_allowed_now:true`. Root exact unchanged profiles (released pin first in PATH, `GODEBUG` unset, default Go cache, normal parallelism):

- `bash scripts/verify_ail.sh` → rc0 (26.17s)
- `bash scripts/verify_go.sh` → **rc1** (243.77s)

The full Go profile failed at the **plain** `go test ./... -count=1` leg (the race leg was never reached under `set -euo pipefail`). Five failing functions across four packages, all timing/deadline-bound:

| Test | Location | Failure |
|---|---|---|
| `TestCLIRealSubprocessEpisode` | `cmd/ailang-worldd/cli_test.go:266` | `/v1/commit` HTTP 503 "commit deadline (3s) exceeded before the durable step" |
| `TestF5WallClockTimeoutHasElapsedBound` | `host/capsule/capsule_test.go:230` | `timeout returned after 2.092940875s, want <= 2s` |
| `TestDispatchDurableDeadline/Commit` | `host/coordinator/coordinator_test.go:624` | `request ran past deadline: 2.504338375s` |
| `TestCommitBudgetAndUncertainReconcile/ii_A_landed_uncertain_then_B_on_top` | `host/daemon/commit_budget_test.go:209` | `commit B: status 503 ... commit deadline (3s) exceeded` |
| `host/daemon` (cascade) | — | "projection: invocation coordinator is unavailable", "sql: database is closed" during teardown |

**Inheritance, not sprint regression.** `git diff --name-only 54e0fb0..741e3c3 -- cmd/ailang-worldd/cli_test.go host/capsule/capsule_test.go host/coordinator/coordinator_test.go host/daemon/commit_budget_test.go` is **empty** — none of the failing test files are touched by the sprint. Independently reproduced in isolation (no parallel load): `TestF5WallClockTimeoutHasElapsedBound` count=3 rc0, `TestDispatchDurableDeadline` count=3 rc0. These are pre-existing wall-clock/commit-deadline tests whose 2s/3s bounds overshoot only under full parallel `go test ./...` saturation. They are **not** MCP regressions, but they still make the unchanged full profile red, which the rubric treats as a hard fail (see below).

## Round 3 repair audit (test-only, verified)

`git diff --stat 59870e58..741e3c3` shows exactly one production-path change: `host/projection/mcp_bounds_test.go` (+59/−0). Everything else is canonical status / sprint metadata / evidence. `go.mod`, `go.sum`, `host/projection/mcp.go`, `host/projection/projection.go`, `host/daemon/daemon.go`, `host/transitionreg/transitionreg.go` are all byte-unchanged from Round 2.

The repair replaces the always-on `delayedResolver` (30ms) with an `initialDelayResolver` (80ms, first call only), widens the aggregate budget to 300ms, CallbackTimeout to 2s, and adds explicit `returns`/`remaining` counters plus a same-handler one-slot recovery phase with an independent 3s budget:

```go
resolver := &initialDelayResolver{inner: cfg.Resolver, delay: 80 * time.Millisecond}
const initialBudget = 300 * time.Millisecond
cfg.InvokeWait = initialBudget
cfg.CallbackTimeout = 2 * time.Second
...
h.invokeWait = 3 * time.Second   // recovery: independent live test budget
```

Independent runtime (own sandbox, released pin, normal Go):

- `go test ./host/projection -run '^TestMCPInvokeDeadlineFreesSlot$' -count=10 -v -timeout=180s` → **rc0**, 10/10. Each run logs `initial aggregate=300ms ... entry remaining≈216ms ... entered=2 return-path=2 resolver=2` — the aggregate deadline propagates (~216ms left after the 80ms resolver delay), and exactly two executions/two returns/two resolver calls confirm the slot is released and re-acquired.
- `... -count=10 -race -timeout=180s` → **rc0**, no data race.

**Corrective causal calibration confirmed.** The prior Round 2 flakiness (always-on 30ms resolver delay + 80ms outer budget ⇒ the *second* request also delayed and hit the callback timeout) is removed by confining the artificial delay to the first request and giving recovery a fresh 3s budget. `remaining` (≈216ms) is bounded in `(0, initialBudget−40ms]`, proving the callback sees the aggregate deadline rather than a lost/fresh context.

## Sandbox limitation (documented, not a verdict)

My cross-vendor pi sandbox forbids TCP bind. `go test ./host/daemon -run '^TestA2AResendAfterWriteTimeout$' -count=3` fails at `httptest.NewUnstartedServer` with `listen tcp6 [::1]:0: bind: operation not permitted`. This is the same `broker EPERM` class Round 2 already classified as source-unchanged; it is an environment restriction, not a product defect. Consequently the **repaired socket** (`TestA2AResendAfterWriteTimeout`) and the **unchanged real-capsule sibling** (`TestA2AInvokeDeadlineClosesSocket`) could not be independently re-run here. Root's full-profile log shows host/daemon failed only on the commit-budget test (not the socket fixtures), consistent with those fixtures passing in the root's unrestricted environment; I rely on the root's raw log for them and do not claim my own socket run.

## Acceptance evidence

Sprint JSON carries 5 features / 22 acceptance criteria; `passes=true` on A1(5), A2(6), B(3), C(4) = 18 met; D `passes=null` (4 criteria not falsely marked accepted). Score `18/22 × 30 = 24.545`. Production budgets 3/10/20/30 and the same-handler/8-callback return-bound vs 8-process reap-bound split are preserved because no production file changed. Module `github.com/sunholo-data/ailang v0.47.2` remains distinct from the pinned execution binary `v0.41.0`; `go.mod`/`go.sum` unchanged from Round 2 (dependency closure still 54→56, only `hostcall`+`mcphttp`). The A1/A2/B/C per-milestone semantic non-vacuity and the literal whole-API-revert compile-failure boundedness established in Rounds 1/2 are retained unchanged (no production diff in this round).

## Rubric score

| Category | Score | Evidence |
|---|---:|---|
| Tests | **0/20** | Full unchanged Go profile rc1 (hard fail). |
| Lint | 10/10 | `go vet ./host/projection` rc0; `gofmt -l` clean; full-repo vet rc0 in evidence. |
| Acceptance | 24.545/30 | 18/22 criteria; D null. |
| Code quality | 0/15 | Literal rubric −5/file: `daemon.go` 966, `daemon_test.go` 1343, `projection_test.go` 1553 (all >800, inherited size; no cosmetic splitting requested). |
| Documentation | 15/15 | `docs/HOST_CHANGELOG.md` "Unreleased — candidate; acceptance pending" entry; runnable Go API / real-daemon guide payloads (equivalent host-example); design-doc header corrected to "Implemented locally; acceptance pending" (Round 2's "Planned" contradiction resolved). |
| Design fidelity | 9/10 | Thin upstream-owned wire architecture matches; full runbook execution narrower than plan wording. |

Total **58.545 → 59/100**, below the 70 bar **and** carrying a hard fail. No compiler-infrastructure or performance bonus category is triggered.

## Nonblocking findings / evidence limits

1. `host/projection/mcp.go:19` (`hostcall.New(cfg.CallbackTimeout, cfg.MaxCallbacks)`) and the `projection.go` scalar mounts are source-unchanged and correct; the shared-handler eight-slot saturation / ninth-entry-refusal-until-RETURN-then-recovery path remains a queued coverage follow-up (per verification 3n), not expanded here.
2. Nil-coordinator refusal branch (`mcp.go` "invocation coordinator is unavailable") remains correct but unpinned; queued, not a product defect.
3. Design-doc header "acceptance pending **final normal profiles**" is now slightly stale — the profiles are terminal (and the Go profile failed); wording only, no contract change.
4. `TestMCPInvokeDeadlineFreesSlot` recovery asserts `resolver.calls==2` (second resolver call is real, undelayed); the about-to-return channel alone is not proof of release — the successful second run (entered=2/return-path=2) is what demonstrates released capacity, as the directive requires.

## Blockers (concrete, with proposed fix)

- **BLOCK-1 (hard)**: unchanged full Go profile rc1. Files/lines above. These are inherited timing-bound tests (2s/3s) that overshoot only under full parallel saturation. Proposed fix is a *test-infrastructure* change outside this sprint's frozen production scope: (a) run the plain full profile with bounded serialization (e.g. `-p 1`) in CI, or (b) widen the wall-clock/commit margins in the pre-existing capsule/coordinator/daemon tests, or (c) run the full profile on a load-isolated runner. None is a production-constant relaxation; none is owned by M108. This blocker is **not** introduced by, nor fixable within, the MCP dispatch change.

## Terminal disposition

EVALUATION_RESULT: fail
EVALUATION_SCORE: 59/100 (exact 58.545)
EVALUATION_ROUND: 3 (final)
CANDIDATE: 741e3c371de3e558243e4ee77a09eeaddbac50a3

All own test subprocesses terminal. No production edits; only this evaluator report directory is untracked. The Round 3 repair is correct and verified for its target (MCP slot recovery), but the full unchanged Go profile remains red, so the hard test gate fails. Final round reached → park; no fourth judge, no waiver, no forced grade-up.

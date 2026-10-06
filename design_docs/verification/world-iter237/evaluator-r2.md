VERDICT: PASS 96/100

# Row 136 `w-mcp-effects-unrecorded-is-untyped`: evaluator round 2 (iteration 237)

- Delta judged: `18a4e9c..8b4440a`, one commit `8b4440a`. It adds a test and moves the design doc and sprint plan to `design_docs/implemented/`.
- Worktree: my own detached checkout at `/Users/voightkampff/dev/sunholo-data/.wt-world-iter237-eval-r2`. It was clean (0 tracked changes) after every mutant and has been removed.
- Evidence: `~/.ailang/state/mission-world-iter237-evidence/evaluator/r2_*`.

## Rubric (r1 → r2)

| Category | r1 | r2 | Change |
|---|---|---|---|
| Tests pass | 18 | 18 | No new failures. The r1 deduction for the flake-driven CI re-run stands. |
| Vet | 10 | 10 | rc=0 |
| Acceptance criteria | 28 | 28 | AC3.1's controller-owned records are still pending (out of scope for r2) |
| Code quality | 13 | 15 | NB1 fixed: all three untested MCP refusal sites are now pinned by a test CI runs |
| Documentation | 13 | 15 | NB3 fixed: Status reads IMPLEMENTED, and the design and plan moved together |
| Design fidelity | 10 | 10 | The delta is test-only, so this is unchanged |
| **Total** | 92 | **96** | PASS |

## BLOCKING findings

None.

## NON-BLOCKING findings

1. The new Status line hard-codes "judge PASS 92", which was the r1 score. This is cosmetic, and the controller may update it to the final score when it records the result.
2. r1 NB2 and NB4–NB6 are deferred to the controller, as instructed.

## Mutant table

The new test is `TestMCPPreDispatchRefusalLinesAreLabelled` in `host/projection/mcp_conformance_test.go`. It does not need the tool binary, so CI runs it. Every row below ran with `WORLD_TOOL_AILANG_BIN` unset.

| Mutant (applied to `host/projection/mcp.go`) | Killing subtest (`mcp_conformance_test.go:474`) | Killed? | CI-protected? |
|---|---|---|---|
| `:107` admission refusal → `logRefusal("a2a","mcp invoke","-",err)` (r1 survivor) | `/tools_call_admission` | YES | yes |
| `:76` and `:81` together → `logRefusal("a2a","mcp tools","-",err)` (r1 survivor) | `/tools_list_registry_read` and `/tools_list_unprojectable` | YES | yes |
| `:76` only (registry-read site) | `/tools_list_registry_read` | YES | yes |
| `:81` only (`mcpDescriptors` site) | `/tools_list_unprojectable` | YES | yes |
| OWN: `:107` keeps its label but drops the cause (logs `errors.New("x")`) | `/tools_call_admission` | YES | yes |

Each of the three sites is killed by its own subtest, so the test does not reach all three through one path.

## Gates (at `8b4440a`)

- `go vet ./...`: **rc=0**.
- `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race -count=1 -v ./host/projection/ ./host/daemon/`: **rc=0**, with **557 `=== RUN`**, 535 PASS, **0 `--- FAIL`**, 22 SKIP (the se-tools e2e tests, which need the tool binary).
  - The new test passes with all 3 subtests.

## Move and references

- Git records the plan as a pure rename (R100) and the design as a rename with one line changed (R098, the Status line).
- The target is flat `design_docs/implemented/`, which matches the existing layout (95 entries, flat).
- The plan's relative link `[w-mcp-effects-unrecorded-is-untyped.md](w-mcp-effects-unrecorded-is-untyped.md)` still resolves, because the two files are siblings.
- `design_docs/planned/w-mcp-effects-unrecorded-…` still appears only in the append-only charter (`world-mission.md:960`) and log (`world-mission-log.md:1327`).
- `world-mission-index.md:15` names the doc without a path.
- No scripts, code or website file refers to it.
- No broken references.

R2 JUDGED

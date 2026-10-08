# Iteration 242 — row 152 evaluator r2 (narrow: N4 docs fix)

Judge: claude-opus-5-5. Worktree `.wt-world-iter242-eval` detached at branch head `0b2d51e` (fix commit `b14916c`).

## Scope check
`git diff --stat d8e9c72 b14916c` → `docs/QUICKSTART.md | 14 +++++++++-----`, 1 file. Docs only; no code changed
since r1, so the r1 gates, mutations and broker verdict carry over unchanged.

## The new paragraph, re-derived from the code at HEAD

- **Tool ids.** `packages/se-tools/transitions.json` declares these effects:
  - `ailang-read` → Workspace.Read
  - `ailang-write` and `ailang-edit` → Workspace.Write
  - `ailang-check` → Ailang.Check
  - `ailang-run` → Ailang.Run/RunEnv/RunNet
  - `builtins-search` and `examples-search` → Ailang.Discover
  - `ailang-cli` → Ailang.CLI
  - `workspace-exec` → Workspace.Exec

  The eight listed tools are exactly the transitions whose declared effects all lie in `ailangToolEffects`
  (`workspace.go:64`), and `workspace-exec` is the only one outside it. N4(a) is fixed.
- **Exec refusal cases** (`workspace.go:338–392`). There are four:
  - The worktree fails `episodeRoot`: an empty registry, then R8.
  - `w.exec == nil`: `broker.ExecUnconfiguredHandler{}`.
  - `profileFor(ep) == nil`: `ExecRefusalHandler{ExecNoEpisodeProfileRefusal}`.
  - `NewExecHandler` fails: the `workspace-exec unavailable` line plus `execFailedHandler`.

  The paragraph names all four ("missing or unsafe", "no exec profile is configured (or none maps to the episode)",
  "building the episode's exec handler fails"). N4(b) is fixed.
- **The 3 s retry wait.**
  - `registry()` runs `ailangFamily` synchronously before `execHandler`, on every call. Failures are not cached:
    the in-flight entry is deleted on return.
  - The build runs under `context.WithTimeout(…, workspaceHandlerBudget)` with 3 s in production (`workspace.go:455`). A
    caller that joins the in-flight build waits at most the rest of that budget.
  - So each call in a persistently broken episode, `workspace-exec` included, waits up to the 3 s budget. "fix the cause
    rather than retrying" is the right advice. N4(c) is fixed.
  - Nuance, not a defect: the mkdirs and the lock check run inside that window but take no ctx. They are µs–ms, and
    "up to" holds in practice.
- The unchanged phrase "(every declared effect fails)" predates this row. In reality the transition is refused (R8) before
  its plan. Left as is: it is the existing operator vocabulary and not introduced here.

## Verdict
**N4 CLOSED.** No new finding. Score 91 → **94/100** (docs 7 → 10). PASS, zero blocking.
N1, N2, N3 and N5 from r1 stand as non-blocking.

MISSION-ROLE: planner

You are the SPRINT PLANNER for Ailang World mission iteration 372, queue row 257
`w-exec-cache-seed-at-bind`. Follow the sprint-planner skill's shape, adapted to this repo. You plan;
you do not implement, commit, push or open PRs.

## Where you work
- Worktree: /Users/voightkampff/.ailang/state/world-iter372-row257 (branch sprint/world-iter372-row257,
  HEAD 99fc663d = origin/dev dc907675 + the committed design doc). Never touch
  /Users/voightkampff/dev/sunholo-data/ailang-world (shared main checkout) or other worktrees.
- Input: design_docs/planned/w-exec-cache-seed-at-bind.md (r3 — read ALL of it, including §11; the
  [r3] text is binding: D-257-1 = A, so NO session_self.go / episode.go / episode_preset.go / warmExec
  changes; unique per-incarnation staging; boundExecHandler accessor; tests S1–S13, D1, D3, D4, C1).
- Deliverables (two files, uncommitted):
  1. design_docs/planned/w-exec-cache-seed-at-bind-sprint-plan.md
  2. design_docs/planned/sprint_w-exec-cache-seed-at-bind.json — model it on
     design_docs/planned/sprint_w-exec-failure-carries-a-reason.json (read it for the schema).
  At the end `git status --short` must list exactly those two files.

## What the plan must contain
- Milestones M1 (broker) and M2 (daemon + CLI) as the doc §7 says; split further only if a milestone
  exceeds ~300 production+test LOC. The EXECUTOR is a smaller model (claude-haiku-5-5) working from
  your plan, so the plan must be PRESCRIPTIVE: for each milestone list exact files, exact
  function/type/const names and signatures, the exact failure-line strings from the doc, each test's
  name and assertions, and the exact gate commands. Leave no design choice to the executor.
- For each test, the mutation it must kill (from doc §5) — the controller and judge will run them.
- Gate commands per milestone (narrow `go test ... -run ... -count=1`), plus the race test for D4.
  Note in the plan: the controller (not the executor) runs `./scripts/verify_go.sh` and
  `./scripts/verify_ail.sh` with `AILANG_BIN=$HOME/.pinned-ailang/ailang`.
- Name the gate `go test ./host/store/ -run '^TestProductionContextRoots$' -count=1` explicitly
  (the diff adds a production context root — charter rule).
- Acceptance criteria list (numbered) the judge will score against.
- A short risk section (from doc §6) and the #515 overlap note (one workspace.go hunk).
- Before writing, VERIFY the doc's load-bearing line references against the current tree
  (handlers_exec.go prepareCache/execSeedBudget/seedMu, workspace.go execHandler :468-506 and
  execMu :123-126, failure_test.go seed rows ~:208-218, context_roots_test.go ~:80,
  episode_report.go :692-724, newExecFixture). If any is wrong, record the correct one in the plan.
- Estimate effort honestly.

## Constraints
- Go is `go` on PATH. Use `AILANG_BIN=$HOME/.pinned-ailang/ailang` for any test that needs ailang.
  Set neither TMPDIR nor GOTMPDIR. You may run narrow read-only checks (`go build ./...`,
  `go vet ./host/broker/`) but do not write production code.
- SECURITY: never run `env`, `printenv`, `set` or anything that prints the environment.
- Every wait bounded and in-turn; never end your turn announcing a wait. Budget ~30 minutes.

## Return
Short report: the two paths, milestone list with LOC estimates, any doc line references you found
wrong, and anything left ambiguous that the controller must decide.

MISSION-ROLE: designer

You are the DESIGNER for Ailang World mission iteration 372, queue row 257 `w-exec-cache-seed-at-bind`.
You write ONE design document. You do not implement the row, commit, push, or open PRs.

## Where you work
- Worktree (your CWD for everything): /Users/voightkampff/.ailang/state/world-iter372-row257
  (branch `sprint/world-iter372-row257`, base = origin/dev `ad13d835`). Never touch
  /Users/voightkampff/dev/sunholo-data/ailang-world (the shared main checkout) or any other worktree.
- Deliverable: `design_docs/planned/w-exec-cache-seed-at-bind.md` in that worktree. You MAY also write a
  prototype patch `design_docs/planned/w-exec-cache-seed-at-bind.patch` (as 360's designer did), but any
  prototype code you apply to the tree must be reverted before you finish: at the end, `git status --short`
  must list only those one or two files.
- Raw measurement outputs: quote them in the doc's Verification Log. Scratch files go under
  /tmp/iter372-designer/ only.

## Read first (binding)
1. design_docs/world-mission.md — the queue row 257 text (grep `^257\. `), rows 255 and 256, and the
   Routing policy section. Row 257 verbatim: the GOCACHE seed is cloned on the episode's first exec call,
   inside that call's 9-second budget; `prepareCache` (host/broker/handlers_exec.go) clones the
   operator's seed (225 MB for ailang-world, `cp -cR` / `cp -a --reflink=auto`) under
   `execSeedBudget(ctx)`. Build: clone the seed when the episode's exec handler is bound, or at
   `episode start`, before the agent runs. Report a seed failure there, with the row-255 reason, and keep
   the per-call path as a check that the copy exists. Measure the clone time at bind under
   package-parallel load, and say what bound it gets. Tests: the first call spawns no seed copy when bind
   succeeded, and a bind-time failure is visible in the episode report.
2. design_docs/coding-standards.md (binding): pure core with contracts where encodable, effects at the
   host boundary, non-vacuous gates (every test must have a mutation it kills).
3. design_docs/planned/w-exec-failure-carries-a-reason.md — the row-255 design that LANDED (failure
   classes such as FailSeedCache / FailExecCache, the 200-byte failure line). Copy its doc STRUCTURE
   (header bullets: Status, Item, Designer lane, Estimated, Measurement base, Instruments; then
   1 What the row says and what you measured, 2 Mechanism, 3 Design, 4 By-construction argument,
   5 Test plan and mutation table, 6 Residuals and risks, 7 Milestones, 8 Files, 9 What this is NOT
   doing, 10 Verification Log, 11 Quorum log (leave empty for the controller)).
4. Code: host/broker/handlers_exec.go (NewExecHandler, prepareCache, execSeedBudget, seedMu),
   host/broker/exec_profile.go (Caches/Seed), host/broker/exec_startup.go and exec_probe.go (they also
   touch seeds), and wherever the daemon binds an episode's exec handler and runs `episode start`
   (find it; cite file:line).

## Constraints
- Row 256 (`w-daemon-startup-under-load`, draft PR #515) is TAKEN BY AN ATTENDED SESSION right now and
  changes: host/broker/exec_probe.go, host/broker/exec_startup.go, host/daemon/daemon.go,
  host/daemon/workspace.go, cmd/ailang-worldd/doctor*.go, docs/QUICKSTART.md. Design row 257 to AVOID
  editing those files where at all possible; if one is unavoidable, name the exact hunk and say how a
  rebase over #515 would resolve. Do not design anything row 256 owns (startup probes, registry
  bootstrap deadline).
- Fail-closed must hold: a bind-time seed failure must never let a call run against a missing/partial
  cache, and must never be silently retried forever. Say what the per-call path does after a bind
  failure (re-attempt within the call's budget? refuse with the row-255 class?) and justify it.
- Say what bound the bind-time clone gets and why (it no longer has the call's 9 s), measured.
- Measurements: Go is `go` on PATH (go1.26.x). For any ailang call use
  `AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0); never the PATH ailang (a -dirty dev build).
  Set neither TMPDIR nor GOTMPDIR (leave the rig defaults; a temp path under ~/dev trips MOD010 and long
  paths overflow the 200-byte failure line). Narrow tests: `go test ./host/broker/ -run '<regex>' -count=1`.
  Do not run the full `./scripts/verify_go.sh`; it is slow and the controller runs it.
- For the "under package-parallel load" measurement: time the real clone (`cp -cR` of a 200+ MB tree,
  e.g. a copy of your own GOCACHE subset under /tmp/iter372-designer/) idle and while
  `go test ./host/... -count=1 -p 8` runs concurrently; report N, min/median/max. Bounded: no single
  measurement run longer than 10 minutes. Delete your scratch copies afterwards.
- SECURITY (standing rule, D-WORLD-112): never run `env`, `printenv`, `set`, or anything else that
  prints the process environment, and never print a secret's value. To test a variable use
  `[ -n "$VAR" ] && echo SET || echo UNSET`.
- Every wait is bounded and ACTIVE: poll an artifact inside a bounded `date +%s` loop. A Monitor is an
  event stream, not a wait; nothing will re-invoke you when it fires, so never end your turn announcing a
  wait. Prefer foreground commands. Your total budget is about 45 minutes.
- Plain English in the doc: short sentences, no marketing words. Every load-bearing claim gets a
  Verification Log row (command + output excerpt).

## Return
End with a short report: doc path, line count, sha256 prefix, the design decision in three sentences,
the measured bind-time clone numbers, the files the build would touch (and any overlap with #515), and
anything you could not measure (say UNMEASURED, never guess).

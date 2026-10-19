# Iteration 281 executor evidence report

Evidence-only deliverable completed. Equivalence=false; M1c remains incomplete.
Worktree: /Users/voightkampff/dev/sunholo-data/.wt-world-iter281-row156-m1c
Base: a0c73053f10092fd511ad22fcf6b95d6a3f896aa.

Created only these checkout files:
- design_docs/verification/world-iter281/repro.py
- design_docs/verification/world-iter281/result.json
- design_docs/verification/world-iter281/README.md

Read CLAUDE.md, coding standards, sprint-executor SKILL.md, charter preamble/Repo Profile,
DESIGN §§1/14, ratified design §§4.2/12, sprint M1c and Environment/PF6, inherited designer
and planner reports (treated as unverified), current renderer and mirror fixture source.
Explicit evidence-only instructions superseded skill full-sprint, bookkeeping, commit,
message and generic gate actions. No production/test/design changes or checkout Git
add/commit/stash; all Git writes confined to disposable fixture repositories.

Commands (all shell commands exported /opt/homebrew/bin PATH and pinned worktree cwd):
- python3 design_docs/verification/world-iter281/repro.py: final rc0; Git 2.54.0.
  Four independent repos, ASCII selected controls true, Unicode selected false;
  exact unfiltered status and ls-files controls present in each repo.
- Recorded focused command: go test -count=1 ./host/broker/... -run
  '^TestProtectedMatchFixture$/(long-s|kelvin)$' -timeout 110s -v: rc0;
  both existing named long-s and kelvin tests PASS, expected true.
- Result JSON assertions: rc0, valid=true, equivalence=false, M1c=incomplete;
  row selection vector [true,false,true,false].
- git diff --check: rc0. git status --short: only evidence directory untracked.

Instrumentation adjustments: initial private-cache run with GOSUMDB=off refused Go
checksum verification (rc1); no result written. A direct provisioned Go binary with
private cache then measured successfully within 10s. Controller subsequently requested
portable PATH Go resolution, existing Go caches, checksum verification on, no stale
result and 120s focused-Go allowance. Final script follows that direction, captures
Go cache configuration before fake HOME, retains caches, GOPROXY=off, and bounds Git
commands to 10s, focused Go to 120s (110s test timeout). No global state seeded.
Scratch fixture repos/homes/temp are worktree-local and removed automatically. Existing
.scratch-iter281 left unchanged. Source literals were separately checked by reading
current gitPathspecs lines 327–337, never parser-dependent source introspection.

Limitations: this is successful measurement of failed equivalence, not a permanent
red acceptance test, full M1c matrix/mutations, full broker baseline, CI pass, or other
platform/engine proof. Broad broker tests were not run because PF6 remains unresolved.
No .ail changes, verify_ail, driver/skill edits, external messages, release or push.
Preserve Unicode semantics; next work is bounded designer revision + review/replanning,
not an invented needs-human-review question or a weakened fixture.

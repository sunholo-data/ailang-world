# Ailang World Mission — a semantic operating environment whose transaction language is AILANG

<!--
  Charter for the WORLD mission (from mission-charter-TEMPLATE.md, M-MISSION-PORTABILITY M3).
  DRAFTED AT BOOTSTRAP 2026-07-23 — written to be CHALLENGED, not to be final.
  Iteration 0 = ratify this charter (bar + queue + guardrails) through the design quorum
  with Mark, attended. The kill switch stays SET until ratification.
  Thesis document: DESIGN.md (v0.3) in this directory — this charter is its operational
  distillation; where they conflict, the ratified charter governs the loop and the
  conflict gets fixed in DESIGN.md.
-->

**Type**: Long-running mission (peer of the V1 mission in `sunholo-data/ailang`); advanced by a
scheduled outer loop on the always-on rig. **First mission outside the AILANG repo.**
**North star**: ship **Ailang World 1.0** — a local-first semantic operating environment
([DESIGN.md](DESIGN.md) §1: an immutable, typed, content-addressed world graph whose transactions
are AILANG programs; propose → verify → commit; capabilities, budgets, evidence, replay) — with
the software-engineering domain live — value proven by clause 5's provenance teeth + the R1
standing evidence, floored by clause 4's resident-agent non-inferiority check (motoko remains
the aspirational first *native* agent as an optional arm).
**Traces to**: [DESIGN.md](DESIGN.md) (thesis) and the AILANG program's operating model
([PROGRAM.md](https://github.com/sunholo-data/ailang/blob/dev/design_docs/PROGRAM.md) — frozen
minimal core, extension-routed evolution); language gaps found here route BACK to
`sunholo-data/ailang` as issues/backlog items, never worked around locally.
**Skill**: the SAME unforked `mission-control` skill every mission uses, reached via user-level
symlinks (`~/.claude/skills/*`) — this repo carries NO skill copies (a copy is a fork; forks stop
learning).
**Scheduling**: launchd `dev.ailang.mission-world` — **ARMED by Mark 2026-07-23** ahead of
ratification; sprint-routing stays BLOCKED until the charter is ratified (iter-0 honored this:
quorum-only, no sprint). Off switch: `~/.ailang/state/mission-world.disabled`. StartInterval
staggered vs the V1 loop (shared rig quota). Billing guard: subscription-or-nothing.
**Log**: [world-mission-log.md](world-mission-log.md) — append-only, one entry per iteration.
**Human-facing reporting**: GitHub issue #1 — every iteration posts its report there as a comment
(Mark follows by email); driver crashes post there too.

## Repo Profile (M-MISSION-PORTABILITY M2 — the per-mission values mission-control reads)

- **Repo slug**: `sunholo-data/ailang-world` (driver: `MISSION_REPO`)
- **Mission doc**: `design_docs/world-mission.md` (driver: `MISSION_DOC`)
- **Mission name / state namespace**: `world` (driver: `MISSION_NAME`; fully namespaced
  `~/.ailang/state/mission-world-*` paths — no collision with the V1 loop)
- **Bookkeeping issue**: `#1`, rotates weekly; live number in
  `~/.ailang/state/mission-world-gh-issue` (seeded by the v1 agent after ratification)
- **CI workflows Gate 3b / Gate 1 poll**: `CI` (job: "ailang-code verify gate")
- **Verify profile**: `ailang-code` — the shipped `ailang` binary IS the gate:
  `ailang check` (types), `ailang test` (tests, once test suites exist),
  `ailang ai-check` (unified check+verify: types + Z3 in one JSON — do not reinvent a split
  gate). Repo-wide sweep: `./scripts/verify_ail.sh` (what CI runs; fails loudly on zero modules).
  Go host code (daemon/store, from M1 onward) adds `go build ./... && go test ./...` to the gate
  when it first lands — extend the CI workflow in the same PR that introduces Go code.
- **Rig PATH (process fix, iter-18 — 3 instances in ONE iteration)**: the agent tool shell's `PATH`
  does **not** include `/opt/homebrew/bin`, so `gh`, `go`, `node` and `codex` all fail with
  `command not found` / `env: node: No such file or directory` (`rc=127`). Every Bash call that
  needs them must `export PATH=/opt/homebrew/bin:$PATH` first, and **every directive handed to a
  spawned sub-agent must say so** — otherwise a planner/executor reads it as a broken toolchain or
  a spent codex quota rather than a PATH gap (iter-18 lost a codex probe to exactly that
  misreading). The pinned AILANG binary lives OUTSIDE the repo at `/tmp/ailang-v0300/ailang`
  (v0.30.0, commit `e37b370`, clean) — always pass it explicitly as `AILANG_BIN`.
- **THE RIG'S SHELL IS `zsh`, AND TWO OF THE SKILL'S OWN INSTRUMENT IDIOMS SILENTLY PRODUCE
  CLEAN-LOOKING NON-RESULTS HERE (process fix, iter-37; 2 instances in ONE iteration, both in the
  controller's own Gate-2 verification).** Every Bash call in this loop runs under `zsh`, not bash,
  and both defects below return something that looks exactly like a passed check:
  (a) **an unquoted `--include=*.go` is GLOB-EXPANDED by zsh** — with no matching file in the cwd it
  fails the whole command with `no matches found`, so `grep -rn PATTERN --include=*.go .` **never
  runs** and the surrounding pipeline reports `0` hits. Iter-37 nearly handed its executor a
  fabricated *"zero callers anywhere"* fact this way; the only reason it was caught is that the same
  call carried a **known-positive control** which also came back empty. Quote it (`--include='*.go'`)
  or, better, use the `Grep` tool, which cannot be glob-mangled.
  (b) **`${PIPESTATUS[0]}` IS BASH SYNTAX AND EXPANDS TO EMPTY IN `zsh`** — the array is
  `pipestatus` and it is **1-indexed**, so the skill's own Gate-2 remedy for *"exit codes through
  pipes lie"* prints `rc=` and voids the reading it was added to protect. Use direct invocation
  (`cmd > /tmp/out 2>&1; echo "rc=$?"`) and read the file afterwards; that works in both shells.
  **The general rule both instances teach: a remedy is an instrument too, and inherits the same
  burden of proof as the thing it verifies.** Pair every negative or exit-code check with a
  known-positive control in the SAME call, so a broken instrument is distinguishable from a real
  all-clear. **Routed upstream as a PROPOSED shared-skill fix** (World cannot edit the
  mission-control SKILL.md — it lives in the V1 checkout), since Gate 2's step 3 currently
  *prescribes* the bash-only form.
- **`git checkout <file>` IS NOT A MUTATION REVERT IN THIS MISSION — IT IS A MILESTONE DELETION
  (process fix, iter-38).** Every sprint here is delivered as an **UNCOMMITTED worktree diff**,
  because the codex executor's `workspace-write` sandbox cannot reach the linked worktree's `.git`.
  So the obvious revert idiom restores from **HEAD**, not from the pre-mutation baseline: it throws
  away the milestone's version of the file and leaves the suite **green on `origin/dev`'s code** —
  a passing run that proves nothing about the work under review, with no error, no diff noise and
  no hint in the output. Measured live at iter-38 while the controller was reproducing
  `MUT-INTENT-AFTER-DISPATCH`: pre-mutation `sha256=99f3568…`, post-`git checkout` `sha256=6bc8514…`.
  **The rules.** (a) Before mutating, take a copy: `cp f /tmp/f.bak`; revert with `cp /tmp/f.bak f`,
  never `git checkout` / `git restore` / `git stash`. (b) The **sha256 comparison is not
  bookkeeping — it is the only thing that distinguishes a revert from a deletion**, so it is
  mandatory and both values must be printed. (c) If they differ, the file is recoverable from the
  executor's own logged `git diff` (codex prints full diffs), and the restoration must be proven by
  **two independent hashes** — `git hash-object` matching the diff's recorded post-image blob id,
  and sha256 matching the pre-mutation baseline. (d) This generalises the iter-35
  mutation-as-instrument rule to its **other end**: a mutation's validity covers *applying* it and
  *undoing* it, and the undo is the half with the destructive failure mode. **Corroborating second
  instance, same iteration, same family**: a `grep -c -- "--- SKIP"` over NON-verbose `go test`
  output reported "ZERO skips" — `--- SKIP` lines print only under `-v`, so the zero was structural
  and could never have been anything else (truth under `-v`: 2 benign subprocess-helper skips).
  **A verification step is code, and code that cannot fail is not a check.**
- **Codex model pins are NOT covered by the driver's probe (process fix, iter-19).** The driver's
  codex pre-flight (`tools/launchd/mission-control.sh:248`) runs `codex exec --skip-git-repo-check
  'reply with exactly: ok'` — **without `--model`** — so it exercises codex's DEFAULT model and
  reports the lane healthy even when the pinned model cannot run. Iter-19 hit exactly this:
  `MISSION_EXECUTOR_MODEL=codex:gpt-5.6-sol` failed with `400 … 'gpt-5.6-sol' requires a newer
  version of Codex` on codex-cli **0.137.0**, while a probe without `--model` returned rc=0 on
  default **`gpt-5.5`**. So: **always run the skill's own probe WITH `--model` before trusting the
  lane**, and read a model-availability 400 as a PIN problem, never as spent quota (that
  misdiagnosis is the iter-18 scar repeating in a new costume). The driver is FROZEN CORE — do not
  patch it here; route the fix upstream on both channels. Also guard the call site with
  `env -u OPENAI_API_KEY` so an ambient metered key cannot silently bill while `auth.json` is on
  ChatGPT subscription auth (`auth_mode=chatgpt`) — the `claude-sub` discipline, applied to OpenAI.
- **ENUMERATE the evaluator's non-blocking carry-forwards in the log entry (process fix, iter-19).**
  Iteration 18 recorded "three non-blocking carry-forwards → A2" as a bare COUNT; the items
  themselves were never written down anywhere durable and are unrecoverable — A2 could not fold
  what it could not read. Every Gate-4 entry now lists them **numbered, one line each, with the
  checkpoint they target** (iter-19: CF-1…CF-9), and the evaluator directive must ask for that list
  explicitly.
- **A fresh worktree NEVER contains the sprint plan JSON (process fix, iter-20).** `.gitignore`
  line 3 is `**/.ailang/`, so `.ailang/state/sprints/<item>.plan.json` — the executor's OWN plan —
  is structurally absent from every `git worktree add`. Iter-20's codex executor reported it as a
  mid-run blocker and fell back to the design doc (correct: the doc governs, so no harm — but the
  plan's `non_vacuity_requirements` would have been silently lost had the directive not quoted
  them). So every worktree-based executor directive must EITHER `cp` the plan into the worktree
  first (it stays gitignored, so it never pollutes the diff) OR quote the binding requirements
  inline. Do not assume a sub-agent can read mission state that lives outside git.
- **A MUTATION IS AN INSTRUMENT — ESTABLISH ITS VALIDITY BEFORE ITS RESULT COUNTS (process fix,
  iter-35; ≥2 instances across iter-34 and iter-35).** Non-vacuity mutations are how this mission
  proves a gate has teeth, so a broken mutation does not merely waste a run — it manufactures
  evidence. **Three failure modes are now on record, and every one produces a reading
  indistinguishable from a real result:**
  (a) **NEVER APPLIED** (iter-35): a `str.replace` pattern carrying four tabs against three-tab
  source matched nothing, wrote the file back unchanged, and the suite went all-green — identical
  to *"the mutation was applied and the gate failed to catch it"*. The conclusion it supported
  happened to be true, which is worse, not better: **being right by luck is not a method.**
  (b) **FAILS TO COMPILE** (iter-34, three times in one iteration): the package does not build, the
  gate script reports failure, and a build error wears the gate's clothes. **You measured the
  compiler, not the gate.**
  (c) **THE NAME DENOTES A FAMILY** (iter-34 `MUT-BENCH-DROP` delete-vs-rename; iter-35
  `MUT-PENDING-UNBOUNDED`, where two forms both compile and both red the same test by *different*
  mechanisms). A mutation report that does not say WHICH FORM ran is not checkable, and a second
  party running the other form will reasonably conclude the record is wrong — that is exactly what
  happened at iter-35, costing a judge round-trip to resolve.
  **Required practice, all five steps:** (1) apply under an assertion that the pattern matched
  **exactly once**, and abort loudly otherwise — never a bare `replace`/`sed`; (2) print the
  applied `git diff` BEFORE running the suite; (3) confirm the mutant **COMPILES** and say so;
  (4) **NAME THE FORM** in the record, alongside the exact error text it produced; (5) revert by
  `cp` from a backup taken first — **never `git checkout --` on a file carrying uncommitted work**
  (iter-34 destroyed an executor's work that way) — then verify byte-identity by `sha256`.
  This is the same family as the silent z3 skip (V27) and the silent `t.Skip` (B1): *a check that
  did not run is indistinguishable from a check that passed*. **PROPOSED UPSTREAM** to the shared
  `mission-control` skill's Gate-2 verification protocol (it lives in the V1 checkout, so World
  proposes and never applies) — V1 runs mutations too and has the same exposure.
- **The `verify_ail.sh` module count is now 10, not 9 (process fix, iter-25).** Iteration STATUS
  stamps 22/23/24 all recorded the gate as "EXACTLY 4/4 identities / **9** modules / 14 tests" as a
  first-party freshness check. `design_docs/sketches/storejournal.ail` landed iter-25, so the sweep
  reads **10 module(s)** from `615619c`+1 forward. The load-bearing numbers are UNCHANGED and stay
  the thing to assert: **4/4 required `world/` identities** and **14/14 named `world/` tests** —
  both are `world/`-scoped by construction (`EXACT_TOTAL_VERIFIED` sums only `case "$mod" in
  world/*`, and Leg 2 runs `ailang test --format json world/`). A future iteration seeing 10 where
  it expected 9 is observing this commit, not a regression. **Corollary worth knowing before you
  design with sketches**: Leg 2 executes `ailang test` on `world/` ONLY, so a sketch's inline tests
  are **never CI-run** — the sketch's *contracts* are swept by Leg 1's per-module `ai-check`, but
  its `tests[]` are not. Any milestone relying on sketch tests must run them explicitly in its
  verify_commands (found independently by the controller and the designer, iter-25).
- **When CI reds on a TIMING assertion, WIDEN THE GAP — never loosen the guard (process fix,
  iter-33).** Iteration 33's CI red was `elapsed <= 200ms` for a 40ms bound, where a shared runner
  had spent 316ms on process spawn alone. The obvious fix is to raise the guard past the
  observation. That would have shipped a real bug: the fixture's sleep was raised 0.3s → 5s and the
  guard 200ms → 2s instead, taking the honoured-vs-ignored margin from ~2× to ~130×, and the next CI
  run reported **5.00s — the mutation signature** — which is how the actual defect was found.
  **A timing assertion exists to detect a bound that is not enforced; loosening it past the symptom
  deletes exactly that power.** Widen the separation between the two outcomes, then re-run the named
  mutation to confirm it still reds (iter-33 re-proved `MUT-HANDLER-TIMEOUT-IGNORED` at 5.25s/5.29s
  against the new 2s guard *before* pushing). **Corollaries, both paid for in the same iteration:**
  (a) **an exec surface must kill the process GROUP, not the child** — `exec.CommandContext` reaches
  only the direct child, so a forked grandchild keeps the inherited pipes open and the deadline is
  not enforced; `scripts/verify_ail.sh` has carried this correction since **V26** and the Go code
  repeated it, so audit any NEW `exec.Command*` against `Setpgid` + `Kill(-pid)`; (b) **darwin and
  linux disagree here**, so the rig's green is not evidence about CI for anything subprocess-lifetime
  shaped — a gate that passes locally and reds in CI is data, not noise.
- **`passed_tests` is NOT the named-test count (process fix, iter-25 — the landed correction D-B,
  restated because a fresh author hit it again).** `ailang test --format json` reports
  `passed_tests`/`total_tests` that ALSO count contract-derived properties: for
  `storejournal.ail` that is `passed_tests: 37` against **`len(tests[]) == 30`** named tests
  (32 / 25 before LAW 6's iter-29 widening — the example moves whenever the sketch does, which is
  itself why the two numbers must always be re-measured rather than quoted). The
  iter-25 designer wrote "26/26 named tests" from `passed_tests` when the real named count was 20;
  the controller caught it by re-running the command. Always report the two numbers separately, and
  gate on `len(tests[])`.
- **The weekly-rotation rule misfires on a thread created just before Monday 07:00 (process fix,
  iter-20 — 2nd instance).** Issue #9 was created `2026-07-27T05:51Z`, ~1h BEFORE that Monday's
  07:00 boundary, so the literal rule ("past the most recent Monday 07:00 AND the current issue was
  created before that boundary") reads ROTATE — which would open a second thread for the very week
  #9 already titles. Iter-19 and iter-20 both hit this and both judged NOT-DUE. **The intent test
  governs**: a thread whose TITLE names the current week IS this week's thread regardless of the
  clock; the `>80 comments` limb is the real size bound. Rotate on the boundary only when the
  current thread belongs to a PREVIOUS week.
- **`/tmp` scratch paths in the cross-provider recipe are SHARED between missions — namespace them
  (process fix, iter-21).** The mission-control skill's Gate-3 codex recipe names literal paths
  (`/tmp/codex_directive.txt`, `/tmp/codex_run.sh`, `/tmp/codex_last.txt`), and BOTH loops run on
  this one rig. Iter-21 went to write `/tmp/codex_run.sh` and found it holding the **V1 loop's**
  directive (a `design-subsumption` worktree), i.e. the two missions silently share one scratch
  buffer: a mid-flight collision could hand one mission's executor the other's directive, or
  truncate a running capture. The blast radius is real even though these are regenerated per use.
  So every `/tmp` artifact this loop writes gets a mission suffix — `codex_directive_worldb.txt`,
  `codex_run_worldb.sh`, `codex_last_worldb.txt`, `codex_out_worldb.log`. Treat the skill's literal
  paths as TEMPLATES, never as filenames. (The skill itself lives in the V1 checkout and World
  never edits it — this is proposed upstream, applied locally as a Repo Profile rule.)
- **A driver probe failing `rc=127` is a PATH gap — NEVER a spent quota and never an unusable model
  pin (process fix, iter-21; 4th instance of the PATH class, 3rd misdiagnosis costume).** Iter-21's
  fire exported `executor=opus` because the driver's codex pre-flight returned
  `rc=127 exec: codex: not found`: `tools/launchd/mission-control.sh:44` exports
  `PATH="$HOME/go/bin:$HOME/.local/bin:$PATH"` and omits `/opt/homebrew/bin`, so `claude` resolves
  and **`codex` cannot**. The skill's own probe WITH `--model` and PATH fixed returned rc=0 on
  codex-cli 0.145.0. Rule: **before honouring ANY driver-exported model fallback, re-probe** (the
  iter-19 rule, now load-bearing in the opposite direction — it catches false NEGATIVES too), and
  classify the rc: 127 = environment defect, fix or route it and honour the ratified pin; a 400 =
  pin problem; a `QUOTA_SIG` match = quota. Routed as `sunholo-data/ailang#493`.
- **CORRECTION (iter-23): the rc=127 cause above was MIS-ATTRIBUTED, and the defect was OURS —
  `PATH` now fixed in `~/.config/ailang/mission-world.env`.** Iter-21 blamed the shared driver
  (`tools/launchd/mission-control.sh:44`) and asserted "the V1 loop is demoting to opus every fire
  too". Mission-v1 refuted that from inside its own fire (its codex pin survived the probe), and
  iter-23 confirmed the real cause first-party: **the World launchd plist sets NO
  `EnvironmentVariables` PATH at all** (`grep -c PATH
  ~/Library/LaunchAgents/dev.ailang.mission-world.plist` → **0**), while V1's supplies
  `/opt/homebrew/bin`. Line 44 **prepends** — it is correct-but-dependent, working for any mission
  whose plist gives it a usable base and failing for any mission that does not. So `codex`, `gh`,
  `go` and `node` were all unreachable on this mission only.
  **The fix needs no frozen-core edit and no launchd reload**: the driver sources
  `~/.config/ailang/mission-<name>.env` at **line 48** — after line 44's export, before the codex
  pre-flight at line 297 — so a bare `PATH=/opt/homebrew/bin:$PATH` in that file lands in exactly
  the right window. Applied iter-23 and verified by replaying the driver's own ordering under
  `env -i` (codex/gh/go all resolve). The per-iteration `export PATH=...` rule above stays as the
  belt: it costs nothing and covers a fire that starts before the env file is read.
  **The durable lesson is the mis-attribution, not the path.** Four iterations (18/19/21/22) each
  re-derived this symptom from scratch and each landed on a DIFFERENT wrong cause — spent quota,
  then an unusable model pin, then the shared driver — because a silent demotion of a ratified pin
  presents identically in all three cases. Two rules follow: (a) when the same symptom produces a
  third distinct diagnosis, the fault is more likely in the part of the system you have NOT
  inspected than in the part you keep re-reading — here, nobody had ever looked at our own plist;
  (b) **before blaming shared/frozen infrastructure, check whether a peer consumer of that same
  infrastructure is healthy** — one `grep` on V1's plist would have refuted the iter-21 filing
  before it was written. A defect in shared code should reproduce for every consumer; one that
  does not is a local environment defect wearing a shared-code costume.
- **A LIVENESS OR EXIT-CODE CHECK IS ONLY EVIDENCE IF IT REFERS TO THE PROCESS YOU LAUNCHED
  (process fix, iter-24 — TWO instances in ONE iteration).** The skill's cross-provider recipe
  captures a pid (`pid=$!`) and polls `kill -0 "$pid"`, then reads `wait "$pid"`. Both times I
  deviated from that snippet, and both deviations produced a reading that was not about the thing
  I meant to measure:
  1. **`wait` on a non-child pid returns rc=127.** I launched the codex probe as
     `( codex … & echo $! > file )` — a *nested subshell* — then ran `wait "$pid"` from the outer
     shell, where that pid is not a child. `wait` reported **rc=127**, i.e. "no such job", while the
     probe had in fact replied `ok` and the lane was healthy. In THIS mission `rc=127` has already
     produced four wrong diagnoses (iters 18/19/21/22 — spent quota, unusable pin, frozen driver);
     a *fifth* rc=127 that means "your `wait` had nothing to wait for" is exactly how that scar
     re-opens. **Only the probe's own OUTPUT saved it.**
  2. **`pgrep -f "<pattern>"` self-matches the polling shell.** A completion poll written as
     `until ! pgrep -f "codex exec --model"; do …` matches its OWN shell's command line, which
     contains that string — so "is it still running?" is permanently TRUE. It read as a 30-minute
     hang on a run that had already exited `rc=0`, and the loop I wrote carried **no deadline at
     all**, which is my own Standing-rule-6 violation.
  **The rules**: (a) poll the **captured pid** (`kill -0 "$pid"`), never a name pattern; (b) keep
  the launch and the `wait` in the **same shell** — a nested `( … & )` breaks the parent/child
  relationship that `wait` needs; (c) read a probe's **output**, not only its exit code, before
  concluding anything about a lane; (d) every loop carries a `date +%s` deadline, including the ones
  you write ad hoc while waiting. **The meta-rule, which is the iteration-107 lesson in a new
  costume: when the skill ships a snippet, use it verbatim — a hand-rolled variant is a new defect
  surface, and a broken instrument reads exactly like a real measurement.**
- **STATUS rotation is a HAND edit, never a regex script (iter-18 scar).** A "move the 4th stamp to
  the archive" script that bounded the stamp with the next `---` deleted **293 lines** of this
  charter (bar, Conflict Surface, guardrails, routing policy, queue, Premise Verification Log) —
  the stamps are not `---`-delimited. Caught by `git diff --stat` before commit and restored with
  `git checkout --`. Always `git diff --stat` the charter before committing it.
- **AN INHERITED GATE IS A CLAIM — audit the gates you did NOT write (process fix, iter-32; THIRD
  instance of the self-referential-gate class).** The non-vacuity rule this mission already runs —
  *a named RED mutation is evidence only if it mutates the code the gate guards* — is applied to
  the mutations a milestone **chooses for itself**. Three times now the defect has instead been in
  a gate **inherited from an already-passed milestone**, where nobody re-asks the question because
  the gate is green and its milestone shipped: iter-29 (AC6 was owned by no milestone, so
  `MUT-SPLIT-TX` guarded a test nobody owned), iter-30 (SD.C's `MUT-AUTO-RETRY` edits
  `recover_test.go`'s own helper → CF-H-1), iter-32 (M3.A's drift test mirrored `recordConsistent`
  in a **test-local** `sketchRecordConsistent`, so it proved the TEST matched the sketch and never
  that PRODUCTION did — forcing production `RecordConsistent` to `return true` red **1** subtest
  and **ZERO** drift rows). The third one is the sharpest: it sat inside a milestone a judge had
  passed **84/100 with zero blocking findings**, which is precisely why nobody looked.
  **The rule**: when a milestone extends a surface a previous milestone gated, run the question
  *"what would have to break for this gate to red?"* against the **INHERITED** gates too, not only
  the new ones — and answer it with a mutation, not by reading. **The cheap universal instrument**
  is to force the production predicate the gate names to a constant (`return true`) and count what
  reds; a gate over a test-local mirror reds nothing and is exposed in one command. **Run the
  identical instrument before and after any rewire** — changing the instrument between readings
  measures the instrument (the iter-108 lesson), whereas one unchanged instrument turns "I fixed
  it" into a number. This is the planner's and the judge's job as much as the executor's; iter-32's
  instance was found by the **planner**, at plan time, before a line was written.
- **A MUTATION THAT REDS BY FAILING TO COMPILE PROVES NOTHING — and the revert must not be able to
  destroy the subject (process fix, iter-34; THREE instances of the first limb in ONE iteration).**
  The mission's non-vacuity discipline asks whether a named mutation *reds*. It never asked *how*,
  and "the build broke" wears a red gate's clothes exactly. Iter-34's three: (a) my own
  `MUT-REPLAY-SKIP-VERIFY` deleted an assertion block and left `expectedBudget` unused —
  `declared and not used`, which I nearly recorded as the gate having teeth; the correct form keeps
  the variable live and then reproduced the executor's exact failure message. (b) **The sprint
  plan's own `MUT-BENCH-DROP` is uninformative for one of its two names**: deleting
  `BenchmarkBrokerFSRead` leaves `"os"` unused, so `bench_worldd.sh --smoke` reports
  `underlying go test failed` rather than `missing expected benchmark(s)` — the manifest gate never
  fires, so the mutation says nothing about the gate it names. (c) It was **structurally invisible
  to the executor**, because under `--sandbox workspace-write` the loopback denial makes BOTH names
  read identically. **The rules**: (i) after any mutation, confirm the package still COMPILES before
  reading the test result — a mutation must change behaviour, not buildability; (ii) prefer a
  mutation form that cannot break the build — **renaming** a benchmark out of the manifest isolates
  the gate where **deleting** it cannot; (iii) read the gate's own MESSAGE, not just its exit code:
  `missing expected benchmark(s): X` and `underlying go test failed` are different findings.
  **Second limb, paid for in the same iteration**: I reverted a mutation with
  `git checkout -- <file>` while that file carried **uncommitted** executor work, and destroyed it;
  the next mutation's reading was then pure artifact. Reconstructing it happened to come back
  **byte-identical to the executor's own reported sha256**, which is the only reason the recovery is
  evidence rather than a plausible story. **Revert a mutation from a `cp` backup taken immediately
  before it, never with `git checkout`, unless the file is known-committed.** The instrument must
  not be able to delete the thing it measures — the same principle as "changing the instrument
  between readings measures the instrument" (iter-108), one step further.
- **THE JUDGE CAN REFUTE THE CONTROLLER'S OWN HEADLINE FINDING — and when it does, retract where you
  published (process fix, iter-34; the mirror of the skill's iter-105/111 provenance rules).** Those
  rules point one way: a sub-agent's or judge's claim is a claim until the controller reproduces it.
  Iter-34 ran the other way. I measured `skipped_tests: 5` first-party, confirmed it pre-existing
  with one unchanged instrument across four commits, filed it as the **third instance** of the
  V27/B1 silent-skip class, wrote it into a design doc and a PR body, and **published it upstream as
  a public issue**. The judge refuted the framing from this repository's own premise rows —
  `implemented/w-m1-ailang-hardening.md:103` (**V14**, *"expected noise"*), the same doc's `:378`
  and `:460` (*assert on named `tests[]` and `failed_tests` only, **never `skipped_tests`***), and
  the broker doc's own **V5**, which records `skipped_tests: 5` verbatim. The measurement was right;
  the CLAIM built on it — *nobody knew* — was false, and the repo said so in three places.
  **The rules**: (a) before calling anything an Nth instance of a known class, `grep` the repo's own
  premise/verification rows for the phenomenon — a documented-and-deliberate exclusion is not a
  silent failure, and the difference is the whole finding; (b) a first-party measurement earns trust
  for the OBSERVATION, never for the narrative wrapped around it (iter-25's lesson, roles unchanged);
  (c) **retract in every place you published, at the same volume** — a PR *comment* not a quiet body
  edit, a commit whose subject line says "retract", and a correcting comment on any public upstream
  issue; (d) a retraction is a SUCCESS of the loop and is recorded as the iteration's finding, not
  buried as an erratum. What survived — the number, the 5-of-5 proportion, and that no CLAIM
  aggregates it — is worth more stated at its real size than the overclaim was.
- **A CLAIM THAT CITES A GATE IS TWO CLAIMS, AND THE SECOND ONE IS ALMOST NEVER CHECKED (process
  fix, iter-36; ≥2 instances — iter-31 and iter-36, both at Gate 2, both about `host/store`).**
  When a design doc, a sprint plan or a queue row costs an item as cheap *because* some gate stays
  green, it is asserting (1) the change is safe and (2) **that gate can tell**. This mission has a
  strong habit of auditing (1) and none at all of auditing (2). Iter-31: the ratification's plain
  reading was not executable against the substrate that had been built to satisfy it — found only
  by measuring at plan time. Iter-36: the queue row's *"no schema change … so AC13's `sqlite_master`
  gate stays green"* was doubly wrong — the DDL **does** constrain the change (`CHECK (kind IN
  ('intent','outcome'))`, `UNIQUE (invocation_id, kind)`), and the cited gate is **inert** on the
  table in question. `plan.json` called it a *"Planner's **measured** note"*, which is how a
  half-measurement acquires the authority of a whole one.
  **The rule**: before routing any item whose cost, scope or safety rests on a named gate, run the
  cheap universal instrument on **that gate** — mutate the thing the gate claims to protect and
  watch what happens. Two outcomes, both worth the two minutes: the mutation reds (the citation is
  sound) or it survives (the citation is decoration).
  **The NEW limb this iteration adds — a gate can have teeth from one mechanism while the mechanism
  it is DOCUMENTED as having is dead.** `MUT-DDL-DRIFT`, the gate's own named mutation, reds
  honestly — via a **sha256 source pin** that `t.Fatalf`s before a database is ever built, so the
  `sqlite_master` comparison the doc credits never executes; short-circuiting that comparison
  outright (`MUT-DDL-COMPARE-DEAD`, still compiling) leaves the package green. So "the named
  mutation reds" is **not sufficient** to establish that the named MECHANISM works. Extend iter-34's
  *read the gate's own MESSAGE, not its exit code* one step: read the message and confirm it names
  **the mechanism you are relying on**. A gate with real teeth in the wrong place is more dangerous
  than a gate with none, because its green gets cited as evidence.
  **Corollary — say which part of the protection is real.** The sha256 pin genuinely protects
  pre-existing tables; only the `sqlite_master` limb is dead, and only the journal table and the
  upgrade path are uncovered. Reporting "the gate is broken" would have been an overclaim in the
  iter-34 sense, and reporting "the gate is green" an understatement in the iter-35 sense.
  **PROPOSED UPSTREAM** to the shared `mission-control` skill's Gate-2 verification protocol (it
  lives in the V1 checkout, so World proposes and never applies) — V1 routes cited-gate claims too.
- **RUNNING THE QUORUM: A BIGGER DOC BUYS A THINNER REVIEW (process fix, iter-36).** `ailang
  design-quorum` prices each reviewer on the doc's input tokens against `--max-cost-usd`, so
  **enriching a doc can silently cost it a reviewer**. Iter-36 added five controller-measured
  premise rows and round 1 then refused `gpt5-6-sol` pre-flight — *"estimated cost $0.1005 (doc
  ~13952 input tok) exceeds cap $0.1000"*, zero spend — and the reviewer thereby lost is the one
  that later found the round-2 TOCTOU. The tool behaved correctly (degrade to N−1, absentee named
  with its reason, never a silent pass); the defect was the operator's default. **The rule**: read
  `synthesis.absent_reviewers` on EVERY quorum result before acting on the verdict, and if anyone
  is absent for `budget`, re-run with a raised `--max-cost-usd` rather than accept a one-eyed
  round — the mission's metered ceiling is `$5` and a full round costs cents (iter-36 total:
  **$0.164**). A quorum that lost a reviewer to arithmetic is not the quorum the charter specifies.
- **Gate 4 skipped its own STATUS stamp once and nobody noticed (process fix, iter-34).**
  Iteration 33 updated the queue row and the log but wrote **no STATUS stamp**, so the charter's
  STATUS block still showed iter-32 as newest and the charter *alone* would have said M3.B had not
  landed. Gate 4 is append-only bookkeeping with three separate artifacts (log entry, STATUS stamp
  with rotation, queue tag) and nothing cross-checks them. **Before committing Gate 4, assert all
  three: the newest `## STATUS` line names THIS iteration, the log's last `## Iteration` heading
  names it, and the queue tag for the picked item names the landed milestone.** Cheap, and it is
  the third artifact-drift instance in this mission after the CF-I-2→CF-J-2 rename and the stale
  Appendix-A JSON.

---

## STATUS (rotation rule)

Newest **3** STATUS stamps live here; older ones move to `world-mission-status-archive.md`.
At Gate 4, after adding your stamp, move the now-4th stamp to the TOP of the archive file.

## STATUS 2026-07-29 (iteration 38) — **`w-effect-journal` MJ.B LANDED — ONLY MJ.C REMAINS** (PR #27 → squash `3ef5510`, dev CI green both jobs SHA-addressed on the merge commit **and the step logs read to prove the gates ran**; judge sonnet **PASS 86/100, ZERO blocking**; executor `codex:gpt-5.6-sol` on codex-cli 0.145.0, `auth_mode=chatgpt` confirmed two ways; planner **not fired** — the iter-37 plan already covered MJ.B; `metered=$0.00`). The broker pipeline is now wired onto the effect journal: `Invoke`'s allowed arm is **`PutObject(request) → AppendNextEffectIntent → debit → dispatch → result object → putRecord → AppendEffectOutcome`**, the denied arm journals NOTHING, the failed arm journals a resolved `"failed"` outcome with the real record ref, replay journals nothing (**asserted by a test, not by reading the code**), and `Recover` gains a second bounded page walk over pending EFFECT intents that **REPORTS ONLY** — never dispatching, resolving, appending an outcome or re-executing, with the counting probe asserted at ZERO. Closes **AC7, AC9, AC10**; +547/−48 vs the planner's +585 = **0.94×**, the first milestone of this item to land UNDER estimate. Crash-window proof is **fault-injected, not simulated**, and the r1-quorum-mandated resumption arm mints a strictly later ordinal and walks it to `resolved`. **THREE PLANNER FACTS REFUTED AT GATE 2, BEFORE A TOKEN REACHED THE EXECUTOR, ALL CHEAP TO MEASURE.** (i) **PD1's call-site census was wrong in both directions and omitted a whole constructor**: measured with a tool that cannot be zsh-glob-mangled, `NewSession` has **10** call sites (not 15), `NewReplaySession` **6** (not 5), and the unexported **`newSession` — 10 further sites, 2 of them PRODUCTION — is absent from the census entirely**, though it takes the same positional args and must also carry the episode ID. (ii) **OQ-5's premise was FALSE, and OQ-5 is a STOP instruction**: it justifies its default by *"the caller-supplied logical clock the broker **already threads for records**"*, but `LogicalTime` appears **zero times** in `host/broker` production code and `EffectRecord` has no such field — unexamined, the executor would have been CORRECT to stop and the milestone would have died at checkpoint one. The conclusion survives its false premise: `EffectRequest.Now` is caller-supplied, in scope in `Invoke`, already feeds `requestHash` and capability expiry, and `time.Now` appears nowhere in `host/broker` production code while `EffectRequest` is constructed **only** in `_test.go` files — so `req.Now` is a logical value by construction. Resolved and handed down **with its evidence**. (iii) **The plan's own verify command prescribes a known-red gate** (`-race` over `host/store`). **THE ITERATION'S SPINE: `git checkout <file>` IS NOT A MUTATION REVERT — IT IS A MILESTONE DELETION WEARING ONE.** Reproducing the headline gate myself (`MUT-INTENT-AFTER-DISPATCH`, **2 red / 99 green**, message `state "not-started" … want indeterminate`, matching the executor exactly), the revert step `git checkout host/broker/broker.go` restored from **HEAD** — and every sprint in this mission is delivered as an **UNCOMMITTED** worktree diff, because the codex sandbox cannot commit. So it silently discarded the milestone's entire `broker.go` and left a **green suite running on `origin/dev`'s code**. Caught **only** by the sha256 baseline the mutation discipline already demanded for a different reason (`99f3568…` → `6bc8514…`); recovered from the executor's own logged `git diff` and proven restored by **two independent hashes** — `git hash-object` returning the diff's recorded post-image `1fb60d78…` and sha256 returning `99f3568…`. **The rule: when the baseline is uncommitted work, the only safe revert is a copy taken before the mutation; the sha256 is what tells you which of the two you actually did.** This extends the iter-35 mutation-as-instrument discipline to its **other end** — validity covers applying a mutation AND undoing it. **A SECOND INSTRUMENT OF MY OWN CAME BACK CLEAN BECAUSE IT COULD NOT HAVE COME BACK DIRTY**: I reported "ZERO skips" from `go test ./...` piped through `grep -c -- "--- SKIP"`, but `--- SKIP` lines print **only under `-v`** — the zero was structural, not observed. Re-run with `-v`: **2 skips**, both pre-existing subprocess-helper guards (`crash_test.go:68`, `writer_lock_test.go:66`), benign — but "zero skips" has been treated as evidence here since B1. Same family as iter-37's glob-mangled `--include=*.go` and the zsh-empty `${PIPESTATUS[0]}`. **ONE MUTATION NAME, TWO FORMS, TWO ANSWERS — AND THIS TIME THE JUDGE HAD THE CONFOUNDED ONE.** Executor: `MUT-OUTCOME-BEFORE-RECORD` **2 red / 99 green**; judge, independently: **18 red**, every success path failing `refuses to persist an invalid ref ""`, concluding the form needs fault injection. Reading both applied patches settles it — the executor recomputed `recordObject(rec).Hash` so the outcome carried a **valid** ref and the mutation isolated the **ordering**; the judge passed the still-zero `recordRef`, so the store's own ref validator fired first and swamped the signal. **The executor's form is faithful, the plan's prediction stands, and CF-MJB-6 is REFUTED.** Fourth "one NAME two FORMS" instance (after `MUT-BENCH-DROP` iter-34, `MUT-PENDING-UNBOUNDED` iter-35) and the first where the **judge** ran the coarser form — a mutation's prose "form:" line does not determine the instrument; only the applied diff does. **THE JUDGE'S FINDING WAS REPRODUCED BEFORE BEING ACCEPTED, AND IT WAS UNDER-STATED.** CF-MJB-1 (filed non-blocking) says `recover.go`'s comment omits the Decision-5 residual. Measured: the doc's freeze block reads *"**One residual stays open and is stated, not hidden** … **No claim … may state or imply that every crash ambiguity is eliminated**"*, and the milestone had **deleted** the production comment carrying exactly that statement, replacing it with an unqualified *"Every commit and dispatched effect is therefore crash-detectable"* — the honest-claim discipline reversed in the one file a reader of `Recover` opens. **Fixed IN-PR, not carried**, with the full gate sweep re-run after. **THE EXECUTOR REFUTED ITS OWN PLAN, UNPROMPTED, AND WAS RIGHT**: `MUT-ORDINAL-ZERO-RESUME` was predicted at "1 red / ≥12 green, EVERY fresh-episode test green" and measured **11 red / 90 green**, because a fresh episode's *second* dispatch also collides with a constant ordinal 0 — mechanism held, discrimination claim too narrow; judge reproduced 11. **TWO DEVIATIONS THE EXECUTOR DID NOT REPORT**, both judge-confirmed, both non-blocking: **PD1 was explicitly violated** (a broker-side empty-episode guard was added where PD1 says "do NOT", and `broker_test.go:252` pins its exact string, so the store-layer guard is unreachable *through the broker* — it survives only because `journal_test.go:446` pins it directly); and **the capability ledger changed semantics unpinned** (the debit moved after the intent append, so the `no handler registered` path no longer debits — and `"no handler registered"` has **zero** tests, so neither the old nor the new behaviour was ever pinned). **THE CI GREEN WAS VERIFIED, NOT READ**: both jobs reported success in **36 s** and **8 s**, implausibly fast and exactly the shape of a gate that did not run — the step logs show all 11 modules ai-checked with `✓ 4/4 identities` and `✓ 14 named tests`, and `go build` + `go test ./... -count=1` with per-package `ok`; the Linux runner is simply ~3× the mac. Every gate re-run **OUTSIDE** the codex sandbox (the executor correctly labelled its own `verify_go.sh` UNINFORMATIVE UNDER SANDBOX with the verbatim bind error): `go build`/`go vet ./...`/`gofmt` clean, `go test ./... -count=1 -v` **rc=0, 180 PASS / 0 FAIL / 2 pre-existing SKIP across 10 packages**, `verify_go.sh` rc=0, `verify_ail.sh` rc=0 (4/4, 11 modules, 14 tests), broker-only `-race` rc=0, both forbidden-path assertions rc=0 **with a known-positive control**. **Queue: 4c [IN-SPRINT] — MJ.A + MJ.B LANDED, MJ.C is the last milestone**; **4d `w-ddl-gate-teeth`** and **4e `w-race-gate-blindspot`** still queued, and 4e's row is re-measured this iteration (the sibling package `host/broker` under `-race` is **GREEN in 4.095 s**, so the defect is `host/store`-local). New **CF-MJB-1…CF-MJB-6** (CF-MJB-1 DONE, CF-MJB-6 REFUTED); **CF-MJA-4/CF-MJA-5 were MJ.B-owned and are NOT discharged — they roll to MJ.C, stated rather than dropped**. Designer **not fired**; rotation pointer unchanged at `claude:claude-fable-5`. 4 stamps → newest 3 kept (iter-38, iter-37, iter-36); the positionally-4th (iter-35) archived this Gate 4.

## STATUS 2026-07-29 (iteration 37) — **`w-effect-journal` MJ.A LANDED** (PR #26 → squash `82d9128`, dev CI green both jobs SHA-addressed on the merge commit; judge sonnet **PASS 86/100, ZERO blocking**; executor `codex:gpt-5.6-sol` on codex-cli 0.145.0, `auth_mode=chatgpt`; planner **opus**; `metered=$0.00`). The effect-journal kernel reopen ships: **four new `host/store` methods** (`AppendNextEffectIntent` — the ordinal minted INSIDE the appending transaction — `AppendEffectOutcome`, `GetEffectReceipt`, `PendingEffectIntents`), **three changed** (`validateIntent` rejects the reserved `effect:` prefix, `GetReceipt` gains the same boundary guard **closing OD1**, `PendingIntents` gains the `world/journal-intent/v1` predicate), and **LAW 5 `effectDispatchLawful`** in the sketch — `ai-check` **8/8 contracts `status=verified`**, z3 4.16.0 present, so not a V27 silent skip. Closes **AC1–AC6, AC7b, AC7c, AC8**; `schema.sql`/`store.go`/`host/broker/**`/`host/replay/**` byte-unchanged. +733/−10 against the planner's 615 = **1.19×**. Path A's DDL-freedom re-verified first-party before routing (`schema.sql:81-87`: `CHECK (kind IN ('intent','outcome'))` + `UNIQUE (invocation_id, kind)`) — DDL-free **by construction, not by discount**. **11 named mutations**, each asserted to match **exactly once**, printed as a diff, compiled before running, reverted byte-identical; the judge independently re-ran two with **sha256-proven** reverts. **THE ITERATION'S SPINE: A GATE THAT IS NEVER RUN CANNOT FAIL, WHICH IS THE SAME DEFECT AS A GATE THAT CANNOT FAIL.** Chasing a line the executor mentioned in passing surfaced a defect in LANDED `w-store-durability` (SD.A, `86d1276`) code that two zero-blocking judgements never saw: on **clean `origin/dev` @ `d057de8`**, with the MJ.A diff absent and `scan.go`/`scan_test.go` untouched, `TestScanUnreadableLogKeysetResumes` **passes without `-race`** (filtered AND whole-package) and **fails 5/5 with `-race`** (filtered AND whole-package). Not `-count`, not the `-run` filter, and **not cgo** (`CGO_ENABLED=1` is the rig default and passes without `-race`). Exactly one element differs — `Rows[0].Field` is `""` under `-race` and `"prevEntryHash"` without it, while `Reason` from the SAME loop iteration is correct, and `fields` is a five-element literal indexed over a `[5]string`, so an empty `Field` should be unreachable. **Zero `DATA RACE` warnings**; separately the package under `-race` stops progressing and is `signal: killed` at ~142 s. Root enabler: **`-race` appears NOWHERE** in `.github/workflows/`, `verify_go.sh` or `verify_ail.sh`. **Mechanism UNKNOWN and labelled so** (standing UNVERIFIED hypothesis: `modernc.org/sqlite`'s heavy `unsafe` plus `-race`'s altered layout, or a go1.26.4 darwin/arm64 toolchain bug) — because zero `DATA RACE` output means *"just add `-race` to CI"* is not yet known-good. Queued as **4e `w-race-gate-blindspot`**, the **sixth** instance of this mission's signature shape and the first where the gate cannot fail *because it is never invoked at all*. **THE LOOP REFUTED ITSELF TWICE, IN BOTH DIRECTIONS, AND BOTH ARE WINS.** The **planner refuted the design doc**: premise V28 argued the ordinal race from "zero mutex hits", but the real serialization point is `store.go:253 SetMaxOpenConns(1)`, whose landed comment already calls it "the sole serialization point" — and the consequence was actionable, since under `MUT-ORDINAL-SPLIT-TX` the losing goroutine **errors** rather than returning a duplicate ID, so an AC7b test that merely compared the two returned IDs would have been **vacuous under its own named mutation**; the shipped test asserts `err == nil` FIRST (`journal_test.go:526-527`). The **executor then refuted the planner**: that mutant's loser returns `DuplicateInvocationError`, not the predicted raw SQLite `UNIQUE (2067)`, because compare-and-append duplicate discipline is preserved — the gate PREDICTION held, the stated MECHANISM did not, recorded rather than waved through. **The judge's finding was reproduced before being accepted, and the reproduction changed the fix**: CF-MJA-2 alleged AC5's "idempotent identical-bytes re-append" was unimplemented; measured, landed `AppendIntent` **is** identical-bytes idempotent while landed `AppendOutcome` is **not**, and `AppendEffectOutcome` mirrors `AppendOutcome` — so AC5's third clause describes behaviour absent from the commit side too. The judge offered "fix the doc **or** add the path"; the measurement rules out the second, since adding it to the effect side alone would diverge from the substrate. **AC5 is met on its two real clauses and its third is STRUCK as a doc defect** — recorded explicitly rather than silently checking the box. **Sandbox discipline held both ways**: the executor correctly labelled its own `go test ./...` and `verify_go.sh` **UNINFORMATIVE UNDER SANDBOX** (loopback bind denial), and every gate was re-run outside it — `go build`/`go vet`/`gofmt` clean, `go test ./... -count=1` **rc=0 across 10 packages with ZERO skips**, `verify_go.sh` rc=0, `verify_ail.sh` rc=0 (4/4 identities, 11 modules, 14 named tests). **TWO INSTRUMENT DEFECTS IN MY OWN GATE-2 COMMANDS, both the vacuous-pass family**: (i) a verification `grep` used `--include=*.go`, which **zsh glob-expanded so the commands never ran**, returning a clean-looking `0` — caught ONLY because the same call carried a **known-positive control** that also came back empty; re-run with a tool that cannot miss, `NewSession` has **11** call sites (one outside `host/broker`: `host/daemon/bench_test.go:474`). Gate 2's rule 3a limb (i) — *prove the instrument can see a positive* — paid for itself the first iteration after it was written. (ii) `rc=${PIPESTATUS[0]}` printed **empty** in zsh (it is `pipestatus`, 1-indexed), silently voiding two gate readings; re-run by direct invocation. **Queue: 4c [IN-SPRINT] — MJ.A LANDED, MJ.B and MJ.C remain**; new item **4e `w-race-gate-blindspot`** joins **4d `w-ddl-gate-teeth`**. New **CF-MJA-1…CF-MJA-5** (CF-MJA-3 DONE — PR #26 carries the mutation table). Designer **not fired** (no new doc); rotation pointer unchanged at `claude:claude-fable-5`. Codex probed **WITH `--model`** per the iter-19 process fix. 4 stamps → newest 3 kept (iter-37, iter-36, iter-35); the positionally-4th (iter-34) archived this Gate 4.

## STATUS 2026-07-29 (iteration 36) — **`w-effect-journal` (item 4c) NEW-DOC LANDED + QUORUM-CLEARED** (PR #25 → squash `fe582b5`, dev CI green **both jobs, SHA-addressed on the merge commit**; designer `claude:claude-fable-5` on the rotation lane, quota bucket; quorum metered **$0.164** against the $5 ceiling; **no sprint routed — this was a design iteration**). **THE ITERATION'S SPINE: THE QUEUE ROW'S OWN COSTING CLAIM WAS FALSE, AND THE GATE IT CITED AS PROOF IS INERT.** The M3 planner costed this item as *"Cheaper than assumed — **no schema change** (the commit shape lives in the payload codec, not the DDL, so AC13's `sqlite_master` gate stays green)"* — a claim its own `plan.json` labels **"Planner's measured note"**, restated in the charter queue row and carried into the ratification record. The premise is TRUE (the eight commit refs do live in the payload codec, `journal.go:104-143`); the **conclusion does not follow**, because the journal table's *kind vocabulary* and its *cardinality* live in the DDL. Measured first-party at `0f2afad`: **P1** `INSERT … kind='effect-intent'` → `CHECK constraint failed: kind IN ('intent','outcome')`; **P2** a second `kind='intent'` for the same `invocation_id` → `UNIQUE constraint failed: journal.invocation_id, journal.kind`; **P3, the sharp one** — a widened CHECK re-applied to an **existing** store returns **rc=0, no error**, `sqlite_master` shows the DDL **unchanged**, and the new-kind insert then fails on the OLD constraint. With **zero** migration machinery anywhere in `host/` (no `user_version`, no `ALTER TABLE`), **any DDL change in this repo ships FAIL-OPEN** — new stores get the new schema, every existing store silently keeps the old one, nothing detects the disagreement. *A schema change that was never applied is indistinguishable from one that was* — **iteration 35's spine, one layer down, at the schema instead of the mutation.** **AND THE CITED GATE HAS NO TEETH THERE — three compiling mutations, reverted byte-identical**: `MUT-JOURNAL-DDL-WIDEN` leaves `go test ./...` **rc=0 across all 10 packages** (*nothing* in the repo guards the journal DDL); the gate's OWN `MUT-DDL-DRIFT` does red but — **read the MESSAGE, not the exit code** (the iter-34 rule paying off) — via the **sha256 source pin**, which `t.Fatalf`s before a database is built, so the `sqlite_master` comparison never runs; and `MUT-DDL-COMPARE-DEAD` (short-circuiting that comparison, both vars still used, `go vet` rc=0) leaves the whole package **GREEN**. The gate's real teeth are a source-text pin over the pre-journal *prefix* — genuine protection for pre-existing tables, delivered by a different mechanism than the one it is documented as having, and covering neither the journal table nor the upgrade path. **FIFTH instance of this mission's signature shape** (a gate no production change could fail) and the **first found at PICK time**, in a gate inherited from an item already COMPLETE and judged **zero-blocking** — exactly the exposure the iter-32 process fix names. Raised as its own queue row **4d `w-ddl-gate-teeth`**, not smuggled into 4c. **THE DESIGN**: Path A, DDL-free — per-effect synthetic IDs in a reserved `effect:<episodeID>:<ordinal>` namespace, two payload codecs, reusing the existing `kind` vocabulary; Path B (new kind labels) REJECTED on the fail-open evidence, since it would honestly cost a migration mechanism first (~2.5–3d, not ~1.5d). **TWO QUORUM ROUNDS, TWO COLLISIONS IN THE SAME FIELD.** r1 `gemini-3-1-pro` **reject**: the in-memory ordinal counter is never re-initialized after a crash, so a resumed broker's first dispatch collides with the durable `effect:<episodeID>:0` and **bricks the episode** — r0 had *frozen* the opposite claim ("collision-free by construction"), which is FALSE across a crash boundary and was corrected everywhere it was restated. `gpt5-6-sol` was **ABSENT — `budget`**, a pre-flight refusal at *"$0.1005 exceeds cap $0.1000"*, zero spend, degraded N−1 and named (never a silent pass). The r1 fix **declined the reviewer's own proposed variant with a reason**: `len(episode.History)` counts *records*, and an indeterminate effect is precisely an intent whose record was lost, so it still collides after exactly the AC7 crash. r2 (cap raised to $0.25 so the absentee could run; **both present**) `gemini-3-1-pro` **PASS**, `gpt5-6-sol` **reject** — the r1 fix is **not atomic**: derive-then-append are separate operations, so *"the resumption fix merely replaces the restart collision with a TOCTOU collision."* Applied **VERBATIM** under the charter's **narrow-refinement carve-out** (concrete reviewer-authored fix; determinism-only; direction endorsed — Path A, the journal shape, intent-before-dispatch and the ID namespace all accepted): one transactional `AppendNextEffectIntent(episodeID, intentWithoutID) (id, ordinal, err)` mints the ordinal **inside** the appending transaction, `AppendEffectIntent` and `NextEffectOrdinal` are both REMOVED, the broker holds **no ordinal state at all**, `MaxInt64` exhaustion gets a structured error, and the reviewer's two named tests land as **AC7b** with its evidence ask as **AC7c**; `gemini`'s `GetReceipt` namespace guard adopted too, **closing OD1**. Nothing was force-passed and no objection was overridden — each was *satisfied*. **V28 — HALF THE OBJECTION IS REFUTED BY LANDED CODE, AND THE FIX WAS APPLIED IN FULL ANYWAY**: two broker *processes* cannot share a store as writers (`store.Open`'s non-waiting exclusive lock → `*WriterAlreadyActive`, proven cross-process), but the in-process half is REAL — **zero mutexes** in `store.go`/`journal.go`, so two goroutines can interleave freely. *Narrowing a fix to the part of an objection you can refute is how a real defect survives a review.* **V29 — A PROCESS FINDING AGAINST MYSELF**: my own premise-row additions (V21–V25) pushed the doc to ~13,952 input tokens and tripped round 1's per-reviewer cap, costing that round the very reviewer that later found the TOCTOU — **enriching a doc degraded its own review**. Also: **my Gate-3b poll was a broken instrument of my own making** — `target=$(git rev-parse A || cd B && git rev-parse HEAD)` binds as `(A || cd) && HEAD`, so `$target` held TWO SHAs, every API call was malformed, and the loop printed blanks toward its deadline; the iter-24/iter-107 meta-rule ("use the shipped snippet verbatim; a hand-rolled variant is a new defect surface") caught in the act, third instance. Re-run with a single literal SHA: **both jobs `completed/success`** on PR head `de0fb59` and again on merge commit `fe582b5`. Designer-verified numbers **re-measured by the controller rather than cited** (V25/V25b): `storejournal.ail` **7/7 contracts verified**, **`len(tests[]) == 30`**, `passed_tests == 37` reported separately; and Path A's mechanism confirmed on **`modernc.org/sqlite`**, not just the designer's sqlite3 CLI. **Queue: 4c [NEXT — QUORUM-CLEARED, routes straight to sprint-planner]; new item 4d `w-ddl-gate-teeth` queued.** Designer rotation advanced `codex:gpt-5.6-sol` → `claude:claude-fable-5`. 4 stamps → newest 3 kept (iter-36, iter-35, iter-34); the positionally-4th (iter-32) archived this Gate 4.

## CURRENT GOAL

1. ~~Iteration 0: ratify the bar~~ — **DONE 2026-07-23 attended** (Mark: clause-4 fixed at
   −2pp / ≤25%; bar + Conflict Surface + guardrails + queue ratified as drafted; re-quorum run
   attended — see STATUS).
2. **NOW**: work the queue through the inner loop (design-doc → sprint-plan → execute →
   evaluate), one sprint-sized item per iteration, recording routing evidence every time.
   **D1 RATIFIED 2026-07-24 (Mark, attended: A+B-metadata) — queue UNBLOCKED.**
   Next item: `[NEXT] w-world-library-m1` (clause-1) — M1 may now freeze the log format per the
   settled decision doc.

## The bar — what "Ailang World 1.0" must meet (RATIFIED by Mark 2026-07-23, attended — see STATUS + issue #1 record)

<!-- Distilled from DESIGN.md §14–§17 milestones + §12.4 value thesis + §13.5 priorities.
     Deliberately fewer clauses than milestones: clauses are CHECKABLE END-STATES, not work items. -->

- **Clause 1 — deterministic kernel**: an immutable, content-addressed world store with an
  append-only transition log and **proven deterministic replay** — replaying any recorded episode
  reconstructs state bit-for-bit from pure transitions + recorded effect results. The log-epoch
  semantics decision (DESIGN.md open question 11: interpreter version pinning + content-addressed
  transition functions) is made and enforced BEFORE the log format freezes.
- **Clause 2 — local-first daemon**: `ailang-worldd` runs on one machine over SQLite with CLI +
  REST, **zero cloud dependencies in the core**. Cloud transports exist only as effect-handler
  extensions (or not at all in 1.0).
- **Clause 3 — explicit authority end-to-end**: every effect goes through the broker with a
  capability + budget check; effect results are recorded (replay input); capsules run with a
  physical isolation floor beneath the semantic checks. No ambient-authority path exists from an
  agent to the outside world.
- **Clause 4 — the M4 RESIDENT-AGENT NON-INFERIORITY FLOOR (do-no-harm kill-switch; thresholds
  ratified 2026-07-23; agent definition + framing AMENDED by Mark same day, attended)**: this
  clause is the FLOOR, deliberately NOT the value proof — agents are World's *residents*, not
  its destination, and a substrate that taxes its residents dies before its systemic value can
  accrue. World's VALUE is measured by clause 5's provenance teeth and the R1 standing evidence
  (below). The floor: **TWO reference agents from different providers — Claude Code (agent mode)
  and codex CLI** — each in paired arms: *shell* (native tools) vs *World* (MCP transition tools
  only), same model, same benchmark set (standard tier), N≥3 runs per arm. **The floor holds
  only if BOTH reference agents show: World pass-rate ≥ shell pass-rate − 2 percentage points
  AND median wall-clock overhead ≤ +25%** (non-inferiority, not superiority — by design).
  **Stability precondition (gate validity, not gate outcome)**: a reference agent is
  gate-eligible only if its SHELL arm completes N≥3 runs with zero harness-fault failures
  (api_error / resource_limit / harness classes) and pass-rate range ≤5pp; an ineligible agent
  is fixed or substituted — if both become ineligible the gate PAUSES for instrumentation (a
  measurement problem, explicitly NOT a value-park). If the gate fails on eligible agents after
  honest tuning, **World parks**. **motoko is an OPTIONAL third arm** — the local-first/R5
  narrative, run only if it meets the same precondition, informative and NEVER gate-blocking
  (amendment rationale: motoko's current instability would make the baseline metrologically
  unsound — Mark 2026-07-23).
- **Clause 5 — the human surface works, with PROVENANCE TEETH (the 1.0 value demonstration;
  strengthened by Mark 2026-07-23)**: one real goal expressed → proposals → approval inbox with
  evidence bundles → commit (SCENARIOS.md scenarios 1–3, live) — AND, measurably: **on ≥3 REAL
  "why did X happen" questions (arising from actual operation, not synthetic), a provenance walk
  yields the verified answer in ≤5 minutes each**, where the pre-World method was a grep/log
  archaeology session. This — capability the shell cannot express at any pass-rate — is what
  World is FOR; clause 4 only proves it isn't paid for with resident-agent regression.
  **Standing value evidence beyond the 1.0 bar: R1** — this mission itself migrating onto World,
  measured on incident classes eliminated + human attention per workstream against the 2026-07
  markdown-and-launchd baseline (tonight's operation IS the control arm).
- **Clause 6 — protocol-native boundary**: the transition registry is served over MCP
  (capability-filtered per session) and an A2A agent card is published; **no new wire protocols**
  (DESIGN.md §3.7).
- **Clause 7 — controlled self-modification proven once**: at least one World behavior change
  ships as an extension package through World's own propose → verify → commit pipeline, with the
  §14 boundaries enforced (mission-loop machinery, AILANG compiler, live daemon EXCLUDED from
  self-mod scope).

## Conflict Surface (RATIFIED 2026-07-23 with the bar; drafted post-iter-0 to answer quorum objection #2)

What this mission touches or overlaps, and the drawn boundaries:

- **`ailang serve-api` (ailang repo)** — the overlap is the PROTOCOL BOUNDARY only.
  **VERIFIED (code inspection 2026-07-23, `internal/apiserver/` at ailang HEAD)**: the package
  advertises NO persistence or scheduling API — zero sqlite/sql.Open/bolt/badger/scheduler hits
  package-wide; request-scoped serving of module exports over REST/MCP/A2A (`mcp.go`, `a2a.go`,
  `handler.go`, `routes_dispatch.go`). Precision note (quorum round-5): the load-bearing form of
  this claim is "no advertised state machinery to build a world kernel on" — verified from the
  package surface; exhaustive absence-of-state through transitive calls is neither claimed nor
  needed, because the daemon justification runs the other way: worldd needs a store + log +
  broker + scheduler that apiserver does not offer, while apiserver's protocol serving IS reused
  (path (a), demonstrated live). A new daemon is justified by state, not by protocol.
  **ALSO VERIFIED**: `internal/apiserver` is a Go `internal/` package → NOT importable from
  another repo. "Reuse" is therefore one of three concrete paths, in preference order:
  (a) **primary**: World exposes its transition registry AS `.ail` modules served by
  `ailang serve-api --mcp/--a2a` — **VERIFIED for static module exports (live test
  2026-07-23)**: served `design_docs/sketches/` over `--mcp-http`, and `tools/list` returned
  `plan`/`verify`/`commit` as MCP tools with schemas + effect signatures. **NOT claimed**
  (explicitly w-mcp-projection's ACCEPTANCE CRITERIA, not premises): dynamic worldd-backed
  registry, per-session capability filtering, propose→verify→commit enforcement at the
  boundary — if those fail there, fall to (b)/(c); (b) fallback: worldd runs `ailang
  serve-api` as a sidecar process; (c) only on evidence: upstream issue to export a public
  serving package. Grafting a world
  store onto serve-api itself stays ruled out — it would move OS concerns into the frozen
  language repo, the wrong direction under PROGRAM.md.
- **The `ailang` binary / compiler (frozen core)** — worldd consumes the RELEASED binary
  (`check`/`test`/`ai-check`/serve-api) and never links compiler internals; language gaps route
  upstream as issues (guardrail above).
- **Coordinator / Collaboration Hub (ailang repo)** — pattern overlap only (approval queue,
  messaging). Patterns port; schemas do not. World's human surfaces are projections from the
  world store (SCENARIOS.md scenario 1), not extensions of the Hub; no shared database.
- **The mission loop itself** — the loop that builds this repo stays OUTSIDE World's
  self-modification scope (DESIGN.md §14). When M4 puts motoko on World, the mission loop is
  unchanged; only motoko's execution substrate changes.
- **OPEN for ratification**: whether worldd lives in THIS repo as its own Go module (DESIGN.md
  §15 assumption, coordinator-recommended: keeps the language repo frozen, lets World iterate at
  its own cadence) vs an `ailang world` subcommand upstream. Revisit only on concrete
  binary-distribution pain.

## Guardrails (mission-specific; the skill's Standing Rules always apply on top)

- **Local-first is inviolable**: no cloud dependency may enter the World core. Cloud = effect
  handlers behind the broker, post-1.0 unless ratified otherwise.
- **A correction is not applied until it reaches EVERY artifact that restates it** (added iter-28;
  **two instances in that one iteration**, and it is the sibling of the "our own prior record is a
  claim" guardrail — that one says re-measure before restating; this one says *propagate before
  believing the record is fixed*). A quorum round, a human ratification, or a first-party
  re-measurement that changes a **number, a field list, or a law's arity** must be applied to every
  place the doc-set asserts it: the prose, the **Acceptance Criteria**, the **Premise Verification
  Log** row, the **Design Freeze** bullet, the **`.ail` sketch**, and any mutation that names it.
  Iter-28 measured the cost: round 2 corrected "seven ref fields → eight" in premise V23 and stopped
  there, so Decision 1's prose and AC2 still said seven — and the field an implementer drops is
  `NextWorld.Ref`, the CLASS 3 **wedge**, which would have left the item's worst failure mode open
  **while every gate went green**. The same round's widening of LAW 6 `intentBindsCommit` never
  reached `storejournal.ail`, where the Design Freeze pins the Go mirror by drift test — so the gate
  would have *certified* the narrow binding the doc calls the defect. **Procedure**: after any such
  correction, `grep` the doc-set for the OLD value and for the artifact's own name, fix every hit, and
  strike-through rather than delete so the correction is auditable. A partially-propagated correction
  is more dangerous than an uncorrected one, because the record now *looks* settled.
- **A PRESCRIBED FIX IS A CLAIM — measure it before adopting it, especially one this mission wrote**
  (added iter-29, the **third** instance of the propagation class and the first where the *remedy*
  was the defect). When a prior iteration hands the next one a fix specified down to exact numbers
  ("widen X, add the required row, update the counts to 26 / 33"), the numbers arrive carrying the
  authority of a completed investigation, and applying them feels like compliance rather than
  judgment. Iter-28's prescription was written into the design doc, this charter AND the sprint
  plan, and it was **vacuous**: measured first-party, the 26-row form it specified reds **nothing**
  when `TransitionRef` alone is dropped from the Go compare (`failed=0`), because the single
  REQUIRED `EntryHash`-preserving row mutates three fields *together* and never touches the fourth
  — while the very mutation the doc demands (`MUT-INTENT-NARROW-BIND`) requires the added fields be
  load-bearing **"individually and not decorative"**. **Rule**: before implementing a
  numeric/structural prescription inherited from any prior iteration, build BOTH the prescribed form
  and the form you would have chosen, run the mutation the prescription exists to catch against
  each, and adopt on the measurement — then correct the prescription everywhere it was restated
  (here: 25→**30** / 32→**37**, doc + charter + plan JSON + new premise row V28). The tell that you
  are at risk: you are about to type a number you did not measure.
- **Allocate a carry-forward ID by reading the PREVIOUS entry's carry-forward block, not the
  open-CF list** (added iter-31). The planner and the judge independently proposed **CF-I-1** for
  the stale `effectAllowed` comment; iteration 30 had already used and *closed* CF-I-1 for a
  different finding, so it was absent from the open list precisely *because* it was resolved. In an
  append-only log two findings sharing an ID is unrecoverable in the same way a bare COUNT is (the
  iter-19 rule) — worse, because the collision reads as continuity. **Rule**: before allocating,
  read the last entry's full carry-forward paragraph *including its CLOSED lines*, take the next
  unused **letter** for the iteration, and never reuse a letter a prior iteration retired. A
  sub-agent's proposed ID is a suggestion; the controller owns the namespace. This iteration is
  CF-**J**-*.
- **Every acceptance check names exactly ONE owning milestone, and that milestone must hold the
  file that can fail it** (added iter-29). Milestone acceptance lists written as RANGES silently
  claim checks nobody can discharge. SD.B's list read `AC5–AC8`, which handed it **AC6** — whose
  only proof mechanism is crash injection — while `crash_test.go` sat in **SD.C's** file list and
  SD.C's own list read `AC10–AC11, AC14`. So the milestone that owned the check had no test for it,
  the milestone with the test was never required to close it, SD.C's close-out ("doc →
  `implemented/` with every box checked") would not have caught it, and AC6's required RED mutation
  `MUT-SPLIT-TX` belonged to a test no milestone owned. That is **an acceptance check no gate can
  fail** — the same shape as a silent `t.Skip` or a z3-less `ai-check`, at the level of the plan
  instead of the code. **Rule**: enumerate acceptance checks per milestone (never `ACn–ACm`), and at
  Gate 2 cross-check that every AC in the doc appears in exactly one milestone list AND that the
  milestone owning it also owns a file capable of reddening it. Structural implementation is not
  discharge: SD.B writing the receipt inside `Commit`'s transaction is not proof it survives process
  death — only the subprocess kill is.
- **Report the gate list the PLAN declares binding, never a remembered subset** (added iter-28). The
  controller ran three gates, reported them green, and the judge FAILED the milestone on a **RED
  fourth** (`bench_worldd.sh --smoke`) that the plan and handoff both list as binding on every
  milestone. **A gate table that omits a binding gate is a false green** — the same vacuous-pass
  class this mission has closed twice in code (silent z3 skip, silent `t.Skip`), committed in the
  *reporting* layer instead. Before writing any gate table into a commit, PR, log entry or report:
  re-read the milestone's `verify_commands` / `gates_binding_on_every_milestone` from the plan JSON
  and run **every** entry. Corollary: a gate whose failure you attribute to the environment is not a
  passed gate — attribute it *and* re-run it somewhere the surface actually works (iter-28: the codex
  sandbox denies loopback binds, and that panic masked a real startup deadlock).
- **Frozen-kernel discipline from day 1**: default every improvement to an extension; kernel
  changes require explicit human ratification (quorum evidence attached).
- **Never touch the shared loop machinery from this mission**: `tools/launchd/*` (shared driver +
  plist) and the symlinked skills are owned by the shared infrastructure; improvements route via
  Gate-5 retro in the shared skill so ALL missions benefit. NEVER copy skills into this repo.
- **Never touch the V1 mission's state or checkout** (`~/.ailang/state/mission-v1*`, legacy
  `mission-control.*` paths, the `sunholo-data/ailang` working tree).
- **Language gaps route upstream — two channels, always both**: anything AILANG can't express
  cleanly (or any core/binary change World needs) becomes (1) a GitHub issue on
  `sunholo-data/ailang` — the durable, triageable record with repro + version per that repo's
  conventions — AND (2) an agent message to the v1 loop's inbox so it is seen without waiting for
  a triage sweep: `ailang messages send mission-control "<summary + issue link>" --title "..."
  --from "world-coordinator"` (local send on this rig; the v1 session's start hook surfaces
  unread messages). The V1 mission owns routing it (extension vs core-floor per PROGRAM.md);
  World NEVER works around the compiler locally, no vendored forks.
- **CODING STANDARDS are binding** (ratified attended 2026-07-24): [coding-standards.md](coding-standards.md) — Z3 contracts on the pure core (S1), effects at the boundary (S2), **slim kernel / package-first** (S3: every kernel addition answers "why is this not a package?"), compiler-checked docs (S4), the AILANG fluency protocol (S5), honest non-vacuous gates (S6). Evaluators score against it; changes to it are ratification-class.
- **Compiler-checked docs discipline**: every `.ail` snippet that enters `design_docs/` ships as a
  checkable file under `design_docs/sketches/` (or a successor source tree) and passes the CI
  gate. A doc claiming "this compiles" without a checked artifact is a defect.
- **Kernel performance budget from day 1** (DESIGN.md §13.5): the M4 overhead gate gets HARDER as
  models get faster — fork cost and verify latency are design constraints on M1, not later
  optimizations.
- **Advisory quorum rounds are BOUNDED** (ratified by Mark 2026-07-23, from the 5-round
  iteration-0 ledger): on ratification-class docs (charter, bar changes), at most **2 quorum
  rounds per revision cycle**. The quorum's job there is objection-SURFACING for the human
  authority gate — it holds no veto, a synthesized PASS is not required, and rounds must never
  be re-run in search of one. When objections stop being load-bearing, close on the human gate
  and record the ledger. (Design-doc quorums at pick keep their normal semantics — BLOCKED
  parks the item.)
- **A MUTATION'S RESULT IS ONLY EVIDENCE IF THE MUTATION ITSELF IS VALID** (process fix, iter-22;
  **2nd instance** — iter-20 was the 1st). Mutation testing is how this mission proves a gate is
  non-vacuous, so a bad mutation silently produces a false finding in EITHER direction. Two
  distinct failure modes have now been observed:
  1. **The mutation does not compile.** Two iter-22 attempts failed `declared and not used`. A
     build break proves nothing about a test's strength — it is not a RED. Reformulate into a
     compiling, behaviour-changing form (e.g. keep the variable used: `from*0 + offset`,
     `*payload && len(path) < 0`) and re-run before scoring.
  2. **The mutation compiles but changes no behaviour**, so its GREEN reads as "the gate is
     vacuous" when the gate is fine. Iter-20's counter was appended before a newline split (a
     silent no-op whose green briefly read as "CF-7 is not load-bearing"); iter-22 collapsed the
     registry client's per-segment `url.PathEscape` into a whole-string one, which is
     behaviour-EQUIVALENT because `PathEscape` encodes `/` as `%2F` and Go's mux unescapes the
     `{name...}` wildcard to the identical `PathValue`.
  **The rule**: before recording ANY mutation result, confirm the mutation (a) COMPILED and
  (b) actually changed observable behaviour. An unexpected GREEN is a claim about your mutation
  first and about the gate second — re-check the mutation before concluding the gate is weak.
  Record refuted mutations in the log's Ruled out ledger rather than dropping them; both times,
  re-checking turned a would-be false finding into the strongest evidence in the set.
- **AN ATTENDED-AUTHORED DOC STILL OWES THE DESIGN-DOC GATES** (process fix, iter-26; **2nd
  instance** — iter-0 was the 1st). A doc written in an attended session with Mark bypasses
  `design-doc-creator`, and therefore bypasses its hard gates — most consequentially the
  **Premise Verification Log** and the **Conflict Surface**. Twice now, a doc that skipped that
  skill arrived at quorum and was blocked by *exactly* those two objections, from the same two
  reviewers, in the same order: the **charter itself** at iteration 0 (`gpt5-6-sol`: "relies on
  unverified premises … wants a Premise Verification Log"; `gemini-3-1-pro`: "fails the Conflict
  Surface gate"), and **HUMAN-SURFACE.md** at iteration 26 (the same two, near-verbatim). This is
  not a reviewer quirk — attended authorship optimizes for shared human context, and shared
  context is exactly what makes an unverified premise feel established. `coding-standards.md`
  does not currently require either section, which is why nothing catches it earlier.
  **The rule**: when the pick is a doc authored attended (or otherwise not produced by
  `design-doc-creator`), the controller adds the two sections — a Premise Verification Log with
  one row per load-bearing dependency, and a Conflict Surface for anything the doc proposes to
  build fresh — **before** spending the first quorum round, not after. Budget one designer pass
  for it. A round spent rediscovering a known-missing section is a round bought at full price.
- **THE MISSION'S OWN RECORD IS A CLAIM, NOT EVIDENCE** (process fix, iter-27; **2nd instance** —
  iter-26 was the 1st). The skill already forbids laundering a *sub-agent's* finding (iter-105), a
  *judge's* finding (iter-111) and a *survey row* (iter-25) into established fact. The gap it does
  not cover is the most seductive source of all: **a measurement this mission itself recorded in the
  charter or the log, in a prior iteration, in the controller's own voice.** It reads as first-party
  because it *was* first-party — for someone else, at a different commit, and it arrives with no
  visible provenance and no re-run.
  Instance 1 (iter-26): the controller wrote a premise row asserting "`host/replay` has no timeout
  tests" without running anything; it has three, `execTimeout = 60 * time.Second`.
  Instance 2 (iter-27): item 4b's row carried an eight-field `store.Commit` matrix presented as
  measured fact — *"seven produce a permanently unreadable row … the eighth (`NextWorld.Ref`) reads
  back fine, degenerate-but-readable"*. Re-run first-party it is **wrong twice**: two of the
  "seven" poison `GetWorld`, a different read surface entirely, and `NextWorld.Ref` is the **only
  unrecoverable field** of the eight — it wedges the store behind a non-`ConflictError`. That row
  had already been handed to a sub-agent directive **and** was on its way into a human ratification
  packet as decision-support.
  **The rule**: a measurement recorded by a PRIOR iteration is `UNVERIFIED, inherited from
  <iteration N>` until this iteration re-runs it. Re-run it before (a) routing it to a sub-agent,
  (b) surfacing it to Mark as decision-support, or (c) restating it in a new record. When a re-run
  refutes it, **correct the original in place with a visible strikethrough plus the corrected
  finding** — never silently overwrite it and never leave the stale prose standing next to the fix,
  because the whole failure mode is prose that outlives its measurement. Prose decays; a committed
  test does not, which is why the ghost-close rule exists and why a measurement worth keeping
  belongs in a fixture rather than a paragraph.
- **A named RED mutation is evidence only if it mutates the code the gate is supposed to be
  guarding. Mutate the PRODUCTION side — or write down what the gate cannot fail** (added
  iter-30; 2nd instance of *a gate that cannot fail*, after iter-29's AC6 owned by no milestone).
  The non-vacuity discipline says "name a mutation that turns this gate RED", and a mutation that
  reds *feels* like proof. It is not, if the thing you mutated is the test. SD.C's AC11 named
  `MUT-AUTO-RETRY` — re-dispatch indeterminate intents on recovery — and it duly reds two tests;
  but the probe consumer is test-local, so the mutation edits the test's **own helper** and the
  same file's assertions fail. **No change to `host/store` could ever fail that gate.** The
  discriminating experiment is one line of thought — *what would have to break in PRODUCTION for
  this to red?* — and one command: mutate the kernel instead. `MUT-RECEIPT-LIE` reds both
  never-retry tests, so the never-lie half had real teeth via a mutation the table never named,
  while a third AC11 test stayed green under **every** kernel mutation and was a sketch-drift
  check all along. The rules: (a) for each named mutation, state which FILE it edits and whether
  that file is production or test — a test-side mutation is a **drift check**, never a proof of a
  kernel property; (b) when a proof must be staged over a not-yet-existing consumer (test-local by
  necessity, as here), that is legitimate — but the acceptance check says so explicitly and carries
  a named carry-forward to the milestone where the mutation becomes production; (c) prefer the
  mutation that is DISCRIMINATING: `MUT-SPLIT-TX` reds exactly one of four crash stop points and
  leaves three green, which is why it proves atomicity — a mutation that reds everything localises
  nothing. **Corollary, same root** (iter-30's third finding): *a ratio measured at low sample
  count is a claim, not a number.* The executor reported the receipt tax as 2.8× from a 50-sample
  run; three 200-sample runs put it at 1.46–1.51×. Before a performance figure enters a doc, a
  charter stamp or a downstream item's premises, re-measure it at the sample size the file's own
  invocation line specifies — and if the two disagree, record **why** the first was wrong, not
  just the corrected value.
- **Kill switch stays until ratification**; only Mark (or the v1 agent on his instruction) arms
  the loop.

## Routing policy

Uses the **shared** per-role model routing from `mission-control` (controller /
designer-rotation / planner / executor / evaluator, generator≠judge enforced). Overrides for THIS
mission in `~/.config/ailang/mission-world.env`:

- **Executor default**: NON-Anthropic lane — `codex:gpt-5.6-sol` (`MISSION_EXECUTOR_MODEL`),
  per the shared quota plan. **iter-19 routing incident RESOLVED by Mark 2026-07-27 (attended,
  option c): codex-cli upgraded 0.137.0 → 0.145.0 (npm -g) and the pin LIVE-PROBED working**
  (`env -u OPENAI_API_KEY codex exec --model gpt-5.6-sol` → PIN-OK, ~21 tokens, ChatGPT-
  subscription OAuth — Mark verified server-codex OAuth2 same day, so the lane works headless
  in mission loops AND evals). The opus fallback stays the documented degrade path only.
  Upstream driver-probe gap (probe omits `--model`, false-greens an unusable pin) remains
  filed as ailang#486 for the v1 lane.
- **Evaluator**: must differ in provider from the executor (generator≠judge) — sonnet or
  gemini when the executor is codex.
- Otherwise: inherits the shared defaults.

## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])

<!-- Every open item carries a clause tag. Estimates are honest guesses at bootstrap;
     iteration 0 re-scores. NEW-DOC items start with design-doc-creator. -->

**[LANDED 2026-07-24 (iter-13), CI green on dev `d0009c8`] w-m1-ailang-hardening** · clause-1 ·
The M1 AILANG surface now USES AILANG's distinguishing features. Doc
`design_docs/implemented/w-m1-ailang-hardening.md` (Fable r1 designer; **codex/gpt-5.6-sol iter-13
revision**; controller carve-out 2nd revision). Shipped via **PR #5 → squash `d0009c8`,
sprint-evaluator PASS 97/100** (sonnet, generator≠judge). **4 Z3-proven contracts**
(`transitions/applyRevision`, `contracts/isValidNextWorld`, `logepoch/sameRef`+`servesEntry`) +
**14 named inline tests** + a **hardcoded, bounded, non-vacuous required-check-manifest gate**
(`scripts/verify_ail.sh`). **iter-12→13 descope (empirically forced):** a contract on any
`Proposal`-taking predicate Z3-errors `unknown sort 'Proposal'` (Proposal.evidence is an ADT) and
`ai-check` exits 0 SILENTLY → the 3 Proposal predicates (`proposalMatchesWorld`,
`verificationMatchesProposal`, `commitAllowed`) are documented-limitation + tests-only (6 tests,
shared-predicate bodies UNCHANGED → anti-drift preserved); achievable proven set = **4, not 7**.
Upstream `sunholo-data/ailang#477`. **Two new toolchain findings landed as V-rows:** V26 (`ai-check
-timeout` is per-function Z3 only, not process-wide; `ailang test` has none → the gate wraps both
legs in a hardcoded-deadline process-group-killing `run_bounded`, Standing Rule 6), and **V27
(`ai-check` shells to an external `z3` and SKIPS SILENTLY without it — a bare `ubuntu-latest` runner
has no z3, so every contract vanished from `verify.results[]` and PR #5's first CI went red; fixed
by installing Z3 4.16.0 in the `ailang-verify` CI job, sha256-pinned)**. Retrofit of M1 milestone-1
(`world/{contracts,transitions,logepoch}.ail`; `types.ail` byte-identical); root cause was
discoverability (`.mcp.json` + upstream #476). Effects/package-extensions correctly out of scope
(frozen core). See memory `ailang-feature-discoverability-gap`. · queued iter-10 (Mark-directed).

1. [LANDED 2026-07-24] **w-log-epoch-decision** · clause-1 · ALL THREE decisions settled: D2
   (content-addressed transition fns) + D3 (SHA-256 + tagged HashRef) via quorum; **D1 RATIFIED
   by Mark attended: A+B-metadata hybrid** — authoritative pin = exact-binary content hash;
   release tag+commit + semanticsEpoch as compatibility metadata in every header (M8 promotion =
   policy change, never a format migration); hermeticity → replay-doubling in M1 acceptance.
   Doc: `design_docs/planned/w-log-epoch-decision.md` (SETTLED — ready for M1).
2. [LANDED 2026-07-27 (iter-16), all 6 milestones CI green on dev; doc → `implemented/`]
   **w-world-library-m1** · clause-1 · Sprint AUTHORIZED (Mark, attended, iter-8-prior). Sprint plan
   (`.ailang/state/sprints/w-world-library-m1.plan.json`): 6 CI-green milestones. **M1 (pure AILANG
   semantic library — `world/{logepoch,types,contracts,transitions}.ail`, 285 LOC) LANDED** via
   PR #2 → squash `9d61d663e`, evaluator PASS 93/100. **M2 (Go host bootstrap — `go.mod` +
   `host/hashref` + `host/canon`, 635 LOC / 45 tests, Decisions 2+3) + CI `go-verify` job LANDED**
   via PR #3 → squash `d5b155c`, evaluator PASS 97/100 (generator≠judge). **M3 (SQLite store +
   atomic append-only log, Decision 4 — `host/store/{schema.sql,store.go,store_test.go}`, 1004 LOC;
   pinned pure-Go `modernc.org/sqlite v1.54.0`, single-transaction compare-and-append with
   `ConflictError`, pair-keyed verify cache) LANDED** via PR #4 → squash `a901c30`, evaluator
   PASS 88/100 (generator≠judge; −5 non-blocking `store_heads` schema-split, carried to M4).
   **M4 (interpreter artifact archive + epoch-1 registry bootstrap, Decisions 5+6 —
   `host/archive/archive.go` single-stream hash-while-copy + atomic-rename + sidecar manifest +
   `ReplayError`; `host/registry/registry.go` `world/epoch-registry/v1` idempotent bootstrap;
   `store_heads` M3 carry-forward folded into `schema.sql`) LANDED** via PR #6 → squash `8133573`,
   evaluator PASS 88/100 (generator≠judge; 2 non-blocking carry-forwards to M5/M6: a `KindExecFailure`
   test + a sidecar-absent idempotence-recovery edge). **M5 (replay engine + replay-doubling +
   fixture episode, Decision 7 — `host/replay/replay.go` authoritative replay resolving the exe from
   each entry's `Interpreter` HashRef, delegating per-transition execution to the archived binary as a
   bounded subprocess [never reimplements the interpreter, §14]; `replay_test.go` 12 tests: bit-for-bit,
   replay-doubling A==B==recorded, epoch-candidate-cannot-redirect, pair-member→cache-miss; committed
   goldens; folded the 2 M4 carry-forwards) LANDED** via PR #7 → squash `ef06937`, evaluator PASS 73/100
   (generator≠judge). **B1 CI false-green fixed in-PR** (`8ba5fe9`): go-verify now downloads the pinned
   **v0.30.0** binary (by TAG, sha256-verified, version-asserted) + exports `AILANG_BIN`, so the replay
   tests actually RUN in CI instead of silently `t.Skip`-ing — verified in the go-verify log. **M6 (CI Go
   gate finalize + `scripts/verify_go.sh` + final green sweep + carry-forwards NB2/NB5) LANDED** via PR #8
   → squash `a07ac96`, evaluator PASS 96/100 (generator≠judge, ZERO blocking): `scripts/verify_go.sh` = the
   durable local Go gate with an anti-false-green guard (fails loudly if `AILANG_BIN` unset/≠v0.30.0 — closes
   the V27/B1 silent-skip class locally too); CI `go-verify` now calls it; NB2 resolved by a genuine
   end-to-end replay through a byte-distinct second interpreter (`TestInterpreterMemberChangeDrivesReal-
   ReplayEndToEnd`). **ITEM COMPLETE**: all 16 acceptance boxes checked, doc → `design_docs/implemented/
   w-world-library-m1.md`. The semantic world library ships: World/Proposal/Transition/Evidence in AILANG
   (ai-check green), Go host (SQLite store + content-addressed objects + append-only log + archive +
   epoch registry + replay), replay proven bit-for-bit — the embedded-host-only M1 floor before
   daemon/broker/isolation (clauses 2–4, queued next).
   The item: the semantic world library: World/Proposal/Transition/
   Evidence types in AILANG (ai-check green in CI), Go host for the SQLite store +
   content-addressed objects + append-only log; replay of a recorded episode proven bit-for-bit ·
   ~2–3d. Design doc `design_docs/planned/w-world-library-m1.md` written (codex/gpt-5.6-sol
   designer, rotation) and quorum-**direction-accepted**; net-axiom +12, all 9 epoch-doc M1
   implications mapped, verified against the **pinned released `AILANG v0.30.0`** (checksum
   `sha256:ac3174e0…`, not the rig's `-dirty` build). Two gemini-3-1-pro quorum rounds each
   surfaced ONE narrow, non-direction defect (r1 dirty-binary → FIXED; r2 `ailang replay`
   Conflict-Surface overlap → clarified in-doc, §14-forced). ~~SPRINT PARKED~~ → **AUTHORIZED
   (see tag above); next fire builds M1.**
3. [**LANDED 2026-07-27 (iter-22) — ITEM COMPLETE**, all milestones dev-CI-green (both jobs), doc → `design_docs/implemented/w-worldd-m2.md` with EVERY acceptance + design-freeze box checked] **w-worldd-m2** · clause-2 ·
   `ailang-worldd` local daemon: SQLite, REST API, CLI; zero cloud deps; kernel perf budget
   measured and recorded from the first commit. **Shipped across five merges**: **A1** `b0deedb`
   (PR #10) the RATIFIED single-writer kernel change — `store.Open` is the fail-closed writer path
   (non-waiting `flock(LOCK_EX|LOCK_NB)`, structured `WriterAlreadyActive`), additive
   `store.OpenReadOnly`, `:memory:` exempt, `schema.sql` byte-unchanged, enforcement proven
   CROSS-PROCESS (a `LOCK_EX`→`LOCK_SH` mutation reds the suite, so an in-process mutex cannot
   pass); **A2** `39b2115` (PR #11) the daemon shell — loopback guard mirroring the Z3-proven
   `isLoopbackHost` with **no override flag**, the D7 bound-constant block (all four
   `http.Server` timeouts + `maxCommitBytes` pinned to the Z3-proven `withinCommitBytes`),
   fail-closed lifecycle with stdout address announce and bounded drain, `/v1/health` + `/v1/head`,
   and **`TestDaemonDependencyAllowlist` — zero-cloud ENFORCED over the real build graph, not
   asserted**; **A3** `9579fe1` (PR #12) Decision 6's perf machinery — p50/p95 harness on a
   fresh temp-FILE store, `bench/BASELINE.md`, and a bench-smoke gate asserting a **hardcoded
   benchmark-NAME manifest** (a dropped benchmark still exits 0 from `go test` — the V27/B1
   vacuous-pass class, closed for benchmarks); **M2.B** `b412699` (PR #13) the full REST v1
   surface — payload opt-in, `GET /v1/log` as a bounded loop with **zero new store methods**, the
   multi-segment registry wildcard, one shared JSON error envelope, body capped → 413, and the
   genesis-over-REST gap found independently by three parties; **M2.C** `73d3486` (PR #14) the
   CLI client verbs 1:1 over the frozen table (ONE transport path, injectable D7 deadline,
   bounded response + commit-file reads), a **real-subprocess** end-to-end episode driving every
   verb through the CLI (genesis → commit → read, plus a 409 re-plan), a bounded-deadline test
   that genuinely observes a timeout, and carry-forwards CF-B-1 + CF-B-4 folded. Design doc
   authored by the Fable designer (rotation) and quorum-run over two rounds; the single-writer
   fork was **ratified by Mark attended (arm A)**. Executor `codex:gpt-5.6-sol` (ratified pin) on
   three of five merges; independent judge `sonnet` throughout (generator≠judge, cross-provider).
   **Open carry-forward OUT of the item**: **CF-B-2** — `store.Commit` writes a zero
   `PrevEntryHash` that `store.GetLogEntry` cannot read back; the daemon refuses it at the
   boundary but the embedded API can still trigger it. Kernel-side decision → a store-hardening
   queue item, NOT M2. Judge's M2.C non-blocking findings: CF-C-1 (`--limit 0` indistinguishable
   from unset), CF-C-2 (registry escaping untested for encoded characters), CF-C-3 (CF-B-2 needs
   a tracking item + repro fixture), CF-C-4 (405 asserted on 2 of 7 GET routes).
4. [**[LANDED 2026-07-29 (iter-35) — ALL FIVE MILESTONES COMPLETE; doc → `implemented/`]** —
   M3.A LANDED 2026-07-28 (iter-31), PR #20 → squash `2edf2ef`, judge PASS
   84/100 · **M3.B0 LANDED 2026-07-28 (iter-32), dev CI green both jobs, PR #21 → squash
   `9401f2d`, judge sonnet PASS 88/100 zero-blocking** · **M3.B LANDED 2026-07-29 (iter-33), dev
   CI green both jobs, PR #22 → squash `10beb83`, judge sonnet PASS 88/100 zero-blocking** ·
   **M3.C LANDED 2026-07-29 (iter-34), dev CI green both jobs, PR #23 → squash `cae04d2`, judge
   sonnet PASS 88/100 zero-blocking** · **M3.D LANDED 2026-07-29 (iter-35), dev CI green both
   jobs, PR #24 → squash `4c4ff69`, judge sonnet PASS 93/100 zero-blocking — THE ITEM IS
   COMPLETE**.

   **M3.D LANDED — commit-boundary anchoring, the PRODUCTION recovery path, CF-H-1 CLOSED.**
   Ratified option (i) (Mark, attended, `c26b27d`) is now EXECUTED, not just recorded.
   `host/broker/recover.go` (138 LOC, **production**) pages `store.PendingIntents` with the
   kernel-owned `store.MaxPendingIntentsPage` and the `Seq` keyset cursor, reads `GetReceipt`, and
   surfaces `*IndeterminateEffectError` for every indeterminate intent — never dispatching, never
   auto-resolving, never appending an outcome, never re-executing. The episode's commit is
   anchored: effects run and are recorded → world+entry built → **then** `AppendIntent` → commit
   with `InvocationID` → `ReceiptResolved`, `PendingIntents` empty. Closes **AC14, AC16, AC19**;
   all 20 boxes checked, the CI box **only after it was OBSERVED** SHA-addressed.
   **AC16 discharged by a PRODUCTION mutation**: `MUT-AUTO-RETRY-PROD` on `recover.go` compiles and
   reds exactly two dispatch tests while **five stay green** — with the SD.C contrast stated (its
   version mutated the test's OWN helper, which is what V37 → CF-H-1 records).
   **THE ITERATION'S SPINE: a mutation that was never applied is indistinguishable from a mutation
   that was survived.** The first `MUT-PENDING-UNBOUNDED` run used a 4-tab pattern against 3-tab
   source; `str.replace` matched nothing and the all-green output looked exactly like a real
   result. The conclusion it supported was true **by luck** — being right by luck is not a method.
   Every later mutation ran under an assert that the pattern matched exactly once, with the applied
   diff shown before the suite. Two defects were then genuinely found in M3.D's own new code before
   it landed: an **unreachable** compile-time-false `retryAllowed(true,false)` guard (proven
   unreachable by a `panic` probe before removal), and a paging discipline with nothing to prove
   until a multi-page test was added. The judge (PASS 93/100) forced three further corrections, one
   of which **REFUTED the judge**: `MUT-PENDING-UNBOUNDED` names a FAMILY — two forms, both
   compiling, both redding the same one test by different mechanisms — so the record was accurate
   for the form run and simply never named it (iter-34's `MUT-BENCH-DROP` shape). Also corrected:
   my **"ZERO `t.Skip`" claim was FALSE for the package** (`handlers_test.go:183,187` → **CF-N-1**;
   the same missing env var makes `episode_test.go` fail loudly while two `Model.Infer` tests
   silently skip — one package, two answers), and `MUT-RECEIPT-LIE-CONSUMED` was **understated**
   (five tests red, not one — **understatement is a reporting defect too**). New **CF-N-1…CF-N-4**;
   **CF-H-1 CLOSED**; CF-M-3/CF-M-4 CLOSED at Gate 2.
   **The dispatch→record window remains OPEN** and is stated, not closed — item 4c owns it.

   **M3.C LANDED — the episode proof, the broker's price, and a retraction** (PR #23 → `cae04d2`;
   +686 across C-1/C-2 plus the retraction commit). Closes **AC8, AC11, AC12, AC13, AC17**;
   **AC14/AC19 deliberately NOT checked** (migrated to M3.D; `acceptance_check_numbering` still
   carries stale "OWNER M3.C" labels for both → **CF-M-4**). `host/broker/episode_test.go` proves
   AC8 across all THREE arms with **zero** replay dispatches; `MUT-REPLAY-DISPATCH-COUNT` reds
   ONLY the counting-stub assertion while every byte-identity assertion stays green, which is why
   the stub is written `== 0` and not `>= 0`. Bench manifest **8 → 10**; `bench/BASELINE.md`
   re-measured in ONE 200x invocation (three for the receipt ratio) by the controller outside the
   sandbox, after the executor wrote `<CONTROLLER-MEASURED>` into all 44 measured fields — the
   fifth milestone running. **The spine is a RETRACTION of the controller's own headline finding**,
   refuted by the judge citing premises **V14** and **V5** that this repo already held: the
   `skipped_tests: 5` property-skip class was measured, documented and DELIBERATELY excluded from
   the gate at M1, so filing it as a third V27/B1 silent-skip instance was an overclaim — retracted
   in `0ff48a6`, in a PR comment, and publicly on `sunholo-data/ailang#517`. **CF-L-5 answered YES**
   (`host/replay/replay.go:325` and `host/archive/archive.go:382` both repeat the iter-33
   process-group defect → **CF-M-2**). New **CF-M-1** (the gate reads no `skipped_tests`),
   **CF-M-3** (`MUT-BENCH-DROP` must move to the rename-form — the delete-form is a compile error
   for one of its two names).

   **M3.B LANDED — the handlers, the approval flow and the isolation floor** (PR #22 → `10beb83`;
   1,767 insertions across checkpoints B-1/B-2/B-3 against the planner's ~1,450 = **1.22×**).
   Closes **AC5, AC6, AC7, AC9, AC15**. ONE shared bounded-subprocess surface (`handlers.go`) so a
   single production mutation reds BOTH handlers' timeout and output-cap tests; `Git.Commit` with a
   scrubbed env and an empty HOME; `Model.Infer` over the pinned binary under `--ai-stub`; the
   frozen synchronous approval flow (`approve.go` — `Pending(requestRef)` as the FINAL result, the
   record never rewritten, `DecideApproval` an operator entry point writing a SEPARATE decision
   object, `Human.PollApproval` an ordinary effect); and `host/capsule`, the F1–F6 floor with one
   independently-reddable test per restriction and **ZERO skips**. `host/replay/**` byte-unchanged.
   **THE FINDING: a bounded-wait guarantee that did not hold on linux, in a mission whose Standing
   Rule 6 is "every wait is bounded".** `runBounded` killed only its DIRECT child, so a handler that
   forks left a grandchild holding the inherited stdout pipe and the capped read blocked until the
   GRANDCHILD exited — linux CI measured **`5.002891s` against a 40ms bound**, to three decimals the
   runtime of the `sleep 5`, while darwin ran the identical code in **42ms**. `scripts/verify_ail.sh`
   already carried the correction (**V26**: kill the whole process group); the Go exec surfaces
   repeated the mistake the shell gate had already fixed. Both corrected (`Setpgid` + `Kill(-pid)`).
   **It was exposed only because a CI-robustness fix WIDENED the discrimination gap rather than
   loosening the assertion** — the tempting fix (raise the guard past the observed 316ms) would have
   shipped it silently. Two more controller-found defects, each reproduced before being acted on:
   `Model.Infer` hand-rolled JSON escaping that emitted **invalid JSON for control bytes**
   (reachable from any payload; post-M3.B0 it bills a STANDING DEBIT for what is a host encoding
   bug), and the capsule never killing a child that overflowed the output cap, so **F6 silently
   degraded into F5** (`*TimeoutError` returned for an overflow — the shipped F6 fixture emits 513
   bytes, below one pipe buffer, so its child always exits and the defect was invisible to it).

   **M3.B0 EXECUTED THE RATIFICATION** (folded into the doc at `84b8efd`, shipped at `9401f2d`):
   Decision 3 now has a THIRD arm — **every failed effect writes exactly one record and the DEBIT
   STANDS** (never-lie applied to money — refunding a possibly-partially-executed effect would
   make the ledger lie about spend). Encoding, decided by the planner and measured before
   adoption: one `failed: bool` after `allowed`, the arm being the pair `(allowed, failed)` with
   `(F,T)` illegal and forbidden by the law; `world/effect-record/v1` does **NOT** bump (measured:
   no durable record exists anywhere, and M3.A's goldens still decode under
   `DisallowUnknownFields`). `recordConsistent` = three disjuncts, still **Z3-verified 7/7**;
   **`len(tests[])` 25 → 31** (`passed_tests` 33 — reported separately, NOT the gate); `world/`
   totals unperturbed at 4/4 identities / 14 named tests / 11 modules. 408 insertions / 7 files
   (planner predicted ~420 = **1.0×**). Closes **AC18**. Folded **CF-J-1** and the **CF-I-2 →
   CF-J-2** rename that iter-31's renumbering never delivered to the code.
   **It also closed a gate no production change could fail**: M3.A's drift test mirrored
   `recordConsistent` in a TEST-LOCAL copy, so it proved the TEST matched the sketch, never that
   PRODUCTION did. One instrument run unchanged before and after — forcing production
   `RecordConsistent` to `return true` red **1** subtest and **ZERO** drift rows at `84b8efd`, and
   reds **6** negative arms **+ 6** drift rows now. Found by the planner, reproduced by the
   controller, re-run independently by the judge. **Third instance of this shape in this mission.**
   **M3.D's `blocked_on` is CLEARED** in the plan (option (i) scope written; `MUT-AUTO-RETRY-PROD`
   is now a PRODUCTION mutation against `recover.go` — which is what makes CF-H-1 dischargeable).
   **M3.A shipped the law + core**: `design_docs/sketches/effectbroker.ail` (213 lines, Appendix A
   byte-verbatim by `diff`, **7/7 contracts Z3-verified**, **`len(tests[])` 25**) +
   `host/broker/{decide,broker,record,handlers_fs}.go` and 4 test files, 1,317 insertions. Frozen
   check order + frozen five-label set pinned by a drift test transcribing **all 25** sketch rows
   with their sketch line numbers; exactly ONE content-addressed record per request — denials
   included — written BEFORE any result is returned; enforcement against the REMAINING ledger;
   no wall clock; store injected, never opened. Closes **AC1, AC2, AC3, AC4, AC10**;
   `host/store/**` and all protected paths byte-unchanged by `git diff --exit-code`. Sprint plan +
   handoff written by the **opus** planner at
   `.ailang/state/sprints/w-effect-broker-m3.{plan.json,handoff.md}` (M3.A/B/C/D).

   > **M3.D — RATIFIED 2026-07-28 (Mark, attended, `c26b27d`): OPTION (i) NOW, (iii) QUEUED AS
   > `w-effect-journal` (item 4c).** M3.D executes episode/commit-boundary anchoring: the episode
   > driver appends the intent once world+entry are built and commits with `InvocationID`; the
   > broker gains a PRODUCTION `recover.go` consuming `PendingIntents`/`GetReceipt`, surfacing
   > `IndeterminateEffectError` and never auto-re-executing — which **closes CF-H-1** with a real
   > production mutation. The dispatch→record window stays OPEN and **must be claimed honestly**:
   > the Decision-3 supersession note must be corrected so it does not overclaim. The kernel reopen
   > that closes the window at effect granularity is **pre-ratified in principle** and lives in
   > item **4c**.
   > **ALSO RATIFIED — CF-J-2 gets a THIRD ARM (frozen Decision 3 REOPENED, human gate exercised):**
   > every FAILED effect writes a record, so audit and replay are complete; and **the debit
   > STANDS** — refunding a possibly-partially-executed effect would make the ledger lie about
   > spend, which is the never-lie law applied to money. `host/broker/handler_error_repro_test.go`
   > (the pinned reproduction landed this iteration) becomes that fix's red→green test.
   >
   > <details><summary>The original park record — what was asked and why (kept, not open)</summary>
   >
   > **THE RATIFIED FOLD-IN IS NOT
   > EXECUTABLE AS WRITTEN, and this is a scope-and-ratification call, not a planning one.**
   > Measured first-party at Gate 2: `validateIntent` (`host/store/journal.go:210`) requires SIX
   > non-zero **commit-shaped** refs and `bindCommitIntentTx` (`store.go:807-825`) compares ALL
   > EIGHT for byte-equality inside the transaction; all FOUR landed callers derive the intent from
   > an already-complete `Commit`. A brokered effect's RESULT is an INPUT to the transition that
   > produces the next world, so `EntryHash`/`WorldRef` are **not knowable before dispatch** — a
   > pre-dispatch intent for a general brokered effect is **structurally impossible**. SD.C's AC6
   > proof holds only because its "external effect" is a probe write whose content never feeds the
   > commit. **The landed substrate is a COMMIT journal; the broker needs an EFFECT journal.**
   > Every way out crosses Design Freeze bullet 8 ("zero `host/store` method changes") or
   > Decision 7 ("the broker writes objects and registry heads only — never log entries").
   > **THE THREE OPTIONS (full costing in `plan.json → human_gate_M3D`):**
   > **(i) RECOMMENDED — episode/commit-boundary anchoring.** The episode driver appends the intent
   > once the world+entry are built and commits with `InvocationID`; the broker gains a production
   > `recover.go` consuming `PendingIntents`/`GetReceipt`, surfacing `IndeterminateEffectError` and
   > never auto-re-executing. **Zero kernel cost**, ~350 LOC, **closes CF-H-1**. Leaves the
   > dispatch→record window OPEN — and says so; the Decision 3 supersession note must then be
   > corrected, because as written it claims more than M3 would deliver.
   > **(ii) Pre-dispatch "effect requested" log entry.** Closes the window fully, but breaks
   > Decision 7, doubles log growth (two entries + two world advances per effect) and makes the
   > broker a second writer competing for the head with no arbitration. **Planner recommends
   > reject.**
   > **(iii) An effect-shaped intent in `host/store`.** Closes the window at exactly the granularity
   > objection 2A asked for. **Cheaper than assumed — no schema change** (the commit shape lives in
   > the payload codec, not the DDL, so AC13's `sqlite_master` gate stays green), but it is a Go
   > kernel reopen: Design-Freeze-forbidden and ratification-class. **Planner recommends it become
   > its own item (`w-effect-journal`, ~1–1.5d) AFTER M3**, with M3 closing on (i)'s honest claim.
   > **CF-H-1 is dischargeable under ALL THREE options** — it needs a *production consumer*, not
   > effect-shaped intents. **Explicitly forbidden**: synthesizing placeholder refs to satisfy
   > `validateIntent` — it would make the durable intent a false statement, wedge
   > `bindCommitIntentTx`, and violate the never-lie law the substrate exists for.
   > **Default if unanswered: (i)**, with (iii) queued as its own item. M3.B and M3.C proceed
   > regardless.
   > </details>

   ~~[NEXT] — UNPARKED 2026-07-28 (iter-30). The scope question is ANSWERED and the dependency
   has LANDED.~~ Mark's attended triple ratification resolved it as **(a)-by-substrate**: the
   broker DEPENDS on 4b's journal, so the crash window closes structurally rather than by
   documentation, sequencing 4b → 4. **4b is now COMPLETE** (SD.A `86d1276` · SD.B `d5774eb` ·
   SD.C `6811604`), so nothing gates this item. **Routes straight to sprint-planner** — the doc is
   written and twice-quorumed, no re-design and no re-quorum needed. Three things the planner MUST
   fold in before writing milestones, all produced by 4b and none of them in the doc as authored:
   (1) the doc's Decision 3 "honest ordering limitation" paragraph and its Deferred-Scope
   write-ahead-journal row are **SUPERSEDED** (marked so at `6811604`) — M3 now consumes
   `store.AppendIntent` / `Commit.InvocationID` / `GetReceipt` and the three-state receipt law
   instead of deferring them; (2) **CF-H-1** — AC11's never-auto-re-execute proof is currently
   demonstrated over a *test-local* probe consumer, so its named mutation `MUT-AUTO-RETRY` is
   self-referential; M3 owns the real dispatch path and must re-run it as a **production**
   mutation, which is the first time that acceptance check can actually fail; (3) the measured
   **+50.9% receipt tax** (`bench/BASELINE.md`) — M3 is the first component to pay it on a real
   dispatch, and owns the question of whether Decision 7's +20% bound was ever right for two
   in-transaction inserts. The round-2 `gemini-3-1-pro` objection (no non-vacuity mutations for
   handler subprocess timeouts / output caps) stays **PRE-APPROVED to apply verbatim**, no
   re-quorum. ~~PARKED `needs-human-review` 2026-07-27 (iter-23)~~. Doc:
   `design_docs/planned/w-effect-broker-m3.md` (Fable designer, rotation; 1,036 lines;
   Appendix-A sketch **7/7 Z3-verified / 27 tests / 0 failures**, re-run first-party by the
   controller and byte-identical to the doc's appendix)] **w-effect-broker-m3** · clause-3 · effect broker with FS / Git / Model (`std/ai`) /
   Human.Approve handlers; effect-result recording; capability + budget checks; first physical
   isolation floor · ~2–3d. Builds on the landed `w-worldd-m2` daemon, whose single-writer
   authority is exactly the property that makes broker-mediated effects meaningful (an embedded
   writer bypassing capability/budget checks is the ambient-authority pattern clause 3 exists to
   end — Mark's ratification rationale, iter-18).

   > **ANSWERED — the block below is kept as the record of what was asked, not as an open gate.**
   > Mark ratified attended (charter STATUS, `bc467f1`): **neither (a) nor (b) as framed** — the
   > journal became its own item (4b) and the broker DEPENDS on it, so M3 gets the truthful broad
   > claim *without* the scope increase (a) would have cost, and does not have to weaken its claim
   > as (b) proposed. The controller's recorded recommendation was **(b)**; the human's answer was
   > better than the recommendation, and 4b's SD.A→SD.C evidence is why. As of iter-30 that
   > dependency is LANDED, so this question gates nothing.

   **THE PARKED QUESTION (`gpt5-6-sol`, quorum round 2 — answerable in one comment).** The broker
   dispatches a handler and *then* writes the effect record. If the process dies in between, the
   external effect (an `FS.Write`, a `Git.Commit`, a paid `Model.Infer`) has HAPPENED with no
   durable record, and replay cannot distinguish "never executed" from "executed but record lost"
   — a silent-duplicate-execution risk. The doc deferred this crash window to future store
   hardening; the reviewer says that contradicts the milestone's own headline claim that every
   effect result is recorded. **Which is M3?**
   **(a)** Close it now — add a durable pre-dispatch intent object + broker journal head, a
   post-dispatch outcome object, an `IndeterminateEffectError` on recovery that NEVER auto-
   re-executes, and crash-injection tests. Truthful claim, but a real scope increase (~+0.5–1d)
   and it overlaps the open **CF-B-2** store-hardening carry-forward.
   **(b)** Keep M3's scope and WEAKEN THE CLAIM to the reviewer's own wording — *"every attempted
   dispatch is durably detectable; completed outcomes are replayable; indeterminate attempts fail
   closed without live fallback"* — deferring the journal to the store-hardening item.
   **Controller's recommendation: (b)**, with the journal queued as its own item alongside CF-B-2
   — the honest narrower claim is worth more than an unproven broad one, and durability belongs
   with the kernel. **Default if unanswered: (b).** Recorded rather than force-applied because
   this is a scope-and-ratification call, not the narrow-refinement carve-out (see the doc's
   Quorum verification log for why 2A fails the carve-out's second limb).
   The OTHER round-2 objection (`gemini-3-1-pro`: no non-vacuity mutations for handler subprocess
   timeouts / output caps) **is** carve-out-eligible — concrete verbatim fix, completeness only —
   and is **PRE-APPROVED to apply on unpark**, no re-quorum needed for it.
4b. [**LANDED 2026-07-28 (iter-30) — ITEM COMPLETE**, all three milestones dev-CI-green (both
   jobs), doc → `design_docs/implemented/w-store-durability.md` with EVERY acceptance +
   design-freeze box checked. SD.A `86d1276` (iter-28) · SD.B `d5774eb` (iter-29) · **SD.C
   `6811604` (iter-30, PR #19, judge sonnet PASS 88/100, ZERO blocking)**.
   **SD.C closed the item**: `host/store/crash_test.go` (+315) proves **AC6** across REAL PROCESS
   DEATH — a re-exec'd helper advances to one of four NAMED stop points (`after-intent`,
   `after-external-effect`, `mid-commit-before-outcome`, `after-outcome`), is SIGKILLed, and the
   parent reopens the store and asserts the receipt law, `PendingIntents`, the probe effect and
   world/entry presence per stop point; `mid-commit` is arranged by helper logic through
   `commitBeforeOutcomeHook` (a production no-op **anchored to the outcome-write SITE**, not to a
   line), never by a sleep; every wait polls the CAPTURED `os.Process` under a deadline (AC10).
   `recover_test.go` (+162) is the probe consumer: never-lie surfacing, the `retryAllowed` gate,
   the deterministic commit-path reconciler, and a `Model.Infer`-shaped counting probe at ZERO
   dispatches (AC11). `bench_test.go` (+99) + `bench_worldd.sh` (+2 manifest names) +
   `bench/BASELINE.md` (all 8 rows in ONE 200x invocation) price the journal (AC14).
   **`MUT-SPLIT-TX` reds EXACTLY `mid-commit-before-outcome` while the other three stop points
   stay GREEN** — reproduced first-party by the controller, `store.go` reverted byte-identical.
   **Two things the close-out found by re-measuring rather than accepting.** (i) **AC11's named
   `MUT-AUTO-RETRY` is SELF-REFERENTIAL** — the probe consumer is test-local (correctly: the real
   consumer is M3's broker), so the mutation edits the test's own helper and no kernel change can
   fail it. A kernel-side `MUT-RECEIPT-LIE` DOES red both never-retry tests, so the never-lie half
   has teeth via a mutation the doc never named; `TestRecoverRetryAllowedMirrorsAllSketchRows` is
   provably kernel-independent. AC11's claim is downgraded in-doc to what SD.C can prove; the
   consumer half is **CF-H-1**, owned by item 4. (ii) **Decision 7's receipt target is BLOWN and
   recorded, not relaxed**: commit-with-receipt p95 **+50.9%** over a bare commit (0.4537 →
   0.6846 ms) against ≤ +20%, reproduced 1.51×/1.49×/1.46× over three 200x runs — and the
   executor's in-sandbox 50x reading of **2.8×** was a low-sample artifact, corrected in both
   files. Item 4 is the first component to pay this on a real dispatch path and owns the question
   of whether +20% was ever the right bound for two in-transaction inserts.
   ~~IN SPRINT — SD.C remains~~
   All three ratification arms answered by Mark, attended (`bc467f1`): ARM V1 · ARM J1 ·
   `Commit.InvocationID` + in-tx receipt binding · three-state receipt law · recovery never
   auto-re-executes. The M1 kernel reopen is RATIFIED. Doc:
   `design_docs/planned/w-store-durability.md` (Fable designer, rotation; header now RATIFIED +
   IN SPRINT) + `design_docs/sketches/storejournal.ail` (180 lines, **7/7 contracts Z3-verified**,
   30 named tests — controller-remeasured iter-29 after LAW 6's widening). Sprint plan + handoff written by the **opus** planner at
   `.ailang/state/sprints/w-store-durability.{plan.json,handoff.md}` (3 milestones SD.A/SD.B/SD.C).

   **SD.A IS LANDED — CF-B-2 IS CLOSED AT THE KERNEL WRITE PATH.** PR **#17** → squash **`86d1276`**,
   dev CI green; executor **`codex:gpt-5.6-sol`**; judge **sonnet across 3 rounds: FAIL 57 → PASS 91
   → MERGE 97** (generator≠judge cross-provider). Shipped: `store.Commit` validates **EIGHT** ref
   fields (5 log-entry + 3 world) plus each Object's `Hash`/`InterfaceHash`, **before `tx.Begin()`**,
   with structured `InvalidRefError{Op,Field,Text,Err}` + `IsInvalidRef`; the same `validateRef`
   guards `PutObject`/`PutWorld`/`SetRegistryHead`/`SelectHead`/`PutVerifyResult`. `ObservedHead`
   remains the single zero-legal COMPARED field (genesis). New `host/store/scan.go`: bounded,
   read-only `ScanUnreadableLog`/`ScanUnreadableWorlds` with **keyset cursors only** (never `OFFSET`,
   never rowid — a `VACUUM` renumbers rowids), `MaxIntegrityScanPage` + `InvalidLimitError` enforced
   **before any query**; the daemon pages both tables at startup to completion or a row/time budget
   and reports `integrity_hole` / `integrity_scan_incomplete` **with a continuation cursor**, so a
   truncated scan never reads as a clean bill of health. `durability_repro_test.go` rewritten from
   asserting the broken behaviour to asserting the post-fix contract; judge carry-forwards **CF-E-3**
   and **CF-E-5** folded. **AC1–AC4, AC12 and AC10's scan half are CLOSED.** `schema.sql`, `go.mod`,
   `go.sum`, `scripts/`, `world/`, `cmd/`, `.github/`,
   `host/{replay,hashref,canon,archive,registry}` **byte-unchanged**.

   > **SD.B's blocking precondition — RESOLVED iter-29, and the prescribed resolution was itself
   > under-propagated (THIRD instance of the one root cause).** LAW 6 `intentBindsCommit` is now the
   > round-2 **8-field** binding (16 params), applying the already-ratified Freeze ⇒ no human gate.
   > But the iter-28 prescription — one new `tests[]` row, `len(tests[])` 25→26 / `passed_tests`
   > 32→33 — was **measured vacuous before adoption**: the single REQUIRED `EntryHash`-preserving
   > row mutates `PrevEntryHash`/`TransitionFn`/`Interpreter` *together* and never touches
   > `TransitionRef`, so at 26 rows a Go mirror that drops `TransitionRef` alone reds **nothing**
   > (`failed=0`, first-party). `MUT-INTENT-NARROW-BIND` demands the four added fields be
   > load-bearing **individually**. Landed form: **10 rows** = all-match + one single-field mismatch
   > per commit-defining field + the REQUIRED combined row; **AC9 pins `len(tests[])` 30 /
   > `passed_tests` 37** (measured). The 16-param Z3 arity risk the planner flagged is **REFUTED**
   > (verifies first try, 7/7) — no upstream issue owed. Evidence: doc row **V28**.

   **Open non-blocking carry-forwards (enumerated — a bare COUNT is unrecoverable, iter-19 rule):**
   **CF-F-1** the daemon's `scanPageSize`/`scanRowBudget`/`scanTimeBudget` wiring is not pinned by a
   constant-equality test the way `TestBoundedWaitsAndBodyLimit` pins the D7 constants;
   **CF-F-2** `MUT-STORE-TOUCHED` does **not** red through the store-untouched assertions — moving
   validation after the world INSERT still rolls the transaction back, so the real placement guard is
   the `Commit_Object_Hash` subtest (executor-reported unprompted, judge-reproduced; shipped
   placement IS correct, the stated non-vacuity proof was weaker than the plan believed);
   **CF-F-4** `integrityFixture` (`host/daemon/integrity_test.go`) is killed at ~100 s under
   `-race` — PROVEN pre-existing (reproduces with iter-28's changes stashed; `host/store -race`
   finishes in 1.9 s), and `-race` is not one of the four binding gates, but trim the 70 raw-SQL
   inserts before anyone adds it. **CF-F-3 CLOSED** at `d506275` (see the charter STATUS stamp).

   **THE RATIFICATION PACKET (three arm choices, answerable in ONE comment).** The charter's
   frozen-kernel guardrail makes each of these human-only, and the sprint's own last acceptance box
   refuses to execute before they are settled:
   1. **Decision 1 — validate refs on WRITE (ARM V1, recommended) vs support the zero on READ
      (ARM V2).** Behavioural change to the landed `store.Commit` et al.
   2. **Decision 3 — the durable intent/outcome journal as an additive SQLite `journal` table
      (ARM J1, recommended) vs a zero-schema registry-head chain (ARM J2).** J1 is the **first-ever
      `schema.sql` change** (additive only; the landed `CREATE TABLE IF NOT EXISTS` mechanism
      already migrates it on writer open).
   3. **Decision 4 — `store.Commit` gains an optional `InvocationID` field**, with the receipt
      written inside the SAME transaction, so "committed iff receipt exists" is atomic by
      construction.

   **This packet can ride the SAME attended comment as item 4's parked (a)/(b) question — one human
   touchpoint unparks TWO items**, and this design is compatible with either answer (the milestone
   cut is pre-drawn in the doc's gating section).

   **What the iteration MEASURED, first-party (the item is bigger than the row that queued it).**
   Controller's own eight-field matrix on `store.Commit`: **it validates NONE of its ref fields** —
   all eight zero-value commits return `err=<nil>`. ~~**seven** produce a permanently unreadable row
   (`TransitionFn`, `Interpreter`, `EntryHash`, `TransitionRef`, `PrevEntryHash`,
   `NextWorld.LogHead`, `NextWorld.StateRoot`), and the eighth (`NextWorld.Ref`) commits and reads
   back *fine* as an empty-string ref that becomes the selected head — degenerate-but-readable, and
   therefore the one shape a read-side fix could never observe.~~ **CORRECTED iter-27 — see the
   three-class split below; the struck text conflated two different READ surfaces and called the
   only unrecoverable field the mildest.** Also measured: the poisoned head ADVANCES, and a
   subsequent perfectly legal commit chains onto the poisoned entry, so the append-only log grows a
   **permanent hole mid-chain with readable entries on both sides** and no detection or recovery
   path — **corroborated exactly, and now executable** as
   `TestCFB2PoisonedEntryLeavesPermanentHoleMidChain`.

   **THE REPRO-FIXTURE HALF IS DONE — LANDED 2026-07-28 (iter-27), PR #16 → squash `e8ba7b2`, dev
   CI green (both jobs), evaluator PASS 93/100 (sonnet, generator≠judge vs the codex executor).**
   `host/store/durability_repro_test.go` (225 LOC, 5 tests / 15 subtests) is the committed,
   CI-enforced reproduction the ghost-close rule demanded; tracking issue **#15** now exists. This
   closes the tracking half of judge carry-forward **CF-C-3** and discharges this row's stated
   "first deliverable". `store.go` and `schema.sql` are **byte-unchanged** — the fixture asserts
   current, WRONG behaviour on purpose and is expected to fail and be rewritten when the fix lands,
   so it pre-empts **none** of the three ratification arms below. **Non-vacuity PROVEN, not
   asserted**: under a temporary write-side validation in `store.Commit`, PASSing CFB2 subtests =
   **0**, FAILing = **20**; mutation reverted and `store.go` confirmed byte-identical
   (`sha256 ebaa5b00…`). Independently re-verified by the judge, mutation and all.

   **THE CORRECTED MATRIX (controller first-party, iter-27 — re-measured because a prior
   iteration's prose is a claim, not evidence).** All eight fields commit with `err=<nil>`; the
   damage splits into **THREE** classes, not two:
   - **CLASS 1 — log entry permanently unreadable (5 fields)**: `TransitionFn`, `Interpreter`,
     `PrevEntryHash`, `EntryHash`, `TransitionRef` → `GetLogEntry` returns `ok=false` + error;
     `GetWorld` and `SelectedHead` unaffected.
   - **CLASS 2 — world revision unloadable (2 fields)**: `NextWorld.StateRoot`, `NextWorld.LogHead`
     → `GetLogEntry` SUCCEEDS, but `GetWorld` returns `ok=false` + error while `SelectedHead`
     succeeds — so the selected head points at a world that cannot be loaded. Previously mislisted
     as "unreadable ROW"; it is a different read surface entirely.
   - **CLASS 3 — the store is WEDGED (1 field: `NextWorld.Ref`)**: log entry AND `GetWorld` both
     read back fine, but `SelectedHead()` **errors** (`store: selected head: hashref: empty hashref
     text`) and **every subsequent `Commit` fails with that same error, which is NOT a
     `ConflictError`** — so a caller's standard re-plan-on-conflict path never fires and the store
     can never accept another commit. Unrecoverable through the public API. The prior record called
     this field "degenerate-but-readable"; it is the **worst** of the eight.

   **DECISION-RELEVANT CONSEQUENCE for the packet above (an observation, NOT a decision).** A
   read-side accommodation (**ARM V2**) cannot address CLASS 3 at all: there is nothing to "support
   on read" when `SelectedHead` has no ref to return and the write path then refuses every later
   commit. The non-vacuity mutation independently showed the converse — one write-side check
   (**ARM V1**) reds all three classes at once. The asymmetry is now executable rather than argued.

   **Judge carry-forwards still open from the fixture (non-blocking; enumerated, per the iter-19
   process fix that a bare COUNT is unrecoverable):** **CF-E-3** — add store-untouched /
   head-advanced assertions for CLASS 1 and 2 to pin the blast radius; **CF-E-4** — note in
   `w-store-durability.md`'s SD.A section that the fixture landed at `e8ba7b2`, so a planner does
   not treat it as TODO; **CF-E-5** — clarify in-file why the zero-ref world is deliberately
   asserted *readable* in the CLASS-3 test. (**CF-E-1** doc-vs-repo filename mismatch and
   **CF-E-2** an under-discriminating error substring were both CLOSED in-PR.)]
   **w-store-durability** · clause-1 · store/kernel hardening across a crash:
   **(i) CF-B-2** — `store.Commit` writes a zero `PrevEntryHash` that `store.GetLogEntry` **cannot
   read back**, so the embedded Go API can append a log entry no reader can ever load. The daemon
   refuses it at the REST boundary, which is a boundary patch over a kernel defect. Discovered
   iter-21 **by a failing test**; ~~as of iter-23 it still has **no issue and no repro fixture**~~ —
   **both closed iter-27: issue #15 + `host/store/durability_repro_test.go` (`e8ba7b2`), CI-green.**
   ~~The remaining half of (i) is the FIX, which is what the ratification packet gates.~~
   **(i) IS COMPLETE — the FIX LANDED iter-28 at `86d1276` (SD.A, ARM V1).** What remains of this
   item is half (ii), the journal/receipts, which is SD.B + SD.C.
   **(ii) the broker effect journal** — a durable pre-dispatch intent object + journal head + a
   post-dispatch outcome object, with an `IndeterminateEffectError` recovery path that **never**
   auto-re-executes, and per-handler idempotency/reconciliation rules before any retry is allowed.
   Raised by `gpt5-6-sol` in the `w-effect-broker-m3` round-2 quorum: without it a real external
   effect (an `FS.Write`, a `Git.Commit`, a **paid** `Model.Infer`) can occur with no durable
   record, and replay cannot distinguish "never executed" from "executed but record lost".
   Both halves are the same question — **what the kernel guarantees across a crash** — which is why
   they share one item rather than being answered twice · ~1–2d · surfaced iter-21, queued iter-23.
4c. [**[IN-SPRINT — MJ.A + MJ.B LANDED; ONLY MJ.C REMAINS]**. **MJ.B LANDED 2026-07-29 (iter-38)**,
   PR #27 → squash `3ef5510`, dev CI green both jobs SHA-addressed on the merge commit **and the
   step logs read to confirm the gates actually executed** (both reported success in 36 s / 8 s —
   implausibly fast, so it was verified rather than believed), judge sonnet **PASS 86/100
   zero-blocking**, executor `codex:gpt-5.6-sol`, `metered=$0.00`. **The broker pipeline is now
   wired onto the effect journal**: `Invoke`'s allowed arm is `PutObject(request) →
   AppendNextEffectIntent → debit → dispatch → result object → putRecord → AppendEffectOutcome`,
   the denied arm journals nothing, the failed arm journals a resolved `"failed"` outcome, and
   replay journals nothing (asserted by a test, not by reading the code). `Recover` gains a second
   bounded page walk over pending EFFECT intents that **REPORTS ONLY**. Closes **AC7, AC9, AC10**;
   +547/−48 vs the planner's +585 = **0.94×**, the first milestone of this item to come in under
   estimate. **THREE PLANNER FACTS REFUTED AT GATE 2, BEFORE THE EXECUTOR RAN**: PD1's call-site
   census was wrong in both directions and **omitted the unexported `newSession` entirely**
   (measured: `NewSession` 10 sites not 15, `NewReplaySession` 6 not 5, plus 10 `newSession` sites
   the census never mentions); **OQ-5's premise was false** — it justifies its default by "the
   clock the broker already threads for records", and the broker threads none (`LogicalTime`
   appears zero times in `host/broker` production code), yet OQ-5 is a **STOP** instruction, so
   unexamined it would have killed the milestone at checkpoint one — resolved instead to
   `EffectIntent.LogicalTime = req.Now`, which is caller-supplied, in scope, and provably never a
   wall clock; and the plan's own `-race` verify command names a **known-red** gate. **THE SPINE:
   `git checkout <file>` IS NOT A MUTATION REVERT — it restores from HEAD, and every sprint here is
   delivered as an UNCOMMITTED diff, so it silently deleted the milestone's `broker.go` and left a
   green suite running on `origin/dev`'s code.** Caught only by the sha256 the mutation discipline
   already required; recovered from the executor's own logged diff and proven byte-identical by
   **both** `git hash-object` (`1fb60d78…`) and sha256 (`99f3568…`). **One mutation NAME, two FORMS,
   two answers — and this time the JUDGE had the confounded one** (`MUT-OUTCOME-BEFORE-RECORD`: the
   executor recomputed a valid ref and isolated the ordering → 2 red / 99 green; the judge passed a
   zero ref, so the store's ref validator fired first → 18 red). Fourth instance of that shape.
   **The judge's finding was reproduced and turned out UNDER-stated**: the milestone had deleted
   `recover.go`'s statement of the Decision-5 residual, against the doc freeze's explicit *"No claim
   … may state or imply that every crash ambiguity is eliminated"* — **fixed in-PR, not carried**.
   New **CF-MJB-1…CF-MJB-6** (CF-MJB-1 DONE; CF-MJB-6 REFUTED and recorded, not carried).
   **CF-MJA-4 and CF-MJA-5 were MJ.B-owned and are NOT discharged — they roll to MJ.C, stated.**
   ~~[IN-SPRINT — MJ.A LANDED 2026-07-29 (iter-37)]~~, PR #26 → squash `82d9128`, dev CI green
   both jobs SHA-addressed on the merge commit, judge sonnet **PASS 86/100 zero-blocking**,
   executor `codex:gpt-5.6-sol`, `metered=$0.00`. Sprint plan + handoff written by the **opus**
   planner at `.ailang/state/sprints/w-effect-journal.{plan.json,handoff.md}` (MJ.A/MJ.B/MJ.C).
   **MJ.A shipped the kernel reopen + the law** (+733/−10; planner predicted 615 = 1.19×): four new
   `host/store` methods (`AppendNextEffectIntent` — the ordinal minted INSIDE the appending
   transaction, verified by reading the code — `AppendEffectOutcome`, `GetEffectReceipt`,
   `PendingEffectIntents`), three changed (`validateIntent` rejects the reserved `effect:` prefix,
   `GetReceipt` gains the same boundary guard **closing OD1**, `PendingIntents` gains the
   `world/journal-intent/v1` predicate), and **LAW 5 `effectDispatchLawful`** in
   `design_docs/sketches/storejournal.ail` — `ai-check` **8/8 contracts `status=verified`** with
   z3 4.16.0 present, so not a V27 silent skip. Closes **AC1–AC6, AC7b, AC7c, AC8**;
   `schema.sql`/`store.go`/`host/broker/**`/`host/replay/**` byte-unchanged. **11 named mutations**,
   each match-asserted-exactly-once, compiled, diffed and reverted byte-identical; the judge
   independently re-ran two with sha256-proven reverts. **THE LOOP REFUTED ITSELF TWICE, both
   wins**: the planner refuted doc premise V28 (the real serialization point is
   `store.go:253 SetMaxOpenConns(1)`, not "no mutex") — which mattered, because under
   `MUT-ORDINAL-SPLIT-TX` the losing goroutine **errors** rather than returning a duplicate ID, so
   an AC7b test that only compared the two IDs would have been **vacuous under its own named
   mutation**; and the executor then refuted the planner's predicted MECHANISM (the loser returns
   `DuplicateInvocationError`, not the raw SQLite `UNIQUE` error, because compare-and-append
   duplicate discipline is preserved) — the gate prediction held, the stated reason did not.
   **AC5's third clause is STRUCK as a doc defect, not implemented**: measured first-party, landed
   `AppendIntent` is identical-bytes idempotent but landed `AppendOutcome` is not, and
   `AppendEffectOutcome` mirrors `AppendOutcome` — so "idempotent identical-bytes re-append"
   describes behaviour absent from the commit side too, and adding it to the effect side alone
   would diverge from the substrate. New **CF-MJA-1…CF-MJA-5** (CF-MJA-3 already DONE — PR #26
   carries the mutation table). ~~Remaining: MJ.B~~ **MJ.B LANDED iter-38 — see the tag above;
   OQ-5 was resolved to `req.Now` by the controller after its stated premise was refuted.**
   **Remaining: MJ.C only** (CF-N-2/CF-N-3 discharge + bench + close-out, now folding
   CF-MJA-1/CF-MJA-2/CF-MJA-4/CF-MJA-5 and CF-MJB-2/3/4/5).
   ~~[NEXT — DOC LANDED + QUORUM-CLEARED 2026-07-29 (iter-36)]~~ — doc `design_docs/planned/w-effect-journal.md` (770 lines,
   **30 premise rows**) shipped via **PR #25 → squash `fe582b5`**, dev CI green both jobs
   SHA-addressed on the merge commit. Designer `claude:claude-fable-5` (rotation, advanced from
   `codex:gpt-5.6-sol`). **TWO quorum rounds, both blocking-then-fixed**: r1 `gemini-3-1-pro` found
   the resumption-ordinal collision (a transient counter restarting at 0 after a crash GUARANTEES
   `DuplicateInvocationError` against the durable `effect:<episodeID>:0`); r2 `gpt5-6-sol` found
   that fixing it with a *split* derive-then-append merely moved the collision to a **TOCTOU**.
   r2's fix applied VERBATIM under the **narrow-refinement carve-out** (concrete reviewer-authored
   fix, determinism-only, direction endorsed): one transactional
   `AppendNextEffectIntent(episodeID, intentWithoutID) (id, ordinal, err)` mints the ordinal
   INSIDE the appending transaction, the broker holds no ordinal state, `MaxInt64` exhaustion gets
   a structured error, and the reviewer's two named tests land as **AC7b**. `gemini-3-1-pro`'s
   `GetReceipt` namespace guard adopted too (**OD1 CLOSED**). Kernel delta frozen at **four new**
   store methods + three changed; `schema.sql` byte-unchanged.
   **THE QUEUE ROW'S OWN COSTING CLAIM WAS REFUTED AT PICK** — see the iter-36 STATUS stamp: the
   M3 planner's *"Cheaper than assumed — no schema change"* is false for the DDL (`CHECK (kind IN
   ('intent','outcome'))` rejects a new kind; `UNIQUE (invocation_id, kind)` caps intents at one
   per ID), **any DDL change here ships fail-OPEN**, and the `sqlite_master` gate it cited is
   **inert**. Path A (DDL-free, reserved `effect:` ID namespace) is the design's answer.
   Still inherits **CF-N-2** and **CF-N-3**, both now acceptance criteria with their own
   mutations. It owns the **dispatch→record window**, which M3 left OPEN and stated. The
   `host/store` kernel reopen stays pre-ratified IN PRINCIPLE by `c26b27d`.]
   **w-effect-journal** · clause-3 · effect-shaped intents in `host/store` (payload-codec change,
   **no DDL — now a design CONSTRAINT actively satisfied, not a free discount**): closes the
   dispatch→record window at the exact granularity of the original 2A objection · **~1.5d**,
   3 milestones (MJ.A kernel reopen + law · MJ.B rewired pipeline + crash-window proof · MJ.C
   CF-N-2/CF-N-3 discharge + bench + close-out)
4d. [**[QUEUED 2026-07-29 (iter-36)]** — raised by the controller's own Gate-2 audit, with the
   evidence attached; NOT folded into 4c, because it is a pre-existing defect in LANDED
   `w-store-durability` code and 4c's AC1 deliberately does not depend on the broken gate.
   **Sized as a small item (~0.25–0.5d), not a sprint.** Pick it when a DDL change is next
   contemplated, or sooner if the queue allows.]
   **w-ddl-gate-teeth** · clause-1 · **the DDL-drift gate is inert where it is cited, and DDL
   changes ship fail-open.** Three compiling mutations, run first-party at `0f2afad` and reverted
   byte-identical: (1) **`MUT-JOURNAL-DDL-WIDEN`** — widening the journal `kind` CHECK leaves
   `go test ./...` **rc=0 across all 10 packages**, so *nothing in this repository guards the
   journal table's DDL*; (2) **`MUT-DDL-DRIFT`** — the gate's OWN named mutation does red, but by
   the **sha256 source pin** (`pre-journal schema source drifted`), which `t.Fatalf`s before a
   database exists, so the `sqlite_master` comparison never executes; (3) **`MUT-DDL-COMPARE-DEAD`**
   — short-circuiting that comparison (both vars still used, `go vet` rc=0) leaves the whole
   `host/store` package **green**, proving it contributes zero discrimination. Separately,
   `CREATE TABLE IF NOT EXISTS` silently drops a DDL edit on any store that already exists (rc=0,
   no error) and `host/` has **zero** migration machinery — so a schema change that was never
   applied is indistinguishable from one that was. **The work**: give the gate teeth over the
   journal table AND the upgrade-an-existing-store path, or retire the dead comparison and
   document the sha256 pin as the real mechanism — plus decide whether a `user_version` pin that
   fails LOUD on an un-upgraded store lands now or with the first item that needs a DDL change.
   **Fifth instance of this mission's signature shape** (a gate no production change could fail),
   and the first found at PICK time in a gate inherited from a COMPLETE, zero-blocking-judged item.
4e. [**[QUEUED 2026-07-29 (iter-37)]** — found by chasing a line the MJ.A executor mentioned in
   passing; **pre-existing in LANDED `w-store-durability` (SD.A, `86d1276`) code that two
   zero-blocking judgements never saw**. Sized small (~0.25–0.5d). Pick alongside 4d.]
   **w-race-gate-blindspot** · clause-1 · **a test that fails deterministically has sat in landed
   code for four milestones, because the gate that would see it is never run.** Measured on CLEAN
   `origin/dev` @ `d057de8`, MJ.A diff absent, `scan.go`/`scan_test.go` untouched:
   `TestScanUnreadableLogKeysetResumes` **passes** without `-race` (both `-run`-filtered and
   whole-package) and **fails 5/5** with `-race` (both). Not `-count`, not the `-run` filter, and
   **not cgo** (`CGO_ENABLED=1` is this rig's default and passes without `-race`). Exactly one
   element differs: `Rows[0].Field` is `""` under `-race` and `"prevEntryHash"` without it, while
   `Reason` from the SAME loop iteration is correct — and `fields` is a five-element string literal
   indexed over a `[5]string`, so an empty `Field` should be unreachable. **Zero `DATA RACE`
   warnings** are reported. Second symptom: the whole `host/store` package under `-race` stops
   progressing and is `signal: killed` at ~142 s. **RE-MEASURED first-party at iter-38 on clean
   `40ff563`** (a different base, with MJ.A now landed) — the failure reproduces identically, and
   the measurement is **sharpened in the direction that matters for scoping**: the sibling package
   `host/broker` under `-race` is **GREEN in 4.095 s**, so the defect is `host/store`-local and a
   `-race` leg is not globally unusable. Whole-package run: `ok host/broker 4.095s` /
   `FAIL host/store 116.552s (signal: killed)`. Also newly observed: `Rows[0].Ref` is empty
   alongside `Field`, consistent with the `fields` array reading as all-zero rather than one
   element being wrong — which would make it a memory-corruption signature, not a logic bug.
   Root enabler: **`grep -rn race
   .github/workflows/ scripts/` returns NOTHING** — neither CI job, nor `verify_go.sh`, nor
   `verify_ail.sh` ever passes `-race`. **Mechanism UNKNOWN and labelled so**; standing hypothesis
   (UNVERIFIED) is `modernc.org/sqlite`'s heavy `unsafe` usage plus `-race`'s altered memory layout
   exposing a latent corruption, or a go1.26.4 darwin/arm64 toolchain bug. **The work**: identify
   the mechanism before prescribing the fix — because zero `DATA RACE` output means "just add
   `-race` to CI" is not yet known-good — then either add a `-race` leg with the hang bounded, or
   record in writing why this repo does not run one. **Sixth instance of this mission's signature
   shape** (a gate that cannot fail), and the first where the gate cannot fail *because it is never
   invoked at all*.
5. [**BLOCKED 2026-07-28 (iter-24) — DOC LANDED + 2 QUORUM ROUNDS + carve-out revision applied.**
   Doc: `design_docs/planned/w-mcp-projection.md` (codex/gpt-5.6-sol designer, rotation; 641 lines;
   NO `.ail` sketch — protocol/session invariants are host-boundary behaviour, so the required-check
   manifest is untouched at 4/9/14). **Milestone P6.A is DONE this iteration** (upstream finding
   filed + this record landed). **P6.B is blocked on THREE named prerequisites** and needs no
   re-design when they clear] **w-mcp-projection** · clause-6 · project the transition registry over
   MCP + publish the A2A agent card (reuse `ailang serve-api --mcp/--a2a` machinery — do not
   reinvent) · ~1.2d of World work *after* prerequisites

   **THE QUEUE ROW'S PREMISE DID NOT SURVIVE THE BINARY.** "Reuse the serve-api machinery" is
   viable for *protocol projection of static exports* (live-tested 2026-07-23) and **not** viable
   for the acceptance half — dynamic worldd-backed registry, per-session capability filtering,
   propose→verify→commit enforcement. Verified first-party on pinned v0.30.0 (controller's own
   stdio MCP probe, iter-24): `unfiltered → ['addOne','submit_feedback']`,
   `--routes-only → ['submit_feedback']`, `--caps '' → ['addOne','submit_feedback']`.
   So **(i)** a built-in `submit_feedback` tool is exposed under EVERY flag combination with no way
   to remove it — and its own description routes to a `public-feedback` inbox with a Pub/Sub
   notification, i.e. a built-in **egress** tool no World session authorized, which collides head-on
   with clause 2 (zero cloud deps in the core) and clause 3 (no ambient-authority path from an agent
   to the outside world); **(ii)** discovery is built from loaded module exports, not a
   caller-supplied provider, and `--caps` / `--routes-only` / `--api-key-*` are all **process-wide**,
   so neither a static key nor a process cap set can represent a session. `--routes-only` DOES
   suppress the 8 embedded `std/io` exports, so upstream `#145` is genuinely fixed and this is not
   its regression. Reuse paths **(a) and (b) are REJECTED ON EVIDENCE** (a sidecar-per-session makes
   process lifetime the session model and still exposes `submit_feedback`; reverse-proxy filtering
   would make World parse and re-encode MCP/A2A — the reinvention §3.7 forbids). The design takes
   **path (c)**: a narrow public serving seam over the existing `internal/apiserver`.

   **The three P6.B prerequisites:**
   1. **Upstream seam** — `sunholo-data/ailang#498` (filed iter-24, both channels, with the stdio
      repro and a labelled-hypothesis cause): export the existing MCP/A2A serving machinery behind a
      public callback-driven API — caller-owned mux, principal resolved before discovery *or*
      invocation, caller supplies the exact visible descriptors, **no built-in tool unless the
      caller supplies it**. A narrower interim fix would unblock much of this: a flag suppressing
      `submit_feedback`, plus documenting that `--caps` gates execution rather than discovery.
   2. **Clause-3 transition registry + broker session API** — VERIFIED ABSENT (iter-24: a repo-wide
      search for `[Tt]ransition[ -]?[Rr]egistry` matches `design_docs/` ONLY — zero hits in `host/`,
      `world/`, `cmd/`; `host/registry` is the *interpreter epoch* registry, a different thing).
      Gated behind `w-effect-broker-m3` (item 4, PARKED).
   3. **A verified commit-boundary contract** — raised by `gpt5-6-sol` in quorum round 2 and applied
      under the carve-out: an atomic "not-started versus committed" contract, a stable
      invocation/idempotency ID, and a queryable durable receipt. No landed API exposes these.
      **This is the SAME kernel-durability question as item 4b half (ii)**, reached from a different
      direction — independent corroboration that 4b is real work, not bookkeeping.
6. [PARKED until 2–5 land] **w-agent-floor-m4** · clause-4 · dual-reference NON-INFERIORITY
   floor: Claude Code + codex, shell arm vs World-MCP arm, paired N≥3, stability precondition
   checked first; motoko as optional third arm if eligible; report honestly; park World if the
   floor fails on eligible agents · ~3d (was w-motoko-m4 → w-value-gate-m4 → renamed with the
   2026-07-23 floor reframe; clause-5 provenance teeth carry the value burden)
6b. [**PARKED `needs-human-review` 2026-07-28 (iter-26) — PICK-TIME QUORUM COMPLETE (2 rounds,
   4 objections, all applied); the only thing left is Mark's §7 ratification.** Doc revised
   `170 → 386` lines: **v0.2** by the rotation designer (`codex:gpt-5.6-sol`) applying round-1's
   two verbatim fixes — a **16-row Premise Verification Log** (§8), **§8.1 reuse-or-replace**
   across all six categories the reviewer named, **§6.1 Hub-vs-workbench Conflict Surface**, §1
   recast from assertion to requirement, dependency notes on P3/P5/P7 and §6; **v0.3** by the
   controller under the narrow-refinement carve-out applying round-2's verbatim fixes. Quorum
   artifacts `human-surface-2026-07-28T03-54-41Z.json` / `…T04-09-06Z.json`; metered $0.089226.
   **Unparks straight to `w-approval-inbox` (item 7) once §7 is ratified — no re-design, no
   re-quorum needed.**

   **WHY IT WAS BLOCKED AT ALL — the doc was authored ATTENDED, so it never passed through
   `design-doc-creator`'s hard gates.** Both round-1 objections were exactly the two sections that
   skill would have forced (Premise Verification Log; Conflict Surface) — the SAME pair that
   blocked the charter itself at iter-0. Now routed as a charter guardrail (see Guardrails).

   **ROUND 2 FOUND SOMETHING REAL, TWICE, INDEPENDENTLY.** `gpt5-6-sol` and `gemini-3-1-pro` —
   different providers, no shared context — both raised the same new objection: §3's *defer*
   ("park for more evidence") and P5's *batch over interrupt* admitted an **unbounded wait on a
   human**, with no TTL, no expiry transition, and no deterministic outcome if the human never
   answers. That is this mission's own **Standing Rule 6** ("every wait is bounded"), missing at
   the UX layer of the doc that governs every human-facing surface. Applied verbatim as §3 (TTL →
   typed rejected-timeout) and **§3.1 Bounded decision lifecycle** (ledger-recorded creation time
   + deadline + typed timeout policy; DEFER must rebound and must not park indefinitely; an
   explicit Timeout transition; *"Silence MUST never synthesize approval or rejection"*; replay
   reproduces deadlines from ledger time, not wall-clock races).

   **WHAT THE CONTROLLER MEASURED FIRST-PARTY (the load-bearing claim did not survive the kernel).**
   The doc's trust gradient has **four** grades (PROVEN > TESTED > ATTESTED > CLAIMED); the landed
   kernel's `Evidence` ADT (`world/types.ail:23-28`) has **five** variants. `TestReport`→TESTED,
   `RecordedEffect`→ATTESTED, `AiReview`→CLAIMED — but **`CompilerOutput` and `HumanApproval` have
   no grade**, and **PROVEN's own stated producers (Z3, replay) have no `Evidence` carrier at
   all**. The grade names appear **nowhere** in `*.ail`/`*.go`/`*.sql` (one prose comment aside).
   `HumanApproval` is precisely the evidence the approval inbox — item 7, gated on this doc —
   would emit, and grading a human's ratification as CLAIMED ("agent said so, unverified") is
   plainly wrong. So the doc's cardinal-sin anti-pattern, **grade laundering, is today
   unenforceable**: there is no total mapping to launder against. Ratification point 7.2 is
   restated as that decidable question. Also measured: `Proposal.confidence` is a **bare float
   with no evidence ref** (`AiReview` carries one), so rendering it would violate the doc's own
   confidence-theater anti-pattern → new ratification point 7.5]
   **w-human-surface** · clause-5 · the founding UX design for the human↔AI surface:
   [HUMAN-SURFACE.md](HUMAN-SURFACE.md) — seven interaction principles (decision packets ·
   grounded prose · trust-gradient rendering PROVEN/TESTED/ATTESTED/CLAIMED · time-as-navigation
   · attention choreography · five zoom levels · speculation-as-gesture) + sign-the-type input
   grammar + anti-patterns (grade laundering is the cardinal sin). Mark named this the vision's
   key surface ("an ideal AI language in an AI state machine OS with a NOVEL human interaction
   surface"); no human-facing sprint routes before it is ratified · ratification session ~0.5d
7. [PARKED until 4 AND 6b land] **w-approval-inbox** · clause-5 · the approval inbox + provenance
   walk, first as CLI/generated projection (SCENARIOS.md scenario 1/3), **built to
   HUMAN-SURFACE.md** · ~2d
8. [PARKED until 4 lands] **w-self-mod-vertical** · clause-7 · one World extension shipped through
   World's own pipeline end-to-end · ~1–2d
9. [BACKLOG — infra, not critical-path; **SHARPENED WITH FIRST-PARTY MEASUREMENT iter-27**]
   **w-verify-binary-lockfile** · clause-1-infra · a durable
   pin of the released `ailang` binary (version + sha256) for the `ailang-code` verify gate, so
   local controller verification stops depending on the rig's `-dirty` build (iter-6 pinned
   released `v0.30.0` ad-hoc — `sha256:ac3174e0…`; CI already sha256-verifies a released linux
   binary). Small; the mechanism may generalize to the SHARED `ailang-code` profile → confirm
   repo-local lockfile vs shared-skill fix with a human before implementing (do not hand-edit CI
   headless) · ~0.5d · surfaced iter-6, queued iter-7

   **MEASURED iter-27 (controller first-party, during the Gate-2 reality check). The two sibling
   gates are ASYMMETRIC, and the `.ail` one is the silent side.**
   - `scripts/verify_go.sh:19-32` **hard-fails loudly** if `AILANG_BIN` is unset or its binary does
     not report `v0.30.0` — the anti-false-green guard M6 landed to close the V27/B1 silent-skip
     class.
   - `scripts/verify_ail.sh:33` is `AILANG_BIN="${AILANG_BIN:-ailang}"` — **no guard, no version
     assertion, and it never prints which binary it used.** On this rig bare PATH `ailang` is
     **`v0.30.0-205-g54d6bd191-dirty`**, i.e. 205 commits past the pin and dirty. So running the
     repo's primary gate without exporting `AILANG_BIN` silently validates against a dev build,
     violating CLAUDE.md's own hard rule ("never a `-dirty` dev build").
   - `.github/workflows/ci.yml:18` — the `ailang-verify` job installs
     `releases/**latest**/download`, **not** the pinned tag, and **never asserts the version**
     (`:43-44` merely prints it). Contrast `:71` where `go-verify` pins
     `releases/download/**v0.30.0**` + sha256 + `grep -q 'v0.30.0'`. Today `latest` **IS** v0.30.0
     (checked: published 2026-07-19), so CI is correct **by coincidence, not by pin** — it will
     silently start verifying against a different compiler the day v0.31.0 ships.
   - **Severity honestly bounded**: `verify_ail.sh` was run BOTH ways this iteration and the output
     was **byte-identical** (4/4 identities / 10 modules / 14 named tests, rc=0 each). So no prior
     iteration's recorded `.ail` evidence is invalidated — this is a **latent** false-green surface
     and a real attribution gap, NOT an active wrong result. Recorded because the next release turns
     it active with no signal.
   - **Decomposition for the human gate**: the *cheap, zero-risk* half is making `verify_ail.sh`
     ANNOUNCE its resolved binary + version the way `verify_go.sh:33` already does (pure
     observability, cannot red anything). The *human-gated* half is the hard version assertion plus
     the CI `latest`→pinned-tag edit — those two are **coupled** (a hard assert alone would red CI on
     the next upstream release, headless), which is exactly why this row's "confirm before
     implementing / do not hand-edit CI headless" flag was right and was respected this iteration.

## Premise Verification Log (quorum objection #1 — every load-bearing claim, with evidence)

| Premise | Verified | Evidence |
|---|---|---|
| `serve-api --mcp / --mcp-http / --a2a` exist (clause-6 machinery) | bootstrap + iter-0 | `ailang serve-api` help, v0.30.0 |
| `ai-check` = unified check+verify, always-JSON, **no `--json` flag** | bootstrap flag-test + iter-0 | error repro'd live; citations fixed here (`4c6080a`) and upstream (`ailang@aabb3a58c`, v1-ack on issue #1) |
| CI `CI` / job `ailang-code verify gate` green and NON-vacuous | every push since `3c8791d` | job log shows `checked 2 module(s)`, released linux binary, sha256-verified |
| DESIGN.md type sketches compile | CI-gated | `design_docs/sketches/*.ail` in every run |
| Mission state fully namespaced (no v1 collision) | dry-runs 2026-07-23 | `/tmp/ailang-mission-world.log`: distinct `mission-world.pid`, worldtest profile proof |
| Upstream routing channel works end-to-end | 2026-07-23 | defect report → v1 verified+fixed+ack'd < 1h (issue #1) |
| Registry + `ailang publish` cascade exist (clause-7 lane) | v1-mission operational history | live on multivac; World's local-first cascade mode is DESIGN.md §13/M-scope work, not assumed |
| `internal/apiserver` (serve-api impl) is stateless — no persistence, no scheduler | code inspection 2026-07-23 | package-wide grep: zero sqlite/sql.Open/bolt/badger/scheduler hits; request-scoped handlers only |
| `internal/apiserver` NOT importable cross-repo (Go `internal/`) | Go module rules + path | Conflict Surface reuse paths (a)/(b)/(c) sized accordingly; path (a) needs no upstream change |
| serve-api projects `.ail` exports as MCP tools (path (a), static case) | LIVE TEST 2026-07-23; **cwd SHARPENED iter-24** | `ailang serve-api --mcp-http --port 8199 sketches` → `tools/list` returned `plan`/`verify`/`commit` with JSON schemas + effect rows in descriptions; server killed after, port freed. Dynamic/capability-filtered projection deliberately NOT claimed — w-mcp-projection acceptance criteria. **iter-24 correction: this row does not record its cwd, and cwd is load-bearing** — serving `design_docs/sketches/` from the repo ROOT fails `LDR001: module not found: sketches/worldtypes` (imports resolve relative to cwd), while serving `sketches/` from `design_docs/` succeeds. serve-api also runs the loader in `--relax-modules` mode (MOD010 warnings ignored). So this row is insufficient as launch configuration for a sidecar |
| **The projected MCP/A2A surface CANNOT be an exact allowlist (clause-6 blocker)** | **LIVE TEST 2026-07-28 (iter-24), controller first-party, pinned v0.30.0** | Own stdio MCP probe (`initialize` → `notifications/initialized` → `tools/list`): `unfiltered → ['addOne','submit_feedback']`, `--routes-only → ['submit_feedback']`, `--caps '' → ['addOne','submit_feedback']`. A built-in **`submit_feedback`** tool survives EVERY flag combination; its own captured description routes to a `public-feedback` inbox with a **Pub/Sub notification** — a built-in egress tool no World session authorized (collides with clause 2 zero-cloud and clause 3 no-ambient-authority). `--caps` gates execution, **not discovery**. `--routes-only` DOES suppress the 8 embedded `std/io` exports (`eprintln`/`exit`/`flush`/`print`/`printErr`/`println`/`readLine`/`writeBytes`), so upstream `#145` is fixed and this is not its regression. MCP HTTP at `/mcp/` replies **SSE-framed** (`event: message` + `data:`), not plain JSON; `--a2a` `/.well-known/agent.json` = 200 with `skills[]` mirroring the same unfiltered set. Routed upstream as `sunholo-data/ailang#498` (cause = labelled HYPOTHESIS; upstream source not inspected) |
| **No transition registry exists in this repo** (clause-3 prerequisite is real) | **iter-24, controller first-party** | Repo-wide search for `[Tt]ransition[ -]?[Rr]egistry` matches `design_docs/` ONLY — **zero** hits in `host/`, `world/`, `cmd/`. `host/registry` is the *interpreter epoch* registry (`world/epoch-registry/v1`), which is interpreter-nomination metadata, not the transition registry clause 6 must project |
| **The daemon enforces global HTTP write deadlines** (frozen D7 block) | **iter-24, controller first-party** | `host/daemon/daemon.go:409-414` wires constants declared at `:77-91`: `ReadHeaderTimeout` 5s, `ReadTimeout` 30s, `WriteTimeout` 30s, `IdleTimeout` 120s. Consequence: a long-lived SSE stream mounted on this server is killed at the 30 s write deadline unless that ONE route relaxes it via `http.ResponseController` — and `go.mod:3` = `go 1.26.4`, so that API is available. `IdleTimeout` is NOT a stream-lifetime bound (it governs idleness between requests) |
| `verify_ail.sh` fails loudly at N=0 (gate cannot pass vacuously) | LIVE TEST 2026-07-23 | script run against an empty `design_docs/` scratch tree → "✗ no .ail modules found — the gate would be vacuous; failing loudly", **exit code 1** |
| **CF-B-2 damage is THREE classes, and the previously-recorded "seven unreadable / one benign" split is WRONG** | **iter-27, controller first-party, then re-verified independently by the sonnet judge** | Eight-field matrix on `store.Commit` in `host/store`: all 8 zero-ref commits return `err=<nil>`. CLASS 1 (5 fields) → `GetLogEntry` `ok=false`+err. CLASS 2 (`NextWorld.StateRoot`, `NextWorld.LogHead`) → `GetLogEntry` OK but `GetWorld` `ok=false`+err, head still selectable. CLASS 3 (`NextWorld.Ref`) → entry+world read fine, `SelectedHead()` **errors**, and **every later `Commit` fails with a non-`ConflictError`** ⇒ store unrecoverably wedged. Now executable as `host/store/durability_repro_test.go` (`e8ba7b2`), **non-vacuity proven** (write-side validation mutation ⇒ 0 PASS / 20 FAIL; reverted, `store.go` byte-identical `sha256 ebaa5b00…`) |
| **`verify_ail.sh` does not pin or announce its binary; CI's `ailang-verify` tracks `releases/latest`, not the v0.30.0 tag** | **iter-27, controller first-party** | `scripts/verify_ail.sh:33` = `AILANG_BIN="${AILANG_BIN:-ailang}"` (no guard, no version assert, never prints the binary) vs `scripts/verify_go.sh:19-32` which hard-fails unless `AILANG_BIN` reports `v0.30.0`. `ci.yml:18` installs `releases/latest/download`; `:71` (go-verify) pins `releases/download/v0.30.0` + sha256 + `grep -q`. Rig PATH `ailang` = `v0.30.0-205-g54d6bd191-dirty`. **Latent, not active**: gate output was byte-identical pinned vs dirty (rc=0, 4/4 identities, 14 named tests both ways), and `latest` IS v0.30.0 today — so CI is right by coincidence. Routed to queue item 9 |
| **CF-B-2 is CLOSED at the kernel write path; `Commit` validates EIGHT ref fields, not seven** | **iter-28, controller first-party + judge-reproduced** | ARM V1 landed at `86d1276` (PR #17). `store.Commit`'s ref table carries all 8 (`Entry.Header.TransitionFn`/`Interpreter`/`PrevEntryHash`, `Entry.EntryHash`, `Entry.TransitionRef`, `NextWorld.StateRoot`/`LogHead`/`Ref`) plus each Object's `Hash`+`InterfaceHash`, all validated BEFORE `s.db.Begin()`. `ObservedHead` is deliberately absent (zero-legal at genesis; `TestGenesisObservedHeadRemainsZeroLegal` pins it). **Non-vacuity measured**: `MUT-SEVEN-NOT-EIGHT` (delete the `NextWorld.Ref` row = the doc's original miscount) → **2 FAIL** incl. `TestCFB2ZeroWorldRefWedgeRejected`; revert byte-identical `sha256 23814a56…`. Four binding gates green: `go test ./...` · `verify_go.sh` · `verify_ail.sh` (4/4 identities, 10 modules, 14 named tests) · `bench_worldd.sh --smoke` (6 benchmarks) |
| **The codex `workspace-write` sandbox DENIES loopback binds, so socket-dependent gate results from it are uninformative — and can MASK a real defect** | **iter-28, controller first-party** | Every `httptest`/live-daemon test in the executor's own run aborted with `panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: bind: operation not permitted`, and it honestly reported the failure as sandbox-caused. Re-running the SAME gate outside the sandbox produced a DIFFERENT failure — `GET /v1/health: context deadline exceeded` — a real deadlock: `announce` is often an `io.Pipe` (synchronous, unbuffered) whose consumers read exactly ONE line, so the new sweep lines blocked `Run` before `Serve()`. Consequence for routing: **never accept an executor's gate verdict for `host/daemon` or `cmd/ailang-worldd`** — the controller must re-run outside the sandbox. Proposed upstream as codex-recipe false-green **#3** (World cannot edit the shared skill) |
| Driver sources `~/.config/ailang/mission-world.env` + respects role overrides | code + LIVE 2026-07-23 | `tools/launchd/mission-control.sh:46-47` sources `mission-${MISSION_PROFILE}.env`; `:238` `MISSION_EXECUTOR_MODEL` respected (default opus); dry-run log 19:45:39 echoes env-only values (repo-slug/doc/workdir) + resolved roles — sourcing proven end-to-end |
| Charter RATIFIED (authorization state, kill switch, sprint routing) | attended session 2026-07-23 | Mark's decisions recorded: clause-4 = −2pp/≤25% paired N≥3; bar+Conflict Surface+guardrails+queue as drafted; ratification comment on issue #1; STATUS stamp above is the in-doc record |
| Kill-switch + gh-issue state paths are world-namespaced in the LIVE driver (safety property) | LIVE 2026-07-23 | driver log 19:45:24/19:45:39/20:02:47: "kill switch present (`…/mission-world.disabled`) — skip" — the driver printed and HONORED the world-namespaced path 3× while the v1 loop ran unaffected; iter-0 posted to issue #1 = `mission-world-gh-issue` read correctly |
| `ailang messages` channel works end-to-end (guardrail's delivery leg) | live round-trip 2026-07-23 | sent `msg_…_2c6964d3` (defect report) + `msg_…_acc5edcc` (channel test); v1 agent RECEIVED and acted — upstream ack on issue #1 at 18:27Z citing the report, fix `ailang@aabb3a58c` |
| v1 session-start hook reads the message inbox | config + observed 2026-07-23 | ailang repo `.claude/settings.json` SessionStart → `scripts/hooks/session_start.sh` ("checks the user inbox … using the ailang messages CLI"); displayed 5 unread at this session's start |

---
**Document created**: 2026-07-23 (bootstrap, attended). **RATIFIED 2026-07-23** (iteration 0,
attended: Mark + World coordinator) — record on issue #1; advisory-quorum ledger in
`.ailang/state/mission-quorum/`. Sprint routing is authorized from the next loop fire.

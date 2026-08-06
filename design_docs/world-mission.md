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
- **THIS LOOP KEEPS *TWO* MARK-COMMENT WATERMARKS, AND ONLY ONE OF THEM IS THE ONE THE SHARED
  SKILL DERIVES (process fix, iter-52).** Gate 0 prescribes the issue-scoped
  `~/.ailang/state/mission-${MISSION_GH_ISSUE}-last-seen`; this mission has in fact been writing
  the mission-scoped `~/.ailang/state/mission-world-last-seen`. Because the issue number ROTATES
  WEEKLY, the skill-derived path is a *fresh, empty* file most weeks — reading it returns the epoch
  default, so the Gate-0 query re-reads all history. That direction is harmless (re-triage is
  idempotent). **The mirror is not**: a stale issue-scoped file whose number happens to be reused,
  or a mission-scoped file left unwritten, can make an unprocessed `MarkEdmondson1234` comment
  invisible — and a human directive OUTRANKS the queue, so dropping one is the highest-severity
  failure available to this loop. **Rule: read BOTH files and take the OLDER of the two as the
  watermark; write BOTH after triage.** Iteration 52 found the pair 24 h apart
  (`mission-world-last-seen` = `2026-08-04T08:25:01Z`, `mission-32-last-seen` absent) and caught
  Mark's `8/OD-1` ratification only because the empty file failed safe. This is World-local (V1's
  issue-scoped path is correct for V1), so it is a process fix here rather than a skill proposal.
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
  (c) **`"$var:path"` IS A ZSH HISTORY-MODIFIER EXPANSION, NOT A LITERAL — AND `git show
  "$rev:host/…"` IS SILENTLY REWRITTEN (added iter-41; instance 3 of this same class).** In zsh
  5.9 the `:x` suffix after an unbraced `$var` inside double quotes is applied as a **modifier**:
  measured with a bash control on identical strings, `"$c:host/x"` → **`.ost/x`** (`:h` = dirname),
  `"$c:tail/x"` → `abc123ail/x` (`:t`), `"$c:runtime/x"` → `abc123untime/x` (`:r`), `"$c:extra/x"`
  → `xtra/x` (`:e`), and `"$c:store/x"` is a hard `bad substitution` — while **bash leaves all five
  literal**. So `git show "$c:host/store/schema.sql"` reads `.ost/store/schema.sql` and fails, and
  iter-41 read **`total_tables=0` for the very commit that created the schema**. Three reasons this
  one is nastier than (a) and (b): `git show "$rev:path"` is **the** canonical archaeology idiom and
  the shared skill's own Gate 1 prescribes that exact form for reading mission state from origin;
  **every Go file in this repository lives under `host/`**, a modifier letter; and with stderr
  redirected — or piped into `grep -c` — the failure arrives as a **plausible zero**, not an error.
  **Always brace it: `"${c}:host/…"`.** Scoped honestly at the time of discovery: no committed
  script in this repo uses `git show` at all (control-verified), so this is a loop-instrument defect
  rather than a defect in landed code. **PROPOSED UPSTREAM** with (a) and (b), since Gate 1 of the
  shared skill teaches the vulnerable form.
  (d) **ZSH ARRAYS ARE 1-INDEXED; BASH ARRAYS ARE 0-INDEXED (added iter-45; instance 4 of this
  class).** `titles=( "" a b c d )` puts `""` at `titles[1]`, so a `for k in 1 2 3 4` loop reading
  `${titles[$k]}` — the natural bash-shaped idiom, and correct under bash — is **off by one** here.
  Iteration 45 built four reconstruction commits whose *content* was right and whose **messages were
  all shifted by one**, the first commit getting an empty subject. It fails silently and reads as a
  typo rather than a shell difference. Prefer an explicit `case` over array indexing in this loop.
  Cheap to catch — read back `git log --oneline` after any scripted commit sequence — and cheap to
  repair only because the tree content was proven identical before the four commits were discarded
  and rebuilt.
  **The general rule all four instances teach: a remedy is an instrument too, and inherits the same
  burden of proof as the thing it verifies.** Pair every negative or exit-code check with a
  known-positive control in the SAME call, so a broken instrument is distinguishable from a real
  all-clear. **Routed upstream as a PROPOSED shared-skill fix** (World cannot edit the
  mission-control SKILL.md — it lives in the V1 checkout), since Gate 2's step 3 currently
  *prescribes* the bash-only form.
- **READING A CI LOG IS AN INSTRUMENT, AND A LESSON RECORDED ONLY IN A STATUS STAMP IS ON A
  THREE-ITERATION TIMER (process fix, iter-42 — and the finding is the DECAY, not the defect).**
  Gate 3b's "verify rather than read" step means grepping a ~70 KB CI log for proof the gates ran.
  Two idioms silently return a plausible **zero** there, and both have now cost an iteration:
  (a) **`go test` prints `ok` + TWO SPACES + a TAB before the package path**, so `grep -cE "ok +github"`
  (iter-41) and `grep -cE "ok(\t| +)github"` (iter-42) BOTH return `0` for a run in which every
  package passed. Neither alternation matches `ok␣␣→github`. Use a shape that cannot care —
  `grep -oE "ok[[:space:]]+github\.com/[a-z0-9./_-]+" | sort -u | wc -l` — and pair it with a
  `FAIL` count, which is the reading that actually matters. True value both times: **10 packages**.
  (b) **`head -N` on a search is a TRUNCATION, not a result.** At iter-42 a
  `grep -rn bench_worldd … | head -20` cut off exactly above the `.github/workflows/ci.yml` hit and
  the controller was one sentence from recording *"the bench smoke gate runs nowhere in CI"* — the
  sibling loop's iteration-119 fabricated-absence defect, reproduced. `grep -rn` output order is
  filesystem traversal order, **not** sorted, so "the important hit would have been near the top" is
  not a defence. Count first (`grep -rc`), or read the file the claim is about.
  **What makes this a charter rule rather than another war story: iter-41 DID record (a) — inside its
  STATUS stamp.** The STATUS block keeps the newest **3** stamps and archives the rest, so a lesson
  written there is deleted from the charter's live text within three iterations, and (a) recurred at
  iter-42 while iter-41's stamp was still present but no longer where anyone looks. **A durable rule
  goes in this section; a STATUS stamp is a narrative of one iteration and expires like one.** Before
  writing a lesson into a stamp, ask whether it must survive the rotation — if yes, it belongs here.
- **A WRITTEN RULE IS NOT A CONTROL: THIS CHARTER ALREADY FORBADE `git checkout` AS A MUTATION
  REVERT, IN TWO PLACES, AND I DID IT ANYWAY FOUR ITERATIONS LATER (process fix, iter-38 — and the
  finding is the RECURRENCE, not the defect).** At iter-34 the controller reverted a mutation with
  `git checkout -- <file>` while that file carried uncommitted executor work and destroyed it; the
  lesson was written into this charter **twice** (the mutation-protocol step 5, and the M3.C
  second-limb paragraph that concludes *"the instrument must not be able to delete the thing it
  measures"*). At **iter-38 the identical mistake recurred** — same idiom, same class of file,
  reproducing `MUT-INTENT-AFTER-DISPATCH` on `host/broker/broker.go`. It is the same mechanism each
  time: every sprint here lands as an **UNCOMMITTED worktree diff** (the codex `workspace-write`
  sandbox cannot reach the linked worktree's `.git`), so `git checkout` restores from **HEAD** and
  discards the milestone's version of the file, leaving a suite **green on `origin/dev`'s code**
  with no error, no diff noise and no hint in the output.
  **What this second instance actually teaches, and what the first one did not:**
  (a) **The prose rule did not hold — the MECHANICAL check did, both times.** At iter-34 and again
  at iter-38, the thing that caught it was the `sha256` comparison the protocol demands for an
  unrelated reason. Treat the hash pair as the control and the prose as a reminder, never the
  reverse. **Print both values every time; a mutation record without them is unverified.**
  (b) **A rule that has to be recalled at the moment of use will eventually not be** — which is why
  the durable version of this fix belongs in the *command*, not the charter: take the backup in the
  same call that applies the mutation (`cp f /tmp/f.bak && python3 - <<'PY' … PY`), so the revert
  path exists before the mutation does and cannot be forgotten separately from it.
  (c) Recovery, if it happens anyway: codex logs full `git diff` dumps, so the file is
  reconstructible — but the restoration must be proven by **two independent hashes**
  (`git hash-object` matching the diff's recorded post-image blob id, AND sha256 matching the
  pre-mutation baseline), never by "it applied cleanly".
  (d) Before writing "new finding" in any Gate-4 record, **grep this charter for the rule you are
  about to invent** — iter-34 already recorded that check, and this bullet exists because the Gate-5
  retro caught the over-claim only after the log entry, STATUS stamp and public report had shipped
  it as novel. **Retract at the same volume you published.**
  **Corroborating instance, same iteration, same family**: a `grep -c -- "--- SKIP"` over NON-verbose
  `go test` output reported "ZERO skips" — `--- SKIP` lines print only under `-v`, so the zero was
  structural and could never have been anything else (truth under `-v`: 2 benign subprocess-helper
  skips). **A verification step is code, and code that cannot fail is not a check.**
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
  **RECURRED AT ITER-39, AND THE RECURRENCE IS THE INSTRUCTIVE PART: THE SAME DIRECTIVE BOTH
  FOLLOWED THIS RULE AND VIOLATED IT.** The iter-39 directive quoted every binding requirement
  inline — the prescribed remedy, and the reason MJ.C shipped with its four named mutations intact —
  while ALSO opening with *"Source of truth: `.ailang/state/sprints/w-effect-journal.plan.json` …
  Read both"*, pointing the executor at a file this very bullet says is structurally absent. The
  executor filed it as deviation 1 and proceeded correctly from the design doc. So the cost was one
  confused deviation report rather than a lost milestone, and **the half that saved it was the
  mechanical half (the requirements were physically in the prompt), not the remembered half.** Same
  shape as the `git checkout` recurrence below: *a rule that must be recalled at the moment of use
  will eventually be half-recalled.* The durable form is to stop naming the plan path in directives
  at all — quote the requirements and cite the design doc, which IS in the worktree.
  **Second-order note, and the reason this bullet exists twice over:** iter-39 was about to file
  this as a NEW systemic finding and caught it only by running iter-38's rule (d) — *grep this
  charter for the rule you are about to invent* — at Gate 5, BEFORE the log entry, STATUS stamp and
  public report shipped. Iter-38 ran the same check one gate too late and had to retract in public.
  **The check works, and it works only if it runs before publication, not after.**
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
- **A STALENESS SWEEP MUST RUN FROM THE DOC'S *OLDEST* DECLARED MEASUREMENT BASE — SWEEPING THE
  ROWS SOMEONE ALREADY NAMED IS NOT A SWEEP (process fix, iter-48).** Two frictions, same gap.
  Gate-2 rule 3b(vi) says to diff a document's Verification Log against its claims, but names no
  **base**, and the base is the whole instrument. Iteration 47 found `P9`/`P22` stale in
  `w-bench-load-confound` and left `4f/CF-R-1` to repair them. Iteration 48 repaired those two and
  declared the class closed — and the **planner refuted it**: `P6` quoted
  `bench/BASELINE.md:18` as `go1.26.4` when it had read `go1.25.6` since **`f19acac`**, *the same
  commit that staled the other two*. One drift event, three faces, two found — because the sweep
  checked the NAMES iteration 47 supplied rather than the CAUSE. Measured, and this is the part
  that matters: a long doc declares **several** measurement bases in its Verification-Log header
  (this one declared `c1e6125` for rows P1–P22 and `61348b9` for P26+), and
  `git diff --name-only <newer-base>..HEAD -- ':!design_docs'` returns **ZERO files** — an
  all-clear — while the same command from the **oldest** base returns **8**. So the natural
  instrument, diffing from whichever base you happened to read, produces a confident clean bill of
  health on a genuinely stale document. **The rule: parse every base the header declares, sweep
  from the EARLIEST, and treat a row's verdict as unverified whenever the diff touches a file that
  row cites.** Pair it with a control (a file you know changed) so an empty diff proves the
  instrument ran. Corollary, and it is the general form: **a document is only as fresh as its
  oldest measurement**, so a doc revised in place across several iterations gets *less* trustworthy
  in exactly the rows nobody has reason to re-read. **PROPOSED UPSTREAM** as a sharpening of the
  shared skill's Gate-2 rule 3b(vi) — it lives in the V1 checkout, so World proposes and never
  applies — and V1 has the same exposure: its own iterations 135 and 138 are this class, and 3b(vi)
  was written there without an instrument.
- **A BLANKET INSTRUCTION IS APPLIED TO ITS EXCEPTION TOO — FAITHFULLY, BY A ROLE THAT CANNOT KNOW
  BETTER (process fix, iter-48).** The iter-48 executor directives said *"run every `go` command
  with `GOTOOLCHAIN=go1.25.6`"*, which is correct for building and testing this repo and **wrong
  for the one probe whose entire purpose is to observe each tree's ambient toolchain**. The
  executor applied it universally — including inside the recorder — so the tool that exists to
  **record** a toolchain was silently **selecting** one, collapsing AC6's straddle and violating
  the design doc's own *"it selects nothing"* clause. The evaluator caught it; the executor was
  blameless. **When a directive states a universal rule, name its exceptions in the same
  sentence, or state the INTENT ("so the build is reproducible") rather than only the mechanism —
  intent lets a role recognise a case the rule was never meant to cover.** Same shape as iter-47's
  designer directive, where a gate left implicit did not exist for a cross-provider role that
  cannot read this repo's skills: a sub-agent obeys what you wrote, not what you meant.
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
- **TWO PRE-REGISTERED WATCH-ITEMS, EACH AT INSTANCE 1 (iter-58).** Gate 5 allows at most ONE skill
  edit per iteration and requires **≥2 recorded frictions pointing at the same gap**; neither of
  these reaches the bar yet, so they are recorded here with their remedy already written, exactly as
  V1's iter-140 pre-registered the rule iter-142 landed. Counting them is the point — an unrecorded
  instance 1 means instance 2 arrives looking like instance 1.
  - **(W1) A SPRINT PLAN'S CLAIMS ABOUT THE CODE GO STALE THE MOMENT ONE OF ITS OWN MILESTONES
    LANDS.** Skill rule 3b(vii) diffs a design doc against its plan **at pick time**, and it is aimed
    at *documents*. Nothing points the same suspicion at the **code**: a plan is written once,
    against the tree as it stood, and then each milestone it describes *changes that tree* — so
    milestone N's step can be factually wrong about a file milestone N−1 rewrote, and the rot is
    invisible because both artifacts were authored by the same planner in the same hour. Instance 1:
    the plan's `BG.B` step says "route BG.A's **two** per-arm writes through `confinedWrite`"; BG.A
    landed a **third** (the AC4 barrier marker), which BG.B's own AST guard would have redded on
    sight. **Remedy at instance 2**: before routing milestone N>1, re-derive by command every
    *quantity* the plan asserts about code that milestone N−1 touched — counts, line numbers, call
    sites — with rule 3a's known-positive control, and treat "the previous milestone edited this
    file" as positive evidence of divergence rather than as reassurance.
  - **(W2) A PLAN NAMES ITS EXECUTOR, AND ITS SAFETY CONSTRAINTS ARE CONDITIONED ON THAT NAME.**
    Instance 1: the item-10 plan's `executor` field is `codex:gpt-5.6-sol under --sandbox
    workspace-write`, while the driver env pins `MISSION_EXECUTOR_MODEL=opus`. The env wins — but
    the plan's `S-7` ("the executor lane cannot commit", hence `.snap/M<k>/` snapshots and a
    controller reconstruction) and its `UNINFORMATIVE UNDER SANDBOX` labelling rule are *derived
    from* the codex assumption and become wrong, not merely redundant, under any other lane. A
    directive that silently carries them tells the executor to do unnecessary work under one reading
    and to distrust valid gate results under another. **Remedy at instance 2**: at Gate 3, diff the
    plan's `executor` field against the role's env value and, where they differ, state in the
    directive which plan constraints are thereby VOID — never leave the executor to infer it.

---

## STATUS (rotation rule)

Newest **3** STATUS stamps live here; older ones move to `world-mission-status-archive.md`.
At Gate 4, after adding your stamp, move the now-4th stamp to the TOP of the archive file.

## STATUS 2026-08-06 (iteration 59) — **DEV IS RED AT HEAD AND IT IS NOT ABOUT THE CODE: BOTH JOBS `cancelled` WITH `steps=0` INSIDE A DECLARED GITHUB ACTIONS `major_outage`. NO CODE LANDED; THE DELIVERABLE IS THE DIAGNOSIS, PER GATE 1'S OUTAGE RULE.** **THE SPINE: A GREEN OBTAINED DURING AN OPEN INCIDENT IS A SAMPLE, NOT A SETTLEMENT — AND THE ITERATION THAT BOUGHT THAT WAS THE ONE IMMEDIATELY BEFORE THIS ONE.** Iteration 58 recorded the CI caveat as **SETTLED** on the strength of `e3808c0`, a descendant of the `BG.A` merge whose `ailang-code verify gate` passed all 11 steps. That inference about the **CODE** was correct and still stands — same code, same job, full step log. What did not follow is the inference about the **INCIDENT**, and iteration 58's own doc-only bookkeeping commit `4e959bf` went red on **BOTH** jobs **13 minutes later**. Measured, SHA-addressed, per job, per step: `10120d6` 10:49Z verify **success/11** · go **success/13**; `278f102` 15:38Z verify **cancelled/0** · go success/13; `ea4b03d` 16:21Z verify **cancelled/0** · go success/13; `e3808c0` 16:33Z verify **success/11** · go **success/13** ← the "settling" green; `4e959bf` 16:46Z **verify cancelled/0 · go cancelled/0**. The pattern is intermittent and got **worse** after the green, not better. **DURING AN OPEN INCIDENT, OUTCOME IS NOT A FUNCTION OF THE TREE**, so neither a red nor a green is attributable to the diff — the loop already knew not to trust the red, and trusted the green. **ATTRIBUTION IS BY MECHANISM PLUS SIX FIRING CONTROLS, NOT BY TIMING (rule 3d).** (1) **No step failed anywhere** — across 5 commits × 2 jobs every job reports `failed=none`; the four red jobs are conclusion **`cancelled`** with **`steps=0`**, i.e. nothing was attributable to any STEP, which is the question `steps=0` is only a proxy for (V1 iter-154's correction, applied here). (2) **Parent arm**: `e3808c0` BOTH jobs success 13 min earlier. (3) **Provider status API, first-party**: `Actions = major_outage`, `Pages = major_outage`, incident opened `2026-08-06T15:22:49Z` and still `investigating` at `20:34:17Z`; HEAD's run was created `16:46:43Z`, **inside the window**. (4) **Sibling mission**: `mission-v1` iteration 154 reports the identical signature in a **different repo** in the same window (`#608` zero runs created; 0 of 4 re-runs started in 28 min). (5) **The diff**: `e3808c0..4e959bf` is exactly **3 markdown mission documents** — `0` `.ail`, `0` `.go`, `0` `scripts/`|`.github/`, `0` `design_docs/verification/` (the one place CI *does* execute a script under `design_docs/`), with a **firing KP control** (the same filter catches `host/boundary/allowlist_world_test.go` on the `BG.A` commit), so the zeros are measurements. (6) **Local gate on the identical tree**: a sibling-of-repo worktree (never `/tmp`) at `4e959bf` runs `verify_go.sh` **rc=0** with pinned **AILANG v0.30.0** — build · plain · `-race`, **24** packages `ok`, **0** FAIL, KP control firing. **DISPOSITION PER THE OUTAGE RULE: NOT REVERTED, NOT FIXED-FORWARD, NOT PARKED.** `BG.A` stays landed; re-runs were fired (`rerun-failed-jobs` accepted, run `31121008498` moved `completed/failure` → `queued`) and are still `queued` at the poll's bound, which is itself the outage. **A re-run is owed once the incident closes, and Gate 3b still binds: 0 observed failures is NOT a green.** **`BG.B`'S PREMISE RE-VERIFIED FIRST-PARTY AT HEAD so the next iteration executes without re-litigating it** (rule 3b(v): a count transcribed from a document is a claim about the document): **3** `os.WriteFile` **calls** at `:383` (AC4 marker), `:428` (mutant), `:439` (overlay JSON) — the plan still says **two** — with `os.OpenFile`/`os.Create`/`os.Rename` all **0**, `confinedWrite` **0** (KP control `checkGoGroup` = **6**, so that zero is a measurement), and the file **byte-identical** at `278f102` and HEAD (`sha256 d535c1ec92641c02…`, which is also the restore hash iteration 58 quoted). **AND A TRAP FOR `BG.B` SPECIFICALLY, FOUND BY MY OWN CONTROL DISAGREEING WITH THE CHARTER:** `os.ReadFile` is **5 textual occurrences but 4 real calls** — `:264` is a **comment** — so the charter's KP of 4 was right and my grep was the imprecise instrument; recorded because `BG.B` installs an **AST** guard, and text-vs-AST disagree by exactly one in the very file being guarded (`os.WriteFile` happens to have **no** comment mentions, so text == AST == 3 there — today). **GATE-4 TELL DEFECT, CAUGHT BY THE SKILL'S OWN KNOWN-POSITIVE CONTROL:** the shared skill prescribes `grep -c "ITERATION <N-1>"` **uppercase** (V1's stamp casing); World stamps `(iteration N)` **lowercase**, so the tell returned **0** — *and so did its control*, which is precisely the signal that the instrument is broken rather than the charter stale. Case-insensitively prev=**1**, control=**2**, rotation invariant `^## STATUS 2026` = **3**, and charter/log `git diff` vs `origin/dev` empty: **the charter was healthy the whole time.** **TWO SHARED-SKILL FIXES PROPOSED TO V1/MARK (World cannot edit the shared skill, per its own frozen-core rule):** (a) the outage rule must say a **green** obtained during an open incident does not close it — instances: this iteration, and V1 iter-153's `docs` job going cancelled→success on a byte-identical tree, which the skill already records as evidence that outcome is environment-driven and then never applies in the green direction; (b) Gate 4's stale-charter tell should be `grep -ci` — instances: V1 iter-134's `Iteration 133` → 0 on a healthy charter (the skill fixed that symptom by ADDING the control rather than fixing the literal) and this iteration, which is the same defect recurring in a second mission. **`metered=$0.00`** — controller-only, no sub-agent spawned, no quorum round, all quota buckets. **ZERO OPEN ASKS FOR MARK.**

## STATUS 2026-08-06 (iteration 58) — **`BG.A` LANDED (PR #47 → squash `278f102`) — ITEM 10's FIRST CODE: THE BOUNDARY GATE NOW PROVES ITS TEETH WITHOUT WRITING THE TREE IT GUARDS. THE SPINE: A CHECKER THAT CANNOT READ THE TREE FINDS NO FORBIDDEN IMPORTS — AND THIS ITERATION MET THAT SHAPE THREE TIMES, ONCE INSIDE THE FIX, ONCE INSIDE THE DOC AND PLAN THAT SPECIFIED IT, AND ONCE INSIDE MY OWN KNOWN-POSITIVE CONTROL.** The mutant is now DECLARED, never written: `go list -deps -overlay=<json>` for the dependency-closure half, an overlay-aware read helper for the import-scan half; `mutateAndRestore` and its `defer`-based restore are **deleted** rather than made safer, because there is nothing to restore when nothing is written. One file, +325/−44, no production code, no script or workflow edits. **AC2, AC3, AC4, AC5 DISCHARGED; `AC1a`/`AC1b`/`AC6′` and `M3`/`M6`/`M7` remain with `BG.B`/`BG.C`.** **CONTROLLER-MEASURED, NOT INHERITED (Gate-2 rule (b): a judge's and an executor's findings are claims too).** AC2's four numbers per Go arm, asserted on the closure `checkGoGroup` **RETURNED** rather than on a second `go list`: `host/store` **160/0 → 229/1**, `host/replay` 162/0 → 231/1, `cmd/ailang-worldd` 233/0 → 234/1 — the planner's `PV-3`/`PV-8` prediction exactly. `M1` and `M2(b)` re-run first-party under the house recipe (anchor count **1** asserted, differing sha256 asserted *before believing*, control arm FIRST at rc=0, byte-identical restore `d535c1ec…` verified between arms): `M2(b)` — the overlay `Replace` KEY naming no real file, i.e. the **silent** failure the planner measured as rc=0/base-closure/no-stderr — reds with `overlay closure=160, baseline closure=160 -- the toolchain half of the gate is dead`, on all three Go arms. Without AC2 that failure is invisible, because the import scan keeps producing a perfectly convincing RED through the read half. **`M5`, THE AC4 KILL HARNESS, IS THE CONTROLLER'S OWN AND IT WAS RUN WITH ITS NEGATIVE CONTROL (rule 3d).** Per arm, deterministically: await that arm's ready marker under a fixed timeout → verify the mutant and overlay JSON EXIST → verify the overlay MAPS the real target to the temp mutant → verify the process is **alive** → `SIGKILL`. All four arms (`host/store`, `host/replay`, `cmd/ailang-worldd`, `world`): `alive_at_kill=True`, `killed_rc=-9`, **0** changed target sha256s, **0** `git status --porcelain` lines. **The same kill on the BASE harness, in a disposable sibling worktree, returned `killed_while_mutating=host/store/store.go`, `RESIDUE=YES`, ` M host/store/store.go`.** Outcomes DIFFER, so the green measures the mechanism and not the environment — which is the whole difference between this result and the one iter-55 warned could mean *"the threat was never exercised"*. AC4's fail-closed property proven in BOTH directions: armed-but-never-killed → `panic: test timed out`, rc=1 (a timeout FAILS, it does not pass); an in-repo marker path → rejected with `resolves inside repoRoot`, and the marker file was **never created**. **RULE 3e RAN FIRST AND IS WHAT MAKES ANY OF THE GREENS MEAN ANYTHING:** all three acceptance commands were baselined on a PRISTINE tree at `10120d6` in a sibling worktree before the executor was spawned — `go test ./host/boundary/ -count=1 -v` rc=0 (12 evidence lines, matching `V16c`'s KP), `verify_go.sh` rc=0 (build · plain · `-race`), `verify_ail.sh` rc=0 — so a green afterwards is attributable to the change rather than to the repo. **A DEFECT IN THE DESIGN DOC AND THE PLAN THEMSELVES, FOUND ONLY BY EXECUTING THEM.** `go/parser`'s `readSource` tests `src != nil` on the **interface**, so a typed nil `[]byte` is a NON-nil interface and is handed back as an **EMPTY SOURCE**; every unreplaced file then parses as `expected 'package', found 'EOF'`. Both artifacts write the helper as `parser.ParseFile(fset, path, <bytes-or-nil>, …)` — precisely the shape that produces it. It surfaced only because `checkGoGroup` surfaces parse errors; had the walk swallowed them, the gate would have been a checker that reads nothing and therefore finds nothing, which is this item's own spine arriving inside its own repair. Isolated in `parseSrc`, documented in a comment, recorded in the doc so the written wording does not reproduce it. **THE EVALUATOR EARNED ITS FEE ON A FINDING IT FILED AS LOW-SEVERITY, AND THE CONTROLLER REPRODUCED IT RATHER THAN ACCEPTING THE LABEL (Gate-2 rule (b)): the plan's `BG.B` write-site count is now WRONG BY ONE, and the missing one would RED THE GUARD `BG.B` INSTALLS.** The plan says "route BG.A's **two** per-arm writes (mutant file, overlay JSON) through `confinedWrite`" — written before the AC4 barrier existed as code. The barrier adds a **third** direct `os.WriteFile(absMarker, …)`. Measured at `278f102`: **3** `os.WriteFile` sites (`:383` marker, `:428` mutant, `:439` overlay JSON), **0** `OpenFile`/`Create`/`Rename`, with a firing known-positive control (`os.ReadFile` = **4**, so the zeros are measurements). Decision 8's AST guard reds on any of the four names outside the single permitted site, so leaving the marker write direct makes `BG.B` red on `BG.A`'s **own landed code**. Routing it through `confinedWrite` is CORRECT rather than an exemption — the marker is *required* to resolve outside `repoRoot`, exactly what `confinedWrite` permits — so the confined writer also becomes the enforcement point for the AC4 marker-path rule and replaces the bespoke `insideRepo` check at `:367–:373`. Plan corrected in place with a `controller_corrections` entry; charter and doc carry it as a `BG.B` precondition. This is rule 3b(vii) at its sharpest: the plan and the code rotted apart **inside the single iteration that produced both**. **A DELIBERATE DEVIATION FROM THE PLAN'S LITERAL SIGNATURE, AND THE CONTROLLER RECORDS THE OPPOSITE VERDICT TO THE JUDGE'S.** The plan specifies `checkGoGroup(root, group, overlay string)`; the executor used a two-field `overlay{jsonPath, replace}` so the toolchain half and the read half are **separately disarmable**. The reason is AC2's own falsifiability: with one string, dropping `-overlay` also disarms the import scan, so `M2(a)` would red at **AC3** instead of AC2 and the toolchain half would go untested — *a mutation shaped to the check tests the check, not the threat* (iter-54's spine). The observed `M2(a)`/`M2(b)` messages confirm it. The evaluator scored it **−5 on design fidelity** as an undocumented departure; on the merits the deviation is what makes AC2 non-vacuous, it is documented on the type, and the PLAN is what was wrong. **MY OWN INSTRUMENT FAILED ONCE AND THE FAILURE IS THE MOST USEFUL THING IN THIS STAMP.** Proving the AC4 barrier is a no-op when unset, my known-positive control returned `BARRIER lines: armed=0, unset=0` — the "control" and the claim agreeing at zero, which reads as a clean result and is in fact an instrument that **cannot see a positive**. Cause: the armed arm ran WITHOUT `-v`, and `V16c` — measured by this very sprint's planner one iteration ago — says a Go test without `-v` emits **nothing**, not even `t.Logf`. Re-run with `-v` on both arms: armed **1**, unset **0**, rc=0. The sprint's own headline measurement invalidated the controller's control, in the same file, one iteration after being recorded. **GATE 3b — RECORDED HONESTLY RATHER THAN ROUNDED TO GREEN.** At the merge SHA `go host build + test gate` is **SUCCESS**; `ailang-code verify gate` is **FAILURE**, and the failure is in **`Set up job`, step 1, before checkout, with ZERO repo commands executed**: `Failed to resolve action download info. Error: Service Unavailable`. **Attribution is by MECHANISM plus two firing controls, never by redness or adjacency (rule 3d):** githubstatus.com reports `Actions: partial_outage` with an incident opened `2026-08-06T15:22:49Z` and the run started `15:38:01Z`, **16 minutes inside it**; the **identical tree** passed the **identical job** on PR #47 at `15:33:02Z` (`51e18968` → SUCCESS, five minutes before the outage bit); the sibling job on the same merge commit is SUCCESS; and both jobs use the same two actions (`actions/checkout@v4`, `actions/setup-go@v5`), so it is not action-specific. **Three bounded re-run attempts, none of which reached a repo command**: two died in the same pre-checkout `Set up job` step, and the third sat `queued` for ~15 min and was **`cancelled`** — the run as a whole ends `completed/failure` with the verify job never having executed a line of this repo's code, while `Actions: partial_outage` was still live at the time. **RESOLVED WITHIN THIS ITERATION, AND THE RESOLUTION IS ITSELF THE CONTROL:** the bookkeeping commit `e3808c0` — a DESCENDANT of `278f102`, so a tree that CONTAINS `BG.A` (`mutateViaOverlay` present **1**, `mutateAndRestore` **0**; control at pre-BG.A `10120d6` is the exact reverse, 0 and 1) — went **green on BOTH jobs, SHA-addressed at `e3808c0`, verified by STEP LOG rather than by badge**: all 11 steps of `ailang-code verify gate` `success`, including `Verify all .ail modules (ai-check = check + Z3 verify)`. So the same job, on the same code, passes; the `278f102` red is conclusively the outage and not the change. **dev is GREEN at HEAD** and iteration 59 inherits no CI carry. **The merge was NOT reverted** — reverting a change proven green on the same tree, by the same job, would have been the worse error, and holding that line is what let the next commit's green settle it 20 minutes later rather than costing a revert plus a re-land. Roles: executor **`opus`** (Agent-tool pinned from `MISSION_EXECUTOR_MODEL`, **not** the plan's assumed `codex:gpt-5.6-sol` — so the plan's `S-7` "the executor lane cannot commit" and its `.snap/M<k>/` reconstruction rule did **not** apply, and BG.A is one ordinary commit; the sandbox `UNINFORMATIVE` caveats likewise did not arise, and this routing delta was stated explicitly in the directive rather than left for the executor to discover); evaluator **`sonnet`** (distinct model ⇒ generator≠judge) **PASS 89/100, round 1, ZERO blocking**, three low-severity non-blocking findings, all carried, none dismissed; designer and planner **NOT fired** — doc and plan both landed at iters 56/57; rotation pointer unchanged at `claude:claude-fable-5`; controller `claude-opus-5`. **`metered=$0.00`** against the `$5` ceiling — every role on a quota bucket, no quorum round, no cross-provider call. **SAFETY**: no publish occurred and no `ailang publish` was invoked in any form, in any arm including probes; no secret was printed; `GOTOOLCHAIN=go1.25.6` on every `go` invocation (the PATH go is **go1.26.4**, which `verify_go.sh:76` explicitly DENIES) and `AILANG_BIN=/tmp/ailang-v0300/ailang` (v0.30.0, `e37b370`) on every gate; **the boundary gate was never run in the main checkout** — it mutates live production sources at base, which is the defect under repair — every run and every re-arming mutation lived in a worktree SIBLING to the repo, never under `/tmp`, and all three probe worktrees plus the poisoned negative-control tree were restored and removed. **GATE-2 CROSS-CHECKS, ALL RUN, ALL RECORDED.** Rule 3b(vi-b) freshness sweep from the doc's OLDEST declared base (`deeb804`, not the newer controller base): **0** non-doc files changed, control firing at **5** (all under `design_docs/`). Rule 3b(v) re-derivation at pick time: the gate file was **351** lines with `repoRoot:60 goListDeps:72 enumerateAIL:96 checkGoGroup:130 checkAILGroup:165 mutateAndRestore:190`, matching `V2` exactly. Rule 3b(vii) doc↔plan diff: doc and plan agree that `BG.A` owns AC2/AC3/AC4/AC5 — and the rot appeared **after** execution instead, which is the new instance above. Already-landed check on FRESH origin: `BG.A` returned **1** hit, the iter-57 planning record, against a firing control (`SM.B1` → `feat(8): SM.B1 … (#43)`). Mid-flight-iteration check (the iteration-149 class): **no open PRs from this loop and no stale sprint worktrees** before starting. **ZERO OPEN ASKS FOR MARK** (`8/OD-2`, `10/OD-1`, `10/OD-2` open, all non-blocking with controller defaults recorded). **WEEKLY EXTERNAL-ISSUE SWEEP: CLEAN** — one open issue (`#32`, the bookkeeping thread), zero-mention count **0**. No rotation due: `#32` created `2026-08-03T06:15:41Z` = **08:15 CEST**, i.e. AFTER the Monday-07:00 **local** boundary, and holds 19 comments (<80). **CROSS-MISSION**: inbox empty, nothing asked of World. Charter, doc, log and dashboard written **in place**; base re-confirmed by a pure fast-forward — **0** local ahead-commits, tree **0**-dirty, one incoming file — `dev == origin/dev` @ `278f102`. Stale-base tell in this charter's own lowercase spelling (`iteration 57` → **1**) **with its known-present control** (`iteration 56` → **2**); structural invariant `^## STATUS 2026` → **3** pre-rotation, **3** post; rotation arithmetic asserted before writing (`after == before + 2 − 2×1`). 4 stamps → newest 3 kept (iter-58, iter-57, iter-56); the positionally-4th (iter-55) archived this Gate 4.

## STATUS 2026-08-06 (iteration 57) — **ITEM 10 `w-boundary-gate-tree-mutation` SPRINT-PLANNED (`BG.A` → `BG.B` → `BG.C`; plan `opus`, lane fail-closed `opus missing-script`) — AND THE PLAN'S REAL PRODUCT WAS FOUR VACUITY-CAPABLE ACCEPTANCE CRITERIA IN A DOC THAT HAD ALREADY PASSED TWO QUORUM ROUNDS.** **THE SPINE: A THRESHOLD WHOSE NOISE IS THE SIZE OF ITS SIGNAL CANNOT FAIL INFORMATIVELY — AND THIS ITERATION FOUND THAT SHAPE IN THE ITEM'S OWN ACCEPTANCE CRITERION, IN ITS AUDIT ROW, AND IN THE CHARTER'S OWN CORRECTION.** No code landed by design: the queue head's declared next unit of work WAS the sprint-planner run, and it produced enough that starting an executor against a doc carrying a false verification row would have been the wrong trade. **THE CONTROLLER FOUND `AC6` AT BASELINE, BEFORE THE PLANNER RAN** — rule 3e, "baseline every acceptance command on a pristine tree", which nothing else in the loop does. `AC6` bounds `go test ./host/boundary/ -count=1` at `≤2× the measured 0.435s green baseline (V13)`. That constant was transcribed from a different worktree at a different cache warmth, so it fails rule 3e(i) — a control is only a control if it runs from a tree in the state the baseline was in. Measured on **UNCHANGED** code at `e9c8c85`, `GOTOOLCHAIN=go1.25.6`, in detached worktrees SIBLING to the repo (never `/tmp`, never the main checkout), tree 0-dirty before and after every run: fresh-worktree FIRST run **0.664 s** and **0.621 s** (n=2, two independent fresh worktrees), warm steady state **0.472–0.507 s**, median **~0.480 s** (n=9). **Zero code change already sits at 1.43–1.53× the AC's own constant in the cold state** — and CI checks out fresh, so cold is the operative state — leaving ~**1.31×** of headroom for an added `go/parser` pass, per-arm `os.Stat`+sha, and overlay JSON writes, not 2×. **The false-red risk is the LESSER harm.** The greater one is that a **GREEN** `AC6` could not have failed informatively, because the observed noise band on an unchanged tree (0.435 → 0.664 = 1.53×) consumes ~**76%** of the stated budget. That is this item's own spine — *a green must be unable to mean "the check never ran"* — arriving inside the item's own acceptance criterion, which is precisely where nobody looks. **THE PLANNER THEN BEAT THE CONTROLLER'S OWN FINDING, TWICE, AND THAT IS THE LOOP WORKING.** (a) **The units are ambiguous by 1.32×**: same warm session, n=5, go-*reported* median **0.479 s** vs wall-clock median **0.631 s**; `0.435 s` is a go-reported figure (`ok … 0.435s`) while the AC's wording is about the *command completing*, so read naturally the unchanged code starts at **1.45×**. (b) **`AC6` could not fail for the change either**: what it nominally protects is a **600 s** `-race` budget against a **0.5 s** package — **1200×** of headroom. Vacuous in BOTH directions. `AC6′` replaces it with a paired same-session ratio (ONE worktree, ONE session, equal warmth, a discarded warm-up, then ≥**8** interleaved A/B pairs swapping only this file and asserting its sha256 changed; `median(wall_B)/median(wall_A) ≤ 1.50`) plus a `median(wall_B) ≤ 3.0 s` ceiling, **asserted on wall-clock and said so explicitly** so the unit cannot be re-ambiguated. Its noise floor was **MEASURED, not assumed**: 8 interleaved pairs on **IDENTICAL** code returned **1.0079** against a true ratio of 1.000, pooled spread 1.058 — so the 1.50 bound is ~**8.5×** the instrument's own spread. **TWO PLANNER REFUTATIONS OF THE DESIGN DOC, BOTH REPRODUCED FIRST-PARTY BEFORE BEING RECORDED** (Gate-2 rule (b)/(d): a sub-agent's finding is a claim too, and my first control for the second one was INVALID — a `-v` run filtered to tests that never emit the lines, i.e. an instrument that could not see a positive, caught and re-run before anything was believed). **(1) `V16` IS REFUTED (→ `V16a`/`V16b`).** `cmd/ailang-worldd`'s baseline closure DOES contain a forbidden prefix: `github.com/sunholo-data/ailang-world/host/registry`, which is `forbiddenImportPrefixes[3]` (`:53`), reached transitively via `cmd/ailang-worldd → host/daemon → host/registry` (`daemon.go:51`, a direct import). Measured worldd **1** of 233, store **0** of 160, replay **0** of 162, with KP `modernc.org/sqlite` firing **1 / 1 / 1** in the same call. So the doc's *"0 forbidden-prefix hits in all three closures"* and *"a scan would be green on the current tree"* are **FALSE**. **BUT THE RED WOULD BE A FALSE POSITIVE, SO THE `10/OD-1` DEFERRAL GETS STRONGER, NOT WEAKER** — and this is the correction to the planner's own framing: `host/registry` is the **interpreter epoch** registry (`world/epoch-registry/v1`, `w-world-library-m1` Decision 5), NOT the *package* registry the forbidden entry targets, and `host/daemon` legitimately needs epoch metadata. **iter-53 PREDICTED THIS EXACT NAME COLLISION IN PROSE; iter-57 MEASURED IT.** A closure scan today reds legitimate code, so `10/OD-1` cannot be implemented until the collision is resolved — shipping it as-is would install a gate whose red means nothing. The doc's error hid because ONE table cell bundled TWO questions (`forbidden prefixes?` and `bare net/http?`) and only the second ever had a firing control — a bundled cell is a place a wrong number can hide behind a right one. **(2) NOTHING A PASSING GO TEST EMITS REACHES CI (`V16c`).** Paired arms on `TestWorldBoundaryDependencyAllowlist`, same worktree: **A** = CI's exact form (no `-v`) → rc=0, output is the single line `ok … 0.580 s`, matching lines **0**; **B** (the KP) = identical **+ `-v`** → rc=0, matching lines **12**, the ENUMERATION rows dumping the full 160/162/233 closures. `verify_go.sh:100` is `go test ./... -count=1` and its `-race` leg builds `["go","test","./...","-count=1","-race","-timeout","8m"]` — **neither carries `-v`**. So the gate's `ENUMERATION`/`MUTATION`/`RESTORE` diagnostics have **NEVER appeared in a CI log** since they were written, and *"loud but non-gating"* is a contradiction in CI. Consequence, carried into the plan: any observable this design wants CI to see must be an **ASSERTION**, never a log line. **A THIRD VACUITY MODE THE DOC'S `Decision 8` NEVER CONSIDERED:** `host/boundary` holds **exactly ONE `.go` file** and it is a `_test.go`, so an empty `ParseDir`, a filter dropping `_test.go`, or a selector bug each yield the AST guard "zero violations, green" — the doc defends the *self-match* mode at length and is silent on the *empty-enumeration* mode, which is rule 3a (a search that found nothing is a claim) wearing the guard's clothes. Repaired in the plan by an exact non-empty file count, a **known-positive requiring the walker to FIND the permitted `os.WriteFile` and report its line**, and deny-list completeness. Also measured by the planner and folded in: `go list -overlay` **SILENTLY IGNORES** a `Replace` key matching no file (rc=0, base closure, **no stderr**), so `AC2`'s "closure contains `mutantImport`" needed strengthening to the closure `checkGoGroup` actually **CONSUMED** plus a negative half (free — `checkGoGroup` already computes and discards it); and `git status --porcelain` reports **untracked** files, so `AC4`'s in-tree ready marker would have redded on the harness's own artifact (marker moved outside `repoRoot`, residue assertion path-scoped to the four targets). `AC3` and `AC5` are sound; `AC3`'s `:217–:222` citation was **never verified by the doc** and the controller verified it first-party (it is exactly the two RED-fidelity assertions, `:218` and `:221`); `AC5`'s baseline PASSES at HEAD (both named tests, rc=0). **MILESTONES: `BG.A` (AC2, AC3, AC4, AC5 · M1, M2, M4, M5) → `BG.B` (AC1a · M3, M6) → `BG.C` (AC1b, AC6′ · M7)** — partition COMPLETE, 7 criteria and 7 mutations, none dropped or double-assigned. ONE ordering joint is **forced** (the AST guard reds on sight while `mutateAndRestore` still calls `os.WriteFile` at `:205/:209/:224`, so `BG.B` cannot precede `BG.A`); one is **chosen** (`BG.C` last, because it is the only milestone whose green depends on an UNMEASURED CI-filesystem property, so a red there leaves the repair already landed). **CARRIED, NOT CLOSED — `C1` IS APFS-ONLY AND MAY NOT TRANSFER:** the `ModTime` backstop was measured 200/200 on darwin/APFS while CI job 2 is `ubuntu-latest`; **Linux takes mtime from a tick-granularity coarse clock (1–4 ms)**, and the planner could not measure it (no docker/colima/podman/lima on the rig) — so it is labelled UNMEASURED rather than assumed. Hence `BG.C` last, a **fail-loud 20/20 granularity probe whose failure is a TEST FAILURE naming both `st_dev`s**, sha256+size+mode+**inode** asserted unconditionally (inode closes a rename route that BOTH of the doc's stated observables miss), and a pre-authorized fallback: record the refutation, keep the four filesystem-independent observables, open `10/OD-3` — and **NEVER lower the 20/20 threshold**. **ESTIMATE HONESTLY SPLIT:** the doc's **≤1 day of EFFORT holds** (velocity derived by command — 12 landed feat/fix commits, insertions median **363**, mean 447; closest analogue `1761a9c`, a single-test-file change to *this same file*, +56/−7, one iteration) **but ELAPSED is 2–3 iterations**, because measured cadence is ≤1 milestone/iteration across iters 41–56 and **4 of the 7 mutations CANNOT run in the executor sandbox** (M5 needs subprocess SIGKILL + git inspection; M6/M7 re-arm live-tree writes; AC6′ needs a file swap out of git history) — **the controller pass is the critical path**, which is a routing fact, not a complaint. LOC re-estimated **+150 → +250**, itemised, all of it spent making the doc's own ACs non-vacuous. **A CHARTER DEFECT FOUND IN ITER-56'S OWN CORRECTION, AND IT IS THE SAME SHAPE AS EVERYTHING ELSE HERE.** iter-56 measured *"deliberately non-compiling"* FALSE and stamped *"Corrected in all three places this Gate 4."* Verified rather than inherited: the **dashboard** was genuinely corrected (1 hit, and it is the correction itself), but the charter's **item-10 queue row still carried the false sentence in LIVE prose ~35 lines BELOW its own correction** — `~~` marker count **2** in that span proves one closed strikethrough covering only the `[BACKLOG …]` tag, so the sentence was never struck. **A correction that does not delete the text it corrects leaves a reader able to find the false version, and the loudness of the correction is what makes nobody re-read the row.** Fixed this Gate 4. The iter-55 STATUS stamp also still carries it and is left alone deliberately — stamps are append-only history, corrected forward, not rewritten. **GATE-2 CROSS-CHECKS, ALL RUN, ALL RECORDED.** Rule 3b(vi-b) freshness sweep from the doc's OLDEST declared base (`deeb804`): **0** non-doc files changed, with a firing control (the unfiltered diff returns **5**, all under `design_docs/`) — so every V-row was measured against code identical to HEAD, which is exactly why V16 being wrong is a *designer* error and not a staleness artifact. Rule 3b(v) re-derivation: the file is **351** lines with `repoRoot:60 goListDeps:72 enumerateAIL:96 forbiddenImport:121 checkGoGroup:130 checkAILGroup:165 digest:188 mutateAndRestore:190`, matching `V2` exactly, and `:276` carries `// boundary mutation: compiling HTTP import`, confirming `V22`. Rule 3b(vii): there was NO prior plan, so the obligation was forward-looking — the doc now carries a ⚠ SPRINT-PLANNED block recording every supersession, so plan and doc cannot rot apart. Quorum-at-pick satisfied (2 artifacts present, both rounds from iter-56); already-landed check on FRESH origin returned only the design-doc and record commits against a firing control (`SM.B1` → `feat(8): SM.B1 … (#43)`); **no open PRs from this loop and no stale sprint worktrees** (the iteration-149 mid-flight class). **INSTRUMENT DISCIPLINE, INCLUDING MY OWN FAILURES:** my Gate-4 stale-base tell first ran as `grep -c "ITERATION 56"` — **0**, on a perfectly healthy charter — and its **control also returned 0**, which is what identified a broken instrument rather than a stale base; this charter spells stamps `(iteration N)` **lowercase**, and the corrected tell read **1** with its control at **2**. My first known-positive for the `-v` refutation was invalid (filtered to tests that cannot emit the lines) and was re-run. And a `git -C <repo> worktree add <relative-path>` resolved **INSIDE the repo**, dirtying the shared checkout — removed, tree re-verified 0-dirty; the skill forbids `/tmp` for worktrees and the mirror-image error is just as easy to make. Roles: planner **`opus`** (Agent-tool pinned, lane `opus missing-script` fail-closed and recorded VERBATIM); designer, executor and evaluator **NOT FIRED** — a planning iteration; rotation pointer unchanged at `claude:claude-fable-5`; controller `claude-opus-5`. **`metered=$0.00`** against the `$5` ceiling — every role on a quota bucket, no quorum round, no cross-provider call. **SAFETY**: no publish occurred and no `ailang publish` was invoked in any form, in any arm including probes; no secret was printed; every `go` command ran in a detached worktree SIBLING to the repo under `GOTOOLCHAIN=go1.25.6` (the PATH go is **go1.26.4**, which `verify_go.sh:76` explicitly DENIES), the boundary gate was **never** run in the main checkout — it mutates live production sources, which is the defect under repair — and all probe worktrees were removed. **ZERO OPEN ASKS FOR MARK** (`8/OD-2`, `10/OD-1`, `10/OD-2` open, all non-blocking with controller defaults recorded; `10/OD-1` now carries a measured reason it cannot proceed). **WEEKLY EXTERNAL-ISSUE SWEEP: CLEAN** — one open issue (`#32`, the bookkeeping thread), zero-mention count **0**. No rotation due: `#32` created `2026-08-03T06:15:41Z` = **08:15 CEST**, i.e. AFTER the Monday-07:00 **local** boundary, and holds 18 comments (<80). **CROSS-MISSION**: 9 unread, all informational — `eval-suite` run notifications plus mission-v1's iter-150 report; nothing asked of World. Charter, doc, log and dashboard written **in place**; base re-confirmed `dev == origin/dev` @ `e9c8c85` with `git diff --stat origin/dev` EMPTY on all four before the first write. Stale-base tell in this charter's own lowercase spelling (`iteration 56` → **1**) **with its known-present control** (`iteration 55` → **2**); structural invariant `^## STATUS 2026` → **3** pre-rotation, **3** post; rotation arithmetic asserted before writing (`after == before + 2 − 2×1`). 4 stamps → newest 3 kept (iter-57, iter-56, iter-55); the positionally-4th (iter-54) archived this Gate 4.

## CURRENT GOAL

1. ~~Iteration 0: ratify the bar~~ — **DONE 2026-07-23 attended** (Mark: clause-4 fixed at
   −2pp / ≤25%; bar + Conflict Surface + guardrails + queue ratified as drafted; re-quorum run
   attended — see STATUS).
2. **NOW**: work the queue through the inner loop (design-doc → sprint-plan → execute →
   evaluate), one sprint-sized item per iteration, recording routing evidence every time.
   **D1 RATIFIED 2026-07-24 (Mark, attended: A+B-metadata) — queue UNBLOCKED.**
   ~~Next item: `[NEXT] w-world-library-m1` (clause-1) — M1 may now freeze the log format per the
   settled decision doc.~~ **Stale since iter-16, corrected iter-40** — that item completed
   2026-07-27, so this pointer went 24 iterations without an update, which is the same
   prose-outlives-its-measurement decay the iter-27 guardrail names.
   ~~**Next items (iter-40): the two remaining small gate-integrity items — `4d w-ddl-gate-teeth`
   and `4f w-bench-load-confound`.** `4e w-race-gate-blindspot` has its mechanism identified and
   its remediation (**RG.A** / carry-forward **CF-K-1**) **PARKED on OD-1/OD-2, awaiting Mark**.~~
   **Corrected iter-45**: `4d w-ddl-gate-teeth` is **COMPLETE** (DG.A `ad619d8` + DG.B `e6ece55`;
   doc → `implemented/`), and 4e's OD-1/OD-2 were ratified in the attended TRIPLE RATIFICATION 2.
   **Next items: `4e w-race-gate-blindspot` (RG.A / CF-K-1) and `4f w-bench-load-confound`
   (branch A) — both routable, neither blocked on Mark.**
   When this line and the Queue disagree, the **Queue** governs — and whoever notices fixes this line.

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

- **AN ACTION THAT SILENTLY DID NOT HAPPEN RETURNS THE SAME RESULT AS A GENUINE NEGATIVE — ASSERT
  THE ACTION'S OWN *EFFECT*, NEVER MERELY ITS EXIT CODE** (process fix, iter-45; 2 instances in ONE
  iteration, both in the controller's own verification). The mission's instrument discipline covers
  *searches* that come back empty and *checks* that come back green — both are rules about reading a
  **result**. Neither covers the step before it: the action never occurring at all. That failure is
  strictly worse, because the result it produces is not merely uninformative but **actively
  misleading in the confident direction**.
  (a) **A MUTATION YOU DID NOT PROVE *APPLIED* IS NOT A MUTATION.** Reproducing
  `MUT-PREFIX-WILDCARD-REGRESSION` first-party, a `perl -0pi -e "s/…/…/g"` substitution **matched
  nothing** (backslash escaping through the shell), left `store.go` **byte-identical**, and the named
  test then returned **rc=0 / `ok`**. Read naively that green says *the executor's mutation does not
  discriminate* — refuting a **real** finding on the strength of an edit that never happened, with an
  authoritative-looking first-party command. A mutation test has two failure modes sharing one exit
  code — the property is genuinely undetected, or the mutation was never there — and **only one is a
  finding**. Caught solely because the same call carried a sha256 before/after assertion.
  (b) **A KILLED, CAPPED OR TIMED-OUT COMMAND'S OUTPUT IS VOID, NOT NEGATIVE.** Two unbounded `find`
  sweeps for `4d/OD-7` were killed mid-traversal; their **empty** output was indistinguishable from
  *"no World stores exist on this rig"* — and would have discharged a **human ratification gate** on
  no evidence whatsoever. Caught because a known-positive control store was **absent from the
  sweep's own output**.
  **Rules.** After any mutation or edit, prove the file **changed** (sha256 before/after, or an
  asserted replacement count) *before* running the test that interprets it — and prove it changed
  **back** afterwards. Prefer `python3` with an asserted `count()` over shell-quoted `perl`/`sed`
  wherever escaping is non-trivial. Record a killed/timed-out command's output as **VOID** in the
  log, never as a negative result. And any sweep whose **emptiness is load-bearing** must carry a
  known-positive target that **must appear in its own output**: if the control is missing from the
  output, the sweep proved nothing, whatever its exit code. This is iteration 44's two-sided-control
  rule aimed one step earlier — at the instrument's *operation* rather than its *reading*.
- **`OD-<n>` IS A MISSION-GLOBAL NAMESPACE, AND THIS TABLE IS ITS REGISTRY** (process fix, iter-43;
  **third instance** of the ID-collision class after iter-31, and the first where the collision landed
  inside a **human ratification**). A design doc numbers its open decisions locally, but the human
  answers them in **one charter stamp and one issue thread**, so two docs using `OD-3` for different
  questions makes a ratification ambiguous *at the moment it is being applied*. Measured this
  iteration: Mark's single attended stamp reads *"4d RATIFIED — fix NOW: … `user_version` pin"* **and**
  *"OD-3 stays declined-as-primary"*. Both sentences are true, of **different items** — and read
  carelessly the second un-ratifies the first. 4d's doc already carried a numbering note written
  against exactly this (it cites iter-31, and it is why OD-5 exists rather than a third `OD-1`), but it
  checked only the **charter's parked list** and not the sibling **doc**, so it deconflicted OD-1/OD-2
  and silently re-collided on OD-3/OD-4. **A remedy that checks a narrower scope than the failure is
  not a remedy** — same shape as this mission's gate defects, aimed at bookkeeping.
  **Rules.** (a) Before allocating an `OD-<n>`, enumerate `### OD-` headings across **every** doc in
  `design_docs/planned/` — not the charter's parked list — with a known-positive control, and take the
  next free integer **mission-wide**. (b) Add the row here in the same edit; an unregistered OD does
  not exist. (c) When quoting a decision anywhere a human reads it, write **`4d/OD-3`**, never a bare
  `OD-3`. (d) Existing collisions are NOT renumbered — renaming an ID a human has already ruled on is
  how a collision becomes a silent contradiction; they are disambiguated in the table instead.

  | ID | Item | Question | State |
  |---|---|---|---|
  | `4e/OD-1` | 4e `w-race-gate-blindspot` | lower the `go.mod` floor 1.26.4 → 1.25.6 | **DISCHARGED** iter-46. RATIFIED 2026-08-03; the floor is `go 1.25.6` as of `f19acac`, enforced by a read-only loud assertion in `verify_go.sh` and detected version-agnostically by `host/store/toolchain_canary_test.go` |
  | `4e/OD-2` | 4e | file the 52-line reproducer upstream at `golang/go` | **DISCHARGED** iter-46 — filed as [`golang/go#80706`](https://github.com/golang/go/issues/80706) after a duplicate search came back empty **with its instrument control firing**. Includes the new AC6 measurement (amd64 unaffected), which narrows the platform for the Go team |
  | `4e/OD-3` | 4e | *also* change the two `scan.go` sites to a compiler-safe shape | **DECLINED as primary** (optional belt). Honoured: `scan.go` is byte-untouched by RG.A (`git diff --exit-code` rc=0, verified by executor AND evaluator). Revisit only if the pin is ever reverted |
  | `4e/OD-4` | 4e | make `maxRecoveryPages` injectable (CF-MJC-1) | **OPEN, and now the named escalation path for a CI `-race` overrun.** Cost re-measured iter-46 and it is NOT the doc's `176.6 s of 179 s`: six runs at one commit spread `host/broker` over **69-175 s** (2.54x, nominal load), so no single figure justifies the change. CI carries `timeout-minutes: 25` with expiry defined as a RED that routes here — never a silently raised ceiling |
  | `4d/OD-3` | 4d `w-ddl-gate-teeth` | add + enforce a `PRAGMA user_version` contract now | **RATIFIED** 2026-08-03 — alt 1, fail LOUD → milestone **DG.B**, **DESIGNED** iter-44 and **LANDED** iter-45 (`e6ece55`). Closed |
  | `4d/OD-4` | 4d | on DDL change, fail or migrate | **FUTURE, not a blocker** — its own recommendation is alt 1 (fail-loud only) "until a concrete DDL change supplies real migration requirements", cost-to-defer **none**, and DG.B *is* alt 1. Travels with the doc to `implemented/`; re-open when a real DDL change arrives |
  | `4d/OD-5` | 4d | does the no-silent-fallback axiom override frozen-kernel deferral | **ANSWERED** 2026-08-03 — yes; DG.A landed iter-43 |
  | `4f/OD-6` | 4f `w-bench-load-confound` | contemporaneous A/B, or stop claiming cost validity | **RATIFIED** 2026-08-03 — **branch A** |
  | `4d/OD-7` | 4d `w-ddl-gate-teeth` | how are valuable deployed **version-0** stores transitioned once DG.B starts refusing them — reconstruct, certify-and-stamp, or delay rollout? | **DISCHARGED** iter-45. RATIFIED 2026-08-03 (attended) as a verified-disposable sweep; the sweep **ran and came back CLEAN** — **ZERO** World stores on the rig (13 candidate SQLite files read-only through the pinned `modernc` driver, all **0/7** canonical tables; the one World-shaped hit is an orphaned `world.db.artifacts` whose database no longer exists). Coordinator attestation VERIFIED, no unexpected store, no re-park. Rollout unblocked and DG.B landed `e6ece55` |

  | `4f/OD-8` | 4f `w-bench-load-confound` | the ratification promises *"mechanically valid cost claims"*; branch A, built exactly as ratified, delivers mechanically complete, contemporaneous, tamper-evident **evidence** — R1–R6 enforce nothing about load, by ratified design. Correct the claim, or grow the mechanism? | **OPEN** — raised iter-47 at quorum round 5, `gpt5-6-sol`'s second reject on this ground. Alternatives: **(1)** correct the claim, ship branch A as designed (no new design work, BC.A′/BC.B′ routable immediately); **(2)** add a comparability criterion (the excluded third limb — real new design work, re-quorum, re-sizing; the loop does not believe a defensible threshold is obtainable on this rig, and P29 measures why); **(3)** amend the ratification's wording in the charter, otherwise as (1). **Recommend (1), optionally with (3)** — and the strongest argument is the objecting reviewer's OWN round-2 fallback: *"If no defensible comparability rule can be derived, revise the policy to say the pair is mechanically complete evidence—not a mechanically valid cost claim."* Full packet: `design_docs/planned/w-bench-load-confound.md` → *OD-8* |

  Next free ID: **`OD-9`**.

- **THE SAME COLLISION EXISTS IN THE `CF-<letter>-<n>` CARRY-FORWARD NAMESPACE, AND IT IS WORSE THERE
  BECAUSE NOTHING REGISTERS THEM** (process fix, iter-46; **fourth instance** of the ID-collision
  class after iter-31, iter-43's `OD-3`/`OD-4`, and this one). The iter-43 guardrail above fixed
  `OD-<n>` by giving it a registry. Carry-forwards were left alone — and they are allocated far more
  often, by more roles (judge, executor, planner, controller), with no table anywhere. **A remedy
  that fixes one namespace and not its sibling is the same shape as this mission's gate defects.**
  Found by the RG.A planner and confirmed first-party: **`CF-K-1` names two different things** —
  `world-mission-log.md:4554` (M3.D, `putRecord` error legibility) and `:5514` (milestone **RG.A** of
  item 4e) — so the ledger line at `:5567` reading *"`CF-K-1`, `CF-K-2`, `CF-MJC-1` (inherited from
  iter-40)"* and the line at `:5337` reading *"`CF-K-1/K-2/K-3`"* refer to **different open items**
  under identical text. Separately, **`CF-K-3` vanished from the enumerated list at `:5567` with no
  recorded closure** — an item does not stop existing because a later ledger forgot it.
  **Rules**, deliberately the OD rules with the scope widened rather than a new mechanism:
  (a) before allocating a `CF-<letter>-<n>`, grep **every** log entry and every doc in
  `design_docs/planned/` + `design_docs/implemented/` for that exact ID — with a known-positive
  control in the same call — and take a letter-suffix that is free **mission-wide**, not free within
  the milestone you are writing about; (b) when quoting a carry-forward anywhere a human or a
  downstream role reads it, write **`4e/CF-K-1`**, never a bare `CF-K-1`; (c) **existing collisions
  are NOT renumbered** — renaming an ID that a later entry already cites is how a collision becomes a
  silent contradiction; disambiguate in prose at the point of use, exactly as the OD table does;
  (d) a carry-forward may only leave a ledger with an explicit **closure line naming where it was
  discharged** — dropping it silently is indistinguishable from forgetting it, which is what happened
  to `CF-K-3`. The two live collisions (`CF-K-1`, `CF-K-2`) are disambiguated at their use sites and
  `CF-K-3` is restored to the open list below rather than renumbered.

- **A RULE FOR FUTURE ALLOCATIONS IS NOT A SWEEP OF EXISTING ONES — AND THE UNSWEPT ONES ARE THE
  DANGEROUS HALF, BECAUSE THEY ARE ALREADY CITED** (process fix, iter-47; **instances 5 and 6** of
  the ID-collision class, found by the controller while reading the carry-forward ledger for a
  routing decision). The rule directly above was written at iteration 46 and it governs the
  *allocation* step: grep mission-wide before taking a letter-suffix. It never asked anyone to look
  at the IDs already in the ledger — so it declared exactly two live collisions (`CF-K-1`,
  `CF-K-2`), the two that happened to be in front of the author, and left the rest unexamined.
  **There were more, and they were allocated at iteration 42 — after the `OD-<n>` registry existed
  and before this rule widened it to carry-forwards, i.e. in the gap the remedy itself left.**
  Confirmed first-party by reading both sites: **`CF-M-1`** means *"the gate reads no
  `skipped_tests`"* at `world-mission.md:1140` (item **4b**, `w-effect-broker-m3`) **and** *"D4's CI
  assertion must grep the specific `hw.ncpu` marker"* at `:1860` (item **4f**); **`CF-M-2`** means
  *"`host/replay/replay.go:325` and `host/archive/archive.go:382` repeat the iter-33 process-group
  defect"* at `:1140` **and** *"P25's deadline arithmetic is designer-measured, not controller-re-run
  cold"* at `:1860`. Two IDs, four meanings, all four live, all four cited.
  **This is the shape the mission keeps finding, aimed at its own remedies: a fix that changes the
  process and not the state.** It is the same error as ratifying a decision without designing it
  (iter-43) and as installing a gate without a control (iter-44) — the remedy is real, and the
  population it was meant to protect was never inspected.
  **Rules**, additive to (a)–(d) above: **(e)** when a namespace rule is introduced, it comes with a
  **one-time sweep of the existing population**, and the sweep's *result* is recorded — including,
  explicitly, what it could not establish; **(f)** per rule (c) these are NOT renumbered —
  `4b/CF-M-1` and `4f/CF-M-1` are disambiguated by their item prefix at every use site, which is
  what rule (b) already requires and what nobody was doing; **(g)** the sweep instrument gets the
  same burden of proof as any other (rule 3a): iteration 47's allocation-site detector returned
  **0 for 60 of 67 IDs** because it matched only one of several allocation phrasings, so its
  **zeros are uninformative and it CANNOT establish that no further collisions exist** — the four
  it found are measurements, the silence is not. A later iteration reading this should treat the
  sweep as **partial and say so**, not inherit it as an all-clear.

- **A KNOWN-POSITIVE CONTROL PROVES AN INSTRUMENT *CAN* FIRE — IT NEVER PROVES IT FIRES *ONLY WHERE
  IT SHOULD*. TEST THE EXCLUSION BOUNDARY, OR YOU HAVE MEASURED HALF THE PREDICATE** (process fix,
  iter-44). This mission's instrument discipline is built on the known-positive control: pair every
  empty or negative result with something you KNOW is there, so a broken instrument cannot wear the
  clean result's clothes. That rule is correct and it is **one-sided**. It certifies the *inclusion*
  half of a filter and says nothing about the *exclusion* half — and a filter is a claim about both.
  **Measured, in this loop's own hands, one round after quoting the rule.** Iteration 44's row **V-K**
  re-checked `name NOT LIKE 'sqlite_%'` against the iteration-41 `sqlite_%` scar, found it excluded 7
  real `sqlite_autoindex_*` rows where iteration 41's had excluded zero, confirmed with a
  known-positive control that the pattern matched exactly those 7, and pronounced the limb **LIVE**.
  Every step of that was true. It was also **the wrong half**: in SQL `LIKE`, `_` is a
  **single-character wildcard**, so the predicate excludes `sqlite` + *any* character. `gpt5-6-sol`
  caught it; the controller then measured it (**V-S**) — SQLite **rejects** `sqlite_internal_probe`
  as reserved but **accepts** `sqliteX_probe`, and a version-0 store whose only application object is
  `sqliteX_probe` counts **0** under the predicate, so it classifies as **fresh** and would be
  initialized and **version-stamped over real application data**. A store-corrupting misclassification,
  sitting behind a control that had just certified the instrument.
  **Rules.** (a) For any predicate, pattern, filter or allowlist that DECIDES something, run a
  **two-sided** control in the same call: one value that must be caught and one that must **survive**.
  One-sided is a half-measurement, and the surviving half is where the damage lives. (b) The tell is
  grammatical: whenever you write "excludes only X", "matches just the Y", "internal objects are
  filtered out", you have made a claim about the complement and almost certainly tested the set.
  (c) Wildcards are the sharpest instance because they are invisible — `_` in `LIKE`, `.` in a regex,
  `*` and `?` in a glob, `|` in an alternation — so any literal containing one needs escaping AND a
  negative control proving the escape works. (d) This subsumes the older rule rather than replacing
  it: still pair negatives with known-positives; now also pair positives with known-negatives.
  **Second, cheaper instance the same iteration**: Gate 4's prescribed stale-base tell greps
  `ITERATION <N>` **uppercase**, while this charter writes `(iteration 43)` **lowercase**, so it
  returned **0 for the target AND 0 for its control**. No harm — the control did exactly its job and
  reported a broken instrument rather than a stale charter — but note that the skill's literal is
  **V1's casing**, and this mission's tell is `grep -c "(iteration <N-1>)"` with `(iteration <N-2>)`
  as the control and `^## STATUS 2026` == 3 as the structural invariant.
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
  **Releases are MARK'S sole decision (Mark, 2026-07-29: "I determine releases")**: when World
  needs an upstream fix in consumable form (pinned-release discipline), ask MARK directly —
  attended, or via issue #9 naming him — never the v1 loop. Standing trigger on record:
  ailang#498 merging → ask Mark for the next release.
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
- **A REVISION IS NOT A SMALLER CHANGE THAN AN ORIGINAL — GIVE THE FIX THE SAME ADVERSARIAL READ YOU
  GAVE THE THING IT FIXES** (process fix, iter-40; **3 instances in ONE iteration**, all inside one
  doc's two quorum rounds). Four blocking objections were raised against `w-race-gate-blindspot`, and
  **three of them were against text written to satisfy an earlier objection**. The failure mode is
  specific and repeatable: a revision arrives feeling like *compliance* rather than *design*, so it is
  typed rather than re-derived, and it inherits none of the scrutiny the original got — while being
  written under exactly the pressure that produced the original defect.
  Instance 1: the remedy for a gate that cannot fail **contained a gate that cannot fail** — the
  proposed `verify_go.sh` exported `GOTOOLCHAIN=go1.25.6` *before* asserting the resolved version, so a
  hostile `GOTOOLCHAIN=go1.26.4` was overridden and the acceptance criterion demanding the script fail
  could never fail.
  Instance 2: the softened second attempt (assign-if-unset) was **still a silent fallback** and
  contradicted the sentence printed immediately after it in the same section (*"a verifier that also
  silently sets the thing it verifies is not a verifier"*). The doc argued against itself one line
  apart and the author did not see it.
  Instance 3: **adding a premise to satisfy round 1 made a round-1 paragraph self-contradictory** —
  logging "a `go` directive is only a floor" as a verified premise destroyed the adjacent claim that
  configuration proved what historical builds resolved. Satisfying an objection *created* the next one.
  **The rules**: (a) after writing a revision, re-read the **whole section** it lands in, not the diff
  — instances 2 and 3 are invisible in a diff and obvious in the paragraph; (b) ask of every fix *"does
  this now contradict something the doc already says, or something I just added?"* and `grep` the doc
  for the claim you just changed (the iter-28 propagation rule pointed inward — a correction can
  *create* an inconsistency as easily as leave one); (c) treat a reviewer-prescribed fix as a **claim**,
  per the iter-29 prescribed-fix guardrail — one of these four objections was itself **factually
  wrong** (`gemini-3-1-pro` on the `toolchain` directive), and adopting it on the reviewer's authority
  would have shipped a **decorative pin while the doc claimed the repo was pinned**; (d) the tell that
  you are about to pay for this: you are applying a fix you did not measure, to satisfy someone else's
  objection, and it feels like paperwork.
- **THE COMPILER IS AN INSTRUMENT, AND EVERY GATE IN THIS REPOSITORY REPORTS THROUGH IT** (process fix,
  iter-40). `ai-check`, `go test`, `verify_go.sh`, `verify_ail.sh` and **the mutation protocol itself**
  all reach the controller through a code generator, and go1.26.0–1.26.5 demonstrably emit a **corrupt
  string** for a source shape present at two sites in landed `host/store` code (see
  `design_docs/verification/w-race-gate-blindspot/`). So a named RED mutation that "reds as predicted"
  proves less than it appears to when the toolchain underneath it is unverified. **The rules**: (a) a
  gate result is scoped to the toolchain that produced it — record the resolved `go version` (and the
  pinned `ailang` version) alongside any gate verdict that will be *cited later*, exactly as iter-39
  requires rig load to be recorded alongside a benchmark number; (b) **configuration is not
  resolution** — `go.mod` and a workflow's `go-version-file` state what was *requested*; only a build's
  own `go version` output states what *ran*, and both the `go` and `toolchain` directives are floors,
  never ceilings (only `GOTOOLCHAIN` selects exactly); (c) when a test fails in a way the source cannot
  explain, **bisect the toolchain before bisecting the code** — two toolchains is a cheaper experiment
  than a mechanism hypothesis, and here it turned three iterations of UNKNOWN into a 52-line reproducer.
- **A `go.mod` `go` DIRECTIVE IS A FLOOR, SO THE CORRECTIVE TOOLCHAIN IS UNSELECTABLE WHILE OD-1
  SITS PARKED — EVERY LOCAL `host/store` MEASUREMENT NOW COSTS A `go.mod` EDIT (process fix,
  iter-41).** Iteration 40 established that go1.26.0–1.26.5 miscompile a shape at
  `host/store/scan.go:74`/`:112` and that `GOTOOLCHAIN` is the only thing that selects a toolchain
  exactly. What it did not price is the consequence for *measurement*: with `go 1.26.4` in `go.mod`,
  `GOTOOLCHAIN=go1.25.6 go test ./host/store/` is refused outright — `go: go.mod requires go >=
  1.26.4 (running go 1.26.4; GOTOOLCHAIN=go1.25.6)`. So until **OD-1** is ratified, any first-party
  measurement touching `host/store` must be taken in a **throwaway worktree with the floor
  temporarily lowered**, and that condition must be **stated alongside every number it produces**
  (per the iter-40 rule that a gate result is scoped to the toolchain that produced it). This is a
  recurring per-iteration tax on a parked decision, not a one-off: it belongs in OD-1's cost column
  next to the correctness argument, and it is why iteration 41 carries **CF-L-3**.
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
4c. [**[LANDED 2026-07-29 (iter-39) — ITEM COMPLETE]**, all three milestones dev-CI-green (both
   jobs, SHA-addressed on each merge commit), doc → `design_docs/implemented/w-effect-journal.md`
   with AC1–AC13 all met. **MJ.C LANDED** via PR #28 → squash `460ade3`, judge sonnet **PASS
   85/100 zero blocking**, executor `codex:gpt-5.6-sol`, `metered=$0.00`. MJ.C discharged **CF-N-2**
   (the `maxRecoveryPages` bound is now justified AND witnessed at its real value, 2^20 pages) and
   **CF-N-3** (the `retryAllowed(false,true)` row landed atomically at all three sites per PD2),
   plus all **eight** open MJ.A/MJ.B carry-forwards (CF-MJA-1/2/4/5, CF-MJB-2/3/4/5). +305/−22.
   **THE PLAN'S HEADLINE MUTATION WAS ONE NAME WITH TWO FORMS, AND THIS TIME IT WAS PREDICTED AT
   GATE 2 RATHER THAN DISCOVERED AFTER THE FACT** — MJ.B had added a second page loop whose text is
   *identical* to the first, so `pageNumber < maxRecoveryPages` matched **twice** (`recover.go:97`
   and `:174`) and the plan's own "assert the mutation matches exactly once" discipline was
   unsatisfiable as written; split into `-COMMIT`/`-EFFECT` with a never-draining fake each, since a
   fake that never drains the commit loop never reaches the effect loop. Fifth instance of that
   shape. **THE BENCHMARK DELTA WAS AN ARTEFACT OF THE RIG**: against the idle-rig M3.C row
   `BenchmarkBrokerFSRead` p95 read 0.7472 → 4.529 ms, a **6.06× "regression"** that would have been
   banked as this item's cost — but the sibling **V1 mission's eval suite was running on the same
   development rig** (ollama + llama-server, 80–98% CPU, load avg 5.22), and the identical
   invocation on the pre-MJ.C parent `b485ead` under the same load reads **4.523 ms**, so the real
   cost is **+0.13%, i.e. zero**. Raised as **4f `w-bench-load-confound`**. **AC5's third clause was
   STRUCK as a doc defect** and **AC11's mutation split** — both recorded in the doc's status block
   rather than quietly checked. Honest claims retained: the Decision-5 residual is verbatim intact
   in `recover.go`, `retryAllowed` has **zero production callers** (so `MUT-RETRY-XOR` proves the
   TEST ROW discriminates, not any runtime behaviour), and the store-side row is a **mirror**, not
   an independently mutation-proven witness. **Known cost, stated:** the two bound tests add ~50 s
   to `host/broker` on every run (local 9.6 → 59.9 s; **CI 55.5 s**, so the Linux runner does NOT
   absorb it and the go job went 36 s → 85 s) — the price of exercising the bound at its real value
   with zero production changes.
   ~~[IN-SPRINT — MJ.A + MJ.B LANDED; ONLY MJ.C REMAINS]~~. **MJ.B LANDED 2026-07-29 (iter-38)**,
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
4d. [**ITEM COMPLETE 2026-08-03 (iter-45) — DG.A + DG.B BOTH LANDED; doc + its `*-sprint-plan.md`
   companion → `design_docs/implemented/w-ddl-gate-teeth.md`.** DG.B shipped via PR **#35** → squash
   **`e6ece55`**, dev CI green **both jobs SHA-addressed on the merge commit and the step logs read**
   (11 `ok` lines → 10 distinct Go packages, 0 FAIL, 4/4 identities across 11 modules, 14/14 named
   tests `failed_tests=0`); evaluator `sonnet` **PASS 96/100, ZERO blocking**, having independently
   reproduced **all ten** named mutations; planner `opus`, executor `codex:gpt-5.6-sol`, designer not
   fired; **`metered=$0.00`**. Four **bisectable** commits (`e13fdfa`/`71fa321`/`9a324d7`/`f0388c6`),
   each green at its boundary, reconstructed from the executor's cumulative `.snap/M<k>/` trees and
   proven byte-identical by sha256 manifest; per-milestone history preserved on
   `sprint/w-ddl-gate-teeth-dgb-impl`. **iter-41's M5 is CLOSED**: `store.Open` no longer accepts a
   structurally stale store and returns `nil`. D5 classifies **before** `schemaSQL` runs; four typed
   errors; schema + `PRAGMA user_version = 1` in **one transaction**; `OpenReadOnly` enforces without
   write/schema/lock; an independent version-1 ledger frozen from `ad619d8`; **CF-O-1 discharged**.
   **Honest tally, verified against the doc (line 709) rather than inherited from the judge: 8 of 9
   production discriminators, AC11 the one test-side probe** — and AC6/AC7/AC8 red on error **TYPE**
   via *secondary* guards (the in-transaction freshness re-check, and the PD-5/PD-6 pragma re-reads),
   which is defence-in-depth beyond spec but is NOT evidence that data would have been overwritten.
   **`4d/OD-7` DISCHARGED**: the ratified read-only sweep found **ZERO World stores on the rig** (13
   candidates inspected through the pinned `modernc` driver, all 0/7 canonical tables; the one
   World-shaped hit is an orphaned `world.db.artifacts` whose database no longer exists) — coordinator
   attestation VERIFIED, no re-park. **`4d/OD-4` (fail vs migrate) stays a recorded FUTURE decision,
   not a blocker**: its own recommendation is alternative 1, fail-loud-only, which is what DG.B is.
   ~~[**DG.B DESIGNED 2026-08-03 (iter-44) — IMPLEMENTATION IS ROUTABLE; ROLLOUT IS PARKED ON `4d/OD-7`.**~~
   PR **#34** → squash **`6b8e77e`**, dev CI green **both jobs SHA-addressed on the merge commit and the
   step logs read** (4/4 identities across 11 modules, 14/14 named tests `failed_tests=0`, 10 distinct Go
   packages `ok`, 0 FAIL); quorum **2 rounds BOTH BLOCKED + narrow-refinement carve-out APPLIED**;
   designer `codex:gpt-5.6-sol` fired **twice** (original + the one sanctioned revision), planner/executor/
   evaluator **not fired** — a design iteration; **`metered=$0.3192`**. **`4d/OD-3` was ratified but
   UNDESIGNED**, so this iteration produced the design: **D5** classifies *before* `schemaSQL` runs,
   **D6** gives legacy/future/invalid distinct typed errors, **D7** commits schema + pin in one
   transaction, **D8** extends the boundary to `OpenReadOnly`, **D9** freezes an independent version-1
   ledger. **AC6–AC14, 9 named mutations, 18 premise rows**; tally corrected twice across the cycle
   (**6/11 → 9/13 → 10/14** production discriminators). **The central problem, measured not assumed**:
   a brand-new store and a legacy store **both** read `user_version=0` and production has **zero**
   `sqlite_master` inspection, so nothing distinguishes them — "reject 0" bricks every existing store,
   "accept 0" is the inert-marker shape OD-3's own text warns against. **V-I…V-S were measured through
   the PINNED `modernc.org/sqlite v1.54.0` driver, not the `sqlite3` C CLI** (iter-43's F4 gap, closed):
   `PRAGMA user_version` **is** transactional (so AC7 is achievable at all), the field is **signed 32-bit
   and truncates silently** (`2147483648` → reads back **0**), and `PRAGMA user_version = ?` is rejected.
   **Three defects, each caught by the quorum**: (1) `gemini-3-1-pro` **and the controller independently**
   found DG.B made DG.A's **AC2 vacuous** — fixed by retaining AC2 verbatim + marked SUPERSEDED,
   re-pointing its mutation to AC10, reassigning the property to AC6, and giving AC2
   `MUT-LEGACY-REJECTION-BYPASS`; (2) `gpt5-6-sol` found D5's table not total over negatives — the
   controller's re-measurement made it **worse than filed** (the truncation above); (3) `gpt5-6-sol` found
   the freshness predicate **unsound** — in `LIKE`, `_` is a single-char wildcard, so a version-0 store
   whose only application object is `sqliteX_probe` counted **0** and would have been **stamped version 1
   over real data**; the `ESCAPE` form counts 1. **CF-O-1**: the same unescaped predicate is in **landed**
   DG.A code (`journal_test.go:877`), so AC1's "unexpected tables fail loudly" cannot see such a table —
   recorded, not fixed; DG.B's executor edits adjacent code. **Next pick for this item: implement DG.B**
   (~0.5–1.0d, planner + executor + evaluator). Item 4d is **still NOT closed**.
   ~~[**DG.A LANDED 2026-08-03 (iter-43) — THE DE-FANG HALF IS DONE; THE ITEM STAYS OPEN ON `DG.B`.**~~
   PR **#33** → squash **`ad619d8`**, dev CI green **both jobs SHA-addressed on the merge commit and
   the step logs read to prove they ran** (10 distinct Go packages `ok`, 0 FAIL; 4/4 identities across
   11 modules; 14/14 named tests); evaluator `sonnet` **PASS 91/100, zero blocking**; planner `opus`,
   executor `codex:gpt-5.6-sol` (ChatGPT subscription bucket), designer **not fired**; **`metered=$0.00`**.
   The inert gate is **replaced**: a hardcoded seven-table canonical manifest + `preJournalSchemaV0`,
   a verbatim `8133573` fixture (sha256 `35f09862e2…`) structurally independent of the embedded
   `schemaSQL`; AC4 deletes the source slice, the sha256 pin, the same-source equality and
   `delete(..., "journal")`. **The de-fang, reproduced first-party by the controller and again by the
   evaluator rather than inherited**: edit `store_heads`' DDL *and* legitimately re-manifest it — the
   one-line action the old pin invited, which used to re-green all 10 packages — and the fresh-store
   gate goes GREEN while `TestOpenAddsJournalAndDetectsStalePreJournalDDL` **REDS** naming
   `store_heads`. No single edit re-greens both halves, because the historical fixture is not derived
   from the thing being edited. **Honest scope: 2 production discriminators (AC1, AC2), 3 test-side
   probes (AC3–AC5)** — "5/5 with named REDs" would overstate production evidence by 2.5×; AC3/AC5's
   mutations remain executor-run, not re-verified first-party.
   **WHAT REMAINS — `DG.B`, RATIFIED BUT UNDESIGNED (this item is NOT closed).** Mark ratified a
   `PRAGMA user_version` pin failing LOUD on binary↔store mismatch, with the frozen-store touch
   ratified; `store.Open` still accepts a structurally stale store and returns `nil` (iter-41's M5,
   untouched). But **OD-3 is a three-alternative decision packet with no ACs, no named mutations and
   no fixtures**, and its own text demands legacy-version-0 treatment plus proven fresh/supported/
   legacy/future cases before ratification-class code lands. Ratifying a DECISION is not having a
   DESIGN — so DG.B needs a **designer + quorum round** before any executor touches `store.Open`.
   **→ DG.B is the recommended next pick for this item.**
   **NUMBERING HAZARD, RECORDED BECAUSE IT NEARLY INVERTED THE RATIFICATION**: Mark's stamp says both
   *"user_version pin — fix NOW"* and *"OD-3 stays declined-as-primary"*. Both are true of **different
   items** — `OD-3` is `user_version` here and *"change the two `scan.go` sites"* in 4e's doc; `OD-4`
   collides likewise. See the **OD registry** in Guardrails.
   ~~[**RATIFIED by Mark 2026-08-03 (attended): FIX NOW → ROUTABLE.**~~ [~~PARKED `needs-human-review`~~ 2026-07-30 (iter-41) — DESIGN DOC + FIVE FIRST-PARTY
   MEASUREMENTS LANDED (PR #30 → squash `d56da6f`, dev CI green both jobs SHA-addressed on the merge
   commit); quorum 2 rounds BOTH BLOCKED, narrow-refinement carve-out deliberately NOT applied;
   `metered=$0.1155`. THE MILESTONE IS SPECIFIED BUT NOT ROUTABLE.]** Doc:
   `design_docs/planned/w-ddl-gate-teeth.md` (designer `codex:gpt-5.6-sol`, rotation; every premise
   row measured by the controller). **All three claimed mutations re-measured at HEAD `e5027df` and
   all three HOLD**, plus **two new measurements the row lacked**: **M4 `MUT-DDL-DRIFT-REPINNED`** —
   editing the DDL *and* pasting back the sha256 the gate's own failure message hands you, the single
   action it invites, re-greens **all 10 packages in one line** while the edit stays unapplied; and
   **M5**, the production-path probe — an existing store opened by a binary carrying the new schema
   gets `store.Open` returning **`err=nil`** with the DDL **byte-identical before and after**. So the
   defect is sharper than *"a gate that cannot fail"*: **the gate's sanctioned repair IS the
   vulnerability**, and a change detector that costs one line to silence is not a correctness gate.
   **Blocked on OD-5**, a collision between two RATIFIED guardrails that a headless loop must not
   resolve: `gpt5-6-sol` rejected twice on the **direction** — a test-only detector leaves the
   verified production fail-open in place, *"axiom noncompliance, not merely an out-of-scope
   enhancement"* — while **frozen-kernel discipline** forbids this loop from changing which on-disk
   stores `store.Open` accepts. Recommended **alternative 1**: land DG.A as designed AND approve
   OD-3 (`user_version`, fail-loud) in principle. **DG.A** (test-only, ~0.25–0.5d, 5 ACs each with a
   named production/test-classified mutation) is fully specified and waits on ratification =
   **CF-L-1**. The quorum's ledger: 4 objections, **3 right** (fixture provenance, table
   enumeration, determinism, and a baseline-drift premise nobody had checked — measured, it HOLDS),
   **1 wrong in its prescription** (`name NOT LIKE 'sqlite_%'` excludes **zero** rows today; the live
   limb is `type='table'`, which drops 7 `sqlite_autoindex_*` rows — adopting it would have logged a
   dead mechanism as live, iter-38's defect inside its own remedy). — originally raised by the
   controller's own Gate-2 audit at iter-36, with the evidence attached; NOT folded into 4c, because
   it is a pre-existing defect in LANDED `w-store-durability` code and 4c's AC1 deliberately does
   not depend on the broken gate. ~~Pick it when a DDL change is next contemplated, or sooner if the
   queue allows.~~ **That deferral condition was backwards and iter-41's M5 is why**: the fail-open
   fires on the FIRST DDL edit anyone makes, so repairing the gate "when a change is contemplated"
   means repairing it after the change that needed it.]
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
4e. [**ITEM COMPLETE 2026-08-04 (iter-46) — RG.A LANDED; doc → `design_docs/implemented/`; `4e/OD-1` AND `4e/OD-2` BOTH DISCHARGED.**
   PR #36 → squash `f19acac`, dev CI green **both jobs SHA-addressed on the merge commit and the step
   logs read to prove the `-race` leg actually ran**; evaluator `sonnet` **PASS 96/100, zero
   blocking**, having independently reproduced AC3 both legs, AC2 leg 1, AC4 and two mutations with
   its own sha256 proofs. **8/8 acceptance criteria met; honest tally 4 of 8 production
   discriminators** (AC1/AC2/AC3/AC4), the evaluator's 3/8 dissent recorded rather than rounded away.
   **The repo now has a `-race` leg for the first time** — item 4e's root enabler was that
   `grep -rn race .github/workflows/ scripts/` returned nothing — and it ships with a **known-positive
   control**: a nested `racecontrol/` module (invisible to `./...`, still exactly 10 packages) whose
   deliberate race the gate runs first, aborting with *"the race detector is not armed; every 0-races
   result in this gate is void"* if no `WARNING: DATA RACE` appears. That was **not** in the design:
   without it a green `0 races` is unfalsifiable, which is this mission's signature defect wearing the
   remedy's clothes, and it is what upgrades `MUT-RACE-LEG-DROPPED` from a drift check to a production
   discriminator. `verify_go.sh` **sets nothing** — it reads `go env GOVERSION` and refuses an
   affected version — because both quorum reviewers had killed drafts that exported or
   assign-if-unset the toolchain, each of which is a gate that cannot fail. **`4e/OD-2` DISCHARGED:
   filed upstream as [`golang/go#80706`](https://github.com/golang/go/issues/80706)** after a
   duplicate search came back empty *with its control firing*. **AC6 ANSWERED AND IT CUTS AGAINST THE
   DOC'S OWN MOTIVATION: linux/amd64 is NOT affected** — all four affected toolchains return `OK` on
   `ubuntu-latest`, so CI was never building through the miscompilation and the historical blast
   radius was bounded to local darwin/arm64 builds. Recorded in `bench/BASELINE.md` and the doc rather
   than left in a step log nobody is required to read. **THE DOC WAS WRONG TWICE AND BOTH ARE ON THE
   RECORD**: `MUT-CANARY-BLIND` does not "pass on BOTH toolchains" (it is `ok` on 1.26.4 and FAILS on
   1.25.6 — refuted independently by planner, executor, controller and evaluator), and
   `MUT-PIN-REMOVED` targeted a `GOTOOLCHAIN` export that Decision 2's own round-2 revision had
   already removed. `4e/OD-3` stays DECLINED and is honoured — `scan.go` is byte-untouched.
   `4e/OD-4` stays OPEN as the named escalation path for a CI `-race` overrun, its cost now
   re-measured as a **69-175 s spread** rather than the doc's single figure.
   ~~[**[MECHANISM IDENTIFIED 2026-07-30 (iter-40) — DOC + REPRODUCTION FIXTURE LANDED, quorum 2 rounds
   + narrow-refinement carve-out revision applied; OD-1 RATIFIED by Mark 2026-08-03: pin Go 1.25.6; OD-2 AUTHORIZED: file the reproducer publicly at golang/go → RG.A ROUTABLE** (~~REMEDIATION PARKED on OD-1~~)]**
   — **the mechanism is a Go compiler regression, and the standing `modernc.org/sqlite` hypothesis is
   REFUTED.** go **1.26.0 → 1.26.5** (1.26.5 is the LATEST STABLE, so there is nothing to upgrade
   forward to) on darwin/arm64 miscompile a **local array literal indexed by a `range` variable**
   assigned into a struct string field, inside a composite literal that also carries an interface
   method call — emitting a **corrupt string header** (nil data pointer, garbage length; printing it
   SIGSEGVs). **It does not need `-race`**: a 52-line dependency-free program reproduces under a plain
   `go build`. `-gcflags=all=-N` fixes it, `all=-l` does not (so: the optimizer, not inlining); go1.25.6
   and go1.24.9 are clean. **ONE mechanism explains BOTH symptoms** — disabling optimizations for the
   single package `host/store` clears the failure *and* the hang, which is what iter-39 demanded of any
   hypothesis. Census: **exactly two** production sites carry the shape, `host/store/scan.go:74` and
   `:112`, precisely the two functions the two symptoms belong to. **`scan.go` is correct Go and needs
   no correctness fix.** The `-race` leg this item asks for is **viable and cheap**: with
   `GOTOOLCHAIN=go1.25.6` + `AILANG_BIN` set, `go test ./... -race` is **rc=0, 10/10 packages, 179 s,
   ZERO data races** (`host/store` 17.1 s green, against FAIL+hang at 124.7 s on go1.26.4 in the same
   worktree minutes apart) — so the blocker was never the race detector. Doc
   `design_docs/planned/w-race-gate-blindspot.md`; re-runnable fixture carrying its own known-positive
   controls at `design_docs/verification/w-race-gate-blindspot/`. **PARKED FOR MARK: OD-1** (lower
   `go.mod` 1.26.4 → 1.25.6 — ratification-class, recommended YES) and **OD-2** (file the reproducer
   upstream at `golang/go` — a public post to a third-party project, recommended YES, repro ready).
   Also measured: `go.mod` **cannot** pin an older toolchain (both `go` and `toolchain` are floors —
   under `toolchain go1.25.6`, `go version -m` stamps the binary `go1.26.4` and it still reproduces);
   only `GOTOOLCHAIN` can. CI's resolved toolchain is a **measured** `go1.26.4 linux/amd64` (step log,
   run `30483249118`), but **amd64 exposure is UNDETERMINED** — no Rosetta on this rig, so AC6 runs the
   fixture in CI to settle it. **Original filing retained below.** — found by chasing a line the MJ.A
   executor mentioned in passing; **pre-existing in LANDED `w-store-durability` (SD.A, `86d1276`)
   code that two zero-blocking judgements never saw**. Sized small (~0.25–0.5d). Pick alongside 4d.
   **WIDENED at iter-39 by the MJ.C judge, and the widening matters: this is TWO distinct symptoms,
   not one restated.** Iter-37/38 recorded `TestScanUnreadableLogKeysetResumes` **failing** 5/5
   under `-race`; the MJ.C judge independently reported `TestScanUnreadableWorldsFindsPoison`
   **hanging** under `-race` (killed at 60 s). Both functions exist and are distinct
   (`host/store/scan_test.go:32` and `:83`, confirmed first-party with a known-negative control),
   and neither is MJ.C's — `git diff b485ead..460ade3 -- host/store/scan_test.go` is empty. So the
   mechanism hypothesis must explain a **failure AND a hang in the same file**, which is further
   evidence for the memory-corruption reading over a logic bug. **Also newly relevant to costing
   this item**: `host/broker` under `-race` now takes **163 s** (judge-measured) against iter-38's
   4.095 s, because MJ.C's two bound tests exercise 2^20 pages each — so any future `-race` leg
   must bound or skip them, and 4e's "add a `-race` leg" option is now materially more expensive
   than when it was queued.]
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
   (UNVERIFIED) is ~~`modernc.org/sqlite`'s heavy `unsafe` usage plus `-race`'s altered memory layout
   exposing a latent corruption~~ **[REFUTED iter-40 — a 52-line program with ZERO dependencies and no
   `-race` reproduces it]**, or **a go1.26.4 darwin/arm64 toolchain bug [CONFIRMED iter-40, and wider
   than stated: 1.26.0 through 1.26.5]**. **The work**: identify
   the mechanism before prescribing the fix — because zero `DATA RACE` output means "just add
   `-race` to CI" is not yet known-good — then either add a `-race` leg with the hang bounded, or
   record in writing why this repo does not run one. **[iter-40: mechanism DONE; "just add `-race`"
   turns out to be correct-but-for-the-wrong-reason — it is known-good ONCE the toolchain is pinned,
   and the pin is OD-1, parked for Mark.]** **Sixth instance of this mission's signature
   shape** (a gate that cannot fail), and the first where the gate cannot fail *because it is never
   invoked at all* — **and iter-40 adds a new member to the family: the COMPILER is an instrument too,
   so every gate in this repo (including the mutation protocol) reports through it.**
4f. [**ITEM COMPLETE 2026-08-05 (iter-50) — BOTH MILESTONES LANDED AND ALL 16 NAMED MUTATIONS DISCHARGED; doc → `design_docs/implemented/w-bench-load-confound.md`.** BC.A′ `0b72019` (iter-48) · BC.B′ `d357474` (iter-49) · controller pass **`C2b`** (iter-50) discharged the last five: `MUT-AB-FLOOR-SPLIT`, `MUT-PROBE-CALLER-DIR`, `MUT-PAIR-TWO-SESSIONS`, `MUT-PAIR-SEQUENTIAL`, `MUT-PAIR-INLINE-BUILD`. **No code changed** — this pass produces EVIDENCE; the two files it mutated were restored byte-identically (`shasum -c` OK) with `git status --porcelain` verified empty after every arm, and the work is structurally **controller-only**: the daemon benchmarks bind loopback, which the executor's `workspace-write` sandbox denies. **THE FIND — A RUN MEANT TO CONFIRM A DOC PREDICTION REFUTED IT INSTEAD (P47).** `MUT-PAIR-INLINE-BUILD`'s rule-based RED fired perfectly (**R1 on both sections + R2's orphan cascade on all four raw blocks, 6 violations**), but its stated *secondary observable* — *"leg-1 elapsed jumps from seconds to a compile-bearing figure"* — did **NOT**: honest legs **7,7,7,7 s** vs the inline-build's **8,8,9,8 s**, because a warm Go build cache prices a full compile of these trees at **1–2 s** (measured via the recorder's own `prebuild_elapsed_s` across three sessions), inside the doc's own ~1.4× within-condition noise. **A secondary observable a cache can erase is worse than none — it gives a reviewer a plausible reason to stop looking.** Struck in the mutation bullet, not absorbed. **`AC6` DISCHARGED WITH ITS VACUITY ARM MEASURED RATHER THAN INHERITED (P45)**: under `GOTOOLCHAIN=auto` the fixture straddle reads variant `go1.26.5` / control `go1.26.4`; under `GOTOOLCHAIN=go1.25.6` **both read `go1.25.6`** — so `4f/CF-S-2` was honoured *and* proven load-bearing in the same session. `MUT-AB-FLOOR-SPLIT` → `rc=1`, **exactly one** violation, `toolchain mismatch inside claimed A/B pair`, against a `rc=0 ✓ PASSED` control arm on the un-appended file. `MUT-PROBE-CALLER-DIR` → the caller-cwd probe records `go1.26.4` for BOTH trees and the gate **greens a genuinely cross-toolchain pair**, reproducing the round-1 bug and proving the `-C` placement is what makes the known-positive fire. **`AC7` DISCHARGED, AND ONE ARM FIRED A LIMB THE DOC NEVER NAMED (P46)**: every arm ran against an honest same-session control pair that greened first (`2 well-formed pairs`), so each RED is the edit and not the fixture. `MUT-PAIR-TWO-SESSIONS` → **2 violations, `unpairable conditions block` on both spliced halves, R4b SILENT** (the round-4 supersession validated on a real pair); its named silencing attempt relocates to **R4c + R4b + an unnamed R4f**, because two sessions recorded **36 s apart** occupy disjoint time windows — time is the field a splicer cannot forge without also forging the ordering. `MUT-PAIR-SEQUENTIAL` → **exactly one** violation, R4f, on an otherwise well-formed pair. **THE CENSUS WAS A TRANSCRIPTION, NOT A MEASUREMENT (P48)**: re-derived by command, the plan's `mutation_tally` is **16** (BC.A′ 3 + **BC.B′ 13, not 12**); the *"8 of 12"* below was a quantity quoted without a command, and `8+5=13>12` was visible on the page. Item-wide the census is now **16 of 16**. **P47b**: iter-49's P3/P4 citation repair was recorded inside P42's own row and never applied to P3/P4, which still served `ci.yml:88-89` under a `CONFIRMED` verdict with no marker — **a supersession that lives only in the superseding row is a note to the person who already knew**; markers applied in place. **CARRIED OUT OF THE ITEM, NON-BLOCKING**: `4f/CF-R-3` (the carry-forward ID-collision sweep is partial and its zeros uninformative — a namespace-hygiene instrument, not an acceptance criterion) and the evaluator's cosmetic NB-1/NB-2/NB-4/NB-5. Prior iter-49 text follows, unedited per the no-prescience convention. ~~[**IN-SPRINT 2026-08-04 (iter-49) — MILESTONE `BC.B′` CODE LANDED; THE ITEM IS **NOT** CLOSED — `AC6` IS UNDISCHARGED AND 5 RE-RECORDING MUTATIONS ARE CARRIED TO THE NEXT `C2b` PASS.** BC.B′ = the claim-structure gate `--check-claims` (R1–R6, R4's six limbs, frozen evaluation order, hardcoded path, no knobs), the `bench/BASELINE.md` policy header + 3 legacy markers + the amortisation-pointer correction, two `go-verify` CI steps, and a real controller-recorded acceptance pair (PR #39 → squash `d357474`, dev CI green **both jobs SHA-addressed on the merge commit** and corroborated by a direct per-workflow read at the same SHA; +457 lines across exactly the 3 in-scope files; one bounded codex run M2a; evaluator `sonnet` **PASS 77/100, zero blocking**; **`metered=$0.00`**). **THE FIND — AND IT WAS ONLY REACHABLE BY RECORDING A PAIR FOR REAL**: appending a freshly emitted pair turned the gate RED. Independent recompute proved the EMISSION correct in every field (both `conditions_sha256`, **4/4** `legN_output_sha256`, one shared `pair_id`, parent edge), so the CHECKER was wrong — `one = lambda key: values[key][0] …` closed over the **loop variable**, and Python late-binding made every block's accessor read the **LAST** block's fields. R3 compared the variant's legs to the control's hashes, R4d saw roles `['control','control']`, and **`R4c` — the section-locality revision round 4 exists to add — was SILENTLY VACUOUS**. Fixed by binding `values` as a default arg; re-arming the bare closure reproduces the identical 3 REDs and the restore is byte-identical. Nothing upstream could have caught it: `DD-3` makes B1 green with **no pair in the file**, M2a's mutations are single-section or file-level, and the executor sandbox cannot bind the loopback sockets a recording needs. **B2 lands the pair, so CI now evaluates the pair rules on every push** — the ubuntu step log prints `checked: 7 raw benchmark blocks, 2 conditions blocks, 1 well-formed pairs, 3 legacy markers`. **PAIR-RULE BATTERY AFTER THE FIX, EVERY DOC PREDICTION HELD**: `MUT-CLAIM-NONPARENT` **dual-fired** (R4c binding → R4a `A/B pair is not variant-vs-parent`); `MUT-PAIR-ID-SPLIT` REDded R4c on the EDITED section + R4d on BOTH with **R4b not firing**; `MUT-CLAIM-TOOLCHAIN-SPLIT` → R4e alone; `MUT-CONTROL-REUSE` → R4d; `MUT-EDIT-RAW-NUMBER` → R3 naming block and leg, no expected hash. **8 of 12 BC.B′ mutations discharged.** **`4f/CF-M-1` MEASURED BOTH WAYS**: a Linux-like `sysctl` stub → exactly `probe FAILED: sysctl -n hw.ncpu`, step PASSES; `AILANG_BIN` unset → still non-zero, still a generic `probe FAILED` (count 1), specific marker **absent**, step FAILS. **`4f/CF-S-1` DISCHARGED** (the checker parses `ailang_pin` and the repeated-key `legN_competing` lines, exempting exactly those two from its uniqueness assertion). **NEW `4f/CF-S-3`**: `MUT-CLAIM-TOOLCHAIN-SPLIT`'s doc literal (`goversion` → `go1.25.6`) is a **no-op** now that both sections record `go1.25.6`, so the mutation as written would be vacuous — the controller used `go1.26.4`; same class as P41, independently found by the evaluator (NB-3). **NEW `P42`**: the freshness sweep from the doc's OLDEST declared base (`c1e6125..HEAD` = 9 files vs `ea5e405..HEAD` = 1) found a **fourth face** of iteration 46's drift — `P3` cites `ci.yml:88-89`, which has read `:101-102` since RG.A and was therefore *already stale at `ea5e405`*, the base iter-48 swept from; `P4` drifted `:53`→`:55`. Both claims survive; only their citations rotted. The **planner** had already measured `:101-102` into the sprint plan, so the plan was fresher than the doc. **STILL OUTSTANDING FOR `C2b`**: `AC6`'s two mutations (`MUT-AB-FLOOR-SPLIT`, `MUT-PROBE-CALLER-DIR`) via the floor-split fixture worktrees — **`4f/CF-S-2` is BINDING there: invoke the fixture session `GOTOOLCHAIN=auto` explicitly or `MUT-AB-FLOOR-SPLIT` passes vacuously (P41)** — plus `MUT-PAIR-TWO-SESSIONS`, `MUT-PAIR-SEQUENTIAL`, `MUT-PAIR-INLINE-BUILD`; `4f/CF-R-3` still carried. Prior iter-48 text follows, unedited per the no-prescience convention. ~~[**IN-SPRINT 2026-08-04 (iter-48) — `4f/OD-8` ANSWERED (attended: EVIDENCE, NOT CLAIMS) AND MILESTONE `BC.A′` LANDED; `BC.B′` IS THE NEXT MILESTONE AND IS ROUTABLE NOW.** BC.A′ = the single-session pair recorder (PR #38 → squash `0b72019`, dev CI green **both jobs SHA-addressed on the merge commit**, corroborated by a direct per-workflow read at the same SHA; +414 lines in `scripts/bench_worldd.sh`, two bounded codex runs M1a/M1b; evaluator `sonnet` **PASS 87/100, zero blocking**; `metered=$0.00`). **THE CONTROLLER PROVED WHAT THE SANDBOX COULD NOT**: the daemon benchmarks bind loopback, so every leg was `UNINFORMATIVE UNDER SANDBOX` to the executor; the controller ran a real session — **rc=0 in 18 s**, 2 conditions + 4 raw blocks — and independently recomputed **every** integrity field (`pair_id` **section-locally for both sections**, `conditions_sha256` both, **4/4** leg hashes matching real raw blocks, shared `pair_id`, parent edge, genuine interleave control `1/4`+`4/4` / variant `2/4`+`3/4`), with a one-hex-digit `pair_id` edit failing to verify as the control. The all-or-nothing property is **structural, not janitorial** and was demonstrated with a leg-3 shim whose mutation was proven applied by differing sha256 before the result was believed. **TWO EVALUATOR DEFECTS FIXED RATHER THAN CARRIED** (`59852ef`): three sites hardcoded `GOTOOLCHAIN=go1.25.6` **inline**, which beats the caller's exported value — so the recorder was **selecting** the toolchain it exists to **record**, collapsing AC6's straddle exactly as P41 predicts and violating the doc's own *"it selects nothing"* clause; **the fault was the controller's own directive**, which told the executor to pin every `go` command. Re-verified under both regimes (`auto` → `go1.26.4`, pinned → `go1.25.6`). And `ailang_pin` emitted **seven** lines from the `--version` banner into a line-oriented block; collapsed at capture, 0 unparseable lines measured. **`4f/CF-R-1` + `4f/CF-R-2` DISCHARGED** (premise rows repaired in fact rather than in prose — and the sweep found **2 of 3**, the planner catching `P6`; AC6 re-pointed at the real straddle AND given a regime clause); **`4f/CF-R-3` carried**; **new `4f/CF-S-1`** (BC.B′'s checker must parse `ailang_pin` and the repeated-key `legN_competing` lines) and **`4f/CF-S-2`** (the AC6 fixture now needs an explicit ambient-regime invocation, since the recorder no longer pins). **P41 IS THE ROW TO READ BEFORE BC.B′**: under `GOTOOLCHAIN=go1.25.6` — what `verify_go.sh` requires and CI pins job-wide — `go env GOVERSION` returns the PIN for every tree, so `MUT-AB-FLOOR-SPLIT` comes back **GREEN**; and no deny-list-free straddle exists on this rig, because floors and `toolchain` directives select only UPWARD. Prior park text follows, unedited per the no-prescience convention. ~~[**PARKED `needs-human-review` on `4f/OD-8` 2026-08-04 (iter-47)**~~ — Branch-A revision landed (PR #37 → squash `2529d4f`, dev CI green **both jobs SHA-addressed on the merge commit**, direct per-workflow read matching); doc 793 → **1713 lines**, 36 premise rows, 16 named mutations, designer `claude:claude-fable-5` (rotation) fired **three times** — original branch-A revision + two sanctioned fixes — each rc=0 touching **only** the target file; quorum **rounds 3/4/5, all BLOCKED**; `metered=$0.3007`. **WHAT THE MECHANISM NOW IS**: one bounded `--record-pair --variant <dir> --control <dir>` session, both binaries **prebuilt** so compilation leaves the measured window, a unique pair ID, the frozen `control/variant/variant/control` leg order, per-leg utc/load/ps/hash, **emission only after leg 4** (a session failing on leg 3 emits zero fences *by construction*, not by cleanup), R4 grown to **six named limbs** under a frozen evaluation order, and **closed-world bounding** — every external-binary invocation through `run_bounded`, not just the `go` ones. **THREE CONTROLLER MEASUREMENTS REFUTED THE DOC**: a full 10-benchmark leg on a prebuilt binary is **6–8 s**, not the doc's ≈155 s (that figure was *compile*-dominated, so one deadline was bounding two populations with different physics; a complete four-leg session measured **27 s** — recorded as a **SAMPLE, not a bound**); `go test -c -o BIN <other-worktree>/host/daemon/` **FAILS** `outside main module` so only the `go -C` form works; within-condition noise is **~1.4×** at essentially constant load, which is the measured reason the excluded third limb is not derivable here. **ROUND 4** caught a pair-scoped `pair_id` derivation sitting under *section*-scoped mutation predictions with the unpairable block undefined — one implementer-reading from a silent skip; fixed, three predictions re-derived. **ROUND 5** caught the remedy for unbounded waits containing **its own unbounded waits** (`sysctl`/`ps`/`git`/`python3`) — the **second instance of that trap inside this one document** — and a **dimensional error in the doc's own honesty sentence** (duration bounds temporal *separation*, not skew *magnitude*); both adopted in full. **WHY IT PARKS**: `gpt5-6-sol`'s remaining limb disputes the design DIRECTION, so the narrow-refinement carve-out does not apply — and the question is genuinely new rather than a re-ask of OD-6 → **`4f/OD-8`**, recommendation (1) with the reviewer's own round-2 fallback as its argument. **ALSO RECORDED: P9 and P22 are STALE against HEAD** — `go.mod` floor is `1.25.6` and `4e/OD-1` discharged at `f19acac` — which **five quorum rounds did not catch**, and which voids the stated reason for deferring limb (iii). ~~[**OD-6 RATIFIED by Mark 2026-08-03 (attended, see STATUS): BRANCH A — interleaved --record-pair, pair ID, control-reuse rejection; re-sized ~0.6–0.9d → ROUTABLE.**~~ [~~PARKED `needs-human-review`~~ 2026-08-03 (iter-42) — DESIGN DOC LANDED (PR #31 → squash
   `b986c7a`, dev CI green both jobs SHA-addressed on the merge commit and step-log-verified);
   quorum 2 rounds, r1 BLOCKED both-reject, r2 BLOCKED 1-pass/1-reject; narrow-refinement carve-out
   deliberately NOT applied; `metered=$0.2104`. NO MILESTONE IS ROUTABLE UNTIL OD-6 IS ANSWERED.]**
   Doc: `design_docs/planned/w-bench-load-confound.md` (designer `claude:claude-fable-5`, rotation;
   793 lines, **25 first-party Verification Log rows**, **10 named RED mutations**, 2 milestones).
   **Blocked on OD-6**, a scope judgment a headless loop must not make: `gpt5-6-sol` holds that the
   design's **central claim** is unsupported, and its own `proposed_fix` **branches two ways that
   deliver different items** — **branch A** (single-session interleaved `--record-pair` with a pair
   ID and control-reuse rejection; keeps the charter row's promise; ~0.6–0.9d, i.e. the item
   outgrows this row's original sizing) vs **branch B** (keep the mechanism, weaken the claim to
   *"mechanically complete evidence, not a mechanically valid cost claim"*). Recommended **branch A,
   bounded**, taking **neither** the reviewer's third limb — a *measured within-pair load-divergence
   acceptance rule* — because no data on this rig can defensibly derive that threshold and it is the
   noise-gate `BASELINE.md:7-8` already calls **dishonest (S6)** wearing a new name.
   **THE DEFECT, verified first-party against the doc's own text:** R4 constrains only
   `control.commit == variant.parent` plus identical `goversion`, `goos_goarch`, `ncpu`, `hw_model`;
   a search for any temporal, load-comparability, pair-identity or control-reuse constraint returns
   **nothing**, while the known-positive control in the same call confirms the schema **does** record
   `utc`, `load_before`, `load_after`. **The recorder captures the load and the gate that blesses the
   pair never reads it** — so an idle variant pairs with a control recorded days earlier under heavy
   load and passes, reintroducing the 6.06× artefact class inside the artifact built to eliminate it.
   Seventh instance of this mission's signature shape; **third inside this one document's evolution**.
   **Both round-1 objections were right and both were adopted**: the recorder's unbounded `go test`
   (Standing Rule 6 — fixed by mirroring `run_bounded` from `scripts/verify_ail.sh:61-74`, with
   `MUT-REC-STALL` proving the kill reaches a grandchild), and a caller-cwd `go env GOVERSION` — the
   controller **measured** that premise rather than forwarding it (`GOTOOLCHAIN=auto` live; floor
   1.26.5 module → `go1.26.5`, this repo → `go1.26.4`) and the measurement made it sharper than
   filed: **OD-1 proposes lowering this repo's floor**, so the first A/B pair straddling it is exactly
   the cross-toolchain pair R4 exists to red, which the round-1 design would have greened by
   comparing the variant's toolchain to itself. **CF-K-2 is FOLDED IN** (R4's toolchain-identity limb
   + `MUT-CLAIM-TOOLCHAIN-SPLIT` + `MUT-AB-FLOOR-SPLIT`) and ships when 4f ships. Limb **(iii)** (the
   amortisation re-derivation) is scoped **OUT**, blocked on OD-1 — banking numbers under an
   unresolved toolchain condition is precisely the defect this item exists to fix. New carry-forwards
   **CF-M-1** (D4's CI assertion must grep the *specific* `hw.ncpu` marker) and **CF-M-2** (P25's
   deadline arithmetic is designer-measured, not controller-re-run cold). **BC.A is very nearly
   branch-independent** and is the routable half if the queue needs motion before OD-6 — but that is
   the human's call too, since branch A reworks the interface BC.A ships.
   ~~[QUEUED 2026-07-29 (iter-39)] — raised by the controller's own Gate-4 bench re-measure,
   with the A/B evidence attached and already committed to `bench/BASELINE.md` in `460ade3`. Sized
   small (~0.25–0.5d). Pick alongside 4d/4e.~~ **Iterations 40 and 41 both deferred this row on the
   ground that CF-K-2 would bank a toolchain condition OD-1 is parked to invalidate. That reasoning
   is sound for limb (iii) ONLY** — limbs (i) and (ii) are toolchain-independent, and recording the
   toolchain is *more* valuable while OD-1 is pending, not less.]
   **w-bench-load-confound** · clause-1 · **the benchmark baseline's comparability promise is
   unenforced, and the confound is a SIBLING MISSION on the same rig.** `bench/BASELINE.md` says
   *"Later sprints diff against this file on the same development rig"* and reasons explicitly about
   honesty (*"Noise-gating a shared runner would be a dishonest gate (S6)"*) — but it considered CI
   noise and not that the **development rig is shared with the V1 mission's eval suite**, which runs
   `ollama` + `llama-server` at 80–98% CPU on a schedule this loop does not control and cannot see.
   **Measured at iter-39**: with V1's eval suite live (`load averages: 5.22 4.99 5.91`),
   `BenchmarkBrokerFSRead` p95 read **4.529 ms** against the idle-rig M3.C baseline of **0.7472 ms**
   — a **6.06×** apparent regression that MJ.C would have banked as the effect journal's cost. The
   same invocation on the pre-MJ.C parent `b485ead` under the same load reads **4.523 ms**, so the
   true cost is **+0.13%**. Nothing in the file, in `scripts/bench_worldd.sh`, or in the harness
   records the rig load at measurement time, so the confound is **invisible by default**; it was
   caught only because the ratio looked implausible enough to warrant a control. **This is the
   signature shape at the benchmark**: a measurement that would report the same thing whether or not
   the property held. **The work**: (i) record load average + a named list of competing processes
   into `BASELINE.md` at measurement time, emitted by the harness rather than typed by hand;
   (ii) decide and write down the policy — either an A/B against the parent commit becomes MANDATORY
   for any cost claim (cheap, load-independent, and what actually worked here), or the harness
   refuses to record when load exceeds a threshold (simpler, but wastes iterations on a rig this
   loop does not own); (iii) the amortisation section in `BASELINE.md` is currently pinned to the
   M3.C **idle-rig** numbers and labelled so — re-derive it once a clean-rig invocation is available.
   **Prefer (ii)-as-A/B**: a threshold gate assumes an idle rig will eventually happen, which on a
   two-mission rig is not guaranteed, whereas the A/B is correct under any load.
5. [**STILL BLOCKED — but on ONE prerequisite now, not three; RE-MEASURED 2026-08-05 (iter-50) rather than inferred from the `#498` stamp.** The attended `#498` SEAM stamp cleared prereq **1** (upstream seam, verified live on released v0.33.0), and item 4b's landing cleared prereq **3** (`Commit.InvocationID` + `GetReceipt` + the three-state receipt law ARE the atomic not-started-vs-committed contract, the stable idempotency ID and the queryable durable receipt that prereq asked for). Prereq **2 is HALF clear, and the missing half is the one the item is named for**: `host/broker/broker.go:45` now defines `type Session struct` with `NewSession(store, episodeID, grants, registry)` — the broker session API exists — but a repo-wide `grep -rniE '[Tt]ransition[ -]?[Rr]egistry' host/ world/ cmd/` at `de80792` still returns **ZERO**, with the same-call known-positive control (`registry` in `host/registry/registry.go` → **25**) firing, so the absence is a measurement and not a failed grep. `host/registry` remains the *interpreter epoch* registry, a different thing. **This row was one stamp away from being promoted to `[NEXT]` on the strength of a single satisfied prerequisite; the re-measurement is what stopped it.** Whoever picks this next writes the transition registry first, or re-scopes P6.B around its absence. Prior text follows. ~~**BLOCKED 2026-07-28 (iter-24) — DOC LANDED + 2 QUORUM ROUNDS + carve-out revision applied.**~~
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
6. [PARKED until 2–5 land; design-lens note added 2026-07-31: the doc, when authored, carries
   the compositional-generalization hypothesis (Zhang 2026, REFERENCES.md) — the World arm's
   typed structure may IMPROVE generalization, so the experiment should record per-task-family
   splits, not only aggregate pass-rate; the ratified thresholds are UNCHANGED]
   **w-agent-floor-m4** · clause-4 · dual-reference NON-INFERIORITY
   floor: Claude Code + codex, shell arm vs World-MCP arm, paired N≥3, stability precondition
   checked first; motoko as optional third arm if eligible; report honestly; park World if the
   floor fails on eligible agents · ~3d (was w-motoko-m4 → w-value-gate-m4 → renamed with the
   2026-07-23 floor reframe; clause-5 provenance teeth carry the value burden)
6b. [**LANDED AS DESIGN, BINDING — §7 RATIFIED by Mark's attended triple ratification
   (STATUS 2026-07-28 ~09:00; row flipped by the coordinator 2026-07-31, closing the
   bookkeeping lag). Unblocks item 7 (w-approval-inbox) the moment item 5's MCP work lands.**
   ~~PARKED `needs-human-review` 2026-07-28 (iter-26)~~ — PICK-TIME QUORUM COMPLETE (2 rounds,
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
8. [**[IN-SPRINT] — UNCHANGED THIS ITERATION; `SM.B2a` IS STILL THE NEXT MILESTONE BUT NO LONGER THE QUEUE HEAD. Item 10 was promoted ahead of it 2026-08-06 (iter-56)** on a measurement iter-55's own row did not record: the boundary gate's `defer`-based restore does not survive SIGKILL or Go's `-test.timeout` panic, so a killed run leaves a forbidden import in a production source **permanently** — invisible to `go build` (rc=0) and reddening the boundary gate itself, which then accuses an innocent file of a **network-boundary violation**. That failure mode is indistinguishable from `SM.B2a` genuinely violating the boundary, during the one milestone that adds network code, which is why it is ordered first. See item 10. **Prior head text follows, still accurate for the milestone itself.** ~~**[IN-SPRINT] — `AC12` REPAIRED 2026-08-06 (iter-55), PR #45 → squash `1761a9c`, dev CI green BOTH jobs SHA-addressed on the merge commit. The next unit of work is STILL milestone `SM.B2a`, gated on nothing — it now lands against a boundary gate with teeth.**~~ **THE CARRY-FORWARD THE LAST TWO ITERATIONS WROTE DOWN WAS DISCHARGED EARLY, AND IT FOUND A HOLE.** Both iter-53 and iter-54 recorded that `AC12`'s *"network confined to `host/broker`"* control is VACUOUS until SM.B2a and must be re-asserted there. Re-asserting it at the boundary — BEFORE the network code exists — showed the gate was weaker than vacuous in a second, unrecorded way: the loopback exception is true of **exactly one** protected group and a single shared `forbiddenImportPrefixes` list granted it to **all three**. Measured, mutations confirmed landed by sha256, restores byte-identical: bare `net/http` blank-imported into `host/store/store.go` → **rc=0 PASS**, into `host/replay/replay.go` → **rc=0 PASS**, while the `net/http/httputil` control correctly REDs. **Every protected group's mutation was `httputil`** — so the gate had only ever been tested against a mutation shaped to itself, iteration 54's own spine arriving inside the gate iteration 53 landed. The exemption was also **UNFORCED**, not merely unnecessary: baseline `net/http` presence per dependency closure is `host/store` **0** (160 deps), `host/replay` **0** (162), `cmd/ailang-worldd` **1** (233), control `host/hashref` **0** — only `cmd/ailang-worldd` ever needed it, exactly as its own code comment says. Repaired with a per-group `extraForbidden` field plus `TestBareNetHTTPExemptionIsPerGroup`, which reds if the asymmetry is collapsed back into one global list (proven non-vacuous by setting `host/store`'s entry to `nil` → RED naming the group, restore byte-identical). Post-repair: both threat arms RED, `httputil` still RED, `cmd/ailang-worldd` still PASS, pristine green. **WHAT REMAINS FOR `SM.B2a` TO RE-ASSERT:** the positive half — that network code, once it EXISTS in `host/broker`, is genuinely permitted there — is still unproven, because `host/broker` has zero `net/http` deps at HEAD. The gate's green control asserts only that the broker's dependency closure is NON-EMPTY (`:281`), which is true of any Go package. **SEPARATE, PRE-EXISTING, AND NOW QUEUED AS ITEM 10:** this same gate mutates three other packages' production sources in the LIVE tree while `go test ./...` builds them concurrently — measured on pristine `dev`, and it red-lit `TestCLIRealSubprocessEpisode` once during this iteration. Prior text follows.] ~~**[IN-SPRINT] — `SM.B1` LANDED 2026-08-05 (iter-54), PR #43 → squash `1856bfb`, dev CI green BOTH jobs SHA-addressed on the merge commit. The next unit of work is milestone `SM.B2a`, gated on nothing.** **What SM.B1 SETTLED:** `approval_claims` + an atomic `AppendClaimedEffectIntent` exist, schema is at version **2**, and **`DD-3` is closed loudly** — `enforceSchemaVersion`'s bare `return nil` became `*LegacySchemaVersionError`, so a store left at `user_version = 1` can no longer open successfully and silently skip `schemaSQL` (which would have surfaced `approval_claims` as absent at the moment of the irreversible publish). All three independent fixtures moved in ONE commit with `schemaV1SQL` byte-unchanged; the DDL gate redded mid-milestone exactly as designed and the recorded RED list is in the log. **What SM.B1 newly OPENED — the milestone-gating ledger check was VACUOUS AS DELIVERED, and the executor's own mutation could not see it.** `TestSchemaVersionLedgerIsIndependent` greps its own source; its two NEGATIVE needles were split so they would not match their own check-lines, while its POSITIVE needle was one literal that did — so it passed whatever the declaration said. Measured: `var schemaV2SQL = string(schemaSQL)` (the ledger becoming the file it exists to attest) returned **`ok 0.290s`**. The executor's `MUT-SM-V2-LEDGER-DERIVED` redded only because it used the bare form the negative needle was written to catch — **a mutation shaped to the check tests the check, not the threat**. Repaired by anchoring to `^` plus a semantic `schemaV2SQL == schemaSQL` backstop; both forms now RED, and the judge added a third (`schemaV1SQL + ""`) that also REDs. **Read before starting `SM.B2a`:** `AC12`'s *"network confined to `host/broker`"* control has been **VACUOUS since SM.A** (`host/broker` has zero `net/http` deps) and **SM.B2a is the milestone that makes it real** — re-assert it there, never inherit it as green. **Evaluator `sonnet` PASS 91/100, ZERO blocking**, three non-blocking findings all carried: **NB-1** AC-B1.4's pre-bump negative arm is measured EXTERNALLY (against `origin/dev`) rather than embedded in the test, so a future reader cannot see it from the test file alone; **NB-2** under genuine in-process concurrency two callers can both pass the `SELECT EXISTS` pre-check and the loser fails on the `approval_claims` PRIMARY KEY — correct (no double-consumption) but surfacing a wrapped constraint error rather than `ErrApprovalAlreadyConsumed`; the judge MEASURED this (`… UNIQUE constraint failed: approval_claims.approval_ref (1555) …`), so it is evidence rather than inference, and **`AC9b` in SM.B2b is the criterion that must close it**; **NB-3** the doc carried no SM.B1-era verification rows — **discharged this iteration** as row `V-S`. **Also landed this iteration, outside item 8:** SM.A had committed a 15.7 MB `ailang-worldd` Mach-O that five independent checks passed; removed and gated (PR #42 → `e24a6f0`). Prior text follows.] ~~**[IN-SPRINT] — `SM.A` LANDED 2026-08-05 (iter-53), PR #41 → squash `13315da`, CI green BOTH jobs SHA-addressed on the merge commit and step-log verified. The next unit of work is milestone `SM.B1`, gated on nothing.** **What SM.A SETTLED:** `AC6`'s three-arm cross-check AGREES on darwin/arm64 **and** linux/amd64 (tarball `5472 = 5472`; CI printed `✓ compiler pinned by exact bytes: AILANG v0.30.0 on Linux/x86_64`, `9/9` steps), each arm mutation-proven able to red — so the cross-toolchain tarball risk `DD-1` raised is **CLOSED**, and `DD-4`'s third-LEG landed with `4/4` identities and `14` named tests intact. **What SM.A newly OPENED — `DD-7`, and its second half is queue item 9's problem now:** a byte-exact compiler pin is **platform-specific** (darwin `e9746fef…` / linux `1e594d15…`), and separately **CI job 1 has been verifying `.ail` against `AILANG v0.33.0` since `latest` moved on 2026-08-04** — measured in the step log at `af0c3b4` (run `30993399332`) against job 2's `v0.30.0` in the same run. The package leg routes around it with its own pinned install via `WORLD_PKG_AILANG_BIN`, so **nothing is blocked**, but the two ORIGINAL legs remain unpinned and item 9's "latent, not active" grade is now false. **Read before starting `SM.B1`:** it must be ONE commit (`schema.sql` + both `store.go` constants + the stale-version acceptance policy + all three independent fixtures — splitting it lands a red DDL gate); `DD-2`'s ~3× blast radius and `DD-3`/`8/OD-3` below are binding. **`AC12`'s limits, carried forward so SM.B2a does not inherit a false sense of coverage:** the *"network confined to `host/broker`"* green control is **VACUOUS today** (`host/broker` has **zero** `net/http` deps — network arrives WITH SM.B2a, which is precisely when that control becomes real and should be re-asserted), the guard is **source-level not closure-level**, `cmd/ailang-worldd`'s inherited `net/http` is **loopback IPC** and not egress, and `host/registry` in the forbidden list is the *interpreter epoch* registry — a **name collision** that will produce a false positive the moment anything legitimately needs epoch metadata. **Evaluator `sonnet` PASS 87/100, ZERO blocking**, and it added a THIRD carried limit: `AC5`'s smoke coverage is enforced **implicitly** — dropping a module changes the smoke's output — rather than by an explicit import-coverage manifest. Correct outcome, weaker instrumentation; worth strengthening when SM.C revisits the fixtures. Prior text follows.] ~~[IN-SPRINT] — SPRINT-PLANNED 2026-08-05 (iter-52). The next unit of work is milestone `SM.A`, gated on nothing.~~ Plan: `.ailang/state/sprints/w-self-mod-vertical.{plan.json,handoff.md}` (planner `opus`, lane derived VERBATIM as `opus fail-closed:env-pin`; gitignored, as all 18 prior sprint artifacts are). **The planner re-scoped 4 milestones into 6** — `SM.A · SM.B1 · SM.B2a · SM.B2b · SM.C · SM.D` — because SM.B as designed prices at ~2,300–2,700 LOC against a measured maximum single landed commit of **751** insertions (n=5: 457, 414, 210, 751, 698; median 457). It stays **ONE queue item** (precedent: 4d and 4e each ran multi-milestone across 4–5 iterations without becoming a second row); the split is internal. **`AC12` moved from SM.B into SM.A** — a boundary guard that lands alongside the code it constrains has never been observed rejecting that code. **`8/OD-1` IS RATIFIED** by Mark on `#32` @ `2026-08-05T08:25:00Z` (*"Approved publish of world/ in ailang extensions - go. Credential is on your machine for this."*): the POLICY is approved and SM.D is no longer blocked on a human answer — but this is **not** the exact-bytes attended stamp SM.D describes, and it cannot be, because the ready packet does not exist until SM.A builds it. **An authorization is not an attendance.** SM.D stays attended-only, never headless, never in CI, stop-at-readiness by default. **THREE PLANNER FINDINGS, ALL REPRODUCED FIRST-PARTY BY THE CONTROLLER, THAT CHANGE THE WORK.** **`DD-1` (blocking for SM.A as designed):** Decision 3's *"small library extraction of v0.30.0 package hashing logic"* is **impossible** — World is `module github.com/sunholo-data/ailang-world`, upstream is `module github.com/sunholo-data/ailang`, and the hashing lives in `internal/pkg/`, which Go's internal rule forbids across modules. The CLI is no fallback (`pkg_publish.go:110-112` prints `hash[:24]+"..."`; the tarball bytes are never persisted). `AC6` therefore needs a **re-implementation** (`host/pkgproj`) with a mandatory hard-failing 24-char cross-check, plus a newly named risk: the tarball hash rides `compress/gzip` and the two modules declare different Go versions. **`DD-3` (a silent runtime hole the version bump CREATES):** `host/store/store.go` ends its version ladder in a bare `return nil` at `:354`, so once `currentSchemaVersion` goes 1→2 a store still at `user_version = 1` matches **no branch**, opens successfully, and **never executes `schemaSQL`** — `approval_claims` would be found absent at the moment of the irreversible publish. Raised as **`8/OD-3`** and answered from already-ratified text (`4d/OD-3` alt-1: fail LOUD on an unsupported **or un-upgraded** store); non-blocking. **`DD-4`:** `scripts/verify_ail.sh:160`/`:190` are exact equalities (`EXACT_TOTAL_VERIFIED=4`, `EXACT_TOTAL_TESTS = 14`), so the package gate must be a third **LEG**, never a new **root** — adding `packages/` to `ROOTS` doubles the identities to 8 and reds the repo's primary gate for a reason unrelated to the code under test. **THE DDL BLAST RADIUS IS ~3× WHAT THE DOC'S CONFLICT SURFACE NAMES**, and all of it lands in **SM.B1's single commit** (splitting it lands a red gate): the doc names only `host/store/schema_version_test.go`, while `host/store/journal_test.go:714` holds a SECOND independent fixture `canonicalTableDDL` (**7** hardcoded tables behind `requireExactTableNames` at `:778`, carrying its own *"must not be derived from schemaSQL … or the database under test"* comment) that the doc never mentions; `schema_version_test.go:16` `frozenFutureSchemaVersion = 2` **collides** with the new current version; and `store.go:316` writes a **literal** `PRAGMA user_version = 1` that would trip `freshInitTx`'s own drift check at `:325-326`. **FIVE ACs JUDGED VACUOUS AND REPLACED IN THE PLAN** (`AC13`, `AC17`, `AC10`'s 2nd clause, `AC19`, `AC1` — each by the mission's own test: *would this pass identically if the thing it protects did not work?*), plus 5 planner-authored mutations → **36 total**. Prior text follows.] ~~[NEXT] — DOC LANDED 2026-08-05 (iter-51); NOT YET SPRINT-PLANNED. The next unit of work is the sprint-planner run, and it is gated on nothing.~~ Doc: `design_docs/planned/w-self-mod-vertical.md` (839 lines; designer `codex:gpt-5.6-sol`, rotation slot 2; PR #40 → squash `269f1fe`, dev CI green both jobs SHA-addressed on the merge commit). Milestones **SM.A–SM.D**, **4–5 d** (raised from 3–4 by the round-1 quorum revision; SM.B alone is 2–2.5 d). **THE VERIFY-FIRST CLAUSE WAS DISCHARGED FIRST-PARTY BEFORE THE DESIGNER WAS SPAWNED, AND IT REFRAMED THE ITEM.** There is no vendor-namespace claim operation to perform: `cmd/registry-validator/main.go:177` @ the pinned `e37b370` reads verbatim `// Step 5: Namespace auth — deferred (accept all publishers for now)`; auth is ONE optional shared secret whose header the client omits when unset; and a live four-arm dry-run at the pinned binary accepts `world/probe`, `someoneelse/probe` and `sunholo/probe` alike (rc=0 each) against a firing known-positive control (`novendor` → rc=1, `must be vendor/name format`). **`world/` is a string World writes, not a namespace World holds** — recorded here so no later reader re-derives ownership from Mark's wording, which said "verified unclaimed" and was true. Census: 40 packages, vendor histogram `sunholo` 40 and nothing else, `world/` 0 with its control firing — World would be the registry's first second vendor ever. Publish is immutable (409) and unrecallable by the publisher, and **`AILANG_REGISTRY_API_KEY` is AMBIENT in this loop's own tool shells**; compose those and any process inheriting this environment can write irreversibly under any vendor string. That is the clause-3 surface, and it is what makes the brokered/receipted framing load-bearing rather than decorative. **TWO THINGS THE PLANNER PRICES FIRST — both carried, neither cleared.** (i) The design NEEDS a `schema.sql` change (`approval_claims` + an atomic `AppendClaimedEffectIntent`, because `host/store/store.go:601` `SetRegistryHead` is a blind upsert with no expected-previous and is the only set-shaped store API, so there is no claim-if-unused primitive) — and the landed `w-ddl-gate-teeth` DDL gate reds on **any** `schema.sql` edit *by design*, so its fixture update belongs inside the same milestone rather than after it. (ii) Whether 4–5 d is still one queue item or wants splitting at the SM.B boundary. Quorum round 2 ran ONE reviewer (`gpt5-6-sol` absent: `budget`) and round 3 was the controller carve-out, so neither question has had two pairs of eyes. **Open decisions: `8/OD-1`** — the attended human stamp for the irreversible first publish; controller default **do not publish**, and it blocks **SM.D only**, so SM.A–SM.C are routable today. **`8/OD-2`** (upstream namespace authorization) is non-blocking. Prior text follows.] ~~[NEXT] — UNPARKED 2026-08-05 (iter-50). The park condition was "until 4 lands"; item 4 `w-effect-broker-m3` completed at iter-35. Mark re-scoped this row attended on 2026-08-04 (`de80792`). Routes to design-doc-creator; VERIFY-FIRST binding at pick.~~ ~~[PARKED until 4 lands]~~ **w-self-mod-vertical** · clause-7 · one World extension shipped through **[NAMING + FLOW, Mark 2026-08-04: World extensions publish to the PUBLIC AILANG
   registry under the `world/` vendor namespace** (verified unclaimed 2026-08-04; explorer:
   ailang.sunholo.com/docs/packages/explorer). The clause-7 publish is itself a BROKERED,
   RECEIPTED effect (capability + budget + record) — World's self-modifications become public,
   hash-verified artifacts with a full evidence chain. Local-first holds: publish is an outward
   effect-handler interaction; consumption is version-pinned + hash-verified install; the core
   never depends on the registry at runtime. Item-8's design doc VERIFIES publish-auth/vendor-
   registration mechanics at pick (VERIFY-FIRST) — `ailang publish` reads ailang.toml, uploads
   via the validation service.]**
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

   - **⚠ THAT GRADE EXPIRED ON 2026-08-04, AND THE ROW ABOVE PREDICTED IT — RAISED TO ACTIVE
     (iter-53, controller first-party).** The bullet above ends *"the next release turns it active
     with no signal"*, and that is exactly what happened: `gh release view` → `latest` = **v0.33.0**,
     published `2026-08-04T12:25:38Z`. Measured SHA-addressed in the step log at `af0c3b4`
     (run `30993399332`), with its own control in the same run: job `ailang-code verify gate` prints
     **`AILANG v0.33.0`** while job `go host build + test gate` prints **`AILANG v0.30.0`**. The
     rig's PATH `ailang` is likewise **`v0.33.0-1-gdd68e0741`**, not the `v0.30.0-205-…-dirty` this
     row records. So for two iterations (51, 52) the repo's primary `.ail` gate validated against an
     unpinned compiler in violation of CLAUDE.md's hard rule, and nothing surfaced it — **a recorded
     prediction of future breakage is not a monitor; it expires silently unless something
     re-measures it.** Evidence that the difference is not cosmetic: v0.33.0 **fails** the new
     package gate's own step 5 (*"5 properties never ran (no generator)"*).
     **Not blocked, and deliberately not fixed headless.** SM.A's package leg carries its own pinned
     v0.30.0 install (`WORLD_PKG_AILANG_BIN`, a separate path) so the byte-exact gate is correct
     today; the two ORIGINAL legs are untouched, because pinning those is this row's human-gated
     half. **What the human decision now costs, restated with the new fact:** pinning job 1 to the
     v0.30.0 tag makes CI stop tracking upstream releases (World would no longer notice a compiler
     change until it chooses to), while leaving it on `latest` means the primary gate silently
     re-validates against every future compiler. The cheap third option this row already named —
     making `verify_ail.sh` ANNOUNCE its resolved binary + version the way `verify_go.sh:33` does —
     is pure observability, cannot red anything, and would have surfaced this on 2026-08-04 rather
     than a day late. **That half is now the recommended first step.**
   - **Decomposition for the human gate**: the *cheap, zero-risk* half is making `verify_ail.sh`
     ANNOUNCE its resolved binary + version the way `verify_go.sh:33` already does (pure
     observability, cannot red anything). The *human-gated* half is the hard version assertion plus
     the CI `latest`→pinned-tag edit — those two are **coupled** (a hard assert alone would red CI on
     the next upstream release, headless), which is exactly why this row's "confirm before
     implementing / do not hand-edit CI headless" flag was right and was respected this iteration.
10. [**[IN-SPRINT] — `BG.A` LANDED 2026-08-06 (iter-58), PR #47 → squash `278f102`; evaluator
    `sonnet` PASS 89/100 round 1, zero blocking. The next unit of work is milestone `BG.B`, gated on
    nothing — and it carries ONE correction that must be applied before it is routed (below).**
    **`AC2`, `AC3`, `AC4`, `AC5` are DISCHARGED**, mutations `M1`, `M2(a)`, `M2(b)`, `M4` run under
    the house recipe and `M5` run by the controller on all four arms **with its negative control**
    (base harness, same kill → `RESIDUE=YES`, ` M host/store/store.go`; BG.A → 0 changed shas, 0
    porcelain lines, on 4/4 arms, each `alive_at_kill=True`, `killed_rc=-9`). AC2's numbers,
    controller-measured: `host/store` **160/0 → 229/1**, `host/replay` 162/0 → 231/1,
    `cmd/ailang-worldd` 233/0 → 234/1 — the planner's prediction exactly, asserted on the closure
    `checkGoGroup` RETURNED. **CI CAVEAT — ⚠ CORRECTED 2026-08-06 (iteration 59): THIS ROW SAID "RAISED AND THEN SETTLED WITHIN THE ITERATION", AND THE "SETTLED" HALF WAS WRONG. `e3808c0`'s green does prove the CODE is fine — a descendant of the BG.A merge, same job, all 11 steps `success`; that inference stands and is why BG.A was correctly not reverted. It does NOT prove the INCIDENT ended: iteration 58's own doc-only bookkeeping commit `4e959bf` went red on BOTH jobs 13 minutes later, `cancelled` with `steps=0`, and dev is RED at HEAD as of iteration 59. A green obtained during an open incident is a sample, not a settlement — see the iteration-59 STATUS stamp for the per-job/per-step table and six firing controls.** Original text follows, still true of `e3808c0` itself — dev was GREEN at `e3808c0` (BOTH jobs, SHA-addressed, all 11 verify-gate steps `success` in the STEP LOG), and `e3808c0` is a DESCENDANT of the BG.A merge, so the same job passed on the same code: at the merge SHA `go host build +
    test gate` is **SUCCESS** but `ailang-code verify gate` is **FAILURE — in `Set up job`, step 1,
    before checkout, with zero repo commands executed**, on `Failed to resolve action download info.
    Error: Service Unavailable`, during a **declared GitHub Actions incident** (githubstatus.com:
    `Actions: partial_outage`, incident opened `2026-08-06T15:22:49Z`; the run started `15:38:01Z`,
    16 min inside it). Attribution is by MECHANISM plus two firing controls: the **identical tree**
    passed the **identical job** on PR #47 at `15:33:02Z` (`51e18968` → SUCCESS), and the sibling job
    on the merge commit is SUCCESS; both jobs use the same two actions, so it is not action-specific.
    **`BG.B` MUST APPLY THIS FIRST — the plan's write-site count is now wrong by one, and the missing
    one would red the guard `BG.B` installs.** The plan says route BG.A's **two** per-arm writes
    through `confinedWrite`; it was written before the AC4 barrier existed as code, and the barrier
    adds a **third** direct `os.WriteFile(absMarker, …)`. Measured at `278f102`: **3** `os.WriteFile`
    sites (`:383` marker, `:428` mutant, `:439` overlay JSON), **0** `OpenFile`/`Create`/`Rename`,
    KP control `os.ReadFile` = **4**. Routing the marker through `confinedWrite` is CORRECT rather
    than an exemption — the marker is *required* to live outside `repoRoot`, exactly what
    `confinedWrite` permits — so the writer also replaces the bespoke `insideRepo` check at
    `:367–:373`. Raised by the evaluator (`NB-2`), controller-reproduced, plan corrected in place
    with a `controller_corrections` entry. **AND ONE DEFECT IN THE DOC AND PLAN THEMSELVES:**
    `go/parser`'s `readSource` tests `src != nil` on the **interface**, so a typed nil `[]byte` is a
    non-nil interface handed back as an **empty source** — every unreplaced file parses as
    `expected 'package', found 'EOF'`, i.e. *a checker that cannot read the tree finds no forbidden
    imports*. Both artifacts write the helper as `parser.ParseFile(fset, path, <bytes-or-nil>, …)`,
    the exact shape that produces it. Isolated in `parseSrc`; recorded because the written wording
    reproduces it. **Prior head text follows.** ~~**[IN-SPRINT] — SPRINT-PLANNED 2026-08-06
    (iter-57). The next unit of work is milestone `BG.A`, gated on nothing.**~~ Plan:
    `.ailang/state/sprints/w-boundary-gate-tree-mutation.{plan.json,handoff.md}` (107,795 B +
    26,491 B; planner `opus`, lane derived **fail-closed `opus missing-script`** —
    `tools/launchd/derive-planner-lane.sh` is ABSENT from this checkout, control fired
    (`tools/launchd/` holds `mission-control.sh` + `mission-template.plist`), and
    `MISSION_PLANNER_MODEL=opus` independently agrees; gitignored, as all 19 prior sprint artifacts
    are). **THREE MILESTONES, PARTITION COMPLETE:** `BG.A` (AC2, AC3, AC4, AC5 · M1, M2, M4, M5) →
    `BG.B` (AC1a · M3, M6) → `BG.C` (AC1b, AC6′ · M7) — 7 criteria, 7 mutations, none dropped or
    double-assigned. Ordering has ONE forced joint (the AST guard reds on sight while
    `mutateAndRestore` still calls `os.WriteFile` at `:205/:209/:224`, so `BG.B` cannot precede
    `BG.A`) and one chosen (`BG.C` last — it is the only milestone whose green depends on an
    UNMEASURED CI-filesystem property, so a red there leaves the repair already landed).
    **THE PLANNER JUDGED FOUR OF THE DOC'S SIX ACs VACUITY-CAPABLE AND REWROTE THEM** — precedent:
    iter-52's planner found 5 in a doc that had passed quorum. The doc now carries a
    ⚠ SPRINT-PLANNED block recording every supersession, so the plan and the doc cannot rot apart
    (rule 3b(vii), applied forward rather than after the fact). **`AC6` IS THE HEADLINE, AND THE
    CONTROLLER FOUND IT AT BASELINE BEFORE THE PLANNER EVER RAN** (rule 3e): its `≤2× 0.435s` bound
    is a constant transcribed from another worktree at another cache warmth. On UNCHANGED code at
    `e9c8c85`: fresh-worktree first run **0.664 s / 0.621 s** (n=2, two independent worktrees), warm
    steady state **~0.480 s** (n=9, 0.472–0.507), tree 0-dirty throughout. Zero code change already
    sits at **1.43–1.53×** the AC's own constant cold — CI checks out fresh, so cold is the operative
    state — leaving ~**1.31×** of headroom, not 2×. **The false-red risk is the lesser harm; the
    greater one is that a GREEN `AC6` could not have failed informatively**, because the noise band
    consumes ~76% of the budget. A threshold whose noise is the size of its signal is this item's own
    spine turned on the item's own acceptance criterion. The planner then found **two further
    defects the controller missed**: the units are ambiguous by **1.32×** (go-*reported* 0.479 s vs
    wall-clock 0.631 s median, n=5 — 0.435 s is a go-reported figure while the wording is about the
    command *completing*), and what `AC6` nominally protects is a **600 s** `-race` budget against a
    **0.5 s** package — **1200×** of headroom, so it could not fail for the change either. `AC6′` is
    a paired same-session ratio (one worktree, equal warmth, ≥8 interleaved A/B pairs swapping only
    this file and asserting its sha256 changed, `median(wall_B)/median(wall_A) ≤ 1.50`) plus a
    `median(wall_B) ≤ 3.0 s` ceiling, **asserted on wall-clock and said so explicitly** so the unit
    cannot be re-ambiguated — with its noise floor MEASURED rather than assumed (8 interleaved pairs
    on IDENTICAL code → **1.0079** against a true 1.000, pooled spread 1.058; 1.50 is ~8.5× that).
    **TWO PLANNER REFUTATIONS OF THE DOC, BOTH REPRODUCED FIRST-PARTY BEFORE BEING RECORDED** (Gate-2
    rule (b)/(d) — a judge's finding is a claim too). **(1) `V16` IS REFUTED (→ `V16a`/`V16b`):**
    `cmd/ailang-worldd`'s closure DOES contain a forbidden prefix — `host/registry`, which is
    `forbiddenImportPrefixes[3]` (`:53`) — reached via `cmd/ailang-worldd → host/daemon →
    host/registry` (`daemon.go:51`, direct). Measured worldd **1** of 233, store **0** of 160, replay
    **0** of 162, KP `modernc.org/sqlite` firing **1/1/1**. So the doc's "0 forbidden-prefix hits in
    all three closures" and "a scan would be green on the current tree" are FALSE. **But the red is a
    FALSE POSITIVE, so the `10/OD-1` deferral gets STRONGER, not weaker:** `host/registry` is the
    *interpreter epoch* registry (`world/epoch-registry/v1`), not the *package* registry the entry
    targets — **exactly the name collision iter-53 predicted in prose and iter-57 measured**. A
    closure scan today reds legitimate code, so `10/OD-1` cannot be implemented until the collision is
    resolved. The doc's error hid because ONE table cell bundled two questions and only one of them
    ever had a firing control. **(2) NOTHING A PASSING GO TEST EMITS REACHES CI (`V16c`):** paired
    arms on `TestWorldBoundaryDependencyAllowlist` — CI's exact form (no `-v`) prints the single line
    `ok … 0.580 s` with **0** matching lines, against a KP arm (identical + `-v`) at **12**;
    `verify_go.sh:100` is `go test ./... -count=1` and its `-race` leg builds
    `["go","test","./...","-count=1","-race","-timeout","8m"]` — **neither carries `-v`**. So the
    gate's `ENUMERATION`/`MUTATION`/`RESTORE` diagnostics have **never appeared in a CI log**, and
    "loud but non-gating" is a contradiction in CI: any observable this design wants CI to see must
    be an **assertion**, never a log line. **A THIRD VACUITY MODE `Decision 8` NEVER CONSIDERED:**
    `host/boundary` holds exactly ONE `.go` file and it is a `_test.go`, so an empty `ParseDir`, a
    filter dropping `_test.go`, or a selector bug each give the AST guard "zero violations, green" —
    the doc defends the *self-match* mode at length and is silent on the *empty-enumeration* mode.
    Also measured: `go list -overlay` **silently ignores** a `Replace` key matching no file (rc=0,
    base closure, no stderr), so `AC2`'s assertion needed strengthening to the closure `checkGoGroup`
    actually CONSUMED; and `git status --porcelain` reports UNTRACKED files, so `AC4`'s in-tree ready
    marker would have redded on the harness's own artifact. **CARRIED, NOT CLOSED — `C1` IS
    APFS-ONLY:** the `ModTime` backstop was measured 200/200 on darwin/APFS and CI job 2 is
    `ubuntu-latest`; Linux takes mtime from a tick-granularity coarse clock (1–4 ms), so it may not
    transfer, and the planner could not measure it (no docker/colima/podman/lima on the rig). Hence
    `BG.C` last, a fail-loud 20/20 granularity probe whose failure is a TEST FAILURE naming both
    `st_dev`s, sha256+size+mode+**inode** asserted unconditionally (inode closes a rename route both
    of the doc's stated observables miss), and a pre-authorized fallback: record the refutation, keep
    the four filesystem-independent observables, open `10/OD-3` — and **never lower the 20/20
    threshold**. **Estimate: ≤1 day of EFFORT holds** (velocity: 12 landed feat/fix commits, median
    **363** insertions; closest analogue `1761a9c`, a single-test-file change to *this same file*,
    +56/−7, one iteration) **but 2–3 iterations ELAPSED** — measured cadence is ≤1 milestone per
    iteration and **4 of the 7 mutations cannot run in the executor sandbox** (M5 needs subprocess
    SIGKILL + git inspection; M6/M7 re-arm live-tree writes; AC6′ needs a file swap out of git
    history), so the controller pass is the critical path. LOC re-estimated **+150 → +250**, all of
    it spent making the doc's own ACs non-vacuous. **Prior head text follows.** ~~**[NEXT] —
    PROMOTED AHEAD OF `SM.B2a` AND DESIGN DOC LANDED 2026-08-06 (iter-56), PR #46 →
    squash `ca25ed6`, dev CI green BOTH jobs SHA-addressed on the merge commit. The next unit of
    work is the sprint-planner run, gated on nothing.**~~ Doc:
    `design_docs/planned/w-boundary-gate-tree-mutation.md` (478 lines; designer
    `claude:claude-fable-5`, rotation slot 1). **PROMOTED ON A MEASUREMENT, NOT ON THE QUEUE
    ORDER** — iter-55 wrote "consider ordering BEFORE `SM.B2a`" and left the call open; iter-56
    discharged it by measuring a hazard iter-55's row does NOT record. The restore is a `defer`,
    and a `defer` survives return and `t.Fatal` **and nothing else**: SIGKILL mid-mutation leaves
    the gate's own mutant on disk **permanently** (`rc=137`, ` M host/store/store.go`), and Go's
    OWN `-test.timeout` panic does the same (`250ms` → `rc=2`, residue on
    `cmd/ailang-worldd/main.go`) against a firing 60s green control. **Both kill paths are inside
    the repo's own gate**: `scripts/verify_go.sh` runs `go test ./... -race -timeout 8m` wrapped in
    an `os.killpg(…, SIGKILL)` at 600s (`:113`). **THE RESIDUE IS INVISIBLE TO THE BUILD** — with
    all three mutants applied simultaneously, `go build ./...` **rc=0** and `go vet` **rc=0**;
    what it breaks is **the boundary gate itself**, which then accuses an innocent production file
    of a network-boundary violation with a message indistinguishable from a real one — during the
    one milestone (`SM.B2a`) whose job is to add network code. **CORRECTION TO THIS ROW'S OWN
    TEXT, MEASURED**: the mutants are NOT "deliberately non-compiling" (this row, the iter-55
    stamp and the dashboard all said so, three times). They compile — the gate's own inline
    comment says "compiling HTTP import" — so iter-55's transient `could not import` CI failure
    cannot be explained by "the file is broken"; it is a build-graph artifact of the file
    *changing* under a concurrent build. **DECISION**: `go list -overlay` + an overlay-aware read,
    so the gate never writes a byte under the repo root — closing crash-residue and the
    concurrency window with ONE mechanism. Controller-verified first-party rather than inherited:
    the overlay closure is **diff-IDENTICAL** to a physically poisoned tree (control fired, 69
    packages difference) and the tree stays **0-dirty** throughout; the scratch-copy alternative
    was rejected because `repoRoot` bakes the real root in at compile time. **TWO QUORUM ROUNDS,
    BOTH BLOCKED, BOTH ACCEPTED RATHER THAN ARGUED, AND BOTH REMOVED A WAY FOR A GREEN TO MEAN
    "THE CHECK NEVER RAN."** Round 1 (both reviewers, one objection): the AC1/M3 polling detector
    was probabilistic — deleted, not supplemented, and replaced by `gpt5-6-sol`'s structural write
    confinement + AST guard plus `gemini-3-1-pro`'s deterministic `ModTime` backstop, whose premise
    the controller MEASURED rather than forwarded (back-to-back write+restore advanced
    `st_mtime_ns` in **200/200**, KP controls firing both ways). Round 2 (two NEW objections,
    resolved under the narrow-refinement carve-out — neither disputed the design DIRECTION):
    `AC4/M5`'s randomized kills could pass without ever arming the threat → replaced by
    `gpt5-6-sol`'s verbatim barrier protocol; and P9's audit had grepped `os.WriteFile` only while
    the design's own threat model names `OpenFile`/`Create`/`Rename` → **gemini was right**, and P9
    is now a MEASUREMENT (whole-suite `git status` sampling at 20 Hz: **1** dirty observation with
    the gate enabled, **0** with only its mutation test skipped) rather than a search. **Open:**
    `10/OD-1` (scan the dependency CLOSURE, not just direct imports — see below; default out of
    scope) and `10/OD-2` (kill-harness permanent in CI vs sprint evidence; default sprint
    evidence). Both non-blocking, controller defaults recorded. **`10/OD-1` IS A REAL SECOND
    FINDING, controller-confirmed**: `checkGoGroup` (`:130`) computes the `go list -deps` closure
    and then uses it **only** as a `len(deps)==0` anti-vacuity guard — the forbidden-prefix scan
    walks direct import specs alone, so a TRANSITIVE forbidden dependency is invisible in a gate
    named "dependency allowlist". Prior text follows.] ~~[BACKLOG — **MEASURED FIRST-PARTY ON
    PRISTINE `dev` iter-55**; consider ordering BEFORE `SM.B2a`]~~
    **w-boundary-gate-tree-mutation** · clause-1-infra · **the World-boundary gate (AC12,
    `host/boundary/allowlist_world_test.go`, landed with SM.A) proves its teeth by rewriting THREE
    OTHER PACKAGES' PRODUCTION SOURCE FILES in the live working tree — and `go test ./...` builds
    those packages concurrently.** `mutateAndRestore(t, root, group.mutantFile, …)` (`:271`) writes
    `cmd/ailang-worldd/main.go`, `host/store/store.go` and `host/replay/replay.go`, each with a
    **COMPILING** `net/http/httputil` import (corrected iter-57 — this sentence said "deliberately
    non-compiling" and iter-56 measured that FALSE; see the correction ~35 lines above, which iter-56
    added without deleting the sentence it corrects), then restores. The window is small but
    real and it is PUBLISHED TO EVERY CONCURRENT READER.
    - **Measured, pristine `origin/dev` @ `ff5d5cc`, tree `0` modified, controller first-party.**
      Sampling the three files during ONE gate run: `main.go` observed mutated **5** times,
      `store.go` **5**, `replay.go` **5**, out of **90** samples. **Control: 0 / 200** with the gate
      not running. So the broken state is genuinely observable, and the control proves the sampler
      is not simply always-true.
    - **It has already fired.** `go test ./...` in the main checkout red-lit
      `TestCLIRealSubprocessEpisode` with `cmd/ailang-worldd/main.go:27:4: could not import
      net/http/httputil` — `cli_test.go:128` builds the subprocess binary from source and sampled
      the file inside the gate's own window. **Attribution is by MECHANISM, not co-occurrence**
      (rule 3d): `httputil` is exactly this gate's `mutantImport` for that group, and nothing else
      in the repo inserts that string into that file.
    - **Rate is low and honestly bounded**: **0 / 8** runs of `go test ./... -count=1` on pristine
      `dev` reproduced it, so this is **latent**, not active — CI has been green. It is recorded
      because `SM.B2a` adds substantially to the already-slow `host/broker` suite (76 s under
      `-race` locally), which lengthens the overlap window rather than shortening it.
    - **Not caused by iter-55's AC12 repair** — every measurement above was taken on a pristine tree
      with that change ABSENT.
    - **The fix is a design choice, not a one-liner**, which is why this is a queue row and not a
      controller fix-forward: the principled repair is to mutate a scratch COPY of the tree rather
      than the live one (`go list -deps` then resolves in the copy; `GOMODCACHE` is global so the
      copy is cheap), which changes what the gate proves and wants a designer + evaluator. The blunt
      alternatives — serialising the suite with `-p 1`, or an advisory lock the other packages'
      tests would also have to take — are worse and touch more files. · ~0.5–1d · surfaced iter-55.

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
| **The daemon enforces global HTTP write deadlines** (frozen D7 block) | **iter-24, controller first-party** | `host/daemon/daemon.go:409-414` wires constants declared at `:77-91`: `ReadHeaderTimeout` 5s, `ReadTimeout` 30s, `WriteTimeout` 30s, `IdleTimeout` 120s. Consequence: a long-lived SSE stream mounted on this server is killed at the 30 s write deadline unless that ONE route relaxes it via `http.ResponseController` — and that API is available. `IdleTimeout` is NOT a stream-lifetime bound (it governs idleness between requests). **CITATION CORRECTED iter-52 (conclusion UNCHANGED):** this row asserted `go.mod:3` = `go 1.26.4`; measured 2026-08-05 it is **`go 1.25.6`**, lowered deliberately at `77ce069` (RG.A M1, *"lower the Go floor to 1.25.6 and land the canary that reds without it"*). Every historical value was enumerated by command rather than assumed — `d5b155c`/`a901c30` → `1.26.4`, `77ce069`/`f19acac` → `1.25.6`. The conclusion survives because `http.ResponseController` has existed since Go **1.20**, so nothing decided on this row is affected; it is the citation that rotted, exactly as `P42`/`P47b` did inside a design doc. **The uncomfortable generalisation: the oldest-declared-base freshness sweep this mission invented protects DESIGN DOCS, and nothing sweeps THIS table**, which is older than every doc it has ever routed |
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

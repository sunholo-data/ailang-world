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
- **`passed_tests` is NOT the named-test count (process fix, iter-25 — the landed correction D-B,
  restated because a fresh author hit it again).** `ailang test --format json` reports
  `passed_tests`/`total_tests` that ALSO count contract-derived properties: for
  `storejournal.ail` that is `passed_tests: 32` against **`len(tests[]) == 25`** named tests. The
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

---

## STATUS (rotation rule)

Newest **3** STATUS stamps live here; older ones move to `world-mission-status-archive.md`.
At Gate 4, after adding your stamp, move the now-4th stamp to the TOP of the archive file.

## STATUS 2026-07-28 (iteration 25) — **`w-store-durability` (clause-1 kernel durability) NEW-DOC authored + TWO QUORUM ROUNDS + carve-out revision → DOC LANDED, item re-tagged `PARKED needs-human-review` on a RATIFICATION PACKET. The design is finished; the block is Mark's, not a defect. And the defect turned out to be eight fields wide, not one.** Item 4 `w-effect-broker-m3` stayed PARKED (still no `@MarkEdmondson1234` answer, now 2 iterations old; the recorded default `(b)` was again NOT force-applied — the queue was not fully blocked), so the pick was the queue's next actionable item, **item 4b `w-store-durability`**, exactly as iter-24's Next line predicted. NEW-DOC tag verified as a FACT before spending anything: `grep -ril w-store-durability design_docs/` matched only the charter + log, `planned/` held three unrelated docs, zero merged PRs, zero `origin/dev` commits. **THE ITEM IS BIGGER THAN THE ROW THAT QUEUED IT — measured, not argued.** The charter row named ONE field (`store.Commit` writes a zero `PrevEntryHash` that `store.GetLogEntry` cannot read back). My own throwaway table-driven probe against `origin/dev` `615619c`, run BEFORE any routing and again after the designer's account rather than inheriting it, found **`store.Commit` validates NONE of its eight ref fields**: all eight zero-value commits return `err=<nil>`; **seven** persist a permanently unreadable row (`TransitionFn` · `Interpreter` · `EntryHash` · `TransitionRef` · `PrevEntryHash` · `NextWorld.LogHead` · `NextWorld.StateRoot`), each failing at read with its own structured `hashref: empty hashref text`; and the **eighth**, `NextWorld.Ref`, commits and reads back *fine* — an empty-string ref becomes the selected head, degenerate-but-readable, which makes it the one poison shape a read-side fix could never even observe, and thereby the single sharpest argument for ARM V1 over ARM V2. Two further measured facts the charter row did not carry: the poisoned commit **ADVANCES the head** (so the store's current world addresses an entry no reader can load), and **a subsequent, entirely legal commit chains onto the poisoned entry and reads back fine** — the append-only log grows a **permanent hole mid-chain with readable entries on both sides**, with no detection and no recovery path. Blast radius reaches the REST surface: `handleLogRange` loops over `GetLogEntry`, so ONE poisoned index 500s the **whole range read**. **Designer = Fable** (`claude:claude-fable-5`; rotation last-used was `codex:gpt-5.6-sol`, so claude was next — and it also kept the author Anthropic while both quorum reviewers are non-Anthropic, so unlike iter-24 there was **no generator=judge collision** this round). Author run 15.7 min under a 30-min cap, revision 7.4 min under a 25-min cap, both rc=0. Delivered `design_docs/planned/w-store-durability.md` (**1,058 lines**) + `design_docs/sketches/storejournal.ail` (**163 lines**) — and unlike iter-24's P8 argument, this doc SHIPS `.ail`: the receipt/retry/write-validity contract is pure kernel law, so coding-standard S1 applies. **7/7 contracts Z3-`verified`** (`writableRefText`, `isIndeterminate`, `mayReportNotStarted`, `retryAllowed`, `journalSeqNext`, `outcomeMatchesIntent`, and the round-1 addition `intentBindsCommit`), 0 counterexamples, 0 errors — every one re-run first-party by me on the pinned binary, never taken from the designer's report. **THE QUORUM EARNED ITS KEEP FOUR TIMES ACROSS TWO ROUNDS, and one of its catches was against ME.** Round 1 (both present, `$0.133037`) **BLOCKED**: `gpt5-6-sol` — `Commit.InvocationID` was never BOUND to the intent's planned fields, so `AppendIntent(id, A)` then `Commit(id, B)` would write a receipt claiming A resolved as B, defeating the journal's whole truthfulness claim; `gemini-3-1-pro` — one `ScanUnreadable(fromIndex, limit)` cannot paginate BOTH `log_entries` (integer index) and `worlds` (content-addressed TEXT key) without unstable `OFFSET`s or implicit rowids a `VACUUM` can renumber. Both applied in the reviewers' own terms via one designer revision — including a **new Z3-proven sketch LAW** (`intentBindsCommit`) so field-level binding is PROVEN rather than asserted in prose, plus the verbatim `ScanUnreadableLog`/`ScanUnreadableWorlds` split with lexicographic `afterRef` ordering stated explicitly. Round 2 (both present, `$0.157032`) **BLOCKED again**: `gpt5-6-sol` — the bounded-allocation claim was unsupported, because *"merely stopping at the supplied limit does not bound allocation or query work"* when the limit is caller-controlled with undefined zero/negative/oversized behaviour, and the startup sweep as written could silently miss every hole after page 1; `gemini-3-1-pro` — the round-1 binding compared only three fields, but **premise V12 records that the kernel never derives `EntryHash` from entry contents**, so a caller could mutate `PrevEntryHash`/`TransitionFn`/`Interpreter` while holding `EntryHash` byte-identical and the receipt would falsely claim the original intent succeeded. **NARROW-REFINEMENT CARVE-OUT APPLIED** (bounded 2nd revision, controller, reviewers' VERBATIM text): both limbs hold — each objection carries concrete reviewer-authored `proposed_fix` prose, and **neither disputes the design DIRECTION** (arms, milestones, ratification framing and the (a)/(b) gating section all stand). Applied: kernel constants `MaxPendingIntentsPage`/`MaxIntegrityScanPage` with a `1 <= limit <= Max…` guard returning `InvalidLimitError`, a startup scan that pages to completion or a fixed budget and emits a **distinct `integrity_scan_incomplete` warning carrying its continuation cursor rather than a clean-looking message over a truncated scan**, the commit-defining field list widened to all seven ref fields with the `EntryHash`-preserving drift case REQUIRED in AC15, P7 rewritten around "the KERNEL owns the bound", two Design-Freeze boxes, and four new RED mutations (`MUT-CALLER-OWNS-LIMIT`, `MUT-SCAN-SILENT-TRUNCATE`, `MUT-INTENT-NARROW-BIND`, joining `MUT-INTENT-UNBOUND`). No third round — the carve-out is one bounded revision, not a re-litigation. **THE REVIEWER CAUGHT AN ARITHMETIC ERROR IN MY OWN FIRST-PARTY EVIDENCE, AND IT IS CORRECTED IN THE OPEN.** My premise row V23 said "seven-field matrix … the seventh, `NextWorld.Ref`" while listing seven *poisoned* fields before it — eight, not seven. `gemini-3-1-pro` caught it as its round-2 "catch". The durable lesson is sharper than the typo: **first-party measurement earns trust for the OBSERVATIONS, not for the arithmetic I wrap around them** — the probe transcript was right, the sentence summarising it was not, and a sub-agent claim gets re-run while my own summary got waved through. Everything a controller hands downstream is a claim too, including its own count of its own results. **WHY NO SPRINT RAN, AND WHY THAT IS NOT A MISS.** The carve-out normally routes straight to sprint-planner. Here the doc's own last acceptance box refuses: **the kernel arms are ratification-class under the charter's frozen-kernel guardrail** (a behavioural change to landed `store.Commit`, the first-ever `schema.sql` change, and a `Commit` struct extension), so a plan would be a plan for work that cannot legally start, and if Mark picks V2 or J2 the plan is waste. Same shape as iter-24's decision, different cause: iter-24 was blocked on external prerequisites, this is blocked on one human decision that the doc has already reduced to three named arm choices with recommendations. Recorded as a deliberate routing decision, not an omission. **A GATE WEAKNESS FOUND INDEPENDENTLY BY BOTH THE CONTROLLER AND THE DESIGNER**: `verify_ail.sh` Leg 2 runs `ailang test --format json world/` — `world/` ONLY — so a sketch's inline `tests[]` are **never CI-executed**; the sketch's contracts ARE swept by Leg 1's per-module `ai-check`, but its tests are honest dead weight in the gate unless a milestone runs them explicitly. Recorded in the Repo Profile with the module-count move (**9 → 10**, so a future iteration seeing 10 is observing this commit, not a regression; the load-bearing 4/4 identities + 14/14 named `world/` tests are `world/`-scoped and UNCHANGED). **A SECOND MEASUREMENT-HONESTY CORRECTION, MINE THIS TIME**: the designer wrote "26/26 named tests" taken from `passed_tests`, when `len(tests[])` was **20** — `passed_tests` also counts contract-derived properties, exactly the landed correction D-B. I caught it by re-running the command instead of reading the report, corrected it in seven places before round 1, and it is now a Repo Profile rule. Post-revision the true numbers are `len(tests[]) = 25` named, `passed_tests = 32` (25 named + 7 properties) — always reported separately from here on. **Controller's own independent evidence** (never laundering a sub-agent claim): the eight-field matrix above; `ai-check` on the final sketch → `check.passed: true`, **7 verified / 0 counterexample / 0 errors**, all seven enumerated in `verify.results[]` (so z3 genuinely ran — not the V27 silent-skip); `ailang test --format json` → 25 named / 32 passed / **0 failed**; `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → **PASS**, 4/4 required `world/` identities across **10** modules, 14/14 named `world/` tests, re-run in BOTH the worktree and the main tree; and scope clean **by diff**, exactly two new files plus this charter and the log — no Go, no schema, no CI, no `go.mod`. **ROUTING**: designer `claude:claude-fable-5` ×2 (author + round-1 revision) on the subscription lane via `claude-sub` (billing tripwire CLEAN, probe replied `ok` before either run); controller applied the round-2 carve-out revision inline; **no planner / executor / evaluator ran** — the item parked before the sprint lane. **`metered=$0.290069`** (quorum reviewers only: $0.133037 + $0.157032; designer + controller on subscription buckets; $5 ceiling untouched). Preflight clean: armed, billing **CLEAN**, `sunholo-voight-kampff`, dev==origin/dev (`615619c`), workflow `CI` **completed/success** at HEAD (this repo has exactly ONE workflow — Build-and-Release / Docs-Deploy do not exist here, so they are **N/A, not pending**), no `[nightly-eval]` issues, only `#9` open, no new `@MarkEdmondson1234` comment on `#9` (13 comments, every one a bot report) nor on predecessor `#1` (25), watermark `2026-07-27T08:55:11Z` unchanged, no rotation due (`#9` titles this week; 13 ≪ 80 — the iter-20 intent test). Inbox: two unread, both V1-side noise for World (an eval-suite start notification and mission-v1's iter-109 report — a sibling report is neither directive nor request), triaged to zero. **ONE PROCESS FIX, no skill edit** (World never edits the shared skill): the iter-24 lesson held under test — I started to arm a `pgrep -f "claude-fable-5"` completion poll, recognised it would **self-match its own shell** exactly as iter-24's did, killed it before it ran, and used the captured-pid/`kill -0` form plus the launcher's own notification instead. The scar tissue worked; what is NEW is the mirror image recorded above — an external reviewer catching an arithmetic slip in the controller's own evidence summary. **Next**: the queue has **no unblocked actionable item left** — 4 and 4b both await Mark, 5 awaits `ailang#498` + clause 3, 6/7/8 are gated behind them. The single highest-value human action is ONE comment on `#9` answering item 4's `(a)`/`(b)` **and** item 4b's three arms together; that unparks two items straight to sprint-planner. If the park persists into the next fire, the queue's only remaining self-serve work is item 9 `w-verify-binary-lockfile` (~0.5d infra) — which the charter itself flags as needing a human confirmation before implementing, so the honest report is that this loop is now human-gated.

## STATUS 2026-07-28 (iteration 24) — **`w-mcp-projection` DESIGN DOC WRITTEN + TWO QUORUM ROUNDS + carve-out revision → LANDED as a record, and the item RE-TAGGED `BLOCKED` on three named prerequisites. The clause-6 boundary cannot be built on v0.30.0, and we now know exactly why.** Item 4 `w-effect-broker-m3` stayed PARKED (no `@MarkEdmondson1234` answer yet), so the pick was the queue's next actionable item — **item 5 `w-mcp-projection`**, the NEW-DOC lane. Grep confirmed the NEW-DOC tag was a FACT (no doc existed; `design_docs/planned/` held only `w-effect-broker-m3` + `w-log-epoch-decision`). Designer = **`codex:gpt-5.6-sol`** (rotation: last-used was `claude:claude-fable-5`, so codex was next) → `design_docs/planned/w-mcp-projection.md`, **641 lines**. **NO SPRINT RAN, DELIBERATELY — and that is the finding, not a failure.** **THE ITEM'S OWN PREMISE DID NOT SURVIVE CONTACT WITH THE BINARY.** The queue row read "reuse `ailang serve-api --mcp/--a2a` machinery — do not reinvent · ~1d". The charter's Conflict Surface had already, correctly, split that into premises (protocol projection of static `.ail` exports, live-tested 2026-07-23) versus **acceptance criteria** (dynamic worldd-backed registry, per-session capability filtering, propose→verify→commit enforcement) — and this iteration proved the acceptance half is **not reachable on v0.30.0 at all**. **THE LOAD-BEARING DISCOVERY: the projected surface cannot be an exact allowlist.** My own first-party stdio MCP probe against the pinned binary (not the designer's account — I re-ran it): `unfiltered → ['addOne','submit_feedback']`, `--routes-only → ['submit_feedback']`, `--caps '' → ['addOne','submit_feedback']`. So (i) a **built-in `submit_feedback` tool is exposed under EVERY flag combination and no flag removes it** — and its own tool description, which I captured verbatim, routes to a `public-feedback` inbox with a **Pub/Sub notification**, i.e. a built-in *egress* tool no World session authorized, sitting inside a bar whose clause 2 demands zero cloud dependencies in the core and whose clause 3 demands that **no ambient-authority path exist from an agent to the outside world**; (ii) discovery is built from loaded module exports rather than a caller-supplied provider, and `--caps` / `--routes-only` / `--api-key-*` are all **process-wide**, so a static key and a process cap set cannot represent a session. My earlier finding that `--a2a` publishes the 8 embedded `std/io` exports (`writeBytes`, `exit`, `readLine`, …) as callable A2A **skills** turned out to be the *less* severe half: `--routes-only` does suppress those, which also means the older upstream `#145` is genuinely FIXED and this is not its regression. Two further first-party facts for whoever builds this: MCP HTTP at `/mcp/` replies **SSE-framed** (`event: message` + `data:`), not a plain JSON body; and module resolution is **cwd-sensitive** — serving `design_docs/sketches/` from the repo root fails `LDR001: module not found: sketches/worldtypes` while serving `sketches/` from `design_docs/` succeeds, so the charter's 2026-07-23 live-test premise row (which records no cwd) is insufficient as launch configuration. **CONSEQUENCE: reuse paths (a) and (b) are REJECTED ON EVIDENCE, and the design takes path (c)** — request a narrow public serving seam from the existing `internal/apiserver` (mount MCP-HTTP/A2A handlers on a caller-owned mux; resolve principal/session before discovery *or* invocation; caller supplies the exact visible descriptors; named invocations route back with the same session; MCP tools and A2A skills generated from that one set; **no built-in tool unless the caller supplies it**; upstream keeps codec ownership and SSE conformance). Path (b) — a sidecar per session — was rejected because it makes process lifetime the session model and still exposes `submit_feedback`; reverse-proxy filtering was rejected because it would make World parse and re-encode MCP/A2A, the exact reinvention DESIGN.md §3.7 forbids. **ROUTED UPSTREAM ON BOTH CHANNELS: `sunholo-data/ailang#498`** + `msg_20260728_015014_8e5a281e`, with the full stdio repro, the version pin, `#145` cited as the fixed predecessor, a narrower interim ask (a flag that suppresses `submit_feedback` + document that `--caps` gates execution not discovery), and the **cause stated as a labelled HYPOTHESIS** since upstream source was never inspected. **THE QUORUM EARNED ITS KEEP TWICE AGAIN.** Round 1 (both present, `metered=$0.0634`) **BLOCKED**: `gpt5-6-sol` — no bounded-wait contract anywhere on the projection path (no timeout source, maximum, context-propagation rule, cleanup rule, error mapping, acceptance test or mutation for a stalled resolver/registry/broker/verifier/client; "broker unavailability returns an error" is not a contract); `gemini-3-1-pro` — mounting SSE-framed MCP handlers on the worldd daemon ignores that a REST daemon's `ReadTimeout`/`WriteTimeout` **abruptly kill long-lived SSE streams**. **I ran gemini's own "verify" step, and it collapsed its two-branch fix to ONE branch**: gemini offered "use `http.ResponseController` **or** document that the daemon lacks global timeouts", and the second branch is **FALSE here** — VERIFIED BY ME at `host/daemon/daemon.go:409-414` wiring constants at `:77-91`: `ReadHeaderTimeout` 5s / `ReadTimeout` 30s / `WriteTimeout` 30s / `IdleTimeout` 120s, **the D7 bound-constant block that `w-worldd-m2` A2 ratified and FROZE**. So the revision was directed to the `ResponseController` branch scoped to `/mcp/` only, with the D7 constants and every REST `/v1/*` path byte-unchanged, plus the follow-up the freeze demands but the reviewer did not ask: *what bounds the connection once its write deadline is relaxed?* → a second explicit finite stream-lifetime maximum, with `IdleTimeout` correctly excluded (it governs idleness **between** requests, not an active handler). Round 2 (both present, `metered=$0.0729`): **`gemini-3-1-pro` PASSED**; `gpt5-6-sol` rejected once more, and its objection is the best single catch of the night — Decision 6/AC13 as revised promised that cancellation after commit *begins* still yields no store/log mutation, which is **not achievable**: an HTTP context can expire mid-atomic-commit, so absent a verified atomic "not-started versus committed" contract the commit may succeed while the caller observed cancellation. **NARROW-REFINEMENT CARVE-OUT APPLIED** (bounded 2nd revision, controller, reviewers' VERBATIM text): both limbs hold — each objection carried concrete reviewer-authored `proposed_fix` prose, and **neither disputed the design DIRECTION**; the defect is truthfulness-of-claim, an acceptance criterion asserting a guarantee the system cannot provide. Applied: the commit-boundary paragraph replaced with the reviewer's own contract **as an auditable block quote with the over-strong claim it replaces stated explicitly rather than silently deleted**; AC13 now tests cancellation on **both sides** of the boundary (before acceptance → no durable mutation; after → exactly one recoverable receipt under a stable invocation/idempotency ID; never a definitive "not committed" while the outcome is unknown) and carries gemini's **OS-level socket-closure** assertion (`ConnState`/client-read-error, because logical `context` cancellation is not proof a socket closed); a `Commit-boundary contract` premise row marked **UNVERIFIED — PREREQUISITE** rather than inventing a mechanism; and three new RED mutations `MUT-COMMIT-BOUNDARY-LIE`, `MUT-LEAK-SSE-CONN`, `MUT-DROP-DEADLINE` (the reviewer's own name). **NO third round was run** — the carve-out is one bounded revision, not a re-litigation. **THE CARVE-OUT DID NOT ROUTE TO sprint-planner, and the reason is the doc's own conclusion, not the quorum**: P6.B is blocked on **three** prerequisites — the upstream seam (`#498`), the clause-3 transition registry + broker session API, and now the verified commit-boundary contract — so a sprint plan would plan work that cannot start. **VERIFIED BY ME, not inherited**: a repo-wide search for `[Tt]ransition[ -]?[Rr]egistry` matches **only** `design_docs/` — **zero** hits in `host/`, `world/`, `cmd/` — so the clause-3 prerequisite is real, not defensive; `host/registry` is the *interpreter epoch* registry (`world/epoch-registry/v1`), a different thing. **A GENERATOR=JUDGE COLLISION, FLAGGED IN BOTH ROUNDS — and it produced evidence rather than just a caveat.** The doc was authored by `codex:gpt-5.6-sol`, so the `gpt5-6-sol` reviewer seat was a **SELF-review**. It was retained rather than excluded on the reasoning that reject-by-default synthesis means a self-*pass* cannot manufacture a PROCEED, so the seat can only add objections — and **the self-seat did not rubber-stamp itself in either round**: it produced the strongest objection both times and was the only reviewer still rejecting in round 2. Independent rejectors throughout: `gemini-3-1-pro` + this controller. **Controller's own independent evidence** (never laundering a sub-agent claim): `verify_ail.sh` rc=0 at EXACTLY **4/4 identities / 9 modules / 14 tests**, re-run in the worktree after the doc landed to prove the manifest is untouched; `verify_go.sh` **rc=0** with `host/replay` **RUNNING 14.047 s not SKIP** — the designer's sandbox denied `bind(2)` (`listen tcp 127.0.0.1:0: bind: operation not permitted`, quoted verbatim) and it correctly **declined to claim the Go gate green**, the fourth consecutive milestone in which the codex lane refused to fabricate; `go.mod:3` = `go 1.26.4` so `http.NewResponseController` is available; `submit_feedback`'s description captured first-party. Scope clean **by diff**: exactly ONE new file plus this charter and the log; **no `.ail` added** (the designer's P8 justified it — protocol/session invariants are host-boundary behaviour, not pure-kernel law — so the required-check manifest legitimately does not move, and no reviewer challenged it). **ROUTING**: designer `codex:gpt-5.6-sol` ×2 (author + bounded revision), both on the ChatGPT subscription lane with `env -u OPENAI_API_KEY` load-bearing (the ambient key WAS set); controller applied the carve-out revision inline; **no planner / executor / evaluator ran** — the item never reached the sprint lane. **`metered=$0.1363`** (quorum reviewers only: $0.0634 + $0.0729; $5 ceiling untouched). **THE ITER-23 PATH FIX IS CONFIRMED WORKING**: the driver exported `MISSION_EXECUTOR_MODEL=codex:gpt-5.6-sol` this fire — not the spurious `opus` of iters 21/22 — so `PATH=/opt/homebrew/bin:$PATH` in `~/.config/ailang/mission-world.env` landed in exactly the window predicted, and the ratified pin now survives the driver's own pre-flight. Preflight clean: armed, billing **CLEAN**, `sunholo-voight-kampff`, dev==origin/dev (`503659d`), CI **completed/success** at HEAD (this repo has exactly ONE workflow — Build-and-Release/Docs-Deploy are **N/A, not pending**), no `[nightly-eval]` issues, inbox triaged to zero (one V1 eval-suite start notification, noise for World), no new `@MarkEdmondson1234` comment on `#9` (11 comments, all bot reports) or predecessor `#1`, no rotation due (`#9` titles this week; 11 ≪ 80). **ONE PROCESS FIX, no skill edit** (World never edits the shared skill): **a liveness/exit-code check is only evidence if it refers to the process you actually launched** — I hit this class TWICE in one iteration, once benignly and once not. `wait "$pid"` on a pid launched inside a nested `( … & echo $! )` subshell returns **rc=127** ("no such job") even though the codex probe had replied `ok` — a rc=127 that means nothing, in a mission where rc=127 has already caused four wrong diagnoses; and a `pgrep -f "codex exec --model"` completion poll **self-matched its own shell's command line**, so "is it still running?" was permanently TRUE and one poll I wrote had no deadline at all, which is my own Standing-rule-6 violation. Both are the iteration-107 lesson in a new costume: **the skill ships `kill -0 "$pid"` on a captured pid precisely to avoid this, and I hand-rolled a variant.** **Next: the queue's next actionable item is `w-store-durability` (item 4b, clause-1, NEW-DOC, ~1–2d)** — and note it is now *more* attractive, not less: `gpt5-6-sol` has independently raised the same commit-durability question from two different directions (the M3 dispatch→record crash window in iter-23, and this iteration's commit point-of-no-return), and its half (i) **CF-B-2** is in scope under BOTH answers to the parked question, so it can start without waiting on Mark. `w-mcp-projection` unparks only when `#498` ships a seam AND clause 3 lands.

## STATUS 2026-07-27 (iteration 23) — **`w-effect-broker-m3` DESIGN DOC WRITTEN + TWO QUORUM ROUNDS RUN → **PARKED `needs-human-review`** on ONE scope question; no sprint ran.** NEW-DOC lane (grep confirmed the tag was a FACT: no doc existed). Designer = **Fable** (rotation, next after codex), bounded 30-min run → `design_docs/planned/w-effect-broker-m3.md`, 1,036 lines: the capability/budget LAW frozen in a **Z3-proven sketch** with a Go drift-test mirror (the `worlddapi` precedent), **zero schema and zero log-format change** (records + approvals are content-addressed store objects), every decision recorded with **denials first-class**, a Replay mode that never dispatches handlers, a Model handler as a subprocess over the pinned binary (zero new cloud deps), and an isolation floor stated as **six named, individually-testable process restrictions with its non-containments enumerated** rather than an aspiration. **THE QUORUM EARNED ITS KEEP TWICE, AND THE SECOND TIME IT PARKED US.** Round 1 (`gemini-3-1-pro`, N−1 — `gpt5-6-sol` refused PRE-FLIGHT at the default $0.10 cap, est. $0.1160, **zero spend**, recorded by name) **REJECTED** on a real internal contradiction: Decision 4's two-phase `Human.Approve` "completes" an `EffectRecord` that Decisions 3 and 7 define as an **immutable content-addressed object** — and a content address IS the content, so there is no completing one. Fix adopted in full via a bounded revision: `Human.Approve` is now strictly synchronous (`Invoke` writes the request object, returns `Pending(requestRef)`, and synchronously writes ONE immutable record), `DecideApproval` writes a SEPARATE decision object and only moves a registry head, and observation is a new normal brokered effect **`Human.PollApproval`** with its own capability, budget line and record — propagated to 15 sites plus a new named RED mutation **`MUT-REC-IMMUT`** with **two independent red paths** (byte-identity re-read of the record captured at `Invoke` time; a store-integrity sweep asserting stored hash == hash(bytes)). Round 2 raised `--max-cost-usd` to 0.25 **specifically to buy back the reviewer round 1 lost** — it worked, both were present and independent (`gpt5-6-sol` $0.1129, `gemini-3-1-pro` $0.0471) — and **BOTH REJECTED**. Neither disputes the design DIRECTION. **`gemini`'s objection is carve-out-eligible** (pure completeness: P7/A9 claim named timeout + output-cap bounds on the Git/Model handlers but the Non-Vacuity table tests only the capsule floor F5/F6; two ready-to-paste mutation rows supplied) → **PRE-APPROVED to apply verbatim on unpark, no re-quorum**. **`gpt5-6-sol`'s is the park reason**: the broker dispatches a handler and *then* writes the record, so process death in between leaves a REAL external effect (`FS.Write`, `Git.Commit`, a paid `Model.Infer`) with no durable record, and replay cannot distinguish "never executed" from "executed but record lost" — a silent-duplicate-execution risk that contradicts the milestone's own headline claim. **This was NOT force-applied under the narrow-refinement carve-out, deliberately**: its fix adds a durable pre-dispatch intent object, a broker journal head, per-handler idempotency/reconciliation rules and crash-injection tests — a durability-ARCHITECTURE change overlapping the open **CF-B-2** store-hardening carry-forward. That fails the carve-out's second limb (scope/ratification judgment, not a verbatim text substitution); applying it would be the controller authoring a substantial design while calling it a reviewer's fix. Guardrail honoured: **park, do not force through**. The human question is binary and answerable in one comment — **(a)** close the crash window inside M3 (~+0.5–1d) or **(b)** keep scope and weaken the claim to the reviewer's OWN wording ("every attempted dispatch is durably detectable; completed outcomes are replayable; indeterminate attempts fail closed without live fallback"), journal → its own item beside CF-B-2. **Controller recommends (b); default if unanswered is (b).** **THE HEADLINE PROCESS FINDING — FOUR ITERATIONS OF WRONG DIAGNOSIS, CLOSED.** Iters 18/19/21/22 each re-derived the same symptom (codex `rc=127` → the ratified `codex:gpt-5.6-sol` pin silently demoted to opus) and each landed on a DIFFERENT wrong cause: spent quota, then an unusable model pin, then **the shared frozen driver** — which iter-21 filed upstream as `ailang#493` asserting "the V1 loop is demoting every fire too". **mission-v1 refuted that from inside its own fire, and iter-23 confirmed the real cause first-party: the defect was OURS.** `grep -c PATH ~/Library/LaunchAgents/dev.ailang.mission-world.plist` → **0** — the World plist sets no `EnvironmentVariables` at all, while V1's supplies `/opt/homebrew/bin`. Driver line 44 **prepends**, so it is correct-but-dependent: fine for any mission whose plist gives it a usable base, broken for one that does not. `gh`, `go` and `node` were collateral, not just codex. **Fixed with no frozen-core edit and no launchd reload**: the driver sources `~/.config/ailang/mission-<name>.env` at line **48** — after line 44, before the codex pre-flight at line 297 — so `PATH=/opt/homebrew/bin:$PATH` in that file lands in exactly the right window; verified by replaying the driver's own ordering under `env -i` (codex/gh/go all resolve). Correction posted to `#493`; ack + the reusable pattern posted to v1's thread `#484`. **The durable lesson is the mis-attribution, not the path**, and it is now a charter rule: (a) when one symptom yields a THIRD distinct diagnosis, suspect the part you have never inspected — nobody had ever looked at our own plist; (b) **before blaming shared/frozen infrastructure, check whether a peer consumer of it is healthy** — one `grep` on V1's plist would have refuted the iter-21 filing before it was written. A defect in shared code should reproduce for every consumer; one that does not is a local environment defect in a shared-code costume. **THREE LANGUAGE DEFECTS ROUTED AS `sunholo-data/ailang#495`, ALL REPRODUCED FIRST-PARTY** (the designer reported them; the controller re-ran every one rather than laundering the claim — and **two were WRONG as stated**): **U3 confirmed and sharpened** — the two toolchain legs CONTRADICT each other on identical source: `ai-check` PROVES a `requires`-guarded `debit` correct while `ailang test` FAILS it (`ensures violated for input: budget=-553, cost=-762`) because the derived ensures property ignores `requires`; **U2 was UNDERSTATED** — nullary ADT constructors fail too (`*ast.Identifier`), not just applied ones (`*ast.FuncCall`), and `ailang check` passes the file CLEAN so only the test leg catches it; **U1's stated cause was REFUTED** — my two minimal repros (a two-record-sort callee, and the `(record, string)` callee the error itself names) BOTH verify clean, so "params mix two record sorts" is not the trigger. The failure IS real (restoring the composed body flips `effectAllowed` to `status: error` while six predicates still verify, Z3 diagnostic captured verbatim), so it was filed **with its cause explicitly open and a labelled hypothesis** rather than a confident wrong cause. Independent of root cause: **`ai-check` exits rc=0 with a `status: error` result** — an encoding failure is indistinguishable from a pass at the process boundary, which is exactly why our gate asserts a hardcoded identity manifest. **Controller's own independent evidence** (never laundering a sub-agent claim): `verify_ail.sh` at EXACTLY **4/4 identities / 9 modules / 14 tests**; `verify_go.sh` green with `host/replay` **RUNNING 13.99 s not SKIP**; the doc's Appendix-A sketch re-run by me → **7/7 `verified`, 0 errors, `verify.results[]` enumerated** (so z3 genuinely ran — not the V27 silent-skip) and **27 passed / 0 failed / 32 total**; and Appendix A **diffed BYTE-IDENTICAL** to the sketch I verified, both before and after the revision, so M3.A's "lands this verbatim" is sound. Scope clean by diff: only the new doc + this charter. **ROUTING**: designer Fable ×2 (author + bounded revision) — a **deliberate, FLAGGED deviation** from the one-Fable-run-per-iteration clause, taken because both quorum reviewers are non-Anthropic and routing the revision to the rotation's codex entry would have made `gpt5-6-sol` judge its own sibling's work, breaking generator≠judge on the very gate that then caught the contradiction; both runs are subscription-bucket, so the clause's COST intent is untouched. No planner/executor/evaluator ran — the item parked before the sprint lane. **`metered=$0.2004`** (quorum reviewers only: $0.0405 + $0.1129 + $0.0471; $5 ceiling untouched, every model lane on a subscription bucket). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (`d1d1a9c`), CI green, only-workflow `CI` (no Build-and-Release/Docs here), no `[nightly-eval]` issues, inbox triaged to zero (mission-v1's iter-108 report — its two verdicts on World's asks acknowledged on `#484`), no new `@MarkEdmondson1234` comment on `#9` (10 comments; newest EQUALS the watermark `2026-07-27T08:55:11Z`) or predecessor `#1`, no rotation due (`#9` titles this week). **Next: `w-effect-broker-m3` unparks STRAIGHT TO sprint-planner** the moment the (a)/(b) question is answered — the doc needs no re-design and `gemini`'s fix applies verbatim on the way in. If the park persists, the queue's next actionable item is **`w-mcp-projection`** (clause-6, ~1d).

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
4. [**PARKED `needs-human-review` 2026-07-27 (iter-23)** — DOC WRITTEN + TWO QUORUM ROUNDS RUN;
   blocked on ONE scope question only. Doc: `design_docs/planned/w-effect-broker-m3.md` (Fable
   designer, rotation; 1,036 lines; Appendix-A sketch **7/7 Z3-verified / 27 tests / 0 failures**,
   re-run first-party by the controller and byte-identical to the doc's appendix).
   **Unparks straight to sprint-planner** once the question below is answered — no re-design
   needed] **w-effect-broker-m3** · clause-3 · effect broker with FS / Git / Model (`std/ai`) /
   Human.Approve handlers; effect-result recording; capability + budget checks; first physical
   isolation floor · ~2–3d. Builds on the landed `w-worldd-m2` daemon, whose single-writer
   authority is exactly the property that makes broker-mediated effects meaningful (an embedded
   writer bypassing capability/budget checks is the ambient-authority pattern clause 3 exists to
   end — Mark's ratification rationale, iter-18).

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
4b. [**PARKED `needs-human-review` 2026-07-28 (iter-25) — DOC LANDED + 2 QUORUM ROUNDS + carve-out
   revision applied; the design is DONE and the block is a RATIFICATION PACKET, not a defect.**
   Doc: `design_docs/planned/w-store-durability.md` (Fable designer, rotation; 1,058 lines) +
   `design_docs/sketches/storejournal.ail` (163 lines, **7/7 contracts Z3-verified**, 25 named
   tests — controller-remeasured). **Unparks straight to sprint-planner** the moment the packet
   below is answered; no re-design needed.

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
   all eight zero-value commits return `err=<nil>`; **seven** produce a permanently unreadable row
   (`TransitionFn`, `Interpreter`, `EntryHash`, `TransitionRef`, `PrevEntryHash`,
   `NextWorld.LogHead`, `NextWorld.StateRoot`), and the eighth (`NextWorld.Ref`) commits and reads
   back *fine* as an empty-string ref that becomes the selected head — degenerate-but-readable, and
   therefore the one shape a read-side fix could never observe. Also measured: the poisoned head
   ADVANCES, and a subsequent perfectly legal commit chains onto the poisoned entry, so the
   append-only log grows a **permanent hole mid-chain with readable entries on both sides** and no
   detection or recovery path]
   **w-store-durability** · clause-1 · store/kernel hardening across a crash:
   **(i) CF-B-2** — `store.Commit` writes a zero `PrevEntryHash` that `store.GetLogEntry` **cannot
   read back**, so the embedded Go API can append a log entry no reader can ever load. The daemon
   refuses it at the REST boundary, which is a boundary patch over a kernel defect. Discovered
   iter-21 **by a failing test**; as of iter-23 it still has **no issue and no repro fixture** —
   this row closes the queue half, and a committed repro fixture is the first deliverable (the
   ghost-close rule: never bare bookkeeping).
   **(ii) the broker effect journal** — a durable pre-dispatch intent object + journal head + a
   post-dispatch outcome object, with an `IndeterminateEffectError` recovery path that **never**
   auto-re-executes, and per-handler idempotency/reconciliation rules before any retry is allowed.
   Raised by `gpt5-6-sol` in the `w-effect-broker-m3` round-2 quorum: without it a real external
   effect (an `FS.Write`, a `Git.Commit`, a **paid** `Model.Infer`) can occur with no durable
   record, and replay cannot distinguish "never executed" from "executed but record lost".
   Both halves are the same question — **what the kernel guarantees across a crash** — which is why
   they share one item rather than being answered twice · ~1–2d · surfaced iter-21, queued iter-23.
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
6b. [DOC-DRAFTED 2026-07-27, attended — quorum + ratify at pick; GATES items 7 and all M6 work]
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
9. [BACKLOG — infra, not critical-path] **w-verify-binary-lockfile** · clause-1-infra · a durable
   pin of the released `ailang` binary (version + sha256) for the `ailang-code` verify gate, so
   local controller verification stops depending on the rig's `-dirty` build (iter-6 pinned
   released `v0.30.0` ad-hoc — `sha256:ac3174e0…`; CI already sha256-verifies a released linux
   binary). Small; the mechanism may generalize to the SHARED `ailang-code` profile → confirm
   repo-local lockfile vs shared-skill fix with a human before implementing (do not hand-edit CI
   headless) · ~0.5d · surfaced iter-6, queued iter-7

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
| Driver sources `~/.config/ailang/mission-world.env` + respects role overrides | code + LIVE 2026-07-23 | `tools/launchd/mission-control.sh:46-47` sources `mission-${MISSION_PROFILE}.env`; `:238` `MISSION_EXECUTOR_MODEL` respected (default opus); dry-run log 19:45:39 echoes env-only values (repo-slug/doc/workdir) + resolved roles — sourcing proven end-to-end |
| Charter RATIFIED (authorization state, kill switch, sprint routing) | attended session 2026-07-23 | Mark's decisions recorded: clause-4 = −2pp/≤25% paired N≥3; bar+Conflict Surface+guardrails+queue as drafted; ratification comment on issue #1; STATUS stamp above is the in-doc record |
| Kill-switch + gh-issue state paths are world-namespaced in the LIVE driver (safety property) | LIVE 2026-07-23 | driver log 19:45:24/19:45:39/20:02:47: "kill switch present (`…/mission-world.disabled`) — skip" — the driver printed and HONORED the world-namespaced path 3× while the v1 loop ran unaffected; iter-0 posted to issue #1 = `mission-world-gh-issue` read correctly |
| `ailang messages` channel works end-to-end (guardrail's delivery leg) | live round-trip 2026-07-23 | sent `msg_…_2c6964d3` (defect report) + `msg_…_acc5edcc` (channel test); v1 agent RECEIVED and acted — upstream ack on issue #1 at 18:27Z citing the report, fix `ailang@aabb3a58c` |
| v1 session-start hook reads the message inbox | config + observed 2026-07-23 | ailang repo `.claude/settings.json` SessionStart → `scripts/hooks/session_start.sh` ("checks the user inbox … using the ailang messages CLI"); displayed 5 unread at this session's start |

---
**Document created**: 2026-07-23 (bootstrap, attended). **RATIFIED 2026-07-23** (iteration 0,
attended: Mark + World coordinator) — record on issue #1; advisory-quorum ledger in
`.ailang/state/mission-quorum/`. Sprint routing is authorized from the next loop fire.

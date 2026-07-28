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

---

## STATUS (rotation rule)

Newest **3** STATUS stamps live here; older ones move to `world-mission-status-archive.md`.
At Gate 4, after adding your stamp, move the now-4th stamp to the TOP of the archive file.

## STATUS 2026-07-28 (iteration 29) — **`w-store-durability` SD.B LANDED: THE DURABLE JOURNAL + IN-TX COMMIT RECEIPTS** (PR #18 → squash `d5774eb`, dev CI green both jobs; judge sonnet **PASS 94/100, zero blocking**; executor `codex:gpt-5.6-sol`). ARM J1 + Decision 4 shipped: additive `journal` table (every existing `CREATE TABLE` byte-unchanged), `host/store/journal.go` (+470 — types, deterministic codecs with golden bytes, `AppendIntent`/`AppendOutcome`/`GetReceipt`/`PendingIntents`, in-tx gapless `seq`), `store.go` (+83 — `Commit.InvocationID`; `bindCommitIntentTx` compares ALL EIGHT commit-defining fields **inside the existing transaction, before any mutation**; the receipt is written in that SAME transaction), `journal_test.go` (+390, 7 tests incl. the 10-row drift test). `InvocationID == ""` is byte-compatible with every landed caller. Closes **AC5, AC7, AC8, AC9, AC13, AC15, AC10-PendingIntents-half, AC12**. **TWO findings, both a prescribed fix that was itself vacuous.** (1) **Iter-28's prescription for the sketch precondition was WRONG, and measuring it before adopting it is what caught that (THIRD instance of the one root cause).** It specified ONE new `tests[]` row and pinned `len(tests[])` 25→26 / `passed_tests` 32→33 — numbers that reached the doc, this charter and the plan JSON. Measured on the pinned binary: at **26 rows, dropping `TransitionRef` alone from the Go compare reds NOTHING** (`failed=0`), because the round-2 REQUIRED `EntryHash`-preserving row mutates `PrevEntryHash`/`TransitionFn`/`Interpreter` *together* and never touches `TransitionRef` — while `MUT-INTENT-NARROW-BIND` itself demands the four added fields be load-bearing **"individually and not decorative"**. Landed form is **10 rows** (all-match + one single-field row per commit-defining field + the REQUIRED combined row); **AC9 now pins 30 / 37** (premise row **V28**). It was load-bearing end-to-end: `MUT-DROP-TRANSITIONREF` reds exactly `row-8-transition-ref`, **reproduced first-party** by the controller, reverted byte-identical. The planner's "16 params is 2× the widest arity ever proven" risk is **REFUTED** — 7/7 verified, 0 counterexamples, first try; no upstream issue owed. (2) **AC6 was owned by NO milestone** — found by re-checking a NON-BLOCKING judge nit instead of filing it. `AC5–AC8` gave SD.B ownership of AC6, whose only proof is crash injection, whose file `crash_test.go` is in **SD.C's** list, while SD.C's list read `AC10–AC11, AC14`; SD.C's close-out would never have forced it, and `MUT-SPLIT-TX` belonged to a test nobody owned — **an acceptance check no gate can fail**. Reassigned to SD.C at `9316286` in all three places. Also: the codex sandbox again denied loopback binds (6 tests + 4 benchmarks), so `verify_go.sh` and `bench_worldd.sh --smoke` were **re-run outside it** — both PASS; and a **baseline** bench was taken at `857a912` BEFORE the sprint so the after-reading is a comparison (`BenchmarkStoreCommit` 0.400 → 0.405 ms, no tax on the empty-ID path). **Queue: 4b → SD.C remains (AC6 now explicitly owned there), then 4.** Carry-forwards CF-F-1/CF-F-2/CF-F-4 + new CF-G-1/CF-G-2/CF-G-3 open. 4 stamps → newest 3 kept, iter-27 archived this Gate 4.

## STATUS 2026-07-28 (iteration 28) — **`w-store-durability` SD.A LANDED: CF-B-2 IS CLOSED AT THE KERNEL WRITE PATH** (PR #17 → squash `86d1276`, dev CI green; judge sonnet across 3 rounds FAIL 57 → PASS 91 → MERGE 97; executor `codex:gpt-5.6-sol`). ARM V1 shipped: `store.Commit` validates **EIGHT** ref fields (5 entry + 3 world) plus each Object's `Hash`/`InterfaceHash`, before `tx.Begin()`, with structured `InvalidRefError`; the same helper guards `PutObject`/`PutWorld`/`SetRegistryHead`/`SelectHead`/`PutVerifyResult`. `ObservedHead` stays the one zero-legal COMPARED field. New `host/store/scan.go`: bounded read-only `ScanUnreadableLog`/`ScanUnreadableWorlds`, **keyset cursors only** (never `OFFSET`, never rowid), `MaxIntegrityScanPage` + `InvalidLimitError` enforced before any query; daemon pages both tables at startup and reports `integrity_hole` / `integrity_scan_incomplete` with a continuation cursor. Closes **AC1–AC4, AC12, AC10-scan-half**; `schema.sql` byte-unchanged. **THREE findings, all of them our own artifacts disagreeing with each other**: (1) quorum round 2's "seven→eight" correction reached premise V23 and STOPPED — Decision 1's prose and AC2 still said seven, and V23 still carried the *superseded* "degenerate-but-readable" claim iter-27 retracted; implementing seven drops `NextWorld.Ref`, the CLASS 3 **wedge**, leaving the worst failure mode open while every gate goes green (fixed `a25d87f`; now EXECUTABLE — `MUT-SEVEN-NOT-EIGHT` reds 2 tests). (2) **SD.B BLOCKING PRECONDITION**: `sketches/storejournal.ail:132` LAW 6 `intentBindsCommit` still declares the round-1 NARROW 4-field binding (8 params) while the ratified Freeze/Decision 4/AC15 require the 8-field form (16 params) — the Freeze pins the Go mirror to the sketch by drift test, so SD.B as written would certify the very binding the doc calls the defect. Applies the ratified Freeze ⇒ no human needed, but MUST precede SD.B's drift test (AC9 counts 25→26 / 32→33). (3) **I reported 3 gates green while 4 were binding** — the judge FAILED round 1 at 57 on a RED `bench_worldd.sh --smoke`, caused by `BenchmarkStoreCommit` seeding a zero `PrevEntryHash`, i.e. **a benchmark relying on the defect**. Also: the codex sandbox denies loopback binds, and that panic MASKED a real `io.Pipe` startup deadlock — my round-1 fix then only covered the healthy path and left it wedgeable for stores WITH holes/truncation (CF-F-3, closed structurally at `d506275`). **Queue: 4b → SD.B/SD.C remain (sketch arity first), then 4.** Carry-forwards CF-F-1/CF-F-2/CF-F-4 open. 6 stamps → newest 3 kept, iters 24/25/26 archived this Gate 4.

## STATUS 2026-07-28 (TRIPLE RATIFICATION, attended ~09:00) — **Mark answered the one-reply gate: items 4, 4b, 6b ALL RATIFIED as recommended.** (1) **4b `w-store-durability` packet RATIFIED IN FULL**: ARM V1 (validate every ref at the kernel WRITE path — iter-27's fixture evidence decisive: V2 cannot reach CLASS 3, one write-side check reds all three classes) · ARM J1 (durable intent/outcome journal as additive SQLite table + content-addressed payloads) · `Commit` stable invocation ID with in-tx receipt binding · three-state receipt law (not-started/indeterminate/resolved, never-lie, Z3-proven) · recovery NEVER auto-re-executes indeterminate intents. The M1 kernel reopen this entails is RATIFIED (human gate exercised). (2) **4 `w-effect-broker-m3` scope RESOLVED: broker DEPENDS on 4b's journal** — the crash window closes structurally, not by documentation; sequencing 4b → 4. (3) **6b `w-human-surface` v0.3 §7 FORMALLY RATIFIED**: principles + anti-patterns binding; trust grades stay PROVEN/TESTED/ATTESTED/CLAIMED; packet schema freezes when the inbox build STARTS; renderer FRESH (Hub = pattern donor); dialect deferred to M6 (common-core emitter); §7.5 resolved consistent with the ratified anti-patterns — CLAIMED-grade `AiReview` content renders ONLY with its evidence ref attached (anything else is confidence theater). **Queue: 4b → [NEXT] (sprint the fix half), then 4; 6b → [LANDED as design, binding]; item 7 (approval inbox) unblocks on 4.** Recorded in-charter per the ratification-channel pattern; 5 stamps present — next Gate-4 hand-rotates two.


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
4b. [**IN SPRINT — RATIFIED; SD.A LANDED (iter-28, `86d1276`) + SD.B LANDED (iter-29, PR #18 →
   squash `d5774eb`, dev CI green both jobs, judge PASS 94/100 zero-blocking); ONLY SD.C REMAINS.**
   SD.B shipped ARM J1 + Decision 4 — the additive `journal` table, `journal.go` (+470),
   `Commit.InvocationID` with the eight-field in-tx binding and the receipt written in the SAME
   transaction, and the 10-row drift test mirroring the frozen sketch — closing **AC5, AC7, AC8,
   AC9, AC13, AC15, AC10-`PendingIntents`-half, AC12**. **SD.C is the last milestone**: crash
   injection at named kill points (real subprocess kills), the probe-consumer recovery proof
   (never auto-re-execute), the two journal benchmarks into the hardcoded smoke manifest +
   `bench/BASELINE.md` re-measured in ONE invocation, **and AC6 — which iter-29 found was owned by
   NO milestone and reassigned here** (`9316286`). No human gate is outstanding for SD.C.
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

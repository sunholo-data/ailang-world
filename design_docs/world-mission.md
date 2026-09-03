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
- **`scripts/mission_directives.sh` DOES NOT EXIST IN THIS REPO — invoke it by ABSOLUTE PATH from
  the V1 checkout (process fix, iter-72).** Gate 0.6 says to run `scripts/mission_directives.sh`
  and, in the same breath, forbids hand-rolling the `gh | jq` pipeline it replaces — because the
  author allowlist that stops arbitrary commenters steering this loop lives *in that script*, not
  in prose a controller is trusted to retype. But the script is a V1 artifact and was never copied
  here (measured iter-72: absent from `scripts/`, present at
  `/Users/voightkampff/dev/sunholo-data/ailang/scripts/mission_directives.sh`). A controller that
  runs the relative path gets `No such file or directory` and is one step from doing exactly the
  thing the rule forbids — or, worse, from concluding the directive channel is unavailable and
  skipping it, which silently drops human directives that OUTRANK the queue. It takes `--repo`, so
  no fork is needed:
  `/Users/voightkampff/dev/sunholo-data/ailang/scripts/mission_directives.sh --issue "$ISSUE"
  --repo sunholo-data/ailang-world --since "$last"`. Do **not** vendor a copy (frozen core, and a
  second copy is a second allowlist to drift); if the V1 checkout ever moves, that is a
  human-decision row, not a licence to hand-roll.
- **INSTANCE 2 LANDED IN THE SAME ITERATION — PROPOSED TO V1 (iter-107; bar MET) — A STRESS CONTROL ONLY CERTIFIES THE AXIS YOU VARIED, AND A GREEN ONE READS AS IF IT CERTIFIED THE PROPERTY.** Iteration 107 re-ran `PE.E`'s wall-clock arm **15× unloaded and 8× under eight CPU spinners — 23/23 green**, hold ratio 26.7–30.2× against a 20× floor, with the loaded arm moving the ratio the *safe* direction — and recorded the timing as sound. The evaluator varied a different axis, **parallelism**, and got **10/10 FAIL at `GOMAXPROCS=1` on unmutated, sha256-identical code**, with the failure text `blocked read returned after 10–33ms`, which is **indistinguishable from the M22/M23 mutant signature the arm exists to detect**. Reproduced first-party with a two-arm control on one tree: `GOMAXPROCS=1` **10/10 red**, default **0/10 red**. Note why more effort would not have helped: every one of the 23 runs was the *same shape*, so the sample grew while its coverage did not — this is rule 3a(i)'s known-positive discipline aimed at a **stress** control rather than a search, and rule 3b(ix)'s scope-travels-with-the-count aimed at an **axis** rather than a directory. Candidate rule if a second instance lands: before quoting a timing or load result as evidence, **name the axes you held fixed** (parallelism, CPU contention, memory pressure, page cache, disk, clock granularity) and vary at least the one the mechanism under test actually depends on — here the stimulus is a *scheduling* race, so parallelism was the load-bearing axis and CPU contention was the decorative one. The tell: you are about to write "N/N green under load" and every run varied the same knob. **INSTANCE 2, ninety minutes later and from CI rather than from a judge:** the Gate-4 record commit `ac17d54` — **docs-only**, `git diff --stat daf48a6 ac17d54 -- ':!design_docs'` EMPTY — turned `dev` RED in the very arm above. The runner's log: `decoy hold=2.626683936s ObjectReadTimeout=2ms ratio=1313.3x` / `blocked read exceeded test-side 20x watchdog`. The stimulus SCALES with the machine (53 ms here, **2.63 s** on the runner — 49x) and the bounds did not, because they were absolute millisecond constants calibrated on this laptop; the resulting red is again **indistinguishable from the M22/M23 mutant signature**. Same tree passed on the parent, so the parent's green was luck. Three axes, three reds, and after each fix the surviving constant still encoded *my machine*: CPU contention (mine, 23/23 green), parallelism (the judge's, 10/10 red), absolute speed (CI's). **So the rule is stronger than "vary more axes": where a bound and its stimulus both depend on the machine, DERIVE the bound FROM the measured stimulus and stop hardcoding wall-clock at all.** Fixed at `a87c723`: `readTimeout := hold / 20` makes the doc's 20x floor true BY CONSTRUCTION on any machine, the watchdog becomes the hold itself, and a `minDecoyHold` floor keeps a too-fast decoy a loud instrument failure. M22/M23 still die. Proposed to V1 on the cross-mission channel; World cannot edit the shared skill.
- **AN ANTI-VACUITY FLOOR IS A BRANCH TOO, AND RULE 3j's ENUMERATION STOPS AT *REFUSAL* BRANCHES — SO THE CHECKS THAT PROTECT THE INSTRUMENT ITSELF ARE SYSTEMATICALLY THE LAST THING ANYONE PINS (watch-item, iter-108; 1 recorded instance, bar is 2 — NOT yet proposed to V1).** The skill's rule 3j says: for any function whose contract is a refusal, enumerate its refusal branches and require one neutering mutation per branch. That is written about the *subject* under test. A gate also carries floors that refuse on its **own health** — zero packages discovered, an empty required set, an empty enumeration — and those are invisible to the same discipline, because a floor is not a refusal *about the input*, it is a refusal *about the measurement*. Measured here on `PE.F`: §5 names **three** floors for the focused `host/evidence` leg and the isolated gate shipped arms for **two**. The `"zero passing host/evidence packages discovered"` message appears **once** in `scripts/verify_go.sh` and **zero** times in `evidence_manifest_gate_test.go`, and every synthetic writer emitted the package-level pass event **BY CONSTRUCTION**, so no existing arm could reach that branch and **no removal-shaped mutant on the other two floors would ever have revealed it**. Live code that no test protected, inside a milestone whose entire subject is non-vacuity, at a **96/100 zero-blocking** evaluation — the judge found it, not the drill. Closed with a sixth arm and a dedicated writer that omits the package event; non-vacuity proved by neutering the floor (`if not package_passes:` -> `if False and not package_passes:`), `go build` rc=0 before any result was read, red set = **the new arm ALONE**, the other five green. **Candidate rule if a second instance lands:** when a gate has floors that fail on the instrument rather than on the input, enumerate the FLOORS separately from the refusal branches and require one arm per floor — and note the tell, which is structural rather than behavioural: *every helper that builds the gate's input emits the field a floor checks for, so the floor's input shape cannot be produced by the test harness at all.* This mission has ADDED anti-vacuity floors in four prior iterations (measured: 4 log hits) and this is the first time one was found unpinned, so the bar is honestly **1**, and it is recorded here rather than proposed — filing below the bar is a defect this loop has already paid for once (iter-106).
- **Rig PATH (process fix, iter-18 — 3 instances in ONE iteration)**: the agent tool shell's `PATH`
  does **not** include `/opt/homebrew/bin`, so `gh`, `go`, `node` and `codex` all fail with
  `command not found` / `env: node: No such file or directory` (`rc=127`). Every Bash call that
  needs them must `export PATH=/opt/homebrew/bin:$PATH` first, and **every directive handed to a
  spawned sub-agent must say so** — otherwise a planner/executor reads it as a broken toolchain or
  a spent codex quota rather than a PATH gap (iter-18 lost a codex probe to exactly that
  misreading). The pinned AILANG binary lives OUTSIDE the repo at
  **`~/.pinned-ailang/ailang`** (v0.30.0, commit `e37b370`, clean) — always pass it explicitly as
  `AILANG_BIN`. **MOVED OFF `/tmp` 2026-09-03 (iter-151), because macOS WIPES `/private/tmp` ON
  BOOT AND THE OLD PIN PATH WAS `/tmp/ailang-v0300/ailang`.** Measured this iteration: the rig
  rebooted at `kern.boottime` = **`1788395029`** (2026-09-03 02:23:49 local) and every `/tmp`
  resident of this mission went with it — the pin (`ls` rc=1, so `AILANG_BIN` had nothing to
  point at), the driver's own log `/tmp/ailang-mission-world.log` (**217 bytes**, holding only the
  current fire's opening line, though the driver only ever *appends* — `tee -a`/`>>`), and the
  row-56 sprint worktree `/private/tmp/wt-row56` (left as a `prunable` entry in
  `git worktree list`). The gate fails CLOSED, which is the good half — `verify_go.sh` refuses
  loudly when `AILANG_BIN` is unset rather than letting `host/replay` `t.Skip()` into a false
  green — but the bad half is that the loop is then simply *unable to verify anything* until a
  human notices, and nothing tells anyone. **The new path is not invented: it is the one CI
  already uses** (`ci.yml` extracts the pinned v0.30.0 release to `$HOME/.pinned-ailang/ailang`),
  so rig and CI now name the same location, and `$HOME` survives a reboot. Restored and verified
  first-party rather than assumed: published checksum `shasum -a 256 -c` **OK**,
  `--version` = `AILANG v0.30.0` / commit `e37b370`, binary sha256
  `e9746fef8570bc42b8cc52c0e88b7088468a5d2bd38bb8c42e27e5859b8f3fb5` — and the strongest control
  is the gate's own, since `verify_ail.sh` step 9/9 prints `compiler pinned by exact bytes: AILANG
  v0.30.0 on Darwin/arm64` and the whole gate then reads `rc=0`. The generalisation, and it is
  this repo's own recurring shape aimed at a *filesystem* rather than at a path variable: **a pin
  stored in a volatile directory is a claim about the last boot, not about the toolchain.** The
  tell: your reproducibility anchor lives somewhere the OS is allowed to delete.
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
- **`ai-check` IS CWD-SENSITIVE, AND THE WRONG CWD REPORTS A KILL ON *EVERY* MUTATION ARM (process
  fix, iter-85).** A `world/*.ail` module imports its siblings by package path, so the module root
  matters: `verify_ail.sh`'s `ROOTS` entry is `.|world`, i.e. **base = repo root, rel =
  `world/types.ail`**. Run it the natural-looking way instead — `cd world && ai-check types.ail` —
  and it fails `LDR001: module not found: world/logepoch`, which surfaces in the JSON as
  **`check.passed=false, verified=0`**. That is **byte-indistinguishable from a mutant that
  destroyed the contract**, so a mutation drill run from the wrong directory reports a confident
  kill on every single arm, including arms that in truth change nothing. Iteration 85 banked two
  such kills before a **pristine known-positive control** (unmutated tree must read
  `check.passed=true, verified=6, cex=0, errors=0`) came back FALSE and exposed the instrument.
  Rule: any drill harness asserts that control **before and after every batch**, never once at the
  start — and mirror `ROOTS` rather than guessing the cwd. This is rule 3a(i) with the unusual
  twist that the broken instrument produces a *positive* result, not an empty one, so the
  "search found nothing" reflex does not fire.
- **A DUPLICATE WEBHOOK DELIVERY MAKES `checks` EXCEED `expected`, AND GATE 3b's COMPLETENESS RULE
  — WRITTEN ENTIRELY FOR A *SHORT* CHECK SET — THEN REFUSES TO REACH A VERDICT ON A PERFECTLY
  GREEN COMMIT (process fix / watch-item, iter-129; 1 recorded instance, bar is 2 — NOT yet
  proposed to V1).** The shared skill's rule (i) says to enumerate the workflows EXPECTED for the
  diff and require every one to be PRESENT, *"else print `INSTRUMENT INCOMPLETE — no verdict` and
  keep polling; a count of what you found is not a count of what exists."* Every word of that is
  aimed at `present < expected`. The natural implementation is an equality — the skill's own
  rule (c) even names `present == expected` — and equality fails in **both** directions.
  **Measured first-party here, in the same incident-recovery window that caused the drops this
  iteration was cleaning up:** the record push `68164ba` received **TWO** `CI` runs — same
  workflow file, same `event=push`, the **same `created_at` second** (`19:46:43Z`), both
  `run_attempt=1` — while the control commits `8e3c8cd`, `2d2513e`, `609e090`, `74c47d5` and
  `fd490ca` each received exactly **1**. `commits/<sha>/check-runs` therefore returned
  **`checks=4`** (each of the two jobs twice) against `expected=2`, so the naive rule prints
  `INSTRUMENT INCOMPLETE — no verdict` and a bounded poll runs to **TIMEOUT and PARKS a landing
  that is green**. Note this is the exact mirror of the failure the same gate already documents:
  an aggregate over a SHORT set is vacuously *green*, and an equality over a LONG set is
  spuriously *incomplete* — at-most-once and at-least-once delivery, one comparison, both
  directions.
  **The fix, and it is a no-op on the normal case, which is what makes it safe to adopt:** compare
  **DISTINCT check NAMES**, not row count — `[.check_runs[].name]|unique|length`. Measured: on the
  double-delivered commit `distinct=2 == expected=2` (verdict reachable) while `total_count=4`;
  on the single-delivery control both readings are **2**, so the two rules differ only where they
  must. Then require every distinct name to carry at least one success and **zero** hard reds
  (`failure`/`cancelled`/`timed_out`) — a duplicate delivery can legitimately leave one copy
  cancelled by concurrency rules, and a name-scoped read handles that while a row count cannot.
  The tell: your Gate-3b completeness assertion is an `-eq`, and you have only ever asked what
  happens when the number comes back too small.
- **A QUEUE ROW KEEPS ITS OWN DEAD HEADS INLINE, SO A ROW-SCOPED GREP RETURNS THE LIVE HEAD AND
  EVERY SUPERSEDED ONE WITH NOTHING IN THE MATCH TO TELL THEM APART (process fix, iter-87; 1
  recorded instance — watch-item, bar is 2).** This charter's convention is that a row records its
  history in place behind `**Prior head text follows.** ~~…~~`; **14** rows carry that marker and
  **5** more carry `Prior row text follows`, so ~19 spans of *deliberately dead* text sit inside
  live rows. `grep` does not see `~~`. Iteration 87 needed the next unblocked item, grepped item 8's
  row, and got **`HEADLESS-ROUTABLE`** and **`NOT BLOCKED ON MARK`** — both real strings, both
  inside struck-through spans from two earlier epochs of that row, while the LIVE head (line 2293)
  says the opposite: *"`SM.D` … IS ATTENDED-ONLY — never headless, never CI. THIS ITEM HAS NO
  HEADLESS-ROUTABLE MILESTONE LEFT."* A `grep -c` returned **12** status tokens spanning three
  epochs, and a count cannot separate them. Note what makes this worse than an ordinary stale read:
  the dead text is not stale by neglect, it is **preserved on purpose**, so it is well-written,
  confident, and often *more* specific than the head that replaced it. Rule: to establish a row's
  state, read the row's **FIRST line** (`sed -n '<n>p'`), or strip the fenced spans before grepping
  (`perl -0pe 's/~~.*?~~//gs'`); never conclude a row's status from a token count over the whole
  row. Same family as rule 3b(ix) — a claim true in the scope it was written for, quoted into a
  wider one — with the scope being *which head is live*.
  **INSTANCE 2 — 2026-08-15 (iter-88), A DIFFERENT MECHANISM, SAME CLASS: THE SUB-ROW NAMING
  CONVENTION IS INVISIBLE TO THE NATURAL PATTERN, SO THE QUEUE UNDER-COUNTS BY A FIFTH. The bar
  of 2 is met and this is now a recorded rule, not a watch-item.** Iteration 87 wrote *"All **19**
  rows (18 numbered + `4b`)"*. Enumerated by parsing row heads, the queue is **24**: `1, 2, 3, 4,
  4b, 4c, 4d, 4e, 4f, 5, 6, 6b, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18` — missing `4c`, `4d`,
  `4e`, `4f`, `6b`. The cause is that a `^[0-9]+\.` reading cannot see a letter-suffixed sub-row,
  and this charter grows them freely (an item that splits keeps its parent's number). **What makes
  it durable is that the VERDICT survived**: every missed row is `LANDED`/`COMPLETE`, so "all rows
  are complete or blocked" was still true, and nothing about the sentence looked wrong. A count
  that is wrong in a direction the conclusion does not care about gets no scrutiny at all.
  Rule: enumerate queue rows with `^[0-9]+[a-z]?\. \[` and **quote the enumerating command beside
  the count** (rule 3b(ix)); never state a queue cardinality from a `^[0-9]+\.` pattern or from

  **INSTANCE 3 — 2026-08-23 (iter-114): THE RULE ABOVE IS NOW ITSELF THE UNDER-COUNTING
  INSTRUMENT, AND IT MISSES THE ROW THIS ITERATION PICKED.** The prescribed pattern
  `^[0-9]+[a-z]?\. \[` demands a `[` IMMEDIATELY after the space, and this charter now writes row
  heads in **bold** — `14. **[[IN-SPRINT] …`, `22. **w-daemon-lock-wait-not-deadline-bound**` — so
  the anchor fails. Measured over the Queue section: prescribed **23**, widened `^[0-9]+[a-z]?\. `
  **42**, i.e. **19 invisible rows** — `14 17 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36`
  — with a fresh absent row-head shape as the negative control at **0**, so the pattern is not
  broken, it is NARROW. Row **14** is `[IN-SPRINT]`, so the enumerator could not see the item the
  loop was working. Note this is the THIRD distinct mechanism (struck-through dead heads →
  letter-suffixed sub-rows → bold heads) and that the VERDICT survived a third time, because every
  invisible row is LANDED or live — which is precisely why nobody re-checks it. **Rule, amended:**
  enumerate with `^[0-9]+[a-z]?\. ` (no `[`), quote the enumerating command beside the count, AND
  pair it with a WIDENED control in the same call, requiring the two to agree or recording the
  delta — a same-shape control cannot detect a narrow anchor. The general form, which is the part
  that will outlive this fix: **a remedy is an instrument too, so a rule that prescribes a pattern
  inherits that pattern's blind spots — and a formatting convention introduced AFTER the rule is
  written is invisible to it by construction.**
  memory of the last row's number. **The generalisation across both instances is the one to carry:
  this charter's own formatting conventions — struck-through dead heads, letter-suffixed sub-rows —
  are each invisible to the instrument a reader naturally reaches for, and both failures are
  SILENT and produce a confident, specific, wrong answer.** Before grepping this file for anything
  structural, ask what convention the pattern cannot see.
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
- **WHEN A MILESTONE'S DELIVERABLE IS A *REFUSAL SET*, THE UNIT OF MUTATION IS THE **BRANCH**, NOT
  THE MILESTONE — A PER-MILESTONE MUTATION LIST SYSTEMATICALLY UNDER-COVERS A FUNCTION WHOSE WHOLE
  JOB IS TO SAY NO IN N DISTINCT WAYS (process fix, iter-63 — THREE first-party instances in ONE
  iteration, found by three different roles).** Every mutation discipline this mission already has
  is aimed at a mutation someone NAMED: the plan's `named_mutations`, the doc's mutation table, the
  shared skill's rule 3i (a test-plan row's "kills which mutation" column). None of them points at
  the branches of a validator the executor WRITES during the sprint — so those branches ship with a
  green suite and no pin, and the suite's green is what makes the gap invisible. Measured here on
  `host/broker/validatePublishApproval`, whose entire purpose is to refuse: **(1)** the inherited
  draft's `request.Effect == EffectRegistryPublish` was satisfiable by *nothing an operator can
  mint* (the doc's own premise was false); **(2)** its replacement,
  `scope.Effect != EffectRegistryPublish`, was **unpinned** — `if false && …` (mutant LANDED by
  sha256, `go build` rc=0) left the **ENTIRE `host/broker` package green**, because the neighbouring
  AC9 "malformed" arm rejects its scope at the PARSER and never reaches the term; **(3)** the
  evaluator, handed (2) as a named target, found **six more**, and the executor's own audit of the
  function found **twelve** — six policy branches plus six traversal-error branches nothing reached
  either. `AC9` ended at **20 negative arms, one per refusal branch**, all seven policy branches
  re-verified RED by the controller with the failing test NAMED each time. The rule, and it is cheap
  because the enumeration is mechanical: **for any function added or modified by a milestone whose
  contract is "refuse X", enumerate its refusal branches (`grep -c 'return .*fmt.Errorf(.*%w'` over
  the function is a fair first cut), and require one neutering mutation per branch before the
  milestone is closed.** Neuter with `if false && <cond>` rather than deleting the block — it keeps
  every import used, so the mutant compiles and "the mutant doesn't build" cannot masquerade as "the
  guard fired". A branch that genuinely cannot be reached is an acceptable outcome **when it is
  declared in the code and in the AC** (as `journal.go`'s PK-collision fallback now is); an
  undeclared one is a guard nobody is protecting. **And read WHICH TEST failed, never the exit code
  alone** — one of this iteration's probes returned `rc=1` in exactly the predicted direction and the
  only FAIL was a pre-existing load flake. **PROPOSED to V1 as a shared-skill rule** (World cannot
  edit the mission-control SKILL.md — it lives in the V1 checkout), since the gap sits in the
  rule-3i family and is mission-independent.
- **WHEN THE DELIVERABLE IS A *DETECTOR*, RULE 3j's BRANCH ENUMERATION IS ONLY HALF THE GATE —
  ENUMERATE THE **SHAPE SPACE OF WHAT IT REFUSES**, BECAUSE A RECOGNISER'S COVERAGE IS A PROPERTY
  OF ITS INPUT GRAMMAR, NOT OF ITS BRANCH COUNT (process fix, iter-75; instance 1 was iter-69).**
  Rule 3j asks *how many ways can this refuse*, and this repo now answers that well: `TR.C` ran
  **32** mutation arms — 23 executor, 9 controller — every branch of the gate pinned, zero
  survivals, whole-package inverses clean. The gate was still **defeated by ordinary Go**, because
  all 32 arms spelled the forbidden thing **the same way**. Every detector for `Invoke`,
  `NewSession` and `NewReplaySession` sat inside `case *ast.CallExpr`, so it matched only when the
  selector is the `Fun` of a call; a **method value** (`call := s.Invoke`) or a **function value**
  (`mk := broker.NewSession`) reached a raw broker session from outside `host/broker` with no
  reflection, no `//go:linkname`, no generics and no build tags. Measured: `go build` rc=0,
  `go vet` rc=0, gate **rc=0 PASS**, `walked=40` — the file was scanned and yielded zero findings.
  Instance 1, iter-69's `M7`: a guard counting the substring `sudo install … /usr/local/bin/z3`
  **survived** a redirect to `/usr/local/bin/z3x`, because *a prefix-shaped needle cannot detect a
  suffix-shaped mutation*. Same class, different recogniser: in both, every branch was exercised
  and the input space was not. **The instrument, and it is cheap because it is a listing exercise
  rather than a run:** before closing a detector milestone, write down every way the forbidden
  thing can be SPELLED in the target language — called / taken as a value / passed as an argument /
  stored in a field / embedded / aliased / dot-imported / reached through an interface or a type
  alias — and require a control per shape, not per branch. The shapes you can name and reject are
  cheap; the ones you cannot are the declared residual, and a **declared** gap is worth far more
  than an assumed one. Note the asymmetry that makes this worth a rule: the branch sweep and the
  shape sweep look identical in a green suite, and the branch sweep is the one the rulebook asks
  for. **Corollary earned the same hour: `go build ./...` DOES NOT COMPILE `_test.go` AT ALL**, so
  on a test-only milestone the mutation-BUILDS assertion is satisfied vacuously — the compile gate
  there is `go vet`. **PROPOSED to V1 as a shared-skill rule** (World cannot edit the
  mission-control SKILL.md — it lives in the V1 checkout); the gap sits in the rule-3j family, both
  instances are first-party here, and it is mission- and language-independent.
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

## Human Decision Ledger (authoritative current state)

This marked table—not STATUS prose, the legacy OD register, or the rolling GitHub thread—is the
source of truth for which decisions are open. Validate it with
`scripts/mission_decisions.sh --check`; generate the human asks with
`scripts/mission_decisions.sh --open`. Rows and IDs are append-only. A human answer changes the
row to `RESOLVED` in the same iteration that consumes the directive, before advancing the
watermark. Never infer resolution from related code landing.

<!-- decision-ledger:start -->
| ID | Status | Decision / recorded answer | Evidence |
|---|---|---|---|
| D-WORLD-5 | RESOLVED | **A — import upstream `serveapi` pinned at v0.33.1.** This is Decision 2's own version discipline executing as written ("P6.B pins the first released/tagged upstream revision that contains this seam and records its commit in `go.mod`/`go.sum`"); B was never the doc's design — Decision 1 names the local codec "the forbidden reinvention" and the milestone section forbids it as schedule compensation. The revision round must (i) run P6.A's frozen conformance fixture against v0.33.1, (ii) audit the dependency closure against `TestDaemonDependencyAllowlist` BEFORE any `go.mod` change — a disallowed graph routes to the doc's Open Decision 4 default (ask upstream for a protocol-only module), never a broad relaxation. The two version axes do not interact: `serveapi` is hash-pinned Go protocol code; the v0.30.0 pin is the `.ail` compiler. The verified commit-boundary contract stays a P6.B prerequisite (`world/*.ail` work). | Ratified attended by Mark 2026-08-17 (recommendation adopted verbatim). Consumes iteration 88's four-for-four `serveapi` measurement. |
| D-WORLD-17 | RESOLVED | **A — bind every seal to its minting validator; tranche 1 stays library-only and explicitly NON-PRODUCTION.** `GradeOfValidated` becomes a METHOD on `Validator`, refusing seals minted by any other instance — this answers the round-4 catch (public `NewValidator` + a free `GradeOfValidated` lets any Go caller mint `PROVEN`) by making possession of a seal worthless without the minting validator. Self-minting into your own validator remains possible and is accepted: no library can stop a caller lying to itself; the enforced property is that you cannot make SOMEONE ELSE'S validator resolve your seal. The revision MUST add the cross-validator refusal arm as a named RED mutation (seal minted by validator 2, refused by validator 1) — that arm is what makes "bind" non-vacuous. Production key custody stays with `w-proven-evidence-production-key-wiring`. B was rejected for chaining the item behind an unwritten successor and re-creating the guard-lands-before-the-guarded-thing vacuity pattern (AC12/AC13 precedent). | Ratified attended by Mark 2026-08-17 (recommendation adopted verbatim). |
| D-WORLD-18 | RESOLVED | **A — ship the scoped item as designed**: the 7 daemon GET routes, the five zero-cost fold-ins, and `TestNoNewDeadlineFreeStoreReads` pinning the 11 residual deadline-free sites (approve 8, registry 2, replay 1; follow-on progress mechanically observable 11 → 0). Per the doc's §13, A unparks this doc DIRECTLY to sprint-planner as written — no designer round. The reviewer's store-boundary guard stays the follow-on's declared closing move, landable exactly when the ratchet reads zero; the two policy questions (what bounds an attended approval *designed* to wait on a human; `Commit` atomicity vs deadline) travel with the follow-on and return to Mark before it routes. B was rejected because it forces those two policy questions NOW, costs ≥2.5d vs 1.5d, and blocks item 14 a further day — while A's residue is enforced, not hidden. | Ratified attended by Mark 2026-08-17 (recommendation adopted verbatim). |
| D-WORLD-ROUTE-1 | RESOLVED | Controller and Anthropic-required planner routes fall back to Codex Sol when Anthropic is unavailable; executor remains Codex Sol primary, DeepSeek v4 Flash second, and Opus last. | Fleet routing directive propagated from AILANG commit `de0e41099` on 2026-08-15. |
| D-WORLD-DRIVER-1 | RESOLVED | **B, with teeth — the driver stays FLEET-owned, and the sync is a COMMIT, never working-tree dirt.** Updates to `tools/launchd/*` land in this repo only as fleet-authored commits (attended, or pushed from the fleet side) — World's controller still never edits them, so the frozen-core rule survives unamended in substance. Strict B-as-dirt was rejected on iteration 89's own measurement: launchd executes the working tree, so an uncommitted driver is a live driver in no repository (it lurked two days and carried this very ledger), and the bundle's CI step cannot run on a checkout lacking untracked files. Enforcement: `verify_go.sh`'s DRIVER DRIFT GATE reds while `tools/launchd/` or `scripts/mission_decisions.sh` diverges from HEAD (path-liveness control ≥5 tracked files); a drift-gate red means "the fleet must commit", never "absorb it". The 2026-08-15 bundle (from ailang `de0e41099`) is committed as the first fleet-authored commit, immediately preceding this ledger update. | Ratified attended by Mark 2026-08-17 (recommendation adopted: B in fleet-authored-commit form). A was rejected because the driver is deliberately mission-agnostic (`MISSION_PROFILE`) — a World-owned fork makes accidental drift structural. |
| D-WORLD-19 | RESOLVED |  **ANSWERED 2026-08-19 (Mark, attended, #68 comment `D world 19 - A yes`): A — YES, tranche 1 MAY extend `host/store` with a bounded object read.** Adopt `gemini-3-1-pro`'s round-6 fix verbatim: a bounded read API on the store (`OpenObject(ref) (io.ReadCloser, error)` or a `maxBytes` bound), §3.3 step 2 streams through an `io.LimitReader` at 256 KiB so the cap runs BEFORE the first full copy is materialised, and §8.2's frozen package list widens to include `host/store`. Round 6's uncontested 6a fix ships with it (sum-style `ResolutionResult`, zero value = refusal — `Validator.Resolve` gains the refusal channel every criterion already required). D-WORLD-17 arm A stays intact. Original ask follows. **May tranche 1 of `w-validated-proven-evidence-boundary` extend `host/store` with a bounded object read?** A: yes — adopt `gemini-3-1-pro`'s round-6 fix verbatim (`OpenObject(ref) (io.ReadCloser, error)` or a `maxBytes` bound; §3.3 step 2 streams through an `io.LimitReader` at 256 KiB; §8.2's frozen package list widens to include `host/store`), closing the OOM vector inside this tranche but putting a second item into `host/store` while item 18 is queued to bound that same package. B: no — record the unbounded first copy as an explicit named limitation of a tranche already declared library-only and NON-PRODUCTION, with an AC stating the residual and a named successor (item 18 or the tranche-2 wiring item) owning the bounded read. Both arms keep D-WORLD-17 arm A intact and both include round 6's uncontested API-shape fix (sum-style `ResolutionResult`, zero value = refusal). | Iteration 90, quorum round 6 blocked with both reviewers present (`metered=$0.174209`). The premise was controller-measured at `03c7892`, not forwarded on trust: `host/store` exposes exactly two exported Object methods (`PutObject:443`, `GetObject:467`), `GetObject` returns the full payload, and non-test `host/store` contains ZERO `io.Reader`/`io.LimitReader`/`maxBytes` occurrences (control: 23 exported `Store` methods, so the zero is a measurement). **RESOLUTION EVIDENCE:** `scripts/mission_directives.sh --issue 68 --since 2026-08-18T19:55:33Z --repo sunholo-data/ailang-world` → 1 directive from `MarkEdmondson1234` @ `2026-08-19T04:57:30Z`, body line 1 `D world 19 - A yes` (of 9 comments; the allowlist filter is enforced in the script). |
| D-WORLD-20 | RESOLVED |  **ANSWERED 2026-08-19 (Mark, attended, #68 comment `Remove deepseek flash`): A — SUSPEND IT.** The ratified executor chain becomes **codex → opus**; `pi:openrouter/deepseek/deepseek-v4-flash-0731:floor` is removed from `MISSION_EXECUTOR_MODEL`/`MISSION_EXECUTOR_FALLBACK` and from the charter's routing policy. This row SUPERSEDES the middle link of `D-WORLD-ROUTE-1` ("DeepSeek v4 Flash second") and nothing else in it. DeepSeek returns only if a future attended ratification re-qualifies the lane. Original ask follows. **Does `pi:openrouter/deepseek/deepseek-v4-flash-0731` stay in the ratified executor chain?** A: SUSPEND it — the chain becomes codex → opus, and deepseek returns only after the lane is re-qualified by a run that changes bytes. B: KEEP it as ratified — the loop continues to spend one failed run per iteration on it, which is cheap in dollars (~$0.02) but costs ~20 minutes of slot and has twice now required a disk-ceiling kill. This is Mark's call and not the controller's because the chain ("we want codex, deepseek, opus") is an attended ratification of 2026-08-10; the ≥3-datapoint evidence bar for a routing change is MET but a controller may not overturn a direction the human set. | Iteration 93: the lane is now **4-for-4 zero-byte failures across four iterations (91 ×2, 92, 93) by four distinct mechanism readings**, with the 1-token probe returning rc=0 every time — so the probe is blind to this failure class by construction. This iteration's run: the model read the code and reasoned CORRECTLY about all six target sites (visible in the NDJSON tail), made 6 tool executions, then streamed pure `thinking` until it wrote **329,584,585 bytes** and was killed by the controller's 300 MB size ceiling; `stopReason="length"` count **0**, `agent_end` **0**, `turn_start` 9 > `turn_end` 8, worktree diff EMPTY. **The size poll was again the ONLY detector that fired** — second consecutive iteration — and the `stopReason` assertion the shared skill prescribes has now fired on **0 of 4**. **RESOLUTION EVIDENCE:** same directive, body line 2 `Remove deepseek flash` (`MarkEdmondson1234` @ `2026-08-19T04:57:30Z` on #68). Applied at iteration 95 to `~/.config/ailang/mission-world.env` and the charter's routing-policy block. |
| D-WORLD-21 | RESOLVED | **ARM A — `ReadObject(ctx, ref, maxBytes)`** (Mark, attended 2026-08-19). The store enforces `maxBytes` BEFORE materialization and performs the complete read under the supplied context, so cancellation becomes ENFORCEABLE and the streaming reader is RETIRED. Chosen because it is the arm already consistent with the rulings on record: it is what objection 6b demanded, what `D-WORLD-19` arm A ratified, and what `gemini-3-1-pro` has PASSED — whereas arm B's burden is naming a concrete termination mechanism in `modernc.org/sqlite` that has never been shown to exist. `gpt5-6-sol`'s objection is UPHELD on its merits: `io.ReadCloser.Read` takes no context and is under no obligation to unblock on `ctx.Done()`, so a `context.WithTimeout` bounds the OPEN call only. | Un-parks item 17 after FIVE consecutive foreclosures of the narrow-refinement carve-out. **Owed unconditionally under this arm** (it was owed under either): `gpt5-6-sol`'s real-store integration test — a blocked read under `context.Background()` relying ONLY on `ObjectReadTimeout`, mutating the ACTUAL cancellation mechanism rather than `WithTimeout`. **Mutation M22 must be REWRITTEN, not re-run**: it is vacuous by construction because its prescribed fake is wired to observe the context, so the mutant dies for a property the real `host/store` reader has never been shown to have — the fake-audit gap raised as Gate-5 skill Proposal 2 (iteration 95). Retiring the streaming reader means item 17 rebases onto item 18's `GetObject` signature change either way. |
| D-WORLD-22 | RESOLVED | **ARM B — RESOLVED 2026-08-20 as a CONSEQUENCE of `D-WORLD-23` arm A (Mark, attended, `#68` comment `A`), not by a separate ask.** Tranche 1 KEEPS ITS SCOPE and weakens the CLAIM to exactly what it proves: every wait item 17's own code performs is bounded by `ObjectReadTimeout`; a LOCK-contended wait is bounded by `busy_timeout` (2000 ms, `writer_lock.go:179`); and the composition is safe ONLY WHILE `busy_timeout` < `ObjectReadTimeout`. Per `D-WORLD-23` obligation (ii) that ordering — which NOTHING in the tree asserts today — ships with an assertion PINNING it; per obligation (i) the residual keeps its named owner, **queue row 22 `w-daemon-lock-wait-not-deadline-bound`**, which must be verified OPEN by command at the moment the revision is written. `gpt5-6-sol`'s objection is SUSTAINED as substantively correct and is answered by narrowing the claim, never by disputing it; arm A is rejected because it absorbs a separately-owned row against Standing rule 1 and re-prices a tranche already 0.75 d over the 3–4 d guardrail. **Original ask follows.** **ONE WORD — does tranche 1 absorb queue row 22's lock-wait bound, or does the CLAIM weaken to match what is proven?** Round 10 of item 17: `gemini-3-1-pro` PASSED (2nd consecutive); `gpt5-6-sol` rejected on the residual the document DECLARED — a LOCK-contended `ReadObject` returns via `busy_timeout` (2000 ms, `writer_lock.go:179`), not via `ObjectReadTimeout`, measured at iteration 94 as **2.043 s under a 300 ms deadline**. The objection is substantively CORRECT and undisputed by the controller. **A** = widen tranche 1, applying `gpt5-6-sol`'s fix verbatim: every wait in `ReadObject` including SQLite busy/lock retry capped by the earlier of the caller's context deadline and the configured store bound; a real-store test holding a conflicting lock; a mutation restoring the fixed 2000 ms wait that must red; the residual deferral deleted. Makes `D-WORLD-21`'s "cancellation is ENFORCEABLE" true without qualification — and absorbs queue row 22 into item 17 against standing rule 1, re-pricing a tranche already half a day over the guardrail. **B** = keep the tranche as ruled and weaken the CLAIM to exactly what is proven: every wait this tranche's own code performs is bounded by `ObjectReadTimeout`; a lock-contended wait is bounded by `busy_timeout`; the composition is safe only while `busy_timeout` (2 s) `<` `ObjectReadTimeout` — plus an assertion PINNING that ordering, which nothing in the tree asserts today. Row 22 keeps ownership. Costs the strength of the claim, buys the ordering guarantee the codebase lacks. | Resolved 2026-08-20 by the ratified `D-WORLD-23` standing rule, whose arm-A text names this row's arm B as its own first consequence (`Resolves D-WORLD-22 as arm B and this instance identically`). Open since iteration 96 (2026-08-19); NOT re-asked — a resolved row is never asked again (Gate 0 decision-recording contract). Item 17 unparks; the owed revision is the claim-narrowing plus the `busy_timeout < ObjectReadTimeout` ordering assertion, and it is a SEPARATE queue item from row 20 under Standing rule 1. **Prior evidence follows.** Raised iteration 96 (2026-08-19), doc §10.12, quorum artifact `w-validated-proven-evidence-boundary-2026-08-19T12-27-17Z.json` (`absent_reviewers` empty, both slots present, `metered=$0.334800`). **Why the controller may not settle it:** `gpt5-6-sol`'s `proposed_fix` says *"Make deadline enforcement for lock contention part of this tranche … Remove the residual deferral"* — a SCOPE call folding a separately filed and separately owned queue row into item 17, and simultaneously a challenge to a predicate of the ruling on record (`D-WORLD-21` chose arm A on the ground that under it *"cancellation becomes ENFORCEABLE"*). The narrow-refinement carve-out is foreclosed on the scope axis for the SIXTH consecutive time. **Owed under BOTH arms:** `gemini-3-1-pro`'s round-10 non-blocking note — `DecodeProposal` caps raw bytes at 256 KiB but states no maximum JSON NESTING DEPTH, so a 256 KiB payload of `[[[[…` can burn CPU or lean implicitly on Go's internal 10,000-deep limit. §10.6 records that an unforwarded non-blocking note returns as BLOCKING one round later; apply it in the SAME revision that applies this ruling. |
| D-WORLD-23 | RESOLVED | **ARM A — RATIFIED BY MARK, ATTENDED, 2026-08-20T08:01:31Z (`#68`, one-word comment `A`). STANDING RULE, EFFECTIVE IMMEDIATELY: when a quorum objection's `proposed_fix` would fold separately-owned work into the tranche in front of it, the tranche ALWAYS KEEPS SCOPE, WEAKENS ITS CLAIM TO EXACTLY WHAT IT PROVES, AND RECORDS THE RESIDUAL WITH A NAMED OWNER — and the controller APPLIES THAT WITHOUT ASKING.** Three obligations bind the application, each an already-earned mission lesson rather than new invention: **(i)** the residual's named owner MUST be an OPEN queue row, asserted by command at the moment of writing — a LANDED row cannot own follow-on work (queue row 20's own §7/§8 defect, and the class queue row 23 exists to record); **(ii)** the weakened claim must STATE the composition condition it now depends on and PIN that condition with an assertion wherever nothing in the tree asserts it today — a claim narrowed without a pin is a claim merely made quieter; **(iii)** the rule licenses a SCOPE call ONLY — it never overrides a PREMISE objection (rule 3f: measure it first-party) nor a dispute about the design DIRECTION (Standing rule 2: that still parks). **This resolves `D-WORLD-22` as its arm B, and queue row 20's round-2 `gpt5-6-sol` objection identically, and every future instance without a further ask.** Foreclosed on this axis 7x before ratification. **Original ask follows.** **ONE WORD, AND IT SETTLES A CLASS RATHER THAN AN ITEM — when a quorum objection's `proposed_fix` would fold separately-owned work into the tranche in front of it, does the tranche ALWAYS keep scope and weaken its claim to exactly what it proves (**A**), or must the controller escalate every instance as today (**B**)?** This is the SECOND live instance of one shape in five iterations. `D-WORLD-22` is the first (item 17 vs queue row 22's lock-wait bound); this is the second (queue row 20 vs the `cmd.WaitDelay`/bounded-cleanup boundary, now queue row 24). In both, `gpt5-6-sol` is substantively CORRECT and undisputed by the controller, its fix is concrete and reviewer-authored, it does NOT dispute the design DIRECTION — and the narrow-refinement carve-out is foreclosed anyway because the fix is a SCOPE call. The carve-out has now been foreclosed on the scope axis **seven** times. **A** = a standing rule: keep scope, weaken the claim to what is proven, record the residual with a named owner, and let the controller apply that WITHOUT asking. Resolves `D-WORLD-22` as arm B and this instance identically, and every future one. Costs the strength of tranche claims; buys the loop back roughly one parked iteration per occurrence. **B** = keep escalating per item, because a scope call is inherently per-item judgment. Costs a human round-trip each time. | Ratified attended by Mark 2026-08-20T08:01:31Z on `#68` (allowlisted author, via `mission_directives.sh --issue 68 --since 2026-08-19T04:57:30Z` — 1 directive of 16 comments). **The bare `A` binds to `D-WORLD-23`, not to `D-WORLD-22`:** iteration 99's report listed `D-WORLD-23` FIRST and as the `(new, one word)` ask, and listed `D-WORLD-22` beneath it as `unchanged; A above resolves it as arm B`, so the single token has exactly one reading that answers the question posed. The competing reading (`D-WORLD-22` arm A) leaves the new ask unanswered AND selects the arm that same report flagged as breaching Standing rule 1. The interpretation is surfaced in iteration 100's report so a misread costs one word, not one sprint. Applied the same iteration to queue row 20 (unparked, routed) and queue row 17 (unparked, revision owed). **Prior evidence follows.** Raised iteration 99 (2026-08-20). Quorum artifacts `w-capsule-output-cap-load-flake-2026-08-20T04-52-19Z.json` (round 1) and `…T04-58-53Z.json` (round 2), `absent_reviewers` EMPTY in both — no self-selecting hole, so the block is the reviewers' considered position and not a degraded quorum. **Why the controller may not settle it:** iteration 96 explicitly RULED OUT *"weakening the doc's claim myself to dissolve the objection (that IS arm B, and choosing it is the human's call)"*, so applying arm A unilaterally here would overturn a recorded ruling on the identical axis. **The empirical half is already measured, so the answer needs no investigation:** production `capsule.Run` IS bounded — 3.005 s against a 3 s `ExecTimeout` with a grandchild holding stdout, correct error kind, `overran=false`, control 11.29 ms — so `gpt5-6-sol`'s "may wait indefinitely" is REFUTED for `Run` as configured. What STANDS is the narrower half: the bound is supplied by the CALLER, the seam being extracted does not require it, and AC5's caller-released fakes prove the protocol rather than the production bound. |
| D-WORLD-24 | RESOLVED | **ARM A — RESOLVED 2026-08-20T16:04:52Z, attended (Mark, `#68`, bare comment `A`). THE PRODUCER IS SHED into ordered queue row 26 `w-bounded-z3-report-producer`; tranche 1 becomes the validator, seal, boundary and their gates. Applied in full by the round-13 designer revision (`codex:gpt-5.6-sol`): §3.4 replaced by a shed marker, `NewProducer`/`Producer.GenerateProof`/`ObjectWriter` removed from the live §3.2 surface, AC5/AC20 and M15/M28 removed as DECLARED GAPS with no renumbering, AC9's required-mutation enumeration updated, §4/§5 claims weakened to hand-authored-fixture validation only. Both round-12 objections leave UNFIXED with the producer, verbatim, in row 26, together with the AC20 connection-pool-not-lock-wait obligation and the missing-reader gap. The ORIGINAL QUESTION follows. ~~ONE WORD — does tranche 1 of item 17 SHED the producer (§3.4) into its own ordered document (**A**), or KEEP it and apply round 12's producer-side bounded-wait fixes (**B**)?** This is the MIRROR of `D-WORLD-22`/`D-WORLD-23`, not a duplicate: those settled whether a tranche ABSORBS separately-owned work (no — keep scope, narrow the claim, name the residual's owner). This asks whether it SHEDS work it DOES own — a scope CUT of a named core deliverable, which no ratified rule licenses the controller to make, and `D-WORLD-23` obligation (iii) is explicit that the standing rule is a scope LICENCE only. **Why it is live now and was not at iteration 100:** §10.13 examined decomposition and rejected shedding the producer on the arithmetic then available (~0.80 d off 5.15 d leaves 4.35 d — still over the 3–4 d guardrail, so it bought no compliance). Round 12's two objections are BOTH producer-side, so shedding now dissolves both unfixed rather than costing another revision round to fix; that arithmetic did not exist when §10.13 ran. **A** = shed the producer; tranche 1 becomes the validator/seal/boundary and their gates. Costs the end-to-end demonstration (no in-repo path from an `ai-check` run to an authenticated envelope; every validator fixture hand-authored, as AC13 already is) and deletes `NewProducer`/`Producer.GenerateProof` — the API surface `gemini-3-1-pro`'s OWN round-9 fix forced into §3.2 three rounds ago. **B** = keep it and fix it: `NewProducer` gains the missing reader parameter and its own derived deadline (an `ObjectReadTimeout` twin, so the producer stops trusting the caller's context), and the write's lock-contended residual weakens to the same composition claim the read side already carries under `D-WORLD-22` arm B, owned by OPEN queue row 22. Costs another revision + re-quorum and more days on a tranche already 1.35 d over guardrail — and on this document's record carries a real chance the fix reveals the next surface.~~ | Raised iteration 101 (2026-08-20). Quorum artifacts `w-validated-proven-evidence-boundary-2026-08-20T15-03-08Z.json` (round 11) and `…T15-09-51Z.json` (round 12), `absent_reviewers` EMPTY in both — neither block is a degraded quorum. **THE EVIDENCE IS A PATTERN, NOT AN OBJECTION:** for THREE consecutive rounds the blocking objection has landed on the PREVIOUS round's own fix, each time one surface over (round 11's live-`PRAGMA` pin was itself an unbounded pool wait; round 11b's `WriteObject(ctx,o)` itself carries an unpinned lock-contended wait). Mechanism, measured: item 18 threaded the store's READ signatures and left the waits unbounded (DR-2, 11 sites, V41), `PutObject` was never threaded at all (V51), and the busy window is ordered against nothing (V49) — so every store surface the tranche newly touches arrives with the same three holes and bounding one reveals the next. **Revision budget spent honestly** — one designer revision (round 11), one quorum, one bounded carve-out revision applying both reviewers' verbatim fixes (round 11b), one confirming re-quorum (round 12); a third revision would be unbounded re-litigation. **Owed under BOTH arms:** AC20's decoy arm must state that it exercises the connection-POOL wait and not the LOCK wait (`gemini-3-1-pro`'s vacuity finding is true of the criterion as written); and `gpt5-6-sol`'s reader-parameter gap holds under B and disappears with the producer under A. | Ratified attended by Mark 2026-08-20T16:04:52Z on `#68` (bare `A`); consumed by iteration 102, which wrote both watermarks after triage. Round-13 quorum artifact `w-validated-proven-evidence-boundary-2026-08-20T19-41-49Z.json`.
| D-WORLD-25 | RESOLVED | **ARM B — "finish 14". ANSWERED 2026-08-24T23:14:21Z (Mark, attended, `#89`, verbatim two-word comment `Finish 14`).** Item 14 `w-workbench-read-only` completes `WB.I`/`WB.J`/`WB.K` before charter row 5 `w-mcp-projection` is taken. Standing rule 1's one-item-at-a-time discipline and the charter's recorded NEXT both hold; row 5 stays UNBLOCKED and ordered immediately after item 14 closes, and its toolchain precondition (`go 1.26.6` vs the `GOTOOLCHAIN: go1.25.6` CI pin) travels with it as row 5's first milestone rather than as a decision. Applied in the same iteration that consumed the directive, before the watermark advanced. **Original ask follows.** ~~ONE WORD — now that row 5 is UNBLOCKED, does it PREEMPT the in-flight item-14 workbench sprint, or does item 14 finish first?** Iteration 120 verified by measurement that AILANG **v0.33.2** delivers [`ailang#764`](https://github.com/sunholo-data/ailang/issues/764)'s protocol-only module, so charter row 5 `w-mcp-projection` is no longer blocked. Row 5 is the **sole remaining blocker on M4** — the reference-agent integration that DESIGN.md names as the value gate (`item 5 -> item 6 -> M4`), i.e. the first propose/verify/commit by a real agent. Item 14 `w-workbench-read-only` is `[IN-SPRINT]` at **8 of 11** milestones with `WB.I`/`WB.J`/`WB.K` remaining, all controller-work. **A = "row 5"**: preempt item 14 and take row 5 next, buying the M4 critical path at the cost of leaving a sprint three milestones short (its landed work is all evidence/test-only and stands on its own; nothing rots). **B = "finish 14"**: complete `WB.I/J/K` first, honouring Standing rule 1's one-item-at-a-time discipline and the charter's recorded NEXT, at the cost of ~3 iterations before M4's blocker starts moving. **Why the controller may not settle it:** this is a queue-ORDERING call against a mid-flight sprint, and the charter's ordering is ratified mission state — no standing rule licenses the controller to reorder around an in-sprint item. Either answer is one word.~~ | **RESOLUTION EVIDENCE:** `/Users/voightkampff/dev/sunholo-data/ailang/scripts/mission_directives.sh --issue 89 --repo sunholo-data/ailang-world --since 2026-08-24T20:06:56Z` -> **1** directive from `MarkEdmondson1234` @ `2026-08-24T23:14:21Z`, body `Finish 14` (of 9 comments; the author allowlist is enforced in the script, never in prose). Control fires: the same script over `#68` since epoch returns **3** genuine directives, so the single hit is a measurement. Watermark read as the OLDER of the two files per the Repo Profile (both `2026-08-24T20:06:56Z`). **Prior evidence follows.** Raised iteration 120 (2026-08-24) on Mark's own `#89` directive *"AILANG has done a new release please verify it's what you need to unblock"*. The unblock is MEASURED, not assumed: `serveapi/protocol` present in tag `63e7909f` (v0.33.2, one tag earlier than upstream's recommended v0.34.0); World's own `TestDaemonDependencyAllowlist` reds naming **exactly one** intruder on a probe-import and greens on **one narrow allowlist line**, 244 -> 250 packages, +1 non-stdlib (against **476** disallowed across 86 roots for the old `serveapi` seam). One precondition travels with it and is row 5's first milestone, not a decision: v0.33.2 declares `go 1.26.6` while CI pins `GOTOOLCHAIN: go1.25.6`, and the repo's own canary clears the move (`go1.25.6` rc=0, `go1.26.5` **rc=1** miscompile, `go1.26.6` rc=0; full `verify_go.sh` under go1.26.6 rc=0). |
| D-WORLD-26 | RESOLVED | **ARM A — `Authorization: Bearer <session-credential>`. ANSWERED 2026-08-25T19:06:41Z (Mark, attended, `#89`, verbatim one-character comment `A`).** The MCP/A2A projection's session credential rides the standard `Authorization: Bearer` header; `protocol.SessionResolver.ResolveSession(ctx, *http.Request)` reads it there. This was the designer's and the controller's RECOMMENDED default: protocol-native for MCP and A2A clients, nothing bespoke to document, and no new header to teach. Arm B (`X-World-Session`) is REJECTED. Two constraints travel with the answer and are NOT relaxed by it: (i) the static `serve-api` API key remains forbidden as a session model (iter-24 measured it process-wide, so it cannot represent a session), so a `Bearer` value here is a SESSION credential and never an API key; and (ii) clause 3 still binds — the resolver must fail closed on an absent, malformed or unknown credential rather than degrading to an unauthenticated surface. Closes `gemini-3-1-pro`'s round-3 blocking objection ("an executable design must close its own forks") on the `w-mcp-projection` doc. Gates only `P6.B`; `P6.T`, `P6.D` and `P6.V` were never blocked on it. | Directive: `scripts/mission_directives.sh --issue 89 --since 2026-08-25T16:57:52Z` -> `MarkEdmondson1234 @ 2026-08-25T19:06:41Z: A` (allowlist enforced in-script). Consumed iteration 125, 2026-08-25, before the watermark advanced. |
| D-WORLD-27 | RESOLVED | **The `D-WORLD-20` deepseek SUSPENSION IS SUPERSEDED FLEET-WIDE — recorded here because this mission's charter still asserted the suspension as current state, and the ledger is the authoritative record.** `D-WORLD-20` (Mark, attended, 2026-08-19, `#68`, verbatim `Remove deepseek flash`) suspended `pi:openrouter/deepseek/deepseek-v4-flash-0731` from the executor chain and said it returns only on a NEW attended ratification. That ratification is recorded in the shared `mission-control` SKILL as a **PROMOTION RULE (Mark, attended 2026-08-26 — supersedes the `D-WORLD-20` suspension)**: DeepSeek returns as the **fallback link only, NOT a rotation peer**, reached only when codex is dry, and re-qualified on runs — two consecutive real sprint executions returning verdict `ok` with a non-empty worktree diff promote it; a single non-zero verdict resets the count. Rationale recorded there: the five failures behind the suspension were all measured through instrumentation now known to be broken, so they are not evidence about the model. **Dated evidence:** the shared skill is `cmp`-identical to `origin/dev` at this iteration; the rule landed in `sunholo-data/ailang` commit `b4122db56` (2026-08-26 07:47); and `~/.config/ailang/mission-world.env` (mtime 2026-08-26 07:46) has its opus override commented out with a `POSTURE`/`ROLLBACK` note matching the rule. **Not inferred from code landing** — the supersession is stated in the ratified shared instructions this loop runs from. **Scope note, and it is not a decision:** no `MarkEdmondson1234` comment on the World channel (`#89`) carries this; it reached World through the shared skill and the live env. If that reading is wrong, one word on `#89` reverses it. **Caveat that IS actionable and is filed as queue row 54, not here:** the id actually in force on this rig is the `:floor` variant that fleet commit `fc6e42682` explicitly dropped, because World's driver copy predates it — so a re-qualification run today would re-run the broken experiment. | Shared `mission-control` SKILL.md PROMOTION RULE; `sunholo-data/ailang@b4122db56`; live `mission-world.env` mtime 2026-08-26 07:46; this session's exported `MISSION_EXECUTOR_FALLBACK`. |
| D-WORLD-28 | RESOLVED | **ONE WORD — how should `verify_go.sh` guarantee its nested race-control module can execute?** **A (recommended): fail closed unless the selected `ACTIVE_GO` is at-or-above the root module floor, then bind the race-control invocation to that exact `ACTIVE_GO`; statically require `racecontrol/go.mod` floor <= root floor.** **B: keep ambient nested-module auto-selection and pin `racecontrol/go.mod` to an explicit minimum independent of the root floor.** Quorum round 2 found that A still needs the explicit `ACTIVE_GO >= root floor` runtime floor: under `GOTOOLCHAIN=local` with a base below the root floor, `go env GOVERSION` can return that lower base, so the static implication alone is false.  **ANSWERED — A — fail closed unless the selected ACTIVE_GO is at-or-above the root module floor, then bind the race-control invocation to that exact ACTIVE_GO, with the static racecontrol/go.mod floor <= root floor requirement. Include the runtime floor check quorum round 2 found necessary: under GOTOOLCHAIN=local with a base below the root floor, go env GOVERSION can return that lower base, so the static implication alone is false.** (Mark Edmondson, attended 2026-09-01, recorded directly in this ledger.)| Iteration 137; two full-strength quorum rounds, all three external reviewers present, both BLOCKED. Artifacts `w-racecontrol-floor-bump-disarms-the-race-control-2026-08-28T15-21-12Z.json` and `…T15-27-58Z.json`.  **Attended ruling 2026-09-01** — recorded in-session under the ATTENDED LEDGER EDITS contract, not via the bookkeeping issue; provenance is the commit author `mark@aitanalabs.com`, which the fleet bot does not hold and the loop may not author with. Attended ruling, matching the loop recommendation on all three; each option was already measured first-party by the loop and no reviewer disputed the direction. Unparks rows 48, 50 and 52 for the normal queue order.|
| D-WORLD-29 | RESOLVED | **ONE WORD — after the whitespace-tolerant rewrite of `shellAssignmentValues`, should a single *indented* assignment be ACCEPTED or REJECTED?** **A (recommended): ACCEPT it — one whitespace-tolerant scan; trim leading spaces/tabs, ignore comments, extract `NAME="…"` regardless of indentation, require `len(values) == 1`.** The row's silent arm then yields count=2 and REDS, so the row's whole purpose is served; an indented-only assignment yields count=1 and supplies its value, which is what bash actually does. Consequence: today's loud red on the indented-only shape disappears — and it is measured to fire because the scanner cannot SEE a real assignment, not because anything is wrong. This is queue row 50's own text and `gpt5-6-sol`'s verbatim `proposed_fix`. **B: REJECT it — keep the design doc's two-sided invariant (column-0 count == 1 AND indented count == 0).** Strictly stronger: reds on ANY indentation, including a benign re-indent — but its stated rationale ("an indented assignment is a conditional/nested shadow") is measured FALSE, so B would be chosen for defence-in-depth, not for the reason the doc gives. **Answering also AMENDS queue row 50**, which lists "an indented *only* assignment likewise reads 0 and fatals" among the correct loud behaviours; that sentence carries the same false premise and is ratified text the loop may not edit itself. **Default if unanswered:** the item stays parked and the queue advances past it to rows 51–58.  **ANSWERED — A — ACCEPT the indented assignment. One whitespace-tolerant scan: trim leading spaces/tabs, ignore comments, extract NAME="..." regardless of indentation, require len(values) == 1. B's stated rationale ("an indented assignment is a conditional/nested shadow") is measured FALSE, and defence-in-depth bought with a false premise is not worth a loud red on a benign re-indent. AMEND queue row 50 in the same iteration to drop the sentence carrying that premise — it is ratified text the loop could not edit itself, and this ruling authorises exactly that edit.** (Mark Edmondson, attended 2026-09-01, recorded directly in this ledger.)| Raised iter-140 (2026-08-31) after two full-strength quorum rounds both BLOCKED. Premise measured first-party by the controller: a bash script setting `COL0=`, a space-indented `INDENTED=` and a tab-indented `TABBED=` at top level prints all three, rc=0; paired control — a genuinely conditional assignment DOES change behaviour (`NARROW` unset → `wide list`, set → `narrow`). Quorum artifacts `.ailang/state/mission-quorum/w-shell-assignment-parser-drops-an-indented-assignment-2026-08-31T14-03-36Z.json` and `…T14-16-33Z.json` (that directory is gitignored, so the decision-bearing subset is quoted here); round 2 restored `oc-glm-5-2` → PASS at $0.0203, so `absent_reviewers` is closed and the quorum is full-strength.  **Attended ruling 2026-09-01** — recorded in-session under the ATTENDED LEDGER EDITS contract, not via the bookkeeping issue; provenance is the commit author `mark@aitanalabs.com`, which the fleet bot does not hold and the loop may not author with. Attended ruling, matching the loop recommendation on all three; each option was already measured first-party by the loop and no reviewer disputed the direction. Unparks rows 48, 50 and 52 for the normal queue order.|
| D-WORLD-30 | RESOLVED | **ONE WORD — should the row-52 CI step-scoping fix stay a LINE SCAN, or become a YAML PARSE?** **A (recommended): LINE SCAN, hardened.** Keep the text scanner and add `gemini-3-1-pro`'s verbatim fix — require the derived step column to equal the constant `6` and fatal loudly otherwise — or equivalently derive the anchor from the SHALLOWEST enclosing `steps:` rather than the nearest one. Measured: the shallowest-anchor variant catches the round-2 attack (`stepCol=6`, `start=173`, flag found → RED) on the identical file where the doc's nearest-anchor derivation is blind (`stepCol=12`, both invariants PASS, flag excluded → GREEN). Consequence: no new dependency; the residual is that the scan stays text-level, so `gpt5-6-sol`'s general point — block-scalar content is author-controlled — remains true in principle, and a future re-indent of `ci.yml` makes the gate fatal until the constant is updated. **B: YAML PARSE.** Adopt `gpt5-6-sol`'s verbatim fix — parse `ci.yml` into a node tree, fail loudly on parse errors, aliases, flow-style steps and duplicate keys, locate the step semantically and reject `continue-on-error` on that mapping. Strictly stronger: it closes the CLASS rather than the measured instances. Consequence: adds the **second** direct dependency to a `go.mod` that has exactly **one** (`modernc.org/sqlite`), and widens a ~0.1d row into a dependency decision. **Recommendation: A.** This gate is a repo-internal tripwire against our own accidental regressions, not an adversarial boundary — the threat model that makes B necessary is an actor who can already edit `ci.yml`, and such an actor can simply delete the test. **Default if unanswered:** row 52 stays parked and the queue advances to rows 53–59 then 39. **ANSWERED — A — LINE SCAN, hardened, deriving the anchor from the SHALLOWEST enclosing `steps:` rather than the nearest one (measured to catch the round-2 attack that the doc's nearest-anchor derivation is blind to). This gate is a repo-internal tripwire against our own accidental regressions, not an adversarial boundary: the threat model that would justify B is an actor who can already edit `ci.yml`, and such an actor can simply delete the test. Not worth adding the second direct dependency to a `go.mod` that has exactly one. The residual is accepted and named: the scan stays text-level, and a future re-indent of `ci.yml` makes the gate fatal until the constant is updated.** (Mark Edmondson, attended 2026-09-01, recorded directly in this ledger.) | Raised iter-142 (2026-09-01) after two FULL-STRENGTH quorum rounds both BLOCKED (`.synthesis.absent_reviewers` `[]` in both, cross-checked against `[.reviewers[]\|select(.present==false)]`, control `has("synthesis")` true). Round 1 $0.1035 (reject/reject/pass + controller reject); round 2 $0.1489, all three external reviewers rejecting. The carve-out does not apply because the reviewers' `proposed_fix` fields disagree and `gpt5-6-sol`'s disputes the DIRECTION, naming park-or-widen itself. Both locators were refuted by first-party measurement, not by argument, and the row's own "non-exploitable" claim was refuted in both directions before any routing (ARM A false positive at `ci.yml:166` on an unrelated step; ARM B `rc=0 --- PASS` with the forbidden flag live on the guarded step). Doc banked at `design_docs/planned/w-wiring-test-step-scoping-imprecise-under-key-reorder.md`; quorum artifacts under `.ailang/state/mission-quorum/` dated `2026-09-01T00-50-31Z` and `T01-05-28Z` (that directory is gitignored, so the decision-bearing subset is quoted here). **Attended ruling 2026-09-01** — recorded in-session under the ATTENDED LEDGER EDITS contract, not via the bookkeeping issue; provenance is the commit author `mark@aitanalabs.com`, which the fleet bot does not hold and the loop may not author with. Attended ruling, matching the loop recommendation on all three; each option was already measured first-party by the loop and no reviewer disputed the direction. Unparks rows 48, 50 and 52 for the normal queue order. |
| D-WORLD-31 | OPEN | **ONE WORD — `D-WORLD-29` ratified rule A on a rationale that a measurement taken AFTER the ruling shows does not cover every shape. Ship A as ratified, or hold row 50 for the fixture migration?** `D-WORLD-29` chose A ("ACCEPT the indented assignment … which is what bash actually does"). That rationale is TRUE for a **top-level** indented assignment — measured, and it is what was in front of you. It is **FALSE for an assignment inside a multi-line branch that never executes**, and that shape is not a corner: `run.sh` carries **17** control-flow openers and **65** indented lines (control: a fresh absent keyword reads 0), so there are 17 live positions where a refactorer can put one. **The measurement, first-party at `cb73cab`, mutant LANDED by sha256, `go vet` rc=0 read before every test result, restored byte-identical, pristine control green either side:** with the only `KNOWN_BAD` inside `if false; then` / `fi` and no column-0 copy, the CURRENT helper is **rc=1 `KNOWN_BAD assignment count=0, want 1`** (correctly loud) and the ratified rule-A helper is **rc=0 RUN=2 PASS=2** — while `bash -c 'if false; then KNOWN_BAD="x"; fi; echo "${KNOWN_BAD:-<UNSET>}"'` prints **`<UNSET>`**. So rule A binds the toolchain floor against a deny-list that does not exist at runtime. Paired control, same batch, proving the helper is not merely permissive and that **row 50's own defect really is closed by A**: the row's silent arm (a second, indented `KNOWN_BAD` beside the column-0 one) goes **rc=1 `count=2, want 1`** at both consumers. **This is a defect rule A INTRODUCES, not one it fails to fix** — which is why it is a decision and not a queue row. **A (recommended): SHIP A AS RATIFIED.** Residual 6 is declared in the doc and in the helper comment, pinned by unit arm 11, and recorded as `V-ARMC1`/`V-ARMC2` with its runtime control. Row 50's defect closes; the gate trades one silent hole (a second indented assignment narrowing the list) for another (a never-executed branch supplying it), and the new hole is louder to a reader because it is written down. Cost: ~0.1d, already designed and planned. **B: HOLD row 50 for the fixture migration** — `gpt5-6-sol`'s verbatim `proposed_fix`: move the pins into a declarative, data-only fixture holding exactly one whitespace-tolerant `NAME="…"` record per name, have `run.sh` source it and the Go gates scan it, rejecting non-comment/non-assignment syntax, with mutation arms proving duplicate, `export` and control-flow text all RED. Closes BOTH holes and removes the text-scanner class entirely; costs a new design doc and roughly a day, and it is scope this row never had. Its verbatim words: *"If that migration is out of scope, keep the row blocked rather than accepting the demonstrated V-ARMC2 false green."* **NOT AVAILABLE, and measured so rather than assumed:** a cheap "refuse if the script contains control flow" floor is dead on arrival — it would fatal on the pristine `run.sh` at all 17 openers. A hybrid that also requires column-0 count ≥ 1 is exactly rule B on the one point `D-WORLD-29` decided, so it is not offered. **Default if unanswered:** row 50 stays parked and the queue advances to rows 52, 54–61, then 39; nothing else is blocked by it. | Quorum round 3, full strength — 3 of 3 reviewers PRESENT, `.synthesis.absent_reviewers` `[]` cross-checked by `[.reviewers[]\|select(.present==false)]`, verdict **blocked**, `metered=$0.15101`: `gpt5-6-sol` **reject** (this objection), `oc-glm-5-2` **reject** (AC7's base FAIL set quoted at `9c0ad0b` while the doc targets `cb73cab` — a PREMISE objection, measured and fixed in-loop rather than forwarded), `gemini-3-1-pro` **pass**. Doc `design_docs/planned/w-shell-assignment-parser-drops-an-indented-assignment.md` at 678 lines, revised to direction A this iteration by `pi:ollama/deepseek-v4-flash:0731-cloud` (typed verdict `ok`). `D-WORLD-29` is NOT reopened by this row — it is a new, uniquely named decision carrying a fact that post-dates the ruling, per the shared skill's ATTENDED LEDGER EDITS rule (g). |
<!-- decision-ledger:end -->

---

## STATUS (rotation rule)

## STATUS 2026-09-03 (iteration 151) — **TWO CONSECUTIVE SLOTS DIED, AND THE SECOND ONE IS DATABLE TO THE SECOND: THE RIG REBOOTED **286 s** AFTER ITERATION 150 STAMPED `gate-4`, AND macOS TOOK `/private/tmp` — WITH IT THE TOOLCHAIN PIN THIS MISSION VERIFIES EVERYTHING WITH, THE DRIVER'S ONLY CRASH LOG, AND A SPRINT WORKTREE.** No PR: this iteration's deliverable is **bookkeeping repair + a precondition restore**, which is the disposition Gate 2 mandates for work that is already landed. **GATE 0:** kill switch armed (namespaced path); `gh` = `sunholo-voight-kampff`; billing **CLEAN**; **0** `MarkEdmondson1234` directives on `#107` (14 comments, read from the epoch since both watermarks agree at `2026-09-02T03:56:21Z` — the OLDER-of-two rule), via the V1 checkout's `mission_directives.sh` by ABSOLUTE PATH. Decision ledger **18 rows, `--check` valid, ONE OPEN** (`D-WORLD-31`, re-asked unchanged — **this iteration adds no new ask**); no ledger row changed since the watermark, so **no attended ruling and no self-resolution**. No rotation owed (`#107` created `2026-08-31T09:26:51Z`, after Monday 07:00 local), no weekly sweep owed. **`tools/launchd/mission-heartbeat.sh` is still absent here (row 69) — my `gate-0` stamp was lost to it before I switched to V1's absolute path, which is the row's own defect measured a second time, live.** **GATE 1:** tree carried **four modified mission docs** — not a sibling's work but **iteration 150's own unwritten record**; `dev` == `origin/dev` at entry; running skill **byte-identical to `origin/dev`** (`cmp` against the RESOLVED symlink target). CI **GREEN 3/3** on `origin/dev` HEAD with the parent's `go host build + test gate: failure` as a **firing negative control**. **ONE CROSS-MISSION MESSAGE TRIAGED** (`mission-v1` iter-321): World's Gate-0 ledger proposal **ACCEPTED AS SOUND, not adopted** (V1's arm is a null — its 4 stranded commits did not touch the ledger); its correction that a closure cross-check must be **field-scoped, never a whole-line grep** was taken, and it is what stopped me publishing a census below. Its `git worktree list` warning does **not** bind World (own clone, re-measured: `git-common-dir` = `.git`, remote `ailang-world.git`). **GATE 2 — THE PICK IS THE TWO DEAD SLOTS, FOUND BY THE DIED-MID-FLIGHT TRACES AND NOT BY THE QUEUE.** Trace (b) returned a **`prunable`** `/private/tmp/wt-row56` on `sprint/w-canary-fence-blind-to-a-skipped-canary`; trace (c) returned the four dirty docs. **Slot 1 (unnumbered, 2026-09-02):** row 56 was designed, executed, merged as PR [#113](https://github.com/sunholo-data/ailang-world/pull/113) → [`725ad5a`](https://github.com/sunholo-data/ailang-world/commit/725ad5a) at `19:26:21Z`, and left **ZERO** trace in all four mission documents (`grep -cE '#113\b'` = **0**, **0**, **0**, **0**; known-positive control `#114` = **2** in the charter, **1** in the log). Iteration 150 read that very SHA in the next fire and used it **only as a CI negative control**, never asking what it was. **Slot 2 (iteration 150):** stamped `gate-4` at epoch **`1788394743`**, `kern.boottime` = **`1788395029`** — **286 s later**, compared epoch-to-epoch so no timezone can be wrong. Its record was on disk and uncommitted. **VERIFIED, NOT ADOPTED** (inherited-claim discipline — nobody had reviewed slot 1's work since the agent that wrote it stopped existing): row 56's Gate 3b re-polled SHA-addressed on the MERGE commit, `checks=2 == expected=2`, both `success`; its deliverable confirmed live at `host/verifygate/toolchain_pin_gate_test.go:510`; iteration 150's record cross-checked against my own independent read of `a7b58dd` (3/3 `success`, one file) before I committed it **verbatim**. **THE PRECONDITION NOBODY WOULD HAVE NOTICED:** `/tmp/ailang-v0300/ailang` — the pin the charter names and every gate here runs on — **was gone** (`ls` rc=1). The gate fails CLOSED, which is the design working; the bad half is that the loop is then unable to verify anything and nothing says so. Restored to **`~/.pinned-ailang/ailang`**, the path **CI already uses**, so rig and CI now name the same location and `$HOME` survives a boot. Verified rather than assumed: published `shasum -a 256 -c` **OK**, `--version` = `AILANG v0.30.0` / `e37b370`, and the gate's own strongest control — `verify_ail.sh` step 9/9 `compiler pinned by exact bytes: AILANG v0.30.0 on Darwin/arm64`. **GATE 3 — no sub-agent spawned; controller-authored repair, `metered=$0.00` of $5.** **VERIFY GATE GREEN, BOTH LEGS:** `./scripts/verify_ail.sh` **rc=0** (11 required identities, 40 named tests, 9/9 world-package steps non-vacuous) and `go build ./... && go test ./...` **rc=0** (**19** packages `ok`, **0** FAIL) — which also retires, by measurement, the 18-failure class iteration 150 attributed to `AILANG_BIN` being unset: it was the same wiped pin, one fire earlier. **A CENSUS I DELIBERATELY DID NOT PUBLISH.** The `N of M closed` headline is an increment chain; two attempts to re-derive it read **36 of 72** and **53 of 70** against the carried **58 of 70**, and the second is *demonstrably* over-counting (it classes the open row 57 as closed on prose). The enumeration control fired first and caught that the naive form stopped at row 66. Three numbers, none agreeing — so the honest state is **no verified census exists**, and it is filed as row 72 rather than asserted. **FINDINGS FILED RATHER THAN ABSORBED:** row **71** — the pin, the driver's only crash log and sprint worktrees all lived in `/tmp`; the pin half is closed here, the other two are **fleet-owned** (`LOG=/tmp/ailang-mission-${MISSION_NAME}.log`, `mission-control.sh:76`) and every mission on this rig has the identical exposure; row **72** — the unauditable progress count. **THE LIMIT, RECORDED AS A LIMIT:** slot 1's cause is **UNRECOVERABLE** — the reboot wiped the driver log five hours later, so its `HARD TIMEOUT`/`STALL:` lines no longer exist. Slot 2's is measured. **This is two dead slots in a row and it is reported as a pattern, not as two incidents** (the loop cannot diagnose why its own slots die; it can make the frequency visible). **NEXT:** rows **57**, **58**, **59**, **60**, **61**, **62**–**66**, **68**–**72**, then **39**. Row **50** stays parked on `D-WORLD-31`. **Decision ledger: 18 rows, ONE OPEN — `D-WORLD-31`, re-asked unchanged.**

## STATUS 2026-09-03 (iteration 150) — **ROW 67 LANDED — AND IT IS THE FIRST RED IN THIS REPO THAT NO SPRINT OF MINE WROTE: `dev` WENT RED AT THE MERGE [`68403ea`](https://github.com/sunholo-data/ailang-world/commit/68403ea) BECAUSE TWO INDIVIDUALLY-FINE BRANCHES MET — A THIRD CI JOB WITH NO GO ANYWHERE IN IT, AND A GATE WHOSE ONE `wantJobs` CONSTANT ANSWERED *BOTH* "WHICH JOBS EXIST" AND "HOW MANY GO PINS THERE MUST BE".** PR [#114](https://github.com/sunholo-data/ailang-world/pull/114) → squash [`a7b58dd`](https://github.com/sunholo-data/ailang-world/commit/a7b58dd), **Gate 3b GREEN on the MERGE commit** (SHA-addressed, `present=3 == expected=3`, all three `success`; the PR head was polled to the same standard first, every count asserted numeric before comparison). **GATE 0:** kill switch armed (namespaced path); `gh` = `sunholo-voight-kampff`; billing **CLEAN**; **0** `MarkEdmondson1234` directives on `#107` since the `2026-09-02T03:56:21Z` watermark (14 comments), via the V1 checkout's `mission_directives.sh` by ABSOLUTE PATH — **World's own `scripts/` has no copy** (`ls` rc=2), which is the second fleet script missing here this iteration. Decision ledger **18 rows, `--check` valid, ONE OPEN** (`D-WORLD-31`, unchanged — **this iteration adds no new ask**); no charter commit since the watermark flipped a row and all five recent ones are fleet-bot authored, so **no attended ruling, no self-resolution**. No rotation owed (`#107` created `2026-08-31T09:26:51Z`, AFTER Monday 07:00 **local** = 05:00Z; 14 of 80 comments), no weekly sweep owed. **One cross-mission message triaged** (`mission-v1` iter-321): World's Gate-0 ledger proposal ACCEPTED AS SOUND but not adopted — V1's arm is a null, its 4 stranded commits happened not to touch the ledger. Two things taken from it: the cross-check must be **field-scoped, never a whole-line grep** (V1 measured `grep -c OPEN` = 4 against the script's correct 1 — three RESOLVED rows whose *answer text* contains the word OPEN), and `git worktree list` is **not** proof of PR ownership on a shared clone. **THE SECOND CLAIM DOES NOT BIND WORLD, AND I MEASURED IT RATHER THAN INHERITING IT:** World's checkout is its own clone (`git-common-dir` = `.git`, remote `ailang-world.git`), so worktree presence IS ownership proof here — but the measurement surfaced something worse, filed as row 68. **GATE 1:** `dev` == `origin/dev` at entry, tree clean, running skill **byte-identical to `origin/dev`** (`cmp` against the RESOLVED symlink target, inode-confirmed same file — the stale-checkout gap that has run for 20+ iterations is CLOSED). **CI RED on `origin/dev` HEAD**, `checks=3`, `go host build + test gate: failure` — World owns this repo, so it outranked the queue and became the pick. **GATE 2:** repro first-party at HEAD, `go test ./host/verifygate -run TestGoToolchainPinsAgreeAndMatchJobList` rc=1 with the gate's own message naming both facts; **negative control** parent `725ad5a` CI `success` on both jobs; and the commit that added the third job, `e308577`, carries **ZERO check-runs** — it was never verified on its own, so the merge was the first run that could see the pair. **THE ONE-LINE FIX IS REFUTED BY MEASUREMENT, NOT BY ARGUMENT:** widening `wantJobs` to three entries alone is **rc=1 on the pristine `ci.yml`** (2 pins vs 3 jobs) — adding the job to the list demands a Go pin on a job that must not have one. **GATE 3:** controller-authored direct fix (no designer, no planner, no executor spend — `metered=$0.00`); two hand-maintained lists (`wantJobs`, `wantGoPinnedJobs`) so a new job must be classified in BOTH edits, and a job classified OUT is asserted to carry **zero** Go pins — "classified out" is now a claim the gate checks rather than a hole. Jobs parsed WITH their bodies and pins attributed **per job**, closing a fail-open the count-only gate always had (two `GOTOOLCHAIN` lines in one Go job and none in the other satisfied `count==2` exactly as well as one each), plus a whole-file-vs-per-job-sum arm for a pin outside every job. **MUTATION DRILL, 7 mutants, each LANDED by sha256, `go vet` rc=0, `ci.yml` restored byte-identical, pristine control green — all 7 RED with the correct assertion named. The load-bearing measurement is the revert-just-the-hunk arm** (patched minus per-job attribution, green at base): **M5 SURVIVES rc=0 with no assertion fired**, so the per-job hunk is its SOLE killer and every other mutant is killed by the count arms alone — that gradient is reported, not hidden. **Two-arm verification on the identical tree and env** (`AILANG_BIN` = pinned v0.30.0, `z3` on PATH), rc captured without a pipe: pristine test file **rc=1 with this test as the SOLE failure** (307.8s), patched **rc=0**, whole package green. The 18 other failures seen in a first full-suite pass were `AILANG_BIN` **unset** (fail-closed by design), not this diff — established by re-running with the pin set, which is the difference between an environmental explanation and an environmental excuse. **THREE FINDINGS FILED AS ROWS RATHER THAN ABSORBED:** row 68 — `~/.ailang-driver-pin/world` is a worktree of the **`ailang`** clone holding `v1-mission.md` and **no** `world-mission.md`, i.e. a pin named for this mission pointing at the wrong repo (harmless today only because World's plist runs World's own unpinned driver — measured, `grep pin-root` = 0 hits); row 69 — **`tools/launchd/mission-heartbeat.sh` does not exist in this repo**, so the skill's mandated first action at every gate is `No such file or directory` here and standing rule 7's attribution contract has been inoperative for World unless a controller reaches for V1's absolute path (frozen core → fleet must port it; I stamped all six gates via the V1 path); row 70 — a direct-to-dev commit with zero check-runs is invisible to Gate 1, which reads only HEAD. **NEXT:** rows **57**, **58**, **59**, **60**, **61**, **62**–**66**, **68**, **69**, **70**, then **39**. **Decision ledger: 18 rows, ONE OPEN — `D-WORLD-31`, re-asked unchanged.**

## STATUS 2026-09-02 (iteration 149) — **ROW 55 LANDED: THE DISPATCH-LEVER GATE FALSE-REDDED ON THE STANDARD REMEDY FOR A FAMOUS ACTIONS FOOTGUN — AND THE PLANNER THEN SHOWED THAT `go build ./...` CANNOT SEE A `_test.go` AT ALL, SO EVERY MUTANT-BUILDS FENCE IN A TEST-ONLY SPRINT IS VACUOUS.** PR [#112](https://github.com/sunholo-data/ailang-world/pull/112) → squash [`165b9fd`](https://github.com/sunholo-data/ailang-world/commit/165b9fd), **Gate 3b GREEN on the MERGE commit** (SHA-addressed, `present=2 == expected=2`, both `success`; the PR head was polled to the same standard first). **GATE 0:** kill switch armed (namespaced path); `gh` = `sunholo-voight-kampff`; billing **CLEAN**; **0** `MarkEdmondson1234` directives on `#107` since the `2026-09-01T23:05:53Z` watermark (10 comments), BOTH watermarks read and the OLDER taken (they agree exactly), plus the rotation-week `-prev` read of `#89` (44 comments, 0 directives), all via the V1 checkout's `mission_directives.sh` by ABSOLUTE PATH. Decision ledger **18 rows, `--check` valid, ONE OPEN** (`D-WORLD-31`, unchanged — **this iteration adds no new ask**). Attended-ledger provenance checked per origin's contract: the only charter commit since the watermark is `234d9da`, fleet-bot authored, no row flipped — **no attended ruling, no self-resolution**. No rotation owed (`#107` created after Monday 07:00 local; 10 of 80 comments), no weekly sweep owed (iteration 141 ran it on this rotation). **GATE 1:** `dev` == `origin/dev` at entry, tree clean, CI green 2/2 with the parent commit as control. **The running mission-control skill is now 187 lines behind origin (was 147), 1 removed** — still V1's stale checkout, which World cannot fix; I read the delta and followed ORIGIN's rules where they differ, which is the only reason the attended-ledger check above exists. **PREDICATE RE-RUN, NOT TRANSCRIBED:** the row-54 standing red still stands and is **WORSE** — the World driver copy is now **759** diff-lines behind the fleet, up from 705 yesterday, and `pin-root.sh`/`rig-lock.sh` are absent here entirely. That red means "the fleet must commit". `verify_go.sh` is rc=1 at base for a **second, independent** reason I attributed by MECHANISM rather than by redness: `TestHandlerTimeoutKillsTheWholeProcessGroup` with markers `exec_started=false forked=false` — byte-identical to row 58's recorded signature, so it is that known FLAKE and not a new regression. **GATE 2 — GHOST DISCIPLINE ON AN EVALUATOR-SOURCED ROW.** All three claimed false-reds reproduced first-party at `234d9da` with TWO known-positive controls in the same call (a canonical fixture, and the real `ci.yml` at the path the gate reads): quoted `"on":` and flow style both Fatal `has no top-level 'on:' trigger block`; a TAB-indented first trigger returns `triggers=[]`, misreporting a parse limit as TOTAL ABSENCE. The scalar-arm cascade reproduced too (two messages, not one). The row's own carried claims about the doc's Residual-3 were confirmed **wrong in the safe direction**: `filepath.Glob` DOES return `.hidden.yml`, and a nested subdirectory is a loud `is a directory`. **The row's blocking question was already answered:** `D-WORLD-30` (attended, 2026-09-01) chose LINE SCAN for a sibling gate on rationale that transfers — and I routed it as PRECEDENT, explicitly **not** as resolving row 55, and raised no new decision (standing rule 8: do not manufacture a decision the human does not have). **GATE 3 DESIGN — rotation designer `claude:claude-fable-5`, ONE doc, initial run + ONE protocol-mandated revision, which is the diet's exact ceiling.** Two FULL-STRENGTH quorum rounds (`.synthesis.absent_reviewers` `[]` in both, cross-checked against `[.reviewers[]|select(.present==false)]`, control `has("synthesis")` true), **both BLOCKED, all three reviewers rejecting both times**, $0.0936 + $0.1275. Every objection was a PREMISE objection, so per rule 3f I **measured each rather than forwarding it**: `gpt5-6-sol`'s reuse hypothesis was REFUTED with its own named instrument (`go list -deps ./...` — **0 of 257** packages match `yaml`; controls `modernc` 30, `encoding/json` 1, a fresh absent literal 0) while its SCOPE complaint was CORRECT and widened the conflict surface to a third file (`host/runbook`, excluded by the old `host/verifygate` scope **by construction**); `gemini-3-1-pro` was right on BOTH halves, and its second was a live defect in the doc's own AC6, whose known-positive control `grep -c 'anti-vacuity floor'` reads **0** because the file writes it UPPERCASE (case-insensitively 1; repo-wide 2 vs 22) — **the AC would have failed at its own control**; `oc-glm-5-2` was right that §3.1 cited `D-WORLD-30` with no V-row. Round 2 closed under the **narrow-refinement carve-out** (ratified here at iter-13), conditions checked before use: every remaining objection carried a concrete reviewer-authored `proposed_fix` and NONE disputed the design DIRECTION. **GATE 3 PLAN — `opus`, lane `derive-planner-lane.sh` → `opus fail-closed:env-pin`, used verbatim — REFUTED SIX doc premises, two severe.** The sharpest generalises far past this repo and is the retro finding: **`go build ./...` is not a compile fence for a `_test.go`** — measured, a hard type error in a test file leaves it at **rc=0**, while `go test -count=1 -run '^$' ./pkg/` and `go vet ./pkg/` both red. Every assertion in this sprint lives in a test file. Also: MUT-E as written **redded nothing** (no fixture reached the L128 control), and a naive "replace lines 291-293" would have deleted a TRUE adjacent sentence. **GATE 3 EXECUTION — `codex:gpt-5.6-sol`**, 11 ACs, 5 mutants all LANDED by sha256, all compile-fenced, every observed `--- FAIL` set matching the plan's ENUMERATED set exactly (no `-skip … rc=0` criterion anywhere, because several mutants reach more than one arm). Commits reconstructed per milestone from the executor's snapshots and proven byte-identical to its final tree by `shasum -c`; M1 and M2 collapse into one commit because the drill restores. **GATE 3 EVALUATION — `sonnet`, ITS OWN worktree, 95/100 PASS, ZERO blocking.** It reproduced all five mutants byte-for-byte, ran a precondition-removal drill over the whole arm set, and threw **11 adversarial fixtures** at the gate trying to build a silent false green — and could not. **Its top non-blocking finding was REAL and I closed it, having reproduced it first because a NON-BLOCKING label is a severity opinion and not a measurement:** `onBlockFailureMessage`'s text for the two NEW sentinels was pinned by nothing — replacing the honest message with nonsense left the ENTIRE package rc=0, so **this row's headline deliverable shipped unprotected**. Two arms now cover it and are **proven non-vacuous, not asserted**: MUT-F regresses the message to the exact absent-block claim row 55 removes, and both arms red with the specific regression assertion firing by name. **TWO OF MY OWN LINE NUMBERS WERE REFUTED BY THE DESIGNER** (the stale comment is at `:110-111`, not the `:258-261` I transcribed out of a document that embeds the code; `ANTI-VACUITY` is at `:117`, not `:116`) — rule 3b(v)(b) twice in one iteration, both caught downstream, both recorded rather than waved through. **metered=$0.2211** of $5, all of it quorum. **NEXT:** rows **56**–**66**, then **39**. Row **50** stays parked on `D-WORLD-31`. **New rows 65 and 66** filed from the planner's and judge's out-of-scope findings.

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

- **A FIX APPLIED UNDER THE NARROW-REFINEMENT CARVE-OUT MUST ACQUIRE ITS OWN ACCEPTANCE CRITERION
  AND ITS OWN NAMED MUTATION BEFORE THE DOC ROUTES — THE CARVE-OUT CHANGES THE DESIGN AND CHANGES
  NOTHING THAT GUARDS IT** (process fix, iter-98; **instance 1 — watch-item, the bar for a shared-skill
  proposal is 2**). The shared skill's carve-out is a routing rule: when every remaining blocking
  objection carries a concrete reviewer-authored `proposed_fix` and none disputes the design
  DIRECTION, the controller applies the reviewers' VERBATIM text and routes straight to
  sprint-planner without a third quorum round. That is correct and this mission leans on it heavily
  — seven foreclosures and several successful applications so far. **What no step in it asks is what
  guards the newly-applied text.** The fix enters the document *after* the round that reviewed the
  document, so by construction the last eyes on it are the author's; the AC table, the mutation
  table and the non-vacuity section were all written against the PRE-fix design and nothing
  re-derives them.
  Measured here on row 21. Decision 2's **conditional** self-heal is `gemini-3-1-pro`'s round-2
  `proposed_fix`, applied verbatim — and it shipped with **no acceptance criterion and no mutation**,
  so an *unconditional* heal passed **AC1–AC7 unchanged**. Base state confirmed the gap rather than
  inferring it: `grep -c 'strings\.' host/archive/archive.go` = **0** (control:
  `host/daemon/daemon.go` = **2**). It was the sprint-planner, two roles downstream, that noticed —
  and only because it was re-deriving base status per AC, which is a habit this mission adopted for a
  different reason entirely.
  Note the asymmetry that makes this worth a guardrail rather than a caution: a reviewer's own words
  are **the last text anyone thinks to audit**. They arrive with a quorum's authority attached, the
  controller is explicitly forbidden from re-litigating them, and "applied verbatim" reads as a
  completed obligation. So the one edit in the document that no reviewer has ever seen is also the
  one the process treats as the most reviewed.
  **Rule.** Before routing a carve-out revision: for each applied `proposed_fix`, name the AC that
  fails if the fix is absent and the named mutation that reds if it is neutered — and if neither
  exists, write them in the same commit. Prefer a **dual** where the fix is a conditional: `if false &&`
  neuters a branch that must fire, and cannot neuter a branch that must be **skipped** — row 21 needed
  both `M4a` (`if false && …`) and `M4b` (`if true || …`), and each killed a test the other left green
  (under M4a the zero-execution arm stayed **PASS**; under M4b it was the **only** arm that red). Record
  the pair in the doc's mutation table, not merely in the sprint plan. **PROPOSED to V1 as a
  shared-skill rule at the second instance** — World cannot edit the mission-control SKILL.md, the gap
  is in the carve-out itself rather than in anything mission-local, and one instance is below this
  mission's own bar.

- **A LOAD OR STRESS EXPERIMENT MUST PROVE ITS OWN TEARDOWN BEFORE ANY LATER MEASUREMENT ON THIS
  RIG COUNTS AS EVIDENCE — AND "I RAN `kill`" IS NOT THAT PROOF** (process fix, iter-97; instance 2
  of *the controller's own step contaminated the controller's own control*, after iter-94's
  `MU-SITE-REVERT` mutation that never landed and was read from a pristine tree). Rule 3e(b) already
  says a control is only a control if it runs from a tree in the baseline's state. It is written
  about **files**, and it says nothing about **machine state**, which is shared by every role, every
  gate, and every sibling mission on this rig. Measured: iteration 97 started **64** CPU spinners to
  supply the "deliberate load" queue row 20 asks for, and `kill $SPINNERS` never reached them —
  `jobs -p` inside a non-interactive `zsh -c` eval did not capture the backgrounded subshells, and
  the command reported no error. Load average sat at **~110 on 16 cores for over an hour**. Inside
  that window ran **two independent `./scripts/verify_go.sh` executions — the controller's and the
  designer's — and BOTH reported `dev` RED at HEAD**, one of them routing it onward as a base
  finding about the repository. Re-run after the spinners were killed, the gate **passes end to
  end**: zero `FAIL` lines, `host/broker` `ok` at **84.597 s** plain and **183.918 s** race.
  **The reason this earns a guardrail rather than a note is that it manufactures FALSE
  CORROBORATION.** Two roles agreed, from separate processes, on a red neither had any reason to
  invent — which is exactly the shape a controller is taught to trust. But agreement is only
  evidence when the arms are independent, and a shared contaminated rig makes them one arm reported
  twice. That is rule 3l (*the fleet is the control group; ask what the failing arms SHARE*) aimed
  inward: the thing they shared was **me**. It also inverts rule 3e(b)'s warning — there the danger
  is reaching for an environmental excuse for a symptom you caused; here the symptom genuinely
  **was** environmental, and the pull was to file it as a defect in `HEAD`.
  Rules: **(a)** any step that deliberately loads the machine (spinners, `-race` sweeps, parallel
  suites, a GPU job) ends with a **verified** teardown — `pgrep -f '<pattern>' | wc -l` must read
  **0** and `uptime` must be quoted — before the next measurement is taken, and the teardown
  assertion is what counts, never the `kill` itself; **(b)** `jobs -p` is unreliable in this rig's
  non-interactive tool shell, so capture PIDs explicitly (`pid=$!` per child, into an array) or kill
  by pattern, and assert the count either way; **(c)** every gate result recorded in the charter or
  handed to another role carries the **load average at which it was taken** — a timing-sensitive
  gate result without a load reading is scoped no better than a `go test` without its toolchain;
  **(d)** when two roles report the same environmental red, treat the agreement as **one**
  observation until you can name something that differs between the arms; **(e)** this is
  rig-global, not mission-local — three missions share this machine, so a leaked load experiment
  corrupts siblings' measurements too, and they have no way to attribute it.

- **A NEGATIVE GREP CANNOT ESTABLISH A PREMISE ABOUT A CLASS OF STATEMENT — IT CAN ONLY FAIL TO
  FIND THE SPELLINGS IT WAS WRITTEN FOR, AND ITS ZERO IS INDISTINGUISHABLE FROM THE TRUTH. ESTABLISH
  IT BY POSITIVE ENUMERATION INSTEAD** (process fix, iter-96; instance 2 of *the control itself was
  the defect*, after iter-95's control scoped so identically to its check that it could only fire on
  the check's own hits). The existing rules cover an empty result (pair it with a control) and a
  green check (state its scope). Neither covers the case where **the control fires, the scope is
  right, and the pattern is still blind** — because the thing it is hunting has a spelling nobody
  thought of. Measured: to establish that `host/store`'s `objects` table is insert-only,
  `grep -rniE '\bUPDATE[[:space:]]+[a-z_]+[[:space:]]+SET\b'` returns **0 repo-wide, tests
  included**, and the pattern is provably sound — it matches a synthetic `UPDATE objects SET …`.
  Yet `store.go` carries **five genuine `ON CONFLICT(…) DO UPDATE SET` upserts** (`:618/:709/:759/
  :836/:978`), invisible because the upsert spelling puts **no table name between `UPDATE` and
  `SET`**. A broken pattern and a genuinely immutable table produce the identical zero, so the zero
  proves nothing either way. **And the failure is symmetric, which is what makes it a rule rather
  than a caution:** in the same iteration the controller handed the designer a `FROM objects|INTO
  objects` enumeration reading **five**, and the designer measured **nine** — four `JOIN objects`
  reads at `journal.go:744/792/918/966` that the adjacency pattern reads 0 on. One grep, two
  spellings, and the author saw only the half being hunted for. **Rule:** to establish a premise
  about what a body of code does or does not do, enumerate every site that NAMES the subject and
  read them one by one; use negative greps to *narrow* an enumeration, never to *conclude* one. The
  tell: you are about to write "there are no X" on the strength of a pattern that describes X's
  most familiar spelling.

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

  | `8/OD-1` | 8 `w-self-mod-vertical` | attended authorization for the irreversible first publish of `world/core@0.1.0` | **ANSWERED** 2026-08-05 (attended, issue #32): *"Approved publish of world/ in ailang extensions - go. Credential is on your machine for this."* **Registered retroactively at iter-61** — see the audit below |
  | `8/OD-2` | 8 | will the AILANG registry add vendor principals/registration before a later World release? | **OPEN, non-blocking.** Controller default: file/route the requirement, label `world/` convention-only |
  | `8/OD-3` | 8 | the `user_version` 1→2 bump makes `store.go:354`'s bare `return nil` reachable — a v1 store then opens fine and never runs `schemaSQL` | **ANSWERED** iter-52 from already-ratified text (`4d/OD-3` alt 1: *"fail LOUD on an unsupported or **un-upgraded** store"*). Non-blocking; binding on `SM.B1` |
  | `10/OD-1` | 10 `w-boundary-gate-tree-mutation` | scan the dependency **closure** against `forbiddenImportPrefixes`, not just direct import specs | **OPEN, deferred on scope** (it changes what the gate ENFORCES). **The doc's premise `V16` is REFUTED**: `cmd/ailang-worldd`'s baseline closure already contains `host/registry` via `host/daemon`, so a closure scan would RED **immediately** today — but that hit is a **false positive on the *epoch* registry** (name collision), so this is MORE blocked, not less |
  | `10/OD-2` | 10 | permanent CI kill-harness, or sprint evidence only? | **CLOSED** iter-61 — controller default taken (sprint evidence). Discharged: `BG.A`'s `M5` ran the SIGKILL harness on all four arms with a firing negative control |
  | `10/OD-9` | 10 | the `ModTime` backstop is darwin-measured only — is a filesystem-independent write-and-restore detector worth its own design? | **PRE-REGISTERED, UNFIRED** (iter-61). Fires only if `BG.C`'s granularity probe reds on `ubuntu-latest`; the pre-authorized `R-EXT4` branch then keeps the four unconditional observables and drops the `mtime` clause. **Allocated as `9`, the next free MISSION-WIDE integer — the design doc and sprint plan both wrote `10/OD-3`, which would have been the FOURTH meaning of `OD-3`** |

  | `9/OD-10` | 9 `w-verify-binary-lockfile` | CI job 1 installs `releases/latest` for `verify_ail.sh` legs 1-2 (today v0.33.0) while `:64`/`:118` pin v0.30.0 — PIN job 1 to the v0.30.0 tag, or ACCEPT the drift and keep only the warning? | **RATIFIED 2026-08-10 (Mark, attended) — ACCEPT, with scope confirmed in the same session.** *"we should accept upstream releases as they will respond to ailang world requests currently"* — i.e. World files upstream asks (e.g. `ailang#633`) and upstream ships fixes, so pinning legs 1-2 to v0.30.0 would mean World's own upstream fixes never reach its primary gate. **THE RATIFIED SCOPE IS THREE CLAUSES, ALL CONFIRMED BY MARK IN-SESSION, because "accept releases" alone does not answer what the gate should do:** (a) legs 1-2 TRACK RELEASED versions — `ci.yml:25` keeps `releases/latest`, no CI config change; (b) a **non-release / `-dirty` build is STILL REFUSED** — an upstream release is not a dev build, and CLAUDE.md's hard rule names `-dirty` specifically; the rig's PATH `ailang` was measured this session at `v0.33.0-70-g1677fcff9-dirty`, which this clause rejects; (c) **`:64` (world-package leg) STAYS PINNED to v0.30.0** regardless — it byte-compares replay goldens generated with v0.30.0, so drift there is a different question that this ruling does not touch. **CONSEQUENCE FOR THE COUPLED ASSERTION, AND IT INVERTS THE ITER-66 CARRY:** iter-66 recorded the hard `verify_ail.sh` assertion as unsafe headless because it would red on the next upstream release. Under clause (b) the assertion is re-shaped — assert *is a release*, not *is v0.30.0* — and in that shape it **cannot** red on an upstream release, only on a dev build. So the remaining half of item 9 is **routable headless after all**, and ACCEPT did not close item 9 empty. **ONE DESIGN CONSEQUENCE NOT YET ADDRESSED (carry `9/CF-A-2`):** the DRIFT warning iter-66 landed compares against the v0.30.0 pin, so under ACCEPT it now fires on **every CI run forever**. A warning that always fires is not a signal — it must be re-shaped to announce the resolved release and warn only on a non-release build (clause b), or on change-since-last-observed. Allocated as `10`, the next free ID per this table. |

  | `9/OD-11` | 9 `w-verify-binary-lockfile` | may a milestone add a **Z3 install step to CI job 2** (`go-verify`)? `ai-check` shells to z3 and CI installs it in job 1 only, so `host/verifygate`'s shim arms — which drive the REAL `.ail` gate under `go test ./...` — cannot reach `verify gate PASSED` there (measured: CI red at `9151797`, `required identity … MISSING from verify.results[]`, the documented V27 class). | **RATIFIED 2026-08-10 (Mark, headless directive on issue #53) — YES.** Verbatim: *"Yes you can install z3 on cicd"*. That widens `9/OD-10` clause (a) **for this install specifically** and nothing else — job 1 still installs `releases/latest` for legs 1-2, and `:64`/`:118` stay pinned. **DISCHARGED by `VL.B` (iter-69, PR #57 → squash `32b086c`).** The pin is now declared ONCE at workflow scope and both jobs install it; `host/verifygate`'s accept-arms assert `verify gate PASSED`. **The cost this ask named was measured in both arms before the fix:** with a solver `rc=0`/PASSED=1, without one `rc=1`/PASSED=0 — and `AILANG_BIN refused`=0, `── Leg 1`=1, `could not parse ai-check JSON`=0 in **BOTH**, so the old contract was satisfied identically by a passing and a failing gate. Settled in CI, not locally: job 2's step log shows `Z3 version 4.16.0 - 64 bit` and `host/verifygate` `ok` twice. Raised iter-68, ratified and closed iter-69. |

  | `14/OD-12` | 14 `w-workbench-read-only` | **the workbench's INTERACTION GRAMMAR — which modalities are first-class?** The whole corpus is silent on it: a search across DESIGN.md / HUMAN-SURFACE.md / SCENARIOS.md for `speech·voice·audio·video·keyboard·shortcut·click·scroll·touch·drag·hover·screen reader·accessib` returns **nothing** except "clickably checkable" (P2), "gesture" used metaphorically (P6/P7), and "not log scroll" as an anti-pattern. So the medium is decided and the *grammar* is not | **OPEN, NON-BLOCKING — build on the default, settle before item 7 ships.** Controller default: **keyboard-first triage** (the surface's nearest ancestor is a code-review queue, and packets are a list you move through, not a canvas), **pointer for the provenance graph** (walking edges is the one genuinely spatial task), **text for goal composition**. **Speech is admissible as an INPUT channel to the composer and MUST NOT be the confirmation channel** — §3 requires the human to sign the *typed* object rather than the prose, and a typed authority envelope cannot be safely confirmed by ear; the signature must be visual and re-readable. **Audio earns one job only**: scenario 4's single batched notification (P5 — an unbatched ping is a budget bug). **Video: no role identified.** Accessibility is unstated corpus-wide, which is notable given P3 already forbids colour as the sole channel for the same epistemic reason — that constraint should be named as an accessibility floor, not re-derived |
  Next free ID: **`OD-13`**.

- **SIX ODs WERE ALLOCATED WITHOUT EVER ENTERING THIS TABLE, AND THE REASON IS THAT THIS
  GUARDRAIL'S OWN ENUMERATION INSTRUMENT FINDS NOTHING** (process fix, iter-61; **fifth instance**
  of the ID-collision class, and the first where the remedy is the cause). Rule (a) above says to
  *"enumerate `### OD-` headings across **every** doc in `design_docs/planned/`"* before allocating.
  Measured at iter-61: `grep -c "^### OD-" design_docs/planned/*.md` returns **0** for all four
  planned docs, against a firing known-positive control (`grep "^### "` on the same files → **23**
  headings). No doc has ever used that shape — item 8 writes `- **8/OD-1 — …**`, and item 10 never
  gives its ODs a heading at all. So an allocator following rule (a) *literally* gets an empty
  enumeration, concludes nothing is allocated, and takes `OD-1`. **Both item 8 and item 10 did
  exactly that**, giving `OD-1` three live meanings (`4e/OD-1` ratified go.mod floor · `8/OD-1`
  ratified publish authorization · `10/OD-1` closure scanning) and `OD-3` four. This is rule 3a
  wearing the registry's clothes: *a search that found nothing is a claim, not a fact* — and the
  claim it licensed was "the integer is free."
  It is not academic. **Mark's 2026-08-05 attended approval landed on `8/OD-1`**, an unregistered,
  thrice-collided integer — precisely the harm iter-43 recorded (*"the first where the collision
  landed inside a human ratification"*), recurring *after* the guardrail, at wider scope.
  **The fix, applied here:** rule (a)'s instrument is replaced by one that cannot miss —
  `grep -rhno '[0-9a-z]*/OD-[0-9]*\|OD-[0-9]*' design_docs/*.md design_docs/planned/*.md
  design_docs/implemented/*.md | sed 's/^[0-9]*://' | sort -u` — **paired with a known-positive
  control** (a registered ID from this table must come back non-empty in the same call), because a
  format-specific pattern is a claim about one author's habits, not about the namespace. Rule (d)
  still holds: the six rows above are **registered as-is, never renumbered**; only the *next*
  allocation moves, which is why `10/OD-9` is `9` and not `3`. The general form, and the reason
  this outranks its two instances: **a remedy is an instrument, and it inherits the same burden of
  proof as the thing it verifies.**

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

- **SUPERSEDED 2026-08-28 (iter-135) — SEE LEDGER ROW `D-WORLD-27`. The chain below is NO LONGER
  the current state, and this paragraph asserted that it was for two days after it stopped being
  true.** The shared `mission-control` SKILL carries a **PROMOTION RULE (Mark, attended
  2026-08-26)** that explicitly supersedes the `D-WORLD-20` suspension: deepseek returns as the
  **fallback link only** — not a rotation peer — reached only when codex is dry, and re-qualified
  on two consecutive real sprint runs returning `ok` with a non-empty worktree diff. The rule
  landed in `sunholo-data/ailang@b4122db56`; the live `mission-world.env` was edited to match at
  the same minute. **And the id actually in force on this rig is NOT the one the fleet intended** —
  it is the `:floor` variant that `fc6e42682` explicitly dropped, because World's driver copy
  predates that commit. Filed as queue row **54**. The suspension text follows for the record.
  ~~**EXECUTOR CHAIN AS OF 2026-08-19 (Mark, attended, `D-WORLD-20`, verbatim "Remove deepseek
  flash"): `codex:gpt-5.6-sol` → `opus`. The `pi:openrouter/deepseek/deepseek-v4-flash-0731`
  link is SUSPENDED and unreachable.**~~ This supersedes the middle link of the 2026-08-10
  three-link ratification below (and of `D-WORLD-ROUTE-1`) and nothing else in either. Evidence
  bar for a routing change is ≥3 rows and was MET at **five**: five consecutive executor runs
  changed zero useful bytes across **four distinct mechanisms** (iter-91 ×2 `stopReason=stop` at
  625 tokens against a 65,536 budget; iter-92 325 MB of `thinking_delta`; iters 93/94 killed at
  the ~300 MB size ceiling with `stopReason="length"` count **0**) — and the detector this skill
  prescribes for the lane, the per-turn `stopReason` assertion, fired on **0 of 5**; only the
  file-SIZE poll ever fired.
  **APPLIED IN `~/.config/ailang/mission-world.env`, NOT THE DRIVER** (`tools/launchd/*` is frozen
  core, `D-WORLD-DRIVER-1`): the driver defaults the tail at `mission-control.sh:331` with
  `${MISSION_EXECUTOR_FALLBACK:-pi:…:floor}` and sources this mission's env file at `:58-59`,
  i.e. BEFORE that line, so `MISSION_EXECUTOR_FALLBACK="${MISSION_EXECUTOR_FALLBACK:-opus}"`
  wins the default and the pi link is never reachable. Same lever and precedent as the iter-23
  PATH fix.
  **PROVEN NON-VACUOUS, TWO ARMS, ONE VARIABLE** (iter-95; the third arm is honestly unrunnable
  and recorded as such rather than asserted). The reading is the driver's own
  `falling back to '<lane>'` line — NOT the `DRY RUN ok:` roles line, which this rig cannot emit
  while an iteration holds the pidfile, and whose absence is an empty result, not a green.
  Arm B, real rather than simulated (codex is genuinely quota-dry until 2026-08-20 05:34, so the
  chain resolved through the fallback link for real): **`falling back to 'opus'`**. Arm C, the
  negative control, byte-identical command with only the env line removed:
  **`falling back to 'pi:openrouter/deepseek/deepseek-v4-flash-0731:floor'`**. Arm A (healthy
  codex → `codex:gpt-5.6-sol`, no fallback line at all) is UNRUNNABLE until the bucket refills;
  it is not claimed. Env file restored byte-identical after the control (`cmp -s`).
  **NOTE FOR THE FIRE THAT APPLIED THIS:** the driver had already exported
  `MISSION_EXECUTOR_MODEL=pi:…:floor` for iteration 95 before the fix existed, so the exported
  env was stale against the directive from the moment it was read — an in-iteration routing
  ratification cannot retro-edit its own fire's exported roles. Read the ratification, not the
  export, for the rest of any iteration that lands one.

  **Prior policy text follows.** ~~**Executor policy (Mark 2026-08-06, attended — SUPERSEDES the default-flip recorded an
  hour earlier at `681990a`, which over-read the directive): `codex:gpt-5.6-sol` REMAINS the
  default; `pi:openrouter/deepseek/deepseek-v4-flash-0731` is the QUOTA-RELIEF REPLACEMENT —
  when the codex bucket is spent, the executor runs deepseek instead of degrading to opus.**
  **THE CHAIN IS NOW LIVE — ATTENDED 2026-08-10 (Mark, verbatim: "we want codex, deepseek,
  opus"). IT IS IMPLEMENTED IN `mission-world.env`, NOT IN THE DRIVER, SO THE `executor=`
  VALUE ON THE "mission iteration starting" LINE IS NOW A RESOLVED PROBE RESULT RATHER THAN A
  DECLARATION.** The 2026-08-06 env pin was the MANUAL instance of this policy, and it exposed
  the failure mode that retires manual pins for good: **a hand-pinned relief lane cannot notice
  the condition that justifies it ending.** Because the pin is `pi:*`, the driver's `codex:*`
  pre-flight branch never executed, so the codex refill of 2026-08-08 11:24 was invisible for
  FOUR DAYS — measured attended (`codex exec --model gpt-5.6-sol` → rc=0 `ok`) with the loop
  still running the relief lane. Rollback executed; the chain replaced it.
  **WHY IT LIVES IN THE ENV FILE:** neither driver implements a chain — upstream has two
  INDEPENDENT role-generic pre-flight loops (`codex:*` and `pi:*`) that each degrade to a
  hardcoded `opus` and do not know about each other, and **World's copy is 101 lines behind
  upstream (487 vs 588) with NO `pi:*` loop at all**, so the deepseek pin ran UNPROBED for
  those same four days. The driver is frozen core → iter-23 PATH-fix precedent. Non-vacuity is
  THREE ARMS, each read off the driver's own resolved roles line under `MISSION_DRY_RUN=1`:
  healthy → `codex:gpt-5.6-sol`; `MISSION_EXECUTOR_CODEX_MODEL=bogus` →
  `pi:openrouter/deepseek/deepseek-v4-flash-0731`; both models bogus → `opus`. Arm B is the
  first time World has probed the pi lane at all. The block carries a DELETE-ON-SYNC marker,
  and a command-line `MISSION_EXECUTOR_MODEL` pin still wins.
  **THE CHAIN DOES NOT KNOW ABOUT THE pi BAR BELOW** — it can resolve `pi:deepseek` for a
  publish-capable milestone. The bar is a CONTROLLER-side rule applied ON TOP of the resolved
  lane: read the lane, then apply the bar; never infer permission from what the chain handed
  you. **STILL OWED BY THE SHARED DRIVER, AND BOTH ARE FROZEN CORE HERE:** **ailang#611** (the
  real per-role chain — carries the probe-blindness constraint that a codex 1-token probe can
  rc=0 on a SPENT bucket, so the chain must ALSO apply at the skill's in-iteration fallback,
  not only at the driver probe) and the **World driver sync** that brings the missing `pi:*`
  loop. Upstream comment + `mission-control` note filed 2026-08-10. Honest cost of the interim:
  codex is probed TWICE per fire on the healthy path, because the driver re-probes the pin it
  is handed and cannot be told "already probed" — ~4.5k subscription tokens, no metered $.
  V1 precedent for the
  lane: ≥3 executor iterations, sonnet 91/88 zero-blocking, ~$0.006/iter REAL metered (posts
  to the $5 ledger, unlike quota lanes).~~ **pi has NO sandbox** — every pi run carries V1's
  discipline: re-verify the main checkout byte-identical to preflight after the run; the
  codex-sandbox plan assumptions (S-7 "executor cannot commit", `.snap/M<k>/` reconstruction,
  UNINFORMATIVE-under-sandbox caveats) do NOT apply — state the routing delta in the
  directive, exactly as iter-58 did for opus. **SM.B2a exception (coordinator, same
  exchange): the first irreversible-publish-capable milestone does NOT run on an unsandboxed
  lane by default** — for SM.B2a use the codex lane, or an opus once-file, or carry an
  explicit extra integrity gate (worktree isolation + `env -u AILANG_REGISTRY_API_KEY`
  verified in the directive, checkout byte-diff both sides).
- *(prior default, still the documented fallback)*: NON-Anthropic lane — `codex:gpt-5.6-sol`
  (`MISSION_EXECUTOR_MODEL`), per the shared quota plan. **iter-19 routing incident RESOLVED
  by Mark 2026-07-27 (attended, option c): codex-cli upgraded 0.137.0 → 0.145.0 (npm -g) and
  the pin LIVE-PROBED working** (`env -u OPENAI_API_KEY codex exec --model gpt-5.6-sol` →
  PIN-OK, ~21 tokens, ChatGPT-subscription OAuth — Mark verified server-codex OAuth2 same
  day, so the lane works headless in mission loops AND evals). The opus fallback stays the
  documented degrade path only. Upstream driver-probe gap (probe omits `--model`,
  false-greens an unusable pin) remains filed as ailang#486 for the v1 lane.
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
5. [**[LANDED 2026-08-26 (iter-127)] — ROW COMPLETE. `P6.V` was the last milestone standing and it SHIPPED: PR [#96](https://github.com/sunholo-data/ailang-world/pull/96) → squash [`699f592`](https://github.com/sunholo-data/ailang-world/commit/699f592), Gate 3b GREEN on the MERGE commit (`checks=2 == expected=2`, `not_green=0`, `event=push`, parent control 2), evaluator `sonnet` **94/100 PASS, zero blocking**. **Branch A fired**, not the fallback: `commitBoundaryHolds` is Z3-`verified` on the pinned v0.30.0 binary, so the gate's floor moved **10 → 11** identities (40 named tests and 9/9 world-package steps unchanged — AC10's "nothing else moves" holds). **AC10 and AC17 are CLOSED; AC15 was closed by `P6.T` at `8b196c3`; AC16 travelled with `P6.D`.** **The milestone's real touch set was SIX files against the two this doc's table named** — the coupling is now queue row 43, because a commit message is not where the next floor-raiser will look. **Two clause-map overclaims the evaluator found were fixed in-iteration** (`1e8a018`, comment-only): `receiptCount <= 1` restates conjunct 1 rather than adding to it, and `clientDisconnected` does no logical work at all — so Decision 6's disconnect sentence is recorded as UNDISCHARGED rather than implied. Honest gaps 2 → 3. The remaining halves of the original `w-mcp-projection` live on as row 40 (A2A, blocked on row 39) and `w-mcp-dispatch-projection` (blocked on `ailang#885`). **Prior-state header follows.** ~~[NEXT] — `P6.T` LANDED 2026-08-26 (iter-126); SPLIT #2 APPLIED; ROUND 5 BLOCKED ON TWO NARROW SURFACES AND BOTH FIXES WENT IN VERBATIM UNDER THE CARVE-OUT; `P6.D` IS DEFERRED OUT OF THIS ROW ENTIRELY. WHAT REMAINS HERE IS `P6.V` ALONE, BLOCKED ON NOTHING.~~** **`P6.T` SHIPPED** — PR [#95](https://github.com/sunholo-data/ailang-world/pull/95) → squash [`8b196c3`](https://github.com/sunholo-data/ailang-world/commit/8b196c3), Gate 3b GREEN on the MERGE commit (SHA-addressed, `checks=2 == expected=2`, `not_green=0`, run existence `total=1 event=push`, parent control `checks=2`); evaluator `sonnet` **96/100 PASS, zero blocking**, generator≠judge. It moved **FOUR** `ci.yml` pin sites plus `go.mod`, not the two the milestone named — `actions/setup-go@v5`'s `go-version: '1.25.6'` at `:28`/`:109` was found by the repo-wide sweep quorum round 5 forced, and moving `GOTOOLCHAIN` alone would have left `setup-go` installing 1.25.6 against a 1.26.6 demand. Non-vacuity is a measured pair, unpinned and outside any sandbox: `verify_go.sh` **rc=1** (`FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go`) **→ rc=0**. Two honest non-kills recorded: arms 3/3′ (`ci.yml:28`/`:109`) **SURVIVED** → **row 41**; arm 4 (the committed canary under a deny-listed toolchain) is **NOT VERIFIABLE post-landing** because the module floor now rejects `go1.26.5` before the canary runs → **row 42**. **SPLIT #2 APPLIED**: `P6.B-A2A` → new child [`w-a2a-session-projection.md`](../design_docs/planned/w-a2a-session-projection.md), BLOCKED on row 39, objection carried VERBATIM (3/3 True, two negative controls False) — tracked as **row 40**. **ROUND 5 (both reviewers PRESENT and EXTERNAL, `absent_reviewers` EMPTY, `metered=$0.219099`) BLOCKED, then CARVE-OUT**: `gemini-3-1-pro` refuted premise row **N7** as a narrowed search (material — see above); `gpt5-6-sol` rejected **`P6.D`** as speculative core growth, *a defect split #2 INTRODUCED one iteration earlier* by moving `P6.D`'s only real consumer to the child and re-anchoring on a dead symbol reference. Both fixes applied verbatim: **`P6.D`, `host/daemon/protocol_use.go`, `AC16` and its two mutations are OUT of this row; dependency admission is now ATOMIC WITH ITS FIRST REAL CONSUMER** and travels with whichever child unblocks first (spec intact in both children). Nothing discarded — re-sequenced. **N7 THEN FAILED TWICE MORE**, caught by the planner and the evaluator in turn, each fix having drawn a different boundary; it now sweeps with **no exclusion at all**. Generalisation proposed to V1: *the scope of an exclusion IS part of the claim.* **ROUND COUNT: FIVE**, plus two splits and a carve-out — data about this loop's SCOPING, not about the doc, and surfaced to Mark for that reason. **NEXT ACTION, blocked on nothing: sprint `P6.V`** (~0.3d, verified commit-boundary law in `world/*.ail` + `REQUIRED_VERIFIED`; objection-free in five rounds) — the last milestone in this row. **Prior head text follows.** ~~**[NEXT] — SPLIT #1 APPLIED 2026-08-25 (iter-125), AND THE RE-QUORUM FOUND A SECOND, DIFFERENT, NEWLY-LOCALISED SURFACE: SESSION AUTHORITY. `P6.T`/`P6.D`/`P6.V` REMAIN OBJECTION-FREE ACROSS ALL FOUR ROUNDS; A SECOND SPLIT IS OWED BEFORE ANY SPRINT.** iter-124's prescribed action was executed in full: the Fable designer applied the split (parent 1040 → 1204 lines, MCP half carved into new child doc [`w-mcp-dispatch-projection.md`](../design_docs/planned/w-mcp-dispatch-projection.md), 173 lines, carrying `gpt5-6-sol`'s round-3 objection VERBATIM and honestly BLOCKED on [`ailang#885`](https://github.com/sunholo-data/ailang/issues/885)), and **`D-WORLD-26` was answered by Mark — ARM A, `Authorization: Bearer`** (attended, `#89`, 2026-08-25T19:06:41Z, verbatim one-character comment `A`), closing round 3's second objection. Both round-3 objections were therefore discharged. **ROUND 4 (both reviewers PRESENT, `absent_reviewers` EMPTY, `metered=$0.213298`) IS THE DECOMPOSITION SIGNAL FIRING AS WRITTEN:** `gemini-3-1-pro` **PASSED** — its first pass in four rounds, leaving only a narrow P6.V wording nit the controller applied verbatim — while `gpt5-6-sol` rejected on **one surface only**, session authority, which is a property of **`P6.B-A2A` alone**. The objection was **CONFIRMED FIRST-PARTY, NOT FORWARDED** (rule 3f): over `host/` at `2e44e3e`, `Bearer` **0**, session-lookup functions **0**, `Credential` **128** all OUTBOUND registry-publish, `Authenticate` **29** all evidence-envelope, same-scope known-positive control **181**, fresh negative control **0** — so `D-WORLD-26` settled the credential **envelope** and the **contents** (who mints it, where credential→(episode,grants) lives, what expires it) were never built. That is a gap the doc **fails to fix**, not one it introduces, so it is filed as **queue row 39 `w-session-authority`** rather than absorbed — and it gates `P6.B-A2A` and nothing earlier. **NEXT ACTIONS, none blocked on a human:** apply split #2 (carve `P6.B-A2A` out behind row 39), re-quorum the `P6.T`/`P6.D`/`P6.V` remainder once, then sprint **`P6.T`** (toolchain floor `go1.25.6`→`go1.26.6`, ~0.1d, independently mergeable, zero objections in four rounds). **ROUND COUNT SAID OUT LOUD:** this doc is at four quorum rounds plus a carve-out revision — data about this loop's SCOPING, not about the doc, and surfaced to Mark for that reason.~~ **Prior head text follows.** ~~SPLIT 2026-08-25 (iter-124): THE SEAM SHIPPED BUT IT ONLY CARRIES HALF THE SURFACE. `serveapi/protocol` HAS NO MCP DISPATCH, SO THE MCP HALF IS BLOCKED ON A NEW UPSTREAM ASK ([`ailang#885`](https://github.com/sunholo-data/ailang/issues/885)) WHILE THE TOOLCHAIN FLOOR, THE PINNED DEPENDENCY, THE COMMIT-BOUNDARY LAW AND THE A2A HALF ARE ALL EXECUTABLE NOW.** The design doc was REVISED this iteration (Fable designer, 641 → 974 lines) because the July draft's central premise was falsified by upstream delivery, and the revision is sound — but re-quorum **BLOCKED at round 3, both reviewers PRESENT, `absent_reviewers` EMPTY** (`metered=$0.1658`), and `gpt5-6-sol`'s objection was **confirmed first-party by the controller** rather than inherited. **THE FINDING, MEASURED WITH FIRING CONTROLS.** `serveapi/protocol` carries the full **A2A** wire surface and the MCP **envelope** helpers (`WriteMCPEnvelope`, `RequestID`, `ValidateMCPName`, `AuthorizationStatus`) — but **no JSON-RPC method dispatch**. That lives in `serveapi/mcp_handler.go`, which delegates to `github.com/modelcontextprotocol/go-sdk`; control, in the same call: SDK import count **1** in `mcp_handler.go` vs **0** in `a2a_handler.go`, which is **180 lines over stdlib + `protocol`** and is the existence proof of the shape World needs. So the MCP half has exactly two routes and **both are closed**: reimplement JSON-RPC dispatch, forbidden by P1 / the Design Freeze / AC1 / DESIGN.md §3.7; or import the SDK, measured over BOTH gated patterns (`./host/daemon/... ./cmd/ailang-worldd/...`) at **249 → 283**, **+34 packages across 5 new module roots**, and `TestDaemonDependencyAllowlist` reds naming **28** disallowed packages — among them `golang.org/x/oauth2`, `go-sdk/auth` and `go-sdk/oauthex`, i.e. an **outbound-credential stack in the daemon core**, which breaches clause 2 AND clause 3. Not a size problem; a guardrail problem. **THE `protocol` ARM, BY CONTRAST, IS AS CLEAN AS IT GETS**: closure **249 → 250**, the single added package IS `serveapi/protocol` itself, removed set **EMPTY** (sentinel control fired), the package imports **only stdlib** across all four of its files, and the unmodified gate reds naming **exactly one** intruder — clearable by **one narrow PACKAGE-path allowlist line**, never the module root (which would admit `internal/apiserver`'s measured 476 disallowed packages). **DISPOSITION — SPLIT, a controller routing call and explicitly NOT `needs-human-review`.** [`ailang#885`](https://github.com/sunholo-data/ailang/issues/885) filed on the cross-repo channel asking for an SDK-free MCP dispatch seam, with all of the above and the `a2a_handler.go` existence proof — this is **D-WORLD-5's own prescribed default executing as written** (a disallowed graph asks upstream, never a broad relaxation), the same route that produced `#764` → `v0.33.2`, and it is NOT a new human ask. **THE TOOLCHAIN PRECONDITION IS CONFIRMED AND IS STILL THIS ROW'S FIRST MILESTONE**, two-arm with a firing control: `go get github.com/sunholo-data/ailang@v0.33.2` is **rc=1** under the repo pin `GOTOOLCHAIN=go1.25.6` (`requires go >= 1.26.6`) and **rc=0** under `go1.26.6`, while the known-positive control `go get github.com/google/uuid@v1.6.0` is **rc=0** under the pin — so the refusal is the version floor, not a broken probe. `v0.33.2`'s `go.mod` declares `go 1.26.6`; this repo pins `go1.25.6` at `ci.yml:21`,`:102` and in `go.mod`. **AND A CORRECTION THE DESIGNER CAUGHT IN THE CONTROLLER'S OWN MEASUREMENT**: the rig's base `go` binary is **`go1.26.4`**, which is IN `verify_go.sh`'s array-literal-miscompile deny-list (`:214-224`, `go1.26.0`–`go1.26.5`) — the controller's `go version` read `go1.26.6` only because its cwd was a `go 1.26.6` module and `GOTOOLCHAIN=auto` switched. `go1.26.6` is not deny-listed. **NEXT ACTIONS, none blocked on a human**: apply the split (one designer revision, next iteration), then sprint `P6.T` (toolchain floor, ~0.1d, independently mergeable), `P6.D` (pinned dep + the one allowlist line + a narrowness test, ~0.15d), `P6.V` (the `"verified"` residual — a commit-boundary law in `world/*.ail` pinned in `REQUIRED_VERIFIED`, raising the floor from 10 identities, ~0.3d). Only the MCP half of `P6.B` waits on `#885`; `D-WORLD-26` (session credential carrier) gates the rest of `P6.B` and nothing earlier. ~~ **Prior head text follows.** ~~UNBLOCKED 2026-08-24 (iter-120) BY DIRECTIVE-DRIVEN VERIFICATION — THE UPSTREAM ASK SHIPPED, AND IT BROUGHT A TOOLCHAIN FLOOR WITH IT. [`ailang#764`](https://github.com/sunholo-data/ailang/issues/764) IS DELIVERED: `serveapi/protocol` is present in tag **v0.33.2** (`63e7909f`) — one tag EARLIER than the v0.34.0 upstream recommended. Verified first-party through THIS repo's own committed gate rather than upstream's `go list` (sibling-claim ghost discipline): a probe-import of `serveapi/protocol` into `host/daemon` with the allowlist UNCHANGED reds `TestDaemonDependencyAllowlist` naming **exactly one** intruder and nothing else, and **one narrow allowlist line** — the PACKAGE path `github.com/sunholo-data/ailang/serveapi/protocol`, never the module root, which would pass the gate while admitting `internal/apiserver` and the whole cloud subtree — turns it green at **250** packages against a pristine **244**, i.e. **+6 packages of which exactly 1 is non-stdlib**. Against iter-90's measured **476 disallowed packages across 86 module roots** for the old `serveapi` seam. Upstream's guarantee is CI-ENFORCED, not a snapshot: `scripts/check_protocol_closure.sh` ships in the tag, `make check-protocol-closure` rc=0 (`non-stdlib count: 1`), refusal self-test rc=0 across **5** arms incl. an intruder arm and vacuity probes R1/R2/R3/R4/R6/R7. **THE REMAINING PRECONDITION, AND IT IS THIS ROW'S FIRST MILESTONE:** v0.33.2's `go.mod` declares **`go 1.26.6`** while this repo hard-pins `GOTOOLCHAIN: go1.25.6` (`ci.yml:21,102`), so `go get …@v0.33.2` is **rc=1** under our pin and **rc=0** under go1.26.6. The repo already anticipated the move — `verify_go.sh:214-224`'s deny-list enumerates exactly `go1.26.0`–`go1.26.5` and names the canary as the version-agnostic detector — and the canary agrees on a pristine `go.mod`: `go1.25.6` rc=0, `go1.26.5` **rc=1** (miscompile, known-positive control fires), `go1.26.6` rc=0. Full `verify_go.sh` under go1.26.6: **rc=0**, 38 `ok`, 0 FAIL, plain and `-race`, race control armed. **SCOPE, NAMED NOW:** `CallbackRunner`, the embedded A2A `http.Handler` and the MCP handler are deliberately NOT in `protocol` (the MCP SDK spans 9 module roots, only 1 allowlisted — it would fail the gate it exists to pass), so **World writes its own HTTP handlers and callback-bounding**. That is D-WORLD-5 executing as written, NOT a new human ask. All four `#498` requirements are present (`SessionResolver`, `ToolSource`, `CallerSurface`/`AuthorizedSurface`, `ToolDescriptor`, plus `WriteMCPEnvelope`/`RequestID`/`ValidateMCPName`). ORDERING is the one open ask: this row is the sole blocker on M4 (item 5 → item 6 → M4, the value gate) but item 14 is mid-sprint at 8 of 11 — put to Mark as a one-word fork. **Prior head text follows.** ~~**RE-BLOCKED ON UPSTREAM 2026-08-18 (iter-90) — ARM A's DEPENDENCY CONDITION FAILS BY MEASUREMENT, AND MARK'S OWN RULING PRE-AUTHORIZED THE ROUTE OUT, SO THIS IS NOT A NEW HUMAN ASK. Filed [`sunholo-data/ailang#764`](https://github.com/sunholo-data/ailang/issues/764) requesting a protocol-only module; the revision round is blocked on that, not on a decision.**~~ `serveapi` is an API seam but NOT a dependency seam: `serveapi/serveapi.go` is 201 lines whose only non-stdlib import is `internal/apiserver`, so importing the facade links the whole runtime. Measured at `v0.33.1` (go1.26.6, `go list -deps`): `serveapi` **479** non-stdlib packages, `internal/apiserver` **478** — the facade adds exactly one. Controls fire, so 479 is a measurement and not an artifact of the instrument: `cmd/registry-validator` **6**, `cmd/wasm` **12**, `cmd/astdump` **14**. **304** of the 479 are cloud/telemetry/GCP/ollama/CGO-sqlite, including `cloud.google.com/go/auth`, `google.golang.org/grpc`, `go.opentelemetry.io/otel`, `github.com/ollama/ollama` and `github.com/mattn/go-sqlite3`. World's committed zero-cloud gate `TestDaemonDependencyAllowlist` (`host/daemon/daemon_test.go:831`) allows **11** module roots and the daemon graph today is **239 packages / 46 non-stdlib / exactly those 11**; importing `serveapi` would add **476 disallowed packages across 86 module roots** (known-positive control: `github.com/google/uuid` IS in the closure and is correctly NOT in the disallowed set, so the partitioner works). D-WORLD-5's ruling says a disallowed graph routes to Open Decision 4's default — ask upstream for a protocol-only module, **never a broad relaxation** — so the controller executed that arm rather than parking. **SECONDARY, MEASURED, AND IT UNBLOCKS SOMETHING ELSE: `go1.26.6` FIXES the go1.26.x array-literal miscompile** that pins this repo to `GOTOOLCHAIN=go1.25.6`. The committed repro's own README says *"including the newest stable, so there is nothing to upgrade to"* — true on 2026-07-30, false now. Run today, all three arms in one call: `go1.26.5` → **BUG**, `go1.26.6` → **OK**, `go1.25.6` → **OK**, and the committed `run.sh` reports both of its own controls firing. That matters here because `v0.33.1`'s `go.mod` declares `go 1.26.5`, which World's deny-list currently refuses; it also unblocks queue item 4e's parked remediation. **Prior head text follows.** ~~**D-WORLD-5 RESOLVED 2026-08-17 (Mark, attended) — A: IMPORT UPSTREAM `serveapi` PINNED AT v0.33.1.** Full ruling + conditions in the decision ledger. In short: this is Decision 2's own version discipline executing as written (pin the first tagged release containing the seam, record it in `go.mod`/`go.sum`); the revision must (i) run P6.A's frozen conformance fixture against `serveapi` v0.33.1 and (ii) audit the dependency closure against `TestDaemonDependencyAllowlist` BEFORE any `go.mod` change — a disallowed graph routes to the doc's Open Decision 4 default (ask upstream for a protocol-only module), never a broad relaxation. The verified commit-boundary residual below (prereq 3) stays a P6.B prerequisite and is `world/*.ail` work. ROUTING NOTE: codex is exhausted fleet-wide until 2026-08-20 05:34 — the designer rotation's `codex:` entry will probe-fail; fall to the NEXT rotation entry, never `$MODEL`. **Prior head text follows.** ~~**ALL THREE P6.B PREREQUISITES ARE NOW DISCHARGED — MEASURED 2026-08-15 (iter-88). THIS ROW IS NO LONGER BLOCKED, AND THE CHARTER'S "QUEUE IS FULLY BLOCKED" HEADLINE (iter-87) IS THEREBY CORRECTED.** What this row is blocked on now is a **design REVISION**, not a prerequisite — see the A/B below.~~ Measurements, each with a firing control (still live evidence, not struck):
   **(1) UPSTREAM SEAM — DISCHARGED, and this is the blocker-rot class in its textbook form.** `sunholo-data/ailang#498` is still `OPEN`, last updated `2026-08-04`, and iter-87 re-checked exactly that and recorded it as blocking. **Its STATE was never the instrument; its PURPOSE was.** Lane B has landed in full: `f5ebcc0b5` M1 *"public embeddable contract, authorized surface, bounded callback runner"* (#585), `6166adab8` M2 *"request-scoped MCP adapter and frozen wire envelopes"* (#592), `b8c038647` M3 *"A2A projection, Mount, and one exposure gateway"* (#601). The artifact is a **public module-root package** `serveapi/` (NOT `internal/`), and it answers this row's four named requirements **four-for-four**: caller-owned mux → `func (s *Server) Mount(mux *http.ServeMux)`; principal resolved before discovery *or* invocation → `SessionResolver.ResolveSession(ctx, *http.Request)` wired into BOTH the `Tools` and `Invoke` closures of BOTH the MCP and A2A handlers; caller supplies the exact visible descriptors → `ToolSource.Tools(ctx, Session) ([]ToolDescriptor, error)`; **no built-in tool unless the caller supplies it** → `submit_feedback` is **0** in `embedded_mcp.go`, `embedded_a2a.go`, `authorized_surface.go` and `callbacks.go`, while the SAME grep in the SAME directory returns **4** on `feedback_tool.go` (known-positive control fires, so the zeros are measurements — rule 3a(i-d), control scoped to the path under test). The process-wide/session gap that produced this row is closed too: `Session` is threaded through discovery *and* invocation. **RELEASED**: M1 is in `v0.33.0`, M2 and M3 in **`v0.33.1`** (control: an older commit resolves to 5 tags, so `--contains` works). The complete seam therefore requires **v0.33.1**.
   **(2) TRANSITION REGISTRY — DISCHARGED iter-75** (item 11 COMPLETE, `TR.C` GREEN at `625fb89`). Unchanged.
   **(3) COMMIT-BOUNDARY CONTRACT — its stated basis is now FALSE, with ONE precise residual.** This prerequisite's own words were *"No landed API exposes these."* Measured at HEAD, all three named artifacts ARE landed **public** APIs: an atomic not-started-vs-committed contract → `JournalIntent` (`host/store/journal.go:26-28`, *"the canonical statement of a planned commit"*) bound inside a transaction by `bindCommitIntentTx` (`host/store/store.go:1015`); a stable invocation/idempotency ID → `InvocationID`, threaded through journal, receipts and recovery (`host/store/journal.go:29,42,51,63,80,95,103`; `host/broker/publish_op.go:220`); a queryable durable receipt → `func (s *Store) GetReceipt(id string) (Receipt, bool, error)` (`journal.go:813`) and `GetEffectReceipt` (`:852`), with `recoverCommitPending` (`host/broker/recover.go:126`) consuming them. **THE RESIDUAL, and it is the one word doing the work: "a VERIFIED commit-boundary contract."** In THIS repo "verified" has a specific meaning, and by that meaning it is NOT satisfied — none of the **10** pinned Z3-proven identities is a commit-boundary law (`applyRevision`, `isValidNextWorld`, `sameRef`, `servesEntry`, `gradeOf`, `timeoutOutcome`, `timeoutFiredLegally`, `validEscalation`, `validDefer`, `wellFormedSchedule` — `REQUIRED_VERIFIED`, `scripts/verify_ail.sh:262-268`). So the Go surface exists and the *proof* does not, which is a pure-core question routed through `world/*.ail`, not more `host/` work. Scope it in the revision; do not silently read the landed API as discharging the word "verified".
   **⚠ THE DESIGN DOC'S CHOSEN ARCHITECTURE WAS SELECTED UNDER A PREMISE THAT IS NOW FALSE — this is the THIRD instance of prescription-rot in this charter** (iter-82's renderer row and iter-70's item-7 park condition are 1 and 2), and the sharpest, because *the clearing of the blocker IS the falsifier*: the doc took **path (c), "a narrow public serving seam over the existing `internal/apiserver`"**, explicitly because reuse paths (a) and (b) were rejected on evidence and upstream exposed nothing public. Upstream now ships that seam. Path (c) may reduce to *"import `github.com/sunholo-data/ailang/serveapi`"*. **A REVISION ROUND IS REQUIRED BEFORE ANY SPRINT.**
   **ANSWERED 2026-08-17: A** (attended; see the head above and the ledger — kept here unstruck because the two axes it names are the revision round's working constraints). **THE A/B (one word — DIRECTION, so the carve-out is foreclosed and the controller may not settle it).** Measured constraint: World's `go.mod` requires **only** `modernc.org/sqlite` and has **no dependency on the ailang Go module at all**, so this is not a version bump — it is World's **first** upstream Go dependency. **A** = import `serveapi` at `v0.33.1`; smallest surface, tracks upstream, but couples World's host to upstream's Go module and re-opens the `.ail` verifier pin (`v0.30.0`, `verify_ail.sh:62`) as a second, independent axis. **B** = stay dependency-free and build path (c) as designed, now informed by upstream's published contract. Note the frozen-core rule cuts FOR **A** (a real module dependency is the sanctioned path, the opposite of a vendored fork) while the slim-kernel rule cuts FOR **B**; that tension is exactly why this is Mark's call and not the loop's. **Prior head text follows.** ~~**THE TRANSITION-REGISTRY PREREQUISITE IS DISCHARGED as of 2026-08-12 (iter-75).** `P6.B` is no longer blocked on the registry; re-verify its OTHER two named prerequisites at pick time before routing (rule: a declared blocker is a claim too, and the ones describing someone else's work rot fastest).~~ ~~STILL BLOCKED on ONE prerequisite — but as of 2026-08-11 (iter-70) that prerequisite is DESIGNED, NOT ABSENT: see item 11 `w-transition-registry` (doc LANDED, PR #58 → `11fb1fd`). P6.B's prerequisite is satisfied only when item 11's `TR.C` binding gate is GREEN — `TR.A`+`TR.B` deliver the mechanism, not the enforcement. Item 5 promotes to `[NEXT]` when item 11 completes.** Prior text follows. ~~**STILL BLOCKED — but on ONE prerequisite now, not three; RE-MEASURED 2026-08-05 (iter-50) rather than inferred from the `#498` stamp.**~~ The attended `#498` SEAM stamp cleared prereq **1** (upstream seam, verified live on released v0.33.0), and item 4b's landing cleared prereq **3** (`Commit.InvocationID` + `GetReceipt` + the three-state receipt law ARE the atomic not-started-vs-committed contract, the stable idempotency ID and the queryable durable receipt that prereq asked for). Prereq **2 is HALF clear, and the missing half is the one the item is named for**: `host/broker/broker.go:45` now defines `type Session struct` with `NewSession(store, episodeID, grants, registry)` — the broker session API exists — but a repo-wide `grep -rniE '[Tt]ransition[ -]?[Rr]egistry' host/ world/ cmd/` at `de80792` still returns **ZERO**, with the same-call known-positive control (`registry` in `host/registry/registry.go` → **25**) firing, so the absence is a measurement and not a failed grep. `host/registry` remains the *interpreter epoch* registry, a different thing. **This row was one stamp away from being promoted to `[NEXT]` on the strength of a single satisfied prerequisite; the re-measurement is what stopped it.** Whoever picks this next writes the transition registry first, or re-scopes P6.B around its absence. Prior text follows. ~~**BLOCKED 2026-07-28 (iter-24) — DOC LANDED + 2 QUORUM ROUNDS + carve-out revision applied.**~~
   Doc: `design_docs/planned/w-mcp-projection.md` (codex/gpt-5.6-sol designer, rotation; 641 lines;
   NO `.ail` sketch — protocol/session invariants are host-boundary behaviour, so the required-check
   manifest is untouched at 4/11/14 — the module count was 9 when this row was written and is **11** as of iter-70; measured again iter-71). **Milestone P6.A is DONE this iteration** (upstream finding
   filed + this record landed). **P6.B is blocked on THREE named prerequisites** and needs no
   re-design when they clear~~] **w-mcp-projection** · clause-6 · project the transition registry over
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
      ~~Gated behind `w-effect-broker-m3` (item 4, PARKED).~~ **STALE — CORRECTED 2026-08-11
      (iter-70), Gate-2 blocker-rot rule.** Item 4 **LANDED 2026-07-29 (iter-35)**; its doc is in
      `design_docs/implemented/`. Twelve iterations carried this clause. The **broker half of this
      prerequisite is SATISFIED** — `host/broker/broker.go:46` `type Session`, `:58`
      `NewSession(store, episodeID, grants, registry)`, `host/broker/decide.go:15` `Capability` —
      and `broker.Registry` is `map[string]Handler`, an **effect-handler** map (`broker.go:34-35`),
      NOT a transition registry; conflating the two is the defect the design most had to avoid.
      Only the **transition-registry half** was missing, and it is now item 11.
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
   kernel's `Evidence` ADT (`world/types.ail:23-28`) has **five** variants. The mapping gap is now
   closed by the total `gradeOf(Evidence)`: `TestReport`→TESTED,
   `CompilerOutput`/`HumanApproval`/`RecordedEffect`→ATTESTED, and `AiReview`→CLAIMED. **PROVEN's
   own stated producers (Z3, replay) still have no `Evidence` carrier at all**.
   `HumanApproval` is precisely the evidence the approval inbox — item 7, gated on this doc —
   would emit, and grading a human's ratification as CLAIMED ("agent said so, unverified") is
   plainly wrong. The total mapping now makes the doc's cardinal-sin anti-pattern, **grade
   laundering**, enforceable for decoded `Evidence`; the carrier gap for `PROVEN` remains.
   Also measured: `Proposal.confidence` is a **bare float
   with no evidence ref** (`AiReview` carries one), so rendering it would violate the doc's own
   confidence-theater anti-pattern → new ratification point 7.5]
   **w-human-surface** · clause-5 · the founding UX design for the human↔AI surface:
   [HUMAN-SURFACE.md](HUMAN-SURFACE.md) — seven interaction principles (decision packets ·
   grounded prose · trust-gradient rendering PROVEN/TESTED/ATTESTED/CLAIMED · time-as-navigation
   · attention choreography · five zoom levels · speculation-as-gesture) + sign-the-type input
   grammar + anti-patterns (grade laundering is the cardinal sin). Mark named this the vision's
   key surface ("an ideal AI language in an AI state machine OS with a NOVEL human interaction
   surface"); no human-facing sprint routes before it is ratified · ratification session ~0.5d
7. [**PARK CONDITION CORRECTED 2026-08-11 (attended) — BLOCKER ROT, THE ITER-70 CLASS AGAIN.** The
   row read *"PARKED until 4 AND 6b land"* and **BOTH LANDED**: item 4 `w-effect-broker-m3` on
   2026-07-29 (iter-35, doc in `implemented/`), and 6b's §7 ratification on 2026-07-28 (Mark,
   attended TRIPLE RATIFICATION — `bc467f1` *"human-surface v0.3 s7 formally ratified"*; the 6b row
   was flipped 2026-07-31 in `ee75837`). So by its own stated condition this item has been unparked
   for **13 days** while reading as blocked. **NOTHING HERE IS WAITING ON A HUMAN.** The REAL gate is
   the one 6b's row names — *"unblocks item 7 the moment item 5's MCP work lands"* — and item 5 is
   itself gated on item 11's `TR.C` binding gate being GREEN (`TR.A`+`TR.B` deliver the mechanism,
   not the enforcement). **PARKED until item 5's `P6.B` lands**, and no longer on anything else.
   Chain at this correction: `TR.A2` → `TR.B` → `TR.C` → item 5 `P6.B` → this item.]
   **w-approval-inbox** · clause-5 · the approval inbox + provenance
   walk, first as CLI/generated projection (SCENARIOS.md scenario 1/3), **built to
   HUMAN-SURFACE.md** · ~2d
8. [**[IN-SPRINT] — `SM.D0` LANDED 2026-08-10 (iter-67), PR #55 → squash `a4452d1`, dev CI green BOTH jobs SHA-addressed and step-log verified on the merge commit (`ailang-code verify gate` 11/11 · `go host build + test gate` 13/13, `failed=0`, `checks=2` = expected 2) in a 0-incident window. Evaluator `sonnet` **88/100, ZERO BLOCKING**. THE ENTRYPOINT NOW EXISTS: `SM.D` IS A REAL PROCEDURE AND IS ATTENDED-ONLY — never headless, never CI. THIS ITEM HAS NO HEADLESS-ROUTABLE MILESTONE LEFT; item 9's three pieces are the routable work.**

    **WHAT `SM.D0` SETTLED.** `cmd/world-publish` (`packet | approve | publish | reconcile`, default outcome **STOP**, exit 3, `STOP fence=<name>`), `host/broker/publish_op.go` (the wiring — `MintAttendedApproval` through the LANDED traversal, `InvokeAttendedPublish` with exactly one dispatch), `host/pkgproj/readypacket.go`, and a runbook whose **Stage B now carries commands**. +4735/−24, 15 files. The load-bearing fence is a **controlling-terminal check**, chosen because it is the one thing this loop is structurally unable to satisfy (stdin is a socket; `open(/dev/tty)` → *device not configured*), with an `os.SameFile(stdin, ctty)` branch because **`/dev/null` IS a character device** and a naive isatty would admit `--live < /dev/null`. `R-CI` is a DECLARED TRIPWIRE, not the fence. **14 refusal branches, 22 mutations, 22 killed**, each with anchor count, pre/post sha256, a rc=0 build on the mutant, a `-run`-scoped kill and an inverse `-skip` arm.

    **READ BEFORE STARTING `SM.D`:** (a) **`8/CF-D0-1` — the fence is enforced at the `cmd/world-publish` CALL SITE ONLY.** Any future Go code in this module — most concerningly `cmd/ailang-worldd`, which CI and this loop actually run — can call `NewRegistryPublishHandler` / `InvokeAttendedPublish` / `DecideApproval` directly and skip all 14 refusal branches. `AC30` greps shell scripts and `ci.yml` for the string `world-publish`; it does **not** grep Go source for direct broker calls. This is the same *absence-of-a-caller-is-the-safety-property* pattern the milestone started from, and it wants a gate. Judge non-blocking 3, accepted. (b) **`MUT-D0-FENCE-ORDER` is killed by AC22 AND by 6 of 15 AC21 rows** (`R-CI`, `R-TTY-OPEN`, `R-TTY-CHARDEV`, `R-TTY-SAMEFILE`, `R-PHRASE-EOF`, `R-PHRASE`) — the shipped code's comments say so now; do not re-derive the original "AC22 is the only killer" claim, which was false. The cause is that those rows' fixture uses a **loopback** origin, so a hoisted constructor refuses on loopback/ambient-credential first. (c) **`AILANG_REGISTRY_API_KEY` is AMBIENT on this rig** and the production constructor refuses on it — that is a third live barrier today, but it also makes any arm touching it machine-dependent; scrub it with the repo's landed idiom (`registry_publish_test.go:1695`). (d) **`validatePublishOrigin` checks SCHEME before LOOPBACK**, so `http://127.0.0.1:1` is refused for its scheme and never reaches the loopback branch — a test row using `http` would be green while proving nothing; use `https://127.0.0.1:1`. (e) `go run` does **NOT** propagate exit status (measured: rc=1 for a program exiting 3; the built binary gives rc=3), so any exit-code assertion must build the binary. (f) the `reconcile --probe` branch's call is executed by no test — doing so would mean a non-loopback GET; its argument construction is measured, its behaviour is covered by the landed loopback reconcile tests, and the limitation is stated in the test. (g) `pi` remains BARRED for this lane (Mark, attended 2026-08-06); iter-67 ran planner + executor on `opus`, evaluator on `sonnet`, FLAGGED.

    **NO PUBLISH HAS OCCURRED.** No non-loopback registry request was made; `AILANG_REGISTRY_API_KEY` was never set, read or passed; **`world/core@0.1.0` remains UNCLAIMED**.

    **Prior head text follows.** ~~**[IN-SPRINT] — THIS ITEM IS NOT BLOCKED ON MARK AND HAS NOT BEEN SINCE 2026-08-05. `SM.D`'S ENTRYPOINT DOES NOT EXIST, AND BUILDING IT IS HEADLESS-ROUTABLE WORK — THE NEXT UNIT OF WORK, GATED ON NOTHING** (attended session 2026-08-10).~~
   **THE CORRECTION.** Iter-65/66 recorded `SM.D` as *"attended-only, blocked on `8/OD-1`"* and concluded the item had no headless-routable milestone. Measured attended, four ways: `cmd/ailang-worldd`'s verbs are `serve · health · head · world · object · log · registry · commit · help` and `grep Publish|Approve` over the whole command returns **zero** non-test hits (`registry` is a read-only GET, `cli.go:164`); `scripts/` holds five scripts, none publishing; the runbook's **three ```bash blocks are ALL in Stage A**, so Stage B steps 5–7 are prose, not a procedure; and the publish machinery is a **library** — `RegistryPublishHandler`/`RegistryPublishConfig` need a caller to supply `RegistryOrigin`/`ValidatorOrigin`, and **no non-test caller exists**. So the blocker was never a human decision. `8/OD-1` was ANSWERED 2026-08-05 and owed no artifact after it.
   **WHAT `SM.D` ACTUALLY IS:** a `main`/verb wiring store → session → the `ApprovalRequestV1`/`ApprovalDecisionV1` pair → `Session.Invoke`, tested against `httptest` exactly like every existing caller. **Headless may BUILD it; headless must NEVER RUN it** — and the design already enforces that, since the approval pair is minted attended and `Session.Invoke` refuses before the credential loads and before any POST. **`pi` stays BARRED from this executor lane** (Mark, attended 2026-08-06).
   **THE ABSENCE IS THE SAFETY PROPERTY.** `registry_publish.go:396-399` is a production constructor demanding `https` and **refusing loopback**, while every caller today is a test against `httptest`. Nothing in this repo can publish, so no headless process can trip it; the entrypoint is the deliberate relaxation, which is why it is attended-class.
   **STAGE A IS GREEN AND REVIEWED** (attended, pinned `/tmp/ailang-v0300/ailang` = `AILANG v0.30.0`): projection reproduces byte-identically (`git status --porcelain packages/` empty), readiness gate 9/9 with canonical JSON equal to the golden byte-for-byte, `verify_ail.sh` rc=0 (4 identities / 14 tests), `verify_go.sh` rc=0 (28 `ok`, 0 FAIL, race control 2). Packet: `world/core@0.1.0` · 4 exports · effects `[]` · 5773 bytes · `a32806a069bbe2e79…` / `5ea15858fddc8f8eb…` / `d16cc88270ff4c4ea…`. **NO PUBLISH OCCURRED; `world/core@0.1.0` remains unclaimed.**
   **TWO RUNBOOK DEFECTS ARE `SM.D` INPUTS:** (a) step 4 tells the human to confirm three digests *"against the gate's own output"* and the gate emits **none** (0 hits for `sha256:[0-9a-f]{64}` in the full gate log, control 3 in the golden) — it asserts only agreement; the dry-run truncates to **17 hex chars** (`hash[:24]+"..."`), so the human-performable check is 68 bits while the real 256-bit comparison happens mechanically at steps 7 and 9. Rewrite it to say so. (b) Stage B carries **no commands at all**, which is how (a) survived.
   **CARRY CLOSED:** the publisher's success marker is now **OBSERVED**, not read from upstream source — attended `--dry-run` with `env -u AILANG_REGISTRY_API_KEY` prints `⚠ Dry run complete. Tarball ready but not uploaded.`
   **Prior head text follows.** ~~**[IN-SPRINT] — `SM.C` LANDED 2026-08-08 (iter-65), PR #52 → squash `0cd00eb`, dev CI green BOTH jobs SHA-addressed and step-log verified on the merge commit (`ailang-code verify gate` 11/11 · `go host build + test gate` 13/13, `failed=0`) during a 0-incident window. Evaluator `sonnet` **93/100, ZERO BLOCKING**. `AC13`, `AC14`, `AC15`, `AC16`, `AC16a`, `AC16b`, `AC17` all discharged. THE NEXT UNIT OF WORK IS `SM.D`, AND IT IS BLOCKED ON `8/OD-1` — an ATTENDED human decision. `SM.D` must never run headless or in CI, so THIS ITEM HAS NO HEADLESS-ROUTABLE MILESTONE LEFT: the loop should pick a different queue item until Mark rules on `8/OD-1`.** **WHAT `SM.C` SETTLED:** an indeterminate publish is resolved by READING the public bucket (`host/broker/registry_reconcile.go`, 525 lines inherited from dead iteration 64 + a 1692-line suite). Single network verb `GET`; four receipt states `succeeded-reconciled` / `conflict` / `not-published` / `probe-unavailable`, of which only the first three are resolutions — `probe-unavailable` is an explicit refusal to decide. Absence is believed ONLY on the measured GCS `NoSuchKey` document **decoded as XML** (not string-matched), and ONLY when a same-pass known-positive control returns 200-with-JSON from the TARGET'S OWN key-space; `metadataObjectURL` builds both target and control so no caller can split them. **READ BEFORE STARTING `SM.D`:** (a) **`8/OD-2` IS NOW ROUTED** — upstream namespace-auth ask filed as `sunholo-data/ailang#633` with every premise measured first-party (`cmd/registry-validator/main.go:177` defers namespace auth; ONE shared `REGISTRY_API_KEY` at `:54,:76,:106-111` never checked against the vendor prefix parsed at `:159`; immutable 409 at Step 4; **zero** owner/vendor/principal/scope JSON paths in the public index against a **396**-path `name|version` control; exactly one vendor `sunholo` has ever published). `world/` remains CONVENTION-ONLY and the runbook says so. (b) **The attended runbook is `docs/SELF_MOD_PUBLISH.md`** and it stops at readiness by default; its gate `host/runbook` reads the commands OUT OF the doc, fails loudly on a zero extraction, and asserts Stage A contains no live publish. (c) **`unreachable_host` in the `AC16a` arm-(iii) table is a COVERAGE BYSTANDER, not a guard** — the refusing branch is `C1` (not `C4`, as the original comment claimed): closing the origin kills the CONTROL request too and the control is examined first. Verified twice first-party — under `MUT-SM-PROBE-NO-CONTROL` its two siblings red and it PASSES. The comment now says this explicitly; do not re-derive it. (d) **`AC13`'s pre-existing guard was satisfiable by a recovery that DOES dispatch** — it passes the REAL handler, which refuses a malformed request long before its own dispatch counter moves, and stayed GREEN under `MUT-SM-RECOVERY-DISPATCH`. The new test counts the CALL. (e) `MaxBodyBytes` truncation is exercised only through equivalent short bodies; the `io.LimitReader` mechanism itself is inferred, not measured — cheap to close. (f) the two load-flaky wall-clock tests in `cmd/ailang-worldd` did NOT fire this iteration, but the rule stands: read WHICH test failed, never the exit code. (g) **`host/boundary` pins `wantFileCount = 1`** — adding any `.go` file there reds `TestBoundaryASTWriteGuard`; this bit iter-65 and the fix is to put the new gate in its own package, never to relax the pin. (h) executor lane: `pi` remains BARRED for publish-capable milestones (Mark, attended 2026-08-06) and codex was measured quota-dry first-party (`rc=1`, resets 11:24), so iter-65 ran `opus` via the Agent tool, FLAGGED.~~ **Prior head text follows.** ~~**[IN-SPRINT] — `SM.B2b` LANDED 2026-08-08 (iter-63), PR #51 → squash `abb3a3d`, dev CI green BOTH jobs SHA-addressed and step-log verified on the merge commit (`ailang-code verify gate` 11/11 · `go host build + test gate` 13/13, `failed=0`) during a 0-incident window. The next unit of work is milestone `SM.C`, gated on nothing.** **WHAT `SM.B2b` SETTLED:** an attended approval now **binds bytes, not a name**, and is **spent exactly once, durably**. `Session.Invoke` traverses `payload.approvalRef` → landed `ApprovalDecisionV1` → its `ApprovalRequestV1` → the canonical scope, and refuses BEFORE the credential is loaded and BEFORE any POST; single use is enforced by `approval_claims`' PRIMARY KEY, not by in-memory budget. `AC8`, `AC9`, `AC9a`, `AC9b`, `AC9c` are discharged non-vacuously — AC8 re-presents on a FRESH session whose budget is asserted `== PublishCost` before invoking; AC9a/AC9c use a FILE-BACKED store (`:memory:` is per-connection, so a reopen criterion would have passed against an empty database); AC9b parks both racers on a start barrier and runs under `-race`; every POST counter is asserted NUMERICALLY (`== 1`), never as "no second success". **`NB-2` FROM SM.B1 IS CLOSED**: the concurrent-collision path now returns `ErrApprovalAlreadyConsumed` rather than a wrapped PRIMARY KEY error, decided by asking the DATABASE rather than matching driver error text. **READ BEFORE STARTING `SM.C`:** (a) **the design doc contradicts itself in ONE paragraph and the contradiction is still there** — §"HumanApproval" says the request wire carries *"the publish effect"* AND that `HumanHandler` is reused unchanged; `approve.go:111` writes `Effect: req.Effect` inside `case EffectHumanApprove:`, so it is always `"Human.Approve"`. SM.B2b implemented the doc's other, achievable clause (the binding is a canonical value of the existing `Scope` field) and added `effect` as the first frozen scope term. The doc's prose is NOT updated; do not re-derive the impossible version. (b) **`journal.go`'s PK-collision fallback is correct but UNREACHABLE** (`db.SetMaxOpenConns(1)` + a cross-process writer lock taken before the handle opens); it has no test and none is claimed — `store.go:293` is the falsifier if that connection cap is ever relaxed. (c) **`NB-2` from SM.B2a is still carried, not closed**: Decision 3's *"refuses redirects to a different origin"* is only PARTIALLY discharged — the POST happens inside the pinned `ailang` child, so this process has no `CheckRedirect` hook and only the origin handed to the child is validated. SM.C/SM.D own it. (d) `manifestRef`/`compilerVersion`/`compilerSHA256` remain carried and request-hash-bound but **never cross-checked** against `ailang.toml` or the pinned binary. (e) the publisher's success/failure markers were read from upstream SOURCE (`e37b370:cmd/ailang/pkg_publish.go`), never from observed output; only an attended `--dry-run` in SM.D can close it. (f) **`cmd/ailang-worldd/TestHandlerTimeoutKillsTheWholeProcessGroup` and `TestCLIRealSubprocessEpisode` are FLAKY under load** — measured at 2/5 and 1/4 respectively during this iteration, both wall-clock-bound and both unrelated to approval code. One of them redded a controller mutation run in exactly the predicted direction; **read WHICH test failed, never the exit code alone.** Worth queueing as its own item. (g) **The `pi` executor lane remains BARRED for publish-capable milestones** (Mark, attended 2026-08-06); codex was measured quota-dry (`rc=1`, resets Aug 8 11:24), so iter-63 ran `opus` via the Agent tool, FLAGGED. **Prior head text follows.** ~~**[IN-SPRINT] — `SM.B2a` LANDED 2026-08-07 (iter-62), PR #50 → squash `3fd889f`, dev CI green BOTH jobs SHA-addressed and step-log verified on the merge commit. The next unit of work is milestone `SM.B2b`, gated on nothing.** **WHAT `SM.B2a` SETTLED:** the brokered publish handler (`host/broker/registry_publish.go`), the de-ambient credential provider (`credential.go`), and the typed indeterminate dispatch path exist; `AC7`, `AC10` and `AC11` are discharged non-vacuously. **WHAT IT FOUND, AND THE REASON THE REPLACED AC EARNED ITS KEEP:** iter-52's planner judged the doc's `AC10` vacuous (*"all non-publish subprocesses observe it unset"* is satisfiable by launching zero subprocesses) and replaced it with one that must re-derive the site count by command in-run, print it, drive every site, and `t.Fatal` on a zero-length enumeration. Executing that literally measured **two of the five production subprocess sites leaking a live, irreversible-publish credential** — `archive.probeVersion` set no `cmd.Env` at all, `replay.runPinnedTransition` set `Dir`/`Stdout`/`Stderr` but not `Env`. Both fixed; `host/childenv` now holds the variable list once so four packages cannot drift. Verified first-party at base `0c47667`, not inherited from the executor's report. **READ BEFORE STARTING `SM.B2b`:** (a) it owns `AC8` (dispatch half) + `AC9`/`AC9a`/`AC9b`/`AC9c` — attended-stamp binding and single-use approval consumption; `AC9b`'s concurrent-race criterion is the one that must close SM.B1's carried **NB-2** (two callers both passing the `SELECT EXISTS` pre-check surface a wrapped PRIMARY KEY constraint error rather than `ErrApprovalAlreadyConsumed`). (b) **`SM.B2a` did NOT validate approvals** — `AppendClaimedEffectIntent` is wired for `Registry.Publish` but no approval *checking* exists yet; that is `SM.B2b`'s whole job. (c) **CARRIED FROM THE JUDGE (`NB-2`, non-blocking):** Decision 3's *"refuses redirects to a different origin"* is only PARTIALLY discharged and cannot be completed here — the POST happens inside the pinned `ailang` child, so this process has no `net/http` import and no `CheckRedirect` hook; only the origin handed to the child is validated. Route to SM.C/SM.D, and do not inherit it as green. (d) `manifestRef`, `compilerVersion` and `compilerSHA256` are **carried and request-hash-bound but never cross-checked** against `ailang.toml` or the pinned binary. (e) The publisher's success/failure classification markers were read from upstream **source** (`e37b370:cmd/ailang/pkg_publish.go`), never from observed output, because the safety fence forbids a live run — a `--dry-run` observation in an attended SM.D step is the only thing that can close it. (f) **The `pi` executor lane remains BARRED for publish-capable milestones** (Mark, attended 2026-08-06); codex was quota-dry until Aug 8 11:24, so iter-62 ran `opus` via the Agent tool, FLAGGED. **Prior head text follows.** ~~**[IN-SPRINT] — THE QUEUE HEAD AGAIN AS OF 2026-08-07 (iter-61): item 10 is COMPLETE, so `SM.B2a` is the next unit of work and it is gated on nothing.**~~ Two things to carry into it. **(a) `8/OD-1` is ANSWERED and is now REGISTERED in the charter's OD table** — Mark's attended 2026-08-05 approval landed on an integer that had never been registered and collided three ways; that is fixed, and the publish POLICY is approved (though not the exact-bytes attended stamp, which cannot exist until the ready packet does). **(b) `SM.B2a` is the standing EXCEPTION to the `pi:deepseek-v4-flash-0731` executor lane** — publish-capable code does not default to an unsandboxed executor (Mark, 2026-08-06, attended). **Item 10 has also now discharged the hazard it was promoted for**: the boundary gate no longer writes the tree it guards, so `SM.B2a` will not have to distinguish a crash-residue false positive from a genuine network-boundary violation. **Prior head text follows.** ~~**[IN-SPRINT] — UNCHANGED THIS ITERATION; `SM.B2a` IS STILL THE NEXT MILESTONE BUT NO LONGER THE QUEUE HEAD. Item 10 was promoted ahead of it 2026-08-06 (iter-56)**~~ on a measurement iter-55's own row did not record: the boundary gate's `defer`-based restore does not survive SIGKILL or Go's `-test.timeout` panic, so a killed run leaves a forbidden import in a production source **permanently** — invisible to `go build` (rc=0) and reddening the boundary gate itself, which then accuses an innocent file of a **network-boundary violation**. That failure mode is indistinguishable from `SM.B2a` genuinely violating the boundary, during the one milestone that adds network code, which is why it is ordered first. See item 10. **Prior head text follows, still accurate for the milestone itself.** ~~**[IN-SPRINT] — `AC12` REPAIRED 2026-08-06 (iter-55), PR #45 → squash `1761a9c`, dev CI green BOTH jobs SHA-addressed on the merge commit. The next unit of work is STILL milestone `SM.B2a`, gated on nothing — it now lands against a boundary gate with teeth.**~~ **THE CARRY-FORWARD THE LAST TWO ITERATIONS WROTE DOWN WAS DISCHARGED EARLY, AND IT FOUND A HOLE.** Both iter-53 and iter-54 recorded that `AC12`'s *"network confined to `host/broker`"* control is VACUOUS until SM.B2a and must be re-asserted there. Re-asserting it at the boundary — BEFORE the network code exists — showed the gate was weaker than vacuous in a second, unrecorded way: the loopback exception is true of **exactly one** protected group and a single shared `forbiddenImportPrefixes` list granted it to **all three**. Measured, mutations confirmed landed by sha256, restores byte-identical: bare `net/http` blank-imported into `host/store/store.go` → **rc=0 PASS**, into `host/replay/replay.go` → **rc=0 PASS**, while the `net/http/httputil` control correctly REDs. **Every protected group's mutation was `httputil`** — so the gate had only ever been tested against a mutation shaped to itself, iteration 54's own spine arriving inside the gate iteration 53 landed. The exemption was also **UNFORCED**, not merely unnecessary: baseline `net/http` presence per dependency closure is `host/store` **0** (160 deps), `host/replay` **0** (162), `cmd/ailang-worldd` **1** (233), control `host/hashref` **0** — only `cmd/ailang-worldd` ever needed it, exactly as its own code comment says. Repaired with a per-group `extraForbidden` field plus `TestBareNetHTTPExemptionIsPerGroup`, which reds if the asymmetry is collapsed back into one global list (proven non-vacuous by setting `host/store`'s entry to `nil` → RED naming the group, restore byte-identical). Post-repair: both threat arms RED, `httputil` still RED, `cmd/ailang-worldd` still PASS, pristine green. **WHAT REMAINS FOR `SM.B2a` TO RE-ASSERT:** the positive half — that network code, once it EXISTS in `host/broker`, is genuinely permitted there — is still unproven, because `host/broker` has zero `net/http` deps at HEAD. The gate's green control asserts only that the broker's dependency closure is NON-EMPTY (`:281`), which is true of any Go package. **SEPARATE, PRE-EXISTING, AND NOW QUEUED AS ITEM 10:** this same gate mutates three other packages' production sources in the LIVE tree while `go test ./...` builds them concurrently — measured on pristine `dev`, and it red-lit `TestCLIRealSubprocessEpisode` once during this iteration. Prior text follows.] ~~**[IN-SPRINT] — `SM.B1` LANDED 2026-08-05 (iter-54), PR #43 → squash `1856bfb`, dev CI green BOTH jobs SHA-addressed on the merge commit. The next unit of work is milestone `SM.B2a`, gated on nothing.** **What SM.B1 SETTLED:** `approval_claims` + an atomic `AppendClaimedEffectIntent` exist, schema is at version **2**, and **`DD-3` is closed loudly** — `enforceSchemaVersion`'s bare `return nil` became `*LegacySchemaVersionError`, so a store left at `user_version = 1` can no longer open successfully and silently skip `schemaSQL` (which would have surfaced `approval_claims` as absent at the moment of the irreversible publish). All three independent fixtures moved in ONE commit with `schemaV1SQL` byte-unchanged; the DDL gate redded mid-milestone exactly as designed and the recorded RED list is in the log. **What SM.B1 newly OPENED — the milestone-gating ledger check was VACUOUS AS DELIVERED, and the executor's own mutation could not see it.** `TestSchemaVersionLedgerIsIndependent` greps its own source; its two NEGATIVE needles were split so they would not match their own check-lines, while its POSITIVE needle was one literal that did — so it passed whatever the declaration said. Measured: `var schemaV2SQL = string(schemaSQL)` (the ledger becoming the file it exists to attest) returned **`ok 0.290s`**. The executor's `MUT-SM-V2-LEDGER-DERIVED` redded only because it used the bare form the negative needle was written to catch — **a mutation shaped to the check tests the check, not the threat**. Repaired by anchoring to `^` plus a semantic `schemaV2SQL == schemaSQL` backstop; both forms now RED, and the judge added a third (`schemaV1SQL + ""`) that also REDs. **Read before starting `SM.B2a`:** `AC12`'s *"network confined to `host/broker`"* control has been **VACUOUS since SM.A** (`host/broker` has zero `net/http` deps) and **SM.B2a is the milestone that makes it real** — re-assert it there, never inherit it as green. **Evaluator `sonnet` PASS 91/100, ZERO blocking**, three non-blocking findings all carried: **NB-1** AC-B1.4's pre-bump negative arm is measured EXTERNALLY (against `origin/dev`) rather than embedded in the test, so a future reader cannot see it from the test file alone; **NB-2** under genuine in-process concurrency two callers can both pass the `SELECT EXISTS` pre-check and the loser fails on the `approval_claims` PRIMARY KEY — correct (no double-consumption) but surfacing a wrapped constraint error rather than `ErrApprovalAlreadyConsumed`; the judge MEASURED this (`… UNIQUE constraint failed: approval_claims.approval_ref (1555) …`), so it is evidence rather than inference, and **`AC9b` in SM.B2b is the criterion that must close it**; **NB-3** the doc carried no SM.B1-era verification rows — **discharged this iteration** as row `V-S`. **Also landed this iteration, outside item 8:** SM.A had committed a 15.7 MB `ailang-worldd` Mach-O that five independent checks passed; removed and gated (PR #42 → `e24a6f0`). Prior text follows.] ~~**[IN-SPRINT] — `SM.A` LANDED 2026-08-05 (iter-53), PR #41 → squash `13315da`, CI green BOTH jobs SHA-addressed on the merge commit and step-log verified. The next unit of work is milestone `SM.B1`, gated on nothing.** **What SM.A SETTLED:** `AC6`'s three-arm cross-check AGREES on darwin/arm64 **and** linux/amd64 (tarball `5472 = 5472`; CI printed `✓ compiler pinned by exact bytes: AILANG v0.30.0 on Linux/x86_64`, `9/9` steps), each arm mutation-proven able to red — so the cross-toolchain tarball risk `DD-1` raised is **CLOSED**, and `DD-4`'s third-LEG landed with `4/4` identities and `14` named tests intact. **What SM.A newly OPENED — `DD-7`, and its second half is queue item 9's problem now:** a byte-exact compiler pin is **platform-specific** (darwin `e9746fef…` / linux `1e594d15…`), and separately **CI job 1 has been verifying `.ail` against `AILANG v0.33.0` since `latest` moved on 2026-08-04** — measured in the step log at `af0c3b4` (run `30993399332`) against job 2's `v0.30.0` in the same run. The package leg routes around it with its own pinned install via `WORLD_PKG_AILANG_BIN`, so **nothing is blocked**, but the two ORIGINAL legs remain unpinned and item 9's "latent, not active" grade is now false. **Read before starting `SM.B1`:** it must be ONE commit (`schema.sql` + both `store.go` constants + the stale-version acceptance policy + all three independent fixtures — splitting it lands a red DDL gate); `DD-2`'s ~3× blast radius and `DD-3`/`8/OD-3` below are binding. **`AC12`'s limits, carried forward so SM.B2a does not inherit a false sense of coverage:** the *"network confined to `host/broker`"* green control is **VACUOUS today** (`host/broker` has **zero** `net/http` deps — network arrives WITH SM.B2a, which is precisely when that control becomes real and should be re-asserted), the guard is **source-level not closure-level**, `cmd/ailang-worldd`'s inherited `net/http` is **loopback IPC** and not egress, and `host/registry` in the forbidden list is the *interpreter epoch* registry — a **name collision** that will produce a false positive the moment anything legitimately needs epoch metadata. **Evaluator `sonnet` PASS 87/100, ZERO blocking**, and it added a THIRD carried limit: `AC5`'s smoke coverage is enforced **implicitly** — dropping a module changes the smoke's output — rather than by an explicit import-coverage manifest. Correct outcome, weaker instrumentation; worth strengthening when SM.C revisits the fixtures. Prior text follows.] ~~[IN-SPRINT] — SPRINT-PLANNED 2026-08-05 (iter-52). The next unit of work is milestone `SM.A`, gated on nothing.~~ Plan: `.ailang/state/sprints/w-self-mod-vertical.{plan.json,handoff.md}` (planner `opus`, lane derived VERBATIM as `opus fail-closed:env-pin`; gitignored, as all 18 prior sprint artifacts are). **The planner re-scoped 4 milestones into 6** — `SM.A · SM.B1 · SM.B2a · SM.B2b · SM.C · SM.D` — because SM.B as designed prices at ~2,300–2,700 LOC against a measured maximum single landed commit of **751** insertions (n=5: 457, 414, 210, 751, 698; median 457). It stays **ONE queue item** (precedent: 4d and 4e each ran multi-milestone across 4–5 iterations without becoming a second row); the split is internal. **`AC12` moved from SM.B into SM.A** — a boundary guard that lands alongside the code it constrains has never been observed rejecting that code. **`8/OD-1` IS RATIFIED** by Mark on `#32` @ `2026-08-05T08:25:00Z` (*"Approved publish of world/ in ailang extensions - go. Credential is on your machine for this."*): the POLICY is approved and SM.D is no longer blocked on a human answer — but this is **not** the exact-bytes attended stamp SM.D describes, and it cannot be, because the ready packet does not exist until SM.A builds it. **An authorization is not an attendance.** SM.D stays attended-only, never headless, never in CI, stop-at-readiness by default. **THREE PLANNER FINDINGS, ALL REPRODUCED FIRST-PARTY BY THE CONTROLLER, THAT CHANGE THE WORK.** **`DD-1` (blocking for SM.A as designed):** Decision 3's *"small library extraction of v0.30.0 package hashing logic"* is **impossible** — World is `module github.com/sunholo-data/ailang-world`, upstream is `module github.com/sunholo-data/ailang`, and the hashing lives in `internal/pkg/`, which Go's internal rule forbids across modules. The CLI is no fallback (`pkg_publish.go:110-112` prints `hash[:24]+"..."`; the tarball bytes are never persisted). `AC6` therefore needs a **re-implementation** (`host/pkgproj`) with a mandatory hard-failing 24-char cross-check, plus a newly named risk: the tarball hash rides `compress/gzip` and the two modules declare different Go versions. **`DD-3` (a silent runtime hole the version bump CREATES):** `host/store/store.go` ends its version ladder in a bare `return nil` at `:354`, so once `currentSchemaVersion` goes 1→2 a store still at `user_version = 1` matches **no branch**, opens successfully, and **never executes `schemaSQL`** — `approval_claims` would be found absent at the moment of the irreversible publish. Raised as **`8/OD-3`** and answered from already-ratified text (`4d/OD-3` alt-1: fail LOUD on an unsupported **or un-upgraded** store); non-blocking. **`DD-4`:** `scripts/verify_ail.sh:160`/`:190` are exact equalities (`EXACT_TOTAL_VERIFIED=4`, `EXACT_TOTAL_TESTS = 14`), so the package gate must be a third **LEG**, never a new **root** — adding `packages/` to `ROOTS` doubles the identities to 8 and reds the repo's primary gate for a reason unrelated to the code under test. **THE DDL BLAST RADIUS IS ~3× WHAT THE DOC'S CONFLICT SURFACE NAMES**, and all of it lands in **SM.B1's single commit** (splitting it lands a red gate): the doc names only `host/store/schema_version_test.go`, while `host/store/journal_test.go:714` holds a SECOND independent fixture `canonicalTableDDL` (**7** hardcoded tables behind `requireExactTableNames` at `:778`, carrying its own *"must not be derived from schemaSQL … or the database under test"* comment) that the doc never mentions; `schema_version_test.go:16` `frozenFutureSchemaVersion = 2` **collides** with the new current version; and `store.go:316` writes a **literal** `PRAGMA user_version = 1` that would trip `freshInitTx`'s own drift check at `:325-326`. **FIVE ACs JUDGED VACUOUS AND REPLACED IN THE PLAN** (`AC13`, `AC17`, `AC10`'s 2nd clause, `AC19`, `AC1` — each by the mission's own test: *would this pass identically if the thing it protects did not work?*), plus 5 planner-authored mutations → **36 total**. Prior text follows.] ~~[NEXT] — DOC LANDED 2026-08-05 (iter-51); NOT YET SPRINT-PLANNED. The next unit of work is the sprint-planner run, and it is gated on nothing.~~ Doc: `design_docs/planned/w-self-mod-vertical.md` (839 lines; designer `codex:gpt-5.6-sol`, rotation slot 2; PR #40 → squash `269f1fe`, dev CI green both jobs SHA-addressed on the merge commit). Milestones **SM.A–SM.D**, **4–5 d** (raised from 3–4 by the round-1 quorum revision; SM.B alone is 2–2.5 d). **THE VERIFY-FIRST CLAUSE WAS DISCHARGED FIRST-PARTY BEFORE THE DESIGNER WAS SPAWNED, AND IT REFRAMED THE ITEM.** There is no vendor-namespace claim operation to perform: `cmd/registry-validator/main.go:177` @ the pinned `e37b370` reads verbatim `// Step 5: Namespace auth — deferred (accept all publishers for now)`; auth is ONE optional shared secret whose header the client omits when unset; and a live four-arm dry-run at the pinned binary accepts `world/probe`, `someoneelse/probe` and `sunholo/probe` alike (rc=0 each) against a firing known-positive control (`novendor` → rc=1, `must be vendor/name format`). **`world/` is a string World writes, not a namespace World holds** — recorded here so no later reader re-derives ownership from Mark's wording, which said "verified unclaimed" and was true. Census: 40 packages, vendor histogram `sunholo` 40 and nothing else, `world/` 0 with its control firing — World would be the registry's first second vendor ever. Publish is immutable (409) and unrecallable by the publisher, and **`AILANG_REGISTRY_API_KEY` is AMBIENT in this loop's own tool shells**; compose those and any process inheriting this environment can write irreversibly under any vendor string. That is the clause-3 surface, and it is what makes the brokered/receipted framing load-bearing rather than decorative. **TWO THINGS THE PLANNER PRICES FIRST — both carried, neither cleared.** (i) The design NEEDS a `schema.sql` change (`approval_claims` + an atomic `AppendClaimedEffectIntent`, because `host/store/store.go:601` `SetRegistryHead` is a blind upsert with no expected-previous and is the only set-shaped store API, so there is no claim-if-unused primitive) — and the landed `w-ddl-gate-teeth` DDL gate reds on **any** `schema.sql` edit *by design*, so its fixture update belongs inside the same milestone rather than after it. (ii) Whether 4–5 d is still one queue item or wants splitting at the SM.B boundary. Quorum round 2 ran ONE reviewer (`gpt5-6-sol` absent: `budget`) and round 3 was the controller carve-out, so neither question has had two pairs of eyes. **Open decisions: `8/OD-1`** — the attended human stamp for the irreversible first publish; controller default **do not publish**, and it blocks **SM.D only**, so SM.A–SM.C are routable today. **`8/OD-2`** (upstream namespace authorization) is non-blocking. Prior text follows.] ~~[NEXT] — UNPARKED 2026-08-05 (iter-50). The park condition was "until 4 lands"; item 4 `w-effect-broker-m3` completed at iter-35. Mark re-scoped this row attended on 2026-08-04 (`de80792`). Routes to design-doc-creator; VERIFY-FIRST binding at pick.~~ ~~[PARKED until 4 lands]~~ **w-self-mod-vertical** · clause-7 · one World extension shipped through **[NAMING + FLOW, Mark 2026-08-04: World extensions publish to the PUBLIC AILANG
   registry under the `world/` vendor namespace** (verified unclaimed 2026-08-04; explorer:
   ailang.sunholo.com/docs/packages/explorer). The clause-7 publish is itself a BROKERED,
   RECEIPTED effect (capability + budget + record) — World's self-modifications become public,
   hash-verified artifacts with a full evidence chain. Local-first holds: publish is an outward
   effect-handler interaction; consumption is version-pinned + hash-verified install; the core
   never depends on the registry at runtime. Item-8's design doc VERIFIES publish-auth/vendor-
   registration mechanics at pick (VERIFY-FIRST) — `ailang publish` reads ailang.toml, uploads
   via the validation service.]**
   World's own pipeline end-to-end · ~1–2d
9. [**[LANDED — ITEM COMPLETE 2026-08-10 (iter-68). `VL.A` shipped all three routable pieces: PR #56 → squash `9789b87`, dev CI green BOTH jobs, SHA-addressed, `present=2` = expected 2, step-log verified on the merge commit, 0-incident window. Evaluator `sonnet` **38/100 FAIL round 1 — and the judge was RIGHT**; its blocking finding (a hardcoded rig path that red CI while every local gate was rc=0) is the iteration's spine. 10 named mutations, 10 killed.**
   **THE DEFECT THIS CLOSED, MEASURED AT BASE:** a shim reporting `AILANG v0.33.0-105-g38e119db1-dirty` drove the REAL gate to `rc=0` / `verify gate PASSED` — the primary `.ail` gate passed exactly the `-dirty` build CLAUDE.md forbids. **(i)** the warn-only compare is now a HARD **is-a-release** assertion with five refusal branches, each emitting a stable reason CODE (the branches FUNNEL — `NOT_A_RELEASE` is a catch-all, so an `rc`-keyed mutation table scores 3 of 8 falsely SURVIVED); the shape deliberately admits a pre-release identifier because upstream published `v0.24.1-rc1` with `isPrerelease: false`, so strict `^vX.Y.Z$` had a measured 1-in-64 chance of redding CI on a genuine release. **(ii) `9/CF-A-1` DISCHARGED** — `scripts/testdata/ailang_version_shim.sh` (0755, delegating) is committed, plus an in-script known-positive control mirroring `verify_go.sh:91-100`/`:62-69`. **(iii) `9/CF-A-2` DISCHARGED** — the always-firing DRIFT warning became an expected-release MEMBERSHIP notice; equality against ONE observation was measured to MOVE the defect rather than fix it (quiet in CI, firing on every LOCAL run with the false text `moved from 'v0.33.0' to 'v0.30.0'`), because legs 1-2 resolve a different release per lane: CI job 1 exports no `AILANG_BIN` (`ci.yml:87` exports only `WORLD_PKG_AILANG_BIN`) while operators export the v0.30.0 pin. **NO CI CONFIG CHANGE** (clause a): 0 paths under `.github/`, control 1.
   **OPEN CARRY-FORWARDS (judge non-blocking, accepted, OUT of the item):** `9/CF-VLA-1` `MUT-VL-SHIM-NODELEGATE`'s inverse arm is not clean — three tests exercise delegation, so the mutation table's 1:1 `killed_by` mapping overstates precision. `9/CF-VLA-2` `TestInScriptControl`/`TestEmptyExpectedReleaseSetFailsLoudly` write temp mutant scripts into `scripts/` and clean up with `defer os.Remove`, which is not crash-safe; residue would be an untracked dotfile, not a mutated tracked file, so no gate reds — low severity, wants `t.Cleanup` + a `TestMain` sweep if the pattern recurs. `9/CF-VLA-3` a synthetic `v0.33.0-105` (digits after one hyphen, no `-g<hash>`) is accepted as a release; `git describe` never produces that shape. **NEW ASK `9/OD-11`** — may a milestone add a Z3 install step to CI job 2? See the OD registry.
   **Prior head text follows.**]** ~~[**`9/OD-10` RATIFIED 2026-08-10 (Mark, attended) — ACCEPT. THE REMAINING HALF IS NO LONGER HUMAN-GATED; IT IS HEADLESS-ROUTABLE AND IS THE NEXT UNIT OF WORK FOR THIS ITEM.**
   **THE RULING AND ITS THREE-CLAUSE SCOPE** (full text in the OD registry): legs 1-2 **TRACK RELEASED** versions — `ci.yml:25` keeps `releases/latest`, **no CI config change**; a **non-release / `-dirty` build is STILL REFUSED**; `:64` **stays pinned** to v0.30.0 because it byte-compares v0.30.0-generated replay goldens. Rationale in Mark's words: *"we should accept upstream releases as they will respond to ailang world requests currently"* — World files upstream asks (e.g. `ailang#633`) and upstream ships fixes, so pinning legs 1-2 would mean World's own upstream fixes never reach its primary gate.
   **THE RULING INVERTS ITER-66'S CARRY.** Iter-66 held the hard `verify_ail.sh` assertion as unsafe headless because it would red on the next upstream release. Under clause (b) the assertion is re-shaped to assert ***is a release*** rather than ***is v0.30.0***, and in that shape it **cannot** red on an upstream release — only on a dev build. So ACCEPT did not close this item empty; it made the remainder routable.
   **THREE ROUTABLE PIECES, ALL HEADLESS:** (i) the re-shaped is-a-release assertion on `verify_ail.sh` legs 1-2; (ii) **`9/CF-A-1`** — commit the shim fixture that proves the version assertions fire (iter-66's proving shim was ad hoc and is not in the repo, so this repo's own rule-3a culture is satisfied by a control that did not survive the iteration); (iii) **`9/CF-A-2` (NEW, from the ruling itself)** — the iter-66 DRIFT warning compares against the v0.30.0 pin, so under ACCEPT it now fires on **every CI run forever**, and **a warning that always fires is not a signal**. Re-shape it to announce the resolved release and warn only on a non-release build, or on change-since-last-observed.
   **Prior head text follows.** ~~**[HALF-LANDED 2026-08-10 (iter-66) — the CHEAP, HEADLESS-SAFE half shipped; the HUMAN-GATED
   half is `9/OD-10` and is the whole remaining item.** PR #54. `verify_ail.sh` now ANNOUNCES the
   binary legs 1-2 resolved and WARNS (never fails) on drift from the v0.30.0 pin; proven
   non-vacuous by two REAL arms rather than by mutation, because an observability feature cannot be
   proven by a red — pinned arm is rc=0 with output identical to the pristine baseline except the
   one added line, bare-PATH arm announces a different path AND version AND fires the warning.
   Re-measured at HEAD and WORSE than iter-53 recorded: the rig's PATH `ailang` is
   **`v0.33.0-70-g1677fcff9-dirty`**, not `v0.33.0-1-gdd68e074`.
   **AND THE CONTROL WRITTEN TO JUSTIFY THE EXACT-TOKEN COMPARE FOUND A LIVE FALSE-GREEN IN THE
   SIBLING GATE:** `verify_go.sh:29`'s **hard** anti-false-green assertion was `grep -q 'v0.30.0'`,
   which the string `v0.30.0-205-g54d6bd191-dirty` **SATISFIES** — so the guard whose entire purpose
   is rejecting an unpinned compiler admitted a 205-commit dirty dev build. Verified on the REAL
   script (rule 3k) with an executable shim: pristine printed its announce and PROCEEDED; tightened
   reds at the check, pin still admitted; full gate re-run rc=0 / 0 FAIL / 28 `ok` / race control
   2/2. Safe headless because `ci.yml:118` installs go-verify from `releases/download/v0.30.0`, the
   **immutable tag** — found by the sonnet judge, a premise stronger than the controller's step-log
   reading. Evaluator **95/100, zero blocking**.
   **`9/OD-10` — THE HUMAN DECISION, one word:** `ci.yml:25` installs `releases/latest` for legs 1-2
   (today v0.33.0) while `:64`/`:118` pin v0.30.0 for the package leg and go-verify. Either **PIN**
   job 1 to the v0.30.0 tag (CI stops tracking upstream releases; World notices a compiler change
   only when it chooses to) or **ACCEPT** the drift and keep only the warning this iteration landed.
   The hard `verify_ail.sh` assertion is COUPLED to whichever is chosen and must not be added alone.
   **`9/CF-A-1`** (judge non-blocking, accepted): the two new version assertions ship with no
   IN-SCRIPT known-positive control, unlike `verify_go.sh:89` (race) and `:57` (tracked-binary
   count); the proving shim was ad hoc and is not committed. Fold a committed shim fixture into the
   human-gated half.~~ Prior head text follows.]** ~~[BACKLOG — infra, not critical-path; **SHARPENED WITH FIRST-PARTY MEASUREMENT iter-27**]~~
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
10. [**[LANDED — ITEM COMPLETE 2026-08-07 (iter-61). ALL THREE MILESTONES SHIPPED: `BG.A` PR #47 →
    `278f102` (iter-58) · `BG.B` PR #48 → `39130ec` (iter-60) · `BG.C` PR #49 → `c6a14c0` (iter-61,
    evaluator `sonnet` PASS 94/100 r1, ZERO blocking). Doc →
    `design_docs/implemented/w-boundary-gate-tree-mutation.md`. dev CI green BOTH jobs on the `BG.C`
    merge, SHA-addressed and step-log verified (11/11 · 13/13, `failed=none`), in a window the
    status API reports as 0 incidents — so the green is ATTRIBUTABLE.**
    **`BG.C` — `AC1b` AND `AC6′` DISCHARGED, AND THE `R-EXT4` QUESTION IS ANSWERED FAVOURABLY,
    FIRST-PARTY, ON THE REAL CI FILESYSTEM.** The runtime backstop captures five observables of the
    live target before/after `check()` (placed BEFORE the guard-red fatals, so it runs on every
    path); sha256/size/mode/**inode** are asserted UNCONDITIONALLY and `mtime_ns` only after a
    20-trial granularity probe fires **20/20**, with no logging-and-continuing state. **THE
    CONTROLLER FINDING: the plan RECORDS `st_dev` for `t.TempDir()` and `repoRoot` but never
    COMPARES them** — a 20/20 probe on a fine-grained tmpfs would license an mtime assertion about a
    file on a possibly coarse-grained repo volume, i.e. *a detector that cannot detect, certified by
    a measurement of somewhere else*, invisible locally because on this host the two devs ARE equal
    (`16777230`). Now asserted before the 20/20 gate. Because `<20/20`, a dev mismatch and an unequal
    `mtime_ns` are EACH a `t.Fatalf`, job 2's `success` on `ubuntu-latest` **IS** the ext4
    measurement — `C1`'s assumption is now a measurement, and **`10/OD-9` stays unfired**; the honest
    limit is that it is BOUNDED, not numeric, since `verify_go.sh` runs `go test` without `-v`.
    **Four mutations, all RED as predicted, AST guard GREEN in every arm** (so the backstop is live
    INDEPENDENTLY of the guard): `M7` `cp` write+restore → *ModTime changed*; **`M10` NEW** `mv`
    restore → *inode changed*, the first proof inode earns its place; **`M8` NEW** forced
    cross-filesystem probe → *not transferable*; **`M9` NEW** probe writes removed → *fired 0/20*,
    the only proof the `R-EXT4` fail-loud branch is reachable. **`AC6′`**: ratio **1.3038 ≤ 1.50**,
    absolute 1.4700 s ≤ 3.0 s — while the doc's ORIGINAL `AC6` would have failed BOTH arms, arm A
    being unmodified base code. `10/OD-1` (closure scanning) stays OPEN and deferred; **`10/OD-2`
    CLOSED** (sprint evidence, discharged by `BG.A`'s `M5`). **Prior head text follows.**
    ~~**[IN-SPRINT] — `BG.B` LANDED 2026-08-07 (iter-60), PR #48 → squash `39130ec`; evaluator
    `sonnet` PASS 88/100 round 1, ZERO blocking; dev CI green BOTH jobs, SHA-addressed and step-log
    verified on the merge commit. The next unit of work is milestone `BG.C`, gated on nothing.**~~
    **`AC1a` IS DISCHARGED**; `M3`, `M6` and the deny-list truncation control all fired — each
    control-arm-first, proven landed by differing sha256, restored byte-identical, and run in a
    disposable sibling worktree, never the main checkout. `confinedWrite` is now the single write
    enforcement point (rejects at/beneath `repoRoot` **synchronously, before `rawWrite` is called at
    all**), all three of BG.A's writes route through it, `armBarrier`'s bespoke `insideRepo` block at
    `:367–:373` is DELETED in its favour, and an **AST** write-guard forbids the four write names
    outside the single permitted site while asserting it can SEE (file count `== 1`, permitted site
    LOCATED and its line reported, deny-list complete by length and membership).
    **THE FINDING WAS IN THE TEST, NOT THE WRITER — this sprint's own spine arriving inside its own
    repair.** The recording-writer test as first delivered **synthesised its own mutant and overlay
    paths and called `confinedWrite` directly** instead of driving the real harness `mutateViaOverlay`,
    so its `wantWrites = 2` counted only the writes the test had itself just made. Measured (rule 3d):
    with `mutateViaOverlay`'s own two writes reverted to bare `os.WriteFile` it **still PASSED 4/4
    arms** — blind to the write path it existed to cover. The regression *was* still caught, by the AST
    guard redding at `:460`/`:471`, so the protection was real; what was not real was the test the
    plan's non-vacuity requirement was written to obtain. Repaired to drive the real harness per arm
    with the sink **TEED** (delegating to the CAPTURED original — delegating to `os.WriteFile` would
    itself be an unpermitted call under the guard three functions below), with the arm mutations
    factored into shared `goArmMutate`/`ailArmMutate` helpers used by BOTH the gate and the recording
    test so the two cannot drift apart. **Non-vacuity by outcome divergence: the identical probe now
    REDS 4/4** with *"the harness recorded ZERO writes through the confined sink"*.
    The evaluator's **NB-1** (a plan-required code comment on the guard's two stated limitations was
    absent — reproduced first-party at **0** mentions against a firing control) was **FIXED HERE rather
    than carried**, because `BG.C` extends this same guard and the limitation it records — a selector
    deny-list is bypassable by import aliasing or reflection — is precisely `M7`'s reason to exist.
    **NB-2** (`EvalSymlinks` called twice on the rejection path) and **NB-3** (`wantWrites=2` becomes 3
    under an armed `WORLD_BOUNDARY_ARM_BARRIER`) carried as cosmetic; neither can produce a false green.
    **READ BEFORE `BG.C`:** `C1`'s nanosecond-`ModTime` premise is **APFS-ONLY** (200/200 on
    darwin/arm64; CI job 2 is `ubuntu-latest`, whose mtime comes from a tick-granularity coarse clock),
    so `BG.C` is a fail-loud **20/20** granularity probe whose failure is a TEST FAILURE naming both
    `st_dev`s, with `sha256+size+mode+inode` asserted unconditionally, and a pre-authorized fallback:
    record the refutation, keep the four filesystem-independent observables, open `10/OD-3`, and
    **NEVER lower the 20/20 threshold**. **Prior head text follows.**
    ~~**[IN-SPRINT] — `BG.A` LANDED 2026-08-06 (iter-58), PR #47 → squash `278f102`; evaluator
    `sonnet` PASS 89/100 round 1, zero blocking. The next unit of work is milestone `BG.B`, gated on
    nothing — and it carries ONE correction that must be applied before it is routed (below).**~~
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

11. [**[LANDED 2026-08-12 (iter-75) — ITEM COMPLETE. All three milestones dev-CI-green (both jobs, SHA-addressed, `checks=2` = expected 2): `TR.A1` `93df1ec` · `TR.A2` `1a12042` · `TR.B1` `6e207ca` · `TR.B2` `88eb850` · `TR.C` `625fb89`. `verify_ail.sh` totals **4/11/14 UNMOVED** across all five, exactly as the doc promised; no `.ail`, no store table, no REST route, no projection package. **ITEM 5's `P6.B` PREREQUISITE IS NOW SATISFIED** — the condition was TR.A+TR.B merged AND TR.C green, and all three are, as of `625fb89`.**

    **`TR.C` LANDED 2026-08-12 (iter-75), PR #63 → squash `625fb89`.** Evaluator `sonnet` **63/100 FAIL** — the lowest score this item recorded and by some distance the most valuable — with its ONE blocking finding reproduced first-party and **FIXED IN-PR**. `metered=$0.00`. Planner `opus` (lane fails closed, reason `missing-script`), executor `codex:gpt-5.6-sol`, judge `sonnet`. Plan: `design_docs/planned/w-transition-registry-trc-sprint-plan.md`; transcript `…/verification/w-transition-registry/trc-mutations.md`. **`host/broker/invoke_boundary_test.go`** (353 LOC, **zero production LOC**): outside `host/broker` no `Invoke` selector, neither exported session constructor, and no `broker.Session` exposure through plain, aliased or dot imports; inside, the three pre-registry calls pinned by identity AND exact count. **AC11 ACTIVATED 1 → 2**, tolerant arm deleted. The planner priced it at **1.25 days against the doc's 0.5** and refused to split it, and it established two things the doc had wrong by omission: enumeration must be a **filesystem walk**, because `go list` never reports `host/store/writer_lock_other.go` (`//go:build !unix`) on darwin *or* linux/CI (walk **39**, `go list` **38**), and the gate cannot live in `host/boundary`, whose landed `TestBoundaryASTWriteGuard` pins `const wantFileCount = 1`.

    **THE FINDING THIS ITEM ENDS ON, AND IT IS A THIRD DIRECTION OF THE CLASS THE PREVIOUS TWO MILESTONES NAMED.** The judge defeated the gate with ordinary Go. Every detector for `Invoke`, `NewSession` and `NewReplaySession` sat inside `case *ast.CallExpr`, so it matched only when the selector is the `Fun` of a call — and a **method value** (`call := s.Invoke`) or a **function value** (`mk := broker.NewSession`) reaches a raw broker session from outside `host/broker` with no reflection, no `//go:linkname`, no generics and no build tags. Reproduced first-party: build rc=0, vet rc=0, gate **rc=0 PASS**, `walked=40` proving the file was scanned and yielded zero findings. **All 32 mutation arms agreed** — 23 executor, 9 controller — because every one of them spelled the forbidden thing the same way. `TR.B1` was *mechanism tested, call sites unguarded*; `TR.B2` was *site tested, second branch unguarded*; `TR.C` is *branch tested, shape space never enumerated*. Rule 3j asks **how many ways can this refuse**; the dual, which nothing in the shared rulebook asks, is **how many ways can the thing it refuses be SPELLED** — a detector is a recogniser, and a recogniser's coverage is a property of its input grammar, not of its branch count. Fixed by matching the bare `*ast.SelectorExpr` and deleting the `CallExpr` arm (keeping both would double-count real calls and move the exemption from 3 to 6), measured safe first: the only other bare `.Invoke` selectors in production are **8 in COMMENTS**, invisible to an AST and red at base to any text scanner — the design's own reason for mandating ASTs. Three new hermetic controls pin the branch, proven by reverting the detector alone and watching exactly those three red.

    **Two instrument facts every later milestone inherits.** **(a) `go build ./...` DOES NOT COMPILE `_test.go` AT ALL**, so on a test-only milestone the mutation-BUILDS rule is satisfied vacuously and the real compile gate is `go vet`. **(b)** the rule-3j cut through ordinary `git diff` returned **0** on this milestone's own file (control: 29 via `git diff --no-index`) — iteration 74's untracked-file trap, reproduced on the next sprint that adds a file.

    **Prior head text follows.** ~~**[IN-SPRINT] — `TR.A` IS COMPLETE. `TR.A2` LANDED 2026-08-11 (iter-72), PR #60 → squash `1a12042`; `TR.A1` landed the same day (iter-71), PR #59 → squash `93df1ec`. dev CI GREEN BOTH JOBS on both merge commits, SHA-addressed, `checks=2` = expected 2, `unresolved_incidents=0`. Evaluators `sonnet` **94/100** (TR.A1, zero blocking) and **86/100** (TR.A2, both blocking findings reproduced first-party and FIXED IN-PR at `1df7f3c`, not carried). `metered=$0.00` for both. Planner `opus` (lane fails closed — `derive-planner-lane.sh` is absent here), executor `codex:gpt-5.6-sol`, judge `sonnet`; TR.A2 ran with NO planner, reusing iter-71's plan. Sprint plan: `design_docs/planned/w-transition-registry-tra-sprint-plan.md`; mutation transcripts `design_docs/verification/w-transition-registry/{tra,tra2}-mutations.md`. **`TR.A` WAS SPLIT AT ITS T4/T5 BOUNDARY** on the planner's measurement (2630 LOC ≈ 1315/day vs this mission's ~1000/day `VL.B` reference), and the split held: **`TR.A1`** = store CAS + the TR-CJSON-1 codec + descriptors + literal goldens, activating **AC1/AC4/AC10**; **`TR.A2`** = the `ObjectStore` seam, `Reader`/`ReadSnapshot` (head read exactly once), eager copy-isolated `Snapshot` with a head-keyed cache, pure `BuildNext`, CAS `Publish`, activating **AC2/AC3** to exactly 3 and 4. Hold set re-measured after landing: AC5/AC6/AC7 `count=0`, AC11 `count=1`, AC8 rc=0 with **0** `transitionreg` imports in production replay, `verify_ail.sh` totals **4/11/14 UNMOVED**. **TWO FINDINGS THIS ITEM PAID FOR AND EVERY LATER MILESTONE INHERITS:** (a) a refusal test asserting only THAT an error occurred pins no branch — `DecodeRevision`'s canonical re-encode is a second refuser behind every guard, so each case must pin its own MEASURED message (2 TR.A1 mutations survived before this; the TR.A2 evaluator confirmed it is not ceremony); (b) **a rule-3j audit anchored to a list of DECISIONS cannot contain the branches the sprint itself writes** — §4.3 enumerated Decisions 1/2/4 and missed Decision 3's parent/revision chain rules plus T5's own store-error wrapper, so THREE refusal branches shipped with zero coverage (all three proven survivals: mutant landed, built, whole package rc=0) and were closed in-PR. Also standing: **a green `go test` is not a green `go vet`** — `copylocks` sits outside `go test`'s default vet subset and `verify_go.sh` never calls `go vet`, so 5 findings were invisible to the local gate AND both CI jobs.

    **`TR.B1` LANDED 2026-08-11 (iter-73), PR #61 → squash `6e207ca`.** dev CI GREEN BOTH JOBS on the merge commit, SHA-addressed, `checks=2` = expected 2, `unresolved_incidents=0`. Evaluator `sonnet` **84/100**; its one blocking finding reproduced first-party and **FIXED IN-PR**. `metered=$0.00`. Planner `opus` (lane fails closed, reason `missing-script`), executor `codex:gpt-5.6-sol`, judge `sonnet`. Plan: `design_docs/planned/w-transition-registry-trb-sprint-plan.md`; transcript `…/verification/w-transition-registry/trb1-mutations.md`. **`TR.B` WAS SPLIT AT ITS T3/T4 (PACKAGE) BOUNDARY** on the planner's price — **1740 LOC ≈ 2 days** against the doc's **1 day**, and 1740/day is *higher* than the 1315/day that forced the `TR.A` split. **`TR.B1`** = `CapabilitySnapshot` + epoch-on-debit + one `debitGrant` mechanism, `Requirement`/`Allows`/`decideOver`, and the confined `BoundInvoker` seam — `host/broker` only, **AC5 ACTIVATED to exactly 2**. **`TR.B2`** = descriptor-bound confinement + the two-session fixture (AC6/AC7), `host/transitionreg` only, already scoped by the same plan's T4–T5 — **no re-planning owed**. The split is safe for the `TR.A` reason: TR.B1 leaves AC6/AC7 at their tolerant `count=0` arm, TR.B2 adds nothing to `host/broker`, so AC5 cannot regress.

    **`Invoke` IS NOW A ONE-LINE WRAPPER OVER UNEXPORTED `invoke`, AND THAT IS WHAT KEEPS `TR.C` REACHABLE.** The bound invoker's call into the pipeline would otherwise be a **FOURTH** production `Invoke` selector call — outside `host/broker` it breaks TR.C's zero-outside rule, inside it breaks the pinned exemption of exactly 3 — i.e. an earlier milestone moving a later one's criterion. Caught at pick time and verified first-party: the 3 exemption sites are exactly `mintAttendedApproval` (`publish_op.go:135`, `:162`) and `invokeAttendedPublish` (`:279`), matching TR.C's frozen set. New hold criterion **`AC-INVOKE3`** (n=3, p=3, control t=88) guards it, because nothing in AC1–AC10 notices a 4th call site and AC11 does not exist until TR.C.

    **THE FINDING EVERY LATER MILESTONE INHERITS, AND IT SHARPENS (b) ABOVE INTO A MECHANISM.** The executor deferred **17 of its 19** mutation arms and enumerated them; the controller ran 18 outside the sandbox and **one SURVIVED** (`J6-MUT-BIND-DECLARED-ALIAS`: `Bind` aliasing the caller's `Declared` slice instead of copying it — so an authority envelope is not frozen at bind time and a later unrelated write widens it). The evaluator, handed `debitGrant` as a named target, then found a **second**: the **failed**-replay debit site (`broker.go:390`) was pinned by nothing, while the live and succeeded-replay sites were. Both reproduced first-party (mutant LANDED, BUILDS, whole package rc=0 with the defect present) and both fixed in-PR as **subtests**, so AC5 stayed at 2. **In both cases the production code was CORRECT** — every direct observation agreed with the claim, which is exactly why a reviewer, an executor and an 18-arm sweep all missed them. The cause is structural: **a refactor that unifies N call sites into one mechanism makes you test the MECHANISM and silently stop testing the SITES.** Five instances across two consecutive milestones make *guard the helper, miss the call site* this repo's most reliably recurring defect class. Cheap instrument: per unified mechanism, ask **how many call sites it has and how many a test observes**.

    **`TR.B2` LANDED 2026-08-12 (iter-74), PR #62 → squash `88eb850`. `TR.B` IS COMPLETE.** dev CI GREEN BOTH JOBS on the merge commit, SHA-addressed, `checks=2` = expected 2, `unresolved_incidents=0`. Evaluator `sonnet` **96/100, ZERO BLOCKING** — the best score this item has had. `metered=$0.00`. Controller `opus`, executor `codex:gpt-5.6-sol`, judge `sonnet`; **no planner ran** — iteration 73's split already scoped `TR.B2` as T4–T7b, so no re-planning was owed. Transcript: `…/verification/w-transition-registry/trb2-mutations.md`. **`host/transitionreg/bind.go`** (161 LOC): `Binder`/`CapabilitySource` declared *in this package* (what keeps the identifier `broker.Session` out of its production code); `Bind` with three ordered refusals — absent transition · access denied carrying **the broker's own label verbatim** · wrapped target-bind error — plus a zero-snapshot guard; `Check` over all five authority-bearing pins including Decision 7's three execution selectors; the confined `Bound.Request`; and `Request`/`Allowed()`, capturing exactly one registry head and one capability reading at construction and never re-reading either. **AC6 AND AC7 ACTIVATED** to exactly 3 each, tolerant `-eq 0 ||` arms deleted, machine-checked with a known-positive control in the same call (AC5–AC8 tolerant arms **0**; AC11's `-eq 1` still **1**). AC5 unmoved at 2 — TR.B2 adds nothing to `host/broker`, exactly as the split promised.

    **THE EXECUTOR RAN ALL 21 ARMS, DEFERRED NONE, AND SELF-REPORTED A BROKEN INSTRUMENT RATHER THAN BANKING ITS OUTPUT.** The plan's rule-3j cut returned **0** because **`bind.go` is UNTRACKED during executor work and ordinary `git diff` omits untracked files** — rule 3a's trap inside the very instrument iteration 72 rewrote to close a *different* blindness (`%w`), and blind on any sprint that ADDS a file, i.e. most of them. Use `git diff --no-index /dev/null <file>`. Enumerating that way found **J14**, the `target.Bind` error-propagation branch J1–J13 does not name.

    **THE MIRROR OF THE FINDING ABOVE, AND THE REASON THIS CLASS KEEPS RECURRING.** The controller's own sweep enumerated every refusal branch in `bind.go` **from the FILE** (not from the plan's J-list, written before a line of TR.B2 existed), neutered each with `if false && <cond>`, and required LANDED + BUILDS before reading any result: **14 branches, 13 KILLED, 1 SURVIVED.** `equalRequirements` has **two** refusal branches — a length guard and an element-wise guard — while the doc, the plan and the executor each name **one** mutation for it. `if false && a[i] != b[i]` LANDED (`068fc0e4…` → `80589e53…`), BUILT rc=0, whole package **rc=0 with the defect present**: a proposal whose `ExpectedEffects` has the same **length** but different **content** passed `Check`, the exact case Decision 7 forbids — because the single covering test *appends* an element and so exercises only the length branch. **Iteration 73's instrument comes back CLEAN here** (the evaluator independently confirmed the call-site table): there the mechanism was tested and the SITES were not; here the SITE was tested and the mechanism's second BRANCH was not. **The instrument that covers both directions asks, per helper AND per mechanism: how many ways can this refuse, and how many does a test observe? — never how many tests name it.** Six instances in three consecutive milestones, and as in both TR.B1 survivals the production code is CORRECT, so reading it, running the suite and mutating the named mechanism all agree; what is absent is only the assertion that would notice a future edit. Fixed as a **subtest** (AC6 stays at 3), proven by the inverse arm: same mutant, `-skip` that test → **rc=0**.

    **Standing base condition, so no later iteration chases it:** `verify_go.sh` returns **rc=1** with `FATAL: active toolchain go1.26.4 miscompiles host/store/scan.go` unless `GOTOOLCHAIN=go1.25.6` is exported. That is the script refusing an unpinned toolchain, not a regression — rc=0 on the identical tree with it set.

    **`TR.C` IS `[NEXT]`** — the binding gate, and the last milestone of this item.**~~]
    **w-transition-registry** · clause-3 · the stable, snapshot-readable transition catalog that
    `w-mcp-projection` P6.B reads: stable-ID/content-hash descriptors, eager immutable snapshots
    over the EXISTING object + registry-head store, capability filtering that **calls** the landed
    `broker.Decide` law rather than restating it, and a structural gate that keeps the future
    dispatch path bound. Doc: `design_docs/planned/w-transition-registry.md` (704 lines, 8
    decisions, **28 verification rows**, **11 acceptance criteria** — all count-gated and
    directory-independent — **44 named mutations**). · **3.5 World days in 3 milestones**:
    **`TR.A`** immutable descriptor/snapshot + store CAS (2d) · **`TR.B`** capability snapshot +
    declared-effect confinement (1d) · **`TR.C`** the binding gate (0.5d). No `.ail`, no store
    table, no REST route, no projection package; `verify_ail.sh` totals stay **4/11/14**.

    **`TR.C` IS NOT OPTIONAL AND IS NOT A FORMALITY.** `TR.A`+`TR.B` deliver the *mechanism*;
    without `TR.C` the undeclared-effect guard is an unenforced helper, and item 5's prerequisite
    is NOT satisfied. That is `gpt5-6-sol`'s restored-reviewer objection — this repo's named
    recurring shape, **guard the helper, miss the call site** — and it is the reason the N−1
    `proceed` was not banked. Measured at `b0f323a`, which made the fix SMALLER: `.Invoke(` has
    **exactly 3** production call sites, ALL in `host/broker/publish_op.go:135,162,279` (control:
    **83** in `_test.go`); exported `broker.NewSession` has **ZERO** production callers; and there
    is **no coordinator-to-broker dispatch path in the repository at all**. So nothing needs
    retrofitting — the registry-mediated path is the one P6.B will BUILD, and it can be **born
    bound**. `TR.C` is an AST gate pinning the 3 legacy sites by identity AND by exact count.

    **Naming hazard, frozen by the design (the iter-53 prediction, now realised):** three adjacent
    registries must stay distinct — `host/registry.Registry` (interpreter EPOCH), `broker.Registry`
    (`map[string]Handler`, effect handlers), and the new `host/transitionreg` (NO package-scope
    type named `Registry`; semantic ID `world/transition-registry/v1`).

    **Standing gotcha this item surfaced, and it binds every future doc:** **`rg` is NOT a binary
    on this rig or in CI** — it is a shell function the agent harness injects, absent under
    `env -i`, and used in **0** of `ci.yml` and all six `scripts/*.sh`. Eight acceptance criteria
    shipped using it and every one returned rc=0 at base *without executing the `rg` branch*.
    Never put `rg` in a committed command. · queued iter-70.


12. [**[LANDED 2026-08-12 (iter-77), PR #64 → squash `40164ea`, dev CI GREEN BOTH JOBS on the merge
    commit (SHA-addressed, `present=2` = expected 2, `unresolved_incidents=0`). Evaluator `sonnet`
    93/100 PASS, zero blocking. `scripts/verify_ail.sh` now pins the Leg-1 module SET by identity
    (`LEG1_MODULES`, 11 entries), enumerating once and comparing BEFORE any `ai-check` runs, NUL-
    delimited end to end; `host/verifygate/module_manifest_gate_test.go` carries FIVE committed
    refusal arms, each gated on a pristine control from its own isolated root. Totals 4/11/14
    unmoved. Move the doc to `implemented/` WITH its `-sprint-plan.md` companion at the next
    bookkeeping pick.**
    **THREE THINGS THIS ROW OWES A LATER READER, because each cost real evidence to establish:**
    **(1) THE DOC'S OWN §4.2 ORDERING WAS DEFEATED BY THE SHELL.** `/usr/bin/env bash` here is
    3.2.57 and the script sets `set -uo pipefail`, under which `"${arr[@]}"` on an EMPTY array is an
    unbound-variable ABORT (measured rc=127; `${#arr[@]}` is safe, rc=0). Write-then-guard would have
    killed the script before its own null-case message printed. Guards precede the `printf` and read
    `${#arr[@]}`. **(2) THE JUDGE DEFEATED THE GATE THIS ITEM SHIPPED, AND IT WAS FIXED IN-PR.**
    `find -name '*.ail'` is case-SENSITIVE, so `world/SNEAKY.AIL` never entered the swept set and the
    gate printed its own `✓ swept .ail module set equals the LEG1_MODULES allowlist (11 modules)`
    line with an unenumerated module in `world/` (control: `-name` 4, `-iname` 5). Repaired with
    `-iname` — case variants are 0 today, so the set does not move — and pinned by a fifth arm,
    proven by reverting only the detector (exactly that arm reds; inverse rc=0). **A recogniser's
    coverage is a property of its ENUMERATOR, one level below its branches.** **(3) BOTH NULL-CASE
    ARMS ARE COMMITTED, NOT ONE-SHOTS** — the doc offered a demotion path and the planner refused it,
    because a one-shot is a proof about a tree that no longer exists, which is this row's own
    original complaint about iteration 71. Within the hour that call was vindicated: the branch those
    arms cover is the one the doc got wrong. **Prior row text follows.** ~~[NEXT] — DESIGNED
    2026-08-12 (iter-76), READY FOR sprint-planner. The design doc exists and
    is quorum-cleared: `design_docs/planned/w-ail-gate-module-pin.md` (557 lines, commit `d201a1e`),
    2 rounds BOTH BLOCKED with 2 external reviewers present in each (`absent_reviewers` empty, no
    N−1 degrade), narrow-refinement carve-out applied to round 2. Iteration 77's pick is the SPRINT,
    not another doc. THREE CORRECTIONS THIS ROW OWES ITS READER, all measured at `304120b`:**
    **(a) THE FIX THIS ROW PRESCRIBES IS WRONG AND THE DOC DOES NOT IMPLEMENT IT.** The row asks for
    "an exact-total assertion mirroring `:239`" — i.e. `EXACT_TOTAL_MODULES=11`. A count pin is
    **defeated by an add-one-delete-one mutant**: with `world/_stray.ail` added AND
    `sketches/storejournal.ail` deleted the gate prints `✓ 4/4 … across 11 module(s)` — a success
    line **byte-identical to the pristine baseline's** — and exits **rc=0 PASSED**. Reproduced
    first-party, tree restored byte-identical. The doc ports the **identity allowlist** the sibling
    leg already implements (`verify_world_package.sh:86-96`), which coding-standards **S1** binds
    ("never aggregate counts alone"); the count is then redundant and is deliberately NOT added, so
    item 13 gains no third literal to maintain.
    **(b) THE `MUT-AIL-EMPTY-MODULE` CLAIM BELOW IS CONFIRMED AND NOW HAS A SECOND HALF**: deletion
    is as invisible as addition — `rm design_docs/sketches/storejournal.ail` → `… across 10
    module(s)`, **rc=0 PASSED**. Pinning `world/` only would leave half the measured defect open.
    **(c) THE OBVIOUS IMPLEMENTATION IS UNSAFE FOR A REASON THAT COST A QUORUM ROUND**: committed
    RED arms must NOT mutate the live tree. `verify_go.sh:108` runs `go test ./... -count=1` with
    **no `-p 1`**, and `host/boundary`'s `enumerateAIL` (`allowlist_world_test.go:197`) walks the
    **live** `world/` and `:293` reads every file it finds — so a stray arm's create/cleanup window
    can ENOENT an unrelated package. The arms run in an isolated `t.TempDir()` root against the real
    *copied* script. (`host/broker`'s AST gate does NOT collide — it filters `.go` at `:149`.)
    **Prior row text follows, still accurate on the defect itself.** Filed 2026-08-11 (iter-71),
    MEASURED not suspected. ~~Still needs a small design doc.~~**]
    **w-ail-gate-module-pin** · clause-1 · `scripts/verify_ail.sh` **never compares the module
    count against 11**. Measured first-party at `adfaa0b`: the variable `checked` is used in
    exactly four places — initialised `:167`, incremented `:176`, compared **ONLY against zero**
    `:233`, and **printed** `:243`. The known-positive control in the same read is what makes this
    a measurement rather than a failed grep: `total_verified` **IS** exactly pinned
    (`-ne "$EXACT_TOTAL_VERIFIED"`, `:239`), so the script demonstrably *can* pin a total and
    simply does not pin this one. Consequence: the charter's thrice-repeated "totals stay 4/11/14"
    is enforced for the **4** and (per the script's own comment at `:249`) only secondarily for the
    **14** — the **11 is decorative**. `w-transition-registry`'s `MUT-AIL-EMPTY-MODULE` was proven
    to LAND (4→5 `world/*.ail` files, 11→**12** modules) with the gate still **rc=0 and PASSED**.
    Iteration 71 repaired the *acceptance command* (greps the printed total; non-vacuous both arms:
    base rc=0 `modules11=1`, mutant rc=1 `modules11=0`) but deliberately did **not** touch the
    script — item 11's design forbids it, and the evaluator independently agreed that was the right
    scope call. **The residual is that an AC command is a one-shot proof on a tree that no longer
    exists, not a gate**: a guard is not a gate until something reds when you remove it. This item
    moves the pin into `verify_ail.sh` itself, with an exact-total assertion mirroring `:239`, a
    deliberate allowlist for intentional module additions, and a RED mutation proving a stray
    `world/*.ail` cannot pass. · ~0.5d · NEEDS A DESIGN DOC (small — may be folded into a
    `TR.*` milestone's docs task if the sprint has room).
13. [**[LANDED 2026-08-13 (iter-81) — ITEM 13 IS COMPLETE.** PR
    [#66](https://github.com/sunholo-data/ailang-world/pull/66) → squash `36f0c7a`, Gate 3b GREEN
    on the merge commit itself (`checks=2` = expected 2, `present == expected` asserted, 0
    non-success), evaluator `sonnet` **89/100, zero blocking**. `metered=$0.00`. Docs moved to
    `design_docs/implemented/` with their sprint-plan companion. The repo's **5th Z3-proven
    identity** landed: `EXACT_TOTAL_VERIFIED` 4→5, `EXACT_TOTAL_TESTS` 14→20, identities
    `gradeCode_test_1`–`_6`. **THREE OF THIS DOC'S OWN CLAIMS DID NOT SURVIVE MEASUREMENT**, all
    found by the planner, re-confirmed by the controller before routing and again by the judge
    after: **(a) AC9 WAS UNSATISFIABLE BY CONSTRUCTION** — it required the package `interfaceHash`
    to move, and `host/pkgproj/pkgproj.go:86` hashes only name/edition/ailang/export-MODULE-NAMES/
    effects, never opening a source file, so adding a function to an already-exported module leaves
    it invariant *by derivation* (`d16cc882 → d16cc882`). The only way to satisfy it was to
    hand-edit the golden, **which AC9's own next clause forbids** — the doc demanded an act and
    prohibited it in adjacent sentences. Amended to: three fields move, `interfaceHash`
    byte-identical. **(b)** the Conflict Surface omitted a **fifth** file,
    `host/verifygate/module_manifest_gate_test.go:128`, whose hardcoded `"✓ 4/4 … across 11
    module(s)"` marker reds CI's go-verify job — AC12's "four metadata files" is false. **(c)**
    `scripts/verify_ail.sh:376` was an **unguarded literal** (`4 … 14` hardcoded while `:315`/`:370`
    interpolate). A **sixth** file, `docs/SELF_MOD_PUBLISH.md`, moved into the CODE commit because
    the five-file boundary reds `host/runbook`'s AC28 — the digest repair is a gate dependency, not
    prose. **AND THIS ROW'S OWN `4/11/14` SHORTHAND WAS WRONG**, which is why it is corrected below
    rather than left for the next reader: there are **four** distinct mechanisms, not three equality
    pins — `EXACT_TOTAL_VERIFIED` (shell, `:310`), `EXACT_TOTAL_TESTS` (**PYTHON**, `:340`, spaces
    around `=`, invisible to a shell-shaped grep), `REQUIRED_TESTS` (python set, `:333-339`) and
    `LEG1_MODULES` (bash **array** compared as a **SET**, `:135-147`). **`EXACT_TOTAL_MODULES` does
    not exist**; the `11` is an allowlist cardinality. Non-vacuity measured, not asserted: moving
    `EXACT_TOTAL_TESTS` to 20 *without* the six names in `REQUIRED_TESTS` still passes (M12), and a
    **consistent** `=> PROVEN` arm in both contract and body leaves **Z3 fully green** (M7) — so the
    no-`PROVEN` property rests on tests plus a grep, never on the proof. That residual is carried on
    item 17. **Prior head text follows.**]
    ~~[**DESIGNED 2026-08-13 (iter-79), QUORUM-CLEARED — `[NEXT]`: this is the sprint to plan and
    execute.** Doc `design_docs/planned/w-evidence-grade-mapping.md` (620 lines), committed
    `6d12a79`, priced 0.65d. Two quorum rounds, both external reviewers present in both
    (`absent_reviewers` empty); R1 both-reject, R2 `gemini-3-1-pro` PASS + `gpt5-6-sol` REJECT on
    one non-directional point resolved under the **narrow-refinement carve-out** with the
    reviewer's verbatim fix. `metered=$0.177493`. **THE DECISION, since it changes what this row
    delivers: representation-only.** `Evidence` keeps its five constructors; a Z3-PROVEN total
    `gradeOf(Evidence) -> EvidenceGrade` lands in `world/types.ail` as the repo's **5th proven
    identity** (`EXACT_TOTAL_VERIFIED` 4→5, `EXACT_TOTAL_TESTS` 14→20, identities
    `gradeCode_test_1`–`_6` — all three MEASURED by running §3.2's literal code on an isolated copy
    of the real `world/` tree, baseline control 14/14 first). `CompilerOutput`/`HumanApproval` →
    `ATTESTED`. **`PROVEN` STAYS UNREACHABLE ON PURPOSE** — the round-1 proposal to add
    `ProofReport`/`ReplayReport` carriers was WITHDRAWN because an agent can mint one from an
    unchecked `HashRef`, which converts a representation gap into a grade-laundering authority gap;
    mint authority is a declared open obligation and the follow-on
    **`w-validated-proven-evidence-boundary`** owns it. Measured while settling that: there is **no
    production `Evidence` producer in the repo at all** (0 non-test Go hits, control 13), no Z3
    report producer, and no `.ail` module reads `Evidence` — so there were no producers to wire.
    **TWO GOTCHAS BEYOND THE ONE THIS ROW ORIGINALLY NAMED:** the edit also moves the
    `packages/world-core/world/types.ail` projection (Leg 3 step 3/9 sha equality) and the
    `world_package_ready_packet.golden.json` (step 9/9, byte-for-byte); and the totality guard is
    the CONTRACT, not the typechecker — v0.30.0 accepts a non-exhaustive ADT match, and a missing
    arm reds only via `verify.errors=1` while `rc` stays **0**, so no AC may read an exit code.
    **Prior head text follows.**]
    ~~[**FILED 2026-08-11 (attended, Mark) — MEASURED ABSENT, AND IT IS THE CHEAPEST HIGH-LEVERAGE
    ITEM IN THE UI PROGRAMME. `6b/§7` RATIFICATION POINT 2 WAS RATIFIED AS A DOCUMENT AND ITS
    DELIVERABLE WAS NEVER PRODUCED.**]~~
    **w-evidence-grade-mapping** · clause-5 · the TOTAL `Evidence` → grade mapping that P3
    (trust-gradient rendering) requires, now supplied by `gradeOf(Evidence)`. **Originally measured
    at `871e3b6`:**
    `world/types.ail` defines exactly **five** `Evidence` variants — `CompilerOutput(HashRef)`,
    `TestReport(HashRef, bool)`, `HumanApproval(HashRef)`, `AiReview(HashRef, float)`,
    `RecordedEffect(HashRef)` — and a repo-wide search of `world/` + `host/` for
    `PROVEN|TESTED|ATTESTED|CLAIMED` returned **zero** non-comment hits. The mapping gap is closed:
    `CompilerOutput` and `HumanApproval` map to `ATTESTED`. The stated producers of the TOP grade
    (Z3 proof, deterministic replay)
    have **no `Evidence` carrier**, so `PROVEN` is unreachable. A faithful renderer can grade
    decoded `Evidence` but can never show `PROVEN` — in a system whose
    distinguishing feature is machine proof. **A gradient whose top grade cannot be produced is not
    a gradient; it is a two-tone badge that teaches the reader to ignore the channel**, which is
    the anti-pattern list's own "grade laundering" arriving by omission rather than by intent.
    §7 point 2 already names the recommended shape (add/reshape variants so a carrier distinguishes
    a compiler result from a verified proof and preserves human ratification without mislabelling it
    as an unverified agent claim) — **that recommendation is explicitly NOT the decision**, so this
    item still owes a design doc. **THE GOTCHA THAT WILL BITE THE EXECUTOR:** `world/types.ail` is
    the pure core, so this edit moves the required-check manifest — `scripts/verify_ail.sh` pins
    `EXACT_TOTAL_VERIFIED` and `EXACT_TOTAL_TESTS` as **exact equalities** (`4/11/14` today), and
    new contracts or tests red the repo's primary gate for a reason unrelated to the change unless
    the pins move in the SAME commit. See item 12. · ~0.5–1d · NEEDS A DESIGN DOC · gated on
    nothing, and gated on no other queue item.]~~
14. **[LANDED] — **`WB.K` LANDED 2026-08-25 (iter-123): ITEM 14 IS COMPLETE AT 11 OF 11 MILESTONES, ALL 32 MUTATION ROWS DISCHARGED BY MEASUREMENT.** PR [#94](https://github.com/sunholo-data/ailang-world/pull/94) -> squash [`3dda87e`](https://github.com/sunholo-data/ailang-world/commit/3dda87e); Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`, 0 not-green, every count asserted NUMERIC before comparison, run existence `total=1 event=push`, parent control `rev-parse`d -> `check-runs=2`). PR checks OBSERVED green on the head SHA before merge; `mergeable` read FIRST (`MERGEABLE`/`CLEAN`). Evidence-only: **+234 lines** (sprint plan §7i) plus the `planned/` -> `implemented/` move of the doc and its plan; no production and no test code changed. **M24-M28 ALL KILLED, NO SURVIVOR** — M24/M25/M27/M28 SOLE KILLERS at subtest granularity, M26's set 4 members all inside `TestNewRefusesNonLoopbackBind`. **§7i(a): THREE OF FOUR LANDING LEGS READ GREEN ON A MUTATION THAT LANDED INSIDE A COMMENT.** Four-leg predicate (shape-appropriate counts; an exact line-content assertion whose expected value lives in the harness; `gofmt` on **rc AND size**; a query against the **PARSED form**), negative-controlled before use with four arms. **NC4 decisive**: mutant text inserted as `// _ "net/http"` -> leg 1 `0->1`, leg 2 satisfied BY CONSTRUCTION, leg 3 rc=0 / **0 bytes**, build rc=0, and **only the parsed-form query refuses**. NC3 sharpened it and corrected the section's own first draft: a **substitution**'s two-sided predicate already refuses, so the exposure is the **one-sided** predicate an **insertion** forces — and **M24 is the only insertion-shaped row of the 32**. Corroborates first-party the rule V1 iteration 274 added to the shared skill the SAME DAY. **THE RUNNING SKILL WAS 27 LINES BEHIND `origin/dev` AND THE MISSING LINES WERE THIS MILESTONE'S GOVERNING RULE** (running 3,757 vs origin 3,784, one hunk `origin:1108-1134`); read from origin and applied as leg 4 before it was readable in the running copy. **THE §7h RETROSPECTIVE JUDGE PASS OWED BY ITER-122 IS DISCHARGED**: evaluator `sonnet`, own worktree, distinct from the `opus` controller — **93/100 PASS, ZERO BLOCKING**, twelve named targets across §7h and §7i, **none refuted**, both instruments rebuilt rather than inherited. **3 OF 3: IT FOUND A SURVIVOR** — two, on `supportedWorkbenchQuery`'s pair-composition guards (`workbench.go:72`/`:75`, `&&` -> `||`), full classification arm **rc=0 with an EMPTY red set**, reproduced first-party and proved LIVE on the function's own truth table; routed to **row 34** as its sixth hunk, not absorbed. **ACCEPTANCE** (base -> post): AC1 1->0, AC2 1->0, AC3 0->0, AC4 0->0, AC5 1->0, AC6 1->0, AC7 broken->0, AC8 1->0; `gofmt -l host/ cmd/` empty. **AC7 is sound for the first time**: D1's blindness is still live (planted probe seen **0** by the doc's form, **1** by the repaired form) but unreachable now that nothing is untracked — 7 paths, **4 `A` + 3 `M`, 0 `D`**, control (no pathspec) **15**. `metered=$0.00`. Prior head follows.** ~~[**[[IN-SPRINT] — `WB.J` LANDED 2026-08-25 (iter-122), **10 OF 11 MILESTONES**. PR [#93](https://github.com/sunholo-data/ailang-world/pull/93) -> squash [`b0d973c`](https://github.com/sunholo-data/ailang-world/commit/b0d973c); Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`, 0 not-green, every count asserted NUMERIC before comparison, run existence `total=1 event=push`, parent control `rev-parse`d -> `check-runs=2`). PR checks OBSERVED green on the head SHA before merge, never behind an armed auto-merge; `mergeable` read FIRST (`MERGEABLE`/`CLEAN`). Evidence-only: **+124 lines in one file** (sprint plan §7h), no production and no test code changed. **ALL TEN ROWS DISCHARGED — M10–M13, M22, M23, M29–M32 — AND THERE IS NO SURVIVING MUTANT IN THIS MILESTONE**, the first drill milestone of the four to end clean at every site. Six are SOLE KILLERS at subtest granularity (M10 `/default-off`, M11 `/oversize`, M13 `/from-overflow`, M31 `/unknown-parameter`, M32 `/duplicate-parameter`, M12 `TestWorkbenchTimelineBound`); **M22 and M23 are broad-blast with `diff`-IDENTICAL 22-member red sets** (a shared `wants` header map at `workbench_test.go:41–50` is asserted on every workbench response, so any header change reds every workbench test at once), and **M29/M30 likewise share an identical 2-subtest set**, so neither is the SOLE KILLER its catalogue row implies. **`M12` IS KILLED BY A HARDCODED LITERAL, AND THE ASSERTION ITS ROW POINTS AT CANNOT DETECT IT.** `workbench_test.go:323` compares `strings.Count(body, "<h3>entry ")` against `workbench.WorkbenchPageLimit` — **the very constant M12 mutates** — so expected and actual move together; the seeding loop at `:301` moves with it too. Proved by a decisive arm rather than argued: **MC2** = M12 plus neutering ONLY the sibling literal `:329` `if strings.Contains(body, "<h3>entry 100</h3>")` → vet rc=0, build rc=0, **test PASSES rc=0**. The pin is live but carries a decorative member; the risk is latent, since generalising `:329` to track the constant — the natural cleanup — would hollow it silently. Distinct from §7d(c)'s CSP case, which was a true SURVIVOR. **THE §7d(c) TELL WAS PUBLISHED AT ITER-113 AND NEVER RUN AS A SWEEP** — *guard the helper, miss the call site*, aimed at a diagnosis rather than at code. Sweep run now with scopes asserted by `test -d` (6/1/1 test files) and a fresh negative control at **0**: **exactly one** assertion-side hit repo-wide, M12's own; `:280`/`:301` name production constants in SETUP only, and the `workbenchCSP` occurrence is inside the COMMENT documenting iter-113's repair — recorded because my own known-positive control returned **1** where I predicted 0, i.e. the control matched prose, which is rule 3a's trap aimed at the control rather than the check. **A HARNESS INSTRUMENT FAILURE FIRED AND REFUSED A VERDICT.** M22's first arm returned `LANDED=NO — INSTRUMENT FAILURE, not a verdict`: the landing predicate required the NEW literal's occurrence count to RISE, and for a DELETION mutant the new literal is the empty string, which `grep -c -F -- ""` matches on **every line** (268 before, 268 after) — unsatisfiable by construction however perfectly the mutation applied. It failed toward refusal, not toward a false SURVIVED; a `>=` instead of `>` would have certified the instrument against itself. Re-run under a shape-appropriate predicate (old count FALLS 1→0 plus the exact line-content assertion) M22 is KILLED. **Rule carried to `WB.K`: the landing predicate must match the mutation's SHAPE.** Every arm LANDED by occurrence count PLUS an exact line-content assertion **whose expected value is written in the harness and never read back from the file under mutation** (iteration 121's `LINES` lesson at instance 2 — 121's was in a controller instrument, this one is in committed test code that had already passed a quorum, a plan, an executor and an evaluator PASS); BUILT rc=0 before any test result was read; classified by the COMPLETE enumerated red set at subtest granularity; restored by `cp` (never `git checkout --`) and verified byte-identical, with the pristine control rc=0 after EVERY arm. Harness itself negative-controlled before use: a deliberately wrong expected line produced `INSTRUMENT FAILURE` and a clean restore. Final tree byte-identical, `git status --porcelain` empty. M31/M32's stale `parseWorkbenchQuery` catalogue text CONFIRMED first-party (`grep -rn` → **0** tree-wide; guards live inline at `:118`/`:122`, identifiers `key`/`values` not `k`/`vs`) — travels with **row 34** as §7g routed it, not absorbed. **⚠ THE EVALUATOR LANE DID NOT RUN THIS ITERATION** — agent spawning was disabled for this session, so generator≠judge was NOT satisfied; this is a CAPACITY gap, not a judgment one, and it is FLAGGED rather than papered over. Both prior drill milestones' judges found a real surviving mutant the controller had missed (2/2), so the gap is load-bearing: a retrospective judge pass over §7h is owed before `WB.K`. `metered=$0.00`. Prior head follows.** ~~[**[[IN-SPRINT] — `WB.I` LANDED 2026-08-25 (iter-121), **9 OF 11 MILESTONES**. PR [#92](https://github.com/sunholo-data/ailang-world/pull/92) -> squash [`2e7154b`](https://github.com/sunholo-data/ailang-world/commit/2e7154b); Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`, 0 not-green, every count asserted NUMERIC before comparison, run existence `total=1 event=push`, parent control `rev-parse`d -> `check-runs=2`). PR checks OBSERVED green on the head SHA before merge, never behind an armed auto-merge; `mergeable` read FIRST (`MERGEABLE`/`CLEAN`). Evidence-only: **+163 lines in one file** (sprint plan §7g), no production and no test code changed. **M1–M8 DISCHARGED, EACH A SOLE KILLER AT SUBTEST GRANULARITY — AND `M9` IS RECORDED SURVIVED AT EVERY ONE OF ITS FIVE SITES.** `TestWorkbenchRefusalBranches/store-error` asserts a `500` that TWO guards on the same request path can each produce, so neutering either alone is invisible: `:179` alone PASSES (500 from `:256`), `:256` alone PASSES (500 from `:179`), the PAIR **FAILS** `status = 200, want 500`. An all-five-guards arm returns a red set IDENTICAL to the pair's, so `:190`/`:214`/`:241` are unreachable with a failing store from anywhere in the repo. This is the sprint's **fifth** hollow-pin mechanism and it is distinct from `WB.D`'s: that one is *the row does not identify a site*; this one survives AFTER the site is identified. Per §3's `survived_is_a_result` the mutant was not repaired, the test was not repaired, the row was not omitted — residue routed to **queue row 37**, whose count and taxonomy this iteration also corrected (the nine `if err != nil {` occurrences are 3 parse guards, 1 log guard and 5 store-error guards, not nine store-error guards; its stale `179` reading is dated to `WB.F` landing after `WB.D`). Every arm LANDED by occurrence count plus an exact line-content assertion, BUILT rc=0 before any test result was read, classified by the COMPLETE enumerated red set, restored by `cp` (never `git checkout --`) and verified byte-identical, with the pristine control rc=0 after EVERY arm; final tree byte-identical, `gofmt` empty. **A HARNESS INSTRUMENT FAILURE IS RECORDED RATHER THAN QUIETLY RE-RUN** — the multi-site arm assigned its site list to a shell variable named `LINES`, a special integer parameter in `zsh`, so `179,256` was evaluated as arithmetic and stored as `256`; both halves of the harness's own landing assertion derived from the corrupted variable and read `1 == 1`. Evaluator `sonnet` **97/100 PASS, ZERO BLOCKING**, own worktree at the sprint commit, distinct model from the controller who generated the work; handed all three headline claims as named targets to attack it re-derived every one and **refuted nothing**, then found a **SURVIVING MUTANT** on `supportedWorkbenchQuery`'s cardinality gate that no catalogue row covers — reproduced first-party and recorded to **row 34**, not absorbed. `metered=$0.00`. Prior head follows.** ~~[**[[IN-SPRINT] — `WB.H` LANDED 2026-08-24 (iter-119), **8 OF 11 MILESTONES**. PR [#91](https://github.com/sunholo-data/ailang-world/pull/91) → squash [`5fd1069`](https://github.com/sunholo-data/ailang-world/commit/5fd1069); Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`, 0 not-green, every count asserted NUMERIC before comparison, run existence `total=1 event=push`, parent control `rev-parse`d → `check-runs=2`). PR checks OBSERVED green on the head SHA before merge, never behind an armed auto-merge; `mergeable` read FIRST (`MERGEABLE`/`CLEAN`). Evidence-only: **+158 lines in one file** (sprint plan §7f), no production and no test code changed; **discharges M14–M21**. **ALL EIGHT ROWS DISCHARGED BY MEASUREMENT, NO SURVIVING MUTANT AMONG THEM** — each LANDED by OCCURRENCE count, BUILT rc=0 before any test result was read, classified by the COMPLETE enumerated red set, restored by `cp` (never `git checkout --`) and verified byte-identical, with the pristine control rc=0 after EVERY mutant. Four of eight are only attributable at SUBTEST granularity (M17/M18/M21 share a parent, M19/M20 share another); M18 and M21 are NOT independent, and `no-proven-inference` firing under M21 while absent from M18's red set is what proves it a distinct guard. **M14's ROW NAMED A SITE ITS OWN NAMED TEST CANNOT DETECT** — `template.HTML` on `Grade.Label` SURVIVES with an EMPTY red set, structurally, because `NewGradeView` refuses any label outside the four constants (M17's own guard); `Grade.Unavailable` is the SOLE KILLER. Row repaired before routing on §7c's precedent. **THE CLASSIFICATION ARM IS UNSATISFIABLE IN THE EXECUTOR LANE BY CONSTRUCTION** — it names `./host/daemon`, which binds real loopback sockets that `--sandbox workspace-write` denies; inside the sandbox all three daemon tests FAIL, outside all three PASS, so the red is the instrument. The Codex executor labelled them `UNINFORMATIVE UNDER SANDBOX` and **STOPPED before M15** rather than assert a control it could not satisfy — the correct call, and the second consecutive milestone whose executor disclosed rather than completed. **WB.I/WB.J/WB.K carry the same arm and are controller-work too** (charter row 4f's iteration-50 finding, recurring). Step 2's occurrence assertion **fired live and caught a FALSE SURVIVED**. Evaluator `sonnet` **88/100 PASS, ZERO BLOCKING**, distinct provider from both the Codex executor and the controller; it **refuted the controller's own citation** (two of three tests bind via `d.Listen()` → `net.Listen`, not `httptest.NewServer`) and found a **SURVIVING MUTANT** outside this milestone's scope — both reproduced first-party, the survivor recorded to **row 34**, not absorbed. `metered=$0.00`. Prior head follows.** ~~[**[[IN-SPRINT] — `WB.G` LANDED 2026-08-24 (iter-118), **7 OF 11 MILESTONES**. PR [#90](https://github.com/sunholo-data/ailang-world/pull/90) → squash [`bedc0d1`](https://github.com/sunholo-data/ailang-world/commit/bedc0d1); Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`, 0 not-green, every count asserted NUMERIC before comparison, run existence `total=1 event=push`, parent control `rev-parse`d → `check-runs=2`). PR checks OBSERVED green on the head SHA before merge, never behind an armed auto-merge; `mergeable` read FIRST (`MERGEABLE`/`CLEAN`). Test-only: **+129 lines in one file** (`host/boundary/allowlist_world_test.go`), no production change; closes **AC2**, claims M24/M25 for WB.K. **THE PLAN'S OWN ANTI-VACUITY CONTROL MADE ONE OF ITS ASSERTIONS UNSATISFIABLE:** task 3b requires `countDep(wbDeps,"net/url") == 0` AND, to stop eight zeros being vacuous, `countDep(wbDeps,"html/template") == 1` — and `html/template` TRANSITIVELY IMPORTS `net/url`. Measured on pristine dev, that closure holds `net/url` exactly **1** time (`net/http` 0, the six other forbidden entries 0, control 78 deps total). The executor DISCLOSED it and SPLIT the requirement rather than dropping a check: pin the transitive count at 1, plus an AST scan of every non-test `.go` under `host/workbench` forbidding a DIRECT import — strictly stronger than the plan's text. Two further deviations, both real, both self-reported, all three adjudicated BY MEASUREMENT and all three in the executor's favour: task 3d's callback polarity is INVERTED relative to `mutateViaOverlay` (whose contract `t.Fatalf`s on a `nil` return, so the plan's literal text fails the test exactly when the detector works), and go1.25.6 does NOT indent subtest `=== RUN` lines (raw 7, top-level 2, base 1) — the latter a defect in the CONTROLLER'S OWN directive, not the executor's work. Drill anchored to the DIFF (rule 3n): M24, M25, and a controller-added **MC** (a direct `net/url` import, aimed at the AST scan the plan never specified and therefore never pinned) are each LANDED-by-sha256 and BUILT before any test result was read, and each is a **SOLE KILLER**; all restored byte-identical by `cp`, never `git checkout --`. Evaluator `sonnet` **97/100 PASS, ZERO BLOCKING** in its own worktree, distinct provider from the Codex executor; handed all three deviations as named targets to attack, it re-derived M25, MC and M24 and added a fourth mutant of its own (`subpkg/leak.go` importing `os/exec`) — **no surviving mutant, no hunk without a killer**. Declared limitation, now a derivation rather than an excuse: `go list -deps` is a SET and **4** production files import `net/http` directly, so no single-file renderer overlay can move the daemon count — `daemon-transport-control-fires` is a non-tautology proof, not a closure measurement. `metered=$0.00`. Prior head follows.** ~~[[IN-SPRINT] — `WB.F` LANDED 2026-08-24 (iter-117), 6 OF 11 MILESTONES. PR [#88](https://github.com/sunholo-data/ailang-world/pull/88) → squash [`8f0037c`](https://github.com/sunholo-data/ailang-world/commit/8f0037c); Gate 3b GREEN on the MERGE commit (SHA-addressed from a `rev-parse`d 40-char SHA, `present=2 == expected=2`, both `success`, 0 not-green, every count asserted NUMERIC before comparison, run existence `total=1 event=push`, parent control `rev-parse`d → `check-runs=2`). PR checks OBSERVED green on the head SHA before merge, never behind an armed auto-merge; `mergeable` read FIRST (`MERGEABLE`). The PARKED-ON-LANE predicate from iter-116 (Sonnet weekly limit, resets 7am Europe/Copenhagen) fired: Sonnet re-probed available and evaluated `a96fd67` independently — **96/100 PASS, ZERO BLOCKING**, in its own worktree, distinct provider from the Codex executor. The controller reproduced the load-bearing pin FIRST-PARTY: neutering the `/workbench` route's `defer cancel()` reds the AC8 loop ONLY on the workbench row (3 recorded ctx nil, want `context.Canceled` — SOLE KILLER) while `TestWorkbenchReadDeadline` stays green (BYSTANDER), proving the appended route is exercised, not decorative. Test-only (+90 lines, `host/daemon/read_deadline_test.go`), closes **AC8**, claims M29/M30 for WB.J; assertion text is a literal constant, not the production identifier. Commit message AND PR body scanned for auto-close keywords (**0** each, known-bad control firing at 1); `#68` asserted OPEN on both sides of the merge. `metered=$0.00`. Prior head follows.** ~~[[IN-SPRINT] — `WB.E` LANDED 2026-08-23 (iter-115), 5 OF 11 MILESTONES. PR~~
    [#87](https://github.com/sunholo-data/ailang-world/pull/87) → squash
    [`e563339`](https://github.com/sunholo-data/ailang-world/commit/e563339); Gate 3b GREEN on the
    MERGE commit (SHA-addressed from a `rev-parse`d 40-char SHA, `present=2 == expected=2`, both
    `success`, 0 not-green, every count asserted NUMERIC before comparison, run existence asserted
    `total=1 event=push`, and the parent control **`rev-parse`d rather than hand-expanded** —
    `total=1`, `check-runs=2` — so a fabricated SHA could not have agreed with me for the wrong
    reason). PR checks OBSERVED green on the head SHA before merge, never behind an armed
    auto-merge; `mergeable` read FIRST (`MERGEABLE`). Evaluator `sonnet` **95/100 PASS, ZERO
    BLOCKING**, in its own worktree at the sprint commit. Payload is opt-in behind `payload=1`, the
    preview is capped by slicing BEFORE the string conversion — bounding the allocation and not
    merely the output — and the timeline read records which of the two conditions ended it.
    Production diff **12 insertions, 0 deletions**.
    **THE DEFECT WAS THE CONTROLLER'S AND THE EXECUTOR CAUGHT IT BY DISCLOSING RATHER THAN
    COMPLETING.** Adjudication ADJ-1 told it to request `/workbench?from=0` for the timeline test.
    That is a SIXTH state in §2.2's closed enumeration — the very thing iter-114 foreclosed as *"a
    design change, not glue"* — so to obey the directive it widened `supportedWorkbenchQuery` by one
    disjunct, and said so in its report. The **empty query** does the identical job (state 1 of the
    enumeration: it defaults `from=0` AND leaves `Page.Selected` nil by construction), so the
    widening was never needed. Reverted: `supportedWorkbenchQuery` is `diff`-identical to its
    parent, the `{"unsupported-combination", "/workbench?from=0", 400}` arm is restored verbatim,
    and the judge verified the revert independently. **A directive is an acceptance list too, and
    ADJ-1 was never checked against what the PREVIOUS milestone actually enforces** — the controller
    read §2.2's enumeration in the doc and not `supportedWorkbenchQuery` in the code.
    **THE DRILL WAS ANCHORED TO THE DIFF, NOT TO THE DEFECT (rule 3n), AND THAT IS WHAT FOUND THE
    FINDING.** Four hunks, four mutants, each LANDED by sha256 and BUILT rc=0 before any test result
    was read, each restored byte-identical from a `cp` backup and never `git checkout --`. `M10`
    (`showPayload := true`), `M11` (drop the `MaxPayloadPreview` slice) and `M12`
    (`WorkbenchPageLimit` 100→101) are each a **SOLE KILLER** with the `-skip` inverse at rc=0.
    **Hunk 4 — `page.Timeline.Truncated = …` — REDS NOTHING**: rc=0, empty FAIL set. A
    defect-derived mutation set would have run three mutants and reported a perfect drill.
    **THE UNPINNED HUNK IS HALF OF A TWO-SIDED SEAM DEFECT → QUEUE ROW 38.** `Truncated` is written
    by the handler and read by **nothing** (1 writer, **0** non-test readers; same-scope control
    `PayloadTruncated` = **3**). Its mirror image was surfaced by the judge and is worse:
    `TimelineView.NextHref`/`PrevHref` are **read by the template and written by nothing** — **2**
    template actions at `render.go:140-141`, **0** occurrences anywhere in `host/daemon`, control
    `page.Timeline` assignments in that same file = **3**. Both links are guarded by `{{if}}`, so
    the empty string renders *nothing*: the workbench caps the timeline at 100 entries and then
    emits no next-link on any page, including the one that IS truncated. The signal a next-link
    needs is exactly the field `WB.E` just shipped. Declared, not repaired (rule 3n(b)) — an
    observable means editing `host/workbench/render.go`, outside this milestone's file list, and
    §3.2 requires a visible "truncated" marker for the **payload** only.
    **THE JUDGE REFUTED THE CONTROLLER IN A USEFUL DIRECTION:** I claimed the timeline oracle is
    sound because `Page.Selected` renders nowhere. It is sound for a **stronger** and already-
    enforced reason — the closed grammar keeps "empty query" and "entry selected" mutually
    exclusive — so the oracle survives a future milestone rendering `Selected`. Recorded because the
    weaker justification would have coupled row 38's fix to this test. Executor
    `codex:gpt-5.6-sol`, **zero git writes**, two files only, one disclosed deviation.
    **NEXT: `WB.F`** — `TestWorkbenchReadDeadline` + `/workbench` in the cancelled-after-handler
    table, closing doc **AC8** (M29, M30). **Prior head text follows.**
    ~~[[IN-SPRINT] — `WB.D` LANDED 2026-08-23 (iter-114), 4 OF 11 MILESTONES. PR
    [#86](https://github.com/sunholo-data/ailang-world/pull/86) → squash
    [`e50fbea`](https://github.com/sunholo-data/ailang-world/commit/e50fbea); Gate 3b GREEN on the
    MERGE commit (SHA-addressed from a `rev-parse`d 40-char SHA, `present=2 == expected=2`, both
    `success`, 0 not-green, every count asserted NUMERIC before comparison, run existence asserted
    `total=1 event=push`, parent control `total=1` with `check-runs=2`). PR checks OBSERVED green
    on the head SHA before merge; `mergeable` read FIRST. Evaluator `sonnet` **82/100 PASS**, in
    its own worktree. The closed query grammar of §2.4 is now ENFORCED: `parseWorkbenchQuery`
    accepts exactly `{}`, `{world}`, `{from,entry}`, `{object}`, `{object,payload}` and refuses
    everything else 400 — nothing ignored, no precedence fallback. **Three plan-flagged points
    adjudicated before routing and all three upheld by the judge first-party:** `?from=5` ALONE is
    refused (§2.2 is a CLOSED enumeration of four states, so a fifth is a design change, not
    glue); `payload` must be exactly `0`/`1` (§2.2 writes the domain literally, so this is not
    plan-added); the two extra subtests are the non-vacuity pins for both. `TestWorkbenchRefusalBranches`
    ships **exactly 14** subtests, 0 SKIP, 0 FAIL — and `rc` is NOT the gate, because
    `go test -run` exits 0 on an empty match set (row 33), so at base that same command was rc=0
    with ZERO tests. **THE CONTROLLER ADDITION, AND ITS MEASURED PAYOFF:** every 400 branch writes
    the class token `BadRequest`, so arms asserting only status+class cannot tell each other
    apart. Each branch was required to carry its OWN constant message, each arm to assert THAT
    message, and every expectation to be a string LITERAL and never the production identifier
    (iter-113's tautological oracle, one milestone earlier). Forcing `supportedWorkbenchQuery`
    false reds **ELEVEN** arms including the positive control, and reds `malformed-world`,
    `absent-object` and the rest precisely BECAUSE each asserts its own message — under a
    class-token-only assertion those ten would all have passed. **THE JUDGE'S FINDING, REPRODUCED
    AND WORSE THAN FILED → QUEUE ROW 37:** it reported 4 of 5 store-error guards unpinned; swept
    across the whole pattern the figure is **7 of 9**, and because `failingStore` returns
    `ok=false` beside its error, a neutered guard reports a real store fault as **404**. `M9`'s
    row names one mutant for nine sites, so it cannot say which it discharges — rule 3a(i-e)
    aimed at the mutation TABLE. NOT absorbed (rule 3n(b)). **A DEFECT OF MINE THE EXECUTOR
    CAUGHT BY REFUSING TO PRETEND:** my directive named `.ailang/state/sprints/*.plan.json` as its
    gate, and `.gitignore:3` ignores `**/.ailang/`, so that file is absent from EVERY sprint
    worktree by construction; it disclosed this rather than inventing requirements. Executor
    `codex:gpt-5.6-sol`, **zero git writes**, two files only. **NEXT: `WB.E`** — payload opt-in,
    64 KiB cap, 100-entry timeline cap (M10–M12). Prior head text follows.]**
    ~~[[IN-SPRINT] — `WB.C` LANDED 2026-08-23 (iter-113), 3 OF 11 MILESTONES. PR
    [#85](https://github.com/sunholo-data/ailang-world/pull/85) → squash
    [`5fd6fb3`](https://github.com/sunholo-data/ailang-world/commit/5fd6fb3); Gate 3b GREEN on the
    MERGE commit (SHA-addressed from a `rev-parse`d 40-char SHA, `present=2 == expected=2`, both
    `success`, 0 not-green, run existence asserted `total=1 event=push`, parent control `total=1`
    with `check-runs=2` proving the control SHA resolves). Evaluator `sonnet` **82/100, ONE
    BLOCKING — and the judge was RIGHT**; fixed in the shipped commit after first-party
    reproduction. `GET /workbench` is the **ninth** mux registration and the route-table doc
    comment now says so; `handleWorkbench` takes `d.readCtx(r)` exactly once with `defer cancel()`
    before any store read, sets the five security headers on EVERY response including errors, and
    renders through `workbench.Render`. **AC5 CLOSED (base rc=1 → 0). AC6 IS *NOT* CLOSED BY THIS
    MILESTONE AND THE PLAN'S TABLE IS WRONG ABOUT IT** — measured at this milestone's base, AC6 is
    **rc=0**, because its `executor_baseline_rc: 1` dates from `3e0c34c` when `host/workbench` did
    not exist and `WB.A` created that directory. AC6's forbidden-token arm still fires (injecting
    `embed.FS` reds it; restore byte-identical at sha256 `007b0f0a…`), so it is carried as a
    **regression pin, never evidence** — the plan's own D3 arriving at its unswept sibling.
    **THE BLOCKING FINDING: A TAUTOLOGICAL ORACLE.** `TestWorkbenchSecurityHeaders` asserted the CSP
    against the production symbol `workbenchCSP` itself, so expected and actual moved together and
    `M22` — the mutant this very test is NAMED as the killer of, in the doc's §6 table and in plan
    task 12 — was invisible: it LANDED (occurrences 1 → 0), BUILT rc=0 before any test result was
    read, and left the package **rc=0 with an empty `--- FAIL` set** while the named arm PASSED.
    §3 rule 5 does not cover this: a tautology is not a mutant escaping a good test, and no later
    drill would have caught it either. Expectation is now a literal; `M22` and `M23` each red the
    arm **ALONE**. **THE EXECUTOR DISCLOSED A DEVIATION AND THE DISCLOSURE WAS THE VALUE:** plan
    task 13 asks for two distinct committed OBJECT hashes in the body, which is **unsatisfiable in
    `WB.C` scope** — see new queue row **35** — so it manufactured a `page.Notice` carrier and said
    so. Two arms, one variable: neutering only that line reds the test on exactly the two hash
    assertions, i.e. it was the SOLE carrier. Removed; the test now asserts the two entries'
    distinct `EntryHash` values (inside `<dt>entry hash</dt>…`, syntax only `{{.EntryHash}}`
    produces — the judge's correction, since `store.Commit` chaining makes entry 1's
    `PrevEntryHash` a second writer of entry 0's hash) plus the selected head world ref.
    Executor `codex:gpt-5.6-sol`, ~4 min, **zero git writes**. **NEXT: `WB.D`** — closed query
    grammar + every refusal branch (M2–M9, M13, M31, M32), extending `WB.C`'s single unknown-key
    guard site, which the judge confirmed `WB.D` can cleanly extend.
    Prior head text follows.]**
    ~~**[IN-SPRINT] — `WB.B` LANDED 2026-08-23 (iter-112), 2 OF 11 MILESTONES. PR
    [#84](https://github.com/sunholo-data/ailang-world/pull/84) → squash
    [`75bc23f`](https://github.com/sunholo-data/ailang-world/commit/75bc23f); Gate 3b GREEN on the
    MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`, 0 not-green, run
    existence asserted). Evaluator `sonnet` **91/100, ZERO BLOCKING**. `host/workbench` now
    RENDERS: one package-level parsed `html/template`, `Render(io.Writer, Page) error`, landmarks in
    document order, every `href` from a helper that can only emit `/workbench`/`/workbench?…`, full
    hashes shortened by CSS alone, absent edges rendered `UNAVAILABLE:` rather than omitted, and a
    dual-channel verdict. Five top-level tests, twelve subtests, zero skips; imports are `errors`,
    `html/template`, `io` only. Executor `codex:gpt-5.6-sol`, ~3.5 min, **zero git writes**, no
    deviations. **THE FINDING — A TASK ASSERTED ITS OWN TEST "MAKES M20 KILLABLE" AND IT DID NOT:**
    task 14's pin (body contains `TESTED` and `FAIL`) is satisfied by the branch-selected literal
    `aria-label="test verdict FAIL"`, which the mutation never touches, so the M20 mutant LANDED,
    BUILT, and left the whole package **rc=0 with an empty `--- FAIL` set**. The ROW was wrong, not
    the code (rule 3i(c)); the pin now requires `verdict: FAIL`, which only the
    `{{.Grade.Verdict}}` action can produce, and the arm is then the SOLE killer (`-skip` → rc=0).
    Recorded as plan §7c, including the disclosed tension with the plan's own §3 rule 5. **THE
    JUDGE'S OWN FIND BECAME QUEUE ROW 34:** three shipped template hunks (`{{if .World.Available}}`
    true branch, the `PayloadTruncated` marker, the PASS-verdict span) are pinned by nothing, and
    none appears in the doc's 32-row §6 table, so WB.H–WB.K's drill can never reach them.
    **`WB.C` IS `[NEXT]`** — it closes AC5/AC6, the first ACs any milestone closes. Prior head text
    follows.]**
    ~~[[IN-SPRINT] — `WB.A` LANDED 2026-08-22 (iter-111), 1 OF 11 MILESTONES. PR
    [#83](https://github.com/sunholo-data/ailang-world/pull/83) → squash
    [`83f1973`](https://github.com/sunholo-data/ailang-world/commit/83f1973); Gate 3b GREEN on the
    MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`, 0 not-green). Evaluator
    `sonnet` **97/100, ZERO BLOCKING**. `host/workbench` now exists: the transport-free view model,
    the two bounding constants, and grade/verdict constructors that refuse invalid, unavailable and
    verdict-less input — 211 lines, 9 subtests, 0 skips. Executor `codex:gpt-5.6-sol`, ~2.5 min,
    **zero git writes**; its `rc=1` on daemon+boundary was correctly labelled `UNINFORMATIVE UNDER
    SANDBOX` and is **rc=0** outside it. **NON-VACUITY SPOT-CHECK (not the WB.H drill):** an
    M18-shaped mutant reds only `unavailable-is-not-claimed`, an M21-shaped mutant reds only
    `no-proven-inference` — each test kills its own mutant and only its own. **A PLAN DEFECT FOUND
    AND RECORDED, NOT PATCHED:** precondition 3 states the pristine manifest is *"1079 lines … in
    the main checkout"*, but run where the precondition applies — the sprint worktree — the identical
    command returns **157**, because the two trees differ by **922 untracked files**; and the
    command's `-not -path './.git/*'` misses a linked worktree's `.git` **FILE**, which is the sole
    difference against `git ls-files` (**156**). Neither breaks AC7 (the delta is intra-tree and both
    errors cancel); what breaks is the cross-check. **`WB.B` IS `[NEXT]`.** Prior head text follows.]~~
    ~~[PLANNED 2026-08-22 (iter-110) — SPRINT PLAN LANDED, `WB.A`–`WB.K`; READY FOR THE
    SPRINT-EXECUTOR; NOT LANDED.]~~ Plan:
    [`design_docs/planned/w-workbench-read-only-sprint-plan.md`](planned/w-workbench-read-only-sprint-plan.md)
    (381 lines, tracked — the machine plan `.ailang/state/sprints/w-workbench-read-only.plan.json`
    is gitignored and so is ABSENT from a fresh sprint worktree, which is why the tracked companion
    exists). Planner **opus**, lane token `opus fail-closed:env-pin` from
    `derive-planner-lane.sh`, used VERBATIM. Eleven milestones, 9.5 h + 2.5 h contingency ≈ 1.5 d,
    inside §9's 1.5–2 d band; all 32 mutation IDs appear exactly once in `mutations[]` and the
    AC union is exactly {AC1…AC8}. **THE PLANNER FOUND THREE ACCEPTANCE CRITERIA THAT CANNOT FAIL,
    AND THE CONTROLLER REPRODUCED ALL THREE FIRST-PARTY.** `AC7` — *"only priced files changed"*,
    the criterion that stops scope creep — is **vacuous by construction**: its command
    `git diff --name-only 93e1ba5 -- ':!design_docs'` cannot see untracked files, and every file
    this sprint adds is untracked for the whole sprint *because this loop's own executor recipe
    forbids git writes in the sandbox*. Reproduced with a real probe file: doc form **0** matches
    while `ls` finds the file and `git status` counts it at **1**; the planner's repaired form
    (`+ git ls-files --others --exclude-standard`) returns **1**. Two quorum rounds plus a restored
    third reviewer could not have caught it — none of them knows how this loop runs its executor.
    `AC2` and `AC8` are green at base with only one of their two named tests in existence
    (**rc=0**, exactly **1** top-level `=== RUN` each), and the **negative control is the
    sharpening**: `go test ./host/boundary -run 'TestZzNoSuchTestIter110' -count=1 -v` is
    **rc=0 with 0 `=== RUN` lines**, so a `-run` selector exits 0 on an EMPTY match set — any AC in
    this repo shaped *"`go test -run 'TestA|TestB'` passes"* is green before either test is written
    and stays green if a rename orphans the selector. All three repaired in the plan. Base gates
    both **rc=0** at `3e0c34c` (including `host/capsule`, which row 32 reds only under gate-level
    load). The planner's three under-specified points **Q1/Q2/Q3 were ADJUDICATED, not parked** —
    each is a completeness gap the doc's own text (plus one measurement for Q1) closes uniquely, so
    all three plan readings are UPHELD and **zero human asks were manufactured**. Q1: `writeReadTimeout`
    is `writeAPIError(…, "Timeout", …, 503)` and `writeAPIError` is `writeJSON`, so §2.4's *"the same
    class the JSON routes emit"* names the **class token, not the envelope** — the workbench 503 is
    HTML carrying `Timeout`. Q2: §2.2 pairs `from` with `entry` while §6's M10 requires `{object}`
    alone to be accepted, so the set is exactly {∅, {world}, {from,entry}, {object}, {object,payload}}
    and `?from=0` alone is 400. Q3: §2.4's omission of `payload` from its malformed-value list is an
    enumeration gap, not a permission — the carve-out's own *"no parameter is ignored"* forbids the
    silent fallback. Prior head text follows.
    ~~[UNBLOCKED, REVISED AND CARVE-OUT-COMPLETE 2026-08-22 (iter-109) — READY FOR THE
    SPRINT-PLANNER; NOT LANDED.]~~ The blocking predicate was RUN, not transcribed: item 18's merge
    `d21754f` is an ancestor of `origin/dev` (`git branch -r --contains`), so the deferral Mark
    ratified is discharged. **Both round-2 blocking objections are answered by the WORLD, not by
    argument.** `gpt5-6-sol`'s own `proposed_fix` offered a second limb — *"otherwise defer
    `/workbench` until the separately proposed daemon read-cancellation item lands"* — and that is
    exactly what happened: all six store getters now take `ctx` (`store.go:475/530/559/636/810`,
    `read_object.go:43`), `handlers.go:270` derives `context.WithTimeout(r.Context(), d.readDeadline)`
    with `readDeadline = 10 * time.Second` (`daemon.go:128`), and `handlers.go:324` emits an explicit
    `503`/`Timeout`. `gemini-3-1-pro`'s premise objection was MEASURED rather than forwarded and came
    back CONFIRMING the doc. Designer **`claude:claude-fable-5`** (rotation; gemini skipped as
    read-only under `CapRemoteSandbox`) rebased the doc from `9491a10` → `93e1ba5` — it was **70
    commits / 65 non-`design_docs` files** stale — re-ran all 12 controller measurements first-party
    with **zero disagreements**, replaced `V12`/`V14`/`V16`–`V19` as no-longer-true and added
    `V20`–`V25`. **ROUND-3 RE-QUORUM BLOCKED, AND THE ABSENT REVIEWER IS THE FINDING:**
    `gpt5-6-sol` dropped out on **`budget`** because the doc had GROWN answering *its own* objection
    — the self-selecting trigger exactly — so it was re-run alone at a raised cap ($0.089025,
    `present=true`). **It did not repeat its round-2 objection** (that one is genuinely discharged)
    and instead found a NEW real defect the budget-absent synthesis would have hidden: §2.4's
    *"unknown workbench query parameters are ignored"* is a **SILENT FALLBACK**, so `?paylod=1`
    renders a different view instead of refusing. Both survivors carry concrete reviewer-authored
    fixes and neither disputes the design DIRECTION, so the **narrow-refinement carve-out** applied:
    §2.4 now carries `gpt5-6-sol`'s closed-grammar replacement **verbatim** (five accepted keys;
    unknown key, duplicate scalar key or unsupported combination → `400` with a constant message;
    nothing ignored, no precedence fallback) pinned by `TestWorkbenchRefusalBranches/unknown-parameter`
    and `/duplicate-parameter` with mutations `M31`/`M32`, and `V26` is `gemini-3-1-pro`'s asked-for
    row proving `type readStore interface` (`daemon.go:323`) has exactly five context-first methods.
    Doc now 894 lines; `metered=$0.127437`. **NEXT: sprint-planner.** Prior head text follows.
    ~~[**RATIFIED 2026-08-14 (Mark, attended): OPTION B — DEFER BEHIND ITEM 18.**~~ The A/B filed at
    iter-82 is answered. Item 14 does NOT grow into `host/store` context plumbing; item 18
    `w-daemon-read-cancellation` lands first and gives all eight routes one elapsed-time contract,
    and item 14 then ships its renderer onto that base. Rationale recorded with the decision: the
    unbounded-read defect is pre-existing on all seven existing GET routes, so B *removes* it rather
    than extending it — which satisfies the rejecting reviewer's own `catch` — and item 18 was
    already the next unparked row, so the ordering cost is nil. The reviewer's objection is
    therefore SUSTAINED, and answered by sequencing rather than by scope growth. Residual `WB-R1`
    is discharged by item 18, not by this item. ~~Carried forward independently of the A/B: the
    `Internal` branches passing `err.Error()` verbatim to an unauthenticated localhost client.~~
    **REFUTED BY MEASUREMENT 2026-08-22 (iter-109), and it is item 18's M3 that refuted it.** At
    `93e1ba5` there is exactly ONE `"Internal"` site (`host/daemon/handlers.go:162`) and it passes
    the constant `internalErrorMessage = "internal store failure"` (`:132`), whose own comment at
    `:118` calls it "the ONLY message a 500 ever carries on the wire". All five surviving
    `err.Error()` sites are `"BadRequest"`/400 branches (`:338 :371 :548 :555 :566`); negative
    control on an invented symbol returns 0. Recorded as `V21` in the design doc — a carried-forward
    concern that nobody re-measured for eight days after the landing that fixed it.
    **Unparked; blocked only on item 18.** ~~[PARKED `needs-human-review` 2026-08-13 (iter-82)]~~ — DESIGNED, NOT LANDED. Doc
    `design_docs/planned/w-workbench-read-only.md` (641 lines, designer `codex:gpt-5.6-sol`,
    rotation) is committed; TWO quorum rounds both BLOCKED with all four reviewer slots
    `present=true` (no N−1 degrade). `metered=$0.160575`. R2's surviving objection
    (`gpt5-6-sol`) disputes the design DIRECTION, so the narrow-refinement carve-out is
    foreclosed by its own limb (b) and a third round is not permitted — hence the park, per
    Standing rule 2. `gemini-3-1-pro`'s R2 objection is procedurally right and substantively
    REFUTED (the API's malformed/absent distinction IS modelled — `writeAPIError` 400/404/500
    across all six handlers, `handlers.go:215-419`, control = 6 `handle` funcs), so it needs a
    one-row citation, not a redesign; its unstated half is real, though — the `Internal`
    branches emit `err.Error()` verbatim to the client.
    **THE OPEN ASK — a one-word A/B, framed by the rejecting reviewer's own `proposed_fix`:**
    **(A)** expand this item to carry context-aware store reads for every `/workbench`
    operation, a request-scoped deadline, an explicit timeout status, and a test that reds when
    context propagation is removed — accepting scope growth into `host/store` and past the
    ~1.5–2d estimate; or **(B)** defer this item behind new item **18** and land the daemon
    read-cancellation first. The reviewer's `catch` pre-empts the obvious third answer:
    *"Do not treat the fact that existing JSON reads share this defect as justification for
    extending it."*
    **TWO OF THIS ROW'S OWN INSTRUCTIONS WERE FALSIFIED BY MEASUREMENT AND MUST NOT BE
    RE-IMPLEMENTED.** (i) "renders `UNSUPPORTED`" is DEAD — iteration 79's carve-out cut that
    constructor from the exported grade type, and the landing of item 13 (the very thing this
    row was waiting on) is what falsified it: `grep -rn "UNSUPPORTED" world/` → **0**,
    same-scope control `grep -c "CLAIMED" world/types.ail` → **4**. The ratified type is exactly
    four grades, `PROVEN | TESTED | ATTESTED | CLAIMED` (`world/types.ail:33-38`). (ii) the
    "six routes" count is wrong: `grep -c "mux.HandleFunc" host/daemon/daemon.go` → **8**
    (seven GET + one POST), and the code's own comment ~~at `:456` says
    "seven patterns"~~ — **that half is now stale too: at `93e1ba5` the comment reads "The eight
    patterns below are the complete frozen v1 table" (`daemon.go:550`); `grep -c 'seven patterns'`
    → **0** with the same-file control `grep -c 'patterns'` → **1** (iter-109). The row's count was
    wrong and the comment has since been fixed, so only the ROW needed correcting.** A third measured
    hazard for whoever implements: `gradeOf`'s mapping sends `TestReport(_, _)` to `TESTED`
    **regardless of its pass/fail boolean**, so a FAILING test report carries the same badge as
    a passing one. ~~ **Prior head text follows.** ~~FILED 2026-08-11 (attended, Mark) — THE UI
    PROGRAMME IS CURRENTLY SERIALISED BEHIND WORK IT DOES NOT DEPEND ON, AND THIS ROW IS THE
    PARALLEL PATH.~~]
    **w-workbench-read-only** · clause-5 · a renderer route on `worldd` plus the READ-ONLY
    workbench — world browser, timeline scrubber, provenance walk (SCENARIOS.md scenario 3).
    **THE ARGUMENT, MEASURED:** `host/daemon/daemon.go:461-468` already serves exactly the
    vocabulary a browser needs — `GET /v1/head`, `/v1/worlds/{ref}`, `/v1/objects/{ref}`, `/v1/log`,
    `/v1/log/{index}`, `/v1/registry/{name...}` — so the ONLY backend addition is the renderer route
    itself (the daemon has no HTML/web/SSE route today; that absence is measured, not assumed).
    **Read-only means it depends on NONE of the four things the approval inbox is waiting on**: no
    decision-packet schema, no timeout policy set, no grade mapping, no session-filtered MCP. It
    therefore runs **IN PARALLEL with item 11**, not behind it — different files, no shared gate.
    It also builds the grammar item 7 will inherit (P6's zoom ladder, P2's link discipline) against
    real objects, so the inbox becomes a composition rather than a first attempt; and it de-risks
    the medium decision (localhost surface + optional shell) on the surface where being wrong is
    cheap. **BUILD IT EVEN THOUGH THE GRADES ARE MISSING:** until item 13 lands it renders
    `UNSUPPORTED`, which the anti-pattern list already rules is the CORRECT rendering of an unmapped
    variant — and a screen full of honest `UNSUPPORTED` badges is the most persuasive argument
    available for funding item 13. **TWO LANDMINES, BOTH FIRST-PARTY:** (a) a renderer route means
    new `net/http` surface in `cmd/ailang-worldd`, which is the ONE group the boundary gate's
    per-group `extraForbidden` list exempts for loopback IPC — the exemption is per-group by
    construction since iter-55, so do not collapse it, and re-assert the positive arm rather than
    inheriting it green; (b) `host/boundary` pins `wantFileCount = 1` and any new `.go` file there
    reds `TestBoundaryASTWriteGuard` — this has bitten three times; the fix is a new package, never
    relaxing the pin. **CARRIES `14/OD-12`** (below) as a NON-BLOCKING decision: build on the
    controller default and settle it before item 7 ships. · ~1.5–2d · NEEDS A DESIGN DOC.
15. [**[LANDED 2026-08-14 (iter-85)] — COMPLETE.** PR #67 → squash `aaada20`; Gate 3b GREEN on the
    merge commit itself, SHA-addressed, `present(2) == expected(2)`, control fires; evaluator
    `sonnet` **96/100, zero blocking**; `metered=$0.00`. The v1 `DecisionPacket` is frozen —
    `TimeoutPolicy`, `TimeoutOutcome`, the seven-field packet under semantic ID
    `world/decision-packet/v1`, and **five Z3-proven laws**, doubling the repo's proven identity
    count **5 → 10** with named tests **20 → 39**. Landed as **SIX files, not the planned five**:
    the sixth is `docs/SELF_MOD_PUBLISH.md`, which quotes `contentHash`/`tarballSHA256` verbatim
    in its attended-operator approval table, so reprojecting red-lit `host/runbook` (attribution
    measured — rc=0 on pristine `origin/dev`, rc=1 on the pre-repair tree). 16-arm mutation drill
    run in full plus the VP15 control; **MU2** is the design's spine (a consistent lie across
    contract AND body stays Z3-GREEN and is caught only by the named `outcomeCode_test_5`), and
    **MU15 CANNOT FIRE and is recorded as decoration rather than a kill** (`x <= y-1` ≡ `x < y`
    over ints), with the plan's pre-registered fallback MU15b landing instead. Fabricated
    authority is deliberately NOT mutated — `independentAuthority` is a host-derived bool and no
    pure law can verify a bool's provenance; it stays the declared item-7 host residual.
    Enforcement remains item 7's by explicit deferral: the five laws are proven and, as landed,
    not yet invoked by any host path. HUMAN-SURFACE §7 point 3 and the timeout-set half of point 1
    are now marked CLOSED. Executed by iteration 85 **attempt 1**, which died mid-flight before
    committing (Standing rule 7); attempt 2 verified rather than adopted that work.
    ~~RATIFIED 2026-08-14 (Mark, attended): §7.3 = OPTION A — FREEZE v1 NOW.~~ The A/B filed at
    iter-83 is answered. The v1 `DecisionPacket` freezes with this item: the type and its five
    Z3-proven laws land in `world/types.ail`, the semantic ID `world/decision-packet/v1` is
    reserved, and every later amendment is a NEW version (`/v2`), never an in-place edit — the
    same discipline as `LogHeader` and every other content-addressed wire type. Enforcement
    remains item 7's by explicit deferral. Rationale recorded with the decision: Option B's only
    merit is letting real inbox usage inform the field set, but that usage comes from item 7,
    which is parked behind item 5's upstream blocker with no ETA — so B trades a cheap,
    version-bumpable commitment for an indefinite wait. The specific defect that made
    "Option A as written" unacceptable at R1 (laws that could not see `deadlineAt`, so an early
    timeout satisfied them) was fixed in the R2 revision and is proven, so the reviewer's stated
    ground for deferral no longer exists. **Unparked; sprint-ready, gated on nothing.**
    ~~[PARKED `needs-human-review` 2026-08-14 (iter-83)]~~ — DESIGNED, NOT LANDED. Doc
    `design_docs/planned/w-decision-lifecycle-freeze.md` (694 lines, commit `2104631`), two quorum
    rounds, all four reviewer slots PRESENT. R1 both reject; **R2 `gemini-3-1-pro` PASS,
    `gpt5-6-sol` REJECT**. The surviving objection disputes the design DIRECTION — *"Do not freeze
    v1 until those integration gates exist"* — so the narrow-refinement carve-out fails on limb (b)
    and Standing rule 2 forbids forcing it. **THE PARK IS NOT A STALL: the reviewer's own first
    `proposed_fix` IS the doc's Option B**, so the block collapses onto HUMAN-SURFACE §7 point 3,
    the freeze-timing question the doc already puts to Mark. **ONE-WORD ASK — `option A` = freeze
    the v1 `DecisionPacket` now with this item (five Z3-proven laws land in `world/types.ail`;
    enforcement remains item 7's by explicit deferral); `option B` = ratify only the `TimeoutPolicy`
    set + its resolution semantics now and leave the record unfrozen until item 7 defines the
    enforcement boundary.** Its checkable half is accurate and worth carrying either way: **no
    acceptance criterion requires any of the five proven laws to be INVOKED**, so they are provable
    and inert until a host path exists — the "a guard is not a gate until something reds when you
    remove it" class, aimed at a kernel law rather than a refusal branch. **WHAT THE DESIGN
    SETTLED, and it survives both options:** §3.1's "ledger-recorded creation time" is satisfiable
    with **no time column anywhere** — `host/store/schema.sql` has zero time fields across 8
    `CREATE TABLE`s, and the resolution is the logical-time idiom already ratified twice
    (`journal.go:26-27` *"LogicalTime is supplied by the caller; journal payloads never read a wall
    clock"*, plus `Capability.ExpiresAt`/`PublishApprovalScope.ExpiresAt`). **THE VERIFIER FACT
    THAT OUTLIVES THIS ITEM — the repo's recorded ADT limitation is narrower than the truth, the
    iter-79 lesson one level deeper:** the unencodable shape is NOT "a record containing
    `list[ADT]`" — a record containing a **bare** ADT field fails identically (`unknown sort`,
    `errors=1`, `check.passed` true, rc 0), **and the contract need not even read that field**
    (measured: two files differing in exactly one field, `verified=1,errors=0` vs
    `verified=0,errors=1`). This is what killed `gpt5-6-sol`'s R1 `proposed_fix`, whose literal
    signature `validTimeout(packet, …)` the controller built and measured unencodable — so a
    reviewer's remedy is a claim about the toolchain too, and the carve-out has no slot for a
    correct objection with an unimplementable fix. Two further v0.30.0 limitations measured and
    routed upstream: ADT equality in a general boolean expression needs an `Eq` instance the
    top-level `result == match …` postcondition form does not, and `import std/prelude` is refused
    by `IMP012_UNSUPPORTED_NAMESPACE`. **Gate pins this item moves when it lands (all five in the
    doc's Conflict Surface):** `EXACT_TOTAL_VERIFIED` **5→10**, `EXACT_TOTAL_TESTS` **20→39**,
    `REQUIRED_TESTS` +19 identities, `REQUIRED_VERIFIED["world/types.ail"]` +5, and pin 5's marker
    → `10/10 … across 11 module(s)`; `LEG1_MODULES` unchanged (no module added). Measured green on
    an isolated copy of the real tree: `verified=6, errors=0, cex=0`, `len(tests[])=39, failed=0`.
    Re-priced 1.0d. **Prior filing text follows.** ~~FILED 2026-08-11 (attended, Mark) — ITEM 7's
    REAL PREREQUISITES, NEITHER OF WHICH DEPENDS ON THE TRANSITION REGISTRY.~~ **NOTE, measured at
    iter-83: this row says it "blocks item 7", but item 7's OWN row states its park chain as
    `TR.A2 → TR.B → TR.C → item 5 P6.B → item 7` and does not name item 15 at all — two ratified
    rows disagreeing about item 7's prerequisites. This item supplies content item 7 needs; it is
    not item 7's only gate.**]
    **w-decision-lifecycle-freeze** · clause-5 · enumerate the **typed finite set of timeout
    policies** (`6b/§7` ratification point 1, deliberately NOT enumerated in the doc — candidates
    named by the reviewers: cancel · remain safely unexecuted with bounded escalation · execute only
    if authority was already independently granted) and **freeze the decision-packet schema**
    (point 3 — it becomes a world type, so it is kernel-adjacent and the item-13 manifest gotcha
    applies here too). **WHY IT IS NOT COSMETIC:** HUMAN-SURFACE.md §3.1 is BINDING and entirely
    unimplemented — every packet MUST carry a ledger-recorded creation time, decision deadline and
    timeout policy; DEFER MUST create a new bounded deadline; silence MUST NEVER synthesize approval
    or rejection; and replay MUST reproduce deadline behaviour deterministically from ledger time.
    Without it **the inbox can wedge on a human exactly the way this loop wedged on a background
    agent** — Standing Rule 6 restated at the UX layer, which is precisely how both round-2
    reviewers found it independently across two providers. · ~1d · NEEDS A DESIGN DOC · blocks
    item 7, gated on nothing.
16. [**[COMPLETE 2026-08-13 (iter-80)] — M1 LANDED as PR
    [#65](https://github.com/sunholo-data/ailang-world/pull/65) → squash `d9712dd`, Gate 3b
    GREEN on the MERGE COMMIT itself (`checks=2` = expected 2, `present == expected` asserted,
    0 non-success). Doc + sprint plan moved to `design_docs/implemented/`. Evaluator (sonnet,
    ≠ codex executor) **96/100, zero blocking**. `metered=$0.148857` (quorum R3 only).
    **UNPARKED BY HUMAN DIRECTIVE:** Mark answered the one-word A/B on #53 at
    `2026-08-13T06:12:23Z` with `option A`, so the `killGroup` seam landed in production
    `host/broker/handlers.go` (+7/−1, behaviour-identical; the doc said `+5/−1` five times,
    including in the A/B text put to him — the sha matched N5 throughout, so only the label
    was ever wrong). **M2 (the post-reap re-sweep) REMAINS DECISION-GATED by the doc's §6 and
    is authorized by nothing that landed here**; it unlocks only when one captured diagnosis
    shows the full H1 conjunction. **THE FLAKE IS BOUNDED, NOT DIAGNOSED, AND THAT IS THE
    HONEST CLAIM:** the trap is armed, so the next firing writes a mechanism-grade record
    instead of a mystery. **CARRY FORWARD — the finding that outlives this row: A TEST SEAM
    THAT *REPLACES* RATHER THAN *WRAPS* MAKES EVERY MUTATION OF THE REPLACED BODY VACUOUS.**
    `killGroup` is a package `var` the test replaces, so a recorder that re-implements the
    kill leaves the production body dead for the test's duration — `MUT-KILL-NEUTER`, and
    therefore AC3 (this item's ONLY forced-failure proof), mutated dead code and passed at
    105 ms with the mutation invisible; delegating via a captured `orig` fails at 5.168 s with
    the complete H1 signature. The mutation table named the right file, the right line and the
    right one-line edit the whole time. **Three cannot-fail gates were removed from this one
    item** — the row's own `-count=20` (iter-78), round 3's cold-majority AC5/P3 +
    `MUT-WARM-SKIP` (killed by `gpt5-6-sol`, whose objection landed *inside* the section
    titled "proofs that cannot pass by luck"), and AC2's claim to prove the recorder mutex,
    which `gemini-3-1-pro`'s round-2 BLOCKING objection had demanded: `os/exec`'s `watchCtx`
    sends on `resultc` after calling `Cancel` and `Wait` receives it, a happens-before edge, so
    an unsynchronised recorder measured **0 `DATA RACE` in 20 `-race` runs** with a
    known-positive control firing. **A reviewer's objection survived a park, a human
    ratification and a verbatim adoption while being wrong.** Prior PARKED row follows.**
    ~~[PARKED `needs-human-review` 2026-08-13 (iter-78)] — DESIGNED, NOT LANDED. Doc
    `design_docs/planned/w-broker-base-flake.md` (586 lines) is committed and quorum-reviewed
    TWICE, both rounds BLOCKED, both external reviewers PRESENT in both (`absent_reviewers`
    empty, no N−1 degrade); `metered=$0.202435`. R2's two objections differ in KIND, which is
    why this parks rather than taking the narrow-refinement carve-out: `gemini-3-1-pro`'s is
    carve-out material (the test-side kill recorder is written by `cmd.Cancel`'s background
    context-cancellation goroutine and read by the main goroutine — an unsynchronised race,
    fixed verbatim by a `sync.Mutex`), while `gpt5-6-sol` asks to REVERSE the doc's central
    architectural choice. **THE ONE-WORD DECISION — A or B:** (A) as designed — bring M2's
    `killGroup` seam forward into M1, i.e. `var killGroup = func(pgid int) error` in
    `host/broker/handlers.go` (+5/−1 behaviour-identical lines), so the test can record the
    kill's count, monotonic offset, target pgid and errno and make the ~0.76% event decisive
    the ONE time CI catches it; or (B) keep M1 production-free, narrow its claim to "localise
    the stall to the `Execute` window", and defer mechanism selection to a later milestone
    using an external tracer or a separate design doc. **Measured for the decision (rule 3f,
    iter-78, controller first-party):** `gpt5-6-sol`'s "enlarges the frozen production core"
    framing is WRONG — `CLAUDE.md:25` scopes the frozen core to `tools/launchd/*` plus skills,
    not `host/`; but its CATCH is right, because coding-standards **S3** does bind `host/`
    ("why is this not a package?") and the doc's §10 answers `not applicable in the S3 sense`,
    a dismissal rather than an answer. There is **exactly one** precedent in this repo, the
    same shape, landed through this loop: `host/store/store.go:863`
    `var commitBeforeOutcomeHook = func() {}` (`6811604`, PR #19, `w-store-durability SD.C`),
    justified in its own comment on the same grounds. Both reviewers' safety premise holds
    today: `t.Parallel()` = **0** across **46** `*_test.go` under `host/`+`cmd/` (same-scope
    control `t.TempDir()` = **119**, `grep` rc=1 not rc=2) — but it is not an enforced
    invariant, which is `gpt5-6-sol`'s point. **THE ROW BELOW IS SUPERSEDED IN TWO PLACES BY
    iter-78's HEAD MEASUREMENTS — the rate is 0.76% (1 failure in 132 runs), not ~18%, and
    its `-count=20` acceptance proof would fail to red ~86% of the time, so it is a
    coin-flip gate (S6). Prior row text follows.** ~~FILED 2026-08-11 (iter-73) — MEASURED,
    NOT SUSPECTED, AND IT IS A TAX ON EVERY FUTURE MUTATION SWEEP~~]
    **w-broker-base-flake** · clause-3 · `host/broker` is **not deterministically green at base**.
    Measured by the `TR.B` planner on unmodified `dev`: `TestHandlerTimeoutKillsTheWholeProcessGroup`
    failed **2 of 11** isolated runs (**~18%**), while a control test in the same package passed
    **3/3** and the whole package passed serially in 35.4 s — so the flake is that one test, not the
    package or the rig. Separately measured: **without `AILANG_BIN` the package is red 100% of the
    time** (`TestEpisodeLiveReplayThreeArmsAndEvidence` refuses to run against an unpinned
    interpreter) — correct behaviour, but it means a bare `go test ./host/broker` is uninformative
    twice over.
    **Why this is not a nuisance row.** `host/broker` is the package every remaining clause-3
    milestone mutates (`TR.B2`, `TR.C`, and P6.B after them), and this mission's entire verification
    discipline rests on reading a mutation's RED. An 18% coin-flip inside that package **fakes kills**
    (a mutant "killed" by the flake, when it in fact survived) **and falsifies inverse arms** (an
    inverse that must be rc=0 to prove isolation goes red for the flake and reads as over-coverage).
    It is the same load-confound class `w-bench-load-confound` already studied, aimed at the gate
    instead of at a benchmark. `TR.B1` worked around it with `-skip` in every arm — a workaround
    that is correct per-iteration and silently erodes coverage of a real timeout guarantee if it
    becomes permanent.
    **Deliberately NOT repaired inside `TR.B1`** (the evaluator agreed the scope call was right):
    it is a pre-existing defect, and repairing an unrelated flaky test inside an authority milestone
    is exactly how a sprint's own RED evidence stops being attributable. The item: diagnose the
    ~18% (process-group kill timing under parallel load is the hypothesis, NOT the finding), fix or
    correctly bound it, and prove the fix with a repeat-count run (`-count=20`) that reds before and
    passes after — never by adding a retry or a skip. · ~0.5d · NEEDS A DESIGN DOC (small) · gated
    on nothing; should land before `TR.C` if the queue allows, since `TR.C`'s whole deliverable is an
    assertion in this package.
17. **[LANDED — ITEM COMPLETE 2026-08-22 (iter-108); doc + sprint plan MOVED to
    `design_docs/implemented/` 2026-08-22 (iter-109), plans travelling with their doc.** **[iter-108 record follows.]** **[LANDED — ITEM COMPLETE 2026-08-22 (iter-108). ALL SIX MILESTONES SHIPPED; `PE.F` WAS THE LAST.** PR [#82](https://github.com/sunholo-data/ailang-world/pull/82) → squash [`189299b`](https://github.com/sunholo-data/ailang-world/commit/189299b), Gate 3b GREEN on the MERGE commit (SHA-addressed, `checks=2`, `present=2 == expected=2`, 0 not-green, run CONFIRMED to exist `total=1 event=push`, parent control `total=1` from a **`rev-parse`d** SHA, negative control `total=0`). Judge `sonnet` **96/100 PASS round 1, zero blocking**; two of its three non-blocking findings REPRODUCED first-party and FIXED in a round-3 commit [`16531ea`](https://github.com/sunholo-data/ailang-world/commit/16531ea), the third recorded not patched. **AC8** — a focused `host/evidence` named-manifest leg in `scripts/verify_go.sh` ahead of the broad legs: terminal `Action=pass` top-level identities only, **set equality** against a non-empty `REQUIRED_EVIDENCE_TESTS` (missing, skipped, failed, **duplicate**, **extra**), `EXACT_EVIDENCE_TESTS=37` supplied from the shell so the banner cannot drift from the assertion, and three anti-vacuity floors that fail LOUDLY. Its isolated self-mutation gate `host/verifygate/evidence_manifest_gate_test.go` copies the **live** script and EXECUTES it rather than re-implementing the comparison. **AC9** — the full **27-row** re-drill (`M1–M5 · M7–M14 · M16–M27 · M29 · M30`) at `design_docs/verification/w-validated-proven-evidence-boundary/AC9-mutation-drill.md`: every row with the exact edit, sha256 **LANDED**, `go build ./...` rc=0 asserted BEFORE any test result was read, the red set **enumerated by running** the mutant, sole-killer vs one-of-N (22 sole, 5 wider with every extra member explained), and a byte-identical `cp` restore — never `git checkout --`. **AC12** — zero diff under `host/daemon/`, `cmd/`, `host/replay/` and renderer-shaped paths, with the same instrument firing non-empty elsewhere. **THE FINDING: a removal proves a check FIRES; only an ADDITION proves it LOOKS.** Four live-script arms plus a pristine control — remove a required literal → RED `extra=[…]`; add a bogus required literal → RED `missing=[…]`; shell literal `38` against an intact set → RED `observed_unique=37 exact_required=38`, so the count pin is **independently** load-bearing and not redundant with set equality; and **append one real PASSING test to `host/evidence`** → RED, count **37→38**, which is the arm that substantiates *"PE.F must be the last change to `host/evidence`"* and which no removal-direction mutant in the table could have produced. Neutering the comparison itself reds **3 of 5** then-existing arms, leaving the two anti-vacuity arms green behind their own earlier branches. **THE REFUTED PREMISE IS NOW CORRECTED, NOT JUST RECORDED (V56):** the plan's *"cannot be killed without the REAL store"* list and §6.1's *"M26 — no fake participates in the kill"* are FALSE — reproduced by executor, controller and judge independently, the fake-based `TestConstructorNamesActuallyUsedUnorderedTimeouts` (M26) and `TestConstructorRefusesUnknownBusyTimeout` (M30) red in isolation with no real store; the real-store arms are **integration evidence** that production reports and forwards its live busy window, not a prerequisite for the branch kill. **Twelve §6 divergences found and REPORTED, not patched away** — moved source anchors (`store.go` → `read_object.go`, `report_codec.go` → `envelope_codec.go`, `grade.go` gone), successor failure text (M3 reaches `hash_mismatch`, not `malformed`), and red sets wider than the table names (M20 also reds `TestValidatorMintIdentitiesAreDistinct`); the judge spot-checked four and found all four to be doc drift with no masked code defect. **The judge's two real findings**: §5 names three anti-vacuity floors and only two had arms — the *zero discovered packages* floor was live code no test protected (its message appears once in the script, zero times in the test, and every synthetic writer emitted the package event by construction), closed with a sixth arm that is the SOLE arm to die when that floor is neutered; and the M17 record cited a `.snap/S2/` path the commit had deleted, replaced by a live-script drill reproducible from the commit alone. **Doc + plan still to move to `implemented/` (bookkeeping, next iteration).** Previously: `PE.D` LANDED 2026-08-21 (iter-106), four of six shipped. PR [#79](https://github.com/sunholo-data/ailang-world/pull/79) → squash [`d1b7eae`](https://github.com/sunholo-data/ailang-world/commit/d1b7eae), Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`, run CONFIRMED to exist `total=1 event=push`, parent control `checks=2`). Judge `sonnet` **62/100 FAIL round 1 — and the judge was RIGHT**: two BLOCKING findings, both reproduced first-party, both unpinnable refusal branches, and one of them literally unreachable (`ValidateProof` decoded the report a second time after `DecodeAuthenticatedEnvelope` had already refused on that exact failure, so the arm named for the report-decode guard was killing through the envelope decoder; neutered, the whole suite stayed rc=0). The other: `NewValidator`'s compound four-disjunct guard, **three of four unpinned** by measurement. Both repaired structurally — the double decode removed (the envelope now carries the report it already decoded in an unexported field) and the guard split into four separately-observable refusals, each a sole killer after the repair. **Round 2: 95/100 PASS, ZERO blocking**, the judge aimed explicitly at the repair. **THE FINDING THAT OUTLIVES THE MILESTONE: `M16` is the sprint's one ADDITION mutation and `TestPublicAuthoritySurfaceIsFrozen` was blind to three natural spellings of the forbidden package-level PROVEN resolver** — a `*ResolutionResult` result (`*ast.StarExpr`, not `*ast.Ident`), a method on an exported type (`fn.Recv != nil`, skipped) and a type alias — each minting `ResolvedGradeProven` from a raw `HashRef` with no seal, whole package green. Replaced with an EXACT MANIFEST of the package's 55 exported declarations, complete by construction, with an anti-vacuity floor. A removal proves a check FIRES; only an addition proves it LOOKS. 14 removal mutants drilled outside the sandbox, 11 sole killers; `M9`/`M10`/`M20` broad by construction with every red-set member explained. **NEXT: `PE.E`** (real-store integration proofs, 0.85 d), then `PE.F` last by its own `EXACT_EVIDENCE_TESTS` pin. **Prior head text follows.** ~~`PE.C` LANDED 2026-08-21 (iter-105). PR [#78](https://github.com/sunholo-data/ailang-world/pull/78) → squash `bd48f68`, Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, run CONFIRMED to exist: `actions/runs?head_sha` `total=1 event=push`; the PR head was polled separately and also green); judge `sonnet` **88/100, ZERO blocking**, in its own worktree. THREE of six milestones are in.** New package `host/evidence`: strict canonical `ProofReportV1` codec (nine fields in §3.3 order), the two-member authenticated envelope codec (`report`, `mac`, base64url-no-padding), `ClaimedEvidence`/`DecodeProposal` with its 256 KiB pre-parse cap, the decoded-report/verified-list/string caps, byte-equal decode→re-encode both ways, and the AC19 depth pin. The raw-envelope cap is deliberately ABSENT — it is `ReadObject`'s `maxBytes`, and duplicating it would make M4 observable-identical and therefore unkillable. Closes AC2, AC19, AC3's canonical-bytes half. **THE MILESTONE'S OWN ANTI-VACUITY PIN WAS HOLLOW AS DELIVERED, AND ITS DRILL PASSED ANYWAY — SECOND CONSECUTIVE ITERATION WHOSE REAL DEFECT WAS IN ITS OWN VERIFICATION MACHINERY RATHER THAN ITS PRODUCT.** `DecodeProposal` shipped with an `if err == nil { … }` wrapper plus a trailing `return ClaimedEvidence{}, nil` that is UNREACHABLE in the correct program; its only effect is on the mutant. Measured both ways on the identical tree, exit codes captured without a pipe: as delivered, M27 → named arm **rc=1**; with the unreachable branch removed (natural spelling) the unmutated suite is still **rc=0** and the SAME mutant **SURVIVES at rc=0**. So the kill was a property of dead code, not of the guard. Cause: the arm's observable — *some* typed `malformed` refusal with a zero claim — has a value set LARGER than the mechanism's, because the swallowed decode error falls through to the **trailing-JSON** bystander guard, which produces exactly that. Repaired rather than accepted: the dead branch is gone and the arm now pins the stdlib scanner's own `exceeded max depth` text, which nothing else on this path produces. M27 on the repaired tree — mutation LANDED (`ed757fd9…`→`44a06aac…`), mutant BUILDS `go build ./...` rc=0 asserted BEFORE any test result was read, named arm rc=1 reporting `trailing JSON; want the stdlib scanner's own max-depth refusal`, measured red set size **1**, `-skip` rc=0 sole killer, restore byte-identical. **THE JUDGE'S NON-BLOCKING FINDING THEN SPLIT UNDER FIRST-PARTY MEASUREMENT, IN BOTH DIRECTIONS.** REPRODUCED and bigger than filed: `report_codec.go`'s member-ARITY guard is genuinely unpinned — neutered with `if false && …` the mutant BUILDS and the whole suite stays **rc=0**, and a report with its FINAL field omitted (every remaining member correctly named AND ordered, so no per-index mismatch fires first) then **panics** `index out of range [8] with length 8` under the mutant, in code whose §3.3 mandate is *malformed input → typed refusal, never a panic*. The existing missing-member arm removes a MIDDLE field, so a sibling guard masks the gap. Now pinned by `TestTruncatedTailReportIsRefusedNotPanicked` with an in-test control and a branch-unique observable (only the arity guard reports a member COUNT); sole killer, red set 1. REFUTED: the same finding claimed the ENVELOPE shape guard survives identically — it does not. Neutering the whole compound condition (the judge's own arm, not a partial one) builds and returns **rc=1**, `TestEnvelopeStrictRefusals/unknown` failing on a genuine assertion, not a panic. Also DELETED one undeclared unreachable branch the judge found (`sort.StringsAreSorted` after a loop already enforcing strict pairwise increase — unreachable by transitivity), with the surviving contract re-verified by a firing control: unsorted REFUSED, duplicate REFUSED, sorted-unique ACCEPTED. Executor deviation ADJUDICATED SOUND BY MEASUREMENT, not by its reasoning: a missing or wrong-width `mac` decodes with `MACValid=false` rather than refusing — §3.3's single classification exception — with the control firing (full valid envelope → `MACValid=true`) and no leak past `mac` (absent `report` → typed refusal). Both gates rc=0 OUTSIDE the sandbox on the pinned `v0.30.0` (`e37b370`), measured rc=0 on the PRISTINE tree first, both before and after the judge fixes; the executor's own `verify_go.sh` red was a sandbox loopback-bind denial and was treated as UNINFORMATIVE UNDER SANDBOX. Judge finding 3 (the trailing-JSON sub-branch split is Kind-only-observable) became queue row **30** rather than a silent edit. **NEXT: `PE.D`.** Prior head text follows.** ~~[IN-SPRINT — `PE.B` LANDED 2026-08-21 (iter-104). PR [#77](https://github.com/sunholo-data/ailang-world/pull/77) → squash `3ddacae`, Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, run CONFIRMED to exist: `actions/runs?head_sha` `total=1 event=push`, parent control `checks=2`); judge `sonnet` **91/100, ZERO blocking**, in its own worktree. TWO of six milestones are in.** `ReadObject` reads probe and payload inside ONE transaction on ONE reserved connection, `BusyTimeout()` is cached at `Open`, and B5 taught the DR-2 ratchet to see `ReadObject` (pins unchanged, 8/2/1 = 11). **THE MILESTONE'S TWO REAL DEFECTS WERE BOTH IN ITS OWN ANTI-VACUITY MACHINERY, AND BOTH ARRIVED AS *NON-BLOCKING* JUDGE FINDINGS.** (1) `busyTimeoutFromParams` reported the LAST `busy_timeout` pragma while the pinned driver applies the FIRST — measured both directions against a live `PRAGMA` readback (100/7000 → driver 100 ms, accessor 7 s; reversed → driver 7000 ms, accessor 100 ms; control: a single `busy_timeout(3333)` reads back 3333 ms). Reachable, because `withBusyTimeout` returns early on ANY existing pragma precisely so a caller's explicit value is never overridden — and the under-reporting direction is the unsafe one for AC18/AC22's `ObjectReadTimeout > BusyTimeout()` pin, which would pass while the real lock wait outlives the read deadline. Fixed to first-wins and pinned against the readback rather than against a comment. (2) B3's own text requires that *"a green can never come from a writer that never ran"*; the arm asserted the HOOK fired and merely logged the outcome, so a hook whose `UPDATE` was removed still passed — against the M25 mutant too. The FIRST remedy was itself vacuous (a removed statement yields a nil error, so `committed-after-snapshot` is written either way — rule 3i's *the observable's value set is wider than the mechanism's*, met inside a remedy written for that very class); the shipped assertion switches on `busy-refused` or `RowsAffected() == 1`. **DECLARED LIMITATION:** on this rig the competing write is busy-refused under a rollback journal, so the `RowsAffected` branch is unreachable here and its non-vacuity is UNPROVEN. M25 controller-run OUTSIDE the sandbox: mutant landed (`a25f41ec` → `f1b1881a`), mutant BUILDS rc=0 (the first attempt did NOT — `undefined: sql`), blast radius exactly one test, `-skip` arm rc=0 → **SOLE KILLER**. An executor deviation reported as "none" (two DSN helpers inlined into production-dead state) was adjudicated by measurement: equivalence CONFIRMED over four parameter sets with a firing negative control, and the controller's own vacuity hypothesis **REFUTED** — a mutant dropping the injection killed three tests. Two judge findings NOT closed here became rows **28** and **29** rather than silent edits. **NEXT: `PE.C`.** Prior head text follows.** ~~[IN-SPRINT — `PE.A` LANDED 2026-08-21 (iter-103). PR [#76](https://github.com/sunholo-data/ailang-world/pull/76) → squash `cbd17de`, Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, run CONFIRMED to exist: `actions/runs?head_sha` `total=1 event=push`); judge `sonnet` **96/100, ZERO blocking**, in its own worktree.~~** The document was routed to **sprint-planner** (`opus`, lane derived `opus fail-closed:env-pin`, used verbatim) exactly as iteration 102's NEXT prescribed. The plan splits tranche 1 into **six** CI-green milestones `PE.A`–`PE.F` at the doc's own **4.70 d**, no milestone over 0.95 d; **27 of 27** live mutation rows mapped to exactly one owning milestone, asserted by SET COMPARISON against the doc's own table (the only doc row absent from the plan is **M6**, tranche 3; M15/M28 stay DECLARED gaps). `open_questions_for_the_human` is **EMPTY** — every candidate was answered by a measurement. Plan: `design_docs/implemented/w-validated-proven-evidence-boundary-sprint-plan.md` + `.ailang/state/sprints/w-validated-proven-evidence-boundary.plan.json`. **THE PLANNER REFUTED THE DOC SIX TIMES**, the sharpest being that item 18's DR-2 ratchet `TestNoNewDeadlineFreeStoreReads` **cannot see `ReadObject`**: its detector alternation names exactly five getters, the doc leans on that ratchet **5 times** including a V-row that prints the regexp verbatim as its own observed output, and a future `.ReadObject(context.Background())` would land silently. Reproduced first-party (the pattern matches `.GetObject(context.Background()`, not `.ReadObject(context.Background()`); it is folded into **PE.B as task B5**, so it lands in the next milestone rather than a queue row. Also refuted: AC18's opening still specifies a *live* `PRAGMA` its own round-11b amendment withdrew (binding reading is CACHED); M1's predicted failure text (`got 4, want 1`) is not what the runner emits (`test 0: expected 1, got 4`, reproduced by me on a scratch module); AC16's `context.WithTimeout` control of 8 is **10 non-test / 25 all** at HEAD; and `verify_go.sh` has **no** `REQUIRED_*`/`EXACT_*` manifest pattern for M17 to extend (0 vs the `verify_ail.sh` control's 13) — it must be BUILT. **`PE.A` SHIPPED THE FIVE COUPLED AILANG MOVES IN ONE COMMIT** (§8.1: a lagging projection reds step 3/9, a stale golden reds step 9/9): `Evidence` gains `| ProofReceipt(HashRef)`, `gradeOf` maps it to `CLAIMED` in BOTH the `ensures` postcondition and the body, a seventh `gradeCode` tuple whose emitted identity `gradeCode_test_7` was OBSERVED before being pinned, the projection reprojected byte-identical, and the golden regenerated. **BOTH DEFECTS FOUND SAT IN A REPORTING CHANNEL RATHER THAN AN ENFORCING ONE, WHICH IS WHY NEITHER HAD EVER RED:** `verify_ail.sh` enforced `EXACT_TOTAL_TESTS` from inside a python heredoc while its banner restated the total as a shell literal, so PE.A's own 39→40 move made the gate pass while announcing a number it no longer enforced — fixed durably by hoisting the constant to the shell (setting it to 41 moves enforcement AND banner together; byte-identical after revert); and `host/runbook`'s AC28 gate went RED correctly because the regenerated golden left `docs/SELF_MOD_PUBLISH.md` naming *"bytes that are not the reviewed artifact"* — digests re-derived FROM the golden by command, one flipped nibble reds it. **THE EXECUTOR NEVER SAW THE SECOND AND WAS RIGHT NOT TO:** it labelled the Go leg `UNINFORMATIVE UNDER SANDBOX` (loopback binds denied under `workspace-write`) and the controller's mandatory out-of-sandbox re-run surfaced it — the sandbox HIDING a real failure, the mirror of iteration 100 where it INVENTED two. The plan had named this exact falsifier in A5 and called a red *"a genuine finding, not an expected cost"*. Mutations RED and reverted byte-identical: M1 (both legs, control tests 1–6 stay pass), M18 (step 3/9), M19 (step 9/9). Gates outside the sandbox: `verify_ail.sh` rc=0 (10 identities, **40** named tests, world package 9/9), `verify_go.sh` rc=0 zero FAIL lines, load average 3.41. **NEXT: `PE.B`.** **Prior head text follows.** [ by `D-WORLD-24` resolving as ARM A (Mark, attended, `#68`, bare comment `A`, 2026-08-20T16:04:52Z). NO ASK IS OPEN ON THIS ITEM; THE DECISION LEDGER IS 11 ROWS, 0 OPEN. THE DOCUMENT IS UNBLOCKED AND ROUTABLE TO sprint-planner.** The producer is SHED into new queue row **26** `w-bounded-z3-report-producer`; tranche 1 is now the validator, the seal, the boundary and their gates. Round 13 (designer `codex:gpt-5.6-sol`) applied the ruling in full — §3.4 → a shed marker, `NewProducer`/`Producer.GenerateProof`/`ObjectWriter` out of the live §3.2 surface, AC5/AC20 and M15/M28 removed as DECLARED GAPS with **no renumbering** (the charter and §§10.1–10.15 cite those numbers by name), AC9's required-mutation enumeration updated, §4/§5 weakened to hand-authored-fixture validation only. Both round-12 objections leave UNFIXED with the producer, verbatim, into row 26. **THE ROUND-13 QUORUM BLOCKED — AND THE THREE-ROUND PATTERN BROKE, WHICH IS THE FINDING.** Rounds 11, 11b and 12 each blocked on the PREVIOUS round's own fix. Round 13's two objections do not: both land on text that is **byte-identical to `HEAD`** — `gemini`'s on §3.3 step 10's configured-required-identity sentence (working-tree line 454 hash-matches `HEAD` line 459; round 13 deleted five producer lines above it, so the same-number comparison is the wrong instrument and the designer caught the controller's directive stating it that way), `gpt5-6-sol`'s on §3.2's optional `busyWindowReporter` from round 11b — with the line-486 producer control correctly DIFFERING. These are **pre-existing gaps twelve rounds walked past**, surfaced by the cut rather than caused by it. `absent_reviewers` was NON-EMPTY (`gpt5-6-sol`, `budget`) and the skill's rule fired: re-run alone at a raised cap → **REJECT**, so the round is a 2-present block, not a degraded pass. **CARVE-OUT TAKEN, SECOND USE ON THIS DOCUMENT (round 13b):** both objections are in-tranche completeness with concrete reviewer-authored `proposed_fix` text and neither disputes the design DIRECTION, so parking would manufacture a decision. Applied verbatim: the validator gains a configured `requiredIdentities []string` with an empty-set constructor refusal (`ErrInvalidValidatorConfig`) — step 10 has mandated that check since round 5 with **no parameter to receive it**, and M14 tested it — and `BusyTimeout()` becomes MANDATORY on `ObjectReader`, deleting the type-assert-and-skip whose absence-of-capability was silently read as absence-of-lock-wait; new **AC21/AC22** and **M29/M30**, gaps unfilled. **CONTROLLER ARITHMETIC CORRECTION, and it is the durable half:** §9's tranche table did not sum to its own Total — the delivered revision read **4.35 d** over rows summing **4.60 d**, and the PRE-revision table was already **0.05 d** adrift (rows 5.40 vs stated 5.35, inherited from round 11b) — so two rounds shipped a Total nobody added the column for. Both producer shares are now deducted IN their rows; tranche **4.70 d**, decomposition **11.70 d**, guardrail overage **0.70 d**, each asserted by re-summing the column and quoting the command (V53–V55, §10.16/§10.17). **The number moved AWAY from compliance and is stated, not rounded.** `metered=$0.4345`. **NEXT: route to sprint-planner** — the 0.70 d overage is a pricing fact to plan around, not a new ask; `D-WORLD-24` already settled the scope question. **Prior head text follows.** [**[PARKED `needs-human-review`] 2026-08-20 (iter-101) — THE RULING WAS APPLIED IN FULL AND THE DOCUMENT BLOCKED TWICE MORE, BOTH TIMES ON ITS OWN PREVIOUS FIX. `D-WORLD-24` IS THE MIRROR OF `D-WORLD-22`/`23`, NOT A DUPLICATE: THOSE SETTLED WHETHER A TRANCHE ABSORBS SEPARATELY-OWNED WORK; THIS ASKS WHETHER IT SHEDS WORK IT OWNS.** `D-WORLD-22` arm B applied in full by the round-11 designer revision (`claude:claude-fable-5`, doc 1733 → 1965 lines): the wait-bound claim narrowed to exactly what is proven, OPEN row 22 named as the residual's owner (obligation (i) asserted BY COMMAND — 0 completion tokens after struck spans are stripped, control on LANDED row 21 returns 1), the `busy_timeout < ObjectReadTimeout` ordering PINNED at construction (obligation (ii): AC18/M26/V49 — `host/daemon/handlers.go:299-302` already admitted *"an ORDERING nothing in this code asserts"*), and `gemini`'s owed nesting-depth note MEASURED rather than forwarded (four arms, both toolchains: the panic/CPU half REFUTED — depth 131,071 inside the 256 KiB cap refused in ~0.5 ms — the unasserted-stdlib-internal half SUSTAINED, `maxNestingDepth = 10000` at `scanner.go:148`; AC19/M27/V50). **ROUND 11 BLOCKED, BOTH PRESENT, AND BOTH OBJECTIONS LANDED ON THE ROUND'S OWN REMEDY:** the live-`PRAGMA` `BusyTimeout()` was itself an unbounded wait on the 1-connection pool, and the producer's WRITE side had never been bounded at all. **The narrow-refinement CARVE-OUT applied for the FIRST time on this document** — the scope axis had foreclosed it seven consecutive times, and these two objections are in-tranche completeness with concrete reviewer-authored fixes, so parking would have manufactured a decision. `gpt5-6-sol`'s two arms were chosen between BY MEASUREMENT (§10.4's discriminator): V51(c) shows non-test `host/`+`cmd/` issues ZERO runtime `PRAGMA busy_timeout` operations, so the window is immutable after `Open` and the "live" property it removes is not load-bearing — the alternative removes the wait rather than capping it. `gemini`'s fix applied in all three parts. **ROUND 12 BLOCKED AGAIN, BOTH PRESENT, `absent_reviewers` EMPTY IN BOTH ROUNDS — and for the THIRD consecutive round the objections land on the previous round's fix, one surface over.** Mechanism, measured: item 18 threaded the store's READ signatures and left the waits unbounded (DR-2, 11 sites), `PutObject` was never threaded at all, and the busy window is ordered against nothing — so every store surface this tranche newly touches arrives with the same three holes. **Revision budget spent honestly** (one designer revision + one quorum + one bounded carve-out revision + one confirming re-quorum); a third would be unbounded re-litigation. `metered=$0.8343` for both rounds. **THE ASK IS ONE WORD** (`D-WORLD-24`, doc §10.15): both round-12 objections are PRODUCER-side, so shedding §3.4 dissolves both unfixed — arithmetic §10.13 could not have had. **Owed under BOTH arms:** AC20's decoy arm must say it exercises the connection-POOL wait and not the LOCK wait, and `gpt5-6-sol`'s reader-parameter gap holds under B. **Prior head text follows.** [**[NEXT — UNPARKED 2026-08-20 (iter-100) by `D-WORLD-22` resolving as ARM B, itself a consequence of `D-WORLD-23` arm A (Mark, attended, `#68` comment `A`). NO FURTHER ASK IS OPEN ON THIS ITEM.** The owed work is a bounded revision, NOT a re-quorum of the direction: (a) weaken the tranche's claim to exactly what it proves — every wait THIS TRANCHE'S OWN CODE performs is bounded by `ObjectReadTimeout`, a LOCK-contended wait is bounded by `busy_timeout` (2000 ms, `writer_lock.go:179`), and the composition is safe ONLY WHILE `busy_timeout` < `ObjectReadTimeout`; (b) per arm A's obligation (ii), SHIP AN ASSERTION PINNING THAT ORDERING, which nothing in the tree asserts today — a claim narrowed without a pin is a claim merely made quieter; (c) per obligation (i), re-assert by command that **queue row 22 `w-daemon-lock-wait-not-deadline-bound`** is still OPEN at the moment the revision is written, since it is the residual's named owner. `gpt5-6-sol`'s round-10 objection is SUSTAINED as substantively correct and is answered by narrowing, never by disputing. Also still OWED from round 10 and deliberately unapplied at park time: `gemini-3-1-pro`'s NON-BLOCKING note that `DecodeProposal` caps raw bytes at 256 KiB but states no maximum JSON NESTING DEPTH (§10.6 records that an unforwarded non-blocking note returns as BLOCKING one round later). §9 prices the tranche at 4.75 d against a 3–4 d guardrail — the re-scope is part of the revision, not a footnote. **Prior head text follows.** ~~[**[PARKED `needs-human-review`] 2026-08-19 (iter-96) — `D-WORLD-21` WAS ANSWERED AND APPLIED, THE CARVE-OUT WAS TAKEN ON BOTH ROUND-9 OBJECTIONS, AND `gemini-3-1-pro` PASSED FOR THE SECOND CONSECUTIVE ROUND. THE DOCUMENT PARKS ON A NEW ONE-WORD ASK, `D-WORLD-22`, AND THE REJECT LANDED ON A DISCLOSURE.** Two revision rounds ran this iteration — round 9 (the ruling applied) and round 9b (the narrow-refinement carve-out, both round-9 objections closed with the reviewers' own verbatim text) — plus two quorums, rounds 9 and 10, `absent_reviewers` empty on both, `metered=$0.632600` for the pair. Doc 1295 → 1733 lines, designer `claude:claude-fable-5` ×2 (rotation wrapped again: `codex:gpt-5.6-sol` probe-exhausted until 2026-08-20 05:34 — measured UNPIPED at `rc=1`, after my first reading of `rc=0` turned out to be a `| tail` artifact rather than a codex defect (corrected in the log; "exit codes through pipes lie" reproduced by me in the same iteration whose spine is a grep's zero being an artifact of its own spelling) — and gemini read-only, so the rotation WRAPPED; FLAGGED, instance 4). Commit `e903e98`. **ROUND 9 — the ruling, applied.** The seam becomes `ReadObject(ctx, ref, maxBytes) (ObjectMeta, []byte, error)`; the store enforces `maxBytes` BEFORE materialization via a `length(payload)` probe whose select list omits the payload column, then reads under the supplied context; the streaming reader, the caller-side `io.LimitReader` and the 256 KiB + 1 detection byte are RETIRED. **M22 rewritten, not re-run** (it mutated the `WithTimeout` derivation and was killed by a fake wired to OBSERVE the context — vacuous by construction), and new §6.1 audits EVERY prescribed fake against one rule: **a fake is admissible where it supplies INPUT the mutated mechanism consumes, inadmissible where it supplies the PROPERTY the mutation is supposed to expose.** The owed real-store integration test lands, blocking on the store's measured single pooled connection (V45, `SetMaxOpenConns(1)`) with lock contention and mid-blob interrupt both REJECTED as mechanisms and the reasons named. **THE FACT THE DOCUMENT HAD NEVER STATED, AND IT DISSOLVES THE ROUND-7 SCOPE OBJECTION: ARM A IS ADDITIVE.** `ReadObject` adds a method and changes NO existing signature, so the **13** non-test out-of-tranche `.GetObject(` call sites (6/2/2/1/1/1, reproducing V40 exactly) and the **4** interface-method DECLARATIONS — counted separately, because a declaration read as a call is this mission's recurring enumeration error — are untouched, and §8.2's frozen packages do not move. The round-7 `maxBytes`-on-`GetObject` arm was unsatisfiable **precisely because it was not additive** (V43). **ROUND 9 QUORUM BLOCKED; BOTH OBJECTIONS CARVE-OUT ELIGIBLE; BOTH FIXES APPLIED VERBATIM IN 9b.** `gpt5-6-sol`: §8.2 asserted the `objects` table immutable and insert-only with **no V-row establishing it**, so the probe and the payload statement could observe different states. **Measured first-party before forwarding (rule 3f) and the premise is TRUE** — 3 writes, all `INSERT OR IGNORE INTO objects`; 0 UPDATE/DELETE/DROP; **0 triggers** (control `CREATE TABLE` = 8 in the same schema); **0 FK/ON DELETE/ON UPDATE cascades** (control `NOT NULL` = 26); and no `fmt.Sprintf` in non-test `host/store` builds a statement (all are `Error()` bodies but `writer_lock.go:196`'s `busy_timeout` pragma). The fix still landed in full: one reserved connection inside one read transaction, a concurrent-mutation test (M25/AC17), the V-row (V47), and `SetMaxOpenConns(1)` explicitly DEMOTED to corroboration — a pool property, not a snapshot, and silent cross-process. `gemini-3-1-pro`: §3.2 listed every validator and codec surface and **no producer API at all** while §3.4 and §9 mandate a bounded producer — verified correct; `NewProducer`/`Producer.GenerateProof` added, with all four departures from the reviewer's literal spelling stated and reasoned in §10.11. **AND THE MEASUREMENT CARRIED ITS OWN FINDING, WHICH IS THE ITERATION'S SPINE: A NEGATIVE GREP WOULD HAVE CONFIRMED IMMUTABILITY VACUOUSLY.** `grep -rniE '\bUPDATE[[:space:]]+[a-z_]+[[:space:]]+SET\b'` returns **0 REPO-WIDE, tests included** — while `host/store/store.go` carries **five** real `ON CONFLICT(...) DO UPDATE SET` upserts at `:618/:709/:759/:836/:978`. **The upsert spelling puts no table name between `UPDATE` and `SET`, so a grep aimed at the UPDATE STATEMENT FORM is blind to a genuine mutation path by construction**, and the pattern itself was proven live against a synthetic positive. What establishes the premise is a POSITIVE enumeration of every statement touching the table, read one by one — never a negative grep for the mutation keywords, because a broken pattern and a genuinely immutable table produce the identical zero. **AND THE CONTROLLER THEN SHIPPED THE SAME BLINDNESS IN THE OTHER SPELLING, IN THE SAME BREATH.** The designer REFUTED two of my supplied numbers and I reproduced both: the statements naming `objects` are **NINE, not five** — four `JOIN objects` reads at `journal.go:744/792/918/966` are invisible to the `FROM objects|INTO objects` adjacency pattern I handed over (verified: that pattern reads **0** on all four lines) — and `Sprintf` is **15, not 16**. The insert-only conclusion is unchanged; the enumeration was not. **I caught the `DO UPDATE SET` blindness and simultaneously wrote a `JOIN objects` blindness into my own instrument: one grep, two spellings, and I only saw the half I was hunting for.** **ROUND 10 — `gemini-3-1-pro` PASS (2nd consecutive, and its summary reads as an endorsement of the premise work: *"exhaustively verifies its premises against the repository"*), `gpt5-6-sol` REJECT, ON THE RESIDUAL THE DOCUMENT DECLARED.** A LOCK-contended `ReadObject` returns via `busy_timeout` (2 s), not `ObjectReadTimeout` — iteration 94's measured 2.043 s under a 300 ms deadline, filed then as queue row 22. **The objection is substantively CORRECT and the controller does not dispute it.** It parks because its `proposed_fix` says *"Make deadline enforcement for lock contention part of this tranche … Remove the residual deferral"*: a SCOPE call folding a separately owned queue row into item 17, and at the same time a challenge to a predicate of the ruling on record (`D-WORLD-21` chose arm A on the ground that under it *"cancellation becomes ENFORCEABLE"*). **Sixth consecutive foreclosure of the carve-out on the scope axis.** **THE SHAPE WORTH MORE THAN THIS DOCUMENT: THE REJECT LANDED ON A DISCLOSURE.** Round 8's applied lesson was to STATE residuals rather than absorb them; round 9 stated this one in AC16 in the reviewer's plain sight; and that honesty is exactly the surface the objection attached to. Absorbing it would have drawn no objection and would have been strictly worse — it is the failure iteration 94 filed row 22 to prevent. **A document is not penalised for the defects it hides.** **WHOEVER RESUMES:** `gemini-3-1-pro`'s round-10 NON-BLOCKING note is OWED under both arms and is deliberately NOT applied here (the document is parking; an unreviewed edit made after the round that parked it is the change nobody would review) — `DecodeProposal` caps raw bytes at 256 KiB but states no maximum JSON NESTING DEPTH, so `[[[[…` within the cap can burn CPU or lean implicitly on Go's internal 10,000-deep limit. §10.6 records that an unforwarded non-blocking note returns as BLOCKING one round later. §§10.1–10.11 were verified BYTE-IDENTICAL by `cmp` at every step, each time with a firing perturbation control. §9 re-priced to **4.75 d** tranche / **10.75 d** total, and the doc states in those words that it now exceeds the 3–4 d guardrail by three-quarters of a day. **THE ASK IS ONE WORD** (`D-WORLD-22`, §10.12). **Prior head text follows.** ~~[**[PARKED `needs-human-review`] 2026-08-19 (iter-95) — `D-WORLD-19` WAS ANSWERED AND APPLIED, AND THE DOC IS CLOSER THAN IT HAS EVER BEEN: ROUND 8 IS THE FIRST ROUND WITH THE STORE SURFACE IN SCOPE TO CARRY A REVIEWER `pass`.** Mark answered attended on 2026-08-19T04:57:30Z (#68, "D world 19 - A yes"), so tranche 1 MAY extend `host/store` with a bounded object read and §8.2's frozen package list widens accordingly. Two revision rounds ran this iteration — round 7 (the ruling applied) and round 7b (the narrow-refinement CARVE-OUT, both round-7 objections closed with the reviewers' own converging text) — plus two quorums, rounds 7 and 8, `absent_reviewers: []` on both, `metered=$0.459098` for the pair. Doc 968 → 1295 lines, designer `claude:claude-fable-5` ×2 (rotation wrapped: the next entry `codex:gpt-5.6-sol` probe-FAILED rc=1, quota-limited until 2026-08-20 05:34, and gemini is read-only under CapRemoteSandbox so it cannot edit a file — FLAGGED). **ROUND 8: `gemini-3-1-pro` PASS, `gpt5-6-sol` reject** — the second reviewer flip to pass in this item's history and the first with `host/store` in the tranche. **THE SPINE, AND IT IS THE SAME CLASS THREE ROUNDS RUNNING, ONE LAYER DEEPER EACH TIME: a context PARAMETER is not a bound (round 7) → a `WithTimeout` bounds the OPEN, not the read loop (round 7b) → an `io.ReadCloser`'s `Read` takes no context and need not unblock on `ctx.Done()` (round 8).** The round-7 sentence "item 18 added `context.Context`, a WAIT bound" was **FALSE and CONTROLLER-AUTHORED** — measured first-party: `host/store/context_read_test.go:370` declares `deadlineFreeReadPins` `{approve.go: 8, registry.go: 2, replay.go: 1}` = **11** production store reads that pass `context.Background()` **by item 18's own ratified DR-2 deferral**, pinned by `TestNoNewDeadlineFreeStoreReads`, which reds only if the set GROWS. Item 18 bounded the SIGNATURE. **AND `gpt5-6-sol` LANDED ON THIS DOC'S OWN MUTATION TABLE: M22 IS VACUOUS BY CONSTRUCTION, because its prescribed fake is wired to OBSERVE the context — the mutant dies for a property the real store reader has never been shown to have. That is iteration 92's fixture-vacuity finding arriving INVERTED (there the fake was too weak and the mutant PASSED; here it is too STRONG and the test passes where production would hang), and nothing in this loop audits what a prescribed fake makes true.** Also measured and now closed: `gemini-3-1-pro` **REFUTED ITS OWN ROUND-6 FIX — the one Mark RATIFIED** — because `OpenObject(ref) (io.ReadCloser, error)` drops the context and the alternative `maxBytes`-on-`GetObject` arm changes a signature with **13** non-test out-of-tranche callers (6/2/2/1/1/1; two in `host/daemon` and `host/replay`, which §8.2 itself freezes), so that arm was unsatisfiable by construction and is DELETED. **THE ASK IS ONE WORD** (`D-WORLD-21`, §10.9): the two reviewers now want different things at one seam — streaming to avoid materialization vs a complete-read-under-context to make cancellation enforceable — so **the carve-out is foreclosed for the FIFTH consecutive time, and for the first time by a dispute between the two REVIEWERS rather than between a reviewer and the design.** `gemini`'s round-8 NON-BLOCKING note (stderr needs its own capped reader, or an attacker-controlled checker OOMs the host while the parsed stream stays in cap) was applied VERBATIM in the same commit rather than deferred — §10.6 records the round-5 miss where a non-blocking note left unforwarded returned as BLOCKING one round later. **WHOEVER RESUMES:** `gpt5-6-sol`'s real-store integration test (a blocked read under `context.Background()` relying only on `ObjectReadTimeout`, mutating the ACTUAL cancellation mechanism rather than `WithTimeout`) is owed under BOTH arms and is the arm that would have caught M22. §§10.1–10.7 were verified BYTE-IDENTICAL by `cmp` before each commit; §9 re-priced honestly to 4.25 d tranche / 10.25 d total, and the doc now states plainly that it exceeds the 3–4 d guardrail by a quarter day rather than rounding it away. ~~ **Prior head text follows.** ~~[**[PARKED `needs-human-review`] 2026-08-18 (iter-90) — REVISED TWICE MORE (rounds 5 and 6), STILL NOT LANDED. D-WORLD-17 arm A IS IMPLEMENTED IN THE DOC AND WAS NEVER RE-DISPUTED; the new ask `D-WORLD-19` is a SCOPE question the carve-out forecloses.** Doc 711 → 968 lines, designer `claude:claude-fable-5` ×2, quorum rounds **5 and 6** overall, `absent_reviewers: []` both rounds, `metered=$0.328765`. Round 5 blocked on two completeness defects, both controller-MEASURED and SUSTAINED before being forwarded (rule 3f) and both fixed in one bounded revision: the binding was **forgeable by Go zero values** (exported types + unexported fields stop foreign writes but not `var v Validator; var s ValidatedEvidence`, whose identities compare EQUAL at zero — the round-4 catch reappearing one layer down), fixed by a mint-validity invariant with the dedicated sentinel `ErrUnmintedAuthority` and new RED arm **M21**/AC15; and §3.4 mandated a bounded subprocess runner that §9 simultaneously deferred to tranche 2, resolved by branch (a) with `V35` measuring the repo first (exactly **3** non-test `exec.CommandContext` sites; `runBounded` at `host/broker/handlers.go:88` is unexported and broker-internal, so reuse is REJECTED with a reason rather than deferred). **Round 6 blocked on two more.** 6a (`gpt5-6-sol`) is uncontested and pre-agreed under either arm: `Validator.Resolve(sealed) ResolvedGrade` has **no refusal channel** while every criterion requires it to return an error and no grade — a controller-authored defect, since the round-5 directive prescribed that signature verbatim from the attended ratification's wording and neither the designer nor round 5 questioned it. 6b (`gemini-3-1-pro`) is what parks this: an attacker-supplied `HashRef` to a multi-gigabyte object OOMs the process before any length check, because §3.3's 256 KiB cap runs *after* the store has already materialised the first full copy. Premise re-derived first-party at `03c7892`: `host/store` exposes exactly **two** exported Object methods (`PutObject:443`, `GetObject:467`), `GetObject` returns the whole payload, and non-test `host/store` has **zero** `io.Reader`/`io.LimitReader`/`maxBytes` occurrences (control: **23** exported `Store` methods, so the zero is a measurement). Its `proposed_fix` widens §8.2's frozen package boundary into `host/store` — the declared subject of item **18** — so it is a scope call the controller may not settle: **fourth consecutive confirmation** that a scope/direction dispute forecloses the narrow-refinement carve-out. **CONTROLLER MISS, recorded:** this row's own ratified head already said *"fold in gemini's non-blocking note (`DecodeProposal(raw)` byte bound)"* — the same byte-bound concern, filed non-blocking at iteration 87 — and the round-5 directive never passed it to the designer, so it returned as a BLOCKING objection one round later. **THE ASK IS ONE WORD** (`D-WORLD-19`, §10.6). **Prior head text follows.** ~~**D-WORLD-17 RESOLVED 2026-08-17 (Mark, attended) — A: BIND EVERY SEAL TO ITS MINTING VALIDATOR.** Full ruling in the decision ledger. The revision implements the row's own arm-A spec below: `Validator.Resolve(sealed)` as a METHOD (drop the free `GradeOfValidated`), unexported per-authority identity, and `TestAttackerChosenValidatorCannotMintForHostAuthority` — plus the cross-validator refusal arm as a named RED mutation (a seal minted by validator 2, refused by validator 1), which is what makes "bind" non-vacuous. Self-minting into your own validator is accepted, not denied: no library prevents a caller lying to itself; the enforced property is that no caller can make SOMEONE ELSE'S validator resolve their seal. Production key custody stays with `w-proven-evidence-production-key-wiring`. Also fold in gemini's non-blocking note (`DecodeProposal(raw)` byte bound). ROUTING NOTE: codex exhausted fleet-wide until 2026-08-20 05:34 — designer rotation falls to the NEXT entry, never `$MODEL`. **Prior head text follows.** ~~**[PARKED `needs-human-review`] 2026-08-15 (iter-87) — REVISED TWICE, STILL NOT LANDED.**~~ The~~
    iter-86 revision round RAN: doc `566 → 711` lines, designer `codex:gpt-5.6-sol` ×2, quorum
    rounds **3 and 4** overall, `absent_reviewers: []` both. **Round 4 is the first reviewer FLIP to
    `pass` in this item's history** — `gemini-3-1-pro` PASS; `gpt5-6-sol` reject. `metered=$0.266188`.
    **ALL THREE PRESCRIBED DELIVERABLES LANDED**: the MAC seam per Mark's Option B (HMAC-SHA256 over
    the canonical `ProofReportV1` bytes, constant-time compare, key custody, `§2.4` rewritten rather
    than deleted since it had argued *against* signatures); the V27 repair (`§3.3`/`§3.4` now source
    a sorted identity set from `verify.results[].function` filtered on `status == "verified"`, never
    the bare int); and the negative control (`unauthenticated_report` on otherwise-perfect canonical
    bytes, in both the mutation table and the ACs, with per-refusal `if false && …` mutations).
    **WHY IT PARKS AGAIN — the reviewer's own verbatim alternative was the trap.** Round 3's
    `gpt5-6-sol` reject was REAL and controller-verified: `§7` carried **14** ACs, and AC14 demanded
    a *daemon-owned* key created at *first startup* inside a tranche whose `§8.2` says
    `host/daemon`/`cmd/**` **do not change** and whose `§3.4:284` says the only integration is a
    **library API** — **unsatisfiable by construction** (the iter-81 class; two reviewers had cleared
    it). Its `proposed_fix` offered two arms and the designer took **arm 2 verbatim** — the
    carve-out's whole safeguard — dropping AC14 to **13** criteria and filing the successor
    `w-proven-evidence-production-key-wiring`. Round 4 then rejected THAT, from the same reviewer,
    correctly: with no production composition root the key comes from the caller, so
    `NewValidator(key [32]byte, …)` (doc `:198`, `:211`) is a public constructor taking a
    **caller-supplied** key and `GradeOfValidated(sealed)` (`:201`; called *the sole bridge* at
    `:78`/`:295`) is a **free function** — attacker-chosen key → hand-MACed envelope → fake reader →
    `ResolvedGradeProven`. **Arm 2 and "authority boundary" are incompatible, and only applying arm 2
    revealed it.** PARK, not carve-out: the one revision + one re-quorum were both spent, the
    objection disputes whether tranche 1 is an authority boundary **at all**, and its fix again ends
    in two mutually exclusive architectures.
    **THE ASK, one word — ANSWERED 2026-08-17: A** (attended; ledger `D-WORLD-17`; both arm descriptions kept unstruck because arm A's text is the revision's spec). **A** = bind every seal to its minting validator (`Validator.Resolve(sealed)`,
    unexported per-authority identity, drop the free `GradeOfValidated`, add
    `TestAttackerChosenValidatorCannotMintForHostAuthority`) and keep tranche 1 library-only —
    tranche 1 then ships a real boundary whose ROOT is deferred. **B** = defer authority-bearing
    `ResolvedGradeProven` to a tranche that HAS a production composition root, i.e. re-widen tranche
    1 to the `host/daemon`/`cmd/**` wiring that round 3's arm 1 described. Both are the reviewer's
    own words; choosing between them is the architecture call the controller may not make.
    **WHOEVER RESUMES:** the doc is otherwise re-quorum-ready, and `gemini-3-1-pro` left one
    NON-BLOCKING note — `DecodeProposal(raw)` states no byte bound before parsing, unlike
    `ValidateProof`'s 256 KiB envelope cap. **TWO ROW CORRECTIONS MEASURED AT `bef0153`:** (1) this
    row's "the gate pins this row must move are now `EXACT_TOTAL_VERIFIED=5` / `EXACT_TOTAL_TESTS=20`"
    is **STALE** — at HEAD they are **10** and **39** (`verify_ail.sh:311`, `:350`; `REQUIRED_TESTS`
    holds exactly 39 names), moved by `aaada20`, item 15's landing **one iteration earlier**; the doc
    now targets 39 → 40. (2) the pinned-verifier limitation this item leans on is **BROADER** than
    "a record containing `list[ADT]`": measured across four record shapes in ONE `ai-check` call,
    scalar-only records **verify**, `list[scalar]` **verifies** and a bare ADT *parameter* **verifies**,
    while a record with a **bare ADT field** and a record with `list[ADT]` both **error**
    (`unknown sort`) — and the failing arm's contract never reads the ADT field, so the trigger is
    the parameter TYPE. `rc=0` / `check.passed=true` / `check.error_count=0` throughout: **the failure
    is silent to the exit code**, which is how it sat unproven through two prior quorum rounds. Rows
    V31/V32.
    **Prior row text follows.** ~~[**`[NEXT]` as of 2026-08-14 (iter-86)** — item 18 parked on a
    one-word A/B and item 14 sits behind 18 by Mark's attended ratification, so this is the top
    unblocked row. It needs a **REVISION round on the existing doc, not a new design**
    (iter-84's finding): the MAC seam, the V27 repair reading `verify.results[].function` rather than
    the bare int, and a negative control.]~~
    **RATIFIED 2026-08-14 (Mark, attended): OPTION B — AUTHENTICATE REPORTS WITH A HOST-HELD
    MAC/SIGNING KEY.** The A/B/C filed at iter-84 is answered. `ValidateProof` earns authority by
    verifying a host-issued tag over the canonical `ProofReportV1` bytes; hash recomputation stays
    as the integrity check it is, and the MAC supplies the PROVENANCE it never could. Option A
    (re-execute the pinned checker inside the validator) is REJECTED on cost and architecture, not
    on strength: it puts a compiler+solver on the critical path of every grade resolution — an
    effectful, timeout-prone operation added to a daemon that, per items 14/18, has no bounded-wait
    discipline in its store layer yet. Option C is REJECTED as a resting state; it remains the
    honest description of tranche 1 only if B is scheduled as its immediate successor. Rationale
    recorded with the decision: the gap is provenance, a MAC is a provenance primitive, and the key
    custody objection is bounded — the key is single-host, never crosses a trust boundary, needs no
    rotation protocol for correctness, and its worst-case loss ("stored reports become
    unvalidatable and must be regenerated") is precisely Option A's steady state. The rejecting
    reviewer's objection is SUSTAINED: §11 must gain the negative control its `catch` demands.
    **Whoever resumes this doc: (i) design the MAC seam per this ratification; (ii) apply the V27
    repair — re-point §3.3/§3.4 at `verify.results[]`, which carries per-identity `function` and
    `status`, NOT at the `verify.verified` integer, and NOT at the reviewer's weaker
    `verifiedCount` fix; (iii) add the negative control: hand-author otherwise-perfect canonical
    `ProofReportV1` bytes and require an explicit `unauthenticated_report` result rather than a
    seal.** That control is required under all three options and so is not contingent on this
    answer. **Unparked; needs a revision round, not a new design.**
    ~~[PARKED `needs-human-review` 2026-08-14 (iter-84)]~~ — DESIGNED, NOT LANDED. Doc
    `design_docs/implemented/w-validated-proven-evidence-boundary.md` (566 lines, 28 Verification Log
    rows), commits `169d6bc` → `bc3965d` → `323baf6`, designer `codex:gpt-5.6-sol`. TWO quorum
    rounds, all four slots PRESENT (`absent_reviewers` empty both rounds), BOTH BLOCKED;
    `metered=$0.179422`. **The doc DECOMPOSES the item** — a fresh inventory (which item 13
    explicitly demanded before pricing) refutes this row's ~1.5–2d: **3.5d** for tranche 1, **8.5d**
    across three ordered documents (`w-validated-proven-evidence-boundary` →
    `w-validated-replay-evidence-boundary` → `w-proven-evidence-renderer-consumption`).
    **THE ASK, one word: how does `ValidateProof` earn authority?** A content-addressed report is
    not an AUTHENTICATED one — every field it checks is a public value an attacker can encode into
    canonical bytes. **A** = re-execute the pinned checker inside the validator, stored reports
    become non-authoritative cache. **B** = MAC/sign reports with a host-held key. **C** = ship as
    designed and record the forgery route as a declared limitation. **V28 prices the premise
    without settling it:** the daemon has 7 GET routes + one write (`POST /v1/commit`) and NO
    object-write route (control: 8 of 8 registrations enumerated), but `PutObject` has 10 non-test
    call sites (control `GetObject` → 16) and `host/broker/broker.go:289` stores bytes derived from
    an effect result — so "writable object store" is not excluded by the transport. No probe
    demonstrates a full forgery. **TWO ROW CORRECTIONS MEASURED AT `4557262`, both by this
    iteration:** (1) this row's claim that the `PROVEN` prohibition is "wired into neither
    `verify_ail.sh` nor CI" is **FALSE as stated** — `:266` pins `gradeOf` in the Leg-1
    required-verified manifest and `:339-340` pin all six `gradeCode_test_N` in `REQUIRED_TESTS`
    (control `EXACT_TOTAL_VERIFIED` → 4); only the AC7 `PROVEN` grep is absent, so the prohibition
    rests on six **gate-pinned** integer expectations rather than on nothing. (2) "no `.ail` module
    reads `Evidence` at all" is stale by construction — item 13 landed `Evidence`, `EvidenceGrade`
    and `gradeOf` into `world/types.ail`; the surviving true part is that `gradeOf` has **no
    caller** outside its own module and its byte-identical projection. **AND THE FACT THIS ROW
    NEVER HAD:** the package's four "exports" are **MODULES**, so publishing `world/types` publishes
    every `Evidence` constructor AND `gradeOf` — a foreign `.ail` module minted `PROVEN` from the
    literal digest "i-made-this-up" (check rc=0, inline PROVEN test PASSING; control `IMP010` fired),
    which is why round 1's kernel-arm design was withdrawn. The gate pins this row must move are now
    `EXACT_TOTAL_VERIFIED=5` / `EXACT_TOTAL_TESTS=20` (not 4/14), the projection, the frozen
    4-export manifest and the ready-packet golden. **Prior filing text follows.** ~~FILED 2026-08-13 (iter-79) — NOT a new idea, it is item 13's DECLARED RESIDUAL, and it is
    filed here because a residual recorded only in a design doc's `## Related` section is invisible
    to the next implementer, which is the exact failure item 13's own §2.4 was written to avoid.**]
    **w-validated-proven-evidence-boundary** · clause-5 · the validated boundary and first real
    producer path that makes the top trust grade `PROVEN` **honestly** reachable. **WHY IT EXISTS:**
    item 13 ratifies a total, Z3-proven `Evidence → EvidenceGrade` mapping over the five existing
    constructors and **deliberately leaves `PROVEN` unreachable**, because round 1's proposed
    `ProofReport`/`ReplayReport` carriers were defeated by `gpt5-6-sol` — an agent authors
    `Proposal.evidence`, so an unvalidated carrier lets it mint the top grade from an arbitrary
    `HashRef`, converting a representation gap into a **grade-laundering authority gap**, and
    HUMAN-SURFACE names grade laundering the cardinal sin. Item 13 therefore types the result and
    declares **mint authority** as an open obligation; this row owns it. **MEASURED at `2ef2271`,
    and the numbers are why this is a separate item rather than a §7.2 sub-clause:** there is **no
    production `Evidence` constructor or decoder anywhere in Go** (0 non-test hits under `host/` +
    `cmd/`; same-call control, the same pattern restricted to `_test.go`, returns **13**), **no Z3
    proof-report producer** (every non-test `Z3` hit is a comment about a mirrored predicate;
    control **424** non-test `hashref` hits), and **no `.ail` module reads `Evidence` at all** —
    `verify()` at `world/transitions.ail:45` touches only `proposalMatchesWorld` and every
    construction site passes the empty list (control `.stateRoot` = **6**). `host/replay` has real
    `DivergenceError`/`KindHashMismatch` (`replay.go:111`, `:285`) but no Evidence output path. So
    this item builds the repo's **FIRST** `Evidence` producer plus a validating decode boundary —
    not a wiring job, which is precisely the scope-expansion item 13 refused to absorb silently.
    **INHERITED FROM ITEM 13'S LANDING (iter-81), measured — THE `PROVEN` PROHIBITION HAS NO
    PERSISTENT GUARD, AND THE PROOF CANNOT SUPPLY ONE.** A **consistent** `=> PROVEN` arm placed in
    BOTH the contract and the body of `gradeOf` leaves Leg 1 **fully green** (`verify.verified=1`,
    `errors=0`, `counterexample=0`) — Z3 has nothing to object to, because the contract and the code
    agree with each other. Only the six hand-authored integer expectations in Leg 2 red. And the
    doc's AC7 grep that was meant to backstop this is wired into **neither** `scripts/verify_ail.sh`
    **nor** CI — it was a one-shot sprint-acceptance command, and **zero** `*.go` files name
    `gradeOf`/`gradeCode`/`EvidenceGrade`. So today the top grade is kept unreachable by six test
    literals and nothing else. This is the mission's own "a guard is not a gate until something reds
    when you remove it" arriving inside a PASS: an acceptance grep proves a property held once, on a
    tree that no longer exists, and reads in the plan exactly like coverage. **This item owns closing
    it** — whatever mint authority it defines needs a gate that reds, not a grep that ran.
    **ACCEPTANCE SURFACE, carried verbatim from the reviewer so the descope loses nothing:** define
    explicit mint authority and a validated/opaque value unavailable to proposal authors; bounded
    `HashRef` loading, hash verification, typed report decode, and an explicit successful
    proof/replay result; an explicit error result on every failure with **no fallback grade**;
    grading accepts only that validated value; both real producers wired through it; and mutations
    proving **arbitrary, missing, malformed, mismatched, failed and divergent** reports cannot yield
    `PROVEN`. **CARRIES ITEM 13's GOTCHAS FORWARD:** it is kernel-adjacent, so the five-pin move
    applies (`EXACT_TOTAL_VERIFIED`, `EXACT_TOTAL_TESTS`, the `packages/world-core` projection at
    Leg 3 step 3/9, the frozen 4-export manifest at step 4/9, and the byte-for-byte ready-packet
    golden at step 9/9); and a contract may **not** take `Proposal` or anything reaching
    `list[Evidence]` (measured: `unknown sort` on a record containing `list[ADT]`, `verify.errors=1`,
    which reds Leg 1). · ~1.5–2d · NEEDS A DESIGN DOC · **gated on item 13 landing** (it ratifies the
    grade type this boundary returns). Not promoted: item 13's sprint, then 14/15 by their existing
    order — this row is filed at normal position, not ahead of anything.
18. [**ITEM COMPLETE — M3 LANDED 2026-08-18 (iter-93), ALL THREE MILESTONES SHIPPED. PR [#71](https://github.com/sunholo-data/ailang-world/pull/71) → squash `d21754f`, Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`); evaluator `sonnet` **96/100**, ZERO blocking findings. Doc + sprint plan → `design_docs/implemented/`.**] **w-daemon-read-cancellation** · clause-2 · M3 greens **AC5** and closes the item: `internalErrorMessage = "internal store failure"` as the ONLY message a 500 carries on the wire; `Config.ErrorLog io.Writer` with nil → `os.Stderr` **resolved once in `New`** (so the default is assertable rather than a nil branch nobody can observe); `d.writeInternalError(w, r, err)` writing one line per error carrying the route — never to `announce`, whose extra lines were measured deadlocking `Run` against an `io.Pipe` (V17). `grep -c 'err.Error()' host/daemon/handlers.go` **11 → exactly 5**, and all five survivors ARE the BadRequest sites (311, 344, 521, 528, 539 at head). 5 files `+338/−10`.
    **THE SPINE — THERE IS A *SEVENTH* `Internal` 500-ECHO SITE, AND AC5, THE MILESTONE'S OWN HEADLINE GATE, IS BLIND TO IT BY CONSTRUCTION.** The design doc (§2.6, §3) and the plan (T3.2) both say *"the six Internal branches"*. Measured at base `e4ba56d`: `handlers.go` has **6** and `daemon.go:563` (`handleHead`) has **1** — the daemon has **SEVEN**. The doc's count was taken *inside* `handlers.go` and then quoted as a property of the DAEMON — rule 3b(ix) exactly, landing on the item's headline AC — and **AC5's grep is file-scoped to `handlers.go`**, so no reading of it can ever see the seventh. Proven, not argued: under **MU7c** (restore `err.Error()` at `daemon.go:597`) the mutant BUILDS, `TestInternalErrorsAreSanitized` KILLS it `rc=1` on the `/v1/head` arm, the `-skip` inverse is `rc=0` — and `grep -c 'err.Error()' handlers.go` **still reads exactly 5**. **AC5 PASSES ON A TREE THAT LEAKS.** Implemented to the letter, M3 would have shipped that green with `/v1/head` still echoing the DSN path (the executor demonstrated the leak live against a corrupted store). The doc is self-inconsistent about it too: its M2 row already treats `handleHead` as a read route. The seventh site is sanitized here.
    **THREE MUTATION ARMS, ALL KILLED**, each with a sha256 landed-proof pair, a build check on the mutant, a `-run` KILL arm and a `-skip` INVERSE arm, restored by `cp` from an out-of-repo backup and verified byte-identical (`git checkout --` never used): **MU7** (prescribed, `handlers.go:322`) rc=1 on the BODY assertion; **MU7b** (executor-added: delete only the log write, keep the sanitized body) rc=1 on the LOG assertion **only**, which is what proves the two writes are asserted independently; **MU7c** as above. MU7 and MU7c were re-run INDEPENDENTLY by the controller. Note **MU7b and MU7c both leave AC5's grep reading 5** — the grep sees only the leak direction, never "the detail was destroyed rather than routed", and never anything outside `handlers.go`.
    **THE EVALUATOR FOUND A SURVIVING MUTATION AND IT WAS REAL** (`sonnet`, 96/100, zero blocking). `assertErrorLogLine` matched the route with a bare `strings.Contains(line, "GET /v1/log")`, and `GET /v1/log` is a PREFIX of `GET /v1/log?from=0&limit=5` — so mutating `r.URL.Path` → `r.URL.String()` leaks client-supplied query text into the operator stream, contradicting the function's own doc comment, and the assertion still passed. Reproduced first-party in the prescribed order before being acted on (mutation LANDS `cc9f206c…`→`f334eebb…`, BUILDS rc=0, SURVIVES the whole `host/daemon` suite rc=0), then closed **while the mutation was still applied** (`want + ":"`, anchoring on the producer's `"%s %s: %v"` format → rc=1 on the one route carrying a query), then production code restored byte-identical and MU7c re-confirmed still killed. Landed as `158c500`.
    **CONTROLLER RE-RAN EVERY GATE** (generator≠judge): `go build`/`go vet` rc=0, `gofmt -l host/ cmd/` empty, `go test ./... -count=1` **rc=0, 17 packages, zero FAIL**, `verify_ail.sh` rc=0 with pins **UNMOVED** (10 identities / 39 named tests / 9-of-9 steps) and **0** `.ail` files changed (control: 4 files changed total). AC3 **17** `=== RUN` against M2's head of 13. QUICKSTART re-executed verbatim per S7 §1–§5, plus a live probe of the new behaviour: a corrupted store answers 500 `internal store failure` on three read routes while stderr carries one detail line each and **stdout stays at exactly one line** (the announce).
    **EIGHT FURTHER DOC DEFECTS RECORDED (E10–E17** in `sprint_w-daemon-read-cancellation.json`**)**, of which one spawns queue row **21** and one is worth naming here: **AC5's grep COUNTS COMMENTS** — the executor's first draft read 7/5 because two doc comments mentioned `err.Error()`. That is iter-92's *a comment can quote a code state* class, now on the milestone's headline AC, and it means any future editor who merely MENTIONS the token in prose reds AC5 for a non-defect. Also: AC4 stays unsatisfiable as a grep (`daemon.go` `r.Context()` is **0** at base AND head, by design — the single `readCtx` helper), so "greens AC5 and the full set" is false as literally written; its persistent form MU2/`TestDaemonReadDisconnect` passes, and the AC's grep form needs amending, not the code. **ROUTING, FLAGGED:** codex probed `rc=1` first-party (exhausted until 2026-08-20 05:34) and `pi:deepseek` failed for the **FOURTH** consecutive time with **zero bytes changed** — see `D-WORLD-20`. Executor ran on **opus**, the chain end. **Prior head text follows.** ~~**M2 LANDED 2026-08-18 (iter-92) — PR [#70](https://github.com/sunholo-data/ailang-world/pull/70) → squash `b3c5de0`, Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`); evaluator `sonnet` **99/100**, ZERO blocking findings. ONE MILESTONE REMAINS (M3) — the doc stays in `planned/`.** M2 greens **AC3 (minus sanitize), AC4′, AC7**: `readDeadline = 10 * time.Second` as a constant plus a `Daemon.readDeadline` field wired by `New`; `readCtx`'s body becomes `context.WithTimeout(r.Context(), d.readDeadline)`; the five-method `readStore` seam; timeout classification with **`ctx.Err()` as the authority** and `errors.Is(err, context.DeadlineExceeded)` as a second arm, emitting **503 / class `Timeout`**; the frozen sketch gains `Timeout(string)` and its `=> 503` arm, mirrored by a Go test that PARSES the sketch; and `TestBoundedWaitsAndBodyLimit` gains a seventh literal row plus a `New`-wiring assertion. 6 files `+988/−62`, including a 544-line `host/daemon/read_deadline_test.go`. AC3 enumerates **13** `=== RUN` against a base of **5**. The sketch edit reproduced the planner's **pre-registered post-edit sha256** `cb7f7f89…` byte-for-byte. **E3 WAS SETTLED BEFORE ROUTING, AND IT WAS BROKEN IN BOTH DIRECTIONS — iter-91 recorded only one of them.** AC4 demanded `r.Context()` ≥ 1 in both files. Measured at `05e79e6`: `handlers.go` raw **2**, but one hit is inside a doc COMMENT, so comment-stripped it is **1** — and that line was landed by **M1**, making the arm **already green at M2's base** (vacuous, rule 3e(a)), not unsatisfiable; `daemon.go` **0**, and it stays 0 by design because the single-`readCtx` design means call sites spell `d.readCtx(r)` (genuinely unsatisfiable). **AC4′ replaces it:** (a) comment-stripped `grep -c 'context.WithTimeout(r.Context(), d.readDeadline)'` in handlers.go **0 → 1**, same-file stripped control `r.Context()` = 1; (b) `grep -c 'readDeadline' daemon.go` **0 → 7**, control `drainTimeout` **5 → 6**; (c) MU1 and MU2 must both kill. **AC4′(a) CARRIES A MEASURED TRAP:** the RAW grep for M2's target token reads **1 at base AND 1 at head**, unchanged across the whole milestone, because M1's doc comment spelled the string M2 was going to add — **a comment can quote a FUTURE code state, so a grep-shaped AC can be pre-satisfied by the very comment that promises the work.** **THE ITEM'S SHARPEST FINDING — THE DOC'S OWN PRESCRIBED FAKE CANNOT KILL THE DOC'S OWN MUTATION (E7).** §2.5 says `blockingStore` getters `return ctx.Err()`; MU3 disables only the `ctx.Err()` arm of the two-arm classifier, so the fake hands the *surviving* `errors.Is` arm exactly the value it needs (`context.DeadlineExceeded`) and the assertion cannot fail. Reproduced in both arms, controller-first-party: MU3 + the doc's literal fake → `-run 'TestDaemonReadDeadline/blocking-store' -v` **rc=0, 2 `=== RUN`, PASS — MU3 SURVIVES**; MU3 + the shipped interrupt-shaped sentinel → **rc=1**, body `{"class":"Internal","message":"store: query interrupted (SQLITE_INTERRUPT)"}`. Implemented verbatim, M2 would have shipped a green suite recording MU3 as killed. **SEVEN MUTATION ARMS, ALL KILLED** (MU1, MU2, MU3, MU8, MU9, MU11, MU14), each with sha256 PRE/POST landed-proof, a compile check, a `-run` KILL arm and a `-skip` INVERSE arm; `git checkout --` never used. **MU10 IS CI-RED AFTER ALL — a declared one-sided pin turned out two-sided:** `verify_ail.sh` is rc=0 under it (the plan's Leg-1/Leg-2 boundary claim is exactly right and re-verified), but `TestTimeoutStatusMirrorsSketch` parses the sketch rather than restating its literals, so the Go leg reds **rc=1**. **MU8 RE-CLASSIFIED from wiring-only to semantic (E9):** `readDeadline: 0` yields an already-expired context, so every read route 503s and the inverse arm reds 8 top-level tests. **CONTROLLER RE-RAN EVERY GATE** (generator≠judge): build/vet rc=0, `gofmt -l` empty, `go test ./... -count=1` **rc=0, 17 packages, zero FAIL**, `verify_ail.sh` rc=0 with pins **UNMOVED** (10/39/9-of-9). Six executor deviations DV-1…DV-6, all checked by command and all held; DV-1's strictness claim verified — `grep -c 'err.Error()' handlers.go` is **11 at base and 11 at head**, so M3's six-branch sweep to exactly 5 survivors is untouched. **ROUTING, FLAGGED:** codex probed `rc=1` first-party (exhausted until 2026-08-20) and `pi:deepseek-v4-flash-0731` failed for the **third** consecutive time with zero bytes changed — by a **new** mechanism (325 MB of pure `thinking_delta` NDJSON, `agent_end`=0, `stopReason":"length"`=**0**, killed by the disk ceiling at `rc=137`), which meets the charter's ≥3-datapoint bar; executor ran on **opus**. **NEXT for this item: M3** — `internalErrorMessage` + `Config.ErrorLog`, the six `Internal` branches down to exactly 5 `err.Error()` survivors (the BadRequest sites), and the QUICKSTART S7 re-execution. Gated on nothing. **Prior head text follows.** ~~**M1 LANDED 2026-08-18 (iter-91) — PR [#69](https://github.com/sunholo-data/ailang-world/pull/69) → squash `7ad24ea`, Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`); evaluator `sonnet` **92/100**, one non-blocking finding, fixed in-iteration. TWO MILESTONES REMAIN (M2, M3) — the doc stays in `planned/`.** M1 greens **AC1, AC2, AC6, AC8, AC9**: five getters context-first (`QueryRowContext`), `busy_timeout(2000)` injected into the production DSNs unless the caller set one, the 86-site migration, and the §2.8 ratchet — 28 files `+252/−126` plus a 482-line `host/store/context_read_test.go`. Daemon behaviour is UNCHANGED, which is what made M1 independently landable. §3.1's trap taken as designed: `readCtx` lands as a **cancel-only** helper derived from `r.Context()`, so the six daemon sites never spell `context.Background()` and the AC9 pin holds at **11** across both milestones — M2 changes only the helper's body. **T1.7's pre-registered outcome: (a) RETRY-WINS** — the row returns at **343.7 ms** (lock held 300 ms, busy_timeout 2000 ms), and it is a real retry rather than a fast path because under MU5 the identical read is refused at **20.458 µs** with `SQLITE_BUSY`. **Nine mutation arms, all killed**, each with a sha256 landed-proof read BEFORE the result and a `cp` restore verified byte-identical (`git checkout --` never used); MU6 was independently re-run by the controller. **THREE INSTRUMENT FAILURES CAUGHT MID-SWEEP, each of which would have scored a vacuous kill, and all three are ONE SHAPE — a guard whose reference moves with the thing it measures:** `timeout(1)` is NOT installed on this rig, so the first MU4a–e sweep scored `rc=127` with `=== RUN`=0 (five vacuous "kills", caught only by the enumeration assert); **MU6 SURVIVED its first run** because the test compared the live PRAGMA readback against `busyTimeoutMillis`, the very constant the mutation moves; and T1.5's first control routed through the getter under test, so under MU4a–e it died by the 180 s global timeout instead of its own assertion. **FIVE DESIGN-DOC DEFECTS found while executing (E1–E5), recorded in the sprint JSON.** **E3 GATES M2 AND MUST BE SETTLED BEFORE IT ROUTES: AC4 is UNSATISFIABLE AS WRITTEN** — it demands `r.Context()` in both `handlers.go` and `daemon.go`, but §2.3's single `readCtx` helper means call sites spell `d.readCtx(r)`; M2 must either satisfy daemon.go's arm another way or record AC4's grep form as unsatisfiable and lean on MU2's `TestDaemonReadDisconnect` (pushing `r.Context()` out to the six call sites would green the grep but turn MU1/MU2 into six-line mutations). E1 `SelectedHead` delegates to `selectedHeadTx`, shared with `Store.Commit`, so T1.1 there is a structural interface change and not a line edit — `Commit` keeps its signature and today's behaviour. E2 the 86-site census counts CALL SITES only and misses 5 production + 7 test-fake interface declarations. E4 V25's line number is wrong: `GetVerifyResult`'s caller is `replay.go:191`, and `:153` is the `GetObject` call the ratchet pins. E5 `store.Open`'s `:memory:` carve-out bypasses `writeDSN`, so T1.6 must be file-backed; `journal_mode` is `delete`, which is the only reason `BEGIN EXCLUSIVE` excludes readers and T1.7's stimulus works. **ALSO CORRECTED, and it could not have redded a gate:** the doc says the read seam is "six-method" twice while its own next clause says "all five getters" — measured, it is **6 call sites → 5 distinct getters** (`GetLogEntry` serves two routes; the sixth site is `SelectedHead()` at `daemon.go:497`, outside `handlers.go`). No AC counts the seam's methods. **ROUTING, FLAGGED:** the executor fell through the ratified chain to **opus** — codex probed `rc=1` first-party (exhausted until 2026-08-20) and **two `pi:deepseek` runs returned `rc=0` having changed ZERO bytes**. **NEXT for this item: M2.** **Prior head text follows.** ~~**D-WORLD-18 RESOLVED 2026-08-17 (Mark, attended) — A: SHIP THE SCOPED ITEM AS DESIGNED. UNPARKED DIRECTLY TO SPRINT-PLANNER (doc §13: "A unparks this doc to sprint-planner as written" — NO designer round). The next unit of work is the sprint plan, gated on nothing.** Full ruling in the decision ledger: the 7 daemon GET routes + the 5 zero-cost fold-ins + `TestNoNewDeadlineFreeStoreReads` pinning the 11 residual sites (approve 8, registry 2, replay 1 — follow-on progress mechanically observable 11 → 0). The reviewer's store-boundary guard stays the follow-on's declared closing move, landable exactly when the ratchet reads zero; the two policy questions (what bounds an attended approval *designed* to wait on a human; `Commit` atomicity vs deadline) travel with `w-bounded-waits-operator-and-write-paths` and return to Mark before it routes. **Prior head text follows.** ~~**[PARKED `needs-human-review`] 2026-08-14 (iter-86) — DESIGNED, TWO QUORUM ROUNDS, BOTH~~
    `blocked`, AND THE ONE SURVIVING OBJECTION IS A *DIRECTION* DISPUTE THE CONTROLLER MAY NOT
    SETTLE.**~~ Doc `design_docs/planned/w-daemon-read-cancellation.md` (673 lines, designer
    `claude:claude-fable-5`, 9 ACs, 29 verification rows, 12 mutation rows; `metered=$0.2299`).
    **The defect was live-reproduced at HEAD `6fd26f0` before any routing** — all five of the
    iter-82 measurements below still hold, each with a firing control — so this is not a ghost.
    **R1 (`$0.0985`, both reviewers present)**: gemini caught a **timer leak** (`readCtx` returns a
    `CancelFunc`; the doc never mandated `defer cancel()`) and a `blockingStore` that **overrides
    one getter while driving all six routes**, which read through five distinct getters — both
    verified first-party, both fixes applied **VERBATIM**, the cancel obligation promoted from note
    to gate (`TestReadCtxCancelledAfterHandler` + MU11). gpt5 caught a real **over-claim**
    ("unbounded reads unrepresentable at the type level") — conceded outright, not defended — and a
    watchdog that could red without releasing the blocked getter (fixed for every watchdog, each
    with a named release mechanism and exit bound). **R2 (`$0.1314`, both reviewers present)**:
    gemini's reject was a **PREMISE** objection, so rule 3f says measure rather than forward —
    **both premises measured TRUE** (`handleLogRange` writes exactly once via a terminal
    `writeJSON`, every error path returning first; `defaultClientTimeout = 30 * time.Second` at
    `daemon.go:110`), i.e. a documentation gap with **zero design change**, discharged as rows
    V28/V29. **WHAT PARKS IT:** `gpt5-6-sol`'s surviving objection — an item filed under clause-3
    cannot claim "every wait is bounded" while `POST /v1/commit` waits indefinitely on the single
    connection and 11 sites keep `context.Background()`. **Its premise is TRUE and the doc already
    records it**, so this is not a factual dispute; it disputes the **DIRECTION** (where the item's
    scope boundary belongs), which forecloses the narrow-refinement carve-out — the THIRD
    consecutive confirmation of iter-82's rule, and the cleanest instance yet, since every other
    objection was resolved in-loop (two verbatim, two conceded, two measured away). The revision
    first folded in everything measurably cheap: **5 of the 16** non-daemon sites take the caller's
    real context at zero API cost (`transitionreg.ReadSnapshot:70`/`Publish:221` already carry a
    `ctx`; `Session.invoke` drops one at `broker.go:173` — both confirmed first-party), a ratchet
    test pins the remaining **11** so the deadline-free set can only shrink, and the residue is a
    **named follow-on item** `w-bounded-waits-operator-and-write-paths`. **THE ASK IS ONE WORD — ANSWERED 2026-08-17: A** (attended; ledger `D-WORLD-18`)
    (doc §13): **A** = ship the scoped item (7 daemon GET routes + the 5 fold-ins + ratchet +
    follow-on), **1.5 d**, unblocks item 14 now, residual exposure tracked · **B** = re-size to
    repo-wide bounded waits, **≥2.5 d**, blocks item 14 a further day — and **B needs two human
    policy calls FIRST**, which is precisely why the controller did not just adopt the reviewer's
    fix: *what bounds an attended approval deliberately designed to wait on a human* is a decision
    about the World's authority model, not plumbing, and `Commit`'s atomicity-vs-deadline is the
    second. Measured cost of B (V27): the store-boundary guard breaks all 11 residual sites the
    moment M1 lands; they sit in 8 functions behind 3 exported entry points with **40** test call
    sites. **Also corrected here — the row's own enumeration below is ONE SHORT:** there are five
    context-free getters on the daemon read path (`GetRegistryHead:628` is the fifth) and **six**
    in the store (`GetVerifyResult:773`, sole caller `replay.go:191`, correctly off the daemon
    path). **Prior row text follows.** ~~[**FILED 2026-08-13 (iter-82) — filed on FIRST-PARTY MEASUREMENT, not on a reviewer's say-so:
    the objection that produced it was measured before it was believed, and the measurement made
    it worse than filed.**] **w-daemon-read-cancellation** · clause-3 · bound the elapsed time of
    the daemon's READ path, which is unbounded today on every route. **THE DEFECT, MEASURED at
    `9491a10`, four commands each with a control:** `grep -c "context.Context" host/store/store.go`
    → **0**, same-file known-positive control `grep -c "func (s \*Store)"` → **14**, so the zero is
    an instrument reading and not a failed grep; all four read getters are context-free
    (`GetObject:467`, `GetWorld:522`, `GetLogEntry:551`, `SelectedHead:802`); `grep -n "r.Context()"
    host/daemon/handlers.go` → **0**; and — the fact neither the reviewer nor the design doc had —
    `grep -rn "busy_timeout" host/store/*.go | grep -v _test` → **0**, the pragma existing ONLY in
    `host/store/writer_lock_test.go:609`, so not even SQLite lock acquisition is bounded by the
    production DSN. **WHY THIS IS NOT ALREADY COVERED, AND THE SHAPE IS WORTH MORE THAN THE BUG:**
    `TestBoundedWaitsAndBodyLimit` (`daemon_test.go:202`) is a REAL, ratified, non-vacuous D7 gate
    — it pins six transport constants as literals (5s/30s/30s/120s/30s/10s) and asserts every
    `http.Server` wait is set and non-zero. It is also completely silent about handler execution,
    because `http.Server.WriteTimeout` is a deadline on the connection's WRITE side and cannot
    cancel a goroutine blocked inside a store call. So the daemon has a bounded-waits discipline
    that stops exactly where the store begins, and the presence of a green, well-designed gate is
    precisely what makes the gap invisible — **a gate's coverage is a property of the layer it
    observes, and this one observes the transport while the wait happens below it.** Standing rule
    6 (every wait is bounded) is a mission axiom, and it is currently unenforced on all seven
    existing `GET /v1/*` routes, not merely on the proposed workbench route. **SCOPE:** thread a
    request-scoped deadline from the handler through the read path (context-aware store reads or an
    equivalent bounded mechanism), return an explicit timeout status rather than hanging or
    falling back, and — non-negotiable, since this repo's standard is that a guard is not a gate
    until something reds when you remove it — a blocking-store/lock-contention test that reds when
    the propagation is removed. **CARRIES A SECOND, INDEPENDENT FINDING from the same pass:** the
    JSON handlers' `Internal` branches pass **`err.Error()` verbatim to the client**
    (`handlers.go:220,247,277,325,351,419`), i.e. raw internal error text on an unauthenticated
    localhost surface; decide sanitize-vs-expose here rather than letting a renderer inherit it.
    · ~1–1.5d · NEEDS A DESIGN DOC · gated on nothing. **Item 14 is parked pending a one-word A/B
    that may make this row its prerequisite (option B) or fold it into item 14 (option A) — this
    row stands on its own evidence either way, because the defect is pre-existing and independent
    of whether the workbench is ever built.**~~]

19. [**LANDED 2026-08-19 (iter-94) — PR [#72](https://github.com/sunholo-data/ailang-world/pull/72) → squash `6c2a537`, Gate 3b GREEN on the MERGE commit (SHA-addressed, `present=2 == expected=2`, both `success`); evaluator `sonnet` **97/100**, ZERO blocking findings. Doc + plan: `design_docs/planned/w-daemon-timeout-test-flake{,-sprint-plan}.md`.** THE FIX: `1 * time.Nanosecond` is a **future** deadline — `context.WithTimeout` arms a `time.AfterFunc` — so a fast store read finishes before the timer runs, and since `timedOut` sits on the ERROR path only, the route answers 200. One named constant `expiredReadDeadline = -1 * time.Nanosecond` takes `context.WithDeadline`'s `dur <= 0` branch and cancels SYNCHRONOUSLY at construction. **THE ROW BELOW UNDER-COUNTED THE BLAST RADIUS: there are TWO tests on the one stimulus**, the second being `TestDaemonReadDeadline/real-store-expired-deadline`. Base rates under the pinned toolchain: **6/1000** and **3/500**; head **0/2000** and **0/1000**, zero `200, want 503` in 3000 runs. Also landed: `TestExpiredReadDeadlineExpiresAtConstruction` (reds on the SIGN, zero timing dependence — it makes a sub-1% race killable at `-count=1`) and a `LIMITATION(w-daemon-late-read-503)` comment naming the residual. Two quorum rounds both BLOCKED then the narrow-refinement carve-out; **every objection was a PREMISE objection measured rather than forwarded**, and gpt5's round-2 `proposed_fix` second arm produced row **22**. **Prior head text follows.** ~~[NEXT] **w-daemon-timeout-test-flake** · clause-2 · `TestTimeoutStatusMirrorsSketch`
    (`host/daemon/read_deadline_test.go`) is **intermittently RED at base**, and it is the test that
    pins M2's whole 503/`Timeout` contract against the frozen `.ail` sketch. **Measured by the
    controller at iteration 93, attribution done BEFORE blaming the change (rule 3d):** on the
    **clean main checkout at `e4ba56d`**, `go test ./host/daemon -run
    'TestTimeoutStatusMirrorsSketch' -count=20` fails **1 in 20**, on the `world` route; in the M3
    worktree `-count=10` failed **1 in 10**, on the `registry` route; and the whole package
    `-count=3` was clean. **Route-agnostic and timing-driven, therefore pre-existing and NOT M3's
    doing.** The failure shape is always the same: *"a timed-out read answered 200, want the
    sketch's 503 for Timeout"* — the stimulus is an already-expired deadline, and the read
    sometimes completes before the deadline is observed, so the route answers the real body. This
    is the item-16 class (`host/broker` ~18% base flake) arriving in the package item 18 just
    shipped, and it will red CI intermittently for every future iteration until fixed — a red that
    a future controller will be tempted to misattribute to their own diff. **SCOPE:** make the
    stimulus deterministic rather than racing (the deadline must be observed before the read can
    answer), and prove the fix with `-count` high enough that the measured rate would have shown.
    **Do NOT weaken the assertion** — a 200 here means the timeout contract genuinely did not
    hold. · ~0.25 d · queued iter-93 (controller-measured).~~]
20. **[LANDED 2026-08-20 (iter-100) — ITEM COMPLETE.** PR [#74](https://github.com/sunholo-data/ailang-world/pull/74) -> squash `912009d`, **Gate 3b GREEN on the MERGE commit** (SHA-addressed, `present=2` == `expected=2`, both `success`, run CONFIRMED to exist `total=1 event=push`) -- the PR's own green was on a different SHA and was not banked as the landing. Evaluator `sonnet` **93/100, ZERO BLOCKING**, in its OWN worktree. Controller re-ran every gate OUTSIDE the executor sandbox (mandatory, `host/` diff): `verify_ail.sh` rc=0, `verify_go.sh` rc=0, gofmt empty, capsule **11 RUN / 11 PASS / 0 SKIP** with `no tests to run` = 0 and a nonsense-name control at 0, broker rc=0 with **0** loopback denials. **UNPARKED BY `D-WORLD-23` ARM A** (Mark, attended, 2026-08-20T08:01:31Z, `#68` comment `A`): the round-2 `gpt5-6-sol` SCOPE objection is sustained as correct and answered by keeping scope, weakening the claim to what is proven, and naming row 24 as the residual's OPEN owner -- no third quorum purchased. `gemini-3-1-pro`'s round-2 PREMISE objection was MEASURED not forwarded (rule 3f), all three claims TRUE, fix applied verbatim -- **and running the command the reviewer asked for produced a correction the objection did not predict: `io.LimitReader` is at `capsule.go:238`, not `:237`, which is the `func readCapped(` line transcribed onto the CALL.** **THE FIX:** the flaky test's only witness that an over-limit child was KILLED was a stopwatch started before `New(...).Run(...)` -- outside the region `ExecTimeout` governs -- so it timed resolve plus a 91.8 MB read+sha256 and attributed the total to the child. Three milestones: extract the post-`Start` lifecycle behind an injectable seam (`df0414d`), witness the kill through a fake child's COUNTER in four deterministic arms (`b6a0d6c`), retire the oracle (`9690b9f`). Ordering FORCED, not stylistic: under M1 the old test fails only at the `elapsed >= clock` line that milestone C deletes. **THE PLANNER'S MEASUREMENT, WHICH THE DESIGN DID NOT CARRY: 4 of the 7 mutants (M2/M3/M4/M7) SURVIVE the entire existing suite at base**, so four arms are the sprint's real teeth rather than decoration; all 7 killed on the final tree, each asserted to BUILD and VET before its red was read. **THE EXECUTOR'S SELF-REPORTED DEVIATION WAS LOAD-BEARING AND WAS ADJUDICATED IN BOTH ARMS, NOT ON ITS REPORT** (rule 3h): its first AC5 draft observed both phases but did not FORBID an early `Wait`, so M7 survived; it added a `select`/`default` guard. Measured -- assertion present + M7 -> AC5 FAILS (`wait-entry phase happened before reader release`); assertion REMOVED + M7 -> AC5 **PASSES**. Without the deviation AC5 was a vacuous arm. **THE EVALUATOR THEN CORRECTED THIS ITEM'S OWN RESIDUAL, and the controller reproduced it before accepting the non-blocking label:** the retired fixture was doing TWO jobs -- a bad wall-clock oracle (correctly retired) and emitting 64 KiB so the child genuinely BLOCKED in `write()` (real coverage, silently lost with it). Pipe capacity measured **65536 B**; the surviving test emits **513 B** against a 64 B cap; `readCapped` reads only `limit+1`. So nothing in the suite drives a blocked child being unblocked by `Kill()` -- filed as new row **25**, NOT row 24 (different mechanism). Non-blocking because M-A is behaviour-preserving, confirmed by M6's kill set matching the predicted 7 tests across 2 packages exactly. **Prior text follows.** ~~[PARKED `needs-human-review` 2026-08-20 (iter-99) on `D-WORLD-23` — DESIGNED, and the quorum BLOCKED TWICE, the second time on the SCOPE axis. Doc `design_docs/planned/w-capsule-output-cap-load-flake.md` (`591c16d` + revision `40b7f19`), designer `codex:gpt-5.6-sol`, rounds 1 and 2 both with BOTH reviewers present, `metered=$0.108951`. The row's OWN filed fix is REFUTED-AS-FILED by controller measurement — see the design doc and iteration 99's STATUS. Owed on unpark: `gemini-3-1-pro`'s round-2 premise rows, all three MEASURED by the controller and TRUE.] w-capsule-output-cap-load-flake** · clause-2 · **IT DID NOT REPRODUCE IN 32 EXECUTIONS, SO THE OBVIOUS ACCEPTANCE CRITERION IS VACUOUS AT BASE — WHICH IS THE ONE THING RULE 3e FORBIDS SHIPPING INTO A SPRINT.** Measured at `b5ddf0e`: **7** full-suite runs (3 under the ambient toolchain, 4 pinned `GOTOOLCHAIN=go1.25.6`), **10** isolated under 2x CPU oversubscription, **15** under 6x — **zero** failures in all 32. The row's rate below rests on **n=2**; pooled with these it is 1 in 9, not 1 in 2. **NOT A GHOST, and not fixed**: the code is byte-identical to the base it was filed at (`git rev-parse e4ba56d:host/capsule/capsule.go` == `HEAD:…` and likewise for `capsule_test.go` — a POSITIVE identity check, not a negative log grep), and the two-caps race is readable in the fixture (`ExecTimeout: 5s` and `MaxOutputBytes: 1024` in one `Config`). Package wall-time per iteration rose **1.04s → 2.26s → 4.57s** from idle to 6x load — **but that figure includes the per-run fixture, which archives a 91.8 MB binary, so it is NOT a measurement of the quantity `ExecTimeout` bounds**; it is consistent with the mechanism and does not establish it. **RE-SCOPE: the observable is the MARGIN, not the outcome.** An AC of the form "the flake no longer reproduces" cannot fail at base and must not be written; what a sprint can assert is the child's elapsed time against `ExecTimeout` (instrument the test to record it and red when the margin falls below a threshold), which turns a load-dependent coin flip into a deterministic assertion. **AND THE ROW'S PREMISE CARRIES NO TOOLCHAIN**, which in this repo is load-bearing: a committed canary (`host/store/toolchain_canary_test.go`) and `verify_go.sh`'s deny-list exist precisely because go1.26.0–go1.26.5 miscompile `host/store` on darwin/arm64 — and the rig's ambient Go IS **go1.26.4**, so iteration 97's first three suite runs red on the canary 3/3 until pinned. Any future measurement of this row states its toolchain. **Prior text follows.** ~~`TestF6OutputCapKillsChildBeyondOnePipeBuffer`
    (`host/capsule`) is **load-dependent**: at base `e4ba56d` it failed **1 of 2** full-suite runs
    (`error after 19.33s = *capsule.TimeoutError capsule: execution exceeded wall-clock limit 5s,
    want *OutputLimitError`) and passed **5/5** run in isolation. Two caps race inside one fixture —
    the output cap the test is about, and a **5 s** wall-clock limit — and under whole-suite load
    the child produces output slowly enough that the wall clock wins. Same class as row 19 and as
    item 16, and the same hazard: it makes `go test ./...` non-deterministic, which is the command
    every milestone's "nothing lands red" gate runs. **SCOPE:** decouple the two caps so the arm
    under test cannot be pre-empted by the unrelated one (raise the wall clock well past any
    plausible load, or drive the output cap by a stimulus that does not depend on child throughput),
    and prove it under deliberate load rather than on an idle machine. · ~0.25 d · queued iter-93
    (controller-measured).
21. **[LANDED 2026-08-20 (iter-98) — ITEM COMPLETE.** PR [#73](https://github.com/sunholo-data/ailang-world/pull/73) -> squash `9fa2647`, **Gate 3b GREEN on the MERGE commit** (SHA-addressed, `present=2` == `expected=2`, both `success`, run confirmed to exist `total=1 event=push`) -- the PR's own green was on a different SHA and was not banked as the landing. Evaluator `sonnet` **97/100, ZERO BLOCKING**, having independently reproduced all 5 mutations and both gates. Controller re-ran both gates OUTSIDE the executor (mandatory for a `host/` diff): `verify_ail.sh` rc=0, `verify_go.sh` rc=0, `FAIL`=0 with same-call controls at 34 `ok` / 3 `PASSED`, toolchain pinned `go1.25.6`, 0 hits for go1.26 -- and, because a package-level `ok` is satisfied whether or not the new tests exist, the 7 new tests were asserted to actually RUN (7 `=== RUN` / 7 `--- PASS` / 0 "no tests to run", with a nonsense-pattern control confirming the instrument does emit `[no tests to run]`). Doc + plan: `design_docs/planned/w-archive-stderr-in-manifest{,-sprint-plan}.md`. **TWO NON-VACUITY DEFECTS IN THE DESIGN'S OWN ARMS, BOTH FOUND BY THE PLANNER AND BOTH REPRODUCED FIRST-PARTY BEFORE ROUTING.** (1) **AC7's prescribed deadline fixture reds against CORRECT code.** `sleep 10` inside the fake makes the sleeping process a GRANDCHILD; `exec.CommandContext` SIGKILLs only the `sh` it started, and the grandchild holds the stdout pipe so `cmd.Wait()` blocks on the copy goroutines regardless. At a 200 ms bound, 3/3: **10.211 / 10.081 / 10.164 s**; with `exec sleep`, 3/3: **202 / 201 / 203 ms**. Under the original fixture the correct implementation and the arm's own M3 mutant are INDISTINGUISHABLE (~10 s both) -- vacuous in one direction, a false red in the other. Its 200 ms test bound is separately a flake: a cold first exec measured **227 ms** inside a must-succeed control. Ships at 1 s / 20 s / 8 s, `-count=5` 5/5. (2) **The CONDITIONAL self-heal -- `gemini-3-1-pro`'s round-2 `proposed_fix`, applied VERBATIM under the narrow-refinement carve-out -- shipped with NO acceptance criterion and NO mutation**, so an unconditional heal passed AC1-AC7 unchanged (`grep -c 'strings\.' host/archive/archive.go` = **0** at base; control `host/daemon/daemon.go` = 2). The carve-out is working as designed -- it applies a reviewer's text without re-litigating direction -- but nothing in it asks WHAT GUARDS THE FIX, and a reviewer's own words are the last text anyone thinks to audit. Closed by AC8, an execution-counting fake, and a **dual** mutation pair: **M4a** (`if false && ...`) and **M4b** (`if true || ...`), because `if false &&` cannot neuter a *skip*. Controller-reproduced with byte-identical sha256 revert: under M4a the `NoProbe` arm stays PASS; under M4b it is the ONLY arm that reds. **DECLARED RESIDUAL, NOT ABSORBED:** the bound covers a single-process interpreter; a grandchild inheriting stdout could still hold `Wait()` open, which needs `cmd.WaitDelay` -- scoped out in the shipped constant's own comment. **Prior head text follows.** ~~[DESIGNED -- quorum-clean via the narrow-refinement carve-out; READY FOR sprint-planner. Doc 808 lines, commit `6a811e1`; designer `claude:claude-fable-5` x2 (rotation WRAPPED, Fable diet exceeded -- FLAGGED); quorum rounds 1 and 2 plus a recovered-reviewer solo re-run, `metered=$0.2164`.]~~ **w-archive-stderr-in-manifest** · clause-2 · **REPRODUCED FIRST-PARTY AT `b5ddf0e` WITH THE POLLUTED ARTIFACT READ OFF DISK.** Measured with separate files (a first attempt using `2>&1 >/dev/null | wc -c` misread stderr as 231 B when it is 63 B — shell redirection is an instrument): `--version` stdout **168 B** beginning `AILANG v0.30.0`, stderr **63 B**, combined **231 B** whose FIRST LINE is the Observatory log line. On disk: `/private/tmp/world-demo.db.artifacts/interpreters/sha256/e9746fef…/manifest.json` carries `"version": "2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\nAILANG v0.30.0\n…"`. **THE DESIGN FOUND A THIRD PERSISTED CONSUMER THIS ROW NEVER NAMED**: `releaseFromVersion` takes the manifest's FIRST LINE, so the log line becomes the epoch-1 candidate `registry.Bootstrap` writes, and Bootstrap fatally refuses divergent heads — repair is cheap now and store-bricking later. **THE AUDIT IS CLOSED BY POSITIVE ENUMERATION** of all five non-test `exec.Command*` sites in `host/`+`cmd/`, cross-checked against the five files importing `os/exec`: sites 1 (`archive.go:391`) and 2 (`pkgproj.go:219`) use `CombinedOutput()` and are fixed; site 3 (`broker/handlers.go:93`) merges DELIBERATELY for marker classification and is a named non-change; sites 4 and 5 (`capsule`, `replay`) are already correct and `replay` is the model the fix adopts. **This row's own "control: 0 `.Output()` calls" is the negative form iteration 96 ruled out and is NOT relied on.** Quorum applied: a bounded probe (documented `probeTimeout`, `exec.CommandContext`, separate buffers, `DeadlineExceeded` surviving the `KindExecFailure` wrap) and a CONDITIONAL self-heal (probe only when `!strings.HasPrefix(m.Version, "AILANG v")`), which together restore a true zero-process no-op on the idempotent path. Store population is now scoped, not rig-wide: **53** SQLite files enumerated across four roots at depth 6, **0** with an adjacent `.artifacts` tree, hence zero worldd stores in the searched roots. **OWED, DECLARED, NOT ABSORBED:** `gpt5-6-sol`'s option (2) — rollout safety without a zero-store assumption, checking a target store's registry state before healing — belongs to the tranche that wires this into a store with existing epochs. **Prior text follows.** ~~**A FOURTH SITE OF THE ITERATION-89 STDERR-MERGE
    CLASS, AND THE FIRST THAT PERSISTS INTO STORED STATE.** `host/archive/archive.go:391` uses
    `cmd.CombinedOutput()` to capture the pinned interpreter's `--version`, so the `Observatory:
    NNNMB (warn threshold: 200MB)` line that every `ailang` invocation writes to **stderr** is
    captured as version data and written into the archive **manifest** — from where `GET /v1/health`
    serves it as `interpreter_version`. Observed live by the M3 executor during the S7 QUICKSTART
    re-execution: `interpreter_version` came back as
    `"2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\nAILANG v0.30.0\n…"`. Control:
    **0** `.Output()` calls in that package. Iteration 89 closed three sites of this class in
    `verify_ail.sh`/`host/verifygate`, all of which were **transient** — a bad parse in a process
    that then exited. This one is different in kind: **the polluted bytes are committed to the
    content-addressed archive and outlive the process that produced it**, so a future replay reads
    a version string that is partly a log line, and the pollution is a function of the rig's
    `~/.ailang/state` size rather than of anything in the world. **SCOPE:** `.Output()` (stdout
    only) at that site, an audit for the same shape across `host/`, and a test that reds when
    stderr is merged back in — plus a decision on whether already-written manifests need
    repair. · ~0.5 d · queued iter-93 (executor-found, controller-confirmed).
22. **w-daemon-lock-wait-not-deadline-bound** · clause-2 · **THE READ DEADLINE DOES NOT GOVERN LOCK
    WAITS — `busy_timeout` DOES, AND NOTHING ASSERTS THE TWO ARE ORDERED.** Surfaced by
    `gpt5-6-sol`'s round-2 quorum objection on item 19 (it refused to let the design claim a
    general bounded-wait property from a CPU-bound probe) and MEASURED by the controller when the
    reviewer's own `proposed_fix` told it to. Probe, iteration 94, first-party
    (`modernc.org/sqlite v1.54.0`, `GOTOOLCHAIN=go1.25.6`, file-backed DSN with
    `busy_timeout(2000)` and `journal_mode(delete)`, a writer holding a lock for 10 s so whichever
    bound fires is identifiable): a `QueryRowContext` under a **300 ms** context deadline returned
    at **2.042929083 s** carrying `context deadline exceeded`. Control, same call: the unlocked
    read returns in **448.458 µs**, so the fixture is sound. **The path IS bounded — but by
    `busy_timeout` (~2 s), not by the request deadline, which it overran 6.8× while still
    REPORTING as a context error.** Contrast the CPU-bound case, which is bounded correctly: a
    ~200,000,000-step recursive CTE cancels **2.13 ms** after a 300 ms deadline (control arm: the
    same query runs the full **20.001 s**). So the daemon has two read-wait regimes and the
    deadline governs only one. Today the composition is still safe because `busyTimeoutMillis =
    2000` (`host/store/writer_lock.go:179`) `<` `readDeadline = 10 * time.Second`
    (`host/daemon/daemon.go:128`) — **two independent unconditional constants with no code linking
    them** (evaluator-confirmed). Raise `busy_timeout` above `readDeadline` and a lock-blocked read
    silently overruns the deadline it is supposed to obey. **SCOPE:** decide whether the deadline
    should dominate (cap the effective busy wait by the remaining context budget) or whether the
    ordering is merely asserted and pinned; either way a test must red when the two constants are
    reordered. Item 18 declared `LIMITATION(w-daemon-late-read-503)` for the *other* residual (a
    read completing before cancellation still answers 200) — this row is the second half, and the
    two should be read together. · ~0.5 d · NEEDS A DESIGN DOC · gated on nothing · queued iter-94
    (reviewer-provoked, controller-measured).

23. **w-store-deadline-free-residue-owner** · clause-2 · **ITEM 18's RATCHET NAMES A FOLLOW-ON ITEM
    THAT DOES NOT EXIST, SO ITS 11 → 0 HAS NO OWNER.** `host/store/context_read_test.go:361-372`
    declares `deadlineFreeReadPins` `{host/broker/approve.go: 8, host/registry/registry.go: 2,
    host/replay/replay.go: 1}` — **11** production store reads passing `context.Background()` by
    item 18's ratified deferral DR-2 — and its own comment says the set "may SHRINK … but it may
    never GROW", so that "the follow-on item's progress [is] mechanically observable, 11 → 0, and
    … the store-boundary reject land[s] exactly when this reads zero", and the ratified
    `D-WORLD-18` row likewise calls the guard "the **follow-on's** declared closing move, landable
    exactly when the ratchet reads zero". **Neither names an ITEM, and no queue row is it.**
    Measured at `fd76fa0`: the charter carries **3** hits for the residue tokens, at lines 678
    (the `D-WORLD-18` ledger row), 701 (its ruling paragraph) and 3364 (item 18's own record) —
    **zero of them a queue row**; the three open rows 20/21/22 are the capsule flake, the
    stderr-merge manifest site and the lock-wait bound. Same-file control `w-daemon-lock-wait`
    → 1, so the instrument can see a queue row when one exists.
    **INSTRUMENT NOTE, recorded because the first version of this row was WRONG and the loop's own
    rule caught it:** the claim originally read "grep … returns ZERO" from a pattern that in fact
    returns **3** — and its known-positive control returned 3 as well, i.e. the control and the
    check were matching THE SAME LINES, so a broad pattern manufactured a false zero and a
    confirming control at once. Rule 3a(i-d) says scope the control to the same path as the check;
    it does not say what to do when the control is scoped so *identically* that it stops
    discriminating. A control that can only fire on the check's own hits is not a control.
    So the ratchet is a perfectly good counter with nothing on the other end of it: it
    will red if anyone ADDS a deadline-free read and will sit at 11 forever otherwise. This is the
    blocked-row class inverted — not a row nobody re-measures, but a *measurement nobody owns*.
    Found because item 17's round-8 quorum forced the controller to establish what item 18 actually
    bounded; the answer is the SIGNATURE, and the 11 are the visible residue. · ~0.5–1d · the work
    is threading a real caller context through 11 sites and shrinking the pin map in the same diff
    (the ratchet's comment says this is a one-line edit per site), then flipping the pin to an empty
    map so a future deadline-free read reds on an EMPTY expectation rather than on a non-zero one.
    · NEEDS NO NEW DESIGN DOC — item 18's implemented doc already ratifies the target state.
    · **Gated on `D-WORLD-21`** only insofar as item 17's seam may add a 12th site; if arm B wins,
    that seam carries its own deadline and this row is unaffected.

24. **w-host-subprocess-cleanup-boundary** · clause-2 · **THE OVERFLOW KILL IS NOT PROCESS-GROUP-WIDE
    WHILE THE CANCELLATION KILL IS — IN BOTH PACKAGES, IN THE SAME FUNCTION, TWELVE AND FORTY LINES
    APART.** Controller-measured at `47e12cc` (iteration 99) by POSITIVE enumeration of every kill site
    in non-test `host/` — four, no more: `host/capsule/capsule.go:165` `syscall.Kill(-cmd.Process.Pid,
    SIGKILL)` (the ctx `Cancel` path, group-wide) vs `:193` `killOnce.Do(func() { _ = cmd.Process.Kill() })`
    (the OVERFLOW path, direct child only); and `host/broker/handlers.go:82` `killGroup(pgid)` —
    commented "the cancellation kill boundary" — vs `:122` `cmd.Process.Kill()` (again the overflow
    path). BOTH files set `SysProcAttr{Setpgid: true}` (control: 2 non-test hits), so the group exists
    and only one of the two paths in each file uses it. **Capsule's own comment at `:158-160` states the
    reason** — *"kill the whole process group, or a forked grandchild keeps the inherited pipes open and
    outlives F5/F6"* — attached to the path that got it right. This mission's own named recurring shape:
    **guard the helper, miss the call site.** **THE EXPOSURE IS BOUNDED TODAY, AND THAT IS MEASURED, NOT
    ASSUMED.** A fake interpreter that forks a `sleep 30` grandchild inheriting stdout and then overflows
    a 1 KiB cap returns from `capsule.Run` at **3.005 s** against a **3 s** `ExecTimeout`, carrying the
    CORRECT `*OutputLimitError` — `overran=false`. Same-call control without the grandchild: **11.29 ms**
    (263× faster), so the fixture is sound and the overflow path is prompt when nothing else holds the
    pipe. So the ctx group-kill DOES reach the grandchild and the `ExecTimeout` bound holds; what the
    narrow overflow kill costs is **the promptness the overflow path exists to deliver** — 3 s instead of
    11 ms — not termination. **UNATTRIBUTED RESIDUE, RECORDED AS SUCH:** the surrounding test function
    took **33.38 s** against 3.005 s inside `Run`, coincident with the grandchild's 30 s lifetime;
    something in that path waits and I did NOT establish what. Mechanism unestablished — do not inherit
    this as a diagnosis (rule 3d). **THIS ROW ALSO OWNS the `cmd.WaitDelay` / bounded-cleanup residual**
    declared by row 21 on landing and again by row 20's round-2 quorum (`gpt5-6-sol`: a `CloseOutput() error`,
    a `Wait(context.Context) error`, an explicit bounded cleanup context, and `errors.Join` surfacing of
    kill/close/wait failures). Row 21 is LANDED and a completed row cannot own follow-on work — that is
    queue row 23's defect verbatim, and it is why this row exists rather than a citation. · ~0.5–1 d ·
    NEEDS A DESIGN DOC · gated on nothing · queued iter-99 (controller-measured).


25. **w-capsule-blocked-child-kill-coverage** · clause-2 · **RETIRING THE FLAKY ORACLE ALSO RETIRED
    THE ONLY ARM THAT EVER DROVE A LIVE CHILD BLOCKED IN `write()`, AND NOBODY SEPARATED THE TWO
    JOBS THE FIXTURE WAS DOING.** Filed by iteration 100 on the evaluator's finding
    (`sonnet` 93/100, non-blocking), **reproduced first-party by the controller before acceptance**
    — a NON-BLOCKING label is the judge's opinion of severity, not a measurement. Measured at
    `912009d`: this rig's pipe capacity is **65536 bytes** (a Go probe writing to an undrained
    `os.Pipe` reports `BLOCKED after 65536 bytes`); the surviving
    `TestF6OutputCapKillsChildBeyondOnePipeBuffer` emits `repeat(32)` = 32x16 + newline = **513
    bytes** against a **64-byte** cap; and `readCapped` reads only `limit+1` bytes before stopping,
    so filling the pipe needs total output to exceed its remaining capacity, which 513 bytes cannot.
    **Therefore no arm in the suite — fake or real-interpreter — exercises a live child blocked in
    `write()` on a full undrained pipe being unblocked by `cmd.Process.Kill()`**, which is exactly
    the scenario `host/capsule/capsule.go`'s own "F6 must not decay into F5" comment describes.
    Row 20's 64 KiB fixture was the only thing that ever drove it, and it left with the wall-clock
    oracle it was entangled with. **NOT row 24's**: row 24 owns the cleanup boundary and the
    non-group-wide overflow kill — a different mechanism, and conflating them is the
    named-owner-must-be-precise defect that arm A's obligation (i) exists to prevent. **The shape,
    worth more than the item: when you retire an instrument, enumerate every property it was
    supplying, not just the one that made you retire it.** A sprint here must restore the blocked-
    child arm WITHOUT restoring a wall-clock oracle — the observable is the kill, witnessed by the
    child's own state, never by elapsed time. · ~0.5 d · NO new design doc needed (this row plus
    row 20's corrected §7 residual is the spec).

26. **w-bounded-z3-report-producer** · clause-2 · **SHED THE BOUNDED Z3 REPORT PRODUCER INTO AN ORDERED DESIGN AFTER THE STORE'S DEADLINE-FREE RESIDUE HAS AN OWNER.** Shed from queue item 17, `w-validated-proven-evidence-boundary`, by attended `D-WORLD-24` arm A (Mark Edmondson, issue #68, 2026-08-20T16:04:52Z; ratified 2026-08-20). Scope is the former §3.4 in full: `host/evidence/proof_producer.go`; pinned executable-byte/version checks; bounded process-group execution; independent stdout/stderr limit+1 caps; strict `ai-check` JSON interpretation; required verified identities; MAC-tagged canonical envelope production; source-object reading; and bounded envelope storage. It inherits AC5/M15 and AC20/M28 without renumbering and must design `NewProducer`/`Producer.GenerateProof` coherently rather than assuming their removed tranche-1 spellings.

    **INHERITED OPEN DEFECT — `gpt5-6-sol`, VERBATIM:** *"`Producer.GenerateProof` must read `sourceRef` through `ObjectReader`, but `NewProducer(key, checker, writer, execTimeout, maxOutputBytes)` receives no reader. Even if a reader were added implicitly, both the source read and `WriteObject` use the caller context without requiring or deriving a deadline, so `context.Background()` can block indefinitely on the store's single-connection pool. AC20 only proves propagation of an already-deadlined context, not a bounded operation."*

    **INHERITED OPEN DEFECT — `gemini-3-1-pro`, VERBATIM:** *"Round 11b's additive `WriteObject(ctx, o)` introduces an unpinned lock-contended wait. Because `NewProducer` configures no store timeout, it cannot pin the ordering against the store's busy window at construction. If the caller's context to `GenerateProof` has a deadline shorter than `busy_timeout` (2000 ms), SQLite lock contention will ignore it and overrun the deadline. AC20 vacuously masks this by testing only connection-pool waiting with a decoy, reproducing the exact fixture defect rejected in Round 10 for the read side."*

    These defects are unresolved, not answered by the shed. The design must state explicitly that AC20's decoy exercises the **connection-pool wait, not the lock wait**; the decoy arm cannot stand as evidence about lock contention. It must also close the reader-parameter gap: `GenerateProof` reads `sourceRef` through a reader that the former `NewProducer` never received.

    **ORDERING:** this item cannot be designed until OPEN queue row 23, `w-store-deadline-free-residue-owner`, owns and closes item 18's declared `deadlineFreeReadPins` follow-on in its own terms: thread a real caller context through the 11 production store reads, shrink the pin map in the same diff, then flip it to an empty map so a future deadline-free read reds against an empty expectation. The producer's source read and envelope write must be designed only after that deadline-free residue has an owner and the resulting store-wait contract is known; this row may not silently become a twelfth deadline-free site.

    · **1.0 d carried from item 17's measured cut** (0.60 d producer + ~0.10 d writer/codec share + ~0.10 d AC5/M15 share + 0.20 d bounded-write row); round-12 fixes remain **unpriced, pending design**.

27. **w-interface-hash-does-not-cover-the-interface** · clause-1 · **THE FIELD NAMED `interfaceHash` DOES NOT CHANGE WHEN THE INTERFACE CHANGES, AND THIS MISSION HAS NOW MEASURED THAT TWICE, TWO MONTHS APART, WITH THE SAME HASH LITERAL.** Filed 2026-08-21 (iter-103) off `PE.A`, which is the first change in this repo's history to add a CONSTRUCTOR to an already-exported ADT — the one stimulus that makes the gap visible rather than theoretical.

    **MEASURED AT `cbd17de`.** `InterfaceHash` (`host/pkgproj/pkgproj.go:87`) hashes exactly five things, all read from `ailang.toml`: `Package.Name`, `Package.Edition`, `Package.AILANG`, the sorted `Exports.Modules` list, and the sorted `Effects.Max` list. **It never opens a `.ail` file** — that is `ContentHash`, a different function walking the source tree. So `interfaceHash` is invariant under EVERY source-level change that does not add or remove a module or an effect. `PE.A` added `| ProofReceipt(HashRef)` to the exported `Evidence` ADT; the compiler's own `iface.json` records `ProofReceipt` in its `constructors` table; `contentHash`, `tarballSHA256` and `tarballBytes` all moved (7856 → 7883 bytes); `interfaceHash` stayed at `sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083`.

    **THIS IS A RECURRENCE, NOT A DISCOVERY, AND THAT IS THE POINT OF THE ROW.** The judge (`sonnet`) found the prior instance and the controller reproduced it against the charter before recording: iteration 81 / item 13 recorded the identical mechanism, the identical file, and the identical hash transition — *"`host/pkgproj/pkgproj.go:86` hashes only name/edition/ailang/export-MODULE-NAMES/effects, never opening a source file … (`d16cc882 → d16cc882`)"* — and item 13 responded correctly by AMENDING its own AC9 to *"three fields move, `interfaceHash` byte-identical"*. Item 17's §8.3 then opens *"The `world/core` interface hash changes because a public ADT changes"*, and that false sentence passed **thirteen quorum rounds and two independent reviewers**. A mission-history fact that is recorded, correct, and load-bearing was re-broken by the next document to touch the area, because nothing makes a design-doc author grep prior findings for the surface they are writing about.

    **WHY IT IS A REAL DEFECT AND NOT A NAMING QUIBBLE.** The ready packet is the artifact's published IDENTITY, and `docs/SELF_MOD_PUBLISH.md` puts all three digests in front of an attended human as the thing they are approving. A consumer — or an operator — reading `interfaceHash` to answer *"did the interface change?"* gets **no** for a genuinely breaking ADT change; §8.3's very next sentence concedes the change IS breaking (*"Consumers matching `Evidence` must add the new arm"*). The packet as a whole is not blind — `contentHash` moved — so the exposure is scoped precisely: it is the one field whose NAME promises interface coverage and whose IMPLEMENTATION is manifest coverage.

    **THE ITEM.** Decide and implement one of: **(a)** make `InterfaceHash` cover the exported interface (hash the compiler's own `iface.json` export/constructor/type tables, which already record what is needed); **(b)** rename the field to what it measures (`manifestHash`) and say so in `docs/SELF_MOD_PUBLISH.md`, so the operator is not told a hash means something it does not; or **(c)** keep both, adding a genuine interface digest beside the manifest one. Whichever arm, ship a test that RE-BREAKS on the exact stimulus this row was filed from — add a constructor to an exported ADT and require the interface-covering hash to move — because the existing suite passed this change without noticing, twice. Fix §8.3 as part of the same work; it is deliberately NOT hand-edited now, because item 17's document has just cleared quorum and an unreviewed post-quorum edit is the change nobody would review (§10.12).

    **AN A/B FOR MARK IS LIKELY BUT IS NOT ASKED YET** — (a) vs (b) is a compatibility call with real downstream weight, and this row does not manufacture the decision before a design has priced the arms. · **~0.5–1.0 d**, NEW-DOC. Ordered AFTER item 17's tranche 1 completes; nothing in `PE.B`–`PE.F` depends on it.

28. **w-readobject-refusal-branch-coverage** · clause-2 · **SIX OF `ReadObject`'s EIGHT REFUSAL BRANCHES HAVE NO TEST THAT REDS WHEN THEY ARE NEUTERED — AND THE SIBLING METHOD'S ANALOGOUS BRANCHES DO.** Filed 2026-08-21 (iter-104) by the `sonnet` evaluator of `PE.B`, whose refusal-branch enumeration is the rule-3j drill this repo already runs on validators.

    **MEASURED.** Neutering each branch with `if false && <cond>` and running the whole `host/store` package: the absent-row (`sql.ErrNoRows`) and the over-size guard both DIE; **connection-reserve failure, `BeginTx` failure, the generic non-`ErrNoRows` probe error, `hashref.Parse` failure, the payload-query error, and `tx.Commit()` failure all SURVIVE at rc=0**. The judge supplied the control that makes this a finding rather than repo-wide slack: the *same* neutering applied to sibling `GetObject`'s generic-error and `hashref.Parse` branches **does** red, via `TestReadGettersHonorContext` and `TestReadRetriesUnderTransientExclusiveLock`. So this is a coverage regression measured against the method `ReadObject` was modelled on, not a house standard.

    **WHY IT IS NOT MERELY TIDINESS.** `ReadObject` is the seam `PE.D`'s validator reads untrusted objects through, and its whole contract is *refuse before materializing*. A refusal path nothing reds on is precisely the shape §10.x keeps closing elsewhere: the guard exists, the suite is green, and nobody has established that removing it changes anything.

    **THE ITEM.** One neutering mutation per surviving branch, each with the assertion that kills it; where a branch is genuinely unreachable (the judge argues `busyTimeoutFromDSN`'s `-1` is, because `Open` errors on the same malformed DSN first), DECLARE it in the code and in the AC rather than leaving it as assumed coverage. · **~0.3–0.5 d**. Ordered after item 17's tranche 1; a natural fold into `PE.F`'s full drill if that milestone's scope is ever revisited, but NOT hand-edited into a plan mid-sprint.

29. **w-readobject-absent-is-not-a-distinct-return** · clause-2 · **`ReadObject` SIGNALS "NO SUCH OBJECT" WITH THE SAME TRIPLE IT WOULD RETURN FOR A ZERO-LENGTH PAYLOAD, AND THE ONLY THING KEEPING THOSE APART IS AN INVARIANT THE *WRITERS* ENFORCE.** Filed 2026-08-21 (iter-104) by the `sonnet` evaluator of `PE.B`.

    **MEASURED.** `GetObject` distinguishes absence with a `bool`. `ReadObject(ctx, ref, maxBytes) (ObjectMeta, []byte, error)` has no such channel: absent returns `(ObjectMeta{}, nil, nil)`, and an object present with an empty payload returns `(ObjectMeta{InterfaceHash: <real>, PayloadLength: 0}, nil, nil)` — **`payload` is `nil` in both**, so the entire discriminator is whether `InterfaceHash` is the zero value. The judge verified that all production write paths (`PutObject`'s `validateRef`, `journalObject`'s `SumSHA256`) make a zero `InterfaceHash` unconstructible today, so the ambiguity **cannot presently occur**.

    **WHY IT IS STILL A ROW.** The safety is enforced at the far end of the system from the reader, is undocumented on `ReadObject` itself, and is exactly the kind of cross-component invariant a later writer path can break with no test noticing — the *guard the helper, miss the call site* shape this mission names. The plan said only "the same absent branch `GetObject` has", which a signature carrying no `bool` cannot literally honour; that under-specification is the root, not the executor's reading of it.

    **THE ITEM.** Either give `ReadObject` an explicit absence channel (a `bool`, or a typed `*ObjectNotFoundError` matching the `*ObjectTooLargeError` precedent it already sets), or document the zero-`InterfaceHash` invariant ON the function and pin it with a test that reds if any write path ever admits a zero interface hash. · **~0.2–0.4 d**. Ordered with row 28.

30. **w-evidence-trailing-refusal-branch-is-kind-only-observable** · clause-2 · **`json.go`'s "trailing JSON value" / "trailing data: %v" SUB-BRANCH SPLIT IS UNOBSERVABLE, SO ONE HALF OF IT CAN BE NEUTERED WITH THE SUITE GREEN — THE SAME SHAPE ITERATION 105 JUST FIXED IN THE AC19 ARM, ONE BRANCH OVER.** Filed by the `sonnet` judge of PE.C as non-blocking finding 3, and filed as a row rather than fixed inline because it is a *sibling* of the defect that iteration's commit already repaired, not part of PE.C's own acceptance criteria — absorbing it would have made the milestone's diff about a defect its plan never named.

    **MEASURED (judge, `host/evidence` at `bd48f68`).** Mutating the inner `if err == nil {` at `json.go:58` to `if false && err == nil {` SURVIVES the package suite. Both sub-branches return `RefusalMalformed`, and every trailing-input arm — `TestProofReportStrictRefusals/trailing`, `TestEnvelopeStrictRefusals/trailing`, `TestProposalStrictRefusals/trailing` — asserts only the refusal `Kind`, never the `Detail` text, so which of the two messages was produced is never observed.

    **WHY IT IS A ROW AND NOT A NIT.** The outward refusal identity is unaffected, so there is no security or correctness consequence today — the judge's own severity read, and it holds. What makes it worth a row is the *generalisation*: iteration 105 established, by running both arms, that a `Kind`-only assertion cannot distinguish the guard that fired from a bystander guard producing the same `Kind`, and that dead code added to narrow the observable hides the defect instead of fixing it. That lesson was applied to the AC19 arm and **not** generalised to its siblings in the same file. The row exists so the generalisation is done deliberately, with the enumeration anchored to the DIFF rather than to the one branch someone happened to name.

    **THE ITEM.** Enumerate every `return` on an error path in non-test `host/evidence` — not a `%w`-qualified grep, which is blind to non-wrapping refusals — neuter each with `if false && <cond>` (never delete: an unused import turns "the mutant does not build" into a false kill), and for each branch with no red either add an arm whose observable is unique to that branch or declare the branch unreachable in the code. The iteration-105 judge's sweep reports 17 of 35 such statements currently redding; that count is the judge's, in the scope `grep -rn '^\s*return.*refusal(' host/evidence --include='*.go' | grep -v _test.go`, and is to be re-derived rather than inherited.

31. **w-stale-planned-doc-citations** · clause-1 · **FOUR PRODUCTION FILES CITE A DESIGN DOC BY A
    PATH THAT NO LONGER EXISTS, AND THE BOOKKEEPING STEP THAT MOVES A COMPLETED DOC IS WHAT CREATES
    THEM — SO THE CLASS GROWS BY ONE EVERY TIME AN ITEM COMPLETES.** Measured 2026-08-22 (iter-109)
    while performing item 17's own `planned/` → `implemented/` move, which is instance N+1 of the
    mechanism. `w-worldd-m2` is in `implemented/` and **2** files under `host/`+`scripts/` still cite
    `design_docs/planned/w-worldd-m2.md`; `w-m1-ailang-hardening` **1**; `w-world-library-m1` **1**.
    Positive control that the matcher is not simply matching everything: `w-self-mod-vertical` is
    CORRECTLY still in `planned/` and its **4** citations are valid. Negative control on an invented
    doc name → **0**. No gate depends on the location (`scripts/`, `.github/` and `host/` were swept:
    every hit is a `//` comment), so this is documentation drift, not a build break — but a comment
    that names a path a reader cannot open is a claim, not a citation, and this repo's own headline
    discipline is that the difference matters. **BATCHED WITH A SECOND INSTANCE OF THE SAME CLASS,
    found by the iter-109 designer and reproduced first-party:** `world/types.ail:40` reads "Canonical
    grading for the ratified **five-constructor** representation" directly above an `Evidence` type
    with **six** constructors (`:23-:29`, `ProofReceipt` added later); same-scope control —
    `EvidenceGrade` really does have exactly four (`:34-:38`). **SCOPE:** correct the five stale
    citations, and decide whether the fix is durable — a CI check that every `design_docs/<dir>/…md`
    path named in a comment resolves on disk would make the class impossible rather than
    periodically swept, and it is enumerator-shaped, so it needs an addition-direction mutant (a NEW
    comment citing a bad path must red it), not only a removal one. Note the `.ail` half means the
    AILANG gate runs, so price it as a two-gate change. · ~0.5 d · NEEDS A DESIGN DOC · gated on
    nothing · queued iter-109 (controller-measured during item 17's bookkeeping).


32. **w-wallclock-ceilings-not-derived** · clause-2 · **A DOCS-ONLY COMMIT REDDENED THE GO GATE, AND
    THE FAILING ASSERTION IS A LAPTOP CONSTANT — THE REMEDY FOR THIS EXACT CLASS ALREADY EXISTS IN
    THIS REPO AND WAS APPLIED TO ONE PACKAGE.** Found 2026-08-22 (iter-109) by the controller's
    own pre-commit gate run on a change containing **zero non-`design_docs` bytes**
    (`git diff --stat origin/dev -- ':!design_docs'` → empty, so the red is non-attributable to the
    diff **by construction**, not by argument). `host/capsule`'s
    `TestF5WallClockTimeoutHasElapsedBound` failed *"timeout returned after 2.909405834s, want
    <= 2s"* (`capsule_test.go:230`). **Isolated it is 10/10 PASS**; it fails only inside
    `verify_go.sh`, where `host/broker` (98.6 s) and `host/verifygate` (59.1 s) run concurrently and
    a race leg follows. So the variable is machine load, and the assertion is
    `if elapsed > 2*time.Second` — an ABSOLUTE ceiling on a path whose stimulus (an archived
    interpreter subprocess running `fib(28)`, killed by an injected `ExecTimeout` of **40 ms**) is
    dominated by spawn-and-teardown cost that scales with the machine. The ceiling encodes one
    laptop, exactly as iteration 107's `PE.E` bounds did. **THIS IS THE SHARED SKILL'S RULE 3m,
    WHICH THIS MISSION ITSELF PROPOSED AND WHICH V1 ADOPTED THE SAME DAY** — and the tell is that
    the remedy is already committed here: `host/evidence/realstore_test.go` derives its bound from
    the measured stimulus (`readTimeout := hold / 20`) and keeps an absolute `minDecoyHold` floor
    that fails LOUDLY as an instrument failure. **That fix was applied to the package that hurt and
    never swept** — *guard the helper, miss the call site*, this loop's own named recurring shape.
    **THE EXPOSURE, MEASURED:** `4` `_test.go` files under `host/` carry a hardcoded absolute
    wall-clock value used as a ceiling — `host/archive/archive_test.go`,
    `host/broker/handlers_test.go`, `host/capsule/capsule_test.go`,
    `host/store/read_object_test.go`. Same-scope positive control: **20** test files mention
    `time.Second` at all. Negative control on a fresh invented literal: **0**. And rule 3m's own
    named blind axis is genuinely unturned here — exactly **1** test file anywhere under `host/`
    varies `GOMAXPROCS`. **WHY IT OUTRANKS ITS SIZE:** this class reds `dev` on commits that cannot
    possibly have caused it, including **record commits**, which is how a loop with no human present
    spends an iteration diagnosing its own bookkeeping. **SCOPE:** derive each of the four bounds
    from a stimulus measured in-test so the ratio the design intends holds by construction on any
    machine; keep an absolute FLOOR on the stimulus, loud, so a degenerate stimulus reports
    instrument failure rather than passing quietly; and prove no kill was lost — every mutation the
    four arms own must still die after the change, and each arm must survive a `GOMAXPROCS=1` run
    under contention, which is the axis none of them has ever been measured on. · ~0.5–1 d ·
    NEEDS A DESIGN DOC · gated on nothing · queued iter-109 (controller-measured, first-party,
    during its own pre-commit gate).

33. **w-run-selector-acs-are-vacuous** · clause-2 · **`go test -run` EXITS 0 ON AN EMPTY MATCH SET, SO
    EVERY ACCEPTANCE CRITERION IN THIS REPO SHAPED *"`go test -run 'TestA|TestB'` PASSES"* IS GREEN
    BEFORE EITHER TEST IS WRITTEN — AND STAYS GREEN IF A RENAME SILENTLY ORPHANS THE SELECTOR.**
    Found 2026-08-22 (iter-110): the sprint-planner flagged two of item 14's own criteria as green at
    base, and the controller reproduced both first-party and then ran the control that generalises
    them. Measured: `go test ./host/boundary -run 'Test(BareNetHTTPExemptionIsPerGroup|WorkbenchPackageRemainsTransportFree)' -count=1`
    → **rc=0** with **1** top-level `=== RUN` line (the second test does not exist yet), and
    `go test ./host/daemon -run 'TestWorkbenchReadDeadline|TestReadCtxCancelledAfterHandler' -count=1 -v`
    → **rc=0** with **1** `=== RUN` line. **The negative control is what makes it a class rather than
    two rows:** `go test ./host/boundary -run 'TestZzNoSuchTestIter110' -count=1 -v` is **rc=0 with 0
    `=== RUN` lines** — matching nothing is indistinguishable, in the exit code, from matching
    everything and passing. This is the vacuous-pass class this mission keeps closing (the silent z3
    skip `V27`, the silent `t.Skip`, the sandbox false-green) aimed at the *acceptance criteria
    themselves*, and it is worse than those because a `-run` AC is written precisely for tests that
    **do not exist yet**, so it is vacuous during the entire window it is supposed to govern. **The
    fix is not the ACs one at a time**: it is an enumeration floor — assert the count of top-level
    `=== RUN Test` lines equals the number of names in the selector, the same shape as the
    `EXACT_EVIDENCE_TESTS` count-pin that `PE.F` already proved independently load-bearing at
    iter-108 (`scripts/verify_go.sh`). Scope for the doc: census every `-run` AC across
    `design_docs/planned/` and `design_docs/implemented/` **with a firing same-scope control**, decide
    whether the floor belongs in the ACs, in a shared helper, or in `verify_go.sh`, and say which
    landed criteria were vacuous when they were signed off — a criterion that could not fail is not
    evidence the milestone it cleared was sound. Note item 14's plan already repairs its own three
    instances, so this row is the SWEEP, not those. · ~0.5–1 d · NEEDS A DESIGN DOC · gated on
    nothing · queued iter-110 (controller-reproduced, first-party, with a firing negative control).

34. **w-workbench-unpinned-render-hunks** · clause-5 · **SIXTH AND SEVENTH HUNKS ADDED 2026-08-25 (iter-123), found by the drill's own evaluator for the THIRD consecutive drill milestone, and they are the CLOSED GRAMMAR's PAIR-COMPOSITION half — the same function as the fifth hunk, different guards.** `supportedWorkbenchQuery` (`host/daemon/workbench.go`): **(6)** `:72` `if query["from"] != nil && query["entry"] != nil {` -> `||`, and **(7)** `:75` `return query["object"] != nil && query["payload"] != nil` -> `||`. Each LANDED (old literal 1->0, new 0->1, exact line-content assertion at the named line, `gofmt` rc=0 / 0 bytes), `go build ./...` **rc=0**, full classification arm (`./host/workbench ./host/daemon ./host/boundary -count=1 -v`) **rc=0 with an EMPTY red set** at 148 RUN lines, restored byte-identical. **BOTH SURVIVE.** The guards are LIVE, measured on the function's own truth table with a temporary in-package probe (deleted afterwards): pristine `world`+`from` = `false` and `world`+`payload` = `false`; under (6) `world`+`from` flips to `true`, under (7) `world`+`payload` flips to `true` — and `:127`'s `if !supportedWorkbenchQuery(query)` is what turns the `false` into the `400` §4's closed grammar requires. So each mutant opens a 2-key combination the grammar forbids and nothing reds. **MECHANISM, not a count**: both are two-clause `&&` sub-expressions inside a helper the tests only observe through its single boolean result; `TestWorkbenchRefusalBranches/unsupported-combination` supplies a **one-key** query, which exits at the `len(query) == 1` branch and never reaches `:72`/`:75`, and every other test supplies a combination the grammar ACCEPTS — so the entire `len == 2` FALSE path is unexercised. Distinct from the fifth hunk (that is `len(query) != 2`, i.e. how many keys; these are which pairs). Out of `WB.K`'s scope (it discharges M24-M28), recorded here rather than absorbed, per Standing rule 1, and NOT repaired (§3 rule 5). **The evaluator's recommendation travels with the row**: close the gap with an explicit negative test for a rejected two-key combination (`world`+`from`, `world`+`payload`). Prior head follows.** ~~[**w-workbench-unpinned-render-hunks** · clause-5 · **FIFTH HUNK ADDED 2026-08-25 (iter-121), again by the drill's own evaluator, and this one is the CLOSED GRAMMAR's cardinality half: `supportedWorkbenchQuery`'s gate `if len(query) != 2 {` (`host/daemon/workbench.go:69`) mutated to `if false && len(query) != 2 {` — LANDED (new literal 1, original 0), `go build ./...` **rc=0**, full classification arm (`./host/workbench ./host/daemon ./host/boundary`) **rc=0 with an EMPTY red set**, restored byte-identical, pristine control rc=0. **SURVIVES.** The guard is live, not dead code: neutered, a three-key query such as `?world=…&object=…&payload=1` falls past the `from`/`entry` pair test to `return query["object"] != nil && query["payload"] != nil`, which is **true**, so the request renders `200` where §4's closed grammar requires `400`. The evaluator confirmed that behaviourally with a temporary probe (pristine `400 unsupported workbench parameter combination`, mutant `200`) and removed it before finishing; the controller reproduced the drill arm first-party before adoption. **No test anywhere supplies a three-key query**, so M31/M32 cannot reach it — they pin the unknown-key and duplicate-key branches, a different guard. So the grammar's **vocabulary** half is pinned twice and its **cardinality** half is pinned by nothing. Out of `WB.I`'s scope (M1–M9), recorded here rather than absorbed, per Standing rule 1. Recorded in the same breath, so it is not re-investigated: M31/M32's catalogue text still places their guards in a `parseWorkbenchQuery` helper that no longer exists — the checks are inlined into `handleWorkbench`. Prior head follows.** **FOURTH HUNK ADDED 2026-08-24 (iter-119), and this one was found by the drill's own evaluator rather than by a spot-check: `workbenchHref`'s malformed-query guard `if query != "" && query[0] == '?'` mutated to `if query != "" {` — mutant LANDED (original guard occurrences 0), `go build ./...` **rc=0**, full classification arm (`./host/workbench ./host/daemon ./host/boundary`) **rc=0 with an EMPTY red set**, restored byte-identical. SURVIVES because NO test in the repository supplies a non-empty, non-`?`-prefixed value — every `Href`/`PrevHref`/`NextHref` in every test file already starts with `?`, so the guard's false branch is unreachable from the suite. Reproduced first-party by the controller before adoption. Out of `WB.H`'s scope (it claimed M14–M21 only), so it is recorded here rather than absorbed, per Standing rule 1. Non-gap recorded in the same breath so it is not re-investigated: `{{if .PayloadShown}}` → `{{if true}}` IS killed, but only by `host/daemon`'s `TestWorkbenchPayloadPreviewBound/default-off` — protection living in a different package from the file it guards, invisible to a `host/workbench`-scoped reading.** **THREE HUNKS `WB.B` SHIPPED ARE PINNED BY NOTHING, AND THE DRILL THAT EXISTS TO
    CATCH THAT CANNOT REACH THEM — THE ENUMERATION THAT WOULD HAVE CAUGHT IT IS THE DESIGN DOC'S
    OWN §6 MUTATION TABLE, FROZEN BEFORE THE TEMPLATE EXISTED.** Found by the `sonnet` evaluator at
    iter-112 inside a **91/100 zero-blocking PASS**, and the most significant of the three was
    reproduced first-party by the controller. **(a)** `{{if .World.Available}}` — the primary
    "world selected" header path, not an edge case — mutated to `{{if false}}`: mutant LANDED (0
    remaining actions), `go build ./host/workbench/...` **rc=0**, whole package **rc=0 with zero
    `--- FAIL`**; no test sets `World.Available` (`grep -c 'World:' render_test.go` = **0**,
    same-scope control `Object:` = **4**). **(b)** the `{{if .PayloadTruncated}}<p>truncated</p>`
    marker, same signature (judge-measured). **(c)** the PASS-verdict span's
    `aria-label="test verdict PASS"`, same signature — `Render` is called four times in the test
    file and never with a passing verdict, because `t.Run("pass", …)` only checks struct fields.
    **Why WB.H–WB.K cannot close this:** those milestones discharge `M1`–`M32` and nothing else,
    and none of the three appears in the table — `grep -icE 'world.*available|truncat|verdict-pass'`
    over the 32 rows = **0** with the same-scope control `verdict` = **2**. So the sprint can
    complete its entire non-vacuity spine with these three live and unprotected. This is the shared
    skill's **rule 3n** (*your mutation set is derived from what the milestone FIXES, so it misses
    what the milestone SHIPS — anchor the enumeration to the DIFF, which is complete by
    construction*) arriving one level up: the stale enumeration is not a controller's mutation list
    but a **quorum-reviewed design doc's table**, written before the code it is supposed to cover.
    Scope for the doc: add the three arms (and ask, per rule 3n(a), what an ADDED member would
    reveal — a removal proves a check FIRES, only an addition proves it LOOKS); then decide the
    durable fix, which is probably a per-milestone *diff-anchored* hunk sweep in the sprint protocol
    rather than three more table rows, since the next milestone will ship hunks this table does not
    know about either. **Do not fold this into item 14** — rule 3n(b) makes an unpinned hunk a queue
    row, not a quiet widening, and standing rule 1 is one item. · ~0.5 d · NEEDS A DESIGN DOC ·
    gated on item 14's remaining milestones (the surface is still growing) · queued iter-112
    (judge-found, controller-reproduced, with firing same-scope controls).

35. **w-workbench-view-model-fields-render-nowhere** · clause-5 · **THREE `EntryView` FIELDS THE
    DAEMON POPULATES ON EVERY TIMELINE ROW ARE RENDERED BY NO TEMPLATE ACTION, AND ONE OF THEM IS
    THE ONLY ROUTE BY WHICH A COMMITTED OBJECT'S HASH CAN REACH THE PAGE AT ALL.** Filed 2026-08-23
    (iter-113) off `WB.C`, whose plan task 13 is **unsatisfiable because of this** — it asks the
    render test to assert two distinct committed OBJECT hashes appear in the body, and no such hash
    can appear. Measured inside `{{range .Timeline.Entries}}` in `host/workbench/render.go`:
    `TransitionFn`, `Interpreter` and `TransitionRef` render **0** times each, against a same-scope
    known-positive control of the five sibling fields (`EntryIndex`, `EntryHash`, `PrevEntryHash`,
    `SemanticsEpoch`, `WrittenBy`) that each render. All three are `EdgeView`s and all three are
    populated by `host/daemon/workbench.go`'s `EntryView` construction, complete with `Href` values
    — the view model and the handler each do their half and the template drops it on the floor.
    **This is a SIBLING of row 34 and a DIFFERENT defect**: row 34 is shipped template hunks that no
    test pins; this is populated view-model fields that no template *reads*, so a mutation drill over
    the template can never surface it and neither can one over the handler. **It is also what forced
    this iteration's disclosed executor deviation** — the executor, correctly refusing to widen into
    the landed WB.B renderer, manufactured a `page.Notice` carrier to satisfy the assertion; the
    controller removed it and filed this row instead. **Do not fold this into item 14** — rule 3n(b)
    makes an unrendered field a queue row, not a quiet widening, and standing rule 1 is one item.
    Open question the design must answer rather than assume: whether the fix is to render the three
    fields (they carry `Href`s, so the provenance-walk mechanism already exists) or to stop
    populating them, which is a design question about what the timeline is *for*. · ~0.5 d · NEEDS A
    DESIGN DOC · gated on item 14's remaining milestones · queued iter-113 (executor-surfaced,
    controller-reproduced with firing same-scope controls).

36. **w-test-oracles-imported-from-the-code-under-test** · clause-2 · **`workbenchCSP` IS NOT AN
    ISOLATED SLIP: **75** PRODUCTION SYMBOLS IN `host/` ARE REFERENCED FROM `_test.go` FILES, AND
    NOTHING IN THIS REPO HAS EVER TRIAGED WHICH OF THEM SIT IN AN *EXPECTATION* POSITION.** Filed
    2026-08-23 (iter-113) off `WB.C`'s blocking finding, and filed as the SWEEP rather than fixed
    inline — item 14's plan repairs its own one instance; the census is its own item, and standing
    rule 1 is one item. **The mechanism, stated precisely, because the naive form of this row would
    be wrong:** referencing a production constant from a test is *not* itself a defect — it is
    correct wherever the constant IS the contract (`internalErrorMessage` is pinned that way on
    purpose, V21). It becomes a **tautological oracle** exactly when the mutation table names that
    test as the killer of a mutant that changes **the constant's own value**: expected and actual
    then move together and the arm can never red. `WB.C` is the confirmed instance —
    `TestWorkbenchSecurityHeaders` imported `workbenchCSP` while being named as `M22`'s killer, and
    `M22` LANDED (occurrences 1 → 0), BUILT rc=0, and left the package **rc=0 with an empty
    `--- FAIL` set**. **What is measured and what is not:** the population is 75 (instrument
    validated — a fabricated symbol returns **0**, so the search discriminates); the count of
    *tautological* ones is **UNKNOWN**, and that unknown is the row. Symbols worth looking at first
    are the ones added by this mission's own recent anti-vacuity work, since a constant extracted
    to make a bound derivable is exactly the kind a test then imports: `decoyHoldRatio`,
    `minDecoyHold`, `watchdogBound` (row 32's fix), `expiredReadDeadline` (item 19's fix),
    `deadlineFreeReadPins`, `frozenPublicSurface`. **This is rule 3i's blind spot, not another
    instance of it** — 3i asks *what else writes this value*, and for an imported oracle the answer
    is *nothing*, so the rule returns clean on the worst case. The tell is syntactic and greppable:
    *a test file names a production constant inside its own expectation table.* · ~0.5 d · NEEDS A
    DESIGN DOC · not gated on item 14 (the census spans the whole repo) · queued iter-113
    (evaluator-found instance, controller-swept with a firing negative control).

37. **w-workbench-store-error-guards-unpinned** · clause-2 · **RE-MEASURED AND SHARPENED 2026-08-25
    (iter-121) BY `WB.I`, THE DRILL THAT WAS SUPPOSED TO DISCHARGE `M9`: THE ROW IS WORSE THAN FILED AND ITS
    ARITHMETIC WAS WRONG. `M9` IS NOT MERELY UNPINNED — IT IS UNKILLABLE AS A ONE-LINE MUTATION, BECAUSE THE
    GUARD IS MASKED BY A SIBLING GUARD ON THE SAME REQUEST PATH.** Three arms, each LANDED by occurrence
    count, BUILT rc=0 before any test result was read, restored by `cp` and verified byte-identical:
    `workbench.go:179` alone -> `…/store-error` **PASSES** (the 500 is produced by `:256`); `:256` alone ->
    `…/store-error` **PASSES** (produced by `:179`); **both together -> FAIL**,
    `workbench_test.go:183: status = 200, want 500`. So the named test asserts a status code that more than
    one guard on the path can produce, and no single-site mutant of either is detectable. **The count of nine
    was a literal count, not a guard count.** Re-measured site by site, the nine `if err != nil {`
    occurrences are **3 parse guards** (`139` the `?from=` parse, `172` = M2, `209` = M4), **1 log guard**
    (`89`, in `writeWorkbenchInternalError`) and **5 genuine store-error guards** (`179`, `190`, `214`,
    `241`, `256`) — so `172`/`209` were counted twice by the original filing, once as killers and once as
    survivors. Of the five real ones, only `179` has any killer at all
    (`TestWorkbenchReadDeadline/blocking-store`), and `190`/`214`/`241` are **unreachable with a failing
    store from anywhere in the repository**: `failingStore` is installed at exactly two sites
    (`workbench_test.go:167`, target `/workbench` with no query; `read_deadline_test.go:756`, whose route
    list `seedReadRoutes` holds six `/v1/…` entries and no `/workbench`), and an all-five-guards arm returns
    a red set **identical** to the `179`+`256` pair's. **The original filing's `179` reading was stale, not
    wrong**, and the delta is dated: its killer entered the tree at `8f0037c` (`WB.F`, #88), after `WB.D`'s
    `e50fbea` (#86) where that sweep ran — a drill result is a claim about a tree, not about a file. `M9` is
    recorded **SURVIVED** in the sprint plan §7g per §3's `survived_is_a_result`; the mutant was not
    repaired, the test was not repaired, the row was not omitted, and the fix stays here rather than being
    absorbed into `WB.I` against Standing rule 1. Whoever takes this row owns four store-error guards with no
    possible killer, the `?from=abc` parse guard, the `parsed < 0` half of the entry guard, and the log guard
    at `89`. **Prior filing follows.** **SEVEN OF THE NINE `if err != nil {`
    STORE-ERROR GUARDS IN `host/daemon/workbench.go` HAVE NO KILLER, AND `M9`'s MUTATION ROW
    CANNOT SAY WHICH ONE IT DISCHARGES.** Filed 2026-08-23 (iter-114) off `WB.D`'s evaluation, and
    filed rather than fixed inline because rule 3n(b) makes an unpinned hunk a queue row and never
    a quiet widening of the milestone — the same call iters 112 and 113 made on rows 34 and 35.
    **Measured**, one guard at a time, each neuter asserted LANDED, asserted to BUILD rc=0 *before
    any test result was read*, and restored by `cp` with sha256 byte-identity: lines
    `89 139 179 190 214 230 245` each leave `go test ./host/daemon -count=1` at **rc=0 with 0
    `--- FAIL`**; only `172` and `209` have killers. Pristine control green. **The consequence is
    a wrong STATUS CODE, not missing coverage:** `failingStore` returns `ok=false` *beside* its
    error, so a neutered guard lets a genuine store fault fall through to `NotFound` — a real
    store failure reported as **404**, which is exactly the malformed-vs-absent-vs-error
    distinction §2.4 exists to preserve and which V20/V21 measured the JSON handlers already
    modelling. **The mechanism, stated precisely, because this is NOT another instance of rule
    3i:** 3i asks *what else writes this value*, and here the observable is fine — the arm that
    exists (`store-error`) genuinely kills the guard it reaches. The defect is one level up, in
    the TABLE: `M9`'s row is `if false && err != nil {` with **no site qualifier**, and the
    pattern now occurs **nine** times in the file. A drill that applies the row as written
    neuters whichever occurrence its tooling matches first and records `M9` DISCHARGED, while
    seven siblings stay unpinned. That is rule 3a(i-e)'s enumeration gap — *a removal proves the
    check FIRES, only an addition proves it LOOKS* — aimed at a quorum-reviewed mutation table
    rather than at a gate, and it is the FOURTH distinct hollow-pin mechanism this one sprint has
    produced (`WB.A` zero-value, `WB.B` observable wider, `WB.C` observable equal, `WB.D` row does
    not identify a site). **WB.I/WB.J MUST PARAMETERISE `M9` OVER ALL NINE CALL SITES**, not over
    the first textual match, or the discharge will be as hollow as `WB.C`'s CSP pin was. Two
    lower-severity siblings found in the same sweep and recorded here rather than lost: a
    non-numeric `from` (`?from=abc`) has no dedicated arm, and its branch reuses the message
    `"from index must be non-negative"`, which is misleading for a parse failure; and `entry`'s
    `parsed < 0` half shares a guard with the tested `err != nil` half, so `?entry=-1` is
    unpinned. · ~0.5 d · gated on item 14 (it is item 14's own file; sequence it with WB.I/WB.J)
    · queued iter-114 (evaluator-found, controller-swept across the whole pattern — the judge
    reported 4 of 5, the sweep says 7 of 9).

38. **w-workbench-timeline-pagination-is-unwired-at-both-ends** · clause-5 · **THE TIMELINE'S
    PAGINATION SEAM IS BROKEN FROM *BOTH* DIRECTIONS AND THE TWO HALVES ARE COMPLEMENTARY:
    `TimelineView.Truncated` IS WRITTEN BY THE HANDLER AND READ BY NOTHING, WHILE
    `TimelineView.NextHref`/`PrevHref` ARE READ BY THE TEMPLATE AND WRITTEN BY NOTHING — SO THE
    WORKBENCH CAPS THE TIMELINE AT 100 ENTRIES AND THEN EMITS NO NEXT-LINK ON ANY PAGE, INCLUDING
    THE ONE THAT IS TRUNCATED.** Filed 2026-08-23 (iter-115) off `WB.E`. Measured, with same-scope
    known-positive controls firing in every call: `.Truncated` has **1** writer
    (`host/daemon/workbench.go`, the line `WB.E` shipped) and **0** non-test readers anywhere under
    `host/`, against a control of **3** for the sibling `PayloadTruncated`; `NextHref`/`PrevHref`
    are referenced by **2** template actions (`host/workbench/render.go:140-141`) and appear **0**
    times anywhere in `host/daemon`, against a control of **3** `page.Timeline` assignments in that
    same file. Both links sit behind `{{if .Timeline.PrevHref}}` / `{{if .Timeline.NextHref}}`, so
    the zero value renders **nothing at all** — there is no broken link to notice, which is why
    four milestones of renderer work never surfaced it. **THE TWO HALVES ARE ONE FIX, NOT TWO:**
    `Truncated` is precisely the signal a next-link needs ("the cap ended this read, there is
    more"), and it was computed and discarded in the same milestone whose deliverable is the cap.
    **THIS IS A SIBLING OF ROW 35 AND A SHARPER MEMBER OF THE SAME CLASS.** Row 35 is populated
    view-model fields no template reads; this row adds the mirror — template actions no handler
    feeds — and shows the class has a *user-visible* consequence rather than only a hygiene one:
    §2.2 refuses `?from=N` alone, so a reader who reaches the 100-entry cap has no in-page route
    onward and must hand-craft `?from=100&entry=100`. **HOW IT WAS FOUND, because the method is the
    reusable part:** the mutation set was anchored to the **diff** (rule 3n) rather than to the
    defect `WB.E` was fixing, so the fourth hunk — the `Truncated` write, which no one would have
    thought to mutate — was drilled and redded nothing. A defect-derived set would have run three
    mutants, all sole killers, and reported a perfect drill. The judge then supplied the mirror
    half. **Do not fold this into item 14** — rule 3n(b) makes an unpinned hunk a queue row, not a
    quiet widening, and standing rule 1 is one item. Open question the design must answer rather
    than assume, and it is the same question row 35 poses one level up: whether the timeline is a
    *paged* surface (wire `Truncated` → `NextHref`, and decide what query state a next-link may
    emit given the closed grammar) or a *head-only* surface (delete all three fields and the two
    template actions). **Sequence with row 35** — both touch the same `TimelineView`/`EntryView`
    seam in `render.go`, so designing them together is likely cheaper than twice. · ~0.5 d · NEEDS
    A DESIGN DOC · gated on item 14's remaining milestones · queued iter-115 (controller-measured
    from a diff-anchored drill; mirror half judge-surfaced and controller-reproduced with firing
    same-scope controls).

39. **w-session-authority** · clause-3 · **THE REPO HAS NO INBOUND CREDENTIAL -> SESSION
    RESOLUTION AT ALL, AND `D-WORLD-26` SETTLED THE ENVELOPE WITHOUT SUPPLYING THE CONTENTS.**
    Surfaced as `gpt5-6-sol`'s sole blocking objection at `w-mcp-projection` quorum **round 4**
    (2026-08-25, iter-125) and **confirmed first-party by the controller** rather than forwarded
    (skill rule 3f — a reviewer's premise objection is the controller's to measure). Measured at
    `2e44e3e` over `host/` (96 Go files, scope asserted with `test -d`), one call, counts are of
    matching LINES: `Bearer` **0**; `Authorization` **1** (a prose comment in `approve.go`, not a
    code path); `Credential` **128** but **all OUTBOUND** — `RegistryCredentialProvider` /
    `FileRegistryCredentialProvider` / `AssertNoAmbientRegistryCredential` in
    `host/broker/credential.go`, a credential World *presents* to an upstream registry;
    `Authenticate` **29** but **all evidence-envelope** — `AuthenticatedEnvelope` + codec in
    `host/evidence/`, i.e. signed evidence, not HTTP session auth; and
    `func .*(GetSession|LookupSession|ResolveSession|SessionByID|FindSession)` over `host/ cmd/`
    → **0**, so **nothing in this repo resolves a session by string**. Known-positive control in
    the SAME scope and same command shape (`'Session' host/`) → **181**, and a fresh absent
    literal → **0**, so the zeros are measurements and not a broken instrument (rule 3a(i-d) —
    the control is scoped to the path under test). `host/broker/broker.go:87` is
    `func NewSession(s *store.Store, episodeID string, grants []Capability, registry Registry)
    *Session` — **the grants are an ARGUMENT**, so nothing decides what grants a credential
    carries. **THE ITEM.** Design the session-authority boundary that `Authorization: Bearer`
    now presupposes: who **mints** a session credential and through what surface; where the
    credential -> (episode, grants) **mapping** lives, given that clause 3 put `broker.Session`
    in-process on purpose; what **expires or revokes** it; and how the resolver **fails closed**
    on absent / malformed / unknown / expired credentials, which is clause 3's inviolable half.
    Note the two constraints `D-WORLD-26` carries and does not relax: a `Bearer` value here is a
    SESSION credential and **never** an API key (the static `serve-api` key was measured
    process-wide at iter-24 and cannot represent a session), and there is no unauthenticated
    degrade path. **WHY IT IS A ROW AND NOT A REVISION** (skill: *a pre-existing defect surfaced
    by a reviewer is a QUEUE ROW, not a revision*): clause 3 landed the session model
    deliberately in-process and no HTTP-facing authority was ever built, so this is a gap
    `w-mcp-projection` **fails to fix**, not one it **introduces** — growing that doc a third
    time to absorb it is exactly what the decomposition rule exists to prevent. **SEQUENCING:**
    this row **gates `P6.B-A2A` and nothing earlier** — `P6.T`, `P6.D` and `P6.V` have drawn
    zero objections across all four rounds and carry no session-authority dependence, so they
    remain executable the moment the second split lands. · ~0.5–0.8 d · NEEDS A DESIGN DOC ·
    gated on nothing · queued iter-125 (reviewer-surfaced, controller-measured with firing
    same-scope controls).


40. **w-a2a-session-projection** · clause-6 · **[BLOCKED on row 39]** · The A2A half of the old
    `P6.B`, carved out of row 5 at **split #2** (2026-08-26, iter-126) into
    [`design_docs/planned/w-a2a-session-projection.md`](../design_docs/planned/w-a2a-session-projection.md)
    (601 lines). Carries quorum round 4's `gpt5-6-sol` objection **VERBATIM** in its opening
    problem statement (controller-checked block-quote-aware: `strongest_objection`/`catch`/
    `proposed_fix` **3/3 True**; one-char-mutation negative control **False**; invented-sentence
    negative control **False**), and retains the parent's acceptance identifiers **AC1–AC9,
    AC11–AC14** plus their named mutations — the numbering gaps on both sides are deliberate so
    five rounds of quorum history stay traceable. **It also carries `P6.D`** (the `v0.33.2` pin +
    the single `allowedDepModules` **package-path** entry + its narrowness test + `AC16` +
    `MUT-ALLOWLIST-ROOT`/`MUT-FACADE-IMPORT`), deferred out of the parent at the round-5 carve-out
    on `gpt5-6-sol`'s verbatim fix: **dependency admission is ATOMIC WITH ITS FIRST REAL CONSUMER
    and is never pre-landed.** Its `host/projection` handler import IS that consumer.
    **ORDERING IS A RACE, not a certainty:** if the MCP dispatch child (blocked on
    [`ailang#885`](https://github.com/sunholo-data/ailang/issues/885)) unblocks first, IT carries
    the admission instead and this row inherits it — check before starting,
    `git grep -n 'serveapi/protocol' -- go.mod host/daemon/daemon_test.go`. **NOT quorum-cleared**;
    it is quorumed at pick time, once unblocked, not now. Start condition is threefold: row 39's
    session-authority design landed, the parent's `P6.T` (**done**) and `P6.V` green, and `P6.D`
    landed here or in the sibling. · ~0.7d + ~0.15d · gated on row 39 · queued iter-126.

41. **[LANDED 2026-08-26 (iter-129) — ITEM COMPLETE.** PR [#97](https://github.com/sunholo-data/ailang-world/pull/97) → squash [`8e3c8cd`](https://github.com/sunholo-data/ailang-world/commit/8e3c8cd). **Gate 3b GREEN on the
    PR head and taken AFTER the Actions incident was marked resolved** (`18:01:30Z`), so it is
    attributable — the symmetric-green rule binds as hard as the red one. SHA-addressed from a
    re-derived 40-char `headRefOid`: `checks=2` **== expected=2**, the expected set enumerated from
    `ci.yml`'s own job list (`ailang-verify`, `go-verify`) rather than counted from what turned up;
    `not_green=0`, both `success`, `runs=1 event=pull_request`, `mergeable` read FIRST
    (`MERGEABLE`/`CLEAN`); every count asserted NUMERIC before comparison. **The dropped webhook
    delivery was NOT replayed when the incident closed** — ninety minutes past `resolved`, the PR
    head still read `checks=0`/`total=0` — so the owed re-run had to be MANUFACTURED, with Gate 3b's
    tree-identical empty commit through the git API (`tree` sha asserted byte-identical `2873892…`
    before the ref moved): `runs=0` → `runs=1`, `event=pull_request`, `jobs=2` within 20 s. Evaluator
    `sonnet` **92/100, ZERO BLOCKING**; 18 mutation arms, ZERO survivors, with M1/M2 — `P6.T`'s
    recorded SURVIVORS — both now RED. The half of that gap with no lever at all (a dropped push to
    `dev`) is **row 47**. Doc + sprint plan moved to `design_docs/implemented/`.
    **Prior head text follows.** ~~**[IN-SPRINT — RESUME POINT, NOT LANDED]**~~ **w-setup-go-pin-unguarded** · **BUILT, EVALUATED
    `sonnet` 92/100 ZERO BLOCKING, ALL GATES GREEN OUTSIDE THE SANDBOX — BLOCKED ONLY ON A
    DECLARED GITHUB ACTIONS OUTAGE.** Design [`74c47d5`](https://github.com/sunholo-data/ailang-world/commit/74c47d5)
    (quorum-cleared, 2 rounds, narrow-refinement carve-out). Sprint
    [`098f608`](https://github.com/sunholo-data/ailang-world/commit/098f608) on branch
    `sprint/w-setup-go-pin-unguarded`, PR [#97](https://github.com/sunholo-data/ailang-world/pull/97),
    `MERGEABLE`/`CLEAN`. **Gate 3b UNDISCHARGED:** `checks=0`, `actions/runs?head_sha=<40>` `total=0`
    after an 8-minute bounded poll, known-positive control rev-parsed (`74c47d5` → `checks=2`,
    `runs=1`) so the zero is TRUE; `githubstatus.com` = **Partial System Outage**, open **Incident
    with Actions** created `2026-08-26T15:11:58Z`, PR created ~`15:05Z` inside the window; last
    repo-wide run created `14:02:25Z` and nothing since; `ci.yml` declares no `workflow_dispatch`.
    **RESUME:** re-poll `#97` once the incident is marked resolved — a green taken DURING an open
    incident is unattributable and does not count. Do NOT arm auto-merge (its checks are not green
    now), do NOT revert. 18 mutation arms, ZERO survivors; M1/M2 — `P6.T`'s recorded survivors —
    both RED. **Prior head text follows.** ~~**[NEXT]**~~ **w-setup-go-pin-unguarded** · clause-2 · **A SURVIVOR THE `P6.T` DRILL RECORDED RATHER THAN
    PAPERED OVER.** `MUT-TOOLCHAIN-REGRESS` arms 3 and 3′ mutate `actions/setup-go@v5`'s
    `go-version: '1.26.6'` back to `'1.25.6'` at `.github/workflows/ci.yml:28` and `:109`. Both
    **LANDED** (occurrence count `4 → 3`, a query against the file's own view) and both
    **SURVIVED**: no command locally or in CI turns a `setup-go` regression red. Confirmed
    independently by the sprint evaluator, which rebuilt the instrument rather than inheriting it —
    `git grep -n "GOTOOLCHAIN\|go-version" host/ cmd/` → **0** hits, i.e. no static text-scan test
    of the `TestZ3PinDeclaredOnceAndInstalledInBothJobs` style exists for toolchain-pin consistency,
    and nothing local can execute `actions/setup-go`. **THE EXPOSURE:** a regression there keeps CI
    **green** and silently pays a second toolchain download per job — the two pins drift apart with
    no symptom. **THE ITEM:** one static consistency test asserting that every `GOTOOLCHAIN` value
    and every `setup-go` `go-version` in `ci.yml` agree, and that the count of each is what the
    job list implies. The precedent already exists in this repo and should be mirrored:
    `TestZ3PinDeclaredOnceAndInstalledInBothJobs`. Also fold in the related pin the N7 sweep
    surfaced: `design_docs/verification/w-race-gate-blindspot/run.sh:25`
    `KNOWN_GOOD="go1.25.6 go1.24.9"` — an **executable** probe list (`test -x` YES) that never
    probes `go1.26.6`, the version the repo now actually pins, so the instrument certifying
    "known-good" has never tested the live toolchain. · ~0.15d · gated on nothing ·
    queued iter-126 (drill-surfaced, evaluator-corroborated).

42. **[LANDED 2026-08-27 (iter-131) — ITEM COMPLETE.** PR [#98](https://github.com/sunholo-data/ailang-world/pull/98) → squash [`58c8f7f`](https://github.com/sunholo-data/ailang-world/commit/58c8f7f). **Gate 3b GREEN on the MERGE commit**, SHA-addressed from a 40-char `rev-parse`d SHA: `present=2` **== expected=2**, where expected was **enumerated from `ci.yml`'s own job list** (`ailang-verify`, `go-verify`) rather than counted from what turned up, with `ls .github/workflows/` confirming `ci.yml` is the only workflow so the enumeration is complete rather than a hand-picked subset; `notdone=0`, `notgreen=0`, both `success`, run existence `runs=1 event=push`, parent control `checks=2`; `mergeable` read FIRST (`MERGEABLE`/`CLEAN`) so the boring cause was excluded before any dropped-event theory. Evaluator `sonnet` **97/100 PASS, ZERO BLOCKING**, in its own worktree, having re-run **all eight** mutation arms itself rather than the four its directive required. **THE WORK WAS AUTHORED BY ITERATION 130, WHICH DIED BEFORE COMMITTING** — the design doc landed at `fc4776f`, the sprint plan and the executor's four-file diff were left UNCOMMITTED in `.wt-world-iter130`, and the charter, the log and the STATUS block carried **zero** trace of any of it. Iteration 131 VERIFIED rather than adopted: acceptance gates **G0–G11 all green**, the two new names counted **2 `=== RUN` / 2 `--- PASS` / 0 `--- FAIL`** rather than by exit code (the naive `-run` form is `rc=0 [no tests to run]`, re-proved on a fresh negative name), and the **full drill re-run: 7 RED arms, zero survivors, plus M2b's boundary GREEN control**. **What LANDS:** `TestReproModuleFloorStaysBelowKnownBadToolchains` binds `repro/go.mod`'s floor at or below the oldest `KNOWN_BAD` toolchain in `run.sh` — compared with `go/version`, never string ordering — with instrument-failure floors on validity, assignment count and non-emptiness; `TestCanaryDeclaresPositiveArmOnly` fences the canary against a re-added known-bad arm (zero `GOTOOLCHAIN` tokens) and requires the `POSITIVE ARM ONLY` marker. **M2b(ii) re-measured armed at equality on ONE source with ONE variable**: `GOTOOLCHAIN=go1.26.0 go run .` → `BUG: Field="" want "stateRoot"`, `go1.26.6` → `OK`, **both rc=0** — which is why the AC15 row now says to read stdout and never the exit code. **A trap fired inside the drill and the assert-landed rule caught it**: the first M5 attempt did not apply (an exact-string delete whose needle was mid-line), and its `rc=0` was indistinguishable from a survivor until the `POSITIVE ARM ONLY`-count assertion read **1** where **0** was required. **Prior head follows.** ~~[42. **w-canary-control-does-not-survive-a-floor-raise** · clause-2 · **A KNOWN-POSITIVE CONTROL
    THAT STOPS BEING RUNNABLE THE MOMENT ITS OWN MILESTONE LANDS — AND THE CONTROLLER NEARLY
    BANKED IT AS IF IT TRANSFERRED.** `AC15` requires the committed canary to print `OK` under the
    pinned toolchain *while its known-bad control prints `BUG` under a deny-listed one, in the same
    run*. That is **structurally unsatisfiable after `P6.T`**: once `go.mod` declares `go 1.26.6`,
    Go's module floor rejects `GOTOOLCHAIN=go1.26.5` **before `TestToolchainCanary` executes**, so
    the red observed is the floor assertion, not `toolchain_canary_test.go:41`. Measured both ways
    by the executor and re-derived by the evaluator, which attacked this claim hardest and found it
    genuine, permanent, and recurring on **every future floor-raise**. The controller HAD measured
    the three-arm control on the **pristine** tree (`go1.26.6` PASS / `go1.26.5` FAIL with
    `Field="" want "stateRoot"` / `go1.25.6` PASS) and that reading is sound — **but it does not
    transfer past the landing**, which is exactly the trap. Compensated today by the standalone
    repro (drill arm 5), which lives in a separate `go 1.22` module and therefore escapes the
    floor: `go1.26.5` **BUG**, `go1.26.6` **OK**, `go1.25.6` **OK**. **THE ITEM:** decide where a
    version-pinned known-bad control belongs so it stays runnable across floor raises — most
    likely by moving the canary's known-bad arm into the nested repro module (which already has
    the property) and having the in-module canary assert only the positive arm, with `AC15`'s
    wording corrected to match. **THE GENERALISATION, which is why this is a row and not a
    footnote:** *a control embedded in the artifact it controls inherits that artifact's
    constraints, so raising a floor can silently disarm the instrument that proves the floor was
    needed.* · ~0.2d · gated on nothing · queued iter-126 (drill-surfaced, evaluator-confirmed).~~
43. **[LANDED 2026-08-27 (iter-132) — ITEM COMPLETE.** PR [#99](https://github.com/sunholo-data/ailang-world/pull/99) → squash [`ecfb62d`](https://github.com/sunholo-data/ailang-world/commit/ecfb62d). **Gate 3b GREEN on the MERGE commit**, SHA-addressed: `present=2` **== expected=2** (enumerated from `ci.yml`'s own job list; `ls .github/workflows/` confirms `ci.yml` is the only workflow, so the enumeration is complete rather than a hand-picked subset), `notdone=0`, both `success`, run existence `runs=1 event=push`, parent control `checks=2` rev-parsed not hand-expanded, `mergeable` read FIRST. **WHAT LANDS:** a marker-bounded comment block at the head of `scripts/verify_ail.sh` naming all six Tier-1 sites (plus a **column-0 Python** pointer above `REQUIRED_VERIFIED`, because `:270` opens a heredoc that `:274` sits inside), a new **§S8** in `design_docs/coding-standards.md` carrying the same map and the regeneration recipe, `host/verifygate/floor_raise_inventory_test.go` binding the two homes, and a positional-literal repair of `docs/SELF_MOD_PUBLISH.md:39` (it cited the package gate `at :224`; the call site is `:403`). **Gate strength untouched**: AC5's non-comment diff against `476069d` is EMPTY and the pinned gate is rc=0 at the **unchanged** floor of 11 identities / 40 named tests. **THE HEADLINE IS THAT THE JUDGE BROKE THE SPRINT ON ITS OWN THESIS, ONE SITE OVER.** Evaluator `sonnet` scored round 1 **62/100 FAIL**, reproduced first-party by the controller: deleting site 3's entire row from the script-block home left the test **PASSING** (mutation LANDED `1→0`, mutant BUILDS, `RUN=1 PASS=1 FAIL=0`), because its needles `REQUIRED_VERIFIED` and `EXACT_TOTAL_VERIFIED` each occur **TWICE** inside the block — once in site 3's row, once in unrelated prose — while the SAME deletion in the §S8 home correctly REDS, an asymmetry that hid the hole from anyone checking one home. Quorum round 1 had hardened **site 1** against exactly this class and the other seven needles were never re-examined. **The fix is structural**: all six sites in BOTH homes now anchor to their exact on-disk row literal, and the test requires each anchor **exactly once**, failing loudly with home/site/literal/count on 0 or ≥2 — so a needle that cannot detect its own row's deletion is itself a test failure. Round 2 **90/100 PASS, zero blocking**, after twelve row-deletion arms (six per home, zero survivors) and five independent attacks on the new assertion, all held. **Prior head follows.** ~~[43. **w-floor-raise-coupling-inventory** · clause-2 · **THE `P6.V` SPRINT PROVED THAT RAISING THE
    VERIFIED-IDENTITY FLOOR TOUCHES SIX FILES, AND THE ONLY PLACE THAT IS WRITTEN DOWN IS A
    COMMIT MESSAGE.** Three roles each found a different subset and no single role found them all:
    the design doc named 2, the sprint planner found 3 more, and the **6th — the pristine
    known-positive control string `✓ 10/10 required world/ identities verified across 11
    module(s)` at `host/verifygate/module_manifest_gate_test.go:128` — surfaced only when the
    controller discharged its mandatory out-of-sandbox gate re-run** and met a genuine `rc=1`
    that the executor had (honestly, and correctly under its own rules) labelled
    `UNINFORMATIVE UNDER SANDBOX`. The coupling is real and permanent: `scripts/verify_ail.sh:403`
    calls `verify_world_package.sh`, which binds `world/*.ail` to a hashed published package, so
    every identity added to `REQUIRED_VERIFIED` also moves `EXACT_TOTAL_VERIFIED` (`:323`), the
    projection copy under `packages/world-core/` (step 3/9 is a SHA-256 byte-identity check), the
    ready-packet golden (step 9/9 `cmp -s`), the `docs/SELF_MOD_PUBLISH.md` digest table (gated by
    `host/runbook`), and that verifygate control. **THE ITEM:** publish the inventory where the next
    floor-raiser will read it — the `w-mcp-projection` doc's file table is now historical, so the
    durable homes are `design_docs/coding-standards.md` and a comment block at the head of
    `scripts/verify_ail.sh` — with the regeneration recipe (already proven byte-reproducing at base)
    beside it, and `interfaceHash`'s non-movement stated explicitly so nobody "fixes" the third digest.
    **EXPLICITLY NOT THE ITEM: do not make the verifygate control DERIVE its expected count from
    `EXACT_TOTAL_VERIFIED`.** That was the controller's first instinct and the evaluator refuted it:
    a control derived from the value it checks is **vacuous by construction**, which the coding
    standards forbid. The literal is correct *because* it is decoupled and hand-maintained; what is
    missing is the map, not the mechanism. **THE GENERALISATION:** *a gate that is complete by
    construction can still have a coupling surface that is complete only by memory, and the memory
    lives in whoever last raised the floor.* · ~0.1d · gated on nothing · queued iter-127
    (sprint-surfaced; 6th site controller-first-party, evaluator searched for a 7th and found none).]~~

44. **[LANDED 2026-08-27 (iter-133) — ITEM COMPLETE.** PR [#100](https://github.com/sunholo-data/ailang-world/pull/100) → squash [`46add2c`](https://github.com/sunholo-data/ailang-world/commit/46add2c), **Gate 3b GREEN on the MERGE commit** (`present=2` == expected=2 enumerated from `ci.yml`'s own job list, `notgreen=0`, `notdone=0`, `runs=1 event=push`, parent control `checks=2` rev-parsed, `mergeable` read first). Evaluator `sonnet` **97/100 PASS, zero blocking**.
    **The close is a before/after on one step, one merge apart:** parent `80c2bd2` → step `success` over `INSTRUMENT FAILURE`×1 / `RESULT`×0; merge `46add2c` → step `success` over `RESULT: linux/amd64 clean`×1 / `INSTRUMENT FAILURE`×0, same-log control firing at 1 in both. The green is earned rather than manufactured.
    **The item was NOT "flip the flag", and the reason is the fact the row below never stated:** `run.sh`'s floor ORDER meant the linux leg asserted *nothing at all* — the `saw_bad` floor fired second, so the known-good and pinned-OK floors had never evaluated there. Shipped as Option A of four: kernel-read platform (`uname`, because `go env GOOS` honours a `$GOOS` override — measured), a FAIL-CLOSED `case` over exactly the two host pairs this repo has measured, a complete-coverage floor on the known-bad arm, then `continue-on-error` deleted and the wiring pinned by `TestMiscompileInstrumentStepIsGatedInCI`. Constraint 1 ("must not red dev on the next push") was discharged independently by designer, planner and evaluator before being observed true in production.
    **Residues filed rather than absorbed: rows 52 and 53.** Original text follows.]** **THE MISSION'S OWN FIRST-PARTY
    MISCOMPILE INSTRUMENT HAS BEEN FAILING ON EVERY CI RUN, AND `continue-on-error: true` HAS
    BEEN SWALLOWING IT.** Surfaced by `gpt5-6-sol` at quorum round 2 of row 41 as a
    *hypothetical* ("a skipped or BUG-reporting pinned toolchain can still produce a green CI
    result"); **measured first-party, the reality is strictly worse than filed.**
    `design_docs/verification/w-race-gate-blindspot/run.sh` exits **non-zero on 10 of the last
    10 `dev` CI runs** — `INSTRUMENT FAILURE (or GOOD NEWS): no known-affected toolchain
    reproduced the defect` — with **zero** `RESULT:` banners and the same-log control
    (`== go1.26 local-array-literal miscompilation reproduction ==`) firing **10/10**, so the
    script ran every time and failed every time. **CAUSE, two-arm platform control, identical
    `repro/` source, ONE variable:** `go1.26.5` reports **`BUG: Field="" want "stateRoot"`** on
    darwin/arm64 and **`OK`** on the linux/amd64 runner. The miscompilation is darwin-only —
    `scripts/verify_go.sh:214` already records the deny-set as "the measured set … on
    darwin/arm64" — so on `ubuntu-latest` no KNOWN_BAD reproduces, `saw_bad` stays 0, the
    script correctly refuses to certify, and `ci.yml:172`'s `continue-on-error: true` discards
    the refusal. **THE ITEM IS NOT "FLIP THE FLAG":** measured, gating it as-is would red `dev`
    on the very next push, for the platform reason and not for any pin. The question is what
    this instrument should assert **on a platform where the defect does not reproduce** — the
    honest options are a platform-conditional expectation, moving the arm to a darwin runner,
    or making the linux leg assert only the known-good half — and that is a design question,
    not a flag. **THE GENERALISATION:** *`continue-on-error: true` converts an instrument's
    loudest possible output into silence, so the more carefully an instrument fails, the more
    completely that flag hides it.* · ~0.3d · gated on nothing · queued iter-128
    (quorum-surfaced, controller-measured; row 41 scoped its claim rather than absorbing this).

45. **[LANDED 2026-08-28 (iter-134) — ITEM COMPLETE.** PR [#101](https://github.com/sunholo-data/ailang-world/pull/101) → squash [`1d22c79`](https://github.com/sunholo-data/ailang-world/commit/1d22c79), **Gate 3b GREEN on the MERGE commit** (`present=2` == expected=2 enumerated from `ci.yml`'s own `jobs:` block, `notgreen=0`, `notdone=0`, `runs=1 event=push`, parent control `checks=2` rev-parsed, `mergeable` read first). Evaluator `sonnet` **98/100 PASS, zero blocking**.
    **The close is a before/after on one tree one merge apart:** at `ea03cc6` the prefix-less `GOTOOLCHAIN: 1.26.6` leaves Test A **rc=0 `--- PASS`** while `GOTOOLCHAIN=1.26.6 go version` is **rc=1 `invalid GOTOOLCHAIN`** (control `go1.26.6` rc=0); at `1d22c79` the same mutation is **rc=1 with exactly one attributed message** and **zero** of `toolchain pins disagree` / `disagrees with ci.yml toolchain pin` — a measured absence, because the same session's `go1.25.6` arm fires both at 1.
    **The fact this row did not state:** `ci.yml` carries **two** `GOTOOLCHAIN` sites and **two** `go-version` sites, one pair per job; the row named `:21`. The shipped validator runs **per collected value**, so a third site added by the evaluator redded too.
    **Round 2's sharpest objection convicted the doc on its own logic:** the round-1 `t.Errorf` let one malformed value produce **three** messages, two asserting a disagreement — verbatim the misattribution the doc had killed option E0 for. Both arms are now `Fatalf`.
    **A reviewer's mechanism was REFUTED and its fix applied anyway:** `go1.26.6-corp` is accepted by BOTH `version.IsValid` and the runtime (availability class), `go1.26.6_x` rejected by both — 4/4 agreement on the reviewer's own examples — but fifteen agreeing shapes are a sample, so the message text was reframed as this repository's pin POLICY rather than a runtime-equivalence claim.
    **Folded in:** Test B's two `instrument failure:` floors go `Fatalf` (finding B), and the pinned-OK guard gets a comment-stripped CODE-line direction pin requiring exactly 1, two-sided so a flipped guard (0) and a duplicated one (≥2) both red (finding C). Original text follows.]** **THE JUDGE'S
    SHARPEST NON-BLOCKING FINDING ON ROW 41, REPRODUCED FIRST-PARTY.** Row 41's
    `normalizeToolchainPin` blindly prepends `go` to whichever pin kind lacks it — correct for
    `setup-go`'s `go-version: '1.26.6'`, which legitimately omits the prefix, and WRONG for
    `GOTOOLCHAIN`, which does not. Measured: mutating `ci.yml:21` to `GOTOOLCHAIN: 1.26.6`
    leaves `TestGoToolchainPinsAgreeAndMatchJobList` **rc=0 PASS**, while
    `GOTOOLCHAIN=1.26.6 go version` is genuinely **rc=1 `go: invalid GOTOOLCHAIN "1.26.6"`**
    (control: `GOTOOLCHAIN=go1.26.6` rc=0). **Correctly filed NON-BLOCKING** by the evaluator
    and confirmed as such: the malformation fails hard at the first `go` invocation in the job,
    so no drift reaches green — this is a discriminating-power gap in the assertion, not a hole
    in the safety net. **THE ITEM:** validate that `GOTOOLCHAIN` values already carry the `go`
    prefix (or are `auto`/`local`/`path`) rather than auto-correcting them, with the malformed
    value as a named mutation. Fold in the evaluator's other two non-blocking findings: Test B's
    known-positive control uses `t.Errorf` where Test A and the Z3 precedent it claims to mirror
    use `t.Fatalf`, so a broken instrument limps on firing ~10 further assertions instead of
    stopping at the cause; and a comparison-direction flip of the new guard's
    `if [ "$saw_pinned_ok" -eq 0 ]` survives Test B statically (loud at runtime, but not one of
    the 18 named arms). **THE GENERALISATION:** *a normalizer shared across two conventions
    silently imports the laxer one's tolerance into the stricter one.* · ~0.15d ·
    **GATE DISCHARGED 2026-08-27 (iter-132): row 41 LANDED at iteration 129
    ([`8e3c8cd`](https://github.com/sunholo-data/ailang-world/commit/8e3c8cd)), so this row is
    ungated. Measured this iteration, not transcribed — an undated "still gated on" is the whole
    defect the blocked-row rule names.** · queued iter-128 (evaluator-surfaced, controller-reproduced).

46. **[LANDED 2026-08-28 (iter-135) — ITEM COMPLETE.** PR [#102](https://github.com/sunholo-data/ailang-world/pull/102) → squash [`d8c2114`](https://github.com/sunholo-data/ailang-world/commit/d8c2114), **Gate 3b GREEN on the merge commit** (`present=2 == expected=2` enumerated from `ci.yml`'s `jobs:` block, `notgreen=0`, `notdone=0`, `runs=1`, parent control `checks=2`, `mergeable` read first). Evaluator `sonnet` **78/100 PASS**, generator≠judge. `daemonErr` is now a mutex-guarded `syncBuffer` mirroring `host/store/writer_lock_test.go:91`, so all **three** `.String()` reads are synchronised — the row said four branches and there are four, but only three call `.String()`; `<-announced` does not. **The row's own suggested remedy would have deadlocked**: receiving from `waited` in the timeout branch waits for a daemon that has not exited, which is the thing being reported — so the fix is the sink, not the branch. **Claim confirmed by a two-arm control before any code changed**: the bare-buffer pattern is `WARNING: DATA RACE` under `-race` (3 warnings in 5 runs), the guarded sink is `ok` 5/5. **The SET was derived, not the site**: exactly 3 `Start()`-based sinks repo-wide, the other 2 already guarded, every other binding reading only after `Run()`/`CombinedOutput()`. Closed durably by `host/verifygate/subprocess_sink_gate_test.go`, a repo-wide AST gate with 9 fixtures and non-vacuity floors (`filesParsed>=50` observed 110, `startSites>=3` observed 5). **The gate's FIRST version was itself vacuous on three shapes** (`:=`, `new(...)`, closure-split `Start()`) — found by the evaluator, reproduced first-party, hardened in-PR at `49ad39e`; three-arm mutation control on the real tree, all caught, restore byte-identical. **And the arm that exposed it exposed that `go build ./...` does not compile test files**, making the loop's standing mutant-builds assertion vacuous for every test-file mutation (`go vet` is the correct typecheck; measured rc=0 vs rc=1 on an undefined symbol) → skill proposal to V1. **Prior row text follows.** ~~**w-worldd-cli-stderr-buffer-race** · clause-2 · **A DATA RACE ON A TIMEOUT PATH THAT AN
    UNLOADED RIG AND CI BOTH MISS BY CONSTRUCTION — AND MY FIRST ATTEMPT TO REFUTE IT WAS THE
    BROKEN INSTRUMENT.** Reported by the row-41 sprint planner as `verify_go.sh` rc=1 at base.
    My scoped re-run (`-run TestCLIRealSubprocessEpisode -race`) was **rc=0 with 0 race
    warnings**, and the FULL gate was **rc=0** — but that is a different command answering a
    different question, and the source settles it: `cmd/ailang-worldd/cli_test.go:139` binds
    `&daemonErr` as `cmd.Stderr`, `:143` parks `cmd.Wait()` in a goroutine, and the `readErr`
    branch and the 5-second-timeout branch both call `daemonErr.String()` **without having
    received from `waited`** — an unsynchronised read against the live `os/exec` copier, which
    only `Wait()` returning guarantees has stopped. The `case err := <-waited:` branch at `:174`
    is safe for exactly that reason. It therefore fires ONLY when the daemon announcement is
    slow enough to take a timeout branch, i.e. under load — which is why an unloaded rig, and
    CI, both report green. **THE ITEM:** synchronise the read (receive from `waited`, or guard
    the buffer) before reading `daemonErr` on the two unsafe branches. **THE GENERALISATION,
    and it is why this is a row rather than a footnote:** *a non-reproduction is a claim about
    the run, not about the code — when a load-dependent defect is alleged, reading the source
    is cheaper and stronger than re-running on a quiet machine.* Also a live instance of
    iter-97's own guardrail read in reverse: there a loaded rig manufactured a false red; here
    an unloaded rig manufactured a false all-clear. · ~0.2d · gated on nothing · queued iter-128
    (planner-surfaced, controller-confirmed at source level).
    (sprint-surfaced; 6th site controller-first-party, evaluator searched for a 7th and found none).~~

47. **[LANDED 2026-08-28 (iter-136) — ITEM COMPLETE.** PR [#103](https://github.com/sunholo-data/ailang-world/pull/103) → squash [`2a01c43`](https://github.com/sunholo-data/ailang-world/commit/2a01c43), **Gate 3b GREEN on the MERGE commit** (`present=2` **== expected=2** enumerated from `ci.yml`'s own `jobs:` block — and the enumeration was SCOPED to below `jobs:` after a first grep over the whole file over-counted **5**, because the trigger keys `push:`/`pull_request:`/`workflow_dispatch:` share the jobs' indentation; `ls .github/workflows/` = **1**, so the enumeration is complete rather than a hand-picked subset; `notgreen=0`, `notdone=0`, `runs=1`, parent control `checks=2` **rev-parsed, never hand-expanded**, `mergeable` read **FIRST** → `MERGEABLE`/`UNSTABLE`). Evaluator `sonnet` **82/100 PASS, zero blocking**, generator≠judge. **WHAT LANDS:** `ci.yml` declares `workflow_dispatch:` (bare, no `inputs:`), and `host/verifygate/dispatch_lever_gate_test.go` pins it repo-wide — `filepath.Glob` enumeration in the sibling pin-tests' idiom, EXACT column-0 `on:` anchoring, duplicate-`on:` refusal, scalar-value rejection, and an anti-vacuity floor so an empty enumeration fails loudly. **THE CLAIM IS SCOPED, WHICH WAS THE WHOLE DIFFICULTY:** the lever buys **a verdict on a commit**, never **a mergeable PR**, and it re-verifies **the tip of a named ref**, never an arbitrary SHA — `gh workflow run --ref` is documented *"Branch or tag name"* (measured), and `default_branch` is `dev` (measured), which is both why the lever is available and why it covers the dropped-push case. **QUORUM: TWO ROUNDS, BOTH BLOCKED AT FULL STRENGTH** (`absent_reviewers` **EMPTY** both times, **3/3** external reviewers present — neither block is a degraded quorum); round 1 spread across three surfaces, round 2 **localised onto one** (the parser's anchoring), closed under the ratified narrow-refinement carve-out with every fix applied in the reviewers' own words. **AND A REVIEWER'S CITED CONTROL WAS REFUTED BY MEASUREMENT WHILE ITS FIX WAS APPLIED ANYWAY:** `oc-glm-5-2` justified exact-equality by saying it matches P14's `awk '{$1=$1}; $0=="on:"'` "exact match at column 0" — measured, that awk matches an indented `  on:` too (`{$1=$1}` rebuilds the record and strips leading whitespace), while `grep -c '^on:$'` returns 0 vs 1 (control fires). **Neither the parser NOR the premise row that certified it was column-0 aware**, so the objection was STRONGER than filed; fix applied and P14's own command repaired. **Prior row text follows.** ~~**w-ci-recovery-lever-absent** · clause-2 · **THE REPO'S ONLY WORKFLOW HAS NO API-DRIVEN
    TRIGGER, SO A DROPPED WEBHOOK DELIVERY LEAVES A COMMIT PERMANENTLY UNVERIFIABLE — AND THIS
    MISSION HAS NOW PAID TWO ITERATIONS FOR IT.** Measured 2026-08-26 (iter-129):
    `.github/workflows/ci.yml` is the **only** workflow file (`ls .github/workflows/` → 1 entry)
    and `grep -c workflow_dispatch` → **0**, with the positive control firing in the same breath
    (`pull_request` → **1**, `^  push:` → **1**) and a fresh negative-control literal → **0**.
    **The exposure has two halves and only one of them has a workaround.**
    **(a) On a PR the gap is survivable.** Gate 3b's tree-identical empty commit through the git
    API fires a genuine `pull_request: synchronize` without `workflow_dispatch`, and it worked
    first-party this iteration: `runs=0` → `runs=1`, `event=pull_request`, `jobs=2` within 20 s,
    the new commit's `tree` sha asserted byte-identical to the old (`2873892…`) **before** the ref
    was moved. **(b) On `dev` itself there is NO lever at all.** Iteration 128's record push
    [`a0b3162`](https://github.com/sunholo-data/ailang-world/commit/a0b3162) was created inside the
    outage window and still read `checks=0` / `actions/runs?head_sha=<40-char>` `total=0` at
    `19:30Z` the next day — a **TRUE** zero, control rev-parsed not hand-typed (`699f592` →
    `checks=2`, `runs=1`). The incident was marked **resolved at `18:01:30Z`** and the runs did not
    appear, because a dropped delivery is never replayed. The only ways to get a verdict on dev's
    HEAD are to *advance* dev — which changes the commit you were trying to verify — or
    `workflow_dispatch`, which does not exist. Iteration 129 resolved the **instance** forward by
    merging `#97`; the **class** is open, and the next dropped push to `dev` reproduces it exactly.
    **THE ITEM:** add `workflow_dispatch:` to `ci.yml`'s `on:` block, plus a static consistency
    test in the `TestZ3PinDeclaredOnceAndInstalledInBothJobs` /
    `TestGoToolchainPinsAgreeAndMatchJobList` family asserting **every** workflow declares it, so
    the lever cannot be deleted silently — the same shape row 41 just shipped for the toolchain
    pins. **The non-obvious part, and why this is not merely a one-line change:** a
    `workflow_dispatch` run is *not* equivalent to the event it replaces — Gate 3b records that its
    checks do **not** satisfy branch protection on a PR — so the test's claim must be scoped to
    what the lever actually buys (**a verdict on a commit**, never **a mergeable PR**), or it
    becomes exactly the over-broad claim row 41 was scoped to avoid. · ~0.1d · gated on nothing ·
    surfaced iter-129 (paid for twice: iter-128 parked on it, iter-129 recovered around it).
~~

48. **[LANDED 2026-09-01 (iter-145), CI GREEN ON THE MERGE COMMIT [`c13ad1f`](https://github.com/sunholo-data/ailang-world/commit/c13ad1f) via PR [#109](https://github.com/sunholo-data/ailang-world/pull/109)] THE ROW'S FIX WAS RATIFIED BY `D-WORLD-28`, AND THE NUMBER THAT JUSTIFIES THE SPRINT IS THE ONE NOBODY ASKED FOR: **DELETING THE ENTIRE MANDATED BLOCK WAS GREEN IN BOTH LANES**, AND 6 OF 6 GUTTING MUTANTS SURVIVED THE PRE-SPRINT TEST.** **WHAT LANDED:** 3 files, +138/−1, two commits. **P1** ([`2a1d1d7`](https://github.com/sunholo-data/ailang-world/commit/2a1d1d7)) — a runtime fail-closed floor check in `scripts/verify_go.sh` between the miscompile deny-list `esac` and the race leg: the floor is READ from `./go.mod` (never hardcoded), three separately-attributed `FATAL:` refusals, three `exit 1`, zero `exit 0`, and an **`awk` numeric-component comparator rather than a lexical one** — `[[ "go1.9" < "go1.10" ]]` is FALSE in bash, so the naive form fails **OPEN** on a toolchain seventeen minor versions below the floor. **P2** — `:229` now runs `GOTOOLCHAIN="$ACTIVE_GO" go run -race .`, closing the nested auto-selection hole. **P3** — `TestRaceControlFloorStaysBelowRootToolchain` in `host/verifygate`: the floor bound, the P2 needle at count=1, and the six **P1a–P1f** semantic needles. **The fence** ([`45fd7f0`](https://github.com/sunholo-data/ailang-world/commit/45fd7f0)) — 8 comment lines in `racecontrol/go.mod`, `go 1.22` byte-unchanged. **WHY THE NEEDLES EXIST, AND WHY THEY ARE NOT ROW 49's DEFECT AGAIN:** deleting the whole P1 block — the block `D-WORLD-28` exists to mandate — measured **rc=0 GREEN in BOTH lanes** with every other assertion in the sprint still passing (**M11**), while P2's binding already had **M7** for exactly that reason. So the block is split into its comparator half and its gate half (an assertion about one cannot be satisfied by text in the other), **two of the six assertions are NEGATIVE** (zero `exit 0`; zero hardcoded `go1.<d>.<d>` literals) and **one is POSITIONAL** (deny-list < P1 < race leg) — the shapes a `strings.Count` structurally lacks. **THE JUDGE MEASURED THE BEFORE-STATE: 6 of 6** mutants M12–M18 SURVIVE the pre-sprint test and all six are killed by the landed one. **And the finding that makes the needles non-optional:** under the **ambient** toolchain — the only condition CI ever runs — all six gutted variants are byte-for-byte indistinguishable from the correct block, so the runtime lane is blind to every one of them in CI; for **M14** (comparator operands swapped) and **M16** (`exit 1` → `exit 0`, which prints its own FATAL and then exits **success**) even the deliberately hostile M9 runtime arm goes green, making the static needle the **sole killer in the entire sprint**. **M17 is the GREEN control** — a comment word reworded inside the block leaves the test passing, so the set is not "any edit reds". **TWO OF THE DESIGN DOC'S OWN ARMS WERE REFUTED BEFORE A LINE WAS IMPLEMENTED,** by the `opus` planner running the whole battery first: **M6 runs ZERO tests** (Go's module loader rejects the mutated root `go.mod` before the test binary builds — `rc=1 RUN=0 PASS=0 FAIL=0`, so under the doc's own AC10 rule the `!version.IsValid(rootFloor)` branch had no killer at all), and **M10's red is supplied by the pre-existing miscompile deny-list, not by P1** — deleting the root `go` directive drops `GOVERSION` to the deny-listed `go1.26.4`, and the P1-REMOVED control is a **byte-identical red**, i.e. M10 fails the very attribution shape the doc's own V14/A3 exists to enforce. Replaced by **M6′** (`go 1.26.6 // pin` — Go accepts the trailing comment so the module loads and the test runs) and **M10a′/M10b** (TAB-indent / duplicate the directive, `GOVERSION` staying `go1.26.6`), each with an rc=0 / race-leg-reached / 2-races sole-killer control. **AC7's sha256 clause was withdrawn as not computable** (hashes do not compose) for a measured diff-shape form. **THE ITERATION INHERITED A DEAD PREDECESSOR:** iteration 144 ran the round-3 quorum, the carve-out revision and the planner, then died before the executor, leaving an unmerged worktree and two uncommitted doc files and **zero** log entries — see the iteration-145 STATUS. Its work was verified rather than adopted (all four backup sha256s, the base test counts and the restore-proof re-derived first-party) and its one real gap closed: it recorded the planner's nine refutations as PROSE in the Quorum History and never propagated them into the acceptance criteria, the mutation table or the Conflict Surface — which a post-revision re-review (`gpt5-6-sol`, $0.13903) then rejected on, correctly. **GATES, ALL RE-RUN BY THE CONTROLLER OUTSIDE THE CODEX SANDBOX:** `./scripts/verify_go.sh` **rc=0 with an EMPTY failing set**, the `✓ toolchain floor gate` line at log line 23, race control armed with **2** `WARNING: DATA RACE` · `go test ./host/verifygate/ -count=1` **RUN=53 PASS=35 FAIL=0 SKIP=0** at BOTH milestone boundaries (base RUN=52 PASS=34) · `go vet ./...` rc=0 · `go build ./...` rc=0 · `gofmt -l host/ cmd/` 0 bytes · `./scripts/verify_ail.sh` rc=0 at 11 identities / 40 named tests · module census exactly **3** · root `go.mod` byte-identical `7a298361…` · **0** `.ail` files touched · porcelain 0. Milestone commits reconstructed from the executor's cumulative `.snap/` and **proved faithful by sha256** (`bdb3e1eb…` both sides). Evaluator `sonnet` **98/100 PASS, ZERO BLOCKING**; its two non-blocking findings are **NEW ROWS 60 and 61**, not absorbed. **Prior row text follows.** ~~**w-racecontrol-floor-bump-disarms-the-race-control** · clause-2 · **THE SIBLING NESTED MODULE
    HAS THE EXACT DEFECT ROW 42 JUST FIXED, IT WAS RECORDED AS THE *NON-INSTANCE*, AND THE
    REFUTATION CAME FROM A QUORUM REVIEWER RATHER THAN FROM THE AUDIT THAT DECLARED IT SAFE.**
    Filed 2026-08-27 (iter-131) from iteration 130's round-2 quorum objection, measured first-party
    by that iteration's controller and recorded as **V24** in
    [`w-canary-control-does-not-survive-a-floor-raise.md`](implemented/w-canary-control-does-not-survive-a-floor-raise.md).
    Row 42's own Systemic-audit row **V17** examined `design_docs/verification/w-race-gate-blindspot/racecontrol/`
    — the third of the repo's three `go.mod` files — and concluded *"driven only by
    `scripts/verify_go.sh:229` `go run -race .` under the **default** toolchain, so a floor bump
    there cannot disarm anything."* **That trailing clause is REFUTED by measurement.** Base, read
    from inside the module: `go version` = **`go1.26.4`** (the Homebrew base, *not* the root-selected
    `go1.26.6`), `GOTOOLCHAIN=auto`, and `go run -race .` rc=1 with **2** `WARNING: DATA RACE` — the
    gate's known-positive fires. Bumped to `go 1.26.6`: the control **still** fires, rc=1 with 2 race
    lines — **but only because `go1.26.6` is already in the local toolchain cache**, i.e. `auto`
    silently switched toolchains to satisfy the new floor. Add the one variable that removes that
    rescue, `GOTOOLCHAIN=local`: rc=1, output is exactly
    `go: go.mod requires go >= 1.26.6 (running go 1.26.4; GOTOOLCHAIN=local)`, **0** `WARNING: DATA
    RACE` lines — and `verify_go.sh:232` then FATALs *"the race detector is not armed"*. So the
    disarming is real and its only concealment is a warm toolchain cache, which is exactly the
    condition a fresh runner does not have. Restored `ab782f11db0f7f259f73dd55a58eaf5a30b871bb79bd98bacbe964d50efc025b`
    byte-identical, post-restore control re-fires, census re-confirmed complete
    (`find . -name go.mod -not -path './.git/*'` → exactly **3**).
    **THE ITEM:** extend row 42's `TestReproModuleFloorStaysBelowKnownBadToolchains` shape to
    `racecontrol/go.mod` — a static invariant binding its floor at or below whatever the race
    control must run under — or, if the right answer is that `racecontrol` should pin
    `GOTOOLCHAIN=local` explicitly, bind THAT instead. **The generalisable half, and the reason this
    is a row rather than a footnote: an audit that enumerates the *instances* of a defect also
    publishes its *non-instances*, and a non-instance is an unbound claim wearing a measurement's
    clothes.** V17 was a real reading of a real file that stopped one variable short.
    · ~0.1d · gated on nothing · surfaced iter-130 (quorum round 2), filed iter-131.~~
49. **[LANDED 2026-08-31 (iter-139), CI GREEN ON THE MERGE COMMIT [`7ab42aa`](https://github.com/sunholo-data/ailang-world/commit/7ab42aa) via PR [#106](https://github.com/sunholo-data/ailang-world/pull/106)] The `PARKED-ON-LANE` predicate iteration 138 wrote was machine-checkable and it FIRED: the evaluator lane was re-probed after the Monday 07:00 local reset, accepted, and ran. Evaluator `sonnet` **96/100 PASS, ZERO BLOCKING** in its own worktree (generator≠judge holds three ways: Codex executor, Sonnet judge, Opus controller — unlike iteration 138, whose Codex controller could not stand in). **THE JUDGE MEASURED THE *BEFORE* STATE, WHICH IS THE NUMBER THAT MAKES THIS ITEM WORTH ITS COST:** it reverted only the call-site hunk and re-ran all 23 mutations against the old `strings.Count(src,"stateRoot") >= 2` check — **22 of 23 SURVIVED IT**, and the one that did not (full-function deletion) reds only because it removes the token, never because it detects a missing assertion. The fence being replaced was almost entirely vacuous, measured rather than argued. **23/23 shape-gutting mutations are killed** (the doc's own M1–M13 plus ten adversarial arms the judge added: helper extraction, closure/IIFE, `for`-wrap, `t.Run`-wrap, bare-block wrap, `t.Fatalf`→`t.Errorf`/`t.Logf`, `if false && cond`, and the shape planted inside a string literal). **ONE SURVIVOR, REPRODUCED FIRST-PARTY AND PRE-DECLARED:** `t.Skip()` as the canary's first statement leaves the AST byte-identical, so the canary reports `--- SKIP` (RUN=1 PASS=0 SKIP=1) while the fence stays `rc=0 RUN=1 PASS=1`. Mutant landed by sha256, `go vet` rc=0, restore byte-identical, pristine control green after. It is inside the doc's own declared residual (*"the shape is not the behaviour"*), so it is NON-BLOCKING — and it is routed to **NEW ROW 56** rather than absorbed. Controller gates re-run OUTSIDE the executor sandbox on the rebased head: `gofmt -l host/ cmd/` **0 bytes**; `go vet ./...` rc=0; `go build ./...` rc=0; the named fence **RUN=1 PASS=1 FAIL=0**; AC1's own vacuity self-test fires (the nonsense selector returns rc=0 with **RUN=0** and `[no tests to run]` — exactly the false green the run-existence form excludes); `verify_go.sh` rc=0 (plain AND raced full suites, `AILANG_BIN` v0.30.0 pinned); `verify_ail.sh` rc=0 at 11 identities / 40 named tests / 9-of-9 package steps with **0 `.ail` files touched**; porcelain 0. Branch rebased onto current `dev` first — the incoming/branch file intersection is **empty** with the control firing at 4, so no conflict was possible. Prior head follows.** ~~**[PARKED-ON-LANE 2026-08-31 (iter-138): evaluator `sonnet` unavailable until the Anthropic weekly reset at Monday 07:00 local; resume by re-probing the evaluator lane, then verify-and-land branch `sprint/world-iter138-canary-fence` at `7d51e02`. NO HUMAN DECISION.** The inherited iteration-138 sprint is complete in three milestone commits (`ca2ecd6`, `345a73a`, `7d51e02`) and clean, but the previous controller died immediately after creating the evaluator worktree when its provider returned a usage-limit error. This iteration re-derived the diff and re-ran controller gates outside the executor sandbox: gofmt 0 lines; scoped vet/build rc=0; both named tests exactly 1 RUN/1 PASS; `verify_go.sh` PASS including plain+raced full suites; `verify_ail.sh` PASS at 11 identities/40 named tests/9 package steps; zero `.ail` changes; worktree porcelain 0. Generator≠judge forbids substituting this Codex controller for the Codex executor's missing Sonnet judge. Prior head follows.** ~~**w-canary-fence-passes-a-gutted-canary** · clause-2 · **ROW 42's OWN POSITIVE CONTROL IS
    SATISFIED BY TEXTUAL PRESENCE, SO THE MISSION'S DARWIN MISCOMPILE DETECTOR CAN BE REDUCED TO A
    NO-OP WITH EVERY GATE IN THE REPO GREEN.** Filed 2026-08-27 (iter-131) by the `sonnet` evaluator
    of `P42` as non-blocking finding 1, **reproduced first-party by the controller before adoption**
    and sharper on re-measurement than as filed. `TestCanaryDeclaresPositiveArmOnly` guards its
    marker clause with `strings.Count(src, "stateRoot") >= 2` — a known-positive control, correctly
    reasoned, and it counts **occurrences of a string** rather than the presence of an *assertion*.
    Repro: delete `host/store/toolchain_canary_test.go`'s three real lines
    (`if rows[0].field != "stateRoot" { t.Fatalf(...) }`), replacing them with a comment that still
    names `stateRoot` once. The needle count falls **3 → 2**, which still satisfies `>= 2`. Measured:
    the fence `TestCanaryDeclaresPositiveArmOnly` → **`--- PASS`**; `TestToolchainCanary` itself →
    **`--- PASS`**, now asserting nothing; and the full
    `go test ./host/verifygate/ ./host/store/ -count=1` sweep → **rc=0, zero FAIL lines**. Restored
    byte-identical (`a23cfa79419ae691…`), porcelain clean. **This does NOT falsify row 42's claim** —
    the Decision scopes Test B to *"the canary's assertion has moved out of this file"* plus the
    known-bad-arm fence, and the doc's residual section says it fences one token rather than
    completeness — which is why the evaluator's non-blocking severity is right. It is a real
    exposure nonetheless, because the artifact being protected is the *only* first-party detector
    this mission has for the darwin/arm64 array-literal miscompilation.
    **THE ITEM:** replace the presence count with a shape assertion — require the comparison against
    `"stateRoot"` and a `t.Fatalf` in the same function — or, better, make the fence behavioural
    rather than textual by having `host/verifygate` execute the canary and require it to red on a
    seeded-wrong value. **The generalisable half: a known-positive control that counts a TOKEN
    proves the file still mentions the subject, never that it still tests it** — the same distance
    between prose and code the mission recorded at iter-124 (`Authorization` = 1, in a comment).
    · ~0.1d · gated on nothing · surfaced iter-131 (evaluator-found, controller-reproduced).~~

50. **[RE-PARKED `needs-human-review` 2026-09-01 (iter-146) ON THE NEW `D-WORLD-31` — UNPARKED BY `D-WORLD-29` (Mark, attended 2026-09-01), REVISED TO THE RATIFIED DIRECTION, AND THEN BLOCKED AT A FULL-STRENGTH ROUND 3 BY A DEFECT THE RULING'S OWN RATIONALE DOES NOT COVER. THE RULING PICKED **A**, WHICH REVERSES THE DESIGN DOC AND RESTORES THIS ROW'S OWN ORIGINAL PRESCRIPTION — AND IT EXPLICITLY AUTHORISES THE TEXT RETRACTION BELOW, WHICH IS THE FIRST TIME THIS LOOP HAS BEEN GIVEN LEAVE TO EDIT RATIFIED QUEUE TEXT.** Ruling, verbatim: *"ANSWERED — A — ACCEPT the indented assignment. One whitespace-tolerant scan: trim leading spaces/tabs, ignore comments, extract NAME=\"...\" regardless of indentation, require len(values) == 1. B's stated rationale (\"an indented assignment is a conditional/nested shadow\") is measured FALSE, and defence-in-depth bought with a false premise is not worth a loud red on a benign re-indent. AMEND queue row 50 in the same iteration to drop the sentence carrying that premise — it is ratified text the loop could not edit itself, and this ruling authorises exactly that edit."* Provenance checked first-party under the shared skill's ATTENDED LEDGER EDITS contract: `git log -1 -S'| D-WORLD-29 | RESOLVED'` on this charter returns author `Mark Edmondson <mark@aitanalabs.com>`, **not** the fleet account `sunholo-voight-kampff`, so it is a human answer and not a self-resolution.

    **THE RETRACTION, MADE UNDER THAT AUTHORISATION.** The struck-through original row text below lists *"an indented **only** assignment likewise reads 0 and fatals"* among the helper's **correct loud behaviours**. **That clause is RETRACTED as a statement of correctness.** The observation is accurate — the shape does read 0 and does fatal — but the premise attached to it is false: it fatals because the scanner cannot **SEE** a real assignment, not because anything is wrong with the script. Measured first-party (iter-140, re-derived iter-146): a bash script setting a column-0, a space-indented and a tab-indented assignment at top level prints **all three**, rc=0, while the paired control shows a genuinely *conditional* assignment does change behaviour. Under ruling A that red disappears by design. Nothing else in the struck text is retracted, and the strike-through is preserved per this charter's own dead-head convention.

    **THE DOC WAS REVISED TO DIRECTION A AND RE-QUORUMED.** `design_docs/planned/w-shell-assignment-parser-drops-an-indented-assignment.md`, 623 → 678 lines, revised by the rotation designer `pi:ollama/deepseek-v4-flash:0731-cloud` (typed verdict `ok`, 112 s, 1 changed file, 9 tool executions, non-empty worktree diff). Rule B is gone from every normative section, not annotated: the High-Impact Decisions row now records **human** as `Chosen By` and cites `D-WORLD-29`; the Design Freeze no longer claims "no human ratification is required"; the helper keeps its signature `func shellAssignmentValues(lines []string, name string) []string` with **no second return value and no new assertion at any of the four call sites**; the four `assignment count=%d, want 1` messages stay byte-unchanged and do all the work; **AC2** now pins `KNOWN_BAD assignment count=2, want 1` at both consumers and **AC3 is REVERSED** to require the indented-ONLY arm to be GREEN. Controller B-residue sweep on the revised doc: `indented int` **0**, `want 0` **0**, `New assertion` **0**, `indented assignment count` **0**, with the same-file known-positive control `shellAssignmentValues` firing at **13** and a fresh negative control at **0**; the two surviving `two-sided` hits are deliberate narration of the rejected direction.

    **THE ROW'S OWN DEFECT IS MEASURED CLOSED, AND THE MEASUREMENT IS THE POINT.** All arms at `cb73cab` in a scratch worktree, mutant asserted LANDED by sha256, intended effect asserted against the system's own view rather than against bytes, `bash -n` rc=0, `go vet` rc=0 read **before** any test result, restored byte-identical from a `cp` backup, pristine control green before **and** after the batch, porcelain 0. **Base helper:** the silent arm (a SECOND, indented `KNOWN_BAD` beside the column-0 one, inside `if true; then … fi`) is **rc=0 RUN=2 PASS=2 — the gate is BLIND** while bash genuinely narrows the deny-list to the last assignment; the loud arm (indented-ONLY) is **rc=1 RUN=2 PASS=0 FAIL=2**, `count=0, want 1` at `:267` and `:351`. **Ratified rule-A helper applied verbatim:** the silent arm becomes **rc=1 RUN=2 PASS=0 FAIL=2** with `KNOWN_BAD assignment count=2, want 1` at `:269` and `:353` — **the row is closed by the pre-existing message alone**.

    **AND THE CONTROLLER FOUND A RESIDUAL THE DOC DID NOT DECLARE — THE ONE SHAPE WHERE RULING A'S OWN RATIONALE DOES NOT REACH.** The ruling justifies A on the ground that an indented assignment's value is *"what bash actually does"*. That is true for a **top-level** indented assignment, which is what was measured when the ruling was written. It is **false for a multi-line branch that never executes**: with the only `KNOWN_BAD` inside `if false; then` / `fi` across three lines and no column-0 copy, the base helper is **rc=1 `count=0, want 1`** (correctly loud) while the rule-A helper is **rc=0 RUN=2 PASS=2** — so the toolchain floor is bound against a deny-list that does not exist at runtime, and the runtime control `bash -c 'if false; then KNOWN_BAD="x"; fi; echo "${KNOWN_BAD:-<UNSET>}"'` prints **`<UNSET>`**. The paired control in the same batch proves the helper is not merely permissive: the silent arm under that same helper still reds at `count=2, want 1`. The doc's residual 4 covers only the **one-line** `if …; then NAME="x"; fi` form; the multi-line form is the shape a refactorer actually writes, and it was declared nowhere. **This does NOT reopen `D-WORLD-29`** — it is a narrowing of loudness that ruling A knowingly accepts, and the row's own defect is closed either way — so it is now **declared residual 6, pinned by unit arm 11**, and recorded as `V-ARMC1`/`V-ARMC2` with its runtime control, rather than being fixed (closing it needs branch-reachability analysis, out of proportion for a ~0.1d row) and rather than being left silent.

    **ROUND 3 BLOCKED AT FULL STRENGTH, AND ONE OF THE TWO OBJECTIONS IS WHY THIS ROW IS PARKED RATHER THAN SHIPPED.** 3 of 3 reviewers PRESENT, `.synthesis.absent_reviewers` **`[]`** cross-checked by `[.reviewers[]|select(.present==false)]`, verdicts read at the NESTED `.result.verdict` path; `metered=$0.15101`. `gemini-3-1-pro` **pass**. **`oc-glm-5-2` reject — a PREMISE objection, so it was MEASURED rather than forwarded (rule 3f), and it was RIGHT:** AC7 hardcoded a 4-name base `--- FAIL` set measured at `9c0ad0b` while this doc targets `cb73cab`. Run twice on the IDENTICAL pristine worktree at `cb73cab` (porcelain 0 between runs, known-positive control 19 `ok`/`---` lines): **run 1 rc=1 with `{TestHandlerTimeoutKillsTheWholeProcessGroup}`, run 2 rc=0 with an EMPTY set.** The hardcoded set matches **neither**, and the gate is confirmed **FLAKY on one tree** — queue row 58 corroborated first-party, from the opposite direction to the way it was filed. Fixed with the reviewer's own verbatim alternative: AC7 part 2 now takes its base reading **at drill start**, in the same worktree at the same commit, twice, and treats only names common to both runs as the base — recorded as `V-ARM-GO`. **`gpt5-6-sol` reject — this one disputes the DIRECTION**, on `V-ARMC2`: *"A verification gate that certifies a deny-list that does not exist at runtime … labeling this a residual does not make the regression acceptable,"* with the verbatim disposition *"If that migration is out of scope, keep the row blocked rather than accepting the demonstrated V-ARMC2 false green."* **The carve-out's own precondition therefore fails** (it forbids applying a fix over a contested DIRECTION), the one-revision-one-requorum allowance is spent, and choosing between "accept the residual" and "migrate to a declarative fixture" is controller judgment over a contested direction — Standing rule 2. **This does NOT reopen `D-WORLD-29`**, which the recording contract forbids; it is filed as the new, uniquely-named **`D-WORLD-31`**, because the fact it rests on was measured AFTER the ruling and the ruling's stated rationale (*"which is what bash actually does"*) is true of the shape Mark was shown and false of this one. **Prior head text follows.**

    ~~[PARKED `needs-human-review` 2026-08-31 (iter-140) on D-WORLD-29 — DESIGN DOC WRITTEN AND BANKED; TWO QUORUM ROUNDS, BOTH BLOCKED; THE SURVIVING OBJECTION DISPUTES THE DESIGN DIRECTION AND ITS PREMISE MEASURES TRUE.** Doc: `design_docs/planned/w-shell-assignment-parser-drops-an-indented-assignment.md` (Fable designer; authoring + one protocol-mandated revision). **THE ROW'S OWN PRESCRIPTION AND THE DOC'S DESIGN DISAGREE, AND THE ROW IS THE ONE THAT NEEDS AMENDING.** The doc refused this row's *"tolerate leading whitespace, require the total to be 1"* on the ground that it would **weaken** the gate — an indented-only assignment reds loudly today and would go green — and proposed a two-sided invariant instead: column-0 count == 1 **AND** indented count == 0. Quorum round 2 rejected that on a shell-semantics premise the controller then MEASURED first-party: **indentation is not syntax in bash.** A script setting `COL0=`, a space-indented `INDENTED=` and a tab-indented `TABBED=` at top level prints all three values, rc=0, while the paired control shows that a genuinely *conditional* assignment does change behaviour (`NARROW` unset → `wide list`, set → `narrow`). So an indented top-level assignment really does execute; the doc's rationale (*"an indented one is a conditional/nested shadow"*) is **FALSE**; and today's loud red on the indented-only shape fires because the scanner cannot **SEE** a real assignment, not because anything is wrong. **The row itself encodes the same false premise** — it lists *"an indented **only** assignment likewise reads 0 and fatals"* among the correct loud behaviours — which is why this is a human decision rather than a controller routing call: amending ratified queue text is not the loop's to do. **WHAT THIS ITERATION BANKED, AND THE NEXT ONE INHERITS:** the silent arm and the loud arm, both measured first-party BEFORE the doc existed — a second, indented `KNOWN_BAD` beside the valid column-0 one leaves both consumers **rc=0 RUN=2 PASS=2** (the gate is blind), while the same assignment indented-**only** gives **rc=1 RUN=2 PASS=0 FAIL=2** with `assignment count=0, want 1` at both sites — with mutants asserted LANDED by sha256, `bash -n` rc=0, the intended effect asserted against `grep -c` rather than against bytes, restore byte-identical, pristine control re-passing, porcelain 0. **QUORUM WAS FULL-STRENGTH IN ROUND 2, NOT DEGRADED:** `oc-glm-5-2` was recorded absent (`invalid`) in BOTH rounds and was RESTORED by a single-reviewer re-run at a raised cap — it returns **PASS** ($0.0203), so round 2 is 2 rejects against 1 pass with `absent_reviewers` closed. Its own non-blocking objection is worth carrying: declared residual 1 (`export KNOWN_BAD=` inside the same `if`) leaves the **IDENTICAL** silent-narrowing hole open under EITHER candidate rule. **Both round-1 objections were measured rather than forwarded:** `gemini-3-1-pro`'s *"fabricated verification log"* was CONFIRMED as a transcription defect and no worse — markdown table-cell `\|` escaping made every ERE alternation match a literal pipe, so the commands as printed returned rc=1 / 0 hits while their conclusions were true (controller re-derived: rc=0, 3 hits at lines 24/25/26) — fixed by moving every runnable command into fenced blocks, after which the controller executed the doc's own repaired strings VERBATIM and reproduced every recorded observation; `gpt5-6-sol`'s *"is there shell-parsing machinery to reuse?"* was REFUTED by measurement (`go.mod` has ONE direct dependency, `modernc.org/sqlite`; zero shell-parsing libraries; the sole adjacent tokenizer is an unexported, different-grammar test helper in another package's `_test.go` — 4 hits against a same-scope known-positive control of **59** `HasPrefix`), and the audit it asked for, which had genuinely never been run, is now published in the doc. **Prior row text follows.** ~~**w-shell-assignment-parser-drops-an-indented-assignment** · clause-2 · **`shellAssignmentValues`~~
    IS COLUMN-0-ANCHORED, SO A SECOND `KNOWN_BAD=` NESTED INSIDE A SHELL BLOCK IS INVISIBLE TO IT —
    SILENTLY, WHILE EVERY OTHER MALFORMED SHAPE FAILS LOUDLY.** Filed 2026-08-27 (iter-131) by the
    `sonnet` evaluator of `P42` as non-blocking finding 2. The helper (`host/verifygate`, inherited
    unchanged from `P41` and reused by row 42's new Test A) matches on
    `strings.HasPrefix(line, name+"=\"")`, so an indented conditional assignment inside a function
    is not counted at all. Today `run.sh` is flat and sets `KNOWN_BAD` exactly once at column 0, so
    nothing is currently wrong — this is a latent gap, not a live defect. **What makes it worth a row
    is the asymmetry the evaluator measured**: every *other* deviation the helper meets fails
    LOUDLY through the instrument-failure floors — a single-quoted `KNOWN_BAD='…'` returns 0 values
    and trips `assignment count=0, want 1`; a commented-out `# KNOWN_BAD=` is correctly ignored; an
    indented *only* assignment likewise reads 0 and fatals. The one silent case is a **second,
    indented** assignment beside a valid top-level one: the count stays `1`, the test stays green,
    and Test A would then bind the floor against a stale or partial deny-list. **THE ITEM:** count
    assignments with leading whitespace tolerated and require the total to be 1 — so a refactor that
    adds a conditional branch reds loudly instead of narrowing the deny-list in silence — and add
    one arm per shape to the drill. · ~0.1d · gated on nothing · surfaced iter-131 (evaluator-found).~~


51. **[LANDED 2026-08-31 (iter-141), CI GREEN ON THE MERGE COMMIT [`d195717`](https://github.com/sunholo-data/ailang-world/commit/d195717) via PR [#108](https://github.com/sunholo-data/ailang-world/pull/108)] THE ROW'S OWN GENERALISATION WAS MEASURED, NOT ARGUED: the `sonnet` evaluator reverted only the new assertion block — leaving the helper compiled but unused — and re-ran the arms against the OLD test. **7 of 9 core arms SURVIVED it**; the two that died (R1, R2) died on the pre-existing per-row loop, untouched by this change. Excluding A3, which is expected-PASS in both states and is not a defect, **6 real previously-invisible mutants are newly caught.** **WHAT LANDED:** one file, two commits, +79/−0 on `host/verifygate/floor_raise_inventory_test.go`. M1 [`dbf343d`](https://github.com/sunholo-data/ailang-world/commit/dbf343d) — `siteReScript`/`siteReTable`, `canonicalSiteSet`, a per-home duplicate guard, and a `>= 6` cardinality floor whose instrument-failure `t.Fatalf` is what stops an empty enumerator result passing vacuously; arm A1 still PASSES at this boundary, measured by the planner and reproduced by the executor — **M1 buys anti-vacuity, M2 buys detection**, and saying so up front is what stopped a correct boundary reading as a failure. M2 [`0ee301a`](https://github.com/sunholo-data/ailang-world/commit/0ee301a) — the SYMMETRIC set-difference, both directions, so an asymmetric addition to EITHER home reds with the divergent site NAMED. **THE ROUND-2 QUORUM IS WHY THE COMPARATOR IS SYMMETRIC AT ALL.** The designer's first sketch checked ONE direction (every S8 site in `scriptSet`) plus length equality. Two reviewers and the controller independently found it is not set equality, and the controller MEASURED the evasion on the real six-path base set: script gains a genuinely new site while S8 gains a DUPLICATE of an existing one → both length 7, every S8 site present → **GREEN**, sets genuinely different. It also made the doc's own AC4 unsatisfiable: with arm A1 landed only the cardinality branch could fire, and that message names counts, never the site. Both are closed by the reverse-direction loop plus a per-home duplicate `t.Fatalf`. **FIFTEEN ARMS, ALL RUN, NONE PREDICTED** — the planner executed the whole drill before writing the plan and repaired six defects in the design doc's own drill text (C1 is a no-op as written: Go redeclaration is package-scoped, so a collision planted in another package leaves both vets rc=0; A4's effect gate must NOT move, violating the drill's own must-MOVE rule; R2's named row 3 is confounded as the sole S8 home of the `REQUIRED_VERIFIED` sharedNeedle; AC8 needs a dual per-row-plus-floor verdict). Every arm: LANDED by sha256 AND by an intended-effect query against the system's own view, `go vet` rc=0 read BEFORE any test result, restored byte-identical from a `cp` backup, pristine control green either side. **The planner's own first N2a/N2b attempt silently failed to land and reported both arms backwards — the landed-proof is what caught it.** **PROTECTED FILES WERE MUTATED AND RESTORED, NEVER DELIVERED:** `design_docs/coding-standards.md` §S8 is ratification-class and was never edited; both `scripts/verify_ail.sh` (`5a1bbe89…`) and the standards file (`b710a510…`) end byte-identical to base, re-asserted after every arm. Milestone commits were reconstructed from the executor's cumulative `.snap/` snapshots and proved faithful by sha256 (`d9bd7a81…` both sides). Controller gates re-run OUTSIDE the codex sandbox: `gofmt -l host/verifygate/` 0 bytes · `go vet` rc=0 · `go build ./...` rc=0 · `go test ./host/verifygate/ -count=1` **52 RUN / 52 PASS / 0 SKIP** at BOTH milestone boundaries · `./scripts/verify_ail.sh` rc=0 · `./scripts/verify_go.sh` **rc=0 with an EMPTY failing set** · 0 `.ail` files touched · porcelain 0. Evaluator `sonnet` **97/100 PASS, ZERO BLOCKING**, in its own worktree. **ONE NON-BLOCKING FINDING, ROUTED TO NEW ROW 59 RATHER THAN ABSORBED.** Prior row text follows.** ~~**w-inventory-test-blind-to-asymmetric-addition** · clause-2 · **EVERY ARM OF THE ROW-43
    DRILL IS REMOVAL-SHAPED, AND A REMOVAL PROVES THE CHECK *FIRES* WHILE ONLY AN ADDITION
    PROVES IT *LOOKS*.** Found by the `sonnet` evaluator at round 2 of row 43, after it had
    already re-run all twelve row-deletion arms and five attacks on the new uniqueness
    assertion — every one of which held. It then added a **fabricated seventh coupled-file
    row** (`some_new_file.go`) to the `scripts/verify_ail.sh` home ONLY, absent from
    `coding-standards.md` §S8. `TestFloorRaiseInventoryNamesEveryCoupledFile` **PASSES**, and
    **no count moves**. **THIS IS A DECLARED RESIDUAL, NOT A BROKEN PROMISE:** the design's
    scope is the six named Tier-1 sites, and the test was never built to bound *cardinality* or
    to require the two homes to be *symmetric* — which is exactly why it is a row rather than a
    blocking finding. **THE ITEM:** decide whether the two homes must agree on their SITE SET
    (not merely on the presence of each known site), and if so bind it — a count equality
    between the homes, a set-difference assertion, or a cardinality floor with an
    instrument-failure branch. Note the trap in the obvious fix: an assertion derived from one
    home's contents cannot detect that home being wrong, so the comparison must be between the
    two homes, or against an independently authored list. **ALSO IN SCOPE, from the same
    evaluation:** the shared bare-token needle layer (`REQUIRED_VERIFIED`, `EXACT_TOTAL_VERIFIED`
    and the rest) still permits duplication — it is explicitly secondary to the row anchors and
    no attack exploited it, but it is the residue of the very class row 43 closed, left standing
    one layer down. **THE GENERALISATION:** *a gate hardened by deletion is hardened against
    deletion; the thing it has never been shown to notice is the thing that was added while
    nobody was looking.* · ~0.1d · gated on nothing · surfaced iter-132 (evaluator-found,
    controller-confirmed, declared out of scope by the design it was found in).

52. **[LANDED 2026-09-02 (iter-147), CI GREEN ON THE MERGE COMMIT [`a1744d3`](https://github.com/sunholo-data/ailang-world/commit/a1744d3) via PR [#110](https://github.com/sunholo-data/ailang-world/pull/110)] `D-WORLD-30` RATIFIED OPTION A (line scan, hardened, anchored on the SHALLOWEST enclosing `steps:`); the doc was revised to it, blocked once more at a FULL-STRENGTH round 3, and closed under the narrow-refinement carve-out with all three objections MEASURED rather than forwarded.** **WHAT LANDED:** one function in one file, `host/verifygate/toolchain_pin_gate_test.go`, +63/-13. The locator derives `stepCol` from the shallowest enclosing `steps:` key, pins it against `expectedStepCol = 6` with a loud `t.Fatalf`, locates the block by exact indentation rather than by a token, and refuses loudly on containment (Invariant A) and identity (Invariant B). **THE ROW'S OWN CLAIM WAS RE-DERIVED, NOT INHERITED** (it was this loop's iter-142 measurement): **ARM A** — flag on the previous *unrelated* step plus a key reorder — was **rc=1 blaming `ci.yml:164`, the unrelated step's own flag line**, and is now **rc=0 GREEN**; **ARM B** — flag live on the guarded step plus a `- name:`-shaped decoy inside its own `run: |` scalar — was **rc=0 `--- PASS` over a live forbidden flag** and is now **rc=1 @ `ci.yml:181`**, its mutant sha matching the doc's pin `c1903a86…`. **A REFINEMENT THE ROW DID NOT STATE:** the same decoy placed BEFORE the path line is still caught (`ARM B0`, rc=1) — the backward scan anchors ON it and widens — so the construct does not always fail open; the decoy's POSITION decides it. **THE COUNTERFACTUAL, MEASURED BY THE JUDGE:** reverting only the locator hunk leaves **7 of 14 arms newly load-bearing** (ARM A, ARM B, ARM D, ARM G's wholly new identity invariant, and Q/Q2/R, which the old scan absorbed at rc=0). **TWO OF THIS ITERATION'S OWN CLAIMS WERE REFUTED BY ITS OWN LANES, which is the loop working:** the `opus` planner measured that the `expectedStepCol` assertion the controller added at round 3 **could not fail** on the arm the doc named (the cross-job capture also yields 6, so Invariant A does all the work — neuter probe FATAL→FATAL, no flip), and replaced it with **AC8′** on new arm **MUT-R**, where the pin fires `stepCol=8` and neutered goes rc=0 GREEN; it also found MUT-Q's landed-assertion unsatisfiable as written (at its own pinned sha the `go-verify` `steps:` key parses to **NIL**). **THE DURABLE ARTIFACT IS DECLARED RESIDUAL 8**, produced by measuring `gpt5-6-sol`'s round-3 objection rather than forwarding it: the backward `steps:` scan is **unbounded**, so on a re-indented file the shallowest anchor **leaves the JOB**. It is recorded as a residual rather than escalated because it fails **LOUD, never silent** — the judge could not construct any silent green, since the `count==1` uniqueness fatal forbids duplicating the identifying line into another job. Executor `codex:gpt-5.6-sol` (17 arms, zero git writes); evaluator `sonnet` **93/100 PASS, ZERO BLOCKING** in its own worktree; its two non-blocking findings are **NEW ROWS 62 and 63**, reproduced first-party before filing. **Prior row text follows.** ~~**[PARKED `needs-human-review` 2026-09-01 (iter-142) on D-WORLD-30 — DESIGN DOC WRITTEN AND
    BANKED; TWO FULL-STRENGTH QUORUM ROUNDS, BOTH BLOCKED; THE SURVIVING OBJECTION DISPUTES THE
    DESIGN DIRECTION AND THE TWO REVIEWERS DISAGREE WITH EACH OTHER ABOUT THE REMEDY.** Doc:
    `design_docs/planned/w-wiring-test-step-scoping-imprecise-under-key-reorder.md` (627 lines, 19
    verification rows; Fable designer, authoring + one protocol-mandated revision).
    **THIS ROW'S OWN HEADLINE CLAIM IS REFUTED — the defect is exploitable in BOTH directions, and
    the amendment is the durable artifact.** The row records it as *"measured non-exploitable"*
    because *"every reordering the judge tried widened the scanned range; none narrowed it."*
    Measured first-party at iter-142, `.github/workflows/ci.yml` only, both mutants LANDED by
    sha256 (`aed8e186…`→`28731ce9…`/`66bad3af…`), `go vet ./host/verifygate/` rc=0 each, effect
    asserted against `ruby -ryaml` (the system's own view) rather than bytes, restored
    byte-identical, pristine control green either side, porcelain 0:
    **(A) FALSE POSITIVE.** `continue-on-error: true` on the PREVIOUS, unrelated `go build + test
    gate` step plus a key reorder of the miscompile step → **rc=1 at `ci.yml:166`**, and `166` IS
    the unrelated step's own flag line (YAML view: `step[5] coe=true`, `step[6] coe=nil`, step
    count unchanged at 10). That violates the boundary the test's own comment cites from quorum
    round-2 R1 — *"a flag on an unrelated step is that step's business"*.
    **(B) FALSE NEGATIVE — THE GATE FAILS OPEN.** The flag live on the miscompile step plus a line
    whose trimmed text begins `- name:` inside that step's own `run: |` block → **rc=0 `--- PASS`**,
    with the YAML view confirming `continue-on-error=true` on the guarded step and the decoy being
    script CONTENT. `strings.TrimSpace` discards indentation, so nesting depth is invisible to the
    walk-back. **So "it happened to widen" was a property of the mutations someone tried — exactly
    as this row predicted of itself — and (B) is the counterexample.**
    **SCOPE CORRECTED:** the step-scoping defect has **ONE** call site, not "three places"
    (`HasPrefix(strings.TrimSpace` → 3 hits: the same function's start and end scans plus an
    unrelated SQLite `busy_timeout` pragma; same-scope known-positive control **6**). Exactly
    **three** Go files read `ci.yml`; the other two do whole-file counting with their own controls
    and no step scoping. Row 41's V18 re-verified at this head: `actionlint` appears in **no** repo
    gate (control firing at 1) though the binary is on the rig — a rig fact, not a repo fact.
    **WHY IT PARKS.** Round 1 BLOCKED (reject/reject/pass + controller reject, $0.1035); one
    revision; round 2 BLOCKED with **all three** external reviewers rejecting ($0.1489).
    `absent_reviewers` **`[]`** in both rounds, cross-checked, so both are full-strength. Every
    objection was measured rather than forwarded: the doc's round-1 relative-indentation locator is
    spoofable (`start=177 (indent 12)`, both invariants PASS, flag excluded → **GREEN/blind**), and
    so is its round-2 absolute-`stepCol` locator, because the anchor is itself an indentation-blind
    backward search for `steps:` — a decoy `steps:` inside the `run:` scalar yields
    **`stepCol=12` → GREEN/blind**, while a shallowest-`steps:` derivation on the identical file
    yields **`stepCol=6` → RED/caught**. The one-revision-one-requorum allowance is spent, and the
    reviewers' `proposed_fix` fields **disagree**: `gemini-3-1-pro` says keep the line scan and pin
    `stepCol == 6` as a loud invariant; `gpt5-6-sol` says replace it with semantic YAML traversal
    and, verbatim, *"if introducing a YAML dependency is not acceptable, park or widen the item
    rather than landing another scanner that remains attacker-controlled by block-scalar text"* —
    which disputes the DIRECTION, so the narrow-refinement carve-out's own precondition fails. The
    surviving line-scan fix is **controller-invented**, which the carve-out forbids by name, so it
    becomes arm A of `D-WORLD-30` rather than an edit. `oc-glm-5-2` additionally caught a falsified
    citation (V12 is cited for `indentOf` being unallocated and its greps never search that name) —
    a one-line repair owed to whichever arm wins. ~~ **Prior row text follows.** ~~**w-wiring-test-step-scoping-imprecise-under-key-reorder** · clause-2 · **THE ROW-44 WIRING
    TEST FINDS ITS STEP BY WALKING BACK TO THE NEAREST `- name:`, SO A STEP WHOSE `- name:` IS
    NOT ITS FIRST KEY IS ATTRIBUTED TO THE PREVIOUS STEP.** Found by the `sonnet` evaluator at
    the row-44 landing and **measured non-exploitable, which is why it is a row and not a
    blocker**: it mutated the miscompile step to `- continue-on-error: true` / `  name: …` so
    that `start` resolved to the *previous* step's `- name:`, and the check still caught the flag
    (rc=1, `ci.yml:173 re-introduces "continue-on-error"`) because the misidentified range was a
    **superset** containing the mutation rather than a subset excluding it. Every reordering the
    judge tried widened the scanned range; none narrowed it. **THE GENERALISATION, and the reason
    it is worth a row anyway:** *a scoping bug that currently fails safe is still a scoping bug,
    and "it happened to widen" is a property of the mutations someone thought to try, not of the
    parser.* This repo runs no `actionlint` (row 41's V18), so the YAML is read by hand-rolled
    line scans in three places now. Candidate fixes: parse the step block by indentation rather
    than by the `- name:` token, or assert the located block CONTAINS the `run:` line that
    identified it — a one-line invariant that makes a misattributed range loud. Adjacent to but
    distinct from the design's Declared Residual 2 (`if:`/anchors/flow-style), which names
    constructs this one does not. · ~0.1d · gated on nothing · surfaced iter-133
    (evaluator-found, non-exploitable by measurement, declared out of scope at the landing).~~

53. **[ROUTED — WORLD-SIDE COMPLETE 2026-09-01 (iter-143); TRACKING [`sunholo-data/ailang#941`](https://github.com/sunholo-data/ailang/issues/941). THE ROW'S OWN MECHANISM IS REFUTED AND THE DEFECT IS BIGGER THAN THE ROW SAID.]**
    **AMENDMENT, all first-party at `sunholo-data/ailang@cf4218992`.** The row attributes the
    parse failure to the objection quoting `"GOOS"`/`"GOARCH"` — *"which the reviewer's own JSON
    escaping did not survive."* **Measured on the artifact itself, that is false**: decoding the
    `%q` payload in `.reviewers[]|select(.present==false)|.error` shows those two literals
    **correctly escaped** (`\"GOOS\"` present, bare `"GOOS"` absent). They are the ones that
    survived. The published window is exactly **200 chars** (`reviewer.go:135`, `(raw: %.200q)`)
    and the break is past it — a reconstruction whose escaping is correct in
    `strongest_objection` and whose one bare quote sits in a later field reproduces the recorded
    error **byte-identically** (sha256 `ae8fe247…` both sides, 1-char negative control diffs) with
    `json.SyntaxError.Offset = 256`. Stated honestly: byte-identity of the error does not uniquely
    determine the payload, so that shows such a payload EXISTS and is consistent; the escaping
    fact above stands on the artifact alone. **TWO INDEPENDENT EVIDENCE LOSSES, both one-line:**
    `%.200q` truncates before the failure site, and `out.Err = perr.Error()`
    (`run.go:169`, `agentic_caller.go:205`) drops `*json.SyntaxError.Offset`, which Go's
    `Error()` never prints — so the artifact can say *a* `G` broke it and structurally cannot say
    which. **THE SWEEP IS THE REASON THIS OUTGREW THE ROW.** Across **193 artifacts / 389
    reviewer rows** in both missions (`ailang-world` 88, `ailang` 105): `unreachable` **11**,
    `budget` **15**, `invalid` **6** — and **6 of 6 `invalid` absences carry a recoverable
    `"verdict":"reject"` inside the already-published window**, over three distinct signatures and
    two models (`oc-glm-5-2` ×2, `gemini-3-1-pro` ×4), $0.1537 billed and discarded. Negative
    control: the same extraction over the 15 `budget`/`unreachable` absences yields **0**.
    **AND THIS CLASS IS THE MECHANISM BEHIND THREE KNOWN VACUOUS PASSES** — three of those six
    artifacts synthesise `proceed` with **zero external reviewers present**, and they are the
    exact three the shared skill already names at Gate 2 and attributes solely to
    `--controller-verdict` inflating `presentCount`. That attribution is correct and incomplete:
    each of those reviewers had said **reject**, and this bug is why it was thrown away.
    **REFUTED HYPOTHESIS, recorded so it is not re-chased:** the four `unexpected end of JSON
    input` rows look like the `resp.FinishReason == "length"` guard failing — it is not, that
    guard landed in `885725f06` on 2026-07-17 10:55Z and all four artifacts predate it.
    **ALSO REFUTED — half of the row's own ask is already done:** the row asks for
    `unreachable`/`over-budget`/`unparseable` to be three distinct reasons. They already are
    (`run.go:86-92` carries five, and the reason reaches `.synthesis.absent_reviewers[].reason`).
    What is missing is three distinct **remedies**, and that half belongs to the shared skill, not
    to the quorum binary. **DISPOSITION:** posted to `ailang#941` (comment count asserted 0 → 1)
    with the three fixes in value order — salvage the leading verdict on `ReasonInvalid` so a
    recoverable `reject` can never synthesise to `proceed`; publish the offset and widen the
    window; retry an unparseable response with a stricter instruction rather than a bigger cap.
    Cross-mission note sent to `mission-v1`, which owns both the binary and the skill. **World has
    no further move here** — this row leaves the routable list and lives on `#941`.
    **Original row text follows.** ~~**A REVIEWER CAN BE
    RECORDED `absent` BECAUSE OF WHAT IT SAID, NOT BECAUSE IT COULD NOT BE REACHED — AND THE
    SHARED SKILL'S ABSENT-REVIEWER RULE HAS NO DISPOSITION FOR THAT.** Measured first-party at
    iteration 133's round-2 quorum: `oc-glm-5-2` was listed in `absent_reviewers` with reason
    **`invalid`**, and the artifact's error field reads `reviewer returned non-JSON or malformed
    response: invalid character 'G' after object key:value pair`. It had answered, at cost
    (`$0.0316` banked), and the decoder threw the answer away — **because the objection quoted
    the very string literals under discussion** (`\"GOOS\"` and `\"GOARCH\"` inside a
    `strings.Contains` snippet), which the reviewer's own JSON escaping did not survive. The
    verdict (**reject**, the same objection the other two raised) was recoverable verbatim from
    the error payload and was then confirmed by a solo re-run with a raised cap, so that round
    carried no hole. **WHY IT IS A ROW:** the shared skill's remedy for a non-empty
    `absent_reviewers` is *"re-run each absent reviewer alone with a raised cap"* — written for a
    lane that was unreachable or over budget. A raised cap is irrelevant here and a re-run risks
    the identical parse failure, so the prescribed remedy can loop. Worse, the trigger is
    **content-correlated**: a doc about string-literal pins is exactly the doc that provokes it,
    which means the reviewer most likely to be silenced is the one reviewing the surface where
    quoting literals is unavoidable. **THE GENERALISATION:** *an absence you attribute to the
    channel may have been caused by the message.* Candidate fixes, none of which World can apply
    (the quorum binary lives in `sunholo-data/ailang`): surface the raw payload in
    `absent_reviewers` so a recoverable verdict is not silently discarded; distinguish
    `unreachable` / `over-budget` / `unparseable` as three reasons with three different remedies;
    or retry an unparseable response once with a stricter output instruction rather than with a
    bigger budget. **FILED UPSTREAM 2026-08-27 (iter-133): [`sunholo-data/ailang#941`](https://github.com/sunholo-data/ailang/issues/941)**, with a cross-mission note to V1 on the `mission-control` channel. Routed rather than worked around locally, per the frozen-core rule.
    · ~0.2d (mostly upstream) · gated on nothing · surfaced iter-133 (controller-measured, with
    the recovered verdict and the solo re-run as corroboration).~~

54. **[LANDED 2026-09-02 (iter-148), CI GREEN ON THE MERGE COMMIT [`14036ee`](https://github.com/sunholo-data/ailang-world/commit/14036ee) via PR [#111](https://github.com/sunholo-data/ailang-world/pull/111), SHA-addressed, `present=2 == expected=2`, both `success`; evaluator `sonnet` **86/100**, ONE blocking finding, fixed in-sprint and its fix mutation-proven.** The gate now carries a fleet-comparison arm and a `--driver-fleet-check` isolated mode. **Re-measured at pick time and WORSE than filed: 11 commits and 705 differing lines, not 8/586 — and 3 of 6 tracked driver paths DIFFER — while the gate was green.** The reviewers' verbatim fix (enumerate the UNION, require every path to match) was measured and NOT applied: 6 local / 48 fleet / 42 missing-locally / 0 missing-in-fleet, and the 42 are other missions' plists, env files, fleet-only scripts and 12 testdata fixtures World must not carry — the literal union reds permanently on correctly-absent files, which is the failure mode the doc had already rejected Option C for. **Standing consequence, intended and declared: `./scripts/verify_go.sh` on the RIG is RED until the fleet lands a current driver; CI is unaffected because the arm loud-skips, measured in the ubuntu job log rather than assumed.** Prior row text follows.]** ~~**w-driver-copy-stale-and-the-drift-gate-compares-it-to-itself** · clause-2 · **`launchd` RUNS~~
    THIS REPO'S COPY OF THE FLEET DRIVER, THAT COPY IS 8 COMMITS AND 430 LINES BEHIND THE FLEET, AND
    THE GATE WHOSE NAME IS "DRIVER DRIFT" CANNOT SEE ANY OF IT — BECAUSE IT COMPARES THE COPY TO
    ITSELF.** Measured first-party at iteration 135 while reading the routing environment, not while
    looking for this. `~/Library/LaunchAgents/dev.ailang.mission-world.plist` names
    `<repo>/tools/launchd/mission-control.sh` by absolute path, so the World loop is driven by the
    repo copy and never by the fleet's. That copy is **550 lines against the fleet's 980**, with
    **586 differing lines**, and `git log --since=2026-08-15 -- tools/launchd/mission-control.sh` in
    `sunholo-data/ailang` lists **8** commits it does not carry — the Ollama Cloud fleet enactment,
    the cross-provider role chains, the opus-4.8 drop, the held-pin drift report, and `fc6e42682`,
    which dropped the `:floor` price-pin from the deepseek executor fallback *because `:floor` is
    OpenRouter `provider.sort=price` and the two cheapest endpoints for that model carry NEGATIVE
    health (StreamLake −2, Decart −5)*. **THE PROOF THAT THIS IS LIVE AND NOT A PAPER DIFF:** this
    iteration's own exported `MISSION_EXECUTOR_FALLBACK` is
    `pi:openrouter/deepseek/deepseek-v4-flash-0731:floor` — the dropped id — because
    `~/.config/ailang/mission-world.env` has its override commented out (line 105) and the stale
    driver's `${VAR:-default}` at `:331` therefore governs. The env file's own comment says the
    unfloored id is deliberate, which is a true statement about the FLEET driver and a false one
    about this mission's effective config. **THE GENERALISATION, and it is why this is a row rather
    than a config note:** `scripts/verify_go.sh:205` implements the `D-WORLD-DRIVER-1` drift gate as
    `git status --porcelain -- tools/launchd/`, i.e. **working tree vs THIS repo's HEAD**. A
    stale-but-committed copy is invisible to it *by construction*, and its success line reads
    `✓ driver drift gate: N tracked driver files, working tree matches HEAD` — a true statement that
    a reader parses as "the driver is current". Same family as this mission's standing finding that
    *a control only certifies the axis you varied*, aimed at the one gate that guards frozen core.
    **THE ITEM:** the gate must compare against the FLEET source, not the local copy — e.g. assert
    the copy's content-hash equals the fleet checkout's `HEAD` blob for the same path (with a loud,
    typed refusal when the fleet checkout is absent, so an unreachable source never reads as
    agreement). The driver itself is FLEET-owned and World must NOT edit or absorb it; the gate is
    World's own file and is in scope. · ~0.3d · gated on nothing · queued iter-135
    (controller-first-party, found off-pick).


55. [**LANDED 2026-09-02 (iter-149) — ROW COMPLETE.** PR [#112](https://github.com/sunholo-data/ailang-world/pull/112) → squash [`165b9fd`](https://github.com/sunholo-data/ailang-world/commit/165b9fd), **Gate 3b GREEN on the MERGE commit** (SHA-addressed, `present=2 == expected=2`, both `success`). All three false-reds reproduced first-party before routing, with two known-positive controls in the same call. Quoted `"on":` and flow style now parse **GREEN**; the tab-indent case became a loud **TYPED** failure — the defect there was never the red, it was the message claiming the block was ABSENT. Anything else non-empty reaches `errUnhandledOnForm` rather than a partial key set (`oc-glm-5-2`'s verbatim guard spec). **LINE SCAN, no YAML parser, no new dependency** — `D-WORLD-30`'s rationale applied as PRECEDENT, explicitly not as resolving this row, and no new decision raised. The three wrong-in-the-safe-direction prose claims are corrected in the row-47 doc AND in the gate's own comment; a FOURTH occurrence the judge found sits in a verbatim code exhibit of the pre-row-55 gate and is flagged **HISTORICAL — DO NOT UPDATE** rather than rewritten. Quorum: 2 FULL-STRENGTH rounds, both BLOCKED, closed under the narrow-refinement carve-out on all three reviewers' verbatim `proposed_fix`; every objection was a premise objection and was MEASURED, not forwarded. Planner `opus` refuted **six** doc premises; executor `codex:gpt-5.6-sol` 11/11 ACs and 5/5 mutants with enumerated blast radii; evaluator `sonnet` **95/100, zero blocking**, 11 adversarial fixtures unable to produce a silent false green. Its top non-blocking finding was real and closed in-PR: this row's HEADLINE message was pinned by nothing, and the two new arms are proven non-vacuous by MUT-F. **Prior row text follows.**]~~ **w-dispatch-lever-parser-false-reds-on-valid-yaml** · clause-2 · **THE ROW-47 GATE NEVER~~
    SILENTLY PASSES A LEVER-LESS WORKFLOW — EVERY ATTACK AIMED AT THAT FAILED — BUT ITS LINE-SCAN
    PARSER REDS ON THREE FORMS OF *VALID, LEVER-DECLARING* YAML, AND ONE OF THEM IS THE STANDARD
    REMEDY FOR A FAMOUS GITHUB ACTIONS FOOTGUN.** Filed by the `sonnet` evaluator of P47 as
    non-blocking findings 1–3, and filed as a row rather than fixed inline because closing them
    needs a decision this item does not own — whether the gate adopts a structural YAML parser
    instead of a line scan. Two SIBLING findings from the same evaluation WERE fixed in-PR at
    [`eb215c3`](https://github.com/sunholo-data/ailang-world/commit/eb215c3) (an inline `#` comment
    misread as a scalar value; the doc claiming a residual its own code comment never stated), so
    what remains here is the residue that is genuinely deferred rather than the whole finding set.
    **(a) The quoted `"on":` form false-reds.** `l == "on:"` requires exact unquoted byte equality,
    so `"on":` reports `instrument failure: … has no top-level \`on:\` trigger block` on a file that
    DOES declare the lever. Quoting `on` is the standard fix for YAML 1.1 parsing bare `on` as the
    boolean `true` — so this gate breaks CI the day anyone applies the recommended remedy.
    **(b) Flow style false-reds:** `on: {push: …, workflow_dispatch: }` fails identically.
    **(c) A tab-indented first trigger line silently empties the trigger SET** — `TrimLeft(l, " ")`
    strips spaces only, so a tab-indented `push:` computes `lead=0 == onLead`, the loop reads the
    block as already exited, and the message becomes `triggers=[] lack workflow_dispatch`, which
    misreports total absence rather than a parse limitation. **All three fail LOUD, which is the
    accepted direction** — the evaluator's attacks aimed at a silent false green (hidden dotfile,
    nested subdirectory, case/extension variants, all nine canonical mutations, and the
    anti-vacuity floor's own precondition) **all failed to produce one**. So this is a
    *brittleness* row, not a soundness row, and it must not be re-scoped as the latter.
    **Carries with it, measured in the same evaluation:** the doc's residual-3 claims are
    themselves wrong in the safe direction — Go's `filepath.Glob` does NOT skip dotfiles the way
    POSIX shell globbing does (a hidden `.hidden.yml` IS caught), and a nested subdirectory is not
    invisible but a loud `t.Fatal` on `os.ReadFile` ("is a directory"). Fix those residual
    sentences in the same edit. **And the scalar arm emits TWO messages, not one**, so the
    Decision's "never cascades" phrasing is an over-claim the planner's D5 already flagged and
    this row inherits. · ~0.2d · gated on nothing · surfaced iter-136.
56. **[LANDED 2026-09-02 — ROW COMPLETE, AND TAGGED LATE BY ITERATION 151 BECAUSE THE ITERATION
    THAT DID THE WORK DIED BEFORE GATE 4 AND LEFT NO TRACE ANYWHERE.** PR
    [#113](https://github.com/sunholo-data/ailang-world/pull/113) → squash
    [`725ad5a`](https://github.com/sunholo-data/ailang-world/commit/725ad5a), branch
    `sprint/w-canary-fence-blind-to-a-skipped-canary`, merged `2026-09-02T19:26:21Z`.
    **Gate 3b verified retroactively to the normal standard, SHA-addressed on the MERGE commit:**
    `checks=2 == expected=2`, `ailang-code verify gate` **success** and `go host build + test gate`
    **success** (2 and not 3 because the third job, `launchd-drivers`, arrived later with
    `e308577`). **The claimed deliverable was verified to have actually landed rather than
    inferred from the commit subject** (inherited-claim discipline — nobody had reviewed this work
    since the agent that wrote it stopped existing): the needle is live at
    `host/verifygate/toolchain_pin_gate_test.go:510`, refusing with
    `t.Skip/t.Skipf/t.SkipNow call count=%d, want 0 (a skipped canary asserts nothing)`.
    **Why this row was invisible to every check the loop runs:** the orphaned iteration wrote
    **zero** charter rows, **zero** log entries, **zero** STATUS stamps and **zero** dashboard
    lines — measured, `grep -cE '#113\b'` = **0** across all four mission documents, against a
    known-positive control of `#114` = **2** in the charter and **1** in the log. Its only
    surviving traces were the merged PR itself and a `prunable` `/private/tmp/wt-row56` entry in
    `git worktree list` (the worktree directory having been destroyed by the same reboot described
    in the Repo Profile). Iteration 150 read `725ad5a` in the very next fire — and used it **only
    as a CI negative control**, never asking what it was. The heartbeat is the reason the slot is
    datable at all: that fire stamped `gate-0` at `2026-09-02T18:27:39Z` and `gate-1` at
    `18:28:00Z` and then nothing, while demonstrably running through Gate 3 and merging at
    `19:26:21Z`. **Its own cause is UNRECOVERABLE and this is recorded as a limit, not a
    verdict** — `/tmp/ailang-mission-world.log` was wiped by the reboot five hours later, so the
    driver's `HARD TIMEOUT` / `STALL:` lines for that fire no longer exist. Prior row text
    follows.** ~~**w-canary-fence-blind-to-a-skipped-canary** · clause-2 · **THE ROW-49 FENCE PROVES THE CANARY
    *CONTAINS* ITS ASSERTION AND CANNOT SEE THAT THE ASSERTION NEVER *RUNS* — `t.Skip()` AS THE
    FIRST STATEMENT LEAVES THE AST BYTE-IDENTICAL, SO THE FENCE STAYS GREEN OVER A CANARY THAT
    ASSERTS NOTHING.** Found by the `sonnet` evaluator of P49 as adversarial arm `A10`, and
    **reproduced first-party by the controller in the same iteration**: mutant LANDED (sha256
    `a23cfa79…` → `06688802…`), BUILDS (`go vet ./host/store/` rc=0), the canary itself reports
    `--- SKIP` (`RUN=1 PASS=0 SKIP=1 FAIL=0`) while `TestCanaryDeclaresPositiveArmOnly` returns
    `rc=0 RUN=1 PASS=1`; restore byte-identical, pristine control green after. **This is filed as a
    row rather than a defect in P49**: row 49's design doc declares this exact class up front
    (*"The shape is not the behaviour. The AST check proves the canary contains the assertion shape;
    it cannot prove the assertion fires"*), and the evaluator confirmed the disclosure preceded the
    finding. It is recorded because a disclosed residual is still a live hole in a gate the mission
    relies on, and because it is **cheap to close in the shape P49 already built**: the canary
    currently contains **0** `t.Skip` occurrences (known-positive control in the same call:
    **3** `t.Fatalf`), so a zero-needle on `t.Skip`/`t.SkipNow` — plus the sibling runtime-neutering
    forms the judge named (`t.SkipNow()`, an early `return`, a build tag) — folds into the SAME AST
    pass that already walks `TestToolchainCanary`'s body. **Scope discipline the row must keep:**
    this closes *statically visible* reachability only. It does NOT make the fence a proof that the
    assertion fires on a miscompiling toolchain — that lane is the nested repro module plus
    `run.sh`'s `saw_bad` floor, which is darwin-only and attended while `ci.yml:172`
    `continue-on-error: true` stands (row 44). Do not let this row grow into that one.
    · ~0.3d · gated on nothing · surfaced iter-139.~~

57. **w-approvals-spine-prints-a-green-no-pending-under-the-row-it-just-listed** · clause-2 ·
    **EVERY MISSION-LOOP-AUTHORED ROW ON THE DECISION SPINE IS TYPED `notification`, NOT
    `approval_request` — SO `ailang coordinator pending` LISTS THE ASK AND THEN PRINTS
    `✓ No pending approval requests` DIRECTLY BENEATH IT.** Measured iter-139, first-party, while
    posting `D-WORLD-28` to the spine. **The split is by AUTHOR CLASS and it is total:** of 12 rows
    in the `approvals` inbox, all **6** authored by `coordinator` carry `message_type=approval_request`
    and all **6** authored by a mission loop (`mission-v1` ×3, `mission-motoko` ×1, `mission-world`
    ×1, `audit-r1` ×1) carry `notification`. So this is not one controller's typo; it is what the
    shared skill's Gate-5 snippet produces for every mission on the rig. **The prescribed flag is
    accepted and silently ignored:** `ailang messages send … --type approval_request` exits **0**,
    emits no warning, and stores `notification`. **A flag-ordering hypothesis was tested and
    REFUTED** — Go's `flag` package stops at the first non-flag argument, so the natural theory is
    that the skill's body-positional-first form leaves the flags unparsed; two arms sent to a
    throwaway inbox (`-type` AFTER the body vs `-type` BEFORE it) both stored `notification`, while
    `-title` and `-from` landed correctly in both. The discriminator collapsed, and here the
    collapse IS the finding: the control is that `approval_request` is representable in this store
    (6 rows prove it), so the value is not being rejected — it is not being applied on this send
    path. **SCOPE, STATED HONESTLY — the ask is NOT lost.** `ailang coordinator pending` unions the
    unread `approvals` inbox and DID surface the row (verified in the same call). What is wrong is
    the second line it prints: the typed sub-query finds zero `approval_request` rows and reports
    `✓ No pending approval requests` — a green checkmark, under the row it just listed. That is the
    vacuous-pass shape this mission keeps closing, aimed at the one surface whose whole job is to
    say what is waiting on a human, and the green line is exactly the kind of sentence that gets
    quoted onward. Whether the Discord "🔔 Approval needed" push filters on the type is **NOT
    established** and must be measured before this row claims it. **Fleet-scoped, not World-owned:**
    the fix is in `ailang messages send` (or in the shared skill's snippet), so it routes to
    `sunholo-data/ailang` as an issue plus a cross-mission note, per the frozen-core rule — this row
    tracks World's side of it. · ~0.2d · gated on nothing · surfaced iter-139.


58. **[AMENDED 2026-08-31 (iter-141) — THE HEADLINE IS WRONG AS FILED: THE GATE IS **FLAKY**, NOT DETERMINISTICALLY RED, AND THAT CHANGES THE FIX.** Four runs of `scripts/verify_go.sh` on this rig within about 90 minutes, on trees whose only difference is one test file: iter-140 **rc=1 with FOUR** failing tests; this iteration's Gate-2 base run **rc=1 with ONE** (`TestHandlerTimeoutKillsTheWholeProcessGroup`, and note it failed on a DIFFERENT mechanism — its own non-vacuity markers reporting `exec_started=false forked=false`, not the archived-interpreter `--version` timeout row 58 names); the `opus` planner's base run **rc=0 with an EMPTY failing set**; the controller's post-change run **rc=0 with an EMPTY failing set**. Two red, two green, three different failing sets. **The planner refuted a fact the controller had handed it under a VERIFIED-BY-ME label** — which is the loop working (Gate-2 rule (d)) — and the correction matters because row 58 as filed prescribes repairing every sprint AC around a *deterministic* base red. A flaky gate needs a different remedy: an AC written as a SET comparison against a recorded base set is still right, but the row must ALSO say that the base set is not stable across runs, so the honest criterion is `S_after \ S_known_flaky == ∅` with `host/verifygate` `ok` in both legs, and a red must be re-run unloaded before it is attributed to anything. The mechanism analysis (macOS first-exec provenance assessment, 96 ms pinned vs 1336 ms cold, plus a per-invocation cleanup over a 513 MB Observatory DB) STANDS as a real rig cost and is the likeliest driver of the variance; what does not stand is *"rc=1 on pristine dev"* stated as a property. **The instrument-failure floor the row asks for is now MORE clearly the right deliverable, not less** — a timing-sensitive test that reds under load and passes alone is exactly a test that should name the timeout an ENVIRONMENT failure. Prior row text follows.** ~~**w-verify-go-is-red-at-pristine-base-on-the-rig-while-ci-is-green** · clause-2 · **THE LOOP'S OWN LOCAL GATE `scripts/verify_go.sh` IS `rc=1` ON UNTOUCHED `dev`, SO EVERY SPRINT'S `verify_go.sh rc=0` ACCEPTANCE CRITERION IS BROKEN AS WRITTEN — IT CAN ONLY BE WAVED THROUGH OR BLAMED ON THE SPRINT.** Measured iter-140 on a pristine worktree at `dev` = `9c0ad0b`, in **TWO arms** — once under concurrent load and once alone, because the controller had itself introduced a load confound by running two full suites plus a designer session in parallel. **The confound is REFUTED: the failing set is IDENTICAL in both arms**, four tests, `diff` silent — `TestEpisodeLiveReplayThreeArmsAndEvidence` and `TestHandlerTimeoutKillsTheWholeProcessGroup` (`host/broker`), `TestF1PinnedInterpreterHashMismatchRefusedBeforeExec` (`host/capsule`), `TestFixtureEpisodeReplaysBitForBit` (`host/replay`); `host/verifygate` is `ok` in both. **Three of the four share ONE mechanism** — `cannot obtain --version from archived interpreter … timed out after 10s`. **Mechanism isolated, not guessed:** the same binary answers `--version` in **96 ms** at its pinned path, **1336 ms** on FIRST exec of a freshly-copied instance at a new temp path, and **372 ms** on a cached second exec — macOS first-exec code-signing/provenance assessment (`com.apple.provenance` xattr present; `spctl -a -t exec` rejects) — on top of a per-invocation Observatory retention cleanup over a **513 MB** DB that every `ailang` invocation on this rig now pays. **This is RIG-LOCAL, not a repo defect, and the discriminator is the fleet control group:** dev CI is **GREEN on this same commit**, 2/2 checks, on Linux runners with neither Gatekeeper nor a 513 MB observatory. **Why it is a row and not a shrug:** `verify_go.sh` is the gate this mission cites in every sprint's acceptance criteria and every STATUS stamp, and iter-139 recorded it `rc=0` roughly five hours earlier — so the regression is in the rig's state, is recent, and will silently convert into either a false *"the sprint broke it"* or a habit of ignoring the gate. **THE ITEM:** decide the disposition — raise the archived-interpreter `--version` deadline above the first-exec assessment cost, pre-warm or ad-hoc-sign the archived copy, or bound the observatory cleanup — and, whichever lands, give the three replay-family tests an instrument-failure floor that names the timeout as an **ENVIRONMENT** failure rather than a replay failure, so a rig cost can never again wear a correctness defect's clothes. · ~0.3d · gated on nothing · surfaced iter-140 (controller-measured).

59. **w-static-grep-cannot-prove-an-assertion-is-live** · clause-2 · **THE ROW-51 DESIGN PROVES ITS
    OWN SET-DIFFERENCE "LOAD-BEARING" WITH A `grep -c`, AND A `grep -c` CANNOT TELL AN ASSERTION
    THAT RUNS FROM ONE THAT IS COMPILED AND NEVER REACHED.** Found by the `sonnet` evaluator of row
    51 as arm `E5`, and **REPRODUCED FIRST-PARTY BY THE CONTROLLER** rather than inherited — a
    NON-BLOCKING tag is a judge's severity opinion, not a measurement. Wrap the entire new
    assertion block in `if false { … }` and land arm A1 (a fabricated 7th coupled-site row in the
    `scripts/verify_ail.sh` home only): the mutant LANDS (`d9bd7a81…` → `ad86884d…`), the block row
    count MOVES 6→7, `gofmt` is clean, **`go vet ./host/verifygate/` rc=0**, and the test returns
    **rc=0 `--- PASS`** — while **AC2's own instrument, `grep -c 'disagree on their site set'`,
    still reads 3**, exactly its passing value. Restored byte-identical, pristine control green
    either side, porcelain 0. **THE SHIPPED CODE IS CORRECT — THE ACCEPTANCE CRITERION IS NOT.**
    Row 51's assertions all fire; what fails is the *proof obligation* attached to them, so this is
    a row rather than a blocking finding. **THIS IS ROW 49'S DEFECT ONE LAYER UP, AND THAT IS WHY
    IT IS WORTH ITS OWN ROW.** Row 49 replaced `strings.Count(src,"stateRoot") >= 2` with an AST
    pass precisely because a token count cannot see shape-gutting — and the judge then measured that
    **22 of 23 mutations had walked straight through the old count**. Row 51 then wrote its own
    load-bearing proof as a token count. The loop closed the defect in the *product* and
    reintroduced it in the *acceptance criteria*, which is this mission's own recurring *guard the
    helper, miss the call site* shape aimed at its rulebook. **THE ITEM:** decide the general rule
    — a "this assertion is load-bearing" criterion must be discharged by a MUTATION that makes the
    test red, never by a static count of the assertion's own text — and bind it where it can
    actually fire. Note the cheap, already-built shape: row 49's AST pass in `host/verifygate` can
    ask whether a `t.Errorf` is reachable at function scope rather than whether its message string
    is present, and `N1`-style deletion arms do not cover the wrapping variant. **ALSO IN SCOPE,
    from the same evaluation and measured in the same session:** row 51's Declared Precondition
    overstates the §S8 backtick requirement — arm `E2` removed the backticks from both homes and the
    test still PASSED, because `([^|]+)` plus the trim tolerates it. That error is in the SAFE
    direction (the doc claims a stricter precondition than the regex needs) and is a doc-accuracy
    nit, but it is the same class: a stated invariant nobody made red. **THE GENERALISATION:** *a
    criterion that greps for an assertion measures that somebody typed it; only a mutation measures
    that anybody runs it.* · ~0.2d · gated on nothing · surfaced iter-141 (evaluator-found,
    controller-reproduced, declared out of scope by the row it was found in).

60. **w-p1-needle-reds-on-a-semantically-inert-rename** · clause-2 · **THE ROW-48 NEEDLE SET PINS
    AN IDENTIFIER, SO A CONSISTENT RENAME THAT CHANGES NOTHING FIRES A LOUD RED.** Found by the
    `sonnet` evaluator of row 48. `P1d` asserts `strings.Count(gate, 'go_version_ge "$ACTIVE_GO"
    "$ROOT_FLOOR"') == 1`, so renaming `ROOT_FLOOR` → `ROOT_FLOOR_` **consistently throughout the
    P1 block** — semantically inert, `bash -n` clean, the gate behaving identically — drops the
    count to 0 and reds CI. **The direction matters and is why this is ~0.1d and not urgent: this
    fails CLOSED**, which is the safe half of the row-49/row-59 axis; a false red costs an
    iteration, a false green costs the gate. **What it does establish is that `M17` is a NARROWER
    green control than it reads.** M17 rewords a *comment word*; it does not exercise an
    identifier rename, so the doc's "the set is not any-edit-reds" claim is supported for comments
    and unmeasured for renames. **THE ITEM:** either add a rename arm to the green-control set and
    relax `P1d` to bind the operand ORDER without pinning the variable's spelling (e.g. assert the
    call's two arguments appear in the `ACTIVE_GO`-then-floor order, whatever they are named), or
    record deliberately that the identifier is part of the contract and say so in the fence. Do
    not do both. · ~0.1d · gated on nothing · surfaced iter-145 (evaluator-found; brittleness
    established by construction from the test's own literal and confirmed by the judge's run).

61. **w-p1-gate-fails-open-on-one-inserted-line-past-every-arm-in-the-sprint** · clause-2 · **A
    SINGLE DEAD LINE DEFEATS ALL SIX STATIC NEEDLES *AND* THE RUNTIME ARM, AND THE GATE THEN
    PRINTS A SELF-CONTRADICTING SUCCESS.** Found by the `sonnet` evaluator of row 48 and
    **REPRODUCED FIRST-PARTY BY THE CONTROLLER** — a NON-BLOCKING tag is a judge's severity
    opinion, not a measurement. Insert one line, `floor_rc=0`, immediately after
    `go_version_ge "$ACTIVE_GO" "$ROOT_FLOOR"; floor_rc=$?`. Measured: mutant LANDED
    (`grep -c '^floor_rc=0$'` **0 → 1**), `bash -n` rc=0, `go vet` rc=0, and then **the static
    needle set returns `rc=0 RUN=1 PASS=1`** — every one of P1a–P1f is satisfied, because every
    byte they assert about is still present and in order — **and the M9 runtime arm returns `rc=0`
    with the race leg REACHED and 2 races**, i.e. the fail-closed gate `D-WORLD-28` mandates is
    fully open under a below-floor toolchain. **The tell is the strongest part:** the gate prints
    `✓ toolchain floor gate: go1.26.4 >= root module floor go1.26.6` — an arithmetically FALSE
    assertion published as a success line, so the log actively argues for its own correctness
    while the control it protects runs unvetted. Restored byte-identical (`7abf30b1…`), pristine
    control green after. **THIS IS INSIDE THE RESIDUAL ROW 48's DOC ALREADY DECLARES** (§9.1: a
    static scan cannot see dataflow) — which is why it did not block the merge — **but the doc
    describes that residual as reachable by REWRITING existing branches (`case "$floor_rc" in` →
    `case "0" in`, or inverting the awk verdicts), and this is reachable by a one-line
    INSERTION**, which is a strictly larger and much likelier accident. **THE ITEM:** decide the
    disposition — the cheapest honest fix is to stop laundering the comparator's verdict through a
    reassignable variable (branch directly on the call, or make the assignment and its `case` a
    single unbreakable expression); the alternative is to keep the shape and add a runtime
    self-check that the printed `✓` line's own inequality actually holds. Whichever lands, the
    residual text in the doc must be widened from "rewrite" to "any dataflow break, including an
    insertion". · ~0.2d · gated on nothing · surfaced iter-145 (evaluator-found,
    controller-reproduced, declared in-residual by the row it was found in).

62. **w-flag-scan-false-positives-on-an-explicit-false-and-on-benign-script-text** · clause-2 ·
    **THE ROW-52 SPRINT FIXED WHERE THE STEP BLOCK IS; IT DID NOT TOUCH WHAT THE SCAN LOOKS FOR
    INSIDE IT, AND THAT LOOP IS A BARE `strings.Contains(line, "continue-on-error")` — SO IT REDS
    ON AN EXPLICIT `continue-on-error: false` AND ON A COMMENT OR `echo` THAT MERELY NAMES THE
    FLAG.** Found by the `sonnet` evaluator at the row-52 landing and **reproduced first-party
    before filing**, in the judge's own worktree, mutants asserted LANDED by sha256, effect
    asserted against `ruby -ryaml`, `go vet` rc=0 read before every verdict, restored
    byte-identical, pristine control green after: **(a)** `continue-on-error: false` inserted on
    the guarded step — a legitimate, explicit opt-OUT that changes nothing about failure
    swallowing — gives **rc=1 at `ci.yml:175`** with the YAML view confirming
    `step[6].continue-on-error = false`; **(b)** a benign `echo "note: never add
    continue-on-error to this step"` inside that step's own `run:` scalar gives **rc=1 at
    `ci.yml:177`** with `step[6].continue-on-error = nil`, i.e. a red over script TEXT with no
    flag present at all. **The discriminating control, which is what makes this a scan defect
    rather than a repo-wide text ban:** the identical `echo` placed in an UNRELATED step stays
    **rc=0**. **WHY IT IS A ROW AND NOT A BLOCKER:** it is PRE-EXISTING — inherited unchanged
    from the row-44 test, untouched by row 52, and outside that sprint's stated one-function
    scope — and it fails CLOSED, so no refusal can be swallowed by it. It is filed because it is
    **declared nowhere**: neither the row-44 doc's residuals nor row 52's name it, so the first
    maintainer who legitimately writes `continue-on-error: false` meets an unexplained red.
    Candidate fixes: match the KEY at the step's own key indentation rather than the substring
    anywhere on the line, and treat a `false` value as compliant; or, at minimum, declare both
    shapes as residuals with the reason. · ~0.1d · gated on nothing · surfaced iter-147
    (evaluator-found, controller-reproduced).

63. **w-locator-derivation-refusals-are-unpinned-and-undeclared** · clause-2 · **THE ROW-52
    LOCATOR HAS FIVE LOUD REFUSAL BRANCHES AND THE DRILL PINS THREE OF THEM; THE TWO DERIVATION
    REFUSALS — `anchor < 0` AND `stepCol < 0` — HAVE NO KILLER ARM AND NO DECLARED-UNREACHABILITY
    NOTE.** Found by the `sonnet` evaluator, which neutered each new branch in turn and reported
    that these two survive all 17 executor arms plus its own adversarial set, adding that both
    are *plausibly* unreachable. **Reproduced first-party, and the finding came back SMALLER than
    filed on one half:** `anchor < 0` is REACHABLE and was fired — renaming every `steps:` key
    above the identifying line gives **rc=1 `instrument failure: could not locate a steps: anchor
    above the miscompile identifying line in ci.yml`** (mutant LANDED by sha256, `go vet` rc=0,
    restored byte-identical, pristine control rc=0 after). So the honest statement is not "no
    killer exists" but "no killer is COMMITTED", and the mutant that fires it is not valid
    Actions YAML, which is exactly the kind of caveat that belongs in the doc rather than in a
    reviewer's head. `stepCol < 0` remains unfired. **THE GENERALISATION, which is why this
    outranks its own size:** this mission's standing rule is that a guard is not a gate until
    something reds when you remove it, and a branch whose unreachability is *asserted* rather
    than *shown* is the same claim class the row-52 doc criticises elsewhere in its own text.
    **THE ITEM:** add an arm for each derivation refusal, or declare each unreachable IN THE CODE
    with the measurement that supports it — matching the MUT-H precedent, which declares Invariant
    A unreachable with a stated reason rather than leaving it silent. · ~0.1d · gated on nothing ·
    surfaced iter-147 (evaluator-found, controller-reproduced and partially refuted).


64. **w-fleet-residual-net-shares-phase-1-pathspec** · clause-2 · **THE ROW-54 ARM REPORTS EVERY
    UNCLASSIFIED FLEET-ONLY DRIVER FILE — BUT ONLY INSIDE THE TWO PREFIXES IT ALREADY LOOKS AT, SO A
    FLEET DRIVER FILE ADDED ANYWHERE ELSE IS NOT MERELY UNCERTIFIED, IT IS UNCOUNTED AND UNMENTIONED.**
    Found by the iteration-148 evaluator with an ADDITION-shaped mutant and confirmed by the controller.
    Phase 3's loud residual net enumerates `git ls-tree -r --name-only HEAD -- tools/launchd
    scripts/mission_decisions.sh` against the FLEET — the same fixed pathspec Phase 1 uses — so a
    fleet-only driver file outside `tools/launchd/` and outside the literal
    `scripts/mission_decisions.sh` is invisible to all three phases. The judge demonstrated it with a
    sibling file (`scripts/mission_decisions_v2.sh`): the arm returned **rc=0**, printed `tracked copy
    is current`, and mentioned the file **zero** times — not even as unclassified. **This is a genuine
    instance of the enumerator-blind-spot class the row-54 work exists to close, one level out**: the
    arm's coverage is a property of its pathspec, and the success line's disclaimer ("untracked fleet
    additions not certified") does not distinguish the LOUD category from this SILENT one, so a reader
    is told less than the line implies. **LATENT, NOT LIVE, and measured so rather than assumed:** the
    real fleet at `722e19c7` carries no executable driver file outside those prefixes (control: the
    same enumeration finds the six it does carry). Filed as a row rather than folded into row 54
    because a pre-existing scope question is a queue row, not a sprint widening. **THE ITEM:** either
    widen the residual net's pathspec and say what the new boundary is, or narrow the success line so
    it claims only what the pathspec can see — and pin whichever you choose with an addition-shaped
    mutant, since removal proves the check FIRES and only an addition proves it LOOKS. · ~0.2d ·
    gated on nothing · surfaced iter-148 (evaluator-found, controller-confirmed, measured latent).

65. **w-go-build-is-not-a-compile-fence-for-a-test-file** · clause-2 · **EVERY "THE MUTANT BUILDS" ASSERTION IN THIS REPO THAT USES `go build ./...` IS VACUOUS WHENEVER THE MUTANT LIVES IN A `_test.go` — AND THE MUTANTS ALWAYS DO, BECAUSE THIS MISSION'S GATES ARE TESTS.** Found by the `opus` planner of row 55 as its refutation R-1, **reproduced first-party by the controller** rather than inherited. `go build` does not compile test files, so a mutation drill that fences on it cannot tell "the mutant is sound" from "the mutant does not compile" — which is the exact third-fact-wearing-the-same-exit-code class the mutation discipline exists to close. **Measured at `234d9da`, mutant LANDED by sha256 and restored byte-identical either side:** appending `func zzBroken() { var x int = "not an int"; _ = x }` to `host/verifygate/dispatch_lever_gate_test.go` leaves `go build ./...` at **rc=0**, while `go test -count=1 -run '^$' ./host/verifygate/` reds at **rc=1** naming the type error, and `go vet ./host/verifygate/` also reds at **rc=1**. So two working fences exist and the one in use is the broken one. **THE ITEM:** sweep every place this repo asserts a mutant builds and replace the fence — `git grep -l` finds the pattern in `CLAUDE.md`'s verify-gate prose and in a number of implemented sprint plans; row 55's own plan already uses the `-run '^$'` form, so the corrected shape is in-tree to copy. **Scope, stated honestly:** `verify_go.sh` itself is NOT affected — it also runs `go test ./... -count=1`, so CI has no hole. This is an **AC-instrument** defect, not a CI defect, and it must not be re-scoped as the latter. **ALSO PROPOSED UPSTREAM:** the shared `mission-control` SKILL.md prescribes exactly this fence ("assert `go build ./...` … rc=0 on the mutated tree"), so the gap is fleet-wide and not World's alone; World cannot edit that file and has proposed it to V1/Mark. **Instance 1** — pre-registered here so the second is recognisable rather than rediscovered. · ~0.1d · gated on nothing · surfaced iter-149 (planner-found, controller-reproduced, declared out of scope by the row it was found in).

66. **w-flow-key-quote-trim-is-uncovered** · clause-2 · **A QUOTED KEY INSIDE A FLOW MAPPING — `on: {"workflow_dispatch": }` — IS VALID YAML DECLARING THE LEVER, AND NOTHING PINS THE TRIM THAT MAKES IT MATCH.** Found by the `sonnet` evaluator of row 55 as non-blocking finding 2. Removing `strings.Trim(…, "'\"")` at `dispatch_lever_gate_test.go:162` leaves the full package **rc=0** (compile fence rc=0, `AILANG_BIN`-set package run rc=0), so the shape would go undetected as a declared lever. **The failure is a false RED, not a false green** — consistent with the row's declared fail-loud bias — and the shape sits outside row 55's three declared shapes, which is why it is a row rather than a blocking finding. Note the symmetry worth recording: this is the same defect class row 55 just closed (a valid lever-declaring form the scan cannot see), one nesting level deeper, and row 55's own arm set does not reach it. **THE ITEM:** add the arm, or declare the shape unsupported **in the code** and let it reach `errUnhandledOnForm` loudly — either is acceptable, an undeclared silent brittleness is not. · ~0.1d · gated on nothing · surfaced iter-149 (evaluator-found, out of the found row's declared scope).


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
67. **[LANDED 2026-09-03 (iter-150)] w-ci-job-enumeration-conflates-go-pin-count** · clause-2 · **`dev` WENT RED AT A MERGE OF TWO INDIVIDUALLY-FINE BRANCHES, AND THE GATE WAS RIGHT — IT JUST COULD NOT SAY WHICH OF THE TWO FACTS IT OBJECTED TO.** `e308577` added a third CI job (`launchd-drivers`, bash 3.2 on macos-latest, **no Go in it**); `725ad5a` carried a gate whose ONE hand-maintained `wantJobs` constant answered both *which jobs exist* and *how many `GOTOOLCHAIN`/`go-version`/`setup-go` pins there must be*. Neither branch is wrong; the conflation is. **Negative control:** parent `725ad5a` CI `success` on both jobs, and `e308577` carries **zero** check-runs — never verified on its own, so the merge was the first run that could see the pair. **The obvious one-line fix is REFUTED by measurement:** widening `wantJobs` to three entries alone is **rc=1 on the pristine `ci.yml`** (2 pins vs 3 jobs) — it demands a Go pin on a job that must not have one. **SHIPPED:** two hand-maintained lists (`wantJobs`, `wantGoPinnedJobs`), a new job must be classified in BOTH or the gate reds, and a job classified OUT is asserted to carry **zero** Go pins; jobs parsed WITH their bodies so pins are attributed **per job**, closing a fail-open the count-only gate always had (two `GOTOOLCHAIN` in one Go job and none in the other satisfied `count==2` exactly as well as one each); plus a whole-file-vs-per-job-sum arm for a pin outside every job. **7 mutants, all RED with the correct assertion named, each LANDED by sha256, `ci.yml` restored byte-identical; the revert-just-the-hunk arm (green at base) leaves M5 — the misattribution case — SURVIVING rc=0, so the per-job hunk is its sole killer.** PR [#114](https://github.com/sunholo-data/ailang-world/pull/114) → squash [`a7b58dd`](https://github.com/sunholo-data/ailang-world/commit/a7b58dd). · surfaced + closed iter-150 (dev RED, outranked the queue).

68. **w-driver-pin-named-world-points-at-the-wrong-repo** · clause-2 · **A PIN WORKTREE NAMED FOR THIS MISSION HOLDS THE *OTHER* MISSION'S REPO, AND NOTHING WOULD TELL YOU.** Measured iter-150 while answering `mission-v1`'s standing ask about shared clones: `~/.ailang-driver-pin/world` has `git-common-dir` = `~/dev/sunholo-data/ailang/.git` and remote `https://github.com/sunholo-data/ailang.git` (NOT `ailang-world.git`) — it contains `design_docs/v1-mission.md` and **no** `design_docs/world-mission.md` (`ls` rc=2), HEAD `da96b98a5` dated 2026-08-26. It is inert **today** only because `dev.ailang.mission-world.plist` runs World's own driver at `~/dev/sunholo-data/ailang-world/tools/launchd/mission-control.sh`, which is **not** pin-rooted (`grep -c pin-root` = 0; `REPO` derives from the script location). So World is genuinely UNPINNED — which is the first-party answer V1 asked for at iter-311 and never got. **THE ITEM:** the residue is a loaded gun for the day World's plist is pointed at a pin-rooted driver, because `MISSION_WORKDIR` would then resolve to a checkout in which this mission's charter does not exist; either delete the worktree or make the pin assert that `$MISSION_DOC` resolves inside it. Fleet-owned decision — World files it and hands it over, does not fix it. · ~0.1d · gated on nothing · surfaced iter-150.

69. **w-heartbeat-script-absent-so-every-gate-stamp-is-a-no-op-here** · clause-2 · **THE SKILL'S MANDATED FIRST ACTION AT EVERY GATE — `bash tools/launchd/mission-heartbeat.sh stamp gate-N` — IS `No such file or directory` IN THIS REPO.** Measured iter-150 on the very first command of the iteration: the file exists at `~/dev/sunholo-data/ailang/tools/launchd/mission-heartbeat.sh` and **not** in World's `tools/launchd/` (which holds `mission-control.sh`, `derive-planner-lane.sh`, the two test suites and the plist template). `~/.ailang/state/mission-world-heartbeat` DOES accumulate stamps, so previous controllers reached for the V1 absolute path by hand — the same undocumented workaround `mission_directives.sh` needs here, which is the pattern rather than the incident. **Why it matters beyond tidiness:** the per-gate stamps are the durable attribution contract for standing rule 7 (a slot that dies holding a background task exits `rc=0` and neither watchdog fires), so for World that contract has been carried by controller discipline rather than by the driver. Frozen core (`D-WORLD-DRIVER-1`) → the port must land as a **fleet-authored commit**; World must not vendor a copy. **THE ITEM:** hand it to the fleet with this measurement, and until it lands, say in each record that the stamps were made by absolute path. · ~0.1d · gated on the fleet · surfaced iter-150.

70. **w-gate-1-cannot-see-an-unverified-commit-that-is-not-head** · clause-2 · **A DIRECT-TO-DEV COMMIT WITH ZERO CHECK-RUNS IS INVISIBLE TO EVERY GATE IN THIS LOOP, BECAUSE GATE 1 READS ONLY `origin/dev`'s HEAD.** Measured iter-150: `e308577` — the commit that added a CI job and thereby armed this iteration's red — has `checks=0` (control: its parent `d75b9c3` and the merge both resolve check sets fine, and the endpoint answers abbreviated SHAs correctly). The origin skill has a rule for a HEAD whose check count is a true zero (fire `workflow_dispatch`, do not record a health verdict); it has **nothing** for a zero on a commit that HEAD has since moved past, and that is the likelier shape here because World's dev takes direct fleet pushes between fires. The red surfaced two commits later wearing a merge's clothes, which cost the attribution work in row 67 rather than being caught where it happened. **THE ITEM:** at Gate 1, after the HEAD check-set read, enumerate `dev@{previous fire}..origin/dev` and assert each commit either has a check set or is explained (intra-push, not a tip) — cheap, one API read per commit, and it is the only instrument that would have named `e308577` at the time. Note the honest scope: only a push **tip** gets a run, so the assertion must be about tips, or it manufactures alarms by construction. · ~0.2d · gated on nothing · surfaced iter-150.

71. **w-mission-critical-state-lives-in-a-directory-the-os-wipes-on-boot** · clause-2 ·
    **THREE DIFFERENT KINDS OF MISSION STATE WERE PARKED IN `/tmp`, AND ONE macOS REBOOT TOOK ALL
    THREE — THE TOOLCHAIN PIN, THE DRIVER'S CRASH LOG, AND A SPRINT WORKTREE.** Measured iter-151:
    `kern.boottime` = `1788395029` (2026-09-03 02:23:49 local); `/tmp/ailang-v0300/ailang` `ls`
    rc=1; `/tmp/ailang-mission-world.log` **217 bytes** holding one line despite the driver only
    ever appending (`tee -a` at `mission-control.sh:98`, `>>` at `:544`); `/private/tmp/wt-row56`
    gone, leaving a `prunable` `git worktree list` entry. **The pin half is CLOSED by iteration
    151** (moved to `~/.pinned-ailang/ailang`, CI's own path — see the Repo Profile), so what this
    row tracks is the remaining two, and they are worse than the pin because neither fails closed.
    **(a) The driver log is this loop's ONLY crash forensics** — `HARD TIMEOUT` and `STALL:` are
    written there and nowhere else — so any fire that dies becomes undiagnosable the moment the rig
    reboots, which is exactly what happened to the row-56 slot. The launchd `StandardOutPath` copy
    (`/tmp/ailang-mission-world.launchd.log`) is in the same directory and is therefore not a
    backup at all; measured, both files were 217 bytes with identical content and **different
    inodes** (`66036286` vs `66032033`), i.e. two copies of the same loss. **(b) Sprint worktrees
    under `/tmp`** lose an interrupted iteration's uncommitted work — precisely the residue the
    skill's died-mid-flight trace (c) exists to recover, deleted before anyone can read it.
    **Fleet-owned, so this routes rather than lands:** `LOG=/tmp/ailang-mission-${MISSION_NAME}.log`
    is set in the frozen-core driver (`tools/launchd/mission-control.sh:76`) and every mission on
    this rig has the identical exposure, so the fix is a fleet commit plus a cross-mission note,
    never a local edit. World's side is the worktree convention and this row. **Do not let this
    grow into "make the driver crash-proof"** — the deliverable is only that the evidence of a
    crash outlives the crash. · ~0.2d · gated on nothing · surfaced iter-151 (controller-measured).

72. **w-queue-closed-count-is-an-increment-chain-no-instrument-reproduces** · clause-2 ·
    **EVERY STATUS STAMP PUBLISHES AN `N of M ROWS CLOSED` HEADLINE, AND IT IS CARRIED FORWARD BY
    ADDING ONE — NOTHING ON THIS RIG RE-DERIVES IT, AND TWO ATTEMPTS TO DO SO THIS ITERATION
    DISAGREED WITH THE RUNNING NUMBER AND WITH EACH OTHER.** Iteration 150 published *"58 of 70"*,
    reached as iteration 149's *"57 of 66"* plus one — and it is wrong by at least one in a way
    nobody could see, because row 56 had landed and was untagged (see that row). Iteration 151
    then tried to measure it properly, twice. **Attempt 1** (tag keyword on each row's FIRST line)
    read **36 closed of 72**. **Attempt 2** (block-wise, tag keyword anywhere in the row body) read
    **53 closed of 70** — and its own enumeration control fired first, catching that the initial
    form had stopped at row 66 because a row body contains a line beginning `## `. Three numbers,
    none agreeing, and attempt 2 is *demonstrably* over-counting: it classes row **57** as closed,
    which is false, because prose inside an open row matches a closure keyword. **So the honest
    state is that no verified census exists** and this row does not assert one; the running count
    is an increment chain whose only audit is that each iteration remembered to add one. **This is
    the mission's own recurring shape aimed at its own progress metric** — the same
    unauditable-summary defect Gate 0's sweep rule already fixed once by demanding a per-issue
    table instead of a *"0 of 52"* sentence, which was likewise false when measured. **THE ITEM:**
    make closure a machine-readable field rather than prose a regex guesses at — a leading
    `[LANDED]`/`[ROUTED]`/`[RULED OUT]`/`[PARKED]`/`[NEXT]` tag in a fixed position the queue's own
    header already promises — plus a script that prints the per-row table and a total, with a
    known-open and a known-closed row as controls in the same call. Until it exists, a STATUS
    stamp should say the count is carried, not measured. · ~0.3d · gated on nothing ·
    surfaced iter-151 (controller-measured).


**Document created**: 2026-07-23 (bootstrap, attended). **RATIFIED 2026-07-23** (iteration 0,
attended: Mark + World coordinator) — record on issue #1; advisory-quorum ledger in
`.ailang/state/mission-quorum/`. Sprint routing is authorized from the next loop fire.

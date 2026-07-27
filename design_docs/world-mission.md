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
- **STATUS rotation is a HAND edit, never a regex script (iter-18 scar).** A "move the 4th stamp to
  the archive" script that bounded the stamp with the next `---` deleted **293 lines** of this
  charter (bar, Conflict Surface, guardrails, routing policy, queue, Premise Verification Log) —
  the stamps are not `---`-delimited. Caught by `git diff --stat` before commit and restored with
  `git checkout --`. Always `git diff --stat` the charter before committing it.

---

## STATUS (rotation rule)

Newest **3** STATUS stamps live here; older ones move to `world-mission-status-archive.md`.
At Gate 4, after adding your stamp, move the now-4th stamp to the TOP of the archive file.

## STATUS 2026-07-27 (iteration 19) — **`w-worldd-m2` M2.A checkpoint **A2 LANDED**: the `ailang-worldd` daemon shell — PR #11 → squash `39b2115`, dev CI green (both jobs).** Mid-sprint EXECUTE ("Plan exists" lane): no new doc, no quorum, no planner. Delivered `cmd/ailang-worldd/main.go` + `main_test.go` + `host/daemon/daemon.go` + `daemon_test.go` (1,754 insertions across 2 rounds, zero deletions): **loopback guard** mirroring the Z3-proven `isLoopbackHost` exactly (non-loopback bind REFUSES startup; M2 ships no override flag, so local-first is structural not advisory), the **D7 bound-constant block** (all four `http.Server` timeouts set at construction + `maxCommitBytes = 8388608` pinned to the Z3-proven `withinCommitBytes` so the Go constant and the frozen semantic bound cannot drift), the **lifecycle** (fail-closed `store.Open` → optional interpreter archive → idempotent `registry.Bootstrap` with divergent-head FATAL never silent → listen → stdout announce of the resolved address → serve → SIGINT/SIGTERM → bounded drain then hard close), and **`/v1/health` + `/v1/head`**. **ROUTING INCIDENT — the pinned executor model does not exist for this CLI**: `$MISSION_EXECUTOR_MODEL=codex:gpt-5.6-sol` failed its pre-flight probe with `400 … 'gpt-5.6-sol' requires a newer version of Codex` (codex-cli 0.137.0) — **NOT quota**: a second bounded probe without `--model` returned rc=0 on default `gpt-5.5`, so the lane is healthy and the MODEL PIN is unreachable. **The driver's own probe (`tools/launchd/mission-control.sh:248`) omits `--model`, so it false-greens a model pin** — FROZEN CORE, so NOT edited locally; routed upstream as a proposal on both channels. Executor fell back to `$MODEL` (opus), **FLAGGED**; `codex:gpt-5.5` deliberately NOT substituted (an unratified model swap is a routing-policy change, which the charter gates behind evidence). **The independent judge earned its keep**: sprint-evaluator (sonnet; generator≠judge holds) returned **PASS-WITH-CONDITIONS 79/100 with ONE BLOCKING condition** — `TestDaemonDependencyAllowlist` had been **dropped**, leaving Decision 4's "zero-cloud is enforced, not asserted" and the charter's *local-first is inviolable* guardrail as prose with no gate behind them. **The controller verified the finding rather than adopting it**, and found the judge's own framing overstated: the plan JSON is **internally inconsistent** (says "the **M2.B** dependency-allowlist test" twice, "**M2.A**" once, and omits it from A2's file list), so the executor's reading was defensible — but the **quorum-reviewed design doc governs** and places the test in M2.A four times over. Closed in a bounded round-2 fix (`a1cc5fa`): walks `go list -deps` over BOTH daemon-core trees (237 transitive packages), **`t.Fatalf`s rather than `t.Skip`s** when `go` is missing (a skip would be the V27/B1 silent-skip class), fails on an **empty** list (S6 null case), and **names offenders**. Judge re-verified round 2 narrowly → **BLOCK-1 CLOSED, revised 92/100**, zero new blocking. **Controller's own independent evidence** (never laundering a sub-agent claim): both gates green on the pinned v0.30.0 with `host/replay` **RUNNING 12.1–12.3s not SKIP** and `verify_ail.sh` at EXACTLY 4/4 identities / 9 modules / 14 tests; scope clean **by diff not by claim** (all six other `host/*` packages + `world/` + `scripts/` + `.github/` + `go.mod`/`go.sum` byte-unchanged); **mutation 1** — subtly widening `isLoopbackHost` to also accept `0.0.0.0` → RED at two independent tests in 1.2s; **mutation 2** — dropping `golang.org/x/sys` from the allowlist → RED naming `golang.org/x/sys/unix`, proving it walks the REAL build graph; and **LIVE with real OS processes**: ephemeral-port announce, `/v1/health` returning the real archived interpreter HashRef + pinned `AILANG v0.30.0` manifest, **a SECOND OS PROCESS on the same DB failing closed in 0s rc=2** with `WriterAlreadyActive` as `StartupError{Stage: store-open}` — A1's ratified single-writer invariant holding at its first real consumer — `--bind 0.0.0.0` refused rc=2, SIGTERM clean. `metered=$0.00` (codex probes were the ChatGPT subscription lane, invoked with `env -u OPENAI_API_KEY` so an ambient metered key could not silently bill — the `claude-sub` call-site discipline applied to the OpenAI lane; executor+evaluator on subscription Agent-tool pins). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (`b1e2b33`→now `39b2115`), CI green, no `[nightly-eval]` issues, no new `@MarkEdmondson1234` comment on #9 or predecessor #1 (watermark unchanged at `2026-07-27T08:55:11Z`), inbox = 3 eval-suite FYIs (not outranking), no rotation due. **Two process fixes, no skill edit** (the one skill-adjacent finding is a DRIVER bug — the skill's recipe correctly passes `--model`): the Repo Profile now records the codex model-pin reality, and the log entry must now ENUMERATE the evaluator's non-blocking carry-forwards (iter-18 recorded "three" as a bare count and they were LOST; this iteration's nine are enumerated CF-1…CF-9 in the log). No routing-policy change. **Next: checkpoint A3** — bench harness (`BenchmarkStoreCommit`/`HeadRead`/`Health` with p50/p95 via `b.ReportMetric`), `bench/BASELINE.md` (day-1 budget; REST-commit + log-range rows explicitly PENDING M2.B), `scripts/bench_worldd.sh --smoke` failing on a MISSING BENCHMARK NAME (not a zero line count — `go test -bench` exits 0 on no-match, the V27/B1 class), and the CI bench-smoke step. That completes M2.A; then M2.B, M2.C.

## STATUS 2026-07-27 (HUMAN-SURFACE FOUNDING DOC, attended) — **Mark named the human↔AI surface THE key novel piece of the vision; coordinator authored [HUMAN-SURFACE.md](HUMAN-SURFACE.md) v0.1 attended**: seven interaction principles (decision packets · grounded prose · trust-gradient PROVEN/TESTED/ATTESTED/CLAIMED · time-as-navigation · attention choreography · five zoom levels · speculation-as-gesture), the MEDIUM answer (a **workbench over the world graph, not a transcript** — chat is a lens, CLI is plumbing-truth, desktop = thin Tauri shell over the same localhost surface justified by P5 only), sign-the-type input grammar, anti-patterns (grade laundering = cardinal sin). Queue: new item **6b w-human-surface** GATES item 7 (w-approval-inbox) + all M6 work — quorum + full §7 ratification at pick. **UPDATE (same day, attended): Mark PROVISIONALLY RATIFIED the v0.1 principles ("good principles for now, yes good start") → binding working basis effective now; mockup design fixtures committed at `design_docs/mockups/{approval-inbox,grounded-prose}.html` and referenced from the doc (§6.5) — fixtures preserve GRAMMAR, not styling.** (The earlier 4-stamp note is resolved — iter-19 completed the hand-rotation.)

## STATUS 2026-07-27 (iteration 18) — **`w-worldd-m2` UNPARKED on Mark's ratification, and the RATIFIED single-writer kernel change LANDED: PR #10 → squash `b0deedb`, dev CI green (both jobs).** Gate 0 caught the human directive — `@MarkEdmondson1234` on #9 @ `2026-07-27T08:55:11Z`, body **"A"** = ENFORCE single-writer (watermark advanced `2026-07-23T20:13:54Z`→`2026-07-27T08:55:11Z` BEFORE routing; predecessor #1 re-checked per the rotation-week catch — nothing new). **No third quorum round**: a HUMAN RATIFICATION outranks the reviewer verdict and the charter's unpark instruction is r3 → sprint-planner (this is NOT the narrow-refinement carve-out and did not need it). **Designer rotation advanced `claude:claude-fable-5` → `codex:gpt-5.6-sol`** — the very reviewer that raised the objection, now authoring its own ratified fix (probe rc=0 after a PATH miss). **r3** rewrote **Decision 2** (fail-closed `store.Open`: canonical DB identity → non-waiting OS exclusive lock on `<canonical-db>.writer.lock` → only then the SQLite handle; structured `WriterAlreadyActive`; additive `OpenReadOnly`; crash/stale-lock semantics; cross-process proof) and carried it through Status/Park box, P6, High-Impact table, Design Freeze, D1, **D5 (`--addr` promoted to a global client flag — the r2 `gemini-3-1-pro` non-blocking nit)**, M2.A/B/C (day split re-cut 1.0/0.6/0.4d, still ~2d), the aggregate file table (~2,135 LOC), Conflict Surface, Acceptance Criteria, A4/A6/A10/A11. **Controller review found ONE defect → bounded r3b fix pass** (same lane): non-file DSNs were unspecified, so `Open(":memory:")` — which has TWO landed call sites — would have created a literal `./:memory:.writer.lock` and self-contended, **falsifying the doc's own "landed suites stay green unmodified" claim**; in-memory DSNs now carve out pre-canonicalization (`5b105b5`). **sprint-planner (opus)** kept the 3 milestones and added a `ratified_kernel_change` block, BINDING per-milestone `non_vacuity_requirements`, and an A1/A2/A3 split of M2.A with **A1 = `safe_landing_point`** (M2.A's ~1,550 planned LOC exceeds M1's largest landed milestone; M1's doc estimate ran 2.2× low). **sprint-executor (opus, isolated worktree)** shipped **A1 only**: 1,017 insertions across 5 files, `host/store` ONLY — `syscall.Flock(LOCK_EX|LOCK_NB)` from **stdlib** (not `x/sys`, so M2.B's dependency-allowlist test stays simple); non-unix arm fails closed and is honestly marked untested. **Controller INDEPENDENTLY re-verified**: both gates green on the pinned v0.30.0 with **`host/replay` RUNNING 12.2s (not SKIP** — the V27/B1 false-green class stayed closed) and `verify_ail.sh` at EXACTLY 4/4 identities / 9 modules / 14 tests; `schema.sql` **byte-for-byte unchanged**; `world/` + the other five `host/*` packages + all five landed `store.Open` call sites untouched **by diff, not by claim**; **and an independent 5th mutation (`LOCK_EX`→`LOCK_SH`) turned the suite RED** naming a live helper PID — an in-process mutex is structurally incapable of passing, which is the entire point of arm A. **sprint-evaluator (sonnet; generator≠judge: opus≠sonnet) PASS 97/100, ZERO blocking**; 3 non-blocking carry-forwards → A2. Gate 3b: PR run `30256646182` + dev-merge run `30256701072` both **completed/success**, both jobs each; worktree removed. `metered=$0.00` — mid-iteration an attended session flipped codex to **ChatGPT subscription auth** (`auth_mode=chatgpt`; metered key moved aside 11:33 local), so both designer runs were a quota lane; honest residual: the r3 run launched ~2 min after that flip, worst case ≈$0.23 if it raced it — far under the $5 ceiling. **Critical Principle 0 exercised**: `tools/launchd/mission-control.sh` (FROZEN CORE) appeared modified+uncommitted mid-iteration (an attended session flipping `MISSION_EXECUTOR_MODEL`'s default to codex); left **completely untouched**, and local `dev` was reconciled with `git reset --mixed` + a scoped `git checkout -- host/store/` precisely so that sibling work survived. **Controller self-inflicted scar, caught before commit**: a regex STATUS-rotation script deleted 293 lines of this charter (it bounded a stamp by the next `---`); `git diff --stat` caught it, `git checkout --` restored it, and the Repo Profile now forbids scripted rotation. Preflight otherwise clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (`6fba10d`→now `b0deedb`), CI green, inbox empty, no `[nightly-eval]` issues, no rotation due (#9 created today, 3 comments). Item 3 `w-worldd-m2` → **[IN-SPRINT]** (r3 applied, plan approved, A1 of M2.A's 3 checkpoints landed). **ONE process fix, no skill edit**: a 3-instance-in-one-iteration PATH class (`codex` rc=127 `env: node: not found`; `go: command not found`; `gh` already in memory) → the Repo Profile now records that the tool shell omits `/opt/homebrew/bin` and that EVERY sub-agent directive must say so. No routing-policy change (7th consecutive clean landed sprint corroborates opus-executor / sonnet-judge / generator≠judge). **Next: checkpoint A2** — `cmd/ailang-worldd` + `host/daemon` shell (config, loopback guard, D7 bound constants + the four `http.Server` timeouts, bounded shutdown, `/v1/health`, `/v1/head`), folding the 3 carry-forwards; then A3 (bench harness + `bench/BASELINE.md` + `scripts/bench_worldd.sh --smoke` + CI bench-smoke) completes M2.A.


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
3. [IN-SPRINT — r3 APPLIED, plan approved, **M2.A/A1 LANDED 2026-07-27 (iter-18), PR #10 → squash `b0deedb`; M2.A/A2 LANDED 2026-07-27 (iter-19), PR #11 → squash `39b2115`; both dev CI green, both jobs**] **w-worldd-m2** · **single-writer fork resolved: ENFORCE (arm A), ratified by Mark attended and now SHIPPED — `store.Open` is the fail-closed writer path (non-waiting `flock(LOCK_EX|LOCK_NB)` on `<canonical-db>.writer.lock`, structured `WriterAlreadyActive`), additive `store.OpenReadOnly`, in-memory DSNs exempt, `schema.sql` byte-unchanged. Enforcement is proven CROSS-PROCESS (subprocess helper + READY handshake + <2s non-waiting assertion + negative control); a controller mutation (`LOCK_EX`→`LOCK_SH`) turns the suite red, so an in-process mutex cannot pass. evaluator PASS 97/100, zero blocking. **A2 (iter-19) LANDED the daemon shell**: `cmd/ailang-worldd` + `host/daemon` (loopback guard mirroring the Z3-proven `isLoopbackHost` — non-loopback bind REFUSES startup, no override flag; the D7 bound-constant block with all four `http.Server` timeouts + `maxCommitBytes` pinned to the Z3-proven `withinCommitBytes`; fail-closed lifecycle, stdout address announce, bounded drain then hard close; `/v1/health` + `/v1/head`), plus **`TestDaemonDependencyAllowlist` — zero-cloud now ENFORCED, not asserted** (walks `go list -deps` over both daemon-core trees, fails loudly if `go` is missing, fails on an empty list, names offenders). The evaluator BLOCKED round 1 because that allowlist test had been dropped; closed in a bounded round-2 fix and re-verified → **92/100**. A1's writer lock verified holding at its first real consumer: a second OS process on the same DB fails closed in 0s. Remaining: A3 (bench harness + `bench/BASELINE.md` + `scripts/bench_worldd.sh --smoke` + CI bench-smoke) completes M2.A; then M2.B, M2.C.** ·
   clause-2 · `ailang-worldd` local daemon: SQLite, REST API, CLI; zero cloud deps; kernel perf
   budget measured and recorded from the first commit · ~2d. **Design doc AUTHORED + quorum-run
   (2 rounds)**: `design_docs/planned/w-worldd-m2.md` (Fable `claude:claude-fable-5` designer,
   rotation; checked sketch `sketches/worlddapi.ail` = 3 Z3 contracts, gate green on pinned
   v0.30.0). Both reviewers ACCEPT the direction; r1's 2 objections (bounded-waits/body-limit +
   log-range N+1 hidden from the perf budget) RESOLVED in revision r1. Re-quorum r2 (the one
   allowed) → `gemini-3-1-pro` PASS, `gpt5-6-sol` REJECT on **single-writer enforcement**
   (asserted-not-enforced; A6). Carve-out N/A (offers a FORK needing controller judgment; one arm
   opens the LANDED M1 `host/store` → kernel change = ratification-class per guardrail). **Parked
   for @MarkEdmondson1234**: (A) enforce via a store writer-lock (`OpenWriter`/`WriterAlreadyActive`,
   ratification-class) vs (B) withdraw the sole-handle claim + downgrade A6 to SQLite's own
   transaction serialization. See the doc's Park box. Unpark = apply the chosen arm as r3 →
   sprint-planner.
4. **w-effect-broker-m3** · clause-3 · effect broker with FS / Git / Model (`std/ai`) /
   Human.Approve handlers; effect-result recording; capability + budget checks; first physical
   isolation floor · ~2–3d
5. **w-mcp-projection** · clause-6 · project the transition registry over MCP + publish the A2A
   agent card (reuse `ailang serve-api --mcp/--a2a` machinery — do not reinvent) · ~1d
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
| serve-api projects `.ail` exports as MCP tools (path (a), static case) | LIVE TEST 2026-07-23 | `ailang serve-api --mcp-http --port 8199 sketches` → `tools/list` returned `plan`/`verify`/`commit` with JSON schemas + effect rows in descriptions; server killed after, port freed. Dynamic/capability-filtered projection deliberately NOT claimed — w-mcp-projection acceptance criteria |
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

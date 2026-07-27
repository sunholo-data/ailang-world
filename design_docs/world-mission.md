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

---

## STATUS (rotation rule)

Newest **3** STATUS stamps live here; older ones move to `world-mission-status-archive.md`.
At Gate 4, after adding your stamp, move the now-4th stamp to the TOP of the archive file.

## STATUS 2026-07-27 (iteration 17) — **`w-worldd-m2` (clause-2 local daemon) NEW-DOC authored + quorum-run (2 rounds) → PARKED needs-human-review on ONE ratification-class decision.** Picked the queue head `w-worldd-m2` (genuine NEW-DOC: no doc/PR/quorum-artifact existed). Rotation designer resolved to `claude:claude-fable-5` (last-used codex → gemini is read-only/CapRemoteSandbox, can't author files → wrap to claude; Fable probe rc=0). Fable authored `design_docs/planned/w-worldd-m2.md` + checked sketch `sketches/worlddapi.ail`; controller independently re-verified on the **pinned v0.30.0**: check.passed, 2→3 Z3 contracts (isLoopbackHost, clampLimit, withinCommitBytes), 0 counterexamples/errors/skips, full `verify_ail.sh` sweep green (9 modules, 4/4 identities, 14 tests). **Pick-time quorum (Gate 2)**: r1 BLOCKED (gpt5-6-sol bounded-waits/body-limit; gemini-3-1-pro log-range N+1 hidden from the day-1 perf budget) → designer revision r1 (new Decision 7 bounded server/client/shutdown timeouts + 8 MiB MaxBytesReader→413 with Z3-proven `withinCommitBytes` + `PayloadTooLarge`; `BenchmarkLogRange` + baseline N+1 rows) → **re-quorum r2 (the one allowed)**: `gemini-3-1-pro` **PASS**, `gpt5-6-sol` **REJECT** on **single-writer enforcement** (D2/A6 asserted-not-enforced: no DB lease/process lock, `store.Open` open to embedded writers). **Narrow-refinement carve-out N/A** — the objection offers a FORK needing controller judgment, and one arm opens the LANDED M1 `host/store` (kernel change = ratification-class by the mission guardrail; carve-out first-use also needs Mark's OK, impossible headless). **PARKED for @MarkEdmondson1234** (Standing rule 2 — never force a guardrail): (A) enforce via a store writer-lock vs (B) withdraw the sole-handle claim + downgrade A6 to SQLite's own transaction serialization — see the doc's Park box. `metered=$0.186` (quorum r1 $0.0817 + r2 $0.104; designer Fable = subscription quota-bucket, $0). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (`6fbbafd`), CI green, no `[nightly-eval]` issues, no `@MarkEdmondson1234` comment on #9 or predecessor #1 (watermark `2026-07-23T20:13:54Z`), inbox empty. Weekly rotation already done this morning (issue #9, week of 2026-07-27; predecessor #1). Doc+sketch+bookkeeping committed to dev. **Next: w-worldd-m2 unparks on Mark's A/B answer (apply chosen arm as r3 → sprint-planner); else the queue head advances.** No skill edit; no routing-policy change. ⟨prior iter-16⟩: **`w-world-library-m1` M6 (CI Go gate + `scripts/verify_go.sh` + final green sweep + NB2) LANDED on dev — PR #8 → squash `a07ac96`, dev CI green (both jobs). This LANDS the entire `w-world-library-m1` item** (all 6 milestones; doc → `implemented/`, all 16 acceptance boxes checked). Mid-sprint EXECUTE, "Plan exists" lane → routed straight to **sprint-executor (opus, isolated worktree `sprint/w-world-library-m1-m6`)**; no new design/quorum/planner. Delivered: **`scripts/verify_go.sh`** — the durable LOCAL `go build ./... && go test ./... -count=1` gate, with an **anti-false-green guard** that FAILS LOUDLY (rc=1) if `AILANG_BIN` is unset or ≠`v0.30.0`, so the load-bearing `host/replay` tests can never silently `t.Skip` into a false green (the iter-13 V27 / iter-15 B1 class, now closed at the local gate too); **`.github/workflows/ci.yml`** — the `go-verify` job's two inline build/test steps replaced by one `./scripts/verify_go.sh` call (pinned-v0.30.0 download + version-assert + Z3 install + `ailang-verify` job all retained untouched); **`host/replay/replay_test.go`** — **NB2 resolved via a GENUINE end-to-end test** (`TestInterpreterMemberChangeDrivesRealReplayEndToEnd`): builds a byte-distinct WORKING second interpreter (thin `exec` wrapper → distinct content hash → distinct HashRef), drives all six replay steps through it, and asserts genuine end-to-end cache-MISS + authoritative resolution following the CHANGED member + original `(fn,interp1)` row intact + faithful byte-identical result/world-hash (NOT a tautology; the semantically-divergent-release negative case honestly documented as env-constrained → upstream multi-release scope). **Controller INDEPENDENTLY re-verified** in the worktree: `go build` rc=0; `verify_go.sh` green with `ok …/host/replay 13.4s` (replay RUNS, all 6 host pkgs `ok`); loud-fail guard rc=1 (unset + wrong-version); `verify_ail.sh` 4/4 identities + 14 tests / 8 modules; `go vet` + `gofmt` + `actionlint` clean. **sprint-evaluator (sonnet, generator≠judge: opus≠sonnet) PASS 96/100**, ZERO blocking conditions; independently confirmed the guard fires and replay runs. **Gate 3b: the anti-false-green wiring WORKS remotely** — PR #8 and dev-merge CI both show `AILANG_BIN=…/ailang (AILANG v0.30.0)` + `ok …/host/replay 9.0s` + `✓ go gate PASSED` (tests RAN in CI via the new script, no SKIP). `metered=$0.00` (executor opus + evaluator sonnet on subscription Agent-tool pins; no designer/quorum/codex/gemini spend). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff` (gh at `/opt/homebrew/bin`, prepended per-call), dev==origin/dev (`f3c73c9`→now `a07ac96`), CI green, no `[nightly-eval]` issues, no new `@MarkEdmondson1234` comment on #1 (watermark `2026-07-23T20:13:54Z`), inbox = 2 eval-suite status FYIs (not outranking, marked read). No weekly rotation (issue #1, created 2026-07-23; re-evaluate the Monday-07:00 boundary next iteration; <80 comments). Item 2 `w-world-library-m1` → **[LANDED]**. **Next: pick the top open queue item — `w-worldd-m2`** (clause-2, `ailang-worldd` local daemon: SQLite/REST/CLI, zero cloud deps, kernel perf budget measured from first commit; NEW-DOC — needs a design doc via the rotation designer + a Conflict-Surface vs `ailang serve-api`). No skill edit; no routing-policy change (6th consecutive clean-landed sprint on opus-executor/sonnet-judge/generator≠judge corroborates the pattern).

## STATUS 2026-07-24 (iteration 15) — **`w-world-library-m1` M5 (replay engine + replay-doubling + fixture episode, Decision 7) LANDED on dev — PR #7 → squash `ef06937`, dev CI green (both jobs).** Mid-sprint EXECUTE (doc quorum-cleared, plan approved, M1–M4 landed) → routed straight to **sprint-executor (opus, isolated worktree)**; no new design/quorum. Built `host/replay/replay.go` (384 LOC — recorded-episode manifest + authoritative replay `Engine`: resolves the executable from **each log entry's `Interpreter` HashRef** [not the epoch registry], verifies exe↔hash, consults the `(TransitionFn,Interpreter)` verify cache, invokes the archived released binary as a bounded `exec.CommandContext` subprocess on pinned source, byte-compares to committed goldens — **delegates per-transition execution to the artifact, never reimplements the interpreter** per §14) + `replay_test.go` (640 LOC, 12 tests: bit-for-bit, replay-doubling A==B==recorded w/ divergence-fails, authoritative-resolution, epoch-candidate-cannot-redirect, pair-member-change→cache-miss) + `testdata/` fixture (drives the REAL `world/transitions` lib) + folded-in M4 carry-forwards (`KindExecFailure`, sidecar-present/exe-absent `Resolve` edge). **Controller INDEPENDENTLY re-verified** in the worktree: `go build`/`go test ./host/...`(6/6 pkgs, replay 12/12, 11s driving the pinned binary)/`go vet`/`gofmt -l host/` all clean; spot-checked genuine delegation + committed-golden anchoring (not tautological). **sprint-evaluator (sonnet, generator≠judge: opus≠sonnet) PASS 73/100 round 1** — confirmed structural hermeticity (zero `registry` import in `replay.go`), pair-only cache key; raised ONE **BLOCKING merge-condition B1**: the CI `go-verify` job ran `go test` **without `AILANG_BIN`**, so all 12 replay tests silently `t.Skip`-ed in CI (a false-green, iter-13 V27 class). **B1 fixed in-PR** (bounded opus follow-up, `8ba5fe9`): go-verify now downloads the **pinned v0.30.0** binary (by TAG not `latest` — goldens are v0.30.0-scoped), sha256-verifies, asserts version, exports `AILANG_BIN`; +NB3 comment fix +5 delivered acceptance boxes checked +NB1 documented. **Gate 3b verified the fix WORKED**: go-verify CI log shows `AILANG v0.30.0` + `ok …/host/replay 3.306s` (tests RAN, no SKIP). `metered=$0.00` (executor opus + evaluator sonnet on subscription Agent-tool pins; no designer/quorum/codex/gemini spend). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (`f116174`→now `ef06937`), CI green, no `[nightly-eval]` issues, no new `@MarkEdmondson1234` comment (watermark `2026-07-23T20:13:54Z`), inbox = 1 FYI (V1 iter-102 report — cross-mission STATUS, no World demand, marked read). No weekly rotation (issue #1, <80 comments). Item 2 → **[IN-SPRINT]** (milestones 1–5 of 6). **Next: `w-world-library-m1` M6** — CI Go gate finalize + `scripts/verify_go.sh` + final green sweep → then `[LANDED]` + doc→`implemented/`; folds carry-forwards NB2 + NB5. No skill/routing change (B1 is a recurring code/CI pattern caught by the gate, below the ≥2-instance skill bar; 5th clean landed sprint corroborates opus-executor/sonnet-judge/generator≠judge).

## STATUS 2026-07-24 (iteration 14) — **`w-world-library-m1` M4 (interpreter artifact archive + epoch-1 registry bootstrap, Decisions 5+6) LANDED on dev — PR #6 → squash `8133573`, dev CI green (both jobs: `go host build + test gate`, `ailang-code verify gate`).** Mid-sprint EXECUTE (doc quorum-cleared, plan approved, M1–M3 landed) → routed straight to **sprint-executor (opus, isolated worktree)**; no new design/quorum. Built `host/archive/archive.go` (387 LOC — single opened byte stream hashed-while-copied via `io.TeeReader`, `fsync`→`chmod 0o555`→atomic `os.Rename`, JSON sidecar manifest {hash, `ailang --version`, size, OS, arch}, `HashRef`→path resolver; structured `ReplayError` with `errors.As`-distinguishable Kinds absent/mismatch/unsupported/exec-failure) + `host/registry/registry.go` (159 LOC — epoch registry `world/epoch-registry/v1`; `Bootstrap` creates epoch 1 naming the M1 release string as first advisory candidate through the ordinary `PutObject`+`SetRegistryHead` object/head mechanism, idempotent, divergent-head→structured error not silent overwrite) + the M3 carry-forward (`store_heads` folded from an inline `db.Exec` in `store.Open()` into `schema.sql`, behavior identical). **Controller INDEPENDENTLY re-verified** in the worktree: `go build`/`go test ./host/...`/`go vet`/`gofmt -l host/` all clean (archive 8 + registry 5 + store 7 tests green), `store_heads` fully removed from `Open()`, archive tests CI-safe via synthetic shell-script fixtures (the pinned-binary test skips when `/tmp/ailang-v0300/ailang` is absent → linux CI stays green with no Z3/binary need). **sprint-evaluator (sonnet, generator≠judge: opus≠sonnet) PASS 88/100 round 1**, no blocking issues; two non-blocking carry-forwards to M5/M6 (a `KindExecFailure` test; a sidecar-present/executable-absent idempotence-recovery edge — `Resolve()` still catches it as `KindAbsentArtifact`). Acceptance #1 (log entries stamp interpreter `HashRef`+version in `writtenBy`) is caller-contract, deferred to M5 replay — the archive side is proven. `metered=$0.00` (executor opus + evaluator sonnet on subscription Agent-tool pins; no designer/quorum/codex/gemini spend). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (`d690e45`→now `8133573`), CI green, no `[nightly-eval]` issues, no new `@MarkEdmondson1234` comment (watermark `2026-07-23T20:13:54Z`), inbox = 4 FYIs (eval-suite start, V1 iter-101 report, own 2 iter-13 reports — none outranking, no cross-mission demand). No weekly rotation (issue #1, <80 comments). Item 2 → **[IN-SPRINT]** (milestones 1–4 of 6). **Next: `w-world-library-m1` M5** — replay + replay-doubling (bit-for-bit episode reconstruction; folds the M4 carry-forwards). No skill/routing change (4th clean landed sprint corroborates opus-executor/sonnet-judge/generator≠judge).

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

- **Executor default**: NON-Anthropic lane — `codex` (e.g. `MISSION_EXECUTOR_MODEL=codex:gpt-5.6-sol`)
  per the shared quota plan, once the fleet lanes are wired for this repo; until then the profile's
  documented interim default (opus, subscription) applies. Wiring the codex lane is part of
  iteration 0's checklist.
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
3. [PARKED — needs-human-review (iter-17): ONE ratification-class decision] **w-worldd-m2** ·
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
7. [PARKED until 4 lands] **w-approval-inbox** · clause-5 · the approval inbox + provenance walk,
   first as CLI/generated projection (SCENARIOS.md scenario 1/3) · ~2d
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

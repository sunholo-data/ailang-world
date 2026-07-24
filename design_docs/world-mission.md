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

## STATUS 2026-07-24 (iteration 11) — **`w-m1-ailang-hardening` design doc DONE + quorum-cleared via the RATIFIED narrow-refinement carve-out; auditable reproduction fixtures committed; sprint PLAN produced. Item → [IN-SPRINT]; EXECUTE is next iteration.** Picked the top `[NEXT]` item (Mark-directed AILANG-feature retrofit of the M1 surface). Reality-check: no doc/plan/eval; M1 `.ail` confirmed to carry only decorative `bool` predicates + 0 tests (gap real). **Controller empirically grounded the syntax on pinned v0.30.0** BEFORE routing (the discoverability root-cause fix): `requires`/`ensures` Z3-verify via `ai-check`, inline `tests [((args),exp)]` run via `ailang test` — both viable. Routed to **design-doc-creator on the rotation designer `claude:claude-fable-5`** (probe rc=0, subscription-only; rotation next after codex) with an adapting brief (known repo friction) + the empirical findings + the ADT-return design question. Doc `design_docs/planned/w-m1-ailang-hardening.md`: resolves the `CommitResult` sum-type contract question via a Z3-proven `applyRevision` helper, 7 proven contracts + 8 tests, 22-row empirical verification log; found real v0.30.0 limits (V8 `plan` unprovable, V10 Z3-encoding-errors-exit-0-SILENTLY). **Quorum r1 BLOCKED** (gpt5-6-sol: aggregate-floor gate too weak; gemini pass) → gate-mandated **Fable revision** rewrote D5 to a hardcoded required-check MANIFEST. **Re-quorum r2 BLOCKED** on two NEW narrow, direction-preserving objections with concrete `proposed_fix` (gpt5-6-sol: commit an auditable verification-fixture dir; gemini-3-1-pro: route Leg-1 python error to stderr). **NARROW-REFINEMENT CARVE-OUT applied** — RATIFIED for the world mission at the M1 GO (attended, `world-mission-status-archive.md` L3), so "later iterations apply without re-asking"; both objections satisfy (a) verbatim `proposed_fix` + (b) no design-direction dispute; the carve-out is a CONTROLLER action → no 3rd Fable run (Fable-discipline preserved). Controller 2nd revision: committed `design_docs/verification/w-m1-ailang-hardening/` (4 fixtures + `run.sh` + captured `OUTPUTS.md`, pinned binary sha256 `e9746fef…`) + the stderr fix. **The reviewer-demanded fixtures caught & corrected two first-draft inaccuracies** (D-A: V3 "no user calls" overstated — an encodable-bodied callee verifies; inlining still safe; D-B: leg-2 secondary must be `len(tests[])==8`, not `passed_tests` which counts contract properties → flaky). `metered=$0.149` (two quorum rounds $0.067+$0.082; Fable designer+revision on subscription = $0.00). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (b0a632a), CI green, no `[nightly-eval]` issues, no new Mark comment (watermark `2026-07-23T20:13:54Z`), inbox = own prior discoverability finding (marked read). Also committed the untracked `.mcp.json` (discoverability MCP-wiring local fix, queue-referenced) + `.gitignore` for `.codex/`. **Next: EXECUTE — sprint-executor (opus, worktree) implements the 3-module retrofit + the manifest gate per the plan → evaluator (sonnet) → PR → CI green.** Designer rotation advanced codex→claude (write-back `claude:claude-fable-5`). No weekly rotation (issue #1, <80 comments). No skill/routing change this iteration (carve-out was a ratified-mechanism APPLICATION, not a new gate change).

## STATUS 2026-07-24 (iteration 10) — **M1 milestone 3 (SQLite store + atomic append-only log, Decision 4) LANDED on dev, CI green (2 jobs).** Item 2 `w-world-library-m1` M3 built by **sprint-executor** (opus, isolated worktree): `host/store/schema.sql` (66 LOC — 5 tables `objects`/`worlds`/`log_entries`/`epoch_registry_heads`/`verification_cache`; each HashRef one canonical `algo:digest` TEXT column; frozen 6-field LogHeader verbatim; `transition_ref` outside the frozen header; cache PK `(transition_fn_ref, interpreter_ref)`) + `host/store/store.go` (576 LOC — `database/sql` + pinned pure-Go `modernc.org/sqlite v1.54.0`; content-verified immutable object insert; single-transaction compare-and-append `Commit` with structured `*ConflictError` on stale head; pair-keyed verify cache) + `store_test.go` (362 LOC, 7 tests). The pure-Go SQLite driver deferred from M2 was pinned here. Plan pre-existed (iter-8) → no planner/designer. Controller independently verified `go build`/`go test ./...`/`gofmt` + **`CGO_ENABLED=0 go build` rc=0** (driver genuinely CGo-free → CI needs no C toolchain). **Sprint-evaluator** (sonnet; generator≠judge: opus≠sonnet) **PASS 88/100 round 1**, no blocking defects — confirmed all 3 acceptance criteria (verbatim frozen-header round-trip; one-transaction compare-and-append with `errors.As`-matchable ConflictError + clean rollback; cache key provably the exact pair, epoch metadata-only). Design-fidelity −5 for the non-frozen `store_heads` helper table created in `Open()` outside `schema.sql` (single-source-of-truth split) — judged **moderate, non-blocking, trivially fixable**, carried to M4 rather than force-fixed (respects generator≠judge). PR **#4** → both CI jobs green → squash **`a901c30`**, worktree cleaned, post-merge dev CI green (both jobs). `metered=$0.00` (all roles on subscription Agent-tool pins). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, no `[nightly-eval]` open issues, watermark `2026-07-23T20:13:54Z` unchanged (no Mark comment), inbox = own iter-9 report + 1 V1 eval-suite FYI (not-outranking). Item 2 → **[IN-SPRINT]** (milestones 1–3 of 6). **Next: milestone 4 — interpreter artifact archive + epoch-1 registry bootstrap (`world/epoch-registry/v1`); fold the M3 carry-forward recs (move `store_heads` into `schema.sql`; add a store-layer `entry_hash_ref` derivation test).** No weekly rotation (issue #1, <80 comments). No skill/routing change (3rd clean landed sprint on opus-executor/sonnet-judge/generator≠judge — corroborates keeping the pattern; the `store_heads` nit is a code-quality note inside a PASS, not a ≥2-gap skill signal).

## STATUS 2026-07-24 (iteration 9) — **M1 milestone 2 (Go host bootstrap) LANDED on dev, CI green — first Go code of the mission.** Item 2 `w-world-library-m1` M2 (Go module + `host/hashref` + `host/canon`, 2 dependency-free leaf packages, 635 LOC / 45 tests) built by **sprint-executor** (opus, isolated worktree) implementing Decision 2 (8-step UTF-8/LF canonicalization + `CanonicalizationError`) and Decision 3 (tagged `HashRef` `algo:digest`, sha256 via crypto/sha256, structured `HashError`, golden vectors). Plan pre-existed (iter-8) → no planner/designer this iteration. **Charter reconciliation**: the Repo Profile mandates extending CI in the PR where Go first lands, but the plan scheduled the CI Go gate at M6 — controller directed the executor to add the `go-verify` job (`go build`+`go test`) NOW, so Go code is never un-gated (plan's M6 `verify_go.sh` stands as a superset). Controller independently verified build/test/vet/gofmt/tidy + the ailang gate still green (`checked 8 module(s)`). **Sprint-evaluator** (sonnet; generator≠judge) **PASS 97/100 round 1**, no blocking defects (independently checksum-verified the golden sha256 vectors). PR **#3** → both CI jobs green → squash **`d5b155c`**, worktree cleaned, post-merge dev CI green (both jobs). **Deliberate deviation**: pure-Go SQLite driver pin DEFERRED to M3 (its first importing code — `go mod tidy` strips an unused require; evaluator-endorsed). `metered=$0.00` (all roles on subscription Agent-tool pins). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, no `[nightly-eval]` open issues, watermark `2026-07-23T20:13:54Z` unchanged (no Mark comment), inbox = 1 V1 eval-suite FYI (not-outranking). Item 2 → **[IN-SPRINT]** (milestone 2 of 6). **Next: milestone 3 — SQLite schema + store (immutable objects, worlds, atomic append-only log); pin the pure-Go SQLite driver here.** No weekly rotation (issue #1, 17 comments). No skill/routing change (planner-vs-charter CI-gate-timing logged as instance 1, below the ≥2 bar; 2nd landed sprint = 2 datapoints on the opus-executor/sonnet-judge pattern, below the ≥3 routing bar).


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

**[IN-SPRINT — doc quorum-cleared via ratified carve-out (iter-11); PREEMPTS w-world-library-m1 M4]
w-m1-ailang-hardening** · clause-1 · **Design doc DONE + carve-out-cleared**
(`design_docs/planned/w-m1-ailang-hardening.md`, Fable designer; 7 Z3-proven contracts + 8 inline
tests + a hardcoded required-check-manifest non-vacuous gate; empirically grounded on pinned v0.30.0,
22-row verification log + committed reproduction fixtures). Quorum: r1 BLOCKED (aggregate-floor gate
weak) → Fable revision to a required-check manifest → r2 BLOCKED on two narrow, direction-preserving
objections (auditability + Leg-1 stderr) with concrete reviewer `proposed_fix`. **NARROW-REFINEMENT
CARVE-OUT applied** (ratified for the world mission at the M1 GO): controller 2nd revision committed
the auditable fixture dir `design_docs/verification/w-m1-ailang-hardening/` + the stderr fix; the
fixtures caught & corrected two doc inaccuracies (V3 overstated; leg-2 secondary → `len(tests[])`).
Routed to **sprint-planner** → executor (opus, worktree) → evaluator (sonnet). **AILANG-knowledge review + retrofit of M1 milestone-1 (`world/{contracts,transitions,logepoch,types}.ail`)**
— the flagship AILANG showcase shipped using **none** of AILANG's distinguishing features
(0 Z3 contracts, 0 tests, only decorative `bool` predicates); root cause = discoverability, now
partially corrected (`.mcp.json` wiring the `ailang-docs` MCP; upstream fix → `mission-control`
msg `msg_20260724_114812_208ab38d` + issue `sunholo-data/ailang#476`). **The reviewer/executor
MUST load the version-locked syntax first** (MCP `prompt_get` or `ailang prompt`) before touching
`.ail`. **Scope (all empirically grounded against pinned `v0.30.0` + Z3 4.16.0):**
(1) add `requires`/`ensures` contracts to the **int/bool/record invariants** and let `ai-check`
actually prove them — e.g. `commit` `ensures` the `Applied` world's `revision == w.revision + 1`
(Contract 4's core; VERIFIED-provable — a `nextRevision`-shaped postcondition proves live);
(2) add inline `tests [(in,exp)]` to the pure functions, **especially the string-rendering ones
Z3 SKIPS** — `renderRef`/`sameRef`/`cacheKey` use interpolation → builtin `show` (no SMT
encoding), so tests are their only machine check; (3) make the verify gate **non-vacuous** (assert
≥1 proven contract + tests present, not just module count). **OUT of scope (correctly absent — do
NOT add):** effects (M1 is intentionally pure `! {}`; shells belong to the Go host + M3 broker) and
package extensions (frozen core; clause 7 later). Go milestones M2/M3 are out of scope (Go, already
tested + evaluator-passed). Cheaper now than after M4/M5 build on these types. See memory
`ailang-feature-discoverability-gap`. · ~0.5–1d · queued iter-10 (Mark-directed).

1. [LANDED 2026-07-24] **w-log-epoch-decision** · clause-1 · ALL THREE decisions settled: D2
   (content-addressed transition fns) + D3 (SHA-256 + tagged HashRef) via quorum; **D1 RATIFIED
   by Mark attended: A+B-metadata hybrid** — authoritative pin = exact-binary content hash;
   release tag+commit + semanticsEpoch as compatibility metadata in every header (M8 promotion =
   policy change, never a format migration); hermeticity → replay-doubling in M1 acceptance.
   Doc: `design_docs/planned/w-log-epoch-decision.md` (SETTLED — ready for M1).
2. [IN-SPRINT — milestones 1–3 of 6 LANDED 2026-07-24 (iter-8/9/10), CI green]
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
   **Remaining: M4** interpreter artifact archive + epoch-1 registry bootstrap (`world/epoch-registry/v1`)
   → M5 replay+replay-doubling → M6 CI Go gate finalize (`verify_go.sh`). → [LANDED] + doc→implemented/
   when M6 lands. **NOTE: `w-m1-ailang-hardening` (top of queue) runs BEFORE M4 — the M1 AILANG
   surface gets its Z3 contracts + tests retrofit first, since M4/M5 build on these types.**
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
3. **w-worldd-m2** · clause-2 · `ailang-worldd` local daemon: SQLite, REST API, CLI; zero cloud
   deps; kernel perf budget measured and recorded from the first commit · ~2d
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

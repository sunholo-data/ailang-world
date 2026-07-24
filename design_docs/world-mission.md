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

## STATUS 2026-07-24 (iteration 12) — **`w-m1-ailang-hardening` EXECUTE attempted; Phase 1 (logepoch) LANDED on branch, Phases 2–4 BLOCKED by a newly-found v0.30.0 encoder limit that invalidates doc claim V5 → item PARKED for a designer revision + re-quorum (autonomous next iteration; not human-blocked).** Routed the quorum-cleared doc: **sprint-planner (opus)** → 4-phase plan JSON + handoff (`.ailang/state/sprints/w-m1-ailang-hardening.{plan.json,handoff.md}`). **sprint-executor (opus, isolated worktree `/tmp/wt-w-m1-hardening`)** built Phase 1 (D3 logepoch): `sameRef` field-eq body + proven `ensures`, `servesEntry` `ensures`, 8 named inline tests on `renderRef`/`sameRef`/`cacheKey`/`servesEntry` — **ai-check `verified:2` (sameRef, servesEntry), 8/8 tests pass, errors:0**; committed `35c3133`. Executor then **STOPPED per the design's own STOP-on-contradiction rule**: applying D2 verbatim to `contracts.ail` gave `verified:1, errors:3` — the 3 `Proposal`-taking predicates Z3-error `unknown sort 'Proposal'`. **Controller independently reproduced** (single-predicate + a minimal self-contained module: a record with an ADT-typed field → `unknown sort`, `errors:1`, **exit 0 silent**) and bisected the trigger to **any user-ADT-typed field** in the record (`Proposal.evidence: list[Evidence]`). Design claim **V5 ("all 4 predicates verify") is empirically FALSE against production types** — the committed fixture used a toy 2-field `Proposal`. Achievable proven set = **4** (`applyRevision`, `isValidNextWorld`, `sameRef`, `servesEntry`), not 7 → D2/D4/D5-manifest/`EXACT_TOTAL_VERIFIED=7` invalidated. **Did NOT force a shrunk manifest through** (Standing rule 2 — the quorum is the guardrail; a 7→4 gate-strength cut is exactly the reviewer's r1 concern and must be re-blessed, not controller-decided). Preserved Phase 1 by **pushing branch `sprint/w-m1-ailang-hardening`** (durable WIP, unmerged). Filed upstream **`sunholo-data/ailang#477`** (encoder: declare Z3 datatypes for ADT-bearing records; + make `ai-check` exit non-zero on `verify.errors>0`) + `mission-control` note. `metered=$0.00` (planner+executor on opus Agent-tool subscription pins; no quorum this iteration). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (`a4ec887`), CI green, no `[nightly-eval]` issues, no new Mark comment (watermark `2026-07-23T20:13:54Z`), inbox = eval-suite FYIs + own prior report (not-outranking). No weekly rotation (issue #1, <80 comments). **Next: rotation designer revises the doc to the achievable-4 scope → re-quorum ONCE → resume Phases 2–4 → evaluator → PR → CI; then w-world-library-m1 M4.** No skill/routing change (the STOP-and-report worked exactly as designed; the finding is a design-fidelity + upstream-escalation outcome, not a loop-process gap).

## STATUS 2026-07-24 (STANDARDS RATIFIED, attended) — **Mark's two findings closed as standing structure: (1) the zero-contracts gap and (2) slim-core/package-first were unwritten house style → now `design_docs/coding-standards.md` (S1–S6, binding, evaluator-scored) + repo `CLAUDE.md` (standing agent instructions: charter→standards→thesis reading order, pinned-binary rule, fluency protocol via the `.mcp.json` ailang-docs server).** Complements the in-sprint `w-m1-ailang-hardening` (code retrofit + manifest gate — the loop's lane); this is the never-again lane. Coordinator stood DOWN from the attended retrofit on sync (in-sprint collision avoidance).

## STATUS 2026-07-24 (iteration 11) — **`w-m1-ailang-hardening` design doc DONE + quorum-cleared via the RATIFIED narrow-refinement carve-out; auditable reproduction fixtures committed (`aa542a1`). Item → [IN-SPRINT]; sprint-planner → EXECUTE is next iteration (the doc's §Implementation Plan is the 4-phase basis).** Picked the top `[NEXT]` item (Mark-directed AILANG-feature retrofit of the M1 surface). Reality-check: no doc/plan/eval; M1 `.ail` confirmed to carry only decorative `bool` predicates + 0 tests (gap real). **Controller empirically grounded the syntax on pinned v0.30.0** BEFORE routing (the discoverability root-cause fix): `requires`/`ensures` Z3-verify via `ai-check`, inline `tests [((args),exp)]` run via `ailang test` — both viable. Routed to **design-doc-creator on the rotation designer `claude:claude-fable-5`** (probe rc=0, subscription-only; rotation next after codex) with an adapting brief (known repo friction) + the empirical findings + the ADT-return design question. Doc `design_docs/planned/w-m1-ailang-hardening.md`: resolves the `CommitResult` sum-type contract question via a Z3-proven `applyRevision` helper, 7 proven contracts + 8 tests, 22-row empirical verification log; found real v0.30.0 limits (V8 `plan` unprovable, V10 Z3-encoding-errors-exit-0-SILENTLY). **Quorum r1 BLOCKED** (gpt5-6-sol: aggregate-floor gate too weak; gemini pass) → gate-mandated **Fable revision** rewrote D5 to a hardcoded required-check MANIFEST. **Re-quorum r2 BLOCKED** on two NEW narrow, direction-preserving objections with concrete `proposed_fix` (gpt5-6-sol: commit an auditable verification-fixture dir; gemini-3-1-pro: route Leg-1 python error to stderr). **NARROW-REFINEMENT CARVE-OUT applied** — RATIFIED for the world mission at the M1 GO (attended, `world-mission-status-archive.md` L3), so "later iterations apply without re-asking"; both objections satisfy (a) verbatim `proposed_fix` + (b) no design-direction dispute; the carve-out is a CONTROLLER action → no 3rd Fable run (Fable-discipline preserved). Controller 2nd revision: committed `design_docs/verification/w-m1-ailang-hardening/` (4 fixtures + `run.sh` + captured `OUTPUTS.md`, pinned binary sha256 `e9746fef…`) + the stderr fix. **The reviewer-demanded fixtures caught & corrected two first-draft inaccuracies** (D-A: V3 "no user calls" overstated — an encodable-bodied callee verifies; inlining still safe; D-B: leg-2 secondary must be `len(tests[])==8`, not `passed_tests` which counts contract properties → flaky). `metered=$0.149` (two quorum rounds $0.067+$0.082; Fable designer+revision on subscription = $0.00). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (b0a632a), CI green, no `[nightly-eval]` issues, no new Mark comment (watermark `2026-07-23T20:13:54Z`), inbox = own prior discoverability finding (marked read). Also committed the untracked `.mcp.json` (discoverability MCP-wiring local fix, queue-referenced) + `.gitignore` for `.codex/`. **Next: sprint-planner (opus) turns the doc's §Implementation Plan into a sprint JSON → sprint-executor (opus, worktree) implements the 3-module retrofit + the manifest gate → evaluator (sonnet) → PR → CI green.** Designer rotation advanced codex→claude (write-back `claude:claude-fable-5`). No weekly rotation (issue #1, <80 comments). No skill/routing change this iteration (carve-out was a ratified-mechanism APPLICATION, not a new gate change).

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

**[PARKED — designer revision + re-quorum needed (iter-12); NOT human-blocked; still PREEMPTS w-world-library-m1 M4]
w-m1-ailang-hardening** · clause-1 · **Phase 1 (logepoch) LANDED on branch; Phases 2–4 BLOCKED by a v0.30.0 encoder limit that invalidates doc claim V5.**
(`design_docs/planned/w-m1-ailang-hardening.md`, Fable designer; doc claimed 7 Z3-proven contracts + 8
inline tests + a hardcoded required-check-manifest non-vacuous gate.) **iter-12 finding (reproduced by
executor AND controller on pinned v0.30.0):** a `requires`/`ensures` contract on a function whose
parameter is a record transitively containing a user **ADT** (sum type) fails with Z3 `unknown sort
'<Record>'` and `ai-check` **exits 0 SILENTLY** (the V10 class). `Proposal` has `evidence: list[Evidence]`,
so the 3 `Proposal`-taking predicates (`proposalMatchesWorld`, `verificationMatchesProposal`,
`commitAllowed`) **cannot be Z3-proven**; only `isValidNextWorld` (World/HashRef) + `applyRevision`
verify. **Achievable proven set = 4, not 7** → the ratified D2/D4/D5-manifest/`EXACT_TOTAL_VERIFIED=7`/V5
are invalidated. Filed upstream `sunholo-data/ailang#477` (+ `mission-control` msg
`msg_20260724_143026_0b2a75a0`); minimal repro + bisection attached. **Phase 1 preserved**:
`world/logepoch.ail` (`sameRef`+`servesEntry` verified, 8 named tests pass) committed `35c3133` on pushed
branch `sprint/w-m1-ailang-hardening` (NOT merged to dev — lands with the revised sprint). **Next iteration
(autonomous, no human needed):** rotation designer revises V5/D2/D4/D5 to the achievable-4 scope (the 3
`Proposal` predicates → documented-limitation rows, inline tests as their machine check per the doc's own
V8/§5 pattern) → **re-quorum ONCE** → resume Phases 2–4 (transitions `applyRevision` + contracts
`isValidNextWorld` + the corrected manifest gate + NT1/NT2) → evaluator → PR → CI. Surfaced to Mark FYI
(his directed item; descopes a ratified claim) but does not block on him. **AILANG-knowledge review + retrofit of M1 milestone-1 (`world/{contracts,transitions,logepoch,types}.ail`)**
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

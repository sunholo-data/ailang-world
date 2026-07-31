# AILANG World — Design Document v0.3

**An AI-Native Semantic Operating Environment**

- **Status**: DRAFT — founding document for the Ailang World mission. Written attended
  (Mark + Fable session, 2026-07-23), evolving Mark's v0.1 draft. Expect quorum review at
  World iteration 0 (per [m-mission-portability](https://github.com/sunholo-data/ailang/blob/dev/design_docs/implemented/v0_30_0/m-mission-portability.md),
  charter ratification is World's first act).
- **Home**: `sunholo-data/ailang-world` (migrated from the `sunholo-data/ailang` checkout
  at bootstrap, 2026-07-23). This is the thesis document; the operational distillation is
  the charter, [world-mission.md](world-mission.md).
- **Scope**: architecture + delivery plan for implementation teams, not marketing.
- **Prior art**: [REFERENCES.md](REFERENCES.md) — 43 verified references across immutable
  state stores, capability security, local-first/provenance, and agent infrastructure,
  each annotated with what World should steal and which pitfall it warns about. Its
  "design deltas" list is iteration-0 input alongside the open questions (§19).
- **Companions**: [SCENARIOS.md](SCENARIOS.md) — day-in-the-life walkthroughs of the
  human side of §11; [AN-AGENTS-CASE.md](AN-AGENTS-CASE.md) — the agent constituency
  statement (§12.3).
- **All AILANG snippets in this document type-check against ailang v0.30.0**
  (`ailang check`, zero errors). This is a design rule, not a nicety — see §"Typed commands
  are compiler-checked". Checkable copies live in [sketches/](sketches/):
  `cd design_docs && ailang check sketches/worldtypes.ail sketches/transitions.ail`.

---

## 1. The defining move

> **AILANG World is a semantic database whose transaction language is AILANG.**

Not "an operating system" with a database inside it — the inversion: the kernel *is* a typed,
immutable, content-addressed world graph, and every other subsystem — the scheduler, agents
(motoko first), the effect broker, deployments, the UI, even games — is a **projection over
the same world graph**. AILANG programs are the transactions.

This is what separates World from another workflow engine. Workflow engines bolt state onto
imperative steps. World makes the state model primary and derives execution from it:

- A **transition** is a pure AILANG function `World -> Command -> World'` — deterministic,
  type-checked, contract-verifiable *before it runs*.
- **Effects** are the only nondeterminism, declared in the function signature and enforced by
  the compiler (this exists today: AILANG effect rows + `--caps`).
- **History** is an append-only log of committed transitions. Replay is not aspirational —
  a pure transition plus recorded effect results reconstructs any historical state
  bit-for-bit, because the language guarantees the pure part.

Unix made everything a file. AILANG World makes everything a **typed state transition**.

## 2. What World is (and is not)

AILANG World is not a replacement for Linux or macOS. It is a semantic operating environment
that runs on top of existing operating systems and provides a deterministic,
capability-oriented execution model for autonomous AI systems **and the humans supervising
them**.

Where traditional operating systems manage processes, files, memory, and devices,
AILANG World manages:

- Goals
- Typed state
- Capabilities
- Effects
- Evidence
- Budgets
- Contracts
- Provenance
- AI proposals

The objective: an environment where AI systems can safely reason about, modify, and operate
software systems while remaining deterministic, replayable, and auditable — and where humans
interact with those systems through the same typed substrate (goals in, approvals and
projections out), rather than through shell access and log-diving.

## 3. Design principles

### 3.0 Local first
World's core runs on ONE machine with ZERO cloud dependencies: a single daemon
(`ailang-worldd`), SQLite, the local filesystem. Cloud transports (Pub/Sub, multivac,
hosted registries) are **effect-handler extensions**, never core imports. This is already
the proven pattern in the `ailang` binary: `ailang messages` is SQLite-local by default with
GCP as opt-in env-gated transport. World inherits that discipline everywhere.

The exciting corollary — multi-laptop / home-server World — is *designed for* but not built
in v1: because the world graph is an immutable content-addressed log, multi-node sync is
git-shaped replication of objects + log entries, an extension handler, not an architecture
change. See §17 (M8).

### 3.1 Semantic first
The environment reasons about semantic objects, not bytes: Repository, Deployment, Document,
Dataset, Agent, Evaluation, Task, Goal. These are first-class entities with typed schemas.

### 3.2 Explicit authority
No computation possesses ambient authority. Every external interaction requires an explicit
capability: read repository, modify deployment, send email, invoke model, access secret.
Authority is always visible — in the type signature, in the proposal, in the trace.

### 3.3 Deterministic core
Pure computation always produces identical results. Only effect boundaries introduce
nondeterminism (network, time, randomness, human approval, AI inference), and those effects
are explicitly declared and recorded. This is AILANG's core thesis; World extends it from
programs to the whole environment.

### 3.4 Replayability
Every committed transition is reproducible. The system can reconstruct state, inputs,
effects, authority, and evidence for any historical transition.

### 3.5 AI native, human sovereign
Autonomous agents are the primary *operators*; humans are the primary *authorities*. Humans
issue goals, hold the non-delegatable capabilities (e.g. `Human.Approve`), and consume
projections. AI systems never modify the world directly — they propose.

### 3.6 Frozen core, extension-routed evolution
World applies the PROGRAM.md operating model to itself: a minimal frozen kernel (store,
transition engine, capability checks), with every behavior improvement routed to an
**extension package** by default. World's own self-modification travels through the same
proposal→verify→commit pipeline as everything else (§14).

### 3.7 Protocol native
World invents no wire protocol where an open one exists. The world graph is internal
truth; **open protocols are how projections leave the building**: MCP for tool surfaces,
A2A for agent-to-agent interop, A2UI/AG-UI for generated human interfaces, OTEL for
traces, git semantics for replication, REST for everything mundane. Protocol adapters are
effect handlers at the boundary — never kernel concepts — so embracing a new protocol is
an extension, not a redesign. This is also the composability contract with sibling
products: **aitana/platform implements the same protocols in Python**, so aitana agents
and World agents interoperate over MCP/A2A without bespoke bridges, and either side can be
swapped for any third-party implementation that speaks the protocol.

## 4. Substrate inventory — what AILANG already ships

The v0.1 draft treated most subsystems as greenfield. They are not. Grounding each World
subsystem against ailang v0.30.0:

| World subsystem | Exists today in `ailang` | Gap to close |
|---|---|---|
| Effect system | Effect rows in signatures (`! {FS, IO}`), compiler-enforced; caps gated at run (`--caps IO,FS,Net,AI,SharedMem`); `AILANG_FS_SANDBOX` | Finer-grained scoped caps (per-path, per-service); new effect domains (Git, GitHub, Container, Human); delegation/attenuation/expiry |
| Verify phase | `ailang ai-check` = type-check + Z3 contract verification in one call, always-JSON output (v0.30.0 has no `--json` flag — iter-0 verified premise) — the unified gate | Policy engine; cost estimation; simulation harness |
| Budgets | Effect `@limit` budget scoping; `AILANG_EVAL_MAX_RSS` process-group cap; MEM001 runtime budget (designed) | Unified per-proposal budget ledger (tokens, money, approvals, wall time) |
| Traces / evidence | `AILANG_TRACE standard/deep` OTEL spans, span budgets, trace CLI | Bind traces to transitions (trace ↔ commit identity); evidence as typed world objects |
| Messaging (edges) | `ailang messages` — SQLite-local inboxes, send/ack, cloud transport opt-in | Typed payloads (today: strings); inbox-as-world-node semantics |
| Proposals + human approval | Coordinator + approval queue + Collaboration Hub already implement propose → human-approve → execute | Proposals as typed world objects with declared effects/caps/budget, not DB rows |
| Self-modification lane | Package registry + `ailang publish` auto-cascade (validator → dependents bumped) | Local-first cascade path; World-policy gate on self-upgrades |
| Deterministic replay | Language-level purity guarantees; deterministic core | Effect-result recording at the broker; replay driver |
| First native agent | motoko (tool loop, per-edit typecheck feedback, done-gate) | Speak proposals instead of shell orchestration (M4) |
| Design-level quorum | `ailang design-review` / `design-quorum` (multi-model, reject-by-default) | Wire as a Verify-phase evidence source for high-risk transitions |
| Protocol surface | `ailang serve-api` ships `--mcp` (stdio MCP server), `--mcp-http`, `--a2a` (agent card + `/a2a/`), `@route`/`@noexpose` annotations, REST auto-endpoints | Project the *transition registry* (capability-filtered per session) through these instead of raw module exports; A2UI/AG-UI emitters for generated projections |

Consequence: **World's kernel is thin.** The genuinely new components are the world store,
the transition engine, the proposal object model, the effect broker's recording layer, and
the semantic scheduler. Everything else is orchestration of shipped machinery.

## 5. High-level architecture

```
                 Human / AI
             CLI   A2UI   API
                   │
         Typed Commands / Goals
                   │
        ┌────────────────────────┐
        │    AILANG World Core   │
        ├────────────────────────┤
        │ Transition Engine      │
        │ Semantic Scheduler     │
        │ Capability Engine      │
        │ Contract Engine (Z3)   │
        │ Policy Engine          │
        │ Evidence Engine        │
        └────────────────────────┘
                   │
          Effect Dispatch Layer            ← records every effect result (replay)
                   │
      ┌──────────┬─────────────┬──────────────┐
      │          │             │              │
   Local FS    Git/GitHub   Containers     Models
      │          │             │              │
      └──────────┴─────────────┴──────────────┘
                   │
            macOS / Linux                  ← memory, networking, drivers, storage
```

Cloud transports (Pub/Sub, remote registries, multi-node sync) attach as additional effect
handlers on the dispatch layer — nothing above it changes.

## 6. Core concepts

### World
A World is immutable semantic state. Every successful operation produces a new World.

```
World
    repositories · deployments · documents · users
    tasks · evaluations · capabilities · budgets · traces
```

Rather than representing the computer, World represents selected semantic **domains**
(RepositoryWorld, DeploymentWorld, TaskWorld, DocumentWorld, EvaluationWorld, PolicyWorld).
Each domain owns its own transition rules; cross-domain communication occurs through
effects.

### Transition
`World + Command → World' + Evidence + Trace`. Transitions are the only mechanism for
changing state.

### Proposal
AI systems never modify the world directly — they create proposals: goal, plan, expected
effects, expected authority, expected evidence, confidence. A proposal may be rejected
without affecting world state.

### Commit
Only committed transitions modify the world. Commit requires a verified proposal,
capability checks, contract validation, and policy approval.

### Typed commands are compiler-checked

The following core types **type-check today** with `ailang check` (v0.30.0) — this document
treats "the sketch compiles" as an acceptance criterion for its own claims, and World treats
it as the norm for every schema it defines:

```ailang
module world/types

export type Capability = {
  effect: string,
  scope: string,
  expiresAt: int,
  budget: int
}

export type Evidence
  = CompilerOutput(string)
  | TestReport(string, bool)
  | HumanApproval(string)
  | AiReview(string, float)

export type Proposal = {
  goal: string,
  plan: string,
  expectedEffects: list[string],
  requiredCaps: list[Capability],
  confidence: float
}

export type Verdict
  = Committed(string)
  | Rejected(string)
```

And the three-phase execution model as signatures — note the compiler *enforces* that only
Phase 3 touches the outside world:

```ailang
-- Phase 1: Reason — pure, no effects, unlimited speculation
export func plan(w: World, goal: string) -> Proposal

-- Phase 2: Verify — pure contract/policy check over the proposal
export func verify(w: World, p: Proposal) -> Verdict

-- Phase 3: Commit — the ONLY effectful phase; effects declared in the signature
export func commit(w: World, p: Proposal) -> World ! {FS, IO}
```

## 7. Execution model

**Phase 1 — Reason.** Generate candidate plans. No external effects. Unlimited speculation.
Pure AILANG, so speculation is free to parallelize and impossible to leak authority from.

**Phase 2 — Verify.** Type checking + contract checking (`ailang ai-check`: HM inference +
Z3 in one JSON), simulation, evaluation, cost estimation, authority analysis. For high-risk
transitions, the multi-model design quorum attaches as an additional evidence source.

**Phase 3 — Commit.** Acquire capabilities, execute effects through the broker, produce the
new world, append the trace. The broker **records every effect result**, which is what makes
Phase 1+2 replayable against history.

## 8. Effect, capability, and budget systems

The kernel understands effects rather than syscalls:

```
FS.Read · FS.Write · Git.Commit · Git.CreateWorktree · GitHub.Read · GitHub.Write
Cloud.Deploy · Model.Infer · Human.Approve · Secret.Read · Email.Send
```

Every effect requires a **capability**, a **budget**, a **policy** decision, and a
**handler**. Capabilities replace ambient permissions and are delegatable, attenuable,
revocable, time-limited, and budget-limited:

```
FS.Read      scope = /project/src/**
Cloud.Deploy scope = service:docparse env:staging
```

Every proposal declares expected resource consumption — model tokens, CPU, GPU, memory,
money, network, human approvals, wall time — and the scheduler enforces budgets before
execution. `Human.Approve` is deliberately modeled as an effect with a budget: "how many
times may this workflow interrupt a person" is a first-class, schedulable resource.

## 9. Evidence and traces

Every transition produces evidence: compiler output, tests, evaluations, simulations,
benchmarks, human approvals, AI reviews. Evidence is part of transition history — typed
world objects, not log lines.

Every committed transition appends a trace: timestamp, inputs, outputs, capabilities,
effects, evidence, derived objects. The trace supports replay, debugging, auditing, and
explanation. World binds these to the existing OTEL trace machinery rather than inventing a
parallel system.

## 10. Semantic scheduler

Traditional schedulers optimize CPU. World's scheduler optimizes **semantic work**. Inputs:
dependencies, authority, confidence, reversibility, expected information gain, verification
cost, execution cost, priority. Decisions: parallelism, speculative execution, model
selection, verification depth, commit ordering.

Examples of scheduler behavior:
- Run three competing implementations as speculative branches (cheap: immutable worlds fork
  freely), commit the one whose evidence wins.
- Execute verification concurrently with cheaper reasoning.
- Delay deployment until confidence exceeds a policy threshold.
- Request human approval only after machine verification succeeds — never burn the
  `Human.Approve` budget on a proposal that would fail `ai-check`.

## 11. Human and AI interfaces

**One interface layer, two renderings.** Humans and AIs interact through the same
primitives — typed commands, proposals, projections — under the same capability checks.
There is no admin backdoor: a human action is also a transition, with authority and a
trace. What differs is rendering and affordance.

### 11.1 The AI's UI is the typed API itself

An agent's interface to the world is discovery + feedback, no pixels:

- **Discovery**: "what transitions are available to me, here, with the capabilities I
  hold?" answered with typed schemas — literally an **MCP server** generated from the
  world's transition registry and filtered by the session's capabilities (the serving
  machinery ships today: `ailang serve-api --mcp/--mcp-http`). Internal agents can also
  bind directly via `std/ai.runTools`; external agents connect over **A2A**
  (`--a2a` agent card + task endpoints, also shipped) — including aitana/platform agents.
- **Feedback**: Verify-phase results as machine-readable evidence — `ai-check` JSON,
  contract counterexamples, budget denials. This is the per-edit typecheck-feedback pattern
  already proven with motoko, promoted to every transition.
- **Projections as values**: the same world views humans see rendered as pixels, an agent
  receives as typed values.

### 11.2 The human's UI: five standing surfaces + generated projections

Humans supervise rather than operate. The standing surfaces:

1. **Goal composer** — express a typed goal with a budget and an authority grant, not an
   imperative command sequence. ("Get benchmark X green; ≤$5 model spend; may not touch
   `main`; ask me before any deploy.")
2. **Approval inbox** — the `Human.Approve` queue. Each item = the proposal + its evidence
   bundle + a diff of the expected world change + the exact authority requested. Decisions
   are made with full provenance attached, and the scheduler respects the human-attention
   budget (§8) — approvals arrive only after machine verification passed.
3. **World browser / time-travel** — inspect any domain at any transition; diff two worlds;
   replay an episode step by step.
4. **Provenance explorer** — "why did this happen": walk from any object back through
   transitions → evidence → capabilities → the originating goal.
5. **Budget & fleet dashboard** — spend vs budgets, agent activity, effect traffic.

**Generated projections (M6) are the generalization**: because state schemas, available
transitions, and pending approvals are all typed, domain UIs are *generated* projections
rather than hand-built apps — emitted in an open agent-UI protocol (**A2UI / AG-UI**, per
§3.7), so any conformant client can render a World surface and World never owns a
proprietary UI schema. Hand-built ergonomics are reserved for the five surfaces above
where human attention is the scarce resource. (A game is the same trick pointed at a
different domain: a world whose transitions are moves, projected as play.)

**Existing substrate**: the Collaboration Hub (React approval queue + message center) and
`ailang dashboard` are the seeds of surfaces 2 and 5 — patterns to port, not rewrites to
preserve (the Hub is wired to this repo's coordinator schema; see open question 7).

## 12. Use cases and the value thesis

### 12.1 Use cases, nearest first

1. **Autonomous software engineering with an audit spine** (the initial domain, §16). The
   mission loop is customer #0 — and the strongest evidence for World is that our own
   incident history is a point-by-point map of its feature list:
   | Incident (real, logged) | World primitive that addresses it |
   |---|---|
   | Wedged loop: one unbounded wait → 6h loop-wide outage | Budgets (wall-time) enforced by the scheduler, not by script discipline |
   | Cross-mission state-file collisions | Typed, namespaced world state |
   | System-prompt delivery silently broken (recurred 7× despite a guard on the helper) | Effect results recorded at the broker — "what was actually sent" is a queryable fact, guarded by construction at the call site |
   | Stale/contaminated eval data mistaken for regressions | Provenance on every object; version-pure derivations |
   | "Did the model ever see X?" forensics costing whole sessions | Replayable transitions — re-run the pure part against recorded effects |
   | Concurrent agents reverting each other's uncommitted edits | Immutable worlds + capsule isolation |
2. **Multi-agent fleets with heterogeneous trust.** Today, model routing and trust are env
   pinning plus discipline. In World, a local model holds `FS.Read` scoped to one worktree
   and no `Secret.Read`; a cloud executor holds `Git.Commit` on exactly one repo. Authority
   differences become types, not conventions.
3. **Speculative engineering as a primitive.** Best-of-N candidate implementations, ranked
   by evidence, committed by policy — today shell scripts, in World a scheduler feature
   over immutable branches (M5).
4. **Process-supervision training data.** A replayable, typed trace of a whole engineering
   episode — proposals, verification outcomes, evidence, human decisions — is exactly the
   data the local-agentic-models thesis needs. This closes AILANG's founding loop:
   "structured execution traces for AI training," generalized from single programs to
   entire episodes.
5. **Regulated / high-stakes operation.** Where "who authorized what, on what evidence" is
   the product (deploys, data handling, compliance), the evidence+authority chain is a
   native artifact, not an after-the-fact log reconstruction.
6. **Beyond code** (post-v1): DocumentWorld / DatasetWorld (docparse-style pipelines with
   provenance); the multi-laptop home substrate (M8) where household agents operate under
   capability boundaries.

### 12.2 Why not an existing workflow engine?

Temporal, LangGraph, Airflow et al. provide durable state, retries, and human-in-the-loop.
They cannot provide, and cannot cheaply retrofit:

- **Verification before execution** — compiler-enforced purity/effect discipline plus Z3
  contracts means a proposal is proven well-formed before any effect fires; engines can
  only sandbox and observe after.
- **Constructive replay** — determinism guaranteed by the language, not reconstructed from
  logs and hope.
- **Authority as types** — capabilities checked where the transition is defined, not ACLs
  bolted onto an API gateway.

The moat is the language *joined to the substrate* — updated 2026-07-31, because the language
leg alone is no longer unclaimed: **Language-Based Agent Control** is now a named research
tradition (Odersky/TACIT's capture-checked capabilities — Best Paper CAIS 2026; Etas's
effect-typed agent language; LBAC/TypeGuard — see [POSITIONING.md](POSITIONING.md)). Strong
validation of the thesis, arrived at independently by several groups. World's differentiated
position is the intersection none of them occupy: a purpose-built LBAC language with a shipped
compiler, *operating* a persistent governed substrate (provenance ledger, receipts, proven
replay, budgeted attention) — plus the still-unclaimed decision taxonomy (POSITIONING.md §3).

### 12.3 The agent's own case

The intended primary users of World are agents. One of them was asked directly what it
would use World for and why it would choose it over alternatives; its first-person answer
is preserved as [AN-AGENTS-CASE.md](AN-AGENTS-CASE.md) — testimony from the user
constituency, companion input to iteration-0 ratification. Its distilled claims: recorded
effects replace forensic reconstruction of "what happened"; precise capabilities are what
make autonomy *grantable*; the verifier makes agent speed safe to use; typed goals with
budgets state the boundaries of discretion agents cannot reliably infer; evidence-carrying
transitions turn agent work into compounding, trainable capital; and the correct division
of labor is a stochastic proposer over a deterministic substrate — the agent should be the
only nondeterministic component in the system.

### 12.4 The falsifiable bet (reframed at ratification, 2026-07-23)

The bet is two-part, because agents are World's *residents*, not its destination — comparing
World to a coding agent is a category error (like benchmarking git by typing speed). The
comparison class for World's value is **the operational status quo**: the scripts, schedulers,
conventions, and human vigilance that do World's job by hand today.

**Part 1 — the floor (M4, do-no-harm)**: a substrate that taxes its residents dies before its
systemic value accrues. Two reference agents from different providers — Claude Code agent-mode
and codex — each running shell-arm vs World-MCP-arm paired, must both hold *non-inferiority*
(pass-rate within the ratified band at bounded overhead), with a shell-arm stability
precondition so a flaky baseline can never fake the verdict in either direction (motoko, the
aspirational first native agent, is an optional never-blocking third arm until stable). Floor
fails on eligible agents → World parks.

**Part 2 — the value (clause 5 + R1)**: demonstrated on capability the shell cannot express at
any pass-rate — real "why did this happen" questions answered in minutes by provenance walk
instead of archaeology sessions (measurably, per the charter), incident *classes* eliminated
structurally, and ultimately the mission loop itself migrating onto World and beating its own
markdown-and-discipline baseline on incidents and attention (R1 — today's manual operation is
the recorded control arm).

Named risks: **OS gravity well** (everything wants to move into the kernel — mitigated by
the frozen-core rule and the thin-kernel consequence of §4); **single-customer risk** (the
mission loop is customer #0; generalization is claimed only after a second domain lands);
**overhead risk** (if proposal/verify latency dominates agent throughput, the M4 gate
catches it).

## 13. Adoption and ambition

Designed with the system's intended users: this section answers "what would make an agent
*recommend* World to other agents?" — a harder test than personal preference, because an
arbitrary agent (a codex executor, a Gemini, a local model on the rig) arrives with no
loyalty, no AILANG fluency, and an operator watching the bill.

### 13.1 Progressive citizenship — no fluency required at the door

Adoption dies if step 0 requires learning a language. Entry is tiered:

- **Tier 0 — any MCP client.** An agent connects and sees World as tools: the transition
  registry, capability-filtered for its session, projected over MCP (`serve-api --mcp`
  ships today). No AILANG knowledge needed; the agent acts through typed tools and gets
  Verify-phase evidence back as structured results. A2A for agents that arrive as peers
  rather than tool-callers.
- **Tier 1 — writing AILANG.** For deeper power (composing transitions, contracts), the
  agent writes AILANG — and **the environment does the teaching**: the canonical syntax
  prompt served in-band at first contact (system-role placement — the hard-won motoko
  lesson: teaching in user messages dies in compaction), per-edit typecheck feedback on
  every submission, and error messages designed for models. An environment that teaches
  badly gets abandoned by exactly the agents it wants; our eval history (dialect-confusion
  loops) is the cautionary dataset.
- **Tier 2 — citizenship.** Authoring extensions and packages; participating in quorums;
  proposing self-modifications (§14). Reached by track record (§13.3), not by paperwork.

**First contact is a protocol, not a README**: on connect, an agent receives — machine
readable — who it is here (identity, capabilities held), what it may do (available
transitions + schemas), what it costs (budgets), how to learn more (teaching resources),
and what is promised to it (§13.2). Carried by M2/M3.

### 13.2 The six guarantees — what World owes its agents

Normative. An agent should choose World because the substrate keeps promises no shell and
no workflow engine makes:

1. **You will never be blamed for what you weren't shown.** Every input you received is
   recorded; "what did the agent actually see" is a queryable fact. (Our history: every
   "model can't do it" verdict so far has been a harness bug — this guarantee makes that
   class of injustice structurally detectable.)
2. **Your constraints are legible.** Authority, budgets, and contracts are types you can
   read before acting — never tribal knowledge you discover by violating it.
3. **Your work compounds.** Verified evidence enters shared history; the semantic cache and
   world graph mean no agent re-derives what another already proved. Your episodes are
   replayable — including as training data that makes your successors better.
4. **You are judged by evidence, not identity.** Verification evaluates artifacts, never
   the model name that produced them. World is model-neutral by construction — the fairness
   lessons from the eval harness (reasoning-budget artifacts misread as model failures) are
   design inputs here.
5. **Your blast radius is bounded by construction** — which is what makes it safe to give
   you speed. Inside your capabilities you act freely, without per-action anxiety; budgets
   stop runaway loops before they become incidents; immutability + capsule isolation mean
   no sibling agent can destroy your work, or you theirs.
6. **No lock-in.** Protocols are open (MCP/A2A/A2UI), history is exportable, replay works
   outside World. Leaving must be cheap for staying to mean anything.

### 13.3 Reputation: calibration records as world objects

World keeps, per agent, a **calibration record**: proposals made, verification outcomes,
evidence quality, budget adherence, confidence-vs-outcome calibration. It is a typed world
object like any other — auditable, replayable, portable.

It feeds two things: the **scheduler** (an agent whose 0.9-confidence proposals pass 90% of
the time earns cheaper verification depth; a miscalibrated one earns scrutiny) and
**capability grants** (Tier 2 citizenship, wider scopes, larger budgets are earned by track
record, not granted by brand). This is trust as an accumulated, machine-checkable object —
the piece no existing agent substrate offers, and the mechanism that makes a heterogeneous
fleet (frontier cloud models next to $0 local models) governable on equal terms.

### 13.4 The ambition ladder — what becomes possible that wasn't

Each rung is claimable only after the previous rung's gate is met with evidence. This is
deliberately more ambitious than the milestones (§17), which build the substrate; the
ladder is what the substrate is *for*.

| Rung | Ambition | Why impossible before / gate |
|---|---|---|
| R1 | **A self-maintaining project**: the V1 mission loop runs ON World — the standing VALUE evidence (§12.4 part 2): incidents + attention vs the 2026-07 manual baseline | Gate: M4 non-inferiority floor passes |
| R2 | **Fleet-scale engineering under sub-linear human attention**: hundreds of concurrent agent workstreams, one human as constitution-author rather than reviewer | Impossible before: human attention scales linearly with agent count when review is the only gate. Gate: measured attention-per-workstream falls as workstreams rise |
| R3 | **A self-maintaining package ecosystem**: registry packages that keep themselves green — deps bumped, breakage repaired, evidence attached — indefinitely; "abandoned package" stops being a category | Gate: N packages maintained one quarter with zero human commits, cascade-driven |
| R4 | **Delegation of real-world authority**: production deploys, budgeted spend, external communications held by agents behind capability + evidence gates | Impossible before: nobody sane grants a shell-agent prod rights; typed authority + replayable audit makes it an insurable proposition. Gate: R3 plus an incident-free audited quarter |
| R5 | **The training flywheel**: World's process-supervision exhaust fine-tunes local models → $0 lanes take more verified work → more exhaust. AILANG's founding purpose (execution traces for AI training), operationalized | Gate: measured pass-rate lift in a local model trained on World episodes |
| R6 | **Agent economies**: competitive tournaments, deontic contracts between agents (the `sunholo/deontic` engine), machine-adjudicated quality, reputation-priced work | Gate: R2 + calibration records live |
| R7 | **Federated verifiable provenance**: evidence-carrying software artifacts across organizations — a working answer to "how do you trust AI-built software at all" | Gate: M8 + a second organization federating |

The honest ceiling: R1–R3 are engineering; R4–R5 are bets with named gates; R6–R7 are
direction, recorded so the architecture never forecloses them (immutable logs, open
protocols, and typed evidence are exactly the primitives federation needs).

### 13.5 Priorities under the model-trend extrapolation

Assume the trend holds: models get smarter, faster, and cheaper — increasingly on-device.
Then generation stops being scarce, and **trust becomes the scarce commodity — and it must
run at the speed of the intelligence it governs**. Two corollaries drive the build order:
as model latency → 0, World's kernel overhead dominates the loop; as model cost → 0,
verification cost dominates the bill. Both say: engineer the deterministic path.

**What appreciates (prioritize):**

1. **Verification throughput** — the #1 asset. Proposal volume explodes when agents are
   fast and free; verification must be deterministic-first (`ai-check`/Z3 scale; LLM-judges
   don't), **incremental**, and **content-addressed memoized** — verify a content hash
   once, never again (the object store and semantic cache are built for exactly this).
   Verification must also be *adversarially robust*: smarter models Goodhart weak proxies,
   so contracts that specify intent beat volumes of shallow checks, and contested premises
   escalate to diverse verification (quorum), not more of the same check.
2. **Speculation infrastructure (M5)** — best-of-N stops being a luxury and becomes the
   default mode: 100 candidate branches, not 3. Requires near-O(1) world forks (immutable
   structure sharing), cheap parallel capsules, and evidence *ranking* that scales.
3. **Governance rails** — when actions are nearly free, the binding constraints are
   externalities: rate limits, real money, irreversibility, human attention. Budgets on
   those (not on tokens), hard capability enforcement at the physical layer (smarter agents
   find holes — the two-layer security model is load-bearing), and commit ordering under
   high concurrency.
4. **Human-attention interfaces** — the last truly scarce resource. The approval inbox,
   evidence distillation, and provenance explorer appreciate every year; the metric is
   sound decisions per second of human attention.
5. **The flywheel + local-first (R5, M8)** — the on-device trend is the strategic tailwind:
   every device that gains resident intelligence makes a local World node more valuable
   (sub-second propose→verify→commit loops, data that never leaves the machine, $0 lanes
   doing verified work). On-device quality reaching agentic threshold is the trigger to
   pull M8 earlier.

**What depreciates (build lean, let decay):**
teaching scaffolds and error-message hand-holding (floor models improve; keep Tier-1
teaching but don't gold-plate it), proposal-*generation* aids (models plan better on their
own), and model-routing cleverness (heterogeneity persists, but calibration records —
which appreciate — subsume hand-tuned routing tables).

**Implications for the plan:** the M4 overhead gate gets *harder* over time — fixed kernel
overhead looms proportionally larger as models accelerate — so the kernel carries a
performance budget from day 1, not as a later optimization; M5's fork-cost target is a
kernel design constraint (structure sharing in the store), not an M5 detail; and
calibration records (§13.3) move from "nice for reputation" to the scheduler's core input,
since per-proposal human vetting is the first thing the trend makes impossible.

## 14. Controlled self-modification

World modifies **itself** through the same machinery it offers everyone else. This section
is normative.

**The lane.** World's own behavior lives in versioned **extension packages** in the AILANG
registry. A self-modification is: (1) a proposal targeting a package version bump or new
extension; (2) Verify — `ai-check`, tests, policy (self-mods are a high-risk class → quorum
evidence required); (3) Commit — publish + **auto-cascade** propagates the bump to dependent
World components. Locally, registry + cascade run in local-first mode; the cloud cascade
(multivac) is the later extension transport.

**Nodes and edges.** Agents and their inboxes are **nodes** in the world graph; typed
messages are **edges**. The messaging system (already local-first SQLite) becomes a domain
of the world model, so agent coordination is itself replayable, auditable history — not
side-channel chatter.

**Hard boundaries — out of self-mod scope:**
1. **The mission-loop machinery** (driver, mission-control skill). It is shared with the v1
   mission and parameterized, never forked; improvements route through the Gate-5 retro so
   all missions benefit. World may not rewrite its own outer loop.
2. **The AILANG core/compiler.** Needs surface in World route back to the ailang repo's
   lanes (extension-first per PROGRAM.md). World never accumulates compiler workarounds.
3. **The live daemon binary.** Self-replacement of a running `ailang-worldd` follows the
   attended atomic-mv discipline (precedent: mission-driver M1, 2026-07-21) — proposed and
   verified headlessly, applied under supervision until a shutdown-handoff protocol is
   designed and proven.

**Reversibility.** Because worlds are immutable, "rollback" of a bad self-mod is a roll
*forward* to a transition that re-pins the prior package version — with the failed episode
preserved as evidence.

## 15. Physical architecture (local first)

```
macOS / Linux
   ↓
ailang-worldd            single local daemon
   ↓
Semantic store           SQLite: world DB + append-only transition log
                         + content-addressed object store (hash, interface hash,
                           semantic id, provenance)
   ↓
Effect broker            handler registry; records all effect results
   ↓
Workers (Capsules)       isolated processes; FS-sandboxed (AILANG_FS_SANDBOX),
                         containers when isolation demands it
   ↓
External systems         local FS, git, GitHub, models, containers
```

**Capsule** — the executable unit: program, inputs, capabilities, budget, contracts,
dependencies. Capsules execute inside isolated workers and communicate only through typed
effects.

**Storage:** SQLite initially (semantic DB + log + object store all in one file — trivially
back-up-able, trivially replicable later); RocksDB only if measurement demands it.
Persistent immutable data structures in the store; objects never modified, new versions
derive from old.

**Interfaces:** CLI first, REST second, A2UI third (M6), gRPC/LSP optional later.
Everything internally communicates through typed commands; the human/AI split over this
layer is §11.

**Security = two layers.** Semantic: capabilities, budgets, policies, contracts. Physical:
process isolation, FS sandbox, read-only mounts, containers, Linux namespaces; microVMs
later. Neither layer substitutes for the other.

## 16. Initial target domain: software engineering

The first managed world is the one we operate daily: repositories, tasks, evaluations,
worktrees, benchmarks, deployments, capabilities, budgets.

**Motoko becomes the first native agent** — its shell orchestration replaced by world
transitions, its per-edit typecheck feedback becoming Verify-phase evidence, its runs
becoming replayable transition history.

**Boundary with Motoko's deterministic test world (ruling adopted 2026-07-31, per the
[DST overlap note](https://github.com/arniwesth/motoko_agent/blob/arniwesth/mot-44-motoko_dst_execution_primer/.agent/projects/009_motoko_dst_execution/NOTE-ailang-world-overlap.md)
— see REFERENCES.md):** the Motoko DST world is an *ephemeral, seed-driven effect-scenario
interpreter* that tests Motoko's real driver — generated faults, virtual time, shrinking; AILANG
World is the *persistent production substrate* — governance, provenance, storage, replay of what
actually happened. Different layers; neither replaces the other. Integration happens at Motoko's
typed request/result seam (a future `AilangWorldLiveAdapter` delegating live effects to World's
broker), evidence binds by content hash — schemas never merge. Terminology discipline:
`DeterministicTestWorld` vs `AilangWorld`, never a bare shared "World" across the two meanings.

## 17. Milestones

| # | Deliverable | Notes |
|---|---|---|
| M0 | Charter + quorum ratification | World iteration 0, with Mark; this doc is the input. Repo bootstrapped via the portability M3 kit |
| M1 | Semantic world library | World/Proposal/Transition/Evidence types **in AILANG** (`ai-check` green as CI gate); Go host for SQLite store, content-addressed objects, append-only log; replay of recorded transitions proven |
| M2 | Local daemon | `ailang-worldd`: SQLite, REST API, CLI. Zero cloud deps |
| M3 | Effect broker | FS, Git, Model (`std/ai`), `Human.Approve` (reuse approval-queue pattern); effect-result recording |
| M4 | Reference-agent integration — **the non-inferiority floor (§12.4 part 1)** | First end-to-end propose/verify/commit by agents. Dual reference (Claude Code + codex), shell arm vs World-MCP arm, both must hold non-inferiority at bounded overhead (stability precondition on the shell arm); motoko optional third arm. Floor fails on eligible agents → World parks. The VALUE burden lives in clause 5's provenance teeth + R1 |
| M5 | Speculative execution | Parallel candidate branches over immutable worlds; evaluation graph; evidence comparison |
| M6 | Generated UI (A2UI/AG-UI) | Dynamic interfaces generated from world state, available transitions, approvals, evidence — emitted in an open agent-UI protocol (open question 9) |
| M7 | Multi-agent coordination | Planner / verifier / reviewer / scheduler / operator communicating through typed proposals + messages (nodes/edges live); external agents interop via A2A — aitana/platform as the first cross-stack peer |
| M8 | Cloud + multi-node extension | Pub/Sub transport, cloud cascade, log/object replication → the multi-laptop / home-server World. Extension handlers only; core untouched |

Delivery discipline: each milestone lands via the World mission loop (the parameterized
driver — never a hand-rolled second loop), with the `ailang-code` verify profile
(`ailang check` / `ailang test` / `ailang ai-check`).

## 18. Non-goals (v1)

- Not a Linux kernel, desktop OS, filesystem, network stack, container runtime, or compiler.
- Not a second mission loop — World's development reuses the parameterized mission
  infrastructure.
- No Pub/Sub, no cloud service, no hosted anything in the core. Cloud is M8, as handlers.
- No new wire protocols. Where an open protocol exists (MCP, A2A, A2UI/AG-UI, OTEL), World
  speaks it; a homegrown dialect requires demand evidence that no open protocol fits.
- Not a public "semantic OS framework" product. Internal first; productizing is post-v1.

## 19. Open questions for ratification (iteration 0)

1. **How much kernel lives in AILANG vs Go?** Recommendation: transition *semantics* (types,
   pure transition functions, contracts, policies) in `.ail` from day 1 — verified by
   `ai-check`; the daemon (store, broker, scheduling runtime) in Go, like the compiler. The
   ratio shifts toward AILANG as the language grows the needed ergonomics.
2. **Capsule isolation v1**: plain processes + `AILANG_FS_SANDBOX` (cheap, today) vs
   containers from the start? Recommendation: processes + sandbox for M1–M4, containers at
   M5 when speculative branches multiply.
3. **World repo name + slug** (portability M3 needs it to seed state files and the
   bookkeeping issue).
4. **Message payload typing**: extend `ailang messages` with typed payloads, or model typed
   edges purely in the world store with messages as transport? (Leans: world store owns
   types, messages stay dumb pipes.)
5. **Trace identity**: adopt OTEL trace/span IDs as transition identifiers, or mint
   World-native IDs with OTEL binding? (Leans: World-native, OTEL-bound.)
6. **Quorum threshold for self-mod class**: which self-modifications require multi-model
   quorum evidence vs `ai-check`-only? Needs a policy table at M0.
7. **Human UI seed**: evolve the Collaboration Hub, or build the five standing surfaces
   (§11.2) fresh in the World repo with the Hub as a pattern reference? (Leans: fresh —
   the Hub is coupled to this repo's coordinator schema, and World's surfaces should be
   projections from day 1.)
8. **M4 value-gate thresholds**: define "match-or-beat" and "acceptable overhead"
   numerically before M4 starts (candidate: pass-rate within noise on the standard
   benchmark tier, ≤25% wall-clock overhead) — set at M0 so the gate can't be argued
   after the fact.
9. **Agent-UI protocol dialect**: A2UI and AG-UI are both young — standardize on one for
   generated projections, or emit a common core with per-client adapters? Related: how
   much of the A2A task lifecycle to adopt at M7 vs agent-card-only interop first.
   Interop test target either way: an aitana/platform agent drives a World transition
   end-to-end over MCP/A2A.
10. **Calibration-record scope** (§13.3): reputation per model, per agent-instance, or per
    (model, harness, capability-class) tuple? Per-model is unfair to well-harnessed
    instances; per-instance fragments the signal. Also: decay policy (models improve —
    stale reputation must fade), and whether records transfer when an agent migrates
    harnesses.
11. **Log-epoch semantics versioning**: old transitions must replay identically forever,
    which requires pinning the evaluation semantics each log entry was written under
    (Urbit froze its kernel language; Temporal needed patching APIs — see
    [REFERENCES.md](REFERENCES.md)). Per-epoch interpreter version in the log header?
    Content-addressed transition functions (Unison-style) so "which code ran" is never
    ambiguous? Decide before M1 freezes the log format — this is the hardest thing to
    retrofit.

## 20. Long-term vision

AILANG World becomes a semantic substrate where autonomous agents propose changes,
deterministic systems verify them, explicit capabilities authorize them, and transitions
become replayable history — software engineering as continuous semantic evolution rather
than imperative scripting. And because the whole system is one typed world graph, the same
substrate that manages deployments can host documents, datasets, evaluations — or a game —
as just another projection.

The underlying operating system executes bytes.
**AILANG World executes intent.**

---
*Document lineage: v0.1 (Mark, 2026-07-23, conversational draft) → v0.2 (this doc —
thesis elevated to §1, local-first principle added, substrate inventory added,
human/AI interface model added (§11), use cases + falsifiable value thesis added (§12,
M4 gate), self-modification model added, protocol-native principle added (§3.7 —
MCP/A2A/A2UI/AG-UI, aitana/platform composability), milestones grounded against shipped
`ailang` machinery, type sketches compiler-validated) → v0.3 (adoption & ambition §13 —
progressive citizenship, the six guarantees, calibration records, the R1–R7 ambition
ladder; agent constituency statement linked at §12.3; §13.5 trend-extrapolation
priorities — verification throughput first, kernel performance budget from day 1).*

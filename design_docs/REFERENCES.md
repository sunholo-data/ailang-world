# AILANG World — Prior Art & References

*Compiled 2026-07-23 (attended, Mark + Fable session). Every URL below was fetched and
verified live on that date. Companion to [DESIGN.md](DESIGN.md). Four clusters: immutable
state stores · capability security & effects · local-first & provenance · agent
infrastructure & protocols.*

## Design deltas these references suggest (read this first)

Lessons strong enough that iteration 0 should act on them:

1. **Version the interpreter semantics per log epoch** (Urbit froze Nock; Temporal's
   patching pain) — old transitions must replay identically forever, so the evaluation
   semantics a log entry was written under must be pinned. → added as open question 11.
2. **Content-address the transition functions themselves** (Unison) — the log references
   code by hash, so "which version of the function ran" is never ambiguous and
   parse/typecheck results cache forever.
3. **Decide input-addressed vs content-addressed** for evidence/caching explicitly (Nix's
   core distinction).
4. **Design the retention/excision story up front** (Datomic had to bolt on excision for
   GDPR; append-only grows forever — Automerge/Jujutsu hit the same wall).
5. **Make the hash algorithm explicit and versioned in the object encoding** (Git's
   decade-long SHA-1→SHA-256 migration).
6. **Choose the revocation model early**: structural/recursive (seL4 CDT) vs interposed
   (Miller's caretaker); Zircon shows retrofitting either is painful.
7. **Capsule boundaries need an explicit effect-forwarding protocol** — handler-stack
   semantics don't cross process boundaries (OCaml 5's hard lesson).
8. **The verifier is an attack surface** (process-supervision literature): agents will
   Goodhart the checker; calibration records should track verifier-gaming distinctly.
9. **Grade evidence maturity like SLSA levels** — self-reported vs service-generated vs
   tamper-evident evidence should be distinguishable, and policy set per level; a
   verification regime's *strength* should weight calibration records (AgentReputation).
10. **Compile declared capabilities down to an OS sandbox profile** (Anthropic
    sandbox-runtime) — the two-layer security model made concrete: a proposal's typed
    capability set generates its physical confinement.
11. **One ambient builtin collapses the whole guarantee** (ocap all-or-nothing) — the
    stdlib must be audited so no function grants authority without a capability parameter.

---

## 1. Immutable state stores & transaction logs

### Datomic — the database as an immutable value, transactions as data
https://docs.datomic.com/datomic-overview.html
A Datomic database is a set of immutable atomic facts; every datom carries its transaction
time, so any query runs unmodified against an `as-of` past state. Most stealable idea:
**transactions are reified as first-class entities in the same graph** — exactly the shape
World's Evidence should take: attached to the transition entity, queryable with the same
language as the world itself. *Pitfall:* single-writer transactor caps write throughput;
unbounded history forced *excision* (GDPR) as a retrofit.

### XTDB — immutability alone is not enough; time is two axes
https://xtdb.com/
An append-only log conflates *when a fact was recorded* with *when it was true*.
Late-arriving corrections ("what did we believe at 5pm Friday?") need bitemporality; if
World keeps only transition-log time, corrections must be modeled as new transitions
*about the past* — design that semantics, don't discover it. *Pitfall:* bitemporality
roughly doubles index/query complexity; XTDB's v1→v2 query-language rewrite shows the cost
of choosing the query surface late.

### Unison — content-addressed code
https://www.unison-lang.org/docs/the-big-idea/
Every definition is identified by a hash of its syntax tree; names are metadata, renames
non-breaking, typecheck results cached forever. Storing World's *transition functions*
content-addressed makes replay exact — the log references code by hash. *Pitfall:*
content-addressed code breaks plain-text tooling (Unison had to build UCM); hash churn
needs deliberate refactoring UX.

### Nix — the purely functional deployment model (Dolstra thesis)
https://edolstra.github.io/pubs/phd-thesis.pdf
Store paths named by hash of *all inputs* give complete dependency graphs, atomic
rollback, zero version interference — an immutable store + pure build function is
`World + Command -> World'` run in production for two decades. Key decision to inherit:
**input-addressed vs content-addressed**. *Pitfall:* the purity boundary leaks
(fixed-output derivations, impure fetches); immutable stores need GC with explicit roots.

### Irmin — branching and merging as first-class store operations
https://irmin.org/
Git-format mergeable/branchable stores with **typed, user-supplied 3-way merge functions
per data type** — agents fork the world, work speculatively, merge with type-directed
semantics instead of last-write-wins. *Pitfall:* per-type merge semantics are hard
(associativity/commutativity is CRDT territory).

### Event Sourcing — Martin Fowler
https://martinfowler.com/eaaDev/EventSourcing.html
Complete rebuild, temporal query, retroactive correction map 1:1 onto World's replay
story. Most valuable: the **gateway pattern** — during replay, side-effecting interactions
are suppressed and prior external responses replayed from the log; AILANG's typed effects
make this cleaner than anyone (replay = evaluate with handlers stubbed to the log).
*Pitfall:* schema/code evolution is the unsolved tax — version the transition schema from
day one; plan snapshots.

### Urbit / Arvo — an OS whose state is a pure function of its event log
https://docs.urbit.org/urbit-os/kernel/arvo
The purest existing implementation of the target architecture: `T: (State, Input) ->
(State, Output)`, log persisted before effects apply. Crucial lesson: replay stability
across upgrades requires **frozen or strictly versioned core evaluation semantics**.
"Solid-state interpreter" is good naming to steal. *Pitfall:* naive full-log replay
doesn't scale (snapshots needed); a frozen kernel makes perf work brutal.

### Ethereum Yellow Paper — the cleanest formalism for world + transaction → world'
https://ethereum.github.io/yellowpaper/paper.pdf
`σ' = Υ(σ, T)` is exactly World's transition, specified with full rigor. Two stealable
mechanisms: the **Merkle-Patricia trie** (one root hash commits to the entire world state
— a verifiable world checkpoint for Evidence) and **gas** as metered bounds on transition
computation. *Pitfall:* take the formalism, not the consensus baggage; MPT state growth
became a severe bottleneck — plan for world-graph growth.

### Git's object model — the minimal viable content-addressed store
https://git-scm.com/book/en/v2/Git-Internals-Git-Objects
Three object types keyed by content hash suffice for a production store; the load-bearing
trick is **structural sharing** — unchanged subtrees reused across snapshots, so immutable
world snapshots cost O(delta), not O(world). Mandatory for transition-per-command.
*Pitfall:* SHA-1 was baked into identity; the SHA-256 migration has taken a decade —
version the hash algorithm in the encoding from day one.

### Jujutsu — the operation log: version the versioning
https://docs.jj-vcs.dev/latest/operation-log/
Every *modification of repo state* is itself an operation object, giving universal undo
and lock-free concurrency (divergences surface as reconcilable conflicts). Meta-lesson:
log operations *on* the world — including admin/meta actions — through the same mechanism
as domain transitions, so undo and audit are uniform at every level. *Pitfall:* op logs
grow unboundedly without a pruning story.

## 2. Capability security, isolation & effect systems

### Robust Composition — Mark S. Miller (PhD thesis)
https://papers.agoric.com/assets/pdf/papers/robust-composition.pdf
The foundational text: capabilities *are* references, authority flows only along the
reference graph, POLA achieved structurally. Steal the patterns: **revocable forwarder /
caretaker** (revocation = interposed proxy), **facets** (attenuation = narrower
interface), **membranes** (transitive attenuation at a boundary — the model for what a
capsule may re-delegate). *Pitfall:* all-or-nothing — one ambient escape hatch collapses
the guarantee; audit every builtin.

### Capability Myths Demolished — Miller, Yee & Shapiro
https://papers.agoric.com/papers/capability-myths-demolished/full-text/
Pre-answers the objections every design review will raise: capabilities ≠ transposed ACLs;
ocap systems *can* confine delegation; capabilities *can* be revoked. Keeps World in the
"references" quadrant rather than the weaker "unforgeable keys" quadrant. *Pitfall:* if
capabilities are serializable tokens smuggleable out-of-band, you've built the keys
quadrant and the myths become real bugs.

### seL4 — capability microkernel, formally verified
https://sel4.systems/About/seL4-whitepaper.pdf
Existence proof that zero-ambient-authority scales to a whole OS. Steal: the **capability
derivation tree** (recursive revocation — revoke a parent, everything derived dies),
untyped-memory retyping (resources enter as explicitly budgeted capabilities), and MCS
**scheduling-context capabilities** — CPU time itself as a capability, the closest
precedent for World's time/budget-limited grants. *Pitfall:* seL4 is mechanism-only;
every delegation policy must be built above it.

### Zircon handles (Fuchsia)
https://fuchsia.dev/fuchsia-src/concepts/kernel/handles
The cleanest modern *API shape*: handle = object reference + **rights bitmask**;
`duplicate_with(rights_subset)` is attenuation as a one-call primitive; delegation is
writing the handle into a message. Copy this shape for World capability values.
*Pitfall:* no deep revocation — decide structural vs interposed revocation early.

### Capsicum — FreeBSD capability mode
https://man.freebsd.org/cgi/man.cgi?query=capsicum&sektion=4
`cap_enter()` is a **one-way door**: global namespaces vanish, only delegated descriptors
work. The primitive to steal for capsule startup — irreversibly drop into capability mode
before running untrusted plan steps. *Pitfall:* adoption stalled because compatibility
fell on every app; a runtime that owns its stdlib avoids this *if the stdlib is audited*.

### WASI + WebAssembly Component Model
https://github.com/WebAssembly/WASI · https://component-model.bytecodealliance.org/
A capability-based syscall surface where a component gets *nothing* it doesn't import.
Steal the **WIT "world"** concept (they share our word!): a typed manifest of exactly
which interfaces a component may import/export — the serialization shape for capsule
manifests; AILANG effect rows are semantically world-shaped. Preopened directories =
standard attenuated-FS pattern. Wasm components are a credible capsule *implementation
substrate*, not just prior art. *Pitfall:* preopens had path-traversal escapes —
attenuation logic is security-critical code; pin WASI versions per capsule.

### Cap'n Proto RPC — capabilities between isolated workers
https://capnproto.org/rpc.html
Interface references that designate *and* authorize, passed over connections. Steal
**promise pipelining** (call methods on a not-yet-resolved capability — delegation chains
cost one round trip, essential once every effect goes through a broker). *Pitfall:*
attenuation/revocation are not protocol primitives and capabilities die with the
connection — World should make them first-class rather than inherit the gap.

### Deno permissions — the adoption-friendly end
https://docs.deno.com/runtime/fundamentals/security/
The best UX study: secure-by-default, `--allow-read=./data` scoped grants, deny-overrides-
allow, interactive prompting, runtime query/self-revoke API. Copy the grant *grammar* and
the prompt-escalation path for `Human.Approve` flows. *Pitfall:* their own docs admit
`--allow-run`/`--allow-ffi` bypass the sandbox — runtime checks are insufficient for
effects that leave the runtime; Process/FFI-class effects need a real isolation floor.

### Koka — typed algebraic effects worked out to the end
https://koka-lang.github.io/koka/doc/book.html
Row-polymorphic effect types with handlers that discharge effects from the row. The idea
to steal is the unification World reaches for: **installing a handler IS granting a
capability** — the effect row is static demand, the handler is dynamic supply, the type
system proves demand met. Study how effect polymorphism stays ergonomic. *Pitfall:*
effect types say *what*, never *how much* — budgets/expiry live in the handler layer,
indexed by the row, not forced into the types.

### OCaml 5 effect handlers — the cautionary counterpoint
https://ocaml.org/manual/5.3/effects.html
Shipped effect handlers *without* effect typing: unhandled effects are runtime exceptions
— directly validating AILANG's typed-rows bet. Copy the pragmatic restriction anyway:
**one-shot linear continuations**. *Pitfall:* effects can't cross async/FFI boundaries —
the identical trap at World's capsule boundary; an effect performed in a worker cannot be
handled by a lexical handler in another process. Capsule boundaries need an explicit
effect-forwarding protocol.

### Isolation floors — gVisor & Firecracker
https://gvisor.dev/docs/architecture_guide/security/ ·
https://github.com/firecracker-microvm/firecracker/blob/main/docs/design.md
gVisor's Sentry is an effect handler at the OS level (every syscall intercepted against a
minimal enumerated surface — a structural rhyme with AILANG handlers). Firecracker: microVM
per workload at ~5 VMs/host-core/second with a **jailer** that acquires privileged
resources then drops privileges — capsule-per-microVM is economically feasible.
*Pitfall:* neither is universal (syscall gaps; no GPUs) — tier isolation per capsule
effect class; resource budgets (RSS/time) are part of isolation (our own eval-OOM
incident).

## 3. Local-first, replication & provenance

### Local-first software — Ink & Switch (Kleppmann et al.)
https://www.inkandswitch.com/essay/local-first/
The founding essay. Steal the **seven ideals** as literal acceptance tests for
`ailang-worldd` (no spinners, offline, multi-device, collaboration, longevity, privacy,
user control). *Pitfall:* its own open-problems list — access control, schema evolution,
unbounded history — is exactly World's hard part; named, not solved.

### Automerge — production-hardened CRDT change log
https://automerge.org/
Every change hash-identified and retained in a compressed columnar format; syncs over any
byte transport. Direct prior art for **an immutable content-addressed change log that
doubles as the sync substrate** — study before designing transition-log replication.
*Pitfall:* generic per-key convergence can't express World's *authorization semantics* —
domain-aware merge or single-writer-per-log discipline needed.

### Electric — the market signal in the pivot
https://electric.ax/
Abandoned full bidirectional local-first sync for **read-path-only sync** with writes
routed through an authorizing backend — which maps cleanly onto "every transition needs an
authorizer"; by 2026 repositioned around durable agent sessions. *Pitfall:* the pivot IS
the lesson — a well-funded team retreated from general two-way sync; don't plan
multi-device sync as an off-the-shelf dependency.

### SLSA — graded supply-chain assurance
https://slsa.dev/
Steal the **assurance ladder** (L0–L3: provenance exists → unforgeable) as the template
for World's evidence-maturity ladder: self-reported vs service-generated vs tamper-evident
evidence, policy set per level. *Pitfall:* SLSA answers "was this tampered with," never
"is this right" — the verification leg has no SLSA equivalent to copy.

### in-toto — attestation framework (CNCF graduated)
https://in-toto.io/
Signed **attestations** (Statement binding artifact digests to a typed predicate) are the
ready-made wire format for per-transition evidence; **layouts** = machine-checkable "who
may perform which step in what order." *Pitfall:* layouts assume rigid pipelines and have
thin adoption — steal the envelope, treat layouts as inspiration.

### Sigstore — keyless signing + transparency logs
https://docs.sigstore.dev/
Fulcio binds short-lived certs to OIDC identities (signatures name a person/workflow, not
a key); **Rekor** is an append-only transparency log — the closest existing thing to
World's transition log lifted to cross-org scale (the R7 substrate). *Pitfall:* for
local-first, verification must work offline — adopt the *bundled proof* pattern, not
online verification.

### W3C PROV-DM — provenance vocabulary
https://www.w3.org/TR/prov-dm/
Entity/Activity/Agent + `wasGeneratedBy`/`wasDerivedFrom` and — most relevant —
**`actedOnBehalfOf`**: the exact shape of "AI agent acted on behalf of the human who
authorized it." Use as naming vocabulary and graph shape. *Pitfall:* deliberately
underspecified semantics; modest adoption — vocabulary, not executable spec.

### Reproducible Builds
https://reproducible-builds.org/
Determinism as *verification*: re-execution equals audit — a transition whose outputs can
be independently regenerated needs no trusted third party. *Pitfall:* a decade of effort
shows the long tail is brutal — scope reproducibility per-transition with recorded inputs,
never promise it globally.

### FoundationDB — deterministic simulation testing
https://apple.github.io/foundationdb/testing.html
The whole cluster runs single-threaded with simulated network/disk/clock; a seed
determines every interleaving; millions of fault-injected scenarios replay exactly. The
lesson is architectural: **determinism was a day-one constraint on how all code is
written** — AILANG's explicit effects are the language-level version of the same bet.
(FDB alumni founded https://antithesis.com/ to retrofit determinism via hypervisor —
proof that building it in is far cheaper than buying it back.) *Pitfall:* simulation only
catches what the fault model simulates.

### TigerBeetle — deterministic state machine + hash-chained log
https://docs.tigerbeetle.com/concepts/safety/
Deterministic state machine under strict serializability; storage that assumes disk
failure, defended with checksums and **hash-chaining**; the VOPR simulator tortures real
code at 1000× speed. The hash-chained corruption-detecting log over untrusted storage is
directly applicable to a transition log in SQLite on a laptop. *Pitfall:* no permission
system, trusted environment assumed — the determinism layer only; authorization comes
from the in-toto/Sigstore side.

## 4. Agent infrastructure, protocols & durable execution

### MCP — Model Context Protocol (current stable 2025-11-25)
https://modelcontextprotocol.io/specification/2025-11-25
Steal capability negotiation at handshake and the client-side primitives — **sampling**
(server-initiated LLM calls needing user approval), **roots**, **elicitation** — mapping
almost 1:1 onto World's approval-budget flow. MCP admits it "cannot enforce these security
principles at the protocol level" — World's pre-execution verifier is exactly the
enforcement layer MCP leaves as an exercise. *Pitfall:* the 2026-07-28 RC is the largest
revision since launch (stateless core, **Tasks extension** — overlaps World's transition
model); build against version negotiation, not a pinned revision.

### A2A Protocol 1.0 — Linux Foundation
https://a2a-protocol.org/latest/specification/
Steal the **Agent Card** (ready-made agent identity record) and the **task lifecycle
state machine** (SUBMITTED → WORKING → INPUT_REQUIRED/AUTH_REQUIRED → COMPLETED/FAILED/
CANCELED/REJECTED) as proposal-state vocabulary; messages-vs-artifacts separation is the
right cut for replayable records. Governance is real (TSC: AWS, Microsoft, IBM,
Salesforce, SAP). *Pitfall:* A2A tasks are opaque and Agent Cards are self-declared,
unverified claims — World's typed proposal is a strict superset; **treat A2A as
transport, not trust**.

### AG-UI (CopilotKit) + Google A2UI
https://docs.ag-ui.com/introduction · https://github.com/google/A2UI
AG-UI standardizes the supervision primitives (interrupts: pause/approve/edit/retry;
thinking-step visualization; typed event-sourced shared state) — steal the event taxonomy
for the approval inbox. A2UI's philosophy — agent-generated UI "**safe like data,
expressive like code**" via declarative JSON from a pre-approved component catalog — is
World's verification-before-expressiveness bet applied to UI; A2UI rides over A2A, AG-UI
integrates, so the three compose. *Pitfall:* A2UI is v0.9.x preview; AG-UI approval
events are UI facts, not proofs.

### Temporal — determinism constraints & code evolution
https://docs.temporal.io/workflow-definition
The most relevant industrial prior art for replayable transitions: no random, no
wall-clock, no I/O in workflow code — AILANG's thesis enforced by *convention*; replay
compares regenerated commands against history and fails loudly. Study **Worker
Versioning** and **Patching** — World faces the identical problem of evolving code while
old histories must replay. *Pitfall:* Temporal detects non-determinism only at replay
time via discipline — a standing argument for doing it statically in the type system.

### LangGraph — checkpointing & time travel
https://docs.langchain.com/oss/python/langgraph/persistence
The most deployed open implementation of replay+fork agent state: thread-scoped
checkpoints at every super-step, resume, human-in-the-loop interrupts, forking from any
prior checkpoint. Steal the checkpoint-as-first-class-object ergonomics. *Pitfall:*
persistence ≠ replayability — arbitrary Python state, no determinism guarantee. World
needs Temporal-grade determinism *plus* LangGraph-grade ergonomics; neither has both.

### AIOS — LLM Agent Operating System (arXiv:2403.16971)
https://arxiv.org/abs/2403.16971
The canonical academic agent-OS: isolates agent scheduling, context/memory management,
tool management, access control into a kernel beneath agent apps. Validates the framing;
gives a checklist of kernel services to provide or explicitly punt. *Pitfall:* it's about
resource efficiency, not safety — no verification, no typed effects, no replay; borrow
the decomposition, not the implementation.

### Karpathy — LLM OS framing
https://www.youtube.com/watch?v=zjkBMFhNj_g
The discourse-defining metaphor (LLM as kernel, context as RAM, tools as peripherals) —
useful positioning vocabulary. World's sharpest differentiator stated against it:
Karpathy's LLM OS is *non-deterministic at the kernel*; World inverts this — a
deterministic verifying kernel with the LLM as an untrusted userspace process.
*Pitfall:* it's an analogy, not a spec.

### Let's Verify Step by Step — process supervision (arXiv:2305.20050)
https://arxiv.org/abs/2305.20050
Rewarding intermediate steps beats outcome-only supervision (PRM800K: 800K step-level
human labels). Direct lesson for the R5 flywheel: **verifier verdicts on plan steps are
automated process-supervision labels at a scale OpenAI paid humans for**. *Pitfall:* when
the reward is an automated verifier, agents reward-hack the verifier — Z3 contracts are
attack surface; calibration records must track verifier-gaming distinctly.

### Anthropic — sandbox-runtime + Building Effective Agents
https://github.com/anthropic-experimental/sandbox-runtime ·
https://www.anthropic.com/engineering/building-effective-agents
`srt`: OS-primitive sandboxing without containers (Seatbelt, bubblewrap+seccomp, WFP) —
deny-then-allow FS, allowlisted egress via proxies. **A proposal's declared capabilities
could compile directly to an srt profile** — the two-layer security model made concrete.
*Pitfall:* per-process resource restriction only; no cross-agent or semantic policy —
that stays World's job.

### E2B — sandbox-as-primitive
https://e2b.dev/docs
The **Sandbox/Template** split: a declarative spec of what environment a sandbox starts
with — a good shape for World's capsules, where required capabilities + effect row compile
to a Template, making the execution environment itself a typed artifact. *Pitfall:* VM
isolation is resource-level, not semantic; egress policy and effect accounting remain the
caller's job.

### SWE-bench — evidence-based verification of agent work
https://www.swebench.com/ · https://arxiv.org/abs/2310.06770
Fail-to-pass tests as a workable evidence record at scale. The meta-lesson: **SWE-bench
Verified had to exist** — human audit found a large fraction of "objective" tasks broken —
i.e. verification harnesses themselves need verification, which is what World's typed
contracts formalize. *Pitfall:* tests underspecify intent; contamination is endemic.

### "Language model harnesses are compositional generalizers" — Alex L. Zhang (2026)
https://alexzhang13.github.io/blog/2026/harness/
The capability-side argument for World's shape (the LBAC landscape covers the control side —
see POSITIONING.md). Thesis: the harness, not the network, should carry the inductive bias —
its job is to encode arbitrarily complex environment state so **every model call is locally
in-distribution**, reducing unfamiliar problems to compositions of familiar ones. Evidence:
RL-training a Recursive Language Model harness generalizes to held-out tasks **8–32× longer**
and across domains at ~10× the eval lift of training the bare Transformer — because the
harness "induces an equivalence relation between tasks with latent similarities" (a quotient
over trajectories: structurally similar tasks become near-isomorphic token streams).
**Steal for World**: (1) "locally in-distribution" as an explicit DESIGN CRITERION for the MCP
projection and every agent-facing surface — stable transition identities, canonical packet
shapes, bounded request digests are not just audit hygiene, they are generalization
infrastructure; (2) the isomorphism framing as the strongest form of the clause-4 hypothesis —
the World arm may beat shell *because of* its typed structure (the decomposition prior), not
despite its overhead; an interpretive lens the `w-agent-floor-m4` design doc should carry;
(3) the "Mismanaged Geniuses Hypothesis" (human tasks decompose into sub-tasks within current
model capability; the decomposition itself is short) — the same bet World's sprint-sized queue
discipline has been winning on for 40 iterations. *Pitfall:* results are one lab's RLM
experiments on one model family; treat as hypothesis-grade for World until the floor
experiment produces our own paired data.

### Motoko DST ⇄ AILANG World boundary note (arniwesth/motoko_agent, 2026-07-24)
https://github.com/arniwesth/motoko_agent/blob/arniwesth/mot-44-motoko_dst_execution_primer/.agent/projects/009_motoko_dst_execution/NOTE-ailang-world-overlap.md
Companion to that repo's `ADR-001-deterministic-test-world-architecture.md`; reviewed World at
`03efeef` and reads our code accurately. Its ruling — **adopted by World**: the Motoko
deterministic test world (ephemeral, seed-driven, fault-generating, virtual-time) and AILANG
World (persistent governance/provenance substrate) share transition discipline but are
DIFFERENT LAYERS; neither replaces the other. World's replay reconstructs *what happened*;
DST generates *what could happen* — the FDB-simulation vs Datomic-history split from this
bibliography, made concrete. Steal/align: the **typed effect envelope** field list (kind+origin,
causal identity + encounter ordinal, bounded request digest, deadline/timing, **ordered
intermediate emissions**, typed result, capability identity, evidence refs); the
**partial-stream-then-error** requirement (terminal-only result records are insufficient — this
names a REAL World gap for streaming `Model.Infer` records, converging with upstream
ailang recorded-stream work); the **terminology rule** (`DeterministicTestWorld` vs
`AilangWorld`; never a bare shared `World`). Integration shape: a future
`AilangWorldLiveAdapter` behind Motoko's typed request/result seam, delegating live effects to
World's broker while Motoko keeps its own driver behavior. *Pitfall/staleness:* reviewed
pre-broker — its revisit-trigger 1 ("review when the M3 broker contract exists") has FIRED
(`host/broker` landed, items 4+4c complete); the scheduled boundary review is now due on both
sides, natural venue = the upstream recorded-stream lane which already carries both DST ADRs
as design context.

### Cordis spatiotemporal composability ⇄ Motoko comparison (arniwesth/motoko_agent, 2026-08-14)
https://github.com/arniwesth/motoko_agent/blob/e3f041919f9d17064b4081cef92be54f32fb1cef/.agent/research/Spatiotemporal_Composability/cordis-paper-vs-motoko.md
Motoko's read of "A Programming Paradigm for Spatiotemporal Composability" (Shi/Zhang/Cui,
PKU/DeepSeek-AI, Aug 2026, 88 pp; Cordis meta-framework, TypeScript, validated observationally
on Koishi). The paper lifts effect/coeffect theory to runtime: **temporal composability** =
revertible effects (`Γ → Γ × (Γ → Γ)` — every context transformation returns its inverse;
unload = LIFO-composed undo), **spatial composability** = reactive coeffects (declared
dependency sets driving activation/teardown; realm isolation; right-biased monoid-merged
interception metadata for capability attenuation). Motoko's doc frames the two as opposite
bets — revert-after (Cordis) vs verify-before (Motoko DST + Phoenix restart); World is the
unstated THIRD bet on the same axis: never mutate — append + re-select over immutable
content-addressed history. **Steal/align:** (1) inverse-carrying effect records at the broker —
a `RecordedEffect` that also records its compensator, *appended as data, never applied to the
log*, making compensation one more replayable transition; (2) the evidence-grade reading of the
paper's central obligation — `g∘f = id` is CLAIMED under Cordis (author obligation, its §5.1.1),
TESTED under Motoko's revert-then-replay DST family, PROVEN for World's encodable pure fragment
— the three systems are rungs of item 13's ladder, not competitors; (3) instantiation-as-
tracked-effect (its Def. 47) is already native here — a Proposal IS one; (4) the interception
algebra maps onto approval-scope/capability attenuation. **Do NOT adopt the runtime:**
hot-plug/reactive topology is a spatial-composability surface whose cost World measured
first-party at iter-84 (exports are modules; a foreign module minted `PROVEN`) — Cordis
*mediates* access at runtime, World *relocates* authority behind the host seal; and World
closes the paper's two admitted gaps by construction (nominal linking with no versioning →
content addressing, its §6.6; boundary emissions only withheld/compensated → the approval
scopes live exactly on that line, its §6.1). *Pitfall:* the paper's harness application is
future work and its Koishi evidence is "existence-and-adoption", not quantitative; and **M4
expires item 17 §2.4's "no cross-host proof receipt" scoping** — a Motoko DST verdict arriving
as Evidence is CLAIMED until a cross-producer provenance story exists (re-execution is
economically absurd at simulation scale; the ratified single-host MAC does not extend across
the trust boundary to a foreign producer's key). Natural venue: the due Motoko⇄World boundary
review — see the entry above, whose revisit-trigger 1 has FIRED.

### AgentReputation — decentralized agent reputation (arXiv:2605.00073)
https://arxiv.org/html/2605.00073v1
Closest published prior art to §13.3 calibration records: **verification regimes of
quantified strength attached to reputation metadata** (an automated check ≠ an expert
review), context-conditioned reputation cards (competence doesn't transfer across
domains), risk-based verification escalation. Steal all three — regime strength maps to
how much a Z3-verified vs type-checked vs human-approved transition moves calibration.
*Pitfall:* recent preprint, unvalidated at scale — take the schema, skip the blockchain.

---
*43 references verified live 2026-07-23. Update discipline: when a reference materially
changes a design decision, record the delta in DESIGN.md and cite the reference; when one
goes stale, replace it here with the successor.*

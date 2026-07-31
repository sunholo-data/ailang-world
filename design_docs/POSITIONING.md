# AILANG World — Positioning & Landscape (v1.0)

**Status**: adopted 2026-07-31 (attended: Mark supplied the landscape analysis; coordinator
spot-verified the Tier-1 anchors before adoption per the premise discipline). Living document —
market-facing positioning, distinct from [DESIGN.md](DESIGN.md) (architecture) and
[REFERENCES.md](REFERENCES.md) (prior art we build on). Update discipline: claims about
competitors carry a verification flag; an unverified survey row is a claim, not a fact.

## The category has a name: Language-Based Agent Control (LBAC)

As of mid-2026, "agents express intentions as code that must type-check against effect/capability
constraints, and unsafe programs are rejected BEFORE execution" is a recognized research
tradition, not a lone bet. Verified anchors (2026-07-31):

- **Odersky et al., "Securing Agents With Tracked Capabilities"** — **Best Paper, ACM CAIS 2026**;
  project **TACIT** (Tracked Agent Capabilities In Types); Scala 3 capture checking statically
  tracks what an agent may touch, including proving sub-computations pure. His argument against
  boundary-only enforcement (a policy language bolted beside an unchecked tool interface = two
  descriptions of the same interaction that must be kept in sync) is our §12.2 argument, made by
  the creator of Scala. **[VERIFIED: arXiv 2603.00991, ACM DOI 10.1145/3786335.3813127, CAIS
  program, Substack post "Tracked Capabilities for Safer Agents" July 2026]**
- **Etas: An Effect-Typed Language for Agent Systems** (HKUST, arXiv **2607.17780**, 2026-07-20) —
  the closest published design to World's language leg: model calls, tools, prompts, typed
  memory, human approvals, policies, and execution traces as *semantic program elements*; effect
  rows + a persistent abstraction of the typed action trace. **[VERIFIED: arXiv abstract]**
- **LBAC / TypeGuard** (Zhou et al., May 2026) — static conformance of all agent-generated code
  (recursive sub-agents included) to effect/data-flow policies, closed under composition;
  prompt-injection resistance at comparable utility. **[UNVERIFIED-INDIVIDUALLY — Mark-supplied
  survey row; density of the tier independently confirmed (see also SIGIL, arXiv 2607.27309)]**

**Consequences accepted:** (1) strong validation — several groups arrived at World's thesis
independently within ~90 days; (2) any "nobody else verifies in the language" phrasing is now
FALSE and has been rewritten (README, DESIGN §12.2) — positioning against named peers beats
claiming an empty field, and being caught claiming disproven novelty is worse than having
competitors.

## The tier map (condensed; Mark-supplied 2026-07-31, spot-verified at Tier 1)

| Tier | What checks, where | Occupants | World's relation |
|---|---|---|---|
| 1 — in the language, pre-execution | type/effect/capability systems reject unsafe agent programs before they run | TACIT/Odersky · Etas · LBAC/TypeGuard · SIGIL | **Peers on the language leg.** All are research prototypes, papers, or embeddings in a host language |
| 2 — at the wire, runtime | policy/boundary guards on proposed tool calls | Bedrock AgentCore+Cedar · MS Agent Governance Toolkit · AgentSpec · Progent/Conseca/PCAS · CaMeL · AgentArmor · MI9 · Aegis · SAFi · VIGIL | The crowded tier World is *structurally* differentiated from (Odersky's argument, ours since §12.2) |
| 3 — hardware attestation | TEEs, cryptographic runtime receipts | EQTY · AgenTEE · EHV · Notarized Agents | Complementary floor, not competing semantics |
| 4 — sovereignty kernels & verifiable history | Merkle logs, proof-gated execution, personal-hardware framing | PunkGo ("Right to History") · DTF/Verifiable Agentic Infrastructure · Sovereign Agents | **Peers on the substrate leg** — history/authority without the language. PunkGo's energy budget is structurally adjacent to our budgets; its EU AI Act Art. 12/14 framing overlaps our compliance story |
| 5 — agent operating systems | kernel-shaped agent runtimes | AIOS · AEROS · Fiserv agentOS · Agent libOS | Category cousins; none verify in a language, none carry a proven-replay ledger |

## World's actual position: the intersection nobody else occupies

The landscape's own structure makes the differentiated claim precise:

1. **Tier 1 has the language and no substrate.** No persistent world graph, no provenance
   ledger, no content-addressed replay-from-store, no approval economics, no calibration
   records, no self-modification governance. Papers and embeddings, not places where work
   lives. **Tier 4 has the substrate and no language.** Proof objects and Merkle history over
   opaque actions — nothing proves the *plan* before it runs.
   **World is the only project in the map holding both legs**: an LBAC language (purpose-built,
   shipped compiler, Z3 contracts, effect rows) *operating* a persistent governed substrate
   (immutable world graph, effect broker with receipts, bit-for-bit replay, budgeted human
   attention) — built by an autonomous mission whose every landing is itself evidence.
2. **Running code beats position papers.** A shipped language + product results upstream
   (Parse; measured benchmark deltas), and in this repo: a daemon, a broker with 7/7 Z3-proven
   decision law, crash-injected durability proofs, and 40+ iterations of publicly-logged
   autonomous construction. Nothing in Tier 1 runs anything comparable.
3. **The decision taxonomy is the unclaimed wedge.** Everyone in the map does *authority over
   effects* (what may the agent touch). Nobody does **decisions as governed objects**:
   permitted ambiguity, designated resolver, collapse deadline, **decision budget**. World
   already carries the embryo: `Human.Approve` as a budgeted effect (§8, HUMAN-SURFACE P5 —
   "how many times may this workflow interrupt a person"), the *defer* verb, sign-the-type
   discretion boundaries. PunkGo prices energy; nobody prices decisions. **Formalizing the
   decision taxonomy as typed World objects is a candidate design item** (post-1.0 lane or
   pulled forward on commercial demand) — and per the analysis, the strongest differentiator
   is the one sellable soonest; it also maps to EU AI Act Article 14 human oversight, where
   the compliance framing will crowd fast.

## Standing positioning sentence (use everywhere)

> Odersky is doing this in Scala; Etas did it as a paper. AILANG World is a purpose-built
> language with a working compiler, Z3-verified contracts in CI, a running governed substrate
> with proven bit-for-bit replay — and it is being built, in public, by the autonomous process
> it exists to govern.

## The second axis: control AND capability (added 2026-07-31, same day)

The tier map above is entirely about **control** — proving what an agent may do. Zhang's
["Language model harnesses are compositional generalizers"](https://alexzhang13.github.io/blog/2026/harness/)
names the orthogonal **capability** axis: a well-designed harness makes every model call
*locally in-distribution*, collapsing structurally similar tasks into near-isomorphic
trajectories — and his RLM experiments show harness-level structure generalizing 8–32× beyond
training length where the bare network fails (see REFERENCES.md).

World's position on this axis is the same discipline wearing its other hat: **typed
transitions, canonical decision packets, and stable tool identities are simultaneously the
proof surface (control) and the normalization surface (capability)**. Verification structure
and generalization structure are the same structure. Nobody in Tier 1 makes the capability
claim; nobody in Zhang's framing makes the control claim; World's architecture happens to be
both — and the clause-4 floor experiment will produce the first paired data on whether the
capability effect is real for production agents. If it is, the floor was conservative: the
honest sales line stops at "no worse" until our own data says "better".

## Risks & open actions

- **Vocabulary risk**: academia will own the category's language within ~a year. Publish the
  plain-English "what is LBAC and why should a CTO care" piece early (unclaimed; companion to
  the entropy-cliff post) or become "an implementation of LBAC" instead of a definer. (Mark's
  lane; recorded here so the loop's docs stay consistent with whatever ships.)
- **Convergence watch**: Etas's "persistent abstraction of the typed action trace" and TACIT's
  purity proofs are ideas to LEARN from, not just position against — candidates for upstream
  AILANG language evolution via the normal lanes.
- **Revisit triggers**: any Tier-1 project ships a runtime/substrate; TACIT lands in a Scala
  release; the decision-taxonomy wedge gets a claimant. Re-verify tiers quarterly.

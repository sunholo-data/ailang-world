# The Human Surface — founding UX design (v0.3)

**Status**: **RATIFIED AND BINDING. §7 was formally ratified by Mark on 2026-07-28** (attended
TRIPLE RATIFICATION — `bc467f1`, *"human-surface v0.3 s7 formally ratified"*), and the charter's
6b queue row was flipped to match on 2026-07-31 (`ee75837`). **This header was not, and read
`PARKED on §7 human ratification` for 14 days — corrected 2026-08-11, attended.** It is the header
a designer reads first, so the stale line said *blocked on a human* about the one thing that was
already decided. Nothing in this document is waiting on a human decision.

Pick-time quorum completed at iteration 26 (2026-07-28): two rounds, both BLOCKED, all four
objections applied — r1 (unverified backend premises · missing Hub conflict surface) in **v0.2** by
the rotation designer; r2 (no bounded wait on a human — both reviewers, independently) in **v0.3**
under the narrow-refinement carve-out, using the reviewers' verbatim text. Full record in §9. The
v0.1 principles were provisionally ratified attended on 2026-07-27 ("good principles for now, yes
good start") and are now covered by the §7 ratification above. Authored attended (Mark + World
coordinator).

**What gates the work is now a DEPENDENCY, not a decision**: `w-approval-inbox` (queue item 7)
unblocks when item 5's MCP projection lands, which is itself gated on item 11's `TR.C` binding
gate. The §7 outputs (total evidence-grade mapping · timeout-policy set · packet-schema freeze
timing) are ratified inputs to that work, not open questions — implement to them.
This is the experience-layer peer of [DESIGN.md](DESIGN.md): §11 answered *what surfaces exist*;
this answers *how humans and the AI state machine actually meet*. **Binding input** for
`w-approval-inbox` (clause 5) and M6 (generated projections) — neither routes to sprint before
this doc is quorum-reviewed and Mark-ratified. Companion: [SCENARIOS.md](SCENARIOS.md).

---

## 1. The premise — what no UI could assume before

Every prior interface renders one of two things: **mutable opaque state** (classic GUIs — the
screen shows a value; where it came from and whether it's true are unknowable) or **token
streams** (chat UIs — prose whose relationship to reality is vibes). The established premise
today is narrower: the store has an immutable append-only log and bit-for-bit replay, verified
by `scripts/verify_go.sh` in CI (§8). Universal provenance and a total evidence-grade mapping
are **not established**: `Store.Commit` permits unreadable references, and the five-variant
`Evidence` ADT has no total mapping to the four proposed grades.

The intended surface is the first that may eventually assume **everything behind the glass is
a typed, verified, replayable ledger**. That is a backend requirement, not a present fact.
**The surface MUST NOT assume universal provenance or replayability until the referenced
ledger schemas and tests establish those properties; unavailable grades or links MUST be
rendered explicitly as unsupported, never silently synthesized.** Until then, each renderer
MUST preserve the narrower guarantees recorded in §8 and expose missing or unreadable evidence.

The design question is therefore NOT "how do we make a nice dashboard" — it is: **what
interaction grammar does verified state make possible that was impossible before?** The answer
below is seven principles. The metric they all serve (DESIGN.md §13.5): **sound decisions per
second of human attention** — attention is the one resource the model-trend never makes cheaper.

## 2. The seven principles (the novel grammar)

### P1 — Decision packets, not messages
The unit of human-AI interaction is not the chat message; it is the **decision packet**:
proposal + evidence bundle + exact authority requested + budget state + reversibility class,
composed so one glance answers "what do I need to know to decide?" Chat carries intent INTO the
system and narrative OUT of it; decisions never live in prose. (The approval-inbox mockup from
the design session is P1 rendered.)

### P2 — Grounded prose
When the AI writes to a human, **every load-bearing noun MUST be a link to a typed world object** —
the deployment it mentions, the evidence it cites, the goal it traces to. Prose that cannot
link is *visibly ungrounded* (rendered distinctly), so unsupported claims are exposed by the
renderer itself. This is the anti-hallucination idiom: not "trust the model," but **text whose
claims are clickably checkable** — and the model knows its ungrounded sentences will show.

### P3 — Trust-gradient rendering
Every fact on screen carries its **epistemic grade**, visually: `PROVEN` (Z3/replay — machine
certainty) > `TESTED` (named checks passed) > `ATTESTED` (recorded effect results) >
`CLAIMED` (agent said so, unverified). One consistent visual channel (weight/badge, never color
alone), everywhere. A human should be able to see at a glance which parts of a screen the
machine *proved* and which parts an agent *asserted* — no UI has ever shown that distinction,
because no backend could supply it.

**Dependency — PARTIAL (`w-approval-inbox`, gated on this document):** `world/types.ail`
defines five evidence variants, but no total grade mapping exists. The gradient is binding
only after ratification point 7.2 supplies that mapping; until then render an unavailable
grade as `UNSUPPORTED`, not as a guessed grade.

### P4 — Time is navigation
When state is backed by the established immutable log, "where am I?" = **(domain × time)**,
not a page. The timeline
scrubber and the **world-diff view** ("show me now vs before that transition") are chrome, not
features — as fundamental as the back button. Scenario 3's provenance walk and clause 5's
5-minute "why" answer are P4 exercised backwards; the R5 replay is P4 exercised forward.

### P5 — Attention choreography
`Human.Approve` MUST become a budgeted effect (DESIGN.md §8), so the surface treats
interruption as spending:
batch over interrupt; rank packets by **irreversibility × novelty**; never ask what
verification settled; deliver ONE digest, not N pings; "review the night" after auto-commits
rather than gating them. Quiet is the default state of a healthy World — the UI must make
quiet feel like health, not absence.

**Dependency — NOT BUILT (`w-effect-broker-m3`, queue item 4, PARKED):** the landed daemon
explicitly has no effect broker or capability/budget authority
(`host/daemon/daemon.go:34`). P5 remains the required interaction policy; enforcement cannot
be claimed until that item lands.

### P6 — Five zoom levels, everywhere
Every projected object MUST support the same ladder: **glance** (chip/count) → **decide** (packet) →
**inspect** (evidence bundle) → **walk** (provenance graph) → **replay** (re-run the episode).
One gesture-grammar for goals, deployments, documents, agents, worlds. Learn it once.

### P7 — Speculation as a gesture
Worlds MUST support near-O(1) forks through immutable structure sharing so **"what if" can be
a primitive**: fork → let agents run against the
branch → compare evidence side-by-side → commit or discard. The compare-arms view (scenario 5)
is a standing screen, not an expert mode. Immutability makes counterfactuals safe enough to
hand to a human as a button.

**Dependency — NOT BUILT (M5 speculation infrastructure; no routed queue item yet):**
`host/store/store.go` has one selected head and no fork, branch, or compare-arms API.
Content addressing makes cheap forks plausible, not implemented.

## 2.5 The medium — a workbench, not a transcript (Mark's question, answered plainly)

**Is this a chat UI? A terminal? Neither is the primary surface.** The ChatGPT-era pattern
makes the *transcript* the workspace: state, decisions, and evidence all flatten into scrolling
prose, and authority accumulates in a conversation log that nobody can query. World inverts
this: **the WORLD is the workspace; conversation is a lens you can open on any object** — ask a
question from a deployment, a goal, a transition, and the answer arrives as grounded prose (P2)
anchored to where you are. The transcript never accumulates authority; the ledger does.

The primary human surface is therefore a **workbench over the world graph** — its closest
ancestors are not chat apps but a code-review queue (P1 packets), a version-control history
browser (P4 time), an ops console (P5 choreography), and a spreadsheet (direct manipulation of
typed values). Concretely, the reference renderer's chrome is: the **decision-packet stream**
(what needs you), the **world browser** (what is), the **timeline scrubber** (what was/why),
with a **conversation lens summonable anywhere**.

Where chat and CLI actually sit:
- **Chat** = one input channel of the goal composer (§3) + the question lens + the *remote*
  rendering (an MCP client like Claude Code MUST eventually render packets and grounded prose
  when you're away from the workbench — scenario 4; §6 and §8 record the current blocker).
  Casual and mobile use may be chat-first; the workbench
  remains the place decisions are *designed* to happen.
- **Terminal/CLI** = plumbing-truth access and scripting for power users; never the designed
  experience, always available (no capability the workbench has may be CLI-impossible, and
  vice versa — one typed layer beneath both).

## 3. The input side — signing the type, not the vibe

The goal composer MUST accept natural language and echo back BOTH renderings: the prose
paraphrase AND the **typed Goal object** (budget, capability grant, contracts, approval gates)
— and the human confirms the *typed* one. What you sign is what the machine enforces; the
prose was only ever the request. Decision verbs beyond approve/reject are first-class:
**attenuate** (approve narrower), **defer** (park for more evidence), **reject-with-reason**
(typed feedback to the agent), **amend-policy** (this class never asks again).

To enforce bounded waits, every Decision Packet and `Human.Approve` effect MUST carry a strict
TTL. If the human does not decide before expiration, the packet deterministically resolves to a
typed rejected-timeout (or a predefined safe fallback), unblocking the agent.

### 3.1 Bounded decision lifecycle (binding principle)

Every decision packet MUST carry a ledger-recorded creation time, decision deadline, and timeout
policy selected from a typed finite set. DEFER MUST create a new bounded deadline and record the
required evidence; it MUST NOT park indefinitely. At deadline, the system emits an explicit
Timeout transition and follows the packet's declared policy, such as cancel, remain safely
unexecuted with bounded escalation, or execute only if authority was already independently
granted. Silence MUST never synthesize approval or rejection. Replay MUST reproduce deadline and
escalation behaviour deterministically from ledger time, without dependence on wall-clock race
ordering.

*Provenance: this subsection and the paragraph above it are the round-2 reviewers' own proposed
fixes, applied verbatim under the charter's narrow-refinement carve-out (§9). Both reviewers —
independently, across two providers — found the same gap: the interaction grammar admitted an
unbounded wait on a human, which is the mission's Standing Rule 6 restated at the UX layer. The
typed finite set of timeout policies is NOT enumerated here; fixing it is part of ratification
point 1 (see §7).*

## 4. The five surfaces, restated as compositions

| Surface (§11.2) | Composition |
|---|---|
| Goal composer | §3 + P2 (the echo is grounded prose over the typed object) |
| Approval inbox | P1 + P3 + P5 (packets, graded evidence, choreographed arrival) |
| World browser | P4 + P6 (scrub, diff, zoom) |
| Provenance explorer | P4 backwards + P2 (every hop is a typed link) |
| Budget/fleet dashboard | P3 + P5 (graded facts, quiet-is-health) |

## 5. Anti-patterns (the surface must never)

- **Streaming prose walls** as the primary channel — prose is narrative, not state.
- **Confidence theater** — percentages/vibes without an evidence link (P2 violation).
- **Unlogged human actions** — every human decision is itself a transition (no admin backdoor).
- **Interrupt spam** — an unbatched notification is a budget bug (P5 violation).
- **Modal interrogation** — blocking the human on questions verification could settle.
- **Grade laundering** — rendering CLAIMED facts in PROVEN clothing is the cardinal sin; the
  trust gradient is load-bearing or it is nothing.
- **Silently synthesizing an unavailable grade** — an unmapped evidence variant or absent
  proof carrier is rendered `UNSUPPORTED`; choosing the nearest-looking grade is grade
  laundering.

## 6. Build strategy (unchanged from §11/§3.7, now sequenced)

Projections MUST emit an open agent-UI protocol (A2UI/AG-UI — dialect = open question 9,
decided at M6). ONE hand-built reference renderer is required: the **workbench**, ultimately a
localhost web surface served by worldd — browser-opened, cross-platform, zero install beyond
the binary. Chat clients (MCP) MUST eventually receive P1/P2 degraded-but-consistent: packets
as structured tool results, grounded prose as resource links. Phone = same packets, smaller
(scenario 4).

**Dependencies — NOT BUILT / BLOCKED:** worldd currently serves eight JSON REST routes and no
HTML, web, or SSE surface (`host/daemon/daemon.go:355-362`); the localhost renderer therefore
belongs to future implementation. MCP/UI projection is queue item 5, `w-mcp-projection`,
**BLOCKED** on its named prerequisites: the pinned v0.30.0 MCP server always discovers
`submit_feedback`, an ambient public-feedback/Pub/Sub egress, and offers no per-session
discovery filtering (upstream `sunholo-data/ailang#498`). No renderer may claim MCP packet or
grounded-link delivery until an exact, session-filtered protocol surface is verified.

### 6.1 Conflict Surface: Hub vs Local Workbench

DESIGN.md:141 says the out-of-repo Coordinator + approval queue + Collaboration Hub already
implement propose → human-approve → execute over DB rows. DESIGN.md:373-375 names the Hub as a
React approval queue + message center and an existing substrate. DESIGN.md open question 7
(`:728-730`) leans **FRESH** because that Hub is coupled to the `sunholo-data/ailang`
coordinator schema, whereas World proposals are typed world objects and its surfaces must be
projections over a content-addressed log.

No Hub source exists in this repository (repo-wide search outside `design_docs/` found no Hub
implementation), and the named source repository is unavailable in this sandbox. Its
authentication, routing, notification, dependency, and deployment internals are therefore
**UNVERIFIED**; this document does not claim they are cloud-coupled or unusable. Closing that
gap requires inspecting the Hub package boundary and build graph in `sunholo-data/ailang`,
then testing whether its renderer can consume World protocol packets without its coordinator
DB schema or any clause-2 cloud dependency.

Pending that inspection, choose **FRESH code, Hub as pattern donor**. There is no local import
path to reuse here; the known data boundary differs (coordinator DB rows versus typed ledger
objects); and the workbench must render evidence grades, content-addressed links, time, and
replay while preserving clause 2's zero-cloud core. Reuse the Hub's approval-queue and message-
center interaction patterns. Reconsider code reuse only if the out-of-repo inspection proves
a dependency-free renderer boundary compatible with World's protocol and authority model.

**Desktop**: a thin shell (e.g. Tauri) wrapping the SAME localhost surface is the sanctioned
path when wanted — it is packaging, not architecture, and it earns its keep on exactly one
principle: **P5**. A menubar/tray presence with a packet-count badge and native notifications
is a better attention-choreography instrument than a browser tab (glanceable quiet-is-health,
OS-level batched alerts). Post-daemon, pullable early on demand; a fully native macOS app
remains an optional post-1.0 renderer, never the definition.

## 6.5 Mockups (design fixtures, committed like the compiled sketches)

Open locally in any browser (self-contained, no dependencies, light/dark aware):

- [mockups/approval-inbox.html](mockups/approval-inbox.html) — **P1 + P5 rendered**: the
  decision-packet stream — evidence bundle, exact authority ask, budget, reversibility, the
  five decision verbs, and the quiet-is-health footer (auto-commits reviewable, never gating).
- [mockups/grounded-prose.html](mockups/grounded-prose.html) — **P2 + P3 + P6 rendered**: AI
  narrative whose nouns are graded typed links (PROVEN/TESTED/ATTESTED/CLAIMED), an ungrounded
  span exposed by the renderer, and the world panel one zoom level deeper with walk/diff/replay.

These are fixtures, not pixel specs: the reference renderer must preserve their *grammar*
(grades visible, packets one-glance-decidable, verbs present), not their styling.

## 7. Ratification points (Mark)

1. The seven principles + anti-patterns + the §3.1 bounded decision lifecycle as BINDING on all
   human-facing work (evaluator-scored, like coding-standards.md). §3.1 requires a **typed finite
   set of timeout policies**; that set is deliberately NOT enumerated in this document and must be
   fixed as part of this ratification (candidate members named by the reviewers: cancel · remain
   safely unexecuted with bounded escalation · execute only if authority was already independently
   granted).
2. **Total evidence-grade mapping.** The kernel has
   `CompilerOutput(HashRef)`, `TestReport(HashRef, bool)`, `HumanApproval(HashRef)`,
   `AiReview(HashRef, float)`, and `RecordedEffect(HashRef)`. The proposed gradient has
   `PROVEN > TESTED > ATTESTED > CLAIMED`. Under this document's definitions,
   `TestReport → TESTED`, `RecordedEffect → ATTESTED`, and `AiReview → CLAIMED`;
   `CompilerOutput` and `HumanApproval` have no grade, while the stated `PROVEN` producers
   (Z3 proof and replay) have no `Evidence` carrier. Ratification MUST produce a **total**
   mapping. Neutral options: add grades; add/reshape `Evidence` variants; or define a
   documented total function with an explicit lowest-grade default. **Recommendation:
   add/reshape evidence variants**, because a carrier should distinguish a compiler result
   from a verified proof and preserve human ratification without mislabelling it as an
   unverified agent claim. This recommendation is not the decision.
3. Decision-packet schema freeze timing (it becomes a world type — kernel-adjacent).
4. Open questions 7 (Hub vs fresh: recommend FRESH, Hub as pattern donor) and 9 (dialect —
   defer to M6 with a common-core emitter).
5. **Proposal confidence.** `Proposal.confidence: float` has no adjacent evidence `HashRef`;
   `AiReview(HashRef, float)` does. Is bare proposal confidence renderable at all, and if so
   under which ratified grade? Until answered, P2 and the confidence-theater anti-pattern
   require renderers to omit it or label it `UNSUPPORTED`, never present it as evidence.

## 8. Premise verification and conflict surface

Status vocabulary in this table is fixed: **BUILT / PARTIAL / NOT BUILT / BLOCKED /
UNVERIFIED**. Controller-supplied observations are identified as such; local commands below
were run in this sandbox on 2026-07-28.

| Premise (as this doc states it) | Status | Where it lives / would live | Validating command or reference | Observed result |
|---|---|---|---|---|
| State has an immutable append-only log and bit-for-bit replay | BUILT | `host/store`, `host/replay` | Controller verification at origin/dev `7a3e7c6`; `scripts/verify_go.sh`; landed iter-16 PRs #2, #3, #4, #6, #7, #8 | SQLite compare-and-append is single-transaction/single-writer; replay doubling establishes A == B == recorded, and CI runs rather than skips against sha256-pinned AILANG v0.30.0. |
| Every fact and log entry has usable provenance | PARTIAL | `host/store/store.go:575` (`Store.Commit`); queue item 4b `w-store-durability` | Controller first-party iter-25 zero-ref commit probe | None of eight reference fields is validated. Seven zero refs create an unreadable row; an empty `NextWorld.Ref` becomes selected head, and later legal commits can extend the poisoned chain. The surface MUST expose unreadable links. |
| Every `Evidence` value has one of the four display grades | NOT BUILT | `world/types.ail:23-28`; mapping to be ratified at §7.2 and implemented with approval-inbox work | `sed -n '20,58p' world/types.ail`; controller search of `*.ail`, `*.go`, `*.sql` for `PROVEN|ATTESTED|CLAIMED` | Five variants exist and no mapping exists. `TestReport → TESTED`, `RecordedEffect → ATTESTED`, `AiReview → CLAIMED`; `CompilerOutput` and `HumanApproval` are unmapped, and PROVEN has no carrier. |
| `Proposal.confidence` is grounded evidence suitable for display | NOT BUILT | `world/types.ail:42-51`; ratification point §7.5 | `sed -n '42,52p' world/types.ail` | It is a bare `float` with no evidence ref; rendering it as epistemic evidence would be confidence theater. `AiReview`, by contrast, carries `HashRef` with its float. |
| `Human.Approve` is enforced as a budgeted effect | NOT BUILT | DESIGN.md §8 intent; future effect broker; queue item 4 `w-effect-broker-m3` (PARKED) | `host/daemon/daemon.go:34`; DESIGN.md:296 | The daemon states there is no effect broker or capability/budget authority. This is design intent, not landed enforcement. |
| Worlds fork near-O(1), enabling branch and compare-arms gestures | NOT BUILT | Future M5 speculation infrastructure; `host/store/store.go` | `rg` of the public `Store` methods; DESIGN.md §13.5 item 2 | Public store API has one `SelectedHead`/`SelectHead` and no fork, branch, or compare API. Content-addressed storage makes sharing plausible only; DESIGN.md correctly states this as an M5 requirement. |
| worldd serves the localhost workbench | NOT BUILT | Future worldd renderer route; `host/daemon/daemon.go:355-362` today | `sed -n '355,365p' host/daemon/daemon.go` | Exactly eight JSON REST routes are registered: health, head, world, object, log entry/range, registry, and commit. No HTML/web/SSE renderer route exists. |
| Open projections can carry the proposed P1/P2 UI protocol | BLOCKED | Queue item 5 `w-mcp-projection`; open question 9; upstream `sunholo-data/ailang#498` | Controller stdio MCP probe on pinned v0.30.0, iter-24 | `submit_feedback` survives unfiltered, `--routes-only`, and `--caps ''`; it routes to public-feedback with Pub/Sub notification. Flags are process-wide, `--caps` gates execution rather than discovery, and per-session filtering is unavailable. |
| MCP clients can safely render decision packets and grounded links | BLOCKED | Same projection item plus a ratified packet/link schema and renderer adapter | Controller iter-24 probe; DESIGN.md open question 9 | Existing MCP projection cannot expose an exact session-authorized allowlist. No verified packet/grounded-link renderer exists in this repo; P1/P2-over-MCP remains a requirement. |
| Collaboration Hub is reusable implementation machinery for World | UNVERIFIED | Out-of-repo `sunholo-data/ailang`; local design references only | `rg -n 'Collaboration Hub|React approval queue|message center|Hub' --glob '!design_docs/**' .`; DESIGN.md:141, 373-375, 728-730 | No Hub source exists here. DESIGN.md verifies its role and coordinator-schema coupling, not its internal package boundaries. Treat it as a pattern donor until its source/build graph can be inspected. |
| Approval machinery already demonstrates propose → human-approve → execute | PARTIAL | Out-of-repo Coordinator + approval queue + Hub; future typed World approval inbox | DESIGN.md:141 and 373-375 | Existing flow operates over coordinator DB rows; World's proposal and `HumanApproval` types exist locally, but broker enforcement, total grading, and the World inbox are not built. |
| Notification machinery satisfies P5 and clause 3 | BLOCKED | Future brokered local notifier; Hub notification internals out of repo; MCP built-in upstream | `host/daemon/daemon.go:34`; controller MCP probe; Hub-source search above | No local notification broker exists. The only measured MCP notification path is ambient Pub/Sub egress from `submit_feedback`, which conflicts with clauses 2 and 3. Hub internals remain unverified. |
| Provenance machinery can support P2/P4/P6 universally | PARTIAL | `world/types.ail`, content-addressed store/log/replay; durability item 4b | Evidence/schema inspection plus controller zero-ref probe | Typed hash references and replay exist, but commit admits broken references and the grade mapping is incomplete. Renderers must preserve `UNSUPPORTED`/unreadable states. |
| UI protocol machinery supports generated World projections | BLOCKED | AILANG `serve-api` outside this repo; future World transition registry and adapters | DESIGN.md §3.7/open question 9; controller iter-24 MCP probe | Static export projection exists per controller verification, but exact per-session World projection does not; the unsafe built-in tool blocks the proposed boundary. |
| A local reference renderer exists | NOT BUILT | Future workbench in this repo | `find`/`rg` inspection of this repo; daemon route inspection | Only two self-contained design fixtures exist. They demonstrate grammar, not a live renderer. |
| Linked approval-inbox mockup exists as attributed | BUILT | `design_docs/mockups/approval-inbox.html` | `find design_docs/mockups ... -exec wc -c`; controller citation check | File exists and is 4270 bytes; it renders the P1/P5 decision-packet fixture described in §6.5. |
| Linked grounded-prose mockup exists as attributed | BUILT | `design_docs/mockups/grounded-prose.html` | Same local `find`/`wc`; controller citation check | File exists and is 4136 bytes; it is the P2/P3/P6 fixture described in §6.5. |
| Decision packets carry a ledger-recorded creation time, deadline and typed timeout policy (§3.1) | NOT BUILT | Future packet schema (kernel-adjacent, ratification point 3); `world/types.ail` has no packet type | Controller inspection of `world/types.ail` exported types: `Capability`, `CommitResult`, `Evidence`, `Proposal`, `RecordedTransition`, `Transition`, `Verification`, `World` | No decision-packet type exists, so no creation time, deadline or timeout policy field exists. §3.1 is a requirement on the schema freeze, not a description. |
| A `Human.Approve` deadline expiry is a deterministic, replayable Timeout transition (§3.1) | NOT BUILT | Future effect broker (queue item 4 `w-effect-broker-m3`, PARKED) + replay | `host/daemon/daemon.go:34`; controller search of `host/replay/*.go` for `timeout` | No effect broker exists, so no approval Timeout transition is emitted or replayed. The only timeout in `host/replay` is `execTimeout = 60 * time.Second` (`replay.go:47-49,321`), a bound on each archived-interpreter subprocess — unrelated to approval deadlines, and recorded here so the two are never conflated. Deterministic-timeout tests are an acceptance criterion for the broker item. |
| DESIGN.md and SCENARIOS.md citations say what this document attributes to them | BUILT | DESIGN.md §§11.2, 13.5 and open questions 7/9; SCENARIOS.md scenarios 1, 3, 4, 5 | `sed`/`rg` local inspection; controller citation check | Five surfaces, attention metric, Hub question, protocol question, approval inbox, provenance walk, budgeted human approval, and speculative branches are present as cited. |

### 8.1 Existing machinery — reuse or replace

| Category | What exists | Reuse or replace before ratification |
|---|---|---|
| Approval | Out-of-repo Coordinator/queue/Hub flow over DB rows; local `HumanApproval(HashRef)` type. World inbox and broker enforcement are not built. | **Reuse interaction patterns; replace the data/authority seam.** World approval must be a typed transition with a total evidence grade and brokered budget, not a coordinator-row mutation. |
| Notification | Hub message-center pattern is named but its implementation is UNVERIFIED here. No local broker exists; pinned MCP exposes an unauthorized Pub/Sub feedback path. | **Reuse batching/message-center patterns only. Replace ambient delivery with a local, capability- and budget-checked notifier.** Do not reuse `submit_feedback`. |
| Provenance | Content-addressed types, immutable log, and replay are built; reference validation and total evidence grading are incomplete. | **Reuse store/log/replay after hardening.** Add ref-integrity enforcement through `w-store-durability`; do not build a parallel provenance database. |
| Projection | Out-of-repo AILANG static export projection exists; dynamic, session-filtered World transition projection does not. | **Reuse upstream codecs/standards only after the exact-provider seam is safe.** Keep World schemas and authority decisions local; do not reverse-proxy and re-encode an unsafe tool list. |
| MCP/UI-protocol | MCP/A2A machinery exists upstream; dialect question 9 and packet/link schemas remain open, and v0.30.0 has a verified ambient-tool blocker. | **Reuse protocol machinery conditionally; replace neither standard with a proprietary wire format.** Block delivery until `#498`, per-session filtering, and conformance tests establish the boundary. |
| Renderer | Two local HTML fixtures; out-of-repo Hub patterns; no live renderer or Hub source here. | **Build fresh code in this repo, reuse fixtures and Hub interaction patterns.** Revisit Hub code reuse only after source inspection proves a dependency-free protocol-renderer boundary compatible with typed ledger objects and clause 2. |

## 9. Quorum verification log

Pick-time quorum, iteration 26 (2026-07-28). Reviewers `gpt5-6-sol` + `gemini-3-1-pro`, both
present in both rounds, cap $0.25/reviewer. Artifacts:
`.ailang/state/mission-quorum/human-surface-2026-07-28T03-54-41Z.json` (r1) and
`…T04-09-06Z.json` (r2). Total metered $0.089226.

### Round 1 — BLOCKED (2/2 reject, controller PASS)

Neither objection disputed the design direction; both were the two sections this document never
had, because v0.1 was authored **attended** rather than through `design-doc-creator` and so
bypassed that skill's hard gates.

- **`gpt5-6-sol`** — "The foundational premise is asserted rather than verified … P2–P7
  consequently depend on unverified backend behavior." Fix: add a premise-verification section,
  one row per dependency, and *"Rewrite every unsupported current-tense claim as a requirement."*
- **`gemini-3-1-pro`** — "The document fails the Conflict Surface gate. It mandates building a
  'FRESH' localhost reference renderer and explicitly dismisses reusing the existing 'Hub' UI …
  but completely omits the required conflict-surface analysis."

Applied in **v0.2** by the rotation designer (`codex:gpt-5.6-sol`): §8 premise table (16 rows),
§8.1 reuse-or-replace across all six named categories, §6.1 Hub conflict surface, §1 recast from
assertion to requirement, dependency notes on P3/P5/P7 and §6, the unavailable-grade anti-pattern,
and ratification points 7.2/7.5 restated as decidable questions.

**The designer declined one thing on purpose**, and it is worth recording: `gemini-3-1-pro`'s fix
supplied an *example* rationale ("the Hub is inextricably tied to cloud authentication and
multi-tenant routing"). That is a hypothesis about a codebase in another repository. §6.1 marks
Hub internals **UNVERIFIED** and names the inspection that would close the gap, rather than
adopting a convenient reviewer-supplied reason as fact.

### Round 2 — BLOCKED (2/2 reject, controller PASS) → narrow-refinement carve-out applied

Both round-1 objections were accepted as satisfied. Both reviewers then raised, **independently
and across two providers, the same new objection**: the interaction grammar admits an unbounded
wait on a human. `§3`'s *defer* ("park for more evidence") and P5's *batch over interrupt* had no
TTL, no expiry transition, and no deterministic outcome when the human never answers — the
mission's own Standing Rule 6 ("every wait is bounded"), missing at the UX layer.

Both objections carry concrete reviewer-authored fixes and neither disputes the design direction
(completeness + determinism only), so the charter's **narrow-refinement carve-out** applies: the
controller made a bounded second revision applying the reviewers' **verbatim** text — no
controller-invented resolution, no third round.

| Applied | Source | Where |
|---|---|---|
| "To enforce bounded waits, every Decision Packet and `Human.Approve` effect MUST carry a strict TTL. If the human does not decide before expiration, the packet deterministically resolves to a typed rejected-timeout (or a predefined safe fallback), unblocking the agent." | `gemini-3-1-pro`, verbatim | §3, appended as its fix directed |
| The full binding principle: ledger-recorded creation time + deadline + typed timeout policy; DEFER MUST rebound and MUST NOT park indefinitely; an explicit Timeout transition at deadline; "Silence MUST never synthesize approval or rejection"; replay reproduces deadlines from ledger time, not wall-clock races. | `gpt5-6-sol`, verbatim | new §3.1 |
| "Add verification-log rows for the packet lifecycle schema and deterministic timeout tests." | `gpt5-6-sol`, verbatim | two new §8 rows, both **NOT BUILT** |

Controller-added routing only (not a resolution): ratification point 1 now names the **typed
finite set of timeout policies** as something the human must fix, since the applied principle
requires such a set and deliberately does not enumerate it.

**Generator≠judge — FLAGGED.** The v0.2 designer and the `gpt5-6-sol` reviewer seat are the same
model family, so that seat was a self-review. Retained rather than excluded, per this mission's
precedent (iter-24): reject-by-default synthesis means a self-*pass* cannot manufacture a PROCEED,
so retaining the seat can only *add* objections. It did not rubber-stamp itself — it rejected its
own revision in round 2 and produced the round's strongest objection. Independent rejector
throughout: `gemini-3-1-pro`, a different provider, which reached the same conclusion on its own.

**Controller's own first-party measurements** (not inherited from any sub-agent): the five-variant
`Evidence` ADT vs the four-grade gradient and the absence of any grade string in `*.ail`/`*.go`/
`*.sql`; the eight exported types in `world/types.ail` (no packet type); the eight daemon routes;
the thirteen public `Store` methods with one selected head; `verify_ail.sh` **rc=0**, 4/4 required
identities across 10 modules, 14/14 named tests, unchanged by this doc-only change.

---
*v0.1 drafted attended 2026-07-27. v0.2 applies quorum round-1 blocking fixes 2026-07-28. v0.3
applies the round-2 carve-out fixes verbatim. Quorum is COMPLETE; what remains is human
ratification (§7). Mark ratifies before `w-approval-inbox` or any M6 work routes to sprint.*

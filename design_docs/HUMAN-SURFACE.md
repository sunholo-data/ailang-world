# The Human Surface — founding UX design (v0.1)

**Status**: **v0.1 principles PROVISIONALLY RATIFIED** (Mark, attended, 2026-07-27: "good
principles for now, yes good start") — binding working basis for all human-facing work; formal
pick-time quorum + full §7 ratification (trust-grade taxonomy naming, packet-schema freeze
timing) still runs when queue item 6b is picked. Authored attended (Mark + World coordinator).
This is the experience-layer peer of [DESIGN.md](DESIGN.md): §11 answered *what surfaces exist*;
this answers *how humans and the AI state machine actually meet*. **Binding input** for
`w-approval-inbox` (clause 5) and M6 (generated projections) — neither routes to sprint before
this doc is quorum-reviewed and Mark-ratified. Companion: [SCENARIOS.md](SCENARIOS.md).

---

## 1. The premise — what no UI could assume before

Every prior interface renders one of two things: **mutable opaque state** (classic GUIs — the
screen shows a value; where it came from and whether it's true are unknowable) or **token
streams** (chat UIs — prose whose relationship to reality is vibes). World's surface is the
first that can assume **everything behind the glass is a typed, verified, replayable ledger**:
every fact has a provenance chain, every claim has an evidence grade, every state has a
navigable history, every action is a transition.

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
When the AI writes to a human, **every load-bearing noun is a link to a typed world object** —
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

### P4 — Time is navigation
State is an immutable log, so "where am I?" = **(domain × time)**, not a page. The timeline
scrubber and the **world-diff view** ("show me now vs before that transition") are chrome, not
features — as fundamental as the back button. Scenario 3's provenance walk and clause 5's
5-minute "why" answer are P4 exercised backwards; the R5 replay is P4 exercised forward.

### P5 — Attention choreography
`Human.Approve` is a budgeted effect (§8), so the surface treats interruption as spending:
batch over interrupt; rank packets by **irreversibility × novelty**; never ask what
verification settled; deliver ONE digest, not N pings; "review the night" after auto-commits
rather than gating them. Quiet is the default state of a healthy World — the UI must make
quiet feel like health, not absence.

### P6 — Five zoom levels, everywhere
Every object supports the same ladder: **glance** (chip/count) → **decide** (packet) →
**inspect** (evidence bundle) → **walk** (provenance graph) → **replay** (re-run the episode).
One gesture-grammar for goals, deployments, documents, agents, worlds. Learn it once.

### P7 — Speculation as a gesture
Because worlds fork O(delta), **"what if" is a primitive**: fork → let agents run against the
branch → compare evidence side-by-side → commit or discard. The compare-arms view (scenario 5)
is a standing screen, not an expert mode. Immutability makes counterfactuals safe enough to
hand to a human as a button.

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
  rendering (an MCP client like Claude Code renders packets and grounded prose when you're away
  from the workbench — scenario 4). Casual and mobile use may be chat-first; the workbench
  remains the place decisions are *designed* to happen.
- **Terminal/CLI** = plumbing-truth access and scripting for power users; never the designed
  experience, always available (no capability the workbench has may be CLI-impossible, and
  vice versa — one typed layer beneath both).

## 3. The input side — signing the type, not the vibe

The goal composer accepts natural language and echoes back BOTH renderings: the prose
paraphrase AND the **typed Goal object** (budget, capability grant, contracts, approval gates)
— and the human confirms the *typed* one. What you sign is what the machine enforces; the
prose was only ever the request. Decision verbs beyond approve/reject are first-class:
**attenuate** (approve narrower), **defer** (park for more evidence), **reject-with-reason**
(typed feedback to the agent), **amend-policy** (this class never asks again).

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

## 6. Build strategy (unchanged from §11/§3.7, now sequenced)

Projections emit an open agent-UI protocol (A2UI/AG-UI — dialect = open question 9, decided at
M6). ONE hand-built reference renderer: the **workbench**, a localhost web surface served by
worldd — browser-opened, cross-platform, zero install beyond the binary. Chat clients (MCP) get
P1/P2 degraded-but-consistent: packets as structured tool results, grounded prose as resource
links. Phone = same packets, smaller (scenario 4).

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

1. The seven principles + anti-patterns as BINDING on all human-facing work (evaluator-scored,
   like coding-standards.md).
2. The trust-grade taxonomy (PROVEN/TESTED/ATTESTED/CLAIMED) — names and cut-lines.
3. Decision-packet schema freeze timing (it becomes a world type — kernel-adjacent).
4. Open questions 7 (Hub vs fresh: recommend FRESH, Hub as pattern donor) and 9 (dialect —
   defer to M6 with a common-core emitter).

---
*v0.1 drafted attended 2026-07-27. Quorum review at pick; Mark ratifies before
`w-approval-inbox` or any M6 work routes to sprint.*

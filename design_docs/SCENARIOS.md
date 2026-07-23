# AILANG World — Human Interaction Scenarios

*Companion to [DESIGN.md](DESIGN.md) §11 (Human and AI interfaces). Written attended
(Mark + Fable, 2026-07-23). These are day-in-the-life walkthroughs of the HUMAN side of
World — what it feels like to operate, not how it's built. Personas: primarily the
operator (Mark); scenario 6 shows a non-engineer.*

---

## The one-paragraph answer

You never operate World imperatively; you govern it. Your three verbs are **express a
goal** (with budget and authority attached), **decide what verification couldn't settle**
(the approval inbox), and **ask the world why** (provenance, time-travel). Everything else
— planning, execution, verification, evidence-gathering — happens between agents and the
kernel, and reaches you only as distilled, decision-ready items. And crucially: you keep
talking to a conversational AI the whole time. The chat window is itself a projection of
the world — World doesn't replace the conversation, it gives the conversation something
solid to act through.

## Scenario 1 — Morning: the approval inbox

You open World (desktop surface, or just ask your chat agent "what's waiting on me?").
Overnight: 17 transitions. 14 auto-committed under standing policy (minor dep bumps with
green evidence — you review the policy, not the PRs). 2 were auto-rejected by the verifier
before ever reaching you (a contract counterexample; a budget overrun) — you only see
those if you ask. **3 need you**, because they hit gates verification cannot settle:

1. *Deploy docparse 0.9.2 to staging* — traced to YOUR goal from Tuesday. The card shows:
   evidence bundle (ai-check green with 4 Z3-verified contracts, 214 tests pass, the
   "billing untouched" contract held — enforced, not promised), the exact authority
   requested (`Cloud.Deploy` scoped to docparse/staging), budget consumed ($3.20 of your
   $10), and reversibility (roll-forward to 0.9.1). You approve in one glance because the
   machine already did the reading.
2. *Widen a capability* — the sonnet evaluator's calibration record (0.94 over 60
   proposals) qualifies it for broader read scope. Capability changes always ask.
3. *Send release notes externally* — Email/External class always asks, regardless of
   evidence.

Time spent: ~4 minutes. The design goal (§11.2): maximize sound decisions per second of
your attention. Verification filters; you judge only the remainder.

## Scenario 2 — Expressing a goal

You say — in chat, CLI, or the composer: *"Add YAML export to docparse. Don't touch the
billing module. Budget $10. Staging deploys need my approval; never prod."*

World turns that sentence into a typed Goal object: a budget ledger, a capability grant
(FS scoped to the repo, `Git.Commit` on a branch, `Cloud.Deploy` staging-only and gated
behind `Human.Approve`), and a **contract** — "billing module untouched" becomes a
machine-checked invariant on every transition under this goal, not a hope. You watch
progress as world transitions (proposed → verified → committed), not log scroll. The
difference from today: your sentence became an enforceable envelope. Nothing inside it
needs you; nothing outside it can happen.

## Scenario 3 — "Why is staging behaving oddly?" (the provenance walk)

Two weeks later something's off in staging. Today this costs a grepping session and an
essay you have to trust. In World: select the deployment object → **provenance** → the
chain walks itself: deployed-by transition → proposal → agent → the evidence *as it stood
that day* → the originating goal (your YAML request) — and the surprise: a config change
arrived via cascade from a package bump, auto-committed under standing policy 7.
Time-travel: diff the world now against the world before that transition. Your decision is
governance, not archaeology: tighten policy 7 to exclude config-touching cascades.
Ninety seconds, and the answer came from the state itself, not from an agent's
reconstruction.

## Scenario 4 — Away from the desk (approval as a budgeted effect)

You're out. An agent hits an approval gate. World does NOT ping you per-event: the
scheduler holds it, batches it with two other pending decisions, and sends ONE
notification through your channel (phone/Discord). On your phone, the same typed proposal
renders as a small generated card (A2UI) — same evidence, same authority ask, smaller
screen. Approve from the beach. If you don't respond within the goal's wall-time budget,
the proposal parks gracefully and the agent moves to other work — no unbounded wait, no
wedged loop; the 6-hour-outage lesson is kernel behavior now, not script discipline.

## Scenario 5 — "What if?" (speculative branches for humans)

You wonder: *"what would upgrading the whole fleet to the new model actually do?"* You
fork the world. Agents run the upgrade against the fork — benchmarks, costs, breakage —
while the real world runs untouched. You compare evidence between branches
side-by-side, then commit the winner or discard the experiment. Immutability makes
counterfactuals free: best-of-N becomes a decision tool for humans, not just a scheduler
trick for agents.

## Scenario 6 — A non-engineer at the boundary (post-v1, M8 flavor)

A colleague on another laptop asks the document world: *"assemble the Q3 client report
from these sources."* Same machinery, different domain: agents propose document
transitions, evidence is citation-coverage and template conformance, and their approval
inbox shows a rendered preview diff. Capability boundaries mean the document agents
cannot see code repositories, and your engineering agents cannot read client documents.
They never learn AILANG; they express goals and make decisions. Citizenship tiers apply
to humans too.

## How the interface is actually handled — the mechanics

1. **One typed layer, many renderings** (§11). CLI, chat, desktop surfaces, phone cards
   are all projections of the same transitions. A human action is itself a transition
   with authority and a trace — there is no admin backdoor to bypass the record.
2. **The conversation is a projection too.** You'll still talk to a conversational agent
   (Claude, or any MCP-speaking assistant) as your primary interface. The difference is
   what's underneath: the agent translates your intent into typed goals, and renders
   world state back as prose — but the goals, approvals, and evidence are solid objects
   in the world, not turns in a chat log that compaction will eat.
3. **Attention is budgeted and defended.** `Human.Approve` is an effect with a budget
   (§8). Verification always runs before you're asked. Batching over interrupting.
   Evidence arrives distilled with drill-down, never as raw logs.
4. **The trust dial moves with evidence.** Standing policies + calibration records
   (§13.3) migrate decision classes from ask-me to auto over time — you set the
   constitution; the dial is moved by track record, and every auto-commit stays
   reviewable after the fact (scenario 1's "review the night").
5. **Approve/reject is never the only choice.** Reject carries a reason back to the
   agent as typed feedback; you can also attenuate (approve with narrower scope), defer
   (park until more evidence), or edit the policy so the class never asks again.

*The recurring pattern across all six: World's offer to the human is a change of job —
from operator and archaeologist to goal-setter, judge, and constitution-author.*

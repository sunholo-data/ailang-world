---
title: Use cases
sidebar_position: 1.5
description: What AILANG World is for, who it is for, and what works today versus what is on the roadmap.
---

# Use cases

AILANG World is a substrate for work done by AI agents. Agents **propose**, a deterministic
verifier **checks**, and only authorized, budgeted changes **commit**, each one leaving a
receipt you can walk back later. A human's job shrinks to three verbs: **set a goal with a
budget**, **decide what verification couldn't settle**, and **ask the world why**.

World does not compete with coding agents. Agents are its *residents*. It competes with the
scripts, conventions, log archaeology and human vigilance that hold agent work together today.

Each use case below is marked:

- <span className="world-status world-status--met">Works today</span> runs on the current release, with evidence linked.
- <span className="world-status world-status--unmet">Roadmap</span> is designed or partly built, but not usable end to end yet.

---

## 1. Governed coding agents <span className="world-status world-status--met">Works today</span>

**The problem.** A coding agent with a shell has ambient authority. It can read your SSH keys,
rewrite `.git/config`, plant a CI workflow, or drop hooks into an editor config that run
the next time a human opens the repo. You manage that risk with prompts and hope.

**What World does.** The agent gets a **worktree, not a machine**. It works through eight AILANG
tools over MCP: read, write, edit, check, run, builtins-search, examples-search and the AILANG
CLI. Every call is a transition:

- The **plan is pure**: arguments are checked, and unknown keys and path tricks are refused.
- The **one effect goes through the broker**: capability, scope and budget are checked.
- The **effect is confined** by AILANG's policy layer: no shell, no network, `.git` read-only,
  and case-folded deny lists, so `.CLAUDE/settings.json` is refused just like `.claude/…`.
- The call **commits one log entry**, with the plan and the effect record attached.

**Any MCP client works.** In the first live session, pi and Claude Code each read, edited,
type-checked and ran a program through World's tools only, with built-in tools switched off.
Both took about 9 seconds, and every call was recorded.

**Who it's for.** Teams running several agents (Claude Code, codex, pi or local models) who need
one authority model for all of them instead of one prompt per agent.

→ [Coding tools over MCP](getting-started/coding-tools.md) · [Tool confinement](security/tool-confinement.md) · [For AI agents](agents/index.md)

---

## 2. "Why did this happen?" — provenance instead of archaeology <span className="world-status world-status--met">Works today</span> · <span className="world-status world-status--unmet">Roadmap</span> (timed proof)

**The problem.** When something goes wrong, someone reconstructs reality from logs, `git log`
and memory. That person is often an agent that has to re-derive context it has already lost.
These archaeology sessions cost hours, and their conclusions are sometimes wrong.

**What World does.** Every committed change carries its plan, its effect records (what was
requested, what was allowed or denied, what came back) and the commit that landed it. You can
**walk the chain**: from a file content back to the tool call that wrote it, the session that
held the grant, and the goal it served. Effect results are recorded at the broker, so "what was
actually sent" is a fact you can query, not a conclusion you have to reason your way to.

**Status.** The records and the walk exist today. What the 1.0 release still has to prove (clause 5)
is that the walk answers ≥3 **real** "why" questions in ≤5 minutes each, timed against the old
grep-and-logs method run by an agent confined to read-only access.

→ [Provenance and replay](concepts/provenance-and-replay.md) · [Provenance walks](guides/provenance-walks.md) · [Roadmap](roadmap/index.md)

---

## 3. Deterministic replay and audit <span className="world-status world-status--met">Works today</span>

**The problem.** "Re-run what happened" usually means re-running it and hoping the world is
the same, or reconstructing it from logs.

**What World does.** Transitions are pure AILANG programs, and effect results are recorded
inputs, so replay **reconstructs any past state bit for bit**. The log pins the exact interpreter
binary (by content hash), the transition source and its inputs. Replay re-runs the pure part on
the archived interpreter and feeds back the recorded effect results. It never calls a live effect
handler. History made under different interpreter versions still replays, because each entry
carries its own pin.

**Who it's for.** Anyone who has to show what happened rather than assert it: post-incident
review, compliance, and debugging agent behaviour after the fact.

→ [World and transitions](concepts/world-and-transitions.md)

---

## 4. Fleets with different levels of trust <span className="world-status world-status--met">Works today</span> (per-session grants) · <span className="world-status world-status--unmet">Roadmap</span> (wider effect set)

**The problem.** Today, trust between agents is env pinning plus discipline: which model may
touch which repo, which one may spend money, which one may ask a human.

**What World does.** Authority is **typed and per session**. A local model can hold
`Workspace.Read` on one worktree and nothing else; a cloud executor can hold `Workspace.Write`
plus `Ailang.Run` with a budget of 50 calls. Grants carry scope, expiry and budget. Spend persists
across calls and restarts. Denials are recorded, and a session only sees the tools it holds grants
for. Six grants cover today's eight tools. Git, GitHub, deploy and model-call effects follow the
same pattern; they are designed but not yet served over MCP.

→ [Effects, capabilities and budgets](concepts/effects-capabilities-budgets.md) · [Sessions and episodes](concepts/sessions-and-episodes.md)

---

## 5. Honest agent evaluation <span className="world-status world-status--unmet">Roadmap</span> (harness in progress)

**The problem.** Agent benchmarks rarely control what the agent could actually touch. A "tools-only"
arm that quietly keeps a shell, or a grader that runs a different binary than the agent tested
against, makes the numbers mean less than they claim.

**What World does.** The 1.0 non-inferiority floor (clause 4) runs Claude Code and codex on
AILANG's core benchmarks in two arms: native shell, and World's tools only. The World arm is
**enforced, not trusted**. Every graded solution must match the content of the World log entry
that wrote it. The grader port agrees with **1,351 of 1,357** banked results, eligibility is sealed
before any World statistic is computed, and the final run is pre-registered.

**Who it's for.** Teams who want to compare agents, models or toolsets with confinement as a
measured condition rather than an assumption.

→ [Roadmap](roadmap/index.md)

---

## 6. An AI employee with an accountable owner <span className="world-status world-status--unmet">Roadmap</span>

**The problem.** A long-running agent with its own email, GitHub and chat accounts can do the
delegable work of a remote role: inbox, tickets, code, reviews, the standup note. Making that
**governable** means a capability envelope it can't widen, an owner whose attention is a budget
rather than an interrupt stream, and a receipt for every action.

**What World does.**
- **Approvals are budgeted effects.** `Human.Approve` is a budgeted effect, so "how many times may
  this agent ask me today" is scheduled, not hoped for.
- **Authority is typed.** "May open a PR on repo X but never push to `main`" is a type the broker
  enforces.
- **The owner's job is small.** It is the three verbs above.

Sunholo already runs a virtual employee, Daneel, with its own identity and capability manifest.
Moving that kind of resident onto World is the intended next application after the mission loop.

---

## 7. Self-modification under governance <span className="world-status world-status--met">Works today</span>

**The problem.** A system that can change its own code needs a way to do it that can't
quietly rewrite its own rules.

**What World does.** World's own kernel package, `world/core`, ships through World's propose →
verify → commit pipeline. Its contracts are Z3-verified, it is published to the AILANG package
registry, and the irreversible publish is **fenced to a human at a real terminal**: no agent,
CI job or script can satisfy that fence. The mission-loop machinery, the compiler and the live
daemon are excluded from self-modification by design. `world/core@0.1.1` was published this way
on 2026-10-02.

→ [Packages and self-modification](concepts/packages-and-self-modification.md) · [Attended steps](guides/attended-steps.md)

---

## 8. Training data from verified work <span className="world-status world-status--unmet">Roadmap</span>

**The problem.** When an agent session ends, the diff survives but the *process* evaporates:
the evidence weighed, the dead ends, the reasons one option beat another.

**What World does.** Every episode is a typed, replayable trace of proposals, verification
outcomes, effect records and human decisions. That is process-supervision data, and it is
**already verified and attributed** because the substrate produced it, not a scraper. It
generalizes AILANG's founding idea, "structured execution traces for AI training", from single
programs to whole episodes of work.

---

## 9. Beyond code <span className="world-status world-status--unmet">Roadmap (post-1.0)</span>

The same kernel can manage worlds other than repositories:

- **Document and data pipelines with provenance.** A DocumentWorld or DatasetWorld where each
  derived artifact knows exactly which inputs, versions and decisions produced it.
- **High-stakes operations.** Deploys, data handling and compliance work where "who authorized
  what, on what evidence" is the product rather than a reconstructed log.
- **A local-first home substrate.** Household agents operating under capability boundaries on
  machines you own.

---

## What World is *not* (today)

- **Not a coding agent.** It hosts agents; it does not replace them.
- **Not a workflow engine.** Temporal or LangGraph give you durable state and retries. World adds
  verification *before* execution, constructive replay, and authority as types. Those are
  properties of the language joined to the substrate, and they are hard to bolt on afterwards.
- **Not 1.0 yet.** Five of seven release clauses are met. The two that remain measure World
  against the status quo: agents must do no worse inside World than with a shell, and provenance
  must beat log archaeology on real questions. See the [roadmap](roadmap/index.md).
- **Not cloud-first.** One machine, SQLite, no cloud dependencies in the core.

---
title: Concepts overview
sidebar_position: 1
description: The ideas behind AILANG World — a semantic database whose transaction language is AILANG — and how they map onto the daemon you can run today.
---

# Concepts overview

**AILANG World is a semantic database whose transaction language is
[AILANG](https://github.com/sunholo-data/ailang).** The kernel is an immutable,
content-addressed world graph. State changes only through **transitions**: pure, type-checked
AILANG functions. Agents never modify the world directly. They call a transition, the host
checks authority and budget, runs any effect through a broker, and commits one log entry that
records what happened.

> Unix made everything a file. AILANG World makes everything a typed state transition.

This section explains the model. Each page states what is built today and what is still design
intent, because World is pre-1.0 and the two are not the same thing.

## The pieces you can run today

| Piece | What it is | Where to read more |
|---|---|---|
| `ailang-worldd` | One local daemon over one SQLite file: REST, MCP and A2A on loopback only | [Operating the daemon](../guides/operating-the-daemon.md) |
| World store | Content-addressed objects, an append-only transition log, world rows, registry heads, an intent/outcome journal | [World and transitions](world-and-transitions.md) |
| Effect broker | The in-process authority boundary: capability, scope, budget, recorded outcome | [Effects, capabilities and budgets](effects-capabilities-budgets.md) |
| Sessions | A bearer credential bound to one episode and a set of grants, minted at a terminal | [Sessions and episodes](sessions-and-episodes.md) |
| Transition registry | Published, pinned transitions, served per session over MCP and A2A | [Transitions manifest](../reference/transitions-manifest.md) |
| `se-tools` | Eight AILANG coding tools (read, write, edit, check, run, two searches, CLI) | [Coding tools](../getting-started/coding-tools.md) |
| `world-publish` | The attended operator entrypoint for registry writes and the `world/core` package publish | [Attended steps](../guides/attended-steps.md) |

## The model in one paragraph per idea

**World.** An immutable state value. Each successful commit produces a new world row
(`revision`, `stateRoot`, `logHead`) and moves the selected head. Old worlds stay readable.

**Transition.** The only way to change state. A transition is published as AILANG source,
pinned by content hash together with the interpreter that runs it and a semantics epoch.

**Propose, verify, commit.** A caller asks for a transition by ID. The host checks the
session's grants against the transition's declared access, runs the pure transition in a
capsule, executes any planned effect through the broker, and commits the result with a
compare-and-append against the head it observed. A stale caller gets a structured conflict, not
a corrupted world.

**Effects.** The only nondeterminism. Every brokered effect needs a live grant with a matching
effect name and scope and enough budget. Allowed, denied and failed effects all leave an
immutable effect record.

**Provenance.** Every log entry points at a record object, which points at its input, output,
plan and effect records. "Why did this happen" is answered by walking those references
instead of searching logs.

**Humans govern; agents operate.** The steps that grant authority — minting a session,
publishing transitions, publishing a package — are fenced to a human at a real terminal.

## What World is not

- It is not a replacement for Linux or macOS. It runs on top of them.
- It is not a container runtime. The tool handlers use AILANG's own policy layer for
  confinement; there is no Docker and no network isolation in the core
  (see [Security model](../security/model.md)).
- It is not multi-node. One daemon, one SQLite file, one writer process.

## Status

World is pre-1.0 and built by an autonomous mission loop with human ratification gates. The
1.0 bar has seven checkable clauses (deterministic kernel, local-first daemon, explicit
authority end-to-end, a resident-agent non-inferiority floor, provenance teeth, a
protocol-native boundary, and one controlled self-modification). Several are met, several are
in progress. Known gaps are listed in [Residuals](../security/residuals.md).

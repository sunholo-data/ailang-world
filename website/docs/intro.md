---
title: What is AILANG World
sidebar_position: 1
slug: /intro
description: AILANG World is a semantic database whose transaction language is AILANG.
---

# What is AILANG World

> **AILANG World is a semantic database whose transaction language is [AILANG](https://ailang.sunholo.com).**

Unix made everything a file. AILANG World makes everything a **typed state transition**.

World is not a replacement for Linux or macOS. It is a semantic operating environment that
runs on top of an existing operating system. A traditional OS manages processes, files and
devices. World manages **goals, typed state, capabilities, effects, evidence, budgets,
contracts, provenance and AI proposals**.

## The defining move

The kernel is an immutable, content-addressed world graph. State changes only through
**transitions**: pure, type-checked AILANG functions that are checked *before* they run.

- A **transition** is a pure function from a world and a command to a new world. It is
  deterministic, type-checked, and its contracts can be verified with Z3 before it runs.
- **Effects** are the only source of nondeterminism. They are declared in the function
  signature, enforced by the compiler, and carried out by an effect broker that checks a
  capability and a budget for every one.
- **History** is an append-only log of committed transitions. A pure transition plus its
  recorded effect results reconstructs any historical state bit-for-bit, because the
  language guarantees the pure part.

## Propose, verify, commit

AI agents never change the world directly:

1. An agent **proposes** a change as an AILANG program.
2. A deterministic verifier (types plus Z3 contracts) **checks** it.
3. Only verified, authorized, budgeted proposals **commit**, and each one leaves evidence
   and a replayable trace.

Agents are World's operators and humans are its authorities. A human expresses goals with
budgets attached, decides what verification could not settle, and asks the world *why*
anything happened. A provenance walk over the log answers that question.

## What it runs on

- **Local first.** One daemon (`ailang-worldd`), SQLite and the local filesystem, with no
  cloud dependency in the core. The REST listener is loopback-only.
- **Protocol native.** Transitions are served as tools over **MCP** and published as skills
  on an **A2A** agent card. World adds no new wire protocol.
- **Agents as residents.** A software-engineering domain serves eight AILANG coding tools
  (read, write, edit, check, run, builtins search, examples search, CLI) over MCP and A2A.
  Every tool's effect goes through the broker and commits one log entry.

## Where it stands

World is **pre-1.0**. Five of the seven clauses in its 1.0 bar are met; the two that remain
are the ones that measure value against the status quo. See the [roadmap](./roadmap/index.md)
for the honest state of each clause.

## Where to go next

- [Concepts](./concepts/index.md): the world graph, transitions, capabilities and replay.
- [Get started](./getting-started/index.md): build the daemon and commit a first transition.
- [For AI agents](./agents/index.md): connect an agent over MCP or A2A.
- [AILANG](https://ailang.sunholo.com): the language World's transactions are written in.

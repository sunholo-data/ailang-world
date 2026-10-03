---
title: Roadmap to 1.0
sidebar_position: 1
description: The seven clauses AILANG World 1.0 must meet, and the current state of each.
---

# Roadmap to 1.0

World 1.0 is defined by a **bar** of seven clauses, ratified by the project owner on
2026-07-23. Each clause is a checkable end state, not a work item. The authoritative text
lives in the
[mission charter](https://github.com/sunholo-data/ailang-world/blob/dev/design_docs/world-mission.md)
("The bar"). This page summarizes it and says where each clause stands.

**Current state: 5 of 7 clauses met.** The two that remain are the ones that measure World
against the status quo. No release date is promised.

| # | Clause | State |
|---|---|---|
| 1 | Deterministic kernel | <span className="world-status world-status--met">Met</span> |
| 2 | Local-first daemon | <span className="world-status world-status--met">Met</span> |
| 3 | Explicit authority end to end | <span className="world-status world-status--met">Met</span> |
| 4 | Resident-agent non-inferiority floor | <span className="world-status world-status--unmet">Not yet met</span> |
| 5 | The human surface works, with provenance teeth | <span className="world-status world-status--unmet">Not yet met</span> |
| 6 | Protocol-native boundary | <span className="world-status world-status--met">Met</span> |
| 7 | Controlled self-modification, proven once | <span className="world-status world-status--met">Met</span> |

## The clauses

### 1. Deterministic kernel <span className="world-status world-status--met">Met</span>

An immutable, content-addressed world store with an append-only transition log and **proven
deterministic replay**. Replaying any recorded episode reconstructs state bit-for-bit from
pure transitions plus recorded effect results. The log-epoch semantics (interpreter version
pinning and content-addressed transition functions) were decided and enforced before the log
format froze.

### 2. Local-first daemon <span className="world-status world-status--met">Met</span>

`ailang-worldd` runs on one machine over SQLite, with a CLI and a REST API, and has **zero
cloud dependencies in the core**. Cloud transports may exist only as effect-handler
extensions.

### 3. Explicit authority end to end <span className="world-status world-status--met">Met</span>

Every effect goes through the broker with a capability and budget check. Effect results are
recorded as replay input. Capsules run with a physical isolation floor beneath the semantic
checks. No ambient-authority path exists from an agent to the outside world.

### 4. Resident-agent non-inferiority floor <span className="world-status world-status--unmet">Not yet met</span>

This is a do-no-harm kill switch, not the value proof. Two reference agents from different
providers, **Claude Code and codex CLI**, each run the same benchmark set in two arms: *shell*
(native tools) and *World* (MCP transition tools only), with the same model, three or more
runs per arm. The floor holds only if **both** agents show a World pass rate no more than
2 percentage points below the shell pass rate **and** a median wall-clock overhead of at most
25%. A stability precondition on the shell arm keeps a flaky harness from faking the verdict
either way. If the floor fails on eligible agents after honest tuning, World parks, and says
so.

**Where it stands.** Tracked as queue row 93. Its prerequisite, row 134 (the
software-engineering tools agents use in the World arm), landed on 2026-10-03, and the
floor's grader port, classification and statistics landed the same day (row 93 M1 and M3,
no agent run yet). Row 135 widened the `ailang-run` tool (stdin, argv and per-task
capabilities) so all 23 core tasks can run through World; it landed on 2026-10-03 on the
v0.52.1 tool binary (#192), with an attended four-task smoke still to run. The ratified order
to release is now the rest of the
developer CLI (row 138), a toolchain effect for non-AILANG projects (row 140), the floor run
(row 93), provenance teeth (row 114) and an interpreter epoch-upgrade path (row 139).

### 5. The human surface works, with provenance teeth <span className="world-status world-status--unmet">Not yet met</span>

This is the 1.0 value demonstration. One real goal is expressed, becomes proposals, reaches
an approval inbox with evidence bundles, and commits. And, measurably: on **at least three
real "why did this happen" questions**, arising from actual operation rather than written for
the test, a provenance walk yields the verified answer in **five minutes or less** each,
where the pre-World method was a grep and log archaeology session.

**Where it stands.** Tracked as queue row 114. The walks so far proved retrieval of answers
that had been recorded in advance, and the pre-World baselines were not timed. The clause
needs unseeded questions from real operation, with the pre-World method timed alongside. Its
corpus is the operation World records now that row 134 has landed.

### 6. Protocol-native boundary <span className="world-status world-status--met">Met</span>

The transition registry is served over MCP, filtered per session by capability, and an A2A
agent card is published. **No new wire protocols.**

### 7. Controlled self-modification, proven once <span className="world-status world-status--met">Met</span>

At least one World behavior change ships as an extension package through World's own
propose, verify and commit pipeline, with the self-modification boundaries enforced: the
mission-loop machinery, the AILANG compiler and the live daemon are outside self-modification
scope.

## Beyond 1.0

The standing value evidence beyond the bar is the mission loop that builds this repository
migrating onto World, measured on incident classes eliminated and human attention per
workstream against the July 2026 markdown-and-scheduler baseline.

## Following along

Every mission iteration reports on the project's weekly
[bookkeeping issue](https://github.com/sunholo-data/ailang-world/issues?q=is%3Aissue+%22mission+bookkeeping%22),
and the full
per-iteration log lives in
[`design_docs/world-mission-log.md`](https://github.com/sunholo-data/ailang-world/blob/dev/design_docs/world-mission-log.md).

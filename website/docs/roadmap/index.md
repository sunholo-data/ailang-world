---
title: The 1.0 bar and release
sidebar_position: 1
description: The seven clauses AILANG World 1.0 had to meet, the evidence for each, and the remaining release steps.
---

# The 1.0 bar and release

World 1.0 is defined by a **bar** of seven clauses, ratified by the project owner on
2026-07-23. Each clause is a checkable end state, not a work item. The authoritative text
lives in the
[mission charter](https://github.com/sunholo-data/ailang-world/blob/dev/design_docs/world-mission.md)
("The bar"). This page summarizes it and says where each clause stands.

**Current state: all 7 clauses met; `v1.0.0-rc.1` is the current release candidate.** Clause 4
was closed by D-WORLD-77 and clause 5 by D-WORLD-80. Under the release plan (D-WORLD-81),
`v1.0.0` follows once three more queue rows land: 156 (protected test paths), 150 (marketing
capture) and 151 (agent-generated projections). See [Release plan](#release-plan).

| # | Clause | State |
|---|---|---|
| 1 | Deterministic kernel | <span className="world-status world-status--met">Met</span> |
| 2 | Local-first daemon | <span className="world-status world-status--met">Met</span> |
| 3 | Explicit authority end to end | <span className="world-status world-status--met">Met</span> |
| 4 | Resident-agent non-inferiority floor | <span className="world-status world-status--met">Met</span> |
| 5 | The human surface works, with provenance teeth | <span className="world-status world-status--met">Met</span> |
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

### 4. Resident-agent non-inferiority floor <span className="world-status world-status--met">Met</span>

This is a do-no-harm kill switch, not the value proof. Two reference agents from different
providers, **Claude Code and codex CLI**, each run the same benchmark set in two arms: *shell*
(native tools) and *World* (MCP transition tools only), with the same model, three or more
runs per arm. The floor holds only if **both** agents show a World pass rate no more than
2 percentage points below the shell pass rate **and** a median wall-clock overhead of at most
25%. A stability precondition on the shell arm keeps a flaky harness from faking the verdict
either way. If the floor fails on eligible agents after honest tuning, World parks, and says
so.

**Where it stands.** Met (queue row 93, D-WORLD-77). The pre-registered FINAL run of
2026-10-15 holds:

| Agent | Pass-rate difference (World − shell) | Median overhead |
|---|---|---|
| Claude Code | 0 (69/69 vs 69/69) | −17.9% |
| codex | 0 (69/69 vs 69/69) | −5.7% |

All 414 graded cells passed (both agents, 23 core tasks, three runs, the shell eligibility
phase and both arms of the paired phase). The rows it depended on all landed first: the coding
tools (row 134), the wider `ailang-run` (row 135), the developer CLI (row 138) and
`workspace-exec` (row 140). Evidence:
[`design_docs/verification/world-floor-final-2026-10-15/`](https://github.com/sunholo-data/ailang-world/blob/dev/design_docs/verification/world-floor-final-2026-10-15/README.md).

### 5. The human surface works, with provenance teeth <span className="world-status world-status--met">Met</span>

This is the 1.0 value demonstration. One real goal is expressed, becomes proposals, reaches
an approval inbox with evidence bundles, and commits. And, measurably: on **at least three
real "why did this happen" questions**, arising from actual operation rather than written for
the test, a provenance walk yields the verified answer in **five minutes or less** each,
where the pre-World method was a grep and log archaeology session.

**Where it stands.** Met (queue row 114, D-WORLD-80). Three unseeded questions harvested from
real operation were answered by World walks verified against sealed ground truth: Q4 in
54.6 s, Q6 in 25.3 s and Q8 in 39.0 s. A live pre-World baseline was timed on each and also
verified, more slowly (67.0 s, 87.5 s and 116.8 s). Attempts that did not qualify are kept in
the record. Evidence:
[`design_docs/verification/w-provenance-teeth-unseeded-baselined-walk/`](https://github.com/sunholo-data/ailang-world/blob/dev/design_docs/verification/w-provenance-teeth-unseeded-baselined-walk/index.md).

### 6. Protocol-native boundary <span className="world-status world-status--met">Met</span>

The transition registry is served over MCP, filtered per session by capability, and an A2A
agent card is published. **No new wire protocols.**

### 7. Controlled self-modification, proven once <span className="world-status world-status--met">Met</span>

At least one World behavior change ships as an extension package through World's own
propose, verify and commit pipeline, with the self-modification boundaries enforced: the
mission-loop machinery, the AILANG compiler and the live daemon are outside self-modification
scope.

## Release plan

`v1.0.0-rc.1` is the current GitHub pre-release. Under D-WORLD-81, `v1.0.0` follows after
these rows:

- **Row 156, protected test paths** (next). Today an agent can edit the tests that
  `workspace-exec` runs to judge its work; this row protects those paths.
- **Row 150, marketing capture.**
- **Row 151, agent-generated projections.**

Row 139, the interpreter epoch-upgrade path, is in progress: its release manifest
(`pins/ailang-releases.json`) and `world-publish epoch nominate` have landed.

## Beyond 1.0

The standing value evidence beyond the bar is the mission loop that builds this repository
migrating onto World, measured on incident classes eliminated and human attention per
workstream against the July 2026 markdown-and-scheduler baseline.

## All agent work through World

The project owner's goal (D-WORLD-94, 2026-10-26): every agent task runs through World,
headless or interactive, on any repo, with nothing used alongside it. The roadmap is
[`w-all-agent-work-through-world.md`](https://github.com/sunholo-data/ailang-world/blob/dev/design_docs/planned/w-all-agent-work-through-world.md),
rows 217–227, in this order (its §3):

| Row | Item | State |
|---|---|---|
| 223 | Agent host adapters: `hosts init` launchers and `doctor` classes (223a), the herdr plugin (223b), key bindings and the World panel (223c), session context and `whoami` (223d) | 223a–223c landed; 223d merging |
| 227 | Error clarity: every failure says what, why and the next command | in part, with 223a |
| 226 | Value scorecard per episode | planned |
| 219 | World as a supervised user service | planned |
| 217 | One-command episode start from a preset | planned |
| 218 | `episode finish`: checks, commit, PR | planned |
| 220 | The surface index follows merges and running episodes | planned |
| 221 | Confined dev tools in exec profiles, with row 199 | planned |
| 224 | Standing session allowances | planned |
| 225 | `ailang-worldd run`: a wrapper for agents with no adapter | planned |
| 222 | An attended session as an episode | planned |

Rows 206 and 207 (the live agent tab and agent states) follow 222, then 196 (the mission loop's
executors) when its ruling exists. To use what has landed, see
[Run your agents through World](../guides/agents-through-world.md).

| Ruling | Resolved 2026-10-27 as |
|---|---|
| D-WORLD-95 | Rows 206 and 207 are un-parked, after row 222. Rows 219, 217 and 218 come before row 188's N2d. |
| D-WORLD-96 | Standing allowances (row 224): at most 10 sessions a day, 12 h each, approved for 30 days, revocable, never for the mission loop. |
| D-WORLD-97 | Confinement stays per launch. It goes machine-wide only when the scorecard (row 226) holds over 10 episodes with zero unclear errors, the daemon is a supervised service (219), and the loop and attended sessions run in World (196, 222). |

## Following along

Every mission iteration reports on the project's weekly
[bookkeeping issue](https://github.com/sunholo-data/ailang-world/issues?q=is%3Aissue+%22mission+bookkeeping%22),
and the full
per-iteration log lives in
[`design_docs/world-mission-log.md`](https://github.com/sunholo-data/ailang-world/blob/dev/design_docs/world-mission-log.md).

---
title: Provenance and replay
sidebar_position: 5
description: What every commit records, how the records link together, what replay can and cannot reconstruct today, and the interpreter pin that makes it possible.
---

# Provenance and replay

## What a tool call leaves behind

Every MCP `tools/call` or A2A `tasks/send` of a published transition commits exactly one log
entry, including calls that were refused by the transition's own plan. For an effectful
transition (all eight coding tools) the commit writes these objects:

| Object (`semanticId`) | Contents |
|---|---|
| `world/invocation-input/v1` | The caller's arguments |
| `world/effect-plan/v1` | The plan the pure transition returned |
| `world/invocation-output/v1` | The result returned to the caller; also the new world's `stateRoot` |
| `world/invocation-record/v2` | `invocationId`, `episodeId`, `skillId`, `transitionFn`, `interpreter`, `semanticsEpoch`, `input`, `output`, `plan`, and the ordered `effects` (effect record refs) |

and the broker has already written, before the commit:

| Object | Contents |
|---|---|
| `world/effect-request/v1` | The effect, scope, cost, decision time and the request payload |
| `world/effect-record/v1` | The decision and its accounting: `allowed`, `failed`, `denial`, `budgetBefore`, `budgetAfter`, `requestRef`, `resultRef` |
| the result object | The handler's output bytes |

Plus the journal's effect intent and outcome. A pure transition writes the
`world/invocation-record/v1` form, which has no `plan` or `effects`.

The log entry ties it together:

```text
log entry (entryIndex N)
  header.transitionFn ─▶ world/transition-source/v1   (the AILANG source that ran)
  header.interpreter  ─▶ archived interpreter         (the exact binary, by hash)
  transitionRef       ─▶ world/invocation-record/v2
                           input   ─▶ world/invocation-input/v1
                           output  ─▶ world/invocation-output/v1
                           plan    ─▶ world/effect-plan/v1
                           effects ─▶ world/effect-record/v1
                                        requestRef ─▶ world/effect-request/v1
                                        resultRef  ─▶ handler output
```

The call's result also carries `world.plan` and `world.effects[].record`, so a caller holds the
two refs it needs to start a walk. How to walk it with the API is in
[Provenance walks](../guides/provenance-walks.md).

## The interpreter pin

`serve --ailang-bin PATH` archives the interpreter by content hash next to the store. Its hash
is reported as `interpreter_ref` by `GET /v1/health`, and every log entry the daemon writes
carries it in `header.interpreter`. Registry descriptors pin the same hash. Before each run the
capsule re-hashes the archived bytes. This is the replay pin: the exact binary that ran a
transition is named by every record of it.

The interpreter's release also selects the **semantics epoch** through the epoch registry
(`world/epoch-registry/v1`), which the daemon bootstraps on first start for the release it
serves.

## Replay

Replay of an effectful invocation is defined as:

1. re-run the plan phase and require it to equal the recorded plan object byte for byte;
2. feed the same requests to a broker **replay** session over the recorded effect records — it
   returns recorded results or denials and never calls a handler;
3. re-run the finish phase, if any, on the reconstructed results;
4. require the recomposed output to equal the recorded output byte for byte.

:::warning Replay has no operator entrypoint yet
`Coordinator.Replay` and the replay binder exist and are exercised by tests, but no daemon
route or CLI verb runs a replay today (residual R-SE-14). Claims that depend on an operator
replaying an effectful invocation are not yet exercisable outside tests.
:::

What replay reconstructs is **World's record**, not the worktree:

- Write and edit payloads are kept in full in the request objects, but the worktree files are
  not in World's content store (residual R-SE-7).
- An `ailang-run` program's own file operations are captured only as the run envelope
  (`admitted`, `exit_code`, `decision`, `limit`, `stdout`, `stderr`), not effect by effect
  (residual R-SE-11).

## Integrity scan

At startup the daemon scans persisted log and world rows within a 20,000-row and 2 s budget. A
clean, complete scan prints nothing extra. It prints `integrity_hole …` lines for unreadable
historic rows, and `integrity_scan_incomplete …` with resume cursors when a budget stopped the
scan; an incomplete scan is not a clean bill of health. Committed history is never rewritten
to repair a hole; affected reads keep failing loudly.

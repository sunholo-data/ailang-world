---
title: World and transitions
sidebar_position: 2
description: The content-addressed world store, the append-only log, transitions as pinned AILANG programs, and the propose-verify-commit path.
---

# World and transitions

## The store

The daemon keeps everything in one SQLite file (`serve --db`). It holds:

- **Objects.** Content-addressed blobs. Each object has a `hash` (`sha256:<hex>` of its payload
  bytes), an `interfaceHash`, a `semanticId` (for example `world/effect-record/v1`) and a
  `provenance` string. The store verifies content against hash on write. Objects are never
  modified.
- **The log.** An append-only sequence of entries. Each entry has a frozen six-field header
  (`entryIndex`, `semanticsEpoch`, `transitionFn`, `interpreter`, `prevEntryHash`,
  `writtenBy`), plus its `entryHash` and a `transitionRef` that points at the object recording
  what the transition did.
- **Worlds.** Rows of `{ref, revision, stateRoot, logHead}`. One of them is the **selected
  head**.
- **Registry heads.** Named pointers such as `world/epoch-registry/v1` and
  `world/transition-registry/v1`.
- **The journal.** Durable intents and outcomes for commits and effects, used for receipts and
  budget accounting.

One process may hold write authority for a store file at a time. This is enforced with a
non-waiting OS lock on `<db>.writer.lock`, not by convention: a second daemon, or a
`session mint` or `world-publish` against the same file while the daemon runs, is refused.

## A commit

A commit is a compare-and-append. The caller states the head it observed, the new objects, the
next world row, and the new log entry. The store accepts it only if the observed head is still
the selected head; otherwise it answers with both heads so the caller can re-plan.

```json
{
  "observedHead": "",
  "objects": [{"hash": "sha256:…", "interfaceHash": "sha256:…",
               "semanticId": "world/demo/genesis-goal", "provenance": "quickstart",
               "payload": "<base64>"}],
  "nextWorld": {"ref": "sha256:…", "revision": 0, "stateRoot": "sha256:…", "logHead": "sha256:…"},
  "entry": {"header": {"entryIndex": 0, "semanticsEpoch": 1, "transitionFn": "sha256:…",
                       "interpreter": "sha256:…", "prevEntryHash": "sha256:…", "writtenBy": "quickstart"},
            "entryHash": "sha256:…", "transitionRef": "sha256:…"}
}
```

`observedHead` is empty only for genesis. The full shape and every error class are in the
[HTTP API reference](../reference/http-api.md#post-v1commit). A hand-written commit like this
is how you create the genesis world; after that, transitions commit for you.

## Transitions

A transition is an AILANG module published to the **transition registry**. A published
descriptor pins:

| Field | Meaning |
|---|---|
| `id` | Stable transition ID; also the MCP tool name (with a reversible encoding for `_ . /`) and the A2A skill ID |
| `transitionFn` | Content hash of the canonical source object (`world/transition-source/v1`) |
| `interpreter` | Content hash of the archived AILANG interpreter that runs it |
| `semanticsEpoch` | The epoch the epoch registry nominates for that interpreter release |
| `inputSchema`, `outputSchema` | JSON Schemas carried to MCP and A2A clients |
| `access` | The `{effect, scope, cost}` a session must hold to see and call the transition |
| `declaredEffects` | The exact `{effect, scope, cost}` triples the transition may request |

Publishing runs the pinned interpreter's `check` on the source in an empty scratch root, so a
transition must be **self-contained**: it may import `std/*` but not `world/*` (that fails with
`LDR001`). See the [manifest reference](../reference/transitions-manifest.md).

### The calling convention

Every transition exports `main(input: string) -> string`.

- A **pure** transition (no `declaredEffects`) receives the caller's JSON arguments as a string
  and returns a JSON object as a string. That object is the result.
- An **effectful** transition is called in two phases. It first receives
  `{"phase":"plan","args":…}` and returns an **effect plan**
  (`world/effect-plan/v1`). If the plan asks for it, it is called again with
  `{"phase":"finish","args":…,"results":[…]}` after the effect has run.

An effect plan looks like this (from the design, and checked by the plan parser):

```json
{"plan":"world/effect-plan/v1",
 "effects":[{"id":"e1","effect":"Workspace.Read","scope":"worktree","cost":1,
             "payload":{"op":"read","path":"src/main.ail"}}],
 "finish":false,
 "result":null}
```

The host enforces the plan laws: at most one effect per plan in v1, IDs unique and
`[a-z0-9]{1,16}`, each `(effect, scope, cost)` must be one of the descriptor's
`declaredEffects`, the payload is a JSON object of at most 1 MiB, a zero-effect plan must
carry a non-null `result`, and the key `world` is reserved in outputs. A refusal is a valid
pure plan (`"effects":[]` with `"result":{"ok":false,"refused":"…"}`) and still commits: a bad
call is provenance too.

## Propose, verify, commit — as implemented

For one MCP `tools/call` or A2A `tasks/send`:

1. **Admit.** Resolve the bearer session, build a fresh snapshot of the registry and the
   session's grants, and check the requested ID is among the transitions this session may see.
2. **Bind.** Check the descriptor's `access` against the grants.
3. **Plan.** Run the transition's plan phase in the capsule with the pinned interpreter and no
   capabilities (`--caps ""`), bounded to 2 s.
4. **Execute.** Send the planned request through the broker: capability, scope, liveness and
   budget are checked, the intent is journaled, the handler runs, the outcome and an effect
   record are written.
5. **Finish.** If requested, run the finish phase (2 s).
6. **Commit.** Write the input, output, plan and invocation record (`world/invocation-record/v2`)
   objects and append one log entry, under the per-episode lock.

Every phase shares one 20 s invocation deadline. Everything after the effect runs on its own
4 s budget so a completed effect is never stranded untyped. If the commit fails after the
effect ran (for example another writer moved the head), the caller receives
`EffectsUnrecordedError` naming the effect record refs; the effect's intent, outcome and record
are already durable. See [MCP and A2A](../reference/mcp-and-a2a.md#error-mapping).

:::note What "verify" means today
Verification today is the type check at publish time, the plan laws, the authority and budget
checks, and the store's compare-and-append. The broader Verify phase in the design (Z3
contracts over proposals, simulation, quorum evidence for high-risk transitions) is design
intent, not something every call runs.
:::

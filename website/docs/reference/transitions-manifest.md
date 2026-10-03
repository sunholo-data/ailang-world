---
title: Transitions manifest
sidebar_position: 4
description: The descriptor manifest format world-publish transitions reads, how each field is pinned and validated, and worked examples from se-tools and the quickstart echo.
---

# Transitions manifest

`world-publish transitions` publishes a **manifest** — a JSON array of descriptors — into a
store's transition registry (`world/transition-registry/v1`). The checked-in example is
`packages/se-tools/transitions.json`.

```bash
world-publish transitions --store <db> --manifest <file.json> (--ailang-bin <PIN> | --interpreter-ref <sha256:…>)
```

This is a **local, attended** registry write. The daemon must be stopped, and the verb asks for a
typed confirmation at a controlling terminal (see [Attended steps](../guides/attended-steps.md)).
`--live` is refused.

## Fields

| Field | Required | Meaning |
|---|---|---|
| `id` | yes | The transition ID (validated by the registry): A2A skill ID verbatim, MCP tool name after escaping `_ . /` |
| `title` | recommended | Shown as the A2A skill `name` |
| `description` | recommended | Shown to agents on the card and in `tools/list` |
| `transitionFnFile` | one of these two | Path to the AILANG source file (relative to the current directory) |
| `transitionFn` | one of these two | Ref of a source object already in the store |
| `inputSchema` | yes | JSON Schema for the arguments; canonicalised |
| `outputSchema` | yes | JSON Schema for the result; canonicalised |
| `access` | yes | `{effect, scope, cost}` a session's grants must satisfy to see and call it |
| `declaredEffects` | yes (may be `[]`) | The exact `{effect, scope, cost}` triples its plan may request |
| `semanticsEpoch` | no | Omit it: it is derived from the epoch registry. If given, it must be one of the epochs nominating the interpreter's release |

There is deliberately **no interpreter field**. The interpreter pin comes from the command line:
`--ailang-bin B` archives `B` next to the store exactly as the daemon does at startup;
`--interpreter-ref R` names an interpreter already archived there. Exactly one is required, so a
descriptor can never name an interpreter nobody verified.

## What publishing checks

For each entry, in order:

1. **Epoch.** The interpreter's release must be nominated by an epoch in
   `world/epoch-registry/v1`. The daemon bootstraps epoch 1 for the release it serves, so start
   `serve --ailang-bin <same binary>` once against the store before publishing. Omitted
   `semanticsEpoch` is derived when exactly one epoch nominates the release, and refused as
   ambiguous otherwise. It is never defaulted.
2. **Source.** `transitionFnFile` is canonicalised and **checked** under the pinned interpreter
   (a bounded `check` in an empty scratch root) **before** it is stored as a
   `world/transition-source/v1` object. A source that does not load is never stored. Because the
   scratch root is empty, imports of `world/*` fail with `LDR001`; publishable modules must be
   self-contained and import only `std/*`.
3. **Schemas** are canonicalised.
4. **The registry write** is a compare-and-swap on the registry head. The publisher re-verifies
   every source exists and loads under its own pinned interpreter before moving the head.

Outcomes:

- `published transition registry revision N (head sha256:…); the daemon's A2A card will list these skills on its next read`
- `transition registry UNCHANGED at revision N (head sha256:…); nothing written` — an identical
  republish writes nothing.
- `semantics epoch 1 derived from world/epoch-registry/v1 for interpreter release "…"` precedes
  either when the epoch was derived.
- A same-ID conflict with different bytes, or a head that kept moving, reports that nothing was
  written.

Publication establishes that the source object exists, passes `check` under the pinned
interpreter, and has an epoch the registry nominates. It does **not** establish that the source
implements the calling convention; invocation reports an incompatibility if it does not.

## The calling convention

Every transition exports `main(input: string) -> string` and returns a JSON object as a string.

- **Pure** (`declaredEffects: []`): `input` is the caller's arguments as JSON.
- **Effectful**: `input` is `{"phase":"plan","args":…}`; the result is a
  `world/effect-plan/v1` object. If the plan sets `"finish": true`, the transition is called again
  with `{"phase":"finish","args":…,"results":[…]}`. See
  [World and transitions](../concepts/world-and-transitions.md#the-calling-convention).

## Example: a pure transition

From `docs/QUICKSTART.md` §6 (marked "attended — pending first verbatim run" there). The source:

```ailang
module quickstart/echo

export func main(input: string) -> string {
  "{\"echo\":${input}}"
}
```

The manifest:

```json
[{"id": "tools.echo", "title": "Echo", "description": "quickstart transition",
  "transitionFnFile": "/tmp/echo.ail",
  "inputSchema": {"type": "object"}, "outputSchema": {"type": "object"},
  "access": {"effect": "world.apply", "scope": "world", "cost": 1},
  "declaredEffects": []}]
```

A session needs a grant such as `world.apply=world:10` to see it. Its MCP name is
`tools_decho`.

## Example: an effectful tool

The first entry of `packages/se-tools/transitions.json`, with the description shortened:

```json
{
  "id": "ailang-read",
  "title": "Read (sandboxed)",
  "description": "Read a file inside this episode's worktree (policy-tool `read`). …",
  "transitionFnFile": "packages/se-tools/se_tools/read.ail",
  "inputSchema": {
    "type": "object",
    "properties": {"path": {"type": "string", "description": "File path, relative to the worktree root"}},
    "required": ["path"],
    "additionalProperties": false
  },
  "outputSchema": {
    "type": "object",
    "properties": {
      "content": {"type": "string"},
      "ok": {"type": "boolean"},
      "refused": {"type": "string", "description": "Why the call was refused (plan, policy or handler)"},
      "world": {"type": "object", "description": "World's provenance: the plan ref and the effect record refs"}
    },
    "required": ["world"]
  },
  "access": {"effect": "Workspace.Read", "scope": "worktree", "cost": 0},
  "declaredEffects": [{"effect": "Workspace.Read", "scope": "worktree", "cost": 1}]
}
```

Note the pattern: `access` at cost 0 and one declared effect at cost 1, so one grant per effect
name both admits the tool and pays for each call. `transitionFnFile` paths are relative to the
**repo root**, so run the publish from there.

All eight `se-tools` entries:

| `id` | Source | `access.effect` | Required args | Optional args |
|---|---|---|---|---|
| `ailang-read` | `se_tools/read.ail` | `Workspace.Read` | `path` | |
| `ailang-write` | `se_tools/write.ail` | `Workspace.Write` | `path`, `content` | |
| `ailang-edit` | `se_tools/edit.ail` | `Workspace.Write` | `path`, `old_text`, `new_text` | |
| `ailang-check` | `se_tools/check.ail` | `Ailang.Check` | `path` | |
| `ailang-run` | `se_tools/run.ail` | `Ailang.Run` | `path` | `args_json` |
| `builtins-search` | `se_tools/builtins_search.ail` | `Ailang.Discover` | | `query`, `module` |
| `examples-search` | `se_tools/examples_search.ail` | `Ailang.Discover` | `query` | |
| `ailang-cli` | `se_tools/cli.ail` | `Ailang.CLI` | `op` | `path`, `module`, `query`, `package`, `flags` |

Every input schema sets `"additionalProperties": false`, and every tool's plan also refuses
unknown keys itself, because the underlying `policy-tool` silently ignores them.

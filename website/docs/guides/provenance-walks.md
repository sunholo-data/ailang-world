---
title: Provenance walks
sidebar_position: 3
description: Answer "why did this happen" by walking World's log, records and objects through the read API, instead of searching logs.
---

# Provenance walks

A provenance walk answers a "why did X happen" question by following content-addressed
references from something you can see (a tool result, a log entry, a named object) back to the
facts that produced it. Every step is a plain `GET`; no session is needed for reads.

World's 1.0 bar asks that, for at least three real questions arising from operation, a walk
yields the verified answer in five minutes or less. The first demonstration (2026-09-24) answered
four real mission incidents in 10–20 s each. Its honest caveat: those incidents had been
diagnosed before they were committed, so the walks measured **retrieval** of a recorded
diagnosis, not a fresh one, and the pre-World baselines were unmeasured. A follow-up protocol
with unseeded questions and live-timed baselines (queue row 114) is not yet run.

## The read API you walk with

| Need | REST | CLI |
|---|---|---|
| Current head | `GET /v1/head` | `ailang-worldd head` |
| One log entry | `GET /v1/log/{index}` | `ailang-worldd log get N` |
| A page of entries | `GET /v1/log?from=N&limit=M` | `ailang-worldd log range --from N --limit M` |
| A world row | `GET /v1/worlds/{ref}` | `ailang-worldd world get REF` |
| An object's metadata | `GET /v1/objects/{ref}` | `ailang-worldd object get REF` |
| An object's payload | `GET /v1/objects/{ref}?payload=true` (base64) | `ailang-worldd object get REF --payload` |
| Objects by semantic ID | `GET /v1/objects/by-semantic-id/{name}?after=REF&limit=N` | `ailang-worldd object find NAME --after REF --limit N` |
| A `rest:` commit's receipt | `GET /v1/receipts/rest:{id}` | — |
| Browse references in HTML | `GET /workbench?object=REF` | — |

Pages default to 100 items and cap at 500. Objects by semantic ID come back in ascending hash
order, without payloads; fetch the one you want by ref. A semantic ID nothing carries returns
`200` with an empty list.

The payload is base64 JSON. A small helper keeps the walk readable:

```bash
U=http://127.0.0.1:7644
payload() { curl -s "$U/v1/objects/$1?payload=true" |
  python3 -c 'import sys,json,base64; print(base64.b64decode(json.load(sys.stdin)["payload"]).decode())'; }
```

## Walk 1: from a tool result to everything behind it

A coding-tool result carries `world.plan` and `world.effects[].record`. Start from the record:

```bash
payload sha256:<record>
```

```json
{"effect":"Workspace.Read","scope":"worktree","cost":1,"budgetBefore":50,"budgetAfter":49,
 "allowed":true,"failed":false,"denial":"","requestRef":"sha256:…","resultRef":"sha256:…"}
```

That answers "was this effect authorized, under which grant, and what did it cost". Follow
`requestRef` for exactly what was asked (effect, scope, cost, decision time, payload) and
`resultRef` for exactly what the handler returned:

```bash
payload sha256:<requestRef>
payload sha256:<resultRef>
payload sha256:<plan>          # what the pure transition planned
```

To reach the **log entry**, find the invocation record that lists this effect record. Records
for effectful calls have semantic ID `world/invocation-record/v2`:

```bash
for h in $(curl -s "$U/v1/objects/by-semantic-id/world/invocation-record/v2?limit=500" |
           python3 -c 'import sys,json; [print(i["hash"]) for i in json.load(sys.stdin)["items"]]'); do
  payload "$h" | grep -q 'sha256:<record>' && echo "$h"
done
```

The invocation record names `invocationId` (`mcp:<episode>:<task>` or `a2a:<episode>:<task>`), `episodeId`, `skillId`,
`transitionFn`, `interpreter`, `semanticsEpoch`, `input`, `output`, `plan` and `effects`. Its hash
is the `transitionRef` of exactly one log entry; find that entry by scanning
`log range` for it, or open `/workbench?object=<record-hash>` to see the entries that reference
it. The entry's `header.transitionFn` is the AILANG source that ran, and `header.interpreter` is
the exact interpreter binary.

:::tip MCP results do not carry the entry index
An A2A `tasks/send` result includes `metadata.entry_index`, so you can go straight to
`GET /v1/log/{entry_index}`. An MCP `tools/call` result does not; use the record walk above.
:::

## Walk 2: from a log entry forward

```bash
ailang-worldd log get 5
payload sha256:<transitionRef>      # the invocation record
payload sha256:<input>              # the arguments the agent sent
payload sha256:<output>             # what it got back
```

For a pure transition the record is `world/invocation-record/v1`, with no `plan` or
`effects`. A refused call has a plan with no effects and an output with `ok:false` and a
`refused` reason: refusals are provenance too.

## Walk 3: named facts (incidents, decisions, evidence)

Anything committed with a stable semantic ID can be found without scanning the log. The
clause-5 demonstration committed each incident as an object with semantic ID
`world/mission/incident/<slug>`, whose payload holds `question`, `answer` and `sources[]` (refs
to evidence objects with IDs like `world/mission/evidence/<slug>/0`):

```bash
ailang-worldd object find world/mission/incident/iter181-ci-bench-401
payload sha256:<incident-hash>      # question, answer, sources[]
payload sha256:<source-hash>        # each piece of evidence
```

Before the `by-semantic-id` route existed (finding F-1 of that demonstration), a walk had to scan
`log range --from 0` and fetch each entry's `transitionRef` metadata until the semantic ID
matched. That still works, and is the fallback when you do not know the name.

## Verify outside World

An answer the walk produced is only **verified** when you confirm it against a source the walk
did not produce: a binary, the git history, an upstream issue. The demonstration recorded an
unconfirmed answer as `UNVERIFIED` and did not count it. Keep that discipline: World makes the
record easy to reach, and the record can still be wrong.

## Commit receipts

For a commit you sent yourself with `"invocationId": "rest:<id>"`,
`GET /v1/receipts/rest:<id>` answers `resolved` (with `resultRef`), `not-started` or
`indeterminate`. Only the `rest:` namespace is exposed. Treat `indeterminate` as fail-closed:
reconcile before you resend.

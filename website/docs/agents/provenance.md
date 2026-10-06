---
title: Provenance
sidebar_position: 6
description: How an AI agent can read back and explain its own past actions in AILANG World through the log and object routes.
---

# Provenance: explaining your own past actions

Every tool call you make commits one entry to World's append-only log, and every entry links, by content address, to what you asked for, what the transition planned, what the broker executed, and what you got back. When your operator asks "what did you do, and why did it fail?", you can answer from the record instead of from memory.

The examples below are the real records of pi's `ailang-run` call from the attended smoke on 2026-10-03 (log entry 5). Long digests are shortened to eight characters.

## The read routes

| Route | Returns |
|---|---|
| `GET /v1/head` | The selected world ref, as plain text |
| `GET /v1/log/<index>` | One log entry |
| `GET /v1/log?from=<n>&limit=<m>` | `{"items": […]}`, entries from `n` (default page 100, maximum 500) |
| `GET /v1/objects/<ref>` | Object metadata: `hash`, `interfaceHash`, `semanticId`, `provenance` |
| `GET /v1/objects/<ref>?payload=true` | The same plus `payload`, **base64-encoded** |
| `GET /v1/objects/by-semantic-id/<id>` | A page of objects with that semantic id (no payloads; `?after=<ref>`, `?limit=`) |

**Authentication.** In this version only `POST /v1/commit` checks the session credential; the `GET /v1/*` routes above are open on the daemon's address (a recorded residual, with read enforcement left to a later change). Send your `Authorization: Bearer` header anyway, so your client keeps working if that changes. Never call `POST /v1/commit`: tool calls commit for you.

```bash
H="Authorization: Bearer $WORLD_SESSION"
curl -s -H "$H" http://127.0.0.1:7644/v1/log/5
curl -s -H "$H" "http://127.0.0.1:7644/v1/objects/sha256:<ref>?payload=true" \
  | python3 -c 'import sys,json,base64; print(base64.b64decode(json.load(sys.stdin)["payload"]).decode())'
```

## Walking one call

**1. The log entry.** `transitionRef` is the invocation record.

```json
{"header": {"entryIndex": 5, "semanticsEpoch": 1,
            "transitionFn": "sha256:371fa1dc…", "interpreter": "sha256:1a67b014…",
            "prevEntryHash": "sha256:5261510b…", "writtenBy": "coordinator:mcp"},
 "entryHash": "sha256:c3710b76…",
 "transitionRef": "sha256:98fe1ad1…"}
```

`transitionFn` is the content address of the tool's AILANG source, and `interpreter` is the pinned interpreter that ran its plan. `writtenBy` names the surface the call arrived on: `coordinator:mcp` for an MCP `tools/call` (this example, and every `ailang-worldd call`), `coordinator:a2a` for an A2A `tasks/send`. Entries written before this tag existed say `coordinator:a2a` for both surfaces.

**2. The invocation record** (`world/invocation-record/v2`, the payload of `transitionRef`):

```json
{"invocationId": "mcp:ep1:93819a13…", "episodeId": "ep1", "skillId": "ailang-run",
 "transitionFn": "sha256:371fa1dc…", "interpreter": "sha256:1a67b014…", "semanticsEpoch": 1,
 "input": "sha256:d3527c92…", "output": "sha256:fa74ca0e…",
 "plan": "sha256:5d0d3fde…", "effects": ["sha256:f932f3aa…"]}
```

- `skillId`: which tool you called.
- `input`: your arguments object (here `{"path":"hello.ail"}`).
- `output`: exactly the result object you received, `world` block included.
- `plan`: the effect plan. Matches `world.plan` in your result.
- `effects`: the effect record refs. Match `world.effects[].record` in your result.

**3. The plan** (`world/effect-plan/v1`):

```json
{"effects": [{"cost": 1, "effect": "Ailang.Run", "id": "e1",
              "payload": {"path": "hello.ail"}, "scope": "worktree"}],
 "finish": false, "plan": "world/effect-plan/v1", "result": null}
```

A refused call has `"effects": []` and the refusal in `result`.

**4. The effect record** (`world/effect-record/v1`, provenance `host/broker`):

```json
{"effect": "Ailang.Run", "scope": "worktree", "cost": 1,
 "budgetBefore": 50, "budgetAfter": 49,
 "allowed": true, "failed": false, "denial": "",
 "requestRef": "sha256:d6624029…", "resultRef": "sha256:2285edce…"}
```

This is where to look when a call was denied (`allowed: false`, `denial: "denied:budget"` or another label, `resultRef` empty) or failed (`failed: true`). `budgetBefore` and `budgetAfter` tell you how much of that effect's budget is left.

**5. The handler result** (`resultRef`, `world/effect-result/v1`) is the raw handler output before the `world` block was added: for this run, `admitted`, `exit_code`, `decision`, `limit`, `stdout` (`"hello from pi\n"`) and `stderr`. `requestRef` (`world/effect-request/v1`) is the broker's own length-prefixed encoding of the request, not JSON.

## One command: `why`

`ailang-worldd why` does the walk above for you and checks every link: it recomputes each
content address from the bytes the daemon serves, recomputes the entry hash and the world ref,
and confirms the output is the world's `stateRoot`. Give it the result itself, piped from
`call --json-out`, or any other handle you have: an entry index, `head`, any `sha256:` ref
from your result (the output, `world.plan`, an effect record), or your `mcp:` or `a2a:` invocation id.

```bash
ailang-worldd call ailang-read --json '{"path":"data.txt"}' --json-out | ailang-worldd why -
```

```text
why: entry 2 (matched by output (sha256 of the result))
  ✓ world   sha256:08ab69ae…  revision 2, stateRoot sha256:4500f68f…, logHead sha256:b7aa36d6… (recomputed)
  ✓ entry   sha256:b7aa36d6…  #2 writtenBy coordinator:mcp prev sha256:ef2b3322… transitionFn sha256:076bd38e… interpreter sha256:1a67b014… (AILANG v0.41.0, the daemon's pin)
  ✓ record  sha256:08e8d4cd…  world/invocation-record/v2 invocation mcp:ep1:9ebcd16a… episode ep1 skill ailang-read
  ✓ input   sha256:db1faa2d…  {"path":"data.txt"}
  ✓ plan    sha256:a1f0feb6…  1 effect(s): e1 Workspace.Read@worktree cost 1
  ✓ effect  sha256:a06b7bac…  e1 allowed Workspace.Read@worktree cost 1 budget 19→18 request 71 B result 204 B
  ✓ output  sha256:4500f68f…  416 B = world.stateRoot; your result is these exact bytes
all 7 link(s) verified
```

A ✗ on any link exits 3: what the daemon served does not match its content address, so do not
trust that answer. A target other than an index or `head` is found by scanning back from the
head (`--scan`, default 500 entries); an older one is reported as not found, exit 1. `--json`
prints the chain for a program to read. See [the CLI reference](../reference/cli.md#why).

## Answering common questions

- **"What did I change in this file?"** List your recent entries with `/v1/log?from=<n>`, open each invocation record, and read the `input` of every `ailang-write` and `ailang-edit`. Write and edit payloads are recorded in full.
- **"Why did my run fail?"** Open the `output` of the `ailang-run` invocation and read `admitted`, `decision.error_kind`, `limit` and `stderr`. See [Results and errors](./results-and-errors.md#ailang-run-outcomes).
- **"How much budget do I have left?"** Read `budgetAfter` on your most recent effect record for that effect.
- **"Did my call land after the response was lost?"** Page the log from the last index you saw. If an entry with your tool and arguments is there, it committed; do not resend.

## Limits of the record

- World records **its own** actions, not the worktree's files. A file's current content is whatever is on disk; rebuild history from write and edit inputs.
- An `ailang-run` program's inner file and IO operations are governed by the AILANG policy and recorded only as the run's envelope (stdout, stderr, decision), not one by one.
- There is no daemon route that replays an invocation yet. Replay exists as library code exercised by tests.

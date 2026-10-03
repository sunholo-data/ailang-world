---
title: Your first world
sidebar_position: 2
description: Start the daemon, mint a session, commit a genesis world, read it back, and watch the guarantees refuse things.
---

# Your first world

This walkthrough follows `docs/QUICKSTART.md` §1–§5, updated for two changes since that
section was first run: the interpreter pin is now v0.41.0, and `POST /v1/commit` now requires a
session. It takes about five minutes.

You need the two programs built and `PIN` set, as in [Install](install.md):

```bash
unset AILANG_REGISTRY_API_KEY
export PIN=$HOME/.pinned-ailang/ailang
go build -o /tmp/ailang-worldd ./cmd/ailang-worldd
```

## 1. Start the daemon

```bash
/tmp/ailang-worldd serve --db /tmp/world-demo.db --ailang-bin $PIN &
```

It prints one line when the socket is bound:

```text
ailang-worldd listening on http://127.0.0.1:7644
```

`serve` is loopback-only: a non-loopback `--bind` is refused, with no override.
`--ailang-bin` archives the interpreter at startup. Its content hash becomes the replay pin that
every log entry carries. The first start also bootstraps the epoch registry for that
interpreter's release.

```bash
/tmp/ailang-worldd health
```

```json
{"status":"ok","daemon_version":"0.1.0","db_path":"/tmp/world-demo.db","interpreter_ref":"sha256:1a67b014…","interpreter_version":"AILANG v0.41.0\nCommit: 24ee108\n…"}
```

There is no world yet:

```bash
/tmp/ailang-worldd head
```

```text
ailang-worldd: /v1/head returned HTTP 404 NotFound: no world head has been selected yet
```

Every read is bounded. A store read that has not answered within 10 s returns
`503` with class `Timeout`. A genuine internal failure returns `500` with the fixed body
`internal store failure`; the cause goes to the daemon's stderr, one line per error. The
terminal running `serve` is where you read why.

## 2. Mint a session (attended)

Committing needs a session. Minting opens the store for writing, and the store allows one
writer process, so **stop the daemon first**:

```bash
kill %1
```

Then mint, at a real terminal:

```bash
/tmp/ailang-worldd session mint --db /tmp/world-demo.db --episode quickstart \
  --grant world.apply=world:10 --out /tmp/qs-session
```

It asks on the terminal:

```text
Confirm mint for episode quickstart (1 grant(s), expiry +3600s) with a session credential? [y/N]
```

Answer `y`. The token is written once to `/tmp/qs-session` at mode `0600`, and the
`credential_id` (for revocation) goes to stderr. The `world.apply=world:10` grant is not needed
to commit, which only needs a valid session. It is the grant the optional
[published transition](#7-optional-publish-and-call-a-pure-transition) below requires.

:::note Agents cannot do this step
`session mint` opens `/dev/tty` and refuses without one:
`refusing: no controlling terminal (…); minting a session credential requires one human act at a terminal`.
That is the fence working. See [Attended steps](../guides/attended-steps.md).
:::

## 3. Commit the genesis world

Restart the daemon:

```bash
/tmp/ailang-worldd serve --db /tmp/world-demo.db --ailang-bin $PIN &
```

A commit is JSON: `observedHead` (empty for genesis), content-addressed `objects` (payload in
base64; `hash` must be `sha256:<hex>` of the payload bytes, and the store verifies it),
`nextWorld`, and the log `entry` with its six-field header. Generate a valid one:

```bash
python3 - <<'EOF'
import json, hashlib, base64, os
def sha(b): return "sha256:" + hashlib.sha256(b).hexdigest()
payload = json.dumps({"goal": "hello, World"}).encode()
interp  = open(os.environ["PIN"], "rb").read()
eh = sha(b"genesis-entry-1")
c = {"observedHead": "",
     "objects": [{"hash": sha(payload), "interfaceHash": sha(b"iface-v1"),
                  "semanticId": "world/demo/genesis-goal", "provenance": "quickstart",
                  "payload": base64.b64encode(payload).decode()}],
     "nextWorld": {"ref": sha(b"world-1"), "revision": 0,
                   "stateRoot": sha(b"state-1"), "logHead": eh},
     "entry": {"header": {"entryIndex": 0, "semanticsEpoch": 1,
                          "transitionFn": sha(payload), "interpreter": sha(interp),
                          "prevEntryHash": sha(b"genesis"), "writtenBy": "quickstart"},
               "entryHash": eh, "transitionRef": sha(payload)}}
open("/tmp/genesis.json","w").write(json.dumps(c, indent=2))
EOF
/tmp/ailang-worldd commit --file /tmp/genesis.json --session "$(cat /tmp/qs-session)"
```

It returns the new head:

```json
{"selectedHead":"sha256:…"}
```

A commit has 3 s of store work. If it answers `503`, read the class:

- `Timeout`: the budget ended before the durable step. Nothing landed; resending the same file
  is safe.
- `CommitUncertain`: it ended after the durable step began. The commit either landed whole or
  not at all. **Do not resend yet.** Read `log get 0` and compare it with the `entry` you sent:
  equal in every field means it landed; a 404 or a different row means it did not. The current
  head is not evidence either way.

For an exact answer, add `"invocationId": "rest:<your-id>"` to the JSON. The daemon then records
the commit's intent first, and `curl -s http://127.0.0.1:7644/v1/receipts/rest:<your-id>`
answers `resolved` (with `resultRef` equal to your `nextWorld.ref`), `not-started` or
`indeterminate`.

## 4. Read everything back

```bash
/tmp/ailang-worldd head
/tmp/ailang-worldd log get 0
/tmp/ailang-worldd world get "$(/tmp/ailang-worldd head)"
```

Compare `log get 0`'s `header.interpreter` with `health`'s `interpreter_ref`: they are
identical. That is the replay pin, live. Object payloads read back with
`object get <hash> --payload` (base64). Read routes need no session.

## 5. See the guarantees refuse things

Send the same commit again:

```bash
/tmp/ailang-worldd commit --file /tmp/genesis.json --session "$(cat /tmp/qs-session)"
```

You get HTTP `409 HeadConflict` with both `observedHead` and `selectedHead`: the structured
conflict a caller re-plans from. A stale writer gets facts, not corruption.

Without the session:

```bash
/tmp/ailang-worldd commit --file /tmp/genesis.json
```

```text
ailang-worldd: /v1/commit returned HTTP 401 SessionAbsent: a session credential is required: no Authorization Bearer header was present
```

Start a second writer on the same file:

```bash
/tmp/ailang-worldd serve --db /tmp/world-demo.db --bind 127.0.0.1:7645
```

```text
ailang-worldd: daemon startup failed at store-open: another process already holds writer authority for this database (single-writer is enforced, not conventional): …
```

## 6. Stop

`SIGTERM` (Ctrl-C or `kill %1`) drains with a bound and releases the writer lock. A clean drain
prints nothing.

## 7. Optional: publish and call a pure transition

`docs/QUICKSTART.md` §6–§8 publish a one-line `echo` transition and call it over A2A and MCP.
Those sections are marked "attended — pending first verbatim run" in the runbook. The publish
step is covered in [Transitions manifest](../reference/transitions-manifest.md) and the calls in
[MCP and A2A](../reference/mcp-and-a2a.md). The session you minted above already holds the
`world.apply` grant the echo transition's `access` names.

Next: [Coding tools](coding-tools.md), the eight real transitions.

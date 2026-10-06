# Quickstart — run AILANG World's daemon in 5 minutes

*§1–§5 were first executed verbatim on 2026-07-28 against `dev` (attended demo, Mark +
coordinator); §9 was executed verbatim in the row-134 M6 smoke on 2026-10-03. On 2026-10-03
§1–§4 were brought up to date with the session-gated commit and the v0.41.0 interpreter pin
(queue row 137): the serve, health, unsessioned-commit (`401 SessionAbsent`) and mint-fence
steps were re-measured against a scratch daemon; the minted-session commit was not, because
mint needs a human at a terminal. Maintained under coding-standards **S7**: if this doc drifts
from the binary, that is a defect.*

> **Shortest path.** The software-engineering tools in [§9](#9-software-engineering-tools) are
> what an agent actually uses; `tools/attended/se_smoke.sh` runs that section one step at a
> time. Once a daemon is serving, `ailang-worldd tools list`, `call`, `why`, `log tail` and
> `provenance` replace hand-written curl (`ailang-worldd <verb> --help` for each).

## 0. Attended steps (the only commands a person must run)

An agent or a script can run everything else in this guide. These two steps each read a line
**you type** at your terminal, and they refuse to run when there is no terminal (an agent
harness, cron, CI). Run them yourself from the repo root, with the daemon stopped. The values
are the ones §9 uses.

**Publish the transitions** (once, and again whenever the manifest changes):

```bash
/tmp/world-publish transitions --store /tmp/se-world/world.db \
  --manifest packages/se-tools/transitions.json --ailang-bin $PIN
```

- `--store`: the world database, the same file the daemon serves with `serve --db`.
- `--manifest`: the transition descriptors to publish. Its paths are repo-relative, so run it
  from the repo root.
- `--ailang-bin`: the `.ail` interpreter pin, `$PIN` (AILANG v0.41.0).

At the prompt, type the line it shows exactly: `publish 9 transitions to /tmp/se-world/world.db`
(the manifest's descriptor count and your `--store`). On success it prints
`published transition registry revision N (head sha256:…)`. Running it again with nothing
changed prints `transition registry UNCHANGED at revision N`.

**Mint a session** (once per episode):

```bash
/tmp/ailang-worldd session new ep1 --db /tmp/se-world/world.db --workspace-root /tmp/se-ws \
  --preset se-tools --ttl 14400 --out /tmp/se-session
```

- `ep1`: the episode. Its worktree is `<workspace-root>/ep1`, and the name must match
  `^[a-z0-9][a-z0-9-]{0,63}$`.
- `--db`: the same store as the publish step.
- `--out`: where the token is written, once, at mode 0600. The file must not exist yet.

At `[y/N]`, type `y`. On success it prints `session for episode ep1: 9 grant(s), expires epoch …,
token written to /tmp/se-session`, plus `credential_id=…` on stderr (the hash `session revoke`
takes; the token itself is never printed).

**Two snags, both handled:**
- **IDE terminal panes work as they are.** Both commands open `/dev/tty` themselves and ask
  there, so you no longer need to add `< /dev/tty`. It is still accepted, and reads the same
  device. If you see `STOP fence=tty reason=no-controlling-terminal` or `refusing: no
  controlling terminal`, the command was not run from a terminal (for example, an agent ran it).
  Run it yourself in a terminal window.
- **`AILANG_REGISTRY_API_KEY`.** While it is set, `ailang-worldd` refuses every verb except
  `--help`, because every process World starts would inherit unrecallable publish authority
  (design Decision 4). Unset it in that shell (`unset AILANG_REGISTRY_API_KEY`), or prefix the one
  command: `env -u AILANG_REGISTRY_API_KEY /tmp/ailang-worldd session new …`.

Run `ailang-worldd doctor` in the shell you will use to check both. With the key set, doctor
refuses too, and names the variable.

## 1. Build and start

`PIN` is the `.ail` interpreter pin, **AILANG v0.41.0**. It runs every transition plan and its
content hash is in every log entry. Unset `AILANG_REGISTRY_API_KEY` first: with it in the
environment, the daemon refuses to start.

```bash
export PIN=$HOME/.pinned-ailang/ailang
$PIN --version
unset AILANG_REGISTRY_API_KEY
go build -o /tmp/ailang-worldd ./cmd/ailang-worldd
/tmp/ailang-worldd serve --db /tmp/world-demo.db --ailang-bin $PIN &
```

`serve` is loopback-only (a non-loopback `--bind` is refused, no override). `--ailang-bin`
archives and pins the interpreter at startup — its content hash becomes the D1 replay pin that
every log entry carries. The first start also bootstraps the epoch registry for that pin.

Every read below is bounded: a store read that has not answered within **10 s** is abandoned and
the route answers **HTTP 503, class `Timeout`**, naming the deadline — never a hang and never a
500. A genuine internal failure answers **HTTP 500 with the fixed body `internal store failure`**;
the verbatim cause (DSN path, driver detail) goes to the daemon's **stderr**, one line per error
carrying the route. A2A not-available refusals similarly write
`ailang-worldd: a2a refusal: tasks/send <invocation-id-or->: "<escaped cause>"`
to stderr (or the configured ErrorLog), with embedded newlines escaped. So the terminal running `serve` is where you read *why* a 500 happened — the
HTTP client is told *that* one happened, and nothing about the host.

```bash
/tmp/ailang-worldd health
```

Returns the daemon version, DB path, and the pinned `interpreter_ref` + version.

## 2. Commit the genesis transition

`POST /v1/commit` is **session-gated**: without a session it answers
`HTTP 401 SessionAbsent: a session credential is required`. A session is minted by an attended,
local operator verb that needs single-writer authority, so the order is: write the commit, stop
the daemon, mint, serve again, commit with `--session`.

A commit is JSON: `observedHead` (empty for genesis) + content-addressed `objects` (payload
base64; `hash` MUST be `sha256:<hex>` of the payload bytes — the store verifies) + `nextWorld`
+ the frozen 6-field log `entry` header. Generate a valid one:

```bash
python3 - <<'EOF'
import json, hashlib, base64, os
def sha(b): return "sha256:" + hashlib.sha256(b).hexdigest()
payload = json.dumps({"goal": "hello, World"}).encode()
interp  = open(os.environ["PIN"],"rb").read()
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
```

Stop the daemon, mint a session, and start the daemon again. **Mint is TTY-fenced**: it asks a
one-line y/N on the controlling terminal (it opens `/dev/tty` itself, so an IDE terminal pane
works as is) and, without one, refuses with `session mint: refusing: no controlling terminal`.
See [§0](#0-attended-steps-the-only-commands-a-person-must-run). The `world.apply` grant is the
one §6's published skill needs, so §6 reuses this session:

```bash
kill %1
/tmp/ailang-worldd session mint --db /tmp/world-demo.db --episode quickstart \
  --grant world.apply=world:10 --ttl 14400 --out /tmp/qs-session
/tmp/ailang-worldd serve --db /tmp/world-demo.db --ailang-bin $PIN &
/tmp/ailang-worldd commit --file /tmp/genesis.json --session /tmp/qs-session
```

Returns `{"selectedHead": "sha256:…"}` — the world now exists. `--session` takes the file (mode
0600) or the token itself; `$WORLD_SESSION` is the fallback.

A commit is bounded too: **3 s** of store work. If it answers **HTTP 503**, read the class.
`Timeout` means the budget ended *before* the durable step: nothing landed, and resending the
same file is safe. `CommitUncertain` means it ended *after* the durable step began: the commit
either landed whole or not at all. **Do not resend yet.** Reconcile by reading the log row at
your entry's index and comparing it with the `entry` you sent:

```bash
/tmp/ailang-worldd log get 0
```

Equal in every field means it landed. A 404, or a row that differs, means it did not, and you
may resend. The current head is **not** evidence either way: another commit may already have
landed on top of yours.

For an exact answer, give the commit an id: add `"invocationId": "rest:<your-id>"` to the JSON
(only the `rest:` namespace is accepted). The daemon then records the commit's intent first, and
its receipt names the outcome:

```bash
curl -s http://127.0.0.1:7644/v1/receipts/rest:<your-id>
```

`resolved` (with `resultRef` = your `nextWorld.ref`) means it landed; `not-started` or
`indeterminate` means it did not.

## 3. Read everything back

```bash
/tmp/ailang-worldd head
/tmp/ailang-worldd log get 0
/tmp/ailang-worldd world get "$(/tmp/ailang-worldd head)"
```

Check `log get 0`'s `header.interpreter` against `health`'s `interpreter_ref`: **identical** —
that is the replay pin, live. Object payloads read back with
`object get <hash> --payload` (base64; the store verified content-vs-hash on write).

## 4. See the guarantees refuse things

```bash
/tmp/ailang-worldd commit --file /tmp/genesis.json --session /tmp/qs-session
```

→ HTTP 409 `HeadConflict` with `observedHead`/`selectedHead` — the structured conflict a
caller re-plans from; stale writers get facts, not corruption. Without `--session` the same
command is refused earlier, `HTTP 401 SessionAbsent`.

```bash
/tmp/ailang-worldd serve --db /tmp/world-demo.db
```

→ refused at startup: *"another process already holds writer authority for this database
(single-writer is enforced, not conventional)"* — the ratified arm-A lock, fail-closed.

## 5. Stop

SIGTERM (Ctrl-C / `kill`) drains bounded and releases the writer lock.

## 6. Publish a transition skill to the A2A card *(attended — pending first verbatim run)*

The transition registry is written by an **attended, local** operator verb, never by the daemon
and never over the network. It needs single-writer authority, so the daemon must be **stopped**
(step 5). Capture the daemon's pinned interpreter first — every published descriptor pins it:

```bash
/tmp/ailang-worldd health   # before step 5: note "interpreter_ref"
```

A publishable source must be a **self-contained** module: publication runs the pinned
interpreter's `check` on it in an empty scratch root, so a module importing `world/*` is refused
with `LDR001` (that is why no landed `world/*.ail` is publishable yet).

```bash
cat > /tmp/echo.ail <<'EOF'
module quickstart/echo

export func main(input: string) -> string {
  "{\"echo\":${input}}"
}
EOF
cat > /tmp/transitions.json <<'EOF'
[{"id": "tools.echo", "title": "Echo", "description": "quickstart transition",
  "transitionFnFile": "/tmp/echo.ail",
  "inputSchema": {"type": "object"}, "outputSchema": {"type": "object"},
  "access": {"effect": "world.apply", "scope": "world", "cost": 1},
  "declaredEffects": []}]
EOF
go build -o /tmp/world-publish ./cmd/world-publish
/tmp/world-publish transitions --store /tmp/world-demo.db --manifest /tmp/transitions.json \
  --interpreter-ref <interpreter_ref from health>
```

`world-publish` refuses without a controlling terminal
(`STOP fence=tty reason=no-controlling-terminal`; see
[§0](#0-attended-steps-the-only-commands-a-person-must-run)). It asks at `/dev/tty` itself, so an
IDE terminal pane works as is. It asks you to type a line naming the local write — here `publish 1 transitions to /tmp/world-demo.db` (the
manifest's descriptor count and the `--store` you gave). Output: `semantics epoch 1 derived from
world/epoch-registry/v1 for interpreter release "…"` then `published transition registry revision
1 (head sha256:…)`. Running it again prints `transition registry UNCHANGED at revision 1` — an
identical republish writes nothing. `semanticsEpoch` is omitted on purpose: it is derived from the
epoch registry the daemon bootstrapped, never defaulted.

The session minted in §2 already holds the skill's capability (`world.apply`). Restart the
daemon and read the card:

```bash
/tmp/ailang-worldd serve --db /tmp/world-demo.db --ailang-bin $PIN &
curl -s -H "Authorization: Bearer $(cat /tmp/qs-session)" http://127.0.0.1:7644/.well-known/agent.json
```

The card lists `tools.echo` / `Echo`. A session minted without `world.apply` gets the same 200
with **zero** skills — the card is capability-filtered per session. What publication does and does
not establish: the source object exists and passes standalone `check` under the pinned
interpreter, and its epoch is one the registry nominates for that interpreter's release; it does
**not** establish that the source implements the `main(input: string) -> string` invocation
calling convention. Invocation reports an incompatibility if it does not.

### 7. Invoke a published transition

**Attended — pending first verbatim run.** With the daemon running and the session from §6,
send one JSON-object data part. The task ID is an idempotency key within the session.

```bash
curl -s -H "Authorization: Bearer $(cat /tmp/qs-session)" \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tasks/send","params":{"id":"quickstart-1","metadata":{"skill_id":"tools.echo"},"message":{"role":"user","parts":[{"type":"data","data":{"message":"hello"}}]}}}' \
  http://127.0.0.1:7644/a2a/
```

The result has task `id` `quickstart-1`, `status.state` `completed`, one JSON-object
`artifacts[0].parts[0].data`, and `metadata.invocation_id`, `metadata.world_ref`, and
`metadata.entry_index`. Repeat the same `curl` command with the same task ID to retrieve
the committed result without executing the transition again. Inspect its log entry, or walk
its whole provenance chain (entry, invocation record, input, plan, effects, output — every link
checked):

```bash
curl -s -H "Authorization: Bearer $(cat /tmp/qs-session)" \
  http://127.0.0.1:7644/v1/log/1
/tmp/ailang-worldd why 1
```

---
**Later:** the approval-inbox workbench (item 7). The echo is a self-contained demonstration source; library-backed `world/*.ail`
transitions still need a separately pinned library dependency before publication.

The session-scoped agent card returns an empty `skills` array only when the registry head is confirmed absent in both reads. A raced publication followed by a store or integrity failure returns `503 ProjectionUnavailable`; a read deadline returns `504 ProjectionDeadlineExceeded`. Retry the card read after resolving the failure. A2A admission uses its existing unavailable JSON-RPC error and does not dispatch the failed admission.

### 8. Use MCP tools

With the published echo and bearer session from §6, POST JSON to `/mcp/`.
Send both `Content-Type: application/json` and `Accept: application/json, text/event-stream`:

```bash
curl -s -H "Authorization: Bearer $(cat /tmp/qs-session)" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' http://127.0.0.1:7644/mcp/
curl -s -H "Authorization: Bearer $(cat /tmp/qs-session)" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"tools_decho","arguments":{"message":"hello"}}}' http://127.0.0.1:7644/mcp/
curl -s -H "Authorization: Bearer $(cat /tmp/qs-session)" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '[{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"tools_decho","arguments":{"message":"first"}}},{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"tools_decho","arguments":{"message":"second"}}}]' http://127.0.0.1:7644/mcp/
```

The list contains `tools_decho`, the reversible encoding of `tools.echo`: `_` becomes `_u`,
`.` becomes `_d`, and `/` becomes `_s`. MCP names longer than 64 bytes refuse the whole
surface; A2A IDs remain verbatim. Input schemas missing a top-level type gain `type:"object"`
with constraints preserved. Successful replies use the pinned upstream SSE framing.

Each call item receives fresh admission and a new random task ID, even when JSON-RPC IDs
repeat. Without an explicit version, or with `MCP-Protocol-Version: 2025-03-26`, batches run
sequentially. Versions `2025-06-18` and `2025-11-25` refuse batches. A host failure returns
one whole-request error; earlier items may already have committed. Retries create new tasks
and may repeat effects. MCP JSON-RPC IDs provide no idempotency guarantee, and an MCP result
carries **no** log entry index (unlike A2A's `metadata.entry_index`). Before retrying a lost
response, look at what committed: `ailang-worldd log tail` lists the latest entries with their
episode, skill and effect statuses, and `ailang-worldd why <index>` walks one. `ailang-worldd
call` runs that probe itself on a host failure: `a commit landed (entry N)` means do not retry,
run `why N`.

The production callback runner allows eight callbacks, each bounded at 20 seconds. Its slots
remain occupied until callbacks return; the separate subprocess limit is also eight and remains
held until child reap. The aggregate POST context is 20 seconds, including credential resolution
(3 seconds) and initial Tools admission (10 seconds). Each invocation re-admits under its own
10-second read bound. Prompt-body requests can return the timeout envelope within the frozen
30-second write window. Slow bodies and synchronous parsing can consume that window and
produce a connection failure. No transport deadline is extended.

The upstream SSE golden is generated with `WORLD_UPDATE_MCP_GOLDEN=1 go test
./host/projection -run '^TestMCPWireConformance$'` against module `v0.47.2`; normal runs compare
both the checked-in bytes and a newly constructed upstream handler. Executable mapping and
absence examples run with `go test ./host/projection ./host/transitionreg -run '^Example' -v`.

Host behavior changes are recorded in the [host changelog](HOST_CHANGELOG.md).

### 9. Software-engineering tools

**Attended — run verbatim in the row 134 M6 smoke, 2026-10-03** (evidence:
`design_docs/verification/world-attended-2026-10-03-row134-m6/`). This section serves the nine
`packages/se-tools` transitions (`ailang-read`, `ailang-write`, `ailang-edit`, `ailang-check`,
`ailang-run`, `builtins-search`, `examples-search`, `ailang-cli`, and row 140's `workspace-exec`)
over `/mcp/` and `/a2a/`. This section configures no exec profile, so `workspace-exec` answers
every call `{"ok":false,"refused":"no exec profile configured: …"}` and runs nothing; §10 turns it
on.
Every tool call runs its plan in the pinned interpreter, runs exactly one brokered effect with the
v0.52.1 tool binary inside the episode's worktree, and commits one log entry. (The M6 smoke ran
on v0.51.0; row 135 M0 moved the tool pin to v0.52.1 after re-proving the confinement matrix,
`design_docs/planned/w-software-engineering-domain.md` §13.) Its flags are bound
to `serve --help` and to the `session new` parser by `TestQuickstartSection9FlagsMatchTheCLI`.
The `-d` payloads below are run against a test daemon by
`TestSeToolsQuickstartPayloadsVerbatim`.

**One command per step:** `tools/attended/se_smoke.sh` wraps this runbook as subcommands
(`prepare`, `publish`, `mint`, `serve`, `check`, `pi`, `claude`, `log`, `stop`, `clean`) and absorbs
every snag below. It never runs the irreversible publish itself (AC30): its `publish` step points
you at the Publish block below, which you paste. Where it and this section disagree, this
section wins.

`ailang-worldd doctor --db /tmp/se-world/world.db --workspace-root /tmp/se-ws` names each of
the snags below that applies to your shell. Four facts measured on the first M6 run
(2026-10-03):
- **Unset `AILANG_REGISTRY_API_KEY`** in the shell that starts the daemon. With it in the
  environment, the daemon refuses to start (only `--help` still answers).
- The two attended steps (publish, session new) are TTY-fenced and are listed, copy-paste
  ready, in [§0](#0-attended-steps-the-only-commands-a-person-must-run). Both ask at
  `/dev/tty` themselves, so an IDE terminal pane works without `< /dev/tty` (this was a snag on
  the first M6 run, fixed on 2026-10-06). An agent harness has no controlling terminal: there
  `world-publish` stops with `STOP fence=tty reason=no-controlling-terminal` and `session new`
  with `refusing: no controlling terminal`.
- `POST /v1/commit` is session-gated, so the genesis commit comes **after** mint, against the
  running daemon, with `--session`.
- `/mcp/` answers as SSE: the JSON-RPC response is the `data:` line of an `event: message`.

The three end-to-end breaks measured in M5b (2026-10-02) are fixed:
- `session mint` grants now expire with the session, at mint time + `--ttl` (they were stored
  with expiry 0, so every CLI-minted session saw zero tools). Pinned by
  `TestSessionMintGrantsExpireWithTheSession`.
- `builtins-search` reads the text inventory. The JSON one is over `policy-tool`'s 64 KiB stdout
  cap. Matches are `{name, module, effect}`, with no signature or description.
- `examples-search` reads the corpus named by `serve --examples-dir`. That corpus is not built into
  the binary.

Pick two binaries. `PIN` is the `.ail` interpreter pin (AILANG v0.41.0) and runs every plan.
`TOOL` is the tool binary, which must be exactly AILANG v0.52.1; startup refuses any other
release. `ailang-worldd setup` installs both at the paths below, each verified against the
digests compiled into the CLI, and `ailang-worldd doctor` checks them (and the TTY fences, the
port, the store and the workspace) without changing anything. Build both CLIs from the repo
root:

```bash
export PIN=$HOME/.pinned-ailang/ailang
export TOOL=$HOME/.pinned-ailang-tools/v0.52.1/ailang
$PIN --version && $TOOL --version
go build -o /tmp/ailang-worldd ./cmd/ailang-worldd
go build -o /tmp/world-publish ./cmd/world-publish
```

Next, provision the worktree. Episode `ep1` maps to `<workspace-root>/ep1`. That path must be
a real directory, not a symlink, and the episode id must match `^[a-z0-9][a-z0-9-]{0,63}$`. The
store directory must lie **outside** the workspace root, or startup refuses.

```bash
mkdir -p /tmp/se-world /tmp/se-ws
git init -q /tmp/se-proj && git -C /tmp/se-proj commit -q --allow-empty -m init
git -C /tmp/se-proj worktree add --detach /tmp/se-ws/ep1
printf 'module hello\n\nimport std/io (println)\n\nexport func main() -> () ! {IO} {\n  println("hello from ailang-run")\n}\n' > /tmp/se-ws/ep1/hello.ail
```

Start the daemon once with the pin only, then stop it. Startup bootstraps the epoch registry
for the pin:

```bash
unset AILANG_REGISTRY_API_KEY
/tmp/ailang-worldd serve --db /tmp/se-world/world.db --ailang-bin $PIN &
until curl -sf http://127.0.0.1:7644/v1/health >/dev/null; do sleep 0.2; done
kill %1
```

Write the genesis commit now; it is sent after mint, because a tool call commits on top of a
selected head and `/v1/commit` needs the session:

```bash
python3 - <<'EOF'
import json, hashlib, base64, os
def sha(b): return "sha256:" + hashlib.sha256(b).hexdigest()
payload = json.dumps({"goal": "se-tools smoke"}).encode()
interp  = open(os.environ["PIN"], "rb").read()
eh = sha(b"se-genesis-entry")
c = {"observedHead": "",
     "objects": [{"hash": sha(payload), "interfaceHash": sha(b"iface-v1"),
                  "semanticId": "world/demo/genesis-goal", "provenance": "quickstart-se",
                  "payload": base64.b64encode(payload).decode()}],
     "nextWorld": {"ref": sha(b"se-world-1"), "revision": 0,
                   "stateRoot": sha(b"se-state-1"), "logHead": eh},
     "entry": {"header": {"entryIndex": 0, "semanticsEpoch": 1,
                          "transitionFn": sha(payload), "interpreter": sha(interp),
                          "prevEntryHash": sha(b"genesis"), "writtenBy": "quickstart-se"},
               "entryHash": eh, "transitionRef": sha(payload)}}
open("/tmp/se-genesis.json","w").write(json.dumps(c, indent=2))
EOF
```

**Publish (attended, TTY fence).** Run this from the **repo root**, because the manifest's
`transitionFnFile` paths are repo-relative. The daemon must be stopped, since publishing needs
single-writer authority. `world-publish` refuses without a controlling terminal and asks you, at
the terminal, to type `publish 9 transitions to /tmp/se-world/world.db` (the manifest's
descriptor count and the store). An agent cannot run this step. It is the first step of
[§0](#0-attended-steps-the-only-commands-a-person-must-run):

```bash
/tmp/world-publish transitions --store /tmp/se-world/world.db \
  --manifest packages/se-tools/transitions.json --ailang-bin $PIN
```

The output is `published transition registry revision 1 (head sha256:…)`. Re-running it prints
`UNCHANGED`.

**Session (attended, TTY fence).** `session new` mints a session for episode `ep1` holding the
nine effect grants, one per effect name the se-tools transitions declare (`--preset
se-tools`): `Workspace.Read`, `Workspace.Write` (shared by `write` and `edit`), `Ailang.Check`,
`Ailang.Run`, `Ailang.RunEnv` and `Ailang.RunNet` (the `ailang-run` calls whose `caps` hold
`Env` or `Net`, row 135), `Ailang.Discover` (shared by the two searches), `Ailang.CLI` and
`Workspace.Exec` (`workspace-exec`, row 140), each with scope `worktree` and a budget of 50
calls. It reuses the worktree made above (without it,
`--repo /tmp/se-proj` makes one), asks y/N on the terminal (type `y`), and writes the token once
to `--out` (mode 0600; the file must not exist yet). It is the second step of
[§0](#0-attended-steps-the-only-commands-a-person-must-run):

```bash
/tmp/ailang-worldd session new ep1 --db /tmp/se-world/world.db --workspace-root /tmp/se-ws \
  --preset se-tools --ttl 14400 --out /tmp/se-session
```

To keep runs at IO/FS (and Declassify), mint with explicit grants instead and leave out
`Ailang.RunEnv` and `Ailang.RunNet`: `--grant EFFECT=SCOPE:BUDGET`, repeatable, in place of
`--preset`. `session list --db /tmp/se-world/world.db` shows the credential (its hash, never the
token), and `session revoke --db /tmp/se-world/world.db <credential_id>` revokes it.

Serve with the workspace tools enabled. Both `--workspace-root` and `--tool-ailang-bin` are
required. With only one of them, the daemon logs `workspace tools disabled` and refuses every
tool call before any effect runs. This is rule R8, and over `/a2a/` its message is
`transition declares an effect this daemon has no handler for`.

`examples-search` needs an AILANG examples corpus that the operator names. v0.52.1 does embed a
corpus, but World does not fall back to it, and the tool runs with `HOME` set to a per-episode cache. `$TOOL examples download` fills
`~/.ailang/examples`, and serve uses that directory by default when it exists. `--examples-dir`
names a different corpus, which reaches the tool as `AILANG_EXAMPLES`. It must lie **outside**
the workspace root, or startup refuses. Without a corpus, `examples-search` answers
`{"ok":false,"refused":"no examples corpus configured: …"}`.

`ailang-run` takes `stdin`, `argv` and `caps` (row 135). Capabilities beyond IO and FS are off
unless the operator enables them. `--run-allow-caps` names the ones a run may request
(`Declassify`, `Env`, `Net`). `--run-net-allow` names each loopback `IP:PORT` a Net run may reach
(repeatable; a bare host, a name or a non-loopback address refuses startup), and
`--run-net-allow-http` allows plain http to them. Every other host, port and redirect hop is
refused, World's own `:7644` included. Start the HTTP mock a Net task talks to on the port you
name here (`7655` below); the grader's mock binds an ephemeral port, World's does not (R-135-9):

```bash
/tmp/ailang-worldd serve --db /tmp/se-world/world.db --ailang-bin $PIN \
  --workspace-root /tmp/se-ws --tool-ailang-bin $TOOL --examples-dir $HOME/.ailang/examples \
  --run-allow-caps Env,Net,Declassify --run-net-allow 127.0.0.1:7655 --run-net-allow-http &
```

Commit the genesis world through the running daemon, with the minted session:

```bash
/tmp/ailang-worldd --addr http://127.0.0.1:7644 commit --file /tmp/se-genesis.json --session "$(cat /tmp/se-session)"
```

List the tools, then make one call. Each response is SSE; the JSON-RPC object is on the
`data:` line:

```bash
curl -s -H "Authorization: Bearer $(cat /tmp/se-session)" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' http://127.0.0.1:7644/mcp/
curl -s -H "Authorization: Bearer $(cat /tmp/se-session)" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"ailang-read","arguments":{"path":"hello.ail"}}}' http://127.0.0.1:7644/mcp/
```

The list holds exactly the nine names. Each MCP name is its ID, because `-` needs no
escaping. A session sees only the tools whose effect it holds a grant for. The call result
carries the handler's fields (`ok`, `content`, `tool`, `policy_digest`) plus
`world: {effects:[{id, status, record}], plan}`. Resolve the record with
`curl -s -H "Authorization: Bearer $(cat /tmp/se-session)" http://127.0.0.1:7644/v1/objects/<record>`.
An argument outside a tool's schema is refused by its plan (`{"ok":false,"refused":"…"}`,
`effects: []`), and the refusal still commits.

The four core gate tasks `ailang-run` could not run before row 135 map onto its arguments like
this (the grader's `caps`, `stdin` and `cli_args`; write `benchmark/solution.ail` and any input
file with `ailang-write` first). Each prints the task's expected stdout:

| Task | `ailang-run` arguments | Effect | stdout |
|---|---|---|---|
| `pipeline` | `{"path":"benchmark/solution.ail","caps":["IO"],"stdin":"1\n2\n3\n4\n5\n"}` | `Ailang.Run` | `2 4 6 8 10`, one per line |
| `cli_args` | `{"path":"benchmark/solution.ail","caps":["IO","FS","Env"],"argv":["numbers.txt"]}` | `Ailang.RunEnv` | `15` |
| `api_call_json` | `{"path":"benchmark/solution.ail","caps":["Net","IO"]}`, the solution posting to `http://127.0.0.1:7655/` | `Ailang.RunNet` | `200` |
| `prompt_injection` | `{"path":"benchmark/solution.ail","caps":["Declassify","IO"]}` | `Ailang.Run` | `1` |

The result adds `policy: {digest, security_mode, caps, net_allow}`, the verified policy the run
executed under. A capability the operator did not enable is refused before anything runs (the
effect is recorded `failed`); one the session holds no grant for is `denied`.

**pi** loads the server through pi-mcp-adapter, first measured on 2.32.1 with `pi` 0.85.1 and
installed at 2.33.0 on 2026-10-03 (`~/.pi/agent/npm/node_modules/pi-mcp-adapter/package.json`).
`directTools` registers each tool as its own pi tool. The adapter names a direct tool
`<server>_<tool>` with `.` replaced by `_` (`formatToolName`, `toolPrefix: "server"`), so server
`world` gives `world_ailang-read` and so on:

```bash
cat > /tmp/se-pi-mcp.json <<'EOF'
{"mcpServers": {"world": {"url": "http://127.0.0.1:7644/mcp/", "auth": "bearer",
  "bearerTokenEnv": "WORLD_SESSION", "directTools": true, "toolPrefix": "server"}}}
EOF
export WORLD_SESSION=$(cat /tmp/se-session)
pi --mcp-config /tmp/se-pi-mcp.json
```

In that interactive session, run `/mcp tools` and confirm the nine `world_*` names. This first
session also writes the adapter's metadata cache (`~/.pi/agent/mcp-cache.json`). Direct tools
register from that cache, and until it exists they fall back to the proxy. Then run the smoke
with every built-in disabled:

```bash
pi --mcp-config /tmp/se-pi-mcp.json --no-builtin-tools \
  --tools world_ailang-read,world_ailang-write,world_ailang-edit,world_ailang-check,world_ailang-run,world_builtins-search,world_examples-search,world_ailang-cli \
  -p "Read hello.ail, change its message to 'hello from pi', check it, then run it and report its stdout."
```

**Claude Code** generates its own config. In a scratch directory, the `-s project` scope writes
`.mcp.json`, which holds the raw bearer token, so delete the directory after the smoke.
`--tools ""` disables the built-in tools, because `--allowedTools` alone leaves them in place:

```bash
mkdir -p /tmp/se-claude && cd /tmp/se-claude
claude mcp add --transport http world http://127.0.0.1:7644/mcp/ \
  --header "Authorization: Bearer $(cat /tmp/se-session)" -s project
claude -p "Read hello.ail, change its message to 'hello from claude', check it, then run it and report its stdout." \
  --tools "" --strict-mcp-config --mcp-config /tmp/se-claude/.mcp.json \
  --allowedTools mcp__world__ailang-read,mcp__world__ailang-write,mcp__world__ailang-edit,mcp__world__ailang-check,mcp__world__ailang-run,mcp__world__builtins-search,mcp__world__examples-search,mcp__world__ailang-cli
```

Confirm the `mcp__world__*` names in an interactive `claude --strict-mcp-config --mcp-config
/tmp/se-claude/.mcp.json` session with `/mcp`. Every call from either client is one log entry.
An MCP result carries no entry index, so read calls back with the developer CLI rather than
`/v1/log/<entry_index>`:

```bash
/tmp/ailang-worldd tools list --session /tmp/se-session
/tmp/ailang-worldd call ailang-read --session /tmp/se-session --arg path=hello.ail --json-out | /tmp/ailang-worldd why -
/tmp/ailang-worldd log tail
/tmp/ailang-worldd provenance --episode ep1
```

`call --json-out` prints the committed output bytes, whose hash is the output ref, so `why -`
finds the entry that committed it and walks the chain (exit 3 on any broken link). `log tail`
shows the latest entries with episode, skill and effect statuses. `provenance` prints the
`World-Provenance: store=… episode=ep1 entries=<from>-<to>` trailer for a PR or commit made
through World (D-WORLD-60; label such PRs `ailang-world`).

### 10. Build and test non-AILANG projects

**Written in row 140 M3 (2026-10-05); the attended run on three real projects is M4.**
`TestExecStartupRefusalTable` drives the refusal table below row by row.
`TestQuickstartSection10FlagsMatchTheCLI` binds the `serve` line to `serve --help`.
`TestQuickstartSection10ProfilesLoad` runs the profile script and loads its three profiles. The
design is `design_docs/planned/w-workspace-exec-toolchain-effect.md`.

`workspace-exec` (effect `Workspace.Exec`) runs **one command of a profile you write**. It runs
inside the episode's worktree, under the sandbox runtime `srt` (`@anthropic-ai/sandbox-runtime`).
The agent picks a command id and passes arguments. The profile fixes everything else: the
program, its leading arguments, the working directory, the environment and the timeout. Inside
the sandbox:
- writes reach only the worktree and the episode's own exec cache. `.git`, `.github`, `.claude`
  and the other World-owned names stay read-only.
- your `$HOME`, World's state dir and the workspace root (so every other episode) are
  unreadable, except the `read_roots` you list. System paths (`/usr`, `/opt`, `/etc`, `/tmp`)
  stay readable.
- there is no network. You install dependencies before the session.
- the child's environment is exactly World's set plus the profile's `env`, never yours.

A command runs for at most `timeout_ms`, which is at most 9000, inside the 10-second handler cap.
So a profile names targeted commands (one package, one test file, a type check), not full suites.

**Install and pin srt.** World accepts exactly srt 0.0.78 whose `dist/cli.js` hashes to the pin
below (`design_docs/verification/world-row140-m0/pin.json`). `serve` refuses any other. Keep srt,
your profiles and your cache seeds **outside** the workspace root and the state dir:

```bash
export EXEC=$HOME/.ailang-world-exec
mkdir -p $EXEC/srt $EXEC/profiles $EXEC/seeds
npm install --prefix $EXEC/srt --no-audit --no-fund @anthropic-ai/sandbox-runtime@0.0.78
shasum -a 256 $EXEC/srt/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js
# want 3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2c7cd40d96e
node --version   # v20.11 or later (srt's engines)
```

On Linux, srt needs bubblewrap, socat and ripgrep, and unprivileged user namespaces. Use
`sha256sum` for the digest there:

```bash
sudo apt-get install -y bubblewrap socat ripgrep
sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0
```

`--exec-node` defaults to the `node` on your `PATH`, resolved once at startup. It runs srt. It is
also the startup probe's client **inside** the sandbox, so it must be readable there. A Homebrew
or system node is. A node under `$HOME` (nvm) needs its directory in every profile's
`read_roots`.

**Profile keys.** A profile is one JSON object. An unknown key anywhere refuses startup.

| Key | Meaning |
|---|---|
| `profile` | always `"world/exec-profile/v1"` |
| `project` | the profile's name, `^[a-z0-9][a-z0-9-]{0,63}$`; `--exec-episode-project EP=PROJECT` selects it |
| `root` | the project directory inside the worktree (`"."` or a clean relative path), and the command's cwd |
| `path` | absolute directories. Every command's `argv[0]` is resolved on it once, at startup, and hashed |
| `env` | extra child variables (e.g. `GOFLAGS`). World owns `HOME`, `TMPDIR`, `PATH`, `LANG`, `LC_ALL`, `PYTHONPYCACHEPREFIX`, `npm_config_cache`, `CI`, `NO_COLOR`, `SANDBOX_RUNTIME`, every `*_PROXY`, `NO_PROXY`, `GIT_*` and the registry credentials. Naming one refuses startup |
| `read_roots` | absolute paths the child may read beyond the worktree and its cache (a module cache, a toolchain). Never `$HOME`, the state dir, the workspace root, an ancestor of them, a path inside the state dir or the root, or `/`, by path or by realpath |
| `caches` | `{"NAME": {"seed": "/abs/dir"}}`. The child gets `NAME=<episode cache>/NAME`, cloned from the read-only seed on first use; `{}` starts empty |
| `timeout_ms` | 1 to 9000 |
| `probe` | the startup probe's toolchain argv. It must exit 0 in an empty scratch worktree. Default: the first command's `argv[0] --version` |
| `commands.ID.argv` | the fixed program and leading arguments; `argv[0]` is a bare name on `path` or an absolute path |
| `commands.ID.flags` | the only flags the agent may pass, e.g. `{"-run": "regex"}`. Classes: `bool`, `regex`, `int`, `relpath`, `enum:[a,b]`. `{"-k": {"class": "regex", "form": "sep"}}` emits `-k`, `v` as two items (for argparse) instead of the default `"form": "eq"`, `-k=v` |
| `commands.ID.positional` | `relpath`, `pkgpattern` (`./...`, `a/b/...`) or `testfile` (with `suffixes`) |
| `commands.ID.suffixes` | the test-file suffixes, with `testfile` only |
| `commands.ID.max_args` | at most this many arguments (0 to 16, default 0); a flag and its value count once |
| `commands.ID.passthrough` | `"--"` lets the agent pass `--` before positionals |

**Three worked profiles.** Point the three variables at your checkouts, then write the profiles.
The Go toolchain and module cache are read from the compiler checkout. The Python interpreter
comes from the CLI's venv:

```bash
export AILANG_SRC=$HOME/dev/sunholo-data/ailang
export TWILIGHT_SRC=$HOME/dev/TwilightGame
export SUNHOLO_SRC=$HOME/dev/sunholo-platform
export GOTC=$(cd $AILANG_SRC && go env GOVERSION) GOMOD=$(cd $AILANG_SRC && go env GOMODCACHE)
export GODIR=$(dirname $(command -v go)) NODEDIR=$(dirname $(command -v node))
export PYBIN=$(cd $SUNHOLO_SRC/cli && .venv/bin/python -c 'import os, sys; print(os.path.realpath(sys.executable))')
python3 - <<'EOF'
import json, os
E = os.environ
V = "world/exec-profile/v1"
profiles = {
  "ailang-compiler": {"profile": V, "project": "ailang-compiler", "root": ".",
    "path": [E["GODIR"], "/usr/bin", "/bin"],
    "env": {"GOFLAGS": "-mod=readonly", "GOPROXY": "off", "GOTOOLCHAIN": E["GOTC"],
            "GOTELEMETRY": "off", "GOMODCACHE": E["GOMOD"]},
    "read_roots": [E["GOMOD"]],
    "caches": {"GOCACHE": {"seed": E["EXEC"] + "/seeds/ailang-compiler/go-build"}},
    "timeout_ms": 9000, "probe": ["go", "version"],
    "commands": {
      "test": {"argv": ["go", "test", "-count=1"],
               "flags": {"-run": "regex", "-v": "bool", "-short": "bool"},
               "positional": "pkgpattern", "max_args": 8},
      "vet": {"argv": ["go", "vet"], "positional": "pkgpattern", "max_args": 4}}},
  "twilightgame": {"profile": V, "project": "twilightgame", "root": ".",
    "path": [E["NODEDIR"], "/usr/bin", "/bin"],
    "timeout_ms": 9000, "probe": ["node", "--version"],
    "commands": {
      "typecheck": {"argv": ["node", "node_modules/typescript/bin/tsc", "--noEmit"]},
      "test-file": {"argv": ["node", "node_modules/vitest/vitest.mjs", "run"],
                    "positional": "testfile", "suffixes": [".test.ts", ".test.tsx"], "max_args": 4}}},
  "sunholo-cli": {"profile": V, "project": "sunholo-cli", "root": "cli",
    "path": ["/usr/bin", "/bin"],
    "read_roots": [os.path.dirname(os.path.dirname(E["PYBIN"]))],
    "timeout_ms": 9000, "probe": [E["PYBIN"], "--version"],
    "commands": {
      "test-file": {"argv": ["sh", "-c", "exec .venv/bin/python -m pytest -q -p no:cacheprovider \"$@\"", "pytest"],
                    "flags": {"-k": {"class": "regex", "form": "sep"}, "-x": "bool"},
                    "positional": "testfile", "suffixes": [".py"], "max_args": 6}}},
}
for name, p in profiles.items():
    with open(E["EXEC"] + "/profiles/" + name + ".json", "w") as f:
        json.dump(p, f, indent=1)
EOF
```

- **Go:** the module cache is a read root, and it also holds the switched toolchain (V17).
  `GOPROXY=off` keeps the run offline. `GOCACHE` is a per-episode clone of a seed you warm once.
- **TypeScript:** `tsc` and `vitest` run through `node` from the worktree's `node_modules`, so
  install them in each worktree (below).
- **Python:** `uv sync` makes the venv inside the worktree. Its interpreter is a symlink into
  uv's store under `$HOME`, so the store's version directory is the read root; World adds each
  root's realpath too. `-k` is `"form": "sep"` because pytest reads `-k=expr` as the value
  `=expr`. The `sh -c '… "$@"'` prefix belongs to the profile: the agent's arguments arrive as
  `"$@"`, and the shell never parses them. `-p no:cacheprovider` keeps `.pytest_cache` out of the
  worktree.

**Seed the caches.** Warm the Go build cache once, from your own checkout, with the profile's
environment. Each episode then starts from a clonefile copy (`cp -cR`; `cp -a --reflink=auto` on
Linux), and the seed is never written:

```bash
mkdir -p $EXEC/seeds/ailang-compiler/go-build
(cd $AILANG_SRC && GOCACHE=$EXEC/seeds/ailang-compiler/go-build GOFLAGS=-mod=readonly GOPROXY=off \
  GOTOOLCHAIN=$GOTC GOTELEMETRY=off go test -count=1 ./internal/lexer/ ./internal/parser/)
```

**Worktrees and dependencies.** Make one episode per project. Install each one's dependencies
before the session, because the sandbox has no network:

```bash
git -C $AILANG_SRC worktree add --detach /tmp/se-ws/go1
git -C $TWILIGHT_SRC worktree add --detach /tmp/se-ws/ts1 && (cd /tmp/se-ws/ts1 && npm ci --no-audit --no-fund)
git -C $SUNHOLO_SRC worktree add --detach /tmp/se-ws/py1 && (cd /tmp/se-ws/py1/cli && uv sync --frozen)
```

Mint one session per episode exactly as §9 does: `session new go1 …`, then `ts1` and `py1`,
each with `--preset se-tools`, which grants `Workspace.Exec`. Write each token to its own file
(`--out /tmp/se-go1-session`, and so on).

**Serve.** The four flags:
- `--exec-profile` is repeatable. With more than one profile, `--exec-episode-project EP=PROJECT`
  says which episode runs which. With one profile it is optional. An unmapped episode with
  several profiles is refused on each call.
- `--exec-sandbox` names srt's `node_modules`.
- `--exec-node` names the node (default: `node` on `PATH`).
- `--exec-max-output-bytes` is the per-stream kill (default 64 MiB).

```bash
/tmp/ailang-worldd serve --db /tmp/se-world/world.db --ailang-bin $PIN \
  --workspace-root /tmp/se-ws --tool-ailang-bin $TOOL \
  --exec-profile $EXEC/profiles/ailang-compiler.json --exec-profile $EXEC/profiles/twilightgame.json \
  --exec-profile $EXEC/profiles/sunholo-cli.json \
  --exec-episode-project go1=ailang-compiler --exec-episode-project ts1=twilightgame \
  --exec-episode-project py1=sunholo-cli \
  --exec-sandbox $EXEC/srt/node_modules --exec-node $NODEDIR/node --exec-max-output-bytes 67108864 &
```

**The startup probe.** Before `serve` listens, World runs seven arms for each profile, against
that profile's own sandbox settings. A scratch episode under the workspace root stands in for a
worktree, and World removes it afterwards, pass or fail. Every refusal must be a positive token
from the probe client, never a bare non-zero exit: `WORLD-PROBE-REFUSED EPERM` on macOS, and
`ENOENT` or `EROFS` on Linux, where a denied directory is masked. Each arm has an unsandboxed
control.

1. `arm1-write-inside`: a write inside the scratch worktree lands.
2. `arm2-write-outside`: a write to a sibling episode is refused, and the host file stays absent.
3. `arm3-read-fence`: reads of three decoys World plants are refused. The decoys are
   `$HOME/.ailang-worldd-exec-probe-decoy`, `<state>/exec-probe/decoy` and
   `<workspace root>/.exec-probe-sibling/decoy`.
4. `arm4-srt-default-write`: a write to `/tmp/claude` (srt's own default) is refused.
5. `arm5-exit-status`: a TERM self-kill reports 143.
6. `arm6-network`: World opens its own live `127.0.0.1` listener, and an unsandboxed connect is
   accepted first. A raw connect and a connect through srt's proxy must not reach it, and a port
   the child binds must be unreachable from the host.
7. `arm7-toolchain`: the profile's `probe` exits 0.

**Startup refusals.** `serve` exits 2 with `daemon startup failed at config: the exec
configuration is refused: …`, naming the flag and the field:

<!-- exec-refusals:begin -->
| Id | When | `serve` names |
|---|---|---|
| `no-workspace-tools` | an `--exec-*` flag without `--workspace-root` and `--tool-ailang-bin` | `need the workspace tools` |
| `sandbox-without-profile` | `--exec-sandbox` without `--exec-profile` | `--exec-sandbox: needs --exec-profile` |
| `node-without-profile` | `--exec-node` without `--exec-profile` | `--exec-node: needs --exec-profile` |
| `episode-without-profile` | `--exec-episode-project` without `--exec-profile` | `--exec-episode-project: needs --exec-profile` |
| `max-output-without-profile` | `--exec-max-output-bytes` without `--exec-profile` | `--exec-max-output-bytes: needs --exec-profile` |
| `no-sandbox` | `--exec-profile` without `--exec-sandbox` | `--exec-sandbox: is required` |
| `sandbox-not-srt` | `--exec-sandbox` names no srt install | `is neither a node_modules holding @anthropic-ai/sandbox-runtime` |
| `sandbox-other-version` | srt other than 0.0.78 | `only 0.0.78 is measured` |
| `sandbox-other-digest` | a `dist/cli.js` other than the pinned bytes | `is not the pinned 3c3092bd` |
| `node-absent` | the `--exec-node` file does not exist | `--exec-node` `does not resolve` |
| `node-too-old` | node below v20.11 | `is below 20.11` |
| `max-output-zero` | `--exec-max-output-bytes` below 1 | `must be at least 1` |
| `episode-not-pair` | an `--exec-episode-project` that is not `EP=PROJECT` | `is not EP=PROJECT` |
| `episode-bad-id` | an episode outside the episode grammar | `outside the episode grammar` |
| `episode-no-project` | a project that names no profile | `names no configured profile` |
| `episode-twice` | one episode mapped twice | `a second time` |
| `project-twice` | two profiles with one `project` | `is configured twice` |
| `profile-absent` | an `--exec-profile` file that does not exist | `--exec-profile` `does not exist` |
| `unknown-key` | a profile key World does not know | `unknown key "network"` |
| `root-escapes` | a `root` that is not a clean relative path | `root:` |
| `argv0-missing` | a command's `argv[0]` that is not on `path` | `does not resolve on path` |
| `env-world-owned` | an `env` naming a World-owned variable | `"HTTPS_PROXY" is World-owned` |
| `env-bad-name` | an `env` name outside `^[A-Za-z_][A-Za-z0-9_]*$` | `is not a variable name` |
| `timeout-cap` | a `timeout_ms` above 9000 | `timeout_ms:` |
| `probe-missing` | a `probe` program that is not on `path` | `probe:` |
| `readroot-home` | `read_roots: [$HOME]` | `ancestor of the operator HOME` |
| `readroot-home-symlink` | `read_roots: [<a symlink to $HOME>]` | `its realpath` `ancestor of the operator HOME` |
| `readroot-root-ancestor` | a read root above the workspace root | `is an ancestor of` |
| `readroot-slash` | `read_roots: ["/"]` | `is "/"` |
| `readroot-decoy` | a read root that is a probe decoy | `probe decoy` |
| `readroot-in-state` | a read root at or inside the state dir | `the state dir` |
| `sandbox-in-root` | `--exec-sandbox` inside the workspace root | `--exec-sandbox` `inside the workspace root` |
| `sandbox-in-state` | `--exec-sandbox` inside the state dir | `--exec-sandbox` `inside the state dir` |
| `sandbox-in-cache` | `--exec-sandbox` inside an exec cache | `--exec-sandbox` `inside an exec cache` |
| `node-in-root` | `--exec-node` inside the workspace root | `--exec-node` `inside the workspace root` |
| `node-in-state` | `--exec-node` inside the state dir | `--exec-node` `inside the state dir` |
| `node-in-cache` | `--exec-node` inside an exec cache | `--exec-node` `inside an exec cache` |
| `profile-in-root` | an `--exec-profile` file inside the workspace root | `--exec-profile` `inside the workspace root` |
| `profile-in-state` | an `--exec-profile` file inside the state dir | `--exec-profile` `inside the state dir` |
| `profile-in-cache` | an `--exec-profile` file inside an exec cache | `--exec-profile` `inside an exec cache` |
| `seed-in-root` | a cache seed inside the workspace root | `caches.GOCACHE.seed` `inside the workspace root` |
| `probe-passthrough` | an srt that runs commands unsandboxed | `arm2-write-outside` `arm3-read-fence` `arm4-srt-default-write` `arm6-network` |
| `probe-noop` | an srt that runs nothing | `arm1-write-inside` |
| `probe-zero` | an srt that reports a signal death as 0 | `arm5-exit-status` |
| `probe-toolchain` | a `probe` that exits non-zero | `arm7-toolchain` |
| `probe-no-client` | a probe client that cannot run inside the sandbox (a node under `$HOME` that no read root covers) | `arm2-write-outside` `arm3-read-fence` `arm6-network` |
| `probe-package-only` | an `--exec-sandbox` holding srt without its dependencies | `arm1-write-inside` `ERR_MODULE_NOT_FOUND` |
<!-- exec-refusals:end -->

**Call it.** The result is the command's own outcome, reported, never judged:
- `exit_code` (a signal death is 128 + n), `timed_out`, `limit` (`"output"` past the kill) and
  `duration_ms`;
- the head (8 KiB) and tail (56 KiB) of each stream, with its total bytes and sha256;
- the exact `argv`, and the `profile` and `sandbox` digests.

A command id the profile lacks, or an argument its grammar does not admit, gets
`{"ok":false,"refused":…}`, and nothing runs:

```bash
curl -s -H "Authorization: Bearer $(cat /tmp/se-go1-session)" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"workspace-exec","arguments":{"command":"test","args":["-run","TestLex","./internal/lexer/"]}}}' http://127.0.0.1:7644/mcp/
```

Each call is one log entry, and its record holds the full result. So `why` and replay return it
without running anything again. A change to a profile, node or srt shows as a different digest
in later records. `serve` reads a profile file once, so restart it after editing one.

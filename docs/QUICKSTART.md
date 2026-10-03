# Quickstart — run AILANG World's daemon in 5 minutes

*Every command below was executed verbatim on 2026-07-28 against `dev` (attended demo, Mark +
coordinator). Maintained under coding-standards **S7**: if this doc drifts from the binary,
that is a defect.*

## 1. Build and start

```bash
go build -o /tmp/ailang-worldd ./cmd/ailang-worldd
/tmp/ailang-worldd serve --db /tmp/world-demo.db --ailang-bin /tmp/ailang-v0300/ailang &
```

`serve` is loopback-only (a non-loopback `--bind` is refused, no override). `--ailang-bin`
archives and pins the interpreter at startup — its content hash becomes the D1 replay pin that
every log entry carries.

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

A commit is JSON: `observedHead` (empty for genesis) + content-addressed `objects` (payload
base64; `hash` MUST be `sha256:<hex>` of the payload bytes — the store verifies) + `nextWorld`
+ the frozen 6-field log `entry` header. Generate a valid one:

```bash
python3 - <<'EOF'
import json, hashlib, base64
def sha(b): return "sha256:" + hashlib.sha256(b).hexdigest()
payload = json.dumps({"goal": "hello, World"}).encode()
interp  = open("/tmp/ailang-v0300/ailang","rb").read()
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
/tmp/ailang-worldd commit --file /tmp/genesis.json
```

Returns `{"selectedHead": "sha256:…"}` — the world now exists.

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
/tmp/ailang-worldd commit --file /tmp/genesis.json
```

→ HTTP 409 `HeadConflict` with `observedHead`/`selectedHead` — the structured conflict a
caller re-plans from; stale writers get facts, not corruption.

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

Type the confirmation phrase when asked. Output: `semantics epoch 1 derived from
world/epoch-registry/v1 for interpreter release "…"` then `published transition registry revision
1 (head sha256:…)`. Running it again prints `transition registry UNCHANGED at revision 1` — an
identical republish writes nothing. `semanticsEpoch` is omitted on purpose: it is derived from the
epoch registry the daemon bootstrapped, never defaulted.

Mint a session that holds the skill's capability, restart the daemon, and read the card:

```bash
/tmp/ailang-worldd session mint --db /tmp/world-demo.db --episode quickstart \
  --grant world.apply=world:10 --out /tmp/qs-session
/tmp/ailang-worldd serve --db /tmp/world-demo.db --ailang-bin /tmp/ailang-v0300/ailang &
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
the committed result without executing the transition again. Inspect its log entry:

```bash
curl -s -H "Authorization: Bearer $(cat /tmp/qs-session)" \
  http://127.0.0.1:7644/v1/log/1
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
and may repeat effects. Inspect the daemon journal at `/v1/log/<entry_index>` before retrying
a lost response; MCP JSON-RPC IDs provide no idempotency guarantee.

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
`design_docs/verification/world-attended-2026-10-03-row134-m6/`). This section serves the eight
`packages/se-tools` transitions (`ailang-read`, `ailang-write`, `ailang-edit`, `ailang-check`,
`ailang-run`, `builtins-search`, `examples-search`, `ailang-cli`) over `/mcp/` and `/a2a/`.
Every tool call runs its plan in the pinned interpreter, runs exactly one brokered effect with the
v0.52.1 tool binary inside the episode's worktree, and commits one log entry. (The M6 smoke ran
on v0.51.0; row 135 M0 moved the tool pin to v0.52.1 after re-proving the confinement matrix,
`design_docs/planned/w-software-engineering-domain.md` §13.) Its flags are bound
to `serve --help` and to the `session mint` parser by `TestQuickstartSection9FlagsMatchTheCLI`.
The `-d` payloads below are run against a test daemon by
`TestSeToolsQuickstartPayloadsVerbatim`.

**One command per step:** `tools/attended/se_smoke.sh` wraps this runbook as subcommands
(`prepare`, `publish`, `mint`, `serve`, `check`, `pi`, `claude`, `log`, `stop`, `clean`) and absorbs
every snag below. It never runs the irreversible publish itself (AC30): its `publish` step points
you at the Publish block below, which you paste. Where it and this section disagree, this
section wins.

Four facts measured on the first M6 run (2026-10-03):
- **Unset `AILANG_REGISTRY_API_KEY`** in the shell that starts the daemon. With it in the
  environment, the daemon refuses to start.
- The two attended steps (publish, mint) are TTY-fenced. In an embedded terminal (an IDE pane,
  an agent harness) they fail `fence=tty reason=stdin-is-not-the-controlling-terminal`; append
  `< /dev/tty` to those two commands.
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
release. Build both CLIs from the repo root:

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
single-writer authority. `world-publish` refuses without a controlling terminal and asks for the
typed confirmation phrase. An agent cannot run this step:

```bash
/tmp/world-publish transitions --store /tmp/se-world/world.db \
  --manifest packages/se-tools/transitions.json --ailang-bin $PIN
```

The output is `published transition registry revision 1 (head sha256:…)`. Re-running it prints
`UNCHANGED`.

**Mint (attended, TTY fence)** a session holding the eight effect grants, one per effect name. A
grant is `EFFECT=SCOPE:BUDGET`, with scope `worktree` and a budget counted in calls. `write`
and `edit` share `Workspace.Write`, and the two searches share `Ailang.Discover`.
`Ailang.RunEnv` and `Ailang.RunNet` are the `ailang-run` calls whose `caps` hold `Env` or `Net`
(row 135); drop them to keep runs at IO/FS (and Declassify):

```bash
/tmp/ailang-worldd session mint --db /tmp/se-world/world.db --episode ep1 \
  --grant Workspace.Read=worktree:50 --grant Workspace.Write=worktree:50 \
  --grant Ailang.Check=worktree:50 --grant Ailang.Run=worktree:50 \
  --grant Ailang.RunEnv=worktree:50 --grant Ailang.RunNet=worktree:50 \
  --grant Ailang.Discover=worktree:50 --grant Ailang.CLI=worktree:50 \
  --ttl 14400 --out /tmp/se-session
```

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

The list holds exactly the eight names. Each MCP name is its ID, because `-` needs no
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

**pi** loads the server through pi-mcp-adapter, measured on 2.32.1 with `pi` 0.85.1.
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

In that interactive session, run `/mcp tools` and confirm the eight `world_*` names. This first
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
/tmp/se-claude/.mcp.json` session with `/mcp`. Every call from either client is one log entry:
read them back with `/v1/log/<entry_index>` as in §7.

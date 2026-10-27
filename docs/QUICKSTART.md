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

> **Section numbers.** §7–§11 are the `###` subsections of the transition and tool runbook
> (§9 the coding tools, §10 non-AILANG projects, §11 module roots and registry packages);
> "Watch the world live", "Open the live workbench" and "Prepare the release demo seed" are
> the top-level `##` sections numbered 10–12 after them. A bare "§10" or "section 11" in these
> docs and in `serve --help` means the subsection. The heading texts are bound by tests, so
> they keep their numbers.

## Run your agents through World

World launchers now tell Claude Code, Codex and pi their episode, worktree, session
lifetime, remaining grants and tool map. Claude adds a World status segment while
keeping your user status command; pi adds a footer. Refresh generated files with
`ailang-worldd hosts init --replace`. Inspect your own session with
`ailang-worldd --addr http://127.0.0.1:47330 whoami --json` using `WORLD_SESSION`
or `whoami --session <file>`; never put a token in argv. This calls the authenticated
`GET /v1/session/self` endpoint. It reports only the caller's session, with no token
or credential hash; exit 10 means daemon unavailable and 11 means session invalid.


With World serving and an approved episode session (see the attended steps below),
generate private launch settings once:

```bash
ailang-worldd hosts init --addr http://127.0.0.1:47330
export WORLD_SESSION="$(cat /path/to/private/episode.session)"
ailang-worldd --addr http://127.0.0.1:47330 doctor --db /path/to/world.db \
  --session /path/to/private/episode.session --host claude --quick
~/.ailang-world/hosts/bin/claude-world
# Or choose one of:
~/.ailang-world/hosts/bin/codex-world
~/.ailang-world/hosts/bin/pi-world # requires pi install npm:pi-mcp-adapter before hosts init
```

Use `--host codex` or `--host pi` for the corresponding doctor check. Drop
`--quick` to include the machine inventory and exec tools; `--profile <file.json>`
selects an exec profile explicitly. `--json` prints one object per check and a
failure count. Doctor returns daemon 10, session 11, host 12 or exec 13 for a
failed dependency, with a named cause and a command or exact action to fix it.
The quick pass reads daemon health, session expiry and launch digests. Its only
subprocess is the host's `--version`, with only PATH and HOME and a 5 s timeout.

Settings live under `~/.ailang-world/hosts` (override with `hosts init --out-dir`
and `doctor --hosts-dir`). Changed files refuse until `hosts init --replace`;
unchanged files are left alone. Directories and launchers are 0700, configs 0600.
Codex uses its own CODEX_HOME there and only World's MCP server; its auth.json is
a symlink to your existing ~/.codex/auth.json, never a copied credential. If that
file is absent, Codex asks you to log in. Each launcher refuses an unset
WORLD_SESSION with instructions to start an approved session.

In herdr, the plugin in `integrations/herdr/ailang-world/` (its README) does these steps from a popup: you type the mint's `y`, and it starts the confined pane, or a labelled plain fallback when World is down.

## 0. Attended steps (the only commands a person must run)

An agent or a script can run everything else in this guide. These two steps each read a line
**you type** at your terminal, and they refuse to run when there is no terminal (an agent
harness, cron, CI). Run them yourself from the repo root. The values are the ones §9 uses.
Publishing writes the world store, so the daemon must be stopped for it. Minting does **not**
need that (row 189): it writes the session credentials store, a separate file beside the
world store with its own lock, so it works while the daemon serves, and the daemon accepts the
new token from its next request.

One more command is attended, but this guide never needs it: `world-publish epoch nominate
--new-epoch`, which opens a new semantics epoch for a new AILANG release (you type a phrase
naming the epoch, the release and the store). Its `--candidate` form is headless. See
`world-publish` with no arguments for its usage.

**Publish the transitions** (once, and again whenever the manifest changes):

```bash
/tmp/world-publish transitions --store /tmp/se-world/world.db \
  --manifest packages/se-tools/transitions.json --ailang-bin $PIN
```

- `--store`: the world database, the same file the daemon serves with `serve --db`.
- `--manifest`: the transition descriptors to publish. Its paths are repo-relative, so run it
  from the repo root.
- `--ailang-bin`: the `.ail` interpreter pin, `$PIN` (AILANG v0.41.0).

At the prompt, type the line it shows exactly: `publish 10 transitions to /tmp/se-world/world.db`
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
- `--db`: the same store as the publish step. The credential itself goes to the credentials
  store beside it, `/tmp/se-world/credentials.db` (`--credentials-db` names another file; give
  `serve` the same one). Its first use with this build moves any credentials an older build
  kept inside `world.db` there, once.
- `--out`: where the token is written, once, at mode 0600. The file must not exist yet.

At `[y/N]`, type `y`. On success it prints `session for episode ep1: 9 grant(s), expires epoch …,
token written to /tmp/se-session`, plus `credential_id=…` on stderr (the hash `session revoke`
takes; the token itself is never printed).

`episode start` (row 189) is the other way to mint: it runs the same fence, y/N and mint for a
whole episode once a daemon is serving. See [Start a whole episode in one
act](#start-a-whole-episode-in-one-act) at the end of §1.

To run a coding agent on your own repo, follow [Use World on your repo](USE_WORLD_ON_YOUR_REPO.md),
which puts these steps in order.

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
`ailang-worldd: a2a refusal: tasks/send <invocation-id-or->: "<escaped cause>"`, and MCP refusals write
`ailang-worldd: mcp refusal: tools/call <invocation-id-or->: "<escaped cause>"` (or `tools/list -`)
to stderr (or the configured ErrorLog), with embedded newlines escaped. So the terminal running `serve` is where you read *why* a 500 happened — the
HTTP client is told *that* one happened, and nothing about the host. A transition whose **plan
phase** runs out of its own time budget (the caller still waiting) answers an A2A `-32603` that
names the phase and the budget (`transition plan phase exceeded its <budget> budget; nothing ran
or was committed; resend the same task id`) — nothing ran, so resending the same task id is safe;
the matching operator line reads `… plan phase exceeded its <budget> budget after <elapsed>`. MCP
clients see the released handler's fixed `host callback timed out` for the same event.

Publishing a transition (the `world-publish transitions` verb, or anything that calls the registry
publisher) also **builds that transition's compile-cache template** under
`<db>.artifacts/capsule-cache/<interpreter-digest>/<source-sha256>/`; every capsule run starts from a
copy of it, which is what keeps the plan phase fast under load. A publish whose interpreter passes
the check but writes no template is refused (retried three times first) rather than accepted into a
registry that would run cold forever. A run that finds no template (a store published before this
change, or a deleted cache) still works, but is never silent: it writes
`capsule: cold compile <interpreter-digest>/<source-sha256> (template missing)` to the daemon's
stderr (or the configured ErrorLog) once, and its success leaves the template behind for the next run.

At startup the daemon states, per transition in the registry head, whether its template is there:
`ailang-worldd: capsule template <id>: ready` or `… <id>: missing` (a stat each; it never builds
anything before it serves). Once the listener is up, a background pass rebuilds each `missing` one
(two at a time, each bounded by the 10 s check limit) and writes `capsule template <id>: rebuilt in
<duration>` or `… <id>: rebuild failed: <reason>`; a call that arrives first simply runs cold and is
labelled as above. The same pass prunes `capsule-cache/<digest>/` directories of interpreters no
transition pins any more (`capsule template gc: pruned <n>: <digests>`) and any `.tmp-*` build
directory older than twenty seconds.

```bash
/tmp/ailang-worldd health
```

Returns the daemon version, DB path, and the pinned `interpreter_ref` + version.

### Start a whole episode in one act

`episode start` (row 189) is attended like §0's two steps: run it yourself in a terminal. With a
daemon of the same build already serving, it validates everything, checks the daemon, runs the
same fence, y/N and mint as `session new`, and prints the agent's launch line. It never starts,
restarts or stops a daemon, and it never runs the agent.

```bash
/tmp/ailang-worldd episode start ep2 --repo ~/dev/my-project --profile my-project \
  --task ~/tasks/ep2.md --db /tmp/se-world/world.db --workspace-root /tmp/se-ws \
  --addr 127.0.0.1:7644
```

- `--profile my-project` binds the session to `~/.ailang-world-exec/profiles/my-project.json`
  (`--exec-profiles-dir` names another directory). The y/N names the project **and the sha256
  of the file's exact bytes**, and the credential records both. The daemon runs the episode's
  `workspace-exec` under that profile from its first call, **with no restart and no
  `--exec-episode-project` flag**. If the profile file changes later and the daemon is restarted
  on the new bytes, every exec of that episode refuses, naming both digests: re-mint it. Without
  `--profile` the session holds no `Workspace.Exec` grant.
- Before the terminal opens it reads `GET /v1/health` on `--addr` (default `127.0.0.1:47330`) and
  refuses a daemon of another build, one serving another tool binary than this build pins, one
  reading another credentials store, or one that has not loaded the profile's current bytes. Two
  builds are the same when their executables have the same sha256, or, failing that, when both
  carry the same clean (unmodified) VCS revision. A build in a git worktree carries no revision,
  so there only the same binary file matches. Each refusal names the fix, usually a rebuild and
  restart. With no daemon serving it refuses and prints the `serve` line to run.
- The token goes once to `~/.ailang/world/sessions/ep2.session` (0600, directory 0700;
  `--sessions-dir` names another), and for claude also to a private 0600
  `~/.ailang/world/sessions/ep2.mcp.json`. Delete the `.mcp.json` after the run. Either file
  existing already is refused before the terminal opens. The token is printed nowhere.
- The launch line is the bare `claude --tools '' --strict-mcp-config --mcp-config <file>
  --allowedTools <the World tools the grants cover>` (the `trial.py` shape), or with
  `--agent codex` the `WORLD_SESSION="$(cat <session file>)" codex exec …` form. Run it from an
  empty directory, never the worktree.

The daemon side, once: start `serve` with `--exec-sandbox <srt node_modules>` and no
`--exec-profile`, and it loads every `<project>.json` in `~/.ailang-world-exec/profiles/`
(`--exec-profiles-dir` names another; `*.proposed.json` is skipped until you rename it). A
profile there that fails to parse or probe is left out, named on stderr and in `/v1/health`
(`exec_profiles_excluded`), and the rest serve. `ailang-worldd episode start --help` has every flag.

**Propose a profile instead of writing one.** `profile init` (row 189) drafts the profile from
the repo's own files. It needs no daemon and no terminal, and it never activates anything:

```bash
/tmp/ailang-worldd profile init --repo ~/dev/my-project --project my-project
```

- **What it reads.** `go.mod` gives `go-test`, `go-vet` and `gofmt-l`. `pyproject.toml` or
  `uv.lock` gives `pytest-file` and `ruff-file`. `package.json` gives `vitest-file` or
  `jest-file`, whichever its `devDependencies` name. A `Makefile`'s targets are listed, but no
  command is proposed from them.
- **Where it looks.** It reads the repo root first. A repo with no such file at its root is read
  one directory down, and that directory becomes the profile's `root` (for example `backend`).
- **What each command looks like.** Each command has the argv, flags, positional and `max_args`
  of the measured profile's command of the same name.
- **Where the paths come from.** Every path comes from this machine: `go env` run in the root,
  the venv interpreter's link chain (or `uv python find`), and `node` on `PATH`.
- **What it prints and writes.** It prints one line per command saying which file triggered it
  and which measured command it copies. Only after the profile parser accepts the file does it
  write `~/.ailang-world-exec/profiles/my-project.proposed.json`, once. An existing file is
  refused, never overwritten (`--out` and `--profiles-dir` name others).
- **How to activate it.** Review the file, then run the `mv` it prints to rename it to
  `my-project.json`. `serve` skips `*.proposed.json`, and it validates and probes the renamed
  profile at startup.

Seed `GOCACHE` yourself if you want warm builds (§10). `ailang-worldd profile init --help` has
every flag.

**The Claude Code plugin.** The `ailang-world` plugin (row 189) teaches an agent inside an
episode how to work through World, and teaches you, in any repo, the commands above:
`/world-start <episode>` prints the `episode start` line for that repo (you run it and type
`y`), and `/world-watch` and `/world-why` print the watching commands. It runs none of them and
carries no MCP server or token. Install it once:

```bash
claude plugin marketplace add sunholo-data/ailang-world
claude plugin install ailang-world@ailang-world-marketplace
```

## 2. Commit the genesis transition

`POST /v1/commit` is **session-gated**: without a session it answers
`HTTP 401 SessionAbsent: a session credential is required`. A session is minted by an attended,
local operator verb. It writes the credentials store beside the world store
(`/tmp/credentials.db` here), not the world store the daemon holds, so the daemon keeps
serving: write the commit, mint, commit with `--session`.

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

Mint a session while the daemon serves. **Mint is TTY-fenced**: it asks a one-line y/N on the
controlling terminal (it opens `/dev/tty` itself, so an IDE terminal pane works as is) and,
without one, refuses with `session mint: refusing: no controlling terminal`. See
[§0](#0-attended-steps-the-only-commands-a-person-must-run). The `world.apply` grant is the
one §6's published skill needs, so §6 reuses this session:

```bash
/tmp/ailang-worldd session mint --db /tmp/world-demo.db --episode quickstart \
  --grant world.apply=world:10 --ttl 14400 --out /tmp/qs-session
/tmp/ailang-worldd commit --file /tmp/genesis.json --session /tmp/qs-session
```

The daemon reads the credentials store read-only for each request, so the commit right after
the mint is accepted with no restart, and a `session revoke` takes effect on the next request.

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
`design_docs/verification/world-attended-2026-10-03-row134-m6/`). This section serves the ten
`packages/se-tools` transitions (`ailang-read`, `ailang-write`, `ailang-edit`, `ailang-check`,
`ailang-run`, `builtins-search`, `examples-search`, `ailang-cli`, row 140's `workspace-exec` and
row 191a's `ailang-delete`)
over `/mcp/` and `/a2a/`. This section configures no exec profile, so `workspace-exec` answers
every call `{"ok":false,"refused":"no exec profile configured: …"}` and runs nothing; §10 turns it
on.
Every tool call runs its plan in the pinned interpreter, runs exactly one brokered effect with the
v0.54.15 tool binary inside the episode's worktree, and commits one log entry. (The M6 smoke ran
on v0.51.0; row 135 M0 moved the tool pin to v0.52.1 after re-proving the confinement matrix,
`design_docs/planned/w-software-engineering-domain.md` §13, and row 139 D7 moved it to v0.54.15,
the release carrying World's upstream fixes ailang#1872–#1875.) Its flags are bound
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
`TOOL` is the tool binary, which must be exactly AILANG v0.54.15; startup refuses any other
release. `ailang-worldd setup` installs both at the paths below, each verified against the
digests compiled into the CLI, and `ailang-worldd doctor` checks them (and the TTY fences, the
port, the store and the workspace) without changing anything. Build both CLIs from the repo
root:

```bash
export PIN=$HOME/.pinned-ailang/ailang
export TOOL=$HOME/.pinned-ailang-tools/v0.54.15/ailang
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
the terminal, to type `publish 10 transitions to /tmp/se-world/world.db` (the manifest's
descriptor count and the store). An agent cannot run this step. It is the first step of
[§0](#0-attended-steps-the-only-commands-a-person-must-run):

```bash
/tmp/world-publish transitions --store /tmp/se-world/world.db \
  --manifest packages/se-tools/transitions.json --ailang-bin $PIN
```

The output is `published transition registry revision 1 (head sha256:…)`. Re-running it prints
`UNCHANGED`.

**Session (attended, TTY fence).** `session new` mints a session for episode `ep1` holding the
ten effect grants, one per effect name the se-tools transitions declare (`--preset
se-tools`): `Workspace.Read`, `Workspace.Write` (shared by `write` and `edit`), `Ailang.Check`,
`Ailang.Run`, `Ailang.RunEnv` and `Ailang.RunNet` (the `ailang-run` calls whose `caps` hold
`Env` or `Net`, row 135), `Ailang.Discover` (shared by the two searches), `Ailang.CLI`,
`Workspace.Exec` (`workspace-exec`, row 140) and `Workspace.Delete` (`ailang-delete`, row 191a:
deleting is its own grant, never part of `Workspace.Write`), each with scope `worktree` and a
budget of 50 calls. It reuses the worktree made above (without it,
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
token), and `session revoke --db /tmp/se-world/world.db <credential_id>` revokes it. Both work
while the daemon serves: they read and write `/tmp/se-world/credentials.db`, never the world
store's writer lock, and a later episode is minted the same way without stopping this daemon.

Serve with the workspace tools enabled. Both `--workspace-root` and `--tool-ailang-bin` are
required. With only one of them, the daemon logs `workspace tools disabled` and refuses every
tool call before any effect runs. This is rule R8, and over `/a2a/` its message is
`transition declares an effect this daemon has no handler for`.

`examples-search` needs an AILANG examples corpus that the operator names. The tool binary may
embed a corpus (v0.52.1 did), but World does not fall back to it, and the tool runs with `HOME` set to a per-episode cache. `$TOOL examples download` fills
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

The list holds exactly the ten names. Each MCP name is its ID, because `-` needs no
escaping. A session sees only the tools whose effect it holds a grant for. The call result
carries the handler's fields (`ok`, `content`, `tool`, `policy_digest`) plus
`world: {effects:[{id, status, record}], plan}`. A large file can be read in pages: `ailang-read`
takes optional integer `offset` (1-based first line) and `limit` (line count, capped at 2000 by
the tool binary), and a ranged read also returns `total_lines`, `offset` and `limit` (row 190). Resolve the record with
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

In that interactive session, run `/mcp tools` and confirm the ten `world_*` names. This first
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
/tmp/ailang-worldd episode report ep1
/tmp/ailang-worldd --addr 127.0.0.1:7644 episode report ep1 --json
curl -s http://127.0.0.1:7644/v1/episodes/ep1/report
```

`episode report` (row 198) summarises the episode in one request to the daemon: total calls,
calls by tool, effect decisions by effect (`allowed`, `failed`, `denied:<why>`), budget per grant
(`minted`, `used`, `left`), plan refusals (the plan named no effect and answered `ok:false`) and
handler refusals (the effect was allowed and ran, but the handler answered `ok:false`, which the
effect record alone shows as allowed), with the first and last entries and their times. `--json`
prints the `GET /v1/episodes/{episode}/report` body unchanged; `--follow` re-prints it whenever the
store's head moves, until Ctrl-C. It reads only that episode's journal rows and their records, so
its cost does not grow with other episodes. Like every `/v1` read it is unauthenticated today
(residual R1, row 199). `--addr` takes a bare `host:port` as well as `http://…`, for every client
verb.

`call --json-out` prints the committed output bytes, whose hash is the output ref, so `why -`
finds the entry that committed it and walks the chain (exit 3 on any broken link). `log tail`
shows the latest entries with episode, skill and effect statuses. `provenance` prints the
`World-Provenance: store=… episode=ep1 entries=<from>-<to>` trailer for a PR or commit made
through World (D-WORLD-60; label such PRs `ailang-world`).

When `why` elides a long input, it prints the short top-level fields (`path`, `command`, `op`, …)
first, so you can see which file an `ailang-edit` touched. A toolchain call (`ailang-run`,
`ailang-check`, `ailang-cli`, `workspace-exec`) that changed the worktree carries a
`worktree_delta` in its effect's output (and in the result, except `ailang-check`'s, which is its
finish report): the files it created, modified or deleted (path, size, sha256;
at most 256, then `truncated` with a `total`), on the `git-status` basis in a git worktree.
`why` prints it as `worktree: N created, M modified, K deleted (basis git-status)` and the
first paths; `--json` holds the whole list. A call that changes nothing has no `worktree_delta`.

An operator harness that grades tasks can record each outcome next to the episode's effects with
`ailang-worldd outcome record --file <outcome.json|-> --session /tmp/se-session`: one
`world/task-outcome/v1` object (the task, the arm, the entry range, the grading rule and the
result), committed under the episode's session, so a provenance walk reaches how a task was
judged. Exit 1 means refused before sending, 2 the commit was refused, and 3 the outcome is
uncertain (reconcile with `object get` before any resend).

`session new` refuses an existing `<workspace-root>/<episode>` that is not a git worktree. With
`--allow-non-git` it adopts the directory as a **degraded** episode: no branch diff, provenance
trailer or reset to a base, and worktree deltas on the bounded-walk basis (`walk`, or
`walk-capped` with nothing listed past 20,000 files). `session list` and `/v1/health` show it.

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
| `protected` | optional list of worktree-relative directory forms and base-name patterns; see below. No default |
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
export TWILIGHT_SRC=<your TwilightGame checkout>      # M4 ran ~/dev/markedmondson1234/TwilightGame
export SUNHOLO_SRC=<your sunholo-platform checkout>    # M4 ran ~/dev/sunholo-data/sunholo-platform
export GOTC=$(cd $AILANG_SRC && go env GOVERSION) GOMOD=$(cd $AILANG_SRC && go env GOMODCACHE)
export GODIR=$(dirname $(command -v go)) NODEDIR=$(dirname $(command -v node))
export PYBIN=$(cd $SUNHOLO_SRC/cli && .venv/bin/python -c 'import os, sys; print(os.path.realpath(sys.executable))')
# the venv links to uv's version-ALIAS dir (cpython-3.12-…), which links to the real one (cpython-3.12.13-…);
# both must be read roots (V24, M4 finding F1)
export PYLINK=$(readlink $SUNHOLO_SRC/cli/.venv/bin/python)
python3 - <<'EOF'
import json, os
E = os.environ
V = "world/exec-profile/v1"
profiles = {
  "ailang-compiler": {"profile": V, "project": "ailang-compiler", "root": ".",
    "protected": ["*_test.go", "*.golden", "tests/**"],
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
    "protected": ["tests/**", "*.test.ts", "*.test.tsx", "vitest.config.ts"],
    "path": [E["NODEDIR"], "/usr/bin", "/bin"],
    "timeout_ms": 9000, "probe": ["node", "--version"],
    "commands": {
      "typecheck": {"argv": ["node", "node_modules/typescript/bin/tsc", "--noEmit"]},
      "test-file": {"argv": ["node", "node_modules/vitest/vitest.mjs", "run"],
                    "positional": "testfile", "suffixes": [".test.ts", ".test.tsx"], "max_args": 4}}},
  "sunholo-cli": {"profile": V, "project": "sunholo-cli", "root": "cli",
    "protected": ["cli/tests/**", "conftest.py", "pyproject.toml"],
    "path": ["/usr/bin", "/bin"],
    "read_roots": sorted({os.path.dirname(os.path.dirname(E["PYLINK"])), os.path.dirname(os.path.dirname(E["PYBIN"]))}),
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
  uv's store under `$HOME`: first to the version-alias directory (`cpython-3.12-…`), which links to
  the real one (`cpython-3.12.13-…`). Both are read roots. With only the real one, the sandbox
  fences the alias path and every call exits 126 (`Operation not permitted`; M4 finding F1). `-k` is `"form": "sep"` because pytest reads `-k=expr` as the value
  `=expr`. The `sh -c '… "$@"'` prefix belongs to the profile: the agent's arguments arrive as
  `"$@"`, and the shell never parses them. `-p no:cacheprovider` keeps `.pytest_cache` out of the
  worktree.

**Protect the tests.** Every pattern has one of two forms: `D/**` protects a clean
relative directory and everything beneath it; `G` has no `/` and matches a base name using
`*`, `?` or `[…]`. Base-name patterns match at any depth. Thus `pyproject.toml` and
`conftest.py` protect every such file. There is no single-file form.
Paths are relative to the worktree, not the profile's `root`: the CLI's `root` is `cli`,
so its tests need `cli/tests/**`. Using `tests/**` would protect `<worktree>/tests` instead.
World's write tools work inside the episode's module root (§11), so World re-bases the set
onto it for them: with module root `cli`, `cli/tests/**` becomes `tests/**`. A directory
outside the module root is dropped there, because those tools cannot reach it.
Patterns are NFC-normalized, compared by case fold, deduplicated and sorted by fold key.

<!-- exec-protected-patterns:begin -->
| Pattern | Load | Why |
|---|---|---|
| `tests/**` | accepted | directory form |
| `cli/tests/**` | accepted | worktree-relative directory form |
| `*_test.go` | accepted | base-name form, any depth |
| `*.test.tsx` | accepted | base-name form, any depth |
| `pyproject.toml` | accepted | every file with this base name |
| `**/x_test.go` | refused | use the base-name form `x_test.go` |
| `tests/**/x` | refused | `**` is only valid as a final `/**` |
| `src/*.ts` | refused | slash patterns must be directory forms |
| `*` | refused | whole tree |
| `**` | refused | whole tree |
| `./**` | refused | whole tree |
| `{a,b}` | refused | braces are unsupported |
| `a\b` | refused | backslashes are unsupported |
| `!tests` | refused | leading negation is unsupported |
| `.git/hooks/**` | refused | already denied by World |
| `.github/**` | refused | already denied by World |
| ` tests/**` | refused | leading whitespace |
| `tests/** ` | refused | trailing whitespace |
<!-- exec-protected-patterns:end -->

At most 64 entries are allowed, each 1–256 bytes, without NUL or control characters.
The directory part must have no glob, leading `/`, empty, `.` or `..` segment, and must
not name or lie inside a World deny name or contain a `.git` segment.

The exec sandbox enforces these forms as follows:

<!-- exec-protected-os:begin -->
| OS | `enforced` | Directory form | Base-name form |
|---|---|---|---|
| darwin | glob | refused | refused |
| linux | dirs-only | refused | detected only |
<!-- exec-protected-os:end -->

On Linux, base-name writes, including creation, are detected by the integrity check.
They are not prevented by srt. A changed protected path is reported with `intact_after: false`;
the call still returns the command's result, and later calls are not refused.
Protected bytes do not prove that every test ran or that project code kept test semantics.
There is no default: absent `protected` or `[]` means no protected set, no `protected` key
in the exec result, and `why` prints `protected: none declared`.

The result's `protected` block records the base commit, its protected content digest and
file count, the episode-start ignored-path snapshot, and both integrity checks:

<!-- exec-protected-record:begin -->
```json
{"protected": {
  "digest": "sha256:…", "patterns": ["*.test.ts", "tests/**"], "enforced": "glob",
  "base": "37d6641b308a3e8ffb4d5a2023ba4faa89a5254a", "base_digest": "sha256:…", "base_files": 247,
  "base_ignored_digest": "sha256:…", "base_ignored_files": 0,
  "intact_before": true, "intact_after": false,
  "changes": [{"path": "tests/zz_new.test.ts", "status": "untracked", "phase": "after"}],
  "truncated": false, "total": 1, "error": ""
}}
```
<!-- exec-protected-record:end -->

A failed or timed-out check records an error and a null verdict, never an intact verdict.
`why` prints one of these lines (CHANGED can say `before`, `after` or `before and after`):

```text
protected: intact (247 files @ 37d6641, digest sha256:…, enforced glob)
protected: CHANGED after the call — tests/zz_new.test.ts (untracked)
protected: not checked — <error>
protected: none declared
```

On the v0.54.15 tool pin, `ailang-write` and `ailang-edit` refuse symlink and hard-link
aliases upstream. They refuse any multiply-linked file, even unprotected ones.
World's own alias pre-check is not built. `ailang-run` FS aliases are unmeasured on this pin;
a write through one is detected at the next exec call. An episode can end without that call.

For a protected profile whose baseline was captured, the command gets `min(timeout_ms, time
left before the 10 s handler deadline − 3 751 ms)`, reserving time for the after-check. If the
baseline failed, both checks are skipped, the record carries the error, and nothing is
reserved. Unprotected profiles keep their own `timeout_ms`.
Install dependencies before the session: later ignored protected files report as changes.
Changes to pre-existing ignored protected files are not detected by the ignored-path snapshot.

**Seed the caches.** Warm the Go build cache once, from your own checkout, with the profile's
environment. Each episode then starts from a clonefile copy (`cp -cR`; `cp -a --reflink=auto` on
Linux), and the seed is never written:

```bash
mkdir -p $EXEC/seeds/ailang-compiler/go-build
(cd $AILANG_SRC && GOCACHE=$EXEC/seeds/ailang-compiler/go-build GOFLAGS=-mod=readonly GOPROXY=off \
  GOTOOLCHAIN=$GOTC GOTELEMETRY=off go test -count=1 ./internal/lexer/ ./internal/parser/)
```

**Warm each worktree once.** The first command in a freshly created worktree can overrun the
9-second cap from cold disk alone (M4 measured `go test` at 14 s, then 2.3 s, then 0.3 s, outside
the sandbox). The call reports `timed_out: true`, and a second call passes. Run each profile's
command once, unsandboxed, in each new worktree before the session.

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
  several profiles is refused on each call. (Row 189: a session minted with
  `session new`/`episode start --profile PROJECT` carries its own binding, which comes first and
  needs no flag and no restart; `--exec-episode-project` is for sessions minted without one.
  `--exec-profiles-dir` loads a whole directory of `<project>.json` instead, leaving out, rather
  than refusing startup for, a profile that fails; see §0.)
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

**The startup probe.** Before `serve` listens, World runs seven arms, plus an eighth when `protected` is declared, against
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
8. `arm8-protected-write` (protected profiles only): directory writes are refused. Base-name
   writes are refused on darwin; on Linux they succeed and must be detected by the integrity
   check. The arm checks host bytes and has an unsandboxed control. Linux refusals also accept
   `EBUSY` and `EXDEV`.

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
| `protected-double-star` | a protected pattern with a leading double star, such as `**/x_test.go` | `protected[0]` `**/x_test.go` |
| `protected-slash-glob` | a protected slash glob, such as `src/*.ts` | `protected[0]` `src/*.ts` |
| `protected-over-limit` | 65 protected entries (the limit is 64) | `protected[64]` `65 entries exceeds the limit of 64` |
| `protected-malformed-value` | a protected string instead of an array | `protected: must be an array of strings` |
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
in later records. `serve` reads a profile file once, so restart it after editing one. A session
minted with `--profile` is bound to the digest of the bytes confirmed at its y/N: after an edit
and a restart, its every exec refuses naming both digests until the episode is re-minted
(`session revoke --episode EP`, then `episode start`/`session new --profile` again).

### 11. Projects with a module root and registry packages

**Written in row 141 M3; the attended verbatim run is pending.**
`TestQuickstartSection11FlagsMatchTheCLI` binds the `serve` line below to `serve --help`. The
design is `design_docs/implemented/w-workspace-project-layouts.md`.

Bare imports between a project's modules (`import daneel_reading`) resolve against the AILANG
sandbox root, and registry packages (`import pkg/sunholo/oauth/token`) resolve under the tool's
`HOME`, which World points at an empty per-episode directory. Two flags fix those two gaps. Both
are off by default, and with neither set `serve` behaves exactly as in §9.

**Provision the worktree.** Make the episode as in §9 (the example project keeps its modules in
`tools/`, so its module root is `tools`):

```bash
export WS=/tmp/se-ws SNAP=$HOME/.ailang-world-pkgs
/tmp/ailang-worldd session new dev1 --db /tmp/se-world/world.db --workspace-root $WS \
  --repo /tmp/se-proj --preset se-tools --ttl 14400 --out /tmp/se-dev1-session
```

**The lock is an operator step.** The tool binary needs an `ailang.lock` in the module root before
it resolves a registry import, and since v0.54.15 it checks each registry dependency against the
lock's `content_hash` (a hand-written or stale lock is refused with `dependency … content
changed`; `$TOOL lock` writes the real hashes). The agent cannot create one: `lock` is not an admitted `ailang-cli`
op. A project that gitignores its lock has none in a fresh worktree. The lock is also the input
to the snapshot below, so run it first, as the operator, with the project's own `ailang.toml`:

```bash
(cd $WS/dev1/tools && $TOOL lock)
```

**Build the package snapshot.** `SNAP` is a read-only copy of exactly the registry packages the
lock names, laid out `<ns>/<name>/<version>/`. It must lie outside the workspace root and the
state directory. The loop stops at the first package your own cache (`~/.ailang/cache/registry`)
lacks:

```bash
LOCK=$WS/dev1/tools/ailang.lock
( set -euo pipefail
jq -r '.packages[] | select(.source == "registry") | "\(.name)/\(.version)"' "$LOCK" | while read -r p; do test -d "$HOME/.ailang/cache/registry/$p" || { echo "not in operator cache: $p" >&2; exit 1; }; mkdir -p "$SNAP/$(dirname "$p")" && cp -R "$HOME/.ailang/cache/registry/$p" "$SNAP/$p"; done
chmod -R a-w "$SNAP" )
```

The recipe runs in a subshell, so `set -euo pipefail` does not leak into your terminal and a
missing package stops the build without closing it.

- **More than one episode.** `--workspace-package-cache` is daemon-global, but each worktree has
  its own lock. Build `SNAP` as the union: run the `jq … | while …` loop once per episode's `LOCK`
  into the same `SNAP`, then `chmod -R a-w "$SNAP"` once at the end.
- **Refreshing.** Run `chmod -R u+w "$SNAP"`, rebuild, `chmod -R a-w "$SNAP"`, and **restart
  `serve`**: the snapshot's digest is computed once, at startup.

**Serve.** One extra flag per gap. `--workspace-module-root` is the default for every episode,
and `--workspace-episode-module-root EP=REL` (repeatable) overrides it for one episode, e.g.
`other1=.` for an episode whose modules sit at the worktree root:

```bash
/tmp/ailang-worldd serve --db /tmp/se-world/world.db --ailang-bin $PIN \
  --workspace-root $WS --tool-ailang-bin $TOOL \
  --workspace-module-root tools --workspace-episode-module-root dev1=tools \
  --workspace-package-cache $SNAP &
```

**What the agent sees.** The sandbox is `<workspace-root>/dev1/tools`, so every path in a
`read`, `write`, `ai-check` or `ailang-run` payload is relative to the module root
(`daneel_brief.ail`, not `tools/daneel_brief.ail`), and the agent cannot write anything outside
it (`.git`, `.github` and the rest of the worktree are out of reach). Each result carries
`policy_digest`, which changes with the sandbox, and, with a snapshot, `package_cache`: the
digest of the snapshot that `ai-check`, `test` and `run` resolved against. The effect record
therefore names the packages a result depended on. An exec profile's `root` (§10) is independent
of the module root: it only places a toolchain command's cwd.

**Operator lines.** Each is written once per refused call. The episode's AILANG tools
(`ailang-read`, `ailang-write`, `ailang-edit`, `ailang-delete`, `ailang-check`, `ailang-run`,
`builtins-search`, `examples-search`, `ailang-cli`) stay refused (every declared effect fails) until you fix the
cause. `workspace-exec` is not bound to that build: it needs only the episode worktree, not the
module root, the package cache or the policy. It is still refused when the worktree is missing or
unsafe, when no exec profile is configured (or none maps to the episode), and when building the
episode's exec handler fails. While an AILANG cause persists, each call in that episode,
`workspace-exec` included, first waits for the failed AILANG build to be retried (up to the 3 s
construction budget), so fix the cause rather than retrying:

| Line | Cause | What to do |
|---|---|---|
| `workspace tools unavailable for episode "dev1": module root "tools": …` | the module root is missing, not a directory, or reached through a symlink | create it as a real directory in the worktree (World never creates it) |
| `workspace package cache refused for episode "dev1": <dir> is not empty and was not provisioned by --workspace-package-cache; clear it manually` | the episode's `HOME` already holds a registry cache World did not link (for example a package fetched by `pkg_docs` before the flag was set) | inspect it and remove it yourself; World never deletes content |
| `workspace package cache does not cover episode "dev1": lock requires <ns>/<name>@<version>, absent from <dir>; …` | the worktree's `ailang.lock` needs a registry package the snapshot lacks (or the lock is unreadable, not a regular file or not JSON) | rebuild the snapshot from this episode's lock, then restart `serve` |
| `workspace package cache refused for episode "dev1": <link> is a symlink to <old>, not to --workspace-package-cache <dir>; remove it if the snapshot moved` | an earlier run linked the episode's `HOME` to a different snapshot path | `rm` that link yourself (World never deletes it) and retry |

`serve` itself refuses, at startup and naming the flag, a `--workspace-package-cache` that is not
a directory, lies inside the workspace root or the state directory, holds a symlink, a special
file or anything writable (the directory itself included), holds no `<ns>/<name>/<ver>/ailang.toml`
package, or when the daemon runs as root (mode bits do not stop root).

**Two things this does not change.** The admitted `pkg_docs` op still fetches from the registry
when it is given a package that is not cached and no snapshot is set; with the snapshot linked it
cannot store one (the tree is read-only). Closing that egress is a separate follow-up. And
`ailang.lock` stays writable by the agent: it can only make a package resolve if the snapshot
already holds it, since coverage is checked against the lock at episode construction.

## 10. Watch the world live

With the daemon from §1 serving and entries 0 and 1 committed (§2–§3), open a
read-only AG-UI 1.0 event stream. The complete request payload is:

```bash
curl -N -X POST -H 'Content-Type: application/json' -d '{"threadId":"t","runId":"r","messages":[]}' http://127.0.0.1:7644/agui/
```

The first three frames look like this (the CUSTOM value is an example committed
log entry; your hashes and writer will differ):

```text
data: {"type":"RUN_STARTED","threadId":"t","runId":"r","protocolVersion":"1.0"}

data: {"type":"STATE_SNAPSHOT","snapshot":{"schema":"world/agui-state/v1","lastIndex":-1,"logHead":null}}

data: {"type":"CUSTOM","name":"world.entry.committed","value":{"header":{"entryIndex":0,"semanticsEpoch":1,"transitionFn":"sha256:0000000000000000000000000000000000000000000000000000000000000064","interpreter":"sha256:0000000000000000000000000000000000000000000000000000000000000065","prevEntryHash":"sha256:0000000000000000000000000000000000000000000000000000000000000064","writtenBy":"fixture \u003c\u003e\u0026"},"entryHash":"sha256:0000000000000000000000000000000000000000000000000000000000000001","transitionRef":"sha256:00000000000000000000000000000000000000000000000000000000000000c8"}}

```

Each committed entry produces a `CUSTOM world.entry.committed` frame followed by
`STATE_DELTA` replacing `lastIndex` and `logHead`. Only that final delta frame has
an SSE `id`, acknowledging the complete entry. The stream crosses index gaps;
its CUSTOM value is the same JSON as `GET /v1/log/{index}`.

A run ends after **18 s** from handler entry with `RUN_FINISHED` and
`result.lastIndex`. To keep watching, **re-POST** with the last complete cursor.
For example, resume after entry 1:

```bash
curl -N -X POST -H 'Content-Type: application/json' -H 'Last-Event-ID: 1' -d '{"threadId":"t","runId":"r-resume","messages":[]}' http://127.0.0.1:7644/agui/
```

A stock AG-UI client sends its snapshot/delta state back in the standard `state`
field: `{"schema":"world/agui-state/v1","lastIndex":1,"logHead":"…"}`. The daemon
uses `lastIndex` and re-derives `logHead`; absent state or `state:{}` starts from
genesis. An explicit state cursor conflicting with `Last-Event-ID` gets 400, and
an unknown cursor gets 404 before any stream bytes.

After a severed connection, discard the unfinished trailing frame and resume
from the last complete delta's state or SSE id. An interrupted CUSTOM frame may
be replayed; the replacement deltas are idempotent. `threadId` and `runId` are
opaque strings echoed by the run; `messages`, `tools`, `context` and
`forwardedProps` are accepted and ignored. This route has the same read posture
as `GET /v1/log`, including ignoring any Authorization header. It adds no CORS
headers. At most 16 streams run concurrently; the 17th gets 503 `StreamLimit`
and `Retry-After: 1`. Send at most 64 KiB and finish the body within 2 s.

## 11. Open the live workbench

Open **http://127.0.0.1:7644/workbench** while your daemon is serving. The page
is read-only; checked links open stored entries and objects.

For static captures of recorded history, start `ailang-worldd serve` with
`--workbench-still` (default off). It suppresses only the `live.js` script;
the complete server-rendered page then never updates. Reload to refresh, as
the footer says: “Live updates: off. Reload to refresh.”

- **Live list:** the newest ten entries, newest first, plus the log position at
  render time. This tail updates even when the oldest-first timeline page is full.
- **Footer:** with JavaScript enabled, the last verified check, the cursor checked,
  and the next expected check (within 18 s); new entries and retry reasons are
  explicit. An empty log is a positive state, not evidence of a dead connection.
- **Graph:** the newest five entries on the home page, or the selected entry,
  linked to checked transitionFn, interpreter and transitionRef objects.
  Dashed UNAVAILABLE targets have no link. The graph carries no grades; on an
  object page it asks you to select an entry.
- **Decisions:** read-only approval requests with effect, scope, requester, cost
  and pending/approve/deny status (including who decided). It reads at most 40
  approval heads and keeps 20 summaries; truncation is labelled. No requests is
  stated explicitly, and malformed chains show UNAVAILABLE. Approvals are registry
  writes, so this pane refreshes on a log-triggered page update, with an as-of entry.

- **Object grade:** an object page shows `PROVEN` only when the daemon holds a
  proof-envelope validator (`serve --proof-key-file <32-byte file>
  --proof-required-identity NAME`, both needing `--ailang-bin`) and a stored
  `world/proof-report/v1` envelope, authenticated under that key, names the object
  as its subject. Otherwise the page says `GRADE UNAVAILABLE —` and one reason: no
  validator configured (`no canonical host projection…`), no stored envelopes
  (`no stored evidence…`), envelopes bound to other objects, an envelope that does
  not authenticate, a different compiler, a failed proof, a missing identity, or a
  malformed candidate. At most 100 candidates are checked, and a longer scan says
  so. Nothing in this release writes envelopes, so a configured daemon shows
  `no stored evidence` until a producer does.

**JavaScript disabled or blocked:** the complete server-rendered page still works;
its footer says “Live updates: off. Reload to refresh.” With JavaScript enabled,
new commits appear **without reload**: the client follows `/agui/`, then re-fetches
this page's own URL and replaces five server-rendered regions. It resumes bounded
runs at their terminal `lastIndex`, including across log index gaps.

Each visible tab uses one of the daemon's **16** global stream slots. A hidden
tab aborts its stream and resumes when shown. At the cap the footer says
“paused: stream limit”; retries use Retry-After and exponential backoff up to
30 s. Network/run errors are also shown as paused, and failed refreshes retain
the existing content. “page layout changed, reload required” requires a reload.

For a reproducible rehearsal, use a disposable store, not your live database.
Run these commands from the repository root in a terminal. Port 7644 must be
free. The session mint is attended: type `y` at its prompt (§0); the credential
stays in the temporary directory. This is the same genesis/commit wire recipe
as §2, with every file placed under `REHEARSAL` and both payloads constructed here.

```bash
export PIN=$HOME/.pinned-ailang/ailang
export REHEARSAL=$(mktemp -d)
unset AILANG_REGISTRY_API_KEY
go build -o "$REHEARSAL/ailang-worldd" ./cmd/ailang-worldd
"$REHEARSAL/ailang-worldd" serve --db "$REHEARSAL/world.db" --ailang-bin "$PIN" \
  >"$REHEARSAL/daemon.log" 2>&1 &
REHEARSAL_PID=$!
```

Wait for `listening` in `$REHEARSAL/daemon.log`. This first start creates the
store and pins the interpreter. Mint while it serves (the credential goes to
`$REHEARSAL/credentials.db`, which the daemon reads per request):

```bash
"$REHEARSAL/ailang-worldd" session mint --db "$REHEARSAL/world.db" \
  --episode workbench-rehearsal --grant world.apply=world:10 --ttl 3600 \
  --out "$REHEARSAL/session"
```

Then construct genesis entry 0
and the second commit, entry 1. Object hashes cover the exact base64 payload bytes;
`observedHead` and `prevEntryHash` in the second commit reference the first.

```bash
python3 - <<'PY'
import base64, hashlib, json, os
from pathlib import Path
def sha(b): return "sha256:" + hashlib.sha256(b).hexdigest()
root = Path(os.environ["REHEARSAL"])
interpreter = sha(Path(os.environ["PIN"]).read_bytes())
previous_world, previous_entry = "", sha(b"genesis")
for index, filename in [(0, "genesis.json"), (1, "second.json")]:
    payload = json.dumps({"goal": "live workbench", "step": index}).encode()
    obj, entry = sha(payload), sha(f"workbench-entry-{index}".encode())
    world = sha(f"workbench-world-{index}".encode())
    commit = {
        "observedHead": previous_world,
        "objects": [{"hash": obj, "interfaceHash": sha(b"iface-v1"),
                     "semanticId": "world/demo/workbench", "provenance": "quickstart",
                     "payload": base64.b64encode(payload).decode()}],
        "nextWorld": {"ref": world, "revision": index,
                      "stateRoot": sha(f"workbench-state-{index}".encode()), "logHead": entry},
        "entry": {"header": {"entryIndex": index, "semanticsEpoch": 1,
                             "transitionFn": obj, "interpreter": interpreter,
                             "prevEntryHash": previous_entry, "writtenBy": "quickstart"},
                  "entryHash": entry, "transitionRef": obj}}
    (root / filename).write_text(json.dumps(commit, indent=2))
    previous_world, previous_entry = world, entry
PY
"$REHEARSAL/ailang-worldd" commit --file "$REHEARSAL/genesis.json" \
  --session "$REHEARSAL/session"
curl --fail --silent --show-error http://127.0.0.1:7644/workbench >/dev/null
```

Open http://127.0.0.1:7644/workbench in a visible browser tab. Confirm the live
list says entry 0 and the footer shows a live check. Keep this tab open and run:

```bash
"$REHEARSAL/ailang-worldd" commit --file "$REHEARSAL/second.json" \
  --session "$REHEARSAL/session"
"$REHEARSAL/ailang-worldd" log get 1
```

The commit returns `selectedHead`; `log get 1` reports entry 1. Without reload,
the live list now puts entry 1 above entry 0, the graph gains entry 1, and the
footer reports the new entry. With JavaScript disabled, the page stays at entry 0
until you reload. This visual observation is the S7 browser rehearsal; curl alone
does not prove the live update.

Stop the disposable daemon when finished; retain its transcript for review.

```bash
kill "$REHEARSAL_PID"
wait "$REHEARSAL_PID"
```

The local automated request-trace drill is available to the controller on a
socket-capable host; it creates its own temporary store, profiles and logs:

```bash
WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" \
  go test -count=1 -timeout 90s -v ./host/daemon -run '^TestWorkbenchLiveChromeDrill$'
```

It observes the initial cursor POST, commits the next entry, observes the page
re-fetch and a later POST with the new cursor, then checks the JavaScript-off
control. Chrome is optional for ordinary tests/CI; a configured invalid path fails.

**Check the grade validator wiring** (disposable store, port 7645, from the repository root):

```bash
export PIN=$HOME/.pinned-ailang/ailang
GRADE=$(mktemp -d)
go build -o "$GRADE/ailang-worldd" ./cmd/ailang-worldd
head -c 32 /dev/urandom > "$GRADE/proof.key"
"$GRADE/ailang-worldd" serve --db "$GRADE/world.db" --ailang-bin "$PIN" --bind 127.0.0.1:7645 \
  --proof-key-file "$GRADE/proof.key" --proof-required-identity world/types.ail &
GRADE_PID=$!
for i in $(seq 1 60); do curl --silent --fail http://127.0.0.1:7645/v1/health >/dev/null && break; sleep 1; done
HEAD=$(curl --silent --fail http://127.0.0.1:7645/v1/registry/world/epoch-registry/v1 | python3 -c 'import json,sys; print(json.load(sys.stdin)["head"])')
curl --silent --fail "http://127.0.0.1:7645/workbench?object=$HEAD" | grep -o 'GRADE UNAVAILABLE — [^<]*'
kill $GRADE_PID; wait $GRADE_PID 2>/dev/null; rm -rf "$GRADE"
```

Expected: `GRADE UNAVAILABLE — no stored evidence: the store holds no object with semantic ID world/proof-report/v1`.

**Inspect a decision packet** (row 172). The invocation coordinator mints a `world/decision-packet/v1` object only when a skill's plan requests exactly one `Human.Approve` effect, the binder in use binds `Human.Approve`, and a live grant for that exact effect and scope authorizes it. The daemon binds `Human.Approve` for every episode and serializes it; it never decides a request (decision `D-WORLD-73`, ruled bind 2026-10-14). On a fresh store, until an invocation requests approval under a live grant, the listing below returns `{"items":[]}`; for example `--grant Human.Approve=worktree:1` on `session mint` (§0) authorizes that effect and scope. REST commits cannot store an object with this semantic ID either (`400 BadRequest`). A request never decides anything: decisions stay operator-only. When a packet exists, it rides in its invocation's own commit: `createdAt` is the invocation's logical time, `deadlineAt` is the authorizing grant's expiry, `policy` is `Cancel` and `escalationsRemaining` is `0`.

List packets and open the first one's object page while the daemon serves on 7644:

```bash
curl --silent --fail http://127.0.0.1:7644/v1/objects/by-semantic-id/world/decision-packet/v1
PACKET=$(curl --silent --fail http://127.0.0.1:7644/v1/objects/by-semantic-id/world/decision-packet/v1 | python3 -c 'import json,sys; items=json.load(sys.stdin)["items"]; print(items[0]["hash"] if items else "")')
[ -n "$PACKET" ] && curl --silent --fail "http://127.0.0.1:7644/workbench?object=$PACKET" | grep -o 'decision packet validated[^<]*\|not a validated decision packet[^<]*'
```

A produced packet reads `decision packet validated: committed with its invocation record at entry #N; …`. The page links the proposal (the invocation's `world/effect-plan/v1` plan), the approval request and entry N. It shows `approvals chain: pending` until an operator decides, and two evidence lines graded only by the object-grade validator above. Any other object labelled `world/decision-packet/v1` reads `not a validated decision packet: …` with the first failing check and no links. The checks are bytes, carrying commit, proposal, request, approvals-chain anchor, lifecycle fields and duplicate pair. The anchor check reads the newest 40 approval heads and keeps at most the newest 20 distinct requests; beyond that it reports "not established" rather than a verdict. The duplicate check reads the first 100 stored packets. The times are logical; this host emits no Timeout transition, so the page does not say whether the deadline has passed.

## 12. Prepare the release demo seed (incremental tooling)

`ailang-worldd demo seed --help` describes the shipped seed: it runs S1-S6 and
writes the public success manifest only after every check passes. It records domain
objects, read/test/check evidence, a proof and a pending approval packet. Success
returns exit 0. Inspect a partial store before any further action; never blindly
retry it. Release capture documentation follows in a later milestone.

Run from the repository root. Create a disposable worktree explicitly; the seed
never runs a hidden git operation. Keep the store, session, key and manifest outside
the worktree. The capsule/proof interpreter is v0.41.0; workspace tools use the
separate v0.54.15 tool pin. These pins are not interchangeable.

```bash
export PIN=$HOME/.pinned-ailang/ailang
export TOOLPIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang
export DEMO_DIR=$(mktemp -d /tmp/world-demo-release.XXXXXX)
export DEMO_WS=$DEMO_DIR/workspaces
export DEMO_BIN=$DEMO_DIR/ailang-worldd
mkdir -p "$DEMO_WS"
git worktree add --detach "$DEMO_WS/demo" HEAD
test -f "$DEMO_WS/demo/world/logepoch.ail"
go build -o "$DEMO_BIN" ./cmd/ailang-worldd
go build -o /tmp/world-publish ./cmd/world-publish
python3 - <<'PY'
import json, os
from pathlib import Path
root = Path(os.environ['DEMO_DIR'])
entries = json.loads(Path('packages/se-tools/transitions.json').read_text())
entries += json.loads(Path('packages/world-demo/transitions.json').read_text())
(root / 'transitions.json').write_text(json.dumps(entries, indent=2) + '\n')
key = root / 'proof.key'
fd = os.open(key, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
with os.fdopen(fd, 'wb') as f:
    f.write(os.urandom(32))
PY
```

With the daemon **stopped**, a person publishes the combined descriptor manifest
and types the exact confirmation phrase printed by the publisher. Then the person
mints one demo session and types `y` at the prompt (typed confirmation is required;
never pipe it or impersonate an operator):

```bash
/tmp/world-publish transitions --store "$DEMO_DIR/world.db" \
  --manifest "$DEMO_DIR/transitions.json" --ailang-bin "$PIN"
"$DEMO_BIN" session mint --db "$DEMO_DIR/world.db" --episode demo \
  --grant world.apply=world:10 --grant Workspace.Read=worktree:1 \
  --grant Ailang.CLI=worktree:2 --grant Human.Approve=worktree:1 \
  --ttl 3600 --out "$DEMO_DIR/session"
```

Start the sole writer in another terminal, exporting the same DEMO_DIR, DEMO_WS,
DEMO_BIN, PIN and TOOLPIN values there and using the same private key:

```bash
"$DEMO_BIN" serve --db "$DEMO_DIR/world.db" --ailang-bin "$PIN" \
  --workspace-root "$DEMO_WS" --tool-ailang-bin "$TOOLPIN" \
  --proof-key-file "$DEMO_DIR/proof.key" \
  --proof-required-identity sameRef --proof-required-identity servesEntry
```

Then, from the repository root in the prepared shell:

```bash
"$DEMO_BIN" demo seed --addr http://127.0.0.1:7644 \
  --db "$DEMO_DIR/world.db" --session "$DEMO_DIR/session" --ailang-bin "$PIN" \
  --proof-key-file "$DEMO_DIR/proof.key" --manifest "$DEMO_DIR/manifest.json"
```

The seed publishes a new manifest atomically after manifest and roster validation.
An unminted run refuses without opening `/dev/tty` or printing credentials. Preflight checks the absent manifest,
empty public world/log, daemon/database identity, interpreter hash/version and the
session's remaining live grants using an existing **read-only** authority view.
Bootstrap registries and session rows are allowed. Public writes always go through
session-authenticated REST commits with observed-head CAS; a changed head, uncertain
commit or refusal stops the command without automatic retries. The proof key
authenticates the proof stage and must also configure the daemon
with `--proof-key-file` and both `--proof-required-identity sameRef` and
`--proof-required-identity servesEntry`; the seed never exposes it.

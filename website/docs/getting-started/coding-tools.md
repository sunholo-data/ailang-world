---
title: Coding tools over MCP
sidebar_position: 3
description: Publish the eight se-tools transitions, mint a session, serve them with a worktree, and call them from curl, pi or Claude Code.
---

# Coding tools over MCP

The `packages/se-tools` package holds eight transitions that give an agent a complete AILANG
coding loop inside one git worktree, with every tool effect going through World's broker:

| Tool (ID = MCP name) | Arguments | Effect (scope `worktree`, cost 1) |
|---|---|---|
| `ailang-read` | `path` | `Workspace.Read` |
| `ailang-write` | `path`, `content` | `Workspace.Write` |
| `ailang-edit` | `path`, `old_text`, `new_text` | `Workspace.Write` |
| `ailang-check` | `path` | `Ailang.Check` |
| `ailang-run` | `path`, optional `args_json` | `Ailang.Run` |
| `builtins-search` | optional `query`, `module` | `Ailang.Discover` |
| `examples-search` | `query` | `Ailang.Discover` |
| `ailang-cli` | `op`, optional `path`, `module`, `query`, `package`, `flags` | `Ailang.CLI` |

Unknown argument keys are refused, never ignored. Every call, including a refused one, commits
one World log entry.

This page follows `docs/QUICKSTART.md` §9, which was run verbatim in the attended smoke on
2026-10-03: pi and Claude Code each read, edited, checked and ran `hello.ail` through World's
MCP tools only. Where this page and the runbook disagree, the runbook wins.

:::tip One command per step
`tools/attended/se_smoke.sh` wraps this whole runbook as subcommands — `prepare`, `publish`,
`mint`, `serve`, `check`, `pi`, `claude`, `log`, `stop`, `clean` — and absorbs every snag listed
below. It keeps its state under `~/.ailang/se-smoke` and serves on port 7644. Paste one line at
a time; interactive zsh does not treat `#` as a comment. It never runs the publish itself: its
`publish` step tells you to paste the Publish block yourself.
:::

## Before you start

Four facts measured on the first attended run:

- **Unset `AILANG_REGISTRY_API_KEY`** in the shell that starts the daemon, or it refuses to
  start.
- The two attended steps, publish and mint, need a controlling terminal. In an IDE pane or agent
  harness, `world-publish` fails with `STOP fence=tty reason=stdin-is-not-the-controlling-terminal`;
  append `< /dev/tty` to the publish command. `session mint` reads its confirmation from
  `/dev/tty` itself.
- `POST /v1/commit` is session-gated, so the genesis commit comes **after** mint, against the
  running daemon, with `--session`.
- `/mcp/` answers as SSE: the JSON-RPC response is on the `data:` line of an `event: message`.

## 1. Binaries

```bash
export PIN=$HOME/.pinned-ailang/ailang
export TOOL=$HOME/.pinned-ailang-tools/v0.51.0/ailang
$PIN --version && $TOOL --version
go build -o /tmp/ailang-worldd ./cmd/ailang-worldd
go build -o /tmp/world-publish ./cmd/world-publish
```

`PIN` (v0.41.0) runs every plan. `TOOL` must be exactly v0.51.0; startup refuses any other
release. See [Install](install.md) for why there are two.

## 2. Provision the worktree

Episode `ep1` maps to `<workspace-root>/ep1`. That path must be a real directory (not a
symlink) and the episode ID must match `^[a-z0-9][a-z0-9-]{0,63}$`. The store directory must lie
**outside** the workspace root, or startup refuses.

```bash
mkdir -p /tmp/se-world /tmp/se-ws
git init -q /tmp/se-proj && git -C /tmp/se-proj commit -q --allow-empty -m init
git -C /tmp/se-proj worktree add --detach /tmp/se-ws/ep1
printf 'module hello\n\nimport std/io (println)\n\nexport func main() -> () ! {IO} {\n  println("hello from ailang-run")\n}\n' > /tmp/se-ws/ep1/hello.ail
```

That writes this program into the worktree:

```ailang
module hello

import std/io (println)

export func main() -> () ! {IO} {
  println("hello from ailang-run")
}
```

## 3. Bootstrap the epoch registry

Start the daemon once with the pin only, then stop it:

```bash
unset AILANG_REGISTRY_API_KEY
/tmp/ailang-worldd serve --db /tmp/se-world/world.db --ailang-bin $PIN &
until curl -sf http://127.0.0.1:7644/v1/health >/dev/null; do sleep 0.2; done
kill %1
```

## 4. Write the genesis commit

It is sent later, after mint:

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

## 5. Publish the transitions (attended)

Run this from the **repo root**, because the manifest's `transitionFnFile` paths are
repo-relative. The daemon must be stopped: publishing needs single-writer authority.

```bash
/tmp/world-publish transitions --store /tmp/se-world/world.db \
  --manifest packages/se-tools/transitions.json --ailang-bin $PIN
```

It asks you to type a confirmation phrase at the terminal (see
[Attended steps](../guides/attended-steps.md#the-confirmation-phrase)), then prints
`published transition registry revision 1 (head sha256:…)`. Running it again prints
`UNCHANGED`: an identical republish writes nothing.

## 6. Mint a session (attended)

Six grants, one per effect name. A grant is `EFFECT=SCOPE:BUDGET`, with scope `worktree` and a
budget counted in calls. `ailang-write` and `ailang-edit` share `Workspace.Write`, and the two
searches share `Ailang.Discover`:

```bash
/tmp/ailang-worldd session mint --db /tmp/se-world/world.db --episode ep1 \
  --grant Workspace.Read=worktree:50 --grant Workspace.Write=worktree:50 \
  --grant Ailang.Check=worktree:50 --grant Ailang.Run=worktree:50 \
  --grant Ailang.Discover=worktree:50 --grant Ailang.CLI=worktree:50 \
  --ttl 14400 --out /tmp/se-session
```

A session sees only the tools whose effect it holds a grant for. Budgets persist per episode,
across restarts and across new sessions for the same episode.

## 7. Serve with the tools enabled

Both `--workspace-root` and `--tool-ailang-bin` are required. With only one, the daemon logs
`workspace tools disabled` and refuses every tool call before any effect runs.

`examples-search` needs an AILANG examples corpus. `$TOOL examples download` fills
`~/.ailang/examples`, which `serve` uses by default when it exists; `--examples-dir` names a
different one. It must lie outside the workspace root. Without a corpus, `examples-search`
answers `{"ok":false,"refused":"no examples corpus configured: …"}`.

```bash
/tmp/ailang-worldd serve --db /tmp/se-world/world.db --ailang-bin $PIN \
  --workspace-root /tmp/se-ws --tool-ailang-bin $TOOL --examples-dir $HOME/.ailang/examples &
```

Commit the genesis world through the running daemon, with the session:

```bash
/tmp/ailang-worldd --addr http://127.0.0.1:7644 commit --file /tmp/se-genesis.json --session "$(cat /tmp/se-session)"
```

## 8. Call the tools with curl

```bash
curl -s -H "Authorization: Bearer $(cat /tmp/se-session)" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' http://127.0.0.1:7644/mcp/
curl -s -H "Authorization: Bearer $(cat /tmp/se-session)" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"ailang-read","arguments":{"path":"hello.ail"}}}' http://127.0.0.1:7644/mcp/
```

The list holds exactly the eight names. The call answers as SSE. This is the `ailang-read`
result from the attended run (`structuredContent`, abbreviated):

```text
event: message
data: {"jsonrpc":"2.0","id":2,"result":{"content":[{"type":"text","text":"…"}],"structuredContent":{"content":"module hello\n…","ok":true,"policy_digest":"87d52ed5…","tool":"sha256:e55ff71c…","world":{"effects":[{"id":"e1","record":"sha256:76404771…","status":"ok"}],"plan":"sha256:ce4340a1…"}}}}
```

The result carries the handler's fields (`ok`, `content`, `tool`, `policy_digest`) plus
`world: {effects: [{id, status, record}], plan}`. Resolve the effect record with:

```bash
curl -s "http://127.0.0.1:7644/v1/objects/<record>?payload=true"
```

An argument outside a tool's schema is refused by its plan (`{"ok":false,"refused":"…"}` with
`effects: []`), and the refusal still commits.

## 9. Connect pi

pi loads the server through pi-mcp-adapter (measured on adapter 2.32.1 with pi 0.85.1).
`directTools` registers each tool as its own pi tool named `<server>_<tool>`, so server `world`
gives `world_ailang-read` and so on:

```bash
cat > /tmp/se-pi-mcp.json <<'EOF'
{"mcpServers": {"world": {"url": "http://127.0.0.1:7644/mcp/", "auth": "bearer",
  "bearerTokenEnv": "WORLD_SESSION", "directTools": true, "toolPrefix": "server"}}}
EOF
export WORLD_SESSION=$(cat /tmp/se-session)
pi --mcp-config /tmp/se-pi-mcp.json
```

Run `/mcp tools` in that session and confirm the eight `world_*` names. This first session also
writes the adapter's cache (`~/.pi/agent/mcp-cache.json`); direct tools register from it. Then
run with every built-in tool disabled:

```bash
pi --mcp-config /tmp/se-pi-mcp.json --no-builtin-tools \
  --tools world_ailang-read,world_ailang-write,world_ailang-edit,world_ailang-check,world_ailang-run,world_builtins-search,world_examples-search,world_ailang-cli \
  -p "Read hello.ail, change its message to 'hello from pi', check it, then run it and report its stdout."
```

## 10. Connect Claude Code

Claude Code generates its own config. In a scratch directory, `-s project` writes `.mcp.json`,
which holds the raw bearer token, so delete the directory afterwards. `--tools ""` disables the
built-in tools; `--allowedTools` alone leaves them in place.

```bash
mkdir -p /tmp/se-claude && cd /tmp/se-claude
claude mcp add --transport http world http://127.0.0.1:7644/mcp/ \
  --header "Authorization: Bearer $(cat /tmp/se-session)" -s project
claude -p "Read hello.ail, change its message to 'hello from claude', check it, then run it and report its stdout." \
  --tools "" --strict-mcp-config --mcp-config /tmp/se-claude/.mcp.json \
  --allowedTools mcp__world__ailang-read,mcp__world__ailang-write,mcp__world__ailang-edit,mcp__world__ailang-check,mcp__world__ailang-run,mcp__world__builtins-search,mcp__world__examples-search,mcp__world__ailang-cli
```

Both `pi` and `claude` read a non-TTY stdin as extra prompt input and wait for EOF; when you run
them from a script, give them `</dev/null` (the smoke script does).

## 11. Read what happened

Every call from any client is one log entry. In the attended run the store held 10 entries:
genesis plus nine tool calls.

```bash
/tmp/ailang-worldd log range --from 0
```

Tool-call entries have `writtenBy: "coordinator:a2a"`. To trace a call back from its result, see
[Provenance walks](../guides/provenance-walks.md).

## What the tools can and cannot do

The tools are confined by AILANG's own policy layer, not by new Go code: paths stay inside the
worktree, `.git` and a deny list of control directories are read-only, there is no shell, and
`ailang-run` has an 8 s supervisor timeout. The full list of measured facts is in
[Tool confinement](../security/tool-confinement.md), and the known gaps are in
[Residuals](../security/residuals.md).

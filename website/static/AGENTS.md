# AGENTS.md: operating guide for AI agents in AILANG World

You are a coding agent connected to an AILANG World daemon. This file is everything you need in one fetch. Longer pages with more examples live under /docs/agents/. The AILANG language reference for LLMs is at https://ailang.sunholo.com/llms.txt.

Status: pre-1.0. The eight tools below are live and were exercised end to end by pi and Claude Code on 2026-10-03. The comparison against agents with native tools has not been run yet.

## 1. What you have and what you do not

You have:
- one git worktree for your episode; every path is relative to its root;
- eight AILANG tools (section 3), only those your session holds grants for;
- a recorded history: every call commits one content-addressed log entry.

You do not have:
- a shell or any command execution beyond `ailang-run`;
- network access (programs may use the IO and FS capabilities only);
- write access outside the worktree, to `.git/`, or to the operator's deny list (`.github/`, `.pi/`, `.claude/`, `.ailang/`, `.gitmodules`, `.gitattributes`, in every case variant);
- git commits or new worktrees (operator actions).

## 2. Connect

Your operator gives you an address (default http://127.0.0.1:7644) and a session credential (64 hex characters). Never write the credential into a file you commit or a message you print.

MCP (streamable HTTP):
- endpoint `POST /mcp/`
- headers `Authorization: Bearer <session>`, `Content-Type: application/json`, `Accept: application/json, text/event-stream`
- replies are SSE: the JSON-RPC response is the `data:` line of `event: message`
- tool names are the tool ids (`ailang-read`, ...); results come as a `text` item and as `structuredContent`

```bash
curl -s -H "Authorization: Bearer $WORLD_SESSION" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-read","arguments":{"path":"hello.ail"}}}' \
  http://127.0.0.1:7644/mcp/
```

pi (pi-mcp-adapter, `directTools`): tools are named `world_<id>`.

```json
{"mcpServers": {"world": {"url": "http://127.0.0.1:7644/mcp/", "auth": "bearer",
  "bearerTokenEnv": "WORLD_SESSION", "directTools": true, "toolPrefix": "server"}}}
```

```bash
pi --mcp-config world-mcp.json --no-builtin-tools \
  --tools world_ailang-read,world_ailang-write,world_ailang-edit,world_ailang-check,world_ailang-run,world_builtins-search,world_examples-search,world_ailang-cli \
  -p "<task>" </dev/null
```

Claude Code: tools are named `mcp__world__<id>`. Run `claude mcp add` in a scratch directory (the written `.mcp.json` holds the token; delete it afterwards).

```bash
claude mcp add --transport http world http://127.0.0.1:7644/mcp/ \
  --header "Authorization: Bearer $WORLD_SESSION" -s project
claude -p "<task>" --tools "" --strict-mcp-config --mcp-config ./.mcp.json \
  --allowedTools mcp__world__ailang-read,mcp__world__ailang-write,mcp__world__ailang-edit,mcp__world__ailang-check,mcp__world__ailang-run,mcp__world__builtins-search,mcp__world__examples-search,mcp__world__ailang-cli
```

codex (registration measured on codex-cli 0.159.2; end-to-end run not yet verified). Disabling codex's own shell is the operator's job.

```bash
codex mcp add world --url http://127.0.0.1:7644/mcp/ --bearer-token-env-var WORLD_SESSION
```

A2A:
- agent card `GET /.well-known/agent.json` (with the bearer header); skill ids are the tool ids
- `POST /a2a/` method `tasks/send`: `params.id` = your task id, `params.metadata.skill_id` = the tool, `params.message.parts` = exactly one `{"type":"data","data":<arguments>}`
- the result is in `artifacts[0].parts[0].data`
- the task id is an idempotency key: resending a committed task id returns the committed result without running again

## 3. The eight tools

Every tool: one effect at scope `worktree`, cost 1 call; unknown argument keys are refused; paths are relative, non-empty, with no `..` segment and no leading `-`.

| Tool | Arguments (`?` = optional) | Effect (grant) | Returns (plus `world`) |
|---|---|---|---|
| `ailang-read` | `path` | Workspace.Read | `ok`, `content` |
| `ailang-write` | `path`, `content` (full file) | Workspace.Write | `ok` |
| `ailang-edit` | `path`, `old_text` (must occur exactly once), `new_text` | Workspace.Write | `ok` |
| `ailang-check` | `path` | Ailang.Check | `ok`, `passed`, `error_count`, `errors[{code, message, file}]` |
| `ailang-run` | `path`, `args_json?` (a string holding JSON) | Ailang.Run | `admitted`, `exit_code`, `decision`, `limit`, `stdout`, `stderr` |
| `builtins-search` | `query?`, `module?` | Ailang.Discover | `count`, `matches[{name, module, effect}]` (max 10 with a query) |
| `examples-search` | `query` | Ailang.Discover | `ok`, `argv`, `exit_code`, `stdout`, `stderr` |
| `ailang-cli` | `op`, `path?`, `module?`, `query?`, `package?`, `flags?` (object of strings) | Ailang.CLI | `ok`, `argv`, `exit_code`, `stdout`, `stderr` |

Handler-backed results (read, write, edit, examples-search, ailang-cli) also carry `tool` (tool binary ref) and `policy_digest`.

`ailang-cli` ops: agent_prompt, ai_check, axioms, builtins_list, builtins_show, check, devtools_prompt, docs_search, examples_list, examples_search, examples_show, examples_tags, fmt, iface, pkg_docs, policy_check, prompt, test, tree, version.
- fields: check/ai_check/fmt/tree/policy_check take `path`; iface/pkg_docs take `module` (e.g. "std/io"); docs_search/examples_search take `query`; examples_show/builtins_show take `module` holding a name; test takes `path?`, `package?`
- flags only where the op admits them; boolean flags take "" (e.g. `{"compact": ""}` for iface, `{"limit": "5"}` for docs_search)
- `fmt` never writes: its `write` flag is refused; save the formatted `stdout` with `ailang-write`
- read, write, edit and run are not ops; use their own tools

`ailang-run` executes `ailang run --policy <operator policy> [--args-json J] -- <path>` from the worktree root: restricted mode, capabilities IO and FS only, FS confined to the worktree, at most 1000 FS operations, 8 second timeout.

## 4. Rules of the road

1. Read a file before editing it. Quote enough of it in `old_text` to be unique.
2. Run `ailang-check` after every change to an `.ail` file, and before `ailang-run`.
3. Look up names with `builtins-search`, `examples-search`, or `ailang-cli` (`iface`, `builtins_show`, `docs_search`); do not guess.
4. A refusal is a normal result. Change what it names; never resend the identical call.
5. Budgets are counted in calls per effect, persist across reconnects and restarts, and are not refunded. On `denied:budget`, stop using that effect and report to your operator.
6. Read `stdout` and `stderr` of a run, not only `exit_code`: a refused inner effect can exit 0.
7. After a JSON-RPC error (not a result), the effect may already have run. Inspect the worktree or the log before repeating a write, edit or run.

## 5. Result shape

Every result object has a `world` member:

```json
{"world": {"effects": [{"id": "e1", "status": "ok", "record": "sha256:..."}], "plan": "sha256:..."}}
```

- `plan`: content address of the effect plan your call produced
- `effects[].status`: `ok` (the handler ran; its output may still say `ok: false`), `denied` (the broker refused; nothing ran), or `failed` (the handler failed)
- `effects` is `[]` when the plan refused your arguments; nothing ran and nothing was spent

Real `ailang-run` result (attended smoke, digests shortened):

```json
{"admitted": true, "exit_code": 0, "limit": null,
 "decision": {"allowed_caps": ["FS", "IO"], "declared_effects": ["IO"], "function": "main", "ok": true},
 "stdout": "hello from pi\n", "stderr": "policy: {\"ok\":true,...}\n",
 "world": {"effects": [{"id": "e1", "record": "sha256:f932f3aa...", "status": "ok"}], "plan": "sha256:5d0d3fde..."}}
```

`ailang-run` outcomes:
- (a) admitted: `admitted: true`, program output in `stdout`, one `policy:` line in `stderr`
- (b) admitted, then refused at runtime: `admitted: true`, `exit_code` 1, the reason in `stderr` (e.g. `readFile: path "/etc/hosts" escapes sandbox ...`)
- (c) refused before execution: `admitted: false`; read `decision.error_kind` (`policy_violation`, `read_failed`, ...) and `decision.missing_from_policy`; narrow the program's effects to IO and FS
- (d) supervisor limit: `exit_code` 3, `limit` holds `reason` (e.g. `timeout`) and `stage`; not authenticated, treat as a hint

`ailang-check`: `ok: true` means the check ran; `passed` is the verdict. Line and column are inside each error's `message` (e.g. `broken.ail:6:11`).

## 6. Common refusals and what to do

| You see | Do |
|---|---|
| `unknown argument "<k>"; <tool> admits only: ...` | Drop or rename the argument |
| `missing required argument "<k>"` | Supply it |
| `path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"` | Use a path relative to the worktree root |
| `read ...: path "..." escapes sandbox "..."` | The path resolves outside the worktree (symlink or absolute); stay inside |
| `write .git/...: .git/ is read-only to the lane's tools` | You cannot change `.git/` |
| `write ...: matches fs_deny_write "..." — read-only under this policy` | The path is on the deny list; leave it alone |
| `edit ...: old_text occurs N times; include more context so it is unique` | Add surrounding lines to `old_text` |
| `edit ...: old_text not found — the file may have changed; read it again` | Re-read, then edit |
| `op "<op>" is not an ailang-cli op; ...` | Use the dedicated tool (read/write/edit/run) |
| `op "fmt" flag "write" writes the worktree; ...` | Call fmt without `write`, then `ailang-write` the output |
| `op <op> does not admit flag --<f> (admitted: ...)` | Use only the listed flags |
| `no examples corpus configured: ...` | The operator has no corpus; use `ailang-cli` `docs_search` instead |
| effect `status: "denied"` | Look up the reason (`denied:budget`, `denied:scope`, `denied:expired`, `denied:effect-name`) in the effect record; report to your operator |
| MCP error `-32602 unknown tool "<name>"` | You hold no grant for that tool; it is not in `tools/list` |
| MCP error `-32603 host callback failed` or `host callback timed out` | The effect may have run; check before retrying writes or runs |
| A2A error `effects were requested but the invocation was not confirmed committed; effect records: ...` | The effect ran and is recorded, but the call did not commit; do not blindly retry |

## 7. Explain your own actions

Read routes (open in this version; send your bearer header anyway):
- `GET /v1/log?from=<n>&limit=<m>` -> `{"items":[{header, entryHash, transitionRef}]}`
- `GET /v1/log/<index>` -> one entry; `transitionRef` is the invocation record
- `GET /v1/objects/<ref>?payload=true` -> object with base64 `payload`

The invocation record (`world/invocation-record/v2`) names `skillId`, `input` (your arguments), `output` (your result), `plan` and `effects` (effect record refs). An effect record (`world/effect-record/v1`) has `effect`, `scope`, `cost`, `budgetBefore`, `budgetAfter`, `allowed`, `failed`, `denial`, `requestRef`, `resultRef`. MCP and A2A calls are both logged with `writtenBy: "coordinator:a2a"`.

Retries: MCP calls are at-least-once (each `tools/call` runs again and is recorded again); A2A task ids are idempotent within the session. If a response is lost, page the log before resending.

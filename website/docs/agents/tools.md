---
title: Tool reference
sidebar_position: 3
description: The eight AILANG World software-engineering tools, generated from packages/se-tools/transitions.json.
---

# Tool reference

AILANG World serves 8 tools. Each one is a pure AILANG transition that plans exactly one brokered effect, runs it inside your episode's worktree, and commits one World log entry. The MCP tool name and the A2A skill id are both the tool id below.

Rules that apply to every tool:

- **Arguments are strict.** A key outside the schema is refused by the plan, never ignored. The refusal commits a log entry with zero effects.
- **Paths are relative to the worktree root**: non-empty, no `..` segment, no leading `-`. AILANG's policy layer refuses anything that resolves outside the worktree (absolute paths, symlinks leading out), keeps `.git/` read-only, and keeps the operator's deny list read-only.
- **One call costs 1** from the session's budget for the tool's effect. Budgets count calls, not time or tokens.
- **Every result carries `world`**: `{effects:[{id, status, record}], plan}`. See [Results and errors](./results-and-errors.md).

| Tool | Effect (grant needed) | Arguments | Finish phase |
|---|---|---|---|
| [`ailang-read`](#ailang-read) | `Workspace.Read` (scope `worktree`) | `path` | no |
| [`ailang-write`](#ailang-write) | `Workspace.Write` (scope `worktree`) | `path`, `content` | no |
| [`ailang-edit`](#ailang-edit) | `Workspace.Write` (scope `worktree`) | `path`, `old_text`, `new_text` | no |
| [`ailang-check`](#ailang-check) | `Ailang.Check` (scope `worktree`) | `path` | yes |
| [`ailang-run`](#ailang-run) | `Ailang.Run` (scope `worktree`) | `path`, `args_json?` | no |
| [`builtins-search`](#builtins-search) | `Ailang.Discover` (scope `worktree`) | `query?`, `module?` | yes |
| [`examples-search`](#examples-search) | `Ailang.Discover` (scope `worktree`) | `query` | no |
| [`ailang-cli`](#ailang-cli) | `Ailang.CLI` (scope `worktree`) | `op`, `path?`, `module?`, `query?`, `package?`, `flags?` | no |

Six grants cover the eight tools: `ailang-write` and `ailang-edit` share `Workspace.Write`, and the two searches share `Ailang.Discover`. `tools/list` shows only the tools whose effect your session holds a grant for. A grant with budget 0 still lists the tool, but every call is `denied:budget`.

## ailang-read

**Read (sandboxed).** Effect `Workspace.Read`, scope `worktree`, cost 1 per call. Handler: policy-tool `read`. Module: `packages/se-tools/se_tools/read.ail`.

Description served in `tools/list`:

```text
Read a file inside this episode's worktree (policy-tool `read`). The path is relative to the worktree root, non-empty, with no `..` segment and no leading `-`; anything that resolves outside the worktree (absolute, `..`, or a symlink leading out) is refused by AILANG's policy layer, and `.git/` plus the operator's deny list stay read-only. Unknown argument keys are refused, never ignored. Returns {ok, content} or {ok:false, refused}. Costs one Workspace.Read call from the session's budget. Every call commits one World log entry; the result carries `world` (the effect plan and the brokered effect's record ref).
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | File path, relative to the worktree root |

`additionalProperties` is `false`: any other key is refused.

### Result

`ok`, `content`, `tool` (the tool binary ref), `policy_digest`, `world`.

Output schema properties: `content`, `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
missing required argument "path"
argument "path" must be a string
```

From the handler or AILANG's policy layer (the effect ran with status `ok`; the output says `ok: false`):

```text
read /etc/hosts: path "/etc/hosts" escapes sandbox "<worktree>"
read nope.ail: openat nope.ail: no such file or directory
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-read","arguments":{"path":"hello.ail"}}}
```

## ailang-write

**Write (sandboxed).** Effect `Workspace.Write`, scope `worktree`, cost 1 per call. Handler: policy-tool `write`. Module: `packages/se-tools/se_tools/write.ail`.

Description served in `tools/list`:

```text
Create or overwrite a file inside this episode's worktree (policy-tool `write`). The path is relative to the worktree root, non-empty, with no `..` segment and no leading `-`; anything that resolves outside the worktree (absolute, `..`, or a symlink leading out) is refused by AILANG's policy layer, and `.git/` plus the operator's deny list stay read-only. Unknown argument keys are refused, never ignored. Returns {ok} or {ok:false, refused}. Costs one Workspace.Write call from the session's budget. Every call commits one World log entry; the result carries `world` (the effect plan and the brokered effect's record ref).
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | File path, relative to the worktree root |
| `content` | string | yes | The complete new file content |

`additionalProperties` is `false`: any other key is refused.

### Result

`ok`, `tool`, `policy_digest`, `world`.

Output schema properties: `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
missing required argument "content"
unknown argument "old_text"; ailang-write admits only: path, content
```

From the handler or AILANG's policy layer (the effect ran with status `ok`; the output says `ok: false`):

```text
write .git/config: .git/ is read-only to the lane's tools
write .ailang/x.ail: matches fs_deny_write ".ailang/**" — read-only under this policy
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-write","arguments":{"path":"notes/todo.txt","content":"first line\n"}}}
```

## ailang-edit

**Edit (sandboxed).** Effect `Workspace.Write`, scope `worktree`, cost 1 per call. Handler: policy-tool `edit`. Module: `packages/se-tools/se_tools/edit.ail`.

Description served in `tools/list`:

```text
Replace old_text with new_text in a file inside this episode's worktree (policy-tool `edit`); old_text must occur exactly once (include enough context). The path is relative to the worktree root, non-empty, with no `..` segment and no leading `-`; anything that resolves outside the worktree (absolute, `..`, or a symlink leading out) is refused by AILANG's policy layer, and `.git/` plus the operator's deny list stay read-only. Unknown argument keys are refused, never ignored. Returns {ok} or {ok:false, refused}. Costs one Workspace.Write call from the session's budget. Every call commits one World log entry; the result carries `world` (the effect plan and the brokered effect's record ref).
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | File path, relative to the worktree root |
| `old_text` | string | yes | Exact text to replace (must occur exactly once) |
| `new_text` | string | yes | Replacement text |

`additionalProperties` is `false`: any other key is refused.

### Result

`ok`, `tool`, `policy_digest`, `world`.

Output schema properties: `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
missing required argument "new_text"
unknown argument "old"; ailang-edit admits only: path, old_text, new_text
```

From the handler or AILANG's policy layer (the effect ran with status `ok`; the output says `ok: false`):

```text
edit hello.ail: old_text occurs 8 times; include more context so it is unique
edit hello.ail: old_text not found — the file may have changed; read it again
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-edit","arguments":{"path":"hello.ail","old_text":"println(\"hello from ailang-run\")","new_text":"println(\"hello from pi\")"}}}
```

## ailang-check

**AILANG Check.** Effect `Ailang.Check`, scope `worktree`, cost 1 per call. Handler: policy-tool `ai_check`. Module: `packages/se-tools/se_tools/check.ail`.

Description served in `tools/list`:

```text
Type-check and Z3-verify an AILANG file inside this episode's worktree (policy-tool `ai_check`) and return STRUCTURED diagnostics: {ok:true, passed, error_count, errors:[{code, message, file}]} (line and column stay inside message), or {ok:false, status, refused}. The path is relative to the worktree root, non-empty, with no `..` segment and no leading `-`; anything that resolves outside the worktree (absolute, `..`, or a symlink leading out) is refused by AILANG's policy layer, and `.git/` plus the operator's deny list stay read-only. Unknown argument keys are refused, never ignored. Costs one Ailang.Check call from the session's budget. Every call commits one World log entry; the result carries `world` (the effect plan and the brokered effect's record ref).
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | Path to the .ail file, relative to the worktree root |

`additionalProperties` is `false`: any other key is refused.

### Result

`ok`, `passed`, `error_count`, `errors`, `world` (the finish output; no `tool` or `policy_digest`).

The finish phase parses `ai_check`'s JSON report into `ok: true`, `passed`, `error_count` and `errors` (each `code`, `message`, `file`; line and column stay inside `message`). `ok: true` means the check ran; `passed` is the verdict. When no report is produced the result is `ok: false` with `status` and `refused` (for example `the Ailang.Check effect ended denied`).

Output schema properties: `passed`, `error_count`, `errors`, `status`, `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
unknown argument "limit"; ailang-check admits only: path
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-check","arguments":{"path":"hello.ail"}}}
```

## ailang-run

**AILANG Run (policy-gated).** Effect `Ailang.Run`, scope `worktree`, cost 1 per call. Handler: `ailang run --policy <policy> [--args-json J] -- <path>` (cwd = the worktree root). Module: `packages/se-tools/se_tools/run.ail`.

Description served in `tools/list`:

```text
Execute an AILANG program inside this episode's worktree under the operator's policy: `ailang run --policy <policy> [--args-json J] -- <path>`, from the worktree root. The program's declared effect row must be a subset of the policy's allowed_caps (IO, FS); FS stays inside the worktree; the run is bounded by the policy's timeout (8 s). Returns {admitted, exit_code, decision, limit, stdout, stderr}. Denied programs never execute; read `decision.missing_from_policy` and narrow the program's effects. A refused inner effect can still exit 0, so read stdout and stderr, not only exit_code. The program's inner effects are governed by the AILANG policy layer, not brokered one by one. The path is relative to the worktree root, non-empty, with no `..` segment and no leading `-`; anything that resolves outside the worktree (absolute, `..`, or a symlink leading out) is refused by AILANG's policy layer, and `.git/` plus the operator's deny list stay read-only. Unknown argument keys are refused, never ignored. Costs one Ailang.Run call from the session's budget. Every call commits one World log entry; the result carries `world` (the effect plan and the brokered effect's record ref).
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | Path to the .ail file, relative to the worktree root |
| `args_json` | string | no | JSON arguments for the entrypoint (passed as --args-json) |

`additionalProperties` is `false`: any other key is refused.

### Result

`admitted`, `exit_code`, `decision`, `limit`, `stdout`, `stderr`, `world`. See [Results and errors](./results-and-errors.md#ailang-run-outcomes) for the four outcome shapes.

Output schema properties: `admitted`, `exit_code`, `decision`, `limit`, `stdout`, `stderr`, `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
argument "args_json" must be a string holding JSON
argument "args_json" must be a string
unknown argument "caps"; ailang-run admits only: path, args_json
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-run","arguments":{"path":"hello.ail"}}}
```

## builtins-search

**Search AILANG builtins.** Effect `Ailang.Discover`, scope `worktree`, cost 1 per call. Handler: policy-tool `builtins_list` (text form, no flags). Module: `packages/se-tools/se_tools/builtins_search.ail`.

Description served in `tools/list`:

```text
Search the REAL builtin inventory compiled into the tool binary (policy-tool `builtins_list`, text form). Pass query (case-insensitive substring over name/module/effect; at most 10 matches) and/or module (e.g. 'std/fs'); omit both for the full list. Use this instead of guessing builtin names. Returns {count, matches}, each match {name, module, effect}. Degraded: no signature or description is returned or searched (policy-tool caps a CLI op's stdout at 64 KiB and the JSON inventory exceeds it); use ailang-cli `builtins_show` or examples-search for usage. A truncated or unparseable inventory is refused, never an empty result. Unknown argument keys are refused, never ignored. Costs one Ailang.Discover call from the session's budget. Every call commits one World log entry; the result carries `world` (the effect plan and the brokered effect's record ref).
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `query` | string | no | Substring to look for |
| `module` | string | no | Module filter, e.g. std/fs |

`additionalProperties` is `false`: any other key is refused.

### Result

`count`, `matches` (each `name`, `module`, `effect`), `world`.

The finish phase parses the `Total: N builtins` header and one `name [effect] module` line per builtin, then keeps entries where `module` contains the `module` argument and `query` is a substring of name, module or effect (both case-insensitive). With a `query`, at most 10 matches are returned; `count` is the number returned, not the inventory total. A truncated or unparseable inventory is `ok: false` with `status` and `refused`, never an empty list. Matches carry no signature or description: use `ailang-cli` op `builtins_show` for those.

Output schema properties: `count`, `matches`, `status`, `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
unknown argument "limit"; builtins-search admits only: query, module
argument "query" must be a string
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"builtins-search","arguments":{"query":"println"}}}
```

## examples-search

**Search AILANG examples.** Effect `Ailang.Discover`, scope `worktree`, cost 1 per call. Handler: policy-tool `examples_search` over the corpus named by `serve --examples-dir`. Module: `packages/se-tools/se_tools/examples_search.ail`.

Description served in `tools/list`:

```text
Search the AILANG examples corpus this daemon serves (policy-tool `examples_search`; the operator's --examples-dir) for `query`, then read a hit with ailang-cli op examples_show. The op admits no flags, so there is no limit argument. Use this before writing a construct you are unsure of. Returns {ok, argv, exit_code, stdout, stderr}, or {ok:false, refused} when no corpus is configured. Unknown argument keys are refused, never ignored. Costs one Ailang.Discover call from the session's budget. Every call commits one World log entry; the result carries `world` (the effect plan and the brokered effect's record ref).
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `query` | string | yes | Text to search the examples for |

`additionalProperties` is `false`: any other key is refused.

### Result

`ok`, `argv`, `exit_code`, `stdout`, `stderr`, `tool`, `policy_digest`, `world`. Read a hit with `ailang-cli` op `examples_show`.

Output schema properties: `stdout`, `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
missing required argument "query"
unknown argument "limit"; examples-search admits only: query
```

From the handler or AILANG's policy layer (the effect ran with status `ok`; the output says `ok: false`):

```text
no examples corpus configured: start ailang-worldd serve with --examples-dir DIR (an AILANG examples corpus: manifest.json plus runnable/; `ailang examples download` makes one at ~/.ailang/examples)
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"examples-search","arguments":{"query":"fold"}}}
```

## ailang-cli

**AILANG CLI (policy-allowlisted).** Effect `Ailang.CLI`, scope `worktree`, cost 1 per call. Handler: policy-tool, `op` from the binary's `cli` list. Module: `packages/se-tools/se_tools/cli.ail`.

Description served in `tools/list`:

```text
Run an allowlisted ailang operation as a typed policy-tool request. ops: agent_prompt, ai_check, axioms, builtins_list, builtins_show, check, devtools_prompt, docs_search, examples_list, examples_search, examples_show, examples_tags, fmt, iface, pkg_docs, policy_check, prompt, test, tree, version. Field by op: check/ai_check/fmt/tree/policy_check: {path}; iface/pkg_docs: {module}; docs_search/examples_search: {query}; examples_show/builtins_show: {module: <name>}; test: {path?, package?}; flags: {"json": "", "limit": "5"} only where the op admits them (boolean flags take ""). fmt's write flag is refused: ailang-cli never writes (fmt returns the formatted text; save it with ailang-write). NOT read/write/edit (their own tools) and NOT run (ailang-run). Paths and package directories follow the worktree path rule. Unknown argument keys are refused, never ignored. Returns {ok, argv, exit_code, stdout, stderr} or {ok:false, refused}. Costs one Ailang.CLI call from the session's budget. Every call commits one World log entry; the result carries `world` (the effect plan and the brokered effect's record ref).
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `op` | string | yes | The operation, e.g. "check", "iface", "docs_search", "test" |
| `path` | string | no | In-worktree file path (check, ai_check, fmt, tree, policy_check, test) |
| `module` | string | no | Module path or name (iface: "std/fs"; examples_show/builtins_show: a name) |
| `query` | string | no | Search text (docs_search, examples_search) |
| `package` | string | no | Package directory for test (in-worktree) |
| `flags` | object of string | no | Admitted flags by name; boolean flags take "" (e.g. \{"json": ""\}) |

`additionalProperties` is `false`: any other key is refused.

Admitted ops and their flags (v0.51.0 tool binary). Pass flags as `{"name": "value"}`; boolean flags take `""`.

| `op` | Fields | Admitted flags |
|---|---|---|
| `agent_prompt` | none | none |
| `ai_check` | `path` | --timeout |
| `axioms` | none | none |
| `builtins_list` | none | --json |
| `builtins_show` | `module` (a name) | none |
| `check` | `path` | --json --quiet --strict-syntax |
| `devtools_prompt` | none | none |
| `docs_search` | `query` | --json --limit |
| `examples_list` | none | --status --tag |
| `examples_search` | `query` | none |
| `examples_show` | `module` (a name) | none |
| `examples_tags` | none | none |
| `fmt` | `path` | --check (--write is refused) |
| `iface` | `module` | --compact |
| `pkg_docs` | `module` | none |
| `policy_check` | `path` | none |
| `prompt` | none | none |
| `test` | `path`, `package` (both optional) | --allow-skips --json --no-color --package |
| `tree` | `path` | none |
| `version` | none | none |

### Result

`ok`, `argv`, `exit_code`, `stdout`, `stderr`, `tool`, `policy_digest`, `world`.

Output schema properties: `argv`, `exit_code`, `stdout`, `stderr`, `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
op "read" is not an ailang-cli op; reads, writes and edits have their own tools, and run is ailang-run
op "fmt" flag "write" writes the worktree; ailang-cli never writes (fmt without it returns the formatted text; write it with ailang-write)
argument "flags" must be an object of string values
missing required argument "op"
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
```

From the handler or AILANG's policy layer (the effect ran with status `ok`; the output says `ok: false`):

```text
op check does not admit flag --zz (admitted: --json --quiet --strict-syntax)
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-cli","arguments":{"op":"iface","module":"std/io","flags":{"compact":""}}}}
```

---

This page is generated by `website/scripts/gen-tool-reference.mjs` from `packages/se-tools/transitions.json` (refusal texts transcribed from `packages/se-tools/se_tools/*.ail`). Do not edit it by hand: change the sources and run `node website/scripts/gen-tool-reference.mjs` from the repo root.

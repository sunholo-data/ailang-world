---
title: Tool reference
sidebar_position: 3
description: The 10 AILANG World software-engineering tools, generated from packages/se-tools/transitions.json.
---

# Tool reference

AILANG World serves 10 tools. Each one is a pure AILANG transition that plans exactly one brokered effect, runs it inside your episode's worktree, and commits one World log entry. The MCP tool name and the A2A skill id are both the tool id below.

Rules that apply to every tool:

- **Arguments are strict.** A key outside the schema is refused by the plan, never ignored. The refusal commits a log entry with zero effects.
- **Paths are relative to the worktree root**: non-empty, no `..` segment, no leading `-`. AILANG's policy layer refuses anything that resolves outside the worktree (absolute paths, symlinks leading out), keeps `.git/` read-only, and keeps the operator's deny list read-only.
- **One call costs 1** from the session's budget for the tool's effect. Budgets count calls, not time or tokens.
- **Every result carries `world`**: `{effects:[{id, status, record}], plan}`. See [Results and errors](./results-and-errors.md).

| Tool | Effect (grant needed) | Arguments | Finish phase |
|---|---|---|---|
| [`ailang-read`](#ailang-read) | `Workspace.Read` (scope `worktree`) | `path`, `offset?`, `limit?` | no |
| [`ailang-write`](#ailang-write) | `Workspace.Write` (scope `worktree`) | `path`, `content` | no |
| [`ailang-edit`](#ailang-edit) | `Workspace.Write` (scope `worktree`) | `path`, `old_text`, `new_text` | no |
| [`ailang-delete`](#ailang-delete) | `Workspace.Delete` (scope `worktree`) | `path` | no |
| [`ailang-check`](#ailang-check) | `Ailang.Check` (scope `worktree`) | `path` | yes |
| [`ailang-run`](#ailang-run) | `Ailang.Run` (scope `worktree`) | `path`, `args_json?`, `stdin?`, `argv?`, `caps?` | no |
| [`builtins-search`](#builtins-search) | `Ailang.Discover` (scope `worktree`) | `query?`, `module?` | yes |
| [`examples-search`](#examples-search) | `Ailang.Discover` (scope `worktree`) | `query` | no |
| [`ailang-cli`](#ailang-cli) | `Ailang.CLI` (scope `worktree`) | `op`, `path?`, `module?`, `query?`, `package?`, `flags?` | no |
| [`workspace-exec`](#workspace-exec) | `Workspace.Exec` (scope `worktree`) | `command`, `args?` | no |

8 grants cover the 10 tools: `ailang-write` and `ailang-edit` share `Workspace.Write`, the two searches share `Ailang.Discover`, and `ailang-delete` and `workspace-exec` have their own `Workspace.Delete` and `Workspace.Exec` (deleting is a different authority from writing). Two more, `Ailang.RunEnv` and `Ailang.RunNet`, let `ailang-run` take the `Env` or `Net` capability when the operator enables it. `tools/list` shows only the tools whose effect your session holds a grant for. A grant with budget 0 still lists the tool, but every call is `denied:budget`.

## ailang-read

**Read (sandboxed).** Effect `Workspace.Read`, scope `worktree`, cost 1 per call. Handler: policy-tool `read`. Module: `packages/se-tools/se_tools/read.ail`.

Description served in `tools/list`:

```text
Read an episode-worktree file. Path: non-empty, relative, no `..` segments or leading `-`; policy confines symlinks and protected paths. Unknown keys refuse. Returns content or refusal plus World plan/effect refs; costs one Workspace.Read and commits one entry.
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | File path, relative to the worktree root |
| `offset` | integer | no | First line, 1-based |
| `limit` | integer | no | Line count, at most 2000 |

`additionalProperties` is `false`: any other key is refused.

### Result

`ok`, `content`, `tool` (the tool binary ref), `policy_digest`, `world`. A ranged read (`offset` and/or `limit`) also returns `total_lines` (the file's line count), `offset` and `limit` (the lines returned). The tool binary clamps `limit` to 2000 (and to the end of the file), and an `offset` past the end returns no `content`.

Output schema properties: `content`, `total_lines`, `offset`, `limit`, `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
missing required argument "path"
argument "path" must be a string
argument "offset" must be an integer
argument "limit" must be an integer >= 1 (got 0)
unknown argument "lines"; ailang-read admits only: path, offset, limit
```

From the handler or AILANG's policy layer (the effect ran with status `ok`; the output says `ok: false`):

```text
read /etc/hosts: path "/etc/hosts" escapes sandbox "<worktree>"
read nope.ail: openat nope.ail: no such file or directory
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-read","arguments":{"path":"docs/big.md","offset":201,"limit":100}}}
```

## ailang-write

**Write (sandboxed).** Effect `Workspace.Write`, scope `worktree`, cost 1 per call. Handler: policy-tool `write`. Module: `packages/se-tools/se_tools/write.ail`.

Description served in `tools/list`:

```text
Create/overwrite an episode-worktree file with complete content. Path: non-empty, relative, no `..` segments or leading `-`; policy confines symlinks and protected paths. Unknown keys refuse. Returns ok or refusal plus World plan/effect refs; costs one Workspace.Write and commits one entry.
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
Replace old_text with new_text in an episode-worktree file; old_text must occur exactly once (include enough context). Path: non-empty, relative, no `..` segments or leading `-`; policy confines symlinks and protected paths. Unknown keys refuse. Returns ok or refusal plus World plan/effect refs; costs one Workspace.Write and commits one entry.
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

## ailang-delete

**Delete (sandboxed).** Effect `Workspace.Delete`, scope `worktree`, cost 1 per call. Handler: policy-tool `delete`. Module: `packages/se-tools/se_tools/delete.ail`.

Description served in `tools/list`:

```text
Delete one episode-worktree file; a directory refuses, and a symlink is removed itself, never its target. Path: non-empty, relative, no `..` segments or leading `-`; policy confines symlinks and protected paths. Unknown keys refuse. Returns ok or refusal plus World plan/effect refs; costs one Workspace.Delete and commits one entry.
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | File path, relative to the worktree root |

`additionalProperties` is `false`: any other key is refused.

### Result

`ok`, `tool`, `policy_digest`, `world`. A symlink is removed itself, never its target.

Output schema properties: `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
missing required argument "path"
unknown argument "recursive"; ailang-delete admits only: path
```

From the handler or AILANG's policy layer (the effect ran with status `ok`; the output says `ok: false`):

```text
delete .github/ci.yml: matches fs_deny_write ".github/**" — read-only under this policy
delete .git: .git/ is read-only to the lane's tools
delete tmp: is a directory — delete removes one file at a time
delete nope.txt: statat nope.txt: no such file or directory
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-delete","arguments":{"path":"tmp/helper_test.go"}}}
```

## ailang-check

**AILANG Check.** Effect `Ailang.Check`, scope `worktree`, cost 1 per call. Handler: policy-tool `ai_check`. Module: `packages/se-tools/se_tools/check.ail`.

Description served in `tools/list`:

```text
Type-check and Z3-verify an episode-worktree AILANG file. Returns structured {ok, passed, error_count, errors:[{code,message,file}]} (line/column in message), or {ok:false,status,refused}. Path: non-empty, relative, no `..` segments or leading `-`; policy confines symlinks and protected paths. Unknown keys refuse. Costs one Ailang.Check; commits one entry with World plan/effect refs.
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | Path to the .ail file, relative to the worktree root |

`additionalProperties` is `false`: any other key is refused.

### Result

`ok`, `passed`, `error_count`, `errors`, `world` (the finish output; no `tool`, `policy_digest` or `worktree_delta`). When the check changed the worktree, the `worktree_delta` is in the effect's own output: resolve the effect record.

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

**AILANG Run (policy-gated).** Effect `Ailang.Run, Ailang.RunEnv, Ailang.RunNet`, scope `worktree`, cost 1 per call. Handler: `ailang run --policy <policy> [--args-json J] -- <path> [-- <argv...>]` with `stdin` piped in (cwd = the worktree root). The policy is the episode's base policy for an IO/FS run, otherwise a per-cap-set variant World renders and AILANG's `policy-tool summary` verifies before its first use. Module: `packages/se-tools/se_tools/run.ail`.

Description served in `tools/list`:

```text
Run an episode-worktree AILANG program under policy (8 s timeout). caps is exact, default IO/FS; never Env with Net. Extra caps require operator enablement; Env/Net also require Ailang.RunEnv/Ailang.RunNet grants, otherwise Ailang.Run. Net reaches only operator-allowed loopback host:port pairs; other targets/redirects refuse. args_json supplies entrypoint JSON; stdin is text (absent=empty); argv follows inner `--` (schema limits). Declared effects must fit policy; unrequested FS has budget 0 and fails at use. Denied programs never execute: inspect decision.missing_from_policy, narrow effects or request authorized caps. An inner refusal may exit 0: inspect stdout/stderr. Inner effects use policy, not individual broker calls. Use a non-empty relative path without `..` segments or leading `-`; policy confines symlinks/protected paths. Unknown keys refuse. Returns admission, exit/limit, stdout/stderr, policy and World plan/effect refs; costs one call of its effect and commits one entry.
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | Path to the .ail file, relative to the worktree root |
| `args_json` | string | no | JSON arguments for the entrypoint (passed as --args-json) |
| `stdin` | string | no | The program's whole standard input, at most 65536 UTF-8 bytes (absent = empty) |
| `argv` | array of string | no | Program arguments (getArgs), passed after an inner `--`; at most 32, each at most 1024 bytes, no NUL |
| `caps` | array of Declassify / Env / FS / IO / Net | no | The run's exact capability set (default IO, FS); never Env with Net |

`additionalProperties` is `false`: any other key is refused.

### Result

`admitted`, `exit_code`, `decision`, `limit`, `stdout`, `stderr`, `policy`, `worktree_delta` (only when the run changed the worktree), `world`. `policy` is `{digest, security_mode, caps, net_allow}` of the policy the run executed under. See [Results and errors](./results-and-errors.md#ailang-run-outcomes) for the four outcome shapes.

Which effect a call spends follows `caps`: with `Env` it is `Ailang.RunEnv`, with `Net` it is `Ailang.RunNet`, otherwise `Ailang.Run` (a `Declassify` run included). Each needs its own grant, and each capability beyond `IO` and `FS` must also be enabled by the operator (`serve --run-allow-caps`), or the handler refuses the call before anything runs (the effect is recorded `failed`). A Net run reaches only the loopback `host:port` pairs the operator named with `--run-net-allow`. See [Tool confinement](../security/tool-confinement.md#ailang-run-capabilities).

Output schema properties: `admitted`, `exit_code`, `decision`, `limit`, `stdout`, `stderr`, `ok`, `refused`, `world`, `policy` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
argument "args_json" must be a string holding JSON
argument "args_json" must be a string
argument "stdin" is 65537 bytes; the limit is 65536
argument "argv" has 33 items; the limit is 32
an argv item is 1025 bytes; the limit is 1024
an argv item contains a NUL byte
unknown capability "Process"; caps admits only: Declassify, Env, FS, IO, Net
capability "IO" is listed twice
caps names both Env and Net; one run takes at most one of them
unknown argument "env"; ailang-run admits only: path, args_json, stdin, argv, caps
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ailang-run","arguments":{"path":"benchmark/solution.ail","caps":["IO"],"stdin":"1\n2\n3\n4\n5\n"}}}
```

## builtins-search

**Search AILANG builtins.** Effect `Ailang.Discover`, scope `worktree`, cost 1 per call. Handler: policy-tool `builtins_list` (text form, no flags). Module: `packages/se-tools/se_tools/builtins_search.ail`.

Description served in `tools/list`:

```text
Search the tool binary's real builtin inventory; query is a case-insensitive name/module/effect substring (max 10 matches), module filters; omit both for all. Returns {count,matches:[{name,module,effect}]}; degraded: no signatures/descriptions returned or searched. Use ailang-cli builtins_show or examples-search for usage. Truncated/invalid inventory refuses, never empty success. Unknown keys refuse. Costs one Ailang.Discover; commits one entry with World plan/effect refs.
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
Search the configured examples corpus for query; read hits with ailang-cli examples_show. No flags or limit argument. Use before guessing syntax. Returns CLI stdout/stderr or refusal if no corpus is configured. Unknown keys refuse. Costs one Ailang.Discover; commits one entry with World plan/effect refs.
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
Run an allowlisted AILANG op: agent_prompt, ai_check, axioms, builtins_list, builtins_show, check, devtools_prompt, docs_search, examples_list, examples_search, examples_show, examples_tags, fmt, iface, pkg_docs, policy_check, prompt, test, tree, version. Arguments: path for check/ai_check/fmt/tree/policy_check; module for iface/pkg_docs/examples_show/builtins_show; query for searches; optional path/package for test. Op-admitted flags only; booleans take an empty string. Never writes: fmt write refuses; save formatted stdout with ailang-write. Paths/packages obey worktree confinement. Unknown keys refuse. Returns CLI argv/exit/stdout/stderr or refusal plus World plan/effect refs; costs one Ailang.CLI and commits one entry.
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

Admitted ops and their flags (World's audited set, first audited on the v0.52.1 tool binary;
the operator's compatibility gate re-ran the audit on v0.54.15 with no change). Pass flags as `{"name": "value"}`; boolean flags take `""`.

| `op` | Fields | Admitted flags |
|---|---|---|
| `agent_prompt` | none | none |
| `ai_check` | `path` | --timeout |
| `axioms` | none | none |
| `builtins_list` | none | --by-effect --by-module --json --module --query --verbose |
| `builtins_show` | `module` (a name) | none |
| `check` | `path` | --json --quiet --strict-syntax |
| `devtools_prompt` | none | none |
| `docs_search` | `query` | --json --limit |
| `examples_list` | none | --status --tags |
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

`ok`, `argv`, `exit_code`, `stdout`, `stderr`, `tool`, `policy_digest`, `worktree_delta` (only when the op changed the worktree), `world`.

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

## workspace-exec

**Exec (operator-profiled, sandboxed).** Effect `Workspace.Exec`, scope `worktree`, cost 1 per call. Handler: one command of the operator's exec profile (`serve --exec-profile`), run as `<node> <archived srt cli.js> --settings <episode settings> -- /bin/sh -c '/usr/bin/env -i "$@"; exit $?' world-exec <K=V…> <argv0> <prefix…> <args…>` (cwd = the profile's `root` in the worktree; host/broker/handlers_exec.go). Module: `packages/se-tools/se_tools/exec.ail`.

Description served in `tools/list`:

```text
Run one command of the operator's exec profile (e.g. `test`, `vet`, `typecheck`) in this episode's worktree. Commands come only from the operator's exec profile (`serve --exec-profile`): `command` is the id of one of them, never an executable or a shell line, and the profile fixes the program, its leading arguments, the working directory, the environment and the timeout. `args` (at most 16 strings, each 1 to 512 bytes, no NUL; absent = none) must match that command's argument grammar: only the flags it lists, as `-f=v` or `-f v`, and positionals of its declared class (relative paths, package patterns or test files); `--` only where the command declares it. The command runs sandboxed: writes stay inside the worktree and the episode's caches, your home directory, World's state and other episodes are unreadable, and there is no network. No shell sees an argument. Without a configured profile the call is refused (`no exec profile configured`). Unknown argument keys are refused, never ignored. Returns {exit_code, timed_out, limit, duration_ms, stdout, stdout_bytes, stdout_truncated, stdout_sha256, stderr, stderr_bytes, stderr_truncated, stderr_sha256, argv, profile, sandbox} or {ok:false, refused}; long output keeps its head and tail. exit_code is the command's own status, reported, never judged. Costs one Workspace.Exec call from the session's budget. Every call commits one World log entry; the result carries `world` (the effect plan and the brokered effect's record ref).
```

### Arguments

| Name | Type | Required | Description |
|---|---|---|---|
| `command` | string | yes | The id of one command in the operator's exec profile, e.g. "test" |
| `args` | array of string | no | Arguments matched against the command's grammar; at most 16, each 1 to 512 bytes, no NUL |

`additionalProperties` is `false`: any other key is refused.

### Result

`exit_code`, `timed_out`, `limit`, `duration_ms`, `stdout`, `stdout_bytes`, `stdout_truncated`, `stdout_sha256`, `stderr`, `stderr_bytes`, `stderr_truncated`, `stderr_sha256`, `argv`, `profile`, `sandbox`, `worktree_delta` (only when the command changed the worktree), `world`. `exit_code` is the command's own status (a signal death is 128 + n), reported, never judged; each stream keeps its first 8 KiB and last 56 KiB with its total size and sha256.

The operator writes the profile: the program, its fixed leading arguments, the flags and positional class `args` must match, the environment, the read roots and the timeout (at most 9 s). The command runs under `srt` with writes confined to the worktree and the episode's exec cache, your home directory, World's state and every other episode unreadable, and no network. See the operator guide, `docs/QUICKSTART.md` §10.

Output schema properties: `exit_code`, `timed_out`, `limit`, `duration_ms`, `stdout`, `stdout_bytes`, `stdout_truncated`, `stdout_sha256`, `stderr`, `stderr_bytes`, `stderr_truncated`, `stderr_sha256`, `argv`, `profile`, `sandbox`, `ok`, `refused`, `world` (required: `world`).

### Refusals you will see

From the plan (zero effects, `world.effects` is `[]`, still committed):

```text
command "Test" refused: a command id must match ^[a-z][a-z0-9-]{0,31}$ (the id of a command in the operator's exec profile)
argument "args" has 17 items; the limit is 16
an args item contains a NUL byte
an args item is 0 bytes; each must be 1 to 512
argument "args" must be an array of strings
missing required argument "command"
unknown argument "cwd"; workspace-exec admits only: command, args
```

From the handler or AILANG's policy layer (the effect ran with status `ok`; the output says `ok: false`):

```text
no exec profile configured: start ailang-worldd serve with --exec-profile FILE
command "deploy" is not a command of exec profile "ailang-compiler"
command "test": argument 0 "-exec=sh" refused: the command does not list this flag
```

### Example call

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"workspace-exec","arguments":{"command":"test","args":["-run","TestLex","./internal/lexer/"]}}}
```

---

This page is generated by `website/scripts/gen-tool-reference.mjs` from `packages/se-tools/transitions.json` (refusal texts transcribed from `packages/se-tools/se_tools/*.ail`). Do not edit it by hand: change the sources and run `node website/scripts/gen-tool-reference.mjs` from the repo root.

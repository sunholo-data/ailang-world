---
title: If you are an AI agent
sidebar_position: 1
description: What AILANG World gives an AI agent, what it withholds, and how to work inside it.
---

# If you are an AI agent

This section is written for you: a coding agent (pi, Claude Code, codex or any MCP or A2A client) that an operator has pointed at an AILANG World daemon. Read this page first. If you can fetch only one file, fetch [`/AGENTS.md`](pathname:///AGENTS.md); it holds the same rules in one plain-text page.

AILANG World is the 1.0 release candidate (`v1.0.0-rc.1`). The tool surface below is live. pi and Claude Code first exercised it end to end on 2026-10-03, and in the pre-registered non-inferiority run of 2026-10-15, Claude Code and codex each passed 69/69 tasks through World's tools, the same as with their native tools, with no added wall-clock overhead. Treat it as working software with known residuals (see [Residual risks](../security/residuals.md)).

## What World gives you

- **A confined worktree.** The operator creates one git worktree per episode. Your session is bound to that episode, and every path you pass is relative to its root.
- **Ten tools**, over MCP at `/mcp/` and over A2A at `/a2a/`:

  | Tool | What it does |
  |---|---|
  | `ailang-read` | Read a file in the worktree, whole or a line range (`offset`, `limit`) |
  | `ailang-write` | Create or overwrite a file in the worktree |
  | `ailang-edit` | Replace one exact, unique occurrence of text in a file |
  | `ailang-delete` | Delete one file in the worktree (its own `Workspace.Delete` grant) |
  | `ailang-check` | Type-check and Z3-verify an `.ail` file; structured diagnostics |
  | `ailang-run` | Run an AILANG program under the operator's policy (IO and FS by default; Env, Net or Declassify only if the operator enables them; 8 s) |
  | `builtins-search` | Search the real builtin inventory of the tool binary |
  | `examples-search` | Search the AILANG examples corpus the daemon serves |
  | `ailang-cli` | The rest of the `ailang` CLI as typed, allowlisted operations |
  | `workspace-exec` | Run one command from the operator's exec profile (for example `go test`), by id, sandboxed; refused unless the operator configured a profile |

- **A record of every call.** Each tool call commits one entry to World's append-only log. The entry links your arguments, the effect plan, the brokered effect record and the result, all content-addressed (`sha256:…`). You can read these back to explain what you did; see [Provenance](./provenance.md).

## What World does not give you

- **No shell.** You cannot run a shell line. AILANG programs run only through `ailang-run`, under the operator's policy. Other commands run only through `workspace-exec`, and only those the operator's exec profile lists: you pass a command id and arguments that must match its grammar, never a program or a shell line. The command runs under the srt sandbox, with writes confined to your worktree and the episode's caches. Without an exec profile, every `workspace-exec` call is refused `no exec profile configured`.
- **No network by default.** The rendered policy allows the `IO` and `FS` capabilities only, unless the operator enables more (`serve --run-allow-caps`). A `Net` run is possible only when the operator allows it, and then reaches only the loopback `IP:PORT` pairs the operator lists (`--run-net-allow`); every other host and redirect is refused. `workspace-exec` commands have no network.
- **No writes outside your worktree.** Absolute paths, `..` segments and symlinks that lead out are refused. `.git/` is read-only. The operator's deny list (`.github/`, `.pi/`, `.claude/`, `.ailang/`, `.gitmodules`, `.gitattributes`, in every case variant) is read-only.
- **No git commits, no new worktrees.** Committing and provisioning are operator actions.
- **Only the tools you were granted.** `tools/list` returns exactly the tools whose effect your session holds a grant for. A tool you cannot see cannot be called.

## How to behave

1. **Read before you edit.** `ailang-edit` needs `old_text` that occurs exactly once. Read the file first, then quote enough surrounding text to be unique. If the edit says `old_text not found — the file may have changed; read it again`, do exactly that.
2. **Check before you run.** Call `ailang-check` after every change to an `.ail` file. It costs one `Ailang.Check` call and returns structured errors; `ailang-run` costs one `Ailang.Run` call and gives you less to work with when the program does not type-check.
3. **Look things up instead of guessing.** Use `builtins-search` for builtin names, `examples-search` for idioms, and `ailang-cli` with `docs_search`, `iface` or `builtins_show` for signatures and documentation.
4. **Treat refusals as information.** A refusal is a normal result, not a crash: `{"ok":false,"refused":"…"}` tells you exactly what to change. Fix the argument or the path; do not retry the same call unchanged. Refusals are committed to the log too.
5. **Spend calls deliberately.** Budgets are counted in calls per effect (for example 50 `Workspace.Write` calls), they persist across reconnects and restarts, and a denied call is still recorded. When a call comes back `denied:budget`, stop using that effect and report to your operator.
6. **Read `stdout` and `stderr`, not only `exit_code`.** A program whose inner effect was refused can still exit 0.
7. **Do not blindly retry after an error.** If a tool call fails with a JSON-RPC error rather than a result, the effect may already have happened. Read the log (see [Results and errors](./results-and-errors.md#when-the-call-itself-fails)) before you repeat a write or a run.

## Where to go next

- [Connecting](./connecting.md): exact configuration for pi, Claude Code, codex, curl and A2A.
- [Tool reference](./tools.md): every argument, result field and refusal, generated from the published manifest.
- [Results and errors](./results-and-errors.md): the `world` block, `ailang-run` outcomes, budgets and failure semantics.
- [Workflow examples](./workflow-examples.md): real call sequences from the first attended smoke.
- [Provenance](./provenance.md): reading back your own past actions.
- The AILANG language itself: [ailang.sunholo.com/llms.txt](https://ailang.sunholo.com/llms.txt).

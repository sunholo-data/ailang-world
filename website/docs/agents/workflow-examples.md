---
title: Workflow examples
sidebar_position: 5
description: Worked tool-call sequences for AI agents in AILANG World, from the first attended smoke and measured tool output.
---

# Workflow examples

The first example is the exact sequence pi and Claude Code ran in the attended smoke on 2026-10-03 (evidence: `design_docs/verification/world-attended-2026-10-03-row134-m6/`). Results are copied from the committed invocation outputs in that smoke's store, with long digests shortened to eight characters. The other examples use output measured directly on the v0.51.0 tool binary, reshaped the way World's tools shape it; each one says so.

Calls are shown as MCP `tools/call` params. Over A2A, send the same `arguments` object as the `data` part with `metadata.skill_id` set to the tool name.

## 1. Read, edit, check, run a hello program

The worktree starts with this file, `hello.ail`:

```ailang
module hello

import std/io (println)

export func main() -> () ! {IO} {
  println("hello from ailang-run")
}
```

The prompt was: *Read hello.ail, change its message to 'hello from pi', check it, then run it and report its stdout.* pi made four calls, and each one became one log entry (entries 2 to 5).

**Read the file first.**

```json
{"name": "ailang-read", "arguments": {"path": "hello.ail"}}
```

```json
{"content": "module hello\n\nimport std/io (println)\n\nexport func main() -> () ! {IO} {\n  println(\"hello from ailang-run\")\n}\n",
 "ok": true, "policy_digest": "87d52ed5…", "tool": "sha256:e55ff71c…",
 "world": {"effects": [{"id": "e1", "record": "sha256:0067bdd0…", "status": "ok"}], "plan": "sha256:ce4340a1…"}}
```

**Edit with a unique `old_text`.** The whole `println(…)` call occurs once, so it is a safe anchor.

```json
{"name": "ailang-edit", "arguments": {"path": "hello.ail",
  "old_text": "println(\"hello from ailang-run\")", "new_text": "println(\"hello from pi\")"}}
```

```json
{"ok": true, "policy_digest": "87d52ed5…", "tool": "sha256:e55ff71c…",
 "world": {"effects": [{"id": "e1", "record": "sha256:0af7f905…", "status": "ok"}], "plan": "sha256:fd356d64…"}}
```

**Check before running.**

```json
{"name": "ailang-check", "arguments": {"path": "hello.ail"}}
```

```json
{"error_count": 0, "errors": [], "ok": true, "passed": true,
 "world": {"effects": [{"id": "e1", "record": "sha256:a89b0538…", "status": "ok"}], "plan": "sha256:414ed668…"}}
```

**Run it.**

```json
{"name": "ailang-run", "arguments": {"path": "hello.ail"}}
```

```json
{"admitted": true, "exit_code": 0, "limit": null,
 "decision": {"allowed_caps": ["FS", "IO"], "declared_effects": ["IO"], "function": "main", "ok": true},
 "stdout": "hello from pi\n",
 "stderr": "policy: {\"ok\":true,…}\n",
 "world": {"effects": [{"id": "e1", "record": "sha256:f932f3aa…", "status": "ok"}], "plan": "sha256:5d0d3fde…"}}
```

Claude Code then ran the same four steps (entries 6 to 9) and printed `hello from claude`. Both clients finished in about 9 seconds with no built-in tools enabled.

Two details worth copying. The `plan` ref of the two `ailang-read` calls is identical (`sha256:ce4340a1…`), because the plan is a pure function of the arguments; the effect `record` differs, because each execution is recorded separately. And the run's `decision.declared_effects` is `["IO"]`: the program asked for less than the policy allows, which is what you want.

## 2. Fix a type error from `ailang-check` diagnostics

The diagnostic below was measured with `policy-tool` `ai_check` on the v0.51.0 tool binary. The result shape is what `ailang-check`'s finish phase produces from that report.

Write a file that passes a number where `println` wants a string:

```json
{"name": "ailang-write", "arguments": {"path": "broken.ail",
  "content": "module broken\n\nimport std/io (println)\n\nexport func main() -> () ! {IO} {\n  println(42)\n}\n"}}
```

Check it:

```json
{"name": "ailang-check", "arguments": {"path": "broken.ail"}}
```

```json
{"ok": true, "passed": false, "error_count": 1,
 "errors": [{"code": "ERROR", "file": "broken.ail",
   "message": "type error in broken (decl 0): at literal at broken.ail:6:11: No instance for Num[string] in scope. Arithmetic operators (+, -, *, /) need numbers, but this is a string. Use ++ to concatenate strings, or stringToInt to convert a string to a number."}],
 "world": {"effects": [{"id": "e1", "record": "sha256:…", "status": "ok"}], "plan": "sha256:…"}}
```

How to read it:

- `ok: true` means the check ran. `passed: false` is the verdict. Do not confuse the two.
- The position is inside `message`: `broken.ail:6:11` is line 6, column 11, the literal `42`.
- The message says a numeric literal met a `string`. `println` takes a `string`, so pass one.

Fix it at the position the diagnostic names, then check again:

```json
{"name": "ailang-edit", "arguments": {"path": "broken.ail", "old_text": "println(42)", "new_text": "println(show(42))"}}
```

```json
{"name": "ailang-check", "arguments": {"path": "broken.ail"}}
```

When `passed` is `true` and `error_count` is `0`, move on to `ailang-run`.

## 3. Look up builtins and examples instead of guessing

**`builtins-search`** searches the inventory compiled into the tool binary. Matching `println` against the v0.51.0 inventory (378 builtins) with the finish phase's filter gives:

```json
{"name": "builtins-search", "arguments": {"query": "println"}}
```

```json
{"count": 2,
 "matches": [{"effect": "io", "module": "std/io", "name": "_io_eprintln"},
             {"effect": "io", "module": "std/io", "name": "_io_println"}],
 "world": {"effects": [{"id": "e1", "record": "sha256:…", "status": "ok"}], "plan": "sha256:…"}}
```

Matches carry no signatures. To narrow by module, pass `module` (`{"module": "std/fs"}` lists every `std/fs` builtin; a `query` caps the list at 10).

**`examples-search`** searches the examples corpus the operator configured. Its result is the CLI's own text output in `stdout`:

```json
{"name": "examples-search", "arguments": {"query": "fold"}}
```

If the operator started the daemon without a corpus, you get `{"ok": false, "refused": "no examples corpus configured: …"}`. That is not something you can fix; carry on with `ailang-cli` `docs_search`.

## 4. Use `ailang-cli` for signatures and docs

`ailang-cli` exposes the rest of the AILANG CLI as typed requests. The outputs below were measured with `policy-tool` on the v0.51.0 tool binary; World adds `tool`, `policy_digest` and `world`.

**Module interface** (`iface`, with the admitted `--compact` flag):

```json
{"name": "ailang-cli", "arguments": {"op": "iface", "module": "std/io", "flags": {"compact": ""}}}
```

```json
{"ok": true, "argv": ["iface", "--compact", "std/io"],
 "stdout": "module std/io\neprintln : (string)->()!{IO}\nexit : (int)->()!{IO}\nflush : (())->()!{IO}\nprint : (string)->()!{IO}\nprintErr : (string)->()!{IO}\nprintln : (string)->()!{IO}\nreadLine : (())->string!{IO}\nwriteBytes : (bytes)->()!{IO}\n",
 "tool": "sha256:…", "policy_digest": "…",
 "world": {"effects": [{"id": "e1", "record": "sha256:…", "status": "ok"}], "plan": "sha256:…"}}
```

**One builtin's documentation** (`builtins_show`; the builtin name goes in `module`):

```json
{"name": "ailang-cli", "arguments": {"op": "builtins_show", "module": "_io_println"}}
```

The `stdout` holds the signature (`_io_println: string -> () ! {IO}`), a description, parameters and an example.

**Documentation search** (`docs_search`, with the admitted `--limit` flag):

```json
{"name": "ailang-cli", "arguments": {"op": "docs_search", "query": "println", "flags": {"limit": "5"}}}
```

**Formatting without writing.** `fmt` returns the formatted source in `stdout`. Its `write` flag is refused, so save the text with `ailang-write`:

```json
{"name": "ailang-cli", "arguments": {"op": "fmt", "path": "hello.ail"}}
```

A flag the op does not admit is refused by the policy layer with the admitted list, for example `op check does not admit flag --zz (admitted: --json --quiet --strict-syntax)`. The full table is in the [Tool reference](./tools.md#ailang-cli).

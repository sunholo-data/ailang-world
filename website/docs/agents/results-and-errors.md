---
title: Results and errors
sidebar_position: 4
description: How to read a World tool result, the world provenance block, ailang-run outcomes, budget denials and failed calls.
---

# Results and errors

A World tool call ends in one of two ways:

1. **A result.** The JSON-RPC `result` (MCP) or the task artifact (A2A) holds one JSON object. This covers successes, refusals and denials alike, and every one of them is committed to the log.
2. **A JSON-RPC error.** The call failed in the host. The effect may or may not have run. See [When the call itself fails](#when-the-call-itself-fails).

## The `world` block

Every result object carries a reserved `world` member:

```json
{"world": {"effects": [{"id": "e1", "status": "ok", "record": "sha256:…"}],
           "plan": "sha256:…"}}
```

| Field | Meaning |
|---|---|
| `plan` | Content address of the effect plan your call's transition produced (`world/effect-plan/v1`) |
| `effects` | The planned effects, in order. In this version a plan has zero or one effect. |
| `effects[].id` | The effect's id inside the plan (always `e1` today) |
| `effects[].status` | `ok`, `failed` or `denied` |
| `effects[].record` | Content address of the broker's effect record, or `null` |

Everything else in the object is the tool's own output, described per tool in the [Tool reference](./tools.md). Resolve any `sha256:` ref with `GET /v1/objects/<ref>?payload=true`; see [Provenance](./provenance.md).

## How a call runs

Each call runs four steps, under a per-episode lock (calls in the same episode run one at a time):

1. **Plan.** The tool's pure AILANG transition reads your arguments and returns a plan: either one effect, or zero effects plus a refusal.
2. **Execute.** The broker checks the planned effect against your grant (effect name, scope, expiry, remaining budget), records the intent, runs the handler inside your worktree, and records the outcome.
3. **Finish** (only `ailang-check` and `builtins-search`). A second pure run reshapes the handler output.
4. **Commit.** World appends one log entry whose record names your input, the plan, the effect record and the output.

## Plan refusals: zero effects, still committed

If your arguments are wrong, the plan refuses and no effect runs. The result is the refusal plus a `world` block with an empty `effects` list:

```json
{"ok": false,
 "refused": "unknown argument \"limit\"; ailang-read admits only: path",
 "world": {"effects": [], "plan": "sha256:…"}}
```

Nothing was read, written or run, and no budget was spent. The call is still in the log, because a bad call is provenance too. Fix the argument named in `refused` and call again.

Common plan refusals:

```text
unknown argument "<key>"; <tool> admits only: <keys>
missing required argument "<key>"
argument "<key>" must be a string
path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"
```

## Handler and policy refusals: the effect ran

When the plan is fine but the operation is not allowed or cannot succeed, the effect runs with status `ok` and its output says `ok: false`. These are AILANG's policy layer speaking, measured on the v0.51.0 tool binary:

```text
read /etc/hosts: path "/etc/hosts" escapes sandbox "<worktree>"
write .git/config: .git/ is read-only to the lane's tools
write .ailang/x.ail: matches fs_deny_write ".ailang/**" — read-only under this policy
edit hello.ail: old_text occurs 8 times; include more context so it is unique
edit hello.ail: old_text not found — the file may have changed; read it again
```

These calls cost one unit of budget, because the handler did run.

## Budget and grant denials

If the broker refuses the effect, the status is `denied`, and nothing executes. For tools without a finish phase the result holds only the `world` block:

```json
{"world": {"effects": [{"id": "e1", "status": "denied", "record": "sha256:…"}],
           "plan": "sha256:…"}}
```

For `ailang-check` and `builtins-search` the finish phase turns this into `{"ok": false, "status": "denied", "refused": "the Ailang.Check effect ended denied", …}`.

The reason is in the effect record's `denial` field: `denied:budget`, `denied:scope`, `denied:expired` or `denied:effect-name`. Budgets are counted in calls, persist across reconnects and daemon restarts, and are spent only by effects that executed; a denial spends nothing. When you see `denied:budget`, stop calling tools that need that effect and tell your operator. A status of `failed` means the handler itself failed (for example, an operation outside the effect's op allowlist); the record says so, and the result again holds only `world`.

## `ailang-run` outcomes

`ailang-run` returns `{admitted, exit_code, decision, limit, stdout, stderr}`. The handler recognises exactly four outcome shapes. Anything else is recorded as a `failed` effect, never guessed at.

**(a) Admitted.** The program's declared effects fit the policy, and it ran. `decision.ok` is `true`, `limit` is `null`, the program's output is in `stdout`, and `stderr` holds one `policy:` line. Recorded in the attended smoke:

```json
{"admitted": true, "exit_code": 0, "limit": null,
 "decision": {"allowed_caps": ["FS", "IO"], "declared_effects": ["IO"], "function": "main", "ok": true},
 "stdout": "hello from pi\n",
 "stderr": "policy: {\"ok\":true,…,\"security_mode\":\"restricted\",\"timeout_ms\":8000,…}\n",
 "world": {"effects": [{"id": "e1", "record": "sha256:f932f3aa…", "status": "ok"}], "plan": "sha256:5d0d3fde…"}}
```

**(b) Admitted, then refused at runtime.** The program type-checks and its effects are allowed, but an operation breaks confinement while it runs. `admitted` is `true` and `exit_code` is `1`. The reason is in `stderr`, before the `policy:` line. This program reads outside the worktree:

```ailang
module escape

import std/fs (readFile)
import std/io (println)

export func main() -> () ! {IO, FS} {
  println(readFile("/etc/hosts"))
}
```

```text
Error: execution failed: readFile: path "/etc/hosts" escapes sandbox "<worktree>"
policy: {"ok":true,…}
```

A refused inner effect does not always produce a non-zero exit code, so always read `stdout` and `stderr`.

**(c) Refused before execution.** The program never ran. `admitted` is `false`, and `decision` carries `error_kind`, `message` and, for a capability problem, `missing_from_policy`. Measured on v0.51.0:

```json
{"ok": false, "error_kind": "policy_violation", "message": "declared effects exceed policy",
 "function": "main", "declared_effects": ["IO", "Net"], "allowed_caps": ["FS", "IO"],
 "missing_from_policy": ["Net"]}
```

That is the `decision` for a program declaring `! {IO, Net}`, with exit code 2. Narrow the program's effects to what the policy allows (`IO`, `FS`). A wrong path gives `"error_kind": "read_failed"` with `"message": "open nope.ail: no such file or directory"` and exit code 1. A type error is also refused here (exit code 2); run `ailang-check` first to get structured diagnostics instead.

**(d) Supervisor limit.** The AILANG supervisor stopped the program, for example at the 8-second timeout. `exit_code` is `3`, `stdout` holds whatever was printed, and `limit` is the parsed `policy-result:` line, carrying `reason` (such as `timeout`) and `stage`. A program can print its own `policy-result:` line, so treat `limit` as a hint about your run, not as an authenticated fact.

The program's inner file operations are governed by the AILANG policy (sandbox root, deny list, `FS` budget of 1000 operations, 8 s timeout). They are not brokered or recorded one by one: World records the run as one `Ailang.Run` effect with its envelope.

## When the call itself fails

Some failures happen in the host, not in the tool. On MCP, every such failure arrives as the same JSON-RPC error envelope, as plain JSON rather than SSE:

```json
{"jsonrpc": "2.0", "id": 2, "error": {"code": -32603, "message": "host callback failed"}}
```

or `"host callback timed out"` when the deadline passed (`id` is `null` when the request was a batch). The message does not say which failure occurred, so assume the worst case: **the effect may have run.**

The worst case has a name in World: **effects unrecorded**. The broker already executed the effect and durably journaled its intent, outcome and record, but committing the invocation failed (a concurrent head move, a finish-phase failure, or a post-effect write error). The budget was debited exactly once. On A2A the error message says so explicitly and lists the record refs:

```text
effects were requested but the invocation was not confirmed committed; effect records: sha256:…
```

On MCP you only see `host callback failed`. Either way:

1. **Do not blindly retry** a write, an edit or a run. Read the worktree state first (`ailang-read` the file you were changing) and decide from what is actually there.
2. Reads, checks and searches are safe to repeat; they only cost budget.
3. Other A2A error messages tell you what to do directly: `invocation outcome is not confirmed; resend the same task id`, `invocation was not committed; send a new task id`, `world head moved during invocation; not committed; send a new task id`, `task id already used in this session`.

## Retries are at-least-once

- **MCP** has no idempotency. Every `tools/call` is a new invocation with a fresh random task id, even when your JSON-RPC id repeats. If you resend after a lost response, the tool runs again and both runs are recorded.
- **A2A** treats your task id as an idempotency key within the session. Resending a task id that committed returns the committed result without running the tool again.

If a response is lost, check the log before resending (`GET /v1/log?from=<n>`, see [Provenance](./provenance.md)). Your call is there if it committed.

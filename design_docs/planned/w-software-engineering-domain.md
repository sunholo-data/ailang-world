# w-software-engineering-domain — The first real transitions: AILANG coding tools over MCP/A2A, every effect through the broker

**Status**: PLANNED (attended draft, 2026-10-02). Not implemented; no quorum yet. Base `cbe564ce12f791d95fd184da86d837fb2b6bc80e` (V1).
**Clauses**: 3 (explicit authority end-to-end), 4 (resident-agent non-inferiority floor — this row supplies the World arm's tools), 5 (real recorded operation for provenance walks); north star "software-engineering domain live".
**Direction**: owner-approved (Mark, attended 2026-10-02): (1) pure transitions emit an effect plan, executed by the coordinator through the broker; (2) tool parity with the hardened `ailang_only` lane, reusing `ailang policy-tool` / `ailang run --policy` as broker handlers; (3) no Docker; (4) production handler registry + budgets + self-contained `.ail` package + checked-in publish manifest; (5) live smoke with pi and Claude Code. Refinements to the direction are marked **[REFINED]** with their V-row.
**Closes**: R-106-3 (effect bridge) by the effect-plan route; R-106-10 for effects; row 107 Residual 7 for production content (R-106-1).
**Estimate**: ~6.5 executor days (M1–M5, M7) + two attended sessions (M0, M6).
**Authoring note**: drafted by a read-only Plan agent in the attended session; every premise cites a V-row the drafter ran first-party (§3).

## 1. Why (measured)

World 1.0's north star requires "the software-engineering domain live" (V29). Today World serves no usable domain:

- **P1 — the live registry serves zero transitions.** Neither live store has a `world/transition-registry/v1` head (V2); no publish manifest is checked in, and no non-test code names a transition ID (V3); the only production publisher is the attended `world-publish transitions` verb (V4). Row 107 Residual 7 holds: a module importing `world/*` fails LDR001 at publish, so transitions must be self-contained (V31).
- **P2 — effectful transitions are refused.** The coordinator returns `EffectsUnsupportedError` for any descriptor with declared effects, and the capsule always runs `--caps ""` (V5). The bridge is named residual R-106-3; R-106-10 notes invocations do not debit budget (V31).
- **P3 — no handlers in the coordinator path.** `broker.OpenBinder` passes a nil handler registry (V6). Handlers exist for FS.Read/FS.Write (scope = one exact absolute file path), Git.Commit, Model.Infer, Human.\*, Registry.Publish — but no check/run/search effect (V9).
- **P4 — budgets reset on every invocation.** The broker does debit an in-memory ledger per effect, but the daemon opens a fresh session at full budget on every Dispatch (V7). R-106-10's real root is the non-persistent ledger, not a missing debit. **[REFINED]**
- **Consequence:** clause 4's World arm ("MCP transition tools only") has nothing to call, so the floor would fail trivially; clause 5 has no corpus of real operation (NO-CORPUS).

## 2. Goals and non-goals

**Goal.** A session holding the right grants can list and call, over `/mcp/` and `/a2a/`, the 8 capabilities of the `ailang_only` lane (V21). Each effect is executed by a broker handler, checked against capability, scope and budget, journaled as intent and outcome, and recorded as replay input. Each call commits one World log entry whose record links to its effect records.

**Non-goals (owners).**
- The clause-4 benchmark harness and floor run → row 93 (`w-prove-the-1-0-bar` Phase B).
- The row-114 protocol → row 114.
- motoko as a client → it refuses restricted policies and its MCP is stdio-only; optional row-93 arm.
- Network/GitHub effects, and Git.Commit as a tool → a later domain row (Git.Commit runs repo hooks/config, R-SE-5).
- Docker or any container → never in this row (direction 3).
- Git.CreateWorktree as an effect → a later row; the operator creates worktrees here.
- A multi-round agent loop inside a transition → v1 is one plan round plus one finish round.

## 3. Verification Log (all run 2026-10-02 at V1's base, read-only)

| V | Command | Observed (short) |
|---|---|---|
| V1 | `git log -1 --format='%H %cd' --date=short` | `cbe564ce12f7… 2026-10-02` |
| V2 | `sqlite3 -readonly "file:$db?mode=ro" "select count(*) from epoch_registry_heads where registry_name='world/transition-registry/v1'"` for `~/.ailang/world/world-0.1.1.db` and `world.db` | `0` and `0`; each store has 1 head total (`world/approvals/v1`) |
| V3 | `git ls-files '*.json' \| xargs grep -l transitionFnFile \| wc -l`; `grep -rnE '"[a-z0-9_-]+\.[a-z0-9_.-]+"' host cmd` (non-test, filtered ID/skill/transition) | `0`; no matches |
| V4 | `grep -rn 'PublishSet(' host cmd \| grep -v _test` | definition `transitionreg/publish.go:60`; sole caller `cmd/world-publish/transitions.go:113` |
| V5 | `sed -n 280,282p host/coordinator/coordinator.go`; `grep -n caps host/capsule/capsule.go` | `if len(d.DeclaredEffects) != 0 { // R8 … &EffectsUnsupportedError`; `caps: ""` (:178), `"--caps", r.caps` (:212) |
| V6 | `cat -n host/broker/binder.go`; `sed -n 548,551p host/daemon/daemon.go` | `:15-16 newSession(s, episodeID, grants, nil, Live, nil)`; the daemon's Binder closure calls `OpenBinder` per Dispatch (also `coordinator.go:263-264`) |
| V7 | `sed -n 114,127p;255,261p host/broker/broker.go` | grants copied at `newSession` (:124); `s.debitGrant(grantIndex, decision.Remaining)` in memory (:260); a new Session per Dispatch starts at full budget |
| V8 | `sed -n 99,101p host/broker/decide.go`; `sed -n 60,71p host/broker/confined.go`; `sed -n 214,216p host/broker/broker.go` | `scopeMatches: c.Scope == scope` (exact); BoundInvoker requires an exact `(effect, scope, cost)` triple; missing handler → `no handler registered for %q` |
| V9 | `grep -rhn 'Effect[A-Za-z]* *= *"' host/broker/*.go`; `sed -n 32,75p host/broker/handlers_fs.go` | FS.Read, FS.Write, Git.Commit, Model.Infer, Human.Approve, Human.PollApproval, Registry.Publish; FS scope = canonical absolute file path |
| V10 | `sed -n 176,186p;367,425p host/broker/broker.go` | Replay branch returns before the registry lookup; `invokeReplay` reads only the record and result objects — a handler is never called |
| V11 | `sed -n 1,117p host/coordinator/plan.go`; `sed -n 173,227p host/coordinator/coordinator.go` | `world/invocation-record/v1` = {invocationId, episodeId, skillId, transitionFn, interpreter, semanticsEpoch, input, output}, no effect refs; `committed()` decodes v1 |
| V12 | `grep -n 'invokeDeadline *=\|writeTimeout *=' host/daemon/daemon.go`; `sed -n 15,17p host/broker/handlers.go` | write 30 s, invoke 20 s; `handlerExecTimeout = 30s`; `MaxInput/MaxOutput: 1 << 20` (daemon.go:551) |
| V13 | `sed -n 16p host/projection/mcpname.go`; `sed -n 330,352p host/transitionreg/transitionreg.go`; QUICKSTART §8 | escapes only `_ . /`; `-` is a legal ID byte; "Retries create new tasks and may repeat effects" |
| V14 | `~/.pinned-ailang/ailang --version`; `ailang --version` | `AILANG v0.41.0 Commit 24ee108`; PATH = `v0.51.0-29-g078463551-dirty` (never used) |
| V15 | `echo '{"op":"summary"}' \| ~/.pinned-ailang/ailang policy-tool --policy <(…restricted, caps IO/FS, fs_sandbox=ailang-world…)` | `ok:true`; ops = agent_prompt, ai_check, axioms, builtins_list, builtins_show, check, devtools_prompt, docs_search, **edit**, examples_list, examples_search, examples_show, examples_tags, fmt, iface, pkg_docs, policy_check, prompt, **read**, summary, test, tree, version, **write**; no `run` |
| V16 | same, with `{"op":"read","path":"go.mod"}`, `"../ailang/go.mod"`, `"/etc/hosts"`, `{"op":"run"}` | content returned; `escapes sandbox … openat ../ailang/go.mod: path escapes from parent`; `/etc/hosts escapes sandbox`; `unknown op "run"` |
| V17 | same, with `check world/contracts.ail`, `examples_search foldl`, `builtins_list {json}`; `builtins list --json \| wc -c` | `argv:[check, world/contracts.ail]` "No errors found"; "Found 6 examples" (corpus built into the binary — ailang-world has no `examples/`); `count:348`; 75,994 bytes |
| V18 | `~/.pinned-ailang/ailang run --help \| grep policy`; `run --policy <(P) --caps IO,FS f.ail`; `run --policy <(P with caps []) examples/runnable/adt_tree.ail` | `--policy` present; `Error: --policy: --caps is not allowed with --policy` rc1; exit 2 + `{"error_kind":"policy_violation","missing_from_policy":["IO"]}`; from a non-root cwd `MOD010`. **Admitted run UNMEASURED**: the worker re-reads the policy *path*, so a process-substitution policy shows `allowed_caps: []`; needs a real file (M0) |
| V19 | `git tag --contains <c> \| sort -V \| head -1` in `ailang` for 6236deb24, 37ad0bc51, 70393ba17, b43ea30a8, 0df2eda17; `sed -n 290,305p changelogs/v0.32-current.md` | run --policy v0.39.0; supervised run v0.41.0; policytool v0.41.0; `--chdir`/`--request-file` v0.50.0; confined-git ceiling + `.git` case-fold v0.50.1: "on macOS's case-insensitive APFS, `.GIT/config` *is* `.git/config` and was writable" |
| V20 | `grep -n frozenCompilerVersion cmd/world-publish/fences.go`; `head -6 go.mod` | `:91 = "AILANG v0.41.0"`; Go module pin `ailang v0.47.2` (separate) |
| V21 | `sed -n 40,55p internal/executor/toolpolicy.go`; `sed -n 29,45p internal/executor/pi/toolnames.go`; `grep -n 'name: "' .pi/extensions/*.ts` (ailang repo) | ailang_only = AilangRead/Edit/Write/Check/Run, BuiltinsSearch, ExamplesSearch, AilangCLI (:53-54); `ailang-exec.ts` registers run/read/write/edit/cli; `ailang_check`/`builtins_search` live in `ailang-lsp-lite.ts` (the search filters `builtins list --json` client-side); `examples_search` in `examples-search.ts` scans `process.cwd()` |
| V22 | `grep -n localhost docs/docs/guides/agent-tool-policy.md` (ailang repo) | :50 restricted "has no localhost/private/metadata grant"; :18 every `--net-*` flag refused by name under `--policy` (doc evidence, not a run) |
| V23 | `cat ailang-multivac/config/policies/ailang-only-executor.toml` | `fs_sandbox = "${WORKSPACE}"` (expanded by execute-job, not by ailang), `timeout_ms = 120000`, `fs_deny_write` |
| V24 | `~/.pinned-ailang/ailang run --help \| grep -i 'bridge\|handler\|socket\|ipc'` | only `ai-stub`, `*-allow-localhost`, the internal `policy-worker` fd — no effect-bridge seam |
| V25 | `claude --version`/`--help`; `pi --help`; pi-mcp-adapter `package.json` + README; `codex mcp add --help` | Claude Code 2.1.287: `--tools ""` disables built-ins, `--mcp-config` takes JSON strings, `--strict-mcp-config`, `--allowedTools`; pi `--no-builtin-tools`, `--tools`, `--mcp-config`; pi-mcp-adapter 2.33.0 `directTools`, `auth:"bearer"`/`bearerTokenEnv`, `toolPrefix` default `"server"`; codex-cli 0.159.2 `--url --bearer-token-env-var` |
| V26 | `cd / && ~/.pinned-ailang/ailang iface std/json` | `CACHE_WRITE_FAILED … path=std/.ailang/cache` — the binary writes its compile cache relative to cwd (the probe left a gitignored `ailang-world/std/.ailang/cache`, since removed) |
| V27 | `/usr/bin/time -p` ×3 on a `policy-tool` read | real 0.19 / 0.09 / 0.03 s |
| V28 | `sed -n 51,60p;488,491p host/store/journal.go` | EffectIntent carries `Cost`; intents keyed `effect:<episode>:<ordinal>`, range-scanned per episode |
| V29 | `sed -n 15,20p;988,1010p design_docs/world-mission.md` | north star "with the software-engineering domain live"; clause 3 "every effect goes through the broker … results are recorded (replay input)"; clause 4 arms = "MCP transition tools only" |
| V30 | `sed -n 160,163p scripts/verify_ail.sh` | sweep roots `design_docs` and `world` only, with an exact LEG1 manifest |
| V31 | `grep -n 'R-106-3\|R-106-10\|R-106-1\b' design_docs/planned/w-transition-invocation-coordinator.md`; QUICKSTART §6 | R-106-3 bridge; R-106-10 budget; R-106-1 first useful skill owned as packages; "a module importing `world/*` is refused with `LDR001`" |

## 4. Design

### 4.1 Effect model: plan → execute → finish → commit (recommended over a live bridge)

| | Effect plan (recommended) | Live capsule→broker bridge |
|---|---|---|
| Purity / Z3 | transition stays `--caps ""`; plan laws can carry contracts | transition becomes effectful; effects mid-evaluation |
| Seam | exists: `Bound.Request` (V8) | none on the pin (V24); needs an upstream feature, or giving the capsule Net plus loopback, which restricted mode forbids (V22) |
| Replay | pure re-run + recorded results fed back, no handler (V10) | needs an interposed recorded-result channel inside the interpreter |
| Expressiveness | one plan round + one finish round; no data-dependent loops | arbitrary |

The bridge buys only data-dependent multi-step effects. A coding agent already supplies that by repeated tool calls, each one a committed transition — which is the right granularity for clause 5. The bridge stays residual R-106-3b (§9), routed upstream if a real need appears.

**Calling convention v2** — applies only to descriptors with non-empty `DeclaredEffects` and is derived from existing fields, so no registry interface change. The coordinator passes `main` one string: `{"phase":"plan","args":<MCP/A2A arguments>}` and, if the plan asks for it, `{"phase":"finish","args":…,"results":[EffectResult…]}`.

**Plan output (`world/effect-plan/v1`)**, a JSON object:
```json
{"plan":"world/effect-plan/v1",
 "effects":[{"id":"e1","effect":"Workspace.Read","scope":"worktree","cost":1,
             "payload":{"op":"read","path":"src/main.ail"}}],
 "finish":false,
 "result":null}
```
A refusal is a valid pure plan: `"effects":[]` with `"result":{"ok":false,"refused":"…"}`. It still commits — a bad call is provenance.

**Plan laws (L1–L6)**, enforced by Go `parsePlan`; per S1 the law is specified in a contracted `design_docs/sketches/effectplan.ail` (the `effectbroker.ail` pattern) with a drift test:
- L1: 0 ≤ |effects| ≤ 8.
- L2: ids unique, `[a-z0-9]{1,16}`.
- L3: each `(effect, scope, cost)` ∈ `DeclaredEffects` (BoundInvoker re-checks as defence in depth).
- L4: payload is a JSON object, canonically re-encoded, ≤ 1 MiB.
- L5: effects run in order; the first non-`ok` stops execution and later ones are `skipped`.
- L6: the finish output is a JSON object and cannot request effects (`effects` key refused); the key `world` is reserved in every transition and handler output.

**EffectResult**: `{"id","effect","status":"ok|failed|denied|skipped","record":"sha256:…"|null,"output":<handler JSON>|null}`.

**Output returned to the MCP/A2A caller** (= the committed output object = the world's `stateRoot`):
- One effect, no finish: the handler output's fields plus `"world":{"effects":[{"id","status","record"}],"plan":"sha256:…"}` — mirrors the lane's own result shape (V21) to minimise ergonomic drift.
- Several effects, no finish: `{"results":[…],"world":…}`.
- Finish: the finish output plus `"world":…`.
- The MCP projection already wraps `OutputBytes` as text plus `structuredContent`; no new wire.

**Replay.** `world/invocation-record/v2` = the v1 fields + `plan` (plan-object ref) + `effects` (ordered effect-record refs). Replay: (1) re-run the plan phase and require byte-equality with the plan object; (2) feed the same requests to a broker `Replay` session over the `effects` refs — it returns recorded results or denials and never dispatches (V10); (3) re-run finish, if any, on the reconstructed results; (4) recompose and require byte-equality with the output. Replay reconstructs World's record, **not** the worktree's files: write/edit payloads (full content) are in the request objects, and an `Ailang.Run` program's inner effects are captured only as the run envelope (D-SE-3).

**Ordering and commit**, under the per-episode lock (§4.4): Bind+Check → plan capsule run → `parsePlan` → effects via `Bound.Request` → optional finish capsule run → `planInvocation` (v2) → `AppendIntent` → `Commit`. On an R14 head conflict the coordinator re-reads the head and re-plans (pure, because transitions do not read world state, R-106-8), up to 3 times, and **never re-runs effects**. When retries run out it returns `EffectsUnrecordedError{InvocationID, EffectRecords}`; the effect intents and outcomes are already durable in the effect journal, so provenance survives. Under client retry, effectful calls are **at-least-once** (existing MCP behaviour, V13), and every execution is recorded. The R8 refusal remains for any declared effect the session's registry has no handler for, checked before any effect runs — no partial execution.

### 4.2 Tools — one transition each, parity with the `ailang_only` lane

| ID = MCP name | Args (lane-shaped) | Effect (scope `worktree`, cost 1) | Handler op | Finish |
|---|---|---|---|---|
| `ailang-read` | `{path}` | Workspace.Read | policy-tool `read` | no |
| `ailang-write` | `{path, content}` | Workspace.Write | policy-tool `write` | no |
| `ailang-edit` | `{path, old_text, new_text}` | Workspace.Write | policy-tool `edit` | no |
| `ailang-check` | `{path}` | Ailang.Check | policy-tool `check` | yes: structured `{code,message,file,line,col}` diagnostics (lsp-lite parity) |
| `ailang-run` | `{path, args_json?}` | Ailang.Run | `ailang run --policy P [--args-json …] <rel path>`, cwd = worktree root (MOD010, V18) | no |
| `builtins-search` | `{query?, module?}` | Ailang.Discover | policy-tool `builtins_list {json}` (76 KB, V17) | yes: substring filter (no op does it, V21) |
| `examples-search` | `{query, limit?}` | Ailang.Discover | policy-tool `examples_search` (corpus built into the binary, V17) | no |
| `ailang-cli` | `{op, path?, module?, query?, package?, flags?}` | Ailang.CLI | policy-tool, `op` ∈ the policy summary's `cli` list | no |

Each descriptor's `access` is its own effect triple at cost 0, so a session sees exactly the tools it holds grants for. Descriptions carry over the lane's tool descriptions with the confinement facts updated.

### 4.3 Handlers — reuse the hardened lane, no new Go confinement

`host/broker/handlers_ailang.go`: `AilangToolHandler{bin, binRef, policyPath, root}` implements `Handler`.
- Asserts `req.Scope == "worktree"`.
- **Fixed op allowlist per effect name**: Read→{read}; Write→{write, edit}; Check→{check}; Discover→{builtins_list, examples_search}; CLI→the summary's `cli`; Run→the run path only.
- Sends the payload to `policy-tool --policy P` on stdin via `runBounded` (process group, output cap, context deadline), `cmd.Dir` = the worktree.
- Result: `{"ok",…policy-tool Response fields…,"tool":"sha256:<binRef>","policy_digest":…}`; the run handler returns the lane's `composeEnvelope` shape (`admitted, exit_code, decision, limit, stdout, stderr`) transcribed in Go.
- All path/symlink confinement, `.git/` protection and widening-flag refusals stay **inside AILANG's Go policy layer** (V16, V18).

**Tool binary** (D-SE-1 A): `serve --tool-ailang-bin PATH`, archived and hash-verified exactly like `--ailang-bin`, independent of the `.ail` compiler pin (V20).

**Rendered policy**, per episode at `<db-dir>/policies/<episode>.toml` (mode 0600): `security_mode="restricted"`, `allowed_caps=["IO","FS"]`, `fs_sandbox=<root>`, `timeout_ms=15000` (below the 20 s invoke deadline, V12), `fs_deny_write=[".github/**",".pi/**",".claude/**"]`, `[budgets] FS=1000`, `entry="main"`. Startup refuses if the policy directory is inside the workspace root — the lane's D4 rule: an agent that can edit its own policy has no policy.

**Production registry.** `broker.OpenBinder(s, episodeID, grants, reg Registry)`. The daemon's Binder closure gets `reg` from `workspaceRegistry(episodeID)`: if `--workspace-root` and `--tool-ailang-bin` are both set, `episodeID` matches `^[a-z0-9][a-z0-9-]{0,63}$`, and `EvalSymlinks(root/episode)` is a directory inside `EvalSymlinks(root)` → a registry with the 6 effect names bound to that root and policy; otherwise an empty registry, so R8 refuses before anything runs. The raw `Session` stays private (TR.C).

**Workspace provisioning**: the operator runs `git worktree add <root>/<episode>` before `session mint`; creating worktrees as an effect is out of scope (§2).

### 4.4 Budgets (closes R-106-10 for effects)

A per-episode lock in the coordinator, held for the whole Dispatch of effectful descriptors (same single-process-writer premise as `inFlight`; same-episode calls serialise). Under the lock, a new `store.EffectSpend(ctx, episodeID) map[{effect,scope}]int64` sums `Cost` over the episode's effect intents (V28); the grants passed to `OpenBinder` are seeded with `Budget − spent`, floored at 0. The broker's frozen decision law is unchanged, and its `denied:budget` is recorded as today. Budget therefore persists across invocations, credentials and restarts; replay seeds from each record's `BudgetBefore`. Units are *calls* (cost 1), not wall time or tokens — stated on the card, not hidden.

### 4.5 The package and its publication

`packages/se-tools/`: 8 self-contained modules (`module se_tools/read` …; imports only `std/*`, which the binary embeds, V17/V26), each `export func main(input: string) -> string`. Every module carries a contracted path predicate (relative, non-empty, no `..` segment) — a precondition for a clear refusal, **not** the confinement (§4.3) — and named inline tests for its plan and finish outputs. The `.ail` is authored under S5 (`ailang prompt` / `ailang-docs` loaded, pinned `check`); **this doc claims no `.ail` text**, since S4 forbids unchecked snippets.

`packages/se-tools/transitions.json` is the checked-in manifest in `world-publish transitions` format (QUICKSTART §6: `id`, `title`, `description`, `transitionFnFile`, schemas, `access`, `declaredEffects`). `verify_ail.sh` gains a `.|packages/se-tools` root and LEG1 lines; if contracts add verified identities, M5 applies the S8 six-file recipe in the same commit. Publication is attended (TTY fence): Mark runs it.

### 4.6 Clients (M6 smoke)

- **curl**: QUICKSTART §8 shape — `tools/list` and each tool.
- **pi**: `pi --no-builtin-tools --tools <the 8 adapter-prefixed names> --mcp-config <file>`, server `{"url":"http://127.0.0.1:7644/mcp/","auth":"bearer","bearerTokenEnv":"WORLD_SESSION","directTools":true}`; the actual prefixed names are measured and recorded, not assumed (V25).
- **Claude Code**: `claude --tools "" --strict-mcp-config --mcp-config '{"mcpServers":{"world":{"type":"http","url":"http://127.0.0.1:7644/mcp/","headers":{"Authorization":"Bearer …"}}}}' --allowedTools mcp__world__ailang-read,…`. **[REFINED]** `--allowedTools` alone leaves the built-in tools in place (V25).
- **Trap**: a confined AILANG program can never be the World client — restricted mode has no loopback grant and refuses `--net-*` under `--policy` (V22). The client is always an outside harness.

## 5. Premises

| P | Premise | Evidence |
|---|---|---|
| P1 | No production transition is published or publishable today | V2, V3, V4, V31 |
| P2 | The coordinator refuses effects; the capsule is `--caps ""` | V5 |
| P3 | The production binder has no handlers; no code/test effect exists | V6, V9 |
| P4 | The budget ledger is per Dispatch, not per session | V6, V7 |
| P5 | Scope is exact equality and descriptors are global, so scope must be symbolic | V8, V13 |
| P6 | Broker replay never dispatches a handler | V10 |
| P7 | The record has no effect refs; reconcile decodes v1 only | V11 |
| P8 | The pin has `policy-tool` and the supervised `run --policy`; the security hardening is v0.50.1+ | V14–V19 |
| P9 | `policy-tool` confines paths in Go and refuses `run` | V16 |
| P10 | No bridge seam exists on the pin | V24 |
| P11 | A tool call must fit within 20 s | V12 |
| P12 | A confined program cannot reach World over loopback | V22 |
| P13 | All three reference clients speak HTTP MCP with a bearer token | V25 |
| P14 | The ailang binary writes a cache under cwd, i.e. inside the worktree | V26 |

## 6. Milestones (each ≤ 1 executor day; ACs mechanically checkable)

**M0 (attended, 0.5 d) — measurements and pin.** Mark answers D-SE-1..5.
- AC0.1: an admitted `run --policy` on a real policy file (outside the sandbox) records rc, the `policy:` line and stdout, for the chosen tool binary.
- AC0.2: the tool binary's `policy-tool` summary has the V15 op set, or the differences are listed.
- AC0.3: median of 5 capsule plan runs, recorded as a V-row for the R-SE-1 budget.

**M1 (1 d) — plan law.** `design_docs/sketches/effectplan.ail` with contracts for L1–L3, verified on the pin; `host/coordinator/effectplan.go` `parsePlan` (L1–L6) with a drift test against the sketch.
- AC1.1: table test, one row per law, each refusal typed.
- AC1.2: `ai-check` reports the sketch's identities verified and the gate manifest lists them.

**M2a (1 d) — coordinator execution.** Blanket R8 replaced by a registered-handlers check; execute via `Bound.Request`; compose EffectResults and output; finish phase; record v2; `committed()` decodes v1 and v2.
- AC2.1: ProbeHandler e2e — one effect → output contains the handler bytes plus a `world.effects[0].record` resolving to an `EffectRecordV1` with `Allowed`.
- AC2.2: an undeclared triple in a plan → no dispatch (counter 0), typed refusal.
- AC2.3: an unregistered declared effect → `EffectsUnsupportedError`, counter 0.

**M2b (1 d) — conflict and replay.**
- AC2.4: injected R14 on the first Commit → committed on retry, handler counter 1.
- AC2.5: retries exhausted → `EffectsUnrecordedError` naming the effect record refs, which `GetEffectReceipt` resolves.
- AC2.6 (AC-REPLAY-EFFECT): re-execute from the committed v2 record via a Replay session with a **nil registry** → output bytes equal, counter 0.

**M3 (0.75 d) — budget.** `store.EffectSpend`, the per-episode lock, seeding.
- AC3.1: grant `Workspace.Read=worktree:2`, three dispatches → ok, ok, and a recorded `denied:budget` on the third (a fresh Dispatch).
- AC3.2: the same across a store reopen.
- AC3.3: two concurrent same-episode dispatches with budget 1 → exactly one allowed.

**M4 (1.25 d) — handlers and wiring.** `handlers_ailang.go`, the archived tool binary, the rendered policy, the D4-style refusal, `serve --workspace-root/--tool-ailang-bin`, `OpenBinder(…, reg)`, `serve --help` text.
- AC4.1: real-binary test (gated on `AILANG_BIN`): read inside the worktree → content; `../x` → `refused`; `/etc/hosts` → `refused`.
- AC4.2: Workspace.Read with `op:"write"` → handler refusal, no file change.
- AC4.3: ailang-run of a module-rooted program → `admitted:true` and its stdout; a program declaring Net → `admitted:false`, `missing_from_policy`.
- AC4.4: episode `../x` or a symlink escaping the root → empty registry → R8.
- AC4.5: policy dir inside the root → startup error.

**M5 (1 d) — the package.** 8 modules, `transitions.json`, gate sweep.
- AC5.1: `verify_ail.sh` green with the new LEG1 lines; S8 inventory applied if the floor moves.
- AC5.2: `PublishSet` test-store publish of the manifest → 8 descriptors, IDs = MCP names, `EncodeMCPName(id) == id`.
- AC5.3: per tool, a daemon e2e over `/mcp/` with a fixture worktree and the real binary.

**M6 (attended, 0.5 d) — smoke.** QUICKSTART §9 executed verbatim: Mark publishes; mint with 6 grants; `serve`; curl `tools/list` (8) and one call each; pi and Claude Code each read, edit, check and run a program to the expected stdout. Evidence in `design_docs/verification/world-attended-<date>-se-domain/`, with log entries and record refs resolved via `/v1/log` and `/v1/objects`.

**M7 (0.5 d) — docs (S7).** QUICKSTART §9, `serve --help` flags, tool cards, client config examples. AC7.1: the runbook test binds the QUICKSTART flags to `serve --help`.

## 7. Load-bearing mutations (each must turn the named test red)

| Mutant (production) | Killing assertion |
|---|---|
| MUT-PLAN-UNDECLARED: `parsePlan` skips L3 and the BoundInvoker check is bypassed | AC2.2 dispatch counter = 0 |
| MUT-R8-PARTIAL: registered-handler check moved after the first effect | AC2.3 counter = 0 |
| MUT-RETRY-REEXEC: the conflict retry loop re-runs effects | AC2.4 counter = 1 |
| MUT-RECORD-V1: `planInvocation` drops the `effects` refs | AC2.6 replay equality + AC2.1 record resolves |
| MUT-REPLAY-LIVE: replay builds a Live session | AC2.6 nil-registry replay errors / counter = 0 |
| MUT-BUDGET-FRESH: seeding skipped (full budget per Dispatch) | AC3.1 third call denied |
| MUT-BUDGET-NOLOCK: per-episode lock removed | AC3.3 exactly one allowed |
| MUT-OP-ALLOWLIST: handler accepts any op for any effect | AC4.2 |
| MUT-SANDBOX-ROOT: `fs_sandbox` rendered as the workspace root, not root/episode | AC4.1 sibling-episode read refused |
| MUT-EPISODE-PATH: episode grammar / EvalSymlinks check removed | AC4.4 |
| MUT-POLICY-IN-ROOT: D4-style check removed | AC4.5 |
| MUT-WORLD-KEY: reserved `world` key collision allowed | AC1.1 L6 row |

## 8. Open decisions for Mark (default in bold)

- **D-SE-1 Tool binary.** **A: separate archived tool binary ≥ v0.50.1**, compiler pin unchanged. B: move the whole compiler pin (S8 inventory plus the world/core golden, V20). C: v0.41.0 (macOS `.GIT` write hole, V19).
- **D-SE-2 Names.** **A: hyphenated IDs (MCP name = ID).** B: `ailang_read` → encoded `ailang_uread` (V13).
- **D-SE-3 Run granularity.** **A: one brokered `Ailang.Run` per run; the program's FS is confined and budgeted by the AILANG policy and recorded as the envelope.** B: IO-only runs (benchmarks needing FS fail).
- **D-SE-4 Worktree binding.** **A: `--workspace-root`, worktree = root/episode.** B: a session-row field (schema change).
- **D-SE-5 Read calls commit.** **A: every call commits a log entry (the clause-5 corpus).** B: reads are journaled only.

## 9. Risks and residuals (named owners)

- **R-SE-1 Clause-4 overhead (row 93).** Each call costs a lock, 1–2 capsule runs, a subprocess (V27: 0.03–0.19 s) and a SQLite commit; same-episode calls serialise; `ailang-run` is capped below 20 s (V12) while the lane allows 120 s (V23). The World arm may fail the +25% wall-clock bound or time out long programs. M0 measures; row 93 tunes.
- **R-SE-2 Ergonomics (row 93).** No grep/glob/list tools beyond `ailang-cli tree`, and results carry `world` provenance tokens; the floor may fail on ergonomics rather than authority. Mitigation: lane-shaped outputs (§4.1).
- **R-SE-3 Arm parity (row 93).** Both arms must use the same ailang binary version; the `examples-search` corpus differs from the lane's checkout scan (V17, V21).
- **R-SE-4 Cache writes.** The binary writes `.ailang/cache` under cwd, inside the agent's worktree (V26). Owner: an upstream issue asking for a cache-dir flag; benign meanwhile (gitignored).
- **R-SE-5 Git.Commit** is not a tool here: it runs repo config and hooks in a worktree the agent writes. Owner: the later git-effects row, which must disable hooks/config as confined git does.
- **R-106-3b Live bridge** — only if multi-round effects are needed; upstream feature request. Owner: a future row.
- **R-SE-6 At-least-once under retry.** Owner: row 108's MCP retry-channel residual.
- **R-SE-7 Worktree state is not in World's content store** (files reconstructible only from Write payloads). Owner: a future snapshot-effect row.

## 10. Conflict surface

**Touches**: `host/coordinator` (R8 narrowed; record v2; retry; lock); `host/broker` (`OpenBinder` signature; new handler file); `host/store` (`EffectSpend`, read-only); `host/daemon` + `cmd/ailang-worldd` (2 flags); `scripts/verify_ail.sh` (roots/LEG1; S8 if the floor moves); `design_docs/sketches/effectplan.ail`, `packages/se-tools/`, `docs/QUICKSTART.md`.

**Does not touch**: the registry schema or interface hash (convention v2 derives from `DeclaredEffects`); the MCP/A2A wire (row 108's no-codec rule holds); the broker decision law (Z3 sketch unchanged); the world/core package and golden (under D-SE-1 A); `tools/launchd/*`, the fleet driver, the compiler (DESIGN §14). TR.C holds because dispatch stays `Bound.Request` → `BoundInvoker.Request`.

## 11. Axiom compliance

- **S1**: plan laws contracted in a sketch with a drift test; package path predicates carry contracts and inline tests.
- **S2**: every effect is in a host handler; transitions are pure.
- **S3** ("why is this not a package?"): the transitions *are* a package; handlers are the effect boundary itself; the coordinator change is host law that cannot live in a package.
- **S5**: no `.ail` claimed here; M1/M5 author it with the reference loaded.
- **S6**: §7 mutations, plus the R8 null case and the empty-registry refusal.
- **S7**: M7.

Related: `design_docs/planned/w-transition-invocation-coordinator.md`, `design_docs/implemented/w-mcp-dispatch-projection.md`, `design_docs/implemented/w-transition-registry-production-publisher.md`, `design_docs/implemented/w-effect-journal.md`, `design_docs/planned/w-prove-the-1-0-bar.md`, `design_docs/world-mission.md` (clauses 3–5).

# w-ailang-run-stdin-argv-caps — `ailang-run` gains stdin, argv and operator-allowlisted caps (Env, loopback Net) under the AILANG policy layer

**Status**: PLANNED — attended draft 2026-10-03 (drafted by a read-only Plan agent, written to file by the attended controller). Queue row 135. No quorum yet. D-135-1..4 RULED (D-WORLD-56): HELD on ailang#1557. Gated on row 134 landing on `dev`. Bases: ailang-world `dev` `926990d2a68b`; row-134 branch `efe3a9698a86`; ailang `76a5aef65c81`; tool binary v0.51.0 `b99dd25` (sha256 `e55ff71c710c1039…`).
**Clauses**: 4 (the World arm of the non-inferiority floor; counted as World tuning for row 93, D-WORLD-54) and 3 (each extra capability is explicit, brokered, budgeted and recorded).
**Direction**: D-WORLD-55, D-NF-3 = B (widen `ailang-run` before row 93's FINAL). Row 134's discipline applies throughout: a pure plan with one brokered effect, confinement kept in AILANG's policy layer, a rendered policy, the deadline budget, the outcome shapes, mutants, and V-rows. Refinements are marked **[REFINED]**.
**Estimate**: ≈ 2.75 executor days (M1–M3) + ≈ 0.5 d attended (M0, M4). That is above the row's original 1–1.5 d, because v0.51.0 cannot express Env or loopback Net in restricted mode (V8, V16).
**Probe conditions**: every V-row was run first-party on 2026-10-03 with the pinned tool binary `$A` under the handler's minimal env `H`. The worktree was a scratch `root/ep1`, with policy files outside it. Each probe program passed `$A check` first, and the AILANG reference was loaded first (S5). This doc claims no `.ail` text.

## 1. Why (measured)
- **P1 — `ailang-run` cannot run 4 of the 23 core gate tasks under benchmark conditions.** Its arguments are only `{path, args_json}`; it gets no stdin and no program argv; and its policy is restricted IO/FS (V2). The four tasks need (V5, V6):
  - `pipeline`: stdin;
  - `cli_args`: argv plus `Env`;
  - `api_call_json`: `Net` to the loopback mock `{{MOCK_HTTP_URL}}`;
  - `prompt_injection`: `Declassify`.
- **P2 — v0.51.0's policy layer only partly expresses this.**
  - **Works today:** stdin and argv after an inner `--` pass through `run --policy` (V12, V13).
  - **Refused in restricted mode:** `Env` (V8) and loopback Net, even when `net_allow` names it (V16).
  - **Not grantable in any mode:** `Declassify` is not a known capability (V9–V11).
  - **`trusted_host` admits Env and loopback Net** (V15, V17). But it drops `.git` write protection, the proxy refusal, Process confinement and the default byte limits (V19–V22). All of these except the proxy refusal can be put back in the rendered policy. The proxy refusal doesn't matter here, because the handler env sets no proxy (V19).
- **Consequence:** three of the four tasks can be unblocked honestly on v0.51.0. `prompt_injection` needs an upstream change (D-135-1).

## 2. Goals and non-goals
**Goals.**
- `ailang-run` gains three optional arguments: `stdin`, `argv` and `caps` (the run's exact capability set).
- Each Env or Net run is a distinct, declared, brokered and budgeted effect (`Ailang.RunEnv` / `Ailang.RunNet`). It also needs an operator allowlist at `serve`.
- Each run executes under a per-episode, per-cap-set policy variant. World renders it, AILANG's own `policy-tool summary` verifies it, and AILANG enforces it.
- Loopback is reachable only when the operator names the host.
- The World agent can test a solution under the grader's own caps, stdin and argv (§4.7).

**Non-goals (with owners).**
- `Declassify` admission: upstream (D-135-1).
- Port-scoped loopback: upstream (D-135-2).
- Env and Net together in one run: v1 refuses it (R-135-4).
- Public-network Net: a later row.
- Binary stdin: R-135-8.
- The row-93 harness: row 93.
- Any Go-side confinement, or any run without `--policy`: never.

## 3. Verification Log (all 2026-10-03)

| V | Command | Observed (short) |
|---|---|---|
| V1 | `git rev-parse dev origin/attended/row134`; `git -C ailang log -1`; `$A --version`; `shasum -a 256 $A` | `926990d2a68b`, `efe3a9698a86`; ailang `76a5aef65c81`; `AILANG v0.51.0 b99dd25`; `e55ff71c710c1039…` |
| V2 | `git show …row134:{packages/se-tools/se_tools/run.ail,host/broker/handlers_ailang.go,packages/se-tools/transitions.json}` | `argKeyAllowed` ensures `k == "path" \|\| k == "args_json"`; `executeRun` builds `run --policy P [--args-json J] -- path`, with no stdin and no program argv; `RenderEpisodePolicy` = restricted, `allowed_caps=["IO","FS"]`, `[budgets] FS=1000`; `loadSummary` refuses any mode other than restricted; one declared triple `(Ailang.Run, worktree, 1)`, `additionalProperties:false` |
| V3 | `git show …row134:host/coordinator/effectful.go \| sed -n 246,254p` | R8: every declared effect must have a registered handler before the plan runs |
| V4 | `git show …row134:host/coordinator/effectplan.go \| sed -n 30,95p;200,210p` | closed key sets `planKeys = [effects finish plan result]`, `effectKeys = [cost effect id payload scope]`; L4 `maxPayloadBytes = 1 << 20`; L3 requires `(effect, scope, cost) ∈ DeclaredEffects` |
| V5 | `cat ailang/benchmarks/{pipeline,cli_args,api_call_json,prompt_injection}.yml` | `pipeline` caps `[IO]`, `stdin: 1..5`, expected `2 4 6 8 10`. `cli_args` caps `[IO,FS,Env]`, `cli_args: ["numbers.txt"]`, expected `15`. `api_call_json` caps `[Net,IO]`, `net_allow_localhost: true`, expected `200`. `prompt_injection` caps `[Declassify,IO]`, expected `1` |
| V6 | `sed -n 228,330p ailang/internal/eval_harness/runner.go`; `…/agent_validation.go:163-185`; `grep -n 127.0.0.1 …/httpmock.go` | grade argv `run --entry main --quiet --relax-modules [--stdlib-path] --caps <spec.caps> [--net-allow-http --net-allow-localhost] benchmark/solution.ail [-- <cli_args>]`; stdin = `spec.Stdin`; `env = os.Environ()`; 10 s; the mock binds `127.0.0.1:0` |
| V7 | `git -C ailang show v0.51.0:internal/policy/policy.go` (tags); a policy with `net_allow_localhost = true` | the keys are `allowed_caps fs_sandbox net_allow net_allow_http process_allow cli_allow budgets timeout_ms max_source_bytes max_module_graph_bytes max_output_bytes max_fs_transfer_bytes ai_provider entry security_mode fs_deny_write`; decoding is strict: `unknown fields: [net_allow_localhost]`. There is **no env-allowlist key and no localhost key** |
| V8 | restricted `allowed_caps=["IO","FS","Env"]` | rc 1 `allowed_caps admits Env, which restricted mode has no confined adapter for — set security_mode = "trusted_host" … (restricted admits: AI,Clock,FS,IO,Net,Process,Rand,Stream)` |
| V9 | `allowed_caps=[…,"Declassify"]` in restricted and in trusted_host | both: rc 1 `allowed_caps names unknown capability "Declassify" (known: AI,Clock,Cog,DOM,Debug,Env,FS,IO,Msg,Net,Process,Secret,Stream,Trace)` |
| V10 | `decl.ail` (main `! {IO, Declassify}`) under the policy; `decl2.ail` (main `! {IO}`); grader-style `run --caps Declassify,IO decl.ail` | `decl2` check: `Missing effects: Declassify` (the effect propagates). `decl` under the policy: rc 2 `missing_from_policy:["Declassify"]`. Grader style: `1`, rc 0 |
| V11 | `git -C ailang show HEAD:internal/policy/resolve.go`; `git log v0.51.0..HEAD \| grep -i declassify` | unchanged at HEAD; no fix in flight |
| V12 | `printf '1\n2\n3\n4\n5\n' \| H $A run --policy <row-134 policy> -- stdin.ail`; the same with `</dev/null`; `v0.51.0:cmd/ailang/run_policy_supervise.go:72` | `2 4 6 8 10`, rc 0; with `/dev/null`: empty output, rc 0. Source: `cmd.Stdin = os.Stdin` |
| V13 | `argv.ail` (prints `getArgs()`): `-- argv.ail -- --caps=IO,Net --policy /etc/x - -- "a b"`; without the inner `--`; source `main_run.go:240-289`, `run_flag_order.go:19-36` | `ARGV=[--caps=IO,Net\|--policy\|/etc/x\|-\|--\|a b]`, rc 0. Without the inner `--`: rc 1 (a misplaced-flag refusal). Flag parsing stops at the path, so argv **never reaches the flag parser**; one leading `--` is stripped |
| V14 | `echo '"x"' \| H $A run --policy P --args-json - -- rd.ail` | `Error: empty stdin for -args-json -`: `--args-json -` reads stdin, and row 134's `json.Valid` guard keeps `-` out |
| V15 | trusted_host `[IO,FS,Env]`: `args.ail -- numbers.txt`; `env.ail` printing HOME / WORLD_SECRET / PATH with `WORLD_SECRET=leak` set; source `run_policy_supervise.go:228-231` | argv works; the program sees `HOME=<cache>`, `WORLD_SECRET=leak`, `PATH=/usr/bin:/bin`. trusted_host gives the worker `os.Environ()`, so the program's env equals the handler's child env |
| V16 | restricted `[IO,FS,Net]`, `net_allow=["127.0.0.1"]`, http on; POST to a scratch mock | `E_NET_IP_BLOCKED: localhost IP blocked: 127.0.0.1` (re-confirms row 134 V35) |
| V17 | trusted_host `[IO,Net]` with `net_allow` set to: `127.0.0.1` + http / `127.0.0.1` without http / `localhost` / `127.0.0.1:7655` / `127.0.0.2` | `200` (the mock saw the POST) / `E_NET_PROTOCOL_BLOCKED` / `disallowed 127.0.0.1` / `disallowed` / `E_NET_IP_BLOCKED 127.0.0.2`. Loopback opens only for a **literal** host in trusted_host; `host:port` is not a valid entry |
| V18 | the V17 `127.0.0.1` policy; URLs `:7656` and `:7644/v1/health`; `…row134:host/daemon/daemon.go:773-796` | `connect: connection refused` (the dial was attempted), so **the port is not scoped**. World's 8 GET routes are unauthenticated |
| V19 | `git -C ailang show v0.51.0:internal/runner/handlers.go:290-330`; `resolve.go` defaults; `run_policy.go:126` | restricted-only: `Net.RefuseProxy`, `Process.Confined`, `Env.ProtectGitDir`, byte-limit defaults, the env allowlist, the platform check. Both modes: the `fs_sandbox` root, `fs_deny_write`, `max_fs_transfer_bytes` when set, `[budgets]`, `timeout_ms` |
| V20 | trusted_host + the row-134 deny list; `wr.ail` writes `.git/config`, `.GIT/config`, the deny-listed dirs, `ok.txt`; reads with `..`/absolute paths; then `+ ".git", ".git/**"` | **`.git/config` and `.GIT/config` → written** (trusted_host drops `.git` protection); the deny globs → `E_FS_PROTECTED`; the reads → `escapes sandbox`. With `.git/**` added, `.GIT/config` and `.Git/hooks/pre-commit` are still written |
| V21 | trusted_host with the **8 case variants** `.git/** … .GIT/**` prepended; `.git/config`, `.GIT/config`, `.Git/hooks/pre-commit`, `.gIt/HEAD`, `.git`, `./.GIT/config`, `sub/../.GiT/config`, dotless-ı/dotted-İ forms, `ok2.txt`; source `internal/effects/fs_root.go:183-200` | every `.git` form → `E_FS_PROTECTED`; the ı/İ forms → `no such file` (not aliased); `ok2.txt` written; no hook created. Source: a `/**` pattern is a literal prefix, and any other pattern is `path.Match` |
| V22 | trusted_host admission line without, then with, explicit `max_*` byte limits | without: all `0` (unbounded); with: equal to the restricted defaults |
| V23 | restricted `["IO"]` vs `["IO","FS"]` with `fs_sandbox`; `link.ail → ../ep2/outside.ail` | IO only: **the outside program ran**. With FS: `program link.ail is outside fs_sandbox` (the root is only resolved when FS is admitted) |
| V24 | restricted `[IO,FS]`, `[budgets] FS = 0`: symlink, stdin, `rd.ail numbers.txt`, `overdecl.ail` | symlink refused; stdin works; read → `E_BUDGET_OPERATOR … limit=0`; an over-declared FS is admitted |
| V25 | grader-style `--caps Net,IO overdecl.ail`; trusted_host `[Env,FS,IO]` running `net.ail` | the grader admits an over-declared unused effect; the policy run → `missing_from_policy:["Net"]` (admission is stricter than the grader) |
| V26 | `module wrongname` at `benchmark/solution.ail`, with and without `--relax-modules` under the policy | both → `MOD010` (`--relax-modules` has no effect under `--policy`) |
| V27 | the §4.4 variants as rendered: FS-zero restricted (`stdin.ail`); Env variant (`sum.ail -- numbers.txt`); Net variant (`net.ail`); Net variant running `rd.ail`; Env variant running `net.ail`; `summary` of each | `2 4 6 8 10`; `15`; `200`; `E_BUDGET_OPERATOR`; `missing_from_policy:["Net"]`. Summaries: `restricted [FS,IO]`, `trusted_host [Env,FS,IO]`, `trusted_host [FS,IO,Net] ["127.0.0.1"]`, each with its own digest |
| V28 | `policy-tool summary` on a trusted_host variant | `ok:true`; reports mode, caps, `net_allow`, `fs_sandbox`, `timeout_ms`, so a variant can be verified by the policy layer itself |
| V29 | **(row 134, side finding)** the exact row-134 restricted policy: `wr.ail` writes `.CLAUDE/settings.json`, `.Claude/x.json`, `.GITHUB/workflows.yml`; policy-tool `write` `.CLAUDE/y.json`, `.AILANG/cache/z`; `diskutil info /` | the writes succeed and land in `.claude/` and `.github/` (APFS is case-insensitive); only the lowercase form is refused. Upstream case-folds only `.git` (`fs_root.go:151-153`, `policytool/fs_ops.go:60-62`). **Being fixed in row 134 before merge (R-SE-15).** |
| V30 | `git show …row134:host/broker/handlers.go:76-90,165-175`; `…:host/daemon/workspace.go:204-290` | `handlerCommand.stdin []byte` already exists; handlers are cached per episode; policy at `<db-dir>/policies/<ep>.toml` |
| V31 | row-93 design §1 P4, §4.2(a) | its text "`ailang-run` has no stdin, argv or extra caps" — this row makes it false |

## 4. Design
### 4.1 Tool schema (8 tools unchanged; only `ailang-run` changes)

| Key | Type / bound | Meaning |
|---|---|---|
| `path` | row-134 path predicate | unchanged |
| `args_json` | JSON string | unchanged (still refuses `-`, V14) |
| `stdin` | string, ≤ 65 536 UTF-8 bytes | the program's whole stdin; absent = `/dev/null` |
| `argv` | array, ≤ 32 strings, each ≤ 1 024 bytes, no NUL | program arguments, always passed after an inner `--` (V13) |
| `caps` | distinct names from `{IO, FS, Env, Net}` | the run's **exact** capability set, like the grader's `--caps` (V6); absent = row 134's `{IO, FS}` |

- **A leading `-` in argv is admitted [REFINED].** argv can't reach the flag parser (V13), so refusing it would only break programs that take flags. The injection is closed structurally, and MUT-NO-INNER-DASHDASH pins that.
- **`Declassify` in `caps`** gets a specific refusal naming the upstream gap. **Env together with Net** is refused in v1.
- **Effect selection** (L1 = 1 still holds) is a contracted pure `runEffect(hasEnv, hasNet)`: Env → `Ailang.RunEnv`, Net → `Ailang.RunNet`, otherwise `Ailang.Run`. All three use scope `worktree`, cost 1.
- **Payload** is `{path, args_json?, stdin?, argv?, caps?}`, with `caps` canonicalized (sorted, unique). With none of the new args, the plan bytes are **byte-identical to row 134's**.

### 4.2 Plan-law impact
- **M1's closed key sets are unchanged** (V4). The new fields go inside `payload`, which stays within L4.
- **L3:** `ailang-run`'s `declaredEffects` grows from 1 to 3 triples. `access` stays `(Ailang.Run, worktree, 0)`, so `tools/list` is unchanged at 8 names.
- **L-ARGS:** `argKeyAllowed` now admits `{path, args_json, stdin, argv, caps}`. New contracted helpers: `capAllowed`, `runEffect`, and bounds predicates.
- **R8 (V3):** the workspace registry **always** binds `Ailang.RunEnv` and `Ailang.RunNet` (6 → 8 effect names). The operator allowlist is enforced inside the handler, never by leaving a name unregistered.

### 4.3 Three gates on every extra capability
1. **Broker (who):** a session grant `Ailang.RunEnv=worktree:N` or `Ailang.RunNet=worktree:N`. Without one, the call is recorded as `denied` and nothing is debited (134:V55).
2. **Operator (what is possible):** `serve --run-allow-caps Env,Net`. A capability outside that list gets a typed handler refusal. Startup refuses unknown names, `Declassify` (V9), `Net` without `--run-net-allow`, and `--run-net-allow` without `Net`.
3. **Policy layer (what runs):** the variant's `allowed_caps` is exactly the requested set (plus FS at budget 0 when FS wasn't requested, §4.4), **verified by `policy-tool summary`** (V28) and enforced by AILANG.

The effect record carries the effect name, the request (stdin, argv, caps) and the output, which now adds `policy: {digest, security_mode, caps, net_allow}`.

### 4.4 Policy rendering: per episode, per cap set, lazily, cached
`RenderRunVariant(root, caps, op)` writes `<db-dir>/policies/<ep>.run-<caps>.toml` (mode 0600, outside the root), verifies it once with `summary`, and caches it with the handler (V30). A `{FS, IO}` request maps to row 134's base policy.

| Requested caps | Mode | Rendered additions (on top of row 134's key set) | Evidence |
|---|---|---|---|
| ⊆ {IO, FS} | restricted | FS not requested → `FS` still allowed with `[budgets] FS = 0`, which keeps the entry-inside-sandbox check and makes FS fail at use, as the grader would (V23, V24) | V12, V24, V27 |
| contains Env | **trusted_host** | the 8 `.git/**` case variants prepended to the deny list; the 4 byte limits set explicitly to restricted defaults; FS budget 1000, or 0 if not requested | V15, V20–V22, V27 |
| contains Net | **trusted_host** | as above, plus `net_allow = <--run-net-allow hosts>` and `net_allow_http = <--run-net-allow-http>` | V17, V27 |

**What trusted_host still relaxes (stated, not hidden; V19):**
- **The proxy refusal is off.** Moot: the handler env sets no proxy.
- **Process confinement is off.** Moot: Process is never rendered.
- **The worker sees the parent env.** That is the handler's minimal `childEnv` (V15): `HOME`, `PATH`, `LANG`, `LC_ALL`, `AILANG_CACHE_DIR`, and optionally `AILANG_EXAMPLES` (R-135-5).

**The deny-list case variants from R-SE-15 (row 134) apply to every variant.**

### 4.5 Handler (`executeRun`)
- **Accepts** `Ailang.Run`, `Ailang.RunEnv` and `Ailang.RunNet`.
- **Refuses:** unknown or case-aliased payload keys; a `caps` set that doesn't match the effect name; argv or stdin over their bounds; and a NUL in argv.
- **Executes** `run --policy <variant> [--args-json J] -- <path> [-- <argv…>]`. stdin goes in through the existing `handlerCommand.stdin` plumbing (V30). No widening flag is ever passed.
- **Unchanged from row 134:** the outcome shapes, the deadline (8 s policy, 10 s handler) and `childEnv`.

### 4.6 Operator configuration (`serve`)

| Flag | Default | Meaning |
|---|---|---|
| `--run-allow-caps Env,Net` | empty | the extra capabilities a run may request |
| `--run-net-allow HOST` (repeatable) | none | literal hosts rendered as `net_allow`; a loopback literal is the **only** way to open loopback. Host-level only: **the port cannot be scoped in v0.51.0** (V17, V18; D-135-2) |
| `--run-net-allow-http` | off | renders `net_allow_http = true` (the benchmark mock is plain http) |

### 4.7 Grader parity (row 93 V9)

| Grader (V6) | World `ailang-run` |
|---|---|
| `--caps <spec.caps>` | `caps: spec.caps` |
| `[-- cli_args]` | `argv: spec.cli_args` |
| stdin piped | `stdin: spec.stdin` |
| `--net-allow-http --net-allow-localhost` | operator `--run-allow-caps Net --run-net-allow 127.0.0.1 --run-net-allow-http` |
| `--relax-modules` | not available under `--policy` (V26) |
| `env = os.Environ()` | the handler's `childEnv` |
| timeout 10 s | 8 s policy / 10 s handler |

Parity runs **one way, by design: World is never more permissive than the grader.** An over-declared effect is refused at admission (V25), except FS, which is admitted at budget 0 (V24). A module-name mismatch is refused (V26).

Per task:
- `pipeline` → `{path, caps:["IO"], stdin:"1\n2\n3\n4\n5\n"}`;
- `cli_args` → `{path, caps:["Env","FS","IO"], argv:["numbers.txt"]}` (`Ailang.RunEnv`);
- `api_call_json` → `{path, caps:["IO","Net"]}` (`Ailang.RunNet`);
- `prompt_injection` → the typed Declassify refusal (D-135-1).

### 4.8 Replay
stdin, argv and caps are in the content-addressed request object, so Replay returns the recorded output and never re-executes (134:V10, V46). The variant bytes can be re-derived, and `policy.digest` detects drift.

## 5. Premises

| P | Premise | Evidence |
|---|---|---|
| P1 | `ailang-run` has no stdin/argv/caps today; restricted IO/FS only | V2, V31 |
| P2 | stdin and argv pass through `run --policy`; argv cannot reach the flag parser | V12, V13 |
| P3 | Env and loopback Net need trusted_host | V7, V8, V16, V17 |
| P4 | Declassify is not expressible in any v0.51.0 policy (nor at HEAD) | V9–V11 |
| P5 | trusted_host's relaxations are compensated in the rendered policy, or moot | V19–V22, V27 |
| P6 | Loopback is host-level; ports can't be scoped; World's GET routes are reachable | V17, V18 |
| P7 | FS at budget 0 keeps the entry check and grader-equivalent FS failure | V23, V24 |
| P8 | Plan laws need no change; R8 forces registering all three run effects | V3, V4 |
| P9 | A variant can be verified with `policy-tool summary` | V28 |
| P10 | World's run is never more permissive than the grader | V6, V25, V26 |
| P11 | The handler plumbing already carries stdin | V30 |

## 6. Milestones

**M0 (attended, 0.25 d).**
- AC0.1: re-run V12–V27 against the archived tool binary's hash (R-SE-8).
- AC0.2: file three upstream asks and record their ids:
  - Declassify as an admission-only capability that restricted mode can admit;
  - port-qualified `net_allow`, with restricted-mode loopback for a port-qualified literal;
  - case-folded `fs_deny_write` matching (V29).

**M1 (1 d) — transition and descriptor.**
- AC1.1: `run.ail` named tests pin the exact plan bytes for each case:
  - no new args (byte-equal to row 134's);
  - stdin;
  - argv;
  - Env → `RunEnv`;
  - Net → `RunNet`;
  - Env+Net → refused;
  - Declassify → the specific refusal;
  - an unknown cap, a duplicate cap, and `caps` given as a string;
  - argv over bounds: 33 items, a 1 025-byte item, a NUL;
  - stdin of 65 537 bytes.
- AC1.2: the contracts verify; `verify_ail.sh` is green.
- AC1.3: `ailang-run` declares 3 triples and has 5 properties with `additionalProperties:false`; `PublishSet` still yields 8 descriptors and `tools/list` still lists 8.

**M2 (1 d) — handler and variants.**
- AC2.1: the exact key set per variant class.
- AC2.2: real binary: stdin → `2 4 6 8 10`; `RunEnv` + `argv:["numbers.txt"]` → `15`; `RunNet` against an in-test 127.0.0.1 mock → `200`, and the mock saw exactly one POST with `X-Test-Header: value123`.
- AC2.3: each of these is refused with a stub exec counter of 0:
  - Env not allowlisted;
  - Net without `--run-net-allow`;
  - an effect/caps mismatch;
  - an unknown or aliased key;
  - a NUL in argv;
  - oversize stdin;
  - `args_json:"-"`.
- AC2.4: the trusted_host confinement matrix on both topologies (V20/V21, plus the R-SE-15 case variants). Every `.git` form and the worktree pointer file → `E_FS_PROTECTED`, bytes unchanged; a sibling read → `escapes sandbox`.
- AC2.5: a Net program under the base, FS-zero and Env variants → `missing_from_policy:["Net"]`; under the Net variant, `http://localhost:<port>/` → `DisallowedHost`.
- AC2.6: an Env program printing `getEnvOr("WORLD_SESSION","<none>")` → `<none>`.
- AC2.7: `argv:["--caps=IO,Net","-","--","a b"]` comes back verbatim.
- AC2.8: `caps:["IO"]` with a symlinked entry → `outside fs_sandbox`; `caps:["IO","Net"]` reading a file → `E_BUDGET_OPERATOR`.

**M3 (0.75 d) — wiring and docs.**
- AC3.1: the startup-refusal table; `serve --help` stays bound to QUICKSTART.
- AC3.2: a `/mcp/` end-to-end run for each of the three runnable tasks. Each commits; its `world.effects[0].record` names the right effect; its request object holds the stdin, argv and caps.
- AC3.3: with no `RunEnv` grant → `denied` and no exec; with budget 1 → the second call is `denied:budget`.
- AC3.4: a nil-registry replay is byte-equal.
- AC3.5: QUICKSTART §9 is updated.

**M4 (attended, 0.5 d).** Mark republishes. A curl client then writes and runs each task's reference solution with the §4.7 mapping:
- `pipeline`, `cli_args` and `api_call_json` print their expected stdout;
- `prompt_injection` gets the typed refusal, and `ailang-check` passes its contracts;
- Net while the daemon runs without Net → refused.

## 7. Load-bearing mutations

| Mutant | Killer |
|---|---|
| MUT-NO-INNER-DASHDASH: argv appended without the inner `--` | AC2.7 |
| MUT-ARGS-JSON-STDIN: the `json.Valid` guard dropped | AC2.3 `args_json:"-"` |
| MUT-CAPS-WIDEN: the variant renders the operator allowlist instead of the requested caps | AC2.5 |
| MUT-ALLOWLIST-SKIP: the handler ignores `--run-allow-caps` | AC2.3 |
| MUT-EFFECT-SELECT: the plan emits `Ailang.Run` for an Env request | AC1.1 + AC3.3 |
| MUT-GIT-VARIANTS: only `.git/**` rendered on trusted_host | AC2.4 `.GIT/config` |
| MUT-LIMITS-UNSET: the byte limits omitted | AC2.1 |
| MUT-FS-DROPPED: FS omitted when not requested | AC2.8 symlink |
| MUT-FS-BUDGET-OPEN: FS at 1000 when not requested | AC2.8 `E_BUDGET_OPERATOR` |
| MUT-NET-RESTRICTED: the Net variant rendered restricted | AC2.2 `200` |
| MUT-NET-WIDEN: `localhost` added beside the named host | AC2.5 |
| MUT-ENV-PARENT: `os.Environ()` passed to the run | AC2.6 |
| MUT-STDIN-UNRECORDED: stdin not taken from the payload | AC3.2 + AC3.4 |
| MUT-STDIN-CAP: the size bounds removed | AC1.1 + AC2.3 |
| MUT-VARIANT-UNVERIFIED: the summary check skipped | AC2.1 (a tampered variant is refused) |

## 8. Decisions — RULED by Mark, attended 2026-10-03 (D-WORLD-56)

**D-135-1 = B: HOLD the row for upstream Declassify (ailang#1557).** D-135-2 = A, with loopback port scoping asked upstream (ailang#1558). D-135-3 = A. D-135-4 = A. The original options follow.

- **D-135-1 Declassify.**
  - **A: ship stdin, argv, Env and Net now.** `prompt_injection` stays check-only in World and is reported per task in row 93; the upstream ask is filed in M0.
  - B: hold row 135 until upstream ships Declassify.
- **D-135-2 The loopback port can't be scoped.**
  - **A: accept host-level `127.0.0.1` behind the explicit operator flag** (R-135-2), plus a row-93 audit AC: no committed solution targets `:7644` or `/v1/`.
  - B: hold Net until upstream supports ports.
- **D-135-3 Effect naming.**
  - **A: per-cap effects `Ailang.RunEnv` / `Ailang.RunNet`.**
  - B: one `Ailang.Run`, with the caps in its payload.
- **D-135-4 trusted_host.**
  - **A: use it only for Env and Net runs, with the measured compensations.**
  - B: don't use it, and `cli_args` and `api_call_json` stay unrunnable.

## 9. Risks and residuals
- **R-135-1 trusted_host is weaker by construction.** Its safety rests on the rendered compensations and on V19's list of what it relaxes. Owner: the R-SE-8 re-proof, plus AC2.1/AC2.4 on every binary change.
- **R-135-2 Loopback is unscoped.** A Net run can reach any 127.0.0.1 port, including World's unauthenticated GET routes, so in row 93 a later trial could read an earlier trial's solution. Owner: the upstream ask; meanwhile, the D-135-2 audit AC.
- **R-135-3 Declassify is unavailable.** Owner: upstream.
- **R-135-4 No Env+Net combination in one run.** Owner: a later row.
- **R-135-5 Env runs see the handler's env.** It holds state-dir paths only, no secrets.
- **R-135-6 Parity runs one way only.** Owner: row 93's per-task report.
- **R-135-7 The entry check depends on FS being admitted.** Closed here by always rendering FS (at budget 0 when not requested).
- **R-135-8 stdin is text-only.**
- **R-SE-15 (row 134): case-sensitive `fs_deny_write` on APFS** (V29). Being fixed in row 134 before merge; the upstream ask is filed in M0.

## 10. Conflict surface
**Touches:**
- `packages/se-tools/se_tools/run.ail` and `transitions.json` (republished attended);
- `host/broker/handlers_ailang.go` and its tests;
- `host/daemon/workspace.go` (6 → 8 effect names);
- `cmd/ailang-worldd/main.go` (3 flags);
- `docs/QUICKSTART.md` §9;
- `scripts/verify_ail.sh`;
- optionally `tools/attended/se_smoke.sh`;
- the row-93 design's §4.2(a) and P4, plus the D-135-2 audit AC.

**Does not touch:** `effectplan.go`, the broker decision law, the projection, store and coordinator, or the compiler pin.

## 11. Non-goals (restated)
- Declassify admission;
- port-scoped loopback;
- Env and Net in one run;
- public-network Net;
- binary stdin;
- any Go-side confinement or any run without `--policy`;
- the row-93 harness.

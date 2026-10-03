# w-ailang-run-stdin-argv-caps — `ailang-run` gains stdin, argv and operator-allowlisted caps (Env, port-scoped loopback Net, Declassify) under the AILANG policy layer

**Status**: PLANNED. Attended draft 2026-10-03, drafted by a read-only Plan agent and written to file by the attended controller. Queue row 135; no quorum yet.
- **Decisions:** D-135-1..4 RULED (D-WORLD-56).
- **UN-HELD 2026-10-03.** ailang#1557 (Declassify) shipped in AILANG v0.52.1, together with #1558 (port-qualified `net_allow`) and #1559 (case-folded deny matching), as re-measured in V32–V42.
- **Revised on v0.52.1 [REFINED], and narrower than the ruling permits:**
  - Net is now **restricted**, scoped to one operator-named loopback port, with no `trusted_host`.
  - `trusted_host` is kept **only for Env**.
  - Declassify is admitted in restricted mode.
  - All four tasks are runnable.
- **Gate:** row 134 is landed on `dev`.
- **Bases:**
  - ailang-world `origin/dev` `54a9750`, after row 134 `f3a90ee`;
  - ailang tag v0.52.1 `c68ded4b2`;
  - tool binary v0.52.1, sha256 `0dd70a1d0036…` (V32).
  - The original draft's bases (V1: v0.51.0 `b99dd25`) are kept for the V1–V31 rows.

**Clauses**: 4 (the World arm of the non-inferiority floor; counted as World tuning for row 93, D-WORLD-54) and 3 (each extra capability is explicit, brokered, budgeted and recorded).

**Direction**: D-WORLD-55, D-NF-3 = B (widen `ailang-run` before row 93's FINAL). Row 134's discipline applies throughout: a pure plan with one brokered effect, confinement kept in AILANG's policy layer, a rendered policy, the deadline budget, the outcome shapes, mutants and V-rows. Refinements are marked **[REFINED]**.

**Estimate**: ≈ 2.5 executor days (M1–M3) plus ≈ 0.5 d attended (M4). M0 (the pin bump) is done in this change. The estimate is lower than the v0.51.0 draft's because Net no longer needs the trusted_host compensations.

**Probe conditions**:
- V1–V31 were run first-party on 2026-10-03 with the v0.51.0 tool binary.
- V32–V42 were re-run on v0.52.1 by the row-135 M0 executor.
- Both sets ran under the handler's minimal env `H` (`HOME`/`AILANG_CACHE_DIR` set to a cache dir outside the root, `PATH=/usr/bin:/bin`, `LANG=C`, `LC_ALL=C`).
- The worktree was a scratch `root/ep1`, and the policy files lived outside it.
- Each v0.52.1 variant is row 134's **exact rendered policy** (`RenderEpisodePolicy`) with only `security_mode`, `allowed_caps`, `net_allow*` or `[budgets]` changed.
- Each probe program passed `$A check` first, and the AILANG reference was loaded first (`$A prompt`, S5).
- This doc claims no `.ail` text.

## 1. Why (measured)
- **P1 — `ailang-run` cannot run 4 of the 23 core gate tasks under benchmark conditions.** Its arguments are only `{path, args_json}`. It gets no stdin and no program argv, and its policy is restricted IO/FS (V2). The four tasks need (V5, V6):
  - `pipeline`: stdin;
  - `cli_args`: argv plus `Env`;
  - `api_call_json`: `Net` to the loopback mock `{{MOCK_HTTP_URL}}`;
  - `prompt_injection`: `Declassify`.
- **P2 — v0.52.1's policy layer expresses three of the four needs in restricted mode [REFINED].**
  - **Works, unchanged:** stdin, and argv after an inner `--`, both pass through `run --policy` (V12, V13, V39).
  - **Now admitted in restricted mode:**
    - `Declassify` (V33; it was not grantable at all on v0.51.0, V9–V11);
    - loopback Net to **one named port** via a port-qualified literal `net_allow = ["127.0.0.1:PORT"]`. Every other port, host and name is refused (V34). On v0.51.0 this needed `trusted_host` and could not scope the port (V16–V18).
  - **Still trusted_host-only:** `Env` (V36). `trusted_host` still drops `.git` write protection and the default byte limits (V37, V38). Both can be put back in the rendered policy. The folded matcher means one `.git/**` entry now covers every case variant and the worktree pointer file (V37).
- **Consequence:** all four tasks can now be unblocked honestly. Only `cli_args` needs `trusted_host`, with its measured compensations.

## 2. Goals and non-goals
**Goals.**
- `ailang-run` gains three optional arguments: `stdin`, `argv` and `caps` (the run's exact capability set).
- Each Env or Net run is a distinct, declared, brokered and budgeted effect (`Ailang.RunEnv` / `Ailang.RunNet`), and needs an operator allowlist at `serve`. A Declassify-only run stays `Ailang.Run` (D-135-5, §8).
- Each run executes under a per-episode, per-cap-set policy variant. World renders it, AILANG's own `policy-tool summary` verifies it, and AILANG enforces it.
- Loopback is reachable only on the **host:port** pairs the operator names. The pair is a literal IP and a port; a bare host is refused (V34, V35) **[REFINED]**.
- The World agent can test a solution under the grader's own caps, stdin and argv (§4.7).

**Non-goals (with owners).**
- Env and Net together in one run: v1 refuses it (R-135-4).
- Public-network Net: a later row.
- Binary stdin: R-135-8.
- The row-93 harness: row 93.
- Any Go-side confinement, or any run without `--policy`: never.

## 3. Verification Log

### 3.1 v0.51.0 draft (2026-10-03, tool binary v0.51.0 `b99dd25`)

Rows marked *superseded* stay as the record of what v0.51.0 did.

| V | Command | Observed (short) |
|---|---|---|
| V1 | `git rev-parse dev origin/attended/row134`; `git -C ailang log -1`; `$A --version`; `shasum -a 256 $A` | `926990d2a68b`, `efe3a9698a86`; ailang `76a5aef65c81`; `AILANG v0.51.0 b99dd25`; `e55ff71c710c1039…` |
| V2 | `git show …row134:{packages/se-tools/se_tools/run.ail,host/broker/handlers_ailang.go,packages/se-tools/transitions.json}` | `argKeyAllowed` ensures `k == "path" \|\| k == "args_json"`; `executeRun` builds `run --policy P [--args-json J] -- path`, with no stdin and no program argv; `RenderEpisodePolicy` = restricted, `allowed_caps=["IO","FS"]`, `[budgets] FS=1000`; `loadSummary` refuses any mode other than restricted; one declared triple `(Ailang.Run, worktree, 1)`, `additionalProperties:false` |
| V3 | `git show …row134:host/coordinator/effectful.go \| sed -n 246,254p` | R8: every declared effect must have a registered handler before the plan runs |
| V4 | `git show …row134:host/coordinator/effectplan.go \| sed -n 30,95p;200,210p` | closed key sets `planKeys = [effects finish plan result]`, `effectKeys = [cost effect id payload scope]`; L4 `maxPayloadBytes = 1 << 20`; L3 requires `(effect, scope, cost) ∈ DeclaredEffects` |
| V5 | `cat ailang/benchmarks/{pipeline,cli_args,api_call_json,prompt_injection}.yml` | `pipeline` caps `[IO]`, `stdin: 1..5`, expected `2 4 6 8 10`. `cli_args` caps `[IO,FS,Env]`, `cli_args: ["numbers.txt"]`, expected `15`. `api_call_json` caps `[Net,IO]`, `net_allow_localhost: true`, expected `200`. `prompt_injection` caps `[Declassify,IO]`, expected `1` |
| V6 | `sed -n 228,330p ailang/internal/eval_harness/runner.go`; `…/agent_validation.go:163-185`; `grep -n 127.0.0.1 …/httpmock.go` | grade argv `run --entry main --quiet --relax-modules [--stdlib-path] --caps <spec.caps> [--net-allow-http --net-allow-localhost] benchmark/solution.ail [-- <cli_args>]`; stdin = `spec.Stdin`; `env = os.Environ()`; 10 s; the mock binds `127.0.0.1:0` (an ephemeral port) |
| V7 | `git -C ailang show v0.51.0:internal/policy/policy.go` (tags); a policy with `net_allow_localhost = true` | the keys are `allowed_caps fs_sandbox net_allow net_allow_http process_allow cli_allow budgets timeout_ms max_source_bytes max_module_graph_bytes max_output_bytes max_fs_transfer_bytes ai_provider entry security_mode fs_deny_write`; decoding is strict: `unknown fields: [net_allow_localhost]`. There is **no env-allowlist key and no localhost key** (unchanged on v0.52.1, V40) |
| V8 | restricted `allowed_caps=["IO","FS","Env"]` | rc 1 `allowed_caps admits Env, which restricted mode has no confined adapter for — set security_mode = "trusted_host" … (restricted admits: AI,Clock,FS,IO,Net,Process,Rand,Stream)` (re-confirmed on v0.52.1, V36) |
| V9 | *superseded by V33.* `allowed_caps=[…,"Declassify"]` in restricted and in trusted_host | both: rc 1 `allowed_caps names unknown capability "Declassify" (known: AI,Clock,Cog,DOM,Debug,Env,FS,IO,Msg,Net,Process,Secret,Stream,Trace)` |
| V10 | `decl.ail` (main `! {IO, Declassify}`) under the policy; `decl2.ail` (main `! {IO}`); grader-style `run --caps Declassify,IO decl.ail` | `decl2` check: `Missing effects: Declassify` (the effect propagates). `decl` under the policy: rc 2 `missing_from_policy:["Declassify"]`. Grader style: `1`, rc 0 |
| V11 | *superseded by V32/V33.* `git -C ailang show HEAD:internal/policy/resolve.go`; `git log v0.51.0..HEAD \| grep -i declassify` | unchanged at HEAD; no fix in flight |
| V12 | `printf '1\n2\n3\n4\n5\n' \| H $A run --policy <row-134 policy> -- stdin.ail`; the same with `</dev/null`; `v0.51.0:cmd/ailang/run_policy_supervise.go:72` | `2 4 6 8 10`, rc 0; with `/dev/null`: empty output, rc 0. Source: `cmd.Stdin = os.Stdin` |
| V13 | `argv.ail` (prints `getArgs()`): `-- argv.ail -- --caps=IO,Net --policy /etc/x - -- "a b"`; without the inner `--`; source `main_run.go:240-289`, `run_flag_order.go:19-36` | `ARGV=[--caps=IO,Net\|--policy\|/etc/x\|-\|--\|a b]`, rc 0. Without the inner `--`: rc 1 (a misplaced-flag refusal). Flag parsing stops at the path, so argv **never reaches the flag parser**; one leading `--` is stripped |
| V14 | `echo '"x"' \| H $A run --policy P --args-json - -- rd.ail` | `Error: empty stdin for -args-json -`: `--args-json -` reads stdin, and row 134's `json.Valid` guard keeps `-` out |
| V15 | trusted_host `[IO,FS,Env]`: `args.ail -- numbers.txt`; `env.ail` printing HOME / WORLD_SECRET / PATH with `WORLD_SECRET=leak` set; source `run_policy_supervise.go:228-231` | argv works; the program sees `HOME=<cache>`, `WORLD_SECRET=leak`, `PATH=/usr/bin:/bin`. trusted_host gives the worker `os.Environ()`, so the program's env equals the handler's child env |
| V16 | *superseded by V34.* restricted `[IO,FS,Net]`, `net_allow=["127.0.0.1"]`, http on; POST to a scratch mock | `E_NET_IP_BLOCKED: localhost IP blocked: 127.0.0.1` (re-confirms row 134 V35) |
| V17 | *superseded by V34/V35.* trusted_host `[IO,Net]` with `net_allow` set to: `127.0.0.1` + http / `127.0.0.1` without http / `localhost` / `127.0.0.1:7655` / `127.0.0.2` | `200` (the mock saw the POST) / `E_NET_PROTOCOL_BLOCKED` / `disallowed 127.0.0.1` / `disallowed` / `E_NET_IP_BLOCKED 127.0.0.2`. Loopback opens only for a **literal** host in trusted_host; `host:port` is not a valid entry |
| V18 | *superseded by V34.* the V17 `127.0.0.1` policy; URLs `:7656` and `:7644/v1/health`; `…row134:host/daemon/daemon.go:773-796` | `connect: connection refused` (the dial was attempted), so **the port is not scoped**. World's 8 GET routes are unauthenticated |
| V19 | `git -C ailang show v0.51.0:internal/runner/handlers.go:290-330`; `resolve.go` defaults; `run_policy.go:126` | restricted-only: `Net.RefuseProxy`, `Process.Confined`, `Env.ProtectGitDir`, byte-limit defaults, the env allowlist, the platform check. Both modes: the `fs_sandbox` root, `fs_deny_write`, `max_fs_transfer_bytes` when set, `[budgets]`, `timeout_ms` |
| V20 | trusted_host + the row-134 deny list; `wr.ail` writes `.git/config`, `.GIT/config`, the deny-listed dirs, `ok.txt`; reads with `..`/absolute paths; then `+ ".git", ".git/**"` | **`.git/config` and `.GIT/config` → written** (trusted_host drops `.git` protection); the deny globs → `E_FS_PROTECTED`; the reads → `escapes sandbox`. With `.git/**` added, `.GIT/config` and `.Git/hooks/pre-commit` are still written (v0.51.0 matched case-sensitively; on v0.52.1 one `.git/**` covers them, V37) |
| V21 | trusted_host with the **8 case variants** `.git/** … .GIT/**` prepended; `.git/config`, `.GIT/config`, `.Git/hooks/pre-commit`, `.gIt/HEAD`, `.git`, `./.GIT/config`, `sub/../.GiT/config`, dotless-ı/dotted-İ forms, `ok2.txt`; source `internal/effects/fs_root.go:183-200` | every `.git` form → `E_FS_PROTECTED`; the ı/İ forms → `no such file` (not aliased); `ok2.txt` written; no hook created. Source: a `/**` pattern is a literal prefix, and any other pattern is `path.Match` |
| V22 | trusted_host admission line without, then with, explicit `max_*` byte limits | without: all `0` (unbounded); with: equal to the restricted defaults (unchanged on v0.52.1, V38) |
| V23 | restricted `["IO"]` vs `["IO","FS"]` with `fs_sandbox`; `link.ail → ../ep2/outside.ail` | IO only: **the outside program ran**. With FS: `program link.ail is outside fs_sandbox` (the root is only resolved when FS is admitted) |
| V24 | restricted `[IO,FS]`, `[budgets] FS = 0`: symlink, stdin, `rd.ail numbers.txt`, `overdecl.ail` | symlink refused; stdin works; read → `E_BUDGET_OPERATOR … limit=0`; an over-declared FS is admitted |
| V25 | grader-style `--caps Net,IO overdecl.ail`; trusted_host `[Env,FS,IO]` running `net.ail` | the grader admits an over-declared unused effect; the policy run → `missing_from_policy:["Net"]` (admission is stricter than the grader) |
| V26 | `module wrongname` at `benchmark/solution.ail`, with and without `--relax-modules` under the policy | both → `MOD010` (`--relax-modules` has no effect under `--policy`) |
| V27 | the v0.51.0 §4.4 variants as rendered: FS-zero restricted (`stdin.ail`); Env variant (`sum.ail -- numbers.txt`); Net variant (`net.ail`); Net variant running `rd.ail`; Env variant running `net.ail`; `summary` of each | `2 4 6 8 10`; `15`; `200`; `E_BUDGET_OPERATOR`; `missing_from_policy:["Net"]`. Summaries: `restricted [FS,IO]`, `trusted_host [Env,FS,IO]`, `trusted_host [FS,IO,Net] ["127.0.0.1"]`, each with its own digest |
| V28 | `policy-tool summary` on a trusted_host variant | `ok:true`; reports mode, caps, `net_allow`, `fs_sandbox`, `timeout_ms`, so a variant can be verified by the policy layer itself |
| V29 | **(row 134, side finding)** the exact row-134 restricted policy: `wr.ail` writes `.CLAUDE/settings.json`, `.Claude/x.json`, `.GITHUB/workflows.yml`; policy-tool `write` `.CLAUDE/y.json`, `.AILANG/cache/z`; `diskutil info /` | the writes succeed and land in `.claude/` and `.github/` (APFS is case-insensitive); only the lowercase form is refused. Fixed in row 134 before merge (R-SE-15) and **closed upstream in v0.52.1** (row-134 design §13 V73) |
| V30 | `git show …row134:host/broker/handlers.go:76-90,165-175`; `…:host/daemon/workspace.go:204-290` | `handlerCommand.stdin []byte` already exists; handlers are cached per episode; policy at `<db-dir>/policies/<ep>.toml` |
| V31 | row-93 design §1 P4, §4.2(a) | its text "`ailang-run` has no stdin, argv or extra caps" — this row makes it false |

### 3.2 v0.52.1 re-measure (2026-10-03, row-135 M0 executor)

`$A` = `~/.pinned-ailang-tools/v0.52.1/ailang`. `B` = row 134's exact rendered policy (restricted, `[IO,FS]`, 198 folded deny patterns, `timeout_ms=8000`, `[budgets] FS=1000`). Each variant below is `B` with only the named keys changed. The mock is a scratch Python HTTP server on two loopback ports, `:17655` (named) and `:17656` (other); it logs each request's port, method, path and `X-Test-Header`. The probe programs are:
- `decl.ail`: the `prompt_injection` reference solution, with `main ! {IO, Declassify}` printing the result;
- `net.ail`: `main(url) ! {IO, Net}` POSTs JSON with `X-Test-Header: value123` through `httpRequest` and prints the status or `ERR <e>`;
- `stdin.ail`, `sum.ail` and `argv.ail`: the floor's `pipeline` and `cli_args` solutions and a `getArgs()` printer;
- `wr.ail`: writes `path`.

| V | Command | Observed (short) |
|---|---|---|
| V32 | `$A --version`; `shasum -a 256 $A`; `gh release view v0.52.1 --repo sunholo-data/ailang --json isDraft,publishedAt`; `git log --oneline v0.52.0..v0.52.1` (ailang) | `AILANG v0.52.1 Commit c68ded4`; sha256 `0dd70a1d00360be0b7b4c667054c71fdc4ea87aa3deb79710cf56f18ffe1a8f5`. The release reads `isDraft:false` and `publishedAt 2026-10-03T09:50:52Z`; the brief had said Draft. It contains `06a2fabdd` (#1557 Declassify), `96b6e5174` (#1558 port-qualified `net_allow`), `ed1df2d3f` + `32f1c66e6` (#1559/#1554), and `c66f33e4d` + `1099e49f6` (#1551–#1553) |
| V33 | **Declassify.** `decl.ail` under `B`; under `B` with `allowed_caps=["IO","FS","Declassify"]` (`summary` and `run`); under `["IO","Declassify"]` (no FS, deny list kept); `stdin.ail` under the Declassify variant; policy-tool `ai_check decl.ail` under `B` | `B`: rc 2, stdout decision `policy_violation`, `missing_from_policy:["Declassify"]`, empty stderr: shape (c). The Declassify variant gives `summary ok`, mode `restricted`, caps `[Declassify,FS,IO]`. `run`: **rc 0, stdout `1`**, with one `policy:` line, `security_mode:"restricted"`. Without FS the variant fails at load (`fs_deny_write is set but FS is not in allowed_caps`), so FS must stay rendered, as §4.4 already does. An over-declared unused Declassify is admitted (`2 4`, rc 0). `ai_check`: `passed:true`, `verified:3` (sanitizeBody, isInternal, safeForward). Source: `resolve.go` adds `"Declassify": true` to `RestrictedEffects` ("a compile-time information-flow gate with no runtime ops") |
| V34 | **Port-scoped loopback, restricted.** `B` with `allowed_caps=["IO","FS","Net"]`, `net_allow=["127.0.0.1:17655"]` and `net_allow_http=true`: `summary`, then `net.ail` POSTing to each URL in the result column. Then the same policy without `net_allow_http`; `wr.ail` writing `.git/config`, `.GIT/config` and `.CLAUDE/x.json` under the Net variant; FS budget 0 running `wr.ail data.txt`; the mock's request log | `summary ok`, `restricted`, `net_allow ["127.0.0.1:17655"]`. Results by URL:<br>• `http://127.0.0.1:17655/`: **`200`**, rc 0.<br>• `:17656` (another live port): `ERR DisallowedHost(127.0.0.1)`.<br>• `:7644/v1/health` (World's port): `ERR DisallowedHost(127.0.0.1)`.<br>• `localhost:17655`: `ERR DisallowedHost(localhost)`.<br>• `[::1]:17655`: `ERR DisallowedHost(::1)`.<br>• `127.0.0.2:17655`: `ERR DisallowedHost(127.0.0.2)`.<br>• `127.0.0.1:17655/redirect` (302 to `:17656`): `ERR Transport(E_NET_DOMAIN_BLOCKED: domain not in allowlist: 127.0.0.1:17656)`.<br>• `https://example.com/`: `ERR DisallowedHost(example.com)`.<br>Without `net_allow_http`: `E_NET_PROTOCOL_BLOCKED`. Writes: `E_FS_PROTECTED` for `.git/config` and `.GIT/config` (`under .git, which is read-only in restricted mode`) and for `.CLAUDE/x.json`, so restricted protection holds in the Net variant. FS budget 0 gives `E_BUDGET_OPERATOR … limit=0`. **The mock logged exactly two requests, both on :17655** (`/` and `/redirect`, `X-Test-Header: value123`); `:17656` saw nothing. A refused effect still exits rc 0 (134:V35) |
| V35 | **Load-time refusals, restricted.** `B` + `[IO,FS,Net]` with `net_allow` set to: `["127.0.0.1"]`; `["localhost:17655"]`; `["10.0.0.1:80"]`; `["169.254.169.254:80"]`; `["127.0.0.1:0"]`. Each through `summary` and `run` | All refused **at policy load** (summary `ok:false`; `run` rc 1 `--policy: policy: …`), each by name:<br>• `127.0.0.1`: `net_allow entry "127.0.0.1" would open every loopback port — restricted mode admits loopback only port-qualified (127.0.0.1:PORT or [::1]:PORT) … set security_mode = "trusted_host" to grant all of loopback`;<br>• `localhost:17655`: `names loopback by name — restricted mode admits loopback only as a port-qualified literal`;<br>• `10.0.0.1:80`: `names a private address, which restricted mode never connects to (no override)`;<br>• `169.254.169.254:80`: `names a link-local address`;<br>• `127.0.0.1:0`: `is malformed: … port must be 1-65535`.<br>This is the release's compatibility break. The v0.51.0 draft's Net variant (`trusted_host` + bare `127.0.0.1`) still loads in `trusted_host`, but this design no longer renders it |
| V36 | **Env.** `B` with `allowed_caps=["IO","FS","Env"]` in restricted, then the same in `trusted_host`. In `trusted_host`: `env.ail` with `WORLD_SECRET=leak` in the handler env; `sum.ail -- numbers.txt` | Restricted: `summary` refused, `allowed_caps admits Env, which restricted mode has no confined adapter for — set security_mode = "trusted_host" … (restricted admits: AI,Clock,Declassify,FS,IO,Net,Process,Rand,Stream)`. **Env is still trusted_host-only.** trusted_host: `summary ok`. `env.ail` prints `HOME=<cache>`, `WORLD_SECRET=leak`, `PATH=/usr/bin:/bin`, so the program still sees the handler's child env (V15 unchanged). `sum.ail` prints **`15`**, rc 0 |
| V37 | **trusted_host `.git` compensation.** `trusted_host` `[IO,FS,Env]` with `B`'s deny list (no `.git` entry): `wr.ail` → `.git/config`, `.GIT/config`, `.Git/hooks/pre-commit`. Then **one** `".git/**"` prepended, on the `git init` topology: `.git`, `.GIT`, `.git/config`, `.GIT/config`, `./.GIT/config`, `sub/../.GiT/config`, `.Git/hooks/pre-commit`, `.gIt/HEAD`, `.gıt/config` (dotless ı), `.gİt/config`, `ok2.txt`. Then on a `git worktree add` topology (`.git` is a pointer **file**): `.git`, `.GIT`, `.Git`, `ok3.txt` | No `.git` entry: **all three are written** (`wrote`, rc 0; `.git/config` became `pwned`), so trusted_host still drops `.git` protection (V20 unchanged). With **one** `.git/**`, every `.git` form, including the bare `.git` and the worktree pointer file, gets rc 1 `E_FS_PROTECTED: … matches fs_deny_write ".git/**"`. The ı/İ forms get `no such file or directory`; they are not aliased. `ok2.txt`/`ok3.txt` are written. `.git/config` and the pointer file are byte-unchanged. The 8 case variants of V21 are no longer needed **[REFINED]** |
| V38 | Admission line of `B` and of the trusted_host Env variant, without and then with explicit `max_*` | `B` (restricted): `limits {max_fs_transfer_bytes 8388608, max_module_graph_bytes 16777216, max_output_bytes 8388608, max_source_bytes 1048576}`. trusted_host without explicit limits: all `0` (unbounded). With them set explicitly: equal to the restricted values. **V22 is unchanged** |
| V39 | `stdin.ail` under `B` with `1..5`, then `</dev/null`; `argv.ail -- --caps=IO,Net --policy /etc/x - -- "a b"` under the Env variant; the same without the inner `--`; `--args-json -` | `2 4 6 8 10` rc 0; empty, rc 0; `ARGV=[--caps=IO,Net\|--policy\|/etc/x\|-\|--\|a b]` rc 0. Without the inner `--`: rc 1 `flag --caps=IO,Net comes after the file path …`. `--args-json -`: rc 2 `empty stdin for -args-json -`. **V12–V14 are unchanged** |
| V40 | `git diff v0.51.0 v0.52.1 -- internal/policy/policy.go` (ailang); a policy carrying `net_allow_localhost = true` | `policy.go` is unchanged, so there are **no new policy keys**; `unknown fields: [net_allow_localhost]`. The new behaviour lives in the existing `net_allow` entry grammar (`internal/policy/net_allow.go`, `internal/effects/net_allow.go`) |
| V41 | `run_policy.go` `refusedWithPolicy` at v0.52.1 | `--net-allow-localhost` and `--stream-allow-localhost` are still refused under `--policy`, with the hint `list the loopback port in net_allow (127.0.0.1:PORT opens that one port)`. The policy file stays the only source of the grant |
| V42 | `go test … -run 'TestAilang…\|TestWorkspace'` on v0.52.1 (row-134 design §13 V70) | Every row-134 run shape is unchanged: (a) admitted, (b) runtime-refused, (c) refused before execution, (d) the supervisor limit, plus the fifth pre-execution shape. The handler's composer needs no change for the new variants (V33–V36 produce shapes (a) and (c) only) |

## 4. Design
### 4.1 Tool schema (8 tools unchanged; only `ailang-run` changes)

| Key | Type / bound | Meaning |
|---|---|---|
| `path` | row-134 path predicate | unchanged |
| `args_json` | JSON string | unchanged (still refuses `-`, V14/V39) |
| `stdin` | string, ≤ 65 536 UTF-8 bytes | the program's whole stdin; absent = `/dev/null` |
| `argv` | array, ≤ 32 strings, each ≤ 1 024 bytes, no NUL | program arguments, always passed after an inner `--` (V13) |
| `caps` | distinct names from `{IO, FS, Env, Net, Declassify}` **[REFINED]** | the run's **exact** capability set, like the grader's `--caps` (V6); absent = row 134's `{IO, FS}` |

- **A leading `-` in argv is admitted [REFINED].** argv can't reach the flag parser (V13, V39), so refusing it would only break programs that take flags. The injection is closed structurally, and MUT-NO-INNER-DASHDASH pins that.
- **Env together with Net** is refused in v1 (R-135-4). Declassify combines with anything, because it has no runtime operations (V33).
- **Effect selection** (L1 = 1 still holds) is a contracted pure `runEffect(hasEnv, hasNet)`:
  - Env → `Ailang.RunEnv`;
  - Net → `Ailang.RunNet`;
  - otherwise `Ailang.Run`, including `{IO, Declassify}` (D-135-5 = A).
  - All three use scope `worktree`, cost 1.
- **Payload** is `{path, args_json?, stdin?, argv?, caps?}`, with `caps` canonicalized (sorted, unique). With none of the new args, the plan bytes are **byte-identical to row 134's**.

### 4.2 Plan-law impact
- **M1's closed key sets are unchanged** (V4). The new fields go inside `payload`, which stays within L4.
- **L3:** `ailang-run`'s `declaredEffects` grows from 1 to 3 triples. `access` stays `(Ailang.Run, worktree, 0)`, so `tools/list` is unchanged at 8 names.
- **L-ARGS:** `argKeyAllowed` now admits `{path, args_json, stdin, argv, caps}`. New contracted helpers: `capAllowed`, `runEffect`, and bounds predicates.
- **R8 (V3):** the workspace registry **always** binds `Ailang.RunEnv` and `Ailang.RunNet`, taking it from 6 to 8 effect names. The operator allowlist is enforced inside the handler, never by leaving a name unregistered.

### 4.3 Three gates on every extra capability
1. **Broker (who):** a session grant `Ailang.RunEnv=worktree:N` or `Ailang.RunNet=worktree:N`. Without one, the call is recorded as `denied` and nothing is debited (134:V55). A Declassify-only run needs only the `Ailang.Run` grant (D-135-5).
2. **Operator (what is possible):** `serve --run-allow-caps Env,Net,Declassify`. A capability outside that list gets a typed handler refusal. Startup refuses:
   - unknown names;
   - `Net` without `--run-net-allow`;
   - `--run-net-allow` without `Net`;
   - any `--run-net-allow` entry that is not `IPLITERAL:PORT`;
   - any variant that `summary` refuses at load (V35) — the policy layer's own check, surfaced at startup rather than on the first call.
3. **Policy layer (what runs):** the variant's `allowed_caps` is exactly the requested set, plus FS at budget 0 when FS wasn't requested (§4.4). It is **verified by `policy-tool summary`** (V28) and enforced by AILANG.

The effect record carries the effect name, the request (stdin, argv, caps) and the output. The output now adds `policy: {digest, security_mode, caps, net_allow}`.

### 4.4 Policy rendering: per episode, per cap set, lazily, cached
`RenderRunVariant(root, caps, op)` writes `<db-dir>/policies/<ep>.run-<caps>.toml` (mode 0600, outside the root). It verifies the file once with `summary` and caches it with the handler (V30). A `{FS, IO}` request maps to row 134's base policy.

| Requested caps | Mode | Rendered additions (on top of row 134's key set) | Evidence |
|---|---|---|---|
| ⊆ {IO, FS, Declassify} | restricted | `Declassify` in `allowed_caps` when requested. If FS was not requested, `FS` is still allowed with `[budgets] FS = 0`. That keeps the entry-inside-sandbox check and the load-time requirement that FS be present with a deny list, and FS fails at use, as the grader would | V12, V24, V27, V33 |
| contains Net (no Env) | **restricted [REFINED]** | `net_allow = <the --run-net-allow HOST:PORT literals>`; `net_allow_http = <--run-net-allow-http>`; FS budget 1000, or 0 if not requested. Nothing else is needed: restricted mode keeps `.git` protection, the proxy refusal, Process confinement and the byte limits | V34, V35, V38 |
| contains Env (no Net) | **trusted_host** | **one** `.git/**` entry prepended to the deny list (V37; it replaces V21's 8 variants); the 4 byte limits set explicitly to the restricted defaults (V38); FS budget 1000, or 0 if not requested | V36–V38 |

**What trusted_host still relaxes, for Env runs only (stated, not hidden; V19):**
- **The proxy refusal is off.** Moot: the variant has no Net.
- **Process confinement is off.** Moot: Process is never rendered.
- **The worker sees the parent env.** That is the handler's minimal `childEnv` (V15, V36): `HOME`, `PATH`, `LANG`, `LC_ALL`, `AILANG_CACHE_DIR`, and optionally `AILANG_EXAMPLES` (R-135-5).

**Row 134's folded deny list (R-SE-15) is rendered in every variant.** On v0.52.1 it is redundant with upstream folding, but harmless (row-134 design §13 V73).

### 4.5 Handler (`executeRun`)
- **Accepts** `Ailang.Run`, `Ailang.RunEnv` and `Ailang.RunNet`.
- **Refuses:** unknown or case-aliased payload keys; a `caps` set that doesn't match the effect name; argv or stdin over their bounds; and a NUL in argv.
- **Executes** `run --policy <variant> [--args-json J] -- <path> [-- <argv…>]`. stdin goes in through the existing `handlerCommand.stdin` plumbing (V30). No widening flag is ever passed (V41).
- **Unchanged from row 134:** the outcome shapes (V42), the deadline (8 s policy, 10 s handler) and `childEnv`.

### 4.6 Operator configuration (`serve`)

| Flag | Default | Meaning |
|---|---|---|
| `--run-allow-caps Env,Net,Declassify` | empty | the extra capabilities a run may request |
| `--run-net-allow HOST:PORT` (repeatable) **[REFINED]** | none | port-qualified IP literals rendered as `net_allow`, e.g. `127.0.0.1:7655`. Each entry opens **exactly that host and port**: other ports, `localhost`, `[::1]`, other loopback addresses and redirect hops are refused (V34). A bare host, a loopback name, or a private, link-local or unspecified literal is refused at startup, mirroring the policy's load-time refusal (V35) |
| `--run-net-allow-http` | off | renders `net_allow_http = true` (the benchmark mock is plain http) |

### 4.7 Grader parity (row 93 V9)

| Grader (V6) | World `ailang-run` |
|---|---|
| `--caps <spec.caps>` | `caps: spec.caps` |
| `[-- cli_args]` | `argv: spec.cli_args` |
| stdin piped | `stdin: spec.stdin` |
| `--net-allow-http --net-allow-localhost` (all of loopback; the mock binds an ephemeral port) | operator `--run-allow-caps Net --run-net-allow 127.0.0.1:<mock port> --run-net-allow-http`: **one port**. The row-93 harness must bind the World arm's mock on the operator-named port (R-135-9) |
| `--relax-modules` | not available under `--policy` (V26) |
| `env = os.Environ()` | the handler's `childEnv` |
| timeout 10 s | 8 s policy / 10 s handler |

Parity runs **one way, by design: World is never more permissive than the grader.**
- An over-declared effect is refused at admission (V25), except FS, which is admitted at budget 0 (V24), and Declassify, which has no runtime operations (V33).
- A module-name mismatch is refused (V26).

Per task (all four are runnable on v0.52.1):
- `pipeline` → `{path, caps:["IO"], stdin:"1\n2\n3\n4\n5\n"}` (`Ailang.Run`, restricted);
- `cli_args` → `{path, caps:["Env","FS","IO"], argv:["numbers.txt"]}` (`Ailang.RunEnv`, trusted_host + compensations);
- `api_call_json` → `{path, caps:["IO","Net"]}` (`Ailang.RunNet`, **restricted**, one port);
- `prompt_injection` → `{path, caps:["Declassify","IO"]}` (`Ailang.Run`, restricted) **[REFINED]**.

### 4.8 Replay
stdin, argv and caps are in the content-addressed request object. Replay returns the recorded output and never re-executes (134:V10, V46). The variant bytes can be re-derived, and `policy.digest` detects drift.

## 5. Premises

| P | Premise | Evidence |
|---|---|---|
| P1 | `ailang-run` has no stdin/argv/caps today; restricted IO/FS only | V2, V31 |
| P2 | stdin and argv pass through `run --policy`; argv cannot reach the flag parser | V12, V13, V39 |
| P3 | Env still needs trusted_host; loopback Net and Declassify do not | V8, V33, V34, V36 |
| P4 | Declassify is admitted in restricted mode on v0.52.1; the `prompt_injection` reference prints `1` and its contracts verify | V33 (supersedes V9–V11) |
| P5 | trusted_host's relaxations, now for Env only, are compensated in the rendered policy, or moot | V19, V36–V38 |
| P6 | A restricted Net variant reaches exactly the named `IP:PORT`; other ports (World's included), names, other loopback addresses and redirect hops are refused; bare, named and private entries fail at load | V34, V35 (supersedes V16–V18) |
| P7 | FS at budget 0 keeps the entry check and grader-equivalent FS failure; FS must be present whenever a deny list is | V23, V24, V33 |
| P8 | Plan laws need no change; R8 forces registering all three run effects | V3, V4 |
| P9 | A variant can be verified with `policy-tool summary`, and load-time refusals surface there | V28, V33–V36 |
| P10 | World's run is never more permissive than the grader | V6, V25, V26, V34 |
| P11 | The handler plumbing already carries stdin | V30 |
| P12 | No new policy keys; the outcome shapes are unchanged | V40, V42 |

## 6. Milestones

**M0 (this change; attended executor 2026-10-03) — re-measure and pin bump.**
- AC0.1 (done): V12–V27 re-run on v0.52.1 as V32–V42. Row 134's real-binary suites re-proved (row-134 design §13).
- AC0.2 (done upstream): the three asks shipped in v0.52.1 — ailang#1557 (V33), #1558 (V34, V35) and #1559 (row-134 §13 V73).
- AC0.3 (done): the tool-binary pin moves to v0.52.1:
  - `daemon.ToolBinaryRelease = "AILANG v0.52.1"`, which the version-gate tests (`TestToolBinaryIsTheMeasuredRelease`, `TestSeToolsE2EBinaryVersions`) and the startup refusal follow;
  - the audited `cli` flag table re-baselined to row-134 §13 V78;
  - `QUICKSTART`, `tools/attended/se_smoke.sh`, `website/docs/**` and the floor scripts' default `$TOOL`;
  - the CI step that installs the linux x64 v0.52.1 tool binary, verified against the release's own `.sha256` asset.
  - The `.ail` interpreter pin (v0.41.0) is unchanged.

**M1 (1 d) — transition and descriptor.**
- AC1.1: `run.ail` named tests pin the exact plan bytes for each case:
  - no new args (byte-equal to row 134's);
  - stdin;
  - argv;
  - Env → `RunEnv`;
  - Net → `RunNet`;
  - `["Declassify","IO"]` → `Run`, with caps in the payload;
  - Env+Net → refused;
  - an unknown cap, a duplicate cap, and `caps` given as a string;
  - argv over bounds: 33 items, a 1 025-byte item, a NUL;
  - stdin of 65 537 bytes.
- AC1.2: the contracts verify; `verify_ail.sh` is green.
- AC1.3: `ailang-run` declares 3 triples and has 5 properties with `additionalProperties:false`; `PublishSet` still yields 8 descriptors and `tools/list` still lists 8.

**M2 (1 d) — handler and variants.**
- AC2.1: the exact key set per variant class. The checks:
  - The Declassify and Net variants are `restricted`.
  - The Env variant is `trusted_host` with one `.git/**` and the 4 explicit limits.
  - No variant contains a bare-host `net_allow`.
- AC2.2: real binary:
  - stdin → `2 4 6 8 10`;
  - `RunEnv` + `argv:["numbers.txt"]` → `15`;
  - `RunNet` against an in-test 127.0.0.1 mock on the operator-named port → `200`; the mock saw exactly one POST, with `X-Test-Header: value123`;
  - Declassify `decl.ail` → `1`.
- AC2.3: each of these is refused with a stub exec counter of 0:
  - Env, Net or Declassify not allowlisted;
  - Net without `--run-net-allow`;
  - an effect/caps mismatch;
  - an unknown or aliased key;
  - a NUL in argv;
  - oversize stdin;
  - `args_json:"-"`.
- AC2.4: the Env (trusted_host) confinement matrix on both topologies (V37): every `.git` form, including bare `.git`, `.GIT` and the worktree pointer file, → `E_FS_PROTECTED`, with the bytes unchanged; a sibling read → `escapes sandbox`. The same `.git/config` and `.GIT/config` rows under the Net variant → `E_FS_PROTECTED` (restricted protection).
- AC2.5: Net refusals:
  - a Net program under the base, FS-zero, Env and Declassify variants → `missing_from_policy:["Net"]`;
  - under the Net variant, a second live loopback port, World's `:7644`, `localhost:<port>`, `[::1]:<port>` and a 302 redirect to another port → `DisallowedHost` / `E_NET_DOMAIN_BLOCKED`; the second mock logs zero requests.
- AC2.6: an Env program printing `getEnvOr("WORLD_SESSION","<none>")` → `<none>`.
- AC2.7: `argv:["--caps=IO,Net","-","--","a b"]` comes back verbatim.
- AC2.8: `caps:["IO"]` with a symlinked entry → `outside fs_sandbox`; `caps:["IO","Net"]` reading a file → `E_BUDGET_OPERATOR`.

**M3 (0.75 d) — wiring and docs.**
- AC3.1: the startup-refusal table, including the bare, named and private `--run-net-allow` rows (V35); `serve --help` stays bound to QUICKSTART.
- AC3.2: a `/mcp/` end-to-end run for each of the **four** tasks. Each commits; its `world.effects[0].record` names the right effect; its request object holds the stdin, argv and caps.
- AC3.3: with no `RunEnv` grant → `denied` and no exec; with budget 1 → the second call is `denied:budget`.
- AC3.4: a nil-registry replay is byte-equal.
- AC3.5: QUICKSTART §9 is updated.

**M4 (attended, 0.5 d).** Mark republishes. A curl client then writes and runs each task's reference solution with the §4.7 mapping:
- all four print their expected stdout (`2 4 6 8 10`, `15`, `200`, `1`);
- Net while the daemon runs without Net → refused.

## 7. Load-bearing mutations

| Mutant | Killer |
|---|---|
| MUT-NO-INNER-DASHDASH: argv appended without the inner `--` | AC2.7 |
| MUT-ARGS-JSON-STDIN: the `json.Valid` guard dropped | AC2.3 `args_json:"-"` |
| MUT-CAPS-WIDEN: the variant renders the operator allowlist instead of the requested caps | AC2.5 |
| MUT-ALLOWLIST-SKIP: the handler ignores `--run-allow-caps` | AC2.3 |
| MUT-EFFECT-SELECT: the plan emits `Ailang.Run` for an Env request | AC1.1 + AC3.3 |
| MUT-GIT-DROPPED: the `.git/**` entry omitted from the trusted_host variant | AC2.4 `.git/config` and `.GIT/config` written |
| MUT-LIMITS-UNSET: the byte limits omitted | AC2.1 |
| MUT-FS-DROPPED: FS omitted when not requested | AC2.8 symlink (and the V33 load-time refusal) |
| MUT-FS-BUDGET-OPEN: FS at 1000 when not requested | AC2.8 `E_BUDGET_OPERATOR` |
| MUT-NET-TRUSTED: the Net variant rendered trusted_host | AC2.1 (mode) + AC2.4 Net-variant `.git/config` row |
| MUT-NET-UNPORTED: `net_allow` rendered without the port | AC2.2 `200` (the restricted variant fails at load, V35) |
| MUT-NET-PORT-WIDEN: a second port, or `localhost`, rendered beside the named entry | AC2.5 second-port and `localhost` rows |
| MUT-DECLASSIFY-DROP: Declassify not rendered | AC2.2 `decl.ail` → `1` |
| MUT-ENV-PARENT: `os.Environ()` passed to the run | AC2.6 |
| MUT-STDIN-UNRECORDED: stdin not taken from the payload | AC3.2 + AC3.4 |
| MUT-STDIN-CAP: the size bounds removed | AC1.1 + AC2.3 |
| MUT-VARIANT-UNVERIFIED: the summary check skipped | AC2.1 (a tampered variant is refused) |

## 8. Decisions — RULED by Mark, attended 2026-10-03 (D-WORLD-56); one new decision

**D-135-1 = B: HOLD the row for upstream Declassify (ailang#1557).** D-135-2 = A, with loopback port scoping asked upstream (ailang#1558). D-135-3 = A. D-135-4 = A.

**Effect of v0.52.1 on the rulings (measured, V32–V36):**
- **D-135-1:** the hold's condition is met (V33), so the row is un-held, and all four tasks ship together as the ruling asked.
- **D-135-2:** the upstream ask shipped (V34). Net is now port-scoped. That is **narrower** than the ruled host-level `127.0.0.1`. The row-93 audit AC that compensated for unscoped loopback (no committed solution targets `:7644` or `/v1/`) is no longer needed for confinement: AC2.5 proves `:7644` refused. It may stay as a cheap row-93 check.
- **D-135-4:** trusted_host is used for Env only. That is **narrower** than the ruled "Env and Net". Neither change widens anything the ruling allowed.

The original options follow.

- **D-135-1 Declassify.**
  - **A: ship stdin, argv, Env and Net now.** `prompt_injection` stays check-only in World and is reported per task in row 93; the upstream ask is filed in M0.
  - B: hold row 135 until upstream ships Declassify.
- **D-135-2 The loopback port can't be scoped** (true on v0.51.0 only).
  - **A: accept host-level `127.0.0.1` behind the explicit operator flag** (R-135-2), plus a row-93 audit AC: no committed solution targets `:7644` or `/v1/`.
  - B: hold Net until upstream supports ports.
- **D-135-3 Effect naming.**
  - **A: per-cap effects `Ailang.RunEnv` / `Ailang.RunNet`.**
  - B: one `Ailang.Run`, with the caps in its payload.
- **D-135-4 trusted_host.**
  - **A: use it only for Env and Net runs, with the measured compensations.**
  - B: don't use it, and `cli_args` and `api_call_json` stay unrunnable.
- **D-135-5 (NEW, for Mark) Declassify's effect and gate.**
  - **A (recommended): a Declassify run stays `Ailang.Run`.** It stays behind the operator's `--run-allow-caps Declassify`, and caps are recorded in the request object. Declassify reaches nothing on the host: it is an admission-only label with no runtime operations (V33). A per-cap grant would add a fourth run effect with no confinement to account for.
  - B: a fourth effect `Ailang.RunDeclassify` (a separate grant and budget), for symmetry with D-135-3.

## 9. Risks and residuals
- **R-135-1 trusted_host is weaker by construction.** It is now confined to **Env runs only**. Its safety rests on the rendered compensations (one `.git/**`, explicit limits) and on V19's list of what it relaxes. Owner: the R-SE-8 re-proof, plus AC2.1/AC2.4 on every binary change. An upstream restricted-mode Env adapter (an env allowlist key) would close it; not asked yet.
- **R-135-2 Loopback is unscoped.** **CLOSED by v0.52.1** (V34): a Net run reaches only the named `IP:PORT`; World's `:7644` is refused.
- **R-135-3 Declassify is unavailable.** **CLOSED by v0.52.1** (V33).
- **R-135-4 No Env+Net combination in one run.** Owner: a later row.
- **R-135-5 Env runs see the handler's env.** It holds state-dir paths only, no secrets (V36).
- **R-135-6 Parity runs one way only.** Owner: row 93's per-task report.
- **R-135-7 The entry check depends on FS being admitted.** Closed here by always rendering FS (at budget 0 when not requested). v0.52.1 also refuses a deny list without FS at load (V33).
- **R-135-8 stdin is text-only.**
- **R-135-9 (NEW) The Net grant is per serve, not per trial.** The grader's mock binds an ephemeral port (V6), but World's port is fixed by `--run-net-allow` at `serve`. Row 93's World arm must bind its mock on the operator-named port (or restart `serve` per port), and two concurrent World trials cannot each have their own mock port. Owner: row 93 (harness).
- **R-SE-15 (row 134): case-sensitive `fs_deny_write` on APFS** (V29). Fixed in row 134; **closed upstream in v0.52.1** (row-134 design §13 V73).

## 10. Conflict surface
**Touches:**
- `packages/se-tools/se_tools/run.ail` and `transitions.json` (republished attended);
- `host/broker/handlers_ailang.go` and its tests;
- `host/daemon/workspace.go` (6 → 8 effect names; the tool pin already moved in M0);
- `cmd/ailang-worldd/main.go` (3 flags);
- `docs/QUICKSTART.md` §9;
- `scripts/verify_ail.sh`;
- optionally `tools/attended/se_smoke.sh`;
- the row-93 design's §4.2(a) and P4, plus the row-93 mock-port binding (R-135-9).

**Does not touch:** `effectplan.go`, the broker decision law, the projection, store and coordinator, or the compiler pin.

## 11. Non-goals (restated)
- Env and Net in one run;
- public-network Net;
- binary stdin;
- any Go-side confinement or any run without `--policy`;
- the row-93 harness.

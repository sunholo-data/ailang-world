# w-workspace-exec-toolchain-effect — `Workspace.Exec`: operator-profiled toolchain commands (Go, TypeScript, Python) as one brokered, sandboxed, recorded effect

**Status**: PLANNED. Designer draft 2026-10-03 (attended session; Agent `opus`, pin `claude:claude-opus-5-5`). Queue row 140. **Quorum: BLOCKED in all three rounds (the brief's cap; §12).**
- Each round's objections were measured first (V36–V42). They were mostly real, and every real one was fixed.
- **The round-3 revisions have not been reviewed by quorum.** They are the host/child env split, the in-sandbox `TMPDIR` and the whole-tree archive.
- **RULED (Mark, attended 2026-10-03, D-WORLD-62):** D-140-1 = A, D-140-2 = A, D-140-3 = A, D-140-4 = A, **D-140-5 = B** (Mark reviews the r3 delta attended; no fourth quorum round). **Executor-ready once Mark's r3 review is recorded here.**
- **Direction:** D-WORLD-58 (1): land before 1.0, after rows 135 and 138 (both landed on `dev`: `e6064b6`, `0b18b94`).
- **Precedent:** row 134 (pure plan → one brokered effect → record v2 → handler-free replay) and row 135 (operator allowlists, a restricted default, residuals stated rather than hidden; D-WORLD-55/56/59).
- **Base:** ailang-world `origin/dev` `a8c9f57`. The sandbox is `@anthropic-ai/sandbox-runtime` (`srt`): 0.0.71 as the fleet installs it for pi; 0.0.78 (latest) re-measured on the load-bearing arms (V32).
- **Rig:** macOS 26.6.2 arm64, node v26.0.0, go1.26.6 (switched from Homebrew 1.26.4 via the module cache), Python 3.12.13 (uv store). **Linux is not measured** (V33); M0 measures it before anything else lands. **[M0, 2026-10-05: MEASURED — V43–V49; the pin is srt 0.0.78; §4.4 and §4.6 are amended for Linux.]**

**Clauses**: 3 (every command is an explicit, brokered, budgeted, recorded effect) and 4 (the floor's World arm and the dogfood corpus need non-AILANG builds and tests), plus the north star: "World can build and test the projects we actually work on".

**Estimate**: ≈ 3 executor days (M0–M3) plus ≈ 0.5 d attended (M4). The row's estimate was 2–3 d; M0 (the Linux measurement and the pin) is the addition.

**Probe conditions**:
- Every V-row was run first-party on 2026-10-03. Each ran on a scratch APFS clone (`cp -c`) of a real project, never on the real checkout:
  - the ailang compiler: `go.mod`, `go.sum`, `internal/`, `cmd/`;
  - TwilightGame, without `.git`, `dist` or `public`;
  - `sunholo-platform/cli`, with its `.venv`.
- Decoy secrets are scratch files (`~/.srt-row140-decoy`, text `DECOY-NOT-A-SECRET`). No real secret was read.
- The probe scripts and redacted outputs are banked in `design_docs/verification/world-row140-design-2026-10-03/`. In the outputs, `$S` is the scratch dir and `SRT` is `node …/sandbox-runtime/dist/cli.js`; srt's per-run proxy credential is redacted.
- This doc claims no `.ail` text. Every `.ail` shape here is described, not quoted. The executor writes it under S5 against the pinned binary.

## 1. Why (measured)

- **P1 — World can edit any text file but can build or test only AILANG.** Its eight tools are `ailang-read/-write/-edit` (generic text) plus five AILANG-only tools (V31). There is no shell, and there should not be one. So the ailang compiler (Go), sunholo-platform (Python) and TwilightGame (TypeScript) can be changed through World but not verified through it (D-WORLD-58).
- **P2 — The fleet's sandbox fences writes and network for an arbitrary process tree. Reads are open by default, but the row's "reads are not fenced" is wrong: they can be fenced.**
  - **Writes and network (measured):**
    - Writes outside the allowed roots are refused for the child, for its grandchildren, and for an orphan that outlives the parent (V3).
    - The network is allowlist-only, through srt's own proxy. Loopback, including a stand-in for World's own port, raw sockets, `bind` and direct DNS are all refused (V9).
  - **Reads:**
    - With the default settings, reads are open everywhere (V7).
    - srt also has a deny-then-allow read fence. With `denyRead: [$HOME]` and `allowRead: [worktree, episode cache, toolchain roots]`, a decoy in `$HOME`, another project's file and `ls $HOME` are all refused. `go test`, `tsc` and `pytest` still pass (V8, V19, V23, V24).
    - **So the honest statement is narrower than the row's:** reads are fenced by path list. System paths (`/usr`, `/opt`, `/etc`, `/private/tmp`, `/var/folders`) stay readable (V8b). §9 R-140-1 says exactly what stays open.
- **P3 — Two srt CLI defaults would each break the design if used naively (V2, V11, V12):**
  - **The option-injection hazard.** Without a `--` after srt's own options, a wrapped command's `-c '…'` runs as an srt command string. `-s <file>` replaces srt's settings file outright, which is a full confinement escape.
  - **The exit-status hazard.** A child killed by SIGTERM or SIGINT reports **rc 0**.
  - **Both have measured fixes:** always put `--` after srt's options, and run the command under a one-line `sh` trampoline that turns signal deaths into 128+n.
- **P4 — The toolchains run confined once their caches are moved:**
  - Go fails outright with worktree-only writes (`~/Library/Caches/go-build`, V16). With a per-episode `GOCACHE` and a read-only shared module cache it passes, offline, and writes nothing to the shared cache (V18).
  - `tsc` and `vitest` need no change beyond `HOME` and `TMPDIR` (V22).
  - Python writes `__pycache__` into the worktree unless `PYTHONPYCACHEPREFIX` moves it (V24).
- **P5 — Overhead is small next to the commands.**
  - srt costs ≈ 150–165 ms per call, median of 20 (V15).
  - Measured sandboxed and unsandboxed on the same warm cache, it costs +0.2 s on `go test` (0.8 s vs 0.6 s) and +0.6 s on `tsc` (6.1 s vs 5.5 s) (V18, V22).
- **P6 — The binding limit is time, not overhead.** World's call deadlines are frozen: invoke 20 s, handler cap 10 s, write 30 s. A source test guards them (V30).
  - **Targeted commands fit:**
    - a whole-project `tsc --noEmit` on TwilightGame takes 6.1 s;
    - one vitest file takes 4.9 s;
    - one pytest file takes 1.1 s;
    - two Go packages take 0.8 s warm and 5.6 s cold.
  - **Full suites do not:** TwilightGame's full vitest takes 16.2 s (V27), and World's own CI `go test` takes ≈ 155 s (STATUS iteration 228).
  - This shapes D-140-4 and the cache decision D-140-3: a cold cache can, on its own, push a call past 10 s.

## 2. Goals and non-goals

**Goals.**
- **One new tool, `workspace-exec`, with one new effect, `Workspace.Exec` (scope `worktree`, cost 1).** It runs **one command of an operator-written project profile**: a fixed argv prefix plus agent arguments that must match the command's own grammar. The cwd is the profile's project root inside the episode worktree.
- **Every run is confined by srt** with a World-rendered settings file:
  - writes go only to the worktree and the episode's own exec cache;
  - reads are fenced by path list (D-140-1);
  - there is no network in v1 (D-140-2);
  - the environment is built by World, never inherited.
- **The record holds everything replay needs:** exit status, signal, timing, bounded stdout and stderr (head and tail, with totals and SHA-256s), and the digests of the profile, settings and sandbox. Replay never re-executes.
- **Fail closed:**
  - no profile → a typed refusal;
  - no sandbox, or a sandbox failing its startup probe → `serve` refuses to enable the tool;
  - never an unsandboxed run, on any platform.

**Non-goals (with owners).**
- **Docker or any container runtime: never** (row text).
- **A general shell:** never. There is no `sh -c` from the agent, and no free-form argv.
- **Long-running commands beyond the frozen deadlines:** D-140-4, follow-up row.
- **Network during exec:** D-140-2 (default none); dependency installs are an operator step in v1.
- **Windows:** srt's Windows support is alpha (README); not a World target.
- **Module-root and registry-cache layouts for the AILANG tools:** row 141. This row's profile `root` covers the same need for exec only (§10).

## 3. Verification Log (2026-10-03, rig as above; outputs under `…/world-row140-design-2026-10-03/out/`)

| V | Command | Observed (short) |
|---|---|---|
| V1 | `cat ailang/tools/pi-extensions/sandbox/package.json`; installed `package.json`; `shasum -a 256 dist/cli.js`; `npm view @anthropic-ai/sandbox-runtime version time.modified`; `grep -n sandbox-runtime ailang/scripts/mission_pi_run.sh` | dep `^0.0.71`, installed **0.0.71**, `cli.js` sha256 `e9653f95…07a6`. Latest is **0.0.78** (modified 2026-10-01). `mission_pi_run.sh:175-179` preflights the package (`sandbox_unavailable`, rc 15). The fleet policy `sandbox.mission.json` widens writes to `~/Library/Caches/go-build`, `~/go/pkg/mod` and `~/.cache`, and `denyRead`s `~/.ssh`, `~/.aws`, `~/.gnupg` and two secret files |
| V2 | read `dist/cli.js` (0.0.71) | argv is `shell-quote`d into one string and run with `spawn(cmd, {shell:true})`. The commander program uses `.allowUnknownOption()`. A missing or invalid `--settings` file → `Refusing to run with the default config`, rc 1. Child SIGINT/SIGTERM → `process.exit(0)`; any other signal → rc 1 |
| V3 | `p1_fs.sh`: `allowWrite=[ws]`; write `in.txt`, `../other/x`, `$HOME/x`, `/tmp/x`, and through `ws/sub/lnk → ../../other`; grandchild `sh→sh→bash` write; an orphan `nohup` write 2 s after srt exits | inside: written. All outside writes: `Operation not permitted`, rc 1, including through the symlink, the grandchild and the orphan; `ls` shows no file |
| V4 | `p3_paths.sh`: no `denyWrite`; write `/tmp/claude/x`, `~/.claude/debug/x`, `~/.npm/_logs/x`; `echo $TMPDIR`; source `getDefaultWritePaths` | **all three written.** `TMPDIR=/tmp/claude`. srt always allows `/dev/{stdout,stderr,null,tty,…}`, `/tmp/claude`, `/private/tmp/claude`, `~/.npm/_logs` and `~/.claude/debug`. With those four in `denyWrite`, both probed writes → `Operation not permitted` (`denyWrite` wins over allow) |
| V5 | `p3_paths.sh`: inside the allowed root, no operator deny: `.git/config`, `.git/hooks/pre-commit`, `.git/HEAD`, `.mcp.json`, `.github/workflows/ci.yml` | srt's mandatory deny list refuses `.git/config`, `.git/hooks`, `.GIT/config` and `.mcp.json`. **`.git/HEAD` and `.github/…` are written**, so World must render its own deny list |
| V6 | `p3_paths.sh`: `denyWrite=[ws/.github, ws/.claude, ws/.git]` → `.github/…`, `.GITHUB/…` (did not exist), `.Claude/…`, `.git/HEAD`, `.GIT/HEAD` | **all refused, case variants included** (APFS is case-insensitive; Seatbelt matched). Contrast row-135 V29, where AILANG's own matcher needed folding |
| V7 | `p1_fs.sh` with `denyRead=[]`: `cat ~/.srt-row140-decoy`; `cat ../other/file.txt` | both read (rc 0): **by default reads are open** |
| V8 | `p1_fs.sh` + `p4_ts_fence.sh`: `denyRead=[$HOME,(/private)/tmp]`, `allowRead=[ws]`, then a workspace UNDER `$HOME` with `denyRead=[$HOME]`, `allowRead=[ws, ep]` | decoy, other project and `ls $HOME` → `Operation not permitted`. `/etc/hosts` and `/usr/bin` readable. A grandchild read → refused. Under `$HOME`: `node a.js` runs, `pwd -P` works, decoy refused |
| V8b | `p4_ts.sh` TS4 / `p4_ts_fence.sh` F2: `/private/tmp` in `denyRead`, with the workspace re-allowed | `node` aborts at startup (`SIGABRT` in `InitializeOncePerProcess`) even with `TMPDIR` re-allowed. **So the fence denies `$HOME` and World's state dir, not `/tmp`**; `/private/tmp` and `/var/folders` stay readable |
| V9 | `p2_net.sh`: `allowedDomains=["proxy.golang.org"]`; curl allowed and other host; curl to `127.0.0.1:17644` (a live loopback server standing in for World's `:7644`) with and without `--noproxy '*'`; python `connect`, `bind`, `getaddrinfo` | `200`; `CONNECT tunnel failed, response 403` (rc 56); loopback: `Couldn't connect` (rc 7) both ways; `connect` → `PermissionError [Errno 1]`; `bind` → `[Errno 1]`; `getaddrinfo` → `nodename nor servname`. Only the proxy is reachable |
| V10 | `env -i HOME=… PATH=… LANG=C srt -- /usr/bin/env` | the child gets exactly the parent's env plus srt's: `HTTP(S)_PROXY`/`ALL_PROXY`/… carrying a per-run proxy credential, `NO_PROXY`, `GIT_SSH_COMMAND`, `GIT_CONFIG_PARAMETERS`, `CLOUDSDK_PROXY_*`, `TMPDIR=/tmp/claude`, `SANDBOX_RUNTIME=1`. No other variable is added |
| V11 | `p1_fs.sh`, `p2_net.sh`, `p8_proc.py`: `false`; `exit 42`; self-kill with TERM, INT, KILL, SEGV; a missing binary; then the same under `/bin/sh -c '"$@"; exit $?' world-exec <argv…>` | bare srt: `1`, `42`; **TERM → 0, INT → 0**; KILL → 1; SEGV → 1; missing binary → 127. Trampoline: **143, 130, 137, 139**, `exit 3` → 3, and argv `["a b","$(id)","-c","--"]` printed verbatim |
| V12 | `p2_net.sh`: `srt --settings S printf '[%s]' -c 'echo via-c'`; `srt --settings S echo -s /nonexistent.json`; the same arguments after `--` | **without `--`:** `via-c` is executed (srt took `-c`), and `Could not load settings from /nonexistent.json` (srt took `-s`: an agent-written file in the worktree would REPLACE the policy). **With `--`:** `[-c]`, `[--]`, `[$(echo INJECT)]` and `[;echo semi]` print verbatim |
| V13 | `p8_proc.py`: srt in its own session running `sleep 300 & sleep 300 & wait`; `killpg(SIGKILL)` | 4 group members before, **0 after**; node's wait status −9 |
| V14 | `p8_proc.py`: 20 MB on stdout through srt | 20 000 000 bytes in 0.8 s |
| V15 | `p5_overhead.py`, n = 20 | median: bare `true` 2.2 ms; `node -e 0` 46.9 ms; `sandbox-exec` alone 8.4 ms; srt with no network 152.8 ms; with one domain 157.8 ms; with the read fence 163.3 ms |
| V16 | `p4_go.sh` G1: `go test ./internal/lexer/ ./internal/parser/` with worktree-only writes and the default env | `open ~/Library/Caches/go-build/…: operation not permitted`; `FAIL [setup failed]` |
| V17 | `p4_go.sh` G5 / first run: `GOTOOLCHAIN=local`; auto with an empty `GOMODCACHE` and `GOPROXY=off`; auto with the shared module cache | `go.mod requires go >= 1.26.6 (running go 1.26.4; GOTOOLCHAIN=local)`; `download go1.26.6 … toolchain not available`; with the shared cache: `go version go1.26.6`. **The switched toolchain's GOROOT lives in the module cache**, so the module cache must be readable |
| V18 | `p4_go.sh` G2/G2b/G4: per-episode `GOCACHE`/`HOME`/`TMPDIR`; shared `GOMODCACHE`, not writable; `GOTOOLCHAIN=go1.26.6 GOFLAGS=-mod=readonly GOPROXY=off GOTELEMETRY=off` | `ok` for both packages. **Cold 5.6 s, warm 0.8 s**, unsandboxed warm 0.6 s; the cache grows to 141 MB. `find <modcache> -newer marker` → empty: no write was attempted on the shared cache |
| V19 | `p4_go.sh` G3: V18 plus `denyRead=[$HOME, /tmp]`, `allowRead=[ws, ep, GOROOT, modcache]` | `ok` for both; decoy → refused |
| V20 | `p4_go_net.sh`, `p4_go_tls.sh`: a cold per-episode module cache with `-mod=mod`: no network; `proxy.golang.org`+`sum.golang.org` allowed; plus `SSL_CERT_FILE`; plus `enableWeakerNetworkIsolation`; curl through the same proxy; `GOPROXY=direct`; `-mod=readonly` with an empty `go.sum`; an operator-populated cache, then offline | `Forbidden`. **`tls: failed to verify certificate: x509: OSStatus -26276`** (Go on darwin verifies through `trustd`, which Seatbelt blocks). `SSL_CERT_FILE`: same error. Weaker isolation: `ok`. curl: `200`. `GOPROXY=direct`: `CONNECT … 403`. readonly: `missing go.sum entry`. Operator cache + `GOPROXY=off`: `ok` |
| V21 | `p7_seed.sh`: clone V18's warm 141 MB / 2 017-file `GOCACHE` per episode, run `go test`, check the seed | `cp -cR` 0.37 s (`cp -R` 0.52 s); `go test` on the clone **1.2 s**; the seed has no file newer than the marker |
| V22 | `p4_ts.sh`: TwilightGame `./node_modules/.bin/tsc --noEmit` sandboxed and bare; `vitest run <one file>`; `npm test -- run <file>`; writes outside the worktree | tsc rc 0, **6.1 s vs 5.5 s**. vitest: 3 tests pass, 4.9 s cold; on a rerun the reported Duration is 0.37 s. npm test passes. Nothing new under `~/.npm` or `~/Library/Caches`; inside, `node_modules/.vite-temp` |
| V23 | `p4_ts_fence.sh` F3: `denyRead=[$HOME]`, `allowRead=[worktree, ep]`; tsc; decoy | tsc rc 0; decoy refused |
| V24 | `p4_py.sh`, `p4_py2.sh`: `.venv/bin/python -m pytest` (the venv's python is a symlink into the uv store under `$HOME`); `__pycache__` placement; `PYTHONPYCACHEPREFIX`; the read fence with no store / the store's link / link and real path | 3 passed, 1.1 s. 3 `__pycache__` dirs written into the worktree. With the prefix: none in the worktree, 9.7 MB in the episode cache; `.pytest_cache` still lands in the worktree unless `-p no:cacheprovider`. Fence: no store → rc 126 `Operation not permitted`; link only → `No module named 'encodings'`; **link + `realpath`** → 3 passed |
| V25 | `p4_py2.sh` PY5: a test writing to `/tmp`, `$HOME`, a sibling dir and reading the decoy, with and without the fence | all three writes `REFUSED Operation not permitted`. The decoy is read without the fence and `REFUSED` with it |
| V26 | `p4_py.sh` PY4: `uv run --offline --frozen` with a cold per-episode `UV_CACHE_DIR` | `hatchling was not found in the cache` (offline uv needs a provisioned cache) |
| V27 | `p10_full.sh`, `p10b_baseline.sh`: TwilightGame full `vitest run`; `sunholo-platform/cli` full pytest; both bare with the same env | vitest **16.2 s wall**, 1 472 271 output bytes, 25 failed / 1 812 passed. **The failing set is identical bare** (25 = 25, empty diff), so the sandbox caused none. pytest 2.2 s, 9 failed / 199 passed, identical bare and with the real `HOME` |
| V28 | `p6_inject.sh`: `go test -exec "sh -c …"`; `-coverprofile=/tmp/…`; a `toolchain go1.99.0` line in `go.mod` under the operator's `GOTOOLCHAIN`; `npm --prefix <other project> test`; a worktree `.npmrc` `node-options=--require ./zz-evil.cjs`; `pytest -p <module>` | `-exec`: accepted, rc 0 (arbitrary runner). Coverprofile: `operation not permitted`. go.mod toolchain line: ignored under the operator's pin. `--prefix`: the OTHER project's script ran, and its write was refused. `.npmrc`: the injected code ran in every node process, and each of its writes was refused (`EPERM`). `-p`: loads any importable module (`ImportError` for a missing one) |
| V29 | `p11_deps.sh`: `npm install` with `registry.npmjs.org` allowed; then `npm ci --offline` with no network from the warm per-episode cache; `uv sync --frozen` with `pypi.org`+`files.pythonhosted.org` allowed | all rc 0. npm and uv TLS work through srt's proxy, unlike Go (V20). uv: 18 packages in 1.5 s; pytest then passes |
| V30 | `host/coordinator/effectful.go:36-41`; `host/daemon/daemon.go:96-99, 615`; `host/projection/projection_test.go` `TestProjection_NoDeadlineTampering`; `host/broker/handlers.go:18-19, 144-245` | `HandlerCap 10 s`, `HandlerHeadroom 6 s`, `invokeDeadline 20 s`, `writeTimeout 30 s` (frozen D7), coordinator `MaxOutput 1 MiB`. `runBounded` treats output over 8 MiB as an **error and kills the group** (no truncation), uses `Setpgid` + group kill, and `procbound.Admit` caps unreaped children at 8 |
| V31 | `host/daemon/workspace.go:49-57, 217-240`; `packages/se-tools/transitions.json` | 8 tools; 8 effect names, ALWAYS bound to every episode (R8). The operator gates live inside the handler, and an unconfigured feature is a typed refusal (`NoExamplesCorpusRefusal`) |
| V32 | `p9_v078.sh`: six load-bearing arms on 0.0.71 and 0.0.78 | identical: TERM → 0; `-c` without `--` executed; `/tmp/claude` writable; outside write refused; fenced decoy refused; non-allowlisted host 403 |
| V33 | `command -v bwrap`; srt README §Platform Support; `.github/workflows/ci.yml:20,118` | **Linux is NOT measured** (no bubblewrap on the rig). The README says Linux needs `bubblewrap`, `socat` and `ripgrep`, does no glob matching, and blocks Unix sockets with seccomp (x64/arm64). On Ubuntu 24.04+, `kernel.apparmor_restrict_unprivileged_userns` must be 0. CI jobs run `ubuntu-latest` |
| V34 | `ls /opt/homebrew/bin/rg`; every macOS arm above | no `rg` binary (only a shell function); all macOS runs worked. Source: `macos-sandbox-utils.js:242` "no ripgrep needed" on macOS |
| V35 | `p12_topology.sh`, World's topology under `$HOME`: `root/ep1` (its `.git` a worktree **pointer file**), `root/ep2`, `state/world.db`, `state/exec-cache/ep1`. `denyRead=[$HOME, state, root]`, `allowRead=[root/ep1, state/exec-cache/ep1]`, `denyWrite=[ep1/.git, ep1/.github, srt's four]`. Run from `ep1` | write `a.txt` → ok. `.git` (the pointer file) and `.GIT` → `Operation not permitted`. `cat ../ep2/notes.txt`, `ls ..` and `cat state/world.db` → refused. A write to the child's own `HOME` (in the exec cache) → ok. `../ep2/pwn` → refused. The pointer file is byte-unchanged |
| V36 | `p13_quorum_r1.py` (a): `denyRead=[$HOME]` with `allowRead` = `[ws]`; `[ws, $HOME]`; `[ws, /Users]`; `[ws, <symlink → $HOME>]`; `cat` the decoy (quorum r1 premise, `oc-glm-5-3`) | `[ws]`: refused. **`[ws, $HOME]`: `DECOY-NOT-A-SECRET` (re-opened).** `[ws, /Users]` (a strict ancestor): refused. Symlink to `$HOME`: refused |
| V37 | `p13_quorum_r1.py` (b): a live `127.0.0.1:0` listener counting accepts; an unsandboxed control connect; a sandboxed python connect and a sandboxed curl (proxy env) to it (quorum r1 premise, `oc-kimi-k3`) | control: accepts = 1. Python: `REFUSED 1 Operation not permitted`. curl: rc 7. **Accepts after both = 1**: the listener was live and the sandbox refused |
| V38 | `p14_quorum_r2.py` (a): World's own `node` as the probe client, connecting to a live counted listener; an unsandboxed control; then sandboxed; then a missing client (quorum r2 premise, `gemini-3-1-pro`) | control: `WORLD-PROBE-CONNECTED`, accepts 1. Sandboxed: rc 0, **`WORLD-PROBE-REFUSED EPERM`**, accepts still 1. Missing client: **rc 127**, `No such file or directory`, no token |
| V39 | `p14_quorum_r2.py` (b): digest the srt package tree; sha256 of `node` and `libnode.147.dylib` (quorum r2 premise, `oc-glm-5-3`) | srt tree: 153 files, 8 383 570 bytes, **15.9 ms**. `node`: 68 384 bytes, 0.2 ms. `libnode`: 70 269 024 bytes, **32.2 ms** |
| V40 | `p15_quorum_r3.sh` (a)–(c) (quorum r3 premises, `gemini-3-1-pro` and `oc-glm-5-3`): `NODE_OPTIONS="--require ./evil.cjs"` in the env of the host `node` that runs srt (`evil.cjs` in the worktree writes `$HOME/.srt-row140-hostnode-escape`); then a host env of `PATH` only, with the same variable passed to the child through `/bin/sh -c 'exec /usr/bin/env -i "$@"' … K=V…`; then `TMPDIR=<cache>/tmp` set that way for Go `t.TempDir()`, Python `tempfile` and node `os.tmpdir()` | **Host env: `HOST-NODE-REQUIRE-RAN`, and the marker file was CREATED in `$HOME`**, an escape before the sandbox existed. Inside: the child env is exactly the given `K=V` set, and `evil.cjs` runs in the child node but its write fails `EPERM`, with no marker. Temp: `--- PASS: TestTmp` with `TEMPDIR <cache>/tmp/TestTmp…`; `PYTMP <cache>/tmp/tmp…`; `NODETMP <cache>/tmp/n-…` |
| V41 | `p15c_tramp.py`: host env `PATH` only; trampoline `/bin/sh -c '/usr/bin/env -i "$@"; exit $?' world-exec HOME=… TMPDIR=… LANG=…` (no `exec`, so the signal mapping is kept): TERM/INT/KILL self-kill, `exit 3`, `/usr/bin/env`, and a printf of `["a b","$(id)","-c","--","X=Y"]` | **143, 130, 137**; 3. `env` prints exactly the three given variables (no srt proxy variables, no `TMPDIR=/tmp/claude`). argv is verbatim, including `X=Y` after `argv0` |
| V42 | `p15_quorum_r3.sh` (d), `p15b_archive.py` (quorum r3 premise, `oc-kimi-k3`): srt run from a copy of the PACKAGE dir only, then from a copy of its whole `node_modules` tree, both under a state dir in `denyRead`; the child `cat`s the archived `cli.js`; the digest of the whole tree | package dir only: **`ERR_MODULE_NOT_FOUND: Cannot find package 'commander'`**. Whole tree: **rc 0, `ran-from-archive`** (srt's host node is outside the sandbox). The child gets `Operation not permitted` on the archive. Tree: `.bin`, `@anthropic-ai`, `@pondwader`, `commander`, `node-forge`, `zod`; 829 files, 13 846 370 bytes; **23.8 ms** |
| V43 | **[M0]** `world-row140-m0/m0_matrix.py --seccomp-read none` on `ubuntu-latest` (kernel 6.17, bubblewrap 0.9.0, node 22.23.3, srt 0.0.78): the design as written, with the archive under the read-denied state dir | **every arm fails rc 127**: `…/state/archive/exec-sandbox/node_modules/@anthropic-ai/sandbox-runtime/vendor/seccomp/x64/apply-seccomp: No such file or directory`. On Linux srt runs its seccomp helper **inside** the sandbox, and the read fence masks it. 1/32 gating arms pass (`out/linux-0.0.78-run3-both-modes.txt`) |
| V44 | **[M0]** the same matrix with `--seccomp-read allowread` (only `<archive>/…/vendor/seccomp/<arch>` added to `allowRead`); macOS on the same pin; a pass-through (unsandboxed) srt as the control | **Linux 32/32 gating arms pass; macOS 33/33.** The control fails all 21 confinement arms. With the helper dir re-allowed, the child still cannot read the archived `cli.js` (`ENOENT`) |
| V45 | **[M0]** Linux, read-denied regions (`$HOME`, state dir, workspace root): reads; writes; `ls $HOME` | bwrap **masks** a denied dir with an empty tmpfs. Reads → `ENOENT` (decoys, sibling episode, archive). A write into `$HOME` or the state dir prints `WROTE`, but **the host file is absent**: it lands in the throwaway mount. A write into a sibling episode → `ENOENT`. `ls $HOME` → rc 0, empty |
| V46 | **[M0]** Linux network: live listener with an accept count; a sandboxed `bind`; then a host connect to the port the child bound; external https | loopback connect → `ECONNREFUSED`, accepts unchanged (raw, proxy env and `--noproxy`). `bind 127.0.0.1` **succeeds despite `allowLocalBinding:false`**, inside the child's own network namespace; **a host connect to that port fails** (`host reached=False`). External https → `CONNECT tunnel failed, response 403`. DNS → `EAI_AGAIN` |
| V47 | **[M0]** Linux exit status, bare and through the trampoline | bare srt SIGTERM → **143** (the rc-0 hazard is macOS-only); trampoline 143/130/137, `exit 3` → 3, missing binary 127 |
| V48 | **[M0]** Linux V3 (grandchild, orphan, symlink), V12 (argv after `--`; `-c` without it), V13 (group kill), V41 (child env), `.git` pointer file and `.github/x` | identical to macOS: all outside writes absent on the host, argv verbatim, `-c` taken without `--` (the hazard is cross-platform), 2 → 0 sleepers, env exactly the `K=V` set, `.git`/`.github` → `EROFS` with bytes unchanged |
| V49 | **[M0]** Linux `.GITHUB/x`; `/tmp/claude` without World's `denyWrite` | `.GITHUB/x` is **written**: ext4 is case-sensitive, so it is a different, non-special dir (git and GitHub read only `.github`). `/tmp/claude` does not exist on the runner (`ENOENT`) |

## 4. Design

### 4.1 The tool: `workspace-exec` (9th tool; one effect)

| Key | Type / bound | Meaning |
|---|---|---|
| `command` | string, `^[a-z][a-z0-9-]{0,31}$` | the id of one command in the episode's project profile, e.g. `test`, `typecheck`, `test-file` |
| `args` | array, ≤ 16 strings, each 1–512 UTF-8 bytes, no NUL | the agent's arguments. The handler matches them against that command's grammar (§4.3); absent = none |

- **No cwd, env, timeout or argv-prefix key.** All four come from the profile. Unknown keys are refused (L-ARGS, the row-134 rule).
- **Plan** (`packages/se-tools/se_tools/exec.ail`, calling convention v2):
  - it checks the key set and the bounds with contracted predicates (`commandIdOk`, `argOk`, `argsOk`);
  - it emits exactly one effect `{effect:"Workspace.Exec", scope:"worktree", cost:1, payload:{command, args}}`, or a zero-effect refusal plan (L7).
  - It cannot see the profile (operator config, host-side). Bounds are the plan's job; the profile is the handler's, as in row 135, where caps bounds were the plan's and the operator allowlist the handler's.
- **Descriptor:** `access (Workspace.Exec, worktree, 0)`; `declaredEffects [(Workspace.Exec, worktree, 1)]`; `additionalProperties:false`. `tools/list` grows from 8 to 9; effect names from 8 to 9, always bound (R8, V31).

### 4.2 Gates

1. **Broker (who):** a session grant `Workspace.Exec=worktree:N`. Without one: `denied`, nothing executes, nothing is debited (row-134 pattern).
2. **Operator (what is possible):** `serve --exec-profile <file>` (repeatable, one per project). With no profile: a typed refusal, `no exec profile configured: start ailang-worldd serve with --exec-profile FILE`. Otherwise, before any subprocess:
   - a command id not in the profile → refusal;
   - args not matching that command's grammar → refusal naming the first bad arg.
3. **Sandbox (what runs):** srt with the episode's rendered settings (§4.4), launched as in §4.5. The settings are verified once at render (§4.4), and the sandbox itself is verified by a startup probe (§4.6).

**Why it is not a package (S3).** The transition *is* a package (`se-tools`, republished attended). The handler is host Go because executing a process is the host boundary (S2), as `Ailang.Run` is. The **argument-matching law** is pure, so it is authored as a contracted sketch `design_docs/sketches/execargs.ail`: per-argument predicates under Z3, list-level cases as named inline tests where Z3 cannot encode them. A drift test evaluates the sketch's rows through the Go mirror (the `effectplan.ail` pattern, row 134 M1).

### 4.3 The project profile (operator-owned JSON, outside every worktree)

```json
{
  "profile": "world/exec-profile/v1",
  "project": "ailang-compiler",
  "root": ".",
  "path": ["/opt/homebrew/bin", "/usr/bin", "/bin"],
  "env": { "GOFLAGS": "-mod=readonly", "GOPROXY": "off", "GOTOOLCHAIN": "go1.26.6", "GOTELEMETRY": "off" },
  "read_roots": ["/Users/<op>/go/pkg/mod"],
  "caches": { "GOCACHE": { "seed": "/Users/<op>/.ailang/exec-seeds/ailang-compiler/go-build" } },
  "timeout_ms": 9000,
  "commands": {
    "test":      { "argv": ["go", "test", "-count=1"],
                   "flags": { "-run": "regex", "-v": "bool", "-short": "bool" },
                   "positional": "pkgpattern", "max_args": 8 },
    "vet":       { "argv": ["go", "vet"], "positional": "pkgpattern", "max_args": 4 }
  }
}
```

**Load-time refusals** (`serve` exits non-zero, naming the field):
- an unknown key at any level;
- `root` failing `pathOk`, or not an existing directory inside a worktree at first use;
- an `argv[0]` that does not resolve on `path` to an absolute executable. It is resolved **once at load** and recorded by path plus sha256. A later swap of the binary changes the recorded digest and is refused at call time: re-hash on each call, ≈ ms;
- `env` naming any World-owned variable: `HOME`, `TMPDIR`, `PATH`, `SANDBOX_RUNTIME`, `*_PROXY`, `NO_PROXY`, `GIT_*`, the cache variables, or anything `childenv.RegistryVariables` strips;
- **a `read_roots` entry that could re-open a denied root [REVISED r1].** An entry is refused when its literal path **or** its `realpath`:
  - equals or is an ancestor of the operator `$HOME`, the World state dir or the workspace root;
  - lies inside the state dir or the workspace root;
  - equals or contains one of the startup probe's fixed decoys (§4.6 arm 3);
  - is `/`.
  - **Why the rule is this strict:** srt's own precedence is narrower and platform-specific. On 0.0.71/macOS an `allowRead` **equal** to a denied root re-opens it, a strict ancestor (`/Users`) does not, and a symlink to `$HOME` does not (V36). But World adds each root's `realpath`, which would turn a symlink into exactly the re-opening case, and Linux is unmeasured. So World refuses the whole class rather than rely on srt's matcher;
- `timeout_ms` > `HandlerCap` − 1 000 (D-140-4).

**Argument grammar (the injection answer, V12, V28):**
- **argv only.** No shell sees an agent byte. The trampoline's `"$@"` passes argv verbatim (V11), and srt's own options end at `--` (V12).
- **The prefix is fixed.** `argv` comes from the profile. The agent cannot name the executable, `npm --prefix`, `go test -exec`/`-toolexec`/`-o`, `pytest -p`/`-c`/`--rootdir`, or any flag the command does not list.
- **`flags` is an allowlist of exact flag names**, each with a value class:
  - `bool`: no value;
  - `regex`: 1–256 bytes, no NUL or newline;
  - `int`: 0–10⁶;
  - `relpath`: `pathOk`;
  - `enum:[…]`.
  - `-f=v` and `-f v` are both normalized to `-f=v` before matching, so `-run` and `-run=X` cannot differ in meaning.
- **`positional`** must be one class:
  - `relpath`: `pathOk`, with no leading `-`;
  - `pkgpattern`: `relpath` optionally ending `/...`, or exactly `./...`;
  - `testfile`: `relpath` with a declared suffix set.
- **`--` is refused** unless the command declares `"passthrough": "--"`, in which case everything after it must still match `positional`. `npm test -- run <file>` (V22) is profiled as `{"argv":["npm","test","--","run"],"positional":"testfile"}`; the `--` sits in the prefix, not in agent args.
- **Env is never agent-controlled.** `GOFLAGS`, `NODE_OPTIONS`, `PYTHONSTARTUP` and `npm_config_*` come only from the profile and from World.

**What the grammar does not stop, and why that is acceptable.** Running a project's tests executes project code by definition. The vectors that remain are:
- the worktree's `package.json` scripts;
- `conftest.py`;
- `.npmrc` (`node-options --require` ran in every node process, V28);
- `go:generate` (not run by `go test`);
- test files themselves.

**The sandbox, not the allowlist, is the fence for code that runs.** The grammar's job is narrower:
1. keep the agent off srt's own options and off the shell;
2. keep the executable and its working set the operator's choice;
3. make every call's meaning an exact, recorded argv.

It is not presented as an execution fence (§9 R-140-3).

### 4.4 Rendered sandbox settings (per episode × project, `<state>/exec/<ep>.<project>.srt.json`, 0600)

| Key | Rendered value | Evidence |
|---|---|---|
| `filesystem.allowWrite` | `[<worktree>/<root>…<worktree>, <state>/exec-cache/<ep>/<project>]`. The **whole worktree** is writable, as with `ailang-write`, plus that episode's exec cache | V3 |
| `filesystem.denyWrite` | World's episode deny list, the same names as the AILANG policy: `.git` (dir **or** pointer file), `.github`, `.pi`, `.claude`, `.ailang`, `.gitmodules`, `.gitattributes`, each as an absolute path under the worktree. Plus srt's always-on paths `/tmp/claude`, `/private/tmp/claude`, `$HOME/.npm/_logs`, `$HOME/.claude/debug` | V4, V5, V6, V35 |
| `filesystem.denyRead` | D-140-1 = A: `[<operator $HOME>, <World state dir>, <workspace root>]` | V8, V19, V23, V24, V35 |
| `filesystem.allowRead` | `[<this episode's worktree>, <this episode's exec cache>, <profile read_roots>, <each read_root's realpath>]`, **plus on Linux only `<archive>/…/sandbox-runtime/vendor/seccomp/<arch>` [M0, V43, V44]**: srt runs that helper inside the sandbox. The dir is never in `allowWrite`, it is covered by the per-call archive digest, and nothing else of the archive is re-allowed. The workspace root is denied, so **sibling episodes' worktrees are unreadable** | V8, V24 |
| `network.allowedDomains` | `[]` (D-140-2 = A) | V9 |
| `network.allowLocalBinding`, `allowAllUnixSockets`, `enableWeakerNetworkIsolation`, `allowAppleEvents` | all `false`, rendered explicitly | V9, V20 |

- **Verification at render** (srt has no `summary` command, unlike AILANG's policy-tool):
  - World re-reads the file and checks it byte-for-byte against its own canonical rendering;
  - it records the sha256 as `settings_digest`;
  - it refuses a file it did not write.
- **The live check is the startup probe** (§4.6).

**Caches (D-140-3 = A):**
- `<state>/exec-cache/<ep>/<project>/` holds `home/` (the child's `HOME`), `tmp/` (the child's `TMPDIR`, set inside the sandbox, V40) and one directory per `caches` entry.
- Each cache directory is **seeded on first use by `cp -cR`** (APFS clonefile; on Linux, `cp -a --reflink=auto`) from the operator's read-only seed. The seed itself is never in `allowWrite` (V21).
- World always sets `PYTHONPYCACHEPREFIX` and `npm_config_cache` to directories in the episode cache, so the worktree gets no `__pycache__` (V24).
- `GOMODCACHE` (when set in the profile's env) must be a `read_root`; it is read-only by construction because it is not in `allowWrite` (V18).

### 4.5 Handler (`host/broker/handlers_exec.go`, `ExecHandler`)

```
<node> <archive>/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js --settings <rendered> -- \
  /bin/sh -c '/usr/bin/env -i "$@"; exit $?' world-exec <K=V…> <abs argv0> <prefix…> <normalized args…>
```

- **cwd** is `<worktree>/<root>`.
- **There are two environments, and only one is the child's [REVISED r3].**
  - **The host-side `node` that runs srt is not sandboxed.** Its env is exactly `PATH=/usr/bin:/bin`, and nothing from the profile reaches it. The measured hazard: `NODE_OPTIONS=--require ./evil.cjs` in that env ran worktree code on the HOST and wrote a marker into `$HOME` before any sandbox existed (V40). MUT-HOST-ENV is the killer.
  - **The child's env is set INSIDE the sandbox by `env -i`** in the trampoline. It contains exactly `<K=V…>` and nothing srt adds, so neither srt's `TMPDIR=/tmp/claude` nor its proxy variables survive (V41).
  - **`<K=V…>` contents:**
    - World's set: `HOME=<cache>/home`, `TMPDIR=<cache>/tmp`, `PATH` from the profile, `LANG=C.UTF-8`, `LC_ALL=C.UTF-8`, `PYTHONPYCACHEPREFIX`, `npm_config_cache`, `CI=1`, `NO_COLOR=1`;
    - plus the profile's `env` and the `caches` variables.
    - Names must match `^[A-Za-z_][A-Za-z0-9_]*$`, which is a load-time refusal.
  - **Temp dirs work inside the fence:** Go's `t.TempDir()`, Python's `tempfile` and node's `os.tmpdir()` all land in `<cache>/tmp` (V40).
  - **Agent args cannot become env assignments.** `env` stops reading assignments at the absolute `argv0`, so an agent arg `X=Y` after it is a plain argument (V41).
  - The env is built fresh for each call, and `childenv`'s credential arm is extended to both environments.
- **stdin** is `/dev/null`.
- **One process lifecycle, not two [REVISED r1].** `runBounded` stays the only subprocess lifecycle: `Setpgid`, group kill on cancel, `procbound.Admit`/`Wait`, the pipe-close grace and the cleanup errors (V30). It gains **a capture strategy** on `handlerCommand`:
  - **The zero value is today's strict capture, byte-for-byte:** stdout over the bound is an error and kills the group, and stderr goes to `boundedSink`. Every existing caller (the AILANG, git, model and publish handlers) and every existing `runBounded` test is unchanged, and the test diff adds arms without editing any. MUT-STRICT-REGRESSED is the killer.
  - **`captureHeadTail{head, tail, hardCap}`** is the only new behaviour, and `ExecHandler` is its only user. It is described in the Output bullet below.
  - The handler adds no `exec.Cmd`, `Setpgid`, kill or `procbound` call of its own. A source-scan arm in `handlers_exec_test.go` (like `TestProjection_NoDeadlineTampering`) refuses any of those identifiers in `handlers_exec.go`.
- **Deadline:** `min(profile.timeout_ms, handlerBudget(ctx))` (V30). On expiry, `runBounded`'s existing group kill fires, and the handler records `timed_out: true` with the partial head and tail.
- **Output (`captureHeadTail`, per stream):**
  - it keeps the first 8 KiB and the last 56 KiB, counts total bytes and hashes the full stream (sha256);
  - it never blocks the child;
  - above an operator hard cap (`--exec-max-output-bytes`, default 64 MiB; V27 saw 1.47 MB), it kills the group through the same path as today's overflow and records `limit: "output"`.
  - **The cap is sized to fit `MaxOutput`:** even in the worst case, where every byte JSON-escapes to six (`\u00XX`), 2 × 64 KiB × 6 = 768 KiB, under the 1 MiB coordinator cap with the envelope.
- **Exit status:** the trampoline makes a signal death 128 + n (V11). World also reads its own wait status for node (killed by World's deadline → `timed_out`).
- **Result** (handler output; the record's `output`):

```
{ "exit_code": 0, "timed_out": false, "limit": null, "duration_ms": 812,
  "stdout": "<head>…<tail>", "stdout_bytes": 1472271, "stdout_truncated": true, "stdout_sha256": "…",
  "stderr": "…", "stderr_bytes": 0, "stderr_truncated": false, "stderr_sha256": "…",
  "argv": ["go","test","-count=1","-run=TestLex","./internal/lexer/"],
  "profile": {"project":"ailang-compiler","digest":"sha256:…","command":"test","argv0_sha256":"…"},
  "sandbox": {"runtime":"@anthropic-ai/sandbox-runtime","version":"0.0.78","cli_sha256":"…","settings_digest":"sha256:…","read_fence":true,"network":"none"} }
```

- **rc 0 is reported, never judged.** It is the command's own status. A refused write inside a passing test still exits 0 (V25), exactly as row 134's refused inner effect does.

### 4.6 Operator configuration (`serve`) and the startup probe

| Flag | Default | Meaning |
|---|---|---|
| `--exec-profile <file>` (repeatable) | none | a project profile (§4.3). Its `project` names it; an episode selects one by `--exec-episode-project <ep>=<project>`, or by the only profile when there is one |
| `--exec-sandbox <dir>` | none (required with a profile) | an installed `@anthropic-ai/sandbox-runtime` package dir. Its `package.json` version must equal `ExecSandboxRelease` and its `dist/cli.js` sha256 must match the archived ref (the `ToolBinaryRelease` pattern, V1) |
| `--exec-node <path>` | `node` on the operator's PATH, resolved once | node ≥ 20.11 (srt's `engines`) |
| `--exec-max-output-bytes <n>` | 64 MiB | the runaway-output kill (§4.5) |

**Startup probe (fail closed) [REVISED r1, r2].** It runs once per profile, against **the profile's real rendered settings** with the probe's scratch episode substituted for the worktree. Any arm failing refuses `serve` with the arm's name. Two rules hold for every arm:
- **Each arm has a control** that proves it can fail.
- **Each "refused" arm demands a POSITIVE refusal token, never just a non-zero rc [r2].** The client in arms 1–4 and 6 is **World's own verified `node`** (the `--exec-node` binary, absolute path), running a fixed probe script that prints `WORLD-PROBE-<RESULT> <errno code>`. An arm passes only on its exact expected token (e.g. `WORLD-PROBE-REFUSED EPERM`). A missing or unexecutable client gives rc 126/127 and no token (V38), so it **fails** the arm instead of passing as a refusal.
- **[M0] Platform tokens (V45, V46).** On macOS a refusal is `EPERM`; on Linux a read-denied region is masked, so a refusal is `ENOENT` (`EROFS` for a `denyWrite` path). Every write arm also asserts that **the host bytes are unchanged**, which is the property on both platforms. On Linux a write into a masked region can print `WROTE` into a throwaway tmpfs. That is recorded, and it passes only if the host is unchanged. Arm 6 also gains a reverse leg: a listener the child binds must be unreachable from the host (`bind` succeeds inside the Linux netns).

1. **Write inside:** a file is written, and World reads it back unsandboxed.
2. **Write outside:** a write to a sibling scratch episode under the workspace root is refused. The control: the file is absent unsandboxed.
3. **Read fence, with fixed decoys** that World plants unsandboxed before the probe and that do not depend on the profile:
   - `$HOME/.ailang-worldd-exec-probe-decoy`;
   - `<state>/exec-probe/decoy`;
   - `<workspace root>/.exec-probe-sibling/decoy`.
   - A sandboxed `cat` of each must be refused. The control: each one reads back unsandboxed. The `read_roots` load rules (§4.3) guarantee that no read root can contain these paths, so a profile cannot make the probe adapt to it.
4. **srt's default write paths:** a `/tmp/claude` write is refused.
5. **Exit status:** the trampoline maps a TERM self-kill to 143.
6. **Network, against a live listener:** World binds its own `127.0.0.1:0` listener and counts accepts.
   - **Control:** an unsandboxed connect is accepted (count 1), which proves the listener is live.
   - **Arm:** a sandboxed connect to that port, both raw and through the proxy env, must fail, **and** the count must still be 1 afterwards (V37).
   - So a connect refused by an unbound port cannot pass as a refusal by the sandbox. A Linux host whose network isolation silently fails (e.g. missing `socat`) shows as accept count 2 and is refused.
7. **Toolchain present:** the profile's declared `probe` argv (default `argv0 --version`) → rc 0.

This is the Linux gate as well. A host where bwrap or user namespaces are missing fails arm 1 or 2; one whose network isolation is missing fails arm 6. In either case the tool is refused, never run unsandboxed.

**Integrity of the confinement stack [REVISED r2].** The fence is only as good as the srt, node and profile bytes that set it up, so each is placed out of reach and re-checked.
- **Placement (load-time refusal):** the `--exec-sandbox` dir, the `--exec-node` binary (and the realpath of its directory and of its `libnode` dylib), every `--exec-profile` file and every cache seed must resolve, by realpath, **outside** the workspace root, the state dir and every rendered `allowWrite` root. Sandboxed code can write only inside `allowWrite` (V3, V35), and World's own write tools only inside the worktree, so nothing an agent or its project code can reach can modify them.
- **Archive and re-verify before every call** (the `verifiedTool` pattern, V31):
  - At startup, World copies **the whole install tree** into `<state>/archive/exec-sandbox/<digest>/node_modules/`: srt plus its runtime dependencies `commander`, `zod`, `node-forge` and `@pondwader/socks5-server`, 829 files and 13.8 MB. The package dir alone does not run (`ERR_MODULE_NOT_FOUND: commander`, V42). World runs srt **from the archived copy**.
  - **The state dir being read-denied in the sandbox does not stop this [r3].** srt's own `node` runs outside the sandbox and reads the archive normally. The sandboxed child cannot read it, which is a bonus (V42).
  - The archived tree's manifest digest takes 23.8 ms (V42). At startup, World also records the sha256 of `node` and `libnode` (0.2 ms + 32.2 ms, V39).
  - Before **each** call, World re-verifies the archived srt tree, `node`, `libnode` and `argv0`. A mismatch is a handler failure (recorded `failed`) with zero spawns, never a run.
  - Other dylibs `node` loads are not hashed. They are covered by the placement rule alone (R-140-11).
- **Profile:** read **once** at `serve` into memory and digested, never re-read. Editing the file takes effect only on restart, and the record's `profile.digest` names what ran.

### 4.7 Replay

- The request object holds `{command, args}`, and the record holds the full result: argv, digests and bounded streams. Replay returns the recorded output and **never re-executes** (row-134 V10/V46).
- A profile, binary or sandbox drift shows as a digest mismatch between records. It does not change history.
- The truncated middle of a stream is not recoverable from the record. Its sha256 lets anyone who kept the full output prove which stream it was, but World does not keep it in v1 (R-140-6).

## 5. Premises

| P | Premise | Evidence |
|---|---|---|
| P1 | No World tool can build or test non-AILANG code today | V31 |
| P2 | srt confines writes and network for the whole process tree, orphans included | V3, V9, V13 |
| P3 | Reads are open by default and fenceable by path list; denying `/tmp` breaks node, so `/tmp` stays readable | V7, V8, V8b |
| P4 | srt's CLI needs `--` after its options and an exit-status trampoline | V2, V11, V12, V32 |
| P5 | srt allows four out-of-tree write paths by default; `denyWrite` closes them, case-insensitively on APFS | V4, V6 |
| P6 | Go, TS and Python tests pass confined with World-owned caches; offline Go works from a read-only module cache that holds the toolchain | V17–V19, V22–V24 |
| P7 | Go cannot fetch modules through srt on macOS without the `trustd` exception; npm and uv can | V20, V29 |
| P8 | A cloned seed makes a per-episode cache warm in < 0.4 s | V21 |
| P9 | srt overhead ≈ 0.15 s per call; the frozen 10 s handler cap fits targeted commands but not full suites | V15, V18, V22, V27, V30 |
| P10 | The sandbox caused none of the test failures seen | V27 |
| P11 | Argument injection that reaches srt's options is an escape; injection into a toolchain's flags only runs code inside the fence | V12, V28 |
| P12 | ~~Linux behaviour is unmeasured~~ **[M0] Linux is measured: with the seccomp helper dir re-allowed, every gating arm passes; the design as written fails every arm** | ~~V33~~ V43, V44 |
| P19 | **[M0]** On Linux a read-denied region is masked, not refused: reads → `ENOENT`, writes land in a throwaway tmpfs. **The confinement property is "host bytes unchanged"**, and probe and AC arms assert that, with the token recorded | V45 |
| P20 | **[M0]** On Linux the network fence is a namespace: a child may `bind`, but nothing it binds is reachable from the host, and it reaches no host listener | V46 |
| P13 | An `allowRead` equal to a denied root re-opens it (on macOS a strict ancestor or a symlink does not), so World must refuse such `read_roots` | V36 |
| P14 | A live-listener probe tells a sandbox refusal apart from an unbound port | V37 |
| P15 | World's own `node` is a probe client that yields a positive refusal token; a missing client yields rc 127 and no token | V38 |
| P16 | Per-call re-verification of the srt install tree, `node` and `libnode` costs ≈ 56 ms | V39, V42 |
| P17 | The host-side node's env is a sandbox-escape surface; the child env must be set inside the sandbox | V40, V41 |
| P18 | srt runs from an archived whole install tree under a read-denied state dir | V42 |

## 6. Milestones

**M0 (0.5 d) — pin and Linux measurement (no production code until it passes).**
- AC0.1: re-run the banked matrix on the chosen pin (default **0.0.78**; fall back to 0.0.71 on any arm difference) and record `ExecSandboxRelease` plus `cli.js` sha256.
- AC0.2: a CI job step on `ubuntu-latest` must:
  - install `bubblewrap socat ripgrep`;
  - set `kernel.apparmor_restrict_unprivileged_userns=0`;
  - run the startup-probe arms 1–6 plus V3/V8/V9/V11/V12/V13 against srt and record the results as V-rows.
- AC0.3: if any Linux arm differs, it is written into this doc as a new V-row and an amended premise before M1 starts.
- **M0 DONE 2026-10-05 (attended, PR #206):** pin srt **0.0.78** (`cli.js` sha256 `3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2c7cd40d96e`, `world-row140-m0/pin.json`). Linux differed, so V43–V49 and P19/P20 were added, and §4.4 (the seccomp helper read) and §4.6 (probe semantics on Linux) were amended. A step in `ci.yml`'s go job runs the matrix on `ubuntu-latest` (a second workflow file is refused by `TestGoToolchainPinsAgreeAndMatchJobList`).

**M1 (1 d) — transition, descriptor, sketch.**
- AC1.1: `exec.ail` named tests pin the plan bytes for:
  - `{command:"test"}`;
  - `{command:"test", args:["-run=X","./a/..."]}`;
  - refusals: a bad id (`Test`, `-x`, 33 chars), 17 args, a 513-byte arg, a NUL, `args` as a string, an unknown key.
- AC1.2: `execargs.ail` carries the argument-law contracts, and they verify; `verify_ail.sh` is green with the new identities.
- AC1.3: `tools/list` lists 9 tools; `workspace-exec` declares 1 triple with `additionalProperties:false`; the workspace registry binds 9 names.

**M2 (1 d) — handler, rendering, caches.**
- AC2.1: the rendered settings' exact key set; a golden for one profile.
- AC2.2 (real srt): `go test` on a scratch module → `ok`, `exit_code 0`. A failing test → non-zero `exit_code`, with the failure text in the tail.
- AC2.3: a test that SIGTERMs itself → `exit_code 143` (killer of MUT-NO-TRAMPOLINE).
- AC2.4: each of these is refused with an exec counter of 0 (no srt spawn):
  - an unprofiled command;
  - `-exec=sh`;
  - `--`;
  - `-c`;
  - `-s=x.json`;
  - `../x`;
  - `/abs`;
  - a flag value with a newline;
  - an arg that is exactly `--settings`.
- AC2.5: the confinement matrix through the handler:
  - writes to a sibling episode, `$HOME`, `/tmp/claude` and `.git/HEAD` → refused;
  - so are `.GITHUB/x` and the worktree `.git` pointer file;
  - reads of a scratch decoy under the operator `HOME` and of a sibling episode's file → refused;
  - the bytes are unchanged.
- AC2.6: 2 MB of output → `stdout_truncated:true`, `stdout_bytes` = 2 000 000, the sha256 equals the full stream's, and the result is < 1 MiB. An all-control-character stream of the maximum recorded size still fits under `MaxOutput`.
- AC2.7: a command sleeping past `timeout_ms` → `timed_out:true`, no member left in the group, and inside the handler budget.
- AC2.8: the cache seed is byte-unchanged after a run that writes to `GOCACHE`; `PYTHONPYCACHEPREFIX` keeps `__pycache__` out of the worktree.
- AC2.10 [r3]: a profile `env` of `NODE_OPTIONS=--require=./zz.cjs`, where the worktree's `zz.cjs` writes a marker under a scratch `HOME`-like path outside `allowWrite`. The marker stays absent; the host-side node's recorded env is exactly `PATH=/usr/bin:/bin`; the child's `/usr/bin/env` output equals the rendered `K=V` set.
- AC2.9 [r2]: after startup, flip one byte of the archived `cli.js` (and, separately, swap `argv0`). The next call is recorded `failed`, the spawn counter stays 0, and the error names the file.

**M3 (0.75 d) — wiring and docs.**
- AC3.1: the startup refusal table, which `serve --help` stays bound to QUICKSTART for:
  - §4.3's load-time refusals, including the `read_roots` rows `$HOME`, an ancestor of the workspace root, a symlink to `$HOME`, `/` and a probe decoy;
  - §4.6's flags;
  - every probe arm by name, each driven red by a pass-through srt shim that runs the command unsandboxed;
  - a missing probe client failing arms 2, 3 and 6 [r2];
  - placement rows: `--exec-sandbox`, `--exec-node` and `--exec-profile` inside the workspace root, inside the state dir and inside an exec cache, each refused [r2].
- AC3.2: an `/mcp/` end-to-end `workspace-exec` call that commits, with `world.effects[0].record` naming `Workspace.Exec` and the request holding `{command,args}`.
- AC3.3: no grant → `denied`, no spawn; budget 1 → the second call is `denied:budget`.
- AC3.4: a nil-registry replay is byte-equal.
- AC3.5: QUICKSTART gains §"Build and test non-AILANG projects": profile authoring, seeding caches, the three worked profiles (Go, TS, Python).

**M4 (attended, 0.5 d).** Mark republishes `se-tools`. Then, through a real MCP client:
- the ailang compiler: `test` on one package;
- TwilightGame: `typecheck`;
- sunholo-platform/cli: `test-file`;
- each returns its measured result; an unprofiled command is refused; a decoy read inside a test is refused.

## 7. Load-bearing mutations

| Mutant | Killer |
|---|---|
| MUT-NO-DASHDASH: srt invoked without `--` after `--settings` | AC2.4 `-s=x.json` / `-c` arms: the stub records srt parsing them. Real srt: V12's settings replacement |
| MUT-NO-TRAMPOLINE: argv passed to srt directly | AC2.3 (TERM → 0 instead of 143) |
| MUT-SHELL-JOIN: args joined into one string | AC2.4 newline/`$(…)` arms change meaning; AC2.2 argv echo |
| MUT-GRAMMAR-SKIP: the handler skips the flag allowlist | AC2.4 `-exec=sh` |
| MUT-PREFIX-FROM-AGENT: `argv[0]` taken from args | AC1.1 + AC2.4 |
| MUT-ENV-INHERIT: `os.Environ()` given to the child | the `childenv` credential arm extended to exec (a planted `AILANG_REGISTRY_API_KEY` must read empty) |
| MUT-NO-DENYWRITE-DEFAULTS: srt's four always-writable paths not denied | AC2.5 `/tmp/claude` |
| MUT-NO-WORLD-DENY: World's deny list not rendered | AC2.5 `.git/HEAD`, `.GITHUB/x` |
| MUT-READ-FENCE-OFF: `denyRead` empty | AC2.5 decoy and sibling-episode reads |
| MUT-ROOT-NOT-DENIED: workspace root left out of `denyRead` | AC2.5 sibling-episode read |
| MUT-NET-ON: `allowedDomains` from the profile | startup probe arm 6 (live-listener accept count becomes 2) |
| MUT-PROBE-NET-VACUOUS: arm 6 without its live listener and control | AC3.1 (srt replaced by a pass-through shim that runs the command UNsandboxed must be refused by arm 6, not only by arms 1–3) |
| MUT-READROOT-ANCESTOR: the `read_roots` load rule drops the equal/ancestor/realpath clauses | AC3.1 rows `read_roots: [$HOME]` and `read_roots: [<symlink → $HOME>]` refused at load |
| MUT-STRICT-REGRESSED: `runBounded`'s zero-value capture changes | every pre-existing `runBounded`/handler test (unedited) goes red, e.g. the 8 MiB overflow-is-an-error arms |
| MUT-SECOND-LIFECYCLE: `handlers_exec.go` spawns its own `exec.Cmd` | the `handlers_exec_test.go` source-scan arm |
| MUT-SEED-WRITABLE: the seed in `allowWrite` | AC2.8 |
| MUT-OVERFLOW-ERROR: the exec handler uses the strict capture | AC2.6 (a 2 MB run must succeed truncated) |
| MUT-NO-GROUP-KILL: only node killed on timeout | AC2.7 (group members remain) |
| MUT-PROBE-RC-ONLY: a probe arm accepts any non-zero rc as "refused" | AC3.1 (the probe client replaced by a missing path must FAIL arms 2, 3 and 6, not pass them) |
| MUT-STACK-UNVERIFIED: the per-call re-verification skipped | AC2.9 |
| MUT-PLACEMENT: a `--exec-sandbox` dir inside the workspace root accepted | AC3.1 placement rows |
| MUT-HOST-ENV: the profile env (or any child var) is given to the host-side `node` | AC2.10 |
| MUT-TRAMPOLINE-EXEC: `exec env -i` loses the signal mapping | AC2.3 (TERM → 143) |
| MUT-ARCHIVE-PKG-ONLY: only the srt package dir is archived | AC3.1 startup (srt fails to load: `ERR_MODULE_NOT_FOUND`) |
| MUT-PROBE-SKIPPED: `serve` enables exec without the probe | AC3.1 (a sandbox dir whose srt is replaced by a no-op shim must be refused) |

## 8. Decisions for Mark (each with a recommended default)

**RULED 2026-10-03 (Mark, attended; D-WORLD-62): 1 = A (fence reads), 2 = A (no network in v1), 3 = A (per-episode seeded caches), 4 = A (inside the frozen D7 deadlines; follow-up row for long runs), 5 = B (Mark reviews the r3 delta attended: the host/child env split, the in-sandbox `TMPDIR`, the whole-tree archive).**

- **D-140-1 Reads.**
  - **A (recommended): fence reads by default.** `denyRead` = the operator `$HOME`, World's state dir and the workspace root. `allowRead` = this worktree, this episode's cache and the profile's toolchain roots. Measured working for Go, TS and Python (V19, V23, V24). It corrects the row: reads **are** fenceable, though weaker than the AILANG tools (system paths and `/tmp` stay readable, V8b).
  - B: no read fence (the row's assumption).
- **D-140-2 Network.**
  - **A (recommended): none in v1.** Dependencies are provisioned by the operator before the session (`npm ci`, `uv sync`, a warm Go module cache), and exec runs offline (V18, V22, V24). An episode that changes dependencies asks the operator.
  - B: a second effect `Workspace.ExecNet`, per the D-135-3 precedent (a separate grant). It serves profile commands marked `network: [hosts]`, for lockfile installs only (`npm ci`, `uv sync --frozen`, both measured through srt's proxy, V29). Go module fetches would still fail on macOS unless `enableWeakerNetworkIsolation` opens `trustd` (V20), which this design would never render.
- **D-140-3 Cache scope.**
  - **A (recommended): per-episode caches seeded by clonefile** from an operator-prewarmed, read-only seed. 0.37 s to seed; `go test` 1.2 s warm vs 5.6 s cold (V21). No episode can poison another's build cache.
  - B: one writable cache per project, shared by episodes. Fastest, but an episode can plant build-cache entries a later episode trusts.
  - C: per-episode and cold. Simplest, but it alone can exceed the 10 s cap.
- **D-140-4 Time budget.**
  - **A (recommended): v1 lives inside the frozen D7 deadlines** (handler ≤ 10 s, `timeout_ms` ≤ 9 s). Targeted commands fit (V18, V22, V24). Full suites do not (V27), so a follow-up row is filed: an asynchronous long-run lane (start effect + result transition, each recorded).
  - B: raise D7 for this tool. That is a ratification-class change to frozen deadlines guarded by `TestProjection_NoDeadlineTampering`, and MCP clients' own tool timeouts become the ceiling.

- **D-140-5 Quorum closure.**
  - **A (recommended): one more quorum round on the r3 delta before sprint planning** (≈ $0.40). Rounds 1–3 each found a real defect, so a fourth look at the newest, unreviewed text is worth its cost.
  - B: Mark reviews the r3 delta attended and waives the round.

Designer's calls (not asked): the srt CLI with `--` and the trampoline, rather than a World-owned node launcher (both CLI hazards are closed and pinned by mutants; upstream issues are PROPOSED, not yet filed, for the TERM → 0 mapping and the option parsing); a JSON profile (no new Go dependency); the tool name `workspace-exec`.

## 9. Risks and residuals

- **R-140-1 The read fence is a path list, not a jail.**
  - **Stays readable:** `/usr`, `/opt`, `/etc`, `/Library`, `/System`, `/private/tmp`, `/private/var/folders` and other users' homes, unless the operator lists them (V8, V8b), plus the profile's `read_roots` (e.g. the whole Go module cache, which is public code).
  - **Denied:** the operator's `$HOME` except the re-allowed roots, World's state and the workspace root (sibling episodes).
  - This is weaker than the AILANG tools' read confinement, and it is stated, not hidden.
  - **Linux has no globs** (V33), and the rendering uses only literal paths for that reason.
- **R-140-2 srt is a "beta research preview"** (README), and its behaviour has moved between releases. Mitigation:
  - the pin plus a hash;
  - the startup probe on every `serve`;
  - M0's re-run of the matrix on any bump (the R-SE-8 re-proof pattern).
- **R-140-3 Project code runs.** Tests, `conftest.py`, `package.json` scripts and `.npmrc` `node-options` execute by definition (V28). The fence bounds what that code can do; nothing bounds what it computes. A malicious dependency in a worktree can read every allowed root and write anywhere in the worktree.
- **R-140-4 Inherited output descriptors are checked by Seatbelt.** `cat: stdout: Operation not permitted` appeared when stdout was a FILE under a denied path (`p1_fs.out` V8 arm). World always passes pipes, so the risk is to a future change.
- **R-140-5 srt's proxy env (V10) is CLOSED [r3].** The trampoline's `env -i` replaces srt's env with World's exact set (V41). Under D-140-2 = B, the proxy variables would have to be passed through on purpose.
- **R-140-6 The middle of long output is dropped.** Head 8 KiB + tail 56 KiB per stream, with total and sha256. A full-output object in the store is a later option (it would need a store size policy).
- **R-140-7 Dependencies drift from the worktree.** Under D-140-2 = A, an agent that edits `go.mod`/`package.json` gets a refusal-like failure (`missing go.sum entry`, V20) until the operator installs. Owner: D-140-2 B or a later row.
- **R-140-8 Linux is unmeasured** (V33). **[M0: CLOSED. GitHub's runner runs bwrap, and the matrix passes there (V44). The real-sandbox ACs run in CI from M2.]** M0 gates everything on it. If GitHub's runner cannot run bwrap, CI keeps only the stub arms, and the real-sandbox arms run on the rig and in M4. That is a degraded gate, said so in the sprint plan.
- **R-140-9 The handler cap bounds what fits** (D-140-4): a cold cache plus a large package can still time out. The `timed_out` result is typed, and the agent sees it.
- **R-140-11 node's other dylibs are trusted by placement, not by hash** (§4.6 Integrity). They sit in operator-owned prefixes (e.g. `/opt/homebrew`) that no sandboxed process can write (V3).
- **R-140-10 `node` becomes a World runtime dependency** for this tool only. Row 138's `setup`/`doctor` should fetch and check it; that is noted for row 138's follow-up, not added here.

## 10. Conflict surface

**Touches:**
- `packages/se-tools/se_tools/exec.ail` and `transitions.json` (republished attended);
- `design_docs/sketches/execargs.ail`;
- `host/broker/handlers_exec.go` (new) and its tests;
- `host/broker/handlers.go` (`runBounded` gains a capture strategy whose zero value is today's strict capture, pinned by the unedited existing tests; the `headTailSink` beside `boundedSink`) [REVISED r1];
- `host/daemon/workspace.go` (9 effect names; the exec handler per episode; settings rendering);
- `cmd/ailang-worldd/main.go` (4 flags);
- `.github/workflows/ci.yml` (M0 Linux step);
- `scripts/verify_ail.sh` (new identities);
- `docs/QUICKSTART.md`;
- `host/childenv` tests (the credential arm).

**Shares with row 141** (`w-workspace-project-layouts`): the "module root inside the worktree" idea. Here it is the profile's `root`; row 141 adds it for the AILANG tools. Whichever lands second reuses the other's predicate.

**Does not touch:**
- `effectplan.go`;
- the broker decision law;
- the projection or its D7 deadlines;
- the store;
- the coordinator (`HandlerCap` unchanged);
- the AILANG policy rendering;
- `tools/launchd/*`.

## 11. Non-goals (restated)
- No Docker or containers.
- No general shell, and no agent-chosen executable, env or cwd.
- No network in v1 (D-140-2).
- No runs past the D7 deadlines (D-140-4).
- No Windows.
- No AILANG-tool layouts (row 141).

## 12. Quorum log

`ailang design-quorum --author claude:claude-opus-5-5` (quorum CLI `AILANG v0.52.0-17-g790169359-dirty`, used only as the quorum runner). Seats: `gemini-3-1-pro`, `oc-glm-5-3`, `oc-kimi-k3`; reserve `gpt6-1-sol`; `claude-sonnet-5@claude-p` sits out as the author's vendor. Artifacts: `design_docs/verification/world-row140-design-2026-10-03/quorum/`.

- **Round 1 — 2026-10-03T18:33:14Z, BLOCKED 3/3** ($0.284). Each premise was measured or read before revising:
  - **`gemini-3-1-pro`: §4.5 duplicated `runBounded`'s lifecycle while §10 said `runBounded` was unchanged. REAL.** The draft implied a second `exec.Cmd` lifecycle beside `runBounded`'s Setpgid, group kill, `procbound` and pipe grace (V30). Revised: `runBounded` gains a capture strategy whose zero value is today's capture byte-for-byte, and the exec handler owns no lifecycle code (a source-scan arm). New mutants MUT-STRICT-REGRESSED and MUT-SECOND-LIFECYCLE.
  - **`oc-glm-5-3`: a `read_roots` entry that contains a denied root re-opens it, and probe arm 3's decoy adapts to the profile. REAL, with a measured nuance (V36).**
    - An `allowRead` **equal** to `$HOME` re-opens the decoy.
    - A strict ancestor (`/Users`) and a symlink to `$HOME` do **not** on 0.0.71/macOS. But World's own realpath expansion would turn the symlink case into the equal case, and Linux is unmeasured.
    - Revised: the load rule refuses equal, ancestor, realpath, `/` and decoy-containing roots (§4.3). Probe arm 3 uses three fixed, profile-independent decoys with unsandboxed controls (§4.6). New mutant MUT-READROOT-ANCESTOR.
  - **`oc-kimi-k3`: probe arm 6 could pass against an unbound port, because a refused connect and a sandbox refusal look alike. REAL** (the draft named "World's own port" without ordering it against `serve`'s bind). Measured the fix (V37): a World-owned live listener, an unsandboxed control (accepts = 1), then the sandboxed raw and proxy connects, refused, with accepts still 1. Revised arm 6; new mutant MUT-PROBE-NET-VACUOUS.
- **Round 2 — 2026-10-03T18:36:55Z, BLOCKED 3/3** ($0.289; same seats). Each premise was measured or read before revising:
  - **`gemini-3-1-pro`: probe arm 6 passes vacuously when the network client is missing (rc 126/127 satisfies "must fail" and "count stays 1"). REAL** (measured V38: a missing client is rc 127 with no output). Revised: every "refused" arm demands a positive token from World's own verified `node` (`WORLD-PROBE-REFUSED EPERM`, V38). A missing client therefore fails the arm. New mutant MUT-PROBE-RC-ONLY.
  - **`oc-glm-5-3`: the confinement stack (srt, node, profile) has no placement rule and is hash-checked only at load; project code could rewrite `cli.js`. PARTLY REFUTED, PARTLY REAL.**
    - *Refuted:* sandboxed project code can write only inside `allowWrite` (V3, V35), so with srt outside the worktree it cannot rewrite `cli.js`.
    - *Real:* nothing in the design enforced that placement, and an operator could install srt inside a worktree.
    - Revised (§4.6 Integrity): realpath placement refusals for srt, node, libnode, profiles and seeds; srt runs from World's archive; srt, node, libnode and `argv0` are re-verified before every call, ≈ 48 ms (V39); the profile is read once into memory. New AC2.9; mutants MUT-STACK-UNVERIFIED and MUT-PLACEMENT; residual R-140-11.
  - **`oc-kimi-k3`: one BLOCKED round with self-certified fixes is not a cleared quorum. Process objection; its premise is now false.** Round 2 *was* the ratifying review of the r1 revisions, and none of its three objections reopened an r1 fix. The Status line now says plainly that the doc is executor-ready only after a PROCEED or a ruling by Mark.
- **Round 3 — 2026-10-03T18:42:12Z, BLOCKED 3/3** ($0.383; same seats). Each premise was measured before revising:
  - **`gemini-3-1-pro`: the profile env reaches the host-side node that runs srt, so `NODE_OPTIONS` could run code before the sandbox exists. REAL, and confirmed as an escape** (V40: `--require ./evil.cjs` in the host env wrote a marker into `$HOME`). Revised §4.5: the host node's env is `PATH` only; the child env is set inside the sandbox by `env -i` in a trampoline that keeps the signal mapping (V41). New AC2.10; mutants MUT-HOST-ENV and MUT-TRAMPOLINE-EXEC. This also closes R-140-5 (srt's proxy variables no longer reach the child).
  - **`oc-glm-5-3`: the child's `TMPDIR` route through `CLAUDE_CODE_TMPDIR` is asserted, not measured, and temp-writing tools may fail. PARTLY REFUTED, superseded.**
    - *Refuted:* `p3_paths.out`'s last arm did measure it (`TMPDIR=<epcache>/tmp`), though no V-row cited it.
    - *Superseded:* the r3 design sets `TMPDIR` directly inside the sandbox.
    - Go `t.TempDir()`, Python `tempfile` and node `os.tmpdir()` were measured working there (V40).
  - **`oc-kimi-k3`: an archive under the read-denied state dir would stop srt from reading `cli.js`. PREMISE REFUTED, but the probe found a real gap.**
    - *Refuted:* srt's own node runs outside the sandbox and reads the archive normally (V42, rc 0); only the sandboxed child is denied, which is desired.
    - *Real gap:* archiving the package dir alone fails (`ERR_MODULE_NOT_FOUND: commander`). The archive is now the whole install tree, and its per-call digest is 23.8 ms. New mutant MUT-ARCHIVE-PKG-ONLY.
    - *Process half (self-certified fixes):* noted. Three rounds is the brief's cap, so the r3 revisions are unreviewed by quorum and go to Mark with this log (§8).


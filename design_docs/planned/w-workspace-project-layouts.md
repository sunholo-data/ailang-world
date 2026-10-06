# w-workspace-project-layouts — a module-root subdirectory sandbox and a read-only registry package cache (row 141)

**Status:** PLANNED. Designed in iteration 238 (designer lane claude-opus-5-5, unattended). **Quorum round 1 BLOCKED 3/3** (gemini-3-1-pro, oc-glm-5-3, oc-kimi-k3). None disputed the direction; all raised completeness or premise defects. **This is r2**, the single revision pass, applying the five r1 fixes listed under "r2 changes" at the end of the doc.
**Target:** next World patch. **Priority:** clause-2/4 (dogfood, D-WORLD-58). Charter position: after 136, before 93 (D-WORLD-61 = A).
**Estimated:** ~1 d of host Go and docs (M1 ~0.35 d, M2 ~0.4 d, M3 ~0.25 d). **Dependencies:** none. Everything runs on the pinned tool binary, AILANG v0.52.1. Two upstream asks (M3) are hardening only and gate nothing in this row.
**Measured base:** worktree `sprint/row141-workspace-project-layouts` at `499105b`. Tool binary `~/.pinned-ailang-tools/v0.52.1/ailang` (V21). Upstream source read at tag `v0.52.1` (`c68ded4`) in `~/dev/sunholo-data/ailang`. Daneel read-only at `ebe26c0`. Every probe ran on a `git archive` copy under `~/.ailang/state/world-iter238-row141/`.

## Problem

`serve --workspace-root W` binds episode E to the worktree `W/E` and makes that directory the AILANG sandbox. It renders the policy with `fs_sandbox = W/E` and runs every tool with cwd `W/E` (V1). Real AILANG projects do not all fit that layout. The first dogfood (Daneel #201, commit `3b2162d`) hit two problems:

- **(a) Module root.** Daneel's modules live in `tools/` and import each other by bare name (`import daneel_reading`). AILANG resolves a bare import against the process cwd (V5). policy-tool always runs its CLI child in the sandbox root (V4), so with the sandbox at the worktree root, `ai_check tools/daneel_brief.ail` fails `LDR001 module not found: daneel_reading`. Starting policy-tool from `tools/` changes nothing (V7, V8). With the sandbox rooted at `tools/` the same module passes with 6 contracts verified (V7). The dogfood got around this by making the episode id `tools`. World has no real option for it (V3).
- **(b) Registry packages.** 12 of Daneel's `tools/` modules import registry packages (`pkg/sunholo/oauth/token`, …). Two things block them:
  - **The lock file.** Daneel's `tools/ailang.lock` is gitignored, so a fresh worktree has none, and the agent cannot run `lock` (V9, V12).
  - **The package cache.** With a lock present, v0.52.1 looks registry packages up under `$HOME/.ailang/cache/registry/<ns>/<name>/<ver>`. Under World, `HOME` is the empty per-episode cache dir (V2, V10, V11).

  `AILANG_CACHE_DIR` does not cover registry packages (V11), and v0.52.1 has no policy key for a package path (V6).

There is a second finding under (b), and it is a defect today. The admitted `pkg_docs` CLI op makes an unconfined child that **fetches from the registry over the network and writes the package into the per-episode `HOME`**. A later `ai_check` then resolves imports against that agent-fetched copy (V13). So the per-episode package cache is already agent-writable, and it is filled from the network through a "documentation" op.

## Measured facts

Each row is backed by the Verification Log (§ Verification Log) entry of the same id.

| # | Claim |
|---|---|
| V1 | `episodeHandler` sets `cacheDir := <state>/cache/E` (`host/daemon/workspace.go:380`). It renders `broker.RenderEpisodePolicy(epRoot)` (`:387`) and builds the handler with `Root: epRoot, CacheDir: cacheDir` (`:397`). `episodeRoot` requires `EvalSymlinks(W/E) == W/E` and a directory (`:355–368`). |
| V2 | Every tool child gets `HOME=<cacheDir>` and `AILANG_CACHE_DIR=<cacheDir>`, with `PATH=/usr/bin:/bin` (`host/broker/handlers_ailang.go:522–533`). Each result is stamped with `tool` and `policy_digest` only (`:496–500`). |
| V3 | No module-root or package-cache option exists. The (valid-ERE) grep over non-test Go in `host/` and `cmd/` returns exactly 1 hit, `host/pinfetch/pinfetch.go:183` ("release subdirectory", unrelated). The positive control (`AILANG_CACHE_DIR`) returns 3. |
| V4 | policy-tool (v0.52.1) runs every `needsRoot` CLI op with `cmd.Dir = h.root.Dir()`, the sandbox root (`internal/policytool/cli_ops.go:316`). The child is "a separate, UNCONFINED process" (`:156`). Its `AILANG_CACHE_DIR` is overridden with a private temp dir (`:405`). Its `HOME` is inherited from World. |
| V5 | The module loader's base dir is `"."`, the cwd (`internal/pipeline/pipeline_module_phases.go:153`). |
| V6 | The v0.52.1 operator policy has 16 TOML keys. None is a module root or a package path. `cli_allow` exists (`internal/policy/policy.go:43,51`). |
| V7 | On a Daneel copy with World's exact child env and rendered policy: sandbox at the root plus `tools/daneel_brief.ail` → `LDR001 module not found: daneel_reading`. Sandbox at `tools/` plus `daneel_brief.ail` → `passed`, 6 verified. |
| V8 | Starting policy-tool with cwd `tools/` under the root-sandbox policy still gives `LDR001` (the child's cwd is reset to the sandbox root, V4). |
| V9 | Daneel's `tools/ailang.lock` is ignored by `.gitignore:16` and is not tracked. Without it, every module that imports `pkg/…` fails with `package "…" not found in ailang.lock`. This happens even unconfined with the operator's real `HOME`. |
| V10 | With the lock copied in and an empty confined `HOME`, path dependencies resolve, even though `../ext/*` lies **outside** the `tools/` sandbox (`daneel_ext` passes). Registry imports fail with `registry package sunholo/oauth cache not found at <HOME>/.ailang/cache/registry/sunholo/oauth/0.1.0`. The child also created `<HOME>/.ailang/cache/registry`. |
| V11 | The registry cache is `os.UserHomeDir()+"/.ailang/cache/registry"`, created with `MkdirAll(…, 0755)` (`internal/pkg/registry.go:213–214`). A `registry` lock entry resolves there, and falls back to the lock's `path` when that directory is missing (`internal/pkg/loader.go:171–175`). `AILANG_CACHE_DIR` is documented as the compile and prompt cache only (`internal/config/paths.go:21`). |
| V12 | Under World's rendered policy, the summary's `cli` list has 20 ops. `pkg_docs` is one of them; `lock` and `design_quorum` are not. |
| V13 | `pkg_docs` (`cmd/ailang/pkg_docs.go:71–90`) fetches the registry index and the package tarball when the package is not cached (`:73`), and extracts it into `$HOME/.ailang/cache/registry` (`:88`). Measured under World's env: `pkg_docs sunholo/oauth` returned `ok:true`, and `<HOME>/.ailang/cache/registry/sunholo/oauth/0.1.0/` appeared. `ai_check daneel_sources.ail` in that `HOME` then **passed**, with 4 verified. |
| V14 | Read-only snapshot test: the 8 locked registry packages were copied into a tree and made `chmod -R a-w` (0 writable entries), and `<HOME>/.ailang/cache/registry` was a symlink to that tree. **All 12** `pkg/`-importing Daneel modules then pass `ai_check`. `pkg_docs` for an absent package fails (`has no AGENT.md`), and the `find`-listing digest of the snapshot plus `HOME` is unchanged before and after. |
| V15 | `ailang run --policy` on a module that imports `pkg/sunholo/oauth/token` fails with an empty `HOME` (`cache not found`) and prints `pkg-ok` with the read-only `HOME`. A run that calls `writeFile` on a path inside the snapshot is refused with `escapes sandbox`, and nothing is written. |
| V16 | An offline synthetic fixture shows all three arms. It has a project in `ep1/tools` (`a` imports bare `b` and `pkg/acme/util/greet`), a generated lock with one `registry` entry, and a read-only snapshot. Root sandbox → `LDR001 module not found: b`. `tools/` sandbox with an empty `HOME` → `cache not found`. `tools/` sandbox with the snapshot `HOME` → `passed`. The snapshot is unchanged. The lock's `content_hash` was a placeholder (`sha256:000…`), and `ai_check` accepted it. |
| V17 | The M2 snapshot recipe, run verbatim against Daneel's lock, produced a tree identical to V14's (`diff -r` empty), with 0 writable entries and 8 package dirs. |
| V18 | `--exec-episode-project` is refused without `--exec-profile` (`host/broker/exec_startup.go:96–98`). A row-140 profile already has `root` ("the project root inside the worktree"). It is validated by `execPathOk` plus `Clean` (`exec_profile.go:154–160`, `execargs.go:148–155`) and resolved symlink-free by `projectDir` (`handlers_exec.go:350–361`). |
| V19 | `fakeToolBin` tests run without any tool binary: `TestWorkspacePolicyRenderedPerEpisode` passes. The real-binary tests SKIP when `WORLD_TOOL_AILANG_BIN` is unset (`TestWorkspaceSandboxIsTheEpisodeNotTheRoot`). The fake binary reports `fs_sandbox = $(pwd -P)` and logs its env (`workspace_test.go:91–121`). |
| V20 | No upstream issue covers either ask (5 searches, none relevant). Positive control: `examples policy-tool` finds #1552. |
| V21 | `--version` reports `AILANG v0.52.1`, commit `c68ded4`, which equals `v0.52.1^{commit}`. The binary's sha256 begins `0dd70a1d00360be0`. |
| V22 | The operator's own `~/.ailang/cache/registry` is writable (6173 entries with `u+w`) and is 233 MB. The V14 snapshot of Daneel's 8 packages is 480 KB. |
| V23 | Where `ailang run --policy` writes its compile cache, under World's exact child env with the `tools/` sandbox and a fresh `HOME`. It writes **inside the sandbox**, at `tools/.ailang/cache/compile/{manifest.json,modules/…}`, and ignores `AILANG_CACHE_DIR` (#1547 confirmed). The rendered `.ailang/**` `fs_deny_write` entry does **not** stop the runtime's own cache write. It **does** refuse an agent `write` to `.ailang/cache/compile/planted.txt` (`matches fs_deny_write ".ailang/**"`; control: `write probe_ok.txt` → `ok: true`). In `HOME` the run creates only the empty `.ailang/cache/registry`. Control: `ai_check` of the same module through policy-tool adds nothing to the sandbox or to `HOME`. |
| V24 | **Every writer under `<state>/cache/E/.ailang/**`.** World creates only `<state>/cache/E` itself (`workspace.go:380–384`). Its 4 non-test `".ailang"` sites are the deny list (`handlers_ailang.go:57`), the `setup` defaults (`setup.go:82–83`) and the examples default (`main.go:423`), none of which write under the cache dir. v0.52.1 writes there in two places. `RegistryCacheDir` creates the **empty** `.ailang/cache/registry` (V10, V11, V23). `pkg_docs` extracts packages into it (V13). The compile cache is not one of them: `run` writes it into the sandbox (V23), and CLI children write it to a private temp dir (V4). Neither is the prompt cache: it lives at `$AILANG_CACHE_DIR/prompts/…` (`internal/prompt/fresh.go:195–200`), outside `.ailang`. So a non-empty `.ailang/cache/registry` under `<state>/cache/E` can only have come from `pkg_docs`, or from a person. |

## Options

### (a) Where the module root is declared

**(a1) On row 140's exec profile `root` / `--exec-episode-project` map. Rejected.** Row 140 already validates a per-project `root` (V18), so sharing it is tempting. But `--exec-episode-project` refuses to start without `--exec-profile` (V18). An operator who only wants `ai_check` to work would have to configure the whole `workspace-exec` stack (srt, node, the startup probe) to get a module root. That couples two unrelated features. The two roots are also different things: exec runs a toolchain from its project root, while the AILANG tools need the sandbox to *be* the module root (V4, V5). The design does reuse row 140's path grammar and its symlink-free resolution (V18), so the two cannot drift apart.

**(a2) New serve flags, independent of exec. RECOMMENDED.**
- `--workspace-module-root REL` sets the module root for every episode.
- `--workspace-episode-module-root EP=REL` (repeatable) overrides it for one episode.

When neither is set, behaviour is exactly today's.

**(a3) Upstream: a `module_root` policy key, so the sandbox stays at the worktree root. Not pursued.** v0.52.1 has no such key (V6). The row asks for "a subdirectory sandbox", which v0.52.1 can do today (V7). A subdirectory sandbox is also strictly *narrower*: an agent working in `tools/` cannot write `ext/`, `.github/` or the worktree's `.git`. If a later dogfood needs edits outside the module root, that becomes its own row and its own upstream ask.

### (b) The registry package cache

**(b1) Read-only operator snapshot, linked into the per-episode `HOME`. RECOMMENDED; works on v0.52.1.**
- New flag `--workspace-package-cache DIR`. The operator provisions DIR as a read-only snapshot of exactly the locked registry packages, using the M2 recipe (V17).
- For each episode, World makes `<state>/cache/E/.ailang/cache/registry` a symlink to DIR's canonical path.
- The child env does not change. `HOME` is still `<state>/cache/E` (V2), and v0.52.1 resolves registry imports through it (V11).

Why this is safe:
- `ai_check`, `test` and `run` resolve against the snapshot (V14, V15).
- `pkg_docs` can no longer store a fetched package (V14).
- A confined run cannot write the snapshot (V15), and `Workspace.Write` is confined to the sandbox (row 134 §4.3).

**Leftover content:** World never deletes anything that holds content under the episode `HOME`. An unprovisioned `registry` directory refuses the episode with a named operator line (M2, V24).

What it depends on: the `.ailang/cache/registry` layout under `HOME` is upstream *source*, not a documented interface (V11). Two things cover that. The rig e2e test (AC2.6) is the tripwire. Upstream ask U2 asks for a first-class, read-only package root, and World would drop the symlink once that ships.

**(b2) Point `HOME` straight at a shared, operator-owned home. Rejected.** `HOME` is also where the unconfined CLI children write: the registry root `MkdirAll` and `pkg_docs` extraction (V10, V13). A home shared across episodes would make those writes visible from one episode to the next.

**(b3) The lock's `path` fallback for `registry` entries (V11). Rejected.** `ailang.lock` sits in the worktree, and the agent can write it (it is not in the deny list). A cache location chosen by the lock is a location chosen by the agent.

**(b4) Copy the snapshot into each episode's `HOME`. Rejected.** This is the pattern of row 140's exec caches. Those are seeded from an operator's read-only `seed`, but into `<state>/exec-cache/<id>/` (`host/broker/exec_probe.go:315–321`, with dirs `home`, `tmp`, `pycache`, `npm`), not into the AILANG tools' `<state>/cache/E`. The copy would be World-owned and writable, so `pkg_docs` could add to it again (V13), unless World also re-applied `a-w` to every copy. The symlink gets read-only for free from a tree that is checked read-only at startup.

**Why host code, not a package (S3):** this is serve-time configuration of the sandbox root and the child's `HOME` at the host boundary (S2), the same layer as `--workspace-root` and `--examples-dir`. No `world/` kernel or `.ail` change.

## Design

### M1 — Module root: a subdirectory sandbox (~0.35 d)

**Flags and config.**
- `Config` gains two fields:
  - `WorkspaceModuleRoot string` (`--workspace-module-root`).
  - `WorkspaceEpisodeModuleRoots []string` (`--workspace-episode-module-root EP=REL`, repeatable).
- Both need `--workspace-root`. Without it, startup is refused with a `StartupError`, the same refusal shape as the `--run-*` flags (`workspace.go:428–431`).

**Startup validation (pure, no filesystem).** Each `REL`:
- passes `execPathOk` (no leading `/` or `-`, no `..` segment; V18);
- equals `filepath.Clean(REL)`;
- holds no character `renderPolicy` would refuse. Check this by rendering `RenderEpisodePolicy(filepath.Join(w.root, REL))` once at startup.

`"."` is allowed and means "the worktree root". Each `EP` must match `episodeIDPattern`. An episode mapped twice is refused, with the same refusal shape as `exec_startup.go:162–163`.

**Per-episode resolution.** A new `w.sandboxRoot(episodeID, epRoot) (string, error)` picks `REL` in this order: the episode's override, else the default, else `"."`. It then requires:
- `want := filepath.Join(epRoot, REL)`;
- `filepath.EvalSymlinks(want) == want`;
- `want` is a directory.

This is the `projectDir` rule (V18), applied after `episodeRoot`. A module root that is missing, a symlink, or not a directory makes `registry()` return an **empty** registry, so the coordinator refuses the declared effect (R8) before any effect runs. `registry()` also writes one operator line, `ailang-worldd: workspace tools unavailable for episode "E": module root "REL": <why>`, through the existing error path (`workspace.go:304–307`). Nothing is ever created: World never `MkdirAll`s a module root.

**Wiring.**
- `episodeHandler` renders `RenderEpisodePolicy(sandbox)` and passes `Root: sandbox`.
- It caches `episodeTool{root: sandbox}`.
- `loadSummary`'s existing `fs_sandbox == Root` check (`handlers_ailang.go:341–343`) then holds by construction.
- The policy path (`policies/E.toml`) and `cacheDir` are unchanged.
- `execHandler` keeps `epRoot`. An exec profile's own `root` is independent (a2 above), and the docs say so.

**What the agent sees.** Paths in `read`/`write`/`ai_check`/`run` payloads are relative to the module root (`daneel_brief.ail`, not `tools/daneel_brief.ail`). Every result's `policy_digest` (V2) changes with `fs_sandbox`, so the effect record already says which root ran. No record-format change is needed.

**Confinement (row 134 §4.3 is unchanged and AILANG's own).**
- The sandbox shrinks from `W/E` to `W/E/REL`. Writes outside it, including the worktree's `.git`, `.github/` and `.claude/` at `W/E`, are now outside `fs_sandbox` and refused by AILANG's path confinement.
- The deny globs are sandbox-relative, so they also cover any `.github/`, `.claude/`, `.pi/` or `.ailang/` inside REL.
- Measured on the `tools/` sandbox (V23): `run --policy` writes its compile cache to `REL/.ailang/cache/compile` (#1547), and the `.ailang/**` deny entry does **not** stop that runtime write. What the entry does stop is an agent `write` into `.ailang/**`. This is the same behaviour as at the worktree root today. The subdirectory sandbox neither causes it nor fixes it, and nothing here relies on it as a safety property.
- Reads by the unconfined module loader already reach outside the sandbox (path dependencies, V10). This row does not change that (see Risks).

### M2 — Read-only registry package cache (~0.4 d)

**Flag and config.** `Config.WorkspacePackageCache string` (`--workspace-package-cache DIR`). It needs `--workspace-root`.

**Startup checks** live in a new `resolvePackageCache(dir, canonicalRoot, stateDir) (canonical string, digest string, n int, err error)`. It refuses a `StartupError` naming the flag when any of these holds:
1. DIR is not absolute after `Abs`, does not resolve (`EvalSymlinks`), or is not a directory.
2. DIR lies inside `--workspace-root` (`broker.CheckPolicyOutsideRoot`) or inside the state dir (lexical `within`, as in `resolveWorkspaceRoot`).
3. **Any** entry under DIR, including DIR itself, is writable by anyone (`Lstat().Mode().Perm() & 0o222 != 0`), whether directory or file.
4. Any entry is a symlink. A link inside a read-only tree can point at a writable one.
5. Fewer than one package exists. A package is a `<ns>/<name>/<ver>/ailang.toml` path. A vacuous cache fails loudly (S6).
6. The daemon runs as uid 0, where mode bits do not stop writes. A `geteuid` seam var lets tests reach this case.

**Snapshot digest.** The check also computes a tree digest over the walk, in byte order. A directory contributes `relpath NUL "d"`. A file contributes `relpath NUL "f" NUL sha256(content)`. Entries are separated by newlines, and the whole is SHA-256 hashed, giving `sha256:<hex>`. The digest is computed **once, at startup**, and never per episode (gemini r1). Startup logs one line: `ailang-worldd: workspace package cache DIR: N packages, read-only, sha256:…`.

**Per episode.** In `episodeHandler`, after the existing `MkdirAll(cacheDir)`, `linkPackageCache(cacheDir)` runs.
- **Flag set:**
  - `MkdirAll(cacheDir/.ailang/cache, 0700)`.
  - Let `link := cacheDir/.ailang/cache/registry`. **Nothing that holds content is ever deleted** (kimi r1).
    - A symlink whose `Readlink` equals the canonical DIR is kept.
    - If `link` is absent, or is an **empty** directory, it is removed (`os.Remove` on the empty directory only) and replaced by `Symlink(DIR, link)`. An empty directory is what v0.52.1's own `RegistryCacheDir` leaves behind (V24).
    - **Anything else** is a handler-construction error: a non-empty directory, a file, or a symlink to anywhere else. Nothing is deleted. The error takes the existing empty-registry path (R8), and the operator line is `ailang-worldd: workspace package cache refused for episode "E": <link> is not empty and was not provisioned by --workspace-package-cache; clear it manually`. Per V24, such content can only have come from `pkg_docs` or from a person, and deciding what to do with it is the operator's call.
  - There is **no** per-episode re-walk of DIR. Episodes are stamped with the startup digest (gemini r1: a per-episode walk is unbounded I/O on the dispatch path).
- **Flag unset:** if `link` is a **symlink**, it is removed. Only an earlier run with the flag set creates one there (V24). A real directory, empty or not, is never touched, which is today's behaviour. Nothing else changes.

**Stamp.** `AilangToolConfig` gains `PackageCacheDigest string`. When it is non-empty, `stamp` adds `"package_cache": "<digest>"` next to `policy_digest`. When it is empty, the key is absent, so results with the flag unset are byte-identical to today's. This records *which* packages an `ai_check` saw. Without it, the result would depend on an input the record does not name.

**Lock file: an operator step, not code.** v0.52.1 needs `ailang.lock` in the module root (V9, V16). The agent cannot create it (`lock` is not admitted, V12). Daneel ignores it in git (V9). The docs give the step, run after `session new` and before the episode starts: run the project's own `ailang lock` in `W/E/REL`. The lock is then the input to the snapshot recipe.

**Snapshot recipe (S7; executed verbatim in M3, measured V17).**

```bash
LOCK=<workspace-root>/<episode>/<module-root>/ailang.lock
SNAP=<a dir outside the workspace root and the db dir>
jq -r '.packages[] | select(.source == "registry") | "\(.name)/\(.version)"' "$LOCK" | while read -r p; do mkdir -p "$SNAP/$(dirname "$p")" && cp -R "$HOME/.ailang/cache/registry/$p" "$SNAP/$p"; done
chmod -R a-w "$SNAP"
```

To refresh the snapshot, run `chmod -R u+w "$SNAP"`, rebuild it, and restart `serve`: the digest is fixed at startup.

**More than one episode.** `--workspace-package-cache` is daemon-global, but locks belong to each worktree. A deployment with several episodes must therefore build the one `SNAP` as the **union** of every episode's locked registry packages: run the recipe's `jq … | while …` loop once per episode's `LOCK` into the same `SNAP`, then run `chmod -R a-w` once at the end.

### M3 — Docs, upstream asks, follow-up rows (~0.25 d)

- **S7 docs.**
  - `serve --help` documents all three flags.
  - A new `docs/QUICKSTART.md` §11, "Projects with a module root and registry packages", shows the Daneel shape: `session new`, the operator's `ailang lock`, the recipe, and one `serve` line. A binding test, `TestQuickstartSection11FlagsMatchTheCLI`, is modelled on `quickstart_exec_test.go:95–130`.
  - Mirror lines go in `website/docs/guides/operating-the-daemon.md` and `website/docs/reference/cli.md`.
  - The docs state the independence of the exec profile's `root` (M1).
- **Upstream asks.** File both on `sunholo-data/ailang`, then send `ailang messages send mission-control` with the links (body positional). The text is below.
- **Charter follow-up rows**, proposed for the controller to file:
  - **R-A:** "`pkg_docs` is an unrecorded network egress and an agent-steered cache write". Options: render `cli_allow` without `pkg_docs` (V6), or refuse the op in the handler as `cliWriteFlags` does, or wait for U1. Rendering `cli_allow` changes the policy bytes and the audited 20-op table, which is why it is not in this row.
  - **R-B:** "`ailang.lock` is agent-writable, and the unconfined loader follows its `path` entries outside the sandbox". This is an observation from V10, not an exploit test. Measure it first.

**U1 — issue text (area: policy-tool):**
> **`pkg-docs` under an operator policy fetches from the registry and writes `$HOME/.ailang/cache/registry`.** On v0.52.1, `policy-tool` admits `pkg_docs` in its default read-only `cli` set. When the package is not cached, `cmd/ailang/pkg_docs.go:71–90` calls `FetchIndex`/`FetchPackage` and `ExtractTarball`s into the registry cache. The child is unconfined (`cli_ops.go:156`), so a confined agent gets network egress and a persistent write. A later `ai-check` resolves imports against that agent-fetched copy. Repro: run `policy-tool` with a restricted policy and an empty `HOME`, send `{"op":"pkg_docs","module":"sunholo/oauth"}`, and `<HOME>/.ailang/cache/registry/sunholo/oauth/0.1.0/` appears. Ask: when `AILANG_AGENT_POLICY` is set, `pkg-docs` should read the cache only and never fetch. That is the same principle as #1552 for the examples corpus.

**U2 — issue text (area: pkg):**
> **No way to point package resolution at a read-only, operator-provisioned package root.** Registry packages resolve only under `os.UserHomeDir()/.ailang/cache/registry` (`internal/pkg/registry.go:208–217`), which `RegistryCacheDir` also `MkdirAll`s. `AILANG_CACHE_DIR` covers the compile and prompt caches only (`internal/config/paths.go:21`), and the policy has no package key. A host that confines tools (AILANG World) has to fake a `HOME` containing a symlink to make a read-only snapshot visible. Ask: `AILANG_PACKAGE_ROOT` (or a policy key) that the `PackageLoader` reads read-only and never writes, plus an optional content-hash check of each registry package against the lock's `content_hash`. v0.52.1 accepts a placeholder hash at check time.

## Acceptance criteria

Every criterion names one test and the mutant it kills (S6). **CI** means the test runs in `go test ./...` with no tool binary: it uses `fakeToolBin` (V19) or is a pure unit test. **RIG** means it needs `WORLD_TOOL_AILANG_BIN`, so it SKIPs in CI until row 160 lands. Every RIG test has a CI companion that kills the same mutant. All fixtures are created by the test itself, using `t.TempDir`/`canonicalTemp` and `os.WriteFile`/`os.Symlink`/`os.Chmod`. The RIG fixture follows V16's generator: a two-module project, a lock marshalled from a Go struct, and a snapshot made read-only by the test.

| AC | Test | Kind | Asserts | Mutant it kills |
|---|---|---|---|---|
| AC1.1 | `host/daemon` `TestWorkspaceModuleRootGrammar` (table) | CI | `tools`, `tools/sub` and `.` are accepted. `""` as a value, `/abs`, `../x`, `tools/../../x`, `-x`, `tools/`, `a"b` and `EP` outside the grammar are refused at startup, as is an episode mapped twice. | **MUT-MR-DOTDOT**: drop the `..`-segment check, and the `tools/../../x` row fires. **MUT-MR-DUP**: last mapping wins, and the duplicate row fires. |
| AC1.2 | `TestWorkspaceModuleRootIsTheSandbox` | CI | With `--workspace-module-root tools` and `ep1/tools` present: `policies/ep1.toml` equals `RenderEpisodePolicy(<root>/ep1/tools)` byte for byte. The fake binary's summary and dispatch `pwd` both equal `<root>/ep1/tools`. The registry binds exactly the 9 names of `workspaceEffects` (`workspace.go:55–66`): `Workspace.Read`, `Workspace.Write`, `Ailang.Check`, `Ailang.Run`, `Ailang.RunEnv`, `Ailang.RunNet`, `Ailang.Discover`, `Ailang.CLI`, `Workspace.Exec`. | **MUT-MR-POLICY-ROOT**: the policy is rendered with `epRoot`, and the byte-equality assertion fires. **MUT-MR-CWD**: `Root: epRoot` is passed, and the `pwd` assertion fires. |
| AC1.3 | `TestWorkspaceModuleRootSymlinkOrMissingIsR8` | CI | `ep1/tools` as a symlink (to `ep1/real`, and to `../ep2/tools`), as a file, or missing → empty registry, zero summaries and dispatches, one operator line naming `module root "tools"`. The missing directory is **not** created. | **MUT-MR-NOSYMLINK**: the `EvalSymlinks == want` check is skipped, and the symlink row fires. **MUT-MR-MKDIR**: the root is `MkdirAll`ed, and the not-created assertion fires. |
| AC1.4 | `TestWorkspaceEpisodeModuleRootOverrides` | CI | With default `tools` and `ep2=.`, ep1's sandbox is `ep1/tools` and ep2's sandbox is `ep2` (logged `pwd`). | **MUT-MR-OVERRIDE-IGNORED**: the override is ignored, and the ep2 assertion fires. |
| AC1.5 | `TestWorkspaceModuleRootNeedsWorkspaceRoot` | CI | Either flag without `--workspace-root` → `StartupError`, and the store is not opened. | **MUT-MR-SILENT**: the flags are ignored when the root is unset, and the error assertion fires. |
| AC1.6 | existing `TestWorkspacePolicyRenderedPerEpisode`, **byte-unchanged** | CI | With the flags unset, the policy equals `RenderEpisodePolicy(epRoot)` and the env holds `AILANG_CACHE_DIR=<state>/cache/ep1`. | **MUT-MR-DEFAULT**: an unset root defaults to anything but `.`, and the existing byte-equality check fires. |
| AC1.7 | `host/daemon/setools_e2e_test.go` `TestSeToolsModuleRootResolvesBareImports` | RIG | V16 fixture. Without the flag, `ai_check tools/a.ail` → `LDR001 module not found: b`. With `--workspace-module-root tools`, `ai_check a.ail` → `passed`. | MUT-MR-POLICY-ROOT and MUT-MR-CWD (both also killed in CI by AC1.2) |
| AC2.1 | `TestWorkspacePackageCacheStartupChecks` (table) | CI | A read-only tree with one package is accepted, and its digest is stable across two calls. Refused, one row each: not a dir; inside `--workspace-root`; inside the state dir; a writable **file**; a writable **dir**; DIR itself writable; a symlink inside; zero packages; uid 0 (seam); no `--workspace-root`. | **MUT-PC-DIRS-ONLY**: only directories' modes are checked, and the writable-file row fires. **MUT-PC-SYMLINK**: symlinks are not refused, and the symlink row fires. **MUT-PC-EMPTY**: the package count is dropped, and the zero-packages row fires. |
| AC2.2 | `TestWorkspacePackageCacheLinkedIntoEpisodeHome` | CI | After `registry("ep1")`, `Lstat(<state>/cache/ep1/.ailang/cache/registry)` is a symlink to the canonical DIR. The fake binary's logged env still has `HOME=<state>/cache/ep1`, and the policy bytes equal the flag-unset render. | **MUT-PC-NOLINK**: the link step is skipped, and the `Lstat` assertion fires. |
| AC2.3 | `TestWorkspacePackageCacheRefusesUnprovisionedRegistry` | CI | Flag set. (i) A non-empty directory pre-planted at `link`, holding `sunholo/x/9.9.9/ailang.toml`, a file at `link`, and a symlink to another dir each give an **empty registry**, zero summaries, and exactly one operator line `ailang-worldd: workspace package cache refused for episode "ep1": <link> is not empty and was not provisioned by --workspace-package-cache; clear it manually`. The planted file **survives**, byte-identical. (ii) An **empty** directory at `link` is replaced by the symlink. | **MUT-PC-SILENT-DELETE**: `RemoveAll` and relink any entry, and the planted-file-survives assertion fires. **MUT-PC-KEEP-EMPTY**: an empty dir is refused instead of replaced, and arm (ii) fires. |
| AC2.4 | `TestWorkspacePackageCacheUnsetLeavesHomeAlone` | CI | Flag unset: a pre-existing non-empty `registry` directory **and** an empty one are both untouched, and a registry still forms. A pre-existing symlink is removed, and its target is untouched. | **MUT-PC-UNSET-REMOVES-DIR**: anything but a symlink is removed, and the untouched-dir assertions fire. **MUT-PC-UNSET-FOLLOWS-LINK**: `RemoveAll` on the symlink's target, and the target-untouched assertion fires. |
| AC2.5 | `TestWorkspacePackageCacheDigestIsStamped` (fake-binary dispatch through the handler) | CI | Flag set: the result carries `"package_cache":"sha256:…"` equal to the startup digest. Flag unset: the result has **no** `package_cache` key. The digest equals the value recomputed by the test from its own fixture under the r2 rule (directory: `relpath NUL d`; file: `relpath NUL f NUL sha256(content)`). | **MUT-PC-STAMP-ALWAYS**: the key is stamped when empty, and the unset assertion fires. **MUT-PC-DIR-HASH**: directories are hashed like files, and the recomputed-digest assertion fires. |
| AC2.6 | `TestSeToolsPackageCacheResolvesRegistryImports` | RIG | V16 fixture with `--workspace-module-root tools`. Without `--workspace-package-cache`, `ai_check a.ail` → `cache not found`. With it → `passed`. The snapshot's tree digest is unchanged after `ai_check`, `test` and `ailang-run`. This is the tripwire for the upstream `HOME` layout (V11). | **MUT-PC-NOLINK** (also killed in CI by AC2.2) |
| AC3.1 | `cmd/ailang-worldd` `TestQuickstartSection11FlagsMatchTheCLI` | CI | `serve --help` documents the three flags. §11's single `serve` line uses them, and every flag on it is documented and parses (it stops at an appended bad `--bind`). | **MUT-HELP-DROP**: one flag is removed from the usage text, and the documented-flag assertion fires. |
| AC3.2 | the U1/U2 issue URLs and the `ailang messages` id, recorded in the iteration log, plus the R-A/R-B rows | — | the asks exist and can be addressed | n/a: evidence, not a test (a labelled instrument) |

Gate: `./scripts/verify_ail.sh` (with `AILANG_BIN`; no `.ail` changes) and `go vet ./... && go test -race ./...` with `AILANG_BIN=$HOME/.pinned-ailang/ailang`. Run once with `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.52.1/ailang` so AC1.7 and AC2.6 run, and check each prints `--- PASS`, not `--- SKIP`. Use `-run '^(TestSeToolsModuleRootResolvesBareImports|TestSeToolsPackageCacheResolvesRegistryImports)$'`.

## Conflict Surface

| File | Change | Unset behaviour |
|---|---|---|
| `host/daemon/daemon.go` (`Config`, startup at `:561–574`) | 3 fields. Module-root grammar and package-cache checks run next to `resolveExamplesDir`, before the store opens. | No new check runs. |
| `host/daemon/workspace.go` (`workspaceTools`, `episodeHandler`, `registry`) | `sandboxRoot`, `linkPackageCache`, and two new struct fields. The handler cache key root becomes the sandbox. | Sandbox = `epRoot`. Policy and env bytes are identical (AC1.6). The only side effect is removing a stale symlink at `<state>/cache/E/.ailang/cache/registry` (AC2.4). |
| `host/broker/handlers_ailang.go` (`AilangToolConfig`, `stamp`) | Adds `PackageCacheDigest`, and `package_cache` in results only when it is set. | Results are byte-identical (AC2.5). |
| `cmd/ailang-worldd/main.go` | 3 flags, usage, and `serve --help`. | No change unless the flags are given. |
| `docs/QUICKSTART.md` §11; `website/docs/guides/operating-the-daemon.md`; `website/docs/reference/cli.md` | New text. | Existing sections are unchanged. §9 and §10 binding tests are unaffected. |
| tests | New tests in `host/daemon/workspace_test.go`, `setools_e2e_test.go` and `cmd/ailang-worldd/quickstart_*_test.go` | Existing tests are byte-unchanged. |

Not touched: `tools/launchd/*`, `world/`, `packages/`, `scripts/verify_ail.sh`, the row-140 exec path (`execHandler` keeps `epRoot`), `configureRunCaps` (still checked against `w.root`), `session new`, and `doctor`.

## Non-goals

- Editing files outside the module root while the AILANG tools are rooted there (a3; a future row plus an upstream ask if a dogfood needs it).
- Automating `ailang lock` or the snapshot. Both are operator steps, documented in §11.
- Git-source packages (`$HOME/.ailang/cache/git`). Daneel has none (V9 lock listing).
- Fixing `pkg_docs` egress (R-A, U1) or `ailang.lock` writability (R-B). Note, though, that M2 already stops `pkg_docs` from *storing* a package whenever the flag is set (V14).
- `doctor` reporting the new flags.

## Risks

- **R1: the upstream `HOME` layout moves.** The symlink depends on `internal/pkg/registry.go`, not a documented interface (V11). Mitigations: AC2.6 reds on a pin bump that moves it, and `ToolBinaryRelease` already refuses any other release at startup. U2 removes the dependency.
- **R2: read-only depends on mode bits.** Root bypasses them, which is why check 6 refuses uid 0. The operator could also `chmod u+w` the tree and change it while serve runs. The digest is computed once at startup (gemini r1), so such a change goes undetected for **every** episode, including ones built after it, and their results carry a stamp that no longer describes the tree. This risk is accepted and documented: the snapshot is operator state outside the workspace root, and the documented refresh is a rebuild followed by a restart (M2).
- **R3: the unconfined loader reads outside the sandbox.** This covers path dependencies (V10), the stdlib, and the snapshot. It is pre-existing AILANG behaviour, recorded as R-B. A subdirectory sandbox does not widen it.
- **R4: packages are not verified against the lock.** `ai_check` accepted a placeholder `content_hash` (V16), so the snapshot is trusted as provisioned. The digest stamp (M2) makes what was used auditable, and U2 asks upstream to verify.

## Verification Log

All probes used `B=~/.pinned-ailang-tools/v0.52.1/ailang` and `S=~/.ailang/state/world-iter238-row141`. The Daneel copy is `git -C ~/dev/daneel archive HEAD | tar -x -C $S/daneel` (HEAD `ebe26c0`). Policies were rendered by World's own `broker.RenderEpisodePolicy`, through a scratch `main.go` built with a `go.work` that uses this worktree (`$S/render <root> > policy.toml`). Every tool call used World's exact child env: `env -i HOME=$H PATH=/usr/bin:/bin LANG=C LC_ALL=C AILANG_CACHE_DIR=$H $B policy-tool --policy P`, with cwd = the sandbox. Run 2026-10-06.

| V | Command (trimmed) | Observed (trimmed) |
|---|---|---|
| V1 | `grep -n 'cacheDir := \|RenderEpisodePolicy(epRoot)\|Root: epRoot' host/daemon/workspace.go` | `380: cacheDir := filepath.Join(w.stateDir, "cache", episodeID)`, `387: … RenderEpisodePolicy(epRoot)`, `397: … Root: epRoot, CacheDir: cacheDir` |
| V2 | `sed -n 496,500p;522,533p host/broker/handlers_ailang.go` | `"HOME=" + h.cacheDir`, `"AILANG_CACHE_DIR=" + h.cacheDir`; `resp["tool"]`, `resp["policy_digest"]` |
| V3 | r2 re-run; exact command in the block below the table | 1 hit: `host/pinfetch/pinfetch.go:183:// The release subdirectory keeps the two pins apart` (unrelated); control → `handlers_ailang.go:229`, `handlers_ailang.go:528`, `workspace.go:8` |
| V4 | `git show v0.52.1:internal/policytool/cli_ops.go \| grep -nF 'dir = h.root.Dir()'` (and the other two) | `316`; `156` (`separate, UNCONFINED process`); `405` (`config.EnvCacheDir+"="+cacheDir`) |
| V5 | `… pipeline_module_phases.go \| grep -nF 'loaderBaseDir := "."'` | `153` |
| V6 | `git grep -n 'toml:"' v0.52.1 -- internal/policy/ \| sed … \| sort -u` | `ai_provider allowed_caps budgets cli_allow entry fs_deny_write fs_sandbox max_fs_transfer_bytes max_module_graph_bytes max_output_bytes max_source_bytes net_allow net_allow_http process_allow security_mode timeout_ms` (16) |
| V7 | `(cd $S/daneel && … <<<'{"op":"ai_check","path":"tools/daneel_brief.ail"}')`; `(cd $S/daneel/tools && … policy-tools.toml <<<'{"op":"ai_check","path":"daneel_brief.ail"}')` | root: `"ok": false … "code": "LDR001", "message": "module not found: daneel_reading"`; tools: `"ok": true … "passed": true … "verified": 6` |
| V8 | `cd $S/daneel/tools && … --policy policy-daneel.toml <<<'{"op":"ai_check","path":"tools/daneel_brief.ail"}'` | `False ['LDR001: module not found: daneel_reading']` |
| V9 | `git -C ~/dev/daneel check-ignore -v tools/ailang.lock`; `git ls-files tools/ailang.lock \| wc -l`; `$B ai-check daneel_calendar.ail` in the copy (real HOME) | `.gitignore:16:ailang.lock tools/ailang.lock`; `0`; `package "sunholo/oauth" not found in ailang.lock; run 'ailang lock'` |
| V10 | lock copied in; `ai_check` daneel_ext / daneel_calendar with `H=$S/cacheE` (empty); `find $S/cacheE` | `daneel_ext: ok=True passed=True`; `daneel_calendar: … registry package sunholo/oauth cache not found at $S/cacheE/.ailang/cache/registry/sunholo/oauth/0.1.0`; `$S/cacheE/.ailang/cache/registry` now exists |
| V11 | `git show v0.52.1:internal/pkg/registry.go \| sed -n 207,217p`; `loader.go` `grep -nF 'case "registry":'`; `paths.go:21` | `dir := fmt.Sprintf("%s/.ailang/cache/registry", home)` (213), `os.MkdirAll(dir, 0755)` (214); `171`; `Root of the compile cache (<dir>/compile) and the prompt cache` |
| V12 | `… <<<'{"op":"summary"}' \| python3 -I -c '…print(s["cli"], len(…))'` | 20 ops: `agent_prompt ai_check axioms builtins_list builtins_show check devtools_prompt docs_search examples_list examples_search examples_show examples_tags fmt iface pkg_docs policy_check prompt test tree version` |
| V13 | `H=$S/cacheP` (empty) `<<<'{"op":"pkg_docs","module":"sunholo/oauth"}'`; `find $S/cacheP`; then `ai_check daneel_sources.ail` with the same H | `ok= True argv= ['pkg-docs','sunholo/oauth'] stdout='# sunholo/oauth …'`; `…/registry/sunholo/oauth/0.1.0`; `ok= True passed= True verified= 4` |
| V14 | build `$S/pkgro` (8 packages, `chmod -R a-w`; `find -perm -u+w \| wc -l` = 0), `ln -s $S/pkgro $S/homeR/.ailang/cache/registry`; `ai_check` over `grep -l '^import pkg/' *.ail`; `pkg_docs sunholo/auth`; tree digest before/after | all 12 `ok= True passed= True` (verified 0,4,12,0,0,3,6,0,10,6,1,2); `ok= False 'Error: package sunholo/auth@0.4.1 has no AGENT.md'`; `9bce9a56…` = `9bce9a56…` |
| V15 | `$B run --policy policy-tools.toml -- probe_pkg.ail` with `H=cacheE`, then `H=homeR`; `probe_write.ail` writes `$S/pkgro/sunholo/planted.txt` | cacheE: `"error_kind": "typecheck_failed" … cache not found`; homeR: `pkg-ok`; `writeFile: path "…/pkgro/sunholo/planted.txt" escapes sandbox "…/daneel/tools"`; `ls planted.txt`: No such file |
| V16 | synthetic `$S/syn` (generated by shell `printf` + `python3 json.dump`): S1 root policy, S2 tools + empty H, S3 tools + snapshot H | S1 `LDR001: module not found: b`; S2 `registry package acme/util cache not found at …/emptyhome/…`; S3 `ok= True passed= True []`; tree `26c138c7effe27b3 -> 26c138c7effe27b3`; lock `content_hash: "sha256:000…0"` accepted |
| V17 | the M2 recipe verbatim with `LOCK=$S/daneel/tools/ailang.lock SNAP=$S/snap-recipe`; `diff -r $S/pkgro $SNAP`; `find $SNAP -perm -u+w \| wc -l`; `find $SNAP -mindepth 3 -maxdepth 3 -type d \| wc -l` | `IDENTICAL`; `0`; `8` |
| V18 | `sed -n 96,98p host/broker/exec_startup.go`; `grep -n 'Root' host/broker/exec_profile.go`; `sed -n 350,361p host/broker/handlers_exec.go` | `needs --exec-profile (it configures workspace-exec, which is off without a profile)`; `58: Root string // the project root inside the worktree`, `159: !execPathOk(p.Root) \|\| filepath.Clean(p.Root) != p.Root`; `got, err := filepath.EvalSymlinks(want); if err != nil \|\| got != want` |
| V19 | `env -u AILANG_BIN -u WORLD_TOOL_AILANG_BIN go test ./host/daemon/ -run '^(TestWorkspacePolicyRenderedPerEpisode\|TestWorkspaceSandboxIsTheEpisodeNotTheRoot\|TestWorkspaceRegistryRefusesEpisodesOutsideTheGrammarAndRoot)$' -v` | `--- PASS: …RefusesEpisodesOutside…`, `--- PASS: TestWorkspacePolicyRenderedPerEpisode`, `--- SKIP: TestWorkspaceSandboxIsTheEpisodeNotTheRoot` (`WORLD_TOOL_AILANG_BIN is unset`) |
| V20 | `gh issue list -R sunholo-data/ailang --state all --search` `pkg-docs` / `registry cache read-only` / `AILANG_CACHE_DIR registry` / `module root policy-tool` / `package cache policy` | no relevant hit (closest: #1547 `run --policy ignores AILANG_CACHE_DIR`, a different issue); control `examples policy-tool` → `#1552` |
| V21 | `$B --version`; `git rev-parse 'v0.52.1^{commit}'`; `shasum -a 256 $B \| cut -c1-16` | `AILANG v0.52.1 Commit: c68ded4`; `c68ded4b2d5e…`; `0dd70a1d00360be0` |
| V22 | `find ~/.ailang/cache/registry -perm -u+w \| wc -l`; `du -sh ~/.ailang/cache/registry $S/pkgro` | `6173`; `233M`, `480K` |
| V23 | `T=$S/daneel/tools`; `rm -rf $T/.ailang`; fresh `H=$S/h2`; `find` listings of `$T` and `$H` before and after `(cd $T && env -i HOME=$H … AILANG_CACHE_DIR=$H $B run --policy policy-tools.toml -- probe_hello.ail)`; then the same for `policy-tool` `ai_check probe_hello.ail`; then `policy-tool` `{"op":"write","path":".ailang/cache/compile/planted.txt",…}` and the control `{"op":"write","path":"probe_ok.txt",…}` | run: `hello rc=0`; new in sandbox: `./.ailang/cache/compile/manifest.json`, `./.ailang/cache/compile/modules/probe_hello/{artifacts.json,constructors.json,core.gob,coretypeinfo.gob}`; new in HOME: `./.ailang/cache/registry` only. ai_check: `"ok": true`, nothing new in either. write: `"refused": "write .ailang/cache/compile/planted.txt: matches fs_deny_write \".ailang/**\" — read-only under this policy"`; `ls planted.txt`: No such file; control `{"ok": true}` |
| V24 | `grep -rn '"\.ailang"' host cmd --include='*.go' \| grep -v _test.go`; `sed -n 380,384p host/daemon/workspace.go`; `git grep -n 'config.CacheDir()' v0.52.1 -- '*.go' \| grep -v _test`; `git show v0.52.1:internal/prompt/fresh.go \| sed -n 195,210p`; positive controls: V10 and V23 (the registry MkdirAll seen), V13 (pkg_docs extraction seen) | World: 4 hits, `handlers_ailang.go:57` (deny list), `setup.go:82,83` (operator-HOME defaults), `main.go:423` (examples default), none under `<state>/cache`; `cacheDir := filepath.Join(w.stateDir, "cache", episodeID)` and then `os.MkdirAll(dir, 0o700)`. Upstream `AILANG_CACHE_DIR` users: `pipeline/cache_runtime.go:50`, `pipeline/cache_store.go:67` (both `<dir>/compile`), `prompt/fresh.go:199` (`cacheBaseDir` returns `$AILANG_CACHE_DIR`, used as `…/prompts/<version>/<kind>.md`). None is under `.ailang`. |

Exact V3 command (r2; valid ERE with bare `|`):

```bash
grep -rniE 'module.?root|subdir|package.?cache|pkg_cache|AILANG_PKG' host cmd --include='*.go' | grep -v _test
grep -rn 'AILANG_CACHE_DIR' host cmd --include='*.go' | grep -v _test   # positive control
```

**Escaping note (r2, glm r1).** In the table cells above, `\|` is Markdown's escaped pipe inside a table, and every command was executed with a bare `|`. The r1 V3 cell also showed a `\|` alternation, which is invalid inside `-E`. That row is replaced by the block above. No other V-log command uses `-E` with an alternation. V19's `-run '^(A\|B)$'` is a Go regexp, executed with a bare `|`.

**UNMEASURED (named, not claimed):**
- The exact semantics of v0.52.1's `cli_allow` beyond its source comment (R-A's option).
- Whether a writable `ailang.lock` lets the loader read arbitrary directories through a `path` entry (R-B).
- `test` op behaviour under the snapshot. AC2.6 measures it at implementation.
- Behaviour on Linux CI runners of the mode-bit checks. They are pure `Lstat` mode reads, so it is expected to be identical; AC2.1 proves it.

## r2 changes (quorum round 1 → this revision)

| # | Reviewer | Fix | Section | V-rows |
|---|---|---|---|---|
| 1 | oc-kimi-k3 | The non-destructive `linkPackageCache`. Keep a matching symlink. Replace an absent or empty dir. Refuse a non-empty dir, a file or a foreign symlink, with a named operator line and no delete. The unset arm removes only a stale symlink. AC2.3 is rewritten (MUT-PC-SILENT-DELETE), and AC2.4 is tightened. (b4) is corrected: row 140 seeds `<state>/exec-cache/<id>`, not the episode `HOME` | M2 "Per episode"; Options (b1)/(b4); AC2.3, AC2.4 | V24 (every writer under `<state>/cache/E/.ailang/**`) |
| 2 | oc-glm-5-3 | The #1547 compile-cache claim is now measured. `run --policy` writes `REL/.ailang/cache/compile` despite the `.ailang/**` deny, which only stops agent writes. The M1 bullet is rewritten, and it is no longer used as a safety rationale | M1 "Confinement" | V23 |
| 3 | oc-glm-5-3 | V3 re-run with valid ERE (1 unrelated hit, `pinfetch.go:183`). Escaping note added for the other cells | Measured facts V3; Verification Log | V3 (replaced) |
| 4 | oc-glm-5-3 | AC1.2 names the 9 effect names. Multi-episode deployments union each episode's locked registry packages into one `SNAP` | AC1.2; M2 "More than one episode" | — (`workspace.go:55–66`) |
| 5 | gemini-3-1-pro | No per-episode re-walk: the digest is computed once at startup and stamped. R2 is widened to all episodes. MUT-PC-NO-RECHECK and the changed-snapshot arm are dropped. Digest rule: directories contribute `relpath NUL d`, files `relpath NUL f NUL sha256(content)`. MUT-PC-DIR-HASH is added | M2 "Snapshot digest", "Per episode"; AC2.5; R2 | — |

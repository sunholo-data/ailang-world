# w-toolchain-worktree-delta — record the worktree changes a toolchain call makes (row 184, clause 5) — r1

**Status:** r1 draft, attended 2026-10-16 (Mark attending: "yes fix them up first", after row 114 Q3 finding
F-4). Ratification is Mark's.

## §1 Problem (measured)

Row 114 Q3 (trial episode tps-a): the walk could not say where `.ailang/` (41 files) came from; the baseline
could, from the files' mtimes. World records an `ailang-run` / `ailang-check` / `ailang-cli` /
`workspace-exec` call's request and result envelope, but not the files the toolchain itself created or
changed in the worktree (SE-domain R-SE-11: inner FS effects not brokered or journaled).

| id | measurement |
|---|---|
| V1 | Handlers: Ailang.Check/CLI via `policyToolWith` (`host/broker/handlers_ailang.go:555`, `dir: h.root` = the episode worktree); Ailang.Run* via `executeRun` (`:573`, `dir: h.root`); Workspace.Exec via srt (`handlers_exec.go:207`, cwd `<worktree>/<profile.root>`, `allowWrite [Worktree, Cache]`). |
| V2 | Results are stored verbatim as `world/effect-result/v1` (`broker.go:311`); the record (`record.go:19`) holds only refs; with no finish phase the handler JSON becomes the committed invocation output (`effectful.go:209-257`). `InterfaceHash` = sha256 of the semantic id, so new result fields leave it unchanged. |
| V3 | No code snapshots or diffs the episode worktree (no `git status`, no before/after). Reusable walkers: `execManifest`/`manifestDigest` (`exec_sandbox.go:248,293`), `packageCacheDigest` (`workspace_layout.go:208`). |
| V4 | Caps: handler output 8 MiB (`handlers.go:22`); coordinator output 1 MiB (`daemon.go:719`). Breaks: `TestExecHandlerResultShape` (`handlers_exec_m2_test.go:201`) pins the key set; `execDecode` uses `DisallowUnknownFields`. Precedent: `package_cache` was added omit-when-empty so unchanged results stay byte-identical (`handlers_ailang.go:255`). |
| V5 | The specific Q3 instance is an upstream bug fixed 2026-10-12 (sunholo-data/ailang#1547, `run --policy` ignored `AILANG_CACHE_DIR`). Measured with the broker's exact argv on the trial's `zz_ls.ail`: v0.52.1 (today's tool pin) writes 36 cache files into the sandbox root, 0 to `AILANG_CACHE_DIR`; v0.54.11 writes 0 into the sandbox, 36 to `AILANG_CACHE_DIR`. Row 139's first tool-pin move removes it. |
| V6 | `why` prints per effect only verdict, effect@scope, cost, budget and request/result sizes (`why.go:752-756`), never result fields. |

## §2 Design

**D1 — `worktree_delta` in toolchain results.** For Ailang.Run/RunEnv/RunNet, Ailang.Check, Ailang.CLI and
Workspace.Exec, the broker snapshots the worktree before the subprocess and again after it, and adds to the
handler result (omit-when-empty, so a call that changes nothing stays byte-identical, V4 precedent):

```json
"worktree_delta": {"changes": [{"path": "zz_ls.ail", "op": "created", "size": 812, "sha256": "…"},
                               {"path": ".ailang/cache/compile/manifest.json", "op": "modified", "size": 25770, "sha256": "…"}],
                   "truncated": false, "total": 2, "basis": "git-status"}
```

`op` ∈ created/modified/deleted; `sha256`/`size` of the after-state (absent for deleted). At most 256 entries,
then `truncated: true` with `total` (bounded far under the 1 MiB output cap).

**D2 — the snapshot.** In a git worktree (every World episode today: `session new` creates one): `git status
--porcelain=v1 -z --untracked-files=all --ignored=no` plus a sha256 of each listed path, before and after; the
delta is the set difference by (path, sha256). Gitignored paths are out of scope and stated (`basis`), which
keeps virtualenvs out. Outside git: a bounded lexical walk (reusing `execManifest`'s walker) with a file cap;
over the cap ⇒ `basis: "walk-capped"`, never silent. `.git` is never read beyond what `git status` reads. The
snapshot runs host-side (broker), outside the sandbox and outside the child's timeout.

**D3 — `why` shows it.** `why <entry>` prints, per effect with a delta, `worktree: N created, M modified, K
deleted` and the first paths; `--json` includes the list. The walk reaches it with `object get <output>` too.

**D4 — tests that pin shapes** (`TestExecHandlerResultShape`, `execDecode`) learn the optional field; an
unchanged call's result stays byte-identical (asserted).

**D5 — not in scope:** brokering inner effects (deny/allow per write) — that stays R-SE-11's later row; this
row RECORDS what happened. Episode-end summaries (one delta for the whole episode) are a possible add-on, not
built here.

**Why not a package (S3)?** Host-boundary capture of subprocess side effects; no `world/` change.

## §3 Acceptance criteria

| AC | criterion |
|---|---|
| AC1 | An `ailang-run` of a program that writes a file records it as `created` with size and sha256; a run that writes nothing has no `worktree_delta` and a byte-identical result to today. |
| AC2 | `workspace-exec` of a command that modifies a tracked file records `modified`; deleting records `deleted`. |
| AC3 | Over 256 changes ⇒ 256 entries, `truncated: true`, `total` exact. |
| AC4 | Gitignored paths are not listed; `basis` says `git-status`; a non-git worktree uses the capped walk and says so. |
| AC5 | `why <entry>` shows the delta line; existing why/provenance tests pass. |
| AC6 | On the trial's case (v0.52.1 tool pin, `zz_ls.ail`), the 36 cache files appear in the run's delta — a walk can answer Q3's `.ailang/` half from the store. |

## §4 Decisions for Mark

- **D-184-1:** D2's git-status basis (excludes gitignored paths, fast) vs always a full walk (sees ignored
  paths such as caches, but slow on worktrees with virtualenvs). Recommendation: git-status, with the basis
  recorded.

## §5 Risks

| risk | mitigation |
|---|---|
| Snapshot cost on large worktrees | git status is incremental; sha256 only for listed (changed/untracked) paths; measured at build time on platform-source |
| A concurrent writer between snapshots | episodes are single-agent; calls are serialized per episode |
| Result shape change breaks clients | omit-when-empty; shape tests updated |

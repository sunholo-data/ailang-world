# w-agent-search-and-read-cap (queue row 228): a World agent can find things — r0

**Status:** r0, written attended 2026-10-28 through World itself (episode `h2610280901` on the
trial rig; Mark: *"Write the design doc for row 228 (w-agent-search-and-read-cap) in
design_docs/planned/. Start from the charter row and friction 34–36. Grep first: you have grep
and ls now. Check every premise against the code. Leave the rulings open with defaults, as you
did for 219."*). It is a design, not a build: no code exists yet. It is not quorum-reviewed. The
pick-time designer re-measures every premise marked **verify** in §1 (they sit upstream in
AILANG's `policy-tool`, or in the agent host, which this session cannot read) and runs the
quorum as usual. Roadmap: [`w-all-agent-work-through-world`](w-all-agent-work-through-world.md)
§2 row 6d, **ahead of 226** so the scorecard measures confinement, not missing tools. Rulings it
needs: §5 R1–R5.

## §0 What is wrong today

Row 228's own words: **"A WORLD AGENT CAN'T FIND THINGS."** Measured in `h2610271506`
(friction 34–36) and re-measured in this episode, the first one run with the quick half:

- **No search or listing (friction 34).** In `h2610271506` the agent could read a file only by its exact
  path. Finding row 219 in the 6,289-line charter took two forks reading it in chunks (~400k
  tokens), and guessed filenames missed (`host/daemon/health.go`; the website guide as `.md`
  and `.mdx`).
  **With the quick half** (`grep` and `ls` in the exec profile, Mark's attended change), this
  episode found row 228 (`design_docs/world-mission.md:6291`) in **one** call, friction 34–36 in
  one call, and the row 219 design in one more. The quick half meets the exit criterion's
  3-call bound for search. It leaves four gaps (P14–P17): its grammar is narrow, its output is
  capped at 64 KiB a stream by head and tail, not refused, it lists one directory only, and it
  exists only in profiles that declare it.
- **No way to tell how old the worktree is (friction 35).** The worktree is cut from the main
  checkout's `HEAD` when the episode is minted. Neither the context block nor `whoami` names
  that commit, and the agent has no git, so a premise missing from the brief can't be told
  apart from a stale worktree. **This episode is no better:** its context block (top of this
  session) names the episode, the worktree path, the TTL and the grants, but no commit. This
  doc can't say what it was written against (§1 P10).
- **One read can overflow the agent (friction 36).** `ailang-read` of `docs/QUICKSTART.md`
  returned 102k characters, which the host spilled to a file the agent couldn't open. Line
  ranges don't help when single lines are huge: the charter's rows are thousands of characters
  each, and row 228 alone is about 2,000 (its line in this episode's first `grep`).

**Friction 38** (added with this doc's PR) is the gap §3.5 leaves to row 218: this session was
asked to commit, push and open a PR, and has no tool that does any of the three.

**Exit criterion (row 228):** re-run `h2610271506`'s two tasks (design a row from the charter;
refresh a status section) in a fresh World pane. Locating a charter row takes at most 3 calls,
no read exceeds the cap, and the agent states its worktree's base and lag without guessing.

## §1 Premise verification log

Each premise is **measured** (read in this worktree's code on 2026-10-28, or observed in this
episode's own calls), **quoted** (from a charter row or friction entry, not re-measured), or
**verify** (upstream AILANG or agent-host behaviour this session can't read).

| # | premise | source | status |
|---|---|---|---|
| P1 | `ailang-read` is an AILANG plan: it admits only `path`, `offset` and `limit` (`argKeyAllowed`, a contracted Z3 core function with its own `ensures`), and emits one `Workspace.Read` effect (scope `worktree`, cost 1) whose payload is the policy-tool `read` request `{op, path, offset?, limit?}`. It has no finish phase, so policy-tool's answer is the call's result. | `packages/se-tools/se_tools/read.ail` | measured |
| P2 | The 2000-line clamp on `limit` is policy-tool's (ailang#1872), and it is silent. Nothing in World caps characters: a no-range read returns the whole file. | `read.ail` (`lineArgOk` comment); no `max_chars`/char cap anywhere in `host/` or `cmd/` (grep) | measured / **verify** the clamp |
| P3 | The broker runs `Workspace.Read` as policy-tool op `read` only (`fixedToolOps`). It decodes the payload, refuses a key that case-folds onto `op` or `flags`, re-encodes it and passes it to `ailang policy-tool --policy <episode policy>`. The response must be an object with a boolean `ok` and none of the reserved keys, and then it is stamped. **The broker holds the whole response in memory before the agent sees it.** | `host/broker/handlers_ailang.go` (`fixedToolOps`, `policyToolRequestKeys`, `Execute`, `stamp`) | measured |
| P4 | policy-tool v0.54.15 refuses any request key that is not one of its own. A new argument therefore can't be passed through, and the broker must consume it. | `handlers_ailang.go` comment on `policyToolRequestKeys` | quoted / **verify** |
| P5 | The episode policy confines reads by `fs_sandbox` = the worktree root, with symlinks out refused by AILANG's policy layer. World renders `fs_deny_write` only: there is no deny-read list, so whatever the episode may read, a search may read. | `RenderEpisodePolicy`, `loadSummary` (`fs_sandbox` must equal the root); the `ailang-read` descriptor baseline | measured |
| P6 | Precedent for a new upstream op: `delete` was added to policy-tool upstream (ailang#1875) and given its own effect, `Workspace.Delete`. CLAUDE.md routes language gaps upstream: "no local workarounds, no vendored forks". | `handlers_ailang.go`; `CLAUDE.md` | measured |
| P7 | A se-tools tool is a `.ail` module under `packages/se-tools/se_tools/`, a `transitions.json` entry and a descriptor in `host/coordinator/testdata/se-tools-descriptor-baseline.json`, held equal by the coordinator's tests. The se-tools preset grants the union of `transitions.json`'s effects (`session_new.go`). A publish is Mark's act (row 228). | those files | measured |
| P8 | `whoami --json` is `GET /v1/session/self`: `{episode, worktree, expiresAt, secondsLeft, grants, entries?, addr, surface}`. It has no commit, branch or repo field. `worktree` is `workspace-root/<episode>`, a path join, not a git query. | `host/daemon/session_self.go` | measured |
| P9 | The context block is `bin/world-context`, generated by `hosts init` (row 223). It runs `whoami --json` and prints the episode/worktree/TTL/daemon line, the tool map, the grants, the surface line and the boundary line. Its guard requires only `.episode`, `.secondsLeft` and `.grants`, so new fields are additive. | `cmd/ailang-worldd/hosts_awareness.go` | measured |
| P10 | The worktree is `git worktree add -b <branch> <dir>` from `--repo` with no start point, so its base is **whatever the main checkout's `HEAD` was at mint**, not `dev`. `episode start`'s default branch is `world/<episode>`. Nothing records the base: not the credential, not the episode, not `/v1/health`. | `host/worktree/worktree.go` (`Add`); `cmd/ailang-worldd/episode.go`, `session_new.go` | measured |
| P11 | The worktree's `.git` is a **file**, `gitdir: /Users/…/ailang-world/.git/worktrees/h2610280901`. The object store and every ref, including other episodes' `world/*` branches, live in the main repository outside the worktree. The agent can read the `.git` file but nothing it points to (absolute paths refuse at the plan, P1). | `ailang-read .git` in this episode; `worktree.IsLinked` | measured |
| P12 | World already runs read-only git against an episode's worktree in two hardened ways. (a) `host/worktree.Output`: GIT_* and registry variables scrubbed, 30 s timeout, stdout over 8 MiB refused, not truncated (used by `surface-index`). (b) The broker's `runWorktreeGit`: `-C <root> -c core.fsmonitor=false -c core.quotePath=false`, env exactly `PATH=/usr/bin:/bin LC_ALL=C GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null GIT_OPTIONAL_LOCKS=0 GIT_TERMINAL_PROMPT=0`, bounded (used by worktree deltas and protected-path records, which already take `rev-parse HEAD` as `base`). The protected check refuses when any `filter.*` config exists (`WorktreeFilterRefused`), because a clean filter runs on `status`/`diff`. | `host/worktree/worktree.go`; `host/broker/worktree_delta.go`, `exec_protected_check.go`, `exec_protected_record.go` | measured |
| P13 | `dev` already has one meaning in World: `surface-index --rev` defaults to the **local** branch `dev`. | `cmd/ailang-worldd/surface_index.go` | measured |
| P14 | The quick half's grammar, as measured in this episode: `grep` runs as `grep -rnI --color=never` and admits `-l` and `--include=GLOB`, but refuses `-n` (already fixed), `-E` and `--exclude=`. A basic regex with `\|` works. `ls` runs as `ls -la` with one relpath and no recursion. The profile has no `git`: `workspace-exec git` is refused as "not a command of exec profile". | this episode's `workspace-exec` calls (World log) | measured |
| P15 | `workspace-exec` keeps the first 8 KiB and the last 56 KiB of each stream (`execHeadBytes`, `execTailBytes`), with a 64 MiB runaway kill. A grep over the charter can return ~64k characters of stdout, middle elided: bounded, but over any agent-sized cap, and truncated rather than refused. | `host/broker/handlers.go`, `handlers_exec.go` | measured |
| P16 | Exec runs under srt with the home directory unreadable (this session's context block). The main repository is under `$HOME`, so `git` in the exec profile would fail on the gitdir unless the sandbox read set opened the main `.git`. That would also open every other episode's branches and worktree metadata. | context block; P11 | measured / **verify** the denial |
| P17 | A broker `Git.Commit` handler exists (`git commit --no-edit -m`, empty HOME, fixed identity). It isn't in the se-tools preset, no se-tools tool emits it, and nothing pushes. Commit, push and PR are row 218's (`episode finish`, an outward effect Mark approves). | `host/broker/handlers_git.go`; `session_new.go`'s preset list; roadmap §2 row 3 | measured |
| P18 | Claude Code spills an MCP tool result above its output limit (`MAX_MCP_OUTPUT_TOKENS`, about 25k tokens by default) to a file, which a confined agent has no tool to open. The 102k-character read in friction 36 was over it. | friction 36; the agent host | quoted / **verify** the number |
| P19 | policy-tool is Go, so a regex it compiles is RE2: linear time, no backtracking blow-up. | AILANG | **verify** |

**P10 and P11 together are friction 35's cause.** The base is a fact only the main repository
holds, and it sits outside everything the agent can read. It must come from the daemon, which
already runs bounded git (P12).

## §2 Conflict surface

- **Row 223 (built: `hosts init`, the context block, `whoami`).** 228 adds fields to
  `/v1/session/self` and one line to `world-context`. The guard and the existing lines stay
  (P9). The herdr plugin reads the same `whoami --json` and ignores unknown fields.
- **Row 227 (error clarity).** The read-cap refusal (§3.2) and every refusal of the new tools use
  227's three-part form: what failed, why, and the next call to make.
- **Row 156 (protected test paths, in-sprint)** works in the broker's exec and protected-path
  code. §3.4 reuses its `filter.*` refusal and `runWorktreeGit`, and changes neither. M3 waits
  for 156 to land, or the two sequence through one controller.
- **Row 218 (`episode finish`)** owns commit, push and PR. §3.4 is read-only by construction and
  must never grow a write verb. Friction 38 is 218's.
- **Row 217 (presets)** is where a repo's "dev" ref would be named per preset (R3).
- **Row 226 (scorecard)** counts "checks the agent could run" and unclear errors. 228 lands first
  so that a missing search tool isn't scored as confinement (roadmap §2).
- **Row 225 (`ailang-worldd run`).** An adapter-less agent under the wrapper has its own shell
  and doesn't need these tools. Nothing here is required by 225.
- **The quick half (the operator's exec profile).** It stays the operator's file: 228 doesn't
  edit or remove it (R5).
- **Upstream AILANG (`policy-tool`).** R1 (A) needs two new ops and a pin bump, routed as an
  issue on `sunholo-data/ailang` with a `mission-control` note, as with ailang#1875 (P6).

## §3 The design

### §3.1 `workspace-search` and `workspace-list`

Two se-tools tools, each one plan emitting **one `Workspace.Read` effect, cost 1**, like
`ailang-read` (P1, P7). The broker admits two more policy-tool ops under that effect,
`fixedToolOps[Workspace.Read] = {read, search, list}` (R1 A), so confinement stays in AILANG's
policy layer (P5) and is not re-implemented in Go.

**`workspace-search`** args (L-ARGS: any other key refuses):

| arg | type | default | meaning |
|---|---|---|---|
| `pattern` | string, 1–512 bytes | required | what to find |
| `regex` | bool | `false` | `false` is a fixed string; `true` is RE2 (P19) |
| `ignore_case` | bool | `false` | |
| `path` | relpath | `.` | a file or a directory under the worktree (`pathOk`, P1) |
| `include` | glob | none | basename glob, e.g. `*.go`, `*.md*` |
| `max_results` | int 1–500 | 100 | matches returned |

Result: `{ok, matches:[{path, line, col, text}], total_matches, files_searched, truncated}`.
`text` is the matched line cut to a **240-character window around the match**, with `…` at
cut ends, so a charter row is a pointer and never a payload. The whole result is also held
under the read cap (§3.2): when it would exceed it, fewer matches are returned and `truncated`
is true, with the next call named (`narrow path or include`). Binary files are skipped, and so
is `.git`. Every other path the episode may read is searched (P5). An open question is whether
to skip `node_modules`-sized trees (§6 Q1).

**`workspace-list`** args: `path` (default `.`), `depth` (1–6, default 2), `max_entries`
(1–1000, default 300), `include` (glob). The result is `{ok, entries:[{path, type: file|dir|symlink,
size}], total, truncated}`, sorted by path. A symlink is listed, never followed.

Both answer in policy-tool's own shape, `ok:false, refused` with rc 0, for every refusal. Each call commits
one World entry, and `ailang-worldd why <entry>` explains it like a read.

### §3.2 `ailang-read`'s character cap

`read.ail` admits a fourth key, `max_chars` (int, 1,000–200,000), and `argKeyAllowed`'s
`ensures` grows to match. That changes a contracted function and its named tests (P1). **The
broker enforces it** (P3, P4). It removes `max_chars` from the payload before policy-tool sees it,
and adds `max_chars` to `policyToolRequestKeys` so a case-folded alias refuses. Once the
response is in hand, it measures `content`. The cap applies **whether or not `max_chars` is
given** (R2 A): the default is **40,000 characters**, chosen to sit under P18's limit with room for
the envelope.

Over the cap, the read is **refused, not truncated**: `ok:false`, three-part (row 227).

```
refused: docs/QUICKSTART.md is <C> characters over <L> lines; the cap is 40,000
  why: one read may not exceed max_chars (default 40,000) so the result fits the agent's context
  next: ailang-read {path, offset: 1, limit: <n1>}, then offset <o2>, <o3> … (each under the cap)
```

The `next` line is computed from the response's own line lengths: the largest `limit` from
`offset` that fits, repeated to the end. When **one line alone** exceeds the cap (a charter row
of about 2,000 characters doesn't, but a minified file can), `next` names that line number and
says `workspace-search` will return a window of it. A character offset is §6 Q2. The refusal still
costs one `Workspace.Read` and is recorded: the read happened, and the result was withheld.
Under the cap, nothing changes, and `total_lines` already comes back from policy-tool.

### §3.3 Base and lag in `whoami` and the context block

`/v1/session/self` gains one additive object, computed per request by the daemon through
`host/worktree.Output` (P12a) against the episode's worktree. It is omitted for a degraded
(non-git) episode:

```json
"git": {"branch": "world/h2610280901", "head": "<sha>", "base": "<sha>",
        "base_subject": "<subject>", "base_date": "<RFC 3339>",
        "ref": "dev", "ref_head": "<sha>", "behind": <n>, "ahead": <m>, "dirty": <bool>}
```

`base` = `git merge-base HEAD <ref>`, `behind` = `rev-list --count base..<ref>`, `ahead` =
`rev-list --count base..HEAD`, `dirty` = `status --porcelain` non-empty, all with
`runWorktreeGit`'s environment and `-c core.fsmonitor=false`. Where `<ref>` is "dev" is R3
(default: the local branch `dev`, as `surface-index` reads it, P13). If any call fails or takes
more than 2 s, the object is `{"error": "<one line>"}` and `whoami` still answers: the context
block must not fail because git is slow. `world-context` prints one more line after the
episode line:

```
Base <7 hex> (<date>, "<subject>") · dev is <n> commits ahead of it · worktree has uncommitted changes
```

and, on error, `Base unknown: <error>`, so an agent never guesses. `whoami` (human form) prints
the same line.

### §3.4 Read-only git: `workspace-git`

A se-tools tool whose plan emits one effect, **`Workspace.GitRead`** (scope `worktree`, cost 1),
a **new grant** (R4 A). It reads the shared object store outside the worktree (P11), which is a
different authority from reading worktree files, as `Workspace.Delete` is from writing (P6). The broker runs it,
not the exec sandbox (P16), through `runWorktreeGit` (P12b), with:

- **Subcommands:** `log`, `show`, `diff`, `status`, and nothing else.
- **A fixed flag grammar** per subcommand. For `log`: `--oneline`, `-n 1–200`, `--stat`,
  `--format` from an enum. For `show`: `--stat`, `--name-only`. For `diff`: `--stat`, `--name-only`,
  `--cached`. For `status`: `--porcelain` only. Paths after a literal `--`, each `pathOk`.
  `-c`, `--output`, `--ext-diff`, `--exec-path`, `--git-dir`, `--work-tree` and any flag not
  listed refuse at the plan.
- **Always added by the broker:** `--no-pager` (as `-c core.pager=cat`), `--no-ext-diff`,
  `--no-textconv`, `-c core.fsmonitor=false`, `-c diff.external=`. The worktree's
  `.gitattributes` is agent-writable, so no attribute may select a driver.
- **Refused outright** when `git config --get-regexp ^filter\.` finds any filter, reusing the
  protected check's refusal (P12). A clean filter runs on `status` and `diff`.
- **Revisions limited** to `HEAD`, the base, the R3 ref, and commits reachable from `HEAD`
  (`merge-base --is-ancestor`), each given as a 7–40 hex sha or one of those names. Another
  episode's `world/*` branch, any other ref, and a `rev:path` outside the worktree's own tree
  refuse. Episodes share one object store, and this tool must not become a window into the
  others' work.
- **Output** is held to the same 40,000-character cap as §3.2, refused with a `next` (`-n`
  fewer, `--stat`, or a path).

### §3.5 What this design does not do

No write verb of any kind in §3.4: commit, push and PR stay row 218's, as an outward effect Mark
approves (P17, friction 38). No change to the exec profile (R5). No network: `behind` is against
the local ref as last fetched by whoever fetches the main checkout, and §3.3's line says "dev"
in that sense.

## §4 Failure modes

| # | condition | what the agent sees |
|---|---|---|
| F1 | read over the cap | §3.2's three lines; one entry recorded |
| F2 | one line over the cap | the line number and a `workspace-search` window as `next` |
| F3 | search hits `max_results` or the cap | the matches it holds, `truncated:true`, `next: narrow path or include` |
| F4 | policy-tool too old for `search`/`list` (pin not bumped) | `refused: workspace-search needs ailang ≥ <v>`, `why: the pinned policy-tool has no search op`, `next: the operator bumps ~/.pinned-ailang-tools` (a broker startup check, not a per-call 500) |
| F5 | git slow or failing in `whoami` | `"git":{"error":…}`, the `Base unknown:` line; `whoami` exit unchanged |
| F6 | `workspace-git` with `filter.*` configured | the protected check's refusal text |
| F7 | a revision outside §3.4's set | `refused: <rev> is not HEAD, the base, dev or an ancestor of HEAD` |
| F8 | session without `Workspace.GitRead` | the ordinary budget refusal, naming the grant |

## §5 Rulings (proposed; recorded as D-WORLD-NN when made)

- **R1, where search and list run.** **(A, default)** as two new ops of AILANG's `policy-tool`
  (an upstream issue, a release and a pin bump), so confinement stays in one place (P5, P6).
  **(B)** a broker-native Go walk under the worktree root, with symlink and `.git` handling
  re-implemented and kept equal to the policy layer by test. It is faster to land but makes two
  confinement implementations.
- **R2, when the read cap applies.** **(A, default)** always, 40,000 characters unless
  `max_chars` (1,000–200,000) says otherwise. **(B)** only when `max_chars` is given, so the
  default reads are unchanged and friction 36 recurs for any agent that doesn't ask.
- **R3, what "dev" means for lag.** **(A, default)** the repository's local branch `dev`, the
  same meaning as `surface-index --rev` (P13), with a per-preset override when row 217 lands.
  **(B)** `refs/remotes/origin/dev` as last fetched, with the fetch time shown.
- **R4, the authority for read-only git.** **(A, default)** a new grant, `Workspace.GitRead`,
  in the se-tools preset, and a separate line on the mint prompt. **(B)** an op under
  `Workspace.Read`, so no new grant, but a worktree read and an object-store read become one
  authority. **(C)** defer §3.4 until after row 226 and ship §3.1–§3.3.
- **R5, the quick half after the native tools land.** **(A, default)** leave `grep`/`ls` in the
  operator's profile: it's the operator's file, and an exec command is a fine second path.
  **(B)** remove them when M2 lands, so search has one recorded shape (`Workspace.Read`).

## §6 Open questions

1. **Generated trees.** Should `workspace-search` skip what `.gitignore` names (needs git, so a
   broker pre-filter under R1 A) or a fixed list (`node_modules`, `vendor`, `.ailang`)? Default:
   neither in r0. `include` and `path` narrow, and `files_searched` shows the cost.
2. **A character offset for `ailang-read`.** F2 sends an agent to `workspace-search` for a
   window. A `char_offset`/`char_limit` pair would read any byte range, but needs upstream support.
   Default: not in 228. Revisit if F2 shows up in friction.
3. **Exec output under the cap.** P15's 64 KiB per stream is already over §3.2's cap. Should
   `workspace-exec`'s capture shrink, or its result refuse like a read? Default: no change in
   228. Exec output is a log, not a file, and head and tail is the right shape for it. It is
   recorded for row 227.
4. **Base at mint.** Should `episode start` also record `HEAD` at mint on the credential, so the
   base is a recorded fact and not one recomputed later? Default: no. `merge-base` is stable for an
   episode that cannot commit, and 218's commit path is where a recorded base matters.

## §7 Milestones

| M | what | gate |
|---|---|---|
| **M0** | §3.2: `read.ail` `max_chars` (contract, tests, descriptor), broker enforcement and refusal with computed `next` | `./scripts/verify_ail.sh`; Go tests over a 102k-character fixture and a single-huge-line fixture; ACs 3–5 |
| **M1** | §3.3: `/v1/session/self.git`, `whoami`, `world-context`'s Base line | Go tests with a fixture repo (dev moved 3 commits; a degraded episode; a git that sleeps 5 s); `hosts_awareness_test.go`; ACs 6–7 |
| **M2a** | upstream: the AILANG issue for policy-tool `search`/`list` and its release (R1 A), plus a `mission-control` note | the released ops measured against §3.1's shapes |
| **M2b** | §3.1: `workspace-search.ail`, `workspace-list.ail`, broker op allowlist, descriptors, F4's startup check; **publish is Mark's act** | `verify_ail.sh`; coordinator descriptor tests; daemon se-tools e2e; ACs 1–2, 8 |
| **M3** | §3.4 `workspace-git`, **after row 156**; its own judge (confinement-adjacent) | Go tests: every refused flag and revision, a hostile `.gitattributes` and `filter.*` config, the env pin; ACs 9–11 |
| **M4** | the exit re-run, attended | AC12 |

M0 and M1 need nothing upstream and no publish of new tools. M0 changes `ailang-read`, so it
is a se-tools publish. Together they close friction 35 and 36 and can land before M2a returns.

## §8 Acceptance criteria (no model call; an operator or CI can check each)

1. **Search finds a row in one call.** With the real charter checked into a fixture,
   `workspace-search {pattern: "w-agent-search-and-read-cap", path: "design_docs/world-mission.md"}`
   returns a match at line 6291 whose `text` is ≤ 240 characters, and the whole result is
   ≤ 40,000 characters.
2. **List is bounded.** `workspace-list {path: ".", depth: 6}` on this repo returns
   `truncated:true` with exactly `max_entries` entries, no entry under `.git/`, and every symlink
   typed `symlink` with no entries beneath it.
3. **The cap refuses.** `ailang-read {path: <a fixture over 100,000 characters, friction 36's size>}` returns `ok:false`,
   with `refused`, `  why:` and `  next:` lines. Every `{offset, limit}` the `next` line names
   reads successfully and under 40,000 characters, and together they cover every line once.
4. **The cap is consumed, not forwarded.** The bytes the policy-tool stub receives carry no
   `max_chars` key. `{"MAX_CHARS": 5}` refuses as an alias, and `max_chars: 999` and `200001`
   refuse at the plan.
5. **Under the cap, unchanged.** A 10-line read's response bytes equal today's, apart from the
   stamp.
6. **Base and lag.** In a fixture repo where `dev` moves 3 commits after the worktree is
   added: `whoami --json | jq .git` shows `behind: 3`, `ahead: 0`, `base` equal to the
   pre-move `dev`. `world-context` prints `Base <7 hex> (…) · dev is 3 commits ahead of it`.
7. **Base never blocks.** With a git stub that sleeps 5 s, `whoami` answers within 3 s with
   `git.error` set, and `world-context` prints `Base unknown:` and exits 0. A degraded episode
   has no `git` key.
8. **Recorded like a read.** Each of `workspace-search`, `workspace-list` and a refused
   `ailang-read` commits one entry with effect `Workspace.Read`, cost 1, and `why <entry>`
   names the tool and its args.
9. **Git is read-only.** A Go test pins `workspace-git`'s argv set: no argv the plan can produce
   contains `commit`, `push`, `fetch`, `checkout`, `reset`, `config` (except the broker's own
   `--get-regexp ^filter\.`), `-c` from the agent, `--output` or `--ext-diff`.
10. **Git is confined to this episode.** In a repo with two episode worktrees, `workspace-git
    log world/<other>` and `show <other-branch-sha>` (not an ancestor of `HEAD`) refuse with F7.
    A `.gitattributes` naming a `textconv` driver configured in the repo's config shows plain
    diff output, and the driver's marker file is never created.
11. **Git env.** The child environment equals `runWorktreeGit`'s six variables exactly. A
    `GIT_DIR` decoy in the daemon's environment changes nothing (the existing
    `worktree_delta_test.go` arm, extended).
12. **The exit re-run** (attended, M4). In a fresh World pane, on `h2610271506`'s two tasks:
    every charter row the task needs is located in ≤ 3 calls, no tool result in the World log
    exceeds 40,000 characters of content, and the agent's first message states its base and lag
    from the context block.

## Sources

- Charter row 228 (`design_docs/world-mission.md` line 6291), and rows 156, 217, 218, 219, 223,
  225, 226 and 227. The mission-log note at line 2011.
- [`w-all-agent-work-through-world`](w-all-agent-work-through-world.md) §2 (rows 3, 2b, 2c, 2d,
  6d) and §3; [`w-worldd-home-service`](w-worldd-home-service.md) for this doc's shape.
- [`friction.md`](../verification/world-dogfood/friction.md) entries 5, 34, 35, 36, and 38
  (added with this PR).
- Code: `packages/se-tools/se_tools/read.ail`, `packages/se-tools/transitions.json`,
  `host/coordinator/testdata/se-tools-descriptor-baseline.json`;
  `host/broker/{handlers_ailang,handlers_git,handlers,handlers_exec,worktree_delta,exec_protected_check,exec_protected_record,exec_profile}.go`;
  `host/daemon/session_self.go`; `host/worktree/worktree.go`;
  `cmd/ailang-worldd/{hosts_awareness,session_new,episode,surface_index}.go`; `CLAUDE.md`.
- This episode's own `workspace-exec` and `ailang-read` calls (World log, episode `h2610280901`)
  for P11 and P14.

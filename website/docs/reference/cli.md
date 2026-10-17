---
title: CLI reference
sidebar_position: 1
description: Every ailang-worldd and world-publish command and flag, taken from the binaries' own help output.
---

# CLI reference

Two programs, both built from this repository. Help text below is the binaries' own output,
brought up to date with the help strings in `cmd/ailang-worldd` and `cmd/world-publish` on
2026-10-16 (`v1.0.0-rc.1`). Sample outputs show placeholder paths.

## `ailang-worldd`

```text
ailang-worldd — AILANG World local daemon (loopback only)

Usage:
  ailang-worldd serve --db <path> [--bind host:port] [--ailang-bin <path>]
                      [--workspace-root <dir> --tool-ailang-bin <path>]
                      [--examples-dir <dir>]
                      [--workspace-module-root REL] [--workspace-episode-module-root EP=REL ...]
                      [--workspace-package-cache <dir>]
                      [--run-allow-caps Env,Net,Declassify]
                      [--run-net-allow 127.0.0.1:PORT ...] [--run-net-allow-http]
                      [--exec-profile <file> ... --exec-sandbox <dir>]
                      [--exec-episode-project EP=PROJECT ...] [--exec-node <path>]
                      [--exec-max-output-bytes N]
                      [--proof-key-file <file> --proof-required-identity NAME ...]
  ailang-worldd [--addr <url>] health
  ailang-worldd [--addr <url>] head
  ailang-worldd [--addr <url>] world get <ref>
  ailang-worldd [--addr <url>] object get <ref> [--payload]
  ailang-worldd [--addr <url>] object find <semanticId> [--after <ref>] [--limit N]
  ailang-worldd [--addr <url>] log get <index>
  ailang-worldd [--addr <url>] log range --from N [--limit M]
  ailang-worldd [--addr <url>] registry get <name>
  ailang-worldd [--addr <url>] log tail [--from N] [--follow] [--interval 1s] [--raw]
  ailang-worldd [--addr <url>] commit --file <commit.json> [--session <file|token>]
  ailang-worldd [--addr <url>] outcome record --file <outcome.json|-> [--session <file|token>]
  ailang-worldd [--addr <url>] tools list [--session <file|token>] [--json]
  ailang-worldd [--addr <url>] call <tool> [--session <file|token>]
                    [--arg k=v]... [--arg-json k=<json>]... | --json <obj>|@file|-
                    [--json-out] [--strict]
  ailang-worldd [--addr <url>] why <index|head|sha256:<ref>|a2a:<id>|mcp:<id>|rest:<id>|->
                    [--result <file>] [--scan N] [--json]
  ailang-worldd [--addr <url>] provenance [--since <entry>] [--episode <ep>] [--scan N]
  ailang-worldd setup [--interpreter-dir <dir>] [--tools-dir <dir>] [--db <path>]
                    [--workspace-root <dir>] [--from-dir <dir>] [--replace]
  ailang-worldd [--addr <url>] doctor [--db <path>] [--workspace-root <dir>]
                    [--interpreter-dir <dir>] [--tools-dir <dir>]
                    [--examples-dir <dir>] [--online]
  ailang-worldd session new <episode> [--db <path>] [--workspace-root <dir>]
                    [--repo <dir>] [--preset se-tools] [--grant EFFECT=SCOPE:BUDGET]...
                    [--budget 50] [--ttl 3600] [--out <file>] [--branch <name>]
                    [--allow-non-git]
  ailang-worldd session list [--db <path>] [--workspace-root <dir>] [--episode <ep>] [--json]
  ailang-worldd session revoke --db <path> <credential_id-hash> | --episode <ep>
  ailang-worldd session mint --db <path> --episode <ep> --grant EFFECT=SCOPE:BUDGET...
                    [--ttl 3600] [--out <file>]
  ailang-worldd [--addr <url>] demo seed --db PATH --session FILE --ailang-bin PIN
                    --proof-key-file KEY --manifest FILE [--addr URL]
  ailang-worldd help [<verb>]

  '<verb> --help' and 'help <verb>' print a verb's help (exit 0).

Session credential (tools, call, commit, outcome): --session <file> (a file holding
the 64-hex token, mode 0600) or the token itself (warns: visible on argv),
else $WORLD_SESSION. The token is never printed.

Global client flag:
  --addr <url>   base URL of the daemon (default http://127.0.0.1:7644).
                 Applies to client verbs only; it is NOT a 'serve' flag.

serve flags:
  --db <path>          world store database (required)
  --bind host:port     loopback listen address (default 127.0.0.1:7644);
                       a non-loopback host is refused — there is no override
  --ailang-bin <path>  interpreter to archive and pin at startup (optional)
  --workspace-root <dir>
                       episode worktrees live at <dir>/<episode> (made by the
                       operator with git worktree add before session mint); it
                       must not contain the store, its archive, its rendered
                       policies or its tool cache — startup refuses otherwise
  --tool-ailang-bin <path>
                       AILANG binary the workspace tools run (must be
                       AILANG v0.54.15); archived and hash-verified like
                       --ailang-bin. The Workspace.*/Ailang.* tools are served
                       only when both this and --workspace-root are set.
                       Workspace.Exec (workspace-exec) is bound with them; it
                       refuses every call "no exec profile configured" unless
                       --exec-profile is given
  --examples-dir <dir> AILANG examples corpus examples-search reads (passed
                       to the tool as AILANG_EXAMPLES; World never falls back
                       to a corpus in the binary). Default: ~/.ailang/examples
                       when it exists, else examples-search refuses "no
                       examples corpus configured". Must be outside
                       --workspace-root
  --workspace-module-root REL
                       the AILANG sandbox and module root of every episode is
                       <workspace-root>/<episode>/REL, so bare imports among a
                       project's modules resolve (default . = the worktree).
                       REL is clean and relative, with no .. segment, and
                       must already exist as a real directory or that
                       episode's tools are refused; it is independent of an
                       exec profile's root. Needs --workspace-root
  --workspace-episode-module-root EP=REL
                       episode EP's module root, overriding the default
                       (repeatable)
  --workspace-package-cache <dir>
                       a read-only snapshot of registry packages
                       (<ns>/<name>/<ver>/ailang.toml; QUICKSTART section 11),
                       linked as every episode's package cache so pkg/ imports
                       resolve. It must lie outside --workspace-root and the
                       state dir and hold no symlink and nothing writable
                       (chmod -R a-w); serve refuses to start as uid 0 with
                       it. Restart serve after rebuilding it
  --run-allow-caps Env,Net,Declassify
                       extra capabilities an ailang-run may request beyond
                       IO and FS (default none). An Env run is the effect
                       Ailang.RunEnv, a Net run Ailang.RunNet; each needs
                       its own session grant
  --run-net-allow 127.0.0.1:PORT
                       a loopback IP:PORT a Net run may reach (repeatable;
                       required with Net). Every other host, port and
                       redirect hop is refused; a bare host, a name, or a
                       non-loopback address is refused at startup
  --run-net-allow-http allow plain http to the --run-net-allow pairs
  --exec-profile <file>
                       an operator exec profile (world/exec-profile/v1 JSON;
                       repeatable, one per project): the commands
                       workspace-exec may run, sandboxed by srt. Needs the
                       workspace tools and --exec-sandbox. Read once at
                       startup; it, srt, node and every cache seed must lie
                       outside --workspace-root and the state dir. Each
                       profile must pass the startup probe (seven arms,
                       arm1-write-inside … arm7-toolchain) or serve refuses
  --exec-episode-project EP=PROJECT
                       run episode EP's workspace-exec under the profile
                       whose project is PROJECT (repeatable; optional with
                       one profile, which every episode then uses)
  --exec-sandbox <dir> the node_modules holding @anthropic-ai/sandbox-runtime
                       0.0.78 (its dist/cli.js must hash to the pin);
                       World archives the whole tree and re-verifies it
                       before every call. Required with --exec-profile
  --exec-node <path>   the node (>= 20.11) that runs srt and the probe
                       client (default: node on PATH, resolved once)
  --exec-max-output-bytes N
                       kill a workspace-exec command whose stdout or stderr
                       passes N bytes (default 67108864, 64 MiB)
  --proof-key-file <file>
                       a file holding exactly 32 bytes: the MAC key the
                       workbench grade validator authenticates stored proof
                       envelopes under. Needs --ailang-bin and at least one
                       --proof-required-identity; read once at startup and
                       never logged. Unset, every object page reads "no
                       canonical host projection" and reads no evidence
  --proof-required-identity NAME
                       an identity every validated proof must verify
                       (repeatable; needs --proof-key-file)

Exit codes: 0 ok, 1 usage or client error, 2 fatal startup,
            3 integrity refusal (why: a broken link; call --strict: ok:false;
              setup: a digest mismatch, nothing installed).
            outcome record: 1 refused before sending, 2 commit refused,
              3 commit outcome uncertain (reconcile, never resend blindly).
```

`ailang-worldd help` prints the same text. Every verb prints its help with `<verb> --help` or
`ailang-worldd help <verb>` (exit 0); the verbs without a help text of their own print their
lines of the usage above.

`--addr` is refused with the `session` verbs, `setup`, and `serve`, because they act on files
directly, not on a running daemon.

Every invocation refuses to start (exit 2) if `AILANG_REGISTRY_API_KEY` is set in its
environment, except a help request (`--help`, `-h` or `help <verb>`), which is answered from the
arguments alone. The refusal's last line is the `env -u AILANG_REGISTRY_API_KEY …` command to
run instead.

### `serve`

See the flags in the usage text above, and [Operating the daemon](../guides/operating-the-daemon.md).

### Client verbs

All client verbs send one bounded request (30 s client timeout) to `--addr` and print the
response body on success. On an HTTP error they print
`ailang-worldd: <path> returned HTTP <status> <Class>: <message>` to stderr (with
`(observedHead=… selectedHead=…)` for a conflict) and exit 1.

| Verb | Request |
|---|---|
| `health` | `GET /v1/health` |
| `head` | `GET /v1/head` (prints the ref as plain text) |
| `world get <ref>` | `GET /v1/worlds/{ref}` |
| `object get <ref> [--payload]` | `GET /v1/objects/{ref}`, with `?payload=true` when `--payload` is given |
| `object find <semanticId> [--after <ref>] [--limit N]` | `GET /v1/objects/by-semantic-id/{semanticId}` |
| `log get <index>` | `GET /v1/log/{index}` |
| `log range --from N [--limit M]` | `GET /v1/log?from=N&limit=M` |
| `registry get <name>` | `GET /v1/registry/{name}` |
| `commit --file <commit.json> [--session <file\|token>]` | `POST /v1/commit` |
| `outcome record --file <outcome.json\|-> [--session <file\|token>]` | `POST /v1/commit` of one `world/task-outcome/v1` object (see [below](#outcome-record)) |

Flags come **after** the positional argument for `object get`, `object find`, and `log range`.

#### `commit`

```text
Usage of ailang-worldd commit:
  -file string
    	commit JSON file
  -session string
    	session credential for the session-gated /v1/commit: a file holding the 64-hex token (preferred) or the token itself; default $WORLD_SESSION
```

The session is resolved like `call`'s (see [Session credential](#session-credential)) and sent as
`Authorization: Bearer <token>`. With neither `--session` nor `WORLD_SESSION` the daemon answers
`401 SessionAbsent`.

#### `outcome record`

```text
usage: ailang-worldd [--addr <url>] outcome record --file <outcome.json|-> [--session <file|token>]
```

Row 182. An operator harness records how one graded task turned out, and how it was judged, as
one `world/task-outcome/v1` object committed next to the episode's effects (written by
`harness:floor`), so a provenance walk can reach the outcome, the grading rule and the task's
entry range. It adds no route and no authority: the object is committed under the episode's own
session credential. `--file -` reads stdin.

The payload (at most 16 KiB) is strict JSON with exactly these keys: `agent`, `arm` (`shell` or
`world`), `entries` (`{from, to}` for a `world` row, `null` for a `shell` row), `episodeId`,
`grading` (`argv`, `expectedStdout`, `graderVersion`, `inputFiles`, `rule`, `toolSha256`),
`label`, `outcome` (`class`, `compileOk`, `errorCategory`, `runtimeOk`, `solutionSha256`,
`solutionState`, `stderrHead`, `stdoutHead`, `stdoutOk`), `run` (≥ 1) and `task`. Unknown keys
are refused. On success it prints
`recorded world/task-outcome/v1 <agent>/<arm>/<task> run <n>: object <ref> entry <N>`.

Exit codes: 0 recorded; 1 refused before sending (usage, an unreadable file or an invalid
payload); 2 the daemon refused the commit; 3 the commit outcome is uncertain: reconcile with
`object get <ref>` before any resend.

### Session credential

`tools`, `call`, `commit` and `outcome record` resolve the session the same way:

1. `--session <value>`, else the `WORLD_SESSION` environment variable, else an error.
2. A 64-hex value is the token itself. Given on the command line it works but warns, because
   other local processes can read a process's arguments.
3. Any other value is a file (at most 256 bytes) holding the token, which is how
   `session new --out <file>` and `session mint --out <file>` write it. A file readable by group or other gets a warning
   (`chmod 600` it).

The token is only ever sent as the `Authorization: Bearer` header; no command prints it.

### `tools list`

```text
usage: ailang-worldd [--addr <url>] tools list [--session <file|token>] [--json]

Lists the tools the session may call (MCP tools/list on /mcp/): name,
description and required arguments. --json prints the tools array verbatim.

The session is --session (a file holding the 64-hex token, preferred; or the
token itself, which warns) or $WORLD_SESSION.
```

```bash
ailang-worldd tools list --session ~/.ailang/world/ep1.session
```

```text
ailang-check       Type-check and Z3-verify an AILANG file inside this episode's worktree (policy-tool `ai_check`) and …
                   required: path
ailang-read        Read a file inside this episode's worktree (policy-tool `read`). The path is relative to the worktre…
                   required: path
…
9 tool(s)
```

The list follows the session's grants: a session granted only `Workspace.Read` lists only
`ailang-read`.

### `call`

```text
usage: ailang-worldd [--addr <url>] call <tool> [--session <file|token>]
           [--arg k=v]... [--arg-json k=<json>]... | --json '<obj>'|@file|-
           [--json-out] [--strict]

Calls one tool (MCP tools/call on /mcp/). Arguments come from repeated --arg
(a string value) and --arg-json (any JSON value), or from one --json object
(inline, @file, or - for stdin); the two forms do not mix.

Output: the result's fields (long values elided with their byte counts), then
the world block: the plan ref and, per effect, its id, status and record ref.
--json-out prints the committed output bytes exactly (one trailing newline),
so sha256 of stdout minus that newline is the output ref, and the output can
be piped to 'ailang-worldd why -'.

Exit: 0 committed (including a committed refusal, ok:false); 3 with --strict
when the committed result has ok:false; 1 on a tool error, a session denial
or a host failure. A host failure is probed: "no world head" means commit a
genesis first; "a commit landed (entry N)" means do not retry — run why N.
```

```bash
ailang-worldd call ailang-read --session ~/.ailang/world/ep1.session --arg path=data.txt
```

```text
content:     hello from ep1\n
ok:          true
policy_digest: 4b7541a98c773068896571d54c5182bc5601b335f2c5e663743c07eda077b121
tool:        sha256:e55ff71c710c10395ebf949f6777ef9114bed8bba859e3ba719ed73c6d1b9c7f
output:      sha256:a70447b8e8c490f547b021f257ee0c3799857d3d70ec6eaeccfb9a6293ee1572 (416 bytes)
world:
  plan       sha256:a1f0feb6b661e46646dab346043e8c144bf3d2be7733df0e520506119b2b773f
  effect     e1 ok sha256:4770647ef5f9b08c8ccb51773631a3460a070d2be9a06cffb609b4eb094df1df
```

`/mcp/` answers in three shapes, and `call` tells them apart by HTTP status and `Content-Type`:

| Answer | Meaning | What `call` does |
|---|---|---|
| `200 text/event-stream` | the JSON-RPC response | prints the result; a JSON-RPC `error` inside it is a tool error (exit 1) |
| `200 application/json` with `-32603` | a host failure, which can follow a commit that did land | reads `/v1/head` before and after the call: no head → "no world head" (commit a genesis); the head moved → "a commit landed (entry N) — do not retry"; otherwise the message as sent |
| `401 text/plain` | the session was refused | prints the reason and a one-line fix (absent, unknown, malformed or expired) |

### `why`

```text
usage: ailang-worldd [--addr <url>] why <target> [--scan N] [--json]
       ailang-worldd [--addr <url>] why --result <file> [--scan N] [--json]

Walks one committed invocation back to its whole provenance chain, checking
every link (✓ or ✗):

  world    the world ref, recomputed from the entry and the output, and
           confirmed by GET /v1/worlds/<ref>
  entry    index, writtenBy, prev, transitionFn and the interpreter (with its
           pin name); the entry hash is recomputed
  record   invocation id, episode, skill (its hash and its fields agree with
           the entry)
  input    your arguments, decoded and elided; when elided, the short
           top-level fields (path, command, op, …) are printed first
  plan     the effects planned, or the refusal result
  effect   per effect record: allowed/failed/denial, budget before→after,
           request and result sizes; a toolchain call that changed the
           worktree adds "worktree: N created, M modified, K deleted" and
           its first paths (--json: the whole worktree_delta)
  output   the committed output; it must equal world.stateRoot
```

(The help goes on with the targets and the scan bound, as the table and paragraph below say.)

Two lines help a walk read a call without fetching payloads. When the input is long, its short
top-level fields come first, for example
`old_text="}\n" path="host/worktree/worktree_test.go" · {"new_text":…`, so the file an
`ailang-edit` touched is visible (PR #292). When an `ailang-run`, `ailang-check`, `ailang-cli` or
`workspace-exec` result carries a `worktree_delta` (row 184), the effect gets a line such as
`worktree: 6 created, 1 modified, 1 deleted (basis git-status): c .ailang/cache/f0, …`; a
truncated delta says `first N of M changes listed`, and a failed snapshot prints
`worktree: not recorded — <reason>`.

| Target | Resolves by |
|---|---|
| `<index>` or `head` | the log entry directly |
| `sha256:<ref>` | the object's `semanticId`: a record, input, output, plan, effect record, effect request or result; or a world ref |
| `mcp:<id>` or `a2a:<id>` | the record's `invocationId` |
| `rest:<id>` | its receipt's world |
| `-` or `--result <file>` | the output of `call --json-out`: its sha256 is the output ref (the result's `world.plan` is the fallback) |

Every target except an index or `head` is found by scanning the log backwards from the head,
one object read per entry, at most `--scan N` entries (default 500, maximum 5000) within 60
seconds. Past either limit it prints `not found in the last N entries` and exits 1. An entry
not written by the coordinator (a REST genesis, say) is shown with its object only, with the
reason.

Exit codes: 0 when every link is ✓; **3 when any link is ✗** (the bytes the daemon served do
not match a content address); 1 when the target is not found. `--json` prints the chain as
JSON. See [Provenance](../agents/provenance.md) for a walked example.

### `log tail`

```text
usage: ailang-worldd [--addr <url>] log tail [--from N] [--follow] [--interval 1s] [--raw]

Prints log entries, one line each: index, short entry hash, writtenBy, and for
a coordinator invocation its episode, skill and effect statuses.

  --from N      first entry (default: head-19, the last 20 entries)
  --follow      keep polling for new entries until Ctrl-C; a daemon restart
                is retried with backoff (up to 5 s) and the outage is named
  --interval D  poll interval with --follow (default 1s)
  --raw         print each entry as its JSON
```

```text
#1 ef2b3322 coordinator:mcp ep1 ailang-read [Workspace.Read ok]
#2 b7aa36d6 coordinator:mcp ep1 ailang-read [Workspace.Read ok]
#3 c0cd3ae0 coordinator:mcp ep1 ailang-read [no effects]
```

With `--follow` the next read starts at the entry after the last one printed, so each entry
is printed exactly once.

### `provenance`

```text
usage: ailang-worldd [--addr <url>] provenance [--since <entry>] [--episode <ep>] [--scan N]
```

Prints the trailer that a pull request or commit made through World carries:

```text
World-Provenance: store=/Users/you/.ailang/world/world.db episode=ep1 entries=1-4
```

`store` is the daemon's database path from `/v1/health`; `entries` are the episode's first and
last coordinator entries from `--since` (default: the last 500 entries) to the head. Without
`--episode` it prints one trailer per episode in the range. Any entry in the range can be walked
back with `ailang-worldd why <entry>`. It exits 1 when the range holds no coordinator entry for
the episode.

### `setup`

```text
usage: ailang-worldd setup [--interpreter-dir <dir>] [--tools-dir <dir>]
           [--db <path>] [--workspace-root <dir>] [--from-dir <dir>] [--replace]

Installs the two pinned AILANG binaries World runs, verified byte for byte:
  interpreter  v0.41.0  -> <interpreter-dir>/ailang  (serve --ailang-bin)
  tool         <tool>   -> <tools-dir>/<tool>/ailang  (serve --tool-ailang-bin)
where <tool> is the tool-binary release this daemon is built for.

For each pin: a file already hashing to the pin is left alone (no network
request). Otherwise the release's .sha256 and tarball are fetched over https
from github.com/sunholo-data/ailang (the URL is compiled in), and the install
happens only when the release .sha256, the computed tarball sha256 and the
compiled-in digest all agree and the extracted ailang hashes to the
compiled-in binary digest. No tarball is kept; pin.json records what was
installed. A mismatching existing file is refused unless --replace, which
keeps it as ailang.prev-<sha8>. setup never runs a downloaded byte.

Then it creates the store directory and the workspace root (mode 0700),
refusing a workspace root that contains the store, and prints the attended
steps that follow.

  --interpreter-dir <dir>  default ~/.pinned-ailang
  --tools-dir <dir>        default ~/.pinned-ailang-tools
  --db <path>              world store (default ~/.ailang/world/world.db);
                           only its directory is created, never the store
  --workspace-root <dir>   default ~/.ailang/world-ws
  --from-dir <dir>         install offline from <dir>/<release>/<asset> and
                           its .sha256, verified the same way
  --replace                replace a mismatching existing binary

Platforms: darwin/arm64 and linux/amd64.
Exit: 0 ok; 1 usage, unsupported platform, network refusal or an existing
mismatching file; 3 a digest mismatch (nothing installed).
```

```bash
ailang-worldd setup
```

```text
✓ interpreter AILANG v0.41.0 installed /Users/you/.pinned-ailang/ailang (sha256 1a67b0146858…)
✓ tool        AILANG v0.54.15 installed /Users/you/.pinned-ailang-tools/v0.54.15/ailang (sha256 d821dad0d6e9…)
✓ store directory /Users/you/.ailang/world
✓ workspace root  /Users/you/.ailang/world-ws

next (the attended steps; docs/QUICKSTART.md §9 has the full walk):
…
```

A second run prints `present` for both pins and makes no network request. The compiled-in
digests come from the release manifest `pins/ailang-releases.json`, embedded through
`host/pins` (`cmd/ailang-worldd/pins.go` derives the pin table from it since row 139 M1):

| Release | Platform | Tarball sha256 | Binary sha256 |
|---|---|---|---|
| v0.41.0 (interpreter) | darwin/arm64 | `b08f3cde…598e0b` | `1a67b014…5b9f` |
| v0.41.0 (interpreter) | linux/amd64 | `fa0045de…faa56` | `8e7a275d…25fb5` |
| v0.54.15 (tool) | darwin/arm64 | `710945da…281382` | `d821dad0…ac043` |
| v0.54.15 (tool) | linux/amd64 | `87731698…c4945` | `db0f9dac…ddd63c` |

v0.52.1, the tool before row 139, stays in the manifest as a `retired` row with its digests;
`setup` and `doctor` install and check only the current pins.

A tarball whose own `.sha256` agrees with it but not with the compiled-in pin is refused with
exit 3: `tarball digest mismatch: release .sha256 …, computed …, compiled-in pin …`, and nothing
is installed. Release signatures (`.sig`/`.pem`) are not checked.

### `doctor`

```text
usage: ailang-worldd [--addr <url>] doctor [--db <path>] [--workspace-root <dir>]
           [--interpreter-dir <dir>] [--tools-dir <dir>] [--examples-dir <dir>] [--online]

Checks this machine for everything World needs, read-only, and prints one
line per check: ✓ fine, ! worth knowing (with a fix), ✗ broken (with a fix).

  api key      AILANG_REGISTRY_API_KEY is unset (doctor could not run otherwise:
               every verb refuses while it is set)
  tty          /dev/tty opens, so the attended steps (world-publish, session
               new/mint) can read the confirmation you type
  pins         both pinned binaries present and hashing to the compiled-in
               digests (never executed); stale tarballs and old binaries noted
  daemon       what answers at --addr: a worldd (its store and interpreter),
               a foreign listener, or nothing
  store        exists (never created), writer lock held or free, world head
               and transition registry present
  examples     the examples corpus examples-search reads
  workspace    the workspace root; each child a git worktree with a valid
               episode name
  --online     the release .sha256 of both pins is reachable and matches

Defaults: --db ~/.ailang/world/world.db, --workspace-root ~/.ailang/world-ws,
the pin directories of 'setup', --examples-dir ~/.ailang/examples.
Exit: 0 no ✗; 1 at least one ✗ (or a usage error).
```

```bash
ailang-worldd doctor --online
```

```text
✓ api key    AILANG_REGISTRY_API_KEY is unset
! tty        no controlling terminal (open /dev/tty: device not configured): world-publish and session new/mint will refuse here
             fix: run the attended steps yourself from a terminal window (an IDE terminal pane works); an agent cannot run them
✓ pins       interpreter AILANG v0.41.0 at /Users/you/.pinned-ailang/ailang (sha256 1a67b0146858…)
✓ pins       tool AILANG v0.54.15 at /Users/you/.pinned-ailang-tools/v0.54.15/ailang (sha256 d821dad0d6e9…)
! daemon     nothing answers at http://127.0.0.1:7644
             fix: ailang-worldd serve --db … (when you need it running)
! store      /Users/you/.ailang/world/world.db does not exist yet (doctor never creates it)
             fix: world-publish transitions --store /Users/you/.ailang/world/world.db … creates it (the attended publish)
! examples   no examples corpus (~/.ailang/examples is absent): examples-search will refuse
             fix: pass serve --examples-dir <dir>
✓ workspace  /Users/you/.ailang/world-ws: 0 episode worktree(s)
✓ online     interpreter AILANG v0.41.0: release .sha256 reachable and matches the pin
✓ online     tool AILANG v0.54.15: release .sha256 reachable and matches the pin
no check failed
```

doctor creates nothing: the writer lock is probed with a shared, non-blocking lock on the lock
file opened read-only, and the store is read through a read-only handle, so it runs safely
beside a live daemon. With `AILANG_REGISTRY_API_KEY` set, doctor (like every verb) refuses with
exit 2 and names the variable — that refusal is the finding.

### `session new`

```text
usage: ailang-worldd session new <episode> [--db <path>] [--workspace-root <dir>]
           [--repo <dir>] [--preset se-tools] [--grant EFFECT=SCOPE:BUDGET]...
           [--budget 50] [--ttl 3600] [--out <file>] [--branch <name>]
           [--allow-non-git]

ATTENDED. Provisions one episode: its git worktree at <workspace-root>/<episode>
and a session credential bound to it, after a y/N typed at /dev/tty (any
terminal window, an IDE terminal pane included; stdin is never read). With no
controlling terminal (an agent, cron, CI) it refuses before touching anything.

  <episode>            ^[a-z0-9][a-z0-9-]{0,63}$ (the daemon's grammar)
  --db <path>          an EXISTING world store (default ~/.ailang/world/world.db);
                       never created. The daemon must be stopped.
  --workspace-root <dir>
                       default ~/.ailang/world-ws; must not contain the store
  --repo <dir>         the git repository the worktree is added from
                       (git worktree add --detach); not needed when
                       <workspace-root>/<episode> is already a worktree
  --branch <name>      add the worktree on a new branch instead of detached
  --allow-non-git      adopt an existing <workspace-root>/<episode> that is
                       not a git worktree (refused without it). The episode
                       is DEGRADED: no branch diff, provenance trailer or
                       reset to a base, and its worktree deltas use the
                       bounded walk (basis walk / walk-capped), not
                       git-status. session list and /v1/health show it.
  --preset se-tools    one grant per effect the se-tools transitions declare:
                       Workspace.Read, Workspace.Write, Ailang.Check, Ailang.Run,
                       Ailang.RunEnv, Ailang.RunNet, Ailang.Discover, Ailang.CLI,
                       Workspace.Exec, each scope worktree with --budget calls
  --grant EFFECT=SCOPE:BUDGET
                       an explicit grant (repeatable); it overrides a preset
                       grant with the same effect and scope
  --budget N           budget of each preset grant (default 50)
  --ttl N              session lifetime in seconds (default 3600); every grant
                       expires with the session
  --out <file>         write the token once to <file> (mode 0600); default:
                       print it to stdout once

The credential_id (the hash 'session revoke' takes) goes to stderr; the token
never does. If the mint fails, the worktree this run created is removed.
```

```bash
ailang-worldd session new ep1 --repo ~/src/my-project --preset se-tools --out ~/.ailang/world/ep1.session
```

On the terminal it asks:

```text
Create worktree <workspace-root>/ep1 from <repo> and mint a session for episode ep1 with 9 grant(s) [Workspace.Read=worktree:50 Workspace.Write=worktree:50 Ailang.Check=worktree:50 Ailang.Run=worktree:50 Ailang.RunEnv=worktree:50 Ailang.RunNet=worktree:50 Ailang.Discover=worktree:50 Ailang.CLI=worktree:50 Workspace.Exec=worktree:50], expiry +3600s? [y/N]
```

and, after `y`:

```text
session for episode ep1: 9 grant(s), expires epoch 1791060000, token written to <session-file>
ailang-worldd session new: worktree <workspace-root>/ep1 created
ailang-worldd session new: credential_id=<64-hex> (episode ep1, worktree <workspace-root>/ep1, 9 grant(s), expires 1791060000)
```

With `--allow-non-git` and an existing non-git directory, the prompt instead begins
`Use the existing NON-GIT directory <workspace-root>/<ep> as a DEGRADED episode (no branch diff,
provenance trailer or reset to a base; walk-basis worktree deltas)`, and success adds
`ailang-worldd session new: episode <ep> is DEGRADED: <dir> is not a git worktree
(--allow-non-git); session list and /v1/health show it` on stderr. Without the flag, such a
directory is refused (exit 1) with a message naming `--allow-non-git`.

Without a controlling terminal it refuses before touching anything:
`refusing: no controlling terminal (open /dev/tty: device not configured); provisioning a session
credential requires one human act at a terminal. fix: run it yourself from a terminal window (an
IDE terminal pane works); an agent cannot run this step` (exit 1). Every validation (the episode name,
an existing `--db`, the workspace root, the grants, `--out` not existing yet) runs before that
fence; the writer lock is taken before the worktree is created, so with the daemon running it
stops at `stop the daemon first` with nothing created.

### `session list`

```text
usage: ailang-worldd session list [--db <path>] [--workspace-root <dir>] [--episode <ep>] [--json]

Lists the session credentials in a store (at most 500): credential_id (the
hash, never the token), episode, live or expired, expiry and grants. Reads
through a read-only handle, so it works while the daemon is running. A --db
that does not exist is refused, never created.

An episode whose <workspace-root>/<episode> exists but is not a git worktree
(adopted with session new --allow-non-git) is marked DEGRADED ("degraded":
true with --json). --workspace-root defaults to ~/.ailang/world-ws.
```

```text
3f9c…e21a  ep1          live    expires 2026-10-03T21:40:00Z  Ailang.CLI=worktree:50 Ailang.Check=worktree:50 …
7b10…04cd  plain        live    expires 2026-10-03T21:40:00Z  Workspace.Read=worktree:3  DEGRADED (not a git worktree)
2 session credential(s)
```

`--json` prints the same rows as JSON (`credentialId`, `episode`, `live`, `createdAt`,
`expiresAt`, `grants`, `degraded`).

### `session revoke`

```text
usage: ailang-worldd session revoke [--db <path>] <credential_id>
       ailang-worldd session revoke [--db <path>] --episode <ep>

Deletes a session credential's row, so the next request with its token is an
unknown credential. <credential_id> is the 64-hex hash session mint/new print
to stderr (and session list shows) — never the token itself. --episode
revokes every credential of that episode. Flags and the id may come in either
order. An id (or episode) with no credential in the store is refused (exit 1),
and so is a --db that does not exist (it is never created). The daemon must be
stopped: revoke needs the writer lock.
```

```bash
ailang-worldd session revoke --db ~/.ailang/world/world.db <credential_id>
ailang-worldd session revoke <credential_id> --db ~/.ailang/world/world.db   # either order
ailang-worldd session revoke --db ~/.ailang/world/world.db --episode ep1     # every credential of ep1
```

An id with no credential in the store prints
`refusing: no session credential <id> in <db> (already revoked, or minted in another store);
nothing was revoked` and exits 1. With the daemon running it stops at `stop the daemon first`.

### `session mint`

```text
Usage of ailang-worldd session mint:
  -db string
    	world store database (required)
  -episode string
    	episode the session binds to (required)
  -grant value
    	EFFECT=SCOPE:BUDGET, repeatable, at least one
  -out string
    	write the credential ONCE to <path> at mode 0600 (default: print to stdout exactly once)
  -ttl int
    	lifetime in whole seconds (default 3600)
```

The lower-level form of `session new`: it mints a session for an episode whose worktree you made
yourself, with explicit `--grant` flags. Attended: opens `/dev/tty` and asks
`Confirm mint for episode … [y/N]` there; refuses without a controlling terminal. The store must
exist (it is never created) and the daemon must be stopped (single writer). Each grant's expiry
is set to the session's (`now + ttl`). With `--out` it prints
`minted session credential for episode <ep>: <n> grant(s), expires epoch <t>, written to <path>`;
it always prints `credential_id=<hash>` to stderr. See
[Sessions and episodes](../concepts/sessions-and-episodes.md).

### `demo seed`

```text
usage: ailang-worldd [--addr URL] demo seed --db PATH --session FILE
     --ailang-bin PIN --proof-key-file KEY --manifest FILE [--addr URL]

Prepare a disposable demo worktree and publish descriptors, then mint the demo
session at a controlling terminal (QUICKSTART section 12). Seed never mints,
publishes, creates a worktree or opens another writable database.
--db identifies the daemon's store and permits read-only grant inspection.
--session is a private session file (default WORLD_SESSION).
--ailang-bin is the capsule/proof pin (v0.41.0); serve's separate
--tool-ailang-bin must be v0.54.15 for the workspace tools.
--proof-key-file is the private operator key required by later proof stages.
--manifest is a new output path; existing files are refused.
--addr is the daemon URL, also accepted as the global client flag.

This incremental build implements preflight and S1 domain commits only.
S2-S6 incomplete: the command refuses completion and writes no success manifest.
Do not retry a partial seed; inspect its recorded history first.
```

Release-demo tooling (queue row 150), still incremental. `--db`, `--ailang-bin`,
`--proof-key-file` and `--manifest` are required. The preparation steps are
`docs/QUICKSTART.md` section 12, "Prepare the release demo seed".

## `world-publish`

```text
world-publish — the attended entrypoint for an IRREVERSIBLE public write

  world-publish packet     [--package-dir D] [--golden G]
  world-publish approve    --store S [--registry-origin O] --now N --expires E
  world-publish publish    --store S --registry-origin O --publisher P \
                           --credential-file C --approval-ref R --now N --expires E (--live | --dry-run)
  world-publish reconcile  --store S [--registry-origin O] [--probe]
  world-publish transitions --store S --manifest M (--ailang-bin B | --interpreter-ref R)  [local registry write]
  world-publish epoch nominate --store S --release vX.Y.Z (--candidate | --new-epoch) --gate-verdict V
                           --binary B  [local epoch-registry write; --new-epoch is ATTENDED]

Exit codes: 0 done · 1 failed · 2 usage · 3 STOP (a fence refused; nothing happened)
```

This usage is printed when the program is run with no verb or an unknown one
(`world-publish --help` prints `unknown verb "--help"` followed by it).

### Verbs

| Verb | Kind | What it does |
|---|---|---|
| `packet` | Read-only, headless-permitted | Recomputes the `world/core` ready-packet from `--package-dir` and compares every field with `--golden`. `STOP fence=packet reason=drift` on any difference. |
| `approve` | **Attended** | Mints a one-shot `ApprovalDecisionV1` for the current packet and prints its ref. |
| `publish` | **Attended**; irreversible with `--live` | Spends the approval, exactly once. Exactly one of `--live` or `--dry-run`. `--dry-run` runs every fence and sends nothing. |
| `reconcile` | Read-only, headless-permitted | Lists publish intents with no outcome; `--probe` issues read-only metadata GETs to resolve them. `--live` is refused. |
| `transitions` | **Attended**; local registry write | Publishes a descriptor manifest into the store's transition registry. `--live` is refused. See [Transitions manifest](transitions-manifest.md). |
| `epoch nominate` | `--candidate` headless-permitted; `--new-epoch` **attended**; local epoch-registry write | Appends one immutable `world/epoch-registry/v1` revision naming `--release`, either as a candidate of the current epoch (`--candidate`: replay goldens byte-identical) or as a new epoch (`--new-epoch`: semantics moved; typed phrase `open epoch <N> for <release> in <store>`). Exactly one of the two. It refuses unless the release is a row of the compiled-in manifest `pins/ailang-releases.json` (there is deliberately no flag naming another manifest), `--gate-verdict` is a passing `world/ailang-release-verdict/v1` verdict for that release on this platform, and `--binary` hashes to the manifest's and the verdict's binary digest and prints the release's `--version` line. The store must already hold an epoch registry (start the daemon on it once), and the daemon must be stopped. `--live` is refused. |

### Flags

All verbs share one flag set. Each verb uses the subset shown in the usage above. `-h` on any
verb prints this list (and exits 2):

```text
  -ailang-bin string
    	transitions: interpreter binary to archive now, as the daemon does at startup
  -approval-ref world-publish approve
    	the ApprovalDecisionV1 hashref minted by world-publish approve
  -binary string
    	epoch nominate: the release's ailang binary, hashed against the compiled-in release manifest
  -candidate
    	epoch nominate: add the release to the CURRENT epoch (goldens byte-identical; headless-permitted)
  -credential-file string
    	file outside the working tree holding the registry API key
  -decided-by string
    	who is granting the approval
  -dry-run
    	rehearse the publish: every fence, no request
  -episode string
    	durable episode ID (default "attended-publish")
  -expires int
    	logical time the approval expires at
  -gate-verdict string
    	epoch nominate: the release's compatibility-gate verdict (world/ailang-release-verdict/v1)
  -golden string
    	the committed ready-packet golden (default "scripts/world_package_ready_packet.golden.json")
  -interpreter-ref string
    	transitions: archived interpreter HashRef pinned into every descriptor
  -live
    	PERFORM THE IRREVERSIBLE PUBLIC WRITE
  -manifest string
    	transitions: descriptor manifest file (JSON array)
  -new-epoch
    	epoch nominate: open epoch N+1 for the release (ATTENDED: typed phrase on /dev/tty)
  -now int
    	logical time of this act
  -package-dir string
    	the projected package directory (default "packages/world-core")
  -probe
    	reconcile: issue the read-only metadata GETs
  -publisher string
    	path to the pinned released ailang binary
  -registry-origin string
    	the read-only public bucket origin
  -release string
    	epoch nominate: the AILANG release to nominate (vX.Y.Z)
  -requester string
    	who is requesting the approval
  -store string
    	path to the world database this publish is recorded in
```

There is deliberately no flag that sets the registry validator origin; it is a compiled
constant. A test compares the flag set to a frozen list, so adding a flag fails the build.

### STOP lines

A refusal prints `STOP fence=<name>[ reason=<reason>]` and exits 3. Fence names: `mode`, `ci`,
`tty`, `confirmation`, `store`, `approval`, `credential`, `packet`, `handler`, and for
`epoch nominate` also `release` (absent or not in the manifest), `verdict` (absent, unreadable,
malformed or not admitted), `binary` (absent, unreadable, digest or version mismatch) and
`registry` (no epoch registry yet, or the release already nominated). See
[Attended steps](../guides/attended-steps.md#how-the-fences-work).

Build `world-publish` to a binary for attended use rather than using `go run`: `go run` exits 1
for a child that exited 3, which hides the STOP contract.

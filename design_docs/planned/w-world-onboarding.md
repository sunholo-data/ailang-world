# w-world-onboarding — start a World episode on any repo in one attended act, and let agents and operators learn World from a plugin

**Status**: DRAFT r1 (attended design, 2026-10-19). Five decisions are OPEN for Mark (§6). No milestone may start before they are ruled.
- **RULED (Mark, attended 2026-10-19, D-WORLD-85):** C, A, C, A, A (D-189-1…5) — every recommendation accepted ("thank you yes go with all recommendations"); ledger PR #317.

**Direction:**
- Queue row 189, filed attended 2026-10-16. Mark asked: "how do I now start to use this on other projects? do we have an ailang world skill that needs to be globally available so I can ask the agents in there?"
- The row 189 scope note (attended 2026-10-19, PR #311): the skill also covers the operator side, watching an episode while it runs.
- Taken attended on 2026-10-19 ("I can run some attended sessions now to help move it forward"), ahead of the D-WORLD-81 post-tag order. The loop keeps 156 → 150 → 151.

**Base:** ailang-world `dev` `78f2860`.

**Rig:**
- macOS arm64; Claude Code 2.1.303; codex at `/opt/homebrew/bin/codex`.
- The trial rig `~/.ailang/state/trial-2026-10-16` with its daemon on 127.0.0.1:47330. This doc only read it (health, files). Nothing was stopped or minted.

**Clauses**: 5 (human surface: adoption) and 3 (authority: the mint stays one human act).

**Estimate**: ≈ 2 executor days across M0–M7 (§8), each milestone ≤ ~30 min of executor work and ≤ ~250 LOC.

**Probe conditions**:
- Every V-row is a command run or a file:line read on 2026-10-19 at the base above. Byte evidence that is not a repo file is banked in `design_docs/verification/world-row189-design-2026-10-19/`.
- No session was minted and no daemon was stopped.

## 1. Why (measured)

- **P1 — Starting one episode takes about eleven manual operator steps (V1, V2, V3).** In order:
  1. write a task file and a trial spec;
  2. refresh a prep clone;
  3. stop the daemon;
  4. confirm it stopped;
  5. rebuild or copy the daemon binary;
  6. mint with nine explicit `--grant` flags at `/dev/tty`;
  7. restart `serve` with ~14 flags, including one `--exec-episode-project` for every episode the store has ever served;
  8. run `scripts/trial/trial.py run` with a hand-built environment;
  9. read the graders;
  10. commit from the operator checkout;
  11. open the PR.

  Steps 3–7 are friction #1, #2 and #13. Only the operator who built the rig can do this today.
- **P2 — Minting requires a stopped daemon, structurally (V4, V5).**
  - `session new` takes the store's exclusive, non-waiting `flock` writer lock.
  - `serve` holds that same lock.
  - So a running daemon makes every mint fail, with "stop the daemon first".
  - Stopping the daemon stops every other episode in flight (friction #2).
- **P3 — A new episode on a multi-project daemon needs a restart even with no mint (V7).**
  - With more than one exec profile loaded, `profileFor` returns nil for an episode that has no `--exec-episode-project` flag.
  - That mapping is fixed at startup, so `workspace-exec` stays unconfigured for the new episode until `serve` restarts with one more flag.
- **P4 — A stale daemon build silently serves an old pin (V8, friction #13).** The rig binary from before row 139 kept serving the v0.52.1 tool binary after #308 merged. Nothing in the procedure compares the rig build with `dev`.
- **P5 — No agent in another repo can discover World (V9, V10).**
  - No `ailang-world` plugin or skill exists on this machine.
  - The marketplaces installed here are `ailang-marketplace` (the `use-ailang` family), `ailang-parse-marketplace` and two Anthropic ones.
  - The only way an agent reaches World today is a launch made by `trial.py`, with a private MCP config and a strict tool set.
- **P6 — An exec profile is hand-written JSON in a grammar only the broker documents (V11, V12).**
  - The five profiles on this machine run from 36 to 159 lines.
  - `platform-source.json` needs absolute interpreter roots for `uv`, `sh -c` argv wrappers, and per-command flag classes (`regex` with `form: sep`, `bool`).
  - `ailang-world.json` needs a seeded `GOCACHE`.
  - The project's own `pyproject.toml`/`go.mod` already says most of what a first draft needs.

## 2. What already exists and is reused

- **The mint (V4):** `session new` already does validate → TTY fence → y/N naming the worktree and every grant → writer lock → `git worktree add` (with `--branch`) → `authority.Mint` → rollback of the worktree it created. `episode start` reuses this path, never a copy.
- **`trial.py` (V3):** spawn, private 0600 MCP config, strict tool set, log range, graders, outcome records and PR body. `episode start` produces the launch; grading stays in `trial.py` until a later `episode finish`, which this row does not include.
- **Other verbs:**
  - `setup` fetches and verifies the pins and makes the store directory and the workspace root. It "never runs a human act (publish, mint)" (`cmd/ailang-worldd/setup.go:1-6`).
  - `doctor` checks pins and paths (`cmd/ailang-worldd/doctor.go:59`).
  - `episode start` calls neither, but tells the operator to run them when the store or pins are missing.
- **Watching an episode (row 198):** `provenance --episode`, `why`, and later `episode report`. The skill teaches these. Row 198 builds the report, not this row.

## 3. The three items

### 3.1 `ailang-worldd episode start`

```
ailang-worldd episode start <episode> --repo <dir> [--branch <name>] [--profile <project>]
    [--task <file>] [--agent claude|codex] [--model <m>] [--preset se-tools] [--ttl N]
    [--db …] [--workspace-root …] [--addr 127.0.0.1:47330]
```

ATTENDED, like `session new`. Its steps:
1. Validate everything first: episode grammar, repo, profile, task file. Then compare the daemon's build and the manifest pins with this binary's, and refuse a mismatch, naming the rebuild (P4).
2. Run `session new`'s path unchanged: TTY fence, y/N, worktree on `--branch` (default `world/<episode>`), mint. The y/N line also names the exec profile and its digest (D-189-2).
3. Write the token to `~/.ailang/world/sessions/<episode>.session` (0600, written once).
4. Reuse a serving daemon on `--addr`, or start one in the background with the standard flags. A restart is never needed (D-189-1, D-189-2).
5. Print the launch line: the `trial.py run …` invocation, or the bare `claude`/`codex` argv with a private 0600 MCP config, the same shape `trial.py` builds (D-189-4). It never execs the agent itself.

It never commits, pushes or opens a PR.

### 3.2 `ailang-worldd profile init`

```
ailang-worldd profile init --repo <dir> --project <name> [--out ~/.ailang-world-exec/profiles/<name>.proposed.json]
```

It reads `pyproject.toml`/`uv.lock`, `go.mod`, `package.json` and `Makefile`, and proposes:
- `path` and `read_roots` (resolved interpreter or toolchain roots);
- `probe`;
- the commands the measured profiles use: `pytest-file`, `ruff-file`, `go-test`, `go-vet`, `gofmt-l`, `vitest-file`, …, with the flag classes the broker accepts (`broker.ParseExecProfile`, `host/broker/exec_profile.go:130`);
- `caches` seeds.

It writes `*.proposed.json` and prints why each command is there. **It never activates a profile**; the operator renames the file (D-189-5). The proposal must pass `ParseExecProfile` before it is written, and `serve`'s startup probe stays the real check.

### 3.3 The `ailang-world` plugin and skill

**A plugin (D-189-3) holding:**
- **Skill `ailang-world` for the agent side:**
  - use World's tools, never a shell, inside an episode;
  - `workspace-exec` commands come from the profile; ask the operator for a missing one;
  - read with `offset`/`limit` (row 190);
  - an `ok:false` refusal is a fact to report, not to work around (friction #11);
  - the 9 s exec cap and what to hand back (row 145);
  - finish with a summary the operator can grade.
- **Skill `ailang-world` for the operator side:**
  - `episode start`;
  - watching a run: `provenance --episode`, `why <entry|mcp:…>`, and `episode report` when row 198 lands;
  - grading with `trial.py`;
  - publishing (`world-publish transitions`) and minting are human acts. **The skill never mints and never types a publish phrase**; it prints the command for the human.
- **Slash commands** that only print commands: `/world-start` (the `episode start` line for this repo), `/world-watch <ep>`, `/world-why <entry>`.
- **No MCP entry in v1, if D-189-4 = A.**

## 4. Minting while a daemon serves (the authority question)

**The invariant every option must keep:** a credential exists only after a human typed `y` at a controlling terminal, read from `/dev/tty`, which an agent harness, cron or CI cannot open (V5, V6). The skill can never mint.

- **A — A daemon mint endpoint.**
  - The CLI runs the TTY fence and y/N as today, then sends the confirmed request to the daemon, which writes the row under its own writer lock.
  - The daemon cannot see the TTY, so the endpoint needs its own operator credential:
    - not a `Bearer` session token: D-WORLD-26 constraint (i) makes a Bearer value a session credential only;
    - so a unix socket with a peer-uid check, or a 0600 operator key file.
  - **Effect on the rule:** the human act moves from "structurally unfakeable" to "guarded by a key or a socket". Any process running as the operator's uid, including a native-shell agent in another repo, could mint without a human.
  - This is row R13 of `w-session-authority.md:307`, left PENDING "if ever needed".
- **B — A writer-lock handoff.**
  - The CLI asks the daemon to drain in-flight commits and release the lock, mints exactly as today, and the daemon reacquires the lock.
  - **Effect on the rule:** none; the mint code path is unchanged.
  - **Costs:**
    - a write pause on every mint;
    - an unauthenticated "pause" request becomes a denial-of-service lever;
    - it changes the ratified single-writer lifecycle (w-worldd-m2 Decision 2, `host/store/writer_lock_unix.go:1-24`).
- **C — A separate credentials store.**
  - `session_credentials` moves out of `world.db` into `~/.ailang/world/credentials.db`, with its own writer lock.
  - Mint and revoke take that lock at the TTY. The daemon resolves through a read-only handle per request, as `session list` already does while the daemon runs (V5).
  - **Effect on the rule:** none.
  - **Costs:**
    - the table is referenced from `host/store/store.go`, `host/authority/resolver.go` and `host/daemon/daemon.go` (V13);
    - a schema-version bump and a one-time migration (current version 4, `host/store/store.go:44`);
    - revocation becomes a write to the second DB.
  - **M0 must measure two things:**
    - whether budget accounting reads anything from the credential row beyond grants and expiry;
    - whether any world-store transaction joins `session_credentials`.

  No `REFERENCES` or `JOIN` turned up in the grep (V13). That is not proof.

Option A is what row 188 (a desktop approval) and row 196 (delegated minting for loop executors) will eventually need. Those rows are gated on their own rulings. Choosing B or C now does not block A later.

## 5. How an agent in another repo finds World

- **Today:** the platform-source episodes ran from a prep clone at `~/.ailang/state/trial-prep/platform-source`, whose `origin` is the operator's checkout `~/dev/sunholo-data/platform-source` (V14). The agent saw only World's MCP tools, through a private config.
- **With this row:**
  - The plugin is installed once, globally.
  - In any repo, the operator says "start a World episode for X". The agent, through the skill, prints the `episode start` line, and the operator runs it at a terminal and types `y`.
  - The new agent runs from the launch line `episode start` printed, so it holds World's tools and no shell. That is what makes its record complete: the evidence-grade mode.
- **Interactive sessions:** an ordinary interactive session with native tools plus World tools would give an incomplete record. That is D-189-4 option B, and it is not v1.
- **Daemon address:** `--addr`, default 127.0.0.1:47330, printed into the launch line.
- **Token:** never in the repo, the plugin or settings. It lives only in `~/.ailang/world/sessions/<ep>.session` (0600) and in the per-launch 0600 MCP config, which is deleted after the run as `trial.py` does today (`scripts/floor/arms.py:112-118`, `run_task.py:64-66`). Codex gets it only through `WORLD_SESSION` (`arms.py:77, 145`).

**Related rows:**
- **188:** a desktop approval could later replace the TTY y/N. That needs option A plus a new ruling.
- **196:** loop executors through World need `episode start` and a delegated-minting ruling, so this row is a prerequisite.
- **198:** its `episode report` is what the operator skill teaches for watching an episode.

## 6. OPEN decisions

- **D-189-1 — Minting while a daemon serves.**
  - **A:** a daemon mint endpoint behind an operator key or socket. The human act becomes key-guarded.
  - **B:** a writer-lock handoff.
  - **C:** a separate credentials store with a read-only resolver. The mint path is unchanged.
  - **Recommendation: C.** It keeps the TTY rule byte for byte, adds no new authority surface, and needs no write pause. B is the fallback if M0 finds a join; A waits for 188/196.
- **D-189-2 — Binding an episode to its exec profile without a restart.**
  - **A:** the mint records the project, and the y/N names it with the profile digest. `serve` loads every profile in `~/.ailang-world-exec/profiles/` at startup, and `profileFor` reads the binding from the credential.
  - **B:** `serve` re-reads a mapping file on SIGHUP.
  - **C:** keep the restart, with `episode start` restarting the daemon.
  - **Recommendation: A.** Which commands an agent may run is authority, so the human should confirm it at the mint. B adds an unconfirmed side channel, and C kills episodes in flight (friction #2).
- **D-189-3 — Where the plugin lives.**
  - **A:** its own repo and marketplace (like `sunholo-data/docparse-skill`).
  - **B:** a plugin inside `sunholo-data/ailang_bootstrap`'s `ailang-marketplace`, next to `use-ailang`.
  - **C:** inside this repo (`plugins/ailang-world/` plus a root `.claude-plugin/marketplace.json`), versioned and CI-tested with the CLI it documents.
  - **Recommendation: C.** The skill's commands can be checked against `ailang-worldd` help in the same CI, and a release tag versions both together. Mark must confirm that CLAUDE.md's "never copy skills into this repo" covers the shared mission skills, not World's own product plugin.
- **D-189-4 — How agents get World's tools.**
  - **A:** `episode start` prints a launch with a private 0600 config and a strict tool set, the `trial.py` shape. The plugin carries skills and commands only.
  - **B:** additionally ship a global plugin MCP entry using Claude Code's `headersHelper` (V6) to read the token for the current worktree's episode. That makes a non-evidence-grade mix of native and World tools available in every session.
  - **Recommendation: A for v1;** B as a later opt-in after an M0 probe of `headersHelper`'s output contract.
- **D-189-5 — How a proposed profile becomes active.**
  - **A:** `profile init` writes `*.proposed.json`, and the operator renames it. No TTY fence; `serve`'s startup probe validates it.
  - **B:** a TTY y/N inside `profile init`, like the mint.
  - **Recommendation: A.** A profile is operator configuration, checked at `serve` startup and named and digested in the mint (D-189-2 A). A second fence adds a step without adding authority.

## 7. Acceptance criteria

- **AC1:** with a daemon serving episode X, `episode start Y` mints Y and Y's agent makes a `workspace-exec` call, with no restart, while X's next call still succeeds. Measured on the rig.
- **AC2:** `episode start` with no controlling terminal refuses before any worktree, row or file exists. This is MUT-FENCE-BYPASS kept, with a test.
- **AC3:** `episode start` refuses a daemon whose tool pin or build differs from its own, naming the rebuild.
- **AC4:** the token appears in no file other than the 0600 session file and the per-launch MCP config. A test greps the worktree, evidence and stdout/stderr.
- **AC5:** `profile init` on `platform-source` and on `ailang-world` proposes files that pass `ParseExecProfile`. Its command set covers the measured profiles' file-scoped test, lint and vet commands.
- **AC6:** the plugin installs from its marketplace with `claude plugin install`. CI fails if a command the skill names is missing from `ailang-worldd` help.
- **AC7, the re-measure:** one attended episode on platform-source, started with `episode start` alone. Count the operator steps and compare with P1's eleven; friction #1, #2 and #13 are each closed or re-scored.

## 8. Milestones (sketch; replanned after the rulings)

- **M0:** probes.
  - Credential-table joins and budget reads (for D-189-1 C).
  - The `headersHelper` contract, only if D-189-4 B is ever wanted.
  - A daemon build/pin self-report for AC3: what `/v1/health` already returns.
- **M1a/M1b:** the credentials store split (C), covering schema, migration, resolver and `session new/list/revoke`.
- **M2:** profile binding at mint (D-189-2 A). `serve` loads the profiles directory.
- **M3a/M3b:** `episode start`: fence and mint reuse, the daemon check and start, the launch line.
- **M4 — DONE (2026-10-19, episode w189c; §8c):** `profile init` for Python and Go, plus `package.json`.
- **M5 — DONE (2026-10-20, episode w189d; §8d):** the plugin: two skills, three commands,
  marketplace manifest, and the CI help-check.
- **M6:** docs (QUICKSTART §0 and a new "Use World on your repo" page), then the AC7 attended run.
  - **Docs half DONE (2026-10-20, episode w189e):** `docs/USE_WORLD_ON_YOUR_REPO.md`, linked
    from `README.md`, QUICKSTART §0 and the plugin README. Its `ailang-worldd` lines are checked
    against the CLI help by `TestUseWorldPageCommandsMatchTheCLIHelp`
    (`cmd/ailang-worldd/use_world_page_test.go`, reusing `checkWorlddLine`). It counts six
    human acts per episode against P1's eleven.
  - **Still to come: AC7**, the attended run on platform-source that follows the page verbatim
    and counts the operator's steps.

## 8a. M0 findings (row 189, 2026-10-19, read before the M1 change; lines are pre-change)

- **No join (C stands, B not needed).** The only statements on `session_credentials` were the four
  credential methods, each alone in its own statement or transaction: INSERT `store.go:1317`
  (`MintSession`, `:1300`), point SELECT `:1341` (`ResolveSession`), DELETE `:1370`
  (`RevokeSession`), list SELECT `:1398` (`ListSessions`). `Commit` (`store.go:1076`) and the
  journal never name it; `schema.sql:107` has no `REFERENCES` to or from it (the only foreign keys
  are `commit_objects`, `schema.sql:121-122`) and the schema has no triggers.
- **Budget does not live on the credential row.** `grants_json` (with each grant's minted `Budget`)
  is written once at mint (`authority/mint.go:52-64`) and never updated. The remaining budget is
  computed in Go, never in SQL: `coordinator/effectful.go:105-117` (`seedGrants`) subtracts the
  episode's durable spend, which `store/journal.go:1079-1111` (`EffectSpend`) sums from the world
  journal's effect intents. `broker/broker.go:94-97` debits only the in-memory session.
  `cmd/ailang-worldd/demo_seed.go:151-169` does the same subtraction. Spend stays in the world store
  and nothing moves with the credential.
- **`/v1/health` today** (`host/daemon/daemon.go:569-590`, handler `:973-981`): `status`,
  `daemon_version` (the constant `Version = "0.1.0"`, `daemon.go:70`, not a build id), `db_path`,
  `interpreter_ref`, `interpreter_version` and `degraded_episodes`. It reports **no** tool-binary
  pin (the daemon knows `workspace.binRef`, `daemon.go:693`, but does not return it), no VCS revision
  or build id, no exec-profile digests and no credentials path. AC3 needs at least the tool ref and
  a build identity added.

## 8b. M2 + M3a/M3b as built (row 189, `episode start` and the profile binding)

```
ailang-worldd episode start <ep> --repo <dir> [--branch world/<ep>] [--profile <project>] [--task <file>]
    [--agent claude|codex] [--model <m>] [--ttl N] [--preset|--grant|--budget as session new]
    [--db …] [--credentials-db …] [--workspace-root …] [--sessions-dir ~/.ailang/world/sessions]
    [--exec-profiles-dir ~/.ailang-world-exec/profiles] [--addr 127.0.0.1:47330]
```

- **Order** (`cmd/ailang-worldd/episode.go`, `runEpisodeStart`): validate (grammar, `--agent`,
  `--addr` loopback, `--task`, the session/MCP files absent, then `prepareMint`: store, workspace
  root, repo, `--profile` parsed with `ParseExecProfile`, never probed) → `GET /v1/health` and
  `checkDaemon` (build identity, tool ref against the manifest pin, credentials path, the profile's
  digest loaded and not left out; no daemon = refusal printing the `serve` line) →
  `provisionEpisode` (`session_new.go`), the one fence-and-mint function `session new` also runs:
  TTY → y/N → credentials lock → worktree → `authority.MintBound` → token file → MCP config →
  rollback → launch line. It never starts a daemon (M3c), runs the agent, commits or pushes.
- **Default grants:** the se-tools preset without `Workspace.Exec`, plus it with `--profile`.
- **Binding (D-189-2 A):** credentials schema 2 adds nullable `project` and `profile_digest`
  (version-1 rows stay, unbound). Only `serve` upgrades a version-1 store, so a session verb of
  this build never hands an older serving daemon a store it would refuse. `profileFor` resolves
  per request: the episode's live credentials' binding → `--exec-episode-project` → the only
  profile. A bound digest the daemon did not load refuses every exec ("re-mint").
- **`/v1/health` adds** `build` (`vcs_revision`, `vcs_modified`, `executable_sha256`,
  `go_version`), `tool_ref`, `credentials_db`, `exec_profiles` and `exec_profiles_excluded`.
- **Profiles directory:** `serve --exec-profiles-dir` (default `~/.ailang-world-exec/profiles`
  when `--exec-sandbox` is given without `--exec-profile`) loads every `<project>.json`, skips
  `*.proposed.json`, and leaves out — never refuses startup for — a profile that fails to parse or
  probe; the `--exec-profile` flags keep their all-or-nothing rule.

## 8c. M4 as built (row 189, `profile init`, D-189-5 = A)

```
ailang-worldd profile init --repo <dir> --project <name> [--out <file>] [--profiles-dir ~/.ailang-world-exec/profiles]
```

- **Where:** `cmd/ailang-worldd/profile_init.go`. `main.go` gains only the `profile` dispatch
  case, and `host/broker` exports `ExecProjectGrammar`/`ExecProjectOk`.
- **Detection:**
  - `go.mod` → `go-test`, `go-vet`, `gofmt-l`.
  - `pyproject.toml`/`uv.lock` → `pytest-file`, `ruff-file`.
  - `package.json` → `vitest-file`/`jest-file`, by `devDependencies`.
  - `Makefile`: its targets are listed, and nothing is proposed from them.
  - The repo root wins. With no marker there, init reads one directory down, ranking Go, then
    Python, then node, and that directory becomes `root` (platform-source's `backend`).
- **Paths are resolved, never guessed:**
  - **Go:** `go env GOROOT GOMODCACHE GOVERSION` run in the root. `path` is `GOROOT/bin`, the
    toolchain go.mod switched to. `env` is `GOFLAGS`/`GOPROXY`/`GOTOOLCHAIN`/`GOTELEMETRY`/
    `GOMODCACHE`, and `read_roots` is GOMODCACHE, plus GOROOT when it lies outside it.
    `caches.GOCACHE` starts empty.
  - **Python:** the install directory of every hop on `<root>/.venv/bin/python`'s link chain
    (or `uv python find`'s) is a read root, and so is the realpath's. The probe is
    `<realpath> --version`.
  - **Node:** node's directory on `PATH`. Its install directory is a read root when it lies
    under HOME.
  - A toolchain that does not resolve is skipped, with its reason.
- **Shapes:** `initShapes` holds the command entries, bound by
  `TestProfileInitShapesMatchTheMeasuredProfiles` to
  `cmd/ailang-worldd/testdata/profile-init/{ailang-world,platform-source}.commands.json`.
  - The ailang-world copy was observed through the live profile (digest `sha256:ca1abce3…`).
  - The platform-source copy is a reconstruction (see its `source` field). The episode could
    not read the file, so the operator owes the real entries.
- **Writing:** the proposal must pass `ParseExecProfile` before any byte or directory is
  written. It is written with O_EXCL, and init prints the digest and the activating `mv`.
- **Refusals:**
  - a `--project` outside the grammar;
  - an existing output;
  - an `--out` naming the active `<project>.json`;
  - no detectable or resolvable toolchain;
  - a parser refusal.

## 8d. M5 as built (row 189, the `ailang-world` plugin, D-189-3 = C, D-189-4 = A)

```
.claude-plugin/marketplace.json                       marketplace ailang-world-marketplace, owner Sunholo
plugins/ailang-world/.claude-plugin/plugin.json       ailang-world 1.0.0-rc.1, author Sunholo; no mcpServers
plugins/ailang-world/skills/ailang-world/SKILL.md            the agent side
plugins/ailang-world/skills/ailang-world-operator/SKILL.md   the human side
plugins/ailang-world/commands/{world-start,world-watch,world-why}.md   print-only
plugins/ailang-world/README.md                        the two install lines
```

- **No MCP entry** (D-189-4 A): no `.mcp.json`, no `mcpServers`. An agent gets World's tools only
  from the launch line `episode start` prints. The plugin holds no token, address book, personal
  name or home path.
- **Human acts stay human:** the operator skill and `/world-start` print `episode start` and say
  the human runs it at a terminal and types `y`; they never run it, feed it input, or run
  `session new`/`session mint`. `world-publish` lines are printed, never run. `/world-start`'s
  `allowed-tools` is `Bash(git rev-parse:*)` and `Bash(ls:*)`; the other two have none.
- **CLAUDE.md** reads "never copy skills into this repo" as the shared mission and fleet skills,
  not this plugin (D-WORLD-85 / D-189-3).
- **The help-check (AC6's CI half):** `TestPluginCommandsMatchTheCLIHelp`, in
  `cmd/ailang-worldd/plugin_help_test.go` and `cmd/world-publish/plugin_help_test.go`, over
  `host/plugindoc` (a test helper beside `host/proctest` and `host/ptytest`, not under `cmd/`,
  which `TestEveryCommandDirectoryIsAProtectedGroup` treats as commands):
  - it walks every `.md` under `plugins/ailang-world/` and takes each fenced line (joined across
    a trailing `\`) and each inline code span that starts with the program name, after a leading
    `$ ` and `VAR=value` assignments;
  - the verb path is the longest path in the top-level usage the line starts with (the global
    `--addr` before the verb is skipped), and every `--flag` must be in that verb's own `--help`
    (ailang-worldd, run in-process) or its usage block (world-publish);
  - controls: the walk must find at least 16 ailang-worldd and 2 world-publish lines, and in-test
    unknown verbs and flags (`episode frobnicate`, `profile frobnicate`, `episode report
    --no-such-flag`, `world-publish frobnicate`, …) must refuse through the same checker.
    Appending the four grader mutants to the operator skill made both tests fail, naming the
    file and line; removing them made both pass;
  - known limit: the check proves every printed flag exists in the verb's help, not that the
    verb's required flags are present, so a line missing `--db` (as `session revoke` needs) passes;
  - the top-level usage now lists `profile init` (`TestUsageListsProfileInit`), and `profile`
    prints its own help.
- **AC6's install half** (`claude plugin validate --strict`, `claude plugin install` into a
  throwaway config dir, `plugin details`) needs the `claude` CLI, which the episode could not
  run; the operator's grader measures it.

## 9. V-rows

| # | claim | evidence |
|---|---|---|
| V1 | Today's per-episode mint script: stop the daemon via `serve.pid`; confirm `/v1/health` is down; copy the rebuilt binary; `session new` with 9 `--grant` flags, `--ttl 43200`, `--out`, `--branch`, `< /dev/tty` | `~/.ailang/state/trial-2026-10-16/MINT_BLOCK_7.sh` (read) |
| V2 | The serving daemon carries 2 `--exec-profile` and 8 `--exec-episode-project` flags plus tool pin, sandbox, node and output cap | `ps -o command=` on the rig daemon (earlier in the attended session, 2026-10-19), and the restart command used for `w190` |
| V3 | `trial.py run`: spec → log range → spawn with a private MCP config → graders unsandboxed → outcome records; `pr-body` never pushes | `scripts/trial/trial.py:1-27`, `:472`, `:555` |
| V4 | `session new` order: validate → `env.openTerminal()` → y/N → `openStoreForSession` → `worktree.Add` → `sessionMint` → rollback | `cmd/ailang-worldd/session_new.go:1-17, 230-300` |
| V5 | The writer lock is `flock(LOCK_EX\|LOCK_NB)` on `<db>.lock`; contention → `WriterAlreadyActive` → "stop the daemon first"; `session list` reads through a read-only handle while the daemon runs | `host/store/writer_lock_unix.go:32-48`; `cmd/ailang-worldd/session.go:110-121`; `session_new.go:98-100` |
| V6 | The TTY fence opens `/dev/tty`, which an agent harness cannot (ENXIO measured); Claude Code 2.1.303 has an MCP server field `headersHelper` ("Headers helper script", absolute path, requires user consent) | `cmd/ailang-worldd/session.go:39-46`; `cmd/world-publish/tty.go:26-40`; `verification/world-row189-design-2026-10-19/v6-claude-2.1.303-mcp-headershelper.txt` |
| V7 | With several profiles, an unmapped episode gets no profile; the mapping is a `serve` flag | `host/daemon/workspace.go:152-160`; `host/broker/exec_startup.go:55, 150`; `cmd/ailang-worldd/main.go:417` |
| V8 | A stale rig build served the v0.52.1 tool pin after #308 | friction #13 (`design_docs/verification/world-dogfood/friction.md`) |
| V9 | Marketplaces installed: `claude-plugins-official`, `anthropic-agent-skills`, `ailang-parse-marketplace` (`sunholo-data/docparse-skill`), `ailang-marketplace` (`sunholo-data/ailang_bootstrap`, ref `dev`); no `ailang-world` plugin | `~/.claude/plugins/known_marketplaces.json` (read) |
| V10 | The ailang-parse plugin layout: `.claude-plugin/marketplace.json` → `plugins/ailang-parse/{.claude-plugin/plugin.json, .mcp.json (type http, url), skills/}` | `find` over `~/.claude/plugins/marketplaces/ailang-parse-marketplace -maxdepth 3` |
| V11 | Profile sizes: 36–159 lines across 5 files; `platform-source` commands use `sh -c` wrappers, `regex`/`sep` and `bool` flag classes, `testfile` positionals, uv interpreter roots | `wc -l ~/.ailang-world-exec/profiles/*.json`; `pytest-file` read |
| V12 | The profile schema is `world/exec-profile/v1`, parsed by `ParseExecProfile`; positional kinds include relpath, pkgpattern, testfile | `host/broker/exec_profile.go:27-28, 130, 371` |
| V13 | `session_credentials` is referenced in `host/store/store.go` (mint/get/revoke/list, `:1284-1400`), `host/authority/resolver.go`, `host/daemon/daemon.go:385`; DDL `host/store/schema.sql:107`; schema version 4 | `grep -rln session_credentials`; `grep -n 'CREATE TABLE' host/store/schema.sql`; `host/store/store.go:44` |
| V14 | The platform-source episodes' repo is a prep clone whose origin is the operator's checkout | `trial/trial*.json` `repo` fields; `git -C …/trial-prep/platform-source remote get-url origin` |

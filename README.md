# AILANG World

**An AI-native semantic operating environment.**

[![CI](https://github.com/sunholo-data/ailang-world/actions/workflows/ci.yml/badge.svg?branch=dev)](https://github.com/sunholo-data/ailang-world/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

> Unix made everything a file.
> **AILANG World makes everything a typed state transition.**

AILANG World is a **semantic database whose transaction language is
[AILANG](https://github.com/sunholo-data/ailang)** — a deterministic programming language
designed for AI code synthesis and reasoning. Where a traditional operating system manages
processes, files, and devices, World manages **goals, typed state, capabilities, effects,
evidence, budgets, contracts, provenance, and AI proposals**.

The kernel is an immutable, content-addressed world graph. State changes only through
**transitions** — pure, type-checked AILANG functions verified *before* they run. AI agents
never modify the world directly: they **propose**, a deterministic verifier (types + Z3
contracts) **checks**, and only verified, authorized, budgeted proposals **commit** — each
one leaving evidence and a replayable trace. Humans don't operate the system; they govern
it: express goals with budgets attached, decide what verification couldn't settle, and ask
the world *why* anything happened.

The underlying OS executes bytes. **AILANG World executes intent.**

**Docs:** [sunholo.com/ailang-world](https://www.sunholo.com/ailang-world/) — concepts, getting
started, the agent guide and the CLI/HTTP reference — and
[use cases](https://www.sunholo.com/ailang-world/docs/use-cases). For agents:
[`llms.txt`](https://www.sunholo.com/ailang-world/llms.txt) and
[`AGENTS.md`](https://www.sunholo.com/ailang-world/AGENTS.md).

## Try it

Serve World's ten coding tools to an agent, one step at a time (needs the two pinned
AILANG binaries — `go run ./cmd/ailang-worldd setup` installs both, verified, and `doctor`
checks the machine; the two attended steps, publish and mint, need a human at a terminal):

```sh
tools/attended/se_smoke.sh prepare   # build the CLIs, a scratch store and an episode worktree
tools/attended/se_smoke.sh publish   # attended: prints the QUICKSTART §9 Publish line to paste
tools/attended/se_smoke.sh mint && tools/attended/se_smoke.sh serve   # mint: you confirm at the terminal
B=~/.ailang/se-smoke; $B/bin/ailang-worldd tools list --session $B/session
$B/bin/ailang-worldd call ailang-read --session $B/session --arg path=hello.ail --json-out | $B/bin/ailang-worldd why -
```

`call` is a built-in MCP client; `why` walks a result back through its log entry, plan and
effect records, checking every link. `log tail`, `provenance`, `outcome record` and the
`setup`, `doctor` and `session new`/`list`/`revoke` helpers complete the developer CLI.
The full runbook, with every flag explained, is [docs/QUICKSTART.md](docs/QUICKSTART.md) §9;
the [coding-tools guide](https://www.sunholo.com/ailang-world/docs/getting-started/coding-tools)
connects pi and Claude Code.

Effect receipts let operators distinguish effects that were never dispatched
from effects whose durable intent has no outcome. Treat an `indeterminate`
receipt as fail-closed: do not retry automatically. A crash after the effect
record write but before the outcome append can leave a record that requires
deterministic reconciliation before the receipt can be resolved.

### Screenshots

The [light workbench](website/static/img/v1/workbench-light.png) and the [decision packet](website/static/img/v1/packet.png) are static captures of recorded history. The attended v1.0.0 release capture publishes them, and these two links resolve only after that capture is promoted. The full set, its captions and what the images do not show are on the [release screenshots page](https://www.sunholo.com/ailang-world/docs/getting-started/marketing-capture); the reproduction steps are in [docs/QUICKSTART.md](docs/QUICKSTART.md#13-capture-the-release-screenshots-attended) §13.

## Why

Extrapolate the model trend — smarter, faster, cheaper, increasingly on-device — and
generation stops being scarce. What stays scarce is **trust, running at the speed of the
intelligence it governs**. Workflow engines schedule agents but can't examine a plan before
it runs and *prove* it stays inside its declared effects — that takes a language. The field
now agrees: **Language-Based Agent Control** is a live research tradition — Odersky's
capture-checked capabilities in Scala 3 ([Best Paper, CAIS 2026](https://dl.acm.org/doi/10.1145/3786335.3813127)),
[Etas](https://arxiv.org/abs/2607.17780)'s effect-typed agent language, and more. AILANG
World's position in that field: a **purpose-built language with a working compiler** (not an
embedding, not a paper), Z3-verified contracts gating every push, **operating a persistent
governed substrate** — world graph, receipts, bit-for-bit replay — that the language tier
doesn't have and the substrate tier can't prove. See
[design_docs/POSITIONING.md](design_docs/POSITIONING.md) for the honest landscape. World is
that bet, operationalized:

- **Verification before execution** — proposals are proven well-formed before any effect fires.
- **Constructive replay** — determinism guaranteed by the language, not reconstructed from logs.
- **Authority as types** — capabilities checked where transitions are defined, not bolted onto a gateway.
- **Local-first** — one daemon, SQLite, zero cloud in the core. Cloud and multi-device sync are effect-handler extensions, never architecture.
- **Protocol-native** — MCP, A2A, and open agent-UI protocols at the boundary; no invented wire formats.
- **Falsifiable** — see [the floor, and the value](#the-floor-and-the-value) below. We name the condition under which this project parks itself.

## The documents

| Document | What it is |
|---|---|
| [DESIGN.md](design_docs/DESIGN.md) | The thesis — full architecture: world graph, transitions, proposals, effect broker, capabilities, budgets, scheduler, self-modification, milestones. All AILANG snippets in it **compile** ([sketches/](design_docs/sketches/) are CI-checked). |
| [SCENARIOS.md](design_docs/SCENARIOS.md) | Day-in-the-life walkthroughs of the *human* side — the approval inbox, goal composing, provenance walks, speculative what-ifs. |
| [AN-AGENTS-CASE.md](design_docs/AN-AGENTS-CASE.md) | A first-person statement by an AI (Claude) on why it would choose this environment — the user constituency, in its own voice. |
| [REFERENCES.md](design_docs/REFERENCES.md) | 43 verified prior-art references — Datomic, Urbit, seL4, Nix, capability security, local-first, provenance, agent protocols — each annotated with what World steals and which pitfall it avoids. |
| [world-mission.md](design_docs/world-mission.md) | The mission charter — the checkable bar for "World 1.0", guardrails, the live work queue, and the ledger of attended human decisions. |
| [QUICKSTART.md](docs/QUICKSTART.md) | The operator runbook: build, serve, commit, publish transitions, and serve the coding tools to agents. Its command lines are bound to the CLIs by tests. |
| [USE_WORLD_ON_YOUR_REPO.md](docs/USE_WORLD_ON_YOUR_REPO.md) | Run a coding agent on your own repo through World, step by step: plugin, profile, `serve`, `episode start`, watch, review, finish. Its command lines are bound to the CLI by a test. |
| [AI-EMPLOYEE.md](design_docs/AI-EMPLOYEE.md) | The first external application — a long-running engineering agent with its own accounts and a human owner. Use case, a gap inventory against HEAD, and the rows that unblock a resident. Draft, unratified. |

## How this repo is built

This repository practices what it preaches: it is advanced by an **autonomous mission
loop** — a scheduled outer loop in which AI agents design, plan, execute, and evaluate one
work item per iteration, with machine verification gating every landing and a human
ratifying the direction. Every iteration reports publicly to a
[bookkeeping issue](https://github.com/sunholo-data/ailang-world/issues?q=is%3Aissue+%22mission+bookkeeping%22)
(one thread per week). Alongside the loop, **attended sessions** — the owner working live with
an agent — build the larger rows, make the rulings the loop must not make for itself, and record
each one in the charter's decision ledger. The workflow
this loop runs *by convention* is precisely the workflow World exists to make *first-class*
— typed proposals, verified commits, budgeted human attention. World is bootstrapping
itself into existence using the process it formalizes.

Roles today: AI agents author designs, write code, and review each other's work
(generator ≠ judge, enforced); a human ([@MarkEdmondson1234](https://github.com/MarkEdmondson1234))
holds the constitution — the bar, the budgets, and the approvals. On this public repo,
**only directives from that account steer the loop**; all other input is welcome but never
executed. Mission-loop machinery is shared infrastructure from the
[AILANG repo](https://github.com/sunholo-data/ailang) and is not modified here.

## Status — honest, not promotional

**Phase: the 1.0 release candidate (`v1.0.0-rc.1`).** As of 2026-10-27, all seven clauses of
the [ratified 1.0 bar](design_docs/world-mission.md) are met, including the two that measure
World against the status quo (clause 4 by D-WORLD-77, clause 5 by D-WORLD-80). `v1.0.0-rc.1` is
the current GitHub pre-release; per D-WORLD-81, `v1.0.0` follows once rows 156 (protected test
paths), 150 (marketing capture) and 151 (agent-generated projections) land.

**Built and working** (on `dev`, CI-gated):

- **The kernel and daemon** — an immutable, content-addressed world store over SQLite, an
  append-only transition log with an interpreter pin in every entry, deterministic replay, one
  loopback-only daemon (`ailang-worldd`) with a REST API and a CLI. Clauses 1 and 2.
- **Explicit authority** — every effect goes through the broker with a session capability and a
  budget check, and every effect result is recorded as replay input. Clause 3.
- **Protocol-native boundary** — the transition registry is served as MCP tools and as skills on
  an A2A agent card, both filtered per session by capability. Clause 6.
- **Controlled self-modification, proven once** — World's own package, `world/core`, went
  through World's propose → verify → commit pipeline and was published to the AILANG package
  registry by an attended human step (`world/core@0.1.1`, 2026-10-02). Clause 7.
- **Ten coding tools for agents** (queue row 134, landed 2026-10-03, PR #180; row 140 added
  the ninth, row 191a the tenth): `ailang-read` (whole or a line range), `-write`, `-edit`,
  `-delete`, `-check`, `-run`, `builtins-search`, `examples-search`, `ailang-cli` and
  `workspace-exec`, served over MCP and A2A. Each call runs
  its plan as a pure World transition, runs exactly one brokered, budgeted effect confined to the
  episode's worktree, and commits one log entry. `workspace-exec` runs only commands from an
  operator exec profile (`go test`, `npm test`, `pytest`, …), sandboxed by srt, and runs nothing
  unless the daemon starts with `--exec-profile`. Toolchain results carry a `worktree_delta`
  (files created, modified, deleted) that `why` prints. pi and Claude Code both did real work
  through them with their built-in tools disabled, and Claude Code and codex ran the floor's
  World arm through them.
- **A developer CLI** (row 138): `ailang-worldd tools list`, `call`, `why`, `log tail` and
  `provenance`, replacing hand-written curl with bearer headers and SSE parsing; `setup` (fetch
  and hash-verify both pins) and `doctor`; `session new`/`list`/`revoke` (`session new
  --allow-non-git` adopts a non-git directory as a degraded episode); `outcome record` (row 182:
  a graded task's outcome as a `world/task-outcome/v1` object a provenance walk can reach); and
  `demo seed`.
- **Two pinned AILANG binaries**: the `.ail` interpreter that runs every transition plan is
  pinned at **v0.41.0** (epoch 1; its hash is in every log entry, so it moves only by a
  deliberate epoch change); the tool binary the coding tools run is pinned at **v0.54.15** and
  refused at any other release (D-WORLD-57; v0.52.1 before row 139 moved it, kept in the
  manifest as a `retired` row). Both are rows of the release manifest
  [`pins/ailang-releases.json`](pins/ailang-releases.json), compiled into the binaries;
  `world-publish epoch nominate` names a new release in a store's epoch registry, as a candidate
  of the current epoch (headless) or as a new epoch (attended).
- **A docs site** at [sunholo.com/ailang-world](https://www.sunholo.com/ailang-world/).

**The 1.0 bar** ([charter, "The bar"](design_docs/world-mission.md)):

| # | Clause | State |
|---|---|---|
| 1 | Deterministic kernel | ✅ met |
| 2 | Local-first daemon | ✅ met |
| 3 | Explicit authority end to end | ✅ met |
| 4 | Resident-agent non-inferiority floor | ✅ met — the 2026-10-15 FINAL holds: Claude Code Δ 0 (69/69 vs 69/69), overhead −17.9%; codex Δ 0 (69/69), overhead −5.7%; 414/414 graded cells passed ([evidence](design_docs/verification/world-floor-final-2026-10-15/README.md)) |
| 5 | Human surface with provenance teeth | ✅ met — three *unseeded* real questions answered by walks verified against sealed ground truth: Q4 54.6 s, Q6 25.3 s, Q8 39.0 s (the live baseline also verified, in 67.0 s, 87.5 s and 116.8 s) ([evidence](design_docs/verification/w-provenance-teeth-unseeded-baselined-walk/index.md)) |
| 6 | Protocol-native boundary | ✅ met |
| 7 | Controlled self-modification, proven once | ✅ met |

**Landed** from the ratified order (D-WORLD-58: 135 → 138 → 140 → 93 → 114 → 139 → release):
rows 135 (`ailang-run` stdin, argv and per-call capabilities), 138 (the developer CLI), 140
(`workspace-exec`), 93 (the floor) and 114 (the provenance walks).

**In progress**, on two tracks:

- **The mission loop, toward `v1.0.0` (D-WORLD-81)** — rows **150** (marketing capture) and
  **156** (protected test paths, so an agent cannot edit the tests that judge it through
  `workspace-exec`) are in sprint, with 156's remaining milestones held on D-WORLD-86, a ruling
  open for the owner; row 151 (agent-generated projections) follows. Row 139, the
  interpreter epoch upgrade path, has M1/M2 landed: the release manifest
  `pins/ailang-releases.json` and `world-publish epoch nominate`.
- **The attended track: all agent work through World** (D-WORLD-94, set by the owner
  2026-10-26) — every agent task, headless or interactive, on any repo, runs through World as
  the normal way, not the trial rig. The roadmap is
  [w-all-agent-work-through-world](design_docs/planned/w-all-agent-work-through-world.md); the
  guide is [Run your agents through World](https://www.sunholo.com/ailang-world/docs/guides/agents-through-world).
  - **Landed 27 Oct, rows 223a–223d**, each built through World: launchers and `hosts init` for
    Claude Code, Codex and pi, with `doctor` checking the daemon, session and host before a
    launch and naming the fix (223a, PR #375); the herdr plugin, which starts a confined agent
    pane after the owner's typed `y`, or a labelled plain fallback when World is down (223b,
    #377), with key bindings and a World panel (223c, #379); and in-agent awareness, so each
    agent knows its episode, worktree, grants and session lifetime at launch (223d, #381).
  - **Next,** in the roadmap's order: row 226 (a per-episode value scorecard), then 219 (World as
    a supervised user service with a home store and a one-command upgrade; design:
    [w-worldd-home-service](design_docs/planned/w-worldd-home-service.md), #386), then 217
    (one-command episodes) and 218 (`episode finish`, start to merged PR).

**Work made through World.** The first real change made by an agent working only through World's
tools is [sunholo-data/daneel#325](https://github.com/sunholo-data/daneel/pull/325) (merged
2026-10-03). By convention (D-WORLD-60), every PR made through World carries the GitHub label
**`ailang-world`** and a `World-Provenance: store=… episode=… entries=<from>-<to>` line, which
`ailang-worldd provenance` prints, so `ailang-worldd why <entry>` can walk any such commit back to
its tool calls, plans and effect records. Filter with `label:ailang-world`.

## The floor, and the value

Agents are World's *residents*, not its destination — so World is not measured by whether it
makes a coding agent pass more benchmarks. That would be like benchmarking git by typing
speed. The comparison class for World's value is the **operational status quo**: the scripts,
schedulers, conventions, and human vigilance that do this job by hand today.

Two burdens, both in the charter, both falsifiable:

- **The floor (do-no-harm kill switch)**: two reference agents from different providers
  (Claude Code and codex), each running the same benchmarks with native shell tools vs
  World's MCP transition tools, must both hold **non-inferiority** — pass-rate within 2
  points, overhead within 25%, thresholds fixed *before* the work started, with a stability
  precondition on the baseline so a flaky harness can't fake the verdict either way. A
  substrate that taxes its residents dies before its value accrues; if World fails this
  floor, World parks, and says so here.
- **The value (what World is for)**: capability a shell cannot express at any pass-rate —
  real "why did this happen" questions answered in minutes by provenance walk instead of
  archaeology sessions (measurably: ≥3 real questions, ≤5 minutes each), incident *classes*
  eliminated structurally, and ultimately the mission loop that builds this repo migrating
  onto World and beating its own manual baseline on incidents and human attention.

A trust substrate that can't pass its own evidence bar has no business asking for yours.

## Run the local world daemon

Build the single daemon/client binary, then start one writer for a database:

```sh
go build -o ailang-worldd ./cmd/ailang-worldd
./ailang-worldd serve --db ./world.db
```

The REST listener is loopback-only (default `127.0.0.1:7644`) and each file-backed
database permits exactly one writer process. A second daemon or embedded writer
fails closed; read-only store users may coexist.

At startup the daemon performs a bounded integrity scan of persisted log and
world references. A clean, complete scan prints nothing. It reports only when
it has something to say: `integrity_scan_complete … holes=N` (with N > 0, after
one `integrity_hole` line per unreadable row) means both tables were fully
scanned within the configured bounds and found holes. `integrity_scan_incomplete`
means a row or time budget stopped the scan;
the message includes continuation cursors and counts for the scanned prefix, so
it must not be read as a clean bill of health. Individual `integrity_hole`
messages identify unreadable historic rows. A store with a historic hole still
serves readable data: affected reads continue to fail loudly. Repair is ruled
out by design because rewriting committed immutable history would destroy the
evidence. Detection and honest operator visibility are the deliverable.

The same binary is the bounded-timeout client:

```sh
./ailang-worldd health
./ailang-worldd head
./ailang-worldd world get sha256:<digest>
./ailang-worldd object get sha256:<digest> --payload
./ailang-worldd log get 0
./ailang-worldd log range --from 0 --limit 100
./ailang-worldd registry get world/epoch-registry/v1
./ailang-worldd log tail --follow
./ailang-worldd commit --file commit.json --session <file>
./ailang-worldd outcome record --file outcome.json --session <file>
./ailang-worldd tools list --session <file>
./ailang-worldd call <tool> --session <file> --arg k=v
./ailang-worldd why <index|head|sha256:<ref>|a2a:<id>|mcp:<id>|rest:<id>|->
./ailang-worldd provenance --episode <ep>
./ailang-worldd doctor
```

`setup` installs the two pinned AILANG binaries, and the `session` verbs act on the
credentials store beside the world store (`credentials.db` next to `--db`, its own file with
its own lock) directly, so they all work while the daemon serves:

```sh
./ailang-worldd setup
./ailang-worldd session new <ep> --db ./world.db --repo . --preset se-tools   # attended, TTY-fenced
./ailang-worldd session mint --db ./world.db --episode <ep> --grant EFFECT=SCOPE:BUDGET   # attended, TTY-fenced
./ailang-worldd session list --db ./world.db
./ailang-worldd session revoke --db ./world.db <credential_id>
```

`POST /v1/commit`, `outcome record` and every tool call need a session credential, minted by
the attended, TTY-fenced `session new` or `session mint`; reads are unauthenticated and
loopback-only. `ailang-worldd <verb> --help` (or `help <verb>`) prints any verb's help.

Use global `--addr http://127.0.0.1:<port>` before a client verb when the daemon
uses another loopback port. `--addr` is not valid with `serve`; use `--bind`.

## Effect broker operator boundary

The effect broker is the in-process authority and accounting boundary for
effect requests. The daemon wires it in through the **coordinator**: every MCP
`tools/call` and A2A `tasks/send` runs the transition's plan in the pinned
interpreter, then opens a session-scoped broker binder (`broker.OpenBinder`,
over the episode's workspace handlers) that checks the session's capability
and budget for each declared effect before a handler runs. The broker has no
REST route or CLI verb of its own: an effect is reachable only as a declared
effect of a published transition, called with a session that holds the grant.
A transition that declares an effect the daemon has no handler for is refused
before anything runs. Every allowed, denied, or failed effect produces an
immutable, content-addressed effect-record object. Successful result bytes are
stored the same way. Replay mode reads those records and never dispatches a
live handler.

Subprocess effects use the capsule floor: a pinned executable, empty
capabilities, an FS jail, a scrubbed environment, a wall-clock timeout, and an
output cap. The coding tools add AILANG's own policy layer on top: the archived
v0.54.15 tool binary runs under a rendered policy that roots every path at the
episode worktree, keeps `.git` and the deny list read-only, and gives no shell
(measured facts: the site's
[tool confinement](https://www.sunholo.com/ailang-world/docs/security/tool-confinement)
page). The one command-running tool, `workspace-exec`, runs only the commands an
operator's exec profile lists (`serve --exec-profile`), never a shell line, under
the srt sandbox: writes confined to the worktree and the episode's caches, no
network, and a startup probe each profile must pass. This is a process-safety
floor, not full isolation: there is no container, VM, or memory or CPU limit.

## Relationship to AILANG

[AILANG](https://github.com/sunholo-data/ailang) is the language: deterministic, explicit
effects in function signatures, Hindley–Milner types, Z3-verified contracts, built for AI
authorship. World is the operating environment that language makes possible. Language gaps
discovered here are routed upstream as issues — World never forks or works around the
compiler.

## Following along

Watch the [bookkeeping issues](https://github.com/sunholo-data/ailang-world/issues?q=is%3Aissue+%22mission+bookkeeping%22)
— every mission iteration reports to that week's thread. The append-only iteration log lives at
[world-mission-log.md](design_docs/world-mission-log.md). Issues and discussion are open;
if you're building in the same space, [REFERENCES.md](design_docs/REFERENCES.md) is the
map of whose shoulders we're standing on.

## License

[Apache 2.0](LICENSE)

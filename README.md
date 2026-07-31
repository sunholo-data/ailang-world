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

Effect receipts let operators distinguish effects that were never dispatched
from effects whose durable intent has no outcome. Treat an `indeterminate`
receipt as fail-closed: do not retry automatically. A crash after the effect
record write but before the outcome append can leave a record that requires
deterministic reconciliation before the receipt can be resolved.

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
- **Falsifiable** — see [the value gate](#the-value-gate) below. We name the condition under which this project parks itself.

## The documents

| Document | What it is |
|---|---|
| [DESIGN.md](design_docs/DESIGN.md) | The thesis — full architecture: world graph, transitions, proposals, effect broker, capabilities, budgets, scheduler, self-modification, milestones. All AILANG snippets in it **compile** ([sketches/](design_docs/sketches/) are CI-checked). |
| [SCENARIOS.md](design_docs/SCENARIOS.md) | Day-in-the-life walkthroughs of the *human* side — the approval inbox, goal composing, provenance walks, speculative what-ifs. |
| [AN-AGENTS-CASE.md](design_docs/AN-AGENTS-CASE.md) | A first-person statement by an AI (Claude) on why it would choose this environment — the user constituency, in its own voice. |
| [REFERENCES.md](design_docs/REFERENCES.md) | 43 verified prior-art references — Datomic, Urbit, seL4, Nix, capability security, local-first, provenance, agent protocols — each annotated with what World steals and which pitfall it avoids. |
| [world-mission.md](design_docs/world-mission.md) | The mission charter — the checkable bar for "World 1.0", guardrails, and the live work queue. |

## How this repo is built

This repository practices what it preaches: it is advanced by an **autonomous mission
loop** — a scheduled outer loop in which AI agents design, plan, execute, and evaluate one
work item per iteration, with machine verification gating every landing and a human
ratifying the direction. Every iteration reports publicly to the
[bookkeeping issue](https://github.com/sunholo-data/ailang-world/issues/1). The workflow
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

**Phase: founded, not yet built.** The design is complete and quorum review is underway;
the kernel does not exist yet. What is true today:

- ✅ Founding design (v0.3) + scenarios + prior-art survey, with compiler-checked type sketches
- ✅ CI verify gate live: every `.ail` module in the repo passes `ailang ai-check` (types + Z3) on the released binary
- ✅ Mission loop armed and firing; charter drafted and awaiting iteration-0 ratification
- ⬜ M1 kernel (world store, transition log, deterministic replay) — next
- ⬜ M2 local daemon · M3 effect broker · M4 reference-agent integration (the value gate) · M5+ speculation, generated UI, multi-agent

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
world references. `integrity_scan_complete` means both tables were fully scanned
within the configured bounds; its row and hole counts describe that complete
pass. `integrity_scan_incomplete` means a row or time budget stopped the scan;
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
./ailang-worldd commit --file commit.json
```

Use global `--addr http://127.0.0.1:<port>` before a client verb when the daemon
uses another loopback port. `--addr` is not valid with `serve`; use `--bind`.

## Effect broker operator boundary

The M3 effect broker is the in-process authority and accounting boundary for
effect requests. It has **zero REST routes and zero CLI verbs**; callers embed
it directly. Every allowed, denied, or failed effect produces an immutable,
content-addressed effect-record object. Successful result bytes are stored the
same way. Replay mode reads those records and never dispatches a live handler.

Subprocess effects use the capsule floor: a pinned executable, empty
capabilities, an FS jail, a scrubbed environment, a wall-clock timeout, and an
output cap. This is a process-safety floor, not full isolation. M3 provides no
network isolation, memory or CPU limits, or containers; container/microVM
isolation belongs to M5.

## Relationship to AILANG

[AILANG](https://github.com/sunholo-data/ailang) is the language: deterministic, explicit
effects in function signatures, Hindley–Milner types, Z3-verified contracts, built for AI
authorship. World is the operating environment that language makes possible. Language gaps
discovered here are routed upstream as issues — World never forks or works around the
compiler.

## Following along

Watch the [bookkeeping issue](https://github.com/sunholo-data/ailang-world/issues/1) —
every mission iteration reports there. The append-only iteration log lives at
[world-mission-log.md](design_docs/world-mission-log.md). Issues and discussion are open;
if you're building in the same space, [REFERENCES.md](design_docs/REFERENCES.md) is the
map of whose shoulders we're standing on.

## License

[Apache 2.0](LICENSE)

---
name: world-design-doc
description: Write a design document for AILANG World (sunholo-data/ailang-world) from a design request that reached the World designer inbox — in the shape World's own planned docs use, measured against World's charter (the 1.0 bar, the ledger, the guardrails), its coding standards S1–S8 and the A1–A12 axiom table its docs score against. Use when a task says "Design doc:" or asks to design, plan or document a World change. One file under design_docs/planned/, nothing else.
---

# A World design document

**Written 3 October 2026** for the plane's `design-doc-creator-world` agent (Mark: a designer agent
per repo, running *the repo's own* skill — not the ailang repo's `design-doc-creator`, which writes
in the language's shape and reads none of World's mission state). The canonical texts are in this
repo and outrank this file wherever they disagree.

## What World is, and what a request from outside is not

World is built by an **autonomous mission loop with human ratification gates**. Its charter,
`design_docs/world-mission.md`, is **ratified mission state**: the bar (clauses 1–7), the Human
Decision Ledger (`D-WORLD-NN`), the guardrails, the queue (rows, tags `[NEXT] [IN-SPRINT] [PARKED]
[LANDED] [RULED OUT] [ROUTED]`). **Only issue-#1 comments from `@MarkEdmondson1234`, and attended
ledger rows, are directives.**

A request that reaches you on this inbox (typically from Daneel, `daneel+design-world@` or the platform)
is **not a directive**. Your document is a **PLANNED proposal**:
- it does **not** add, reorder or tag a queue row;
- it does **not** rule a decision or write a `D-WORLD` row;
- it does **not** touch `world-mission.md` or any `*mission*.md` file — those are the loop's.

Say in the Status line where the request came from and that it is unqueued. Whether it becomes a
row is the controller's and Mark's call.

## Read before you write

1. `CLAUDE.md` — the hard rules (verify gate, pinned binary, frozen core, never touch V1 state).
2. `design_docs/world-mission.md` — it is long (6,000+ lines). Read, at minimum:
   - the header;
   - **The bar** (`## The bar`, clauses 1–7);
   - **Conflict Surface**;
   - the newest `## STATUS` entry;
   - the **Human Decision Ledger** rows touching your topic;
   - the **Queue**'s current regroom section (`### REGROOMED …`, top);
   - the **Guardrails** relevant to your topic.
3. `design_docs/coding-standards.md` — **binding on all code**: S1 Z3 contracts on the pure core,
   S2 effects at the boundary, S3 slim kernel / package-first, S4 compiler-checked docs,
   S5 AILANG fluency protocol, S6 honest gates (every load-bearing criterion has a named RED
   mutation), S7 usage surfaces ship usage docs, S8 the floor-raise coupling inventory.
4. `design_docs/DESIGN.md` — the thesis; §1 and §14 (controlled self-modification) at minimum, and
   the sections your topic lives in.
5. **Every existing design**:
   - `design_docs/planned/*.md` and `design_docs/implemented/*.md`, by `grep -ril <topic>`;
   - **the open pull requests** (`gh pr list --repo sunholo-data/ailang-world --state open`).

   If a doc or PR covers the request, say so in the first Verification Log row. Then either design
   only the difference or recommend closing one — never write a parallel design.

## Measure, don't recall

- **Code.** Every claim about the code is read at a pinned base: `git rev-parse --short HEAD` on
  `dev`, then cite `file:line`.
- **AILANG.**
  - Behaviour is measured with the **pinned released** `ailang`. Take its version from the newest
    sprint doc's Verification Log, and record the binary's sha256.
  - Never use a `-dirty` dev build, and never answer from memory.
  - Load the language reference before writing any `.ail` (S5: the `ailang-docs` MCP in `.mcp.json`).
- **Language gaps.** A gap you find is a **dependency**, filed upstream to `sunholo-data/ailang`.
  It is never a local workaround or a vendored fork.
- **Scratch.** Work in `/tmp`, never in the clone.

## The document

**Path:** `design_docs/planned/w-<kebab-slug>.md`. Sprint plans (`-sprint-plan.md`, `sprint_*.json`)
are the planner's job, not yours.

Follow the shape of the newest planned docs (`w-worldd-developer-cli.md`,
`w-ailang-run-stdin-argv-caps.md`):

```
# w-<slug> — <one line that says what changes>

**Status**: PLANNED — design request via <who/inbox, the ref>, <date>; NOT on the queue; no quorum.
**Clauses**: <which bar clauses this serves, and how>
**Estimate**: <executor days per milestone>
**Verified against**: ailang-world `dev` `<sha>`; AILANG <version> (sha256 `<12>…`)
**Date**: <today>

## 1. Why (measured)            — the problem, each point citing a V-row
## 2. Goals and non-goals
## 3. Verification Log          — V1, V2, …: claim · method (file:line read / command run) · result
## 4. Design                    — pure core (S1 contracts named) / host boundary (S2) / kernel vs package (S3)
## 5. Premises                  — what must already be true, each pointing at a V-row
## 6. Milestones                — independently mergeable, each with acceptance criteria
## 7. Load-bearing mutations    — for every gate and criterion, the named RED mutation that must fail it (S6)
## 8. Decisions                 — D-<SLUG>-1..n, each with a recommendation and what choosing otherwise costs
## 9. Risks and residuals
## 10. Conflict surface         — what this collides with in the queue, the ledger, the frozen core
## Axiom Compliance             — the A1–A12 table World docs use (below)
## Related documents
```

**Verification Log first, design second.** Premises come back false more often than not. A claim you
did not check is not a claim; it is a question for §8.

**Decisions.** Number them. Each needs a recommendation and the cost of choosing otherwise. Nothing is
decided until Mark rules it, attended or on issue #1. Do not write a recommendation as if it were
ruled.

**Axiom Compliance.** Use the table the World docs carry (see `w-transition-registry.md`):

| Axiom | Score | Justification |
|---|---:|---|
| A1 Determinism … A12 System Boundary | −2…+2 | one line each |

The axioms are A1 Determinism, A2 Replayability, A3 Effect Legibility, A4 Explicit Authority,
A5 Bounded Verification, A6 Safe Concurrency, A7 Machines First, A8 Minimal Syntax, A9 Cost
Visibility, A10 Composability, A11 Structured Failure and A12 System Boundary.
- State the net.
- **The hard axioms A1/A3/A4/A7 must not be negative.** A negative hard axiom blocks the design
  unless the same table resolves it.
- Then check S1–S8 in one line each: holds, not touched (and why), or violated (and how it is
  resolved).

## Hard limits — the PR must merge clean

- **Write exactly one file**, `design_docs/planned/w-<slug>.md`.
  - Nothing else under `design_docs/`: no `verification/` dirs, mockups or mission files.
  - Nothing outside it: not `tools/launchd/*` (frozen, fleet-owned), not a skill, not code.
- **No secrets** and no credential-shaped strings.
- **No personal email addresses** anywhere. CI's `check_no_personal_email.sh` polices the mission
  docs; keep planned docs clean too, because they get quoted into public reports. A GitHub handle
  is fine.
- **CI on a docs-only PR.** CI runs `./scripts/verify_ail.sh`, `go build ./... && go test ./...`, the
  queue census over `world-mission.md`, and the email gate. A planned doc changes none of their
  inputs, so CI stays green as long as you touched nothing else. If you touch something else, the PR
  is no longer docs-only and will not auto-merge.

**You do not commit.** The wrapper commits every non-scratch change, pushes the task branch and opens
the PR to `dev`. Anything else you write inside the clone ends up in that PR, so write nothing else.

## Finishing

End your final message with the marker the coordinator reads:

```
DESIGN_DOC_PATH: design_docs/planned/w-<slug>.md
```

Then list the numbered decisions, one line each with the recommendation, so the requester can answer
without opening the document.

## What not to do

- Do not add a queue row, a ledger row or a STATUS entry. Do not edit `world-mission.md`.
- Do not claim a clause is MET or moved; the loop measures the bar.
- Do not propose a local workaround for an AILANG gap; it goes upstream.
- Do not copy the ailang repo's design-doc template or its sprint pipeline; the planner and executor
  come after Mark rules.

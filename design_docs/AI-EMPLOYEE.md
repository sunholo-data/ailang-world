# The AI Employee — World's first external application (v0.1)

**Status**: DRAFT, use-case statement + gap inventory. Authored attended (Mark + Claude Fable 5.1,
2026-09-01). **Not a queue row, not quorum-reviewed, not ratified** — it proposes rows (§8) and a
directive (§9) for Mark to place on issue #1; until then nothing here steers the loop. Verified
against `dev` at `e9a0efa` (2026-09-01, iteration 146). Every "exists / does not exist" claim
below carries a file:line or a charter row; a claim without one is a claim, not a fact.

Companion to [SCENARIOS.md](SCENARIOS.md) (the human side), [AN-AGENTS-CASE.md](AN-AGENTS-CASE.md)
(the agent side), [HUMAN-SURFACE.md](HUMAN-SURFACE.md) (the ratified interaction grammar) and
[POSITIONING.md](POSITIONING.md) (the wedge this application sells). Where this document conflicts
with [DESIGN.md](DESIGN.md) or the charter, they govern.

---

## 1. The one-paragraph answer

An **AI employee** is a long-running engineering agent that holds its own accounts — email,
GitHub, chat — does the delegable work of a remote engineering role (inbox, tickets, code, PRs,
reviews, docs, the standup note), and reports to a **human owner** who is the accountable party.
World is the layer that makes that arrangement governable rather than merely possible: the agent
*proposes*, a deterministic verifier *checks*, only authorised and budgeted proposals *commit*,
every action leaves a receipt, and the owner's attention is a metered resource the agent spends
against a budget. The owner's job is exactly the three verbs the README already gives humans —
**express a goal with a budget, decide what verification couldn't settle, ask the world why** —
which is why this is World's natural first application rather than a new product: the mission
loop is customer #0 (DESIGN.md §12.1), the AI employee is customer #1.

## 2. What an AI employee is — and what it is not

**It is** an agent with:

- a **declared identity** — its own email, GitHub, and chat accounts, a bot badge, signed commits,
  `Co-Authored-By` on every change; a **session credential** in World that binds it to an episode
  and an explicit set of capability grants;
- a **capability envelope** — typed grants with scope, expiry and budget (`world/types.ail:13-18`;
  `host/broker/decide.go:15`), so "may open a PR on repo X but never push to `main`" is a type the
  broker enforces, not a rule the agent remembers;
- an **interrupt budget** — `Human.Approve` is a budgeted effect (`host/broker/approve.go:17`),
  so "how many times per day may this agent ask its owner" is scheduled, not hoped for
  (DESIGN.md §8);
- a **ledger** — every effect it performs is a content-addressed record with a receipt; every
  approval is spent exactly once (`host/store/journal.go:573,626`, `approval_claims` PRIMARY KEY);
  "why did it send that email" is a provenance walk, not an interview.

**It is not** a disguised human. The identity is declared, never impersonated — this is a design
axiom, not a preference, for three reasons: (1) the buyer of a governed agent is buying the
audit trail, and an audit trail on an undeclared actor is worthless to them; (2) EU AI Act
Article 50 transparency obligations (people must be told they are interacting with an AI) make
the disguised form unsaleable in the market POSITIONING.md already targets under Article 14;
(3) World's own guarantee 4 (DESIGN.md §13.2, "judged by evidence, not identity") only works when
identity is honest. An AI employee that hides what it is forfeits the one thing World gives it.

## 3. Why this is World's first application

| World already says | The AI employee is that, applied |
|---|---|
| "Humans don't operate the system; they govern it" (README) | The owner is a manager of agents: goals, approvals, why-questions — not a co-worker pretending the agent is one |
| Scenario 1, the morning approval inbox: 17 transitions, 14 auto-committed, 3 need you, 4 minutes (SCENARIOS.md) | The owner's morning with an AI employee, verbatim |
| "Decisions as governed objects … *how many times may this workflow interrupt a person* … the one sellable soonest" (POSITIONING.md §"actual position", item 3) | The product's headline feature: an agent with a metered claim on its owner's attention |
| Use case 5, regulated / high-stakes operation: "who authorized what, on what evidence, *is* the product" (DESIGN.md §12.1) | The client-facing compliance story for a delegated engineer |
| Use case 2, heterogeneous trust: a local model holds `FS.Read` on one worktree, a cloud executor holds `Git.Commit` on one repo | The employee's own fleet: cheap triage model, capable coding model, one envelope each |
| Calibration records earn wider grants by track record (DESIGN.md §13.3) | The employee's probation → trusted progression, as a typed object |
| Clause 5's value bar: ≥3 REAL "why did X happen" questions answered in ≤5 min by provenance walk | A manager asking "why did it merge that / email them" is the real question the charter is waiting for; the loop's own CI archaeology is not |
| Clause 4's floor: a reference agent, shell arm vs World arm, non-inferiority | The employee's real task set is a better benchmark than a synthetic tier, and it is the arm a customer would actually run |

The last two rows are the argument. The charter measures World's *value* against the operational
status quo on real questions, and its *floor* on real agents. An AI employee doing real work for a
real owner generates both kinds of evidence as a by-product of existing; the mission loop, by
construction, generates only evidence about itself.

## 4. Scenario 7 — a day with a resident employee

*(In the style of SCENARIOS.md; the human is the owner, the agent is a World resident. Every
step names the primitive it rides on and its status in §6.)*

**08:10 — Intake.** Overnight the employee read its inbox (`Email.Read`, §6 row E1 — NOT BUILT)
and the tracker. It triaged 23 messages into 4 tasks and 19 no-actions; the triage is a
transition with `Model.Infer` evidence attached (`host/broker/handlers_model.go:12`, BUILT). No
human was interrupted: nothing in triage requires authority the agent does not hold.

**09:30 — A task becomes a proposal.** "Fix the flaky retry test in `docparse`." The agent forks a
worktree (`FS.*`, BUILT), edits, runs the suite, and proposes a commit under its `Git.Commit`
grant scoped to a branch (`host/broker/handlers_git.go:10`, BUILT). The verifier checks types
and contracts before anything fires (`world/contracts.ail:71 commitAllowed`, BUILT). The PR opens
under the agent's own GitHub identity (`GitHub.Write`, row E2 — NOT BUILT).

**11:00 — The gate verification cannot settle.** The fix needs a reply to the client who reported
it. Email to an external party is the "External class always asks" case (SCENARIOS.md scenario 1
item 3). The broker mints a `Human.Approve` request (BUILT), debits the agent's interrupt budget,
and a **decision packet** is created — proposal hash, evidence bundle, the exact authority asked,
a deadline, and a timeout policy (`world/types.ail:159-167`, BUILT as types; Z3-proven timeout
law `timeoutOutcome`, `world/types.ail:172`). The packet does not ping the owner: it is batched
with two others for the owner's next digest (HUMAN-SURFACE P5; the inbox itself is row 7 — NOT
BUILT).

**13:45 — The owner decides.** Three packets, one glance each: approve the client email
(attenuated — "send, but cc me"), approve the PR merge, defer the dependency bump for more
evidence with a new bounded deadline (`validDefer`, `world/types.ail:188`, BUILT). Every decision
is itself a transition; there is no admin backdoor (HUMAN-SURFACE §5).

**17:00 — The standup.** The agent's end-of-day note is not prose it wrote from memory; it is a
projection over the ledger — every noun a typed link with its grade (`PROVEN / TESTED / ATTESTED /
CLAIMED`, `world/types.ail:42 gradeOf`, BUILT; rendered by the read-only workbench,
`host/daemon/daemon.go:565`, BUILT). Ungrounded sentences show as ungrounded (P2).

**Two weeks later — "Why did the client get two emails?"** The owner selects the second email's
effect record → provenance → the proposal → the deferred packet whose escalation re-fired after
the first reply landed. Ninety seconds. That is clause 5's ≥3-real-questions evidence, arising
from actual operation.

## 5. What the human is still for — as World objects

The question that motivates this document is "what is a human still needed for in a remote
engineering job once an agent has the accounts?" World's answer is precise, and each item is (or
will be) a typed object rather than a convention:

| Human function | World object | Status |
|---|---|---|
| **Accountability** — the name on the contract, the one who can be fired | Goal owner on every `Proposal.goal`; `HumanApproval(HashRef)` evidence names who decided | BUILT (types) |
| **Authority minting** — who is allowed to give the agent its powers | Session credential → (episode, grants) resolver | NOT BUILT — row 39 |
| **Judgement verification cannot settle** — external comms, deploys, spend, capability widening | `Human.Approve` gate on the External / Irreversible / Capability classes | BUILT (effect); policy per class NOT BUILT |
| **Bounded attention** — the owner is not on call for the agent | Interrupt budget on `Human.Approve`; batching; TTL + timeout policy so silence never blocks | BUILT (budget, laws); inbox NOT BUILT — row 7 |
| **Trust calibration** — probation, then wider scope | Calibration record (DESIGN.md §13.3) | NOT BUILT — no code references outside tests |
| **Relationships, meetings, politics** | Deliberately outside World — the agent prepares the owner (a "questions for Thursday" projection) and follows up; it does not attend as a face | Not a World concern; employee-package scope (§7 B) |

## 6. Premise verification — what exists at `e9a0efa`

Status vocabulary as in HUMAN-SURFACE.md §8: **BUILT / PARTIAL / NOT BUILT / BLOCKED /
UNVERIFIED**. Rows E1–E3 are new effect handlers this application needs and are numbered so §8
can name them.

| # | Premise | Status | Where | Evidence |
|---|---|---|---|---|
| 1 | Immutable content-addressed log with bit-for-bit replay | BUILT | `host/store`, `host/replay` | Clause 1 landed (charter rows 2, 16); `scripts/verify_go.sh` |
| 2 | Broker: capability grants, budgets, effect records, undeclared-effect refusal | BUILT | `host/broker/broker.go:87 NewSession`, `decide.go:15-33`, `Session.Bind` | Charter row 4 (`implemented/w-effect-broker-m3.md`); 7/7 Z3-proven decision law |
| 3 | `Human.Approve` as a budgeted, spend-once effect | BUILT | `host/broker/approve.go:17`; `host/store/journal.go:573,626` | `approval_claims` PRIMARY KEY; budget debited at request time |
| 4 | Decision packet with deadline, escalation budget, typed timeout policy | BUILT as types + laws; **not invoked by any host path** | `world/types.ail:138-220` | HUMAN-SURFACE §7.3: "five laws are proven and, as landed, not yet invoked by any host path" |
| 5 | Total evidence-grade mapping | BUILT | `world/types.ail:42 gradeOf` | Charter row 15; `PROVEN` still has no minting constructor (`types.ail:41`) |
| 6 | Read-only workbench rendering grades and provenance edges | BUILT (read-only) | `host/daemon/daemon.go:565`, `host/workbench/render.go` | Charter row 14 complete at 11/11 |
| 7 | Effects available to a resident | PARTIAL | `handlers_fs.go:12-13`, `handlers_git.go:10`, `handlers_model.go:12`, `approve.go:17-18`, `registry_publish.go:25` | Exactly seven: `FS.Read`, `FS.Write`, `Git.Commit`, `Model.Infer`, `Human.Approve`, `Human.PollApproval`, `Registry.Publish`. None reach email, GitHub, calendar, or the network |
| 8 | A way for an external agent to use the broker | **NONE BY DESIGN at M3** | README §"Effect broker operator boundary" | "zero REST routes and zero CLI verbs; callers embed it directly" |
| 9 | MCP projection of the transition registry (Tier-0 citizenship, DESIGN.md §13.1) | BLOCKED (upstream) | `planned/w-mcp-dispatch-projection.md` | `sunholo-data/ailang#885` open; `serveapi/protocol` has no JSON-RPC dispatch |
| 10 | A2A agent card + session-scoped projection | BLOCKED (local) | `planned/w-a2a-session-projection.md`, charter row 40 | Blocked on row 39 |
| 11 | Session authority: who mints a credential, credential → (episode, grants), expiry, fail-closed | NOT BUILT | Charter row 39 | `NewSession` takes `grants` as an argument; `Bearer` 0 / session-lookup 0 across `host/` (row 39, controller-measured) |
| 12 | Approval inbox (P1 + P3 + P5) | NOT BUILT | Charter row 7; `mockups/approval-inbox.html` | Parked behind item 5's `P6.B`, which was split into the two BLOCKED projections in table rows 9–10 |
| 13 | M4 reference-agent floor run | NOT RUN | Charter row 6 | "PARKED until 2–5 land" — 2–5 all landed by iter-127; row 6 is unparked and unpicked |
| 14 | Calibration record per agent (§13.3) | NOT BUILT | — | `rg -il 'calibration' host/ world/` → test files only |
| 15 | First-contact protocol: identity, grants, transitions, budgets on connect (§13.1) | NOT BUILT | — | Depends on rows 9/10 and 11 |
| 16 | Fork / compare-arms (M5) | NOT BUILT | `host/store` has one selected head | HUMAN-SURFACE §8 |
| E1 | `Email.Read` / `Email.Send` handlers | NOT BUILT | would follow `handlers_git.go`'s pattern behind `broker.Registry` (`broker.go:35`) | eparse is the obvious read-side donor; send is External class, always gated |
| E2 | `GitHub.Read` / `GitHub.Write` handlers | NOT BUILT | same | Named in DESIGN.md §8's effect list; no handler exists |
| E3 | `Doc.Parse` handler | NOT BUILT | same | docparse as a capsule-floor subprocess (`host/capsule`) |

Rows 8–12 are the whole story: the kernel and broker are real, and **nothing outside a Go
process can be a resident yet**. The floor (row 13) has never been run because there has never
been anything to run it on.

## 7. What we need to do

Two tracks, deliberately separated. Track A is kernel work and goes through the mission loop with
its full discipline. Track B is the application and **must not** go through the loop — it is a
consumer of World, developed package-first (coding-standards S3) in its own repository, routing
kernel gaps here as rows exactly as World routes language gaps to `sunholo-data/ailang`. Putting
the employee's feature list through quorum would be the slowest possible way to build it and
would grow the kernel it is supposed to sit on top of.

### Track A — World rows, in the order that unblocks a resident soonest

| Step | Row | What it gives the employee | Estimate (charter's where one exists) |
|---|---|---|---|
| A1 | **39 `w-session-authority`** | The identity primitive: a credential that resolves, fail-closed, to (episode, grants) | ~0.5–0.8d, needs a design doc, gated on nothing |
| A2 | **40 `w-a2a-session-projection`** | The employee connects as a peer with an agent card, acts through propose → verify → commit | ~0.7d after A1 |
| A2′ | *(alternative while A2 or `ailang#885` is open)* an **embedded resident**: a `cmd/ailang-worldres` binary that runs the agent loop in-process with the broker, the way `ailang-worldd` embeds it today | Unblocks Track B immediately. Honest cost: it is not the "World-MCP arm" clause 4 names, so it does not count toward the floor; it counts toward clause 5 | unmeasured; the daemon is the pattern |
| A3 | **7 `w-approval-inbox`** | The owner's surface — packets, grades, five verbs, digest batching; invokes the four proven laws for the first time | ~2d, built to HUMAN-SURFACE (ratified) |
| A4 | **new: `w-external-effect-handlers`** (§8 R1) | E1–E3 above, behind the existing `Registry`; the External-class always-ask policy as a typed rule | unmeasured — three handler pairs following `handlers_git.go`; the design doc prices it |
| A5 | **6 `w-agent-floor-m4`** | Run the floor with the employee's real task set as (or alongside) the standard tier | ~3d |
| A6 | *(post-1.0 lane, pull forward on demand)* `w-calibration-record` (§13.3) and `w-decision-budget` (the POSITIONING wedge: per-owner, per-day interrupt budget as a world object) | The probation → trusted dial and the sellable feature | unpriced |

The reorder this implies is the substance of §9: the dashboard's "Next picks" at iteration 146
are rows 52, 54, 55 and 56–61 — all clause-2 instrumentation hygiene (CI step scoping, the
driver-copy drift gate, the lever parser, the canary fence, the approvals spine, a
`verify_go.sh` flake, grep-cannot-prove-live, an inert-rename false red, a one-inserted-line
fail-open) — and only "then row 39". Rows 48, 49 and 51 of the same class landed and row 53 was
routed upstream across iterations 138–146, while rows 6, 7, 39 and 40 did not move. They are
real defects the loop found in its own gates; none of them is on the path to a resident.

### Track B — the employee package (outside this repo)

| Piece | What it is | Rides on |
|---|---|---|
| B1 Provisioning kit | One command creates the agent's email / GitHub / chat identities with the bot flag, signed-commit key, and naming convention; writes the World session grant request | A1 |
| B2 Resident loop | The agent's outer loop: intake → task → proposal → (verify/commit by World) → report; long-running on the studio host | A2 or A2′ |
| B3 Work journal & standup | A projection over the ledger: grounded prose (P2), grades visible (P3), one digest per day (P5) | rows 6, 5 of §6 |
| B4 Owner digest & mobile card | The batched decision packets, rendered small (SCENARIOS.md scenario 4) | A3 |
| B5 Meeting proxy | Calendar-aware; produces "questions I need you to answer before Thursday"; ingests transcripts as `RecordedEffect` evidence | E-row handlers |
| B6 Metering | Tokens, wall-clock, cost per task, approvals spent vs budget — the invoice line items | broker effect records (BUILT) |

Nothing in Track B is a kernel change. If any piece turns out to need one, that is a row here,
not a fork.

## 8. Proposed queue rows (for Mark to place; text is a draft, not a directive)

**R1 — `w-external-effect-handlers` · clause-3 ·** THE BROKER'S SEVEN EFFECTS REACH THE
FILESYSTEM, GIT, ONE MODEL CALL, THE HUMAN AND THE REGISTRY, AND NOTHING ELSE A RESIDENT DOING
REAL WORK TOUCHES. Add `Email.Read`, `Email.Send`, `GitHub.Read`, `GitHub.Write`, `Doc.Parse` as
extension handlers behind `broker.Registry` (`host/broker/broker.go:35`), each under the capsule
floor where it spawns a subprocess, each producing effect records like the existing seven. Carry
the External-class policy from SCENARIOS.md scenario 1 ("Email/External always asks, regardless
of evidence") as a typed rule with a Z3 contract, not a handler-local convention. DESIGN.md §8
already lists `GitHub.Read · GitHub.Write · Email.Send` as kernel-understood effects; this row
makes that list true. · NEEDS A DESIGN DOC · gated on nothing (handlers), on row 39 for the
grants that would authorise them in a real session · unpriced until designed.

**R2 — `w-resident-employee-floor-arm` · clause-4 ·** ROW 6'S FLOOR HAS NEVER RUN BECAUSE
NOTHING HAS EVER BEEN A RESIDENT. When row 6 is picked, add the AI-employee task family (intake
→ PR → gated external reply) as a recorded per-task-family split alongside the standard tier, per
row 6's own 2026-07-31 design-lens note (record splits, not only aggregate pass-rate).
Thresholds UNCHANGED. · amends row 6's scope, not its bar · gated on row 6.

## 9. Proposed directive (issue #1, Mark's account only)

> Reorder: take row 39, then row 40, then row 7, then row 6, ahead of rows 52, 54, 55 and 56–61.
> Batch those hygiene rows into one sprint after row 7 lands. Add rows R1 and R2 from
> `design_docs/AI-EMPLOYEE.md` §8 to the queue after row 40.

(The one decision parked on Mark at iteration 146, `D-WORLD-31`, concerns row 50's canary rule and
is independent of this document; it is not addressed here.)

## 10. Conflict surface — what this document does not decide

- **herdr vs World.** The studio-host session manager (`herdr`) keeps agent *processes* alive
  across laptop disconnects; World governs their *effects*. Different layers, same ruling as the
  Motoko DST boundary (DESIGN.md §16): herdr hosts the resident loop, World brokers what it may
  do. Neither replaces the other; this document does not propose merging them.
- **eparse / docparse as tools vs as handlers.** Today they are CLIs the agent could run inside
  the capsule under `FS.*`. Wrapping them as effect handlers (E1, E3) is what makes their results
  `RecordedEffect` evidence with receipts instead of bytes in a transcript. The row's design doc
  decides whether the handler shells out to the CLI or links the library.
- **Collaboration Hub.** Still a pattern donor, not a code donor, per HUMAN-SURFACE §6.1; the
  owner digest (B4) reuses its approval-queue interaction, not its schema.
- **The floor.** Nothing here relaxes clause 4. If the employee task family makes the World arm
  look worse than the shell arm on eligible agents, that is the floor doing its job.
- **Commercial framing.** "An engineering agent with its own accounts, a verified permission
  policy, an audited ledger, and a human owner who is accountable" — never "an AI employee you
  cannot tell from a human." The second framing is out of scope of this document and of World.

## 11. What is deliberately not claimed

No number in this document is a measurement of value: the ≥3-questions bar and the
non-inferiority thresholds are the charter's and are unchanged; whether the employee task family
clears them is what running rows 6 and R2 will tell us. The estimates in §7 are the charter's own
where a row exists and are marked unmeasured where it does not. The scenario in §4 is a design
fixture in the SCENARIOS.md sense — it names which primitive each step rides on so the reader can
check §6 rather than trust the prose.

---
*v0.1 drafted attended 2026-09-01 (Mark + Claude Fable 5.1). Next step is a human one: place §9
on issue #1, or not.*

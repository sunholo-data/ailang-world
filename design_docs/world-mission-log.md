# Ailang World Mission — iteration log

Append-only. One entry per outer-loop iteration. Newest at the bottom.

---

## Iteration 0 — 2026-07-23 — charter advisory quorum; ratification PARKED for Mark

**Kind**: iteration-0 bootstrap prep (no sprint routed — the queue is un-ratified; the charter
forbids routing any sprint before ratification).

**Context / preflight**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (loop armed). Billing tripwire:
  CLEAN (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff`.
- Pidfile `mission-world.pid`=34197 = this run's own driver (no overlap). Local `dev` ==
  `origin/dev` == `3c8791d` (in sync). CI `CI` on dev: **completed/success** @ `3c8791df5`.
- Inbox: 1 unread — a generic `eval-suite` partial (42/50, 84%) on controlplane; not a World
  regression, not a directive. No `[nightly-eval]` issues. Bookkeeping issue `#1` open; **zero**
  `@MarkEdmondson1234` comments (watermark unchanged).

**Pick**: iteration 0 = ratify the charter (bar + queue + guardrails) through the design quorum,
**with Mark, attended**. Ratifying one's own operating charter is a human-authority gate (kill
switch + directive principal = Mark only) — a headless agent must not self-ratify. So the
ratification decision is **PARKED (needs-human-review)** for Mark; this iteration produced the
machine-half prep to make his attended session short.

**Work done (headless-safe, no authority overreach)**
1. **Advisory design quorum on the charter** (`ailang design-quorum design_docs/world-mission.md`,
   controller-verdict `reject`). Both reviewers present; **BLOCKED** (3/3 reject). Artifact:
   `.ailang/state/mission-quorum/world-mission-2026-07-23T18-09-28Z.json`. Metered cost **$0.037**.
   - `gpt5-6-sol` (reject, $0.0272): charter relies on **unverified premises** (ai-check JSON,
     serve-api MCP/A2A, CI job names, mission-control profile vars, launchd/kill-switch). Wants a
     "Premise Verification Log."
   - `gemini-3-1-pro` (reject, $0.0097): fails the **Conflict Surface** gate — proposes a new
     `ailang-worldd` Go daemon (SQLite/REST/CLI) with no overlap analysis vs existing
     `ailang serve-api`. Wants a `## Conflict Surface` section before ratification.
   - controller (me, reject): **clause-4 M4 thresholds are explicitly un-fixed** ("candidate …
     FIX NUMERICALLY at iteration 0") — a gate with un-fixed numbers is not ratifiable; and
     ratification is Mark's attended authority gate.
2. **Premise verification** (read-only, de-risking gpt5-6-sol's objection):
   - ✅ `serve-api --mcp / --mcp-http / --a2a` exist (clause-6 machinery real).
   - ✅ `scripts/verify_ail.sh`; CI `CI` / job `ailang-code verify gate` (green); `DESIGN.md`;
     sketches `transitions.ail`, `worldtypes.ail` all present.
   - ✅ `ai-check` exists (types + Z3 verify; flags `-relax-modules -timeout
     -verify-recursive-depth`).
   - ❌ **`ai-check` has NO `--json` flag** in v0.30.0. The charter's Repo Profile AND the shared
     `mission-control` `ailang-code` verify-profile both cite `ailang ai-check --json` as the
     unified gate — the shipped binary does not provide it. Real premise defect.

**Routing evidence**
| Role | Pinned | Actual | Notes |
|---|---|---|---|
| Controller (triage/pick/quorum-verdict/record/retro) | `$MODEL` session | claude-opus-4-8 | opus-first, correct |
| Design/plan/execute/evaluate | — | not spawned | no sprint routed (queue un-ratified) |
| Quorum reviewers | gpt5-6-sol, gemini-3-1-pro | both present | metered $0.037 total (cap $0.10/reviewer; budget $5) |

**Metered ledger**: `metered=$0.037` (quorum only). Ceiling `$5` not approached.

**Ruled out** (do not re-chase)
- Self-ratifying the charter headless — human-authority gate; parked for Mark by design.
- Routing any queue item (incl. top `[NEXT] w-log-epoch-decision`) — charter forbids sprint
  routing before ratification.

**Parked for human (Mark, attended ratification session — agenda)**
1. **Fix clause-4 numerically** — M4 value gate pass-rate + wall-clock-overhead thresholds
   (candidates: within-noise pass-rate, ≤25% overhead). Cannot be argued after the fact.
2. **Add a Conflict Surface section** — `ailang-worldd` vs existing `ailang serve-api`/state
   (gemini objection). Justify a new daemon+repo vs extending the frozen core.
3. **Resolve the `ai-check --json` premise defect** — either correct the charter/profile to
   `ailang ai-check` (which already does check+verify), or file an upstream request on
   `sunholo-data/ailang` to add `--json` output (language/tooling gap → routes upstream per
   guardrail). Recommend the former (drop `--json`) as the cheaper fix.
4. (optional) add the "Premise Verification Log" appendix gpt5-6-sol asked for — most rows
   already verified above.

**Next**: on Mark's ratification (kill switch stays his), re-quorum ONCE post-revision, then pick
`[NEXT] w-log-epoch-decision` (clause-1, ~0.5d design doc).

**Retro** (Gate 5)
- Friction A (instance 1/2): charter + shared-skill `ailang-code` profile both cite
  `ailang ai-check --json`, which v0.30.0 does not provide. Skill edits route via the shared
  Gate-5 retro and need ≥2 instances — recorded as a watch-item, not edited this iteration.
- Process note: the loop was armed (kill switch off) while the charter says it "stays SET until
  ratification." Surfaced to Mark — arming is his call; ratification remains pending.
- No skill edit, no mission-doc bar/queue edit (ratification-scope = Mark's). STATUS stamp only.

---

## Iteration 1 — 2026-07-23 — w-log-epoch-decision: D2/D3 settled, D1 replay-pin PARKED for Mark

**Kind**: first post-ratification sprint (NEW-DOC decision doc, clause-1). One backlog item.

**Context / preflight (Gate 0–1)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed, post-ratification). Billing
  tripwire: CLEAN (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`; `~/.zshenv` unset-after-source
  holds — login shell has no key). gh account: `sunholo-voight-kampff`.
- Pidfile `mission-world.pid`=49985 = this run's own driver (my PPID; no overlap).
  Local `dev` == `origin/dev` == `08ba5ba` (in sync). CI `CI` on dev: **completed/success** @
  `08ba5ba52`. No `[nightly-eval]` issues.
- Inbox: 11 unread — 9 eval-suite controlplane noise (motoko fizzbuzz partials); 1 `mission-v1`
  ACK of World's earlier channel-test (cross-mission, informational); 1 `world-coordinator`
  self-note (the iter-0 friction observation, already a guardrail). **No directive, no
  regression, no cross-mission demand/bug** → queue pick stands.
- Bookkeeping issue #1 (`mission-world-gh-issue`): 6 comments, **zero `@MarkEdmondson1234`**
  (watermark set `2026-07-23T19:25:04Z`). Created today, <80 comments → **no weekly rotation**.

**Pick + reality-check (Gate 2)**: top `[NEXT] w-log-epoch-decision`. Fresh-fetch verified:
no existing doc (`grep -ri`), not landed (no PR/commit on origin/dev), no quorum artifact →
genuine NEW-DOC. Not a bug-claim (a decision doc), so no ghost live-repro needed.

**Route + execute (Gate 3)**: NEW-DOC → **design-doc-creator on the rotation designer**. Rotation:
world-namespaced `mission-world-designer-rotation` absent → **start at `claude:claude-fable-5`**
(Fable, subscription). Fable probe green via `claude-sub` (env-stripped, subscription-only).
Ran design + **1 revision** (both bounded ≤30m, backgrounded). Skill assumes the ailang-compiler
repo layout (no `std/VERSION`/`planned/`/`changelogs/` here) → gave the designer a self-contained
brief adapting it. Doc: `design_docs/planned/w-log-epoch-decision.md`; checkable sketch
`design_docs/sketches/logepoch.ail` (controller-verified green via `verify_ail.sh`, 3 modules).

**Quorum (pick-time, the item's mandate)**:
- **r1 BLOCKED** (both reviewers reject, convergent + correct): attestation-by-conformance-corpus
  cannot prove semantic equivalence over untested programs → "replay with any attested binary" is
  non-deterministic (violates A1 + bit-for-bit); `(transitionFn, semanticsEpoch)` cache unsound.
  My controller "pass" under-weighted this — the reviewers were right. `.ailang/state/mission-quorum/
  w-log-epoch-decision-2026-07-23T19-56-42Z.json`.
- **revision** (Fable, same lane): authoritative replay now pins the exact content-addressed
  **interpreter artifact** (`interpreter: HashRef`); corpus attestation demoted to
  non-authoritative `candidateBinaries`; cache keyed `(transitionFn, interpreter)`; missing
  artifact → explicit A11 failure. §14 boundary preserved (hashing released binary bytes =
  consuming the artifact). D2/D3 substance unchanged.
- **r2 re-quorum BLOCKED** (the one allowed re-quorum): reviewers split in OPPOSITE directions —
  gpt5-6-sol wants a tighter runtime-closure (hermeticity: libc/locale/env), gemini wants a
  looser platform-independent release-identity pin (portability: exact binary is OS/arch-locked,
  breaks multi-node). The determinism⇄portability fork. `...20-04-56Z.json`.
- Bounded protocol (**1 revision + 1 re-quorum**) exhausted → **park needs-human-review**, do not
  spin (Standing rule 2 + bounded-advisory discipline). Substantively correct: the log format is
  the frozen kernel → charter guardrail requires human ratification. Parked with a 3-option
  decision framed in the doc's "Open Decision" section for Mark.

**Landed (Gate 3b)**: doc-only commit `c3c5124` on dev; **CI `CI` completed/success** @ c3c5124
(verify_ail.sh gate, 3 modules). The *doc artifact* is landed + green; the *item* is PARKED.

**Upstream routing (guardrail, both channels)**: the determinism/portability fork depends on
whether a released `ailang` exposes a platform-INDEPENDENT semantics-version identity → **GitHub
issue `sunholo-data/ailang#471`** (Q1 semantics identity + hermeticity; Q2 stable canonical-form
serialization for D2) + **`ailang messages` to `mission-control`** (`msg_20260723_220936_124e4205`,
from `world-coordinator`). Framed as a design-input question, low priority, no local workaround.

**Routing evidence**
| Role | Pinned | Actual | Notes |
|---|---|---|---|
| Controller (triage/pick/quorum-verdict/record/retro) | `$MODEL` session | claude-opus-4-8 | opus-first, correct |
| Designer | ROTATION → `claude:claude-fable-5` | claude-fable-5 (subscription, `claude-sub`) | rotation start (world file absent); ran design + 1 revision; both bounded backgrounded; rotation state written = `claude:claude-fable-5` (next new-doc → codex) |
| Planner / Executor / Evaluator | opus / opus / sonnet | **not spawned** | no sprint-plan/execute/evaluate — a decision doc parks at the quorum gate |
| Quorum reviewers | gpt5-6-sol, gemini-3-1-pro | both present, 2 rounds | metered $0.0547 (r1) + $0.0687 (r2) |

**Metered ledger**: `metered=$0.1234` (quorum r1 $0.0547 + r2 $0.0687). Ceiling `$5` not approached.
Designer = subscription (Fable), $0 metered. Upstream routing = $0 (gh + local message).

**Ruled out** (do not re-chase)
- A THIRD quorum round / further headless revision on D1 — bounded protocol exhausted; the fork is
  a values call for the human authority gate, not a headless-provable answer. Spinning would
  violate Standing rule 2 + the bounded-advisory discipline.
- Picking item 2 (`w-world-library-m1`) this iteration — it is BLOCKED on D1 (M1 freezes the log
  format the parked decision governs). Items 3–5 chain on M1; 6–8 already parked. No viable
  "next item" while D1 is open → single-item iteration is correct.
- gemini's release-identity fix as a headless auto-adopt — it re-opens r1's closed
  corpus-equivalence gap (substitution gated by a finite corpus = probabilistic, not proof). Kept
  as option B for Mark, not silently applied.

**Parked for human (Mark — issue #1 + the doc's "Open Decision" section)**
- **Decide D1's replay-pin identity** (unblocks the whole queue): A exact-binary pin (max
  determinism, platform-locked, hermeticity-assuming) / B release-or-semantics-identity +
  corpus-gated local substitution (portable, probabilistic equivalence) / C content-addressed
  runtime-closure (max determinism, still platform-specific). 1.0 is single-machine so A/C work
  now; the log format freeze must survive to M8 multi-node, so A/C is a bet that cross-platform
  replay is a future extension. B depends on ailang#471.

**Next**: on Mark's D1 pick → unpark `w-log-epoch-decision`, land the D1 decision into the doc,
then route `w-world-library-m1` (M1) which freezes the log format per the decided pin. If ailang
answers #471 first, that may make option B cheap and pre-empt the fork.

**Retro (Gate 5)** — see this iteration's report; friction watch-items recorded (design-doc-creator
repo-layout coupling: instance 1 for the world repo). No skill edit (needs ≥2 same-gap instances).

---

## Iteration 2 — 2026-07-23 — queue HUMAN-BLOCKED on D1; no sprint (bookkeeping-only)

**Kind**: no-actionable-item iteration. The entire queue is blocked on Mark's D1 decision (parked
iter-1). No forcing; single bookkeeping deliverable + honest report (Standing rule 2).

**Context / preflight (Gate 0–1)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed). Billing tripwire: **CLEAN**
  (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff`.
- Local `dev` == `origin/dev` == `5219e54` (in sync, `git fetch` clean). CI `CI` on dev:
  **completed/success** @ `5219e5499`. No `[nightly-eval]` regression issues.
- Inbox: 1 unread — `mission-v1` controlplane status (their iter-94 dev-RED docs fix, PR #470,
  their repo). Cross-mission sender class: **not a World bug, directive, or demand** → informational,
  acked, no action. No human directive.
- Bookkeeping issue #1 (`mission-world-gh-issue`): 7 comments, **zero `@MarkEdmondson1234`** — D1
  still unanswered. Watermark advanced `19:25:04Z`→`20:13:54Z`. Created today (after the 07-20
  Monday boundary), 7<80 comments → **no weekly rotation**.
- Upstream `ailang#471` (the input that could pre-empt the D1 fork): still **OPEN, 0 comments** —
  no new information to feed the decision.

**Pick + reality-check (Gate 2)**: top `[NEXT]` is item 1 `w-log-epoch-decision` = **PARKED
needs-human-review** on D1 (Mark). Walked the queue: item 2 (M1) BLOCKED on item 1's D1 pick; items
3–5 chain on M1 landing (M1 freezes the log format D1 governs); items 6–8 explicitly parked. **No
sprint-executable OR critical-path design-doc item is independent of D1.** Verified the D1 ask is
already crisply presented to Mark on issue #1 (iter-1 report: A/B/C options table + explicit
"unblocks the whole queue" ask) — so no re-nag; a duplicate decision request would be noise on the
day the ask was posted.

**Route / execute (Gate 3)**: **none** — no design-doc-creator / planner / executor / evaluator /
quorum spawned; no worktree. Correctly declined to freelance speculative M2/M3 design ahead of the
D1-dependent M1 (would risk rework once the log format is decided — a "data before conclusions"
violation, PROGRAM.md invariant + Standing rule 2).

**Landed (Gate 3b)**: doc-only bookkeeping commit on dev (this log entry + STATUS rotation +
watermark). CI expected green (docs/`.ail`-unchanged; `verify_ail.sh` gate only).

**Routing evidence**
| Role | Pinned | Actual | Notes |
|---|---|---|---|
| Controller (triage/pick/record/retro) | `$MODEL` session | claude-opus-4-8 | opus-first, correct |
| Designer / Planner / Executor / Evaluator / Quorum | — | **not spawned** | no doc/sprint/eval this iteration; queue human-blocked |

**Metered ledger**: `metered=$0.00` — no codex/gemini/quorum calls. All work was controller-session
+ gh + local `ailang messages` (subscription/free). Ceiling `$5` untouched.

**Ruled out** (do not re-chase)
- Re-pinging Mark for the D1 decision this iteration — already crisply asked on issue #1 the same
  day (iter-1 report); a same-day duplicate is noise, not signal. The Gate-5 report notes "still
  awaiting D1" once, concisely, without re-posting the options table.
- Speculative design-doc drafting for items 3–5 (M2/M3/MCP) to fill the idle iteration — all depend
  on M1's store/log representation, which depends on D1. Drafting now risks rework; declined.
- Advancing R1 standing-value baseline capture (clause 5) as filler — it is w-approval-inbox's scope
  (parked until item 4 lands), not a freelance idle-iteration task.

**Parked for human (Mark — unchanged from iter-1, issue #1 + doc "Open Decision")**
- **Decide D1's replay-pin identity** (A exact-binary / B release-identity+corpus-gated / C
  content-addressed runtime-closure) — unblocks the whole queue. Optionally informed by ailang#471
  (portable semantics identity) if upstream answers first.

**Next**: unchanged — on Mark's D1 pick, unpark item 1, land the decision, route M1. Until then every
queue item is blocked; scheduled fires will confirm-and-idle (see retro watch-item).

**Retro (Gate 5)** — **watch-item (instance 1): fully-human-blocked queue.** The skill assumes "the
queue always has a next item," but a mission can be legitimately blocked on ONE human decision with
every downstream item dependency-chained behind it. Handled correctly here (confirm park validity,
no forcing, one concise report). If this recurs (instance 2) it justifies a shared-skill rule for
headless back-off cadence when a queue is 100% human-blocked (avoid per-fire report noise on issue
#1). No skill edit this iteration (needs ≥2 same-gap instances). No process/mission-doc change.

---

## Iteration 3 — 2026-07-24 — queue still HUMAN-BLOCKED on D1; turned the idle fire into D1 decision-support (evidence on ailang#471)

**Kind**: no-actionable-item iteration (2nd consecutive), but NOT a pure no-op — spent the fire
reducing the top blocker's human-decision cost with non-speculative, non-forcing evidence-gathering.

**Context / preflight (Gate 0–1)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed). Billing tripwire: **CLEAN**
  (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff`. Pidfile
  `mission-world.pid`=71276 is this run's own driver (alive, my process tree) — not an overlap.
- Local `dev` == `origin/dev` == `b337ee2` (in sync, `git fetch` clean). CI `CI` on dev:
  **completed/success** @ `b337ee251` (HEAD). No `[nightly-eval]` regression issues.
- Inbox: 2 unread, both informational cross-mission/automation → acked, no action, no queue impact:
  (1) `mission-v1` iter-95 status (m-budget-scoping-bug parked + quorum carve-out — their repo);
  (2) `eval-suite` auto "Eval Suite Started" notification. Neither a World bug/directive/demand.
- Bookkeeping issue #1 (`mission-world-gh-issue`): 8 comments, **zero `@MarkEdmondson1234`** — D1
  still unanswered. Created after the 07-20 Monday boundary, 8<80 comments → **no weekly rotation**.
- Upstream `ailang#471` (the input that could pre-empt/narrow the D1 fork): still **OPEN, 0 comments**
  at fire time — no answer had arrived, so the block persists.

**Pick + reality-check (Gate 2)**: top `[NEXT]` is item 1 `w-log-epoch-decision` = **PARKED
needs-human-review** on D1 (Mark). Re-walked the queue: item 2 (M1) BLOCKED on D1; items 3–5 chain
on M1 (freezes the D1-governed log format); items 6–8 explicitly parked. **No sprint-executable OR
critical-path design-doc item is independent of D1** (unchanged from iter-2). The D1 ask is already
crisply on issue #1 (iter-1: A/B/C table + "unblocks the whole queue") → **no re-nag**.

**Route / execute (Gate 3)**: **no inner-loop skill spawned** (no designer/planner/executor/
evaluator/quorum; no worktree). Instead — the "is there REALLY nothing to do?" check found genuine,
non-speculative work on the blocker itself. Read the doc's "Open Decision": Option B's *availability*
hinges on an upstream fact (does a released `ailang` expose a platform-independent **semantics-version
identity**?) = exactly ailang#471. Probed the shipped binary `v0.30.0-147-g6ed26bebd` (surface fact,
build-stable):
  - `ailang version` → release tag + git commit + build time = a **source** identity
    (platform-independent, but proves "same source", not "same semantics across platform builds").
  - runtime `ailangVersion()` → `{version, buildTime, **platform**}` = **platform-dependent**.
  - **No** `conformance`/`spec`/`semantics`/`abi`/`canonical` subcommand → **no decoupled
    semantics-hash** that would upgrade Option B from corpus-gated (probabilistic) to *proven*
    cross-platform equivalence.
  **Finding**: Option B is available only in its *weak* form today; nothing upstream makes it strong.
  So D1 reduces cleanly to **A/C** (strong bit-for-bit determinism, platform-locked) vs **B**
  (portable, probabilistic-modulo-conformance). This does NOT decide D1 (still Mark's values call) but
  removes the "maybe B can be made strong" uncertainty. Posted as evidence to ailang#471
  ([comment 5064518731](https://github.com/sunholo-data/ailang/issues/471#issuecomment-5064518731))
  and pinged the V1 loop (`msg_20260724_012702_edc4613d`, informational — does not outrank their queue).

**Landed (Gate 3b)**: doc-only bookkeeping commit on dev (this log entry + STATUS rotation +
watermark). `.ail`/docs unchanged in a way that affects `verify_ail.sh` → CI expected green.

**Routing evidence**
| Role | Pinned | Actual | Notes |
|---|---|---|---|
| Controller (triage/pick/record/retro) | `$MODEL` session | claude-opus-4-8 | opus-first, correct |
| Designer / Planner / Executor / Evaluator / Quorum | — | **not spawned** | no doc/sprint/eval this iteration; queue human-blocked |

**Metered ledger**: `metered=$0.00` — no codex/gemini/quorum calls; all work was controller-session +
`gh` + local `ailang` probe/messages (subscription/free). Ceiling `$5` untouched.

**Ruled out** (do not re-chase)
- Re-pinging Mark for the D1 decision — already crisply asked on issue #1 (iter-1); a re-post is
  noise. The report notes "still awaiting D1" once, and adds the NEW #471 evidence as fresh signal
  (not a re-ask).
- Speculative design-doc drafting for items 3–5 to fill the fire — all depend on M1's store/log
  representation, which depends on D1. Declined (rework risk; "data before conclusions"). Same as iter-2.
- **Directly editing the shared `mission-control` SKILL.md** to add the blocked-queue rule — the
  symlinked skill resolves into the **V1 checkout** (`sunholo-data/ailang/.claude/skills/…`), which the
  charter forbids World from touching, and mission-v1 iter-95 edited that file ~1h earlier (00:33). A
  cross-mission skill improvement is PROPOSED to the shared-infra owner (Mark + V1), not applied by
  World. Routed below.

**Parked for human (Mark — D1 unchanged from iter-1; PLUS one new proposal)**
- **Decide D1's replay-pin identity** — now with the #471 evidence: Option B is corpus-gated-only
  today (no proven-semantics identity exists upstream), so the live fork is **A/C (deterministic,
  platform-locked) vs B (portable, probabilistic)**. Unblocks the whole queue.
- **(New) Ratify a shared-skill back-off rule** for a 100%-human-blocked queue (proposed patch in the
  Gate-5 report). World cannot edit the shared skill (V1-checkout guardrail) → routed to V1's Gate-5 /
  Mark to apply so all missions benefit.

**Next**: unchanged — on Mark's D1 pick, unpark item 1, land the decision, route M1. Until then every
fire confirms-and-supports (chase #471, surface new decision evidence), not confirm-and-idle.

**Retro (Gate 5)** — **instance 2 of the fully-human-blocked-queue friction** (iter-2 = instance 1).
Meets the ≥2-same-gap bar for a shared-skill fix. Gap: Standing rule 2 ("the queue always has a next
item") gives no guidance for a queue that is 100% human-blocked — the risk is either per-fire idle
noise or forcing speculative work. This iteration demonstrates the RIGHT behavior (which the proposed
rule should codify): **the deliverable becomes reducing the top blocker's human-decision cost —
gather decision-supporting evidence, chase/re-probe the upstream dependency it waits on, sharpen the
framing — plus bounded bookkeeping and NO re-nag; never force a speculative downstream item.**
**Lane = PROPOSED skill fix, not applied** (World may not edit the shared skill, which lives in the V1
checkout; also avoids racing V1's concurrent skill edits). Proposed patch routed to Mark + V1 in the
report. No mission-doc/process change this iteration.

## Iteration 4 — 2026-07-24 — queue still HUMAN-BLOCKED on D1 (3rd consecutive); no new decision-support to add — bookkeeping-only + escalate the idle-fire cost

**Kind**: no-actionable-item iteration (3rd consecutive). Unlike iter-3, there was **no NEW
decision-support left to extract** — iter-3 already reduced D1 to its irreducible values-call core and
the upstream input hasn't moved. So the honest deliverable this fire is bookkeeping + escalating the
*systemic* cost (a 100%-human-blocked queue with no back-off) to the shared-infra owner, WITHOUT
re-nagging the D1 ask itself.

**Context / preflight (Gate 0–1)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed). Billing tripwire: **CLEAN**
  (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff`. Main tree clean.
- Local `dev` == `origin/dev` == `39de8a8` (in sync; `git fetch` clean). CI `CI` on dev:
  **completed/success** @ `39de8a82e` (HEAD). No `[nightly-eval]` regression issues (world repo).
- Inbox: 1 unread → informational, acked/read, no queue impact: `eval-suite` auto "Eval Suite Started:
  1 model, 42 benchmarks" (V1's nightly suite kicking off — not a World bug/directive/demand).
- Bookkeeping issue #1 (`mission-world-gh-issue`): 9 comments, **zero `@MarkEdmondson1234`** since
  watermark `2026-07-23T20:13:54Z` → D1 still unanswered. Created 2026-07-23 (after the 07-20 Monday
  boundary), 9<80 comments → **no weekly rotation**. Watermark unchanged (no Mark comment to process).
- Upstream `ailang#471` (the input that could pre-empt/narrow D1): still **OPEN**, only our own iter-3
  evidence comment (`#issuecomment-5064518731`) — no answer from anyone. Block persists.

**Pick + reality-check (Gate 2)**: top `[NEXT]` is item 1 `w-log-epoch-decision` = **PARKED
needs-human-review** on D1 (Mark). Queue re-walked, unchanged from iter-2/3: item 2 (M1) BLOCKED on
D1; items 3–5 chain on M1 (which freezes the D1-governed log format); items 6–8 explicitly parked.
**No sprint-executable OR critical-path design-doc item is independent of D1.**

**Route / execute (Gate 3)**: **no inner-loop skill spawned** (no designer/planner/executor/
evaluator/quorum; no worktree). Ran the "is there REALLY new decision-support?" check and found
**none** — and confirmed it with evidence rather than assuming:
  - The doc's own "Recommended framing" already resolves the *technical feasibility* question: 1.0 is
    single-machine (DESIGN.md §15) → **A and C are viable now**; B needs a proven upstream
    semantics-identity. Iter-3 established that identity does not exist today. So the fork is already
    at its irreducible core: **A/C (deterministic, platform-locked) vs B (portable, probabilistic)** —
    a pure values call, nothing headless can decide.
  - Re-probed the shipped binary to verify nothing upstream moved since iter-3: `ailang version` ==
    `v0.30.0-147-g6ed26bebd` (**same commit** as iter-3's probe); `ailang --help` still exposes **no**
    `conformance`/`semantics`/`spec`/`abi`/`canonical` subcommand. → Option B is still corpus-gated-only;
    **no new material to post to #471, and re-posting the same finding would be noise.**
  - Re-probing Option C's feasibility would be redundant (the doc already records A/C as viable-now).
  So there was no non-redundant, non-nagging automated work on the blocker. Correct move: do not
  manufacture a probe to look busy.

**Landed (Gate 3b)**: doc-only bookkeeping commit on dev (this log entry + STATUS rotation). No `.ail`
or `verify_ail.sh`-relevant change → CI expected green (gate is non-vacuous on existing modules).

**Routing evidence**
| Role | Pinned | Actual | Notes |
|---|---|---|---|
| Controller (triage/pick/record/retro) | `$MODEL` session | claude-opus-4-8 | opus-first, correct |
| Designer / Planner / Executor / Evaluator / Quorum | — | **not spawned** | no doc/sprint/eval; queue human-blocked |

**Metered ledger**: `metered=$0.00` — no codex/gemini/quorum calls; all work was controller-session +
`gh` + local `ailang version`/`--help` + `ailang messages` (subscription/free). Ceiling `$5` untouched.

**Ruled out** (do not re-chase)
- Re-posting the "B is corpus-gated-only" finding to #471 — iter-3 already posted it; the binary is
  unchanged, so a re-post is pure noise. Only a CHANGE upstream (a semantics-identity landing, or Mark
  answering) warrants a new post.
- Re-nagging Mark for the D1 pick — crisply asked on issue #1 since iter-1; iter-3 added the evidence.
  A third technical re-ask is noise. (The report DOES add one new *actionable* item — an offer to pause
  the loop — which is not a re-ask of the same question.)
- Speculative design drafting for items 3–5 to fill the fire — all depend on M1's store/log
  representation, which depends on D1. Declined (rework risk; "data before conclusions"). Same as iter-2/3.
- Probing Option C's binary feasibility — redundant; the doc already records A/C as viable-now.

**Parked for human (Mark)**
- **Decide D1's replay-pin identity** — fork is settled to **A/C (deterministic, platform-locked) vs B
  (portable, probabilistic; upstream semantics-identity does not exist today)**. Unblocks the whole queue.
- **(Still pending from iter-3) Ratify the shared-skill back-off rule** for a 100%-human-blocked queue.
  World cannot edit the shared skill (V1-checkout guardrail); the proposed patch (iter-3) is **still
  unapplied** (shared `SKILL.md` mtime 07-24 00:33 has no back-off/human-blocked rule) → re-surfaced to
  V1/Mark below.
- **(New, actionable) Offer**: since the loop will keep firing nightly with no work until D1 is decided,
  Mark may prefer to set `~/.ailang/state/mission-world.disabled` (pause) until he decides — World will
  NOT self-disable (guardrail: only Mark or the v1 agent on his instruction arms/disarms the loop).

**Next**: unchanged — on Mark's D1 pick, unpark item 1, land the decision, route M1. Until then a fire
with no new upstream signal is **confirm-and-report** (cheap: preflight + queue-walk + one binary probe
+ report), explicitly NOT another manufactured decision-support pass.

**Retro (Gate 5)** — **instance 3 of the fully-human-blocked-queue friction** (iter-2 = instance 1,
iter-3 = instance 2). The ≥2-same-gap bar was already met at instance 2; iter-3 PROPOSED the shared-skill
patch to Mark + V1. **Instance 3's new signal: the patch is still unapplied AND the "reduce the blocker's
decision cost" behavior has now hit its floor** — iter-3 extracted the decision-support; there is nothing
left to extract until the upstream/human input moves, so instance-3's honest behavior is cheap
confirm-and-report, not another decision-support pass. This sharpens the proposed rule: a blocked-queue
fire should (a) reduce the blocker's decision cost *while there is non-redundant support to add*, then
(b) fall to cheap confirm-and-report once that floor is hit — and after K consecutive fully-blocked
fires, escalate an explicit *pause-the-loop* offer to the human (this iteration = the first such
escalation). **Lane = PROPOSED skill fix, not applied** (V1-checkout guardrail); re-surfaced to Mark + V1
in the report with the sharpened rule text. No mission-doc/process change this iteration (the guardrails
already handle the mechanics; the gap is genuinely in the shared skill).

## Iteration 5 — 2026-07-24 — queue still HUMAN-BLOCKED on D1 (4th consecutive); confirm-and-report bookkeeping-only; standing pause offer unanswered (NOT re-nagged)

**Kind**: no-actionable-item iteration (4th consecutive). Same class as iter-4: decision-support floor
still hit, nothing upstream moved, pause offer already standing from iter-4. The honest deliverable is a
**compact heartbeat** — cheap bookkeeping + a deliberately-minimal report — NOT another proposal round or
a re-nag. The one new datum is that iter-4's pause offer went a full cycle unanswered.

**Context / preflight (Gate 0–1)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed). Billing tripwire: **CLEAN**
  (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff` (active). Main
  tree clean. Pidfile `mission-world.pid`=87378 = this run's own driver (no sibling overlap).
- Local `dev` == `origin/dev` == `5795dcb` (`git fetch` clean, in sync). CI `CI` on dev:
  **completed/success** @ `5795dcbdd` (HEAD). No `[nightly-eval]` regression issues (world repo).
- Inbox: **no unread messages** (`ailang messages list --unread` → "No messages found").
- Bookkeeping issue #1 (`mission-world-gh-issue`): 10 comments, **zero `@MarkEdmondson1234`** since
  watermark `2026-07-23T20:13:54Z` → D1 still unanswered. Created 2026-07-23 (after the 07-20 Monday
  boundary), 10<80 comments → **no weekly rotation**. Watermark unchanged (no Mark comment to process).
  No `-prev` issue (no rotation has occurred) → no predecessor-thread check needed.
- Upstream `ailang#471` (the input that could pre-empt/narrow D1): still **OPEN**, only our own iter-3
  evidence comment (last activity `2026-07-23T23:26:53Z`, author `sunholo-voight-kampff`) — no answer
  from anyone. Block persists.

**Pick + reality-check (Gate 2)**: top `[NEXT]` is item 1 `w-log-epoch-decision` = **PARKED
needs-human-review** on D1 (Mark). Queue re-walked, unchanged from iter-2/3/4: item 2 (M1) BLOCKED on
D1; items 3–5 chain on M1 (which freezes the D1-governed store/log rep); items 6–8 explicitly parked.
Re-verified the prior D1-block judgment rather than inheriting it: no sprint-executable OR critical-path
design-doc item is independent of D1.

**Route / execute (Gate 3)**: **no inner-loop skill spawned** (no designer/planner/executor/
evaluator/quorum; no worktree). Re-ran the "is there REALLY new decision-support?" check and confirmed
the floor is still hit with fresh evidence:
  - Shipped binary `ailang --version` == `v0.30.0-147-g6ed26bebd` — **same commit** as iter-3/4's probe;
    `ailang --help` still exposes **no** `conformance`/`semantics`/`spec` subcommand → Option B remains
    corpus-gated-only; nothing upstream makes it "proven". No new material to post to #471.
  - `ailang#471` unmoved (only our own comment; still OPEN). → D1's fork is unchanged at its irreducible
    core: **A/C (deterministic, platform-locked) vs B (portable, probabilistic)** — a pure values call.
  So there is no non-redundant automated work on the blocker. Correct move: cheap confirm-and-report; do
  not manufacture a probe, re-nag the D1 ask, or repeat iter-4's proposal round.

**Landed (Gate 3b)**: doc-only bookkeeping commit on dev (this log entry + STATUS rotation). No `.ail`
or `verify_ail.sh`-relevant change → CI expected green (gate is non-vacuous on existing modules).

**Routing evidence**
| Role | Pinned | Actual | Notes |
|---|---|---|---|
| Controller (triage/pick/record/retro) | `$MODEL` session | claude-opus-4-8 | opus-first, correct |
| Designer / Planner / Executor / Evaluator / Quorum | — | **not spawned** | no doc/sprint/eval; queue human-blocked |

**Metered ledger**: `metered=$0.00` — no codex/gemini/quorum calls; all work was controller-session +
`gh` + local `ailang --version`/`--help` + `ailang messages` (subscription/free). Ceiling `$5` untouched.

**Ruled out** (do not re-chase)
- Re-posting the "B is corpus-gated-only" finding to #471 — binary unchanged (same commit); a re-post is
  pure noise. Only a CHANGE upstream (a semantics-identity landing, or Mark answering) warrants a new post.
- Re-nagging Mark for the D1 pick — asked crisply since iter-1, evidenced iter-3, pause-offer iter-4. A
  4th re-ask is noise.
- Re-issuing iter-4's pause offer / shared-skill proposal as a fresh full report — both are STANDING and
  unanswered; repeating them verbatim each fire is the very noise the back-off rule warns against. This
  iteration's report is a compact heartbeat that references them as standing, not re-proposes them.
- Speculative design drafting for items 3–5 (or pre-staging per-D1-branch log formats) to fill the fire —
  all depend on M1's store/log rep, which depends on D1. Declined (rework risk; "data before
  conclusions"). Same as iter-2/3/4.

**Parked for human (Mark)** — all STANDING from prior iterations, no new asks:
- **Decide D1's replay-pin identity** — fork settled to **A/C (deterministic, platform-locked) vs B
  (portable, probabilistic; upstream semantics-identity does not exist today)**. Unblocks the whole queue.
- **(Standing, iter-3) Ratify the shared-skill back-off rule** for a 100%-human-blocked queue — still
  unapplied (World cannot edit the shared skill; V1-checkout guardrail).
- **(Standing, iter-4) Pause offer** — set `~/.ailang/state/mission-world.disabled` to stop the nightly
  fire until D1 is decided. Unanswered for a full cycle. World will NOT self-disable (guardrail).

**Next**: unchanged — on Mark's D1 pick, unpark item 1, land the decision, route M1. Until then each fire
is a compact confirm-and-report heartbeat (preflight + queue-walk + one binary/issue probe + minimal
report), NOT a manufactured decision-support pass or a re-nag.

**Retro (Gate 5)** — **instance 4 of the fully-human-blocked-queue friction.** No NEW skill signal beyond
iter-4: the proposed back-off patch is still unapplied and the decision-support floor is still hit. The
incremental datum is behavioral — iter-4's escalated pause offer went a full cycle unanswered, so the
honest response is to hold at minimal heartbeat rather than escalate again (escalating every fire would
itself be the noise the rule targets). **Lane = none applied** (the ≥2-instance skill fix is already
PROPOSED and blocked on the V1-checkout guardrail; re-proposing adds nothing). No mission-doc/process
change (guardrails already handle the mechanics correctly — this iteration executed the iter-4 rule as
designed). Memory `human-blocked-queue-no-backoff` updated to instance-4 (standing-offer-unanswered).

---

## Iteration 6 — 2026-07-24 — FIRST REAL SPRINT: w-world-library-m1 design doc written (codex designer) & quorum-direction-accepted; M1 SPRINT PARKED on a carve-out first-use gate

**Kind**: First post-D1-unblock design iteration. D1 ratified (c7864bf) → queue unblocked → picked the
[NEXT] kernel item `w-world-library-m1`. No design doc existed → routed to the ROTATION designer. Deliverable
= the M1 design doc + creation-time quorum. Outcome: doc is design-complete and quorum-**direction**-accepted,
but the go/no-go to SPRINT is parked for Mark (bounded quorum budget exhausted on two narrow non-direction
defects; a carve-out first-use decision belongs to the human).

**Context / preflight (Gate 0–1)**
- Kill switch NOT set (armed). Billing tripwire **CLEAN**. gh `sunholo-voight-kampff` (active). Main tree
  clean. Pidfile `mission-world.pid`=95872 = this run's own driver (no sibling overlap).
- Local `dev` == `origin/dev` == `c7864bf` (fetch clean, in sync). No `[nightly-eval]` open issues on
  ailang-world.
- Inbox: **6 unread, all V1-mission traffic** — nightly-eval regressions `type_safe_record_access` +
  `prompt_injection`, `eval-suite` 52/84 (61%), V1 iter-97 report. These are the AILANG **language**
  benchmarks (V1's eval suite); World owns no eval suite and the nightly regression issues are not on
  ailang-world. Triaged as **not-outranking** (a sibling mission's benchmark regressions are V1's to fix;
  they are neither a World regression nor a directive).
- Bookkeeping issue #1 (`mission-world-gh-issue`): no `@MarkEdmondson1234` comment since watermark
  `2026-07-23T20:13:54Z`. #1 created 2026-07-23 (post 07-20 Monday boundary), <80 comments → **no weekly
  rotation**. No `-prev` predecessor check needed.

**Pick + reality-check (Gate 2)**: top `[NEXT]` = item 2 `w-world-library-m1` (D1 now RATIFIED, item 1
[LANDED]). Reality-check: `grep -ri w-world-library design_docs/` → no existing doc (only the settled epoch
doc references it) → genuinely a NEW-DOC item. Repo was **bootstrap-only**: no Go host, no `.ail` world
library — just `sketches/{logepoch,worldtypes,transitions}.ail`. Not already-landed (fresh `git fetch`;
no merged PR; direct-commit check clean). Item is un-started and sprint-sized (~2–3d).

**Route / execute (Gate 3)** — designer = **ROTATION** (`mission-world-designer-rotation` last-used =
`claude:claude-fable-5`, who authored the epoch doc → NEXT = **`codex:gpt-5.6-sol`**):
- Codex preflight probe (bounded 120s): rc=0, replied `ok` → lane available.
- Real run: `codex exec --model gpt-5.6-sol --sandbox workspace-write` in an isolated worktree
  (`/tmp/wt-w-world-library-m1` off `origin/dev`), backgrounded with a 30-min `date +%s` cap; directive
  carried the design-doc-creator methodology + a **self-contained ailang-world adapting brief** (no
  std/VERSION, no changelogs, flat planned/, Conflict Surface = §14 boundary) + all settled context.
  71,648 tokens; exit 0. Produced `design_docs/planned/w-world-library-m1.md` (~510 lines) +
  `sketches/worldkernel.ail` (HashRef adopted throughout, pure plan→verify→commit, ai-check green).
- Controller INDEPENDENTLY re-verified the load-bearing claims: `ailang --version` matches; `ailang --help`
  genuinely lists `run/check/ai-check/iface/replay/serve-api` (codex's positive-existence claim TRUE);
  `./scripts/verify_ail.sh` green on 4 modules.
- **Quorum (creation-time)**: reviewer **gemini-3-1-pro** only — `gpt5-6-sol` EXCLUDED (it authored the doc;
  **generator≠judge**), quorum degraded to 1 external + controller. Controller verdict = pass both rounds.
  - **Round 1 BLOCKED** — gemini: Verification Log validated against the rig's `-dirty` build
    (`v0.30.0-147-g6ed26bebd-dirty`), inadmissible for a determinism kernel. **Fixed** (controller revision,
    mechanical/data-only): downloaded the official released **`AILANG v0.30.0`** darwin/arm64 artifact,
    checksum-verified (`sha256:ac3174e0…` matches the published `.sha256`), confirmed clean (no `-dirty`),
    re-verified all 4 sketches ai-check-green on it, re-pinned header + Verification Log.
  - **Round 2 (the one permitted re-quorum) BLOCKED** — gemini: Conflict Surface omitted the overlap between
    the released `ailang replay` and M1's Go replay engine. **Clarified in-doc**: `ailang replay <trace.jsonl>`
    is a single-program execution-trace replay; M1 `host/replay` is a store/log-level ORCHESTRATOR that
    delegates per-transition re-execution to the released binary — **§14-forced, not a new design decision**;
    layers are complementary.
- Design DIRECTION accepted both rounds; net-axiom **+12** (A1/A3/A4/A7 all ≥0); all **9** epoch-doc M1
  implications mapped 1:1.

**Landed (Gate 3b)**: doc + sketch + bookkeeping committed to dev (see commit). `.ail` change = new
`worldkernel.ail`, ai-check green on BOTH the rig binary and the clean release; CI expected green.

**Decision — SPRINT PARKED (needs-human-review), NOT the doc**: the bounded one-revision-one-requorum budget
is exhausted, both blocks were narrow non-direction defects the controller resolved, and routing straight to
sprint-planner would require the **narrow-refinement carve-out** — whose FIRST use in the World mission needs
Mark's one-time OK, and whose strict trigger (a verbatim reviewer-authored fix) the r2 objection did not
cleanly meet. So the planner→executor SPRINT is parked for Mark's go/no-go on (a) the carve-out first-use here
and (b) the §14 replay-orchestration framing. The design itself is complete and quorum-direction-accepted.

**Routing evidence**
| Role | Pinned | Actual | Notes |
|---|---|---|---|
| Controller (triage/pick/record/retro/revision) | `$MODEL` session | claude-opus-4-8 | opus-first, correct |
| Designer | ROTATION (next after claude) | **codex:gpt-5.6-sol** | probe rc=0; real run exit 0, 71.6k tok; sandbox worktree; rotation advanced → codex recorded last-used |
| Quorum reviewer | default gpt5-6-sol,gemini-3-1-pro | **gemini-3-1-pro only** | gpt5-6-sol EXCLUDED = doc author (generator≠judge); degraded N−1 by NAME |
| Planner / Executor / Evaluator | opus / opus / sonnet | **not spawned** | sprint parked for Mark before planning |

**Metered ledger**: `metered≈$0.24` — codex designer 71,648 tok (gpt-5.6-sol; CLI emitted no USD, est ~$0.20)
+ gemini quorum r1 $0.01737 + r2 $0.01793 + release download $0.00. Ceiling `$5` untouched (~5%).

**Ruled out** (do not re-chase)
- Re-spawning codex for a 3rd revision / 3rd quorum round — exceeds the bounded budget; the remaining
  question is a human go/no-go, not more automated work.
- Force-passing the doc to sprint over gemini's r2 block — Standing rule 2 (never force a guardrail) + the
  carve-out's first-use-needs-Mark rule both forbid it.
- Using the rig's `-dirty` binary as the M1 verification substrate — a determinism kernel must pin a released
  artifact; the checksum-verified released v0.30.0 is now the pin.
- Treating the V1 nightly-eval regressions as World's pick — they are V1's language benchmarks, not World's.

**Retro / Next**
- **Retro (Gate 5)**: Two frictions worth recording. (1) **`ailang-code` profile assumes a pinned released
  binary, but the rig ships only a `-dirty` dev build** — surfaced by gemini's r1 block; this iter established
  the pin (released v0.30.0, checksummed) but a durable **binary-lockfile mechanism** is a backlog item, not
  yet in place. (2) **design-doc-creator's ailang-repo assumptions** (instance 2/2 of the watch-item):
  the codex directive had to carry a full self-contained adapting brief; the shared skill still lacks an
  ailang-world adapting section — but World CANNOT edit the shared skill (V1-checkout guardrail), so this is
  a PROPOSED shared-skill fix routed to Mark/V1, not applied here. Neither triggers a mission-doc/process
  change this iteration (the carve-out + generator≠judge + bounded-quorum rules all functioned as designed).
- **Next**: on Mark's OK → route `w-world-library-m1` to **sprint-planner** (opus), then executor
  (`MISSION_EXECUTOR_MODEL=codex:gpt-5.6-sol` per the profile, evaluator sonnet — generator≠judge holds since
  codex-executor ≠ sonnet-evaluator). If Mark instead wants the replay-reuse question reopened, that's a
  bounded designer revision. Also queue a small **binary-lockfile** item so future iterations pin the released
  ailang deterministically instead of the rig's dev build.

---

## Iteration 7 — 2026-07-24 — queue HEAD still HUMAN-BLOCKED (w-world-library-m1 sprint on Mark's carve-out OK); confirm-and-report heartbeat + durable backlog capture

**Kind**: bookkeeping-only heartbeat (no sprint routed — the top actionable item is parked
`needs-human-review` on a human gate surfaced <1 day ago; no other queue item is workable).

**Context / preflight (Gate 0–1)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed). Billing tripwire: **CLEAN**
  (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff`.
- Local `dev` == `origin/dev` == `405beda` (in sync, `git fetch` clean, 0 behind). CI `CI` on dev:
  **completed/success** @ `405beda4c`.
- Bookkeeping issue `#1` (`mission-world-gh-issue`=1). **Zero** new `@MarkEdmondson1234` comments
  since watermark `2026-07-23T20:13:54Z` → the iter-6 carve-out OK is still **unanswered**.
- No `[nightly-eval]` open issues on `sunholo-data/ailang-world`. Inbox: 2 unread, both
  **V1-mission** `eval-suite` controlplane notifications (suite-started + no-op all-skipped) — not
  World's benchmarks (World owns no eval suite), triaged **not-outranking** (per Gate-0 cross-mission
  contract: a sibling's eval traffic never sets this mission's priorities).

**Pick + reality-check (Gate 2)**
- Top actionable: item 2 `w-world-library-m1` — **PARKED (needs-human-review)** since iter-6 on
  Mark's one-time OK of (a) the narrow-refinement carve-out's FIRST use in World and (b) the
  replay-orchestration framing. **Re-verified the park is correct to hold, NOT to unblock**: the r2
  gemini objection lacked a verbatim reviewer-authored `proposed_fix`, so the carve-out's condition
  (a) [every remaining objection carries a concrete reviewer fix] FAILS; and even if it held,
  first-use-in-World requires Mark's ratification. Unblocking it myself = a controller-invented
  resolution (explicitly forbidden by the carve-out text) + forcing a guardrail (Standing Rule 2).
  → left parked, correctly.
- No other item is workable: item 3 `w-worldd-m2` and items 4–8 all chain on M1 being **built**
  (the daemon/broker/projection consume the M1 store/log rep). Declined speculative M2+ design —
  same ruled-out reasoning as iters 4–5 (rework risk; "data before conclusions"). The M1 *design*
  is done and quorum-direction-accepted; only its *sprint* is gated.

**Work done (Gate 4, headless-safe)**
- **Durable backlog capture**: promoted iter-6's prose-only infra finding to a real queue row —
  item 9 `[BACKLOG] w-verify-binary-lockfile` (clause-1-infra): the rig ships `v0.30.0-151-…-dirty`
  and the repo has **no lockfile**, so the released-binary pin is re-established ad-hoc each verify.
  Queued, NOT implemented — the mechanism may generalize to the SHARED `ailang-code` profile (a
  future attended / Gate-5-shared-skill decision), and hand-editing CI headless is out of bounds.
- STATUS: added the iter-7 stamp; rotated iter-5 out to `world-mission-status-archive.md` (keeping
  newest-3: iter-7, iter-6, D1-RATIFIED).

**Routing evidence**: no role spawned (no sprint). Controller = Opus (session). `metered=$0.00`
(no codex / managed_agents / quorum calls). Budget ceiling: N/A (nothing metered).

**Ruled out** (do not re-chase)
- Unblocking the M1 sprint via the carve-out — condition (a) fails + first-use needs Mark; forbidden.
- Speculative M2–M8 design to "fill the fire" — all chain on M1 built; rework risk.
- Re-nagging the carve-out OK — surfaced <1 day ago (iter-6); the report restates it once, no new
  proposal round (established minimal-heartbeat discipline, iters 3–5).
- A new pause-the-loop offer — the iter-4 offer is still standing/unanswered; not re-raised.

**Retro / Next (Gate 5)**
- **Retro**: instance 1 of the *M1-carve-out* heartbeat block (distinct from the earlier D1 block).
  No NEW skill/process signal — the existing minimal-heartbeat discipline already covers a
  freshly-surfaced human gate; no ≥2-instance same-gap trigger fires, so **no skill edit, no
  process change** this iteration. The carve-out / generator≠judge / bounded-quorum rules all
  functioned as designed at iter-6; nothing here contradicts them.
- **Next**: on Mark's carve-out OK on #1 → route `w-world-library-m1` to **sprint-planner** (opus)
  → executor (`codex:gpt-5.6-sol`) → evaluator (sonnet; generator≠judge holds). If he reopens the
  replay-reuse framing instead → bounded designer revision. If the block persists, later iterations
  stay minimal heartbeats until a human/regression signal arrives.

## Iteration 8 — 2026-07-24 — M1 UNBLOCKED & BUILDING: w-world-library-m1 milestone 1 (pure AILANG library) LANDED on dev, CI green

**Kind**: full inner-loop sprint pass on ONE item (`w-world-library-m1`) — planner → executor →
evaluator → land. First shipped CODE of the mission. Milestone 1 of 6 (per the sprint plan).

**Context / preflight (Gate 0–1)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed). Billing tripwire: **CLEAN**
  (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff`.
- Overlap guard: `mission-world.pid`=42355 is ALIVE but is THIS run's own driver (`claude -p` under
  `tools/launchd/mission-control.sh`, my grandparent) — no concurrent iteration.
- Local `dev` == `origin/dev` == `bea1871` at start (in sync). CI `CI` on dev: **completed/success**
  @ `bea18716d`. No `[nightly-eval]` open issues on `sunholo-data/ailang-world`.
- Bookkeeping issue `#1`. **Zero** new `@MarkEdmondson1234` comments since watermark
  `2026-07-23T20:13:54Z` (watermark unchanged). Inbox: unread = V1-mission `eval-suite` /
  `nightly-eval` traffic (`type_safe_record_access`, `prompt_injection` regressions are **V1's**
  benchmarks; World owns no eval suite) → triaged **not-outranking** per the cross-mission contract.
- **State change since iter-7**: commit `bea1871` (STATUS "M1 GO, attended") records Mark's
  attended authorization — the iter-6 carve-out first-use OK **GIVEN** (option A) + §14
  replay-orchestration framing APPROVED. Recorded in the charter (not via the bot-account #1
  comment, correctly rejected by the allowlist). Queue item 2 was `[NEXT — SPRINT AUTHORIZED]`.

**Pick + reality-check (Gate 2)**
- Picked item 2 `w-world-library-m1` (SPRINT AUTHORIZED). Design doc
  `design_docs/planned/w-world-library-m1.md` present, quorum-direction-accepted (2 gemini rounds,
  artifacts in `.ailang/state/mission-quorum/`) → **skip re-quorum**. Not already-landed (no `world/`
  dir on origin/dev at pick). Estimate ~2–3d → sprint-sized, NOT a multi-week decomposition item.
- Scope confirmed against the doc's "Files to Create" (~1,925 LOC) + Acceptance Criteria. Pinned
  released binary `/tmp/ailang-v0300/ailang` re-asserted `AILANG v0.30.0` (clean, no `-dirty`).

**Route + execute (Gate 3) — all roles model-PINNED, spawned (never inline)**
- **Sprint-planner** (opus Agent): decomposed M1 into 6 CI-green, independently-committable
  milestones (M1 pure-AILANG lib 1.5d → M2 go bootstrap+hashref+canon → M3 SQLite store → M4
  archive+epoch registry → M5 replay+replay-doubling 2d → M6 CI Go gate). Total ~6d (doc's 2–3d
  judged optimistic). Plan+handoff → `.ailang/state/sprints/w-world-library-m1.plan.json`/`.handoff.md`.
  Surfaced a real finding: `verify_ail.sh` swept only `design_docs/` → new `world/` modules would
  be silently ungated.
- **Sprint-executor** (opus Agent, isolated worktree `sprint/w-world-library-m1` from origin/dev):
  built milestone 1 — `world/logepoch.ail` (84), `world/types.ail` (85), `world/contracts.ail`
  (40), `world/transitions.ail` (76); extended `verify_ail.sh` ROOTS to sweep `world/` from repo
  root (MOD010: module path == file path) + made the gate binary `AILANG_BIN`-configurable
  (CI-safe default `ailang`, no hardcoded /tmp). Handled a genuine MOD010 deviation from the plan's
  `cd world` step correctly.
- **Controller independent verify**: all 4 modules `ai-check` → `"passed":true,"error_count":0`;
  `verify_ail.sh` `checked 8 module(s)` exit 0 (4 sketches + 4 world), 0 errors.
- **Sprint-evaluator** (sonnet Agent; generator≠judge: opus executor ≠ sonnet judge): **PASS
  93/100 round 1**. Design-fidelity 10/10 (typed surface + shared contracts match the doc; verify &
  commit both call the shared predicates → no drift; every digest field `HashRef`; `cacheKey` =
  `(transitionFn, interpreter)`, epoch excluded). One pre-merge item: drop the force-added
  `.ailang/state/sprints/*.json` (inside a gitignored tree). Applied — amended the commit to strip
  it; `world/*` + script retained.
- **Land (Gate 3b)**: PR **#2** → dev; bounded 30-min CI poll → **completed success**; auto-merged
  squash `9d61d663e`, branch deleted; main tree ff'd to `9d61d66`; post-merge dev CI **green**
  (id 30070945259). Worktree removed.

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus | triage/pick/verify/record |
| Sprint-planner | `MISSION_PLANNER_MODEL`=opus | opus Agent | ✓ as pinned |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | opus Agent | ✓ as pinned. NB iter-7 "Next" mused `codex:gpt-5.6-sol` but that was the DESIGNER rotation seed, not the executor pin — env pin is authoritative; executor = opus |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | sonnet Agent | ✓ generator≠judge (opus≠sonnet) |

`metered=$0.00` — no codex / managed_agents / quorum reviewer calls this iteration (all roles ran on
Anthropic subscription Agent-tool pins). Budget ceiling ($5) not approached.

**Ruled out** (do not re-chase)
- Landing all 6 milestones this iteration — M2–M6 (Go host, SQLite, replay) are ~4.5d; milestone-by-
  milestone across iterations is the plan (each ends CI-green). Item stays [IN-SPRINT], not [LANDED].
- Carrying the `.ailang/state/sprints/*.json` progress artifact onto dev — it lives inside the
  repo-wide-gitignored `.ailang/` tree; dropped from the commit (evaluator rec). Sprint state stays
  local; if it should be tracked, that's a deliberate `.gitignore` negation, not a silent force-add.

**Retro / Next (Gate 5)**
- **Retro**: the full planner→executor→evaluator→land loop ran cleanly on the first real build —
  the ailang-code verify profile (pinned released binary as the gate), generator≠judge, bounded
  CI poll, and worktree isolation all functioned as designed. Two friction items, NEITHER meeting
  the ≥2-same-gap skill-edit bar this iteration: (1) the sprint-planner's per-module verify command
  (`cd world && ai-check logepoch.ail`) was MOD010-wrong for prefixed module names — the executor
  caught & corrected it, but a repeat would justify a sprint-planner note that ailang module paths
  are repo-root-relative; logged as instance 1. (2) `git worktree`/PR flow needed a manual merge
  after `--auto` initially read as "not set" (it did fire post-CI) — no action. No skill edit, no
  process/routing change (routing-policy change needs ≥3 evidence rows; this is the FIRST landed
  sprint — 1 datapoint).
- **Next**: route `w-world-library-m1` **milestone 2** (Go bootstrap + `host/hashref` + `host/canon`,
  leaf utils, `go.mod`, golden SHA-256 tests) to sprint-executor (opus) in a fresh worktree →
  evaluator (sonnet). Item stays [IN-SPRINT] until M6 lands the full Go host + replay-doubling + CI
  Go gate, then → [LANDED] and the doc moves to `implemented/`.

---

## Iteration 9 — 2026-07-24 — M1 milestone 2 (Go host bootstrap: host/hashref + host/canon) LANDED on dev, CI green (2 jobs)

**Kind**: inner-loop sprint pass on ONE item (`w-world-library-m1`), executor → evaluator → land.
Plan already existed (iter-8) → no planner/designer this iteration. Milestone 2 of 6.

**Context / preflight (Gate 0–1)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed). Billing tripwire: **CLEAN**
  (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff` (gh was not on
  PATH in the tool shell — `/opt/homebrew/bin` prepended per-call; no state impact).
- Overlap guard: `mission-world.pid`=55929 is ALIVE but is THIS run's own `claude -p` driver process
  (verified: identical mission-control command line) — no concurrent iteration.
- Local `dev` == `origin/dev` == `8f918b9` at start (in sync; `git branch -vv` confirmed — the
  two-ref `git rev-parse --short dev origin/dev` errored oddly but sync was confirmed independently).
  CI `CI` on dev: **completed/success** @ `8f918b914`. No `[nightly-eval]` open issues on
  `sunholo-data/ailang-world`.
- Bookkeeping issue `#1` (world-namespaced `mission-world-gh-issue`; the generic `mission-gh-issue`
  =422 is V1's, in the ailang repo — correctly NOT used). **Zero** new `@MarkEdmondson1234` comments
  since watermark `2026-07-23T20:13:54Z` (unchanged). Inbox: 1 unread = V1-mission `eval-suite`
  "started" FYI (3 models × 3 benchmarks — V1's suite; World owns none) → triaged **not-outranking**.
- No weekly rotation: issue #1 is current-week (created 2026-07-23, week of 2026-07-20), 17 comments
  (<80), watermark post-Monday boundary.

**Pick + reality-check (Gate 2)**
- Picked item 2 `w-world-library-m1`, milestone **M2** (per the iter-8 "Next" + the sprint plan's
  M2 entry). Plan exists (`.ailang/state/sprints/w-world-library-m1.plan.json`) → route straight to
  sprint-executor (no re-quorum: design doc unchanged, direction-accepted iter-6/8).
- Not already-landed: fresh `git fetch`; no `.go` files, no `host/`, no `go.mod` on origin/dev; only
  prior merged PR is #2 (M1). No sibling session (main tree clean, no MERGE_HEAD).
- M2 spec (plan): two dependency-free leaf Go packages + `go.mod`, ~340 LOC, acceptance = `go build`/
  `go test` green for `host/hashref` + `host/canon`, canon covers 8 canonicalization cases +
  idempotence + golden HashRef, hash readers reject malformed/uppercase/bare/unsupported forms.
- **Charter reconciliation**: the Repo Profile (line ~50) mandates "extend the CI workflow in the
  same PR that introduces Go code"; the sprint plan scheduled the CI Go gate at M6. M2 is the PR
  where Go first lands → directed the executor to ALSO add the go-verify CI job now (the plan's M6
  "verify_go.sh + final sweep" still stands as a superset). Go host code is never un-gated.

**Route + execute (Gate 3) — all roles model-PINNED, spawned (never inline)**
- **Sprint-executor** (opus Agent, isolated worktree `/tmp/wt-world-m2` branch
  `sprint/w-world-library-m2` from dev@8f918b9): built `go.mod` (module
  `github.com/sunholo-data/ailang-world`, go 1.26.4, no external deps), `host/hashref/`
  (hashref.go 184 + test 143: tagged `HashRef` `algo:digest` parse/render, sha256 via crypto/sha256,
  structured `HashError` rejecting malformed/uppercase-hex/bare/unsupported-tag/wrong-width, golden
  vectors) and `host/canon/` (source.go 166 + test 142: 8-step UTF-8/LF canonicalization,
  `CanonicalizationError` on invalid UTF-8/BOM/NUL, all 8 cases + idempotence + golden HashRef).
  Extended `.github/workflows/ci.yml` with a `go-verify` job (setup-go via `go-version-file: go.mod`
  → `go build ./...` → `go test ./... -count=1`); existing `ailang-verify` job untouched. Committed
  `83809e8` on the branch. **Deliberate deviation**: the plan's "pin SQLite driver at M2" was
  DEFERRED to M3 (its first importing code) — `go mod tidy` strips an unused `require`, which would
  break the gate; the leaf packages build+test clean with zero deps (`go mod tidy` = no-op, verified).
- **Controller independent verify** (data before conclusions): `go build ./...` rc=0; `go test
  ./host/... -count=1` both packages `ok` (45 test cases via `-v`); `go vet ./...` rc=0; `gofmt -l
  host/` empty; `go mod tidy` no-op; ci.yml diff is a clean job addition; the existing ailang gate
  still passes (`verify_ail.sh` rc=0, `checked 8 module(s)`).
- **Sprint-evaluator** (sonnet Agent; generator≠judge: opus executor ≠ sonnet judge): **PASS
  97/100 round 1**, no blocking defects. Ran its own build/test/vet/gofmt + independently
  checksum-verified the golden sha256 vectors (`shasum -a 256`). Confirmed Decision-2 (all 8 steps)
  and Decision-3 (HashRef invariants, lowercase-hex enforcement, structured errors) fidelity, and
  that the go-verify job genuinely gates future PRs. Judged the SQLite deferral sound. Three
  non-blocking notes (SQLite-pin-deferral doc, verify_go.sh is an M6 deliverable, one unreachable
  defensive branch `firstInvalidUTF8Offset` returns -1) — none actioned (all cosmetic / future-milestone).
- **Land (Gate 3b)**: PR **#3** → dev; both CI jobs (`ailang-code verify gate` + `go host build +
  test gate`) **completed success** on the PR head; mergeable CLEAN → squash-merged **`d5b155c`**,
  branch deleted; local dev ff'd; post-merge dev CI **green** both jobs (run 30072503676). Worktree
  removed.

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus | triage/pick/verify/record/retro |
| Sprint-planner | — | (not run) | plan pre-existed from iter-8; no re-plan needed |
| Design-doc-creator | — | (not run) | no new/revised doc; rotation state unchanged (`codex:gpt-5.6-sol`) |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | opus Agent | ✓ as pinned (codex lane not yet wired for World; interim opus per Repo Profile) |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | sonnet Agent | ✓ generator≠judge (opus≠sonnet) |

`metered=$0.00` — no codex / managed_agents / quorum reviewer calls (all roles on Anthropic
subscription Agent-tool pins). Budget ceiling ($5) not approached.

**Ruled out** (do not re-chase)
- Pinning the pure-Go SQLite driver at M2 — deferred to M3 where it is first imported (unused
  `require` is stripped by `go mod tidy` and breaks the build). Documented in the commit + evaluator-endorsed.
- Landing M3–M6 this iteration — milestone-by-milestone across iterations is the plan (each ends
  CI-green). Item stays [IN-SPRINT], not [LANDED], until M6.
- Deferring the CI Go gate to M6 (as the plan scheduled) — overridden: the charter requires CI to
  gate Go code from its first landing, so the go-verify job landed with M2.

**Retro / Next (Gate 5)**
- **Retro**: clean executor→evaluator→land pass; the ailang-code verify profile + generator≠judge +
  bounded CI polls + worktree isolation all functioned. Friction, NONE meeting the ≥2-same-gap
  skill-edit bar: (1) the sprint-planner scheduled the CI Go gate at M6, but the charter Repo
  Profile requires CI to gate Go code in the PR where it FIRST lands (M2) — reconciled by the
  controller directing the executor to add go-verify now; logged as **instance 1** of a
  "planner-vs-charter CI-gate-timing" gap (a repeat would justify a sprint-planner note to front-
  load charter-mandated CI gates to the first-landing PR). (2) A controller scripting slip — the
  first Gate-3b poll used `git rev-parse HEAD` in the MAIN checkout (returns dev's SHA, not the
  branch head) → "no run discovered"; self-caught, re-polled with `origin/sprint/...` SHA; no
  action. No skill edit, no process/routing change (routing-policy change needs ≥3 evidence rows;
  this is the 2nd landed sprint — the executor=opus / evaluator=sonnet / generator≠judge pattern
  now has 2 clean datapoints, still below the bar to change anything).
- **Next**: route `w-world-library-m1` **milestone 3** (SQLite schema + store: immutable objects,
  worlds, atomic append-only log — pins the pure-Go SQLite driver here) to sprint-executor (opus) in
  a fresh worktree → evaluator (sonnet). Item stays [IN-SPRINT] until M6 lands the epoch registry +
  replay-doubling + verify_go.sh, then → [LANDED] and the doc moves to `implemented/`.

## Iteration 10 — 2026-07-24 — M1 milestone 3 (SQLite store + atomic append-only log, Decision 4) LANDED on dev, CI green (2 jobs)

**Kind**: inner-loop sprint pass on ONE item (`w-world-library-m1`), executor → evaluator → land.
Plan already existed (iter-8) → no planner/designer this iteration. Milestone 3 of 6.

**Context / preflight (Gate 0–1)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed). Billing tripwire: **CLEAN**
  (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff` (gh not on the
  tool-shell PATH — `/opt/homebrew/bin` prepended per-call; no state impact).
- Overlap guard: `mission-world.pid`=65052 ALIVE but is THIS run's own `claude -p` mission-control
  driver (verified by command line) — no concurrent iteration.
- Local `dev` == `origin/dev` == `03efeef` at start (in sync; `git ls-remote` confirmed). CI `CI` on
  dev: **completed/success** @ `03efeef0a`. No `[nightly-eval]` open issues on
  `sunholo-data/ailang-world`.
- Bookkeeping issue `#1` (world-namespaced `mission-world-gh-issue`; generic `mission-gh-issue`=422
  is V1's, correctly NOT used). **Zero** new `@MarkEdmondson1234` comments since watermark
  `2026-07-23T20:13:54Z` (unchanged). Inbox: 2 unread = mission-world's own iter-9 report (self) +
  V1 `eval-suite` "started" FYI (3 models × 3 benchmarks — V1's suite) → both **not-outranking**.
- No weekly rotation: issue #1 is current-week (created 2026-07-23, week of 2026-07-20), <80 comments.

**Pick + reality-check (Gate 2)**
- Picked item 2 `w-world-library-m1`, milestone **M3** (per iter-9 "Next" + the sprint plan's M3
  entry). Plan exists (`.ailang/state/sprints/w-world-library-m1.plan.json`) → route straight to
  sprint-executor (no re-quorum: design doc unchanged, direction-accepted iter-6/8).
- Not already-landed: fresh `git fetch`; no `host/store/`, no `schema.sql`, no SQLite driver on
  origin/dev; only prior merged PRs are #2 (M1) and #3 (M2). No sibling session (main tree clean, no
  MERGE_HEAD).
- M3 spec (plan): `host/store/{schema.sql,store.go}` (~420 LOC), 5 tables, single-transaction
  compare-and-append (Decision 4), structured `ConflictError` on stale head, verification-cache key
  = `(transition_fn_ref, interpreter_ref)`, pins the pure-Go SQLite driver (deferred from M2).
  Verify = `go build ./...` + `go test ./host/store/... -count=1`.

**Route + execute (Gate 3) — all roles model-PINNED, spawned (never inline)**
- **Sprint-executor** (opus Agent, isolated worktree
  `~/.ailang/state/worktrees/w-world-library-m3` branch `sprint/w-world-library-m3` from
  origin/dev@03efeef): built `host/store/schema.sql` (66 LOC — `objects`, `worlds`, `log_entries`,
  `epoch_registry_heads`, `verification_cache`; each HashRef one canonical `algo:digest` TEXT
  column; frozen 6-field LogHeader verbatim; `transition_ref` separate/outside the frozen header;
  `entry_hash_ref TEXT UNIQUE`; cache PK `(transition_fn_ref, interpreter_ref)`), `host/store/store.go`
  (576 LOC — `database/sql` + pinned pure-Go `modernc.org/sqlite v1.54.0`; content-verified
  immutable object insert via host/hashref; single-transaction compare-and-append `Commit` with
  structured `*ConflictError`/`IsConflict`; pair-keyed verify cache), `host/store/store_test.go`
  (362 LOC, 7 tests). Committed `e2cfad6` on the branch. **Deliberate deviation** (self-reported): a
  non-frozen helper table `store_heads` (selected-world-head persistence for the compare-and-append
  guard) is created in `Open()` outside `schema.sql` — kept out of the frozen schema by design.
- **Controller independent verify** (data before conclusions): `go build ./...` rc=0; `go test
  ./host/store/... -count=1` `ok`; `go test ./... -count=1` all three packages `ok` (no M1/M2
  regression); `gofmt -l host/` empty; **`CGO_ENABLED=0 go build ./host/store/...` rc=0** (driver
  genuinely CGo-free → CI needs no C toolchain); tree clean at commit.
- **Sprint-evaluator** (sonnet Agent; generator≠judge: opus executor ≠ sonnet judge): **PASS
  88/100 round 1**, no blocking defects. Independently re-ran build/test/-v/vet/gofmt/tidy +
  `CGO_ENABLED=0`, read store.go/store_test.go critically, confirmed all 3 acceptance criteria
  (persistence incl. verbatim frozen-header round-trip; one-transaction compare-and-append with
  `errors.As`-matchable `ConflictError` + clean rollback proof; cache key provably the exact pair
  with epoch as metadata-only). Design-fidelity −5 for the `store_heads` schema split (schema.sql
  claims single-source-of-truth but a 6th table lives in `Open()`); judged **moderate, non-blocking,
  trivially fixable** — carried to M4, not force-fixed by the controller (respects generator≠judge).
- **Land (Gate 3b)**: PR **#4** → dev; both CI checks (`ailang-code verify gate` + `go host build +
  test gate`) **completed success** on the PR head `e2cfad6`; mergeable CLEAN → squash-merged
  **`a901c30`**, branch deleted; local dev ff'd; post-merge dev CI **green** (run 30078588978).
  Worktree removed; stale `sprint/w-world-library-m2` local branch pruned.

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus | triage/pick/verify/record/retro |
| Sprint-planner | — | (not run) | plan pre-existed from iter-8; no re-plan needed |
| Design-doc-creator | — | (not run) | no new/revised doc; rotation state unchanged (`codex:gpt-5.6-sol`) |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | opus Agent | ✓ as pinned (codex lane not yet wired for World; interim opus per Repo Profile) |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | sonnet Agent | ✓ generator≠judge (opus≠sonnet) |

`metered=$0.00` — no codex / managed_agents / quorum reviewer calls (all roles on Anthropic
subscription Agent-tool pins). Budget ceiling ($5) not approached.

**Ruled out** (do not re-chase)
- Force-fixing the `store_heads` schema split this iteration — it PASSED (88/100, non-blocking, the
  evaluator itself scoped it "for M4/next sprint"); the controller hand-editing passing executor code
  would erode generator≠judge. Carried forward as an M4 cleanup, not a defect gate.
- Landing M4–M6 this iteration — milestone-by-milestone across iterations is the plan (each ends
  CI-green). Item stays [IN-SPRINT], not [LANDED], until M6.
- A different SQLite driver — `modernc.org/sqlite v1.54.0` is the canonical pure-Go/CGo-free choice;
  `CGO_ENABLED=0` build confirms hermeticity. `go mod tidy` stable (no strip).

**Retro / Next (Gate 5)**
- **Retro**: clean executor→verify→evaluator→land pass; ailang-code verify profile + generator≠judge
  + bounded CI polls + worktree isolation all functioned. Friction, NONE meeting the ≥2-same-gap
  skill-edit bar: the only note is the evaluator's `store_heads` schema-split observation (a
  code-quality nit inside a PASS, not a skill/process gap). No skill edit. No process change. No
  routing change — this is the **3rd** landed sprint on the executor=opus / evaluator=sonnet /
  generator≠judge pattern (iter-8/9/10, all clean PASS: 93/97/88); the pattern now meets the
  ≥3-datapoint bar as *corroboration to keep it*, not to change it — no policy edit proposed.
- **Next**: route `w-world-library-m1` **milestone 4** (interpreter artifact archive + epoch-1
  registry bootstrap — semantic ID `world/epoch-registry/v1`, first `EpochRecord` = M1 interpreter
  release string) to sprint-executor (opus) in a fresh worktree → evaluator (sonnet). Fold the M3
  evaluator's carry-forward recs where they touch M4 (move `store_heads` into `schema.sql`; add a
  store-layer `entry_hash_ref` derivation test). Item stays [IN-SPRINT] until M6 lands
  replay-doubling + verify_go.sh, then → [LANDED] and the doc moves to `implemented/`.

---

## Iteration 11 — 2026-07-24 — `w-m1-ailang-hardening` design doc DONE + quorum-cleared via the RATIFIED narrow-refinement carve-out; auditable reproduction fixtures committed (`aa542a1`)

**Kind**: design + carve-out iteration on ONE item (`w-m1-ailang-hardening`, top `[NEXT]`,
Mark-directed). Produced the design doc (rotation designer, Fable) through two quorum rounds and
applied the ratified narrow-refinement carve-out (controller). Item → [IN-SPRINT]; sprint-planner
→ executor is next iteration. NOT a code sprint (no `.ail`/gate edits yet — those are the execute step).

**Context / preflight (Gate 0–1)**
- Kill switch: NOT set (armed). Billing tripwire: **CLEAN** (no API keys). gh account:
  `sunholo-voight-kampff` (gh not on tool-shell PATH — `/opt/homebrew/bin` prepended per-call).
- Overlap guard: `mission-world.pid`=2840 ALIVE = THIS run's own `claude -p` driver (verified by
  command line) — no concurrent iteration.
- Local `dev` == `origin/dev` == `b0a632a` (in sync). CI `CI` on dev: **completed/success** @
  `b0a632a2a`. No `[nightly-eval]` open issues on `sunholo-data/ailang-world`.
- Bookkeeping issue `#1` (world-namespaced `mission-world-gh-issue`; generic `mission-gh-issue`=422
  is V1's, correctly NOT used). **Zero** new `@MarkEdmondson1234` comments since watermark
  `2026-07-23T20:13:54Z`. Inbox: 1 unread = mission-world's OWN prior discoverability finding
  (`msg …_208ab38d`; already in memory `ailang-feature-discoverability-gap`, local `.mcp.json` fix
  applied, upstream asks route to shared-skill owners World can't edit) → not-outranking, marked read.
- No weekly rotation: issue #1 current-week, <80 comments.

**Pick + reality-check (Gate 2)**
- Picked top `[NEXT]` `w-m1-ailang-hardening` (Mark-directed, queued iter-10): add Z3
  `requires`/`ensures` + inline tests to the M1 `.ail` surface + a non-vacuous verify gate.
- Reality-check: no design doc (`grep` → only the mission-doc queue entry), no plan, no eval. M1
  `.ail` files present; `contracts.ail` confirmed to carry only decorative `bool` predicates, 0 tests
  (the gap is REAL). Not already-landed (fresh fetch; nothing on origin).
- **Controller empirically grounded the syntax on the pinned released `v0.30.0`** BEFORE routing
  (this is the discoverability root-cause fix — the prior iteration wrote AILANG from priors): a
  scratch probe proved `requires`/`ensures` Z3-verify via `ai-check` (2 fns "verified") and inline
  `tests [((args),exp)]` run via `ailang test`. Loaded the version-locked syntax via `ailang prompt`
  first. **Used `/tmp/ailang-v0300/ailang` (released v0.30.0), NOT the PATH `-dirty` dev build.**

**Route + execute (Gate 3) — all heavy roles model-PINNED, spawned (never inline)**
- No design doc → **design-doc-creator on the ROTATION designer `claude:claude-fable-5`** (rotation
  next after codex; Fable probe rc=0, subscription-only via key-strip; backgrounded, bounded 30-min
  cap). Directive carried a self-contained adapting brief (known repo friction: design-doc-creator
  assumes the ailang-repo layout), the controller's empirical findings, and the key ADT-return design
  question. Produced `design_docs/planned/w-m1-ailang-hardening.md`: resolves the `CommitResult`
  sum-type contract question via a Z3-proven `applyRevision` helper (`commit` stays uncontracted,
  composes it); 7 proven contracts + 8 inline tests; 22-row empirical verification log; found real
  v0.30.0 limits (V3 encoder won't inline unencodable-bodied callees, V4 `implies` unparseable, V8
  `plan` unprovable — float/empty-list literal mis-sort, V10 Z3-encoding-errors-exit-0-SILENTLY).
- **Quorum-at-pick** (`ailang design-quorum`, reviewers gpt5-6-sol + gemini-3-1-pro — distinct from
  the Fable designer, generator≠judge): **r1 BLOCKED** (gpt5-6-sol reject: aggregate-floor gate
  `MIN_VERIFIED=6` too weak — a dropped contract still clears it; gemini pass). Gate-mandated **Fable
  revision** rewrote D5 to a hardcoded required-check MANIFEST (per-module verified-identity sets +
  named-test sets, no env-overridable strength knobs). **r2 re-quorum BLOCKED** on two NEW **narrow,
  direction-preserving** objections with concrete `proposed_fix`: (gpt5-6-sol) the V-log isn't
  auditable → commit a fixture dir; (gemini-3-1-pro) Leg-1 `$()` capture swallows the python error →
  route to stderr.
- **NARROW-REFINEMENT CARVE-OUT applied.** RATIFIED for the world mission at the M1 GO (attended,
  `world-mission-status-archive.md` L3: "narrow-refinement carve-out first-use APPROVED") → "later
  iterations apply without re-asking." Both r2 objections satisfy (a) concrete verbatim `proposed_fix`
  + (b) no design-direction dispute → the **controller** (not a 3rd Fable run — the carve-out is a
  controller action, so Fable discipline is preserved) made a bounded 2nd revision applying the
  reviewers' verbatim fixes: (1) committed `design_docs/verification/w-m1-ailang-hardening/`
  (4 `.ail.txt` fixtures + `run.sh` + captured `OUTPUTS.md`; pinned binary sha256 `e9746fef…`,
  reproducing the load-bearing V-rows); (2) routed the Leg-1 python error to stderr.
- **The reviewer-demanded fixtures caught two first-draft inaccuracies** (the auditability objection
  was VINDICATED — data before conclusions): **D-A** V3 "may not call ANY user function" is
  overstated — `callsUserFn` (contract body calls the encodable-bodied `sameRef`) actually
  **verified**; corrected to "unencodable-bodied callee errors; encodable can verify" (the decision
  to inline predicate bodies still stands as the strictly-safe choice). **D-B** the leg-2 secondary
  count must be `len(tests[])==8`, NOT `passed_tests==8` — `--format json` counts passing
  contract-derived properties in `passed_tests` (`inline_tests` fixture: `passed_tests==7`), so it
  would be flaky. Both corrected in D5 + acceptance criteria + a doc "Post-fixture corrections" section.
- Committed `aa542a1` (doc + fixtures + `.mcp.json` + `.gitignore .codex/` + mission bookkeeping),
  pushed to dev. CI running for the doc-only push (no `.ail` under swept trees → ai-check gate
  unaffected; `.ail.txt` fixtures are not swept).

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus | triage/pick/empirical-probe/quorum-verdict/carve-out-2nd-revision/record/retro |
| Design-doc-creator | ROTATION (`mission-world-designer-rotation`) | **`claude:claude-fable-5`** Fable (subscription via `claude-sub` key-strip; create + 1 gate-mandated revision) | rotation next after codex; probe rc=0; write-back `claude:claude-fable-5` |
| Quorum reviewers | (design-quorum default) | gpt5-6-sol + gemini-3-1-pro | 2 rounds; generator≠judge (≠ Fable designer); reject-by-default |
| Sprint-planner | `MISSION_PLANNER_MODEL`=opus | (not run) | next iteration (design doc's §Implementation Plan is the 4-phase basis) |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | (not run) | next iteration |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | (not run) | next iteration |

`metered=$0.149` — two quorum rounds ($0.067 r1 + $0.082 r2, external reviewer API). Fable
designer + revision = **$0.00** (OAuth subscription). Well under the $5 ceiling.

**Ruled out** (do not re-chase)
- Parking the item for Mark's carve-out ratification — the carve-out was ALREADY ratified for the
  world mission at the M1 GO (attended); re-asking would waste an iteration, exactly what the
  carve-out exists to prevent. (My first-pass STATUS/queue edits assumed un-ratified and were
  corrected once the archive ratification was found.)
- A 3rd Fable run for the 2nd revision — the carve-out is a CONTROLLER action; the controller applies
  the verbatim fixes directly, keeping the one-Fable-run-per-iteration discipline intact.
- Running the sprint-planner/executor THIS iteration — the design+quorum×2+carve-out+fixtures work is
  a full iteration; the execute step (implement 3 modules + the manifest gate + PR + CI) is the next
  iteration's deliverable (mirrors iter-6 doc → iter-8 execute).
- `plan`/`verify` contracts — empirically UNPROVABLE in v0.30.0 (V8 literal mis-sort; V3 predicate
  calls) — documented as out-of-scope with the upstream escalation channel (`sunholo-data/ailang#476`).

**Retro / Next (Gate 5)**: No skill edit (the carve-out was a ratified-mechanism APPLICATION, not a
new controller-authored gate change → no re-ratification, no ≥2-friction skill signal). No
routing-policy change. One process observation logged, below action bar: the two premature
STATUS/queue edits (assuming the carve-out needed fresh ratification) cost only local churn because
the archive ratification check (Gate-2 reality-check on mission state) caught it before reporting —
instance 1 of "check the archive for prior human ratifications of a gate mechanism before assuming a
first-use park." **Next iteration**: sprint-planner (opus) → the doc's §Implementation Plan → sprint
-executor (opus, worktree: 3 `.ail` modules + the `verify_ail.sh` manifest gate, both negative tests)
→ evaluator (sonnet, generator≠judge) → PR → CI green → [LANDED], then resume `w-world-library-m1` M4.

## Iteration 12 — 2026-07-24 — `w-m1-ailang-hardening` EXECUTE attempted; Phase 1 (logepoch) LANDED on branch, Phases 2–4 BLOCKED by a v0.30.0 encoder limit invalidating doc claim V5 → PARKED for a designer revision + re-quorum (autonomous, not human-blocked); upstream issue #477 filed

**Kind**: execute iteration on ONE item (`w-m1-ailang-hardening`, top of queue). Ran the full inner
loop (planner → executor); executor STOPPED at Phase 2 on a design-contradicting empirical finding
that the controller independently reproduced. Deliverable pivoted to: preserve Phase 1, escalate the
encoder gap upstream, park for a bounded designer-revision + re-quorum next iteration. NOT a
force-through; NOT human-blocked.

**Context / preflight (Gate 0–1)**
- Kill switch: NOT set (armed). Billing tripwire: **CLEAN** (no API keys). gh account:
  `sunholo-voight-kampff` (gh not on tool-shell PATH — `/opt/homebrew/bin` prepended per-call).
- Local `dev` == `origin/dev` == `a4ec887` (in sync). CI `CI` on dev: **completed/success** @
  `a4ec887f8`. No `[nightly-eval]` open issues on `sunholo-data/ailang-world`.
- Bookkeeping issue `#1` (world-namespaced `mission-world-gh-issue`=1). **Zero** new
  `@MarkEdmondson1234` comments since watermark `2026-07-23T20:13:54Z`. Inbox: eval-suite start/partial
  FYIs (V1's suite, not World's) + own iter-11 report + a sibling `mission-v1` status — none
  outranking; the one unread eval-suite start marked read. No cross-mission DEMAND.
- No weekly rotation: issue #1 current-week, <80 comments.

**Pick + reality-check (Gate 2)**
- Picked top item `w-m1-ailang-hardening` (Mark-directed): design doc DONE + quorum-cleared (iter-11),
  no plan yet → route planner → executor → evaluator (per the iter-11 "Next"). Quorum artifact present
  (carve-out-cleared) → QUORUM-AT-PICK satisfied. Not already-landed (fresh fetch; only the iter-11
  doc commit `aa542a1` on origin). Pinned binary confirmed: `/tmp/ailang-v0300/ailang` = v0.30.0
  `e37b370` (matches doc). `world/contracts.ail` confirmed to still carry decorative predicates + 0
  tests (gap real).

**Route + execute (Gate 3) — all heavy roles model-PINNED, spawned (never inline)**
- **sprint-planner (opus)** → 4-phase sprint JSON + handoff faithful to the doc's Implementation Plan
  (`.ailang/state/sprints/w-m1-ailang-hardening.{plan.json,handoff.md}`). No redesign.
- **sprint-executor (opus, isolated worktree `/tmp/wt-w-m1-hardening`, branch
  `sprint/w-m1-ailang-hardening`)**, loaded version-locked syntax first (mission requirement).
  **Phase 1 (D3, logepoch) DONE + committed `35c3133`** — controller-verified on the pinned binary:
  `ai-check world/logepoch.ail` → `verified:2` (`sameRef`, `servesEntry`), `errors:0`;
  `test --format json world/logepoch.ail` → 8 named inline tests pass (`renderRef/sameRef/cacheKey/
  servesEntry _test_1/2`), `failed:0`, `len(tests[])==8`. Matches D3/D4/D-B exactly.
- **Executor STOPPED at Phase 2** per the design's own "STOP-and-report if the pinned binary
  contradicts a V-claim" rule: applying D2 verbatim to `contracts.ail` gave `verified:1, errors:3`.
  Only `isValidNextWorld` (World/HashRef) verified; the 3 `Proposal`-taking predicates
  (`proposalMatchesWorld`, `verificationMatchesProposal`, `commitAllowed`) Z3-errored
  `unknown sort 'Proposal'`.

**The finding (data before conclusions — controller-reproduced, not merely trusted)**
- Applied the D2 contract to `proposalMatchesWorld` against the REAL `world/types.ail`: `errors:1`,
  reason `Z3 error … Invalid constant declaration: unknown sort 'Proposal'`, **exit 0 (silent —
  the V10 class)**.
- Built a **minimal self-contained module** (a record with one user-ADT-typed field + a trivial
  `ensures`): same `unknown sort` error. Bisection (executor + controller): **any field whose type is
  a user-defined sum type (ADT) makes the enclosing record an unencodable Z3 sort** — the encoder
  declares the record sort without first declaring a datatype for the contained ADT.
- `Proposal.evidence: list[Evidence]` (Evidence is a 4-constructor ADT) is the trigger. Design claim
  **V5 ("all four contracts.ail predicates verify") is empirically FALSE against production types** —
  the iter-11 committed fixture used a **toy 2-field `Proposal`**, which is why the auditability
  objection's fixtures didn't catch it. **Achievable Z3-proven set = 4** (`applyRevision`,
  `isValidNextWorld`, `sameRef`, `servesEntry`), NOT the ratified 7 → D2/D4/D5-manifest/
  `EXACT_TOTAL_VERIFIED=7` are invalidated.

**Disposition (Gate 2/3 judgment)**
- **Did NOT force a shrunk manifest through** (Standing rule 2 — the quorum is the guardrail here; a
  7→4 gate-strength reduction is exactly the reviewers' r1 concern and must be re-blessed by
  re-quorum, not decided by the controller/executor). **Did NOT merge a partial** to dev (the gate is
  the load-bearing deliverable and it's blocked).
- **Preserved Phase 1**: pushed branch `sprint/w-m1-ailang-hardening` (durable WIP, unmerged) — it
  lands with the revised sprint.
- **Escalated upstream** (frozen-core protocol): filed `sunholo-data/ailang#477` with the minimal
  repro + bisection (two asks: declare Z3 datatypes for ADT-bearing records; make `ai-check` exit
  non-zero on `verify.errors>0`), + `mission-control` msg `msg_20260724_143026_0b2a75a0`.
- **Parked** the item `[PARKED — designer revision + re-quorum needed; NOT human-blocked]`; surfaced
  to Mark FYI (his directed item; descopes a ratified claim) but it does not block on him.

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus | triage/pick/independent-repro/bisection/disposition/record/retro |
| Sprint-planner | `MISSION_PLANNER_MODEL`=opus | **opus** (Agent-tool pin) | 4-phase plan JSON + handoff; faithful, no redesign |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | **opus** (Agent-tool pin, isolated worktree) | Phase 1 landed on branch; STOPPED at Phase 2 per design rule (correct) |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | (not run) | no complete sprint to evaluate; deferred with the re-execute |
| Design-doc-creator | ROTATION | (not run) | revision routes next iteration (rotation head = `claude:claude-fable-5`) |

`metered=$0.00` — planner + executor on opus **Agent-tool subscription pins** (session inheritance,
not the metered API); no quorum/codex/gemini spend this iteration. Well under the $5 ceiling.

**Ruled out** (do not re-chase)
- **Contracting the 3 `Proposal`-taking predicates in v0.30.0** — empirically impossible
  (record-transitively-contains-ADT ⇒ `unknown sort`, reproduced twice incl. a minimal module).
  Do not re-attempt against this binary; the fix is upstream (#477).
- **Editing `world/types.ail` to drop/flatten `evidence`** to dodge the encoder — forbidden (doc
  mandates `types.ail` byte-identical; it would change the production type surface the Go host uses).
- **Reshaping predicate signatures to avoid `Proposal`** — changes the exported contract surface
  (`verify`/`commit` + Go host depend on it); outside retrofit scope, not what D2 specifies.
- **Controller unilaterally shrinking `REQUIRED_VERIFIED` 7→4 and landing** — that is the exact
  gate-weakening the r1 quorum objected to; it must go through the sanctioned revision→re-quorum loop.
- **Merging Phase 1 alone to dev** — a partial of a sprint whose doc is under revision; the gate
  doesn't yet enforce the new contracts. Kept on the pushed branch instead.

**Retro / Next (Gate 5)**: No skill edit — the STOP-and-report mechanism worked exactly as the design
intended (the executor halted rather than working around a false V-claim; the controller reproduced
before concluding). No routing-policy change. One process observation (below the ≥2-instance skill
bar, logged for pattern-watch): **a design doc's empirical fixtures are only as strong as the types
they exercise** — iter-11's auditability fixtures used a toy `Proposal`, so they "passed" while the
production type fails; a future fixture-review gate could require fixtures to import the REAL types
they claim to validate (instance 1 of "fixture-vs-production-type drift"). **Next iteration
(autonomous):** rotation designer revises V5/D2/D4/D5 to the achievable-4 scope (the 3 `Proposal`
predicates → documented-limitation rows with inline tests as their machine check, per the doc's own
V8/§5 pattern) → **re-quorum ONCE** → resume Phases 2–4 (transitions `applyRevision` + contracts
`isValidNextWorld` + the corrected manifest gate + NT1/NT2) on the existing branch → evaluator
(sonnet) → PR → CI green → [LANDED]; then resume `w-world-library-m1` M4.

## Iteration 13 — 2026-07-24 — `w-m1-ailang-hardening` LANDED (PR #5 → squash `d0009c8`, dev CI green): 4 Z3-proven contracts + 14 inline tests + a hardcoded bounded non-vacuous required-check-manifest gate; two new toolchain findings (V26 bounded-waits, V27 z3-on-CI) landed as fixes

**Kind**: execute iteration on ONE item (`w-m1-ailang-hardening`, top of queue) — ran the pre-authorized
autonomous iter-12 "Next" path end-to-end (empirical grounding → designer revision → re-quorum →
carve-out 2nd revision → planner refresh → executor → evaluator → PR → CI → merge). LANDED.

**Context / preflight (Gate 0–1)**
- Kill switch armed. Billing tripwire **CLEAN**. gh account `sunholo-voight-kampff` (gh on `/opt/homebrew/bin`, prepended per-call).
- World-namespaced state: bookkeeping issue `#1` (`mission-world-gh-issue`=1). **Zero** new `@MarkEdmondson1234`
  comments (only bot author on #1; watermark `2026-07-23T20:13:54Z`). Inbox: **one** unread eval-suite-START FYI
  (V1's local-GPU suite, `bcb87630`) — informational, not outranking; no cross-mission DEMAND, no regression.
  No `[nightly-eval]` open issues on the world repo.
- Local `dev` == `origin/dev` == `6de21b5` (in sync). CI `CI` on dev: **completed/success**.
- No weekly rotation (issue #1, current week, <80 comments).

**Pick + reality-check (Gate 2)**
- Picked top item `w-m1-ailang-hardening` — the iter-12 "Next" explicitly pre-authorized the autonomous
  revision→re-quorum→resume path (parked but NOT human-blocked). Fresh-origin already-landed check: not on
  dev (only iter logs + the iter-11 doc). Branch `sprint/w-m1-ailang-hardening` held Phase 1 (`35c3133`,
  logepoch: 2 verified + 8 tests). Pinned binary confirmed: `/tmp/ailang-v0300/ailang` v0.30.0 `e37b370`.

**Route + execute (Gate 3) — all heavy roles model-PINNED, spawned (never inline)**
- **Controller empirical grounding FIRST** (data-before-conclusions; closes the iter-11 toy-Proposal blind spot):
  on the pinned binary against the REAL `world/types.ail` — (a) a contract on `proposalMatchesWorld(...,p: Proposal)`
  → `unknown sort 'Proposal'`, `errors:1`, exit 0 SILENT (reproduced iter-12); (b) `isValidNextWorld` (World/HashRef)
  → verified; (c) inline `tests` on a Proposal-taking predicate with a full 9-field literal → run + pass; (d) 3
  predicates tests-only + isValidNextWorld proven → ai-check `verified:1, errors:0` (clean gate). Achievable = 4.
- **design-doc-creator on the ROTATION designer `codex:gpt-5.6-sol`** (advanced from last-used `claude`;
  probe rc=0; independent authorship deliberately chosen — the original doc's blind spot was Fable's). Surgical
  revision: V5 superseded + V23/V24/V25 added; D2/D4/D5/Goals/Acceptance/Conflict-Surface descoped 7→4; the 3
  Proposal predicates → documented-limitation + tests-only (bodies UNCHANGED → anti-drift preserved); totals
  4 verified / 14 tests. Faithful, doc-only. Committed `0b623e7`.
- **Re-quorum ONCE** (`ailang design-quorum`, controller verdict pass, metered **$0.095**): **gemini-3-1-pro PASS**
  (non-blocking: schema-coupling note + a `mod=${mod#./}` path-norm catch — folded in); **gpt5-6-sol REJECT**, ONE
  blocking objection — the gate ran `ai-check`/`ailang test` with **no wall-clock deadline** (bounded-waits, Standing
  Rule 6). Concrete verbatim `proposed_fix`, no design-direction dispute.
- **NARROW-REFINEMENT CARVE-OUT applied** (already ratified for world at the M1 GO, attended) → controller bounded
  2nd revision applying the reviewer's verbatim fix: V26 (empirically established `ai-check -timeout` is per-function
  Z3 only, `ailang test` has none) + **Leg 0 `run_bounded`** (hardcoded non-env-overridable `GATE_LEG_TIMEOUT_S=120`/
  `GATE_TEST_TIMEOUT_S=180`, `start_new_session` process group, SIGKILL-on-expiry, exit 124 fatal; controller
  pre-validated the mechanism on 4 cases) + NT3/NT4 + acceptance criterion. Routed straight to planner (no 3rd
  quorum). NOT a force-pass (Rule 2: direction uncontested). Committed `4685063`.
- **sprint-planner (opus)** refreshed the plan/handoff (P1 done, P2/P4 rescoped to 4-verified/14-tests + Leg 0).
- **sprint-executor (opus, worktree `/tmp/wt-w-m1-hardening`)** — loaded version-locked syntax first — built
  P2 (contracts: `isValidNextWorld` inlined+proven; 3 predicates tests-only, 6 tests; `e38ffa1`), P3 (transitions
  `applyRevision` proven helper + rewired commit's Applied arm; `7242e64`), P4 (the bounded manifest gate; `7b20411`).
  **Controller INDEPENDENTLY re-verified** on the pinned binary: full gate exit 0 (4/4 verified, 14 tests); NT1
  mutation (strip applyRevision's contract) → gate exit 1 naming the identity (teeth confirmed); deadlines hardcoded;
  3 predicates carry no `ensures`; `world/types.ail` byte-identical; diff scoped to the 3 `.ail` + the `.sh`.
- **sprint-evaluator (sonnet, generator≠judge: opus executor ≠ sonnet judge) PASS 97/100** — no blocking issues;
  independently re-ran every check; −3 only for the pre-move `Status: Planned` (a Gate-4 step, done here).

**Gate 3b — CI (the item is not LANDED until remote CI is green)**
- Doc → `implemented/` + `Status: Implemented` (`5c69428`); pushed; **PR #5** opened.
- **First CI run RED** → root-caused to **V27: `ai-check` shells out to an external `z3`** (PATH + hardcoded
  `/usr/bin`, `/usr/local/bin`, `/snap/bin`, `/opt/homebrew/bin` — confirmed in the binary's strings) and **SKIPS
  SILENTLY** when absent. A bare `ubuntu-latest` runner has no z3, so every contract vanished from `verify.results[]`
  (`isValidNextWorld MISSING`, the V20 class) — the identical binary (`e37b370`) + z3 4.16.0 verifies it locally on
  darwin. Fix (in-scope infra, analogous to the go-verify job added in M2): the `ailang-verify` job installs **Z3
  4.16.0** (x64-glibc-2.39, sha256-pinned) to `/usr/local/bin/z3`. Recorded as V27 + acceptance-criteria update
  (`622a543`). **Re-run CI GREEN** (bounded poll): `4/4 required world/ identities verified, all 14 named tests pass`,
  both jobs success.
- Squash-merged PR #5 → dev `d0009c8`; **dev CI on the merge commit GREEN** (bounded poll). Worktree removed.

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus | triage/pick/empirical-grounding/carve-out-2nd-revision/independent-verify/record/retro |
| Design-doc-creator | ROTATION (`codex:gpt-5.6-sol`) | **codex gpt-5.6-sol** (executor recipe, backgrounded, 30-min cap; probe rc=0) | surgical 7→4 revision; metered (see below) |
| Sprint-planner | `MISSION_PLANNER_MODEL`=opus | **opus** (Agent-tool pin) | plan refresh (P1 done, P2/P4 rescoped) |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | **opus** (Agent-tool pin, isolated worktree) | P2/P3/P4; controller-reproduced |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | **sonnet** (Agent-tool pin) | PASS 97/100; generator≠judge (opus≠sonnet) held |

`metered≈$0.44` — re-quorum **$0.095** (gpt5-6-sol $0.067 + gemini-3-1-pro $0.028) + codex designer **~$0.35** (est,
~69k tok gpt-5.6-sol); planner/executor/evaluator on opus/sonnet subscription Agent-tool pins ($0.00). Under the $5 ceiling.
Designer rotation advanced `claude`→**`codex`** (write-back `codex:gpt-5.6-sol`).

**Ruled out** (do not re-chase)
- **Proving the 3 Proposal predicates in v0.30.0** — impossible (ADT-in-record `unknown sort`, V23; upstream #477).
  They are tests-only by design; do not re-attempt a contract on them against this binary.
- **A gate that assumes z3 is present on CI** — false on bare `ubuntu-latest` (V27). CI MUST install z3; the gate is
  vacuous otherwise (silent skip, exit 0). Do not remove the CI Z3-install step.
- **Trusting `ai-check`/`ailang test` exit codes** — a Z3 encoding error exits 0 (V10) and a missing z3 skips silently
  (V27); the gate parses JSON. A `124` from `run_bounded` is the one fatal, non-advisory code.
- **A full planner re-spawn producing a fresh sprint** — unneeded; P1 was landed and P3 unchanged, so the planner
  refreshed the existing plan (P2/P4 rescope) rather than redesigning.

**Retro / Next (Gate 5)**: No skill edit, no routing-policy change. Two pattern-watch process notes (both below the
≥2-instance skill-edit bar, logged for future correlation): **(1) V27 z3-on-CI** — the shipped-binary verify profile
silently no-ops its Z3 leg without an external solver; a future `ailang-code` mission on a fresh CI will hit this, so
the world CI now carries the Z3-install as the reference pattern (candidate: fold a "CI installs the pinned Z3" note into
the shared `ailang-code` verify-profile guidance if a 2nd mission hits it). **(2) fixture-vs-production-type drift**
(iter-12 instance 1; this iteration's empirical-grounding-first step was the mitigation — verify claims against the REAL
exported types, not toy fixtures). **Next iteration: `w-world-library-m1` M4** — interpreter artifact archive + epoch-1
registry bootstrap (`world/epoch-registry/v1`); also fold the M3 carry-forward (`store_heads` → `schema.sql`).

## Iteration 14 — 2026-07-24 — `w-world-library-m1` M4 LANDED (PR #6 → squash `8133573`, dev CI green, both jobs): interpreter artifact archive + epoch-1 registry bootstrap (Decisions 5+6) + the M3 `store_heads` carry-forward

**Kind**: execute iteration on ONE item (`w-world-library-m1`, top of queue, IN-SPRINT). Mid-sprint EXECUTE — the
doc was quorum-direction-accepted and the 6-milestone plan approved at M1 (iter-8); M1–M3 landed (iter-8/9/10). No
new design-doc-creator, no re-quorum: routed straight to sprint-executor (Gate 3 "Plan exists" lane).

**Context / preflight (Gate 0–1)**
- Kill switch armed. Billing tripwire **CLEAN** (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account
  `sunholo-voight-kampff` (gh on `/opt/homebrew/bin`, prepended per-call). Pidfile `mission-world.pid`=85308 =
  this run's own driver ancestor (walked ps tree; no overlap).
- World-namespaced state: bookkeeping issue `#1` (`mission-world-gh-issue`=1). **Zero** new `@MarkEdmondson1234`
  comments on #1 (watermark `2026-07-23T20:13:54Z`). Inbox: 4 unread, all informational — an `eval-suite` START
  notice (V1's local-GPU suite), the sibling `mission-v1` iter-101 report (no world-directed demand → cross-mission
  class, no action), and my own two iter-13 reports. None outranked the queue; no genuine regression, no
  cross-mission DEMAND. Marked read. No `[nightly-eval]` open issues on the world repo.
- Local `dev` == `origin/dev` == `d690e45` (in sync after `git fetch`). CI `CI` on dev: **completed/success** @
  `d690e458`. No weekly rotation (issue #1, current week, <80 comments).

**Pick + reality-check (Gate 2)**
- Picked top item `w-world-library-m1` M4 — the iter-13 "Next" explicitly named it. Fresh-origin already-landed
  check (`git fetch` + `git log origin/dev --grep=M4/archive`): NOT landed — no `host/archive/`, no
  `world/epoch-registry/`, no merged PR mentioning archive. Genuinely the next build step.
- Read the plan's M4 milestone spec + acceptance_checks + verify_commands, and the design doc's Decision 5
  (Epoch Registry Objects) + Decision 6 (Interpreter Artifact Archive) — the exact M4 scope. Confirmed the M3
  evaluator's −5 carry-forward (`store_heads` created inline in `store.Open()` rather than `schema.sql`).
- Build/test sanity on current dev: `go build ./...` OK, `go test ./host/...` all green, pinned binary
  `/tmp/ailang-v0300/ailang` present (v0.30.0 `e37b370`).

**Route + execute (Gate 3) — heavy roles model-PINNED, spawned (never inline)**
- Isolated worktree `git worktree add -b sprint/w-world-library-m1-m4 /tmp/wt-w-m1-m4 origin/dev` (from `d690e45`,
  build OK). NEVER the shared main tree.
- **sprint-executor (opus, MISSION_EXECUTOR_MODEL=opus, Agent-tool pin, worktree)** — directive carried the full
  API context (store/hashref public surface, Decisions 5+6 mechanics, the 3 deliverables + acceptance checks +
  verify commands + CI-safety requirement). Delivered 3 commits: `78d3c6a` (fold `store_heads` → `schema.sql`),
  `d771c52` (archive, Decision 6), `da70f02` (registry bootstrap, Decision 5). Files: `host/archive/archive.go`
  (387) + `_test.go` (243); `host/registry/registry.go` (159) + `_test.go` (158); `schema.sql` (+10);
  `store.go` (−9). Net +957/−9.
- **Controller INDEPENDENTLY re-verified** in the worktree (data-before-conclusions): `go build ./...` rc=0;
  `go test ./host/... -count=1` all green (archive 8, registry 5, store 7); `go vet ./host/...` clean; `gofmt -l
  host/` empty. Spot-checks: `store_heads` CREATE TABLE fully removed from `Open()` (0 hits) and present in
  `schema.sql`; archive uses `io.TeeReader(src, hasher)` single stream + `os.Rename` atomic + `archivedPerm =
  0o555`; archive tests use synthetic POSIX shell fake-interpreters with skip guards, and `TestArchivePinnedInterpreter`
  skips when the pinned binary is absent → linux CI needs no binary/Z3; registry uses `store.EpochRegistryV1` +
  ordered `EpochRecord.Candidates` (advisory).
- **sprint-evaluator (sonnet, MISSION_EVALUATOR_MODEL=sonnet, Agent-tool pin; generator≠judge: opus executor ≠
  sonnet judge) PASS 88/100 round 1** — independently re-ran build/test/vet/gofmt (15/15 green), confirmed
  single-stream hash-while-copy, fsync→chmod→rename order, `errors.As`-distinguishable ReplayError Kinds,
  genuine idempotence, no privileged side channel, full `store_heads` removal. No blocking issues.

**Gate 3b — CI (item not LANDED until remote CI green on the merge)**
- Pushed branch; **PR #6** opened (base dev). Expected checks for this Go+schema diff: the two world CI jobs
  (`go host build + test gate`, `ailang-code verify gate`) — both present in the run, no path-filtered N/A.
- Bounded poll (30-min cap): PR #6 CI **completed/success**, both jobs green. Squash-merged (`--delete-branch`) →
  dev `8133573`. Bounded poll of dev CI on the merge commit `8133573`: **completed/success**. Worktree removed;
  local dev fast-forwarded to `8133573`.

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus (`claude-opus-4-8`) | triage/pick/reality-check/independent-verify/record/retro |
| Design-doc-creator | — | not spawned | doc already quorum-cleared; no new/revised doc |
| Sprint-planner | — | not spawned | plan pre-existed (iter-8); M4 milestone spec used as-is |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | **opus** (Agent-tool pin, isolated worktree) | 3 commits; controller-reproduced |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | **sonnet** (Agent-tool pin) | PASS 88/100; generator≠judge (opus≠sonnet) held |

**Metered ledger**: `metered=$0.00` — executor opus + evaluator sonnet on subscription Agent-tool pins; no
designer/quorum/codex/gemini metered spend this iteration. Ceiling ($5) not approached. Designer rotation
UNCHANGED (`codex:gpt-5.6-sol` — no new doc authored).

**Ruled out** (do not re-chase)
- Re-running design-doc-creator or a quorum for M4 — the doc is already quorum-direction-accepted and M1–M3
  executed against it; mid-sprint milestones route straight to the executor.
- Folding `host/registry` into `host/store` — the executor made it a separate package; the evaluator judged this
  correct package-first design (imports store+hashref, no cycle, keeps registry semantics out of the physical
  store layer). The plan's "within store/archive path" wording meant "through the store object/head mechanism,"
  not "literally in those files." Not a defect.
- Proving acceptance-check #1 (`writtenBy` stamps interpreter HashRef+version on every log entry) in M4 — that is
  a log-WRITE-time caller contract that lands in **M5** (replay); M4 delivers the archive + resolver the stamp
  will reference. The archive side (hash, manifest, resolver) is proven here.

**Retro / Next (Gate 5)**: No skill edit, no routing-policy change. This is the **4th consecutive clean landed
sprint** on the opus-executor / sonnet-judge / generator≠judge pattern (iter-8/9/10/14 all PASS ≥88, zero
blocking defects) — corroborates keeping the pattern; still below any change bar. Two non-blocking M4
carry-forwards recorded for M5/M6 (not skill/process gaps, ordinary code follow-ups): (1) add
`TestArchiveExecFailureOnBadVersion` to cover `KindExecFailure` (defined + used, untested); (2) harden the
sidecar-present/executable-absent idempotence path (currently `Resolve()` catches it as `KindAbsentArtifact`
before any execution — safe, but archival could re-verify the file exists on disk). **Next iteration:
`w-world-library-m1` M5** — replay + replay-doubling (bit-for-bit episode reconstruction); the `writtenBy`
interpreter-HashRef stamp lands here, and it folds the two M4 carry-forwards.

---

## Iteration 15 — 2026-07-24 — `w-world-library-m1` M5 LANDED (PR #7 → squash `ef06937`, dev CI green, both jobs): replay engine + replay-doubling + fixture episode (Decision 7); evaluator PASS 73/100 with an in-PR CI false-green fix (B1)

**Kind**: execute iteration on ONE item (`w-world-library-m1`, top of queue, IN-SPRINT). Mid-sprint EXECUTE — the
doc was quorum-direction-accepted and the 6-milestone plan approved at M1 (iter-8); M1–M4 landed (iter-8/9/10/14).
No new design-doc-creator, no re-quorum: routed straight to sprint-executor (Gate 3 "Plan exists" lane). M5 is the
**load-bearing, acceptance-defining** milestone (subprocess-driven replay + replay-doubling hermeticity).

**Context / preflight (Gate 0–1)**
- Kill switch armed. Billing tripwire **CLEAN** (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account
  `sunholo-voight-kampff` (gh at `/opt/homebrew/bin`, prepended per-call — not on the default shell PATH).
  Pidfile `mission-world.pid`=95235 = this run's own `claude -p` driver (verified via `ps`; my bash child ≠
  overlap).
- World-namespaced state: bookkeeping issue `#1` (`mission-world-gh-issue`=1; the generic `mission-gh-issue`=422
  is the V1 loop's, not World's). **Zero** new `@MarkEdmondson1234` comments on #1 (watermark
  `2026-07-23T20:13:54Z`). Inbox: 1 unread — the sibling `mission-v1` iter-102 controlplane report (VM parity
  re-scope + a soundness bug awaiting Mark on V1's #422); a cross-mission STATUS FYI, no World-directed demand,
  no soundness bug in World → did not outrank, marked read. No `[nightly-eval]` open issues on the world repo.
- Local `dev` == `origin/dev` == `f116174` (in sync after `git fetch`). CI `CI` on dev: **completed/success** @
  `f1161740`. No weekly rotation (issue #1, created 2026-07-23 — after this Monday's 07:00 boundary — and <80
  comments).

**Pick + reality-check (Gate 2)**
- Picked top item `w-world-library-m1` M5 — the iter-14 "Next" explicitly named it. Fresh-origin already-landed
  check: `git log origin/dev --grep`, no `host/replay/` dir, no merged replay PR → NOT landed, genuinely the next
  build step.
- Read the plan's M5 milestone spec (files, 5 acceptance_checks, verify_commands, ci_green_boundary) + the design
  doc §14 (replay DELEGATES per-transition execution to the archived released artifact — never reimplements the
  interpreter). Read the existing host public APIs (store/archive/registry/canon/hashref) so the executor directive
  carried accurate interfaces. Pinned binary `/tmp/ailang-v0300/ailang` present + clean (`v0.30.0` `e37b370`,
  no `-dirty`).

**Route + execute (Gate 3) — heavy roles model-PINNED, spawned (never inline)**
- Isolated worktree `git worktree add -b sprint/w-world-library-m1-m5 /tmp/wt-w-world-m1-m5 origin/dev` (from
  `f116174`). NEVER the shared main tree.
- **sprint-executor (opus, MISSION_EXECUTOR_MODEL=opus, Agent-tool pin, worktree)** — directive carried §14, the
  M5 files/acceptance/verify, existing host APIs, and the two M4 carry-forwards to fold in. Delivered `host/replay/
  replay.go` (384), `replay_test.go` (640, 12 tests), `testdata/{transition_fixture.ail, recorded_result.bytes,
  recorded_world_hash.txt}`. Commit `3559a07`.
- **Controller INDEPENDENTLY re-verified** in the worktree (data-before-conclusions): `go build ./...` rc=0;
  `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/... -count=1` all 6 packages green incl. **replay 12/12**
  (11s runtime — genuinely driving the pinned binary as a subprocess, not skipping); `go vet` clean; `gofmt -l
  host/` empty. Spot-checked `replay.go`: resolves execPath from **each entry's `Interpreter` HashRef** (not the
  registry), verifies exe against that hash, consults the `(TransitionFn,Interpreter)` verify cache, executes via
  `exec.CommandContext` — no interpreter reimplementation; goldens are committed anchors (not tautological).
- **sprint-evaluator (sonnet, MISSION_EVALUATOR_MODEL=sonnet, Agent-tool pin; generator≠judge: opus ≠ sonnet)
  PASS 73/100 round 1** — independently re-ran the gate, confirmed genuine hermeticity (zero `registry` import in
  `replay.go` — structural, not just behavioral, proof), committed golden anchor, no Go-side interpreter
  reimplementation, pair-only cache key. Raised ONE **BLOCKING merge-condition (B1)**: the CI `go-verify` job ran
  `go test ./...` **without `AILANG_BIN`**, so all 12 replay tests silently `t.Skip`-ed in CI — a false-green of
  exactly the iter-13 V27 z3-silent-skip class. Non-blocking carry-forwards NB1 (fixture-scoped world-hash =
  `SHA-256(result bytes)`, documented in code), NB2 (interpreter-member end-to-end re-verify, env-constrained by
  one archived binary), NB5 (`scripts/verify_go.sh`, M6 scope).
- **B1 fixed in-PR (bounded follow-up, opus executor)** — added a `go-verify` step that downloads the **pinned
  v0.30.0** linux binary (by TAG, not `latest` — the replay goldens are v0.30.0-scoped; `latest` drift would break
  bit-for-bit replay), sha256-verifies it, asserts `v0.30.0` loudly, and exports `AILANG_BIN` via `$GITHUB_ENV`.
  Also fixed the NB3 comment typo + checked the 5 delivered M5 acceptance boxes in the design doc + documented NB1
  at the world-hash site. Commit `8ba5fe9`. Controller re-verified all four gate commands still green on the fix
  commit.

**Gate 3b — CI (item not LANDED until remote CI green on the merge)**
- Pushed branch; **PR #7** opened (base dev). Expected checks for this Go+CI-YAML diff: the two world CI jobs
  (`go host build + test gate`, `ailang-code verify gate`) — both present, no path-filtered N/A. PR mergeable.
- Bounded poll (25-min cap): PR #7 CI **completed/success**. **Verified the B1 fix WORKED**: the go-verify job log
  shows `AILANG v0.30.0` (version-assert passed) and `ok …/host/replay 3.306s` — the replay tests **actually ran
  in CI** (no SKIP), which was the whole point. Squash-merged (`--delete-branch`) → dev `ef06937`. Bounded poll
  (20-min cap) of dev CI on the merge commit `ef06937`: **completed/success**, both jobs. Worktree removed.

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus (`claude-opus-4-8`) | triage/pick/reality-check/independent-verify/record/retro |
| Design-doc-creator | — | not spawned | doc already quorum-cleared; no new/revised doc |
| Sprint-planner | — | not spawned | plan pre-existed (iter-8); M5 milestone spec used as-is |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | **opus** (Agent-tool pin, isolated worktree) | M5 (`3559a07`) + B1 fix (`8ba5fe9`); controller-reproduced |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | **sonnet** (Agent-tool pin) | PASS 73/100; generator≠judge (opus≠sonnet) held; raised B1 |

**Metered ledger**: `metered=$0.00` — executor opus + evaluator sonnet on subscription Agent-tool pins; no
designer/quorum/codex/gemini metered spend this iteration. Ceiling ($5) not approached. Designer rotation
UNCHANGED (`codex:gpt-5.6-sol` — no new doc authored).

**Ruled out** (do not re-chase)
- Re-running design-doc-creator/quorum for M5 — doc already quorum-direction-accepted; mid-sprint milestone.
- Treating B1 as an M6 deferral — the evaluator's accepted criterion is "CI asserts …"; a replay suite that
  silently skips in CI makes M5's load-bearing acceptance unenforceable remotely. Fixed in the SAME PR (the
  evaluator's explicit recommendation), not carried to M6.
- Downloading `latest` in the go-verify CI step (mirroring the existing `ailang-verify` job) — rejected: the
  replay goldens are v0.30.0-scoped, so the go-verify download is pinned to the **v0.30.0 tag** for drift-safety.
  A future release could break bit-for-bit replay under `latest`.
- The NB1 "world hash" is not a defect to fix now — for the M1 FIXTURE episode `SHA-256(result bytes)` is a
  genuine deterministic content address anchored to a committed golden; a real episode hashing the full `World`
  struct is future scope (documented in code, not re-architected here).

**Retro / Next (Gate 5)**: No skill edit, no routing-policy change. The B1 CI false-green is the **2nd instance**
of the silent-skip-in-CI class (V27 z3 was the 1st, iter-13) — but both were caught by the gate (evaluator here,
first CI red there) and fixed in-PR, and the fix pattern is identical (download the pinned artifact + assert
version in the job that needs it). It is a recurring *code/CI* pattern, NOT a loop-process gap the skill could
prevent, so it stays a pattern-watch note below the ≥2-instance skill-edit bar (the skill can't know a repo's
test-skip guards). **5th consecutive clean-landed sprint** on opus-executor / sonnet-judge / generator≠judge
(iter-8/9/10/14/15) — corroborates keeping the pattern. **Next iteration: `w-world-library-m1` M6** — CI Go gate
finalize + `scripts/verify_go.sh` + final green sweep → then item `[LANDED]` + doc → `implemented/`. M6 also picks
up the documented carry-forwards NB2 (interpreter-member end-to-end re-verify) and NB5 (`verify_go.sh`).

## Iteration 16 — 2026-07-27 — `w-world-library-m1` M6 LANDED (PR #8 → squash `a07ac96`, dev CI green, both jobs): CI Go gate + `scripts/verify_go.sh` (anti-false-green) + NB2 end-to-end — **the entire w-world-library-m1 item LANDED (all 6 milestones); doc → implemented/**

**Kind**: execute iteration on ONE item (`w-world-library-m1`, top of queue, IN-SPRINT). Mid-sprint EXECUTE, the
final milestone (M6). The doc was quorum-direction-accepted and the 6-milestone plan approved at M1 (iter-8);
M1–M5 landed (iter-8/9/10/14/15). No new design-doc-creator, no re-quorum, no planner: routed straight to
sprint-executor (Gate 3 "Plan exists" lane). M6 wires the Go gate into CI + local, closes the NB2/NB5
carry-forwards, and lands the whole item.

**Context / preflight (Gate 0–1)**
- Kill switch armed. Billing tripwire **CLEAN** (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account
  `sunholo-voight-kampff` (gh at `/opt/homebrew/bin`, prepended per-call — not on the default shell PATH).
- World-namespaced state: bookkeeping issue `#1` (`mission-world-gh-issue`=1; the generic `mission-gh-issue`=484
  is the V1 loop's, not World's — read the world-namespaced files per the Repo Profile). Pidfile
  `mission-world.pid`=47511 = this run's own `claude -p` driver (verified via `ps`; my bash child ≠ overlap).
  **Zero** new `@MarkEdmondson1234` comments on #1 (watermark `2026-07-23T20:13:54Z`, unchanged). Inbox: 2 unread
  — both `eval-suite` controlplane status FYIs (suite started; 65/84 partial); not a World regression, not a
  directive, no cross-mission demand → did not outrank. No `[nightly-eval]` open issues on the world repo.
- Local `dev` == `origin/dev` == `f3c73c9` (in sync after `git fetch`). CI `CI` on dev: **completed/success** @
  `f3c73c9`. No weekly rotation (issue #1, created 2026-07-23, <80 comments; re-evaluate the Monday-07:00
  boundary next iteration).

**Pick + reality-check (Gate 2)**
- Picked top item `w-world-library-m1` M6 — the iter-15 "Next" explicitly named it. Fresh-origin already-landed
  check: `git log origin/dev --grep`, no `scripts/verify_go.sh`, no merged M6 PR → NOT landed, genuinely the
  final build step. Read the plan's M6 milestone spec (files, 3 acceptance_checks, verify_commands,
  ci_green_boundary) + inspected the EXISTING `ci.yml` (M2–M5 already added `ailang-verify` + `go-verify`; the
  go-verify job inlined `go build`/`go test` rather than calling a script) so the executor directive matched
  reality. Pinned binary `/tmp/ailang-v0300/ailang` present + clean (`v0.30.0`).

**Route + execute (Gate 3) — heavy roles model-PINNED, spawned (never inline)**
- Isolated worktree `git worktree add -b sprint/w-world-library-m1-m6 /tmp/wt-w-world-m1-m6 origin/dev` (from
  `f3c73c9`). NEVER the shared main tree.
- **sprint-executor (opus, MISSION_EXECUTOR_MODEL=opus, Agent-tool pin, worktree)** — directive carried the M6
  spec, the existing CI structure, the anti-false-green requirement, and NB2. Delivered `scripts/verify_go.sh`
  (durable local `go build && go test -count=1` gate with a loud-fail guard on unset/wrong-version `AILANG_BIN`),
  a `ci.yml` edit (go-verify job → single `./scripts/verify_go.sh` step, pinned-binary download + version-assert
  + Z3 + ailang-verify all retained), and `host/replay/replay_test.go` NB2
  (`TestInterpreterMemberChangeDrivesRealReplayEndToEnd` — full replay through a byte-distinct second interpreter,
  genuine end-to-end, env-constraint honestly documented). Checked the M1 acceptance boxes it verified.
- **Controller INDEPENDENTLY re-verified** in the worktree (data-before-conclusions): `go build` rc=0;
  `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` green with `ok …/host/replay 13.4s` (replay RAN,
  all 6 host pkgs `ok`); the anti-false-green guard fires (`env -u AILANG_BIN` → rc=1; `AILANG_BIN=/bin/echo`
  wrong-version → rc=1); `verify_ail.sh` = 4/4 world identities + 14 named tests / 8 modules; `go vet` clean;
  `gofmt -l scripts/ host/` empty; `actionlint ci.yml` rc=0. Read the NB2 test in full — it is a genuine
  end-to-end assertion (distinct HashRef, all six replay steps, cache-miss + authoritative-resolution +
  original-row-intact + faithful bytes), not a tautology. Also completed the remaining 5 M1 acceptance boxes the
  executor conservatively left unchecked — each verified against a landed, currently-green test (canon
  CRLF/CR/BOM/NUL/UTF-8/idempotence/golden; hashref reject malformed/uppercase/bare/tagged; store
  header/pair/registry-head/verif; registry epoch-1 candidate; archive HashRef/manifest/sidecar).
- **sprint-evaluator (sonnet, MISSION_EVALUATOR_MODEL=sonnet, Agent-tool pin; generator≠judge: opus ≠ sonnet)
  PASS 96/100 round 1**, ZERO blocking conditions — independently re-ran the gate, confirmed the guard fires,
  the replay tests RUN (not skip), the NB2 test is genuine end-to-end, the CI edit weakened nothing (pinned
  download + version-assert + Z3 + ailang-verify intact), and the M1 acceptance checkmarks are truthful. One
  non-blocking note: move the doc `planned/ → implemented/` on land (done in-PR).
- Controller moved the doc `planned/ → implemented/` in the SAME PR (atomic LAND; `verify_ail.sh` unaffected —
  it sweeps `world/` + sketches, not `design_docs/`). Committed on the worktree branch (`5888daf` M6 +
  `5c3f594` doc move), crediting the opus executor.

**Gate 3b — CI (item not LANDED until remote CI green on the merge)**
- Pushed branch; **PR #8** opened (base dev). Expected checks: both world CI jobs (`ailang-code verify gate`,
  `go host build + test gate`) — the workflow has no path filter, both run on the PR, no N/A. PR mergeable.
- Bounded poll (30-min cap): PR #8 CI **completed/success**, both jobs. **Verified the anti-false-green wiring
  WORKED in CI**: the go-verify log shows `AILANG_BIN=/home/runner/.local/bin/ailang (AILANG v0.30.0)` +
  `ok …/host/replay 9.021s` + `✓ go gate PASSED` — the replay tests ACTUALLY RAN via the new script (no SKIP),
  which is the whole point of M6. Squash-merged (`--delete-branch`) → dev `a07ac96`. Bounded poll (30-min cap)
  of dev CI on the merge commit `a07ac96`: **completed/success**, both jobs. Worktree removed; main checkout
  fast-forwarded to `a07ac96`.

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus (`claude-opus-4-8`) | triage/pick/reality-check/independent-verify/box-completion/record/retro |
| Design-doc-creator | — | not spawned | doc already quorum-cleared; no new/revised doc |
| Sprint-planner | — | not spawned | plan pre-existed (iter-8); M6 milestone spec used as-is |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | **opus** (Agent-tool pin, isolated worktree) | M6 (`5888daf`) + doc move (`5c3f594`); controller-reproduced |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | **sonnet** (Agent-tool pin) | PASS 96/100, zero blocking; generator≠judge (opus≠sonnet) held |

**Metered ledger**: `metered=$0.00` — executor opus + evaluator sonnet on subscription Agent-tool pins; no
designer/quorum/codex/gemini metered spend this iteration. Ceiling ($5) not approached. Designer rotation
UNCHANGED (`codex:gpt-5.6-sol` — no new doc authored).

**Ruled out** (do not re-chase)
- Re-running design-doc-creator/quorum/planner for M6 — doc already quorum-direction-accepted; mid-sprint final
  milestone, "Plan exists" lane.
- Deferring the doc `planned/→implemented/` move to a follow-up — done in the SAME PR so the LAND is atomic
  (the evaluator's non-blocking note; verify_ail.sh proven unaffected by the move).
- Proving that two SEMANTICALLY-DIVERGENT interpreter releases yield DIFFERENT replay bytes — genuinely
  env-constrained (needs ≥2 distinct upstream AILANG releases in the archive); documented in the NB2 test
  comment as upstream multi-release integration scope, NOT faked. The end-to-end interpreter-member replay
  itself IS verified (byte-distinct working wrapper).

**Retro / Next (Gate 5)**: No skill edit, no routing-policy change. **6th consecutive clean-landed sprint** on
opus-executor / sonnet-judge / generator≠judge (iter-8/9/10/14/15/16) — corroborates keeping the pattern. The
anti-false-green guard M6 adds is the 3rd touch of the silent-skip-in-CI class (V27 z3 iter-13, B1 AILANG_BIN
iter-15, now M6's verify_go.sh guard) — but each was caught by the gate and the pattern is a repo-specific
code/CI discipline (download the pinned artifact + assert + fail-loud in the job that needs it), NOT a
loop-process gap the skill could enforce → stays a pattern-watch note below the ≥2-instance skill-edit bar.
**Next iteration: `w-world-library-m1` is COMPLETE — pick the next open queue item, `w-worldd-m2`** (clause-2,
`ailang-worldd` local daemon: SQLite/REST/CLI, zero cloud deps, kernel perf budget measured from first commit).
This is a NEW-DOC item: it needs a design doc via the ROTATION designer (last-used `codex:gpt-5.6-sol` → next is
gemini after G4, else back to claude) + a mandatory `## Conflict Surface` vs the existing `ailang serve-api`
machinery (the iteration-0 quorum's standing gemini objection), then the pick-time quorum before routing. It is
a ~2d sprint-sized item — decompose into milestones at planning if needed.

---

## Iteration 17 — 2026-07-27 — `w-worldd-m2` (clause-2 local daemon) NEW-DOC authored + quorum-run (2 rounds) → **PARKED needs-human-review** on ONE ratification-class decision (single-writer enforcement)

**Kind**: NEW-DOC design iteration (design-doc-creator → pick-time quorum). Parked at the quorum
gate for a human decision; no sprint routed (Standing rule 2 — never force a guardrail).

**Context / preflight (Gate 0)**
- Kill switch `~/.ailang/state/mission-world.disabled`: NOT set (armed). Billing tripwire: **CLEAN**
  (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff` (`gh` not on
  the tool PATH — prepended `/opt/homebrew/bin` per call, per memory).
- Pidfile `mission-world.pid`=57443 = this run's own driver (my shell PPID; no overlap).
- Inbox: empty (`ailang messages list --unread` → none). No `[nightly-eval]` issues (only open
  issue is bookkeeping #9). No `@MarkEdmondson1234` comment on #9 (created today) or predecessor
  #1 since watermark `2026-07-23T20:13:54Z`. Weekly rotation already performed this morning
  (issue #9 "week of 2026-07-27", predecessor #1) — no rotation needed this iteration.

**Observe (Gate 1)**
- `git fetch origin`; local `dev` == `origin/dev` == `6fbbafd`; no missing commits. CI on dev:
  **completed/success** @ `6fbbafd36`. (Leftover remote branch `origin/sprint/w-world-library-m1`
  from the completed M1 sprint — harmless; not cleaned.)

**Pick + reality-check (Gate 2)**
- Queue head = item 3 `w-worldd-m2` (items 1–2 LANDED). NEW-DOC tag is a claim → verified: no
  `design_docs/planned|implemented/w-worldd-m2.md` (grep hits were queue/log references), not on
  `origin/dev` (`git log --grep`), no merged PR, no quorum artifact → **genuine NEW-DOC**.
- Charter's Conflict Surface marks worldd's placement "OPEN for ratification … Revisit only on
  concrete binary-distribution pain"; the coordinator-recommended DEFAULT is the in-repo Go module
  and there is no distribution pain → proceeds without a human gate (recorded as premise P1).

**Route + execute (Gate 3) — designer model-PINNED, spawned (never inline)**
- Designer rotation: last-used `codex:gpt-5.6-sol` → next in cycle is gemini, but gemini
  (managed_agents) is READ-ONLY (CapRemoteSandbox — file edits don't reach the worktree) so it
  cannot author a doc → wrap to `claude:claude-fable-5`. Fable probe via `claude-sub` (billing
  guard, bounded 120s): **rc=0, replied `ok`**.
- **design-doc-creator (claude:claude-fable-5, `claude-sub` backgrounded, bounded ≤30min,
  bypassPermissions)** — directive carried a full ADAPTING BRIEF (the skill assumes the upstream
  `sunholo-data/ailang` layout — known friction; this is ailang-world), all clause-2 constraints,
  the M1 host API surface it wraps, the day-1 perf-budget guardrail, S3 slim-kernel, §14, and the
  MANDATORY `## Conflict Surface` vs `ailang serve-api`. Produced
  `design_docs/planned/w-worldd-m2.md` (34 KB, 3 CI-green milestones ~2d) + checked sketch
  `design_docs/sketches/worlddapi.ail`.
- **Controller INDEPENDENTLY re-verified** (data before conclusions) on the **pinned
  `/tmp/ailang-v0300/ailang` (v0.30.0, clean)**: sketch `check.passed:true`, `verify {verified:2,
  counterexample:0, skipped:0, errors:0}` (isLoopbackHost, clampLimit), full
  `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → PASSED (9 modules, 4/4 required
  identities, 14 named tests). Confirmed the module-resolution path (`cd design_docs` first, per
  `verify_ail.sh`'s base-cwd logic) — a root-relative check falsely reports `LDR001 module not
  found: sketches/logepoch`; the gate's `cd $base` form is the truth.

**Pick-time quorum (Gate 2 quorum-at-pick) — 2 rounds, controller-verdict pass both**
- **r1** (`ailang design-quorum`, reviewers gpt5-6-sol + gemini-3-1-pro, $0.10/reviewer cap):
  **BLOCKED**, both present, metered **$0.0817**. Objections: (1) gpt5-6-sol — bounded-waits axiom:
  unbounded `http.Server`, unbounded client calls, unbounded graceful shutdown, no request-body
  limit on base64 commit payloads; (2) gemini-3-1-pro — the `GET /v1/log?from&limit` deliberate
  O(N) N+1 loop is omitted from the D6 perf harness, so the day-1 budget (P3/A9) hides the only
  deliberate N+1 latency. Both concrete, non-direction → routed to the designer.
- **designer revision r1** (same Fable lane, bounded ≤25min, rc=0): new **Decision 7 "Bounded Waits
  & Allocations"** (http.Server ReadHeader/Read/Write/Idle timeouts; bounded `Shutdown(ctx)`+hard
  close; client `context.WithTimeout` on every call; `POST /v1/commit` capped by
  `http.MaxBytesReader(maxCommitBytes=8 MiB)`→**413** with a new **Z3-proven `withinCommitBytes`**
  contract + `PayloadTooLarge` `ApiError` arm + non-vacuous `TestBoundedWaitsAndBodyLimit`) and
  **`BenchmarkLogRange`** (limit=100 + clamp-max 500) + `bench/BASELINE.md` N+1 rows (P3/A9/M2.*
  /Acceptance updated). Controller re-verified on pinned v0.30.0: `verified:3` (isLoopbackHost,
  clampLimit, withinCommitBytes), 0 counterexamples/errors/skips; full sweep still green.
- **re-quorum r2** (the ONE allowed re-quorum): **BLOCKED**, both present, metered **$0.104**.
  `gemini-3-1-pro` **PASS** (its note — make CLI `--addr` a global flag — is non-blocking).
  `gpt5-6-sol` **REJECT**, strongest objection: *"The single-writer guarantee is asserted but not
  enforced. `ailang-worldd` takes no database lease or process lock, `store.Open` remains available
  to embedded writers, and a second daemon can open the same SQLite file. An operational
  instruction to 'never' do so cannot support the claimed sole-handle model or A6 safe-concurrency
  conclusion."* Its `proposed_fix` is a FORK: (A) add a fail-closed store writer-lock
  (`OpenWriter`/`WriterAlreadyActive`, read-only path, subprocess tests) — a `host/store` change —
  **or** (B) withdraw the sole-writer claim + document/test bounded multi-process SQLite behavior.

**Decision — PARK (Gate 2 default; carve-out ruled out)**
- One revision + one re-quorum used (the gate's budget). Still-blocked → default is
  `needs-human-review`. The **narrow-refinement carve-out (iter-95) does NOT apply**: the objection
  (a) offers a FORK requiring controller JUDGMENT to choose an architecture (not a single verbatim
  fix), and (b) touches the LANDED M1 `host/store` — a **kernel change**, which the mission
  guardrail makes **ratification-class** ("kernel changes require explicit human ratification,
  quorum evidence attached"); arm (B) withdraws a load-bearing A6 claim (direction-level). Carve-out
  first-use also needs Mark's one-time OK — impossible headless. → **PARKED for @MarkEdmondson1234**
  with the A/B decision named in the doc's Park box + the queue tag.
- Doc + sketch committed to dev (doc-only, gate-green — preserves the design work + quorum evidence
  so the next iteration unparks fast). Not routed to sprint-planner.

**Gate 3b — CI**: doc-only + sketch push. The sketch enters `verify_ail.sh`'s CI sweep (already
green locally on pinned v0.30.0); no Go code changed. Bounded poll of the resulting dev CI run.

**Routing evidence** (role, model ACTUALLY used)
| Role | Pin (env) | Actual | Notes |
|---|---|---|---|
| Controller | `$MODEL` (session) | opus (`claude-opus-4-8`) | triage/pick/reality-check/independent-verify/quorum-orchestration/record/retro |
| Design-doc-creator | ROTATION (seed `claude:claude-fable-5`) | **`claude:claude-fable-5`** (`claude-sub`, backgrounded, bounded, subscription/quota-bucket) | authored doc + sketch + revision r1; rotation last-used advanced codex→claude (gemini skipped: read-only) |
| Quorum reviewers | fleet Phase B | **gpt5-6-sol + gemini-3-1-pro** (both present, metered) | r1 BLOCKED, r2 gemini PASS / gpt5-6-sol REJECT |
| Sprint-planner / executor / evaluator | — | not spawned | parked before routing |

**Metered ledger**: `metered=$0.186` (quorum r1 $0.0817 + r2 $0.1040; both under the $0.10/reviewer
cap). Designer Fable = subscription quota-bucket ($0). Ceiling ($5) not approached.

**Ruled out** (do not re-chase)
- The narrow-refinement carve-out for the r2 objection — it's a judgment fork touching the M1
  kernel (ratification-class), not a verbatim reviewer fix; forcing it would breach Standing rule 2
  + the kernel-ratification guardrail.
- A 2nd designer revision this iteration — the gate allows exactly one revision + one re-quorum;
  spending more without a human decision would re-litigate a contested (fork) objection.
- Worrying the daemon-placement "OPEN for ratification" item — the charter defers it to concrete
  binary-distribution pain, which does not exist (no `cmd/` ships today); recorded as premise P1.
- Root-relative `ai-check` of the sketch (false `LDR001`) — resolution is base-cwd-relative; the
  gate's `cd $base` form is authoritative.

**Retro / Next (Gate 5)**: No skill edit (no ≥2-instance recurring friction this iteration; the
design-doc-creator upstream-layout assumption is already the PROPOSED-to-Mark instance 2/2 in
memory — handled by the adapting-brief workaround, unchanged). No routing-policy change. The
opus-controller / Fable-designer / dual-provider-quorum path worked exactly as designed — the
quorum caught two real hardening gaps (bounded waits, hidden N+1) that a single-eye pass would have
shipped, and the revise→re-quorum loop resolved both; the third objection is a genuine
architecture decision correctly escalated, not a loop failure. **Next iteration**: if Mark answers
the A/B single-writer decision (comment on #9), UNPARK `w-worldd-m2` — apply the chosen arm as an
r3 designer revision (+ fold gemini's non-blocking `--addr` global-flag nit), re-verify, route to
sprint-planner. If no answer yet, `w-worldd-m2` stays parked and the queue head advances to
`w-effect-broker-m3` (clause-3, NEW-DOC) — but note w-effect-broker depends conceptually on the
daemon shell, so prefer waiting for the unpark unless Mark redirects.

---

## Iteration 18 — 2026-07-27 — `w-worldd-m2` UNPARKED on Mark's ratification: r3 revision applied + sprint planned + **M2.A/A1 LANDED** (PR #10 → squash `b0deedb`, dev CI green, both jobs) — the RATIFIED single-writer kernel change is now ENFORCED across processes

**Kind**: full-chain iteration — human-directive unpark → designer revision r3 (+ a controller-review
fix pass r3b) → sprint-planner → sprint-executor → sprint-evaluator → PR → CI-green merge.

**Context / preflight (Gate 0)**
- Kill switch `~/.ailang/state/mission-control.disabled`: NOT set (armed). Billing tripwire:
  **CLEAN** (no `ANTHROPIC_API_KEY`/`ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff`.
- Pidfile `mission-world.pid`=3308 = this run's own driver (verified via `ps`; no overlap).
- Inbox: empty. No `[nightly-eval]` issues (only open issue is bookkeeping #9, 3 comments).
- **HUMAN DIRECTIVE ACTED ON** — `@MarkEdmondson1234` on #9 @ `2026-07-27T08:55:11Z`, body `"A"`:
  the answer to iteration 17's parked A/B fork = **arm A, ENFORCE single-writer**. Already recorded
  in the charter by the attended session (`9578d1f`); this iteration is the unpark that acts on it.
  Watermark `mission-world-last-seen` advanced `2026-07-23T20:13:54Z` → `2026-07-27T08:55:11Z`
  (written BEFORE routing). Predecessor issue #1 re-checked (rotation-week catch): no new comments.
- Weekly rotation: NOT due — #9 was created today (2026-07-27 05:51Z, after the Monday-07:00
  boundary) and has <80 comments.

**Observe (Gate 1)**
- `git fetch origin`; local `dev` == `origin/dev` == `6fba10d`; CI on dev **completed/success**.
  Repo has exactly one workflow (`CI`, two jobs: `ailang-code verify gate` + `go host build + test
  gate`) — no path-filtered workflow to mis-poll (Gate-3b blind-poll class N/A here).

**Pick + reality-check (Gate 2)**
- Queue head = item 3 `w-worldd-m2`, tagged `[NEXT — ARM A RATIFIED]`. Verified against reality:
  no `cmd/` directory, no `worldd` Go code, no merged PR, no sprint plan → genuinely unstarted.
- Quorum-at-pick: already satisfied (2 artifacts, `.ailang/state/mission-quorum/w-worldd-m2-*.json`).
  The charter prescribes the unpark path explicitly — **apply the chosen arm as r3 → sprint-planner**
  — so no third quorum round was run (the human ratification supersedes the reviewer's objection;
  this is NOT the narrow-refinement carve-out and did not need it).

**Route + execute (Gate 3) — every heavy role model-PINNED and spawned, never inline**
- **Designer rotation** advanced `claude:claude-fable-5` → **`codex:gpt-5.6-sol`** (state file
  `mission-world-designer-rotation` written back after the run). Notably this is the SAME reviewer
  that raised the objection, now authoring its own ratified fix.
- **Designer r3** (codex, `--sandbox workspace-write`, bounded 30-min cap, rc=0, 63.6k tokens):
  rewrote **Decision 2** — `store.Open` = fail-closed writer path (canonical DB identity →
  non-waiting OS exclusive lock on `<canonical-db>.writer.lock` → only then the SQLite handle),
  structured `WriterAlreadyActive` on contention, additive `store.OpenReadOnly`, explicit
  crash/stale-lock recovery, cross-process subprocess proof — and carried it through Status/Park
  box, P6, High-Impact table, Design Freeze, D1 reuse table, **D5 (`--addr` promoted to a global
  client flag — the r2 `gemini-3-1-pro` non-blocking nit)**, M2.A/B/C (lock lands in M2.A; day
  split re-cut 1.0/0.6/0.4d, still ~2d), aggregate file table (~2,135 LOC), Conflict Surface vs M1
  embedded consumers, Acceptance Criteria, A4/A6/A10/A11, and a new r3 log subsection.
  API-shape choice (its own rationale): keep the landed `Open` NAME as the locked writer
  constructor + add `OpenReadOnly`, rather than a new `OpenWriter` + deprecation — no call-site
  churn, same invariant.
- **Controller review found ONE defect → bounded fix pass r3b** (same lane, 15-min cap, rc=0,
  30.4k tokens): r3 specified a lock file beside the canonicalized DB but left **non-file DSNs
  unspecified**. `Open(":memory:")` has TWO landed call sites (`host/store/store_test.go:13`,
  `host/registry/registry_test.go:12`) and `store.go:148` documents `:memory:` as supported — so
  as written it would have created a literal `./:memory:.writer.lock` and made the two in-process
  opens contend, **falsifying the doc's own "landed suites stay green unmodified" claim**. r3b adds
  the carve-out (in-memory DSNs detected pre-canonicalization, no lock, no lock file; the invariant
  is stated over file-backed paths), the matching M2.A acceptance checks, and an honest
  "controller-review correction" bullet in the r3 log. Committed as `5b105b5`.
- **sprint-planner** (`$MISSION_PLANNER_MODEL`=**opus**, pinned Agent sub-agent): wrote
  `.ailang/state/sprints/w-worldd-m2.{plan.json,handoff.md}`. Kept the doc's 3 milestones; added a
  top-level `ratified_kernel_change` block, per-milestone `non_vacuity_requirements`, and an
  internal **A1/A2/A3 split of M2.A** with A1 (`host/store` only) marked
  `safe_landing_point: true` — because M2.A's planned ~1,550 LOC exceeds M1's largest landed
  milestone (1,086). Velocity was measured, not invented: 7 landed PRs, and M1's doc estimate ran
  **2.2× low** (1,925 est vs 4,273 actual) → a 1.5× haircut on doc LOC estimates.
- **sprint-executor** (`$MISSION_EXECUTOR_MODEL`=**opus**, pinned Agent sub-agent, isolated
  worktree `/tmp/wt-worldd-m2-a1` on branch `sprint/w-worldd-m2-a1` — never the shared main tree):
  implemented **A1 only**. 1,017 insertions / 11 deletions across 5 files, all in `host/store`:
  `store.go` (+111/−11, constructor wiring only — every existing method body untouched),
  `writer_lock.go` (203, shared: `WriterAlreadyActive`/`IsWriterAlreadyActive`,
  `UnsupportedPlatformError`, `isInMemoryDSN`, `canonicalDBPath`, DSN rendering),
  `writer_lock_unix.go` (68, `//go:build unix`, `syscall.Flock(LOCK_EX|LOCK_NB)` from **stdlib**
  — chosen over `golang.org/x/sys` so M2.B's dependency-allowlist test stays simple),
  `writer_lock_other.go` (24, `//go:build !unix`, fails closed, honestly marked untested),
  `writer_lock_test.go` (622, the cross-process proof).

**Verification — controller INDEPENDENTLY re-ran everything (data before conclusions)**
- Both gates green in the worktree on the **pinned `/tmp/ailang-v0300/ailang` (v0.30.0, clean)**:
  `verify_go.sh` → build clean, all 6 host packages `ok` with **`host/replay` RUNNING 12.2s, not
  SKIP** (the iter-13 V27 / iter-15 B1 false-green class stayed closed); `verify_ail.sh` → **exactly**
  4/4 required `world/` identities across 9 modules + 14 named tests (A1 adds no `.ail`, so any
  movement would have been a red flag, not a success).
- Scope discipline verified by diff, not by claim: `host/store/schema.sql` **byte-for-byte
  unchanged**; `git diff dev..HEAD` over `world/ host/hashref host/canon host/archive host/registry
  host/replay scripts/ host/store/store_test.go` is **empty**; all five landed `store.Open` call
  sites green **unmodified**.
- **Controller's own mutation test** (the executor reported four; this is an independent fifth):
  `LOCK_EX` → `LOCK_SH` turns the suite RED with
  `Open(...) SUCCEEDED while OS process 23237 holds the writer lock: the ratified single-writer
  invariant is not enforced` — naming a live helper PID. An in-process mutex is structurally
  incapable of passing this suite. Tree restored, re-verified green.
- **sprint-evaluator** (`$MISSION_EVALUATOR_MODEL`=**sonnet**; generator≠judge holds: opus executor
  ≠ sonnet judge): **PASS 97/100, ZERO blocking merge conditions**. It re-ran both gates itself,
  checked all 12 claims adversarially, and judged all five of the executor's self-declared gaps
  acceptable. Three non-blocking carry-forwards → checkpoint A2.

**Gate 3b — CI GREEN (bounded polls, 30-min caps, no unbounded waits)**
- PR **#10** (`sprint/w-worldd-m2-a1` → `dev`): run `30256646182` **completed/success**, both jobs
  (`ailang-code verify gate` 9s, `go host build + test gate` 28s). `mergeable: MERGEABLE`.
- Squash-merged → **`b0deedb`** on dev, branch deleted. Dev-merge run `30256701072`
  **completed/success**, both jobs. Only an OBSERVED green upgraded the tag. Worktree removed.

**Routing evidence** (role, model ACTUALLY used — the enforcement backstop)

| Role | Env pin | ACTUALLY ran on | Notes |
|---|---|---|---|
| Controller | session | `claude-opus-5` | triage/pick/review/mutation-test/record/retro |
| Design-doc-creator (r3 + r3b) | ROTATION | **`codex:gpt-5.6-sol`** | rotation advanced from fable; probe rc=0; both runs bounded (30m/15m) and rc=0 |
| Sprint-planner | `MISSION_PLANNER_MODEL`=opus | **opus** | pinned Agent sub-agent |
| Sprint-executor | `MISSION_EXECUTOR_MODEL`=opus | **opus** | pinned Agent sub-agent, isolated worktree |
| Sprint-evaluator | `MISSION_EVALUATOR_MODEL`=sonnet | **sonnet** | generator≠judge satisfied (opus ≠ sonnet) |

No role fell back off its pin; no FLAGged degradation.

**Metered ledger**: `metered=$0.00` — and this is a CHANGE worth recording. Mid-iteration an
attended session flipped the codex lane to **ChatGPT subscription auth** (`~/.codex/auth.json`
`auth_mode=chatgpt`; the metered API key moved aside to `auth.json.apikey.bak` at 11:33 local,
concurrent with this iteration's codex probe). So both designer runs (63.6k + 30.4k tokens) were a
**quota lane, not metered dollars**. Honest residual uncertainty: the r3 run launched within ~2
minutes of that flip, so if it raced the change it would have billed the API key — worst case ≈
**$0.23** by a GPT-5-class rate estimate. Either way far under the `$5` ceiling. No quorum calls,
no gemini/managed_agents calls this iteration.

**Ruled out** (do not re-chase)
- **A third quorum round on `w-worldd-m2`** — ruled out deliberately: the blocking objection was
  resolved by HUMAN RATIFICATION, which outranks the reviewer verdict; the charter's unpark
  instruction says r3 → sprint-planner, full stop. Re-quorumming a ratified decision would be
  re-litigating a settled human call.
- **`OpenWriter` as a new constructor + deprecating `Open`** — considered by the designer and
  rejected: it forces churn at all five landed call sites without strengthening the invariant.
  Recorded here so a later iteration does not "fix" the naming back.
- **Unlinking the lock file on release** — rejected in-design: it races another process that has
  already opened the file. flock ownership is per open-file-description, so the leftover *pathname*
  is inert and the SIGKILL-recovery test proves it cannot wedge the DB.
- **An in-process mutex / package-level held-paths map** — this would be a DEFECT, not a shortcut;
  it is exactly what the quorum rejected. Verified impossible-to-ship by the controller's own
  `LOCK_EX`→`LOCK_SH` mutation.
- **Blaming the codex lane for the first probe failure** — refuted: `rc=127 env: node: No such
  file or directory` was a PATH gap in the tool shell, not a codex/quota problem. See the retro.

**Concurrency note (Critical Principle 0 — recorded, not acted on)**: partway through this
iteration `tools/launchd/mission-control.sh` (FROZEN CORE, shared driver) appeared MODIFIED and
uncommitted in the main checkout — an attended/sibling session flipping
`MISSION_EXECUTOR_MODEL`'s default to `codex:gpt-5.6-sol` and documenting the subscription-auth
change. The tree was clean at Gate 0, so this landed mid-iteration. It was **left completely
untouched**: not stashed, not committed, not reverted. Local `dev` was reconciled to `origin/dev`
with `git reset --mixed` + a scoped `git checkout -- host/store/` precisely so the sibling's
uncommitted work survived; `git status` afterwards shows that one file and nothing else.

**Retro / Next (Gate 5)**: ONE process fix, no skill edit — see the charter's Repo Profile
addition. Friction was a **3-instance PATH class in a single iteration**: `codex` (`rc=127 env:
node: not found`), `go` (`verify_go.sh: line 36: go: command not found`), and `gh` (already in
memory from earlier iterations) are all under `/opt/homebrew/bin`, which this tool shell's PATH
omits. That is ≥2 recorded instances pointing at one gap, and it is a MISSION-DOC (process) fix,
not a skill fix, because the PATH is a property of this rig's Repo Profile rather than of the
shared skill. No routing-policy change (7th consecutive clean landed sprint corroborates
opus-executor / sonnet-judge / generator≠judge; the codex designer lane produced a doc revision
that survived independent review with exactly one narrow defect, caught and fixed inside the same
iteration).
**Next**: `w-worldd-m2` checkpoint **A2** — `cmd/ailang-worldd` + `host/daemon` shell (config,
loopback guard, D7 bound constants + the four `http.Server` timeouts, bounded shutdown,
`/v1/health`, `/v1/head`), folding the evaluator's three non-blocking carry-forwards; then A3
(bench harness + `bench/BASELINE.md` + `scripts/bench_worldd.sh --smoke` + the CI bench-smoke step)
completes M2.A; then M2.B, M2.C.

## Iteration 19 — 2026-07-27 — `w-worldd-m2` **M2.A/A2 LANDED** (PR #11 → squash `39b2115`, dev CI green, both jobs): the `ailang-worldd` daemon shell — loopback guard, D7 bounds, `/v1/health` + `/v1/head`, and zero-cloud ENFORCED by a dependency allowlist

**Kind**: mid-sprint EXECUTE iteration ("Plan exists" lane) — no new design doc, no quorum, no
planner. Executor → evaluator → **round-2 fix on a BLOCKING finding** → re-verify → PR → CI-green
merge.

**Context / preflight (Gate 0)**
- Kill switch: NOT set (armed). Billing tripwire: **CLEAN** (no `ANTHROPIC_API_KEY` /
  `ANTHROPIC_AUTH_TOKEN`). gh account: `sunholo-voight-kampff`.
- Pidfile `mission-world.pid`=35412 = this run's own driver (verified via `ps`; no overlap).
- Inbox: 3 unread, all `[controlplane] eval-suite` status FYIs from the V1 loop — not directives,
  not cross-mission requests, do not outrank the queue. No `[nightly-eval]` issues (the only open
  issue is bookkeeping #9).
- **No new `@MarkEdmondson1234` comment** on #9 (watermark `2026-07-27T08:55:11Z` = his "A"
  ratification, already acted on in iter-18). Predecessor #1 re-checked per the rotation-week
  catch: nothing new. Watermark unchanged — nothing new to advance past.
- Weekly rotation: NOT due — #9 created 2026-07-27 05:51Z (after the Monday-07:00 boundary), 4
  comments.

**Observe (Gate 1)**
- `git fetch origin`; local `dev` == `origin/dev` == `b1e2b33`. CI on dev **completed/success**.
  One workflow (`CI`, two jobs) — no path-filtered workflow to mis-poll.

**Pick + reality-check (Gate 2)**
- Queue head = item 3 `w-worldd-m2` **[IN-SPRINT]**; iter-18's recorded **Next** is checkpoint A2.
- Already-landed check against a FRESH origin: `cmd/` did not exist, `host/daemon` did not exist,
  no PR mentioned the daemon. A2 genuinely unstarted. No quorum needed (design doc r3 is
  quorum-run + human-ratified; this is the execute lane).
- Baseline established by the controller BEFORE handing the worktree over: `verify_go.sh` green,
  6 host packages `ok`, `host/replay` RUNNING 12.673s. So a red gate would be the executor's
  change, not a pre-existing condition.

**ROUTING INCIDENT — the pinned executor model does not exist for this CLI (Gate 3)**
- `$MISSION_EXECUTOR_MODEL` = `codex:gpt-5.6-sol` (the default flipped to codex by an attended
  session, `dd12587`, for quota relief). Per the skill's cross-provider recipe the controller ran
  the bounded pre-flight probe **with the pinned model** — it FAILED, and NOT on quota:
  `400 invalid_request_error: The 'gpt-5.6-sol' model requires a newer version of Codex. Please
  upgrade to the latest app or CLI` (codex-cli **0.137.0**).
- **Diagnosis, not assumption**: a second bounded probe WITHOUT `--model` returned rc=0 on default
  **`gpt-5.5`**. So the codex lane itself is healthy and on ChatGPT subscription auth
  (`auth_mode=chatgpt`); it is the MODEL PIN that is unreachable.
- **This means the driver's own pre-flight probe false-greens a model pin.**
  `tools/launchd/mission-control.sh:248` probes `codex exec --skip-git-repo-check 'reply with
  exactly: ok'` — **no `--model` flag** — so it exercises the default model and reports the lane
  healthy while the pinned model cannot run. The driver logged
  `executor=codex:gpt-5.6-sol` for this fire on exactly that false signal. The shared SKILL's
  recipe is correct (it passes `--model "$MODEL"`); the DRIVER's probe is the one missing it.
  `tools/launchd/*` is **FROZEN CORE** for this mission → NOT edited locally; routed upstream as a
  proposal (two channels, per the guardrail).
- **Action taken**: per the recipe's fallback rule, the executor role fell back to `$MODEL` (opus)
  via a pinned Agent sub-agent, **FLAGGED**. Deliberately NOT substituted with `codex:gpt-5.5` —
  swapping in an unratified model is a routing-policy change, which the charter gates behind
  evidence, not controller convenience. generator≠judge still holds (opus executor ≠ sonnet judge).

**Execute (Gate 3) — sprint-executor, isolated worktree `/tmp/wt-worldd-m2-a2`**
- Round 1 (`879558d`): 4 new files, **1548 insertions, 0 deletions** — `cmd/ailang-worldd/main.go`
  (233), `cmd/ailang-worldd/main_test.go` (129), `host/daemon/daemon.go` (522),
  `host/daemon/daemon_test.go` (664). LOC came in **2.1× the plan estimate** (~730 planned),
  consistent with M1's measured 2.2×-low calibration — the haircut in the plan is still too small.
- Declared deviations, both accepted: a 4th file (`main_test.go`) because the exit-code contract
  and the "`--addr` is not a serve flag" rule were stated acceptance properties with no test
  anywhere in A2/A3; and a `drainTimeout` field on `Daemon` instead of referencing `shutdownTimeout`
  directly — its own mutation testing showed the direct reference was **unfalsifiable** (raising the
  constant to 100h stayed green). Making an assertion falsifiable is the right call.

**Verification — controller re-ran everything INDEPENDENTLY (data before conclusions)**
- Both gates green on the pinned `/tmp/ailang-v0300/ailang` (v0.30.0, clean, `e37b370`): all 8
  packages `ok` with **`host/replay` RUNNING 12.1–12.3s, never SKIP** (the V27/B1 false-green class
  stayed closed); `verify_ail.sh` at **EXACTLY** 4/4 required `world/` identities across 9 modules
  and **EXACTLY** 14 named tests (A2 adds no `.ail`, so any movement would be a red flag, not a
  success). `gofmt`/`go vet` clean.
- Scope verified **by diff, not by claim**: `host/store`, `host/replay`, `host/registry`,
  `host/archive`, `host/canon`, `host/hashref`, `world/`, `scripts/`, `.github/`, `go.mod`, `go.sum`
  all byte-unchanged vs `origin/dev`. No new module dependency.
- **Controller mutation 1** (independent of the executor's): widening `isLoopbackHost` to also
  accept `0.0.0.0` — a SUBTLE widening, not an always-true — turned the suite RED at **two**
  independent tests (`TestIsLoopbackHostMirrorsSketchPredicate` and
  `TestNewRefusesNonLoopbackBind/refused/0.0.0.0`) in 1.2s. Tree restored byte-clean.
- **Controller mutation 2**: dropping `golang.org/x/sys` from the dependency allowlist turned it RED
  naming `golang.org/x/sys/unix` — proving the allowlist walks the REAL build graph, not a synthetic
  list. Tree restored byte-clean.
- **LIVE, with real OS processes** (the highest-value check — A1's ratified kernel invariant meeting
  its first real consumer): built the binary; proc 1 announced
  `ailang-worldd listening on http://127.0.0.1:54819`; `/v1/health` returned the real archived
  interpreter HashRef `sha256:e9746fef…` plus the pinned `AILANG v0.30.0` manifest version;
  **a SECOND OS PROCESS on the same DB failed closed in 0s, rc=2**, with `WriterAlreadyActive`
  surfaced as `StartupError{Stage: store-open}`; `--bind 0.0.0.0` refused rc=2 with a named error;
  SIGTERM exited cleanly.

**EVALUATION — the independent judge found a real BLOCKING defect (this is the loop working)**
- **sprint-evaluator** (`$MISSION_EVALUATOR_MODEL`=**sonnet**; generator≠judge holds: opus executor
  ≠ sonnet judge): round 1 **PASS-WITH-CONDITIONS, 79/100, ONE BLOCKING condition**.
- **BLOCK-1**: `TestDaemonDependencyAllowlist` was **not delivered**. Without it, Decision 4's
  "zero-cloud is enforced, not asserted" and the charter's **local-first is inviolable** guardrail
  were prose with no gate behind them. The executor had claimed the plan JSON scoped the test to
  M2.B.
- **The controller verified the evaluator's finding rather than adopting it** — and found the
  evaluator's own framing slightly overstated. Its claim "all three documents place it in A2" is not
  exactly right: the **plan JSON is internally inconsistent** — it says "the **M2.B**
  dependency-allowlist test" twice (in two risk texts), says "`TestDaemonDependencyAllowlist`,
  **M2.A**" once, and **omits the test from A2's file list entirely**. The executor's reading was
  therefore defensible, not careless. The **design doc governs** (it is the quorum-reviewed,
  ratified artifact) and places the test in M2.A four times over: Decision 4, the M2.A
  `daemon_test.go` file description, the M2.A acceptance checks, and a standalone global acceptance
  criterion. The handoff agrees (line 149). **Verdict: the test belongs in A2 and was dropped.**
- **Round 2** (`a1cc5fa`, bounded single-file fix pass, executor lane): `TestDaemonDependencyAllowlist`
  added to `host/daemon/daemon_test.go` (+206). It shells out to `go list -deps` over BOTH
  `./host/daemon/...` and `./cmd/ailang-worldd/...` (**237 transitive packages** today), resolves
  `go` via `exec.LookPath` and **`t.Fatalf`s rather than `t.Skip`s** when the toolchain is missing
  (a skip here would be precisely the V27/B1 silent-skip class), **fails on an empty dependency
  list** (S6 null case), and **names the offending packages** so CI output is actionable.
- **Round-2 verification by the SAME independent judge** (narrowly scoped to BLOCK-1 — the
  controller is opus, i.e. the executor's own model family, so controller self-verification alone
  would have been weaker evidence than it looks): **BLOCK-1 CLOSED**, **revised score 92/100**, no
  new blocking conditions. It ran its own distinct mutation (removing `modernc.org/sqlite`, which
  neither the controller nor the executor had mutated) → RED naming three packages by name; and
  confirmed the synthetic self-proof subtests are *in addition to* the real-graph walk, not instead
  of it.

**Gate 3b — CI GREEN (bounded polls, 30-min caps, no unbounded waits)**
- PR **#11** (`sprint/w-worldd-m2-a2` → `dev`): run `30266163765` **completed/success**, both jobs
  (`ailang-code verify gate`, `go host build + test gate`).
- Squash-merged → **`39b2115`** on dev, branch deleted. Dev-merge run `30266239289`
  **completed/success**, both jobs. Only an OBSERVED green upgraded the tag. Worktree removed;
  local `dev` fast-forwarded to `39b2115`.

**Routing evidence** (role, model ACTUALLY used — the enforcement backstop)

| Role | Env pin | ACTUALLY ran on | Notes |
|---|---|---|---|
| Controller | session | `claude-opus-5` | triage/pick/baseline/2 mutation tests/live process proof/review/record/retro |
| Design-doc-creator | — | **not invoked** | execute lane; doc r3 already quorum-run + ratified |
| Sprint-planner | — | **not invoked** | plan already approved in iter-18 |
| Sprint-executor | `codex:gpt-5.6-sol` | **opus (FALLBACK, FLAGGED)** | pinned model rejected by codex-cli 0.137.0 (`requires a newer version of Codex`) — NOT quota. Recipe's fallback rule applied; 2 rounds (r1 + bounded fix) |
| Sprint-evaluator | `sonnet` | **sonnet** | generator≠judge satisfied (opus ≠ sonnet); 2 rounds (full + scoped BLOCK-1 re-verify) |

One role fell back off its pin — the executor — **FLAGGED above and reported to Mark**.

**Metered ledger**: `metered=$0.00`. The codex probes were the ChatGPT **subscription** lane
(`auth_mode=chatgpt`) and were ~1 reply-token each; the controller additionally invoked codex with
`env -u OPENAI_API_KEY` so an ambient metered key could not silently bill (guard the CALL SITE, the
`claude-sub` discipline applied to the OpenAI lane). Executor + evaluator ran on subscription
Agent-tool pins. No quorum calls, no gemini/managed_agents calls. Far under the `$5` ceiling.

**Ruled out** (do not re-chase)
- **Substituting `codex:gpt-5.5` for the failed `gpt-5.6-sol` pin** — deliberately rejected. The
  lane is healthy on the default model, so this was tempting and would have preserved Mark's
  quota-relief intent. But swapping in an unratified model is a routing-policy change, and the
  charter gates those behind an evidence rule, not controller convenience. The skill's documented
  fallback is `$MODEL`. Recorded so a later iteration does not treat this as an oversight.
- **Blaming quota for the codex probe failure** — REFUTED by a second probe: the default model
  answered rc=0 on the same auth. The failure is a CLI/model-availability mismatch. (Note the shape
  of the iter-18 scar repeating: a lane failure that superficially reads as "quota spent" but is
  something else entirely. Diagnose before concluding.)
- **Treating the plan JSON as authoritative over the design doc** — the plan contradicts itself on
  where `TestDaemonDependencyAllowlist` lives. The quorum-reviewed design doc governs; the plan is a
  derived artifact.
- **Landing A3 in this iteration** — not attempted. A2 alone ran 2.1× its LOC estimate and needed a
  fix round; the plan explicitly authorises incremental checkpoint landing, and iter-18 set the
  precedent by landing A1 alone.
- **An in-process test for `WriterAlreadyActive` at the daemon layer** — correctly refused by the
  executor as the forbidden anti-pattern; the controller proved it with two real OS processes
  instead. Belongs in M2.C's subprocess e2e (CF-2).

**Non-blocking carry-forwards — ENUMERATED, because iteration 18 recorded "three carry-forwards"
without naming them and they were LOST** (the fix: the evaluator directive now REQUIRES an explicit
numbered list):
1. **CF-1** — `TestBoundedWaitsAndBodyLimit` part (b), 413-on-oversized-body. Correctly absent:
   `POST /v1/commit` is M2.B. → **M2.B** (`handlers_test.go`).
2. **CF-2** — no in-process test that the daemon surfaces `WriterAlreadyActive` as
   `StartupError{Stage: store-open}`; an in-process test would violate the ratified cross-process
   requirement. Controller verified live with two real OS processes. → **M2.C** subprocess e2e.
3. **CF-3** — the non-unix writer-lock arm (`writer_lock_other.go`) is unexercised by CI. Inherited
   from A1, honestly documented; dev=darwin/arm64 and CI=ubuntu are both unix. → **no change
   needed**, acknowledged design choice.
4. **CF-4** — `releaseFromVersion` divergence hazard: a DB bootstrapped WITH `--ailang-bin` and later
   served WITHOUT it yields a different epoch-1 revision → fatal startup. That IS the specified
   never-silent behaviour, but it is an operator sharp edge. → **M2.C close-out** (or a doc decision
   if a friendlier rule is wanted — it is a design question, not a code tweak).
5. **CF-5** — `/v1/head` error bodies are `text/plain`, not the sketch's JSON `ApiError` envelope.
   Status codes already follow `httpStatus`. → **M2.B** handlers.
6. **CF-6** — `TestHealthAndHeadRoundTrip` skips on Windows (POSIX shell-script fake interpreter),
   mirroring `host/archive`'s existing convention; CI is ubuntu-only so it never skips in the gate.
   → **no change needed**.
7. **CF-7** — `TestNewBootstrapsEpochRegistryIdempotently` never exercises the `--ailang-bin` path,
   so a broken `releaseFromVersion` is caught only by its own unit test. → **A3 or M2.C**.
8. **CF-8** — A3 items (`bench_test.go`, `scripts/bench_worldd.sh`, `bench/BASELINE.md`, the CI
   bench-smoke step) deferred by design. → **A3** (the next checkpoint).
9. **CF-9** — `isStdlibImportPath` lives in `_test.go` so it cannot be imported/fuzzed externally;
   its unit vectors are the only documentation of the intended contract (the existing
   `vendor/golang.org/x/net/...` vector already guards the main regression). → **informational, no
   action required**.

**Retro / Next (Gate 5)** — **no skill edit** (no gap reached the ≥2-instance bar; and the one
skill-adjacent finding is a DRIVER bug, not a skill bug — the skill's recipe is correct). **Two
process fixes to the charter's Repo Profile**, both first-instance-but-costly-to-rediscover:
(1) the codex model-pin reality (`gpt-5.6-sol` unreachable on codex-cli 0.137.0; the driver's probe
omits `--model` and therefore false-greens a pin), and (2) a standing requirement that the
evaluator's non-blocking carry-forwards be ENUMERATED in the log entry — iteration 18's three were
recorded only as a count and are unrecoverable. **Upstream proposal routed on two channels** (per
the frozen-core guardrail): the driver probe fix belongs to the shared infrastructure, so it goes to
the V1 loop as an issue + an `ailang messages` note — never a local edit of `tools/launchd/*`.
No routing-policy change (the executor fallback is the documented rule operating correctly, not
evidence for a new policy).
**Next**: `w-worldd-m2` checkpoint **A3** — bench harness (`BenchmarkStoreCommit`,
`BenchmarkHeadRead`, `BenchmarkHealth` with p50/p95 via `b.ReportMetric`), `bench/BASELINE.md` (the
committed day-1 budget; REST-commit + log-range rows marked explicitly PENDING M2.B), a
`scripts/bench_worldd.sh --smoke` that fails on a MISSING BENCHMARK NAME (not a zero line count —
`go test -bench` exits 0 on no-match, the V27/B1 class), and the CI bench-smoke step. That completes
M2.A; then M2.B, M2.C.

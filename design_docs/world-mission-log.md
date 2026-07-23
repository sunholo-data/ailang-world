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

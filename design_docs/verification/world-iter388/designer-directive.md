MISSION-ROLE: designer

You are the DESIGNER for Ailang World mission iteration 388 (unattended; no human is present).
Your one deliverable is a NEW child design document:

    /Users/voightkampff/dev/sunholo-data/.wt-world-iter388-a4/design_docs/planned/w-world-approvals-by-key-a4.md

Work ONLY inside the worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter388-a4`
(branch `design/world-iter388-row201-a4`, base `b0fbcd3eee7b5b29190e678e296635a3dd3a0d0c`).
Write that one file. Do NOT edit any other tracked file, do NOT write Go/AILANG code, do NOT run
any git write (no add/commit/push/branch/stash/reset/checkout). The controller commits.
Scratch files go under `/Users/voightkampff/.ailang/state/world-iter388/designer-scratch/` only.

## What A4 is

Queue row 201 (`w-world-approvals-by-key`, an ATTENDED parent design, rulings D-WORLD-89 R1–R5 all
"yes" by Mark) has milestones A0–A4. A0–A3a are LANDED on `dev`: mint requests with canonical
bytes/digest/challenge/expiry; SSHSIG verification + `approver enroll|list|revoke` + `approve`
(`--signature`, `--sign-with`); WebAuthn passkeys (`--webauthn-assertion`, `/surface/passkey`);
the surface flow (`/surface/requests.json`, the card `/surface/approve/{id}`, `POST
/surface/approve/{id}`, the shared judge `host/authority/judge.go`, `approve <id> --await`, Needs
you / Decided / notification in `app.js`). The parent's whole A4 text is one line:

> **A4:** `Human.Approve` decision packets (row 172) approved the same way.

and parent §2.2: "A **mint request**, and in S2 any `Human.Approve` decision packet
(`world/decision-packet/v1`, row 172), carries: the request digest …; a one-time challenge …; an
expiry … The signed message is a fixed-prefix encoding of `("world-approval/v1", request digest,
challenge, expiry)`." Parent §6 R1: approval may be a verified signature "for **mints and
`Human.Approve`**". R4: **publishing stays typed-only**.

Design A4: how an operator decides (approve / deny) a pending `Human.Approve` request that an
agent's invocation raised inside a daemon episode, using an enrolled key (SSHSIG and WebAuthn,
the same enrolled approvers, the same judge) and the typed `y` as root — through the CLI and the
surface, reusing A0–A3a machinery rather than building a parallel one.

## Measured anchors (provenance labelled — re-check any you rely on)

VERIFIED BY CONTROLLER at `b0fbcd3e` (commands run from the worktree):
- `grep -rn "DecideApproval(\|decideApproval(" --include="*.go" . | grep -v _test.go` →
  only `host/broker/approve.go:181` (def), `:188`, `:191` (def), `:696` (inside
  `RecordOperatorApproval`) and `host/broker/publish_op.go:153` (inside `mintAttendedApproval`).
  Control in the same breath: 5 `_test.go` files call `DecideApproval(`. So the exported operator
  entry point `broker.DecideApproval` has NO production caller.
- `RecordOperatorApproval` has one production caller: `host/daemon/episode_finish.go:173`
  (`"operator:episode-finish"`, decidedBy `"operator:terminal"`), which creates AND decides its
  own request in one call. `MintAttendedApproval` has one: `cmd/world-publish/main.go:310`
  (publish; typed-only by R4; `cmd/world-publish/tty.go` explains the `/dev/tty` fence and says
  "a headless loop can mint one by calling broker.DecideApproval").
- Therefore (controller inference — CONFIRM OR REFUTE IT FIRST, as your V1): a decision packet
  produced on the daemon's normal path (row 172 producer `host/coordinator/decisionpacket.go`,
  bound by the daemon since row 176, `d.binder`) shows `approvals chain: pending`
  (docs/QUICKSTART.md ~line 1603) and no operator command or route can decide it today.
- `validatePublishApproval` is at `host/broker/approve.go:517`; publish consumption needs an
  `approve` decision and a publish-shaped scope (D-WORLD-73's text, V44 of the row-172 design).

UNVERIFIED — inherited from earlier iteration records; re-check before relying:
- The shared judge, `surface_approvals` (credentials schema 7), `MintRequestIDPattern`
  `^mr-[0-9a-f]{16}$` (`host/store/mint_requests.go:105`), `DecidedMintRequests`, `surfaceCSP`
  unchanged since A2, the 90-byte signed message, SSHSIG namespace `world-approval@sunholo`.

## Read before writing (in this order)

1. `CLAUDE.md`, `design_docs/coding-standards.md` (binding: pure core / effects at the host
   boundary; slim kernel; non-vacuous gates), `design_docs/DESIGN.md` §1 and §14.
2. The parent `design_docs/planned/w-world-approvals-by-key.md` (93 lines — all of it).
3. `design_docs/planned/w-world-approvals-by-key-a3.md` r2 — your STRUCTURAL TEMPLATE (its §0
   Amendment posture, §1 measured-today, §2 design, §3 pure core, §4 threats, §5 by-construction,
   §6 test plan with mutants, §7 inventories, §8 conflict surface, §9 milestones/sprints/gates
   with a fallback cut, §10 attended exit run, §11–§13, §14 Verification Log, §15 quorum log).
   Skim A0a1 / A1b / A2 child docs only for the request/judge/approver facts you need.
4. `design_docs/planned/w-decision-packet-production-projection.md` (row 172) and its sprint plan
   §C4.2 (row 176), the charter rows 172/176 (`grep -n '^17[26]\. ' design_docs/world-mission.md`)
   and ledger rows D-WORLD-69, D-WORLD-73, D-WORLD-89, D-WORLD-109, D-WORLD-114.
5. Code: `host/broker/approve.go`, `host/broker/approvals_view.go`, `host/coordinator/
   decisionpacket*.go`, `host/coordinator/effectful.go`, `host/daemon/workbench_packet.go`,
   `host/authority/judge.go`, `host/store/mint_requests.go`, `host/store/surface_approvals.go`,
   `cmd/ailang-worldd/` approve/approver commands, `host/surfacepage/` (the card + app.js).

## Questions the design must answer (with measurements, not assertions)

- What object is signed for a `Human.Approve` decision: the approval request ref, the packet, or
  both — and what "the exact request" means here (effect, scope, cost, requester, the proposal
  plan, the evidence links the packet carries). How the parent's `("world-approval/v1", digest,
  challenge, expiry)` message is reused WITHOUT letting a mint-request signature be replayed as a
  `Human.Approve` decision or vice versa (domain separation — measure what the current message
  encoding binds and say exactly what changes, if anything).
- Where the one-time challenge and expiry live for a `Human.Approve` request (today's approval
  request object has none); how consumption is single-use and durable across restart; how this
  interacts with the immutable approvals head chain and the coordinator's serialization
  (`serializedHandler`, row 176), including the row-176 judge residual (a conflict-losing
  invocation leaves a pending request with no committed packet).
- **Deny** as well as approve, by key and by typed `y`. Who decides is authority-derived
  (`ssh:<name>`, `webauthn:<name>`, `operator:terminal`), never caller-supplied.
- **R4 interplay (load-bearing):** a `Human.Approve` request whose scope is publish-shaped must
  NOT become key-approvable — publishing stays typed-only. Measure how `validatePublishApproval`
  and `world-publish` consume an approve decision and show by construction that no key-signed A4
  decision can satisfy a publish. If the clean answer needs a refusal on publish-shaped scopes,
  design it and its test.
- Surfaces: CLI (`ailang-worldd approve …` for a `Human.Approve` request — name the syntax) and the
  surface (Needs you lists it; a card renders the packet's fields as text; Approve/Deny with
  passkey or SSHSIG through the existing POST pattern and judge). Say what is reused verbatim and
  what is new. Keep `surfaceCSP` unchanged unless you measure a need.
- What the agent observes (Human.PollApproval) after a key decision, and that nothing about an
  agent's grants lets it reach any decide path (enrolment typed-only; no key).

## Hard rules for the document

(a) **Every sentence claiming the codebase currently does X needs a §14 Verification Log row with
the exact command AND its observed output.** An empty or negative result is a CLAIM, not a fact:
pair it with a known-positive control in the same command. Cite file:line.
(b) Any browser-capability premise (WebAuthn, fetch, Notification, CSP behaviour) gets its own
V-row measured in this fire, or is explicitly marked as an obligation for the planner to measure —
never "known". (Lesson from the A3 quorum, round 1.)
(c) **Any departure from the attended parent** gets an *Amendment posture* paragraph in §0 AND a
proposed one-word ledger ask (A/B with a default) in §12 — in r0, not later. (Lesson from A3,
round 2.) If you need none, say so explicitly.
(d) Sprint-sized: the build must split into sprints of ≤ ~2,000 changed lines each (prefer one
or two small, fully proven sprints), with a priced fallback cut in §9. Name packages and
approximate line counts.
(e) Test plan: per sprint, test rows each naming the mutation it kills; refusal branches each get
their own mutant (unit = the branch); include must-die and must-live controls; name the exact
`go test -run` patterns the gates will run so no test is silently skipped.
(f) Pure core vs effects: the decision law (what is accepted/refused) is a pure function with
tests; IO only at the host boundary. No new Go dependencies (`TestDaemonDependencyAllowlist`).
(g) An attended exit run section (§10) for Mark — the loop never runs it (D-WORLD-109: build and
judge against fakes and test keys only).
(h) A §15 "Quorum verification log" stub for the controller to fill.
(i) No marketing language; plain, short sentences; measured numbers carry their instrument.

## Process rules (you run headless; nothing will re-invoke you)

- Do NOT background anything. If you start any long command, run it in the foreground with a
  timeout. A `Monitor` is an event stream, not a wait; nothing will re-invoke you when it fires,
  so never end your turn announcing that you are waiting. If something you need cannot finish,
  record it as UNMEASURED in §14 and finish the document.
- Budget: aim to finish within ~30 minutes. A complete r0 with honest UNMEASURED rows beats an
  unfinished one.
- Use absolute paths or `cd <abs> &&` in the same command (the shell does not keep cwd).

## Your final reply (the controller reads only this)

1. The file path and its line count (`wc -l`).
2. The number of §14 V-rows, and any marked UNMEASURED.
3. A 10-line summary: the shape, the sprints with line estimates, the R4 answer, any amendment to
   the parent and its proposed ledger ask.
4. Any anchor above that you found FALSE.

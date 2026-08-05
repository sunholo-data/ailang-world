# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-05, iteration 51 · dev @ `269f1fe` · CI green both jobs (SHA-addressed)

## In flight

- **Item 8 `w-self-mod-vertical` — DOC LANDED, not yet sprint-planned.** PR #40 → `269f1fe`,
  839 lines, designer `codex:gpt-5.6-sol`. Milestones **SM.A–SM.D**, **4–5 d**. `[NEXT]` is its
  sprint-planner run, gated on nothing. **SM.A–SM.C routable today**; only SM.D waits on `8/OD-1`.
- **Planner prices these two first** (carried, not cleared — round 2 had one reviewer, round 3 was
  the carve-out): the design needs a `schema.sql` change, and the landed `w-ddl-gate-teeth` DDL gate
  reds on *any* schema edit **by design**, so its fixture update belongs in the same milestone; and
  whether 4–5 d is one queue item or splits at SM.B.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite: the transition registry is
  absent at HEAD (measured iter-50, control fired). Unchanged this iteration.

## Latest

- Iter-51 headline: the item's binding VERIFY-FIRST clause returned a fact that **reframes it**.
  There is no vendor namespace to claim — `registry-validator/main.go:177` says
  `// Step 5: Namespace auth — deferred (accept all publishers for now)` — and a live 4-arm dry-run
  accepts `world/`, `someoneelse/` **and `sunholo/`** alike, against a firing control. *`world/` is
  a string World writes, not a namespace World holds.* Publish is immutable, unrecallable by the
  publisher, and the key is **ambient in this loop's shells**.
- Quorum blocked twice. Round 1's defect was **mine**: I appended an `approve.go` evidence row after
  the body was written, so the doc's evidence contradicted its own design *by construction*.
  Round 3's carve-out ran the reviewer's own prescribed check and the answer was neither
  right-nor-wrong — the metadata path is a GCS **bucket key**, not a validator route. A bare 404
  there re-authorizes an irreversible POST, so absence now needs a same-pass known-positive control.

## Loop

- Cadence: launchd, `mission-world`. Controller `claude-opus-5`.
- Routing: designer rotation advanced slot 1 → **slot 2 `codex:gpt-5.6-sol`** (pointer written back;
  next new-doc iteration returns to `claude:claude-fable-5`). Executor `codex:gpt-5.6-sol` ·
  evaluator `sonnet` · planner `opus`. Planner/executor/evaluator **not fired** — a design
  iteration has no implementation to judge.
- Verify profile `ailang-code`; AILANG pinned **v0.30.0** at `/tmp/ailang-v0300/ailang`
  (`e9746fef…`); upstream source read only at `e37b370d…`.
- Bookkeeping issue **#32** (week of 2026-08-03).

## Cost

- Iteration 51 **`metered≈$0.96`** against the $5 ceiling: quorum **$0.140** measured exactly
  (r1 $0.103 + r2 $0.037), codex designer **≈$0.82 ESTIMATED** — the CLI reports tokens
  (525k in / 15.5k out, then 40.9k) and **not dollars**, so that half of the ledger is an estimate
  wearing a number. Controller on the opus quota bucket.

## Parked on Mark

- **`8/OD-1`** — attended stamp authorizing the **irreversible first public publish** of
  `world/core@0.1.0`. Controller default: **do not publish**; the loop stops at
  `READY_AWAITING_HUMAN_PUBLISH`. Blocks **SM.D only**. Next free OD number: **`OD-9`**.

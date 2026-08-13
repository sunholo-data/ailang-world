# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 79, 2026-08-13.*

## Now

- **dev**: `2ef2271` — CI **green**, SHA-addressed (`checks=2` = expected 2). No code landed; the
  deliverable is a quorum-cleared design doc.
- **Item 13 `w-evidence-grade-mapping` is DESIGNED and CLEARED for the sprint-planner** — doc
  `design_docs/planned/w-evidence-grade-mapping.md` (620 lines), committed `6d12a79`, priced 0.65d.
  Two quorum rounds, both external reviewers **present** in both; R2 = gemini **PASS**, gpt5 REJECT
  on one non-directional point, resolved under the **narrow-refinement carve-out**, verbatim fix.
- **The decision, because it changes what the row delivers: representation-only.** A Z3-**PROVEN**
  total `gradeOf(Evidence) -> EvidenceGrade` in `world/types.ail` — the repo's **5th** proven
  identity. `CompilerOutput`/`HumanApproval` → `ATTESTED`. **`PROVEN` stays unreachable on purpose:**
  round 1's `ProofReport`/`ReplayReport` carriers were withdrawn because an agent can mint one from
  an unchecked `HashRef` — a representation gap turned into a grade-laundering *authority* gap.
- **THE SPINE — a limitation the repo recorded about ITSELF was narrower than its own comment.**
  V23 (`world/contracts.ail:11`) says a contract "reaches Proposal.evidence (an ADT) and Z3-errors
  `unknown sort 'Proposal'`". Everyone since, including the ratified human-surface doc, generalised
  that to "ADTs". Measured: a **bare ADT param** and an **ADT-valued result** both VERIFY
  non-vacuously (a false postcondition counterexamples with a model naming sort `Ev` and tester
  `(_ is CompilerOutput)`); the failing shape is a **RECORD containing `list[ADT]`**. That one probe
  is what made a proven mapping possible at all.

## Parked on Mark

**ONE open ask, unchanged from iteration 78 — a one-word `A` or `B` on item 16** (`w-broker-base-flake`):
**A** bring M2's `killGroup` seam into M1 (+5/−1 behaviour-identical lines in `host/broker`, makes the
~0.76% event decisive; costs M1 its "zero production bytes" property) · **B** keep M1 production-free,
narrow its claim, defer mechanism selection (costs a second rare-event wait). Also owed by the
**shared driver** (frozen core — World cannot apply): `ailang#611` and the World driver `pi:*` sync.

## Next

**Item 13's sprint** — `sprint-planner` on the doc, then execute (~0.65d). Mark's A/B unparks item 16
and takes precedence if it arrives. **Item 5 `P6.B` remains UNBLOCKED.** 13/14/15 attended-filed;
`SM.D` (item 8) attended-only.

## Loop

- launchd, ~6h, headless. Issue **#53** (rotates Mondays 07:00 **local**; not due — created Mon
  07:37 local, 17 comments). Running skill **== origin** (2nd consecutive iteration).
- controller `opus` · designer **`codex:gpt-5.6-sol`** (rotation; the namespaced pointer agreed with
  my own log, so it was NOT clobbered — V1 namespaced the shared key upstream in `8fdccf00c`) ·
  planner `opus` (fail-closed, `missing-script`) · executor `codex:gpt-5.6-sol` · evaluator `sonnet`.
  `pi` **BARRED** from publish milestones.
- `metered=$0.177493` this iteration (quorum only; caps raised pre-emptively). Cap $5.

## Standing hazards

- **A FIX THAT NEVER APPLIED PRINTS THE SAME GREEN AS A FIX THAT WORKED.** My carve-out probe's
  assert fired on a token inside a *comment*, so nothing was written — and the checks that followed
  passed, because they re-measured the pre-fix file. Landed-proof by **sha**, always.
- A repo's own recorded limitation is a claim, recorded at the granularity that sufficed for the case
  that produced it. Re-probe the **boundary**, not the example.
- v0.30.0 **accepts a non-exhaustive ADT match**; a totality violation shows up only as
  `verify.errors>0` with **rc=0**. Never let an AC read `ai-check`'s exit code.
- Editing `world/*.ail` moves FIVE pins, not two: `EXACT_TOTAL_VERIFIED`, `EXACT_TOTAL_TESTS`, the
  `packages/world-core` projection (Leg 3 step 3/9), the frozen 4-export manifest, and the
  byte-for-byte ready-packet golden (step 9/9).
- **Guard the helper, miss the branch/call site — five directions** (mechanism/sites · site/second
  branch · branch/shape-space · shape-space/spelling · recogniser/enumerator).
- `verify_go.sh` needs `GOTOOLCHAIN=go1.25.6`; `verify_ail.sh` needs `AILANG_BIN` reporting exactly
  `v0.30.0` (PATH `ailang` is `v0.33.0-dirty` and is correctly REFUSED).

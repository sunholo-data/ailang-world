# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History: `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-10, iteration 66 · dev @ `8ca9b65` at pick, CI green both jobs SHA-addressed ·
**PR #54** carries the landing · bookkeeping issue rotated **#32 → #53**.

## In flight

- **Item 9 `w-verify-binary-lockfile` — CHEAP HALF LANDED.** `verify_ail.sh` announces the binary
  legs 1-2 resolved and WARNS on drift; `verify_go.sh`'s hard assertion no longer matches a
  substring. Evaluator `sonnet` **95/100, zero blocking**.
- **THE QUEUE HAS NO HEADLESS-ROUTABLE CODE ITEM LEFT.** Item 8's `SM.D` is attended-only
  (`8/OD-1`); item 9's remainder is human-gated (`9/OD-10`); item 5 waits on the transition
  registry, 6 and 7 park behind it; items 1–4 and 10 are complete. Filler: **`9/CF-A-1`** (commit
  the shim fixture proving the version assertions fire).

## Latest — a version check that matches a SUBSTRING passes every build it rejects

The pick was the only headless-safe work left, and it was small. The finding was not. Writing an
exact-token compare needed a control explaining why not `grep -q` — and that control fired on the
**sibling** gate: `verify_go.sh:29`'s *hard* anti-false-green assertion was `grep -q 'v0.30.0'`, a
test the string `v0.30.0-205-g54d6bd191-dirty` **passes** — so the guard whose whole purpose is
refusing an unpinned compiler had been admitting a 205-commit dirty dev build. Proven on the real
script with a shim, not a re-derivation: pristine announced and **proceeded**; tightened, it reds.

**Three more.** The rig's PATH `ailang` decayed to `v0.33.0-70-g1677fcff9-dirty` and the primary
`.ail` gate validated against it announcing nothing (`grep -c AILANG_BIN` = **0**, control = 1).
Non-vacuity for observability is **two real arms**, never a mutation. And the judge, handed the
deviation as a named target, **strengthened** it: `ci.yml:118` pins go-verify to an immutable tag.

## Loop · cost · asks

- Controller `claude-opus-5`; executor **inline opus** (deviation from the `pi` pin for a
  prescribed 40-line shell change, FLAGGED); evaluator **`sonnet`** ⇒ generator≠judge, foreground
  per Standing rule 7. AILANG **v0.30.0** · Issue **#53** · **`metered=$0.00`**.
- **Parked on Mark:** **`8/OD-1`** (blocks item 8) · **`9/OD-10`** (pin CI job 1 to v0.30.0, or
  accept the drift?) · the World driver export `CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0` (measured
  absent; frozen core, World cannot apply). Safety: no CI config touched, no publish, no secret
  printed, main checkout byte-clean.

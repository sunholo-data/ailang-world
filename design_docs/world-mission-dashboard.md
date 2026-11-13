# World mission dashboard — iteration 392 (2026-11-13)

- **Done this fire:** row 201 A4 sprint 4 MERGED ([#571](https://github.com/sunholo-data/ailang-world/pull/571), merge `842b9331`). Judge (sonnet) PASS 88/100, zero blocking. A4 is code-complete.
  - *Needs you* lists each pending agent decision and links it to its card; a new decidable one raises the notification.
  - The card (`/surface/decide/<ref>`) shows both signed objects, checks the hashes itself, and decides with a passkey or a pasted SSH signature. It is proven in live Chrome for approve, deny, 16 tampered subjects and a real `ssh-keygen`.
  - `approve --list` shows the same pending decisions, read offline from the world store.
  - 1,947 of 2,000 lines, nothing cut.
- **Dev CI is RED (verify gate only):** AILANG v0.60.4 (released 08:23Z) type-checks contract clauses, and three kernel `ensures` use `==` on types with no `Eq` (`world/types.ail`, `world/transitions.ail`). Not from this sprint (no `.ail` changed). Both Go gates are green on the merge. Open PR #565 fixes only the enum half. Row 260 + **D-WORLD-116** filed.
- **The record PR for this fire stays open** until the verify gate is green (nothing lands red); the next fire reads it from `mission/world-iter392-record`.
- **Next loop work:** row 260, the contract fix-forward (default A of D-WORLD-116), then row 243 (address-bar verbs, new design r0).
- **Residuals named for Mark:** RS-VERB (the CLI `--webauthn-link` path opens the frozen A2 passkey page, which says "Approve request …" for a deny; the card shows the right verb); RS-RIG (the card's Chrome tests post to a rig, not a seeded daemon; the attended run's step 4 is the first real one).
- **Attended, pending Mark:** exit runs for A1b (QUICKSTART §15), A2 (§16), A3a (§17), A4 (§18 / design §10, incl. the card on a phone).
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (default A, applied), D-WORLD-115 (default A, applied), D-WORLD-116 (new; default A from iteration 393).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:** planner opus (Agent, 93 min), executor haiku (Agent ×4; one stop on a plan snapshot-command bug, controller-corrected), evaluator sonnet (Agent, 67 min). No designer. Codex unusable at preflight (probe rc 75). Metered $0.00.
- **Instrument notes:**
  - Merge race gate red once on `TestBackoff` (host/surfacefollow, untouched): instance 2; green on re-run.
  - The judge's `approve --list`-with-zero-mint-requests test gap was reproduced by the controller (code correct; logged, not filed).
  - The self-notice script still reports `watermark channel unavailable` until a directive is marked on #524.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).

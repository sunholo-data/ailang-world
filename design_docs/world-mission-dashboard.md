# World mission dashboard — iteration 391 (2026-11-13)

- **Done this fire:** row 201 A4 sprint 3 LANDED ([#569](https://github.com/sunholo-data/ailang-world/pull/569), merge `e6761873`, CI 3/3 on head and merge, both first try). Judge (sonnet) PASS 91/100, zero blocking.
  - An operator can decide a `Human.Approve` packet with a passkey: `approve <request-ref> [--deny] --webauthn-link` prints a link and the message hash; after the device ceremony, `--webauthn-assertion <file> --expiry <printed>` posts it. The daemon checks the counter inside the credentials writer's transaction.
  - Out-of-band SSH decisions work too: `--write-message` writes the 91-byte message, and `--signature <file> --expiry <printed>` posts it.
  - A live-Chrome test drives the passkey decision with Chrome's own bytes; QUICKSTART §18 documents it.
  - The whole surface (decision listing, *Needs you*, the card) moved to a sprint 4. That goes further than the design's fallback, and the controller accepted it. The sprint is 1,493 of 2,000 lines.
- **Known defect, for sprint 4:** the A2 passkey page reads "Approve request …" even for a deny. The signature binds the verb, and a mismatched post is refused with 422.
- **Safety note:** while planning, one test draft made a read-only health GET to the rig's real World service (nothing written). The landed test refuses to run unless its temp service record resolves to the fake daemon.
- **Next loop work:** the A4 sprint-4 planner — the surface (≈ 1,250 lines).
- **Attended, pending Mark:** the A1b exit run (QUICKSTART §15), the A2 real-passkey exit run (§16), the A3a exit run (§17); A4 design §10 steps 1–3 and 5 are runnable now.
- **OPEN decisions:** D-WORLD-105 (default C), D-WORLD-114 (default A, applied), D-WORLD-115 (default A, applied).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 and 151 are not routable.
- **Routing this fire:** planner opus (Agent, 57 min), executor haiku (Agent ×3, all clean), evaluator sonnet (Agent, 37 min). No designer. Codex was unusable at preflight (probe rc 75). Metered $0.00.
- **Instrument notes:**
  - The §6 Chrome leg went red once on an untouched A2 timing test (`TestPasskeyPageAbortsAPendingCeremony/get-late`); the leg passed 7/7 twice on re-run.
  - The judge's one surviving-mutant claim was refuted at package scope by the controller.
  - The self-notice script still reports `watermark channel unavailable` until a directive is marked on #524.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).

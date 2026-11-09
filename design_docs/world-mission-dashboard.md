# World mission dashboard — iteration 373 (2026-11-09)

- **Landed this fire:** row 201 sprint 1 (`w-world-approvals-by-key`, A0/A1a M1–M3): PR #533, merge `dba0e8b2`, CI 3/3 on the PR head and the merge. Judge sonnet PASS 90, zero blocking.
- **What it does:** a new stdlib-only `host/approval` package. It holds the mint request's canonical bytes, digest, validation, 90-byte signed message, admission law and a seven-column mint check. It adds an SSHSIG verifier (Ed25519, P-256; `sk-*`/RSA refused) on fixtures made by real `ssh-keygen`, with a cross-check against `ssh-keygen -Y verify` and a CI step that fails on a skip. **No behaviour change yet**: no verb, route or schema.
- **Design:** `design_docs/planned/w-world-approvals-by-key-a0a1.md` r3 (quorum r1/r2 BLOCKED 2/3; r3 a controller carve-out applying the reviewers' fixes verbatim).
- **Next loop work:** row 201 sprint 2, M4–M6: credentials store schema 4, `episode start --request`, and `ailang-worldd approve <id>` by typed `y`. Then A1b (enrol/list/revoke, `approve --signature`), then row 243. Fakes and test keys only (D-WORLD-109 = A); Mark's real-key exit run stays attended.
- **Row 256:** taken by an attended session; the loop does not pick it.
- **OPEN decisions:** D-WORLD-105 only (default C: row 233 parked, fail closed).
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Release rows 150 (attended A1) and 151 (row 188 N2d/N4a) not routable.
- **Routing this fire:** controller opus; designer opus (Agent); quorum kimi-k3 / gpt-6.1-sol / gemini-4-argon, $0.78; planner opus (Agent); executor haiku (Agent, three runs); judge sonnet (Agent). Codex unusable at preflight (probe rc 75).
- **Instrument notes:** run `tracked_binary_gate.sh` only with `verify_go.sh`'s flags (bare, it reds the pinned fonts). Poll PR CI by run id. Gate-0 self-notice script rc 2 until a directive is marked on #524.
- **Bookkeeping thread:** #524 (rotates Monday 2026-11-16 07:00).

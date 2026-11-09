# Iteration 373 — judge summary (row 201 sprint 1)

Judge: Agent `sonnet` (MISSION-ROLE: evaluator), on a `git archive` copy of 9c8f8744 under
`~/.ailang/state/world-iter373-judge`. Generator families: designer/planner opus, executor haiku.

**PASS 90/100, zero blocking.** Correctness 28/30 · test strength 25/30 · fixtures 15/15 ·
standards and gates 14/15 · scope and hygiene 8/10 (the judge did not re-run full `go test ./...`
or audit every one of the 18 files; the controller's `verify_go.sh` rc 0 and CI 3/3 cover that).

Mutation: 56 run, 48 killed. Surviving, logged as non-blocking (D-WORLD-109: logged, not filed):
1. sha256 hash path in `signedData` is unguarded (every fixture and the cross-check use
   ssh-keygen's default sha512). Controller reproduced: `hashalg` appears in 0 test files.
   The judge's own 120 `-O hashalg=sha256` signatures verified correctly.
2. Three strictness branches have no killing test: trailing bytes after the signature string,
   trailing bytes after the P-256 mpint pair, `base64.StdEncoding.Strict()` vs non-strict.
3. Parse-level `ErrHashAlg` is not asserted directly (verify still refuses sha1 via `signedData`).
4. `DisallowUnknownFields` in `Decode` (equivalent: the canonical re-encode catches the input).
5. Ed25519 `len(Sig) != 64` guard (equivalent: `ed25519.Verify` returns false).

Independent agreement: 240 fresh ssh-keygen signatures (Ed25519 and P-256 × sha256/sha512, a
message with NUL and 0x01 bytes) all verified; every tampered message refused; P-384 and RSA-2048
key lines refused with `ErrKeyType`.

SG4 deviation adjudicated for the executor: a doubled armor block passes the BEGIN/END frame and
fails on the inner armor lines as base64, so `ErrBase64` is the correct typed error; the plan's
`ErrArmor` expectation was over-specified.

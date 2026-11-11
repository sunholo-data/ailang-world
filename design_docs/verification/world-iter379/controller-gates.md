# Controller gates — iteration 379 (row 201 A2 sprint 1)

Sweep = plan §6, run as `gates.sh` in one bash shell at `ca2605d7` (before the judge), AILANG_BIN v0.41.0. **rc 0.**

- verify_ail.sh: ✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass
- gofmt (cmd/ailang-worldd, host/approval, host/sshsign): empty
- go build ./... && go vet ./...: rc 0
- compile fence `go test -run '^$' ./...`: rc 0
- touched packages: 4 ok, 0 FAIL
- -race host/approval/... host/sshsign: 3 ok, DATA RACE lines 0
- -race cmd/ailang-worldd (TestApproverRevoke|TestApproverEnroll): 1 ok
- the 16 named repo-wide inventories with -v: 16 PASS, 0 SKIP/FAIL
- VERIFY_GO_LEGS=plain verify_go.sh: 62 ok lines, 0 FAIL

Commits rebuilt from executor snapshots with commit_ms.sh; numstat vs base `03d78b6a` (cmd host): M0 75, M1 905, M2 1510, M3a 1958 — every value equals the plan's. Final tree `diff -r cmd host` IDENTICAL to the planner's round-trip tree (`.planner-scratch-iter379/roundtrip`).

Judge finding UV-OR-BE reproduced by the controller (mutant `if !ad.UV && !ad.BE` → TestNewPasskeyApprover* rc 0), then pinned by test-only commit `5a4b1224` (3 arms). Controller mutation proof on that tree: `!UV && !BE` rc 1, `!UV && !BS` rc 1, `!UP && !BE` rc 1; production restored byte-identical; host/approval/... rc 0. Numstat after: 1968.

Quorum premise checks (rule 3f), first-party:
- r1 gpt6-1-sol: doc lines 511–521 derive the 40-bit code only from the registration file → secrecy claim false. REAL.
- r1 gemini-4-argon: doc line 833 left blob: save under surfaceCSP UNVERIFIED. REAL.
- r1 oc-kimi-k3: no log row read the parent's ruling text; 16-byte floor UNVERIFIED. REAL (parent lines 32/42/49/66/77 do carry the quotes).
- r2 oc-kimi-k3: Chrome 157 headless and Node WHATWG URL both serialise `https://x.ts.net:443/` → `https://x.ts.net` (control :7644 kept). REAL.
- r2 gpt6-1-sol: r1 text had `timeout: 240000` and no AbortController/signal. REAL.

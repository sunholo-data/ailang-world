# Evaluator — iteration 379 (row 201 A2 sprint 1)

Judge: sonnet via the Agent tool, own detached worktree `.wt-world-iter379-judge` (generator ≠ judge: designer/planner opus, executor haiku).

**Round 1 at `ca2605d7`: PASS 91/100, zero blocking.** Correctness 28, tests 26, security 18, standards 9, scope 10.
- 77 mutants (44 from plan §4, 28+ own, controls): 71 killed, 5 survived (must-live control; CD-EMPTY equivalent; UV-OR-BE; COSE-ALIAS; CRED-ALIAS), 1 no-compile (re-run as RESTSKIP2, killed). Must-die control MUT-K9 killed; must-live MUT-Q1-EQUIV survived.
- Per-milestone non-vacuity: M0, M1, M2, M3a each stubbed in its own production code → its headline tests fail.
- Gates: vet rc 0, compile fence clean, touched packages ok, -race approval/sshsign ok, 16/16 inventories PASS, FuzzDecodeCBOR 20 s PASS.
- Non-blocking: (1) UV-OR-BE survivor — UV pinned only with BE=0; (2) retained-buffer aliasing of `Passkey.Key`/`CredentialID` not pinned (copies today); (3) no Ed25519 point/small-order check; (4) `SOURCE.txt` names `es256.json`/`eddsa.json`, files are `chrome-*.json`; (5) `RegistrationInput.Challenge` zero / `Now` unchecked (caller's job, sprint 2/CLI).
- The judge's own numstat read 906/1511 at M1/M2 (counting `.gitattributes`); the controller's `-- cmd host` reading is 905/1510, the plan's unit.

**Round 2 at `5a4b1224` (test-only UV/UP pin): confirmed, PASS 93/100** (tests 26 → 28). UV-OR-BE, UV-OR-BS, UV-OR-(BE||BS), UP-OR-BE, UP-OR-BS, UP-OR-UV all KILLED; must-die control killed; must-live survived. Worktree clean.

Carried to sprint 2: findings (2) aliasing pins, (4) SOURCE.txt names, (5) zero-challenge check; (3) is design-optional.

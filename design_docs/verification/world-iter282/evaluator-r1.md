# Evaluator round 1 — Agent sonnet, tree `.eval-world-iter282` @ e4810dc

**FAIL 68/100.** Categories: acceptance 22/30 · test quality 10/20 · design fidelity 19/20 · plan/JSON 9/10 · code/standards/gates 8/10.

- **BLOCKING 1** — `TestProtectedCandidateParseFailsClosed/over-cap-{status,ignored,flags}` used `strings.Repeat("x", cap+1)` with no trailing NUL, so they died on the missing-NUL check, not the cap; removing the `len(b) > protectedCheckMaxOutput` block left all three PASSing; a cap-doubled mutant also survived. (Controller reproduced: all three PASS with the cap check deleted.)
- NB2 — §4.5 cost sentence (≈165–195 ms/phase) stale against V23's real-checkout 406–504 ms warm / 2 114 ms cold ignored leg.
- NB3 — one-character widened patterns (`s`, `?`, `é`, `[ab]`) render `:(top,glob,icase)**/*`; correct (mirror filters) and bounded (8 MiB cap → error), but undisclosed.
- NB4 (opinion) — a status stream with no `# branch.oid` returns head "" and nil error; M6's head≠base comparison fails safe.
- NB5 — the 2 500 ms / 3 200 ms bounds are the controller's application of the reviewers' fixes; oc-glm-5-3's split-budget wording was not taken literally; disclosed and deferred to V24.
- Superset argument derived independently from `foldProtectedKey`/`matchProtected` and Git glob semantics: no hole found; Unicode semantics preserved; all original fixture rows retained.
- Own mutants: 18, 16 killed, 2 survived (both the cap). Gates: gofmt empty, `go vet ./...` rc0, `-run '^TestProtected'` 155 PASS / 1 SKIP (pinned arm) / 0 FAIL, `-race` ok, `jq -e` ok.

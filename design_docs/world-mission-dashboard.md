# World mission dashboard — iteration 289 (2026-10-20)

- **This fire:** row 156 V25 measured and M5e built, PR #330. The v0.54.15 tool pin (World's production write-tool binary) already refuses every symlink and hard-link alias into a protected path, where v0.52.1 wrote them all. M5e pins that through World's real handler on the real binary.
- **Scope:** no production change. The design gained V25 (evidence only). The plan gained M5e and holds M5c/M5d (World's own pre-check) on D-WORLD-86.
- **Side findings:** v0.54.15 refuses ANY hard-linked file through `ailang-write`/`ailang-edit`, even with no deny list (R-156-10, a property of the pin). `**` is doublestar on v0.54.15.
- **Judge:** Agent sonnet, own worktree, re-ran V25 itself: r1 PASS 90, zero blocking; wording fixes `5ab76bca`.
- **Verification:** outside the sandbox: real-pin 42/42 PASS ×3, discovered SKIP unset, `go vet ./...` and daemon `-race` rc 0. Drills: the retired v0.52.1 tool reds the 7 alias arms; a nil protected policy reds the 5 symlink arms.
- **Release:** `v1.0.0-rc.1` pre-release. Proof/compiler pin 0.41.0; tool pin 0.54.15.
- **Clauses:** 7/7 MET; goal unmoved. Before v1.0.0 (D-WORLD-81): row 156 (16/24, M5c/M5d held), row 150 (9/17), then row 151.
- **Next:** row 156 M6-pre (V24 real-worktree timing), then M6a.
- **Parked on Mark:** D-WORLD-86 — drop M5c/M5d (A, recommended) or build them anyway (B). Default: hold.
- **Routing/cadence:** planner codex gpt-6.1-sol, executor codex gpt-6.1-sol, evaluator Agent sonnet. No designer: evidence-only amendment.
- **Cost:** metered $0.00. codex ≈228k tok (planner ≈141k, executor ≈87k); Agent subagent ≈101k tok.
- **Retro:** harness share 0/20; no drift alarm. Create a multi-arm fixture's precondition objects only in the arms that need them. Measure every `-run` subtest selector's count.
- **Record:** world-mission-log.md entry 289; design_docs/verification/world-iter289/.

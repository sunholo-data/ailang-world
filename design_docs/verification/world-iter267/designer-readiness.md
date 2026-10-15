# Iteration 267 designer readiness — row 150 M1b only

Verdict: READY for bounded M1b execution. No substantive design blocker found in this milestone. This is a readiness review of existing r4, not a new design or a quorum pass. R1–r3 remain BLOCKED; r4's existing narrow-refinement routing and M1a landing are recorded mission state. No release, mint, publish, ledger resolution or harness authorization is supplied here.

Reviewed base: `57941bf7410c7462d755cd3a94748eaf76d66288` (origin/dev). Authoritative design-doc-creator skill read at `/Users/voightkampff/dev/sunholo-data/ailang/.claude/skills/design-doc-creator/SKILL.md`; binding CLAUDE.md, charter guardrails/row150 state, coding standards and DESIGN §§1/14 read. Existing `design_docs/planned/w-marketing-capture.md` r4 and its sprint plan reviewed. Language reference loaded through released `ailang prompt` (rc0, 94510 output bytes). No source edits or git writes performed.

## Exact scope and criteria

Allowlisted implementation files, exactly as sprint M1b:

- `packages/world-demo/approve.ail`
- `packages/world-demo/transitions.json`
- `packages/world-demo/README.md`
- `host/coordinator/world_demo_test.go`

`TestWorldDemoApprovalExactOnePlan` must execute the real v0.41.0 module with convention-v2 input and decode exactly one `e1`, `Human.Approve`, `worktree`, cost1, requester `demo`, `finish=false`, `result=null`. Descriptor access cost0 and declared cost1 must agree. Unknown keys/phases and malformed roster must refuse with zero effects; require the named `oneApprovalEffect` proof identity and independently pinned happy/refusal inline identities. Roster is four unique kind/ref entries with no caller grade; approved kinds are test-invocation, recorded-effect, authored-claim, proof-envelope. M1b checks shape, not future host evidence authenticity.

The module must be standalone (std imports only); publisher checks one staged source. If package law is repeated, extend executable drift coverage to compare approve and evidence sources. README must document payload construction, convention-v2 plan/finish behavior and bounded commands. No attended confirmation may be exercised. Mutation `M1b-two-effects` must remain valid typed JSON/source and fail the named `len(plan.effects)==1` assertion; a downstream L1 error is not that mutation's acceptance.

Required focused commands remain the plan's approval test, existing multi-effect refusal test, approve check/test/ai-check from its package directory, legacy drift test and `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh`, plus protected before/after observations and the named mutation. No empty selection or skip counts as acceptance. Future approve source/test/descriptor are absent at this base and were not claimed passing.

## Independently remeasured facts and controls

Environment: PATH begins `/opt/homebrew/bin`; AILANG_BIN is `$HOME/.pinned-ailang/ailang`. Live banner: `AILANG v0.41.0`, commit `24ee1088776e21cd06a3781ed18e77f40be06db3`. Go focused command:

```sh
go test -run '^Test(WorldDemoEvidencePinnedManifest|WorldDemoLegacyGradeDrift|ProducerMintsNoPacketForMultiEffectPlan|ProducerMintsPacketForApprovalOnly|DecisionPacketCodec)$' ./host/coordinator -count=1 -timeout 120s -v
```

rc0; all five top-level identities PASS, package 1.602s. Positive and negative controls are actual test bodies, not inferred names:

- M1a manifest: three named verified identities, 21 named inline cases, zero errors/counterexamples/skips; healthy-looking empty verification/test reports refused by the parser.
- M1a drift: executes frozen kernel constructors, including passing and failing TestReport; authored-claim maps CLAIMED while AiReview impersonation refuses; raw ProofReceipt stays CLAIMED, never PROVEN.
- Approval producer: successful single approval yields one packet; denied, non-approval and zero-effect refusal yield none.
- Multi-effect control: explicitly expects `violates L1` and zero packets. This establishes existing downstream confinement, not the new approval module's exact-one assertion.
- Codec: round trip passes; wrong hash, unknown members/policy, missing member, trailing value, malformed ref and invalid deadline refuse.

Existing production convention-v2 reference independently checked using 40s subprocess bounds, from `packages/se-tools`: `ailang check se_tools/read.ail` rc0; `ailang test --format json se_tools/read.ail` rc0, success=true, 38 passes, zero skips; main_test_1 through main_test_10 PASS. Read source bodies establish unknown argument/top-level key, wrong phase/type and finish refusal controls. An initial check from repository root produced MOD010 due to source-root/module-path mismatch; corrected package-root invocation passes. This was an instrument invocation error, not a language blocker. Initial reporting code expected an absent JSON `function` field; corrected inspection uses named test identities.

Code mechanism read: `cmd/world-publish/transitions.go:298` canonicalizes source then calls EnsureSourceLoadable before PutObject; `host/transitionreg/verify.go:82` stages/checks the single source using the archived pin. `host/coordinator/effectful.go:336` only creates a packet for exactly one successful Human.Approve; `decisionpacket.go:189` selects the authorizing grant expiry and parses the pending request. Existing test fake Human.Approve runners remain test-only and cannot satisfy M1b's real-source criterion.

Rig load sampled immediately after focused measurements: `uptime` load averages 1.96, 2.12, 2.08. No stress experiment or persistent background process was started.

## Design versus plan and residual boundaries

No substantive M1b mismatch: design S5 prescribes exactly the same effect/descriptor costs, requester and plan shape. Sprint splits design M1 into landed evidence law M1a and approval source/descriptor/docs M1b, explicitly forbids sibling imports and requires drift coverage for repetition. Existing README still truthfully describes M1a only; updating it is the planned M1b task. Plan status “pending” for M1a is historical metadata; current charter and live source/tests establish landed dependency.

The future host roster authenticates references in M5; M1b must not promote caller kind strings to evidence authority. Evidence-only kind admission currently exists and is tested. Fresh-seed equality, seeded six-shot capture, still mode and attended release preparation are later milestones and remain unproved here. Historical broad-suite failures in the planner baseline are not reclassified; this focused review neither reran nor certified the full repository gate. Controller must require the unchanged mandatory gates before landing and independently verify this readiness conclusion.

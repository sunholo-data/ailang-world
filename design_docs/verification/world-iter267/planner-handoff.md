# Iteration 267 planner handoff — row 150 M1b

Verdict: bounded M1b implementation is ready; landing is not certified. Base `57941bf7410c7462d755cd3a94748eaf76d66288`. Existing r4 is the authority, through the iteration-266 narrow-refinement carve-out; r1–r3 remain BLOCKED. This handoff creates no new design, human decision, ledger resolution or release authorization. M1b remains pending until implementation, independent review and unchanged CI pass.

Authoritative sprint-planner skill, CLAUDE.md, coding standards, mission guardrails, DESIGN §§1/14, existing design/plan/JSON and designer-readiness.md reviewed. Designer statements were independently checked against actual source, design S5, sprint M1b and executed focused gates; future source/proof/test success is not inferred. M1a is landed in `28ff4db371056929b9ddfa0c1dd9e87f95c8cfe6` (commit message records PASS88, zero blocking); only its JSON completion metadata and plan status were reconciled. Other milestones are preserved.

## Exact bounded executor directive

25 minutes target, stop at 30 minutes with an honest incomplete report. Only these implementation files may change:

- `packages/world-demo/approve.ail`
- `packages/world-demo/transitions.json`
- `packages/world-demo/README.md`
- `host/coordinator/world_demo_test.go`

No source changes were made by this planner. Executor performs no git writes, mint, publish, release, harness edit or ledger change. No new kernel constructor or host product implementation. Preserve frozen sources and recursive tools/launchd; compare before/after SHA256 and newline counts using the existing plan C11 recipe, plus scoped read-only diff. Do not restore protected originals. Snapshot changed allowlisted files and evidence under `.snap/M1b/` without credentials.

Implement a standalone convention-v2 `demo-approve` module (std imports only; do not import sibling evidence.ail). Publisher stages one source. Load the language prompt and explore builtins before writing AILANG. The released binary is v0.41.0; its active teaching prompt reports v0.16.6 (loaded with `ailang prompt`, distinct from binary version). Contracts and independently pinned inline identities precede pure implementation. Pure planning helpers use `! {}`; main uses only the named IO row needed for line protocol. Human.Approve is returned as plan data, executed by the existing host boundary.

Test-first `TestWorldDemoApprovalExactOnePlan` must invoke the real pin on actual source and decode plan JSON before any downstream coordinator guard. Happy input names goal, two candidates and exactly four unique kind/ref roster entries, no caller grade: test-invocation, recorded-effect, authored-claim, proof-envelope. Output has exactly one effect: id e1, Human.Approve@worktree, cost1, requester demo; finish=false and result=null. Descriptor access cost0 and declared cost1 must agree. Unknown input/argument keys, wrong/unknown phase, malformed/duplicate/missing/extra roster entries or caller grade refuse with zero effects. Do not authenticate evidence here or derive PROVEN from kind strings; future M5 owns reference validation. Document the complete runnable payload and plan/finish conventions in README, including refusal behavior.

Require nonempty healthy check/test/ai-check reports, named `oneApprovalEffect` VERIFIED, independent happy/refusal inline identities, zero errors/counterexamples/skips/vacuous skips. Existing manifest parsing helpers in world_demo_test.go may be reused and extended; do not replace named identities with counts or derive expected identities from source. If approval repeats evidence-kind predicates, extend the executable legacy drift probe to compare both real sources as well as the frozen kernel law. No absent binary, SKIP or empty test selection is green.

Mutation M1b-two-effects: append a second valid effect e2, retain valid types/JSON, and require the approval test itself to fail at `len(plan.effects)==1`. Compile errors, panic, timeout, pin absence, proof rejection before the target assertion or downstream L1 refusal do not kill this mutation. Keep the source contract/inline checks separate enough that the mutation drill demonstrably reaches the direct plan assertion. Save/restore only the allowlisted source byte-for-byte; rerun the pristine named test green and record mutant/assertion/rc/restoration evidence.

## Exact commands

```sh
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export GOTOOLCHAIN=go1.26.6
run_ail() { python3 -c 'import subprocess,sys; sys.exit(subprocess.run(sys.argv[1:],timeout=40).returncode)' "$@"; }
run_ail "$AILANG_BIN" --version
run_ail "$AILANG_BIN" prompt --version-active
run_ail "$AILANG_BIN" prompt
```

Run each focused command below with its exact named top-level PASS, not exit0 alone:

```sh
go test -run '^TestWorldDemoApprovalExactOnePlan$' ./host/coordinator -count=1 -timeout 120s -v
go test -run '^TestProducerMintsNoPacketForMultiEffectPlan$' ./host/coordinator -count=1 -timeout 45s -v
(cd packages/world-demo && run_ail "$HOME/.pinned-ailang/ailang" check approve.ail)
(cd packages/world-demo && run_ail "$HOME/.pinned-ailang/ailang" test --format json approve.ail)
(cd packages/world-demo && run_ail "$HOME/.pinned-ailang/ailang" ai-check -timeout 5s approve.ail)
go test -run '^TestWorldDemoLegacyGradeDrift$' ./host/coordinator -count=1 -timeout 120s -v
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh
```

Repeat the first named approval command for the valid two-effect mutant, then pristine restoration. Run before/after C11 observations. `verify_ail` omits world-demo; its green result cannot replace the Go source manifest test. Controller retains unchanged mandatory verify_ail/verify_go and CI before landing. Do not broaden executor work to repair unrelated baseline failures.

## C1 consistency and reuse audit

Design S5 and M1/C1 prescribe the same effect id/type/scope/cost/requester and plan shape as sprint M1b and JSON. M1a supplies landed package law; M1b supplies the self-contained plan, descriptor/docs and direct assertion/mutation. README currently describes M1a only, as expected before implementation. Approve source/descriptor and the new named test are absent at baseline; none was run as an empty green selection. No substantive C1 mismatch found.

Registry reuse decision for M1b: `none`. Released-pin `ailang pkg search approval` and `ailang pkg search world-demo` each rc0, no packages found. New demo-specific standalone capsule must use existing host convention and producer, not invent a registry release. Existing sprint predates registry_reuse and has no field (skill explicitly says such sprints warn); scope authorization permits only M1a metadata repair, so this populated audit is recorded here without changing other sprint fields.

## Independently measured baseline and pending boundary

Planner command:

```sh
go test -run '^(TestWorldDemoEvidencePinnedManifest|TestWorldDemoLegacyGradeDrift|TestDecisionPacketCodec|TestProducerMintsNoPacketForMultiEffectPlan)$' ./host/coordinator -count=1 -timeout 120s -v
```

rc0, wall4.653s, package4.337s; all four named top-level tests PASS. Evidence gate reports three named proofs, 21 named inline tests and null controls refused; drift executes real legacy constructors including both TestReport outcomes, raw ProofReceipt=CLAIMED and AiReview impersonation refusal. This confirms the M1a dependency and existing downstream two-effect refusal, not M1b acceptance. Load sampled immediately after measurements: 4.93,3.00,2.42. No load/stress experiment started; concurrent controller broad run means timing is not an isolated performance result.

Controller `~/.ailang/state/world-iter267/base-ail.log` currently ends required16/named40/se-tools444 and world-package9/9 PASS. `base-go.log` is still running at handoff preparation and already reports TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter and TestCheckSourceUnderPinnedInterpreter FAIL: the observed pin output includes a Google Cloud OpenTelemetry deprecation line and truncated module-loading error, lacking required LDR001. This is a concrete baseline finding, not established cause or permission to adjust tests/pin/harness. Controller owns completed baseline verdict and resolution; handoff readiness does not authorize landing red.

Controller subsequently measured the exact two failing pristine tests in two environment arms; planner inspected `~/.ailang/state/world-iter267/baseline-telemetry-results.json`: inherited telemetry rc1/2.21s; only GOOGLE_CLOUD_PROJECT, OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_RESOURCE_ATTRIBUTES and OTLP_GOOGLE_CLOUD_PROJECT removed rc0/1.93s. Controller reports both named tests PASS in the latter arm. This supports tracing environment contamination; it does not certify the broad gate. Clean-environment unchanged verify_go is running under600s bound, final verdict UNMEASURED here. Executor local child test/gate processes must remove ONLY those four tracing keys (use subprocess env or `env -u` prefixes), retain all other required pins and settings; messaging commands retain their canonical GCP environment. No source/test/gate changes or waiver. Protected scoped read-only diff rc0 after planner edits.

Changed planner artifacts: this handoff; `design_docs/planned/w-marketing-capture-sprint-plan.md` (M1a status only); `design_docs/planned/sprint_w-marketing-capture.json` (M1a status and completion evidence only). Token usage is not exposed by this role transport.

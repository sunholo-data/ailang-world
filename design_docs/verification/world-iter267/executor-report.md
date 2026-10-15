# Iteration 267 executor — M1b

IMPLEMENTATION_COMPLETE: true (M1b implementation only; independent review, candidate full Go gate and CI remain controller-owned). Base `57941bf7410c7462d755cd3a94748eaf76d66288`; no git writes, release, publication, session mint, ledger or milestone-state edits. Declared `codex:gpt-6.1-sol` model pin preserved; operator-authorized native transport override FLAGGED. Token usage unavailable.

Implemented only `packages/world-demo/approve.ail`, `packages/world-demo/transitions.json`, `packages/world-demo/README.md`, and `host/coordinator/world_demo_test.go`. The real standalone source returns one e1 Human.Approve@worktree cost1, requester demo, finish=false/result=null. Thirty-eight real-input refusals require zero effects. Descriptor access cost0/declared cost1 and the independently specified input schema are checked. Evidence-kind strings admit roster names only; no grade, reference authentication or PROVEN resolution is manufactured.

Boundary clarification: `host/broker/approve.go` approvalInputWire and HumanHandler.Execute strictly accept the requester payload; therefore the effect carries only `{requester:"demo"}`, while goal/candidates/roster remain recorded invocation input. Controller independently confirmed this existing-host adaptation. No production host change.

Fluency/checklist: released binary v0.41.0/24ee108; teaching prompt v0.16.6 loaded before AILANG authoring; builtins inspected. Scalar/string contracts precede implementation. Independently pinned `oneApprovalEffect` and `rosterKindAllowed` both VERIFIED; 12 inline identities and generated properties healthy, zero errors/counterexamples/skips/vacuous skips. JSON dispatch/rendering has independent happy/refusal `main` byte identities plus real-input cases. Main is pure `! {}` because convention-v2 returns a string; the runner renders that return, and Human.Approve remains host-executed plan data. No IO operation occurs in the module. Only std imports; executable drift compares both real package sources plus frozen kernel mappings. Existing nonempty report parsers and their null controls remain exercised by the evidence manifest gate.

Baseline isolation: unchanged controller baseline terminal rc0 was inspected before any worktree edit. Test-first exact approval command then failed rc1 on absent approve.ail (`test-first.log`); this is development red, not a claimed mutation kill. All scoped child processes removed ONLY GOOGLE_CLOUD_PROJECT, OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_RESOURCE_ATTRIBUTES, OTLP_GOOGLE_CLOUD_PROJECT, retaining other settings and pins. No env dump or secrets captured. No sandbox denial occurred.

Charter supersedes stale plan snapshot recipe: all snapshots/evidence are external at `~/.ailang/state/world-iter267/snap/M1b/`, per world-mission.md's executor snapshot guardrail. No .snap directory was created in the worktree.

## Final commands

Exact argv, cwd, rc, seconds, load and absolute log paths: external `final-commands.json`. Every Go command used count1, anchored exact top-level selection, timeout120s and -v; each named top-level PASS was checked. AILANG calls had40s outer bounds (ai-check solver5s); verify_ail outer300s. The README payload shell was extracted and executed verbatim, then its actual plan decoded and checked.

| Command / log stem | rc | seconds |
|---|---:|---:|
| `binary-version` | 0 | 0.017 |
| `prompt-version` | 0 | 0.015 |
| `TestWorldDemoApprovalExactOnePlan` | 0 | 6.297 |
| `TestProducerMintsNoPacketForMultiEffectPlan` | 0 | 0.54 |
| `TestWorldDemoLegacyGradeDrift` | 0 | 0.826 |
| `TestWorldDemoEvidencePinnedManifest` | 0 | 0.853 |
| `approve-check` | 0 | 0.028 |
| `approve-tests` | 0 | 0.454 |
| `approve-ai` | 0 | 0.04 |
| `README-quickstart` | 0 | 0.085 |
| `verify-ail` | 0 | 35.596 |
| `protected-diff` | 0 | 0.006 |

Logs use the table's stem plus `.log` under the external snapshot directory. `verify_ail` omits world-demo and is reported only as the unchanged repository AIL regression gate; the real-source approval/evidence Go manifests discharge this package's proof/tests. No executor full verify_go result is claimed. Controller's pristine baseline rc0 is not candidate certification.

## Mutation enumeration and measured red sets

Each of 26 mutations changed SHA256, typechecked approve.ail and built `go build ./host/coordinator` rc0 BEFORE reading its test verdict. Every mutant command rc1 reached its specified assertion; no build/proof failure before the target, panic, timeout, missing pin or empty selection counted. Each source/descriptor was restored from byte copies, SHA256 equality asserted, and the original complete named gate rerun green. `mutation-ledger.json` records exact edits, commands, rc/times/logs, mutant/restored hashes; `full-mutation-redsets.json` records full approval-test red sets; `support-mutations.json` records supporting-hunk drills.

M1b-two-effects appends a type-valid e2 Human.Approve. Exact approval command fails immediately at **`len(plan.effects)==1: got 2`**, before manifest proofs/inline tests and before downstream L1; pristine green. This is an M1b source-plan mutation, not the inherited producer refusal. The unchanged producer test is regression evidence only.

Every terminal validation reason below was neutered by replacing its sole `Some(reason)` with `None`. Decode and final-refusal rows change their refusal to approvalPlan. All displayed red sets are MEASURED with the complete approval command; parent top-level failure is omitted below. They are some-killer claims across this acceptance test, not claims about unrelated package tests.

| Refusal branch | Actual failing approval subtests |
|---|---|
| `refusal-roster-count` | missing-roster-entry |
| `refusal-roster-entry-schema` | duplicate-entry-key, caller-grade, roster-nonobject |
| `refusal-roster-kind` | missing-kind, unknown-kind, kind-type, extra-roster-entry |
| `refusal-roster-ref` | missing-ref, empty-ref, ref-type |
| `refusal-duplicate-kind` | duplicate-kind |
| `refusal-duplicate-ref` | duplicate-ref |
| `refusal-missing-roster` | missing-roster |
| `refusal-roster-array` | roster-type |
| `refusal-missing-candidates` | missing-candidates |
| `refusal-candidates-array` | candidates-type |
| `refusal-candidate-a-string` | candidate-a-type |
| `refusal-candidate-b-string` | candidate-type |
| `refusal-distinct-nonempty-candidates` | duplicate-candidates, empty-candidate |
| `refusal-candidate-count` | one-candidate, three-candidates |
| `refusal-args-schema` | nonobject-args, unknown-args, duplicate-arg-key |
| `refusal-goal-string` | missing-goal, empty-goal, goal-type |
| `refusal-input-schema` | nonobject, unknown-input, duplicate-input-key |
| `refusal-plan-phase-required` | unknown-phase, finish-phase, missing-phase, wrong-phase-type, pinned-manifest |
| `refusal-missing-args` | missing-args |
| `json-decode-refusal` | malformed-json |
| `terminal-refusal` | all 37 decoded-input refusal cases plus pinned-manifest (malformed-json stays green) |

Supporting mutations: kind predicate body→true kills only TestWorldDemoLegacyGradeDrift in that named command, at `roster kind AiReview`; scalar admitted result1→2 kills the direct approval count assertion (`got 0`); descriptor declared cost1→2 kills `approval descriptor`; inputSchema.additionalProperties=false→true kills `approval descriptor input schema drift`. All four typecheck/build0, mutant1, exact restore/pristine0. Their measured selected top-level red sets are in support-mutations.json; no whole-suite sole-killer claim.

## Diff-hunk audit / milestone ownership

| Shipped hunk | Load-bearing evidence |
|---|---|
| New approval scalar/string laws | Named proof/inline manifest; scalar-decision and kind-law compiling mutants; real-source drift |
| Strict key/schema and JSON/list validators | 38 runtime cases; all 19 terminal reasons + decode/final refusal neutered; full red sets above include unknown/duplicate keys and types |
| Approval/refusal JSON construction and pure main entry | Direct real-source plan envelope/payload/count assertions; required two-effect and terminal-refusal drills; independent main byte identities |
| New descriptor / input schema | Direct independent descriptor and schema assertions; compiling cost/schema mutants; absent descriptor cannot pass the new gate |
| Drift test's staging/import/output additions | Actual approval source executed alongside evidence and frozen kernel; kind-law mutant reaches drift comparison |
| New acceptance rig, manifests and refusal table | Test-first missing-source red; compiling M1b production mutations trigger this M1b gate, with healthy real-pin check/test/proof reports |
| README introduction and runnable approval section | Complete verbatim payload executed rc0 and actual JSON asserted; prose-only metadata has no separate mutation killer |

No unrelated supporting production hunk shipped and no measured behavioral survivor is concealed. Reverting M1b's new source makes its new acceptance gate fail as measured test-first; inherited M1a checks and downstream producer safety are not credited as M1b implementation. Reviewer should judge against the M1b diff, not the whole sprint.

C11: 12 protected path/hash/newline triples compare BYTE-IDENTICAL before/after; scoped `git diff --exit-code` rc0. Evidence: protected-before.json, protected-after.json, protected-diff.log. Observation only; no protected originals seeded, mutated or restored. Exact four final source copies and SHA256/bytes are in source-manifest.json; this report is also copied to the external snapshot.

Residuals: M1b does not authenticate roster refs (later M5), publish/mint or seed a world, validate release screenshots, run Windows, or certify candidate full Go/CI/independent review. Those are outside this authorized milestone. No load-bearing focused check remains untested. Timing/load readings are recorded, not presented as isolated performance claims.

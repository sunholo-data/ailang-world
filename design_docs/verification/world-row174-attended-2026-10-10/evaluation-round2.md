# Row 174 M1 + row 166 prerequisite — round 2 independent evaluation

**FAIL — 68/100.** A pristine independent invocation fails the new fixture's finite-budget guard; the root normal Go run independently fails the same guard. No failed gate is waived. Descriptor M1 semantics pass review; fixture prerequisite AC4 and M1 full-gate acceptance remain unmet. This report evaluates these bounded scopes, not the complete two-milestone sprint.

Product commit: `38c61e479c84e8b1a57c39f6766f4086ea0122c9`; baseline: `5ca111547fdd80480a8d42578fb59942ee8a0b63`. Isolated reviewer checkout: `/private/tmp/ailang-world-row174-evaluator-20261010`. Requested evaluator lane: `gpt-6.1-sol`, native fresh separate context, **FLAG same-model-fresh-context**, under D-WORLD-48 amended October 2. Cross-vendor independence is not claimed; actual runtime model attestation is unavailable. No pi transport or acknowledgement is fabricated. No external messages, PR, merge, publish, mint, shared-harness modification or mission-state write occurred.

## Blocking finding: calibrated fixture cannot fit observed valid startup

Independent pristine focused test, using isolated writable GOCACHE, performed all three real no-sleep controls successfully. Measured control durations were 2.353960541s, 4.666375ms, 4.4415ms. Its arithmetic produced **C=2.353960541s, B=25.54s, H=79s, R=27.636286083s**, grace=1s. It failed `finite fixture guard: B+H+grace must fit remaining caller deadline` before the exercised deadline arm. This is a discovered test failure, not a permission exception, compilation error, empty-output null result, skip, or mutation. The mechanism of the first slow startup is not established. Later successful runs do not erase it.

The controller's still-running normal `verify_go.sh` at this exact product head independently reports the same guard red: C=756.807791ms, B=9.569s, H=31s, R≈29.234s, all three controls PASS. The normal plain test leg is therefore red and race will not be reached. Final process return code and complete output remain pending at report time; **no full normal Go PASS is credited**. The known failed plain assertion already establishes the hard-fail condition. Completed evidence will be appended by the controller without changing this verdict.

Action: revise the fixture allowance through the separate approved prerequisite design, with measured finite-budget reasoning, retain all exact output/readiness/real-kill/reap assertions, rerun mutation and full normal gates on a new immutable head, then obtain a fresh review. Do not suppress the guard, retry/select a favorable run, increase production or caller bounds, or substitute old green logs. The proposed C+2-grace revision is not evaluated by this report.

## Independent verification and mutation

Evidence is banked in `round2-independent/` beside this report, with SHA-256 manifest.

- Initial default-cache run was blocked by sandbox cache access; its setup result is not a product verdict. An isolated writable-cache run then reached the test and produced the genuine pristine guard RED above.
- Separate approved standard-cache control PASS: C165.527375ms, B3.656s, H13s; actual timeout, readiness=true at real group kill, Outstanding=0 baseline=0, exact both streams. It is characterization evidence, not a waiver.
- Independent deadline-only stdout-loss mutant changes only the runHeadTail deadline return branch to replace the stdout sink. Three controls PASS; C93.587708ms/B2.936s/H11s; real kill readiness=true; exact timeout reached; reservation reaped to baseline. Only then stdout retention assertion fails with empty text/0 bytes/empty SHA while stderr remains exact. Mutant rc=1. Production file restored byte-exact, fresh pristine control PASS, zero SKIP. An initial mutation-construction assertion prevented editing; that associated successful invocation is explicitly banked as standard-pristine, not misrepresented as mutation evidence.
- Independent descriptor/source/live tests PASS, zero SKIP, under released interpreter v0.41.0 and tool v0.52.1. Live real HTTP exact-eight tools/list =10,368 SSE bytes. An earlier sandboxed live test could not bind localhost and is permission-limited, not descriptor evidence.
- Independently compare complete parsed manifest entries against `git show baseline`: only eight top-level description strings differ; workspace-exec is wholly equal. Immutable fixture bytes equal baseline exactly. All 28 tracked `.ail` files are byte-identical to baseline. Thus existing pinned AILANG verification may be carried forward for those unchanged sources; no new source proof is invented.

Root/executor banks are corroborating evidence, distinguished from own executions: focused plain count10 10PASS/0SKIP, race count5 5PASS/0SKIP/0race; eight fixture mutants intended RED with exact restore GREEN. Discard/either-printf/readiness mutants terminate in startup controls and are honestly early-control evidence. The coherent two-line WithCancel/caller-deadline-minus-grace mutant naturally exits, real Wait reaps, no kill, then timeout assertion RED. Separate stdout/stderr deadline sink resets reach three green controls, exact timeout and ready real kill before affected-stream RED. The independent stdout reset reproduces that load-bearing exercised assertion. These successes cannot waive the pristine/full normal guard failure.

Fourteen descriptor mutants remain banked with nonzero rc and intended identities: deletion/duplication/tenth entry census; effects/additionalProperties/output-world/transitionFnFile/unknown-field/workspace full-entry equivalence; old prose byte cap; source bytes; stale live projection; interpreter/tool verification bypasses. No production transition/projection/binary-verification changes occurred after round1, and current independent descriptor/source/live controls recheck this head. Existing zero-skip real-capsule/grant/budget/replay/security scoped bank has 29 top-level/52 subtests; it is prior unchanged-surface corroboration, not a new independent run. The prior banked AILANG gate has 16 required identities, 40 named tests, 444 se-tools tests, 9/9 package steps.

## Source and fixture identities

- Baseline manifest fixture SHA-256: `240840ca9a0f8f4c60ac159dd78fce26a6c681200e8a10c64fe0e08c0a228868`.
- Current fixture test SHA-256: `8d41e1de3034309c64a178556c4e3ea96818f0cab03d8e04c8ed65c5a518ae73`.
- Unchanged production handlers.go SHA-256 after own restoration: `e4659a6d23b74a7334b255cbb7ad05ea465994f7fe77e6339c1fdc1ce279096c`.
- Candidate manifest SHA-256: `d171b69548ba99bade5d74b34e9b7e0ce676829fc62a7f54587ce054bd138c97`.

## Full current diff inventory and design fidelity

Baseline→product head: 20 tracked files, 2210 insertions/18 deletions. Production change is exactly eight descriptions in `packages/se-tools/transitions.json`. Verification additions are coordinator descriptor test (141 lines), immutable 671-line baseline fixture, daemon live descriptor test (71 lines), and existing broker capture test (137-line diff). The latter is the only source-code delta since the first reviewed descriptor candidate; production Go is unchanged.

Design/plan additions: w-codex-world-overhead and companion sprint plan; separately approved w-head-tail-timeout-fixture and companion plan. Evidence/support additions: row174 CONTINUE.md, README.md, attended-mint.sh, attended-publish-candidate.sh, daemon-probe.log, daemon_probe_test.go.txt, final-profile.json, profile_final.py, runtime-before.json, verify-ail-wrong-env.log, verify-ail.log. Attended helpers were reviewed as prepared operator instructions, not executed fences. All product comparisons exclude the review report and new evidence artifacts; no other tracked product file changed in this checkout.

Canonical descriptor bytes 13,947→10,301; description UTF-8 bytes 7,430→3,794 meet ≤10,500/≤3,800. Wire and canonical bytes remain separate. Eight unique floor IDs and complete nine-entry census are pinned. Full schemas/property prose/access/effects/source references/unknown fields equal parent; workspace-exec complete equality is enforced. The 6-byte description-budget margin is nonblocking maintenance debt.

Semantic review retains IO/FS defaults, Env/Net exclusion, separate operator/grant authority, loopback/redirect restrictions, stdin/argv/args_json, 8-second timeout, denied non-execution, policy-effect admission and unrequested FS budget0, decision.missing_from_policy action, exit-zero inner-refusal warning, costs/provenance, exact-once edit, structured diagnostics, builtin degraded inventory/invalid-or-truncated refusal, CLI allowlist/write refusal, confined paths and unknown-key refusal. Existing QUICKSTART supplies operator setup and payload usage; no new syntax or example module is needed.

Fixture review: captures attach before real Start; WithTimeout precedes construction; real process-group kill is delegated; readiness follows both printf; observations are mutex protected; test is not parallel; reservation cleanup precedes global-seam restoration. Overflow checks and rounded B/H arithmetic are finite and reject invalid values. This is honest failure, but the arbitrary tenfold expansion prevents completion at observed startup timings. Missing restoration on incomplete reap remains a loud failed process-local test rather than silent continued shared-seam mutation; no race observed in scoped bank. No kernel/package architectural growth or §14 boundary change. Read CLAUDE, binding standards, charter profile/bar/D48, DESIGN §1/14 and relevant normal-role guidance. No compilation regression conditional is triggered. Performance verification belongs to M2 and is not credited.

## Round1 findings disposition and rubric

Round1 FAIL68: old partial-output fixture red in full plain/race; unchanged parent13PASS/7FAIL. Those logs remain banked, not waived. New observed-state fixture addresses the old false assumption but introduces an independently reproduced finite-budget guard failure. Metadata findings remain resolved. Completion artifacts stay pending until accepted M1; M2 deliberately pending is not dishonestly marked complete.

| Category | Score | Reason |
|---|---:|---|
| Tests | 0/20 | HARD FAIL: pristine guard RED; exact normal plain gate same RED |
| Lint | 10/10 | Added tests gofmt clean; no production Go change |
| Acceptance | 25/30 | Five M1 criteria met; sixth normal-gate/review criterion unmet; prerequisite AC4 unmet |
| Code quality | 13/15 | Narrow coherent fixture; calibration practical failure and incomplete completion artifacts |
| Documentation | 10/15 | Design/plan and usage adequate; no changelog; truthful pending statuses |
| Design fidelity | 10/10 | Metadata-only package repair plus separately designed fixture prerequisite |
| Total | **68/100** | Hard fail independently controls rejection |

M2 smoke has **not run**. No latency efficacy, row174 closure, clause4 satisfaction, complete-sprint PASS, fresh FINAL authority or benchmark/model/N/threshold change is granted. All design-quorum triggers were false and no quorum is claimed. Human publish/mint fences and fixed smoke models/arms/tasks remain prerequisites for M2.

EVALUATION_RESULT: fail
EVALUATION_SCORE: 68/100
EVALUATION_ROUND: 2
EVALUATION_REPORT_PATH: design_docs/verification/world-row174-attended-2026-10-10/evaluation-round2.md
FEEDBACK_SUMMARY: Descriptor M1 semantics pass; pristine calibrated timeout fixture and exact normal Go plain leg fail finite-budget guard, blocking prerequisite AC4/M1 gate. Revise finite calibration through approved design, preserve failures, rerun exact full gates and fresh review. M2 unrun.

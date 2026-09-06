# W-DE-FORK-CI-OWNERSHIP — Verify World After Driver Centralization

**Status**: Needs Human Review; two bounded unattended quorum rounds blocked before planning.
**Date**: 2026-09-06
**Target**: World dev, iteration 162
**Priority**: P0 — current dev CI regression
**Estimated**: 0.25 day, including controls and independent evaluation
**Dependencies**: Attended DE-FORK decision embodied by `e92594c23d76f51a195f36f31448ecedb5eddcfd`
**Planner-Lane**: codex-ok

## Problem Statement

The attended DE-FORK commit removed World's copied mission driver and routing test, but two World-owned verification entry points still invoke them. The macOS CI job fails extracting watchdog functions from the deleted driver; the Go verification script exits before Go build/test at the deleted routing test. This is a shared ownership-transition gap: deleting the implementation did not retire its consumers or update the CI job inventory assertion.

The repair must honor the new centralized owner. Recreating the deleted files or redirecting CI into a mutable sibling checkout would reinstate the dependency being removed. World should verify its own mission inputs and product gates, and explicitly decline to certify the installed shared runtime.

## Verification Log

All designer reads below used the isolated iteration worktree at `e92594c`; no production file was mutated for a probe. Controller-supplied observations are labelled separately.

| ID | Premise | Evidence and result |
|---|---|---|
| V1 | Centralized driver ownership is an attended decision | `git show --stat e92594c` names Mark's option (a), `MISSION_PROFILE=world`, `MISSION_WORKDIR` pointing to World, and removal of local driver/routing tests. `tools/launchd/README.md` directs runtime changes to the ailang registry. This establishes intended ownership; it is not a fresh live-runtime attestation. |
| V2 | Driver and routing test are absent at the base | `git ls-tree HEAD tools/launchd/mission-control.sh tools/launchd/test_mission_routing.sh` returns no entries; V1's deletion diff is the positive control. |
| V3 | Two active consumers still invoke deleted files | Read `.github/workflows/ci.yml:202-231` and `scripts/verify_go.sh:347-349`; workflow invokes stall/routing suites, script invokes routing and driver syntax checks. `rg` over `.github/workflows` and `scripts` confirms those active sites. |
| V4 | Stall failure reaches an absent driver, not a watchdog behavioral regression | Read `tools/launchd/test_mission_stall.sh:22-39`: `DRIVER` resolves beside the suite; awk extracts three functions and nonempty extraction is mandatory. Controller reproduced stall rc1 and routing rc127. |
| V5 | The job inventory test will reject removal without an explicit update | Read `host/verifygate/toolchain_pin_gate_test.go:115-245`: both positive controls and `wantJobs` include `launchd-drivers`; `wantGoPinnedJobs` contains only the other two jobs. Per-job counts and whole-file accounting are independent. |
| V6 | Default verification also assumes a local driver copy after the missing invocation | Read `scripts/verify_go.sh:351-384`: minimum-five tracked-path floor, working-tree drift check, and `check_driver_fleet` all precede Go work. Function at lines 125-259 compares local HEAD blobs to a sibling fleet HEAD; required pin-root path is a local-copy requirement. These facts do not identify the installed centralized driver. |
| V7 | A reusable World-local non-vacuous input validator exists | Read `scripts/mission_decisions.sh`: explicit `--file`, unreadable-file refusal, one-block cardinality, nonzero row floor, duplicate ID/status/empty-field validation. `/bin/bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md` returns rc0 and `decision ledger valid: 20 rows`. |
| V8 | Legacy fleet comparison has independent fixture coverage | Read `host/verifygate/driver_fleet_scope_gate_test.go`: synthetic World/fleet git repositories invoke `--driver-fleet-check`, assert three-match positive control and bounded residual disclosure. Keep this diagnostic and coverage intact. |
| V9 | Regression is observed in remote CI | Controller reports parent `b7e4a8e` green and tip run `34028333254` with ailang-code green, launchd and go-host red at V3/V4. Planner must retain exact run evidence; this row is controller-supplied, not designer-fetched. |
| V10 | A 10-second subprocess deadline is established test precedent | Controller measured, designer independently read `host/verifygate/toolchain_pin_gate_test.go:496` and `driver_fleet_scope_gate_test.go:221`: `context.WithTimeout(..., 10*time.Second)` with `exec.CommandContext`; timeout is a harness failure, not an ordinary expected refusal. |
| V11 | The explicit legacy diagnostic is pre-existing and red on the actual base | Controller ran `AILANG_FLEET_REPO=/Users/voightkampff/dev/sunholo-data/ailang /bin/bash scripts/verify_go.sh --driver-fleet-check`: rc1 at derive-planner-lane drift. V6 also establishes the invalidated local-copy premise. Charter queue row 76 explicitly calls the fleet-comparison arm frozen core. Retention is required by that authority constraint, not evidence that this dead diagnostic remains useful or currently passes. |
| V12 | Downstream references extend beyond the two active consumers | Controller search, independently repeated by designer with `rg -l 'verify_go\\.sh' host --glob '*test.go'`, finds eight files: `host/runbook/runbook_commands_test.go`, `host/runbook/runbook_stageb_test.go`, `host/transitionreg/transitionreg_test.go`, `host/broker/handlers_test.go`, `host/verifygate/toolchain_pin_gate_test.go`, `host/verifygate/driver_fleet_scope_gate_test.go`, `host/verifygate/ail_binary_gate_test.go`, `host/verifygate/evidence_manifest_gate_test.go`. Controller classified active default-flow consumers as CI plus toolchain-pin tests and runbook-stageb presence checking; evidence/AIl-binary/module-manifest checks use isolated modes. Preserve both modes and default-flow expectations. |
| V13 | Focused package requires the pinned binary even when the new mode does not | Controller's base focused-package test with AILANG_BIN unset returned rc1 at the intended anti-vacuity assertion. Required rerun is `AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang go test ./host/verifygate -count=1`; the unset result is not a product regression or a passing baseline. |

No AILANG syntax, compiler capability, diagnostic namespace, schema, or language semantics claims are made by this design.

## Goals and Scope

Restore both product verification jobs while retaining an executable, nonempty World mission-input check. Preserve independent Go/AIl gates and all unrelated integrity controls. Product gates must run with the central checkout absent, and their output must not claim shared runtime currency.

Frozen scope stays frozen: no edits, deletions, restoration, or copying under `tools/launchd/*`; no edits to symlinked skills, the V1 checkout, launchd state, central registry, or installed runtime. Cleaning up residual frozen artifacts is a fleet follow-up, not this repair.

The fleet-comparison arm inside the otherwise World-owned script is also frozen by charter row 76. Preserve its existing preamble/configuration, entire function, and `--driver-fleet-check` dispatcher byte-identical; do not add comments inside those regions. The attended DE-FORK authorizes removing its invocation from the default World product flow. It does not authorize rewriting or deleting the frozen diagnostic. Record this inherited dead/red interface as a limitation; remediation belongs to fleet. If independent review requires changing that frozen interface as a prerequisite, park the sprint for authority rather than inventing it.

## High-Impact Decisions

| Decision | Why it matters | Chosen by | Deadline | Change cost |
|---|---|---|---|---|
| Retire the local driver CI job | Ownership moved centrally under V1 | human, e92594c | resolved | low |
| Validate World's ledger with its existing validator | Keeps real mission-input evidence after removing copied routing tests | agent | design | low |
| Remove local-copy checks from default Go flow, preserve frozen dead diagnostic byte-identical | Default flow must verify World; charter row 76 withholds authority to repair or delete the old interface | agent, implementing attended ownership transition within frozen boundary | design | medium |

### Design Freeze

- [x] Central driver ownership is already attended and authoritative (V1).
- [x] No frozen path or shared runtime mutation is permitted.
- [x] New success wording certifies only World-owned inputs, never deployment or shared-driver freshness.

## Solution Design

1. Remove the obsolete `launchd-drivers` job. Keep both existing product jobs and their pins/steps unchanged. Update `TestGoToolchainPinsAgreeAndMatchJobList`'s explicit inventory and its stale positive control/comment, preserving the two independent sets and per-job attribution even though their members now happen to coincide.
2. Add `check_mission_config` in `scripts/verify_go.sh`, exposed as `--mission-config-check` before the AILANG binary preflight so hermetic tests can exercise it. The function runs the existing validator with the exact World charter path, checks shell syntax of the validator, propagates failures, and emits an explicit scope message: World decision-ledger validated; centralized runtime currency is not certified here. A missing validator or charter must fail, not skip. Reject extra arguments to the new mode.
3. Call that same function from normal verification in place of the old routing/syntax lines and local-copy drift block. Keep tracked-binary hygiene, compiler pin, canary, build, full tests, evidence manifest, and race behavior unchanged. The retired default block comprises tracked driver count, working-tree driver comparison, and default `check_driver_fleet` invocation: all three now measure a local copy that is no longer the runtime.
4. Preserve the entire frozen fleet-comparison region (existing preamble, `AILANG_FLEET_REPO`, `REQUIRED_FLEET_PATHS`, `check_driver_fleet`) and its `--driver-fleet-check` dispatcher byte-identical, as well as its fixture tests. Do not add comments or alter entry-point output there. It is an inherited dead/red local-copy diagnostic (V11), neither introduced nor certified by this sprint. Its synthetic fixture success only protects unchanged behavior; it does not claim a passing live diagnostic. Explain this boundary in this document and the replacement World-owned check's scope message, outside frozen regions.
5. Add World-owned regression tests that execute the real new mode against isolated copied-script fixtures and inspect active CI/default-script wiring. Fixtures include the real validator and deliberately valid/invalid ledger inputs; do not mock the validator into success. Wiring checks require known-positive World gate anchors before trusting absence results. Test default-flow integration with a valid pinned-binary shim and a controlled later stop, asserting the mission check executes before the stop; missing/invalid ledger must stop before Go. Shim results certify wiring only; real full gates certify product behavior. Every new subprocess helper, including fixture preparation commands and mutation runs, must use `context.WithTimeout(context.Background(), 10*time.Second)` and `exec.CommandContext`. Check the context after execution and fail explicitly on deadline expiry before examining an expected nonzero exit or diagnostic. Bound process/output cleanup too so inherited output pipes cannot keep the test hung after cancellation. Add a deliberately blocking fixture control proving timeout cannot count as the expected refusal.

### Files to Modify/Create

- `.github/workflows/ci.yml` (~30 lines removed) — retire obsolete driver job.
- `scripts/verify_go.sh` (~45 lines removed, ~30 added) — World mission-input check and retirement of default legacy invocation; frozen comparison and dispatcher bytes unchanged.
- `host/verifygate/toolchain_pin_gate_test.go` (~10 changed) — inventory update; preserve independent pin invariants.
- `host/verifygate/mission_config_gate_test.go` (~180 new) — real-validator fixture arms, default-flow wiring, no-local-driver consumer regression.
- `design_docs/planned/w-de-fork-ci-ownership.md` — execution evidence/status update.

## Conflict Surface

| Existing surface | Resolution and required preservation |
|---|---|
| CI exact job inventory and per-job Go pin attribution | Update inventory only; both product jobs retain exactly one of each pin kind. Reintroducing an unclassified job and misattributing a pin must still fail. |
| Existing `--evidence-manifest-check` and `--driver-fleet-check` modes | Preserve behavior and their tests. New mode must not fall through to compiler/Go work. |
| Default Go gate order | Replace obsolete driver block with real World ledger validation; preserve remaining gate sequence and failure propagation. |
| Historical driver drift policy and queue row 76 | The attended DE-FORK changes the default gate's subject. Frozen diagnostic/configuration/dispatcher remain byte-identical and pre-existing red; removal of the default invocation does not grant authority to repair that interface. |
| Eight downstream host-test reference files (V12) | Run the complete pinned verifygate package plus full product gate. Preserve isolated-mode consumers and runbook default-command/presence expectations; a grep inventory alone is not a behavioral pass. |
| Frozen leftover stall suite and helpers | Leave byte-identical and stop executing them from product CI. Ownership cleanup belongs to fleet. |
| Central configuration and installed launchd | Out of this repo's CI measurement boundary. Record limitation; do not substitute a local path guess or mutable fleet checkout for proof. |

## Acceptance Criteria and Mutation Plan

Planner must record baseline outcomes and exact diagnostics for each new behavioral acceptance, then executor/controller must record final outcomes. Fixture mutations belong in temporary fixture trees, never frozen working-tree paths.

The 10-second subprocess rule applies to every new helper-backed fixture, focused invocation and mutation arm in AC1-AC6, including positive controls: timeout is explicitly an instrument failure and never satisfies a refusal. The real full-package/full-product commands in AC7 use their own explicit larger top-level test/CI budgets because they intentionally exceed ten seconds; their individual new helper subprocesses still retain the ten-second deadline. Do not run the full product gate inside a ten-second fixture helper.

- [ ] AC1: New mission-config mode on real World input is rc0 and reports a nonzero valid row count plus explicit centralized-runtime non-certification. With AILANG_BIN unset and an absent fleet path it still succeeds, proving independence from those unrelated prerequisites.
- [ ] AC2: Fixture modes reject missing charter, missing validator, absent ledger block, zero-row ledger, duplicate ID, invalid status, and empty required field. Each returns nonzero with the validator/shell diagnostic appropriate to that input; a successful valid fixture is mandatory in the same harness. A deliberately blocking fixture is terminated at the ten-second context deadline and classified explicitly as timeout/instrument failure, never as one of these successful refusal arms.
- [ ] AC3: Normal default flow executes the mission check and proceeds to the next controlled stage when valid; invalid ledger aborts before Go. Record the reached-stage marker for the valid arm. Mutation M1 deletes only the default mission-check call: the integration assertion must turn red. Mutation M2 suppresses validator failure: an invalid-ledger arm must turn red. Mutation M3 bypasses the new function body: positive count/disclosure and bad-input arms must turn red.
- [ ] AC4: CI contains exactly the two retained product jobs and still invokes `./scripts/verify_go.sh` and `./scripts/verify_ail.sh`; no active CI command invokes the deleted local driver/routing/stall suites. Default Go flow does not invoke frozen local-driver paths or fleet comparison. M4 restores an obsolete invocation at each entry point independently; the corresponding wiring test must fail, naming the consumer. Do not scan only comments or accept an empty parsed command set.
- [ ] AC5: Go pin assertions remain non-vacuous. M5 moves one job's GOTOOLCHAIN declaration into the other while preserving whole-file count; per-job assertion fails. M6 adds an unclassified CI job; exact inventory fails.
- [ ] AC6: Existing explicit fleet diagnostic fixture suite stays green; default verification succeeds independently of an absent fleet checkout without setting CI to escape a guard. Compare extracted base/final frozen comparison/configuration and dispatcher bytes and record equal SHA256 digests. No tools/launchd or symlinked/shared file appears in the diff. The live explicit legacy diagnostic remains the pre-existing red limitation (V11), excluded from this sprint's success claim; it is not silently relabelled green by fixture results.
- [ ] AC7: All in-scope tests passing: complete `AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang go test ./host/verifygate -count=1`, real `AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang ./scripts/verify_go.sh`, real `./scripts/verify_ail.sh`, and both retained remote CI jobs. Planner records explicit finite top-level budgets for these complete gates (existing go-host CI ceiling is 25 minutes); deadline expiry is a failure. Any unrelated failure must be reported honestly, not treated as a pass or bypassed. The inherited frozen live diagnostic is the explicitly excluded result in AC6, not an unrelated failure discovered and waived after execution.
- [ ] AC8: Documentation updated with measured outcomes, diagnostic scope and the retirement of local-copy certification. Independently spawned evaluator approves before landing.

Non-vacuity question: deleting every old driver test alone can make CI green. AC1-AC3 require a real World-owned input validator to run and distinguish corrupt input; AC4 keeps the ownership regression visible; AC5 keeps existing job-pin semantics intact. These controls do not claim to test centralized routing, installed state, watchdog behavior, or fleet freshness.

## Examples

Before: `./scripts/verify_go.sh` reaches `tools/launchd/test_mission_routing.sh` and exits on a missing file before Go. After: it validates the World ledger, discloses its runtime boundary, then reaches the existing product checks. Focused mode: `/bin/bash scripts/verify_go.sh --mission-config-check` needs neither an AILANG binary nor the sibling central checkout.

## Deferred Decisions, Risks and Timeline

The implementer may choose test helper names and fixture layout, preserving behavioral criteria. One implementation milestone (~2h), controls/full gates (~2h), independent review and buffer (~2h). Main risks are accidentally weakening job-pin tests, retaining a default legacy call behind the newly repaired missing-file failure, and overstating what ledger validation proves; AC3-AC6 address those separately. Shared runtime doctor/pinning and frozen residue cleanup remain fleet work.

## Axiom Compliance

| Axiom | Score | Reason |
|---|---|---|
| A1 Determinism | +1 | Product CI no longer depends on mutable sibling driver copies. |
| A2 Replayability | 0 | Replay unchanged. |
| A3 Effect Legibility | 0 | Existing subprocess boundary remains explicit. |
| A4 Explicit Authority | +1 | Honors centralized owner and limits World assertions to its inputs. |
| A5 Bounded Verification | +1 | Focused real-validator mode with finite fixtures. |
| A6 Safe Concurrency | 0 | No runtime concurrency change. |
| A7 Machines First | +1 | Typed exits and explicit certification boundary. |
| A8 Minimal Syntax | 0 | No language syntax change. |
| A9 Cost Visibility | 0 | No billing or KPI semantics change. |
| A10 Composability | +1 | Reuses existing validator in both local and CI gate. |
| A11 Structured Failure | +1 | Invalid mission inputs remain failing gate results. |
| A12 System Boundary | +1 | Local input validation is distinguished from shared runtime attestation. |

Net +7; no hard violation of A1/A3/A4/A7.

## Related Documents and Quorum

`ailang docs search 'centralized driver CI verification' --neural --limit 5` actually reported SimHash results (100 docs), not neural scores; neural duplicate thresholds therefore cannot be claimed satisfied numerically. Read top contextual matches: `w-ddl-gate-teeth` concerns schema non-vacuity; `w-wiring-test-step-scoping-imprecise-under-key-reorder` concerns step parser precision; `w-setup-go-pin-unguarded` concerns compiler pin consistency. None repairs the newly introduced DE-FORK consumers. `w-driver-drift-gate-compares-the-copy-to-itself` is historical comparison design; this change preserves its explicit diagnostic and retires its invalidated default-runtime premise.

Scaffold fallback: the shared create script derives PROJECT_ROOT from its own symlinked skill location, which would write outside this task worktree. This document was created with apply_patch using the shared structure instead; no skill was copied or edited. The unattended trigger requires design quorum; controller owns that independent gate and records the exact review artifact/verdict before routing planning. This document is not its own passing verdict.

Round 1: **3/3 external reviewers rejected**, all present, cost **$0.0686**, artifact `.ailang/state/mission-quorum/w-de-fork-ci-ownership-2026-09-06T10-56-14Z.json`. Revision 1 addressed the controller-verified objections: enforce ten-second fixture subprocess deadlines and typed timeout failures (V10); narrow legacy retention to the byte-identical frozen, pre-existing dead/red interface (V11); enumerate downstream test consumers and require full pinned focused-package validation (V12-V13).

Round 2: **BLOCKED**, cost **$0.0213**, artifact `.ailang/state/mission-quorum/w-de-fork-ci-ownership-2026-09-06T11-00-58Z.json`. `gpt6-astra` was absent on budget; `gemini-3-1-pro` and `oc-glm-5-2` rejected. Gemini requires the full Go suite rather than only `host/verifygate`, a narrow mechanical correction. GLM disputes the larger authority split: the attended DE-FORK makes the fleet-copy diagnostic obsolete, while charter row 76 explicitly says World cannot modify the fleet-comparison arm.

The remaining blocking question is human-owned: does the attended DE-FORK supersede row 76 sufficiently to let World remove the obsolete explicit `--driver-fleet-check` mode and its fixture, or must World preserve that frozen-but-permanently-red historical diagnostic and limit this repair to removing it from default product verification? No planner or executor may run until that authority boundary is resolved and the revised design passes its normal gates.

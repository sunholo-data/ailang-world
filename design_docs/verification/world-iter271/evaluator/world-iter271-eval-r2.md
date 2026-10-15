# World iter-271 evaluator r2 — row 150 M2b (3789b7a..ea4b95d)

VERDICT: PASS, score 86/100. No blocking findings. (r1 was FAIL 66.)

## r1 blockers
- B1 (S2 read-bytes guard unproven): RESOLVED with pin. `demoRunStages` takes a task-ID suffix; the mismatch arm uses fresh IDs and asserts the error contains "S2 read bytes differ" and that S3 never dispatched (Test.Task empty). Mutant `if false` now fails at demo_test.go:744 "S2 read-bytes mismatch must be named: err=<nil>".
- B2 (guards without negative control): RESOLVED. New pin-free TestDemoSeedToolGuardControls (17 subtests) plus split, individually-worded refusal reasons. Every r1 survivor is now killed by the specific assertion it should be (table below).

## Non-blocking
N1. (item 2) The read-bytes mutant survives in ordinary CI (no tool pin): pin=0 run showed no FAIL, only SKIPs. Not blocking: the plan's M2b directive explicitly lets TestDemoSeedReadTestCheck's pin-requiring arm SKIP when the pin is unset, and acceptance requires a top-level PASS with the pin. Cheap improvement: a fake-/a2a/ unit control (already have fakeA2A) returning content "x" vs source "y" and asserting "S2 read bytes differ".
N2. (item 3) First-party v0.52.1: `ailang test --json world/contracts.ail` gives success=false, passed 7, failed 0, skipped 4, vacuous_skips 4 (controller is right). On the A2A path, the real nogen.ail (a copy of contracts.ail) and malformed.ail both report exit=1, validator "success=false is refused" (logged in the arm). Both agree success=false; the executor's "failing tests" wording is loose (failed_tests=0), but harmless: refusal is by success=false first. Consequence: the real arms do not reach the skipped/vacuous branches; those are covered only by the pin-free unit controls (killed). Acceptable.
N3. (item 4) The retained pin is asserted in the test from the daemon archive manifest (ReadManifest(sha256 of WORLD_TOOL_AILANG_BIN) contains v0.52.1), not carried in demoToolRun. Weaker than "retain actual pin" in the S3 record; the daemon is only configured with that tool binary and the version is also checked via requireBinaryRelease, so the v0.52.1-vs-v0.41.0 distinction is reasonably proven. Suggest adding a Pin field to the record in M3+. Non-blocking.
N4. Carried from r1: demo.go still prints S1-only success text; docs/QUICKSTART.md not modified (diff = demo_seed.go, demo_test.go); production demoSeedDomain refuses at S2 until descriptors are published/worktree exists (matches plan: publish is attended). Stage functions have no manifest writer.
No forbidden imports (archive added to the test file only, host/registry absent); no new production subprocess.

## Gates (AILANG_BIN=v0.41.0, WORLD_TOOL_AILANG_BIN=v0.52.1)
| command | rc | PASS |
| TestDemoSeedReadTestCheck | 0 | --- PASS (4.05s) |
| TestDemoSeedRejectsToolReportRefusals | 0 | --- PASS (+null control, real arms) |
| TestSeToolsE2EBinaryVersions ./host/daemon | 0 | --- PASS |
| TestSeToolsA2AEndToEnd ./host/daemon | 0 | --- PASS |
| TestCLIAgainstSeToolsDaemon | 0 | --- PASS |
| TestDemoSeedToolGuardControls (no pin) | 0 | --- PASS, 17 subtests |
| go test ./cmd/ailang-worldd WITH pin | 0 | ok 54.9s |
| go test ./cmd/ailang-worldd WITHOUT pin | 0 | ok 39.2s (pin arms SKIP, controls PASS) |
| TestQuickstart* cmd/daemon/broker | 0 | ok x3 |
| go vet ./cmd/ailang-worldd | 0 | ok |

## Mutations (each reverted)
| mutant | killed by |
| M2b-null (drop probe==nil) | RejectsToolReportRefusals/null-report-unit-control (:703 "null report is refused") |
| M2b-cli-op (S3 test->ai_check) | ReadTestCheck :655 "argv[op] is test" |
| read-bytes `if false` | pin=1: Refusals :744 "S2 read-bytes mismatch must be named"; pin=0: SURVIVES (N1, SKIP permitted) |
| drop skipped | GuardControls report-skipped-tests |
| drop vacuous | GuardControls report-vacuous-skips |
| drop zero-pass | GuardControls report-zero-passes |
| drop failed count | GuardControls report-failed-tests |
| drop success | GuardControls report-success-false |
| accept non-completed state | GuardControls a2a-non-completed |
| drop ai verified>0 | GuardControls ai-verified-zero |
| drop ai check.passed | GuardControls ai-check-failed |
| drop ai error_count | GuardControls ai-error-count |
| drop null ai_check guard | GuardControls ai-null |
| drop invocation/effect-ref guard | GuardControls a2a-missing-invocation + a2a-missing-effect-record |
| drop op allowlist | GuardControls op-allowlist-before-dispatch (dispatch count 1 != 0) |
| drop OK/exit check | Refusals :728 "non-test op must be refused" |
| extra read call | ReadTestCheck :689 Workspace.Read delta = 2 |

Worktree: `git status --short` empty, `git diff` 0 lines after all drills.

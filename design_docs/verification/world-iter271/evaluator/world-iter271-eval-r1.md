# World iter-271 evaluator r1 — row 150 M2b (d95cf50 vs 3789b7a)

VERDICT: FAIL (score 66/100; 70 needed; 2 blocking findings)

## Blocking
B1. cmd/ailang-worldd/demo_seed.go demoRunToolStages `if read.Content != string(source)` (S2 read-bytes guard) is UNPROVEN. Mutant `if false` survives all of TestDemoSeedReadTestCheck + TestDemoSeedRejectsToolReportRefusals. The only mismatch control (demo_test.go ~L728, "read bytes mismatch must refuse after S2") passes for the wrong reason: with the guard removed, the error is "demo-test: CLI refused or exited 1" (task id demo-test reused from the earlier arms), not the S2 mismatch. The positive test compares Read.Content to the same file, so it cannot detect a removed guard. Evidence: temp t.Logf printed `mismatch-arm err=demo-test: CLI refused or exited 1`, test PASS.
B2. demoValidateTestReport skipped/vacuous refusal, zero-pass refusal, and failed-count refusal have no negative control (plan C2: "named passes>0 and no failures/skips/vacuity"). Mutants `case false:` for each survive. Only null, malformed.ail and nogen.ail arms exist, and those are caught earlier by CLI exit code or the not-JSON branch. Same for demoValidateAICheck (verified==0, check.passed/error_count, null guard) and the effect-record/invocation-ref absence guard (survive), and the demoTestStage unlisted-op refusal (survives; the "non-test op must be refused" assertion is satisfied by malformed.ail failing ai_check, not by the op allowlist).

## Non-blocking
N1. (limit 1) Production demoSeedDomain now calls demoRunToolStages, but nothing in production publishes se-tools descriptors or creates the demo worktree; a real run refuses at S2 until the operator does so. Plan says test setup publishes via the production publisher core and publish/mint/capture are attended, so acceptable. The S2/S3 path itself is real /a2a/ code, reached by the test through demoRunToolStages (same function production calls). Refusal message ("S2/S3 refused; S2-S6 incomplete; no success manifest written") is truthful.
N2. (limit 2) "never publish a successful manifest" not asserted explicitly; stage functions have no manifest writer, and TestDemoSeedAuthenticatedCommits still expects "S2-S6 incomplete" (passes). Acceptable.
N3. (limit 3) cmd/ailang-worldd/demo.go still prints S1-only success text; outside allowlist; not touched. QUICKSTART.md in allowlist was NOT modified (diff = demo_seed.go, demo_test.go only), though the plan lists it; minor.
N4. (limit 4) v0.52.1 `ailang test --json` first-party: stdout is an object with failed_tests, passed_tests, module_path, properties[]; flags `--json`/`--format json` both exist. Validator field names `success`/`skipped_tests`/`vacuous_skips` are not verified present in the sample (only passed_tests/failed_tests seen); success/skipped absent would default zero => skipped guard and Success check semantics depend on fields not shown. Needs a control.
N5. Retained run lacks the actual pin path/version field (plan: "Retain actual pin"); v0.52.1-vs-v0.41.0 proven only by requireBinaryRelease on the configured binaries, not from the daemon's result.
Budget delta check is real: extra-call mutant killed (Workspace.Read delta 2 != 1). Imports in demo_seed.go: archive, authority, broker, daemon, hashref, store (pre-existing set; no host/registry). No new subprocess in production code; test uses exec only for `--version`.

## Gates (AILANG_BIN=v0.41.0, WORLD_TOOL_AILANG_BIN=v0.52.1)
| command | rc | PASS line |
| TestDemoSeedReadTestCheck | 0 | --- PASS (4.47s) |
| TestDemoSeedRejectsToolReportRefusals | 0 | --- PASS (+both subtests) |
| TestSeToolsE2EBinaryVersions ./host/daemon | 0 | --- PASS |
| TestSeToolsA2AEndToEnd ./host/daemon | 0 | --- PASS |
| TestCLIAgainstSeToolsDaemon | 0 | --- PASS |
| go test ./cmd/ailang-worldd WITH pin | 0 | ok 54s |
| same WITHOUT pin | 0 | ok 39s (ReadTestCheck + real-refusal-arms SKIP, null unit control PASS) |
| TestQuickstart* cmd/daemon/broker | 0 | ok x3 |
| go vet ./cmd/ailang-worldd | 0 | ok |

## Mutations (demo_seed.go; each reverted)
| mutant | killed by |
| M2b-null (drop probe==nil) | RejectsToolReportRefusals/null-report-unit-control |
| M2b-cli-op (S3 test->ai_check) | ReadTestCheck "argv[op] is test" |
| extra read call | ReadTestCheck Workspace.Read delta=2 |
| drop OK/exit check | RejectsToolReportRefusals "non-test op must be refused" |
| drop flags json | ReadTestCheck (not JSON) |
| drop read-bytes equality | SURVIVED (B1) |
| drop skipped/vacuous | SURVIVED (B2) |
| drop zero-pass | SURVIVED |
| drop failed/success | SURVIVED |
| accept non-completed A2A state | SURVIVED |
| drop ai_check verified>0 | SURVIVED |
| drop ai_check check.passed/error_count | SURVIVED |
| drop null ai_check guard | SURVIVED |
| drop invocation/effect-ref guard | SURVIVED |
| drop op allowlist | SURVIVED |

Worktree: `git status --short` empty, `git diff` empty after drills.

## Fix guidance
Add: a mismatched-source arm using fresh task ids (or assert error contains "S2 read bytes differ"); unit controls for demoValidateTestReport (skipped, vacuous, zero-pass, failed, success=false) and demoValidateAICheck (null, verified 0, passed false); a fake-/a2a/ (httptest) control for non-completed state and missing refs; op allowlist unit test.

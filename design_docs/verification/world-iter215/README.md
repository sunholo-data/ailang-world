# Iteration 215 — row 127 parked

Measured pristine base: 4e2600e601711f40dae42479167251d43a8a3bf6. Pinned compiler: v0.41.0, full commit 24ee1088776e21cd06a3781ed18e77f40be06db3. No product code changed.

Both quorum JSON artifacts are decision-bearing records, not passing gates. Author OpenAI was benched; GLM, Kimi and Gemini rejected both rounds; Sonnet absent (quota/unknown usage). R1 $0.2136992, R2 $0.1582822. R2 asks for source-verified card labels, corrected negative controls and a full-suite disposition. The GLM baseline-relative proposal and Kimi blocked-until-remediation-or-explicit-quorum rule require a judgment about acceptance; no controller-invented combined rule was used. D-WORLD-47 requests another design revision and fresh full quorum, default defer.

Commands and terminal results (all on unchanged base, outside role sandbox):

- `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/projection ./host/coordinator ./host/daemon -count=1 -timeout=300s`: rc0 (`base-go.log`).
- `AILANG_BIN=~/.pinned-ailang/ailang ./scripts/verify_ail.sh`: rc0, 16 identities/40 named tests (`base-ail.log`).
- `AILANG_BIN=~/.pinned-ailang/ailang ./scripts/verify_go.sh`: rc1 (`base-verifygo.log`), full suite red in TestF5WallClockTimeoutHasElapsedBound, TestRunContextCancel, TestInterfaceV2MovesWhenAnExportedADTGainsAConstructor, TestUnknownEntryInterpreterFailsAbsent, TestMissionConfigDefaultFlowIntegration/valid_ledger. Race leg did not run after failure.
- `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/capsule ./host/pkgproj ./host/replay ./host/verifygate -p 1 -count=1 -timeout=300s -run '^(TestF5WallClockTimeoutHasElapsedBound|TestRunContextCancel|TestInterfaceV2MovesWhenAnExportedADTGainsAConstructor|TestUnknownEntryInterpreterFailsAbsent|TestMissionConfigDefaultFlowIntegration)$'`: rc0 (`base-isolated.log`).
- `AILANG_BIN=~/.pinned-ailang/ailang go test ./... -p 1 -count=1 -timeout=300s`: rc0 (`base-full-serial.log`). This is a changed execution shape, not proof of a unique cause or green default verify_go.

No scope expansion to remediate baseline reds. Shared harness tickets are fleet-owned; no skill edit. Ignored .ailang/state/mission-quorum copies are local discovery aids; these tracked JSONs are the durable record.

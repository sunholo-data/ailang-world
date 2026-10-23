# Iteration 303 controller verification — darwin/arm64

Candidate `f615a9698370694664fbbc01431e4a2c549377cc`; pristine base `c73c79048daf9c3a19ea05c96e793436eb45e9c5`.
Pins: AILANG_BIN v0.41.0/24ee108; WORLD_TOOL_AILANG_BIN v0.54.15/0384450; Go1.26.6; Z3 4.16.0. Fleet CLI on PATH is separate.

- Focused three PNG tests and TestProductionGoSurface: rc0 outside sandbox.
- `bash scripts/verify_ail.sh`: rc0, 27.4s,16 required identities/40 named tests/496 se-tools tests.
- `bash scripts/verify_go.sh`: rc1,448.0s. Plain leg only; race leg NOT_REACHED. Full output retained in persistent evidence.
- Two unchanged-package failures: TestCheckSourceUnderPinnedInterpreter and TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter. Same exact targeted command failed at candidate and pristine base (2.4s/2.8s), packages byte-identical. Raw pinned check emitted458bytes with LDR001 at offset421, behind a deprecation warning at138; archive/outputHead caps at400. No code fix is included; queue202 owns it.
- Third full plain failure: carried TestMarketingSupervisorFakeChromeTimeoutKillsGroup (missing grandchild PID after8s). Isolated full marketing package rc0 (17.9s); exact inherited test at base rc0 (8.6s); full marketing -race rc0 (35.5s). Timing failure is not relabelled a new M7b defect or a green full suite. Existing M7a/M7c timing residue persists.
- Fresh controller mutation rerun outside sandbox: all8 named M7b edits BUILDS, intended assertion red, SHA-identical restore and originalgreen; receipt controller-mutations.json. Direct-helper reach for uniform/missing does not prove the later sibling end-to-end assertion executed under mutation.
- Protected12 paths matched HEAD and each other before/after stage execution. Main checkout status remains its original untracked .claude/worktrees/ only.
- Historical PR361 CI39880506556 on exact candidate is SUCCESS for all3 jobs. New record-extended PR head and merge require their own observed CI before final LANDED report.

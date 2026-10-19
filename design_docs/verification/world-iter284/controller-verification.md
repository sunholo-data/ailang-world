# Controller verification — iteration 284 (row 156 M3a + M3b)

Outside the codex sandbox, `WORLD_TOOL_AILANG_BIN` unset, `AILANG_BIN` = pinned v0.41.0, real srt 0.0.78 at `~/.ailang-world-exec/srt/node_modules` (darwin arm64).

## Baseline (pristine 1f36947)
- `go vet ./...` rc 0; focused broker tests rc 0; `go test -race -count=1 ./host/broker/...` rc 0 (242.7 s).

## After (executor tree, then rebuilt commits)
- `TestExecSrtProtectedWrites` real srt darwin: 18/18 PASS; parent-dir-rename logged written-and-recorded-as-detection-residual.
- MUT-NO-SRT-PROTECTED (append deleted) vs real srt: rc 1; `/modify` FAIL `handlers_exec_srt_test.go:629: protected operation must be refused: token "WORLD-PROBE-WROTE"` and `:632` host hash changed; 15 other prevention rows FAIL; parent-dir-rename and unprotected-control PASS. Restored → ok.
- Focused gates incl. `TestExecHandler.*` rc 0 (the executor’s in-sandbox red `TestExecHandlerTimeoutKillsTheGroup` passes outside the sandbox; sandbox denies `ps`).
- M3a boundary alone (M3b files at base): vet + focused rc 0. Commits rebuilt from `.snap/`, `shasum -c` OK on all 5 files.
- `go vet ./...` rc 0; broker `-race` rc 0; daemon `-race` rc 0 (313.2 s).

## Judge NB1 (r1) — controller fix
- Mutant `p.Protected = nil` after load (exec_profile.go:175): old `TestExecProtectedSettingsThroughHandler` ok (SURVIVED); after the literal-suffix fix: `/true` FAIL. Handler-drop mutant still FAIL. Clean green.
- First draft of the fix had the wrong literal order (canonical sort puts `*_test.go` before `tests/**`) and failed clean; corrected before commit.

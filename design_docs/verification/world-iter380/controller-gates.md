# Controller gates — iteration 380 (sprint head `1fe7c86d`)

Plan §6 run as one bash script (`gates.sh`, extracted verbatim from the plan) AFTER the last commit and BEFORE the judge: **rc 0**.

- `verify_ail.sh`: PASSED — 16 required identities, 40 named tests, 503 se-tools named tests.
- gofmt clean; `go build ./...`; `go vet ./...`; compile fence `go test -run '^$' ./...`: rc 0.
- Touched packages (`host/approval/...`, `host/store`, `host/authority`, `cmd/ailang-worldd`, `host/daemon`): 6 ok, 0 FAIL.
- `-race` on `host/approval/...`, `host/store`, `host/authority`: 0 DATA RACE, 0 FAIL.
- 22 named repo-wide inventories and reflective tests with `-v`: 22 PASS, 0 SKIP/FAIL.
- `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after commit: rc 0, 62 ok, 0 FAIL. Its two `WARNING: DATA RACE` lines are the known-positive control `design_docs/verification/w-race-gate-blindspot/racecontrol/main.go:16`.

Per-sub-milestone rebuild with `commit_ms.sh` (gofmt, build, vet, `go test ./host/approval/... ./host/store/ ./host/authority/`, plus `./cmd/ailang-worldd/` at M4a): every commit rc 0. Numstat (`git diff --numstat 3b77c294 HEAD -- cmd host`) after each commit: 345 / 672 / 890 / 1,027 / 1,300 / 1,538 / 1,777 / 1,843 = the plan's values. Final `cmd`/`host` tree `diff -r`-identical to the planner's mutant tree.

Plan round trip re-done by the controller: the 8 patches extracted from the plan file with §2's awk and applied to a fresh `git archive` of `3b77c294` reproduce the planner's `work/` tree (`diff -rq cmd host` empty).

Judge's three recompiled mutants reproduced by the controller (`go test -vet=off`): P1 (`v != Admitted && v == Admitted`) killed by `TestApproveByAssertionOrder`; C17 (`… && len(id) < 0`) killed by `TestPasskeyConsumeChecksTheRowInTheTransaction`; A6-count (`int64(count)*0 + 1`) killed by `TestMintForRequestDerivesPasskeyUse`. Judge worktree clean afterwards. Judge N2 (zero-valued `Passkey.Key`) read in code: `COSEKey.verify`'s default arm returns `ErrCOSEAlg`, so it fails closed.

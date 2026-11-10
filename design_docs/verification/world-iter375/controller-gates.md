# Iteration 375 — controller gates (row 201 A1b sprint 1)

Sprint head `048537c3` (branch `sprint/world-iter375-row201-a1b`, base `3a96cf75`). AILANG_BIN = `~/.pinned-ailang/ailang` (AILANG v0.41.0). Script: `gates.sh`.

| Gate | Result |
|---|---|
| `./scripts/verify_ail.sh` | rc 0 (81 s) |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` | rc 0 (439 s). 3 pinned fonts; 0 other binary blobs among 6,605 tracked files; toolchain floor ok; 64 evidence tests (`verify_go-048537c3-tail.txt`). Its race-detector known-positive control fired as designed. |
| `go build ./... && go vet ./...` | rc 0 |
| `TestProductionContextRoots` | rc 0 |
| `TestDaemonDependencyAllowlist` | rc 0 |
| `go test -race ./host/store/ ./host/approval/... ./host/authority/ -count=1` | rc 0 |
| per-milestone reconstruction | Each sub-milestone was restored from `snap/<M>/` and run through gofmt, build, vet and its package tests at the boundary (`commit_ms.sh`), then committed. Final tree sha256-identical to each executor run's output (3 runs, 0 mismatches). |

Per-milestone commits:
- M1a `bbfb6eba`
- M1b `eef98c1b`
- M1c `7b023888`
- M2a `cd91a67b`
- M2b `a28175ba`
- M2c `048537c3`

Design commits: `446e5d93`, `7b1b0366`, `5e2bf772`. Plan: `83bbc24b`.

## Judge findings reproduced by the controller

Each was run in the judge's worktree, and the file was restored afterwards.

- **PAT-UPPER2 survives.** Widening the pattern at `host/approval/approver.go:25` to accept uppercase after the first character leaves `go test ./host/approval/ -count=1` = `ok`. Non-blocking; logged per D-WORLD-109 and carried to sprint 2.
- **RV-DELETE is killed.** Replacing `RevokeApprover`'s guarded UPDATE with a DELETE makes `go test ./host/store/ -run '^TestRevokeApprover$'` FAIL.
- **CR-TOOMANY is killed.** This was iteration 374's survivor. The mutant `|| (false && errors.Is(err, store.ErrTooManyPendingRequests))` at `mint_request.go:142` makes CT1 FAIL.
  - A first attempt printed `ok`, but its sed had matched nothing. The empty-diff assertion caught this, and the run was repeated.

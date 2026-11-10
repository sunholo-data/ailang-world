# Iteration 374 — controller gates (row 201 sprint 2)

Sprint head `86aa2f9c` (branch `sprint/world-iter374-row201-s2`, base `ab4ae3aa`), AILANG_BIN = `~/.pinned-ailang/ailang` (AILANG v0.41.0).

| Gate | Result |
|---|---|
| `./scripts/verify_ail.sh` | rc 0 |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` | rc 0 — 3 pinned fonts, 0 other binary blobs among 6,587 tracked files, toolchain floor ok, 64 evidence tests (see `verify_go-86aa2f9c-tail.txt`) |
| `go test -race ./host/store/ ./host/approval/ ./host/authority/ ./cmd/ailang-worldd/ -count=1` | rc 0, 0 `DATA RACE` |
| per-milestone reconstruction | every milestone restored from the executor's `snap/<M>/`, its package tests re-run at the boundary, committed; final tree sha256-identical to each executor run's output (5 runs, 0 mismatches) |

Per-milestone commits: M3b `55a7cc25`, M4a `62148065`, M4b `3cae9428`, M4c `75f38336`, M5a `5e48bb28`, M5b `5e613e5d`, M6a `87e7540a`, M6b `4ad9248a`, M6c `1d2a3d17`, M6d `86aa2f9c` (plan `1edd723a`).

Instrument note: the first M4c boundary run of `host/daemon` went red on `TestDaemonStartReportsMissingTemplates` ("AILANG_BIN is unset: this test … never skips"). The controller had not exported the pin in that shell. With it exported: `ok` (267 s). Not a code defect.

## Judge finding reproduced by the controller (non-blocking, logged per D-WORLD-109, not filed)

`cmd/ailang-worldd/mint_request.go:142`: in `createMintRequest`, `ErrTooManyPendingRequests` maps to `exitUsage`, and no CLI test pins that. Mutant `|| (false && errors.Is(err, store.ErrTooManyPendingRequests))` → `go test ./cmd/ailang-worldd/ -run 'Request|Approve|EpisodeStart' -count=1` = `ok` (survives). Restored byte-identical. The store-level cap (ST7) is tested.

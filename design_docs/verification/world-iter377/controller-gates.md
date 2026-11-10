# Iteration 377 — controller gates at sprint head 9f014e0d

AILANG_BIN = ~/.pinned-ailang/ailang (v0.41.0). Run by `gates.sh`, sequenced BEFORE the judge (iteration-376 lesson: no concurrent sweeps).

| Gate | rc | Time |
|---|---|---|
| `scripts/verify_ail.sh` | 0 | 67 s |
| gofmt (cmd/ailang-worldd, host/approval, host/store, host/authority) | 0 | 0 s |
| `go build ./... && go vet ./...` | 0 | 4 s |
| `go test -run '^$' ./...` (_test compile fence) | 0 | 3 s |
| touched packages (`cmd/ailang-worldd`, `host/approval/...`, `host/store`, `host/authority`) | 0 | 164 s |
| `TestProductionContextRoots` (unedited) | 0 | 1 s |
| `TestDaemonDependencyAllowlist` | 0 | 2 s |
| `-race` touched host packages | 0 | 34 s |
| `-race` new CLI tests `-count=3` | 0 | 20 s |
| no-SKIP `-v` run of the approve/approver/session tests | 0 | 6 s |
| `VERIFY_GO_LEGS=plain scripts/verify_go.sh` (incl. full `go test ./... -count=1`) | 0 | 502 s |
| clean `git status` | 0 | 1 s |

The two `WARNING: DATA RACE` lines inside `verify_go.sh` are its race-detector KNOWN-POSITIVE control (`race-control-excerpt.txt`; "Found 2 data race(s)", exit 66 expected), not a finding.

Per-commit numstat (`git diff --numstat 545fdaae HEAD -- cmd host`), against the plan's §7 expected values, after each controller commit:
M4a 221 · M4b 323 · G1 511 · M4c1 602 · M4c2 815 · M4d1 993 · M4d2 1125 · M4e 1314 (size gate ≤ 1,404: PASS) · M5a 1577 · M5b 1752 · M5c 1910. All 11 matched exactly.

Controller reproduction of the judge's survivor SIG-NOW (`approve_sig.go:53` `Now: now` → `Now: row.CreatedAt`): `go test ./cmd/ailang-worldd/ -run '^(TestApprove)'` → ok (SURVIVES). Control (`ExpiresAt: row.ExpiresAt` → `row.CreatedAt`, same line) → FAIL (killed). Judge worktree restored clean.

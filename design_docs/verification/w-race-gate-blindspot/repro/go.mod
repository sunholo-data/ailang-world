// Deliberately a SEPARATE, NESTED module: `go build ./...` / `go test ./...` in the
// repository root module must NOT pick this reproducer up. It is a diagnostic
// artifact, not host code, and on an affected toolchain it is expected to misbehave.
module ailang-world/verification/go1_26-arraylit-miscompile

// The `go 1.22` line below is LOAD-BEARING: it must stay at or below the
// oldest toolchain in ../run.sh's KNOWN_BAD list, or every deny-listed probe
// SKIPs and the instrument is disarmed. Enforced by
// TestReproModuleFloorStaysBelowKnownBadToolchains (host/verifygate). Do not
// let `go mod tidy -go=…` or an IDE floor-bump touch it.
go 1.22

// Deliberately a SEPARATE, NESTED module: `go build ./...` / `go test ./...` in the
// repository root module must NOT pick this race-detector control up. It is a
// diagnostic artifact, not host code, and it deliberately contains a data race.
//
// The `go 1.22` line below is LOAD-BEARING: it must stay at or below the repository
// root module floor (`go 1.26.6`, the `go` line of the root go.mod), because
// scripts/verify_go.sh runs this control under GOTOOLCHAIN=$ACTIVE_GO derived from
// the root module's GOVERSION; a floor above that toolchain makes `go run -race .`
// refuse before it can fire (verify_go.sh then FATALs "the race detector is not
// armed"). Enforced by TestRaceControlFloorStaysBelowRootToolchain (host/verifygate).
// Do not let `go mod tidy -go=…` or an IDE floor-bump touch it.
module ailang-world/verification/race-detector-control

go 1.22

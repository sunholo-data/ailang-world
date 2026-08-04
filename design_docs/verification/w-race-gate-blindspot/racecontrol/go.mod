// Deliberately a SEPARATE, NESTED module: `go build ./...` / `go test ./...` in the
// repository root module must NOT pick this race-detector control up. It is a
// diagnostic artifact, not host code, and it deliberately contains a data race.
module ailang-world/verification/race-detector-control

go 1.22

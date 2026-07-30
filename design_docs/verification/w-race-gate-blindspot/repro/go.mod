// Deliberately a SEPARATE, NESTED module: `go build ./...` / `go test ./...` in the
// repository root module must NOT pick this reproducer up. It is a diagnostic
// artifact, not host code, and on an affected toolchain it is expected to misbehave.
module ailang-world/verification/go1_26-arraylit-miscompile

go 1.22

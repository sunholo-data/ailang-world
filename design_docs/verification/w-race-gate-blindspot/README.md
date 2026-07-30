# `w-race-gate-blindspot` — committed reproduction of a go1.26.x code-generation defect

**Status**: reproduction CONFIRMED first-party, iteration 40 (2026-07-30). Remediation is
**parked for human ratification** — see `design_docs/planned/w-race-gate-blindspot.md`.

This directory exists because **prose decays and a committed artifact does not**. Queue item 4e
carried, for three iterations, the hypothesis that a `-race`-only test failure in `host/store` was
`modernc.org/sqlite`'s `unsafe` usage interacting with the race detector's memory layout. That
hypothesis is **refuted**. The real cause is a Go compiler regression, and this is the smallest
program that shows it.

## What the defect is

On **darwin/arm64**, go **1.26.0 → 1.26.5** (1.26.5 being the latest stable release as of
2026-07-30) miscompile the following source shape:

> a **local** array literal declared inside a function, indexed by a `range` index variable, whose
> element is assigned to a struct **string** field that is **not at offset 0**, inside a composite
> literal that also contains an **interface method call** *after* the array access, appended to a
> slice, followed by `break`.

The resulting struct's string field is wrong. Depending on the struct's layout it is either **empty**
or carries a **corrupt string header** — a nil data pointer with a garbage length — in which case
merely printing it segfaults:

```
len(Table) = 6 len(Ref) = 1 len(Field) = 4334851712 len(Reason) = 10
panic: runtime error: invalid memory address or nil pointer dereference
```

Key properties, each measured:

| Property | Result |
|---|---|
| Requires `-race`? | **No.** A default `go build` reproduces. |
| `-gcflags=all=-N` (optimizations off) | **Correct** — so it is the optimizer |
| `-gcflags=all=-l` (inlining off) | **Still wrong** — so it is not inlining |
| go1.25.6, go1.24.9 | **Correct** — so it is a 1.26 regression |
| go1.26.0 / .3 / .4 / .5 | **Wrong** — including the newest stable, so there is nothing to upgrade *to* |
| `GOOS=linux GOARCH=amd64` | **UNDETERMINED.** It compiles, but this rig has no Rosetta so the binary cannot be executed here. Do not infer anything from that. |

A sub-agent's SSA-dump reading (**UNVERIFIED by the controller**, recorded as a lead, not a fact)
attributes the final wrong code to the `late_fuse` pass decomposing a 16-byte `string` store into
separate ptr/len stores at offsets shifted one slot too late, set up by `generic_cse` and `prove`;
disabling any one of those three individually also makes the program correct.

## How this repository was exposed

- `go.mod` declared `go 1.26.4`, and `.github/workflows/ci.yml` selects the toolchain with
  `go-version-file: go.mod`. **CI therefore compiled this repository with an affected toolchain**,
  as did every local `go build` / `go test` / `scripts/verify_go.sh` run on the rig.
- A census of the at-risk shape in production code found **exactly two** sites, both in
  `host/store/scan.go` (lines 74 and 112) — which are precisely the two functions whose tests
  exhibited the two symptoms recorded against item 4e.
- Under `-race`, those two sites produced a **wrong value** (`Field: ""` where `"prevEntryHash"` was
  required) and a **hang**, respectively. Under a default build they happen to produce correct
  values on the paths the tests cover. **That is luck, not safety** — the same compiler
  demonstrably emits a corrupt string header for the same shape in this reproducer's layout.

## Re-running it

```bash
./run.sh
```

`run.sh` is written as an **instrument with its own known-positive controls** and exits non-zero
rather than reporting a clean result when it cannot see what it is supposed to see: it requires that
at least one known-affected toolchain actually reported `BUG` *and* at least one known-good
toolchain actually reported `OK`. If Go fixes this upstream, the first control fails and the script
tells you to re-derive the pin decision — instead of silently printing "all clear". An empty reading
from any probe is treated as an instrument failure, never a pass.

Captured output from the confirming run is in [OUTPUTS.md](OUTPUTS.md).

## Files

| Path | What |
|---|---|
| `repro/main.go` | the 52-line self-checking reproducer |
| `repro/go.mod` | **deliberately a nested, separate module** so the root module's `go build ./...` and `go test ./...` never pick this up — verified: `go list ./...` still returns exactly 10 packages |
| `run.sh` | re-runs the reproduction across toolchains, with controls |
| `OUTPUTS.md` | captured evidence from the iteration-40 confirming run |

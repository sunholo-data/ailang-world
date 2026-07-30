# Captured evidence — iteration 40 (2026-07-30)

All runs first-party by the mission controller on the development rig
(`darwin/arm64`, default toolchain go1.26.4, `CGO_ENABLED=1`), against clean
`origin/dev` @ `8ed04c0543c1bdadedc694fdaad154dbf8dfb38b` unless stated otherwise.

## 1. The original symptom, reproduced in the repository

`TestScanUnreadableLogKeysetResumes`, `-run`-filtered, same command twice with the only
difference being `-race`:

```
=== A: WITHOUT -race (control) ===
rc=0
--- PASS: TestScanUnreadableLogKeysetResumes (0.00s)
ok  	github.com/sunholo-data/ailang-world/host/store	0.320s

=== B: WITH -race ===
rc=1
    scan_test.go:51: first page = {Rows:[{Table:log_entries Index:2 Ref: Field: Reason:hashref: empty hashref text}] Scanned:2 NextIndex:3 NextRef: Done:false}, err=<nil>
--- FAIL: TestScanUnreadableLogKeysetResumes (0.02s)
```

Note `Field:` is empty while `Reason:` and `Index:2` from the same loop iteration are correct.

> **Refutation of an inherited inference.** Iteration 38 recorded that `Rows[0].Ref` was "empty
> alongside `Field`, consistent with the `fields` array reading as all-zero rather than one element
> being wrong — which would make it a memory-corruption signature, not a logic bug." Reading
> `host/store/scan.go:76-79`, `ScanUnreadableLog` **never assigns `Ref` at all** on that path, so
> `Ref` is empty in every run, with or without `-race`. The strongest stated evidence for the
> "all-zero array" reading rested on a field that is empty by construction.

## 2. It is the optimizer, not inlining, and it is local to `host/store`

Same test, varying only compiler flags:

| Command | Result |
|---|---|
| `-race` | **FAIL** (`Field:` empty) |
| `-race -gcflags='all=-N -l'` | ok 1.616s |
| `-race -gcflags='github.com/sunholo-data/ailang-world/host/store=-N -l'` | ok 1.682s |
| `-race -gcflags='github.com/sunholo-data/ailang-world/host/store=-N'` | ok 1.393s |
| `-race -gcflags='all=-l'` | **FAIL** (`Field:` empty) |
| `-gcflags='all=-N -l'`, no `-race` | ok 0.317s |

Disabling optimizations for the single package `host/store` is sufficient; disabling inlining
everywhere is not. So: **the optimizer, in this package.**

## 3. One mechanism explains BOTH of item 4e's symptoms

Item 4e recorded two distinct symptoms in one file — a deterministic *failure* of
`TestScanUnreadableLogKeysetResumes` and a *hang* of `TestScanUnreadableWorldsFindsPoison` — and
required that any mechanism hypothesis explain both.

| `TestScanUnreadableWorldsFindsPoison` | Result |
|---|---|
| `-race` | **`panic: test timed out after 1m0s`**, `FAIL … 62.669s` |
| `-race -gcflags='…/host/store=-N -l'` | ok 1.538s |
| `-race -gcflags='…/host/store=-N'` | ok 1.513s |
| no `-race` | ok 0.184s |

The same one-package optimization switch clears the failure *and* the hang.

## 4. It is a go1.26 regression, and it does not need `-race`

Standalone 52-line program, no dependencies, `go vet` clean, build cache cleared. Per-field
lengths printed before the strings so a corrupt header is visible without dereferencing it:

```
=== A. plain go build + run (NO -race) ===
len(rows) = 1
len(Table) = 6 len(Ref) = 1 len(Field) = 4334851712 len(Reason) = 10
Table="worlds"
Ref="w"
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x2 addr=0x0 pc=0x1025e24b4]
rc=2

=== B. -gcflags=all=-N (opts off, no -race) ===
len(Table) = 6 len(Ref) = 1 len(Field) = 9 len(Reason) = 10
Field="stateRoot"
Reason="empty text"
rc=0

=== C. -race ===
len(Field) = 4312746944   -> same SIGSEGV
rc=2
```

Toolchain sweep on the final minimized form (this directory's `repro/`):

```
go1.24.9 -> OK
go1.25.6 -> OK
go1.26.0 -> BUG: Field="" want "stateRoot"
go1.26.3 -> BUG: Field="" want "stateRoot"
go1.26.4 -> BUG: Field="" want "stateRoot"
go1.26.5 -> BUG: Field="" want "stateRoot"
```

**go1.26.5 is the latest stable release**, so there is no fixed version to upgrade forward to.

## 5. `run.sh`, with both controls firing

```
== go1.26 local-array-literal miscompilation reproduction ==
host: darwin/arm64   default toolchain: go1.26.4

expected BUG (affected):
  go1.26.0   expect=BAD   got: BUG: Field="" want "stateRoot" (rc=0)
  go1.26.3   expect=BAD   got: BUG: Field="" want "stateRoot" (rc=0)
  go1.26.4   expect=BAD   got: BUG: Field="" want "stateRoot" (rc=0)
  go1.26.5   expect=BAD   got: BUG: Field="" want "stateRoot" (rc=0)
expected OK (unaffected):
  go1.25.6   expect=GOOD  got: OK (rc=0)
  go1.24.9   expect=GOOD  got: OK (rc=0)

-- optimization-level control on the default toolchain --
  -gcflags=all=-N        got: OK
  -gcflags=all=-l        got: BUG: Field="" want "stateRoot"

RESULT: reproduction confirmed, and both controls fired.
```

## 6. The `-race` gate item 4e asks about is viable, and cheap

Measured in a worktree off `8ed04c0` with `go.mod`'s `go` directive temporarily lowered to
`1.25.6` and `GOTOOLCHAIN=go1.25.6`, `AILANG_BIN=/tmp/ailang-v0300/ailang` (v0.30.0, `e37b370`):

```
go test ./... -count=1 -race        rc=0   elapsed=179s
ok  cmd/ailang-worldd    4.410s      ok  host/hashref    2.876s
ok  host/archive         8.358s      ok  host/registry   3.245s
ok  host/broker        176.615s      ok  host/replay    28.155s
ok  host/canon           2.134s      ok  host/store     17.069s
ok  host/capsule        16.366s
ok  host/daemon         14.669s
DATA RACE warnings: 0        (control: 10 `ok` lines, so the grep works)
```

Whole-package `host/store` under `-race`, the two toolchains, same worktree, minutes apart:

| Toolchain | Result |
|---|---|
| go1.25.6 | **rc=0, `ok … 11.710s`** |
| go1.26.4 | **rc=1, `--- FAIL: TestScanUnreadableLogKeysetResumes` + `panic: test timed out`, 124.699s** |

So a `-race` leg costs **~179 s wall**, dominated by `host/broker` at 176.6 s — which is
carry-forward **CF-MJC-1** (MJ.C's two tests exercising the `maxRecoveryPages` bound at its real
value of 2^20 pages). Making that bound injectable is the lever if 179 s is judged too slow.

> **Two false starts recorded rather than hidden.** (i) My first full-repo `-race` sweep reported
> two FAILs (`TestCLIRealSubprocessEpisode`, `TestEpisodeLiveReplayThreeArmsAndEvidence`) which I
> was about to attribute to `-race`. A no-`-race` control showed the broker one fails **either
> way**, with `AILANG_BIN must name the pinned released interpreter` — I had not exported
> `AILANG_BIN`. That is the M6/B1 anti-false-green guard working exactly as designed, and both
> FAILs disappeared once it was set. (ii) I ran `GOARCH=amd64 go build && ./amd64 | sed …`, got
> empty output with `rc=0`, and nearly recorded "amd64 is unaffected". The binary had never run —
> `bad CPU type in executable`, no Rosetta on this rig. Architecture scope is **UNDETERMINED**.

## 7. Census of the at-risk shape in production code

```
$ grep -rn ':= \[\.\.\.\]' --include='*.go' . | grep -v '_test.go'
./host/store/scan.go:74:  fields := [...]string{"entryHash", "transitionFn", "interpreter", "prevEntryHash", "transitionRef"}
./host/store/scan.go:112: fields := [...]string{"worldRef", "stateRoot", "logHead"}
```

Exactly two sites, both in the two affected functions. (Run with a known-positive control in the
same call: the control `grep -n 'fields := \[\.\.\.\]' host/store/scan.go` returned both lines,
proving the pattern was not silently glob-mangled by zsh.)

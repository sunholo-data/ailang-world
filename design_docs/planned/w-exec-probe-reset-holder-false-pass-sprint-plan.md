# Sprint plan — `w-exec-probe-reset-holder-false-pass` (iteration 318, row 232, clause 3)

**Design doc:** [`design_docs/planned/w-exec-probe-reset-holder-false-pass.md`](w-exec-probe-reset-holder-false-pass.md)
(r2; quorum PROCEED at N-1, `gpt6-1-sol` absent unreachable). **Parent:**
[`design_docs/implemented/w-exec-probe-reverse-leg-port-race.md`](../implemented/w-exec-probe-reverse-leg-port-race.md) (§4.2).
**Sprint base:** worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter318-row232`, branch
`sprint/world-iter318-row232`, HEAD `e1110207`. Its Go tree equals the design's measurement base
`df5a0671`: `git diff --stat df5a0671 HEAD` lists only `design_docs/` files (13 files, V-P1).
**Planner:** mission-control iteration 318 (planner role). **Executor:** claude-haiku, so every edit
is spelled out below. Copy code from the fenced blocks, never from table cells.
**Estimate:** about 0.25 day. **Two milestones**, each independently green. **Risk:** low. The
whole change and every mutant were prototyped by the planner in a scratch copy
(`/tmp/w318-planner/proto`, outside the repo); the results are in §9.
**Scope (the only files you may edit):**
- `host/broker/exec_probe.go`
- `host/broker/exec_probe_test.go`
- `.github/workflows/ci.yml` (line 247 only)
- `docs/QUICKSTART.md` (line 1153 only)
- `design_docs/planned/w-exec-probe-reset-holder-false-pass.md` (the four narrow amendments in M2 only)
- `design_docs/implemented/w-exec-probe-reverse-leg-port-race.md` (the two amendment notes in M2 only)

**No `.ail` file. No other file.**

---

## 0. The one-paragraph version

`pollLoopback` classifies a reverse-leg dial only when `connect` succeeds. If the peer accepts and
resets before Go reads `SO_ERROR`, `connect` itself returns ECONNRESET, and today that dial records
nothing. T7's `reset-holder` makes one dial, so the attempt looks clean and passes. That is the red
`go race gate` on `dev` (run 40833618366). The fix:
- a new pure `dialErrIsSilent(err) bool` that is true only for `errors.Is(err, syscall.ECONNRESET)`;
- a dial seam on `pollLoopback`, so a test can inject dial errors with no timing;
- an `else if dialErrIsSilent(err) { silent.Store(true) }` branch.

Two new tests:
- **T8** `TestExecProbeReverseDialErrIsSilent`: a table over the classifier.
- **T9** `TestExecProbeReversePollRecordsConnectReset`: the poll with a fake dialer.

Plus the CI list, one QUICKSTART clause and the reviewers' doc amendments.
- **M1:** the code, T8, T9 and the M22–M29 mutation drill.
- **M2:** CI, QUICKSTART, the doc amendments, AC-STABLE and the gates.

---

## 1. Environment (on EVERY command)

```sh
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter318-row232
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang     # reports "AILANG v0.41.0" (measured)
```

| Pin | Measured this session | Why it matters |
|---|---|---|
| `go` | `go version go1.26.6 darwin/arm64` in the worktree (`go.mod:3` `go 1.26.6`) | The CI pin. Outside the worktree the shell's go is go1.26.4, so run every command from the worktree |
| `AILANG_BIN` | `AILANG v0.41.0` | Without it the daemon e2e tests SKIP, and `verify_ail.sh` cannot run |
| `WORLD_EXEC_SRT_NODE_MODULES` | **unset** | The real-srt tests SKIP locally. Do NOT set it |

**Instrument rule.** Never write `cmd | tail; echo rc=$?`. Redirect to a file, put `echo "rc=$?"` on
the next line, then read the file. A `-run` that matches nothing exits 0, so read every verdict
from counted `--- PASS:` / `--- FAIL:` lines.

**Compile fence.** `go build ./...` does **not** compile `_test.go` files. It returns rc 0 even with
a type error in a test file. The compile fence is `go test -run '^$' ./host/broker/` (or
`go vet ./host/broker/`).

**Executor constraints (binding).**
1. **No git write operations of any kind:** no `commit`, `add`, `stash`, `checkout -- file`,
   `restore` or `reset`. The controller commits. `git diff`, `git status` and `git show` are
   allowed.
2. Edit only the files listed in *Scope*, and only at the sites this plan names.
3. Never edit `tools/launchd/*`. Never touch the main checkout
   `/Users/voightkampff/dev/sunholo-data/ailang-world` or `~/.ailang/state/mission-v1*`.
4. Apply each mutation **alone**, observe it red **on its named killer rows** (not the whole suite),
   and revert it byte-identically from a snapshot, checked with `shasum -a 256 -c` (§6).
5. Do not add `t.Parallel()`, and do not set `WORLD_EXEC_SRT_NODE_MODULES`.
6. **Bound every wait** with a wall-clock deadline: wrap long runs in
   `perl -e 'alarm N; exec @ARGV' …`. Never end your turn announcing a wait.

**The controller (not the executor) runs** the full `./scripts/verify_go.sh` (both legs) and the
AC-LINUX-LOAD Linux CI probe (§5). It also reads AC-LINUX on the PR's CI. The executor's local
gate is G1–G9.

---

## 2. Gate list

| ID | Command | Pass |
|---|---|---|
| G1 | `gofmt -l host/broker/` | empty output |
| G2 | `go vet ./host/broker/` | rc 0 |
| G3 | `go test -run '^$' ./host/broker/` | rc 0, `ok … [no tests to run]` (**the compile fence for `_test.go`**) |
| G4 | `go build ./... && go vet ./host/broker/ && go test ./host/broker/ -count=1` | rc 0, `ok  	github.com/sunholo-data/ailang-world/host/broker` |
| G5 | `go test -race ./host/broker/ -run '^TestExecProbeReverse' -count=1` | rc 0, no `DATA RACE` |
| G6 | AC-STABLE (M2): see the code block in M2 | 8 names, each `--- PASS` exactly 20 times, 0 `--- FAIL` |
| G7 | `./scripts/verify_ail.sh` | rc 0 (no `.ail` is touched; the gate must still be green) |
| G8 | `actionlint .github/workflows/ci.yml` | rc 0, no output |
| G9 | `git diff --name-only` | after M1: exactly the two Go files; after M2: exactly the six *Scope* files |

Copy the `-run` regexes from the fenced code blocks. In a markdown table a `\|` renders as `|`.

```sh
# G1
gofmt -l host/broker/
# G2
go vet ./host/broker/; echo "rc=$?"
# G3 compile fence
go test -run '^$' ./host/broker/ > /tmp/w318-g3.txt 2>&1; echo "rc=$?"; cat /tmp/w318-g3.txt
# G4 (about 170-300 s)
perl -e 'alarm 900; exec @ARGV' sh -c 'go build ./... && go vet ./host/broker/ && go test ./host/broker/ -count=1' > /tmp/w318-g4.txt 2>&1; echo "rc=$?"; tail -3 /tmp/w318-g4.txt
# G5
perl -e 'alarm 600; exec @ARGV' go test -race ./host/broker/ -run '^TestExecProbeReverse' -count=1 > /tmp/w318-g5.txt 2>&1; echo "rc=$?"; grep -c 'DATA RACE' /tmp/w318-g5.txt; tail -2 /tmp/w318-g5.txt
# G7
perl -e 'alarm 900; exec @ARGV' ./scripts/verify_ail.sh > /tmp/w318-g7.txt 2>&1; echo "rc=$?"; tail -3 /tmp/w318-g7.txt
# G8
actionlint .github/workflows/ci.yml; echo "rc=$?"
# G9
git diff --name-only
```

---

## 3. Pristine-tree baselines (record these BEFORE the first edit)

```sh
mkdir -p /tmp/w318-exec
shasum -a 256 host/broker/exec_probe.go host/broker/exec_probe_test.go .github/workflows/ci.yml docs/QUICKSTART.md \
  design_docs/planned/w-exec-probe-reset-holder-false-pass.md design_docs/implemented/w-exec-probe-reverse-leg-port-race.md \
  > /tmp/w318-exec/base.sha
cat /tmp/w318-exec/base.sha
go test ./host/broker/ -run '^TestExecProbeReverse' -count=1 -v > /tmp/w318-exec/focused-base.txt 2>&1; echo "rc=$?"
grep '^--- \|^ok\|^FAIL' /tmp/w318-exec/focused-base.txt
```

Expected (planner, this session):

| File | sha256 prefix |
|---|---|
| `host/broker/exec_probe.go` | `6e1b65365459` |
| `host/broker/exec_probe_test.go` | `8c9cb59872ba` |
| `.github/workflows/ci.yml` | `c40e8d653f0f` |
| `docs/QUICKSTART.md` | `a7dc9c61a5e7` |
| `design_docs/planned/w-exec-probe-reset-holder-false-pass.md` | `7477c3529f97` |
| `design_docs/implemented/w-exec-probe-reverse-leg-port-race.md` | `12c6ec7be333` |

Focused base run: rc 0, 8 s, **6** `--- PASS` (`…LegJudge`, `…ChattyForeignHolderIsRetried`,
`…ExhaustionRefusesDistinctly`, `…NoTokenIsNotRetried`, `…SilentHolderRefusesSticky`,
`…RealReachSurvivesAForeignAttempt`), 0 FAIL, 0 SKIP. G1 empty, G2 rc 0, G3 rc 0 (1 s).

If any sha256 differs, **stop** and report: the tree is not the planned base.

### 3.1 Base RED / GREEN

- **Red at base:** T7 `…/reset-holder` is red on **Linux under load** (≥7/300, design V14), and
  rarely on darwin (0/300 in V8). It is not a reliable local red. The deterministic red is M22
  (today's behaviour on the new seam) against T9 `connect-reset`.
- **Cannot compile at base:** T8 and T9 (they call `dialErrIsSilent` and the 3-argument
  `pollLoopback`).
- **Green at base, must stay green:** every `TestExecProbe*` test, `TestQuickstartSection10ProfilesLoad`
  and the whole `host/broker` package.

---

## 4. Milestones

### Milestone M1 — the classifier, the dial seam, T8 and T9 (+19/−6 prod, +93 test)

**Files:** `host/broker/exec_probe.go`, `host/broker/exec_probe_test.go`.

#### M1.1 `host/broker/exec_probe.go`: four edits

**(a) Import.** In the import block (`:26-46`), insert `"syscall"` between `"sync/atomic"` (`:43`)
and `"time"` (`:44`):

```go
	"sync/atomic"
	"syscall"
	"time"
```

**(b) The call site** (`:566`, the only caller; census V-P5). Replace

```go
		poll := pollLoopback(port, "WORLD-PROBE-NONCE "+nonce)
```

with

```go
		poll := pollLoopback(port, "WORLD-PROBE-NONCE "+nonce, nil)
```

**(c) The classifier and the seam.** Replace the doc comment and signature at `:690-694`:

```go
// pollLoopback dials 127.0.0.1:port every 50 ms until stopped and
// classifies every connected dial (classifyReverseDial); the returned stop
// function reports what the dials saw. It dials once before it first checks
// for stop, so an attempt always observes at least one dial.
func pollLoopback(port int, want string) func() reverseSeen {
```

with exactly:

```go
// dialErrIsSilent reports whether a failed reverse-leg dial proves a listener accepted the
// handshake and reset it (row 232): only ECONNRESET. ECONNREFUSED (nothing listens), EINVAL
// (seen on darwin with nothing listening), timeouts and local errors record nothing.
func dialErrIsSilent(err error) bool { return errors.Is(err, syscall.ECONNRESET) }

// pollLoopback dials 127.0.0.1:port every 50 ms until stopped and
// classifies every connected dial (classifyReverseDial); a dial that fails
// with ECONNRESET is silent (dialErrIsSilent). The returned stop function
// reports what the dials saw. It dials once before it first checks for
// stop, so an attempt always observes at least one dial. dial nil means
// net.DialTimeout("tcp", addr, 200ms); tests inject dial errors through it.
func pollLoopback(port int, want string, dial func(addr string) (net.Conn, error)) func() reverseSeen {
	if dial == nil {
		dial = func(addr string) (net.Conn, error) { return net.DialTimeout("tcp", addr, 200*time.Millisecond) }
	}
```

`dialErrIsSilent` must stay a **one-line body**, because the mutants in §M1.4 replace that exact line.

**(d) The dial and the error branch** (`:705` and `:714-715`). Replace

```go
			if c, err := net.DialTimeout("tcp", addr, 200*time.Millisecond); err == nil {
```

with

```go
			if c, err := dial(addr); err == nil {
```

and replace the closing of that `if` (the lines right after the inner `switch`):

```go
				_ = c.Close()
			}
			select {
```

with

```go
				_ = c.Close()
			} else if dialErrIsSilent(err) {
				silent.Store(true)
			}
			select {
```

Leave everything else untouched: the cadence (`time.After(50 * time.Millisecond)`), the 200 ms
timeout, `classifyReverseDial`, `judgeReverse` and every refusal text.

#### M1.2 `host/broker/exec_probe_test.go`: the import and three additions

**(a) Import.** Insert `"syscall"` between `"sync/atomic"` and `"testing"` in the import block.
`context`, `errors`, `fmt`, `net`, `os`, `strconv`, `strings`, `sync/atomic` and `time` are already
imported.

**(b) Append at the end of the file**, verbatim:

```go

// revDialErr is the error net.DialTimeout returns when connect(2) fails with
// errno (Go's non-blocking connect: *net.OpError -> *os.SyscallError -> Errno).
func revDialErr(errno syscall.Errno) error {
	return &net.OpError{Op: "dial", Net: "tcp", Err: os.NewSyscallError("connect", errno)}
}

// T8 (row 232 design §6): only ECONNRESET proves a listener completed the
// handshake; every other dial error records nothing. real-refused and
// timeout are real dials, each checked against its known shape first.
func TestExecProbeReverseDialErrIsSilent(t *testing.T) {
	port, err := freeLoopbackPort()
	if err != nil {
		t.Fatal(err)
	}
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	_, realRefused := net.DialTimeout("tcp", addr, time.Second)
	if !errors.Is(realRefused, syscall.ECONNREFUSED) {
		t.Fatalf("control: a dial to the freed port %s gave %v, want ECONNREFUSED", addr, realRefused)
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	var d net.Dialer
	_, realTimeout := d.DialContext(ctx, "tcp", addr)
	var ne net.Error
	if !errors.As(realTimeout, &ne) || !ne.Timeout() {
		t.Fatalf("control: an expired-context dial gave %v, want a net.Error timeout", realTimeout)
	}
	for _, row := range []struct {
		name string
		err  error
		want bool
	}{
		{"op-reset", revDialErr(syscall.ECONNRESET), true},
		{"bare-reset", syscall.ECONNRESET, true},
		{"wrapped-reset", fmt.Errorf("x: %w", revDialErr(syscall.ECONNRESET)), true},
		{"nil", nil, false},
		{"op-refused", revDialErr(syscall.ECONNREFUSED), false},
		{"real-refused", realRefused, false},
		{"op-einval", revDialErr(syscall.EINVAL), false},
		{"op-aborted", revDialErr(syscall.ECONNABORTED), false},
		{"addrnotavail", revDialErr(syscall.EADDRNOTAVAIL), false},
		{"timeout", realTimeout, false},
		{"deadline-exceeded", &net.OpError{Op: "dial", Net: "tcp", Err: os.ErrDeadlineExceeded}, false},
	} {
		t.Run(row.name, func(t *testing.T) {
			if got := dialErrIsSilent(row.err); got != row.want {
				t.Fatalf("dialErrIsSilent(%v) = %v, want %v", row.err, got, row.want)
			}
		})
	}
}

// T9 (row 232 design §6): pollLoopback records a connect-time reset as
// silent, through its dial seam and with no timing: the poll dials once
// before it first checks stop. Calls after the first return ECONNREFUSED,
// so extra dials before stop are harmless.
func TestExecProbeReversePollRecordsConnectReset(t *testing.T) {
	port, err := freeLoopbackPort()
	if err != nil {
		t.Fatal(err)
	}
	wantAddr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	refused := revDialErr(syscall.ECONNREFUSED)
	for _, row := range []struct {
		name  string
		first error
		want  reverseSeen
	}{
		{"connect-reset", revDialErr(syscall.ECONNRESET), reverseSeen{silent: true}},
		{"refused-only", refused, reverseSeen{}},
		{"einval-then-refused", revDialErr(syscall.EINVAL), reverseSeen{}},
	} {
		t.Run(row.name, func(t *testing.T) {
			var calls atomic.Int32
			var badAddr atomic.Bool
			fake := func(addr string) (net.Conn, error) {
				if addr != wantAddr {
					badAddr.Store(true)
				}
				if calls.Add(1) == 1 {
					return nil, row.first
				}
				return nil, refused
			}
			got := pollLoopback(port, "WORLD-PROBE-NONCE "+strings.Repeat("0", 32), fake)()
			if got != row.want || calls.Load() < 1 || badAddr.Load() {
				t.Fatalf("seen = %+v after %d dials (bad addr %v); want %+v after >=1 dial to %s", got, calls.Load(), badAddr.Load(), row.want, wantAddr)
			}
		})
	}
}
```

Why the details matter:
- `real-refused` and `timeout` are real dials. Each is checked against its known-positive shape
  (ECONNREFUSED; a `net.Error` timeout) **before** the table runs, so neither row can pass
  vacuously on a wrong shape.
- `deadline-exceeded` is needed because a real dial timeout does **not** satisfy
  `errors.Is(err, os.ErrDeadlineExceeded)`. Measured: `is-ErrDeadline=false`,
  `is-ctxDeadline=true`, `Timeout()=true` (D1). Without the row, mutant M28b would survive.
- T9's `calls.Load() < 1` and `badAddr` checks keep `refused-only` and `einval-then-refused` from
  passing if the seam were bypassed.
- The fake's `atomic.Int32` (gemini) makes call 1 return the row's error and every later call
  `op-refused`. So the extra dials the 50 ms loop may make before stop are harmless.

#### M1.3 Acceptance (M1)

```sh
gofmt -l host/broker/                                   # G1: empty
go vet ./host/broker/; echo "rc=$?"                     # G2: rc 0
go test -run '^$' ./host/broker/ > /tmp/w318-exec/g3.txt 2>&1; echo "rc=$?"   # G3: rc 0
go test ./host/broker/ -run '^TestExecProbeReverse(DialErrIsSilent|PollRecordsConnectReset)$' -count=1 -v > /tmp/w318-exec/t89.txt 2>&1; echo "rc=$?"
grep -c '^--- PASS' /tmp/w318-exec/t89.txt              # want 2
grep -c '^    --- PASS' /tmp/w318-exec/t89.txt          # want 14 (11 T8 rows + 3 T9 rows)
grep -c -- '--- FAIL\|--- SKIP' /tmp/w318-exec/t89.txt  # want 0
go test ./host/broker/ -run '^TestExecProbeReverse' -count=1 -v > /tmp/w318-exec/focused-m1.txt 2>&1; echo "rc=$?"
grep -c '^--- PASS' /tmp/w318-exec/focused-m1.txt       # want 8 (the 6 base names + T8 + T9)
grep -c -- '--- FAIL' /tmp/w318-exec/focused-m1.txt     # want 0
```

Then G4 and G5 (§2). `git diff --name-only` must list exactly `host/broker/exec_probe.go` and
`host/broker/exec_probe_test.go`. The planner's prototype gave a `diffstat` of `+19/−6` and `+93`,
and gofmt and vet were clean.

#### M1.4 Mutation drill M22–M29 (+M28b)

Protocol (§6) for each row: snapshot, apply **one** edit, run **only** the killer command in the
row, see every listed `--- FAIL:` line, restore, and check the sha. "Stays green" rows are in-table
controls: they must print `--- PASS:` under the mutant.

The edits name the line to replace. `RET` is the one line
`func dialErrIsSilent(err error) bool { return errors.Is(err, syscall.ECONNRESET) }`. `BR` is the
two lines `} else if dialErrIsSilent(err) {` / `silent.Store(true)` and the `}` that closes them.

| ID | Edit (apply alone) |
|---|---|
| M22 MUT-DIALERR-DROPPED | delete `BR` from `} else if` up to its closing `}`, leaving only the `}` that closes `if c, err := dial(addr); err == nil {` (today's behaviour) |
| M23 MUT-RESET-NOT-SILENT | `RET` → `func dialErrIsSilent(err error) bool { return errors.Is(err, syscall.ECONNRESET) && false }` |
| M24 MUT-RESET-NO-UNWRAP | `RET` → `func dialErrIsSilent(err error) bool { return err == syscall.ECONNRESET }` |
| M25 MUT-REFUSED-IS-SILENT | `RET` → `func dialErrIsSilent(err error) bool { return errors.Is(err, syscall.ECONNRESET) \|\| errors.Is(err, syscall.ECONNREFUSED) }` |
| M26 MUT-ANY-ERR-IS-SILENT | `RET` → `func dialErrIsSilent(err error) bool { return err != nil \|\| errors.Is(err, syscall.ECONNRESET) }` |
| M27 MUT-EINVAL-IS-SILENT | `RET` → `func dialErrIsSilent(err error) bool { return errors.Is(err, syscall.ECONNRESET) \|\| errors.Is(err, syscall.EINVAL) }` |
| M28 MUT-TIMEOUT-IS-SILENT | `RET` → `func dialErrIsSilent(err error) bool { var ne net.Error; return errors.Is(err, syscall.ECONNRESET) \|\| (errors.As(err, &ne) && ne.Timeout()) }` |
| M28b MUT-DEADLINE-IS-SILENT (planner-added, D1) | `RET` → `func dialErrIsSilent(err error) bool { return errors.Is(err, syscall.ECONNRESET) \|\| errors.Is(err, os.ErrDeadlineExceeded) }` |
| M29 MUT-CONNECT-RESET-IS-FOREIGN | inside `BR`, `silent.Store(true)` → `foreign.Store(true)` |

(In the table, `\|\|` is Go's `||`.) **M23 and M26 keep `syscall` referenced on purpose.** The
planner's first forms, `return false` and `return err != nil`, failed to compile
(`"syscall" imported and not used`). A mutant that does not compile is a **protocol error, never a
kill**. Do not run `gofmt` on a mutant: `gofmt -l` flags the long M25 one-liner, which does not matter,
because the snapshot restores the file.

**Killer commands.** Copy them from this block. Each `-run` targets exactly the named rows.

```sh
T8=TestExecProbeReverseDialErrIsSilent; T9=TestExecProbeReversePollRecordsConnectReset
k() { go test ./host/broker/ -run "$1" -count=1 -v > /tmp/w318-exec/mut-$2.txt 2>&1; echo "rc=$?"; grep -- '--- FAIL:\|--- PASS:\|build failed\|imported and not used' /tmp/w318-exec/mut-$2.txt; }
# M22: want FAIL T9/connect-reset; control PASS T9/refused-only
k "^$T9\$/^(connect-reset|refused-only)\$" M22
# M23: want FAIL T8/op-reset, T8/bare-reset, T8/wrapped-reset, T9/connect-reset
k "^$T8\$/^(op-reset|bare-reset|wrapped-reset)\$" M23a; k "^$T9\$/^connect-reset\$" M23b
# M24: want FAIL T8/op-reset, T8/wrapped-reset; control PASS T8/bare-reset
k "^$T8\$/^(op-reset|wrapped-reset|bare-reset)\$" M24
# M25: want FAIL T8/op-refused, T8/real-refused, T9/refused-only, T5/clean, every T2 row
k "^$T8\$/^(op-refused|real-refused)\$" M25a; k "^$T9\$/^refused-only\$" M25b
k '^TestExecProbeReverseLegNoTokenIsNotRetried$/^clean$' M25c
k '^TestExecProbeReverseLegChattyForeignHolderIsRetried$' M25d      # all 8 T2 rows FAIL
# M26: want FAIL T8/op-refused, op-einval, op-aborted, timeout, addrnotavail, and T5/clean; control PASS T8/nil
k "^$T8\$/^(op-refused|op-einval|op-aborted|timeout|addrnotavail|nil)\$" M26a
k '^TestExecProbeReverseLegNoTokenIsNotRetried$/^clean$' M26b
# M27: want FAIL T8/op-einval, T9/einval-then-refused
k "^$T8\$/^op-einval\$" M27a; k "^$T9\$/^einval-then-refused\$" M27b
# M28: want FAIL T8/timeout, T8/deadline-exceeded
k "^$T8\$/^(timeout|deadline-exceeded)\$" M28
# M28b: want FAIL T8/deadline-exceeded; T8/timeout stays PASS (this is why the row exists, D1)
k "^$T8\$/^(deadline-exceeded|timeout)\$" M28b
# M29: want FAIL T9/connect-reset (seen = {foreign}, not {silent})
k "^$T9\$/^connect-reset\$" M29
```

Each `k` call must show its listed `--- FAIL:` subtest lines (the parent test line `--- FAIL:` too),
and an rc ≠ 0. A `build failed` line means a protocol error: restore, fix the mutant's form and
re-run.

**Planner's prototype results** (scratch copy, same code; §9 V-P8): every listed killer FAILED,
and every control PASSED. In detail:
- M22: connect-reset FAIL, refused-only PASS.
- M23: 4/4 FAIL.
- M24: op-reset and wrapped-reset FAIL, bare-reset PASS.
- M25: 6/6 probed rows FAIL, and the whole of T2 run under M25 gave 8 of 8 subtests FAIL.
- M26: 6/6 FAIL, nil PASS.
- M27: 2/2 FAIL.
- M28: 2/2 FAIL.
- M28b: deadline-exceeded FAIL, timeout PASS.
- M29: FAIL.

After the last mutant, `shasum -a 256 -c /tmp/w318-exec/m1.sha` must print `OK` for both Go files.

---

### Milestone M2 — CI list, QUICKSTART clause, the reviewers' doc amendments, AC-STABLE, gates (±1 CI, +1 docs, about +12 doc lines)

**Files:** `.github/workflows/ci.yml`, `docs/QUICKSTART.md`,
`design_docs/planned/w-exec-probe-reset-holder-false-pass.md`,
`design_docs/implemented/w-exec-probe-reverse-leg-port-race.md`.

#### M2.1 `.github/workflows/ci.yml:247` (step 8 `tests=` list)

Line 247 is the parent's reverse-leg line. It ends with
`TestExecProbeReverseLegSilentHolderRefusesSticky"`, a string that occurs exactly once in the file
(V-P6). Replace that ending with:

```
TestExecProbeReverseLegSilentHolderRefusesSticky TestExecProbeReverseDialErrIsSilent TestExecProbeReversePollRecordsConnectReset"
```

Change nothing else. Step 8 (`:238-261`) then runs both new tests with `-v` on linux and fails
unless each prints `^--- PASS: <name> ` (`:254-257`).

#### M2.2 `docs/QUICKSTART.md:1153` (§10 arm 6)

Line 1153 currently reads:

```
   process that never answers held the port, and World cannot tell it from the child: find it.
```

Replace it with these **two** lines (3-space indent, inside list item 6):

```
   process that never answers (or that resets the host's dial) held the port, and World cannot
   tell it from the child: find it.
```

The backticked refusal text on `:1152` is the code's literal string and stays byte-identical (D3).

#### M2.3 The design doc's four narrow amendments (`design_docs/planned/w-exec-probe-reset-holder-false-pass.md`)

These are the only edits allowed to this file. Use the Edit tool with the exact old strings.

**(i) §2 timeout row (kimi).** In the row starting `| timeout (`, replace the cell's last sentence

```
The child's listener is dialled about 20 times per attempt (V15), so one lost dial cannot hide a reach |
```

with

```
The child's listener is dialled about 20 times per attempt (V15), so one lost dial cannot hide a reach. The reset holder is a single-dial shape (T7's `/bin/sh` client exits at once), so a lost timeout dial still hides it: T7 stays a stability read, and T9 is the proof |
```

**(ii) §6 T9 bullet (gemini).** After the sentence ending
`` `einval-then-refused` want `{}`. `` (the end of the **T9** bullet), append to that bullet:

```
  The fake dialer counts calls with an `atomic.Int32`: call 1 returns the row's error and every
  later call returns `op-refused`, so extra loop iterations before stop are harmless.
```

**(iii) §6 AC-LINUX (kimi; controller-verified).** Replace the whole bullet that begins
`- **AC-LINUX [r2].**` (4 lines, ending `amend §2's table with the observed class before any retry.`)
with:

```
- **AC-LINUX [r2, restated in the sprint].** The fix PR's linux/amd64 CI must show, in step 8
  ("Row 140 Workspace.Exec real-srt tests, verbose", the only `tests=` PASS-guard list,
  `ci.yml:245-257`), `--- PASS:` for T7, T8 and T9, plus the subtest line
  `    --- PASS: TestExecProbeReverseLegSilentHolderRefusesSticky/reset-holder` in its `-v` log. In
  the "go race gate" job (`go-race-verify`, `VERIFY_GO_LEGS: race`), which runs
  `go test ./... -count=1 -race` without `-v` and so prints only `ok`/`FAIL` per package, it must
  show `ok  	github.com/sunholo-data/ailang-world/host/broker`. A red in either blocks the merge:
  re-measure the Linux dial-error classes (V14's probe) and amend §2's table with the observed
  class before any retry.
```

**(iv) §9 V16 (glm).** Insert a new row directly after the `| V15 |` row:

```
| V16 | PR #402 (the merge that reddened `dev`) touched no probe file | `gh pr view 402 --repo sunholo-data/ailang-world --json files,mergeCommit --jq '{merge: .mergeCommit.oid, n: (.files\|length), probe: [.files[].path \| select(test("host/broker/exec_probe"))], go: [.files[].path \| select(test("\\.go$"))][0:2]}'` | merge `df5a0671dadd…`; `n` 29; `probe` `[]`; control `go` `["cmd/ailang-worldd/episode.go","cmd/ailang-worldd/episode_preset.go"]`, so the filter matches real paths |
```

(The planner ran this command this session; the observed column is the planner's output, V-P9.)

#### M2.4 The parent's amendment notes (`design_docs/implemented/w-exec-probe-reverse-leg-port-race.md`)

**Decision (glm's point; the planner checked precedent).** The parent's header carries only a
`Status` line, and no implemented doc in this repo has an "Amended by" header. But implemented docs
do get in-place post-implementation notes: `w-boundary-gate-tree-mutation.md:345` ("superseded"),
`w-bench-load-confound.md:1` (title clause SUPERSEDED) and commit `fd998400` (a status header
fixed after the move). So add **two** short notes, and change nothing else in the parent.

**(i) Header.** Directly after the `- **Status**: Implemented (iteration 314: …)` line (`:3`),
insert:

```
- **Amended by**: row 232, `w-exec-probe-reset-holder-false-pass` (iteration 318): §4.2's silent class also takes a dial that fails at connect with ECONNRESET.
```

**(ii) §4.2.** Directly after this line (`:188`, unique in the file):

```
  - zero bytes at the deadline, or a reset with zero bytes: set `seen.silent`.
```

insert:

```
    **Amended by row 232** (`w-exec-probe-reset-holder-false-pass`, iteration 318): a dial that
    fails *at connect* with ECONNRESET is the same TCP event and also sets `seen.silent`
    (`dialErrIsSilent`); every other dial error records nothing.
```

Cite the doc by name, not by path, so the note survives the controller's `planned/` →
`implemented/` move.

#### M2.5 Acceptance (M2)

```sh
G1; G2; G3                                                   # as §2
actionlint .github/workflows/ci.yml; echo "rc=$?"            # G8: rc 0
ruby -ryaml -e 'YAML.load_file(".github/workflows/ci.yml"); puts "ok"'
grep -c 'TestExecProbeReverseDialErrIsSilent TestExecProbeReversePollRecordsConnectReset"' .github/workflows/ci.yml   # want 1
go test ./host/broker/ -run '^TestQuickstartSection10ProfilesLoad$' -count=1 -v > /tmp/w318-exec/qs.txt 2>&1; echo "rc=$?"; grep '^--- ' /tmp/w318-exec/qs.txt   # PASS
# step-8 emulation (the new names only; real-srt names SKIP locally, so they are left out)
go test ./host/broker/ -count=1 -p 1 -timeout 600s -v \
  -run '^(TestExecProbeReverseLegSilentHolderRefusesSticky|TestExecProbeReverseDialErrIsSilent|TestExecProbeReversePollRecordsConnectReset)$' \
  > /tmp/w318-exec/step8.txt 2>&1; echo "rc=$?"
for t in TestExecProbeReverseLegSilentHolderRefusesSticky TestExecProbeReverseDialErrIsSilent TestExecProbeReversePollRecordsConnectReset; do
  grep -q -- "^--- PASS: $t " /tmp/w318-exec/step8.txt && echo "PASS $t" || echo "MISSING $t"
done
grep -c '^    --- PASS: TestExecProbeReverseLegSilentHolderRefusesSticky/reset-holder ' /tmp/w318-exec/step8.txt   # want 1
```

**AC-STABLE (G6).** This takes about 150 s idle (the planner's base run was 8 s per count).

```sh
perl -e 'alarm 900; exec @ARGV' go test ./host/broker/ -run '^TestExecProbeReverse' -count=20 -v > /tmp/w318-exec/c20.txt 2>&1; echo "rc=$?"
for t in TestExecProbeReverseLegJudge TestExecProbeReverseLegChattyForeignHolderIsRetried TestExecProbeReverseLegExhaustionRefusesDistinctly \
         TestExecProbeReverseLegNoTokenIsNotRetried TestExecProbeReverseLegSilentHolderRefusesSticky TestExecProbeReverseLegRealReachSurvivesAForeignAttempt \
         TestExecProbeReverseDialErrIsSilent TestExecProbeReversePollRecordsConnectReset; do
  printf '%s %s\n' "$(grep -c -- "^--- PASS: $t " /tmp/w318-exec/c20.txt)" "$t"     # each 20
done
grep -c -- '--- FAIL' /tmp/w318-exec/c20.txt                                          # want 0
grep -c '^    --- PASS: TestExecProbeReverseLegSilentHolderRefusesSticky/reset-holder ' /tmp/w318-exec/c20.txt   # want 20
```

**T7 under load (the design's stability read).** This is optional for the executor, and the
controller may run it instead. Start 32 `yes` with their PIDs recorded and a `trap` to kill them,
then run:

```sh
perl -e 'alarm 900; exec @ARGV' go test ./host/broker/ -run '^TestExecProbeReverseLegSilentHolderRefusesSticky$/^reset-holder$' -count=300 -v > /tmp/w318-exec/t7load.txt 2>&1; echo "rc=$?"
```

Kill the `yes` PIDs, then confirm none are alive (`pgrep -x yes | wc -l` reads 0). Count over the
**whole** log: 0 `--- FAIL: …/reset-holder`, and exactly 300 `    --- PASS: …/reset-holder`.

Then G4, G5, G7 and G9 (§2). G9 lists exactly the six *Scope* files.

**Hand-off to the controller:** G4 is not the full gate. The controller runs
`./scripts/verify_go.sh` (both legs) and AC-LINUX-LOAD (§5), and reads AC-LINUX on the PR's CI.

---

## 5. Local versus CI-only evidence (controller-owned)

- **AC-LINUX** (restated, M2.3 iii): on the PR's linux CI, step 8 must show `--- PASS:` for T7,
  T8 and T9 plus the `…/reset-holder` subtest PASS line. The race gate must show `ok` for
  `host/broker`.
- **AC-LINUX-LOAD** (design §6; **controller-run, before merge**): use V14's workflow pattern on a
  throwaway branch at the **fixed** head. Start 2×nproc `yes`, run
  `go test ./host/broker/ -run '^TestExecProbeReverseLegSilentHolderRefusesSticky$/^reset-holder$' -count=300 -v`,
  then kill the `yes` PIDs. Count over the whole log: 0 FAIL, exactly 300 PASS. The paired red is
  base ≥7 FAIL (V14).
- **Full `./scripts/verify_go.sh`** (plain and race legs): controller. A failure that appears
  **only** in `host/verifygate TestMissionConfigDefaultFlowIntegration` timing out under load is the
  pre-existing flake the parent recorded. Re-run that test alone; any other failure is a stop.

## 6. Mutation protocol (every mutant)

```sh
mkdir -p /tmp/w318-exec
cp host/broker/exec_probe.go /tmp/w318-exec/exec_probe.go.snap
shasum -a 256 host/broker/exec_probe.go host/broker/exec_probe_test.go > /tmp/w318-exec/m1.sha
#   apply ONE edit from the M1.4 table with the Edit tool
#   run ONLY that mutant's killer line(s) from the M1.4 block
cp /tmp/w318-exec/exec_probe.go.snap host/broker/exec_probe.go
shasum -a 256 -c /tmp/w318-exec/m1.sha        # both lines: OK
```

Take the snapshot **after** M1.1–M1.3 are green. Every mutant edits only `exec_probe.go`.

## 7. Risks

| Risk | Mitigation |
|---|---|
| T8 `real-refused` control on a loaded darwin host returns EINVAL (122 of 2,013,696 no-listener dials under 32 `yes`, design V7) | The control fails loudly with the observed error, never silently. Rate ≈6×10⁻⁵ per run under heavy load. Re-run once; a second hit is a finding, so report it |
| A mutant that does not compile looks like a kill | `k()` greps `build failed` / `imported and not used`. M23 and M26 are written to keep `syscall` used |
| T7 is not strictly deterministic (a 200 ms dial timeout records nothing; Linux 143/179 under 4 spinning dialers, V14) | T9 is the deterministic proof. T7 is a stability read (AC-STABLE, AC-LINUX-LOAD) |
| G4 reds only on `TestWorktreeSnapshotRefusesAConfiguredFilter` (`worktree_delta_test.go:432`, a git clean-filter control) under host load (planner saw it once, V-P11) | Re-run that test alone with `-count=3`. PASS alone means an unrelated load flake: record it, do not touch `worktree_delta*`. Any other failure is a stop. Do not run G4 while a mutant or AC-STABLE run is going |
| R5 foreign-holder exit is classed silent (pre-existing) | Out of scope (design §8, successor `w-exec-probe-holder-exit-reset`) |

## 8. Definition of done

1. M1: G1–G5 green. T8 shows 11 subtest PASS and T9 shows 3. Every M22–M29 and M28b killer row is
   red under its mutant, every in-table control is green, and both Go files restore OK.
2. M2: the CI line is edited (G8 and YAML ok). The QUICKSTART clause is in, and
   `TestQuickstartSection10ProfilesLoad` PASSes. The step-8 emulation PASSes all three names.
   AC-STABLE is 20/20 for each of the 8 names. G4, G5 and G7 are green. G9 shows the six files.
3. The executor's report lists every command, rc and count. The controller then runs §5.

## 9. Design discrepancies found by the planner (resolved in-plan, flagged to the controller)

- **D1 (false in the design: M28's killer cannot see the `os.ErrDeadlineExceeded` form).** The
  §2 timeout row names `net.Error.Timeout()` and `os.ErrDeadlineExceeded`, and §5 says T8
  `timeout` (a real expired-context dial) kills M28. But measured, that real error satisfies
  `Timeout()` and `errors.Is(err, context.DeadlineExceeded)`, and **not**
  `errors.Is(err, os.ErrDeadlineExceeded)` (V-P7). Go maps a context deadline to `errTimeout`,
  so production `DialTimeout` timeouts have the same shape. **Resolved:** T8 gains a
  `deadline-exceeded` row (a synthesised `&net.OpError{…, Err: os.ErrDeadlineExceeded}`), M28 is
  pinned to the `net.Error` form (killed by both rows), and M28b, the `os.ErrDeadlineExceeded`
  form, is added (killed only by `deadline-exceeded`; prototype-verified). This is not one of the
  three permitted design amendments, so the design doc's §5/§6 text is left as is for the
  controller to note.
- **D2 (AC-LINUX as written is unsatisfiable).** The race gate has no `-v` and no `tests=` guard
  (`ci.yml:418` job, `:497` step, `VERIFY_GO_LEGS: race`; `scripts/verify_go.sh:355`
  `cmd = ["go", "test", "./...", "-count=1", "-race", "-timeout", "8m"]`). **Resolved:** restated
  per kimi in M2.3 (iii).
- **D3 (QUICKSTART placement).** Design §7 says "answered nothing to World's hello" gains "(or
  reset the host's dial)". That phrase sits inside the backticked literal refusal text
  (`:1152`), which the code still prints unchanged. **Resolved:** the clause goes outside the code
  span, on `:1153` ("never answers (or that resets the host's dial)"). The design's `:1152` is
  the line that holds the quoted text.
- **D4 (cosmetic citation).** Design §6 cites the "dials once before it checks stop" property as
  `:690-693`. That is the doc comment. The code is the `for {` loop at `:704-720` (dial `:705`,
  `select` `:716-720`). It is true as read.
- **D5 (imports).** Neither `exec_probe.go` nor `exec_probe_test.go` imports `syscall` today (0
  matches, V-P4). Both gain it. The prod diff is `+19/−6` (the design estimated `+12/−3`), because
  the doc comment is rewritten to name the new error branch and the seam.
- **D6 (AC-MUT "each reds exactly its named killer").** That is not the protocol here. Per the
  mission rule, each mutant is run **against the rows that name it**, not the suite. M25 also reds
  T3 `child-eaddrinuse` and the T7 rows with free ports, which is expected and not recorded as a
  defect.
- **D7 (AC-GATE scope).** The design's AC-GATE says `go build ./... && go test ./...`. The
  executor runs `go build ./... && go vet ./host/broker/ && go test ./host/broker/` plus the
  compile fence and `verify_ail.sh`. The repo-wide `go test ./...` is in the controller's
  `verify_go.sh` (§5).
- **D8 (step 8 guards top-level names only).** The `--- PASS:` guard (`:254-257`) greps
  `^--- PASS: $t `, so it does not guard the `…/reset-holder` subtest line. A failing subtest
  fails its parent, and T7's rows never skip, so the top-level PASS implies it. The subtest line is
  read from the log (M2.5 and AC-LINUX). No extra CI guard is added; the design's CI change stays
  one line.

**Nothing else in the design was found false.** These checks held:
- the caller census (13 lines, one `pollLoopback` call site at `:566`; control `runBounded` = 9);
- `:84`, `:98`, `:112`, `:128`, `:660`, `:694`, `:700`, `:705`, `:706` and `:719`;
- T7's failure line `:664`;
- the reset holder `SetLinger(0)` at `:495`;
- `ci.yml:247` as the reverse-leg list line;
- every T8 shape in V10 (re-run).

## 10. Verification notes (planner, this session, worktree `e1110207`)

| # | Claim | Command | Observed |
|---|---|---|---|
| V-P1 | Base: clean, and the Go tree equals `df5a0671` | `git status --porcelain \| wc -l; git rev-parse HEAD; git diff --stat df5a0671 HEAD` | `0`; `e1110207be7f…`; 13 files, all under `design_docs/` |
| V-P2 | Toolchain | `go version` (in the worktree); `$AILANG_BIN --version` | `go1.26.6 darwin/arm64`; `AILANG v0.41.0` |
| V-P3 | Baselines | §3 commands | the sha prefixes in §3; gofmt empty; vet rc 0; fence rc 0; focused 6 PASS 0 FAIL 0 SKIP, 8 s |
| V-P4 | Edit sites | `grep -n 'syscall\|func pollLoopback\|poll := pollLoopback\|net.DialTimeout\|time.After(50' host/broker/exec_probe.go`; `grep -c syscall` on both files | `:566`, `:694`, `:705`, `:719`; syscall `0` and `0` |
| V-P5 | Census | `grep -rn 'pollLoopback\|classifyReverseDial\|judgeReverse' host/ \| wc -l`; `grep -rn 'pollLoopback(' host/`; the `runBounded` control | `13`; call `:566` and def `:694` only; control `9` |
| V-P6 | CI and QUICKSTART anchors | `sed -n 236,262p .github/workflows/ci.yml`; `grep -c 'TestExecProbeReverseLegSilentHolderRefusesSticky"' ci.yml`; `sed -n 410,420p;490,505p ci.yml`; `sed -n 350,362p scripts/verify_go.sh`; `sed -n 1152,1153p docs/QUICKSTART.md` | step 8 `:238`, list `:244-248`, guard `:254-257`; anchor count `1`; `go-race-verify` job `:418`, step `VERIFY_GO_LEGS: race`; race cmd without `-v`; QUICKSTART `:1152` holds the quoted text, `:1153` "process that never answers held the port, …" |
| V-P7 | Error shapes (scratch `/tmp/w318-planner/is`, go1.26.4, cross-compiled for linux/amd64, rc 0) | `go run .` | `real-refused … connect: connection refused is-refused=true is-reset=false`; `timeout … i/o timeout as-net.Error=true Timeout=true is-ErrDeadline=false is-ctxDeadline=true`; `op-reset is=true eq=false`; `bare-reset is=true eq=true`; `wrapped-reset is=true eq=false`; refused, einval, aborted and addrnotavail `is-reset=false`; nil `false`; synthetic `OpError{Err: os.ErrDeadlineExceeded}`: `is-os=true as-timeout=true` |
| V-P8 | The whole M1 prototyped (scratch copy `/tmp/w318-planner/proto`, rsync of the worktree minus `.git`, so the repo is untouched) | M1.1/M1.2 edits applied by script; gofmt; vet; T8+T9 `-v`; mutation runner `/tmp/w318-planner/mut.py` + `mut2.py` | prod `+19/−6`, test `+93`; gofmt empty, vet ok; 2 PASS + 14 subtest PASS; every mutant's named killer FAIL and every control PASS (M1.4); the first M23/M26 forms failed to compile (`"syscall" imported and not used`), so they were rewritten; the scratch file restored `cmp`-identical; the scratch copy was deleted after V-P11 (the runner scripts stay in `/tmp/w318-planner/`) |
| V-P9 | PR #402 touched no probe file (glm's V16) | `gh pr view 402 --repo sunholo-data/ailang-world --json files,mergeCommit,title --jq …` | merge `df5a0671dadd567e…`; 29 files; `probe: []`; control `.go` paths present (`cmd/ailang-worldd/episode.go`, …) |
| V-P10 | Amendment precedent in implemented docs | `grep -n -i 'superseded' design_docs/implemented/*.md`; `git log … -- design_docs/implemented/` | `w-boundary-gate-tree-mutation.md:345`, `w-bench-load-confound.md:1`; commit `fd998400` (a status header fixed post-move); no "Amended by" header anywhere (the new form starts here) |
| V-P11 | Prototype stability and race | in the scratch copy: AC-STABLE `-count=20 -v` block (M2.5); `go test -race ./host/broker/ -run '^TestExecProbeReverse' -count=1` | `-count=20`: rc 0, 98 s, each of the 8 names `--- PASS` 20, 0 FAIL, 20 `…/reset-holder` PASS; race rc 0, 0 `DATA RACE`, `ok … 6.028s`. Full `go test ./host/broker/ -count=1` on the prototype: rc 1 in 198 s, and the ONLY failure was `TestWorktreeSnapshotRefusesAConfiguredFilter` (`worktree_delta_test.go:432`, "control: plain git status did not run the clean filter"), while a mutant run loaded the host. Re-run alone `-count=3`: PASS 3/3 at base AND on the prototype, so it is an unrelated load flake |

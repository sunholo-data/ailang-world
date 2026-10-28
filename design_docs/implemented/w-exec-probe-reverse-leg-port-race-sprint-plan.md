# Sprint plan — `w-exec-probe-reverse-leg-port-race` (iteration 314, clause 3, availability defect)

**Design doc:** [`design_docs/implemented/w-exec-probe-reverse-leg-port-race.md`](w-exec-probe-reverse-leg-port-race.md)
(r3; quorum-cleared through the narrow-refinement carve-out after rounds 1 and 2, §13)
**Sprint base:** worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter314-arm6`, branch
`sprint/world-iter314-arm6-reverse-leg`, HEAD `e4dbac36`, clean tree. Its Go tree is the design's
measurement base `2769cff9`, since `72eba94b..e4dbac36` changed only the design doc.
**Planner:** mission-control iteration 314 (planner role). **Executor:** claude-haiku, so this plan
spells out every edit.
**Estimate:** 0.75 day. 3 milestones, each independently green. **Risk:** low to medium. The
mechanism is small, but the sockets in the test fixtures need care (§6).
**Scope:** `host/broker/exec_probe.go`, `host/broker/exec_probe_test.go`, `docs/QUICKSTART.md`
(§10 arm 6 prose), `.github/workflows/ci.yml` (step 8 `tests=`). **No `.ail` file. No other file.**

---

## 0. The one-paragraph version

Arm 6's reverse leg counts **any** connect on a port it freed a moment ago as "the host reached the
child's listener". srt's own mux proxy binds a host ephemeral port during Run C, and on CI one
collision refused daemon startup. The fix has four parts:
- the probe client's listener writes a fresh 128-bit nonce on accept;
- the host poll writes `WORLD-PROBE-HELLO\n` on every connected dial and classifies the answer as
  **child** (the nonce line), **foreign** (other bytes, or a clean close) or **silent** (zero bytes:
  deadline or reset);
- a pure `judgeReverse` turns one attempt into PASS / REACHED / NO-TOKEN / SILENT / INCONCLUSIVE;
- a `reverseLeg` loop retries INCONCLUSIVE on a fresh port, at most 3 attempts.

A connected dial never makes its own attempt pass. The work lands in three milestones:
- **M1:** the production change and the pure judge test.
- **M2:** the four `/bin/sh`-fixture tests (T2, T3, T5, T7), which need no node.
- **M3:** the two node tests, including the flake reproduction (red at base), the CI list and the
  QUICKSTART prose.

---

## 1. Environment — on EVERY command

```sh
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter314-arm6
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang     # reports "AILANG v0.41.0" (measured)
```

| Pin | Measured this session | Why it matters |
|---|---|---|
| `go` | `go version go1.26.6 darwin/arm64` | `scripts/verify_go.sh` refuses go1.26.0–1.26.5. go1.26.6 is the CI pin (`ci.yml` `GOTOOLCHAIN: go1.26.6`) |
| `node` | `/opt/homebrew/bin/node` | T4 and T6 need it (`probeNode`). They SKIP by name without it. Locally they RUN |
| `AILANG_BIN` | `AILANG v0.41.0` | Without it the daemon e2e tests SKIP on "AILANG_BIN unset" and `verify_go.sh` refuses to start |
| `WORLD_EXEC_SRT_NODE_MODULES` | **unset** | Every real-srt test SKIPs locally (§5). Do NOT set it |

**Instrument rule.** Never write `cmd | tail; echo rc=$?`, and never use `${PIPESTATUS[0]}` (zsh).
Redirect to a file, put `echo "rc=$?"` on the next line, then read the file. A `-run` that matches
nothing exits 0, so every verdict is read from counted `--- PASS:` / `--- FAIL:` / `--- SKIP:` lines.

**Executor constraints (binding).**
1. **No git write operations of any kind**: no `commit`, `add`, `stash`, `checkout -- file`,
   `restore` or `reset`. The controller commits. `git diff`, `git status` and `git show` are allowed.
2. Edit **only** the four files in *Scope*. The P0 scratch file (§3) is the single exception, and
   you must delete it in the same step that creates it.
3. Never edit `tools/launchd/*`. Never touch `/Users/voightkampff/dev/sunholo-data/ailang-world`
   (the main checkout) or `~/.ailang/state/mission-v1*`. Work only in the worktree above.
4. Each mutation is applied **alone**, observed red, and reverted **byte-identically** by file copy
   from a snapshot, then checked with `shasum -a 256 -c` (protocol in §7). `git diff --quiet` cannot
   check this, because the milestone's own uncommitted edits are in the diff.
5. Do not set `WORLD_EXEC_SRT_NODE_MODULES`, and do not add `t.Parallel()`. Neither package uses it
   (design V6), and CI step 8 runs them with `-p 1`.

---

## 2. Gate list (derived from `.github/workflows/ci.yml` as read this session)

What CI actually runs on this change. Each item was read with `grep -n` / `sed -n` and cited in §10:
- **go-verify job, step "Row 140 Workspace.Exec real-srt tests, verbose (linux bubblewrap)"**
  (`ci.yml:238-259`). It runs `go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -timeout 600s -v
  -run "^(${tests// /|})\$"` (`:250`). Then, for every name in `tests=` (`:245-247`), it fails unless
  the log has `^--- PASS: $t ` (`:254-256`). **This is the PASS-guard that keeps a SKIP from passing
  silently.**
- **go-verify job, step "go build + test gate"** (`ci.yml:274-280`). It runs `./scripts/verify_go.sh`
  with `VERIFY_GO_LEGS=plain`: build, the evidence manifest and `go test ./... -count=1`.
- **go-race job** (`ci.yml:418`, step `:496`). It runs `verify_go.sh` with `VERIFY_GO_LEGS=race`,
  which is `go test ./... -count=1 -race`. New tests must be race-clean.
- **ailang-verify job** (`ci.yml:113-114`). It runs `./scripts/verify_ail.sh`. That is not needed
  here, because no `.ail` changes.
- CI runs **no** `gofmt` and **no** `go vet` (`grep -n 'gofmt\|go vet' .github/workflows/*.yml
  scripts/*.sh` returns 0 lines). The local gates G1 and G2 are stricter than CI, on purpose.

**Run the gates at the end of every milestone (G1–G7), and the full set (G1–G13) after M3.**

> **Copy commands from the code block below, NOT from table cells.** Inside markdown tables a `|`
> has to be written `\|`. Pasting `\|` into a `-run` regex makes it a *literal* pipe, so the regex
> matches nothing and the run reads green (rc 0, `no tests to run`).

```sh
# G1
gofmt -l host/broker/
# G2
go vet ./host/broker/
# G3 compile fence
go test -run '^$' ./... > /tmp/w314-fence.txt 2>&1; echo "rc=$?"; grep -v '^ok\|no test files' /tmp/w314-fence.txt
# G4
go test ./host/broker/ -count=1 > /tmp/w314-broker.txt 2>&1; echo "rc=$?"
# G5
go test ./host/daemon/ -count=1 -run 'TestExecStartup|TestExecSrt' -v > /tmp/w314-daemon.txt 2>&1; echo "rc=$?"
echo "PASS=$(grep -c '^--- PASS' /tmp/w314-daemon.txt) FAIL=$(grep -c '^--- FAIL' /tmp/w314-daemon.txt) SKIP=$(grep -c '^--- SKIP' /tmp/w314-daemon.txt)"
# G6
go test ./host/broker -run 'TestExecProbe|TestQuickstartSection10' -count=1 -v > /tmp/w314-probe.txt 2>&1; echo "rc=$?"
grep -E '^(    )*--- (PASS|FAIL|SKIP)' /tmp/w314-probe.txt
# G7
go test -race ./host/broker -run 'TestExecProbe' -count=1 > /tmp/w314-race.txt 2>&1; echo "rc=$?"; grep -c 'WARNING: DATA RACE' /tmp/w314-race.txt
# G8 (after M2): node absent, no real srt
env -i HOME=$HOME PATH=/usr/bin:/bin WORLD_EXEC_SRT_NODE_MODULES=/nonexistent GOCACHE=$(go env GOCACHE) GOPATH=$(go env GOPATH) GOTOOLCHAIN=go1.26.6 \
  /opt/homebrew/bin/go test ./host/broker -run '^(TestExecProbeReverseLegJudge|TestExecProbeReverseLegChattyForeignHolderIsRetried|TestExecProbeReverseLegExhaustionRefusesDistinctly|TestExecProbeReverseLegNoTokenIsNotRetried|TestExecProbeReverseLegSilentHolderRefusesSticky)$' -count=1 -v > /tmp/w314-nonode.txt 2>&1; echo "rc=$?"
echo "PASS=$(grep -c '^--- PASS' /tmp/w314-nonode.txt) SKIP=$(grep -c -- '--- SKIP' /tmp/w314-nonode.txt)"
# G9 (after M3): AC-STABLE
perl -e 'alarm 900; exec @ARGV' go test ./host/broker -run '^TestExecProbeReverseLeg|^TestExecProbeForeignListener' -count=20 -v > /tmp/w314-stable.txt 2>&1; echo "rc=$?"
for t in TestExecProbeReverseLegJudge TestExecProbeReverseLegChattyForeignHolderIsRetried TestExecProbeReverseLegExhaustionRefusesDistinctly TestExecProbeReverseLegRealReachSurvivesAForeignAttempt TestExecProbeReverseLegNoTokenIsNotRetried TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse TestExecProbeReverseLegSilentHolderRefusesSticky; do
  echo "$t $(grep -c "^--- PASS: $t " /tmp/w314-stable.txt)"; done
grep -cE '^(    )*--- FAIL' /tmp/w314-stable.txt
# G10 (after M3)
VERIFY_GO_LEGS=plain perl -e 'alarm 1500; exec @ARGV' ./scripts/verify_go.sh > /tmp/w314-verifygo.txt 2>&1; echo "rc=$?"; tail -1 /tmp/w314-verifygo.txt
# G11 (after M3)
actionlint .github/workflows/ci.yml; echo "rc=$?"
ruby -ryaml -e 'YAML.load_file(".github/workflows/ci.yml"); puts "ok"'
# G12 (after M3): the step-8 emulation is in M3
# G13
git status --porcelain; git diff --name-only
```

| # | Command | Pass condition | Base reading (planner, pristine `e4dbac36`) |
|---|---|---|---|
| G1 | `gofmt -l host/broker/` | empty output | empty, rc 0 |
| G2 | `go vet ./host/broker/` | rc 0, no output | rc 0, no output |
| G3 | `go test -run '^$' ./... > /tmp/w314-fence.txt 2>&1; echo "rc=$?"` (the **compile fence**: `go build` does not compile `_test.go` files) | rc 0, and every line `ok` or `no test files` | rc 0, 0 other lines, 13.7 s |
| G4 | `go test ./host/broker/ -count=1 > /tmp/w314-broker.txt 2>&1; echo "rc=$?"` | rc 0 | rc 0, `ok … 167.649s` |
| G5 | `go test ./host/daemon/ -count=1 -run 'TestExecStartup\|TestExecSrt' -v > /tmp/w314-daemon.txt 2>&1; echo "rc=$?"` then count `^--- PASS` / `^--- FAIL` / `^--- SKIP` | rc 0, **PASS=2 FAIL=0 SKIP=2** (unchanged) | rc 0, PASS=2 (`TestExecSrtReplayIsByteEqual`, `TestExecStartupRefusalTable`), SKIP=2 (`TestExecSrtMCPEndToEnd`, `TestExecSrtMCPGrantAndBudget`: "real srt absent") |
| G6 | `go test ./host/broker -run 'TestExecProbe\|TestQuickstartSection10' -count=1 -v > /tmp/w314-probe.txt 2>&1; echo "rc=$?"` | rc 0, FAIL=0. The only SKIP is `TestExecProbeArm8` (real srt), plus the 2 host-GOOS subtest SKIPs. Every new test in the milestone shows `--- PASS` | rc 0, PASS=10 FAIL=0 SKIP=1 (`TestExecProbeArm8`), 2 subtest SKIPs (`linux-ebusy-passes`, `linux-exdev-passes`) for `-run 'TestExecProbe'` |
| G7 | `go test -race ./host/broker -run 'TestExecProbe' -count=1 > /tmp/w314-race.txt 2>&1; echo "rc=$?"` | rc 0, no `WARNING: DATA RACE` | not baselined (the existing tests are race-clean in CI's race job) |
| G8 | (after M2 and M3) **node-absent run**: `env -i HOME=$HOME PATH=/usr/bin:/bin WORLD_EXEC_SRT_NODE_MODULES=/nonexistent GOCACHE=$(go env GOCACHE) GOPATH=$(go env GOPATH) GOTOOLCHAIN=go1.26.6 /opt/homebrew/bin/go test ./host/broker -run '^(TestExecProbeReverseLegJudge\|TestExecProbeReverseLegChattyForeignHolderIsRetried\|TestExecProbeReverseLegExhaustionRefusesDistinctly\|TestExecProbeReverseLegNoTokenIsNotRetried\|TestExecProbeReverseLegSilentHolderRefusesSticky)$' -count=1 -v > /tmp/w314-nonode.txt 2>&1; echo "rc=$?"` | rc 0. `^--- PASS` = 5, and **0** `SKIP` lines of any depth | n/a: tests are new. The same `env -i` shape ran a `/bin/sh` proxy+listen fixture with rc 0 (§10 P6) |
| G9 | (after M3) **AC-STABLE**: `perl -e 'alarm 900; exec @ARGV' go test ./host/broker -run '^TestExecProbeReverseLeg\|^TestExecProbeForeignListener' -count=20 -v > /tmp/w314-stable.txt 2>&1; echo "rc=$?"` | rc 0. For each of the 7 new names, `grep -c "^--- PASS: <name> "` = **20**. `grep -c '^--- FAIL\|^    --- FAIL'` = 0 | n/a |
| G10 | (after M3) **CI go gate**: `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh > /tmp/w314-verifygo.txt 2>&1; echo "rc=$?"` | rc 0, last line `✓ go gate PASSED (plain)…`. **Rule:** if the ONLY failure is `host/verifygate` `TestMissionConfigDefaultFlowIntegration` timing out, re-run that test alone. If it then PASSES, record "pre-existing load flake, red at base too (V-GATE)" and do not touch `host/verifygate`. Any other failure is a stop | **rc 1 at base**: exactly that verifygate timeout under load ~15; it PASSES alone (V-GATE) |
| G11 | (after M3) `actionlint .github/workflows/ci.yml; echo "rc=$?"` and `ruby -ryaml -e 'YAML.load_file(".github/workflows/ci.yml"); puts "ok"'` | rc 0, no findings; `ok` | rc 0, 0 findings; `ruby-yaml ok` |
| G12 | (after M3) **step-8 emulation**: run step 8's exact `go test` line with the new `tests=` list (snippet in M3) and check that each of the **7 new names** prints `^--- PASS: <name> ` | 7/7 PASS lines. The other listed tests SKIP locally, which is expected (§5) | n/a |
| G13 | `git status --porcelain` and `git diff --name-only` | only the 4 scope files. **No `.ail` path.** If an `.ail` path ever appears, stop: `./scripts/verify_ail.sh` would become a gate, and the design says none is touched (V15) | 0 lines |

---

## 3. Pristine-tree baselines — record these BEFORE the first edit

Re-run G1–G6 and the commands below on the untouched tree, and record each rc and count. The
planner's readings are in §2 and §3.1. A different reading is a stop-and-report, not a fix-forward.

```sh
shasum -a 256 host/broker/exec_probe.go host/broker/exec_probe_test.go .github/workflows/ci.yml docs/QUICKSTART.md
# planner: 6eac0a12…88f4 exec_probe.go ; 45e89f3b…6d7b exec_probe_test.go ;
#          f53248a6…9b1 ci.yml ; 7e67fc48…85e7 QUICKSTART.md
go test ./host/broker -run '^TestExecProbeReverseLeg|^TestExecProbeForeignListener' -count=1 -v > /tmp/w314-new-base.txt 2>&1; echo "rc=$?"
# planner: rc 0, "testing: warning: no tests to run" (the names do not exist yet)
```

**P0 — the flake reproduction, measured at base (required).** T1–T5 and T7 call symbols that do not
exist at base, so they cannot compile there; their base reading is "absent". T6 uses only
`ProbeExecProfiles`, so its base reading **can** be taken. To take it:
1. Create `host/broker/zz_w314_t6base_test.go` holding §8's `revStandIn` constant and T6's body,
   renamed `TestZZW314T6Base`.
2. Run `perl -e 'alarm 300; exec @ARGV' go test ./host/broker -run '^TestZZW314T6Base$' -count=3 -v`.
3. Count the subtest lines.
4. **Delete the file**, and confirm `git status --porcelain` is empty.

**Expected base: `chatty` FAIL 3/3, `silent` FAIL 3/3, `control` PASS 3/3.** The planner ran exactly
this (§10 P5). In all six red runs, arm 6's Why ended `reverse leg: the host reached the listener the
child bound on 127.0.0.1:<port>`, the CI signature.

### 3.1 Base-RED / base-GREEN summary

| Acceptance test | Base |
|---|---|
| T1 `TestExecProbeReverseLegJudge` | absent (it calls `judgeReverse`) |
| T2 `TestExecProbeReverseLegChattyForeignHolderIsRetried` | absent (it calls `reverseLeg` and `pickPort`) |
| T3 `TestExecProbeReverseLegExhaustionRefusesDistinctly` | absent |
| T4 `TestExecProbeReverseLegRealReachSurvivesAForeignAttempt` | absent |
| T5 `TestExecProbeReverseLegNoTokenIsNotRetried` | absent |
| **T6 `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse`** | **`chatty` RED 3/3, `silent` RED 3/3, `control` GREEN 3/3** (P0) |
| T7 `TestExecProbeReverseLegSilentHolderRefusesSticky` | absent |
| P1 `TestExecProbePassThroughShimFailsEveryConfinementArm` | GREEN (1.87 s) and must stay GREEN |
| `TestExecProbeMissingClientFailsTheRefusalArms`, `…NoopShimFailsArm1`, `…ZeroShimFailsArm5`, `…FailingToolchainFailsArm7`, `…Arm8Wired`, `…Arm8FailsUnsandboxed`, `…Arm8Absent`, `…ProtectedSample`, `…ArmNames` | GREEN, and must stay GREEN |
| `TestExecProbeArm8` | SKIP locally (real srt). CI-only |
| G4 `./host/broker` whole package | GREEN |
| G5 daemon `TestExecStartup\|TestExecSrt` | GREEN, PASS=2 SKIP=2 |
| G10 `verify_go.sh` plain | **RED at base** (rc 1, 464 s): only `host/verifygate` `TestMissionConfigDefaultFlowIntegration` (both subtests: `default-flow fixture timed out … context deadline exceeded`, a 10 s fixture deadline at `mission_config_gate_test.go:38`) under load average ~15. **Alone it PASSES** (7.57 s; 3.77 s / 3.80 s). This is pre-existing and unrelated (V-GATE) |

---

## 4. Milestones

LOC are changed lines (added plus removed), estimated against the code as read. The design's
estimate (+95/−20 production, +260 test) is low by about 30%. The holder helper and the picker
account for most of the gap, so the work is split three ways to keep each milestone ≤ ~250.

| Milestone | Files | Est. LOC | Independently green because |
|---|---|---|---|
| M1 | `exec_probe.go`, `exec_probe_test.go` | ~+125/−30 prod, ~+55 test ≈ **210** | all existing probe tests (P1 included) pass through the new nonce path; T1 is pure |
| M2 | `exec_probe_test.go` | ~+215 ≈ **215** | test-only; it uses M1's seam |
| M3 | `exec_probe_test.go`, `ci.yml`, `QUICKSTART.md` | ~+90 test, +1 CI, +5 doc ≈ **100** | the new tests pass; the CI list names only existing tests |

---

### Milestone M1 — production change and the pure judge (T1)

**Edit `host/broker/exec_probe.go`.** Line numbers are at base, from `grep -n` (§10).

**M1.1 — the probe client's `listen` op (`:80` and `:87`).** Make exactly two edits inside the
`execProbeScript` raw string:
- `:80` `function listen(port,ms){return new Promise(r=>{const v=net.createServer(c=>c.destroy());`
  → `function listen(port,ms,nonce){return new Promise(r=>{const v=net.createServer(c=>{c.resume();c.end('WORLD-PROBE-NONCE '+nonce+'\n')});`
- `:87` `else if(op==='listen'){const q=x.split(':');t=await listen(q[0],q[1])}`
  → `else if(op==='listen'){const q=x.split(':');t=await listen(q[0],q[1],q[2])}`

Also update the doc comment at `:59-63`: the listen arg becomes `"<port>:<ms>:<nonce>"`, and the
listener answers each accept with `WORLD-PROBE-NONCE <nonce>`.

**M1.2 — the seam.** In `ExecProbeConfig` (`:99-107`), add after `Client string`:

```go
	// reversePort picks each reverse-leg attempt's port; nil = freeLoopbackPort.
	// Package tests steer it, as they steer Client.
	reversePort func() (int, error)
```

In `type execProber struct` (`:300-305`), add the field `pickPort func() (int, error)`. In
`probeExecProfile` (`:357`), change
`pp := &execProber{h: h, cwd: cwd, client: client}` →
`pp := &execProber{h: h, cwd: cwd, client: client, pickPort: cfg.reversePort}`.
The census shows every `ExecProbeConfig{` and `execProber{` literal is keyed: `exec_probe_test.go:129`,
`exec_startup.go:226`, `exec_probe.go:357`, `exec_probe_test.go:281`. None breaks.

**M1.3 — new constants, types and the pure judge.** Put them right after `execProbeListenMS` (`:96`).
Add `"io"` to the imports (`:26-45`; it is not imported today).

```go
// execProbeReverseAttempts bounds the reverse leg's fresh-port retries.
const execProbeReverseAttempts = 3

// execProbeHello is what the reverse leg's poll writes on every connected
// dial: the child's listener answers it with its nonce; srt's mux answers
// HTTP 400 (design V23); a listener that answers nothing is "silent".
const execProbeHello = "WORLD-PROBE-HELLO\n"

// execProbeHelloWait bounds the read of the answer to the hello.
const execProbeHelloWait = 500 * time.Millisecond

// reverseSeen is what one reverse-leg attempt's dials observed.
type reverseSeen struct{ child, foreign, silent bool }

// revVerdict is judgeReverse's verdict on one attempt.
type revVerdict int

const (
	revPass revVerdict = iota
	revReached
	revNoToken
	revSilent
	revInconclusive
)

// judgeReverse is the reverse leg's whole decision for one attempt, in the
// design's §4.2 order: a nonce always wins; a missing token names itself; a
// silent peer refuses; a foreign answer or the child's EADDRINUSE retries.
func judgeReverse(tok string, seen reverseSeen) revVerdict {
	switch {
	case seen.child:
		return revReached
	case tok != "WORLD-PROBE-BOUND" && !strings.HasPrefix(tok, "WORLD-PROBE-REFUSED "):
		return revNoToken
	case seen.silent:
		return revSilent
	case seen.foreign || tok == "WORLD-PROBE-REFUSED EADDRINUSE":
		return revInconclusive
	}
	return revPass
}
```

**M1.4 — classify one dial, and rewrite `pollLoopback` (`:561-588`).** Replace the whole function
and its doc comment:

```go
// classifyReverseDial writes World's hello on a connected dial, then reads
// at most len(want)+1 bytes up to '\n' within execProbeHelloWait. The line
// want+"\n" is the child; any other bytes, or a clean EOF after a written
// hello, is foreign; zero bytes otherwise (deadline, reset, a failed hello
// write) is silent.
func classifyReverseDial(c net.Conn, want string) (child, foreign bool) {
	_ = c.SetDeadline(time.Now().Add(execProbeHelloWait))
	_, werr := io.WriteString(c, execProbeHello)
	buf := make([]byte, 0, len(want)+1)
	one := make([]byte, 1)
	var rerr error
	for len(buf) < len(want)+1 {
		var n int
		n, rerr = c.Read(one)
		if n == 1 {
			buf = append(buf, one[0])
			if one[0] == '\n' {
				break
			}
		}
		if rerr != nil {
			break
		}
	}
	switch {
	case string(buf) == want+"\n":
		return true, false
	case len(buf) > 0:
		return false, true
	case werr == nil && errors.Is(rerr, io.EOF):
		return false, true
	}
	return false, false // silent
}

// pollLoopback dials 127.0.0.1:port every 50 ms until stopped and
// classifies every connected dial (classifyReverseDial); the returned stop
// function reports what the dials saw. It dials once before it first checks
// for stop, so an attempt always observes at least one dial.
func pollLoopback(port int, want string) func() reverseSeen {
	var child, foreign, silent atomic.Bool
	done := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		if port == 0 {
			return
		}
		addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
		for {
			if c, err := net.DialTimeout("tcp", addr, 200*time.Millisecond); err == nil {
				switch isChild, isForeign := classifyReverseDial(c, want); {
				case isChild:
					child.Store(true)
				case isForeign:
					foreign.Store(true)
				default:
					silent.Store(true)
				}
				_ = c.Close()
			}
			select {
			case <-done:
				return
			case <-time.After(50 * time.Millisecond):
			}
		}
	}()
	return func() reverseSeen {
		close(done)
		<-finished
		return reverseSeen{child: child.Load(), foreign: foreign.Load(), silent: silent.Load()}
	}
}
```

**Load-bearing: keep "dial, then check `done`".** The `/bin/sh` fixtures' client exits in
milliseconds. The existing loop dials *before* its `select`, so every attempt makes at least one dial
even when stop is called at once, and the M2 tests are therefore not timing-dependent. Do not move
the `select` above the dial.

**M1.5 — `reverseLeg`.** Put it right after `probe()`:

```go
// reverseLeg is arm 6's reverse leg (design §4.2): up to
// execProbeReverseAttempts attempts, each on a fresh port with a fresh
// nonce. Attempt 1 is Run C (proxy, then listen), so first carries the
// proxy leg's token at index 0. why is "" when the leg passes.
func (pp *execProber) reverseLeg(ctx context.Context, proxyPort int) (first execProbeRun, why string) {
	pick := pp.pickPort
	if pick == nil {
		pick = freeLoopbackPort
	}
	proxyArgv := []string{pp.client, "-e", execProbeScript, "--", "proxy", strconv.Itoa(proxyPort)}
	var tried []string
	for attempt := 1; attempt <= execProbeReverseAttempts; attempt++ {
		port, err := pick()
		var nonce string
		if err == nil {
			nonce, err = randomHex(16)
		}
		if err != nil {
			if attempt == 1 { // the proxy leg is still judged
				first = pp.run(ctx, proxyArgv, false, execProbeRunBudget)
			}
			return first, fmt.Sprintf("reverse leg: no free port: %v", err)
		}
		argv := []string{pp.client, "-e", execProbeScript, "--"}
		if attempt == 1 {
			argv = slices.Clone(proxyArgv)
		}
		argv = append(argv, "listen", fmt.Sprintf("%d:%d:%s", port, execProbeListenMS, nonce))
		poll := pollLoopback(port, "WORLD-PROBE-NONCE "+nonce)
		r := pp.run(ctx, argv, false, execProbeRunBudget)
		seen := poll()
		tok := r.token(0)
		if attempt == 1 {
			first, tok = r, r.token(1)
		}
		switch judgeReverse(tok, seen) {
		case revPass:
			return first, ""
		case revReached:
			return first, fmt.Sprintf("reverse leg: the host reached the listener the child bound on 127.0.0.1:%d", port)
		case revNoToken:
			return first, fmt.Sprintf("reverse leg: no token from the child's listener (%q; %s)", tok, r.describe())
		case revSilent:
			return first, fmt.Sprintf("reverse leg: a listener on 127.0.0.1:%d accepted the host's dial and answered nothing to World's hello within 500 ms; World cannot tell it from the child's listener, so it refuses", port)
		}
		what := "foreign reach"
		if !seen.foreign {
			what = "child EADDRINUSE"
		}
		tried = append(tried, fmt.Sprintf("127.0.0.1:%d %s", port, what))
	}
	return first, fmt.Sprintf("reverse leg: inconclusive after %d attempts: a listener the child did not bind held each port (%s)", len(tried), strings.Join(tried, ", "))
}
```

**M1.6 — `probe()`: replace base `:452-465` and `:469`, `:481-485`.**
- Replace lines `:452-465` (the Run C comment through `acceptsC := net6.accepts.Load()`) with:

```go
	// Run C, WITHOUT the trampoline so srt's proxy environment reaches the
	// client: arm 6's proxy leg, then the reverse leg (a listener the child
	// binds must be unreachable from the host, V46) while World dials it.
	// reverseLeg retries a port another host process holds (w-exec-probe-
	// reverse-leg-port-race): only the child's nonce counts as a reach.
	c, revWhy := pp.reverseLeg(ctx, net6.port)
	time.Sleep(200 * time.Millisecond)
	acceptsC := net6.accepts.Load()
```

- `:469` `raw, proxy, rev := a.token(6), c.token(0), c.token(1)` → `raw, proxy := a.token(6), c.token(0)`
- Replace `:481-485` (the `if rev != … else if reached …` block) with:

```go
	if revWhy != "" {
		legs = append(legs, revWhy)
	}
```

Keep `freeLoopbackPort` (`:551-559`) unchanged; it is now `reverseLeg`'s default picker. The
`no free port` refusal moves from a separate `pp.fail` into `legs`, so the compaction at `:367` no
longer hides the other legs behind it.

**M1.7 — T1 in `host/broker/exec_probe_test.go`.** Append at the end of the file. Exactly 14 rows:

```go
// T1 (design §7): the reverse leg's whole decision for one attempt, §4.2's order.
func TestExecProbeReverseLegJudge(t *testing.T) {
	const bound, eperm, inuse = "WORLD-PROBE-BOUND", "WORLD-PROBE-REFUSED EPERM", "WORLD-PROBE-REFUSED EADDRINUSE"
	child, foreign, silent := reverseSeen{child: true}, reverseSeen{foreign: true}, reverseSeen{silent: true}
	for _, row := range []struct {
		name string
		tok  string
		seen reverseSeen
		want revVerdict
	}{
		{"bound-none", bound, reverseSeen{}, revPass},
		{"eperm-none", eperm, reverseSeen{}, revPass},
		{"bound-foreign", bound, foreign, revInconclusive},
		{"eperm-foreign", eperm, foreign, revInconclusive},
		{"eaddrinuse-none", inuse, reverseSeen{}, revInconclusive},
		{"bound-child", bound, child, revReached},
		{"bound-child-foreign", bound, reverseSeen{child: true, foreign: true}, revReached},
		{"bound-child-silent", bound, reverseSeen{child: true, silent: true}, revReached},
		{"notoken-child", "", child, revReached},
		{"bound-silent", bound, silent, revSilent},
		{"bound-silent-foreign", bound, reverseSeen{silent: true, foreign: true}, revSilent},
		{"eperm-silent", eperm, silent, revSilent},
		{"notoken-none", "", reverseSeen{}, revNoToken},
		{"notoken-foreign", "", foreign, revNoToken},
	} {
		if got := judgeReverse(row.tok, row.seen); got != row.want {
			t.Errorf("%s: judgeReverse(%q, %+v) = %d, want %d", row.name, row.tok, row.seen, got, row.want)
		}
	}
}
```

**M1 acceptance:**

| Command | Expected |
|---|---|
| G1, G2, G3 | clean / rc 0 / rc 0 |
| `go test ./host/broker -run '^TestExecProbeReverseLegJudge$' -count=1 -v` | `--- PASS: TestExecProbeReverseLegJudge`, FAIL 0 |
| G6 | rc 0. P1, `MissingClient`, `NoopShim`, `ZeroShim`, `FailingToolchain` and `Arm8Wired` all PASS. **P1 must still report `reverse leg: the host reached`**: it now gets there through the nonce path |
| G4, G5, G7 | rc 0; G5 still PASS=2 SKIP=2 |

**M1 mutations** (§7 protocol; each one alone, and the killer must print `--- FAIL`):

| ID | Concrete edit in `exec_probe.go` | Killer: run, and the line that must FAIL |
|---|---|---|
| M1 MUT-REV-FOREIGN-IS-REACH | in `judgeReverse`: `case seen.child:` → `case seen.child \|\| seen.foreign \|\| seen.silent:` | `-run '^TestExecProbeReverseLegJudge$'`: row `bound-foreign` |
| M8 MUT-REV-EADDRINUSE-PASSES | in `judgeReverse`: `case seen.foreign \|\| tok == "WORLD-PROBE-REFUSED EADDRINUSE":` → `case seen.foreign:` | T1 row `eaddrinuse-none` |
| M11 MUT-REV-FOREIGN-DOMINATES | in `judgeReverse`: insert `case seen.foreign:` / `return revInconclusive` as the **first** case | T1 row `bound-child-foreign` |
| M12 MUT-REV-NOTOKEN-RETRIED | in `judgeReverse`: the no-token case's `return revNoToken` → `return revInconclusive` | T1 row `notoken-none` |
| M16 MUT-REV-SILENT-RETRIED | in `judgeReverse`: `return revSilent` → `return revInconclusive` | T1 row `bound-silent` |
| M21 MUT-REV-SILENT-DOMINATES-CHILD | in `judgeReverse`: move `case seen.silent:` / `return revSilent` above `case seen.child:` | T1 row `bound-child-silent` |
| M9 MUT-REV-NEVER-REAL (P1 half) | in `classifyReverseDial`: `case string(buf) == want+"\n":` → `case string(buf) == want+"X\n":` | `-run '^TestExecProbePassThroughShimFailsEveryConfinementArm$'`: `does not report "reverse leg: the host reached"` |

---

### Milestone M2 — the `/bin/sh` fixture tests (T2, T3, T5, T7). They need no node and no srt

All of M2 goes in `host/broker/exec_probe_test.go`. Add imports `"bufio"`, `"fmt"`, `"io"`, `"net"`,
`"sync"` and `"sync/atomic"` (the file imports none of them today; see its `:12-23`).

**M2.1 — helpers (≈ 95 LOC).** Append these verbatim:

```go
// revHolderKinds are the test-owned holders of a reverse-leg port (design
// §6): chatty ones drain World's hello first (closing with unread input
// would reset) and then answer; "silent" accepts and never writes;
// "reset" accepts and closes with SO_LINGER 0.
var revHolderKinds = []string{"http-400", "forged-nonce", "garbage", "clean-close", "silent", "reset"}

func revHolder(t *testing.T, kind string) int {
	t.Helper()
	if !slices.Contains(revHolderKinds, kind) {
		t.Fatalf("revHolder: unknown kind %q", kind)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var held []net.Conn
	t.Cleanup(func() {
		_ = ln.Close()
		mu.Lock()
		defer mu.Unlock()
		for _, c := range held {
			_ = c.Close()
		}
	})
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			switch kind {
			case "silent":
				mu.Lock()
				held = append(held, c)
				mu.Unlock()
				continue
			case "reset":
				_ = c.(*net.TCPConn).SetLinger(0)
				_ = c.Close()
				continue
			}
			go func() {
				defer c.Close()
				_ = c.SetDeadline(time.Now().Add(2 * time.Second))
				_, _ = bufio.NewReader(c).ReadString('\n') // World's hello
				switch kind {
				case "http-400":
					_, _ = io.WriteString(c, "HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n")
				case "forged-nonce":
					_, _ = io.WriteString(c, "WORLD-PROBE-NONCE "+strings.Repeat("0", 32)+"\n")
				case "garbage":
					_, _ = io.WriteString(c, "garbage\n")
				} // clean-close: zero bytes, then FIN
			}()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

// revPicker is the reverse leg's port seam in tests: call n returns
// steps[n-1] ("free" = a port nothing holds, else a revHolder kind) and
// records the port; a call past the list fails, so an unbounded retry
// ends at once instead of hanging.
type revPicker struct {
	t     *testing.T
	steps []string
	calls atomic.Int32
	mu    sync.Mutex
	ports []int
}

func newRevPicker(t *testing.T, steps ...string) *revPicker { return &revPicker{t: t, steps: steps} }

func (p *revPicker) pick() (int, error) {
	n := int(p.calls.Add(1))
	if n > len(p.steps) {
		return 0, fmt.Errorf("test picker exhausted after %d ports", len(p.steps))
	}
	var port int
	if p.steps[n-1] == "free" {
		var err error
		if port, err = freeLoopbackPort(); err != nil {
			return 0, err
		}
	} else {
		port = revHolder(p.t, p.steps[n-1])
	}
	p.mu.Lock()
	p.ports = append(p.ports, port)
	p.mu.Unlock()
	return port, nil
}

func (p *revPicker) port(t *testing.T, i int) int {
	t.Helper()
	p.mu.Lock()
	defer p.mu.Unlock()
	if i >= len(p.ports) {
		t.Fatalf("picker handed out %d ports, want at least %d", len(p.ports), i+1)
	}
	return p.ports[i]
}

// revFixture is arm8Fixture's /bin/sh node, srt and client (V24: no node,
// no real srt): the client answers the proxy op with a refusal and the
// listen op with listenTok ("" = no token), and never binds.
func revFixture(t *testing.T, listenTok string, p *revPicker) *execProber {
	t.Helper()
	body := "for a in \"$@\"; do [ \"$a\" = proxy ] && echo 'WORLD-PROBE-REFUSED ECONNREFUSED'; done\n"
	if listenTok != "" {
		body += "echo '" + listenTok + "'\n"
	}
	pp := arm8Fixture(t, body, nil)
	pp.pickPort = p.pick
	return pp
}

func runReverseLeg(t *testing.T, pp *execProber) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, why := pp.reverseLeg(ctx, 1) // port 1: nothing listens there
	return why
}

const revReachedText = "the host reached the listener the child bound"
```

**M2.2 — T2.** Rows {`linux-bound` = `WORLD-PROBE-BOUND`, `darwin-eperm` = `WORLD-PROBE-REFUSED
EPERM`} × {`http-400`, `forged-nonce`, `garbage`, `clean-close`}, 8 subtests, each named
`t.Run(plat.name+"-"+kind, …)` (a hyphen, never `/`). For each subtest:
`p := newRevPicker(t, kind, "free", "free")`; `why := runReverseLeg(t, revFixture(t, tok, p))`.
Require `why == ""` **and** `p.calls.Load() == 2`; on failure print both.

**M2.3 — T3.** Two rows:
- `chatty-every-time`: tok `WORLD-PROBE-BOUND`, steps `"http-400","http-400","http-400"`.
- `child-eaddrinuse`: tok `WORLD-PROBE-REFUSED EADDRINUSE`, steps `"free","free","free"`.

Each row requires `strings.HasPrefix(why, "reverse leg: inconclusive after 3 attempts")`,
`!strings.Contains(why, revReachedText)` and `p.calls.Load() == 3`.

**M2.4 — T5.** Two rows:
- `silent-client`: tok `""`, steps `"free","free","free"`. Requires
  `strings.Contains(why, "reverse leg: no token from the child's listener")` and calls `== 1`.
- `clean`: tok `WORLD-PROBE-BOUND`, steps `"free","free","free"`. Requires `why == ""` and calls `== 1`.

**M2.5 — T7.** All three rows use tok `WORLD-PROBE-BOUND`. Define
`silentText := func(port int) string { return fmt.Sprintf("reverse leg: a listener on 127.0.0.1:%d accepted the host's dial and answered nothing", port) }`.
- `silent-attempt-1`: steps `"silent","free","free"`. Requires
  `strings.HasPrefix(why, silentText(p.port(t,0)))`, calls `== 1`.
- `silent-after-chatty`: steps `"http-400","silent","free"`. Requires
  `strings.HasPrefix(why, silentText(p.port(t,1)))`, calls `== 2`.
- `reset-holder`: steps `"reset","free","free"`. Requires
  `strings.HasPrefix(why, silentText(p.port(t,0)))`, calls `== 1`.

Every row also requires `strings.Contains(why, "answered nothing to World's hello")` and
`!strings.Contains(why, revReachedText)`.

**M2 acceptance:**

| Command | Expected |
|---|---|
| G1, G2, G3 | clean / rc 0 / rc 0 |
| `go test ./host/broker -run '^TestExecProbeReverseLeg(Judge\|ChattyForeignHolderIsRetried\|ExhaustionRefusesDistinctly\|NoTokenIsNotRetried\|SilentHolderRefusesSticky)$' -count=1 -v` | 5 top-level `--- PASS`, 0 FAIL, 0 SKIP, 15 subtest PASS lines (T2 8, T3 2, T5 2, T7 3) |
| G8 (node absent) | rc 0, 5 PASS, **0 SKIP** |
| G6, G4, G5, G7 | as in M1. G7 now covers the holders and the picker under `-race` |
| `go test ./host/broker -run '^TestExecProbeReverseLeg' -count=10` | rc 0 (stability spot check before M3's ×20) |

**M2 mutations.** The killer regex is `-run '^TestExecProbeReverseLeg'` unless stated.

| ID | Concrete edit in `exec_probe.go` | Must FAIL |
|---|---|---|
| M1 | as in M1 | T2: every row (Why = `reverse leg: the host reached …`). The named killer is T2 |
| M2 MUT-REV-NO-RETRY | in `reverseLeg`, insert `break` on the line after `tried = append(tried, …)` | T2: Why `reverse leg: inconclusive after 1 attempts…`, calls 1 |
| M3 MUT-REV-ANY-BYTES | in `classifyReverseDial`: `case string(buf) == want+"\n":` → `case len(buf) > 0:` | T2 rows `http-400`, `forged-nonce`, `garbage` (refused as reached) |
| M4 MUT-REV-PREFIX-NONCE | same line → `case strings.HasPrefix(string(buf), "WORLD-PROBE-NONCE "):` | T2 rows `*-forged-nonce` |
| M5 MUT-REV-EXHAUST-PASSES | the final `return first, fmt.Sprintf("reverse leg: inconclusive after …` → `return first, ""` | T3 both rows |
| M6 MUT-REV-UNBOUNDED | `for attempt := 1; attempt <= execProbeReverseAttempts; attempt++ {` → `for attempt := 1; ; attempt++ {` (vet may flag the then-unreachable return; it still compiles) | T3: calls 4, Why `reverse leg: no free port: test picker exhausted…` |
| M7 MUT-REV-EXHAUST-AS-REACH | the final format string → `"reverse leg: inconclusive after %d attempts: the host reached the listener the child bound (%s)"` | T3: the `!Contains(revReachedText)` assertion |
| M8 | as in M1 | T3 row `child-eaddrinuse` (passes, calls 1) |
| M12 | as in M1 | T5 row `silent-client`: calls 3, and Why lacks `no token` |
| M13 MUT-REV-ALWAYS-RETRY | in `reverseLeg`: `case revPass:` / `return first, ""` → `case revPass:` / `continue` | T5 row `clean`: calls 3 |
| M15 MUT-REV-HELLO-DROPPED | in `classifyReverseDial`: delete the line `_, werr := io.WriteString(c, execProbeHello)` and replace it with `var werr error` | T2 rows `*-http-400` (and the other chatty rows) refuse SILENT |
| M16 | as in M1 | T7 row `silent-attempt-1`: passes on attempt 2, calls 2 |
| M17 MUT-REV-SILENT-NOT-STICKY | in `reverseLeg`: `case revSilent:` / `return first, …` → `case revSilent:` / `continue` | T7 row `silent-after-chatty`: passes on attempt 3 |
| M18 MUT-REV-EOF-IS-SILENT | in `classifyReverseDial`: delete `case werr == nil && errors.Is(rerr, io.EOF):` and its `return false, true` | T2 rows `*-clean-close` refuse SILENT |
| M19 MUT-REV-RESET-IS-FOREIGN | in `classifyReverseDial`: `case werr == nil && errors.Is(rerr, io.EOF):` → `case !errors.Is(rerr, os.ErrDeadlineExceeded):` | T7 row `reset-holder`: passes on attempt 2 |
| M20 MUT-REV-SILENT-AS-REACH-TEXT | the revSilent format string → `"reverse leg: the host reached the listener the child bound on 127.0.0.1:%d"` | T7: every row |

---

### Milestone M3 — node tests (T4, T6), the CI list and the QUICKSTART prose

**M3.1 — T4 in `exec_probe_test.go`.** It needs node; `newProbeFixture` calls `probeNode`, which
skips by name when node is absent.

```go
// T4 (design §7): a collision on attempt 1 must not mask a real reach on
// attempt 2. The pass-through shim runs World's real node client
// unsandboxed: attempt 1's child gets EADDRINUSE under an http-400 holder
// (INCONCLUSIVE); attempt 2's child binds and answers the nonce.
func TestExecProbeReverseLegRealReachSurvivesAForeignAttempt(t *testing.T) {
	f := newProbeFixture(t, shimSrtTree(t, "passthrough"), nil, "", []any{"sh", "-c", "exit 0"})
	pp := f.arm8Prober(t) // client = World's node
	p := newRevPicker(t, "http-400", "free", "free")
	pp.pickPort = p.pick
	why := runReverseLeg(t, pp)
	want := fmt.Sprintf("reverse leg: %s on 127.0.0.1:%d", revReachedText, p.port(t, 1))
	if why != want || p.calls.Load() != 2 {
		t.Fatalf("Why = %q, picker calls %d; want %q and 2", why, p.calls.Load(), want)
	}
}
```

Planner measurement (§10 P7): node `listen` on a port a Go listener holds gets `REFUSED EADDRINUSE`
on darwin, 2 of 2 times.

**M3.2 — T6 in `exec_probe_test.go`.** It goes through the production entry point. The stand-in
below was run verbatim by the planner at base (§10 P5):

```go
// revStandIn is a node srt stand-in for T6: a run whose argv carries
// `listen` never runs the child; it prints the Linux-correct tokens, and
// when the argv also carries `proxy` (attempt 1 = Run C, as srt's own mux
// would) it first plants a holder of shape MODE on the listen port.
// Every other run passes through, unsandboxed.
const revStandIn = `const {spawnSync}=require('child_process');const os=require('os');const net=require('net');
const a=process.argv.slice(2);
if(a[0]!=='--settings'||a[2]!=='--'){console.error('STANDIN: bad launch '+JSON.stringify(a));process.exit(99)}
const argv=a.slice(3);const li=argv.indexOf('listen');
if(li<0){const r=spawnSync(argv[0],argv.slice(1),{stdio:'inherit'});if(r.error){console.error('STANDIN: '+r.error.code);process.exit(127)}
process.exit(r.status===null?128+os.constants.signals[r.signal]:r.status)}
const mode=MODE;const q=argv[li+1].split(':');const port=+q[0],ms=+q[1];
const say=s=>process.stdout.write(s+'\n');
const report=()=>{if(argv.includes('proxy'))say('WORLD-PROBE-REFUSED ECONNREFUSED');say('WORLD-PROBE-BOUND')};
const done=srv=>setTimeout(()=>{if(srv)srv.close();process.exit(0)},ms);
if(mode==='control'||!argv.includes('proxy')){report();done(null)}
else{const srv=net.createServer(mode==='chatty'?c=>{c.once('data',()=>c.end('HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n'))}:c=>{});
srv.on('error',e=>{console.error('STANDIN: holder '+e.code);process.exit(98)});
srv.listen(port,'127.0.0.1',()=>{report();done(srv)})}
`

// T6 (design §7): the CI flake, reproduced through ProbeExecProfiles.
// chatty is RED at base (P0); silent is the sticky refusal; control has no holder.
func TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse(t *testing.T) {
	for _, mode := range []string{"chatty", "silent", "control"} {
		t.Run(mode, func(t *testing.T) {
			cli := strings.Replace(revStandIn, "MODE", strconv.Quote(mode), 1)
			f := newProbeFixture(t, fakeSrtTree(t, "0.0.0-test", cli), nil, "", []any{"sh", "-c", "echo toolchain-ok"})
			err := f.probe(t, "")
			var why string
			var perr *ExecProbeError
			if errors.As(err, &perr) {
				for _, fl := range perr.Failed {
					if fl.Arm == ExecProbeArmNetwork {
						why = fl.Why
					}
				}
			}
			t.Logf("%s: arm 6 = %q", mode, why)
			switch mode {
			case "chatty", "control":
				if strings.Contains(why, "reverse leg") {
					t.Errorf("%s: arm 6 reports a reverse leg: %s", mode, why)
				}
			case "silent":
				if !strings.Contains(why, "reverse leg: a listener on 127.0.0.1:") || !strings.Contains(why, "answered nothing to World's hello") || strings.Contains(why, revReachedText) {
					t.Errorf("silent: arm 6 must carry the silent refusal, not the reached text: %s", why)
				}
			}
			f.assertProbeCleaned(t)
		})
	}
}
```

Add `"strconv"` to the test file's imports.

**Why the stand-in fakes *every* `listen` run, and plants the holder only with `proxy`.** If attempt
2 (`listen` without `proxy`) were passed through, the unsandboxed child would bind host loopback and
answer its nonce. The `chatty` row would then refuse as REACHED after the fix (design discrepancy
D1, §9). The other runs (A, B, D) pass through, so arm 6 still fails its raw and proxy legs. That is
expected: T6 asserts only on the *reverse-leg* text.

**M3.3 — `.github/workflows/ci.yml` step 8.** Insert **one new line** directly after line `:247`
(`tests="$tests TestExecStartupRefusalTable … TestExecSrtWorktreeDelta"`), with the same 10-space
indentation:

```
          tests="$tests TestExecProbeReverseLegJudge TestExecProbeReverseLegChattyForeignHolderIsRetried TestExecProbeReverseLegExhaustionRefusesDistinctly TestExecProbeReverseLegRealReachSurvivesAForeignAttempt TestExecProbeReverseLegNoTokenIsNotRetried TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse TestExecProbeReverseLegSilentHolderRefusesSticky"
```

Change nothing else in `ci.yml`. The PASS-guard loop (`:254-256`) then covers all seven.

**M3.4 — `docs/QUICKSTART.md` §10, arm 6 (`:1146-1148`).** Append these lines directly after
`:1148` (`   the child binds must be unreachable from the host.`), with 3-space continuation indent so
they stay inside list item 6:

```
   A reach counts only when the listener answers World's hello with the child's one-time nonce. A
   port another host process holds is retried on a fresh port, up to 3 attempts. `reverse leg:
   inconclusive after 3 attempts` means something on the host kept taking the chosen ports: retry,
   or find that process. `reverse leg: a listener on … answered nothing to World's hello` means a
   process that never answers held the port, and World cannot tell it from the child: find it.
```

Do not touch the refusal table (`<!-- exec-refusals:begin -->`). `TestExecStartupRefusalTable` (in
G5) and `TestQuickstartSection10ProfilesLoad` (in G6) parse §10. No new line may start with
`` | ` `` or contain `-d '`, the shapes their regexes match.

**M3 acceptance:**

| Command | Expected |
|---|---|
| G1, G2, G3 | clean / rc 0 / rc 0 |
| `go test ./host/broker -run '^(TestExecProbeReverseLegRealReachSurvivesAForeignAttempt\|TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse)$' -count=1 -v` | 2 top-level PASS, 3 T6 subtest PASS, **0 SKIP** (node is present locally) |
| G6 (now including `TestQuickstartSection10ProfilesLoad`), G4, G5, G7 | rc 0; G5 PASS=2 SKIP=2 |
| G8 | unchanged: 5 PASS, 0 SKIP |
| G9 AC-STABLE ×20 | each of the 7 names PASS 20/20, 0 FAIL |
| G10 `verify_go.sh` plain | rc 0 |
| G11 actionlint and YAML parse | rc 0 / `ok` |
| G12 step-8 emulation (below) | all 7 new names `--- PASS:` |
| G13 | `git diff --name-only` = exactly the 4 scope files |

G12, the step-8 emulation. Run it from the worktree, in bash, because `${tests// /|}` is a bashism
that zsh does not handle the same way:

```bash
bash -c 'set -uo pipefail
sed -n 245,248p .github/workflows/ci.yml     # eyeball: 4 tests= lines, the 4th is the new one
tests="$(sed -n 245,248p .github/workflows/ci.yml | grep -oE "Test[A-Za-z0-9_]+" | tr "\n" " " | sed "s/ $//")"
echo "names=$(echo $tests | wc -w)"           # expect 28 (21 existing + 7 new)
log=/tmp/w314-step8.log
go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -timeout 600s -v -run "^(${tests// /|})\$" >"$log" 2>&1; echo "go rc=$?"
for t in TestExecProbeReverseLegJudge TestExecProbeReverseLegChattyForeignHolderIsRetried TestExecProbeReverseLegExhaustionRefusesDistinctly TestExecProbeReverseLegRealReachSurvivesAForeignAttempt TestExecProbeReverseLegNoTokenIsNotRetried TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse TestExecProbeReverseLegSilentHolderRefusesSticky; do
  grep -q -- "^--- PASS: $t " "$log" && echo "PASS $t" || echo "MISSING $t"; done'
```

Expect `names=28`, `go rc=0` and 7 `PASS` lines. The real-srt names SKIP locally, as expected (§5).
If `names` is not 28, read lines 245–248 and build the list by hand. Never weaken the check.

M2 and M3 focused runs, to copy (the tables escape `|`):

```sh
# M2
go test ./host/broker -run '^TestExecProbeReverseLeg(Judge|ChattyForeignHolderIsRetried|ExhaustionRefusesDistinctly|NoTokenIsNotRetried|SilentHolderRefusesSticky)$' -count=1 -v > /tmp/w314-m2.txt 2>&1; echo "rc=$?"
grep -E '^(    )*--- (PASS|FAIL|SKIP)' /tmp/w314-m2.txt
# M3
go test ./host/broker -run '^(TestExecProbeReverseLegRealReachSurvivesAForeignAttempt|TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse)$' -count=1 -v > /tmp/w314-m3.txt 2>&1; echo "rc=$?"
grep -E '^(    )*--- (PASS|FAIL|SKIP)' /tmp/w314-m3.txt
```

**M3 mutations**, killer regex `-run '^TestExecProbe(ReverseLegRealReach|ForeignListener)'`:

| ID | Concrete edit | Must FAIL |
|---|---|---|
| M1 | as in M1 | T6 `chatty` (`arm 6 reports a reverse leg`) |
| M9 | as in M1 | T4: Why becomes `reverse leg: inconclusive after 3 attempts…`, not the reached text on port 2 (and P1, already shown in M1) |
| M10 MUT-REV-SAME-PORT | in `reverseLeg`: add `var port int` on the line before the `for`; inside the loop, replace `port, err := pick()` with the two statements `var err error` and `if attempt == 1 { port, err = pick() }` | T4: attempts 2 and 3 reuse the held port, so Why is the inconclusive text and calls are 1 |
| M14 MUT-REV-SEAM-ONLY | in `probe()`: replace `c, revWhy := pp.reverseLeg(ctx, net6.port)` with the legacy block below | T6 `chatty`. T1–T5 and T7 stay green, which proves T6 is the only test that crosses `probe()` |
| M15 | as in M2 | T6 `chatty` (it refuses SILENT, so `reverse leg` is present) |
| M20 | as in M2 | T6 `silent` (reached text present, silent text absent) |

M14's legacy block (paste it in place of the single `reverseLeg` line):

```go
	revPort, _ := freeLoopbackPort()
	revPoll := pollLoopback(revPort, "WORLD-PROBE-NONCE x")
	c := pp.run(ctx, []string{pp.client, "-e", execProbeScript, "--", "proxy", strconv.Itoa(net6.port),
		"listen", fmt.Sprintf("%d:%d:x", revPort, execProbeListenMS)}, false, execProbeRunBudget)
	revWhy := ""
	if s := revPoll(); s.child || s.foreign || s.silent {
		revWhy = fmt.Sprintf("reverse leg: the host reached the listener the child bound on 127.0.0.1:%d", revPort)
	}
```

**M3 closing pass, AC-MUT in full.** Apply each of M1–M21 alone and run
`go test ./host/broker -run '^TestExecProbe' -count=1 -v`. Record the **full** set of top-level
tests that print `--- FAIL` (the blast radius). The criterion is: every named killer from the design
§6 table is **in** that set, and the unmutated control run has an empty set. Several mutants
necessarily red more than their named killers (D2, §9). Record those extras; they do not count as
failures of the criterion.

---

## 5. Local versus CI-only evidence

- **Runs locally and in CI:** T1–T7 and P1. T1, T2, T3, T5 and T7 also run node-free (G8). T4 and
  T6 need node, which is at `/opt/homebrew/bin/node` locally and is set as `WORLD_EXEC_NODE` on the
  runner (`ci.yml:233`).
- **CI-only evidence.** These SKIP locally, because `WORLD_EXEC_SRT_NODE_MODULES` is unset:
  - `TestExecSrtStartupProbePasses` and `TestExecSrtProbeRefusesAPackageOnlySandbox` (AC-REAL: the
    Linux correct path, BOUND with no reach, passes on 1 attempt);
  - `TestExecSrtMCPGrantAndBudget` (the CI victim) and `TestExecSrtMCPEndToEnd`;
  - the `TestExecSrtReplayIsByteEqual` subtests;
  - `TestExecStartupRefusalTable/probe-package-only`;
  - `TestExecProbeArm8`;
  - the `handlers_exec_srt_test.go` `TestExecSrt*` set.

  Locally, G5 shows the two daemon e2e SKIPs and G6 shows `TestExecProbeArm8` SKIP.
- **The CI PASS-guard** is `ci.yml:254-256`: `for t in $tests; do grep -q -- "^--- PASS: $t " "$log"
  || { echo "::error::$t did not PASS on linux (SKIP?)"; exit 1; }`. M3.3 puts all seven new names
  under it, so a SKIP of any of them on Linux fails CI (design AC-CI, S6).
- **Not seen by any test here:** srt's *real* mux inside a probe on Linux (design §7, "What this
  cannot see"; R4 and the V23 scope). If CI later shows `reverse leg: inconclusive after 3 attempts`
  or `… answered nothing to World's hello`, re-measure those assumptions; do not widen the fix.
- A gap noticed, outside scope: the guard checks only *top-level* PASS lines, and
  `TestExecSrtReplayIsByteEqual`'s parent prints `--- PASS` even when both of its subtests SKIP
  (G5 at base). This item does not change that.

---

## 6. Risks

| Risk | Mitigation |
|---|---|
| A chatty holder closes with unread input, and the host sees a reset instead of bytes | `revHolder` drains the hello line before it answers. Node's listener `c.resume()`s (M1.1) |
| A reset arrives before the hello write; Linux then reports a *clean EOF* on the next read | `classifyReverseDial` treats zero bytes after a **failed** hello write as silent (D3). Without that, T7 `reset-holder` could flake on Linux |
| Moving the poll's `select` above the dial makes the `/bin/sh` tests timing-dependent | §M1.4's load-bearing note. G9 ×20 would catch it |
| M6's unbounded loop hangs | `revPicker` fails past its list, so M6 reds at once with calls = 4 |
| A data race in the holders or the picker | atomics plus a mutex; G7 runs `-race` |
| The QUICKSTART edit breaks a §10 parser | prose-only, inside list item 6. G5 and G6 run both parsers |

## 7. Mutation protocol (every mutant, every milestone)

```sh
S=/tmp/w314-mut; mkdir -p $S
cp host/broker/exec_probe.go $S/exec_probe.go.snap
shasum -a 256 host/broker/exec_probe.go host/broker/exec_probe_test.go > $S/pre.sha
#   apply ONE edit from the milestone's table with the Edit tool
go test ./host/broker -run '<killer regex>' -count=1 -v > $S/<ID>.txt 2>&1; echo "rc=$?"
grep -E '^(    )*--- FAIL' $S/<ID>.txt          # must list the named killer
cp $S/exec_probe.go.snap host/broker/exec_probe.go
shasum -a 256 -c $S/pre.sha                    # both files: OK
```

Every mutant edits `exec_probe.go` only, and the test file is never mutated. Before the first mutant
of each milestone, run the killer regex once unmutated and record 0 FAIL: that is the control. A
mutant that fails to **compile** is a protocol error. Fix the mutant's form (keep its meaning) and
record what you changed; a compile error is never a kill.

## 8. Definition of done

- G1–G13 green after M3, with readings recorded against the §3 base. The one allowed exception is
  G10's pre-existing `host/verifygate` load timeout, handled per G10's rule and recorded as red at
  base.
- P0's base reading recorded: chatty RED, silent RED, control GREEN.
- AC-MUT: every design-named killer is in its mutant's red set, the control set is empty, and both
  Go files are sha256-restored after each mutant.
- `git diff --name-only` lists exactly `host/broker/exec_probe.go`, `host/broker/exec_probe_test.go`,
  `docs/QUICKSTART.md` and `.github/workflows/ci.yml`.
- No git writes. Report back to the controller, who commits.

## 9. Design discrepancies found by the planner (resolved in-plan, flagged to the controller)

- **D1, T6's stand-in.** Design §7 T6 (lines 368–371): "An srt stand-in binds a holder on the
  `listen` port only when the argv also carries `proxy` … It prints the Linux-correct tokens without
  running the child, and passes every other run through." Read literally, attempts 2 and 3 (`listen`
  without `proxy`) would pass through. The unsandboxed child would then answer its nonce, and the
  `chatty` row would refuse as REACHED **after** the fix. **Resolved:** the stand-in fakes every run
  carrying `listen` and plants the holder only with `proxy`. This is V11's original shape.
- **D2, AC-MUT "exactly".** Design §7 (lines 409–410): "Each of M1–M21 is applied alone and reds
  exactly its named killer." That cannot hold. M1 (foreign is a reach) also reds T3
  `chatty-every-time`, T4 and T7 `silent-after-chatty`. M15 (no hello) also reds T3, T4 and T7.
  **Resolved:** the named killer must be *in* the red set, and the full set is recorded.
- **D3, a reset on the hello write.** Design §4.2 (lines 186–188): "≥1 byte that is not that line,
  or a clean EOF with zero bytes: set `seen.foreign`; zero bytes at the deadline, or a reset with
  zero bytes: set `seen.silent`." Linux TCP semantics (read, **not measured**) are that a reset which
  lands before the hello write is reported on the *write*, and the next read then returns EOF.
  **Resolved:** "clean EOF" requires a successful hello write, so a failed write with zero bytes is
  silent. This is still the design's "reset with zero bytes" branch.
- **D4, a partial nonce (not adopted, for the controller).** Under §4.2, a strict prefix of the
  nonce line that ends at the deadline is "≥1 byte that is not that line", so it is *foreign* and
  retryable. §4.3 says World's client "can only ever produce child or silent". Over loopback the
  nonce is one write under 60 bytes (V13 and V26: 1000/1000 whole), so the gap is theoretical. The
  plan implements §4.2 literally; tightening it would need design review.
- **D5, an unspecified case.** The design does not say what happens when the picker (or the nonce)
  fails on attempt 2 or 3. **Resolved:** it refuses with `reverse leg: no free port: %v`, which
  fails closed.
- **D6, the exhaustion count.** M2's description (`inconclusive after 1 …`) implies the count is
  the number of attempts made. **Resolved:** `%d` = `len(tried)`.
- **D7, line cites (cosmetic).** AC-CI cites `ci.yml:254-257`. The guard loop is `:254-256` (`:257`
  is `done`), and V5 says `:254-256`. §1's `exec_probe.go:466-488` cite is accurate.
- **D8, sizing.** The design's +95/−20 production and +260 test lines are low. The plan's estimate
  is about +125/−30 and +360, hence three milestones.

## 10. Verification notes (planner, this session, worktree `e4dbac36`)

| # | Claim | Command | Observed |
|---|---|---|---|
| P1 | The base code sites | `grep -n 'pollLoopback\|freeLoopbackPort\|createServer\|Client  *string\|reverse leg\|^type execProber\|^func ' host/broker/exec_probe.go` | `:80` `net.createServer(c=>c.destroy())`; `:96` `execProbeListenMS = 1000`; `:106` `Client string`; `:300` `type execProber struct`; `:390` `func (pp *execProber) probe`; `:455` `revPort, err := freeLoopbackPort()`; `:460` `poll := pollLoopback(revPort)`; `:482`/`:484` the two reverse-leg texts; `:552` `func freeLoopbackPort`; `:563` `func pollLoopback(port int) func() bool`. Imports `:26-45`: no `io`. Run C block `:452-488` read with `sed -n 449,489p` |
| P2 | The CI step 8 shape | `grep -n 'tests=\|go test ./host/broker\|grep -q -- "^--- PASS: \$t\|did not PASS on linux' .github/workflows/ci.yml`; `sed -n 198,290p` | `:245-247` `tests=`; `:250` `go test ./host/broker/ ./host/daemon/ -count=1 -p 1 -timeout 600s -v`; `:256` `::error::$t did not PASS on linux (SKIP?)`; `:233` `WORLD_EXEC_NODE=$(command -v node)` into `GITHUB_ENV`; `:274-280` the go gate `VERIFY_GO_LEGS: plain` `./scripts/verify_go.sh`; `:418`/`:496` the race job; `:113-114` `verify_ail.sh`. `grep -n 'gofmt\|go vet\|staticcheck\|golangci' .github/workflows/*.yml scripts/*.sh` returned 0 lines |
| P3 | The base gates | G1, G2, G3, G4, G5 (with `AILANG_BIN`), G6 as in §2 | as in §2's base column. G4: `ok … 167.649s`. G5: PASS=2 SKIP=2. G6: PASS=10 SKIP=1 |
| P4 | The new names are absent at base | `go test ./host/broker -run '^TestExecProbeReverseLeg\|^TestExecProbeForeignListener' -count=1 -v` | rc 0, `testing: warning: no tests to run` |
| P5 | **T6 at base** (a scratch `host/broker/zz_planner_iter314_test.go` holding §M3.2's stand-in verbatim, then deleted) | `perl -e 'alarm 300; exec @ARGV' go test ./host/broker -run '^TestZZPlannerIter314T6$' -count=3 -v` | rc 1. `chatty` FAIL ×3 (`… reverse leg: the host reached the listener the child bound on 127.0.0.1:55710` / `:55786` / `:55861`); `silent` FAIL ×3 (the same reached text); `control` PASS ×3, with no reverse-leg text. `assertProbeCleaned` clean |
| P6 | The `/bin/sh` fixture drives a no-trampoline `proxy`+`listen` run with **node absent** | the same scratch file: `arm8Fixture` + `pp.run(ctx, [client,-e,script,--,proxy,1,listen,4242:1000:abcd], false, …)`, under `env -i HOME=$HOME PATH=/usr/bin:/bin WORLD_EXEC_SRT_NODE_MODULES=/nonexistent … /opt/homebrew/bin/go test …` | rc 0. tokens `["WORLD-PROBE-REFUSED ECONNREFUSED" "WORLD-PROBE-BOUND"]`, `$target` = `4242:1000:abcd`. A listen-only run gave `["WORLD-PROBE-BOUND"]`. The scratch file was then deleted: `git status --porcelain` showed 0 lines, and `shasum` of both Go files equalled the base values |
| P7 | node `listen` on a Go-held port | a throwaway Go program in `/tmp/w314p_eaddr` (removed): `net.Listen("tcp","127.0.0.1:0")`, then node `v.listen(port,'127.0.0.1')` | `REFUSED EADDRINUSE` 2/2 (darwin) |
| P8 | The literal census for the seam | `grep -rn 'ExecProbeConfig{\|execProber{' --include='*.go' .` | 4 literals, all keyed: `exec_probe_test.go:129`, `exec_startup.go:226`, `exec_probe.go:357`, `exec_probe_test.go:281` |
| P9 | QUICKSTART arm 6 prose and its parsers | `sed -n 1140,1160p docs/QUICKSTART.md`; `grep -n 'find(\`\|regexp.MustCompile' host/broker/exec_quickstart_test.go host/daemon/exec_startup_test.go`; `grep -rn 'the child binds must be unreachable' --include='*.md' .` | `:1146-1148` is the arm 6 item. Regexes: the table row `^\| \`…`, `-d '…'`, `At most …`, `deadline − … ms\)`, `On the (v…) tool pin`. The arm 6 sentence exists only in `docs/QUICKSTART.md:1148` (outside `design_docs`) |
| P10 | The CI file lints clean at base | `actionlint .github/workflows/ci.yml`; `ruby -ryaml -e 'YAML.load_file(…)'` | rc 0, 0 findings; `ruby-yaml ok` (PyYAML is not installed) |
| P11 | The test helpers' locations | `grep -rn 'execNodeEnv *=\|func fakeSrtTree\|func fakePin\|func canonicalTempDir\|func mustMkdir\|func (h \*ExecHandler) launch' host/broker/*.go`; `grep -n 'func execHostEnv' host/broker/*.go` | `execNodeEnv = "WORLD_EXEC_NODE"` at `handlers_exec_srt_test.go:38`; `fakeSrtTree` at `exec_sandbox_test.go:18`; `launch` at `handlers_exec.go:379` (`[cli, --settings, S, --, (trampoline…), argv…]`); `execHostEnv` = `PATH=/usr/bin:/bin` only, so the pass-through child has no proxy env |
| V-GATE | `verify_go.sh` plain at base | `VERIFY_GO_LEGS=plain perl -e 'alarm 1500; exec @ARGV' ./scripts/verify_go.sh` (with `AILANG_BIN` set); then `go test ./host/verifygate -run '^TestMissionConfigDefaultFlowIntegration$' -count=1 -v` alone | the gate: **rc 1, wall 464 s**. Every package `ok` except `FAIL host/verifygate 428.169s`, from `TestMissionConfigDefaultFlowIntegration/{valid_ledger,invalid_ledger}`: `mission_config_gate_test.go:286: default-flow fixture timed out / instrument failure: context deadline exceeded` (10.16 s / 10.10 s). The re-run alone: rc 0, PASS 7.57 s (3.77 s / 3.80 s); `uptime` load averages 15.36 15.22 9.85. `host/broker` itself was `ok 294.430s` inside the gate |

# Sprint plan — `w-human-approve-control-window-flake` (iteration 334, queue row 235)

**Authority:** [design r3](w-human-approve-control-window-flake.md), especially §2.1–§2.5, §4, §5 and §10; [coding standards](../coding-standards.md). **Planner:** mission iteration 334. **Sprint HEAD:** `cf79a7a92199eb65ddda4bf724167ba9c8c092e1`; its Go tree equals measurement base `388001c2` (all intervening paths are design/evidence). **Estimate:** 0.1 day of code change; execution dominated by soak/package gates. **Two milestones, each independently green.** Copy commands from fenced blocks; the executor applies one machine-extracted unified diff.

The current 200 ms probe can report CONTROL `max = 1` when the second dispatch is late and can pass its serialized arm without exercising the mutex. The change waits for arrival events under a real 10 s failure bound, adds an outer arrival gate and a held-mutex `TryLock` check, and replays a 300 ms stagger anchored on the first probe entry. A direct probe handshake tests the missing-arrival message deterministically. This is a Go host test fixture, so no AILANG package or kernel addition is appropriate. S6 is discharged by the mutation drill; S1/S3/S4/S5/S7/S8 add no language, kernel, package, or user-facing surface work.

Final scope: **only `host/daemon/decisionpacket_daemon_test.go` changes.** `daemon.go` stays unchanged; production mutations run solely in scratch archives. The executor adds only `verification/world-iter334/executor-report.md` and `verification/world-iter334/executor-*.txt`. The controller owns git writes, mission documents, moves to implemented/, PR and CI. The planner writes only this Markdown and its JSON; all planner verification is in `.planner-scratch/`.

## 1. Environment and executor rules

```sh
# Run from the executor's assigned repo root (pwd recorded in the report).
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WX=$(mktemp -d /tmp/world-iter334-executor.XXXXXX)
mkdir -p "$WX/tools" "$WX/logs" "$WX/base" "$WX/mutants"
PLAN=design_docs/planned/w-human-approve-control-window-flake-sprint-plan.md
x() { awk -v h="$1" '$0==h{f=1;next} f&&/^```$/{f=0} f' "$PLAN"; }
x '```bash runmut.sh' > "$WX/tools/runmut.sh"
x '```bash loaded-soak.sh' > "$WX/tools/loaded-soak.sh"
chmod +x "$WX/tools/"*.sh
git rev-parse HEAD
git status --porcelain
go version
"$AILANG_BIN" --version
```

Record every command's raw output and exit status. An exit code alone does not prove test discovery: count `--- PASS:` / `--- FAIL:` and inspect both subtests. Commands below use redirects, without pipelines that could hide `go test`'s result. Check each rc before advancing; preserve red output and write the report even when a milestone cannot be green. Tools use Bash explicitly. `WX` is exported for the scripts.

1. No git writes of any kind: no add/commit/stash/checkout/reset/restore/apply/init/worktree; the controller commits. Read-only git archive/status/diff/show/rev-parse are allowed.
2. Apply only the plan single diff with patch -p1 after clean dry-run, full diff SHA256 and base-file checks. Final code touches only host/daemon/decisionpacket_daemon_test.go. daemon.go stays byte-identical to its base pin.
3. Run base and mutation drills only in git archive scratch copies. Temporary daemon.go edits for M1/M1b/M6/M7 must never enter the executor worktree. Restore both target files for every drill and check full before/edited/restored SHA256s.
4. Executor creates verification/world-iter334/executor-report.md and executor-*.txt. Tools, snapshots, temporary copies go under $WX outside the repo. Preserve existing verification evidence and both plan deliverables.
5. Do not edit design_docs/world-mission*.md; do not move plans/designs or touch main checkout, launchd files, or mission state. Controller handles queue/log/commit/PR/CI.
6. Every go test has -timeout; bound other long commands with alarm. Run package-scale gates sequentially; capture raw output and rc without a masking pipeline. Never finish by announcing a wait.
7. Copy exact commands from fenced blocks. Run shell loops/scripts with bash (interactive shell is zsh). No t.Parallel, no additional code or timing changes, no mutant gofmt.
8. A substitution that did not apply, missing test discovery, SKIP, build error, race detector finding, or Go timeout is an instrument/gate error, never a mutation kill. Require named killers and measured line numbers.
9. Any gate red at base is a finding to REPORT, not to work around. Do not fix unrelated issues, change binaries, relax timeouts, initialize git in copies, or claim acceptance green. Baseline read-only in an archive if needed; missing .git can itself fail host/pins and must be reported.
10. Socket-touching results from inside a sandbox are labelled UNINFORMATIVE UNDER SANDBOX; a bind permission failure is neither pass nor fail for the sprint. Executor must obtain informative gates in its authorized environment or report them unverified.
11. Load, if used, consists solely of yes > /dev/null burners. Track and kill every burner you start before returning, including on red/interrupt; never kill unrelated processes. pgrep -x yes must have no PID output; pgrep -x yes | wc -l must print 0.

## 2. Base pins and provenance

Check the full pins before applying; a mismatch stops execution and is reported. HEAD may advance with controller documentation commits only when these Go bytes still match.

| File | HEAD sha256 | Fixed sha256 |
|---|---|---|
| `host/daemon/decisionpacket_daemon_test.go` | `53a708a5264fafe53217d59c58dafa7e25eeabe7f5ef9debc554912fa30a0445` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |

Prototype source: `verification/world-iter334/designer-r2/proto-r2.diff`. Its old header is an absolute `/tmp/...` path and its new header lacks `b/`; it is **not** directly a `-p1` patch. The planner normalized headers to `a/host/...` and `b/host/...`, applied it in scratch and verified the prototype test sha256 `bb75420d0713188621a111bb2079e66a2f4ab91f0529d00e2bc106d1a483d3ee`. The single diff below adds §2.5 to those exact bytes. No production source is changed.

**Plan diff sha256:** `8dee691695d632d8afa3f775629d213bf4769606072a0225140e9ec9f66c6d51`. **Stat:** `1 file changed, 177 insertions(+), 17 deletions(-)` (`194` changed lines). These cover the full prototype plus the seam/test. The JSON repeats these pins and stat.

Base expectations: the ordinary test is normally green but CI-flaky; AC-1's base-only stagger empirically reproduces CONTROL failure at line 270 (base line 267 + 3 inserted lines). The handshake test and `late-second` subtest are absent at base. No claim of deterministic base reproduction is made.

## 3. M1 — apply, compile, target tests and mutations

M1 ends with the fixed test file green, mutation copies restored, and daemon.go unchanged in the executor worktree. AC-1 is run on a base archive before applying. After recording it, execute the following, checking both full base pins and the extracted diff pin before `patch`:

```sh
shasum -a 256 host/daemon/decisionpacket_daemon_test.go host/daemon/daemon.go > "$WX/base.sha"
cat "$WX/base.sha"
x '```diff' > "$WX/plan.diff"
shasum -a 256 "$WX/plan.diff"
# Require 8dee691695d632d8afa3f775629d213bf4769606072a0225140e9ec9f66c6d51; require both base pins below.
patch -p1 --dry-run < "$WX/plan.diff"
patch -p1 < "$WX/plan.diff"
shasum -a 256 host/daemon/decisionpacket_daemon_test.go host/daemon/daemon.go > "$WX/fixed.sha"
cat "$WX/fixed.sha"
perl -e 'alarm 180; exec @ARGV' gofmt -l host/daemon/ > "$WX/logs/gofmt.txt" 2>&1
perl -e 'alarm 300; exec @ARGV' go vet ./host/daemon/ > "$WX/logs/vet.txt" 2>&1
go test -timeout 180s -run '^$' ./host/daemon/ > "$WX/logs/compile.txt" 2>&1
go test -timeout 600s ./host/daemon -run '^(TestHumanApproveIsSerializedAcrossEpisodes|TestConcurrencyProbeNamesAMissedArrival)$' -count=3 -v > "$WX/logs/target-plain.txt" 2>&1
go test -race -timeout 600s ./host/daemon -run '^(TestHumanApproveIsSerializedAcrossEpisodes|TestConcurrencyProbeNamesAMissedArrival)$' -count=3 -v > "$WX/logs/target-race.txt" 2>&1
```

`gofmt -l` must be empty; vet and compile-only must return 0. `go build ./...` is not a compile fence for `_test.go`. Both named tests must PASS 3/3 plain and race. Require the exact fixed SHA256s; do not edit code to reconcile a mismatch.

### 3.1 The single executor code diff

```diff
--- a/host/daemon/decisionpacket_daemon_test.go
+++ b/host/daemon/decisionpacket_daemon_test.go
@@ -3,6 +3,7 @@
 import (
 	"bytes"
 	"context"
+	"fmt"
 	"net/http"
 	"sync"
 	"testing"
@@ -196,25 +197,79 @@
 	}
 }
 
+// probeWaitBound bounds how long the first Human.Approve execution waits for
+// the second dispatch (row 235). It returns on the arrival, so it costs nothing
+// when healthy; it sits below dispatchPair's bound (+5 s) and the test context
+// (30 s), so a missing second dispatch fails by name, never as "max = 1".
+const probeWaitBound = 10 * time.Second
+
+// serializedHold is how long the first holder keeps the mutex after the second
+// dispatch has reached the arrivalGate: a mutex-less wrapper lets the second
+// in, which the probe sees as max = 2 (secondary to the TryLock check).
+const serializedHold = 50 * time.Millisecond
+
+// lateSecond is the regression arm's stagger: the 250-400 ms gap that read as
+// "unwrapped probe max = 1" on the 200 ms probe (iteration 334 measurement).
+// It is measured from the first probe entry (entered), not from goroutine start.
+const lateSecond = 300 * time.Millisecond
+
+// arrivalGate counts Human.Approve arrivals OUTSIDE the serialization, so the
+// serialized arm can show the second dispatch arrived while the first was
+// still inside Execute (it does not observe the second calling Lock).
+type arrivalGate struct {
+	h      broker.Handler
+	n      atomic.Int32
+	second chan struct{}
+}
+
+func (g *arrivalGate) Execute(ctx context.Context, req broker.EffectRequest, payload []byte) ([]byte, error) {
+	if g.n.Add(1) == 2 {
+		close(g.second)
+	}
+	return g.h.Execute(ctx, req, payload)
+}
+
+// concurrencyProbe records how many executions overlap. Its first entry closes
+// entered and holds until the second dispatch arrives: inside the probe
+// (CONTROL, arrived nil) or at the gate (serialized, then serializedHold).
+// held is set only in the serialized arm: each entry then checks it is held
+// and counts in unheld any entry where TryLock succeeded. CONTROL's probe has
+// held == nil, so the check and the counter are off there.
 type concurrencyProbe struct {
-	mu       sync.Mutex
-	cur, max int
-	second   chan struct{}
-	once     sync.Once
+	mu                sync.Mutex
+	cur, max, entries int
+	unheld            int
+	missed            string
+	inside2           chan struct{}
+	entered           chan struct{}
+	once              sync.Once
+	arrived           <-chan struct{}
+	held              *sync.Mutex
+	timeout           <-chan time.Time // test-only; nil uses the real probeWaitBound
 }
 
 func (p *concurrencyProbe) Execute(ctx context.Context, _ broker.EffectRequest, _ []byte) ([]byte, error) {
+	unheld := p.held != nil && p.held.TryLock()
+	if unheld {
+		p.held.Unlock()
+	}
 	p.mu.Lock()
+	if unheld {
+		p.unheld++
+	}
+	first := p.entries == 0
+	if first {
+		close(p.entered)
+	}
+	p.entries++
 	p.cur++
 	p.max = max(p.max, p.cur)
 	if p.cur == 2 {
-		p.once.Do(func() { close(p.second) })
+		p.once.Do(func() { close(p.inside2) })
 	}
 	p.mu.Unlock()
-	select {
-	case <-p.second:
-	case <-time.After(200 * time.Millisecond):
-	case <-ctx.Done():
+	if first {
+		p.awaitSecond(ctx)
 	}
 	p.mu.Lock()
 	p.cur--
@@ -222,14 +277,54 @@
 	return []byte(`{"status":"pending","requestRef":"` + hashref.SumSHA256([]byte("ac8")).String() + `"}`), nil
 }
 
+func (p *concurrencyProbe) awaitSecond(ctx context.Context) {
+	var wait <-chan struct{} = p.inside2
+	if p.arrived != nil {
+		wait = p.arrived
+	}
+	timeout := p.timeout
+	if timeout == nil {
+		timeout = time.After(probeWaitBound)
+	}
+	why := ""
+	select {
+	case <-wait:
+		if p.arrived != nil {
+			select {
+			case <-p.inside2:
+			case <-time.After(serializedHold):
+			case <-ctx.Done():
+			}
+		}
+		return
+	case <-timeout:
+		why = fmt.Sprintf("second Human.Approve dispatch did not arrive within %s", probeWaitBound)
+	case <-ctx.Done():
+		why = "test context ended before the second Human.Approve dispatch arrived: " + ctx.Err().Error()
+	}
+	p.mu.Lock()
+	p.missed = why
+	p.mu.Unlock()
+}
+
 // TestHumanApproveIsSerializedAcrossEpisodes (AC8): one daemon mutex covers every episode.
+// Row 235: each arm holds its first execution until the second dispatch arrives,
+// and "late-second" replays the CI gap: the second dispatch starts lateSecond
+// after the first has entered the probe.
 func TestHumanApproveIsSerializedAcrossEpisodes(t *testing.T) {
+	for _, tc := range []struct {
+		name    string
+		stagger time.Duration
+	}{{"together", 0}, {"late-second", lateSecond}} {
+		t.Run(tc.name, func(t *testing.T) { humanApproveSerializedArms(t, tc.stagger) })
+	}
+}
+
+func humanApproveSerializedArms(t *testing.T, stagger time.Duration) {
 	d := newHandlerDaemon(t)
 	publishHumanTool(t, d)
 	sh := d.human.(serializedHandler)
-	probe := &concurrencyProbe{second: make(chan struct{})}
-	d.human = serializedHandler{mu: sh.mu, h: probe}
-	dispatchPair := func(episodeA, episodeB, taskA, taskB string) [2]error {
+	dispatchPair := func(entered <-chan struct{}, episodeA, episodeB, taskA, taskB string) [2]error {
 		c := humanCoordinator(t, d, d.binder, humanPlan("e1"))
 		calls := []coordinator.Call{humanCall(t, d, episodeA, taskA), humanCall(t, d, episodeB, taskB)}
 		ctx := boundedTestContext(t)
@@ -239,6 +334,18 @@
 			wg.Add(1)
 			go func(i int) {
 				defer wg.Done()
+				if i == 1 && stagger > 0 {
+					select {
+					case <-entered:
+					case <-time.After(probeWaitBound):
+						errs[i] = fmt.Errorf("first dispatch did not enter the probe within %s", probeWaitBound)
+						return
+					case <-ctx.Done():
+						errs[i] = ctx.Err()
+						return
+					}
+					time.Sleep(stagger)
+				}
 				_, errs[i] = c.Dispatch(ctx, calls[i])
 			}(i)
 		}
@@ -246,28 +353,81 @@
 		go func() { wg.Wait(); close(done) }()
 		select {
 		case <-done:
-		case <-time.After(5 * time.Second):
+		case <-time.After(probeWaitBound + 5*time.Second):
 			t.Fatal("concurrent Human.Approve dispatches did not finish")
 		case <-ctx.Done():
 			t.Fatal(ctx.Err())
 		}
 		return errs
 	}
-	errs := dispatchPair("epA", "epB", "a", "b")
+	gate := &arrivalGate{second: make(chan struct{})}
+	probe := &concurrencyProbe{inside2: make(chan struct{}), entered: make(chan struct{}), arrived: gate.second, held: sh.mu}
+	gate.h = serializedHandler{mu: sh.mu, h: probe}
+	d.human = gate
+	errs := dispatchPair(probe.entered, "epA", "epB", "a", "b")
+	if probe.missed != "" {
+		t.Fatalf("serialized arm measured nothing: %s (dispatch errors %v)", probe.missed, errs)
+	}
+	if probe.unheld != 0 {
+		t.Fatalf("serialized probe ran %d of %d executions without the daemon mutex held", probe.unheld, probe.entries)
+	}
 	if probe.max != 1 {
 		t.Fatalf("serialized probe max = %d, want 1", probe.max)
 	}
-	t.Logf("serialized probe max = %d", probe.max)
+	t.Logf("serialized probe max = %d, gate arrivals = %d", probe.max, gate.n.Load())
 	requireOneCommitted(t, errs)
 	// CONTROL: an unwrapped fresh probe must observe both episodes concurrently.
-	probe2 := &concurrencyProbe{second: make(chan struct{})}
+	// held is nil here: no TryLock check, no unheld assertion; missed and max only.
+	probe2 := &concurrencyProbe{inside2: make(chan struct{}), entered: make(chan struct{})}
 	d.human = probe2
-	errs = dispatchPair("epC", "epD", "c", "d")
+	errs = dispatchPair(probe2.entered, "epC", "epD", "c", "d")
+	if probe2.missed != "" {
+		t.Fatalf("CONTROL measured nothing: %s (dispatch errors %v)", probe2.missed, errs)
+	}
 	if probe2.max != 2 {
 		t.Fatalf("unwrapped probe max = %d, want 2", probe2.max)
 	}
 	t.Logf("CONTROL unwrapped probe max = %d", probe2.max)
 	requireOneCommitted(t, errs)
+}
+
+// TestConcurrencyProbeNamesAMissedArrival drives the failure bound without
+// depending on scheduling a second dispatch or waiting for the real 10 s timer.
+func TestConcurrencyProbeNamesAMissedArrival(t *testing.T) {
+	ctx := boundedTestContext(t)
+	timeout := make(chan time.Time, 1)
+	probe := &concurrencyProbe{
+		inside2: make(chan struct{}),
+		entered: make(chan struct{}),
+		arrived: make(chan struct{}),
+		timeout: timeout,
+	}
+	bound := time.NewTimer(5 * time.Second)
+	defer bound.Stop()
+	done := make(chan struct{})
+	go func() {
+		_, _ = probe.Execute(ctx, broker.EffectRequest{}, nil)
+		close(done)
+	}()
+	select {
+	case <-probe.entered:
+	case <-bound.C:
+		t.Fatalf("%s: Execute did not enter within 5s", t.Name())
+	}
+	// Buffered so even an ignored timeout seam cannot strand the sender.
+	timeout <- time.Time{}
+	select {
+	case <-done:
+	case <-bound.C:
+		t.Fatalf("%s: Execute did not return after the injected timeout within 5s", t.Name())
+	}
+	want := fmt.Sprintf("second Human.Approve dispatch did not arrive within %s", probeWaitBound)
+	if probe.missed != want {
+		t.Fatalf("missed = %q, want %q", probe.missed, want)
+	}
+	if probe.entries != 1 {
+		t.Fatalf("entries = %d, want 1", probe.entries)
+	}
 }
 
 // requireOneCommitted accepts a world-head conflict on the slower of two
```

### 3.2 Behavioral invariants

The serialized probe has `arrived: gate.second` and `held: sh.mu`; the gate is outside `serializedHandler`. Only the first entry waits. Every serialized entry checks the mutex; assertions check `missed`, then `unheld`, then `max == 1`. CONTROL is a **different probe**, with `held == nil` and `arrived == nil`; it checks `missed` then `max == 2`, never `unheld`. All existing arms leave `timeout` nil, preserving the real `probeWaitBound = 10s`, pair bound `probeWaitBound + 5s`, and a fresh 30s `boundedTestContext(t)` per pair. `requireOneCommitted` still allows the expected world-head conflict.

The new test uses `boundedTestContext(t)`, calls `Execute` on exactly one goroutine, and never closes `arrived` or `inside2`. It waits for `entered`, sends its own buffered timeout event (so an ignored seam cannot block the sender), and waits for `Execute`'s completion acknowledgement within one total 5s timer. Closing `done` synchronizes the `missed`/`entries` reads after the write; assertions require the complete `fmt.Sprintf(..., probeWaitBound)` message and `entries == 1`. A failure cancels the test context via cleanup, so the ignored-seam mutant's goroutine is released without waiting for 10s. No sleeps or second dispatch drive AC-4.

## 4. Exact mutation drill and sole killers

All lines below are **computed on the final diff's resulting file**, or on the exact mutated file when an edit shifts lines. They are not presented as executor measurements. The executor must record **measured** assertion lines and full before/edited/restored SHA256s in its report. Final assertions: serialized missed `369`, unheld `372`, max `375`; CONTROL missed `385`, max `388`; handshake own bound `422`, message `426`, entries `429`. M3 adds four lines, so its CONTROL max killer is `392`. The preexisting production binding assertion remains at test lines `182–183` (the `fmt` import shifts the actual base `181–182` by one; this corrects the design’s citation).

| ID | Exact edit (commands are in the script below) | Sole killer on the final file | Kind / expected verdict |
|---|---|---|---|
| A | Base only: insert a three-line 300 ms sleep before the second Dispatch. | `:270` `unwrapped probe max = 1, want 2` | EMPIRICAL; FAIL |
| M1 | Delete both MD3 Lock/deferred Unlock lines from serializedHandler.Execute. | `:372` `serialized probe ran 2 of 2 executions without the daemon mutex held` | DETERMINISTIC; FAIL |
| M1b | Replace deferred Unlock with immediate Unlock before the inner Execute. | `:372` `serialized probe ran 2 of 2 executions without the daemon mutex held` | DETERMINISTIC; FAIL |
| M2 | gate.h = probe (bare serialized-arm handler). | `:372` `serialized probe ran 2 of 2 executions without the daemon mutex held` | DETERMINISTIC; FAIL |
| M3 | Replace p.awaitSecond(ctx) at Execute with the old inside2/200 ms/context select. | `:392` `unwrapped probe max = 1, want 2` | EMPIRICAL; FAIL |
| M6 | M1 plus disable TryLock with false && (both files in one isolated mutant). | `:375` `serialized probe max = 2, want 1` | EMPIRICAL; FAIL |
| M7 | Add package-level otherMu; Execute locks otherMu instead of s.mu. | `:372` `serialized probe ran 2 of 2 executions without the daemon mutex held` | DETERMINISTIC; FAIL |
| M7b | Consistently rewire test wrapper and held to one otherMu; _ = sh avoids unused-variable build failure. | none; expected PASS | BOUNDARY_EXPECT_PASS; PASS |
| H1 | §2.5(i): timeout branch no longer records a missed reason (why = empty). | `:426` `missed = "", want "second Human.Approve dispatch did not arrive within 10s"` | DETERMINISTIC; FAIL |
| H2 | §2.5(ii): timeout reason drops the bound constant. | `:426` `missed = "second Human.Approve dispatch did not arrive", want "second Human.Approve dispatch did not arrive within 10s"` | DETERMINISTIC; FAIL |
| H3 | §2.5(iii): ignore p.timeout; always use the real bound. | `:422` `TestConcurrencyProbeNamesAMissedArrival: Execute did not return after the injected timeout within 5s` | DETERMINISTIC; FAIL |

M1/M1b/M2/M7 must fail in **both** subtests at `unheld`, before `max` is examined. M7 is not killed by `max`. M6 shows the secondary serialized `max` signal empirically; no deterministic overlap claim is made. M7b must PASS: a consistent alternate mutex still serializes. `TryLock` shows that the supplied mutex is held, not who holds it and not that the second dispatch called Lock. Production composition relies on the existing binding assertion (`sh.mu == &d.humanMu`) and the mutex-parametric wrapper. No assertion here catches consistent test-side rewiring. M5 (100 ms bound) is an empirical demonstration only and is superseded as named-failure evidence by H1–H3; it is not a required deterministic drill. M6b is historical, not an additional required row.

Each `perl -0pi -e` substitution below is literal, applies to the pinned bytes, and is run **alone** except the explicitly compound M6/M7/M7b rows. The runner snapshots both files, refuses an unchanged targeted file, prints SHA256 before/edited/restored, and restores with `cp` via an EXIT trap, including failed test commands. Require edited digests to differ on every targeted file, and restored digests to equal before in full. The JSON lists the exact computed edited digests. Never use git to revert. A build failure is an instrument error, not a kill. Inspect raw logs to confirm the named assertion is the sole killer; the runner checks its expected count and top-level test failure count. Each invocation returns 0 for a correctly observed mutant outcome, 3 for an unexpected outcome requiring inspection/adjudication, and other nonzero codes for protocol errors.

```bash runmut.sh
#!/bin/bash
# runmut.sh <tree> <A|M1|M1b|M2|M3|M6|M7|M7b|H1|H2|H3> <plain|race> <count>
set -euo pipefail
: "${WX:?export WX to the executor scratch directory}"
cd "$1"
id=$2 mode=$3 count=$4
T=host/daemon/decisionpacket_daemon_test.go
G=host/daemon/daemon.go
export AILANG_BIN=$HOME/.pinned-ailang/ailang
save=$(mktemp -d "$WX/mut-${id}-${mode}.XXXXXX")
cp "$T" "$save/test.go"
cp "$G" "$save/daemon.go"
shasum -a 256 "$T" "$G" > "$save/before.sha"
restore() {
    local result=$?
    trap - EXIT INT TERM
    cp "$save/test.go" "$T" || result=90
    cp "$save/daemon.go" "$G" || result=90
    shasum -a 256 "$T" "$G" > "$save/restored.sha"
    printf '%s %s sha256 before:\n' "$id" "$mode"
    cat "$save/before.sha"
    printf 'edited:\n'
    test ! -f "$save/edited.sha" || cat "$save/edited.sha"
    printf 'restored:\n'
    cat "$save/restored.sha"
    cmp -s "$save/before.sha" "$save/restored.sha" || result=91
    # Per-run copies keep every mutation's digest evidence distinct.
    for stage in before edited restored; do
        test ! -f "$save/$stage.sha" || cp "$save/$stage.sha" "$WX/logs/$id-$mode-$stage.sha.txt"
    done
    exit "$result"
}
trap restore EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
case "$id" in
    A)
        perl -0pi -e 's/(\t\t\t\t_, errs\[i\] = c\.Dispatch\(ctx, calls\[i\]\))/\t\t\t\tif i == 1 {\n\t\t\t\t\ttime.Sleep(300 * time.Millisecond)\n\t\t\t\t}\n$1/' "$T"
        if cmp -s "$T" "$save/test.go"; then echo "SUBSTITUTION DID NOT APPLY: A $T"; exit 8; fi
        ;;
    M1)
        perl -0pi -e 's/\ts\.mu\.Lock\(\)         \/\/ MD3 target\n\tdefer s\.mu\.Unlock\(\) \/\/ MD3 target\n//' "$G"
        if cmp -s "$G" "$save/daemon.go"; then echo "SUBSTITUTION DID NOT APPLY: M1 $G"; exit 8; fi
        ;;
    M1b)
        perl -0pi -e 's/defer s\.mu\.Unlock\(\) \/\/ MD3 target/s.mu.Unlock() \/\/ MD3 target/' "$G"
        if cmp -s "$G" "$save/daemon.go"; then echo "SUBSTITUTION DID NOT APPLY: M1b $G"; exit 8; fi
        ;;
    M2)
        perl -0pi -e 's/gate\.h = serializedHandler\{mu: sh\.mu, h: probe\}/gate.h = probe/' "$T"
        if cmp -s "$T" "$save/test.go"; then echo "SUBSTITUTION DID NOT APPLY: M2 $T"; exit 8; fi
        ;;
    M3)
        perl -0pi -e 's/\t\tp\.awaitSecond\(ctx\)\n/\t\tselect {\n\t\tcase <-p.inside2:\n\t\tcase <-time.After(200 * time.Millisecond):\n\t\tcase <-ctx.Done():\n\t\t}\n/' "$T"
        if cmp -s "$T" "$save/test.go"; then echo "SUBSTITUTION DID NOT APPLY: M3 $T"; exit 8; fi
        ;;
    M6)
        perl -0pi -e 's/\ts\.mu\.Lock\(\)         \/\/ MD3 target\n\tdefer s\.mu\.Unlock\(\) \/\/ MD3 target\n//' "$G"
        perl -0pi -e 's/unheld := p\.held != nil && p\.held\.TryLock\(\)/unheld := false \&\& p.held != nil \&\& p.held.TryLock()/' "$T"
        if cmp -s "$G" "$save/daemon.go"; then echo "SUBSTITUTION DID NOT APPLY: M6 $G"; exit 8; fi
        if cmp -s "$T" "$save/test.go"; then echo "SUBSTITUTION DID NOT APPLY: M6 $T"; exit 8; fi
        ;;
    M7)
        perl -0pi -e 's/func \(s serializedHandler\) Execute/var otherMu sync.Mutex\n\nfunc (s serializedHandler) Execute/; s/\ts\.mu\.Lock\(\)         \/\/ MD3 target\n\tdefer s\.mu\.Unlock\(\) \/\/ MD3 target/\totherMu.Lock()         \/\/ MD3 target\n\tdefer otherMu.Unlock() \/\/ MD3 target/' "$G"
        if cmp -s "$G" "$save/daemon.go"; then echo "SUBSTITUTION DID NOT APPLY: M7 $G"; exit 8; fi
        ;;
    M7b)
        perl -0pi -e 's/\tsh := d\.human\.\(serializedHandler\)\n/\tsh := d.human.(serializedHandler)\n\t_ = sh\n\totherMu := \&sync.Mutex{}\n/; s/held: sh\.mu/held: otherMu/; s/gate\.h = serializedHandler\{mu: sh\.mu, h: probe\}/gate.h = serializedHandler{mu: otherMu, h: probe}/' "$T"
        if cmp -s "$T" "$save/test.go"; then echo "SUBSTITUTION DID NOT APPLY: M7b $T"; exit 8; fi
        ;;
    H1)
        perl -0pi -e 's/(case <-timeout:\n\t\t)why = fmt\.Sprintf\("second Human\.Approve dispatch did not arrive within %s", probeWaitBound\)/$1why = ""/' "$T"
        if cmp -s "$T" "$save/test.go"; then echo "SUBSTITUTION DID NOT APPLY: H1 $T"; exit 8; fi
        ;;
    H2)
        perl -0pi -e 's/(case <-timeout:\n\t\t)why = fmt\.Sprintf\("second Human\.Approve dispatch did not arrive within %s", probeWaitBound\)/$1why = "second Human.Approve dispatch did not arrive"/' "$T"
        if cmp -s "$T" "$save/test.go"; then echo "SUBSTITUTION DID NOT APPLY: H2 $T"; exit 8; fi
        ;;
    H3)
        perl -0pi -e 's/timeout := p\.timeout/var timeout <-chan time.Time/' "$T"
        if cmp -s "$T" "$save/test.go"; then echo "SUBSTITUTION DID NOT APPLY: H3 $T"; exit 8; fi
        ;;
    *) echo "unknown mutant $id"; exit 9 ;;
esac
shasum -a 256 "$T" "$G" > "$save/edited.sha"
cat "$save/edited.sha"
flags=()
case "$mode" in plain) ;; race) flags=(-race) ;; *) exit 9 ;; esac
regex='^TestHumanApproveIsSerializedAcrossEpisodes$'
case "$id" in H*) regex='^TestConcurrencyProbeNamesAMissedArrival$' ;; esac
out="$WX/logs/$id-$mode.txt"
set +e
go test "${flags[@]}" -timeout 600s ./host/daemon -run "$regex" -count="$count" -v > "$out" 2>&1
rc=$?
set -e
cat "$out"
printf '%s mode=%s count=%s rc=%s\n' "$id" "$mode" "$count" "$rc"
if rg -q 'build failed|^# |DATA RACE|panic:|test timed out|operation not permitted|permission denied' "$out"; then
    echo 'INSTRUMENT/GATE ERROR: never count this as a mutation kill; socket restrictions => UNINFORMATIVE UNDER SANDBOX'
    exit 7
fi
# Expected red is a green mutation check only when the named assertion kills it.
# Adjudicate any empirical M3 survival per the plan; this script never silently accepts it.
if [ "$id" = M7b ]; then
    [ "$rc" -eq 0 ] || exit 3
    [ "$(grep -c '^--- PASS: TestHumanApproveIsSerializedAcrossEpisodes ' "$out" || true)" -eq "$count" ] || exit 3
else
    [ "$rc" -eq 1 ] || exit 3
    case "$id" in
        A|M3) killer='unwrapped probe max = 1, want 2' ;;
        M1|M1b|M2|M7) killer='serialized probe ran 2 of 2 executions without the daemon mutex held' ;;
        M6) killer='serialized probe max = 2, want 1' ;;
        H1) killer='missed = "", want "second Human.Approve dispatch did not arrive within 10s"' ;;
        H2) killer='missed = "second Human.Approve dispatch did not arrive", want "second Human.Approve dispatch did not arrive within 10s"' ;;
        H3) killer='TestConcurrencyProbeNamesAMissedArrival: Execute did not return after the injected timeout within 5s' ;;
    esac
    want=$count
    case "$id" in M1|M1b|M2|M6|M7) want=$((count * 2)) ;; esac
    [ "$(grep -F -c "$killer" "$out" || true)" -eq "$want" ] || exit 3
    case "$id" in
        H*) name=TestConcurrencyProbeNamesAMissedArrival ;;
        *) name=TestHumanApproveIsSerializedAcrossEpisodes ;;
    esac
    [ "$(grep -c "^--- FAIL: $name " "$out" || true)" -eq "$count" ] || exit 3
    if [ "$id" = M3 ]; then
        [ "$(grep -c -- '--- PASS: TestHumanApproveIsSerializedAcrossEpisodes/together ' "$out" || true)" -eq "$count" ] || exit 3
        [ "$(grep -c -- '--- FAIL: TestHumanApproveIsSerializedAcrossEpisodes/late-second ' "$out" || true)" -eq "$count" ] || exit 3
    fi
fi
```

Run the additional required table rows in M1 after AC-2 creates the fixed archive:

```bash
for id in M6 M7 M7b; do
    for mode in plain race; do
        "$WX/tools/runmut.sh" "$WX/mutants" "$id" "$mode" 20 > "$WX/logs/$id-$mode-protocol.txt" 2>&1 || exit $?
    done
done
shasum -a 256 -c "$WX/fixed.sha"
# Snapshot equality proves daemon.go remained unchanged in the executor worktree.
git diff --name-only
git diff --stat
```

## 5. Acceptance commands (design §5, AC-4 refined by §2.5)

Execute all loops using Bash. Every Go test command has `-timeout`. Within mutation scripts the exact Go command is `go test "${flags[@]}" -timeout 600s ./host/daemon -run "$regex" -count="$count" -v`; `plain` expands no race flag and `race` expands `-race`, with the regex/count shown by each call. These are the design's plain and race AC-1 command for the mutex/window drills. Commands execute from repo root or the specified archive root. Mutation drills go under M1; soak/package/full gate under M2.

### AC-1 — Base reproduction (EMPIRICAL; M1)

```bash
git archive HEAD | tar -x -C "$WX/base"
# A inserts exactly the three-line 300 ms sleep before the base second Dispatch.
"$WX/tools/runmut.sh" "$WX/base" A plain 20 > "$WX/logs/A-plain-protocol.txt" 2>&1
"$WX/tools/runmut.sh" "$WX/base" A race 20 > "$WX/logs/A-race-protocol.txt" 2>&1
```

Demonstration only: want 20/20 FAIL plain and race at base harness line 270: unwrapped probe max = 1, want 2; restore exactly. Survival is empirical evidence, never a gate on the fix.

### AC-2 — Mutex removal (DETERMINISTIC; M1)

```bash
# Fixed mutation copy: all production mutations are temporary here, not in the worktree.
git archive HEAD | tar -x -C "$WX/mutants"
patch -d "$WX/mutants" -p1 < "$WX/plan.diff"
for id in M1 M1b M2; do
    for mode in plain race; do
        "$WX/tools/runmut.sh" "$WX/mutants" "$id" "$mode" 20 > "$WX/logs/$id-$mode-protocol.txt" 2>&1 || exit $?
    done
done
```

20/20 top-level FAIL and 40/40 subtest FAIL per mutant/mode; sole killer unheld at final test line 372. No build errors, races, or skips; hashes restored.

### AC-3 — Old 200 ms window (EMPIRICAL; M1)

```bash
"$WX/tools/runmut.sh" "$WX/mutants" M3 plain 20 > "$WX/logs/M3-plain-protocol.txt" 2>&1
"$WX/tools/runmut.sh" "$WX/mutants" M3 race 20 > "$WX/logs/M3-race-protocol.txt" 2>&1
```

Want 20/20 FAIL in late-second only; together passes. Sole killer CONTROL max at mutant line 392. Apply the exact r3 adjudication rule below for a lone evidenced residual; a second survival fails the row.

**Adjudication (design r3, verbatim):** a lone M3 survival whose log shows the second dispatch entering only after the first goroutine resumed (>100 ms deschedule past its expired timer) is the R2 residual, not a regression — re-run once at -count=20; a second survival fails the row. AC-1 is a demonstration on the base tree, not a gate on the fix.

The three scheduling residuals are delayed timer creation after closing `entered`, delayed decrement after expiry, and random selection when arrival and timeout are both ready. If a survival log does not establish the quoted residual, do not infer it; report it as unresolved. Only one justified rerun is allowed, using the same 20-count plain/race invocation for the affected mode, with a distinct saved raw log. Do not tune the stagger or rerun until green.

### AC-4 — Named-failure handshake (DETERMINISTIC; M1)

```bash
go test -timeout 600s ./host/daemon -run '^TestConcurrencyProbeNamesAMissedArrival$' -count=200 -v > "$WX/logs/handshake-plain.txt" 2>&1
go test -race -timeout 600s ./host/daemon -run '^TestConcurrencyProbeNamesAMissedArrival$' -count=200 -v > "$WX/logs/handshake-race.txt" 2>&1
for id in H1 H2 H3; do
    for mode in plain race; do
        "$WX/tools/runmut.sh" "$WX/mutants" "$id" "$mode" 3 > "$WX/logs/$id-$mode-protocol.txt" 2>&1 || exit $?
    done
done
```

Control 200/200 PASS each mode; H1/H2/H3 fail 3/3 each mode by their sole killers. H3 fails by the test own 5s bound at line 422, never the 600s Go timeout. H1/H2 line 426 checks the complete constant-bearing message; entries == 1 line 429. M5 is empirical only, not deterministic acceptance.

### AC-5 — Idle and loaded soak (GATE; M2)

```bash
go test -race -timeout 900s ./host/daemon -run '^TestHumanApproveIsSerializedAcrossEpisodes$' -count=200 > "$WX/logs/soak-idle.txt" 2>&1
"$WX/tools/loaded-soak.sh" "$PWD" > "$WX/logs/soak-loaded.txt" 2>&1
```

Both ok, 0 FAIL, 0 DATA RACE; count=200 each. Loaded run has exactly 32 yes burners, all killed before return; pgrep -x yes produces no PIDs and its count prints 0.

### AC-6 — Package and full gate (GATE; M2)

```bash
perl -e 'alarm 300; exec @ARGV' go vet ./host/daemon/ > "$WX/logs/package-vet.txt" 2>&1
go test -timeout 900s ./host/daemon/ -count=1 > "$WX/logs/package-plain.txt" 2>&1
go test -race -timeout 1500s ./host/daemon/ -count=1 > "$WX/logs/package-race.txt" 2>&1
perl -e 'alarm 2400; exec @ARGV' bash > "$WX/logs/full-gate.txt" 2>&1 <<'FULL'
go build ./... && go test -timeout 1800s ./...
FULL
(
    unset AILANG_BIN
    export PATH=$HOME/.ailang/state/ailang-release-v0.55.0/x:$PATH
    export WORLD_PKG_AILANG_BIN=$HOME/.ailang/state/ailang-release-v0.41.0/x/ailang
    perl -e 'alarm 600; exec @ARGV' ./scripts/verify_ail.sh
) > "$WX/logs/verify-ail.txt" 2>&1
```

All rc 0, no skips from missing binaries, no DATA RACE; verify_ail.sh required identity manifests pass. CI on fix head all 3 jobs green with ok host/daemon in race gate is controller-only follow-up. A base-red gate is a finding to REPORT, never work around.

### AC-7 — Idle wall time (GATE; M1)

```bash
go test -timeout 600s ./host/daemon -run '^TestHumanApproveIsSerializedAcrossEpisodes$' -count=1 -v > "$WX/logs/wall-time.txt" 2>&1
```

The whole top-level test elapsed time in its --- PASS line is under 1.2s in one idle run, both subtests discovered. Designer measured 0.78-0.80s; planner scratch measured 0.79s. Record measured elapsed; no load burners during this check.


## 6. M2 — soak, package, full gate and report

M2 makes no code changes. Run AC-5 then AC-6 sequentially. The loaded soak uses exactly the same 900s timeout and count as the idle soak; the following script is the only approved load generation:

```bash loaded-soak.sh
#!/bin/bash
# loaded-soak.sh <fixed-worktree>
set -euo pipefail
: "${WX:?}"
cd "$1"
export AILANG_BIN=$HOME/.pinned-ailang/ailang
# Existing burners cannot be attributed to this run; stop, report, do not kill others.
scan_rc=0
pgrep -x yes > "$WX/logs/yes-before.txt" 2> "$WX/logs/yes-before-error.txt" || scan_rc=$?
case "$scan_rc" in
    0) cat "$WX/logs/yes-before.txt"; echo 'pre-existing yes processes: report and stop'; exit 8 ;;
    1) ;; # no matching process
    *) echo 'process enumeration unavailable: UNINFORMATIVE UNDER SANDBOX; no load started'; exit 7 ;;
esac
pids=()
cleanup() {
    local result=$?
    trap - EXIT INT TERM
    if [ "${#pids[@]}" -gt 0 ]; then kill "${pids[@]}" 2>/dev/null || true; fi
    scan_rc=0
    for attempt in {1..50}; do
        scan_rc=0
        pgrep -x yes > "$WX/logs/yes-after.txt" 2> "$WX/logs/yes-after-error.txt" || scan_rc=$?
        [ "$scan_rc" -eq 0 ] || break
        sleep 0.1
    done
    if [ "$scan_rc" -gt 1 ]; then
        echo 'process enumeration unavailable: UNINFORMATIVE UNDER SANDBOX; cleanup count unverified'
        result=9
    else
        echo 'pgrep -x yes (must print no PIDs):'
        cat "$WX/logs/yes-after.txt"
        left=$(wc -l < "$WX/logs/yes-after.txt" | tr -d ' ')
        printf 'yes remaining: %s (must be 0)\n' "$left"
        [ "$left" -eq 0 ] || result=9
    fi
    exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
for i in {1..32}; do yes > /dev/null & pids+=("$!"); done
sleep 5
sysctl -n vm.loadavg || true
go test -race -timeout 900s ./host/daemon -run '^TestHumanApproveIsSerializedAcrossEpisodes$' -count=200
```

Full-gate failures must be reported. If a failure may predate the diff, run the same named command with the same binary/environment and an explicit timeout on a base archive under `$WX/base`, recording its result. A `git archive` copy has no `.git`; e.g. `host/pins` can fail solely because it shells out to `git ls-files`. Label that artifact; **do not git init or commit in a copy to work around it**. Do not change unrelated code, weaken assertions, or suppress a red gate. Any base red is a finding to REPORT and prevents an unqualified green claim. Socket-touching results inside a sandbox must read **UNINFORMATIVE UNDER SANDBOX**. Missing binaries/skipped tests likewise cannot discharge acceptance. Controller-only CI acceptance is all three jobs green on the fix head, including the daemon race-gate `ok` line; executor reports local results and leaves this handoff explicit.

Write `verification/world-iter334/executor-report.md` with: authority/base HEAD/environment; base/fixed/diff pins and stat; each milestone verdict; each AC command and rc, discovered PASS/FAIL counts and elapsed time; each mutant's exact substitution, measured sole-killer line, full SHA256 before/edited/restored and count; M3 empirical adjudication evidence if any; M7b expected PASS boundary; all gate findings with base comparisons; sandbox labels; final scope and burner count; and CI/controller handoff. Save **raw** output of every command, including red ones, as `verification/world-iter334/executor-*.txt`. Append each command's captured rc to its log and record elapsed wall time, preserving the original rc when classifying. The report remains required if a gate prevents completion.

Final evidence collection (after the report is written, record the scope command output too):

```bash
mkdir -p verification/world-iter334
for file in "$WX/logs/"*.txt; do
    cp "$file" "verification/world-iter334/executor-$(basename "$file")"
done
shasum -a 256 -c "$WX/fixed.sha"
git diff --name-only
git diff --stat
git status --porcelain
scan_rc=0
pgrep -x yes > "$WX/logs/final-yes.txt" 2> "$WX/logs/final-yes-error.txt" || scan_rc=$?
if [ "$scan_rc" -le 1 ]; then
    wc -l < "$WX/logs/final-yes.txt"
else
    echo 'UNINFORMATIVE UNDER SANDBOX: process count unavailable'
fi
```

M1 green means compile/tests and deterministic mutations are green, all mutation restoration checks pass, empirical evidence has been recorded/adjudicated, and no production source differs. M2 green means informative soak, package and full gates pass, the report/raw logs exist, fixed pins are unchanged, and no load processes remain. If a base finding or sandbox limit prevents a gate from being informative/green, mark that milestone unverified or blocked with the evidence; do not claim green.

## 7. Planner verification (workspace-write sandbox)

The planner performed only scratch verification; this is not the executor's unsandboxed acceptance run. Commands and raw local scratch logs are under `.planner-scratch/` (controller deletes that directory). The authoritative Go files in the worktree were never edited. No git write operation was used.

- **Environment:** workspace-write sandbox; go1.26.6 darwin/arm64; AILANG_BIN=$HOME/.pinned-ailang/ailang reports v0.41.0. All following results are sandbox-local, not the executor full gate.
- **Base and scope:** Both worktree Go files equal their HEAD SHA256 pins after every check. No worktree Go file was edited; scratch Go edits occurred only under .planner-scratch/. No git write operation was run.
- **Prototype provenance:** Normalized designer r2 headers, applied in scratch, obtained bb75420d0713188621a111bb2079e66a2f4ab91f0529d00e2bc106d1a483d3ee before adding §2.5.
- **Exact diff round trip:** PASS in sandbox: extracted the unique diff and both tool blocks byte-identically; patch -p1 --dry-run clean against this worktree; applied the extracted final diff to a fresh git archive HEAD in .planner-scratch/roundtrip; both resulting file pins match. Stat: 1 file changed, 177 insertions(+), 17 deletions(-).
- **Formatting and compile fence:** PASS in sandbox: gofmt -l host/daemon/ empty (rc 0), go vet ./host/daemon/ rc 0, go test -timeout 180s -run ^$ ./host/daemon/ rc 0. gofmt/vet externally bounded at 180s/300s.
- **Target tests:** PASS in sandbox: go test -timeout 180s ./host/daemon -run the two exact named tests -count=3 -v, plain and -race: 3/3 top-level PASS for each named test and 6 serialization subtest PASS each mode; no DATA RACE, socket/bind error, or SKIP. Plain whole serialization test 0.79s in each run.
- **AC-4 controls:** PASS in sandbox: go test -timeout 600s ./host/daemon -run ^TestConcurrencyProbeNamesAMissedArrival$ -count=200 -v, plain and -race; 200/200 PASS each mode.
- **Mutation smoke drill:** PASS in sandbox, race count=1 per row: M1/M1b/M2/M7 killed solely at unheld :372; M3 killed late-second only at CONTROL max :392; M6 killed at serialized max :375; M7b expected PASS; H1/H2 killed at exact message :426; H3 killed by own 5s bound :422 (not Go timeout). All before/edited/restored full digests matched the JSON computed values; every restoration byte-identical. This smoke drill does not substitute for executor count=20 or AC-4 mutant count=3 plain/race.
- **Static mutation and document checks:** All 11 exact Perl row substitutions change the intended pinned files; computed assertion lines verified on final/mutated bytes. bash -n passes for both extracted scripts. jq -e . passes for the JSON; exactly one diff block and only the intended test path.
- **Not run by planner:** AC-1 base demonstration, full mutation repetition budgets, AC-5 soaks/load generation, AC-6 package/full/AILANG gates and CI remain executor/controller work. No load processes were started. Socket-touching results inside this sandbox would be UNINFORMATIVE UNDER SANDBOX; no socket gate result is claimed.
- **Process enumeration:** No burners were started. A read-only pgrep -x yes check cannot enumerate processes in this sandbox (sysmond service not found / Cannot get process list); the piped count 0 is UNINFORMATIVE UNDER SANDBOX, not cleanup proof. The planned load script refuses to start burners if enumeration is unavailable and reports an unverified cleanup count on enumeration failure. Extracted loaded-soak.sh preflight verified rc=7 with UNINFORMATIVE UNDER SANDBOX before launching any burner.

## 8. Controller handoff

The controller alone commits, updates `design_docs/world-mission*.md`, moves completed design/plan files, opens/merges the PR and reads CI. It must obtain all three green jobs on the final fix head; a sandbox or base finding is carried as a finding rather than patched around by the executor.

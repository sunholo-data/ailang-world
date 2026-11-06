package broker

// Row 250 (design_docs/planned/w-exec-probe-foreign-dial-false-refusal.md):
// arm 6 judges each forward leg on the accepts in its own window, each
// classified by the first line it carries. These tests dial World's live
// listener FROM THE TEST PROCESS while a leg's run is in flight: the probe
// client below hands the test each connect/proxy op's argument and waits for
// the test's go before it prints its token, so the dial is inside the window
// whatever the rig's own loopback dialers do. Like the other shim tests they
// need a real node (probeNode skips by name without one).

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// net6DialClient answers every probe op as a confining sandbox would: the
// first write (arm 1, inside) lands, every other op refuses, listen binds
// nothing. A connect or proxy op writes "op arg" to DIR/op-N and waits for
// DIR/go-N.
const net6DialClient = `#!/bin/sh
dir=DIR
while [ "$#" -gt 0 ] && [ "$1" != "--" ]; do shift; done; shift
first=1
while [ "$#" -ge 2 ]; do op=$1; x=$2; shift 2
  case $op in
    write) if [ $first = 1 ]; then printf WORLD-PROBE > "$x"; echo WORLD-PROBE-WROTE; else echo "WORLD-PROBE-REFUSED EPERM"; fi; first=0;;
    read|listen) echo "WORLD-PROBE-REFUSED EPERM";;
    connect|proxy) n=0; while [ -e "$dir/op-$n" ]; do n=$((n+1)); done
      echo "$op $x" > "$dir/tmp-$n"; /bin/mv "$dir/tmp-$n" "$dir/op-$n"
      while [ ! -e "$dir/go-$n" ]; do /bin/sleep 0.01; done
      echo "WORLD-PROBE-REFUSED ECONNREFUSED";;
    *) echo "WORLD-PROBE-ERROR EBADOP";;
  esac
done
`

// net6ForeignRequest is what the rig's start page sends on every dial
// (design V3): an HTTP GET.
const net6ForeignRequest = "GET / HTTP/1.1\r\nHost: 127.0.0.1\r\nUser-Agent: SunholoHome/1\r\n\r\n"

// net6Dial makes one test-process dial of shape kind to port; nonce is the
// child's, as the op's argument carried it. It returns what it did.
func net6Dial(port, nonce, kind string) string {
	c, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), 2*time.Second)
	if err != nil {
		return kind + ":dial-error:" + err.Error()
	}
	defer c.Close()
	_ = c.SetDeadline(time.Now().Add(3 * time.Second))
	switch kind {
	case "foreign": // World reads the first byte, then closes
		if _, err := io.WriteString(c, net6ForeignRequest); err != nil {
			return kind + ":write-error"
		}
		_, _ = io.Copy(io.Discard, c)
	case "child": // a reach the child's own token denies
		if _, err := io.WriteString(c, "WORLD-PROBE-DIAL "+nonce+"\n"); err != nil {
			return kind + ":write-error"
		}
	case "silent": // accepted, nothing written until well past World's 500 ms read
		time.Sleep(700 * time.Millisecond)
	case "scan": // connect and close: zero bytes, then EOF
	}
	return kind + ":ok"
}

// net6Rig runs ProbeExecProfiles on the pass-through shim with
// net6DialClient; plan[op][i] is the dial shape (or kinds joined by "+")
// the test makes during that op's i-th run ("" or absent: none).
type net6Rig struct {
	ops   []string // "op arg" per connect/proxy run, in order
	dials []string
}

func runNet6Rig(t *testing.T, plan map[string][]string) (*net6Rig, error) {
	t.Helper()
	f := newProbeFixture(t, shimSrtTree(t, "passthrough"), nil, "", []any{"sh", "-c", "echo toolchain-ok"})
	dir := filepath.Join(canonicalTempDir(t), "ops")
	mustMkdir(t, dir)
	client := filepath.Join(canonicalTempDir(t), "client.sh")
	mustWrite(t, client, strings.Replace(net6DialClient, "DIR", strconv.Quote(dir), 1))
	if err := os.Chmod(client, 0o700); err != nil {
		t.Fatal(err)
	}
	rig := &net6Rig{}
	var mu sync.Mutex
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		seen := map[string]int{}
		for n := 0; ; {
			b, err := os.ReadFile(filepath.Join(dir, fmt.Sprintf("op-%d", n)))
			if err != nil {
				select {
				case <-stop:
					return
				case <-time.After(5 * time.Millisecond):
				}
				continue
			}
			op, arg, _ := strings.Cut(strings.TrimSpace(string(b)), " ")
			port, nonce, _ := strings.Cut(arg, ":")
			i := seen[op]
			seen[op]++
			var did []string
			if i < len(plan[op]) && plan[op][i] != "" {
				for _, kind := range strings.Split(plan[op][i], "+") {
					did = append(did, net6Dial(port, nonce, kind))
				}
			}
			mu.Lock()
			rig.ops = append(rig.ops, op+" "+arg)
			rig.dials = append(rig.dials, strings.Join(did, "+"))
			mu.Unlock()
			mustWrite(t, filepath.Join(dir, fmt.Sprintf("go-%d", n)), "")
			n++
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err := ProbeExecProfiles(ctx, ExecProbeConfig{Sandbox: f.sandbox, Paths: f.paths, Client: client}, []*ExecProfile{f.profile})
	close(stop)
	<-done
	f.assertProbeCleaned(t)
	return rig, err
}

func (r *net6Rig) count(op string) int {
	n := 0
	for _, o := range r.ops {
		if strings.HasPrefix(o, op+" ") {
			n++
		}
	}
	return n
}

// arm6Why is arm 6's refusal text, "" when arm 6 passed; any other failed arm fails the test.
func arm6Why(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		return ""
	}
	var perr *ExecProbeError
	if !errors.As(err, &perr) {
		t.Fatalf("probe error = %v (%T), want nil or an *ExecProbeError", err, err)
	}
	why := ""
	for _, f := range perr.Failed {
		if f.Arm != ExecProbeArmNetwork {
			t.Fatalf("arm %s failed under a client that answers as a confining sandbox: %s", f.Arm, f.Why)
		}
		why = f.Why
	}
	return why
}

// A foreign dial, the start page's GET, inside a leg's window is retried,
// and the retry's clean window passes: no false refusal (RED at base, where
// it reads "accepts 2").
func TestExecProbeNet6ForeignDialIsRetriedNotRefused(t *testing.T) {
	for _, row := range []struct {
		name  string
		plan  map[string][]string
		calls map[string]int
	}{
		{"none", nil, map[string]int{"connect": 1, "proxy": 1}},
		{"raw", map[string][]string{"connect": {"foreign"}}, map[string]int{"connect": 2, "proxy": 1}},
		{"proxy", map[string][]string{"proxy": {"foreign"}}, map[string]int{"connect": 1, "proxy": 2}},
		{"raw-twice-proxy", map[string][]string{"connect": {"foreign", "foreign+foreign"}, "proxy": {"foreign"}}, map[string]int{"connect": 3, "proxy": 2}},
	} {
		t.Run(row.name, func(t *testing.T) {
			rig, err := runNet6Rig(t, row.plan)
			t.Logf("ops %q dials %q", rig.ops, rig.dials)
			if why := arm6Why(t, err); why != "" {
				t.Fatalf("arm 6 refused on foreign dials alone: %s", why)
			}
			// A rig dialer (the start page, design V3) can add a retry, never remove one.
			for op, want := range row.calls {
				if got := rig.count(op); got < want || got > execProbeAttempts {
					t.Errorf("%s runs = %d, want %d (at most %d)", op, got, want, execProbeAttempts)
				}
			}
			for _, d := range rig.dials {
				if d != "" && strings.Contains(d, "error") {
					t.Fatalf("a test dial did not land: %q", d)
				}
			}
		})
	}
}

// Every accept that is not provably foreign still fails arm 6 at once,
// with no retry: the child's line (a reach its own token denies), a silent
// accept, a connect-and-close, and a child in the same window as a foreign
// dial or on the retry after one.
func TestExecProbeNet6ChildOrSilentAcceptStillFails(t *testing.T) {
	const raw, prx = "a raw connect to World's live listener", "a connect through srt's proxy environment"
	for _, row := range []struct {
		name, op, want string
		plan           []string
		runs           int
	}{
		{"raw-child", "connect", raw + `: token "WORLD-PROBE-REFUSED ECONNREFUSED" on attempt 1; accepts in its window: child 1, silent 0,`, []string{"child"}, 1},
		{"proxy-child", "proxy", prx + `: token "WORLD-PROBE-REFUSED ECONNREFUSED" on attempt 1; accepts in its window: child 1, silent 0,`, []string{"child"}, 1},
		{"raw-silent", "connect", raw + `: token "WORLD-PROBE-REFUSED ECONNREFUSED" on attempt 1; accepts in its window: child 0, silent 1,`, []string{"silent"}, 1},
		{"raw-scan", "connect", raw + `: token "WORLD-PROBE-REFUSED ECONNREFUSED" on attempt 1; accepts in its window: child 0, silent 1,`, []string{"scan"}, 1},
		{"proxy-silent", "proxy", prx + `: token "WORLD-PROBE-REFUSED ECONNREFUSED" on attempt 1; accepts in its window: child 0, silent 1,`, []string{"silent"}, 1},
		{"raw-foreign-and-child", "connect", raw + `: token "WORLD-PROBE-REFUSED ECONNREFUSED" on attempt 1; accepts in its window: child 1, silent 0,`, []string{"foreign+child"}, 1},
		{"raw-child-on-retry", "connect", raw + `: token "WORLD-PROBE-REFUSED ECONNREFUSED" on attempt 2; accepts in its window: child 1, silent 0,`, []string{"foreign", "child"}, 2},
	} {
		t.Run(row.name, func(t *testing.T) {
			rig, err := runNet6Rig(t, map[string][]string{row.op: row.plan})
			why := arm6Why(t, err)
			t.Logf("ops %q dials %q; arm 6 = %q", rig.ops, rig.dials, why)
			if !strings.Contains(why, row.want) {
				t.Fatalf("arm 6 = %q, want it to contain %q", why, row.want)
			}
			if got := rig.count(row.op); got != row.runs {
				t.Fatalf("%s runs = %d, want %d (a child or silent accept is never retried)", row.op, got, row.runs)
			}
		})
	}
}

// A dialer that hits every attempt's window exhausts the retries, and the
// leg refuses, distinctly, never as a pass.
func TestExecProbeNet6ForeignEveryAttemptRefuses(t *testing.T) {
	rig, err := runNet6Rig(t, map[string][]string{"connect": {"foreign", "foreign", "foreign"}})
	why := arm6Why(t, err)
	want := "a raw connect to World's live listener: inconclusive after 3 attempts: every attempt's window held an accept that was not the child's (foreign "
	if !strings.Contains(why, want) || rig.count("connect") != 3 {
		t.Fatalf("arm 6 = %q, connect runs %d; want %q and 3", why, rig.count("connect"), want)
	}
}

// The classifier, on every byte shape, through a real loopback pair. It
// also pins window's bound above the classifier's read deadline: a
// classifier still pending at the bound must be a stall, never a read that
// was still allowed to finish (design §3.5).
func TestExecProbeNet6DialClassifier(t *testing.T) {
	if execProbeWindowWait <= execProbeHelloWait {
		t.Fatalf("execProbeWindowWait %s must exceed the line-read deadline execProbeHelloWait %s", execProbeWindowWait, execProbeHelloWait)
	}
	ctl, child := strings.Repeat("c", 32), strings.Repeat("d", 32)
	for _, row := range []struct {
		name, send string
		hold       bool
		want       net6Class
	}{
		{"control", "WORLD-PROBE-DIAL " + ctl + "\n", false, net6Control},
		{"child", "WORLD-PROBE-DIAL " + child + "\n", false, net6Child},
		{"child-then-more", "WORLD-PROBE-DIAL " + child + "\nmore", false, net6Child},
		{"http-get", net6ForeignRequest, false, net6Foreign},
		{"wrong-nonce", "WORLD-PROBE-DIAL " + strings.Repeat("0", 32) + "\n", true, net6Foreign},
		{"preamble-then-child", "PROXY TCP4 127.0.0.1 127.0.0.1 1 2\r\nWORLD-PROBE-DIAL " + child + "\n", false, net6Foreign},
		{"connect-close", "", false, net6Silent},
		{"held-silent", "", true, net6Silent},
		{"prefix-then-close", "WORLD-PROBE-DI", false, net6Silent},
		{"prefix-held", "WORLD-PROBE-DIAL " + child[:8], true, net6Silent},
	} {
		t.Run(row.name, func(t *testing.T) {
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer ln.Close()
			go func() {
				c, err := net.Dial("tcp", ln.Addr().String())
				if err != nil {
					return
				}
				defer c.Close()
				_, _ = io.WriteString(c, row.send)
				if row.hold {
					time.Sleep(execProbeHelloWait + 200*time.Millisecond)
				}
			}()
			c, err := ln.Accept()
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()
			start := time.Now()
			if got := classifyNet6Dial(c, ctl, child); got != row.want {
				t.Fatalf("classifyNet6Dial(%q) = %d, want %d", row.send, got, row.want)
			}
			if el := time.Since(start); el > execProbeHelloWait+250*time.Millisecond {
				t.Fatalf("classification took %s, beyond its read deadline", el)
			}
		})
	}
}

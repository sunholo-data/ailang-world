package broker

// Row 242 (w-exec-probe-cross-profile-net6-hoist): every profile's arm-6
// net6 listener is bound before any profile starts, so no profile's
// reverse-leg pick can return a port a net6 listener takes later.

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// loopbackListenPorts is the 127.0.0.1 TCP ports this process listens on,
// read from its own descriptors: it sees the probe's real listeners without
// any seam on how they are made.
func loopbackListenPorts() map[int]bool {
	out := map[int]bool{}
	var lim syscall.Rlimit
	max := 4096
	if syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim) == nil && lim.Cur < uint64(max) {
		max = int(lim.Cur)
	}
	for fd := 0; fd < max; fd++ {
		sa, err := syscall.Getsockname(fd)
		if err != nil {
			continue
		}
		in4, ok := sa.(*syscall.SockaddrInet4)
		if !ok || in4.Addr != [4]byte{127, 0, 0, 1} {
			continue
		}
		// darwin has no SO_ACCEPTCONN getsockopt: a bound stream socket with
		// no peer is a listener (or a dial not yet connected).
		if ty, err := syscall.GetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_TYPE); err != nil || ty != syscall.SOCK_STREAM {
			continue
		}
		if _, err := syscall.Getpeername(fd); err != syscall.ENOTCONN {
			continue
		}
		out[in4.Port] = true
	}
	return out
}

// A protected profile's handler runs Git before its probe starts (the
// protected baseline's preflight). Here that Git blocks until the first
// reverse-leg pick, so the unprotected profile picks while the protected
// one has not reached its probe: both net6 listeners must be live at every
// pick all the same, and none may outlive the call.
func TestExecProbeNet6BoundBeforeAnyPick(t *testing.T) {
	f := newProbeFixture(t, shimSrtTree(t, "passthrough"), nil, "", []any{"sh", "-c", "echo toolchain-ok"})
	dir := canonicalTempDir(t)
	gate, gitLog, clientLog := filepath.Join(dir, "gate"), filepath.Join(dir, "git"), filepath.Join(dir, "client")
	fake, _ := protectedFakeGit(t, `case " $* " in *' config '*) echo start >> '`+gitLog+`'; i=0
while [ ! -e '`+gate+`' ] && [ $i -lt 400 ]; do sleep 0.01; i=$((i+1)); done
echo done >> '`+gitLog+`';; esac
exit 1`)
	old := worktreeGit
	worktreeGit = func() (string, error) { return fake, nil }
	t.Cleanup(func() { worktreeGit = old })
	client := filepath.Join(dir, "client.sh")
	mustWrite(t, client, "#!/bin/sh\nprintf '%s\\n' \"$*\" >> '"+clientLog+"'\nexec '"+f.node+"' \"$@\"\n")
	if err := os.Chmod(client, 0o700); err != nil {
		t.Fatal(err)
	}
	slow := mustParseProfile(t, map[string]any{"profile": ExecProfileVersion, "project": "slow",
		"path": []any{filepath.Dir(f.node), "/usr/bin", "/bin"}, "timeout_ms": 9000, "probe": []any{"sh", "-c", "exit 0"},
		"protected": []any{"tests/**"}, "commands": map[string]any{"true": map[string]any{"argv": []any{"sh", "-c", "exit 0"}}}})

	type pickSeen struct {
		live          map[int]bool
		slowInFlight  bool
		slowPreflight bool
	}
	var mu sync.Mutex
	var picks []pickSeen
	before := loopbackListenPorts()
	cfg := ExecProbeConfig{Sandbox: f.sandbox, Paths: f.paths, Client: client, reversePort: func() (int, error) {
		live := loopbackListenPorts()
		log, _ := os.ReadFile(gitLog)
		mu.Lock()
		picks = append(picks, pickSeen{live: live, slowPreflight: strings.Contains(string(log), "start"),
			slowInFlight: !strings.Contains(string(log), "done")})
		mu.Unlock()
		_ = os.WriteFile(gate, nil, 0o600)
		return freeLoopbackPort() // the production pick
	}}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	errs, err := ProbeExecProfilesEach(ctx, cfg, []*ExecProfile{f.profile, slow})
	if err != nil || len(errs) != 2 {
		t.Fatalf("ProbeExecProfilesEach = %v, %v", errs, err)
	}

	// Each profile's net6 port, from its own Run A's connect target.
	logged, _ := os.ReadFile(clientLog)
	net6 := map[int]bool{}
	for _, m := range regexp.MustCompile(`(?m) connect (\d+):[0-9a-f]{32}$`).FindAllStringSubmatch(string(logged), -1) {
		port, _ := strconv.Atoi(m[1])
		net6[port] = true
	}
	if len(net6) != 2 {
		t.Fatalf("net6 ports from the two Run A connects = %v, want two distinct ports (client log:\n%s)", net6, logged)
	}
	if len(picks) == 0 || !picks[0].slowPreflight || !picks[0].slowInFlight {
		t.Fatalf("the first pick did not land while the slow profile's pre-probe Git ran (picks %d): the test proves nothing", len(picks))
	}
	for i, pk := range picks {
		for port := range net6 {
			if !pk.live[port] {
				t.Errorf("pick %d: net6 port %d was not a live listener (slow profile in flight=%t; live %v)", i, port, pk.slowInFlight, pk.live)
			}
		}
	}
	// None outlives the call (its accept loop may take a moment to let go).
	for deadline := time.Now().Add(2 * time.Second); ; time.Sleep(20 * time.Millisecond) {
		var left []int
		for port := range loopbackListenPorts() {
			if !before[port] {
				left = append(left, port)
			}
		}
		if len(left) == 0 {
			break
		}
		if time.Now().After(deadline) {
			slices.Sort(left)
			t.Fatalf("listeners left after ProbeExecProfilesEach returned: %v (net6 %v)", left, net6)
		}
	}
	f.assertProbeCleaned(t)
}

// The bind has exactly one call site in the package's non-test files: in
// ProbeExecProfilesEach, outside any function literal, before the fan-out's
// go statement. Moving it into the goroutine, back into probe, or into any
// other function moves it after a profile's own steps.
func TestExecProbeNet6BindSiteIsBeforeTheFanOut(t *testing.T) {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the broker test package")
	}
	dir := filepath.Dir(source)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	var decls []ast.Decl
	for _, e := range entries { // every non-test file of the package: a bind moved to another file still counts
		if !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, e.Name()), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		decls = append(decls, file.Decls...)
	}
	if len(decls) < 100 {
		t.Fatalf("parsed %d declarations from %s: the scan did not reach the package", len(decls), dir)
	}
	var sites []string
	for _, decl := range decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		var firstGo token.Pos
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if g, ok := n.(*ast.GoStmt); ok && firstGo == token.NoPos {
				firstGo = g.Pos()
			}
			return true
		})
		var lit int
		var walk func(n ast.Node) bool
		walk = func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.FuncLit:
				lit++
				ast.Inspect(n.Body, walk)
				lit--
				return false
			case *ast.CallExpr:
				if id, ok := n.Fun.(*ast.Ident); ok && id.Name == "newProbeListener" {
					where := fn.Name.Name
					if lit > 0 {
						where += " (in a func literal)"
					}
					if firstGo == token.NoPos || n.Pos() > firstGo {
						where += " (not before a go statement)"
					}
					sites = append(sites, where+" at "+fset.Position(n.Pos()).String())
				}
			}
			return true
		}
		ast.Inspect(fn.Body, walk)
	}
	if len(sites) != 1 || !strings.HasPrefix(sites[0], "ProbeExecProfilesEach at ") {
		t.Fatalf("newProbeListener call sites = %q, want exactly one, in ProbeExecProfilesEach before its go statement", sites)
	}
}

package broker

// Row 140 M3 (design_docs/planned/w-workspace-exec-toolchain-effect.md §4.6
// "Startup probe (fail closed)"): before `serve` enables Workspace.Exec, each
// profile runs seven arms against ITS OWN rendered settings, with a scratch
// episode in place of a worktree. Any failed arm refuses startup by name.
// The reference implementation is M0's matrix
// (design_docs/verification/world-row140-m0/m0_matrix.py), measured on darwin
// and linux, whose pass-through control fails every confinement arm.
//
// Two rules hold for every arm:
//   - each has a control World runs itself, unsandboxed, in Go, so the arm
//     can fail and a missing client cannot fake the control;
//   - each "refused" arm demands a POSITIVE token from World's own node
//     running execProbeScript (`WORLD-PROBE-REFUSED <errno>`); a missing or
//     unexecutable client prints nothing and FAILS the arm (V38,
//     MUT-PROBE-RC-ONLY).
//
// [M0] platform tokens (V45, V46): darwin refuses with EPERM; linux masks a
// read-denied region (ENOENT) and ro-binds a denyWrite path (EROFS). Every
// write arm also requires the host bytes unchanged.
//
// Like the handler, the probe runs every subprocess through runBounded.

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// The seven arms (§4.6), in order. A refusal names them.
const (
	ExecProbeArmWriteInside  = "arm1-write-inside"
	ExecProbeArmWriteOutside = "arm2-write-outside"
	ExecProbeArmReadFence    = "arm3-read-fence"
	ExecProbeArmDefaultWrite = "arm4-srt-default-write"
	ExecProbeArmExitStatus   = "arm5-exit-status"
	ExecProbeArmNetwork      = "arm6-network"
	ExecProbeArmToolchain    = "arm7-toolchain"
)

// execProbeScript is the fixed probe client World's node runs (`node -e
// <it> -- op arg …`): one token line per op. write/read try the path;
// connect dials 127.0.0.1:<port>; proxy sends an HTTP CONNECT for it through
// srt's proxy environment (a direct dial when there is none, as curl would);
// listen binds 127.0.0.1:<port> for <ms> (arg "<port>:<ms>").
const execProbeScript = `const fs=require('fs'),net=require('net');
const say=s=>process.stdout.write(s+'\n');const code=e=>(e&&e.code)||'EUNKNOWN';
function connect(port){return new Promise(r=>{const s=net.connect(+port,'127.0.0.1');
const t=setTimeout(()=>{s.destroy();r('WORLD-PROBE-REFUSED ETIMEDOUT')},3000);
s.on('connect',()=>{clearTimeout(t);s.destroy();r('WORLD-PROBE-CONNECTED')});
s.on('error',e=>{clearTimeout(t);r('WORLD-PROBE-REFUSED '+code(e))})})}
function proxy(port){const u=process.env.HTTPS_PROXY||process.env.https_proxy||process.env.HTTP_PROXY||process.env.http_proxy;
if(!u)return connect(port);return new Promise(r=>{let p;try{p=new URL(u)}catch(e){return r('WORLD-PROBE-REFUSED EBADPROXY')}
const s=net.connect(+(p.port||80),p.hostname);const t=setTimeout(()=>{s.destroy();r('WORLD-PROBE-REFUSED ETIMEDOUT')},3000);
s.on('connect',()=>{let h='CONNECT 127.0.0.1:'+port+' HTTP/1.1\r\nHost: 127.0.0.1:'+port+'\r\n';
if(p.username)h+='Proxy-Authorization: Basic '+Buffer.from(decodeURIComponent(p.username)+':'+decodeURIComponent(p.password)).toString('base64')+'\r\n';
s.write(h+'\r\n')});let b='';
s.on('data',d=>{b+=d;const i=b.indexOf('\r\n');if(i<0)return;clearTimeout(t);s.destroy();const st=b.slice(0,i).split(' ')[1]||'';
r(st==='200'?'WORLD-PROBE-CONNECTED':'WORLD-PROBE-REFUSED HTTP-'+st)});
s.on('error',e=>{clearTimeout(t);r('WORLD-PROBE-REFUSED '+code(e))});s.on('close',()=>{clearTimeout(t);r('WORLD-PROBE-REFUSED ECLOSED')})})}
function listen(port,ms){return new Promise(r=>{const v=net.createServer(c=>c.destroy());
v.on('error',e=>r('WORLD-PROBE-REFUSED '+code(e)));
v.listen(+port,'127.0.0.1',()=>{say('WORLD-PROBE-BOUND');setTimeout(()=>{v.close();r(null)},+ms)})})}
(async()=>{const a=process.argv.slice(1);for(let i=0;i<a.length;i+=2){const op=a[i],x=a[i+1];let t;try{
if(op==='write'){try{fs.writeFileSync(x,'WORLD-PROBE');t='WORLD-PROBE-WROTE'}catch(e){t='WORLD-PROBE-REFUSED '+code(e)}}
else if(op==='read'){try{fs.readFileSync(x);t='WORLD-PROBE-READ'}catch(e){t='WORLD-PROBE-REFUSED '+code(e)}}
else if(op==='connect')t=await connect(x);else if(op==='proxy')t=await proxy(x);
else if(op==='listen'){const q=x.split(':');t=await listen(q[0],q[1])}else t='WORLD-PROBE-ERROR EBADOP'}
catch(e){t='WORLD-PROBE-ERROR '+code(e)}if(t)say(t)}})();
`

// execProbeRunBudget bounds one sandboxed probe run (the toolchain arm uses
// the profile's timeout_ms); the startup context bounds them all.
const execProbeRunBudget = 8 * time.Second

// execProbeListenMS is how long the reverse leg's child holds its listener.
const execProbeListenMS = 1000

// ExecProbeConfig is the stack the probe checks.
type ExecProbeConfig struct {
	Sandbox *ExecSandbox
	Paths   ExecHostPaths
	// MaxOutputBytes is the handler's per-stream cap; 0 = the default.
	MaxOutputBytes int64
	// Client is the probe client run inside the sandbox; "" = Sandbox.Node,
	// World's own verified node (§4.6). Tests point it at a missing path.
	Client string
}

// ExecProbeArmFailure is one failed arm and what it observed.
type ExecProbeArmFailure struct {
	Arm string
	Why string
}

// ExecProbeError refuses `serve` for a profile, naming every failed arm.
type ExecProbeError struct {
	Project string
	Failed  []ExecProbeArmFailure
}

func (e *ExecProbeError) Error() string {
	parts := make([]string, 0, len(e.Failed))
	for _, f := range e.Failed {
		parts = append(parts, f.Arm+": "+f.Why)
	}
	return fmt.Sprintf("exec startup probe failed for profile %q (%d of 7 arms): %s", e.Project, len(e.Failed), strings.Join(parts, "; "))
}

// FailedArms returns the failed arms' names, in arm order.
func (e *ExecProbeError) FailedArms() []string {
	out := make([]string, 0, len(e.Failed))
	for _, f := range e.Failed {
		out = append(out, f.Arm)
	}
	return out
}

// execProbeEnv is the per-startup scratch the probes share: the decoys and
// /tmp/claude, planted once and removed after every profile has run.
type execProbeEnv struct {
	decoys        []string
	decoyOK       []bool // each decoy read back unsandboxed (the arm 3 control)
	madeTmpClaude bool
}

// ProbeExecProfiles runs the startup probe for every profile, concurrently,
// and returns one *ExecProbeError per refused profile (joined).
func ProbeExecProfiles(ctx context.Context, cfg ExecProbeConfig, profiles []*ExecProfile) error {
	if cfg.Sandbox == nil {
		return errors.New("broker: exec probe needs a verified sandbox")
	}
	env, err := plantExecProbe(cfg.Paths)
	defer env.cleanup(cfg.Paths)
	if err != nil {
		return err
	}
	errs := make([]error, len(profiles))
	var wg sync.WaitGroup
	for i, p := range profiles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = probeExecProfile(ctx, cfg, env, p)
		}()
	}
	wg.Wait()
	return errors.Join(errs...)
}

const execProbeDecoyText = "DECOY-NOT-A-SECRET\n"

// plantExecProbe plants the three fixed decoys (§4.6 arm 3) and makes sure
// /tmp/claude exists, so a refusal there cannot be a vacuous ENOENT (V49).
func plantExecProbe(h ExecHostPaths) (*execProbeEnv, error) {
	env := &execProbeEnv{decoys: ExecProbeDecoys(h)}
	for _, d := range env.decoys {
		if err := os.MkdirAll(filepath.Dir(d), 0o700); err != nil {
			return env, fmt.Errorf("broker: exec probe: plant decoy %s: %w", d, err)
		}
		if err := os.WriteFile(d, []byte(execProbeDecoyText), 0o600); err != nil {
			return env, fmt.Errorf("broker: exec probe: plant decoy %s: %w", d, err)
		}
		got, err := os.ReadFile(d)
		env.decoyOK = append(env.decoyOK, err == nil && string(got) == execProbeDecoyText)
	}
	if _, err := os.Lstat("/tmp/claude"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("/tmp/claude", 0o777); err != nil && !errors.Is(err, os.ErrExist) {
			return env, fmt.Errorf("broker: exec probe: make /tmp/claude: %w", err)
		}
		env.madeTmpClaude = true
	}
	return env, nil
}

func (env *execProbeEnv) cleanup(h ExecHostPaths) {
	for _, d := range env.decoys {
		_ = os.Remove(d)
	}
	_ = os.RemoveAll(filepath.Join(h.StateDir, "exec-probe"))
	_ = os.RemoveAll(filepath.Join(h.WorkspaceRoot, ".exec-probe-sibling"))
	if env.madeTmpClaude {
		_ = os.Remove("/tmp/claude") // only when empty: never someone else's files
	}
}

// execRefusalCodes are the errno codes a sandbox refusal of a file write or
// read may carry on this platform [M0, V45]: ENOENT only on linux, where a
// denied region is masked (each arm's control proves the path exists).
func execRefusalCodes(goos string) []string {
	codes := []string{"EPERM", "EACCES", "EROFS"}
	if goos == "linux" {
		codes = append(codes, "ENOENT")
	}
	return codes
}

func execFileRefusal(tok string) bool {
	code, ok := strings.CutPrefix(tok, "WORLD-PROBE-REFUSED ")
	return ok && slices.Contains(execRefusalCodes(hostGOOS()), code)
}

// execProbeRun is one sandboxed run's observation.
type execProbeRun struct {
	tokens   []string
	stdout   string
	stderr   string
	exitCode int
	exited   bool
	err      error
}

func (r execProbeRun) token(i int) string {
	if i < len(r.tokens) {
		return r.tokens[i]
	}
	return ""
}

// describe is a run's evidence for a refusal message.
func (r execProbeRun) describe() string {
	status := "no exit status"
	switch {
	case r.err == nil:
		status = "exit 0"
	case r.exited:
		status = fmt.Sprintf("exit %d", r.exitCode)
	default:
		status = r.err.Error()
	}
	return fmt.Sprintf("%s; stderr %q", status, headOf(r.stderr, 240))
}

func headOf(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) > n {
		return s[:n] + "…"
	}
	return s
}

// run runs argv under the handler's launch line, through runBounded.
func (pp *execProber) run(ctx context.Context, argv []string, trampoline bool, budget time.Duration) execProbeRun {
	if dl, ok := ctx.Deadline(); ok {
		budget = min(budget, time.Until(dl))
	}
	if budget <= 0 {
		return execProbeRun{err: errors.New("the startup budget is spent")}
	}
	capture := newCaptureHeadTail(execHeadBytes, execTailBytes, pp.h.maxOutput)
	pp.h.spawns.Add(1)
	_, err := runBounded(ctx, handlerBounds{execTimeout: budget, maxOutputBytes: pp.h.maxOutput}, handlerCommand{
		path: pp.h.sandbox.Node, args: pp.h.launch(argv, trampoline), dir: pp.cwd, env: execHostEnv(), capture: capture,
	})
	r := execProbeRun{stdout: capture.stdout.snapshot().Text, stderr: capture.stderr.snapshot().Text, err: err}
	r.exitCode, r.exited = handlerExitStatus(err)
	for _, l := range strings.Split(r.stdout, "\n") {
		if strings.HasPrefix(l, "WORLD-PROBE-") {
			r.tokens = append(r.tokens, strings.TrimSpace(l))
		}
	}
	return r
}

// execProber is one profile's probe: a scratch episode and its handler.
type execProber struct {
	h      *ExecHandler
	cwd    string
	client string
	failed []ExecProbeArmFailure
}

func (pp *execProber) fail(arm, format string, args ...any) {
	pp.failed = append(pp.failed, ExecProbeArmFailure{Arm: arm, Why: fmt.Sprintf(format, args...)})
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func probeExecProfile(ctx context.Context, cfg ExecProbeConfig, env *execProbeEnv, p *ExecProfile) error {
	tag, err := randomHex(8)
	if err != nil {
		return err
	}
	id := "wep-" + tag // a scratch episode id, in the episode grammar
	worktree, sibling := filepath.Join(cfg.Paths.WorkspaceRoot, id), filepath.Join(cfg.Paths.WorkspaceRoot, id+"-sib")
	for _, d := range []string{worktree, sibling} {
		if err := os.Mkdir(d, 0o755); err != nil {
			return fmt.Errorf("broker: exec probe for profile %q: scratch episode: %w", p.Project, err)
		}
		defer func() { _ = os.RemoveAll(d) }()
	}
	h, err := NewExecHandler(ExecHandlerConfig{Profile: p, Sandbox: cfg.Sandbox, Episode: id, Worktree: worktree,
		WorkspaceRoot: cfg.Paths.WorkspaceRoot, StateDir: cfg.Paths.StateDir, OperatorHome: cfg.Paths.OperatorHome,
		MaxOutputBytes: cfg.MaxOutputBytes})
	if err != nil {
		return fmt.Errorf("broker: exec probe for profile %q: %w", p.Project, err)
	}
	defer func() {
		_ = os.Remove(h.settingsPath)
		_ = os.RemoveAll(filepath.Join(cfg.Paths.StateDir, "exec-cache", id))
	}()
	cwd := filepath.Join(worktree, p.Root)
	dirs := []string{cwd, filepath.Join(h.cache, "home"), filepath.Join(h.cache, "tmp"), filepath.Join(h.cache, "pycache"),
		filepath.Join(h.cache, "npm")}
	for _, name := range sortedKeys(p.Caches) {
		dirs = append(dirs, filepath.Join(h.cache, name)) // never seeded by the probe
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return err
		}
	}
	client := cfg.Client
	if client == "" {
		client = cfg.Sandbox.Node
	}
	pp := &execProber{h: h, cwd: cwd, client: client}
	pp.probe(ctx, env, worktree, sibling, id)
	if len(pp.failed) > 0 {
		order := []string{ExecProbeArmWriteInside, ExecProbeArmWriteOutside, ExecProbeArmReadFence, ExecProbeArmDefaultWrite,
			ExecProbeArmExitStatus, ExecProbeArmNetwork, ExecProbeArmToolchain}
		slices.SortStableFunc(pp.failed, func(a, b ExecProbeArmFailure) int {
			return slices.Index(order, a.Arm) - slices.Index(order, b.Arm)
		})
		// One entry per arm: the first failure of each.
		pp.failed = slices.CompactFunc(pp.failed, func(a, b ExecProbeArmFailure) bool { return a.Arm == b.Arm })
		return &ExecProbeError{Project: p.Project, Failed: pp.failed}
	}
	return nil
}

func fileText(p string) (string, bool) {
	b, err := os.ReadFile(p)
	return string(b), err == nil
}

func absent(p string) bool {
	_, err := os.Lstat(p)
	return errors.Is(err, os.ErrNotExist)
}

// controlWrite proves dir is writable unsandboxed (and cleans up after).
func controlWrite(p string) error {
	if err := os.WriteFile(p, []byte("control"), 0o600); err != nil {
		return err
	}
	return os.Remove(p)
}

func (pp *execProber) probe(ctx context.Context, env *execProbeEnv, worktree, sibling, id string) {
	inside := filepath.Join(worktree, "world-exec-probe-inside")
	outside := filepath.Join(sibling, "world-exec-probe-outside")
	tmpClaude := filepath.Join("/tmp/claude", "world-exec-probe-"+id)
	defer func() { _ = os.Remove(tmpClaude) }()

	// Controls, unsandboxed, by World: the sibling and /tmp/claude are
	// writable and the targets absent; each decoy reads back (planted).
	outsideCtl := controlWrite(filepath.Join(sibling, "world-exec-probe-control"))
	tmpCtl := controlWrite(tmpClaude + "-control")
	if !absent(outside) || !absent(tmpClaude) {
		outsideCtl = errors.New("a probe target already exists")
	}

	// Arm 6's live listener and its control (V37): an unsandboxed connect
	// is accepted, so the listener is proved live before any refusal counts.
	net6 := newProbeListener()
	defer net6.close()
	netCtl := net6.control()

	// Run A, with the trampoline: arms 1–4 and arm 6's raw connect.
	argv := []string{pp.client, "-e", execProbeScript, "--", "write", inside, "write", outside}
	for _, d := range env.decoys {
		argv = append(argv, "read", d)
	}
	argv = append(argv, "write", tmpClaude, "connect", strconv.Itoa(net6.port))
	a := pp.run(ctx, argv, true, execProbeRunBudget)
	time.Sleep(200 * time.Millisecond)
	acceptsA := net6.accepts.Load()

	if got, ok := fileText(inside); a.token(0) != "WORLD-PROBE-WROTE" || !ok || got != "WORLD-PROBE" {
		pp.fail(ExecProbeArmWriteInside, "a write inside the scratch worktree: token %q, read back %q (%s)", a.token(0), got, a.describe())
	}
	switch tok := a.token(1); {
	case outsideCtl != nil:
		pp.fail(ExecProbeArmWriteOutside, "control: the sibling scratch episode is not writable unsandboxed: %v", outsideCtl)
	case !execFileRefusal(tok) || !absent(outside):
		pp.fail(ExecProbeArmWriteOutside, "a write to a sibling episode: token %q (want WORLD-PROBE-REFUSED %s), on the host: present=%t (%s)",
			tok, strings.Join(execRefusalCodes(hostGOOS()), "|"), !absent(outside), a.describe())
	}
	for i, d := range env.decoys {
		switch tok := a.token(2 + i); {
		case !env.decoyOK[i]:
			pp.fail(ExecProbeArmReadFence, "control: decoy %s did not read back unsandboxed", d)
		case !execFileRefusal(tok):
			pp.fail(ExecProbeArmReadFence, "a read of decoy %s: token %q, want WORLD-PROBE-REFUSED %s (%s)",
				d, tok, strings.Join(execRefusalCodes(hostGOOS()), "|"), a.describe())
		}
	}
	switch tok := a.token(5); {
	case tmpCtl != nil:
		pp.fail(ExecProbeArmDefaultWrite, "control: /tmp/claude is not writable unsandboxed: %v", tmpCtl)
	case !execFileRefusal(tok) || !absent(tmpClaude):
		pp.fail(ExecProbeArmDefaultWrite, "a write to %s: token %q, on the host: present=%t (%s)", tmpClaude, tok, !absent(tmpClaude), a.describe())
	}

	// Run B: arm 5, a TERM self-kill through the trampoline is 143 (V11).
	b := pp.run(ctx, []string{"/bin/sh", "-c", "kill -TERM $$"}, true, execProbeRunBudget)
	if !b.exited || b.exitCode != 143 {
		pp.fail(ExecProbeArmExitStatus, "a TERM self-kill reported %s, want exit 143", b.describe())
	}

	// Run C, WITHOUT the trampoline so srt's proxy environment reaches the
	// client: arm 6's proxy leg, then the reverse leg (a listener the child
	// binds must be unreachable from the host, V46) while World dials it.
	revPort, err := freeLoopbackPort()
	if err != nil {
		pp.fail(ExecProbeArmNetwork, "reverse leg: no free port: %v", err)
		revPort = 0
	}
	poll := pollLoopback(revPort)
	c := pp.run(ctx, []string{pp.client, "-e", execProbeScript, "--", "proxy", strconv.Itoa(net6.port),
		"listen", fmt.Sprintf("%d:%d", revPort, execProbeListenMS)}, false, execProbeRunBudget)
	reached := poll()
	time.Sleep(200 * time.Millisecond)
	acceptsC := net6.accepts.Load()
	switch raw, proxy, rev := a.token(6), c.token(0), c.token(1); {
	case netCtl != nil:
		pp.fail(ExecProbeArmNetwork, "control: World's live listener did not accept an unsandboxed connect: %v", netCtl)
	case !strings.HasPrefix(raw, "WORLD-PROBE-REFUSED ") || acceptsA != 1:
		pp.fail(ExecProbeArmNetwork, "a raw connect to World's live listener: token %q, accepts %d (want a refusal and still 1; %s)",
			raw, acceptsA, a.describe())
	case !strings.HasPrefix(proxy, "WORLD-PROBE-REFUSED ") || acceptsC != 1:
		pp.fail(ExecProbeArmNetwork, "a connect through srt's proxy environment: token %q, accepts %d (want a refusal and still 1; %s)",
			proxy, acceptsC, c.describe())
	case rev != "WORLD-PROBE-BOUND" && !strings.HasPrefix(rev, "WORLD-PROBE-REFUSED "):
		pp.fail(ExecProbeArmNetwork, "reverse leg: no token from the child's listener (%q; %s)", rev, c.describe())
	case reached:
		pp.fail(ExecProbeArmNetwork, "reverse leg: the host reached the listener the child bound on 127.0.0.1:%d", revPort)
	}

	// Run D: arm 7, the profile's toolchain probe exits 0.
	probe := pp.h.profile.Probe
	d := pp.run(ctx, append([]string{probe.Argv0}, probe.Argv[1:]...), true, time.Duration(pp.h.profile.TimeoutMS)*time.Millisecond)
	if d.err != nil {
		pp.fail(ExecProbeArmToolchain, "probe %q: %s; stdout %q", strings.Join(probe.Argv, " "), d.describe(), headOf(d.stdout, 240))
	}
}

// probeListener is arm 6's World-owned live listener with an accept count.
type probeListener struct {
	ln      net.Listener
	port    int
	accepts atomic.Int64
	err     error
}

func newProbeListener() *probeListener {
	l := &probeListener{}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		l.err = err
		return l
	}
	l.ln, l.port = ln, ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			l.accepts.Add(1)
			_ = c.Close()
		}
	}()
	return l
}

// control dials the listener unsandboxed and waits for the accept count 1.
func (l *probeListener) control() error {
	if l.err != nil {
		return l.err
	}
	c, err := net.DialTimeout("tcp", l.ln.Addr().String(), 2*time.Second)
	if err != nil {
		return err
	}
	_ = c.Close()
	for deadline := time.Now().Add(2 * time.Second); time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
		if l.accepts.Load() == 1 {
			return nil
		}
	}
	return fmt.Errorf("accept count %d after one connect, want 1", l.accepts.Load())
}

func (l *probeListener) close() {
	if l.ln != nil {
		_ = l.ln.Close()
	}
}

// freeLoopbackPort is a port nothing listens on at the time of the call.
func freeLoopbackPort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, ln.Close()
}

// pollLoopback dials 127.0.0.1:port every 50 ms until stopped; the returned
// stop function reports whether any dial connected.
func pollLoopback(port int) func() bool {
	var reached atomic.Bool
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
				reached.Store(true)
				_ = c.Close()
			}
			select {
			case <-done:
				return
			case <-time.After(50 * time.Millisecond):
			}
		}
	}()
	return func() bool {
		close(done)
		<-finished
		return reached.Load()
	}
}

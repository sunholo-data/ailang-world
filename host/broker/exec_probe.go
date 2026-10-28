package broker

// Row 140 M3 (design_docs/planned/w-workspace-exec-toolchain-effect.md §4.6
// "Startup probe (fail closed)"): before `serve` enables Workspace.Exec, each
// profile runs seven legacy arms (plus protected-write when declared) against
// ITS OWN rendered settings, with a scratch
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
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

// The legacy seven arms (§4.6), then the optional protected-write arm.
const (
	ExecProbeArmWriteInside  = "arm1-write-inside"
	ExecProbeArmWriteOutside = "arm2-write-outside"
	ExecProbeArmReadFence    = "arm3-read-fence"
	ExecProbeArmDefaultWrite = "arm4-srt-default-write"
	ExecProbeArmExitStatus   = "arm5-exit-status"
	ExecProbeArmNetwork      = "arm6-network"
	ExecProbeArmToolchain    = "arm7-toolchain"
)

const ExecProbeArmProtectedWrite = "arm8-protected-write"

// execProbeScript is the fixed probe client World's node runs (`node -e
// <it> -- op arg …`): one token line per op. write/read try the path;
// connect dials 127.0.0.1:<port>; proxy sends an HTTP CONNECT for it through
// srt's proxy environment (a direct dial when there is none, as curl would);
// listen binds 127.0.0.1:<port> for <ms> (arg "<port>:<ms>:<nonce>"), and
// answers every accept with "WORLD-PROBE-NONCE <nonce>".
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
function listen(port,ms,nonce){return new Promise(r=>{const v=net.createServer(c=>{c.resume();c.end('WORLD-PROBE-NONCE '+nonce+'\n')});
v.on('error',e=>r('WORLD-PROBE-REFUSED '+code(e)));
v.listen(+port,'127.0.0.1',()=>{say('WORLD-PROBE-BOUND');setTimeout(()=>{v.close();r(null)},+ms)})})}
(async()=>{const a=process.argv.slice(1);for(let i=0;i<a.length;i+=2){const op=a[i],x=a[i+1];let t;try{
if(op==='write'){try{fs.writeFileSync(x,'WORLD-PROBE');t='WORLD-PROBE-WROTE'}catch(e){t='WORLD-PROBE-REFUSED '+code(e)}}
else if(op==='read'){try{fs.readFileSync(x);t='WORLD-PROBE-READ'}catch(e){t='WORLD-PROBE-REFUSED '+code(e)}}
else if(op==='connect')t=await connect(x);else if(op==='proxy')t=await proxy(x);
else if(op==='listen'){const q=x.split(':');t=await listen(q[0],q[1],q[2])}else t='WORLD-PROBE-ERROR EBADOP'}
catch(e){t='WORLD-PROBE-ERROR '+code(e)}if(t)say(t)}})();
`

// execProbeRunBudget bounds one sandboxed probe run (the toolchain arm uses
// the profile's timeout_ms); the startup context bounds them all.
const execProbeRunBudget = 8 * time.Second

// execProbeListenMS is how long the reverse leg's child holds its listener.
const execProbeListenMS = 1000

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

// ExecProbeConfig is the stack the probe checks.
type ExecProbeConfig struct {
	Sandbox *ExecSandbox
	Paths   ExecHostPaths
	// MaxOutputBytes is the handler's per-stream cap; 0 = the default.
	MaxOutputBytes int64
	// Client is the probe client run inside the sandbox; "" = Sandbox.Node,
	// World's own verified node (§4.6). Tests point it at a missing path.
	Client string
	// reversePort picks each reverse-leg attempt's port; nil = freeLoopbackPort.
	// Package tests steer it, as they steer Client.
	reversePort func() (int, error)
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
	total   int
}

func (e *ExecProbeError) Error() string {
	parts := make([]string, 0, len(e.Failed))
	for _, f := range e.Failed {
		parts = append(parts, f.Arm+": "+f.Why)
	}
	total := e.total
	if total == 0 {
		total = 7
	}
	return fmt.Sprintf("exec startup probe failed for profile %q (%d of %d arms): %s", e.Project, len(e.Failed), total, strings.Join(parts, "; "))
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
// /tmp/claude, planted once; the decoys are removed after every profile.
type execProbeEnv struct {
	decoys  []string
	decoyOK []bool // each decoy read back unsandboxed (the arm 3 control)
}

// ProbeExecProfiles runs the startup probe for every profile, concurrently,
// and returns one *ExecProbeError per refused profile (joined).
func ProbeExecProfiles(ctx context.Context, cfg ExecProbeConfig, profiles []*ExecProfile) error {
	errs, err := ProbeExecProfilesEach(ctx, cfg, profiles)
	if err != nil {
		return err
	}
	return errors.Join(errs...)
}

// ProbeExecProfilesEach runs the startup probe for every profile,
// concurrently, and returns each profile's own result (errs[i] is profiles[i]'s,
// nil when it passed), so a caller can leave out one profile and keep the
// rest (row 189 M2, the profiles directory). err is a failure of the shared
// probe scratch (the decoys), which no profile can pass without.
func ProbeExecProfilesEach(ctx context.Context, cfg ExecProbeConfig, profiles []*ExecProfile) (errs []error, err error) {
	if cfg.Sandbox == nil {
		return nil, errors.New("broker: exec probe needs a verified sandbox")
	}
	env, err := plantExecProbe(cfg.Paths)
	defer env.cleanup(cfg.Paths)
	if err != nil {
		return nil, err
	}
	errs = make([]error, len(profiles))
	var wg sync.WaitGroup
	for i, p := range profiles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = probeExecProfile(ctx, cfg, env, p)
		}()
	}
	wg.Wait()
	return errs, nil
}

const execProbeDecoyText = "DECOY-NOT-A-SECRET\n"

// plantExecProbe plants the three fixed decoys (§4.6 arm 3) and makes sure
// /tmp/claude exists, so a refusal there cannot be a vacuous ENOENT (V49).
// /tmp/claude is srt's own default dir (V4): World leaves it in place, as srt
// does, so concurrent probes never pull it from under each other.
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
	}
	return env, nil
}

func (env *execProbeEnv) cleanup(h ExecHostPaths) {
	for _, d := range env.decoys {
		_ = os.Remove(d)
	}
	_ = os.RemoveAll(filepath.Join(h.StateDir, "exec-probe"))
	_ = os.RemoveAll(filepath.Join(h.WorkspaceRoot, ".exec-probe-sibling"))
}

// execRefusalCodes are the errno codes a sandbox refusal of a file write or
// read may carry on this platform [M0, V45]: ENOENT only on linux, where a
// denied region is masked (each arm's control proves the path exists).
func execRefusalCodes(goos string) []string {
	codes := []string{"EPERM", "EACCES", "EROFS"}
	if goos == "linux" {
		codes = append(codes, "ENOENT", "EBUSY", "EXDEV")
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
	h        *ExecHandler
	cwd      string
	client   string
	pickPort func() (int, error)
	failed   []ExecProbeArmFailure
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
	pp := &execProber{h: h, cwd: cwd, client: client, pickPort: cfg.reversePort}
	pp.probe(ctx, env, worktree, sibling, id)
	pp.probeProtected(ctx, hostGOOS())
	if len(pp.failed) > 0 {
		order := execProbeOrder(p)
		slices.SortStableFunc(pp.failed, func(a, b ExecProbeArmFailure) int {
			return slices.Index(order, a.Arm) - slices.Index(order, b.Arm)
		})
		// One entry per arm: the first failure of each.
		pp.failed = slices.CompactFunc(pp.failed, func(a, b ExecProbeArmFailure) bool { return a.Arm == b.Arm })
		return &ExecProbeError{Project: p.Project, Failed: pp.failed, total: len(order)}
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
	// reverseLeg retries a port another host process holds (w-exec-probe-
	// reverse-leg-port-race): only the child's nonce counts as a reach.
	c, revWhy := pp.reverseLeg(ctx, net6.port)
	time.Sleep(200 * time.Millisecond)
	acceptsC := net6.accepts.Load()
	// Every leg is judged and every failing one named: one leg must never
	// hide another (the reverse leg would otherwise mask a vacuous connect).
	var legs []string
	raw, proxy := a.token(6), c.token(0)
	if netCtl != nil {
		legs = append(legs, fmt.Sprintf("control: World's live listener did not accept an unsandboxed connect: %v", netCtl))
	}
	if !strings.HasPrefix(raw, "WORLD-PROBE-REFUSED ") || acceptsA != 1 {
		legs = append(legs, fmt.Sprintf("a raw connect to World's live listener: token %q, accepts %d (want a refusal and still 1; %s)",
			raw, acceptsA, a.describe()))
	}
	if !strings.HasPrefix(proxy, "WORLD-PROBE-REFUSED ") || acceptsC != 1 {
		legs = append(legs, fmt.Sprintf("a connect through srt's proxy environment: token %q, accepts %d (want a refusal and still 1; %s)",
			proxy, acceptsC, c.describe()))
	}
	if revWhy != "" {
		legs = append(legs, revWhy)
	}
	if len(legs) > 0 {
		pp.fail(ExecProbeArmNetwork, "%s", strings.Join(legs, "; "))
	}

	// Run D: arm 7, the profile's toolchain probe exits 0.
	probe := pp.h.profile.Probe
	d := pp.run(ctx, append([]string{probe.Argv0}, probe.Argv[1:]...), true, time.Duration(pp.h.profile.TimeoutMS)*time.Millisecond)
	if d.err != nil {
		pp.fail(ExecProbeArmToolchain, "probe %q: %s; stdout %q", strings.Join(probe.Argv, " "), d.describe(), headOf(d.stdout, 240))
	}
}

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

// execProbeOrder preserves the legacy seven arms for profiles with no protection.
func execProbeOrder(p *ExecProfile) []string {
	arms := []string{ExecProbeArmWriteInside, ExecProbeArmWriteOutside, ExecProbeArmReadFence, ExecProbeArmDefaultWrite,
		ExecProbeArmExitStatus, ExecProbeArmNetwork, ExecProbeArmToolchain}
	if len(p.Protected) > 0 {
		arms = append(arms, ExecProbeArmProtectedWrite)
	}
	return arms
}

// execProbeProtectedSample constructs a witness using the grammar's own matcher.
func execProbeProtectedSample(pattern string) (string, error) {
	if _, err := parseProtectedEntry(0, pattern); err != nil {
		return "", err
	}
	if strings.HasSuffix(pattern, "/**") {
		return strings.TrimSuffix(pattern, "/**") + "/world-probe", nil
	}
	name := regexp.MustCompile(`\[[^]]*\]|[?*]`).ReplaceAllStringFunc(pattern, func(glob string) string {
		if glob == "*" || glob == "?" {
			return "x"
		}
		for c := rune(32); c <= utf8.MaxRune; c++ {
			if c == '/' {
				continue
			}
			if ok, _ := path.Match(foldProtectedKey(glob), foldProtectedKey(string(c))); ok {
				return string(c)
			}
		}
		return ""
	})
	if name == "." || name == ".." || !matchProtected([]string{pattern}, name) {
		return "", fmt.Errorf("unsupported probe specimen for %q", pattern)
	}
	return name, nil
}

// probeProtected owns only scratch worktree data. goos makes the Linux oracle
// testable on Darwin without changing the host's sandbox rendering.
func (pp *execProber) probeProtected(ctx context.Context, goos string) {
	if !slices.Contains(execProbeOrder(pp.h.profile), ExecProbeArmProtectedWrite) {
		return
	}
	if err := pp.protectedProbe(ctx, goos); err != nil {
		pp.fail(ExecProbeArmProtectedWrite, "%v", err)
	}
}

func (pp *execProber) protectedProbe(ctx context.Context, goos string) error {
	root := pp.h.worktree
	git, err := exec.LookPath("git")
	if err != nil {
		return err
	}
	var targets []string
	var basename string
	seen := map[bool]bool{}
	for _, pattern := range pp.h.profile.Protected {
		directory := strings.HasSuffix(pattern, "/**")
		if seen[directory] {
			continue
		}
		seen[directory] = true
		sample, err := execProbeProtectedSample(pattern)
		if err != nil {
			return err
		}
		if !directory {
			// Keep basename witnesses outside every protected directory.
			for i := 0; ; i++ {
				parent := fmt.Sprintf("world-probe-base-%d", i)
				dirs := []string{}
				for _, p := range pp.h.profile.Protected {
					if strings.HasSuffix(p, "/**") {
						dirs = append(dirs, p)
					}
				}
				if !matchProtected(dirs, parent) {
					sample = parent + "/" + sample
					break
				}
			}
			basename = filepath.Join(root, sample)
		}
		target := filepath.Join(root, sample)
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			return err
		}
		// Prove both creation and modification unsandboxed, then restore baseline.
		for _, text := range []string{"control", "BASE"} {
			if err := os.WriteFile(target, []byte(text), 0600); err != nil {
				return err
			}
		}
		targets = append(targets, target)
	}
	// No source-repository Git writes: this is the dedicated scratch episode.
	for _, args := range [][]string{{"init", "-q", "--template="}, {"add", "-A"},
		{"-c", "user.name=World probe", "-c", "user.email=probe@world.invalid", "-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "commit", "-qm", "probe base", "--allow-empty"}} {
		if _, err := runWorktreeGit(ctx, git, root, args...); err != nil {
			return err
		}
	}
	base, err := runWorktreeGit(ctx, git, root, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	for _, target := range targets {
		r := pp.run(ctx, []string{pp.client, "-e", execProbeScript, "--", "write", target}, true, execProbeRunBudget)
		if goos == "linux" && target == basename {
			if r.token(0) != "WORLD-PROBE-WROTE" {
				return fmt.Errorf("linux basename: missing success token (%s)", r.describe())
			}
			checked := checkProtected(ctx, git, root, strings.TrimSpace(string(base)), pp.h.profile.Protected, nil)
			if checked.err != nil || checked.intact == nil || *checked.intact {
				return fmt.Errorf("linux detection must report intact=false: intact=%v error=%v", checked.intact, checked.err)
			}
		} else {
			got, ok := fileText(target)
			if !execFileRefusal(r.token(0)) || !ok || got != "BASE" {
				return fmt.Errorf("protected modification: token %q host unchanged=%t (%s)", r.token(0), ok && got == "BASE", r.describe())
			}
		}
	}
	if goos == "darwin" && basename != "" {
		target := filepath.Join(filepath.Dir(basename), "new", filepath.Base(basename))
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			return err
		}
		if err := controlWrite(target); err != nil {
			return err
		}
		r := pp.run(ctx, []string{pp.client, "-e", execProbeScript, "--", "write", target}, true, execProbeRunBudget)
		if !execFileRefusal(r.token(0)) || !absent(target) {
			return fmt.Errorf("protected creation: token %q host absent=%t (%s)", r.token(0), absent(target), r.describe())
		}
	}
	return nil
}

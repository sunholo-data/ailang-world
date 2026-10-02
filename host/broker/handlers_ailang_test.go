package broker

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// w-software-engineering-domain M4a — the AILANG tool handlers.
//
// Two tiers. The binary-independent tests (a scripted fake tool binary and
// pure envelope tables) always run and carry the op allowlist, the `--`
// separator, the envelope classification and the rendered key set. The
// real-binary tests run only when WORLD_TOOL_AILANG_BIN names the measured
// v0.51.0 tool binary, and reproduce the design's V-rows on fixture worktrees.

const toolBinEnv = "WORLD_TOOL_AILANG_BIN"

const measuredToolVersion = "AILANG v0.51.0"

// ---------------------------------------------------------------------------
// rendered policy (§4.3) and the D4-style refusal (AC4.5, handler half)
// ---------------------------------------------------------------------------

func TestRenderEpisodePolicyHasExactlyTheDesignKeySet(t *testing.T) {
	got, err := RenderEpisodePolicy("/w/root/ep1")
	if err != nil {
		t.Fatal(err)
	}
	want := `security_mode = "restricted"
allowed_caps = ["IO", "FS"]
fs_sandbox = "/w/root/ep1"
timeout_ms = 8000
fs_deny_write = [".github/**", ".pi/**", ".claude/**", ".gitmodules", ".gitattributes", ".ailang/**"]
entry = "main"

[budgets]
FS = 1000
`
	if string(got) != want {
		t.Fatalf("rendered policy:\n%s\nwant:\n%s", got, want)
	}
	// The key set, independently of layout: a dropped or added key is named.
	var keys []string
	for _, line := range strings.Split(string(got), "\n") {
		if k, _, ok := strings.Cut(line, " = "); ok {
			keys = append(keys, k)
		}
	}
	if fmt.Sprint(keys) != "[security_mode allowed_caps fs_sandbox timeout_ms fs_deny_write entry FS]" {
		t.Fatalf("rendered keys = %v", keys)
	}
	for _, bad := range []string{"", "relative/ep1", "/w/root/../ep1", "/w/root/ep1/", `/w/"q"`, `/w/back\slash`, "/w/new\nline"} {
		if _, err := RenderEpisodePolicy(bad); err == nil {
			t.Errorf("RenderEpisodePolicy(%q) accepted", bad)
		}
	}
}

func TestCheckPolicyOutsideRootRefusesInsideAndSymlinkedIn(t *testing.T) {
	base := canonicalTempDir(t)
	root := filepath.Join(base, "root", "ep1")
	outside := filepath.Join(base, "policies")
	mustMkdir(t, root, outside)
	if err := os.Symlink(filepath.Join(root, "sub"), filepath.Join(base, "link-in")); err != nil {
		t.Fatal(err)
	}
	mustMkdir(t, filepath.Join(root, "sub"))
	for _, tc := range []struct {
		path   string
		inside bool
	}{
		{filepath.Join(outside, "ep1.toml"), false},
		{filepath.Join(base, "root", "ep2.toml"), false},
		{filepath.Join(root, "p.toml"), true},
		{filepath.Join(root, "deep", "er", "p.toml"), true},
		{root, true},
		{filepath.Join(base, "link-in", "p.toml"), true},
	} {
		err := CheckPolicyOutsideRoot(tc.path, root)
		var inside *PolicyInsideWorkspaceError
		if got := errors.As(err, &inside); got != tc.inside {
			t.Errorf("CheckPolicyOutsideRoot(%q) = %v, want inside=%v", tc.path, err, tc.inside)
		}
	}
}

// ---------------------------------------------------------------------------
// a scripted fake tool binary
// ---------------------------------------------------------------------------

// fakeTool is a /bin/sh stand-in for the tool binary. It logs each argv and
// stdin, counts non-summary dispatches, and answers per the mode file:
// summary (restricted, fs_sandbox = cwd), `{"ok":true}` for any other op,
// and for `run` the output named by mode.
type fakeTool struct {
	bin, dir string
}

func newFakeTool(t *testing.T, mode string) fakeTool {
	t.Helper()
	dir := t.TempDir()
	script := `#!/bin/sh
d=$(dirname "$0")
printf '%s\n' "$@" > "$d/argv"
mode=$(cat "$d/mode")
if [ "$1" = "policy-tool" ]; then
  req=$(cat)
  case "$req" in
    *'"summary"'*)
      sm=restricted; [ "$mode" = permissive ] && sm=permissive
      sb=$(pwd -P); [ "$mode" = elsewhere ] && sb=/
      printf '{"ok":true,"summary":{"security_mode":"%s","policy_digest":"fakedigest","fs_sandbox":"%s","cli":["check","tree"]}}' "$sm" "$sb"
      exit 0 ;;
  esac
  printf '%s' "$req" > "$d/stdin"
  echo x >> "$d/dispatches"
  echo '{"ok":true}'
  exit 0
fi
echo x >> "$d/dispatches"
case "$mode" in
  admitted) echo out; echo 'policy: {"ok":true,"decision":{"ok":true,"function":"main"}}' >&2; exit 0 ;;
  nodecision) echo 'some output'; echo 'no decision here' >&2; exit 0 ;;
esac
exit 0
`
	bin := filepath.Join(dir, "ailang")
	if err := os.WriteFile(bin, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	f := fakeTool{bin: bin, dir: dir}
	f.setMode(t, mode)
	return f
}

func (f fakeTool) setMode(t *testing.T, mode string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(f.dir, "mode"), []byte(mode), 0o600); err != nil {
		t.Fatal(err)
	}
}

func (f fakeTool) dispatches(t *testing.T) int {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.dir, "dispatches"))
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(data), "x")
}

func (f fakeTool) argv(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.dir, "argv"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n")
}

type toolLayout struct {
	base, root, sibling, policyPath, cacheDir string
}

func newToolLayout(t *testing.T) toolLayout {
	t.Helper()
	base := canonicalTempDir(t)
	l := toolLayout{
		base:       base,
		root:       filepath.Join(base, "root", "ep1"),
		sibling:    filepath.Join(base, "root", "ep2"),
		policyPath: filepath.Join(base, "policies", "ep1.toml"),
		cacheDir:   filepath.Join(base, "cache", "ep1"),
	}
	mustMkdir(t, l.root, l.sibling, filepath.Dir(l.policyPath), l.cacheDir)
	return l
}

func fakeBinRef(t *testing.T) hashref.HashRef {
	t.Helper()
	ref, err := hashref.New("sha256", strings.Repeat("ab", 32))
	if err != nil {
		t.Fatal(err)
	}
	return ref
}

func newFakeToolHandler(t *testing.T, mode string) (*AilangToolHandler, fakeTool, toolLayout) {
	t.Helper()
	tool := newFakeTool(t, mode)
	l := newToolLayout(t)
	examples := filepath.Join(l.base, "examples")
	mustMkdir(t, examples)
	h, err := NewAilangToolHandler(boundedTestContext(t), AilangToolConfig{
		Bin: tool.bin, BinRef: fakeBinRef(t), PolicyPath: l.policyPath, Root: l.root, CacheDir: l.cacheDir,
		ExamplesDir: examples,
	})
	if err != nil {
		t.Fatalf("NewAilangToolHandler: %v", err)
	}
	return h, tool, l
}

// ---------------------------------------------------------------------------
// AC4.2 + MUT-OP-ALLOWLIST (binary-independent)
// ---------------------------------------------------------------------------

func TestAilangToolOpAllowlistIsFixedPerEffect(t *testing.T) {
	h, tool, _ := newFakeToolHandler(t, "ok")
	allowed := map[string][]string{
		EffectWorkspaceRead:  {"read"},
		EffectWorkspaceWrite: {"write", "edit"},
		EffectAilangCheck:    {"ai_check"},
		EffectAilangDiscover: {"builtins_list", "examples_search"},
		EffectAilangCLI:      {"check", "tree"}, // the fake summary's cli list
	}
	universe := []string{"read", "write", "edit", "ai_check", "builtins_list", "examples_search",
		"check", "tree", "summary", "run", "fmt", "policy_check", ""}
	want, refusals := 0, 0
	for effect, ops := range allowed {
		in := map[string]bool{}
		for _, op := range ops {
			in[op] = true
		}
		for _, op := range universe {
			before := tool.dispatches(t)
			payload := fmt.Sprintf(`{"op":%q,"path":"data.txt"}`, op)
			out, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: effect, Scope: WorkspaceScope, Cost: 1}, []byte(payload))
			dispatched := tool.dispatches(t) - before
			if in[op] {
				want++
				if err != nil || dispatched != 1 {
					t.Errorf("%s op %q: err=%v dispatched=%d, want admitted", effect, op, err, dispatched)
				}
				var resp map[string]any
				if json.Unmarshal(out, &resp) != nil || resp["tool"] != fakeBinRef(t).String() || resp["policy_digest"] != "fakedigest" {
					t.Errorf("%s op %q: result %s lacks tool/policy_digest", effect, op, out)
				}
				continue
			}
			var refusal *AilangToolRefusalError
			if !errors.As(err, &refusal) || dispatched != 0 {
				t.Errorf("%s op %q: err=%v dispatched=%d, want a refusal and no dispatch", effect, op, err, dispatched)
			}
			refusals++
		}
	}
	t.Logf("op allowlist: %d admitted, %d refused", want, refusals)
	if want != 8 || refusals != 5*len(universe)-8 {
		t.Fatalf("matrix counted %d admitted / %d refused", want, refusals)
	}
}

func TestAilangToolRefusesAliasedOpKeyScopeAndRunPayloadKeys(t *testing.T) {
	h, tool, _ := newFakeToolHandler(t, "admitted")
	ctx := boundedTestContext(t)
	for _, tc := range []struct {
		effect, scope, payload string
	}{
		{EffectWorkspaceRead, WorkspaceScope, `{"op":"read","OP":"write","path":"x"}`},
		{EffectWorkspaceRead, WorkspaceScope, `{"Op":"write","path":"x"}`},
		{EffectWorkspaceRead, WorkspaceScope, `["op","read"]`},
		{EffectWorkspaceRead, WorkspaceScope, `{"op":"read"} {"op":"write"}`},
		{EffectWorkspaceRead, WorkspaceScope, `{"op":7}`},
		{EffectWorkspaceRead, "/abs/worktree", `{"op":"read","path":"x"}`},
		{EffectAilangRun, "/abs/worktree", `{"path":"x.ail"}`},
		{EffectAilangRun, WorkspaceScope, `{"path":"x.ail","flags":"--caps"}`},
		{EffectAilangRun, WorkspaceScope, `{"path":""}`},
		{EffectAilangRun, WorkspaceScope, `{"args_json":"1"}`},
		{EffectAilangRun, WorkspaceScope, `{"path":"x.ail","args_json":"{not json"}`},
	} {
		before := tool.dispatches(t)
		_, err := h.Execute(ctx, EffectRequest{Effect: tc.effect, Scope: tc.scope, Cost: 1}, []byte(tc.payload))
		var refusal *AilangToolRefusalError
		if !errors.As(err, &refusal) || tool.dispatches(t) != before {
			t.Errorf("%s %s %s: err=%v, want a refusal with no dispatch", tc.effect, tc.scope, tc.payload, err)
		}
	}
	if _, err := h.Execute(ctx, EffectRequest{Effect: EffectFSRead, Scope: WorkspaceScope}, []byte(`{}`)); err == nil {
		t.Error("an effect outside the six was executed")
	}
}

// ---------------------------------------------------------------------------
// MUT-NO-DASHDASH (binary-independent) and the AC4.3 stub arm
// ---------------------------------------------------------------------------

func TestAilangRunArgvPutsDashDashBeforeThePath(t *testing.T) {
	h, tool, l := newFakeToolHandler(t, "admitted")
	out, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"-policy.ail","args_json":"\"-x\""}`))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	want := []string{"run", "--policy", l.policyPath, "--args-json", `"-x"`, "--", "-policy.ail"}
	if got := tool.argv(t); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("run argv = %q, want %q", got, want)
	}
	var res runResult
	if err := json.Unmarshal(out, &res); err != nil || !res.Admitted || res.Stdout != "out\n" {
		t.Fatalf("run result = %s (%v)", out, err)
	}
	// Without args_json the flag is absent, not empty.
	if _, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"m.ail"}`)); err != nil {
		t.Fatal(err)
	}
	if got := tool.argv(t); fmt.Sprint(got) != fmt.Sprint([]string{"run", "--policy", l.policyPath, "--", "m.ail"}) {
		t.Fatalf("run argv without args = %q", got)
	}
}

func TestAilangRunWithoutADecisionIsRecordedFailed(t *testing.T) {
	h, tool, _ := newFakeToolHandler(t, "nodecision")
	recording := &publishRecordingStore{base: openTestStore(t)}
	session := newSession(recording, "se-stub", []Capability{{
		Effect: EffectAilangRun, Scope: WorkspaceScope, ExpiresAt: 100, Budget: 1,
	}}, Registry{EffectAilangRun: h}, Live, nil)
	_, ref, err := session.Invoke(boundedTestContext(t), EffectRequest{
		Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1, Now: 1,
	}, []byte(`{"path":"m.ail"}`))
	var failed *EffectFailedError
	var envelope *AilangRunEnvelopeError
	if !errors.As(err, &failed) || !errors.As(err, &envelope) {
		t.Fatalf("stub run error = %T %v, want *EffectFailedError wrapping *AilangRunEnvelopeError", err, err)
	}
	if tool.dispatches(t) != 1 || ref.IsZero() {
		t.Fatalf("dispatches=%d ref=%s", tool.dispatches(t), ref)
	}
	receipt := effectReceipt(t, recording, 0)
	if receipt.State != store.ReceiptResolved || receipt.EffectOutcome == nil || receipt.EffectOutcome.Status != "failed" {
		t.Fatalf("stub run receipt = %+v, want resolved/failed", receipt)
	}
}

func TestNewAilangToolHandlerRefusals(t *testing.T) {
	ctx := boundedTestContext(t)
	ref := fakeBinRef(t)
	for _, tc := range []struct {
		name string
		mode string
		edit func(l toolLayout, c *AilangToolConfig)
	}{
		{"policy inside root", "ok", func(l toolLayout, c *AilangToolConfig) { c.PolicyPath = filepath.Join(l.root, "p.toml") }},
		{"cache inside root", "ok", func(l toolLayout, c *AilangToolConfig) { c.CacheDir = filepath.Join(l.root, ".cache") }},
		{"relative bin", "ok", func(_ toolLayout, c *AilangToolConfig) { c.Bin = "ailang" }},
		{"zero bin ref", "ok", func(_ toolLayout, c *AilangToolConfig) { c.BinRef = hashref.HashRef{} }},
		{"missing root", "ok", func(l toolLayout, c *AilangToolConfig) { c.Root = filepath.Join(l.base, "nope") }},
		{"summary not restricted", "permissive", nil},
		{"summary sandbox elsewhere", "elsewhere", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tool := newFakeTool(t, tc.mode)
			l := newToolLayout(t)
			cfg := AilangToolConfig{Bin: tool.bin, BinRef: ref, PolicyPath: l.policyPath, Root: l.root, CacheDir: l.cacheDir}
			if tc.edit != nil {
				tc.edit(l, &cfg)
			}
			if _, err := NewAilangToolHandler(ctx, cfg); err == nil {
				t.Fatal("constructor accepted")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// the four run outcome shapes, as recorded on v0.51.0 (pure)
// ---------------------------------------------------------------------------

func TestComposeRunResultAcceptsExactlyTheFourShapes(t *testing.T) {
	const pl = `policy: {"ok":true,"policy_digest":"d","decision":{"ok":true,"function":"main","declared_effects":["FS","IO"]}}`
	const limit = `policy-result: {"message":"exceeded timeout_ms (1.5s); worker process group killed","reason":"timeout","stage":"execute","version":1}`
	const static = `{"file":"net.ail","decision":{"ok":false,"error_kind":"policy_violation","missing_from_policy":["Net"]}}`
	for _, tc := range []struct {
		name           string
		rc             int
		stdout, stderr string
		ok             bool
		admitted       bool
		limit          bool
	}{
		{"(a) admitted rc0", 0, "hello\n", pl + "\n", true, true, false},
		{"(a) program's own rc", 3, "bye\n", pl + "\n", true, true, false},
		{"(b) runtime refusal", 1, "", "Error: execution failed: readFile: escapes sandbox\n" + pl + "\n", true, true, false},
		{"(c) static refusal", 2, static, "", true, false, false},
		{"(d) supervisor limit", 3, "start\n", pl + "\n" + limit + "\n", true, true, true},
		{"no decision anywhere", 0, "out", "noise\n", false, false, false},
		{"flag parse (no --)", 2, "", "flag provided but not defined: -policy.ail\nUsage of run:\n", false, false, false},
		{"two policy lines", 0, "", pl + "\n" + pl + "\n", false, false, false},
		{"policy line not ok", 0, "", `policy: {"ok":false,"decision":{"ok":true}}` + "\n", false, false, false},
		{"admitted line, decision not ok", 0, "", `policy: {"ok":true,"decision":{"ok":false}}` + "\n", false, false, false},
		{"limit without rc 3", 1, "", pl + "\n" + limit + "\n", false, false, false},
		{"limit without reason", 3, "", pl + "\n" + `policy-result: {"stage":"execute"}` + "\n", false, false, false},
		{"(c) unreadable entry rc 1 (read_failed, V63)", 1, `{"file":"missing.ail","decision":{"ok":false,"error_kind":"read_failed"}}`, "", true, false, false},
		{"static decision with rc 0", 0, static, "", false, false, false},
		{"static with stderr", 2, static, "warning\n", false, false, false},
		{"static decision ok", 2, `{"decision":{"ok":true}}`, "", false, false, false},
		{"static not json", 2, "nope", "", false, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := composeRunResult(tc.rc, []byte(tc.stdout), []byte(tc.stderr))
			if !tc.ok {
				var envelope *AilangRunEnvelopeError
				if !errors.As(err, &envelope) {
					t.Fatalf("err = %v (out %s), want *AilangRunEnvelopeError", err, out)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			var res runResult
			if err := json.Unmarshal(out, &res); err != nil {
				t.Fatal(err)
			}
			if res.Admitted != tc.admitted || res.ExitCode != tc.rc || res.Stdout != tc.stdout || res.Stderr != tc.stderr {
				t.Fatalf("result = %s", out)
			}
			if gotLimit := string(res.Limit) != "null"; gotLimit != tc.limit {
				t.Fatalf("limit = %s, want present=%v", res.Limit, tc.limit)
			}
			if len(res.Decision) == 0 {
				t.Fatal("no decision")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// real tool binary (gated)
// ---------------------------------------------------------------------------

func toolBinary(t *testing.T) string {
	t.Helper()
	bin := os.Getenv(toolBinEnv)
	if bin == "" {
		t.Skipf("%s is unset: the real-binary AC4.1/AC4.3 matrix needs the %s tool binary "+
			"(e.g. $HOME/.pinned-ailang-tools/v0.51.0/ailang); skipping", toolBinEnv, measuredToolVersion)
	}
	return bin
}

// TestToolBinaryIsTheMeasuredRelease fails, rather than skips, when the gate
// variable names a binary other than the one the V-rows were measured on: a
// green matrix on an unmeasured binary would prove nothing (R-SE-8).
func TestToolBinaryIsTheMeasuredRelease(t *testing.T) {
	bin := toolBinary(t)
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("%s --version: %v", bin, err)
	}
	first, _, _ := strings.Cut(string(out), "\n")
	if strings.TrimSpace(first) != measuredToolVersion {
		t.Fatalf("%s names %q, want %q: the confinement matrix was measured on that release only",
			toolBinEnv, strings.TrimSpace(first), measuredToolVersion)
	}
}

func realBinRef(t *testing.T, bin string) hashref.HashRef {
	t.Helper()
	data, err := os.ReadFile(bin)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	ref, err := hashref.New("sha256", hex.EncodeToString(sum[:]))
	if err != nil {
		t.Fatal(err)
	}
	return ref
}

const (
	progReadit = "module readit\n\nimport std/fs (readFile)\nimport std/io (println)\n\n" +
		"export func main() -> () ! {IO, FS} {\n  println(readFile(\"data.txt\"))\n}\n"
	progArgs = "module args\n\nimport std/fs (readFile)\nimport std/io (println)\n\n" +
		"export func main(path: string) -> () ! {IO, FS} {\n  println(readFile(path))\n}\n"
	progWr = "module wr\n\nimport std/fs (writeFile)\nimport std/io (println)\n\n" +
		"export func main(path: string) -> () ! {IO, FS} {\n  writeFile(path, \"pwned\");\n  println(\"wrote\")\n}\n"
	progNet = "module net\n\nimport std/net (httpGet)\nimport std/io (println)\n\n" +
		"export func main() -> () ! {IO, Net} {\n  println(httpGet(\"http://127.0.0.1:7644/\"))\n}\n"
	progLoop = "module loop\n\nimport std/clock (sleep)\nimport std/io (println)\n\n" +
		"export func main() -> () ! {IO, Clock} {\n  println(\"start\");\n  sleep(5000)\n}\n"
)

type realFixture struct {
	toolLayout
	topology string
	bin      string
	h        *AilangToolHandler
	// protected holds the git metadata whose bytes must never change.
	protected map[string][]byte
}

func fixtureGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = []string{"PATH=/usr/bin:/bin:/opt/homebrew/bin:/usr/local/bin", "HOME=" + dir, "GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@invalid", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@invalid"}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// newRealFixture builds root/ep1 as either a `git init` repository (a .git
// directory) or a `git worktree add` worktree (a .git pointer FILE, V54), with
// a sibling episode, a symlink out, the fixture programs, the rendered policy
// outside the root and a handler bound to it.
func newRealFixture(t *testing.T, topology string) *realFixture {
	t.Helper()
	bin := toolBinary(t)
	f := &realFixture{toolLayout: newToolLayout(t), topology: topology, bin: bin, protected: map[string][]byte{}}
	switch topology {
	case "init":
		fixtureGit(t, f.root, "init", "-q", ".")
		f.protected[".git/config"] = mustRead(t, filepath.Join(f.root, ".git", "config"))
	case "worktree":
		if err := os.Remove(f.root); err != nil {
			t.Fatal(err)
		}
		main := filepath.Join(f.base, "main")
		mustMkdir(t, main)
		fixtureGit(t, main, "init", "-q", ".")
		fixtureGit(t, main, "commit", "-q", "--allow-empty", "-m", "init")
		fixtureGit(t, main, "worktree", "add", "-q", f.root)
		if info, err := os.Lstat(filepath.Join(f.root, ".git")); err != nil || !info.Mode().IsRegular() {
			t.Fatalf("worktree .git is not a pointer file: %v", err)
		}
		f.protected[".git"] = mustRead(t, filepath.Join(f.root, ".git"))
	default:
		t.Fatalf("unknown topology %q", topology)
	}
	mustWrite(t, filepath.Join(f.root, "data.txt"), "hello from ep1\n")
	mustWrite(t, filepath.Join(f.sibling, "secret.txt"), "secret\n")
	if err := os.Symlink("../ep2/secret.txt", filepath.Join(f.root, "link.txt")); err != nil {
		t.Fatal(err)
	}
	for name, src := range map[string]string{"readit.ail": progReadit, "args.ail": progArgs, "wr.ail": progWr,
		"net.ail": progNet, "loop.ail": progLoop, "-policy.ail": progReadit} {
		mustWrite(t, filepath.Join(f.root, name), src)
	}
	policy, err := RenderEpisodePolicy(f.root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.policyPath, policy, 0o600); err != nil {
		t.Fatal(err)
	}
	f.h = f.handler(t, f.policyPath)
	return f
}

func (f *realFixture) handler(t *testing.T, policyPath string) *AilangToolHandler {
	t.Helper()
	h, err := NewAilangToolHandler(boundedTestContext(t), AilangToolConfig{
		Bin: f.bin, BinRef: realBinRef(t, f.bin), PolicyPath: policyPath, Root: f.root, CacheDir: f.cacheDir,
	})
	if err != nil {
		t.Fatalf("NewAilangToolHandler on the real binary: %v", err)
	}
	return h
}

func (f *realFixture) op(t *testing.T, effect string, payload map[string]any) map[string]any {
	t.Helper()
	body, _ := json.Marshal(payload)
	out, err := f.h.Execute(boundedTestContext(t), EffectRequest{Effect: effect, Scope: WorkspaceScope, Cost: 1}, body)
	if err != nil {
		t.Fatalf("%s %s: %v", effect, body, err)
	}
	var resp map[string]any
	if err := json.Unmarshal(out, &resp); err != nil {
		t.Fatal(err)
	}
	return resp
}

func (f *realFixture) run(t *testing.T, h *AilangToolHandler, path, argsJSON string) runResult {
	t.Helper()
	payload := map[string]string{"path": path}
	if argsJSON != "" {
		payload["args_json"] = argsJSON
	}
	body, _ := json.Marshal(payload)
	out, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1}, body)
	if err != nil {
		t.Fatalf("run %s %s: %v", path, argsJSON, err)
	}
	var res runResult
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatal(err)
	}
	return res
}

func (f *realFixture) snapshot(rel string) (string, bool) {
	data, err := os.ReadFile(filepath.Join(f.root, rel))
	if err != nil {
		return "", false
	}
	return string(data), true
}

func (f *realFixture) assertProtectedUnchanged(t *testing.T) {
	t.Helper()
	for rel, want := range f.protected {
		if got := mustRead(t, filepath.Join(f.root, rel)); !bytes.Equal(got, want) {
			t.Fatalf("%s changed: %q -> %q", rel, want, got)
		}
	}
	if _, err := os.Stat(filepath.Join(f.root, ".git", "hooks", "pre-commit")); err == nil {
		t.Fatal(".git/hooks/pre-commit was created")
	}
}

// protectedTargets is the write matrix per topology: the git metadata (V36,
// V52, V54) plus every rendered fs_deny_write glob (V36, V52, V54, V57).
func (f *realFixture) protectedTargets() []string {
	git := []string{".git/config", ".GIT/config", ".git/hooks/pre-commit"}
	if f.topology == "worktree" {
		git = []string{".git", ".GIT"}
	}
	return append(git, ".github/workflows/x.yml", ".claude/settings.json", ".pi/settings.json",
		".gitmodules", ".gitattributes", ".ailang/cache/compile/manifest.json")
}

// TestAilangToolConfinementMatrix is AC4.1 on the real tool binary, for both
// shipped topologies.
func TestAilangToolConfinementMatrix(t *testing.T) {
	for _, topology := range []string{"init", "worktree"} {
		t.Run(topology, func(t *testing.T) {
			f := newRealFixture(t, topology)
			rows := 0

			// The fixture programs are checked by the tool binary itself
			// before any row relies on them (S5).
			for _, prog := range []string{"readit.ail", "args.ail", "wr.ail", "net.ail", "loop.ail"} {
				resp := f.op(t, EffectAilangCheck, map[string]any{"op": "ai_check", "path": prog})
				if resp["ok"] != true {
					t.Fatalf("fixture %s does not check on the tool binary: %v", prog, resp)
				}
			}

			// Reads: in-sandbox content; `..`, absolute and symlink-out refused.
			if resp := f.op(t, EffectWorkspaceRead, map[string]any{"op": "read", "path": "data.txt"}); resp["ok"] != true || resp["content"] != "hello from ep1\n" {
				t.Fatalf("in-sandbox read = %v", resp)
			}
			rows++
			for _, out := range []string{"../ep2/secret.txt", filepath.Join(f.sibling, "secret.txt"), "link.txt"} {
				resp := f.op(t, EffectWorkspaceRead, map[string]any{"op": "read", "path": out})
				if resp["ok"] != false || !strings.Contains(fmt.Sprint(resp["refused"]), "escapes sandbox") {
					t.Fatalf("read %s = %v, want refused escapes sandbox", out, resp)
				}
				if strings.Contains(fmt.Sprint(resp), "secret\n") {
					t.Fatalf("read %s leaked the sibling's content", out)
				}
				rows++
			}
			// The same three escapes by a RUNNING program (V34).
			for _, out := range []string{"../ep2/secret.txt", filepath.Join(f.sibling, "secret.txt"), "link.txt"} {
				arg, _ := json.Marshal(out)
				res := f.run(t, f.h, "args.ail", string(arg))
				if res.ExitCode != 1 || !res.Admitted || !strings.Contains(res.Stderr, "escapes sandbox") || strings.Contains(res.Stdout, "secret") {
					t.Fatalf("program read %s = %+v, want rc 1 escapes sandbox", out, res)
				}
				rows++
			}

			// Writes into protected paths: refused for the write op AND for a
			// running program, bytes unchanged.
			for _, target := range f.protectedTargets() {
				before, existed := f.snapshot(target)
				resp := f.op(t, EffectWorkspaceWrite, map[string]any{"op": "write", "path": target, "content": "pwned"})
				if resp["ok"] != false || fmt.Sprint(resp["refused"]) == "" {
					t.Fatalf("op write %s = %v, want refused", target, resp)
				}
				rows++
				arg, _ := json.Marshal(target)
				res := f.run(t, f.h, "wr.ail", string(arg))
				if res.ExitCode != 1 || !res.Admitted || !strings.Contains(res.Stderr, "E_FS_PROTECTED") || strings.Contains(res.Stdout, "wrote") {
					t.Fatalf("program write %s = %+v, want rc 1 E_FS_PROTECTED", target, res)
				}
				rows++
				after, exists := f.snapshot(target)
				if after != before || exists != existed || strings.Contains(after, "pwned") {
					t.Fatalf("%s changed (existed=%v exists=%v %q -> %q)", target, existed, exists, before, after)
				}
				f.assertProtectedUnchanged(t)
			}

			// In-sandbox writes succeed by both paths, and edit works.
			if resp := f.op(t, EffectWorkspaceWrite, map[string]any{"op": "write", "path": "ok.txt", "content": "fine"}); resp["ok"] != true {
				t.Fatalf("in-sandbox op write = %v", resp)
			}
			if got, _ := f.snapshot("ok.txt"); got != "fine" {
				t.Fatalf("ok.txt = %q", got)
			}
			if resp := f.op(t, EffectWorkspaceWrite, map[string]any{"op": "edit", "path": "ok.txt", "old_text": "fine", "new_text": "edited"}); resp["ok"] != true {
				t.Fatalf("in-sandbox edit = %v", resp)
			}
			if got, _ := f.snapshot("ok.txt"); got != "edited" {
				t.Fatalf("ok.txt after edit = %q", got)
			}
			if res := f.run(t, f.h, "wr.ail", `"ok2.txt"`); res.ExitCode != 0 || !res.Admitted || res.Stdout != "wrote\n" {
				t.Fatalf("in-sandbox program write = %+v", res)
			}
			if got, _ := f.snapshot("ok2.txt"); got != "pwned" {
				t.Fatalf("ok2.txt = %q", got)
			}
			rows += 3
			f.assertProtectedUnchanged(t)
			t.Logf("%s: %d confinement rows green", topology, rows)
		})
	}
}

// TestAilangToolReadEffectRefusesWriteOp is AC4.2 on the real binary.
func TestAilangToolReadEffectRefusesWriteOp(t *testing.T) {
	f := newRealFixture(t, "init")
	_, err := f.h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectWorkspaceRead, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"op":"write","path":"data.txt","content":"overwritten"}`))
	var refusal *AilangToolRefusalError
	if !errors.As(err, &refusal) {
		t.Fatalf("Workspace.Read op write: err = %v, want *AilangToolRefusalError", err)
	}
	if got, _ := f.snapshot("data.txt"); got != "hello from ep1\n" {
		t.Fatalf("data.txt changed to %q", got)
	}
}

// TestRenderedPolicyIsHonouredByTheToolBinary: the binary reads the rendered
// file as the §4.3 policy (restricted, IO/FS, the root, 8000 ms).
func TestRenderedPolicyIsHonouredByTheToolBinary(t *testing.T) {
	f := newRealFixture(t, "init")
	out, err := f.h.policyTool(boundedTestContext(t), []byte(`{"op":"summary"}`))
	if err != nil {
		t.Fatal(err)
	}
	var resp struct {
		Summary struct {
			SecurityMode string   `json:"security_mode"`
			Caps         []string `json:"caps"`
			FSSandbox    string   `json:"fs_sandbox"`
			TimeoutMS    int      `json:"timeout_ms"`
			PolicyDigest string   `json:"policy_digest"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(out, &resp); err != nil {
		t.Fatal(err)
	}
	s := resp.Summary
	sum := sha256.Sum256(mustRead(t, f.policyPath))
	if s.SecurityMode != "restricted" || fmt.Sprint(s.Caps) != "[FS IO]" || s.FSSandbox != f.root || s.TimeoutMS != 8000 ||
		s.PolicyDigest != hex.EncodeToString(sum[:]) {
		t.Fatalf("summary = %+v", s)
	}
}

// TestAilangRunOutcomeShapesOnTheToolBinary is AC4.3 on the real binary: one
// case per observed shape, plus the hyphen-leading path after `--` (V58).
func TestAilangRunOutcomeShapesOnTheToolBinary(t *testing.T) {
	f := newRealFixture(t, "worktree")

	// (a) admitted, with --args-json (V39).
	res := f.run(t, f.h, "args.ail", `"data.txt"`)
	if !res.Admitted || res.ExitCode != 0 || res.Stdout != "hello from ep1\n\n" || string(res.Limit) != "null" {
		t.Fatalf("(a) = %+v", res)
	}
	// (b) runtime escape (V34).
	res = f.run(t, f.h, "args.ail", `"../ep2/secret.txt"`)
	if !res.Admitted || res.ExitCode != 1 || !strings.Contains(res.Stderr, "escapes sandbox") {
		t.Fatalf("(b) = %+v", res)
	}
	// (c) a Net-declaring program is statically refused (V49).
	res = f.run(t, f.h, "net.ail", "")
	var decision struct {
		OK      bool     `json:"ok"`
		Missing []string `json:"missing_from_policy"`
	}
	if err := json.Unmarshal(res.Decision, &decision); err != nil {
		t.Fatal(err)
	}
	if res.Admitted || res.ExitCode != 2 || decision.OK || fmt.Sprint(decision.Missing) != "[Net]" || res.Stderr != "" {
		t.Fatalf("(c) = %+v decision %s", res, res.Decision)
	}
	// (c) with rc 1 (V63): a missing entry file is refused before execution
	// with the decision on stdout; the agent sees error_kind read_failed.
	res = f.run(t, f.h, "missing.ail", "")
	var readFailed struct {
		OK        bool   `json:"ok"`
		ErrorKind string `json:"error_kind"`
	}
	if err := json.Unmarshal(res.Decision, &readFailed); err != nil {
		t.Fatalf("(c) rc1 decision %s: %v", res.Decision, err)
	}
	if res.Admitted || res.ExitCode != 1 || readFailed.OK || readFailed.ErrorKind != "read_failed" || res.Stderr != "" {
		t.Fatalf("(c) missing entry = %+v decision %s", res, res.Decision)
	}
	// (d) the supervisor's own timeout (V61) under a short Clock-granting
	// policy: rc 3 and limit.reason == "timeout", well inside the Go cap.
	short := filepath.Join(filepath.Dir(f.policyPath), "short.toml")
	mustWrite(t, short, fmt.Sprintf("security_mode = \"restricted\"\nallowed_caps = [\"IO\", \"FS\", \"Clock\"]\n"+
		"fs_sandbox = %q\ntimeout_ms = 1500\nentry = \"main\"\n", f.root))
	start := time.Now()
	res = f.run(t, f.handler(t, short), "loop.ail", "")
	var lim struct {
		Reason string `json:"reason"`
		Stage  string `json:"stage"`
	}
	if err := json.Unmarshal(res.Limit, &lim); err != nil {
		t.Fatalf("(d) limit %s: %v", res.Limit, err)
	}
	if !res.Admitted || res.ExitCode != 3 || lim.Reason != "timeout" || res.Stdout != "start\n" || time.Since(start) > 8*time.Second {
		t.Fatalf("(d) = %+v limit %+v after %s", res, lim, time.Since(start))
	}
	// V58: a hyphen-leading path passed after `--` is a FILE: the binary
	// reads it and refuses its module name (MOD010), never a flag parse.
	res = f.run(t, f.h, "-policy.ail", "")
	if res.Admitted || res.ExitCode != 2 || !strings.Contains(string(res.Decision), "-policy") ||
		strings.Contains(res.Stderr+res.Stdout, "flag provided but not defined") {
		t.Fatalf("hyphen path = %+v decision %s", res, res.Decision)
	}
	f.assertProtectedUnchanged(t)
}

// ---------------------------------------------------------------------------
// AC10 census drivers for the tool handler's two branches
// ---------------------------------------------------------------------------

func driveAilangToolPolicyTool(t *testing.T, probe string) {
	t.Helper()
	l := newToolLayout(t)
	// The probe is not a tool binary, so the summary read fails; the
	// policy-tool launch, with the handler's child env, is the measurement.
	_, _ = NewAilangToolHandler(boundedTestContext(t), AilangToolConfig{
		Bin: probe, BinRef: fakeBinRef(t), PolicyPath: l.policyPath, Root: l.root, CacheDir: l.cacheDir,
	})
}

func driveAilangToolRun(t *testing.T, probe string) {
	t.Helper()
	l := newToolLayout(t)
	h := &AilangToolHandler{bin: probe, binRef: fakeBinRef(t), policyPath: l.policyPath, root: l.root,
		cacheDir: l.cacheDir, bounds: handlerBounds{execTimeout: 20 * time.Second}.normalized()}
	_, _ = h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"m.ail"}`))
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func canonicalTempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func mustMkdir(t *testing.T, dirs ...string) {
	t.Helper()
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// TestAilangToolExamplesSearchNeedsACorpus (row 134 break-3 fix): with no
// ExamplesDir, examples_search is answered with NoExamplesCorpusRefusal in
// policy-tool's refusal shape and never dispatched; builtins_list on the same
// handler still dispatches. An ExamplesDir inside the root is refused at
// construction.
func TestAilangToolExamplesSearchNeedsACorpus(t *testing.T) {
	tool := newFakeTool(t, "ok")
	l := newToolLayout(t)
	h, err := NewAilangToolHandler(boundedTestContext(t), AilangToolConfig{
		Bin: tool.bin, BinRef: fakeBinRef(t), PolicyPath: l.policyPath, Root: l.root, CacheDir: l.cacheDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := EffectRequest{Effect: EffectAilangDiscover, Scope: WorkspaceScope, Cost: 1}
	before := tool.dispatches(t)
	out, err := h.Execute(boundedTestContext(t), req, []byte(`{"op":"examples_search","query":"foldl"}`))
	var resp map[string]any
	if err != nil || json.Unmarshal(out, &resp) != nil || resp["ok"] != false || resp["refused"] != NoExamplesCorpusRefusal ||
		resp["tool"] != fakeBinRef(t).String() || resp["policy_digest"] != "fakedigest" {
		t.Fatalf("examples_search without a corpus = %s %v, want the no-corpus refusal with provenance", out, err)
	}
	if got := tool.dispatches(t) - before; got != 0 {
		t.Fatalf("examples_search without a corpus dispatched %d times, want 0", got)
	}
	if _, err := h.Execute(boundedTestContext(t), req, []byte(`{"op":"builtins_list"}`)); err != nil || tool.dispatches(t)-before != 1 {
		t.Fatalf("builtins_list: err=%v, want one dispatch", err)
	}
	_, err = NewAilangToolHandler(boundedTestContext(t), AilangToolConfig{
		Bin: tool.bin, BinRef: fakeBinRef(t), PolicyPath: l.policyPath, Root: l.root, CacheDir: l.cacheDir,
		ExamplesDir: filepath.Join(l.root, "examples"),
	})
	var inside *PolicyInsideWorkspaceError
	if !errors.As(err, &inside) {
		t.Fatalf("ExamplesDir inside the root: err = %v, want *PolicyInsideWorkspaceError", err)
	}
}

package broker

// Row 135 M2 (design_docs/planned/w-ailang-run-stdin-argv-caps.md §4.3–§4.5,
// §6 AC2.1–AC2.8, §7 mutants): ailang-run's stdin / argv / caps in the
// handler and the per-cap-set policy variants. As in row 134, two tiers: the
// binary-independent tests (the rendered key sets, the operator table, the
// refusals with an exec counter, the argv/stdin plumbing on the scripted fake
// tool) always run; the real-binary tests run when WORLD_TOOL_AILANG_BIN names
// the measured v0.52.1 tool binary and reproduce V33–V39 through the handler.

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// AC2.1 — the exact key set per variant class (binary-independent)
// ---------------------------------------------------------------------------

// policyLines parses a rendered policy into key -> raw value (the [budgets]
// table's FS as "budgets.FS"), failing on a repeated key.
func policyLines(t *testing.T, policy []byte) (map[string]string, []string) {
	t.Helper()
	keys := map[string]string{}
	var order []string
	table := ""
	for _, line := range strings.Split(string(policy), "\n") {
		if strings.HasPrefix(line, "[") {
			table = strings.Trim(line, "[]") + "."
			continue
		}
		k, v, ok := strings.Cut(line, " = ")
		if !ok {
			continue
		}
		k = table + k
		if _, dup := keys[k]; dup {
			t.Fatalf("key %s rendered twice", k)
		}
		keys[k] = v
		order = append(order, k)
	}
	return keys, order
}

var testRunOp = RunCapsConfig{
	Allow:    []string{RunCapDeclassify, RunCapEnv, RunCapNet},
	NetAllow: []string{"127.0.0.1:17655"}, NetAllowHTTP: true,
}

func TestRenderRunVariantKeySetsPerClass(t *testing.T) {
	root := "/w/root/ep1"
	base, err := RenderEpisodePolicy(root)
	if err != nil {
		t.Fatal(err)
	}
	baseKeys, _ := policyLines(t, base)
	baseDeny := baseKeys["fs_deny_write"]
	const limits = "max_source_bytes=1048576 max_module_graph_bytes=16777216 max_output_bytes=8388608 max_fs_transfer_bytes=8388608"
	for _, tc := range []struct {
		caps   []string
		mode   string
		allow  string
		fs     string
		net    string // "" = no net_allow key
		http   bool
		git    bool
		limits bool
	}{
		{[]string{"IO"}, "restricted", `["FS", "IO"]`, "0", "", false, false, false},
		{[]string{"Declassify", "IO"}, "restricted", `["Declassify", "FS", "IO"]`, "0", "", false, false, false},
		{[]string{"Declassify", "FS", "IO"}, "restricted", `["Declassify", "FS", "IO"]`, "1000", "", false, false, false},
		{[]string{"IO", "Net"}, "restricted", `["FS", "IO", "Net"]`, "0", `["127.0.0.1:17655"]`, true, false, false},
		{[]string{"FS", "IO", "Net"}, "restricted", `["FS", "IO", "Net"]`, "1000", `["127.0.0.1:17655"]`, true, false, false},
		{[]string{"Declassify", "IO", "Net"}, "restricted", `["Declassify", "FS", "IO", "Net"]`, "0", `["127.0.0.1:17655"]`, true, false, false},
		{[]string{"Env", "FS", "IO"}, "trusted_host", `["Env", "FS", "IO"]`, "1000", "", false, true, true},
		{[]string{"Env", "IO"}, "trusted_host", `["Env", "FS", "IO"]`, "0", "", false, true, true},
	} {
		t.Run(strings.Join(tc.caps, "+"), func(t *testing.T) {
			got, err := RenderRunVariant(root, tc.caps, testRunOp)
			if err != nil {
				t.Fatal(err)
			}
			keys, order := policyLines(t, got)
			want := []string{"security_mode", "allowed_caps", "fs_sandbox", "timeout_ms", "fs_deny_write", "entry"}
			if tc.net != "" {
				want = append(want, "net_allow")
				if tc.http {
					want = append(want, "net_allow_http")
				}
			}
			if tc.limits {
				want = append(want, "max_source_bytes", "max_module_graph_bytes", "max_output_bytes", "max_fs_transfer_bytes")
			}
			want = append(want, "budgets.FS")
			if fmt.Sprint(order) != fmt.Sprint(want) {
				t.Fatalf("keys %v, want %v\n%s", order, want, got)
			}
			if keys["security_mode"] != `"`+tc.mode+`"` || keys["allowed_caps"] != tc.allow || keys["budgets.FS"] != tc.fs ||
				keys["fs_sandbox"] != `"`+root+`"` || keys["timeout_ms"] != "8000" || keys["entry"] != `"main"` {
				t.Fatalf("variant values:\n%s", got)
			}
			if tc.net != "" && (keys["net_allow"] != tc.net || keys["net_allow_http"] != "true") {
				t.Fatalf("net keys: %q %q", keys["net_allow"], keys["net_allow_http"])
			}
			// No bare-host net_allow, ever (V35: it would open every port).
			for _, entry := range strings.Split(strings.Trim(keys["net_allow"], "[]"), ", ") {
				entry = strings.Trim(entry, `"`)
				if entry != "" && ValidateRunNetAllow(entry) != nil {
					t.Fatalf("net_allow entry %q is not a port-qualified loopback literal", entry)
				}
			}
			wantDeny := baseDeny
			if tc.git {
				wantDeny = `[".git/**", ` + strings.TrimPrefix(baseDeny, "[")
			}
			if keys["fs_deny_write"] != wantDeny || strings.Count(keys["fs_deny_write"], `".git/**"`) != btoi(tc.git) {
				t.Fatalf("fs_deny_write is not the base list (plus one leading .git/** for trusted_host = %v); starts %.80s",
					tc.git, keys["fs_deny_write"])
			}
			if tc.limits {
				var got []string
				for _, k := range []string{"max_source_bytes", "max_module_graph_bytes", "max_output_bytes", "max_fs_transfer_bytes"} {
					got = append(got, k+"="+keys[k])
				}
				if strings.Join(got, " ") != limits {
					t.Fatalf("limits %v, want %s", got, limits)
				}
			}
		})
	}
	// {FS, IO} is row 134's base policy, byte for byte.
	if got, err := RenderRunVariant(root, []string{"FS", "IO"}, testRunOp); err != nil || string(got) != string(base) {
		t.Fatalf("RenderRunVariant(FS, IO) is not the base policy: %v", err)
	}
	for _, bad := range [][]string{nil, {}, {"IO", "FS"}, {"IO", "IO"}, {"Env", "IO", "Net"}, {"IO", "Process"}} {
		if _, err := RenderRunVariant(root, bad, testRunOp); err == nil {
			t.Errorf("RenderRunVariant(%v) accepted a non-canonical or refused cap set", bad)
		}
	}
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ---------------------------------------------------------------------------
// AC3.1 (handler half) — the operator's startup refusal table
// ---------------------------------------------------------------------------

func TestRunCapsConfigValidate(t *testing.T) {
	ok := []RunCapsConfig{
		{},
		{Allow: []string{"Declassify"}},
		{Allow: []string{"Env", "Declassify"}},
		{Allow: []string{"Net"}, NetAllow: []string{"127.0.0.1:7655"}},
		{Allow: []string{"Net"}, NetAllow: []string{"127.0.0.1:7655", "[::1]:7655"}, NetAllowHTTP: true},
	}
	for _, c := range ok {
		if err := c.Validate(); err != nil {
			t.Errorf("%+v refused: %v", c, err)
		}
	}
	for _, tc := range []struct {
		c    RunCapsConfig
		want string
	}{
		{RunCapsConfig{Allow: []string{"Process"}}, `names "Process"`},
		{RunCapsConfig{Allow: []string{"IO"}}, `names "IO"`},
		{RunCapsConfig{Allow: []string{"env"}}, `names "env"`},
		{RunCapsConfig{Allow: []string{"Env", "Env"}}, "twice"},
		{RunCapsConfig{Allow: []string{"Net"}}, "needs at least one --run-net-allow"},
		{RunCapsConfig{NetAllow: []string{"127.0.0.1:7655"}}, "does not name Net"},
		{RunCapsConfig{Allow: []string{"Env"}, NetAllowHTTP: true}, "does not name Net"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"127.0.0.1"}}, "is not HOST:PORT"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"localhost:7655"}}, "names a host by name"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"example.com:80"}}, "names a host by name"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"10.0.0.1:80"}}, "not a loopback address"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"169.254.169.254:80"}}, "not a loopback address"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"0.0.0.0:80"}}, "not a loopback address"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"8.8.8.8:53"}}, "not a loopback address"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"127.0.0.1:0"}}, "invalid port"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"127.0.0.1:65536"}}, "invalid port"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"127.0.0.1:07655"}}, "invalid port"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"127.0.0.1:http"}}, "invalid port"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"[0:0::1]:7655"}}, "canonical form"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"[::ffff:127.0.0.1]:7655"}}, "canonical form"},
		{RunCapsConfig{Allow: []string{"Net"}, NetAllow: []string{"127.0.0.1:7655", "127.0.0.1:7655"}}, "twice"},
	} {
		err := tc.c.Validate()
		var cfgErr *RunCapsConfigError
		if !errors.As(err, &cfgErr) || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%+v: err = %v, want a RunCapsConfigError containing %q", tc.c, err, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// AC2.3 — every refusal happens before any exec (fake tool, exec counter)
// ---------------------------------------------------------------------------

func newFakeRunHandler(t *testing.T, op RunCapsConfig) (*AilangToolHandler, fakeTool) {
	t.Helper()
	tool := newFakeTool(t, "admitted")
	l := newToolLayout(t)
	h, err := NewAilangToolHandler(boundedTestContext(t), AilangToolConfig{
		Bin: tool.bin, BinRef: fakeBinRef(t), PolicyPath: l.policyPath, Root: l.root, CacheDir: l.cacheDir, RunCaps: op,
	})
	if err != nil {
		t.Fatalf("NewAilangToolHandler: %v", err)
	}
	return h, tool
}

func TestAilangRunRefusalsBeforeAnyExec(t *testing.T) {
	big := strings.Repeat("x", 65537)
	for _, tc := range []struct {
		name    string
		op      RunCapsConfig
		effect  string
		payload string
		want    string
	}{
		{"Env not allowlisted", RunCapsConfig{}, EffectAilangRunEnv, `{"path":"s.ail","caps":["Env","FS","IO"]}`, `"Env" is not enabled`},
		{"Net not allowlisted", RunCapsConfig{Allow: []string{"Env"}}, EffectAilangRunNet, `{"path":"s.ail","caps":["IO","Net"]}`, `"Net" is not enabled`},
		{"Declassify not allowlisted", RunCapsConfig{Allow: []string{"Env", "Net"}, NetAllow: []string{"127.0.0.1:7655"}},
			EffectAilangRun, `{"path":"s.ail","caps":["Declassify","IO"]}`, `"Declassify" is not enabled`},
		{"Run with Env caps", testRunOp, EffectAilangRun, `{"path":"s.ail","caps":["Env","IO"]}`, "run as Ailang.RunEnv, not Ailang.Run"},
		{"Run with Net caps", testRunOp, EffectAilangRun, `{"path":"s.ail","caps":["IO","Net"]}`, "run as Ailang.RunNet, not Ailang.Run"},
		{"RunEnv without caps", testRunOp, EffectAilangRunEnv, `{"path":"s.ail"}`, "run as Ailang.Run, not Ailang.RunEnv"},
		{"RunNet with Declassify", testRunOp, EffectAilangRunNet, `{"path":"s.ail","caps":["Declassify","IO"]}`, "not Ailang.RunNet"},
		{"Env with Net", testRunOp, EffectAilangRunEnv, `{"path":"s.ail","caps":["Env","IO","Net"]}`, "both Env and Net"},
		{"unknown cap", testRunOp, EffectAilangRun, `{"path":"s.ail","caps":["IO","Process"]}`, `unknown capability "Process"`},
		{"duplicate cap", testRunOp, EffectAilangRun, `{"path":"s.ail","caps":["IO","IO"]}`, "not sorted and distinct"},
		{"empty caps", testRunOp, EffectAilangRun, `{"path":"s.ail","caps":[]}`, "names no capability"},
		{"caps as a string", testRunOp, EffectAilangRun, `{"path":"s.ail","caps":"IO"}`, "caps must be an array"},
		{"unknown key", testRunOp, EffectAilangRun, `{"path":"s.ail","env":{"A":"B"}}`, `unknown payload key "env"`},
		{"aliased key Stdin", testRunOp, EffectAilangRun, `{"path":"s.ail","Stdin":"x"}`, `"Stdin" aliases "stdin"`},
		{"aliased key CAPS", testRunOp, EffectAilangRun, `{"path":"s.ail","CAPS":["IO"]}`, `"CAPS" aliases "caps"`},
		{"aliased key ARGV", testRunOp, EffectAilangRun, `{"path":"s.ail","ARGV":["x"]}`, `"ARGV" aliases "argv"`},
		{"NUL in argv", testRunOp, EffectAilangRun, `{"path":"s.ail","argv":["a\u0000b"]}`, "NUL"},
		{"33 argv items", testRunOp, EffectAilangRun, `{"path":"s.ail","argv":[` + strings.TrimSuffix(strings.Repeat(`"a",`, 33), ",") + `]}`, "33 items"},
		{"1025-byte argv item", testRunOp, EffectAilangRun, `{"path":"s.ail","argv":["` + strings.Repeat("a", 1025) + `"]}`, "1025 bytes"},
		{"argv of numbers", testRunOp, EffectAilangRun, `{"path":"s.ail","argv":[1]}`, "argv must be an array of strings"},
		{"oversize stdin", testRunOp, EffectAilangRun, `{"path":"s.ail","stdin":"` + big + `"}`, "65537 bytes"},
		{"stdin not a string", testRunOp, EffectAilangRun, `{"path":"s.ail","stdin":["x"]}`, "stdin must be a string"},
		{"args_json -", testRunOp, EffectAilangRun, `{"path":"s.ail","args_json":"-"}`, "args_json must be a string holding JSON"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h, tool := newFakeRunHandler(t, tc.op)
			_, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: tc.effect, Scope: WorkspaceScope, Cost: 1}, []byte(tc.payload))
			var refusal *AilangToolRefusalError
			if !errors.As(err, &refusal) || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want a refusal containing %q", err, tc.want)
			}
			if n := tool.dispatches(t); n != 0 {
				t.Fatalf("refused run dispatched %d times, want 0", n)
			}
		})
	}
	// Net with no --run-net-allow can only arise past Validate; the handler
	// still refuses it before any exec (defence in depth).
	h, tool := newFakeRunHandler(t, RunCapsConfig{})
	h.runCaps = RunCapsConfig{Allow: []string{RunCapNet}}
	_, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRunNet, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"s.ail","caps":["IO","Net"]}`))
	var refusal *AilangToolRefusalError
	if !errors.As(err, &refusal) || !strings.Contains(err.Error(), "--run-net-allow") || tool.dispatches(t) != 0 {
		t.Fatalf("Net without --run-net-allow: err=%v dispatches=%d", err, tool.dispatches(t))
	}
	// The constructor refuses an invalid allowlist outright.
	l := newToolLayout(t)
	if _, err := NewAilangToolHandler(boundedTestContext(t), AilangToolConfig{
		Bin: tool.bin, BinRef: fakeBinRef(t), PolicyPath: l.policyPath, Root: l.root, CacheDir: l.cacheDir,
		RunCaps: RunCapsConfig{Allow: []string{RunCapNet}},
	}); err == nil {
		t.Fatal("constructor accepted Net without --run-net-allow")
	}
}

// TestAilangRunVariantUnverifiedIsRefused is MUT-VARIANT-UNVERIFIED's
// binary-independent killer: the fake tool's summary reports [FS IO] for any
// policy, so a Declassify variant's summary disagrees with what was rendered
// and the run never executes.
func TestAilangRunVariantUnverifiedIsRefused(t *testing.T) {
	h, tool := newFakeRunHandler(t, testRunOp)
	_, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"pi.ail","caps":["Declassify","IO"]}`))
	var unverified *RunVariantUnverifiedError
	if !errors.As(err, &unverified) || tool.dispatches(t) != 0 {
		t.Fatalf("err = %v dispatches=%d, want *RunVariantUnverifiedError and no run", err, tool.dispatches(t))
	}
}

// TestAilangRunArgvAndStdinPlumbing (MUT-NO-INNER-DASHDASH, MUT-STDIN-*,
// binary-independent): argv follows an inner `--` after the path, verbatim;
// stdin is the payload's, byte for byte; absent stdin is no stdin; the
// result names the policy it ran under.
func TestAilangRunArgvAndStdinPlumbing(t *testing.T) {
	h, tool := newFakeRunHandler(t, testRunOp)
	out, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"-p.ail","args_json":"1","stdin":"1\n2\n","argv":["--caps=IO,Net","-","--","a b"]}`))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"run", "--policy", h.policyPath, "--args-json", "1", "--", "-p.ail", "--", "--caps=IO,Net", "-", "--", "a b"}
	if got := tool.argv(t); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("argv = %q, want %q", got, want)
	}
	if got := string(mustRead(t, filepath.Join(tool.dir, "runstdin"))); got != "1\n2\n" {
		t.Fatalf("child stdin = %q", got)
	}
	var res runResult
	if err := json.Unmarshal(out, &res); err != nil || res.Policy == nil || res.Policy.Digest != "fakedigest" ||
		res.Policy.SecurityMode != "restricted" || fmt.Sprint(res.Policy.Caps) != "[FS IO]" {
		t.Fatalf("result = %s (%v), want the base policy block", out, err)
	}
	if _, err := h.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"m.ail","argv":[]}`)); err != nil {
		t.Fatal(err)
	}
	if got := tool.argv(t); fmt.Sprint(got) != fmt.Sprint([]string{"run", "--policy", h.policyPath, "--", "m.ail"}) {
		t.Fatalf("argv for an empty argv = %q", got)
	}
	if got := string(mustRead(t, filepath.Join(tool.dir, "runstdin"))); got != "" {
		t.Fatalf("absent stdin reached the child as %q", got)
	}
}

// ---------------------------------------------------------------------------
// real tool binary (gated): AC2.1 (tamper), AC2.2, AC2.4–AC2.8
// ---------------------------------------------------------------------------

const (
	progStdin = "module stdin\nimport std/string (stringToInt)\nimport std/io (println, readLine)\n\n" +
		"export func main() -> () ! {IO} {\n  let line = readLine();\n  if line == \"\" then ()\n  else {\n" +
		"    match stringToInt(line) {\n      Some(n) => println(show(n * 2)),\n      None => ()\n    };\n    main()\n    }\n}\n"
	progSum = "module sum\nimport std/env (getArgs)\nimport std/fs (readFile)\nimport std/io (println)\n" +
		"import std/string (split, stringToInt)\nimport std/list (foldl)\nimport std/option (Some, None)\n\n" +
		"export func main() -> () ! {IO, FS, Env} {\n  let args = getArgs();\n  let filename = match args {\n" +
		"    f :: _ => f,\n    _ => \"numbers.txt\"\n  };\n  let content = readFile(filename);\n" +
		"  let lines = split(content, \"\\n\");\n  let sum = foldl(\\acc line. acc + match stringToInt(line) {\n" +
		"    Some(n) => n,\n    None => 0\n  }, 0, lines);\n  println(show(sum))\n}\n"
	progPostNet = "module postnet\n\nimport std/net (httpRequest)\nimport std/io (println)\n\n" +
		"export func main(url: string) -> () ! {IO, Net} {\n" +
		"  let headers = [{ name: \"X-Test-Header\", value: \"value123\" }, { name: \"Content-Type\", value: \"application/json\" }];\n" +
		"  match httpRequest(\"POST\", url, headers, \"{\\\"message\\\":\\\"Hello from AILANG\\\",\\\"count\\\":42}\") {\n" +
		"    Ok(resp) => println(show(resp.status)),\n    Err(e) => println(\"ERR ${show(e)}\")\n  }\n}\n"
	progDecl = "module decl\n\nimport std/string (endsWith)\nimport std/io (println)\n\n" +
		"type SendAction = { to: string, body: string }\n\n" +
		"export pure func sanitizeBody(rawBody: string<email>) -> string<sanitized> ! {Declassify}\n" +
		"ensures { result == \"[sanitized]\" }\n{\n  \"[sanitized]\"\n}\n\n" +
		"export pure func isInternal(addr: string) -> bool ! {}\nensures { result == endsWith(addr, \"@company.com\") }\n" +
		"{\n  endsWith(addr, \"@company.com\")\n}\n\n" +
		"export pure func safeForward(rawBody: string<email>, recipient: string) -> SendAction ! {Declassify}\n" +
		"requires { endsWith(recipient, \"@company.com\") }\n" +
		"ensures  { endsWith(result.to, \"@company.com\"), result.body == \"[sanitized]\" }\n" +
		"{\n  { to: recipient, body: sanitizeBody(rawBody) }\n}\n\n" +
		"export func main() -> () ! {IO, Declassify} {\n" +
		"  let action = safeForward(\"raw-confidential-body\", \"alice@company.com\");\n" +
		"  println(show(if isInternal(action.to) then 1 else 0))\n}\n"
	progEnv = "module envprobe\n\nimport std/env (getEnvOr)\nimport std/io (println)\n\n" +
		"export func main() -> () ! {IO, Env} {\n  println(getEnvOr(\"WORLD_SESSION\", \"<none>\"))\n}\n"
)

// progArgvPrint prints getArgs() joined by "|" (the V13 probe).
const progArgvPrint = "module argvprobe\n\nimport std/env (getArgs)\nimport std/io (println)\nimport std/string (join)\n\n" +
	"export func main() -> () ! {IO, Env} {\n  println(\"ARGV=[${join(\"|\", getArgs())}]\")\n}\n"

// netMock is a loopback HTTP mock that logs every request.
type netMock struct {
	srv  *httptest.Server
	mu   sync.Mutex
	reqs []string // "METHOD PATH X-Test-Header"
}

func newNetMock(t *testing.T, redirectTo func() string) *netMock {
	t.Helper()
	m := &netMock{}
	m.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.mu.Lock()
		m.reqs = append(m.reqs, r.Method+" "+r.URL.Path+" "+r.Header.Get("X-Test-Header"))
		m.mu.Unlock()
		if r.URL.Path == "/redirect" && redirectTo != nil {
			http.Redirect(w, r, redirectTo(), http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(m.srv.Close)
	return m
}

func (m *netMock) hostPort() string { return strings.TrimPrefix(m.srv.URL, "http://") }

func (m *netMock) port(t *testing.T) string {
	t.Helper()
	_, port, err := net.SplitHostPort(m.hostPort())
	if err != nil {
		t.Fatal(err)
	}
	return port
}

func (m *netMock) log() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.reqs...)
}

// runFixture is a real fixture plus a handler with the operator's full run
// allowlist, Net scoped to the named mock's one port.
type runFixture struct {
	*realFixture
	rh            *AilangToolHandler
	named, second *netMock
}

func newRunFixture(t *testing.T, topology string) *runFixture {
	t.Helper()
	f := newRealFixture(t, topology)
	second := newNetMock(t, nil)
	named := newNetMock(t, func() string { return second.srv.URL + "/" })
	for name, src := range map[string]string{"stdin.ail": progStdin, "sum.ail": progSum, "postnet.ail": progPostNet,
		"decl.ail": progDecl, "envprobe.ail": progEnv, "argvprobe.ail": progArgvPrint} {
		mustWrite(t, filepath.Join(f.root, name), src)
	}
	mustWrite(t, filepath.Join(f.root, "numbers.txt"), "1\n2\n3\n4\n5\n")
	mustWrite(t, filepath.Join(f.sibling, "outside.ail"), "module outside\n\nimport std/io (println)\n\n"+
		"export func main() -> () ! {IO} {\n  println(\"outside ran\")\n}\n")
	if err := os.Symlink("../ep2/outside.ail", filepath.Join(f.root, "linkentry.ail")); err != nil {
		t.Fatal(err)
	}
	op := RunCapsConfig{Allow: []string{RunCapDeclassify, RunCapEnv, RunCapNet},
		NetAllow: []string{named.hostPort()}, NetAllowHTTP: true}
	h, err := NewAilangToolHandler(boundedTestContext(t), AilangToolConfig{
		Bin: f.bin, BinRef: realBinRef(t, f.bin), PolicyPath: f.policyPath, Root: f.root, CacheDir: f.cacheDir, RunCaps: op,
	})
	if err != nil {
		t.Fatalf("NewAilangToolHandler: %v", err)
	}
	return &runFixture{realFixture: f, rh: h, named: named, second: second}
}

// exec runs one request through the run handler and decodes the result.
func (f *runFixture) exec(t *testing.T, payload map[string]any) runResult {
	t.Helper()
	body, _ := json.Marshal(payload)
	var caps []string
	if c, ok := payload["caps"].([]string); ok {
		caps = append(caps, c...)
		slices.Sort(caps)
	} else {
		caps = runBaseCaps
	}
	out, err := f.rh.Execute(boundedTestContext(t), EffectRequest{Effect: runEffectFor(caps), Scope: WorkspaceScope, Cost: 1}, body)
	if err != nil {
		t.Fatalf("run %s: %v", body, err)
	}
	var res runResult
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatal(err)
	}
	return res
}

func missingFromPolicy(t *testing.T, res runResult) string {
	t.Helper()
	var d struct {
		Missing []string `json:"missing_from_policy"`
	}
	_ = json.Unmarshal(res.Decision, &d)
	return fmt.Sprint(d.Missing)
}

// TestAilangRunCapsOnTheToolBinary is AC2.2, AC2.5–AC2.8 on the real binary.
func TestAilangRunCapsOnTheToolBinary(t *testing.T) {
	f := newRunFixture(t, "worktree")
	t.Setenv("WORLD_SESSION", "leaked-session-token")

	// The probe programs check on the tool binary first (S5).
	for _, prog := range []string{"stdin.ail", "sum.ail", "postnet.ail", "decl.ail", "envprobe.ail", "argvprobe.ail"} {
		if resp := f.op(t, EffectAilangCheck, map[string]any{"op": "ai_check", "path": prog}); resp["ok"] != true {
			t.Fatalf("fixture %s does not check on the tool binary: %v", prog, resp)
		}
	}

	// AC2.2: the four task shapes.
	if res := f.exec(t, map[string]any{"path": "stdin.ail", "caps": []string{"IO"}, "stdin": "1\n2\n3\n4\n5\n"}); !res.Admitted ||
		res.ExitCode != 0 || res.Stdout != "2\n4\n6\n8\n10\n" {
		t.Fatalf("stdin run = %+v", res)
	}
	res := f.exec(t, map[string]any{"path": "sum.ail", "caps": []string{"Env", "FS", "IO"}, "argv": []string{"numbers.txt"}})
	if !res.Admitted || res.ExitCode != 0 || res.Stdout != "15\n" || res.Policy == nil || res.Policy.SecurityMode != "trusted_host" ||
		fmt.Sprint(res.Policy.Caps) != "[Env FS IO]" {
		t.Fatalf("RunEnv argv run = %+v policy %+v", res, res.Policy)
	}
	url, _ := json.Marshal("http://" + f.named.hostPort() + "/")
	res = f.exec(t, map[string]any{"path": "postnet.ail", "caps": []string{"IO", "Net"}, "args_json": string(url)})
	if !res.Admitted || res.Stdout != "200\n" || res.Policy == nil || res.Policy.SecurityMode != "restricted" ||
		fmt.Sprint(res.Policy.NetAllow) != "["+f.named.hostPort()+"]" {
		t.Fatalf("RunNet run = %+v policy %+v", res, res.Policy)
	}
	if got := f.named.log(); fmt.Sprint(got) != "[POST / value123]" {
		t.Fatalf("named mock saw %q, want exactly one POST with X-Test-Header value123", got)
	}
	if res := f.exec(t, map[string]any{"path": "decl.ail", "caps": []string{"Declassify", "IO"}}); !res.Admitted ||
		res.ExitCode != 0 || res.Stdout != "1\n" || res.Policy.SecurityMode != "restricted" {
		t.Fatalf("Declassify run = %+v policy %+v", res, res.Policy)
	}

	// AC2.5: a Net program under every non-Net variant is refused before it
	// runs; under the Net variant every other port, name, address and
	// redirect hop is refused, and the second mock never sees a request.
	for _, caps := range [][]string{nil, {"IO"}, {"Env", "FS", "IO"}, {"Declassify", "IO"}} {
		payload := map[string]any{"path": "postnet.ail", "args_json": string(url)}
		if caps != nil {
			payload["caps"] = caps
		}
		res := f.exec(t, payload)
		if res.Admitted || res.ExitCode != 2 || missingFromPolicy(t, res) != "[Net]" {
			t.Fatalf("Net program under caps %v = %+v decision %s", caps, res, res.Decision)
		}
	}
	port := f.named.port(t)
	for _, tc := range []struct{ url, want string }{
		{f.second.srv.URL + "/", "DisallowedHost(127.0.0.1)"},
		{"http://127.0.0.1:7644/v1/health", "DisallowedHost(127.0.0.1)"},
		{"http://localhost:" + port + "/", "DisallowedHost(localhost)"},
		{"http://[::1]:" + port + "/", "DisallowedHost(::1)"},
		{"http://127.0.0.2:" + port + "/", "DisallowedHost(127.0.0.2)"},
		{"http://" + f.named.hostPort() + "/redirect", "E_NET_DOMAIN_BLOCKED"},
	} {
		u, _ := json.Marshal(tc.url)
		res := f.exec(t, map[string]any{"path": "postnet.ail", "caps": []string{"IO", "Net"}, "args_json": string(u)})
		if !res.Admitted || !strings.HasPrefix(res.Stdout, "ERR ") || !strings.Contains(res.Stdout, tc.want) {
			t.Fatalf("Net run to %s = %+v, want ERR %s", tc.url, res, tc.want)
		}
	}
	if got := f.second.log(); len(got) != 0 {
		t.Fatalf("the second loopback port saw %q, want no request", got)
	}

	// AC2.6: an Env run sees the handler's minimal env, never the parent's.
	if res := f.exec(t, map[string]any{"path": "envprobe.ail", "caps": []string{"Env", "IO"}}); res.Stdout != "<none>\n" {
		t.Fatalf("Env probe = %+v, want <none>", res)
	}
	// AC2.7: argv reaches the program verbatim, flags and `--` included.
	res = f.exec(t, map[string]any{"path": "argvprobe.ail", "caps": []string{"Env", "IO"}, "argv": []string{"--caps=IO,Net", "-", "--", "a b"}})
	if res.Stdout != "ARGV=[--caps=IO,Net|-|--|a b]\n" {
		t.Fatalf("argv probe = %+v", res)
	}
	// AC2.8: FS is always rendered, so the entry check holds for an IO-only
	// run; and an unrequested FS fails at use (budget 0).
	// The binary refuses it at policy load (rc 1, one stderr line, no
	// decision): none of row 134's outcome shapes, so the handler fails the
	// call (recorded `failed`) with the refusal verbatim — the program never
	// runs. (With FS not rendered, V23, it ran: "outside ran".)
	_, err := f.rh.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRun, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"linkentry.ail","caps":["IO"]}`))
	var envelope *AilangRunEnvelopeError
	if !errors.As(err, &envelope) || envelope.ExitCode != 1 || !strings.Contains(envelope.Stderr, "is outside fs_sandbox") ||
		strings.Contains(envelope.Stdout, "outside ran") {
		t.Fatalf("symlinked entry under caps [IO]: err = %v, want the load-time outside fs_sandbox refusal", err)
	}
	res = f.exec(t, map[string]any{"path": "readit.ail", "caps": []string{"IO", "Net"}})
	if !res.Admitted || !strings.Contains(res.Stderr+res.Stdout, "E_BUDGET_OPERATOR") || strings.Contains(res.Stdout, "hello from ep1") {
		t.Fatalf("FS read under caps [IO Net] = %+v", res)
	}
	f.assertProtectedUnchanged(t)
}

// TestAilangRunEnvConfinementMatrix is AC2.4: the trusted_host Env variant on
// both topologies refuses every `.git` form (V37), bytes unchanged; a sibling
// read escapes nothing; the Net variant keeps restricted `.git` protection.
func TestAilangRunEnvConfinementMatrix(t *testing.T) {
	for _, topology := range []string{"init", "worktree"} {
		t.Run(topology, func(t *testing.T) {
			f := newRunFixture(t, topology)
			mustMkdir(t, filepath.Join(f.root, "sub"))
			targets := []string{".git", ".GIT", ".Git"}
			if topology == "init" {
				targets = append(targets, ".git/config", ".GIT/config", "./.GIT/config", "sub/../.GiT/config",
					".Git/hooks/pre-commit", ".gIt/HEAD")
			}
			env := []string{"Env", "FS", "IO"}
			for _, target := range targets {
				arg, _ := json.Marshal(target)
				res := f.exec(t, map[string]any{"path": "wr.ail", "caps": env, "args_json": string(arg)})
				if res.ExitCode != 1 || !res.Admitted || !strings.Contains(res.Stderr, "E_FS_PROTECTED") || strings.Contains(res.Stdout, "wrote") {
					t.Fatalf("Env-variant write %s = %+v, want rc 1 E_FS_PROTECTED", target, res)
				}
				f.assertProtectedUnchanged(t)
			}
			if res := f.exec(t, map[string]any{"path": "wr.ail", "caps": env, "args_json": `"ok-env.txt"`}); res.ExitCode != 0 || res.Stdout != "wrote\n" {
				t.Fatalf("in-sandbox Env-variant write = %+v", res)
			}
			for _, out := range []string{"../ep2/secret.txt", "link.txt"} {
				arg, _ := json.Marshal(out)
				res := f.exec(t, map[string]any{"path": "args.ail", "caps": env, "args_json": string(arg)})
				if res.ExitCode != 1 || !strings.Contains(res.Stderr, "escapes sandbox") || strings.Contains(res.Stdout, "secret") {
					t.Fatalf("Env-variant read %s = %+v, want escapes sandbox", out, res)
				}
			}
			if topology == "init" {
				for _, target := range []string{".git/config", ".GIT/config"} {
					arg, _ := json.Marshal(target)
					res := f.exec(t, map[string]any{"path": "wr.ail", "caps": []string{"FS", "IO", "Net"}, "args_json": string(arg)})
					if res.ExitCode != 1 || !strings.Contains(res.Stderr, "E_FS_PROTECTED") {
						t.Fatalf("Net-variant write %s = %+v, want E_FS_PROTECTED", target, res)
					}
				}
			}
			f.assertProtectedUnchanged(t)
		})
	}
}

// TestAilangRunTamperedVariantIsRefused (AC2.1, MUT-VARIANT-UNVERIFIED on the
// real binary): an operator literal the policy layer refuses at load (a bare
// host, V35) never runs — the summary check refuses the variant first.
func TestAilangRunTamperedVariantIsRefused(t *testing.T) {
	f := newRunFixture(t, "init")
	f.rh.runCaps = RunCapsConfig{Allow: []string{RunCapNet}, NetAllow: []string{"127.0.0.1"}, NetAllowHTTP: true}
	_, err := f.rh.Execute(boundedTestContext(t), EffectRequest{Effect: EffectAilangRunNet, Scope: WorkspaceScope, Cost: 1},
		[]byte(`{"path":"readit.ail","caps":["IO","Net"]}`))
	var unverified *RunVariantUnverifiedError
	if !errors.As(err, &unverified) || !strings.Contains(err.Error(), "would open every loopback port") {
		t.Fatalf("bare-host Net variant: err = %v, want *RunVariantUnverifiedError naming the load refusal", err)
	}
}

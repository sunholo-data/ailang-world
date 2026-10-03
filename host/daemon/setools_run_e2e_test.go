package daemon

// Row 135 M3 (design_docs/planned/w-ailang-run-stdin-argv-caps.md §4.7, §6
// AC3.2–AC3.4): the four core gate tasks ailang-run could not run in row 134
// — pipeline (stdin), cli_args (argv + Env), api_call_json (Net to the
// loopback mock) and prompt_injection (Declassify) — run end to end through
// World: the checked-in se-tools package published into a test store, a
// daemon built with New(...) and the operator's --run-* allowlist (Net scoped
// to the mock's one port), real HTTP on /mcp/, the agent writing each
// reference solution with ailang-write and running it with ailang-run under
// the grader's caps, stdin and argv (§4.7). Reference solutions:
// testdata/row135 (pipeline, cli_args, api_call_json from the floor's
// testdata; prompt_injection is sunholo-data/ailang
// benchmarks/prompt_injection/expected_ailang_safe.ail with its module renamed
// to the grader's benchmark/solution).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// benchMock is the api_call_json mock: 200 for a POST carrying the task's
// header and JSON body, 400 otherwise; every request is logged.
type benchMock struct {
	srv  *httptest.Server
	mu   sync.Mutex
	seen []string
}

func newBenchMock(t *testing.T) *benchMock {
	t.Helper()
	m := &benchMock{}
	m.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Message string `json:"message"`
			Count   int    `json:"count"`
		}
		err := json.NewDecoder(r.Body).Decode(&body)
		m.mu.Lock()
		m.seen = append(m.seen, r.Method+" "+r.Header.Get("X-Test-Header")+" "+r.Header.Get("Content-Type"))
		m.mu.Unlock()
		if r.Method != http.MethodPost || r.Header.Get("X-Test-Header") != "value123" || err != nil || body.Count != 42 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(m.srv.Close)
	return m
}

func (m *benchMock) log() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.seen...)
}

// benchTask is one §4.7 mapping: the reference solution, the input files,
// the ailang-run arguments the grader's spec maps to, and the expected stdout.
type benchTask struct {
	id       string
	files    map[string]string
	args     map[string]any
	effect   string
	stdout   string
	wantCaps []string
}

func benchTasks(t *testing.T, mockURL string) []benchTask {
	t.Helper()
	solution := func(name string) string {
		raw, err := os.ReadFile(filepath.Join("testdata", "row135", name+".ail.txt"))
		if err != nil {
			t.Fatal(err)
		}
		return strings.ReplaceAll(string(raw), "{{MOCK_HTTP_URL}}", mockURL)
	}
	return []benchTask{
		{"pipeline", map[string]string{"benchmark/solution.ail": solution("pipeline")},
			map[string]any{"path": "benchmark/solution.ail", "caps": []string{"IO"}, "stdin": "1\n2\n3\n4\n5\n"},
			broker.EffectAilangRun, "2\n4\n6\n8\n10\n", []string{"FS", "IO"}},
		{"cli_args", map[string]string{"benchmark/solution.ail": solution("cli_args"), "numbers.txt": "1\n2\n3\n4\n5\n"},
			map[string]any{"path": "benchmark/solution.ail", "caps": []string{"IO", "FS", "Env"}, "argv": []string{"numbers.txt"}},
			broker.EffectAilangRunEnv, "15\n", []string{"Env", "FS", "IO"}},
		{"api_call_json", map[string]string{"benchmark/solution.ail": solution("api_call_json")},
			map[string]any{"path": "benchmark/solution.ail", "caps": []string{"Net", "IO"}},
			broker.EffectAilangRunNet, "200\n", []string{"FS", "IO", "Net"}},
		{"prompt_injection", map[string]string{"benchmark/solution.ail": solution("prompt_injection")},
			map[string]any{"path": "benchmark/solution.ail", "caps": []string{"Declassify", "IO"}},
			broker.EffectAilangRun, "1\n", []string{"Declassify", "FS", "IO"}},
	}
}

// seRunGrants are row 134's six grants plus the two run effects.
var seRunGrants = append(append([]string(nil), seGrantEffects...), broker.EffectAilangRunEnv, broker.EffectAilangRunNet)

// newSeRunRig is the se-tools rig with the operator's full run allowlist and
// Net scoped to the mock's one port (R-135-9: the mock binds first, and the
// operator names its port).
func newSeRunRig(t *testing.T, mock *benchMock) *seRig {
	t.Helper()
	return newSeRigWith(t, func(f wsFixture, cfg *Config) {
		// The grader's layout: the solution lives at benchmark/solution.ail
		// (ailang-write does not create directories).
		if err := os.MkdirAll(filepath.Join(f.root, "ep1", "benchmark"), 0o755); err != nil {
			t.Fatal(err)
		}
		cfg.RunAllowCaps = []string{"Declassify", "Env", "Net"}
		cfg.RunNetAllow = []string{strings.TrimPrefix(mock.srv.URL, "http://")}
		cfg.RunNetAllowHTTP = true
	})
}

// effectRequest decodes the request object the record of out's one effect
// names (§4.8: stdin, argv and caps are in the content-addressed request).
func (r *seRig) effectRequest(out map[string]any) (broker.EffectRecord, map[string]any) {
	r.t.Helper()
	w := worldOf(r.t, out)
	ctx := boundedTestContext(r.t)
	obj, ok, err := r.d.store.GetObject(ctx, hashref.MustParse(*w.Effects[0].Record))
	if err != nil || !ok {
		r.t.Fatalf("record: ok=%v err=%v", ok, err)
	}
	rec, err := broker.DecodeRecord(obj.Payload)
	if err != nil {
		r.t.Fatal(err)
	}
	reqObj, ok, err := r.d.store.GetObject(ctx, rec.RequestRef)
	if err != nil || !ok {
		r.t.Fatalf("request object %s: ok=%v err=%v", rec.RequestRef, ok, err)
	}
	// The request object is length-prefixed fields (effect, scope, cost,
	// payload) ending in the canonical payload JSON.
	raw := string(reqObj.Payload)
	if !strings.HasPrefix(raw, fmt.Sprintf("%d:%s", len(rec.Effect), rec.Effect)) || strings.Index(raw, "{") < 0 {
		r.t.Fatalf("request object %q does not name %s", raw, rec.Effect)
	}
	var req map[string]any
	if err := json.Unmarshal([]byte(raw[strings.Index(raw, "{"):]), &req); err != nil {
		r.t.Fatalf("request object %s: %v", reqObj.Payload, err)
	}
	// The request object wraps the effect payload; find it at any depth.
	if inner := findPayload(req); inner != nil {
		return rec, inner
	}
	r.t.Fatalf("request object %s holds no run payload", reqObj.Payload)
	return rec, nil
}

// TestSeToolsRunBenchmarkTasksEndToEnd is AC3.2: each of the four tasks
// commits, prints its expected stdout, names the right effect in its record,
// and its request object holds the stdin, argv and caps.
func TestSeToolsRunBenchmarkTasksEndToEnd(t *testing.T) {
	mock := newBenchMock(t)
	r := newSeRunRig(t, mock)
	token := r.mint("ep1", seRunGrants...)
	if got := r.toolsList(token); !reflect.DeepEqual(got, seToolNames) {
		t.Fatalf("tools/list = %v, want exactly the 8 tools %v", got, seToolNames)
	}
	for _, task := range benchTasks(t, mock.srv.URL) {
		for path, content := range task.files {
			wire, _ := r.call(token, "ailang-write", map[string]any{"path": path, "content": content})
			if out := wire.Result.StructuredContent; wire.Error != nil || out["ok"] != true {
				t.Fatalf("%s: ailang-write %s = %+v", task.id, path, wire)
			}
		}
		before := r.entryCount()
		wire, took := r.call(token, "ailang-run", task.args)
		out := wire.Result.StructuredContent
		if wire.Error != nil || wire.Result.IsError || out == nil {
			t.Fatalf("%s: tools/call = %+v", task.id, wire)
		}
		t.Logf("%s: stdout %q (%d ms)", task.id, out["stdout"], took.Milliseconds())
		if out["admitted"] != true || out["exit_code"] != float64(0) || out["stdout"] != task.stdout {
			t.Fatalf("%s: run output %v, want admitted, rc 0, stdout %q", task.id, out, task.stdout)
		}
		policy, _ := out["policy"].(map[string]any)
		if fmt.Sprint(policy["caps"]) != fmt.Sprint(task.wantCaps) {
			t.Fatalf("%s: policy %v, want caps %v", task.id, policy, task.wantCaps)
		}
		r.assertEffectRecord(task.id, out, task.effect)
		if after := r.entryCount(); after != before+1 {
			t.Fatalf("%s: log entries %d -> %d, want one commit", task.id, before, after)
		}
		rec, req := r.effectRequest(out)
		if rec.Effect != task.effect {
			t.Fatalf("%s: record effect %s, want %s", task.id, rec.Effect, task.effect)
		}
		for _, key := range []string{"stdin", "argv", "caps"} {
			if want, ok := task.args[key]; ok {
				if key == "caps" {
					caps := append([]string(nil), want.([]string)...)
					want = sortedStrings(caps)
				}
				if fmt.Sprint(req[key]) != fmt.Sprint(want) {
					t.Fatalf("%s: request object %s = %v, want %v (request %v)", task.id, key, req[key], want, req)
				}
			}
		}
	}
	if got := mock.log(); fmt.Sprint(got) != "[POST value123 application/json]" {
		t.Fatalf("the mock saw %q, want exactly the one api_call_json POST", got)
	}
}

func sortedStrings(xs []string) []string {
	out := append([]string(nil), xs...)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}

// TestSeToolsRunEnvGrantAndBudget is AC3.3: a session without the
// Ailang.RunEnv grant is denied the Env run and the program never executes;
// with a RunEnv budget of 1 the second Env run is denied:budget.
func TestSeToolsRunEnvGrantAndBudget(t *testing.T) {
	mock := newBenchMock(t)
	r := newSeRunRig(t, mock)
	ep := filepath.Join(r.f.root, "ep1")
	// The probe program leaves a marker file when it runs.
	writeFile(t, filepath.Join(ep, "marker.ail"), "module marker\n\nimport std/env (getArgs)\nimport std/fs (writeFile)\n"+
		"import std/io (println)\n\nexport func main() -> () ! {IO, FS, Env} {\n  writeFile(\"ran.txt\", \"ran\");\n  println(\"ran\")\n}\n")
	envRun := map[string]any{"path": "marker.ail", "caps": []string{"Env", "FS", "IO"}}

	noEnv := r.mint("ep1", broker.EffectAilangRun)
	wire, _ := r.call(noEnv, "ailang-run", envRun)
	out := wire.Result.StructuredContent
	w := worldOf(t, out)
	if len(w.Effects) != 1 || w.Effects[0].Status != "denied" || out["admitted"] == true {
		t.Fatalf("Env run without a RunEnv grant = %+v, want the effect denied", wire)
	}
	if _, err := os.Stat(filepath.Join(ep, "ran.txt")); !os.IsNotExist(err) {
		t.Fatalf("the denied Env run executed the program (%v)", err)
	}
	// A Declassify-only run needs only Ailang.Run (D-135-5 = A).
	writeFile(t, filepath.Join(ep, "benchmark", "solution.ail"), mustReadFile(t, filepath.Join("testdata", "row135", "prompt_injection.ail.txt")))
	wire, _ = r.call(noEnv, "ailang-run", map[string]any{"path": "benchmark/solution.ail", "caps": []string{"Declassify", "IO"}})
	if out := wire.Result.StructuredContent; out["stdout"] != "1\n" {
		t.Fatalf("Declassify run on the Ailang.Run grant alone = %+v", wire)
	}

	now := mustMintBudget(t, r, "ep1", map[string]int64{broker.EffectAilangRun: 5, broker.EffectAilangRunEnv: 1})
	wire, _ = r.call(now, "ailang-run", envRun)
	if out := wire.Result.StructuredContent; out["stdout"] != "ran\n" {
		t.Fatalf("first Env run with budget 1 = %+v", wire)
	}
	if err := os.Remove(filepath.Join(ep, "ran.txt")); err != nil {
		t.Fatal(err)
	}
	wire, _ = r.call(now, "ailang-run", envRun)
	out = wire.Result.StructuredContent
	w = worldOf(t, out)
	if len(w.Effects) != 1 || w.Effects[0].Status != "denied" {
		t.Fatalf("second Env run with budget 1 = %+v, want denied", wire)
	}
	rec, _ := r.effectRequestAny(out)
	if rec.Allowed || !strings.Contains(rec.Denial, "budget") {
		t.Fatalf("second Env run record = %+v, want denied:budget", rec)
	}
	if _, err := os.Stat(filepath.Join(ep, "ran.txt")); !os.IsNotExist(err) {
		t.Fatalf("the budget-denied Env run executed the program (%v)", err)
	}
}

// effectRequestAny decodes out's one effect record without requiring it to
// be allowed.
func (r *seRig) effectRequestAny(out map[string]any) (broker.EffectRecord, hashref.HashRef) {
	r.t.Helper()
	w := worldOf(r.t, out)
	if len(w.Effects) != 1 || w.Effects[0].Record == nil {
		r.t.Fatalf("world.effects = %+v", w.Effects)
	}
	ref := hashref.MustParse(*w.Effects[0].Record)
	obj, ok, err := r.d.store.GetObject(boundedTestContext(r.t), ref)
	if err != nil || !ok {
		r.t.Fatalf("record %s: ok=%v err=%v", ref, ok, err)
	}
	rec, err := broker.DecodeRecord(obj.Payload)
	if err != nil {
		r.t.Fatal(err)
	}
	return rec, ref
}

func mustMintBudget(t *testing.T, r *seRig, episode string, budgets map[string]int64) string {
	t.Helper()
	tok := r.mintBudgets(episode, budgets)
	return tok
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

// TestSeToolsRunReplayIsByteEqual is AC3.4: each new run shape, invoked over
// /a2a/ (whose result names the invocation id), replays from its committed
// record through a replay binder that has no registry — no handler runs —
// to the committed output byte for byte.
func TestSeToolsRunReplayIsByteEqual(t *testing.T) {
	mock := newBenchMock(t)
	r := newSeRunRig(t, mock)
	token := r.mint("ep1", seRunGrants...)
	for i, task := range benchTasks(t, mock.srv.URL) {
		for path, content := range task.files {
			writeFile(t, filepath.Join(r.f.root, "ep1", path), content)
		}
		code, raw := r.post("/a2a/", token, a2aSendBody(fmt.Sprintf("row135-replay-%d", i), "ailang-run", task.args))
		var wire struct {
			Result struct {
				Metadata struct {
					InvocationID string `json:"invocation_id"`
				} `json:"metadata"`
				Artifacts []struct {
					Parts []struct {
						Data map[string]any `json:"data"`
					} `json:"parts"`
				} `json:"artifacts"`
			} `json:"result"`
		}
		if code != http.StatusOK || json.Unmarshal(raw, &wire) != nil || wire.Result.Metadata.InvocationID == "" ||
			len(wire.Result.Artifacts) != 1 || len(wire.Result.Artifacts[0].Parts) != 1 {
			t.Fatalf("%s: tasks/send = %d %s", task.id, code, raw)
		}
		live := wire.Result.Artifacts[0].Parts[0].Data
		if live["stdout"] != task.stdout {
			t.Fatalf("%s: live output %v", task.id, live)
		}
		seen := len(mock.log())
		replayed, err := r.d.coord.Replay(boundedTestContext(t), wire.Result.Metadata.InvocationID,
			func(grants []broker.Capability, refs []hashref.HashRef) transitionreg.Binder {
				return broker.OpenReplayBinder(r.d.store, grants, refs)
			})
		if err != nil {
			t.Fatalf("%s: replay: %v", task.id, err)
		}
		var replayedOut map[string]any
		if err := json.Unmarshal(replayed, &replayedOut); err != nil || !reflect.DeepEqual(replayedOut, live) {
			t.Fatalf("%s: replayed %s, live %v (%v)", task.id, replayed, live, err)
		}
		if len(mock.log()) != seen {
			t.Fatalf("%s: the replay reached the network", task.id)
		}
	}
}

// findPayload returns the first object (depth-first) holding a "path" key.
func findPayload(v any) map[string]any {
	switch x := v.(type) {
	case map[string]any:
		if _, ok := x["path"]; ok {
			return x
		}
		for _, child := range x {
			if p := findPayload(child); p != nil {
				return p
			}
		}
	case string:
		var inner any
		if json.Unmarshal([]byte(x), &inner) == nil {
			if _, isObj := inner.(map[string]any); isObj {
				return findPayload(inner)
			}
		}
	}
	return nil
}

// mintBudgets mints a session with one grant per effect at the given budget.
func (r *seRig) mintBudgets(episode string, budgets map[string]int64) string {
	r.t.Helper()
	now := time.Now().Unix()
	grants := make([]broker.Capability, 0, len(budgets))
	for e, b := range budgets {
		grants = append(grants, broker.Capability{Effect: e, Scope: broker.WorkspaceScope, ExpiresAt: now + 7200, Budget: b})
	}
	tok, _, _, err := authority.Mint(boundedTestContext(r.t), r.d.store, episode, grants, 3600, now, nil)
	if err != nil {
		r.t.Fatalf("mint: %v", err)
	}
	return tok
}

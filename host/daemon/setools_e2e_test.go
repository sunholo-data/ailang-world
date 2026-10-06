package daemon

// Row 134 M5b (w-software-engineering-domain §6 AC5.3, AC5.4): the checked-in
// se-tools package, published into a TEST store through the production
// publisher core, served by a daemon built with New(...) and its real routes,
// and driven over real HTTP on /mcp/ and /a2a/ — the archived pinned
// interpreter runs every plan and finish phase, the archived tool binary
// (ToolBinaryRelease) runs every effect inside the episode's worktree.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang/serveapi/protocol"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/projection"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// seToolsManifest is the checked-in manifest, relative to this package.
const seToolsManifest = "../../packages/se-tools/transitions.json"

// The two measured releases this e2e runs on.
const (
	seInterpreterRelease = "AILANG v0.41.0"
	seToolRelease        = ToolBinaryRelease
)

var seToolNames = []string{
	"ailang-check", "ailang-cli", "ailang-edit", "ailang-read", "ailang-run",
	"ailang-write", "builtins-search", "examples-search", "workspace-exec",
}

// seToolBins returns the pinned interpreter and the tool binary, or skips
// with an explicit message: this e2e needs both real binaries.
func seToolBins(t *testing.T) (interp, tool string) {
	t.Helper()
	interp, tool = os.Getenv("AILANG_BIN"), os.Getenv(toolBinEnv)
	if interp == "" || tool == "" {
		t.Skipf("the se-tools daemon e2e needs both real binaries: AILANG_BIN (pinned %s, plan/finish phases) "+
			"and %s (%s, the tool effects); AILANG_BIN set=%t, %s set=%t; skipping",
			seInterpreterRelease, toolBinEnv, seToolRelease, interp != "", toolBinEnv, tool != "")
	}
	return interp, tool
}

func binaryVersion(t *testing.T, bin string) string {
	t.Helper()
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("%s --version: %v", bin, err)
	}
	return releaseFromVersion(string(out))
}

// TestSeToolsE2EBinaryVersions pins the two binaries the e2e measures: a
// wrong binary in either variable is a red test, never a silent re-baseline.
func TestSeToolsE2EBinaryVersions(t *testing.T) {
	interp, tool := seToolBins(t)
	if got := binaryVersion(t, interp); got != seInterpreterRelease {
		t.Fatalf("AILANG_BIN is %q, want the pinned %q", got, seInterpreterRelease)
	}
	if got := binaryVersion(t, tool); got != seToolRelease {
		t.Fatalf("%s is %q, want %q", toolBinEnv, got, seToolRelease)
	}
}

// ---------------------------------------------------------------------------
// rig
// ---------------------------------------------------------------------------

type seRig struct {
	t      *testing.T
	f      wsFixture
	d      *Daemon
	srv    *httptest.Server
	client *http.Client
}

const (
	seHelloProgram  = "module hello\n\nimport std/io (println)\n\nexport func main() -> () ! {IO} {\n  println(\"hello from ailang-run\")\n}\n"
	seBrokenProgram = "module broken\n\nexport func main() -> int {\n  \"oops\"\n}\n"
)

// seExamplesCorpus writes a one-example AILANG examples corpus in the layout
// v0.51.0 `examples search` reads (measured, V65): manifest.json for the
// descriptions and tags, and the searched .ail files under runnable/.
func seExamplesCorpus(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "runnable"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "manifest.json"),
		`{"examples":[{"path":"fold_demo.ail","status":"working","tags":["fold"],"description":"Sum a list with foldl"}]}`+"\n")
	writeFile(t, filepath.Join(dir, "runnable", "fold_demo.ail"),
		"-- fold_demo.ail - Sum a list with foldl\nmodule fold_demo\n")
}

// newSeRig serves the published se-tools package from a daemon whose
// --workspace-root holds the fixture worktree ep1 and whose --examples-dir is
// a fixture corpus outside that root.
func newSeRig(t *testing.T) *seRig {
	t.Helper()
	return newSeRigWith(t, func(f wsFixture, cfg *Config) {
		// Not at <base>/examples: the binary's CWD fallback (../../examples
		// from <root>/ep1) would find that one even without AILANG_EXAMPLES.
		corpus := filepath.Join(f.outside, "se-corpus")
		seExamplesCorpus(t, corpus)
		cfg.ExamplesDir = corpus
	})
}

func newSeRigWith(t *testing.T, configure func(f wsFixture, cfg *Config)) *seRig {
	t.Helper()
	interp, tool := seToolBins(t)
	f := newWSFixture(t)
	ep := filepath.Join(f.root, "ep1")
	writeFile(t, filepath.Join(ep, "hello.ail"), seHelloProgram)
	writeFile(t, filepath.Join(ep, "broken.ail"), seBrokenProgram)
	writeFile(t, filepath.Join(ep, "notes.txt"), "alpha beta\n")
	cfg := Config{DBPath: f.db, AilangBin: interp, WorkspaceRoot: f.root, ToolAilangBin: tool}
	configure(f, &cfg)
	d := mustWSDaemon(t, cfg)
	if d.coord == nil {
		t.Fatal("no coordinator with --ailang-bin set")
	}
	seGenesis(t, d)
	publishSeTools(t, d, f.db)
	srv := httptest.NewServer(d.Handler())
	t.Cleanup(srv.Close)
	return &seRig{t: t, f: f, d: d, srv: srv, client: &http.Client{Timeout: DefaultClientTimeout}}
}

// seGenesis commits a genesis world (the coordinator commits on a head).
func seGenesis(t *testing.T, d *Daemon) {
	t.Helper()
	obj := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")),
		SemanticID: "test/genesis", Provenance: "setools_e2e_test", Payload: []byte("genesis-state")}
	entryHash := hashref.SumSHA256([]byte("genesis-entry"))
	w := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: obj.Hash, LogHead: entryHash}
	if err := d.store.Commit(boundedTestContext(t), store.Commit{Objects: []store.Object{obj}, NextWorld: w, Entry: store.LogEntry{Header: store.LogHeader{
		EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")), Interpreter: daemonInterpreterRef(t, d),
		PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "setools_e2e_test"}, EntryHash: entryHash, TransitionRef: obj.Hash}}); err != nil {
		t.Fatalf("genesis: %v", err)
	}
}

// publishSeTools publishes the checked-in manifest the way `world-publish
// transitions` does after its attended fences: each transitionFnFile
// (repo-relative) canonicalised and stored as a world/transition-source/v1
// object, the daemon's archived interpreter as the pin, the epoch derived
// from the epoch registry the daemon bootstrapped, schemas through the
// codec's canonicalizer, then the production PublishSet (which re-checks
// sources, epochs and loadability).
func publishSeTools(t *testing.T, d *Daemon, db string) {
	t.Helper()
	ctx := boundedTestContext(t)
	raw, err := os.ReadFile(seToolsManifest)
	if err != nil {
		t.Fatal(err)
	}
	var entries []struct {
		ID, Title, Description, TransitionFnFile string
		InputSchema, OutputSchema                json.RawMessage
		Access                                   transitionreg.EffectRequirement
		DeclaredEffects                          []transitionreg.EffectRequirement
	}
	if err := json.Unmarshal(raw, &entries); err != nil {
		t.Fatal(err)
	}
	pub := transitionreg.NewPublisher(d.store, archive.New(db))
	interp := daemonInterpreterRef(t, d)
	epochs, _, err := pub.EpochsForInterpreter(ctx, interp)
	if err != nil || len(epochs) != 1 {
		t.Fatalf("epochs for the daemon's interpreter = %v %v, want exactly one", epochs, err)
	}
	changes := make([]transitionreg.Change, 0, len(entries))
	for _, e := range entries {
		src, err := os.ReadFile(filepath.Join("..", "..", e.TransitionFnFile))
		if err != nil {
			t.Fatal(err)
		}
		source, err := canon.Source(src)
		if err != nil {
			t.Fatal(err)
		}
		obj := store.Object{Hash: hashref.SumSHA256(source), InterfaceHash: hashref.SumSHA256([]byte("world/transition-source/v1")),
			SemanticID: "world/transition-source/v1", Provenance: "setools_e2e_test", Payload: source}
		if err := d.store.PutObject(ctx, obj); err != nil {
			t.Fatal(err)
		}
		in, err := transitionreg.CanonicalSchema(e.InputSchema)
		if err != nil {
			t.Fatal(err)
		}
		out, err := transitionreg.CanonicalSchema(e.OutputSchema)
		if err != nil {
			t.Fatal(err)
		}
		desc := transitionreg.Descriptor{ID: e.ID, TransitionFn: obj.Hash, Interpreter: interp, SemanticsEpoch: epochs[0],
			InputSchema: in, OutputSchema: out, Access: e.Access, DeclaredEffects: e.DeclaredEffects,
			Title: e.Title, Description: e.Description}
		changes = append(changes, transitionreg.Change{ID: e.ID, Descriptor: &desc})
	}
	res, err := pub.PublishSet(ctx, changes)
	if err != nil || res.Revision != 1 || len(changes) != len(seToolNames) {
		t.Fatalf("PublishSet = %+v %v (%d changes)", res, err, len(changes))
	}
}

// seGrants are the seven effect grants: §4.2's six (one per effect name)
// plus row 140's Workspace.Exec.
var seGrantEffects = []string{
	broker.EffectWorkspaceRead, broker.EffectWorkspaceWrite, broker.EffectAilangCheck,
	broker.EffectAilangRun, broker.EffectAilangDiscover, broker.EffectAilangCLI, broker.EffectWorkspaceExec,
}

func (r *seRig) mint(episode string, effects ...string) string {
	r.t.Helper()
	now := time.Now().Unix()
	grants := make([]broker.Capability, 0, len(effects))
	for _, e := range effects {
		grants = append(grants, broker.Capability{Effect: e, Scope: broker.WorkspaceScope, ExpiresAt: now + 7200, Budget: 20})
	}
	tok, _, _, err := authority.Mint(boundedTestContext(r.t), r.d.store, episode, grants, 3600, now, nil)
	if err != nil {
		r.t.Fatalf("mint: %v", err)
	}
	return tok
}

// post sends one JSON body to route and returns the status and body.
func (r *seRig) post(route, token, body string) (int, []byte) {
	r.t.Helper()
	req, err := http.NewRequest(http.MethodPost, r.srv.URL+route, strings.NewReader(body))
	if err != nil {
		r.t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := r.client.Do(req)
	if err != nil {
		r.t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		r.t.Fatal(err)
	}
	return resp.StatusCode, raw
}

type mcpWire struct {
	Result struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		StructuredContent map[string]any `json:"structuredContent"`
		IsError           bool           `json:"isError"`
	} `json:"result"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (r *seRig) mcp(token, body string) mcpWire {
	r.t.Helper()
	code, raw := r.post("/mcp/", token, body)
	payload := strings.TrimSpace(strings.TrimPrefix(string(raw), "event: message\ndata: "))
	var wire mcpWire
	if code != http.StatusOK || json.Unmarshal([]byte(payload), &wire) != nil {
		r.t.Fatalf("MCP %s = %d %s", body, code, raw)
	}
	return wire
}

func (r *seRig) toolsList(token string) []string {
	r.t.Helper()
	wire := r.mcp(token, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	if wire.Error != nil {
		r.t.Fatalf("tools/list error %+v", wire.Error)
	}
	names := []string{}
	for _, tool := range wire.Result.Tools {
		names = append(names, tool.Name)
	}
	sort.Strings(names)
	return names
}

// call runs one tools/call and returns the structured result (and the wall
// time of the HTTP round trip).
func (r *seRig) call(token, name string, args map[string]any) (mcpWire, time.Duration) {
	r.t.Helper()
	body, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 2, "method": "tools/call",
		"params": map[string]any{"name": name, "arguments": args}})
	if err != nil {
		r.t.Fatal(err)
	}
	start := time.Now()
	wire := r.mcp(token, string(body))
	return wire, time.Since(start)
}

// world returns the reserved provenance block of a committed output.
type seWorld struct {
	Effects []struct {
		ID     string  `json:"id"`
		Status string  `json:"status"`
		Record *string `json:"record"`
	} `json:"effects"`
	Plan string `json:"plan"`
}

func worldOf(t *testing.T, out map[string]any) seWorld {
	t.Helper()
	raw, err := json.Marshal(out["world"])
	if err != nil {
		t.Fatal(err)
	}
	var w seWorld
	if err := json.Unmarshal(raw, &w); err != nil || w.Plan == "" {
		t.Fatalf("output carries no world block: %v (%v)", out, err)
	}
	return w
}

// assertEffectRecord requires world.effects[0] ok and its record to resolve
// in the store as an allowed EffectRecordV1 for effect; the plan ref resolves.
func (r *seRig) assertEffectRecord(tool string, out map[string]any, effect string) {
	r.t.Helper()
	w := worldOf(r.t, out)
	if len(w.Effects) != 1 || w.Effects[0].Status != "ok" || w.Effects[0].Record == nil {
		r.t.Fatalf("%s: world.effects = %+v, want one ok effect with a record", tool, w.Effects)
	}
	ctx := boundedTestContext(r.t)
	obj, ok, err := r.d.store.GetObject(ctx, hashref.MustParse(*w.Effects[0].Record))
	if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 {
		r.t.Fatalf("%s: effect record %s ok=%v err=%v semantic=%q", tool, *w.Effects[0].Record, ok, err, obj.SemanticID)
	}
	rec, err := broker.DecodeRecord(obj.Payload)
	if err != nil || !rec.Allowed || rec.Failed || rec.Effect != effect || rec.Scope != broker.WorkspaceScope {
		r.t.Fatalf("%s: effect record = %+v %v, want an allowed %s@worktree", tool, rec, err, effect)
	}
	if _, ok, err := r.d.store.GetObject(ctx, hashref.MustParse(w.Plan)); err != nil || !ok {
		r.t.Fatalf("%s: plan ref %s does not resolve (%v %v)", tool, w.Plan, ok, err)
	}
}

func (r *seRig) entryCount() int64 {
	r.t.Helper()
	ctx := boundedTestContext(r.t)
	var n int64
	for {
		_, ok, err := r.d.store.GetLogEntry(ctx, n)
		if err != nil {
			r.t.Fatal(err)
		}
		if !ok {
			return n
		}
		n++
	}
}

// ---------------------------------------------------------------------------
// AC5.3 — /mcp/
// ---------------------------------------------------------------------------

func TestSeToolsMCPEndToEnd(t *testing.T) {
	r := newSeRig(t)
	ep := filepath.Join(r.f.root, "ep1")
	token := r.mint("ep1", seGrantEffects...)
	if got := r.toolsList(token); !reflect.DeepEqual(got, seToolNames) {
		t.Fatalf("tools/list = %v, want exactly %v", got, seToolNames)
	}

	timings := map[string]time.Duration{}
	type step struct {
		tool   string
		effect string
		args   map[string]any
		check  func(out map[string]any) error
	}
	str := func(v any) string { s, _ := v.(string); return s }
	steps := []step{
		{"ailang-read", broker.EffectWorkspaceRead, map[string]any{"path": "data.txt"}, func(out map[string]any) error {
			if out["ok"] != true || out["content"] != "hello from ep1\n" {
				return fmt.Errorf("read: want the file content")
			}
			return nil
		}},
		{"ailang-write", broker.EffectWorkspaceWrite, map[string]any{"path": "made.txt", "content": "made by ailang-write\n"}, func(out map[string]any) error {
			if got, err := os.ReadFile(filepath.Join(ep, "made.txt")); out["ok"] != true || err != nil || string(got) != "made by ailang-write\n" {
				return fmt.Errorf("write: file = %q %v", got, err)
			}
			return nil
		}},
		{"ailang-edit", broker.EffectWorkspaceWrite, map[string]any{"path": "notes.txt", "old_text": "beta", "new_text": "gamma"}, func(out map[string]any) error {
			if got, err := os.ReadFile(filepath.Join(ep, "notes.txt")); out["ok"] != true || err != nil || string(got) != "alpha gamma\n" {
				return fmt.Errorf("edit: file = %q %v", got, err)
			}
			return nil
		}},
		{"ailang-check", broker.EffectAilangCheck, map[string]any{"path": "broken.ail"}, func(out map[string]any) error {
			errs, _ := out["errors"].([]any)
			if out["passed"] != false || len(errs) == 0 || out["error_count"] == float64(0) {
				return fmt.Errorf("check: want passed:false with errors")
			}
			return nil
		}},
		{"ailang-run", broker.EffectAilangRun, map[string]any{"path": "hello.ail"}, func(out map[string]any) error {
			if out["admitted"] != true || out["exit_code"] != float64(0) || !strings.Contains(str(out["stdout"]), "hello from ailang-run") {
				return fmt.Errorf("run: want admitted, rc 0 and the program's stdout")
			}
			return nil
		}},
		// Row 134 break-2 fix: v0.51.0 policy-tool caps a CLI op's stdout at
		// 64 KiB (cli_ops.go capBytes) and `builtins list --json` is 82,774
		// bytes, so the plan asks for the 18,827-byte TEXT inventory and the
		// finish parses it (V64). A real query through /mcp/ returns matches.
		{"builtins-search", broker.EffectAilangDiscover, map[string]any{"query": "println"}, func(out map[string]any) error {
			matches, _ := out["matches"].([]any)
			if _, refused := out["refused"]; refused || len(matches) == 0 || out["count"] != float64(len(matches)) {
				return fmt.Errorf("builtins-search: want println matches, got none or a refusal")
			}
			found := false
			for _, m := range matches {
				e, _ := m.(map[string]any)
				if !strings.Contains(strings.ToLower(str(e["name"])), "println") {
					return fmt.Errorf("builtins-search: match %v does not contain the query", e)
				}
				found = found || (e["name"] == "_io_println" && e["module"] == "std/io" && e["effect"] == "io")
			}
			if !found {
				return fmt.Errorf("builtins-search: want _io_println [io] std/io among the matches")
			}
			return nil
		}},
		// Row 134 break-3 fix: the corpus is not in the binary (V65); serve
		// --examples-dir reaches the tool as AILANG_EXAMPLES, and the fixture
		// corpus's one example is found.
		{"examples-search", broker.EffectAilangDiscover, map[string]any{"query": "foldl"}, func(out map[string]any) error {
			if out["ok"] != true || !strings.Contains(str(out["stdout"]), "Found 1 examples") ||
				!strings.Contains(str(out["stdout"]), "fold_demo.ail") {
				return fmt.Errorf("examples-search: want the fixture corpus's fold_demo.ail found")
			}
			return nil
		}},
		{"ailang-cli", broker.EffectAilangCLI, map[string]any{"op": "version"}, func(out map[string]any) error {
			if out["ok"] != true || !strings.Contains(str(out["stdout"]), strings.TrimPrefix(seToolRelease, "AILANG ")) {
				return fmt.Errorf("cli version: want the tool binary's version on stdout")
			}
			return nil
		}},
	}
	for _, s := range steps {
		before := r.entryCount()
		wire, took := r.call(token, s.tool, s.args)
		timings[s.tool] = took
		if wire.Error != nil || wire.Result.IsError || wire.Result.StructuredContent == nil || len(wire.Result.Content) != 1 {
			t.Fatalf("%s: tools/call = %+v", s.tool, wire)
		}
		out := wire.Result.StructuredContent
		if err := s.check(out); err != nil {
			t.Fatalf("%s: %v\noutput: %v", s.tool, err, out)
		}
		var text map[string]any
		if err := json.Unmarshal([]byte(wire.Result.Content[0].Text), &text); err != nil || !reflect.DeepEqual(text, out) {
			t.Fatalf("%s: text content does not carry the structured output (%v)", s.tool, err)
		}
		r.assertEffectRecord(s.tool, out, s.effect)
		if after := r.entryCount(); after != before+1 {
			t.Fatalf("%s: log entries %d -> %d, want one commit", s.tool, before, after)
		}
	}
	spend, err := r.d.store.EffectSpend(boundedTestContext(t), "ep1")
	if err != nil {
		t.Fatal(err)
	}
	for effect, want := range map[string]int64{broker.EffectWorkspaceRead: 1, broker.EffectWorkspaceWrite: 2, broker.EffectAilangCheck: 1,
		broker.EffectAilangRun: 1, broker.EffectAilangDiscover: 2, broker.EffectAilangCLI: 1} {
		if got := spend[store.EffectKey{Effect: effect, Scope: broker.WorkspaceScope}]; got != want {
			t.Fatalf("spend[%s] = %d, want %d (all: %v)", effect, got, want, spend)
		}
	}

	// Plan refusal: an argument outside the tool's schema (L-ARGS) returns the
	// refusal text, runs no effect, and still commits.
	before, spendBefore := r.entryCount(), spend
	wire, took := r.call(token, "ailang-read", map[string]any{"path": "data.txt", "limit": 5})
	timings["ailang-read (L-ARGS refusal)"] = took
	out := wire.Result.StructuredContent
	if wire.Error != nil || out["ok"] != false || !strings.Contains(fmt.Sprint(out["refused"]), "limit") {
		t.Fatalf("refusal = %+v", wire)
	}
	if w := worldOf(t, out); len(w.Effects) != 0 {
		t.Fatalf("refusal ran effects: %+v", w.Effects)
	}
	if after := r.entryCount(); after != before+1 {
		t.Fatalf("refusal: log entries %d -> %d, want one commit", before, after)
	}
	if spendAfter, err := r.d.store.EffectSpend(boundedTestContext(t), "ep1"); err != nil || !reflect.DeepEqual(spendAfter, spendBefore) {
		t.Fatalf("refusal changed spend %v -> %v (%v)", spendBefore, spendAfter, err)
	}

	names := make([]string, 0, len(timings))
	for name := range timings {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		t.Logf("TIMING tools/call %-30s %6d ms", name, timings[name].Milliseconds())
	}
}

// TestSeToolsExamplesSearchWithoutCorpusRefuses: a daemon started with no
// examples corpus answers examples-search with a clear refusal through /mcp/
// (never an empty result, never a search of the agent's own worktree), and the
// effect is still recorded and committed.
func TestSeToolsExamplesSearchWithoutCorpusRefuses(t *testing.T) {
	r := newSeRigWith(t, func(wsFixture, *Config) {})
	// A worktree-local examples/ corpus must not be what answers: the binary
	// falls back to CWD-relative examples/ dirs when AILANG_EXAMPLES is unset.
	seExamplesCorpus(t, filepath.Join(r.f.root, "ep1", "examples"))
	token := r.mint("ep1", broker.EffectAilangDiscover)
	before := r.entryCount()
	wire, _ := r.call(token, "examples-search", map[string]any{"query": "foldl"})
	out := wire.Result.StructuredContent
	if wire.Error != nil || out["ok"] != false || out["refused"] != broker.NoExamplesCorpusRefusal {
		t.Fatalf("examples-search without a corpus = %+v, want ok:false refused %q", wire, broker.NoExamplesCorpusRefusal)
	}
	r.assertEffectRecord("examples-search", out, broker.EffectAilangDiscover)
	if after := r.entryCount(); after != before+1 {
		t.Fatalf("log entries %d -> %d, want one commit", before, after)
	}
}

// TestSeToolsWorkspaceExecWithoutProfileRefuses is row 140 M1 end to end: a
// daemon with no exec profile answers workspace-exec through /mcp/ with the
// typed refusal and still records and commits the effect (that it spawns
// nothing is TestWorkspaceExecUnconfiguredSpawnsNothing's, on the fake tool
// binary); the plan's own refusal of a bad id is a zero-effect commit.
func TestSeToolsWorkspaceExecWithoutProfileRefuses(t *testing.T) {
	r := newSeRig(t)
	token := r.mint("ep1", broker.EffectWorkspaceExec)
	before := r.entryCount()
	wire, _ := r.call(token, "workspace-exec", map[string]any{"command": "test", "args": []string{"-run=X", "./..."}})
	out := wire.Result.StructuredContent
	if wire.Error != nil || out["ok"] != false || out["refused"] != broker.NoExecProfileRefusal {
		t.Fatalf("workspace-exec without a profile = %+v, want ok:false refused %q", wire, broker.NoExecProfileRefusal)
	}
	r.assertEffectRecord("workspace-exec", out, broker.EffectWorkspaceExec)
	if after := r.entryCount(); after != before+1 {
		t.Fatalf("log entries %d -> %d, want one commit", before, after)
	}

	wire, _ = r.call(token, "workspace-exec", map[string]any{"command": "-x"})
	out = wire.Result.StructuredContent
	if wire.Error != nil || out["ok"] != false || !strings.HasPrefix(fmt.Sprint(out["refused"]), `command "-x" refused`) {
		t.Fatalf("workspace-exec -x = %+v, want the plan's id refusal", wire)
	}
	if w := worldOf(t, out); len(w.Effects) != 0 {
		t.Fatalf("the id refusal ran effects: %+v", w.Effects)
	}
}

// TestSeToolsQuickstartPayloadsVerbatim runs QUICKSTART §9's `-d` payloads
// byte for byte against the se-tools daemon with the §9 six-grant session.
func TestSeToolsQuickstartPayloadsVerbatim(t *testing.T) {
	r := newSeRig(t)
	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, section, ok := strings.Cut(string(guide), "### 9. Software-engineering tools")
	if !ok {
		t.Fatal("QUICKSTART §9 is absent")
	}
	// §9 ends where §10 (row 140's exec runbook, bound by its own tests) begins.
	section, _, _ = strings.Cut(section, "### 10. ")
	matches := regexp.MustCompile(`-d '([^']+)'`).FindAllStringSubmatch(section, -1)
	if len(matches) != 2 {
		t.Fatalf("§9 payload count = %d, want 2 (tools/list, one tools/call)", len(matches))
	}
	token := r.mint("ep1", seGrantEffects...)
	list := r.mcp(token, matches[0][1])
	names := []string{}
	for _, tool := range list.Result.Tools {
		names = append(names, tool.Name)
	}
	sort.Strings(names)
	if !reflect.DeepEqual(names, seToolNames) {
		t.Fatalf("§9 tools/list = %v, want %v", names, seToolNames)
	}
	call := r.mcp(token, matches[1][1])
	out := call.Result.StructuredContent
	if call.Error != nil || out["ok"] != true || out["content"] != seHelloProgram {
		t.Fatalf("§9 tools/call = %+v", call)
	}
	r.assertEffectRecord("§9 ailang-read", out, broker.EffectWorkspaceRead)
	// The §9 worktree program is the fixture program, byte for byte.
	if !strings.Contains(section, `println("hello from ailang-run")`) {
		t.Fatal("§9 no longer provisions the fixture's hello.ail")
	}
}

// TestSeToolsMCPVisibilityFollowsGrants: a session granted only
// Workspace.Read lists exactly ailang-read and cannot call ailang-write.
func TestSeToolsMCPVisibilityFollowsGrants(t *testing.T) {
	r := newSeRig(t)
	token := r.mint("ep1", broker.EffectWorkspaceRead)
	if got := r.toolsList(token); !reflect.DeepEqual(got, []string{"ailang-read"}) {
		t.Fatalf("Read-only tools/list = %v, want [ailang-read]", got)
	}
	before := r.entryCount()
	wire, _ := r.call(token, "ailang-write", map[string]any{"path": "x.txt", "content": "no"})
	if wire.Error == nil && !wire.Result.IsError {
		t.Fatalf("ailang-write without Workspace.Write = %+v, want refused", wire)
	}
	if _, err := os.Stat(filepath.Join(r.f.root, "ep1", "x.txt")); !os.IsNotExist(err) {
		t.Fatalf("the refused write created the file (%v)", err)
	}
	if after := r.entryCount(); after != before {
		t.Fatalf("the refused call committed (%d -> %d entries)", before, after)
	}
	// Control: the granted tool works on the same session.
	wire, _ = r.call(token, "ailang-read", map[string]any{"path": "data.txt"})
	if wire.Result.StructuredContent["content"] != "hello from ep1\n" {
		t.Fatalf("granted ailang-read = %+v", wire)
	}
}

// ---------------------------------------------------------------------------
// AC5.4 — /a2a/
// ---------------------------------------------------------------------------

func a2aSendBody(task, skill string, data map[string]any) string {
	raw, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 9, "method": "tasks/send", "params": map[string]any{
		"id": task, "metadata": map[string]any{"skill_id": skill},
		"message": map[string]any{"role": "user", "parts": []any{map[string]any{"type": "data", "data": data}}}}})
	return string(raw)
}

func TestSeToolsA2AEndToEnd(t *testing.T) {
	r := newSeRig(t)
	token := r.mint("ep1", broker.EffectWorkspaceRead)

	req, err := http.NewRequest(http.MethodGet, r.srv.URL+"/.well-known/agent.json", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := r.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var card struct{ Skills []struct{ ID string } }
	err = json.NewDecoder(resp.Body).Decode(&card)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil || len(card.Skills) != 1 || card.Skills[0].ID != "ailang-read" {
		t.Fatalf("Read-only card = %d %+v %v, want exactly [ailang-read]", resp.StatusCode, card, err)
	}

	code, raw := r.post("/a2a/", token, a2aSendBody("a2a-read-1", "ailang-read", map[string]any{"path": "data.txt"}))
	if code != http.StatusOK {
		t.Fatalf("tasks/send = %d %s", code, raw)
	}
	var wire struct {
		Result struct {
			Status    struct{ State string } `json:"status"`
			Artifacts []struct {
				Parts []struct {
					Type string         `json:"type"`
					Data map[string]any `json:"data"`
				} `json:"parts"`
			} `json:"artifacts"`
		} `json:"result"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal(raw, &wire); err != nil || wire.Error != nil || wire.Result.Status.State != "completed" ||
		len(wire.Result.Artifacts) != 1 || len(wire.Result.Artifacts[0].Parts) != 1 {
		t.Fatalf("tasks/send = %s (%v)", raw, err)
	}
	out := wire.Result.Artifacts[0].Parts[0].Data
	if out["content"] != "hello from ep1\n" {
		t.Fatalf("A2A ailang-read output = %v", out)
	}
	r.assertEffectRecord("a2a ailang-read", out, broker.EffectWorkspaceRead)
}

// r14Store is the injected post-effect failure: every Commit loses the head
// race (R14) after the effect has run and been recorded.
type r14Store struct{ *store.Store }

func (r14Store) Commit(context.Context, store.Commit) error { return &store.ConflictError{} }

// serveR14 rebinds r's projection to a coordinator over r14Store, logging
// operator refusals to errLog, and serves it from a fresh server.
func (r *seRig) serveR14(errLog io.Writer) {
	t := r.t
	t.Helper()
	coord, err := coordinator.New(coordinator.Config{Store: r14Store{r.d.store}, Runner: capsule.New(archive.New(r.f.db), capsule.Config{}),
		Binder: r.d.binder, Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	p, err := projection.New(projection.Config{CredentialWait: credentialBudget, CallbackTimeout: invokeDeadline, MaxCallbacks: 8, WriteWait: writeTimeout,
		Resolver: r.d.resolver, Reader: transitionreg.NewReader(r.d.store), Heads: r.d.reads, Deny: writeSessionDenial, Fail: writeAPIError,
		ErrorLog: errLog, Agent: protocol.AgentInfo{Name: "ailang-worldd", Version: Version},
		MaxWait: readDeadline, InvokeWait: invokeDeadline, Coordinator: coord})
	if err != nil {
		t.Fatal(err)
	}
	r.d.projection = p
	// d.Handler() binds the projection when it is built: serve the
	// replacement from a fresh server over the same daemon.
	r.srv.Close()
	r.srv = httptest.NewServer(r.d.Handler())
	t.Cleanup(r.srv.Close)
}

func TestSeToolsA2APostEffectFailureIsMapped(t *testing.T) {
	r := newSeRig(t)
	r.serveR14(io.Discard)
	token := r.mint("ep1", broker.EffectWorkspaceRead)
	before := r.entryCount()
	code, raw := r.post("/a2a/", token, a2aSendBody("a2a-r14", "ailang-read", map[string]any{"path": "data.txt"}))
	var wire struct {
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if code != http.StatusOK || json.Unmarshal(raw, &wire) != nil || wire.Error == nil {
		t.Fatalf("tasks/send under R14 = %d %s, want a JSON-RPC error", code, raw)
	}
	if !strings.HasPrefix(wire.Error.Message, projection.EffectsUnrecordedPrefix+" sha256:") {
		t.Fatalf("error message = %q, want %q + the effect record ref", wire.Error.Message, projection.EffectsUnrecordedPrefix)
	}
	ref := strings.TrimSpace(strings.TrimPrefix(wire.Error.Message, projection.EffectsUnrecordedPrefix))
	obj, ok, err := r.d.store.GetObject(boundedTestContext(t), hashref.MustParse(ref))
	if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 {
		t.Fatalf("record ref %q ok=%v err=%v", ref, ok, err)
	}
	if rec, err := broker.DecodeRecord(obj.Payload); err != nil || !rec.Allowed || rec.Effect != broker.EffectWorkspaceRead {
		t.Fatalf("effect record = %+v %v", rec, err)
	}
	if after := r.entryCount(); after != before {
		t.Fatalf("the conflicted call committed (%d -> %d)", before, after)
	}
}

// syncLog is an ErrorLog the hostcall goroutine writes while the test reads.
type syncLog struct {
	mu  sync.Mutex
	buf strings.Builder
}

func (l *syncLog) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.Write(p)
}

func (l *syncLog) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.String()
}

// mcpRefusalPrefix is the operator line an MCP tools/call refusal opens with
// for episode ep1 (row 136 AC1.1).
const mcpRefusalPrefix = "ailang-worldd: mcp refusal: tools/call a2a:ep1:"

// TestSeToolsMCPPostEffectFailureIsLogged is row 136 AC1.1 (kills
// MUT-MCP-LOG-NOCAUSE and MUT-MCP-LOG-LABEL): under R14 the /mcp/ wire is the
// released handler's frozen envelope (the tripwire for the upstream typed
// error), and the operator line, labelled mcp, names the minted invocation id
// and carries the effect-record ref, which resolves in the store.
func TestSeToolsMCPPostEffectFailureIsLogged(t *testing.T) {
	r := newSeRig(t)
	var log syncLog
	r.serveR14(&log)
	token := r.mint("ep1", broker.EffectWorkspaceRead)
	before := r.entryCount()
	wire, _ := r.call(token, "ailang-read", map[string]any{"path": "data.txt"})
	if wire.Error == nil || wire.Error.Code != -32603 || wire.Error.Message != "host callback failed" {
		t.Fatalf("tools/call under R14 error = %+v, want the frozen -32603 \"host callback failed\"", wire.Error)
	}
	text := log.String()
	if strings.Count(text, "\n") != 1 || !strings.HasPrefix(text, mcpRefusalPrefix) {
		t.Fatalf("operator log = %q, want exactly one line starting %q", text, mcpRefusalPrefix)
	}
	// The line's id is the id the coordinator dispatched: the cause names it.
	id := regexp.MustCompile(`^ailang-worldd: mcp refusal: tools/call (\S+): `).FindStringSubmatch(text)
	if len(id) != 2 || !strings.Contains(text, "coordinator: invocation "+id[1]+" requested effects") {
		t.Fatalf("operator line id is not the dispatched invocation id: %q", text)
	}
	refs := regexp.MustCompile(`sha256:[0-9a-f]{64}`).FindAllString(text, -1)
	if len(refs) == 0 {
		t.Fatalf("operator line carries no effect-record ref: %q", text)
	}
	resolved := false
	for _, ref := range refs {
		obj, ok, err := r.d.store.GetObject(boundedTestContext(t), hashref.MustParse(ref))
		if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 {
			continue
		}
		if rec, err := broker.DecodeRecord(obj.Payload); err == nil && rec.Allowed && rec.Effect == broker.EffectWorkspaceRead {
			resolved = true
		}
	}
	if !resolved {
		t.Fatalf("no ref on the operator line resolves to an allowed Workspace.Read effect record: %q", text)
	}
	if after := r.entryCount(); after != before {
		t.Fatalf("the conflicted call committed (%d -> %d)", before, after)
	}
}

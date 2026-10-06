package main

// Row 138 M1/M2 end to end (w-worldd-developer-cli §5 AC1.2, AC1.3, AC1.4,
// AC2.1 and the D-WORLD-60 trailer): the CLI verbs against a daemon built
// with daemon.New on a store holding the published checked-in se-tools
// package, the pinned interpreter running every plan/finish phase and the
// tool binary running every effect — the same rig as
// host/daemon/setools_e2e_test.go, rebuilt from exported APIs. Every entry
// `why` walks here was written by the real coordinator, so it cross-checks
// the coordinator-shaped mirror in why_test.go.

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

const toolBinEnvVar = "WORLD_TOOL_AILANG_BIN"

// cliSeBins returns the pinned interpreter and the tool binary, or skips
// with an explicit message.
func cliSeBins(t *testing.T) (string, string) {
	t.Helper()
	interp, tool := os.Getenv("AILANG_BIN"), os.Getenv(toolBinEnvVar)
	if interp == "" || tool == "" {
		t.Skipf("the CLI se-tools e2e needs both real binaries: AILANG_BIN (pinned AILANG v0.41.0) and %s (%s); "+
			"AILANG_BIN set=%t, %s set=%t; skipping", toolBinEnvVar, daemon.ToolBinaryRelease, interp != "", toolBinEnvVar, tool != "")
	}
	return interp, tool
}

type cliSeRig struct {
	url, db, session string
}

// newCLISeRig provisions the store (genesis when asked, the published
// se-tools package, a seven-grant ep1 session in a 0600 file) with the daemon
// stopped, then serves it.
func newCLISeRig(t *testing.T, genesis bool) cliSeRig {
	t.Helper()
	interp, tool := cliSeBins(t)
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	db, ws, corpus := filepath.Join(base, "state", "world.db"), filepath.Join(base, "ws"), filepath.Join(base, "outside", "corpus")
	for _, d := range []string{filepath.Dir(db), filepath.Join(ws, "ep1"), filepath.Join(corpus, "runnable")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFileT(t, filepath.Join(ws, "ep1", "data.txt"), "hello from ep1\n")
	writeFileT(t, filepath.Join(corpus, "manifest.json"), `{"examples":[]}`+"\n")
	cfg := daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, AilangBin: interp, WorkspaceRoot: ws,
		ToolAilangBin: tool, ExamplesDir: corpus, ErrorLog: io.Discard}

	// First start: archive the interpreter and bootstrap the epoch registry.
	d, err := daemon.New(boundedTestContext(t), cfg)
	if err != nil {
		t.Fatalf("daemon.New: %v", err)
	}
	srv := httptest.NewServer(d.Handler())
	_, body, err := newClient(srv.URL).get("/v1/health")
	srv.Close()
	if cerr := d.Close(); err != nil || cerr != nil {
		t.Fatalf("health %v close %v", err, cerr)
	}
	var h healthInfo
	if json.Unmarshal([]byte(body), &h) != nil || h.InterpreterRef == "" {
		t.Fatalf("health = %s", body)
	}
	interpRef := hashref.MustParse(h.InterpreterRef)

	st, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	ctx := boundedTestContext(t)
	if genesis {
		obj := store.Object{Hash: hashref.SumSHA256([]byte("genesis-state")), InterfaceHash: hashref.SumSHA256([]byte("test/genesis")),
			SemanticID: "test/genesis", Provenance: "cliwalk_e2e_test", Payload: []byte("genesis-state")}
		entryHash := hashref.SumSHA256([]byte("genesis-entry"))
		w := store.World{Ref: hashref.SumSHA256([]byte("genesis-world")), Revision: 0, StateRoot: obj.Hash, LogHead: entryHash}
		if err := st.Commit(ctx, store.Commit{Objects: []store.Object{obj}, NextWorld: w, Entry: store.LogEntry{Header: store.LogHeader{
			EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: hashref.SumSHA256([]byte("genesis-fn")), Interpreter: interpRef,
			PrevEntryHash: hashref.SumSHA256([]byte("genesis-prev")), WrittenBy: "cliwalk_e2e_test"}, EntryHash: entryHash, TransitionRef: obj.Hash}}); err != nil {
			t.Fatalf("genesis: %v", err)
		}
	}
	publishSeToolsCLI(t, st, db, interpRef)
	now := time.Now().Unix()
	var grants []broker.Capability
	for _, e := range []string{broker.EffectWorkspaceRead, broker.EffectWorkspaceWrite, broker.EffectAilangCheck,
		broker.EffectAilangRun, broker.EffectAilangDiscover, broker.EffectAilangCLI, broker.EffectWorkspaceExec} {
		grants = append(grants, broker.Capability{Effect: e, Scope: broker.WorkspaceScope, ExpiresAt: now + 7200, Budget: 20})
	}
	tok, _, _, err := authority.Mint(ctx, st, "ep1", grants, 3600, now, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	d, err = daemon.New(boundedTestContext(t), cfg)
	if err != nil {
		t.Fatalf("daemon.New (serve): %v", err)
	}
	srv = httptest.NewServer(d.Handler())
	t.Cleanup(func() {
		srv.Close()
		_ = d.Close()
	})
	return cliSeRig{url: srv.URL, db: db, session: sessionFile(t, tok)}
}

// publishSeToolsCLI publishes packages/se-tools/transitions.json as
// host/daemon's publishSeTools does: sources canonicalised and stored, the
// archived interpreter pinned, the production PublishSet.
func publishSeToolsCLI(t *testing.T, st *store.Store, db string, interp hashref.HashRef) {
	t.Helper()
	ctx := boundedTestContext(t)
	raw, err := os.ReadFile("../../packages/se-tools/transitions.json")
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
	pub := transitionreg.NewPublisher(st, archive.New(db))
	epochs, _, err := pub.EpochsForInterpreter(ctx, interp)
	if err != nil || len(epochs) != 1 {
		t.Fatalf("epochs = %v %v", epochs, err)
	}
	var changes []transitionreg.Change
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
			SemanticID: "world/transition-source/v1", Provenance: "cliwalk_e2e_test", Payload: source}
		if err := st.PutObject(ctx, obj); err != nil {
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
			InputSchema: in, OutputSchema: out, Access: e.Access, DeclaredEffects: e.DeclaredEffects, Title: e.Title, Description: e.Description}
		changes = append(changes, transitionreg.Change{ID: e.ID, Descriptor: &desc})
	}
	if _, err := pub.PublishSet(ctx, changes); err != nil || len(changes) != 9 {
		t.Fatalf("PublishSet: %v (%d changes)", err, len(changes))
	}
}

var refRe = regexp.MustCompile(`sha256:[0-9a-f]{64}`)

// TestCLIAgainstSeToolsDaemon is AC1.2, AC1.3 and AC2.1 against the real
// coordinator, plus the provenance trailer and log tail.
func TestCLIAgainstSeToolsDaemon(t *testing.T) {
	r := newCLISeRig(t, true)
	t.Logf("$ ailang-worldd tools list --session <file>")
	code, stdout, stderr := runCLI(t, r.url, "tools", "list", "--session", r.session)
	t.Logf("%s", stdout)
	if code != exitOK || !strings.HasSuffix(stdout, "9 tool(s)\n") {
		t.Fatalf("AC1.3 tools list: exit %d %q %q", code, stdout, stderr)
	}
	for _, name := range []string{"ailang-check", "ailang-cli", "ailang-edit", "ailang-read", "ailang-run", "ailang-write", "builtins-search", "examples-search", "workspace-exec"} {
		if !strings.Contains(stdout, name+" ") {
			t.Fatalf("tools list lacks %s", name)
		}
	}

	t.Logf("$ ailang-worldd call ailang-read --session <file> --arg path=data.txt")
	code, stdout, stderr = runCLI(t, r.url, "call", "ailang-read", "--session", r.session, "--arg", "path=data.txt")
	t.Logf("%s", stdout)
	if code != exitOK || !strings.Contains(stdout, `hello from ep1\n`) || !strings.Contains(stdout, "world:\n  plan       sha256:") ||
		!regexp.MustCompile(`(?m)^  effect     \S+ ok sha256:[0-9a-f]{64}$`).MatchString(stdout) {
		t.Fatalf("AC1.3 call: exit %d\n%s\n%s", code, stdout, stderr)
	}
	// Its effect record resolves: why on the record ref walks entry 1.
	effLine := regexp.MustCompile(`(?m)^  effect .*$`).FindString(stdout)
	record := refRe.FindString(effLine)
	code, stdout, stderr = whyRun(t, r.url, "", record)
	assertChainOK(t, "AC1.3 effect record", code, stdout, stderr, 1)

	// AC1.2 + AC2.1: --json-out piped to why -, and the other target forms.
	code, out, stderr := runCLI(t, r.url, "call", "ailang-read", "--session", r.session, "--json", `{"path":"data.txt"}`, "--json-out")
	if code != exitOK {
		t.Fatalf("--json-out: exit %d %s", code, stderr)
	}
	outRef := hashref.SumSHA256([]byte(strings.TrimSuffix(out, "\n"))).String()
	t.Logf("$ ailang-worldd call ailang-read --json '{\"path\":\"data.txt\"}' --json-out | ailang-worldd why -")
	code, stdout, stderr = whyRun(t, r.url, out, "-")
	t.Logf("%s", stdout)
	assertChainOK(t, "why -", code, stdout, stderr, 2)
	if !strings.Contains(stdout, "✓ output  "+outRef) || !strings.Contains(stdout, "your result is these exact bytes") {
		t.Fatalf("AC1.2: sha256(--json-out − newline) %s is not the walked output ref:\n%s", outRef, stdout)
	}
	inv := regexp.MustCompile(`invocation (mcp:\S+)`).FindStringSubmatch(stdout)
	if len(inv) != 2 {
		t.Fatalf("no invocation id in the chain:\n%s", stdout)
	}
	for _, target := range []string{"2", "head", outRef, inv[1]} {
		code, stdout, stderr = whyRun(t, r.url, "", target)
		assertChainOK(t, "AC2.1 "+target, code, stdout, stderr, 2)
	}

	// A plan refusal (L-ARGS) commits with no effects: 6 links, all ✓; exit
	// 0 by default, 3 under --strict (D-CLI-3).
	if code, _, _ := runCLI(t, r.url, "call", "ailang-read", "--session", r.session, "--arg", "path=data.txt", "--arg-json", "limit=5"); code != exitOK {
		t.Fatalf("committed refusal exit %d", code)
	}
	if code, _, _ := runCLI(t, r.url, "call", "ailang-read", "--session", r.session, "--arg", "path=data.txt", "--arg-json", "limit=5", "--strict"); code != exitIntegrity {
		t.Fatalf("--strict committed refusal exit %d", code)
	}
	code, stdout, _ = whyRun(t, r.url, "", "head")
	if code != exitOK || !strings.Contains(stdout, "all 6 link(s) verified") || !strings.Contains(stdout, "no effects; result") {
		t.Fatalf("why of a refusal: exit %d\n%s", code, stdout)
	}

	t.Logf("$ ailang-worldd provenance --episode ep1")
	code, stdout, stderr = runCLI(t, r.url, "provenance", "--episode", "ep1")
	t.Logf("%s", stdout)
	if want := "World-Provenance: store=" + r.db + " episode=ep1 entries=1-4\n"; code != exitOK || stdout != want {
		t.Fatalf("provenance: exit %d %q %q, want %q", code, stdout, stderr, want)
	}
	code, stdout, _ = runCLI(t, r.url, "log", "tail")
	t.Logf("$ ailang-worldd log tail\n%s", stdout)
	if code != exitOK || len(strings.Split(strings.TrimSpace(stdout), "\n")) != 5 ||
		!strings.Contains(stdout, "#1 ") || !strings.Contains(stdout, "ep1 ailang-read [Workspace.Read ok]") ||
		!strings.Contains(stdout, "ep1 ailang-read [no effects]") {
		t.Fatalf("log tail: exit %d\n%s", code, stdout)
	}
}

// TestCLINoGenesisIsNoWorldHead is AC1.4 on the real daemon: with no genesis
// the call surfaces as the generic -32603 (V16), and the probe names it.
func TestCLINoGenesisIsNoWorldHead(t *testing.T) {
	r := newCLISeRig(t, false)
	code, stdout, stderr := runCLI(t, r.url, "call", "ailang-read", "--session", r.session, "--arg", "path=data.txt")
	if code != exitUsage || stdout != "" || !strings.Contains(stderr, "no world head") {
		t.Fatalf("AC1.4: exit %d stdout %q stderr %q", code, stdout, stderr)
	}
}

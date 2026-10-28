package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Every subprocess has file sinks, including cheap precondition probes.
func marketingChromeCommand(t *testing.T, limit time.Duration, bin string, args ...string) string {
	t.Helper()
	dir := t.TempDir()
	stdout, err := os.Create(filepath.Join(dir, "stdout"))
	if err != nil {
		t.Fatal(err)
	}
	defer stdout.Close()
	stderr, err := os.Create(filepath.Join(dir, "stderr"))
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()
	ctx, cancel := context.WithTimeout(t.Context(), limit)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	if err := cmd.Run(); err != nil {
		se, _ := os.ReadFile(stderr.Name())
		t.Fatalf("%s %q: %v stderr=%s (sinks: %s)", bin, args, err, se, dir)
	}
	out, err := os.ReadFile(stdout.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func TestMarketingChromeSixShots(t *testing.T) {
	chrome := os.Getenv("WORLD_CHROME")
	if chrome == "" {
		t.Skip("WORLD_CHROME unset: controller runs the real Chrome drill")
	}
	info, err := os.Stat(chrome)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		t.Fatalf("WORLD_CHROME must be an executable file: %q (%v)", chrome, err)
	}
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	var version struct{ Verdict, Version string }
	out := marketingChromeCommand(t, 60*time.Second, "python3", filepath.Join(root, "scripts", "marketing_capture.py"), "chrome-version", "--chrome", chrome)
	if err := json.Unmarshal([]byte(out), &version); err != nil || version.Verdict != "ok" || version.Version == "" {
		t.Fatalf("Chrome version must be ok: %s (%v)", out, err)
	}
	interp, tool := os.Getenv("AILANG_BIN"), os.Getenv("WORLD_TOOL_AILANG_BIN")
	if interp == "" || tool == "" {
		t.Fatal("AILANG_BIN and WORLD_TOOL_AILANG_BIN are required")
	}
	rectangles := os.Getenv("WORLD_MARKETING_RECTS")
	if rectangles == "" {
		t.Fatal("WORLD_MARKETING_RECTS is required; never guess rectangles")
	}
	// Mirror Rectangles/ShotRects/Rect; require the complete measured shot set.
	type rect struct {
		Name       string `json:"name"`
		X, Y, W, H int
	}
	var measured struct {
		Schema string `json:"schema"`
		Shots  map[string]struct {
			Background rect   `json:"background"`
			Content    []rect `json:"content"`
		} `json:"shots"`
	}
	raw, err := os.ReadFile(rectangles)
	if err != nil {
		t.Fatalf("WORLD_MARKETING_RECTS: %v", err)
	}
	if err := json.Unmarshal(raw, &measured); err != nil || measured.Schema != "world/demo/capture-rectangles/v1" || len(measured.Shots) != 6 {
		t.Fatalf("WORLD_MARKETING_RECTS invalid schema/shot set: %v", err)
	}
	for _, name := range []string{"workbench-light.png", "workbench-dark.png", "packet.png", "graph.png", "provenance.png", "phone.png"} {
		sr, ok := measured.Shots[name]
		if !ok || len(sr.Content) == 0 {
			t.Fatalf("WORLD_MARKETING_RECTS missing shot/content: %s", name)
		}
		seen := map[string]bool{}
		for _, r := range sr.Content {
			if r.Name == "" || seen[r.Name] {
				t.Fatalf("WORLD_MARKETING_RECTS empty/duplicate content name: %s", name)
			}
			seen[r.Name] = true
		}
		if name == "packet.png" {
			for _, required := range []string{"validation", "roster-proven", "roster-tested", "roster-attested", "roster-claimed"} {
				if !seen[required] {
					t.Fatalf("WORLD_MARKETING_RECTS packet missing %s", required)
				}
			}
		}
	}
	for bin, release := range map[string]string{interp: "v0.41.0", tool: daemon.ToolBinaryRelease} {
		if out := marketingChromeCommand(t, 60*time.Second, bin, "--version"); !strings.Contains(out, release) {
			t.Fatalf("%s must report %s: %s", bin, release, out)
		}
	}
	// All preconditions precede seeding and building.
	rig := newMarketingChromeRig(t, interp, tool)
	dir := t.TempDir()
	worldd, check := filepath.Join(dir, "ailang-worldd"), filepath.Join(dir, "marketingcheck")
	marketingChromeCommand(t, 120*time.Second, "go", "build", "-o", worldd, ".")
	marketingChromeCommand(t, 120*time.Second, "go", "build", "-o", check, filepath.Join(root, "scripts", "marketingcheck"))
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	bind := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	capture := filepath.Join(dir, "capture")
	marketingChromeCommand(t, 660*time.Second, "bash", filepath.Join(root, "scripts", "capture_marketing.sh"),
		"--db", rig.db, "--manifest", rig.manifest, "--proof-key-file", rig.key,
		"--ailang-bin", rig.interp, "--workspace-root", rig.ws, "--tool-ailang-bin", rig.tool,
		"--chrome", chrome, "--out", capture, "--rectangles", rectangles, "--worldd-bin", worldd,
		"--bind", bind, "--marketingcheck-bin", check, "--verify-twice")
	raw, err = os.ReadFile(filepath.Join(capture, "capture-manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]json.RawMessage
	if err := json.Unmarshal(raw, &manifest); err != nil || manifest == nil {
		t.Fatalf("capture-manifest.json must be parseable: %v", err)
	}
	raw, err = os.ReadFile(filepath.Join(capture, "diagnostics.txt"))
	if err != nil || !strings.Contains(string(raw), "chrome-version: "+version.Version+"\n") || !strings.Contains(string(raw), "inventory: equal\n") {
		t.Fatalf("diagnostics must record Chrome version and inventory equality: %s (%v)", raw, err)
	}
}

type marketingChromeRig struct{ db, manifest, key, ws, interp, tool string }

// Compose newDemoApprovalRig with an empty store, retaining and closing the seed
// daemon before capture. See the iteration-305 seeded-still probe, lines 64–144.
func newMarketingChromeRig(t *testing.T, interp, tool string) marketingChromeRig {
	t.Helper()
	t.Setenv("TMPDIR", t.TempDir())
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	db, ws, corpus := filepath.Join(base, "state", "world.db"), filepath.Join(base, "ws"), filepath.Join(base, "outside", "corpus")
	for _, dir := range []string{filepath.Dir(db), filepath.Join(ws, "demo", "world"), filepath.Join(corpus, "runnable")} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	source, err := os.ReadFile(filepath.Join("..", "..", demoLogepochPath))
	if err != nil {
		t.Fatal(err)
	}
	writeFileT(t, filepath.Join(ws, "demo", demoLogepochPath), string(source))
	writeFileT(t, filepath.Join(corpus, "manifest.json"), `{"examples":[]}`+"\n")
	keyPath := filepath.Join(base, "proof.key")
	key := freshProofKey(t)
	if err := os.WriteFile(keyPath, key[:], 0600); err != nil {
		t.Fatal(err)
	}
	errLog, err := os.Create(filepath.Join(base, "daemon-stderr"))
	if err != nil {
		t.Fatal(err)
	}
	defer errLog.Close()
	cfg := daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, AilangBin: interp, WorkspaceRoot: ws,
		ToolAilangBin: tool, ExamplesDir: corpus, ErrorLog: errLog, ProofKeyFile: keyPath, ProofRequiredIdentities: []string{"sameRef", "servesEntry"}}
	d, err := daemon.New(boundedTestContext(t), cfg)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(d.Handler())
	_, body, healthErr := newClient(srv.URL).get("/v1/health")
	srv.Close()
	if err := d.Close(); err != nil || healthErr != nil {
		t.Fatalf("health %v close %v", healthErr, err)
	}
	var h healthInfo
	if json.Unmarshal([]byte(body), &h) != nil || h.InterpreterRef == "" {
		t.Fatalf("health = %s", body)
	}
	st, err := store.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	interpRef := hashref.MustParse(h.InterpreterRef)
	publishSeToolsCLI(t, st, db, interpRef)
	publishDemoApproval(t, st, db, interpRef, false)
	now := time.Now().Unix()
	grants := []broker.Capability{
		{Effect: "world.apply", Scope: "world", ExpiresAt: now + 7200, Budget: 100},
		{Effect: "Human.Approve", Scope: "worktree", ExpiresAt: now + 7200, Budget: 20},
	}
	for _, effect := range []string{broker.EffectWorkspaceRead, broker.EffectAilangCLI} {
		grants = append(grants, broker.Capability{Effect: effect, Scope: broker.WorkspaceScope, ExpiresAt: now + 7200, Budget: 20})
	}
	token, _, _, err := authority.Mint(boundedTestContext(t), credsFor(t, db), "demo", grants, 3600, now, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	d, err = daemon.New(boundedTestContext(t), cfg)
	if err != nil {
		t.Fatal(err)
	}
	srv = httptest.NewServer(d.Handler())
	defer func() { srv.Close(); _ = d.Close() }()
	session, manifest := filepath.Join(base, "session"), filepath.Join(base, "manifest.json")
	if err := os.WriteFile(session, []byte(token), 0600); err != nil {
		t.Fatal(err)
	}
	code, out, se := demoRun(t, &demoTestRig{addr: srv.URL}, []string{"demo", "seed", "--db", db, "--session", session, "--ailang-bin", interp, "--proof-key-file", keyPath, "--manifest", manifest})
	if code != exitOK || !strings.Contains(out, manifest) {
		t.Fatalf("S1-S6 seed: code=%d out=%q err=%s", code, out, se)
	}
	srv.Close()
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	return marketingChromeRig{db: db, manifest: manifest, key: keyPath, ws: ws, interp: interp, tool: tool}
}

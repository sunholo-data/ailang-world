package daemon

// Row 153 M3b (AC3.1 GC, AC3.7): the daemon says, per descriptor, whether its
// compile-cache template is ready; rebuilds a missing one only AFTER the
// listener is up (never blocking serving); and prunes templates of interpreters
// nothing pins. The startup test needs only AILANG_BIN and the fake tool
// binary, so it runs in CI; it never skips.

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

var templateLine = regexp.MustCompile(`(?m)^ailang-worldd: capsule template (\S+): (ready|missing)$`)

// templateStates parses the startup status lines in log into id -> state,
// failing on a repeated id (exactly one line per descriptor).
func templateStates(t *testing.T, log string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, m := range templateLine.FindAllStringSubmatch(log, -1) {
		if _, dup := out[m[1]]; dup {
			t.Fatalf("two startup status lines for %s in:\n%s", m[1], log)
		}
		out[m[1]] = m[2]
	}
	return out
}

func waitFor(t *testing.T, what string, bound time.Duration, ok func() bool) {
	t.Helper()
	deadline := time.Now().Add(bound)
	for !ok() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out after %s waiting for %s", bound, what)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestDaemonStartReportsMissingTemplates(t *testing.T) {
	interp := os.Getenv("AILANG_BIN")
	if interp == "" {
		t.Fatal("AILANG_BIN is unset: this test runs the pinned interpreter's plan phase and never skips")
	}
	f := newWSFixture(t)
	cfg := Config{DBPath: f.db, BindHost: DefaultBindHost, AilangBin: interp, WorkspaceRoot: f.root,
		ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)}

	// Daemon #1 publishes the se-tools (building every template), mints the
	// token the later call uses, and records where ailang-read's template is.
	cfg1 := cfg
	cfg1.ErrorLog = testErrorLog(t)
	d1, err := New(boundedTestContext(t), cfg1)
	if err != nil {
		t.Fatalf("New #1: %v", err)
	}
	var closeOnce sync.Once
	closeD1 := func() { closeOnce.Do(func() { _ = d1.Close() }) }
	t.Cleanup(closeD1)
	seGenesis(t, d1)
	publishSeTools(t, d1, f.db)
	r1 := &seRig{t: t, f: f, d: d1}
	token := r1.mint("ep1", broker.EffectWorkspaceRead)
	_, rev, ok, err := transitionreg.NewReader(d1.store).CurrentRevision(boundedTestContext(t))
	if err != nil || !ok {
		t.Fatalf("read the published head: ok=%v err=%v", ok, err)
	}
	var readDir string
	var readInterp hashref.HashRef
	var readSrc []byte
	arch := archive.New(f.db)
	for _, desc := range rev.Entries {
		if desc.ID != "ailang-read" {
			continue
		}
		obj, found, err := d1.store.GetObject(boundedTestContext(t), desc.TransitionFn)
		if err != nil || !found {
			t.Fatalf("ailang-read source: found=%v err=%v", found, err)
		}
		readInterp, readSrc = desc.Interpreter, obj.Payload
		readDir = arch.CapsuleTemplateDir(readInterp, readSrc)
	}
	if readDir == "" || !arch.CapsuleTemplateReady(readInterp, readSrc) {
		t.Fatalf("ailang-read has no ready template after publication (dir %q)", readDir)
	}
	closeD1()
	if err := os.RemoveAll(readDir); err != nil {
		t.Fatal(err)
	}

	// (ii) New alone reports the states through CapsuleTemplates().
	var newLog syncLog
	cfgN := cfg
	cfgN.ErrorLog = &newLog
	dn, err := New(boundedTestContext(t), cfgN)
	if err != nil {
		t.Fatalf("New (status arm): %v", err)
	}
	got := map[string]string{}
	for _, st := range dn.CapsuleTemplates() {
		got[st.ID] = st.State
	}
	if err := dn.Close(); err != nil {
		t.Fatal(err)
	}
	wantStates := map[string]string{}
	for _, name := range seToolNames {
		wantStates[name] = TemplateReady
	}
	wantStates["ailang-read"] = TemplateMissing
	if !mapsEqual(got, wantStates) {
		t.Fatalf("CapsuleTemplates() = %v, want %v", got, wantStates)
	}
	if fromLog := templateStates(t, newLog.String()); !mapsEqual(fromLog, wantStates) {
		t.Fatalf("startup lines = %v, want the same states %v\n%s", fromLog, wantStates, newLog.String())
	}

	// (i), (iii), (iv): daemon #2 through Run, rebuild held behind the gate.
	var log syncLog
	gate := make(chan struct{})
	cfg2 := cfg
	cfg2.ErrorLog = &log
	cfg2.templateRebuildGate = gate
	pr, pw := io.Pipe()
	runCtx, cancelRun := context.WithCancel(context.Background())
	ran := make(chan error, 1)
	go func() {
		err := Run(runCtx, cfg2, pw)
		_ = pw.Close()
		ran <- err
	}()
	var gateOnce sync.Once
	release := func() { gateOnce.Do(func() { close(gate) }) }
	// Cleanup runs LIFO: release the gate first so a mutant that waits on it
	// cannot leak a blocked goroutine, then stop Run and wait for the store.
	t.Cleanup(func() {
		cancelRun()
		select {
		case <-ran:
		case <-time.After(shutdownTimeout + 15*time.Second):
			t.Error("Run did not return after cancellation")
		}
	})
	t.Cleanup(release)

	announced := make(chan string, 1)
	go func() {
		line, _ := bufio.NewReader(pr).ReadString('\n')
		announced <- line
	}()
	var line string
	select {
	case line = <-announced:
	case err := <-ran:
		t.Fatalf("Run returned %v before announcing", err)
	// The gate holds the rebuild forever, so a blocking rebuild never lets Run
	// announce: any bound kills that mutant, hence a generous one that cannot
	// flake on a loaded -race runner.
	case <-time.After(60 * time.Second):
		t.Fatalf("Run announced nothing within 60s with the template rebuild held: the rebuild is blocking startup (log: %s)", log.String())
	}
	if !strings.HasPrefix(line, ListenAnnouncePrefix) {
		t.Fatalf("announcement = %q", line)
	}
	url := strings.TrimSpace(strings.TrimPrefix(line, ListenAnnouncePrefix))

	// (i) the status lines were written before the announcement was consumed.
	if fromLog := templateStates(t, log.String()); !mapsEqual(fromLog, wantStates) {
		t.Fatalf("Run's startup lines = %v, want %v\n%s", fromLog, wantStates, log.String())
	}

	// (iv) a call to the missing descriptor, while the rebuild is held, runs
	// cold, is served, and is labelled.
	if n := strings.Count(log.String(), coldCompileMarker); n != 0 {
		t.Fatalf("%d cold-compile lines before any call", n)
	}
	r2 := &seRig{t: t, f: f, srv: &httptest.Server{URL: url}, client: &http.Client{Timeout: DefaultClientTimeout}}
	r2.call(token, "ailang-read", map[string]any{"path": "data.txt"}) // the effect's outcome is irrelevant: the plan ran
	if n := strings.Count(log.String(), coldCompileMarker); n != 1 {
		t.Fatalf("%d cold-compile lines after one call to the missing descriptor, want exactly 1\n%s", n, log.String())
	}
	if strings.Contains(log.String(), "rebuilt in") {
		t.Fatalf("the rebuild finished while the gate was held:\n%s", log.String())
	}

	// (iii) releasing the gate lets the rebuild run, after serving began.
	release()
	waitFor(t, "the ailang-read rebuilt line", 60*time.Second, func() bool {
		return strings.Contains(log.String(), "ailang-worldd: capsule template ailang-read: rebuilt in ")
	})
	if !arch.CapsuleTemplateReady(readInterp, readSrc) {
		t.Fatal("ailang-read's template is not ready after the rebuild line")
	}
	if strings.Contains(log.String(), "rebuild failed") {
		t.Fatalf("a rebuild failed:\n%s", log.String())
	}
	t.Logf("daemon #2 operator log:\n%s", log.String())
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// AC3.7 (v) / AC3.1: an interpreter whose `check` passes but writes no cache
// makes publication fail with the typed error, after the bounded retry.
func TestPublishWithoutATemplateIsRefused(t *testing.T) {
	stub := filepath.Join(t.TempDir(), "no-cache-ailang")
	script := "#!/bin/sh\ncase \"$1\" in\n  --version) echo \"" + daemonFakeRelease + "\";;\nesac\nexit 0\n"
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	dbPath := filepath.Join(t.TempDir(), "world.db")
	d, err := New(boundedTestContext(t), Config{DBPath: dbPath, BindHost: DefaultBindHost, AilangBin: stub})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	payload := []byte("transition source for stub")
	src := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("daemon-test/transition-source")),
		SemanticID: "daemon-test/transition-source", Provenance: "templates_test", Payload: payload}
	if err := d.store.PutObject(boundedTestContext(t), src); err != nil {
		t.Fatal(err)
	}
	desc := transitionreg.Descriptor{ID: "tools.stub", TransitionFn: src.Hash, Interpreter: daemonInterpreterRef(t, d), SemanticsEpoch: 1,
		InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`),
		Access:          transitionreg.EffectRequirement{Effect: "world.apply", Scope: "world", Cost: 1},
		DeclaredEffects: []transitionreg.EffectRequirement{{Effect: "world.apply", Scope: "world", Cost: 1}}}
	_, err = transitionreg.NewPublisher(d.store, archive.New(dbPath)).PublishSet(boundedTestContext(t),
		[]transitionreg.Change{{ID: desc.ID, Descriptor: &desc}})
	var build *archive.TemplateBuildError
	if !errors.As(err, &build) {
		t.Fatalf("PublishSet with an interpreter that writes no cache = %v, want *archive.TemplateBuildError", err)
	}
}

// P11 / AC3.1 GC: templates of interpreters no descriptor pins go, the pinned
// interpreter's stay, a stale in-flight .tmp-* goes, a fresh one stays.
func TestTemplateGCPrunesStaleInterpDigests(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	d, err := New(boundedTestContext(t), Config{DBPath: dbPath, BindHost: DefaultBindHost, AilangBin: daemonFakeInterpreter(t),
		ErrorLog: &syncLog{}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	root := archive.New(dbPath).CapsuleCacheRoot()
	stale := filepath.Join(root, strings.Repeat("0", 64), "x", "compile")
	current := filepath.Join(root, d.interpreterHash.Digest(), "y", "compile")
	for _, dir := range []string{stale, current, filepath.Join(root, ".tmp-old"), filepath.Join(root, ".tmp-new")} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(stale, "manifest.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(filepath.Join(root, ".tmp-old"), old, old); err != nil {
		t.Fatal(err)
	}

	d.startTemplateMaintenance(boundedTestContext(t))
	select {
	case <-d.templates.done:
	case <-time.After(20 * time.Second):
		t.Fatal("template maintenance did not finish")
	}
	exists := func(p string) bool { _, err := os.Stat(p); return err == nil }
	if exists(filepath.Join(root, strings.Repeat("0", 64))) {
		t.Error("the stale interpreter digest's templates were not pruned")
	}
	if exists(filepath.Join(root, ".tmp-old")) {
		t.Error("the stale .tmp-old was not pruned")
	}
	if !exists(current) {
		t.Error("the pinned interpreter's templates were pruned")
	}
	if !exists(filepath.Join(root, ".tmp-new")) {
		t.Error("a fresh .tmp-new (a concurrent publisher's) was pruned")
	}
}

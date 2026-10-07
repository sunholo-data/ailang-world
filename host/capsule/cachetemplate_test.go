package capsule

// Row 153 M3a: every run starts from a copy of the publication-built compile
// cache template, and a run without one is loud. These use the REAL pinned
// interpreter (AILANG_BIN) and a real se-tools source: the template is that
// interpreter's own cache, which a fake cannot model. They never skip.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

// syncBuf is a log sink safe for the run's goroutines.
type syncBuf struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuf) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuf) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// templateRig is the real interpreter archived, plus the workspace-exec
// se-tools source (module se_tools/exec, whose header differs from the staging
// path, so a COLD run prints the MOD010 warning and a warm one does not).
type templateRig struct {
	t      *testing.T
	a      *archive.Archive
	ref    hashref.HashRef
	source []byte
}

func newTemplateRig(t *testing.T) *templateRig {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: the capsule template is the pinned interpreter's own cache; never skip")
	}
	a := archive.New(filepath.Join(t.TempDir(), "world.db"))
	ref, err := a.Archive(bin)
	if err != nil {
		t.Fatalf("archive pinned interpreter: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join("..", "..", "packages", "se-tools", "se_tools", "exec.ail"))
	if err != nil {
		t.Fatal(err)
	}
	src, err := canon.Source(raw)
	if err != nil {
		t.Fatal(err)
	}
	return &templateRig{t: t, a: a, ref: ref, source: src}
}

// buildTemplate is publication's step: CheckSource builds and promotes it.
func (r *templateRig) buildTemplate() {
	r.t.Helper()
	res, err := r.a.CheckSource(context.Background(), r.ref, r.source)
	if err != nil || !res.Passed {
		r.t.Fatalf("CheckSource = (%+v, %v), want a pass that built the template", res, err)
	}
	if !r.a.CapsuleTemplateReady(r.ref, r.source) {
		r.t.Fatalf("no ready template at %s after CheckSource", r.a.CapsuleTemplateDir(r.ref, r.source))
	}
}

func (r *templateRig) dropTemplate() {
	r.t.Helper()
	if err := os.RemoveAll(r.a.CapsuleTemplateDir(r.ref, r.source)); err != nil {
		r.t.Fatal(err)
	}
}

var execPlanArgs = func() []byte {
	in, _ := json.Marshal(map[string]any{"phase": "plan", "args": map[string]any{"command": "test", "args": []string{"-run=X"}}})
	s, _ := json.Marshal(string(in))
	return s
}()

func (r *templateRig) run(runner *Runner) Result {
	r.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := runner.RunContext(ctx, Entry{Interpreter: r.ref, Source: r.source, Args: execPlanArgs})
	if err != nil {
		r.t.Fatalf("RunContext: %v (stderr %q)", err, res.Stderr)
	}
	if !bytes.Contains(res.Stdout, []byte(`"world/effect-plan/v1"`)) {
		r.t.Fatalf("stdout %q is not a plan", res.Stdout)
	}
	return res
}

// warmth is the deterministic witness (P-M4): the interpreter warns MOD010 on
// stderr only when it COMPILES a module whose header is not the staging path.
func compiled(res Result) bool { return bytes.Contains(res.Stderr, []byte("MOD010")) }

// digestTree fingerprints a directory by (path, size, sha256, mtime).
func digestTree(t *testing.T, dir string) string {
	t.Helper()
	var lines []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		lines = append(lines, fmt.Sprintf("%s %d %x %d", rel, info.Size(), sha256.Sum256(raw), info.ModTime().UnixNano()))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// AC3.2 (MUT-SHARE): a run never writes the template. The template is made
// deliberately INCOMPLETE (a measured run on a complete one changes no file
// content, P-M2), so a run that shared it would add files and rewrite
// manifest.json.
func TestCapsuleTemplateIsCopiedNotShared(t *testing.T) {
	r := newTemplateRig(t)
	r.buildTemplate()
	dir := r.a.CapsuleTemplateDir(r.ref, r.source)
	modules, err := filepath.Glob(filepath.Join(dir, "compile", "modules", "host__capsule__main*"))
	if err != nil || len(modules) == 0 {
		t.Fatalf("template has no compiled entry module (glob %v, err %v)", modules, err)
	}
	for _, m := range modules {
		if err := os.RemoveAll(m); err != nil {
			t.Fatal(err)
		}
	}
	before := digestTree(t, dir)
	res := r.run(New(r.a, Config{Log: &syncBuf{}}))
	if !compiled(res) {
		t.Fatalf("premise: the incomplete template did not force a recompile (stderr %q)", res.Stderr)
	}
	if after := digestTree(t, dir); after != before {
		t.Fatalf("a run changed the template.\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

var coldLine = regexp.MustCompile(`^capsule: cold compile [0-9a-f]+/[0-9a-f]{64} \(template missing\)$`)

// AC3.6 (MUT-SILENT-COLD, MUT-NO-PROMOTE): a template-less run is counted and
// logged exactly once, and its lazy promotion makes the next run warm.
func TestCapsuleColdRunIsLabelled(t *testing.T) {
	r := newTemplateRig(t)
	var log syncBuf
	runner := New(r.a, Config{Log: &log})

	if res := r.run(runner); !compiled(res) {
		t.Fatalf("premise: a run with no template did not compile (stderr %q)", res.Stderr)
	}
	lines := strings.Split(strings.TrimSuffix(log.String(), "\n"), "\n")
	if len(lines) != 1 || !coldLine.MatchString(lines[0]) || runner.ColdRuns() != 1 {
		t.Fatalf("after one cold run: log %q, ColdRuns %d; want exactly one matching line and 1", log.String(), runner.ColdRuns())
	}
	if !r.a.CapsuleTemplateReady(r.ref, r.source) {
		t.Fatal("the successful cold run did not promote a template")
	}

	if res := r.run(runner); compiled(res) {
		t.Fatalf("the run after promotion still compiled (stderr %q)", res.Stderr)
	}
	if got := log.String(); strings.Count(got, "\n") != 1 || runner.ColdRuns() != 1 {
		t.Fatalf("the warm run was logged or counted: log %q, ColdRuns %d", got, runner.ColdRuns())
	}
}

func median(ds []time.Duration) time.Duration {
	s := append([]time.Duration(nil), ds...)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	return s[len(s)/2]
}

// AC3.4 (MUT-NOTEMPLATE), r2: warm runs do not compile. The deterministic
// MOD010 witness is the gate; the timings are logged as evidence only. The
// original median-warm < median-cold/2 assertion was removed (judge B1): the
// ratio is machine-dependent (CI linux -race 0.53-0.63, idle mac 0.35) because
// fixed per-run costs (100 MB verify, copy-in) do not scale with compile cost.
func TestCapsuleWarmRunsDoNotCompile(t *testing.T) {
	r := newTemplateRig(t)
	runner := New(r.a, Config{Log: &syncBuf{}})
	var cold, warm []time.Duration
	for i := 0; i < 5; i++ {
		r.dropTemplate()
		start := time.Now()
		c := r.run(runner)
		cold = append(cold, time.Since(start))
		if !compiled(c) {
			t.Fatalf("cold run %d did not compile (stderr %q)", i, c.Stderr)
		}
		start = time.Now() // the cold run's lazy promotion made this one warm
		w := r.run(runner)
		warm = append(warm, time.Since(start))
		if compiled(w) {
			t.Fatalf("warm run %d compiled (stderr %q): the template was not copied in", i, w.Stderr)
		}
	}
	mc, mw := median(cold), median(warm)
	t.Logf("cold %v median %v; warm %v median %v (evidence only, not asserted)", cold, mc, warm, mw)
	if got := runner.ColdRuns(); got != 5 {
		t.Fatalf("ColdRuns = %d, want 5 (one per cold arm, none for warm arms)", got)
	}
}

// S8 canary: the pinned interpreter honours AILANG_CACHE_DIR — `check`
// populates it, and a run from the result does not recompile. A floor raise
// that changes either breaks the cache silently, so this is listed in the S8
// floor-raise inventory.
func TestPinnedInterpreterHonoursCacheDir(t *testing.T) {
	r := newTemplateRig(t)
	r.buildTemplate() // fails unless check wrote compile/manifest.json under CacheDirEnv
	manifest := filepath.Join(r.a.CapsuleTemplateDir(r.ref, r.source), "compile", "manifest.json")
	if info, err := os.Stat(manifest); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("%s = (%v, %v), want the manifest `check` wrote under %s", manifest, info, err, archive.CacheDirEnv)
	}
	runner := New(r.a, Config{Log: &syncBuf{}})
	if res := r.run(runner); compiled(res) {
		t.Fatalf("a run from the check-built template compiled (stderr %q)", res.Stderr)
	}
	r.dropTemplate()
	if res := r.run(runner); !compiled(res) {
		t.Fatalf("a cold run did not warn MOD010 (stderr %q): the warmth witness is gone", res.Stderr)
	}
}

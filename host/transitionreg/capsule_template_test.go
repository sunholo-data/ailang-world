package transitionreg

// Row 153 M3a: publication builds the capsule compile-cache template, and a
// pass that could not build one is a retryable, typed hard failure.

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// realPublisherStore is publisherStore over the REAL pinned interpreter
// (AILANG_BIN), bootstrapped with its own release so the epoch derivation
// holds. Never skips: the template is the real interpreter's cache.
func realPublisherStore(t *testing.T) (*store.Store, *archive.Archive, hashref.HashRef) {
	t.Helper()
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: the capsule template is the pinned interpreter's own cache; never skip")
	}
	dbPath := filepath.Join(t.TempDir(), "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	arch := archive.New(dbPath)
	ref, err := arch.Archive(bin)
	if err != nil {
		t.Fatalf("archive pinned interpreter: %v", err)
	}
	m, err := arch.ReadManifest(ref)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if _, _, err := registry.Bootstrap(boundedTestContext(t), s, releaseFromManifest(m.Version)); err != nil {
		t.Fatalf("bootstrap epoch registry: %v", err)
	}
	return s, arch, ref
}

func realSourceDescriptor(t *testing.T, s *store.Store, interpreter hashref.HashRef, id, src string) (Descriptor, []byte) {
	t.Helper()
	d := validDescriptor()
	d.ID = id
	d.Title = "title-" + id
	d.Interpreter = interpreter
	d.SemanticsEpoch = 1
	d.TransitionFn = putSource(t, s, []byte(src))
	return d, []byte(src)
}

// AC3.1 (MUT-PUBLISH-NOTEMPLATE): publishing two real sources leaves a ready
// template for each and no in-flight tmp behind.
func TestPublishBuildsCapsuleTemplates(t *testing.T) {
	s, arch, ref := realPublisherStore(t)
	a, srcA := realSourceDescriptor(t, s, ref, "tools.a", "module transitions/a\n\nexport func main(input: string) -> string {\n  input\n}\n")
	b, srcB := realSourceDescriptor(t, s, ref, "tools.b", "module transitions/b\n\nexport func main(input: string) -> string {\n  \"b:${input}\"\n}\n")
	if _, err := NewPublisher(s, arch).PublishSet(boundedTestContext(t), []Change{{ID: a.ID, Descriptor: &a}, {ID: b.ID, Descriptor: &b}}); err != nil {
		t.Fatalf("PublishSet: %v", err)
	}
	for name, src := range map[string][]byte{"tools.a": srcA, "tools.b": srcB} {
		if !arch.CapsuleTemplateReady(ref, src) {
			t.Fatalf("%s: no ready capsule template at %s after publication", name, arch.CapsuleTemplateDir(ref, src))
		}
	}
	entries, err := os.ReadDir(arch.CapsuleCacheRoot())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".tmp-") {
			t.Fatalf("in-flight template tmp %q left in capsule-cache/", e.Name())
		}
	}
}

// countingInterpreter is a fake whose `check` counts its invocations in
// counter and writes the cache manifest only from invocation writeFrom on
// (0 never). --version prints testFakeRelease so the epoch derivation holds.
func countingInterpreter(t *testing.T, counter string, writeFrom int) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "counting-ailang")
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then echo \"" + testFakeRelease + "\"; exit 0; fi\n" +
		"if [ \"$1\" = \"check\" ]; then\n" +
		"  n=$(cat '" + counter + "' 2>/dev/null || echo 0); n=$((n+1)); echo $n > '" + counter + "'\n" +
		"  if [ " + strconv.Itoa(writeFrom) + " -gt 0 ] && [ $n -ge " + strconv.Itoa(writeFrom) + " ] && [ -n \"$AILANG_CACHE_DIR\" ]; then\n" +
		"    mkdir -p \"$AILANG_CACHE_DIR/compile\" && printf '{}' > \"$AILANG_CACHE_DIR/compile/manifest.json\"\n" +
		"  fi\n" +
		"  exit 0\n" +
		"fi\n" +
		"exit 0\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func countingStore(t *testing.T, bin string) (*store.Store, *archive.Archive, hashref.HashRef) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "world.db")
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	arch := archive.New(dbPath)
	ref, err := arch.Archive(bin)
	if err != nil {
		t.Fatalf("archive fake: %v", err)
	}
	if _, _, err := registry.Bootstrap(boundedTestContext(t), s, testFakeRelease); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	return s, arch, ref
}

func invocations(t *testing.T, counter string) int {
	t.Helper()
	raw, err := os.ReadFile(counter)
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(raw)))
	if err != nil {
		t.Fatalf("counter %q: %v", raw, err)
	}
	return n
}

// AC3.1 retry (MUT-WEDGE): a template build that fails once is retried inside
// EnsureSourceLoadable, so a transient fault does not wedge publication; one
// that never succeeds surfaces the typed error after exactly 3 attempts.
func TestPublishRetriesTemplateBuildError(t *testing.T) {
	t.Run("second_attempt_succeeds", func(t *testing.T) {
		counter := filepath.Join(t.TempDir(), "count")
		s, arch, ref := countingStore(t, countingInterpreter(t, counter, 2))
		d := storedSourceDescriptor(t, s, ref, "tools.echo")
		if _, err := NewPublisher(s, arch).PublishSet(boundedTestContext(t), []Change{{ID: d.ID, Descriptor: &d}}); err != nil {
			t.Fatalf("PublishSet with a once-failing template build: %v", err)
		}
		if n := invocations(t, counter); n != 2 {
			t.Fatalf("check ran %d times, want exactly 2 (one failed build, one retry)", n)
		}
		if !arch.CapsuleTemplateReady(ref, []byte("transition source for tools.echo")) {
			t.Fatal("no ready template after the retry succeeded")
		}
	})
	t.Run("never_succeeds", func(t *testing.T) {
		counter := filepath.Join(t.TempDir(), "count")
		s, arch, ref := countingStore(t, countingInterpreter(t, counter, 0))
		d := storedSourceDescriptor(t, s, ref, "tools.echo")
		_, err := NewPublisher(s, arch).PublishSet(boundedTestContext(t), []Change{{ID: d.ID, Descriptor: &d}})
		var build *archive.TemplateBuildError
		if !errors.As(err, &build) {
			t.Fatalf("PublishSet error = %v, want *archive.TemplateBuildError", err)
		}
		if n := invocations(t, counter); n != templateBuildAttempts {
			t.Fatalf("check ran %d times, want exactly %d", n, templateBuildAttempts)
		}
		if _, _, ok, _ := NewReader(s).CurrentRevision(boundedTestContext(t)); ok {
			t.Fatal("a refused publish must leave no head")
		}
	})
}

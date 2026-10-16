package daemon

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Row 139 D3 / AC3, daemon half: a store whose registry has grown starts on
// any nominated release and refuses an un-nominated one at registry-bootstrap,
// printing the exact `world-publish epoch nominate` command.
func TestStartupOnExistingStoreRequiresANominatedRelease(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "world.db")
	first, _ := fakeInterpreter(t, t.TempDir(), "AILANG v0.98.0\nCommit: aaaaaaa\n")
	second, _ := fakeInterpreter(t, t.TempDir(), "AILANG v0.98.1\nCommit: bbbbbbb\n")
	start := func(bin string) (*Daemon, error) {
		return New(boundedTestContext(t), Config{DBPath: dbPath, BindHost: DefaultBindHost, AilangBin: bin})
	}

	// Fresh store: bootstraps epoch 1 on the first release, as before.
	d, err := start(first)
	if err != nil {
		t.Fatalf("fresh start: %v", err)
	}
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}

	// Existing store, un-nominated release: refused, naming the command.
	_, err = start(second)
	var se *StartupError
	if !errors.As(err, &se) || se.Stage != StageRegistry {
		t.Fatalf("un-nominated start = %v, want a %s StartupError", err, StageRegistry)
	}
	var nn *registry.NotNominatedError
	if !errors.As(err, &nn) {
		t.Fatalf("the refusal is not a typed *registry.NotNominatedError: %T", se.Err)
	}
	want := "world-publish epoch nominate --store " + dbPath + " --release v0.98.1 --candidate --gate-verdict <verdict.json> --binary " + second
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("refusal = %q, want it to carry the exact command %q", err.Error(), want)
	}

	// Nominate the second release (the verb's registry call), then it starts
	// — and so does the first, on the same grown head.
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("store.Open after a refused start: %v (the refusal leaked writer authority)", err)
	}
	_, head, _, err := registry.Load(boundedTestContext(t), s)
	if err != nil {
		t.Fatal(err)
	}
	_, grown, err := registry.Nominate(boundedTestContext(t), s, head, "AILANG v0.98.1", registry.ModeCandidate)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	for _, bin := range []string{second, first} {
		d, err := start(bin)
		if err != nil {
			t.Fatalf("start on a nominated release (%s): %v", bin, err)
		}
		h, ok, err := d.store.GetRegistryHead(boundedTestContext(t), registry.SemanticID)
		if err != nil || !ok || h != grown {
			_ = d.Close()
			t.Fatalf("startup moved the grown head: %s ok=%v err=%v, want %s", h, ok, err, grown)
		}
		if err := d.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

// An unpinned restart of a pinned store has no nomination command: the
// refusal says to serve a nominated release instead.
func TestUnpinnedStartOnPinnedStoreNamesTheNominatedReleases(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "world.db")
	bin, _ := fakeInterpreter(t, t.TempDir(), "AILANG v0.98.0\n")
	d, err := New(boundedTestContext(t), Config{DBPath: dbPath, BindHost: DefaultBindHost, AilangBin: bin})
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = New(boundedTestContext(t), Config{DBPath: dbPath, BindHost: DefaultBindHost})
	if err == nil || strings.Contains(err.Error(), "epoch nominate") || !strings.Contains(err.Error(), `"AILANG v0.98.0"`) {
		t.Fatalf("unpinned restart = %v, want a refusal naming the nominated release and no nominate command", err)
	}
}

package archive

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// pinnedInterpreter is the dev interpreter of record. Tests that need a real
// `ailang --version` skip when it is absent so CI (which has no pinned binary)
// stays green; the core archival/idempotence/mismatch paths use synthetic
// fixtures instead.
const pinnedInterpreter = "/tmp/ailang-v0300/ailang"

// fakeInterpreter writes a synthetic executable that echoes a fixed version
// string when invoked with --version, and returns its path plus the exact bytes
// it holds. body distinguishes different "interpreters" so their hashes differ.
func fakeInterpreter(t *testing.T, dir, body, version string) (string, []byte) {
	t.Helper()
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then printf '%s' '" + version + "'; exit 0; fi\n" +
		"# body marker: " + body + "\n"
	p := filepath.Join(dir, "ailang-fake")
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake interpreter: %v", err)
	}
	return p, []byte(script)
}

// storeDBPath returns a store DB path inside a fresh temp dir so the adjacent
// artifact tree "<store>.db.artifacts" is isolated per test.
func storeDBPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "world.db")
}

func TestArchiveResolvesAndWritesManifest(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	dir := t.TempDir()
	execPath, bytes := fakeInterpreter(t, dir, "alpha", "FakeAILANG v9.9.9\n")
	wantRef := hashref.SumSHA256(bytes)

	a := New(storeDBPath(t))
	ref, err := a.Archive(execPath)
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if ref.String() != wantRef.String() {
		t.Fatalf("archived ref = %q, want %q", ref.String(), wantRef.String())
	}

	// Resolver: HashRef -> archived executable path.
	got, err := a.Resolve(ref)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	// The archived path follows <root>/interpreters/<algo>/<digest>/ailang.
	wantPath := filepath.Join(a.Root(), "interpreters", ref.Algo(), ref.Digest(), "ailang")
	if got != wantPath {
		t.Fatalf("Resolve path = %q, want %q", got, wantPath)
	}

	// Archived bytes are byte-identical to the source stream.
	archived, err := os.ReadFile(got)
	if err != nil {
		t.Fatalf("read archived executable: %v", err)
	}
	if string(archived) != string(bytes) {
		t.Fatalf("archived bytes differ from source stream")
	}

	// Archived executable is read-only executable (0o555): no write bits.
	info, err := os.Stat(got)
	if err != nil {
		t.Fatalf("stat archived executable: %v", err)
	}
	if perm := info.Mode().Perm(); perm != archivedPerm {
		t.Fatalf("archived perm = %o, want %o", perm, archivedPerm)
	}

	// Sidecar manifest carries hash + version + size + OS + arch.
	m, err := a.ReadManifest(ref)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	if m.Hash != ref.String() {
		t.Errorf("manifest hash = %q, want %q", m.Hash, ref.String())
	}
	if m.Version != "FakeAILANG v9.9.9\n" {
		t.Errorf("manifest version = %q, want %q", m.Version, "FakeAILANG v9.9.9\n")
	}
	if m.Size != int64(len(bytes)) {
		t.Errorf("manifest size = %d, want %d", m.Size, len(bytes))
	}
	if m.OS != runtime.GOOS {
		t.Errorf("manifest OS = %q, want %q", m.OS, runtime.GOOS)
	}
	if m.Arch != runtime.GOARCH {
		t.Errorf("manifest arch = %q, want %q", m.Arch, runtime.GOARCH)
	}
}

func TestArchiveIsIdempotent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	dir := t.TempDir()
	execPath, _ := fakeInterpreter(t, dir, "alpha", "v1\n")

	a := New(storeDBPath(t))
	ref1, err := a.Archive(execPath)
	if err != nil {
		t.Fatalf("first Archive: %v", err)
	}
	// Re-archiving the same bytes is a no-op success returning the same ref.
	ref2, err := a.Archive(execPath)
	if err != nil {
		t.Fatalf("second Archive (idempotent): %v", err)
	}
	if ref1.String() != ref2.String() {
		t.Fatalf("idempotent re-archive changed ref: %q -> %q", ref1.String(), ref2.String())
	}
	// Still exactly one manifest and one executable in the slot.
	m, err := a.ReadManifest(ref1)
	if err != nil {
		t.Fatalf("ReadManifest after re-archive: %v", err)
	}
	if m.Hash != ref1.String() {
		t.Fatalf("manifest hash drifted after re-archive: %q", m.Hash)
	}
}

func TestArchiveAbsentArtifactIsReplayError(t *testing.T) {
	a := New(storeDBPath(t))
	_, err := a.Archive(filepath.Join(t.TempDir(), "does-not-exist"))
	re, ok := IsReplayError(err)
	if !ok {
		t.Fatalf("expected *ReplayError, got %v", err)
	}
	if re.Kind != KindUnsupportedPlatform {
		t.Fatalf("Kind = %q, want %q", re.Kind, KindUnsupportedPlatform)
	}
}

func TestArchiveUnsupportedPlatformOnDirectory(t *testing.T) {
	a := New(storeDBPath(t))
	// A directory is not a regular file -> unsupported platform.
	_, err := a.Archive(t.TempDir())
	re, ok := IsReplayError(err)
	if !ok {
		t.Fatalf("expected *ReplayError, got %v", err)
	}
	if re.Kind != KindUnsupportedPlatform {
		t.Fatalf("Kind = %q, want %q", re.Kind, KindUnsupportedPlatform)
	}
}

func TestArchiveHashMismatchAgainstExistingSidecar(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	dir := t.TempDir()
	execPath, bytes := fakeInterpreter(t, dir, "alpha", "v1\n")
	ref := hashref.SumSHA256(bytes)

	a := New(storeDBPath(t))
	// Pre-seed the slot's sidecar with a DIFFERENT recorded hash, simulating a
	// corrupted/tampered archive for this exact content-address slot.
	if err := os.MkdirAll(a.dirFor(ref), 0o755); err != nil {
		t.Fatalf("mkdir slot: %v", err)
	}
	tamperedHash := hashref.SumSHA256([]byte("some other interpreter")).String()
	if err := a.writeManifest(ref, Manifest{Hash: tamperedHash, Version: "x", Size: 1, OS: "x", Arch: "x"}); err != nil {
		t.Fatalf("seed tampered manifest: %v", err)
	}

	_, err := a.Archive(execPath)
	re, ok := IsReplayError(err)
	if !ok {
		t.Fatalf("expected *ReplayError, got %v", err)
	}
	if re.Kind != KindHashMismatch {
		t.Fatalf("Kind = %q, want %q", re.Kind, KindHashMismatch)
	}
}

func TestResolveAbsentIsReplayError(t *testing.T) {
	a := New(storeDBPath(t))
	ref := hashref.SumSHA256([]byte("never archived"))
	_, err := a.Resolve(ref)
	re, ok := IsReplayError(err)
	if !ok {
		t.Fatalf("expected *ReplayError, got %v", err)
	}
	if re.Kind != KindAbsentArtifact {
		t.Fatalf("Kind = %q, want %q", re.Kind, KindAbsentArtifact)
	}
}

func TestResolveZeroRefIsReplayError(t *testing.T) {
	a := New(storeDBPath(t))
	_, err := a.Resolve(hashref.HashRef{})
	re, ok := IsReplayError(err)
	if !ok {
		t.Fatalf("expected *ReplayError, got %v", err)
	}
	if re.Kind != KindAbsentArtifact {
		t.Fatalf("Kind = %q, want %q", re.Kind, KindAbsentArtifact)
	}
}

// TestArchivePinnedInterpreter exercises the real pinned interpreter end to end
// when present. It is skipped in CI (no pinned binary) so the synthetic tests
// above carry the core coverage.
func TestArchivePinnedInterpreter(t *testing.T) {
	if _, err := os.Stat(pinnedInterpreter); err != nil {
		t.Skipf("pinned interpreter %s absent; skipping", pinnedInterpreter)
	}
	a := New(storeDBPath(t))
	ref, err := a.Archive(pinnedInterpreter)
	if err != nil {
		t.Fatalf("Archive(pinned): %v", err)
	}
	if _, err := a.Resolve(ref); err != nil {
		t.Fatalf("Resolve(pinned): %v", err)
	}
	m, err := a.ReadManifest(ref)
	if err != nil {
		t.Fatalf("ReadManifest(pinned): %v", err)
	}
	if !strings.Contains(m.Version, "AILANG") {
		t.Errorf("pinned manifest version = %q, expected to contain AILANG", m.Version)
	}
	if m.OS != runtime.GOOS || m.Arch != runtime.GOARCH {
		t.Errorf("pinned manifest platform = %s/%s, want %s/%s", m.OS, m.Arch, runtime.GOOS, runtime.GOARCH)
	}
}

package archive

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

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

// ---------------------------------------------------------------------------
// w-archive-stderr-in-manifest (queue row 21): the version recorded in a
// sidecar manifest is the interpreter's STDOUT, never its stderr, and the probe
// that obtains it is bounded.
//
// Every fixture below emits its stderr UNCONDITIONALLY from the fake itself.
// That is deliberate: the live stimulus for this item is a released `ailang`
// that logs an Observatory size warning to stderr only while ~/.ailang/state
// exceeds a threshold, so a test shelling out to the REAL binary would pass
// with the defect still in place on any rig under that threshold. An
// environment-dependent red is a vacuous green somewhere else.
// ---------------------------------------------------------------------------

// realPollutedVersionPrefix is the VERBATIM first line of the one genuinely
// polluted manifest found on the authoring rig, at
// /private/tmp/world-demo.db.artifacts/interpreters/sha256/e9746fef…3fb5/manifest.json.
// It is pinned here as a fixture so the acceptance tooth is hermetic and
// reproduces on any rig, rather than depending on that tree surviving.
const realPollutedVersionPrefix = "2026/08/18 21:02:42 Observatory: 301MB (warn threshold: 200MB)\n"

// cleanPinnedVersion is the shape of a healthy `ailang --version` STDOUT: 168
// bytes whose first line begins with versionPrefix.
const cleanPinnedVersion = "AILANG v0.30.0\nCommit: e37b370\nFull:   e37b370d1d7a9c4e7136b319e38bec4d5f2bd9a0\nBuilt:  2026-07-19T09:27:00Z\n\nThe AI-First Programming Language\nCopyright (c) 2025-2026\n"

// fakeInterpreterWithStderr is fakeInterpreter's stderr-emitting sibling: the
// script writes stderrLine to fd 2 and version to fd 1, in that order, from one
// shell. Under CombinedOutput() both fds are the same pipe, so the merged read
// sees stderrLine FIRST -- which is exactly how the real defect stamps a
// wall-clock-varying log line onto the front of a recorded version.
func fakeInterpreterWithStderr(t *testing.T, dir, body, version, stderrLine string) (string, []byte) {
	t.Helper()
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then printf '%s' '" + stderrLine + "' 1>&2; printf '%s' '" + version + "'; exit 0; fi\n" +
		"# body marker: " + body + "\n"
	p := filepath.Join(dir, "ailang-fake-stderr")
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatalf("write stderr-emitting fake interpreter: %v", err)
	}
	return p, []byte(script)
}

// fakeInterpreterBlocking replaces the shell with /bin/sleep via `exec`, so the
// process exec.CommandContext kills IS the sleeper.
//
// This shape is load-bearing and was MEASURED, not assumed. With the naive
// `sleep 20; printf …` form the shell forks a sleep GRANDCHILD which inherits
// the stdout pipe; CommandContext SIGKILLs only the direct child (sh), and
// cmd.Wait() then blocks on the output-copy goroutines until the grandchild
// exits anyway. At a 200 ms bound that form was measured returning after
// 10.211 s / 10.081 s / 10.164 s (3/3) -- i.e. the deadline arm would red
// against a CORRECT implementation. With `exec`, the same 200 ms bound returned
// in 202 / 201 / 203 ms (3/3). Absolute /bin/sleep because cmd.Env is scrubbed
// and a PATH-dependent lookup is one more thing that can fail for the wrong
// reason.
func fakeInterpreterBlocking(t *testing.T, dir, body string) (string, []byte) {
	t.Helper()
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then exec /bin/sleep 20; fi\n" +
		"# body marker: " + body + "\n"
	p := filepath.Join(dir, "ailang-fake-blocking")
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatalf("write blocking fake interpreter: %v", err)
	}
	return p, []byte(script)
}

// fakeInterpreterCounting appends one byte to counterPath on EVERY invocation,
// before doing anything else. The counter file's length is therefore the exact
// number of times the archive executed this interpreter -- the only instrument
// that can distinguish "the heal converged" from "the heal never ran".
func fakeInterpreterCounting(t *testing.T, dir, body, version, counterPath string) (string, []byte) {
	t.Helper()
	script := "#!/bin/sh\n" +
		"printf 'x' >> '" + counterPath + "'\n" +
		"if [ \"$1\" = \"--version\" ]; then printf '%s' '" + version + "'; exit 0; fi\n" +
		"# body marker: " + body + "\n"
	p := filepath.Join(dir, "ailang-fake-counting")
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatalf("write counting fake interpreter: %v", err)
	}
	return p, []byte(script)
}

// execCount reports how many times a fakeInterpreterCounting has been invoked.
// An absent counter file is zero executions, not an error.
func execCount(t *testing.T, counterPath string) int {
	t.Helper()
	data, err := os.ReadFile(counterPath)
	if errors.Is(err, os.ErrNotExist) {
		return 0
	}
	if err != nil {
		t.Fatalf("read exec counter %q: %v", counterPath, err)
	}
	return len(data)
}

// T1 / AC1 (the load-bearing criterion). RED MUTATION M1: restore
// `out, err := cmd.CombinedOutput()` + `return string(out), nil` in
// probeVersion and BOTH assertions below fail.
func TestArchiveManifestVersionExcludesInterpreterStderr(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	const marker = "FAKE-STDERR-MARKER should never reach the manifest"
	const wantVersion = "AILANG v9.9.9-fake\n"

	dir := t.TempDir()
	execPath, _ := fakeInterpreterWithStderr(t, dir, "stderr-emitter", wantVersion, marker)

	a := New(storeDBPath(t))
	ref, err := a.Archive(execPath)
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	m, err := a.ReadManifest(ref)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	// (a) the FULL value is pinned, so even a reordered or partial merge reds.
	if m.Version != wantVersion {
		t.Errorf("manifest version = %q, want the interpreter's stdout %q", m.Version, wantVersion)
	}
	// (b) the marker is named directly, so the failure says WHY.
	if strings.Contains(m.Version, marker) {
		t.Errorf("manifest version contains the fake's STDERR marker %q: %q", marker, m.Version)
	}
}

// T2 arm 1 / AC4a: the idempotent path HEALS a sidecar polluted by a build that
// merged stderr, seeded with the REAL polluted bytes measured on this rig.
// RED MUTATIONS: M1 (the heal writes the merged string, so the equality fails)
// and M4a (`if false && …`, so the pollution is returned unhealed).
func TestArchiveHealsPollutedSidecarVersionOnIdempotentPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	dir := t.TempDir()
	execPath, _ := fakeInterpreterWithStderr(t, dir, "polluted-sidecar",
		cleanPinnedVersion, "2026/08/20 01:13:12 Observatory: 314MB (warn threshold: 200MB)")

	a := New(storeDBPath(t))
	ref, err := a.Archive(execPath)
	if err != nil {
		t.Fatalf("first Archive: %v", err)
	}

	// Rewrite the sidecar the way a pre-fix build left it: stderr merged in
	// front of the real stdout. This is the verbatim on-disk shape.
	polluted, err := a.ReadManifest(ref)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	polluted.Version = realPollutedVersionPrefix + cleanPinnedVersion
	if err := a.writeManifest(ref, polluted); err != nil {
		t.Fatalf("seed polluted sidecar: %v", err)
	}

	// The idempotent path: identical bytes, so nothing about the artifact
	// changes -- but the sidecar must be repaired in passing.
	if _, err := a.Archive(execPath); err != nil {
		t.Fatalf("second Archive (idempotent, healing): %v", err)
	}
	healed, err := a.ReadManifest(ref)
	if err != nil {
		t.Fatalf("ReadManifest after heal: %v", err)
	}
	if healed.Version != cleanPinnedVersion {
		t.Errorf("healed version = %q, want the interpreter's stdout %q", healed.Version, cleanPinnedVersion)
	}
	if strings.Contains(healed.Version, "Observatory") {
		t.Errorf("healed version still carries stderr chatter: %q", healed.Version)
	}
	// The artifact bytes and its content address are NEVER touched by the heal.
	if healed.Hash != ref.String() {
		t.Errorf("heal moved the recorded hash: %q, want %q", healed.Hash, ref.String())
	}
}

// T2 arm 2 (convergence -- asserted directly; no mutation, per the plan). The
// INNER compare `fresh != existing.Version` is what this pins: the stored
// version here fails the well-formedness prefix on EVERY pass, so the outer
// guard admits a probe every time, and only the inner compare stops the sidecar
// being rewritten. Remove that compare and the manifest's mtime moves.
func TestArchiveHealIsConvergentAndDoesNotRewriteAnAgreeingSidecar(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	dir := t.TempDir()
	// "v1\n" deliberately fails the versionPrefix test, so the probe runs on
	// every idempotent pass and the compare is the only thing preventing churn.
	execPath, _ := fakeInterpreter(t, dir, "convergent", "v1\n")

	a := New(storeDBPath(t))
	ref, err := a.Archive(execPath)
	if err != nil {
		t.Fatalf("first Archive: %v", err)
	}
	manifestPath := a.manifestPathFor(ref)
	before, err := os.Stat(manifestPath)
	if err != nil {
		t.Fatalf("stat sidecar: %v", err)
	}
	beforeBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read sidecar: %v", err)
	}

	// Sleep past filesystem timestamp resolution so a rewrite WOULD be visible;
	// without this the assertion could pass vacuously on a fast rig.
	time.Sleep(20 * time.Millisecond)
	if _, err := a.Archive(execPath); err != nil {
		t.Fatalf("second Archive (idempotent): %v", err)
	}
	after, err := os.Stat(manifestPath)
	if err != nil {
		t.Fatalf("stat sidecar after re-archive: %v", err)
	}
	afterBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read sidecar after re-archive: %v", err)
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Errorf("agreeing probe rewrote the sidecar: mtime %v -> %v", before.ModTime(), after.ModTime())
	}
	if string(afterBytes) != string(beforeBytes) {
		t.Errorf("sidecar bytes changed across an idempotent re-archive:\n before %q\n after  %q", beforeBytes, afterBytes)
	}
}

// T2 arm 3 (fail-loud -- a fixture arm, asserted directly; no mutation). An
// archived interpreter that cannot answer --version while its sidecar is being
// healed is a startup FAILURE, not a silently-kept stale value.
func TestArchiveHealFailsLoudWhenTheArchivedInterpreterCannotExecute(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	if os.Geteuid() == 0 {
		t.Skip("root ignores the execute permission bit, so the fixture cannot fail loud")
	}
	dir := t.TempDir()
	// "v1\n" fails versionPrefix, so the idempotent path probes.
	execPath, _ := fakeInterpreter(t, dir, "unexecutable", "v1\n")

	a := New(storeDBPath(t))
	ref, err := a.Archive(execPath)
	if err != nil {
		t.Fatalf("first Archive: %v", err)
	}
	// Strip execute from the ARCHIVED copy only; the source stays intact so the
	// bytes still hash to the same slot and the idempotent path is still taken.
	archived := a.pathFor(ref)
	if err := os.Chmod(archived, 0o444); err != nil {
		t.Fatalf("chmod archived interpreter: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(archived, 0o555) })

	_, err = a.Archive(execPath)
	re, ok := IsReplayError(err)
	if !ok {
		t.Fatalf("expected *ReplayError from the healing probe, got %v", err)
	}
	if re.Kind != KindExecFailure {
		t.Fatalf("Kind = %q, want %q", re.Kind, KindExecFailure)
	}
	if !strings.Contains(re.Detail, "healing") {
		t.Errorf("detail %q does not identify the heal as the failing step", re.Detail)
	}
}

// T2 arm 4 / AC7: the probe is BOUNDED. RED MUTATION M3: replace
// exec.CommandContext(ctx, …) with exec.Command(…) and the probe waits the
// fake's full 20 s (blowing the 8 s wall assertion) and then returns rc=0 with
// EMPTY stdout, so the DeadlineExceeded assertion fails too -- a deterministic
// ~20 s ASSERTION-red, not a hang.
//
// The bound is shrunk to 1 s, not the design doc's suggested 200 ms: a COLD
// first exec of a freshly written script was measured at 227 ms on this rig in
// a must-succeed control, i.e. 200 ms is a flake. 1 s is ~9x the measured cold
// start; the 20 s sleep is 20x the bound; the 8 s wall assertion is 8x the
// bound and 2.5x under the sleep, so an unbounded probe cannot satisfy it.
func TestArchiveProbeIsBoundedByProbeTimeoutDeadline(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	dir := t.TempDir()
	execPath, _ := fakeInterpreterBlocking(t, dir, "blocking")

	a := New(storeDBPath(t))
	// The readDeadline shrink-in-tests idiom: exercise the SHIPPED timeout
	// branch without a ten-second test.
	a.probeTimeout = 1 * time.Second

	start := time.Now()
	_, err := a.Archive(execPath)
	elapsed := time.Since(start)

	if elapsed >= 8*time.Second {
		t.Errorf("Archive against a 20s-blocking interpreter took %v under a 1s probeTimeout; the probe is not bounded", elapsed)
	}
	re, ok := IsReplayError(err)
	if !ok {
		t.Fatalf("expected *ReplayError after the deadline, got %v (elapsed %v)", err, elapsed)
	}
	if re.Kind != KindExecFailure {
		t.Errorf("Kind = %q, want %q", re.Kind, KindExecFailure)
	}
	// The explicit timeout CAUSE, asserted through the public ReplayError chain.
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("errors.Is(err, context.DeadlineExceeded) = false for %v", err)
	}
}

// T2 arm 5 / AC8: a HEALTHY sidecar performs ZERO process executions on the
// idempotent path. RED MUTATION M4b: `if true || …` -- the DUAL form, because
// `if false &&` cannot neuter a SKIP. Under it the counter reads 2, not 1.
//
// Without this arm an UNCONDITIONAL heal passes AC1-AC7 unchanged, which is
// precisely the vacuity the mission's rule 3e forbids: the conditional was the
// round-2 quorum's own load-bearing fix and would ship with nothing asserting
// it.
func TestArchiveHealPerformsNoProbeOnAHealthySidecar(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	dir := t.TempDir()
	// The counter lives OUTSIDE the archive tree; the fake's own bytes (which
	// embed this path) are what get hashed, and they never change mid-test.
	counter := filepath.Join(dir, "exec-count")
	// Version starts with versionPrefix, so the stored value is well-formed and
	// the heal must decline to probe it.
	execPath, _ := fakeInterpreterCounting(t, dir, "counting", cleanPinnedVersion, counter)

	a := New(storeDBPath(t))
	ref, err := a.Archive(execPath)
	if err != nil {
		t.Fatalf("first Archive: %v", err)
	}
	// Known-positive control for the instrument: the FRESH path must have
	// executed the interpreter exactly once. A counter that can only ever read
	// zero would make the assertion below vacuous.
	if got := execCount(t, counter); got != 1 {
		t.Fatalf("control: fresh archival executed the interpreter %d times, want 1 (the counter instrument is broken)", got)
	}
	before, err := a.ReadManifest(ref)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	if !strings.HasPrefix(before.Version, versionPrefix) {
		t.Fatalf("fixture precondition: stored version %q must begin %q", before.Version, versionPrefix)
	}

	// The idempotent path over a HEALTHY sidecar: a true no-op.
	if _, err := a.Archive(execPath); err != nil {
		t.Fatalf("second Archive (idempotent, healthy): %v", err)
	}
	if got := execCount(t, counter); got != 1 {
		t.Errorf("idempotent re-archive of a HEALTHY artifact executed the interpreter; counter = %d, want 1", got)
	}
	after, err := a.ReadManifest(ref)
	if err != nil {
		t.Fatalf("ReadManifest after re-archive: %v", err)
	}
	if after.Version != before.Version {
		t.Errorf("healthy sidecar version changed: %q -> %q", before.Version, after.Version)
	}
}

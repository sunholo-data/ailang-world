package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/pins"
	"github.com/sunholo-data/ailang-world/host/registry"
)

// TestV1ReproductionRestartOnANominatedCandidate reproduces design §1 V1 end
// to end (row 139 M2, AC3 + AC2): a store bootstrapped by the daemon on the
// REAL pinned interpreter (AILANG_BIN, the manifest's interpreter release)
// refuses to restart on another release — printing the nominate command —
// until `world-publish epoch nominate --candidate` (headless, a test-manifest
// row and a passing verdict) adds it; then the daemon starts on it, and still
// starts on the original pin. A third, un-nominated release is still refused.
//
// The candidate binary is a shell fake printing a manifest-shaped version
// line, unless WORLD_V1_CANDIDATE_BIN names a real released binary (the
// measured V1 binary is v0.54.11): then its own version line and its real
// sha256 are used. The pinned interpreter is never faked.
func TestV1ReproductionRestartOnANominatedCandidate(t *testing.T) {
	pin := os.Getenv("AILANG_BIN")
	if pin == "" {
		t.Fatal("AILANG_BIN unset: V1 is reproduced on the pinned released interpreter; never skip")
	}
	if got, want := firstVersionLine(t, pin), pins.VersionLine(pins.RoleInterpreter); got != want {
		t.Fatalf("AILANG_BIN is %q, want the manifest's interpreter release %q", got, want)
	}

	candidate := os.Getenv("WORLD_V1_CANDIDATE_BIN")
	if candidate == "" {
		candidate = filepath.Join(t.TempDir(), "ailang-candidate")
		if err := os.WriteFile(candidate, []byte("#!/bin/sh\nif [ \"$1\" = --version ]; then echo '"+testCandidateLine+"'; exit 0; fi\nexit 0\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	line := firstVersionLine(t, candidate)
	release, ok := pins.ReleaseOfVersionLine(line)
	if !ok {
		t.Fatalf("candidate version line %q is not a release line", line)
	}
	binary, err := os.ReadFile(candidate)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(binary)
	binaryDigest := hex.EncodeToString(sum[:])
	t.Logf("V1: pin %q, candidate %q (%s, sha256 %s)", pins.VersionLine(pins.RoleInterpreter), line, candidate, binaryDigest)

	db := filepath.Join(t.TempDir(), "state", "world.db")
	if err := os.MkdirAll(filepath.Dir(db), 0o755); err != nil {
		t.Fatal(err)
	}
	serve := func(bin string) error {
		d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, AilangBin: bin, ErrorLog: io.Discard})
		if err != nil {
			return err
		}
		return d.Close()
	}

	// 1. The store is bootstrapped on the real pin (epoch 1).
	if err := serve(pin); err != nil {
		t.Fatalf("first start on the pin: %v", err)
	}
	// 2. V1 as measured: the candidate is refused at registry-bootstrap, now
	//    with the exact command to run.
	err = serve(candidate)
	cmd, _ := registry.NominateCommand(db, line, candidate)
	if err == nil || !strings.Contains(err.Error(), "registry-bootstrap") || !strings.Contains(err.Error(), cmd) {
		t.Fatalf("un-nominated restart = %v, want a registry-bootstrap refusal carrying %q", err, cmd)
	}
	t.Logf("V1 refusal: %v", err)

	// 3. The loop's nomination: headless, verdict + --binary, against a
	//    test-manifest row reached only through the runEpochNominateWith seam
	//    (production argv uses the compiled-in manifest, which has no row for
	//    the candidate yet).
	tarball := strings.Repeat("cd", 32)
	flags := map[string]string{
		"store": db, "release": release, "binary": candidate,
		testManifestKey: writeTestManifestFor(t, release, line, binaryDigest, tarball),
		"gate-verdict": writeVerdict(t, func(v map[string]any) {
			v["release"], v["binary_sha256"], v["tarball_sha256"] = release, binaryDigest, tarball
		}),
	}
	res := driveNominate(t, nominateInv(flags, "candidate"), "", noEnv, noTTYProbe(t))
	if res.code != exitOK {
		t.Fatalf("nominate = (%d, %q)", res.code, res.stderr)
	}
	t.Logf("nominate: %s", res.stdout)

	// 4. The candidate now starts, and the original pin still does.
	if err := serve(candidate); err != nil {
		t.Fatalf("restart on the nominated candidate: %v", err)
	}
	if err := serve(pin); err != nil {
		t.Fatalf("restart on the original pin after nomination: %v", err)
	}
	reg := storeRegistry(t, db)
	if len(reg.Epochs) != 1 || strings.Join(reg.Epochs[0].Candidates, "|") != pins.VersionLine(pins.RoleInterpreter)+"|"+line {
		t.Fatalf("registry = %+v, want epoch 1 nominating the pin then the candidate", reg.Epochs)
	}

	// 5. Fail-closed is intact: another release is still refused.
	other := filepath.Join(t.TempDir(), "ailang-other")
	if err := os.WriteFile(other, []byte("#!/bin/sh\necho 'AILANG v0.98.7'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := serve(other); err == nil || !strings.Contains(err.Error(), "--release v0.98.7 --candidate") {
		t.Fatalf("un-nominated third release = %v, want a refusal with its nominate command", err)
	}
}

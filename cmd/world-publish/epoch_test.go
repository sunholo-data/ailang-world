package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/pins"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Row 139 M2, AC2: `world-publish epoch nominate`. The candidate release is a
// TEST-ONLY manifest row (the real pins/ailang-releases.json is not changed
// by this milestone): the compiled-in manifest plus one observed release
// carrying invented, well-formed digests.
const (
	testBaseLine      = "AILANG v0.98.0"
	testCandidate     = "v0.98.1"
	testCandidateLine = "AILANG v0.98.1"
)

var (
	testBinaryDigest  = strings.Repeat("ab", 32)
	testTarballDigest = strings.Repeat("cd", 32)
)

// writeTestManifest writes the compiled-in manifest plus the test candidate
// row and returns its path. pins.Parse validates it exactly like the real one.
func writeTestManifest(t *testing.T) string {
	t.Helper()
	return writeTestManifestFor(t, testCandidate, testCandidateLine, testBinaryDigest, testTarballDigest)
}

// writeTestManifestFor is writeTestManifest for any release row; a release
// the compiled-in manifest already lists is not added twice.
func writeTestManifestFor(t *testing.T, release, line, binaryDigest, tarballDigest string) string {
	t.Helper()
	m := pins.Current()
	assets := map[string]pins.Asset{}
	for _, plat := range m.Platforms {
		name := "darwin.arm64.ailang.tar.gz"
		if plat == "linux/amd64" {
			name = "linux.x64.ailang.tar.gz"
		}
		assets[plat] = pins.Asset{Asset: name, TarballSHA256: tarballDigest, BinarySHA256: binaryDigest}
	}
	if _, listed := m.Lookup(release); !listed {
		m.Releases = append(m.Releases, pins.Release{Release: release, Version: line,
			Roles:  []pins.Role{{Role: pins.RoleObserved, Since: "2026-10-16", Evidence: "row 139 M2 test-only manifest row"}},
			Assets: assets})
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pins.Parse(b); err != nil {
		t.Fatalf("instrument: the test manifest does not parse: %v", err)
	}
	path := filepath.Join(t.TempDir(), "ailang-releases.json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// testVerdict returns a passing verdict for the test candidate on this host.
func testVerdict() map[string]any {
	return map[string]any{
		"schema": pins.VerdictSchema, "release": testCandidate, "platform": hostPlatform,
		"binary_sha256": testBinaryDigest, "tarball_sha256": testTarballDigest,
		"golden_identical": true,
		"checks": []any{
			map[string]any{"name": "verify_ail", "pass": true},
			map[string]any{"name": "replay_goldens", "pass": true},
		},
		"all_pass": true, "produced_at": "2026-10-16T12:00:00Z",
	}
}

func writeVerdict(t *testing.T, mutate func(map[string]any)) string {
	t.Helper()
	v := testVerdict()
	if mutate != nil {
		mutate(v)
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "verdict.json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// noTTYProbe is the agent-harness situation: /dev/tty cannot be opened.
func noTTYProbe(t *testing.T) ttyProbe {
	return ttyProbe{stdin: regularFileInfo(t), cttyErr: fmt.Errorf("device not configured")}
}

// nominateFixture is a store bootstrapped on testBaseLine plus the flags of a
// passing candidate nomination of the test candidate.
func nominateFixture(t *testing.T) map[string]string {
	t.Helper()
	storePath := filepath.Join(t.TempDir(), "world.db")
	bootstrapEpochRegistry(t, storePath, testBaseLine)
	return map[string]string{
		"store": storePath, "release": testCandidate,
		"gate-verdict": writeVerdict(t, nil), "release-manifest": writeTestManifest(t),
	}
}

func nominateInv(flags map[string]string, bools ...string) invocation {
	return invocation{verb: "epoch", flags: flags, bools: bools}
}

// args of an `epoch` invocation: the sub-verb goes before the flags.
func (inv invocation) epochArgs(t *testing.T) []string {
	args := inv.args(t)
	return append([]string{"epoch", "nominate"}, args[1:]...)
}

func driveNominate(t *testing.T, inv invocation, stdin string, getenv func(string) string, probe ttyProbe) armResult {
	t.Helper()
	var out, errw strings.Builder
	code := run(inv.epochArgs(t), strings.NewReader(stdin), &out, &errw,
		environment{getenv: getenv, probe: func() ttyProbe { return probe }})
	return armResult{code: code, stdout: out.String(), stderr: errw.String()}
}

func storeRegistry(t *testing.T, storePath string) registry.Registry {
	t.Helper()
	db, err := store.Open(storePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	reg, _, ok, err := registry.Load(boundedTestContext(t), db)
	if err != nil || !ok {
		t.Fatalf("registry.Load: ok=%v err=%v", ok, err)
	}
	return reg
}

// D-WORLD-78 = a: --candidate runs HEADLESS. The probe here cannot open
// /dev/tty and stdin is empty, and the nomination still lands.
func TestEpochNominateCandidateIsHeadless(t *testing.T) {
	flags := nominateFixture(t)
	res := driveNominate(t, nominateInv(flags, "candidate"), "", noEnv, noTTYProbe(t))
	if res.code != exitOK {
		t.Fatalf("headless candidate = (%d, %q), want exit 0", res.code, res.stderr)
	}
	reg := storeRegistry(t, flags["store"])
	if len(reg.Epochs) != 1 || fmt.Sprint(reg.Epochs[0].Candidates) != fmt.Sprint([]string{testBaseLine, testCandidateLine}) {
		t.Fatalf("registry after candidate = %+v", reg.Epochs)
	}
	if !strings.Contains(res.stdout, "registry revision head sha256:") ||
		!strings.Contains(res.stdout, fmt.Sprintf("epoch 1 candidates %q", []string{testBaseLine, testCandidateLine})) {
		t.Fatalf("stdout = %q, want the new head and the epoch/candidates", res.stdout)
	}
	if strings.Contains(res.stdout, "Type exactly") {
		t.Fatalf("a candidate nomination prompted for a phrase: %q", res.stdout)
	}
	// Re-running is refused (already nominated), nothing written.
	again := driveNominate(t, nominateInv(flags, "candidate"), "", noEnv, noTTYProbe(t))
	if again.code != exitStop || !strings.Contains(again.stderr, "STOP fence=registry reason=already-nominated") {
		t.Fatalf("second nomination = (%d, %q), want STOP already-nominated", again.code, again.stderr)
	}
}

// --new-epoch is a human's act: refused without a controlling terminal and in
// CI, refused on a mistyped phrase, and lands only with the exact phrase.
func TestEpochNominateNewEpochIsAttended(t *testing.T) {
	flags := nominateFixture(t)
	golden := func(v map[string]any) { v["golden_identical"] = false }
	flags["gate-verdict"] = writeVerdict(t, golden)
	inv := nominateInv(flags, "new-epoch")
	phrase := newEpochPhrase(2, testCandidate, flags["store"])
	if phrase != "open epoch 2 for "+testCandidate+" in "+flags["store"] {
		t.Fatalf("instrument: phrase oracle %q", phrase)
	}

	if res := driveNominate(t, inv, phrase+"\n", noEnv, noTTYProbe(t)); res.code != exitStop ||
		!strings.Contains(res.stderr, "STOP fence=tty reason=no-controlling-terminal") {
		t.Fatalf("new epoch without a TTY = (%d, %q), want STOP fence=tty", res.code, res.stderr)
	}
	if res := driveNominate(t, inv, phrase+"\n", ciEnv, satisfiedProbe(t)); res.code != exitStop ||
		!strings.Contains(res.stderr, "STOP fence=ci") {
		t.Fatalf("new epoch in CI = (%d, %q), want STOP fence=ci", res.code, res.stderr)
	}
	for _, wrong := range []string{"yes", newEpochPhrase(3, testCandidate, flags["store"]), newEpochPhrase(2, testCandidate, flags["store"]+".x")} {
		if res := driveNominate(t, inv, wrong+"\n", noEnv, satisfiedProbe(t)); res.code != exitStop ||
			!strings.Contains(res.stderr, "fence=confirmation reason=mismatch") {
			t.Fatalf("phrase %q = (%d, %q), want STOP mismatch", wrong, res.code, res.stderr)
		}
	}
	if reg := storeRegistry(t, flags["store"]); len(reg.Epochs) != 1 {
		t.Fatalf("a refused new epoch wrote: %+v", reg.Epochs)
	}
	res := driveNominate(t, inv, phrase+"\n", noEnv, satisfiedProbe(t))
	if res.code != exitOK {
		t.Fatalf("attended new epoch = (%d, %q / %q), want exit 0", res.code, res.stdout, res.stderr)
	}
	reg := storeRegistry(t, flags["store"])
	if len(reg.Epochs) != 2 || reg.Epochs[1].Epoch != 2 || fmt.Sprint(reg.Epochs[1].Candidates) != fmt.Sprint([]string{testCandidateLine}) {
		t.Fatalf("registry after new epoch = %+v", reg.Epochs)
	}
	if !strings.Contains(res.stdout, fmt.Sprintf("epoch 2 candidates %q", []string{testCandidateLine})) {
		t.Fatalf("stdout = %q", res.stdout)
	}
}

// Every refusal is a STOP (exit 3) and writes nothing; each arm changes ONE
// thing from the passing fixture (TestEpochNominateCandidateIsHeadless is the
// control).
func TestEpochNominateRefusals(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(t *testing.T, f map[string]string) []string // returns the bools
		want   string
	}{
		{"no verdict", func(t *testing.T, f map[string]string) []string {
			delete(f, "gate-verdict")
			return []string{"candidate"}
		},
			"STOP fence=verdict reason=absent"},
		{"unreadable verdict", func(t *testing.T, f map[string]string) []string {
			f["gate-verdict"] = filepath.Join(t.TempDir(), "absent.json")
			return []string{"candidate"}
		}, "STOP fence=verdict reason=unreadable"},
		{"malformed verdict", func(t *testing.T, f map[string]string) []string {
			f["gate-verdict"] = writeVerdict(t, func(v map[string]any) { v["schema"] = "nope" })
			return []string{"candidate"}
		}, "STOP fence=verdict reason=malformed"},
		{"verdict for another release", func(t *testing.T, f map[string]string) []string {
			f["gate-verdict"] = writeVerdict(t, func(v map[string]any) { v["release"] = "v0.98.2" })
			return []string{"candidate"}
		}, "STOP fence=verdict reason=release-mismatch"},
		{"binary digest mismatch", func(t *testing.T, f map[string]string) []string {
			f["gate-verdict"] = writeVerdict(t, func(v map[string]any) { v["binary_sha256"] = strings.Repeat("ef", 32) })
			return []string{"candidate"}
		}, "STOP fence=verdict reason=digest-mismatch"},
		{"not all pass", func(t *testing.T, f map[string]string) []string {
			f["gate-verdict"] = writeVerdict(t, func(v map[string]any) {
				v["checks"] = []any{map[string]any{"name": "replay_goldens", "pass": false}}
				v["all_pass"] = false
			})
			return []string{"candidate"}
		}, "STOP fence=verdict reason=not-all-pass"},
		{"candidate whose goldens differ", func(t *testing.T, f map[string]string) []string {
			f["gate-verdict"] = writeVerdict(t, func(v map[string]any) { v["golden_identical"] = false })
			return []string{"candidate"}
		}, "STOP fence=verdict reason=goldens-differ"},
		{"unknown release", func(t *testing.T, f map[string]string) []string {
			f["release"] = "v0.98.9"
			f["gate-verdict"] = writeVerdict(t, func(v map[string]any) { v["release"] = "v0.98.9" })
			return []string{"candidate"}
		}, "STOP fence=release reason=unknown"},
		{"release in the compiled-in manifest has no candidate row", func(t *testing.T, f map[string]string) []string {
			delete(f, "release-manifest")
			return []string{"candidate"}
		}, "not in the release manifest pins/ailang-releases.json (compiled in)"},
		{"no release flag", func(t *testing.T, f map[string]string) []string { delete(f, "release"); return []string{"candidate"} },
			"STOP fence=release reason=absent"},
		{"no mode", func(t *testing.T, f map[string]string) []string { return nil }, "STOP fence=mode reason=none"},
		{"both modes", func(t *testing.T, f map[string]string) []string { return []string{"candidate", "new-epoch"} },
			"STOP fence=mode reason=both"},
		{"--live", func(t *testing.T, f map[string]string) []string { return []string{"candidate", "live"} },
			"STOP fence=mode reason=epoch nominate-is-read-only"},
		{"no store flag", func(t *testing.T, f map[string]string) []string { delete(f, "store"); return []string{"candidate"} },
			"STOP fence=store reason=absent"},
		{"store with no registry head", func(t *testing.T, f map[string]string) []string {
			f["store"] = filepath.Join(t.TempDir(), "fresh.db")
			return []string{"candidate"}
		}, "STOP fence=registry reason=absent"},
		{"new epoch without a TTY", func(t *testing.T, f map[string]string) []string { return []string{"new-epoch"} },
			"STOP fence=tty reason=no-controlling-terminal"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := nominateFixture(t)
			storePath := f["store"]
			bools := c.mutate(t, f)
			res := driveNominate(t, nominateInv(f, bools...), "", noEnv, noTTYProbe(t))
			if res.code != exitStop || !strings.Contains(res.stderr, c.want) {
				t.Fatalf("= (%d, %q), want STOP containing %q", res.code, res.stderr, c.want)
			}
			if reg := storeRegistry(t, storePath); len(reg.Epochs) != 1 || len(reg.Epochs[0].Candidates) != 1 {
				t.Fatalf("a refused nomination wrote: %+v", reg.Epochs)
			}
		})
	}
	// The unknown-release refusal names the manifest it consulted.
	f := nominateFixture(t)
	f["release"] = "v0.98.9"
	res := driveNominate(t, nominateInv(f, "candidate"), "", noEnv, noTTYProbe(t))
	if !strings.Contains(res.stderr, f["release-manifest"]) {
		t.Fatalf("unknown release refusal %q does not name the manifest %s", res.stderr, f["release-manifest"])
	}
}

func TestEpochNeedsTheNominateSubVerb(t *testing.T) {
	for _, args := range [][]string{{"epoch"}, {"epoch", "--store", "x"}, {"epoch", "grow"}} {
		var out, errw strings.Builder
		if code := run(args, strings.NewReader(""), &out, &errw, environment{getenv: noEnv, probe: func() ttyProbe { return noTTYProbe(t) }}); code != exitUsage {
			t.Fatalf("%v = %d (%q), want usage", args, code, errw.String())
		}
	}
}

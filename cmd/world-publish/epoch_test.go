package main

import (
	"crypto/sha256"
	"encoding/hex"
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
// whose binary digest is the sha256 of a shell fake. The test manifest
// reaches the verb ONLY through runEpochNominateWith (no flag names a
// manifest); arms that go through run() use the compiled-in manifest.
const (
	testBaseLine      = "AILANG v0.98.0"
	testCandidate     = "v0.98.1"
	testCandidateLine = "AILANG v0.98.1"
	// testManifestKey is a PSEUDO-flag of the test driver, never passed to
	// argv: it names the test manifest handed to runEpochNominateWith.
	testManifestKey = "\x00test-manifest"
)

var testTarballDigest = strings.Repeat("cd", 32)

// fakeScript is a shell fake whose --version prints line.
func fakeScript(line string) string {
	return "#!/bin/sh\nif [ \"$1\" = --version ]; then echo '" + line + "'; exit 0; fi\nexit 0\n"
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// writeScript writes an executable script and returns its path and sha256.
func writeScript(t *testing.T, body string) (string, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ailang")
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path, sha256Hex([]byte(body))
}

// writeTestManifestFor writes the compiled-in manifest plus one release row
// (unless the compiled-in manifest already lists it) and returns its path.
// pins.Parse validates it exactly like the real one.
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

// verdictFor returns a passing verdict for release with binaryDigest.
func verdictFor(release, binaryDigest string) map[string]any {
	return map[string]any{
		"schema": pins.VerdictSchema, "release": release, "platform": hostPlatform,
		"binary_sha256": binaryDigest, "tarball_sha256": testTarballDigest,
		"golden_identical": true,
		"checks": []any{
			map[string]any{"name": "verify_ail", "pass": true},
			map[string]any{"name": "replay_goldens", "pass": true},
		},
		"all_pass": true, "produced_at": "2026-10-16T12:00:00Z",
	}
}

// testBinaryDigest is the sha256 of the test candidate's fake binary.
var testBinaryDigest = sha256Hex([]byte(fakeScript(testCandidateLine)))

func writeVerdict(t *testing.T, mutate func(map[string]any)) string {
	t.Helper()
	return writeVerdictMap(t, verdictFor(testCandidate, testBinaryDigest), mutate)
}

func writeVerdictMap(t *testing.T, v map[string]any, mutate func(map[string]any)) string {
	t.Helper()
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
// passing candidate nomination of the test candidate (fake binary, test
// manifest naming its digest, verdict echoing it).
func nominateFixture(t *testing.T) map[string]string {
	t.Helper()
	storePath := filepath.Join(t.TempDir(), "world.db")
	bootstrapEpochRegistry(t, storePath, testBaseLine)
	bin, digest := writeScript(t, fakeScript(testCandidateLine))
	if digest != testBinaryDigest {
		t.Fatal("instrument: fake binary digest drifted")
	}
	return map[string]string{
		"store": storePath, "release": testCandidate, "binary": bin,
		"gate-verdict":  writeVerdict(t, nil),
		testManifestKey: writeTestManifestFor(t, testCandidate, testCandidateLine, testBinaryDigest, testTarballDigest),
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

// driveNominate runs the verb. With testManifestKey set it parses the other
// flags with the production FlagSet and calls runEpochNominateWith with that
// manifest (the test seam); without it, it goes through run(), i.e. argv
// and the compiled-in manifest.
func driveNominate(t *testing.T, inv invocation, stdin string, getenv func(string) string, probe ttyProbe) armResult {
	t.Helper()
	var out, errw strings.Builder
	env := environment{getenv: getenv, probe: func() ttyProbe { return probe }}
	manifestPath, seam := inv.flags[testManifestKey]
	if !seam {
		code := run(inv.epochArgs(t), strings.NewReader(stdin), &out, &errw, env)
		return armResult{code: code, stdout: out.String(), stderr: errw.String()}
	}
	argv := inv.with(func(i *invocation) { delete(i.flags, testManifestKey) }).epochArgs(t)
	var opts options
	fs := newFlagSet(&opts)
	if err := fs.Parse(argv[2:]); err != nil {
		t.Fatalf("instrument: test argv does not parse: %v", err)
	}
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	m, err := pins.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	code := runEpochNominateWith(opts, m, manifestPath, strings.NewReader(stdin), &out, &errw, env)
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
		}, "STOP fence=verdict reason=absent"},
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
		{"verdict digest mismatch", func(t *testing.T, f map[string]string) []string {
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
		{"argv uses the compiled-in manifest, which has no row for the test release", func(t *testing.T, f map[string]string) []string {
			delete(f, testManifestKey)
			return []string{"candidate"}
		}, "not in the release manifest pins/ailang-releases.json (compiled in)"},
		{"no binary", func(t *testing.T, f map[string]string) []string { delete(f, "binary"); return []string{"candidate"} },
			"STOP fence=binary reason=absent"},
		{"unreadable binary", func(t *testing.T, f map[string]string) []string {
			f["binary"] = filepath.Join(t.TempDir(), "no-such-ailang")
			return []string{"candidate"}
		}, "STOP fence=binary reason=unreadable"},
		{"binary whose bytes differ from the manifest", func(t *testing.T, f map[string]string) []string {
			f["binary"], _ = writeScript(t, fakeScript(testCandidateLine)+"# one more byte\n")
			return []string{"candidate"}
		}, "STOP fence=binary reason=digest-mismatch"},
		{"binary whose version line is another release", func(t *testing.T, f map[string]string) []string {
			// Manifest and verdict name THIS binary's digest; only its version differs.
			bin, digest := writeScript(t, fakeScript("AILANG v0.98.2"))
			f["binary"] = bin
			f[testManifestKey] = writeTestManifestFor(t, testCandidate, testCandidateLine, digest, testTarballDigest)
			f["gate-verdict"] = writeVerdictMap(t, verdictFor(testCandidate, digest), nil)
			return []string{"candidate"}
		}, "STOP fence=binary reason=version-mismatch"},
		{"binary whose version probe fails", func(t *testing.T, f map[string]string) []string {
			bin, digest := writeScript(t, "#!/bin/sh\nexit 7\n")
			f["binary"] = bin
			f[testManifestKey] = writeTestManifestFor(t, testCandidate, testCandidateLine, digest, testTarballDigest)
			f["gate-verdict"] = writeVerdictMap(t, verdictFor(testCandidate, digest), nil)
			return []string{"candidate"}
		}, "STOP fence=binary reason=version-probe-failed"},
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
	if !strings.Contains(res.stderr, f[testManifestKey]) {
		t.Fatalf("unknown release refusal %q does not name the manifest %s", res.stderr, f[testManifestKey])
	}
}

// Judge r1 B1: no flag can name another manifest. --release-manifest (and any
// other spelling) is an unknown flag: usage, nothing written.
func TestEpochNominateHasNoManifestFlag(t *testing.T) {
	f := nominateFixture(t)
	for _, name := range []string{"release-manifest", "manifest-file", "pins"} {
		inv := nominateInv(map[string]string{"store": f["store"], "release": testCandidate, "binary": f["binary"],
			"gate-verdict": f["gate-verdict"], name: f[testManifestKey]}, "candidate")
		res := driveNominate(t, inv, "", noEnv, noTTYProbe(t))
		if res.code != exitUsage || !strings.Contains(res.stderr, "flag provided but not defined: -"+name) {
			t.Fatalf("--%s = (%d, %q), want usage: flag not defined", name, res.code, res.stderr)
		}
	}
	// --manifest exists (it is the transitions descriptor manifest) but is NOT
	// read by epoch nominate: naming the test manifest there changes nothing,
	// the compiled-in manifest still refuses the test release.
	inv := nominateInv(map[string]string{"store": f["store"], "release": testCandidate, "binary": f["binary"],
		"gate-verdict": f["gate-verdict"], "manifest": f[testManifestKey]}, "candidate")
	res := driveNominate(t, inv, "", noEnv, noTTYProbe(t))
	if res.code != exitStop || !strings.Contains(res.stderr, "pins/ailang-releases.json (compiled in)") {
		t.Fatalf("--manifest = (%d, %q), want the compiled-in manifest's unknown-release STOP", res.code, res.stderr)
	}
	for _, name := range flagNames {
		if strings.Contains(name, "release-manifest") {
			t.Fatalf("frozen flag surface names %q", name)
		}
	}
	if reg := storeRegistry(t, f["store"]); len(reg.Epochs[0].Candidates) != 1 {
		t.Fatalf("a refused nomination wrote: %+v", reg.Epochs)
	}
}

// Through argv (the COMPILED-IN manifest), against a release the compiled
// manifest really lists — the interpreter release, served by AILANG_BIN:
//   - a verdict whose digests differ from the compiled manifest is refused;
//   - a FORGED verdict echoing the compiled digests, with a binary that is
//     not the release's bytes, is refused at --binary;
//   - CONTROL: the real pinned binary with that verdict nominates.
func TestEpochNominateAgainstTheCompiledManifest(t *testing.T) {
	pin := os.Getenv("AILANG_BIN")
	if pin == "" {
		t.Fatal("AILANG_BIN unset: the compiled-manifest control needs the pinned released interpreter; never skip")
	}
	rel := pins.ReleaseFor(pins.RoleInterpreter)
	asset, ok := rel.Assets[hostPlatform]
	if !ok {
		t.Fatalf("the compiled manifest has no %s asset for %s", hostPlatform, rel.Release)
	}
	compiled := verdictFor(rel.Release, asset.BinarySHA256)
	compiled["tarball_sha256"] = asset.TarballSHA256
	fresh := func() string {
		storePath := filepath.Join(t.TempDir(), "world.db")
		bootstrapEpochRegistry(t, storePath, testBaseLine)
		return storePath
	}
	forged, _ := writeScript(t, fakeScript(rel.Version))

	for _, c := range []struct {
		name, binary string
		mutate       func(map[string]any)
		want         string
	}{
		{"verdict digests differ from the compiled manifest", pin,
			func(v map[string]any) { v["binary_sha256"] = strings.Repeat("ef", 32) }, "STOP fence=verdict reason=digest-mismatch"},
		{"tarball digest differs from the compiled manifest", pin,
			func(v map[string]any) { v["tarball_sha256"] = strings.Repeat("ef", 32) }, "STOP fence=verdict reason=digest-mismatch"},
		{"forged verdict, binary not the release's bytes", forged, nil, "STOP fence=binary reason=digest-mismatch"},
	} {
		storePath := fresh()
		v := map[string]any{}
		for k, x := range compiled {
			v[k] = x
		}
		inv := nominateInv(map[string]string{"store": storePath, "release": rel.Release, "binary": c.binary,
			"gate-verdict": writeVerdictMap(t, v, c.mutate)}, "candidate")
		res := driveNominate(t, inv, "", noEnv, noTTYProbe(t))
		if res.code != exitStop || !strings.Contains(res.stderr, c.want) {
			t.Fatalf("%s = (%d, %q), want %q", c.name, res.code, res.stderr, c.want)
		}
		if reg := storeRegistry(t, storePath); len(reg.Epochs[0].Candidates) != 1 {
			t.Fatalf("%s wrote: %+v", c.name, reg.Epochs)
		}
	}
	storePath := fresh()
	inv := nominateInv(map[string]string{"store": storePath, "release": rel.Release, "binary": pin,
		"gate-verdict": writeVerdictMap(t, compiled, nil)}, "candidate")
	res := driveNominate(t, inv, "", noEnv, noTTYProbe(t))
	if res.code != exitOK {
		t.Fatalf("control: real pin through argv = (%d, %q), want exit 0", res.code, res.stderr)
	}
	if reg := storeRegistry(t, storePath); fmt.Sprint(reg.Epochs[0].Candidates) != fmt.Sprint([]string{testBaseLine, rel.Version}) {
		t.Fatalf("control registry = %+v", reg.Epochs)
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

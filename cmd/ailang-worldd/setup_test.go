package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/pinfetch"
)

// fakeRelease is an in-memory pinfetch.Source: no network at all.
type fakeRelease struct {
	assets map[string][]byte // "<release>/<asset>" -> bytes
	count  atomic.Int64
}

func (f *fakeRelease) Open(_ context.Context, release, asset string, _ time.Duration) (io.ReadCloser, error) {
	f.count.Add(1)
	b, ok := f.assets[release+"/"+asset]
	if !ok {
		return nil, fmt.Errorf("fake release: no %s/%s", release, asset)
	}
	return io.NopCloser(bytes.NewReader(b)), nil
}
func (f *fakeRelease) Describe(release, asset string) string {
	return "fake://" + release + "/" + asset
}
func (f *fakeRelease) Requests() int64 { return f.count.Load() }

func sha256Hex(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

func fixtureTarball(t *testing.T, binary []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "ailang", Typeflag: tar.TypeReg, Mode: 0o755, Size: int64(len(binary))}); err != nil {
		t.Fatal(err)
	}
	_, _ = tw.Write(binary)
	_ = tw.Close()
	_ = gz.Close()
	return buf.Bytes()
}

// setupFixture swaps the pin table, the source and the platform for a small
// fixture release pair, and HOME for a scratch directory.
type setupFixture struct {
	home     string
	src      *fakeRelease
	binaries map[string][]byte // release -> binary
}

func newSetupFixture(t *testing.T) *setupFixture {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	f := &setupFixture{home: home, src: &fakeRelease{assets: map[string][]byte{}}, binaries: map[string][]byte{}}
	var table []pinfetch.Pin
	for _, rel := range []string{interpreterRelease, toolRelease()} {
		bin := []byte("#!/bin/sh\necho fixture AILANG " + rel + "\n")
		tb := fixtureTarball(t, bin)
		asset := "darwin.arm64.ailang.tar.gz"
		f.src.assets[rel+"/"+asset] = tb
		f.src.assets[rel+"/"+asset+".sha256"] = []byte(sha256Hex(tb) + "  " + asset + "\n")
		f.binaries[rel] = bin
		table = append(table, pinfetch.Pin{Release: rel, Platform: "darwin/arm64", Asset: asset,
			TarballSHA256: sha256Hex(tb), BinarySHA256: sha256Hex(bin)})
	}
	oldTable, oldSrc, oldPlat := pinTable, setupSource, setupPlatform
	pinTable, setupPlatform = table, "darwin/arm64"
	setupSource = func() pinfetch.Source { return f.src }
	t.Cleanup(func() { pinTable, setupSource, setupPlatform = oldTable, oldSrc, oldPlat })
	return f
}

func runSetupCLI(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errw bytes.Buffer
	code := run(append([]string{"setup"}, args...), &out, &errw)
	return code, out.String(), errw.String()
}

// AC3.1 + AC3.6 through the verb: both pins installed with pin.json and no
// tarball, the store dir and workspace root made 0700, and a second run makes
// zero requests.
func TestSetupInstallsBothPinsThenMakesNoRequest(t *testing.T) {
	f := newSetupFixture(t)
	code, out, errOut := runSetupCLI(t)
	if code != exitOK {
		t.Fatalf("setup rc %d\nstdout:\n%s\nstderr:\n%s", code, out, errOut)
	}
	interp := filepath.Join(f.home, ".pinned-ailang", "ailang")
	tool := filepath.Join(f.home, ".pinned-ailang-tools", toolRelease(), "ailang")
	for rel, path := range map[string]string{interpreterRelease: interp, toolRelease(): tool} {
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, f.binaries[rel]) {
			t.Fatalf("%s: installed bytes differ (%v)", path, err)
		}
		ents, _ := os.ReadDir(filepath.Dir(path))
		var names []string
		for _, e := range ents {
			names = append(names, e.Name())
		}
		if strings.Join(names, ",") != "ailang,pin.json" {
			t.Fatalf("%s holds %v, want [ailang pin.json]", filepath.Dir(path), names)
		}
	}
	for _, d := range []string{filepath.Join(f.home, ".ailang", "world"), filepath.Join(f.home, ".ailang", "world-ws")} {
		info, err := os.Stat(d)
		if err != nil || !info.IsDir() || info.Mode().Perm() != 0o700 {
			t.Fatalf("%s: %v %v, want a 0700 directory", d, info, err)
		}
	}
	if _, err := os.Stat(filepath.Join(f.home, ".ailang", "world", "world.db")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("setup created the store itself (stat err %v); it must create only the directory", err)
	}
	for _, want := range []string{"installed " + interp, "installed " + tool, "next (the attended steps", "world-publish transitions --store"} {
		if !strings.Contains(out, want) {
			t.Errorf("setup stdout lacks %q:\n%s", want, out)
		}
	}
	if n := f.src.Requests(); n != 4 {
		t.Fatalf("first run made %d requests, want 4 (two .sha256 + two tarballs)", n)
	}

	code, out, errOut = runSetupCLI(t)
	if code != exitOK {
		t.Fatalf("second setup rc %d: %s", code, errOut)
	}
	if n := f.src.Requests(); n != 4 {
		t.Fatalf("second run made %d more request(s), want 0", n-4)
	}
	if strings.Count(out, "present") != 2 {
		t.Fatalf("second run should report both pins present:\n%s", out)
	}
}

// AC3.2 through the verb: a consistently tampered release exits 3 and
// installs nothing.
func TestSetupDigestMismatchExits3AndInstallsNothing(t *testing.T) {
	f := newSetupFixture(t)
	evil := fixtureTarball(t, []byte("trojan"))
	asset := "darwin.arm64.ailang.tar.gz"
	f.src.assets[interpreterRelease+"/"+asset] = evil
	f.src.assets[interpreterRelease+"/"+asset+".sha256"] = []byte(sha256Hex(evil) + "\n")
	code, out, errOut := runSetupCLI(t)
	if code != exitIntegrity {
		t.Fatalf("rc %d, want %d\n%s%s", code, exitIntegrity, out, errOut)
	}
	if !strings.Contains(errOut, sha256Hex(evil)) || !strings.Contains(errOut, "compiled-in pin") {
		t.Fatalf("refusal does not name the digests:\n%s", errOut)
	}
	if _, err := os.Stat(filepath.Join(f.home, ".pinned-ailang", "ailang")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a refused interpreter was installed (stat err %v)", err)
	}
}

func TestSetupRefusesUnsupportedPlatformBeforeAnyRequest(t *testing.T) {
	f := newSetupFixture(t)
	setupPlatform = "windows/amd64"
	code, _, errOut := runSetupCLI(t)
	if code != exitUsage || !strings.Contains(errOut, "windows/amd64 is not supported") || !strings.Contains(errOut, "darwin/arm64 and linux/amd64") {
		t.Fatalf("rc %d stderr %q", code, errOut)
	}
	if f.src.Requests() != 0 {
		t.Fatalf("an unsupported platform still made %d request(s)", f.src.Requests())
	}
}

func TestSetupRefusesAWorkspaceRootThatContainsTheStore(t *testing.T) {
	f := newSetupFixture(t)
	ws := filepath.Join(f.home, "ws")
	code, _, errOut := runSetupCLI(t, "--workspace-root", ws, "--db", filepath.Join(ws, "store", "world.db"))
	if code != exitUsage || !strings.Contains(errOut, "contains the store directory") {
		t.Fatalf("rc %d stderr %q", code, errOut)
	}
	if f.src.Requests() != 0 {
		t.Fatalf("a refused layout still made %d request(s)", f.src.Requests())
	}
}

func TestSetupFromDirInstallsOffline(t *testing.T) {
	f := newSetupFixture(t)
	dir := t.TempDir()
	asset := "darwin.arm64.ailang.tar.gz"
	for _, rel := range []string{interpreterRelease, toolRelease()} {
		_ = os.MkdirAll(filepath.Join(dir, rel), 0o755)
		_ = os.WriteFile(filepath.Join(dir, rel, asset), f.src.assets[rel+"/"+asset], 0o644)
		_ = os.WriteFile(filepath.Join(dir, rel, asset+".sha256"), f.src.assets[rel+"/"+asset+".sha256"], 0o644)
	}
	code, out, errOut := runSetupCLI(t, "--from-dir", dir)
	if code != exitOK {
		t.Fatalf("rc %d\n%s%s", code, out, errOut)
	}
	if f.src.Requests() != 0 {
		t.Fatalf("--from-dir made %d network-source request(s)", f.src.Requests())
	}
	got, _ := os.ReadFile(filepath.Join(f.home, ".pinned-ailang-tools", toolRelease(), "ailang"))
	if !bytes.Equal(got, f.binaries[toolRelease()]) {
		t.Fatal("--from-dir did not install the tool pin from <dir>/<release>/")
	}
}

func TestSetupHelpExitsZero(t *testing.T) {
	code, out, _ := runSetupCLI(t, "--help")
	if code != exitOK || !strings.Contains(out, "usage: ailang-worldd setup") {
		t.Fatalf("rc %d out %q", code, out)
	}
}

// AC3.7: the compiled table covers daemon.ToolBinaryRelease and the
// interpreter release on darwin/arm64 and linux/amd64, and the interpreter
// release and digests are bound to the other places that pin them.
func TestPinTableCoversBothReleases(t *testing.T) {
	if len(pinTable) != 4 {
		t.Fatalf("pin table has %d rows, want 4 (2 releases × 2 platforms)", len(pinTable))
	}
	hex64 := regexp.MustCompile(`^[0-9a-f]{64}$`)
	for _, rel := range []string{interpreterRelease, toolRelease()} {
		for _, plat := range supportedPlatforms {
			p, err := lookupPin(rel, plat)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if !hex64.MatchString(p.TarballSHA256) || !hex64.MatchString(p.BinarySHA256) {
				t.Errorf("%s %s: malformed digests %+v", rel, plat, p)
			}
			wantAsset := map[string]string{"darwin/arm64": "darwin.arm64.ailang.tar.gz", "linux/amd64": "linux.x64.ailang.tar.gz"}[plat]
			if p.Asset != wantAsset {
				t.Errorf("%s %s: asset %q, want %q", rel, plat, p.Asset, wantAsset)
			}
		}
	}
	if fmt.Sprint(supportedPlatforms) != "[darwin/arm64 linux/amd64]" {
		t.Fatalf("supported platforms %v (D-CLI-2)", supportedPlatforms)
	}

	root := filepath.Join("..", "..")
	read := func(rel string) string {
		b, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	if !strings.Contains(read("cmd/world-publish/fences.go"), `const frozenCompilerVersion = "AILANG `+interpreterRelease+`"`) {
		t.Errorf("interpreterRelease %s is not world-publish's frozenCompilerVersion", interpreterRelease)
	}
	ci := read(".github/workflows/ci.yml")
	for _, rel := range []string{interpreterRelease, toolRelease()} {
		if !strings.Contains(ci, "releases/download/"+rel+`"`) {
			t.Errorf("ci.yml installs no %s release", rel)
		}
	}
	pkg := read("scripts/verify_world_package.sh")
	for plat, uname := range map[string]string{"darwin/arm64": "Darwin/arm64", "linux/amd64": "Linux/x86_64"} {
		p, _ := lookupPin(interpreterRelease, plat)
		if !strings.Contains(pkg, uname+") printf '%s' '"+p.BinarySHA256+"'") {
			t.Errorf("verify_world_package.sh does not pin %s's interpreter binary %s", uname, p.BinarySHA256)
		}
	}
}

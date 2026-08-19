package pkgproj

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoldenHashes(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "b.ail"), "beta\n")
	mustWrite(t, filepath.Join(dir, "sub", "a.ail"), "alpha\n")
	mustWrite(t, filepath.Join(dir, "ignored.txt"), "ignored")
	got, err := ContentHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if want := "sha256:7219158e2cded9307b0583e018607b6c853d1a98eb9f5893661abec73587dca7"; got != want {
		t.Fatalf("ContentHash=%s want %s", got, want)
	}
	m := Manifest{Package: Package{Name: "world/core", Edition: "1", AILANG: ">=0.30.0"}, Exports: Exports{Modules: []string{"world/z", "world/a"}}, Effects: Effects{Max: []string{"IO", "FS"}}}
	if got, want := InterfaceHash(m), "sha256:a58f67ff3ec5579b8fe9d5eb2b2a818b29d8773e7662f51d9a7927240a40958e"; got != want {
		t.Fatalf("InterfaceHash=%s want %s", got, want)
	}
}

func TestGoldenTarballBytes(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "ailang.toml"), "manifest\n")
	mustWrite(t, filepath.Join(dir, "world", "x.ail"), "export x\n")
	mustWrite(t, filepath.Join(dir, "tests", "bad.ail"), "bad")
	data, err := CreateTarball(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := TarballHash(data), "sha256:97ee99e4b051d7f85032dc5dde1ea0b57411068acdcbc82da69e0231ff13b3f1"; got != want {
		t.Fatalf("TarballHash=%s want %s (length=%d)", got, want, len(data))
	}
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(gr)
	var names []string
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, h.Name)
	}
	if got := strings.Join(names, ","); got != "ailang.toml,world/x.ail" {
		t.Fatalf("entries=%s", got)
	}
}

func TestCompareFailsLoudlyOnEveryDisagreement(t *testing.T) {
	local := Hashes{Content: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Interface: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", Tarball: "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc", TarballBytes: 101}
	cli := Hashes{Content: "sha256:11111111111111111", Interface: "sha256:22222222222222222", Tarball: "sha256:33333333333333333", TarballBytes: 99}
	err := Compare(local, cli)
	if err == nil {
		t.Fatal("Compare returned nil for three disagreements")
	}
	for _, value := range []string{"content", local.Content, cli.Content, "interface", local.Interface, cli.Interface, "tarball", local.Tarball, cli.Tarball, "101", "99"} {
		if !strings.Contains(err.Error(), value) {
			t.Errorf("error %q does not name %q", err, value)
		}
	}
}

func mustWrite(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// w-archive-stderr-in-manifest (queue row 21), T3: the FIRST test of this
// package's exec seam. Before this, CrossCheck spawned a subprocess that no
// test had ever run.
//
// It demonstrates the sharper half of the stderr-merge hazard: merged stderr
// does not merely add noise to a diagnostic, it can INJECT PARSEABLE DATA. A
// stderr line in dry-run SHAPE is indistinguishable from a stdout one once the
// two fds are the same pipe, and parseDryRun rejects the second "Tarball:" line
// it sees. The line below is matched against dryRunLine's own anchors -- two
// leading spaces, exactly 17 lowercase hex -- so it is a real second data line,
// not decorative noise a laxer parser would ignore.
// ---------------------------------------------------------------------------

// stderrDryRunLine is a byte-for-byte VALID dry-run Tarball line emitted on
// fd 2. Its 17 hex digits and its "999 bytes" are deliberately wrong: if the
// merge is restored, parseDryRun sees it first and errors on the duplicate
// rather than silently believing it -- but either way the test reds, which is
// the point.
const stderrDryRunLine = "  Tarball: 999 bytes (sha256:abcdef0123456789a...)"

// cliPrefix is how the AILANG CLI abbreviates a hash in its dry-run block:
// "sha256:" plus the first 17 digits, which Compare checks by HasPrefix against
// the full local hash.
func cliPrefix(t *testing.T, full string) string {
	t.Helper()
	const want = len("sha256:") + 17
	if len(full) < want {
		t.Fatalf("hash %q is too short to abbreviate", full)
	}
	return full[:want]
}

// fakeAilangCLI writes a stand-in `ailang` that prints a well-formed dry-run
// block on STDOUT and one regex-matching line on STDERR, sequentially from one
// shell.
//
// Two fixture rules are load-bearing:
//   - dir MUST NOT be the package directory. CrossCheck computes
//     ContentHash(packageDir) and CreateTarball(packageDir) over that tree, so a
//     script dropped inside it would change the very hashes it echoes back.
//   - the returned path is absolute, because CrossCheck sets cmd.Dir =
//     packageDir and a relative binary path would resolve against it.
//
// Sequential writes from ONE shell also matter: under CombinedOutput() both fds
// are the same pipe, so the writes cannot interleave mid-line. A split line
// would fail dryRunLine's (?m)^…$ anchors, parseDryRun would see only ONE
// Tarball line, and the M2 red would silently not fire.
func fakeAilangCLI(t *testing.T, dir string, stdoutBlock, stderrLine string) string {
	t.Helper()
	script := "#!/bin/sh\n" +
		"printf '%s\\n' '" + stderrLine + "' 1>&2\n" +
		"printf '%s' '" + stdoutBlock + "'\n" +
		"exit 0\n"
	p := filepath.Join(dir, "ailang-fake")
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake ailang CLI: %v", err)
	}
	if !filepath.IsAbs(p) {
		t.Fatalf("fake CLI path %q is not absolute; CrossCheck sets cmd.Dir", p)
	}
	return p
}

// AC3. RED MUTATION M2: restore `out, err := cmd.CombinedOutput()` at the
// CrossCheck exec and parseDryRun sees TWO Tarball lines, returning
// "duplicate dry-run Tarball line"; this test asserts success, so it reds.
func TestCrossCheckStderrIsNotParsedAsDryRunData(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake CLI is POSIX-only")
	}
	// The package tree the hashes are computed over.
	pkgDir := t.TempDir()
	mustWrite(t, filepath.Join(pkgDir, "ailang.toml"), "[package]\nname = \"world/stderrcheck\"\n")
	mustWrite(t, filepath.Join(pkgDir, "world", "x.ail"), "export func x() -> int { 1 }\n")

	manifest := Manifest{
		Package: Package{Name: "world/stderrcheck", Edition: "1", AILANG: ">=0.30.0"},
		Exports: Exports{Modules: []string{"world/x"}},
		Effects: Effects{Max: []string{"IO"}},
	}

	// Compute the truth with the package's own exported functions, so the fake
	// CLI "agrees" for exactly the reason the real one would.
	content, err := ContentHash(pkgDir)
	if err != nil {
		t.Fatalf("ContentHash: %v", err)
	}
	tarball, err := CreateTarball(pkgDir)
	if err != nil {
		t.Fatalf("CreateTarball: %v", err)
	}
	tarballHash := TarballHash(tarball)
	ifaceHash := InterfaceHash(manifest)

	stdoutBlock := "Dry run for " + manifest.Package.Name + "\n" +
		"  Tarball: " + strconv.Itoa(len(tarball)) + " bytes (" + cliPrefix(t, tarballHash) + "...)\n" +
		"  Content hash: " + cliPrefix(t, content) + "...\n" +
		"  Interface hash: " + cliPrefix(t, ifaceHash) + "...\n"

	// Fixture control: the line we are about to emit on stderr must actually
	// satisfy the production parser's regex. A "stderr line" the parser ignores
	// would make the M2 red vacuous -- the mutation would change nothing.
	if !dryRunLine.MatchString(stderrDryRunLine + "\n") {
		t.Fatalf("control: stderr fixture %q does not match dryRunLine, so it could never be parsed as data", stderrDryRunLine)
	}
	// Same control for the stdout block, which must yield all three labels.
	if got := len(dryRunLine.FindAllStringSubmatch(stdoutBlock, -1)); got != 3 {
		t.Fatalf("control: stdout block yielded %d dry-run lines, want 3:\n%s", got, stdoutBlock)
	}

	// The fake lives OUTSIDE pkgDir (C7) so it cannot perturb the hashes.
	binDir := t.TempDir()
	fakeBin := fakeAilangCLI(t, binDir, stdoutBlock, stderrDryRunLine)

	result, err := CrossCheck(pkgDir, manifest, fakeBin)
	if err != nil {
		t.Fatalf("CrossCheck: %v (CLI hashes parsed: %+v)", err, result.CLI)
	}
	if result.Local.Tarball != tarballHash {
		t.Errorf("local tarball hash = %q, want %q", result.Local.Tarball, tarballHash)
	}
	if result.CLI.TarballBytes != len(tarball) {
		t.Errorf("CLI tarball bytes = %d, want %d (999 would mean the stderr line was parsed as data)", result.CLI.TarballBytes, len(tarball))
	}
	if result.CLI.Tarball != cliPrefix(t, tarballHash) {
		t.Errorf("CLI tarball hash = %q, want %q", result.CLI.Tarball, cliPrefix(t, tarballHash))
	}
	if result.CLI.Content != cliPrefix(t, content) {
		t.Errorf("CLI content hash = %q, want %q", result.CLI.Content, cliPrefix(t, content))
	}
	if result.CLI.Interface != cliPrefix(t, ifaceHash) {
		t.Errorf("CLI interface hash = %q, want %q", result.CLI.Interface, cliPrefix(t, ifaceHash))
	}
}

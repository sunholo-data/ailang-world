package pkgproj

// Row 139 D6: the leg-3 dry-run parsers (parseDryRun here, and the python
// observable check inside scripts/verify_world_package.sh step 7) accept every
// qualified release's `publish --dry-run` shape. Fixtures are the REAL stdout
// of `ailang publish --dry-run` over packages/world-core, captured 2026-10-16
// with the registry variables unset:
//   testdata/dryrun_v0.41.0.txt   - the pinned v0.41.0 binary
//   testdata/dryrun_v0.54.11.txt  - v0.54.11 (row 114's rel-v0.54.11 binary)
// v0.54.11 split the Tarball line ("Tarball: N bytes" + "Tarball hash:") and
// prints FULL 64-hex hashes where v0.41.0 printed 17-hex prefixes; the old
// parser failed it with "missing dry-run Tarball line".

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func readFixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// golden is world/core's committed ready packet: the identity both shapes must
// report.
func goldenPacket(t *testing.T) Hashes {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "scripts", "world_package_ready_packet.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var g struct {
		ContentHash   string `json:"contentHash"`
		InterfaceHash string `json:"interfaceHash"`
		TarballSHA256 string `json:"tarballSHA256"`
		TarballBytes  int    `json:"tarballBytes"`
	}
	if err := json.Unmarshal(b, &g); err != nil {
		t.Fatal(err)
	}
	return Hashes{Content: g.ContentHash, Interface: g.InterfaceHash, Tarball: g.TarballSHA256, TarballBytes: g.TarballBytes}
}

func TestParseDryRunAcceptsEveryQualifiedShape(t *testing.T) {
	golden := goldenPacket(t)
	old, err := parseDryRun([]byte(readFixture(t, "dryrun_v0.41.0.txt")))
	if err != nil {
		t.Fatalf("v0.41.0 shape: %v", err)
	}
	cur, err := parseDryRun([]byte(readFixture(t, "dryrun_v0.54.11.txt")))
	if err != nil {
		t.Fatalf("v0.54.11 shape: %v", err)
	}
	// The v0.41.0 shape abbreviates; the v0.54.11 shape is the full hash.
	for name, h := range map[string]string{"content": old.Content, "interface": old.Interface, "tarball": old.Tarball} {
		if len(h) != len("sha256:")+17 {
			t.Errorf("v0.41.0 %s %q is not a 17-hex prefix", name, h)
		}
	}
	for name, h := range map[string]string{"content": cur.Content, "interface": cur.Interface, "tarball": cur.Tarball} {
		if len(h) != len("sha256:")+64 {
			t.Errorf("v0.54.11 %s %q is not a full hash", name, h)
		}
	}
	// Both shapes report the committed identity, so Compare (the production
	// cross-check) passes both against the golden.
	if err := Compare(golden, old); err != nil {
		t.Errorf("v0.41.0 vs golden: %v", err)
	}
	if err := Compare(golden, cur); err != nil {
		t.Errorf("v0.54.11 vs golden: %v", err)
	}
	if cur != golden {
		t.Errorf("v0.54.11 full hashes %+v != golden %+v", cur, golden)
	}
	// Control: a full hash is checked EXACTLY, not by a laxer prefix — one
	// flipped trailing digit reds.
	bad := golden
	bad.Content = golden.Content[:len(golden.Content)-1] + flip(golden.Content[len(golden.Content)-1])
	if Compare(bad, cur) == nil {
		t.Fatal("control: a one-digit content change passed Compare against a full hash")
	}
}

func flip(c byte) string {
	if c == '0' {
		return "1"
	}
	return "0"
}

func TestParseDryRunLeadingWhitespaceAndRefusals(t *testing.T) {
	src := readFixture(t, "dryrun_v0.54.11.txt")
	// Any leading whitespace (D6): tab-indent every identity line.
	tabbed := regexp.MustCompile(`(?m)^  (Tarball|Content hash|Interface hash)`).ReplaceAllString(src, "\t\t$1")
	if tabbed == src {
		t.Fatal("control: tab rewrite changed nothing")
	}
	got, err := parseDryRun([]byte(tabbed))
	if err != nil || got != goldenPacket(t) {
		t.Fatalf("tab-indented v0.54.11 shape: %+v, %v", got, err)
	}
	cases := map[string]string{
		"no tarball hash":     regexp.MustCompile(`(?m)^.*Tarball hash:.*\n`).ReplaceAllString(src, ""),
		"no content hash":     regexp.MustCompile(`(?m)^.*Content hash:.*\n`).ReplaceAllString(src, ""),
		"two tarball hashes":  src + "  Tarball hash: sha256:" + strings.Repeat("a", 64) + "\n",
		"inline and separate": strings.Replace(readFixture(t, "dryrun_v0.41.0.txt"), "  Content hash:", "  Tarball hash: sha256:"+strings.Repeat("b", 64)+"\n  Content hash:", 1),
		"short hash no dots":  strings.Replace(src, "Content hash: sha256:473517079249f3959b5aba9727f62b40130624c9915612a47736fdc1dae2a780", "Content hash: sha256:473517079249f3959", 1),
	}
	for name, out := range cases {
		if out == src {
			t.Fatalf("control: case %q did not mutate the fixture", name)
		}
		if _, err := parseDryRun([]byte(out)); err == nil {
			t.Errorf("%s: parsed, want refusal", name)
		}
	}
	// "Interface hash (v2):" must never be read as the v1 interface hash.
	v2only := regexp.MustCompile(`(?m)^  Interface hash: .*\n`).ReplaceAllString(src, "")
	if _, err := parseDryRun([]byte(v2only)); err == nil || !strings.Contains(err.Error(), "Interface hash") {
		t.Errorf("v2 line alone satisfied the interface hash: %v", err)
	}
}

// The python observable check of verify_world_package.sh step 7, extracted
// verbatim from the script so the fixture runs the shipped code.
func TestVerifyWorldPackageDryRunCheckAcceptsEveryQualifiedShape(t *testing.T) {
	py, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal("python3 is required by scripts/verify_world_package.sh and so by this test")
	}
	script, err := os.ReadFile(filepath.Join("..", "..", "scripts", "verify_world_package.sh"))
	if err != nil {
		t.Fatal(err)
	}
	const open = "python3 - \"$tmp_dry\" <<'PY' || exit 1\n"
	s := string(script)
	i := strings.Index(s, open)
	if i < 0 {
		t.Fatal("step-7 python heredoc not found in verify_world_package.sh")
	}
	body := s[i+len(open):]
	j := strings.Index(body, "\nPY\n")
	if j < 0 {
		t.Fatal("step-7 python heredoc unterminated")
	}
	prog := filepath.Join(t.TempDir(), "step7.py")
	if err := os.WriteFile(prog, []byte(body[:j+1]), 0o644); err != nil {
		t.Fatal(err)
	}
	run := func(content string) (string, error) {
		f := filepath.Join(t.TempDir(), "dry.txt")
		if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		out, err := exec.Command(py, "-I", prog, f).CombinedOutput()
		return string(out), err
	}
	for _, name := range []string{"dryrun_v0.41.0.txt", "dryrun_v0.54.11.txt"} {
		if out, err := run(readFixture(t, name)); err != nil {
			t.Errorf("%s: step-7 check refused: %v %s", name, err, out)
		}
	}
	tabbed := regexp.MustCompile(`(?m)^  `).ReplaceAllString(readFixture(t, "dryrun_v0.54.11.txt"), "\t")
	if out, err := run(tabbed); err != nil {
		t.Errorf("tab-indented: step-7 check refused: %v %s", err, out)
	}
	// Control: the check is not vacuous — drop the Content hash line and it reds.
	noContent := regexp.MustCompile(`(?m)^.*Content hash:.*\n`).ReplaceAllString(readFixture(t, "dryrun_v0.54.11.txt"), "")
	if out, err := run(noContent); err == nil || !strings.Contains(out, "Content hash") {
		t.Errorf("control: missing Content hash passed the step-7 check: %v %s", err, out)
	}
}

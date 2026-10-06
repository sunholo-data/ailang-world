package broker

// Row 140 M2 (w-workspace-exec-toolchain-effect §4.6 "Integrity of the
// confinement stack"): the whole srt install tree archived under the state
// dir, and the per-call re-verification of the tree, node and libnode.

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeSrtTree is a node_modules dir shaped like srt's install: the package
// with its cli.js and seccomp helper, a runtime dependency and a .bin link.
func fakeSrtTree(t *testing.T, version, cli string) string {
	t.Helper()
	nm := filepath.Join(t.TempDir(), "node_modules")
	pkg := filepath.Join(nm, "@anthropic-ai", "sandbox-runtime")
	files := map[string]string{
		filepath.Join(pkg, "package.json"):                                `{"name":"@anthropic-ai/sandbox-runtime","version":"` + version + `"}`,
		filepath.Join(pkg, "dist", "cli.js"):                              cli,
		filepath.Join(pkg, "vendor", "seccomp", "x64", "apply-seccomp"):   "#!/bin/sh\n",
		filepath.Join(pkg, "vendor", "seccomp", "arm64", "apply-seccomp"): "#!/bin/sh\n",
		filepath.Join(nm, "commander", "index.js"):                        "module.exports = {}\n",
		filepath.Join(nm, "commander", "package.json"):                    `{"name":"commander"}`,
	}
	for p, body := range files {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		mode := os.FileMode(0o644)
		if strings.HasSuffix(p, "apply-seccomp") {
			mode = 0o755
		}
		if err := os.WriteFile(p, []byte(body), mode); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(nm, ".bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../@anthropic-ai/sandbox-runtime/dist/cli.js", filepath.Join(nm, ".bin", "srt")); err != nil {
		t.Fatal(err)
	}
	return nm
}

// fakeNode is an executable standing in for node, with a libnode beside it
// in node's install layout (<prefix>/bin/node, <prefix>/lib/libnode.*).
func fakeNode(t *testing.T, body string) (node, libnode string) {
	t.Helper()
	prefix := t.TempDir()
	node = filepath.Join(prefix, "bin", "node")
	libnode = filepath.Join(prefix, "lib", "libnode.147.dylib")
	for _, p := range []string{node, libnode} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(node, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(libnode, []byte("not really a dylib"), 0o644); err != nil {
		t.Fatal(err)
	}
	real, err := filepath.EvalSymlinks(node)
	if err != nil {
		t.Fatal(err)
	}
	reallib, _ := filepath.EvalSymlinks(libnode)
	return real, reallib
}

func fakePin(t *testing.T, nm string) ExecSandboxPin {
	t.Helper()
	sum, err := sha256File(filepath.Join(nm, "@anthropic-ai", "sandbox-runtime", "dist", "cli.js"))
	if err != nil {
		t.Fatal(err)
	}
	return ExecSandboxPin{Version: "0.0.0-test", CLISHA256: sum}
}

func TestArchiveExecSandboxCopiesTheWholeTreeAndVerifies(t *testing.T) {
	nm := fakeSrtTree(t, "0.0.0-test", "console.log('cli')\n")
	node, libnode := fakeNode(t, "#!/bin/sh\n")
	state := canonicalTempDir(t)
	sb, err := ArchiveExecSandbox(ExecSandboxConfig{NodeModules: nm, Node: node, StateDir: state, Pin: fakePin(t, nm)})
	if err != nil {
		t.Fatalf("ArchiveExecSandbox: %v", err)
	}
	wantDir := filepath.Join(state, "archive", "exec-sandbox", sb.Digest, "node_modules")
	if sb.Dir != wantDir {
		t.Fatalf("archive dir = %q, want %q", sb.Dir, wantDir)
	}
	for _, rel := range []string{"@anthropic-ai/sandbox-runtime/dist/cli.js", "commander/index.js", ".bin/srt",
		"@anthropic-ai/sandbox-runtime/vendor/seccomp/x64/apply-seccomp"} {
		if _, err := os.Lstat(filepath.Join(sb.Dir, rel)); err != nil {
			t.Fatalf("archive lacks %s: the WHOLE install tree must be archived (V42)", rel)
		}
	}
	if info, _ := os.Stat(filepath.Join(sb.Dir, "@anthropic-ai/sandbox-runtime/vendor/seccomp/x64/apply-seccomp")); info.Mode().Perm()&0o111 == 0 {
		t.Fatal("the seccomp helper lost its exec bit in the archive")
	}
	if sb.CLI() != filepath.Join(sb.Dir, "@anthropic-ai", "sandbox-runtime", "dist", "cli.js") {
		t.Fatalf("CLI() = %q", sb.CLI())
	}
	if sb.Node != node || sb.NodeSHA256 == "" || len(sb.Libnode) != 1 || sb.Libnode[0].Path != libnode {
		t.Fatalf("node %q %q libnode %+v", sb.Node, sb.NodeSHA256, sb.Libnode)
	}
	if err := sb.Verify(); err != nil {
		t.Fatalf("Verify on a fresh archive: %v", err)
	}
	again, err := ArchiveExecSandbox(ExecSandboxConfig{NodeModules: nm, Node: node, StateDir: state, Pin: fakePin(t, nm)})
	if err != nil || again.Digest != sb.Digest {
		t.Fatalf("re-archiving the same tree: %v, digest %q vs %q", err, again.Digest, sb.Digest)
	}
	if sb.SeccompDir("linux", "amd64") != filepath.Join(sb.Dir, "@anthropic-ai/sandbox-runtime/vendor/seccomp/x64") ||
		sb.SeccompDir("linux", "arm64") != filepath.Join(sb.Dir, "@anthropic-ai/sandbox-runtime/vendor/seccomp/arm64") ||
		sb.SeccompDir("darwin", "arm64") != "" {
		t.Fatal("SeccompDir maps the platform wrongly")
	}
}

func TestArchiveExecSandboxRefusesAnUnpinnedTree(t *testing.T) {
	node, _ := fakeNode(t, "#!/bin/sh\n")
	nm := fakeSrtTree(t, ExecSandboxRelease, "not the pinned cli\n")
	// The zero pin is the production pin: version 0.0.78 and its cli.js digest.
	if _, err := ArchiveExecSandbox(ExecSandboxConfig{NodeModules: nm, Node: node, StateDir: canonicalTempDir(t)}); err == nil ||
		!strings.Contains(err.Error(), "cli.js") {
		t.Fatalf("a tree whose cli.js is not the pinned one was archived: %v", err)
	}
	nm = fakeSrtTree(t, "0.0.71", "x\n")
	pin := fakePin(t, nm)
	pin.Version = ExecSandboxRelease
	if _, err := ArchiveExecSandbox(ExecSandboxConfig{NodeModules: nm, Node: node, StateDir: canonicalTempDir(t), Pin: pin}); err == nil ||
		!strings.Contains(err.Error(), "0.0.71") {
		t.Fatalf("a 0.0.71 tree was archived: %v", err)
	}
	if _, err := ArchiveExecSandbox(ExecSandboxConfig{NodeModules: nm, Node: "node", StateDir: canonicalTempDir(t), Pin: fakePin(t, nm)}); err == nil {
		t.Fatal("a relative node path was accepted")
	}
}

// The production pin is M0's: pin.json is the one source.
func TestExecSandboxPinIsM0s(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "design_docs", "verification", "world-row140-m0", "pin.json"))
	if err != nil {
		t.Fatal(err)
	}
	var pin struct {
		Package, Version string
		CLISHA256        string `json:"cli_sha256"`
	}
	if err := json.Unmarshal(data, &pin); err != nil {
		t.Fatal(err)
	}
	if pin.Package != ExecSandboxPackage || pin.Version != ExecSandboxRelease || pin.CLISHA256 != ExecSandboxCLISHA256 {
		t.Fatalf("pin.json %+v != broker pin %s %s %s", pin, ExecSandboxPackage, ExecSandboxRelease, ExecSandboxCLISHA256)
	}
}

// AC2.9's unit half: every kind of drift is named by file.
func TestExecSandboxVerifyNamesTheDriftedFile(t *testing.T) {
	setup := func(t *testing.T) (*ExecSandbox, string) {
		nm := fakeSrtTree(t, "0.0.0-test", "console.log('cli')\n")
		node, _ := fakeNode(t, "#!/bin/sh\n")
		sb, err := ArchiveExecSandbox(ExecSandboxConfig{NodeModules: nm, Node: node, StateDir: canonicalTempDir(t), Pin: fakePin(t, nm)})
		if err != nil {
			t.Fatal(err)
		}
		return sb, nm
	}
	for _, tc := range []struct {
		name   string
		mutate func(t *testing.T, sb *ExecSandbox) string // returns the path the error must name
	}{
		{"flipped byte in cli.js", func(t *testing.T, sb *ExecSandbox) string {
			flipByte(t, sb.CLI())
			return sb.CLI()
		}},
		{"removed dependency file", func(t *testing.T, sb *ExecSandbox) string {
			p := filepath.Join(sb.Dir, "commander", "index.js")
			if err := os.Remove(p); err != nil {
				t.Fatal(err)
			}
			return p
		}},
		{"added file", func(t *testing.T, sb *ExecSandbox) string {
			p := filepath.Join(sb.Dir, "commander", "evil.js")
			if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
				t.Fatal(err)
			}
			return p
		}},
		{"retargeted link", func(t *testing.T, sb *ExecSandbox) string {
			p := filepath.Join(sb.Dir, ".bin", "srt")
			_ = os.Remove(p)
			if err := os.Symlink("/bin/sh", p); err != nil {
				t.Fatal(err)
			}
			return p
		}},
		{"node swapped", func(t *testing.T, sb *ExecSandbox) string {
			flipByte(t, sb.Node)
			return sb.Node
		}},
		{"libnode swapped", func(t *testing.T, sb *ExecSandbox) string {
			flipByte(t, sb.Libnode[0].Path)
			return sb.Libnode[0].Path
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sb, _ := setup(t)
			want := tc.mutate(t, sb)
			err := sb.Verify()
			var drift *ExecStackDriftError
			if !errors.As(err, &drift) || drift.Path != want || !strings.Contains(err.Error(), want) {
				t.Fatalf("Verify = %v, want an *ExecStackDriftError naming %s", err, want)
			}
		})
	}
}

// flipByte inverts the first byte of path, making it writable first (the
// archive is read-only).
func flipByte(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, info.Mode().Perm()|0o200); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data[0] ^= 0xff
	if err := os.WriteFile(path, data, info.Mode().Perm()|0o200); err != nil {
		t.Fatal(err)
	}
}

package pkgproj

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
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

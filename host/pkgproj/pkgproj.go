// Package pkgproj re-implements the AILANG v0.30.0 package projection hashes.
//
// This is host-side publication verification code, not a published package:
// it computes the hashes which authorize publication, so packaging it would be
// circular.
package pkgproj

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Manifest struct {
	Package Package
	Exports Exports
	Effects Effects
}

type Package struct{ Name, Edition, AILANG string }
type Exports struct{ Modules []string }
type Effects struct{ Max []string }

func ContentHash(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	var paths []string
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".ail") {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		paths = append(paths, rel)
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(paths)
	h := sha256.New()
	for _, rel := range paths {
		if _, err := fmt.Fprintf(h, "file:%s\n", rel); err != nil {
			return "", err
		}
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			return "", err
		}
		if _, err := h.Write(data); err != nil {
			return "", err
		}
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}

func InterfaceHash(manifest Manifest) string {
	h := sha256.New()
	fmt.Fprintf(h, "name:%s\n", manifest.Package.Name)
	fmt.Fprintf(h, "edition:%s\n", manifest.Package.Edition)
	if manifest.Package.AILANG != "" {
		fmt.Fprintf(h, "ailang:%s\n", manifest.Package.AILANG)
	}
	exports := append([]string(nil), manifest.Exports.Modules...)
	sort.Strings(exports)
	for _, mod := range exports {
		fmt.Fprintf(h, "export:%s\n", mod)
	}
	effects := append([]string(nil), manifest.Effects.Max...)
	sort.Strings(effects)
	for _, effect := range effects {
		fmt.Fprintf(h, "effect:%s\n", effect)
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}

func CreateTarball(packageDir string) ([]byte, error) {
	packageDir, err := filepath.Abs(packageDir)
	if err != nil {
		return nil, err
	}
	files := make(map[string]string)
	err = filepath.Walk(packageDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			switch info.Name() {
			case ".git", "tests", "test":
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(packageDir, path)
		if err != nil {
			return err
		}
		clean := filepath.Clean(rel)
		if strings.HasPrefix(clean, "..") || strings.Contains(clean, "../") {
			return fmt.Errorf("unsafe package path %q", rel)
		}
		forward := filepath.ToSlash(rel)
		if forward == "ailang.toml" || strings.HasSuffix(forward, ".ail") || forward == "AGENT.md" || strings.HasPrefix(forward, "assets/") {
			files[forward] = path
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	for _, name := range names {
		data, err := os.ReadFile(files[name])
		if err != nil {
			return nil, err
		}
		header := &tar.Header{Name: name, Size: int64(len(data)), Mode: 0644, ModTime: time.Time{}}
		if err := tw.WriteHeader(header); err != nil {
			return nil, err
		}
		if _, err := tw.Write(data); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TarballHash(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

type Hashes struct {
	Content, Interface, Tarball string
	TarballBytes                int
}
type CrossCheckResult struct{ Local, CLI Hashes }

var dryRunLine = regexp.MustCompile(`(?m)^  (Tarball|Content hash|Interface hash): (?:([0-9]+) bytes \()?((?:sha256:)[0-9a-f]{17})\.\.\.?\)?$`)

// Compare requires each CLI prefix to match its corresponding full local hash.
func Compare(local, cli Hashes) error {
	var mismatches []string
	if local.TarballBytes != cli.TarballBytes {
		mismatches = append(mismatches, fmt.Sprintf("tarball byte length mismatch: local=%d cli=%d", local.TarballBytes, cli.TarballBytes))
	}
	for _, arm := range []struct{ name, full, prefix string }{
		{"content", local.Content, cli.Content}, {"interface", local.Interface, cli.Interface}, {"tarball", local.Tarball, cli.Tarball},
	} {
		if !strings.HasPrefix(arm.full, arm.prefix) {
			mismatches = append(mismatches, fmt.Sprintf("%s hash mismatch: local=%s cli=%s", arm.name, arm.full, arm.prefix))
		}
	}
	if len(mismatches) != 0 {
		return fmt.Errorf("%s", strings.Join(mismatches, "; "))
	}
	return nil
}

// CrossCheck computes all hashes and compares them with one pinned CLI dry-run.
func CrossCheck(packageDir string, manifest Manifest, ailangBin string) (CrossCheckResult, error) {
	content, err := ContentHash(packageDir)
	if err != nil {
		return CrossCheckResult{}, err
	}
	tarball, err := CreateTarball(packageDir)
	if err != nil {
		return CrossCheckResult{}, err
	}
	local := Hashes{Content: content, Interface: InterfaceHash(manifest), Tarball: TarballHash(tarball), TarballBytes: len(tarball)}
	cmd := exec.Command(ailangBin, "publish", "--dry-run")
	cmd.Dir = packageDir
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "AILANG_REGISTRY_API_KEY=") {
			cmd.Env = append(cmd.Env, env)
		}
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return CrossCheckResult{Local: local}, fmt.Errorf("ailang publish --dry-run: %w: %s", err, out)
	}
	cli, err := parseDryRun(out)
	result := CrossCheckResult{Local: local, CLI: cli}
	if err != nil {
		return result, err
	}
	if err := Compare(local, cli); err != nil {
		return result, err
	}
	return result, nil
}

func parseDryRun(out []byte) (Hashes, error) {
	var got Hashes
	seen := map[string]bool{}
	for _, m := range dryRunLine.FindAllSubmatch(out, -1) {
		label, size, hash := string(m[1]), string(m[2]), string(m[3])
		if seen[label] {
			return Hashes{}, fmt.Errorf("duplicate dry-run %s line", label)
		}
		seen[label] = true
		switch label {
		case "Tarball":
			n, err := strconv.Atoi(size)
			if err != nil {
				return Hashes{}, err
			}
			got.TarballBytes, got.Tarball = n, hash
		case "Content hash":
			got.Content = hash
		case "Interface hash":
			got.Interface = hash
		}
	}
	for _, label := range []string{"Tarball", "Content hash", "Interface hash"} {
		if !seen[label] {
			return Hashes{}, fmt.Errorf("missing dry-run %s line in output: %s", label, out)
		}
	}
	return got, nil
}

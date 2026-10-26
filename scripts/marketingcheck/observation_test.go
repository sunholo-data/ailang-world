package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

// observationPython is the plan's C11 boundary command, copied verbatim from
// w-marketing-capture-sprint-plan.md "C11 observations and mutation discipline".
const observationPython = `import pathlib, hashlib, json
p = pathlib.Path('.')
paths = sorted(set([q for pat in ('world/*.ail','go.mod','go.sum','design_docs/world-mission.md') for q in p.glob(pat) if q.is_file()] + [q for q in p.joinpath('tools/launchd').rglob('*') if q.is_file()]))
print(json.dumps([{'path':str(q),'sha256':hashlib.sha256(q.read_bytes()).hexdigest(),'lines':q.read_bytes().count(b'\n')} for q in paths],sort_keys=True))
`

var protectedFixture = map[string]string{
	"world/a.ail":                      "alpha one\nalpha two\n",
	"world/b.ail":                      "beta one\nbeta two\nbeta three\n",
	"world/sub/c.ail":                  "nested one\nnested two\n",
	"go.mod":                           "module fixture\n\ngo 1.26\n",
	"go.sum":                           "sum one\nsum two\n",
	"design_docs/world-mission.md":     "# charter\nrow one\n",
	"tools/launchd/x.sh":               "#!/bin/sh\necho x\n",
	"tools/launchd/testdata/deep/y.md": "deep one\ndeep two\n",
	"README.md":                        "readme one\nreadme two\n",
}

// protectedPaths is the sorted C11 subset of protectedFixture.
var protectedPaths = []string{"design_docs/world-mission.md", "go.mod", "go.sum", "tools/launchd/testdata/deep/y.md", "tools/launchd/x.sh", "world/a.ail", "world/b.ail"}

func protectedTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for rel, body := range protectedFixture {
		p := filepath.Join(root, filepath.FromSlash(rel))
		pngOK(t, os.MkdirAll(filepath.Dir(p), 0o755))
		pngOK(t, os.WriteFile(p, []byte(body), 0o644))
	}
	return root
}

// flipProtectedByte turns "alpha one" into "blpha one": same length, same newline count.
func flipProtectedByte(t *testing.T, root string) {
	t.Helper()
	p := filepath.Join(root, "world", "a.ail")
	raw, err := os.ReadFile(p)
	pngOK(t, err)
	raw[0] = 'b'
	pngOK(t, os.WriteFile(p, raw, 0o644))
}

func mustObserve(t *testing.T, root string) []ObservationRecord {
	t.Helper()
	records, err := Observe(root)
	pngOK(t, err)
	return records
}

// pythonObserve runs observationPython in root and returns the path of its stdout file.
func pythonObserve(t *testing.T, root string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dir := t.TempDir()
	outPath, errPath := filepath.Join(dir, "python-observation.json"), filepath.Join(dir, "python-stderr.txt")
	out, err := os.Create(outPath)
	pngOK(t, err)
	defer out.Close()
	errFile, err := os.Create(errPath)
	pngOK(t, err)
	defer errFile.Close()
	cmd := exec.CommandContext(ctx, "python3", "-I", "-c", observationPython)
	cmd.Dir = root
	cmd.Stdout = out
	cmd.Stderr = errFile
	if err := cmd.Run(); err != nil {
		msg, _ := os.ReadFile(errPath)
		t.Fatalf("python observation: %v %s", err, msg)
	}
	return outPath
}

// runObservationCLI runs bin with file sinks and returns the exit code and stderr.
func runObservationCLI(t *testing.T, bin string, args ...string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	dir := t.TempDir()
	stdout, err := os.Create(filepath.Join(dir, "stdout"))
	pngOK(t, err)
	defer stdout.Close()
	stderrPath := filepath.Join(dir, "stderr")
	stderr, err := os.Create(stderrPath)
	pngOK(t, err)
	defer stderr.Close()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	rc := 0
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatal(err)
		}
		rc = ee.ExitCode()
	}
	text, err := os.ReadFile(stderrPath)
	pngOK(t, err)
	return rc, string(text)
}

func TestMarketingProtectedObservationComparator(t *testing.T) {
	t.Run("set-coverage", func(t *testing.T) {
		got := mustObserve(t, protectedTree(t))
		var want []ObservationRecord
		for _, rel := range protectedPaths {
			sum := sha256.Sum256([]byte(protectedFixture[rel]))
			want = append(want, ObservationRecord{Path: rel, SHA256: hex.EncodeToString(sum[:]), Lines: strings.Count(protectedFixture[rel], "\n")})
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Observe must cover exactly the C11 protected set: got %+v want %+v", got, want)
		}
	})
	t.Run("equal", func(t *testing.T) {
		root := protectedTree(t)
		if err := CompareObservations(mustObserve(t, root), mustObserve(t, root)); err != nil {
			t.Fatalf("unchanged copy must compare equal: %v", err)
		}
	})
	t.Run("same-line-count-change", func(t *testing.T) {
		root := protectedTree(t)
		before := mustObserve(t, root)
		flipProtectedByte(t, root)
		after := mustObserve(t, root)
		if before[5].Path != "world/a.ail" || after[5].Path != "world/a.ail" || before[5].Lines != after[5].Lines {
			t.Fatalf("fixture must be a genuine same-line-count change: %+v %+v", before[5], after[5])
		}
		err := CompareObservations(before, after)
		wantCode(t, err, "observation-digest", "same-line-count changed copy must red")
		var r *Refusal
		if !errors.As(err, &r) || r.Name != "world/a.ail" {
			t.Fatalf("same-line-count change must name world/a.ail: %v", err)
		}
	})
	t.Run("digest-equal-line-change", func(t *testing.T) {
		before := mustObserve(t, protectedTree(t))
		after := slices.Clone(before)
		after[0].Lines++
		wantCode(t, CompareObservations(before, after), "observation-lines", "line-count mismatch must refuse")
	})
	t.Run("deleted", func(t *testing.T) {
		root := protectedTree(t)
		before := mustObserve(t, root)
		pngOK(t, os.Remove(filepath.Join(root, "tools", "launchd", "testdata", "deep", "y.md")))
		wantCode(t, CompareObservations(before, mustObserve(t, root)), "observation-set", "deleted protected copy must red")
	})
	t.Run("added", func(t *testing.T) {
		root := protectedTree(t)
		before := mustObserve(t, root)
		pngOK(t, os.WriteFile(filepath.Join(root, "world", "z.ail"), []byte("zeta\n"), 0o644))
		wantCode(t, CompareObservations(before, mustObserve(t, root)), "observation-set", "added protected copy must red")
	})
	t.Run("duplicate", func(t *testing.T) {
		root := protectedTree(t)
		before := mustObserve(t, root)
		dup := append(slices.Clone(before), before[0])
		wantCode(t, CompareObservations(dup, mustObserve(t, root)), "observation-set", "duplicate path must red")
	})
	t.Run("empty", func(t *testing.T) {
		wantCode(t, CompareObservations(nil, nil), "observation-empty", "empty observation must not compare equal")
	})
	t.Run("python-parity", func(t *testing.T) {
		root := protectedTree(t)
		raw, err := os.ReadFile(pythonObserve(t, root))
		pngOK(t, err)
		var py []ObservationRecord
		pngOK(t, json.Unmarshal(raw, &py))
		slices.SortFunc(py, func(a, b ObservationRecord) int { return strings.Compare(a.Path, b.Path) })
		if got := mustObserve(t, root); !reflect.DeepEqual(py, got) {
			t.Fatalf("Go observer must equal the plan's Python boundary command: python %+v go %+v", py, got)
		}
	})
	t.Run("cli", func(t *testing.T) {
		dir := t.TempDir()
		bin := filepath.Join(dir, "marketingcheck")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if out, err := exec.CommandContext(ctx, "go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
			t.Fatalf("CLI build: %v %s", err, out)
		}
		root := protectedTree(t)
		f1, f2 := filepath.Join(dir, "f1.json"), filepath.Join(dir, "f2.json")
		expect := func(label string, rc int, stderr string, wantRC int, wantText string) {
			t.Helper()
			if rc != wantRC || !strings.Contains(stderr, wantText) {
				t.Fatalf("%s: want rc %d stderr containing %q, got rc %d stderr %q", label, wantRC, wantText, rc, stderr)
			}
		}
		rc, se := runObservationCLI(t, bin, "observe", "--root", root, "--out", f1)
		expect("(i) observe --out", rc, se, 0, "")
		first, err := os.ReadFile(f1)
		pngOK(t, err)
		var parsed []ObservationRecord
		pngOK(t, json.Unmarshal(first, &parsed))
		if len(parsed) != len(protectedPaths) {
			t.Fatalf("(i) observe output must hold %d records, got %d", len(protectedPaths), len(parsed))
		}
		rc, se = runObservationCLI(t, bin, "observe", "--root", root, "--out", f1)
		expect("(ii) existing --out", rc, se, 1, "marketingcheck:")
		again, err := os.ReadFile(f1)
		pngOK(t, err)
		if string(again) != string(first) {
			t.Fatal("(ii) existing --out must remain unchanged")
		}
		rc, se = runObservationCLI(t, bin, "compare", "--before", f1, "--after", f1)
		expect("(iv) compare equal", rc, se, 0, "")
		rc, se = runObservationCLI(t, bin, "compare", "--before", pythonObserve(t, root), "--after", f1)
		expect("(vi) python output vs observe output", rc, se, 0, "")
		rc, se = runObservationCLI(t, bin, "compare", "--before", f1)
		expect("(v) compare missing --after", rc, se, 2, "")
		bad := filepath.Join(dir, "bad.json")
		pngOK(t, os.WriteFile(bad, []byte("{"), 0o600))
		rc, se = runObservationCLI(t, bin, "compare", "--before", f1, "--after", bad)
		expect("(vii) unparsable JSON", rc, se, 1, "marketingcheck: observation-json")
		rc, se = runObservationCLI(t, bin, "observe")
		expect("(viii) observe missing --root", rc, se, 2, "")
		flipProtectedByte(t, root)
		rc, se = runObservationCLI(t, bin, "observe", "--root", root, "--out", f2)
		expect("(iii) observe after flip", rc, se, 0, "")
		rc, se = runObservationCLI(t, bin, "compare", "--before", f1, "--after", f2)
		expect("(iii) compare after flip", rc, se, 1, "marketingcheck: observation-digest world/a.ail:")
	})
	t.Run("repo-read-only", func(t *testing.T) {
		first, second := mustObserve(t, "../.."), mustObserve(t, "../..")
		if err := CompareObservations(first, second); err != nil {
			t.Fatalf("real protected set observable without writes: %v", err)
		}
		paths := observationPaths(first)
		hasLaunchd := slices.ContainsFunc(paths, func(p string) bool { return strings.HasPrefix(p, "tools/launchd/") })
		if len(first) < 10 || !slices.Contains(paths, "go.mod") || !slices.Contains(paths, "design_docs/world-mission.md") || !hasLaunchd {
			t.Fatalf("real protected set observable without writes: %d records %v", len(first), paths)
		}
	})
}

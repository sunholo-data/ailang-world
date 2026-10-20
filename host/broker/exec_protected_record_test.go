package broker

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

var recordPatterns = []string{"*.test", "*_test.go"}

const recordTreeDigest = "sha256:c81ac713cec2a6c252a5e8f79c8131e00fb251549d9e51a2abab9ad283d40c58"
const recordIgnoredDigest = "sha256:df2826e62a4850b13e1f6ec3d690171f5e45423b78526862bcd3abb412bd4529"
const recordEmptyDigest = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// Write loose fixture objects directly: no git init/add/commit or other Git
// write operations. The resulting checkout is read by the real Git binary.
func recordRepo(t *testing.T, root string) {
	t.Helper()
	gitdir := filepath.Join(root, ".git")
	for _, d := range []string{"objects", "refs/heads"} {
		if err := os.MkdirAll(filepath.Join(gitdir, d), 0700); err != nil {
			t.Fatal(err)
		}
	}
	object := func(kind string, data []byte) string {
		raw := append([]byte(fmt.Sprintf("%s %d\x00", kind, len(data))), data...)
		oid := fmt.Sprintf("%x", sha1.Sum(raw))
		var compressed bytes.Buffer
		w := zlib.NewWriter(&compressed)
		if _, err := w.Write(raw); err != nil {
			t.Fatal(err)
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		dir := filepath.Join(gitdir, "objects", oid[:2])
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, oid[2:]), compressed.Bytes(), 0600); err != nil {
			t.Fatal(err)
		}
		return oid
	}
	var tree []byte
	for _, f := range []struct{ mode, path, data string }{
		{"100644", ".gitignore", "*.test\n"}, {"100644", "a_test.go", "alpha\n"},
		{"100755", "b_test.go", "beta\n"}, {"100644", "ordinary.txt", "noise\n"},
	} {
		oid, _ := hex.DecodeString(object("blob", []byte(f.data)))
		tree = append(tree, []byte(f.mode+" "+f.path+"\x00")...)
		tree = append(tree, oid...)
		mustWrite(t, filepath.Join(root, f.path), f.data)
	}
	// A nested protected file, so the capture's ls-tree must recurse (-r).
	blob, _ := hex.DecodeString(object("blob", []byte("gamma\n")))
	sub, _ := hex.DecodeString(object("tree", append([]byte("100644 c_test.go\x00"), blob...)))
	tree = append(append(tree, []byte("40000 pkg\x00")...), sub...)
	if err := os.MkdirAll(filepath.Join(root, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(root, "pkg", "c_test.go"), "gamma\n")
	commit := object("commit", []byte("tree "+object("tree", tree)+"\nauthor Fixture <fixture@example.invalid> 1 +0000\ncommitter Fixture <fixture@example.invalid> 1 +0000\n\nfixture\n"))
	mustWrite(t, filepath.Join(gitdir, "refs/heads/main"), commit+"\n")
	mustWrite(t, filepath.Join(gitdir, "HEAD"), "ref: refs/heads/main\n")
	mustWrite(t, filepath.Join(gitdir, "config"), "[core]\nrepositoryformatversion = 0\nbare = false\nprecomposeunicode = false\n")
}

func recordEnv(t *testing.T) string {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	return git
}

func recordReadGit(t *testing.T, git, root string, args ...string) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, git, args...)
	cmd.Dir, cmd.Env = root, gitEnv()
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("fixture git %v: %v", args, err)
	}
	return out
}

func recordIndependent(t *testing.T, git, root string) (base, digest string, files int) {
	t.Helper()
	base = strings.TrimSpace(string(recordReadGit(t, git, root, "rev-parse", "HEAD")))
	h := sha256.New()
	out := recordReadGit(t, git, root, "ls-tree", "-r", "-z", "--full-tree", "HEAD")
	for _, line := range bytes.Split(bytes.TrimSuffix(out, []byte{0}), []byte{0}) {
		header, path, _ := strings.Cut(string(line), "\t")
		if matchProtected(recordPatterns, path) {
			fields := strings.Fields(header)
			fmt.Fprintf(h, "%s %s\t%s\n", fields[0], fields[2], path)
			files++
		}
	}
	return base, fmt.Sprintf("sha256:%x", h.Sum(nil)), files
}

func recordAssertTree(t *testing.T, b *protectedBaseline, git, root string) {
	t.Helper()
	base, digest, files := recordIndependent(t, git, root)
	if b == nil {
		t.Fatal("handler baseline missing: capture was not wired")
	}
	if b.err != nil || b.base != base || b.digest != digest || b.digest != recordTreeDigest || b.files != files || files != 3 {
		t.Fatalf("committed tree identity: baseline=%+v error=%v independent=%s/%s/%d fixed=%s", b, b.err, base, digest, files, recordTreeDigest)
	}
}

func recordHandler(t *testing.T, setup func(*execFixture)) *execFixture {
	t.Helper()
	return newExecFixture(t, execFixtureOpts{profile: func(m map[string]any) { m["protected"] = recordPatterns }, setup: setup})
}

func TestExecProtectedBaseline(t *testing.T) {
	oldBound := protectedBaselineTimeout
	protectedBaselineTimeout = 30 * time.Second // Like protectedRelaxBound: these rows test identity, not scheduling.
	t.Cleanup(func() { protectedBaselineTimeout = oldBound })
	for _, name := range []string{"independent-digest", "dirty-worktree", "empty-set", "snapshot-immutable", "through-handler", "ignored-snapshot", "no-protected-no-git"} {
		t.Run(name, func(t *testing.T) {
			git := recordEnv(t)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			if name == "no-protected-no-git" {
				fake, log := protectedFakeGit(t, "exit 99")
				old := worktreeGit
				worktreeGit = func() (string, error) { return fake, nil }
				t.Cleanup(func() { worktreeGit = old })
				for _, empty := range []bool{false, true} {
					f := newExecFixture(t, execFixtureOpts{profile: func(m map[string]any) {
						if empty {
							m["protected"] = []string{}
						}
					}})
					if f.h.baseline != nil || f.h.spawns.Load() != 0 {
						t.Fatal("unprotected handler stored a baseline or counted spawns")
					}
				}
				if _, err := os.Stat(log); !os.IsNotExist(err) {
					t.Fatalf("unprotected construction invoked Git: %v", err)
				}
				return
			}
			if name == "through-handler" {
				f := recordHandler(t, func(f *execFixture) { recordRepo(t, f.worktree) })
				recordAssertTree(t, f.h.baseline, git, f.worktree)
				if f.h.baseline.ignoredDigest != recordEmptyDigest || len(f.h.baseline.ignored) != 0 || f.h.spawns.Load() != 0 {
					t.Fatalf("handler ignored metadata/spawns: %+v/%d", f.h.baseline, f.h.spawns.Load())
				}
				return
			}
			root := t.TempDir()
			recordRepo(t, root)
			if name == "ignored-snapshot" || name == "snapshot-immutable" {
				mustWrite(t, filepath.Join(root, "Case.test"), "case\n")
				mustWrite(t, filepath.Join(root, "é.test"), "accent\n")
				// Ignored and selected by the widened pathspec (*.te*t) but not
				// by the mirror (*.test): only the mirror filter excludes it.
				if err := os.MkdirAll(filepath.Join(root, ".git", "info"), 0o700); err != nil {
					t.Fatal(err)
				}
				mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), "x.text\n")
				mustWrite(t, filepath.Join(root, "x.text"), "widened only\n")
				// APFS cannot hold both case/NFC aliases. Keep real Git's
				// ignored output and append equivalent spellings as on ext4.
				git, _ = protectedFakeGit(t, "case \" $* \" in *' ls-files '*) '"+git+"' \"$@\" || exit $?; printf 'case.test\\000é.test\\000é.test\\000';; *) exec '"+git+"' \"$@\";; esac")
			}
			b := captureProtectedBaseline(ctx, git, root, recordPatterns)
			if name != "dirty-worktree" {
				recordAssertTree(t, &b, git, root)
			}
			if b.err != nil {
				t.Fatalf("baseline capture: %v", b.err)
			}
			switch name {
			case "dirty-worktree":
				mustWrite(t, filepath.Join(root, "a_test.go"), "dirty uncommitted bytes\n")
				dirty := captureProtectedBaseline(ctx, git, root, recordPatterns)
				if dirty.err != nil || dirty.digest != b.digest {
					t.Fatalf("dirty bytes changed committed baseline: before=%s after=%s error=%v", b.digest, dirty.digest, dirty.err)
				}
			case "ignored-snapshot", "snapshot-immutable":
				// Independently enumerate/fold/sort/dedup the candidate output.
				out := recordReadGit(t, git, root, append(strings.Fields(protectedIgnoredCommand), gitCandidatePathspecs(recordPatterns)...)...)
				want := map[string]struct{}{}
				for _, p := range strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00") {
					if matchProtected(recordPatterns, p) {
						want[foldProtectedKey(p)] = struct{}{}
					}
				}
				keys := make([]string, 0, len(want))
				for k := range want {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(strings.Join(keys, "\n")+"\n")))
				if b.ignoredDigest != digest || b.ignoredDigest != recordIgnoredDigest || len(b.ignored) != 2 || !reflect.DeepEqual(b.ignored, want) {
					t.Fatalf("ignored fold snapshot: got %v/%s want %v/%s fixed %s", b.ignored, b.ignoredDigest, want, digest, recordIgnoredDigest)
				}
				if name == "snapshot-immutable" {
					ok := true
					phase := protectedCheck{intact: &ok}
					r1 := buildProtectedRecord(b, recordPatterns, "glob", phase, phase)
					mustWrite(t, filepath.Join(root, "later.test"), "later\n")
					r2 := buildProtectedRecord(b, recordPatterns, "glob", phase, phase)
					r1.Patterns[0] = "changed"
					*r1.IntactBefore = false
					if !reflect.DeepEqual(b.ignored, want) || b.ignoredDigest != digest || r2.BaseIgnoredDigest != digest || r2.BaseIgnoredFiles != 2 || !*r2.IntactBefore || recordPatterns[0] != "*.test" {
						t.Fatal("baseline snapshot or builder aliases changed")
					}
				}
			case "empty-set":
				if b.ignoredDigest != recordEmptyDigest || len(b.ignored) != 0 {
					t.Fatalf("empty ignored snapshot: %s/%d", b.ignoredDigest, len(b.ignored))
				}
			}
		})
	}
}

func TestExecProtectedBaselineFailure(t *testing.T) {
	for _, name := range []string{"missing-git", "non-git", "snapshot-error", "snapshot-over-cap", "construction-timeout", "bad-head", "bad-tree", "tree-nul", "ignored-nul", "preflight-over-cap", "filter-refused"} {
		t.Run(name, func(t *testing.T) {
			git := recordEnv(t)
			oldGit, oldBound := worktreeGit, protectedBaselineTimeout
			protectedBaselineTimeout = 30 * time.Second
			t.Cleanup(func() { worktreeGit, protectedBaselineTimeout = oldGit, oldBound })
			body := "exec '" + git + "' \"$@\""
			switch name {
			case "missing-git":
				worktreeGit = func() (string, error) { return "", errors.New("fixture git missing") }
			case "snapshot-error":
				body = "case \" $* \" in *' ls-files '*) echo snapshot-broken >&2; exit 2;; esac\n" + body
			case "snapshot-over-cap", "preflight-over-cap":
				leg := "ls-files"
				if name == "preflight-over-cap" {
					leg = "config"
				}
				// Non-NUL filler: NUL bytes would fail the record parse with or
				// without the 8 MiB cap, so only overflow proves the cap.
				body = "case \" $* \" in *' " + leg + " '*) head -c 8388609 /dev/zero | tr '\\000' x; exit 0;; esac\n" + body
			case "construction-timeout":
				// The two quick legs start before any deadline even on a loaded
				// rig; the two 1.2 s legs each fit the 2 s bound alone but not
				// together, and fit protectedCheckTimeout, so a capture bounded
				// by the wrong constant, or by a fresh bound per leg, records no
				// error.
				protectedBaselineTimeout = 2 * time.Second
				body = "case \" $* \" in *' config '*) exit 1;; *' rev-parse '*) echo 0000000000000000000000000000000000000000;; *) sleep 1.2;; esac"
			case "bad-head":
				body = "case \" $* \" in *' rev-parse '*) echo not-hex; exit 0;; esac\n" + body
			case "bad-tree":
				body = "case \" $* \" in *' ls-tree '*) printf 'bad\\000'; exit 0;; esac\n" + body
			case "tree-nul", "ignored-nul":
				leg := "ls-tree"
				if name == "ignored-nul" {
					leg = "ls-files"
				}
				body = "case \" $* \" in *' " + leg + " '*) printf unterminated; exit 0;; esac\n" + body
			case "filter-refused":
				body = "case \" $* \" in *' config '*) printf 'filter.bad.clean\\ncat\\000'; exit 0;; esac\n" + body
			}
			fake, log := protectedFakeGit(t, body)
			if name != "missing-git" {
				worktreeGit = func() (string, error) { return fake, nil }
			}
			f := recordHandler(t, func(f *execFixture) {
				if name != "non-git" {
					recordRepo(t, f.worktree)
				}
			})
			if f.h == nil || f.h.baseline == nil || f.h.baseline.err == nil {
				t.Fatalf("handler must construct and retain baseline error: %+v", f.h)
			}
			b := *f.h.baseline
			if b.base != "" || b.digest != "" || b.ignoredDigest != "" || b.files != 0 || b.ignored != nil {
				t.Fatalf("failed baseline trusted partial metadata: %+v", b)
			}
			ok := true
			r := buildProtectedRecord(b, recordPatterns, "glob", protectedCheck{intact: &ok}, protectedCheck{intact: &ok})
			encoded, err := json.Marshal(r)
			if err != nil || r.IntactBefore != nil || r.IntactAfter != nil || r.Error == "" || !bytes.Contains(encoded, []byte(`"intact_before":null`)) || !bytes.Contains(encoded, []byte(`"intact_after":null`)) {
				t.Fatalf("baseline failure must record both null verdicts: %s (%v)", encoded, err)
			}
			if name == "non-git" && r.Error != "worktree is not a git checkout" {
				t.Fatalf("non-git diagnostic: %q", r.Error)
			}
			if strings.HasSuffix(name, "-over-cap") && !errors.Is(b.err, ErrHandlerOverflow) {
				t.Fatalf("over-cap leg must overflow the 8 MiB cap: %v", b.err)
			}
			if name == "snapshot-error" && !strings.Contains(r.Error, "snapshot-broken") {
				t.Fatalf("snapshot error missing: %s", r.Error)
			}
			if name == "construction-timeout" {
				var timeout *HandlerTimeoutError
				if !errors.As(b.err, &timeout) || strings.Count(protectedLog(t, log), "core.fsmonitor=false") < 2 {
					t.Fatalf("shared construction deadline must span multiple legs: %s", r.Error)
				}
			}
			if name == "filter-refused" && (strings.Contains(protectedLog(t, log), "rev-parse") || !strings.Contains(r.Error, WorktreeFilterRefused)) {
				t.Fatalf("preflight did not stop construction: %s", r.Error)
			}
		})
	}
}

func TestExecProtectedWorstCaseFitsOutputCap(t *testing.T) {
	ok := true
	long := strings.Repeat("\x01\x02\x03\x04\x05", 10000)
	before := protectedCheck{intact: &ok, err: errors.New(strings.Repeat("\x01", 1000000))}
	after := protectedCheck{intact: &ok, err: errors.New("second error")}
	for i := 0; i < 40; i++ {
		before.changes = append(before.changes, protectedChange{Path: fmt.Sprint(i) + long, Status: "modified"})
		after.changes = append(after.changes, protectedChange{Path: fmt.Sprint(i) + long, Status: "ignored"})
	}
	patterns := make([]string, 64)
	for i := range patterns {
		patterns[i] = strings.Repeat("p", 250) + fmt.Sprint(i)
	}
	b := protectedBaseline{base: strings.Repeat("a", 64), digest: "sha256:" + strings.Repeat("b", 64), ignoredDigest: "sha256:" + strings.Repeat("c", 64), files: 999999, ignored: map[string]struct{}{"ONE": {}}}
	r := buildProtectedRecord(b, patterns, "dirs-only", before, after)
	out, err := json.Marshal(r)
	if err != nil || len(out) >= 512<<10 {
		t.Fatalf("protected record exceeds cap: %d bytes (%v)", len(out), err)
	}
	if len(r.Changes) != 32 || r.Total != 80 || !r.Truncated {
		t.Fatalf("protected change bound: kept=%d total=%d truncated=%v", len(r.Changes), r.Total, r.Truncated)
	}
	if len(r.Error) > 4096 || r.IntactBefore != nil || r.IntactAfter != nil {
		t.Fatal("error text/verdict bounds")
	}
	for _, c := range r.Changes {
		path, _ := json.Marshal(c.Path)
		if len(path) > worktreeDeltaMaxPathJSON || !c.PathTruncated || c.Phase != "before" {
			t.Fatalf("path or phase bound: bytes=%d cut=%v phase=%s", len(path), c.PathTruncated, c.Phase)
		}
	}
	// A partial before phase must precede after; phase errors affect only their
	// own verdict, and error precedence is baseline, before, after.
	before.changes, before.err = before.changes[:1], nil
	r = buildProtectedRecord(b, patterns, "glob", before, after)
	if r.Changes[1].Phase != "after" || r.IntactBefore == nil || !*r.IntactBefore || r.IntactAfter != nil || r.Error != "second error" {
		t.Fatal("phase ordering/verdict/error precedence")
	}
	before.err = errors.New("first error")
	if r = buildProtectedRecord(b, patterns, "glob", before, after); r.Error != "first error" {
		t.Fatal("before error precedence")
	}
	b.err = errors.New("baseline error")
	if r = buildProtectedRecord(b, patterns, "glob", before, after); r.Error != "baseline error" || r.IntactBefore != nil || r.IntactAfter != nil {
		t.Fatal("baseline error precedence")
	}
}

// TestProtectedV24BoundsApplied pins design V24's bounds (iteration 290: the
// slowest phase 2 501 ms and construction 3 232 ms, each x 1.5 rounded up), so
// reverting either to its provisional 2 500 / 3 200 ms value goes red.
func TestProtectedV24BoundsApplied(t *testing.T) {
	if protectedCheckTimeout != 3751*time.Millisecond || protectedBaselineTimeout != 4848*time.Millisecond {
		t.Fatalf("bounds = %v / %v, want design V24's 3.751s / 4.848s", protectedCheckTimeout, protectedBaselineTimeout)
	}
}

// TestProtectedIgnoredDigestSorted builds the set in reverse order many times:
// a digest over map order instead of sorted keys differs on some iteration.
func TestProtectedIgnoredDigestSorted(t *testing.T) {
	keys := []string{"z.test", "y.test", "x.test", "w.test", "v.test", "u.test", "t.test"}
	want := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("t.test\nu.test\nv.test\nw.test\nx.test\ny.test\nz.test\n")))
	for i := 0; i < 64; i++ {
		set := map[string]struct{}{}
		for _, k := range keys {
			set[k] = struct{}{}
		}
		if got := protectedKeysDigest(set); got != want {
			t.Fatalf("iteration %d: digest %s, want sorted %s", i, got, want)
		}
	}
	if got := protectedKeysDigest(map[string]struct{}{}); got != recordEmptyDigest {
		t.Fatalf("empty set digest %s, want %s", got, recordEmptyDigest)
	}
}

// TestProtectedRecordGolden pins every key and value of the design §4.5 block.
func TestProtectedRecordGolden(t *testing.T) {
	yes, no := true, false
	patterns := []string{"*.test", "tests/**"}
	b := protectedBaseline{base: strings.Repeat("a", 40), digest: "sha256:" + strings.Repeat("b", 64), files: 2,
		ignoredDigest: "sha256:" + strings.Repeat("c", 64), ignored: map[string]struct{}{"k1": {}, "k2": {}}}
	head := `{"digest":"` + protectedDigest(patterns) + `","patterns":["*.test","tests/**"],"enforced":"%s","base":"` + strings.Repeat("a", 40) +
		`","base_digest":"sha256:` + strings.Repeat("b", 64) + `","base_files":2,"base_ignored_digest":"sha256:` + strings.Repeat("c", 64) + `","base_ignored_files":2,`
	for _, c := range []struct {
		name          string
		before, after protectedCheck
		want          string
	}{
		{"changed", protectedCheck{intact: &yes}, protectedCheck{intact: &no, changes: []protectedChange{{Path: "tests/x.py", Status: "modified"}}},
			`"intact_before":true,"intact_after":false,"changes":[{"path":"tests/x.py","status":"modified","phase":"after"}],"truncated":false,"total":1,"error":""}`},
		{"intact", protectedCheck{intact: &yes}, protectedCheck{intact: &yes},
			`"intact_before":true,"intact_after":true,"changes":[],"truncated":false,"total":0,"error":""}`},
	} {
		enforced := map[string]string{"changed": "glob", "intact": "dirs-only"}[c.name]
		want := fmt.Sprintf(head, enforced) + c.want
		got, err := json.Marshal(buildProtectedRecord(b, patterns, enforced, c.before, c.after))
		if err != nil || string(got) != want {
			t.Fatalf("%s record:\n got %s\nwant %s (%v)", c.name, got, want, err)
		}
	}
}

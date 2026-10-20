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
	"runtime"
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

// M6b fixtures use a real index and commit so status can prove a clean tree.
func recordCommittedRepo(t *testing.T, git, root string) {
	t.Helper()
	mustWrite(t, filepath.Join(root, "a_test.go"), "alpha\n")
	mustWrite(t, filepath.Join(root, ".gitignore"), "*.test\n")
	for _, args := range [][]string{{"init", "-q", "--template="}, {"add", "."},
		{"-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "-c", "core.hooksPath=/dev/null", "commit", "-qm", "fixture"}} {
		recordReadGit(t, git, root, args...)
	}
	mustWrite(t, filepath.Join(root, "standing.test"), "ignored at construction\n")
}

func recordSwapGit(t *testing.T, git string) {
	t.Helper()
	old := worktreeGit
	worktreeGit = func() (string, error) { return git, nil }
	t.Cleanup(func() { worktreeGit = old })
}

func recordPhaseBound(t *testing.T, bound time.Duration) {
	t.Helper()
	old := protectedCheckTimeout
	protectedCheckTimeout = bound
	t.Cleanup(func() { protectedCheckTimeout = old })
}

func recordSuccess(t *testing.T, r execResult) *protectedRecord {
	t.Helper()
	if r.ExitCode == nil || *r.ExitCode != 0 || r.Stdout != "[a]" || r.TimedOut {
		t.Fatalf("report-only command result changed: %+v", r)
	}
	if r.Protected == nil {
		t.Fatal("protected result absent")
	}
	return r.Protected
}

func TestExecProtectedIntact(t *testing.T) {
	for _, name := range []string{"intact", "changed", "command-changed"} {
		t.Run(name, func(t *testing.T) {
			git := recordEnv(t)
			recordPhaseBound(t, 15*time.Second) // Identity rows do not test scheduling.
			f := newExecFixture(t, execFixtureOpts{profile: func(m map[string]any) {
				m["protected"] = recordPatterns
				m["commands"].(map[string]any)["modify"] = map[string]any{"argv": []string{"/bin/sh", "-c", "printf changed > a_test.go; printf '[a]'"}}
			}, setup: func(f *execFixture) { recordCommittedRepo(t, git, f.worktree) }})
			base, digest, files := recordIndependent(t, git, f.worktree)
			ignored := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("STANDING.TEST\n")))
			command := "test"
			if name == "changed" {
				mustWrite(t, filepath.Join(f.worktree, "a_test.go"), "host changed\n")
			}
			var r execResult
			if name == "command-changed" {
				command = "modify"
				r = f.mustRun(t, command)
			} else {
				r = f.mustRun(t, command, "a")
			}
			p := recordSuccess(t, r)
			if p.Base != base || p.BaseDigest != digest || p.BaseFiles != files || p.BaseIgnoredDigest != ignored || p.BaseIgnoredFiles != 1 || p.Digest != protectedDigest(recordPatterns) || !reflect.DeepEqual(p.Patterns, recordPatterns) {
				t.Fatalf("independent protected identity mismatch: %+v", p)
			}
			wantEnforced := "glob"
			if runtime.GOOS == "linux" {
				wantEnforced = "dirs-only"
			}
			if p.Enforced != wantEnforced {
				t.Fatalf("enforced = %q, want %q on %s", p.Enforced, wantEnforced, runtime.GOOS)
			}
			if p.IntactBefore == nil || p.IntactAfter == nil || *p.IntactBefore != (name != "changed") || *p.IntactAfter != (name == "intact") || p.Error != "" {
				t.Fatalf("before/after verdict mismatch: %+v", p)
			}
			want := []protectedRecordChange{}
			for _, phase := range []string{"before", "after"} {
				if name == "changed" || (name == "command-changed" && phase == "after") {
					want = append(want, protectedRecordChange{protectedChange: protectedChange{Path: "a_test.go", Status: "modified"}, Phase: phase})
				}
			}
			if !reflect.DeepEqual(p.Changes, want) || p.Total != len(want) || p.Truncated {
				t.Fatalf("phase changes = %+v/%d, want %+v", p.Changes, p.Total, want)
			}
			if f.h.Spawns() != 1 {
				t.Fatalf("phase git counted as command spawn: %d", f.h.Spawns())
			}
		})
	}
}

func recordExecute(t *testing.T, f *execFixture, ctx context.Context, command string) execResult {
	t.Helper()
	payload, _ := json.Marshal(map[string]any{"command": command, "args": []string{}})
	if command == "test" {
		payload = []byte(`{"command":"test","args":["a"]}`)
	}
	out, err := f.h.Execute(ctx, EffectRequest{Effect: EffectWorkspaceExec, Scope: WorkspaceScope}, payload)
	if err != nil {
		t.Fatalf("report-only Execute failed: %v", err)
	}
	return execDecode(t, out)
}

func recordTimeoutError(message string) bool {
	return strings.Contains(message, "timed out") || strings.Contains(message, "deadline exceeded")
}

func TestExecProtectedCheckTimeout(t *testing.T) {
	t.Run("reserve", func(t *testing.T) {
		now := time.Unix(100, 0)
		for _, tc := range []struct {
			name              string
			elapsed, deadline time.Duration
			hasDeadline       bool
			want              time.Duration
		}{
			{"warm", 200 * time.Millisecond, 10 * time.Second, true, 6049 * time.Millisecond},
			{"warm-slow", 650 * time.Millisecond, 10 * time.Second, true, 5599 * time.Millisecond},
			{"cold", 2500 * time.Millisecond, 10 * time.Second, true, 3749 * time.Millisecond},
			{"spent", 7 * time.Second, 10 * time.Second, true, -751 * time.Millisecond},
			{"no-deadline", 0, 0, false, 9 * time.Second},
			{"capped", 0, 30 * time.Second, true, 9 * time.Second},
		} {
			got := execChildBudget(9*time.Second, now.Add(tc.deadline), tc.hasDeadline, now.Add(tc.elapsed), 3751*time.Millisecond)
			if got != tc.want {
				t.Errorf("exact reserve %s: got %s want %s", tc.name, got, tc.want)
			}
		}
		git := recordEnv(t)
		bound := 2 * time.Second
		recordPhaseBound(t, bound)
		f := recordHandler(t, func(f *execFixture) { recordCommittedRepo(t, git, f.worktree) })
		ctx, cancel := context.WithTimeout(context.Background(), 3*bound)
		defer cancel()
		deadline, _ := ctx.Deadline()
		r := recordExecute(t, f, ctx, "sleep")
		if !r.TimedOut || r.Protected == nil {
			t.Fatalf("reserve handler must time out child and retain protected: %+v", r)
		}
		if r.Protected.IntactAfter == nil && (r.Protected.Error == "" || recordTimeoutError(r.Protected.Error)) {
			t.Errorf("reserve must leave after-check time: %+v", r.Protected)
		}
		if !time.Now().Before(deadline) {
			t.Error("reserve handler returned at/after original deadline")
		}
	})
	for _, name := range []string{"shared-phase", "null", "before-only"} {
		t.Run(name, func(t *testing.T) {
			git := recordEnv(t)
			bound := 2 * time.Second
			recordPhaseBound(t, bound)
			f := recordHandler(t, func(f *execFixture) { recordCommittedRepo(t, git, f.worktree) })
			marker := shellQuote(filepath.Join(t.TempDir(), "first-status"))
			body := ""
			sleep := func(d time.Duration) string { return fmt.Sprintf("sleep %.3f\n", d.Seconds()) }
			switch name {
			case "shared-phase":
				body = "if [ \"$1\" = -c ] && { [ \"$3\" = status ] || [ \"$3\" = ls-files ]; }; then\n" + sleep(bound*4/10) + "fi\n"
			case "null":
				body = "if [ \"$1\" = -c ] && [ \"$3\" = status ]; then\n" + sleep(3*bound) + "fi\n"
			case "before-only":
				body = "if [ \"$1\" = -c ] && [ \"$3\" = status ] && [ ! -f " + marker + " ]; then touch " + marker + "; echo BEFORE-ONLY >&2; exit 73; fi\n"
			}
			fake, _ := protectedFakeGit(t, body+"exec "+shellQuote(git)+" \"$@\"")
			recordSwapGit(t, fake) // Baseline is already captured; only phases are injected.
			ctx, cancel := context.WithTimeout(context.Background(), 6*bound)
			defer cancel()
			deadline, _ := ctx.Deadline()
			p := recordSuccess(t, recordExecute(t, f, ctx, "test"))
			if p.IntactBefore != nil || p.Error == "" {
				t.Fatalf("failed before-check must be null with error: %+v", p)
			}
			if name == "before-only" {
				if p.IntactAfter == nil || !*p.IntactAfter || !strings.Contains(p.Error, "BEFORE-ONLY") {
					t.Fatalf("before error must survive successful after-check: %+v", p)
				}
			} else if p.IntactAfter != nil || !recordTimeoutError(p.Error) {
				t.Fatalf("timed-out phases must both be null with timeout error: %+v", p)
			}
			if !time.Now().Before(deadline) {
				t.Fatal("report-only timeout call exceeded original deadline")
			}
		})
	}
	t.Run("after-cap", func(t *testing.T) {
		git := recordEnv(t)
		bound := 3 * time.Second
		recordPhaseBound(t, bound)
		f := recordHandler(t, func(f *execFixture) { recordCommittedRepo(t, git, f.worktree) })
		dir := t.TempDir()
		marker, delta := shellQuote(filepath.Join(dir, "command")), shellQuote(filepath.Join(dir, "delta"))
		// Delay only row 184's first snapshot, which follows budget computation.
		body := fmt.Sprintf("if [ \"$1\" = -C ] && [ ! -f %s ]; then touch %s; sleep %.3f; fi\n", delta, delta, (bound * 11 / 4).Seconds())
		body += fmt.Sprintf("if [ \"$1\" = -c ] && [ \"$3\" = status ] && [ -f %s ]; then sleep %.3f; fi\n", marker, (3 * bound).Seconds())
		fake, _ := protectedFakeGit(t, body+"exec "+shellQuote(git)+" \"$@\"")
		recordSwapGit(t, fake)
		// The command marks the transition; neither baseline nor before-check sleeps.
		m := defaultExecProfile(f)
		m["protected"] = recordPatterns
		m["commands"].(map[string]any)["mark"] = map[string]any{"argv": []string{"/bin/sh", "-c", "touch " + marker + "; printf '[a]'"}}
		data, _ := json.Marshal(m)
		var err error
		f.profile, err = ParseExecProfile(data)
		if err != nil {
			t.Fatal(err)
		}
		f.h, err = NewExecHandler(ExecHandlerConfig{Profile: f.profile, Sandbox: f.sandbox, Episode: "ep1", Worktree: f.worktree, WorkspaceRoot: f.root, StateDir: f.state, OperatorHome: f.home})
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*bound)
		defer cancel()
		deadline, _ := ctx.Deadline()
		p := recordSuccess(t, recordExecute(t, f, ctx, "mark"))
		if p.IntactBefore == nil || !*p.IntactBefore || p.IntactAfter != nil || !recordTimeoutError(p.Error) {
			t.Fatalf("original deadline must cap after-phase: %+v", p)
		}
		if time.Now().After(deadline.Add(bound / 2)) {
			t.Fatal("after-phase extended original deadline by more than cleanup allowance")
		}
	})
	// The before-phase's own time is charged to the child: with the budget
	// read before it, the child ends 0.6 bound later and the 0.7-bound
	// after-phase no longer fits.
	t.Run("reserve-after-before", func(t *testing.T) {
		git := recordEnv(t)
		bound := 3 * time.Second
		recordPhaseBound(t, bound)
		f := recordHandler(t, func(f *execFixture) { recordCommittedRepo(t, git, f.worktree) })
		first := shellQuote(filepath.Join(t.TempDir(), "first-status"))
		sleep := func(d time.Duration) string { return fmt.Sprintf("sleep %.3f", d.Seconds()) }
		fake, _ := protectedFakeGit(t, fmt.Sprintf("if [ \"$1\" = -c ] && [ \"$3\" = status ]; then\n"+
			"if [ -f %s ]; then %s; else touch %s; %s; fi\nfi\nexec %s \"$@\"", first, sleep(bound*7/10), first, sleep(bound*6/10), shellQuote(git)))
		recordSwapGit(t, fake)
		ctx, cancel := context.WithTimeout(context.Background(), 3*bound)
		defer cancel()
		deadline, _ := ctx.Deadline()
		r := recordExecute(t, f, ctx, "sleep")
		if !r.TimedOut || r.Protected == nil {
			t.Fatalf("child must time out inside its reserved budget: %+v", r)
		}
		// Only the after verdict is asserted: a loaded rig may time out the
		// before-phase too, and that spends the same budget.
		if p := r.Protected; p.IntactAfter == nil || !*p.IntactAfter {
			t.Fatalf("after-phase must keep a full bound after a slow before-phase: %+v", p)
		}
		if !time.Now().Before(deadline) {
			t.Fatal("reserve-after-before returned at/after the original deadline")
		}
	})
	// Cancelling the call after the command ran must not abort the after-phase.
	t.Run("after-cancel", func(t *testing.T) {
		git := recordEnv(t)
		bound := 3 * time.Second
		recordPhaseBound(t, bound)
		f := recordHandler(t, func(f *execFixture) { recordCommittedRepo(t, git, f.worktree) })
		dir := t.TempDir()
		first, after := filepath.Join(dir, "first-status"), filepath.Join(dir, "after-status")
		fake, _ := protectedFakeGit(t, fmt.Sprintf("if [ \"$1\" = -c ] && [ \"$3\" = status ]; then\n"+
			"if [ -f %s ]; then touch %s; sleep %.3f; else touch %s; fi\nfi\nexec %s \"$@\"",
			shellQuote(first), shellQuote(after), (bound*3/10).Seconds(), shellQuote(first), shellQuote(git)))
		recordSwapGit(t, fake)
		ctx, cancel := context.WithTimeout(context.Background(), 6*bound)
		defer cancel()
		done := make(chan struct{})
		defer close(done)
		go func() {
			for {
				select {
				case <-done:
					return
				case <-time.After(5 * time.Millisecond):
				}
				if _, err := os.Stat(after); err == nil {
					cancel()
					return
				}
			}
		}()
		p := recordSuccess(t, recordExecute(t, f, ctx, "test"))
		if ctx.Err() != context.Canceled {
			t.Fatalf("instrument: the call was not cancelled during the after-phase (ctx err %v)", ctx.Err())
		}
		if p.IntactAfter == nil || !*p.IntactAfter {
			t.Fatalf("after-phase must survive call cancellation: %+v", p)
		}
	})
}

func TestExecProtectedNotGit(t *testing.T) {
	for _, name := range []string{"not-git", "missing-git"} {
		t.Run(name, func(t *testing.T) {
			git := recordEnv(t)
			fake, log := protectedFakeGit(t, "exec "+shellQuote(git)+" \"$@\"")
			recordSwapGit(t, fake)
			if name == "missing-git" {
				worktreeGit = func() (string, error) { return "", errors.New("git missing") }
			}
			f := recordHandler(t, nil)
			// Baseline errors are sticky even if Git becomes available for Execute.
			worktreeGit = func() (string, error) { return fake, nil }
			mustWrite(t, log, "")
			p := recordSuccess(t, f.mustRun(t, "test", "a"))
			if p.Error == "" || len(p.Error) > 4096 || p.IntactBefore != nil || p.IntactAfter != nil {
				t.Fatalf("baseline failure must be episode error/null: %+v", p)
			}
			if name == "not-git" && !strings.Contains(p.Error, "worktree is not a git checkout") {
				t.Fatalf("not-git baseline error: %s", p.Error)
			}
			if calls := protectedLog(t, log); strings.Contains(calls, "-c core.fsmonitor=false status") || strings.Contains(calls, "ls-files -v") || strings.Contains(calls, ":(top,glob,icase)") {
				t.Fatalf("baseline failure ran protected phase git: %s", calls)
			}
		})
	}
}

func TestExecNoProtectedIsByteIdentical(t *testing.T) {
	for _, name := range []string{"absent", "empty"} {
		t.Run(name, func(t *testing.T) {
			git := recordEnv(t)
			fake, log := protectedFakeGit(t, "exec "+shellQuote(git)+" \"$@\"")
			recordSwapGit(t, fake)
			f := newExecFixture(t, execFixtureOpts{profile: func(m map[string]any) {
				if name == "empty" {
					m["protected"] = []string{}
				}
			}, setup: func(f *execFixture) { recordCommittedRepo(t, git, f.worktree) }})
			out, err := f.call(t, "test", "a")
			if err != nil {
				t.Fatal(err)
			}
			var got map[string]any
			if err := json.Unmarshal(out, &got); err != nil {
				t.Fatal(err)
			}
			if _, ok := got["protected"]; ok {
				t.Fatal("unprotected JSON must omit protected key")
			}
			data, _ := json.Marshal(defaultExecProfile(f)) // Same paths, sandbox and config; no protected key.
			profile, err := ParseExecProfile(data)
			if err != nil {
				t.Fatal(err)
			}
			// Digest identifies exact operator bytes, so [] and absent already
			// differ at M6a. Hold that identity constant to compare exec behavior.
			profile.Digest = f.profile.Digest
			f.h, err = NewExecHandler(ExecHandlerConfig{Profile: profile, Sandbox: f.sandbox, Episode: "ep1", Worktree: f.worktree, WorkspaceRoot: f.root, StateDir: f.state, OperatorHome: f.home})
			if err != nil {
				t.Fatal(err)
			}
			control, err := f.call(t, "test", "a")
			if err != nil {
				t.Fatal(err)
			}
			var want map[string]any
			if err := json.Unmarshal(control, &want); err != nil {
				t.Fatal(err)
			}
			got["duration_ms"], want["duration_ms"] = 0, 0
			gotBytes, _ := json.Marshal(got)
			wantBytes, _ := json.Marshal(want)
			if !bytes.Equal(gotBytes, wantBytes) {
				t.Fatalf("legacy normalized bytes differ: got %s want %s", gotBytes, wantBytes)
			}
			calls := protectedLog(t, log)
			if !strings.Contains(calls, "-C "+f.worktree+" -c core.fsmonitor=false -c core.quotePath=false") {
				t.Fatal("logging positive control: worktree delta git never ran")
			}
			for _, line := range strings.Split(calls, "\n") {
				if strings.HasPrefix(line, "-c core.fsmonitor=false") || strings.Contains(line, ":(top,glob,icase)") {
					t.Fatalf("unprotected profile ran protected git: %s", line)
				}
			}
		})
	}
}

package broker

import (
	"context"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Real Git is only the candidate generator; fixture expectations independently
// pin the mirror and the complete filtered set. No shared subtest deadline.
func TestProtectedGitCandidateEquivalence(t *testing.T) {
	bin, err := exec.LookPath("git")
	if err != nil {
		if os.Getenv("CI") != "" {
			t.Fatal("Git is required in CI:", err)
		}
		t.Skip("Git unavailable locally:", err)
	}
	rows := append([]protectedFixtureRow(nil), protectedFixtureRows...)
	rows = append(rows, []protectedFixtureRow{
		{"greek-question-mark", []string{"a;_test.go"}, "pkg/a\u037e_test.go", true},
		{"greek-varia", []string{"a`_test.go"}, "pkg/a\u1fef_test.go", true},
		{"nfc", []string{"é_test.go"}, "pkg/é_test.go", true},
		{"nfd", []string{"é_test.go"}, "pkg/e\u0301_test.go", true},
		{"qmark-utf8", []string{"?_test.go"}, "pkg/é_test.go", true},
		{"class-utf8", []string{"[^a]_test.go"}, "pkg/é_test.go", true},
		{"sharp-s", []string{"*_test.go"}, "pkg/a_teßt.go", false},
		{"ligature-ff", []string{"buff/**"}, "buﬀ/x.py", false},
		{"ignored-creation", []string{"*_test.go"}, "zzign/a_test.go", true},
		{"ignored-nm-c-txt", []string{"*_test.go"}, "nm/c.txt", false},
		{"directory-replaced-by-file", []string{"tests/**"}, "tests", true},
		{"directory-replaced-by-symlink", []string{"tests/**"}, "tests", true},
		{"deleted", []string{"tests/**"}, "tests/x.py", true},
		{"assume-unchanged", []string{"*_test.go"}, "pkg/a_test.go", true},
		{"skip-worktree", []string{"*_test.go"}, "pkg/a_test.go", true},
	}...)
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			for _, cwd := range []string{"root-CWD", "subdirectory-CWD"} {
				t.Run(cwd, func(t *testing.T) {
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()
					tmp := t.TempDir()
					root, home := filepath.Join(tmp, "repo"), filepath.Join(tmp, "home")
					mustMkdir(t, root)
					mustMkdir(t, home)
					git := func(dir string, args ...string) []byte {
						t.Helper()
						cmd := exec.CommandContext(ctx, bin, append([]string{"-c", "core.fsmonitor=false"}, args...)...)
						cmd.Dir = dir
						cmd.Env = []string{"HOME=" + home, "TMPDIR=" + tmp, "GOTMPDIR=" + tmp, "PATH=/usr/bin:/bin", "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_OPTIONAL_LOCKS=0", "GIT_TERMINAL_PROMPT=0"}
						out, err := cmd.CombinedOutput()
						if err != nil {
							t.Fatalf("git %q: %v: %s", args, err, out)
						}
						return out
					}
					write := func(rel, content string) {
						t.Helper()
						mustMkdir(t, filepath.Dir(filepath.Join(root, rel)))
						mustWrite(t, filepath.Join(root, rel), content)
					}
					git(root, "init", "-q", "--template=")
					for key, value := range map[string]string{"user.name": "fixture", "user.email": "fixture@example.invalid", "core.ignorecase": "false", "core.fsmonitor": "false", "core.precomposeunicode": "false", "core.hooksPath": home} {
						git(root, "config", key, value)
					}
					write("sub/keep.txt", "keep\n")
					git(root, "add", "sub/keep.txt")
					replacement := strings.HasPrefix(row.name, "directory-replaced-by-")
					ignored := strings.HasPrefix(row.name, "ignored-")
					if ignored {
						write(".git/info/exclude", "zzign/\nnm/\n")
						write("nm/c.txt", "standing non-hit\n")
						write("zzign/b.txt", "standing non-hit\n")
						write("zzign/standing_test.go", "standing protected\n")
					} else if replacement {
						// No tracked descendant: losing the bare-D arm cannot hide behind D rows.
						mustMkdir(t, filepath.Join(root, "tests"))
					} else {
						write(row.path, "before\n")
						git(root, "add", "--", row.path)
						if !strings.Contains(string(git(root, "ls-files", "-z")), row.path+"\x00") {
							// Explicit cacheinfo preserves an NFD index spelling on Darwin too.
							blob := strings.TrimSpace(string(git(root, "hash-object", "-w", "--", row.path)))
							git(root, "update-index", "--add", "--cacheinfo", "100644", blob, row.path)
						}
						if !strings.Contains(string(git(root, "ls-files", "-z")), row.path+"\x00") {
							t.Fatalf("index lost raw spelling %q", row.path)
						}
					}
					git(root, "commit", "-q", "-m", "fixture")
					head := strings.TrimSpace(string(git(root, "rev-parse", "HEAD")))
					canon, err := ParseProtectedPatterns(row.patterns)
					if err != nil {
						t.Fatal(err)
					}
					if got := matchProtected(canon, row.path); got != row.want {
						t.Fatalf("mirror(%q)=%v want %v", row.path, got, row.want)
					}
					specs := gitCandidatePathspecs(canon)
					ignoredArgs := strings.Fields(protectedIgnoredCommand)
					baseline := map[string]struct{}{}
					for _, p := range strings.Split(string(git(root, append(append([]string(nil), ignoredArgs...), specs...)...)), "\x00") {
						if p != "" && matchProtected(canon, p) {
							baseline[foldProtectedKey(p)] = struct{}{}
						}
					}
					status := "modified"
					switch row.name {
					case "directory-replaced-by-file":
						if err := os.Remove(filepath.Join(root, "tests")); err != nil {
							t.Fatal(err)
						}
						write("tests", "replacement\n")
						status = "untracked"
					case "directory-replaced-by-symlink":
						if err := os.Remove(filepath.Join(root, "tests")); err != nil {
							t.Fatal(err)
						}
						if err := os.Symlink("sub", filepath.Join(root, "tests")); err != nil {
							t.Fatal(err)
						}
						status = "untracked"
					case "deleted":
						if err := os.Remove(filepath.Join(root, row.path)); err != nil {
							t.Fatal(err)
						}
						status = "deleted"
					case "assume-unchanged", "skip-worktree":
						git(root, "update-index", "--"+row.name, "--", row.path)
						write(row.path, "after\n")
						status = "index-flag"
					default:
						write(row.path, "after\n")
						if ignored {
							status = "ignored"
						}
					}
					dir := root
					prefix := ""
					if cwd == "subdirectory-CWD" {
						dir = filepath.Join(root, "sub")
						prefix = "sub"
					}
					// ls-files reports relative to CWD; status porcelain always reports the root.
					normalize := func(b []byte, flags bool) []byte {
						if prefix == "" || len(b) == 0 {
							return b
						}
						records := strings.Split(strings.TrimSuffix(string(b), "\x00"), "\x00")
						for i, p := range records {
							if flags {
								records[i] = p[:2] + path.Clean(path.Join(prefix, p[2:]))
							} else {
								records[i] = path.Clean(path.Join(prefix, p))
							}
						}
						return []byte(strings.Join(records, "\x00") + "\x00")
					}
					statusArgs := []string{"status", "--porcelain=v2", "-z", "--branch", "--untracked-files=all", "--ignored=no", "--no-renames", "--ignore-submodules=all", "--"}
					statusOut := git(dir, append(statusArgs, specs...)...)
					rawIgnored := git(dir, append(append([]string(nil), ignoredArgs...), specs...)...)
					flagsOut := git(dir, append([]string{"ls-files", "-v", "-z", "--"}, specs...)...)
					if row.name == "ignored-nm-c-txt" && (!reflect.DeepEqual(ignoredArgs, []string{"ls-files", "-z", "-o", "-i", "--exclude-standard", "--"}) || strings.Contains(string(rawIgnored), "! ")) {
						t.Fatalf("ignored source is not per-file ls-files: argv=%q output=%q", ignoredArgs, rawIgnored)
					}
					if row.name == "non-hit-sibling" || row.name == "sharp-s" {
						if !strings.Contains(string(statusOut), row.path+"\x00") {
							t.Fatal("filter control was not selected as a raw candidate")
						}
					}
					gotHead, changes, err := parseProtectedCandidates(canon, statusOut, normalize(rawIgnored, false), normalize(flagsOut, true), baseline)
					if err != nil || gotHead != head {
						t.Fatalf("head=%q want=%q error=%v status=%q", gotHead, head, err, statusOut)
					}
					got, want := map[protectedChange]bool{}, map[protectedChange]bool{}
					for _, c := range changes {
						got[c] = true
					}
					if row.want {
						want[protectedChange{row.path, status}] = true
					}
					if !reflect.DeepEqual(got, want) {
						t.Fatalf("reported=%+v want=%+v; status=%q ignored=%q flags=%q", got, want, statusOut, rawIgnored, flagsOut)
					}
				})
			}
		})
	}
}

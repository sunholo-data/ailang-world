package broker

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func protectedRepo(t *testing.T, ctx context.Context) (string, string, func(...string) string) {
	t.Helper()
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	run := func(args ...string) string {
		t.Helper()
		cmd := exec.CommandContext(ctx, git, args...)
		cmd.Dir, cmd.Env = root, gitEnv()
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %q: %v: %s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	run("init", "-q", "--template=")
	run("config", "user.name", "fixture")
	run("config", "user.email", "fixture@example.invalid")
	run("config", "core.hooksPath", root)
	mustWrite(t, filepath.Join(root, "a_test.go"), "before\n")
	run("add", "a_test.go")
	run("commit", "-qm", "fixture")
	return root, git, run
}

func protectedFakeGit(t *testing.T, body string) (string, string) {
	t.Helper()
	dir := t.TempDir()
	git, log := filepath.Join(dir, "git"), filepath.Join(dir, "spawns")
	mustWrite(t, git, "#!/bin/sh\nprintf '%s\\n' \"$*\" >> '"+log+"'\n"+body+"\n")
	if err := os.Chmod(git, 0700); err != nil {
		t.Fatal(err)
	}
	return git, log
}

func protectedLog(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestProtectedGitFilterPreflight(t *testing.T) {
	for _, name := range []string{"configured-clean", "configured-process", "configured-process-include", "no-filter", "query-error"} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			root, realGit, run := protectedRepo(t, ctx)
			marker := filepath.Join(root, "marker")
			configured := strings.HasPrefix(name, "configured-")
			if configured {
				mustWrite(t, filepath.Join(root, ".gitattributes"), "a_test.go filter=probe\n")
				kind := "clean"
				if strings.Contains(name, "process") {
					kind = "process"
				}
				command := "touch '" + marker + "'; cat"
				if strings.HasSuffix(name, "include") {
					include := filepath.Join(root, "included")
					mustWrite(t, include, "[filter \"probe\"]\n\tprocess = \""+command+"\"\n")
					run("config", "include.path", include)
				} else {
					run("config", "filter.probe."+kind, command)
				}
				mustWrite(t, filepath.Join(root, "a_test.go"), "after!\n")
			}
			body := "exec '" + realGit + "' \"$@\""
			if name == "query-error" {
				body = "case \" $* \" in *' config '*) echo query-broken >&2; exit 2;; esac\n" + body
			}
			git, log := protectedFakeGit(t, body)
			err := protectedGitPreflight(ctx, git, root)
			// The status leg makes a removed refusal observable in M4a alone.
			if err == nil {
				_, _ = runWorktreeGit(ctx, git, root, "status", "--porcelain=v2", "-z")
			}
			spawns := strings.Count(protectedLog(t, log), " status ")
			_, markerErr := os.Stat(marker)
			if name != "no-filter" {
				if spawns != 0 || !os.IsNotExist(markerErr) {
					t.Fatalf("unsafe progression: status spawns=%d marker error=%v", spawns, markerErr)
				}
				if err == nil || (configured && !strings.Contains(err.Error(), WorktreeFilterRefused)) {
					t.Fatalf("missing refusal diagnostic: %v", err)
				}
			} else if err != nil || spawns != 1 {
				t.Fatalf("control: diagnostic=%v status spawns=%d", err, spawns)
			}
		})
	}
}

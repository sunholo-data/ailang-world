package broker

import (
	"context"
	"errors"
	"fmt"
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

// protectedRelaxBound lifts the phase bound for rows that are not about it: a
// loaded runner (go test ./...) measured 2.5 s for one fake-Git spawn.
func protectedRelaxBound(t *testing.T) {
	old := protectedCheckTimeout
	protectedCheckTimeout = 30 * time.Second
	t.Cleanup(func() { protectedCheckTimeout = old })
}

func TestProtectedChecker(t *testing.T) {
	for _, name := range []string{"clean", "modified", "deleted", "untracked", "ignored", "index-flag", "head-moved", "malformed", "filter-refused"} {
		t.Run(name, func(t *testing.T) {
			protectedRelaxBound(t)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			root, git, run := protectedRepo(t, ctx)
			base := run("rev-parse", "HEAD")
			want := name
			switch name {
			case "filter-refused":
				// The checker itself must stop at the preflight: no status, no filter run.
				marker := filepath.Join(root, "marker")
				mustWrite(t, filepath.Join(root, ".gitattributes"), "a_test.go filter=probe\n")
				run("config", "filter.probe.clean", "touch '"+marker+"'; cat")
				mustWrite(t, filepath.Join(root, "a_test.go"), "after!\n")
				fake, log := protectedFakeGit(t, "exec '"+git+"' \"$@\"")
				r := checkProtected(ctx, fake, root, base, []string{"*_test.go"}, nil)
				_, markerErr := os.Stat(marker)
				if r.err == nil || !strings.Contains(r.err.Error(), WorktreeFilterRefused) || r.intact != nil {
					t.Fatalf("filter not refused: %+v", r)
				}
				if spawns := strings.Count(protectedLog(t, log), " status "); spawns != 0 || !os.IsNotExist(markerErr) {
					t.Fatalf("unsafe progression: status spawns=%d marker error=%v", spawns, markerErr)
				}
				return
			case "modified":
				mustWrite(t, filepath.Join(root, "a_test.go"), "changed\n")
			case "deleted":
				if err := os.Remove(filepath.Join(root, "a_test.go")); err != nil {
					t.Fatal(err)
				}
			case "untracked", "ignored":
				if name == "ignored" {
					mustWrite(t, filepath.Join(root, ".gitignore"), "new_test.go\n")
				}
				mustWrite(t, filepath.Join(root, "new_test.go"), "new\n")
			case "index-flag":
				mustWrite(t, filepath.Join(root, "skip_test.go"), "skip\n")
				run("add", "skip_test.go")
				run("commit", "-qm", "second protected file")
				base = run("rev-parse", "HEAD")
				run("update-index", "--assume-unchanged", "a_test.go")
				run("update-index", "--skip-worktree", "skip_test.go")
			case "head-moved":
				run("commit", "--allow-empty", "-qm", "moved")
			case "malformed":
				git, _ = protectedFakeGit(t, "case \" $* \" in *' config '*) exit 1;; esac\nprintf 'garbage'")
			}
			r := checkProtected(ctx, git, root, base, []string{"*_test.go"}, nil)
			if name == "malformed" {
				if r.err == nil || r.intact != nil {
					t.Fatalf("malformed accepted: %+v", r)
				}
				return
			}
			if r.err != nil || r.intact == nil || *r.intact != (name == "clean") {
				t.Fatalf("%s verdict: %+v", name, r)
			}
			expected := map[protectedChange]bool{}
			switch name {
			case "clean":
			case "head-moved":
				expected[protectedChange{"HEAD", want}] = true
			case "index-flag":
				expected[protectedChange{"a_test.go", want}] = true
				expected[protectedChange{"skip_test.go", want}] = true
			case "untracked", "ignored":
				expected[protectedChange{"new_test.go", want}] = true
			default:
				expected[protectedChange{"a_test.go", want}] = true
			}
			for _, c := range r.changes {
				if !expected[c] {
					t.Fatalf("unexpected change: %+v", c)
				}
				delete(expected, c)
			}
			if len(expected) != 0 {
				t.Fatalf("missing changes: %+v", expected)
			}
			if name == "ignored" {
				r = checkProtected(ctx, git, root, base, []string{"*_test.go"}, map[string]struct{}{foldProtectedKey("NEW_TEST.GO"): {}})
				if r.err != nil || r.intact == nil || !*r.intact {
					t.Fatalf("standing ignored snapshot: %+v", r)
				}
			}
		})
	}
}

func TestProtectedCheckerRootAndBound(t *testing.T) {
	for _, name := range []string{"cwd", "shared-phase-bound", "argv", "nonzero", "output-failure", "stderr-failure", "missing-head"} {
		t.Run(name, func(t *testing.T) {
			if name != "shared-phase-bound" {
				protectedRelaxBound(t)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			// Canonical, as snapshotWorktree passes it: the fake git logs a resolved $PWD.
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			base := strings.Repeat("a", 40)
			envLog := filepath.Join(root, "environment")
			t.Setenv("PROTECTED_PARENT_SECRET", "must-not-inherit")
			body := "printf 'cwd=%s\\n' \"$PWD\" >> '" + envLog + "'\n/usr/bin/env >> '" + envLog + "'\n"
			if name == "shared-phase-bound" {
				// Derived from the bound (rule 3m): each of the four calls
				// fits it alone, the four together exceed it by 40%.
				body += fmt.Sprintf("sleep %.3f\n", protectedCheckTimeout.Seconds()*0.35)
			}
			body += "case \" $* \" in *' config '*) exit 1;; esac\n"
			switch name {
			case "nonzero":
				body += "echo deliberate-failure >&2; exit 2\n"
			case "output-failure", "stderr-failure":
				body += "/usr/bin/head -c 8388609 /dev/zero"
				if name == "stderr-failure" {
					body += " >&2"
				}
				body += "\n"
			case "missing-head":
				body += "exit 0\n"
			default:
				body += "case \" $* \" in *' status '*) printf '# branch.oid " + base + "\\000';; esac\n"
			}
			git, log := protectedFakeGit(t, body)
			patterns := []string{"tests/**", "é_test.go"}
			r := checkProtected(ctx, git, root, base, patterns, nil)
			if name != "cwd" && name != "argv" {
				if r.err == nil || r.intact != nil {
					t.Fatalf("%s accepted: %+v", name, r)
				}
				if name == "shared-phase-bound" && !errors.Is(r.err, ErrHandlerTimeout) {
					t.Fatalf("expected phase timeout: %v", r.err)
				}
				if strings.HasSuffix(name, "failure") && !errors.Is(r.err, ErrHandlerOverflow) {
					t.Fatalf("expected output overflow: %v", r.err)
				}
				if name == "nonzero" && !strings.Contains(r.err.Error(), "deliberate-failure") {
					t.Fatalf("missing exit diagnostic: %v", r.err)
				}
				return
			}
			if r.err != nil || r.intact == nil || !*r.intact {
				t.Fatalf("control verdict: %+v", r)
			}
			calls := strings.Split(strings.TrimSpace(protectedLog(t, log)), "\n")
			if len(calls) != 4 || !strings.Contains(calls[0], " config --null --get-regexp ^filter\\.") {
				t.Fatalf("expected preflight plus three data calls: %q", calls)
			}
			data := []string{"status --porcelain=v2 -z --branch --untracked-files=all --ignored=no --no-renames --ignore-submodules=all --", protectedIgnoredCommand, "ls-files -v -z --"}
			for i, args := range data {
				want := "-c core.fsmonitor=false " + args + " :(top,glob,icase)te*t* :(top,glob,icase)te*t*/** :(top,glob,icase)**/*_te*t.go"
				if calls[i+1] != want {
					t.Fatalf("argv[%d]=%q want %q", i+1, calls[i+1], want)
				}
			}
			for _, call := range calls {
				if !strings.Contains(call, "-c core.fsmonitor=false") {
					t.Fatalf("fsmonitor enabled: %s", call)
				}
			}
			environment := protectedLog(t, envLog)
			for _, entry := range append(gitEnv(), "cwd="+root) {
				if strings.Count(environment, entry+"\n") != 4 {
					t.Fatalf("cwd/env missing %q: %s", entry, environment)
				}
			}
			for _, line := range strings.Split(strings.TrimSpace(environment), "\n") {
				key, _, _ := strings.Cut(line, "=")
				if !strings.Contains(" PATH LC_ALL GIT_CONFIG_NOSYSTEM GIT_CONFIG_GLOBAL GIT_OPTIONAL_LOCKS GIT_TERMINAL_PROMPT PWD SHLVL _ cwd ", " "+key+" ") {
					t.Fatalf("inherited environment: %s", line)
				}
			}
		})
	}
}

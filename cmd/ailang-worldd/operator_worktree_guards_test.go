package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func TestOperatorWorktreeSurvivors(t *testing.T) {
	r := newWorktreeStartRig(t, false)
	for _, name := range []string{"store", "episode", "relative"} {
		t.Run(name, func(t *testing.T) {
			before := observeProvision(t, r.f.r, r.creds)
			if name == "relative" {
				if _, err := operatorPhysicalTarget("ep"); !errors.Is(err, store.ErrOperatorConflict) {
					t.Fatalf("relative path admitted: %v", err)
				}
			} else {
				a := r.action
				if name == "store" {
					a.Store += ".other"
				} else {
					a.Episode = "BAD!"
				}
				if err := reserveOperatorStartWorktree(presetTestContext(t), runnerDB(r.creds), a, r.at, r.src, func(string) error {
					t.Fatal("invalid start crossed a worktree boundary")
					return nil
				}); err != store.ErrOperatorConflict {
					t.Fatalf("expected identity refusal before lock and reconstruction: %v", err)
				}
			}
			if !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
				t.Fatal("identity refusal modified store, Git or files")
			}
		})
	}
}

func TestOperatorWorktreeRecoveryGuards(t *testing.T) {
	for _, name := range []string{"store", "branch", "path", "action", "basename", "headless", "observed", "signed_kind", "repo_abs", "common_abs"} {
		t.Run(name, func(t *testing.T) {
			r := newWorktreeStartRig(t, false)
			crash := errors.New("before add")
			err := r.start(t, func(step string) error {
				if step == "before add" && name != "observed" {
					return crash
				}
				return nil
			})
			if name == "observed" {
				if err != nil {
					t.Fatal(err)
				}
				r.state(t, "running", "observed", "created", 2, true)
			} else {
				if !errors.Is(err, crash) {
					t.Fatalf("expected reserved worktree: %v", err)
				}
				a := r.action
				var target operatorWorktreeTarget
				var identity string
				if err := r.db.QueryRow("SELECT identity FROM operator_action_steps WHERE kind='worktree'").Scan(&identity); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(identity), &target); err != nil {
					t.Fatal(err)
				}
				switch name {
				case "store":
					a.Store += ".other"
				case "branch":
					target.Branch += "-other"
				case "path":
					target.Path += "-other"
				case "action":
					target.Action += "-other"
				case "signed_kind":
					a.Kind, a.Task, a.Options.Headless = "finish", "", false
					a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
					if err := a.Validate(); err != nil {
						t.Fatal(err)
					}
				case "repo_abs":
					target.Repo = strings.TrimPrefix(target.Repo, string(filepath.Separator))
				case "common_abs":
					target.Common = strings.TrimPrefix(target.Common, string(filepath.Separator))
				case "headless":
					a.Options.Headless = false
					a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
				case "basename":
					target.Path += "-other"
					a.Resolution = append([]approval.OperatorResolved(nil), a.Resolution...)
					for i := range a.Resolution {
						if a.Resolution[i].Name == "worktree_path" {
							a.Resolution[i].Value += "-other"
						}
					}
				}
				body, err := approval.CanonicalOperatorAction(a)
				if err != nil {
					t.Fatal(err)
				}
				digest := sha256.Sum256(body)
				encoded, err := json.Marshal(target)
				if err != nil {
					t.Fatal(err)
				}
				// Both reads see the same damaged rows. The CAS can agree with them.
				r.exec(t, "UPDATE operator_actions SET action=?, digest=?, kind=?", body, digest[:], a.Kind)
				r.exec(t, "UPDATE operator_action_steps SET identity=? WHERE kind='worktree'", string(encoded))
			}
			before := observeProvision(t, r.f.r, r.creds)
			r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
				t.Error("invalid worktree reservation reached the prompt")
				return nil, errors.New("unexpected prompt")
			}
			want := store.ErrOperatorConflict
			if name == "observed" {
				want = store.ErrOperatorMissing
			}
			if code, msg := r.run(t, worktreeExample); code != exitUsage || !strings.Contains(msg, want.Error()) ||
				!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
				t.Fatalf("invalid worktree reservation changed state or returned wrong class: %d %s", code, msg)
			}
		})
	}
}

func TestOperatorWorktreeReadbackGuards(t *testing.T) {
	for _, name := range []string{"repo_common", "checkout_common", "link_regular", "link_prefix", "registration_count", "relative_registration", "directory", "physical"} {
		t.Run(name, func(t *testing.T) {
			f := newEpisodeRun(t)
			ctx := presetTestContext(t)
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			repo, err := filepath.EvalSymlinks(f.f.repo)
			if err != nil {
				t.Fatal(err)
			}
			target := operatorWorktreeTarget{Repo: repo, Common: filepath.Join(repo, ".git"), Path: filepath.Join(root, "ep"),
				Branch: "world/ep", SHA: presetGit(t, f.f, "rev-parse", "HEAD")}
			git := f.f.git
			setLink := func(dst, src string) {
				t.Helper()
				seedStartFixtureFile(t, dst, readStartFixtureFile(t, src))
			}
			link := filepath.Join(target.Path, ".git")
			wrapper := func(body string) {
				t.Helper()
				git = filepath.Join(root, "git")
				if err := os.WriteFile(git, []byte("#!/bin/sh\n"+body+"\nexec "+shellQuote(f.f.git)+" \"$@\"\n"), 0700); err != nil {
					t.Fatal(err)
				}
			}
			if name == "repo_common" {
				target.Common = root
				if _, err := target.inspect(ctx, git); err == nil || err.Error() != "worktree repository identity differs" {
					t.Fatalf("initial inspect used another repository: %v", err)
				}
				if _, err := os.Lstat(target.Path); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("initial inspection wrote checkout: %v", err)
				}
				return
			}
			presetGit(t, f.f, "worktree", "add", "-b", target.Branch, target.Path, target.SHA)
			if e, err := target.inspect(ctx, git); err != nil || e.state != treeExact {
				t.Fatalf("fixture checkout is not exact: %v", err)
			}
			switch name {
			case "checkout_common":
				other := f.f
				other.repo = filepath.Join(root, "other")
				presetGit(t, f.f, "clone", "--quiet", repo, other.repo)
				decoy := filepath.Join(root, "decoy")
				presetGit(t, other, "worktree", "add", "-b", target.Branch, decoy, target.SHA)
				setLink(link, filepath.Join(decoy, ".git"))
				checkout := f.f
				checkout.repo = target.Path
				if head := presetGit(t, checkout, "rev-parse", "HEAD", "--symbolic-full-name", "HEAD"); head != target.SHA+"\nrefs/heads/"+target.Branch {
					t.Fatalf("decoy would be caught by HEAD: %q", head)
				}
			case "link_regular", "link_prefix":
				saved := filepath.Join(root, "saved-link")
				setLink(saved, link)
				if name == "link_regular" {
					if err := os.Remove(link); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(saved, link); err != nil {
						t.Fatal(err)
					}
					checkout := f.f
					checkout.repo = target.Path
					if head := presetGit(t, checkout, "rev-parse", "HEAD"); head != target.SHA {
						t.Fatalf("symlink fixture has no real Git checkout: %q", head)
					}
				} else {
					// Different reads can see different link contents. The later real Git
					// invocation would succeed if the early gitdir prefix check were skipped.
					seedStartFixtureFile(t, link, []byte("interrupted gitfile\n"))
					wrapper("if [ \"$2\" = " + shellQuote(target.Path) + " ]; then cp " + shellQuote(saved) + " " + shellQuote(link) + " || exit 1; fi")
				}
			case "relative_registration":
				cwd, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				cwd, err = filepath.EvalSymlinks(cwd)
				if err != nil {
					t.Fatal(err)
				}
				relative, err := filepath.Rel(cwd, target.Path)
				if err != nil {
					t.Fatal(err)
				}
				list, err := worktree.Output(ctx, git, []string{"-C", repo, "worktree", "list", "--porcelain", "-z"})
				if err != nil || !strings.Contains(list, "worktree "+target.Path+"\x00") {
					t.Fatalf("cannot read real registration: %v", err)
				}
				list = strings.Replace(list, "worktree "+target.Path+"\x00", "worktree "+relative+"\x00", 1)
				capture := filepath.Join(root, "porcelain")
				seedStartFixtureFile(t, capture, []byte(list))
				wrapper(`if [ "$3" = worktree ]; then exec cat ` + shellQuote(capture) + "; fi")
				before := launchFileSnapshot(t, target.Path)
				e, err := target.inspect(ctx, git)
				if err == nil || err.Error() != "invalid Git worktree registration" || e.state == treeExact ||
					!reflect.DeepEqual(before, launchFileSnapshot(t, target.Path)) {
					t.Fatalf("relative registration admitted or checkout changed: %+v %v", e, err)
				}
			case "registration_count":
				second := filepath.Join(root, "second")
				presetGit(t, f.f, "worktree", "add", "--detach", second, target.SHA)
				secondGitDir := strings.TrimPrefix(strings.TrimSpace(string(readStartFixtureFile(t, filepath.Join(second, ".git")))), "gitdir: ")
				seedStartFixtureFile(t, filepath.Join(secondGitDir, "gitdir"), []byte(link+"\n"))
			case "directory", "physical":
				// List the original registration, then present the changed object after
				// that read. Both branches of this fixture run the actual Git command.
				moved := filepath.Join(root, "moved")
				var replacement string
				if name == "directory" {
					replacement = "touch " + shellQuote(target.Path)
				} else {
					// Replace only the ancestor so the leaf's Lstat remains a directory.
					parent := filepath.Join(root, "parent")
					if err := os.Mkdir(parent, 0700); err != nil {
						t.Fatal(err)
					}
					presetGit(t, f.f, "worktree", "move", target.Path, filepath.Join(parent, "ep"))
					target.Path = filepath.Join(parent, "ep")
					replacement = "ln -s " + shellQuote(moved) + " " + shellQuote(parent)
				}
				victim := target.Path
				if name == "physical" {
					victim = filepath.Dir(target.Path)
				}
				wrapper("if [ \"$3\" = for-each-ref ]; then mv " + shellQuote(victim) + " " + shellQuote(moved) + " && " + replacement + " || exit 1; fi")
			}
			e, err := target.inspect(ctx, git)
			if name == "physical" || name == "relative_registration" {
				want := "worktree physical path differs"
				if name == "relative_registration" {
					want = "invalid Git worktree registration"
				}
				if err == nil || err.Error() != want {
					t.Fatalf("registration or path not refused: %v", err)
				}
			} else if err != nil || e.state != treePartial {
				t.Fatalf("attribution accepted or wrong refusal: state=%d %v", e.state, err)
			}
			if sha := presetGit(t, f.f, "rev-parse", "refs/heads/"+target.Branch); sha != target.SHA {
				t.Fatal("inspection changed pinned branch")
			}
			if name == "link_prefix" && string(readStartFixtureFile(t, link)) != "interrupted gitfile\n" {
				t.Fatal("unchecked link reached later Git read")
			}
			if name == "registration_count" && strings.Count(e.registration, "\x00\x00") != 2 {
				t.Fatalf("expected two real registrations: %q", e.registration)
			}
		})
	}
}

func TestOperatorReadbackErrorIsNeverCreated(t *testing.T) {
	// The inspector's error returns are all before the one assignment of treeExact.
	// Exercise a real command failure at the last possible Git observation.
	f := newEpisodeRun(t)
	repo, err := filepath.EvalSymlinks(f.f.repo)
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	target := operatorWorktreeTarget{Repo: repo, Common: filepath.Join(repo, ".git"), Path: filepath.Join(root, "ep"),
		Branch: "world/ep", SHA: presetGit(t, f.f, "rev-parse", "HEAD")}
	presetGit(t, f.f, "worktree", "add", "-b", target.Branch, target.Path, target.SHA)
	git := filepath.Join(root, "git")
	script := "#!/bin/sh\nif [ \"$4\" = HEAD ]; then exit 1; fi\nexec " + shellQuote(f.f.git) + " \"$@\"\n"
	if err := os.WriteFile(git, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	e, err := target.inspect(presetTestContext(t), git)
	if err != nil || e.state != treePartial || e.common != target.Common {
		t.Fatalf("failed HEAD read must be partial: %+v %v", e, err)
	}
	e, err = target.inspect(presetTestContext(t), f.f.git)
	if err != nil || e.state != treeExact || e.head != target.SHA+"\nrefs/heads/"+target.Branch+"\n" {
		t.Fatalf("control HEAD read must be exact: %+v %v", e, err)
	}
}

func TestOperatorWorktreeLockedReread(t *testing.T) {
	for _, name := range []string{"generation", "error"} {
		t.Run(name, func(t *testing.T) {
			r := newWorktreeStartRig(t, false)
			crash := errors.New("before add")
			if err := r.start(t, func(step string) error {
				if step == "before add" {
					return crash
				}
				return nil
			}); !errors.Is(err, crash) {
				t.Fatalf("expected reserved worktree: %v", err)
			}
			before := observeProvision(t, r.f.r, r.creds)
			want := store.ErrOperatorConflict
			if name == "error" {
				want = errors.New("fixture locked read error")
			}
			calls := 0
			read := func(path, id string) (store.OperatorStepKey, operatorWorktreeTarget, string, error) {
				key, target, episode, err := readRunnerWorktree(path, id)
				calls++
				if name == "error" && calls == 2 && err == nil {
					return key, target, episode, want
				}
				if name == "generation" && calls == 1 && err == nil {
					r.exec(t, "UPDATE operator_actions SET generation=generation+1")
					before = observeProvision(t, r.f.r, r.creds)
				}
				return key, target, episode, err
			}
			r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
				t.Error("a changed reservation reached the prompt")
				return nil, errors.New("unexpected prompt")
			}
			var stderr strings.Builder
			code := r.env.reconcileRunnerWorktree(r.action.ID, io.Discard, &stderr, read)
			if code != exitUsage || calls != 2 || !strings.Contains(stderr.String(), want.Error()) ||
				!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
				t.Fatalf("reread did not preserve the replacement row and its evidence: %d %s", code, &stderr)
			}
		})
	}
}

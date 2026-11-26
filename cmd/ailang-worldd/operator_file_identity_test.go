package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestOperatorFilePublicationIdentityWindows(t *testing.T) {
	t.Run("N04", func(t *testing.T) {
		root, err := filepath.EvalSymlinks(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		dir, moved := filepath.Join(root, "private"), filepath.Join(root, "retained")
		if err := os.Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}
		set := operatorFileSet{Dir: dir, UID: uint64(os.Geteuid())}
		calls := 0
		got := readOperatorFileSetWithOpener(set, "", func(s operatorFileSet) (*os.File, operatorFileProof) {
			calls++
			if calls == 2 {
				if err := os.Rename(dir, moved); err != nil {
					t.Fatal(err)
				}
			}
			return openOperatorFilesDir(s)
		})
		items, listErr := os.ReadDir(moved)
		if calls != 2 || got.Dir.Status != "directory changed during readback" || listErr != nil || len(items) != 0 {
			t.Fatalf("directory swap result %q, calls %d, changed files %v: %v", got.Dir.Status, calls, items, listErr)
		}
	})
	// Restore the signed mint fixture between the two publication windows.
	var r *worktreeStartRig
	var token string
	var minted []byte
	for _, name := range []string{"N12", "N15"} {
		if !t.Run(name, func(child *testing.T) {
			if r == nil {
				r = newWorktreeStartRig(t, false)
				token = mintOperatorFileFixture(t, r)
				if err := r.db.Close(); err != nil {
					child.Fatal(err)
				}
				minted = readStartFixtureFile(t, r.creds)
			}
			t := child
			ctx := presetTestContext(t)
			path, point := r.f.r.sessions, "before ep.task.md"
			if name == "N15" {
				path, point = r.path(), "before observe"
			}
			moved := path + ".retained"
			var before provisionSnapshot
			err := publishOperatorFilesFixture(t, ctx, r, token, r.src.clock, func(at string) error {
				if at != point {
					return nil
				}
				if err := os.Rename(path, moved); err != nil {
					return err
				}
				before = observeProvision(t, r.f.r, r.creds)
				return nil
			})
			if err == nil || err.Error() != "file publication uncertain; "+filesExample {
				t.Fatalf("checkpoint swap refusal: %v", err)
			}
			s, _, err := readRunnerFiles(ctx, r.creds, r.action.ID)
			if err != nil || s.Phase != "reserved" || s.Step.Version != 2 ||
				s.Mint.CredentialID != serviceHash([]byte(token))[len("sha256:"):] {
				t.Fatalf("lost exact uncertainty evidence: %+v %v", s.Step, err)
			}
			after := observeProvision(t, r.f.r, r.creds)
			if before.worktrees == "" || before.worktrees != after.worktrees || before.branches != after.branches ||
				!reflect.DeepEqual(before.workspace, after.workspace) || !reflect.DeepEqual(before.sessions, after.sessions) {
				t.Fatal("refusal changed Git or published files")
			}
			_, missing := os.Lstat(path)
			items, err := os.ReadDir(moved)
			if err != nil || !errors.Is(missing, os.ErrNotExist) || name == "N12" && len(items) != 0 {
				t.Fatalf("moved subject files changed: %v %v", items, err)
			}
			if name == "N12" {
				if err := os.Rename(moved, path); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(r.creds, minted, 0600); err != nil {
					t.Fatal(err)
				}
			}
		}) {
			return
		}
	}
}

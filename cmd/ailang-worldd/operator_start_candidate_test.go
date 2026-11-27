package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorStartCandidateProcess(t *testing.T) {
	const holder = "WORLD_TEST_START_CANDIDATE_TERMINAL"
	if path := os.Getenv(holder); path != "" {
		fd, err := terminalEpisodeLock(path, "ep", nil)
		if err != nil {
			t.Fatal(err)
		}
		defer fd.Close()
		fmt.Println("holding canonical terminal lock")
		if _, err := io.Copy(io.Discard, os.Stdin); err != nil {
			t.Fatal(err)
		}
		return
	}
	creds := filepath.Join(t.TempDir(), "credentials.db")
	canonical, err := store.CanonicalPath(creds)
	if err != nil {
		t.Fatal(err)
	}
	lock := claimFixtureEpisode(t, creds, "ep")
	lock.Close()
	credSQL(t, creds, "UPDATE operator_actions SET state='confirmed', owner=NULL")
	aliasPath := filepath.Join(t.TempDir(), "alias.db")
	if err := os.Symlink(creds, aliasPath); err != nil {
		t.Fatal(err)
	}
	bin, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "-test.run=^TestOperatorStartCandidateProcess$")
	cmd.Env = append(os.Environ(), holder+"="+creds, "GORACE="+strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
	cmd.WaitDelay = time.Second
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	waited := false
	t.Cleanup(func() {
		if !waited {
			cancel()
			stdin.Close()
			cmd.Wait()
		}
	})
	reader := bufio.NewReader(stdout)
	if line, err := reader.ReadString('\n'); err != nil || line != "holding canonical terminal lock\n" {
		t.Fatalf("child barrier: %q %v", line, err)
	}
	before := serviceFileHash(creds)
	worktree := t.TempDir()
	files := t.TempDir()
	treeBefore, filesBefore := provisionFiles(t, worktree), provisionFiles(t, files)
	expectUnchanged := func() {
		t.Helper()
		row := candidateRow(t, creds, 1012)
		if serviceFileHash(creds) != before || row.State != "confirmed" || row.Owner != "" ||
			!reflect.DeepEqual(provisionFiles(t, worktree), treeBefore) || !reflect.DeepEqual(provisionFiles(t, files), filesBefore) {
			t.Fatal("candidate probe changed ownership, files or worktree")
		}
	}
	now := func() int64 { return 1012 }
	err = withOperatorStartCandidate(ctx, aliasPath, now, store.OpenCredentialsReadOnly, store.OpenOperatorEpisodeLock, func(*store.OperatorEpisodeLock, store.OperatorActionRow) error {
		t.Fatal("visit ran while the real terminal process held the lock")
		return nil
	})
	if !errors.Is(err, store.ErrOperatorConflict) || !errors.Is(err, syscall.EWOULDBLOCK) {
		t.Fatalf("terminal holder returned wrong refusal: %v", err)
	}
	expectUnchanged()
	stdin.Close()
	if _, err := io.Copy(io.Discard, reader); err != nil {
		t.Fatal(err)
	}
	err = cmd.Wait()
	waited = true
	if err != nil {
		t.Fatalf("terminal release: %v", err)
	}
	handled := errors.New("test visit stopped")
	var acquired *store.OperatorEpisodeLock
	err = withOperatorStartCandidate(ctx, aliasPath, now, store.OpenCredentialsReadOnly, func(path, episode string) (*store.OperatorEpisodeLock, error) {
		var err error
		acquired, err = store.OpenOperatorEpisodeLock(path, episode)
		return acquired, err
	}, func(lock *store.OperatorEpisodeLock, row store.OperatorActionRow) error {
		fd, err := lock.Verify(canonical, "ep")
		if err != nil {
			t.Fatal(err)
		}
		physical, err := filepath.EvalSymlinks(fd.Name())
		want, wantErr := filepath.EvalSymlinks(operatorEpisodeLockPath(canonical, "ep"))
		if err != nil || wantErr != nil || physical != want || lock != acquired || row.ID != fixtureActionID {
			t.Fatalf("visit did not receive exact locked candidate: %v %v", err, wantErr)
		}
		if err := runnerWrite(canonical, func(*store.Credentials) error { return nil }); err != nil {
			t.Fatalf("candidate held the SQLite writer across the visit: %v", err)
		}
		if contender, err := runnerKernelLock(fd.Name()); !errors.Is(err, syscall.EWOULDBLOCK) {
			if contender != nil {
				contender.Close()
			}
			t.Fatalf("visit lost kernel exclusion: %v", err)
		}
		return handled
	})
	if !errors.Is(err, handled) {
		t.Fatalf("visit outcome: %v", err)
	}
	expectUnchanged()
	requireProvisionLockReleased(t, canonical, "ep")
}

func candidateRow(t *testing.T, path string, now int64) store.OperatorActionRow {
	t.Helper()
	return candidateRowWithContext(t, boundedTestContext(t), path, now)
}

func candidateRowWithContext(t *testing.T, ctx context.Context, path string, now int64) store.OperatorActionRow {
	t.Helper()
	db, err := store.OpenCredentialsReadOnly(runnerDB(path))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	row, err := db.OperatorAction(ctx, fixtureActionID, now)
	if err != nil {
		t.Fatal(err)
	}
	return row
}

func TestOperatorStartCandidateRefusals(t *testing.T) {
	for _, tc := range []struct {
		name, before, interleave string
		advance                  bool
		want                     error
	}{
		{"changed_generation", "", "UPDATE operator_actions SET generation=generation+1", false, store.ErrOperatorConflict},
		{"expired_under_lock", "", "", true, store.ErrOperatorMissing},
		{"empty", "UPDATE operator_actions SET state='requested', approved_by=NULL, artifact=NULL", "", false, store.ErrOperatorMissing},
		{"invalid_episode", "UPDATE operator_actions SET episode_id='ep/invalid'", "", false, store.ErrOperatorConflict},
	} {
		t.Run(tc.name, func(t *testing.T) {
			creds := filepath.Join(t.TempDir(), "credentials.db")
			held := claimFixtureEpisode(t, creds, "ep")
			held.Close()
			canonical, err := store.CanonicalPath(creds)
			if err != nil {
				t.Fatal(err)
			}
			credSQL(t, creds, "UPDATE operator_actions SET state='confirmed', owner=NULL")
			if tc.before != "" {
				credSQL(t, creds, tc.before)
			}
			before, hash := candidateRow(t, creds, 1012), serviceFileHash(creds)
			worktree, files := t.TempDir(), t.TempDir()
			treeBefore, filesBefore := provisionFiles(t, worktree), provisionFiles(t, files)
			now := int64(1012)
			var acquired bool
			err = withOperatorStartCandidate(boundedTestContext(t), creds, func() int64 { return now }, store.OpenCredentialsReadOnly,
				func(path, episode string) (*store.OperatorEpisodeLock, error) {
					acquired = true
					if tc.advance {
						now = 1600
					}
					if tc.interleave != "" {
						credSQL(t, creds, tc.interleave)
						before, hash = candidateRow(t, creds, 1012), serviceFileHash(creds)
					}
					return store.OpenOperatorEpisodeLock(path, episode)
				}, func(*store.OperatorEpisodeLock, store.OperatorActionRow) error {
					t.Fatal("refused candidate reached its visit")
					return nil
				})
			if !errors.Is(err, tc.want) || acquired != (tc.interleave != "" || tc.advance) ||
				!reflect.DeepEqual(candidateRow(t, creds, 1012), before) || serviceFileHash(creds) != hash ||
				!reflect.DeepEqual(provisionFiles(t, worktree), treeBefore) || !reflect.DeepEqual(provisionFiles(t, files), filesBefore) {
				t.Fatalf("refusal changed evidence or returned wrong class: %v", err)
			}
			requireProvisionLockReleased(t, canonical, "ep")
		})
	}
	t.Run("missing_parent", func(t *testing.T) {
		root := t.TempDir()
		before := provisionFiles(t, root)
		missing := filepath.Join(root, "absent", "credentials.db")
		err := withOperatorStartCandidate(boundedTestContext(t), missing, func() int64 { return 1012 }, store.OpenCredentialsReadOnly,
			func(string, string) (*store.OperatorEpisodeLock, error) {
				t.Fatal("canonical path failure reached the lock")
				return nil, nil
			}, nil)
		if !errors.Is(err, os.ErrNotExist) || !reflect.DeepEqual(provisionFiles(t, root), before) {
			t.Fatalf("missing parent: %v", err)
		}
	})
	t.Run("missing_store", func(t *testing.T) {
		root := t.TempDir()
		before := provisionFiles(t, root)
		err := withOperatorStartCandidate(boundedTestContext(t), filepath.Join(root, "credentials.db"), func() int64 { return 1012 }, store.OpenCredentialsReadOnly,
			func(string, string) (*store.OperatorEpisodeLock, error) {
				t.Fatal("missing store reached the lock")
				return nil, nil
			}, nil)
		if err == nil || !strings.Contains(err.Error(), "unable to open database file") || !reflect.DeepEqual(provisionFiles(t, root), before) {
			t.Fatalf("missing store: %v", err)
		}
	})
}

func TestOperatorStartCandidateClosesCredentials(t *testing.T) {
	creds := filepath.Join(t.TempDir(), "credentials.db")
	claimFixtureEpisode(t, creds, "ep").Close()
	credSQL(t, creds, "UPDATE operator_actions SET state='confirmed', owner=NULL")
	before := serviceFileHash(creds)
	ctx := boundedTestContext(t)
	var db *store.Credentials
	open := func(dsn string) (*store.Credentials, error) {
		var err error
		db, err = store.OpenCredentialsReadOnly(dsn)
		return db, err
	}
	entered := false
	err := withOperatorStartCandidate(ctx, creds, func() int64 { return 1012 }, open, store.OpenOperatorEpisodeLock,
		func(_ *store.OperatorEpisodeLock, row store.OperatorActionRow) error {
			entered = row.ID == fixtureActionID
			return nil
		})
	if db != nil {
		t.Cleanup(func() { db.Close() })
	}
	if err != nil || !entered || db == nil || serviceFileHash(creds) != before {
		t.Fatalf("candidate did not visit the confirmed action read-only: %v", err)
	}
	if _, err := db.NextOperatorStart(ctx, 1012); err == nil || err.Error() != "sql: database is closed" {
		t.Fatalf("candidate did not close its credentials handle: %v", err)
	}
}

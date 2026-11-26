package store

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestOperatorEpisodeLockIdentity(t *testing.T) {
	root := t.TempDir()
	credentials, err := CanonicalPath(filepath.Join(root, "credentials.db"))
	if err != nil {
		t.Fatal(err)
	}
	path := OperatorEpisodeLockPath(credentials, "ep")
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(root, alias); err != nil {
		t.Fatal(err)
	}
	l, err := OpenOperatorEpisodeLock(filepath.Join(alias, "credentials.db"), "ep")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	fd, err := l.Verify(credentials, "ep")
	if err != nil || fd != l.file {
		t.Fatalf("expected exact descriptor: %v", err)
	}
	refuse := func(t *testing.T, lease *OperatorEpisodeLock, episode string, reason string) {
		t.Helper()
		before, beforeErr := os.Lstat(path)
		got, err := lease.Verify(credentials, episode)
		after, afterErr := os.Lstat(path)
		if got != nil || !errors.Is(err, ErrOperatorConflict) || err.Error() != ErrOperatorConflict.Error()+": "+reason ||
			beforeErr != nil || afterErr != nil || !os.SameFile(before, after) || before.Mode() != after.Mode() {
			t.Fatalf("lock refusal changed inode or returned wrong class: %v", err)
		}
	}
	t.Run("nil", func(t *testing.T) { refuse(t, nil, "ep", "wrong episode lock") })
	t.Run("zero", func(t *testing.T) { refuse(t, &OperatorEpisodeLock{}, "ep", "wrong episode lock") })
	t.Run("wrong_key", func(t *testing.T) { refuse(t, l, "other", "wrong episode lock") })
	t.Run("uid", func(t *testing.T) {
		before, err := l.file.Stat()
		if err != nil {
			t.Fatal(err)
		}
		got, err := l.verify(credentials, "ep", os.Getuid()+1)
		after, statErr := l.file.Stat()
		if got != nil || !errors.Is(err, ErrOperatorConflict) || statErr != nil || !os.SameFile(before, after) {
			t.Fatalf("wrong owner changed the inode or was admitted: %v %v", err, statErr)
		}
	})
	t.Run("permissions", func(t *testing.T) {
		if err := os.Chmod(path, 0644); err != nil {
			t.Fatal(err)
		}
		refuse(t, l, "ep", "unsafe episode lock")
		if err := os.Chmod(path, 0600); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("hardlink", func(t *testing.T) {
		alias := filepath.Join(root, "alias")
		if err := os.Link(path, alias); err != nil {
			t.Fatal(err)
		}
		refuse(t, l, "ep", "unsafe episode lock")
		if err := os.Remove(alias); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Rename(path, path+".held"); err != nil {
		t.Fatal(err)
	}
	t.Run("removed", func(t *testing.T) {
		got, err := l.Verify(credentials, "ep")
		_, statErr := os.Lstat(path)
		if got != nil || !errors.Is(err, ErrOperatorConflict) || !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("missing pathname recreated or admitted: %v %v", err, statErr)
		}
	})
	t.Run("replaced", func(t *testing.T) {
		if err := os.WriteFile(path, nil, 0600); err != nil {
			t.Fatal(err)
		}
		refuse(t, l, "ep", "episode lock pathname changed")
	})
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}
	if got, err := l.Verify(credentials, "ep"); got != nil || !errors.Is(err, ErrOperatorConflict) || !errors.Is(err, os.ErrClosed) {
		t.Fatalf("closed descriptor admitted: %v", err)
	}
	if next, err := OpenOperatorEpisodeLock(credentials, "ep"); err != nil {
		t.Fatalf("owner did not release the lock: %v", err)
	} else {
		next.Close()
	}
}

func operatorLockFDCount(t *testing.T, f *os.File) int {
	t.Helper()
	var want syscall.Stat_t
	if err := syscall.Fstat(int(f.Fd()), &want); err != nil {
		t.Fatal(err)
	}
	dir, err := os.Open("/dev/fd")
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, name := range names {
		n, err := strconv.Atoi(name)
		if err != nil {
			continue
		}
		var have syscall.Stat_t
		if syscall.Fstat(n, &have) == nil && have.Dev == want.Dev && have.Ino == want.Ino {
			count++
		}
	}
	return count
}

func TestOperatorEpisodeLockAcquireRefusals(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	credentials := filepath.Join(root, "credentials.db")
	for _, ep := range []string{"", "sub/ep", "sub\\ep", "sub\x00ep"} {
		got, err := OpenOperatorEpisodeLock(credentials, ep)
		if got != nil || !errors.Is(err, ErrOperatorConflict) {
			t.Fatalf("unsafe key %q admitted: %v", ep, err)
		}
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatalf("rejected keys left effects: %v %v", entries, err)
	}
	if got, err := OpenOperatorEpisodeLock(filepath.Join(root, "missing", "credentials.db"), "ep"); got != nil || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing store parent was created or returned wrong class: %v", err)
	}
	path := credentials + ".operator-episode-ep.lock" // legacy terminal key
	target := filepath.Join(root, "target")
	if err := os.WriteFile(target, []byte("target"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
	if got, err := OpenOperatorEpisodeLock(credentials, "ep"); got != nil || !errors.Is(err, syscall.ELOOP) {
		t.Fatalf("expected no-follow refusal: %v", err)
	}
	if link, err := os.Readlink(path); err != nil || link != target {
		t.Fatalf("symlink refusal altered lock name: %q %v", link, err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("unsafe"), 0644); err != nil {
		t.Fatal(err)
	}
	if got, err := OpenOperatorEpisodeLock(credentials, "ep"); got != nil || !errors.Is(err, ErrOperatorConflict) {
		t.Fatalf("unsafe lock mode admitted: %v", err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if count := operatorLockFDCount(t, f); count != 1 {
		t.Fatalf("refused unsafe lock leaked a descriptor: %d", count)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "unsafe" {
		t.Fatalf("refusal altered lock bytes: %q %v", data, err)
	}
}

func TestOperatorEpisodeLockFIFORefused(t *testing.T) {
	credentials := filepath.Join(t.TempDir(), "credentials.db")
	canonical, err := CanonicalPath(credentials)
	if err != nil {
		t.Fatal(err)
	}
	path := OperatorEpisodeLockPath(canonical, "ep")
	if err := syscall.Mkfifo(path, 0600); err != nil {
		t.Fatal(err)
	}
	before, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(path, os.O_RDWR|syscall.O_NONBLOCK, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	// Some kernels reject flock on a FIFO before the stat checks. Exercise the
	// regular-file check with that real FIFO descriptor on every kernel.
	got, err := (&OperatorEpisodeLock{f}).Verify(canonical, "ep")
	after, statErr := os.Lstat(path)
	if got != nil || !errors.Is(err, ErrOperatorConflict) || statErr != nil || !os.SameFile(before, after) {
		t.Fatalf("FIFO lock changed or was admitted: %v %v", err, statErr)
	}
	if acquired, err := OpenOperatorEpisodeLock(canonical, "ep"); acquired != nil ||
		(!errors.Is(err, syscall.ENOTSUP) && !errors.Is(err, ErrOperatorConflict)) {
		t.Fatalf("FIFO entered the acquisition path: %v", err)
	}
}

func TestOperatorEpisodeLockContendsWithoutWaiting(t *testing.T) {
	const childKey = "WORLD_TEST_STORE_EPISODE_LOCK"
	if path := os.Getenv(childKey); path != "" {
		canonical, err := CanonicalPath(path)
		if err != nil {
			t.Fatal(err)
		}
		witness, err := os.Open(OperatorEpisodeLockPath(canonical, "ep"))
		if err != nil {
			t.Fatal(err)
		}
		defer witness.Close()
		if count := operatorLockFDCount(t, witness); count != 1 {
			t.Fatalf("descriptor census positive control: %d", count)
		}
		got, err := OpenOperatorEpisodeLock(path, "ep")
		if got != nil {
			got.Close()
		}
		if got != nil || !errors.Is(err, syscall.EWOULDBLOCK) {
			t.Fatalf("unexpected independent acquisition: %v", err)
		}
		if count := operatorLockFDCount(t, witness); count != 1 {
			t.Fatalf("refused contention leaked a descriptor: %d", count)
		}
		return
	}
	path := filepath.Join(t.TempDir(), "credentials.db")
	l, err := OpenOperatorEpisodeLock(path, "ep")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	bin, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	child := exec.CommandContext(ctx, bin, "-test.run=^TestOperatorEpisodeLockContendsWithoutWaiting$")
	child.Env = append(os.Environ(), childKey+"="+path, "GORACE="+strings.TrimSpace(os.Getenv("GORACE")+" atexit_sleep_ms=0"))
	child.WaitDelay = time.Second
	if out, err := child.CombinedOutput(); err != nil {
		t.Fatalf("bounded contention failed: %v %s", err, out)
	}
}

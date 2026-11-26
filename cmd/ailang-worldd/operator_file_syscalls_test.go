package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestOperatorFileSyscalls(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(root, "sessions-fixture")
	if err := os.Mkdir(base, 0700); err != nil {
		t.Fatal(err)
	}
	set := operatorFileSet{Dir: base, UID: uint64(os.Geteuid())}
	dir, _ := openOperatorFilesDir(set)
	if dir == nil {
		t.Fatal("cannot open private fixture")
	}
	defer dir.Close()
	body := []byte(operatorStartTask)
	cloexec := func(t *testing.T, fd int) {
		t.Helper()
		got, err := unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0)
		if err != nil || got&unix.FD_CLOEXEC == 0 {
			t.Errorf("descriptor can leak into a child: %d %v", got, err)
		}
	}
	if _, err := installOperatorFileAt(dir, "task", body, func(fd int, name string, flags int, mode uint32) (int, error) {
		f, err := unix.Openat(fd, name, flags, mode)
		if err == nil {
			cloexec(t, f)
		}
		return f, err
	}); err != nil {
		t.Fatal(err)
	}
	for alias, target := range map[string]string{"link": "task", "dirlink": filepath.Base(base)} {
		parent := base
		if alias == "dirlink" {
			parent = root
		}
		if err := os.Symlink(target, filepath.Join(parent, alias)); err != nil {
			t.Fatal(err)
		}
	}
	before := launchFileSnapshot(t, root)
	refuses := func(t *testing.T, fd int, name string, flags int) {
		t.Helper()
		if got, err := unix.Openat(fd, name, flags, 0); err == nil {
			unix.Close(got)
			t.Errorf("kernel opened %s with publication flags", name)
		}
	}
	t.Run("directory_open", func(t *testing.T) {
		probed := false
		f, p := openOperatorFilesDirAt(set, func(fd int, name string, flags int, mode uint32) (int, error) {
			if name == filepath.Base(base) && !probed {
				probed = true
				refuses(t, fd, "dirlink", flags&^unix.O_DIRECTORY)
				refuses(t, int(dir.Fd()), "task", flags)
			}
			child, err := unix.Openat(fd, name, flags, mode)
			if err == nil {
				cloexec(t, child)
			}
			return child, err
		})
		if f == nil || !probed || p.Status != "private directory" {
			t.Fatalf("directory positive control failed: %+v", p)
		}
		f.Close()
	})
	t.Run("read_open", func(t *testing.T) {
		opened, stats := 0, 0
		p, data := readOperatorFileAt(dir, "task", serviceHash(body), set.UID,
			func(fd int, name string, flags int, mode uint32) (int, error) {
				opened++
				refuses(t, fd, "link", flags)
				raw, err := unix.Openat(fd, name, flags, mode)
				if err == nil {
					cloexec(t, raw)
					got, err := unix.FcntlInt(uintptr(raw), unix.F_GETFL, 0)
					if err != nil || got&unix.O_NONBLOCK == 0 {
						t.Errorf("blocking read-open: %d %v", got, err)
					}
				}
				return raw, err
			}, func(fd int, name string, st *unix.Stat_t, flags int) error {
				stats++
				var link unix.Stat_t
				if err := unix.Fstatat(fd, "link", &link, flags); err != nil || link.Mode&unix.S_IFMT != unix.S_IFLNK {
					t.Errorf("stat %d followed the alias: %o %v", stats, link.Mode, err)
				}
				return unix.Fstatat(fd, name, st, flags)
			})
		if p.Status != "exact" || string(data) != string(body) || opened != 1 || stats != 2 {
			t.Fatalf("independent readback did not run: %s %d/%d", p.Status, opened, stats)
		}
	})
	t.Run("create_alias", func(t *testing.T) {
		// O_CREAT|O_EXCL itself refuses a symlink to an existing file.
		if _, err := installOperatorFile(dir, "link", []byte("foreign")); !errors.Is(err, os.ErrExist) {
			t.Fatalf("exclusive create did not refuse an alias: %v", err)
		}
	})
	if !reflect.DeepEqual(before, launchFileSnapshot(t, root)) {
		t.Fatal("syscall checks changed fixture evidence")
	}
}

func TestOperatorFileLockFlags(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lock")
	hold, err := openRunnerLock(path, os.O_CREATE)
	if err != nil {
		t.Fatal(err)
	}
	defer hold.Close()
	other, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	if err := unix.Flock(int(other.Fd()), unix.LOCK_SH|unix.LOCK_NB); !errors.Is(err, unix.EWOULDBLOCK) {
		t.Fatalf("the episode lock allows a shared entrant: %v", err)
	}
	out := make(chan error, 1)
	go func() {
		second, err := openRunnerLock(path, 0)
		if second != nil {
			second.Close()
		}
		out <- err
	}()
	select {
	case err := <-out:
		if !errors.Is(err, unix.EWOULDBLOCK) {
			t.Fatalf("unexpected competing lock result: %v", err)
		}
	case <-time.After(time.Second):
		hold.Close()
		t.Fatal("episode lock blocked instead of refusing")
	}
}

func TestOperatorFileReopen(t *testing.T) {
	for _, phase := range []string{"open", "last_stat", "unreadable"} {
		t.Run(phase, func(t *testing.T) {
			base, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(base, 0700); err != nil {
				t.Fatal(err)
			}
			d, _ := openOperatorFilesDir(operatorFileSet{Dir: base, UID: uint64(os.Geteuid())})
			if d == nil {
				t.Fatal("missing private directory")
			}
			defer d.Close()
			body := []byte(operatorStartTask)
			for _, name := range []string{"task", "replacement"} {
				if _, err := installOperatorFile(d, name, body); err != nil {
					t.Fatal(err)
				}
			}
			replace := func(fd int, name string) {
				t.Helper()
				if err := unix.Renameat(fd, "replacement", fd, name); err != nil {
					t.Fatal(err)
				}
			}
			stats := 0
			p, got := readOperatorFileAt(d, "task", serviceHash(body), uint64(os.Geteuid()),
				func(fd int, name string, flags int, mode uint32) (int, error) {
					switch phase {
					case "open":
						return unix.Openat(fd, "replacement", flags, mode)
					case "unreadable":
						return -1, unix.EACCES
					}
					return unix.Openat(fd, name, flags, mode)
				}, func(fd int, name string, st *unix.Stat_t, flags int) error {
					stats++
					if phase == "last_stat" && stats == 2 {
						replace(fd, name)
					}
					return unix.Fstatat(fd, name, st, flags)
				})
			if p.Status != "unattributed file" || got != nil || string(readStartFixtureFile(t, filepath.Join(base, "task"))) != string(body) {
				t.Fatalf("unattributed file received a proof or changed: %s", p.Status)
			}
		})
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

const lockHelpAddress = "127.0.0.1:8765"

func writeLockHelpRecord(t *testing.T, home, storePath string, uninstalled bool) {
	t.Helper()
	dir := filepath.Join(home, ".ailang-world")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(map[string]any{
		"schema": "world/service/v1", "store": storePath,
		"address": lockHelpAddress, "uninstalled": uninstalled,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "service.json"), body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func serviceLockHelp(verb string) []string {
	return []string{
		fmt.Sprintf("world-publish %s: the world writer lock is held by the World service (%s)", verb, lockHelpAddress),
		"  why: publishing writes the store, and only one writer may hold it",
		fmt.Sprintf("  next: ailang-worldd service stop; world-publish %s …; ailang-worldd service start", verb),
	}
}

func TestWriterLockHelpNamesTheService(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := t.TempDir()
	storePath := filepath.Join(dir, "world.db")
	if err := os.WriteFile(storePath, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "linked")
	if err := os.Symlink(dir, link); err != nil {
		t.Fatal(err)
	}
	writeLockHelpRecord(t, home, filepath.Join(link, "world.db"), false)
	if got, want := writerLockHelp("transitions", storePath), serviceLockHelp("transitions"); !reflect.DeepEqual(got, want) {
		t.Fatalf("help = %q, want %q", got, want)
	}
}

func TestWriterLockHelpOtherStore(t *testing.T) {
	for _, name := range []string{"no-record", "other-store", "uninstalled", "bad-json", "wrong-schema", "unreadable"} {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			storePath := filepath.Join(t.TempDir(), "world.db")
			switch name {
			case "other-store":
				writeLockHelpRecord(t, home, storePath+".other", false)
			case "uninstalled":
				writeLockHelpRecord(t, home, storePath, true)
			case "bad-json", "wrong-schema":
				writeLockHelpRecord(t, home, storePath, false)
				body := "{"
				if name == "wrong-schema" {
					body = `{"schema":"world/service/v2","store":"` + storePath + `"}`
				}
				if err := os.WriteFile(filepath.Join(home, ".ailang-world", "service.json"), []byte(body), 0o600); err != nil {
					t.Fatal(err)
				}
			case "unreadable":
				// A directory cannot be read as a record, even under elevated test permissions.
				if err := os.MkdirAll(filepath.Join(home, ".ailang-world", "service.json"), 0o700); err != nil {
					t.Fatal(err)
				}
			}
			want := []string{
				"world-publish transitions: the world writer lock is held by another process",
				"  why: publishing writes the store, and only one writer may hold it",
				"  next: stop the daemon serving " + storePath + ", publish, then start it again",
			}
			if got := writerLockHelp("transitions", storePath); !reflect.DeepEqual(got, want) {
				t.Fatalf("help = %q, want %q", got, want)
			}
		})
	}
}

func assertServiceWriterRefusal(t *testing.T, verb string, res armResult) {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(res.stderr, "\n"), "\n")
	want := serviceLockHelp(verb)
	if res.code != exitStop || len(lines) != 5 || !reflect.DeepEqual(lines[:3], want) ||
		lines[3] != "STOP fence=store reason=unopenable" || !strings.HasPrefix(lines[4], "  opening ") {
		t.Fatalf("refusal = (%d, %q), want three help lines then the store STOP", res.code, res.stderr)
	}
}

func TestTransitionsWriterLockedNamesTheService(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	storePath := filepath.Join(t.TempDir(), "world.db")
	db, err := store.Open(storePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	writeLockHelpRecord(t, home, storePath, false)
	res := driveTransitions(t, map[string]string{"store": storePath, "manifest": "unused.json"}, "", noEnv)
	assertServiceWriterRefusal(t, "transitions", res)
}

func TestEpochNominateWriterLockedNamesTheService(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	flags := nominateFixture(t)
	db, err := store.Open(flags["store"])
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	writeLockHelpRecord(t, home, flags["store"], false)
	res := driveNominate(t, nominateInv(flags, "new-epoch"), "", noEnv, noTTYProbe(t))
	assertServiceWriterRefusal(t, "epoch nominate", res)
}

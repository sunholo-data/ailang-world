package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorTerminalCreateLockMissingParent(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Join(root, "gone")
	if err := os.Mkdir(parent, 0700); err != nil {
		t.Fatal(err)
	}
	canonical, err := store.CanonicalPath(filepath.Join(parent, "credentials.db"))
	if err != nil {
		t.Fatal(err)
	}
	// This is the window after canonicalization but before lock open.
	if err := os.Remove(parent); err != nil {
		t.Fatal(err)
	}
	before := provisionFiles(t, root)
	lock, err := openCanonicalTerminalEpisodeLock(canonical, "ep", nil, os.O_CREATE)
	if lock != nil || !errors.Is(err, os.ErrNotExist) || !reflect.DeepEqual(provisionFiles(t, root), before) {
		t.Fatalf("creating lock accepted the missing parent or changed it: %v %v", lock, err)
	}
	lock, err = openCanonicalTerminalEpisodeLock(canonical, "ep", nil, 0)
	if lock != nil || err != nil || !reflect.DeepEqual(provisionFiles(t, root), before) {
		t.Fatalf("read-only legacy absence control: %v %v", lock, err)
	}
}

func TestOperatorWaitReadBudget(t *testing.T) {
	dir := t.TempDir()
	before := provisionFiles(t, dir)
	c := runnerConfig{Credentials: filepath.Join(dir, "credentials.db")}
	clock := time.Unix(1012, 0)
	env := serviceEnv{wait: 250 * time.Millisecond,
		runnerNow:   func() time.Time { return clock },
		runnerSleep: func(step time.Duration) { clock = clock.Add(step) },
	}
	var budgets []time.Duration
	read := func(path string, budget time.Duration) (store.OperatorHeartbeat, error) {
		t.Helper()
		budgets = append(budgets, budget)
		if path != c.Credentials || len(budgets) > 3 {
			t.Fatalf("unexpected heartbeat read: %s %v", path, budgets)
		}
		return store.OperatorHeartbeat{}, store.ErrOperatorRunner
	}
	want := []time.Duration{250 * time.Millisecond, 150 * time.Millisecond, 50 * time.Millisecond}
	if env.waitForRunner(c, 0, read) || !reflect.DeepEqual(budgets, want) ||
		!reflect.DeepEqual(provisionFiles(t, dir), before) {
		t.Fatalf("wait did not refuse with the remaining read budget: %v", budgets)
	}
}

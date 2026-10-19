package main

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

// testCredsByPath caches one writable credentials handle per store: its
// writer lock is real, so a second OpenCredentials in this process would be
// refused.
var (
	testCredsMu     sync.Mutex
	testCredsByPath = map[string]*store.Credentials{}
)

// credsFor is the terminal side of a test: the writable credentials store
// beside the world store db (row 189), closed when t ends. A daemon serving db
// reads it per request, so the writer may stay open while it serves.
func credsFor(t testing.TB, db string) *store.Credentials {
	t.Helper()
	path, err := store.DefaultCredentialsPath(db)
	if err != nil {
		t.Fatal(err)
	}
	testCredsMu.Lock()
	defer testCredsMu.Unlock()
	if c, ok := testCredsByPath[path]; ok {
		return c
	}
	c, err := store.OpenCredentials(path)
	if err != nil {
		t.Fatalf("open the credentials store beside %s: %v", db, err)
	}
	testCredsByPath[path] = c
	t.Cleanup(func() {
		testCredsMu.Lock()
		delete(testCredsByPath, path)
		testCredsMu.Unlock()
		_ = c.Close()
	})
	return c
}

// credRows reads every credential beside db, read-only; none when the
// credentials store does not exist yet.
func credRows(t testing.TB, db string) []store.SessionRow {
	t.Helper()
	path, err := store.DefaultCredentialsPath(db)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	ro, err := store.OpenCredentialsReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := ro.ListSessions(ctx, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// daemonView is how a daemon serving db resolves tokens (row 189): a
// read-only open of the credentials store beside it per lookup.
func daemonView(t testing.TB, db string) authority.Resolver {
	t.Helper()
	path, err := store.DefaultCredentialsPath(db)
	if err != nil {
		t.Fatal(err)
	}
	return authority.New(authority.ReadOnlyAt(path))
}

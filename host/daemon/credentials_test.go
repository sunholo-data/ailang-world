package daemon

import (
	"sync"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

// testCredsByDaemon caches the terminal-side credentials writer per daemon:
// one writable handle per credentials store, because its writer lock is real
// (a second OpenCredentials in this process would be refused).
var (
	testCredsMu       sync.Mutex
	testCredsByDaemon = map[*Daemon]*store.Credentials{}
)

// testCreds is what `session new` holds at the terminal in these tests: the
// writable credentials store d resolves from (row 189), opened beside the
// running daemon's world writer lock and closed when t ends.
func testCreds(t testing.TB, d *Daemon) *store.Credentials {
	t.Helper()
	testCredsMu.Lock()
	defer testCredsMu.Unlock()
	if c, ok := testCredsByDaemon[d]; ok {
		return c
	}
	c, err := store.OpenCredentials(d.credentialsPath)
	if err != nil {
		t.Fatalf("open the credentials store beside the daemon: %v", err)
	}
	testCredsByDaemon[d] = c
	t.Cleanup(func() {
		testCredsMu.Lock()
		delete(testCredsByDaemon, d)
		testCredsMu.Unlock()
		_ = c.Close()
	})
	return c
}

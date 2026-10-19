package daemon

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// resolveCredentialsPath is Config.CredentialsPath, defaulting to the
// credentials store beside the world store (row 189). A credentials store that
// IS the world store is refused: the two need separate writer locks.
func resolveCredentialsPath(cfg Config) (string, error) {
	path := cfg.CredentialsPath
	if path == "" {
		var err error
		if path, err = store.DefaultCredentialsPath(cfg.DBPath); err != nil {
			return "", err
		}
	}
	path, err := store.CanonicalPath(path)
	if err != nil {
		return "", fmt.Errorf("--credentials-db: %w", err)
	}
	same, err := store.SameDatabase(cfg.DBPath, path)
	if err != nil {
		return "", err
	}
	if same {
		return "", fmt.Errorf("--credentials-db %s is the world store itself; session credentials need their own file", path)
	}
	return path, nil
}

// checkCredentialsOutsideRoot refuses a credentials store (or its directory,
// which holds its lock and journal) inside the canonical workspace root: an
// agent writing that file could create a credential no human confirmed.
func checkCredentialsOutsideRoot(credPath, canonicalRoot string) error {
	for _, p := range []string{filepath.Dir(credPath), credPath} {
		if err := broker.CheckPolicyOutsideRoot(p, canonicalRoot); err != nil {
			return &WorkspaceStateError{What: "the credentials store is inside --workspace-root", Err: err}
		}
	}
	return nil
}

// ensureCredentialsStore makes sure the per-request read-only opens will find
// an initialized, current credentials store. A current one is only read; a
// missing one is created, and a version-1 one (before the exec profile
// binding, row 189 M2) upgraded, with a short writable open that is closed
// again at once, so the daemon never keeps the credentials writer lock. Only
// serve upgrades: from here on a daemon of an older build would refuse the
// store, and this daemon is the one serving it. A writer holding the lock
// right now (a `session new` at the terminal) means it exists already; a
// version-1 store then stays version 1 until the next start, which costs
// nothing but bound mints (refused until then, naming the restart).
func ensureCredentialsStore(path string, errLog io.Writer) error {
	if ro, err := store.OpenCredentialsReadOnly(path); err == nil {
		current := ro.Version() == store.CredentialsSchemaVersion
		if err := ro.Close(); err != nil || current {
			return err
		}
	}
	upgraded, err := store.UpgradeCredentials(path)
	if err != nil {
		if store.IsWriterAlreadyActive(err) {
			fmt.Fprintf(errLog, "ailang-worldd: the credentials store %s is busy (a session verb is writing it); "+
				"if it predates the exec profile binding it is upgraded at the next serve start\n", path)
			return nil
		}
		return err
	}
	if upgraded {
		fmt.Fprintf(errLog, "ailang-worldd: upgraded the credentials store %s to schema %d (exec profile binding; every existing credential unchanged and unbound)\n",
			path, store.CredentialsSchemaVersion)
	}
	return nil
}

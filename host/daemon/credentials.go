package daemon

import (
	"fmt"
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
// an initialized credentials store. An existing one is only read; a missing
// one is created with a short writable open that is closed again at once, so
// the daemon never keeps the credentials writer lock. A writer holding it
// right now (a `session new` at the terminal) means it exists already.
func ensureCredentialsStore(path string) error {
	if ro, err := store.OpenCredentialsReadOnly(path); err == nil {
		return ro.Close()
	}
	c, err := store.OpenCredentials(path)
	if err != nil {
		if store.IsWriterAlreadyActive(err) {
			return nil
		}
		return err
	}
	return c.Close()
}

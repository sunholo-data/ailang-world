package authority

import (
	"context"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/store"
)

// ReadOnlyAt is the serving daemon's CredentialStore (row 189, D-WORLD-85 C):
// the path of a credentials store, opened READ-ONLY for each lookup and closed
// after it. The daemon holds no credentials handle and never the credentials
// writer lock, so `session new/mint/revoke` run at the terminal while it
// serves, and because no handle outlives a request, a mint or a revoke
// committed there is what the very next request sees — no restart.
//
// It only reads. Nothing reachable from the daemon creates, extends or alters
// a credential: the writable OpenCredentials is called only by the attended
// session verbs (after the /dev/tty fence and y/N) and by the one-time move.
type ReadOnlyAt string

// ResolveSession opens the store, performs the single point lookup and closes
// it. The deadline is checked first, so a deadline-free caller is refused
// before any file is opened (B4).
func (p ReadOnlyAt) ResolveSession(ctx context.Context, credentialID string) (store.SessionRow, bool, error) {
	if ctx == nil {
		return store.SessionRow{}, false, store.ErrNoDeadline
	}
	if _, ok := ctx.Deadline(); !ok {
		return store.SessionRow{}, false, store.ErrNoDeadline
	}
	c, err := store.OpenCredentialsReadOnly(string(p))
	if err != nil {
		return store.SessionRow{}, false, fmt.Errorf("authority: open the credentials store read-only: %w", err)
	}
	defer func() { _ = c.Close() }()
	return c.ResolveSession(ctx, credentialID)
}

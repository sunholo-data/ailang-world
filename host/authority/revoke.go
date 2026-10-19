package authority

import (
	"context"

	"github.com/sunholo-data/ailang-world/host/store"
)

// Revoke deletes the mapping row for a credential_id (D4 delete-revokes). A
// deleted row means the next resolve of that credential is an UNKNOWN
// credential — there is nothing left to distinguish it from a token that never
// existed. Revocation is total and immediate; the credentials store applies
// the DELETE in a transaction, and a serving daemon's next request (a fresh
// read-only handle, ReadOnlyAt) no longer finds the row. Revoking a
// credential_id that has no row is a no-op, not an error.
func Revoke(ctx context.Context, st *store.Credentials, credentialID string) error {
	return st.RevokeSession(ctx, credentialID)
}

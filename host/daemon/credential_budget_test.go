package daemon

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

// heldCredentials is a CredentialStore over the daemon's REAL store whose
// lookup waits for ctx to end (a held sole connection). Its escape hatch then
// answers from the real store, so a ctx-ignoring mutant resolves the valid
// token and is RED on status and elapsed time, not a hang.
type heldCredentials struct {
	real *store.Store
	err  error // non-nil: fail at once with this error instead of holding
}

func (h heldCredentials) ResolveSession(ctx context.Context, id string) (store.SessionRow, bool, error) {
	if h.err != nil {
		return store.SessionRow{}, false, h.err
	}
	select {
	case <-ctx.Done():
		return store.SessionRow{}, false, ctx.Err()
	case <-time.After(seamEscape):
		rctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return h.real.ResolveSession(rctx, id)
	}
}

// TestCredentialLookupExpiryIsNotADenial is AC6 (row 23 M4, B4): a VALID
// token whose lookup outlives the budget gets 503 Timeout, never 401; a lookup
// that fails for another reason is a sanitized 500, never 401.
func TestCredentialLookupExpiryIsNotADenial(t *testing.T) {
	t.Run("held_lookup_is_503_timeout_within_budget", func(t *testing.T) {
		d := newHandlerDaemon(t)
		genesis := seedGenesisEmbedded(t, d, "cred")
		auth := authHeader(t, d) // a real, valid session
		d.resolver = authority.New(heldCredentials{real: d.store})
		d.credentialBudget = 100 * time.Millisecond
		start := time.Now()
		rec := requestRecorderAuth(t, d, auth, http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(testCommit(genesis, 1, "cred"))))
		elapsed := time.Since(start)
		body := assertErrorClass(t, rec, http.StatusServiceUnavailable, "Timeout")
		if !strings.Contains(body.Error.Message, "credential lookup") {
			t.Fatalf("message = %q; want it to name the credential lookup", body.Error.Message)
		}
		if elapsed > d.credentialBudget+500*time.Millisecond {
			t.Fatalf("answered after %v; bound %v + 500ms", elapsed, d.credentialBudget)
		}
		if got := headOf(t, d); got != genesis.Ref.String() {
			t.Fatalf("head moved to %s behind an unresolved credential", got)
		}
	})
	t.Run("failed_lookup_is_sanitized_500", func(t *testing.T) {
		d := newHandlerDaemon(t)
		genesis := seedGenesisEmbedded(t, d, "cred-err")
		auth := authHeader(t, d)
		var errLog bytes.Buffer
		d.errLog = &errLog
		d.resolver = authority.New(heldCredentials{err: errors.New("private lookup detail")})
		rec := requestRecorderAuth(t, d, auth, http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(testCommit(genesis, 1, "cred-err"))))
		assertErrorClass(t, rec, http.StatusInternalServerError, "Internal")
		if strings.Contains(rec.Body.String(), "private lookup detail") || !strings.Contains(errLog.String(), "private lookup detail") {
			t.Fatalf("500 must be sanitized on the wire and logged: body=%s log=%q", rec.Body, errLog.String())
		}
	})
}

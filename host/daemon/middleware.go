package daemon

import (
	"net/http"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
)

// SessionMiddleware wraps the daemon mux so that requests in the protected set
// require a valid, unexpired, known session credential before their handler
// runs (w-session-authority D6). It lives in the daemon package (round-2
// quorum fix) so the denial→HTTP mapping can call the package-local
// writeAPIError directly — host/authority must NOT import host/daemon, and the
// placement here makes that edge impossible.
//
// It is route-agnostic by construction: a `protected func(*http.Request) bool`
// decides which requests carry the session boundary, so flipping the read routes
// on later is a predicate change, not a rewrite (residual R1). This sprint
// protects only POST /v1/commit.
//
// The middleware NEVER reads the X-World-Session header (D-WORLD-26: rejected,
// even as a fallback) and never treats the static serve-api key as a session
// (constraint i — a Bearer value here is a session credential and never an API
// key). Fail-closed typed denials map to HTTP: absent/unknown/expired → 401
// (distinct classes SessionAbsent/SessionUnknown/SessionExpired), malformed →
// 400 InvalidSession. On success the resolved *authority.SessionBinding is
// carried into the request context via the typed authority.WithBinding /
// authority.FromContext pair.
type SessionMiddleware struct {
	resolver authority.Resolver
}

// NewSessionMiddleware builds a middleware over an authority.Resolver.
func NewSessionMiddleware(resolver authority.Resolver) *SessionMiddleware {
	return &SessionMiddleware{resolver: resolver}
}

// Wrap returns a handler that enforces the session boundary on requests for
// which protected(r) is true; all other requests pass through to mux untouched
// (residual R1 — the read routes pass unauthenticated this sprint).
func (m *SessionMiddleware) Wrap(protected func(*http.Request) bool, mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !protected(r) {
			mux.ServeHTTP(w, r)
			return
		}
		header := r.Header.Get("Authorization")
		out := m.resolver.Resolve(header, time.Now().Unix())
		if out.Denied != nil {
			writeSessionDenial(w, *out.Denied)
			return
		}
		ctx := authority.WithBinding(r.Context(), out.Success)
		mux.ServeHTTP(w, r.WithContext(ctx))
	})
}

// writeSessionDenial is THE ONE session-denial -> HTTP mapping: Wrap uses it
// for /v1/commit, and the daemon injects it into host/projection's Config.Deny
// so a card-route denial is byte-identical to the middleware's (B1) — the
// APIError envelope, its classes, its constant messages and the F2 status
// mapping (absent/unknown/expired -> 401, malformed -> 400) live in exactly
// one place, and host/projection never formats the envelope itself (AC1).
func writeSessionDenial(w http.ResponseWriter, k authority.DenialKind) {
	switch k {
	case authority.DenialAbsent, authority.DenialUnknown, authority.DenialExpired:
		// All three fail closed to 401 Unauthorized but keep distinct
		// APIError classes (F6) so a consumer can tell a bad token from
		// an expired one.
		writeAPIError(w, classFor(k), msgFor(k), http.StatusUnauthorized)
	case authority.DenialMalformed:
		// The header is present but is not a valid Bearer credential.
		writeAPIError(w, classFor(k), msgFor(k), http.StatusBadRequest)
	}
}

// classFor maps a denial kind to its operator-facing APIError class. It reuses
// DenialKind.String() so the middleware, the resolver's diagnostics and the
// wire share one vocabulary: SessionAbsent / InvalidSession / SessionUnknown /
// SessionExpired.
func classFor(k authority.DenialKind) string { return k.String() }

// msgFor returns a fixed, operator-useful message per denial kind. It carries
// no token or credential-hash material (D3/F9) and is a constant string per
// kind (never derived from header content).
func msgFor(k authority.DenialKind) string {
	switch k {
	case authority.DenialAbsent:
		return "a session credential is required: no Authorization Bearer header was present"
	case authority.DenialMalformed:
		return "malformed Authorization header: expected a Bearer <64-hex-credential>"
	case authority.DenialUnknown:
		return "unknown session credential: no session matches this token"
	case authority.DenialExpired:
		return "session credential has expired"
	default:
		return "session denied"
	}
}

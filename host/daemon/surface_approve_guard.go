package daemon

// Row 201 A3 sprint 2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.3, §2.4 R2-R5): the
// approval routes' guards, as pure functions of the header values and the surface origin set O.
// Only the new approval routes apply them; every existing route stays as row 199 leaves it.

import (
	"errors"
	"fmt"
	"mime"
	"slices"
	"strings"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// surfaceOrigins is the surface origin set O (design §2.3 [r1]): the loopback origins at port
// (the port the request arrived on; none when unknown), every --surface-origin, and every active
// passkey's origin. It does not depend on which format is POSTed.
func surfaceOrigins(port string, configured, enrolled []string) []string {
	var o []string
	if port != "" {
		o = append(o, "http://localhost:"+port, "http://127.0.0.1:"+port)
	}
	o = append(o, configured...)
	return append(o, enrolled...)
}

// checkSurfaceOrigin is --surface-origin's grammar: A2's ceremony origin
// (https://<host>.ts.net[:port] or http(s)://localhost[:port], no default port, no IP literal).
func checkSurfaceOrigin(origin string) error {
	if _, err := approval.ParseCeremonyOrigin(origin); err != nil {
		return fmt.Errorf("--surface-origin %q: %w", origin, err)
	}
	return nil
}

// allowedApprovalHost reports whether host (the Host header, any case) is exactly the host[:port]
// of an origin in o (R2: DNS rebinding sends its own name).
func allowedApprovalHost(host string, o []string) bool {
	host = strings.ToLower(host)
	for _, origin := range o {
		if _, hostPort, _ := strings.Cut(origin, "://"); host == hostPort {
			return true
		}
	}
	return false
}

// allowedApprovalOrigin reports whether origin (the Origin header) is exactly a member of o (R3);
// a missing Origin and "null" never are.
func allowedApprovalOrigin(origin string, o []string) bool { return slices.Contains(o, origin) }

// sameOriginFetch is R4: Sec-Fetch-Site, when present, is exactly same-origin (two localhost
// ports are same-site, design V17).
func sameOriginFetch(site string) bool { return site == "" || site == "same-origin" }

// jsonMediaType is R5: the media type is application/json, so a browser preflights every
// cross-origin POST, and the daemon answers no preflight with CORS.
func jsonMediaType(contentType string) bool {
	mt, _, err := mime.ParseMediaType(contentType)
	return err == nil && mt == "application/json"
}

// isRequestRace reports a record refused because the request changed after the checks (R14): it
// was consumed, expired, deleted, or its clock or digest moved, between the judge and the write.
func isRequestRace(err error) bool {
	for _, race := range []error{store.ErrRequestConsumed, store.ErrRequestExpired, store.ErrRequestClockBehind, store.ErrRequestDigestChanged, store.ErrRequestUnknown} {
		if errors.Is(err, race) {
			return true
		}
	}
	return false
}

// misdirectedReason is R2's refusal text; an unlisted *.ts.net Host names the setting that adds it.
func misdirectedReason(host string) string {
	if name, _, _ := strings.Cut(strings.ToLower(host), ":"); strings.HasSuffix(name, ".ts.net") {
		return fmt.Sprintf("this daemon does not serve the approval surface at %s; to approve there, start serve with --surface-origin https://%s", host, strings.ToLower(host))
	}
	return fmt.Sprintf("this daemon does not serve the approval surface at host %q", host)
}

// crossOriginReason is R3's refusal text; a well-formed *.ts.net origin names the setting that adds it.
func crossOriginReason(origin string) string {
	if rp, err := approval.ParseCeremonyOrigin(origin); err == nil && strings.HasSuffix(rp, ".ts.net") {
		return fmt.Sprintf("the origin %s is not one this daemon serves the approval surface at; to approve from it, start serve with --surface-origin %s", origin, origin)
	}
	return fmt.Sprintf("the origin %q is not one this daemon serves the approval surface at (loopback, a configured surface origin or an enrolled passkey's origin)", origin)
}

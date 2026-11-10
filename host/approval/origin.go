package approval

// Row 201 A2 (design §2.2): the one origin a passkey ceremony may run at, and its rpID. Only
// localhost and a tailnet *.ts.net host, never an IP literal; rpID is the host exactly.

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrOrigin     = errors.New("approval: a ceremony origin is scheme://host[:port], with https for *.ts.net and http or https for localhost")
	ErrOriginHost = errors.New("approval: a ceremony origin's host must be localhost or a tailnet *.ts.net name (an IP address cannot be a WebAuthn origin: use http://localhost:<port>)")
)

var originHost = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)

// ParseCeremonyOrigin accepts exactly https://<host>.ts.net[:port] and http(s)://localhost[:port]
// (O1-O3) and returns the rpID, which is the host.
func ParseCeremonyOrigin(s string) (rpID string, err error) {
	scheme, rest, ok := strings.Cut(s, "://")
	host, port, hasPort := strings.Cut(rest, ":")
	if !ok || (scheme != "https" && scheme != "http") || len(host) > 253 || !originHost.MatchString(host) { // O1
		return "", ErrOrigin
	}
	if host != "localhost" && !strings.HasSuffix(host, ".ts.net") { // O2
		return "", ErrOriginHost
	}
	if scheme == "http" && host != "localhost" { // O1: http is for localhost only
		return "", fmt.Errorf("%w (http is accepted for localhost only)", ErrOrigin)
	}
	if hasPort {
		n, perr := strconv.Atoi(port)
		if perr != nil || n < 1 || n > 65535 || port != strconv.Itoa(n) { // O3: a bad port
			return "", fmt.Errorf("%w (the port must be 1 to 65535 with no leading zero)", ErrOrigin)
		}
		if (scheme == "https" && n == 443) || (scheme == "http" && n == 80) { // O3: a default port
			return "", fmt.Errorf("%w (leave out the default port :%d: the browser's origin omits it, so it could never match)", ErrOrigin, n)
		}
	}
	return host, nil
}

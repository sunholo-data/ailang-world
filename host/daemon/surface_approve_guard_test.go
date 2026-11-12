package daemon

// Row 201 A3 sprint 2 M4c (design §2.3, §2.4 R2-R5, §6 RG1): the approval routes' guards and the
// surface origin set O, pure.

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

// RG1: O's three members; the Host, Origin, Sec-Fetch-Site and media-type predicates; the
// --surface-origin grammar, at the predicate and at New.
func TestApprovalGuards(t *testing.T) {
	o := surfaceOrigins("7644", []string{"https://s.tail0.ts.net"}, []string{"https://p.tail0.ts.net:8443"})
	if want := []string{"http://localhost:7644", "http://127.0.0.1:7644", "https://s.tail0.ts.net", "https://p.tail0.ts.net:8443"}; !reflect.DeepEqual(o, want) {
		t.Fatalf("O = %v, want %v", o, want)
	}
	if got := surfaceOrigins("", nil, nil); len(got) != 0 {
		t.Fatalf("O with no port, flag or passkey = %v, want empty", got)
	}
	loopbackOnly := surfaceOrigins("7644", nil, nil)
	for _, tc := range []struct {
		host string
		in   []string
		want bool
	}{
		{"localhost:7644", o, true}, {"127.0.0.1:7644", o, true}, {"LOCALHOST:7644", o, true},
		{"s.tail0.ts.net", o, true}, {"S.Tail0.TS.net", o, true}, {"p.tail0.ts.net:8443", o, true},
		{"rebind.example:7644", o, false}, {"localhost:7645", o, false}, {"localhost", o, false}, {"localhost:76440", o, false},
		{"127.0.0.1:7644.evil", o, false}, {"s.tail0.ts.net.evil", o, false}, {"evil.s.tail0.ts.net", o, false},
		{"p.tail0.ts.net", o, false}, {"s.tail0.ts.net:443", o, false}, {"", o, false}, {"[::1]:7644", o, false},
		{"s.tail0.ts.net", loopbackOnly, false}, {"127.0.0.1:7644", loopbackOnly, true},
	} {
		if got := allowedApprovalHost(tc.host, tc.in); got != tc.want {
			t.Errorf("Host %q in %v = %v, want %v", tc.host, tc.in, got, tc.want)
		}
	}
	for _, tc := range []struct {
		origin string
		in     []string
		want   bool
	}{
		{"http://localhost:7644", o, true}, {"http://127.0.0.1:7644", o, true}, {"https://s.tail0.ts.net", o, true}, {"https://p.tail0.ts.net:8443", o, true},
		{"", o, false}, {"null", o, false}, {"http://localhost:7645", o, false}, {"https://localhost:7644", o, false},
		{"http://s.tail0.ts.net", o, false}, {"https://other.tail0.ts.net", o, false}, {"HTTP://LOCALHOST:7644", o, false},
		{"http://localhost:7644/", o, false}, {"https://p.tail0.ts.net", o, false},
		{"https://s.tail0.ts.net", loopbackOnly, false}, {"http://127.0.0.1:7644", loopbackOnly, true},
	} {
		if got := allowedApprovalOrigin(tc.origin, tc.in); got != tc.want {
			t.Errorf("Origin %q in %v = %v, want %v", tc.origin, tc.in, got, tc.want)
		}
	}
	for site, want := range map[string]bool{"": true, "same-origin": true, "same-site": false, "cross-site": false, "none": false, "Same-Origin": false} {
		if got := sameOriginFetch(site); got != want {
			t.Errorf("Sec-Fetch-Site %q = %v, want %v", site, got, want)
		}
	}
	for ct, want := range map[string]bool{"application/json": true, "application/json; charset=utf-8": true, "Application/JSON": true,
		"text/plain": false, "": false, "application/json-patch+json": false, "multipart/form-data; boundary=x": false,
		"application/x-www-form-urlencoded": false, "application/json;;": false} {
		if got := jsonMediaType(ct); got != want {
			t.Errorf("Content-Type %q = %v, want %v", ct, got, want)
		}
	}
	for _, good := range []string{"https://s.tail0.ts.net", "https://s.tail0.ts.net:8443", "http://localhost:7644", "https://localhost"} {
		if err := checkSurfaceOrigin(good); err != nil {
			t.Errorf("--surface-origin %q refused: %v", good, err)
		}
	}
	for _, bad := range []string{"https://100.64.0.1", "https://x.ts.net:443", "https://example.com", "http://s.tail0.ts.net", "https://s.tail0.ts.net/", "s.tail0.ts.net", ""} {
		if err := checkSurfaceOrigin(bad); err == nil || !strings.Contains(err.Error(), "--surface-origin") {
			t.Errorf("--surface-origin %q: %v, want a refusal naming the flag", bad, err)
		}
		_, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, SurfaceOrigins: []string{"http://localhost:7644", bad}})
		var se *StartupError
		if !errors.As(err, &se) || se.Stage != StageConfig || !strings.Contains(se.Detail, "--surface-origin") {
			t.Errorf("New with --surface-origin %q: %v, want a config refusal naming the flag", bad, err)
		}
	}
	for _, race := range []error{store.ErrRequestConsumed, store.ErrRequestExpired, store.ErrRequestClockBehind, store.ErrRequestDigestChanged, store.ErrRequestUnknown} {
		if !isRequestRace(fmt.Errorf("%w: mr-0123456789abcdef", race)) {
			t.Errorf("%v is not a race (R14 409)", race)
		}
	}
	if isRequestRace(store.ErrBadSurfaceApproval) || isRequestRace(errors.New("disk I/O error")) {
		t.Error("a non-race record failure is a race")
	}
	if r := misdirectedReason("S.tail0.ts.net"); !strings.Contains(r, "--surface-origin https://s.tail0.ts.net") {
		t.Errorf("an unlisted ts.net Host is refused with %q, want the fix named", r)
	}
	if r := misdirectedReason("rebind.example:7644"); strings.Contains(r, "--surface-origin") {
		t.Errorf("a foreign Host is refused with %q, which names --surface-origin", r)
	}
	if r := crossOriginReason("https://s.tail0.ts.net"); !strings.Contains(r, "--surface-origin https://s.tail0.ts.net") {
		t.Errorf("an unlisted ts.net Origin is refused with %q, want the fix named", r)
	}
	for _, other := range []string{"null", "http://localhost:7645", "https://100.64.0.1"} {
		if r := crossOriginReason(other); strings.Contains(r, "--surface-origin") {
			t.Errorf("Origin %q is refused with %q, which names --surface-origin", other, r)
		}
	}
}

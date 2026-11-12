package main

// Row 201 A3 sprint 2 M4c (design §2.3 [r1], §6 CF1): serve --surface-origin.

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

// CF1: --surface-origin is in serve's help, repeatable, wired into Config.SurfaceOrigins (nil when
// unset), and a value outside A2's ceremony-origin grammar ends serve with exitUsage naming the
// flag, before the daemon is started.
func TestServeSurfaceOriginFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"serve", "--help"}, &stdout, &stderr); code != 0 ||
		!regexp.MustCompile(`(?m)^\s+--surface-origin <origin>$`).MatchString(stdout.String()+stderr.String()) {
		t.Fatalf("serve --help: exit %d, lacks the --surface-origin flag", code)
	}
	got := captureServeConfig(t)
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"serve", "--db", "/x/world.db", "--surface-origin", "https://s.tail0.ts.net", "--surface-origin", "http://localhost:7644"}, &stdout, &stderr); code != 0 {
		t.Fatalf("serve with two origins: exit %d, stderr %q", code, stderr.String())
	}
	if want := []string{"https://s.tail0.ts.net", "http://localhost:7644"}; !reflect.DeepEqual(got.SurfaceOrigins, want) {
		t.Fatalf("SurfaceOrigins = %#v, want %#v", got.SurfaceOrigins, want)
	}
	*got = daemon.Config{}
	if code := run([]string{"serve", "--db", "/x/world.db"}, &stdout, &stderr); code != 0 || got.DBPath != "/x/world.db" || got.SurfaceOrigins != nil {
		t.Fatalf("unset: exit %d, DBPath %q, SurfaceOrigins %#v; want nil", code, got.DBPath, got.SurfaceOrigins)
	}
	for _, bad := range []string{"https://100.64.0.1", "https://x.ts.net:443", "https://example.com", "http://s.tail0.ts.net", ""} {
		*got = daemon.Config{}
		stdout.Reset()
		stderr.Reset()
		if code := run([]string{"serve", "--db", "/x/world.db", "--surface-origin", "https://s.tail0.ts.net", "--surface-origin", bad}, &stdout, &stderr); code != exitUsage ||
			!strings.Contains(stderr.String(), `invalid value "`+bad+`" for flag -surface-origin: approval: `) || got.DBPath != "" {
			t.Errorf("--surface-origin %q: exit %d, daemon started %v, stderr %q; want exitUsage naming the flag before any start", bad, code, got.DBPath != "", stderr.String())
		}
	}
}

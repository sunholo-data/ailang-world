package surfacepage

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

func TestSurfaceModelServiceLine(t *testing.T) {
	at := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	n := Now{At: at, StartedAt: at.Add(-6 * time.Minute), Build: "sha256:abcdef012345", Supervisor: "com.sunholo.ailang-world.worldd"}
	for _, tc := range []struct {
		supervisor, installed, line string
		warn, pending               bool
	}{
		{n.Supervisor, "", "build abcdef0 · started 6 min ago", false, false},
		{"none", "", "World is not running as a service. It stops when its terminal does. Run ailang-worldd service install.", true, false},
		{"none", "sha256:123456789abc", "World is running build abcdef0 (started 6 min ago); the installed binary is 1234567. Run ailang-worldd service upgrade.", true, true},
	} {
		n.Supervisor, n.InstalledSHA256 = tc.supervisor, tc.installed
		s := BuildModel(surface.Index{}, n).NowEdge.Service
		if s == nil || s.Build != "abcdef0" || s.StartedAt != "2026-10-30T11:54:00Z" || s.Line != tc.line || s.Warn != tc.warn || s.RestartPending != tc.pending || s.Supervisor != tc.supervisor || s.InstalledSHA256 != tc.installed {
			t.Fatalf("service: %+v, want %+v", s, tc)
		}
	}
	n.Build = ""
	raw, err := json.Marshal(BuildModel(surface.Index{}, n))
	if err != nil || strings.Contains(string(raw), "\"service\"") {
		t.Fatalf("empty build must omit service: %s (%v)", raw, err)
	}
}

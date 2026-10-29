package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

func TestDoctorExecToolsWithAnM3Daemon(t *testing.T) {
	for _, probe := range []string{"pending", "ok", "failed: arm7 toolchain", "does-not-start: --version failed"} {
		t.Run(probe, func(t *testing.T) {
			var out bytes.Buffer
			next := "brew reinstall node, then retry (no restart needed)"
			r := &doctorReport{w: &out, execToolStatus: map[string]broker.ExecToolStatus{
				"node": {Configured: "/opt/homebrew/bin/node", Target: "/opt/homebrew/Cellar/node/26.12.1/bin/node",
					Accepted: "/opt/homebrew/Cellar/node/26.0.0/bin/node", Probe: probe, Next: next},
			}}
			doctorExecTools(r, "", map[string]string{"node": "/opt/homebrew/Cellar/node/26.0.0/bin/node"})
			want := "✓ exec tools"
			if strings.HasPrefix(probe, "failed:") || strings.HasPrefix(probe, "does-not-start:") {
				want = "✗ exec tools"
				if !strings.Contains(out.String(), next) || r.fails != 1 {
					t.Fatalf("missing F15 next: %s", out.String())
				}
			}
			if !strings.Contains(out.String(), want) || strings.Contains(out.String(), "restart World") {
				t.Fatal(out.String())
			}
			if probe == "pending" && (!strings.Contains(out.String(), "resolves to /opt/homebrew/Cellar/node/26.12.1/bin/node") ||
				!strings.Contains(out.String(), "the next exec call re-checks it, no restart needed")) {
				t.Fatal(out.String())
			}
		})
	}
}

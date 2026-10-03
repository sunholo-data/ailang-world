package daemon

// Row 135 M3 AC3.1 (design §4.3 gate 2, §4.6): the daemon's startup refusal
// table for the --run-* allowlist, and the positive case on the real tool
// binary, whose own policy-tool summary verifies every variant at startup.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

func TestRunCapsStartupRefusals(t *testing.T) {
	for _, tc := range []struct {
		name   string
		edit   func(f wsFixture, cfg *Config)
		reason string
	}{
		{"unknown name", func(_ wsFixture, c *Config) { c.RunAllowCaps = []string{"Process"} }, `names "Process"`},
		{"Net without --run-net-allow", func(_ wsFixture, c *Config) { c.RunAllowCaps = []string{"Net"} }, "needs at least one --run-net-allow"},
		{"--run-net-allow without Net", func(_ wsFixture, c *Config) { c.RunNetAllow = []string{"127.0.0.1:7655"} }, "does not name Net"},
		{"bare host", func(_ wsFixture, c *Config) {
			c.RunAllowCaps, c.RunNetAllow = []string{"Net"}, []string{"127.0.0.1"}
		}, "is not HOST:PORT"},
		{"named host", func(_ wsFixture, c *Config) {
			c.RunAllowCaps, c.RunNetAllow = []string{"Net"}, []string{"localhost:7655"}
		}, "names a host by name"},
		{"private address", func(_ wsFixture, c *Config) {
			c.RunAllowCaps, c.RunNetAllow = []string{"Net"}, []string{"10.0.0.1:80"}
		}, "not a loopback address"},
		{"link-local address", func(_ wsFixture, c *Config) {
			c.RunAllowCaps, c.RunNetAllow = []string{"Net"}, []string{"169.254.169.254:80"}
		}, "not a loopback address"},
		{"without the workspace tools", func(_ wsFixture, c *Config) {
			c.WorkspaceRoot, c.ToolAilangBin, c.RunAllowCaps = "", "", []string{"Declassify"}
		}, "need the workspace tools"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newWSFixture(t)
			cfg := Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)}
			tc.edit(f, &cfg)
			_, err := newWSDaemon(t, cfg)
			var startup *StartupError
			if !errors.As(err, &startup) || !strings.Contains(err.Error(), "the run capabilities are refused") ||
				!strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("New = %v, want a StartupError refusing the run capabilities with %q", err, tc.reason)
			}
		})
	}
}

// TestRunCapsStartupVerifiesVariantsOnTheToolBinary: a valid allowlist starts
// on the real binary (each variant passes policy-tool summary) and leaves no
// check files behind; the handler of an episode carries it.
func TestRunCapsStartupVerifiesVariantsOnTheToolBinary(t *testing.T) {
	tool := os.Getenv(toolBinEnv)
	if tool == "" {
		t.Skipf("%s unset: the startup variant check needs the %s tool binary; skipping", toolBinEnv, ToolBinaryRelease)
	}
	f := newWSFixture(t)
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: tool,
		RunAllowCaps: []string{"Declassify", "Env", "Net"}, RunNetAllow: []string{"127.0.0.1:7655", "[::1]:7655"}, RunNetAllowHTTP: true})
	if got := d.workspace.runCaps; len(got.Allow) != 3 || len(got.NetAllow) != 2 || !got.NetAllowHTTP {
		t.Fatalf("workspace run caps = %+v", got)
	}
	if _, err := os.Stat(filepath.Join(f.stateDir, "policies", ".run-caps-check")); !os.IsNotExist(err) {
		t.Fatalf("the startup check left its policies behind (%v)", err)
	}
	if reg := d.workspace.registry("ep1"); reg[broker.EffectAilangRunNet] == nil || reg[broker.EffectAilangRunEnv] == nil {
		t.Fatalf("registry(ep1) lacks the run effects: %v", registryNames(reg))
	}
}

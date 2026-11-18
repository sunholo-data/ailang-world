package main

import (
	"regexp"
	"strconv"
	"testing"
	"time"
)

// serviceDrain is the daemon's shutdown drain that service upgrade announces
// ("expected downtime: 10 s drain plus a restart").
const serviceDrain = 10 * time.Second

// TestLaunchctlBudgetCoversTheStop pins the 2026-11-18 failure: the first real
// `service upgrade` ran `launchctl kickstart -k` under a flat 2 s budget, which
// killed launchctl while launchd was still waiting for the old daemon to drain.
// The upgrade then rolled back, the rollback's kickstart was killed the same way,
// and the service was left stopped. The verbs that stop the daemon must outlast
// the plist's ExitTimeOut (launchd's own ceiling on that wait); the rest stay short.
func TestLaunchctlBudgetCoversTheStop(t *testing.T) {
	plist := renderServicePlist(worldServiceLabel, servicePaths{Log: "/x/worldd.log"}, []string{"/x/ailang-worldd", "serve"}, "/x")
	m := regexp.MustCompile(`<key>ExitTimeOut</key><integer>(\d+)</integer>`).FindStringSubmatch(plist)
	if m == nil {
		t.Fatal("rendered plist carries no ExitTimeOut")
	}
	secs, _ := strconv.Atoi(m[1])
	if got := time.Duration(secs) * time.Second; got != serviceExitTimeout {
		t.Fatalf("plist ExitTimeOut %v, serviceExitTimeout %v: keep them equal", got, serviceExitTimeout)
	}
	for _, verb := range []string{"kickstart", "bootout"} {
		if b := launchctlBudget(verb); b <= serviceExitTimeout || b <= serviceDrain {
			t.Errorf("launchctl %s budget %v does not outlast ExitTimeOut %v (drain %v)", verb, b, serviceExitTimeout, serviceDrain)
		}
	}
	for _, verb := range []string{"print", "bootstrap", "getenv"} {
		if b := launchctlBudget(verb); b != 2*time.Second {
			t.Errorf("launchctl %s budget %v, want 2s", verb, b)
		}
	}
}

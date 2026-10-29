package main

import (
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/toolprobe"
)

// serviceController is the only process boundary for the service supervisor.
var errServiceNotLoaded = errors.New("service not loaded")

const serviceRegistryKeyFix = "launchctl unsetenv AILANG_REGISTRY_API_KEY"

type serviceController interface {
	Bootstrap(domain, plist string) error
	Kickstart(domain, label string) error
	Bootout(domain, label string) error
	Print(domain, label string) (string, error)
	Getenv(name string) (string, error)
}
type systemServiceController struct{}

// command runs launchctl through toolprobe, the AC10-driven subprocess site that
// scrubs the registry credential (a direct exec here failed the AC10 census).
// launchctl reports "Could not find service" on stderr, which a non-zero exit
// carries in its *exec.ExitError; it is joined to stdout as CombinedOutput did.
func (systemServiceController) command(args ...string) (string, error) {
	ctx, cancel := budgetContext(2 * time.Second)
	defer cancel()
	out, err := toolprobe.OutputEnv(ctx, "/bin/launchctl", "", []string{"PATH=/usr/bin:/bin:/usr/sbin:/sbin"}, args...)
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		out = append(out, exit.Stderr...)
	}
	return strings.TrimSpace(string(out)), err
}
func (c systemServiceController) Bootstrap(domain, plist string) error {
	_, err := c.command("bootstrap", domain, plist)
	return err
}
func (c systemServiceController) Kickstart(domain, label string) error {
	_, err := c.command("kickstart", "-k", domain+"/"+label)
	return err
}
func (c systemServiceController) Bootout(domain, label string) error {
	_, err := c.command("bootout", domain+"/"+label)
	return err
}
func (c systemServiceController) Print(domain, label string) (string, error) {
	out, err := c.command("print", domain+"/"+label)
	if err != nil && strings.Contains(out, "Could not find service") {
		return out, errServiceNotLoaded
	}
	return out, err
}
func (c systemServiceController) Getenv(name string) (string, error) {
	return c.command("getenv", name)
}

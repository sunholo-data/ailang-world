package broker

import (
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/childenv"
	"github.com/sunholo-data/ailang-world/host/toolprobe"
)

// AC10 / TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential:
// doctor and launcher tests share toolprobe's existing driven exec site.
// This drives its explicit-environment arm with a deliberately contaminated env.
func TestDoctorVersionUsesDrivenScrubbedSite(t *testing.T) {
	dump := filepath.Join(t.TempDir(), "version.env")
	probe := writeProbeScript(t, dump)
	_, err := toolprobe.OutputEnv(boundedTestContext(t), probe, t.TempDir(),
		[]string{"PATH=/usr/bin:/bin", "HOME=" + t.TempDir(), RegistryCredentialVariable + "=" + ac10Sentinel},
		"--version")
	if err != nil {
		t.Fatal(err)
	}
	env := readEnvDump(t, dump)
	if len(env) == 0 {
		t.Fatal("version probe observed an empty environment")
	}
	if childenv.Has(env, RegistryCredentialVariable) {
		t.Fatal("version probe inherited the registry credential")
	}
}

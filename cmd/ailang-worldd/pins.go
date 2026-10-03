package main

// Row 138 M3 (design_docs/planned/w-worldd-developer-cli.md §3.5): the pin
// table `setup` installs from and `doctor` checks against. Every digest here
// was MEASURED (2026-10-03, row 138 M3 executor): each release asset was
// downloaded from github.com/sunholo-data/ailang, its tarball hashed, its one
// `ailang` entry extracted and hashed. The release's own `.sha256` agreed with
// the computed tarball digest in all four cases. They are compiled in so the
// release cannot vouch for itself (R-CLI-3).
//
// Two pins, two releases (D-WORLD-57):
//   - the INTERPRETER (`serve --ailang-bin`), test-bound to world-publish's
//     frozenCompilerVersion, ci.yml and scripts/verify_world_package.sh;
//   - the TOOL binary (`serve --tool-ailang-bin`), whose release is
//     daemon.ToolBinaryRelease — one source of truth, so a tool bump that
//     forgets this table fails TestPinTableCoversBothReleases (R-CLI-5).
//
// Platforms: darwin/arm64 and linux/amd64 only (D-CLI-2).

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/pinfetch"
)

// interpreterRelease is the pinned interpreter release.
const interpreterRelease = "v0.41.0"

// toolRelease is daemon.ToolBinaryRelease without its "AILANG " prefix.
func toolRelease() string { return strings.TrimPrefix(daemon.ToolBinaryRelease, "AILANG ") }

// supportedPlatforms is D-CLI-2's list.
var supportedPlatforms = []string{"darwin/arm64", "linux/amd64"}

// pinTable is a var only so the setup/doctor tests can serve small fixture
// releases; production never assigns it.
var pinTable = []pinfetch.Pin{
	{Release: "v0.41.0", Platform: "darwin/arm64", Asset: "darwin.arm64.ailang.tar.gz",
		TarballSHA256: "b08f3cdea137222ea8df71c53d1a1353ef7c942fd2a307de5d7a17b8fb598e0b",
		BinarySHA256:  "1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f"},
	{Release: "v0.41.0", Platform: "linux/amd64", Asset: "linux.x64.ailang.tar.gz",
		TarballSHA256: "fa0045dee577942f91f045b2e52d324d34db6087d31802bdabeb4b30a26faa56",
		BinarySHA256:  "8e7a275da6f26c38127352518ce7cdf855375dadff49547216ed90790ed25fb5"},
	{Release: "v0.52.1", Platform: "darwin/arm64", Asset: "darwin.arm64.ailang.tar.gz",
		TarballSHA256: "576236fe152e8f6ee614529d2d12d7212b16743b0443e9774fa5f89c8fa7579e",
		BinarySHA256:  "0dd70a1d00360be0b7b4c667054c71fdc4ea87aa3deb79710cf56f18ffe1a8f5"},
	{Release: "v0.52.1", Platform: "linux/amd64", Asset: "linux.x64.ailang.tar.gz",
		TarballSHA256: "c682c30f629ca2305f5b09c49ee2b9d43bf73beb68b095b51e19fc31a46f883f",
		BinarySHA256:  "97dcd4a51bd4197ef22be05fa0bff7b24591b1c7a2a6a2250413abfbc3070f30"},
}

// lookupPin returns the table row for release on platform.
func lookupPin(release, platform string) (pinfetch.Pin, error) {
	supported := false
	for _, p := range supportedPlatforms {
		if p == platform {
			supported = true
		}
	}
	if !supported {
		return pinfetch.Pin{}, fmt.Errorf("platform %s is not supported: AILANG World pins only %s (D-CLI-2)",
			platform, strings.Join(supportedPlatforms, " and "))
	}
	for _, p := range pinTable {
		if p.Release == release && p.Platform == platform {
			return p, nil
		}
	}
	return pinfetch.Pin{}, fmt.Errorf("no pin-table row for %s on %s (cmd/ailang-worldd/pins.go)", release, platform)
}

// pinRole is one of the two pins and where it is installed.
type pinRole struct {
	name    string // "interpreter" or "tool"
	release string
	dest    string // the binary path
}

// pinRoles resolves both pins' destinations under the given directories.
func pinRoles(interpDir, toolsDir string) []pinRole {
	return []pinRole{
		{name: "interpreter", release: interpreterRelease, dest: filepath.Join(interpDir, "ailang")},
		{name: "tool", release: toolRelease(), dest: filepath.Join(toolsDir, toolRelease(), "ailang")},
	}
}

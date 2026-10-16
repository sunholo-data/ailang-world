package main

// Row 138 M3 (design_docs/planned/w-worldd-developer-cli.md §3.5): the pin
// table `setup` installs from and `doctor` checks against. Since row 139 M1 it
// is DERIVED from the release manifest (pins/ailang-releases.json, embedded by
// host/pins), which carries the digests measured 2026-10-03 by the row 138 M3
// executor: each release asset downloaded from github.com/sunholo-data/ailang,
// its tarball hashed, its one `ailang` entry extracted and hashed, the
// release's own `.sha256` agreeing in all four cases. They are compiled in so
// the release cannot vouch for itself (R-CLI-3).
//
// Two pins, two releases (D-WORLD-57):
//   - the INTERPRETER (`serve --ailang-bin`): the manifest's `interpreter` role;
//   - the TOOL binary (`serve --tool-ailang-bin`), whose release is
//     daemon.ToolBinaryRelease, itself the manifest's `tool` role.
//
// Platforms: darwin/arm64 and linux/amd64 only (D-CLI-2).

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/pinfetch"
	"github.com/sunholo-data/ailang-world/host/pins"
)

// interpreterRelease is the pinned interpreter release (manifest role
// `interpreter`).
var interpreterRelease = pins.ReleaseFor(pins.RoleInterpreter).Release

// toolRelease is daemon.ToolBinaryRelease without its "AILANG " prefix.
func toolRelease() string { return strings.TrimPrefix(daemon.ToolBinaryRelease, "AILANG ") }

// supportedPlatforms is D-CLI-2's list, as the manifest states it.
var supportedPlatforms = pins.Current().Platforms

// pinTable is a var only so the setup/doctor tests can serve small fixture
// releases; production never assigns it.
var pinTable = manifestPinTable(pins.Current(), interpreterRelease, toolRelease())

// manifestPinTable builds the setup/doctor table for the given releases: one
// row per manifest platform, in release-then-platform order.
func manifestPinTable(m pins.Manifest, releases ...string) []pinfetch.Pin {
	var out []pinfetch.Pin
	for _, rel := range releases {
		r, ok := m.Lookup(rel)
		if !ok {
			panic(fmt.Sprintf("pins: release %s is not in the manifest", rel))
		}
		for _, plat := range m.Platforms {
			a := r.Assets[plat]
			out = append(out, pinfetch.Pin{Release: rel, Platform: plat, Asset: a.Asset,
				TarballSHA256: a.TarballSHA256, BinarySHA256: a.BinarySHA256})
		}
	}
	return out
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
	return pinfetch.Pin{}, fmt.Errorf("no pin-table row for %s on %s (pins/ailang-releases.json)", release, platform)
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

// Package pins is the AILANG release manifest (row 139 D1,
// design_docs/planned/w-interpreter-epoch-upgrade-path.md): which AILANG
// releases World uses, in which role, with which measured digests.
//
// The canonical file is pins/ailang-releases.json at the repo root. go:embed
// cannot reach upward out of a package directory, so this package embeds a
// byte-identical copy (ailang-releases.json beside this file);
// TestEmbeddedManifestMatchesCanonical reds on any drift between the two.
//
// Every production release literal derives from here or is checked against it
// by TestReleaseLiteralsExplainedByManifest (AC1). A release bump is one
// manifest edit (both copies) plus whatever that test then names.
//
// The package is pure: it parses compiled-in bytes and imports nothing from
// the host, so both host/daemon and the cmd binaries can depend on it.
package pins

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"sync"
)

// Schema is the manifest's schema id.
const Schema = "world/ailang-releases/v1"

// Roles. The exclusive roles are each held by exactly one release; observed
// (the releases verify_ail.sh legs 1-2 expect to resolve) may be held by many.
const (
	RoleInterpreter    = "interpreter"     // serve --ailang-bin; setup/doctor
	RoleProof          = "proof"           // demo proof checker / seed
	RoleFrozenCompiler = "frozen-compiler" // world-publish ready packet
	RoleGatePrimary    = "gate-primary"    // AILANG_BIN of the verify gates and CI
	RoleTool           = "tool"            // serve --tool-ailang-bin
	RoleObserved       = "observed"        // scripts/testdata/ailang_release_observed.txt
)

// exclusiveRoles lists the roles exactly one release must hold.
var exclusiveRoles = []string{RoleInterpreter, RoleProof, RoleFrozenCompiler, RoleGatePrimary, RoleTool}

// ExclusiveRoles returns a copy of the single-holder role list.
func ExclusiveRoles() []string { return append([]string(nil), exclusiveRoles...) }

//go:embed ailang-releases.json
var embedded []byte

// Asset is one platform's release tarball and the `ailang` binary inside it.
type Asset struct {
	Asset         string `json:"asset"`
	TarballSHA256 string `json:"tarball_sha256"`
	BinarySHA256  string `json:"binary_sha256"`
}

// Role is one qualification a release holds, with when and why.
type Role struct {
	Role     string `json:"role"`
	Since    string `json:"since"`
	Evidence string `json:"evidence"`
}

// Release is one AILANG release.
type Release struct {
	Release string           `json:"release"` // "v0.41.0"
	Version string           `json:"version"` // first --version line, "AILANG v0.41.0"
	Commit  string           `json:"commit,omitempty"`
	Roles   []Role           `json:"roles"`
	Assets  map[string]Asset `json:"assets,omitempty"`
}

// Has reports whether r holds role.
func (r Release) Has(role string) bool {
	for _, x := range r.Roles {
		if x.Role == role {
			return true
		}
	}
	return false
}

// Manifest is the parsed release manifest.
type Manifest struct {
	Schema    string    `json:"schema"`
	Note      string    `json:"note,omitempty"`
	Platforms []string  `json:"platforms"`
	Releases  []Release `json:"releases"`
}

var (
	releaseRe = regexp.MustCompile(`^v0\.[0-9]+\.[0-9]+$`)
	hex64     = regexp.MustCompile(`^[0-9a-f]{64}$`)
	dateRe    = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`)
	knownRole = map[string]bool{RoleInterpreter: true, RoleProof: true, RoleFrozenCompiler: true,
		RoleGatePrimary: true, RoleTool: true, RoleObserved: true}
	// assetName is D-CLI-2's per-platform upstream asset spelling.
	assetName = map[string]string{"darwin/arm64": "darwin.arm64.ailang.tar.gz", "linux/amd64": "linux.x64.ailang.tar.gz"}
)

// Parse decodes and validates a manifest. It refuses unknown fields, a wrong
// schema, a malformed release or digest, an exclusive role held by zero or
// several releases, and a binary-carrying release missing a platform.
func Parse(b []byte) (Manifest, error) {
	var m Manifest
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		return Manifest{}, fmt.Errorf("pins: decode: %w", err)
	}
	if m.Schema != Schema {
		return Manifest{}, fmt.Errorf("pins: schema %q, want %q", m.Schema, Schema)
	}
	if fmt.Sprint(m.Platforms) != "[darwin/arm64 linux/amd64]" {
		return Manifest{}, fmt.Errorf("pins: platforms %v, want [darwin/arm64 linux/amd64] (D-CLI-2)", m.Platforms)
	}
	seen := map[string]bool{}
	holders := map[string]int{}
	for _, r := range m.Releases {
		if !releaseRe.MatchString(r.Release) {
			return Manifest{}, fmt.Errorf("pins: malformed release %q", r.Release)
		}
		if seen[r.Release] {
			return Manifest{}, fmt.Errorf("pins: release %s listed twice", r.Release)
		}
		seen[r.Release] = true
		if r.Version != "AILANG "+r.Release {
			return Manifest{}, fmt.Errorf("pins: %s version line %q, want %q", r.Release, r.Version, "AILANG "+r.Release)
		}
		if len(r.Roles) == 0 {
			return Manifest{}, fmt.Errorf("pins: %s holds no role", r.Release)
		}
		binaryRole := false
		roleSeen := map[string]bool{}
		for _, role := range r.Roles {
			if !knownRole[role.Role] {
				return Manifest{}, fmt.Errorf("pins: %s: unknown role %q", r.Release, role.Role)
			}
			if roleSeen[role.Role] {
				return Manifest{}, fmt.Errorf("pins: %s: role %s twice", r.Release, role.Role)
			}
			roleSeen[role.Role] = true
			if !dateRe.MatchString(role.Since) || role.Evidence == "" {
				return Manifest{}, fmt.Errorf("pins: %s role %s needs a since date and evidence", r.Release, role.Role)
			}
			holders[role.Role]++
			if role.Role != RoleObserved {
				binaryRole = true
			}
		}
		if binaryRole && len(r.Assets) != len(m.Platforms) {
			return Manifest{}, fmt.Errorf("pins: %s holds a binary role but has %d platform assets, want %d", r.Release, len(r.Assets), len(m.Platforms))
		}
		for plat, a := range r.Assets {
			want, ok := assetName[plat]
			if !ok {
				return Manifest{}, fmt.Errorf("pins: %s: unsupported platform %s", r.Release, plat)
			}
			if a.Asset != want {
				return Manifest{}, fmt.Errorf("pins: %s %s: asset %q, want %q", r.Release, plat, a.Asset, want)
			}
			if !hex64.MatchString(a.TarballSHA256) || !hex64.MatchString(a.BinarySHA256) {
				return Manifest{}, fmt.Errorf("pins: %s %s: malformed digests", r.Release, plat)
			}
		}
	}
	for _, role := range exclusiveRoles {
		if holders[role] != 1 {
			return Manifest{}, fmt.Errorf("pins: role %s held by %d releases, want exactly 1", role, holders[role])
		}
	}
	if holders[RoleObserved] == 0 {
		return Manifest{}, fmt.Errorf("pins: no release holds role %s", RoleObserved)
	}
	return m, nil
}

// Holders returns the releases holding role, in manifest order.
func (m Manifest) Holders(role string) []Release {
	var out []Release
	for _, r := range m.Releases {
		if r.Has(role) {
			out = append(out, r)
		}
	}
	return out
}

// Holder returns the one release holding an exclusive role.
func (m Manifest) Holder(role string) (Release, error) {
	h := m.Holders(role)
	if len(h) != 1 {
		return Release{}, fmt.Errorf("pins: role %s held by %d releases", role, len(h))
	}
	return h[0], nil
}

// Lookup returns the release named rel ("v0.41.0").
func (m Manifest) Lookup(rel string) (Release, bool) {
	for _, r := range m.Releases {
		if r.Release == rel {
			return r, true
		}
	}
	return Release{}, false
}

// Names returns every release name, sorted.
func (m Manifest) Names() []string {
	var out []string
	for _, r := range m.Releases {
		out = append(out, r.Release)
	}
	sort.Strings(out)
	return out
}

var (
	once    sync.Once
	current Manifest
	loadErr error
)

// Embedded returns a copy of the compiled-in manifest bytes.
func Embedded() []byte { return append([]byte(nil), embedded...) }

// Current returns the compiled-in manifest. An invalid compiled-in manifest is
// a build defect (its tests red first), so it panics rather than letting a
// binary run with no pins.
func Current() Manifest {
	once.Do(func() { current, loadErr = Parse(embedded) })
	if loadErr != nil {
		panic(loadErr)
	}
	return current
}

// ReleaseFor returns the release holding an exclusive role in the compiled-in
// manifest.
func ReleaseFor(role string) Release {
	r, err := Current().Holder(role)
	if err != nil {
		panic(err)
	}
	return r
}

// VersionLine is ReleaseFor(role).Version, e.g. "AILANG v0.41.0".
func VersionLine(role string) string { return ReleaseFor(role).Version }

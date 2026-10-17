package pins

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The canonical manifest is pins/ailang-releases.json; this package embeds a
// copy because go:embed cannot reach the repo root. Drift reds here.
func TestEmbeddedManifestMatchesCanonical(t *testing.T) {
	canonical, err := os.ReadFile(filepath.Join(repoRoot(t), "pins", "ailang-releases.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(canonical, Embedded()) {
		t.Fatal("host/pins/ailang-releases.json differs from pins/ailang-releases.json: copy the canonical file over the embedded one (cp pins/ailang-releases.json host/pins/)")
	}
}

// Today's pins, restated independently of the manifest: the four v0.41.0
// roles of row 139 M1, the tool role on v0.54.15 since D7's tool half, and
// v0.52.1 (the tool from row 135 M0) kept as retired history with its assets.
func TestCurrentManifestIsTodaysPins(t *testing.T) {
	m := Current()
	for role, want := range map[string]string{RoleInterpreter: "v0.41.0", RoleProof: "v0.41.0",
		RoleFrozenCompiler: "v0.41.0", RoleGatePrimary: "v0.41.0", RoleTool: "v0.54.15"} {
		if got := ReleaseFor(role).Release; got != want {
			t.Errorf("role %s held by %s, want %s", role, got, want)
		}
	}
	var observed, retired []string
	for _, r := range m.Holders(RoleObserved) {
		observed = append(observed, r.Release)
	}
	for _, r := range m.Holders(RoleRetired) {
		retired = append(retired, r.Release)
	}
	if strings.Join(observed, " ") != "v0.41.0 v0.33.0 v0.30.0" {
		t.Errorf("observed releases %v", observed)
	}
	if strings.Join(retired, " ") != "v0.52.1" {
		t.Errorf("retired releases %v, want [v0.52.1] (D7: the old tool stays in the manifest for history)", retired)
	}
	if got := VersionLine(RoleTool); got != "AILANG v0.54.15" {
		t.Errorf("tool version line %q", got)
	}
	if r := ReleaseFor(RoleTool); r.Commit != "0384450" {
		t.Errorf("tool commit %q, want 0384450", r.Commit)
	}
	// The pre-M1 digests (cmd/ailang-worldd/pins.go at c9ff55e), restated, and
	// the v0.54.15 digests the operator measured from the published release
	// and checked against its .sha256 files (row 139 D7, 2026-10-17).
	want := map[string][2]string{
		"v0.41.0 darwin/arm64":  {"b08f3cdea137222ea8df71c53d1a1353ef7c942fd2a307de5d7a17b8fb598e0b", "1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f"},
		"v0.41.0 linux/amd64":   {"fa0045dee577942f91f045b2e52d324d34db6087d31802bdabeb4b30a26faa56", "8e7a275da6f26c38127352518ce7cdf855375dadff49547216ed90790ed25fb5"},
		"v0.52.1 darwin/arm64":  {"576236fe152e8f6ee614529d2d12d7212b16743b0443e9774fa5f89c8fa7579e", "0dd70a1d00360be0b7b4c667054c71fdc4ea87aa3deb79710cf56f18ffe1a8f5"},
		"v0.52.1 linux/amd64":   {"c682c30f629ca2305f5b09c49ee2b9d43bf73beb68b095b51e19fc31a46f883f", "97dcd4a51bd4197ef22be05fa0bff7b24591b1c7a2a6a2250413abfbc3070f30"},
		"v0.54.15 darwin/arm64": {"710945da1baee47fea61bb2eeb9f2018a90c491fa2c6591d5f8dd7ce4e281382", "d821dad0d6e906604182a8b5bb29d60e57c872395ba11ad9e8a7f34e440ac043"},
		"v0.54.15 linux/amd64":  {"877316981bdbf108dacacafc733720333e718a912e0d215564baa37c4cac4945", "db0f9dac038fbc99212c2d3effc30d1e8a98496c6a481384c1c21a8091ddd63c"},
	}
	n := 0
	for _, r := range m.Releases {
		for plat, a := range r.Assets {
			n++
			w, ok := want[r.Release+" "+plat]
			if !ok || a.TarballSHA256 != w[0] || a.BinarySHA256 != w[1] {
				t.Errorf("%s %s digests %+v, want %v", r.Release, plat, a, w)
			}
		}
	}
	if n != len(want) {
		t.Errorf("manifest carries %d platform assets, want %d", n, len(want))
	}
}

func TestParseRefusesMalformedManifests(t *testing.T) {
	if _, err := Parse(Embedded()); err != nil {
		t.Fatalf("control: embedded manifest refused: %v", err)
	}
	mutate := func(f func(map[string]any)) []byte {
		var v map[string]any
		if err := json.Unmarshal(Embedded(), &v); err != nil {
			t.Fatal(err)
		}
		f(v)
		b, _ := json.Marshal(v)
		return b
	}
	// Releases are found by name, not position, so a manifest that gains a
	// row does not silently retarget a case.
	rel := func(v map[string]any, name string) map[string]any {
		for _, r := range v["releases"].([]any) {
			if r := r.(map[string]any); r["release"] == name {
				return r
			}
		}
		t.Fatalf("no release %s in the embedded manifest", name)
		return nil
	}
	role := func(name, since, evidence string) map[string]any {
		return map[string]any{"role": name, "since": since, "evidence": evidence}
	}
	tool := ReleaseFor(RoleTool).Release
	cases := map[string][]byte{
		"wrong schema":     mutate(func(v map[string]any) { v["schema"] = "world/ailang-releases/v0" }),
		"unknown field":    mutate(func(v map[string]any) { v["extra"] = true }),
		"bad version line": mutate(func(v map[string]any) { rel(v, "v0.41.0")["version"] = "AILANG v0.41.1" }),
		"bad digest": mutate(func(v map[string]any) {
			rel(v, "v0.41.0")["assets"].(map[string]any)["linux/amd64"].(map[string]any)["binary_sha256"] = "00"
		}),
		"wrong asset name": mutate(func(v map[string]any) {
			rel(v, tool)["assets"].(map[string]any)["linux/amd64"].(map[string]any)["asset"] = "linux.amd64.tar.gz"
		}),
		"missing platform": mutate(func(v map[string]any) { delete(rel(v, tool)["assets"].(map[string]any), "darwin/arm64") }),
		"tool role twice": mutate(func(v map[string]any) {
			rel(v, "v0.41.0")["roles"] = append(rel(v, "v0.41.0")["roles"].([]any), role("tool", "2026-10-16", "x"))
		}),
		"release with no role": mutate(func(v map[string]any) { rel(v, tool)["roles"] = []any{} }),
		// The tool retires with no successor: an exclusive role held by none.
		"tool role absent":  mutate(func(v map[string]any) { rel(v, tool)["roles"] = []any{role("retired", "2026-10-17", "x")} }),
		"unknown role":      mutate(func(v map[string]any) { rel(v, "v0.33.0")["roles"].([]any)[0].(map[string]any)["role"] = "nightly" }),
		"role without date": mutate(func(v map[string]any) { rel(v, "v0.33.0")["roles"].([]any)[0].(map[string]any)["since"] = "" }),
		"duplicate release": mutate(func(v map[string]any) {
			rel(v, "v0.30.0")["version"] = "AILANG v0.33.0"
			rel(v, "v0.30.0")["release"] = "v0.33.0"
		}),
		"malformed release": mutate(func(v map[string]any) {
			rel(v, "v0.30.0")["version"] = "AILANG 0.30.0"
			rel(v, "v0.30.0")["release"] = "0.30.0"
		}),
		"third platform": mutate(func(v map[string]any) { v["platforms"] = []any{"darwin/arm64", "linux/amd64", "linux/arm64"} }),
	}
	for name, b := range cases {
		if _, err := Parse(b); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	// The retired role (row 139 D7). Each refusal must be for its own reason,
	// so the message is asserted, not just the error.
	retiredCases := map[string]struct {
		b    []byte
		want string
	}{
		"retired release without its assets": {mutate(func(v map[string]any) { delete(rel(v, "v0.52.1"), "assets") }),
			"v0.52.1 holds a binary role but has 0 platform assets"},
		"retired release still holding an exclusive role": {mutate(func(v map[string]any) {
			rel(v, "v0.52.1")["roles"] = []any{role("retired", "2026-10-17", "x"), role("tool", "2026-10-03", "x")}
			rel(v, tool)["roles"] = []any{role("retired", "2026-10-17", "x")}
		}), "v0.52.1 is retired but still holds exclusive role tool"},
		"retired without evidence": {mutate(func(v map[string]any) {
			rel(v, "v0.52.1")["roles"].([]any)[0].(map[string]any)["evidence"] = ""
		}), "v0.52.1 role retired needs a since date and evidence"},
		"retired twice on one release": {mutate(func(v map[string]any) {
			rel(v, "v0.52.1")["roles"] = append(rel(v, "v0.52.1")["roles"].([]any), role("retired", "2026-10-17", "y"))
		}), "v0.52.1: role retired twice"},
	}
	for name, tc := range retiredCases {
		_, err := Parse(tc.b)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: err = %v, want %q", name, err, tc.want)
		}
	}
	// Positive control: retired is NOT exclusive. Two retired releases, one of
	// them also observed, parse; Holders lists both in manifest order.
	several := mutate(func(v map[string]any) {
		r := rel(v, "v0.33.0")
		r["assets"] = rel(v, "v0.52.1")["assets"]
		r["roles"] = append(r["roles"].([]any), role("retired", "2026-10-17", "control"))
	})
	m, err := Parse(several)
	if err != nil {
		t.Fatalf("control: two retired releases refused: %v", err)
	}
	var names []string
	for _, r := range m.Holders(RoleRetired) {
		names = append(names, r.Release)
	}
	if strings.Join(names, " ") != "v0.52.1 v0.33.0" {
		t.Fatalf("control: retired holders %v, want [v0.52.1 v0.33.0]", names)
	}
}

// verify_world_package.sh pins the frozen compiler's BINARY bytes per
// platform; those must be the manifest's frozen-compiler digests.
func TestVerifyWorldPackageDigestsMatchManifest(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(repoRoot(t), "scripts", "verify_world_package.sh"))
	if err != nil {
		t.Fatal(err)
	}
	fc := ReleaseFor(RoleFrozenCompiler)
	re := regexp.MustCompile(`(?m)^\s*(Darwin/arm64|Linux/x86_64)\) printf '%s' '([0-9a-f]{64})' ;;$`)
	got := map[string]string{}
	for _, m := range re.FindAllStringSubmatch(string(src), -1) {
		got[m[1]] = m[2]
	}
	for uname, plat := range map[string]string{"Darwin/arm64": "darwin/arm64", "Linux/x86_64": "linux/amd64"} {
		if got[uname] != fc.Assets[plat].BinarySHA256 {
			t.Errorf("verify_world_package.sh %s pins %q, manifest frozen-compiler %s binary is %s", uname, got[uname], fc.Release, fc.Assets[plat].BinarySHA256)
		}
	}
	if len(got) != 2 {
		t.Errorf("verify_world_package.sh platform table has %d rows, want 2", len(got))
	}
}

// ci.yml installs exactly the frozen-compiler, gate-primary and tool releases,
// each from the manifest's linux/amd64 asset.
func TestCIInstallsMatchManifest(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(repoRoot(t), ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	m := Current()
	base := regexp.MustCompile(`base="https://github.com/sunholo-data/ailang/releases/download/(v0\.[0-9]+\.[0-9]+)"`)
	asset := regexp.MustCompile(`"\$base/([a-z0-9.]+\.tar\.gz)"`)
	var installed []string
	lines := strings.Split(string(src), "\n")
	for i, line := range lines {
		mm := base.FindStringSubmatch(line)
		if mm == nil {
			continue
		}
		installed = append(installed, mm[1])
		r, ok := m.Lookup(mm[1])
		if !ok {
			t.Errorf("ci.yml:%d installs %s, not in the manifest", i+1, mm[1])
			continue
		}
		next := lines[i+1]
		am := asset.FindStringSubmatch(next)
		if am == nil || am[1] != r.Assets["linux/amd64"].Asset {
			t.Errorf("ci.yml:%d fetches %q for %s, manifest linux/amd64 asset is %q", i+2, next, mm[1], r.Assets["linux/amd64"].Asset)
		}
	}
	// Job order: ailang-verify (frozen compiler), go-verify (gate primary, then the floor tool),
	// go-race-verify (gate primary again; row 181 moved the -race leg into its own job).
	want := []string{ReleaseFor(RoleFrozenCompiler).Release, ReleaseFor(RoleGatePrimary).Release, ReleaseFor(RoleTool).Release, ReleaseFor(RoleGatePrimary).Release}
	if strings.Join(installed, " ") != strings.Join(want, " ") {
		t.Errorf("ci.yml installs %v (in order), manifest frozen-compiler/gate-primary/tool/gate-primary(race job) are %v", installed, want)
	}
	// The tool install pins the manifest's linux/amd64 tarball AND binary
	// digests (row 139 D7), on top of the release's own .sha256.
	ta := ReleaseFor(RoleTool).Assets["linux/amd64"]
	for _, line := range []string{
		`echo "` + ta.TarballSHA256 + `  tool.tar.gz" | sha256sum -c -`,
		`echo "` + ta.BinarySHA256 + `  $dir/ailang" | sha256sum -c -`,
	} {
		if strings.Count(string(src), line) != 1 {
			t.Errorf("ci.yml's tool install lacks the manifest digest check %q", line)
		}
	}
	// Every sha256 ci.yml carries is the z3 pin or a manifest digest, so a
	// stale digest left behind by a release move reds here.
	known := map[string]bool{}
	for _, r := range m.Releases {
		for _, a := range r.Assets {
			known[a.TarballSHA256], known[a.BinarySHA256] = true, true
		}
	}
	hex := regexp.MustCompile(`\b[0-9a-f]{64}\b`)
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "Z3_SHA:") {
			continue
		}
		for _, d := range hex.FindAllString(line, -1) {
			if !known[d] {
				t.Errorf("ci.yml:%d carries sha256 %s, which is no manifest digest", i+1, d)
			}
		}
	}
}

// ailang_release_observed.txt lists exactly the manifest's observed releases.
func TestObservedReleasesFileMatchesManifest(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(repoRoot(t), "scripts", "testdata", "ailang_release_observed.txt"))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, line := range strings.Split(string(src), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		got = append(got, line)
	}
	var want []string
	for _, r := range Current().Holders(RoleObserved) {
		want = append(want, r.Release)
	}
	sort.Strings(got)
	sort.Strings(want)
	if len(got) == 0 || strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("ailang_release_observed.txt lists %v, manifest observed role %v", got, want)
	}
}

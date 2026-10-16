package pins

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// gateRelease is the release whose manifest row the admission arms use: the
// interpreter-role holder, read from the manifest (no literal).
func gateRelease(t *testing.T) (Release, string) {
	t.Helper()
	r := ReleaseFor(RoleInterpreter)
	return r, "darwin/arm64"
}

func verdictMap(t *testing.T) map[string]any {
	t.Helper()
	r, plat := gateRelease(t)
	a := r.Assets[plat]
	return map[string]any{
		"schema": VerdictSchema, "release": r.Release, "platform": plat,
		"binary_sha256": a.BinarySHA256, "tarball_sha256": a.TarballSHA256,
		"golden_identical": true,
		"checks": []any{
			map[string]any{"name": "verify_ail", "pass": true, "evidence": "23 modules"},
			map[string]any{"name": "replay_goldens", "pass": true},
		},
		"all_pass": true, "produced_at": "2026-10-16T12:00:00Z",
	}
}

func encodeVerdict(t *testing.T, m map[string]any) []byte {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestParseVerdictControlAndMalformations(t *testing.T) {
	v, err := ParseVerdict(encodeVerdict(t, verdictMap(t)))
	if err != nil {
		t.Fatalf("control verdict refused: %v", err)
	}
	if !v.AllPass || !v.GoldenIdentical || len(v.Checks) != 2 {
		t.Fatalf("control verdict decoded as %+v", v)
	}
	for name, mutate := range map[string]func(map[string]any){
		"wrong schema":         func(m map[string]any) { m["schema"] = "world/other/v1" },
		"unknown field":        func(m map[string]any) { m["extra"] = 1 },
		"malformed release":    func(m map[string]any) { m["release"] = "latest" },
		"unsupported platform": func(m map[string]any) { m["platform"] = "plan9/386" },
		"short digest":         func(m map[string]any) { m["binary_sha256"] = "abc" },
		"absent golden flag":   func(m map[string]any) { delete(m, "golden_identical") },
		"absent all_pass":      func(m map[string]any) { delete(m, "all_pass") },
		"no checks":            func(m map[string]any) { m["checks"] = []any{} },
		"nameless check":       func(m map[string]any) { m["checks"] = []any{map[string]any{"name": "", "pass": true}} },
		"duplicate check": func(m map[string]any) {
			m["checks"] = []any{map[string]any{"name": "a", "pass": true}, map[string]any{"name": "a", "pass": true}}
		},
		"all_pass contradicted": func(m map[string]any) { m["checks"] = []any{map[string]any{"name": "a", "pass": false}} },
		"bad produced_at":       func(m map[string]any) { m["produced_at"] = "yesterday" },
	} {
		m := verdictMap(t)
		mutate(m)
		_, err := ParseVerdict(encodeVerdict(t, m))
		var ve *VerdictError
		if !errors.As(err, &ve) || ve.Reason != VerdictMalformed {
			t.Errorf("%s: err = %v, want a %s VerdictError", name, err, VerdictMalformed)
		}
	}
}

func TestAdmitVerdictReasons(t *testing.T) {
	m := Current()
	r, plat := gateRelease(t)
	admit := func(mutate func(map[string]any), release string, candidate bool) error {
		vm := verdictMap(t)
		if mutate != nil {
			mutate(vm)
		}
		v, err := ParseVerdict(encodeVerdict(t, vm))
		if err != nil {
			t.Fatalf("instrument: mutated verdict did not parse: %v", err)
		}
		return m.AdmitVerdict(v, release, plat, "test-manifest", candidate)
	}
	// CONTROLS: the unmutated verdict is admitted for both modes.
	if err := admit(nil, r.Release, true); err != nil {
		t.Fatalf("control (candidate): %v", err)
	}
	if err := admit(nil, r.Release, false); err != nil {
		t.Fatalf("control (new epoch): %v", err)
	}
	observedOnly := ""
	for _, rel := range m.Releases {
		if len(rel.Assets) == 0 {
			observedOnly = rel.Release
		}
	}
	if observedOnly == "" {
		t.Fatal("instrument: the manifest has no asset-less release")
	}
	other := ReleaseFor(RoleTool)
	if other.Release == r.Release {
		t.Fatal("instrument: tool and interpreter releases coincide")
	}
	failing := func(vm map[string]any) {
		vm["checks"] = []any{map[string]any{"name": "replay_goldens", "pass": false}, map[string]any{"name": "verify_ail", "pass": true}}
		vm["all_pass"] = false
	}
	for _, c := range []struct {
		name      string
		mutate    func(map[string]any)
		release   string
		candidate bool
		reason    string
		detail    string
	}{
		{"names another release", nil, other.Release, true, VerdictReleaseMismatch, r.Release},
		{"other platform", func(vm map[string]any) { vm["platform"] = "linux/amd64" }, r.Release, true, VerdictPlatform, "linux/amd64"},
		{"unknown release", func(vm map[string]any) { vm["release"] = "v0.0.1" }, "v0.0.1", true, VerdictUnknownRelease, "test-manifest"},
		{"no manifest digest", func(vm map[string]any) { vm["release"] = observedOnly }, observedOnly, true, VerdictNoDigest, observedOnly},
		{"binary digest mismatch", func(vm map[string]any) { vm["binary_sha256"] = strings.Repeat("0", 64) }, r.Release, true, VerdictDigestMismatch, "binary"},
		{"tarball digest mismatch", func(vm map[string]any) { vm["tarball_sha256"] = strings.Repeat("0", 64) }, r.Release, false, VerdictDigestMismatch, "tarball"},
		{"a check failed", failing, r.Release, false, VerdictNotAllPass, "replay_goldens"},
		{"goldens differ for a candidate", func(vm map[string]any) { vm["golden_identical"] = false }, r.Release, true, VerdictGoldensDiffer, "NEW epoch"},
	} {
		err := admit(c.mutate, c.release, c.candidate)
		var ve *VerdictError
		if !errors.As(err, &ve) || ve.Reason != c.reason || !strings.Contains(ve.Detail, c.detail) {
			t.Errorf("%s: err = %v, want reason %s naming %q", c.name, err, c.reason, c.detail)
		}
	}
	// golden_identical=false is admissible for a NEW epoch (that is what one is for).
	if err := admit(func(vm map[string]any) { vm["golden_identical"] = false }, r.Release, false); err != nil {
		t.Fatalf("goldens differ for a new epoch: %v", err)
	}
}

func TestReleaseOfVersionLine(t *testing.T) {
	r := ReleaseFor(RoleInterpreter)
	if got, ok := ReleaseOfVersionLine(r.Version); !ok || got != r.Release {
		t.Fatalf("ReleaseOfVersionLine(%q) = %q %v", r.Version, got, ok)
	}
	for _, bad := range []string{"unpinned", r.Release, "AILANG " + r.Release + "-1-gabc-dirty", "AILANG  " + r.Release} {
		if _, ok := ReleaseOfVersionLine(bad); ok {
			t.Errorf("accepted %q", bad)
		}
	}
}

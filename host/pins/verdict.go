package pins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// VerdictSchema is the schema id of a per-release compatibility-gate verdict
// (row 139 D4, design_docs/planned/w-interpreter-epoch-upgrade-path.md). The
// gate (scripts/ailang_release_check.sh, a later milestone) writes one; the
// `world-publish epoch nominate` verb (D2) refuses without one that
// AdmitVerdict accepts.
//
// Wire shape (every field required; unknown fields refused):
//
//	{
//	  "schema":           "world/ailang-release-verdict/v1",
//	  "release":          "vX.Y.Z",
//	  "platform":         "darwin/arm64" | "linux/amd64",
//	  "binary_sha256":    "<64 hex>",   // the `ailang` binary the checks ran
//	  "tarball_sha256":   "<64 hex>",   // the release asset it came from
//	  "golden_identical": true | false, // every replay golden byte-identical
//	  "checks":           [{"name": "...", "pass": true|false, "evidence": "..."}],
//	  "all_pass":         true | false, // must equal the AND of checks[].pass
//	  "produced_at":      "<RFC 3339>"
//	}
const VerdictSchema = "world/ailang-release-verdict/v1"

// VerdictCheck is one named gate check.
type VerdictCheck struct {
	Name     string `json:"name"`
	Pass     bool   `json:"pass"`
	Evidence string `json:"evidence,omitempty"`
}

// Verdict is a parsed, structurally valid gate verdict.
type Verdict struct {
	Schema          string         `json:"schema"`
	Release         string         `json:"release"`
	Platform        string         `json:"platform"`
	BinarySHA256    string         `json:"binary_sha256"`
	TarballSHA256   string         `json:"tarball_sha256"`
	GoldenIdentical bool           `json:"golden_identical"`
	Checks          []VerdictCheck `json:"checks"`
	AllPass         bool           `json:"all_pass"`
	ProducedAt      string         `json:"produced_at"`
}

// verdictWire mirrors Verdict with pointer booleans so an ABSENT
// golden_identical or all_pass is refused rather than read as false.
type verdictWire struct {
	Schema          string         `json:"schema"`
	Release         string         `json:"release"`
	Platform        string         `json:"platform"`
	BinarySHA256    string         `json:"binary_sha256"`
	TarballSHA256   string         `json:"tarball_sha256"`
	GoldenIdentical *bool          `json:"golden_identical"`
	Checks          []VerdictCheck `json:"checks"`
	AllPass         *bool          `json:"all_pass"`
	ProducedAt      string         `json:"produced_at"`
}

// VerdictError is a verdict refusal. Reason is a stable token (the verb
// renders it as `STOP fence=verdict reason=<Reason>`).
type VerdictError struct {
	Reason string
	Detail string
}

func (e *VerdictError) Error() string { return "pins: verdict " + e.Reason + ": " + e.Detail }

// Verdict refusal reasons.
const (
	VerdictMalformed       = "malformed"
	VerdictReleaseMismatch = "release-mismatch"
	VerdictPlatform        = "platform-mismatch"
	VerdictUnknownRelease  = "unknown-release"
	VerdictNoDigest        = "no-manifest-digest"
	VerdictDigestMismatch  = "digest-mismatch"
	VerdictNotAllPass      = "not-all-pass"
	VerdictGoldensDiffer   = "goldens-differ"
)

func malformed(format string, args ...any) *VerdictError {
	return &VerdictError{Reason: VerdictMalformed, Detail: fmt.Sprintf(format, args...)}
}

// ParseVerdict decodes and STRUCTURALLY validates a verdict: schema, release
// and digest shapes, a supported platform, at least one uniquely named
// check, both booleans present, all_pass consistent with the checks, and an
// RFC 3339 produced_at. It does not consult a manifest (AdmitVerdict does).
func ParseVerdict(b []byte) (Verdict, error) {
	var w verdictWire
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&w); err != nil {
		return Verdict{}, malformed("decode: %v", err)
	}
	if dec.More() {
		return Verdict{}, malformed("trailing data after the verdict object")
	}
	if w.Schema != VerdictSchema {
		return Verdict{}, malformed("schema %q, want %q", w.Schema, VerdictSchema)
	}
	if !releaseRe.MatchString(w.Release) {
		return Verdict{}, malformed("release %q is not vX.Y.Z", w.Release)
	}
	if _, ok := assetName[w.Platform]; !ok {
		return Verdict{}, malformed("platform %q is not a supported platform", w.Platform)
	}
	if !hex64.MatchString(w.BinarySHA256) || !hex64.MatchString(w.TarballSHA256) {
		return Verdict{}, malformed("binary_sha256 and tarball_sha256 must be 64 lowercase hex characters")
	}
	if w.GoldenIdentical == nil || w.AllPass == nil {
		return Verdict{}, malformed("golden_identical and all_pass are required")
	}
	if len(w.Checks) == 0 {
		return Verdict{}, malformed("checks is empty: a verdict names every check it ran")
	}
	all := true
	names := map[string]bool{}
	for _, c := range w.Checks {
		if strings.TrimSpace(c.Name) == "" {
			return Verdict{}, malformed("a check has no name")
		}
		if names[c.Name] {
			return Verdict{}, malformed("check %q listed twice", c.Name)
		}
		names[c.Name] = true
		all = all && c.Pass
	}
	if *w.AllPass != all {
		return Verdict{}, malformed("all_pass=%v contradicts the checks (their AND is %v)", *w.AllPass, all)
	}
	if _, err := time.Parse(time.RFC3339, w.ProducedAt); err != nil {
		return Verdict{}, malformed("produced_at %q is not RFC 3339", w.ProducedAt)
	}
	return Verdict{
		Schema: w.Schema, Release: w.Release, Platform: w.Platform,
		BinarySHA256: w.BinarySHA256, TarballSHA256: w.TarballSHA256,
		GoldenIdentical: *w.GoldenIdentical, Checks: w.Checks, AllPass: *w.AllPass,
		ProducedAt: w.ProducedAt,
	}, nil
}

// AdmitVerdict is the nomination gate over a parsed verdict (row 139 D2):
// it names release; it was produced on platform; release is in m with an
// asset for platform whose binary AND tarball digests equal the verdict's;
// every check passed; and, when requireGoldenIdentical (a candidate
// nomination), every replay golden was byte-identical. manifestName is how
// a refusal names the manifest it consulted.
func (m Manifest) AdmitVerdict(v Verdict, release, platform, manifestName string, requireGoldenIdentical bool) error {
	if v.Release != release {
		return &VerdictError{Reason: VerdictReleaseMismatch,
			Detail: fmt.Sprintf("the verdict is for %s, not %s", v.Release, release)}
	}
	if v.Platform != platform {
		return &VerdictError{Reason: VerdictPlatform,
			Detail: fmt.Sprintf("the verdict was produced on %s; this host is %s", v.Platform, platform)}
	}
	r, ok := m.Lookup(release)
	if !ok {
		return &VerdictError{Reason: VerdictUnknownRelease,
			Detail: fmt.Sprintf("%s is not in the release manifest %s (releases: %v); add its row first", release, manifestName, m.Names())}
	}
	a, ok := r.Assets[platform]
	if !ok {
		return &VerdictError{Reason: VerdictNoDigest,
			Detail: fmt.Sprintf("the release manifest %s records no %s digest for %s", manifestName, platform, release)}
	}
	if v.BinarySHA256 != a.BinarySHA256 || v.TarballSHA256 != a.TarballSHA256 {
		return &VerdictError{Reason: VerdictDigestMismatch,
			Detail: fmt.Sprintf("the verdict's %s digests (binary %s, tarball %s) differ from the release manifest %s (binary %s, tarball %s)",
				platform, v.BinarySHA256, v.TarballSHA256, manifestName, a.BinarySHA256, a.TarballSHA256)}
	}
	if !v.AllPass {
		var failed []string
		for _, c := range v.Checks {
			if !c.Pass {
				failed = append(failed, c.Name)
			}
		}
		return &VerdictError{Reason: VerdictNotAllPass,
			Detail: fmt.Sprintf("failed check(s): %s", strings.Join(failed, ", "))}
	}
	if requireGoldenIdentical && !v.GoldenIdentical {
		return &VerdictError{Reason: VerdictGoldensDiffer,
			Detail: "golden_identical is false: a release whose replay goldens differ can only open a NEW epoch (--new-epoch, human-typed), never join the current one as a candidate"}
	}
	return nil
}

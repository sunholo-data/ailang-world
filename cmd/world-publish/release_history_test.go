package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// changelogSections splits CHANGELOG.md into its `## <version>` sections, in
// file order, keyed by the version token and holding the trimmed body.
func changelogSections(text string) (order []string, bodies map[string]string) {
	bodies = map[string]string{}
	parts := strings.Split("\n"+text, "\n## ")
	for _, part := range parts[1:] {
		heading, body, _ := strings.Cut(part, "\n")
		version, _, _ := strings.Cut(heading, " ")
		order = append(order, version)
		bodies[version] = strings.TrimSpace(body)
	}
	return order, bodies
}

// TestChangelogCarriesTheCandidateAndPreservesPublishedHistory binds three
// things an attended release must not get wrong:
//   - the NEWEST section is the version this command publishes (PUB001 reads it);
//   - that version is not one the registry already serves (immutable, 409);
//   - the published 0.1.0 notes are byte-identical to what the registry recorded,
//     read from the preserved served record, never rewritten to suit a release.
func TestChangelogCarriesTheCandidateAndPreservesPublishedHistory(t *testing.T) {
	root := commandRepoRoot(t)
	changelog, err := os.ReadFile(filepath.Join(root, defaultPackageDir, "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	order, bodies := changelogSections(string(changelog))
	if len(order) < 2 {
		t.Fatalf("instrument failure: %d CHANGELOG sections parsed, want >= 2 (%v)", len(order), order)
	}
	if order[0] != frozenPackageVersion || bodies[frozenPackageVersion] == "" {
		t.Fatalf("newest CHANGELOG section is %q (empty=%v), want a non-empty %q",
			order[0], bodies[order[0]] == "", frozenPackageVersion)
	}

	raw, err := os.ReadFile(filepath.Join(root, "host", "broker", "testdata", "metadata_world_core_0.1.0.json"))
	if err != nil {
		t.Fatal(err)
	}
	var served struct {
		Version string `json:"version"`
		Quality struct {
			Release struct {
				Notes string `json:"notes"`
			} `json:"release"`
		} `json:"quality"`
	}
	if err := json.Unmarshal(raw, &served); err != nil {
		t.Fatal(err)
	}
	if served.Version != "0.1.0" || served.Quality.Release.Notes == "" {
		t.Fatalf("instrument failure: preserved record is %q with %d bytes of notes",
			served.Version, len(served.Quality.Release.Notes))
	}
	if frozenPackageVersion == served.Version {
		t.Fatalf("frozenPackageVersion %q is already published; the registry is immutable", served.Version)
	}
	if got := bodies[served.Version]; got != strings.TrimSpace(served.Quality.Release.Notes) {
		t.Fatalf("CHANGELOG `## %s` no longer equals the notes the registry published "+
			"(%d vs %d bytes): published history was rewritten", served.Version, len(got),
			len(served.Quality.Release.Notes))
	}
}

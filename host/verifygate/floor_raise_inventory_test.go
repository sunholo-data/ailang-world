package verifygate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	// siteReScript matches `#   N. path` rows in the verify_ail.sh inventory block.
	siteReScript = regexp.MustCompile(`^#   [0-9]+\.\s+(\S+)`)
	// siteReTable matches `| N | `path` |` rows in the §S8 table; capture up to the
	// next pipe and trim backticks/space so the canonical path is the bare file path.
	siteReTable = regexp.MustCompile(`^\| [0-9]+ \| ([^|]+)`)
)

// canonicalSiteSet extracts the ordered set of coupled-site file paths from a home's
// row lines. scriptHome=true parses `#   N. path` rows; false parses `| N | `path` |`
// table rows. It returns the paths in row order; callers compare them as sets.
func canonicalSiteSet(home string, scriptHome bool) []string {
	re := siteReTable
	if scriptHome {
		re = siteReScript
	}
	var sites []string
	for _, line := range strings.Split(home, "\n") {
		m := re.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		sites = append(sites, strings.Trim(m[1], "` \t"))
	}
	return sites
}

func TestFloorRaiseInventoryNamesEveryCoupledFile(t *testing.T) {
	scriptBytes, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "verify_ail.sh"))
	if err != nil {
		t.Fatalf("read verify script: %v", err)
	}
	script := string(scriptBytes)
	beginMarker := "# ── FLOOR-RAISE COUPLING INVENTORY"
	endMarker := "# ── END FLOOR-RAISE COUPLING INVENTORY"
	beginCount := strings.Count(script, beginMarker)
	endCount := strings.Count(script, endMarker)
	if beginCount != 1 {
		t.Fatalf("begin marker count=%d, want 1", beginCount)
	}
	if endCount != 1 {
		t.Fatalf("END marker count=%d, want 1", endCount)
	}
	begin := strings.Index(script, beginMarker)
	end := strings.Index(script, endMarker)
	if begin >= end {
		t.Fatalf("inventory markers misordered: begin=%d end=%d", begin, end)
	}
	block := script[begin : end+len(endMarker)]
	scriptRows := []string{
		"#   1. world/<module>.ail",
		"#   2. packages/world-core/world/<module>.ail",
		"#   3. scripts/verify_ail.sh",
		"#   4. scripts/world_package_ready_packet.golden.json",
		"#   5. docs/SELF_MOD_PUBLISH.md",
		"#   6. host/verifygate/module_manifest_gate_test.go",
	}
	for site, needle := range scriptRows {
		// This evaluation showed that a needle appearing twice cannot detect its
		// own row's deletion, so every row anchor must be unique in its bounded home.
		if count := strings.Count(block, needle); count != 1 {
			t.Errorf("inventory block site %d row %q count=%d, want exactly 1", site+1, needle, count)
		}
	}
	sharedNeedles := []string{
		"packages/world-core/world/",
		"REQUIRED_VERIFIED",
		"EXACT_TOTAL_VERIFIED",
		"world_package_ready_packet.golden.json",
		"SELF_MOD_PUBLISH.md",
		"module_manifest_gate_test.go",
		"interfaceHash",
		"does not move for",
	}
	for _, needle := range sharedNeedles {
		if !strings.Contains(block, needle) {
			t.Errorf("inventory block omits %q", needle)
		}
	}

	standardsBytes, err := os.ReadFile(filepath.Join(repoRoot, "design_docs", "coding-standards.md"))
	if err != nil {
		t.Fatalf("read coding standards: %v", err)
	}
	standards := string(standardsBytes)
	heading := "## S8"
	s8Start := strings.Index(standards, heading)
	if s8Start < 0 {
		t.Fatalf("%s heading not found", heading)
	}
	s8 := standards[s8Start:]
	if next := strings.Index(s8[len(heading):], "\n## "); next >= 0 {
		s8 = s8[:len(heading)+next]
	}
	s8Rows := []string{
		"| 1 | `world/<module>.ail`",
		"| 2 | `packages/world-core/world/<module>.ail`",
		"| 3 | `scripts/verify_ail.sh`",
		"| 4 | `scripts/world_package_ready_packet.golden.json`",
		"| 5 | `docs/SELF_MOD_PUBLISH.md`",
		"| 6 | `host/verifygate/module_manifest_gate_test.go`",
	}
	for site, needle := range s8Rows {
		// This evaluation showed that a needle appearing twice cannot detect its
		// own row's deletion, so every row anchor must be unique in its bounded home.
		if count := strings.Count(s8, needle); count != 1 {
			t.Errorf("S8 site %d row %q count=%d, want exactly 1", site+1, needle, count)
		}
	}
	for _, needle := range sharedNeedles {
		if !strings.Contains(s8, needle) {
			t.Errorf("S8 omits %q", needle)
		}
	}

	scriptSites := canonicalSiteSet(block, true)
	s8Sites := canonicalSiteSet(s8, false)

	// Duplicate-within-a-home guard: a duplicated coupled-site row is a defect in its own
	// right, and it is what makes a one-directional membership check evadable (a duplicate
	// in one home can mask an asymmetric addition by keeping the counts equal). Assert no
	// path repeats within EITHER home before comparing.
	scriptSet := make(map[string]bool, len(scriptSites))
	for _, s := range scriptSites {
		if scriptSet[s] {
			t.Fatalf("duplicate coupled-site path %q in the verify_ail.sh inventory block — a duplicated row is a defect", s)
		}
		scriptSet[s] = true
	}
	s8Set := make(map[string]bool, len(s8Sites))
	for _, s := range s8Sites {
		if s8Set[s] {
			t.Fatalf("duplicate coupled-site path %q in the S8 table — a duplicated row is a defect", s)
		}
		s8Set[s] = true
	}

	// Cardinality floor with instrument-failure branch: if the enumerator matches nothing,
	// both sets are empty and equal, so the set-equality below would pass vacuously. The
	// floor makes that a LOUD red. It is a MINIMUM (>= 6), never an equality, because the
	// legitimate floor-raise path adds a 7th row to BOTH homes and must pass.
	if len(scriptSites) < 6 {
		t.Fatalf("inventory block enumerator matched %d sites, want >= 6 — instrument failure, not an empty inventory", len(scriptSites))
	}
	if len(s8Sites) < 6 {
		t.Fatalf("S8 enumerator matched %d sites, want >= 6 — instrument failure, not an empty inventory", len(s8Sites))
	}
}

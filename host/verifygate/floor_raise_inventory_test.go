package verifygate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
}

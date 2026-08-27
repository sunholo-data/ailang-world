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
	for _, needle := range append(sharedNeedles, "#   1. world/<module>.ail") {
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
	for _, needle := range append(sharedNeedles, "| 1 | `world/<module>.ail`") {
		if !strings.Contains(s8, needle) {
			t.Errorf("S8 omits %q", needle)
		}
	}
}

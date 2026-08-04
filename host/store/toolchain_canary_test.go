package store

import "testing"

// TestToolchainCanary preserves the compiler-regression shape from
// design_docs/verification/w-race-gate-blindspot/repro. Go 1.26.0 through
// 1.26.5 miscompile it on darwin/arm64; the pinned-good Go 1.25.6 does not.
func TestToolchainCanary(t *testing.T) {
	type row struct {
		n      int
		field  string
		reason string
	}
	type stringer interface {
		String() string
	}
	type value string

	texts := []string{"w", ""}
	fields := [...]string{"worldRef", "stateRoot"}
	var rows []row
	for i, text := range texts {
		if len(text) == 0 {
			var v stringer = canaryString(value("empty"))
			rows = append(rows, row{
				field: fields[i], reason: v.String(),
			})
			break
		}
	}

	if len(rows) != 1 {
		t.Fatalf("len(rows)=%d want 1", len(rows))
	}
	// Check the length before formatting field: the affected compiler may emit
	// a corrupt string header, and printing it can crash the test process.
	if n := len(rows[0].field); n > 1000 {
		t.Fatalf("len(Field)=%d (corrupt string header)", n)
	}
	if rows[0].field != "stateRoot" {
		t.Fatalf("Field=%q want %q", rows[0].field, "stateRoot")
	}
}

type canaryString string

func (s canaryString) String() string { return string(s) }

package verbs

// Row 243 V2b (design §9, V1 judge notes 2 and 4b): VB6, VB7.

import (
	"strings"
	"testing"
)

// VB6: every never and other word's refusal, bare and with an argument, ends with the CLI line.
func TestNeverAndOtherReasonsEndWithTheCLILine(t *testing.T) {
	n := 0
	for _, w := range append(append([]string{}, NeverWords...), OtherWords...) {
		for _, argv := range [][]string{{w}, {w, "zzz"}} {
			v, _ := Lookup(argv)
			if want := "; at a terminal: ailang-worldd " + strings.Join(argv, " "); v.CLI == "" || !strings.HasSuffix(v.Reason, want) {
				t.Errorf("Lookup(%q).Reason = %q, want it to end %q", argv, v.Reason, want)
			}
			n++
		}
	}
	if n != 52 {
		t.Fatalf("%d rows, want 52 (26 words, bare and with an argument)", n)
	}
}

// VB7: the bar's two input bounds, as literals, so a bound changed with its test still fails.
func TestSplitBoundsAreLiterals(t *testing.T) {
	if MaxWords != 32 || MaxLine != 1024 {
		t.Fatalf("MaxWords %d MaxLine %d, want 32 and 1024", MaxWords, MaxLine)
	}
	if _, err := Split(strings.TrimSpace(strings.Repeat("w ", 33))); err == nil {
		t.Error("33 words split")
	}
	if _, err := Split(strings.Repeat("x", 1025)); err != ErrTooLong {
		t.Errorf("1,025 bytes: %v", err)
	}
}

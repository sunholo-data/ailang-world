package verbs

// Row 243 V4 sprint 2 M4 (judge-407 survivor: AdmitOperator without fs.SetOutput(io.Discard) survived
// host/verbs): the pure bar parser writes nothing to stdout or stderr, on the lines whose parse the flag
// package would otherwise report (an undefined flag, a bad bool, a missing value, help).

import (
	"io"
	"os"
	"testing"
)

func TestAdmitOperatorWritesNothing(t *testing.T) {
	capture := func(f func()) string {
		r, w, _ := os.Pipe()
		out, errOut := os.Stdout, os.Stderr
		os.Stdout, os.Stderr = w, w
		f()
		os.Stdout, os.Stderr = out, errOut
		_ = w.Close()
		b, _ := io.ReadAll(r)
		return string(b)
	}
	if got := capture(func() { _, _ = os.Stderr.WriteString("x") }); got != "x" {
		t.Fatalf("control: the harness saw %q", got)
	}
	if got := capture(func() {
		for _, line := range []string{"episode start e1 --bogus", "episode start e1 --preset p --launch=maybe", "episode finish e1 --preset",
			"episode start e1 --help", "episode finish e1 -h", "episode start e1 --preset p --launch --headless"} {
			admitLine(t, line, allYes)
		}
	}); got != "" {
		t.Fatalf("AdmitOperator wrote %q", got)
	}
}

package approval_test

// Row 243 V4 sprint 2 M4 (judge-407's retro rule): every new pure entry point a later sprint calls from
// a server writes nothing to stdout, stderr or the log, on its accepting and its refusing paths.

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// silent runs f with os.Stdout, os.Stderr and the log on one pipe and returns what was written.
func silent(t *testing.T, f func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	out, errOut := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = w, w
	log.SetOutput(w)
	f()
	os.Stdout, os.Stderr = out, errOut
	log.SetOutput(os.Stderr)
	_ = w.Close()
	b, _ := io.ReadAll(r)
	return string(b)
}

func TestOperatorPureEntryPointsWriteNothing(t *testing.T) {
	if got := silent(t, func() { _, _ = os.Stderr.WriteString("x"); _, _ = os.Stdout.WriteString("y"); log.Print("z") }); !strings.HasPrefix(got, "xy") || !strings.HasSuffix(got, " z\n") {
		t.Fatalf("control: the harness saw %q, want x, y and the log line", got)
	}
	got := silent(t, func() {
		for _, b := range [][]byte{opGoldenBytes(t), []byte(`{"a":1,"a":2}`), []byte("\xff")} {
			a, _ := approval.DecodeOperatorAction(b)
			_, _ = approval.OperatorReviewDigest(a)
			_, _ = approval.OperatorSignatureDigest(a)
		}
		_, _ = approval.OperatorActionMessageHash([32]byte{}, [32]byte{}, -1)
		for _, format := range []string{approval.FormatSSH, approval.FormatWebAuthn, "typed"} {
			in := opIn(t, format)
			_, _ = approval.ConfirmOperatorAction(in)
			_, _ = approval.CheckOperatorAction(approval.OperatorCheckInput{Action: in.Action, Challenge: opC, Expiry: opExp, Now: opNow,
				Artifact: []byte("y\n"), Format: format, Approver: in.Approver})
		}
	})
	if got != "" {
		t.Fatalf("a pure entry point wrote %q", got)
	}
}

package authority

// Row 243 V4 sprint 2b (judge-407's retro rule, as sprint 2 M4): the two operator judges S7 and the
// runner will call write nothing to stdout, stderr or the log, accepting or refusing.

import (
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ojSilent runs f with os.Stdout, os.Stderr and the log on one pipe and returns what was written.
func ojSilent(t *testing.T, f func()) string {
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

func TestOperatorJudgesWriteNothing(t *testing.T) {
	if got := ojSilent(t, func() { _, _ = os.Stderr.WriteString("x"); _, _ = os.Stdout.WriteString("y"); log.Print("z") }); !strings.HasPrefix(got, "xy") || !strings.HasSuffix(got, " z\n") {
		t.Fatalf("control: the harness saw %q, want x, y and the log line", got)
	}
	in, ssh, pk := ojIn(t)
	ro, revoked, closed := ojOpen(t), ojOpen(t, `UPDATE approvers SET revoked_at = 1000`), ojOpen(t)
	_ = closed.Close()
	arts := [][]byte{ssh, pk, []byte("y\n"), ojSSH(t, ojStranger, approval.OperatorActionNamespace, in)}
	ctx := boundedTestContext(t)
	got := ojSilent(t, func() {
		for _, h := range []*store.Credentials{ro, revoked, closed} {
			for _, art := range arts {
				_, _ = JudgeOperatorAction(ctx, h, in, art)
				for _, by := range []string{"ssh:mark-ssh", "webauthn:mark-mac"} {
					_, _ = RecheckOperatorAction(ctx, h, in.Action, in.Challenge, in.ExpiresAt, 1500, art, by)
				}
			}
		}
	})
	if got != "" {
		t.Fatalf("an operator judge wrote %q", got)
	}
}

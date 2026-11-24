//go:build unix

package authority

// Row 243 V4 sprint 2b (judge-407's retro rule, as sprint 2 M4): the two operator judges S7 and the
// runner will call write nothing to stdout, stderr or the log, accepting or refusing.

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
	"golang.org/x/sys/unix"
)

// ojSilent captures writes to descriptors 1 and 2 as well as the log writer.
// Run only in sequential tests, because the descriptors are process-wide.
func ojSilent(t *testing.T, f func()) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "fd-output")
	capture, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	defer capture.Close()
	func() {
		for _, fd := range []int{1, 2} {
			saved, err := unix.Dup(fd)
			if err != nil {
				t.Fatal(err)
			}
			defer unix.Close(saved)
			defer func() {
				if err := unix.Dup2(saved, fd); err != nil {
					t.Error(err)
				}
			}()
			if err := unix.Dup2(int(capture.Fd()), fd); err != nil {
				t.Fatal(err)
			}
		}
		logWriter := log.Writer()
		defer log.SetOutput(logWriter)
		log.SetOutput(os.Stderr)
		f()
	}()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestOperatorJudgesWriteNothing(t *testing.T) {
	if got := ojSilent(t, func() {
		_, _ = os.Stderr.WriteString("x")
		_, _ = os.Stdout.WriteString("y")
		log.Print("z")
		println("builtin")
	}); !strings.HasPrefix(got, "xy") || !strings.HasSuffix(got, " z\nbuiltin\n") {
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

func TestOperatorStoreWritesNothing(t *testing.T) {
	r := newOperatorStoreRig(t, true)
	ctx := boundedTestContext(t)
	const now = 1010
	got := ojSilent(t, func() {
		_, _ = r.RequestOperatorAction(ctx, store.OperatorRequest{Action: r.action, OperationID: "1234567890abcdef1234567890abcdef", CreatedAt: 1000, ExpiresAt: 1600}, now)
		_, _ = r.OperatorAction(ctx, r.cas.ID, now)
		_, _ = r.OperatorAction(ctx, "missing", now)
		_, _ = r.ConfirmOperatorAction(ctx, r.cas, r.artifact, nil, now)
		r.confirm(t)
		_ = r.ClaimOperatorAction(ctx, r.cas, store.OperatorOwner{}, now)
		jgExec(t, r.path, `INSERT INTO operator_runner_state VALUES(1, 7, 'runner', 1, 'test', 1010, 'start', 1, NULL)`)
		if err := r.ClaimOperatorAction(ctx, r.cas, store.OperatorOwner{Nonce: "runner", Generation: 7}, now); err != nil {
			t.Fatal(err)
		}
		if err := r.PruneOperatorActions(ctx, now); err != nil {
			t.Fatal(err)
		}
	})
	if got != "" {
		t.Fatalf("an operator store method wrote %q", got)
	}
}

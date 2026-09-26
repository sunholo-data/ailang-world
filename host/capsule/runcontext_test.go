package capsule

import (
	"context"
	"strings"
	"testing"
	"time"
)

// AC-CAP-ARGS: Args reaches main(input) through --args-file.
func TestRunContextArgs(t *testing.T) {
	fixture := archivePinned(t)
	src := source(`export func main(input: string) -> string {
  input
}`)
	res, err := New(fixture.archive, Config{}).RunContext(context.Background(), Entry{
		Interpreter: fixture.ref, Source: src, Args: []byte(`"{\"a\":1}"`),
	})
	if err != nil {
		t.Fatalf("RunContext: %v (stderr %q)", err, res.Stderr)
	}
	if got := strings.TrimSuffix(string(res.Stdout), "\n"); got != `{"a":1}` {
		t.Fatalf("stdout = %q, want the argument echoed", res.Stdout)
	}
}

// AC-CAP-CTX: cancelling the caller's ctx kills the capsule long before the
// exec allowance.
func TestRunContextCancel(t *testing.T) {
	fixture := archivePinned(t)
	// fib(32) runs for tens of seconds on the pinned interpreter (fib(30)
	// measured ~11 s), so a prompt return can only come from the ctx.
	src := source(`func fib(n: int) -> int {
  if n < 2 then n else fib(n - 1) + fib(n - 2)
}

export func main() -> int {
  fib(32)
}`)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := New(fixture.archive, Config{ExecTimeout: 50 * time.Second}).RunContext(ctx, Entry{
		Interpreter: fixture.ref, Source: src,
	})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("RunContext returned nil for a cancelled non-terminating capsule")
	}
	if elapsed > 5*time.Second {
		t.Fatalf("cancellation took %s; the caller's ctx does not bound the child", elapsed)
	}
}

// AC-CAP-RELAX: a module header that does not name the staging path still
// runs — the publish check's shape (V9).
func TestRunContextRelaxedModule(t *testing.T) {
	fixture := archivePinned(t)
	src := []byte("module transitions/echo\n\nexport func main(input: string) -> string {\n  input\n}\n")
	res, err := New(fixture.archive, Config{}).RunContext(context.Background(), Entry{
		Interpreter: fixture.ref, Source: src, Args: []byte(`"{}"`),
	})
	if err != nil {
		t.Fatalf("RunContext with a relaxed module header: %v (stderr %q)", err, res.Stderr)
	}
}

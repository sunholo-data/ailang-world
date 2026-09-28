package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type delayedApprovalInput struct {
	io.Reader
	delay  time.Duration
	waited bool
}

func (r *delayedApprovalInput) Read(p []byte) (int, error) {
	if !r.waited {
		r.waited = true
		time.Sleep(r.delay)
	}
	return r.Reader.Read(p)
}
func budgetSource(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
func requireBudgetOrder(t *testing.T, src string, phrases ...string) {
	t.Helper()
	at := -1
	for _, phrase := range phrases {
		next := strings.Index(src[at+1:], phrase)
		if next < 0 {
			t.Fatalf("missing budget step %q", phrase)
		}
		at += next + 1
	}
}

func TestApprovalBudgetExcludesPromptTime(t *testing.T) {
	src := budgetSource(t, "main.go")
	requireBudgetOrder(t, src, "func runApprove(", "requireAttendedOperator(in, out, env.getenv, env.probe())", "context.WithTimeout(context.Background(), 5*time.Second)", "broker.MintAttendedApproval(activeCtx, db, plan)")
	root := commandRepoRoot(t)
	opts := options{store: filepath.Join(t.TempDir(), "world.db"), packageDir: filepath.Join(root, defaultPackageDir), golden: filepath.Join(root, defaultGolden), publisher: filepath.Join(root, defaultGolden), registryOrigin: "https://registry.example", episode: "budget-approve", requester: "operator", decidedBy: "operator", now: 10, expires: 100}
	input := &delayedApprovalInput{Reader: strings.NewReader(attendedPhrase + "\n"), delay: 10*time.Second + 100*time.Millisecond}
	var out, errw bytes.Buffer
	code := runApprove(opts, input, &out, &errw, environment{getenv: noEnv, probe: func() ttyProbe { return satisfiedProbe(t) }})
	if code != exitOK {
		t.Fatalf("approval after >2×B2 think time: exit=%d stderr=%s", code, errw.String())
	}
	if !strings.Contains(out.String(), "minted approval sha256:") {
		t.Fatalf("missing minted approval: %s", out.String())
	}
}

func TestReconcileScanBudget(t *testing.T) {
	src := budgetSource(t, "main.go")
	requireBudgetOrder(t, src, "func runReconcile(", "context.WithTimeout(context.Background(), 3*time.Second)", "db.PendingEffectIntents(scanCtx, reconcileScanLimit)", "cancelScan()", "reconcileReadOnlyProbe(opts, packet)")
	if !strings.Contains(budgetSource(t, "reconcile_probe.go"), "broker.ReconcileRegistryPublish(context.Background()") {
		t.Fatal("read-only probe root moved without the broker call")
	}
	if strings.Contains(src, "scanCtx := context.Background()") {
		t.Fatal("temporary unbounded reconcile root remains")
	}
}

func TestTransitionsRootBudget(t *testing.T) {
	src := budgetSource(t, "transitions.go")
	requireBudgetOrder(t, src, "func runTransitions(", "readManifest(opts.manifest, errw)", "context.WithTimeout(context.Background(), time.Duration(len(entries))*10*time.Second+3*time.Second)", "buildChanges(ctx, db, arch, entries, interpreter)")
}

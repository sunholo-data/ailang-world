package broker

// Row 140 M3 (w-workspace-exec-toolchain-effect §4.6, §7) against the REAL
// pinned srt: the startup probe passes on a correct stack (darwin Seatbelt,
// linux bubblewrap), and a sandbox dir holding only the srt package — no
// dependencies — is refused at startup (MUT-ARCHIVE-PKG-ONLY: srt cannot
// load, ERR_MODULE_NOT_FOUND: commander, V42). Like the M2 real-srt tests,
// these skip only when srt or node is absent, by name; CI's go job runs them
// verbosely and fails unless each prints --- PASS.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The whole §4.6 probe passes against the real pinned srt with the
// production rendering and the real operator HOME. MUT-NET-ON (a profile
// opening allowedDomains) and the archive mutants make this red.
func TestExecSrtStartupProbePasses(t *testing.T) {
	nm, _ := realSrt(t)
	pin := ExecSandboxPin{Version: ExecSandboxRelease, CLISHA256: ExecSandboxCLISHA256}
	f := newProbeFixture(t, nm, &pin, realHome(t), []any{"sh", "-c", "echo toolchain-ok"})
	start := time.Now()
	if err := f.probe(t, ""); err != nil {
		t.Fatalf("the startup probe refused the real pinned srt: %v", err)
	}
	// It runs inside serve's 9 s startup bound.
	t.Logf("TIMING the startup probe took %s", time.Since(start))
	f.assertProbeCleaned(t)
}

// copyDir copies src to dst (files and dirs; links as links).
func copyDir(t *testing.T, src, dst string) {
	t.Helper()
	if err := copyTree(src, dst); err != nil {
		t.Fatal(err)
	}
}

// MUT-ARCHIVE-PKG-ONLY's startup arm: --exec-sandbox naming a node_modules
// that holds only @anthropic-ai/sandbox-runtime (version and cli.js digest
// both the pin's) is refused at startup by the probe: srt does not load.
func TestExecSrtProbeRefusesAPackageOnlySandbox(t *testing.T) {
	nm, _ := realSrt(t)
	pkgOnly := filepath.Join(canonicalTempDir(t), "node_modules")
	copyDir(t, filepath.Join(nm, "@anthropic-ai", "sandbox-runtime"), filepath.Join(pkgOnly, "@anthropic-ai", "sandbox-runtime"))
	if entries, _ := os.ReadDir(pkgOnly); len(entries) != 1 {
		t.Fatalf("the package-only tree holds %d entries, want only @anthropic-ai", len(entries))
	}
	pin := ExecSandboxPin{Version: ExecSandboxRelease, CLISHA256: ExecSandboxCLISHA256}
	f := newProbeFixture(t, pkgOnly, &pin, "", []any{"sh", "-c", "echo toolchain-ok"})
	err := f.probe(t, "")
	got := failedArms(t, err)
	if len(got) == 0 || got[0] != ExecProbeArmWriteInside {
		t.Fatalf("failed arms %v, want arm 1 first: srt itself must not load", got)
	}
	if !strings.Contains(err.Error(), "ERR_MODULE_NOT_FOUND") {
		t.Fatalf("the refusal does not carry srt's load failure: %v", err)
	}
	f.assertProbeCleaned(t)
}

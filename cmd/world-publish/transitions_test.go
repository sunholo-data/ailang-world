package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/canon"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// driveTransitions runs the transitions verb in-process with an injected
// environment. stdin carries the typed confirmation phrase.
func driveTransitions(t *testing.T, flags map[string]string, stdin string, getenv func(string) string) armResult {
	t.Helper()
	return drive(t, invocation{verb: "transitions", flags: flags}, stdin, getenv, satisfiedProbe(t))
}

func writeTransitionManifest(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// echoManifest builds a one-entry manifest; withEpoch adds an explicit
// semanticsEpoch, without it the epoch must be DERIVED from the registry.
func echoManifest(t *testing.T, sourcePath string, withEpoch bool) string {
	t.Helper()
	epoch := ""
	if withEpoch {
		epoch = ",\n\t\t\"semanticsEpoch\": 1"
	}
	return writeTransitionManifest(t, `[{
		"id": "tools.echo",
		"title": "Echo",
		"description": "test transition",
		"transitionFnFile": "`+sourcePath+`"`+epoch+`,
		"inputSchema": {"type":"object"},
		"outputSchema": {"type":"object"},
		"access": {"effect": "world.apply", "scope": "world", "cost": 1},
		"declaredEffects": [{"effect": "world.apply", "scope": "world", "cost": 1}]
	}]`)
}

// fakeInterpreter writes the house-pattern shell-script fake interpreter
// (the same pattern host/archive/archive_test.go uses) and returns its path.
// It prints the version for --version and exits 0 for everything else, so its
// `check` accepts every source.
func fakeInterpreter(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fake-ailang")
	script := "#!/bin/sh\n" +
		"echo test-interpreter-version\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// refusingInterpreter is the fake whose `check` exits 1: it models the pinned
// interpreter refusing a source (parse error, type error, LDR001).
func refusingInterpreter(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "refusing-ailang")
	script := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then echo test-interpreter-version; exit 0; fi\n" +
		"echo \"Error: this fake refuses the check\" >&2\n" +
		"exit 1\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// bootstrapEpochRegistry seeds a store's epoch registry with a release string,
// the way the daemon does at startup (registry.Bootstrap is idempotent).
func bootstrapEpochRegistry(t *testing.T, storePath, release string) {
	t.Helper()
	db, err := store.Open(storePath)
	if err != nil {
		t.Fatalf("open store for bootstrap: %v", err)
	}
	defer func() { _ = db.Close() }()
	if _, _, err := registry.Bootstrap(context.Background(), db, release); err != nil {
		t.Fatalf("bootstrap epoch registry: %v", err)
	}
}

// The happy path, end to end in-process: a real store, a real AILANG source
// file canonicalised and stored as an object, a REAL interpreter archived from
// the shell-script fake (the probe runs `--version`), the epoch registry
// bootstrapped with the fake's release, and the production PublishSet — with
// the semantics epoch DERIVED from the registry (the manifest omits it). Then
// the IDEMPOTENCE arm: an identical second run must be a no-op, not a new
// revision.
func TestTransitionsVerbHappyPathAndIdempotence(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "world.db")
	srcPath := filepath.Join(t.TempDir(), "echo.ail")
	if err := os.WriteFile(srcPath, []byte("module world/echo\n\nexport func apply(w: World) -> World ! {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest := echoManifest(t, srcPath, false) // semanticsEpoch omitted: derived
	bootstrapEpochRegistry(t, storePath, "test-interpreter-version")
	flags := map[string]string{"store": storePath, "manifest": manifest, "ailang-bin": fakeInterpreter(t)}

	res := driveTransitions(t, flags, attendedPhrase+"\n", noEnv)
	if res.code != exitOK {
		t.Fatalf("happy path = (%d, %s, %s), want exit 0", res.code, res.stdout, res.stderr)
	}
	if !strings.Contains(res.stdout, "published transition registry revision 1") {
		t.Fatalf("happy path stdout = %q, want the revision-1 publication line", res.stdout)
	}
	if !strings.Contains(res.stdout, "semantics epoch 1 derived from world/epoch-registry/v1 for interpreter release") {
		t.Fatalf("happy path stdout = %q, want the derived-epoch line (no typed/default epoch)", res.stdout)
	}
	// The published head is readable through the production reader.
	if !strings.Contains(res.stdout, "sha256:") {
		t.Fatalf("happy path stdout = %q, want the head ref printed", res.stdout)
	}

	again := driveTransitions(t, flags, attendedPhrase+"\n", noEnv)
	if again.code != exitOK {
		t.Fatalf("republish = (%d, %s), want exit 0", again.code, again.stderr)
	}
	if !strings.Contains(again.stdout, "UNCHANGED at revision 1") {
		t.Fatalf("republish stdout = %q, want the idempotent no-op line at revision 1", again.stdout)
	}
}

// The epoch arms (objection A at the verb): an explicit epoch must be one the
// registry nominates for the interpreter's release; an omitted epoch with NO
// nomination is refused (never defaulted to 1); an omitted epoch against a
// store with no epoch registry at all is refused with the typed reason.
func TestTransitionsVerbEpochArms(t *testing.T) {
	srcPath := filepath.Join(t.TempDir(), "echo.ail")
	if err := os.WriteFile(srcPath, []byte("module world/echo\n\nexport func apply(w: World) -> World ! {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Run("explicit mismatched epoch refused", func(t *testing.T) {
		storePath := filepath.Join(t.TempDir(), "world.db")
		bootstrapEpochRegistry(t, storePath, "test-interpreter-version")
		manifest := echoManifest(t, srcPath, true) // semanticsEpoch: 1 — but…
		body, err := os.ReadFile(manifest)
		if err != nil {
			t.Fatal(err)
		}
		// …make it mismatched: epoch 2 against an epoch-1 registry.
		fixed := strings.Replace(string(body), `"semanticsEpoch": 1`, `"semanticsEpoch": 2`, 1)
		manifest = writeTransitionManifest(t, fixed)
		res := driveTransitions(t, map[string]string{"store": storePath, "manifest": manifest, "ailang-bin": fakeInterpreter(t)},
			attendedPhrase+"\n", noEnv)
		if res.code != exitError || !strings.Contains(res.stderr, "semanticsEpoch 2 is not one of the epochs nominating") {
			t.Fatalf("mismatched epoch = (%d, %q), want exit 1 naming the epochs that nominate", res.code, res.stderr)
		}
	})

	t.Run("omitted epoch with no nomination refused (no default)", func(t *testing.T) {
		storePath := filepath.Join(t.TempDir(), "world.db")
		bootstrapEpochRegistry(t, storePath, "some-other-release") // epoch 1 nominates a DIFFERENT release
		manifest := echoManifest(t, srcPath, false)
		res := driveTransitions(t, map[string]string{"store": storePath, "manifest": manifest, "ailang-bin": fakeInterpreter(t)},
			attendedPhrase+"\n", noEnv)
		if res.code != exitError || !strings.Contains(res.stderr, "nominated by NO epoch") {
			t.Fatalf("no-nomination epoch = (%d, %q), want exit 1 refusing (never defaulted)", res.code, res.stderr)
		}
	})

	t.Run("absent epoch registry refused typed", func(t *testing.T) {
		storePath := filepath.Join(t.TempDir(), "world.db") // never bootstrapped
		manifest := echoManifest(t, srcPath, false)
		res := driveTransitions(t, map[string]string{"store": storePath, "manifest": manifest, "ailang-bin": fakeInterpreter(t)},
			attendedPhrase+"\n", noEnv)
		if res.code != exitError || !strings.Contains(res.stderr, "no epoch registry head exists") {
			t.Fatalf("absent registry = (%d, %q), want exit 1 naming the epoch registry's absence", res.code, res.stderr)
		}
	})
}

// TestTransitionsVerbRefusesGarbageSourceBeforePutObject is the objection-B
// verb arm: source bytes that canon.Source ACCEPTS (clean UTF-8, no NUL —
// canonicalisation alone proves nothing) but the pinned interpreter's check
// REFUSES are rejected BEFORE PutObject: the store holds no source object at
// all afterwards.
func TestTransitionsVerbRefusesGarbageSourceBeforePutObject(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "world.db")
	bootstrapEpochRegistry(t, storePath, "test-interpreter-version")
	srcPath := filepath.Join(t.TempDir(), "garbage.ail")
	raw := []byte("module world/garbage\n\nexport func broken(\n")
	if err := os.WriteFile(srcPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := canon.Source(raw); err != nil {
		t.Fatalf("control premise: canon.Source must accept these bytes: %v", err)
	}
	manifest := echoManifest(t, srcPath, false)
	res := driveTransitions(t, map[string]string{"store": storePath, "manifest": manifest, "ailang-bin": refusingInterpreter(t)},
		attendedPhrase+"\n", noEnv)
	if res.code != exitError || !strings.Contains(res.stderr, "not loadable under its pinned interpreter") {
		t.Fatalf("garbage source = (%d, %q), want exit 1 refusing the unloadable source", res.code, res.stderr)
	}
	// BEFORE PutObject: the canonical bytes were never stored.
	db, err := store.Open(storePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	source, err := canon.Source(raw)
	if err != nil {
		t.Fatal(err)
	}
	if _, found, err := db.GetObject(context.Background(), hashref.SumSHA256(source)); err != nil {
		t.Fatal(err)
	} else if found {
		t.Fatal("a refused source must not have been stored (the check runs BEFORE PutObject)")
	}
	if _, _, ok, _ := transitionreg.NewReader(db).CurrentRevision(context.Background()); ok {
		t.Fatal("a refused publish must leave no head")
	}
}

// TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter pins
// the honest scope measured in the design doc (V26b): a real module whose
// imports cannot resolve in the hermetic check root (a copy of the shape of
// world/transitions.ail, which imports world/types) is refused by the REAL
// pinned interpreter with its LDR001 output. Gated on AILANG_BIN like its
// siblings: the pinned released binary is the only interpreter that produces
// honest LDR001 output — never skip.
func TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter(t *testing.T) {
	bin := os.Getenv("AILANG_BIN")
	if bin == "" {
		t.Fatal("AILANG_BIN unset: the import-bearing source check needs the pinned released interpreter; never skip")
	}
	storePath := filepath.Join(t.TempDir(), "world.db")
	srcPath := filepath.Join(t.TempDir(), "transitions-copy.ail")
	if err := os.WriteFile(srcPath, []byte("module world/transitions_copy\n\nimport world/types (World)\n\nexport func apply(w: World) -> World {\n  w\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Bootstrap the epoch registry with the pinned binary's own release line.
	db, err := store.Open(storePath)
	if err != nil {
		t.Fatal(err)
	}
	release := firstVersionLine(t, bin)
	if _, _, err := registry.Bootstrap(context.Background(), db, release); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	manifest := echoManifest(t, srcPath, false)
	res := driveTransitions(t, map[string]string{"store": storePath, "manifest": manifest, "ailang-bin": bin},
		attendedPhrase+"\n", noEnv)
	if res.code != exitError || !strings.Contains(res.stderr, "not loadable under its pinned interpreter") {
		t.Fatalf("import-bearing source = (%d, %q), want exit 1 refusing the unloadable source", res.code, res.stderr)
	}
	if !strings.Contains(res.stderr, "LDR001") {
		t.Fatalf("refusal stderr = %q, want the interpreter's own LDR001 module-not-found output", res.stderr)
	}
}

// firstVersionLine reads the release line of an interpreter binary the same
// way the archive's manifest will (probeVersion) and the daemon's bootstrap
// reduction does (first non-blank line, trimmed).
func firstVersionLine(t *testing.T, bin string) string {
	t.Helper()
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		t.Fatalf("%s --version: %v", bin, err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			return trimmed
		}
	}
	t.Fatalf("%s --version printed no release line", bin)
	return ""
}

// TestPublishErrorLineNamesConflictingID pins the CLI rendering contract for
// the same-ID CAS-retry conflict: the line names the ID and says nothing was
// written (kimi's secondary).
func TestPublishErrorLineNamesConflictingID(t *testing.T) {
	err := &transitionreg.SameIDConflictError{ID: "tools.echo", Revision: 3}
	line := publishSetErrorLine(err)
	if !strings.Contains(line, `"tools.echo"`) {
		t.Fatalf("error line = %q, want the conflicting ID named", line)
	}
	if !strings.Contains(line, "nothing was written") {
		t.Fatalf("error line = %q, want the no-write statement", line)
	}
	var typed *transitionreg.SameIDConflictError
	if !errors.As(err, &typed) || typed.ID != "tools.echo" {
		t.Fatalf("the rendered error must stay detectably typed, got %T", err)
	}
}

// The authority gates, in order: --live is refused (this verb is local, not
// the network package publish); CI environments STOP; the typed phrase must
// match; the interpreter pin must be one of --ailang-bin/--interpreter-ref.
func TestTransitionsVerbFences(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "world.db")
	srcPath := filepath.Join(t.TempDir(), "echo.ail")
	if err := os.WriteFile(srcPath, []byte("module world/echo\n\nexport func apply(w: World) -> World ! {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest := echoManifest(t, srcPath, false)
	base := map[string]string{"store": storePath, "manifest": manifest, "ailang-bin": fakeInterpreter(t)}

	t.Run("live refused", func(t *testing.T) {
		res := drive(t, invocation{verb: "transitions", flags: base, bools: []string{"live"}}, attendedPhrase+"\n", noEnv, satisfiedProbe(t))
		if res.code != exitStop || !strings.Contains(res.stderr, "--live") {
			t.Fatalf("live = (%d, %q), want STOP naming --live (the network-publish flag)", res.code, res.stderr)
		}
	})
	t.Run("ci environment stops", func(t *testing.T) {
		res := driveTransitions(t, base, attendedPhrase+"\n", ciEnv)
		if res.code != exitStop {
			t.Fatalf("ci = (%d, %q), want STOP (never run in CI)", res.code, res.stderr)
		}
	})
	t.Run("wrong phrase stops", func(t *testing.T) {
		res := driveTransitions(t, base, "yes\n", noEnv)
		if res.code != exitStop {
			t.Fatalf("phrase = (%d, %q), want STOP on a mistyped confirmation", res.code, res.stderr)
		}
	})
	t.Run("interpreter pin required", func(t *testing.T) {
		noPin := map[string]string{"store": storePath, "manifest": manifest}
		res := driveTransitions(t, noPin, attendedPhrase+"\n", noEnv)
		if res.code != exitError || !strings.Contains(res.stderr, "interpreter pin is required") {
			t.Fatalf("pin = (%d, %q), want exit 1 naming the missing interpreter pin", res.code, res.stderr)
		}
	})
	t.Run("absent manifest is usage", func(t *testing.T) {
		res := driveTransitions(t, map[string]string{"store": storePath}, attendedPhrase+"\n", noEnv)
		if res.code != exitUsage {
			t.Fatalf("manifest = (%d, %q), want usage exit 2", res.code, res.stderr)
		}
	})
	t.Run("unverifiable interpreter ref refused", func(t *testing.T) {
		pinned := map[string]string{"store": storePath, "manifest": manifest,
			"interpreter-ref": "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
		res := driveTransitions(t, pinned, attendedPhrase+"\n", noEnv)
		if res.code != exitError || !strings.Contains(res.stderr, "is not archived next to the store") {
			t.Fatalf("unverifiable ref = (%d, %q), want exit 1 refusing the unarchived interpreter", res.code, res.stderr)
		}
	})
}

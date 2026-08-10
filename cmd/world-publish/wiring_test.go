package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/pkgproj"
)

// ---------------------------------------------------------------------------
// AC22 — the fence DOMINATES the irreversible path
//
// Every AC21 row proves a fence FIRES. None of them proves the fence is
// REACHED BEFORE the thing it guards, which is what this AST walk asserts:
// MUT-D0-FENCE-ORDER swaps the fence call with the production-constructor
// call, the mutant BUILDS (rc=0), and this test reds naming the ordering.
//
// CORRECTED (controller iter-67, measured). This comment previously claimed
// "every single AC21 row still passes" under that mutant. It does not — 6 of
// 15 red (R-CI, R-TTY-OPEN, R-TTY-CHARDEV, R-TTY-SAMEFILE, R-PHRASE-EOF,
// R-PHRASE), because those rows' fixture supplies a LOOPBACK registry origin,
// so a hoisted constructor refuses on loopback/ambient-credential first and
// the row sees "STOP fence=handler" rather than its documented line. The
// ARGUMENT for this test survives intact — an AC21 row passing proves only
// that a fence fires, never that it fires FIRST, and against a fixture with a
// valid https origin every row WOULD pass while the dominance defect shipped.
// What does not survive is "AC22 is the only killer". Recorded rather than
// quietly fixed: a non-vacuity claim asserted without being run as literally
// described is the class this mission exists to close.
//
// THE PARSER TRAP, recorded because it has cost this repository iterations:
// go/parser tests `src != nil` on the INTERFACE, so a typed-nil []byte is a
// NON-nil interface and parses as an EMPTY source — and a checker that reads
// nothing finds nothing. Every parse below therefore asserts a non-zero Decls
// count before ANY negative result is believed.
// ---------------------------------------------------------------------------

func parseCommandFile(t *testing.T, rel string) (*token.FileSet, *ast.File) {
	t.Helper()
	path := filepath.Join(commandRepoRoot(t), rel)
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	if len(src) == 0 {
		t.Fatalf("instrument failure: %s is zero bytes", rel)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", rel, err)
	}
	// ANTI-VACUITY (the typed-nil trap): a source that parsed as empty would
	// make every "we found no X" assertion below trivially true.
	if len(file.Decls) == 0 {
		t.Fatalf("instrument failure: %s parsed with ZERO declarations", rel)
	}
	return fset, file
}

// commandSourceFiles enumerates the PRODUCTION sources of this package. Test
// files are excluded deliberately: they are allowed to name loopback hosts (the
// AC21 table does), and they are not part of the shipped surface.
func commandSourceFiles(t *testing.T) []string {
	t.Helper()
	dir := filepath.Join(commandRepoRoot(t), "cmd", "world-publish")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read cmd/world-publish: %v", err)
	}
	var files []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		files = append(files, filepath.Join("cmd", "world-publish", name))
	}
	sort.Strings(files)
	if len(files) == 0 {
		t.Fatal("instrument failure: enumerated ZERO production sources in cmd/world-publish")
	}
	return files
}

// findFunc returns the named top-level function declaration.
func findFunc(t *testing.T, file *ast.File, name string) *ast.FuncDecl {
	t.Helper()
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Recv == nil && fn.Name.Name == name {
			return fn
		}
	}
	t.Fatalf("cmd/world-publish/main.go declares no func %s", name)
	return nil
}

// callSitesByStatement returns, for each top-level statement index of fn's
// body, whether that statement contains a call matching pred.
func callSitesByStatement(fn *ast.FuncDecl, pred func(*ast.CallExpr) bool) []int {
	var indices []int
	for i, stmt := range fn.Body.List {
		found := false
		ast.Inspect(stmt, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if ok && pred(call) {
				found = true
			}
			return true
		})
		if found {
			indices = append(indices, i)
		}
	}
	return indices
}

func isPlainCall(name string) func(*ast.CallExpr) bool {
	return func(call *ast.CallExpr) bool {
		ident, ok := call.Fun.(*ast.Ident)
		return ok && ident.Name == name
	}
}

func isSelectorCall(pkg, name string) func(*ast.CallExpr) bool {
	return func(call *ast.CallExpr) bool {
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != name {
			return false
		}
		ident, ok := sel.X.(*ast.Ident)
		return ok && ident.Name == pkg
	}
}

func TestTheFenceDominatesTheProductionHandlerConstruction(t *testing.T) {
	_, file := parseCommandFile(t, filepath.Join("cmd", "world-publish", "main.go"))
	fn := findFunc(t, file, "runPublish")

	if len(fn.Body.List) == 0 {
		t.Fatal("instrument failure: runPublish has an empty body")
	}

	fenceAt := callSitesByStatement(fn, isPlainCall("requireAttendedOperator"))
	ctorAt := callSitesByStatement(fn, isSelectorCall("broker", "NewRegistryPublishHandler"))

	t.Logf("runPublish has %d top-level statements; requireAttendedOperator at %v; "+
		"broker.NewRegistryPublishHandler at %v", len(fn.Body.List), fenceAt, ctorAt)

	// EXACTLY ONE of each. Two fence calls would let one be deleted silently;
	// two constructor calls would let one escape the ordering claim entirely.
	if len(fenceAt) != 1 {
		t.Fatalf("runPublish contains %d statements calling requireAttendedOperator, want exactly 1", len(fenceAt))
	}
	if len(ctorAt) != 1 {
		t.Fatalf("runPublish contains %d statements calling broker.NewRegistryPublishHandler, want exactly 1",
			len(ctorAt))
	}
	// THE ORDERING CLAIM.
	if fenceAt[0] >= ctorAt[0] {
		t.Fatalf("the attended fence is at statement %d and the production publish handler is "+
			"constructed at statement %d: the fence does NOT dominate the irreversible path",
			fenceAt[0], ctorAt[0])
	}

	// KNOWN-POSITIVE CONTROL for the instrument itself: the same walk applied to
	// a call that IS present must find it, and to one that is not must not.
	// Without this, "we found the fence at index 5" is indistinguishable from an
	// inspector that reports whatever it is asked about.
	if got := callSitesByStatement(fn, isPlainCall("requireExactlyOneMode")); len(got) != 1 {
		t.Fatalf("instrument control: requireExactlyOneMode found at %v, want exactly one statement", got)
	}
	if got := callSitesByStatement(fn, isPlainCall("thisFunctionDoesNotExist")); len(got) != 0 {
		t.Fatalf("instrument control: a nonexistent call was 'found' at %v", got)
	}
}

// ---------------------------------------------------------------------------
// AC24 — the production door is the ONLY door
// ---------------------------------------------------------------------------

var loopbackIdentifier = regexp.MustCompile(`(?i)loopback`)

// loopbackHostLiteral matches the host forms broker.isLoopbackHost recognises.
var loopbackHostLiteral = regexp.MustCompile(`127\.|localhost|::1`)

func TestTheCommandNamesOnlyTheProductionConstructor(t *testing.T) {
	files := commandSourceFiles(t)
	t.Logf("AC24(a) scanning %d production source(s): %v", len(files), files)

	productionCalls := 0
	loopbackIdents := []string{}
	loopbackLiterals := []string{}
	identifiersSeen := 0

	for _, rel := range files {
		_, file := parseCommandFile(t, rel)
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.CallExpr:
				if isSelectorCall("broker", "NewRegistryPublishHandler")(node) {
					productionCalls++
				}
			case *ast.Ident:
				identifiersSeen++
				if loopbackIdentifier.MatchString(node.Name) {
					loopbackIdents = append(loopbackIdents, rel+": "+node.Name)
				}
			case *ast.SelectorExpr:
				if loopbackIdentifier.MatchString(node.Sel.Name) {
					loopbackIdents = append(loopbackIdents, rel+": ."+node.Sel.Name)
				}
			case *ast.BasicLit:
				if node.Kind != token.STRING {
					return true
				}
				value, err := strconv.Unquote(node.Value)
				if err != nil {
					return true
				}
				if loopbackHostLiteral.MatchString(value) {
					loopbackLiterals = append(loopbackLiterals, rel+": "+strconv.Quote(value))
				}
			}
			return true
		})
	}

	// ANTI-VACUITY, BOTH HALVES. A scan that walked nothing would report zero
	// loopback identifiers AND zero production calls; requiring the production
	// call to be FOUND is what makes the two zeros above believable.
	if identifiersSeen == 0 {
		t.Fatal("instrument failure: the AST walk saw ZERO identifiers")
	}
	if productionCalls != 1 {
		t.Fatalf("cmd/world-publish names broker.NewRegistryPublishHandler %d time(s), want exactly 1. "+
			"A zero here also means the loopback scan below proves nothing", productionCalls)
	}
	t.Logf("AC24(a): %d identifiers scanned; broker.NewRegistryPublishHandler named exactly %d time(s)",
		identifiersSeen, productionCalls)

	if len(loopbackIdents) != 0 {
		t.Errorf("cmd/world-publish names loopback identifier(s): %v. The test-only door "+
			"newLoopbackRegistryPublishHandler is UNEXPORTED and must stay unreachable from here",
			loopbackIdents)
	}
	if len(loopbackLiterals) != 0 {
		t.Errorf("cmd/world-publish carries loopback host literal(s) in production source: %v",
			loopbackLiterals)
	}

	// The instrument must be able to MATCH: a scanner that never matches would
	// report the same two empty slices.
	if !loopbackIdentifier.MatchString("newLoopbackRegistryPublishHandler") ||
		!loopbackHostLiteral.MatchString("http://127.0.0.1:1") ||
		!loopbackHostLiteral.MatchString("http://localhost:8080") ||
		!loopbackHostLiteral.MatchString("http://[::1]:80") {
		t.Fatal("instrument failure: the loopback matchers do not match known-positive inputs")
	}
}

// TestTheFlagSurfaceIsFrozen is AC24(b). It is EXACT SET EQUALITY, not a subset
// check: the threat is an ADDED flag (a --validator-origin someone reaches for
// while making a test easier), and a subset check would not see one.
func TestTheFlagSurfaceIsFrozen(t *testing.T) {
	var opts options
	fs := newFlagSet(&opts)

	var declared []string
	fs.VisitAll(func(f *flag.Flag) { declared = append(declared, f.Name) })
	sort.Strings(declared)

	frozen := append([]string(nil), flagNames...)
	sort.Strings(frozen)

	if len(declared) == 0 {
		t.Fatal("instrument failure: the FlagSet enumerated ZERO flags")
	}
	if strings.Join(declared, ",") != strings.Join(frozen, ",") {
		t.Fatalf("flag surface drifted:\n declared %v\n frozen   %v", declared, frozen)
	}
	t.Logf("AC24(b): %d flags, frozen and exact: %v", len(declared), declared)

	// The specific flag that must never exist, named so it is a known target
	// rather than a discovery. broker.ApprovedValidatorOrigin is a COMPILED
	// CONSTANT and there is no path from the command line to it.
	for _, name := range declared {
		if strings.Contains(name, "validator") {
			t.Fatalf("flag --%s reaches the validator origin; ApprovedValidatorOrigin is unwidenable "+
				"by design and no flag may widen it", name)
		}
	}
	if fs.Lookup("validator-origin") != nil {
		t.Fatal("--validator-origin exists")
	}
}

// TestTheEntrypointReachesTheProductionRefusal is AC24(c): with EVERY fence
// satisfied by an injected probe and the exact typed phrase, a loopback
// --registry-origin is refused BY THE PRODUCTION CONSTRUCTOR, with the landed
// message.
//
// STATED LIMITATION, recorded rather than dressed up: the plan asks this arm to
// also assert handler.Dispatches() == 0 and handler.CredentialLoads() == 0.
// That is NOT OBSERVABLE HERE and cannot be made so honestly — the constructor
// REFUSES, so it returns no handler whose counters could be read. The landed
// TestPublishDenialsPersistAndNeverDispatchWithLivePositiveControl already
// measures those counters at the handler. What THIS arm adds, and what did not
// exist before this milestone, is that the ENTRYPOINT reaches the refusal at
// all.
func TestTheEntrypointReachesTheProductionRefusal(t *testing.T) {
	// MEASURED WHILE BUILDING THIS MILESTONE, and worth recording: on the rig
	// this loop runs on, AILANG_REGISTRY_API_KEY IS AMBIENT. The first version
	// of this arm reached broker.NewRegistryPublishHandler and was refused by
	// AssertNoAmbientRegistryCredential — which is checked BEFORE the origin —
	// so it never saw the loopback message at all. A third independent barrier
	// is therefore live in this environment today.
	//
	// That also makes the arm ENVIRONMENT-DEPENDENT: CI has the variable unset
	// and would take the loopback branch, the rig takes the ambient branch, and
	// a test whose observed refusal depends on the machine is not a measurement.
	// Scrubbing it to EMPTY for the duration of this test is the repository's
	// own landed idiom (registry_publish_test.go:1695 and :1703) and is what
	// childenv.Has treats as absent. It never reads the value, never exports one
	// and never supplies a credential: it strictly REMOVES authority from this
	// test process.
	t.Setenv(broker.RegistryCredentialVariable, "")

	inv := liveInvocation(t)
	got := drive(t, inv, attendedPhrase+"\n", noEnv, satisfiedProbe(t))

	if got.code != exitStop {
		t.Fatalf("exit code = %d, want %d\nstdout:\n%s\nstderr:\n%s",
			got.code, exitStop, got.stdout, got.stderr)
	}
	line := stopLine(got.stderr)
	if line != "STOP fence=handler reason=refused" {
		t.Fatalf("STOP line = %q, want the handler stage\nstderr:\n%s", line, got.stderr)
	}
	// The LANDED refusal text, verbatim from host/broker/registry_publish.go.
	const landed = "the production constructor refuses a loopback registry origin"
	if !strings.Contains(got.stderr, landed) {
		t.Fatalf("the refusal does not carry the landed loopback message %q\nstderr:\n%s",
			landed, got.stderr)
	}
	// It must have got PAST the fence, or this arm proves nothing about the
	// constructor: the prompt is written only when the confirmation is read.
	if !strings.Contains(got.stdout, attendedPhrase) {
		t.Fatalf("the confirmation prompt was never printed, so the fence stack did not complete "+
			"and the refusal above is not attributable to the constructor\nstdout:\n%s", got.stdout)
	}
	t.Logf("AC24(c): %s\n  %s", line, landed)

	// The same refusal, reached directly, so the message this arm greps for is
	// demonstrably the CONSTRUCTOR's and not a string this package holds.
	_, err := broker.NewRegistryPublishHandler(broker.RegistryPublishConfig{
		PublisherPath:   "/usr/bin/true",
		PackageDir:      t.TempDir(),
		RegistryOrigin:  "https://127.0.0.1:1",
		ValidatorOrigin: broker.ApprovedValidatorOrigin,
	})
	if err == nil || !strings.Contains(err.Error(), landed) {
		t.Fatalf("instrument failure: the landed constructor did not produce %q; got %v", landed, err)
	}
}

// TestCommandNeverNamesTheCredentialVariable is the static half of the
// ambient-credential guarantee.
//
// DEVIATION, STATED: the plan's AC24 arm (b) asks for a child process with
// AILANG_REGISTRY_API_KEY set to a sentinel. The executing directive for this
// milestone forbids setting, reading, exporting or passing that variable under
// any circumstances, so that arm is NOT performed here. The landed
// TestProductionConstructorRefusesAnAmbientCredential already drives the
// refusal at the constructor, and broker.AssertNoAmbientRegistryCredential is
// called INSIDE broker.NewRegistryPublishHandler, which AC22 proves this
// command reaches. What is added instead is the stronger STATIC claim: this
// command cannot read, set or name the variable at all.
func TestCommandNeverNamesTheCredentialVariable(t *testing.T) {
	files := commandSourceFiles(t)
	needles := []string{
		broker.RegistryCredentialVariable,
		"RegistryCredentialVariable",
	}
	matches := []string{}
	scanned := 0
	for _, rel := range files {
		data, err := os.ReadFile(filepath.Join(commandRepoRoot(t), rel))
		if err != nil {
			t.Fatal(err)
		}
		scanned += len(data)
		// The scan is over the AST's identifiers, selectors and string literals
		// rather than raw bytes: the doc comment on `environment` NAMES the
		// variable to explain why it is never read, and a byte grep would call
		// that documentation a violation.
		_, file := parseCommandFile(t, rel)
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.Ident:
				for _, needle := range needles {
					if node.Name == needle {
						matches = append(matches, rel+": ident "+node.Name)
					}
				}
			case *ast.SelectorExpr:
				for _, needle := range needles {
					if node.Sel.Name == needle {
						matches = append(matches, rel+": selector ."+node.Sel.Name)
					}
				}
			case *ast.BasicLit:
				if node.Kind != token.STRING {
					return true
				}
				value, err := strconv.Unquote(node.Value)
				if err != nil {
					return true
				}
				for _, needle := range needles {
					if strings.Contains(value, needle) {
						matches = append(matches, rel+": literal "+strconv.Quote(value))
					}
				}
			}
			return true
		})
	}
	if scanned == 0 {
		t.Fatal("instrument failure: scanned ZERO bytes")
	}
	if len(matches) != 0 {
		t.Fatalf("cmd/world-publish names the registry credential variable in CODE: %v", matches)
	}
	// KNOWN-POSITIVE CONTROL: the same needle, in the same shape, IS found when
	// it is present. Without it, "zero matches" is also what a dead scan reports.
	control := "package p\nvar x = \"" + broker.RegistryCredentialVariable + "\"\n"
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, "control.go", control, 0)
	if err != nil {
		t.Fatal(err)
	}
	found := 0
	ast.Inspect(parsed, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING {
			if value, err := strconv.Unquote(lit.Value); err == nil &&
				strings.Contains(value, broker.RegistryCredentialVariable) {
				found++
			}
		}
		return true
	})
	if found != 1 {
		t.Fatalf("instrument failure: the control source contains the needle once but the scan found %d", found)
	}
	t.Logf("AC24: %d bytes of production source scanned, ZERO references to the credential variable "+
		"(control fired: %d)", scanned, found)
}

// TestWorldCoreManifestMatchesTheCommittedGolden binds the Go LITERAL manifest
// in fences.go to the real artifact.
//
// This repository has no TOML dependency, so packages/world-core/ailang.toml is
// restated in Go. That restatement is the obvious place for silent rot — so it
// is measured against the golden the SHELL gate produced FROM THE REAL FILE, by
// recomputing the whole packet from the real projection directory. If the
// literal's exports, edition or version ever drift from the toml, the interface
// hash or a scalar field moves and this reds.
func TestWorldCoreManifestMatchesTheCommittedGolden(t *testing.T) {
	root := commandRepoRoot(t)
	golden, err := pkgproj.LoadReadyPacket(filepath.Join(root, defaultGolden))
	if err != nil {
		t.Fatalf("load committed golden: %v", err)
	}
	recomputed, err := pkgproj.RecomputeReadyPacket(
		filepath.Join(root, defaultPackageDir), worldCoreManifest, frozenCompilerVersion)
	if err != nil {
		t.Fatalf("recompute the packet from the real projection: %v", err)
	}
	if field, equal := recomputed.Equal(golden); !equal {
		t.Fatalf("the Go manifest literal disagrees with the committed golden at %q: "+
			"recomputed %q, golden %q", field, recomputed.Field(field), golden.Field(field))
	}
	t.Logf("worldCoreManifest reproduces the committed golden in all %d fields",
		len(pkgproj.ReadyPacketFields))

	// The toml on disk must still name what the literal claims, or the two
	// sources agree only because the golden is stale too.
	toml, err := os.ReadFile(filepath.Join(root, defaultPackageDir, "ailang.toml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`name = "` + worldCoreManifest.Package.Name + `"`,
		`version = "` + worldCoreManifest.Package.Version + `"`,
		`edition = "` + worldCoreManifest.Package.Edition + `"`,
		`ailang = "` + worldCoreManifest.Package.AILANG + `"`,
	} {
		if !bytes.Contains(toml, []byte(want)) {
			t.Errorf("packages/world-core/ailang.toml does not contain %q", want)
		}
	}
	for _, module := range worldCoreManifest.Exports.Modules {
		if !bytes.Contains(toml, []byte(`"`+module+`"`)) {
			t.Errorf("packages/world-core/ailang.toml does not export %q", module)
		}
	}
}

// TestReconcileConfigIsBuiltFromTheReadyPacket covers the one production line
// AC26's loopback arm cannot reach from here: the arguments handed to the
// PRODUCTION reconciler. It is pure, so it is asserted without issuing a GET.
//
// STATED LIMITATION: the `--probe` branch's CALL is not executed by any test in
// this repository, because executing it would mean a non-loopback request to the
// public bucket. Its argument construction is measured here; its behaviour is
// measured in host/broker/registry_reconcile_test.go against loopback fakes.
func TestReconcileConfigIsBuiltFromTheReadyPacket(t *testing.T) {
	root := commandRepoRoot(t)
	packet, err := pkgproj.LoadReadyPacket(filepath.Join(root, defaultGolden))
	if err != nil {
		t.Fatal(err)
	}
	cfg := reconcileConfigFor(options{registryOrigin: "https://storage.googleapis.com/ailang-registry"}, packet)

	if cfg.Vendor != "world" || cfg.Name != "core" || cfg.Version != frozenPackageVersion {
		t.Fatalf("reconcile target = %s/%s@%s, want world/core@%s",
			cfg.Vendor, cfg.Name, cfg.Version, frozenPackageVersion)
	}
	if cfg.Expected.TarballSHA256 != packet.TarballSHA256 ||
		cfg.Expected.ContentHash != packet.ContentHash ||
		cfg.Expected.InterfaceHash != packet.InterfaceHash {
		t.Fatalf("reconcile expectations do not come from the ready packet: %+v", cfg.Expected)
	}
	// Empty expectations would make a present document compare equal to
	// anything, which is how `conflict` silently becomes `succeeded-reconciled`.
	if cfg.Expected.TarballSHA256 == "" || cfg.Expected.ContentHash == "" || cfg.Expected.InterfaceHash == "" {
		t.Fatal("reconcile expectations carry an empty digest")
	}
	// There is deliberately no caller-chosen control: host/broker builds the
	// target and its same-pass known-positive control from ONE origin.
	if cfg.ControlVendor != "" || cfg.ControlName != "" || cfg.ControlVersion != "" {
		t.Fatalf("the command chose a probe control (%s/%s@%s); the shipped default must be used "+
			"so the control always travels the target's own key-space",
			cfg.ControlVendor, cfg.ControlName, cfg.ControlVersion)
	}
}

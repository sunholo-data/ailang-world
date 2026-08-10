package runbook

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

// ---------------------------------------------------------------------------
// SM.D0 — Stage B is now a sequence of commands, so it can be GATED
//
// Until this milestone Stage B was prose. That is not a stylistic observation:
// it is the reason TestRunbookStageAPerformsNoPublicWrite has been STRUCTURALLY
// VACUOUS since the day it landed. Measured at 6d1dce0,
// `grep -c 'ailang publish' docs/SELF_MOD_PUBLISH.md` = 0 — so its detection
// loop has NEVER EXECUTED ITS BODY. Its two instrument-failure fatals check the
// REGION being scanned; nothing checked the PREDICATE doing the scanning, and a
// green was therefore indistinguishable from a matcher that can never match.
//
// AC27's repair is structural: ONE predicate, driven in TWO directions. The
// same function that must find ZERO live publishes in Stage A must find AT
// LEAST ONE in Stage B. A broken matcher now yields a Stage-B count of zero and
// REDS. A broken instrument can no longer produce a green.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// extraction
// ---------------------------------------------------------------------------

// fencedCommand is one command line lifted OUT of a ```bash block, with the
// stage it was found in.
type fencedCommand struct {
	stage string // "A" or "B"
	line  int
	text  string
}

const (
	stageAMarker = "## Stage A"
	stageBMarker = "## Stage B"
)

// extractFencedCommands walks the document once and returns every command line
// inside a ```bash fence, tagged with its stage.
//
// It scans FENCED lines rather than all lines on purpose. The predicate below
// is a claim about COMMANDS, and prose in this runbook legitimately discusses
// `--live` (the fence table says the reconcile verb refuses it). A line-wide
// scan would call that documentation a live publish, and the usual repair —
// carving out exceptions — is how a matcher stops matching.
func extractFencedCommands(t *testing.T, doc string) []fencedCommand {
	t.Helper()
	var (
		commands []fencedCommand
		stage    string
		inFence  bool
	)
	for i, raw := range strings.Split(doc, "\n") {
		line := strings.TrimRight(raw, " \t")
		switch {
		case strings.HasPrefix(line, stageAMarker):
			stage = "A"
			continue
		case strings.HasPrefix(line, stageBMarker):
			stage = "B"
			continue
		}
		if strings.HasPrefix(line, "```bash") {
			inFence = true
			continue
		}
		if strings.HasPrefix(line, "```") {
			inFence = false
			continue
		}
		if !inFence || stage == "" || strings.TrimSpace(line) == "" {
			continue
		}
		commands = append(commands, fencedCommand{stage: stage, line: i + 1, text: strings.TrimSpace(line)})
	}
	// ANTI-VACUITY: a document from which no commands were extracted would make
	// every count below zero, and "Stage A contains no live publish" would be
	// true of a blank file.
	if len(commands) == 0 {
		t.Fatalf("instrument failure: extracted ZERO fenced command lines from %s (%d bytes)",
			runbookPath, len(doc))
	}
	return commands
}

// livePublishCommand is THE PREDICATE. There is exactly one, and both
// directions of AC27 are stated against it.
//
// Two forms perform a public write:
//
//	`ailang publish` without --dry-run — the pinned compiler's own publisher,
//	    which the readiness gate genuinely runs in its --dry-run form, so the
//	    exemption is the PRESENCE of --dry-run;
//	a world-publish invocation carrying --live — where the exemption is the
//	    ABSENCE of --live, because --live is the flag that selects the
//	    irreversible path and nothing else in this repository accepts it.
//
// The second form is matched on `--live` rather than on the program name
// because the runbook invokes the built binary through "$WORLD_BIN": a name
// match would silently stop matching the moment the document did the sensible
// thing, which is exactly the rot this predicate exists to detect.
func livePublishCommand(line string) bool {
	if strings.Contains(line, "ailang publish") && !strings.Contains(line, "--dry-run") {
		return true
	}
	return strings.Contains(line, "--live")
}

// TestLivePublishPredicateHasBothDirections drives the predicate itself before
// any count is believed. A predicate that always returned false would make
// AC27's Stage-A half green and its Stage-B half red; a predicate that always
// returned true would do the reverse. Neither can survive this.
func TestLivePublishPredicateHasBothDirections(t *testing.T) {
	for _, positive := range []string{
		`ailang publish`,
		`ailang publish --allow-dotted-tool-names`,
		`"$WORLD_BIN" publish --live --store "$WORLD_STORE"`,
		`./world-publish publish --live`,
	} {
		if !livePublishCommand(positive) {
			t.Errorf("livePublishCommand(%q) = false, want true", positive)
		}
	}
	for _, negative := range []string{
		`ailang publish --dry-run`,
		`./scripts/verify_world_package.sh`,
		`go run ./cmd/world-publish packet`,
		`"$WORLD_BIN" publish --dry-run --store "$WORLD_STORE"`,
		`"$WORLD_BIN" reconcile --store "$WORLD_STORE"`,
		`go build -o "$WORLD_BIN" ./cmd/world-publish`,
	} {
		if livePublishCommand(negative) {
			t.Errorf("livePublishCommand(%q) = true, want false", negative)
		}
	}
}

// TestOnePredicateSeesZeroLivePublishesInStageAAndAtLeastOneInStageB is AC27.
func TestOnePredicateSeesZeroLivePublishesInStageAAndAtLeastOneInStageB(t *testing.T) {
	doc := readRunbook(t)
	commands := extractFencedCommands(t, doc)

	fencesA, fencesB := countStageFences(t, doc)
	t.Logf("AC27 fenced bash blocks: Stage A = %d, Stage B = %d", fencesA, fencesB)
	// Stage A's three blocks are the projection, the readiness gate and the repo
	// gate. An exact count here is what stops a block being deleted to make the
	// zero below easier to achieve.
	if fencesA != 3 {
		t.Fatalf("Stage A carries %d fenced bash blocks, want exactly 3", fencesA)
	}
	if fencesB < 1 {
		t.Fatalf("Stage B carries %d fenced bash blocks: it is prose again, and the attended "+
			"procedure is unexecutable", fencesB)
	}

	var inA, inB []fencedCommand
	for _, cmd := range commands {
		if !livePublishCommand(cmd.text) {
			continue
		}
		if cmd.stage == "A" {
			inA = append(inA, cmd)
		} else {
			inB = append(inB, cmd)
		}
	}

	// DIRECTION 1: the automated stage performs no public write.
	for _, cmd := range inA {
		t.Errorf("%s:%d Stage A (the UNATTENDED stage) contains a live publish: %q. "+
			"Stage A must stop at readiness", runbookPath, cmd.line, cmd.text)
	}

	// DIRECTION 2 — THE REPAIR. The SAME predicate must produce a positive here.
	// Without this half, a matcher that can never match yields a green above and
	// the gate measures nothing at all. That was literally the landed state.
	if len(inB) == 0 {
		t.Fatalf("instrument failure: the SAME predicate that reported zero live publishes in "+
			"Stage A also reports zero in Stage B, across %d extracted command lines. Either the "+
			"attended stage stopped declaring the publish it exists to perform, or "+
			"livePublishCommand no longer matches this document's shape — and in the second case "+
			"the Stage-A zero above proves NOTHING", len(commands))
	}
	t.Logf("AC27: %d fenced command line(s); live publishes: Stage A = %d, Stage B = %d (%q)",
		len(commands), len(inA), len(inB), inB[0].text)
}

func countStageFences(t *testing.T, doc string) (int, int) {
	t.Helper()
	var stage string
	counts := map[string]int{}
	for _, line := range strings.Split(doc, "\n") {
		switch {
		case strings.HasPrefix(line, stageAMarker):
			stage = "A"
		case strings.HasPrefix(line, stageBMarker):
			stage = "B"
		case strings.HasPrefix(line, "```bash") && stage != "":
			counts[stage]++
		}
	}
	return counts["A"], counts["B"]
}

func readRunbook(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), runbookPath))
	if err != nil {
		t.Fatalf("attended publish runbook is unreadable at %s: %v", runbookPath, err)
	}
	if len(data) == 0 {
		t.Fatalf("instrument failure: %s is zero bytes", runbookPath)
	}
	return string(data)
}

// ---------------------------------------------------------------------------
// AC28 — the document's digests ARE the artifact's digests
// ---------------------------------------------------------------------------

var fullDigest = regexp.MustCompile(`sha256:[0-9a-f]{64}`)

const goldenPath = "scripts/world_package_ready_packet.golden.json"

// TestRunbookDigestsAppearVerbatimInTheCommittedGolden is AC28, and it is rule
// 3k: the digits are taken OUT of the document and compared with the artifact,
// never re-derived. A test that recomputed them would verify its own arithmetic.
//
// The EXACT COUNT of 3 is the anti-vacuity guard, and it is not theoretical:
// before this milestone the runbook named ZERO digests (measured: 0 in the doc,
// 3 in the golden), so a "every digest in the doc is in the golden" assertion
// was true of the empty set.
func TestRunbookDigestsAppearVerbatimInTheCommittedGolden(t *testing.T) {
	root := repoRoot(t)
	doc := readRunbook(t)
	golden, err := os.ReadFile(filepath.Join(root, goldenPath))
	if err != nil {
		t.Fatalf("committed golden is unreadable at %s: %v", goldenPath, err)
	}

	inDoc := fullDigest.FindAllString(doc, -1)
	inGolden := fullDigest.FindAllString(string(golden), -1)
	t.Logf("AC28: %d full-length digest(s) in %s, %d in %s", len(inDoc), runbookPath, len(inGolden), goldenPath)

	// KNOWN-POSITIVE CONTROL, in the same call: the extractor must find the
	// three digests in the artifact. If it finds none there, a zero in the
	// document means nothing about the document.
	if len(inGolden) != 3 {
		t.Fatalf("instrument failure: the digest extractor found %d digests in %s, want 3",
			len(inGolden), goldenPath)
	}
	if len(inDoc) != 3 {
		t.Fatalf("%s carries %d full-length digests, want exactly 3 (contentHash, interfaceHash, "+
			"tarballSHA256). A runbook that names none would pass a subset check over an empty set — "+
			"which is exactly the state this criterion repaired", runbookPath, len(inDoc))
	}

	goldenText := string(golden)
	for _, digest := range inDoc {
		if !strings.Contains(goldenText, digest) {
			t.Errorf("%s names digest %s, which does not appear in %s. The runbook is telling an "+
				"attended operator to approve bytes that are not the reviewed artifact",
				runbookPath, digest, goldenPath)
		}
	}
	// And the three must be DISTINCT: the same digest repeated three times would
	// satisfy both the count and the membership check.
	distinct := map[string]bool{}
	for _, digest := range inDoc {
		distinct[digest] = true
	}
	if len(distinct) != 3 {
		t.Fatalf("%s names %d distinct digests among 3 occurrences: %v", runbookPath, len(distinct), inDoc)
	}
	// NEGATIVE CONTROL: a digest that is NOT in the golden must be detectable.
	if strings.Contains(goldenText, "sha256:"+strings.Repeat("f", 64)) {
		t.Fatal("negative control failed: the golden contains an all-f digest, so a wrong digest " +
			"would not be detectable")
	}
}

// TestRunbookStepFourStatesTheGatesRealGuarantee is AC29.
//
// DECLARED WEAK. This is a phrase-presence assertion over prose: it would pass
// identically if the underlying claim were false, because nothing about the
// gate's behaviour is measured by it. It is kept because the text had to change
// (it previously instructed a comparison against output that does not exist —
// measured: 0 digest-shaped strings of ANY length in the gate's own log) and
// something must stop it silently reverting. Its TEETH are AC28, which binds
// the document to the artifact by measurement. It is labelled here rather than
// dressed up, so no reviewer counts it as a gate.
func TestRunbookStepFourStatesTheGatesRealGuarantee(t *testing.T) {
	doc := readRunbook(t)
	for _, phrase := range []string{
		"byte-for-byte",
		"The gate prints no digests",
		"17 hex characters",
		"68 bits",
		"are not the verification",
	} {
		if !strings.Contains(doc, phrase) {
			t.Errorf("WEAK ASSERTION: %s no longer states %q", runbookPath, phrase)
		}
	}
	// The old, impossible instruction must be gone.
	if strings.Contains(doc, "Confirm all three digests against the gate's own output") {
		t.Errorf("%s still instructs a comparison against output the gate does not produce", runbookPath)
	}
}

// ---------------------------------------------------------------------------
// AC30 — CI cannot reach the publish entrypoint
// ---------------------------------------------------------------------------

// TestNoCIStepOrScriptReachesThePublishEntrypoint is AC30. An empty result is a
// CLAIM, not a fact, so the zero is paired with TWO known-positive controls in
// the SAME call that prove the scanner reads those files and can match.
func TestNoCIStepOrScriptReachesThePublishEntrypoint(t *testing.T) {
	root := repoRoot(t)

	scripts, err := filepath.Glob(filepath.Join(root, "scripts", "*.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if len(scripts) == 0 {
		t.Fatal("instrument failure: enumerated ZERO shell scripts under scripts/")
	}
	targets := append([]string{filepath.Join(root, ".github", "workflows", "ci.yml")}, scripts...)

	hits := []string{}
	bytesScanned := 0
	for _, path := range targets {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		bytesScanned += len(data)
		for i, line := range strings.Split(string(data), "\n") {
			if strings.Contains(line, "world-publish") {
				rel, _ := filepath.Rel(root, path)
				hits = append(hits, rel+":"+itoa(i+1)+": "+strings.TrimSpace(line))
			}
		}
	}
	if bytesScanned == 0 {
		t.Fatal("instrument failure: scanned ZERO bytes")
	}

	// THE TWO KNOWN-POSITIVE CONTROLS, in the same call as the zero.
	ciYML := mustRead(t, filepath.Join(root, ".github", "workflows", "ci.yml"))
	verifyAIL := mustRead(t, filepath.Join(root, "scripts", "verify_ail.sh"))
	controlCI := strings.Count(ciYML, "verify_go.sh")
	controlScript := strings.Count(verifyAIL, "verify_world_package.sh")
	t.Logf("AC30: scanned %d bytes across %d file(s); controls: verify_go.sh in ci.yml = %d, "+
		"verify_world_package.sh in verify_ail.sh = %d", bytesScanned, len(targets), controlCI, controlScript)
	if controlCI < 1 {
		t.Fatalf("known-positive control failed: ci.yml does not mention verify_go.sh, so this "+
			"scanner is not reading the workflow and its zero proves nothing (got %d)", controlCI)
	}
	if controlScript < 1 {
		t.Fatalf("known-positive control failed: verify_ail.sh does not mention "+
			"verify_world_package.sh, so this scanner is not reading the scripts (got %d)", controlScript)
	}

	if len(hits) != 0 {
		t.Fatalf("a CI step or repository script reaches the attended publish entrypoint:\n  %s\n"+
			"world-publish performs an IRREVERSIBLE public write and must never be reachable from "+
			"automation. The controlling-terminal fence would refuse it, but a CI step that TRIES "+
			"is a design error, not a near miss", strings.Join(hits, "\n  "))
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

// ---------------------------------------------------------------------------
// AC23 — RULE 3k: execute the DOCUMENT'S OWN tokens
// ---------------------------------------------------------------------------

// SAFETY OF THIS TEST, stated structurally rather than by assurance.
//
// This test runs the runbook's real live-publish argv. FIVE independent things
// make a publish impossible, and no ONE of them is relied on:
//
//  1. The controlling-terminal fence is ordered BEFORE any transport can exist
//     (asserted over the AST by AC22), and it is UNSATISFIABLE here: the child
//     gets piped stdio, so /dev/tty either will not open or is not stdin.
//  2. The child's environment is built with childenv.Scrubbed, so every
//     registry variable — including the credential — is REMOVED. This test
//     never reads, sets or passes AILANG_REGISTRY_API_KEY.
//  3. WORLD_REGISTRY is expanded to a LOOPBACK origin, which the PRODUCTION
//     constructor refuses outright.
//  4. WORLD_CREDENTIAL names a path that does not exist, and the production
//     constructor refuses to build a live handler without a usable credential.
//  5. WORLD_APPROVAL names a well-formed hashref that resolves to no object, so
//     validatePublishApproval refuses before the durable claim is taken.
//
// Any four of these could fail and the fifth would still hold.

// buildOnce compiles cmd/world-publish exactly once per test binary. Rebuilding
// per arm would make this gate slow enough that someone would eventually skip
// it, and a rule-3k arm that is skipped is a rule-3k arm that does not exist.
var buildOnce struct {
	sync.Once
	path string
	err  error
}

func publishBinary(t *testing.T) string {
	t.Helper()
	buildOnce.Do(func() {
		dir, err := os.MkdirTemp("", "world-publish-gate-*")
		if err != nil {
			buildOnce.err = err
			return
		}
		out := filepath.Join(dir, "world-publish")
		cmd := exec.Command("go", "build", "-o", out, "./cmd/world-publish")
		cmd.Dir = repoRootFromCaller()
		if combined, err := cmd.CombinedOutput(); err != nil {
			buildOnce.err = err
			buildOnce.path = string(combined)
			return
		}
		buildOnce.path = out
	})
	if buildOnce.err != nil {
		t.Fatalf("build ./cmd/world-publish: %v\n%s", buildOnce.err, buildOnce.path)
	}
	return buildOnce.path
}

func repoRootFromCaller() string {
	// repoRoot needs a *testing.T; buildOnce runs inside one, but sync.Once's
	// closure must not capture it (a later caller would then be reporting
	// through the first caller's T). Resolve independently.
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

// tokenize splits ONE runbook command line into argv, honouring double quotes
// and expanding $VAR / ${VAR} from the supplied session. It is deliberately a
// tiny, total subset of shell: the runbook's Stage B commands are written as
// single lines of plain tokens for exactly this reason, so the gate can run
// them without interpreting a shell.
func tokenize(t *testing.T, line string, session map[string]string) []string {
	t.Helper()
	var (
		argv    []string
		current strings.Builder
		quoted  bool
		started bool
	)
	flush := func() {
		if started {
			argv = append(argv, current.String())
			current.Reset()
			started = false
		}
	}
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case c == '"':
			quoted = !quoted
			started = true
		case c == ' ' && !quoted:
			flush()
		case c == '$':
			name, width := readVarName(line[i+1:])
			if name == "" {
				current.WriteByte(c)
				started = true
				continue
			}
			value, ok := session[name]
			if !ok {
				t.Fatalf("the runbook command %q references $%s, which the gate does not bind. "+
					"Either the runbook grew a variable or this gate stopped covering one", line, name)
			}
			current.WriteString(value)
			started = true
			i += width
		default:
			current.WriteByte(c)
			started = true
		}
	}
	flush()
	if len(argv) == 0 {
		t.Fatalf("instrument failure: tokenizing %q produced ZERO arguments", line)
	}
	for _, arg := range argv {
		if strings.Contains(arg, "$") {
			t.Fatalf("argument %q still carries an unexpanded variable after tokenizing %q", arg, line)
		}
	}
	return argv
}

func readVarName(rest string) (string, int) {
	if strings.HasPrefix(rest, "{") {
		end := strings.Index(rest, "}")
		if end < 0 {
			return "", 0
		}
		return rest[1:end], end + 1
	}
	end := 0
	for end < len(rest) && (rest[end] == '_' ||
		rest[end] >= 'A' && rest[end] <= 'Z' ||
		rest[end] >= 'a' && rest[end] <= 'z' ||
		end > 0 && rest[end] >= '0' && rest[end] <= '9') {
		end++
	}
	return rest[:end], end
}

// findStageBCommand returns the single Stage B fenced command matching pred.
func findStageBCommand(t *testing.T, commands []fencedCommand, what string, pred func(string) bool) fencedCommand {
	t.Helper()
	var found []fencedCommand
	for _, cmd := range commands {
		if cmd.stage == "B" && pred(cmd.text) {
			found = append(found, cmd)
		}
	}
	if len(found) != 1 {
		t.Fatalf("Stage B of %s carries %d %s command(s), want exactly 1: rule 3k needs an "+
			"unambiguous line to lift out of the document", runbookPath, len(found), what)
	}
	return found[0]
}

func TestTheRunbooksOwnStageBCommandsAreExecutedAgainstTheBuiltBinary(t *testing.T) {
	root := repoRoot(t)
	doc := readRunbook(t)
	commands := extractFencedCommands(t, doc)
	binary := publishBinary(t)

	session := map[string]string{
		"WORLD_BIN": binary,
		// A path store.Open will create. Nothing is published into it.
		"WORLD_STORE": filepath.Join(t.TempDir(), "world.db"),
		// BARRIER 3: loopback, which the production constructor refuses.
		"WORLD_REGISTRY": "https://127.0.0.1:1",
		// BARRIER 4: a path that does not exist.
		"WORLD_CREDENTIAL": filepath.Join(t.TempDir(), "no-such-credential.key"),
		// BARRIER 5: well-formed, resolves to no object.
		"WORLD_COMPILER": publisherStub(t),
		"WORLD_APPROVAL": "sha256:" + strings.Repeat("ab", 32),
		"USER":           "runbook-gate",
	}
	// MEASURED, and the reason WORLD_COMPILER is a stub rather than a missing
	// path: --publisher is READ (to record which binary published) while the
	// plan is built, which is BEFORE the fence. Pointing it at a nonexistent
	// file made the live arm exit 1 with "read --publisher: no such file",
	// never reaching the tty fence — so the arm would have been green on a
	// refusal that had nothing to do with the property it claims to measure.
	//
	// The stub is BARRIER 6, and it is stronger than absence: it is mode 0600,
	// so even a total collapse of the five barriers below could not EXECUTE it.
	if _, err := os.Stat(session["WORLD_CREDENTIAL"]); err == nil {
		t.Fatalf("safety control failed: WORLD_CREDENTIAL (%s) EXISTS", session["WORLD_CREDENTIAL"])
	}
	stubInfo, err := os.Stat(session["WORLD_COMPILER"])
	if err != nil {
		t.Fatalf("the publisher stub is unreadable: %v", err)
	}
	if stubInfo.Mode().Perm()&0o111 != 0 {
		t.Fatalf("safety control failed: the publisher stub %s is EXECUTABLE (mode %v)",
			session["WORLD_COMPILER"], stubInfo.Mode().Perm())
	}

	// --- ARM 1: the live publish command, lifted out of the document ---------
	live := findStageBCommand(t, commands, "live publish", livePublishCommand)
	liveArgv := tokenize(t, live.text, session)
	t.Logf("AC23 live argv (from %s:%d): %q", runbookPath, live.line, liveArgv)
	if liveArgv[0] != binary {
		t.Fatalf("the runbook's live command invokes %q, not the built binary", liveArgv[0])
	}
	liveRun := runArgv(t, root, liveArgv)
	t.Logf("AC23 live arm: exit %d\nstderr:\n%s", liveRun.code, liveRun.stderr)

	if liveRun.code != 3 {
		t.Fatalf("the runbook's OWN live-publish command exited %d, want 3 (STOP).\n"+
			"stdout:\n%s\nstderr:\n%s", liveRun.code, liveRun.stdout, liveRun.stderr)
	}
	if !strings.Contains(liveRun.stderr, "STOP fence=tty") {
		t.Fatalf("the live command stopped, but NOT at the controlling-terminal fence. The fence "+
			"that makes this procedure attended is the tty fence; a stop at any earlier fence means "+
			"this arm is not measuring it.\nstderr:\n%s", liveRun.stderr)
	}

	// --- ARM 2: the reconcile command, THE DISCRIMINATOR ---------------------
	//
	// Without this arm, arm 1 would pass identically if the binary refused
	// EVERYTHING — a broken build, a bad flag, a panic. Same binary, same
	// document, same piped stdio, exit 0: that is what proves the refusal is
	// attached to the LIVE path specifically rather than to the command.
	reconcile := findStageBCommand(t, commands, "reconcile", func(line string) bool {
		return strings.Contains(line, " reconcile ") && !strings.Contains(line, "--probe")
	})
	reconcileArgv := tokenize(t, reconcile.text, session)
	t.Logf("AC23 reconcile argv (from %s:%d): %q", runbookPath, reconcile.line, reconcileArgv)
	reconcileRun := runArgv(t, root, reconcileArgv)
	t.Logf("AC23 reconcile arm: exit %d\nstdout:\n%s", reconcileRun.code, reconcileRun.stdout)

	if reconcileRun.code != 0 {
		t.Fatalf("the runbook's OWN reconcile command exited %d, want 0. The read-only verb is "+
			"deliberately headless-permitted; if it cannot run under piped stdio, the live arm's "+
			"exit 3 above proves only that this binary refuses everything.\nstdout:\n%s\nstderr:\n%s",
			reconcileRun.code, reconcileRun.stdout, reconcileRun.stderr)
	}
	if strings.Contains(reconcileRun.stderr, "STOP fence=") {
		t.Fatalf("the reconcile arm exited 0 but printed a STOP line: %s", reconcileRun.stderr)
	}
}

// publisherStub is a readable, NON-EXECUTABLE stand-in for the pinned
// compiler. world-publish digests it as provenance and never runs it; the
// attended fence stops the invocation long before any subprocess exists.
func publisherStub(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "not-the-pinned-compiler")
	if err := os.WriteFile(path, []byte("this file is not a compiler and is not executable\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

type argvResult struct {
	code           int
	stdout, stderr string
}

// runArgv executes one extracted argv with PIPED stdio and a SCRUBBED
// environment. The stdio is what makes the tty fence unsatisfiable; the scrub
// is what makes a credential unavailable even if it were not.
func runArgv(t *testing.T, root string, argv []string) argvResult {
	t.Helper()
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = root
	// childenv.Scrubbed is the LANDED scrubber: it removes every registry
	// variable, including AILANG_REGISTRY_API_KEY, by name. Building the child
	// environment through it rather than by hand means this test cannot drift
	// away from the list the rest of the repository uses.
	cmd.Env = childenv.Scrubbed(os.Environ())
	// Also strip CI markers: this arm must measure the TTY fence, and the
	// DECLARED TRIPWIRE would otherwise fire first when the gate runs in CI and
	// hide whether the load-bearing layer works at all.
	cmd.Env = withoutVariables(cmd.Env, "CI", "GITHUB_ACTIONS")
	cmd.Stdin = strings.NewReader("")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("run %q: %v", argv, err)
		}
		code = exit.ExitCode()
	}
	return argvResult{code: code, stdout: stdout.String(), stderr: stderr.String()}
}

func withoutVariables(environ []string, names ...string) []string {
	drop := map[string]bool{}
	for _, name := range names {
		drop[name] = true
	}
	kept := make([]string, 0, len(environ))
	for _, assignment := range environ {
		key, _, _ := strings.Cut(assignment, "=")
		if !drop[key] {
			kept = append(kept, assignment)
		}
	}
	return kept
}

// TestTheScrubbedChildEnvironmentCarriesNoRegistryCredential is the control for
// barrier 2. Without it, "the environment was scrubbed" is a claim about a
// function call rather than about the bytes the child receives.
func TestTheScrubbedChildEnvironmentCarriesNoRegistryCredential(t *testing.T) {
	scrubbed := withoutVariables(childenv.Scrubbed(os.Environ()), "CI", "GITHUB_ACTIONS")
	for _, assignment := range scrubbed {
		key, _, _ := strings.Cut(assignment, "=")
		if key == childenv.CredentialVariable {
			t.Fatalf("the scrubbed child environment still assigns %s", childenv.CredentialVariable)
		}
	}
	// KNOWN-POSITIVE CONTROL: the scrubber must be capable of REMOVING
	// something. A no-op scrubber would produce the same empty result above on a
	// machine where the variable happens to be unset — which is most machines,
	// and precisely how this control would otherwise be worthless.
	planted := append([]string{}, "PATH=/usr/bin", childenv.CredentialVariable+"=not-a-real-key")
	if len(childenv.Scrubbed(planted)) != 1 {
		t.Fatalf("known-positive control failed: the scrubber left %v intact",
			childenv.Scrubbed(planted))
	}
	t.Logf("AC23 barrier 2: scrubbed environment carries %d variables and no registry credential "+
		"(control: a planted assignment was removed)", len(scrubbed))
}

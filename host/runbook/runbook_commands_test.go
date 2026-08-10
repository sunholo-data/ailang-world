package runbook

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// The attended publish runbook hands a HUMAN commands to run against an
// IRREVERSIBLE public write. A runbook that names a script which has been
// renamed or deleted is worse than no runbook: it fails at the keyboard, in the
// middle of an attended session, on the one path where improvisation is most
// expensive.
//
// Nothing else in this repo can catch that. `verify_ail.sh` and `verify_go.sh`
// gate CODE; a markdown file naming `./scripts/gone.sh` is green under every
// existing gate. So this test reads the commands OUT OF THE DOCUMENT and checks
// the artifacts they name, rather than re-deriving a list someone maintains by
// hand alongside the doc — a second hand-maintained list would rot in exactly
// the same way and would prove only that the two lists agree.
//
// MUT-SM-RUNBOOK-GHOST: rename any `./scripts/*.sh` occurrence in the runbook
// to a script that does not exist; this test reds naming that exact path.

// repoRoot resolves the repository root from this file's own location, so the
// gate does not depend on the working directory a test runner happens to use.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repository root")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

const runbookPath = "docs/SELF_MOD_PUBLISH.md"

// scriptInvocation matches a repo-relative shell script invocation as it
// appears in the runbook's fenced blocks.
var scriptInvocation = regexp.MustCompile(`\./scripts/[A-Za-z0-9_.-]+\.sh`)

func TestRunbookNamesOnlyScriptsThatExist(t *testing.T) {
	root := repoRoot(t)
	doc, err := os.ReadFile(filepath.Join(root, runbookPath))
	if err != nil {
		t.Fatalf("attended publish runbook is unreadable at %s: %v", runbookPath, err)
	}

	seen := map[string]bool{}
	for _, match := range scriptInvocation.FindAllString(string(doc), -1) {
		seen[match] = true
	}
	named := make([]string, 0, len(seen))
	for path := range seen {
		named = append(named, path)
	}
	sort.Strings(named)

	// ANTI-VACUITY (rule 3a): an empty extraction is indistinguishable from a
	// runbook that names no scripts, and both would let every assertion below
	// pass over nothing. A broken matcher must FAIL here, never pass quietly.
	if len(named) == 0 {
		t.Fatalf("instrument failure: extracted ZERO script invocations from %s (%d bytes). "+
			"Either the runbook stopped handing the operator commands, or scriptInvocation no "+
			"longer matches the document's shape", runbookPath, len(doc))
	}
	t.Logf("runbook names %d distinct script(s): %s", len(named), strings.Join(named, ", "))

	for _, rel := range named {
		info, err := os.Stat(filepath.Join(root, strings.TrimPrefix(rel, "./")))
		if err != nil {
			t.Errorf("%s tells the operator to run %s, which does not exist: %v",
				runbookPath, rel, err)
			continue
		}
		if info.IsDir() {
			t.Errorf("%s names %s, which is a directory", runbookPath, rel)
			continue
		}
		// An attended operator runs these directly, not via `sh <path>`, so the
		// executable bit is part of the claim the runbook is making.
		if info.Mode().Perm()&0o111 == 0 {
			t.Errorf("%s tells the operator to run %s, which is not executable (mode %v)",
				runbookPath, rel, info.Mode().Perm())
		}
	}

	// NEGATIVE CONTROL (rule 3d): the existence predicate above must be capable
	// of REPORTING a miss. Without this, a stat that could never fail — a
	// silently-rewritten path, a root that resolves somewhere unexpected —
	// yields the same green as a genuinely intact runbook.
	ghost := filepath.Join(root, "scripts", "this_script_does_not_exist.sh")
	if _, err := os.Stat(ghost); err == nil {
		t.Fatalf("negative control failed: %s exists, so a missing script would not be detectable", ghost)
	}
}

// The runbook's whole contract is that it STOPS before the irreversible write
// unless a human is present. If the automated stage ever grew a publish
// invocation, the document would still read as safe while instructing an
// unattended operator into a permanent public artifact.
//
// MUT-SM-RUNBOOK-UNATTENDED: add `ailang publish` (without `--dry-run`) to
// Stage A; this test reds.
//
// SM.D0 — WIDENED, NOT REPLACED, and the widening is the repair of a measured
// defect in this very function.
//
// As landed, its detection predicate was `strings.Contains(trimmed, "ailang
// publish")`. Measured at 6d1dce0: `grep -c 'ailang publish'
// docs/SELF_MOD_PUBLISH.md` = 0. THE LOOP BODY HAD NEVER EXECUTED. The two
// instrument-failure fatals above check the REGION being scanned — that the
// Stage B marker exists, that Stage A contains the readiness gate — and nothing
// checked the PREDICATE doing the scanning, so a green here was
// indistinguishable from a matcher that could never match.
//
// It now calls livePublishCommand, the SINGLE shared predicate, which
// TestOnePredicateSeesZeroLivePublishesInStageAAndAtLeastOneInStageB drives in
// the POSITIVE direction against Stage B. A broken matcher therefore yields a
// Stage-B count of zero and reds there, so the zero asserted here can no longer
// be produced by a dead instrument. The `--dry-run` exemption the landed version
// carried is preserved inside that predicate verbatim.
func TestRunbookStageAPerformsNoPublicWrite(t *testing.T) {
	root := repoRoot(t)
	doc, err := os.ReadFile(filepath.Join(root, runbookPath))
	if err != nil {
		t.Fatalf("attended publish runbook is unreadable at %s: %v", runbookPath, err)
	}

	const stageBMarker = "## Stage B"
	idx := strings.Index(string(doc), stageBMarker)
	// ANTI-VACUITY: without the marker the "automated stage" would be the empty
	// string and every assertion below would pass over nothing.
	if idx < 0 {
		t.Fatalf("instrument failure: %s no longer contains %q, so the automated stage cannot be "+
			"delimited and this gate would scan an empty region", runbookPath, stageBMarker)
	}
	stageA := string(doc[:idx])

	// CONTROL: the region we are about to scan must actually contain the
	// automated gate, or we are proving the absence of publish in a region that
	// contains nothing at all.
	if !strings.Contains(stageA, "./scripts/verify_world_package.sh") {
		t.Fatalf("instrument failure: the region before %q does not contain the readiness gate, "+
			"so it is not the automated stage", stageBMarker)
	}

	// CONTROL ON THE PREDICATE ITSELF, in the same call as the zero it produces.
	// Without this, "no line in Stage A matched" is also what a predicate that
	// matches nothing at all reports — which is precisely the state this function
	// shipped in.
	if !livePublishCommand("ailang publish") || livePublishCommand("ailang publish --dry-run") {
		t.Fatal("instrument failure: the shared livePublishCommand predicate does not discriminate " +
			"a live publish from a dry run, so the Stage-A zero below proves nothing")
	}

	for _, line := range strings.Split(stageA, "\n") {
		trimmed := strings.TrimSpace(line)
		if !livePublishCommand(trimmed) {
			continue
		}
		t.Errorf("%s Stage A (the UNATTENDED stage) contains a live publish: %q. "+
			"Stage A must stop at readiness; live publishes belong under %q, which is "+
			"gated on an attended approval", runbookPath, trimmed, stageBMarker)
	}
}

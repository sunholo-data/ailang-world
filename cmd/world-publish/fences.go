package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/pkgproj"
)

// ---------------------------------------------------------------------------
// THE STOP CONTRACT
//
// The DEFAULT OUTCOME OF EVERY NEW COMMAND PATH IS STOP. Every refusal below
// exits 3 and prints exactly one machine-readable line to stderr:
//
//	STOP fence=<name>[ reason=<reason>]
//
// The line is the contract; the human-readable detail follows on its own line
// and is not part of it. A runbook gate, a CI grep and an operator all read the
// same token.
// ---------------------------------------------------------------------------

const (
	// exitOK is a completed read-only verb, or a completed publish.
	exitOK = 0
	// exitError is a failure of something that was allowed to proceed — a
	// dispatch that failed, an indeterminate publish awaiting reconciliation.
	exitError = 1
	// exitUsage is a malformed invocation: unknown verb, unparseable flags.
	exitUsage = 2
	// exitStop is a FENCE. Nothing irreversible happened, and nothing will.
	exitStop = 3
)

// The fence names. They are constants because three separate things read them:
// this package's table test, host/runbook's rule-3k arm, and a human.
const (
	fenceMode         = "mode"
	fenceCI           = "ci"
	fenceTTY          = "tty"
	fenceConfirmation = "confirmation"
	fenceStore        = "store"
	fenceApproval     = "approval"
	fenceCredential   = "credential"
	fencePacket       = "packet"
	// fenceHandler is the REPORTING site for the inherited constructor
	// refusals (origin, ambient credential, missing approval, missing
	// credential provider). It adds no refusal of its own — every one of them
	// is raised inside host/broker — but it needs a name of its own because it
	// is the LAST stage, strictly after the attended fence, and AC21's
	// per-row positive controls are stated in terms of reaching a later stage.
	fenceHandler = "handler"
)

// attendedPhrase is the exact line an attended operator must type. It names the
// package, the version and the irreversibility, so it cannot be typed by
// accident and cannot be reused for a different publish.
//
// It defeats ACCIDENT, not automation — a script can echo it into a pipe. The
// layer that defeats automation is the controlling-terminal check in tty.go.
// Stating that here is deliberate: an operator who believes this phrase is the
// security boundary will eventually route around the one that is.
const attendedPhrase = "publish world/core@0.1.0 irreversibly"

// frozenPackageVersion is the ONE version this command may publish. The scope
// grammar already makes an approval for 0.1.0 unable to authorize 0.1.1; this
// is the same statement at the operator's surface, where a typo lives.
const frozenPackageVersion = "0.1.0"

// frozenCompilerVersion is the pinned compiler the ready packet was projected
// with. It is stated here INDEPENDENTLY of the golden so the comparison in
// R-PACKET-DRIFT has two sources: a golden whose compilerVersion drifted away
// from this constant reds rather than being copied into the recomputation.
const frozenCompilerVersion = "AILANG v0.30.0"

// stopError is a fence refusal. It is a value rather than a printf because the
// STOP line has to be produced identically everywhere it is produced.
type stopError struct {
	Fence  string
	Reason string
	Detail string
}

// Line is the machine-readable contract.
func (e *stopError) Line() string {
	if e.Reason == "" {
		return "STOP fence=" + e.Fence
	}
	return "STOP fence=" + e.Fence + " reason=" + e.Reason
}

func (e *stopError) Error() string {
	if e.Detail == "" {
		return e.Line()
	}
	return e.Line() + ": " + e.Detail
}

// report writes the refusal and returns the exit code, so every refusal site is
// a single `return report(w, err)` and none can print the line differently.
func report(w io.Writer, err *stopError) int {
	fmt.Fprintln(w, err.Line())
	if err.Detail != "" {
		fmt.Fprintln(w, "  "+err.Detail)
	}
	return exitStop
}

// ---------------------------------------------------------------------------
// R-MODE-NONE / R-MODE-BOTH
// ---------------------------------------------------------------------------

// requireExactlyOneMode refuses an ambiguous intent. The irreversible path is
// never the default and never the fallback: it is reached only by naming it,
// and only by naming it alone.
func requireExactlyOneMode(dryRun, live bool) *stopError {
	// R-MODE-BOTH
	if dryRun && live {
		return &stopError{Fence: fenceMode, Reason: "both",
			Detail: "--dry-run and --live are mutually exclusive"}
	}
	// R-MODE-NONE
	if !dryRun && !live {
		return &stopError{Fence: fenceMode, Reason: "none",
			Detail: "publish requires exactly one of --dry-run or --live"}
	}
	return nil
}

// refuseLiveOnReadOnlyVerb guards the deliberate asymmetry that makes the
// read-only verbs headless-permitted.
//
// `packet` and `reconcile` are the two verbs a headless loop MAY run: a loop
// may inspect the projection, and it may reconcile an already-authorized
// indeterminate attempt read-only. That permission is only safe if a read-only
// verb cannot be talked into the write path by adding one flag.
func refuseLiveOnReadOnlyVerb(verb string, live bool) *stopError {
	// R-RECONCILE-LIVE-FLAG
	if live {
		return &stopError{Fence: fenceMode, Reason: verb + "-is-read-only",
			Detail: verb + " performs no public write; --live is refused rather than ignored"}
	}
	return nil
}

// ---------------------------------------------------------------------------
// R-CI — a TRIPWIRE, DECLARED. NOT THE FENCE.
// ---------------------------------------------------------------------------

// ciVariables are the two markers every runner in use here sets.
var ciVariables = []string{"CI", "GITHUB_ACTIONS"}

// refuseAutomationEnvironment is a TRIPWIRE and is labelled one in the code, in
// the acceptance criterion and in the runbook.
//
// It is defeated by `env -u CI`. It is kept because it costs three lines and it
// would stop a future CI step that somehow allocated a pty. IT IS NOT THE
// FENCE. The fence is requireControllingTerminal. Anyone who lets this become
// load-bearing has replaced a structural property with a naming convention.
//
// Measured at 6d1dce0 in the executing agent's shell: CI="" and
// GITHUB_ACTIONS="", so this branch does not fire during development and cannot
// mask a failure of the TTY arms.
func refuseAutomationEnvironment(getenv func(string) string) *stopError {
	for _, name := range ciVariables {
		// R-CI
		if getenv(name) != "" {
			return &stopError{Fence: fenceCI, Reason: strings.ToLower(name),
				Detail: name + " is set; this command is never run in CI"}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// R-PHRASE-EOF / R-PHRASE
// ---------------------------------------------------------------------------

// requireTypedConfirmation reads ONE line and compares it with the exact
// phrase.
func requireTypedConfirmation(in io.Reader, out io.Writer) *stopError {
	fmt.Fprintf(out, "Type exactly, to proceed with an IRREVERSIBLE public write:\n  %s\n> ", attendedPhrase)
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	// R-PHRASE-EOF. A closed stdin is not a silent yes. Note the deliberate
	// order: io.EOF with a complete line already read is NOT this branch, so a
	// terminal that sends the phrase without a trailing newline still works.
	if err != nil && strings.TrimSpace(line) == "" {
		return &stopError{Fence: fenceConfirmation, Reason: "eof",
			Detail: "stdin closed before a confirmation line was typed"}
	}
	// R-PHRASE
	if strings.TrimSpace(line) != attendedPhrase {
		return &stopError{Fence: fenceConfirmation, Reason: "mismatch",
			Detail: "the typed line is not the required confirmation phrase"}
	}
	return nil
}

// requireAttendedOperator is the composed fence stack, and it is the ONE call
// that must dominate the construction of a production publish handler.
//
// AC22 asserts, over the AST, that runPublish names this function exactly once,
// names broker.NewRegistryPublishHandler exactly once, and that this call's
// statement index is strictly lower. That assertion — not this comment — is
// what distinguishes "there is a fence" from "the fence dominates the
// irreversible path".
func requireAttendedOperator(in io.Reader, out io.Writer, getenv func(string) string, probe ttyProbe) *stopError {
	if err := refuseAutomationEnvironment(getenv); err != nil {
		return err
	}
	if err := requireControllingTerminal(probe); err != nil {
		return err
	}
	return requireTypedConfirmation(in, out)
}

// ---------------------------------------------------------------------------
// R-STORE / R-APPROVAL-ABSENT / R-CRED-FLAG
// ---------------------------------------------------------------------------

// requireStorePath refuses the ABSENCE of the flag. Opening it is a separate
// refusal on the same fence, raised by the caller, because "you did not say
// which world" and "that world will not open" are different mistakes.
func requireStorePath(path string) *stopError {
	// R-STORE (the flag half)
	if strings.TrimSpace(path) == "" {
		return &stopError{Fence: fenceStore, Reason: "absent",
			Detail: "--store names the world database this publish is recorded in"}
	}
	return nil
}

func storeOpenFailed(path string, err error) *stopError {
	return &stopError{Fence: fenceStore, Reason: "unopenable",
		Detail: "opening " + path + " failed: " + err.Error()}
}

// requireApprovalRef refuses a live publish with no attended stamp NAMED.
//
// This is NOT a duplicate of validatePublishApproval. That function refuses a
// stamp that is WRONG — missing, expired, already consumed, wrong-scope,
// wrong-hash, denied, malformed — and it needs a store and a payload to do so.
// This refuses the absence of the FLAG, before a store is even consulted. The
// seven wrongness classes stay inherited and are not restated.
func requireApprovalRef(raw string) (hashref.HashRef, *stopError) {
	// R-APPROVAL-ABSENT (the absent half)
	if strings.TrimSpace(raw) == "" {
		return hashref.HashRef{}, &stopError{Fence: fenceApproval, Reason: "absent",
			Detail: "--approval-ref names the ApprovalDecisionV1 minted by `world-publish approve`; " +
				"minting and spending are deliberately two separate invocations"}
	}
	ref, err := hashref.Parse(strings.TrimSpace(raw))
	// R-APPROVAL-ABSENT (the unparseable half — the same branch's second arm)
	if err != nil {
		return hashref.HashRef{}, &stopError{Fence: fenceApproval, Reason: "malformed",
			Detail: "--approval-ref is not a hashref: " + err.Error()}
	}
	return ref, nil
}

// requireCredentialFile refuses a live publish with no credential file NAMED.
// The file's mode, ownership, location and readability refusals are INHERITED
// from broker.NewFileRegistryCredentialProvider and are not restated here.
func requireCredentialFile(path string) *stopError {
	// R-CRED-FLAG
	if strings.TrimSpace(path) == "" {
		return &stopError{Fence: fenceCredential, Reason: "absent",
			Detail: "--credential-file names a file outside the working tree holding the registry API key; " +
				"the key is never read from the environment"}
	}
	return nil
}

// ---------------------------------------------------------------------------
// R-PACKET-DRIFT / R-PACKET-VERSION
// ---------------------------------------------------------------------------

// requireUndriftedPacket is the fence that binds this command to the REVIEWED
// ARTIFACT rather than to whatever happens to be in packages/ right now.
//
// The design doc's AC18 says the publish must "equal the ready packet
// byte-for-byte" and names no artifact — there was no canonical ready-packet
// file when it was written. This function is where that sentence acquires a
// referent: scripts/world_package_ready_packet.golden.json, by name.
func requireUndriftedPacket(packageDir, goldenPath string) (pkgproj.ReadyPacket, *stopError) {
	golden, err := pkgproj.LoadReadyPacket(goldenPath)
	if err != nil {
		return pkgproj.ReadyPacket{}, &stopError{Fence: fencePacket, Reason: "golden-unreadable",
			Detail: err.Error()}
	}
	// R-PACKET-VERSION. Checked BEFORE the drift comparison so an operator who
	// pointed at the wrong golden is told that, rather than being handed a list
	// of nine differing fields.
	if golden.Version != frozenPackageVersion {
		return pkgproj.ReadyPacket{}, &stopError{Fence: fencePacket, Reason: "version",
			Detail: fmt.Sprintf("the golden names version %q; this command publishes only %q",
				golden.Version, frozenPackageVersion)}
	}
	recomputed, recomputeErr := pkgproj.RecomputeReadyPacket(packageDir, worldCoreManifest, frozenCompilerVersion)
	if recomputeErr != nil {
		return pkgproj.ReadyPacket{}, &stopError{Fence: fencePacket, Reason: "unrecomputable",
			Detail: recomputeErr.Error()}
	}
	// R-PACKET-DRIFT
	if field, equal := recomputed.Equal(golden); !equal {
		return pkgproj.ReadyPacket{}, &stopError{Fence: fencePacket, Reason: "drift",
			Detail: fmt.Sprintf("%s differs: recomputed %q, golden %q — the bytes on disk are not "+
				"the artifact that was reviewed", field, recomputed.Field(field), golden.Field(field))}
	}
	return golden, nil
}

// worldCoreManifest is packages/world-core/ailang.toml stated in Go.
//
// It is a LITERAL rather than a parse because this repository has no TOML
// dependency and adding one to read six lines would be a worse trade. The
// literal cannot rot silently: every field of it feeds the recomputed packet,
// which R-PACKET-DRIFT compares to the committed golden that the shell gate
// produced FROM THE REAL FILE. A literal that disagreed with the toml would
// change the interface hash, the exports or the version, and the fence would
// refuse. TestWorldCoreManifestMatchesTheCommittedGolden is that measurement.
var worldCoreManifest = pkgproj.Manifest{
	Package: pkgproj.Package{
		Name:    "world/core",
		Edition: "1",
		AILANG:  ">=0.30.0",
		Version: frozenPackageVersion,
	},
	Exports: pkgproj.Exports{Modules: []string{
		"world/types", "world/contracts", "world/transitions", "world/logepoch",
	}},
	Effects: pkgproj.Effects{Max: []string{}},
}

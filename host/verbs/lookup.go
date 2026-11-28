package verbs

import (
	"regexp"
	"strings"
)

// Class is what the bar may do with a verb (design §2.1(b)).
type Class string

const (
	Address  Class = "address"  // navigation: the page's own parseAddress path
	Read     Class = "read"     // reads World: episode report and why run
	Navigate Class = "navigate" // approve opens row 201's card
	Mutate   Class = "mutate"   // changes World: terminal (V3/V4)
	Host     Class = "host"     // needs a host process: terminal (V4)
	Never    Class = "never"    // typed-only, never from a page
	Other    Class = "other"    // a CLI command that is not an operator verb
	Refused  Class = "refused"  // a global flag, or a seized first word with no row
)

// Verb is one answer of Lookup: the table row's words (empty for an address), its class, the
// reason the bar gives, and the CLI line (empty when the input is not a CLI line).
type Verb struct {
	Words  string
	Class  Class
	Reason string
	CLI    string
}

// Binary is the CLI's name; a pasted line may start with it.
const Binary = "ailang-worldd"

// The table's word sets. They are exactly the CLI switch's arms but world (VB4 asserts it).
var (
	OperatorWords = []string{"why", "episode", "approve"}
	NeverWords    = []string{"allowance", "approver", "session", "service", "serve", "setup"}
	OtherWords    = []string{"app", "commit", "outcome", "call", "demo", "doctor", "head", "health", "help", "hosts", "log", "mcp-stdio", "object",
		"preset", "profile", "provenance", "registry", "run", "surface-index", "tools", "version", "whoami"}
)

// neverWhy is each never verb's reason (D-WORLD-89 R3; host administration stays at a terminal).
var neverWhy = map[string]string{
	"allowance": "granting and revoking a standing allowance is typed at a terminal",
	"approver":  "enrolling and revoking an approver key is typed at a terminal",
	"session":   "minting and revoking a session is typed at a terminal",
	"service":   "host administration stays at a terminal",
	"serve":     "host administration stays at a terminal",
	"setup":     "host administration stays at a terminal",
}

// EpisodeNeedsSubcommand is step 4b's refusal for a bare or unknown episode form.
const EpisodeNeedsSubcommand = Binary + " episode needs a subcommand: start | launch | finish | report"

// AddressHelp is the page's own address help (surface.js ADDRESS_HELP), step 4b's text for world.
const AddressHelp = `unknown address: try "world", "world NAME", "world NAME@N", "world NAME@HASH" or "back"`

func member(set []string, w string) bool {
	for _, s := range set {
		if s == w {
			return true
		}
	}
	return false
}

// Lookup classifies argv (design §2.1(b), steps 1–5) and returns the verb and the words after it.
func Lookup(argv []string) (Verb, []string) {
	if len(argv) > 0 && argv[0] == Binary { // 1
		argv = argv[1:]
	}
	if len(argv) == 0 {
		return Verb{Class: Address}, nil
	}
	cli := Quote(append([]string{Binary}, argv...))
	first := argv[0]
	if strings.HasPrefix(first, "-") { // 2
		return Verb{Class: Refused, Reason: "global flags are not accepted: the app talks to its own daemon", CLI: cli}, nil
	}
	if first == "world" && len(argv) <= 2 { // 3 (back is not a CLI word: step 5 keeps it an address)
		return Verb{Class: Address}, nil
	}
	if first == "world" && argv[1] == "get" { // 3b
		return Verb{Words: "world get", Class: Other, CLI: cli,
			Reason: "world get is a CLI command the app does not run; at a terminal: " + cli}, argv[2:]
	}
	if first == "episode" && len(argv) > 1 { // 4: the two-word rows
		switch argv[1] {
		case "report":
			return Verb{Words: "episode report", Class: Read, CLI: cli}, argv[2:]
		case "start":
			return Verb{Words: "episode start", Class: Mutate, CLI: cli,
				Reason: "episode start changes World. The app cannot run it yet."}, argv[2:]
		case "finish": // row 243 V3 §2.8: still terminal; the reason names the CLI's surface mode
			return Verb{Words: "episode finish", Class: Mutate, CLI: cli,
				Reason: "episode finish changes World. The app cannot run it yet; at a terminal, --confirm-on-surface lets you approve its pull request here with your key."}, argv[2:]
		case "launch":
			return Verb{Words: "episode launch", Class: Host, CLI: cli,
				Reason: "episode launch needs a host process. The app cannot run it yet."}, argv[2:]
		}
	}
	if first == "profile" && len(argv) > 1 && argv[1] == "list" { // row 269 Phase 1: read-only
		return Verb{Words: "profile list", Class: Read, CLI: cli}, argv[2:]
	}
	if first == "preset" && len(argv) > 1 && argv[1] == "list" { // row 269 Phase 1: read-only
		return Verb{Words: "preset list", Class: Read, CLI: cli}, argv[2:]
	}
	if first == "app" && len(argv) > 1 && argv[1] == "list" { // row 274 M1: read-only
		return Verb{Words: "app list", Class: Read, CLI: cli}, argv[2:]
	}
	switch { // 4: the one-word rows
	case first == "why":
		return Verb{Words: "why", Class: Read, CLI: cli}, argv[1:]
	case first == "approve":
		return Verb{Words: "approve", Class: Navigate, CLI: cli}, argv[1:]
	case member(NeverWords, first):
		return Verb{Words: first, Class: Never, CLI: cli,
			Reason: first + " is never run from the app: " + neverWhy[first] + "; at a terminal: " + cli}, argv[1:]
	case member(OtherWords, first):
		return Verb{Words: first, Class: Other, CLI: cli,
			Reason: first + " is a CLI command the app does not run; at a terminal: " + cli}, argv[1:]
	}
	switch first { // 4b: a seized first word with no row
	case "episode":
		return Verb{Words: "episode", Class: Refused, Reason: EpisodeNeedsSubcommand, CLI: cli}, argv[1:]
	case "world":
		return Verb{Words: "world", Class: Refused, Reason: AddressHelp}, argv[1:]
	}
	return Verb{Class: Address}, nil // 5
}

var (
	mintRequestID = regexp.MustCompile(`^mr-[0-9a-f]{16}$`)
	decisionRef   = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	verbRequestID = regexp.MustCompile(`^vr-[0-9a-f]{16}$`)
)

// ApproveTarget is approve's navigate arm: one request id or decision ref, no flag, opens its card.
func ApproveTarget(args []string) (href, reason string) {
	if len(args) == 1 && mintRequestID.MatchString(args[0]) {
		return "/surface/approve/" + args[0], ""
	}
	if len(args) == 1 && decisionRef.MatchString(args[0]) {
		return "/surface/decide/" + args[0], ""
	}
	if len(args) == 1 && verbRequestID.MatchString(args[0]) { // row 243 V3 §2.8: a finish confirmation's card
		return "/surface/confirm/" + args[0], ""
	}
	return "", "the app opens approve <request-id|request-ref> as its card; any flag or other form runs at a terminal"
}

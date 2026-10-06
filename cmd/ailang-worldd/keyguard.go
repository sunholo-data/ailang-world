package main

import (
	"errors"
	"flag"
	"io"
	"strings"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
)

// ---------------------------------------------------------------------------
// The registry-key guard's two operator affordances (2026-10-06).
//
// The guard itself is unchanged and is not weakened: with
// AILANG_REGISTRY_API_KEY in the environment every verb that DOES anything
// refuses (w-self-mod-vertical Decision 4), and nothing here unsets the
// variable. What changes is that:
//
//   - a help request is answered anyway. Help prints a constant and executes
//     nothing — no verb function is called, nothing is forked, no store is
//     opened — so there is no process for the key to leak into. The text is
//     resolved from the argv alone, BEFORE any verb runs, so a help token in a
//     position a verb's own parser would treat as data cannot make the verb run
//     with the key set;
//   - a refusal ends with the exact command to run instead.
// ---------------------------------------------------------------------------

// isHelpToken reports whether an argument asks for help.
func isHelpToken(a string) bool {
	return a == "--help" || a == "-help" || a == "-h"
}

// helpOnly resolves the help text an invocation asks for, without running any
// verb. ok is false when the invocation is not a help request, or names a verb
// this function does not know (the guard then refuses it as before).
func helpOnly(args []string) (string, bool) {
	globals := flag.NewFlagSet("ailang-worldd", flag.ContinueOnError)
	globals.SetOutput(io.Discard)
	globals.String("addr", daemon.DefaultAddr, "")
	if err := globals.Parse(args); err != nil {
		return usage, errors.Is(err, flag.ErrHelp)
	}
	rest := globals.Args()
	if len(rest) == 0 {
		return "", false
	}
	if rest[0] == "help" {
		if len(rest) == 1 {
			return usage, true
		}
		rest = append(append([]string(nil), rest[1:]...), "--help")
	}
	asked := false
	for _, a := range rest[1:] {
		asked = asked || isHelpToken(a)
	}
	if !asked {
		return "", false
	}
	if text, ok := genericVerbHelp(rest); ok {
		return text, true
	}
	return ownVerbHelp(rest)
}

// ownVerbHelp is the help each ownHelpVerbs verb (and `log tail`) prints for
// --help. TestHelpIsNeverBlockedByTheRegistryKeyGuard compares it byte for
// byte with what the verb itself prints, for every verb line.
func ownVerbHelp(rest []string) (string, bool) {
	switch rest[0] {
	case "tools":
		return toolsHelp, true
	case "call":
		return callHelp, true
	case "why":
		return whyHelp, true
	case "provenance":
		return provenanceHelp, true
	case "setup":
		return setupHelp, true
	case "doctor":
		return doctorHelp, true
	case "log":
		if len(rest) > 1 && rest[1] == "tail" {
			return logTailHelp, true
		}
	case "session":
		if len(rest) < 2 || isHelpToken(rest[1]) {
			return sessionHelp, true
		}
		switch rest[1] {
		case "new":
			return sessionNewHelp, true
		case "list":
			return sessionListHelp, true
		case "revoke":
			return sessionRevokeHelp, true
		case "mint":
			return sessionHelp, true
		}
	}
	return "", false
}

// subVerbs are the verbs whose second word is part of the command's name.
var subVerbs = map[string]bool{
	"world": true, "object": true, "log": true, "registry": true, "tools": true, "session": true,
}

// registryKeyFix is the last line of the guard's refusal: the exact command,
// for the verb that was refused, with the variable removed for that one
// process (and so for everything it starts). The rest of the argv is elided
// as "…", never echoed: an argument can be a session token.
func registryKeyFix(args []string) string {
	globals := flag.NewFlagSet("ailang-worldd", flag.ContinueOnError)
	globals.SetOutput(io.Discard)
	globals.String("addr", daemon.DefaultAddr, "")
	name := "ailang-worldd"
	if globals.Parse(args) == nil {
		if rest := globals.Args(); len(rest) > 0 {
			name += " " + rest[0]
			if subVerbs[rest[0]] && len(rest) > 1 && !strings.HasPrefix(rest[1], "-") {
				name += " " + rest[1]
			}
		}
	}
	return "  fix: run it without the key in its environment: env -u " +
		broker.RegistryCredentialVariable + " " + name + " …"
}

// registryKeyRationale is the Decision-4 reason, said in the operator's terms.
const registryKeyRationale = "  why: the public AILANG registry is immutable, so with the key in the environment every " +
	"process World starts would inherit unrecallable publish authority (w-self-mod-vertical Decision 4). " +
	"It is never auto-unset, and only --help runs with it set."

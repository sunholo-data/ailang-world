// Package childenv is the single definition of the registry-related
// environment World refuses to hand a child process.
//
// It exists because the registry credential must not travel by inheritance.
// Decision 4 of design_docs/planned/w-self-mod-vertical.md requires that every
// World-launched process except the one live publish dispatch observe
// AILANG_REGISTRY_API_KEY unset, and that the publish handler strip "every
// registry-related variable" before adding the credential back for exactly one
// child. That is a property of a LIST OF NAMES, and a list copied into four
// packages is a list that drifts silently, so the names live here once and
// every subprocess launch site reads them from here.
//
// This package is a leaf: it imports only the standard library, holds no
// state, and is imported by host/archive, host/broker and host/replay.
package childenv

import "strings"

// CredentialVariable is the registry API key the AILANG v0.30.0 publisher
// reads (e37b370:cmd/ailang/pkg_publish.go:270). It is the one variable whose
// leakage is irreversible, because the public registry is immutable.
const CredentialVariable = "AILANG_REGISTRY_API_KEY"

// RegistryVariables is the exact set of registry-related variables stripped
// from every World-launched child. All four are read by the pinned binary:
// AILANG_REGISTRY selects the read-only bucket, AILANG_REGISTRY_VALIDATOR and
// its deprecated alias AILANG_REGISTRY_API select the publish service
// (e37b370:cmd/ailang/pkg_info.go:19-23), and CredentialVariable authenticates
// the write. Stripping the three non-secret ones matters as much as the
// secret: a child that inherits an unexpected origin publishes, or reports
// absence, against a registry nobody chose.
var RegistryVariables = []string{
	"AILANG_REGISTRY",
	"AILANG_REGISTRY_API",
	CredentialVariable,
	"AILANG_REGISTRY_VALIDATOR",
}

// Scrubbed returns environ without any RegistryVariables assignment. It
// returns a fresh slice and never mutates its argument, so a caller may hand
// os.Environ() straight through. A nil or empty environ yields a nil result,
// which exec's cmd.Env treats as "inherit" — callers that require a stripped
// environment must therefore pass a real environ, and the ones in this
// repository pass os.Environ().
func Scrubbed(environ []string) []string {
	var kept []string
	for _, assignment := range environ {
		if !isRegistryAssignment(assignment) {
			kept = append(kept, assignment)
		}
	}
	return kept
}

// Has reports whether environ assigns name a non-empty value. Emptiness is the
// same "absent" the pinned publisher applies to the credential
// (e37b370:cmd/ailang/pkg_publish.go:270 tests apiKey != ""), so a variable
// present with an empty value is not ambient authority and is not reported as
// one.
func Has(environ []string, name string) bool {
	value, ok := Lookup(environ, name)
	return ok && value != ""
}

// Lookup returns the last assignment of name in environ, matching the
// last-wins rule the operating system applies to a duplicated variable.
func Lookup(environ []string, name string) (string, bool) {
	var (
		value string
		found bool
	)
	for _, assignment := range environ {
		key, rest, ok := strings.Cut(assignment, "=")
		if ok && key == name {
			value, found = rest, true
		}
	}
	return value, found
}

func isRegistryAssignment(assignment string) bool {
	key, _, ok := strings.Cut(assignment, "=")
	if !ok {
		return false
	}
	for _, name := range RegistryVariables {
		if key == name {
			return true
		}
	}
	return false
}

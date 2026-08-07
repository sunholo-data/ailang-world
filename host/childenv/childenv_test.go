package childenv

import (
	"strings"
	"testing"
)

func TestScrubbedRemovesEveryRegistryVariableAndNothingElse(t *testing.T) {
	// The control entries must SURVIVE. A scrubber that returned nil would
	// satisfy "the credential is gone" while breaking every child it touches,
	// so absence alone is not the criterion.
	control := []string{"PATH=/usr/bin:/bin", "HOME=/tmp", "AILANG_FS_SANDBOX=/tmp/x"}
	environ := append([]string(nil), control...)
	for _, name := range RegistryVariables {
		environ = append(environ, name+"=value-for-"+name)
	}
	if len(RegistryVariables) == 0 {
		t.Fatal("RegistryVariables is empty: Scrubbed would be a no-op and every caller vacuous")
	}

	scrubbed := Scrubbed(environ)
	if len(scrubbed) != len(control) {
		t.Fatalf("Scrubbed kept %v, want exactly the %d control entries", scrubbed, len(control))
	}
	for i, want := range control {
		if scrubbed[i] != want {
			t.Errorf("Scrubbed[%d] = %q, want %q", i, scrubbed[i], want)
		}
	}
	for _, name := range RegistryVariables {
		if _, found := Lookup(scrubbed, name); found {
			t.Errorf("%s survived scrubbing", name)
		}
		// Known-positive control for the detector itself: it must SEE the
		// variable in the unscrubbed input.
		if _, found := Lookup(environ, name); !found {
			t.Errorf("the detector could not see %s in the unscrubbed environ", name)
		}
	}
	// Scrubbed must not mutate its argument.
	if !Has(environ, CredentialVariable) {
		t.Error("Scrubbed mutated the environ it was given")
	}
}

// TestScrubbedNeverReturnsNilBecauseNilMeansInherit is the guard for the one
// way this package can fail OPEN. exec reads a nil cmd.Env as "inherit the
// process environment", so a Scrubbed that returned nil for a degenerate input
// would hand a child the very credential it was called to remove — silently,
// and in the direction that costs an irreversible publish. Every production
// caller passes os.Environ() today, so this is unreachable now; it is a guard
// against the next caller, which is exactly when nobody re-reads the comment.
func TestScrubbedNeverReturnsNilBecauseNilMeansInherit(t *testing.T) {
	// Every input that could plausibly reduce to nothing.
	allRegistry := make([]string, 0, len(RegistryVariables))
	for _, name := range RegistryVariables {
		allRegistry = append(allRegistry, name+"=leaked")
	}
	for _, tc := range []struct {
		name   string
		input  []string
		reason string
	}{
		{"nil", nil, "a caller with no environment to hand on"},
		{"empty", []string{}, "an environment already emptied upstream"},
		{"only registry variables", allRegistry, "everything present was stripped"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Scrubbed(tc.input)
			if got == nil {
				t.Fatalf("Scrubbed(%s) returned nil (%s): exec would read that as "+
					"INHERIT and pass %s to the child", tc.name, tc.reason, CredentialVariable)
			}
			if len(got) != 0 {
				t.Errorf("Scrubbed(%s) = %v, want an empty (but non-nil) environment", tc.name, got)
			}
		})
	}

	// Known-positive control for the assertion itself: on an input that DOES
	// have survivors, the same call must come back non-nil AND non-empty. A
	// test that only ever demands non-nil would pass against a function that
	// returned an empty slice unconditionally — i.e. against a scrubber that
	// silently discarded PATH.
	if got := Scrubbed([]string{"PATH=/usr/bin", CredentialVariable + "=leaked"}); len(got) != 1 {
		t.Fatalf("control: Scrubbed kept %v, want exactly [PATH=/usr/bin] — "+
			"the non-nil assertion above is only meaningful if survivors survive", got)
	}
}

func TestHasTreatsAnEmptyAssignmentAsAbsentAuthority(t *testing.T) {
	if Has([]string{CredentialVariable + "="}, CredentialVariable) {
		t.Error("an empty assignment was reported as ambient authority")
	}
	if !Has([]string{CredentialVariable + "=x"}, CredentialVariable) {
		t.Error("a non-empty assignment was not detected")
	}
	if Has(nil, CredentialVariable) {
		t.Error("a nil environ reported the variable present")
	}
}

func TestLookupAppliesLastWins(t *testing.T) {
	value, found := Lookup([]string{"A=1", "A=2", "B=3"}, "A")
	if !found || value != "2" {
		t.Errorf("Lookup = %q, %v, want \"2\", true", value, found)
	}
	if _, found := Lookup([]string{"NOT_AN_ASSIGNMENT"}, "NOT_AN_ASSIGNMENT"); found {
		t.Error("a bare token without '=' was read as an assignment")
	}
}

func TestCredentialVariableIsListed(t *testing.T) {
	var listed bool
	for _, name := range RegistryVariables {
		listed = listed || name == CredentialVariable
		if strings.TrimSpace(name) != name || name == "" {
			t.Errorf("registry variable %q is malformed", name)
		}
	}
	if !listed {
		t.Fatalf("%s is not in RegistryVariables: Scrubbed would leave the secret in place",
			CredentialVariable)
	}
}

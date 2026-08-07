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

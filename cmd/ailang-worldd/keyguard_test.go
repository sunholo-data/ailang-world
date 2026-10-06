package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

const keyGuardSentinel = "keyguard-sentinel-not-a-real-key"

// helpForms is every way to ask for help: `<verb> --help`, `-h`, `-help`,
// `help <verb>`, for every verb line of the CLI, and the top-level `help`.
func helpForms() [][]string {
	forms := [][]string{{"help"}, {"--addr", "http://127.0.0.1:1", "help", "health"}}
	for _, verb := range everyVerb {
		for _, flag := range []string{"--help", "-h", "-help"} {
			forms = append(forms, append(append([]string{}, verb...), flag))
		}
		forms = append(forms, append([]string{"help"}, verb...))
	}
	// A help flag after the verb's other flags is still help.
	forms = append(forms, []string{"serve", "--db", "/tmp/x.db", "--help"})
	return forms
}

// --help never needs the registry-key guard: it prints text and runs nothing.
// With the key set, every help form exits 0 with EXACTLY the help it prints
// without the key (the negative control is the same argv, key absent), and
// neither stream carries the value.
func TestHelpIsNeverBlockedByTheRegistryKeyGuard(t *testing.T) {
	for _, argv := range helpForms() {
		name := strings.Join(argv, " ")
		t.Setenv(broker.RegistryCredentialVariable, "")
		var want, wantErr bytes.Buffer
		wantCode := run(argv, &want, &wantErr)
		if wantCode != exitOK {
			t.Fatalf("%q without the key: exit %d (instrument: every help form exits 0)\n%s", name, wantCode, wantErr.String())
		}

		t.Setenv(broker.RegistryCredentialVariable, keyGuardSentinel)
		var got, gotErr bytes.Buffer
		if code := run(argv, &got, &gotErr); code != exitOK {
			t.Errorf("%q with the key set: exit %d, want 0\n%s", name, code, gotErr.String())
			continue
		}
		if got.String() != want.String() {
			t.Errorf("%q with the key set printed different help:\n got %q\nwant %q", name, got.String(), want.String())
		}
		if strings.Contains(got.String()+gotErr.String(), keyGuardSentinel) {
			t.Fatalf("%q printed the key's value", name)
		}
	}
}

// When the guard refuses, the last line is the exact command to run instead,
// for the verb that was refused, and the Decision-4 rationale is kept. The
// rest of the command is elided, never echoed: an argument can be a session
// token.
func TestRegistryKeyRefusalEndsWithTheExactCommand(t *testing.T) {
	t.Setenv(broker.RegistryCredentialVariable, keyGuardSentinel)
	for _, tc := range []struct {
		argv []string
		want string
	}{
		{[]string{"serve", "--db", "/tmp/x.db"}, "env -u AILANG_REGISTRY_API_KEY ailang-worldd serve …"},
		{[]string{"session", "new", "ep1", "--db", "/tmp/x.db"}, "env -u AILANG_REGISTRY_API_KEY ailang-worldd session new …"},
		{[]string{"--addr", "http://127.0.0.1:1", "commit", "--session", "deadbeef"}, "env -u AILANG_REGISTRY_API_KEY ailang-worldd commit …"},
		{[]string{"doctor"}, "env -u AILANG_REGISTRY_API_KEY ailang-worldd doctor …"},
	} {
		var stdout, stderr bytes.Buffer
		if code := run(tc.argv, &stdout, &stderr); code != exitFatal {
			t.Fatalf("%q with the key set: exit %d, want %d (fatal)", tc.argv, code, exitFatal)
		}
		text := strings.TrimSpace(stderr.String())
		if !strings.HasSuffix(text, tc.want) {
			t.Errorf("%q refusal does not end with %q:\n%s", tc.argv, tc.want, text)
		}
		if !strings.Contains(text, "Decision 4") || !strings.Contains(text, "unrecallable") {
			t.Errorf("%q refusal dropped the Decision-4 rationale:\n%s", tc.argv, text)
		}
		if strings.Contains(text, keyGuardSentinel) || strings.Contains(text, "deadbeef") {
			t.Fatalf("%q refusal echoed a secret:\n%s", tc.argv, text)
		}
	}
}

package main

// Row 201 A2 sprint 4 M7 (design §2.15): host/approval/testdata/webauthn/generate.sh, its recorder
// and its refusals.

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestPasskeyPageRecordsTheFixtures is host/approval/testdata/webauthn/generate.sh's recorder:
// with WORLD_PASSKEY_RECORD_DIR set it records, through the page and a virtual authenticator, an
// ES256 and an EdDSA create and two gets over the fixture message hash (design §2.15), and writes
// chrome-es256.json and chrome-eddsa.json there in the committed fixtures' layout. Unset, it skips.
func TestPasskeyPageRecordsTheFixtures(t *testing.T) {
	dir := os.Getenv("WORLD_PASSKEY_RECORD_DIR")
	if dir == "" {
		t.Skip("not recording: host/approval/testdata/webauthn/generate.sh sets WORLD_PASSKEY_RECORD_DIR")
	}
	b := startPasskeyBrowser(t)
	origin := passkeyPageOrigin(t)
	const hashHex = "29181d224d982f659d1bd4703238d3f8489e26968ffe2a22bd786fa3bf15fef5" // host/approval request_test.go
	h, _ := hex.DecodeString(hashHex)
	challenge := base64.RawURLEncoding.EncodeToString(h)
	record := func(link, dl, name string) json.RawMessage {
		var c pageCeremony
		b.do(t, nil, "open", map[string]any{"url": link})
		b.do(t, &c, "ceremony", map[string]any{"download": dl})
		return json.RawMessage(bytes.TrimSuffix(saved(t, c, origin, dl, name), []byte("\n")))
	}
	for _, alg := range []struct{ name, inject string }{{"es256", ""}, {"eddsa", "eddsa"}} {
		b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": alg.inject})
		reg := record(ceremonyLink{V: 1, Op: "create", RP: "localhost", Origin: origin, Name: "mark-mac", Challenge: challenge}.link(), t.TempDir(), "world-registration-mark-mac.json")
		var r struct{ RawID string }
		if err := json.Unmarshal(reg, &r); err != nil {
			t.Fatal(err)
		}
		get := ceremonyLink{V: 1, Op: "get", RP: "localhost", Origin: origin, Request: "mr-fixture", Expires: 1767225600, Hash16: hashHex[:16], Challenge: challenge, Allow: []string{r.RawID}}.link()
		asr := record(get, t.TempDir(), "world-assertion-mr-fixture.json")
		var asr2 struct {
			Response struct{ AuthenticatorData string }
		}
		if err := json.Unmarshal(record(get, t.TempDir(), "world-assertion-mr-fixture.json"), &asr2); err != nil {
			t.Fatal(err)
		}
		a2, _ := json.Marshal(map[string]string{"authenticatorData": asr2.Response.AuthenticatorData})
		out := fmt.Sprintf("{\n  \"source\": %q,\n  \"chrome\": %q,\n  \"origin\": %q,\n  \"rp_id\": \"localhost\",\n  \"challenge_hex\": %q,\n  \"registration\": %s,\n  \"assertion\": %s,\n  \"assertion2\": %s\n}\n",
			"row 201 A2, generate.sh: Chrome headless, the passkey page, CDP virtual authenticator (a test key, D-WORLD-109); see SOURCE.txt", b.chrome, origin, hashHex, reg, asr, a2)
		if err := os.WriteFile(filepath.Join(dir, "chrome-"+alg.name+".json"), []byte(out), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// generate.sh refuses, writing nothing and running nothing, a call that does not name exactly one
// existing empty directory, and one without WORLD_CHROME. It needs no browser.
func TestPasskeyFixtureGeneratorRefuses(t *testing.T) {
	script, err := filepath.Abs(filepath.Join("..", "..", "host", "approval", "testdata", "webauthn", "generate.sh"))
	if err != nil {
		t.Fatal(err)
	}
	full := t.TempDir()
	if err := os.WriteFile(filepath.Join(full, "keep"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	empty := t.TempDir()
	for _, tc := range []struct {
		name, chrome, want string
		args               []string
	}{
		{"no-directory", "chrome", "usage: sh generate.sh <an empty directory>\n", nil},
		{"two-directories", "chrome", "usage: sh generate.sh <an empty directory>\n", []string{empty, empty}},
		{"missing", "chrome", "generate.sh: " + filepath.Join(empty, "nope") + " is not an empty directory\n", []string{filepath.Join(empty, "nope")}},
		{"not-empty", "chrome", "generate.sh: " + full + " is not an empty directory\n", []string{full}},
		{"no-chrome", "", "generate.sh: set WORLD_CHROME to Chrome's binary (only its virtual authenticator is used)\n", []string{empty}},
	} {
		cmd := exec.CommandContext(boundedTestContext(t), "sh", append([]string{script}, tc.args...)...)
		cmd.Env = append(os.Environ(), "WORLD_CHROME="+tc.chrome, "PATH=/usr/bin:/bin") // a fall-through finds no go and no node there
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		err := cmd.Run()
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 2 || stderr.String() != tc.want || stdout.Len() != 0 {
			t.Errorf("%s: %v, stderr %q, stdout %q; want exit 2 and %q", tc.name, err, stderr.String(), stdout.String(), tc.want)
		}
		if entries, _ := os.ReadDir(empty); len(entries) != 0 {
			t.Fatalf("%s wrote into %s", tc.name, empty)
		}
	}
}

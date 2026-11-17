package surfacepage

// Row 243 V3 sprint 2 M4 (design §2.5, §6 CP0): the confirmation card's data, over fixtures made by
// the real tools (approval.CanonicalVerb, broker.FinishPromptText, broker.FormatFinishApprovalScope).

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
)

// confirmFixture is one stored verb request: the request, its canonical bytes and the row.
type confirmFixture struct {
	name string
	req  approval.VerbRequest
	in   ConfirmationInput
}

// confirmFixtures are a draft finish with protected paths and the repository's committer, and a ready
// one creating its label whose title, body and argv are hostile (markup, quotes, a NUL).
func confirmFixtures(t *testing.T) []confirmFixture {
	t.Helper()
	const created = 1800000000
	var out []confirmFixture
	for i, c := range []struct {
		name, id, title, body, committer string
		argv, protected                  []string
		draft, createLabel               bool
	}{
		{"draft", "vr-0123456789abcdef", "Fix the bar", "## Summary\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7 is not in the prompt\n", "",
			[]string{"--preset", "p", "--title", "Fix the bar"}, []string{"host/x.go", "docs/a b.md"}, true, false},
		{"hostile", "vr-00000000000000ff", "<img src=x onerror=alert(1)> a\x00b & 'q'", "</pre><script>alert(2)</script>\nNUL \x00 here &amp; after — ünï\n", "Mark <m@x>",
			[]string{"--preset", "p", "--title", "it's <b>", "--ready"}, []string{}, false, true},
	} {
		scope := broker.FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/w1", SHA: strings.Repeat("ab", 20),
			Base: "dev", TitleSHA256: strings.Repeat("0", 64), BodySHA256: strings.Repeat("1", 64), Draft: c.draft, CreateLabel: c.createLabel, ExpiresAt: created + 900}
		r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: append([]string{"ailang-worldd", "episode", "finish", "w1"}, c.argv...),
			Episode: "w1", Prompt: broker.FinishPromptText(scope, c.title, c.body, c.protected, c.committer), Scope: broker.FormatFinishApprovalScope(scope),
			Title: c.title, Body: c.body, Committer: c.committer, Protected: c.protected}
		b, err := approval.CanonicalVerb(r)
		if err == nil {
			err = r.Validate() // the fixtures are admissible requests
		}
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		d, ch := sha256.Sum256(b), sha256.Sum256([]byte("challenge "+c.name))
		out = append(out, confirmFixture{c.name, r, ConfirmationInput{ID: c.id, Request: b, Digest: hex.EncodeToString(d[:]), Challenge: hex.EncodeToString(ch[:]),
			Created: created + int64(i), Expires: created + int64(i) + 600, Now: created + 5}})
	}
	return out
}

// CP0: the card carries the stored bytes (base64url), the digest and the raw challenge as stored, the
// times, the state from the row (its outcome, else expired from its expiry as Admit, else pending), the
// approver and decision only for an approved row, the passkeys grouped by origin and the verb namespace.
func TestConfirmationCard(t *testing.T) {
	f := confirmFixtures(t)[0]
	pk := []CardPasskey{{Origin: "http://localhost:1", RPID: "localhost", CredentialID: []byte{1}}, {Origin: "http://127.0.0.1:1", RPID: "127.0.0.1", CredentialID: []byte{2}},
		{Origin: "http://localhost:1", RPID: "localhost", CredentialID: []byte{3}}}
	c := ConfirmationCard(f.in, pk)
	if b, err := base64.RawURLEncoding.DecodeString(c.Request); err != nil || string(b) != string(f.in.Request) {
		t.Fatalf("request %q does not carry the stored bytes (%v)", c.Request, err)
	}
	if u := ConfirmationCard(ConfirmationInput{Request: []byte{0xfb, 0xff, 0xfe}}, nil).Request; u != "-__-" { // base64url, unpadded ("+//+" in base64)
		t.Fatalf("request %q is not base64url", u)
	}
	if c.V != 1 || c.ID != f.in.ID || c.Digest != f.in.Digest || c.Challenge != f.in.Challenge || c.Created != f.in.Created || c.Expires != f.in.Expires || c.Now != f.in.Now {
		t.Fatalf("card %+v, input %+v", c, f.in)
	}
	if want := []CardPasskeys{{Origin: "http://localhost:1", RP: "localhost", Allow: []string{"AQ", "Aw"}}, {Origin: "http://127.0.0.1:1", RP: "127.0.0.1", Allow: []string{"Ag"}}}; !reflect.DeepEqual(c.Passkeys, want) {
		t.Fatalf("passkeys %+v", c.Passkeys)
	}
	if c.SSH.Namespace != "world-verb@sunholo" { // approval.VerbNamespace, as a literal
		t.Fatalf("namespace %q", c.SSH.Namespace)
	}
	if none := ConfirmationCard(f.in, nil); none.Passkeys == nil || len(none.Passkeys) != 0 {
		t.Fatalf("no passkeys: %#v", none.Passkeys)
	}
	for _, tc := range []struct {
		name, outcome     string
		now               int64
		state, by, decide string
	}{
		{"pending", "", f.in.Expires - 1, ConfirmationPending, "", ""},
		{"expired at its expiry", "", f.in.Expires, ConfirmationExpired, "", ""},
		{"expired after", "", f.in.Expires + 100, ConfirmationExpired, "", ""},
		{"approved", ConfirmationApproved, f.in.Created + 1, ConfirmationApproved, "webauthn:mark", "sha256:" + strings.Repeat("cd", 32)},
		{"approved, read after its expiry", ConfirmationApproved, f.in.Expires + 1, ConfirmationApproved, "webauthn:mark", "sha256:" + strings.Repeat("cd", 32)},
		{"superseded", ConfirmationSuperseded, f.in.Created + 1, ConfirmationSuperseded, "", ""},
	} {
		in := f.in
		in.Outcome, in.Now, in.ApprovedBy, in.Decision = tc.outcome, tc.now, "webauthn:mark", "sha256:"+strings.Repeat("cd", 32) // set on every row: only approved shows them
		c := ConfirmationCard(in, nil)
		if c.State != tc.state || c.ApprovedBy != tc.by || c.Decision != tc.decide || c.Now != tc.now {
			t.Errorf("%s: state %q by %q decision %q now %d; want %q %q %q %d", tc.name, c.State, c.ApprovedBy, c.Decision, c.Now, tc.state, tc.by, tc.decide, tc.now)
		}
	}
}

// CP0: the wire is exactly design §2.5's keys in order; the input has no artifact; nothing on the wire
// is the message or its hash in any encoding.
func TestConfirmationCardWire(t *testing.T) {
	var names []string
	for i, ty := 0, reflect.TypeOf(ConfirmationInput{}); i < ty.NumField(); i++ {
		names = append(names, ty.Field(i).Name)
	}
	if want := []string{"ID", "Request", "Digest", "Challenge", "Created", "Expires", "Now", "Outcome", "ApprovedBy", "Decision"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("ConfirmationInput's fields %v, want %v (no artifact, no hash)", names, want)
	}
	for _, f := range confirmFixtures(t) {
		in := f.in
		in.Outcome, in.ApprovedBy, in.Decision = ConfirmationApproved, "ssh:mark", "sha256:"+strings.Repeat("ef", 32)
		b, err := json.Marshal(ConfirmationCard(in, []CardPasskey{{Origin: "http://localhost:1", RPID: "localhost", CredentialID: []byte{9}}}))
		if err != nil {
			t.Fatal(err)
		}
		var keys []string
		dec := json.NewDecoder(strings.NewReader(string(b)))
		_, _ = dec.Token()
		for dec.More() {
			k, _ := dec.Token()
			keys = append(keys, k.(string))
			var skip json.RawMessage
			_ = dec.Decode(&skip)
		}
		if want := []string{"v", "id", "request", "digest", "challenge", "created", "expires", "now", "state", "approvedBy", "decision", "passkeys", "ssh"}; !reflect.DeepEqual(keys, want) {
			t.Fatalf("%s: keys %v, want %v", f.name, keys, want)
		}
		var d, ch [32]byte
		hex.Decode(d[:], []byte(in.Digest))
		hex.Decode(ch[:], []byte(in.Challenge))
		m, _ := approval.VerbMessage(d, ch, in.Expires)
		h, _ := approval.VerbMessageHash(d, ch, in.Expires)
		for _, x := range [][]byte{m, h[:]} {
			for _, enc := range []string{hex.EncodeToString(x), base64.RawURLEncoding.EncodeToString(x), base64.StdEncoding.EncodeToString(x)} {
				if strings.Contains(string(b), enc) {
					t.Fatalf("%s: the wire carries the message or its hash (%s)", f.name, enc)
				}
			}
		}
		if !strings.Contains(string(b), `"state":"approved","approvedBy":"ssh:mark","decision":"sha256:`) || !strings.Contains(string(b), `"ssh":{"namespace":"world-verb@sunholo"}`) {
			t.Fatalf("%s: %s", f.name, b)
		}
	}
}

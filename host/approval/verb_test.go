package approval_test

// Row 243 V3 sprint 1 M1a (design §2.2, §6 VL1-VL3): the verb request's canonical bytes, its strict
// decoder, every Validate rule, and the 86-byte message's layout.

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// vrRequest is a valid verb request: a draft finish of w1 whose title and body need no escaping
// but hold the bytes CanonicalVerb must not HTML-escape.
func vrRequest() approval.VerbRequest {
	return approval.VerbRequest{
		Schema:    approval.VerbRequestSchema,
		Verb:      approval.VerbFinish,
		Argv:      []string{"ailang-worldd", "episode", "finish", "w1", "--preset", "p", "--confirm-on-surface"},
		Episode:   "w1",
		Prompt:    "committer: (repository git config)\nrepository: o/r\n",
		Scope:     "github-pr:o/r#remote=origin&head=world/w1",
		Title:     "Fix <b> & tidy",
		Body:      "line 1\nline <2> & \"3\"\n",
		Committer: "",
		Protected: []string{"a/b", "c"},
	}
}

// VL1: the canonical bytes are compact, in field order, with HTML escaping off and no trailing
// newline; DecodeVerb inverts them exactly.
func TestVerbRequestRoundTrip(t *testing.T) {
	r := vrRequest()
	c, err := approval.CanonicalVerb(r)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"schema":"world/verb-request/v1","verb":"episode finish","argv":["ailang-worldd","episode","finish","w1","--preset","p","--confirm-on-surface"],` +
		`"episode":"w1","prompt":"committer: (repository git config)\nrepository: o/r\n","scope":"github-pr:o/r#remote=origin&head=world/w1",` +
		`"title":"Fix <b> & tidy","body":"line 1\nline <2> & \"3\"\n","committer":"","protected":["a/b","c"]}`
	if string(c) != want {
		t.Fatalf("CanonicalVerb =\n%s\nwant\n%s", c, want)
	}
	got, err := approval.DecodeVerb(c)
	if err != nil {
		t.Fatal(err)
	}
	c2, _ := approval.CanonicalVerb(got)
	if !bytes.Equal(c, c2) || got.Title != r.Title || len(got.Argv) != len(r.Argv) || got.Protected[1] != "c" {
		t.Fatalf("round trip: %+v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("the valid request: %v", err)
	}
	nilProtected := vrRequest()
	nilProtected.Protected = nil
	c, _ = approval.CanonicalVerb(nilProtected)
	if back, err := approval.DecodeVerb(c); err != nil || back.Protected != nil || !strings.HasSuffix(string(c), `"protected":null}`) {
		t.Fatalf("nil protected: %v %q", err, c)
	}
}

// vrEscaped is "<b>" as an HTML-escaping encoder spells it (backslash u003c b backslash u003e).
var vrEscaped = string([]byte{'\\', 'u', '0', '0', '3', 'c', 'b', '\\', 'u', '0', '0', '3', 'e'})

// VL1: unknown fields are ErrDecode; trailing data, reordered keys, whitespace and any other
// non-canonical spelling are ErrNotCanonical.
func TestVerbRequestDecodeRefusals(t *testing.T) {
	c, err := approval.CanonicalVerb(vrRequest())
	if err != nil {
		t.Fatal(err)
	}
	s := string(c)
	for _, tc := range []struct {
		name string
		in   string
		want error
	}{
		{"unknown field", strings.Replace(s, `{"schema"`, `{"typed":"y","schema"`, 1), approval.ErrDecode},
		{"unknown field at the end", strings.TrimSuffix(s, "}") + `,"y":true}`, approval.ErrDecode},
		{"not JSON", "y", approval.ErrDecode},
		{"empty", "", approval.ErrDecode},
		{"trailing data", s + "x", approval.ErrNotCanonical},
		{"trailing newline", s + "\n", approval.ErrNotCanonical},
		{"a second object", s + s, approval.ErrNotCanonical},
		{"reordered keys", `{"verb":"episode finish","schema":"world/verb-request/v1"` + s[len(`{"schema":"world/verb-request/v1","verb":"episode finish"`):], approval.ErrNotCanonical},
		{"whitespace", strings.Replace(s, `"schema":`, `"schema": `, 1), approval.ErrNotCanonical},
		{"html-escaped", strings.Replace(s, "<b>", vrEscaped, 1), approval.ErrNotCanonical},
	} {
		if r, err := approval.DecodeVerb([]byte(tc.in)); !errors.Is(err, tc.want) || r.Schema != "" {
			t.Errorf("%s: (%q, %v), want %v", tc.name, r.Schema, err, tc.want)
		}
	}
}

// VL2: one case per Validate rule (design §2.2), each with its boundary where it has one.
func TestVerbRequestValidate(t *testing.T) {
	argv := func(n int) []string {
		a := []string{"ailang-worldd", "episode", "finish", "w1"}
		for len(a) < n {
			a = append(a, "--allow-failing")
		}
		return a
	}
	for _, tc := range []struct {
		name string
		edit func(r *approval.VerbRequest)
		want error // nil: admissible
		text string
	}{
		{"V1 schema", func(r *approval.VerbRequest) { r.Schema = "world/mint-request/v1" }, approval.ErrSchema, "schema"},
		{"V2 utf-8 title", func(r *approval.VerbRequest) { r.Title = "a\xffb" }, approval.ErrString, ": title"},
		{"V2 utf-8 protected", func(r *approval.VerbRequest) { r.Protected = []string{"ok", "\xc3"} }, approval.ErrString, "protected[1]"},
		{"V3 verb", func(r *approval.VerbRequest) { r.Verb = "episode start" }, approval.ErrVerbRequest, "verb"},
		{"V4 episode", func(r *approval.VerbRequest) { r.Episode, r.Argv[3] = "W1", "W1" }, approval.ErrEpisode, "episode"},
		{"V5 argv short", func(r *approval.VerbRequest) { r.Argv = r.Argv[:3] }, approval.ErrVerbRequest, "argv does not start"},
		{"V5 argv[0]", func(r *approval.VerbRequest) { r.Argv[0] = "world" }, approval.ErrVerbRequest, "argv does not start"},
		{"V5 argv[1]", func(r *approval.VerbRequest) { r.Argv[1] = "session" }, approval.ErrVerbRequest, "argv does not start"},
		{"V5 argv[2]", func(r *approval.VerbRequest) { r.Argv[2] = "start" }, approval.ErrVerbRequest, "argv does not start"},
		{"V5 argv[3]", func(r *approval.VerbRequest) { r.Argv[3] = "w2" }, approval.ErrVerbRequest, "argv does not start"},
		{"V5 argv of exactly four", func(r *approval.VerbRequest) { r.Argv = r.Argv[:4] }, nil, ""},
		{"V6 64 words", func(r *approval.VerbRequest) { r.Argv = argv(64) }, nil, ""},
		{"V6 65 words", func(r *approval.VerbRequest) { r.Argv = argv(65) }, approval.ErrVerbRequest, "65 words"},
		{"V7 newline", func(r *approval.VerbRequest) { r.Argv[5] = "p\nq" }, approval.ErrVerbRequest, "argv[5] holds a control byte"},
		{"V7 0x1f", func(r *approval.VerbRequest) { r.Argv[4] = "\x1f" }, approval.ErrVerbRequest, "argv[4] holds a control byte"},
		{"V7 DEL", func(r *approval.VerbRequest) { r.Argv[6] = "a\x7f" }, approval.ErrVerbRequest, "argv[6] holds a control byte"},
		{"V7 space and tilde", func(r *approval.VerbRequest) { r.Argv[5] = " ~p" }, nil, ""},
		{"V8 prompt", func(r *approval.VerbRequest) { r.Prompt = "" }, approval.ErrVerbRequest, "prompt is empty"},
		{"V8 scope", func(r *approval.VerbRequest) { r.Scope = "" }, approval.ErrVerbRequest, "scope is empty"},
		{"V8 title", func(r *approval.VerbRequest) { r.Title = "" }, approval.ErrVerbRequest, "title is empty"},
		{"V8 committer may be empty", func(r *approval.VerbRequest) { r.Committer, r.Body, r.Protected = "", "", nil }, nil, ""},
		{"V9 body 60000", func(r *approval.VerbRequest) { r.Body = strings.Repeat("b", 60000) }, nil, ""},
		{"V9 body 60001", func(r *approval.VerbRequest) { r.Body = strings.Repeat("b", 60001) }, approval.ErrVerbRequest, "body is 60001 bytes"},
		{"V10 canonical 131072", func(r *approval.VerbRequest) { r.Prompt = strings.Repeat("p", 131072-vrLen(t)+1) }, nil, ""},
		{"V10 canonical 131073", func(r *approval.VerbRequest) { r.Prompt = strings.Repeat("p", 131073-vrLen(t)+1) }, approval.ErrVerbRequest, "131073 bytes"},
	} {
		r := vrRequest()
		tc.edit(&r)
		err := r.Validate()
		if tc.want == nil {
			if err != nil {
				t.Errorf("%s: %v, want admissible", tc.name, err)
			}
			continue
		}
		if !errors.Is(err, tc.want) || !strings.Contains(err.Error(), tc.text) {
			t.Errorf("%s: %v, want %v naming %q", tc.name, err, tc.want, tc.text)
		}
	}
}

// vrLen is the canonical length of vrRequest() with the one-byte prompt "p".
func vrLen(t *testing.T) int {
	t.Helper()
	r := vrRequest()
	r.Prompt = "p"
	c, err := approval.CanonicalVerb(r)
	if err != nil {
		t.Fatal(err)
	}
	return len(c)
}

// VL3: the message is 86 bytes in design §2.2's layout; its hash is the SHA-256 of those bytes; a
// negative expiry has no message; the three prefixes differ at byte 6 and none prefixes another.
func TestVerbMessageLayout(t *testing.T) {
	if approval.VerbPrefix != "world-verb/v1" || approval.VerbNamespace != "world-verb@sunholo" || approval.VerbMessageLen != 86 ||
		approval.VerbRequestSchema != "world/verb-request/v1" || approval.MaxVerbRequestBytes != 131072 || approval.VerbRequestLifetime.Seconds() != 600 {
		t.Fatalf("constants: %q %q %d %q %d %v", approval.VerbPrefix, approval.VerbNamespace, approval.VerbMessageLen,
			approval.VerbRequestSchema, approval.MaxVerbRequestBytes, approval.VerbRequestLifetime)
	}
	var d, c [32]byte
	for i := range d {
		d[i], c[i] = byte(i), byte(0x80+i)
	}
	const expiry = int64(0x0102030405060708)
	m, err := approval.VerbMessage(d, c, expiry)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 86 || string(m[:13]) != "world-verb/v1" || m[13] != 0 || !bytes.Equal(m[14:46], d[:]) || !bytes.Equal(m[46:78], c[:]) ||
		binary.BigEndian.Uint64(m[78:]) != uint64(expiry) || m[78] != 0x01 || m[85] != 0x08 {
		t.Fatalf("VerbMessage = %x", m)
	}
	h, err := approval.VerbMessageHash(d, c, expiry)
	if err != nil || h != sha256.Sum256(m) {
		t.Fatalf("VerbMessageHash = %x, %v", h, err)
	}
	if m, err := approval.VerbMessage(d, c, -1); m != nil || !errors.Is(err, approval.ErrNegativeExpiry) {
		t.Fatalf("negative expiry: %x, %v", m, err)
	}
	if h, err := approval.VerbMessageHash(d, c, -1); h != ([32]byte{}) || !errors.Is(err, approval.ErrNegativeExpiry) {
		t.Fatalf("negative expiry hash: %x, %v", h, err)
	}
	if m, err := approval.VerbMessage(d, c, 0); err != nil || len(m) != 86 {
		t.Fatalf("expiry 0: %v", err)
	}
	mint, _ := approval.SignedMessage(d, c, expiry)
	decision, _ := approval.DecisionMessage(approval.VerbApprove, d, c, expiry)
	if approval.VerbPrefix[:6] != approval.MessagePrefix[:6] || approval.VerbPrefix[6] == approval.MessagePrefix[6] || approval.VerbPrefix[6] == approval.DecisionPrefix[6] {
		t.Fatalf("prefixes: %q %q %q", approval.VerbPrefix, approval.MessagePrefix, approval.DecisionPrefix)
	}
	for name, other := range map[string][]byte{"mint": mint, "decision": decision} {
		if bytes.HasPrefix(m, other[:len(other)-40]) || bytes.HasPrefix(other, m[:20]) || bytes.Equal(m, other) {
			t.Errorf("the verb message and the %s message share a prefix", name)
		}
	}
	if approval.VerbNamespace == approval.Namespace || approval.VerbNamespace == approval.DecisionNamespace {
		t.Fatalf("namespaces: %q %q %q", approval.VerbNamespace, approval.Namespace, approval.DecisionNamespace)
	}
}

package approval

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// approverFixtureSHA256 pins the four generated files of testdata/approvers (AP0). Re-run
// generate.sh only with a re-pin in the same commit.
var approverFixtureSHA256 = map[string]string{
	"SOURCE.txt":  "7a2ed3a96e88820e3a86076b0db4ce8b4f0f858e842b3a02e37c578720c39fe9",
	"generate.sh": "54d00a3e9177751498b8a467929dccb68c6e942f6a9f4b201e0ccc91d06036ab",
	"p384.pub":    "a762526789a43e8f8710d953f3b7312266abe094643c03bbb386656aff65d8b2",
	"rsa3072.pub": "38532b1868366820cc7b8092af1514aaf5c01843c39be7080396dbff2e14d0ec",
}

const (
	fixtureEd25519Fingerprint = "SHA256:pHMYJm+CndqyZGrzMHQQBH0tNeNGj1U9sy0JrVsYdSQ"
	fixtureP256Fingerprint    = "SHA256:yMwGgMZdq1aVNYRG/HO2T69s/Fb14yPZW0GBWtwRyEo"
)

// approverFixtureBytes reads one file of testdata/approvers.
func approverFixtureBytes(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "approvers", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// paddedKeyLine is ed25519.pub's line (no newline) + " " + 'x' padding, exactly n bytes long.
// t.Fatal when n is too small to hold the line and one pad byte.
func paddedKeyLine(t *testing.T, n int) []byte {
	t.Helper()
	line := bytes.TrimSuffix(fixtureBytes(t, "ed25519.pub"), []byte("\n"))
	if n < len(line)+2 {
		t.Fatalf("n=%d is too small for the %d-byte line and a pad byte", n, len(line))
	}
	out := bytes.Clone(line)
	out = append(out, ' ')
	return append(out, bytes.Repeat([]byte("x"), n-len(line)-1)...)
}

// skRetype is ed25519.pub with its key type rewritten to sk-ssh-ed25519@openssh.com in the line
// and in the blob (appendString(nil, "sk-ssh-ed25519@openssh.com") + the blob after its type string).
func skRetype(t *testing.T) []byte {
	t.Helper()
	line := string(bytes.TrimSuffix(fixtureBytes(t, "ed25519.pub"), []byte("\n")))
	f := strings.Fields(line)
	if len(f) < 2 || f[0] != "ssh-ed25519" {
		t.Fatalf("ed25519.pub is %q, want an ssh-ed25519 line", line)
	}
	blob, err := base64.StdEncoding.DecodeString(f[1])
	if err != nil {
		t.Fatal(err)
	}
	r := &sshReader{b: blob}
	typ, err := r.readString()
	if err != nil || string(typ) != "ssh-ed25519" {
		t.Fatalf("blob type %q, %v; want ssh-ed25519", typ, err)
	}
	sk := appendString(nil, []byte("sk-ssh-ed25519@openssh.com"))
	sk = append(sk, r.b...)
	out := "sk-ssh-ed25519@openssh.com " + base64.StdEncoding.EncodeToString(sk)
	for _, extra := range f[2:] {
		out += " " + extra
	}
	return []byte(out + "\n")
}

// derivedPrivateKeyFile is an OpenSSH private-key armor around derived noise (no real key):
// "-----BEGIN OPENSSH PRIVATE KEY-----\n" + noise + "\n-----END OPENSSH PRIVATE KEY-----\n", where
// s := sha512.Sum512([]byte("approver-test-noise")); noise := base64.StdEncoding.EncodeToString(bytes.Repeat(s[:], 4)).
func derivedPrivateKeyFile() (file []byte, noise string) {
	s := sha512.Sum512([]byte("approver-test-noise"))
	noise = base64.StdEncoding.EncodeToString(bytes.Repeat(s[:], 4))
	return []byte("-----BEGIN OPENSSH PRIVATE KEY-----\n" + noise + "\n-----END OPENSSH PRIVATE KEY-----\n"), noise
}

// AP0: the generated fixtures are exactly the pinned bytes; no NUL; text; the two key types.
func TestApproverFixturesAreTheGeneratedBytes(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("testdata", "approvers"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != len(approverFixtureSHA256) {
		t.Errorf("testdata/approvers holds %d entries, want %d", len(entries), len(approverFixtureSHA256))
	}
	for _, e := range entries {
		if _, ok := approverFixtureSHA256[e.Name()]; !ok {
			t.Errorf("unexpected fixture %s", e.Name())
		}
	}
	for name, want := range approverFixtureSHA256 {
		b := approverFixtureBytes(t, name)
		sum := sha256.Sum256(b)
		if got := hex.EncodeToString(sum[:]); got != want {
			t.Errorf("%s sha256 %s, want %s", name, got, want)
		}
		if bytes.IndexByte(b, 0) >= 0 {
			t.Errorf("%s contains a NUL byte", name)
		}
		if !bytes.HasSuffix(b, []byte("\n")) {
			t.Errorf("%s does not end with a newline", name)
		}
	}
	if !bytes.HasPrefix(approverFixtureBytes(t, "rsa3072.pub"), []byte("ssh-rsa ")) {
		t.Errorf("rsa3072.pub is not an ssh-rsa line")
	}
	if !bytes.HasPrefix(approverFixtureBytes(t, "p384.pub"), []byte("ecdsa-sha2-nistp384 ")) {
		t.Errorf("p384.pub is not an ecdsa-sha2-nistp384 line")
	}
}

// AP1: every refusal of the enrolment law, in its fixed E1-E7 order, returns the typed error and
// the zero Approver.
func TestNewApproverRefusals(t *testing.T) {
	const now = 1000
	good := fixtureBytes(t, "ed25519.pub")
	file, noise := derivedPrivateKeyFile()
	type refusal struct {
		name string
		call func() (Approver, error)
		want error
		must string // the error text must contain this, if set
		not  string // the error text must not contain this, if set
	}
	ssh := func(name string, key []byte) func() (Approver, error) {
		return func() (Approver, error) { return NewApprover(name, "ssh", key, now) }
	}
	cases := []refusal{
		{name: "name-empty", call: ssh("", good), want: ErrApproverName},
		{name: "name-upper", call: ssh("Mark", good), want: ErrApproverName},
		{name: "name-space", call: ssh("a b", good), want: ErrApproverName},
		{name: "name-dash-first", call: ssh("-x", good), want: ErrApproverName},
		{name: "name-dot-first", call: ssh(".x", good), want: ErrApproverName},
		{name: "name-49", call: ssh(strings.Repeat("a", 49), good), want: ErrApproverName},
		{name: "name-escape", call: ssh("x\x1b[2J", good), want: ErrApproverName},
		{name: "name-newline", call: ssh("mark\n", good), want: ErrApproverName},
		{name: "webauthn", call: func() (Approver, error) { return NewApprover("mark", "webauthn", good, now) }, want: ErrFormatNotYet},
		{name: "webauthn-empty-key", call: func() (Approver, error) { return NewApprover("mark", "webauthn", nil, now) }, want: ErrFormatNotYet},
		{name: "format-gpg", call: func() (Approver, error) { return NewApprover("mark", "gpg", good, now) }, want: ErrApproverFormat},
		{name: "format-empty", call: func() (Approver, error) { return NewApprover("mark", "", good, now) }, want: ErrApproverFormat},
		{name: "format-SSH", call: func() (Approver, error) { return NewApprover("mark", "SSH", good, now) }, want: ErrApproverFormat},
		{name: "key-nil", call: ssh("mark", nil), want: ErrApproverKeyFile},
		{name: "key-empty", call: ssh("mark", []byte{}), want: ErrApproverKeyFile},
		{name: "key-8193", call: ssh("mark", paddedKeyLine(t, 8193)), want: ErrApproverKeyFile},
		{name: "private", call: ssh("mark", file), want: ErrPrivateKeyFile, not: noise[:24]},
		{name: "private-leading-space", call: ssh("mark", append([]byte("\n  \t"), file...)), want: ErrPrivateKeyFile},
		{name: "two-lines", call: ssh("mark", append(bytes.Clone(good), good...)), want: ErrApproverKeyFile},
		{name: "blank-second-line", call: ssh("mark", append(bytes.Clone(good), '\n')), want: ErrApproverKeyFile},
		{name: "rsa", call: ssh("mark", approverFixtureBytes(t, "rsa3072.pub")), want: ErrKeyType},
		{name: "p384", call: ssh("mark", approverFixtureBytes(t, "p384.pub")), want: ErrKeyType},
		{name: "sk", call: ssh("mark", skRetype(t)), want: ErrKeyType, must: "row 201 Q1"},
		{name: "one-field", call: ssh("mark", []byte("ssh-ed25519\n")), want: ErrAuthorizedKey},
		{name: "bad-base64", call: ssh("mark", []byte("ssh-ed25519 !!!! c\n")), want: ErrBase64},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, err := tc.call()
			if !errors.Is(err, tc.want) {
				t.Fatalf("err %v, want %v", err, tc.want)
			}
			if !reflect.DeepEqual(a, Approver{}) {
				t.Errorf("returned %+v with the refusal", a)
			}
			if tc.must != "" && !strings.Contains(err.Error(), tc.must) {
				t.Errorf("error %q lacks %q", err, tc.must)
			}
			if tc.not != "" && strings.Contains(err.Error(), tc.not) {
				t.Errorf("error %q echoes key material %q", err, tc.not)
			}
		})
	}
}

// AP2: the real fixtures enrol, with or without their trailing newline; the 8192-byte line is the
// largest accepted.
func TestNewApproverAcceptsTheRealFixtures(t *testing.T) {
	const now = 1234
	for _, tc := range []struct{ file, typ, fp string }{
		{"ed25519.pub", "ssh-ed25519", fixtureEd25519Fingerprint},
		{"p256.pub", "ecdsa-sha2-nistp256", fixtureP256Fingerprint},
	} {
		raw := fixtureBytes(t, tc.file)
		for _, b := range [][]byte{raw, bytes.TrimSuffix(raw, []byte("\n"))} {
			a, err := NewApprover("mark-mac", "ssh", b, now)
			if err != nil {
				t.Fatalf("%s: %v", tc.file, err)
			}
			want, err := ParseAuthorizedKey(b)
			if err != nil {
				t.Fatal(err)
			}
			if a.Name != "mark-mac" || a.Format != FormatSSH || a.Key.Type != tc.typ || !bytes.Equal(a.Key.Blob, want.Blob) {
				t.Errorf("%s: approver %+v", tc.file, a)
			}
			if got := a.Key.Fingerprint(); got != tc.fp {
				t.Errorf("%s: fingerprint %s, want %s", tc.file, got, tc.fp)
			}
			if a.EnrolledAt != now || a.LastUsedAt != 0 || a.RevokedAt != 0 {
				t.Errorf("%s: times enrolled %d used %d revoked %d", tc.file, a.EnrolledAt, a.LastUsedAt, a.RevokedAt)
			}
		}
	}
	pad := paddedKeyLine(t, MaxKeyFileBytes)
	if len(pad) != 8192 {
		t.Fatalf("boundary line is %d bytes, want 8192", len(pad))
	}
	if _, err := NewApprover("mark-mac", "ssh", pad, now); err != nil {
		t.Errorf("8192-byte line refused: %v", err)
	}
}

// AP3: the revocation law. A revoked approver is refused whatever the clock says; a clock behind
// the enrolment time is refused.
func TestCheckRevoke(t *testing.T) {
	a := Approver{Name: "x", Format: "ssh", EnrolledAt: 1000}
	revoked := a
	revoked.RevokedAt = 1500
	for _, tc := range []struct {
		name string
		a    Approver
		now  int64
		want error
	}{
		{"active-at-enrolment", a, 1000, nil},
		{"active-later", a, 2000, nil},
		{"behind-enrolment", a, 999, ErrApproverClockBehind},
		{"revoked", revoked, 2000, ErrApproverRevoked},
		{"revoked-and-behind", revoked, 999, ErrApproverRevoked},
	} {
		err := CheckRevoke(tc.a, tc.now)
		if tc.want == nil && err != nil {
			t.Errorf("%s: %v, want nil", tc.name, err)
		}
		if tc.want != nil && !errors.Is(err, tc.want) {
			t.Errorf("%s: %v, want %v", tc.name, err, tc.want)
		}
	}
}

// AP4: the name pattern's boundaries: 1 to 48 characters, a letter or digit first, and the
// anchors (a newline or a space is refused).
func TestApproverNamePatternBoundaries(t *testing.T) {
	for _, ok := range []string{"a", "0", strings.Repeat("a", 48), "a.b_c-d", "mark-mac", "a-", "9."} {
		if !ApproverNamePattern.MatchString(ok) {
			t.Errorf("%q refused, want accepted", ok)
		}
	}
	for _, bad := range []string{strings.Repeat("a", 49), "", ".a", "_a", "-a", "A", "a/b", "a:b", "a b", "ä", "a\n"} {
		if ApproverNamePattern.MatchString(bad) {
			t.Errorf("%q accepted, want refused", bad)
		}
	}
}

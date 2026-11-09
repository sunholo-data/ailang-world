package approval

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// fixtureMessageHex is SignedMessage(fixtureDigest, fixtureChallenge(), fixtureExpiry):
// the 90 bytes every SSHSIG fixture signs. testdata/sshsig/generate.sh reads this line.
const fixtureMessageHex = "776f726c642d617070726f76616c2f7631000ace927ef952f541b1e454b066f6d6b1941fb5f3d710c784c22dd80a53d771eba0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf000000006955b900"

// fixtureExpiry is 2026-01-01T00:00:00Z.
const fixtureExpiry int64 = 1767225600

func fixtureChallenge() [32]byte {
	var c [32]byte
	for i := range c {
		c[i] = byte(0xa0 + i)
	}
	return c
}

// The SSHSIG fixtures under testdata/sshsig/ are produced ONLY by
// testdata/sshsig/generate.sh, with real ssh-keygen, signing from stdin with
// throwaway keys made in a mktemp -d directory that it deletes:
//
//	ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/ed25519" < "$tmp/message.bin" > "$here/ed25519.sig"
//	ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/p256" < "$tmp/message.bin" > "$here/p256.sig"
//	ssh-keygen -Y sign -n git -f "$tmp/ed25519" < "$tmp/message.bin" > "$here/git-ns-ed25519.sig"
//
// SOURCE.txt records the tool version. Whoever runs generate.sh updates this
// table and SOURCE.txt in the same commit (design §2.8, test SG1).
var fixtureSHA256 = map[string]string{
	"SOURCE.txt":         "044d6fb53e14c4512d2232b333239b97bab738368086b92367eb18c4f8b2017d",
	"ed25519.pub":        "1dd2626f4e0986338729a5e7cbe3bb53c1e77530b0973f9569d5a96a7c6502d3",
	"ed25519.sig":        "4376f4f23154b889dcbb955c59d1c1fdd3e5f5df5932573bf3c02409413213be",
	"generate.sh":        "069b19d659c7357d3c29a53821fa086c88cbf999a84c8395a557d2be51723e5e",
	"git-ns-ed25519.sig": "88be2bf82dfa0d6efb447d80ccab332461eb005c0afc66b93953aa24d14a3137",
	"message.hex":        "77bd21faa6199f457f57ed1ecf310207e5c4af9ab3a768f8051cad58427d8ad3",
	"p256.pub":           "f109a7ec4d8d871aec2fe5ca3a78f5f2dbd27a71ffab682614f6c3f3d1d12dde",
	"p256.sig":           "265e0ef0578aa36fa315c811a99a829f744067e3f49a28b9aad5366c89d9b623",
}

// fixtureBytes reads one file of testdata/sshsig.
func fixtureBytes(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "sshsig", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// fixtureMessage is the 90-byte message every fixture signs.
func fixtureMessage(t *testing.T) []byte {
	t.Helper()
	b, err := hex.DecodeString(fixtureMessageHex)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// fixtureKey parses one .pub fixture with the production parser.
func fixtureKey(t *testing.T, name string) SSHPublicKey {
	t.Helper()
	k, err := ParseAuthorizedKey(fixtureBytes(t, name))
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// unarmorTest is a test-side decoder, independent of ParseSSHSIG.
func unarmorTest(t *testing.T, armored []byte) []byte {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(string(armored), "\n"), "\n")
	if len(lines) < 3 || lines[0] != "-----BEGIN SSH SIGNATURE-----" || lines[len(lines)-1] != "-----END SSH SIGNATURE-----" {
		t.Fatal("unarmorTest: not an SSHSIG armor")
	}
	b, err := base64.StdEncoding.DecodeString(strings.Join(lines[1:len(lines)-1], ""))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// armorTest wraps blob in SSHSIG armor: base64 at 70 columns, one trailing newline.
func armorTest(blob []byte) []byte {
	enc := base64.StdEncoding.EncodeToString(blob)
	var out strings.Builder
	out.WriteString("-----BEGIN SSH SIGNATURE-----\n")
	for len(enc) > 0 {
		n := min(70, len(enc))
		out.WriteString(enc[:n] + "\n")
		enc = enc[n:]
	}
	out.WriteString("-----END SSH SIGNATURE-----\n")
	return []byte(out.String())
}

// sshStringTest is an SSH string: uint32 BE length followed by b.
func sshStringTest(b []byte) []byte {
	return append(binary.BigEndian.AppendUint32(nil, uint32(len(b))), b...)
}

// packStrings concatenates SSH strings.
func packStrings(fields ...[]byte) []byte {
	var out []byte
	for _, f := range fields {
		out = append(out, sshStringTest(f)...)
	}
	return out
}

// readStringsTest reads exactly n SSH strings; t.Fatal on a short or trailing input.
func readStringsTest(t *testing.T, b []byte, n int) [][]byte {
	t.Helper()
	out := make([][]byte, 0, n)
	for i := 0; i < n; i++ {
		if len(b) < 4 {
			t.Fatalf("readStringsTest: input ends before field %d", i)
		}
		l := int(binary.BigEndian.Uint32(b[:4]))
		b = b[4:]
		if l > len(b) {
			t.Fatalf("readStringsTest: field %d claims %d bytes, %d remain", i, l, len(b))
		}
		out = append(out, b[:l:l])
		b = b[l:]
	}
	if len(b) != 0 {
		t.Fatalf("readStringsTest: %d trailing bytes", len(b))
	}
	return out
}

// firstStringTest splits off the first SSH string of b.
func firstStringTest(t *testing.T, b []byte) (first, rest []byte) {
	t.Helper()
	if len(b) < 4 {
		t.Fatal("firstStringTest: short input")
	}
	l := int(binary.BigEndian.Uint32(b[:4]))
	if l > len(b)-4 {
		t.Fatal("firstStringTest: length exceeds input")
	}
	return b[4 : 4+l : 4+l], b[4+l:]
}

// sigParts is an SSHSIG blob split into its fields.
type sigParts struct {
	magic                                              []byte
	version                                            uint32
	publicKey, namespace, reserved, hashAlg, signature []byte
}

func splitSig(t *testing.T, armored []byte) sigParts {
	t.Helper()
	blob := unarmorTest(t, armored)
	if len(blob) < 10 {
		t.Fatal("splitSig: blob shorter than magic and version")
	}
	f := readStringsTest(t, blob[10:], 5)
	return sigParts{
		magic:     blob[:6],
		version:   binary.BigEndian.Uint32(blob[6:10]),
		publicKey: f[0], namespace: f[1], reserved: f[2], hashAlg: f[3], signature: f[4],
	}
}

// blob re-encodes the parts in order, without armor.
func (p sigParts) blob() []byte {
	b := append([]byte(nil), p.magic...)
	b = binary.BigEndian.AppendUint32(b, p.version)
	for _, f := range [][]byte{p.publicKey, p.namespace, p.reserved, p.hashAlg, p.signature} {
		b = append(b, sshStringTest(f)...)
	}
	return b
}

// armored re-encodes the parts in order, then armors them.
func (p sigParts) armored() []byte { return armorTest(p.blob()) }

// flipByte returns a copy of b with b[i] XOR 1.
func flipByte(b []byte, i int) []byte {
	c := bytes.Clone(b)
	c[i] ^= 0x01
	return c
}

// highBitPadded reports whether an mpint starts with 0x00 before a byte >= 0x80 (D7).
func highBitPadded(b []byte) bool {
	return len(b) > 1 && b[0] == 0 && b[1] >= 0x80
}

// withVersion rewrites the SSHSIG version field.
func withVersion(t *testing.T, a []byte, v uint32) []byte {
	t.Helper()
	p := splitSig(t, a)
	p.version = v
	return p.armored()
}

// withMagic rewrites the 6-byte magic.
func withMagic(t *testing.T, a []byte, m string) []byte {
	t.Helper()
	if len(m) != 6 {
		t.Fatalf("withMagic: %q is not 6 bytes", m)
	}
	p := splitSig(t, a)
	p.magic = []byte(m)
	return p.armored()
}

// withReserved rewrites the reserved field.
func withReserved(t *testing.T, a []byte, r string) []byte {
	t.Helper()
	p := splitSig(t, a)
	p.reserved = []byte(r)
	return p.armored()
}

// withHashAlg rewrites the hash algorithm field.
func withHashAlg(t *testing.T, a []byte, alg string) []byte {
	t.Helper()
	p := splitSig(t, a)
	p.hashAlg = []byte(alg)
	return p.armored()
}

// withTrailingByte appends one 0x00 after the signature string, inside the blob.
func withTrailingByte(t *testing.T, a []byte) []byte {
	t.Helper()
	p := splitSig(t, a)
	return armorTest(append(p.blob(), 0))
}

// withKeyType rewrites the FIRST string of the publickey blob and keeps the rest.
func withKeyType(t *testing.T, a []byte, typ string) []byte {
	t.Helper()
	p := splitSig(t, a)
	_, rest := firstStringTest(t, p.publicKey)
	p.publicKey = append(sshStringTest([]byte(typ)), rest...)
	return p.armored()
}

// withCurveName rewrites the second string (the curve) of a p256 publickey blob.
func withCurveName(t *testing.T, a []byte, curve string) []byte {
	t.Helper()
	p := splitSig(t, a)
	f := readStringsTest(t, p.publicKey, 3)
	p.publicKey = packStrings(f[0], []byte(curve), f[2])
	return p.armored()
}

// withPointByteFlipped flips Q[33] of a p256 publickey blob, moving the point off the curve.
func withPointByteFlipped(t *testing.T, a []byte) []byte {
	t.Helper()
	p := splitSig(t, a)
	f := readStringsTest(t, p.publicKey, 3)
	p.publicKey = packStrings(f[0], f[1], flipByte(f[2], 33))
	return p.armored()
}

// withShortEd25519Key truncates the ed25519 key to 31 bytes.
func withShortEd25519Key(t *testing.T, a []byte) []byte {
	t.Helper()
	p := splitSig(t, a)
	f := readStringsTest(t, p.publicKey, 2)
	p.publicKey = packStrings(f[0], f[1][:31])
	return p.armored()
}

// withSigType rewrites the sig-format-id inside the signature string.
func withSigType(t *testing.T, a []byte, typ string) []byte {
	t.Helper()
	p := splitSig(t, a)
	f := readStringsTest(t, p.signature, 2)
	p.signature = packStrings([]byte(typ), f[1])
	return p.armored()
}

// withLastSigByteFlipped flips the last byte of the inner signature blob.
func withLastSigByteFlipped(t *testing.T, a []byte) []byte {
	t.Helper()
	p := splitSig(t, a)
	f := readStringsTest(t, p.signature, 2)
	p.signature = packStrings(f[0], flipByte(f[1], len(f[1])-1))
	return p.armored()
}

// withNegativeMpint strips the real leading 0x00 from the first of r, s that has a high-bit
// value, so that the value is now a negative mpint. t.Fatal if neither r nor s has one.
func withNegativeMpint(t *testing.T, a []byte) []byte {
	t.Helper()
	p := splitSig(t, a)
	f := readStringsTest(t, p.signature, 2)
	rs := readStringsTest(t, f[1], 2)
	r, s := rs[0], rs[1]
	switch {
	case highBitPadded(r):
		r = r[1:]
	case highBitPadded(s):
		s = s[1:]
	default:
		t.Fatal("withNegativeMpint: neither r nor s has a leading 0x00 before a high-bit byte")
	}
	p.signature = packStrings(f[0], packStrings(r, s))
	return p.armored()
}

// withNonMinimalMpint prepends 0x00 to s.
func withNonMinimalMpint(t *testing.T, a []byte) []byte {
	t.Helper()
	p := splitSig(t, a)
	f := readStringsTest(t, p.signature, 2)
	rs := readStringsTest(t, f[1], 2)
	s := append([]byte{0}, rs[1]...)
	p.signature = packStrings(f[0], packStrings(rs[0], s))
	return p.armored()
}

// parseSafely runs ParseSSHSIG and fails the test, not the binary, on a panic.
func parseSafely(t *testing.T, in []byte) (s *SSHSig, err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ParseSSHSIG panicked on a %d-byte input: %v", len(in), r)
		}
	}()
	return ParseSSHSIG(in)
}

func TestSSHSIGFixturesAreTheGeneratedBytes(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("testdata", "sshsig"))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range entries {
		got = append(got, e.Name())
	}
	var pinned []string
	for k := range fixtureSHA256 {
		pinned = append(pinned, k)
	}
	sort.Strings(got)
	sort.Strings(pinned)
	if len(pinned) != 8 || strings.Join(got, " ") != strings.Join(pinned, " ") {
		t.Fatalf("testdata/sshsig holds %v, want exactly %v", got, pinned)
	}
	for _, name := range pinned {
		sum := sha256.Sum256(fixtureBytes(t, name))
		if h := hex.EncodeToString(sum[:]); h != fixtureSHA256[name] {
			t.Errorf("%s: sha256 %s, pinned %s", name, h, fixtureSHA256[name])
		}
	}

	mh := fixtureBytes(t, "message.hex")
	if string(mh) != fixtureMessageHex+"\n" {
		t.Errorf("message.hex is not fixtureMessageHex plus one newline")
	}
	wantMsg, err := SignedMessage(fixtureDigest(t), fixtureChallenge(), fixtureExpiry)
	if err != nil {
		t.Fatal(err)
	}
	gotMsg, err := hex.DecodeString(strings.TrimSuffix(string(mh), "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotMsg, wantMsg) {
		t.Errorf("message.hex does not decode to SignedMessage(fixture)")
	}

	for _, name := range []string{"ed25519.sig", "p256.sig", "git-ns-ed25519.sig"} {
		b := fixtureBytes(t, name)
		if len(b) == 0 || !bytes.HasPrefix(b, []byte("-----BEGIN SSH SIGNATURE-----\n")) {
			t.Errorf("%s: not an SSHSIG armor", name)
		}
	}
	for _, name := range pinned {
		b := fixtureBytes(t, name)
		if bytes.IndexByte(b, 0) >= 0 {
			t.Errorf("%s: contains a NUL byte", name)
		}
		if bytes.IndexByte(b, '\r') >= 0 {
			t.Errorf("%s: contains a carriage return", name)
		}
	}

	// D7 precondition: the negative-mpint arm needs a high-bit r or s in p256.sig.
	f := splitSig(t, fixtureBytes(t, "p256.sig"))
	rs := readStringsTest(t, readStringsTest(t, f.signature, 2)[1], 2)
	if !highBitPadded(rs[0]) && !highBitPadded(rs[1]) {
		t.Errorf("p256.sig has no high-bit r or s; re-run generate.sh (D7)")
	}
}

func TestParseSSHSIGFailsLoudly(t *testing.T) {
	assertFails := func(t *testing.T, label string, in []byte) {
		t.Helper()
		s, err := parseSafely(t, in)
		if s != nil || err == nil {
			t.Errorf("%s: got (%v, %v), want nil and an error", label, s, err)
		}
	}
	for _, name := range []string{"ed25519.sig", "p256.sig", "git-ns-ed25519.sig"} {
		blob := unarmorTest(t, fixtureBytes(t, name))
		for n := 0; n < len(blob); n++ {
			assertFails(t, fmt.Sprintf("%s prefix %d", name, n), armorTest(blob[:n]))
		}
		if s, err := parseSafely(t, fixtureBytes(t, name)); err != nil || s == nil {
			t.Errorf("%s: positive control failed: %v", name, err)
		}
	}

	rng := rand.New(rand.NewPCG(201, 373))
	for i := 0; i < 64; i++ {
		raw := make([]byte, 1+rng.IntN(512))
		for j := range raw {
			raw[j] = byte(rng.Uint32())
		}
		assertFails(t, fmt.Sprintf("random %d raw", i), raw)
		assertFails(t, fmt.Sprintf("random %d armored", i), armorTest(raw))
	}
	assertFails(t, "PGP block", []byte("-----BEGIN PGP SIGNATURE-----\n\niQEz\n-----END PGP SIGNATURE-----\n"))
	assertFails(t, "json", []byte(`{"sig":"x"}`))
	assertFails(t, "armorTest(hello)", armorTest([]byte("hello")))

	if s, err := parseSafely(t, nil); s != nil || !errors.Is(err, ErrArmor) {
		t.Errorf("empty input: got (%v, %v), want ErrArmor", s, err)
	}

	ed := fixtureBytes(t, "ed25519.sig")
	lines := strings.Split(strings.TrimSuffix(string(ed), "\n"), "\n")
	emptyBody := strings.Join(append([]string{lines[0], ""}, lines[1:]...), "\n") + "\n"
	faults := []struct {
		name string
		in   []byte
		want error
	}{
		{"CRLF", bytes.ReplaceAll(ed, []byte("\n"), []byte("\r\n")), ErrArmor},
		{"two trailing newlines", append(bytes.Clone(ed), '\n'), ErrArmor},
		{"leading blank line", append([]byte("\n"), ed...), ErrArmor},
		{"leading x", append([]byte("x"), ed...), ErrArmor},
		{"fixture twice", append(bytes.Clone(ed), ed...), ErrBase64},
		{"empty body line", []byte(emptyBody), ErrArmor},
	}
	for _, f := range faults {
		s, err := parseSafely(t, f.in)
		if s != nil || !errors.Is(err, f.want) {
			t.Errorf("%s: got (%v, %v), want %v", f.name, s, err, f.want)
		}
	}
}

func TestFingerprintMatchesSourceTxt(t *testing.T) {
	lines := strings.Split(string(fixtureBytes(t, "SOURCE.txt")), "\n")
	if !strings.HasPrefix(lines[0], "OpenSSH_") {
		t.Fatalf("SOURCE.txt line 1 is %q, want an OpenSSH_ version", lines[0])
	}
	for _, c := range []struct{ marker, pub string }{{"(ED25519)", "ed25519.pub"}, {"(ECDSA)", "p256.pub"}} {
		want := fixtureKey(t, c.pub).Fingerprint()
		if !strings.HasPrefix(want, "SHA256:") || strings.Contains(want, "=") {
			t.Errorf("%s: fingerprint %q is not unpadded SHA256:", c.pub, want)
		}
		found := false
		for _, l := range lines {
			if !strings.Contains(l, c.marker) {
				continue
			}
			found = true
			if f := strings.Fields(l); len(f) < 2 || f[1] != want {
				t.Errorf("%s: SOURCE.txt line %q, want fingerprint %s", c.pub, l, want)
			}
		}
		if !found {
			t.Errorf("SOURCE.txt has no %s line", c.marker)
		}
	}
}

func TestParseAuthorizedKeyRefusals(t *testing.T) {
	ed := fixtureBytes(t, "ed25519.pub")
	p2 := fixtureBytes(t, "p256.pub")
	if k, err := ParseAuthorizedKey(ed); err != nil || k.Type != "ssh-ed25519" {
		t.Fatalf("ed25519.pub: got (%v, %v)", k.Type, err)
	}
	if k, err := ParseAuthorizedKey(p2); err != nil || k.Type != "ecdsa-sha2-nistp256" {
		t.Fatalf("p256.pub: got (%v, %v)", k.Type, err)
	}
	edFields := strings.Fields(string(ed))
	if k, err := ParseAuthorizedKey([]byte(edFields[0] + " " + edFields[1] + "\n")); err != nil || k.Type != "ssh-ed25519" {
		t.Errorf("two-field line: got (%v, %v)", k.Type, err)
	}

	edBlob := fixtureKey(t, "ed25519.pub").Blob
	_, rest := firstStringTest(t, edBlob)
	skBlob := append(sshStringTest([]byte("sk-ssh-ed25519@openssh.com")), rest...)
	p2Fields := strings.Fields(string(p2))

	arms := []struct {
		name string
		line []byte
		want error
	}{
		{"empty", []byte(""), ErrAuthorizedKey},
		{"type only", []byte("ssh-ed25519"), ErrAuthorizedKey},
		{"line twice", append(bytes.Clone(ed), ed...), ErrAuthorizedKey},
		{"bad base64", []byte("ssh-ed25519 !!!notbase64"), ErrBase64},
		{"options prefix", append([]byte(`namespaces="x" `), ed...), ErrBase64},
		{"outer type vs p256 blob", []byte("ssh-ed25519 " + p2Fields[1] + "\n"), ErrKeyType},
		{"sk outer and inner", []byte("sk-ssh-ed25519@openssh.com " + base64.StdEncoding.EncodeToString(skBlob) + "\n"), ErrKeyType},
	}
	for _, a := range arms {
		if k, err := ParseAuthorizedKey(a.line); err == nil || !errors.Is(err, a.want) {
			t.Errorf("%s: got (%v, %v), want %v", a.name, k, err, a.want)
		}
	}
}

package approval

// Row 201 A2 sprint 1 M3 (design §6, PE1-PE4): the registration law over Chrome's recorded
// registrations. A registration carries no signature (fmt "none"), so each refusal arm edits the
// recorded bytes directly; the enrolment code is what binds them (PE4, sprint 3's WN10).

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// regFile is a registration file around rawId, clientDataJSON and attestationObject.
func regFile(rawID, cdj, attObj []byte) []byte {
	e := base64.RawURLEncoding.EncodeToString
	b, _ := json.Marshal(map[string]any{"id": e(rawID), "rawId": e(rawID), "type": "public-key",
		"response": map[string]string{"clientDataJSON": e(cdj), "attestationObject": e(attObj)}})
	return b
}

func attObj(fmtName string, attStmt cbMap, authData []byte) []byte {
	return cb(cbMap{"fmt", fmtName, "attStmt", attStmt, "authData", authData})
}

// chromeRegistration is the recorded registration as NewPasskeyApprover's input.
func chromeRegistration(t *testing.T, alg string) (chromeCeremony, RegistrationInput) {
	t.Helper()
	c := chromeFixture(t, alg)
	var ch [32]byte
	copy(ch[:], unhex(t, c.ChallengeHex))
	return c, RegistrationInput{Name: "mark-mac", Origin: c.Origin, Challenge: ch, Response: c.Registration, Now: 1000}
}

func unhex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// PE1: one arm per R1-R15 on an otherwise valid registration, each with the zero Approver.
func TestNewPasskeyApproverRefusals(t *testing.T) {
	c, base := chromeRegistration(t, "es256")
	ad := regAuthData(t, c)
	cdj := func(old, new string) []byte { return []byte(strings.Replace(string(c.RegCDJ), old, new, 1)) }
	file := func(cd []byte, ao []byte) []byte { return regFile(c.RegRawID, cd, ao) }
	arms := []struct {
		name string
		edit func(*RegistrationInput)
		want error
	}{
		{"R1-name", func(in *RegistrationInput) { in.Name = "Mark" }, ErrApproverName},
		{"R2-ip-origin", func(in *RegistrationInput) { in.Origin = strings.Replace(c.Origin, "localhost", "127.0.0.1", 1) }, ErrOriginHost},
		{"R3-envelope", func(in *RegistrationInput) {
			in.Response = bytes.Replace(c.Registration, []byte(`"public-key"`), []byte(`"x"`), 1)
		}, ErrPasskeyResponse},
		{"R4-client-data", func(in *RegistrationInput) { in.Response = file([]byte("{"), c.AttObj) }, ErrClientData},
		{"R5-type-get", func(in *RegistrationInput) { in.Response = file(cdj("webauthn.create", "webauthn.get"), c.AttObj) }, ErrClientDataType},
		{"R6-challenge", func(in *RegistrationInput) { in.Challenge[0] ^= 1 }, ErrChallenge},
		{"R7-origin", func(in *RegistrationInput) { in.Origin = "http://localhost:7644" }, ErrOriginMismatch},
		{"R8-cross-origin", func(in *RegistrationInput) {
			in.Response = file(cdj(`"crossOrigin":false`, `"crossOrigin":true`), c.AttObj)
		}, ErrCrossOrigin},
		{"R8-top-origin", func(in *RegistrationInput) {
			in.Response = file(cdj(`"crossOrigin":false`, `"crossOrigin":false,"topOrigin":"https://x.example"`), c.AttObj)
		}, ErrCrossOrigin},
		{"R9-trailing-byte", func(in *RegistrationInput) { in.Response = file(c.RegCDJ, append(bytes.Clone(c.AttObj), 0)) }, ErrCBORTrailing},
		{"R9-array", func(in *RegistrationInput) { in.Response = file(c.RegCDJ, cb(cbArr{"none", cbMap{}, ad})) }, ErrAttestationObject},
		{"R9-array-of-six", func(in *RegistrationInput) {
			in.Response = file(c.RegCDJ, cb(cbArr{"fmt", "none", "attStmt", cbMap{}, "authData", ad}))
		}, ErrAttestationObject},
		{"R9-extra-key", func(in *RegistrationInput) {
			in.Response = file(c.RegCDJ, cb(cbMap{"fmt", "none", "attStmt", cbMap{}, "authData", ad, "x", 1}))
		}, ErrAttestationObject},
		{"R9-authData-text", func(in *RegistrationInput) {
			in.Response = file(c.RegCDJ, cb(cbMap{"fmt", "none", "attStmt", cbMap{}, "authData", "x"}))
		}, ErrAttestationObject},
		{"R10-packed", func(in *RegistrationInput) { in.Response = file(c.RegCDJ, attObj("packed", cbMap{}, ad)) }, ErrAttestationFormat},
		{"R10-attStmt", func(in *RegistrationInput) { in.Response = file(c.RegCDJ, attObj("none", cbMap{"alg", -7}, ad)) }, ErrAttestationFormat},
		{"R11-no-AT", func(in *RegistrationInput) {
			in.Response = file(c.RegCDJ, attObj("none", cbMap{}, withFlags(ad, 0, 0x40)))
		}, ErrAuthDataFlags},
		{"R12-rp-id", func(in *RegistrationInput) {
			in.Origin = "https://x.ts.net"
			in.Response = file(cdj(c.Origin, in.Origin), c.AttObj)
		}, ErrRPIDHash},
		{"R13-no-UP", func(in *RegistrationInput) {
			in.Response = file(c.RegCDJ, attObj("none", cbMap{}, withFlags(ad, 0, 0x01)))
		}, ErrUserPresence},
		{"R14-no-UV", func(in *RegistrationInput) {
			in.Response = file(c.RegCDJ, attObj("none", cbMap{}, withFlags(ad, 0, 0x04)))
		}, ErrUserVerification},
		{"R15-raw-id", func(in *RegistrationInput) { in.Response = regFile(c.RegRawID[1:], c.RegCDJ, c.AttObj) }, ErrCredentialIDMismatch},
	}
	for _, a := range arms {
		t.Run(a.name, func(t *testing.T) {
			in := base
			a.edit(&in)
			if got, err := NewPasskeyApprover(in); !errors.Is(err, a.want) || !reflect.DeepEqual(got, Approver{}) {
				t.Fatalf("(%+v, %v), want the zero Approver and %v", got, err, a.want)
			}
		})
	}
	if _, err := NewPasskeyApprover(base); err != nil {
		t.Fatalf("control: %v", err)
	}
	if _, err := NewPasskeyApprover(RegistrationInput{Name: "mark-mac", Origin: base.Origin, Challenge: base.Challenge,
		Response: file(c.RegCDJ, attObj("none", cbMap{}, ad)), Now: 1}); err != nil {
		t.Fatalf("control (re-encoded attestationObject): %v", err)
	}
}

// PE2: both Chrome registrations enrol, every Passkey field as recorded.
func TestNewPasskeyApproverAcceptsChrome(t *testing.T) {
	for _, alg := range []string{"es256", "eddsa"} {
		c, in := chromeRegistration(t, alg)
		k, _ := chromeKey(t, alg)
		want := Approver{Name: "mark-mac", Format: FormatWebAuthn, EnrolledAt: 1000, Passkey: &Passkey{RPID: "localhost",
			Origin: c.Origin, CredentialID: c.RegRawID, Key: k, AAGUID: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}, SignCount: 1}}
		got, err := NewPasskeyApprover(in)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: (%+v %+v, %v), want %+v", alg, got, got.Passkey, err, want.Passkey)
		}
		// The same bytes re-made for a tailnet origin, synced (BE, BS), counter 7.
		ad := withFlags(regAuthData(t, c), 0x18, 0)
		h := sha256.Sum256([]byte("x.ts.net"))
		ad = binary.BigEndian.AppendUint32(append(h[:], ad[32]), 7)[:37:37]
		ad = append(ad, regAuthData(t, c)[37:]...)
		in.Origin, in.Response = "https://x.ts.net:8443", regFile(c.RegRawID, []byte(strings.Replace(string(c.RegCDJ), c.Origin, "https://x.ts.net:8443", 1)), attObj("none", cbMap{}, ad))
		want.Passkey.RPID, want.Passkey.Origin, want.Passkey.BackupEligible, want.Passkey.BackupState, want.Passkey.SignCount = "x.ts.net", in.Origin, true, true, 7
		if got, err := NewPasskeyApprover(in); err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("%s tailnet: (%+v, %v), want %+v", alg, got.Passkey, err, want.Passkey)
		}
	}
}

// passkeyWithFlags is the ES256 registration with authData's backup flags set, enrolled.
func passkeyWithFlags(t *testing.T, set byte) Approver {
	c, in := chromeRegistration(t, "es256")
	in.Response = regFile(c.RegRawID, c.RegCDJ, attObj("none", cbMap{}, withFlags(regAuthData(t, c), set, 0)))
	a, err := NewPasskeyApprover(in)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

// PE3: the enrolment y/N text, golden, and its backup labels.
func TestRenderPasskeyEnrolmentGolden(t *testing.T) {
	_, in := chromeRegistration(t, "es256")
	a, err := NewPasskeyApprover(in)
	if err != nil {
		t.Fatal(err)
	}
	const want = "approver  mark-mac\nformat    webauthn\npasskey   ES256 for localhost at http://localhost:56547\n" +
		"cred      914d608a1de4d7bf\naaguid    01020304050607080102030405060708\nbackup    device-bound\ncounter   1\n" +
		"records   approved_by=webauthn:mark-mac\n\n" +
		"This passkey will approve mint requests with no terminal, each time after your device verifies\n" +
		"you (Face ID, Touch ID, a PIN or a security-key touch). World cannot check which device holds\n" +
		"it; the 26-character enrolment code you typed, which matched this exact registration, is what ties\n" +
		"it to the browser you just used. Enrol only a passkey you created yourself just now.\n"
	if got := RenderPasskeyEnrolment(a); got != want {
		t.Fatalf("RenderPasskeyEnrolment:\n%s\nwant:\n%s", got, want)
	}
	for set, label := range map[byte]string{0x18: "backup    synced passkey\n", 0x08: "backup    backup-eligible\n"} {
		if got := RenderPasskeyEnrolment(passkeyWithFlags(t, set)); !strings.Contains(got, label) {
			t.Errorf("flags %#x: %q lacks %q", set, got, label)
		}
	}
	_, ed := chromeRegistration(t, "eddsa")
	if e, err := NewPasskeyApprover(ed); err != nil || !strings.Contains(RenderPasskeyEnrolment(e), "passkey   EdDSA for localhost") {
		t.Fatalf("EdDSA: %v", err)
	}
}

// PE4: the enrolment code is 26 characters of A-Z2-7, pinned on the recorded registration, and
// changes with one byte of either input.
func TestEnrolmentCodeGolden(t *testing.T) {
	c := chromeFixture(t, "es256")
	code := EnrolmentCode(c.RegCDJ, c.AttObj)
	if code != "YDOR4X36PQNFRTXBHWAV4W2XG4" || !regexp.MustCompile(`^[A-Z2-7]{26}$`).MatchString(code) {
		t.Fatalf("EnrolmentCode = %q", code)
	}
	if g := GroupEnrolmentCode(code); g != "YDOR4-X36PQ-NFRTX-BHWAV-4W2XG4" {
		t.Fatalf("GroupEnrolmentCode = %q", g)
	}
	for name, in := range map[string][2][]byte{
		"clientDataJSON":    {append(bytes.Clone(c.RegCDJ[:len(c.RegCDJ)-1]), '!'), c.AttObj},
		"attestationObject": {c.RegCDJ, append(bytes.Clone(c.AttObj[:len(c.AttObj)-1]), c.AttObj[len(c.AttObj)-1]^1)},
		"swapped":           {c.AttObj, c.RegCDJ},
	} {
		if got := EnrolmentCode(in[0], in[1]); got == code || len(got) != 26 {
			t.Errorf("a changed %s gives %q", name, got)
		}
	}
}

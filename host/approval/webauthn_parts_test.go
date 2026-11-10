package approval

// Row 201 A2 sprint 1 M2 (design §6, AD1-AD3, CD1-CD2, OR1, EN1): authenticator data,
// clientDataJSON, the ceremony origin and the response envelope.

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// withCredentialID rebuilds a registration's authData around an n-byte credential ID.
func withCredentialID(reg []byte, n int) []byte {
	l := int(binary.BigEndian.Uint16(reg[53:55]))
	b := binary.BigEndian.AppendUint16(bytes.Clone(reg[:53]), uint16(n))
	return append(append(b, bytes.Repeat([]byte{7}, n)...), reg[55+l:]...)
}

func withFlags(b []byte, set, clear byte) []byte {
	b = bytes.Clone(b)
	b[32] = b[32]&^clear | set
	return b
}

// AD1: D1-D8, with the credential-ID bounds at 0, 1, 1023 and 1024 bytes.
func TestAuthDataRefusals(t *testing.T) {
	c := chromeFixture(t, "es256")
	reg, asr := regAuthData(t, c), c.AsrAuthData
	key := coseKeyBytes(t, c)
	badKey := append(bytes.Clone(reg[:len(reg)-len(key)]), cb(cbMap{1, 3, 3, -257, -1, 1})...)
	arms := []struct {
		name         string
		in           []byte
		registration bool
		want         error
	}{
		{"D1-36-bytes", bytes.Clone(asr[:36]), false, ErrAuthDataShort},
		{"D2-BS-without-BE", withFlags(asr, 0x10, 0), false, ErrAuthDataFlags},
		{"D3-registration-without-AT", withFlags(reg, 0, 0x40), true, ErrAuthDataFlags},
		{"D3-assertion-with-AT", withFlags(asr, 0x40, 0), false, ErrAuthDataFlags},
		{"D4-inside-aaguid", bytes.Clone(reg[:37+17]), true, ErrAuthDataShort},
		{"D4-inside-credential-id", bytes.Clone(reg[:55+10]), true, ErrAuthDataShort},
		{"D5-0-byte-id", withCredentialID(reg, 0), true, ErrCredentialID},
		{"D5-1024-byte-id", withCredentialID(reg, 1024), true, ErrCredentialID},
		{"D6-RS256-key", badKey, true, ErrCOSEAlg},
		{"D6-truncated-key", bytes.Clone(reg[:len(reg)-1]), true, ErrCBORTruncated},
		{"D7-integer-extensions", append(withFlags(reg, 0x80, 0), 0x01), true, ErrAuthDataExtensions},
		{"D7-no-extensions", withFlags(asr, 0x80, 0), false, ErrAuthDataExtensions},
		{"D8-assertion-trailing", append(bytes.Clone(asr), 0), false, ErrAuthDataTrailing},
		{"D8-registration-trailing", append(bytes.Clone(reg), 0), true, ErrAuthDataTrailing},
	}
	for _, a := range arms {
		t.Run(a.name, func(t *testing.T) {
			if ad, err := ParseAuthData(a.in, a.registration); !errors.Is(err, a.want) || !reflect.DeepEqual(ad, AuthData{}) {
				t.Fatalf("(%+v, %v), want the zero value and %v", ad, err, a.want)
			}
		})
	}
	for _, n := range []int{1, 1023} {
		ad, err := ParseAuthData(withCredentialID(reg, n), true)
		if err != nil || len(ad.CredentialID) != n {
			t.Errorf("a %d-byte credential ID: %v (%d bytes)", n, err, len(ad.CredentialID))
		}
	}
}

// AD2: Chrome's registration and assertions read as recorded (V21).
func TestAuthDataReadsChrome(t *testing.T) {
	for _, alg := range []string{"es256", "eddsa"} {
		c := chromeFixture(t, alg)
		ad, err := ParseAuthData(regAuthData(t, c), true)
		if err != nil {
			t.Fatal(err)
		}
		k, _ := chromeKey(t, alg)
		aaguid := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
		if ad.RPIDHash != sha256.Sum256([]byte("localhost")) || !ad.UP || !ad.UV || ad.BE || ad.BS || ad.SignCount != 1 ||
			ad.AAGUID != aaguid || len(ad.CredentialID) != 32 || !bytes.Equal(ad.CredentialID, c.RegRawID) || !reflect.DeepEqual(ad.Key, k) {
			t.Fatalf("%s registration: %+v", alg, ad)
		}
		for want, b := range map[uint32][]byte{2: c.AsrAuthData, 3: c.Asr2AuthData} {
			as, err := ParseAuthData(b, false)
			if err != nil || len(b) != 37 || !as.UP || !as.UV || as.BE || as.BS || as.SignCount != want || as.CredentialID != nil {
				t.Fatalf("%s assertion %d: %+v, %v", alg, want, as, err)
			}
		}
	}
}

// AD3: an extensions map is parsed past and its contents ignored.
func TestAuthDataIgnoresExtensionsContents(t *testing.T) {
	c := chromeFixture(t, "eddsa")
	ext := cb(cbMap{"credProtect", 2})
	ad, err := ParseAuthData(append(withFlags(regAuthData(t, c), 0x80, 0), ext...), true)
	if k, _ := chromeKey(t, "eddsa"); err != nil || !reflect.DeepEqual(ad.Key, k) {
		t.Fatalf("registration with extensions: %v, %+v", err, ad.Key)
	}
	if _, err := ParseAuthData(append(withFlags(c.AsrAuthData, 0x80, 0), ext...), false); err != nil {
		t.Fatalf("assertion with extensions: %v", err)
	}
}

const cdjBase = `"type":"webauthn.get","challenge":"abc","origin":"http://localhost:7644"`

// CD1: J1-J7, and what is allowed: unknown members, any member order, topOrigin and crossOrigin.
func TestClientDataRefusals(t *testing.T) {
	arms := []struct {
		name, in string
		want     error
	}{
		{"J1-empty", "", ErrClientData},
		{"J1-2049-bytes", "{" + cdjBase + `,"x":"` + strings.Repeat("a", 2049-len(cdjBase)-9) + `"}`, ErrClientData},
		{"J2-not-utf8", `{"type":"webauthn.get","challenge":"abc","origin":"http://localhost:7644","x":"` + "\xff" + `"}`, ErrClientData},
		{"J3-array", `[]`, ErrClientData},
		{"J3-array-of-pairs", `["type","webauthn.get","challenge","abc","origin","http://localhost:7644"]`, ErrClientData},
		{"J3-unterminated", "{" + cdjBase, ErrClientData},
		{"J3-two-objects", "{" + cdjBase + "}{}", ErrClientData},
		{"J3-trailing", "{" + cdjBase + "}x", ErrClientData},
		{"J3-unknown-member-unterminated", "{" + cdjBase + `,"x":[1`, ErrClientData},
		{"J3-number-key", `{1:2}`, ErrClientData},
		{"J3-unknown-member-no-value", "{" + cdjBase + `,"x":}`, ErrClientData},
		{"J3-unknown-member-number-key", "{" + cdjBase + `,"x":{1:2}}`, ErrClientData},
		{"J3-known-member-no-value", `{"type":}`, ErrClientData},
		{"J4-challenge-twice", "{" + cdjBase + `,"challenge":"xyz"}`, ErrClientDataDuplicate},
		{"J4-crossOrigin-twice", "{" + cdjBase + `,"crossOrigin":false,"crossOrigin":true}`, ErrClientDataDuplicate},
		{"J5-no-origin", `{"type":"webauthn.get","challenge":"abc"}`, ErrClientData},
		{"J5-type-number", `{"type":1,"challenge":"abc","origin":"http://localhost:7644"}`, ErrClientData},
		{"J6-crossOrigin-string", "{" + cdjBase + `,"crossOrigin":"false"}`, ErrClientData},
		{"J6-topOrigin-bool", "{" + cdjBase + `,"topOrigin":true}`, ErrClientData},
		{"J7-depth-5", "{" + cdjBase + `,"x":{"a":{"b":{"c":{"d":1}}}}}`, ErrClientData},
	}
	for _, a := range arms {
		t.Run(a.name, func(t *testing.T) {
			if c, err := ParseClientData([]byte(a.in)); !errors.Is(err, a.want) || c != (ClientData{}) {
				t.Fatalf("(%+v, %v), want the zero value and %v", c, err, a.want)
			}
		})
	}
	for in, want := range map[string]ClientData{
		"{" + cdjBase + `,"other_keys_can_be_added_here":"x","x":{"a":{"b":[1]}},"x":2}`:                {"webauthn.get", "abc", "http://localhost:7644", false, false},
		`{"origin":"http://localhost:7644","crossOrigin":true,"challenge":"abc","type":"webauthn.get"}`: {"webauthn.get", "abc", "http://localhost:7644", true, false},
		"{" + cdjBase + `,"topOrigin":"https://x.example"} `:                                            {"webauthn.get", "abc", "http://localhost:7644", false, true},
	} {
		if c, err := ParseClientData([]byte(in)); err != nil || c != want {
			t.Errorf("%s: (%+v, %v), want %+v", in, c, err, want)
		}
	}
	if in := "{" + cdjBase + `,"x":"` + strings.Repeat("a", 2048-len(cdjBase)-9) + `"}`; len(in) != 2048 {
		t.Fatalf("the 2048-byte control is %d bytes", len(in))
	} else if _, err := ParseClientData([]byte(in)); err != nil {
		t.Errorf("2048 bytes: %v, want accepted", err)
	}
}

// CD2: Chrome's clientDataJSON, as recorded.
func TestClientDataReadsChrome(t *testing.T) {
	for _, alg := range []string{"es256", "eddsa"} {
		c := chromeFixture(t, alg)
		for typ, b := range map[string][]byte{"webauthn.create": c.RegCDJ, "webauthn.get": c.AsrCDJ} {
			want := ClientData{Type: typ, Challenge: "KRgdIk2YL2WdG9RwMjjT-EieJpaP_ioivXhvo78V_vU", Origin: c.Origin}
			if got, err := ParseClientData(b); err != nil || got != want {
				t.Errorf("%s %s: (%+v, %v)", alg, typ, got, err)
			}
		}
	}
}

// OR1: the accepted origins and their rpID; every refusal with its error.
func TestCeremonyOrigin(t *testing.T) {
	long := strings.Repeat("a", 63) + "." + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 54) + ".ts.net"
	for in, rp := range map[string]string{
		"http://localhost": "localhost", "http://localhost:7644": "localhost", "https://localhost:8443": "localhost",
		"https://localhost": "localhost", "https://voights-mac-studio.tail97eda0.ts.net": "voights-mac-studio.tail97eda0.ts.net",
		"https://x.ts.net:8443": "x.ts.net", "https://" + long: long, "https://x.ts.net:65535": "x.ts.net",
	} {
		if got, err := ParseCeremonyOrigin(in); err != nil || got != rp {
			t.Errorf("%s: (%q, %v), want %q", in, got, err, rp)
		}
	}
	for in, want := range map[string]error{
		"https://x.ts.net:443": ErrOrigin, "https://localhost:443": ErrOrigin, "http://localhost:80": ErrOrigin,
		"http://127.0.0.1:7644": ErrOriginHost, "https://evil.example": ErrOriginHost, "https://ts.net": ErrOriginHost,
		"https://[::1]": ErrOrigin, "https://LOCALHOST": ErrOrigin, "https://x.ts.net/": ErrOrigin,
		"https://x.ts.net:0": ErrOrigin, "https://x.ts.net:065535": ErrOrigin, "https://x.ts.net:65536": ErrOrigin,
		"https://x.ts.net:+8443": ErrOrigin, "https://x.ts.net:": ErrOrigin, "https://a..ts.net": ErrOrigin,
		"http://x.ts.net": ErrOrigin, "https://u@x.ts.net": ErrOrigin, "ftp://localhost": ErrOrigin, "localhost": ErrOrigin,
		"https://" + strings.Replace(long, ".ts.net", "d.ts.net", 1): ErrOrigin,
	} {
		got, err := ParseCeremonyOrigin(in)
		if !errors.Is(err, want) || got != "" || (want == ErrOrigin && errors.Is(err, ErrOriginHost)) {
			t.Errorf("%s: (%q, %v), want %v", in, got, err, want)
		}
		if strings.HasSuffix(in, ":443") || strings.HasSuffix(in, ":80") {
			if !strings.Contains(err.Error(), "default port") {
				t.Errorf("%s: %v does not name the default-port rule", in, err)
			}
		}
	}
}

// EN1: E2-E4, and Chrome's toJSON() output with its extra members.
func TestPasskeyEnvelope(t *testing.T) {
	const reg = `{"id":"AQI","rawId":"AQI","type":"public-key","response":{"clientDataJSON":"e30","attestationObject":"oA"}}`
	for name, in := range map[string]string{
		"E2-not-json":          `{`,
		"E2-no-response":       `{"id":"AQI","rawId":"AQI","type":"public-key"}`,
		"E2-no-id":             strings.Replace(reg, `"id":"AQI",`, "", 1),
		"E2-no-rawId":          strings.Replace(reg, `"rawId":"AQI",`, "", 1),
		"E2-no-type":           strings.Replace(reg, `"type":"public-key",`, "", 1),
		"E2-no-clientDataJSON": strings.Replace(reg, `"clientDataJSON":"e30",`, "", 1),
		"E2-null-attestation":  strings.Replace(reg, `"oA"`, "null", 1),
		"E3-type-x":            strings.Replace(reg, `"public-key"`, `"x"`, 1),
		"E3-id-not-rawId":      strings.Replace(reg, `"id":"AQI"`, `"id":"AQM"`, 1),
		"E4-padded":            strings.Replace(reg, `"oA"`, `"oA=="`, 1),
		"E4-non-canonical":     strings.Replace(reg, `"e30"`, `"e31"`, 1),
		"E4-rawId-standard":    strings.ReplaceAll(reg, `"AQI"`, `"A+I"`),
	} {
		if r, err := ParseRegistrationResponse([]byte(in)); !errors.Is(err, ErrPasskeyResponse) || !reflect.DeepEqual(r, RegistrationResponse{}) {
			t.Errorf("%s: (%+v, %v), want ErrPasskeyResponse", name, r, err)
		}
	}
	if r, err := ParseRegistrationResponse([]byte(reg)); err != nil || !bytes.Equal(r.RawID, []byte{1, 2}) || string(r.ClientDataJSON) != "{}" || !bytes.Equal(r.AttestationObject, []byte{0xa0}) {
		t.Fatalf("control: (%+v, %v)", r, err)
	}
	asr := `{"id":"AQI","rawId":"AQI","type":"public-key","response":{"clientDataJSON":"e30","authenticatorData":"AQ","signature":"Ag"%s}}`
	for name, in := range map[string]string{
		"E2-no-signature":         strings.Replace(asr, `,"signature":"Ag"`, "", 1),
		"E2-no-authenticatorData": strings.Replace(asr, `"authenticatorData":"AQ",`, "", 1),
		"E4-signature":            strings.Replace(asr, `"Ag"`, `"Ag=="`, 1),
		"E4-userHandle":           strings.Replace(asr, "%s", `,"userHandle":"A="`, 1),
		"E3-assertion-type":       strings.Replace(asr, `"public-key"`, `"x"`, 1),
	} {
		if r, err := ParseAssertionResponse([]byte(strings.Replace(in, "%s", "", 1))); !errors.Is(err, ErrPasskeyResponse) || r.RawID != nil {
			t.Errorf("%s: (%+v, %v), want ErrPasskeyResponse", name, r, err)
		}
	}
	for uh, want := range map[string]AssertionResponse{
		"":                   {RawID: []byte{1, 2}, ClientDataJSON: []byte("{}"), AuthenticatorData: []byte{1}, Signature: []byte{2}},
		`,"userHandle":null`: {RawID: []byte{1, 2}, ClientDataJSON: []byte("{}"), AuthenticatorData: []byte{1}, Signature: []byte{2}},
		`,"userHandle":"bWFyaw"`: {RawID: []byte{1, 2}, ClientDataJSON: []byte("{}"), AuthenticatorData: []byte{1}, Signature: []byte{2},
			UserHandle: []byte("mark"), HasUserHandle: true},
	} {
		if r, err := ParseAssertionResponse([]byte(strings.Replace(asr, "%s", uh, 1))); err != nil || !reflect.DeepEqual(r, want) {
			t.Errorf("userHandle %q: (%+v, %v), want %+v", uh, r, err, want)
		}
	}
	for _, alg := range []string{"es256", "eddsa"} {
		c := chromeFixture(t, alg)
		r, err := ParseRegistrationResponse(c.Registration)
		if err != nil || !bytes.Equal(r.RawID, c.RegRawID) || !bytes.Equal(r.ClientDataJSON, c.RegCDJ) || !bytes.Equal(r.AttestationObject, c.AttObj) {
			t.Fatalf("%s registration toJSON: (%+v, %v)", alg, r, err)
		}
		a, err := ParseAssertionResponse(c.Assertion)
		if err != nil || !bytes.Equal(a.RawID, c.RegRawID) || !bytes.Equal(a.Signature, c.AsrSig) || a.HasUserHandle {
			t.Fatalf("%s assertion toJSON: (%+v, %v)", alg, a, err)
		}
	}
}

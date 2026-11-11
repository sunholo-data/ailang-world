package approval

// Row 201 A2: the recorded Chrome fixtures (testdata/webauthn, design §2.15) and the test CBOR
// encoder the refusal arms are built with.

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// chromeCeremony is one recording, with its base64url fields decoded.
type chromeCeremony struct {
	Origin, RPID, ChallengeHex  string
	Registration, Assertion     []byte // the toJSON() envelopes, as recorded
	RegRawID, RegCDJ, AttObj    []byte
	AsrCDJ, AsrAuthData, AsrSig []byte
	Asr2AuthData                []byte
}

func b64u(t *testing.T, s string) []byte {
	t.Helper()
	b, err := base64.RawURLEncoding.Strict().DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// chromeFixture loads testdata/webauthn/chrome-<alg>.json ("es256" or "eddsa").
func chromeFixture(t *testing.T, alg string) chromeCeremony {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "webauthn", "chrome-"+alg+".json"))
	if err != nil {
		t.Fatal(err)
	}
	type resp struct{ ClientDataJSON, AttestationObject, AuthenticatorData, Signature string }
	var f struct {
		Origin                  string
		RPID                    string `json:"rp_id"`
		ChallengeHex            string `json:"challenge_hex"`
		Registration, Assertion json.RawMessage
		Assertion2              struct{ AuthenticatorData string }
	}
	var reg, asr struct {
		RawID    string
		Response resp
	}
	if json.Unmarshal(raw, &f) != nil || json.Unmarshal(f.Registration, &reg) != nil || json.Unmarshal(f.Assertion, &asr) != nil {
		t.Fatalf("chrome-%s.json does not decode", alg)
	}
	return chromeCeremony{Origin: f.Origin, RPID: f.RPID, ChallengeHex: f.ChallengeHex,
		Registration: f.Registration, Assertion: f.Assertion, RegRawID: b64u(t, reg.RawID),
		RegCDJ: b64u(t, reg.Response.ClientDataJSON), AttObj: b64u(t, reg.Response.AttestationObject),
		AsrCDJ: b64u(t, asr.Response.ClientDataJSON), AsrAuthData: b64u(t, asr.Response.AuthenticatorData),
		AsrSig: b64u(t, asr.Response.Signature), Asr2AuthData: b64u(t, f.Assertion2.AuthenticatorData)}
}

// regAuthData is the authData inside a recorded attestationObject; coseKeyBytes is the COSE key
// inside it (after rpIdHash, flags, counter, aaguid, the length and the credential ID).
func regAuthData(t *testing.T, c chromeCeremony) []byte {
	t.Helper()
	v, err := decodeCBORExact(c.AttObj, attestationLimits)
	if err != nil {
		t.Fatal(err)
	}
	ad, ok := v.get(cborValue{kind: cborText, b: []byte("authData")})
	if !ok {
		t.Fatal("no authData")
	}
	return ad.b
}

func coseKeyBytes(t *testing.T, c chromeCeremony) []byte {
	ad := regAuthData(t, c)
	return ad[55+int(binary.BigEndian.Uint16(ad[53:55])):]
}

// XW0: the fixtures are the recorded bytes, and SOURCE.txt names their origin.
func TestPasskeyFixturesAreTheGeneratedBytes(t *testing.T) {
	for name, want := range map[string]string{
		"chrome-es256.json": "924e9f79e873d487d83e93f0df22da342e1130db85c3e5dddc6bec9129e27639",
		"chrome-eddsa.json": "4c3fe46cb62898bd33c83244728ea84ace5ca25cb7df550e3759842928467b40",
	} {
		b, err := os.ReadFile(filepath.Join("testdata", "webauthn", name))
		if err != nil {
			t.Fatal(err)
		}
		if got := sha256.Sum256(b); hex.EncodeToString(got[:]) != want {
			t.Errorf("%s sha256 %x, want %s (fixtures are recorded, never hand-edited)", name, got, want)
		}
	}
	src, err := os.ReadFile(filepath.Join("testdata", "webauthn", "SOURCE.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Chrome/157.0.8103.50", "virtual authenticator", "D-WORLD-109", "chrome-es256.json  from", "chrome-eddsa.json  from"} {
		if !strings.Contains(string(src), want) {
			t.Errorf("SOURCE.txt does not name %q", want)
		}
	}
	for _, alg := range []string{"es256", "eddsa"} {
		c := chromeFixture(t, alg)
		if c.RPID != "localhost" || c.ChallengeHex != fixtureMessageHashHex || !strings.HasPrefix(c.Origin, "http://localhost:") {
			t.Errorf("%s: rp %q challenge %q origin %q", alg, c.RPID, c.ChallengeHex, c.Origin)
		}
	}
}

// ---- the test encoder: shortest heads only, so a non-canonical arm is written by hand.

type cbArr []any
type cbMap []any // key, value, key, value, ...

func cbHead(major byte, n uint64) []byte {
	switch {
	case n < 24:
		return []byte{major<<5 | byte(n)}
	case n < 1<<8:
		return []byte{major<<5 | 24, byte(n)}
	case n < 1<<16:
		return binary.BigEndian.AppendUint16([]byte{major<<5 | 25}, uint16(n))
	case n < 1<<32:
		return binary.BigEndian.AppendUint32([]byte{major<<5 | 26}, uint32(n))
	}
	return binary.BigEndian.AppendUint64([]byte{major<<5 | 27}, n)
}

func cb(v any) []byte {
	switch x := v.(type) {
	case int:
		if x < 0 {
			return cbHead(1, uint64(-1-x))
		}
		return cbHead(0, uint64(x))
	case []byte:
		return append(cbHead(2, uint64(len(x))), x...)
	case string:
		return append(cbHead(3, uint64(len(x))), x...)
	case bool:
		if x {
			return []byte{0xf5}
		}
		return []byte{0xf4}
	case cbArr:
		b := cbHead(4, uint64(len(x)))
		for _, e := range x {
			b = append(b, cb(e)...)
		}
		return b
	case cbMap:
		b := cbHead(5, uint64(len(x)/2))
		for _, e := range x {
			b = append(b, cb(e)...)
		}
		return b
	}
	panic("cb: unsupported type")
}

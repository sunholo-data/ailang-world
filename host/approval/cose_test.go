package approval

// Row 201 A2 sprint 1 M1 (design §6, CK1-CK4): ES256 and EdDSA COSE keys.

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"reflect"
	"testing"
)

// chromeKey is the parsed COSE key of a recorded registration.
func chromeKey(t *testing.T, alg string) (COSEKey, []byte) {
	t.Helper()
	raw := coseKeyBytes(t, chromeFixture(t, alg))
	v, err := decodeCBORExact(raw, attestationLimits)
	if err != nil {
		t.Fatal(err)
	}
	k, err := ParseCOSEKey(v)
	if err != nil {
		t.Fatalf("%s: %v", alg, err)
	}
	return k, raw
}

func parseCOSEBytes(b []byte) (COSEKey, error) {
	v, err := decodeCBORExact(b, attestationLimits)
	if err != nil {
		return COSEKey{}, err
	}
	return ParseCOSEKey(v)
}

// CK1: Q1-Q8, each refused with its error.
func TestCOSEKeyRefusals(t *testing.T) {
	ec, _ := chromeKey(t, "es256")
	ed, _ := chromeKey(t, "eddsa")
	x, y, ex := ec.X, ec.Y, ed.X
	offCurve := bytes.Clone(y)
	offCurve[31] ^= 1
	arms := []struct {
		name string
		key  any
		want error
	}{
		{"Q1-array", cbArr{1, 2}, ErrCOSEKey},
		{"Q2-no-kty", cbMap{3, -7, -1, 1, -2, x, -3, y}, ErrCOSEKey},
		{"Q2-no-alg", cbMap{1, 2, -1, 1, -2, x, -3, y}, ErrCOSEKey},
		{"Q2-alg-text", cbMap{1, 2, 3, "ES256", -1, 1, -2, x, -3, y}, ErrCOSEKey},
		{"Q2-no-crv", cbMap{1, 1, 3, -8, -2, ex}, ErrCOSEKey},
		{"Q3-RS256", cbMap{1, 3, 3, -257, -1, 1, -2, x, -3, y}, ErrCOSEAlg},
		{"Q3-ES384", cbMap{1, 2, 3, -35, -1, 2, -2, x, -3, y}, ErrCOSEAlg},
		{"Q4-ES256-kty-1", cbMap{1, 1, 3, -7, -1, 1, -2, x, -3, y}, ErrCOSEKey},
		{"Q4-ES256-crv-2", cbMap{1, 2, 3, -7, -1, 2, -2, x, -3, y}, ErrCOSEKey},
		{"Q5-EdDSA-crv-1", cbMap{1, 1, 3, -8, -1, 1, -2, ex}, ErrCOSEKey},
		{"Q5-EdDSA-kty-2", cbMap{1, 2, 3, -8, -1, 6, -2, ex}, ErrCOSEKey},
		{"Q6-EdDSA-31-byte-x", cbMap{1, 1, 3, -8, -1, 6, -2, ex[:31]}, ErrCOSEKey},
		{"Q6-EdDSA-x-text", cbMap{1, 1, 3, -8, -1, 6, -2, string(bytes.Repeat([]byte("a"), 32))}, ErrCOSEKey},
		{"Q6-ES256-no-y", cbMap{1, 2, 3, -7, -1, 1, -2, x}, ErrCOSEKey},
		{"Q7-off-curve", cbMap{1, 2, 3, -7, -1, 1, -2, x, -3, offCurve}, ErrCOSEKey},
		{"Q8-kid", cbMap{1, 2, 2, []byte{1}, 3, -7, -1, 1, -2, x, -3, y}, ErrCOSEKey},
		{"Q8-EdDSA-y", cbMap{1, 1, 3, -8, -1, 6, -2, ex, -3, ex}, ErrCOSEKey},
		{"Q8-text-label", cbMap{1, 1, 3, -8, -1, 6, -2, ex, "x", 1}, ErrCOSEKey},
	}
	for _, a := range arms {
		t.Run(a.name, func(t *testing.T) {
			k, err := parseCOSEBytes(cb(a.key))
			if !errors.Is(err, a.want) || (a.want == ErrCOSEKey && errors.Is(err, ErrCOSEAlg)) || !reflect.DeepEqual(k, COSEKey{}) {
				t.Fatalf("(%+v, %v), want the zero key and %v", k, err, a.want)
			}
		})
	}
	for name, m := range map[string]cbMap{"es256": {1, 2, 3, -7, -1, 1, -2, x, -3, y}, "eddsa": {1, 1, 3, -8, -1, 6, -2, ex}} {
		if _, err := parseCOSEBytes(cb(m)); err != nil {
			t.Errorf("control %s: %v", name, err)
		}
	}
}

// CK2: Chrome's keys verify Chrome's assertion signatures over authData || SHA-256(clientDataJSON).
func TestCOSEKeyVerifies(t *testing.T) {
	for _, alg := range []string{"es256", "eddsa"} {
		c := chromeFixture(t, alg)
		k, _ := chromeKey(t, alg)
		h := sha256.Sum256(c.AsrCDJ)
		signed := append(bytes.Clone(c.AsrAuthData), h[:]...)
		if err := k.verify(signed, c.AsrSig); err != nil {
			t.Fatalf("%s: Chrome's signature: %v", alg, err)
		}
		tampered := bytes.Clone(signed)
		tampered[len(tampered)-1] ^= 1
		for name, in := range map[string][2][]byte{
			"tampered-message": {tampered, c.AsrSig},
			"authData-only":    {c.AsrAuthData, c.AsrSig},
			"long-signature":   {signed, append(bytes.Clone(c.AsrSig), 0)},
			"short-signature":  {signed, c.AsrSig[:len(c.AsrSig)-1]},
		} {
			if err := k.verify(in[0], in[1]); !errors.Is(err, ErrBadSignature) {
				t.Errorf("%s %s: %v, want ErrBadSignature", alg, name, err)
			}
		}
	}
	if err := (COSEKey{Alg: -257}).verify(nil, nil); !errors.Is(err, ErrCOSEAlg) {
		t.Errorf("an unknown algorithm: %v, want ErrCOSEAlg", err)
	}
}

// CK3: Canonical() is Chrome's key bytes, and a reordered map re-encodes to them.
func TestCOSECanonicalIsChromesBytes(t *testing.T) {
	for _, alg := range []string{"es256", "eddsa"} {
		k, raw := chromeKey(t, alg)
		if !bytes.Equal(k.Canonical(), raw) {
			t.Fatalf("%s: Canonical %x, Chrome %x", alg, k.Canonical(), raw)
		}
		m := cbMap{-2, k.X, -1, 6, 3, -8, 1, 1}
		if alg == "es256" {
			m = cbMap{-3, k.Y, -2, k.X, -1, 1, 3, -7, 1, 2}
		}
		p, err := parseCOSEBytes(cb(m))
		if err != nil || !bytes.Equal(p.Canonical(), raw) {
			t.Fatalf("%s: a permuted map: %v, %x", alg, err, p.Canonical())
		}
	}
	if es, ed := (COSEKey{Alg: AlgES256}).AlgName(), (COSEKey{Alg: AlgEdDSA}).AlgName(); es != "ES256" || ed != "EdDSA" {
		t.Fatalf("AlgName %q %q", es, ed)
	}
}

// CK4: COSEKey holds exactly Alg, X and Y, and Canonical reads each.
func TestCOSEKeyHasOnlyDeclaredFields(t *testing.T) {
	typ := reflect.TypeOf(COSEKey{})
	var names []string
	for i := 0; i < typ.NumField(); i++ {
		names = append(names, typ.Field(i).Name)
	}
	if !reflect.DeepEqual(names, []string{"Alg", "X", "Y"}) {
		t.Fatalf("COSEKey fields %v, want [Alg X Y] (a new field must be encoded by Canonical)", names)
	}
	k, _ := chromeKey(t, "es256")
	for name, p := range map[string]func(*COSEKey){
		"Alg": func(c *COSEKey) { c.Alg = AlgEdDSA },
		"X":   func(c *COSEKey) { c.X = bytes.Clone(c.X); c.X[0] ^= 1 },
		"Y":   func(c *COSEKey) { c.Y = bytes.Clone(c.Y); c.Y[0] ^= 1 },
	} {
		c := k
		p(&c)
		if bytes.Equal(c.Canonical(), k.Canonical()) {
			t.Errorf("changing %s does not change Canonical()", name)
		}
	}
}

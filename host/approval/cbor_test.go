package approval

// Row 201 A2 sprint 1 M1 (design §6, CB1-CB4): the bounded CBOR subset.

import (
	"bytes"
	"errors"
	"testing"
)

func nested(depth int) []byte { return append(bytes.Repeat([]byte{0x81}, depth-1), 0x80) }

// intMap is a map of n entries i -> i.
func intMap(n int) cbMap {
	m := cbMap{}
	for i := range n {
		m = append(m, i, i)
	}
	return m
}

// CB1: one arm per refusal K1-K16, each errors.Is its own error.
func TestCBORRefusals(t *testing.T) {
	arms := []struct {
		name  string
		in    []byte
		exact bool
		want  error
	}{
		{"K1-head-argument", []byte{0x18}, false, ErrCBORTruncated},
		{"K1-empty", nil, false, ErrCBORTruncated},
		{"K1-in-a-map", []byte{0xa1, 0x01}, false, ErrCBORTruncated},
		{"K2-reserved", []byte{0x1c}, false, ErrCBORMalformed},
		{"K2-reserved-30", []byte{0x5e}, false, ErrCBORMalformed},
		{"K3-indefinite-bytes", []byte{0x5f, 0x41, 0x00, 0xff}, false, ErrCBORIndefinite},
		{"K3-indefinite-map", []byte{0xbf, 0xff}, false, ErrCBORIndefinite},
		{"K4-one-byte", []byte{0x18, 0x05}, false, ErrCBORNonCanonical},
		{"K4-one-byte-23", []byte{0x18, 0x17}, false, ErrCBORNonCanonical},
		{"K4-two-byte", []byte{0x19, 0x00, 0xff}, false, ErrCBORNonCanonical},
		{"K4-four-byte", []byte{0x1a, 0x00, 0x00, 0xff, 0xff}, false, ErrCBORNonCanonical},
		{"K4-eight-byte", []byte{0x1b, 0, 0, 0, 0, 0xff, 0xff, 0xff, 0xff}, false, ErrCBORNonCanonical},
		{"K4-length", []byte{0x58, 0x01, 0xaa}, false, ErrCBORNonCanonical},
		{"K5-bytes", []byte{0x58, 0x20, 1, 2, 3}, false, ErrCBORTruncated},
		{"K5-text", []byte{0x62, 'a'}, false, ErrCBORTruncated},
		{"K6-2049-bytes", cb(make([]byte, 2049)), false, ErrCBORTooLarge},
		{"K6-2049-text", cb(string(make([]byte, 2049))), false, ErrCBORTooLarge},
		{"K7-17-entry-map", cb(intMap(17)), false, ErrCBORTooLarge},
		{"K7-17-element-array", cb(cbArr{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), false, ErrCBORTooLarge},
		{"K8-five-nested", nested(5), false, ErrCBORTooDeep},
		{"K9-tag", []byte{0xc0, 0x00}, false, ErrCBORUnsupported},
		{"K10-half-float", []byte{0xf9, 0x00, 0x00}, false, ErrCBORUnsupported},
		{"K10-undefined", []byte{0xf7}, false, ErrCBORUnsupported},
		{"K10-simple-19", []byte{0xf3}, false, ErrCBORUnsupported},
		{"K10-break", []byte{0xff}, false, ErrCBORUnsupported},
		{"K11-utf8", []byte{0x61, 0xff}, false, ErrCBORMalformed},
		{"K12-bytes-key", cb(cbMap{[]byte{1}, 1}), false, ErrCBORUnsupported},
		{"K12-bool-key", cb(cbMap{true, 1}), false, ErrCBORUnsupported},
		{"K13-duplicate-int", []byte{0xa2, 0x01, 0x01, 0x01, 0x02}, false, ErrCBORDuplicateKey},
		{"K13-duplicate-text", cb(cbMap{"a", 1, "a", 2}), false, ErrCBORDuplicateKey},
		{"K14-uint-over-int64", []byte{0x1b, 0x80, 0, 0, 0, 0, 0, 0, 0}, false, ErrCBORTooLarge},
		{"K14-uint-max", []byte{0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, false, ErrCBORTooLarge},
		{"K14-negative", []byte{0x3b, 0x80, 0, 0, 0, 0, 0, 0, 0}, false, ErrCBORTooLarge},
		{"K15-65-items", cb(cbArr{intMap(16), intMap(15)}), false, ErrCBORTooLarge},
		{"K16-trailing", []byte{0x00, 0x00}, true, ErrCBORTrailing},
	}
	for _, a := range arms {
		t.Run(a.name, func(t *testing.T) {
			var err error
			if a.exact {
				_, err = decodeCBORExact(a.in, attestationLimits)
			} else {
				_, _, err = decodeCBOR(a.in, attestationLimits)
			}
			if !errors.Is(err, a.want) {
				t.Fatalf("decode(% x) = %v, want %v", a.in, err, a.want)
			}
		})
	}
	if _, rest, err := decodeCBOR([]byte{0x00, 0x00}, attestationLimits); err != nil || len(rest) != 1 {
		t.Fatalf("decodeCBOR leaves trailing bytes to its caller: rest %d, %v", len(rest), err)
	}
}

// CB2: every bound is inclusive, and the subset's values decode to what they encode.
func TestCBORBoundsAreInclusive(t *testing.T) {
	accepted := map[string][]byte{
		"2048-bytes":   cb(make([]byte, 2048)),
		"2048-text":    cb(string(make([]byte, 2048))),
		"16-entries":   cb(intMap(16)),
		"depth-4":      nested(4),
		"64-items":     cb(cbArr{intMap(16), intMap(14), 0}),
		"0-and-empty":  cb(cbMap{0, 1, "", 2}),
		"int64-limits": {0x83, 0x1b, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3b, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x18, 0x18},
	}
	for name, in := range accepted {
		if _, err := decodeCBORExact(in, attestationLimits); err != nil {
			t.Errorf("%s: %v, want accepted", name, err)
		}
	}
	v, err := decodeCBORExact([]byte{0x86, 0x20, 0x17, 0x43, 1, 2, 3, 0x62, 'h', 'i', 0xf5, 0xf6}, attestationLimits)
	if err != nil || len(v.items) != 6 {
		t.Fatalf("decode: %v, %+v", err, v)
	}
	it := v.items
	if it[0].kind != cborInt || it[0].i != -1 || it[1].i != 23 || it[2].kind != cborBytes || !bytes.Equal(it[2].b, []byte{1, 2, 3}) ||
		it[3].kind != cborText || string(it[3].b) != "hi" || it[4].kind != cborBool || !it[4].t || it[5].kind != cborNull {
		t.Fatalf("decoded %+v", it)
	}
	if f, err := decodeCBORExact([]byte{0xf4}, attestationLimits); err != nil || f.kind != cborBool || f.t {
		t.Fatalf("false decoded as %+v, %v", f, err)
	}
}

// CB3: Chrome's attestationObjects and COSE keys decode with nothing left over.
func TestCBORAcceptsChromesBytes(t *testing.T) {
	for alg, keyLen := range map[string]int{"es256": 77, "eddsa": 42} {
		c := chromeFixture(t, alg)
		v, err := decodeCBORExact(c.AttObj, attestationLimits)
		if err != nil || v.kind != cborMap || len(v.items) != 6 {
			t.Fatalf("%s attestationObject: %v (%d items)", alg, err, len(v.items))
		}
		if k := coseKeyBytes(t, c); len(k) != keyLen {
			t.Fatalf("%s COSE key is %d bytes, want %d", alg, len(k), keyLen)
		} else if _, err := decodeCBORExact(k, attestationLimits); err != nil {
			t.Fatalf("%s COSE key: %v", alg, err)
		}
	}
}

// CB4: never panics; an accepted item is never longer than its input.
func FuzzDecodeCBOR(f *testing.F) {
	f.Add([]byte{0x18})
	f.Add(cb(intMap(16)))
	f.Add(nested(5))
	f.Add([]byte{0xa2, 0x01, 0x01, 0x01, 0x02})
	f.Fuzz(func(t *testing.T, b []byte) {
		_, rest, err := decodeCBOR(b, attestationLimits)
		if err == nil && len(rest) > len(b) {
			t.Fatalf("rest %d > input %d", len(rest), len(b))
		}
	})
}

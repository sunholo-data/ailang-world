package approval

// Row 201 A2 sprint 3 M6b (plan §1.4): ParseStoredCOSEKey, which approve --webauthn-assertion
// rebuilds an enrolled passkey's key with, reads back exactly what Canonical() stored.

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestParseStoredCOSEKey(t *testing.T) {
	for _, alg := range []string{"es256", "eddsa"} {
		k, raw := chromeKey(t, alg)
		got, err := ParseStoredCOSEKey(k.Canonical())
		if err != nil || !reflect.DeepEqual(got, k) {
			t.Fatalf("%s: %+v, %v; want Chrome's key back", alg, got, err)
		}
		for name, arm := range map[string]struct {
			b    []byte
			want error
		}{
			"trailing":  {append(bytes.Clone(raw), 0), ErrCBORTrailing},
			"truncated": {raw[:len(raw)-1], ErrCBORTruncated},
			"empty":     {nil, ErrCBORTruncated},
		} {
			if _, err := ParseStoredCOSEKey(arm.b); !errors.Is(err, ErrCOSEKey) || !errors.Is(err, arm.want) {
				t.Fatalf("%s %s: %v, want ErrCOSEKey wrapping %v", alg, name, err, arm.want)
			}
		}
	}
	if _, err := ParseStoredCOSEKey(cb(cbMap{1, 3, 3, -257, -1, 1})); !errors.Is(err, ErrCOSEAlg) {
		t.Fatalf("an RS256 key: %v, want ErrCOSEAlg", err)
	}
}

package daemon

import (
	"encoding/base64"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"strings"
	"testing"
)

func TestWorkbenchReferenceCursor(t *testing.T) {
	raw := func(s string) string { return base64.RawURLEncoding.EncodeToString([]byte(s)) }
	for _, c := range []store.ObjectReferenceCursor{{Kind: store.ReferenceTransitionRef, EntryIndex: 0}, {Kind: store.ReferenceInterpreter, EntryIndex: 10}, {Kind: store.ReferenceStateRoot, WorldRef: hashref.SumSHA256([]byte("world"))}} {
		t.Run("round-trip", func(t *testing.T) {
			got, err := decodeReferenceCursor(encodeReferenceCursor(c))
			if err != nil || got != c {
				t.Fatalf("%+v %v", got, err)
			}
		})
	}
	bad := map[string]string{"version": `{"v":2,"kind":0,"key":"0"}`, "duplicate": `{"v":1,"v":1,"kind":0,"key":"0"}`, "unknown": `{"v":1,"kind":0,"key":"0","extra":1}`, "leading-zero": `{"v":1,"kind":0,"key":"00"}`, "negative": `{"v":1,"kind":0,"key":"-1"}`, "hash": `{"v":1,"kind":3,"key":"bad"}`, "trailing": `{"v":1,"kind":0,"key":"0"} {}`}
	for name, input := range bad {
		t.Run(name, func(t *testing.T) {
			if _, err := decodeReferenceCursor(raw(input)); err == nil {
				t.Fatal("accepted", input)
			}
		})
	}
	t.Run("cap", func(t *testing.T) {
		if _, err := decodeReferenceCursor(strings.Repeat("A", 513)); err == nil {
			t.Fatal("accepted")
		}
	})
}

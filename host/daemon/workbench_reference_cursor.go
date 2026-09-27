package daemon

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func encodeReferenceCursor(c store.ObjectReferenceCursor) string {
	key := c.WorldRef.String()
	if c.Kind != store.ReferenceStateRoot {
		key = strconv.FormatInt(c.EntryIndex, 10)
	}
	b, _ := json.Marshal(struct {
		V    int                 `json:"v"`
		Kind store.ReferenceKind `json:"kind"`
		Key  string              `json:"key"`
	}{1, c.Kind, key})
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodeReferenceCursor(text string) (store.ObjectReferenceCursor, error) {
	bad := func() (store.ObjectReferenceCursor, error) {
		return store.ObjectReferenceCursor{}, fmt.Errorf("invalid refsAfter cursor")
	}
	if len(text) == 0 || len(text) > 512 {
		return bad()
	}
	b, err := base64.RawURLEncoding.DecodeString(text)
	if err != nil || base64.RawURLEncoding.EncodeToString(b) != text {
		return bad()
	}
	d := json.NewDecoder(bytes.NewReader(b))
	tok, err := d.Token()
	if err != nil || tok != json.Delim('{') {
		return bad()
	}
	values := map[string]json.RawMessage{}
	for d.More() {
		tok, err = d.Token()
		if err != nil {
			return bad()
		}
		key, ok := tok.(string)
		if !ok || (key != "v" && key != "kind" && key != "key") {
			return bad()
		}
		if _, exists := values[key]; exists {
			return bad()
		}
		var raw json.RawMessage
		if err = d.Decode(&raw); err != nil {
			return bad()
		}
		values[key] = raw
	}
	tok, err = d.Token()
	if err != nil || tok != json.Delim('}') || len(values) != 3 {
		return bad()
	}
	if _, err = d.Token(); err != io.EOF {
		return bad()
	}
	var v int
	if err = json.Unmarshal(values["v"], &v); err != nil || v != 1 {
		return bad()
	}
	var kind store.ReferenceKind
	if err = json.Unmarshal(values["kind"], &kind); err != nil || kind < store.ReferenceTransitionRef || kind > store.ReferenceStateRoot {
		return bad()
	}
	var key string
	if err = json.Unmarshal(values["key"], &key); err != nil {
		return bad()
	}
	c := store.ObjectReferenceCursor{Kind: kind}
	if kind == store.ReferenceStateRoot {
		c.WorldRef, err = hashref.Parse(key)
		if err != nil || c.WorldRef.String() != key {
			return bad()
		}
	} else {
		c.EntryIndex, err = strconv.ParseInt(key, 10, 64)
		if err != nil || c.EntryIndex < 0 || strconv.FormatInt(c.EntryIndex, 10) != key {
			return bad()
		}
	}
	return c, nil
}

package evidence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf8"
)

type member struct {
	name string
	raw  json.RawMessage
}

func refusal(kind string, format string, args ...any) error {
	return &DecodeRefusal{Kind: kind, Detail: fmt.Sprintf(format, args...)}
}

// objectMembers uses Decoder tokens so duplicate members remain observable.
func objectMembers(raw []byte) ([]member, error) {
	if !utf8.Valid(raw) {
		return nil, refusal(RefusalInvalidUTF8, "input is not valid UTF-8")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	tok, err := d.Token()
	if err != nil {
		return nil, refusal(RefusalMalformed, "%v", err)
	}
	if tok != json.Delim('{') {
		return nil, refusal(RefusalMalformed, "top-level value is not an object")
	}
	seen := make(map[string]bool)
	var out []member
	for d.More() {
		kt, err := d.Token()
		if err != nil {
			return nil, refusal(RefusalMalformed, "%v", err)
		}
		name, ok := kt.(string)
		if !ok {
			return nil, refusal(RefusalMalformed, "object member name is not a string")
		}
		if seen[name] {
			return nil, refusal(RefusalMalformed, "duplicate member %q", name)
		}
		seen[name] = true
		var value json.RawMessage
		if err := d.Decode(&value); err != nil {
			return nil, refusal(RefusalMalformed, "%v", err)
		}
		out = append(out, member{name, value})
	}
	if _, err := d.Token(); err != nil {
		return nil, refusal(RefusalMalformed, "%v", err)
	}
	if _, err := d.Token(); err != io.EOF {
		if err == nil {
			return nil, refusal(RefusalMalformed, "trailing JSON value")
		}
		return nil, refusal(RefusalMalformed, "trailing data: %v", err)
	}
	return out, nil
}

func decodeExact(raw json.RawMessage, dst any) error {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if err := d.Decode(dst); err != nil {
		return refusal(RefusalMalformed, "%v", err)
	}
	if _, err := d.Token(); err != io.EOF {
		return refusal(RefusalMalformed, "trailing member value")
	}
	return nil
}

func checkString(label, s string) error {
	if len(s) > MaxStringBytes {
		return refusal(RefusalLimit, "%s is %d bytes; limit is %d", label, len(s), MaxStringBytes)
	}
	return nil
}

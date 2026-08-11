// Package transitionreg implements the immutable transition descriptor registry.
package transitionreg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

const (
	SemanticIDV1         = "world/transition-registry/v1"
	maxSchemaRaw         = 262144
	maxSchemaCanonical   = 65536
	maxRevisionRaw       = 16777216
	maxRevisionCanonical = 8388608
	maxEntries           = 1024
	maxNumberDigits      = 1024
	maxNumberExponent    = 10000
	maxNumberToken       = 16384
)

var InterfaceHashV1 = hashref.MustParse("sha256:743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac")

type jsonKind uint8

const (
	kindNull jsonKind = iota
	kindBool
	kindNumber
	kindString
	kindArray
	kindObject
)

type jsonMember struct {
	key   string
	value jsonValue
}
type jsonValue struct {
	kind jsonKind
	b    bool
	s    string
	a    []jsonValue
	o    []jsonMember
}

func parseJSON(raw []byte, maxRaw int) (jsonValue, error) {
	if len(raw) > maxRaw {
		return jsonValue{}, fmt.Errorf("raw JSON is %d bytes; limit is %d", len(raw), maxRaw)
	}
	if !utf8.Valid(raw) {
		return jsonValue{}, errors.New("JSON is not valid UTF-8")
	}
	if err := rejectEscapedSurrogates(raw); err != nil {
		return jsonValue{}, err
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	v, err := parseValue(d)
	if err != nil {
		return jsonValue{}, err
	}
	if _, err := d.Token(); err != io.EOF {
		if err == nil {
			return jsonValue{}, errors.New("multiple JSON values")
		}
		return jsonValue{}, fmt.Errorf("trailing JSON: %w", err)
	}
	return v, nil
}

func rejectEscapedSurrogates(raw []byte) error {
	inString, escaped := false, false
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if !inString {
			if c == '"' {
				inString = true
			}
			continue
		}
		if escaped {
			escaped = false
			if c == 'u' && i+4 < len(raw) {
				n, err := strconv.ParseUint(string(raw[i+1:i+5]), 16, 16)
				if err == nil && n >= 0xd800 && n <= 0xdfff {
					return errors.New("JSON string contains an escaped surrogate")
				}
			}
			continue
		}
		if c == '\\' {
			escaped = true
		} else if c == '"' {
			inString = false
		}
	}
	return nil
}

func parseValue(d *json.Decoder) (jsonValue, error) {
	t, err := d.Token()
	if err != nil {
		return jsonValue{}, fmt.Errorf("decode JSON: %w", err)
	}
	switch x := t.(type) {
	case nil:
		return jsonValue{kind: kindNull}, nil
	case bool:
		return jsonValue{kind: kindBool, b: x}, nil
	case string:
		return jsonValue{kind: kindString, s: x}, nil
	case json.Number:
		n, err := normalizeNumber(string(x))
		if err != nil {
			return jsonValue{}, err
		}
		return jsonValue{kind: kindNumber, s: n}, nil
	case json.Delim:
		switch x {
		case '[':
			v := jsonValue{kind: kindArray}
			for d.More() {
				child, err := parseValue(d)
				if err != nil {
					return jsonValue{}, err
				}
				v.a = append(v.a, child)
			}
			if end, err := d.Token(); err != nil || end != json.Delim(']') {
				return jsonValue{}, errors.New("unterminated JSON array")
			}
			return v, nil
		case '{':
			v := jsonValue{kind: kindObject}
			seen := map[string]struct{}{}
			for d.More() {
				kt, err := d.Token()
				if err != nil {
					return jsonValue{}, err
				}
				key, ok := kt.(string)
				if !ok {
					return jsonValue{}, errors.New("object key is not a string")
				}
				if _, ok := seen[key]; ok {
					return jsonValue{}, fmt.Errorf("duplicate object member %q", key)
				}
				seen[key] = struct{}{}
				child, err := parseValue(d)
				if err != nil {
					return jsonValue{}, err
				}
				v.o = append(v.o, jsonMember{key, child})
			}
			if end, err := d.Token(); err != nil || end != json.Delim('}') {
				return jsonValue{}, errors.New("unterminated JSON object")
			}
			return v, nil
		}
	}
	return jsonValue{}, fmt.Errorf("unsupported JSON token %v", t)
}

func normalizeNumber(raw string) (string, error) {
	s := raw
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	e := 0
	if i := strings.IndexAny(s, "eE"); i >= 0 {
		var err error
		e, err = strconv.Atoi(s[i+1:])
		if err != nil {
			return "", fmt.Errorf("invalid number %q", raw)
		}
		s = s[:i]
		if e < -maxNumberExponent || e > maxNumberExponent {
			return "", fmt.Errorf("number exponent exceeds %d", maxNumberExponent)
		}
	}
	frac := 0
	if i := strings.IndexByte(s, '.'); i >= 0 {
		frac = len(s) - i - 1
		s = s[:i] + s[i+1:]
	}
	if len(s) > maxNumberDigits {
		return "", fmt.Errorf("number coefficient has %d digits; limit is %d", len(s), maxNumberDigits)
	}
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return "0", nil
	}
	point := len(s) + e - frac
	var out string
	switch {
	case point <= 0:
		out = "0." + strings.Repeat("0", -point) + s
	case point >= len(s):
		out = s + strings.Repeat("0", point-len(s))
	default:
		out = s[:point] + "." + s[point:]
	}
	if strings.Contains(out, ".") {
		out = strings.TrimRight(out, "0")
		out = strings.TrimSuffix(out, ".")
	}
	if neg {
		out = "-" + out
	}
	if len(out) > maxNumberToken {
		return "", fmt.Errorf("normalized number is %d bytes; limit is %d", len(out), maxNumberToken)
	}
	return out, nil
}

func appendJSON(dst []byte, v jsonValue) []byte {
	switch v.kind {
	case kindNull:
		return append(dst, "null"...)
	case kindBool:
		if v.b {
			return append(dst, "true"...)
		}
		return append(dst, "false"...)
	case kindNumber:
		return append(dst, v.s...)
	case kindString:
		return appendJSONString(dst, v.s)
	case kindArray:
		dst = append(dst, '[')
		for i := range v.a {
			if i > 0 {
				dst = append(dst, ',')
			}
			dst = appendJSON(dst, v.a[i])
		}
		return append(dst, ']')
	case kindObject:
		members := append([]jsonMember(nil), v.o...)
		sort.Slice(members, func(i, j int) bool { return members[i].key < members[j].key })
		dst = append(dst, '{')
		for i, m := range members {
			if i > 0 {
				dst = append(dst, ',')
			}
			dst = appendJSONString(dst, m.key)
			dst = append(dst, ':')
			dst = appendJSON(dst, m.value)
		}
		return append(dst, '}')
	default:
		panic("unknown JSON kind")
	}
}

func appendJSONString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	const hex = "0123456789abcdef"
	for _, r := range s {
		switch r {
		case '"', '\\':
			dst = append(dst, '\\', byte(r))
		case '\b':
			dst = append(dst, "\\b"...)
		case '\t':
			dst = append(dst, "\\t"...)
		case '\n':
			dst = append(dst, "\\n"...)
		case '\f':
			dst = append(dst, "\\f"...)
		case '\r':
			dst = append(dst, "\\r"...)
		default:
			if r < 0x20 {
				dst = append(dst, "\\u00"...)
				dst = append(dst, hex[byte(r)>>4], hex[byte(r)&15])
			} else {
				dst = utf8.AppendRune(dst, r)
			}
		}
	}
	return append(dst, '"')
}

func canonicalSchema(raw []byte) ([]byte, error) {
	v, err := parseJSON(raw, maxSchemaRaw)
	if err != nil {
		return nil, err
	}
	if v.kind != kindObject {
		return nil, errors.New("schema root must be an object")
	}
	out := appendJSON(nil, v)
	if len(out) > maxSchemaCanonical {
		return nil, fmt.Errorf("canonical schema is %d bytes; limit is %d", len(out), maxSchemaCanonical)
	}
	return out, nil
}

func EncodeRevision(r Revision) ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	v, err := revisionValue(r)
	if err != nil {
		return nil, err
	}
	out := appendJSON(nil, v)
	if len(out) > maxRevisionCanonical {
		return nil, fmt.Errorf("canonical revision is %d bytes; limit is %d", len(out), maxRevisionCanonical)
	}
	return out, nil
}

func revisionValue(r Revision) (jsonValue, error) {
	entries := make([]jsonValue, len(r.Entries))
	for i, d := range r.Entries {
		v, err := descriptorValue(d)
		if err != nil {
			return jsonValue{}, err
		}
		entries[i] = v
	}
	return object(
		member("semanticID", str(r.SemanticID)), member("interfaceHash", str(r.InterfaceHash.String())), member("revision", num(r.Revision)), member("parent", str(r.Parent.String())), member("entries", jsonValue{kind: kindArray, a: entries}),
	), nil
}

func descriptorValue(d Descriptor) (jsonValue, error) {
	in, err := parseJSON(d.InputSchema, maxSchemaCanonical)
	if err != nil {
		return jsonValue{}, err
	}
	out, err := parseJSON(d.OutputSchema, maxSchemaCanonical)
	if err != nil {
		return jsonValue{}, err
	}
	effects := make([]jsonValue, len(d.DeclaredEffects))
	for i, e := range d.DeclaredEffects {
		effects[i] = requirementValue(e)
	}
	return object(member("id", str(d.ID)), member("transitionFn", str(d.TransitionFn.String())), member("interpreter", str(d.Interpreter.String())), member("semanticsEpoch", num(d.SemanticsEpoch)), member("inputSchema", in), member("outputSchema", out), member("access", requirementValue(d.Access)), member("declaredEffects", jsonValue{kind: kindArray, a: effects}), member("title", str(d.Title)), member("description", str(d.Description))), nil
}
func requirementValue(e EffectRequirement) jsonValue {
	return object(member("effect", str(e.Effect)), member("scope", str(e.Scope)), member("cost", num(e.Cost)))
}
func object(ms ...jsonMember) jsonValue       { return jsonValue{kind: kindObject, o: ms} }
func member(k string, v jsonValue) jsonMember { return jsonMember{k, v} }
func str(s string) jsonValue                  { return jsonValue{kind: kindString, s: s} }
func num(n int64) jsonValue                   { return jsonValue{kind: kindNumber, s: strconv.FormatInt(n, 10)} }

func DecodeRevision(raw []byte) (Revision, error) {
	v, err := parseJSON(raw, maxRevisionRaw)
	if err != nil {
		return Revision{}, err
	}
	if v.kind != kindObject {
		return Revision{}, errors.New("revision root must be an object")
	}
	canonical := appendJSON(nil, v)
	if len(canonical) > maxRevisionCanonical {
		return Revision{}, fmt.Errorf("canonical revision is %d bytes; limit is %d", len(canonical), maxRevisionCanonical)
	}
	r, err := decodeRevisionValue(v)
	if err != nil {
		return Revision{}, err
	}
	if err := r.Validate(); err != nil {
		return Revision{}, err
	}
	reencoded, err := EncodeRevision(r)
	if err != nil {
		return Revision{}, err
	}
	if !bytes.Equal(canonical, reencoded) {
		return Revision{}, errors.New("revision is not canonical for its typed schema")
	}
	return r, nil
}

func fields(v jsonValue, names ...string) (map[string]jsonValue, error) {
	if v.kind != kindObject {
		return nil, errors.New("expected object")
	}
	m := map[string]jsonValue{}
	allowed := map[string]bool{}
	for _, n := range names {
		allowed[n] = true
	}
	for _, p := range v.o {
		if !allowed[p.key] {
			return nil, fmt.Errorf("unknown key %q", p.key)
		}
		m[p.key] = p.value
	}
	for _, n := range names {
		if _, ok := m[n]; !ok {
			return nil, fmt.Errorf("missing key %q", n)
		}
	}
	return m, nil
}
func decodeRevisionValue(v jsonValue) (Revision, error) {
	m, err := fields(v, "semanticID", "interfaceHash", "revision", "parent", "entries")
	if err != nil {
		return Revision{}, err
	}
	sid, err := asString(m["semanticID"])
	if err != nil {
		return Revision{}, err
	}
	iface, err := asRef(m["interfaceHash"], false)
	if err != nil {
		return Revision{}, err
	}
	rev, err := asInt(m["revision"])
	if err != nil {
		return Revision{}, err
	}
	parent, err := asRef(m["parent"], true)
	if err != nil {
		return Revision{}, err
	}
	if m["entries"].kind != kindArray {
		return Revision{}, errors.New("entries must be an array")
	}
	if len(m["entries"].a) > maxEntries {
		return Revision{}, fmt.Errorf("entries exceeds %d", maxEntries)
	}
	r := Revision{SemanticID: sid, InterfaceHash: iface, Revision: rev, Parent: parent}
	for _, ev := range m["entries"].a {
		d, err := decodeDescriptor(ev)
		if err != nil {
			return Revision{}, err
		}
		r.Entries = append(r.Entries, d)
	}
	return r, nil
}
func decodeDescriptor(v jsonValue) (Descriptor, error) {
	m, err := fields(v, "id", "transitionFn", "interpreter", "semanticsEpoch", "inputSchema", "outputSchema", "access", "declaredEffects", "title", "description")
	if err != nil {
		return Descriptor{}, err
	}
	d := Descriptor{}
	if d.ID, err = asString(m["id"]); err != nil {
		return d, err
	}
	if d.TransitionFn, err = asRef(m["transitionFn"], false); err != nil {
		return d, err
	}
	if d.Interpreter, err = asRef(m["interpreter"], false); err != nil {
		return d, err
	}
	if d.SemanticsEpoch, err = asInt(m["semanticsEpoch"]); err != nil {
		return d, err
	}
	d.InputSchema = appendJSON(nil, m["inputSchema"])
	d.OutputSchema = appendJSON(nil, m["outputSchema"])
	if d.Access, err = decodeRequirement(m["access"]); err != nil {
		return d, err
	}
	if m["declaredEffects"].kind != kindArray {
		return d, errors.New("declaredEffects must be an array")
	}
	for _, x := range m["declaredEffects"].a {
		e, eerr := decodeRequirement(x)
		if eerr != nil {
			return d, eerr
		}
		d.DeclaredEffects = append(d.DeclaredEffects, e)
	}
	if d.Title, err = asString(m["title"]); err != nil {
		return d, err
	}
	if d.Description, err = asString(m["description"]); err != nil {
		return d, err
	}
	return d, nil
}
func decodeRequirement(v jsonValue) (EffectRequirement, error) {
	m, err := fields(v, "effect", "scope", "cost")
	if err != nil {
		return EffectRequirement{}, err
	}
	effect, err := asString(m["effect"])
	if err != nil {
		return EffectRequirement{}, err
	}
	scope, err := asString(m["scope"])
	if err != nil {
		return EffectRequirement{}, err
	}
	cost, err := asInt(m["cost"])
	return EffectRequirement{effect, scope, cost}, err
}
func asString(v jsonValue) (string, error) {
	if v.kind != kindString {
		return "", errors.New("expected string")
	}
	return v.s, nil
}
func asInt(v jsonValue) (int64, error) {
	if v.kind != kindNumber || strings.Contains(v.s, ".") {
		return 0, errors.New("expected signed decimal integer")
	}
	n, err := strconv.ParseInt(v.s, 10, 64)
	if err != nil {
		return 0, errors.New("integer outside int64")
	}
	return n, nil
}
func asRef(v jsonValue, zeroOK bool) (hashref.HashRef, error) {
	s, err := asString(v)
	if err != nil {
		return hashref.HashRef{}, err
	}
	if s == "" && zeroOK {
		return hashref.HashRef{}, nil
	}
	return hashref.Parse(s)
}

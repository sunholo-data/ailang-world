package approval

// Row 201 A2 (design §2.6): clientDataJSON, read token by token so a duplicated member is
// refused (Unmarshal would keep the last one silently). Unknown members are allowed, bounded.

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"unicode/utf8"
)

// MaxClientDataBytes bounds clientDataJSON (J1).
const MaxClientDataBytes = 2048

var (
	ErrClientData          = errors.New("approval: clientDataJSON is not a valid client data object")
	ErrClientDataDuplicate = errors.New("approval: clientDataJSON repeats a member")
)

// ClientData is the part of CollectedClientData the laws read (L3 §5.8.1).
type ClientData struct {
	Type, Challenge, Origin string
	CrossOrigin             bool
	HasTopOrigin            bool
}

// ParseClientData judges clientDataJSON (J1-J7).
func ParseClientData(b []byte) (ClientData, error) {
	if len(b) == 0 || len(b) > MaxClientDataBytes { // J1
		return ClientData{}, ErrClientData
	}
	if !utf8.Valid(b) { // J2
		return ClientData{}, ErrClientData
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') { // J3
		return ClientData{}, ErrClientData
	}
	var c ClientData
	seen := map[string]bool{}
	for dec.More() {
		tok, _ := dec.Token()
		key, ok := tok.(string)
		if !ok { // J3: not a member name (a failed read is nil)
			return ClientData{}, ErrClientData
		}
		strs := map[string]*string{"type": &c.Type, "challenge": &c.Challenge, "origin": &c.Origin}
		if _, known := strs[key]; !known && key != "crossOrigin" && key != "topOrigin" {
			if err := skipJSON(dec, 2); err != nil {
				return ClientData{}, err
			}
			continue
		}
		if seen[key] { // J4
			return ClientData{}, ErrClientDataDuplicate
		}
		seen[key] = true
		val, _ := dec.Token() // a failed read is nil: neither a string nor a bool
		s, isStr := val.(string)
		bv, isBool := val.(bool)
		switch {
		case strs[key] != nil:
			if !isStr { // J5
				return ClientData{}, ErrClientData
			}
			*strs[key] = s
		case key == "crossOrigin":
			if !isBool { // J6
				return ClientData{}, ErrClientData
			}
			c.CrossOrigin = bv
		default: // topOrigin
			if !isStr { // J6
				return ClientData{}, ErrClientData
			}
			c.HasTopOrigin = true
		}
	}
	if _, err := dec.Token(); err != nil { // J3: the closing brace
		return ClientData{}, ErrClientData
	}
	if _, err := dec.Token(); err != io.EOF { // J3: one object, nothing after it
		return ClientData{}, ErrClientData
	}
	if !seen["type"] || !seen["challenge"] || !seen["origin"] { // J5
		return ClientData{}, ErrClientData
	}
	return c, nil
}

// skipJSON reads one value of an unknown member (a member name is read as a value: the token
// stream is the same), refusing nesting over depth 4 (J7) and any token that does not read (J3).
func skipJSON(dec *json.Decoder, depth int) error {
	tok, err := dec.Token()
	if err != nil { // J3
		return ErrClientData
	}
	if tok != json.Delim('{') && tok != json.Delim('[') {
		return nil
	}
	if depth > 4 { // J7
		return ErrClientData
	}
	for dec.More() {
		if err := skipJSON(dec, depth+1); err != nil {
			return err
		}
	}
	if _, err := dec.Token(); err != nil { // J3: the closing delimiter
		return ErrClientData
	}
	return nil
}

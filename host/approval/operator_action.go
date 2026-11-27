package approval

// Row 243 V4 sprint 2 M1 (design_docs/planned/w-world-address-bar-verbs-v4.md §3.1): the operator action
// an enrolled key approves from the World app (host/verbs' OperatorIntent, canonicalized, never a re-parsed
// bar line) and the 97-byte message it signs, with its own schema, prefix and SSHSIG namespace. Pure: no
// I/O, no clock, no randomness. No production caller yet (S3 stores it, S7 serves it).

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// The operator action's schema, signed-message prefix and SSHSIG namespace; the message is
// 24 + 1 + 32 + 32 + 8 bytes; a signed request lives at most 600 s (design §3.1).
const (
	OperatorActionSchema     = "world/operator-action/v1"
	OperatorActionPrefix     = "world-operator-action/v1"
	OperatorActionNamespace  = "world-operator-action@sunholo"
	OperatorActionMessageLen = len(OperatorActionPrefix) + 1 + 32 + 32 + 8
	OperatorActionLifetime   = 10 * time.Minute
)

// MaxOperatorActionBytes bounds an action's canonical bytes; MaxOperatorTaskBytes bounds its task text,
// in UTF-8 bytes, separately (design §3.1).
const (
	MaxOperatorActionBytes = 128 << 10
	MaxOperatorTaskBytes   = 8 << 10
)

// operatorActionIDPattern is host/verbs' OperatorActionIDPattern (host/approval imports only the standard
// library; TestOperatorArgvIsTheBarsIntent joins the two); operatorResolvedName is a resolution fact's
// name; operatorPresetMax is host/verbs' MaxPresetName (the CLI's maxPresetNameLen).
var (
	operatorActionIDPattern = regexp.MustCompile(`^oa-[0-9a-f]{32}$`)
	operatorResolvedName    = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)
)

const operatorResolvedMax, operatorPresetMax = 64, 51

// OperatorOptions are the admitted options of OperatorIntent.
type OperatorOptions struct {
	Preset           string `json:"preset"`
	Headless         bool   `json:"headless"`           // start: --launch --headless
	ConfirmOnSurface bool   `json:"confirm_on_surface"` // finish
	Ready            bool   `json:"ready"`              // finish
}

// OperatorResolved is one frozen resolution fact (S5 names them); entries are in strictly ascending
// name order, so one resolution has one encoding.
type OperatorResolved struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// OperatorAction is one operator action. Its canonical JSON (CanonicalOperatorAction), in this field
// order, is what is digested, stored and signed.
type OperatorAction struct {
	Schema     string             `json:"schema"`     // OperatorActionSchema
	ID         string             `json:"id"`         // ^oa-[0-9a-f]{32}$; "" only inside OperatorReviewDigest
	Kind       string             `json:"kind"`       // "start" or "finish"
	Episode    string             `json:"episode"`    // ^[a-z0-9][a-z0-9-]{0,63}$
	Store      string             `json:"store"`      // the served credentials store, absolute and clean
	Argv       []string           `json:"argv"`       // exactly OperatorArgv(Kind, Episode, Options)
	Options    OperatorOptions    `json:"options"`    // the admitted options
	Resolution []OperatorResolved `json:"resolution"` // the frozen resolution, 1 to 64 facts
	Task       string             `json:"task"`       // start: 1 to 8 KiB of UTF-8; finish: ""
	Review     string             `json:"review"`     // the exact review text the card displays
}

var (
	ErrOperatorAction       = errors.New("approval: bad operator action")
	ErrOperatorTooLarge     = errors.New("approval: operator action is over 128 KiB")
	ErrOperatorTask         = errors.New("approval: operator task is empty, over 8 KiB or on a finish")
	ErrOperatorUTF8         = errors.New("approval: operator action is not valid UTF-8")
	ErrOperatorNUL          = errors.New("approval: operator action holds a NUL byte")
	ErrOperatorShape        = errors.New("approval: operator action is not a JSON object")
	ErrOperatorDuplicateKey = errors.New("approval: operator action repeats a key")
	ErrOperatorUnknownKey   = errors.New("approval: operator action has an unknown key")
	ErrOperatorTrailing     = errors.New("approval: operator action has data after its object")
)

// OperatorArgv is the canonical CLI line of an action: host/verbs AdmitOperator's Intent.Argv.
func OperatorArgv(kind, episode string, o OperatorOptions) []string {
	argv := []string{"ailang-worldd", "episode", kind, episode, "--preset", o.Preset}
	switch kind {
	case "start":
		argv = append(argv, "--launch")
		if o.Headless {
			argv = append(argv, "--headless")
		}
	case "finish":
		if o.ConfirmOnSurface {
			argv = append(argv, "--confirm-on-surface")
		}
		if o.Ready {
			argv = append(argv, "--ready")
		}
	}
	return argv
}

// Validate refuses an action the daemon or the runner must never act on. Each refusal wraps an Err
// above, ErrSchema, ErrString, ErrEpisode or ErrPath, and names its JSON field.
func (a OperatorAction) Validate() error { return a.validate(true) }

func (a OperatorAction) validate(withID bool) error {
	bad := func(format string, args ...any) error {
		return fmt.Errorf("%w: "+format, append([]any{ErrOperatorAction}, args...)...)
	}
	if a.Schema != OperatorActionSchema { // O1
		return fmt.Errorf("%w: schema is %q, want %q", ErrSchema, a.Schema, OperatorActionSchema)
	}
	if err := stringLeaves(reflect.ValueOf(a), "", func(p, s string) error { // O2
		if !utf8.ValidString(s) || strings.ContainsRune(s, 0) {
			return fmt.Errorf("%w: %s", ErrString, p)
		}
		return nil
	}); err != nil {
		return err
	}
	if withID && !operatorActionIDPattern.MatchString(a.ID) { // O3
		return bad("id %q does not match %s", a.ID, operatorActionIDPattern)
	}
	if a.Kind != "start" && a.Kind != "finish" { // O4
		return bad("kind is %q, want start or finish", a.Kind)
	}
	if !episodePattern.MatchString(a.Episode) { // O5
		return fmt.Errorf("%w: episode %q", ErrEpisode, a.Episode)
	}
	if !path.IsAbs(a.Store) || path.Clean(a.Store) != a.Store { // O6
		return fmt.Errorf("%w: store %q", ErrPath, a.Store)
	}
	if p := a.Options.Preset; p == "se-tools" || len(p) > operatorPresetMax || !episodePattern.MatchString(p) { // O7
		return bad("options.preset %q is not a launch preset name", p)
	}
	if a.Kind == "start" && (a.Options.ConfirmOnSurface || a.Options.Ready) { // O8
		return bad("options: a start sets only preset and headless")
	}
	if a.Kind == "finish" && a.Options.Headless {
		return bad("options: a finish has no headless")
	}
	if want := OperatorArgv(a.Kind, a.Episode, a.Options); !slices.Equal(a.Argv, want) { // O9
		return bad("argv is %q, want %q", a.Argv, want)
	}
	if a.Kind == "start" && (len(a.Task) < 1 || len(a.Task) > MaxOperatorTaskBytes) { // O10
		return fmt.Errorf("%w: task is %d bytes, want 1 to %d", ErrOperatorTask, len(a.Task), MaxOperatorTaskBytes)
	}
	if a.Kind == "finish" && a.Task != "" {
		return fmt.Errorf("%w: a finish carries no task", ErrOperatorTask)
	}
	if len(a.Resolution) < 1 || len(a.Resolution) > operatorResolvedMax { // O11
		return bad("resolution has %d facts, want 1 to %d", len(a.Resolution), operatorResolvedMax)
	}
	for i, r := range a.Resolution {
		if !operatorResolvedName.MatchString(r.Name) {
			return bad("resolution[%d].name %q", i, r.Name)
		}
		if i > 0 && a.Resolution[i-1].Name >= r.Name {
			return bad("resolution[%d].name %q is not after %q", i, r.Name, a.Resolution[i-1].Name)
		}
	}
	if a.Review == "" { // O12
		return bad("review is empty")
	}
	c, err := CanonicalOperatorAction(a)
	if err != nil {
		return err
	}
	if len(c) > MaxOperatorActionBytes { // O13
		return fmt.Errorf("%w: %d bytes, at most %d", ErrOperatorTooLarge, len(c), MaxOperatorActionBytes)
	}
	return nil
}

// CanonicalOperatorAction is the compact JSON of a in struct field order, HTML escaping off, no
// trailing newline (CanonicalVerb's settings). It does not call Validate.
func CanonicalOperatorAction(a OperatorAction) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(a); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

// DecodeOperatorAction is the strict inverse of CanonicalOperatorAction. Before any decoding it
// refuses over-limit input, invalid UTF-8, a NUL byte, a non-object, a repeated key and trailing data;
// then an unknown key, a type error, bytes that do not re-encode identically and an invalid action.
func DecodeOperatorAction(b []byte) (OperatorAction, error) {
	if len(b) > MaxOperatorActionBytes { // D1
		return OperatorAction{}, fmt.Errorf("%w: %d bytes, at most %d", ErrOperatorTooLarge, len(b), MaxOperatorActionBytes)
	}
	if !utf8.Valid(b) { // D2
		return OperatorAction{}, ErrOperatorUTF8
	}
	if bytes.IndexByte(b, 0) >= 0 { // D3
		return OperatorAction{}, ErrOperatorNUL
	}
	if err := strictObject(b); err != nil { // D4-D6
		return OperatorAction{}, err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var a OperatorAction
	if err := dec.Decode(&a); err != nil { // D7
		if strings.HasPrefix(err.Error(), "json: unknown field ") {
			return OperatorAction{}, fmt.Errorf("%w: %v", ErrOperatorUnknownKey, err)
		}
		return OperatorAction{}, fmt.Errorf("%w: %v", ErrDecode, err)
	}
	c, err := CanonicalOperatorAction(a)
	if err != nil {
		return OperatorAction{}, err
	}
	if !bytes.Equal(c, b) { // D8
		return OperatorAction{}, ErrNotCanonical
	}
	if err := a.Validate(); err != nil { // D9
		return OperatorAction{}, err
	}
	return a, nil
}

// DecodeOperatorObject reads a closed, single JSON object without duplicate
// keys or case-folded aliases. The caller supplies its body and field bounds.
func DecodeOperatorObject(b []byte, dst any, keys ...string) error {
	if !utf8.Valid(b) {
		return ErrOperatorUTF8
	}
	if bytes.IndexByte(b, 0) >= 0 {
		return ErrOperatorNUL
	}
	if err := strictObject(b); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}
	for name := range fields {
		if !slices.Contains(keys, name) {
			return ErrOperatorUnknownKey
		}
	}
	return json.Unmarshal(b, dst)
}

// strictObject walks b's tokens (host/daemon strictRosterJSON's form): one JSON object (D4), no key
// twice in any object (D5) and nothing after it (D6). encoding/json alone keeps the last of a repeated
// key and stops reading after the first value.
func strictObject(b []byte) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	var walk func(tok json.Token) error
	walk = func(tok json.Token) error {
		d, ok := tok.(json.Delim)
		if !ok {
			return nil
		}
		seen := map[string]bool{}
		for dec.More() {
			if d == '{' {
				key, err := dec.Token()
				if err != nil {
					return fmt.Errorf("%w: %v", ErrDecode, err)
				}
				k, _ := key.(string)
				if seen[k] { // D5
					return fmt.Errorf("%w: %q", ErrOperatorDuplicateKey, k)
				}
				seen[k] = true
			}
			next, err := dec.Token()
			if err != nil {
				return fmt.Errorf("%w: %v", ErrDecode, err)
			}
			if err := walk(next); err != nil {
				return err
			}
		}
		if _, err := dec.Token(); err != nil { // the closing delimiter
			return fmt.Errorf("%w: %v", ErrDecode, err)
		}
		return nil
	}
	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') { // D4
		return ErrOperatorShape
	} else if err := walk(tok); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF { // D6
		return ErrOperatorTrailing
	}
	return nil
}

// OperatorReviewDigest is the SHA-256 of a's canonical bytes with its id blank: the semantic inputs
// before an action id or challenge is assigned (design §3.1), the key S3 deduplicates intents on. It
// validates every field but the id.
func OperatorReviewDigest(a OperatorAction) ([32]byte, error) {
	a.ID = ""
	if err := a.validate(false); err != nil {
		return [32]byte{}, err
	}
	c, err := CanonicalOperatorAction(a)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(c), nil
}

// OperatorSignatureDigest is the SHA-256 of the final envelope's canonical bytes, id included: the
// digest OperatorActionMessage signs. It validates a.
func OperatorSignatureDigest(a OperatorAction) ([32]byte, error) {
	if err := a.Validate(); err != nil {
		return [32]byte{}, err
	}
	c, err := CanonicalOperatorAction(a)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(c), nil
}

// OperatorActionMessage is the 97-byte message a confirmation signs (design §3.1):
// OperatorActionPrefix ‖ 0x00 ‖ signature digest ‖ challenge ‖ expiry (uint64, big-endian).
func OperatorActionMessage(digest, challenge [32]byte, expiry int64) ([]byte, error) {
	if expiry < 0 {
		return nil, ErrNegativeExpiry
	}
	m := make([]byte, 0, OperatorActionMessageLen)
	m = append(m, OperatorActionPrefix...)
	m = append(m, 0)
	m = append(m, digest[:]...)
	m = append(m, challenge[:]...)
	m = binary.BigEndian.AppendUint64(m, uint64(expiry))
	return m, nil
}

// OperatorActionMessageHash is the SHA-256 of OperatorActionMessage: a passkey's WebAuthn challenge.
func OperatorActionMessageHash(digest, challenge [32]byte, expiry int64) ([32]byte, error) {
	m, err := OperatorActionMessage(digest, challenge, expiry)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(m), nil
}

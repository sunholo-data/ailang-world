package approval

// Row 243 V3 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.2): the verb request a
// confirm-on-surface `episode finish` stores, and the 86-byte message an enrolled key signs to
// confirm it. A third kind beside the mint message (SignedMessage) and the decision message
// (DecisionMessage): its own prefix and SSHSIG namespace, so no signature of one kind is accepted
// as another (design §2.2 "Domain separation"). Pure: no I/O, no clock, no randomness.

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
	"unicode/utf8"
)

// VerbRequestSchema is the schema of every stored verb request.
const VerbRequestSchema = "world/verb-request/v1"

// VerbRequestLifetime is a verb request's lifetime, as a mint request's (RequestLifetime).
const VerbRequestLifetime = 10 * time.Minute

// MaxVerbRequestBytes bounds a verb request's canonical bytes: a body of at most 60,000 bytes
// (cmd/ailang-worldd reportMarkdownMax) plus the prompt, scope, title and argv.
const MaxVerbRequestBytes = 128 << 10

// VerbFinish is the only admissible verb.
const VerbFinish = "episode finish"

// VerbPrefix starts every signed verb message (design §2.2).
const VerbPrefix = "world-verb/v1"

// VerbNamespace is the SSHSIG namespace every verb confirmation is signed under.
const VerbNamespace = "world-verb@sunholo"

// VerbMessageLen is the length of VerbMessage: 13 + 1 + 32 + 32 + 8.
const VerbMessageLen = len(VerbPrefix) + 1 + 32 + 32 + 8

// verbBodyMax is the pull request body's cap (cmd/ailang-worldd episode_markdown.go reportMarkdownMax).
const verbBodyMax = 60000

// verbArgvMax bounds the canonical argv's words.
const verbArgvMax = 64

// VerbRequest is a stored, single-use confirmation request. Its canonical JSON (CanonicalVerb), in
// this field order, is what is digested, stored and shown on the card.
type VerbRequest struct {
	Schema    string   `json:"schema"`    // VerbRequestSchema
	Verb      string   `json:"verb"`      // VerbFinish, the only admissible value
	Argv      []string `json:"argv"`      // the canonical argv: ailang-worldd episode finish <episode> …
	Episode   string   `json:"episode"`   // ^[a-z0-9][a-z0-9-]{0,63}$
	Prompt    string   `json:"prompt"`    // the terminal's lines (broker.FinishPromptText), without the question
	Scope     string   `json:"scope"`     // the finish approval scope the daemon records
	Title     string   `json:"title"`     // the pull request's title
	Body      string   `json:"body"`      // the full pull request body
	Committer string   `json:"committer"` // the pinned committer, raw ("" allowed)
	Protected []string `json:"protected"` // the protected paths the commit changed
}

// ErrVerbRequest is a verb request that Validate refuses; the text names the field.
var ErrVerbRequest = errors.New("approval: bad verb request")

// Validate refuses a verb request the daemon or the CLI must never act on (design §2.2). Each
// refusal wraps ErrSchema, ErrString, ErrEpisode or ErrVerbRequest and names its JSON field.
func (r VerbRequest) Validate() error {
	if r.Schema != VerbRequestSchema { // V1
		return fmt.Errorf("%w: schema is %q, want %q", ErrSchema, r.Schema, VerbRequestSchema)
	}
	if err := stringLeaves(reflect.ValueOf(r), "", func(path, s string) error { // V2
		if !utf8.ValidString(s) {
			return fmt.Errorf("%w: %s", ErrString, path)
		}
		return nil
	}); err != nil {
		return err
	}
	if r.Verb != VerbFinish { // V3
		return fmt.Errorf("%w: verb is %q, want %q", ErrVerbRequest, r.Verb, VerbFinish)
	}
	if !episodePattern.MatchString(r.Episode) { // V4
		return fmt.Errorf("%w: episode %q", ErrEpisode, r.Episode)
	}
	if len(r.Argv) < 4 || r.Argv[0] != "ailang-worldd" || r.Argv[1] != "episode" || r.Argv[2] != "finish" || r.Argv[3] != r.Episode { // V5
		return fmt.Errorf("%w: argv does not start ailang-worldd episode finish %s", ErrVerbRequest, r.Episode)
	}
	if len(r.Argv) > verbArgvMax { // V6
		return fmt.Errorf("%w: argv has %d words, at most %d", ErrVerbRequest, len(r.Argv), verbArgvMax)
	}
	for i, w := range r.Argv { // V7
		for j := 0; j < len(w); j++ {
			if w[j] < 0x20 || w[j] == 0x7f {
				return fmt.Errorf("%w: argv[%d] holds a control byte", ErrVerbRequest, i)
			}
		}
	}
	for _, f := range []struct{ name, val string }{{"prompt", r.Prompt}, {"scope", r.Scope}, {"title", r.Title}} { // V8
		if f.val == "" {
			return fmt.Errorf("%w: %s is empty", ErrVerbRequest, f.name)
		}
	}
	if len(r.Body) > verbBodyMax { // V9
		return fmt.Errorf("%w: body is %d bytes, at most %d", ErrVerbRequest, len(r.Body), verbBodyMax)
	}
	c, err := CanonicalVerb(r)
	if err != nil {
		return err
	}
	if len(c) > MaxVerbRequestBytes { // V10
		return fmt.Errorf("%w: the request is %d bytes, at most %d", ErrVerbRequest, len(c), MaxVerbRequestBytes)
	}
	return nil
}

// CanonicalVerb is the compact JSON of r in struct field order, with HTML escaping off and no
// trailing newline (Canonical's encoder settings). It does not call Validate.
func CanonicalVerb(r VerbRequest) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

// DecodeVerb is the strict inverse of CanonicalVerb (Decode's strictness): unknown fields, trailing
// data and any input that does not re-encode to the same bytes are refused. It does not call
// Validate; callers do.
func DecodeVerb(b []byte) (VerbRequest, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var r VerbRequest
	if err := dec.Decode(&r); err != nil {
		return VerbRequest{}, fmt.Errorf("%w: %v", ErrDecode, err)
	}
	c, err := CanonicalVerb(r)
	if err != nil {
		return VerbRequest{}, err
	}
	if !bytes.Equal(c, b) {
		return VerbRequest{}, ErrNotCanonical
	}
	return r, nil
}

// VerbMessage is the 86-byte message a confirmation signs (design §2.2):
// VerbPrefix ‖ 0x00 ‖ digest ‖ challenge ‖ expiry (uint64, big-endian).
func VerbMessage(digest, challenge [32]byte, expiry int64) ([]byte, error) {
	if expiry < 0 {
		return nil, ErrNegativeExpiry
	}
	m := make([]byte, 0, VerbMessageLen)
	m = append(m, VerbPrefix...)
	m = append(m, 0)
	m = append(m, digest[:]...)
	m = append(m, challenge[:]...)
	m = binary.BigEndian.AppendUint64(m, uint64(expiry))
	return m, nil
}

// VerbMessageHash is the SHA-256 of VerbMessage: a passkey's WebAuthn challenge.
func VerbMessageHash(digest, challenge [32]byte, expiry int64) ([32]byte, error) {
	m, err := VerbMessage(digest, challenge, expiry)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(m), nil
}

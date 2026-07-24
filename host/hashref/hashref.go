// Package hashref implements the tagged HashRef value type from Decision 3 of
// the w-world-library-m1 design: a content-address of the form "algo:digest"
// where algo is a lowercase registered algorithm tag and digest is the
// lowercase hexadecimal digest.
//
// Parsing, rendering, and hash calculation are centralized here so that every
// digest-bearing field across the kernel shares one identity representation.
// M1 registers exactly one algorithm, sha256, backed by crypto/sha256.
package hashref

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// AlgoSHA256 is the canonical tag for the SHA-256 algorithm registered in M1.
const AlgoSHA256 = "sha256"

// algoSpec describes a registered algorithm: how many lowercase hex characters
// its digest occupies and how to compute it over payload bytes.
type algoSpec struct {
	// digestHexLen is the exact number of lowercase hex characters a digest
	// for this algorithm must have (2 chars per output byte).
	digestHexLen int
	// sum computes the raw digest bytes for the given payload.
	sum func(payload []byte) []byte
}

// registry maps a lowercase algorithm tag to its specification. M1 registers
// only sha256; the map keeps future tags coexisting behind one dispatcher.
var registry = map[string]algoSpec{
	AlgoSHA256: {
		digestHexLen: sha256.Size * 2, // 64
		sum: func(payload []byte) []byte {
			s := sha256.Sum256(payload)
			return s[:]
		},
	},
}

// HashError is the structured error returned for every malformed, unsupported,
// or non-canonical HashRef the readers reject.
type HashError struct {
	// Input is the offending text (empty when the failure is not text-driven).
	Input string
	// Reason is a stable, human-readable explanation.
	Reason string
}

func (e *HashError) Error() string {
	if e.Input == "" {
		return "hashref: " + e.Reason
	}
	return fmt.Sprintf("hashref: %s: %q", e.Reason, e.Input)
}

// HashRef is a tagged content address. The zero value is invalid; construct
// values through Parse, New, or the Sum* helpers so the canonical invariants
// always hold: Algo is a registered lowercase tag and Digest is lowercase hex
// of the exact width for that algorithm.
type HashRef struct {
	algo   string
	digest string
}

// Algo returns the lowercase algorithm tag (for example "sha256").
func (h HashRef) Algo() string { return h.algo }

// Digest returns the lowercase hexadecimal digest.
func (h HashRef) Digest() string { return h.digest }

// IsZero reports whether h is the invalid zero value.
func (h HashRef) IsZero() bool { return h.algo == "" && h.digest == "" }

// String renders the canonical "algo:digest" text form. It returns the empty
// string for the zero value.
func (h HashRef) String() string {
	if h.IsZero() {
		return ""
	}
	return h.algo + ":" + h.digest
}

// isLowerHex reports whether s is non-empty and consists only of the lowercase
// hexadecimal characters 0-9 and a-f. Uppercase hex is deliberately rejected
// so canonical text has exactly one spelling per digest.
func isLowerHex(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') {
			continue
		}
		return false
	}
	return true
}

// New constructs a HashRef from an algorithm tag and a lowercase hex digest,
// validating both against the registry. It returns a *HashError if the tag is
// unsupported or the digest is not canonical lowercase hex of the expected
// width for that algorithm.
func New(algo, digest string) (HashRef, error) {
	spec, ok := registry[algo]
	if !ok {
		return HashRef{}, &HashError{Input: algo, Reason: "unsupported algorithm tag"}
	}
	if !isLowerHex(digest) {
		return HashRef{}, &HashError{Input: digest, Reason: "digest is not lowercase hexadecimal"}
	}
	if len(digest) != spec.digestHexLen {
		return HashRef{}, &HashError{
			Input:  digest,
			Reason: fmt.Sprintf("digest length %d does not match %s width %d", len(digest), algo, spec.digestHexLen),
		}
	}
	return HashRef{algo: algo, digest: digest}, nil
}

// Parse reads a canonical "algo:digest" HashRef. It rejects, with a structured
// *HashError: empty text, missing or extra colons, bare digests (no tag),
// unsupported tags, uppercase hex, and digests of the wrong width.
func Parse(text string) (HashRef, error) {
	if text == "" {
		return HashRef{}, &HashError{Reason: "empty hashref text"}
	}
	// Exactly one separator: a bare digest (no colon) and an over-segmented
	// string are both malformed.
	i := strings.IndexByte(text, ':')
	if i < 0 {
		return HashRef{}, &HashError{Input: text, Reason: "missing algorithm tag (expected algo:digest)"}
	}
	algo := text[:i]
	digest := text[i+1:]
	if strings.IndexByte(digest, ':') >= 0 {
		return HashRef{}, &HashError{Input: text, Reason: "malformed hashref (multiple colons)"}
	}
	if algo == "" {
		return HashRef{}, &HashError{Input: text, Reason: "empty algorithm tag"}
	}
	return New(algo, digest)
}

// MustParse is Parse for compile-time-known constants; it panics on error.
// Intended for tests and package-level golden values, not runtime input.
func MustParse(text string) HashRef {
	h, err := Parse(text)
	if err != nil {
		panic(err)
	}
	return h
}

// Sum computes the HashRef of payload under the named algorithm. It returns a
// *HashError if algo is not registered.
func Sum(algo string, payload []byte) (HashRef, error) {
	spec, ok := registry[algo]
	if !ok {
		return HashRef{}, &HashError{Input: algo, Reason: "unsupported algorithm tag"}
	}
	digest := hex.EncodeToString(spec.sum(payload))
	return HashRef{algo: algo, digest: digest}, nil
}

// SumSHA256 computes the sha256-tagged HashRef of payload. It is the M1
// convenience entry point and cannot fail.
func SumSHA256(payload []byte) HashRef {
	h, err := Sum(AlgoSHA256, payload)
	if err != nil {
		// sha256 is always registered; a failure here is a programming error.
		panic(err)
	}
	return h
}

// Supported reports whether algo is a registered algorithm tag.
func Supported(algo string) bool {
	_, ok := registry[algo]
	return ok
}

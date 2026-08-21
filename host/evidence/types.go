// Package evidence defines the untrusted proof-evidence wire codecs.
package evidence

import (
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

const (
	MaxBytes              = 256 * 1024
	MaxStringBytes        = 1024
	MaxVerifiedIdentities = 256
	ProofReportSchemaV1   = "world/proof-report/v1"
)

// DecodeRefusal is returned only when untrusted bytes cannot be admitted by a
// strict codec. Kind is stable and permits callers to distinguish the guard
// that refused the input from an operational error.
type DecodeRefusal struct {
	Kind   string
	Detail string
}

func (e *DecodeRefusal) Error() string {
	if e.Detail == "" {
		return "evidence: decode refused: " + e.Kind
	}
	return fmt.Sprintf("evidence: decode refused: %s: %s", e.Kind, e.Detail)
}

const (
	RefusalOversize     = "oversize"
	RefusalInvalidUTF8  = "invalid_utf8"
	RefusalMalformed    = "malformed"
	RefusalNonCanonical = "non_canonical"
	RefusalLimit        = "limit"
)

// ProofReportV1 is the decoded form of the canonical nine-member report.
type ProofReportV1 struct {
	Schema          string
	Subject         hashref.HashRef
	Compiler        hashref.HashRef
	CompilerVersion string
	Verified        []string
	Errors          int
	Counterexamples int
	CheckPassed     bool
	ProofSucceeded  bool
}

// AuthenticatedEnvelope carries report bytes and the supplied tag. MacValid
// reports only whether mac had the canonical base64url spelling and 32-byte
// width; an invalid tag is deliberately data for the authentication step.
type AuthenticatedEnvelope struct {
	Report   []byte
	MAC      []byte
	MACValid bool
}

// ClaimedEvidence is decoded but untrusted proposal data. Its representation
// is intentionally opaque; only validation in a later milestone may confer
// authority on it.
type ClaimedEvidence struct{ canonical string }

// IsZero distinguishes a successfully decoded claim from no claim.
func (c ClaimedEvidence) IsZero() bool { return c.canonical == "" }

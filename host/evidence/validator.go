package evidence

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const ProofSemanticID = ProofReportSchemaV1

var (
	// InterfaceHashV1 is the fixed identity of the canonical authenticated-envelope interface.
	InterfaceHashV1           = hashref.SumSHA256([]byte("world/authenticated-proof-envelope/v1"))
	proofInterfaceHash        = InterfaceHashV1
	ErrInvalidValidatorConfig = errors.New("evidence: invalid validator configuration")
	ErrUnorderedTimeouts      = errors.New("evidence: unordered timeouts")
	ErrUnmintedAuthority      = errors.New("evidence: unminted authority")
	ErrForeignSeal            = errors.New("evidence: foreign seal")
)

type UnsupportedReason string

const (
	UnsupportedInvalidRef            UnsupportedReason = "invalid_ref"
	UnsupportedMissing               UnsupportedReason = "missing"
	UnsupportedOversize              UnsupportedReason = "oversize"
	UnsupportedHashMismatch          UnsupportedReason = "hash_mismatch"
	UnsupportedWrongSemanticID       UnsupportedReason = "wrong_semantic_id"
	UnsupportedWrongInterface        UnsupportedReason = "wrong_interface"
	UnsupportedMalformed             UnsupportedReason = "malformed"
	UnsupportedUnauthenticatedReport UnsupportedReason = "unauthenticated_report"
	UnsupportedSubjectMismatch       UnsupportedReason = "subject_mismatch"
	UnsupportedToolMismatch          UnsupportedReason = "tool_mismatch"
	UnsupportedProofFailed           UnsupportedReason = "proof_failed"
	UnsupportedProofIncomplete       UnsupportedReason = "proof_incomplete"
)

type CompilerConfig struct {
	Compiler          hashref.HashRef
	CompilerVersion   string
	ObjectReadTimeout time.Duration
}

type ObjectReader interface {
	ReadObject(context.Context, hashref.HashRef, int64) (store.ObjectMeta, []byte, error)
	BusyTimeout() time.Duration
}

type Validator struct {
	key                *[32]byte
	id                 *[32]byte
	reader             ObjectReader
	compiler           CompilerConfig
	requiredIdentities []string
}

func NewValidator(key [32]byte, reader ObjectReader, cfg CompilerConfig, requiredIdentities []string) (*Validator, error) {
	if reader == nil || cfg.ObjectReadTimeout <= 0 || cfg.Compiler.IsZero() || cfg.CompilerVersion == "" {
		return nil, fmt.Errorf("%w: reader, compiler identity, compiler version, and positive ObjectReadTimeout are required", ErrInvalidValidatorConfig)
	}
	if len(requiredIdentities) == 0 {
		return nil, fmt.Errorf("%w: required identities are empty", ErrInvalidValidatorConfig)
	}
	window := reader.BusyTimeout()
	if window < 0 {
		return nil, fmt.Errorf("%w: BusyTimeout is unknown (%s)", ErrInvalidValidatorConfig, window)
	}
	if window > 0 && cfg.ObjectReadTimeout <= window {
		return nil, fmt.Errorf("%w: ObjectReadTimeout %s must exceed BusyTimeout %s", ErrUnorderedTimeouts, cfg.ObjectReadTimeout, window)
	}
	keyCopy := new([32]byte)
	*keyCopy = key
	return &Validator{key: keyCopy, id: keyCopy, reader: reader, compiler: cfg, requiredIdentities: append([]string(nil), requiredIdentities...)}, nil
}

type ValidatedEvidence struct {
	reportRef hashref.HashRef
	subject   hashref.HashRef
	mintedBy  *[32]byte
}

type ValidationResult struct {
	validated   ValidatedEvidence
	hasValue    bool
	unsupported UnsupportedReason
	err         error
}

func (r ValidationResult) Validated() (ValidatedEvidence, bool) { return r.validated, r.hasValue }
func (r ValidationResult) Unsupported() (UnsupportedReason, bool) {
	return r.unsupported, r.unsupported != ""
}
func (r ValidationResult) Err() error { return r.err }

func unsupported(reason UnsupportedReason) ValidationResult {
	return ValidationResult{unsupported: reason}
}

func (v *Validator) ValidateProof(ctx context.Context, reportRef, expectedSubject hashref.HashRef) ValidationResult {
	if _, err := hashref.Parse(reportRef.String()); err != nil {
		return unsupported(UnsupportedInvalidRef)
	}
	readCtx, cancel := context.WithTimeout(ctx, v.compiler.ObjectReadTimeout)
	defer cancel()
	meta, payload, err := v.reader.ReadObject(readCtx, reportRef, MaxBytes)
	if err != nil {
		var tooLarge *store.ObjectTooLargeError
		if errors.As(err, &tooLarge) {
			return unsupported(UnsupportedOversize)
		}
		return ValidationResult{err: fmt.Errorf("evidence: read proof report: %w", err)}
	}
	if payload == nil {
		return unsupported(UnsupportedMissing)
	}
	if got := hashref.SumSHA256(payload); got != reportRef {
		return unsupported(UnsupportedHashMismatch)
	}
	if meta.SemanticID != ProofSemanticID {
		return unsupported(UnsupportedWrongSemanticID)
	}
	if meta.InterfaceHash != proofInterfaceHash {
		return unsupported(UnsupportedWrongInterface)
	}
	envelope, err := DecodeAuthenticatedEnvelope(payload)
	if err != nil {
		return unsupported(UnsupportedMalformed)
	}
	report, err := DecodeProofReportV1(envelope.Report)
	if err != nil {
		return unsupported(UnsupportedMalformed)
	}
	mac := hmac.New(sha256.New, v.key[:])
	_, _ = mac.Write(envelope.Report)
	want := mac.Sum(nil)
	if len(envelope.MAC) != sha256.Size || !hmac.Equal(want, envelope.MAC) {
		return unsupported(UnsupportedUnauthenticatedReport)
	}
	if report.Subject != expectedSubject {
		return unsupported(UnsupportedSubjectMismatch)
	}
	toolMismatch := report.Compiler != v.compiler.Compiler || report.CompilerVersion != v.compiler.CompilerVersion
	if toolMismatch {
		return unsupported(UnsupportedToolMismatch)
	}
	proofFailed := !report.CheckPassed || !report.ProofSucceeded || report.Errors != 0 || report.Counterexamples != 0
	if proofFailed {
		return unsupported(UnsupportedProofFailed)
	}
	incomplete := len(report.Verified) == 0
	for _, required := range v.requiredIdentities {
		if !containsIdentity(report.Verified, required) {
			incomplete = true
		}
	}
	if incomplete {
		return unsupported(UnsupportedProofIncomplete)
	}
	return ValidationResult{validated: ValidatedEvidence{reportRef: reportRef, subject: expectedSubject, mintedBy: v.id}, hasValue: true}
}

func containsIdentity(sorted []string, want string) bool {
	i, j := 0, len(sorted)
	for i < j {
		h := int(uint(i+j) >> 1)
		if sorted[h] < want {
			i = h + 1
		} else {
			j = h
		}
	}
	return i < len(sorted) && sorted[i] == want
}

type ResolvedGrade uint8

const ResolvedGradeProven ResolvedGrade = 1

type ResolutionResult struct {
	grade ResolvedGrade
	ok    bool
	err   error
}

func (r ResolutionResult) Proven() (ResolvedGrade, bool) { return r.grade, r.ok }
func (r ResolutionResult) Err() error                    { return r.err }

func (v *Validator) Resolve(sealed ValidatedEvidence) ResolutionResult {
	if v == nil || v.id == nil || sealed.mintedBy == nil {
		return ResolutionResult{err: ErrUnmintedAuthority}
	}
	if sealed.mintedBy != v.id {
		return ResolutionResult{err: ErrForeignSeal}
	}
	// A copied validator deliberately remains the same authority. Likewise, a
	// caller may mint for a validator it constructed; no library can stop a
	// caller lying to itself. The enforced boundary is between identities.
	return ResolutionResult{grade: ResolvedGradeProven, ok: true}
}

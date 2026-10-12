package daemon

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/workbench"
)

const proofKeySize = 32

const gradeCandidateScanLimit = 100

const ( // refusals of configureGrader (row 105 D6, PD4)
	gradeConfigIdentitiesWithoutKey = "--proof-required-identity needs --proof-key-file"
	gradeConfigNoInterpreter        = "the workbench grade validator needs the pinned interpreter; --ailang-bin is unset"
	gradeConfigKeyRefused           = "cannot load the proof key file"
	gradeConfigValidatorRefused     = "cannot construct the workbench grade validator"
)

// proofValidator is the daemon's validator seam (row 105 D2): *evidence.Validator
// satisfies it, and New wires nothing else.
type proofValidator interface {
	ValidateProof(ctx context.Context, reportRef, expectedSubject hashref.HashRef) evidence.ValidationResult
	Resolve(sealed evidence.ValidatedEvidence) evidence.ResolutionResult
}

// proofCompilerVersionLine mirrors the producer's derivation exactly
// (host/evidence/proof_producer.go:328): the RAW first line, trailing "\r"
// stripped, and deliberately NOT releaseFromVersion's TrimSpace.
func proofCompilerVersionLine(verbatim string) string {
	return strings.TrimRight(strings.SplitN(verbatim, "\n", 2)[0], "\r")
}

// loadProofKey reads exactly proofKeySize bytes from a regular file. It stats
// the OPEN file (no path TOCTOU) and reads at most proofKeySize+1 bytes.
// Errors name the path and sizes, never the key bytes.
func loadProofKey(path string) ([proofKeySize]byte, error) {
	var key [proofKeySize]byte
	f, err := os.Open(path)
	if err != nil {
		return key, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return key, err
	}
	if !info.Mode().IsRegular() {
		return key, fmt.Errorf("%s is not a regular file", path)
	}
	raw, err := io.ReadAll(io.LimitReader(f, proofKeySize+1))
	if err != nil {
		return key, err
	}
	if len(raw) > proofKeySize {
		return key, fmt.Errorf("%s holds more than %d bytes; a proof key is exactly %d", path, proofKeySize, proofKeySize)
	}
	if len(raw) < proofKeySize {
		return key, fmt.Errorf("%s holds %d bytes; a proof key is exactly %d", path, len(raw), proofKeySize)
	}
	copy(key[:], raw)
	return key, nil
}

// configureGrader wires the workbench grade validator (row 105 D6), or leaves
// d.grader nil when no key file is configured. Every half-configuration refuses.
func (d *Daemon) configureGrader(cfg Config) error {
	if cfg.ProofKeyFile == "" {
		if len(cfg.ProofRequiredIdentities) > 0 {
			return d.abort(StageConfig, gradeConfigIdentitiesWithoutKey, nil)
		}
		return nil
	}
	if d.interpreterHash.IsZero() {
		return d.abort(StageConfig, gradeConfigNoInterpreter, nil)
	}
	key, err := loadProofKey(cfg.ProofKeyFile)
	if err != nil {
		return d.abort(StageConfig, gradeConfigKeyRefused, err)
	}
	v, err := evidence.NewValidator(key, d.store, evidence.CompilerConfig{
		Compiler:          d.interpreterHash,
		CompilerVersion:   proofCompilerVersionLine(d.interpreterVersion),
		ObjectReadTimeout: readDeadline,
	}, cfg.ProofRequiredIdentities)
	if err != nil {
		return d.abort(StageConfig, gradeConfigValidatorRefused, err)
	}
	d.grader = v
	return nil
}

const ( // frozen reason sentences (design D4 table, byte-exact)
	gradeReasonNoStoredEvidence    = "no stored evidence: the store holds no object with semantic ID world/proof-report/v1"
	gradeReasonBoundElsewhere      = "stored proof envelopes exist; none is bound to this subject"
	gradeReasonToolMismatch        = "subject-bound evidence authenticates under the daemon's key but names a compiler other than the daemon's pinned interpreter"
	gradeReasonProofFailed         = "subject-bound evidence authenticates under the daemon's key but records a failed proof outcome"
	gradeReasonProofIncomplete     = "subject-bound evidence authenticates under the daemon's key but does not verify every identity the daemon requires"
	gradeReasonUnauthenticated     = "a stored proof envelope does not authenticate under the daemon's key"
	gradeReasonStructuralCandidate = "a stored candidate object does not satisfy the proof envelope's identity, content or schema checks"
	gradeReasonNoneValidated       = "stored proof envelopes exist and none validates for this object"
	gradeScanTruncation            = " (scan bounded to the first 100 stored candidates)"
)

var errGradeNoVerdict = errors.New("daemon: the grade validator returned no verdict")
var errGradeResolverRefused = errors.New("daemon: the grade resolver refused the validator's own seal")

var gradeReasonSentences = map[evidence.UnsupportedReason]string{
	evidence.UnsupportedUnauthenticatedReport: gradeReasonUnauthenticated,
	evidence.UnsupportedToolMismatch:          gradeReasonToolMismatch,
	evidence.UnsupportedProofFailed:           gradeReasonProofFailed,
	evidence.UnsupportedProofIncomplete:       gradeReasonProofIncomplete,
	evidence.UnsupportedInvalidRef:            gradeReasonStructuralCandidate,
	evidence.UnsupportedMissing:               gradeReasonStructuralCandidate,
	evidence.UnsupportedOversize:              gradeReasonStructuralCandidate,
	evidence.UnsupportedHashMismatch:          gradeReasonStructuralCandidate,
	evidence.UnsupportedWrongSemanticID:       gradeReasonStructuralCandidate,
	evidence.UnsupportedWrongInterface:        gradeReasonStructuralCandidate,
	evidence.UnsupportedMalformed:             gradeReasonStructuralCandidate,
}

// gradeReasonSentence maps a non-mismatch validator reason to its frozen
// sentence; any token outside the map renders the branch invariant.
func gradeReasonSentence(reason evidence.UnsupportedReason) string {
	if sentence, ok := gradeReasonSentences[reason]; ok {
		return sentence
	}
	return gradeReasonNoneValidated
}

// objectGradeView projects the object page's grade (row 105 D4). Unconfigured,
// it is the constant with zero reads. Configured, it scans at most
// gradeCandidateScanLimit proof envelopes (fetching one extra row to prove
// truncation) and renders PROVEN only through ValidateProof -> Resolve.
func (d *Daemon) objectGradeView(ctx context.Context, ref hashref.HashRef) (workbench.GradeView, error) {
	if d.grader == nil {
		return workbench.NewGradeUnavailable(objectGradeUnavailableReason), nil
	}
	candidates, err := d.reads.ObjectsBySemanticID(ctx, evidence.ProofSemanticID, "", gradeCandidateScanLimit+1)
	if err != nil {
		return workbench.GradeView{}, err
	}
	truncated := len(candidates) == gradeCandidateScanLimit+1
	if truncated {
		candidates = candidates[:gradeCandidateScanLimit]
	}
	if len(candidates) == 0 {
		return workbench.NewGradeUnavailable(gradeReasonNoStoredEvidence), nil
	}
	firstFailure := ""
	for _, candidate := range candidates {
		res := d.grader.ValidateProof(ctx, candidate.Hash, ref)
		if err := res.Err(); err != nil {
			return workbench.GradeView{}, fmt.Errorf("daemon: validate proof candidate %s: %w", candidate.Hash, err)
		}
		if sealed, ok := res.Validated(); ok {
			resolved := d.grader.Resolve(sealed)
			if _, proven := resolved.Proven(); !proven {
				return workbench.GradeView{}, fmt.Errorf("%w: candidate %s: %v", errGradeResolverRefused, candidate.Hash, resolved.Err())
			}
			return workbench.GradeView{Available: true, Label: workbench.GradePROVEN}, nil
		}
		reason, unsupported := res.Unsupported()
		if !unsupported {
			return workbench.GradeView{}, fmt.Errorf("%w: candidate %s", errGradeNoVerdict, candidate.Hash)
		}
		if reason != evidence.UnsupportedSubjectMismatch && firstFailure == "" {
			firstFailure = gradeReasonSentence(reason)
		}
	}
	sentence := gradeReasonBoundElsewhere
	if firstFailure != "" {
		sentence = firstFailure
	}
	if truncated {
		sentence += gradeScanTruncation
	}
	return workbench.NewGradeUnavailable(sentence), nil
}

package daemon

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

const proofKeySize = 32

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

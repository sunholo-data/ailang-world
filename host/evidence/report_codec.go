package evidence

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

var reportFields = [...]string{"schema", "subject", "compiler", "compilerVersion", "verified", "errors", "counterexamples", "checkPassed", "proofSucceeded"}

// EncodeProofReportV1 returns the sole canonical JSON spelling of r.
func EncodeProofReportV1(r ProofReportV1) ([]byte, error) {
	if err := validateReport(r); err != nil {
		return nil, err
	}
	type wire struct {
		Schema          string   `json:"schema"`
		Subject         string   `json:"subject"`
		Compiler        string   `json:"compiler"`
		CompilerVersion string   `json:"compilerVersion"`
		Verified        []string `json:"verified"`
		Errors          int      `json:"errors"`
		Counterexamples int      `json:"counterexamples"`
		CheckPassed     bool     `json:"checkPassed"`
		ProofSucceeded  bool     `json:"proofSucceeded"`
	}
	return json.Marshal(wire{r.Schema, r.Subject.String(), r.Compiler.String(), r.CompilerVersion, r.Verified, r.Errors, r.Counterexamples, r.CheckPassed, r.ProofSucceeded})
}

// DecodeProofReportV1 strictly decodes and proves canonicality by re-encoding.
func DecodeProofReportV1(raw []byte) (ProofReportV1, error) {
	if len(raw) > MaxBytes {
		return ProofReportV1{}, refusal(RefusalOversize, "report is %d bytes; limit is %d", len(raw), MaxBytes)
	}
	ms, err := objectMembers(raw)
	if err != nil {
		return ProofReportV1{}, err
	}
	if len(ms) != len(reportFields) {
		return ProofReportV1{}, refusal(RefusalMalformed, "report has %d members; want %d", len(ms), len(reportFields))
	}
	for i := range reportFields {
		if ms[i].name != reportFields[i] {
			return ProofReportV1{}, refusal(RefusalMalformed, "member %d is %q; want %q", i, ms[i].name, reportFields[i])
		}
	}
	var r ProofReportV1
	var subject, compiler string
	dsts := []any{&r.Schema, &subject, &compiler, &r.CompilerVersion, &r.Verified, &r.Errors, &r.Counterexamples, &r.CheckPassed, &r.ProofSucceeded}
	for i := range ms {
		if err := decodeExact(ms[i].raw, dsts[i]); err != nil {
			return ProofReportV1{}, err
		}
	}
	if r.Subject, err = hashref.Parse(subject); err != nil {
		return ProofReportV1{}, refusal(RefusalMalformed, "subject: %v", err)
	}
	if r.Compiler, err = hashref.Parse(compiler); err != nil {
		return ProofReportV1{}, refusal(RefusalMalformed, "compiler: %v", err)
	}
	canonical, err := EncodeProofReportV1(r)
	if err != nil {
		return ProofReportV1{}, err
	}
	if !bytes.Equal(raw, canonical) {
		return ProofReportV1{}, refusal(RefusalNonCanonical, "report bytes differ from canonical encoding")
	}
	return r, nil
}

func validateReport(r ProofReportV1) error {
	if r.Subject.IsZero() || r.Compiler.IsZero() {
		return refusal(RefusalMalformed, "subject and compiler must be canonical hashrefs")
	}
	strings := []struct{ label, value string }{{"schema", r.Schema}, {"subject", r.Subject.String()}, {"compiler", r.Compiler.String()}, {"compilerVersion", r.CompilerVersion}}
	for _, s := range strings {
		if err := checkString(s.label, s.value); err != nil {
			return err
		}
	}
	if len(r.Verified) > MaxVerifiedIdentities {
		return refusal(RefusalLimit, "verified has %d identities; limit is %d", len(r.Verified), MaxVerifiedIdentities)
	}
	for i, identity := range r.Verified {
		if err := checkString(fmt.Sprintf("verified[%d]", i), identity); err != nil {
			return err
		}
		// Strict pairwise increase gives sortedness AND uniqueness in one branch.
		// A following sort.StringsAreSorted belt was removed at iteration 105: by
		// transitivity it cannot fire on any input that reaches it, so it was an
		// undeclared unreachable branch — a guard nothing can pin, which this
		// milestone has already been bitten by once (see the AC19 arm).
		if i > 0 && r.Verified[i-1] >= identity {
			return refusal(RefusalMalformed, "verified identities are not sorted and unique")
		}
	}
	return nil
}

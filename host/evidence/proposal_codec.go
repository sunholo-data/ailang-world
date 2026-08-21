package evidence

import (
	"bytes"
	"encoding/json"
	"io"
	"unicode/utf8"
)

// DecodeProposal admits canonical JSON only and never confers authority.
func DecodeProposal(raw []byte) (ClaimedEvidence, error) {
	// This guard deliberately precedes UTF-8 validation and JSON parsing.
	if len(raw) > MaxBytes {
		return ClaimedEvidence{}, refusal(RefusalOversize, "proposal is %d bytes; limit is %d", len(raw), MaxBytes)
	}
	if !utf8.Valid(raw) {
		return ClaimedEvidence{}, refusal(RefusalInvalidUTF8, "input is not valid UTF-8")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	var v any
	err := d.Decode(&v)
	if err != nil {
		return ClaimedEvidence{}, refusal(RefusalMalformed, "%v", err)
	}
	if _, err := d.Token(); err != io.EOF {
		return ClaimedEvidence{}, refusal(RefusalMalformed, "trailing JSON")
	}
	canonical, err := json.Marshal(v)
	if err != nil {
		return ClaimedEvidence{}, refusal(RefusalMalformed, "%v", err)
	}
	if !bytes.Equal(raw, canonical) {
		return ClaimedEvidence{}, refusal(RefusalNonCanonical, "proposal bytes differ from canonical encoding")
	}
	return ClaimedEvidence{canonical: string(canonical)}, nil
}

package evidence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
)

var envelopeFields = [...]string{"report", "mac"}

// EncodeAuthenticatedEnvelope returns the canonical two-member envelope.
func EncodeAuthenticatedEnvelope(e AuthenticatedEnvelope) ([]byte, error) {
	if _, err := DecodeProofReportV1(e.Report); err != nil {
		return nil, err
	}
	if !e.MACValid || len(e.MAC) != 32 {
		return nil, refusal(RefusalMalformed, "mac is not a valid 32-byte tag")
	}
	type wire struct {
		Report string `json:"report"`
		MAC    string `json:"mac"`
	}
	return json.Marshal(wire{base64.RawURLEncoding.EncodeToString(e.Report), base64.RawURLEncoding.EncodeToString(e.MAC)})
}

// DecodeAuthenticatedEnvelope strictly decodes the envelope. A missing,
// malformed, padded, or wrong-width mac is retained as MACValid=false rather
// than refused, reserving all tag failures for authentication.
// The caller/store owns the raw-envelope byte cap; it is intentionally absent.
func DecodeAuthenticatedEnvelope(raw []byte) (AuthenticatedEnvelope, error) {
	ms, err := objectMembers(raw)
	if err != nil {
		return AuthenticatedEnvelope{}, err
	}
	if len(ms) < 1 || len(ms) > 2 || ms[0].name != envelopeFields[0] || (len(ms) == 2 && ms[1].name != envelopeFields[1]) {
		return AuthenticatedEnvelope{}, refusal(RefusalMalformed, "envelope members must be report then mac")
	}
	var reportText string
	if err := decodeExact(ms[0].raw, &reportText); err != nil {
		return AuthenticatedEnvelope{}, err
	}
	report, err := base64.RawURLEncoding.Strict().DecodeString(reportText)
	if err != nil || base64.RawURLEncoding.EncodeToString(report) != reportText {
		return AuthenticatedEnvelope{}, refusal(RefusalMalformed, "report is not canonical base64url without padding")
	}
	if len(report) > MaxBytes {
		return AuthenticatedEnvelope{}, refusal(RefusalOversize, "decoded report is %d bytes; limit is %d", len(report), MaxBytes)
	}
	decoded, err := DecodeProofReportV1(report)
	if err != nil {
		return AuthenticatedEnvelope{}, err
	}
	e := AuthenticatedEnvelope{Report: report, decoded: decoded}
	canonical := append([]byte(`{"report":`), mustJSON(reportText)...)
	if len(ms) == 2 {
		var macText string
		if decodeExact(ms[1].raw, &macText) == nil {
			e.MAC, err = base64.RawURLEncoding.Strict().DecodeString(macText)
			e.MACValid = err == nil && len(e.MAC) == 32 && base64.RawURLEncoding.EncodeToString(e.MAC) == macText
			if !e.MACValid {
				e.MAC = nil
			}
		}
		var macValue any
		if err := json.Unmarshal(ms[1].raw, &macValue); err != nil {
			return AuthenticatedEnvelope{}, refusal(RefusalMalformed, "malformed mac member")
		}
		canonical = append(canonical, `,"mac":`...)
		canonical = append(canonical, mustJSON(macValue)...)
	}
	canonical = append(canonical, '}')
	if !bytes.Equal(raw, canonical) {
		return AuthenticatedEnvelope{}, refusal(RefusalNonCanonical, "envelope bytes differ from canonical encoding")
	}
	return e, nil
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err) // strings and decoder-produced JSON values are always marshalable
	}
	return b
}

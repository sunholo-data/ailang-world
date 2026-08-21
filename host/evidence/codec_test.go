package evidence

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func goodReport() ProofReportV1 {
	return ProofReportV1{
		Schema: ProofReportSchemaV1, Subject: hashref.SumSHA256([]byte("subject")),
		Compiler: hashref.SumSHA256([]byte("compiler")), CompilerVersion: "AILANG v0.30.0",
		Verified: []string{"world/example.check"}, Errors: 0, Counterexamples: 0,
		CheckPassed: true, ProofSucceeded: true,
	}
}

func requireRefusal(t *testing.T, err error, kind string) {
	t.Helper()
	var got *DecodeRefusal
	if !errors.As(err, &got) || got.Kind != kind {
		t.Fatalf("error = %T %v; want *DecodeRefusal kind %q", err, err, kind)
	}
}

func TestProofReportCanonicalRoundTrip(t *testing.T) {
	raw, err := EncodeProofReportV1(goodReport())
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeProofReportV1(raw)
	if err != nil {
		t.Fatal(err)
	}
	again, err := EncodeProofReportV1(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(raw, again) {
		t.Fatalf("decode/re-encode changed bytes:\n%s\n%s", raw, again)
	}
	wantPrefix := `{"schema":"world/proof-report/v1","subject":`
	if !strings.HasPrefix(string(raw), wantPrefix) {
		t.Fatalf("field order = %s; want prefix %s", raw, wantPrefix)
	}
}

func TestProofReportStrictRefusals(t *testing.T) {
	good, _ := EncodeProofReportV1(goodReport())
	cases := []struct {
		name string
		raw  []byte
		kind string
	}{
		{"unknown", bytes.Replace(good, []byte(`"proofSucceeded":true`), []byte(`"unknown":true,"proofSucceeded":true`), 1), RefusalMalformed},
		{"duplicate", bytes.Replace(good, []byte(`"schema":`), []byte(`"schema":"world/proof-report/v1","schema":`), 1), RefusalMalformed},
		{"missing", bytes.Replace(good, []byte(`,"errors":0`), nil, 1), RefusalMalformed},
		{"wrong_order", bytes.Replace(good, []byte(`"schema":"world/proof-report/v1","subject":`), []byte(`"subject":"x","schema":"world/proof-report/v1","discard":`), 1), RefusalMalformed},
		{"whitespace_noncanonical", append([]byte(" "), good...), RefusalNonCanonical},
		{"trailing", append(good, []byte("null")...), RefusalMalformed},
		{"invalid_utf8", []byte{'{', 0xff, '}'}, RefusalInvalidUTF8},
		{"over_limit", bytes.Repeat([]byte{'{'}, MaxBytes+1), RefusalOversize},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeProofReportV1(tc.raw)
			requireRefusal(t, err, tc.kind)
		})
	}
	// Anti-vacuity: each arm requires the typed identity, so a generic parser
	// error, nil error, or zero report cannot produce the same observable.
}

// TestTruncatedTailReportIsRefusedNotPanicked pins the report ARITY guard
// (report_codec.go's `len(ms) != len(reportFields)`), which the iteration-105
// evaluator found unpinned: neutered with `if false && …` the mutant BUILDS and
// the whole suite stayed rc=0. The existing missing-member arm removes a MIDDLE
// field, so the per-index name comparison mismatches at i=5 and the arity guard
// is never the thing that fires — the gap is invisible precisely because a
// sibling guard covers the common case. Omitting the FINAL field instead leaves
// every remaining member correctly named AND ordered, so nothing mismatches and
// the loop indexes past the end. §3.3's contract is that malformed untrusted
// input becomes a typed refusal, never a panic; this arm is what makes that
// claim falsifiable.
func TestTruncatedTailReportIsRefusedNotPanicked(t *testing.T) {
	raw, err := EncodeProofReportV1(goodReport())
	if err != nil {
		t.Fatal(err)
	}
	const tail = `,"proofSucceeded":true}`
	if !strings.HasSuffix(string(raw), tail) {
		t.Fatalf("instrument failure: canonical report does not end %s: %s", tail, raw)
	}
	truncated := []byte(strings.TrimSuffix(string(raw), tail) + "}")

	// Control, in the SAME test: the untruncated bytes decode, so a
	// refuse-everything implementation cannot pass this arm vacuously.
	if _, err := DecodeProofReportV1(raw); err != nil {
		t.Fatalf("control: canonical report was refused: %v", err)
	}

	got, err := DecodeProofReportV1(truncated)
	if err == nil {
		t.Fatalf("tail-truncated report decoded as %+v; want a typed refusal", got)
	}
	requireRefusal(t, err, RefusalMalformed)
	// The observable is unique to the ARITY guard: a name/order mismatch reports
	// "member N is ...", and only this branch reports a member COUNT.
	var refused *DecodeRefusal
	if !errors.As(err, &refused) || !strings.Contains(refused.Detail, "report has 8 members; want 9") {
		t.Fatalf("refusal = %v; want the arity guard's own member-count refusal", err)
	}
}

func TestProofReportCaps(t *testing.T) {
	r := goodReport()
	r.CompilerVersion = strings.Repeat("x", MaxStringBytes+1)
	_, err := EncodeProofReportV1(r)
	requireRefusal(t, err, RefusalLimit)
	r = goodReport()
	r.Verified = make([]string, MaxVerifiedIdentities+1)
	for i := range r.Verified {
		r.Verified[i] = strings.Repeat("a", i/26+1) + string(rune('a'+i%26))
	}
	_, err = EncodeProofReportV1(r)
	requireRefusal(t, err, RefusalLimit)
	r = goodReport()
	r.Verified = []string{"same", "same"}
	_, err = EncodeProofReportV1(r)
	requireRefusal(t, err, RefusalMalformed)
}

func TestEnvelopeCanonicalRoundTripAndMACDeferral(t *testing.T) {
	report, _ := EncodeProofReportV1(goodReport())
	e := AuthenticatedEnvelope{Report: report, MAC: bytes.Repeat([]byte{7}, 32), MACValid: true}
	raw, err := EncodeAuthenticatedEnvelope(e)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeAuthenticatedEnvelope(raw)
	if err != nil {
		t.Fatal(err)
	}
	again, err := EncodeAuthenticatedEnvelope(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(raw, again) {
		t.Fatalf("decode/re-encode changed envelope bytes")
	}

	report64 := base64.RawURLEncoding.EncodeToString(report)
	for _, raw := range [][]byte{
		[]byte(`{"report":"` + report64 + `"}`),
		[]byte(`{"report":"` + report64 + `","mac":"%%%"}`),
		[]byte(`{"report":"` + report64 + `","mac":"AA"}`),
	} {
		got, err := DecodeAuthenticatedEnvelope(raw)
		if err != nil || got.MACValid {
			t.Fatalf("invalid mac = (%+v, %v); want admitted with MACValid=false", got, err)
		}
	}
	// Anti-vacuity: successful structural decode plus MACValid=false is unique
	// to deferred authentication; an ordinary decode refusal cannot satisfy it.
}

func TestEnvelopeStrictRefusals(t *testing.T) {
	report, _ := EncodeProofReportV1(goodReport())
	r64 := base64.RawURLEncoding.EncodeToString(report)
	mac := base64.RawURLEncoding.EncodeToString(bytes.Repeat([]byte{1}, 32))
	cases := []struct{ name, raw, kind string }{
		{"unknown", `{"report":"` + r64 + `","unknown":0,"mac":"` + mac + `"}`, RefusalMalformed},
		{"duplicate", `{"report":"` + r64 + `","report":"` + r64 + `","mac":"` + mac + `"}`, RefusalMalformed},
		{"missing_report", `{"mac":"` + mac + `"}`, RefusalMalformed},
		{"wrong_order", `{"mac":"` + mac + `","report":"` + r64 + `"}`, RefusalMalformed},
		{"padded_report", `{"report":"` + r64 + `=","mac":"` + mac + `"}`, RefusalMalformed},
		{"trailing", `{"report":"` + r64 + `","mac":"` + mac + `"}null`, RefusalMalformed},
		{"noncanonical", ` {"report":"` + r64 + `","mac":"` + mac + `"}`, RefusalNonCanonical},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeAuthenticatedEnvelope([]byte(tc.raw))
			requireRefusal(t, err, tc.kind)
		})
	}
	_, err := DecodeAuthenticatedEnvelope([]byte{'{', 0xff, '}'})
	requireRefusal(t, err, RefusalInvalidUTF8)
	// Anti-vacuity: exact typed kinds distinguish codec guards from a zero envelope.
}

func TestDecodeProposalCapsBeforeParse(t *testing.T) {
	raw := bytes.Repeat([]byte{'['}, MaxBytes+1) // also malformed JSON
	claim, err := DecodeProposal(raw)
	if !claim.IsZero() {
		t.Fatal("oversize proposal returned a claim")
	}
	requireRefusal(t, err, RefusalOversize)
	// Anti-vacuity: malformed parsing could also error, but cannot produce the
	// oversize refusal identity asserted above.
}

func TestProposalStrictRefusals(t *testing.T) {
	cases := []struct {
		name string
		raw  []byte
		kind string
	}{
		{"trailing", []byte(`[]null`), RefusalMalformed},
		{"noncanonical", []byte(`[ ]`), RefusalNonCanonical},
		{"duplicate", []byte(`{"x":1,"x":2}`), RefusalNonCanonical},
		{"invalid_utf8", []byte{'"', 0xff, '"'}, RefusalInvalidUTF8},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			claim, err := DecodeProposal(tc.raw)
			if !claim.IsZero() {
				t.Fatal("refusal returned a claim")
			}
			requireRefusal(t, err, tc.kind)
		})
	}
	// Anti-vacuity: both no-claim and the mechanism-specific refusal kind are required.
}

func TestNestingDepthBombWithinByteCapIsRefused(t *testing.T) {
	shallow := []byte(strings.Repeat("[", 10) + strings.Repeat("]", 10))
	claim, err := DecodeProposal(shallow)
	if err != nil || claim.IsZero() {
		t.Fatalf("depth-10 control = zero %v, err %v; want decoded claim", claim.IsZero(), err)
	}

	bomb := []byte(strings.Repeat("[", 131071) + strings.Repeat("]", 131071))
	if len(bomb) != 262142 || len(bomb) >= MaxBytes {
		t.Fatalf("depth-bomb instrument length = %d", len(bomb))
	}
	claim, err = DecodeProposal(bomb)
	if err == nil || !claim.IsZero() {
		t.Fatalf("depth bomb decoded as ClaimedEvidence; want typed decode refusal")
	}
	requireRefusal(t, err, RefusalMalformed)
	// The refusal must come from the DEPTH mechanism, not merely be "some typed
	// refusal". Measured (iteration 105): under M27 with the natural spelling of
	// DecodeProposal, the swallowed decode error falls through to the trailing-JSON
	// guard, which also yields a typed RefusalMalformed with a zero claim — so a
	// Kind-only assertion is satisfied by a bystander guard and M27 SURVIVES. The
	// scanner's own "exceeded max depth" text is produced by nothing else on this
	// path, which is what makes the mutant die for the reason AC19 names.
	var refused *DecodeRefusal
	if !errors.As(err, &refused) || !strings.Contains(refused.Detail, "exceeded max depth") {
		t.Fatalf("depth-bomb refusal = %v; want the stdlib scanner's own max-depth refusal", err)
	}
	// Anti-vacuity: a zero default alone could mimic no claim, but cannot carry
	// DecodeRefusal/malformed; the successful shallow arm defeats refuse-all.
}

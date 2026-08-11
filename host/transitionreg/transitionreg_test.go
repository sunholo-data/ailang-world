package transitionreg

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

var (
	testFn          = hashref.MustParse("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	testInterpreter = hashref.MustParse("sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
)

func validDescriptor() Descriptor {
	return Descriptor{ID: "tools.echo", TransitionFn: testFn, Interpreter: testInterpreter, SemanticsEpoch: 1, InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`), Access: EffectRequirement{Effect: "read", Scope: "world", Cost: 1}, DeclaredEffects: []EffectRequirement{{Effect: "read", Scope: "world", Cost: 1}}, Title: "Echo", Description: "test"}
}
func validRevision(entries ...Descriptor) Revision {
	return Revision{SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1, Revision: 1, Entries: entries}
}

func TestCodecGoldenRoundTrip(t *testing.T) {
	const wantInterface = "sha256:743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac"
	if InterfaceHashV1.String() != wantInterface {
		t.Fatalf("interface hash = %q, want literal %q", InterfaceHashV1, wantInterface)
	}

	const emptyGolden = `{"entries":[],"interfaceHash":"sha256:743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac","parent":"","revision":1,"semanticID":"world/transition-registry/v1"}`
	empty, err := EncodeRevision(validRevision())
	if err != nil {
		t.Fatal(err)
	}
	if string(empty) != emptyGolden {
		t.Fatalf("empty golden bytes:\n got %s\nwant %s", empty, emptyGolden)
	}

	d := validDescriptor()
	d.InputSchema, _ = canonicalSchema([]byte(`{"z":1.0,"a":"<&é"}`))
	d.OutputSchema, _ = canonicalSchema([]byte(`{"n":-0,"array":[3,2,1]}`))
	const descriptorGolden = `{"entries":[{"access":{"cost":1,"effect":"read","scope":"world"},"declaredEffects":[{"cost":1,"effect":"read","scope":"world"}],"description":"test","id":"tools.echo","inputSchema":{"a":"<&é","z":1},"interpreter":"sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","outputSchema":{"array":[3,2,1],"n":0},"semanticsEpoch":1,"title":"Echo","transitionFn":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}],"interfaceHash":"sha256:743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac","parent":"","revision":1,"semanticID":"world/transition-registry/v1"}`
	got, err := EncodeRevision(validRevision(d))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != descriptorGolden {
		t.Fatalf("descriptor golden bytes:\n got %s\nwant %s", got, descriptorGolden)
	}
	round, err := DecodeRevision(got)
	if err != nil {
		t.Fatal(err)
	}
	again, err := EncodeRevision(round)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, again) {
		t.Fatal("Encode -> Decode -> Encode changed bytes")
	}

	numbers, err := canonicalSchema([]byte(`{"a":1,"b":1.0,"c":1e0,"d":-0}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(numbers) != `{"a":1,"b":1,"c":1,"d":0}` {
		t.Fatalf("number normalization = %s", numbers)
	}
	nfc, _ := canonicalSchema([]byte(`{"s":"é"}`))
	nfd, _ := canonicalSchema([]byte("{\"s\":\"e\u0301\"}"))
	if bytes.Equal(nfc, nfd) {
		t.Fatal("NFC and NFD strings were normalized together")
	}
}

func TestDescriptorIdentityAndContentUpdate(t *testing.T) {
	d1 := validDescriptor()
	r1 := validRevision(d1)
	b1, err := EncodeRevision(r1)
	if err != nil {
		t.Fatal(err)
	}
	h1 := hashref.SumSHA256(b1)
	d2 := d1
	d2.TransitionFn = hashref.SumSHA256([]byte("updated transition source"))
	r2 := validRevision(d2)
	r2.Revision = 2
	r2.Parent = h1
	b2, err := EncodeRevision(r2)
	if err != nil {
		t.Fatal(err)
	}
	h2 := hashref.SumSHA256(b2)
	if d1.ID != d2.ID || h1 == h2 {
		t.Fatalf("stable ID/content update invariant failed: ids=%q/%q hashes=%q/%q", d1.ID, d2.ID, h1, h2)
	}
	objects := map[hashref.HashRef][]byte{h1: append([]byte(nil), b1...), h2: append([]byte(nil), b2...)}
	if !bytes.Equal(objects[h1], b1) {
		t.Fatal("old immutable revision is no longer addressable")
	}
	if bytes.Equal(objects[h1], objects[h2]) {
		t.Fatal("content update did not create a new revision object")
	}
}

func TestDescriptorValidationRefusals(t *testing.T) {
	tests := []struct {
		name string
		want string
		run  func() error
	}{
		{"id_grammar", "does not match stable ID grammar", func() error { d := validDescriptor(); d.ID = "Bad ID"; return d.Validate() }},
		{"id_too_long", "length must be 1..128 bytes", func() error {
			d := validDescriptor()
			d.ID = strings.Join([]string{strings.Repeat("a", 32), strings.Repeat("b", 32), strings.Repeat("c", 32), strings.Repeat("d", 32)}, "/")
			return d.Validate()
		}},
		{"segment_too_long", "segment length must be 1..32 bytes", func() error { d := validDescriptor(); d.ID = strings.Repeat("a", 33); return d.Validate() }},
		{"zero_transition_fn", "transitionFn is zero", func() error { d := validDescriptor(); d.TransitionFn = hashref.HashRef{}; return d.Validate() }},
		{"zero_interpreter", "interpreter is zero", func() error { d := validDescriptor(); d.Interpreter = hashref.HashRef{}; return d.Validate() }},
		{"negative_semantics_epoch", "semanticsEpoch is negative", func() error { d := validDescriptor(); d.SemanticsEpoch = -1; return d.Validate() }},
		{"negative_cost", "cost is negative", func() error { d := validDescriptor(); d.Access.Cost = -1; return d.Validate() }},
		{"schema_not_an_object", "schema root must be an object", func() error { d := validDescriptor(); d.InputSchema = []byte(`[]`); return d.Validate() }},
		{"schema_raw_over_262144", "raw JSON is 262146 bytes; limit is 262144", func() error {
			raw := append(bytes.Repeat([]byte(" "), maxSchemaRaw), []byte(`{}`)...)
			_, err := canonicalSchema(raw)
			return err
		}},
		{"schema_canonical_over_65536", "canonical schema is 65544 bytes; limit is 65536", func() error {
			d := validDescriptor()
			d.InputSchema = []byte(`{"x":"` + strings.Repeat("a", maxSchemaCanonical) + `"}`)
			return d.Validate()
		}},
		{"duplicate_schema_key_nested", "duplicate object member \"a\"", func() error {
			d := validDescriptor()
			d.InputSchema = []byte(`{"x":{"a":1,"a":2}}`)
			return d.Validate()
		}},
		{"lone_surrogate", "JSON string contains an escaped surrogate", func() error { d := validDescriptor(); d.InputSchema = []byte(`{"x":"\ud800"}`); return d.Validate() }},
		{"invalid_utf8", "JSON is not valid UTF-8", func() error {
			d := validDescriptor()
			d.InputSchema = []byte{'{', '"', 'x', '"', ':', '"', 0xff, '"', '}'}
			if utf8.Valid(d.InputSchema) {
				t.Fatal("fixture unexpectedly valid")
			}
			return d.Validate()
		}},
		{"number_coefficient_overflow", "number coefficient has 1025 digits; limit is 1024", func() error {
			d := validDescriptor()
			d.InputSchema = []byte(`{"n":` + strings.Repeat("1", maxNumberDigits+1) + `}`)
			return d.Validate()
		}},
		{"unknown_revision_key", "unknown key \"extra\"", func() error {
			_, err := DecodeRevision([]byte(`{"entries":[],"extra":0,"interfaceHash":"` + InterfaceHashV1.String() + `","parent":"","revision":1,"semanticID":"` + SemanticIDV1 + `"}`))
			return err
		}},
		{"missing_descriptor_key", "missing key \"id\"", func() error {
			_, err := DecodeRevision([]byte(`{"entries":[{}],"interfaceHash":"` + InterfaceHashV1.String() + `","parent":"","revision":1,"semanticID":"` + SemanticIDV1 + `"}`))
			return err
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			if err == nil {
				t.Fatal("invalid value was accepted")
			}
			// Pin the BRANCH, not merely the refusal: DecodeRevision's canonical
			// re-encode check (codec.go) refuses these inputs too, so a message-agnostic
			// assertion stays green when the named guard is neutered. Measured: with
			// the escaped-surrogate and unknown-key guards disabled, this test passed.
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("refused by the wrong branch: got %q, want it to contain %q", err, tc.want)
			}
		})
	}
}

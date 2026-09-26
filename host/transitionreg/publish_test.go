package transitionreg

import (
	"testing"
)

func TestCanonicalSchemaExportMatchesValidate(t *testing.T) {
	raw := []byte(`{"z":1.0,"a":"x"}`)
	got, err := CanonicalSchema(raw)
	if err != nil {
		t.Fatal(err)
	}
	d := validDescriptor()
	d.InputSchema = got
	if err := d.Validate(); err != nil {
		t.Fatalf("CanonicalSchema output rejected by Validate: %v", err)
	}
	if string(got) != `{"a":"x","z":1}` {
		t.Fatalf("CanonicalSchema = %s, want the codec's canonical form", got)
	}
}

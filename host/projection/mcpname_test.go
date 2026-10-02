package projection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"github.com/sunholo-data/ailang/serveapi/protocol"
	"strings"
	"testing"
)

func TestMCPNameRoundTrip(t *testing.T) {
	for id, want := range map[string]string{"tools.echo": "tools_decho", "a_b.c/d": "a_ub_dc_sd", "world/recovery-transition/v1": "world_srecovery-transition_sv1"} {
		got, err := EncodeMCPName(id)
		if err != nil || got != want {
			t.Fatalf("exact encoding %q=%q/%v want %q", id, got, err, want)
		}
		decoded, err := decodeMCPName(got)
		if err != nil || decoded != id {
			t.Fatalf("roundtrip %q: %q/%v", id, decoded, err)
		}
	}
	for i := 0; i < 200; i++ {
		id := fmt.Sprintf("tools/group_%d.action-%d", i, i)
		name, err := EncodeMCPName(id)
		if err != nil {
			t.Fatal(err)
		}
		back, err := decodeMCPName(name)
		if err != nil || back != id {
			t.Fatalf("corpus %q -> %q -> %q: %v", id, name, back, err)
		}
		again, err := EncodeMCPName(back)
		if err != nil || again != name {
			t.Fatal("noncanonical roundtrip")
		}
	}
	a, _ := EncodeMCPName("a_b")
	b, _ := EncodeMCPName("a.b")
	if a == b {
		t.Fatal("distinct IDs collide")
	}
}
func TestMCPNameRefusal(t *testing.T) {
	for _, name := range []string{"a_", "a_x", "UPPER", "a.b", "a/b", "", "a__u"} {
		if got, err := decodeMCPName(name); err == nil {
			t.Fatalf("invalid encoding %q accepted as %q", name, got)
		}
	}
	for _, id := range []string{strings.Repeat("a", 32) + "/" + strings.Repeat("b", 32), strings.Repeat("a_", 31) + ".b"} {
		if _, err := EncodeMCPName(id); err == nil {
			t.Fatalf("long encoding accepted: %q", id)
		}
	}
	ds := []transitionreg.Descriptor{descriptor("tools.echo", "alpha"), descriptor("tools.echo", "alpha")}
	if _, err := mcpDescriptors(ds); err == nil {
		t.Fatal("duplicate descriptor surface accepted")
	}
}
func TestMCPSchemaNormalization(t *testing.T) {
	typed := []byte(`{"type":"object","properties":{"x":{"type":"integer"}},"minimum":9007199254740993}`)
	for _, raw := range [][]byte{[]byte(`{}`), []byte(`{"properties":{"x":{"type":"integer"}},"minimum":9007199254740993}`), typed} {
		out, err := normalizeMCPSchema(raw)
		if err != nil {
			t.Fatal(err)
		}
		var before, after map[string]json.RawMessage
		if err := json.Unmarshal(raw, &before); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(out, &after); err != nil {
			t.Fatal(err)
		}
		if string(after["type"]) != `"object"` {
			t.Fatalf("missing object type: %s", out)
		}
		for k, v := range before {
			if !bytes.Equal(after[k], v) {
				t.Fatalf("constraint %s lost/rounded: %s want %s", k, after[k], v)
			}
		}
		if _, ok := before["type"]; ok && !bytes.Equal(raw, out) {
			t.Fatal("typed bytes changed")
		}
	}
	for _, raw := range []string{`{"type":"array"}`, `[]`, `null`, `bad`} {
		if _, err := normalizeMCPSchema([]byte(raw)); err == nil {
			t.Fatalf("invalid schema accepted: %s", raw)
		}
	}
	d := descriptor("tools.echo", "alpha")
	d.InputSchema = []byte(`{"properties":{"x":{"type":"string","x-mcp-header":17}}}`)
	ds, err := mcpDescriptors([]transitionreg.Descriptor{d})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := protocol.CallerSurface(ds); err == nil {
		t.Fatal("invalid annotation accepted")
	}
	d.InputSchema = []byte(`{}`)
	d.OutputSchema = []byte(`{"answer":9007199254740993}`)
	ds, err = mcpDescriptors([]transitionreg.Descriptor{d})
	if err != nil || !bytes.Equal(ds[0].OutputSchema, d.OutputSchema) {
		t.Fatalf("output schema changed %v", err)
	}
}
func ExampleEncodeMCPName() {
	name, _ := EncodeMCPName("tools.echo")
	fmt.Println(name)
	// Output: tools_decho
}

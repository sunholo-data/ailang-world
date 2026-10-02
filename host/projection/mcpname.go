package projection

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"github.com/sunholo-data/ailang/serveapi/protocol"
	"strings"
)

// EncodeMCPName reversibly maps a registry ID to an MCP name. Registry admission
// owns ID validity; names longer than the released protocol's limit refuse the
// whole surface rather than dropping, truncating or aliasing a transition.
func EncodeMCPName(id string) (string, error) {
	name := strings.NewReplacer("_", "_u", ".", "_d", "/", "_s").Replace(id)
	if len(name) > 64 {
		return "", fmt.Errorf("projection: encoded transition name exceeds 64 bytes: %q", id)
	}
	return name, nil
}
func decodeMCPName(name string) (string, error) {
	if name == "" || len(name) > 64 {
		return "", errors.New("projection: invalid MCP name")
	}
	var b strings.Builder
	for i := 0; i < len(name); i++ {
		c := name[i]
		if c == '_' {
			i++
			if i >= len(name) {
				return "", errors.New("projection: invalid MCP escape")
			}
			switch name[i] {
			case 'u':
				b.WriteByte('_')
			case 'd':
				b.WriteByte('.')
			case 's':
				b.WriteByte('/')
			default:
				return "", errors.New("projection: invalid MCP escape")
			}
			continue
		}
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return "", errors.New("projection: invalid MCP name")
		}
		b.WriteByte(c)
	}
	return b.String(), nil
}
func normalizeMCPSchema(raw []byte) (json.RawMessage, error) {
	var members map[string]json.RawMessage
	if err := json.Unmarshal(raw, &members); err != nil || members == nil {
		return nil, errors.New("projection: input schema must be a JSON object")
	}
	if typ, ok := members["type"]; ok {
		var name string
		if json.Unmarshal(typ, &name) != nil || name != "object" {
			return nil, errors.New("projection: explicit input schema type must be object")
		}
		return append(json.RawMessage(nil), raw...), nil
	}
	members["type"] = json.RawMessage(`"object"`)
	return json.Marshal(members)
}
func mcpDescriptors(ds []transitionreg.Descriptor) ([]protocol.ToolDescriptor, error) {
	tools := make([]protocol.ToolDescriptor, 0, len(ds))
	seen := make(map[string]bool, len(ds))
	for _, d := range ds {
		name, err := EncodeMCPName(d.ID)
		if err != nil {
			return nil, err
		}
		if seen[name] {
			return nil, fmt.Errorf("projection: duplicate MCP name for %q", d.ID)
		}
		seen[name] = true
		schema, err := normalizeMCPSchema(d.InputSchema)
		if err != nil {
			return nil, fmt.Errorf("projection: schema for %q: %w", d.ID, err)
		}
		tools = append(tools, protocol.ToolDescriptor{Name: name, Description: d.Description, InputSchema: schema, OutputSchema: append(json.RawMessage(nil), d.OutputSchema...)})
	}
	return tools, nil
}

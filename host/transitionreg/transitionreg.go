package transitionreg

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

type EffectRequirement struct {
	Effect, Scope string
	Cost          int64
}
type Descriptor struct {
	ID                        string
	TransitionFn, Interpreter hashref.HashRef
	SemanticsEpoch            int64
	InputSchema, OutputSchema []byte
	Access                    EffectRequirement
	DeclaredEffects           []EffectRequirement
	Title, Description        string
}
type Revision struct {
	SemanticID    string
	InterfaceHash hashref.HashRef
	Revision      int64
	Parent        hashref.HashRef
	Entries       []Descriptor
}

func (d Descriptor) Validate() error {
	if err := validateID(d.ID); err != nil {
		return fmt.Errorf("descriptor ID: %w", err)
	}
	if d.TransitionFn.IsZero() {
		return errors.New("descriptor transitionFn is zero")
	}
	if d.Interpreter.IsZero() {
		return errors.New("descriptor interpreter is zero")
	}
	if d.SemanticsEpoch < 0 {
		return errors.New("descriptor semanticsEpoch is negative")
	}
	in, err := canonicalSchema(d.InputSchema)
	if err != nil {
		return fmt.Errorf("input schema: %w", err)
	}
	if !bytes.Equal(in, d.InputSchema) {
		return errors.New("input schema is not canonical")
	}
	out, err := canonicalSchema(d.OutputSchema)
	if err != nil {
		return fmt.Errorf("output schema: %w", err)
	}
	if !bytes.Equal(out, d.OutputSchema) {
		return errors.New("output schema is not canonical")
	}
	if err := validateRequirement(d.Access); err != nil {
		return fmt.Errorf("access: %w", err)
	}
	seen := map[EffectRequirement]struct{}{}
	for _, e := range d.DeclaredEffects {
		if err := validateRequirement(e); err != nil {
			return fmt.Errorf("declared effect: %w", err)
		}
		if _, ok := seen[e]; ok {
			return errors.New("duplicate declared effect")
		}
		seen[e] = struct{}{}
	}
	return nil
}
func validateRequirement(e EffectRequirement) error {
	if e.Cost < 0 {
		return errors.New("cost is negative")
	}
	return nil
}
func validateID(id string) error {
	if len(id) < 1 || len(id) > 128 {
		return errors.New("length must be 1..128 bytes")
	}
	for _, seg := range strings.FieldsFunc(id, func(r rune) bool { return r == '.' || r == '/' }) {
		if len(seg) < 1 || len(seg) > 32 {
			return errors.New("segment length must be 1..32 bytes")
		}
	}
	if strings.Contains(id, "//") || strings.Contains(id, "..") || strings.Contains(id, "/.") || strings.Contains(id, "./") {
		return errors.New("empty segment")
	}
	for i, c := range []byte(id) {
		valid := c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '-' || c == '.' || c == '/'
		if !valid {
			return errors.New("does not match stable ID grammar")
		}
		if (i == 0 || i == len(id)-1 || id[i-1] == '.' || id[i-1] == '/') && (c == '_' || c == '-' || c == '.' || c == '/') {
			return errors.New("does not match stable ID grammar")
		}
		if i+1 < len(id) && (id[i+1] == '.' || id[i+1] == '/') && (c == '_' || c == '-') {
			return errors.New("does not match stable ID grammar")
		}
	}
	return nil
}
func (r Revision) Validate() error {
	if r.SemanticID != SemanticIDV1 {
		return fmt.Errorf("semanticID %q is not %q", r.SemanticID, SemanticIDV1)
	}
	if r.InterfaceHash != InterfaceHashV1 {
		return errors.New("wrong interface hash")
	}
	if r.Revision < 0 {
		return errors.New("revision is negative")
	}
	if len(r.Entries) > maxEntries {
		return fmt.Errorf("entries exceeds %d", maxEntries)
	}
	for i := range r.Entries {
		if err := r.Entries[i].Validate(); err != nil {
			return fmt.Errorf("entry %d: %w", i, err)
		}
		if i > 0 && CompareID(r.Entries[i-1].ID, r.Entries[i].ID) >= 0 {
			return errors.New("entries are not strictly ordered by ID")
		}
	}
	return nil
}
func CompareID(a, b string) int { return bytes.Compare([]byte(a), []byte(b)) }
func SortedDescriptors(entries []Descriptor) ([]Descriptor, error) {
	out := append([]Descriptor(nil), entries...)
	sort.Slice(out, func(i, j int) bool { return CompareID(out[i].ID, out[j].ID) < 0 })
	for i := range out {
		if err := out[i].Validate(); err != nil {
			return nil, err
		}
		if i > 0 && out[i-1].ID == out[i].ID {
			return nil, fmt.Errorf("duplicate descriptor ID %q", out[i].ID)
		}
	}
	return out, nil
}

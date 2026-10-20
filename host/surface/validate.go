package surface

import (
	"fmt"
	"time"
)

// Validate checks an index before anything shows it. Each refusal names the
// field and the state or brane it is about:
//
//   - the schema is Schema, built_at is RFC 3339, and there are at most
//     MaxStates states;
//   - brane keys are unique and non-empty, and brane kinds are attended,
//     loop or mixed;
//   - state ids are unique full commit hashes, kinds are shared or step,
//     and times are RFC 3339;
//   - every parent resolves to a state, except the oldest shared state's
//     first parent, which lies outside the window;
//   - a shared state has two parents: the previous shared state (never
//     later than it) and the step it merged; it has no brane;
//   - a step's brane is "" or listed in branes, it has a base, and its merge
//     is a shared state whose second parent is that step, never earlier
//     than the step.
func Validate(ix Index) error {
	if ix.Schema != Schema {
		return fmt.Errorf("surface index: schema %q, want %q", ix.Schema, Schema)
	}
	if _, err := time.Parse(time.RFC3339, ix.BuiltAt); err != nil {
		return fmt.Errorf("surface index: built_at %q is not RFC 3339", ix.BuiltAt)
	}
	if len(ix.States) > MaxStates {
		return fmt.Errorf("surface index: %d states, over the bound of %d", len(ix.States), MaxStates)
	}

	branes := map[string]bool{}
	for i, b := range ix.Branes {
		switch {
		case b.Key == "":
			return fmt.Errorf("surface index: branes[%d].key is empty", i)
		case branes[b.Key]:
			return fmt.Errorf("surface index: branes[%d].key %q is a duplicate", i, b.Key)
		case b.Kind != BraneAttended && b.Kind != BraneLoop && b.Kind != BraneMixed:
			return fmt.Errorf("surface index: branes[%d] %s: kind %q is not attended, loop or mixed", i, b.Key, b.Kind)
		}
		branes[b.Key] = true
	}

	byID := map[string]int{}
	times := make([]time.Time, len(ix.States))
	oldestShared := -1
	for i, s := range ix.States {
		if !commitHash.MatchString(s.ID) {
			return fmt.Errorf("surface index: states[%d].id %q is not a full commit hash", i, s.ID)
		}
		if j, dup := byID[s.ID]; dup {
			return fmt.Errorf("surface index: states[%d].id %s duplicates states[%d]", i, s.ID, j)
		}
		byID[s.ID] = i
		if s.Kind != KindShared && s.Kind != KindStep {
			return fmt.Errorf("surface index: states[%d] %s: kind %q is not shared or step", i, s.ID, s.Kind)
		}
		t, err := time.Parse(time.RFC3339, s.At)
		if err != nil {
			return fmt.Errorf("surface index: states[%d] %s: at %q is not RFC 3339", i, s.ID, s.At)
		}
		times[i] = t
		if s.Kind == KindShared && oldestShared < 0 {
			oldestShared = i
		}
	}

	for i, s := range ix.States {
		for k, p := range s.Parents {
			if _, ok := byID[p]; !ok && !(i == oldestShared && k == 0) {
				return fmt.Errorf("surface index: states[%d] %s: parents[%d] %q does not resolve", i, s.ID, k, p)
			}
		}
	}
	// The shared chain first, so a shared state out of order is named
	// itself rather than through the step it merged.
	for i, s := range ix.States {
		if s.Kind == KindShared {
			if err := validateShared(ix, byID, times, i, i == oldestShared); err != nil {
				return err
			}
		}
	}
	for i, s := range ix.States {
		if s.Kind == KindStep {
			if err := validateStep(ix, byID, branes, times, i); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateShared(ix Index, byID map[string]int, times []time.Time, i int, oldest bool) error {
	s := ix.States[i]
	if len(s.Parents) != 2 {
		return fmt.Errorf("surface index: states[%d] %s: shared state has %d parents, want 2", i, s.ID, len(s.Parents))
	}
	if s.Brane != "" {
		return fmt.Errorf("surface index: states[%d] %s: brane %q on a shared state", i, s.ID, s.Brane)
	}
	if j, ok := byID[s.Parents[0]]; ok {
		if ix.States[j].Kind != KindShared {
			return fmt.Errorf("surface index: states[%d] %s: parents[0] %s is a %s, not the previous shared state", i, s.ID, s.Parents[0], ix.States[j].Kind)
		}
		if times[i].Before(times[j]) {
			return fmt.Errorf("surface index: states[%d] %s: at %s is before its previous shared state %s (%s)", i, s.ID, s.At, s.Parents[0], ix.States[j].At)
		}
	} else if !oldest {
		return fmt.Errorf("surface index: states[%d] %s: parents[0] %q does not resolve", i, s.ID, s.Parents[0])
	}
	if j := byID[s.Parents[1]]; ix.States[j].Kind != KindStep {
		return fmt.Errorf("surface index: states[%d] %s: parents[1] %s is a %s, not the step it merged", i, s.ID, s.Parents[1], ix.States[j].Kind)
	}
	return nil
}

func validateStep(ix Index, byID map[string]int, branes map[string]bool, times []time.Time, i int) error {
	s := ix.States[i]
	if s.Brane != "" && !branes[s.Brane] {
		return fmt.Errorf("surface index: states[%d] %s: brane %q is not in branes", i, s.ID, s.Brane)
	}
	if !commitHash.MatchString(s.Base) {
		return fmt.Errorf("surface index: states[%d] %s: base %q is not a full commit hash", i, s.ID, s.Base)
	}
	j, ok := byID[s.Merge]
	if !ok {
		return fmt.Errorf("surface index: states[%d] %s: merge %q does not resolve", i, s.ID, s.Merge)
	}
	m := ix.States[j]
	if m.Kind != KindShared || len(m.Parents) != 2 || m.Parents[1] != s.ID {
		return fmt.Errorf("surface index: states[%d] %s: merge %s is not a shared state whose second parent is this step", i, s.ID, s.Merge)
	}
	if times[i].After(times[j]) {
		return fmt.Errorf("surface index: states[%d] %s: at %s is after its merge %s (%s)", i, s.ID, s.At, s.Merge, m.At)
	}
	return nil
}

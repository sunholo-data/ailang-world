package surface

import (
	"fmt"
	"path/filepath"
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
//     than the step;
//   - the episode facts hold (validateRig).
func Validate(ix Index) error {
	if err := validateHuman(ix); err != nil {
		return err
	}
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
	return validateRig(ix, byID)
}

// validateRig checks the episode facts (M0c):
//
//   - they are present only with source.rig, which names the reads method
//     and counts the steps it attached;
//   - each steps[] entry has an episode, unique; its state is a step state,
//     used once, whose pr matches; rounds ≥ 1 and 0 ≤ passed ≤ total;
//   - every files, reads and decisions[].step key is a step in steps;
//   - status is A, M or D; via is new, step, commit or outside, with from
//     "", a step state, a commit or "" to match; actor is world, you or
//     loop; kind is grade, refusal or merge; at is RFC 3339;
//   - every path is filepath.IsLocal.
func validateRig(ix Index, byID map[string]int) error {
	rig := ix.Source.Rig
	if rig == nil {
		if len(ix.Steps)+len(ix.Files)+len(ix.Reads)+len(ix.Decisions) > 0 || ix.Source.ReadsMethod != "" {
			return fmt.Errorf("surface index: steps, files, reads, decisions and reads_method need source.rig")
		}
		return nil
	}
	if ix.Source.ReadsMethod != ReadsMethod {
		return fmt.Errorf("surface index: source.reads_method %q, want %q", ix.Source.ReadsMethod, ReadsMethod)
	}
	if rig.Attached != len(ix.Steps) {
		return fmt.Errorf("surface index: source.rig.attached %d, but %d steps", rig.Attached, len(ix.Steps))
	}
	steps := map[string]bool{}
	episodes := map[string]bool{}
	for i, s := range ix.Steps {
		j, ok := byID[s.State]
		switch {
		case s.Episode == "" || episodes[s.Episode]:
			return fmt.Errorf("surface index: steps[%d].episode %q is empty or a duplicate", i, s.Episode)
		case !ok || ix.States[j].Kind != KindStep:
			return fmt.Errorf("surface index: steps[%d] %s: state %q is not a step state", i, s.Episode, s.State)
		case steps[s.State]:
			return fmt.Errorf("surface index: steps[%d] %s: state %s is used twice", i, s.Episode, s.State)
		case s.PR != ix.States[j].PR:
			return fmt.Errorf("surface index: steps[%d] %s: pr %d, but state %s has pr %d", i, s.Episode, s.PR, s.State, ix.States[j].PR)
		case s.Rounds < 1:
			return fmt.Errorf("surface index: steps[%d] %s: rounds %d, want at least 1", i, s.Episode, s.Rounds)
		case s.Graders.Passed < 0 || s.Graders.Passed > s.Graders.Total:
			return fmt.Errorf("surface index: steps[%d] %s: graders %d/%d, want 0 ≤ passed ≤ total", i, s.Episode, s.Graders.Passed, s.Graders.Total)
		}
		steps[s.State], episodes[s.Episode] = true, true
	}
	for id, fc := range ix.Files {
		if !steps[id] {
			return fmt.Errorf("surface index: files[%s] is not a step in steps", id)
		}
		for _, c := range fc.Changes {
			if c.Status != "A" && c.Status != "M" && c.Status != "D" {
				return fmt.Errorf("surface index: files[%s] %s: status %q is not A, M or D", id, c.Path, c.Status)
			}
			if !filepath.IsLocal(c.Path) {
				return fmt.Errorf("surface index: files[%s]: path %q is not local", id, c.Path)
			}
		}
		if len(fc.Changes) > MaxChanges || fc.Truncated < 0 {
			return fmt.Errorf("surface index: files[%s]: %d changes listed, %d truncated", id, len(fc.Changes), fc.Truncated)
		}
	}
	for id, reads := range ix.Reads {
		if !steps[id] {
			return fmt.Errorf("surface index: reads[%s] is not a step in steps", id)
		}
		for _, r := range reads {
			if !filepath.IsLocal(r.Path) {
				return fmt.Errorf("surface index: reads[%s]: path %q is not local", id, r.Path)
			}
			if err := validateVia(ix, byID, r); err != nil {
				return fmt.Errorf("surface index: reads[%s] %s: %v", id, r.Path, err)
			}
		}
	}
	for i, d := range ix.Decisions {
		switch {
		case !steps[d.Step]:
			return fmt.Errorf("surface index: decisions[%d].step %q is not a step in steps", i, d.Step)
		case d.Actor != ActorWorld && d.Actor != ActorYou && d.Actor != ActorLoop:
			return fmt.Errorf("surface index: decisions[%d]: actor %q is not world, you or loop", i, d.Actor)
		case d.Kind != DecisionGrade && d.Kind != DecisionRefusal && d.Kind != DecisionMerge:
			return fmt.Errorf("surface index: decisions[%d]: kind %q is not grade, refusal or merge", i, d.Kind)
		}
		if _, err := time.Parse(time.RFC3339, d.At); err != nil {
			return fmt.Errorf("surface index: decisions[%d]: at %q is not RFC 3339", i, d.At)
		}
	}
	return nil
}

func validateVia(ix Index, byID map[string]int, r Read) error {
	switch r.Via {
	case ViaNew, ViaOutside:
		if r.From != "" {
			return fmt.Errorf("via %s has from %q, want \"\"", r.Via, r.From)
		}
	case ViaCommit:
		if !commitHash.MatchString(r.From) {
			return fmt.Errorf("via commit: from %q is not a full commit hash", r.From)
		}
	case ViaStep:
		if j, ok := byID[r.From]; !ok || ix.States[j].Kind != KindStep {
			return fmt.Errorf("via step: from %q is not a step state", r.From)
		}
	default:
		return fmt.Errorf("via %q is not new, step, commit or outside", r.Via)
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

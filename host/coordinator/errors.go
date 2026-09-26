package coordinator

import "fmt"

// InvalidCallError (R1) reports a call that fails validation before any
// registry, store or capsule access.
type InvalidCallError struct{ Field string }

func (e *InvalidCallError) Error() string {
	return fmt.Sprintf("coordinator: invalid call: %s", e.Field)
}

// SourceError (R6) reports a descriptor whose transition source object is
// absent or does not hash to the descriptor's pin.
type SourceError struct {
	Kind string // "absent" | "corrupt"
}

func (e *SourceError) Error() string { return "coordinator: transition source " + e.Kind }

// WorldAbsentError (R7) reports a store with no selected world head.
type WorldAbsentError struct{}

func (e *WorldAbsentError) Error() string { return "coordinator: no world is selected" }

// EffectsUnsupportedError (R8) reports a descriptor that declares effects:
// the capsule runs with no capabilities, so no effect could be honoured.
type EffectsUnsupportedError struct{ ID string }

func (e *EffectsUnsupportedError) Error() string {
	return fmt.Sprintf("coordinator: transition %q declares effects; effectful invocation is unavailable", e.ID)
}

// IncompatibleError (R9) reports a source that the pinned interpreter runs
// but whose entry does not implement main(input: string) -> string.
type IncompatibleError struct{ Err error }

func (e *IncompatibleError) Error() string {
	return fmt.Sprintf("coordinator: transition does not implement the calling convention: %v", e.Err)
}
func (e *IncompatibleError) Unwrap() error { return e.Err }

// ExecutionError (R10) reports any other capsule failure.
type ExecutionError struct{ Err error }

func (e *ExecutionError) Error() string {
	return fmt.Sprintf("coordinator: transition execution failed: %v", e.Err)
}
func (e *ExecutionError) Unwrap() error { return e.Err }

// OutputError (R12) reports output that is oversized or not a JSON object.
type OutputError struct{ Reason string }

func (e *OutputError) Error() string { return "coordinator: transition output " + e.Reason }

// NotCommittedError (R15) reports a resent task id whose durable intent has no
// outcome. The store is single-connection and single-writer, so by the time
// the receipt read runs no commit for it is in flight: the intent's commit
// did not land. The task id is spent; nothing is re-executed.
type NotCommittedError struct{ InvocationID string }

func (e *NotCommittedError) Error() string {
	return fmt.Sprintf("coordinator: invocation %s was not committed; send a new task id", e.InvocationID)
}

// UnconfirmedError (R16) reports a Commit that returned an untyped error: the
// outcome is not confirmed to this caller. Resending the same task id
// reconciles it from the journal (resolved → the committed result; intent
// only → R15).
type UnconfirmedError struct {
	InvocationID string
	Err          error
}

func (e *UnconfirmedError) Error() string {
	return fmt.Sprintf("coordinator: invocation %s outcome not confirmed: %v", e.InvocationID, e.Err)
}
func (e *UnconfirmedError) Unwrap() error { return e.Err }

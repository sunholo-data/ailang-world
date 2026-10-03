package coordinator

import (
	"fmt"
	"strings"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

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

// EffectsUnsupportedError (R8) reports a descriptor that declares an effect
// the session's handler registry cannot execute. It is decided before any
// effect runs (row 134): a descriptor that could be only partly honoured runs
// none of its effects. Effects names the unhandled declared effects.
type EffectsUnsupportedError struct {
	ID      string
	Effects []string
}

func (e *EffectsUnsupportedError) Error() string {
	return fmt.Sprintf("coordinator: transition %q declares effects with no registered handler %v; effectful invocation is unavailable", e.ID, e.Effects)
}

// EffectsUnrecordedError (row 134 §4.1, §4.3) reports an invocation whose
// effect was requested through the broker — so its intent, outcome and effect
// record are already durable in the effect journal and its budget is debited
// exactly once — but whose invocation was not committed: an R14 head conflict,
// a finish-phase failure, an output that cannot be composed, or a failed
// post-effect write or commit. Cause carries that failure (it may itself be
// an *UnconfirmedError, when the commit's outcome is unknown). EffectRecords
// are the effect-record refs, in plan order, resolvable through the store's
// effect receipts. The coordinator never retries: re-planning would mint a
// different invocation intent for the same ID (V46).
type EffectsUnrecordedError struct {
	InvocationID  string
	EffectRecords []hashref.HashRef
	Cause         error
}

func (e *EffectsUnrecordedError) Error() string {
	refs := make([]string, len(e.EffectRecords))
	for i, r := range e.EffectRecords {
		refs[i] = r.String()
	}
	return fmt.Sprintf("coordinator: invocation %s requested effects (records [%s]) but was not committed: %v",
		e.InvocationID, strings.Join(refs, " "), e.Cause)
}
func (e *EffectsUnrecordedError) Unwrap() error { return e.Cause }

// ReplayDivergenceError reports a replay of a committed effectful invocation
// that does not reproduce its record: Stage names the first step that
// diverged ("receipt", "record", "plan", "effects", "output").
type ReplayDivergenceError struct {
	InvocationID string
	Stage        string
	Reason       string
}

func (e *ReplayDivergenceError) Error() string {
	return fmt.Sprintf("coordinator: replay %s diverged at %s: %s", e.InvocationID, e.Stage, e.Reason)
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

// InFlightError (R17) refuses a concurrent retry of the same invocation ID.
// The caller may resend after the current call finishes to reconcile its receipt.
type InFlightError struct{ InvocationID string }

func (e *InFlightError) Error() string {
	return fmt.Sprintf("coordinator: invocation %s is in flight", e.InvocationID)
}

// IntegrityError (R13, row 122) refuses to reconcile a resolved invocation
// whose recorded object does not match the ref that names it, or whose
// record/output/world rows do not name each other. Nothing is reported as
// reconciled; the stored rows are damaged, so a resend cannot help.
type IntegrityError struct {
	InvocationID string
	Object       string // "record" | "output" | "world"
	Kind         string // "mismatch" | "absent" | "undecodable"
}

func (e *IntegrityError) Error() string {
	phrase := "does not match its reference"
	switch e.Kind {
	case "absent":
		phrase = "is absent"
	case "undecodable":
		phrase = "does not decode"
	}
	return fmt.Sprintf("coordinator: reconcile %s: recorded %s %s", e.InvocationID, e.Object, phrase)
}

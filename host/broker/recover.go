package broker

import (
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const maxRecoveryPages = 1 << 20

// IndeterminateEffectError reports a durable commit intent for which no
// outcome was recorded. Recovery only reports this ambiguity; it never calls
// an effect handler or writes a journal outcome.
type IndeterminateEffectError struct {
	InvocationID     string
	PlannedWorldRef  hashref.HashRef
	PlannedEntryHash hashref.HashRef
}

func (e *IndeterminateEffectError) Error() string {
	return fmt.Sprintf(
		"broker: effect %q is indeterminate: planned world %s and entry %s have no outcome",
		e.InvocationID, e.PlannedWorldRef.String(), e.PlannedEntryHash.String(),
	)
}

// IndeterminateEffect is one recovery finding and its durable intent.
type IndeterminateEffect struct {
	Intent store.JournalIntent
	Err    *IndeterminateEffectError
}

type recoveryStore interface {
	PendingIntents(limit int, fromIndex ...int64) ([]store.PendingIntent, error)
	GetReceipt(id string) (store.Receipt, bool, error)
}

// These consumer rules mirror the SD.C contract. The authoritative law is
// design_docs/sketches/storejournal.ail; these functions are its broker-side
// application, not an independent definition.
func mayReportNotStarted(hasIntent bool) bool {
	return !hasIntent
}

func retryAllowed(indeterminate, reconciled bool) bool {
	return !indeterminate || reconciled
}

// Recover scans durable pending intents and surfaces every indeterminate
// commit. Registries may be supplied by session construction code, but are
// intentionally never consulted: recovery does not dispatch, auto-resolve,
// append an outcome, or re-execute an effect.
//
// Every commit of an effectful episode is therefore crash-detectable and never
// auto-re-executed. The dispatch->record window remains open: an effect whose
// handler completed but whose record write was lost is not durably detectable.
// Closing that effect-granularity window belongs to queue item 4c,
// w-effect-journal.
func Recover(s *store.Store, registries ...Registry) ([]IndeterminateEffect, error) {
	// Accepting a registry makes the no-dispatch policy observable at the
	// production boundary. It is deliberately unused.
	_ = registries
	return recoverPending(s)
}

func recoverPending(s recoveryStore) ([]IndeterminateEffect, error) {
	var (
		findings []IndeterminateEffect
		cursor   int64
	)
	for pageNumber := 0; pageNumber < maxRecoveryPages; pageNumber++ {
		var (
			page []store.PendingIntent
			err  error
		)
		if cursor == 0 {
			page, err = s.PendingIntents(store.MaxPendingIntentsPage)
		} else {
			page, err = s.PendingIntents(store.MaxPendingIntentsPage, cursor)
		}
		if err != nil {
			return nil, fmt.Errorf("broker: recover pending intents: %w", err)
		}
		if len(page) == 0 {
			return findings, nil
		}

		for _, pending := range page {
			if pending.Seq <= cursor {
				return nil, fmt.Errorf(
					"broker: recovery cursor did not advance: got seq %d after %d",
					pending.Seq, cursor,
				)
			}
			cursor = pending.Seq

			receipt, hasIntent, err := s.GetReceipt(pending.InvocationID)
			if err != nil {
				return nil, fmt.Errorf(
					"broker: recover receipt %q: %w", pending.InvocationID, err,
				)
			}
			// PendingIntents just handed us this invocation, so an intent IS
			// durable for it. mayReportNotStarted(true) == false is therefore
			// the law's verdict here: the kernel may NOT report not-started.
			// It is spelled as a constant rather than as a call because a call
			// on a known-true argument folds at compile time — it would read as
			// a runtime consultation of the law while being unable to vary.
			// The law itself is pinned by TestRecoveryConsumerContractMirrorsSketch
			// against design_docs/sketches/storejournal.ail.
			if !hasIntent {
				return nil, fmt.Errorf(
					"broker: pending invocation %q was reported not-started",
					pending.InvocationID,
				)
			}
			if receipt.State != store.ReceiptIndeterminate {
				continue
			}
			// NOTE: there is deliberately no retryAllowed() guard here. Recovery
			// never retries anything, so a guard on retryAllowed(true, false)
			// would be a compile-time false — unreachable code wearing the
			// clothes of a runtime check. Proven unreachable before removal by
			// replacing its body with a panic: the whole package still passed.
			indeterminate := &IndeterminateEffectError{
				InvocationID:     pending.InvocationID,
				PlannedWorldRef:  pending.Intent.WorldRef,
				PlannedEntryHash: pending.Intent.EntryHash,
			}
			findings = append(findings, IndeterminateEffect{
				Intent: pending.Intent,
				Err:    indeterminate,
			})
		}

		if len(page) < store.MaxPendingIntentsPage {
			return findings, nil
		}
	}
	return nil, fmt.Errorf("broker: recovery exceeded %d pages", maxRecoveryPages)
}

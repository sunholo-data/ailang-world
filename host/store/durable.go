package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// UncertainError reports a durable operation whose outcome this caller cannot
// confirm: the caller's context ended after the cancellation cutoff (the
// transaction's COMMIT had been handed to database/sql), or COMMIT itself
// returned a driver error. The outcome is committed or not committed, never
// partial (one SQLite transaction). Reconcile it by reading durable state
// (GetReceipt for an invocation, the log-index witness for a bare
// Commit); never
// retry blindly.
//
// It deliberately does NOT unwrap to the context error, so
// errors.Is(err, context.Canceled|DeadlineExceeded) holds only for the
// definite "not committed" result of cancellation before the cutoff.
type UncertainError struct {
	Op    string
	Cause error
}

func (e *UncertainError) Error() string {
	return fmt.Sprintf("store: %s outcome is uncertain (reconcile before any retry): %v", e.Op, e.Cause)
}

// IsUncertain reports whether err carries an *UncertainError.
func IsUncertain(err error) bool {
	var u *UncertainError
	return errors.As(err, &u)
}

// durableCommitHook is a no-op in production. Tests replace it to stall the
// durable step before tx.Commit (the transaction is then rolled back if ctx
// ends first).
var durableCommitHook = func() {}

// durableAfterCommitHook is a no-op in production. Tests replace it to stall
// after tx.Commit has returned, modelling a slow driver COMMIT whose outcome
// the caller cannot wait for: the commit is durable, the caller sees
// *UncertainError.
var durableAfterCommitHook = func() {}

// beginDurable acquires the store's sole pooled connection and opens a
// transaction under ctx. database/sql honours ctx while waiting for the
// connection (sql.go DB.conn: select on ctx.Done() vs the conn request), so a
// held connection bounds this wait by the caller's deadline, not busy_timeout.
func (s *Store) beginDurable(ctx context.Context, op string) (*sql.Tx, error) {
	if err := s.checkQuarantine(); err != nil {
		return nil, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("store: begin %s: %w", op, err)
	}
	return tx, nil
}

// finishDurable is the cancellation cutoff. Before it, every error leaves the
// store untouched (deferred rollback, or database/sql's own ctx rollback of a
// BeginTx transaction). The caller waits for tx.Commit at most until ctx ends;
// if ctx ends first the result is *UncertainError and the one goroutine still
// finishing COMMIT holds the sole connection, so at most one such goroutine
// exists per store and the next acquisition observes the settled outcome.
//
// No ctx.Err() pre-check is needed here: database/sql's Tx.Commit checks the
// BeginTx ctx before its CompareAndSwap(done) and returns ctx.Err() without
// reaching the driver (the mutation that removed such a check survived).
func (s *Store) finishDurable(ctx context.Context, op string, tx *sql.Tx) error {
	done := make(chan error, 1)
	s.workers.Add(1)
	go func() {
		defer s.workers.Done()
		durableCommitHook()
		err := tx.Commit()
		durableAfterCommitHook()
		done <- err
	}()
	select {
	case err := <-done:
		return classifyCommit(ctx, op, err)
	case <-ctx.Done():
		select {
		case err := <-done:
			return classifyCommit(ctx, op, err)
		default:
			return &UncertainError{Op: op, Cause: ctx.Err()}
		}
	}
}

// classifyCommit maps tx.Commit's result. database/sql returns ctx.Err() or
// ErrTxDone only when it did NOT reach the driver COMMIT (sql.go Tx.Commit
// checks ctx, then CompareAndSwap(done)); those are definite "not committed".
// Any driver error is uncertain by contract.
func classifyCommit(ctx context.Context, op string, err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, sql.ErrTxDone), ctx.Err() != nil && errors.Is(err, ctx.Err()):
		if ctx.Err() == nil {
			return fmt.Errorf("store: %s not committed: %w", op, err)
		}
		return fmt.Errorf("store: %s not committed: %w", op, ctx.Err())
	default:
		return &UncertainError{Op: op, Cause: err}
	}
}

// notCommitted normalises a pre-cutoff failure caused by caller cancellation.
// database/sql rolls a BeginTx transaction back asynchronously when ctx ends,
// so a later statement not bound to ctx reports sql.ErrTxDone; the caller is
// told the cause instead. Every non-nil, non-uncertain error from a durable
// operation is a definite "not committed" either way.
func notCommitted(ctx context.Context, op string, err error) error {
	if err != nil && ctx.Err() != nil && errors.Is(err, sql.ErrTxDone) && !IsUncertain(err) {
		return fmt.Errorf("store: %s not committed: %w", op, ctx.Err())
	}
	return err
}

// CommitLanded is the bare-commit reconciliation rule INV-LOG (row 23 design
// §2.5, ratified with D-WORLD-40): after an uncertain Commit has settled, read
// the log row at c's own entry index. entry_index is the primary key and the
// only writer is Commit's INSERT (no UPDATE or DELETE exists), so:
//
//   - no row, or a row that differs from c.Entry in any stored field: c did
//     NOT commit (false, nil); resubmission is allowed;
//   - a row equal to c.Entry in every stored field: a commit carrying c's
//     exact log row is durable (true, nil) — "entry landed", not "c's world is
//     durable": the store verifies neither EntryHash nor NextWorld.Ref;
//   - a read error (including the caller's deadline): unknown (false, err).
//
// The selected head is NEVER consulted: another commit may have landed on top.
func (s *Store) CommitLanded(ctx context.Context, c Commit) (bool, error) {
	if err := s.checkQuarantine(); err != nil {
		return false, err
	}
	got, ok, err := s.GetLogEntry(ctx, c.Entry.Header.EntryIndex)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return got == c.Entry, nil
}

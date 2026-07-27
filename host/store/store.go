// Package store implements the M1 SQLite store and append-only log from
// Decision 4 of the w-world-library-m1 design.
//
// One SQLite database (opened through database/sql and the pinned pure-Go,
// CGo-free modernc.org/sqlite driver) holds:
//
//   - objects: immutable content-addressed envelopes plus payload bytes;
//   - worlds: immutable world revisions (state root + log head);
//   - log_entries: the frozen six-field LogHeader plus a separate
//     transition-body reference;
//   - epoch_registry_heads: the current immutable registry object reference;
//   - verification_cache: cached verify results keyed exactly by the pair
//     (transition_fn_ref, interpreter_ref).
//
// The load-bearing invariant is the single-transaction compare-and-append
// Commit (Decision 4): it reads the observed world head, and if the store's
// selected head has moved on it returns a structured ConflictError instead of
// appending. Object inserts verify that the stored HashRef matches the payload
// bytes through host/hashref before writing, so content addressing cannot be
// violated by a caller.
package store

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"

	_ "modernc.org/sqlite" // pure-Go, CGo-free SQLite driver ("sqlite")
)

// schemaSQL is the table definition applied on Open. Embedding keeps the schema
// text as the single source of truth shared with any external tooling.
//
//go:embed schema.sql
var schemaSQL string

// EpochRegistryV1 is the semantic ID / registry name of the epoch registry
// (Decision 5). It is the key used in epoch_registry_heads for M1.
const EpochRegistryV1 = "world/epoch-registry/v1"

// Object is the immutable ObjectEnvelope (Decision 3) together with its payload
// bytes. Hash addresses Payload; the store verifies that relationship on insert.
type Object struct {
	Hash          hashref.HashRef
	InterfaceHash hashref.HashRef
	SemanticID    string
	Provenance    string
	Payload       []byte
}

// World is one immutable world revision. Ref addresses the revision; StateRoot
// is the state object and LogHead is the append-only log head at that revision.
type World struct {
	Ref       hashref.HashRef
	Revision  int64
	StateRoot hashref.HashRef
	LogHead   hashref.HashRef
}

// LogHeader is the frozen six-field header stored verbatim in log_entries. The
// field set and order are frozen (world/logepoch.ail LogHeader); the store never
// adds to or reorders them.
type LogHeader struct {
	EntryIndex     int64
	SemanticsEpoch int64
	TransitionFn   hashref.HashRef
	Interpreter    hashref.HashRef
	PrevEntryHash  hashref.HashRef
	WrittenBy      string
}

// LogEntry is one append-only log row: the frozen header, the content address of
// the whole encoded header-plus-body-reference bytes (EntryHash), and the
// separate transition-body object reference (TransitionRef) which is OUTSIDE the
// frozen header.
type LogEntry struct {
	Header        LogHeader
	EntryHash     hashref.HashRef
	TransitionRef hashref.HashRef
}

// VerifyResult is a cached typecheck/verify outcome. The cache key is the pair
// (TransitionFn, Interpreter) EXCLUSIVELY. SemanticsEpoch is carried as
// diagnostic/migration metadata only and never participates in the key, so an
// epoch-only change to an otherwise identical pair preserves the selected row.
type VerifyResult struct {
	TransitionFn   hashref.HashRef
	Interpreter    hashref.HashRef
	SemanticsEpoch int64
	Verified       bool
	Detail         string
}

// Commit describes one compare-and-append kernel commit (Decision 4).
//
//   - ObservedHead is the world head the caller planned against; if the store's
//     selected head differs, Commit returns ConflictError.
//   - Objects are immutable objects to insert (content-verified) as part of the
//     same transaction.
//   - NextWorld is the world revision produced by the transition.
//   - Entry is the append-only log row; its PrevEntryHash must equal the log
//     head observed at ObservedHead.
type Commit struct {
	ObservedHead hashref.HashRef
	Objects      []Object
	NextWorld    World
	Entry        LogEntry
}

// ConflictError is returned by Commit when the observed world head is stale: the
// store's selected head has advanced since the caller planned. Callers may
// re-plan against SelectedHead.
type ConflictError struct {
	ObservedHead hashref.HashRef
	SelectedHead hashref.HashRef
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("store: stale observed world head %q; selected head is %q",
		e.ObservedHead.String(), e.SelectedHead.String())
}

// IsConflict reports whether err is (or wraps) a *ConflictError.
func IsConflict(err error) bool {
	var c *ConflictError
	return errors.As(err, &c)
}

// Store is a handle to the SQLite-backed world store. It is safe for the
// single-connection embedded-library use of M1; concurrency correctness rests on
// the single compare-and-append transaction, not on Go-level locking.
type Store struct {
	db  *sql.DB
	sel selectedHead
	// lock is the held single-writer lock (w-worldd-m2 Decision 2, arm A). It
	// is nil for in-memory databases (nothing to exclude across processes) and
	// for read-only handles (which never take writer authority).
	lock *writerLock
}

// selectedHead tracks the store's currently selected world head. M1 keeps this
// in the DB as an ordinary row so Commit can compare-and-append transactionally.
type selectedHead struct{}

// selectedHeadKey is the fixed key for the single selected-head row M1 uses.
const selectedHeadKey = "selected_world_head"

// Open opens (or creates) the SQLite database at path, taking sole WRITER
// authority over it, and applies the schema. path may be a filename or
// ":memory:".
//
// Single-writer, fail closed (w-worldd-m2 Decision 2 r3 — arm A, ratified
// attended by Mark, 2026-07-27): for a FILE-BACKED database, Open resolves a
// canonical database identity and takes a NON-WAITING OS-backed exclusive lock
// on "<canonical-db>.writer.lock" BEFORE any SQLite write handle exists. If
// another process already holds that lock, Open returns *WriterAlreadyActive
// immediately — it never waits, never retries and never blocks. Readers that do
// not need write authority use OpenReadOnly instead.
//
// An in-memory database ( ":memory:", "file::memory:", any DSN with
// mode=memory ) is per-connection and physically unreachable from another
// process, so it takes NO lock and creates NO lock file. That is decided on the
// resolved DSN BEFORE canonicalization.
//
// The pure-Go driver is registered as "sqlite" by the blank import above.
func Open(path string) (*Store, error) {
	if isInMemoryDSN(path) {
		db, err := openSQLite(path, path, true)
		if err != nil {
			return nil, err
		}
		return &Store{db: db}, nil
	}

	canonical, params, err := resolveDSN(path)
	if err != nil {
		return nil, err
	}
	lock, err := acquireWriterLock(canonical)
	if err != nil {
		return nil, err
	}
	db, err := openSQLite(path, writeDSN(canonical, params), true)
	if err != nil {
		// The lock was taken before SQLite ever opened; drop it, or a retry in
		// this same process would deadlock against its own descriptor.
		_ = lock.release()
		return nil, err
	}
	return &Store{db: db, lock: lock}, nil
}

// OpenReadOnly opens an EXISTING file-backed database in SQLite's read-only
// mode. It acquires NO writer lock and applies NO schema (a read-only handle
// must not attempt DDL), so it succeeds while another process holds writer
// authority over the same database.
//
// It is an error to ask for a read-only view of an in-memory DSN: such a
// database is per-connection, so a fresh read-only connection could only ever
// observe an empty database. Failing loudly beats returning an empty store.
func OpenReadOnly(path string) (*Store, error) {
	if isInMemoryDSN(path) {
		return nil, fmt.Errorf(
			"store: open read-only %q: an in-memory database is per-connection and has no "+
				"shared contents to read", path)
	}
	canonical, params, err := resolveDSN(path)
	if err != nil {
		return nil, err
	}
	db, err := openSQLite(path, readOnlyDSN(canonical, params), false)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// openSQLite opens the driver handle for dsn, reporting errors against the
// caller-facing display path. applySchema is false for read-only handles.
func openSQLite(display, dsn string, applySchema bool) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("store: open %q: %w", display, err)
	}
	// A single open connection keeps a :memory: database (which is per-connection)
	// coherent across the store's lifetime and makes the compare-and-append
	// transaction the sole serialization point.
	db.SetMaxOpenConns(1)
	// database/sql opens lazily, so this first Exec is also what surfaces an
	// unopenable database (a missing file under mode=ro, say) at construction
	// time rather than at first query.
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("store: enable foreign keys: %w", err)
	}
	if applySchema {
		if _, err := db.Exec(schemaSQL); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("store: apply schema: %w", err)
		}
	}
	return db, nil
}

// Close closes the underlying database and then releases the writer lock, if
// this handle holds one. The lock is released even when closing the database
// fails, so a failed Close can never leave writer authority stranded.
func (s *Store) Close() error {
	err := s.db.Close()
	if s.lock != nil {
		if rerr := s.lock.release(); rerr != nil && err == nil {
			err = rerr
		}
		s.lock = nil
	}
	return err
}

// verifyObject checks that o.Hash addresses o.Payload under o.Hash's algorithm.
// A mismatch is a content-addressing violation and blocks insertion.
func verifyObject(o Object) error {
	if o.Hash.IsZero() {
		return &hashref.HashError{Reason: "object hash is the zero HashRef"}
	}
	got, err := hashref.Sum(o.Hash.Algo(), o.Payload)
	if err != nil {
		return fmt.Errorf("store: verify object %q: %w", o.Hash.String(), err)
	}
	if got.Digest() != o.Hash.Digest() {
		return fmt.Errorf("store: object hash %q does not match payload (computed %q)",
			o.Hash.String(), got.String())
	}
	return nil
}

// PutObject inserts one immutable object after verifying its content address.
// Re-inserting an identical object is idempotent (INSERT OR IGNORE), since the
// bytes and hash are unchanged by definition of content addressing.
func (s *Store) PutObject(o Object) error {
	if err := verifyObject(o); err != nil {
		return err
	}
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO objects
			(hash_ref, interface_hash_ref, semantic_id, provenance, payload)
		 VALUES (?, ?, ?, ?, ?);`,
		o.Hash.String(), o.InterfaceHash.String(), o.SemanticID, o.Provenance, o.Payload,
	)
	if err != nil {
		return fmt.Errorf("store: put object %q: %w", o.Hash.String(), err)
	}
	return nil
}

// GetObject loads an immutable object by its HashRef. It returns (Object{}, nil,
// false) — via ok=false — when the object is absent.
func (s *Store) GetObject(ref hashref.HashRef) (Object, bool, error) {
	var (
		ifaceText string
		semantic  string
		prov      string
		payload   []byte
	)
	row := s.db.QueryRow(
		`SELECT interface_hash_ref, semantic_id, provenance, payload
		   FROM objects WHERE hash_ref = ?;`,
		ref.String(),
	)
	switch err := row.Scan(&ifaceText, &semantic, &prov, &payload); {
	case errors.Is(err, sql.ErrNoRows):
		return Object{}, false, nil
	case err != nil:
		return Object{}, false, fmt.Errorf("store: get object %q: %w", ref.String(), err)
	}
	iface, err := hashref.Parse(ifaceText)
	if err != nil {
		return Object{}, false, fmt.Errorf("store: object %q interface hash: %w", ref.String(), err)
	}
	return Object{
		Hash:          ref,
		InterfaceHash: iface,
		SemanticID:    semantic,
		Provenance:    prov,
		Payload:       payload,
	}, true, nil
}

// PutWorld inserts one immutable world revision. Re-inserting the same revision
// is idempotent.
func (s *Store) PutWorld(w World) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO worlds (world_ref, revision, state_root, log_head)
		 VALUES (?, ?, ?, ?);`,
		w.Ref.String(), w.Revision, w.StateRoot.String(), w.LogHead.String(),
	)
	if err != nil {
		return fmt.Errorf("store: put world %q: %w", w.Ref.String(), err)
	}
	return nil
}

// GetWorld loads a world revision by its HashRef; ok=false when absent.
func (s *Store) GetWorld(ref hashref.HashRef) (World, bool, error) {
	var (
		revision  int64
		stateText string
		headText  string
	)
	row := s.db.QueryRow(
		`SELECT revision, state_root, log_head FROM worlds WHERE world_ref = ?;`,
		ref.String(),
	)
	switch err := row.Scan(&revision, &stateText, &headText); {
	case errors.Is(err, sql.ErrNoRows):
		return World{}, false, nil
	case err != nil:
		return World{}, false, fmt.Errorf("store: get world %q: %w", ref.String(), err)
	}
	state, err := hashref.Parse(stateText)
	if err != nil {
		return World{}, false, fmt.Errorf("store: world %q state root: %w", ref.String(), err)
	}
	head, err := hashref.Parse(headText)
	if err != nil {
		return World{}, false, fmt.Errorf("store: world %q log head: %w", ref.String(), err)
	}
	return World{Ref: ref, Revision: revision, StateRoot: state, LogHead: head}, true, nil
}

// GetLogEntry loads one append-only log row by its integer index, round-tripping
// the frozen header verbatim; ok=false when absent.
func (s *Store) GetLogEntry(index int64) (LogEntry, bool, error) {
	var (
		entryHashText string
		epoch         int64
		fnText        string
		interpText    string
		prevText      string
		writtenBy     string
		transRefText  string
	)
	row := s.db.QueryRow(
		`SELECT entry_hash_ref, semantics_epoch, transition_fn_ref, interpreter_ref,
		        prev_entry_hash_ref, written_by, transition_ref
		   FROM log_entries WHERE entry_index = ?;`,
		index,
	)
	switch err := row.Scan(&entryHashText, &epoch, &fnText, &interpText,
		&prevText, &writtenBy, &transRefText); {
	case errors.Is(err, sql.ErrNoRows):
		return LogEntry{}, false, nil
	case err != nil:
		return LogEntry{}, false, fmt.Errorf("store: get log entry %d: %w", index, err)
	}
	entryHash, err := hashref.Parse(entryHashText)
	if err != nil {
		return LogEntry{}, false, fmt.Errorf("store: log entry %d hash: %w", index, err)
	}
	fn, err := hashref.Parse(fnText)
	if err != nil {
		return LogEntry{}, false, fmt.Errorf("store: log entry %d transitionFn: %w", index, err)
	}
	interp, err := hashref.Parse(interpText)
	if err != nil {
		return LogEntry{}, false, fmt.Errorf("store: log entry %d interpreter: %w", index, err)
	}
	prev, err := hashref.Parse(prevText)
	if err != nil {
		return LogEntry{}, false, fmt.Errorf("store: log entry %d prevEntryHash: %w", index, err)
	}
	transRef, err := hashref.Parse(transRefText)
	if err != nil {
		return LogEntry{}, false, fmt.Errorf("store: log entry %d transitionRef: %w", index, err)
	}
	return LogEntry{
		Header: LogHeader{
			EntryIndex:     index,
			SemanticsEpoch: epoch,
			TransitionFn:   fn,
			Interpreter:    interp,
			PrevEntryHash:  prev,
			WrittenBy:      writtenBy,
		},
		EntryHash:     entryHash,
		TransitionRef: transRef,
	}, true, nil
}

// SetRegistryHead upserts the current immutable registry object reference for a
// registry name (Decision 5). M1 bootstrap uses EpochRegistryV1.
func (s *Store) SetRegistryHead(name string, objectRef hashref.HashRef) error {
	_, err := s.db.Exec(
		`INSERT INTO epoch_registry_heads (registry_name, object_ref)
		 VALUES (?, ?)
		 ON CONFLICT(registry_name) DO UPDATE SET object_ref = excluded.object_ref;`,
		name, objectRef.String(),
	)
	if err != nil {
		return fmt.Errorf("store: set registry head %q: %w", name, err)
	}
	return nil
}

// GetRegistryHead returns the current registry object reference; ok=false when
// the registry name has no head.
func (s *Store) GetRegistryHead(name string) (hashref.HashRef, bool, error) {
	var text string
	row := s.db.QueryRow(
		`SELECT object_ref FROM epoch_registry_heads WHERE registry_name = ?;`, name)
	switch err := row.Scan(&text); {
	case errors.Is(err, sql.ErrNoRows):
		return hashref.HashRef{}, false, nil
	case err != nil:
		return hashref.HashRef{}, false, fmt.Errorf("store: get registry head %q: %w", name, err)
	}
	ref, err := hashref.Parse(text)
	if err != nil {
		return hashref.HashRef{}, false, fmt.Errorf("store: registry head %q: %w", name, err)
	}
	return ref, true, nil
}

// PutVerifyResult upserts a cached verify result. The row is keyed EXACTLY by
// (TransitionFn, Interpreter); SemanticsEpoch and the outcome fields are payload.
// An upsert with the same pair but a different epoch updates metadata in place
// and preserves the single selected row for that pair.
func (s *Store) PutVerifyResult(r VerifyResult) error {
	verified := 0
	if r.Verified {
		verified = 1
	}
	_, err := s.db.Exec(
		`INSERT INTO verification_cache
			(transition_fn_ref, interpreter_ref, semantics_epoch, verified, result_detail)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(transition_fn_ref, interpreter_ref) DO UPDATE SET
			semantics_epoch = excluded.semantics_epoch,
			verified        = excluded.verified,
			result_detail   = excluded.result_detail;`,
		r.TransitionFn.String(), r.Interpreter.String(), r.SemanticsEpoch, verified, r.Detail,
	)
	if err != nil {
		return fmt.Errorf("store: put verify result for (%q,%q): %w",
			r.TransitionFn.String(), r.Interpreter.String(), err)
	}
	return nil
}

// GetVerifyResult looks up a cached verify result by the pair
// (transitionFn, interpreter) EXCLUSIVELY; ok=false on a miss.
func (s *Store) GetVerifyResult(transitionFn, interpreter hashref.HashRef) (VerifyResult, bool, error) {
	var (
		epoch    int64
		verified int
		detail   string
	)
	row := s.db.QueryRow(
		`SELECT semantics_epoch, verified, result_detail
		   FROM verification_cache
		  WHERE transition_fn_ref = ? AND interpreter_ref = ?;`,
		transitionFn.String(), interpreter.String(),
	)
	switch err := row.Scan(&epoch, &verified, &detail); {
	case errors.Is(err, sql.ErrNoRows):
		return VerifyResult{}, false, nil
	case err != nil:
		return VerifyResult{}, false, fmt.Errorf("store: get verify result: %w", err)
	}
	return VerifyResult{
		TransitionFn:   transitionFn,
		Interpreter:    interpreter,
		SemanticsEpoch: epoch,
		Verified:       verified != 0,
		Detail:         detail,
	}, true, nil
}

// SelectedHead returns the store's currently selected world head; ok=false
// before any commit has selected a head.
func (s *Store) SelectedHead() (hashref.HashRef, bool, error) {
	return selectedHeadTx(s.db)
}

// selectedHeadTx reads the selected world head using the given querier (the DB
// or an open transaction).
func selectedHeadTx(q interface {
	QueryRow(string, ...any) *sql.Row
}) (hashref.HashRef, bool, error) {
	var text string
	row := q.QueryRow(`SELECT world_ref FROM store_heads WHERE head_key = ?;`, selectedHeadKey)
	switch err := row.Scan(&text); {
	case errors.Is(err, sql.ErrNoRows):
		return hashref.HashRef{}, false, nil
	case err != nil:
		return hashref.HashRef{}, false, fmt.Errorf("store: read selected head: %w", err)
	}
	ref, err := hashref.Parse(text)
	if err != nil {
		return hashref.HashRef{}, false, fmt.Errorf("store: selected head: %w", err)
	}
	return ref, true, nil
}

// SelectHead sets the store's selected world head without a full commit. It is
// used to seed the genesis head before the first Commit; steady-state head
// advancement happens inside Commit.
func (s *Store) SelectHead(ref hashref.HashRef) error {
	_, err := s.db.Exec(
		`INSERT INTO store_heads (head_key, world_ref) VALUES (?, ?)
		 ON CONFLICT(head_key) DO UPDATE SET world_ref = excluded.world_ref;`,
		selectedHeadKey, ref.String(),
	)
	if err != nil {
		return fmt.Errorf("store: select head %q: %w", ref.String(), err)
	}
	return nil
}

// Commit performs the single-transaction compare-and-append of Decision 4:
//
//  1. Read the store's selected world head; if it differs from c.ObservedHead
//     return a *ConflictError (the caller planned against a stale head).
//  2. Insert every required immutable object with content verification.
//  3. Insert the next immutable world row.
//  4. Insert the append-only log row (frozen header + transition-body ref).
//  5. Advance the selected world head to c.NextWorld.
//  6. Commit the SQLite transaction.
//
// Any error rolls the whole transaction back, so a conflict or a bad object
// leaves the store untouched. A nil selected head (genesis) is treated as a
// match only when c.ObservedHead is also the zero HashRef.
func (s *Store) Commit(c Commit) error {
	// Content-verify all objects before opening the transaction so a bad object
	// never even starts a write.
	for _, o := range c.Objects {
		if err := verifyObject(o); err != nil {
			return err
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("store: begin commit: %w", err)
	}
	// Roll back on any early return; a nil-error Commit calls tx.Commit first,
	// making this rollback a no-op.
	defer func() { _ = tx.Rollback() }()

	// Step 1: compare-and-append guard.
	selected, hasSelected, err := selectedHeadTx(tx)
	if err != nil {
		return err
	}
	if hasSelected {
		if selected.String() != c.ObservedHead.String() {
			return &ConflictError{ObservedHead: c.ObservedHead, SelectedHead: selected}
		}
	} else if !c.ObservedHead.IsZero() {
		// No head selected yet, but the caller observed a non-genesis head.
		return &ConflictError{ObservedHead: c.ObservedHead, SelectedHead: hashref.HashRef{}}
	}

	// Step 2: immutable objects.
	for _, o := range c.Objects {
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO objects
				(hash_ref, interface_hash_ref, semantic_id, provenance, payload)
			 VALUES (?, ?, ?, ?, ?);`,
			o.Hash.String(), o.InterfaceHash.String(), o.SemanticID, o.Provenance, o.Payload,
		); err != nil {
			return fmt.Errorf("store: commit object %q: %w", o.Hash.String(), err)
		}
	}

	// Step 3: next immutable world row.
	w := c.NextWorld
	if _, err := tx.Exec(
		`INSERT OR IGNORE INTO worlds (world_ref, revision, state_root, log_head)
		 VALUES (?, ?, ?, ?);`,
		w.Ref.String(), w.Revision, w.StateRoot.String(), w.LogHead.String(),
	); err != nil {
		return fmt.Errorf("store: commit world %q: %w", w.Ref.String(), err)
	}

	// Step 4: append-only log row (frozen header verbatim + separate body ref).
	h := c.Entry.Header
	if _, err := tx.Exec(
		`INSERT INTO log_entries
			(entry_index, entry_hash_ref, semantics_epoch, transition_fn_ref,
			 interpreter_ref, prev_entry_hash_ref, written_by, transition_ref)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?);`,
		h.EntryIndex, c.Entry.EntryHash.String(), h.SemanticsEpoch,
		h.TransitionFn.String(), h.Interpreter.String(), h.PrevEntryHash.String(),
		h.WrittenBy, c.Entry.TransitionRef.String(),
	); err != nil {
		return fmt.Errorf("store: commit log entry %d: %w", h.EntryIndex, err)
	}

	// Step 5: advance the selected world head.
	if _, err := tx.Exec(
		`INSERT INTO store_heads (head_key, world_ref) VALUES (?, ?)
		 ON CONFLICT(head_key) DO UPDATE SET world_ref = excluded.world_ref;`,
		selectedHeadKey, w.Ref.String(),
	); err != nil {
		return fmt.Errorf("store: advance selected head: %w", err)
	}

	// Step 6: commit.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit transaction: %w", err)
	}
	return nil
}

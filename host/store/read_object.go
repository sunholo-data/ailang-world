package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// ObjectMeta is the non-payload metadata returned by a bounded object read.
type ObjectMeta struct {
	InterfaceHash hashref.HashRef
	SemanticID    string
	Provenance    string
	PayloadLength int64
}

// ObjectTooLargeError reports a payload refused before it was materialized.
type ObjectTooLargeError struct {
	Size     int64
	MaxBytes int64
}

func (e *ObjectTooLargeError) Error() string {
	return fmt.Sprintf("store: object payload is %d bytes; maximum is %d", e.Size, e.MaxBytes)
}

const readObjectProbeSQL = `SELECT interface_hash_ref, semantic_id, provenance, length(payload)
FROM objects WHERE hash_ref = ?`

const readObjectPayloadSQL = `SELECT payload FROM objects WHERE hash_ref = ?`

// readObjectBetweenStatements is a test-only scheduling input. Production
// leaves it nil; when installed it wraps, but does not replace, the interval
// between the length probe and payload statement.
var readObjectBetweenStatements func()

// ReadObject reads an object only when its payload fits maxBytes. The length
// probe and payload materialization share one transaction and therefore one
// SQLite snapshot.
func (s *Store) ReadObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (ObjectMeta, []byte, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return ObjectMeta{}, nil, fmt.Errorf("store: read object %q: reserve connection: %w", ref.String(), err)
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return ObjectMeta{}, nil, fmt.Errorf("store: read object %q: begin snapshot: %w", ref.String(), err)
	}
	defer tx.Rollback()

	var ifaceText string
	var meta ObjectMeta
	err = tx.QueryRowContext(ctx, readObjectProbeSQL, ref.String()).Scan(
		&ifaceText, &meta.SemanticID, &meta.Provenance, &meta.PayloadLength,
	)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ObjectMeta{}, nil, nil
	case err != nil:
		return ObjectMeta{}, nil, fmt.Errorf("store: read object %q probe: %w", ref.String(), err)
	}
	meta.InterfaceHash, err = hashref.Parse(ifaceText)
	if err != nil {
		return ObjectMeta{}, nil, fmt.Errorf("store: object %q interface hash: %w", ref.String(), err)
	}
	if meta.PayloadLength > maxBytes {
		return meta, nil, &ObjectTooLargeError{Size: meta.PayloadLength, MaxBytes: maxBytes}
	}

	if readObjectBetweenStatements != nil {
		readObjectBetweenStatements()
	}

	var payload []byte
	if err := tx.QueryRowContext(ctx, readObjectPayloadSQL, ref.String()).Scan(&payload); err != nil {
		return ObjectMeta{}, nil, fmt.Errorf("store: read object %q payload: %w", ref.String(), err)
	}
	if err := tx.Commit(); err != nil {
		return ObjectMeta{}, nil, fmt.Errorf("store: read object %q commit snapshot: %w", ref.String(), err)
	}
	return meta, payload, nil
}

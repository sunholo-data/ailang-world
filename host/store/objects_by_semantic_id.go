package store

import (
	"context"
	"fmt"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// MaxSemanticIDPage is the kernel-owned upper bound for one semanticId lookup
// page. It equals the daemon's clampLimit ceiling, so the route never asks for
// more than the store will serve.
const MaxSemanticIDPage = 500

// objectsBySemanticIDSQL is a keyset page over the objects_by_semantic_id
// index provisioned at writable Open. It never reads payload, so a page's allocation is
// bounded by envelope metadata only, and it never depends on OFFSET or SQLite
// rowids: hash_ref is the PRIMARY KEY, so it is a unique total order.
const objectsBySemanticIDSQL = `SELECT hash_ref, interface_hash_ref, semantic_id, provenance
FROM objects WHERE semantic_id = ? AND hash_ref > ? ORDER BY hash_ref LIMIT ?`

// ObjectsBySemanticID returns up to limit objects whose semantic_id equals id
// and whose hash_ref sorts strictly after after (use "" for the first page),
// ascending by hash_ref. Payload is always nil: callers fetch the one payload
// they want with GetObject. An id no object carries yields an empty slice.
func (s *Store) ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]Object, error) {
	if err := s.checkQuarantine(); err != nil {
		return nil, err
	}
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	const op = "ObjectsBySemanticID"
	if limit < 1 || limit > MaxSemanticIDPage {
		return nil, &InvalidLimitError{Op: op, Limit: limit, Max: MaxSemanticIDPage}
	}
	if !s.lookupIndexAvailable {
		return nil, &LookupIndexUnavailableError{}
	}
	if objectsBySemanticIDBeforeQuery != nil {
		objectsBySemanticIDBeforeQuery()
	}
	rows, err := s.db.QueryContext(ctx, objectsBySemanticIDSQL, id, after, limit)
	if err != nil {
		return nil, fmt.Errorf("store: objects by semantic id %q: %w", id, err)
	}
	defer rows.Close()
	objects := make([]Object, 0, limit)
	for rows.Next() {
		var hashText, ifaceText, storedID, prov string
		if err := rows.Scan(&hashText, &ifaceText, &storedID, &prov); err != nil {
			return nil, fmt.Errorf("store: objects by semantic id %q: scan: %w", id, err)
		}
		hash, err := hashref.Parse(hashText)
		if err != nil {
			return nil, fmt.Errorf("store: object %q hash: %w", hashText, err)
		}
		iface, err := hashref.Parse(ifaceText)
		if err != nil {
			return nil, fmt.Errorf("store: object %q interface hash: %w", hashText, err)
		}
		objects = append(objects, Object{Hash: hash, InterfaceHash: iface, SemanticID: storedID, Provenance: prov})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: objects by semantic id %q: %w", id, err)
	}
	return objects, nil
}

// objectsBySemanticIDBeforeQuery is a test-only hook immediately before lookup SQL.
var objectsBySemanticIDBeforeQuery func()

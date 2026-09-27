package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const lookupIndexUnavailableMessage = "semantic-id lookup index is absent or incompatible on this read-only store; open the store writable once (for example start ailang-worldd serve on it) to provision objects_by_semantic_id"

// lookupIndexProvisionDeadline is a variable so tests can exercise the deadline path.
var lookupIndexProvisionDeadline = 30 * time.Second

type LookupIndexUnavailableError struct{}

func (*LookupIndexUnavailableError) Error() string { return lookupIndexUnavailableMessage }

func provisionLookupIndex(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), lookupIndexProvisionDeadline)
	defer cancel()
	if _, err := db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS objects_by_semantic_id ON objects(semantic_id, hash_ref)"); err != nil {
		return fmt.Errorf("store: provision semantic-id lookup index: %w", err)
	}
	available, err := verifyLookupIndex(db)
	if err != nil {
		return err
	}
	if !available {
		return fmt.Errorf("store: incompatible semantic-id lookup index")
	}
	return nil
}

func verifyLookupIndex(db *sql.DB) (bool, error) {
	rows, err := db.Query("PRAGMA index_list('objects')")
	if err != nil {
		return false, fmt.Errorf("store: list object indexes: %w", err)
	}
	found := false
	for rows.Next() {
		var seq, unique, partial int
		var name, origin string
		if err = rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			break
		}
		if name == "objects_by_semantic_id" {
			found = partial == 0
			break
		}
	}
	if err == nil {
		err = rows.Err()
	}
	rows.Close()
	if err != nil {
		return false, fmt.Errorf("store: inspect object indexes: %w", err)
	}
	if !found {
		return false, nil
	}

	rows, err = db.Query("PRAGMA index_xinfo('objects_by_semantic_id')")
	if err != nil {
		return false, fmt.Errorf("store: inspect semantic-id index: %w", err)
	}
	defer rows.Close()
	wants := []struct{ name string }{{"semantic_id"}, {"hash_ref"}}
	keys := 0
	valid := true
	for rows.Next() {
		var seqno, cid, desc, key int
		var nameText sql.NullString
		var coll string
		if err := rows.Scan(&seqno, &cid, &nameText, &desc, &coll, &key); err != nil {
			return false, fmt.Errorf("store: scan semantic-id index: %w", err)
		}
		if key == 0 {
			continue
		}
		name := nameText.String
		if keys >= len(wants) {
			valid = false
			continue
		}
		want := wants[keys]
		if !(seqno == keys && name != "" && want.name != "" && name == want.name && coll == "BINARY" && desc == 0) {
			valid = false
		}
		keys++
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("store: read semantic-id index: %w", err)
	}
	return valid && keys == len(wants), nil
}

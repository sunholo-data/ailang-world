package store

import (
	"context"
	"database/sql"
	"fmt"
)

type referenceIndexSpec struct{ name, table, relation, sourceKey string }

var referenceIndexSpecs = []referenceIndexSpec{
	{"log_entries_by_transition_ref", "log_entries", "transition_ref", "entry_index"},
	{"log_entries_by_transition_fn_ref", "log_entries", "transition_fn_ref", "entry_index"},
	{"log_entries_by_interpreter_ref", "log_entries", "interpreter_ref", "entry_index"},
	{"worlds_by_state_root", "worlds", "state_root", "world_ref"},
}

// verifyReferenceIndexes checks the complete access path without changing it.
func verifyReferenceIndexes(ctx context.Context, db *sql.DB) (bool, error) {
	for _, spec := range referenceIndexSpecs {
		ok, err := verifyReferenceIndex(ctx, db, spec)
		if err != nil {
			return false, fmt.Errorf("store: inspect reference index %s: %w", spec.name, err)
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func verifyReferenceIndex(ctx context.Context, db *sql.DB, spec referenceIndexSpec) (bool, error) {
	rows, err := db.QueryContext(ctx, "PRAGMA index_list('"+spec.table+"')")
	if err != nil {
		return false, err
	}
	found := false
	for rows.Next() {
		var seq, unique, partial int
		var name, origin string
		if err = rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			break
		}
		if name == spec.name {
			found = partial == 0
		}
	}
	if err == nil {
		err = rows.Err()
	}
	closeErr := rows.Close()
	if err != nil {
		return false, err
	}
	if closeErr != nil {
		return false, closeErr
	}
	if !found {
		return false, nil
	}

	rows, err = db.QueryContext(ctx, "PRAGMA index_xinfo('"+spec.name+"')")
	if err != nil {
		return false, err
	}
	keys := 0
	valid := true
	wants := [2]string{spec.relation, spec.sourceKey}
	for rows.Next() {
		var seq, cid, desc, key int
		var name sql.NullString
		var coll string
		if err = rows.Scan(&seq, &cid, &name, &desc, &coll, &key); err != nil {
			break
		}
		if key == 0 {
			continue
		}
		if keys >= 2 {
			valid = false
			keys++
			continue
		}
		if seq != keys || cid < 0 || !name.Valid || name.String != wants[keys] || !(coll == "BINARY") || !(desc == 0) {
			valid = false
		}
		keys++
	}
	if err == nil {
		err = rows.Err()
	}
	closeErr = rows.Close()
	if err != nil {
		return false, err
	}
	if closeErr != nil {
		return false, closeErr
	}
	return valid && keys == 2, nil
}

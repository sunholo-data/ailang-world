package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
)

func TestReferenceIndexDefinitions(t *testing.T) {
	ddl := func(s *Store, name, statement string) {
		t.Helper()
		if _, err := s.db.Exec("DROP INDEX IF EXISTS " + name); err != nil {
			t.Fatal(err)
		}
		if statement != "" {
			if _, err := s.db.Exec(statement); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Run("all-four", func(t *testing.T) {
		s := openMem(t)
		for _, spec := range referenceIndexSpecs {
			_, err := s.db.Exec("CREATE INDEX IF NOT EXISTS " + spec.name + " ON " + spec.table + "(" + spec.relation + "," + spec.sourceKey + ")")
			if err != nil {
				t.Fatal(err)
			}
		}
		if len(referenceIndexSpecs) != 4 {
			t.Fatalf("spec count %d", len(referenceIndexSpecs))
		}
		for _, name := range []string{"log_entries_by_transition_ref", "log_entries_by_transition_fn_ref", "log_entries_by_interpreter_ref", "worlds_by_state_root"} {
			found := false
			for _, spec := range referenceIndexSpecs {
				if spec.name == name {
					found = true
				}
			}
			if !found {
				t.Fatalf("missing %s", name)
			}
		}
		ok, err := verifyReferenceIndexes(context.Background(), s.db)
		if err != nil || !ok {
			t.Fatalf("all four: %v %v", ok, err)
		}
	})
	tests := []struct{ name, statement string }{
		{"absent", ""},
		{"partial", "CREATE INDEX log_entries_by_transition_ref ON log_entries(transition_ref,entry_index) WHERE transition_ref <> ''"},
		{"nocase", "CREATE INDEX log_entries_by_transition_ref ON log_entries(transition_ref COLLATE NOCASE,entry_index)"},
		{"descending", "CREATE INDEX log_entries_by_transition_ref ON log_entries(transition_ref,entry_index DESC)"},
		{"extra-key", "CREATE INDEX log_entries_by_transition_ref ON log_entries(transition_ref,entry_index,written_by)"},
		{"reversed", "CREATE INDEX log_entries_by_transition_ref ON log_entries(entry_index,transition_ref)"},
		{"expression", "CREATE INDEX log_entries_by_transition_ref ON log_entries(lower(transition_ref),entry_index)"},
		{"wrong-table", "CREATE INDEX log_entries_by_transition_ref ON worlds(state_root,world_ref)"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			ddl(s, "log_entries_by_transition_ref", tc.statement)
			ok, err := verifyReferenceIndex(context.Background(), s.db, referenceIndexSpecs[0])
			if err != nil || ok {
				t.Fatalf("%s: %v %v", tc.name, ok, err)
			}
		})
	}
}

func TestReferenceIndexVerifierContext(t *testing.T) {
	s := openMem(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := verifyReferenceIndexes(ctx, s.db)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled verifier: %v", err)
	}
	// Every PRAGMA must share the caller context, including index_xinfo.
	if _, err := s.db.Exec("CREATE INDEX IF NOT EXISTS log_entries_by_transition_ref ON log_entries(transition_ref,entry_index)"); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	_, err = verifyReferenceIndexes(ctx, db)
	if err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("cancelled fresh connection: %v", err)
	}
}

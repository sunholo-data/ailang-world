package store

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
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
		ok, err := verifyReferenceIndexes(boundedTestContext(t), s.db)
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
			ok, err := verifyReferenceIndex(boundedTestContext(t), s.db, referenceIndexSpecs[0])
			if err != nil || ok {
				t.Fatalf("%s: %v %v", tc.name, ok, err)
			}
		})
	}
}

func TestReferenceIndexVerifierContext(t *testing.T) {
	s := openMem(t)
	ctx, cancel := context.WithCancel(boundedTestContext(t))
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

func TestReferenceIndexLifecycle(t *testing.T) {
	t.Run("memory", func(t *testing.T) {
		s := openMem(t)
		if !s.referenceIndexesAvailable {
			t.Fatal("memory indexes unavailable")
		}
		ok, err := verifyReferenceIndexes(boundedTestContext(t), s.db)
		if err != nil || !ok {
			t.Fatalf("verify: %v %v", ok, err)
		}
	})
	t.Run("read-only-absent", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "refs.db")
		s, err := Open(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = s.db.Exec("DROP INDEX worlds_by_state_root"); err != nil {
			t.Fatal(err)
		}
		if err = s.Close(); err != nil {
			t.Fatal(err)
		}
		ro, err := OpenReadOnly(path)
		if err != nil {
			t.Fatal(err)
		}
		if ro.referenceIndexesAvailable {
			t.Fatal("read-only cache true without all indexes")
		}
		if !ro.lookupIndexAvailable {
			t.Fatal("semantic index lost")
		}
		ro.Close()
		s, err = Open(path)
		if err != nil {
			t.Fatal(err)
		}
		s.Close()
		ro, err = OpenReadOnly(path)
		if err != nil {
			t.Fatal(err)
		}
		defer ro.Close()
		if !ro.referenceIndexesAvailable {
			t.Fatal("reopened read-only cache false")
		}
	})
	t.Run("version", func(t *testing.T) {
		s := openMem(t)
		var version int
		if err := s.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != currentSchemaVersion {
			t.Fatalf("version %d %v", version, err)
		}
	})
}

func TestReferenceIndexFailureCleanup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.db.Exec("DROP INDEX worlds_by_state_root"); err != nil {
		t.Fatal(err)
	}
	if _, err = s.db.Exec("CREATE INDEX worlds_by_state_root ON worlds(log_head,state_root)"); err != nil {
		t.Fatal(err)
	}
	s.Close()
	failed, err := Open(path)
	if failed != nil || err == nil {
		t.Fatalf("incompatible Open=%v %v", failed, err)
	}
	// Lock and SQLite handle must be released immediately, even after failure.
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("DROP INDEX worlds_by_state_root"); err != nil {
		t.Fatal(err)
	}
	db.Close()
	s, err = Open(path)
	if err != nil {
		t.Fatalf("retry after failure: %v", err)
	}
	defer s.Close()
	if !s.referenceIndexesAvailable {
		t.Fatal("retry did not provision")
	}
}

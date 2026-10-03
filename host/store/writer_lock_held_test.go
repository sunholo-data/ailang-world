package store

import (
	"os"
	"path/filepath"
	"testing"
)

// TestWriterLockHeldIsNonMutating pins row 138 M4's lock probe: it reports a
// held and a free lock correctly, creates no file on a missing store, and
// takes nothing (a writer can open right after a probe).
func TestWriterLockHeldIsNonMutating(t *testing.T) {
	dir := t.TempDir()
	db := filepath.Join(dir, "world.db")
	held, err := WriterLockHeld(db)
	if err != nil || held {
		t.Fatalf("missing store: held=%v err=%v", held, err)
	}
	if ents, _ := os.ReadDir(dir); len(ents) != 0 {
		t.Fatalf("WriterLockHeld created %d file(s)", len(ents))
	}
	w, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	if held, err := WriterLockHeld(db); err != nil || !held {
		t.Fatalf("open writer: held=%v err=%v, want held", held, err)
	}
	// The probe must not have taken anything: the writer is still the writer,
	// and once it closes a new writer can open at once.
	_ = w.Close()
	if held, err := WriterLockHeld(db); err != nil || held {
		t.Fatalf("closed writer: held=%v err=%v, want free", held, err)
	}
	w2, err := Open(db)
	if err != nil {
		t.Fatalf("a probe left the lock taken: %v", err)
	}
	_ = w2.Close()
}

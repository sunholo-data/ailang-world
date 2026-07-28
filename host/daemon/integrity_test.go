package daemon

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	_ "modernc.org/sqlite"
)

func integrityFixture(t *testing.T) string {
	t.Helper()
	path := t.TempDir() + "/world.db"
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	ref := hashref.SumSHA256([]byte("integrity genesis"))
	genesis := store.World{
		Ref: ref, StateRoot: hashref.SumSHA256([]byte("integrity state")),
		LogHead: hashref.SumSHA256([]byte("integrity head")),
	}
	if err := s.PutWorld(genesis); err != nil {
		t.Fatal(err)
	}
	if err := s.SelectHead(ref); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	valid := hashref.SumSHA256([]byte("integrity valid")).String()
	for i := int64(1); i <= 70; i++ {
		entry := hashref.SumSHA256([]byte{byte(i), 1}).String()
		prev := valid
		if i == 66 {
			prev = ""
		}
		if _, err := tx.Exec(`INSERT INTO log_entries
			(entry_index,entry_hash_ref,semantics_epoch,transition_fn_ref,interpreter_ref,
			 prev_entry_hash_ref,written_by,transition_ref) VALUES(?,?,?,?,?,?,?,?)`,
			i, entry, 1, valid, valid, prev, "integrity", valid); err != nil {
			t.Fatal(err)
		}
	}
	worldRef := hashref.SumSHA256([]byte("integrity poison world")).String()
	if _, err := tx.Exec(`INSERT INTO worlds(world_ref,revision,state_root,log_head)
		VALUES(?,?,?,?)`, worldRef, 2, "", valid); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestIntegrityStartupSweepPagesBothTables(t *testing.T) {
	d, err := New(Config{DBPath: integrityFixture(t), BindHost: DefaultBindHost, BindPort: 0})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	r := d.IntegrityReport()
	if !r.Complete || r.LogRowsScanned != 70 || r.WorldRowsScanned != 2 || len(r.Holes) != 2 {
		t.Fatalf("report = %+v", r)
	}
	lines := strings.Join(r.Lines(), "\n")
	for _, want := range []string{
		"integrity_hole table=log_entries index=66 field=prevEntryHash",
		"integrity_hole table=worlds ref=",
		" field=stateRoot",
		"integrity_scan_complete log_rows=70 world_rows=2 holes=2",
	} {
		if !strings.Contains(lines, want) {
			t.Fatalf("lines %q do not contain %q", lines, want)
		}
	}
}

func TestIntegrityStartupSweepReportsTruncation(t *testing.T) {
	d, err := New(Config{DBPath: integrityFixture(t), BindHost: DefaultBindHost, BindPort: 0})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	d.scanPageSize = 64
	d.scanRowBudget = 32
	d.scanTimeBudget = time.Second
	d.integrity = d.scanIntegrity()
	r := d.IntegrityReport()
	if r.Complete || r.LogRowsScanned != 32 {
		t.Fatalf("report = %+v", r)
	}
	want := "integrity_scan_incomplete log_rows=32 world_rows=0 holes=0 resume_log_index=33 resume_world_ref="
	if got := r.Lines()[len(r.Lines())-1]; got != want {
		t.Fatalf("line = %q, want %q", got, want)
	}
}

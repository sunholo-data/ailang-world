package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"

	_ "modernc.org/sqlite"
)

func testCommitIntent(id string, c Commit) JournalIntent {
	return JournalIntent{
		InvocationID: id, WorldRef: c.NextWorld.Ref, EntryHash: c.Entry.EntryHash,
		ObservedHead: c.ObservedHead, PrevEntryHash: c.Entry.Header.PrevEntryHash,
		TransitionFn: c.Entry.Header.TransitionFn, TransitionRef: c.Entry.TransitionRef,
		Interpreter: c.Entry.Header.Interpreter, LogicalTime: 42,
	}
}

func journalCommitFixture(t *testing.T, s *Store, id string) Commit {
	t.Helper()
	genesis := seedGenesis(t, s)
	body := obj("journal-body-"+id, "transition/body")
	entryHash := hashref.SumSHA256([]byte("journal-entry-" + id))
	return Commit{
		InvocationID: id, ObservedHead: genesis.Ref, Objects: []Object{body},
		NextWorld: World{
			Ref: hashref.SumSHA256([]byte("journal-world-" + id)), Revision: 1,
			StateRoot: hashref.SumSHA256([]byte("journal-state-" + id)), LogHead: entryHash,
		},
		Entry: LogEntry{
			Header: LogHeader{
				EntryIndex: 1, SemanticsEpoch: 1,
				TransitionFn:  hashref.SumSHA256([]byte("journal-fn-" + id)),
				Interpreter:   hashref.SumSHA256([]byte("journal-interpreter-" + id)),
				PrevEntryHash: genesis.LogHead, WrittenBy: "journal-test",
			},
			EntryHash: entryHash, TransitionRef: body.Hash,
		},
	}
}

func TestReceiptStateDriftAllBooleanCombinations(t *testing.T) {
	tests := []struct {
		name, hasIntent, hasOutcome string
		want                        ReceiptState
		corrupt                     bool
	}{
		{"not-started", "no", "no", ReceiptNotStarted, false},
		{"indeterminate", "yes", "no", ReceiptIndeterminate, false},
		{"resolved", "yes", "yes", ReceiptResolved, false},
		{"corrupt-unrepresentable", "no", "yes", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			id := "receipt-" + tc.name
			c := journalCommitFixture(t, s, id)
			intent := testCommitIntent(id, c)
			if tc.hasIntent == "yes" {
				if _, _, err := s.AppendIntent(id, intent); err != nil {
					t.Fatal(err)
				}
			}
			if tc.hasOutcome == "yes" {
				outcome := JournalOutcome{id, "committed", c.NextWorld.Ref, 43}
				_, _, err := s.AppendOutcome(id, outcome)
				if tc.corrupt {
					if err == nil || !IsInvocationMismatch(err) {
						t.Fatalf("outcome without intent = %v, want structured error", err)
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
			}
			got, ok, err := s.GetReceipt(id)
			if err != nil {
				t.Fatal(err)
			}
			if got.State != tc.want {
				t.Fatalf("state = %q, want %q", got.State, tc.want)
			}
			if tc.hasIntent == "yes" && (!ok || got.State == ReceiptNotStarted) {
				t.Fatal("mayReportNotStarted violated for durable intent")
			}
		})
	}
}

func TestIntentBindingMirrorsAllTenSketchRows(t *testing.T) {
	type mutate func(*Commit)
	different := func(label string) hashref.HashRef {
		return hashref.SumSHA256([]byte("different-" + label))
	}
	tests := []struct {
		name      string
		field     string
		mutate    mutate
		wantMatch bool
	}{
		{"row-1-all-match", "", func(*Commit) {}, true},
		{"row-2-invocation-id", "InvocationID", func(c *Commit) { c.InvocationID += "-B" }, false},
		{"row-3-world-ref", "WorldRef", func(c *Commit) { c.NextWorld.Ref = different("world") }, false},
		{"row-4-entry-hash", "EntryHash", func(c *Commit) { c.Entry.EntryHash = different("entry") }, false},
		{"row-5-observed-head", "ObservedHead", func(c *Commit) { c.ObservedHead = different("head") }, false},
		{"row-6-prev-entry-hash", "PrevEntryHash", func(c *Commit) { c.Entry.Header.PrevEntryHash = different("prev") }, false},
		{"row-7-transition-fn", "TransitionFn", func(c *Commit) { c.Entry.Header.TransitionFn = different("fn") }, false},
		{"row-8-transition-ref", "TransitionRef", func(c *Commit) { c.Entry.TransitionRef = different("transition") }, false},
		{"row-9-interpreter", "Interpreter", func(c *Commit) { c.Entry.Header.Interpreter = different("interpreter") }, false},
		{"row-10-entryhash-preserving-combined", "PrevEntryHash", func(c *Commit) {
			c.Entry.Header.PrevEntryHash = different("combined-prev")
			c.Entry.Header.TransitionFn = different("combined-fn")
			c.Entry.Header.Interpreter = different("combined-interpreter")
		}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := openMem(t)
			id := "binding-" + tc.name
			c := journalCommitFixture(t, s, id)
			if _, _, err := s.AppendIntent(id, testCommitIntent(id, c)); err != nil {
				t.Fatal(err)
			}
			before := snapshotJournalStore(t, s)
			tc.mutate(&c)
			err := s.Commit(c)
			if tc.wantMatch {
				if err != nil {
					t.Fatalf("matching Commit: %v", err)
				}
				receipt, _, err := s.GetReceipt(id)
				if err != nil || receipt.State != ReceiptResolved {
					t.Fatalf("receipt = %+v, err=%v", receipt, err)
				}
				afterFirst := snapshotJournalStore(t, s)
				if err := s.Commit(c); err != nil {
					t.Fatalf("resolved idempotent Commit: %v", err)
				}
				if after := snapshotJournalStore(t, s); !reflect.DeepEqual(after, afterFirst) {
					t.Fatalf("resolved re-commit mutated store:\n%v\n%v", afterFirst, after)
				}
				return
			}
			var mismatch *InvocationMismatchError
			if !errors.As(err, &mismatch) || mismatch.Field != tc.field {
				t.Fatalf("Commit error = %T %v, want mismatch field %s", err, err, tc.field)
			}
			if after := snapshotJournalStore(t, s); !reflect.DeepEqual(after, before) {
				t.Fatalf("mismatch mutated store:\n before=%v\n after=%v", before, after)
			}
		})
	}
}

type storeSnapshot struct {
	Head    string
	Objects int
	Worlds  int
	Logs    int
	Journal []string
}

func snapshotJournalStore(t *testing.T, s *Store) storeSnapshot {
	t.Helper()
	var snap storeSnapshot
	head, ok, err := s.SelectedHead()
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		snap.Head = head.String()
	}
	for table, target := range map[string]*int{
		"objects": &snap.Objects, "worlds": &snap.Worlds, "log_entries": &snap.Logs,
	} {
		if err := s.db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(target); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := s.db.Query(`SELECT printf('%d|%s|%s|%s',seq,kind,invocation_id,object_ref)
		FROM journal ORDER BY seq`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var row string
		if err := rows.Scan(&row); err != nil {
			t.Fatal(err)
		}
		snap.Journal = append(snap.Journal, row)
	}
	return snap
}

func TestAppendIntentIdempotencyDuplicateAndSchemaOutcomeUniqueness(t *testing.T) {
	s := openMem(t)
	c := journalCommitFixture(t, s, "idem")
	intent := testCommitIntent("idem", c)
	seq, ref, err := s.AppendIntent("idem", intent)
	if err != nil {
		t.Fatal(err)
	}
	seq2, ref2, err := s.AppendIntent("idem", intent)
	if err != nil || seq2 != seq || ref2 != ref {
		t.Fatalf("idempotent append = (%d,%s,%v), want (%d,%s,nil)", seq2, ref2, err, seq, ref)
	}
	changed := intent
	changed.LogicalTime++
	if _, _, err := s.AppendIntent("idem", changed); !IsDuplicateInvocation(err) {
		t.Fatalf("different bytes error = %T %v", err, err)
	}
	outcome := JournalOutcome{"idem", "committed", c.NextWorld.Ref, 44}
	_, outcomeRef, err := s.AppendOutcome("idem", outcome)
	if err != nil {
		t.Fatal(err)
	}
	var next int64
	if err := s.db.QueryRow(`SELECT MAX(seq)+1 FROM journal`).Scan(&next); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`INSERT INTO journal(seq,kind,invocation_id,object_ref)
		VALUES (?,'outcome','idem',?)`, next, outcomeRef.String()); err == nil {
		t.Fatal("raw SQL duplicate outcome bypassed UNIQUE(invocation_id,kind)")
	}
}

func TestJournalSequenceGaplessAfterInducedRollback(t *testing.T) {
	s := openMem(t)
	appendID := func(id string) {
		c := journalCommitFixture(t, s, id)
		if _, _, err := s.AppendIntent(id, testCommitIntent(id, c)); err != nil {
			t.Fatal(err)
		}
	}
	appendID("seq-1")
	if _, err := s.db.Exec(`CREATE TRIGGER induce_journal_rollback BEFORE INSERT ON journal
		WHEN NEW.invocation_id = 'rollback' BEGIN SELECT RAISE(ABORT, 'induced'); END`); err != nil {
		t.Fatal(err)
	}
	c := journalCommitFixture(t, s, "rollback")
	if _, _, err := s.AppendIntent("rollback", testCommitIntent("rollback", c)); err == nil {
		t.Fatal("induced rollback unexpectedly succeeded")
	}
	if _, err := s.db.Exec(`DROP TRIGGER induce_journal_rollback`); err != nil {
		t.Fatal(err)
	}
	appendID("seq-2")
	rows, err := s.db.Query(`SELECT seq FROM journal ORDER BY seq`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var got []int64
	for rows.Next() {
		var seq int64
		if err := rows.Scan(&seq); err != nil {
			t.Fatal(err)
		}
		got = append(got, seq)
	}
	if !reflect.DeepEqual(got, []int64{1, 2}) {
		t.Fatalf("seq after rollback = %v, want [1 2]", got)
	}
}

func TestPendingIntentsLimitsAndCursorPagination(t *testing.T) {
	s := openMem(t)
	for i := 0; i < MaxPendingIntentsPage+7; i++ {
		id := fmt.Sprintf("pending-%04d", i)
		c := journalCommitFixture(t, s, id)
		if _, _, err := s.AppendIntent(id, testCommitIntent(id, c)); err != nil {
			t.Fatal(err)
		}
	}
	for _, limit := range []int{0, -1, MaxPendingIntentsPage + 1} {
		if _, err := s.PendingIntents(limit); !IsInvalidLimit(err) {
			t.Fatalf("limit %d error = %T %v", limit, err, err)
		}
	}
	small, err := s.PendingIntents(3)
	if err != nil || len(small) != 3 {
		t.Fatalf("small bounded page len=%d err=%v", len(small), err)
	}
	first, err := s.PendingIntents(MaxPendingIntentsPage)
	if err != nil || len(first) != MaxPendingIntentsPage {
		t.Fatalf("max page len=%d err=%v", len(first), err)
	}
	second, err := s.PendingIntents(MaxPendingIntentsPage, first[len(first)-1].Seq)
	if err != nil || len(second) != 7 {
		t.Fatalf("second page len=%d err=%v", len(second), err)
	}
	seen := map[int64]bool{}
	for _, item := range append(first, second...) {
		if seen[item.Seq] {
			t.Fatalf("overlap at seq %d", item.Seq)
		}
		seen[item.Seq] = true
	}
	if len(seen) != MaxPendingIntentsPage+7 {
		t.Fatalf("pagination covered %d rows", len(seen))
	}
}

func TestJournalPayloadGoldenBytes(t *testing.T) {
	ref := func(label string) hashref.HashRef { return hashref.SumSHA256([]byte(label)) }
	intent := JournalIntent{
		InvocationID: "golden", WorldRef: ref("world"), EntryHash: ref("entry"),
		ObservedHead: ref("observed"), PrevEntryHash: ref("prev"), TransitionFn: ref("fn"),
		TransitionRef: ref("transition"), Interpreter: ref("interpreter"), LogicalTime: 17,
	}
	intentBytes, err := encodeJournalIntent(intent)
	if err != nil {
		t.Fatal(err)
	}
	const wantIntent = "{\"invocationId\":\"golden\",\"worldRef\":\"sha256:486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7\",\"entryHash\":\"sha256:923fe53966c6cd9343e11af776cd4b05be315ea4b200b02e4d5dfb0f929b73bf\",\"observedHead\":\"sha256:604cee807f644af47487bf2bbab442b94212ac5119f36f995f78e9e4694dae8c\",\"prevEntryHash\":\"sha256:84fd9bac333ad79154348296204fa7f8c537a96e08983e5f73b3f5aca8e8edf7\",\"transitionFn\":\"sha256:0f1e18bb4143dc4be22e61ea4deb0491c2bf7018c6504ad631038aed5ca4a0ca\",\"transitionRef\":\"sha256:70dd37c11434d9c571dd83fdd5450e5b0471f7f5ed52943e9409574fea364d33\",\"interpreter\":\"sha256:9666d9e8899447735cb9897b77dbb121754fd4d503609758755bd3fdae3a4b22\",\"logicalTime\":17}\n"
	if string(intentBytes) != wantIntent {
		t.Fatalf("intent golden:\n%q", intentBytes)
	}
	outcomeBytes, err := encodeJournalOutcome(JournalOutcome{"golden", "committed", ref("world"), 18})
	if err != nil {
		t.Fatal(err)
	}
	const wantOutcome = "{\"invocationId\":\"golden\",\"status\":\"committed\",\"resultRef\":\"sha256:486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7\",\"logicalTime\":18}\n"
	if string(outcomeBytes) != wantOutcome {
		t.Fatalf("outcome golden:\n%q", outcomeBytes)
	}
}

func TestPreJournalMigrationPreservesExistingDDL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pre-journal.db")
	oldSchema := schemaSQL[:strings.Index(schemaSQL, "-- Durable ordered index")]
	sum := sha256.Sum256([]byte(oldSchema))
	if got := hex.EncodeToString(sum[:]); got != "43a9c80b4ebbd73dd7f30eb360db4a2a5df7f33466ca5c3993e9ee075e1354b5" {
		t.Fatalf("pre-journal schema source drifted: sha256=%s", got)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(oldSchema); err != nil {
		t.Fatal(err)
	}
	before := tableDDL(t, db)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	after := tableDDL(t, s.db)
	if _, ok := after["journal"]; !ok {
		t.Fatal("journal table absent after writer open")
	}
	delete(after, "journal")
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("existing sqlite_master DDL drifted:\nbefore=%v\nafter=%v", before, after)
	}
}

func tableDDL(t *testing.T, db *sql.DB) map[string]string {
	t.Helper()
	rows, err := db.Query(`SELECT name, sql FROM sqlite_master
		WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var name, ddl string
		if err := rows.Scan(&name, &ddl); err != nil {
			t.Fatal(err)
		}
		got[name] = ddl
	}
	return got
}

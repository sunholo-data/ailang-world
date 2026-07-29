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

func effectIntentFixture(episode string, logicalTime int64) EffectIntent {
	return EffectIntent{
		EpisodeID: episode, Effect: "FS.Write", Scope: "/workspace/out",
		Cost: 7, RequestRef: hashref.SumSHA256([]byte("effect-request")), LogicalTime: logicalTime,
	}
}

func plantEffectIntent(t *testing.T, s *Store, intent EffectIntent) {
	t.Helper()
	payload, err := encodeEffectIntent(intent)
	if err != nil {
		t.Fatal(err)
	}
	object := journalObject(EffectIntentV1, payload)
	if err := s.PutObject(object); err != nil {
		t.Fatal(err)
	}
	var seq int64
	if err := s.db.QueryRow(`SELECT COALESCE(MAX(seq), 0) + 1 FROM journal`).Scan(&seq); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`INSERT INTO journal(seq,kind,invocation_id,object_ref)
		VALUES (?,'intent',?,?)`, seq, intent.InvocationID, object.Hash.String()); err != nil {
		t.Fatal(err)
	}
}

func TestEffectJournalPayloadGoldenBytes(t *testing.T) {
	ref := hashref.SumSHA256([]byte("effect-ref"))
	intentBytes, err := encodeEffectIntent(EffectIntent{
		InvocationID: "effect:episode-1:3", EpisodeID: "episode-1", Ordinal: 3,
		Effect: "FS.Write", Scope: "/workspace/out", Cost: 7, RequestRef: ref, LogicalTime: 91,
	})
	if err != nil {
		t.Fatal(err)
	}
	const wantIntent = "{\"invocationId\":\"effect:episode-1:3\",\"episodeId\":\"episode-1\",\"ordinal\":3,\"effect\":\"FS.Write\",\"scope\":\"/workspace/out\",\"cost\":7,\"requestRef\":\"sha256:546e089e6b9035ee569704416c7d584e208bd789f67fc91bec41444cb021eeae\",\"logicalTime\":91}\n"
	if string(intentBytes) != wantIntent {
		t.Fatalf("effect intent golden:\n%s", intentBytes)
	}
	decodedIntent, err := decodeEffectIntent(intentBytes)
	if err != nil || decodedIntent.InvocationID != "effect:episode-1:3" || decodedIntent.RequestRef != ref {
		t.Fatalf("effect intent round trip = %+v, %v", decodedIntent, err)
	}

	outcomeBytes, err := encodeEffectOutcome(EffectOutcome{
		InvocationID: "effect:episode-1:3", Status: "succeeded", RecordRef: ref, LogicalTime: 92,
	})
	if err != nil {
		t.Fatal(err)
	}
	const wantOutcome = "{\"invocationId\":\"effect:episode-1:3\",\"status\":\"succeeded\",\"recordRef\":\"sha256:546e089e6b9035ee569704416c7d584e208bd789f67fc91bec41444cb021eeae\",\"logicalTime\":92}\n"
	if string(outcomeBytes) != wantOutcome {
		t.Fatalf("effect outcome golden:\n%s", outcomeBytes)
	}
	decodedOutcome, err := decodeEffectOutcome(outcomeBytes)
	if err != nil || decodedOutcome.InvocationID != "effect:episode-1:3" || decodedOutcome.RecordRef != ref {
		t.Fatalf("effect outcome round trip = %+v, %v", decodedOutcome, err)
	}
}

func TestEffectJournalNamespaceDisjointness(t *testing.T) {
	s := openMem(t)
	id := EffectInvocationID("ep-1", 0)
	c := journalCommitFixture(t, s, id)
	if _, _, err := s.AppendIntent(id, testCommitIntent(id, c)); !IsInvocationMismatch(err) {
		t.Fatalf("commit-side effect ID error = %T %v", err, err)
	}
	if _, _, err := s.AppendEffectOutcome("commit-id", EffectOutcome{
		InvocationID: "commit-id", Status: "failed",
		RecordRef: hashref.SumSHA256([]byte("record")), LogicalTime: 1,
	}); !IsInvocationMismatch(err) {
		t.Fatalf("effect-side commit ID error = %T %v", err, err)
	}
}

func TestGetReceiptRejectsEffectNamespaceBeforeDecode(t *testing.T) {
	s := openMem(t)
	id, _, err := s.AppendNextEffectIntent("ep-1", effectIntentFixture("ep-1", 1))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = s.GetReceipt(id)
	var mismatch *InvocationMismatchError
	if !errors.As(err, &mismatch) || mismatch.Field != "InvocationID" {
		t.Fatalf("GetReceipt effect ID = %T %v, want namespace mismatch", err, err)
	}
}

func TestPendingIntentsExcludeRealEffectObjects(t *testing.T) {
	s := openMem(t)
	id := EffectInvocationID("ep-cross", 0)
	plantEffectIntent(t, s, EffectIntent{
		InvocationID: id, EpisodeID: "ep-cross", Ordinal: 0, Effect: "FS.Read",
		Scope: "/workspace/in", Cost: 1, RequestRef: hashref.SumSHA256([]byte("cross-request")),
		LogicalTime: 2,
	})
	pending, err := s.PendingIntents(10)
	if err != nil || len(pending) != 0 {
		t.Fatalf("commit pending contaminated by effect: n=%d err=%v", len(pending), err)
	}
}

func TestAppendNextEffectIntentValidationAndOrdinalDerivation(t *testing.T) {
	s := openMem(t)
	if id, _, err := s.AppendNextEffectIntent("", effectIntentFixture("", 1)); !IsInvocationMismatch(err) || id != "" {
		t.Fatalf("empty episode = (%q,%v), want structured rejection", id, err)
	}
	ep := "episode:with:colon"
	id0, ordinal0, err := s.AppendNextEffectIntent(ep, effectIntentFixture(ep, 10))
	if err != nil || id0 != EffectInvocationID(ep, 0) || ordinal0 != 0 {
		t.Fatalf("fresh mint = (%q,%d,%v)", id0, ordinal0, err)
	}
	if _, _, err := s.AppendEffectOutcome(id0, EffectOutcome{
		InvocationID: id0, Status: "succeeded",
		RecordRef: hashref.SumSHA256([]byte("record-0")), LogicalTime: 11,
	}); err != nil {
		t.Fatal(err)
	}
	id1, ordinal1, err := s.AppendNextEffectIntent(ep, effectIntentFixture(ep, 12))
	if err != nil || id1 != EffectInvocationID(ep, 1) || ordinal1 != 1 {
		t.Fatalf("resumed mint = (%q,%d,%v)", id1, ordinal1, err)
	}
	receipt, ok, err := s.GetEffectReceipt(id1)
	if err != nil || !ok || receipt.State != ReceiptIndeterminate ||
		receipt.EffectIntent == nil || receipt.EffectIntent.LogicalTime != 12 {
		t.Fatalf("minted receipt = %+v, ok=%v err=%v", receipt, ok, err)
	}
}

func TestAppendNextEffectIntentIgnoresAdversarialSuffixes(t *testing.T) {
	s := openMem(t)
	ep := "ep-adversarial"
	for ordinal := int64(0); ordinal < 2; ordinal++ {
		if _, got, err := s.AppendNextEffectIntent(ep, effectIntentFixture(ep, ordinal+1)); err != nil || got != ordinal {
			t.Fatalf("seed ordinal %d = %d, %v", ordinal, got, err)
		}
	}
	for _, id := range []string{"effect:" + ep + ":9x", "effect:" + ep + ":x:0"} {
		plantEffectIntent(t, s, EffectIntent{
			InvocationID: id, EpisodeID: ep, Effect: "FS.Read", Scope: "/x", Cost: 1,
			RequestRef: hashref.SumSHA256([]byte(id)), LogicalTime: 3,
		})
	}
	id, ordinal, err := s.AppendNextEffectIntent(ep, effectIntentFixture(ep, 4))
	if err != nil || ordinal != 2 || id != EffectInvocationID(ep, 2) {
		t.Fatalf("adversarial next = (%q,%d,%v), want ordinal 2", id, ordinal, err)
	}
}

func TestAppendNextEffectIntentOrdinalExhaustion(t *testing.T) {
	s := openMem(t)
	ep := "ep-exhausted"
	id := EffectInvocationID(ep, int64(^uint64(0)>>1))
	plantEffectIntent(t, s, EffectIntent{
		InvocationID: id, EpisodeID: ep, Ordinal: int64(^uint64(0) >> 1),
		Effect: "FS.Read", Scope: "/x", Cost: 1,
		RequestRef: hashref.SumSHA256([]byte("exhausted")), LogicalTime: 1,
	})
	_, _, err := s.AppendNextEffectIntent(ep, effectIntentFixture(ep, 2))
	var exhausted *OrdinalExhaustedError
	if !errors.As(err, &exhausted) || exhausted.EpisodeID != ep {
		t.Fatalf("exhaustion error = %T %v", err, err)
	}
}

func TestAppendNextEffectIntentConcurrentAllocation(t *testing.T) {
	s := openMem(t)
	start := make(chan struct{})
	type result struct {
		id      string
		ordinal int64
		err     error
	}
	results := make(chan result, 2)
	for i := 0; i < 2; i++ {
		logicalTime := int64(i + 1)
		go func() {
			<-start
			id, ordinal, err := s.AppendNextEffectIntent(
				"ep-concurrent", effectIntentFixture("ep-concurrent", logicalTime))
			results <- result{id, ordinal, err}
		}()
	}
	close(start)
	first, second := <-results, <-results
	if first.err != nil || second.err != nil {
		t.Fatalf("both concurrent appends must succeed: first=%v second=%v", first.err, second.err)
	}
	ordinals := map[int64]bool{first.ordinal: true, second.ordinal: true}
	if len(ordinals) != 2 || !ordinals[0] || !ordinals[1] || first.id == second.id {
		t.Fatalf("concurrent allocations = %+v, %+v", first, second)
	}
	var durable int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM journal
		WHERE kind='intent' AND invocation_id >= 'effect:ep-concurrent:'
			AND invocation_id < 'effect:ep-concurrent;'`).Scan(&durable); err != nil {
		t.Fatal(err)
	}
	if durable != 2 {
		t.Fatalf("durable concurrent intents = %d, want 2", durable)
	}
	for _, item := range []result{first, second} {
		receipt, ok, err := s.GetEffectReceipt(item.id)
		if err != nil || !ok || receipt.State != ReceiptIndeterminate {
			t.Fatalf("receipt %q = %+v, ok=%v err=%v", item.id, receipt, ok, err)
		}
	}
}

func TestAppendEffectOutcomeDisciplineAndReceiptWalk(t *testing.T) {
	s := openMem(t)
	missingID := EffectInvocationID("missing", 0)
	outcome := EffectOutcome{
		InvocationID: missingID, Status: "failed",
		RecordRef: hashref.SumSHA256([]byte("missing-record")), LogicalTime: 2,
	}
	if _, _, err := s.AppendEffectOutcome(missingID, outcome); !IsInvocationMismatch(err) {
		t.Fatalf("orphan effect outcome = %T %v", err, err)
	}

	id := EffectInvocationID("ep-walk", 0)
	receipt, ok, err := s.GetEffectReceipt(id)
	if err != nil || ok || receipt.State != ReceiptNotStarted {
		t.Fatalf("not-started = %+v, ok=%v err=%v", receipt, ok, err)
	}
	id, _, err = s.AppendNextEffectIntent("ep-walk", effectIntentFixture("ep-walk", 3))
	if err != nil {
		t.Fatal(err)
	}
	receipt, ok, err = s.GetEffectReceipt(id)
	if err != nil || !ok || receipt.State != ReceiptIndeterminate {
		t.Fatalf("indeterminate = %+v, ok=%v err=%v", receipt, ok, err)
	}
	outcome = EffectOutcome{
		InvocationID: id, Status: "succeeded",
		RecordRef: hashref.SumSHA256([]byte("record")), LogicalTime: 4,
	}
	if _, _, err := s.AppendEffectOutcome(id, outcome); err != nil {
		t.Fatal(err)
	}
	receipt, ok, err = s.GetEffectReceipt(id)
	if err != nil || !ok || receipt.State != ReceiptResolved ||
		receipt.EffectOutcome == nil || receipt.EffectOutcome.Status != "succeeded" {
		t.Fatalf("resolved = %+v, ok=%v err=%v", receipt, ok, err)
	}
	if _, _, err := s.AppendEffectOutcome(id, outcome); !IsDuplicateInvocation(err) {
		t.Fatalf("duplicate outcome = %T %v", err, err)
	}
}

func TestPendingEffectIntentsLimitsPagingAndIsolation(t *testing.T) {
	s := openMem(t)
	c := journalCommitFixture(t, s, "commit-isolation")
	if _, _, err := s.AppendIntent("commit-isolation", testCommitIntent("commit-isolation", c)); err != nil {
		t.Fatal(err)
	}
	var ids []string
	for i := 0; i < 5; i++ {
		id, _, err := s.AppendNextEffectIntent("ep-paging", effectIntentFixture("ep-paging", int64(i+1)))
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	if _, _, err := s.AppendEffectOutcome(ids[1], EffectOutcome{
		InvocationID: ids[1], Status: "failed",
		RecordRef: hashref.SumSHA256([]byte("paging-record")), LogicalTime: 9,
	}); err != nil {
		t.Fatal(err)
	}
	for _, limit := range []int{0, -1, MaxPendingIntentsPage + 1} {
		if _, err := s.PendingEffectIntents(limit); !IsInvalidLimit(err) {
			t.Fatalf("limit %d error = %T %v", limit, err, err)
		}
	}
	first, err := s.PendingEffectIntents(2)
	if err != nil || len(first) != 2 {
		t.Fatalf("first page = %+v, err=%v", first, err)
	}
	second, err := s.PendingEffectIntents(2, first[1].Seq)
	if err != nil || len(second) != 2 {
		t.Fatalf("second page = %+v, err=%v", second, err)
	}
	got := []string{first[0].InvocationID, first[1].InvocationID, second[0].InvocationID, second[1].InvocationID}
	want := []string{ids[0], ids[2], ids[3], ids[4]}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("paged effect IDs = %v, want %v", got, want)
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

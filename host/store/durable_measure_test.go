package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// measureN is the sample size behind the row 23 policy bound table.
const measureN = 50

func quantiles(d []time.Duration) string {
	sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
	p := func(q float64) time.Duration { return d[int(q*float64(len(d)-1))] }
	return fmt.Sprintf("n=%d p50=%v p99=%v max=%v", len(d), p(0.50), p(0.99), d[len(d)-1])
}

// TestMeasureDurableOps is the instrument for the bound table in
// design_docs/planned/w-store-bounded-durable-operations.md. It is skipped
// unless WORLD_BOUND_MEASURE=1; it asserts nothing about speed.
func TestMeasureDurableOps(t *testing.T) {
	if os.Getenv("WORLD_BOUND_MEASURE") != "1" {
		t.Skip("set WORLD_BOUND_MEASURE=1 to measure")
	}
	ctx := context.Background()
	samples := map[string][]time.Duration{}
	timeIt := func(name string, f func() error) {
		start := time.Now()
		if err := f(); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		samples[name] = append(samples[name], time.Since(start))
	}
	dir := t.TempDir()
	for i := 0; i < measureN; i++ {
		path := filepath.Join(dir, fmt.Sprintf("open-%d.db", i))
		timeIt("Open(fresh file, schema)+Close", func() error {
			s, err := Open(path)
			if err != nil {
				return err
			}
			return s.Close()
		})
		timeIt("Open(existing file)+Close", func() error {
			s, err := Open(path)
			if err != nil {
				return err
			}
			return s.Close()
		})
	}
	s := openFileStore(t)
	genesis := seedGenesis(t, s)
	head := genesis
	for i := 0; i < measureN; i++ {
		id := fmt.Sprintf("m-%d", i)
		body := obj("measure-body-"+id, "transition/body")
		entryHash := hashref.SumSHA256([]byte("measure-entry-" + id))
		c := Commit{
			InvocationID: id, ObservedHead: head.Ref, Objects: []Object{body},
			NextWorld: World{Ref: hashref.SumSHA256([]byte("measure-world-" + id)),
				Revision: int64(i + 1), StateRoot: body.Hash, LogHead: entryHash},
			Entry: LogEntry{Header: LogHeader{EntryIndex: int64(i + 1), SemanticsEpoch: 1,
				TransitionFn:  hashref.SumSHA256([]byte("measure-fn")),
				Interpreter:   hashref.SumSHA256([]byte("measure-interp")),
				PrevEntryHash: head.LogHead, WrittenBy: "measure"},
				EntryHash: entryHash, TransitionRef: body.Hash},
		}
		timeIt("GetReceiptContext(absent)", func() error { _, _, err := s.GetReceiptContext(ctx, id); return err })
		timeIt("AppendIntentContext", func() error {
			_, _, err := s.AppendIntentContext(ctx, id, testCommitIntent(id, c))
			return err
		})
		timeIt("CommitContext(invocation)", func() error { return s.CommitContext(ctx, c) })
		timeIt("GetReceiptContext(resolved)", func() error { _, _, err := s.GetReceiptContext(ctx, id); return err })
		timeIt("GetObject", func() error { _, _, err := s.GetObject(ctx, body.Hash); return err })
		timeIt("SelectedHead", func() error { _, _, err := s.SelectedHead(ctx); return err })
		cred := fmt.Sprintf("cred-%d", i)
		if err := s.MintSession(context.Background(), SessionRow{CredentialID: cred, EpisodeID: "ep", GrantsJSON: "[]",
			ExpiresAt: 1 << 40, CreatedAt: 1}); err != nil {
			t.Fatal(err)
		}
		timeIt("ResolveSession", func() error { _, _, err := s.ResolveSession(ctx, cred); return err })
		timeIt("BeginTx+Rollback (uncontended acquisition)", func() error {
			tx, err := s.db.BeginTx(ctx, nil)
			if err != nil {
				return err
			}
			return tx.Rollback()
		})
		head = c.NextWorld
	}
	names := make([]string, 0, len(samples))
	for n := range samples {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		t.Logf("%-45s %s", n, quantiles(samples[n]))
	}
}

// TestMeasureContextFreeSurface times every remaining exported context-free
// Store method (revision 1 inventory). Scans run at their production page
// sizes over full pages: PendingIntents at MaxPendingIntentsPage over 1000
// pending intents, PendingEffectIntents at world-publish's 200 over 200
// effect intents, and both integrity scans at the daemon's 64-row page.
func TestMeasureContextFreeSurface(t *testing.T) {
	if os.Getenv("WORLD_BOUND_MEASURE") != "1" {
		t.Skip("set WORLD_BOUND_MEASURE=1 to measure")
	}
	ctx := context.Background()
	samples := map[string][]time.Duration{}
	timeIt := func(name string, f func() error) {
		start := time.Now()
		if err := f(); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		samples[name] = append(samples[name], time.Since(start))
	}
	s := openFileStore(t)
	genesis := seedGenesis(t, s)
	for i := 0; i < MaxPendingIntentsPage; i++ {
		c := journalCommitFixture(t, s, fmt.Sprintf("p-%d", i))
		if _, _, err := s.AppendIntent(fmt.Sprintf("p-%d", i), testCommitIntent(fmt.Sprintf("p-%d", i), c)); err != nil {
			t.Fatal(err)
		}
	}
	head := genesis
	for i := 0; i < 64; i++ { // one full integrity-scan page of log entries
		body := obj(fmt.Sprintf("scan-body-%d", i), "transition/body")
		entryHash := hashref.SumSHA256([]byte(fmt.Sprintf("scan-entry-%d", i)))
		next := World{Ref: hashref.SumSHA256([]byte(fmt.Sprintf("scan-world-%d", i))),
			Revision: int64(i + 1), StateRoot: body.Hash, LogHead: entryHash}
		if err := s.Commit(Commit{ObservedHead: head.Ref, Objects: []Object{body}, NextWorld: next,
			Entry: LogEntry{Header: LogHeader{EntryIndex: int64(i + 1), SemanticsEpoch: 1,
				TransitionFn: body.Hash, Interpreter: body.Hash, PrevEntryHash: head.LogHead,
				WrittenBy: "measure"}, EntryHash: entryHash, TransitionRef: body.Hash}}); err != nil {
			t.Fatal(err)
		}
		head = next
	}
	for i := 0; i < 200; i++ {
		intent := effectIntentFixture("ep-scan", int64(i))
		if _, _, err := s.AppendNextEffectIntent(context.Background(), "ep-scan", intent); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < measureN; i++ {
		id := fmt.Sprintf("x-%d", i)
		o := obj("surface-"+id, "measure/object")
		timeIt("PutObject", func() error { return s.PutObject(context.Background(), o) })
		w := World{Ref: hashref.SumSHA256([]byte("w-" + id)), Revision: int64(i + 1),
			StateRoot: o.Hash, LogHead: o.Hash}
		timeIt("PutWorld", func() error { return s.PutWorld(context.Background(), w) })
		timeIt("SelectHead", func() error { return s.SelectHead(context.Background(), head.Ref) })
		timeIt("SetRegistryHead", func() error { return s.SetRegistryHead(context.Background(), "measure/reg", o.Hash) })
		o2 := obj("surface-next-"+id, "measure/object")
		if err := s.PutObject(context.Background(), o2); err != nil {
			t.Fatal(err)
		}
		timeIt("CompareAndSetRegistryHead", func() error {
			return s.CompareAndSetRegistryHead(context.Background(), "measure/reg", o.Hash, o2.Hash)
		})
		timeIt("PutVerifyResult", func() error {
			return s.PutVerifyResult(context.Background(), VerifyResult{TransitionFn: o.Hash, Interpreter: genesis.Ref,
				SemanticsEpoch: 1, Verified: true, Detail: "measure"})
		})
		cred := "mint-" + id
		timeIt("MintSession", func() error {
			return s.MintSession(context.Background(), SessionRow{CredentialID: cred, EpisodeID: "ep", GrantsJSON: "[]",
				ExpiresAt: 1 << 40, CreatedAt: 1})
		})
		timeIt("RevokeSession", func() error { return s.RevokeSession(ctx, cred) })
		var effectID string
		timeIt("AppendNextEffectIntent", func() error {
			var err error
			effectID, _, err = s.AppendNextEffectIntent(context.Background(), "ep-m", effectIntentFixture("ep-m", int64(i)))
			return err
		})
		timeIt("AppendEffectOutcome", func() error {
			_, _, err := s.AppendEffectOutcome(context.Background(), effectID, EffectOutcome{InvocationID: effectID,
				Status: "succeeded", RecordRef: o.Hash, LogicalTime: int64(i)})
			return err
		})
		timeIt("GetEffectReceipt", func() error { _, _, err := s.GetEffectReceipt(context.Background(), effectID); return err })
		pid := fmt.Sprintf("p-%d", i)
		timeIt("AppendOutcome", func() error {
			_, _, err := s.AppendOutcome(context.Background(), pid, JournalOutcome{InvocationID: pid, Status: "committed",
				ResultRef: o.Hash, LogicalTime: 42})
			return err
		})
		timeIt("PendingIntents(page 1000)", func() error { _, err := s.PendingIntents(context.Background(), MaxPendingIntentsPage); return err })
		timeIt("PendingEffectIntents(page 200)", func() error { _, err := s.PendingEffectIntents(context.Background(), 200); return err })
		timeIt("ScanUnreadableLog(page 64)", func() error { _, err := s.ScanUnreadableLog(context.Background(), 0, 64); return err })
		timeIt("ScanUnreadableWorlds(page 64)", func() error { _, err := s.ScanUnreadableWorlds(context.Background(), "", 64); return err })
	}
	names := make([]string, 0, len(samples))
	for n := range samples {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		t.Logf("%-45s %s", n, quantiles(samples[n]))
	}
}

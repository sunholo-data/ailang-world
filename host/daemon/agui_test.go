package daemon

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

var _ readStore = (*store.Store)(nil)

func TestAGUIReadStoreSeam(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "seam.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	h := hashref.SumSHA256([]byte("seam"))
	if err := s.Commit(ctx, store.Commit{
		NextWorld: store.World{Ref: h, Revision: 0, StateRoot: h, LogHead: h},
		Entry:     store.LogEntry{Header: store.LogHeader{EntryIndex: 0, SemanticsEpoch: 1, TransitionFn: h, Interpreter: h, PrevEntryHash: h, WrittenBy: "seam"}, EntryHash: h, TransitionRef: h},
	}); err != nil {
		t.Fatal(err)
	}
	var reads readStore = s
	entries, err := reads.LogEntriesAfter(ctx, -1, 100)
	if err != nil || len(entries) != 1 || entries[0].Header.EntryIndex != 0 || entries[0].EntryHash != h {
		t.Fatalf("readStore keyset: %+v err=%v", entries, err)
	}
}

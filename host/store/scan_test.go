package store

import (
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestIntegrityScanLimitGuard(t *testing.T) {
	s := openMem(t)
	for _, limit := range []int{0, -1, MaxIntegrityScanPage + 1} {
		for _, scan := range []func(int) error{
			func(n int) error { _, err := s.ScanUnreadableLog(0, n); return err },
			func(n int) error { _, err := s.ScanUnreadableWorlds("", n); return err },
		} {
			err := scan(limit)
			var invalid *InvalidLimitError
			if !errors.As(err, &invalid) || !IsInvalidLimit(err) || invalid.Limit != limit {
				t.Fatalf("limit %d: got %#v", limit, err)
			}
		}
	}
	if _, err := s.ScanUnreadableLog(0, MaxIntegrityScanPage); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ScanUnreadableWorlds("", MaxIntegrityScanPage); err != nil {
		t.Fatal(err)
	}
}

func TestScanUnreadableLogKeysetResumes(t *testing.T) {
	s := openMem(t)
	ref := hashref.SumSHA256([]byte("scan-log")).String()
	for i := int64(1); i <= 3; i++ {
		entryRef := hashref.SumSHA256([]byte{byte(i)}).String()
		prev := ref
		if i == 2 {
			prev = ""
		}
		if _, err := s.db.Exec(`INSERT INTO log_entries
			(entry_index,entry_hash_ref,semantics_epoch,transition_fn_ref,interpreter_ref,
			 prev_entry_hash_ref,written_by,transition_ref) VALUES(?,?,?,?,?,?,?,?)`,
			i, entryRef, 1, ref, ref, prev, "scan", ref); err != nil {
			t.Fatal(err)
		}
	}
	first, err := s.ScanUnreadableLog(1, 2)
	if err != nil || first.Scanned != 2 || first.Done || first.NextIndex != 3 ||
		len(first.Rows) != 1 || first.Rows[0].Index != 2 || first.Rows[0].Field != "prevEntryHash" {
		t.Fatalf("first page = %+v, err=%v", first, err)
	}
	second, err := s.ScanUnreadableLog(first.NextIndex, 2)
	if err != nil || second.Scanned != 1 || !second.Done || second.NextIndex != 4 {
		t.Fatalf("second page = %+v, err=%v", second, err)
	}
}

func TestScanUnreadableWorldsKeysetStableAcrossEarlierInsert(t *testing.T) {
	s := openMem(t)
	valid := hashref.SumSHA256([]byte("scan-world")).String()
	for _, worldRef := range []string{"b", "d", "f"} {
		if _, err := s.db.Exec(`INSERT INTO worlds(world_ref,revision,state_root,log_head)
			VALUES(?,?,?,?)`, worldRef, 1, valid, valid); err != nil {
			t.Fatal(err)
		}
	}
	first, err := s.ScanUnreadableWorlds("", 2)
	if err != nil || first.NextRef != "d" || first.Done {
		t.Fatalf("first page = %+v, err=%v", first, err)
	}
	// Inserting before the cursor must not shift the next keyset page.
	if _, err := s.db.Exec(`INSERT INTO worlds(world_ref,revision,state_root,log_head)
		VALUES(?,?,?,?)`, "c", 1, valid, valid); err != nil {
		t.Fatal(err)
	}
	second, err := s.ScanUnreadableWorlds(first.NextRef, 2)
	if err != nil || second.Scanned != 1 || !second.Done || second.NextRef != "f" {
		t.Fatalf("second page = %+v, err=%v", second, err)
	}
}

func TestScanUnreadableWorldsFindsPoison(t *testing.T) {
	s := openMem(t)
	valid := hashref.SumSHA256([]byte("scan-world-poison")).String()
	if _, err := s.db.Exec(`INSERT INTO worlds(world_ref,revision,state_root,log_head)
		VALUES(?,?,?,?)`, valid, 1, "", valid); err != nil {
		t.Fatal(err)
	}
	page, err := s.ScanUnreadableWorlds("", 10)
	if err != nil || len(page.Rows) != 1 || page.Rows[0].Ref != valid || page.Rows[0].Field != "stateRoot" {
		t.Fatalf("page = %+v, err=%v", page, err)
	}
}

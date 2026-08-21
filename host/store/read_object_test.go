package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestReadObjectProbeOmitsPayloadAndGuardsMaterialization(t *testing.T) {
	selectList, _, ok := strings.Cut(strings.TrimPrefix(readObjectProbeSQL, "SELECT "), "\nFROM ")
	if !ok {
		t.Fatalf("probe statement has an unrecognized shape: %q", readObjectProbeSQL)
	}
	for _, expression := range strings.Split(selectList, ",") {
		if strings.TrimSpace(expression) == "payload" {
			t.Fatalf("probe SELECT list materializes payload: %q", readObjectProbeSQL)
		}
	}
	if !strings.Contains(selectList, "length(payload)") {
		t.Fatalf("probe SELECT list does not measure payload length: %q", readObjectProbeSQL)
	}

	s, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()
	payload := []byte("payload larger than the bound")
	o := testReadObject(t, s, payload)

	meta, got, err := s.ReadObject(context.Background(), o.Hash, int64(len(payload)-1))
	var tooLarge *ObjectTooLargeError
	if !errors.As(err, &tooLarge) {
		t.Fatalf("ReadObject oversize error = %v; want *ObjectTooLargeError", err)
	}
	if got != nil {
		t.Fatalf("ReadObject oversize returned %d payload bytes; want nil", len(got))
	}
	if tooLarge.Size != int64(len(payload)) || tooLarge.MaxBytes != int64(len(payload)-1) {
		t.Fatalf("ObjectTooLargeError = %+v; want size=%d max=%d", tooLarge, len(payload), len(payload)-1)
	}
	if meta.PayloadLength != int64(len(payload)) || meta.InterfaceHash != o.InterfaceHash || meta.SemanticID != o.SemanticID || meta.Provenance != o.Provenance {
		t.Fatalf("ReadObject metadata = %+v; want metadata and length from probe", meta)
	}

	absent := hashref.SumSHA256([]byte("absent"))
	meta, got, err = s.ReadObject(context.Background(), absent, 1024)
	if err != nil || got != nil || meta != (ObjectMeta{}) {
		t.Fatalf("ReadObject absent = (%+v, %v, %v); want zero, nil, nil", meta, got, err)
	}
}

func TestReadObjectSchedulingHookIsNilInProduction(t *testing.T) {
	if readObjectBetweenStatements != nil {
		t.Fatal("readObjectBetweenStatements is non-nil without a test installation")
	}
}

func TestConcurrentMutationCannotDesyncProbeAndPayload(t *testing.T) {
	t.Run("no-write control", func(t *testing.T) {
		s, _, o := fileReadObjectFixture(t, []byte(strings.Repeat("a", 100)))
		fired := false
		readObjectBetweenStatements = func() { fired = true }
		t.Cleanup(func() { readObjectBetweenStatements = nil })

		meta, payload, err := s.ReadObject(context.Background(), o.Hash, 1024)
		assertSnapshotRead(t, o.Hash, meta, payload, err)
		if !fired {
			t.Fatal("no-write scheduling hook did not fire")
		}
	})

	t.Run("concurrent mutation", func(t *testing.T) {
		s, path, o := fileReadObjectFixture(t, []byte(strings.Repeat("b", 100)))
		writer, err := sql.Open("sqlite", fileDSNWithPragma(path, "busy_timeout(2000)"))
		if err != nil {
			t.Fatalf("open second database handle: %v", err)
		}
		t.Cleanup(func() { writer.Close() })
		writer.SetMaxOpenConns(1)

		fired := false
		writerOutcome := "not-run"
		writerRowsAffected := int64(-1)
		mutated := []byte(strings.Repeat("changed", 50))
		readObjectBetweenStatements = func() {
			fired = true
			res, writeErr := writer.ExecContext(context.Background(),
				`UPDATE objects SET payload = ? WHERE hash_ref = ?`, mutated, o.Hash.String())
			if writeErr != nil {
				writerOutcome = "busy-refused: " + writeErr.Error()
				return
			}
			// RowsAffected is the discriminating observable. "the statement
			// succeeded" is ALSO what a removed statement looks like from the
			// error value alone, so the outcome string must carry something
			// only a real UPDATE can produce.
			affected, affErr := res.RowsAffected()
			if affErr != nil {
				writerOutcome = "committed-after-snapshot: rows unknown: " + affErr.Error()
				return
			}
			writerRowsAffected = affected
			writerOutcome = fmt.Sprintf("committed-after-snapshot: %d row(s)", affected)
		}
		t.Cleanup(func() { readObjectBetweenStatements = nil })

		meta, payload, err := s.ReadObject(context.Background(), o.Hash, 1024)
		if !fired {
			t.Fatal("mutating scheduling hook did not fire; pass would be vacuous")
		}
		// B3: "a green can never come from a writer that never ran." Asserting
		// only that the HOOK fired is not that — a hook whose UPDATE is removed
		// still sets fired, and the arm then passes under the M25 mutant too
		// (measured). Nor is asserting the outcome STRING enough: a removed
		// statement yields a nil error, so "committed-after-snapshot" is written
		// by both a real write and no write at all — the observable's value set
		// is wider than the mechanism's. The discriminator has to be a value
		// only a real UPDATE produces, so this asserts the row count.
		switch {
		case writerOutcome == "not-run":
			t.Fatal("the competing writer never executed; this arm cannot discriminate " +
				"one snapshot from two and would pass under the M25 mutant")
		case strings.HasPrefix(writerOutcome, "busy-refused: "):
			// The rollback-journal outcome: the write was refused at
			// busy_timeout. Real contention, and AC17 accepts it.
		case writerRowsAffected != 1:
			t.Fatalf("the competing writer reported success but touched %d row(s), "+
				"want exactly 1: outcome=%q — a write that changed nothing cannot "+
				"desync anything, so a pass here would be vacuous",
				writerRowsAffected, writerOutcome)
		}
		t.Logf("writer observed outcome: %s", writerOutcome)
		assertSnapshotRead(t, o.Hash, meta, payload, err)
	})
}

func TestBusyTimeoutCachesEffectiveDSNAndDoesNotBlock(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		s, err := Open(t.TempDir() + "/default.db")
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer s.Close()
		if got := s.BusyTimeout(); got != 2*time.Second {
			t.Fatalf("BusyTimeout = %v; want 2s", got)
		}
	})

	t.Run("caller override and occupied connection", func(t *testing.T) {
		path := t.TempDir() + "/override.db"
		dsn := fileDSNWithPragma(path, "busy_timeout(1375)")
		s, err := Open(dsn)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer s.Close()
		if got := s.BusyTimeout(); got != 1375*time.Millisecond {
			t.Fatalf("BusyTimeout = %v; want caller value 1.375s", got)
		}

		decoy, err := s.db.Conn(context.Background())
		if err != nil {
			t.Fatalf("occupy sole connection: %v", err)
		}
		released := make(chan struct{})
		go func() {
			time.Sleep(300 * time.Millisecond)
			decoy.Close()
			close(released)
		}()
		started := time.Now()
		got := s.BusyTimeout()
		elapsed := time.Since(started)
		if elapsed >= 50*time.Millisecond {
			t.Fatalf("BusyTimeout blocked %v behind a 300ms decoy hold", elapsed)
		}
		if got != 1375*time.Millisecond {
			t.Fatalf("BusyTimeout with occupied connection = %v; want 1.375s", got)
		}
		t.Logf("BusyTimeout returned in %v while sole connection was held for 300ms", elapsed)
		<-released
	})
}

func fileReadObjectFixture(t *testing.T, payload []byte) (*Store, string, Object) {
	t.Helper()
	path := t.TempDir() + "/objects.db"
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s, path, testReadObject(t, s, payload)
}

func testReadObject(t *testing.T, s *Store, payload []byte) Object {
	t.Helper()
	o := Object{
		Hash:          hashref.SumSHA256(payload),
		InterfaceHash: hashref.SumSHA256([]byte("interface")),
		SemanticID:    "test/object/v1",
		Provenance:    "read-object-test",
		Payload:       payload,
	}
	if err := s.PutObject(o); err != nil {
		t.Fatalf("PutObject: %v", err)
	}
	return o
}

func assertSnapshotRead(t *testing.T, ref hashref.HashRef, meta ObjectMeta, payload []byte, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("ReadObject: %v", err)
	}
	if int64(len(payload)) != meta.PayloadLength {
		t.Fatalf("payload diverged from probed row: probe=%d bytes, payload=%d bytes; want one snapshot", meta.PayloadLength, len(payload))
	}
	got, err := hashref.Sum(ref.Algo(), payload)
	if err != nil {
		t.Fatalf("hash returned payload: %v", err)
	}
	if got != ref {
		t.Fatalf("payload hash = %s; want requested ref %s", got.String(), ref.String())
	}
}

func fileDSNWithPragma(path, pragma string) string {
	u := url.URL{Scheme: "file", Path: path}
	q := u.Query()
	q.Add("_pragma", pragma)
	u.RawQuery = q.Encode()
	return u.String()
}

// TestBusyTimeoutMatchesTheDriverUnderDuplicatePragmas pins BusyTimeout() against
// what the DRIVER actually applies, not against a comment.
//
// withBusyTimeout returns early when the caller's DSN already carries any
// busy_timeout, so a DSN with TWO of them reaches busyTimeoutFromParams
// unmodified. The pinned driver applies the FIRST; reporting the LAST made the
// accessor under-report the real lock-retry window in one order and over-report
// it in the other — and under-reporting is the unsafe direction for AC18's
// ObjectReadTimeout > BusyTimeout() ordering.
func TestBusyTimeoutMatchesTheDriverUnderDuplicatePragmas(t *testing.T) {
	cases := []struct{ name, first, second string }{
		{"ascending", "busy_timeout(100)", "busy_timeout(7000)"},
		{"descending", "busy_timeout(7000)", "busy_timeout(100)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			params := url.Values{}
			params.Add("_pragma", tc.first)
			params.Add("_pragma", tc.second)
			dsn := fileURI(t.TempDir()+"/dup.db", params)

			db, err := sql.Open("sqlite", dsn)
			if err != nil {
				t.Fatalf("open raw handle: %v", err)
			}
			defer db.Close()
			var appliedMS int64
			if err := db.QueryRow("PRAGMA busy_timeout").Scan(&appliedMS); err != nil {
				t.Fatalf("read back PRAGMA busy_timeout: %v", err)
			}
			applied := time.Duration(appliedMS) * time.Millisecond

			if got := busyTimeoutFromParams(params); got != applied {
				t.Fatalf("busyTimeoutFromParams = %v but the driver applied %v "+
					"(DSN order: %s then %s)", got, applied, tc.first, tc.second)
			}
		})
	}

	// Control: the readback must be able to observe a value we chose, or the
	// agreement above proves nothing.
	params := url.Values{"_pragma": []string{"busy_timeout(3333)"}}
	db, err := sql.Open("sqlite", fileURI(t.TempDir()+"/control.db", params))
	if err != nil {
		t.Fatalf("control open: %v", err)
	}
	defer db.Close()
	var appliedMS int64
	if err := db.QueryRow("PRAGMA busy_timeout").Scan(&appliedMS); err != nil {
		t.Fatalf("control readback: %v", err)
	}
	if appliedMS != 3333 {
		t.Fatalf("CONTROL FAILED: single pragma applied %dms, want 3333 — the "+
			"readback instrument cannot see the DSN, so the arms above are vacuous", appliedMS)
	}
}

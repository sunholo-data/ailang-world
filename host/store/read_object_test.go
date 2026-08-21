package store

import (
	"context"
	"database/sql"
	"errors"
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
		mutated := []byte(strings.Repeat("changed", 50))
		readObjectBetweenStatements = func() {
			fired = true
			_, writeErr := writer.ExecContext(context.Background(),
				`UPDATE objects SET payload = ? WHERE hash_ref = ?`, mutated, o.Hash.String())
			if writeErr != nil {
				writerOutcome = "busy-refused: " + writeErr.Error()
				return
			}
			writerOutcome = "committed-after-snapshot"
		}
		t.Cleanup(func() { readObjectBetweenStatements = nil })

		meta, payload, err := s.ReadObject(context.Background(), o.Hash, 1024)
		if !fired {
			t.Fatal("mutating scheduling hook did not fire; pass would be vacuous")
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

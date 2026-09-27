package daemon

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/workbench"
)

type referenceReadStore struct {
	readStore
	refs                               []store.ObjectReference
	readErr                            error
	missingLog, missingWorld, mismatch bool
	seenCtx                            context.Context
}

func (s *referenceReadStore) ObjectReferences(ctx context.Context, _ hashref.HashRef, _ *store.ObjectReferenceCursor, limit int) ([]store.ObjectReference, error) {
	s.seenCtx = ctx
	if s.readErr != nil {
		return nil, s.readErr
	}
	if len(s.refs) < limit {
		limit = len(s.refs)
	}
	return s.refs[:limit], nil
}
func (s *referenceReadStore) GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error) {
	if s.missingLog && index == 0 {
		return store.LogEntry{}, false, nil
	}
	e, ok, err := s.readStore.GetLogEntry(ctx, index)
	if s.mismatch && ok && index == 0 {
		e.TransitionRef = hashref.SumSHA256([]byte("mismatch"))
	}
	return e, ok, err
}
func (s *referenceReadStore) GetWorld(ctx context.Context, ref hashref.HashRef) (store.World, bool, error) {
	if s.missingWorld && ref == hashref.SumSHA256([]byte("missing-source-world")) {
		return store.World{}, false, nil
	}
	return s.readStore.GetWorld(ctx, ref)
}
func referenceFixture(t *testing.T) (*Daemon, hashref.HashRef, hashref.HashRef) {
	t.Helper()
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "reference-fixture")
	commit := testCommit(genesis, 0, "reference-fixture")
	if err := d.store.Commit(commit); err != nil {
		t.Fatal(err)
	}
	return d, commit.Entry.TransitionRef, commit.NextWorld.Ref
}
func refPage(t *testing.T, d *Daemon, ref hashref.HashRef) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	d.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workbench?object="+ref.String(), nil))
	return rec
}
func TestWorkbenchReferenceFirstPage(t *testing.T) {
	for _, count := range []int{0, 1, 100, 101} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			d, ref, _ := referenceFixture(t)
			rs := &referenceReadStore{readStore: d.store}
			for i := 0; i < count; i++ {
				rs.refs = append(rs.refs, store.ObjectReference{Cursor: store.ObjectReferenceCursor{Kind: store.ReferenceTransitionRef, EntryIndex: 0}})
			}
			d.reads = rs
			rec := refPage(t, d, ref)
			if rec.Code != 200 {
				t.Fatalf("status %d: %s", rec.Code, rec.Body)
			}
			body := rec.Body.String()
			want := count > 100
			if strings.Contains(body, "Showing 100 references; more recorded") != want {
				t.Fatalf("truncation at %d", count)
			}
			if strings.Count(body, "transitionRef: <a") != min(count, 100) {
				t.Fatalf("display count at %d", count)
			}
			if count == 1 {
				link := pageHref(0, 0)
				target := httptest.NewRecorder()
				d.Handler().ServeHTTP(target, httptest.NewRequest(http.MethodGet, "/workbench"+link, nil))
				if target.Code != 200 || !strings.Contains(target.Body.String(), ref.String()) {
					t.Fatalf("link failed: %d", target.Code)
				}
			}
		})
	}
}
func TestWorkbenchMissingReferenceSource(t *testing.T) {
	for _, kind := range []struct {
		name string
		kind store.ReferenceKind
	}{{"log", store.ReferenceTransitionRef}, {"world", store.ReferenceStateRoot}} {
		t.Run(kind.name, func(t *testing.T) {
			d, ref, _ := referenceFixture(t)
			c := store.ObjectReferenceCursor{Kind: kind.kind}
			rs := &referenceReadStore{readStore: d.store}
			if kind.name == "log" {
				rs.missingLog = true
			} else {
				c.WorldRef = hashref.SumSHA256([]byte("missing-source-world"))
				rs.missingWorld = true
			}
			rs.refs = []store.ObjectReference{{Cursor: c}}
			d.reads = rs
			rec := refPage(t, d, ref)
			if rec.Code != 200 || !strings.Contains(rec.Body.String(), "source "+map[bool]string{true: "entry", false: "world"}[kind.name == "log"]+" is no longer stored") {
				t.Fatalf("missing source: %d %s", rec.Code, rec.Body)
			}
			section := rec.Body.String()
			start := strings.Index(section, `<section aria-label="referencedBy">`)
			if start < 0 || strings.Contains(section[start:], `href="/workbench?from=`) || strings.Contains(section[start:], `href="/workbench?world=`) {
				t.Fatal("missing source linked")
			}
		})
	}
}
func TestWorkbenchReferenceMismatch(t *testing.T) {
	d, ref, _ := referenceFixture(t)
	rs := &referenceReadStore{readStore: d.store, mismatch: true, refs: []store.ObjectReference{{Cursor: store.ObjectReferenceCursor{Kind: store.ReferenceTransitionRef}}}}
	d.reads = rs
	rec := refPage(t, d, ref)
	if rec.Code != 500 || strings.Contains(rec.Body.String(), "mismatch") {
		t.Fatalf("status/body: %d %s", rec.Code, rec.Body)
	}
}
func TestWorkbenchReferenceUnavailable(t *testing.T) {
	d, ref, _ := referenceFixture(t)
	d.reads = &referenceReadStore{readStore: d.store, readErr: &store.ReferenceIndexUnavailableError{}}
	rec := refPage(t, d, ref)
	if rec.Code != 503 || !strings.Contains(rec.Body.String(), "ReferenceIndexUnavailable") || strings.Contains(rec.Body.String(), "none recorded") {
		t.Fatalf("status/body: %d %s", rec.Code, rec.Body)
	}
}
func TestWorkbenchReferenceCancellation(t *testing.T) {
	d, ref, _ := referenceFixture(t)
	rs := &referenceReadStore{readStore: d.store, readErr: context.DeadlineExceeded}
	d.reads = rs
	rec := refPage(t, d, ref)
	if rec.Code != 503 || !strings.Contains(rec.Body.String(), "Timeout") {
		t.Fatalf("status/body: %d %s", rec.Code, rec.Body)
	}
	if rs.seenCtx == nil {
		t.Fatal("no context")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(http.MethodGet, "/workbench?object="+ref.String(), nil).WithContext(ctx)
	d.Handler().ServeHTTP(httptest.NewRecorder(), req)
	if !errors.Is(rs.seenCtx.Err(), context.Canceled) {
		t.Fatal("reverse read lost parent cancellation")
	}
}

var _ = math.MaxInt64
var _ = workbench.WorkbenchPageLimit

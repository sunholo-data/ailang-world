package daemon

import (
	"context"
	"errors"
	"fmt"
	"html"
	"math"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

type commitReadStore struct {
	readStore
	indexes       []int64
	membershipErr error
	entryErr      error
	missing       int64
	seenCtx       context.Context
	seenAfter     int64
	seenLimit     int
}

func (s *commitReadStore) ObjectCommits(ctx context.Context, _ hashref.HashRef, after int64, limit int) ([]int64, error) {
	s.seenCtx, s.seenAfter, s.seenLimit = ctx, after, limit
	if s.membershipErr != nil {
		return nil, s.membershipErr
	}
	var out []int64
	for _, n := range s.indexes {
		if n > after {
			out = append(out, n)
			if len(out) == limit {
				break
			}
		}
	}
	return out, nil
}
func (s *commitReadStore) GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error) {
	if s.entryErr != nil && index == s.missing {
		err := s.entryErr
		s.entryErr = nil
		return store.LogEntry{}, false, err
	}
	if index == s.missing {
		return store.LogEntry{}, false, nil
	}
	if index > 0 {
		e, ok, err := s.readStore.GetLogEntry(ctx, 0)
		e.Header.EntryIndex = index
		return e, ok, err
	}
	return s.readStore.GetLogEntry(ctx, index)
}
func commitFixture(t *testing.T) (*Daemon, hashref.HashRef, *commitReadStore) {
	t.Helper()
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "commit-page")
	c := testCommit(genesis, 0, "commit-page")
	if err := d.store.Commit(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	s := &commitReadStore{readStore: d.store, missing: -99}
	d.reads = s
	return d, c.Entry.TransitionRef, s
}
func commitGET(d *Daemon, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	d.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}
func TestWorkbenchCommittedByEdges(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		d, ref, s := commitFixture(t)
		s.indexes = nil
		rec := commitGET(d, "/workbench?object="+ref.String())
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), "no commit carried this object: it was stored outside any commit (PutObject or journal). Entries that only reference it are listed under referencedBy.") {
			t.Fatalf("%d %s", rec.Code, rec.Body)
		}
	})
	t.Run("available", func(t *testing.T) {
		d, ref, s := commitFixture(t)
		s.indexes = []int64{0}
		rec := commitGET(d, "/workbench?object="+ref.String())
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), `committedBy: <a href="/workbench?from=0&amp;entry=0"`) || s.seenAfter != -1 || s.seenLimit != 101 || s.seenCtx == nil {
			t.Fatalf("%d %s after=%d limit=%d", rec.Code, rec.Body, s.seenAfter, s.seenLimit)
		}
		link := commitGET(d, "/workbench?from=0&entry=0")
		if link.Code != 200 || !strings.Contains(link.Body.String(), "selected entry 0") {
			t.Fatalf("link %d", link.Code)
		}
	})
	t.Run("missing-entry", func(t *testing.T) {
		d, ref, s := commitFixture(t)
		s.indexes = []int64{7}
		s.missing = 7
		rec := commitGET(d, "/workbench?object="+ref.String())
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), "UNAVAILABLE: log entry 7 is not stored") || strings.Contains(rec.Body.String(), `committedBy: <a`) {
			t.Fatalf("%d %s", rec.Code, rec.Body)
		}
	})
	t.Run("read-error", func(t *testing.T) {
		d, ref, s := commitFixture(t)
		s.indexes = []int64{7}
		s.missing = 7
		s.entryErr = errors.New("private commit read detail")
		rec := commitGET(d, "/workbench?object="+ref.String())
		if rec.Code != 500 || !strings.Contains(rec.Body.String(), ">Internal<") || strings.Contains(rec.Body.String(), "private commit read detail") {
			t.Fatalf("%d %s", rec.Code, rec.Body)
		}
	})
	t.Run("membership-error", func(t *testing.T) {
		d, ref, s := commitFixture(t)
		s.membershipErr = errors.New("private table detail")
		rec := commitGET(d, "/workbench?object="+ref.String())
		if rec.Code != 500 || strings.Contains(rec.Body.String(), "private table detail") {
			t.Fatalf("%d %s", rec.Code, rec.Body)
		}
	})
	t.Run("deadline", func(t *testing.T) {
		d, ref, s := commitFixture(t)
		s.membershipErr = context.DeadlineExceeded
		rec := commitGET(d, "/workbench?object="+ref.String())
		if rec.Code != 503 || !strings.Contains(rec.Body.String(), ">Timeout<") {
			t.Fatalf("%d %s", rec.Code, rec.Body)
		}
	})
	t.Run("max-index", func(t *testing.T) {
		d, ref, s := commitFixture(t)
		s.indexes = []int64{math.MaxInt64}
		rec := commitGET(d, "/workbench?object="+ref.String())
		want := fmt.Sprintf(`from=%d&amp;entry=%d`, int64(math.MaxInt64-WorkbenchPageLimit), int64(math.MaxInt64))
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), want) {
			t.Fatalf("%d %s", rec.Code, rec.Body)
		}
	})
}
func TestWorkbenchCommitPaging(t *testing.T) {
	for _, n := range []int{0, 100, 101} {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			d, ref, s := commitFixture(t)
			for i := 0; i < n; i++ {
				s.indexes = append(s.indexes, int64(i))
			}
			rec := commitGET(d, "/workbench?object="+ref.String())
			body := rec.Body.String()
			if rec.Code != 200 || (strings.Contains(body, "Showing 100 commits; more recorded") != (n > 100)) || strings.Count(body, `committedBy: <a`) != min(n, 100) {
				t.Fatalf("n=%d code=%d", n, rec.Code)
			}
			if n == 101 {
				if !strings.Contains(body, "commitsAfter=99") {
					t.Fatal("next cursor missing")
				}
				next := commitGET(d, "/workbench?object="+ref.String()+"&commitsAfter=99")
				if next.Code != 200 || s.seenAfter != 99 || strings.Count(next.Body.String(), `committedBy: <a`) != 1 || !strings.Contains(next.Body.String(), "entry 100</a>") {
					t.Fatalf("continuation %d after %d body %s", next.Code, s.seenAfter, next.Body)
				}
				end := commitGET(d, "/workbench?object="+ref.String()+"&commitsAfter=100")
				if end.Code != 200 || !strings.Contains(end.Body.String(), "no further commits carried this object") {
					t.Fatalf("end %d", end.Code)
				}
			}
		})
	}
	t.Run("continuation", func(t *testing.T) {
		d, ref, fake := commitFixture(t)
		for i := 0; i < 101; i++ {
			fake.indexes = append(fake.indexes, int64(i))
		}
		first := commitGET(d, "/workbench?object="+ref.String())
		if first.Code != 200 || !strings.Contains(first.Body.String(), "commitsAfter=99") {
			t.Fatalf("first page %d", first.Code)
		}
		next := commitGET(d, "/workbench?object="+ref.String()+"&commitsAfter=99")
		if next.Code != 200 || fake.seenAfter != 99 || strings.Count(next.Body.String(), `committedBy: <a`) != 1 || !strings.Contains(next.Body.String(), "entry 100</a>") {
			t.Fatalf("next page %d after %d", next.Code, fake.seenAfter)
		}
	})
	t.Run("cursor-coexistence", func(t *testing.T) {
		d, ref, s := commitFixture(t)
		for i := 0; i < 102; i++ {
			s.indexes = append(s.indexes, int64(i))
		}
		refs := &referenceReadStore{readStore: s}
		for i := 0; i < 102; i++ {
			refs.refs = append(refs.refs, store.ObjectReference{Cursor: store.ObjectReferenceCursor{Kind: store.ReferenceTransitionRef, EntryIndex: int64(i)}})
		}
		d.reads = refs
		first := commitGET(d, "/workbench?object="+ref.String()+"&payload=1&commitsAfter=0&refsAfter="+encodeReferenceCursor(store.ObjectReferenceCursor{Kind: store.ReferenceTransitionRef, EntryIndex: 0}))
		if first.Code != 200 || !strings.Contains(first.Body.String(), "payload=1") {
			t.Fatalf("%d", first.Code)
		}
		if !strings.Contains(first.Body.String(), "commitsAfter=0") || !strings.Contains(first.Body.String(), "refsAfter=") {
			t.Fatal("cursor link missing")
		}
		commitNext := regexp.MustCompile(`href="([^"]*commitsAfter=100[^"]*)"[^>]*>next commits</a>`).FindStringSubmatch(first.Body.String())
		if len(commitNext) != 2 || !strings.Contains(commitNext[1], "refsAfter=") || !strings.Contains(commitNext[1], "payload=1") {
			t.Fatal("commit link lost other cursor or payload")
		}
		refsNext := regexp.MustCompile(`href="([^"]*refsAfter=[^"]*)"[^>]*>next references</a>`).FindStringSubmatch(first.Body.String())
		if len(refsNext) != 2 || !strings.Contains(refsNext[1], "commitsAfter=0") || !strings.Contains(refsNext[1], "payload=1") {
			t.Fatal("reference link lost commit cursor or payload")
		}
		for _, link := range []string{commitNext[1], refsNext[1]} {
			target := commitGET(d, html.UnescapeString(link))
			if target.Code != 200 {
				t.Fatalf("link %s status %d", link, target.Code)
			}
		}
	})
}
func TestWorkbenchCommitGrammar(t *testing.T) {
	d, ref, _ := commitFixture(t)
	for _, q := range []string{"commitsAfter=0", "world=" + ref.String() + "&commitsAfter=0", "object=" + ref.String() + "&from=0&commitsAfter=0", "object=" + ref.String() + "&entry=0&commitsAfter=0", "object=" + ref.String() + "&commitsAfter=0&commitsAfter=1", "object=" + ref.String() + "&commitsAfter=", "object=" + ref.String() + "&commitsAfter=-1", "object=" + ref.String() + "&commitsAfter=bad", "object=" + ref.String() + "&commitsAfter=9223372036854775808"} {
		rec := commitGET(d, "/workbench?"+q)
		if rec.Code != 400 {
			t.Errorf("%s: %d", q, rec.Code)
		}
	}
	for _, q := range []string{"object=" + ref.String() + "&commitsAfter=0", "object=" + ref.String() + "&commitsAfter=0&payload=1", "object=" + ref.String() + "&commitsAfter=0&refsAfter=" + encodeReferenceCursor(store.ObjectReferenceCursor{Kind: store.ReferenceTransitionRef, EntryIndex: 0})} {
		rec := commitGET(d, "/workbench?"+q)
		if rec.Code != 200 {
			t.Errorf("%s: %d %s", q, rec.Code, rec.Body)
		}
	}
}
func TestWorkbenchCommitWalk(t *testing.T) {
	start := time.Now()
	d := newHandlerDaemon(t)
	current := seedGenesisEmbedded(t, d, "commit-walk")
	state := workbenchTestObject("shared-state", "incident/row104/state", hashref.SumSHA256([]byte("state-iface")))
	interpreter := workbenchTestObject("shared-interpreter", "incident/row104/interpreter", hashref.SumSHA256([]byte("interpreter-iface")))
	transitions := make([]hashref.HashRef, 3)
	for i := int64(0); i < 3; i++ {
		c := testCommit(current, i, fmt.Sprintf("row104-%d", i))
		c.Entry.Header.Interpreter = interpreter.Hash
		c.NextWorld.StateRoot = state.Hash
		c.Objects = append(c.Objects, state)
		if i == 0 {
			c.Objects[0].SemanticID = "incident/row104/transition"
			c.Objects = append(c.Objects, interpreter)
		}
		transitions[i] = c.Entry.TransitionRef
		if err := d.store.Commit(context.Background(), c); err != nil {
			t.Fatal(err)
		}
		current = c.NextWorld
	}
	located, err := d.store.ObjectsBySemanticID(context.Background(), "incident/row104/transition", "", 10)
	if err != nil || len(located) != 1 || located[0].Hash != transitions[0] {
		t.Fatalf("semantic locate: %v %+v", err, located)
	}
	questions := []struct {
		question string
		ref      hashref.HashRef
		want     []int64
	}{
		{"Which entry carried the transition body?", located[0].Hash, []int64{0}},
		{"Which entries carried the identical state object?", state.Hash, []int64{0, 1, 2}},
		{"Which entry carried the interpreter object?", interpreter.Hash, []int64{0}},
	}
	re := regexp.MustCompile(`<p>committedBy: <a href="([^"]+)"[^>]*>entry ([0-9]+)</a></p>`)
	for _, q := range questions {
		object := commitGET(d, "/workbench?object="+q.ref.String())
		if object.Code != 200 {
			t.Fatalf("object %d", object.Code)
		}
		matches := re.FindAllStringSubmatch(object.Body.String(), -1)
		if len(matches) != len(q.want) {
			t.Fatalf("%s: edges %d", q.question, len(matches))
		}
		var statuses []int
		for j, m := range matches {
			idx, _ := strconv.ParseInt(m[2], 10, 64)
			if idx != q.want[j] {
				t.Fatalf("%s: index %d", q.question, idx)
			}
			entry := commitGET(d, html.UnescapeString(m[1]))
			statuses = append(statuses, entry.Code)
			if entry.Code != 200 || !strings.Contains(entry.Body.String(), fmt.Sprintf("selected entry %d", idx)) || !strings.Contains(entry.Body.String(), "transitionRef: <a href=\"/workbench?object="+transitions[idx].String()) {
				t.Fatalf("%s: entry %d status %d", q.question, idx, entry.Code)
			}
			payload := commitGET(d, "/workbench?object="+transitions[idx].String()+"&payload=1")
			statuses = append(statuses, payload.Code)
			if payload.Code != 200 || !strings.Contains(payload.Body.String(), "raw bytes") {
				t.Fatalf("payload %d", payload.Code)
			}
		}
		t.Logf("question=%s answer=%v refs=object %s transitions %v statuses=object %d followed %v elapsed=%s", q.question, q.want, q.ref, transitions, object.Code, statuses, time.Since(start))
	}
}

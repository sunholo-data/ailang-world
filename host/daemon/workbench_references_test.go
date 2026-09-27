package daemon

import (
	"context"
	"errors"
	"fmt"
	"html"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

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
	dense                              bool
}

func (s *referenceReadStore) ObjectReferences(ctx context.Context, _ hashref.HashRef, after *store.ObjectReferenceCursor, limit int) ([]store.ObjectReference, error) {
	s.seenCtx = ctx
	if s.readErr != nil {
		return nil, s.readErr
	}
	refs := s.refs
	if after != nil {
		i := 0
		for i < len(refs) && refs[i].Cursor.EntryIndex <= after.EntryIndex {
			i++
		}
		refs = refs[i:]
	}
	if len(refs) < limit {
		limit = len(refs)
	}
	return refs[:limit], nil
}
func (s *referenceReadStore) GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error) {
	if s.missingLog && index == 0 {
		return store.LogEntry{}, false, nil
	}
	if s.dense && index > 0 {
		e, ok, err := s.readStore.GetLogEntry(ctx, 0)
		e.Header.EntryIndex = index
		return e, ok, err
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

func TestWorkbenchReferenceGrammar(t *testing.T) {
	d, ref, _ := referenceFixture(t)
	good := encodeReferenceCursor(store.ObjectReferenceCursor{Kind: store.ReferenceTransitionRef, EntryIndex: 0})
	for _, q := range []string{"refsAfter=" + good, "world=" + ref.String() + "&refsAfter=" + good, "object=" + ref.String() + "&from=0&refsAfter=" + good, "object=" + ref.String() + "&refsAfter=bad", "object=" + ref.String() + "&refsAfter=" + good + "&refsAfter=" + good} {
		rec := requestRecorder(t, d, http.MethodGet, "/workbench?"+q, nil)
		if rec.Code != 400 {
			t.Fatalf("%s: %d", q, rec.Code)
		}
	}
	rec := requestRecorder(t, d, http.MethodGet, "/workbench?object="+ref.String()+"&refsAfter="+good, nil)
	if rec.Code != 200 {
		t.Fatalf("valid continuation: %d", rec.Code)
	}
}
func TestWorkbenchReferencePaging(t *testing.T) {
	for _, n := range []int{100, 101, 501} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			d, ref, _ := referenceFixture(t)
			rs := &referenceReadStore{readStore: d.store, dense: true}
			for i := 0; i < n; i++ {
				rs.refs = append(rs.refs, store.ObjectReference{Cursor: store.ObjectReferenceCursor{Kind: store.ReferenceTransitionRef, EntryIndex: int64(i)}})
			}
			d.reads = rs
			path := "/workbench?object=" + ref.String() + "&payload=1"
			seen := 0
			for {
				rec := requestRecorder(t, d, http.MethodGet, path, nil)
				if rec.Code != 200 {
					t.Fatalf("page status %d: %s", rec.Code, rec.Body)
				}
				body := rec.Body.String()
				seen += strings.Count(body, "transitionRef: <a")
				marker := `>next references</a>`
				at := strings.Index(body, marker)
				if at < 0 {
					break
				}
				prefix := body[:at]
				start := strings.LastIndex(prefix, `href="`)
				if start < 0 {
					t.Fatal("no next href")
				}
				escaped := strings.SplitN(prefix[start+6:], `"`, 2)[0]
				escaped = html.UnescapeString(escaped)
				u, err := url.Parse(escaped)
				if err != nil {
					t.Fatal(err)
				}
				if u.Query().Get("object") != ref.String() || u.Query().Get("payload") != "1" {
					t.Fatal("lost object/payload state")
				}
				cursor, err := decodeReferenceCursor(u.Query().Get("refsAfter"))
				if err != nil {
					t.Fatal(err)
				}
				if cursor.EntryIndex != int64(seen-1) {
					t.Fatalf("no-skipped-lookahead cursor %d want %d", cursor.EntryIndex, seen-1)
				}
				path = escaped
			}
			if seen != n {
				t.Fatalf("walk saw %d of %d", seen, n)
			}
		})
	}
	t.Run("no-skipped-lookahead", func(t *testing.T) {
		d, ref, _ := referenceFixture(t)
		rs := &referenceReadStore{readStore: d.store, dense: true}
		for i := 0; i < 101; i++ {
			rs.refs = append(rs.refs, store.ObjectReference{Cursor: store.ObjectReferenceCursor{Kind: store.ReferenceTransitionRef, EntryIndex: int64(i)}})
		}
		d.reads = rs
		rec := refPage(t, d, ref)
		body := rec.Body.String()
		at := strings.Index(body, `>next references</a>`)
		if at < 0 {
			t.Fatal("missing next")
		}
		prefix := body[:at]
		start := strings.LastIndex(prefix, `href="`)
		link := html.UnescapeString(strings.SplitN(prefix[start+6:], `"`, 2)[0])
		u, _ := url.Parse(link)
		cursor, err := decodeReferenceCursor(u.Query().Get("refsAfter"))
		if err != nil || cursor.EntryIndex != 99 {
			t.Fatalf("cursor %+v %v", cursor, err)
		}
		second := requestRecorder(t, d, http.MethodGet, link, nil)
		if second.Code != 200 || !strings.Contains(second.Body.String(), "entry 100") {
			t.Fatal("lookahead skipped")
		}
	})
}

func TestWorkbenchReferenceWalk(t *testing.T) {
	start := time.Now()
	const question = "Which entry references incident/row103/transition, and what transition body does that entry name?"
	d := newHandlerDaemon(t)
	genesis := seedGenesisEmbedded(t, d, "row103-walk")
	state := workbenchTestObject("shared-state", "incident/row103/state", hashref.SumSHA256([]byte("state-iface")))
	interpreter := workbenchTestObject("shared-interpreter", "incident/row103/interpreter", hashref.SumSHA256([]byte("interpreter-iface")))
	current := genesis
	var transition hashref.HashRef
	var transitionPayload []byte
	var worlds []hashref.HashRef
	for i := int64(0); i < 3; i++ {
		c := testCommit(current, i, fmt.Sprintf("row103-%d", i))
		c.Entry.Header.Interpreter = interpreter.Hash
		c.NextWorld.StateRoot = state.Hash
		if i == 0 {
			c.Objects[0].SemanticID = "incident/row103/transition"
			c.Objects = append(c.Objects, state, interpreter)
			transition = c.Entry.TransitionRef
			transitionPayload = append([]byte(nil), c.Objects[0].Payload...)
		}
		if err := d.store.Commit(c); err != nil {
			t.Fatal(err)
		}
		worlds = append(worlds, c.NextWorld.Ref)
		current = c.NextWorld
	}
	found, err := d.store.ObjectsBySemanticID(context.Background(), "incident/row103/transition", "", 10)
	if err != nil || len(found) != 1 || found[0].Hash != transition {
		t.Fatalf("semantic lookup: %v %+v", err, found)
	}
	get := func(path string) (int, string) {
		rec := httptest.NewRecorder()
		d.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec.Code, rec.Body.String()
	}
	objectStatus, objectBody := get("/workbench?object=" + transition.String())
	if objectStatus != 200 {
		t.Fatalf("object status %d", objectStatus)
	}
	re := regexp.MustCompile(`<p>transitionRef: <a href="([^"]+)"`)
	match := re.FindStringSubmatch(objectBody)
	if len(match) != 2 {
		t.Fatal("missing labelled backlink")
	}
	entryLink := html.UnescapeString(match[1])
	entryStatus, entryBody := get(entryLink)
	if entryStatus != 200 || !strings.Contains(entryBody, `<article aria-label="selected entry">`) || !strings.Contains(entryBody, `transitionRef: <a href="/workbench?object=`+transition.String()) {
		t.Fatalf("entry status/body %d %s", entryStatus, entryBody)
	}
	objectLink := regexp.MustCompile(`<p>transitionRef: <a href="([^"]+)"`).FindStringSubmatch(entryBody)
	if len(objectLink) != 2 {
		t.Fatal("selected entry has no transition object link")
	}
	target := html.UnescapeString(objectLink[1])
	if target != `/workbench?object=`+transition.String() {
		t.Fatalf("unexpected target %s", target)
	}
	payloadStatus, payloadBody := get(target + "&payload=1")
	if payloadStatus != 200 || !strings.Contains(payloadBody, string(transitionPayload)) {
		t.Fatalf("payload status/body %d", payloadStatus)
	}
	interpreterStatus, interpreterBody := get("/workbench?object=" + interpreter.Hash.String())
	if interpreterStatus != 200 || strings.Count(interpreterBody, "interpreter: <a") != 3 {
		t.Fatalf("interpreter status/refs %d", interpreterStatus)
	}
	worldStatus, stateBody := get("/workbench?object=" + state.Hash.String())
	if worldStatus != 200 || strings.Count(stateBody, "stateRoot: <a") != 3 {
		t.Fatalf("state status/refs %d", worldStatus)
	}
	for _, ref := range worlds {
		if !strings.Contains(stateBody, ref.String()) {
			t.Fatalf("missing world %s", ref)
		}
		status, body := get("/workbench?world=" + ref.String())
		if status != 200 || !strings.Contains(body, ref.String()) {
			t.Fatalf("world link %s: %d", ref, status)
		}
	}
	t.Logf("question=%s", question)
	t.Logf("answer=entry 0 names transition body %s with payload %q", transition, transitionPayload)
	t.Logf("refs=object %s; entry 0; interpreter %s (3 entries); state %s (worlds %v)", transition, interpreter.Hash, state.Hash, worlds)
	t.Logf("statuses=object %d, entry %d, payload %d, interpreter %d, state %d", objectStatus, entryStatus, payloadStatus, interpreterStatus, worldStatus)
	t.Logf("elapsed=%s", time.Since(start))
	t.Logf("remaining stops=committedBy: %s; registry, journal and object-interface inbound references outside scope", objectCommittedByMissing)
}

func TestWorkbenchReferenceMaxIndexLink(t *testing.T) {
	d, ref, _ := referenceFixture(t)
	rs := &referenceReadStore{readStore: d.store, dense: true, refs: []store.ObjectReference{{Cursor: store.ObjectReferenceCursor{Kind: store.ReferenceTransitionRef, EntryIndex: math.MaxInt64}}}}
	d.reads = rs
	rec := refPage(t, d, ref)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	want := fmt.Sprintf(`href="/workbench?from=%d&amp;entry=%d"`, int64(math.MaxInt64-workbench.WorkbenchPageLimit), int64(math.MaxInt64))
	if !strings.Contains(rec.Body.String(), want) {
		t.Fatalf("missing safe link %s", want)
	}
	link := pageHref(math.MaxInt64-workbench.WorkbenchPageLimit, math.MaxInt64)
	target := requestRecorder(t, d, http.MethodGet, "/workbench"+link, nil)
	if target.Code != 200 || !strings.Contains(target.Body.String(), `selected entry 9223372036854775807`) {
		t.Fatalf("link status %d", target.Code)
	}
}

package transitionreg

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

var (
	testFn          = hashref.MustParse("sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	testInterpreter = hashref.MustParse("sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
)

func validDescriptor() Descriptor {
	return Descriptor{ID: "tools.echo", TransitionFn: testFn, Interpreter: testInterpreter, SemanticsEpoch: 1, InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`), Access: EffectRequirement{Effect: "read", Scope: "world", Cost: 1}, DeclaredEffects: []EffectRequirement{{Effect: "read", Scope: "world", Cost: 1}}, Title: "Echo", Description: "test"}
}

type fakeObjectStore struct {
	mu          sync.Mutex
	head        hashref.HashRef
	hasHead     bool
	headErr     error
	object      store.Object
	hasObject   bool
	objectErr   error
	headReads   int
	objectReads int
	putErr      error
	casErr      error
}

func (f *fakeObjectStore) GetRegistryHead(string) (hashref.HashRef, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.headReads++
	return f.head, f.hasHead, f.headErr
}
func (f *fakeObjectStore) GetObject(hashref.HashRef) (store.Object, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.objectReads++
	return f.object, f.hasObject, f.objectErr
}
// clone returns an independent fake seeded from f's configuration. It exists
// because `g := *f` copies the embedded sync.Mutex — `go vet` reports copylocks
// for that, and `go test`'s default vet subset does not include the copylocks
// analyzer, so the shape is invisible to `verify_go.sh`. Counters deliberately
// start at zero: a clone is a fresh observation, never a continuation of f's.
func (f *fakeObjectStore) clone() *fakeObjectStore {
	return &fakeObjectStore{
		head: f.head, hasHead: f.hasHead, headErr: f.headErr,
		object: f.object, hasObject: f.hasObject, objectErr: f.objectErr,
		putErr: f.putErr, casErr: f.casErr,
	}
}

func (f *fakeObjectStore) PutObject(store.Object) error { return f.putErr }
func (f *fakeObjectStore) CompareAndSetRegistryHead(string, hashref.HashRef, hashref.HashRef) error {
	return f.casErr
}

func storedRevision(t *testing.T, r Revision) store.Object {
	t.Helper()
	payload, err := EncodeRevision(r)
	if err != nil {
		t.Fatal(err)
	}
	return store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: InterfaceHashV1, SemanticID: SemanticIDV1, Payload: payload}
}

func fakeWithRevision(t *testing.T, r Revision) *fakeObjectStore {
	o := storedRevision(t, r)
	return &fakeObjectStore{head: o.Hash, hasHead: true, object: o, hasObject: true}
}

func TestReadSnapshotReadsHeadOnce(t *testing.T) {
	f := fakeWithRevision(t, validRevision(validDescriptor()))
	r := NewReader(f)
	first, err := r.ReadSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	second, err := r.ReadSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if first.Head != second.Head || f.headReads != 2 {
		t.Fatalf("head reads = %d, want one per call; heads %q/%q", f.headReads, first.Head, second.Head)
	}
	if f.objectReads != 1 {
		t.Fatalf("object reads = %d, want 1 across cache hit", f.objectReads)
	}
}

func TestSnapshotIsEagerAndCopyIsolated(t *testing.T) {
	d := validDescriptor()
	d.InputSchema = []byte(`{"in":1}`)
	d.OutputSchema = []byte(`{"out":2}`)
	f := fakeWithRevision(t, validRevision(d))
	r := NewReader(f)
	s, err := r.ReadSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// Mutating construction inputs and store bytes cannot affect the parsed snapshot.
	d.InputSchema[2] = 'X'
	f.object.Payload[0] = '['
	list := s.List()
	list[0].ID = "changed"
	list[0].InputSchema[2] = 'Y'
	lookup, ok := s.Lookup("tools.echo")
	if !ok {
		t.Fatal("lookup missed frozen descriptor")
	}
	lookup.OutputSchema[2] = 'Z'
	lookup.DeclaredEffects[0].Effect = "changed"

	again := s.List()
	got, ok := s.Lookup("tools.echo")
	if !ok || again[0].ID != "tools.echo" || string(got.InputSchema) != `{"in":1}` || string(got.OutputSchema) != `{"out":2}` || got.DeclaredEffects[0].Effect != "read" {
		t.Fatalf("snapshot alias escaped: list=%+v lookup=%+v ok=%v", again, got, ok)
	}
	// Cache returns a fresh deep copy without touching the now-corrupt store bytes.
	s.entries[0].InputSchema[2] = 'Q'
	cached, err := r.ReadSnapshot(context.Background())
	if err != nil || string(cached.List()[0].InputSchema) != `{"in":1}` {
		t.Fatalf("cache copy changed: snapshot=%+v err=%v", cached.List(), err)
	}
	cached.entries[0].InputSchema[2] = 'R'
	third, err := r.ReadSnapshot(context.Background())
	if err != nil || string(third.List()[0].InputSchema) != `{"in":1}` {
		t.Fatalf("cache result aliases cache: snapshot=%+v err=%v", third.List(), err)
	}
}

func TestReadSnapshotRefusals(t *testing.T) {
	injected := errors.New("injected head read failure")
	injectedObject := errors.New("injected object read failure")
	base := fakeWithRevision(t, validRevision(validDescriptor()))
	tests := []struct {
		name string
		want string
		make func() *fakeObjectStore
	}{
		{"absent_head", "head is absent", func() *fakeObjectStore { return &fakeObjectStore{} }},
		{"injected_read_error", "read transition registry head: injected head read failure", func() *fakeObjectStore { return &fakeObjectStore{headErr: injected} }},
		{"injected_object_read_error", "read transition registry object: injected object read failure", func() *fakeObjectStore {
			f := base.clone()
			f.objectErr = injectedObject
			return f
		}},
		{"object_absent", "object \"" + base.head.String() + "\" is absent", func() *fakeObjectStore { f := base.clone(); f.hasObject = false; return f }},
		{"corrupt_object_payload", "object hash mismatch", func() *fakeObjectStore {
			f := base.clone()
			f.object.Payload = append([]byte(nil), f.object.Payload...)
			f.object.Payload[0] = '['
			return f
		}},
		{"wrong_semantic_id", "semantic ID \"wrong/semantic\" is not \"world/transition-registry/v1\"", func() *fakeObjectStore { f := base.clone(); f.object.SemanticID = "wrong/semantic"; return f }},
		{"wrong_interface_hash", "interface hash \"sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\" is not", func() *fakeObjectStore {
			f := base.clone()
			f.object.InterfaceHash = hashref.MustParse("sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")
			return f
		}},
		{"revision_raw_over_limit", "raw JSON is 16777217 bytes; limit is 16777216", func() *fakeObjectStore {
			f := base.clone()
			f.object.Payload = bytes.Repeat([]byte(" "), maxRevisionRaw+1)
			f.head = hashref.SumSHA256(f.object.Payload)
			f.object.Hash = f.head
			return f
		}},
		// Decision 3's parent/revision chain rules. Both branches shipped with NO
		// coverage: neutered together with `if false && …` the mutant built and the
		// WHOLE package stayed green, so a corrupted or tampered object with a broken
		// revision chain was silently accepted as a valid Snapshot.
		{"revision_zero_with_parent", "revision 0 must have no parent", func() *fakeObjectStore {
			return fakeWithRevision(t, Revision{SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1, Revision: 0, Parent: testFn, Entries: []Descriptor{validDescriptor()}})
		}},
		{"revision_after_one_without_parent", "revision after 1 must have a parent", func() *fakeObjectStore {
			return fakeWithRevision(t, Revision{SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1, Revision: 5, Entries: []Descriptor{validDescriptor()}})
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewReader(tc.make()).ReadSnapshot(context.Background())
			if err == nil {
				t.Fatal("invalid snapshot was accepted")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("refused by wrong branch: got %q, want substring %q", err, tc.want)
			}
		})
	}
}

func openTransitionStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func seedRevision(t *testing.T, s *store.Store, r Revision) hashref.HashRef {
	t.Helper()
	o := storedRevision(t, r)
	if err := s.PutObject(o); err != nil {
		t.Fatal(err)
	}
	if err := s.CompareAndSetRegistryHead(store.TransitionRegistryV1, hashref.HashRef{}, o.Hash); err != nil {
		t.Fatal(err)
	}
	return o.Hash
}

func TestPublishCASConflictPreservesWinner(t *testing.T) {
	s := openTransitionStore(t)
	current := validRevision(validDescriptor())
	expected := seedRevision(t, s, current)
	dA := validDescriptor()
	dA.Title = "winner"
	dB := validDescriptor()
	dB.Title = "loser"
	nextA, err := BuildNext(current, []Change{{ID: dA.ID, Descriptor: &dA}})
	if err != nil {
		t.Fatal(err)
	}
	nextB, err := BuildNext(current, []Change{{ID: dB.ID, Descriptor: &dB}})
	if err != nil {
		t.Fatal(err)
	}
	p := NewReader(s)
	winner, err := p.Publish(context.Background(), expected, nextA)
	if err != nil {
		t.Fatal(err)
	}
	orphan, err := p.Publish(context.Background(), expected, nextB)
	if err == nil || !store.IsRegistryCASConflict(err) {
		t.Fatalf("stale publication error = %v, want typed CAS conflict", err)
	}
	head, ok, err := s.GetRegistryHead(store.TransitionRegistryV1)
	if err != nil || !ok || head != winner {
		t.Fatalf("winner was not preserved: head=%q ok=%v err=%v", head, ok, err)
	}
	loserBytes, _ := EncodeRevision(nextB)
	loserRef := hashref.SumSHA256(loserBytes)
	if orphan != (hashref.HashRef{}) {
		t.Fatalf("failed publish returned ref %q", orphan)
	}
	if _, ok, err := s.GetObject(loserRef); err != nil || !ok {
		t.Fatalf("CAS orphan was not preserved: ok=%v err=%v", ok, err)
	}
}

func TestConcurrentPublishHasOneWinner(t *testing.T) {
	s := openTransitionStore(t)
	p := NewReader(s)
	next := validRevision(validDescriptor())
	const racers = 8
	start := make(chan struct{})
	var ready sync.WaitGroup
	var done sync.WaitGroup
	ready.Add(racers)
	done.Add(racers)
	errs := make(chan error, racers)
	for i := 0; i < racers; i++ {
		go func() {
			defer done.Done()
			ready.Done()
			<-start
			_, err := p.Publish(context.Background(), hashref.HashRef{}, next)
			errs <- err
		}()
	}
	ready.Wait()
	close(start)
	done.Wait()
	close(errs)
	winners, conflicts := 0, 0
	for err := range errs {
		switch {
		case err == nil:
			winners++
		case store.IsRegistryCASConflict(err):
			conflicts++
		default:
			t.Fatalf("unexpected racer error: %v", err)
		}
	}
	if winners != 1 || conflicts != racers-1 {
		t.Fatalf("winners=%d conflicts=%d, want 1/%d", winners, conflicts, racers-1)
	}
}

func TestStableIDByteOrder(t *testing.T) {
	if CompareID("a", "ab") >= 0 {
		t.Fatal("prefix did not sort shorter first")
	}
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			got := CompareID(string([]byte{byte(a)}), string([]byte{byte(b)}))
			if (a < b && got >= 0) || (a == b && got != 0) || (a > b && got <= 0) {
				t.Fatalf("unsigned byte order failed for %d/%d: %d", a, b, got)
			}
		}
	}
	current := validRevision()
	da, db := validDescriptor(), validDescriptor()
	da.ID, db.ID = "ab", "a"
	next, err := BuildNext(current, []Change{{ID: da.ID, Descriptor: &da}, {ID: db.ID, Descriptor: &db}})
	if err != nil {
		t.Fatal(err)
	}
	if got := []string{next.Entries[0].ID, next.Entries[1].ID}; got[0] != "a" || got[1] != "ab" {
		t.Fatalf("BuildNext order = %v", got)
	}
}

func TestPublishRefusals(t *testing.T) {
	current := validRevision(validDescriptor())
	currentObject := storedRevision(t, current)
	baseNext, err := BuildNext(current, nil)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		want string
		make func() (*fakeObjectStore, hashref.HashRef, Revision)
	}{
		{"revision_not_n_plus_1", "revision 3 is not expected 2", func() (*fakeObjectStore, hashref.HashRef, Revision) {
			n := baseNext
			n.Revision = 3
			return &fakeObjectStore{object: currentObject, hasObject: true}, currentObject.Hash, n
		}},
		{"parent_not_captured_head", "parent \"\" is not captured head", func() (*fakeObjectStore, hashref.HashRef, Revision) {
			n := baseNext
			n.Parent = hashref.HashRef{}
			return &fakeObjectStore{object: currentObject, hasObject: true}, currentObject.Hash, n
		}},
		{"duplicate_id", "entries are not strictly ordered by ID", func() (*fakeObjectStore, hashref.HashRef, Revision) {
			n := baseNext
			n.Entries = append(n.Entries, cloneDescriptor(n.Entries[0]))
			return &fakeObjectStore{object: currentObject, hasObject: true}, currentObject.Hash, n
		}},
		{"entries_over_1024", "entries exceeds 1024", func() (*fakeObjectStore, hashref.HashRef, Revision) {
			n := baseNext
			n.Entries = make([]Descriptor, maxEntries+1)
			return &fakeObjectStore{object: currentObject, hasObject: true}, currentObject.Hash, n
		}},
		{"injected_put_error", "put object: injected put failure", func() (*fakeObjectStore, hashref.HashRef, Revision) {
			return &fakeObjectStore{object: currentObject, hasObject: true, putErr: errors.New("injected put failure")}, currentObject.Hash, baseNext
		}},
		{"injected_cas_error", "compare and set head: injected CAS failure", func() (*fakeObjectStore, hashref.HashRef, Revision) {
			return &fakeObjectStore{object: currentObject, hasObject: true, casErr: errors.New("injected CAS failure")}, currentObject.Hash, baseNext
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, expected, next := tc.make()
			_, err := NewReader(f).Publish(context.Background(), expected, next)
			if err == nil {
				t.Fatal("invalid publication was accepted")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("refused by wrong branch: got %q, want substring %q", err, tc.want)
			}
		})
	}
}
func validRevision(entries ...Descriptor) Revision {
	return Revision{SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1, Revision: 1, Entries: entries}
}

func TestCodecGoldenRoundTrip(t *testing.T) {
	const wantInterface = "sha256:743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac"
	if InterfaceHashV1.String() != wantInterface {
		t.Fatalf("interface hash = %q, want literal %q", InterfaceHashV1, wantInterface)
	}

	const emptyGolden = `{"entries":[],"interfaceHash":"sha256:743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac","parent":"","revision":1,"semanticID":"world/transition-registry/v1"}`
	empty, err := EncodeRevision(validRevision())
	if err != nil {
		t.Fatal(err)
	}
	if string(empty) != emptyGolden {
		t.Fatalf("empty golden bytes:\n got %s\nwant %s", empty, emptyGolden)
	}

	d := validDescriptor()
	d.InputSchema, _ = canonicalSchema([]byte(`{"z":1.0,"a":"<&é"}`))
	d.OutputSchema, _ = canonicalSchema([]byte(`{"n":-0,"array":[3,2,1]}`))
	const descriptorGolden = `{"entries":[{"access":{"cost":1,"effect":"read","scope":"world"},"declaredEffects":[{"cost":1,"effect":"read","scope":"world"}],"description":"test","id":"tools.echo","inputSchema":{"a":"<&é","z":1},"interpreter":"sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","outputSchema":{"array":[3,2,1],"n":0},"semanticsEpoch":1,"title":"Echo","transitionFn":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}],"interfaceHash":"sha256:743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac","parent":"","revision":1,"semanticID":"world/transition-registry/v1"}`
	got, err := EncodeRevision(validRevision(d))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != descriptorGolden {
		t.Fatalf("descriptor golden bytes:\n got %s\nwant %s", got, descriptorGolden)
	}
	round, err := DecodeRevision(got)
	if err != nil {
		t.Fatal(err)
	}
	again, err := EncodeRevision(round)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, again) {
		t.Fatal("Encode -> Decode -> Encode changed bytes")
	}

	numbers, err := canonicalSchema([]byte(`{"a":1,"b":1.0,"c":1e0,"d":-0}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(numbers) != `{"a":1,"b":1,"c":1,"d":0}` {
		t.Fatalf("number normalization = %s", numbers)
	}
	nfc, _ := canonicalSchema([]byte(`{"s":"é"}`))
	nfd, _ := canonicalSchema([]byte("{\"s\":\"e\u0301\"}"))
	if bytes.Equal(nfc, nfd) {
		t.Fatal("NFC and NFD strings were normalized together")
	}
}

func TestDescriptorIdentityAndContentUpdate(t *testing.T) {
	d1 := validDescriptor()
	r1 := validRevision(d1)
	b1, err := EncodeRevision(r1)
	if err != nil {
		t.Fatal(err)
	}
	h1 := hashref.SumSHA256(b1)
	d2 := d1
	d2.TransitionFn = hashref.SumSHA256([]byte("updated transition source"))
	r2 := validRevision(d2)
	r2.Revision = 2
	r2.Parent = h1
	b2, err := EncodeRevision(r2)
	if err != nil {
		t.Fatal(err)
	}
	h2 := hashref.SumSHA256(b2)
	if d1.ID != d2.ID || h1 == h2 {
		t.Fatalf("stable ID/content update invariant failed: ids=%q/%q hashes=%q/%q", d1.ID, d2.ID, h1, h2)
	}
	objects := map[hashref.HashRef][]byte{h1: append([]byte(nil), b1...), h2: append([]byte(nil), b2...)}
	if !bytes.Equal(objects[h1], b1) {
		t.Fatal("old immutable revision is no longer addressable")
	}
	if bytes.Equal(objects[h1], objects[h2]) {
		t.Fatal("content update did not create a new revision object")
	}
}

func TestDescriptorValidationRefusals(t *testing.T) {
	tests := []struct {
		name string
		want string
		run  func() error
	}{
		{"id_grammar", "does not match stable ID grammar", func() error { d := validDescriptor(); d.ID = "Bad ID"; return d.Validate() }},
		{"id_too_long", "length must be 1..128 bytes", func() error {
			d := validDescriptor()
			d.ID = strings.Join([]string{strings.Repeat("a", 32), strings.Repeat("b", 32), strings.Repeat("c", 32), strings.Repeat("d", 32)}, "/")
			return d.Validate()
		}},
		{"segment_too_long", "segment length must be 1..32 bytes", func() error { d := validDescriptor(); d.ID = strings.Repeat("a", 33); return d.Validate() }},
		{"zero_transition_fn", "transitionFn is zero", func() error { d := validDescriptor(); d.TransitionFn = hashref.HashRef{}; return d.Validate() }},
		{"zero_interpreter", "interpreter is zero", func() error { d := validDescriptor(); d.Interpreter = hashref.HashRef{}; return d.Validate() }},
		{"negative_semantics_epoch", "semanticsEpoch is negative", func() error { d := validDescriptor(); d.SemanticsEpoch = -1; return d.Validate() }},
		{"negative_cost", "cost is negative", func() error { d := validDescriptor(); d.Access.Cost = -1; return d.Validate() }},
		{"schema_not_an_object", "schema root must be an object", func() error { d := validDescriptor(); d.InputSchema = []byte(`[]`); return d.Validate() }},
		{"schema_raw_over_262144", "raw JSON is 262146 bytes; limit is 262144", func() error {
			raw := append(bytes.Repeat([]byte(" "), maxSchemaRaw), []byte(`{}`)...)
			_, err := canonicalSchema(raw)
			return err
		}},
		{"schema_canonical_over_65536", "canonical schema is 65544 bytes; limit is 65536", func() error {
			d := validDescriptor()
			d.InputSchema = []byte(`{"x":"` + strings.Repeat("a", maxSchemaCanonical) + `"}`)
			return d.Validate()
		}},
		{"duplicate_schema_key_nested", "duplicate object member \"a\"", func() error {
			d := validDescriptor()
			d.InputSchema = []byte(`{"x":{"a":1,"a":2}}`)
			return d.Validate()
		}},
		{"lone_surrogate", "JSON string contains an escaped surrogate", func() error { d := validDescriptor(); d.InputSchema = []byte(`{"x":"\ud800"}`); return d.Validate() }},
		{"invalid_utf8", "JSON is not valid UTF-8", func() error {
			d := validDescriptor()
			d.InputSchema = []byte{'{', '"', 'x', '"', ':', '"', 0xff, '"', '}'}
			if utf8.Valid(d.InputSchema) {
				t.Fatal("fixture unexpectedly valid")
			}
			return d.Validate()
		}},
		{"number_coefficient_overflow", "number coefficient has 1025 digits; limit is 1024", func() error {
			d := validDescriptor()
			d.InputSchema = []byte(`{"n":` + strings.Repeat("1", maxNumberDigits+1) + `}`)
			return d.Validate()
		}},
		{"unknown_revision_key", "unknown key \"extra\"", func() error {
			_, err := DecodeRevision([]byte(`{"entries":[],"extra":0,"interfaceHash":"` + InterfaceHashV1.String() + `","parent":"","revision":1,"semanticID":"` + SemanticIDV1 + `"}`))
			return err
		}},
		{"missing_descriptor_key", "missing key \"id\"", func() error {
			_, err := DecodeRevision([]byte(`{"entries":[{}],"interfaceHash":"` + InterfaceHashV1.String() + `","parent":"","revision":1,"semanticID":"` + SemanticIDV1 + `"}`))
			return err
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			if err == nil {
				t.Fatal("invalid value was accepted")
			}
			// Pin the BRANCH, not merely the refusal: DecodeRevision's canonical
			// re-encode check (codec.go) refuses these inputs too, so a message-agnostic
			// assertion stays green when the named guard is neutered. Measured: with
			// the escaped-surrogate and unknown-key guards disabled, this test passed.
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("refused by the wrong branch: got %q, want it to contain %q", err, tc.want)
			}
		})
	}
}

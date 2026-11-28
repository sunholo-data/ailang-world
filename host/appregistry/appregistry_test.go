package appregistry

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

var t0 = time.Date(2026, 11, 28, 18, 0, 0, 0, time.UTC)

func TestInterfaceHashesArePreimages(t *testing.T) {
	if got := hashref.SumSHA256([]byte(revisionPreimage)); got != RevisionInterfaceHash {
		t.Fatalf("revision interface hash %s, preimage gives %s", RevisionInterfaceHash, got)
	}
	if got := hashref.SumSHA256([]byte(manifestPreimage)); got != ManifestInterfaceHash {
		t.Fatalf("manifest interface hash %s, preimage gives %s", ManifestInterfaceHash, got)
	}
}

func exampleJSON(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "docs", "examples", "apps", "docparse-convert.json"))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func ref(s string) hashref.HashRef { return hashref.SumSHA256([]byte(s)) }

func TestCodecIsCanonicalAndClosed(t *testing.T) {
	r, err := BuildPublish(nil, hashref.HashRef{}, "docparse-convert", ref("m1"), []string{"b", "a"}, t0)
	if err != nil {
		t.Fatal(err)
	}
	if r.Revision != 1 || r.Parent != "" || r.Apps[0].Version != 1 || strings.Join(r.Apps[0].Transitions, ",") != "a,b" {
		t.Fatalf("first publish = %+v", r)
	}
	enc, err := Encode(r)
	if err != nil {
		t.Fatal(err)
	}
	back, err := Decode(enc)
	if err != nil || back.Apps[0].Manifest != r.Apps[0].Manifest {
		t.Fatalf("round trip: %v %+v", err, back)
	}
	for name, payload := range map[string][]byte{
		"whitespace":  bytes.Replace(enc, []byte(`"revision":1`), []byte(`"revision": 1`), 1),
		"unknown key": bytes.Replace(enc, []byte(`{"schema"`), []byte(`{"owner":"x","schema"`), 1),
		"trailing":    append(append([]byte{}, enc...), []byte("{}")...),
		"no newline":  bytes.TrimSuffix(enc, []byte("\n")),
	} {
		if bytes.Equal(payload, enc) {
			t.Fatalf("%s: case did not change the encoding", name)
		}
		if _, err := Decode(payload); err == nil {
			t.Errorf("%s: Decode accepted it", name)
		}
	}
}

func TestValidateRefusals(t *testing.T) {
	good := func() Revision {
		return Revision{Schema: RegistryName, Revision: 2, Parent: ref("p").String(), Apps: []Entry{
			{Name: "a-app", Version: 1, Manifest: ref("a").String(), Transitions: []string{}, PublishedAt: "2026-11-28T18:00:00Z"},
			{Name: "b-app", Version: 3, Manifest: ref("b").String(), Transitions: []string{"t1", "t2"}, PublishedAt: "2026-11-28T18:00:00Z"},
		}}
	}
	if err := good().Validate(); err != nil {
		t.Fatalf("the good revision is invalid: %v", err)
	}
	cases := map[string]func(r *Revision){
		"schema":              func(r *Revision) { r.Schema = "world/transition-registry/v1" },
		"revision zero":       func(r *Revision) { r.Revision = 0 },
		"revision 1 parent":   func(r *Revision) { r.Revision = 1 },
		"later without one":   func(r *Revision) { r.Parent = "" },
		"bad parent":          func(r *Revision) { r.Parent = "not-a-ref" },
		"nil apps":            func(r *Revision) { r.Apps = nil },
		"unsorted":            func(r *Revision) { r.Apps[0], r.Apps[1] = r.Apps[1], r.Apps[0] },
		"duplicate":           func(r *Revision) { r.Apps[1].Name = "a-app" },
		"bad name":            func(r *Revision) { r.Apps[0].Name = "A_app" },
		"zero version":        func(r *Revision) { r.Apps[0].Version = 0 },
		"bad manifest":        func(r *Revision) { r.Apps[0].Manifest = "" },
		"nil transitions":     func(r *Revision) { r.Apps[0].Transitions = nil },
		"unsorted transition": func(r *Revision) { r.Apps[1].Transitions = []string{"t2", "t1"} },
		"empty transition":    func(r *Revision) { r.Apps[1].Transitions = []string{""} },
		"local time":          func(r *Revision) { r.Apps[0].PublishedAt = "2026-11-28T19:00:00+01:00" },
		"not a time":          func(r *Revision) { r.Apps[0].PublishedAt = "yesterday" },
	}
	for name, edit := range cases {
		r := good()
		edit(&r)
		if err := r.Validate(); err == nil {
			t.Errorf("%s: Validate accepted it", name)
		}
	}
}

func TestBuildPublishAndRetire(t *testing.T) {
	head1 := ref("head1")
	r1, err := BuildPublish(nil, hashref.HashRef{}, "docparse-convert", ref("m1"), nil, t0)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := BuildPublish(&r1, head1, "docparse-convert", ref("m2"), nil, t0.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if r2.Revision != 2 || r2.Parent != head1.String() || len(r2.Apps) != 1 || r2.Apps[0].Version != 2 || r2.Apps[0].Manifest != ref("m2").String() {
		t.Fatalf("second version = %+v", r2)
	}
	if _, err := BuildPublish(&r2, ref("head2"), "docparse-convert", ref("m2"), nil, t0); err == nil || !strings.Contains(err.Error(), "already pins") {
		t.Fatalf("republishing the same manifest = %v, want a refusal", err)
	}
	r3, err := BuildPublish(&r2, ref("head2"), "a-first", ref("other"), nil, t0)
	if err != nil {
		t.Fatal(err)
	}
	if len(r3.Apps) != 2 || r3.Apps[0].Name != "a-first" || r3.Apps[0].Version != 1 || r3.Apps[1].Version != 2 {
		t.Fatalf("a second app = %+v", r3)
	}
	if r1.Apps[0].Version != 1 || len(r2.Apps) != 1 {
		t.Fatal("building a revision mutated its input")
	}
	r4, err := BuildRetire(&r3, ref("head3"), "docparse-convert")
	if err != nil || len(r4.Apps) != 1 || r4.Apps[0].Name != "a-first" || r4.Revision != 4 {
		t.Fatalf("retire = %+v %v", r4, err)
	}
	if _, err := BuildRetire(&r4, ref("head4"), "docparse-convert"); err == nil {
		t.Error("retiring an absent app succeeded")
	}
	if _, err := BuildRetire(nil, hashref.HashRef{}, "x"); err == nil {
		t.Error("retiring from an empty registry succeeded")
	}
	if _, err := BuildPublish(nil, hashref.HashRef{}, "se-tools", ref("m"), nil, t0); err == nil {
		t.Error("publishing the reserved name se-tools succeeded")
	}
}

func openStore(t *testing.T) (*store.Store, context.Context) {
	t.Helper()
	s, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return s, ctx
}

func TestPublishAgainstTheStore(t *testing.T) {
	s, ctx := openStore(t)
	if snap, err := Read(ctx, s); err != nil || snap.Present {
		t.Fatalf("an empty registry reads as %+v %v", snap, err)
	}

	obj, m, err := ManifestObject(exampleJSON(t))
	if err != nil {
		t.Fatal(err)
	}
	r1, err := BuildPublish(nil, hashref.HashRef{}, m.Name, obj.Hash, nil, t0)
	if err != nil {
		t.Fatal(err)
	}
	head1, err := Publish(ctx, s, hashref.HashRef{}, r1, &obj)
	if err != nil {
		t.Fatal(err)
	}
	snap, err := Read(ctx, s)
	if err != nil || !snap.Present || snap.Head != head1 || snap.Revision.Revision != 1 {
		t.Fatalf("after genesis: %+v %v", snap, err)
	}
	e, ok := snap.Revision.Lookup("docparse-convert")
	if !ok || e.Version != 1 || e.Manifest != obj.Hash.String() {
		t.Fatalf("entry = %+v %v", e, ok)
	}
	pinned, raw, err := ReadManifest(ctx, s, obj.Hash)
	if err != nil || pinned.Name != "docparse-convert" || !bytes.Equal(raw, exampleJSON(t)) {
		t.Fatalf("ReadManifest = %v, exact bytes %v", err, bytes.Equal(raw, exampleJSON(t)))
	}

	// A second version, published over the captured head.
	edited := bytes.Replace(exampleJSON(t), []byte(`"budget": 2000`), []byte(`"budget": 1000`), 1)
	obj2, _, err := ManifestObject(edited)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := BuildPublish(&snap.Revision, snap.Head, "docparse-convert", obj2.Hash, nil, t0.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	// A rival publisher plans against the same head and wins first.
	rival, err := BuildRetire(&snap.Revision, snap.Head, "docparse-convert")
	if err != nil {
		t.Fatal(err)
	}
	headRival, err := Publish(ctx, s, snap.Head, rival, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Publish(ctx, s, snap.Head, r2, &obj2); err == nil || !store.IsRegistryCASConflict(err) {
		t.Fatalf("the losing publish = %v, want a registry CAS conflict", err)
	}
	if now, _ := Read(ctx, s); now.Head != headRival || len(now.Revision.Apps) != 0 {
		t.Fatalf("the lost race moved the head: %+v", now)
	}

	// Revision rules are checked against the captured head before anything is written.
	cur, _ := Read(ctx, s)
	wrong, _ := BuildPublish(&cur.Revision, cur.Head, "docparse-convert", obj2.Hash, nil, t0)
	wrong.Revision = 9
	if _, err := Publish(ctx, s, cur.Head, wrong, &obj2); err == nil || !strings.Contains(err.Error(), "captured head") {
		t.Fatalf("a skipped revision = %v", err)
	}
	if _, err := Publish(ctx, s, hashref.HashRef{}, r1, &obj); err == nil {
		t.Fatal("a genesis over an existing head succeeded")
	}
	unpinned, _ := BuildPublish(&cur.Revision, cur.Head, "docparse-convert", obj2.Hash, nil, t0)
	if _, err := Publish(ctx, s, cur.Head, unpinned, &obj); err == nil || !strings.Contains(err.Error(), "no entry pins") {
		t.Fatalf("a manifest no entry pins = %v", err)
	}
	if after, _ := Read(ctx, s); after.Head != headRival {
		t.Fatal("a refused publish moved the head")
	}

	// Objects of the wrong kind are refused on read.
	if _, _, err := ReadManifest(ctx, s, headRival); err == nil {
		t.Error("ReadManifest accepted a revision object")
	}
	if err := s.CompareAndSetRegistryHead(ctx, RegistryName, headRival, obj.Hash); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(ctx, s); err == nil || !strings.Contains(err.Error(), "not an app registry revision") {
		t.Errorf("Read over a manifest object = %v", err)
	}
}

func TestManifestObjectRefusesInvalid(t *testing.T) {
	bad := bytes.Replace(exampleJSON(t), []byte(`"Workspace.Read"`), []byte(`"Model.Infer"`), 1)
	if _, _, err := ManifestObject(bad); err == nil || !strings.Contains(err.Error(), "breaks the law") {
		t.Fatalf("ManifestObject(Model.Infer) = %v", err)
	}
	if _, _, err := ManifestObject([]byte("{")); err == nil {
		t.Fatal("ManifestObject accepted broken JSON")
	}
}

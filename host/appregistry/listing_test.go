package appregistry

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/appmanifest"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

func TestListMergesRegistryAndCatalog(t *testing.T) {
	s, ctx := openStore(t)
	example := exampleJSON(t)
	dir := t.TempDir()
	write := func(name string, b []byte) {
		if err := os.WriteFile(filepath.Join(dir, name+".json"), b, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	// Publish docparse-convert and a second app, sync-notes, whose local copy will differ.
	obj, _, err := ManifestObject(example)
	if err != nil {
		t.Fatal(err)
	}
	r1, _ := BuildPublish(nil, hashref.HashRef{}, "docparse-convert", obj.Hash, nil, t0)
	h1, err := Publish(ctx, s, hashref.HashRef{}, r1, &obj)
	if err != nil {
		t.Fatal(err)
	}
	notes := bytes.Replace(example, []byte(`"docparse-convert"`), []byte(`"sync-notes"`), 1)
	obj2, _, err := ManifestObject(notes)
	if err != nil {
		t.Fatal(err)
	}
	r2, _ := BuildPublish(&r1, h1, "sync-notes", obj2.Hash, nil, t0)
	if _, err := Publish(ctx, s, h1, r2, &obj2); err != nil {
		t.Fatal(err)
	}

	write("docparse-convert", example)                                                            // same bytes
	write("sync-notes", bytes.Replace(notes, []byte(`"budget": 2000`), []byte(`"budget": 9`), 1)) // a draft
	write("draft-only", bytes.Replace(example, []byte(`"docparse-convert"`), []byte(`"draft-only"`), 1))
	local, err := appmanifest.ReadCatalog(dir)
	if err != nil {
		t.Fatal(err)
	}
	snap, err := Read(ctx, s)
	if err != nil {
		t.Fatal(err)
	}
	rows := List(ctx, s, snap, local)
	got := map[string]Listing{}
	for _, r := range rows {
		got[r.Name] = r
	}
	if len(rows) != 3 || rows[0].Name != "docparse-convert" || rows[1].Name != "draft-only" || rows[2].Name != "sync-notes" {
		t.Fatalf("rows = %+v", rows)
	}
	if d := got["docparse-convert"]; d.State != StatePublished || d.Version != 1 || !d.Valid() || !d.Launchable() || d.Local == nil {
		t.Errorf("same bytes = %+v (the local digest and the pinned ref must compare equal)", d)
	}
	if n := got["sync-notes"]; n.State != StateLocalDiffers || !n.Launchable() || n.Manifest == nil || n.Manifest.Preset.Grants[0].Budget != 2000 {
		t.Errorf("a differing draft = %+v; the pinned manifest (budget 2000) is what launches", n)
	}
	if l := got["draft-only"]; l.State != StateLocal || l.Launchable() || l.Version != 0 || !l.Valid() {
		t.Errorf("local only = %+v; a local draft never launches", l)
	}

	// With no registry at all, every catalog entry is local.
	for _, r := range List(ctx, s, Snapshot{}, local) {
		if r.State != StateLocal || r.Launchable() {
			t.Errorf("no registry: %+v", r)
		}
	}
}

func TestListReportsAMissingPinnedManifest(t *testing.T) {
	s, ctx := openStore(t)
	ghost := ref("never stored")
	r1, _ := BuildPublish(nil, hashref.HashRef{}, "docparse-convert", ghost, nil, t0)
	// Publish without the manifest object, so the entry names an object the store lacks.
	if _, err := Publish(ctx, s, hashref.HashRef{}, r1, nil); err != nil {
		t.Fatal(err)
	}
	snap, _ := Read(ctx, s)
	rows := List(ctx, s, snap, nil)
	if len(rows) != 1 || rows[0].Error == "" || rows[0].Valid() || rows[0].Launchable() {
		t.Fatalf("a missing pinned manifest = %+v, want an error row that cannot launch", rows)
	}
}

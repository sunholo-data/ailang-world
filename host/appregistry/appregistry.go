// Package appregistry is World's app registry (queue row 274 M1b,
// design_docs/planned/w-app-registry.md): a registry head world/app-registry/v1 over
// immutable revision objects. Each entry pins an app's exact app.json bytes as a
// world/app-manifest/v1 object, with a registry-owned version.
//
// Revisions follow host/transitionreg's rules: revision 1 has no parent, every later
// revision is its parent's + 1, the parent is the head the publisher captured, and the
// head moves only by compare-and-set. Rollback is roll-forward: publish an older
// manifest again as a new version. This package holds no authority: the attended
// `world-publish apps` verb (M1b-2) is the only writer.
package appregistry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/sunholo-data/ailang-world/host/appmanifest"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const (
	// RegistryName is the head's name and the revision objects' semantic ID.
	RegistryName = "world/app-registry/v1"
	// ManifestSemanticID is the semantic ID of a pinned app.json object.
	ManifestSemanticID = "world/app-manifest/v1"

	// The frozen interface preimages (ASCII, no newline). Each interface hash is the
	// sha256 of its preimage; TestInterfaceHashesArePreimages re-derives both.
	revisionPreimage = "world/app-registry/v1|APP-REG-JSON-1|revision:{apps,parent,revision,schema}|app:{manifest,name,published_at,transitions,version}|name:app-name-1"
	manifestPreimage = "world/app-manifest/v1|exact-app-json-bytes|schema:world/app/v1"

	maxApps        = 256
	maxTransitions = 64
	maxRevisionRaw = 1 << 20
)

var (
	RevisionInterfaceHash = hashref.MustParse("sha256:67086c90d3699fcb4683b01dae545f8bdb6406540516c9fa8f745de4f09a8b58")
	ManifestInterfaceHash = hashref.MustParse("sha256:91b9bf87219568c86d24deccaec783ce99975542b756d8e7a7d95468e4727f7e")
)

// Entry is one published app.
type Entry struct {
	Name        string   `json:"name"`
	Version     int64    `json:"version"`
	Manifest    string   `json:"manifest"`
	Transitions []string `json:"transitions"`
	PublishedAt string   `json:"published_at"`
}

// Revision is one immutable registry revision.
type Revision struct {
	Schema   string  `json:"schema"`
	Revision int64   `json:"revision"`
	Parent   string  `json:"parent,omitempty"`
	Apps     []Entry `json:"apps"`
}

// Lookup returns the entry for name.
func (r Revision) Lookup(name string) (Entry, bool) {
	for _, e := range r.Apps {
		if e.Name == name {
			return e, true
		}
	}
	return Entry{}, false
}

// Validate checks a revision's shape: schema, a positive revision, a parent exactly when
// revision > 1, at most maxApps entries sorted by name with unique names in the app
// grammar, positive versions, a manifest ref, sorted unique transition IDs, and an
// RFC 3339 UTC publication time.
func (r Revision) Validate() error {
	if r.Schema != RegistryName {
		return fmt.Errorf("app registry: schema is %q, not %s", r.Schema, RegistryName)
	}
	if !revisionParentOk(r.Revision, r.Parent != "") {
		return fmt.Errorf("app registry: revision %d with parent %q (only revision 1 has none)", r.Revision, r.Parent)
	}
	if r.Parent != "" {
		if _, err := hashref.Parse(r.Parent); err != nil {
			return fmt.Errorf("app registry: parent: %w", err)
		}
	}
	if r.Apps == nil {
		return errors.New("app registry: apps must be a list")
	}
	if !appCountOk(len(r.Apps)) {
		return fmt.Errorf("app registry: %d apps; limit is %d", len(r.Apps), maxApps)
	}
	for i, e := range r.Apps {
		if !appmanifest.NameOk(e.Name) {
			return fmt.Errorf("app registry: apps[%d] name %q is not an app name", i, e.Name)
		}
		if i > 0 && r.Apps[i-1].Name >= e.Name {
			return fmt.Errorf("app registry: apps must be sorted by unique name (%q then %q)", r.Apps[i-1].Name, e.Name)
		}
		if !versionOk(e.Version) {
			return fmt.Errorf("app registry: %s version %d is not positive", e.Name, e.Version)
		}
		if _, err := hashref.Parse(e.Manifest); err != nil {
			return fmt.Errorf("app registry: %s manifest: %w", e.Name, err)
		}
		if e.Transitions == nil {
			return fmt.Errorf("app registry: %s transitions must be a list", e.Name)
		}
		if !transitionCountOk(len(e.Transitions)) {
			return fmt.Errorf("app registry: %s lists %d transitions; limit is %d", e.Name, len(e.Transitions), maxTransitions)
		}
		for j, id := range e.Transitions {
			if id == "" || (j > 0 && e.Transitions[j-1] >= id) {
				return fmt.Errorf("app registry: %s transitions must be sorted, unique and non-empty", e.Name)
			}
		}
		if t, err := time.Parse(time.RFC3339, e.PublishedAt); err != nil || t.UTC().Format(time.RFC3339) != e.PublishedAt {
			return fmt.Errorf("app registry: %s published_at %q is not an RFC 3339 UTC time", e.Name, e.PublishedAt)
		}
	}
	return nil
}

// Encode is the canonical encoding: fixed key order, no insignificant whitespace, a
// trailing newline. It refuses an invalid revision.
func Encode(r Revision) ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode reads a revision strictly: unknown keys, trailing content, an invalid revision
// and any non-canonical encoding are refused, so equal revisions have equal hashes.
func Decode(payload []byte) (Revision, error) {
	if len(payload) > maxRevisionRaw {
		return Revision{}, fmt.Errorf("app registry: revision is %d bytes; limit is %d", len(payload), maxRevisionRaw)
	}
	d := json.NewDecoder(bytes.NewReader(payload))
	d.DisallowUnknownFields()
	var r Revision
	if err := d.Decode(&r); err != nil {
		return Revision{}, fmt.Errorf("app registry: decode: %w", err)
	}
	if _, err := d.Token(); !errors.Is(err, io.EOF) {
		return Revision{}, errors.New("app registry: content after the revision")
	}
	canon, err := Encode(r)
	if err != nil {
		return Revision{}, err
	}
	if !bytes.Equal(canon, payload) {
		return Revision{}, errors.New("app registry: revision is not in canonical form")
	}
	return r, nil
}

// BuildPublish returns the revision after publishing manifestRef as name's next version.
// current is nil when the registry has no head; currentHead is that head's ref. A
// publish whose manifest equals the published version's is refused: it would mint a
// version that changes nothing.
func BuildPublish(current *Revision, currentHead hashref.HashRef, name string, manifestRef hashref.HashRef, transitions []string, now time.Time) (Revision, error) {
	if !appmanifest.NameOk(name) {
		return Revision{}, fmt.Errorf("app registry: %q is not an app name", name)
	}
	ids := append([]string{}, transitions...)
	sort.Strings(ids)
	next := nextRevision(current, currentHead)
	previous := int64(0)
	var apps []Entry
	for _, e := range next.Apps {
		if e.Name != name {
			apps = append(apps, e)
			continue
		}
		if !publishChanges(e.Manifest == manifestRef.String()) {
			return Revision{}, fmt.Errorf("app registry: %s version %d already pins %s", name, e.Version, e.Manifest)
		}
		previous = e.Version
	}
	apps = append(apps, Entry{Name: name, Version: nextVersion(previous), Manifest: manifestRef.String(), Transitions: ids,
		PublishedAt: now.UTC().Truncate(time.Second).Format(time.RFC3339)})
	sort.Slice(apps, func(i, j int) bool { return apps[i].Name < apps[j].Name })
	next.Apps = apps
	return next, next.Validate()
}

// BuildRetire returns the revision after removing name. Its history stays in the
// parent chain.
func BuildRetire(current *Revision, currentHead hashref.HashRef, name string) (Revision, error) {
	if current == nil {
		return Revision{}, fmt.Errorf("app registry: no app %s to retire (the registry is empty)", name)
	}
	next := nextRevision(current, currentHead)
	apps := []Entry{}
	found := false
	for _, e := range next.Apps {
		if e.Name == name {
			found = true
			continue
		}
		apps = append(apps, e)
	}
	if !retireOk(found) {
		return Revision{}, fmt.Errorf("app registry: no app %s to retire", name)
	}
	next.Apps = apps
	return next, next.Validate()
}

func nextRevision(current *Revision, currentHead hashref.HashRef) Revision {
	if current == nil {
		return Revision{Schema: RegistryName, Revision: 1, Apps: []Entry{}}
	}
	apps := make([]Entry, len(current.Apps))
	copy(apps, current.Apps)
	return Revision{Schema: RegistryName, Revision: current.Revision + 1, Parent: currentHead.String(), Apps: apps}
}

// ObjectStore is the persistence boundary (the same narrow interface host/transitionreg
// uses), so store failures inject directly in tests.
type ObjectStore interface {
	GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error)
	GetObject(context.Context, hashref.HashRef) (store.Object, bool, error)
	PutObject(context.Context, store.Object) error
	CompareAndSetRegistryHead(context.Context, string, hashref.HashRef, hashref.HashRef) error
}

// Snapshot is one read of the registry: its head and decoded revision. Present is
// false when the registry has no head.
type Snapshot struct {
	Head     hashref.HashRef
	Present  bool
	Revision Revision
}

// Read reads and decodes the current head. The head is never cached.
func Read(ctx context.Context, s ObjectStore) (Snapshot, error) {
	head, ok, err := s.GetRegistryHead(ctx, RegistryName)
	if err != nil {
		return Snapshot{}, fmt.Errorf("app registry: read head: %w", err)
	}
	if !ok {
		return Snapshot{}, nil
	}
	r, err := readRevision(ctx, s, head)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Head: head, Present: true, Revision: r}, nil
}

func readRevision(ctx context.Context, s ObjectStore, ref hashref.HashRef) (Revision, error) {
	obj, ok, err := s.GetObject(ctx, ref)
	if err != nil {
		return Revision{}, fmt.Errorf("app registry: read revision %s: %w", ref, err)
	}
	if !ok {
		return Revision{}, fmt.Errorf("app registry: revision object %s is absent", ref)
	}
	if obj.SemanticID != RegistryName || obj.InterfaceHash != RevisionInterfaceHash {
		return Revision{}, fmt.Errorf("app registry: object %s is %s, not an app registry revision", ref, obj.SemanticID)
	}
	return Decode(obj.Payload)
}

// ReadManifest reads a pinned app.json object and re-validates it.
func ReadManifest(ctx context.Context, s ObjectStore, ref hashref.HashRef) (appmanifest.Manifest, []byte, error) {
	obj, ok, err := s.GetObject(ctx, ref)
	if err != nil {
		return appmanifest.Manifest{}, nil, fmt.Errorf("app registry: read manifest %s: %w", ref, err)
	}
	if !ok {
		return appmanifest.Manifest{}, nil, fmt.Errorf("app registry: manifest object %s is absent", ref)
	}
	if obj.SemanticID != ManifestSemanticID || obj.InterfaceHash != ManifestInterfaceHash {
		return appmanifest.Manifest{}, nil, fmt.Errorf("app registry: object %s is %s, not an app manifest", ref, obj.SemanticID)
	}
	m, err := appmanifest.Parse(obj.Payload)
	if err != nil {
		return appmanifest.Manifest{}, nil, err
	}
	if ps := appmanifest.Validate(m); len(ps) > 0 {
		return appmanifest.Manifest{}, nil, fmt.Errorf("app registry: pinned manifest %s no longer passes the law: %s", ref, ps[0])
	}
	return m, obj.Payload, nil
}

// ManifestObject wraps exact app.json bytes as a pinned object after checking them
// against the manifest law. The file's basename is not known here; the caller checks
// name agreement before publishing.
func ManifestObject(appJSON []byte) (store.Object, appmanifest.Manifest, error) {
	m, err := appmanifest.Parse(appJSON)
	if err != nil {
		return store.Object{}, appmanifest.Manifest{}, err
	}
	if ps := appmanifest.Validate(m); len(ps) > 0 {
		return store.Object{}, appmanifest.Manifest{}, fmt.Errorf("app registry: the manifest breaks the law: %s (and %d more)", ps[0], len(ps)-1)
	}
	return store.Object{Hash: hashref.SumSHA256(appJSON), InterfaceHash: ManifestInterfaceHash,
		SemanticID: ManifestSemanticID, Provenance: "host/appregistry", Payload: appJSON}, m, nil
}

// Publish writes next over expectedHead (zero when the registry has none). It checks the
// revision rules against the captured head, puts the manifest object (when given) and
// the revision object, then compare-and-sets the head. A lost race surfaces as the
// store's *RegistryCASConflict and changes no head.
func Publish(ctx context.Context, s ObjectStore, expectedHead hashref.HashRef, next Revision, manifest *store.Object) (hashref.HashRef, error) {
	if err := next.Validate(); err != nil {
		return hashref.HashRef{}, err
	}
	if expectedHead.IsZero() {
		if !nextRevisionOk(0, next.Revision) || next.Parent != "" {
			return hashref.HashRef{}, fmt.Errorf("app registry: an empty registry takes revision 1 with no parent, not revision %d", next.Revision)
		}
	} else {
		current, err := readRevision(ctx, s, expectedHead)
		if err != nil {
			return hashref.HashRef{}, err
		}
		if !nextRevisionOk(current.Revision, next.Revision) || next.Parent != expectedHead.String() {
			return hashref.HashRef{}, fmt.Errorf("app registry: revision %d over parent %q is not revision %d over the captured head %s",
				next.Revision, next.Parent, current.Revision+1, expectedHead)
		}
	}
	if manifest != nil {
		if manifest.SemanticID != ManifestSemanticID || manifest.InterfaceHash != ManifestInterfaceHash {
			return hashref.HashRef{}, errors.New("app registry: the manifest object is not a world/app-manifest/v1 object")
		}
		naming := 0
		for _, e := range next.Apps {
			if e.Manifest == manifest.Hash.String() {
				naming++
			}
		}
		if !manifestPinned(naming) {
			return hashref.HashRef{}, fmt.Errorf("app registry: no entry pins manifest %s", manifest.Hash)
		}
		if err := s.PutObject(ctx, *manifest); err != nil {
			return hashref.HashRef{}, fmt.Errorf("app registry: put manifest: %w", err)
		}
	}
	payload, err := Encode(next)
	if err != nil {
		return hashref.HashRef{}, err
	}
	ref := hashref.SumSHA256(payload)
	obj := store.Object{Hash: ref, InterfaceHash: RevisionInterfaceHash, SemanticID: RegistryName, Provenance: "host/appregistry", Payload: payload}
	if err := s.PutObject(ctx, obj); err != nil {
		return hashref.HashRef{}, fmt.Errorf("app registry: put revision: %w", err)
	}
	if err := s.CompareAndSetRegistryHead(ctx, RegistryName, expectedHead, ref); err != nil {
		return hashref.HashRef{}, fmt.Errorf("app registry: compare and set head: %w", err)
	}
	return ref, nil
}
